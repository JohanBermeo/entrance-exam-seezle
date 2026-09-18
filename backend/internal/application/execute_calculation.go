package application

import (
	"context"
	"fmt"
	"time"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/expression"
	"back-calculator/internal/domain/operators"
	"back-calculator/internal/engine"
)

// Request represents the input for a calculation.
type Request struct {
	Operations []calculation.Operation `json:"operations,omitempty"`
	Outputs    []string                `json:"outputs,omitempty"`
	Expression string                  `json:"expression,omitempty"`
}

// Response represents the output of a calculation.
type Response struct {
	RequestID    string              `json:"requestId"`
	Results      calculation.Results `json:"results"`
	Outputs      []string            `json:"outputs"`
	OutputValues []float64           `json:"outputValues"`
	OutputErrors []string            `json:"outputErrors,omitempty"`
	DurationMs   int64               `json:"durationMs"`
}

// Options carries the execution dependencies for a calculation.
type Options struct {
	Registry operators.Registry
	// Scheduler runs the validated DAG concurrently. Zero value executes
	// sequentially (a single worker slot).
	Scheduler engine.Scheduler
	// MaxNodes caps the operations compiled from an expression.
	MaxNodes int
	// MaxDepth caps the expression AST depth accepted by the compiler.
	MaxDepth int
}

// ExecuteCalculation executes a calculation request and returns the response.
func ExecuteCalculation(ctx context.Context, req Request, registry operators.Registry) (Response, error) {
	return ExecuteWithOptions(ctx, req, Options{
		Registry:  registry,
		Scheduler: engine.NewScheduler(registry, 1),
		MaxNodes:  100,
		MaxDepth:  50,
	})
}

// ExecuteWithOptions executes a calculation request with explicit execution options.
func ExecuteWithOptions(ctx context.Context, req Request, opts Options) (Response, error) {
	startTime := time.Now()
	requestID := generateRequestID()
	registry := opts.Registry

	// Validate request
	if req.Expression != "" && len(req.Operations) > 0 {
		return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "cannot provide both expression and operations")
	}
	if req.Expression == "" && len(req.Operations) == 0 {
		return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "must provide either expression or operations")
	}
	if len(req.Outputs) == 0 {
		return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "at least one output is required")
	}

	var graph *calculation.Graph
	var compileErrors []expression.ParseError

	if req.Expression != "" {
		if len(req.Outputs) != 1 {
			return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "expression mode requires exactly one output")
		}
		// Compile expression to DAG
		lexer := expression.NewLexer(req.Expression)
		tokens, err := lexer.Tokenize()
		if err != nil {
			return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "lexer error: "+err.Error())
		}

		parser := expression.NewParser(tokens)
		expr, parseErrors := parser.Parse()
		if len(parseErrors) > 0 {
			var msgs []string
			for _, e := range parseErrors {
				msgs = append(msgs, e.Error())
			}
			return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "parse error: "+fmt.Sprintf("%v", msgs))
		}

		compiler := expression.NewCompiler(registry)
		compiler.SetLimits(opts.MaxNodes, opts.MaxDepth)
		graph, compileErrors = compiler.Compile(expr, req.Outputs[0])
		if len(compileErrors) > 0 {
			var msgs []string
			for _, e := range compileErrors {
				msgs = append(msgs, e.Error())
			}
			return Response{}, calculation.NewDomainError(calculation.CodeInvalidInput, "", "compile error: "+fmt.Sprintf("%v", msgs))
		}
	} else {
		// Use explicit DAG
		knownOps := registry.ArityMap()
		var err error
		graph, err = calculation.NewGraph(req.Operations, req.Outputs, knownOps)
		if err != nil {
			return Response{}, err
		}
	}

	// Execute the graph with the concurrent fan-out/fan-in scheduler.
	// On failure the root error surfaces so the HTTP layer maps it to
	// 422 (domain rule) or 408/499 (deadline/cancel) per the plan.
	// Per-operation partial results will be formalized with the
	// OpenAPI contract in M05.
	results, err := opts.Scheduler.Execute(ctx, graph)
	if err != nil {
		return Response{}, err
	}

	return Response{
		RequestID:    requestID,
		Results:      results,
		Outputs:      req.Outputs,
		OutputValues: results.OutputValues(req.Outputs),
		DurationMs:   time.Since(startTime).Milliseconds(),
	}, nil
}

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

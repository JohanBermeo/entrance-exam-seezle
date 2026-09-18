package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/expression"
	"back-calculator/internal/domain/operators"
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

// ExecuteCalculation executes a calculation request and returns the response.
func ExecuteCalculation(ctx context.Context, req Request, registry operators.Registry) (Response, error) {
	startTime := time.Now()
	requestID := generateRequestID()

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
		compiler.SetLimits(100, 50)
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

	// Execute the graph (sequential in M03; concurrent scheduler arrives in M04).
	results, execErrors := executeGraph(ctx, graph, registry)
	if len(execErrors) > 0 {
		// Surface the root domain error so the HTTP layer maps it to
		// 422 (domain rule) or 408/499 (deadline/cancel) per the plan.
		// Per-operation partial results will be formalized with the
		// scheduler + OpenAPI contract in M04/M05.
		for _, execErr := range execErrors {
			var domainErr *calculation.DomainError
			if errors.As(execErr, &domainErr) {
				return Response{}, domainErr
			}
		}
		return Response{}, execErrors[0]
	}

	return Response{
		RequestID:    requestID,
		Results:      results,
		Outputs:      req.Outputs,
		OutputValues: results.OutputValues(req.Outputs),
		DurationMs:   time.Since(startTime).Milliseconds(),
	}, nil
}

func executeGraph(ctx context.Context, graph *calculation.Graph, registry operators.Registry) (calculation.Results, []error) {
	results := make(calculation.Results)
	completed := make(map[string]bool)
	var execErrors []error

	// Simple sequential execution for now (will be replaced with concurrent scheduler in M04)
	// Process operations in topological order using Kahn's algorithm
	remaining := graph.Len()
	for remaining > 0 {
		select {
		case <-ctx.Done():
			execErrors = append(execErrors, ctx.Err())
			return results, execErrors
		default:
		}

		progress := false
		for _, op := range graph.Operations {
			if completed[op.ID] {
				continue
			}

			// Check if all dependencies are met
			allMet := true
			for _, input := range op.Inputs {
				if input.IsRef() && !completed[*input.Ref] {
					allMet = false
					break
				}
			}

			if !allMet {
				continue
			}

			// Execute operation
			var inputs []float64
			for _, input := range op.Inputs {
				if input.IsLiteral() {
					inputs = append(inputs, *input.Value)
				} else {
					res, ok := results.Get(*input.Ref)
					if !ok || res.HasError() {
						// Dependency failed
						results.SetError(op.ID, "dependency "+*input.Ref+" failed")
						execErrors = append(execErrors, fmt.Errorf("dependency %s failed", *input.Ref))
					} else {
						inputs = append(inputs, res.Value)
					}
				}
			}

			if len(inputs) == len(op.Inputs) {
				value, err := registry.Evaluate(op.Op, inputs)
				if err != nil {
					results.SetError(op.ID, err.Error())
					execErrors = append(execErrors, err)
				} else {
					results.Set(op.ID, value)
				}
			}

			completed[op.ID] = true
			remaining--
			progress = true
		}

		if !progress {
			// Cycle or deadlock - should have been caught in validation
			execErrors = append(execErrors, fmt.Errorf("execution deadlock: no progress made"))
			break
		}
	}

	return results, execErrors
}

func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}

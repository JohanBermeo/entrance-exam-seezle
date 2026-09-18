package expression

import (
	"testing"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
)

func TestCompilerSimpleExpression(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("1 + 2")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	graph, errors := compiler.Compile(expr, "result")
	if len(errors) > 0 {
		t.Fatalf("Compiler errors: %v", errors)
	}

	if graph == nil {
		t.Fatal("expected graph, got nil")
	}

	// One compute node (bin) plus the "result" identity alias.
	if graph.Len() != 2 {
		t.Fatalf("expected 2 operations (compute + result alias), got %d", graph.Len())
	}

	// The compiler exposes the final node under the requested output name.
	if len(graph.Outputs) != 1 || graph.Outputs[0] != "result" {
		t.Fatalf("expected outputs [result], got %v", graph.Outputs)
	}

	out := graph.Operation("result")
	if out == nil || out.Op != "add" || len(out.Inputs) != 2 {
		t.Fatalf("expected result alias add/2, got %+v", out)
	}
}

func TestCompilerExpressionWithVariables(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("a + b")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	_, errors = compiler.Compile(expr, "result")
	if len(errors) == 0 {
		t.Fatal("expected error for unbound identifier")
	}
}

func TestCompilerFunctionCall(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("sqrt(16)")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	graph, errors := compiler.Compile(expr, "result")
	if len(errors) > 0 {
		t.Fatalf("Compiler errors: %v", errors)
	}

	if graph.Len() != 2 {
		t.Fatalf("expected 2 operations (sqrt + result alias), got %d", graph.Len())
	}

	op := graph.Operation(graph.Outputs[0])
	if op == nil {
		t.Fatal("output operation not found")
	}

	// Should have a sqrt operation with a literal input (literals are inlined)
	found := false
	for _, o := range graph.Operations {
		if o.Op == "sqrt" {
			found = true
			if len(o.Inputs) != 1 {
				t.Fatalf("sqrt should have 1 input, got %d", len(o.Inputs))
			}
			// Input should be a literal (16), not a reference
			if !o.Inputs[0].IsLiteral() {
				t.Fatal("sqrt input should be a literal")
			}
			if *o.Inputs[0].Value != 16 {
				t.Fatalf("expected literal 16, got %v", *o.Inputs[0].Value)
			}
		}
	}
	if !found {
		t.Fatal("sqrt operation not found")
	}
}

func TestCompilerPercentFunction(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("percent(200, 15)")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	graph, errors := compiler.Compile(expr, "result")
	if len(errors) > 0 {
		t.Fatalf("Compiler errors: %v", errors)
	}

	found := false
	for _, o := range graph.Operations {
		if o.Op == "percent" {
			found = true
			if len(o.Inputs) != 2 {
				t.Fatalf("percent should have 2 inputs, got %d", len(o.Inputs))
			}
			// Inputs should be literals (200, 15), not references
			if !o.Inputs[0].IsLiteral() || !o.Inputs[1].IsLiteral() {
				t.Fatal("percent inputs should be literals")
			}
			if *o.Inputs[0].Value != 200 || *o.Inputs[1].Value != 15 {
				t.Fatalf("expected literals 200, 15, got %v, %v", *o.Inputs[0].Value, *o.Inputs[1].Value)
			}
		}
	}
	if !found {
		t.Fatal("percent operation not found")
	}
}

func TestCompilerComplexExpression(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("sqrt(percent(200, 15)) + 4 ^ 2")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	graph, errors := compiler.Compile(expr, "result")
	if len(errors) > 0 {
		t.Fatalf("Compiler errors: %v", errors)
	}

	// With inlining: percent(200,15), sqrt(...), 4^2, add(...) = 4 operations
	if graph.Len() < 4 {
		t.Fatalf("expected at least 4 operations, got %d", graph.Len())
	}

	// Verify graph is valid (no cycles, all refs exist)
	for _, op := range graph.Operations {
		for _, input := range op.Inputs {
			if input.IsRef() {
				if graph.Operation(*input.Ref) == nil {
					t.Fatalf("operation %s references unknown operation %s", op.ID, *input.Ref)
				}
			}
		}
	}
}

func TestCompilerUnknownFunction(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("unknown(1)")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	_, errors = compiler.Compile(expr, "result")
	if len(errors) == 0 {
		t.Fatal("expected error for unknown function")
	}
}

func TestCompilerWrongArity(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("sqrt(1, 2)")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	_, errors = compiler.Compile(expr, "result")
	if len(errors) == 0 {
		t.Fatal("expected error for wrong arity")
	}
}

func TestCompilerMaxNodesLimit(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)
	compiler.SetLimits(5, 50) // Very low limit

	lexer := NewLexer("1 + 2 + 3 + 4 + 5 + 6 + 7")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	_, errors = compiler.Compile(expr, "result")
	if len(errors) == 0 {
		t.Fatal("expected error for max nodes exceeded")
	}
}

func TestCompilerMaxDepthLimit(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)
	compiler.SetLimits(100, 3) // Very low depth limit

	lexer := NewLexer("1 + (2 + (3 + (4 + 5)))")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	_, errors = compiler.Compile(expr, "result")
	if len(errors) == 0 {
		t.Fatal("expected error for max depth exceeded")
	}
}

func TestCompilerGeneratesValidGraph(t *testing.T) {
	registry := operators.NewRegistry()
	compiler := NewCompiler(registry)

	lexer := NewLexer("sqrt(percent(200, 15)) + 4 ^ 2")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	graph, errors := compiler.Compile(expr, "result")
	if len(errors) > 0 {
		t.Fatalf("Compiler errors: %v", errors)
	}

	// Verify graph can be validated with calculation package
	knownOps := registry.ArityMap()
	validatedGraph, err := calculation.NewGraph(graph.Operations, graph.Outputs, knownOps)
	if err != nil {
		t.Fatalf("Graph validation failed: %v", err)
	}

	if validatedGraph.Len() != graph.Len() {
		t.Fatalf("Graph length mismatch after validation")
	}
}

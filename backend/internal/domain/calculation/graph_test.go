package calculation_test

import (
	"errors"
	"testing"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
)

func knownOps() map[string]int {
	return operators.NewRegistry().ArityMap()
}

func lit(v float64) calculation.Input { return calculation.Literal(v) }
func ref(id string) calculation.Input { return calculation.Ref(id) }

func TestNewGraphAcceptsValidDAG(t *testing.T) {
	ops := []calculation.Operation{
		{ID: "sum", Op: "add", Inputs: []calculation.Input{lit(12), lit(8)}},
		{ID: "pow", Op: "power", Inputs: []calculation.Input{lit(4), lit(2)}},
		{ID: "total", Op: "multiply", Inputs: []calculation.Input{ref("sum"), ref("pow")}},
	}
	g, err := calculation.NewGraph(ops, []string{"total"}, knownOps())
	if err != nil {
		t.Fatalf("NewGraph error = %v", err)
	}
	if g.Len() != 3 {
		t.Fatalf("Len = %d, want 3", g.Len())
	}
	if deps := g.Dependencies("total"); len(deps) != 2 {
		t.Fatalf("Dependencies(total) = %v, want 2 entries", deps)
	}
	if deps := g.Dependents("sum"); len(deps) != 1 || deps[0] != "total" {
		t.Fatalf("Dependents(sum) = %v, want [total]", deps)
	}
	ready := g.ReadyOperations(map[string]bool{})
	if len(ready) != 2 {
		t.Fatalf("initial ready = %d, want 2 (sum, pow)", len(ready))
	}
}

func TestNewGraphRejectsInvalidGraphs(t *testing.T) {
	testCases := []struct {
		name    string
		ops     []calculation.Operation
		outputs []string
		code    calculation.ErrorCode
	}{
		{
			name:    "empty_operations",
			ops:     nil,
			outputs: []string{"total"},
			code:    calculation.CodeInvalidInput,
		},
		{
			name: "empty_outputs",
			ops: []calculation.Operation{
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(1), lit(2)}},
			},
			outputs: nil,
			code:    calculation.CodeInvalidInput,
		},
		{
			name: "duplicate_id",
			ops: []calculation.Operation{
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(1), lit(2)}},
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(3), lit(4)}},
			},
			outputs: []string{"a"},
			code:    calculation.CodeInvalidInput,
		},
		{
			name: "unknown_operation",
			ops: []calculation.Operation{
				{ID: "a", Op: "mod", Inputs: []calculation.Input{lit(1), lit(2)}},
			},
			outputs: []string{"a"},
			code:    calculation.CodeUnknownOperation,
		},
		{
			name: "invalid_arity",
			ops: []calculation.Operation{
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(1)}},
			},
			outputs: []string{"a"},
			code:    calculation.CodeInvalidArity,
		},
		{
			name: "unknown_output",
			ops: []calculation.Operation{
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(1), lit(2)}},
			},
			outputs: []string{"missing"},
			code:    calculation.CodeUnknownReference,
		},
		{
			name: "unknown_reference",
			ops: []calculation.Operation{
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(1), ref("missing")}},
			},
			outputs: []string{"a"},
			code:    calculation.CodeUnknownReference,
		},
		{
			name: "self_reference",
			ops: []calculation.Operation{
				{ID: "a", Op: "sqrt", Inputs: []calculation.Input{ref("a")}},
			},
			outputs: []string{"a"},
			code:    calculation.CodeCycleDetected,
		},
		{
			name: "cycle",
			ops: []calculation.Operation{
				{ID: "a", Op: "add", Inputs: []calculation.Input{lit(1), ref("b")}},
				{ID: "b", Op: "add", Inputs: []calculation.Input{lit(2), ref("a")}},
			},
			outputs: []string{"a"},
			code:    calculation.CodeCycleDetected,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := calculation.NewGraph(tc.ops, tc.outputs, knownOps())
			if err == nil {
				t.Fatal("NewGraph error = nil, want domain error")
			}
			var domainErr *calculation.DomainError
			if !errors.As(err, &domainErr) {
				t.Fatalf("error type = %T, want *calculation.DomainError", err)
			}
			if domainErr.Code != tc.code {
				t.Fatalf("error code = %q, want %q", domainErr.Code, tc.code)
			}
		})
	}
}

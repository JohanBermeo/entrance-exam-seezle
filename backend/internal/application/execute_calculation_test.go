package application_test

import (
	"context"
	"testing"

	"back-calculator/internal/application"
	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
	"back-calculator/internal/engine"
)

func testOptions(workers int) application.Options {
	registry := operators.NewRegistry()
	return application.Options{
		Registry:  registry,
		Scheduler: engine.NewScheduler(registry, workers),
		MaxNodes:  100,
		MaxDepth:  50,
	}
}

func TestExecuteWithOptionsDAG(t *testing.T) {
	for _, workers := range []int{1, 2, 8} {
		req := application.Request{
			Operations: []calculation.Operation{
				{ID: "sum", Op: "add", Inputs: []calculation.Input{calculation.Literal(12), calculation.Literal(8)}},
				{ID: "pow", Op: "power", Inputs: []calculation.Input{calculation.Literal(4), calculation.Literal(2)}},
				{ID: "total", Op: "multiply", Inputs: []calculation.Input{calculation.Ref("sum"), calculation.Ref("pow")}},
			},
			Outputs: []string{"total"},
		}
		resp, err := application.ExecuteWithOptions(context.Background(), req, testOptions(workers))
		if err != nil {
			t.Fatalf("workers=%d: Execute error = %v", workers, err)
		}
		if len(resp.OutputValues) != 1 || resp.OutputValues[0] != 320 {
			t.Fatalf("workers=%d: outputValues = %v, want [320]", workers, resp.OutputValues)
		}
	}
}

func TestExecuteWithOptionsExpression(t *testing.T) {
	req := application.Request{
		Expression: "sqrt(percent(200, 15)) + 4 ^ 2",
		Outputs:    []string{"result"},
	}
	resp, err := application.ExecuteWithOptions(context.Background(), req, testOptions(4))
	if err != nil {
		t.Fatalf("Execute error = %v", err)
	}
	if len(resp.OutputValues) != 1 {
		t.Fatalf("outputValues = %v, want single value", resp.OutputValues)
	}
	if got := resp.OutputValues[0]; got < 21.47 || got > 21.49 {
		t.Fatalf("outputValues = %v, want ≈21.477", resp.OutputValues)
	}
}

func TestExecuteWithOptionsHonorsCompilerLimits(t *testing.T) {
	opts := testOptions(2)
	opts.MaxNodes = 1
	req := application.Request{
		Expression: "1 + 2 + 3",
		Outputs:    []string{"result"},
	}
	if _, err := application.ExecuteWithOptions(context.Background(), req, opts); err == nil {
		t.Fatal("Execute error = nil, want max-nodes error")
	}
}

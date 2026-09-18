package engine

import (
	"fmt"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
)

// NodeExecutor runs a single operation over resolved float64 inputs.
// It is the seam that lets tests instrument execution (latency, concurrency,
// failures) without touching operators or the scheduler.
type NodeExecutor interface {
	Execute(op calculation.Operation, inputs []float64) (float64, error)
}

// RegistryExecutor is the production NodeExecutor backed by the operator registry.
type RegistryExecutor struct {
	Registry operators.Registry
}

// Execute resolves and runs the operator by its stable API name.
func (e RegistryExecutor) Execute(op calculation.Operation, inputs []float64) (float64, error) {
	return e.Registry.Evaluate(op.Op, inputs)
}

// resolveInputs maps an operation's inputs to values using completed results.
// It fails if a referenced result is missing or carries an error; the scheduler
// never launches a node in that state, so this is a safety net.
func resolveInputs(op *calculation.Operation, get func(string) (calculation.Result, bool)) ([]float64, error) {
	inputs := make([]float64, 0, len(op.Inputs))
	for _, input := range op.Inputs {
		if input.IsLiteral() {
			inputs = append(inputs, *input.Value)
			continue
		}
		res, ok := get(*input.Ref)
		if !ok {
			return nil, fmt.Errorf("dependency %s has no result", *input.Ref)
		}
		if res.HasError() {
			return nil, fmt.Errorf("dependency %s failed: %s", *input.Ref, res.Error)
		}
		inputs = append(inputs, res.Value)
	}
	return inputs, nil
}

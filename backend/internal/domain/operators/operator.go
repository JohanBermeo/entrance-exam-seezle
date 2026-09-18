package operators

import (
	"math"

	"back-calculator/internal/domain/calculation"
)

// Operator evaluates one mathematical operation over an ordered set of inputs.
type Operator interface {
	Name() string
	Arity() int
	Evaluate(inputs []float64) (float64, error)
}

func validateInputs(operation string, inputs []float64, arity int) error {
	if len(inputs) != arity {
		return calculation.NewDomainError(calculation.CodeInvalidArity, operation, "unexpected number of inputs")
	}
	for _, input := range inputs {
		if math.IsNaN(input) || math.IsInf(input, 0) {
			return calculation.NewDomainError(calculation.CodeInvalidInput, operation, "inputs must be finite numbers")
		}
	}
	return nil
}

func validateResult(operation string, result float64) (float64, error) {
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, calculation.NewDomainError(calculation.CodeNonFiniteNumber, operation, "result is not a finite number")
	}
	return result, nil
}

package operators

import (
	"math"

	"back-calculator/internal/domain/calculation"
)

type squareRoot struct{}

func (squareRoot) Name() string { return "sqrt" }
func (squareRoot) Arity() int   { return 1 }

func (squareRoot) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("sqrt", inputs, 1); err != nil {
		return 0, err
	}
	if inputs[0] < 0 {
		return 0, calculation.NewDomainError(calculation.CodeNegativeSquareRoot, "sqrt", "input must not be negative")
	}
	return validateResult("sqrt", math.Sqrt(inputs[0]))
}

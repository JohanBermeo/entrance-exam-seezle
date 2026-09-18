package operators

import "back-calculator/internal/domain/calculation"

type division struct{}

func (division) Name() string { return "divide" }
func (division) Arity() int   { return 2 }

func (division) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("divide", inputs, 2); err != nil {
		return 0, err
	}
	if inputs[1] == 0 {
		return 0, calculation.NewDomainError(calculation.CodeDivisionByZero, "divide", "divisor must not be zero")
	}
	return validateResult("divide", inputs[0]/inputs[1])
}

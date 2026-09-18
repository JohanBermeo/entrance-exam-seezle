package operators

import "math"

type power struct{}

func (power) Name() string { return "power" }
func (power) Arity() int { return 2 }

func (power) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("power", inputs, 2); err != nil {
		return 0, err
	}
	return validateResult("power", math.Pow(inputs[0], inputs[1]))
}

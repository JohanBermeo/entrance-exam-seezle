package operators

type subtraction struct{}

func (subtraction) Name() string { return "subtract" }

func (subtraction) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("subtract", inputs, 2); err != nil {
		return 0, err
	}
	return validateResult("subtract", inputs[0]-inputs[1])
}

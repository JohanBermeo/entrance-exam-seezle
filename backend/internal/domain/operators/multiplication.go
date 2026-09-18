package operators

type multiplication struct{}

func (multiplication) Name() string { return "multiply" }

func (multiplication) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("multiply", inputs, 2); err != nil {
		return 0, err
	}
	return validateResult("multiply", inputs[0]*inputs[1])
}

package operators

type percentage struct{}

func (percentage) Name() string { return "percent" }

func (percentage) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("percent", inputs, 2); err != nil {
		return 0, err
	}
	return validateResult("percent", inputs[0]*inputs[1]/100)
}

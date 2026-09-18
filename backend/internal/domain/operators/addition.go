package operators

type addition struct{}

func (addition) Name() string { return "add" }
func (addition) Arity() int   { return 2 }

func (addition) Evaluate(inputs []float64) (float64, error) {
	if err := validateInputs("add", inputs, 2); err != nil {
		return 0, err
	}
	return validateResult("add", inputs[0]+inputs[1])
}

package operators

import "back-calculator/internal/domain/calculation"

// Registry is the catalog of mathematical operators supported by the API.
type Registry struct {
	operators map[string]Operator
}

// NewRegistry builds the default operator catalog.
func NewRegistry() Registry {
	registered := []Operator{
		addition{}, subtraction{}, multiplication{}, division{}, power{}, squareRoot{}, percentage{},
	}
	operators := make(map[string]Operator, len(registered))
	for _, operator := range registered {
		operators[operator.Name()] = operator
	}
	return Registry{operators: operators}
}

// Evaluate resolves and runs an operator by its stable API name.
func (registry Registry) Evaluate(name string, inputs []float64) (float64, error) {
	operator, found := registry.operators[name]
	if !found {
		return 0, calculation.NewDomainError(calculation.CodeUnknownOperation, name, "operation is not supported")
	}
	return operator.Evaluate(inputs)
}

// ArityMap returns a map of operator names to their arities.
func (registry Registry) ArityMap() map[string]int {
	result := make(map[string]int, len(registry.operators))
	for name, op := range registry.operators {
		result[name] = op.Arity()
	}
	return result
}

package calculation

import (
	"math"
	"strconv"
)

// Input represents either a literal number or a reference to another operation's result.
type Input struct {
	Value *float64 `json:"value,omitempty"`
	Ref   *string  `json:"ref,omitempty"`
}

// IsLiteral returns true if the input is a literal number.
func (i Input) IsLiteral() bool {
	return i.Value != nil
}

// IsRef returns true if the input is a reference to another operation.
func (i Input) IsRef() bool {
	return i.Ref != nil
}

// Literal creates a literal input.
func Literal(v float64) Input {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		panic("literal input must be finite")
	}
	return Input{Value: &v}
}

// Ref creates a reference input.
func Ref(id string) Input {
	return Input{Ref: &id}
}

// Operation represents a single node in the calculation DAG.
type Operation struct {
	ID     string  `json:"id"`
	Op     string  `json:"op"`
	Inputs []Input `json:"inputs"`
}

// Validate checks the operation for basic structural validity.
func (o *Operation) Validate(knownOps map[string]int) error {
	if o.ID == "" {
		return NewDomainError(CodeInvalidInput, "", "operation id is required")
	}
	if o.Op == "" {
		return NewDomainError(CodeInvalidInput, "", "operation op is required")
	}
	arity, ok := knownOps[o.Op]
	if !ok {
		return NewDomainError(CodeUnknownOperation, o.Op, "operation is not supported")
	}
	if len(o.Inputs) != arity {
		return NewDomainError(CodeInvalidArity, o.Op, "unexpected number of inputs")
	}
	for i, input := range o.Inputs {
		if input.IsLiteral() && (math.IsNaN(*input.Value) || math.IsInf(*input.Value, 0)) {
			return NewDomainError(CodeInvalidInput, o.Op, "input at index "+strconv.Itoa(i)+" must be a finite number")
		}
		if input.IsRef() && *input.Ref == "" {
			return NewDomainError(CodeInvalidInput, o.Op, "reference id cannot be empty")
		}
		if !input.IsLiteral() && !input.IsRef() {
			return NewDomainError(CodeInvalidInput, o.Op, "input at index "+strconv.Itoa(i)+" must be a value or ref")
		}
	}
	return nil
}

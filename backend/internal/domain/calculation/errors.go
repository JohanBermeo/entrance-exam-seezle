package calculation

import "fmt"

// ErrorCode identifies a business-rule violation that a client can correct.
type ErrorCode string

const (
	CodeUnknownOperation   ErrorCode = "unknown_operation"
	CodeInvalidArity       ErrorCode = "invalid_arity"
	CodeDivisionByZero     ErrorCode = "division_by_zero"
	CodeNegativeSquareRoot ErrorCode = "negative_square_root"
	CodeNonFiniteNumber    ErrorCode = "non_finite_number"
)

// DomainError describes an invalid mathematical operation without exposing transport details.
type DomainError struct {
	Code      ErrorCode
	Operation string
	Message   string
}

func (error *DomainError) Error() string {
	if error.Operation == "" {
		return fmt.Sprintf("%s: %s", error.Code, error.Message)
	}
	return fmt.Sprintf("%s for %s: %s", error.Code, error.Operation, error.Message)
}

// NewDomainError creates a typed domain error for an operation.
func NewDomainError(code ErrorCode, operation, message string) *DomainError {
	return &DomainError{Code: code, Operation: operation, Message: message}
}

package operators_test

import (
	"errors"
	"math"
	"testing"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
)

func TestRegistryEvaluatesSupportedOperations(t *testing.T) {
	registry := operators.NewRegistry()

	testCases := []struct {
		name   string
		inputs []float64
		want   float64
	}{
		{name: "add", inputs: []float64{1.25, 2.75}, want: 4},
		{name: "subtract", inputs: []float64{10, 3}, want: 7},
		{name: "multiply", inputs: []float64{7, 6}, want: 42},
		{name: "divide", inputs: []float64{7, 2}, want: 3.5},
		{name: "power", inputs: []float64{2, 3}, want: 8},
		{name: "sqrt", inputs: []float64{81}, want: 9},
		{name: "percent", inputs: []float64{200, 15}, want: 30},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := registry.Evaluate(testCase.name, testCase.inputs)
			if err != nil {
				t.Fatalf("Evaluate(%q, %v) error = %v", testCase.name, testCase.inputs, err)
			}
			if got != testCase.want {
				t.Fatalf("Evaluate(%q, %v) = %v, want %v", testCase.name, testCase.inputs, got, testCase.want)
			}
		})
	}
}

func TestRegistryRejectsInvalidOperations(t *testing.T) {
	registry := operators.NewRegistry()

	testCases := []struct {
		name     string
		operator string
		inputs   []float64
		code     calculation.ErrorCode
	}{
		{name: "unknown_operation", operator: "unknown", inputs: []float64{1}, code: calculation.CodeUnknownOperation},
		{name: "invalid_arity", operator: "add", inputs: []float64{1}, code: calculation.CodeInvalidArity},
		{name: "division_by_zero", operator: "divide", inputs: []float64{1, 0}, code: calculation.CodeDivisionByZero},
		{name: "negative_square_root", operator: "sqrt", inputs: []float64{-1}, code: calculation.CodeNegativeSquareRoot},
		{name: "infinity_input", operator: "multiply", inputs: []float64{math.Inf(1), 2}, code: calculation.CodeInvalidInput},
		{name: "nan_input", operator: "multiply", inputs: []float64{math.NaN(), 2}, code: calculation.CodeInvalidInput},
		{name: "non_finite_number_result", operator: "power", inputs: []float64{0, -1}, code: calculation.CodeNonFiniteNumber},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := registry.Evaluate(testCase.operator, testCase.inputs)
			if err == nil {
				t.Fatalf("Evaluate(%q, %v) error = nil", testCase.operator, testCase.inputs)
			}

			var domainError *calculation.DomainError
			if !errors.As(err, &domainError) {
				t.Fatalf("error type = %T, want *calculation.DomainError", err)
			}
			if domainError.Code != testCase.code {
				t.Fatalf("error code = %q, want %q", domainError.Code, testCase.code)
			}
		})
	}
}

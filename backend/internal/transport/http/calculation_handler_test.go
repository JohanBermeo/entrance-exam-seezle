package http

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postCalculations(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	router := NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := httptest.NewRequest(http.MethodPost, "/v1/calculations", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestCalculationsDAG(t *testing.T) {
	response := postCalculations(t, `{
		"operations": [
			{"id": "sum", "op": "add", "inputs": [{"value": 12}, {"value": 8}]},
			{"id": "pow", "op": "power", "inputs": [{"value": 4}, {"value": 2}]},
			{"id": "total", "op": "multiply", "inputs": [{"ref": "sum"}, {"ref": "pow"}]}
		],
		"outputs": ["total"]
	}`)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", response.Code, http.StatusOK, response.Body.String())
	}
	var payload struct {
		Outputs      []string  `json:"outputs"`
		OutputValues []float64 `json:"outputValues"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(payload.OutputValues) != 1 || payload.OutputValues[0] != 320 {
		t.Fatalf("outputValues = %v, want [320]", payload.OutputValues)
	}
}

func TestCalculationsExpression(t *testing.T) {
	response := postCalculations(t, `{"expression": "sqrt(percent(200, 15)) + 4 ^ 2", "outputs": ["result"]}`)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", response.Code, http.StatusOK, response.Body.String())
	}
	var payload struct {
		OutputValues []float64 `json:"outputValues"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(payload.OutputValues) != 1 {
		t.Fatalf("outputValues = %v, want single value", payload.OutputValues)
	}
	// sqrt(30) + 16 ≈ 21.477
	if got := payload.OutputValues[0]; got < 21.47 || got > 21.49 {
		t.Fatalf("outputValues = %v, want ≈21.477", payload.OutputValues)
	}
}

func TestCalculationsErrors(t *testing.T) {
	testCases := []struct {
		name   string
		body   string
		status int
	}{
		{
			name:   "invalid_json",
			body:   `{"operations": [`,
			status: http.StatusBadRequest,
		},
		{
			name:   "both_modes",
			body:   `{"expression": "1+1", "operations": [{"id": "a", "op": "add", "inputs": [{"value": 1}, {"value": 1}]}], "outputs": ["a"]}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "unknown_operation",
			body:   `{"operations": [{"id": "a", "op": "mod", "inputs": [{"value": 1}, {"value": 1}]}], "outputs": ["a"]}`,
			status: http.StatusUnprocessableEntity,
		},
		{
			name:   "unknown_reference",
			body:   `{"operations": [{"id": "a", "op": "add", "inputs": [{"value": 1}, {"ref": "missing"}]}], "outputs": ["a"]}`,
			status: http.StatusUnprocessableEntity,
		},
		{
			name:   "cycle",
			body:   `{"operations": [{"id": "a", "op": "add", "inputs": [{"value": 1}, {"ref": "b"}]}, {"id": "b", "op": "add", "inputs": [{"value": 2}, {"ref": "a"}]}], "outputs": ["a"]}`,
			status: http.StatusUnprocessableEntity,
		},
		{
			name:   "division_by_zero",
			body:   `{"operations": [{"id": "a", "op": "divide", "inputs": [{"value": 1}, {"value": 0}]}], "outputs": ["a"]}`,
			status: http.StatusUnprocessableEntity,
		},
		{
			name:   "expression_syntax_error",
			body:   `{"expression": "1 +", "outputs": ["result"]}`,
			status: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := postCalculations(t, tc.body)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d (body: %s)", response.Code, tc.status, response.Body.String())
			}
		})
	}
}

// Package integration exercises the full HTTP stack: router, middlewares,
// calculation handler, application, scheduler and operators.
package integration_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"back-calculator/internal/application"
	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
	"back-calculator/internal/engine"
	transporthttp "back-calculator/internal/transport/http"
)

func testRouter(timeout time.Duration, maxBody int64) http.Handler {
	registry := operators.NewRegistry()
	return transporthttp.NewRouterWithOptions(slog.New(slog.NewTextHandler(io.Discard, nil)), application.Options{
		Registry:  registry,
		Scheduler: engine.NewScheduler(registry, 4),
		MaxNodes:  100,
		MaxDepth:  50,
	}, timeout, maxBody)
}

func post(t *testing.T, router http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/v1/calculations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func TestIntegrationValidDAG(t *testing.T) {
	resp := post(t, testRouter(30*time.Second, 1<<20), `{
		"operations": [
			{"id": "sum", "op": "add", "inputs": [{"value": 12}, {"value": 8}]},
			{"id": "pow", "op": "power", "inputs": [{"value": 4}, {"value": 2}]},
			{"id": "total", "op": "multiply", "inputs": [{"ref": "sum"}, {"ref": "pow"}]}
		],
		"outputs": ["total"]
	}`)
	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", resp.Code, resp.Body.String())
	}
	var payload struct {
		RequestID    string                     `json:"requestId"`
		Results      map[string]json.RawMessage `json:"results"`
		Outputs      []string                   `json:"outputs"`
		OutputValues []float64                  `json:"outputValues"`
		DurationMs   int64                      `json:"durationMs"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("contract decode error: %v", err)
	}
	if payload.RequestID == "" || len(payload.Results) != 3 {
		t.Fatalf("contract fields missing: %+v", payload)
	}
	if len(payload.Outputs) != 1 || payload.Outputs[0] != "total" {
		t.Fatalf("outputs = %v, want [total]", payload.Outputs)
	}
	if len(payload.OutputValues) != 1 || payload.OutputValues[0] != 320 {
		t.Fatalf("outputValues = %v, want [320]", payload.OutputValues)
	}
}

func TestIntegrationValidExpression(t *testing.T) {
	resp := post(t, testRouter(30*time.Second, 1<<20),
		`{"expression": "sqrt(percent(200, 15)) + 4 ^ 2", "outputs": ["result"]}`)
	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", resp.Code, resp.Body.String())
	}
	var payload struct {
		OutputValues []float64 `json:"outputValues"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("contract decode error: %v", err)
	}
	if len(payload.OutputValues) != 1 || payload.OutputValues[0] < 21.47 || payload.OutputValues[0] > 21.49 {
		t.Fatalf("outputValues = %v, want ≈[21.477]", payload.OutputValues)
	}
}

func TestIntegrationErrorContract(t *testing.T) {
	router := testRouter(30*time.Second, 1<<20)
	cases := []struct {
		name   string
		body   string
		status int
		code   string
	}{
		{"invalid_json", `{"operations": [`, 400, "invalid_input"},
		{"unknown_operation", `{"operations": [{"id": "a", "op": "mod", "inputs": [{"value": 1}]}], "outputs": ["a"]}`, 422, "unknown_operation"},
		{"division_by_zero", `{"operations": [{"id": "a", "op": "divide", "inputs": [{"value": 1}, {"value": 0}]}], "outputs": ["a"]}`, 422, "division_by_zero"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp := post(t, router, tc.body)
			if resp.Code != tc.status {
				t.Fatalf("status = %d, want %d (body: %s)", resp.Code, tc.status, resp.Body.String())
			}
			var payload struct {
				Code      string `json:"code"`
				Operation string `json:"operation"`
				Message   string `json:"message"`
			}
			if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
				t.Fatalf("contract decode error: %v", err)
			}
			if payload.Code != tc.code || payload.Message == "" {
				t.Fatalf("error contract = %+v, want code %q", payload, tc.code)
			}
			if strings.Contains(resp.Body.String(), `"Code"`) {
				t.Fatalf("body uses legacy capitalized keys: %s", resp.Body.String())
			}
		})
	}
}

func TestIntegrationPayloadTooLarge(t *testing.T) {
	resp := post(t, testRouter(30*time.Second, 10),
		`{"expression": "1 + 2", "outputs": ["result"]}`)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body: %s)", resp.Code, resp.Body.String())
	}
}

func TestIntegrationTimeout(t *testing.T) {
	registry := operators.NewRegistry()
	router := transporthttp.NewRouterWithOptions(slog.New(slog.NewTextHandler(io.Discard, nil)), application.Options{
		Registry:  registry,
		Scheduler: slowScheduler(registry),
		MaxNodes:  100,
		MaxDepth:  50,
	}, 20*time.Millisecond, 1<<20)
	resp := post(t, router, `{"expression": "1 + 2 + 3 + 4", "outputs": ["result"]}`)
	if resp.Code != http.StatusRequestTimeout {
		t.Fatalf("status = %d, want 408 (body: %s)", resp.Code, resp.Body.String())
	}
}

func TestIntegrationCancelledClient(t *testing.T) {
	router := testRouter(30*time.Second, 1<<20)
	req := httptest.NewRequest(http.MethodPost, "/v1/calculations", strings.NewReader(`{"expression": "1 + 1", "outputs": ["result"]}`))
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req.WithContext(ctx))
	if resp.Code != 499 {
		t.Fatalf("status = %d, want 499 (body: %s)", resp.Code, resp.Body.String())
	}
}

// slowScheduler returns a scheduler whose executor sleeps past short deadlines.
func slowScheduler(registry operators.Registry) engine.Scheduler {
	return engine.Scheduler{Executor: sleepExecutor{Registry: registry}, MaxWorkers: 2}
}

type sleepExecutor struct {
	Registry operators.Registry
}

func (e sleepExecutor) Execute(op calculation.Operation, inputs []float64) (float64, error) {
	time.Sleep(200 * time.Millisecond)
	return e.Registry.Evaluate(op.Op, inputs)
}

func TestIntegrationMetricsReflectTraffic(t *testing.T) {
	router := testRouter(30*time.Second, 1<<20)
	post(t, router, `{"expression": "1 + 1", "outputs": ["result"]}`)
	post(t, router, `{"operations": [`)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	var payload struct {
		Routes map[string]struct {
			Count    int64            `json:"count"`
			ByStatus map[string]int64 `json:"byStatus"`
		} `json:"routes"`
		InFlight int64 `json:"inFlight"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	route, ok := payload.Routes["POST /v1/calculations"]
	if !ok || route.Count != 2 {
		t.Fatalf("routes = %+v, want 2 POST /v1/calculations", payload.Routes)
	}
	if route.ByStatus["200"] != 1 || route.ByStatus["400"] != 1 {
		t.Fatalf("byStatus = %v, want map[200:1 400:1]", route.ByStatus)
	}
}

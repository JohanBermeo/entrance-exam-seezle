package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"back-calculator/internal/application"
	"back-calculator/internal/domain/operators"
	"back-calculator/internal/engine"
	"back-calculator/internal/platform/metrics"
)

// NewRouter returns the public HTTP surface for the calculator API
// with default execution options.
func NewRouter(logger *slog.Logger) http.Handler {
	registry := operators.NewRegistry()
	return NewRouterWithOptions(logger, application.Options{
		Registry:  registry,
		Scheduler: engine.NewScheduler(registry, 8),
		MaxNodes:  100,
		MaxDepth:  50,
	}, 30*time.Second, 1<<20)
}

// NewRouterWithOptions returns the HTTP surface with explicit execution
// options, request timeout and payload size limit.
func NewRouterWithOptions(logger *slog.Logger, options application.Options, timeout time.Duration, maxBody int64) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /readyz", readyHandler)

	calcHandler := NewCalculationHandlerWithOptions(logger, options, timeout, maxBody)
	mux.Handle("POST /v1/calculations", calcHandler)

	recorder := metrics.NewRecorder()
	mux.Handle("GET /metrics", recorder.Handler())

	return cors(requestID(recoverPanic(logger, logRequests(logger, recorder.Track(mux)))))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, statusResponse{Status: "ok"})
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, statusResponse{Status: "ready"})
}

type statusResponse struct {
	Status string `json:"status"`
}

func writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("HTTP request completed", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(startedAt).Milliseconds(), "request_id", requestIDFromContext(r.Context()))
	})
}

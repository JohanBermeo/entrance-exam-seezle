package http

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoints(t *testing.T) {
	router := NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))

	for _, testCase := range []struct {
		path string
		body string
	}{
		{path: "/healthz", body: `{"status":"ok"}` + "\n"},
		{path: "/readyz", body: `{"status":"ready"}` + "\n"},
	} {
		t.Run(testCase.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, testCase.path, nil)
			request.Header.Set("X-Request-ID", "test-request")
			response := httptest.NewRecorder()

			router.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			if response.Body.String() != testCase.body {
				t.Fatalf("body = %q, want %q", response.Body.String(), testCase.body)
			}
			if response.Header().Get("X-Request-ID") != "test-request" {
				t.Fatal("X-Request-ID header was not preserved")
			}
			if !strings.Contains(response.Header().Get("Content-Type"), "application/json") {
				t.Fatal("Content-Type must be JSON")
			}
		})
	}
}

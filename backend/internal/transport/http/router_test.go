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

func TestCORSHeaders(t *testing.T) {
	router := NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))

	t.Run("preflight options request", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodOptions, "/v1/calculations", nil)
		request.Header.Set("Origin", "http://localhost:3000")
		request.Header.Set("Access-Control-Request-Method", "POST")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
		}
		if response.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Fatalf("Access-Control-Allow-Origin = %q, want '*'", response.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("post request cors headers", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/v1/calculations", strings.NewReader(`{"expression":"1+1"}`))
		request.Header.Set("Origin", "http://localhost:3000")
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Fatalf("Access-Control-Allow-Origin = %q, want '*'", response.Header().Get("Access-Control-Allow-Origin"))
		}
	})
}


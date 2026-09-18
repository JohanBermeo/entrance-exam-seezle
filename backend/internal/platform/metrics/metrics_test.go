package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"back-calculator/internal/platform/metrics"
)

func TestRecorderCountsByRouteAndStatus(t *testing.T) {
	recorder := metrics.NewRecorder()
	handler := recorder.Track(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/boom" {
			w.WriteHeader(http.StatusTeapot)
			return
		}
		time.Sleep(time.Millisecond)
	}))

	for _, path := range []string{"/ok", "/ok", "/boom"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		recorder.Handler() // exercise handler construction; no effect
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}

	snapshot := recorder.Snapshot()
	if snapshot.InFlight != 0 {
		t.Fatalf("InFlight = %d, want 0", snapshot.InFlight)
	}
	ok, found := snapshot.Routes["GET /ok"]
	if !found || ok.Count != 2 {
		t.Fatalf("GET /ok = %+v, want count 2", ok)
	}
	if ok.ByStatus[http.StatusOK] != 2 {
		t.Fatalf("GET /ok byStatus = %v, want map[200:2]", ok.ByStatus)
	}
	if ok.AvgDurationMs < 1 {
		t.Fatalf("GET /ok avgDurationMs = %v, want >= 1ms", ok.AvgDurationMs)
	}
	boom, found := snapshot.Routes["GET /boom"]
	if !found || boom.ByStatus[http.StatusTeapot] != 1 {
		t.Fatalf("GET /boom = %+v, want byStatus map[418:1]", boom)
	}
}

func TestHandlerServesJSONSnapshot(t *testing.T) {
	recorder := metrics.NewRecorder()
	tracked := recorder.Track(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
	tracked.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/v1/calculations", nil))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	resp := httptest.NewRecorder()
	recorder.Handler().ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.Code)
	}
	var payload struct {
		Routes map[string]struct {
			Count int64 `json:"count"`
		} `json:"routes"`
		InFlight int64 `json:"inFlight"`
	}
	if err := json.Unmarshal(resp.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Routes["POST /v1/calculations"].Count != 1 {
		t.Fatalf("payload = %s, want 1 POST /v1/calculations", resp.Body.String())
	}
}

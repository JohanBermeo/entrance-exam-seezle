// Package metrics provides in-process HTTP instrumentation with no external
// dependencies: per-route counters, status breakdown, average latency and an
// in-flight gauge, exposed as JSON on GET /metrics.
package metrics

import (
	"encoding/json"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Recorder collects HTTP metrics. The zero value is not usable; use NewRecorder.
type Recorder struct {
	mu       sync.Mutex
	routes   map[string]*routeStats
	inFlight atomic.Int64
}

type routeStats struct {
	count         int64
	byStatus      map[int]int64
	totalDuration time.Duration
}

// NewRecorder builds an empty Recorder.
func NewRecorder() *Recorder {
	return &Recorder{routes: make(map[string]*routeStats)}
}

// Track wraps next, recording one observation per request keyed by
// "METHOD path". All API paths are fixed (no parameters), so the raw path
// is a stable route key.
func (r *Recorder) Track(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.inFlight.Add(1)
		defer r.inFlight.Add(-1)
		startedAt := time.Now()
		captured := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(captured, req)
		r.observe(req.Method+" "+req.URL.Path, captured.status, time.Since(startedAt))
	})
}

func (r *Recorder) observe(route string, status int, duration time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stats, ok := r.routes[route]
	if !ok {
		stats = &routeStats{byStatus: make(map[int]int64)}
		r.routes[route] = stats
	}
	stats.count++
	stats.byStatus[status]++
	stats.totalDuration += duration
}

// RouteMetrics is the JSON view of one route's stats.
type RouteMetrics struct {
	Count         int64          `json:"count"`
	ByStatus      map[int]int64 `json:"byStatus"`
	AvgDurationMs float64        `json:"avgDurationMs"`
}

// Snapshot is the JSON view served on GET /metrics.
type Snapshot struct {
	Routes   map[string]RouteMetrics `json:"routes"`
	InFlight int64                   `json:"inFlight"`
}

// Snapshot copies current metrics. Safe for concurrent use.
func (r *Recorder) Snapshot() Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	snapshot := Snapshot{
		Routes:   make(map[string]RouteMetrics, len(r.routes)),
		InFlight: r.inFlight.Load(),
	}
	for route, stats := range r.routes {
		byStatus := make(map[int]int64, len(stats.byStatus))
		for status, count := range stats.byStatus {
			byStatus[status] = count
		}
		var avg float64
		if stats.count > 0 {
			avg = float64(stats.totalDuration.Milliseconds()) / float64(stats.count)
		}
		snapshot.Routes[route] = RouteMetrics{
			Count:         stats.count,
			ByStatus:      byStatus,
			AvgDurationMs: avg,
		}
	}
	return snapshot
}

// Handler serves the metrics snapshot as JSON.
func (r *Recorder) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(r.Snapshot())
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

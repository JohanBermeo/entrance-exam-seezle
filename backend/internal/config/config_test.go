package config

import (
	"log/slog"
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("HTTP_READ_HEADER_TIMEOUT", "")
	t.Setenv("HTTP_READ_TIMEOUT", "")
	t.Setenv("HTTP_WRITE_TIMEOUT", "")
	t.Setenv("HTTP_IDLE_TIMEOUT", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")
	t.Setenv("CALC_MAX_WORKERS", "")
	t.Setenv("CALC_TIMEOUT", "")
	t.Setenv("CALC_MAX_BODY_BYTES", "")
	t.Setenv("CALC_MAX_NODES", "")
	t.Setenv("CALC_MAX_DEPTH", "")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.Port != "8080" || config.Environment != "development" || config.LogLevel != slog.LevelInfo {
		t.Fatalf("Load() = %+v, want defaults", config)
	}
	if config.ShutdownTimeout != 10*time.Second {
		t.Fatalf("ShutdownTimeout = %v, want 10s", config.ShutdownTimeout)
	}
	if config.MaxWorkers != 8 {
		t.Fatalf("MaxWorkers = %d, want 8", config.MaxWorkers)
	}
	if config.CalculationTimeout != 30*time.Second {
		t.Fatalf("CalculationTimeout = %v, want 30s", config.CalculationTimeout)
	}
	if config.MaxBodyBytes != 1<<20 {
		t.Fatalf("MaxBodyBytes = %d, want 1MiB", config.MaxBodyBytes)
	}
	if config.MaxNodes != 100 || config.MaxDepth != 50 {
		t.Fatalf("MaxNodes/MaxDepth = %d/%d, want 100/50", config.MaxNodes, config.MaxDepth)
	}
}

func TestLoadRejectsInvalidConfiguration(t *testing.T) {
	t.Setenv("PORT", "70000")
	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want invalid port error")
	}
}

func TestLoadRejectsInvalidEngineLimits(t *testing.T) {
	for _, env := range []string{"CALC_MAX_WORKERS", "CALC_MAX_BODY_BYTES", "CALC_MAX_NODES", "CALC_MAX_DEPTH"} {
		t.Run(env, func(t *testing.T) {
			t.Setenv(env, "0")
			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil for %s=0", env)
			}
			t.Setenv(env, "abc")
			if _, err := Load(); err == nil {
				t.Fatalf("Load() error = nil for %s=abc", env)
			}
		})
	}
	t.Run("CALC_TIMEOUT", func(t *testing.T) {
		t.Setenv("CALC_TIMEOUT", "-1s")
		if _, err := Load(); err == nil {
			t.Fatal("Load() error = nil for negative CALC_TIMEOUT")
		}
	})
}

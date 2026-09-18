package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPort              = "8080"
	defaultEnvironment       = "development"
	defaultLogLevel          = "info"
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 15 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
	defaultMaxWorkers        = 8
	defaultCalcTimeout       = 30 * time.Second
	defaultMaxBodyBytes      = int64(1 << 20)
	defaultMaxNodes          = 100
	defaultMaxDepth          = 50
)

// Config contains the runtime settings required by the HTTP server.
type Config struct {
	Port              string
	Environment       string
	LogLevel          slog.Level
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	// MaxWorkers bounds concurrent node executions per calculation request.
	MaxWorkers int
	// CalculationTimeout is the deadline applied to POST /v1/calculations.
	CalculationTimeout time.Duration
	// MaxBodyBytes caps the calculation request payload size.
	MaxBodyBytes int64
	// MaxNodes caps the operations compiled from an expression.
	MaxNodes int
	// MaxDepth caps the expression AST depth accepted by the compiler.
	MaxDepth int
}

// Load reads configuration from environment variables and applies safe defaults.
func Load() (Config, error) {
	port, err := loadPort()
	if err != nil {
		return Config{}, err
	}

	level, err := loadLogLevel()
	if err != nil {
		return Config{}, err
	}

	readHeaderTimeout, err := loadDuration("HTTP_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout)
	if err != nil {
		return Config{}, err
	}
	readTimeout, err := loadDuration("HTTP_READ_TIMEOUT", defaultReadTimeout)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := loadDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := loadDuration("HTTP_IDLE_TIMEOUT", defaultIdleTimeout)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := loadDuration("SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		return Config{}, err
	}
	maxWorkers, err := loadPositiveInt("CALC_MAX_WORKERS", defaultMaxWorkers)
	if err != nil {
		return Config{}, err
	}
	calcTimeout, err := loadDuration("CALC_TIMEOUT", defaultCalcTimeout)
	if err != nil {
		return Config{}, err
	}
	maxBodyBytes, err := loadPositiveInt64("CALC_MAX_BODY_BYTES", defaultMaxBodyBytes)
	if err != nil {
		return Config{}, err
	}
	maxNodes, err := loadPositiveInt("CALC_MAX_NODES", defaultMaxNodes)
	if err != nil {
		return Config{}, err
	}
	maxDepth, err := loadPositiveInt("CALC_MAX_DEPTH", defaultMaxDepth)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:               port,
		Environment:        valueOrDefault("APP_ENV", defaultEnvironment),
		LogLevel:           level,
		ReadHeaderTimeout:  readHeaderTimeout,
		ReadTimeout:        readTimeout,
		WriteTimeout:       writeTimeout,
		IdleTimeout:        idleTimeout,
		ShutdownTimeout:    shutdownTimeout,
		MaxWorkers:         maxWorkers,
		CalculationTimeout: calcTimeout,
		MaxBodyBytes:       maxBodyBytes,
		MaxNodes:           maxNodes,
		MaxDepth:           maxDepth,
	}, nil
}

func loadPort() (string, error) {
	port := valueOrDefault("PORT", defaultPort)
	parsed, err := strconv.Atoi(port)
	if err != nil || parsed < 1 || parsed > 65535 {
		return "", fmt.Errorf("PORT must be an integer between 1 and 65535")
	}
	return port, nil
}

func loadLogLevel() (slog.Level, error) {
	levelName := strings.ToLower(valueOrDefault("LOG_LEVEL", defaultLogLevel))
	var level slog.Level
	if err := level.UnmarshalText([]byte(levelName)); err != nil {
		return 0, fmt.Errorf("LOG_LEVEL must be debug, info, warn, or error: %w", err)
	}
	return level, nil
}

func loadDuration(name string, fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback, nil
	}
	duration, err := time.ParseDuration(raw)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", name)
	}
	return duration, nil
}

func loadPositiveInt(name string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return value, nil
}

func loadPositiveInt64(name string, fallback int64) (int64, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}
	return value, nil
}

func valueOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

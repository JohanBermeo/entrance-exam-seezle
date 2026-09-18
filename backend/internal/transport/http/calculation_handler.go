package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"back-calculator/internal/application"
	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
)

// CalculationHandler handles POST /v1/calculations requests.
type CalculationHandler struct {
	logger   *slog.Logger
	registry operators.Registry
}

// NewCalculationHandler creates a new calculation handler.
func NewCalculationHandler(logger *slog.Logger, registry operators.Registry) *CalculationHandler {
	return &CalculationHandler{
		logger:   logger,
		registry: registry,
	}
}

// ServeHTTP handles the calculation request.
func (h *CalculationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	var req application.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, calculation.CodeInvalidInput, "invalid JSON: "+err.Error())
		return
	}

	resp, err := application.ExecuteCalculation(ctx, req, h.registry)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeError(w, http.StatusRequestTimeout, calculation.CodeInvalidInput, "deadline exceeded")
			return
		}
		if errors.Is(err, context.Canceled) {
			h.writeError(w, 499, calculation.CodeInvalidInput, "client closed request")
			return
		}
		var domainErr *calculation.DomainError
		if errors.As(err, &domainErr) {
			status := h.domainErrorToStatus(domainErr.Code)
			h.writeError(w, status, domainErr.Code, domainErr.Message)
			return
		}
		h.logger.Error("unexpected error", "error", err, "request_id", resp.RequestID)
		h.writeError(w, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *CalculationHandler) domainErrorToStatus(code calculation.ErrorCode) int {
	switch code {
	case calculation.CodeInvalidInput:
		return http.StatusBadRequest
	case calculation.CodeUnknownOperation, calculation.CodeInvalidArity,
		calculation.CodeDivisionByZero, calculation.CodeNegativeSquareRoot,
		calculation.CodeNonFiniteNumber, calculation.CodeUnknownReference,
		calculation.CodeCycleDetected:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusUnprocessableEntity
	}
}

func (h *CalculationHandler) writeJSON(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *CalculationHandler) writeError(w http.ResponseWriter, status int, code calculation.ErrorCode, message string) {
	h.writeJSON(w, status, calculation.DomainError{
		Code:    code,
		Message: message,
	})
}

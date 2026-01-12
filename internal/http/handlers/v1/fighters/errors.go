package fighters

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	dtos "github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	repo "github.com/OJOMB/fightpicker/internal/repo/users"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

// writeError is a helper function to create a JSON formatted error from a user service or handler level error.
func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, err error, logger logs.Logger) {
	reqID, ok := ctx.Value(contextual.KeyRequestID).(string)
	if !ok {
		logger.ErrorContext(ctx, "request ID not found in context")
		reqID = "unknown"
	}

	var status int
	var resp dtos.ErrorEnvelope
	switch {
	case errors.Is(err, service.ErrMissingParameter):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeMissingRequiredParameter, reqID)
	case errors.Is(err, service.ErrInvalidParameter), errors.Is(err, v1.ErrInvalidUUID):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeInvalidParameter, reqID)
	case errors.Is(err, repo.ErrUserNotFound):
		status = http.StatusNotFound
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeResourceNotFound, reqID)
	default:
		status = http.StatusInternalServerError
		logger.ErrorContext(ctx, "internal server error", "error", err)
		resp = dtos.NewErrorEnvelope(v1.ErrInternalServerError, v1.ErrCodeInternalServerError, reqID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode error response", "error", err)
	}
}

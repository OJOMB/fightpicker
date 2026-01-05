package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	repo "github.com/OJOMB/fightpicker/internal/repo/users"
	"github.com/OJOMB/fightpicker/internal/server/dtos"
	service "github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	// ErrMissingRefreshToken indicates that the refresh token was not provided in the request cookie.
	ErrMissingRefreshToken = errors.New("missing refresh token")
)

// writeError is a helper function to create a JSON formatted error from a user service or handler level error.
func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, err error, logger logs.Logger) {
	reqID, ok := ctx.Value(contextual.KeyRequestID).(string)
	if !ok {
		logger.WarnContext(ctx, "request ID not found in context")
		reqID = "unknown"
	}

	var status int
	var resp dtos.ErrorEnvelope
	switch {
	case errors.Is(err, service.ErrMissingParameter), errors.Is(err, ErrMissingRefreshToken):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, "MISSING_REQUIRED_PARAMETER", reqID)
	case strings.Contains(err.Error(), "email: failed to pass regex validation"):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, "INVALID_PARAMETER", reqID)
	case errors.Is(err, repo.ErrUserNotFound):
		status = http.StatusNotFound
		resp = dtos.NewErrorEnvelope(err, "RESOURCE_NOT_FOUND", reqID)
	case errors.Is(err, service.ErrInvalidCredentials):
		status = http.StatusUnauthorized
		resp = dtos.NewErrorEnvelope(err, "INVALID_CREDENTIALS", reqID)
	default:
		status = http.StatusInternalServerError
		logger.ErrorContext(ctx, "internal server error", "error", err)
		resp = dtos.NewErrorEnvelope(fmt.Errorf("internal server error"), "INTERNAL_SERVER_ERROR", reqID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode error response", "error", err)
	}
}

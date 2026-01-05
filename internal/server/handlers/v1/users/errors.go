package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	repo "github.com/OJOMB/fightpicker/internal/repo/users"
	dtos "github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	service "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	// ErrIncompatibleQueryParameters is returned when incompatible query parameters are provided.
	ErrIncompatibleQueryParameters = errors.New("incompatible query parameters provided")
)

// writeError is a helper function to create a JSON formatted error from a user service or handler level error.
func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, err error, logger logs.Logger) {
	reqID, ok := ctx.Value(contextual.KeyRequestID).(string)
	if !ok {
		// This should never happen
		logger.ErrorContext(ctx, "request ID not found in context")
		reqID = "unknown"
	}

	var status int
	var resp dtos.ErrorEnvelope
	switch {
	case errors.Is(err, service.ErrMissingParameter):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeMissingRequiredParameter, reqID)
	case errors.Is(err, v1.ErrUnreadableRequestBody),
		errors.Is(err, v1.ErrInvalidJSONRequestBody):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeMalformedRequestBody, reqID)
	case strings.Contains(err.Error(), "email: failed to pass regex validation"),
		errors.Is(err, ErrIncompatibleQueryParameters),
		errors.Is(err, service.ErrInvalidParameter),
		errors.Is(err, v1.ErrInvalidUUID),
		errors.Is(err, v1.ErrUnreadableRequestBody),
		errors.Is(err, v1.ErrInvalidJSONRequestBody),
		errors.Is(err, v1.ErrMissingRequiredQueryParameter):
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeInvalidParameter, reqID)
	case errors.Is(err, service.ErrUnauthorized):
		status = http.StatusUnauthorized
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeUnauthorized, reqID)
	case errors.Is(err, repo.ErrUserNotFound):
		status = http.StatusNotFound
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeResourceNotFound, reqID)
	case errors.Is(err, repo.ErrEmailTaken), errors.Is(err, repo.ErrUsernameTaken):
		status = http.StatusConflict
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeConflictingResources, reqID)
	case errors.Is(err, repo.ErrDefaultRoleNotFound):
		status = http.StatusInternalServerError
		logger.ErrorContext(ctx, "roles table not setup properly - default user role missing", "error", err)
		resp = dtos.NewErrorEnvelope(v1.ErrInternalServerError, v1.ErrCodeInternalServerError, reqID)
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

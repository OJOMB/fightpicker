package users

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	usersrepo "github.com/OJOMB/fightpicker/internal/repo/users"
	dtos "github.com/OJOMB/fightpicker/internal/server/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/server/handlers/v1"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	// ErrIncompatibleQueryParameters is returned when incompatible query parameters are provided.
	ErrIncompatibleQueryParameters = errors.New("incompatible query parameters provided")
)

// writeError is a helper function to create a JSON formatted error from a user service/repo or handler level error.
// It maps specific errors to appropriate HTTP status codes and logs the errors accordingly.
// log level is determined based on the severity of the error. Error logs are reserved for genuine server-side issues,
// while client-side errors are logged at the debug level.
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
	case errors.Is(err, usersservice.ErrMissingParameter):
		// 400 Bad Request missing parameters
		logger.DebugContext(ctx, "missing required parameter", "error", err)
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeMissingRequiredParameter, reqID)
	case errors.Is(err, v1.ErrUnreadableRequestBody),
		errors.Is(err, v1.ErrInvalidJSONRequestBody):
		// 400 Bad Request malformed request body
		logger.DebugContext(ctx, "malformed request body", "error", err)
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeMalformedRequestBody, reqID)
	case strings.Contains(err.Error(), "email: failed to pass regex validation"),
		errors.Is(err, ErrIncompatibleQueryParameters),
		errors.Is(err, usersservice.ErrInvalidParameter),
		errors.Is(err, v1.ErrInvalidUUID),
		errors.Is(err, v1.ErrMissingRequiredQueryParameter):
		// 400 Bad Request invalid parameters
		logger.DebugContext(ctx, "invalid parameter(s)", "error", err)
		status = http.StatusBadRequest
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeInvalidParameter, reqID)
	case errors.Is(err, usersservice.ErrUnauthorized):
		// 401 Unauthorized request lacks valid authentication credentials
		logger.DebugContext(ctx, "unauthorized access attempt", "error", err)
		status = http.StatusUnauthorized
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeUnauthorized, reqID)
	case errors.Is(err, usersrepo.ErrUserNotFound):
		// 404 Not Found resource does not exist
		logger.DebugContext(ctx, "requested resource not found", "error", err)
		status = http.StatusNotFound
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeResourceNotFound, reqID)
	case errors.Is(err, usersrepo.ErrEmailTaken),
		errors.Is(err, usersrepo.ErrUsernameTaken):
		// 409 Resource conflict
		logger.DebugContext(ctx, "resource conflict", "error", err)
		status = http.StatusConflict
		resp = dtos.NewErrorEnvelope(err, v1.ErrCodeConflictingResources, reqID)
	case errors.Is(err, usersrepo.ErrDefaultRoleNotFound):
		// 500 Internal Server Error - system not setup properly
		logger.ErrorContext(ctx, "system setup error", "error", err)
		status = http.StatusInternalServerError
		resp = dtos.NewErrorEnvelope(v1.ErrInternalServerError, v1.ErrCodeInternalServerError, reqID)
	default:
		// 500 Internal Server Error - generic catch-all for unexpected errors
		logger.ErrorContext(ctx, "internal server error", "error", err)
		status = http.StatusInternalServerError
		resp = dtos.NewErrorEnvelope(v1.ErrInternalServerError, v1.ErrCodeInternalServerError, reqID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.ErrorContext(ctx, "failed to encode error response", "error", err)
	}
}

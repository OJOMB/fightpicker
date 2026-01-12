package users

import (
	"context"
	"errors"
	"net/http"

	dtos "github.com/OJOMB/fightpicker/internal/http/dtos"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	usersrepo "github.com/OJOMB/fightpicker/internal/repo/users"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/contextual"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

type apiError struct {
	Status   int
	Code     string
	LogLevel logs.Level
	LogMsg   string
	Public   error
}

func classifyError(err error) apiError {
	switch {
	case errors.Is(err, usersservice.ErrMissingParameter):
		return apiError{
			Status:   http.StatusBadRequest,
			Code:     v1.ErrCodeMissingRequiredParameter,
			LogLevel: logs.LevelDebug,
			LogMsg:   "missing required parameter",
			Public:   err,
		}

	case errors.Is(err, v1.ErrUnreadableRequestBody),
		errors.Is(err, v1.ErrInvalidJSONRequestBody):
		return apiError{
			Status:   http.StatusBadRequest,
			Code:     v1.ErrCodeMalformedRequestBody,
			LogLevel: logs.LevelDebug,
			LogMsg:   "malformed request body",
			Public:   err,
		}

	case errors.Is(err, v1.ErrIncompatibleParameters),
		errors.Is(err, usersservice.ErrInvalidParameter),
		errors.Is(err, v1.ErrInvalidUUID),
		errors.Is(err, v1.ErrMissingRequiredQueryParameter):
		return apiError{
			Status:   http.StatusBadRequest,
			Code:     v1.ErrCodeInvalidParameter,
			LogLevel: logs.LevelDebug,
			LogMsg:   "invalid parameter(s)",
			Public:   err,
		}

	case errors.Is(err, usersservice.ErrUnauthorized):
		return apiError{
			Status:   http.StatusUnauthorized,
			Code:     v1.ErrCodeUnauthorized,
			LogLevel: logs.LevelDebug,
			LogMsg:   "unauthorized access attempt",
			Public:   err,
		}

	case errors.Is(err, usersrepo.ErrUserNotFound):
		return apiError{
			Status:   http.StatusNotFound,
			Code:     v1.ErrCodeResourceNotFound,
			LogLevel: logs.LevelDebug,
			LogMsg:   "resource not found",
			Public:   err,
		}

	case errors.Is(err, usersrepo.ErrEmailTaken),
		errors.Is(err, usersrepo.ErrUsernameTaken):
		return apiError{
			Status:   http.StatusConflict,
			Code:     v1.ErrCodeConflictingResources,
			LogLevel: logs.LevelDebug,
			LogMsg:   "resource conflict",
			Public:   err,
		}

	case errors.Is(err, usersrepo.ErrDefaultRoleNotFound):
		return apiError{
			Status:   http.StatusInternalServerError,
			Code:     v1.ErrCodeInternalServerError,
			LogLevel: logs.LevelError,
			LogMsg:   "system setup error",
			Public:   v1.ErrInternalServerError,
		}

	default:
		return apiError{
			Status:   http.StatusInternalServerError,
			Code:     v1.ErrCodeInternalServerError,
			LogLevel: logs.LevelError,
			LogMsg:   "internal server error",
			Public:   v1.ErrInternalServerError,
		}
	}
}

// writeError is a helper function to create a JSON formatted error from a user service/repo or handler level error.
// It maps specific errors to appropriate HTTP status codes and logs the errors accordingly.
// log level is determined based on the severity of the error. Error logs are reserved for genuine server-side issues,
// while client-side errors are logged at the debug level.
func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, err error, logger logs.Logger) {
	reqID, _ := ctx.Value(contextual.KeyRequestID).(string)
	if reqID == "" {
		reqID = "unknown"
	}

	apiErr := classifyError(err)

	switch apiErr.LogLevel {
	case logs.LevelDebug:
		logger.DebugContext(ctx, apiErr.LogMsg, "error", err)
	case logs.LevelError:
		logger.ErrorContext(ctx, apiErr.LogMsg, "error", err)
	}

	resp := dtos.NewErrorEnvelope(apiErr.Public, apiErr.Code, reqID)
	h.writeJSON(ctx, w, logger, apiErr.Status, resp)
}

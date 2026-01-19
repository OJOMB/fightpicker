package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	usersrepo "github.com/OJOMB/fightpicker/internal/repo/users"
	usersservice "github.com/OJOMB/fightpicker/internal/service/users"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	ErrContextToolIsNil = fmt.Errorf("contextTool cannot be nil")
	ErrLoggerIsNil      = fmt.Errorf("logger cannot be nil")
	ErrIDToolIsNil      = fmt.Errorf("idTool cannot be nil")
)

func classifyError(err error) *apierr.APIError {
	var (
		status    int
		code      string
		logLevel  logs.Level
		logMsg    string
		publicErr error
	)

	switch {
	case errors.Is(err, usersservice.ErrMissingParameter):
		status = http.StatusBadRequest
		code = v1.ErrCodeMissingRequiredParameter
		logLevel = logs.LevelDebug
		logMsg = "missing required parameter"
		publicErr = err

	case errors.Is(err, v1.ErrUnreadableRequestBody),
		errors.Is(err, v1.ErrInvalidJSONRequestBody):
		status = http.StatusBadRequest
		code = v1.ErrCodeMalformedRequestBody
		logLevel = logs.LevelDebug
		logMsg = "malformed request body"
		publicErr = err

	case errors.Is(err, v1.ErrIncompatibleParameters),
		errors.Is(err, usersservice.ErrInvalidParameter),
		errors.Is(err, v1.ErrInvalidUUID),
		errors.Is(err, v1.ErrMissingRequiredQueryParameter):
		status = http.StatusBadRequest
		code = v1.ErrCodeInvalidParameter
		logLevel = logs.LevelDebug
		logMsg = "invalid parameter(s)"
		publicErr = err

	case errors.Is(err, usersservice.ErrUnauthorized):
		status = http.StatusUnauthorized
		code = v1.ErrCodeUnauthorized
		logLevel = logs.LevelDebug
		logMsg = "unauthorized access attempt"
		publicErr = err

	case errors.Is(err, usersrepo.ErrUserNotFound):
		status = http.StatusNotFound
		code = v1.ErrCodeResourceNotFound
		logLevel = logs.LevelDebug
		logMsg = "resource not found"
		publicErr = err

	case errors.Is(err, usersrepo.ErrEmailTaken),
		errors.Is(err, usersrepo.ErrUsernameTaken):
		status = http.StatusConflict
		code = v1.ErrCodeConflictingResources
		logLevel = logs.LevelDebug
		logMsg = "resource conflict"
		publicErr = err

	case errors.Is(err, usersrepo.ErrDefaultRoleNotFound):
		status = http.StatusInternalServerError
		code = v1.ErrCodeInternalServerError
		logLevel = logs.LevelError
		logMsg = "system setup error"
		publicErr = v1.ErrInternalServerError

	default:
		return apierr.NewInternalError(err)
	}

	return apierr.NewAPIError(
		status,
		code,
		logLevel,
		logMsg,
		publicErr,
	)
}

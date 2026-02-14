package auth

import (
	"errors"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	repo "github.com/OJOMB/fightpicker/internal/repo/users"
	authservice "github.com/OJOMB/fightpicker/internal/service/auth"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	// ErrMissingRefreshToken indicates that the refresh token is missing.
	ErrMissingRefreshToken = errors.New("missing refresh token")
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
	case errors.Is(err, authservice.ErrMissingParameter), errors.Is(err, ErrMissingRefreshToken):
		status = http.StatusBadRequest
		code = v1.ErrCodeMissingRequiredParameter
		logLevel = logs.LevelDebug
		logMsg = "missing required parameter"
		publicErr = err
	case errors.Is(err, openapi_types.ErrValidationEmail):
		status = http.StatusBadRequest
		code = v1.ErrCodeInvalidParameter
		logLevel = logs.LevelDebug
		logMsg = "invalid parameter"
		publicErr = v1.ErrInvalidEmailFormat
	case errors.Is(err, repo.ErrUserNotFound):
		status = http.StatusNotFound
		code = v1.ErrCodeResourceNotFound
		logLevel = logs.LevelDebug
		logMsg = "resource not found"
		publicErr = repo.ErrUserNotFound
	case errors.Is(err, authservice.ErrInvalidCredentials):
		status = http.StatusUnauthorized
		code = v1.ErrCodeInvalidCredentials
		logLevel = logs.LevelDebug
		logMsg = "invalid credentials provided"
		publicErr = authservice.ErrInvalidCredentials
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

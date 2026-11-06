package cards

import (
	"errors"
	"net/http"
	"time"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	cardsrepo "github.com/OJOMB/fightpicker/internal/repo/cards"
	cardsservice "github.com/OJOMB/fightpicker/internal/service/cards"
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

	var parseErr *time.ParseError

	switch {
	case errors.As(err, &parseErr):
		status = http.StatusBadRequest
		code = v1.ErrCodeInvalidParameter
		logLevel = logs.LevelDebug
		logMsg = "invalid time format"
		publicErr = err
	case errors.Is(err, cardsservice.ErrMissingParameter), errors.Is(err, ErrMissingRefreshToken):
		status = http.StatusBadRequest
		code = v1.ErrCodeMissingRequiredParameter
		logLevel = logs.LevelDebug
		logMsg = "missing required parameter"
		publicErr = err
	case errors.Is(err, cardsservice.ErrMissingParameter), errors.Is(err, ErrMissingRefreshToken):
		status = http.StatusBadRequest
		code = v1.ErrCodeMissingRequiredParameter
		logLevel = logs.LevelDebug
		logMsg = "missing required parameter"
		publicErr = err
	case errors.Is(err, cardsrepo.ErrResourceNotFound):
		status = http.StatusNotFound
		code = v1.ErrCodeResourceNotFound
		logLevel = logs.LevelDebug
		logMsg = "resource not found"
		publicErr = cardsrepo.ErrResourceNotFound
	case errors.Is(err, cardsrepo.ErrDuplicateResource):
		status = http.StatusConflict
		code = v1.ErrCodeConflictingResources
		logLevel = logs.LevelDebug
		logMsg = "duplicate resource detected"
		publicErr = cardsrepo.ErrDuplicateResource
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

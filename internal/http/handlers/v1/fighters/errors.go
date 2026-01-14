package fighters

import (
	"errors"
	"net/http"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
	v1 "github.com/OJOMB/fightpicker/internal/http/handlers/v1"
	fightersrepo "github.com/OJOMB/fightpicker/internal/repo/fighters"
	fightersservice "github.com/OJOMB/fightpicker/internal/service/fighters"
	"github.com/OJOMB/fightpicker/pkg/logs"
)

func classifyError(err error) apierr.APIError {
	var (
		status    int
		code      string
		logLevel  logs.Level
		logMsg    string
		publicErr error
	)

	switch {
	case errors.Is(err, fightersservice.ErrMissingParameter):
		status = http.StatusBadRequest
		code = v1.ErrCodeMissingRequiredParameter
		logLevel = logs.LevelDebug
		logMsg = "missing required parameter"
		publicErr = err
	case errors.Is(err, v1.ErrInvalidUUID), errors.Is(err, fightersservice.ErrInvalidParameter):
		status = http.StatusBadRequest
		code = v1.ErrCodeInvalidParameter
		logLevel = logs.LevelDebug
		logMsg = "invalid parameter(s)"
		publicErr = err
	case errors.Is(err, fightersrepo.ErrFighterNotFound):
		status = http.StatusNotFound
		code = v1.ErrCodeResourceNotFound
		logLevel = logs.LevelDebug
		logMsg = "resource not found"
		publicErr = fightersrepo.ErrFighterNotFound
	default:
		status = http.StatusInternalServerError
		code = v1.ErrCodeInternalServerError
		logLevel = logs.LevelError
		logMsg = "internal server error"
		publicErr = v1.ErrInternalServerError
	}

	return apierr.NewAPIError(
		status,
		code,
		logLevel,
		logMsg,
		publicErr,
	)
}

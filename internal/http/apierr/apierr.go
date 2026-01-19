package apierr

import (
	"fmt"
	"net/http"

	"github.com/OJOMB/fightpicker/pkg/logs"
)

var (
	// ErrInternalServerError represents a generic internal server error.
	ErrInternalServerError = fmt.Errorf("internal server error")
	// Err
)

const (
	// ErrCodeInternalServerError
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	// ErrCodeBadRequest represents a bad request error code.
	ErrCodeBadRequest = "BAD_REQUEST"
	// ErrCodeUnauthorized represents an unauthorized error code.
	ErrCodeUnauthorized = "UNAUTHORIZED"
	// ErrCodeForbidden represents a forbidden error code.
	ErrCodeForbidden = "FORBIDDEN"
	// ErrCodeNotFound represents a not found error code.
	ErrCodeNotFound = "NOT_FOUND"
	// ErrCodeConflict represents a conflict error code.
	ErrCodeConflict = "CONFLICT"
)

type APIErrClassifier func(err error) *APIError

type APIError struct {
	Status   int
	Code     string
	LogLevel logs.Level
	LogMsg   string
	Public   error
}

func NewAPIError(
	status int,
	code string,
	logLevel logs.Level,
	logMsg string,
	public error,
) *APIError {
	return &APIError{
		Status:   status,
		Code:     code,
		LogLevel: logLevel,
		LogMsg:   logMsg,
		Public:   public,
	}
}

func NewInternalError(err error) *APIError {
	return NewAPIError(
		http.StatusInternalServerError,
		ErrCodeInternalServerError,
		logs.LevelError,
		err.Error(),
		ErrInternalServerError,
	)
}

func NewBadRequestError(err error) *APIError {
	return NewAPIError(
		http.StatusBadRequest,
		ErrCodeBadRequest,
		logs.LevelDebug,
		err.Error(),
		err,
	)
}

func NewUnauthorizedError(err error) *APIError {
	return NewAPIError(
		http.StatusUnauthorized,
		ErrCodeUnauthorized,
		logs.LevelDebug,
		err.Error(),
		err,
	)
}

func NewForbiddenError(err error) *APIError {
	return NewAPIError(
		http.StatusForbidden,
		ErrCodeForbidden,
		logs.LevelDebug,
		err.Error(),
		err,
	)
}

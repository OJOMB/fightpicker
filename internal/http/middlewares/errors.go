package middlewares

import (
	"errors"
	"fmt"

	"github.com/OJOMB/fightpicker/internal/http/apierr"
)

var (
	ErrSecretKeyIsNilOrEmpty = fmt.Errorf("secret key is nil or empty")
	ErrJWTValidatorIsNil     = fmt.Errorf("jwt validator is nil")
	ErrContextToolIsNil      = fmt.Errorf("context tool is nil")
	ErrLoggerIsNil           = fmt.Errorf("logger is nil")

	ErrInvalidRouteName        = fmt.Errorf("route name invalid")
	ErrMissingToken            = fmt.Errorf("token missing from request")
	ErrInvalidToken            = fmt.Errorf("invalid token")
	ErrExpiredToken            = fmt.Errorf("expired token")
	ErrInsufficientPermissions = fmt.Errorf("insufficient permissions")
)

func classifyError(err error) *apierr.APIError {
	switch {
	case
		errors.Is(err, ErrInvalidToken):
		return apierr.NewBadRequestError(err)
	case errors.Is(err, ErrExpiredToken),
		errors.Is(err, ErrMissingToken):
		return apierr.NewUnauthorizedError(err)
	case errors.Is(err, ErrInsufficientPermissions):
		return apierr.NewForbiddenError(err)
	default:
		return apierr.NewInternalError(err)
	}
}

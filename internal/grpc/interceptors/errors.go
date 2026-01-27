package interceptors

import "fmt"

var (
	ErrSecretKeyIsNilOrEmpty = fmt.Errorf("secret key is nil or empty")
	ErrJWTValidatorIsNil     = fmt.Errorf("JWT validator is nil")
	ErrContextToolIsNil      = fmt.Errorf("context tool is nil")
)

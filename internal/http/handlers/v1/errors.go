package v1

import "fmt"

var (
	// ErrLoggerIsNil is returned when the provided logger is nil.
	ErrLoggerIsNil = fmt.Errorf("logger cannot be nil")
	// ErrContextToolIsNil is returned when the provided context tool is nil.
	ErrContextToolIsNil = fmt.Errorf("context tool cannot be nil")
	// ErrIDToolIsNil is returned when the provided ID tool is nil.
	ErrIDToolIsNil = fmt.Errorf("ID tool cannot be nil")
	// ErrJWTToolIsNil is returned when the provided JWT tool is nil.
	ErrJWTToolIsNil = fmt.Errorf("JWT tool cannot be nil")

	// ErrInvalidUUID is returned when a provided UUID is invalid.
	ErrInvalidUUID = fmt.Errorf("invalid UUID format")
	// ErrEmptyRequestBody is returned when the request body is empty.
	ErrEmptyRequestBody = fmt.Errorf("empty request body")
	// ErrUnreadableRequestBody is returned when the request body cannot be read.
	ErrUnreadableRequestBody = fmt.Errorf("unreadable request body")
	// ErrInvalidJSONRequestBody is returned when the request body contains invalid JSON.
	ErrInvalidJSONRequestBody = fmt.Errorf("invalid JSON in request body")
	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = fmt.Errorf("internal server error")
	// ErrMissingRequiredQueryParameter is returned when a required query parameter is missing.
	ErrMissingRequiredQueryParameter = fmt.Errorf("missing required query parameter")
	// ErrInvalidPageSize is returned when the provided page size is invalid.
	ErrInvalidPageSize = fmt.Errorf("invalid page size")
	// ErrIncompatibleParameters is returned when incompatible parameters are provided.
	ErrIncompatibleParameters = fmt.Errorf("incompatible parameters provided")
)

var (
	ErrCodeMissingRequiredParameter = "MISSING_REQUIRED_PARAMETER"
	ErrCodeInvalidParameter         = "INVALID_PARAMETER"
	ErrCodeResourceNotFound         = "RESOURCE_NOT_FOUND"
	ErrCodeInvalidCredentials       = "INVALID_CREDENTIALS"
	ErrCodeConflictingResources     = "CONFLICTING_RESOURCES"
	ErrCodeInternalServerError      = "INTERNAL_SERVER_ERROR"
	ErrCodeMalformedRequestBody     = "MALFORMED_REQUEST_BODY"
	ErrCodeUnauthorized             = "UNAUTHORIZED"
)

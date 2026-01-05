package v1

import "errors"

var (
	// ErrInvalidUUID is returned when a provided UUID is invalid.
	ErrInvalidUUID = errors.New("invalid UUID format")
	// ErrUnreadableRequestBody is returned when the request body cannot be read.
	ErrUnreadableRequestBody = errors.New("unreadable request body")
	// ErrInvalidJSONRequestBody is returned when the request body contains invalid JSON.
	ErrInvalidJSONRequestBody = errors.New("invalid JSON in request body")
	// ErrInternalServerError is returned when an internal server error occurs.
	ErrInternalServerError = errors.New("internal server error")
	// ErrMissingRequiredQueryParameter is returned when a required query parameter is missing.
	ErrMissingRequiredQueryParameter = errors.New("missing required query parameter")
)

var (
	ErrCodeMissingRequiredParameter = "MISSING_REQUIRED_PARAMETER"
	ErrCodeInvalidParameter         = "INVALID_PARAMETER"
	ErrCodeResourceNotFound         = "RESOURCE_NOT_FOUND"
	ErrCodeConflictingResources     = "CONFLICTING_RESOURCES"
	ErrCodeInternalServerError      = "INTERNAL_SERVER_ERROR"
	ErrCodeMalformedRequestBody     = "MALFORMED_REQUEST_BODY"
	ErrCodeUnauthorized             = "UNAUTHORIZED"
)

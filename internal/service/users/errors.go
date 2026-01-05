package users

import "fmt"

var (
	// ErrNilRepo is returned when a nil repository is provided during service initialization.
	ErrNilRepo = fmt.Errorf("repository cannot be nil")

	// ErrNilLogger is returned when a nil logger is provided.
	ErrNilLogger = fmt.Errorf("logger cannot be nil")

	// ErrEmailRegexCompile is returned when the email regex fails to compile.
	ErrEmailRegexCompile = fmt.Errorf("failed to compile email regex")

	// ErrNilIDGenerator is returned when a nil ID generator function is provided.
	ErrNilIDGenerator = fmt.Errorf("ID generator function cannot be nil")

	// ErrNilNowFunc is returned when a nil time function is provided.
	ErrNilNowFunc = fmt.Errorf("now function cannot be nil")

	// ErrNilPasswordHasher is returned when a nil password hasher is provided.
	ErrNilPasswordHasher = fmt.Errorf("password hasher cannot be nil")

	// ErrNilImageProcessor is returned when a nil image processor is provided.
	ErrNilImageProcessor = fmt.Errorf("image processor cannot be nil")

	// ErrMissingParameter is returned when required data is ommitted from the request.
	ErrMissingParameter = fmt.Errorf("missing parameter")

	// ErrInvalidParameter is returned when provided data is invalid.
	ErrInvalidParameter = fmt.Errorf("invalid parameter")

	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = fmt.Errorf("internal error")

	// ErrEmptyUpdate is returned when an update request contains no fields to update.
	ErrEmptyUpdate = fmt.Errorf("no fields to update")

	// ErrUnauthorized is returned when the requesting user is not authorized to perform the requested action.
	ErrUnauthorized = fmt.Errorf("user is not authorized to perform this action")
)

package users

import "fmt"

// private initialization errors
var (
	// errNilRepo is returned when a nil repository is provided during service initialization.
	errNilRepo = fmt.Errorf("repository cannot be nil")

	// errNilLogger is returned when a nil logger is provided.
	errNilLogger = fmt.Errorf("logger cannot be nil")

	// errNilContextTool is returned when a nil context tool is provided.
	errNilContextTool = fmt.Errorf("context tool cannot be nil")

	// errEmailRegexCompile is returned when the email regex fails to compile.
	errEmailRegexCompile = fmt.Errorf("failed to compile email regex")

	// errNilIDTool is returned when a nil ID generator function is provided.
	errNilIDTool = fmt.Errorf("ID tool cannot be nil")

	// errNilNowFunc is returned when a nil time function is provided.
	errNilNowFunc = fmt.Errorf("now function cannot be nil")

	// errNilPasswordHasher is returned when a nil password hasher is provided.
	errNilPasswordHasher = fmt.Errorf("password hasher cannot be nil")

	// errNilImageProcessor is returned when a nil image processor is provided.
	errNilImageProcessor = fmt.Errorf("image processor cannot be nil")
)

// public runtime errors
var (
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

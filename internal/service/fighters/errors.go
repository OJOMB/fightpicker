package fighters

import "fmt"

// private initialization errors
var (
	// errNilRepo is returned when a nil repository is provided to the service during initialization.
	errNilRepo = fmt.Errorf("repository cannot be nil")

	// errNilLogger is returned when a nil logger is provided to the service during initialization.
	errNilLogger = fmt.Errorf("logger cannot be nil")

	// errNilIDGenerator is returned when a nil ID generator is provided to the service during initialization.
	errNilIDGenerator = fmt.Errorf("ID generator cannot be nil")

	// errNilDatetimeTool is returned when a nil datetime tool is provided to the service during initialization.
	errNilDatetimeTool = fmt.Errorf("now function cannot be nil")
)

// public runtime errors
var (
	// ErrFighterNotFound is returned when a fighter with the specified ID does not exist.
	ErrFighterNotFound = fmt.Errorf("fighter not found")

	// ErrMissingParameter is returned when a required parameter is missing.
	ErrMissingParameter = fmt.Errorf("missing required parameter")

	// ErrInvalidParameter is returned when a provided parameter is invalid.
	ErrInvalidParameter = fmt.Errorf("invalid parameter provided")

	// ErrInternalError is returned when an internal error occurs.
	ErrInternalError = fmt.Errorf("internal error occurred")

	// ErrInvalidFighterID is returned when an invalid fighter ID is provided.
	ErrInvalidFighterID = fmt.Errorf("invalid fighter ID provided")

	// ErrUnauthorized is returned when the requester is not authorized to perform the action.
	ErrUnauthorized = fmt.Errorf("unauthorized to perform this action")
)

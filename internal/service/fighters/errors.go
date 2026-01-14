package fighters

import "fmt"

var (
	// ErrNilRepo is returned when a nil repository is provided to the service during initialization.
	ErrNilRepo = fmt.Errorf("repository cannot be nil")

	// ErrNilLogger is returned when a nil logger is provided to the service during initialization.
	ErrNilLogger = fmt.Errorf("logger cannot be nil")

	// ErrNilIDGenerator is returned when a nil ID generator is provided to the service during initialization.
	ErrNilIDGenerator = fmt.Errorf("ID generator cannot be nil")

	// ErrNilNowFunc is returned when a nil now function is provided to the service during initialization.
	ErrNilNowFunc = fmt.Errorf("now function cannot be nil")

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
)

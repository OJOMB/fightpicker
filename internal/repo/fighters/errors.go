package fighters

import "fmt"

var (
	// ErrNilLogger is returned when a nil logger is provided to the repository during initialization.
	ErrNilLogger = fmt.Errorf("logger cannot be nil")
	// ErrNilDBPool is returned when a nil database pool is provided to the repository during initialization.
	ErrNilDBPool = fmt.Errorf("database pool cannot be nil")
	// ErrNilDBClient is returned when a nil database client is provided to the repository during initialization.
	ErrNilDBClient = fmt.Errorf("database client cannot be nil")
	// ErrNilNowFunc is returned when a nil now function is provided to the repository during initialization.
	ErrNilNowFunc = fmt.Errorf("now function cannot be nil")
)

func dbErrorToServiceError(err error) error {
	// Here you can map specific database errors to service-level errors.
	// For simplicity, we'll just return a generic error for now.
	return fmt.Errorf("database error: %w", err)
}

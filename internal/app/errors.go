package app

import "fmt"

var (
	// ErrNilConfig is returned when a nil configuration is provided.
	ErrNilConfig = fmt.Errorf("nil configuration")
	// ErrInvalidConfig is returned when the configuration passed is invalid.
	ErrInvalidConfig = fmt.Errorf("invalid configuration")
	// ErrAppNotInitialized is returned when the app has not been properly initialized.
	ErrAppNotInitialized = fmt.Errorf("app not initialized")
)

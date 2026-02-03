package app

import "fmt"

var (
	// errNilConfig is returned when a nil configuration is provided.
	errNilConfig = fmt.Errorf("nil configuration")
	// errInvalidConfig is returned when the configuration passed is invalid.
	errInvalidConfig = fmt.Errorf("invalid configuration")
	// errAppNotInitialized is returned when the app has not been properly initialized.
	errAppNotInitialized = fmt.Errorf("app not initialized")
)

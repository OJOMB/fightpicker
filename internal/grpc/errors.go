package grpc

import "fmt"

var (
	////////////////////////////
	// initialisation errors //
	//////////////////////////

	// errNilConfig is returned when a nil config is passed to New function.
	errNilConfig = fmt.Errorf("nil config provided")

	// errNilApp is returned when a nil app is passed to New function.
	errNilApp = fmt.Errorf("nil app provided")

	// errFailedToInitInterceptor is returned when an interceptor fails to be created.
	errFailedToInitInterceptor = fmt.Errorf("failed to initialize interceptor")

	// errListeningFailed is returned when the server fails to start listening on the specified address.
	errListeningFailed = fmt.Errorf("failed to start listening on the specified address")
)

package server

import "fmt"

var (
	// ErrNoDomainSpecified is returned when no domain is specified for the server to bind to.
	ErrNoDomainSpecified = fmt.Errorf("init failure - server domain not specified")
	// ErrPortNotSpecified is returned when no port is specified for the server to listen on.
	ErrPortNotSpecified = fmt.Errorf("init failure - server port not specified")
)

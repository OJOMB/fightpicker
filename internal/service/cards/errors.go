package cards

import (
	"fmt"
)

// private initialization errors
var (
	// errNilLogger is returned when a nil logger is provided to the repository during initialization.
	errNilLogger = fmt.Errorf("logger cannot be nil")
	// errNilRepo is returned when a nil database client is provided to the repository during initialization.
	errNilRepo = fmt.Errorf("repo cannot be nil")
	// errNilDatetimeTool is returned when a nil datetime tool is provided.
	errNilDatetimeTool = fmt.Errorf("datetime tool cannot be nil")
	// errNilContextTool is returned when a nil context tool is provided.
	errNilContextTool = fmt.Errorf("context tool cannot be nil")
	// errNilIDTool is returned when a nil ID tool is provided.
	errNilIDTool = fmt.Errorf("ID tool cannot be nil")
)

var (
	ErrMissingParameter = fmt.Errorf("missing parameter")
)

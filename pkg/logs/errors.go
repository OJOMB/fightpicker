package logs

import "fmt"

// private initialization errors
var (
	// errNilContextTool is returned when a nil context tool is provided.
	errNilContextTool = fmt.Errorf("context tool cannot be nil")
	// errEmptyEnv is returned when an empty environment string is provided.
	errEmptyEnv = fmt.Errorf("environment cannot be empty")
	// errEmptyAppName is returned when an empty application name string is provided.
	errEmptyAppName = fmt.Errorf("application name cannot be empty")
)

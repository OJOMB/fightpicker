package apiresponder

import "fmt"

var (
	ErrNilLogger          = fmt.Errorf("nil logger provided to JSONErrorWriter")
	ErrNilContextTool     = fmt.Errorf("nil context tool provided to JSONErrorWriter")
	ErrNilErrorClassifier = fmt.Errorf("nil error classifier provided to JSONErrorWriter")
)

package id

import "fmt"

var (
	ErrInvalidFormat = fmt.Errorf("invalid uuid format")

	ErrInvalidVersion7 = fmt.Errorf("invalid uuid version - expected v7")

	ErrInvalidVariant = fmt.Errorf("invalid uuid variant - expected RFC4122")
)

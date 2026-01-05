package media

import "fmt"

var (
	// ErrUnknownImageFormat is returned when the image format cannot be determined
	ErrUnknownImageFormat = fmt.Errorf("unknown image format")
	// ErrUnsupportedImageFormat is returned when the image format is not supported
	ErrUnsupportedImageFormat = fmt.Errorf("unsupported image format")
	// ErrImageTooLarge is returned when the image size exceeds the maximum allowed size
	ErrImageTooLarge = fmt.Errorf("image size exceeds maximum allowed size")
	// ErrImageDimensionsTooLarge is returned when the image dimensions exceed the maximum allowed dimensions
	ErrImageDimensionsTooLarge = fmt.Errorf("image dimensions exceed maximum allowed dimensions")
)

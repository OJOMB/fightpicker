package media

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/HugoSmits86/nativewebp"
	"golang.org/x/image/draw"
)

const (
	MaxPixels = 20_000_000 // 20MP hard cap
)

// ImageProcessor handles image processing tasks such as normalization and thumbnail generation.
// It currently supports JPEG, PNG, and WEBP formats. Uses native go image libraries rather than cgo bindings.
type ImageProcessor struct {
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) ProcessUserProfilePicture(input []byte) (normalized []byte, thumbnail []byte, err error) {
	img, format, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, nil, err
	}

	switch format {
	case "jpeg", "jpg", "png", "webp":
	default:
		return nil, nil, ErrUnsupportedImageFormat
	}

	b := img.Bounds()
	if b.Dx()*b.Dy() > MaxPixels {
		return nil, nil, ErrImageTooLarge
	}

	// Normalize + crop square
	square := centerCrop(img)

	// Resize main image
	main := resize(square, 800, 800)

	// Resize thumbnail
	thumb := resize(square, 100, 100)

	normalized, err = encodeWebP(main)
	if err != nil {
		return nil, nil, err
	}

	thumbnail, err = encodeWebP(thumb)
	if err != nil {
		return nil, nil, err
	}

	return normalized, thumbnail, nil
}

func centerCrop(img image.Image) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	size := w
	if min(w, h) == size {
		return img
	}

	x0 := b.Min.X + (w-size)/2
	y0 := b.Min.Y + (h-size)/2

	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(dst, dst.Bounds(), img, image.Point{X: x0, Y: y0}, draw.Src)

	return dst
}

func resize(img image.Image, w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func encodeWebP(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	options := &nativewebp.Options{
		UseExtendedFormat: true,
	}

	err := nativewebp.Encode(&buf, img, options)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

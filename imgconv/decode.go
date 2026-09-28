package imgconv

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

// MaxDecodePixels bounds the images that are decoded, in pixels (width ×
// height). The header of an image states its size and the decoders take
// the memory for it before they read the pixels, so a file of a hundred
// bytes could otherwise ask for gigabytes.
const MaxDecodePixels = 100 << 20

// ErrTooLarge is returned for an image of more than MaxDecodePixels pixels.
var ErrTooLarge = fmt.Errorf("imgconv: the image has more than %d pixels and is not decoded", MaxDecodePixels)

// decoder is the decoder of a format, and the reader of its header.
type decoder struct {
	decode func(io.Reader) (image.Image, error)
	config func(io.Reader) (image.Config, error)
}

var decoders = map[string]decoder{
	"png":  {png.Decode, png.DecodeConfig},
	"jpeg": {jpeg.Decode, jpeg.DecodeConfig},
	"gif":  {gif.Decode, gif.DecodeConfig},
	"bmp":  {bmp.Decode, bmp.DecodeConfig},
	"webp": {webp.Decode, webp.DecodeConfig},
}

// Decode decodes PNG, JPEG, GIF, BMP or WebP bytes (AVIF is not decodable
// here). An image of more than MaxDecodePixels pixels is not decoded: the
// error is ErrTooLarge.
func Decode(data []byte) (img image.Image, err error) {
	format := Sniff(data)
	d, ok := decoders[format]
	if !ok {
		return nil, fmt.Errorf("imgconv: cannot decode %q", format)
	}
	// A header that cannot be read is left to the decoder to report.
	if c, err := d.config(bytes.NewReader(data)); err == nil {
		if c.Width <= 0 || c.Height <= 0 || int64(c.Width)*int64(c.Height) > MaxDecodePixels {
			return nil, ErrTooLarge
		}
	}
	// The decoders are made for damaged files, but a panic of one must
	// not end a conversion.
	defer func() {
		if r := recover(); r != nil {
			img, err = nil, fmt.Errorf("imgconv: %s: %v", format, r)
		}
	}()
	return d.decode(bytes.NewReader(data))
}

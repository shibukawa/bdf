package imgconv

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

// Decode decodes PNG, JPEG, GIF, BMP or WebP bytes (AVIF is not decodable here).
func Decode(data []byte) (image.Image, error) {
	switch Sniff(data) {
	case "png":
		return png.Decode(bytes.NewReader(data))
	case "jpeg":
		return jpeg.Decode(bytes.NewReader(data))
	case "gif":
		return gif.Decode(bytes.NewReader(data))
	case "bmp":
		return bmp.Decode(bytes.NewReader(data))
	case "webp":
		return webp.Decode(bytes.NewReader(data))
	}
	return nil, fmt.Errorf("imgconv: cannot decode %q", Sniff(data))
}

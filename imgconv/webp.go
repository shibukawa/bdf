//go:build !bdf_noconv

package imgconv

import (
	"encoding/binary"
	"errors"
	"image"

	"github.com/shibukawa/bdf/imgconv/internal/webpw"
)

// Available reports whether the WebP codec is compiled in.
func Available() bool { return true }

// encodeWebP encodes NRGBA pixels with libwebp. For lossless, quality is the effort (0-100).
func encodeWebP(img *image.NRGBA, quality, method int, lossless bool) ([]byte, error) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 || w > 16383 || h > 16383 {
		return nil, errors.New("imgconv: image size not supported by WebP")
	}
	pix := img.Pix
	if img.Stride != w*4 {
		pix = make([]byte, w*h*4)
		for y := 0; y < h; y++ {
			copy(pix[y*w*4:(y+1)*w*4], img.Pix[y*img.Stride:y*img.Stride+w*4])
		}
	}
	m := webpw.New()
	webpw.Initialize(m)
	inPtr := webpw.Malloc(m, int32(len(pix)))
	if inPtr == 0 {
		return nil, errors.New("imgconv: webp: out of memory")
	}
	copy(webpw.Memory(m)[inPtr:], pix)
	sizePtr := webpw.Malloc(m, 8)
	ll := int32(0)
	if lossless {
		ll = 1
	}
	out := webpw.Encode(m, inPtr, int32(w), int32(h), sizePtr, int32(quality), int32(method), ll, 0)
	mem := webpw.Memory(m)
	if out == 0 {
		return nil, errors.New("imgconv: webp: encode failed")
	}
	size := binary.LittleEndian.Uint32(mem[sizePtr:])
	res := make([]byte, size)
	copy(res, mem[out:out+int32(size)])
	webpw.Free(m, out)
	webpw.Free(m, sizePtr)
	webpw.Free(m, inPtr)
	return res, nil
}

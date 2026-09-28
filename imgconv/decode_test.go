package imgconv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"runtime"
	"testing"
)

// pngOfSize is a PNG file whose header states a size and whose pixel data
// is a few bytes.
func pngOfSize(w, h uint32) []byte {
	b := []byte("\x89PNG\r\n\x1a\n")
	chunk := func(typ string, data []byte) {
		b = binary.BigEndian.AppendUint32(b, uint32(len(data)))
		b = append(b, typ...)
		b = append(b, data...)
		b = binary.BigEndian.AppendUint32(b, crc32.ChecksumIEEE(b[len(b)-len(data)-4:]))
	}
	ihdr := binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32(nil, w), h)
	chunk("IHDR", append(ihdr, 8, 6, 0, 0, 0)) // 8 bits, RGBA
	chunk("IDAT", []byte{0x78, 0x9c, 0x63, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01})
	chunk("IEND", nil)
	return b
}

// bmpOfSize is a BMP file of 24 bits a pixel without pixels.
func bmpOfSize(w, h int32) []byte {
	b := []byte{'B', 'M'}
	b = binary.LittleEndian.AppendUint32(b, 54)
	b = binary.LittleEndian.AppendUint32(b, 0)
	b = binary.LittleEndian.AppendUint32(b, 54)
	b = binary.LittleEndian.AppendUint32(b, 40)
	b = binary.LittleEndian.AppendUint32(b, uint32(w))
	b = binary.LittleEndian.AppendUint32(b, uint32(h))
	b = binary.LittleEndian.AppendUint16(b, 1)
	b = binary.LittleEndian.AppendUint16(b, 24)
	return append(b, make([]byte, 24)...)
}

// TestDecodeSize decodes images whose headers state sizes that their few
// bytes do not hold: they are refused before memory is taken for them.
func TestDecodeSize(t *testing.T) {
	for name, data := range map[string][]byte{
		"png 20000 x 20000":   pngOfSize(20000, 20000),
		"png 2^29 x 2^29":     pngOfSize(1<<29, 1<<29),
		"png 2^31-1 x 2^31-1": pngOfSize(1<<31-1, 1<<31-1),
		"bmp 16000 x 16000":   bmpOfSize(16000, 16000),
		"bmp 2^31-1 x 2^31-1": bmpOfSize(1<<31-1, 1<<31-1),
	} {
		t.Run(name, func(t *testing.T) {
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			_, err := Decode(data)
			res, oerr := Optimize(data, Options{Mode: Convert})
			runtime.ReadMemStats(&after)
			// a size the decoder itself refuses is its error to report
			if own := name == "png 2^31-1 x 2^31-1"; err == nil || !own && !errors.Is(err, ErrTooLarge) {
				t.Errorf("Decode: %v", err)
			} else if Available() && (oerr == nil || !own && !errors.Is(oerr, ErrTooLarge)) {
				t.Errorf("Optimize: %v", oerr)
			}
			if !bytes.Equal(res.Data, data) || res.Converted {
				t.Error("Optimize does not keep the image as it is")
			}
			if n := after.TotalAlloc - before.TotalAlloc; n > 1<<20 {
				t.Errorf("%d bytes allocated", n)
			}
		})
	}
}

// TestDecodeWithinLimit decodes an image as before.
func TestDecodeWithinLimit(t *testing.T) {
	src := image.NewNRGBA(image.Rect(0, 0, 40, 30))
	for i := range src.Pix {
		src.Pix[i] = byte(i)
	}
	src.SetNRGBA(3, 4, color.NRGBA{1, 2, 3, 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatal(err)
	}
	img, err := Decode(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds() != src.Bounds() || img.At(3, 4) != src.At(3, 4) {
		t.Errorf("bounds %v, pixel %v", img.Bounds(), img.At(3, 4))
	}
}

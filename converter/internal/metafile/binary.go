package metafile

import (
	"encoding/binary"
	"image/color"
	"math"
	"unicode/utf16"
)

func le16(b []byte, i int) uint16 {
	if i < 0 || i+2 > len(b) {
		return 0
	}
	return binary.LittleEndian.Uint16(b[i:])
}

func le32(b []byte, i int) uint32 {
	if i < 0 || i+4 > len(b) {
		return 0
	}
	return binary.LittleEndian.Uint32(b[i:])
}

func lei32(b []byte, i int) float64 { return float64(int32(le32(b, i))) }
func lei16(b []byte, i int) float64 { return float64(int16(le16(b, i))) }
func lef32(b []byte, i int) float64 { return float64(math.Float32frombits(le32(b, i))) }

// slice returns n bytes of a record at off, or nil when they are not there.
func slice(b []byte, off, n uint32) []byte {
	if n == 0 || uint64(off)+uint64(n) > uint64(len(b)) {
		return nil
	}
	return b[off : off+n]
}

// colorRef decodes a COLORREF (0x00BBGGRR).
func colorRef(v uint32) color.NRGBA {
	return color.NRGBA{uint8(v), uint8(v >> 8), uint8(v >> 16), 255}
}

func utf16String(b []byte, n int) []rune {
	if n*2 > len(b) {
		n = len(b) / 2
	}
	u := make([]uint16, n)
	for i := range u {
		u[i] = le16(b, i*2)
	}
	for len(u) > 0 && u[len(u)-1] == 0 {
		u = u[:len(u)-1]
	}
	return utf16.Decode(u)
}

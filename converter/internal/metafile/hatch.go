package metafile

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// hatchBits are 8×8 bitmaps of the GDI hatch styles HS_HORIZONTAL,
// HS_VERTICAL, HS_FDIAGONAL, HS_BDIAGONAL, HS_CROSS and HS_DIAGCROSS (set
// bit = the brush color).
var hatchBits = [...][8]byte{
	{0xff, 0, 0, 0, 0, 0, 0, 0},
	{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80},
	{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01},
	{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80},
	{0xff, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80},
	{0x81, 0x42, 0x24, 0x18, 0x18, 0x24, 0x42, 0x81},
}

type hatchKey struct {
	style  int
	fg, bg color.NRGBA
}

// hatchTile renders a hatch as an 8×8 PNG tile (96 dpi pixels).
func hatchTile(k hatchKey) []byte {
	bits := hatchBits[k.style]
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if bits[y]&(0x80>>x) != 0 {
				img.SetNRGBA(x, y, k.fg)
			} else {
				img.SetNRGBA(x, y, k.bg)
			}
		}
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}

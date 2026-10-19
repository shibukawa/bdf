package raster

import (
	"bytes"
	"image"
	"testing"
)

// FuzzRead checks that any bytes read as PWG Raster give pages or an
// error, without a panic, within the pixels allowed.
func FuzzRead(f *testing.F) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	w.StartPage(Page{Width: 9, Height: 3, XDPI: 300, YDPI: 300, Type: Black1})
	img := image.NewGray(image.Rect(0, 0, 9, 3))
	img.Pix[4] = 255
	w.WriteRows(img)
	w.StartPage(Page{Width: 2, Height: 2, XDPI: 300, YDPI: 300, Type: SRGB8})
	w.WriteRows(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	w.Close()
	f.Add(buf.Bytes())
	f.Add([]byte(Sync))
	f.Fuzz(func(t *testing.T, b []byte) {
		pages, err := Read(bytes.NewReader(b), 1<<16)
		if err != nil {
			return
		}
		for _, p := range pages {
			if p.Width*p.Height > 1<<16 || p.Image.Bounds().Dx() != p.Width || p.Image.Bounds().Dy() != p.Height {
				t.Fatalf("a page of %d × %d with an image of %v", p.Width, p.Height, p.Image.Bounds())
			}
		}
	})
}

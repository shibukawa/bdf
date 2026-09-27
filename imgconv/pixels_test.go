package imgconv

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// pngHeader returns the bit depth and colour type of a PNG file.
func pngHeader(t *testing.T, b []byte) (depth, colorType byte) {
	t.Helper()
	if Sniff(b) != "png" || len(b) < 26 {
		t.Fatalf("not a PNG file: % x", b[:min(8, len(b))])
	}
	return b[24], b[25]
}

func TestEncodePixelsKeep(t *testing.T) {
	bw := image.NewPaletted(image.Rect(0, 0, 64, 32), color.Palette{color.Black, color.White})
	gray := image.NewGray(image.Rect(0, 0, 64, 32))
	grayBW := image.NewGray(image.Rect(0, 0, 64, 32))
	rgbGray := image.NewRGBA(image.Rect(0, 0, 64, 32))
	few := image.NewRGBA(image.Rect(0, 0, 64, 32))
	for i := range gray.Pix {
		gray.Pix[i] = uint8(i)
		grayBW.Pix[i] = uint8(i%2) * 255
		bw.Pix[i] = uint8(i % 2)
		v := uint8(i)
		rgbGray.Pix[4*i], rgbGray.Pix[4*i+1], rgbGray.Pix[4*i+2], rgbGray.Pix[4*i+3] = v, v, v, 255
		few.Pix[4*i], few.Pix[4*i+1], few.Pix[4*i+2], few.Pix[4*i+3] = uint8(i%3)*100, 50, 0, 255
	}
	for _, c := range []struct {
		name             string
		img              image.Image
		depth, colorType byte
	}{
		{"bilevel", bw, 1, 3},
		{"grey", gray, 8, 0},
		{"grey with two levels", grayBW, 1, 3},
		{"RGB that is grey", rgbGray, 8, 0},
		{"three colours", few, 2, 3},
		{"many colours", photoImage(64, 32), 8, 2},
	} {
		r, err := EncodePixels(c.img, true, Options{})
		if err != nil || r.Format != "png" || !r.Lossless {
			t.Fatalf("%s: %+v %v", c.name, r.Format, err)
		}
		if d, ct := pngHeader(t, r.Data); d != c.depth || ct != c.colorType {
			t.Errorf("%s: PNG of depth %d, colour type %d; want %d, %d", c.name, d, ct, c.depth, c.colorType)
		}
		dec, err := png.Decode(bytes.NewReader(r.Data))
		if err != nil {
			t.Fatal(err)
		}
		b := c.img.Bounds()
		for _, p := range []image.Point{{0, 0}, {5, 3}, {63, 31}} {
			want := color.NRGBAModel.Convert(c.img.At(b.Min.X+p.X, b.Min.Y+p.Y))
			if got := color.NRGBAModel.Convert(dec.At(p.X, p.Y)); got != want {
				t.Errorf("%s: pixel %v is %v, want %v", c.name, p, got, want)
			}
		}
	}
	// Pixels from a lossy source become JPEG, unless they are transparent.
	if r, _ := EncodePixels(photoImage(64, 32), false, Options{}); r.Format != "jpeg" || r.Lossless {
		t.Errorf("lossy source stored as %s", r.Format)
	}
	clear := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	if r, _ := EncodePixels(clear, false, Options{}); r.Format != "png" {
		t.Errorf("transparent lossy source stored as %s", r.Format)
	}
}

func TestEncodePixelsConvert(t *testing.T) {
	if !Available() {
		t.Skip("codecs not compiled in")
	}
	bw := image.NewPaletted(image.Rect(0, 0, 256, 256), color.Palette{color.Black, color.White})
	for i := range bw.Pix {
		if (i/256/16+i%256/16)%2 == 0 {
			bw.Pix[i] = 1
		}
	}
	keep, _ := EncodePixels(bw, true, Options{})
	r, err := EncodePixels(bw, true, Options{Mode: Convert})
	if err != nil || r.Format != "webp" || !r.Lossless || !r.Converted || len(r.Data) >= len(keep.Data) {
		t.Errorf("bilevel under Convert: %s lossless %v, %d bytes (PNG %d), %v", r.Format, r.Lossless, len(r.Data), len(keep.Data), err)
	}
}

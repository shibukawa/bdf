package imgconv

import (
	"image"
	"image/color"
	"testing"
)

func TestFitSize(t *testing.T) {
	a4w, a4h := 595.2756, 841.8898 // points
	for _, c := range []struct {
		name         string
		opts         Options
		w, h         int
		wPt, hPt     float64
		wantW, wantH int
	}{
		{"A4 at 300 dpi to the default 192", Options{}, 2480, 3508, a4w, a4h, 1587, 2245},
		{"A4 at 150 dpi is kept", Options{}, 1240, 1754, a4w, a4h, 1240, 1754},
		{"exactly 192 dpi is kept", Options{}, 1587, 2245, 1587.0 / 192 * 72, 2245.0 / 192 * 72, 1587, 2245},
		{"a fax at 204 × 98 dpi: only x is capped", Options{}, 1728, 1078, 1728.0 / 204 * 72, 1078.0 / 98 * 72, 1626, 1078},
		{"A0 at 192 dpi: the pixel cap", Options{}, 6358, 8988, 2383.9, 3370.4, 3229, 4565},
		{"unknown page size: the pixel cap only", Options{}, 10000, 5000, 0, 0, 5430, 2715},
		{"no dpi cap", Options{MaxDPI: -1}, 2480, 3508, a4w, a4h, 2480, 3508},
		{"no caps", Options{MaxDPI: -1, MaxPixels: -1}, 20000, 20000, 100, 100, 20000, 20000},
		{"custom caps", Options{MaxDPI: 100, MaxPixels: 1000}, 300, 300, 72, 72, 31, 31},
	} {
		w, h := c.opts.FitSize(c.w, c.h, c.wPt, c.hPt)
		if w != c.wantW || h != c.wantH {
			t.Errorf("%s: %d×%d, want %d×%d", c.name, w, h, c.wantW, c.wantH)
		}
		if mp := c.opts.maxPixels(); mp > 0 && w*h > mp {
			t.Errorf("%s: %d pixels over the cap %d", c.name, w*h, mp)
		}
	}
}

func TestResizeAverages(t *testing.T) {
	// 3 → 2 pixels: the target pixels cover 1.5 source pixels each.
	src := image.NewGray(image.Rect(0, 0, 3, 1))
	copy(src.Pix, []byte{0, 90, 255})
	got := Resize(src, 2, 1).(*image.Gray)
	// (0·2 + 90·1) / 3 = 30; (90·1 + 255·2) / 3 = 200
	if got.Pix[0] != 30 || got.Pix[1] != 200 {
		t.Errorf("averaged to %v, want [30 200]", got.Pix)
	}
	// RGBA averages premultiplied: a transparent pixel adds no colour.
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 1))
	copy(rgba.Pix, []byte{0, 0, 0, 0, 200, 100, 50, 255})
	one := Resize(rgba, 1, 1).(*image.RGBA)
	if c := one.RGBAAt(0, 0); c != (color.RGBA{100, 50, 25, 128}) {
		t.Errorf("RGBA average %v", c)
	}
	// The same size, or larger, is left alone.
	if Resize(src, 3, 1) != image.Image(src) || Resize(src, 4, 2) != image.Image(src) {
		t.Error("an image at its size was resampled")
	}
}

func TestResizeBilevel(t *testing.T) {
	// Black strokes one source pixel wide every fourth column, scaled to
	// half: each target pixel is half covered by ink, which is more than
	// two fifths, so the strokes stay (thinner, but there).
	src := image.NewPaletted(image.Rect(0, 0, 16, 4), color.Palette{color.Black, color.White})
	for i := range src.Pix {
		src.Pix[i] = 1
	}
	for y := range 4 {
		for x := 0; x < 16; x += 4 {
			src.SetColorIndex(x, y, 0)
		}
	}
	got, ok := Resize(src, 8, 2).(*image.Paletted)
	if !ok || len(got.Palette) != 2 {
		t.Fatalf("a bilevel image did not stay bilevel: %T", got)
	}
	for x := range 8 {
		want := uint8(1)
		if x%2 == 0 {
			want = 0
		}
		if got.ColorIndexAt(x, 0) != want {
			t.Errorf("column %d: index %d, want %d", x, got.ColorIndexAt(x, 0), want)
		}
	}
	// Ink is the less frequent colour: white text on black keeps its text.
	for i := range src.Pix {
		src.Pix[i] ^= 1
	}
	inv := Resize(src, 8, 2).(*image.Paletted)
	for x := range 8 {
		if inv.ColorIndexAt(x, 0) != 1-got.ColorIndexAt(x, 0) {
			t.Fatalf("inverted image: column %d differs", x)
		}
	}
}

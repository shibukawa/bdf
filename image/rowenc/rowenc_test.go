package rowenc

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"testing"

	"golang.org/x/image/tiff"
)

// images makes an RGBA image (with translucent pixels), a gray one and a
// bilevel one of w × h pixels.
func images(w, h int) (*image.RGBA, *image.Gray, *image.Paletted) {
	rng := rand.New(rand.NewSource(1))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	gray := image.NewGray(image.Rect(0, 0, w, h))
	bi := image.NewPaletted(image.Rect(0, 0, w, h), color.Palette{color.Gray{0}, color.Gray{255}})
	for i := 0; i < w*h; i++ {
		a := uint8(rng.Intn(256))
		if rng.Intn(4) == 0 {
			a = 255
		}
		// premultiplied: no channel above the alpha
		rgba.Pix[4*i], rgba.Pix[4*i+1], rgba.Pix[4*i+2], rgba.Pix[4*i+3] = uint8(rng.Intn(int(a)+1)), uint8(rng.Intn(int(a)+1)), a-uint8(rng.Intn(3))*(a/3), a
		gray.Pix[i] = uint8(rng.Intn(256))
		if (i/w)%5 == 0 || rng.Intn(3) == 0 {
			bi.Pix[i] = 1
		}
	}
	return rgba, gray, bi
}

// rows yields an image in bands of rows of different heights.
func rows(img image.Image, w *testing.T, enc Writer) {
	b := img.Bounds()
	for y, n := b.Min.Y, 1; y < b.Max.Y; y, n = y+n, n*2+1 {
		r := image.Rect(b.Min.X, y, b.Max.X, min(y+n, b.Max.Y))
		var band image.Image
		switch img := img.(type) {
		case *image.RGBA:
			band = img.SubImage(r)
		case *image.Gray:
			band = img.SubImage(r)
		case *image.Paletted:
			band = img.SubImage(r)
		}
		if err := enc.WriteRows(band); err != nil {
			w.Fatal(err)
		}
	}
}

func TestPNGOfRows(t *testing.T) {
	for _, size := range [][2]int{{1, 1}, {7, 3}, {33, 70}, {300, 200}} {
		rgba, gray, bi := images(size[0], size[1])
		for _, tc := range []struct {
			kind Kind
			img  image.Image
		}{{RGBA, rgba}, {Gray, gray}, {Bilevel, bi}} {
			var buf, whole bytes.Buffer
			enc, err := NewPNG(&buf, size[0], size[1], tc.kind, &Options{DPI: 600})
			if err != nil {
				t.Fatal(err)
			}
			rows(tc.img, t, enc)
			if err := enc.Close(); err != nil {
				t.Fatal(err)
			}
			if err := enc.Close(); err == nil {
				t.Error("closed twice")
			}
			got, err := png.Decode(&buf)
			if err != nil {
				t.Fatalf("%dx%d kind %d: %v", size[0], size[1], tc.kind, err)
			}
			// the pixels of image/png's own file
			if err := png.Encode(&whole, tc.img); err != nil {
				t.Fatal(err)
			}
			want, err := png.Decode(&whole)
			if err != nil {
				t.Fatal(err)
			}
			if got.Bounds() != want.Bounds() {
				t.Fatalf("bounds %v, want %v", got.Bounds(), want.Bounds())
			}
			for y := 0; y < size[1]; y++ {
				for x := 0; x < size[0]; x++ {
					if g, w := got.At(x, y), want.At(x, y); g != w {
						t.Fatalf("%dx%d kind %d: pixel (%d, %d) is %v, want %v", size[0], size[1], tc.kind, x, y, g, w)
					}
				}
			}
			if tc.kind == Bilevel {
				if p, ok := got.(*image.Paletted); !ok || len(p.Palette) != 2 {
					t.Errorf("a bilevel PNG decodes to %T", got)
				}
			}
		}
	}
	// a gray image as bilevel rows
	_, gray, _ := images(20, 9)
	var buf bytes.Buffer
	enc, err := NewPNG(&buf, 20, 9, Bilevel, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := enc.WriteRows(gray); err != nil {
		t.Fatal(err)
	}
	if err := enc.Close(); err != nil {
		t.Fatal(err)
	}
	got, err := png.Decode(&buf)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range gray.Pix {
		want := color.Gray{0}
		if v >= 128 {
			want = color.Gray{255}
		}
		if g := color.GrayModel.Convert(got.At(i%20, i/20)); g != want {
			t.Fatalf("pixel %d of gray %d is %v", i, v, g)
		}
	}
}

func TestTIFFOfRows(t *testing.T) {
	for _, size := range [][2]int{{1, 1}, {7, 3}, {33, 70}, {300, 200}} {
		rgba, gray, bi := images(size[0], size[1])
		for _, tc := range []struct {
			kind Kind
			img  image.Image
		}{{RGBA, rgba}, {Gray, gray}, {Bilevel, bi}} {
			var buf bytes.Buffer
			enc, err := NewTIFF(&buf, size[0], size[1], tc.kind, &Options{DPI: 600})
			if err != nil {
				t.Fatal(err)
			}
			rows(tc.img, t, enc)
			if err := enc.Close(); err != nil {
				t.Fatal(err)
			}
			got, err := tiff.Decode(bytes.NewReader(buf.Bytes()))
			if err != nil {
				t.Fatalf("%dx%d kind %d: %v", size[0], size[1], tc.kind, err)
			}
			if got.Bounds() != tc.img.Bounds() {
				t.Fatalf("bounds %v, want %v", got.Bounds(), tc.img.Bounds())
			}
			for y := 0; y < size[1]; y++ {
				for x := 0; x < size[0]; x++ {
					g, w := got.At(x, y), tc.img.At(x, y)
					if tc.kind == RGBA {
						// the reader unmultiplies the alpha: compare as RGBA
						g, w = color.RGBAModel.Convert(g), color.RGBAModel.Convert(w)
						gr, gg, gb, ga := g.RGBA()
						wr, wg, wb, wa := w.RGBA()
						if ga != wa || diff(gr, wr) > 0x101 || diff(gg, wg) > 0x101 || diff(gb, wb) > 0x101 {
							t.Fatalf("%dx%d: pixel (%d, %d) is %v, want %v", size[0], size[1], x, y, g, w)
						}
						continue
					}
					if color.GrayModel.Convert(g) != color.GrayModel.Convert(w) {
						t.Fatalf("%dx%d kind %d: pixel (%d, %d) is %v, want %v", size[0], size[1], tc.kind, x, y, g, w)
					}
				}
			}
		}
	}
}

func diff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestRowsErrors(t *testing.T) {
	rgba, gray, _ := images(10, 4)
	if _, err := NewPNG(&bytes.Buffer{}, 0, 4, RGBA, nil); err == nil {
		t.Error("an empty image")
	}
	if _, err := NewPNG(&bytes.Buffer{}, 10, 4, Kind(7), nil); err == nil {
		t.Error("an unknown kind")
	}
	if _, err := NewPNG(&bytes.Buffer{}, 10, 4, RGBA, &Options{Level: 12}); err == nil {
		t.Error("a bad level")
	}
	for _, format := range []string{"png", "tiff"} {
		enc, err := New(&bytes.Buffer{}, format, 10, 4, RGBA, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := enc.WriteRows(gray); err == nil {
			t.Errorf("%s: gray rows for an RGBA image", format)
		}
		enc, _ = New(&bytes.Buffer{}, format, 10, 4, RGBA, nil)
		if err := enc.WriteRows(rgba.SubImage(image.Rect(0, 0, 5, 2))); err == nil {
			t.Errorf("%s: rows too narrow", format)
		}
		enc, _ = New(&bytes.Buffer{}, format, 10, 4, RGBA, nil)
		if err := enc.WriteRows(rgba.SubImage(image.Rect(0, 0, 10, 2))); err != nil {
			t.Fatal(err)
		}
		if err := enc.Close(); err == nil {
			t.Errorf("%s: closed after 2 of 4 rows", format)
		}
		enc, _ = New(&bytes.Buffer{}, format, 10, 4, RGBA, nil)
		if err := enc.WriteRows(rgba); err != nil {
			t.Fatal(err)
		}
		if err := enc.WriteRows(rgba); err == nil {
			t.Errorf("%s: rows past the end", format)
		}
	}
	if _, err := New(&bytes.Buffer{}, "gif", 10, 4, RGBA, nil); err == nil {
		t.Error("an unknown format")
	}
	if KindOf(rgba) != RGBA || KindOf(gray) != Gray {
		t.Error("KindOf")
	}
	var buf bytes.Buffer
	if err := Encode(&buf, "png", gray, Gray, nil); err != nil {
		t.Fatal(err)
	}
	if img, err := png.Decode(&buf); err != nil || img.Bounds() != gray.Bounds() {
		t.Errorf("Encode: %v %v", img, err)
	}
}

// TestBandsWriteTheSameFile checks that an image written in bands is
// the file of the image written whole: the filters look at a row and
// the one before, and the compression at the bytes, not at the writes.
func TestBandsWriteTheSameFile(t *testing.T) {
	rgba, gray, bi := images(257, 301)
	for _, tc := range []struct {
		kind Kind
		img  image.Image
	}{{RGBA, rgba}, {Gray, gray}, {Bilevel, bi}} {
		for _, format := range []string{"png", "tiff"} {
			var whole, banded bytes.Buffer
			if err := Encode(&whole, format, tc.img, tc.kind, &Options{Level: 9}); err != nil {
				t.Fatal(err)
			}
			enc, err := New(&banded, format, 257, 301, tc.kind, &Options{Level: 9})
			if err != nil {
				t.Fatal(err)
			}
			rows(tc.img, t, enc)
			if err := enc.Close(); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(whole.Bytes(), banded.Bytes()) {
				t.Errorf("%s of kind %d: the file of the bands differs from that of the whole", format, tc.kind)
			}
		}
	}
}

package raster

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"testing"

	"github.com/shibukawa/bdf"
)

// pattern is a page of stripes and a gradient, different for each n.
func pattern(w, h, n int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			o := img.PixOffset(x, y)
			v := uint8((x*7 + y*3 + n*40) % 256)
			if (x/5+y/7+n)%3 == 0 {
				v = 0
			}
			img.Pix[o], img.Pix[o+1], img.Pix[o+2], img.Pix[o+3] = v, 255-v, uint8(x*n), 255
		}
	}
	return img
}

// rasterOf writes pages of the patterns in a format.
func rasterOf(t *testing.T, urf bool, typ Type, pages, w, h, dpi int) ([]byte, []image.Image) {
	t.Helper()
	var buf bytes.Buffer
	var wr *Writer
	if urf {
		wr = NewURFWriter(&buf, pages)
	} else {
		wr = NewWriter(&buf)
	}
	var imgs []image.Image
	for n := range pages {
		img := pattern(w, h, n)
		imgs = append(imgs, img)
		if err := wr.StartPage(Page{Width: w, Height: h, XDPI: dpi, YDPI: dpi, Type: typ, PageWidth: w * 72 / dpi, PageHeight: h * 72 / dpi,
			Duplex: true, Tumble: true}); err != nil {
			t.Fatal(err)
		}
		if err := wr.WriteRows(img); err != nil {
			t.Fatal(err)
		}
	}
	if err := wr.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes(), imgs
}

// TestDocument reads raster streams as documents of pictures: the pixels
// of each page are those written (as the reader of the tests reads them),
// and the pages as large as the paper.
func TestDocument(t *testing.T) {
	cases := []struct {
		urf bool
		typ Type
	}{{false, Black1}, {false, SGray8}, {false, SRGB8}, {true, SGray8}, {true, SRGB8}}
	for _, c := range cases {
		b, _ := rasterOf(t, c.urf, c.typ, 3, 150, 97, 150)
		var want []*ReadPage
		var err error
		if c.urf {
			want, err = ReadURF(bytes.NewReader(b), 0)
		} else {
			want, err = Read(bytes.NewReader(b), 0)
		}
		if err != nil {
			t.Fatal(err)
		}
		doc, info, err := Document(bytes.NewReader(b), &DocumentOptions{Title: "scan"})
		if err != nil {
			t.Fatalf("%v %s: %v", c.urf, c.typ, err)
		}
		if info.Pages != 3 || !info.Duplex || !info.Tumble || info.Resolution != 150 || (info.Format == FormatURF) != c.urf {
			t.Errorf("%v %s: %+v", c.urf, c.typ, info)
		}
		v := doc.Views[0]
		if len(v.Pages) != 3 || v.Pages[0].W != 72 || v.Pages[0].H != 46 && v.Pages[0].H != 47 {
			t.Fatalf("%v %s: %d pages of %v × %v", c.urf, c.typ, len(v.Pages), v.Pages[0].W, v.Pages[0].H)
		}
		var got []image.Image
		for _, p := range doc.Parts() {
			if p.Type != bdf.PartImage {
				continue
			}
			img, err := png.Decode(bytes.NewReader(p.Data))
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, img)
		}
		if len(got) != 3 {
			t.Fatalf("%v %s: %d images", c.urf, c.typ, len(got))
		}
		for i := range got {
			if n, most := differ(want[i].Image, got[i]); n > 0 {
				t.Errorf("%v %s page %d: %d pixels differ (by %d at most)", c.urf, c.typ, i+1, n, most)
			}
		}
	}
}

// TestPageReaderBands reads a page in bands that do not divide it, and
// skips the rest of a page not read.
func TestPageReaderBands(t *testing.T) {
	b, _ := rasterOf(t, false, SRGB8, 2, 31, 50, 72)
	want, _ := Read(bytes.NewReader(b), 0)
	pr, err := NewPageReader(bytes.NewReader(b), 0)
	if err != nil {
		t.Fatal(err)
	}
	p, err := pr.Next()
	if err != nil || p.Width != 31 || p.Height != 50 {
		t.Fatal(p, err)
	}
	all := image.NewRGBA(image.Rect(0, 0, 31, 50))
	y := 0
	for {
		band := image.NewRGBA(image.Rect(0, 0, 31, 7))
		n, err := pr.ReadRows(band)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		for r := 0; r < n; r++ {
			copy(all.Pix[all.PixOffset(0, y+r):all.PixOffset(0, y+r+1)], band.Pix[band.PixOffset(0, r):band.PixOffset(0, r+1)])
		}
		y += n
	}
	if y != 50 {
		t.Fatalf("%d rows", y)
	}
	if n, _ := differ(want[0].Image, all); n > 0 {
		t.Errorf("%d pixels differ", n)
	}
	// the second page is skipped: Next ends
	if _, err := pr.Next(); err != nil {
		t.Fatal(err)
	}
	if _, err := pr.Next(); err != io.EOF {
		t.Fatalf("after the last page: %v", err)
	}
}

func TestDocumentErrors(t *testing.T) {
	b, _ := rasterOf(t, false, SGray8, 2, 40, 40, 72)
	for name, in := range map[string][]byte{
		"empty":     nil,
		"not":       []byte("%PDF-1.7\n"),
		"cut":       b[:len(b)-30],
		"header":    b[:4+100],
		"no pages":  []byte(Sync),
		"urf count": append([]byte(URFSync), 0, 0, 0, 2),
	} {
		if _, _, err := Document(bytes.NewReader(in), nil); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, _, err := Document(bytes.NewReader(b), &DocumentOptions{MaxPages: 1}); err == nil {
		t.Error("more pages than MaxPages")
	}
	if _, _, err := Document(bytes.NewReader(b), &DocumentOptions{MaxPixels: 100}); err == nil {
		t.Error("a page larger than MaxPixels")
	}
}

func FuzzDocument(f *testing.F) {
	for _, urf := range []bool{false, true} {
		b, _ := rasterOfF(urf)
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		Document(bytes.NewReader(b), &DocumentOptions{MaxPixels: 1 << 16, MaxPages: 8})
	})
}

func rasterOfF(urf bool) ([]byte, error) {
	var buf bytes.Buffer
	var wr *Writer
	if urf {
		wr = NewURFWriter(&buf, 1)
	} else {
		wr = NewWriter(&buf)
	}
	img := pattern(16, 9, 1)
	if err := wr.StartPage(Page{Width: 16, Height: 9, XDPI: 72, YDPI: 72, Type: SRGB8}); err != nil {
		return nil, err
	}
	if err := wr.WriteRows(img); err != nil {
		return nil, err
	}
	err := wr.Close()
	return buf.Bytes(), err
}

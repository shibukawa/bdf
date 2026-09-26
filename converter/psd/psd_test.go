package psd

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// withoutComposite marks a file as saved without "Maximize Compatibility":
// its version information says the image data is not the composite.
func withoutComposite(t *testing.T, data []byte) []byte {
	t.Helper()
	i := bytes.Index(data, []byte("8BIM\x04\x21\x00\x00"))
	if i < 0 {
		t.Fatal("no version information")
	}
	out := bytes.Clone(data)
	out[i+8+4+4] = 0 // signature, ID, name; size; version; hasRealMergedData
	return out
}

type converted struct {
	res   *Result
	r     *bdf.Reader
	pages []*image.NRGBA
}

func convert(t *testing.T, data []byte, opts *Options) converted {
	t.Helper()
	res, err := Convert(bytes.NewReader(data), int64(len(data)), opts)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	c := converted{res: res, r: r}
	for _, p := range r.Manifest.Views[0].Pages {
		o, err := r.Object(p.Layers[0].Obj)
		if err != nil {
			t.Fatal(err)
		}
		if len(o.Images) != 1 {
			t.Fatalf("page draws %d images", len(o.Images))
		}
		b, err := r.Part(o.Images[0])
		if err != nil {
			t.Fatal(err)
		}
		img, err := png.Decode(bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}
		n := image.NewNRGBA(img.Bounds())
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				n.Set(x, y, img.At(x, y))
			}
		}
		c.pages = append(c.pages, n)
	}
	return c
}

// sameImage compares premultiplied samples: colours under transparent
// pixels do not count.
func sameImage(t *testing.T, name string, got, want *image.NRGBA, tol int) {
	t.Helper()
	if got.Bounds() != want.Bounds() {
		t.Fatalf("%s: %v, want %v", name, got.Bounds(), want.Bounds())
	}
	bad, worst := 0, 0
	for i := range got.Pix {
		pm := func(p []uint8) int {
			if i%4 == 3 {
				return int(p[i])
			}
			return int(p[i]) * int(p[i-i%4+3]) / 255
		}
		d := pm(got.Pix) - pm(want.Pix)
		if d < 0 {
			d = -d
		}
		if d > tol {
			bad++
		}
		worst = max(worst, d)
	}
	if bad > 0 {
		t.Errorf("%s: %d samples differ by more than %d (at most %d)", name, bad, tol, worst)
	}
}

func rgbaAt(img *image.NRGBA, x, y int) [4]uint8 {
	i := img.PixOffset(x, y)
	return [4]uint8(img.Pix[i : i+4])
}

func near(c [4]uint8, want [4]uint8, tol int) bool {
	for i := range c {
		d := int(c[i]) - int(want[i])
		if d < -tol || d > tol {
			return false
		}
	}
	return true
}

func TestConvertLayers(t *testing.T) {
	data := read(t, "layers.psd")
	if f := converter.Detect(bytes.NewReader(data), int64(len(data))); f == nil || f.Name != "psd" {
		t.Fatalf("detected as %v", f)
	}
	c := convert(t, data, nil)
	if len(c.res.Warnings) != 0 || c.res.Composited {
		t.Errorf("warnings %v, composited %v", c.res.Warnings, c.res.Composited)
	}
	m := c.r.Manifest
	if m.Meta.Source != "psd" || m.Meta.DC.Title.First() != "Layers" || m.Meta.DC.Creator.First() != "test/psd/gen.py" ||
		m.Meta.DC.Created.First() != "2026-09-26T12:00:00+09:00" || len(m.Meta.DC.Format) != 0 {
		t.Errorf("meta = %+v", m.Meta)
	}
	// 160 × 120 pixels at 144 ppi
	if p := m.Views[0].Pages; len(p) != 1 || p[0].W != 80 || p[0].H != 60 {
		t.Fatalf("pages = %+v", p)
	}
	if b := c.pages[0].Bounds(); b.Dx() != 160 || b.Dy() != 120 {
		t.Errorf("image %v", b)
	}
	// The hidden red layer covers everything; it is not drawn.
	if px := rgbaAt(c.pages[0], 2, 2); px[0] > 250 && px[1] < 5 {
		t.Errorf("hidden layer drawn: %v", px)
	}

	// Without the composite, the converter composites the layers: groups,
	// clipping, masks, a fill layer, blend modes, fill opacity.
	d := convert(t, withoutComposite(t, data), nil)
	if !d.res.Composited || len(d.res.Warnings) != 1 || !strings.Contains(d.res.Warnings[0], "Maximize Compatibility") {
		t.Errorf("composited %v, warnings %v", d.res.Composited, d.res.Warnings)
	}
	sameImage(t, "layers composited", d.pages[0], c.pages[0], 2)
}

func TestConvertArtboards(t *testing.T) {
	data := read(t, "artboards.psd")
	c := convert(t, data, nil)
	if len(c.res.Warnings) != 0 {
		t.Errorf("warnings: %v", c.res.Warnings)
	}
	// Three visible artboards, bottom up; the hidden one is left out. 72 ppi:
	// a point per pixel.
	pages := c.r.Manifest.Views[0].Pages
	want := [][2]float32{{200, 150}, {180, 120}, {160, 120}}
	if c.res.Artboards != 3 || len(pages) != 3 {
		t.Fatalf("%d artboards, %d pages", c.res.Artboards, len(pages))
	}
	for i, p := range pages {
		if p.W != want[i][0] || p.H != want[i][1] {
			t.Errorf("artboard %d: %v × %v, want %v", i+1, p.W, p.H, want[i])
		}
	}
	for _, e := range []struct {
		page, x, y int
		want       [4]uint8
	}{
		{0, 2, 2, [4]uint8{255, 255, 255, 255}},   // white background
		{1, 2, 110, [4]uint8{40, 70, 140, 255}},   // coloured background
		{2, 2, 2, [4]uint8{0, 0, 0, 0}},           // transparent background
		{2, 80, 60, [4]uint8{255, 255, 255, 255}}, // the hole of the ring
	} {
		if got := rgbaAt(c.pages[e.page], e.x, e.y); !near(got, e.want, 1) && !(e.want[3] == 0 && got[3] == 0) {
			t.Errorf("artboard %d at %d,%d: %v, want %v", e.page+1, e.x, e.y, got, e.want)
		}
	}
	// The disc that runs over the right edge of artboard 1 is cut there.
	if got := rgbaAt(c.pages[0], 199, 90); got[3] != 255 || got[2] < 200 || got[0] > 100 {
		t.Errorf("artboard 1 at its edge: %v, want the blue disc", got)
	}

	// Pages select artboards; without artboards the canvas is one page.
	sel := convert(t, data, &Options{Pages: []int{3, 1}})
	if p := sel.r.Manifest.Views[0].Pages; len(p) != 2 || p[0].W != 160 || p[1].W != 200 {
		t.Errorf("selected pages: %+v", p)
	}
	whole := convert(t, data, &Options{NoArtboards: true})
	if p := whole.r.Manifest.Views[0].Pages; len(p) != 1 || p[0].W != 400 || p[0].H != 260 || whole.res.Artboards != 0 {
		t.Errorf("no artboards: %+v", p)
	}
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), &Options{Pages: []int{4}}); err == nil {
		t.Error("page 4 of 3 converted")
	}

	// Composited by the converter: backgrounds, clipping to the artboard,
	// the hidden artboard left out.
	d := convert(t, withoutComposite(t, data), nil)
	for i := range c.pages {
		sameImage(t, "artboard composited", d.pages[i], c.pages[i], 2)
	}
}

func TestColorModes(t *testing.T) {
	type check struct {
		x, y int
		want [4]uint8
	}
	for _, c := range []struct {
		file   string
		checks []check
		tol    int
	}{
		// 1 bit: 1 is black
		{"bitmap.psd", []check{{0, 0, [4]uint8{0, 0, 0, 255}}, {4, 0, [4]uint8{255, 255, 255, 255}}, {4, 4, [4]uint8{0, 0, 0, 255}}}, 0},
		// the colour table, index 3 transparent
		{"indexed.psd", []check{{0, 0, [4]uint8{255, 0, 0, 255}}, {6, 0, [4]uint8{0, 255, 0, 255}},
			{12, 0, [4]uint8{0, 0, 255, 255}}, {18, 0, [4]uint8{255, 255, 255, 0}}}, 0},
		// Lab: white, black, and red (L 54, a 81, b 70 is about sRGB 255 0 0)
		{"lab.psd", []check{{0, 0, [4]uint8{255, 255, 255, 255}}, {1, 0, [4]uint8{0, 0, 0, 255}}, {2, 0, [4]uint8{255, 0, 0, 255}}}, 20},
		// CMYK with transparency: the yellow disc on nothing
		{"cmyk.psd", []check{{14, 6, [4]uint8{230, 230, 51, 191}}, {0, 0, [4]uint8{0, 0, 0, 0}}}, 6}, // CMYK and the matte round
	} {
		img := convert(t, read(t, c.file), nil).pages[0]
		for _, k := range c.checks {
			got := rgbaAt(img, k.x, k.y)
			if k.want[3] == 0 && got[3] == 0 {
				continue
			}
			if !near(got, k.want, c.tol) {
				t.Errorf("%s at %d,%d: %v, want %v", c.file, k.x, k.y, got, k.want)
			}
		}
	}

	// The same layers in other depths and in the large document format:
	// composite and composited layers agree, and agree with 8 bits.
	ref := convert(t, read(t, "large.psb"), nil).pages[0]
	for _, name := range []string{"large.psb", "rgb32.psd", "gray16.psd"} {
		data := read(t, name)
		c := convert(t, data, nil)
		d := convert(t, withoutComposite(t, data), nil)
		if len(d.res.Warnings) != 1 {
			t.Errorf("%s: warnings %v", name, d.res.Warnings)
		}
		sameImage(t, name+" composited", d.pages[0], c.pages[0], 2)
		if name == "rgb32.psd" {
			sameImage(t, name+" against 8 bits", c.pages[0], ref, 2)
		}
	}
	gray := convert(t, read(t, "gray16.psd"), nil).pages[0]
	if px := rgbaAt(gray, 0, 0); px[0] != px[1] || px[1] != px[2] || px[3] != 255 {
		t.Errorf("gray16: %v", px)
	}
}

func TestNotPSD(t *testing.T) {
	for _, data := range [][]byte{[]byte("8BPS\x00\x03"), []byte("hello"), []byte("8BPS\x00\x01\x00\x00\x00\x00\x00\x00\x00\x03\x00\x00\x00\x10\x00\x00\x00\x10\x00\x08\x00\x03")} {
		if _, err := Convert(bytes.NewReader(data), int64(len(data)), nil); err == nil {
			t.Errorf("%q converted", data)
		}
	}
}

// TestCorrupt converts damaged files: they may fail or warn, never panic.
func TestCorrupt(t *testing.T) {
	seed := uint32(1)
	rand := func(n int) int {
		seed = seed*1664525 + 1013904223
		return int(seed>>8) % n
	}
	for _, name := range []string{"layers.psd", "artboards.psd", "gray16.psd", "rgb32.psd", "large.psb", "indexed.psd"} {
		data := read(t, name)
		for _, variant := range [][]byte{data, withoutComposite(t, data)} {
			for i := 0; i < 60; i++ {
				b := bytes.Clone(variant)
				if i%2 == 0 {
					b = b[:rand(len(b))]
				} else {
					for k := 0; k < 8; k++ {
						b[rand(len(b))] = byte(rand(256))
					}
				}
				Convert(bytes.NewReader(b), int64(len(b)), &Options{Warn: func(string) {}})
			}
		}
	}
}

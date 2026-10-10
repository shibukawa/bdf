package psd

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/image/imgconv"
)

// The layers are composited at the resolution the pages are stored at: the
// tests compare that with the composite of the whole canvas, reduced
// afterwards as the pages of a composite in the file are.

// whole are the options without a resolution cap, and capped returns those
// that store a page with at most so many pixels.
var whole = &Options{Images: imgconv.Options{MaxDPI: -1, MaxPixels: -1}}

func capped(pixels int) *Options {
	return &Options{Images: imgconv.Options{MaxDPI: -1, MaxPixels: pixels}}
}

// reducedAfter composites a document at the size of its canvas and reduces
// the page.
func reducedAfter(t *testing.T, data []byte, w, h int) *image.NRGBA {
	t.Helper()
	small := imgconv.Resize(convert(t, data, whole).pages[0], w, h)
	n := image.NewNRGBA(small.Bounds())
	for y := range h {
		for x := range w {
			n.Set(x, y, small.At(x, y))
		}
	}
	return n
}

// grain gives samples that differ from pixel to pixel.
func grain(c, x, y int) uint8 { return uint8(40 + 50*c + (x*37+y*91)%97) }

// TestCompositeReduced: layers that blend normally, with the edges and the
// transparency of whole reduced pixels, give what the composite of the
// canvas gives when it is reduced.
func TestCompositeReduced(t *testing.T) {
	data := layered(64, 48, []testLayer{
		{bottom: 48, right: 64, shade: grain},
		{top: 8, left: 8, bottom: 32, right: 40, shade: func(c, x, y int) uint8 { return grain(2-c, y, x) }, opacity: 128},
		// transparency that changes from one reduced pixel to the next
		{top: 16, left: 20, bottom: 44, right: 60, rgb: [3]uint8{250, 200, 20}, alpha: func(x, y int) uint8 { return uint8(x / 4 * 16) }},
		// a group with an opacity of its own, which reaches over the canvas
		{section: 3},
		{top: 28, left: -8, bottom: 56, right: 28, shade: grain},
		{section: 1, opacity: 200},
	})
	c := convert(t, data, capped(16*12))
	if b := c.pages[0].Bounds(); b.Dx() != 16 || b.Dy() != 12 || c.res.Scaled != 1 || len(c.res.Warnings) != 1 {
		t.Fatalf("page of %v, %d scaled, warnings %q", b, c.res.Scaled, c.res.Warnings)
	}
	if p := c.r.Manifest.Views[0].Pages[0]; p.W != 64 || p.H != 48 {
		t.Errorf("page of %v × %v points", p.W, p.H)
	}
	sameImage(t, "reduced composite", c.pages[0], reducedAfter(t, data, 16, 12), 2)

	// A canvas whose sides are reduced differently: 144 by 72 pixels per
	// inch, stored at 72.
	res := []byte("8BIM\x03\xed\x00\x00\x00\x00\x00\x10\x00\x90\x00\x00\x00\x01\x00\x02\x00\x48\x00\x00\x00\x01\x00\x02")
	i := bytes.Index(data, []byte("8BIM\x04\x21"))
	wide := bytes.Clone(data[:i-4])
	wide = append(append(append(wide, be32(len(res)+18)...), res...), data[i:]...)
	c = convert(t, wide, &Options{Images: imgconv.Options{MaxDPI: 72}})
	if b := c.pages[0].Bounds(); b.Dx() != 32 || b.Dy() != 48 {
		t.Fatalf("page of %v", b)
	}
	full := convert(t, wide, whole).pages[0]
	side := imgconv.Resize(full, 32, 48)
	n := image.NewNRGBA(side.Bounds())
	for y := range 48 {
		for x := range 32 {
			n.Set(x, y, side.At(x, y))
		}
	}
	sameImage(t, "one side reduced", c.pages[0], n, 2)
}

// TestPlacement: rectangles are placed by their edges, each at the nearest
// boundary of the reduced pixels.
func TestPlacement(t *testing.T) {
	for v, want := range map[int]int{0: 0, 1: 0, 2: 1, 5: 1, 6: 2, 30: 8, 64: 16, -1: 0, -2: 0, -3: -1, -6: -1, -7: -2} {
		if got := edge(v, 16, 64); got != want {
			t.Errorf("edge of %d: %d, want %d", v, got, want)
		}
	}
	red, blue, green := [3]uint8{255, 0, 0}, [3]uint8{0, 0, 255}, [3]uint8{0, 255, 0}
	// Layers that meet between two reduced pixels meet there: no pixel
	// shows what is below them, none has both.
	c := convert(t, layered(64, 48, []testLayer{
		{bottom: 48, right: 64, rgb: green},
		{bottom: 48, right: 30, rgb: red},
		{bottom: 48, left: 30, right: 64, rgb: blue},
	}), capped(16*12))
	for x := range 16 {
		want := [4]uint8{255, 0, 0, 255}
		if x >= 8 {
			want = [4]uint8{0, 0, 255, 255}
		}
		if px := rgbaAt(c.pages[0], x, 5); px != want {
			t.Errorf("layers that meet, at %d,5: %v, want %v", x, px, want)
		}
	}
	// A line thinner than a reduced pixel is as faint as it is thin,
	// wherever it lies.
	for _, x := range []int{12, 13, 14, 15} {
		c = convert(t, layered(64, 48, []testLayer{{bottom: 48, left: x, right: x + 1, rgb: red}}), capped(16*12))
		for px := range 16 {
			want := [4]uint8{}
			if px == 3 {
				want = [4]uint8{255, 0, 0, 64}
			}
			if got := rgbaAt(c.pages[0], px, 5); got != want && (want[3] != 0 || got[3] != 0) {
				t.Errorf("line at %d, pixel %d: %v, want %v", x, px, got, want)
			}
		}
	}
	// A mask is placed as its layer is: here it hides the right half.
	data := layered(64, 48, []testLayer{{bottom: 48, right: 64, rgb: red}})
	f, err := readFile(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatal(err)
	}
	cp := &compositor{f: f, scale: scale{full: image.Pt(64, 48), to: image.Pt(16, 12)}, canvas: image.Rect(0, 0, 16, 12)}
	half := make([]uint8, 34*48)
	for i := range half {
		half[i] = 255
	}
	pix, at := cp.reduceMask(half, image.Rect(-4, 0, 30, 48), 0)
	if at != image.Rect(-1, 0, 8, 12) || len(pix) != 9*12 || pix[0] != 255 || pix[9*12-1] != 255 {
		t.Errorf("mask over %v, %d pixels", at, len(pix))
	}
	// A mask thinner than a pixel covers a part of one, its default value
	// the rest.
	pix, at = cp.reduceMask([]uint8{255, 255}, image.Rect(9, 0, 10, 2), 0)
	if at != image.Rect(2, 0, 3, 1) || len(pix) != 1 || pix[0] != 32 {
		t.Errorf("thin mask over %v: %v", at, pix)
	}
}

// TestBlendModesReduced: blend modes, masks, clipping and edges anywhere
// differ from the reduced composite of the canvas only where a reduced
// pixel stands for pixels that differ; on the whole the pictures agree.
func TestBlendModesReduced(t *testing.T) {
	for _, name := range []string{"layers.psd", "artboards.psd"} {
		data := withoutComposite(t, read(t, name))
		o := capped(60 * 45)
		o.NoArtboards = true
		c := convert(t, data, o)
		got := c.pages[0]
		b := got.Bounds()
		whole := *whole
		whole.NoArtboards = true
		small := imgconv.Resize(convert(t, data, &whole).pages[0], b.Dx(), b.Dy())
		far, sum := 0, 0
		for y := range b.Dy() {
			for x := range b.Dx() {
				r, g, bl, a := small.At(x, y).RGBA() // premultiplied
				px := got.NRGBAAt(x, y)
				for i, v := range []uint32{r, g, bl, a} {
					w := int(px.A)
					if i < 3 {
						w = int([]uint8{px.R, px.G, px.B}[i]) * int(px.A) / 255
					}
					d := max(w-int(v>>8), int(v>>8)-w)
					sum += d
					if d > 48 {
						far++
					}
				}
			}
		}
		// the samples far off are those of the edges of layers and
		// artboards, which the test documents are full of
		n := 4 * b.Dx() * b.Dy()
		if mean := float64(sum) / float64(n); mean > 3 || far > n/25 {
			t.Errorf("%s at %v: mean difference %.2f, %d of %d samples far off", name, b, mean, far, n)
		}
	}

	// The pages of artboards are cut from the reduced composite where the
	// artboards are placed, and have the size they are stored with.
	data := withoutComposite(t, read(t, "artboards.psd"))
	c := convert(t, data, capped(50*40))
	full := convert(t, data, whole)
	if len(c.pages) != 3 || c.res.Scaled != 3 {
		t.Fatalf("%d pages, %d scaled", len(c.pages), c.res.Scaled)
	}
	for i, p := range c.pages {
		fb := full.pages[i].Bounds()
		tw, th := capped(50*40).Images.FitSize(fb.Dx(), fb.Dy(), 0, 0)
		if b := p.Bounds(); b.Dx() != tw || b.Dy() != th {
			t.Errorf("artboard %d: %v, want %d × %d", i+1, b, tw, th)
		}
		if pg := c.r.Manifest.Views[0].Pages[i]; pg.W != float32(fb.Dx()) || pg.H != float32(fb.Dy()) {
			t.Errorf("artboard %d: page of %v × %v points", i+1, pg.W, pg.H)
		}
	}
	// the white background of artboard 1 and the transparent one of 3
	if px := rgbaAt(c.pages[0], 1, 1); px != [4]uint8{255, 255, 255, 255} {
		t.Errorf("artboard 1: %v", px)
	}
	if px := rgbaAt(c.pages[2], 1, 1); px[3] != 0 {
		t.Errorf("artboard 3: %v", px)
	}
}

// TestLimitsReduced: the limits of compositing are those of what is
// allocated, the reduced canvas and its buffers.
func TestLimitsReduced(t *testing.T) {
	limit(t, &maxCompositePixels, 20*15)
	data := layered(64, 48, []testLayer{{section: 3}, {section: 3}, {bottom: 48, right: 64, rgb: [3]uint8{255, 0, 0}}, {section: 1}, {section: 1}})
	// the canvas and two groups are alive at once: 3 × 16 × 12 pixels
	limit(t, &maxLivePixels, 3*16*12)
	limit(t, &maxCompositeWork, 3*16*12)
	c := convert(t, data, capped(16*12))
	if px := rgbaAt(c.pages[0], 5, 5); px != [4]uint8{255, 0, 0, 255} || len(c.res.Warnings) != 1 {
		t.Errorf("reduced: %v at 5,5, warnings %q", px, c.res.Warnings)
	}
	// Without a cap the canvas is composited as large as it is.
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), whole); err == nil || !strings.Contains(err.Error(), "64 × 48 pixels are too large to composite") {
		t.Errorf("not reduced: %v", err)
	}
	// With a cap that leaves it too large, the error says so.
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), capped(32*24)); err == nil || !strings.Contains(err.Error(), "reduced to 32 × 24 pixels too") {
		t.Errorf("reduced to 32 × 24: %v", err)
	}
	// A layer is decoded as large as it is stored, and has its limit there.
	data = layered(64, 48, []testLayer{{bottom: 48, right: 64, rgb: [3]uint8{0, 0, 255}}, {top: 4, left: 4, bottom: 8, right: 8, rgb: [3]uint8{255, 0, 0}}})
	limit(t, &maxLayerPixels, 16*12)
	c = convert(t, data, capped(16*12))
	if px := rgbaAt(c.pages[0], 1, 1); px != [4]uint8{255, 0, 0, 255} || rgbaAt(c.pages[0], 5, 5)[3] != 0 || !hasWarning(c.res.Warnings, "too large to draw") {
		t.Errorf("large layer: %v at 1,1, warnings %q", px, c.res.Warnings)
	}

	// What is allocated goes by the reduced canvas (the megabyte that
	// remains is the encoder's).
	limit(t, &maxCompositePixels, 1<<26)
	limit(t, &maxLivePixels, 1<<28)
	limit(t, &maxCompositeWork, 1<<32)
	data = layered(1024, 1024, []testLayer{{top: 100, left: 100, bottom: 108, right: 108, rgb: [3]uint8{255, 0, 0}}})
	var err error
	n := allocated(func() { _, err = Convert(bytes.NewReader(data), int64(len(data)), capped(16*16)) })
	if err != nil || n > 4<<20 {
		t.Errorf("canvas of 1024 × 1024 pixels stored as 16 × 16: %v, %d bytes allocated", err, n)
	}
}

// TestNotReduced: a document within the resolution cap is composited as
// before, pixel for pixel (the pages of the test documents, hashed when the
// layers were still composited at the size of the canvas in any case).
func TestNotReduced(t *testing.T) {
	for name, want := range map[string]string{
		"layers.psd":    "c082b8cf0188229ea1e56caacde27cb56be655fe064ee3dd03306bf50f951959",
		"artboards.psd": "5b7dfffeaf502558648292f68efa48edac3d0d4616cce182619f9357ea402608",
		"large.psb":     "fc74452f2206b8b78c99fe61678803f22132a70c3534d001116ff5827d68bdbf",
		"gray16.psd":    "e3e0c7a83342d012d61cd7ffc1cb42aab51c9d5f412f96b6387bd2572692d8fe",
	} {
		c := convert(t, withoutComposite(t, read(t, name)), nil)
		h := sha256.New()
		for _, p := range c.pages {
			fmt.Fprintf(h, "%v", p.Bounds())
			h.Write(p.Pix)
		}
		if got := fmt.Sprintf("%x", h.Sum(nil)); got != want || c.res.Scaled != 0 {
			t.Errorf("%s: pages %s, %d scaled", name, got, c.res.Scaled)
		}
	}
}

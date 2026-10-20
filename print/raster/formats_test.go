package raster

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"

	"github.com/shibukawa/bdf/raster/imagebdf"
)

// TestURFRoundTrip writes pages as Apple Raster and reads them back: the
// pixels are those the renderer draws, and the header has the sides, the
// quality and the resolution.
func TestURFRoundTrip(t *testing.T) {
	d, v := doc(t)
	r := imagebdf.New(d, &imagebdf.Options{BandRows: 20, ImageCacheBytes: -1})
	pl := imagebdf.Place{W: 144, H: 216, Scale: 2}
	for _, typ := range []Type{SGray8, SRGB8} {
		var buf bytes.Buffer
		w := NewURFWriter(&buf, len(v.Pages))
		mono := imagebdf.Mono(0)
		if typ == SGray8 {
			mono = imagebdf.Gray
		}
		for page := range v.Pages {
			if err := w.StartPage(Page{Width: pl.W, Height: pl.H, XDPI: 144, YDPI: 144, Type: typ, Duplex: true, Tumble: page == 1, Quality: 5,
				MediaPosition: 2}); err != nil {
				t.Fatal(err)
			}
			if err := r.PlaceRows(v, page, pl, mono, func(_ int, rows image.Image) error { return w.WriteRows(rows) }); err != nil {
				t.Fatal(err)
			}
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		b := buf.Bytes()
		if string(b[:8]) != URFSync || b[11] != 2 {
			t.Fatalf("%s: the stream begins %q", typ, b[:12])
		}
		pages, err := ReadURF(bytes.NewReader(b), 0)
		if err != nil {
			t.Fatalf("%s: %v", typ, err)
		}
		if len(pages) != 2 {
			t.Fatalf("%s: %d pages", typ, len(pages))
		}
		p := pages[1]
		if p.Type != typ || p.Width != 144 || p.Height != 216 || p.XDPI != 144 || !p.Duplex || !p.Tumble || pages[0].Tumble || p.Quality != 5 || p.MediaPosition != 2 {
			t.Fatalf("%s: header %+v", typ, p.Page)
		}
		want, err := r.PageMono(v, 1, 2, max(mono, imagebdf.Gray))
		if mono == 0 {
			want, err = r.Page(v, 1, 2)
		}
		if err != nil {
			t.Fatal(err)
		}
		if n, most := differ(want, p.Image); n > 0 {
			t.Errorf("%s: %d pixels differ (by %d at most)", typ, n, most)
		}
	}
}

func TestURFErrors(t *testing.T) {
	var buf bytes.Buffer
	w := NewURFWriter(&buf, 1)
	if err := w.StartPage(Page{Width: 4, Height: 2, XDPI: 300, YDPI: 300, Type: Black1}); err == nil {
		t.Error("black_1 in Apple Raster")
	}
	// no pages: the stream says none
	buf.Reset()
	if err := NewURFWriter(&buf, 0).Close(); err != nil || buf.String() != URFSync+"\x00\x00\x00\x00" {
		t.Errorf("an empty stream: %q %v", buf.String(), err)
	}
	// more pages than the stream says
	buf.Reset()
	w = NewURFWriter(&buf, 1)
	for i := 0; i < 2; i++ {
		w.StartPage(Page{Width: 4, Height: 1, XDPI: 300, Type: SGray8})
		w.WriteRows(image.NewGray(image.Rect(0, 0, 4, 1)))
	}
	w.Close()
	if _, err := ReadURF(bytes.NewReader(buf.Bytes()), 0); err == nil {
		t.Error("a stream of more pages than it says")
	}
	if _, err := ReadURF(bytes.NewReader([]byte("RaS2")), 0); err == nil {
		t.Error("PWG Raster read as Apple Raster")
	}
}

// readPCLm reads back what PCLmWriter writes: the pages with their
// strips put where the content of the page draws them.
func readPCLm(t *testing.T, b []byte) []image.Image {
	t.Helper()
	if !bytes.HasPrefix(b, []byte("%PDF-1.7\n%PCLm 1.0\n")) {
		t.Fatalf("the file begins %q", b[:min(len(b), 20)])
	}
	// every object is where the cross-reference table says
	m := regexp.MustCompile(`(?s)xref\n0 (\d+)\n0000000000 65535 f \n(.*?)trailer`).FindSubmatch(b)
	if m == nil {
		t.Fatal("no cross-reference table")
	}
	offs := regexp.MustCompile(`(\d{10}) 00000 n `).FindAllSubmatch(m[2], -1)
	objs := map[int][]byte{}
	for i, o := range offs {
		off, _ := strconv.Atoi(string(o[1]))
		head := fmt.Sprintf("%d 0 obj\n", i+1)
		if !bytes.HasPrefix(b[off:], []byte(head)) {
			t.Fatalf("object %d is not at %d", i+1, off)
		}
		end := bytes.Index(b[off:], []byte("endobj\n"))
		objs[i+1] = b[off+len(head) : off+end]
	}
	stream := func(obj []byte) (dict, data []byte) {
		i := bytes.Index(obj, []byte(">>\nstream\n"))
		dict = obj[:i]
		n, _ := strconv.Atoi(string(regexp.MustCompile(`/Length (\d+)`).FindSubmatch(dict)[1]))
		return dict, obj[i+len(">>\nstream\n") : i+len(">>\nstream\n")+n]
	}
	kids := regexp.MustCompile(`(\d+) 0 R`).FindAllSubmatch(regexp.MustCompile(`/Kids \[([^\]]*)\]`).FindSubmatch(objs[2])[1], -1)
	var pages []image.Image
	for _, k := range kids {
		n, _ := strconv.Atoi(string(k[1]))
		page := objs[n]
		box := regexp.MustCompile(`/MediaBox \[0 0 ([\d.]+) ([\d.]+)\]`).FindSubmatch(page)
		if box == nil {
			t.Fatalf("page %d has no media box", n)
		}
		cn, _ := strconv.Atoi(string(regexp.MustCompile(`/Contents (\d+) 0 R`).FindSubmatch(page)[1]))
		_, content := stream(objs[cn])
		xobjects := map[string]int{}
		for _, x := range regexp.MustCompile(`/(S\d+) (\d+) 0 R`).FindAllSubmatch(page, -1) {
			xobjects[string(x[1])], _ = strconv.Atoi(string(x[2]))
		}
		var img *image.RGBA
		for _, d := range regexp.MustCompile(`q (\d+) 0 0 (\d+) 0 (\d+) cm /(S\d+) Do Q`).FindAllSubmatch(content, -1) {
			w, _ := strconv.Atoi(string(d[1]))
			h, _ := strconv.Atoi(string(d[2]))
			y, _ := strconv.Atoi(string(d[3]))
			dict, data := stream(objs[xobjects[string(d[4])]])
			var pixels []byte
			gray := bytes.Contains(dict, []byte("/DeviceGray"))
			switch {
			case bytes.Contains(dict, []byte("/FlateDecode")):
				zr, err := zlib.NewReader(bytes.NewReader(data))
				if err != nil {
					t.Fatal(err)
				}
				pixels, _ = io.ReadAll(zr)
			case bytes.Contains(dict, []byte("/RunLengthDecode")):
				pixels = unRunLength(data)
			case bytes.Contains(dict, []byte("/DCTDecode")):
				ji, err := jpeg.Decode(bytes.NewReader(data))
				if err != nil {
					t.Fatal(err)
				}
				strip := image.NewRGBA(image.Rect(0, 0, w, h))
				draw.Draw(strip, strip.Rect, ji, ji.Bounds().Min, draw.Src)
				for i := 0; i < len(strip.Pix); i += 4 {
					if gray {
						pixels = append(pixels, strip.Pix[i])
					} else {
						pixels = append(pixels, strip.Pix[i:i+3]...)
					}
				}
			}
			if img == nil {
				// the first strip is at the top: its y and height give the page's
				img = image.NewRGBA(image.Rect(0, 0, w, y+h))
			}
			top := img.Rect.Dy() - y - h
			unit := 3
			if gray {
				unit = 1
			}
			if len(pixels) != w*h*unit {
				t.Fatalf("a strip of %d bytes for %d × %d pixels", len(pixels), w, h)
			}
			for row := 0; row < h; row++ {
				for x := 0; x < w; x++ {
					p := pixels[(row*w+x)*unit:]
					o := img.PixOffset(x, top+row)
					if gray {
						img.Pix[o], img.Pix[o+1], img.Pix[o+2] = p[0], p[0], p[0]
					} else {
						img.Pix[o], img.Pix[o+1], img.Pix[o+2] = p[0], p[1], p[2]
					}
					img.Pix[o+3] = 255
				}
			}
		}
		pages = append(pages, img)
	}
	return pages
}

func unRunLength(data []byte) []byte {
	var out []byte
	for i := 0; i < len(data); {
		n := int(data[i])
		switch {
		case n == 128:
			return out
		case n < 128:
			out = append(out, data[i+1:i+2+n]...)
			i += 2 + n
		default:
			out = append(out, bytes.Repeat(data[i+1:i+2], 257-n)...)
			i += 2
		}
	}
	return out
}

// TestPCLm writes pages as PCLm with each compression and reads them
// back: the strips make the pixels the renderer draws (JPEG nearly).
func TestPCLm(t *testing.T) {
	d, v := doc(t)
	r := imagebdf.New(d, &imagebdf.Options{BandRows: 20, ImageCacheBytes: -1})
	pl := imagebdf.Place{W: 144, H: 216, Scale: 2}
	want, err := r.Page(v, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	wantGray, err := r.PageMono(v, 1, 2, imagebdf.Gray)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		typ  Type
		comp string
	}{{SRGB8, "flate"}, {SGray8, "flate"}, {SRGB8, "rle"}, {SGray8, "rle"}, {SRGB8, "jpeg"}, {SGray8, "jpeg"}} {
		var buf bytes.Buffer
		w, err := NewPCLmWriter(&buf, PCLmOptions{StripHeight: 20, Compression: c.comp})
		if err != nil {
			t.Fatal(err)
		}
		mono := imagebdf.Mono(0)
		if c.typ == SGray8 {
			mono = imagebdf.Gray
		}
		for page := range v.Pages {
			if err := w.StartPage(Page{Width: pl.W, Height: pl.H, XDPI: 144, YDPI: 144, Type: c.typ}); err != nil {
				t.Fatal(err)
			}
			if err := r.PlaceRows(v, page, pl, mono, func(_ int, rows image.Image) error { return w.WriteRows(rows) }); err != nil {
				t.Fatal(err)
			}
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		pages := readPCLm(t, buf.Bytes())
		if len(pages) != 2 || pages[1].Bounds().Dx() != 144 || pages[1].Bounds().Dy() != 216 {
			t.Fatalf("%s %s: %d pages", c.typ, c.comp, len(pages))
		}
		if bytes.Count(buf.Bytes(), []byte("/Subtype /Image")) != 2*11 {
			t.Errorf("%s %s: %d strips, want 11 a page (20 rows, the last 16)", c.typ, c.comp, bytes.Count(buf.Bytes(), []byte("/Subtype /Image")))
		}
		ref := image.Image(want)
		if c.typ == SGray8 {
			ref = wantGray
		}
		n, most := differ(ref, pages[1])
		if c.comp == "jpeg" {
			if most > 80 {
				t.Errorf("%s jpeg: %d pixels differ, by %d at most", c.typ, n, most)
			}
		} else if n > 0 {
			t.Errorf("%s %s: %d pixels differ (by %d at most)", c.typ, c.comp, n, most)
		}
		if c.typ == SRGB8 && c.comp == "flate" {
			ghostscript(t, buf.Bytes(), want)
		}
	}
	if _, err := NewPCLmWriter(io.Discard, PCLmOptions{Compression: "lzw"}); err == nil {
		t.Error("an unknown compression")
	}
	w, _ := NewPCLmWriter(io.Discard, PCLmOptions{})
	if err := w.StartPage(Page{Width: 4, Height: 2, XDPI: 300, Type: Black1}); err == nil {
		t.Error("black_1 in PCLm")
	}
}

// ghostscript draws a PCLm file with Ghostscript, when it is installed:
// it reads it as a PDF, and draws the page.
func ghostscript(t *testing.T, file []byte, want image.Image) {
	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Log("no Ghostscript: the PCLm file is not drawn by it")
		return
	}
	dir := t.TempDir()
	in, out := filepath.Join(dir, "in.pdf"), filepath.Join(dir, "page%d.png")
	os.WriteFile(in, file, 0o644)
	if b, err := exec.Command(gs, "-q", "-dNOPAUSE", "-dBATCH", "-dSAFER", "-sDEVICE=png16m", "-r144", "-o", out, in).CombinedOutput(); err != nil {
		t.Fatalf("Ghostscript: %v\n%s", err, b)
	}
	f, err := os.Open(filepath.Join(dir, "page2.png"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if n, most := differ(want, img); n < 0 || n > img.Bounds().Dx()*img.Bounds().Dy()/50 {
		t.Errorf("Ghostscript's page: %d pixels differ (by %d at most)", n, most)
	}
}

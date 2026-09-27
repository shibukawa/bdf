package thumbnail

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
	"github.com/shibukawa/bdf/raster"
	"golang.org/x/image/webp"
)

func open(t *testing.T, path string) *bdf.Document {
	t.Helper()
	r, err := bdf.OpenSingleFile(path)
	if err != nil {
		t.Fatal(err)
	}
	d, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// testRaster draws with the repository's fonts only, the same everywhere.
var testRaster = raster.Options{FontDirs: []string{"../converter/pptx/testdata/fonts", "../fixture/testdata/fonts"}, NoSystemFonts: true}

func TestLayouts(t *testing.T) {
	for _, c := range []struct {
		file string
		mode Mode
		want Mode
		w, h int
	}{
		{"docx/basic", Auto, Crop, 256, 256},       // a text document: the top of page 1
		{"pptx/basic", Auto, Fit, 256, 192},        // a 4:3 slide, whole
		{"pdf/chrome-slides", Auto, Fit, 256, 144}, // landscape PDF pages are slides
		{"pdf/chrome-doc", Auto, Crop, 256, 256},   // portrait ones documents
		{"xlsx/basic", Auto, Crop, 256, 256},       // a sheet from A1
		{"csv/basic", Fit, Fit, 256, 256},          // a sheet has no page to fit
		{"markdown/basic", Auto, Crop, 256, 256},   // a scroll view from its top
		{"epub/basic", Auto, Fit, 180, 256},        // a book's cover, whole
		{"visio/flow", Auto, Fit, 198, 256},
		{"docx/basic", Fit, Fit, 181, 256}, // asked for the whole page
		{"pptx/basic", Crop, Crop, 256, 256},
	} {
		res, err := Make(open(t, "../testdata/"+c.file+".bdf"), &Options{Mode: c.mode, Raster: testRaster})
		if err != nil {
			t.Fatalf("%s: %v", c.file, err)
		}
		if b := res.Image.Bounds(); res.Mode != c.want || b.Dx() != c.w || b.Dy() != c.h {
			t.Errorf("%s (%s): %s %d×%d, want %s %d×%d", c.file, c.mode, res.Mode, b.Dx(), b.Dy(), c.want, c.w, c.h)
		}
		if blank(res.Image) {
			t.Errorf("%s: blank thumbnail", c.file)
		}
	}
}

// blank reports whether an image is one colour.
func blank(img *image.RGBA) bool {
	for i := 4; i < len(img.Pix); i++ {
		if img.Pix[i] != img.Pix[i%4] {
			return false
		}
	}
	return true
}

// TestCropIsTopLeft checks that a cropped thumbnail is the top-left square
// of the page: the same as drawing that square of the page directly.
func TestCropIsTopLeft(t *testing.T) {
	d := open(t, "../testdata/docx/basic.bdf")
	res, err := Make(d, &Options{Size: 128, Raster: testRaster})
	if err != nil {
		t.Fatal(err)
	}
	v := d.Views[0]
	side := v.Pages[0].W // portrait: the square is as wide as the page
	want, err := raster.New(d, &testRaster).Region(v, 0, bdf.Rect{W: side, H: side}, 128, 128)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res.Image.Pix, want.Pix) {
		t.Error("the thumbnail is not the top-left square of page 1")
	}
}

func TestSizeAndView(t *testing.T) {
	d := open(t, "../testdata/drawio/multipage.bdf")
	res, err := Make(d, &Options{Size: 100, View: "details", Raster: testRaster})
	if err != nil {
		t.Fatal(err)
	}
	if b := res.Image.Bounds(); max(b.Dx(), b.Dy()) != 100 {
		t.Errorf("size %v", b)
	}
	if _, err := Make(d, &Options{View: "nope"}); err == nil {
		t.Error("an unknown view was drawn")
	}
}

func TestEncode(t *testing.T) {
	res, err := Make(open(t, "../testdata/pptx/basic.bdf"), &Options{Size: 64, Raster: testRaster})
	if err != nil {
		t.Fatal(err)
	}
	decoders := map[string]func([]byte) (image.Image, error){
		PNG:  func(b []byte) (image.Image, error) { return png.Decode(bytes.NewReader(b)) },
		JPEG: func(b []byte) (image.Image, error) { return jpeg.Decode(bytes.NewReader(b)) },
		WebP: func(b []byte) (image.Image, error) { return webp.Decode(bytes.NewReader(b)) },
	}
	for _, f := range []string{PNG, JPEG, WebP} {
		if f == WebP && !imgconv.Available() {
			continue
		}
		var b bytes.Buffer
		if err := Encode(&b, res.Image, f); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		img, err := decoders[f](b.Bytes())
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if img.Bounds().Size() != res.Image.Bounds().Size() {
			t.Errorf("%s: %v", f, img.Bounds())
		}
	}
	for name, want := range map[string]string{"a.PNG": PNG, "b.jpg": JPEG, "c.jpeg": JPEG, "d.webp": WebP, "e.gif": ""} {
		if got := FormatOf(name); got != want {
			t.Errorf("FormatOf(%q) = %q", name, got)
		}
	}
}

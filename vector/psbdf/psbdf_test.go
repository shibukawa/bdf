package psbdf

import (
	"bytes"
	"compress/zlib"
	"encoding/ascii85"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// testFonts are the fonts of the tests: text is laid out with them alone, so
// that the files do not depend on the machine.
var testFonts = []string{"../../converter/pptx/testdata/fonts"}

func testOptions() *Options {
	return &Options{FontDirs: testFonts, NoSystemFonts: true}
}

func load(t testing.TB, name string) *bdf.Document {
	t.Helper()
	r, err := bdf.OpenSingleFile(filepath.Join("../../testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	doc, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

func write(t testing.TB, doc *bdf.Document, opts *Options) (string, *Result) {
	t.Helper()
	var buf bytes.Buffer
	res, err := Write(&buf, doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	return buf.String(), res
}

func one(t testing.TB, w, h float32, draw func(doc *bdf.Document, o *bdf.Object)) *bdf.Document {
	t.Helper()
	doc := bdf.NewDocument()
	o := bdf.NewObject()
	draw(doc, o)
	o.SetBBox(0, 0, w, h)
	hash, _ := doc.AddObject(o)
	v := doc.NewView("pages", bdf.ViewFixed, "")
	p := v.AddPage(w, h)
	p.Layers = append(p.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	return doc
}

// ghostscript draws the pages of a file at 72 dpi with Ghostscript, and
// fails on what it says; the test is skipped without it.
func ghostscript(t testing.TB, file string) []*image.RGBA {
	t.Helper()
	gs, err := exec.LookPath("gs")
	if err != nil {
		t.Skip("Ghostscript (gs) is not installed")
	}
	dir := t.TempDir()
	cmd := exec.Command(gs, "-q", "-dSAFER", "-dBATCH", "-dNOPAUSE", "-sDEVICE=png16m", "-r72", "-dTextAlphaBits=4", "-dGraphicsAlphaBits=4",
		"-sOutputFile="+filepath.Join(dir, "p%d.png"), "-")
	cmd.Stdin = strings.NewReader(file)
	if out, err := cmd.CombinedOutput(); err != nil || len(out) > 0 {
		t.Fatalf("Ghostscript: %v\n%s", err, out)
	}
	var pages []*image.RGBA
	for n := 1; ; n++ {
		f, err := os.Open(filepath.Join(dir, fmt.Sprintf("p%d.png", n)))
		if err != nil {
			return pages
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatal(err)
		}
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Rect, img, img.Bounds().Min, draw.Src)
		pages = append(pages, rgba)
	}
}

func render(t testing.TB, doc *bdf.Document, v *bdf.View, page int) *image.RGBA {
	t.Helper()
	img, err := imagebdf.New(doc, &imagebdf.Options{FontDirs: testFonts, NoSystemFonts: true}).Page(v, page, 1)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// distance is the mean difference of the channels of two images, of 255,
// where both are (Ghostscript rounds the size of a page down, the image
// renderer up).
func distance(a, b *image.RGBA) float64 {
	w, h := min(a.Rect.Dx(), b.Rect.Dx()), min(a.Rect.Dy(), b.Rect.Dy())
	if w == 0 || h == 0 || max(a.Rect.Dx(), b.Rect.Dx())-w > 1 || max(a.Rect.Dy(), b.Rect.Dy())-h > 1 {
		return 255
	}
	sum := 0.0
	for y := 0; y < h; y++ {
		for x := 0; x < 4*w; x++ {
			sum += math.Abs(float64(a.Pix[y*a.Stride+x]) - float64(b.Pix[y*b.Stride+x]))
		}
	}
	return sum / float64(4*w*h)
}

// save writes the images of a page that differ, where $PSBDF_OUT says.
func save(name string, imgs ...*image.RGBA) {
	dir := os.Getenv("PSBDF_OUT")
	if dir == "" {
		return
	}
	for i, img := range imgs {
		var b bytes.Buffer
		png.Encode(&b, img)
		os.WriteFile(filepath.Join(dir, name+"-"+string(rune('a'+i))+".png"), b.Bytes(), 0o644)
	}
}

// fontsOf returns the definitions of the Type 3 fonts of a file, inflated.
func fontsOf(t testing.TB, file string) string {
	t.Helper()
	i := strings.Index(file, "bdfrun\n")
	if i < 0 {
		return ""
	}
	data := file[i+len("bdfrun\n"):]
	data = data[:strings.Index(data, "~>")]
	data = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' {
			return -1
		}
		return r
	}, data)
	dec, err := io.ReadAll(ascii85.NewDecoder(strings.NewReader(data)))
	if err != nil {
		t.Fatal(err)
	}
	zr, err := zlib.NewReader(bytes.NewReader(dec))
	if err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// TestGhostscript writes documents as PostScript and draws the files with
// Ghostscript: their pages are drawn as those of the documents.
func TestGhostscript(t *testing.T) {
	if !woff2.Available() {
		t.Skip("the fonts of the test documents are WOFF2, which this build does not read")
	}
	for _, c := range []struct {
		name     string
		pages    int
		fonts    int
		pictures int
		// most is the most a page may differ from that of the document
		// (Ghostscript and the image renderer smooth the edges of text
		// differently)
		most float64
	}{
		{"form/fixed.bdf", 1, 2, 0, 4},
		{"form/pages.bdf", 3, 2, 0, 6},
		{"pptx/basic.bdf", 5, 2, 1, 2},
		{"pptx/features.bdf", 5, 2, 17, 2},
		{"docx/basic.bdf", 3, 2, 0, 5.5},
		{"docx/vertical.bdf", 2, 2, 0, 4},
		{"xlsx/print-pages.bdf", 17, 2, 0, 10.5},
		{"pdf/chrome-masks.bdf", 1, 0, 1, 2},
		{"pdf/cairo-cff.bdf", 1, 2, 0, 2},
		{"pdf/reportlab-master.bdf", 3, 3, 0, 1.5},
		{"drawio/showcase.bdf", 2, 2, 2, 4},
		{"image/photo.bdf", 1, 0, 0, 1},
		{"epub/vertical.bdf", 6, 3, 0, 5},
	} {
		t.Run(c.name, func(t *testing.T) {
			doc := load(t, c.name)
			file, res := write(t, doc, testOptions())
			if res.Pages != c.pages || res.Fonts != c.fonts || res.Pictures != c.pictures {
				t.Errorf("%d pages, %d fonts and %d pictures, want %d, %d and %d", res.Pages, res.Fonts, res.Pictures, c.pages, c.fonts, c.pictures)
			}
			for _, w := range res.Warnings {
				t.Errorf("warning: %s", w)
			}
			for _, line := range strings.Split(file, "\n") {
				if len(line) > 255 {
					t.Fatalf("a line of %d bytes: %.80q", len(line), line)
				}
			}
			if n := strings.Count(file, "\n%%Page: "); n != c.pages || !strings.Contains(file, fmt.Sprintf("\n%%%%Pages: %d\n", c.pages)) {
				t.Errorf("%d pages in the comments, want %d", n, c.pages)
			}
			got := ghostscript(t, file)
			if len(got) != c.pages {
				t.Fatalf("Ghostscript drew %d pages, want %d", len(got), c.pages)
			}
			n := 0
			for _, v := range doc.Views {
				if v.Kind != bdf.ViewFixed && v.Kind != bdf.ViewFlow {
					continue
				}
				for i := range v.Pages {
					if n++; n > 3 {
						break
					}
					a := render(t, doc, v, i)
					if d := distance(a, got[n-1]); d > c.most {
						t.Errorf("page %d differs by %.2f from that of the document (at most %.2f)", n, d, c.most)
						save(strings.ReplaceAll(c.name, "/", "_")+"-"+string(rune('0'+n)), a, got[n-1])
					}
				}
			}
		})
	}
}

// TestStructure checks the comments of the file, and what it asks of the
// device.
func TestStructure(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("v", bdf.ViewFixed, "")
	for i := 0; i < 3; i++ {
		o := bdf.NewObject()
		o.FillColor(bdf.RGB(255, 0, 0)).FillRect(0, 0, float32(10*(i+1)), 10)
		h, _ := d.AddObject(o)
		v.AddPage(float32(100+10*i), 200, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	d.Meta.DC.Title = bdf.DCValues{"題名 (1)"}
	opts := testOptions()
	opts.Pages = []int{2, -1, 0, 7}
	opts.Copies, opts.Collate, opts.Duplex, opts.Tumble = 3, true, true, true
	opts.MediaType, opts.MediaPosition = "stationery", 2
	file, res := write(t, d, opts)
	if res.Pages != 3 || res.Fonts != 0 || res.Pictures != 0 {
		t.Errorf("%+v", res)
	}
	for _, want := range []string{
		"%!PS-Adobe-3.0\n", "\n%%Title: (\\351\\241\\214\\345\\220\\215 \\(1\\))\n", "\n%%LanguageLevel: 3\n", "\n%%Pages: 3\n",
		"\n%%BoundingBox: 0 0 120 200\n", "\n%%DocumentData: Clean7Bit\n", "\n%%EndComments\n%%BeginProlog\n",
		"%%BeginSetup\n%%BeginFeature: *PageSize\nmark { << /PageSize [120 200] >> setpagedevice } stopped cleartomark\n%%EndFeature\n",
		"<< /NumCopies 3 >>", "<< /Collate true >>", "<< /Duplex true /Tumble true >>", "<< /MediaType (stationery) >>", "<< /MediaPosition 2 >>",
		"bdf begin\n%%EndSetup\n%%Page: 1 1\nsave\n1 0 0 -1 0 200 cm\n",
		"1 0 0 rg\n", "0 0 30 10 re\nf\n",
		"%%Page: 2 2\nsave\nrestore showpage\n", // the blank page
		"%%Page: 3 3\n%%BeginPageSetup\n%%BeginFeature: *PageSize\nmark { << /PageSize [100 200] >> setpagedevice } stopped cleartomark\n%%EndFeature\n%%EndPageSetup\nsave\n",
		"restore showpage\n%%Trailer\nend\n%%EOF\n",
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	// no line of data looks like a comment
	for _, line := range strings.Split(file, "\n") {
		if len(line) > 255 {
			t.Errorf("a line of %d bytes", len(line))
		}
	}
	// (three copies of each)
	if got := ghostscript(t, file); len(got) != 9 || got[0].Rect.Dx() != 120 || got[8].Rect.Dx() != 100 {
		t.Errorf("Ghostscript drew %d pages", len(got))
	}
	opts = testOptions()
	opts.Pages = []int{5}
	if _, err := Write(io.Discard, d, opts); !errors.Is(err, ErrNoPages) {
		t.Errorf("no page selected: %v", err)
	}
	if _, err := Write(io.Discard, d, &Options{Views: []string{"none"}}); err == nil || !strings.Contains(err.Error(), `no view "none"`) {
		t.Errorf("no view: %v", err)
	}
	// placed on a sheet of another size
	opts = testOptions()
	opts.Place = func(w, h float64) (float64, float64, [6]float64) { return 300, 400, [6]float64{2, 0, 0, -2, 0, 400} }
	opts.Gray = true
	file, _ = write(t, d, opts)
	for _, want := range []string{"<< /PageSize [300 400] >>", "save\n2 0 0 -2 0 400 cm\n", "0.299 g\n"} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	if strings.Contains(file, " rg\n") || strings.Contains(file, "%%BeginPageSetup") {
		t.Error("a gray file has colors, or the sheets differ")
	}
}

// TestPictures draws what PostScript cannot as pictures of the parts of
// the page where it is drawn: the rest stays vectors.
func TestPictures(t *testing.T) {
	doc := one(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.FillColor(bdf.RGB(0, 0, 200)).FillRect(0, 0, 200, 100)
		o.FillColor(bdf.RGB(250, 250, 0)).FillRect(150, 50, 30, 30) // under a picture: left out by the clip
		o.FillColor(bdf.RGBA(255, 0, 0, 128)).FillRect(140, 40, 20, 20)
		o.GroupBegin(0.5, bdf.BlendSourceOver, 10, 10, 30, 30).FillColor(bdf.RGB(0, 255, 0)).FillRect(10, 10, 30, 30).GroupEnd()
		o.GroupBegin(1, bdf.BlendSourceOver, 60, 10, 30, 30).FillColor(bdf.RGB(0, 255, 255)).FillRect(60, 10, 30, 30).GroupEnd()
		o.Save().Shadow(bdf.RGBA(0, 0, 0, 128), 2, 3, 3).FillColor(bdf.RGB(255, 255, 255)).FillRect(100, 70, 10, 10).Restore()
	})
	file, res := write(t, doc, testOptions())
	if res.Pictures != 3 || len(res.Warnings) != 0 {
		t.Errorf("%d pictures, warnings %q", res.Pictures, res.Warnings)
	}
	for _, want := range []string{
		"q -32768 -32768 65536 65536 re 139 39 22 22 re 9 9 32 32 re 98 68 20 20 re W* n\n", // the holes of the pictures, with the shadow
		"0 0 200 100 re\nf\n", "150 50 30 30 re\nf\n",
		"q\n0 1 1 rg\n60 10 30 30 re\nf\nQ\n", // an opaque group is vectors
		"Q\nq 22 0 0 22 139 39 cm\n<< /ImageType 1 /Width 92 /Height 92 /BitsPerComponent 8 /Decode [0 1 0 1 0 1] /ImageMatrix [92 0 0 92 0 0] >> /DeviceRGB /FlateDecode bdfimage\n",
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	for _, not := range []string{"140 40 20 20 re", "10 10 30 30 re", "100 70 10 10 re"} {
		if strings.Contains(file, not) {
			t.Errorf("the file draws %q, which is a picture", not)
		}
	}
	got := ghostscript(t, file)
	if d := distance(render(t, doc, doc.Views[0], 0), got[0]); d > 1 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		save("pictures", render(t, doc, doc.Views[0], 0), got[0])
	}

	// a page that is mostly pictures is one
	doc = one(t, 100, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.FillColor(bdf.RGB(0, 0, 200)).FillRect(0, 0, 100, 100)
		o.Alpha(0.5).FillColor(bdf.RGB(255, 0, 0)).FillRect(10, 10, 80, 80)
	})
	file, res = write(t, doc, &Options{RasterDPI: 72, NoSystemFonts: true})
	if res.Pictures != 1 || strings.Contains(file, " re\nf\n") || !strings.Contains(file, "q 100 0 0 100 0 0 cm\n<< /ImageType 1 /Width 100 /Height 100 ") {
		t.Errorf("%d pictures:\n%s", res.Pictures, file[strings.Index(file, "%%Page:"):][:300])
	}
}

// TestPaints writes gradients as smooth shadings, and what PostScript
// cannot draw of them as pictures.
func TestPaints(t *testing.T) {
	doc := one(t, 300, 100, func(doc *bdf.Document, o *bdf.Object) {
		lin := o.AddPaint(bdf.LinearGradient(0, 0, 100, 0, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 0.5, Color: bdf.RGB(0, 255, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 255)}))
		rad := o.AddPaint(bdf.RadialGradient(150, 50, 0, 150, 50, 40, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 255, 255)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 0)}))
		fade := o.AddPaint(bdf.LinearGradient(200, 0, 300, 0, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGBA(255, 0, 0, 0)}))
		con := o.AddPaint(bdf.ConicGradient(0, 250, 80, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 255)}))
		one := o.AddPaint(bdf.LinearGradient(0, 0, 1, 0, bdf.Stop{Offset: 0.3, Color: bdf.RGB(10, 20, 30)}))
		none := o.AddPaint(bdf.LinearGradient(0, 0, 1, 0))
		o.FillPaint(lin).FillRect(0, 0, 100, 50)
		o.StrokePaint(lin).Line(4, bdf.CapButt, bdf.JoinMiter, 10).StrokeRect(10, 60, 80, 30)
		o.FillPaint(rad).FillPath(o.AddPath((&bdf.Path{}).Circle(150, 50, 40)), bdf.NonZero)
		o.FillPaint(fade).FillRect(200, 0, 100, 40)
		o.FillPaint(con).FillRect(240, 70, 20, 20)
		o.FillPaint(one).FillRect(0, 0, 9, 9).FillPaint(none).FillRect(0, 0, 8, 8)
	})
	file, res := write(t, doc, testOptions())
	if res.Pictures != 2 || len(res.Warnings) != 0 {
		t.Errorf("%d pictures, warnings %q", res.Pictures, res.Warnings)
	}
	for _, want := range []string{
		"q\n0 0 100 50 re\nW n\n<< /ShadingType 2 /ColorSpace /DeviceRGB /Coords [0 0 100 0] /Function << /FunctionType 3 /Domain [0 1] /Functions [\n<< /FunctionType 2 /Domain [0 1] /C0 [1 0 0] /C1 [0 1 0] /N 1 >>\n",
		"/N 1 >>\n] /Bounds [0.5 ] /Encode [0 1 0 1 ] >> /Extend [true true] >> shfill\nQ\n",
		"q\n10 60 80 30 re\nstrokepath W n\n<< /ShadingType 2",
		"<< /ShadingType 3 /ColorSpace /DeviceRGB /Coords [150 50 0 150 50 40] ",
		"0.039216 0.078431 0.117647 rg\n0 0 9 9 re\nf\n", // a gradient of one stop is its color
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	for _, not := range []string{"0 0 8 8 re", "200 0 100 40 re", "240 70 20 20 re"} {
		if strings.Contains(file, not) {
			t.Errorf("the file has %q", not)
		}
	}
	// (Ghostscript does not smooth the edges of the clip of a shading)
	got := ghostscript(t, file)
	if d := distance(render(t, doc, doc.Views[0], 0), got[0]); d > 2.5 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		save("paints", render(t, doc, doc.Views[0], 0), got[0])
	}
	opts := testOptions()
	opts.Gray = true
	file, _ = write(t, doc, opts)
	if !strings.Contains(file, "/ShadingType 2 /ColorSpace /DeviceGray") || !strings.Contains(file, "/C0 [0.299] /C1 [0.587]") || strings.Contains(file, "/DeviceRGB") {
		t.Error("the shadings of a gray file are not gray")
	}
}

// TestText writes text with Type 3 fonts of the outlines of its glyphs.
func TestText(t *testing.T) {
	font, err := os.ReadFile(filepath.Join(testFonts[0], "MPLUS1p-Regular-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	doc := one(t, 300, 200, func(doc *bdf.Document, o *bdf.Object) {
		f := o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 400, bdf.StyleNormal))
		bold := o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 700, bdf.StyleItalic))
		lin := o.AddPaint(bdf.LinearGradient(0, 0, 100, 0, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 255)}))
		o.Font(f, 12).FillText("Hello 日本", 20, 40, 0)
		o.TextStyle(bdf.AlignCenter, bdf.BaselineTop, bdf.DirLTR, 1).FillText("wide", 150, 60, 80)
		o.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirLTR, 0)
		o.Font(bold, 20).FillText("Bold", 20, 100, 0)
		o.Font(f, 12).StrokeColor(bdf.RGB(255, 0, 0)).StrokeText("Line", 20, 130, 0)
		o.FillPaint(lin).FillText("Paint", 20, 160, 0)
		o.FillColor(bdf.RGBA(0, 0, 0, 100)).FillText("Faint", 200, 160, 0)
		o.FillText("", 0, 0, 0)
	})
	file, res := write(t, doc, testOptions())
	if res.Fonts != 1 || res.Pictures != 1 || len(res.Warnings) != 0 {
		t.Errorf("%d fonts, %d pictures, warnings %q", res.Fonts, res.Pictures, res.Warnings)
	}
	for _, want := range []string{
		"/BDF0 [12 0 0 -12 0 0] Tf 20 40 m <0001020203040506>\n[8.52 6.624 3.684 3.684 6.84 3.48 12 12] xshow n\n",
		"/BDF0 [31.003746 0 0 -12 0 0] Tf 110 69.2473 m <", // centred at the top, as wide as asked
		"f\nq 0.7485 w 1 j [] 0 d\n",                       // synthesized bold: the fill and a line
		"S\nQ\n",
		"h\nS\n", // STROKE_TEXT: the outlines
		"h\nW n\n<< /ShadingType 2",
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	fonts := fontsOf(t, file)
	for _, want := range []string{
		"/BDF0 <<\n/FontType 3 /FontMatrix [1 1000 div 0 0 1 1000 div 0 0] /FontBBox [",
		"/Encoding [\n/g0 /g1 ", "/.notdef ]\n/CharProcs <<\n/.notdef {0 0 0 0 0 0 setcachedevice}\n/g0 {",
		" setcachedevice\n", " m\n", "f}\n", ">> definefont pop\n",
	} {
		if !strings.Contains(fonts, want) {
			t.Errorf("the fonts lack %q", want)
		}
	}
	// (the others are outlines, and the faint text is a picture)
	if strings.Count(file, "xshow") != 2 {
		t.Errorf("%d strings of glyphs", strings.Count(file, "xshow"))
	}
	got := ghostscript(t, file)
	if d := distance(render(t, doc, doc.Views[0], 0), got[0]); d > 2 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		save("text", render(t, doc, doc.Views[0], 0), got[0])
	}

	// more glyphs than a font holds
	doc = one(t, 600, 300, func(doc *bdf.Document, o *bdf.Object) {
		f := o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 400, bdf.StyleNormal))
		o.Font(f, 8)
		sf, err := sfnt.Parse(font)
		if err != nil {
			t.Fatal(err)
		}
		var line []rune
		for r := rune(0x20); r < 0x10000; r++ {
			if g := sf.Cmap[uint32(r)]; g != 0 {
				line = append(line, r)
			}
		}
		for y := 0; y*60 < len(line); y++ {
			o.FillText(string(line[y*60:min(len(line), y*60+60)]), 10, float32(10+y*9), 0)
		}
	})
	file, res = write(t, doc, testOptions())
	if res.Fonts < 2 || !strings.Contains(file, "/BDF1 [8 0 0 -8 0 0] Tf") || strings.Count(fontsOf(t, file), "definefont") != res.Fonts {
		t.Errorf("%d fonts", res.Fonts)
	}
}

// TestImages writes images as their data: JPEG files as they are, others
// as their pixels, deflated.
func TestImages(t *testing.T) {
	photo := image.NewRGBA(image.Rect(0, 0, 16, 8))
	gray := image.NewNRGBA(image.Rect(0, 0, 4, 2))
	alpha := image.NewNRGBA(image.Rect(0, 0, 4, 2))
	for i := 0; i < 128; i++ {
		photo.Set(i%16, i/16, color.RGBA{uint8(i * 2), 100, 200, 255})
	}
	for i := 0; i < 8; i++ {
		gray.Set(i%4, i/4, color.NRGBA{uint8(i * 30), uint8(i * 30), uint8(i * 30), 255})
		alpha.Set(i%4, i/4, color.NRGBA{uint8(i * 30), 100, 200, uint8(255 - i*20)})
	}
	var jb, gb, ab bytes.Buffer
	jpeg.Encode(&jb, photo, &jpeg.Options{Quality: 90})
	png.Encode(&gb, gray)
	png.Encode(&ab, alpha)
	doc := one(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		j := o.AddImage(doc.AddImage(jb.Bytes()))
		g := o.AddImage(doc.AddImage(gb.Bytes()))
		a := o.AddImage(doc.AddImage(ab.Bytes()))
		o.Image(j, 10, 10, 64, 32).ImageSub(g, 1, 0, 2, 2, 100, 10, 40, 40)
		o.Smoothing(false, 0).Image(g, 10, 60, 40, 20).Image(a, 150, 60, 40, 20)
	})
	file, res := write(t, doc, testOptions())
	if res.Pictures != 1 || len(res.Warnings) != 0 {
		t.Errorf("%d pictures, warnings %q", res.Pictures, res.Warnings)
	}
	for _, want := range []string{
		"q 64 0 0 32 10 10 cm\n<< /ImageType 1 /Width 16 /Height 8 /BitsPerComponent 8 /Decode [0 1 0 1 0 1] /ImageMatrix [16 0 0 8 0 0] /Interpolate true >> /DeviceRGB /DCTDecode bdfimage\n",
		"q 100 10 40 40 re W n\nq 80 0 0 40 80 10 cm\n<< /ImageType 1 /Width 4 /Height 2 /BitsPerComponent 8 /Decode [0 1] /ImageMatrix [4 0 0 2 0 0] /Interpolate true >> /DeviceGray /FlateDecode bdfimage\n",
		"q 40 0 0 20 10 60 cm\n<< /ImageType 1 /Width 4 /Height 2 /BitsPerComponent 8 /Decode [0 1] /ImageMatrix [4 0 0 2 0 0] >> /DeviceGray",
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	got := ghostscript(t, file)
	if d := distance(render(t, doc, doc.Views[0], 0), got[0]); d > 2 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		save("images", render(t, doc, doc.Views[0], 0), got[0])
	}
}

// TestDeterministic checks that a document is always written as the same
// file.
func TestDeterministic(t *testing.T) {
	if !woff2.Available() {
		t.Skip("the fonts of the test documents are WOFF2, which this build does not read")
	}
	a, _ := write(t, load(t, "pptx/features.bdf"), testOptions())
	b, _ := write(t, load(t, "pptx/features.bdf"), testOptions())
	if a != b {
		t.Error("a document is written as different files")
	}
}

// TestSkipRoles leaves out the layers of roles, in the vectors and in the
// pictures.
func TestSkipRoles(t *testing.T) {
	doc := one(t, 100, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.FillColor(bdf.RGBA(255, 0, 0, 128)).FillRect(10, 10, 20, 20)
	})
	bg := bdf.NewObject()
	bg.FillColor(bdf.RGB(0, 0, 255)).FillRect(0, 0, 100, 100).SetBBox(0, 0, 100, 100)
	h, _ := doc.AddObject(bg)
	p := doc.Views[0].Pages[0]
	p.Layers = append([]bdf.Layer{{Role: bdf.RoleBackground, Obj: h}}, p.Layers...)
	opts := testOptions()
	opts.SkipRoles = []string{bdf.RoleBackground}
	file, res := write(t, doc, opts)
	if res.Pictures != 1 || strings.Contains(file, "0 0 100 100 re") {
		t.Errorf("%d pictures", res.Pictures)
	}
	got := ghostscript(t, file)
	if c := got[0].RGBAAt(20, 20); c.B > 200 || c.R < 200 {
		t.Errorf("the picture has the background: %v", c)
	}
}

// TestLimits checks that a document cannot ask for a file of any size, and
// that what it leaves open is closed.
func TestLimits(t *testing.T) {
	// objects that draw each other twice, forty deep
	doc := bdf.NewDocument()
	leaf := bdf.NewObject()
	leaf.FillRect(0, 0, 1, 1).SetBBox(0, 0, 1, 1)
	h, box := doc.AddObject(leaf)
	for i := 0; i < 40; i++ {
		o := bdf.NewObject()
		r := o.AddObject(h, box)
		o.Use(r).UseAt(r, 1, 1).SetBBox(0, 0, 2, 2)
		h, box = doc.AddObject(o)
	}
	v := doc.NewView("pages", bdf.ViewFixed, "")
	p := v.AddPage(10, 10)
	p.Layers = append(p.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	var buf bytes.Buffer
	res, err := Write(&buf, doc, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() > 256<<20 || len(res.Warnings) == 0 || !strings.Contains(res.Warnings[0], "drawn again too many times") {
		t.Errorf("%d bytes, warnings %q", buf.Len(), res.Warnings)
	}
	// groups without end, states without end, and ends without beginnings,
	// in objects that end what they did not begin
	inner := bdf.NewObject()
	inner.GroupEnd().MaskEnd().Restore().FillRect(0, 0, 2, 2).GroupBegin(1, 0, 0, 0, 5, 5).Save().SetBBox(0, 0, 5, 5)
	ih, ibox := bdf.NewDocument().AddObject(inner)
	doc = one(t, 10, 10, func(doc *bdf.Document, o *bdf.Object) {
		ih, ibox = doc.AddObject(inner)
		r := o.AddObject(ih, ibox)
		o.GroupEnd().MaskEnd().Restore().Restore()
		o.GroupBegin(1, 0, 0, 0, 10, 10).Use(r).MaskBegin(bdf.MaskAlpha, 0, nil).Use(r).FillRect(0, 0, 3, 3)
		for i := 0; i < 100; i++ {
			o.GroupBegin(1, 0, 0, 0, 10, 10).MaskBegin(bdf.MaskAlpha, 0, nil).GroupBegin(1, 0, 0, 0, 10, 10)
		}
		for i := 0; i < 2000; i++ {
			o.Save()
		}
		o.FillRect(0, 0, 1, 1)
	})
	_ = ih
	file, res := write(t, doc, testOptions())
	if q, Q := balance(file); q != Q {
		t.Errorf("%d q and %d Q", q, Q)
	}
	if len(res.Warnings) == 0 {
		t.Error("no warning of what is left out")
	}
	ghostscript(t, file)
}

// balance counts the q and the Q of the pages of a file, outside the data
// of images.
func balance(file string) (q, Q int) {
	data := false
	for _, line := range strings.Split(file[strings.Index(file, "%%Page: "):], "\n") {
		switch {
		case strings.HasSuffix(line, " bdfimage"):
			data = true
		case line == "~>":
			data = false
		}
		if data {
			continue
		}
		for _, tok := range strings.Fields(line) {
			switch tok {
			case "q":
				q++
			case "Q":
				Q++
			}
		}
	}
	return q, Q
}

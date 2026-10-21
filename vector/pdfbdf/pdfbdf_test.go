package pdfbdf

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/vector/internal/vectortest"
)

func testOptions() *Options {
	return &Options{FontDirs: vectortest.Fonts, NoSystemFonts: true}
}

func write(t testing.TB, doc *bdf.Document, opts *Options) ([]byte, *Result) {
	t.Helper()
	var buf bytes.Buffer
	res, err := Write(&buf, doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	return buf.Bytes(), res
}

// render draws a page of a view of a document with the image renderer.
func render(t testing.TB, doc *bdf.Document, view, page int) *image.RGBA {
	t.Helper()
	return vectortest.Render(t, doc, doc.Views[view], page)
}

// back reads a PDF file as a document again.
func back(t testing.TB, file []byte) *bdf.Document {
	t.Helper()
	res, err := pdf.Convert(bytes.NewReader(file), nil)
	if err != nil {
		t.Fatalf("the file cannot be read as a PDF: %v", err)
	}
	return res.Doc
}

func plain(t testing.TB, doc *bdf.Document) string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	for _, v := range st.Views {
		for _, p := range v.Pages {
			b.WriteString(p.Text + "\n")
		}
	}
	return b.String()
}

// TestRoundTrip writes documents as PDF files and reads the files again:
// their pages are drawn as those of the documents, and have their text.
//
// The reader of PDF files (converter/pdf) places gradients that are not
// horizontal wrongly: the pages with such gradients may differ more, and
// the gradients of files are checked by their numbers (TestPaintsAndGroups).
func TestRoundTrip(t *testing.T) {
	if !woff2.Available() {
		t.Skip("the fonts of the test documents are WOFF2, which this build does not read")
	}
	for _, c := range []struct {
		name  string
		pages int
		fonts int
		// most is the most a page may differ from that of the document
		most float64
		text []string
	}{
		{"form/fixed.bdf", 1, 2, 2, []string{"Certificate", "ACME さま", "13,200.00"}},
		{"form/pages.bdf", 3, 2, 2, []string{"Delivery D-1", "Part 24"}},
		{"pptx/basic.bdf", 5, 2, 6, []string{"BDF from PowerPoint", "日本語の文章は"}},
		{"pptx/features.bdf", 5, 2, 3, nil},
		{"docx/basic.bdf", 3, 2, 3, []string{"BDF Word Converter", "文書の行は、変換器が折り返します。"}},
		{"docx/vertical.bdf", 2, 2, 3, nil},
		{"xlsx/print-pages.bdf", 17, 2, 3, []string{"Items by line", "Item 029"}},
		{"pdf/chrome-masks.bdf", 1, 3, 3, []string{"Soft masks fade this heading out"}},
		{"pdf/cairo-cff.bdf", 1, 2, 3, nil},
		{"pdf/reportlab-master.bdf", 3, 3, 3, nil},
		{"drawio/showcase.bdf", 2, 2, 3, nil},
		{"image/photo.bdf", 1, 0, 1, nil},
		{"epub/vertical.bdf", 6, 3, 3, nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			doc := vectortest.Load(t, c.name)
			file, res := write(t, doc, testOptions())
			if res.Pages != c.pages || res.Fonts != c.fonts {
				t.Errorf("%d pages and %d fonts, want %d and %d", res.Pages, res.Fonts, c.pages, c.fonts)
			}
			for _, w := range res.Warnings {
				t.Errorf("warning: %s", w)
			}
			if !bytes.HasPrefix(file, []byte("%PDF-1.7\n")) || !bytes.HasSuffix(file, []byte("%%EOF\n")) {
				t.Fatalf("not a PDF file: %q … %q", file[:min(len(file), 20)], file[max(len(file)-20, 0):])
			}
			// every object that has a number is in the file
			if n := bytes.Count(file, []byte(" 65535 f \n")); n != 1 {
				t.Errorf("%d objects of the table are not written", n-1)
			}
			got := back(t, file)
			if n := len(got.Views[0].Pages); n != c.pages {
				t.Fatalf("the file has %d pages, want %d", n, c.pages)
			}
			// the pages of the file are those of the views of pages, in their order
			n := 0
			for vi, v := range doc.Views {
				if v.Kind != bdf.ViewFixed && v.Kind != bdf.ViewFlow {
					continue
				}
				for i := range v.Pages {
					if n++; n > 3 {
						break
					}
					a, b := vectortest.Render(t, doc, doc.Views[vi], i), vectortest.Render(t, got, got.Views[0], n-1)
					if d := vectortest.Distance(a, b); d > c.most {
						t.Errorf("page %d differs by %.2f from that of the document (at most %.2f)", n, d, c.most)
						vectortest.Save("PDFBDF_OUT", strings.ReplaceAll(c.name, "/", "_")+"-"+string(rune('0'+n)), a, b)
					}
				}
			}
			text := plain(t, got)
			for _, want := range c.text {
				if !strings.Contains(text, want) {
					t.Errorf("the text of the file lacks %q:\n%s", want, text)
				}
			}
		})
	}
}

// TestDeterministic checks that a document is always written as the same
// file.
func TestDeterministic(t *testing.T) {
	doc := vectortest.Load(t, "pptx/basic.bdf")
	a, _ := write(t, doc, testOptions())
	b, _ := write(t, vectortest.Load(t, "pptx/basic.bdf"), testOptions())
	if !bytes.Equal(a, b) {
		t.Error("the same document is written as two files")
	}
}

// content writes a document without compression and returns the file as
// text.
func content(t testing.TB, doc *bdf.Document) (string, *Result) {
	t.Helper()
	o := testOptions()
	o.NoCompression = true
	file, res := write(t, doc, o)
	return string(file), res
}

func TestOperators(t *testing.T) {
	doc := vectortest.One(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.Save().Translate(10, 20).FillColor(bdf.RGBA(255, 0, 0, 128)).FillRect(0, 0, 30, 40).Restore()
		o.StrokeColor(bdf.RGB(0, 0, 255)).Line(2, bdf.CapRound, bdf.JoinBevel, 4).Dash([]float32{3, 1, 2}, 1)
		p := o.AddPath((&bdf.Path{}).MoveTo(0, 0).LineTo(10, 0).QuadTo(20, 0, 20, 10).Close().Rect(1, 2, 3, 4))
		o.StrokePath(p).FillPath(p, bdf.EvenOdd).ClipPath(p, bdf.NonZero).ClipRect(0, 0, 5, 5)
		o.Alpha(0.5).Blend(bdf.BlendMultiply).FillRect(0, 0, 1, 1)
		o.Blend(bdf.BlendDestinationOut).FillRect(0, 0, 1, 1)
		o.Shadow(bdf.RGB(0, 0, 0), 4, 1, 1).FillRect(0, 0, 1, 1)
		o.Link(10, 10, 50, 20, "https://example.com/a b").Link(0, 0, 5, 5, "#page=1").Link(0, 0, 5, 5, "javascript:alert(1)")
	})
	doc.Meta.DC.Title = bdf.DCValues{"題名 <1>"}
	doc.Meta.DC.Creator = bdf.DCValues{"A", "B"}
	doc.Meta.DC.Created = bdf.DCValues{"2026-10-14T09:30:00+09:00"}
	doc.Meta.DC.Language = bdf.DCValues{"ja"}
	file, res := content(t, doc)
	for _, want := range []string{
		"/MediaBox [0 0 200 100]",
		"1 0 0 -1 0 100 cm\n",           // the y axis of a document goes down
		"q\n1 0 0 1 10 20 cm\n1 0 0 rg", // SAVE, TRANSLATE, FILL_COLOR
		"/ca 0.501961",                  // the alpha of a color
		"0 0 30 40 re\nf\nQ\n",
		"2 w\n1 J\n2 j\n4 M\n[3 1 2 3 1 2] 1 d\n0 0 1 RG\n",
		"0 0 m\n10 0 l\n16.6667 0 20 3.3333 20 10 c\nh\n1 2 m\n4 2 l\n4 6 l\n1 6 l\nh\nS\n",
		"h\nf*\n", "h\nW n\n0 0 5 5 re W n\n",
		"/BM /Multiply", "/ca 0.5 ",
		"/Subtype /Link /Rect [10 70 60 90]", "/S /URI /URI (https://example.com/a%20b)", "/Dest [",
		"/Title <FEFF984C540D0020003C0031003E>", "/Author (A, B)", "/CreationDate (D:20261014093000+09'00')", "/Lang (ja)",
		"<dc:title><rdf:Alt><rdf:li xml:lang=\"x-default\">題名 &lt;1&gt;</rdf:li></rdf:Alt></dc:title>", "<xmp:CreateDate>2026-10-14T09:30:00+09:00</xmp:CreateDate>",
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	if strings.Contains(file, "javascript") {
		t.Error("a link that is no link of a document is written")
	}
	want := []string{"the composite operation destination-out is drawn as source-over"}
	if strings.Join(res.Warnings, "\n") != strings.Join(want, "\n") {
		t.Errorf("warnings: %q, want %q", res.Warnings, want)
	}
	back(t, []byte(file))
}

func TestPaintsAndGroups(t *testing.T) {
	var buf bytes.Buffer
	img := image.NewNRGBA(image.Rect(0, 0, 4, 2))
	for i := 0; i < 8; i++ {
		img.Set(i%4, i/4, color.NRGBA{uint8(i * 30), 100, 200, uint8(255 - i*20)})
	}
	png.Encode(&buf, img)
	doc := vectortest.One(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		lin := o.AddPaint(bdf.LinearGradient(0, 0, 100, 0, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 0.5, Color: bdf.RGB(0, 255, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 255)}))
		rad := o.AddPaint(bdf.RadialGradient(50, 50, 0, 50, 50, 40, bdf.Stop{Offset: 0, Color: bdf.RGBA(255, 255, 255, 255)}, bdf.Stop{Offset: 1, Color: bdf.RGBA(0, 0, 0, 0)}))
		con := o.AddPaint(bdf.ConicGradient(0, 50, 50, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 255)}))
		one := o.AddPaint(bdf.LinearGradient(0, 0, 1, 0, bdf.Stop{Offset: 0.3, Color: bdf.RGBA(10, 20, 30, 128)}))
		none := o.AddPaint(bdf.LinearGradient(0, 0, 1, 0))
		pic := o.AddImage(doc.AddImage(buf.Bytes()))
		tile := o.AddPaint(bdf.Pattern(pic, bdf.RepeatBoth))
		once := o.AddPaint(bdf.Pattern(pic, bdf.NoRepeat))
		o.FillPaint(lin).FillRect(0, 0, 100, 50)
		o.Save().Translate(100, 0).FillPaint(lin).FillRect(0, 0, 100, 50).Restore()
		o.FillPaint(rad).FillRect(0, 0, 100, 100)
		o.FillPaint(con).FillRect(0, 0, 100, 100)
		o.FillPaint(one).FillRect(0, 0, 9, 9).FillPaint(none).FillRect(0, 0, 8, 8)
		o.FillPaint(tile).FillRect(0, 0, 50, 50).StrokePaint(lin).StrokeRect(0, 0, 7, 7)
		o.FillPaint(once).FillRect(0, 0, 50, 50)
		o.Image(pic, 10, 10, 40, 20).ImageSub(pic, 1, 0, 2, 2, 60, 10, 20, 20)
		o.GroupBegin(0.5, bdf.BlendScreen, 10, 10, 80, 60).FillColor(bdf.RGB(0, 128, 0)).FillRect(10, 10, 80, 60)
		o.MaskBegin(bdf.MaskLuminosity, bdf.RGB(0, 0, 0), nil).FillColor(bdf.RGB(255, 255, 255)).FillRect(20, 20, 30, 30).MaskEnd().GroupEnd()
		o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 0, 10).FillRect(0, 0, 6, 6).GroupEnd() // nothing of it shows
	})
	file, res := content(t, doc)
	for _, want := range []string{
		"/ShadingType 2 /ColorSpace /DeviceRGB /Coords [0 0 100 0]", "/FunctionType 3", "/Bounds [0.5 ]",
		"/PatternType 2", "/Matrix [1 0 0 -1 0 100]", "/Matrix [1 0 0 -1 100 100]", // a pattern is in the space of the page
		"/ShadingType 3 /ColorSpace /DeviceRGB /Coords [50 50 0 50 50 40]",
		"/ShadingType 3 /ColorSpace /DeviceGray", "/SMask << /Type /Mask /S /Luminosity /G", // the alphas of a gradient
		"1 0 0 -1 0 100 cm /GS", " gs 1 0 0 -1 0 100 cm\n", "1 0 0 -1 0 100 cm /Sh", // set in the space of the page
		"/ShadingType 4", "/BitsPerCoordinate 32",
		"0.0392 0.0784 0.1176 rg", "0 0 9 9 re\nf\n", // a gradient of one stop is its color
		"/PatternType 1 /PaintType 1 /TilingType 1 /BBox [0 0 4 2] /XStep 4 /YStep 2",
		"/Pattern CS /P",
		"0 0 50 50 re\nW n\n4 0 0 -2 0 2 cm /Im", // a pattern that does not repeat: the image in the path
		"/Subtype /Image /Width 4 /Height 2 /ColorSpace /DeviceRGB /BitsPerComponent 8 /SMask",
		"/Subtype /Image /Width 4 /Height 2 /ColorSpace /DeviceGray /BitsPerComponent 8 /Interpolate true",
		"q 40 0 0 -20 10 30 cm /Im", "q 60 10 20 20 re W n 40 0 0 -20 50 30 cm /Im",
		// the form of a group is in the space of the page, with its transform
		"/Subtype /Form /BBox [9 29 91 91]", "/Group << /Type /Group /S /Transparency /I true /CS /DeviceRGB >>", "1 0 0 -1 0 100 cm 10 10 80 60 re W n\n",
		"/ca 0.5 ", "/BM /Screen", "/SMask << /Type /Mask /S /Luminosity /G", "/BC [0 0 0]",
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	for _, not := range []string{"0 0 8 8 re", "0 0 6 6 re"} {
		if strings.Contains(file, not) {
			t.Errorf("the file has %q, which draws nothing", not)
		}
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	back(t, []byte(file))
}

// TestGroupsRoundTrip draws groups, soft masks, images and patterns from
// the file again: as the document draws them.
func TestGroupsRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for i := 0; i < 64; i++ {
		img.Set(i%8, i/8, color.NRGBA{uint8(i * 4), uint8(255 - i*3), 200, uint8(128 + i*2)})
	}
	png.Encode(&buf, img)
	inv := make([]byte, 256)
	for i := range inv {
		inv[i] = byte(255 - i)
	}
	doc := vectortest.One(t, 200, 150, func(doc *bdf.Document, o *bdf.Object) {
		pic := o.AddImage(doc.AddImage(buf.Bytes()))
		tile := o.AddPaint(bdf.Pattern(pic, bdf.RepeatBoth))
		o.FillColor(bdf.RGB(250, 240, 200)).FillRect(0, 0, 200, 150)
		// under the transforms that the pages of converted PDF files have
		o.Save().Transform(1, 0, 0, -1, 0, 150).Transform(0.5, 0, 0, -0.5, 0, 150)
		o.GroupBegin(0.6, bdf.BlendMultiply, 20, 20, 160, 100).FillColor(bdf.RGB(200, 30, 30)).FillRect(20, 20, 160, 100)
		o.FillPaint(tile).FillRect(40, 40, 60, 60)
		o.MaskBegin(bdf.MaskAlpha, 0, nil).FillColor(bdf.RGB(0, 0, 0)).FillRect(20, 20, 100, 100).MaskEnd().GroupEnd()
		o.GroupBegin(1, bdf.BlendSourceOver, 200, 20, 160, 100).FillColor(bdf.RGB(30, 30, 200)).FillRect(200, 20, 160, 100)
		o.MaskBegin(bdf.MaskLuminosity, bdf.RGB(0, 0, 0), inv).FillColor(bdf.RGB(255, 255, 255))
		o.FillPath(o.AddPath((&bdf.Path{}).Circle(280, 70, 40)), bdf.NonZero).MaskEnd().GroupEnd()
		o.Restore()
		o.Image(pic, 20, 90, 40, 40).ImageSub(pic, 2, 2, 4, 4, 80, 90, 40, 40)
		o.Save().Translate(140, 90).Scale(2, 2).FillPaint(tile).FillRect(0, 0, 20, 20).Restore()
	})
	file, res := write(t, doc, testOptions())
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	a, b := render(t, doc, 0, 0), render(t, back(t, file), 0, 0)
	if d := vectortest.Distance(a, b); d > 2 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		vectortest.Save("PDFBDF_OUT", "groups", a, b)
	}
}

func TestText(t *testing.T) {
	font, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Regular-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	doc := vectortest.One(t, 300, 200, func(doc *bdf.Document, o *bdf.Object) {
		f := o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 400, bdf.StyleNormal))
		bold := o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 700, bdf.StyleItalic))
		o.Font(f, 12).FillText("Hello 日本", 20, 40, 0)
		o.TextStyle(bdf.AlignCenter, bdf.BaselineTop, bdf.DirLTR, 1).FillText("wide", 150, 60, 80)
		o.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirLTR, 0)
		o.Font(bold, 20).FillText("Bold", 20, 100, 0)
		o.Font(f, 12).StrokeColor(bdf.RGB(255, 0, 0)).StrokeText("Line", 20, 130, 0)
		o.Mark(bdf.MarkAltText, "fi").FillText("\ufb01", 20, 160, 0)
		o.FillText("", 0, 0, 0)
	})
	file, res := content(t, doc)
	if res.Fonts != 1 || len(res.Warnings) != 0 {
		t.Errorf("%d fonts, warnings %q", res.Fonts, res.Warnings)
	}
	for _, want := range []string{
		"/Subtype /Type0", "/Encoding /Identity-H", "/Subtype /CIDFontType2", "/CIDToGIDMap /Identity", "/FontFile2", "+MPLUS1p-Regular",
		"BT /F", " 12 Tf 1 0 0 -1 20 40 Tm <", "> Tj ET",
		"1 Tc ",                // letter spacing
		" 0 0.25 -1 20 100 Tm", // synthesized italic
		"2 Tr ",                // synthesized bold: the fill and a line
		"1 Tr ",                // STROKE_TEXT
		"/Span << /ActualText (fi) >> BDC", "EMC",
		"beginbfchar", "<65E5>", // 日 in the characters of the codes
	} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	got := back(t, []byte(file))
	if text := plain(t, got); !strings.Contains(text, "Hello 日本") || !strings.Contains(text, "wide") || !strings.Contains(text, "Bold") {
		t.Errorf("the text of the file: %q", text)
	}
	a, b := render(t, doc, 0, 0), vectortest.Render(t, got, got.Views[0], 0)
	if d := vectortest.Distance(a, b); d > 2 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		vectortest.Save("PDFBDF_OUT", "text", a, b)
	}
}

// TestShadows writes the shadows of what is drawn as pictures of them: the
// file is drawn as the document is.
func TestShadows(t *testing.T) {
	font, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Bold-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	doc := vectortest.One(t, 300, 200, func(doc *bdf.Document, o *bdf.Object) {
		o.FillColor(bdf.RGB(245, 245, 245)).FillRect(0, 0, 300, 200)
		// a blurred shadow, and one without a blur; under a scale and a turn,
		// which a shadow does not follow
		o.Save().Shadow(bdf.RGBA(0, 0, 0, 128), 8, 4, 6).FillColor(bdf.RGB(255, 255, 255)).FillRect(20, 20, 80, 50).Restore()
		o.Save().Translate(200, 40).Transform(0, 1, -1, 0, 0, 0).Scale(2, 2).Shadow(bdf.RGB(200, 0, 0), 0, 3, 3)
		o.FillColor(bdf.RGB(0, 120, 200)).FillPath(o.AddPath((&bdf.Path{}).Circle(0, 0, 15)), bdf.NonZero).Restore()
		// of a line, of text and of an image with transparent pixels
		o.Save().Shadow(bdf.RGBA(0, 0, 0, 200), 4, 2, 2).StrokeColor(bdf.RGB(0, 150, 0)).Line(6, bdf.CapRound, bdf.JoinRound, 10)
		o.StrokePath(o.AddPath((&bdf.Path{}).MoveTo(20, 110).LineTo(120, 130).LineTo(60, 170))).Restore()
		o.Save().Shadow(bdf.RGBA(0, 0, 80, 255), 3, 2, 2).Font(o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 700, bdf.StyleNormal)), 28)
		o.FillColor(bdf.RGB(255, 200, 0)).FillText("Shadow", 150, 150, 0).Restore()
		// nothing casts a shadow outside the state that has one
		o.FillColor(bdf.RGB(90, 90, 90)).FillRect(250, 170, 30, 20)
	})
	file, res := content(t, doc)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	// four pictures, each under what casts it: the first 4 and 6 units from the rectangle
	if n := strings.Count(file, "/SMask "); n != 4 {
		t.Errorf("%d pictures of shadows, want 4", n)
	}
	for _, want := range []string{"/ca 0.501961", "1 0 0 1 0 0 cm 104 0 0 -74 12 88 cm /Im", "20 20 80 50 re\nf\n"} {
		if !strings.Contains(file, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	if i, j := strings.Index(file, "74 12 88 cm /Im"), strings.Index(file, "20 20 80 50 re"); i < 0 || j < i {
		t.Error("the shadow is not under what casts it")
	}
	a, b := render(t, doc, 0, 0), render(t, back(t, []byte(file)), 0, 0)
	if d := vectortest.Distance(a, b); d > 1.5 {
		t.Errorf("the page differs by %.2f from that of the document", d)
		vectortest.Save("PDFBDF_OUT", "shadows", a, b)
	}
	// without the shadows the page is another
	plain := vectortest.One(t, 300, 200, func(doc *bdf.Document, o *bdf.Object) {
		o.FillColor(bdf.RGB(245, 245, 245)).FillRect(0, 0, 300, 200)
		o.FillColor(bdf.RGB(255, 255, 255)).FillRect(20, 20, 80, 50)
	})
	if d := vectortest.Distance(vectortest.Render(t, plain, plain.Views[0], 0), b); d < 1.5 {
		t.Errorf("the page of the file differs by %.2f from one without shadows", d)
	}
}

// TestGroupShadow writes the shadow of the state around a group: one
// picture of what the group draws as a whole (the shapes of a diagram that
// are drawn as groups to cast one shadow each), under the group.
func TestGroupShadow(t *testing.T) {
	font, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Bold-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	doc := vectortest.One(t, 300, 200, func(doc *bdf.Document, o *bdf.Object) {
		// an object that the group uses
		star := bdf.NewObject()
		star.FillColor(bdf.RGB(250, 200, 0)).FillPath(star.AddPath((&bdf.Path{}).MoveTo(0, -20).LineTo(18, 16).LineTo(-22, -6).LineTo(22, -6).LineTo(-18, 16).Close()), bdf.NonZero)
		star.SetBBox(-25, -25, 50, 50)
		h, box := doc.AddObject(star)
		o.FillColor(bdf.RGB(245, 245, 245)).FillRect(0, 0, 300, 200)
		o.Save().Translate(20, 20).Shadow(bdf.RGBA(0, 0, 0, 160), 6, 5, 7)
		o.FillColor(bdf.RGB(120, 170, 230)).Font(o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 700, bdf.StyleNormal)), 24)
		o.GroupBegin(0.9, bdf.BlendSourceOver, 0, 0, 220, 150)
		// shapes over each other, a group in the group, text in the font of
		// the state around the group, and the other object
		o.FillRect(10, 10, 120, 60).FillColor(bdf.RGB(230, 120, 120))
		o.FillPath(o.AddPath((&bdf.Path{}).Circle(120, 70, 35)), bdf.NonZero)
		o.GroupBegin(0.5, bdf.BlendSourceOver, 0, 80, 100, 60).FillColor(bdf.RGB(60, 160, 90)).FillRect(20, 90, 60, 40).GroupEnd()
		o.FillColor(bdf.RGB(40, 40, 40)).FillText("Group", 20, 50, 0)
		o.UseAt(o.AddObject(h, box), 180, 110)
		o.GroupEnd().Restore()
		// a group in a state without a shadow casts none
		o.GroupBegin(1, bdf.BlendSourceOver, 240, 150, 40, 40).FillColor(bdf.RGB(90, 90, 90)).FillRect(250, 160, 20, 20).GroupEnd()
	})
	file, res := content(t, doc)
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	// one picture, under the form of the group
	if n := strings.Count(file, "/SMask "); n != 1 {
		t.Errorf("%d pictures of shadows, want 1", n)
	}
	i := strings.Index(file, " cm /Im")
	if i < 0 {
		t.Fatal("no shadow in the page")
	}
	if rest := file[i:]; !strings.Contains(rest[:strings.Index(rest, "endstream")], " Do\nQ\nq\n/GS6 gs\n1 0 0 -1 -20 180 cm\n/Fm7 Do\n") {
		t.Errorf("the group is not drawn over its shadow: %q", rest[:min(len(rest), 120)])
	}
	if !strings.Contains(file, "/ca 0.627451") {
		t.Error("the shadow is not drawn with the alpha of its color")
	}
	a, b := render(t, doc, 0, 0), render(t, back(t, []byte(file)), 0, 0)
	if d := vectortest.Distance(a, b); d > 1.5 {
		t.Errorf("the page differs by %.2f from that of the document", d)
	}
	vectortest.Save("PDFBDF_OUT", "group-shadow", a, b)
	px := func(img *image.RGBA, x, y int) [3]uint8 {
		k := img.PixOffset(x, y)
		return [3]uint8(img.Pix[k : k+3])
	}
	// (under the rectangle of the group, where nothing but its shadow is)
	if c := px(b, 100, 96); c[0] > 200 || c[0] < 60 {
		t.Errorf("no shadow under the group: %v", c)
	}
	if c := px(b, 290, 100); c != [3]uint8{245, 245, 245} {
		t.Errorf("a shadow away from the group: %v", c)
	}
}

// TestRestrictedFont draws the text of a font of the system that may not be
// embedded as the outlines of its glyphs, under the text they stand for.
func TestRestrictedFont(t *testing.T) {
	data, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Regular-subset.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := sfnt.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	os2 := append([]byte(nil), f.Tables["OS/2"]...)
	os2[8], os2[9] = 0, 2 // fsType: restricted license embedding
	f.Tables["OS/2"] = os2
	fsys := fstest.MapFS{"restricted.ttf": {Data: sfnt.Build(f.Tables, false)}}
	doc := vectortest.One(t, 200, 100, func(doc *bdf.Document, o *bdf.Object) {
		o.Font(o.AddFont(bdf.SystemFont("M PLUS 1p", 400, bdf.StyleNormal)), 20).FillText("Hello", 20, 50, 0).StrokeText("Line", 20, 80, 0)
	})
	o := &Options{FontFS: fsys, NoSystemFonts: true, NoCompression: true}
	file, res := write(t, doc, o)
	if res.Fonts != 0 || strings.Contains(string(file), "/FontFile") || strings.Contains(string(file), " Tj") {
		t.Errorf("%d fonts are embedded", res.Fonts)
	}
	for _, want := range []string{"/Span << /ActualText (Hello) >> BDC", " c\n", "h\nf\nEMC", "/Span << /ActualText (Line) >> BDC", "h\nS\nEMC"} {
		if !strings.Contains(string(file), want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	a, b := render(t, doc, 0, 0), render(t, back(t, file), 0, 0)
	// (the renderer of the test has no font for the text: only the file has its outlines)
	if d := vectortest.Distance(a, b); d == 0 {
		t.Error("the page of the file has no text")
	}
	// with the rights to embed it, it is embedded
	o.IgnoreFSType = true
	file, res = write(t, doc, o)
	if res.Fonts != 1 || !strings.Contains(string(file), "/FontFile2") {
		t.Errorf("%d fonts are embedded with IgnoreFSType", res.Fonts)
	}
}

func TestViews(t *testing.T) {
	// the pages of a document with pages; its other views when it has none
	for _, c := range []struct {
		name  string
		opts  Options
		pages int
		box   string
	}{
		{"docx/basic.bdf", Options{}, 3, "/MediaBox [0 0 595.3 841.9]"},
		{"docx/basic.bdf", Options{Views: []string{"scroll"}}, 1, ""},
		{"html/article.bdf", Options{}, 1, ""},
		{"xlsx/basic.bdf", Options{}, 2, ""},
		{"xlsx/print.bdf", Options{}, 5, ""},
		{"csv/basic.bdf", Options{}, 1, ""},
	} {
		doc := vectortest.Load(t, c.name)
		o := c.opts
		o.FontDirs, o.NoSystemFonts, o.NoCompression = vectortest.Fonts, true, true
		file, res := write(t, doc, &o)
		if res.Pages != c.pages || c.box != "" && !strings.Contains(string(file), c.box) {
			t.Errorf("%s: %d pages, want %d (with %s)", c.name, res.Pages, c.pages, c.box)
		}
		if got := back(t, file); len(got.Views[0].Pages) != c.pages {
			t.Errorf("%s: the file has %d pages", c.name, len(got.Views[0].Pages))
		}
	}
	if _, err := Write(&bytes.Buffer{}, vectortest.Load(t, "docx/basic.bdf"), &Options{Views: []string{"none"}}); err == nil || !strings.Contains(err.Error(), `no view "none"`) {
		t.Errorf("a view the document does not have: %v", err)
	}
	if _, err := Write(&bytes.Buffer{}, bdf.NewDocument(), nil); err != ErrNoPages {
		t.Errorf("a document without views: %v", err)
	}
	// a page too large for the units of a PDF has a larger unit
	doc := vectortest.One(t, 30000, 100, func(doc *bdf.Document, o *bdf.Object) { o.FillRect(0, 0, 30000, 100) })
	file, _ := content(t, doc)
	if !strings.Contains(file, "/MediaBox [0 0 10000 33.3333]") || !strings.Contains(file, "/UserUnit 3") || !strings.Contains(file, "0.33333334 0 0 -0.33333334 0 33.333332 cm") {
		t.Error("a page of 30000 units is not written with a unit of 3")
	}
	// layers that pages share are written once, and layers are left out by their roles
	slides := vectortest.Load(t, "pptx/basic.bdf")
	file, _ = content(t, slides)
	if n := strings.Count(file, "/Subtype /Form /BBox [-32768 -32768 32768 32768]"); n == 0 || n > 4 {
		t.Errorf("%d forms of shared layers", n)
	}
	if !strings.Contains(file, "q 1 0 0 -1 0 540 cm /Fm") {
		t.Error("the form of a shared layer is not drawn in the space of the page")
	}
	o := testOptions()
	o.SkipRoles = []string{bdf.RoleMaster, bdf.RoleBackground}
	less, _ := write(t, slides, o)
	all, _ := write(t, slides, testOptions())
	if len(less) >= len(all) {
		t.Errorf("without masters and backgrounds the file has %d bytes, with them %d", len(less), len(all))
	}
}

// TestLimits checks that a document cannot ask for a file of any size.
func TestLimits(t *testing.T) {
	// objects that draw each other twice, forty deep
	doc := vectortest.DeepUses()
	var buf bytes.Buffer
	res, err := Write(&buf, doc, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() > 256<<20 || len(res.Warnings) == 0 || !strings.Contains(res.Warnings[0], "drawn again too many times") {
		t.Errorf("%d bytes, warnings %q", buf.Len(), res.Warnings)
	}
	// groups without end, states without end, and ends without beginnings
	doc = vectortest.One(t, 10, 10, func(doc *bdf.Document, o *bdf.Object) {
		o.GroupEnd().MaskEnd().Restore().Restore()
		for i := 0; i < 100; i++ {
			o.GroupBegin(1, 0, 0, 0, 10, 10).MaskBegin(bdf.MaskAlpha, 0, nil)
		}
		for i := 0; i < 2000; i++ {
			o.Save()
		}
		o.FillRect(0, 0, 1, 1)
	})
	file, res := content(t, doc)
	if q, Q := strings.Count(file, "\nq\n"), strings.Count(file, "\nQ\n"); q != Q {
		t.Errorf("%d q and %d Q", q, Q)
	}
	if len(res.Warnings) == 0 {
		t.Error("no warning of what is left out")
	}
	back(t, []byte(file))
}

// TestPages writes a selection of the pages of a document, in the order
// asked, and leaves out numbers past them.
func TestPages(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("v", bdf.ViewFixed, "")
	for i := 0; i < 3; i++ {
		o := bdf.NewObject()
		o.FillColor(0x000000ff).FillRect(0, 0, float32(10*(i+1)), 10)
		h, _ := d.AddObject(o)
		v.AddPage(float32(100+10*i), 200, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	opts := testOptions()
	opts.Pages = []int{2, 0, 7}
	file, res := write(t, d, opts)
	if res.Pages != 2 {
		t.Fatalf("%d pages", res.Pages)
	}
	got := back(t, file).Views[0].Pages
	if len(got) != 2 || got[0].W != 120 || got[1].W != 100 {
		t.Fatalf("pages of %v", []float32{got[0].W, got[1].W})
	}
	opts.Pages = []int{5}
	if _, err := Write(io.Discard, d, opts); !errors.Is(err, ErrNoPages) {
		t.Fatalf("no page selected: %v", err)
	}
}

package pdfbdf

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/vector/internal/vectortest"
)

// FuzzWrite writes a page of an object of any bytes: whatever the
// instructions of the object, a file is written, and the writing does not
// fail of itself.
func FuzzWrite(f *testing.F) {
	font, err := os.ReadFile(filepath.Join(vectortest.Fonts[0], "MPLUS1p-Regular-subset.ttf"))
	if err != nil {
		f.Fatal(err)
	}
	seed := func(draw func(doc *bdf.Document, o *bdf.Object)) {
		doc := bdf.NewDocument()
		o := bdf.NewObject()
		draw(doc, o)
		o.SetBBox(0, 0, 100, 100)
		f.Add(o.Encode())
	}
	seed(func(doc *bdf.Document, o *bdf.Object) {
		o.Save().Translate(10, 20).FillColor(bdf.RGBA(255, 0, 0, 128)).FillRect(0, 0, 30, 40).Restore()
		p := o.AddPath((&bdf.Path{}).MoveTo(0, 0).LineTo(10, 0).QuadTo(20, 0, 20, 10).Close().Rect(1, 2, 3, 4).Circle(5, 5, 3).ArcTo(1, 1, 2, 2, 1).RoundRect(0, 0, 9, 9, 2))
		o.Line(2, 1, 2, 4).Dash([]float32{3, 1, 2}, 1).StrokePath(p).FillPath(p, 1).ClipPath(p, 0).FillPathAt(p, 0, 3, 3)
		o.FillPathRun(0, []bdf.Glyph{{Path: p, X: 1, Y: 2}, {Path: p, X: 5, Y: 2}})
	})
	seed(func(doc *bdf.Document, o *bdf.Object) {
		lin := o.AddPaint(bdf.LinearGradient(0, 0, 100, 0, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGBA(0, 0, 255, 0)}))
		rad := o.AddPaint(bdf.RadialGradient(50, 50, 0, 50, 50, 40, bdf.Stop{Offset: 0, Color: bdf.RGB(1, 2, 3)}, bdf.Stop{Offset: 0.5, Color: bdf.RGB(9, 8, 7)}))
		con := o.AddPaint(bdf.ConicGradient(0, 50, 50, bdf.Stop{Offset: 0, Color: bdf.RGB(255, 0, 0)}, bdf.Stop{Offset: 1, Color: bdf.RGB(0, 0, 255)}))
		o.FillPaint(lin).FillRect(0, 0, 100, 50).StrokePaint(rad).StrokeRect(0, 0, 7, 7).FillPaint(con).FillRect(0, 0, 9, 9)
		o.GroupBegin(0.5, 2, 10, 10, 80, 60).FillRect(10, 10, 80, 60).MaskBegin(1, 0, nil).FillRect(20, 20, 30, 30).MaskEnd().GroupEnd()
		o.Alpha(0.3).Blend(18).Shadow(255, 4, 1, 1).ClearRect(0, 0, 4, 4).Link(1, 1, 5, 5, "#page=1").Mark(bdf.MarkAltText, "x")
	})
	seed(func(doc *bdf.Document, o *bdf.Object) {
		fn := o.AddFont(bdf.SystemFont("M PLUS 1p", 700, bdf.StyleItalic))
		o.Font(fn, 12).TextStyle(2, 1, 2, 1).FillText("Hello 日本 שלום", 20, 40, 60).StrokeText("x", 1, 2, 0)
	})
	seed(func(doc *bdf.Document, o *bdf.Object) {
		// groups in a state with a shadow, one in the other, with text and a line
		fn := o.AddFont(bdf.SystemFont("M PLUS 1p", 400, bdf.StyleNormal))
		o.Shadow(0x00000080, 3, 2, 2).Font(fn, 10).GroupBegin(0.8, 0, 5, 5, 90, 90).FillRect(10, 10, 30, 30)
		o.Shadow(0xff0000ff, 0, 1, 1).GroupBegin(1, 0, 20, 20, 60, 60).FillText("a", 30, 40, 0).StrokeRect(25, 25, 20, 20).GroupEnd()
		o.MaskBegin(0, 0, nil).FillRect(0, 0, 50, 50).MaskEnd().GroupEnd().GroupEnd()
	})
	seed(func(doc *bdf.Document, o *bdf.Object) {
		// the MARKs of a structure: a heading, a list, a table, a figure, a group in them, ends that open nothing
		fn := o.AddFont(bdf.SystemFont("M PLUS 1p", 400, bdf.StyleNormal))
		o.Font(fn, 10).Mark(bdf.MarkHeading, "1").FillText("h", 5, 10, 0).Mark(bdf.MarkList, "").Mark(bdf.MarkListItem, "").FillText("i", 5, 20, 0)
		o.Save().Mark(bdf.MarkTable, "").Mark(bdf.MarkCell, "A1 col").FillText("c", 5, 30, 0).Restore().Mark(bdf.MarkCell, "B2:C3").FillRect(0, 0, 5, 5)
		o.Mark(bdf.MarkFigure, "alt").GroupBegin(0.5, 0, 0, 0, 50, 50).FillText("g", 5, 40, 0).Mark(bdf.MarkEnd, "").GroupEnd().Mark(bdf.MarkEnd, "").Mark(bdf.MarkEnd, "")
		o.Mark(bdf.MarkEnd, "").Mark(bdf.MarkLang, "ja").Mark(bdf.MarkAltText, "x").FillRect(1, 1, 2, 2).Link(1, 1, 5, 5, "https://example.com/").Mark(bdf.MarkCell, "B5").FillText("s", 5, 50, 0)
	})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<16 {
			return
		}
		doc := bdf.NewDocument()
		fh := doc.AddFont(font)
		_ = fh
		h := doc.AddPart(bdf.PartObject, data)
		v := doc.NewView("pages", bdf.ViewFixed, "")
		p := v.AddPage(100, 100)
		p.Layers = append(p.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: h}, bdf.Layer{Role: bdf.RoleAnnotation, Obj: h})
		var buf bytes.Buffer
		res, err := Write(&buf, doc, testOptions())
		if err != nil {
			t.Fatal(err)
		}
		file := buf.String()
		if !strings.HasPrefix(file, "%PDF-1.7") || !strings.HasSuffix(file, "%%EOF\n") || res.Pages != 1 {
			t.Fatalf("not a file of a page: %d bytes", len(file))
		}
		// as a tagged PDF/A file: whatever the MARKs, the marked content
		// of every stream is closed
		o := testOptions()
		o.PDFA, o.NoCompression = "2a", true
		buf.Reset()
		if _, err := Write(&buf, doc, o); err != nil {
			t.Fatal(err)
		}
		open := 0
		for _, line := range strings.Split(buf.String(), "\n") {
			switch {
			case strings.HasSuffix(line, " BDC") && strings.HasPrefix(line, "/") || line == "/Artifact BMC":
				open++
			case line == "EMC":
				open--
			case line == "endstream":
				if open != 0 {
					t.Fatalf("%d sequences of marked content are left open in a stream", open)
				}
			}
		}
	})
}

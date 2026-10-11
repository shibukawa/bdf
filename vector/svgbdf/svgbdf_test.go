package svgbdf

import (
	"bytes"
	"encoding/xml"
	"image/color"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// png1 is a PNG image of one pixel.
var png1 = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x04\x00\x00\x00\x02\x08\x06\x00\x00\x00\x7f\xa8}c")

// draw returns the SVG of a page of 200 × 100 units with one layer.
func draw(t *testing.T, opts *Options, build func(doc *bdf.Document, o *bdf.Object)) (string, *Renderer) {
	t.Helper()
	doc := bdf.NewDocument()
	o := bdf.NewObject()
	build(doc, o)
	h, _ := doc.AddObject(o)
	v := doc.NewView("main", bdf.ViewFixed, "")
	v.AddPage(200, 100, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	r := New(doc, opts)
	svg, err := r.Page(v, 0)
	if err != nil {
		t.Fatal(err)
	}
	wellFormed(t, "page", svg)
	return string(svg), r
}

// wellFormed checks that an SVG is XML, and that what it refers to by id
// is defined in it once.
func wellFormed(t *testing.T, name string, svg []byte) {
	t.Helper()
	ids := map[string]bool{}
	dec := xml.NewDecoder(bytes.NewReader(svg))
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("%s: not XML: %v", name, err)
		}
		if el, ok := tok.(xml.StartElement); ok {
			for _, a := range el.Attr {
				if a.Name.Local == "id" {
					if ids[a.Value] {
						t.Errorf("%s: id %q is defined twice", name, a.Value)
					}
					ids[a.Value] = true
				}
			}
		}
	}
	for _, m := range regexp.MustCompile(`url\(#([^)]+)\)|href="#([^"]+)"`).FindAllSubmatch(svg, -1) {
		if id := string(m[1]) + string(m[2]); !ids[id] {
			t.Errorf("%s: id %q is referred to and not defined", name, id)
		}
	}
}

func contains(t *testing.T, svg string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(svg, w) {
			t.Errorf("no %s in\n%s", w, svg)
		}
	}
}

func TestShapes(t *testing.T) {
	svg, _ := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		p := o.AddPath((&bdf.Path{}).MoveTo(10, 10).LineTo(50, 10).QuadTo(60, 20, 50, 30).Close())
		o.FillColor(0x3366ff80).StrokeColor(0x112233ff).Line(2, bdf.CapRound, bdf.JoinBevel, 10).Dash([]float32{4, 2}, 1)
		// a fill and the stroke of the same outline are one element
		o.FillPath(p, bdf.EvenOdd).StrokePath(p)
		o.Translate(100, 0).FillRect(0, 0, 20, 10)
		o.StrokeRect(0, 20, 20, 10)
	})
	contains(t, svg,
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 100" width="200pt" height="100pt">`,
		`<path d="M0 0h200v100h-200Z" fill="#fff"/>`,
		`<path d="M10 10L50 10Q60 20 50 30Z" fill-rule="evenodd" fill="#36f" fill-opacity=".501961" stroke="#123" stroke-width="2" stroke-linecap="round" stroke-linejoin="bevel" stroke-dasharray="4 2" stroke-dashoffset="1"/>`,
		`<path d="M0 0h20v10h-20Z" transform="translate(100 0)" fill="#36f" fill-opacity=".501961"/>`,
		`<path d="M0 20h20v10h-20Z" transform="translate(100 0)" fill="none" stroke="#123" stroke-width="2"`,
	)
	if n := strings.Count(svg, "M10 10L50 10"); n != 1 {
		t.Errorf("the outline is written %d times", n)
	}
}

func TestPathData(t *testing.T) {
	for _, c := range []struct {
		p    *bdf.Path
		want string
	}{
		// a line without a sub-path starts one, as in Canvas
		{(&bdf.Path{}).LineTo(1, 2).LineTo(3, 4), "M1 2L3 4"},
		{(&bdf.Path{}).Rect(1, 2, 3, 4).LineTo(9, 9), "M1 2h3v4h-3ZL9 9"},
		{(&bdf.Path{}).CubicTo(1, 2, 3, 4, 5, 6), "M1 2C1 2 3 4 5 6"},
		{(&bdf.Path{}).MoveTo(0, 0).LineTo(-1.5, -0.25), "M0 0L-1.5-.25"},
		// a full circle: three arcs, each one reading
		{(&bdf.Path{}).Circle(0, 0, 10), "M10 0A10 10 0 0 1 -5 8.66A10 10 0 0 1 -5 -8.66A10 10 0 0 1 10 0"},
		{(&bdf.Path{}).RoundRect(0, 0, 10, 10, 2), "M2 0L8 0L8 0A2 2 0 0 1 10 2L10 8L10 8A2 2 0 0 1 8 10L2 10L2 10A2 2 0 0 1 0 8L0 2L0 2A2 2 0 0 1 2 0Z"},
	} {
		if got, _ := pathData(c.p, 2); got != c.want {
			t.Errorf("%v: %s, want %s", c.p.Verbs, got, c.want)
		}
	}
	if _, ok := pathData(&bdf.Path{}, 2); ok {
		t.Error("an empty path has data")
	}
}

func TestText(t *testing.T) {
	font := []byte("wOF2 not a font")
	svg, r := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		f := o.AddFont(bdf.Font{Kind: bdf.FontEmbedded, Hash: doc.AddFont(font), Family: `"Noto Sans", sans-serif`, Weight: 700, Style: bdf.StyleItalic})
		sys := o.AddFont(bdf.SystemFont("serif", 400, bdf.StyleNormal))
		o.Font(f, 12).FillColor(0xff0000ff)
		o.FillText("a < b & c ", 10, 20, 50).FillText("next", 60, 20, 0)
		o.Font(sys, 10).TextStyle(bdf.AlignRight, bdf.BaselineMiddle, bdf.DirLTR, 1.5).StrokeText("tab\there", 190, 40, 0)
	})
	hash := bdf.HashOf(font).String()
	contains(t, svg,
		`@font-face{font-family:"bdf-`+hash[:10]+`";src:url("`+hash+`.woff2")}`,
		// the runs of a font and color share a group
		`<g font-family="bdf-`+hash[:10]+`,&quot;Noto Sans&quot;, sans-serif" font-size="12" font-weight="700" font-style="italic" fill="#f00">`+"\n"+
			`<text x="10" y="20" textLength="50" lengthAdjust="spacingAndGlyphs" xml:space="preserve">a &lt; b &amp; c </text>`+"\n"+
			`<text x="60" y="20">next</text>`+"\n</g>",
		`<g font-family="serif" font-size="10" fill="none" stroke="#000" stroke-miterlimit="10" text-anchor="end" dominant-baseline="central" letter-spacing="1.5">`+"\n"+
			`<text x="190" y="40">tab here</text>`,
	)
	if f := r.Files(); len(f) != 1 || f[0].Name != hash+".woff2" || !bytes.Equal(f[0].Data, font) {
		t.Errorf("files: %v", f)
	}

	svg, r = draw(t, &Options{Inline: true}, func(doc *bdf.Document, o *bdf.Object) {
		o.Font(o.AddFont(bdf.EmbeddedFont(doc.AddFont(font), 400, bdf.StyleNormal)), 12).FillText("x", 0, 0, 0)
	})
	contains(t, svg, `src:url("data:font/woff2;base64,d09GMiBub3QgYSBmb250")`, `<text>x</text>`)
	if len(r.Files()) != 0 {
		t.Errorf("a file beside an SVG that stands alone: %v", r.Files())
	}
}

func TestClips(t *testing.T) {
	svg, _ := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		for i := 0; i < 2; i++ {
			o.Save().ClipRect(0, 0, 50, 50).Translate(5, 5).ClipRect(0, 0, 10, 10)
			o.FillRect(0, 0, 100, 100).FillRect(1, 1, 2, 2)
			o.Restore()
		}
		o.FillRect(0, 0, 1, 1)
	})
	contains(t, svg,
		`<clipPath id="c1"><path d="M0 0h50v50h-50Z"/></clipPath>`,
		// a clip narrows the one before it, in the coordinates it was set in
		`<clipPath id="c2" clip-path="url(#c1)"><path d="M0 0h10v10h-10Z" transform="translate(5 5)"/></clipPath>`,
		// the elements of a clip share its group, across what restores and sets it again
		`<g clip-path="url(#c2)">`+"\n"+strings.Repeat(
			`<path d="M0 0h100v100h-100Z" transform="translate(5 5)"/>`+"\n"+`<path d="M1 1h2v2h-2Z" transform="translate(5 5)"/>`+"\n", 2)+"</g>\n"+
			`<path d="M0 0h1v1h-1Z"/>`,
	)
	if strings.Contains(svg, `id="c3"`) {
		t.Errorf("the same clip is defined again in\n%s", svg)
	}
}

func TestPaints(t *testing.T) {
	svg, r := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		lin := o.AddPaint(bdf.LinearGradient(0, 0, 100, 0, bdf.Stop{Offset: 1, Color: 0x0000ff80}, bdf.Stop{Offset: 0, Color: 0xff0000ff}))
		rad := o.AddPaint(bdf.RadialGradient(10, 10, 2, 20, 20, 30, bdf.Stop{Offset: 0, Color: 0xffffffff}))
		none := o.AddPaint(bdf.LinearGradient(5, 5, 5, 5, bdf.Stop{Offset: 0, Color: 0xffffffff}))
		pat := bdf.Pattern(o.AddImage(doc.AddImage(png1)), bdf.RepeatX)
		pat.Matrix = [6]float32{2, 0, 0, 2, 3, 4}
		o.FillPaint(lin).FillRect(0, 0, 10, 10).FillRect(10, 0, 10, 10)
		o.StrokePaint(rad).StrokeRect(0, 0, 10, 10)
		o.FillPaint(none).FillRect(20, 0, 10, 10)
		o.FillPaint(o.AddPaint(pat)).Alpha(0.5).FillRect(30, 0, 10, 10)
		o.FillPaint(o.AddPaint(bdf.ConicGradient(0, 5, 5, bdf.Stop{Offset: 0, Color: 0x00ff00ff}))).FillRect(40, 0, 10, 10)
	})
	img := bdf.HashOf(png1).String() + ".png"
	contains(t, svg,
		`<linearGradient id="g1" gradientUnits="userSpaceOnUse" x1="0" y1="0" x2="100" y2="0"><stop offset="0" stop-color="#f00"/><stop offset="1" stop-color="#00f" stop-opacity=".501961"/></linearGradient>`,
		`<radialGradient id="g2" gradientUnits="userSpaceOnUse" cx="20" cy="20" r="30" fx="10" fy="10" fr="2"><stop offset="0" stop-color="#fff"/></radialGradient>`,
		`<path d="M0 0h10v10h-10Z" fill="url(#g1)"/>`, `<path d="M10 0h10v10h-10Z" fill="url(#g1)"/>`,
		`fill="none" stroke="url(#g2)"`,
		`<path d="M20 0h10v10h-10Z" fill="none"/>`,
		`<pattern id="g3" patternUnits="userSpaceOnUse" width="4" height="32768" patternTransform="matrix(2 0 0 2 3 4)"><image href="`+img+`" width="4" height="2" preserveAspectRatio="none"/></pattern>`,
		`<path d="M30 0h10v10h-10Z" fill="url(#g3)" fill-opacity=".5"/>`,
		`<path d="M40 0h10v10h-10Z" fill="#0f0" fill-opacity=".5"/>`,
	)
	if w := r.Warnings(); len(w) != 1 || !strings.Contains(w[0], "conic") {
		t.Errorf("warnings: %q", w)
	}
}

func TestImages(t *testing.T) {
	build := func(doc *bdf.Document, o *bdf.Object) {
		i := o.AddImage(doc.AddImage(png1))
		o.Image(i, 10, 20, 40, 20).Smoothing(false, 0).Alpha(0.5).Image(i, 60, 20, -40, 20)
		o.Alpha(1).Smoothing(true, 0).ImageSub(i, 1, 0, 2, 2, 0, 50, 20, 20)
		o.Image(o.AddImage(doc.AddImage([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="8" height="6"/>`))), 0, 0, 8, 6)
	}
	svg, r := draw(t, nil, build)
	img := bdf.HashOf(png1).String() + ".png"
	contains(t, svg,
		`<image href="`+img+`" x="10" y="20" width="40" height="20" preserveAspectRatio="none"/>`,
		// mirrored by a rectangle of a negative width
		`<image href="`+img+`" x="-60" y="20" width="40" height="20" preserveAspectRatio="none" transform="matrix(-1 0 0 1 0 0)" opacity=".5" style="image-rendering:pixelated"/>`,
		`<g><svg y="50" width="20" height="20" viewBox="1 0 2 2" preserveAspectRatio="none"><image href="`+img+`" width="4" height="2" preserveAspectRatio="none"/></svg>`,
		// an SVG image is never a file of its own
		`<image href="data:image/svg+xml;base64,`,
	)
	if f := r.Files(); len(f) != 1 || f[0].Name != img {
		t.Errorf("files: %v", f)
	}

	svg, r = draw(t, &Options{Inline: true}, build)
	contains(t, svg,
		`<symbol id="i1" viewBox="0 0 4 2" preserveAspectRatio="none"><image href="data:image/png;base64,`,
		`<use href="#i1" x="10" y="20" width="40" height="20"/>`,
		`<use href="#i1" width="4" height="2"/></svg>`,
	)
	if n := strings.Count(svg, "data:image/png"); n != 1 || len(r.Files()) != 0 {
		t.Errorf("the image is in the SVG %d times, with files %v", n, r.Files())
	}
}

func TestGroups(t *testing.T) {
	svg, _ := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		o.ClipRect(0, 0, 100, 100)
		// a group that changes nothing is its content
		o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 10, 10).FillRect(0, 0, 1, 1).GroupEnd()
		o.GroupBegin(0.5, bdf.BlendSourceOver, 0, 0, 10, 10).FillRect(0, 0, 2, 2).GroupEnd()
		o.GroupBegin(1, bdf.BlendMultiply, 0, 0, 10, 10).FillRect(0, 0, 3, 3).GroupEnd()
		// a soft mask of luminance with a transfer function
		transfer := make([]byte, 256)
		for i := range transfer {
			transfer[i] = byte(255 - i)
		}
		o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 10, 10).FillRect(0, 0, 4, 4)
		o.MaskBegin(bdf.MaskLuminosity, 0x102030ff, transfer).FillRect(0, 0, 5, 5).MaskEnd().GroupEnd()
		o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 10, 10).FillRect(0, 0, 6, 6)
		o.MaskBegin(bdf.MaskAlpha, 0, nil).FillRect(0, 0, 7, 7).MaskEnd().GroupEnd()
		// what blends does so with what is under it: it is in no group of a clip
		o.Blend(bdf.BlendScreen).FillRect(0, 0, 8, 8)
	})
	contains(t, svg,
		`<g clip-path="url(#c1)">`+"\n"+`<path d="M0 0h1v1h-1Z"/>`+"\n"+`<g opacity=".5">`+"\n"+`<path d="M0 0h2v2h-2Z"/>`+"\n</g>\n</g>\n",
		`<g clip-path="url(#c1)" style="mix-blend-mode:multiply"><g>`+"\n"+`<path d="M0 0h3v3h-3Z"/>`+"\n</g></g>",
		`<feColorMatrix values=".3 .59 .11 0 0 .3 .59 .11 0 0 .3 .59 .11 0 0 0 0 0 1 0"/><feComponentTransfer><feFuncR type="table" tableValues="1 .9961 `,
		`<mask id="m2" maskUnits="userSpaceOnUse" x="0" y="0" width="200" height="100">`+"\n"+`<g filter="url(#f1)">`+"\n"+`<path d="M0 0h200v100h-200Z" fill="#102030"/>`+"\n"+`<path d="M0 0h5v5h-5Z"/>`,
		`<g mask="url(#m2)">`+"\n"+`<path d="M0 0h4v4h-4Z"/>`,
		`<mask id="m3" maskUnits="userSpaceOnUse" x="0" y="0" width="200" height="100" style="mask-type:alpha">`+"\n"+`<path d="M0 0h7v7h-7Z"/>`,
		`<g clip-path="url(#c1)" style="mix-blend-mode:screen"><path d="M0 0h8v8h-8Z"/></g>`,
		`height="100pt" style="isolation:isolate">`,
	)
}

func TestErase(t *testing.T) {
	svg, r := draw(t, &Options{Background: color.Transparent}, func(doc *bdf.Document, o *bdf.Object) {
		o.GroupBegin(1, bdf.BlendSourceOver, 0, 0, 100, 100)
		o.FillColor(0x008000ff).FillRect(0, 0, 100, 100)
		// two holes are one mask of what was drawn before them
		o.Blend(bdf.BlendDestinationOut).FillColor(0xffffff80).FillRect(10, 10, 5, 5).ClearRect(30, 30, 5, 5)
		o.Blend(bdf.BlendSourceOver).FillRect(50, 50, 5, 5)
		o.Blend(bdf.BlendDestinationOver).FillRect(60, 60, 5, 5)
		o.Blend(bdf.BlendXor).FillRect(70, 70, 5, 5)
		o.GroupEnd()
	})
	contains(t, svg,
		`<mask id="m1" maskUnits="userSpaceOnUse" x="0" y="0" width="200" height="100"><path d="M0 0h200v100h-200Z" fill="#fff"/>`+
			`<g><path d="M10 10h5v5h-5Z" fill-opacity=".501961"/></g><g><path d="M30 30h5v5h-5Z"/></g></mask>`,
		`<g><path d="M60 60h5v5h-5Z" fill="#fff" fill-opacity=".501961"/></g>`+"\n"+
			`<g mask="url(#m1)">`+"\n"+`<path d="M0 0h100v100h-100Z" fill="#008000"/>`+"\n</g>\n"+
			`<path d="M50 50h5v5h-5Z" fill="#fff" fill-opacity=".501961"/>`+"\n"+
			`<path d="M70 70h5v5h-5Z" fill="#fff" fill-opacity=".501961"/>`,
	)
	if w := r.Warnings(); len(w) != 1 || !strings.Contains(w[0], "xor") {
		t.Errorf("warnings: %q", w)
	}
}

func TestGlyphsAndShadows(t *testing.T) {
	svg, _ := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		ext := doc.AddPaths([]*bdf.Path{(&bdf.Path{}).Rect(0, 0, 500, 700), (&bdf.Path{}).Rect(0, 0, 300, 300)})
		a, b := o.AddExtPath(ext, 0), o.AddExtPath(ext, 1)
		o.Save().Translate(10, 50).Scale(0.01, -0.01)
		o.FillPathRun(bdf.NonZero, []bdf.Glyph{{Path: a}, {Path: b, X: 600}, {Path: a, X: 1000, Y: 5}})
		o.FillPathAt(b, bdf.NonZero, 1600, 0)
		o.Restore()
		o.Scale(2, 2).Shadow(0x00000080, 4, 1, -1).FillRect(0, 0, 10, 10)
	})
	contains(t, svg,
		`<path id="a1" d="M0 0h500v700h-500Z"/>`, `<path id="a2" d="M0 0h300v300h-300Z"/>`,
		`<g transform="matrix(.01 0 0 -.01 10 50)">`+"\n"+`<use href="#a1"/>`+"\n"+`<use href="#a2" x="600"/>`+"\n"+`<use href="#a1" x="1000" y="5"/>`+"\n"+`<use href="#a2" x="1600"/>`+"\n</g>",
		// a shadow is in the units and directions of the page: twice those of the transform it was set under
		`<feDropShadow dx="2" dy="-2" stdDeviation="4" flood-color="#000" flood-opacity=".501961"/>`,
		`<g filter="url(#f3)"><path d="M0 0h10v10h-10Z" transform="matrix(2 0 0 2 0 0)"/></g>`,
	)
}

func TestLinks(t *testing.T) {
	doc := bdf.NewDocument()
	o := bdf.NewObject()
	o.Translate(5, 5)
	for _, u := range []string{"https://example.com/?a=1&b=2", "#page=3", "#view=two%20b&page=2", "#view=none", "javascript:alert(1)", "data:text/html,x", "#", "#page=0"} {
		o.Link(0, 0, 10, 10, u)
	}
	h, _ := doc.AddObject(o)
	doc.NewView("one", bdf.ViewFixed, "").AddPage(100, 100, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	doc.NewView("two b", bdf.ViewFixed, "").AddPage(100, 100)
	svg, err := New(doc, nil).Page(doc.Views[0], 0)
	if err != nil {
		t.Fatal(err)
	}
	wellFormed(t, "links", svg)
	var got []string
	for _, m := range regexp.MustCompile(`<a href="([^"]*)">`).FindAllSubmatch(svg, -1) {
		got = append(got, string(m[1]))
	}
	if want := "https://example.com/?a=1&amp;b=2 view-1-page-3.svg view-2-page-2.svg"; strings.Join(got, " ") != want {
		t.Errorf("links %q, want %q", got, want)
	}
	contains(t, string(svg), `<path d="M0 0h10v10h-10Z" transform="translate(5 5)" fill="none" pointer-events="all"/></a>`)

	svg, _ = New(doc, &Options{PageURL: func(view, page int) string { return "" }}).Page(doc.Views[0], 0)
	if n := bytes.Count(svg, []byte("<a ")); n != 1 {
		t.Errorf("%d links without the URLs of pages", n)
	}
}

func TestViews(t *testing.T) {
	doc := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillRect(0, 0, 10, 10)
	h, _ := doc.AddObject(o)
	scroll := doc.NewView("scroll", bdf.ViewScroll, "")
	for i := 0; i < 2; i++ {
		p := scroll.AddPage(100, 120, bdf.Layer{Role: bdf.RoleHeader, Obj: h}, bdf.Layer{Role: bdf.RoleBody, Obj: h})
		p.Body = &bdf.RectDef{X: 10, Y: 20, W: 80, H: 90}
	}
	sheet := doc.NewView("sheet", bdf.ViewSheet, "")
	sheet.Tile, sheet.Gridlines = 100, true
	sheet.Cols, sheet.Rows = []bdf.Run{{3, 50}}, []bdf.Run{{2, 20}, {1, 0}, {4, 20}}
	sheet.Tiles = map[string]string{"0,0": h.String(), "1,0": h.String(), "bad": h.String()}

	r := New(doc, nil)
	// a strip is cut to its body, and has the layers of it
	svg, err := r.Page(scroll, 1)
	if err != nil {
		t.Fatal(err)
	}
	contains(t, string(svg), `viewBox="10 20 80 90" width="80pt" height="90pt"`, `<path d="M10 20h80v90h-80Z" fill="#fff"/>`)
	if n := bytes.Count(svg, []byte("M0 0h10v10h-10Z")); n != 1 {
		t.Errorf("%d layers of a strip are drawn", n)
	}
	// the strips stacked: the second begins where the body of the first ends
	svg, err = r.Region(scroll, 0, bdf.Rect{X: 0, Y: 80, W: 80, H: 50})
	if err != nil {
		t.Fatal(err)
	}
	wellFormed(t, "scroll", svg)
	contains(t, string(svg),
		`<clipPath id="c2"><path d="M0 0h80v90h-80Z" transform="translate(0 90)"/></clipPath>`,
		`<g clip-path="url(#c2)">`+"\n"+`<path d="M0 90h80v90h-80Z" fill="#fff"/>`+"\n"+`<path d="M0 0h10v10h-10Z" transform="translate(-10 70)"/>`,
	)

	if b := SheetBounds(sheet); b != (bdf.Rect{W: 150, H: 100}) {
		t.Errorf("sheet bounds %v", b)
	}
	if _, err := r.Page(sheet, 0); err == nil {
		t.Error("a sheet has a page")
	}
	svg, err = r.Region(sheet, 0, SheetBounds(sheet))
	if err != nil {
		t.Fatal(err)
	}
	wellFormed(t, "sheet", svg)
	contains(t, string(svg),
		// a hidden row draws no second line, and the lines end with the tiles
		`<path d="M50 0V100M100 0V100M150 0V100M0 20H150M0 40H150M0 60H150M0 80H150M0 100H150" fill="none" stroke="#d9d9d9" vector-effect="non-scaling-stroke" shape-rendering="crispEdges"/>`,
		`<clipPath id="c2"><path d="M0 0h100v100h-100Z" transform="translate(100 0)"/></clipPath>`,
		`<g clip-path="url(#c2)">`+"\n"+`<path d="M0 0h10v10h-10Z" transform="translate(100 0)"/>`,
	)
}

func TestLimits(t *testing.T) {
	// an object that draws itself: as deep as objects may be, and no deeper
	doc := bdf.NewDocument()
	leaf := bdf.NewObject()
	leaf.FillRect(0, 0, 1, 1)
	h, box := doc.AddObject(leaf)
	for i := 0; i < 12; i++ {
		o := bdf.NewObject()
		ref := o.AddObject(h, box)
		for j := 0; j < 8; j++ {
			o.Use(ref)
		}
		h, box = doc.AddObject(o)
	}
	v := doc.NewView("main", bdf.ViewFixed, "")
	v.AddPage(10, 10, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	r := New(doc, nil)
	svg, err := r.Page(v, 0)
	if err != nil {
		t.Fatal(err)
	}
	wellFormed(t, "limits", svg)
	if n := bytes.Count(svg, []byte("<path")); n > maxReusedInstructions+16 {
		t.Errorf("%d elements", n)
	}
	if w := r.Warnings(); len(w) != 1 || !strings.Contains(w[0], "too many times") {
		t.Errorf("warnings: %q", w)
	}

	// groups and states left open, and ended without a beginning
	svg2, _ := draw(t, nil, func(doc *bdf.Document, o *bdf.Object) {
		o.GroupEnd().MaskEnd().Restore()
		for i := 0; i < maxLayers+4; i++ {
			o.GroupBegin(0.5, bdf.BlendSourceOver, 0, 0, 1, 1).FillRect(0, 0, 1, 1)
		}
		o.MaskBegin(bdf.MaskAlpha, 0, nil).Save().FillRect(0, 0, 2, 2)
	})
	if n := strings.Count(svg2, "<g opacity"); n != maxLayers {
		t.Errorf("%d groups", n)
	}
}

// TestSamples writes the documents of testdata, and checks that what is
// written is XML that refers to nothing it does not define.
func TestSamples(t *testing.T) {
	root := filepath.Join("..", "..", "testdata")
	n := 0
	err := filepath.WalkDir(root, func(path string, e fs.DirEntry, err error) error {
		if err != nil || e.IsDir() || filepath.Ext(path) != ".bdf" || strings.Contains(path, "encrypted") {
			return err
		}
		rd, err := bdf.OpenSingleFile(path)
		if err != nil {
			return nil // not a document of one file
		}
		defer rd.Close()
		doc, err := rd.ToDocument()
		if err != nil {
			return nil
		}
		n++
		for _, inline := range []bool{false, true} {
			dir := t.TempDir()
			res, err := Write(dir, doc, &Options{Inline: inline})
			if err != nil {
				t.Errorf("%s: %v", path, err)
				return nil
			}
			for _, name := range res.Pages {
				svg, err := os.ReadFile(filepath.Join(dir, name))
				if err != nil {
					t.Fatal(err)
				}
				wellFormed(t, path+": "+name, svg)
				// what an SVG refers to beside it is there
				for _, m := range regexp.MustCompile(`(?:href="|url\(")([0-9a-f]{32}\.[a-z0-9]+)"`).FindAllSubmatch(svg, -1) {
					if _, err := os.Stat(filepath.Join(dir, string(m[1]))); err != nil {
						t.Errorf("%s: %s: %v", path, name, err)
					}
				}
			}
			if inline && len(res.Files) != 0 {
				t.Errorf("%s: files beside SVGs that stand alone: %v", path, res.Files)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Skip("no documents in testdata")
	}
}

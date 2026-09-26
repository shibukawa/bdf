package drawio

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
)

// testFonts lays text out with the PowerPoint converter's test fonts only,
// so that results do not depend on the fonts installed.
var testFonts = []string{filepath.Join("..", "pptx", "testdata", "fonts")}

func testOptions() *Options {
	return &Options{FontDirs: testFonts, NoSystemFonts: true}
}

// compress encodes a diagram the way draw.io does (Graph.compress).
func compress(t *testing.T, xml string) string {
	t.Helper()
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, flate.BestCompression)
	w.Write([]byte(url.PathEscape(xml)))
	w.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

const tinyModel = `<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` +
	`<mxCell id="a" value="Hello" style="rounded=0;whiteSpace=wrap;html=1;" vertex="1" parent="1"><mxGeometry x="10" y="10" width="80" height="40" as="geometry"/></mxCell>` +
	`</root></mxGraphModel>`

func TestReadFile(t *testing.T) {
	read := func(name string) *file {
		t.Helper()
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		f, err := readFile(data)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		return f
	}
	f := read("multipage.drawio")
	if len(f.pages) != 3 || f.pages[1].id != "details" || f.pages[1].name != "Details" || f.modified == "" {
		t.Fatalf("multipage: %d pages, %+v, modified %q", len(f.pages), f.pages[1], f.modified)
	}
	for _, name := range []string{"embedded.drawio.svg", "embedded.drawio.png"} {
		f := read(name)
		if len(f.pages) != 1 || f.pages[0].model == nil || len(parseModel(f.pages[0].model).cells) != 5 {
			t.Fatalf("%s: %d pages", name, len(f.pages))
		}
	}

	// compressed pages, a bare model, and junk
	comp := `<mxfile><diagram id="x" name="Zipped">` + compress(t, tinyModel) + `</diagram></mxfile>`
	f, err := readFile([]byte(comp))
	if err != nil || len(f.pages) != 1 || f.pages[0].name != "Zipped" || parseModel(f.pages[0].model).cells["a"].value != "Hello" {
		t.Fatalf("compressed: %v %+v", err, f)
	}
	if f, err := readFile([]byte(tinyModel)); err != nil || len(f.pages) != 1 {
		t.Fatalf("bare model: %v", err)
	}
	for _, bad := range []string{"", "hello", "<html><body/></html>", "<svg/>", "<mxfile></mxfile>"} {
		if _, err := readFile([]byte(bad)); err == nil {
			t.Errorf("%q accepted", bad)
		}
	}
}

func TestParseStyle(t *testing.T) {
	s := parseStyle("ellipse;whiteSpace=wrap;fillColor=none;strokeColor=default;fontColor=#FF0000;", false)
	if s["shape"] != "ellipse" || s["perimeter"] != "ellipsePerimeter" {
		t.Errorf("named style not merged: %v", s)
	}
	if s.has("fillColor") {
		t.Errorf("none should remove the key: %v", s["fillColor"])
	}
	if s["strokeColor"] != "#000000" || s["fontSize"] != "12" {
		t.Errorf("defaults: %v", s)
	}
	e := parseStyle("edgeStyle=orthogonalEdgeStyle;endArrow=none;", true)
	if e.has("endArrow") || e["shape"] != "connector" || e.num("fontSize", 0) != 11 {
		t.Errorf("edge defaults: %v", e)
	}
	if s := parseStyle(";fillColor=#fff;", false); s.has("shape") || s["fillColor"] != "#fff" {
		t.Errorf("leading ; must drop the defaults: %v", s)
	}
	if s := parseStyle("fillColor=default;defaultFillColor=invert;", false); s["fillColor"] != "#000000" {
		t.Errorf("inverted default fill: %v", s["fillColor"])
	}
	if v, ok := parseFloat("12px"); !ok || v != 12 {
		t.Errorf("parseFloat(12px) = %v %v", v, ok)
	}
}

func TestParseColor(t *testing.T) {
	for in, want := range map[string]rgba{
		"#fff": {255, 255, 255, 255}, "#12AbEf": {0x12, 0xab, 0xef, 255}, "#11223380": {0x11, 0x22, 0x33, 0x80},
		"red": {255, 0, 0, 255}, "rgb(1, 2, 3)": {1, 2, 3, 255}, "rgba(10,20,30,0.5)": {10, 20, 30, 128},
		"light-dark(#000000, #ffffff)": {0, 0, 0, 255}, "light-dark(rgb(1,2,3), red)": {1, 2, 3, 255},
	} {
		if got, ok := parseColor(in); !ok || got != want {
			t.Errorf("parseColor(%q) = %v %v, want %v", in, got, ok, want)
		}
	}
	for _, in := range []string{"none", "", "transparent", "#12", "nonsense"} {
		if _, ok := parseColor(in); ok {
			t.Errorf("parseColor(%q) accepted", in)
		}
	}
}

func TestParseHTML(t *testing.T) {
	root := parseHTML(`<div style="text-align: left">a &amp; b<br>c<ul><li>one<li>two</ul><p>x < y</div>`)
	div := root.kids[0]
	if div.tag != "div" || div.attrs["style"] == "" || div.kids[0].text != "a & b" || div.kids[1].tag != "br" {
		t.Fatalf("div: %+v", div)
	}
	var ul *hnode
	for _, k := range div.kids {
		if k.tag == "ul" {
			ul = k
		}
	}
	if ul == nil || len(ul.kids) != 2 || ul.kids[1].kids[0].text != "two" {
		t.Fatalf("list items are closed by the next item: %+v", ul)
	}
	p := div.kids[len(div.kids)-1]
	if p.tag != "p" || p.kids[0].text != "x < y" {
		t.Fatalf("a lone < is text: %+v", p)
	}
	if css := parseCSS("color: red; font-family: 'A; B', serif ; x"); css["color"] != "red" || css["font-family"] != "'A; B', serif" {
		t.Fatalf("css: %v", css)
	}
}

// convertTest converts a test file with the test fonts.
func convertTest(t *testing.T, name string) *Result {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// viewText returns the text of a view's page, runs joined by their separators.
func viewText(t *testing.T, doc *bdf.Document, v *bdf.View) string {
	t.Helper()
	resolve := func(h bdf.Hash) *bdf.ObjectPart {
		p := doc.Part(h)
		if p == nil {
			return nil
		}
		o, err := bdf.DecodeObject(p.Data)
		if err != nil {
			t.Fatal(err)
		}
		return o
	}
	var b strings.Builder
	for _, l := range v.Pages[0].Layers {
		runs, err := bdf.ExtractText(resolve(l.Obj), resolve)
		if err != nil {
			t.Fatal(err)
		}
		for i, r := range runs {
			if i > 0 || b.Len() > 0 {
				switch r.Sep {
				case bdf.SepSpace:
					b.WriteByte(' ')
				case bdf.SepBreak:
					b.WriteByte('\n')
				}
			}
			b.WriteString(r.Text)
		}
	}
	return b.String()
}

func TestConvertMultipage(t *testing.T) {
	res := convertTest(t, "multipage.drawio")
	doc := res.Doc
	if res.Pages != 3 || len(doc.Views) != 3 {
		t.Fatalf("%d pages, %d views", res.Pages, len(doc.Views))
	}
	want := []struct{ id, title string }{{"overview", "Overview"}, {"details", "Details"}, {"layers", "Layers"}}
	for i, v := range doc.Views {
		if v.ID != want[i].id || v.Title != want[i].title || v.Kind != bdf.ViewFixed || len(v.Pages) != 1 || v.TextIndex == "" {
			t.Errorf("view %d: %+v", i, v)
		}
	}
	if doc.Meta.Source != "drawio" || doc.Meta.DC.Modified.First() != "2026-09-26T03:00:00.000Z" {
		t.Errorf("meta: %+v", doc.Meta)
	}
	// the details page has a background layer; the hidden layer is left out
	if l := doc.Views[1].Pages[0].Layers; len(l) != 2 || l[0].Role != bdf.RoleBackground {
		t.Errorf("details layers: %+v", l)
	}
	if l := doc.Views[2].Pages[0].Layers; len(l) != 2 {
		t.Errorf("layers page: %d layers, want 2 (one hidden)", len(l))
	}
	text := viewText(t, doc, doc.Views[0])
	for _, s := range []string{"System overview", "Client", "API server\n(see Details)", "HTTPS", "SQL", "Database"} {
		if !strings.Contains(text, s) {
			t.Errorf("overview text %q lacks %q", text, s)
		}
	}
	if text := viewText(t, doc, doc.Views[2]); strings.Contains(text, "Hidden") || !strings.Contains(text, "On top") {
		t.Errorf("layers text %q", text)
	}
	// links between pages become #view= links; web links stay
	links := map[string]bool{}
	for _, v := range doc.Views {
		for _, l := range v.Pages[0].Layers {
			o, _ := bdf.DecodeObject(doc.Part(l.Obj).Data)
			o.Walk(func(in bdf.Instr) {
				if in.Op == bdf.OpLink {
					links[in.Args[4].(string)] = true
				}
			})
		}
	}
	for _, u := range []string{"#view=details", "#view=overview", "https://www.drawio.com/"} {
		if !links[u] {
			t.Errorf("link %s missing in %v", u, links)
		}
	}
}

func TestPageSelection(t *testing.T) {
	opts := testOptions()
	opts.Pages = conv.PageList(2)
	res, err := ConvertFile(filepath.Join("testdata", "multipage.drawio"), opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Doc.Views) != 1 || res.Doc.Views[0].ID != "details" {
		t.Fatalf("views %+v", res.Doc.Views)
	}
	opts.Pages = conv.PageList(4)
	if _, err := ConvertFile(filepath.Join("testdata", "multipage.drawio"), opts); err == nil {
		t.Fatal("page 4 of 3 accepted")
	}
}

package visio

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// The test drawings are made by test/visio/gen.py; text is laid out with
// the test fonts of the PowerPoint converter.

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func convert(t *testing.T, name string) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), testOptions())
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
	return res, r
}

func plainText(t *testing.T, r *bdf.Reader) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[0].TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	return bdf.PlainText(idx)
}

// instrs lists the instructions of an object with a given opcode.
func instrs(t *testing.T, r *bdf.Reader, h bdf.Hash, op byte) []bdf.Instr {
	t.Helper()
	o, err := r.Object(h)
	if err != nil {
		t.Fatal(err)
	}
	var out []bdf.Instr
	o.Walk(func(in bdf.Instr) {
		if in.Op == op {
			out = append(out, in)
		}
	})
	return out
}

func layer(p *bdf.Page, role string) bdf.Hash {
	for _, l := range p.Layers {
		if l.Role == role {
			return l.Obj
		}
	}
	return bdf.Hash{}
}

func TestConvertShapes(t *testing.T) {
	res, r := convert(t, "shapes.vsdx")
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	m := r.Manifest
	if m.Meta.Source != "vsdx" || m.Meta.DC.Title.First() != "Visio shapes" || m.Meta.DC.Language.First() != "en-US" {
		t.Errorf("meta = %+v", m.Meta)
	}
	pages := m.Views[0].Pages
	if len(pages) != 2 {
		t.Fatalf("%d pages, want 2 (the background page is not one)", len(pages))
	}
	for i, p := range pages {
		if p.W != 792 || p.H != 612 {
			t.Errorf("page %d is %g×%g", i+1, p.W, p.H)
		}
		if len(p.Layers) != 2 || p.Layers[0].Role != bdf.RoleBackground || p.Layers[1].Role != bdf.RoleBody {
			t.Errorf("page %d layers = %+v", i+1, p.Layers)
		}
	}
	if layer(pages[0], bdf.RoleBackground) != layer(pages[1], bdf.RoleBackground) {
		t.Error("the background page is not shared")
	}
	text := plainText(t, r)
	for _, want := range []string{"BDF from Visio: shapes.vsdx", "Theme quick styles", "Turned", "Left aligned text wraps",
		"bold", "日本語の文章は", "縦書きのテキスト", "connector label", "Nested item"} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q missing", want)
		}
	}
	if strings.Contains(text, "Hidden layer") {
		t.Error("a shape on a hidden layer is drawn")
	}
	links := map[string]bool{}
	for _, p := range pages {
		for _, in := range instrs(t, r, layer(p, bdf.RoleBody), bdf.OpLink) {
			links[in.Args[4].(string)] = true
		}
	}
	if !links["#page=2"] || !links["https://example.com/"] {
		t.Errorf("links = %v", links)
	}
	if n := len(instrs(t, r, layer(pages[0], bdf.RoleBody), bdf.OpImage)); n != 1 {
		t.Errorf("%d pictures drawn, want 1", n)
	}
	if n := len(instrs(t, r, layer(pages[0], bdf.RoleBody), bdf.OpShadow)); n == 0 {
		t.Error("no shadow drawn")
	}
}

// TestVDXMatchesVSDX converts the same drawing in both formats.
func TestVDXMatchesVSDX(t *testing.T) {
	a, ra := convert(t, "flow.vsdx")
	b, rb := convert(t, "flow.vdx")
	if len(a.Warnings)+len(b.Warnings) != 0 {
		t.Errorf("warnings: %v %v", a.Warnings, b.Warnings)
	}
	if rb.Manifest.Meta.Source != "vdx" || rb.Manifest.Meta.DC.Title.First() != "Visio flowchart" {
		t.Errorf("vdx meta = %+v", rb.Manifest.Meta)
	}
	pa, pb := ra.Manifest.Views[0].Pages, rb.Manifest.Views[0].Pages
	if len(pa) != 1 || len(pb) != 1 {
		t.Fatalf("pages %d and %d", len(pa), len(pb))
	}
	if len(pa[0].Layers) != 2 || len(pb[0].Layers) != 2 {
		t.Fatalf("layers %+v and %+v", pa[0].Layers, pb[0].Layers)
	}
	for i := range pa[0].Layers {
		if pa[0].Layers[i] != pb[0].Layers[i] {
			t.Errorf("layer %d differs: %v and %v", i, pa[0].Layers[i], pb[0].Layers[i])
		}
	}
	text := plainText(t, rb)
	for _, want := range []string{"Start", "Collect input", "Valid?", "Yes", "No", "Report error", "Flowchart test drawing"} {
		if !strings.Contains(text, want) {
			t.Errorf("text %q missing", want)
		}
	}
	if strings.Contains(text, "Hidden layer") {
		t.Error("a shape on a hidden layer is drawn")
	}
}

// sameValue compares cell values: numbers to 1e-9, others ignoring case.
func sameValue(a, b string) bool {
	if x, err := strconv.ParseFloat(a, 64); err == nil {
		if y, err := strconv.ParseFloat(b, 64); err == nil {
			return math.Abs(x-y) < 1e-9
		}
	}
	return strings.EqualFold(a, b)
}

// shapeWithText finds a shape by its text.
func shapeWithText(shapes []*shape, text string) *shape {
	for _, s := range shapes {
		if n := s.textNode(); n != nil && n.Content() == text {
			return s
		}
		if k := shapeWithText(s.kids, text); k != nil {
			return k
		}
	}
	return nil
}

func TestThemedValues(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "shapes.vsdx"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	st, _ := f.Stat()
	p, err := ooxml.Open(f, st.Size())
	if err != nil {
		t.Fatal(err)
	}
	d, err := readVSDX(p, func(k, m string) { t.Error(m) })
	if err != nil {
		t.Fatal(err)
	}
	c := &converter{opts: &Options{}, doc: bdf.NewDocument(), d: d, warned: map[string]bool{}, pageNum: map[int]int{}}
	pc := c.newPageCtx(d.pages[0], 612)
	for _, tc := range []struct {
		text, cell, want string
	}{
		{"V1", "FillForegnd", "#2E75B6"},                // variant style 1: solid variant color 1
		{"V1", "FillGradientEnabled", "0"},              //
		{"V2", "FillGradientEnabled", "1"},              // variant style 2: a gradient
		{"V3", "FillForegnd", "#1A446A"},                // variant color 3 shaded by 75 % (linear RGB)
		{"V1", "LineWeight", ftoa(1.0 / 72)},            // line style 3: 1 pt
		{"V4", "ShdwPattern", "1"},                      // effect style 4: a shadow
		{"V1", "ShdwPattern", "0"},                      //
		{"V1", "LinePattern", "1"},                      //
		{"Hidden layer", "FillForegnd", "#FF0000"},      // not themed
		{"Shadow", "ShapeShdwOffsetX", ftoa(0.06)},      //
		{"To page 2", "FillForegnd", "#2E75B6"},         // the Normal style's theme values
		{"To page 2", "QuickStyleFillColor", "100"},     //
		{"Theme quick styles", "FillPattern", "0"},      //
		{"FlipX", "FlipX", "1"},                         //
		{"Turned", "Angle", ftoa(math.Pi / 6)},          //
		{"V2", "FillGradientAngle", ftoa(-math.Pi / 2)}, // DrawingML's 90° downwards
		{"V1", "QuickStyleFillMatrix", "100"},           //
		{"V1", "Rounding", "0"},                         //
		{"V1", "LineCap", "1"},                          // flat
		{"V1", "LineColorTrans", "0"},                   //
		{"V1", "FillForegndTrans", "0"},                 //
		{"V1", "FillPattern", "1"},                      //
		{"V1", "FillBkgnd", "#2E75B6"},                  //
		{"V2", "FillBkgnd", "#6FA0D0"},                  // the last stop
		{"V1", "ShdwForegnd", "#2E75B6"},                //
		{"V1", "EndArrow", "0"},                         //
		{"V1", "BeginArrow", "0"},                       //
		{"V1", "LineColor", "#2766A0"},                  // line style 3: variant color 1 shaded by 75 %
		{"V1", "HideText", "0"},                         //
		{"V1", "LangID", "en-US"},                       //
		{"V1", "VerticalAlign", "1"},                    //
		{"V1", "TextDirection", "0"},                    //
		{"V1", "LeftMargin", ftoa(4.0 / 72)},            //
		{"Hidden layer", "LayerMember", "1"},            //
		{"Hidden layer", "LineWeight", ftoa(0.75 / 72)}, //
		{"Hidden layer", "LinePattern", "1"},            //
		{"Hidden layer", "FillBkgnd", "0"},              //
		{"Hidden layer", "QuickStyleFillColor", ""},     // the Normal style is not inherited
		{"Hidden layer", "ShdwPattern", "0"},            //
	} {
		s := shapeWithText(d.pages[0].shapes, tc.text)
		if s == nil {
			s = shapeWithText(d.pages[1].shapes, tc.text)
		}
		if s == nil {
			t.Fatalf("no shape with the text %q", tc.text)
		}
		if got, _ := pc.val(s, tc.cell); !sameValue(got, tc.want) {
			t.Errorf("%s %s = %q, want %q", tc.text, tc.cell, got, tc.want)
		}
	}
	// a connector takes its arrowhead from the connector scheme
	conn := shapeWithText(d.pages[1].shapes, "connector label")
	pt := c.newPageCtx(d.pages[1], 612)
	if !pt.quickStyle(conn).connector {
		t.Error("the connector is not one")
	}
	if got, _ := pt.val(conn, "EndArrow"); got != "13" {
		t.Errorf("connector EndArrow = %q, want 13", got)
	}
	// variant colors 200-206 are those of 100-106
	q := &quickStyle{varColor: 0}
	if pc.themeColor(q, 204) != pc.themeColor(q, 104) || pc.themeColor(q, 104) != 0xED7D31FF {
		t.Errorf("variant color 5 = %08x", uint32(pc.themeColor(q, 204)))
	}
}

func TestInheritance(t *testing.T) {
	root := &style{sheet: newSheet(), id: 0, name: "No Style"}
	root.cells["LineWeight"] = cell{v: "0.01"}
	root.cells["FillForegnd"] = cell{v: "1"}
	root.section("Character", -1).row("0").cells["Size"] = cell{v: "0.1"}
	child := &style{sheet: newSheet(), id: 1, name: "Child", line: root, fill: root, text: root}
	child.cells["FillForegnd"] = cell{v: "#FF0000"}
	m := &shape{sheet: newSheet(), lineStyle: child, fillStyle: child, textStyle: child}
	m.cells["Width"] = cell{v: "2"}
	g := m.section("Geometry", 0)
	g.cells["NoFill"] = cell{v: "1"}
	for i, xy := range [][2]string{{"0", "0"}, {"2", "0"}, {"2", "1"}} {
		r := g.row(itoa(i + 1))
		r.t = "LineTo"
		if i == 0 {
			r.t = "MoveTo"
		}
		r.cells["X"], r.cells["Y"] = cell{v: xy[0], f: "Width*" + ftoa(num(xy[0])/2)}, cell{v: xy[1]}
	}
	inst := &shape{sheet: newSheet(), inh: m}
	inst.cells["Width"] = cell{v: "4"}
	ig := inst.section("Geometry", 0)
	ig.row("2").cells["Y"] = cell{v: "0.5"}
	ig.row("3").del = true
	for _, c := range []struct{ name, want string }{{"Width", "4"}, {"FillForegnd", "#FF0000"}, {"LineWeight", "0.01"}} {
		if got, _ := inst.get(c.name); got.v != c.want {
			t.Errorf("%s = %q, want %q", c.name, got.v, c.want)
		}
	}
	if got, _ := inst.rowCell("Character", "0", "Size"); got.v != "0.1" {
		t.Errorf("Size = %q", got.v)
	}
	geo := inst.geometry()
	if len(geo) != 1 || len(geo[0].rows) != 2 || !geo[0].flag("NoFill") {
		t.Fatalf("geometry = %+v", geo)
	}
	r := geo[0].rows[1]
	r.w, r.h = 4, 1
	// X is inherited: its formula follows the instance's width
	if r.num("X") != 4 || r.num("Y") != 0.5 {
		t.Errorf("row 2 = (%g, %g), want (4, 0.5)", r.num("X"), r.num("Y"))
	}
	// cycles end
	a := &style{sheet: newSheet(), id: 2}
	b := &style{sheet: newSheet(), id: 3, line: a}
	a.line = b
	s := &shape{sheet: newSheet(), lineStyle: a}
	if _, ok := s.get("LineWeight"); ok {
		t.Error("a value out of a cycle")
	}
}

func TestSplitText(t *testing.T) {
	n, err := ooxml.Parse([]byte(`<Text><cp IX="0"/>one <cp IX="1"/>two<fld IX="0">F</fld>` + "\n" + `<pp IX="1"/>three` + " " + `four` + "\n" + `</Text>`))
	if err != nil {
		t.Fatal(err)
	}
	paras := splitText(n)
	if len(paras) != 2 {
		t.Fatalf("%d paragraphs", len(paras))
	}
	var got []string
	for _, pa := range paras {
		s := "pp" + pa.pp + ":"
		for _, r := range pa.runs {
			if r.br {
				s += "|br"
			} else {
				s += "|" + r.char + "=" + r.text
			}
		}
		got = append(got, s)
	}
	want := []string{"pp0:|0=one |1=two|1=F", "pp1:|1=three|br|1=four"}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("paragraphs =\n%s\nwant\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestGeometry(t *testing.T) {
	b := &pathBuilder{m: canvas.Identity}
	// a positive bow bulges to the right of the direction of travel
	b.moveTo(pt{0, 0})
	b.bowArc(pt{2, 0}, 0.5)
	if mid := bezierPoint(b.subs[0].segs, 0.5); math.Abs(mid.x-1) > 1e-6 || math.Abs(mid.y+0.5) > 1e-6 {
		t.Errorf("arc middle = %v, want (1, -0.5)", mid)
	}
	// an elliptical arc runs through its control point
	b = &pathBuilder{m: canvas.Identity}
	b.moveTo(pt{-2, 0})
	b.ellipticalArc(pt{2, 0}, pt{0, 1}, 0, 2)
	if mid := bezierPoint(b.subs[0].segs, 0.5); math.Abs(mid.x) > 1e-6 || math.Abs(mid.y-1) > 1e-6 {
		t.Errorf("elliptical arc middle = %v, want (0, 1)", mid)
	}
	// a clamped B-spline runs from its first control point to its last
	ctrl := []pt{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {2, 0}}
	knots := clampKnots([]float64{0, 0, 0.5, 1, 1, 1}, len(ctrl), 3)
	if len(knots) != 9 || knots[3] != 0 || knots[4] != 0.5 || knots[5] != 1 {
		t.Errorf("knots = %v", knots)
	}
	pts := splinePoints(ctrl, nil, knots, 3)
	if last := pts[len(pts)-1]; last.dist(pt{2, 0}) > 1e-6 {
		t.Errorf("spline ends at %v", last)
	}
	// rounding keeps a figure closed
	sp := &subpath{start: pt{0, 0}, segs: []seg{{p: pt{10, 0}}, {p: pt{10, 10}}, {p: pt{0, 10}}, {p: pt{0, 0}}}}
	roundCorners(sp, 2)
	if !sp.closed() || len(sp.segs) != 8 {
		t.Errorf("rounded square: closed %v, %d segments", sp.closed(), len(sp.segs))
	}
	for i := 1; i < len(arrowDefs); i++ {
		if arrowDefs[i].draw == nil {
			t.Errorf("arrowhead %d is not defined", i)
		}
	}
}

// bezierPoint evaluates the figure made of segments at the middle of its
// middle segment (enough for the symmetric arcs above).
func bezierPoint(segs []seg, _ float64) pt {
	if len(segs)%2 == 1 {
		s := segs[len(segs)/2]
		return s.p
	}
	return segs[len(segs)/2-1].p
}

func TestEvalSize(t *testing.T) {
	for _, c := range []struct {
		f    string
		want float64
		ok   bool
	}{
		{"Width*0.5", 2, true}, {"Height", 3, true}, {"(Width+Height)/5", 1.4, true}, {"-Width*1", -4, true},
		{"72 pt", 1, true}, {"25.4 mm+1", 2, true}, {"1E-1*Width", 0.4, true}, {"Inh", 0, false},
		{"Width*User.Ratio", 0, false}, {"GUARD(Width)", 0, false}, {"Width/0", 0, false}, {"", 0, false},
	} {
		got, ok := evalSize(c.f, 4, 3)
		if ok != c.ok || math.Abs(got-c.want) > 1e-9 {
			t.Errorf("evalSize(%q) = %g, %v; want %g, %v", c.f, got, ok, c.want, c.ok)
		}
	}
}

func TestDeterministic(t *testing.T) {
	for _, name := range []string{"shapes.vsdx", "flow.vdx"} {
		var out [2][]byte
		for i := range out {
			res, err := ConvertFile(filepath.Join("testdata", name), testOptions())
			if err != nil {
				t.Fatal(err)
			}
			var buf bytes.Buffer
			if err := res.Doc.WriteSingle(&buf); err != nil {
				t.Fatal(err)
			}
			out[i] = buf.Bytes()
		}
		if !bytes.Equal(out[0], out[1]) {
			t.Errorf("%s: two conversions differ", name)
		}
	}
}

func TestPagesOption(t *testing.T) {
	opts := testOptions()
	opts.Pages = []int{2}
	res, err := ConvertFile(filepath.Join("testdata", "shapes.vsdx"), opts)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(res.Doc.Views[0].Pages); n != 1 || res.Pages != 1 {
		t.Errorf("%d pages", n)
	}
	opts.Pages = []int{3}
	if _, err := ConvertFile(filepath.Join("testdata", "shapes.vsdx"), opts); err == nil {
		t.Error("page 3 of 2 converted")
	}
}

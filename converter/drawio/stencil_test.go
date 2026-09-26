package drawio

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestStencilBasename(t *testing.T) {
	for _, tc := range []struct{ name, want string }{
		{"mxgraph.flowchart.decision", "flowchart"},
		{"mxgraph.cisco.routers.router", "cisco/routers"},
		{"mxgraph.pid.compressors_-_iso.compressor_(centrifugal)", "pid/compressors_-_iso"},
		{"mxgraph.flowchart", "flowchart"},
		{"mxgraph", ""},
		{"rectangle", ""},
		{"MxGraph.flowchart.decision", ""},
		{"", ""},
	} {
		if got := stencilBasename(tc.name); got != tc.want {
			t.Errorf("stencilBasename(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestLookupStencil(t *testing.T) {
	for _, tc := range []struct {
		name   string
		w0, h0 float64
		aspect string
	}{
		{"mxgraph.flowchart.decision", 100, 100, "variable"}, // no w and h: 100
		{"mxgraph.flowchart.multi-document", 88, 60.28, "variable"},
		{"mxgraph.flowchart.predefined_process", 98, 60, "variable"}, // "Predefined Process"
		{"mxgraph.basic.4_point_star", 92, 92, "variable"},
		{"mxgraph.arrows.arrow_right", 97.5, 70, "variable"},
		{"mxgraph.networks.cloud", 92.1, 48.51, "variable"},
		{"mxgraph.bpmn.error_end", 97, 97, "fixed"},
		{"mxgraph.cisco.routers.router", 0, 0, ""},
		{"mxgraph.rack.general.1u_rack_server", 0, 0, ""},                    // <shapes name="mxgraph.rack.General">
		{"mxgraph.pid.compressors_-_iso.compressor_(centrifugal)", 0, 0, ""}, // pid/compressors_iso.xml
		{"mxgraph.electrical.radio.aerial_-_antenna_1", 80, 100, "variable"},
		{"mxgraph.electrical.logic_gates.and", 100, 60, "variable"},
		{"mxgraph.azure.computer", 0, 0, ""},
	} {
		st := lookupStencil(tc.name)
		if st == nil {
			t.Errorf("%s: not found", tc.name)
			continue
		}
		if tc.aspect != "" && (st.w0 != tc.w0 || st.h0 != tc.h0 || st.aspect != tc.aspect) {
			t.Errorf("%s: w0, h0, aspect = %g, %g, %q, want %g, %g, %q", tc.name, st.w0, st.h0, st.aspect, tc.w0, tc.h0, tc.aspect)
		}
		if st.fg == nil {
			t.Errorf("%s: no foreground", tc.name)
		}
	}
	for _, name := range []string{
		"mxgraph.flowchart.nonexistent",
		"mxgraph.flowchart.Decision",   // registry names are lowercase
		"mxgraph.arrows2.arrow",        // a JavaScript shape
		"mxgraph.aws4.resourceIcon",    // a JavaScript shape (mxAWS4.js)
		"mxgraph.basic.rect",           // a JavaScript shape of a stencil library
		"mxgraph.nosuchlibrary.x",      // no such file
		"mxgraph.../../stencils/x.y.z", // not a file name
		"rectangle", "", "mxgraph",
	} {
		if st := lookupStencil(name); st != nil {
			t.Errorf("%s: found %q", name, st.name)
		}
	}
	first, second := lookupStencil("mxgraph.flowchart.decision"), lookupStencil("mxgraph.flowchart.decision")
	if first == nil || first != second {
		t.Error("lookups return different stencils")
	}
}

// TestStencilLazyLoading checks that library files are parsed on first
// use only: waveforms is used by no other test.
func TestStencilLazyLoading(t *testing.T) {
	loaded := func(file string) bool {
		stencilFilesMu.Lock()
		defer stencilFilesMu.Unlock()
		f := stencilFiles[file]
		return f != nil && f.shapes != nil
	}
	const file = "electrical/waveforms.xml"
	if loaded(file) {
		t.Skipf("%s is loaded already (an earlier run of the test)", file)
	}
	lookupStencil("mxgraph.flowchart.decision")
	if loaded(file) {
		t.Fatalf("%s loaded by another library", file)
	}
	var wg sync.WaitGroup
	res := make([]*stencil, 8)
	for i := range res {
		wg.Go(func() { res[i] = lookupStencil("mxgraph.electrical.waveforms.sine_wave") })
	}
	wg.Wait()
	if !loaded(file) || res[0] == nil {
		t.Fatalf("%s not loaded after use", file)
	}
	for _, r := range res {
		if r != res[0] {
			t.Fatal("concurrent lookups return different stencils")
		}
	}
}

// compressGraph encodes a stencil as draw.io's Graph.compress does.
func compressGraph(t *testing.T, xmlText string) string {
	t.Helper()
	var buf bytes.Buffer
	zw, _ := flate.NewWriter(&buf, flate.BestCompression)
	zw.Write([]byte(strings.ReplaceAll(url.QueryEscape(xmlText), "+", "%20")))
	zw.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func testConverter() *converter {
	return &converter{opts: &Options{}, doc: bdf.NewDocument(), warned: map[string]bool{}, images: map[string]*imageRef{}}
}

func TestInlineStencil(t *testing.T) {
	c := testConverter()
	name := "stencil(" + compressGraph(t, `<shape name="T" w="40" h="20" aspect="fixed" strokewidth="2"><connections/><background><rect x="0" y="0" w="40" h="20"/></background><foreground><fillstroke/></foreground></shape>`) + ")"
	st := c.stencil(name)
	if st == nil {
		t.Fatalf("inline stencil not decoded (warnings %q)", c.warnings)
	}
	if st.w0 != 40 || st.h0 != 20 || st.aspect != "fixed" || st.strokewidth != "2" || st.name != "T" {
		t.Errorf("got w0 %g h0 %g aspect %q strokewidth %q name %q", st.w0, st.h0, st.aspect, st.strokewidth, st.name)
	}
	if c.stencil(name) != st {
		t.Error("inline stencil is not cached")
	}
	// draw.io's output: padding may be missing, line breaks may be present
	trimmed := strings.TrimRight(name[:len(name)-1], "=")
	if c.stencil(trimmed[:20]+"\n"+trimmed[20:]+")") == nil {
		t.Error("unpadded inline stencil not decoded")
	}
	if c.stencil("stencil(not base64!)") != nil || len(c.warnings) != 1 {
		t.Errorf("broken inline stencil: warnings %q", c.warnings)
	}
}

// paintStencil paints a stencil shape of the given style in b and
// returns the decoded object.
func paintStencil(t *testing.T, c *converter, styleStr string, b rect) (*bdf.ObjectPart, *c2d) {
	t.Helper()
	s := &shape{conv: c, st: &cellState{cell: &cell{styleStr: styleStr}}, style: parseStyle(styleStr, false)}
	s.name = s.style.get("shape", "")
	s.stencil = c.stencil(s.name)
	if s.stencil == nil {
		t.Fatalf("no stencil %q (warnings %q)", s.name, c.warnings)
	}
	s.def = &shapeDef{}
	s.apply()
	s.bounds = b
	cv := c.newCanvas()
	c2 := newC2D(cv)
	s.paint(c2)
	part, err := bdf.DecodeObject(cv.obj.Encode())
	if err != nil {
		t.Fatal(err)
	}
	return part, c2
}

// paintedPath returns a path and the transform in effect where it is
// first filled or stroked.
func paintedPath(t *testing.T, part *bdf.ObjectPart, idx int) (*bdf.Path, matrix) {
	t.Helper()
	ins, err := part.Instructions()
	if err != nil {
		t.Fatal(err)
	}
	m := identity
	var stack []matrix
	for _, in := range ins {
		switch in.Op {
		case bdf.OpSave:
			stack = append(stack, m)
		case bdf.OpRestore:
			m, stack = stack[len(stack)-1], stack[:len(stack)-1]
		case bdf.OpTransform:
			var a matrix
			for i := range a {
				a[i] = float64(in.Args[i].(float32))
			}
			m = m.mul(a)
		case bdf.OpTranslate:
			m = m.mul(translateM(float64(in.Args[0].(float32)), float64(in.Args[1].(float32))))
		case bdf.OpFillPath, bdf.OpStrokePath:
			if int(in.Args[0].(uint64)) == idx {
				return part.Paths[idx].Inline, m
			}
		}
	}
	t.Fatalf("path %d is not painted", idx)
	return nil, m
}

// pathBox returns the bounds of the points of a path (and of the boxes of
// its rectangles and ellipses) under m, rounded to 1/1000, and its first
// point.
func pathBox(p *bdf.Path, m matrix) (r rect, first point) {
	var e extent
	defer func() { r = e.r }()
	add := func(x, y float32) {
		tx, ty := m.apply(float64(x), float64(y))
		pt := point{math.Round(tx*1000) / 1000, math.Round(ty*1000) / 1000}
		if !e.ok {
			first = pt
		}
		e.add(pt)
	}
	ai := 0
	for _, v := range p.Verbs {
		a := p.Args[ai:]
		switch v {
		case bdf.VerbMove, bdf.VerbLine:
			add(a[0], a[1])
			ai += 2
		case bdf.VerbQuad:
			add(a[0], a[1])
			add(a[2], a[3])
			ai += 4
		case bdf.VerbCubic:
			add(a[0], a[1])
			add(a[2], a[3])
			add(a[4], a[5])
			ai += 6
		case bdf.VerbRect, bdf.VerbRoundRect:
			add(a[0], a[1])
			add(a[0]+a[2], a[1]+a[3])
			ai += 4
			if v == bdf.VerbRoundRect {
				ai++
			}
		case bdf.VerbEllipse:
			add(a[0]-a[2], a[1]-a[3])
			add(a[0]+a[2], a[1]+a[3])
			ai += 8
		case bdf.VerbArcTo:
			ai += 5
		}
	}
	return r, first
}

func TestStencilDraw(t *testing.T) {
	c := testConverter()
	part, _ := paintStencil(t, c, "shape=mxgraph.flowchart.decision;fillColor=#ff0000;", rect{10, 20, 100, 60})
	if len(part.Paths) != 1 {
		t.Fatalf("got %d paths, want 1", len(part.Paths))
	}
	p := part.Paths[0].Inline
	wantVerbs := []byte{bdf.VerbMove, bdf.VerbLine, bdf.VerbLine, bdf.VerbLine, bdf.VerbClose}
	wantArgs := []float32{60, 20, 110, 50, 60, 80, 10, 50}
	if !bytes.Equal(p.Verbs, wantVerbs) || len(p.Args) != len(wantArgs) {
		t.Fatalf("path %v %v", p.Verbs, p.Args)
	}
	for i := range wantArgs {
		if p.Args[i] != wantArgs[i] {
			t.Fatalf("path args %v, want %v", p.Args, wantArgs)
		}
	}
	ins, _ := part.Instructions()
	var fills, strokes int
	for _, in := range ins {
		switch in.Op {
		case bdf.OpFillPath:
			fills++
		case bdf.OpStrokePath:
			strokes++
		}
	}
	if fills != 1 || strokes != 1 {
		t.Errorf("%d fills and %d strokes, want 1 and 1", fills, strokes)
	}

	// Directions, flips and rotations keep the shape in its bounds. The
	// card (98 x 60) starts at (19, 0), after its cut corner.
	for _, tc := range []struct {
		style string
		first point
	}{
		{"", point{29.388, 20}},
		{"flipH=1;", point{90.612, 20}},
		{"stencilFlipH=1;", point{90.612, 20}},
		{"flipV=1;", point{29.388, 80}},
		{"direction=south;", point{110, 31.633}},
		{"direction=north;", point{10, 68.367}},
		{"direction=west;", point{90.612, 80}},
		{"direction=south;flipV=1;", point{110, 68.367}},
	} {
		part, _ := paintStencil(t, c, "shape=mxgraph.flowchart.card;"+tc.style, rect{10, 20, 100, 60})
		box, first := pathBox(paintedPath(t, part, 0))
		if box != (rect{10, 20, 100, 60}) || first != tc.first {
			t.Errorf("%s: bounds %v, first point %v; want %v", tc.style, box, first, tc.first)
		}
	}
	// fixed aspect: centered with the stencil's proportions
	for _, dir := range []string{"", "direction=south;"} {
		part, _ = paintStencil(t, c, "shape=mxgraph.bpmn.error_end;"+dir, rect{0, 0, 200, 100})
		if box, _ := pathBox(paintedPath(t, part, 0)); box != (rect{50, 0, 100, 100}) {
			t.Errorf("fixed aspect %s: bounds %v", dir, box)
		}
	}
}

func TestComputeAspect(t *testing.T) {
	st := &stencil{w0: 100, h0: 50, aspect: "variable"}
	for _, tc := range []struct {
		dir  string
		want rect
	}{
		{"", rect{10, 20, 2, 2}},
		{"east", rect{10, 20, 2, 2}},
		{"south", rect{60, -30, 1, 4}},
		{"north", rect{60, -30, 1, 4}},
	} {
		if got := st.computeAspect(10, 20, 200, 100, tc.dir); got != tc.want {
			t.Errorf("variable %q: %v, want %v", tc.dir, got, tc.want)
		}
	}
	st.aspect = "fixed"
	if got := st.computeAspect(0, 0, 300, 100, ""); got != (rect{50, 0, 2, 2}) {
		t.Errorf("fixed: %v", got)
	}
	if got := st.computeAspect(0, 0, 300, 100, "south"); got != (rect{100, 25, 1, 1}) {
		t.Errorf("fixed south: %v", got)
	}
}

func TestStencilStrokeWidth(t *testing.T) {
	c := testConverter()
	lineWidths := func(styleStr string, b rect) []float32 {
		part, _ := paintStencil(t, c, styleStr, b)
		ins, _ := part.Instructions()
		var ws []float32
		for _, in := range ins {
			if in.Op == bdf.OpLine {
				ws = append(ws, in.Args[0].(float32))
			}
		}
		return ws
	}
	// strokewidth="inherit": the style's width as it is
	if ws := lineWidths("shape=mxgraph.flowchart.decision;strokeWidth=3;", rect{0, 0, 200, 50}); len(ws) != 1 || ws[0] != 3 {
		t.Errorf("inherit: line widths %v", ws)
	}
	// strokewidth="2": scaled by the smaller scale (50/48.51 for the cloud)
	ws := lineWidths("shape=mxgraph.networks.cloud;strokeWidth=5;", rect{0, 0, 184.2, 48.51})
	if len(ws) == 0 || math.Abs(float64(ws[0])-2) > 1e-4 {
		t.Errorf("numeric: line widths %v", ws)
	}
	// strokewidth elements: scaled, or fixed
	name := "stencil(" + compressGraph(t, `<shape w="10" h="10" strokewidth="inherit"><foreground>`+
		`<strokewidth width="2"/><rect x="0" y="0" w="10" h="10"/><stroke/>`+
		`<strokewidth width="3" fixed="1"/><rect x="0" y="0" w="10" h="10"/><stroke/>`+
		`</foreground></shape>`) + ")"
	if ws := lineWidths("shape="+name+";", rect{0, 0, 40, 20}); len(ws) != 2 || ws[0] != 4 || ws[1] != 3 {
		t.Errorf("strokewidth elements: line widths %v", ws)
	}
}

func TestStencilSaveRestore(t *testing.T) {
	c := testConverter()
	for _, body := range []string{
		`<save/><save/><rect x="0" y="0" w="10" h="10"/><fill/>`,
		`<save/><rect x="0" y="0" w="10" h="10"/><fill/><restore/><restore/><restore/>`,
	} {
		name := "stencil(" + compressGraph(t, `<shape w="10" h="10"><foreground>`+body+`</foreground></shape>`) + ")"
		part, c2 := paintStencil(t, c, "shape="+name+";", rect{0, 0, 10, 10})
		if len(c2.stack) != 0 {
			t.Errorf("%s: canvas stack depth %d", body, len(c2.stack))
		}
		ins, _ := part.Instructions()
		depth := 0
		for _, in := range ins {
			switch in.Op {
			case bdf.OpSave:
				depth++
			case bdf.OpRestore:
				depth--
			}
			if depth < 0 {
				break
			}
		}
		if depth != 0 {
			t.Errorf("%s: unbalanced SAVE/RESTORE (%d)", body, depth)
		}
	}
}

func TestStencilIncludeShape(t *testing.T) {
	c := testConverter()
	name := "stencil(" + compressGraph(t, `<shape w="100" h="100"><foreground>`+
		`<include-shape name="mxgraph.flowchart.decision" x="50" y="0" w="50" h="50"/>`+
		`<include-shape name="mxgraph.nosuch.shape" x="0" y="0" w="50" h="50"/>`+
		`</foreground></shape>`) + ")"
	part, _ := paintStencil(t, c, "shape="+name+";", rect{0, 0, 200, 100})
	if len(part.Paths) != 1 {
		t.Fatalf("got %d paths, want 1", len(part.Paths))
	}
	if got, _ := pathBox(paintedPath(t, part, 0)); got != (rect{100, 0, 100, 50}) {
		t.Errorf("included shape bounds %v", got)
	}
}

func TestStencilParseColor(t *testing.T) {
	st := &stencil{}
	node := func(name, attrs string) *stencilNode {
		root, err := parseStencilXML([]byte("<" + name + " " + attrs + "/>"))
		if err != nil {
			t.Fatal(err)
		}
		return root
	}
	shapeOf := func(styleStr string) *shape {
		s := &shape{conv: testConverter(), st: &cellState{cell: &cell{styleStr: styleStr}}, style: parseStyle(styleStr, false)}
		s.apply()
		return s
	}
	for _, tc := range []struct {
		style, elem, attrs, want string
	}{
		{"", "fillcolor", `color="#123456"`, "#123456"},
		{"", "fillcolor", `color="red"`, "red"},
		{"", "fillcolor", `color="none"`, ""},
		{"", "fillcolor", `color="default"`, "#ffffff"},
		{"", "strokecolor", `color="default"`, "#000000"},
		{"fillColor=#aabbcc;strokeColor=#ddeeff;", "strokecolor", `color="fill"`, "#aabbcc"},
		{"fillColor=#aabbcc;strokeColor=#ddeeff;", "fillcolor", `color="stroke"`, "#ddeeff"},
		{"strokeColor=none;", "fillcolor", `color="stroke"`, ""},
		{"", "fillcolor", `color="strokeColor2" default="#ffffff"`, "#ffffff"},
		{"strokeColor2=#ff0000;", "fillcolor", `color="strokeColor2" default="#ffffff"`, "#ff0000"},
		{"strokeColor2=none;", "fillcolor", `color="strokeColor2" default="#ffffff"`, ""},
		{"strokeColor2=default;", "fillcolor", `color="strokeColor2" default="#00ff00"`, "#00ff00"},
		{"strokeColor2=default;", "strokecolor", `color="strokeColor2"`, "#000000"},
		{"", "fillcolor", `color="strokeColor2"`, ""},
		{"text;", "fillcolor", `color="fillColor" default="#ffffff"`, ""}, // none from a named style
	} {
		if got := st.parseColor(shapeOf(tc.style), node(tc.elem, tc.attrs)); got != tc.want {
			t.Errorf("%s <%s %s>: %q, want %q", tc.style, tc.elem, tc.attrs, got, tc.want)
		}
	}
}

func TestStencilLabelMargins(t *testing.T) {
	st := lookupStencil("mxgraph.flowchart.multi-document")
	s := &shape{style: parseStyle("boundedLbl=1;", false)}
	m := st.labelMargins(s, rect{100, 100, 176, 120.56})
	if m == nil {
		t.Fatal("no label margins")
	}
	want := rect{0, 20, 20, 6.56}
	if math.Abs(m.x-want.x) > 1e-9 || math.Abs(m.y-want.y) > 1e-9 || math.Abs(m.w-want.w) > 1e-9 || math.Abs(m.h-want.h) > 1e-9 {
		t.Errorf("label margins %v, want %v", *m, want)
	}
	if m := st.labelMargins(&shape{style: parseStyle("", false)}, rect{0, 0, 88, 60}); m != nil {
		t.Errorf("label margins without boundedLbl: %v", *m)
	}
	if m := lookupStencil("mxgraph.flowchart.decision").labelMargins(s, rect{0, 0, 88, 60}); m != nil {
		t.Errorf("label margins of a stencil without labelBounds: %v", *m)
	}
}

func TestJSNumber(t *testing.T) {
	for _, tc := range []struct {
		s    string
		want float64
	}{
		{"", 0}, {"  ", 0}, {"12", 12}, {" -1.5 ", -1.5}, {".5", 0.5}, {"5.", 5}, {"1e3", 1000},
		{"0x10", 16}, {"Infinity", math.Inf(1)},
	} {
		if got := numberJS(tc.s); got != tc.want {
			t.Errorf("numberJS(%q) = %g, want %g", tc.s, got, tc.want)
		}
	}
	for _, s := range []string{"12px", "abc", "1e", "inf", "NaN", "0x", "1_000"} {
		if got := numberJS(s); !math.IsNaN(got) {
			t.Errorf("numberJS(%q) = %g, want NaN", s, got)
		}
	}
}

// TestStencilTestdata converts the stencil test diagrams: every shape in
// them must resolve.
func TestStencilTestdata(t *testing.T) {
	files, _ := filepath.Glob("testdata/stencil/*.drawio")
	if len(files) == 0 {
		t.Skip("no test diagrams")
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		res, err := Convert(data, &Options{NoSystemFonts: true, NoTextIndex: true})
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		for _, w := range res.Warnings {
			if strings.Contains(w, "not supported") || strings.Contains(w, "cannot be read") {
				t.Errorf("%s: %s", f, w)
			}
		}
	}
}

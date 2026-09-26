package drawio

import (
	"fmt"
	"html"
	"slices"
	"strings"
	"testing"
)

// basicShapeNames are the shapes mxBasic.js registers.
var basicShapeNames = []string{
	"cross2", "rectCallout", "roundRectCallout", "wave2", "octagon2", "isocube", "acute_triangle",
	"obtuse_triangle", "drop", "cone2", "pyramid", "4_point_star_2", "diag_snip_rect", "diag_round_rect",
	"corner_round_rect", "plaque", "frame", "plaque_frame", "rounded_frame", "frame_corner", "diag_stripe",
	"donut", "layered_rect", "button", "shaded_button", "pie", "arc", "partConcEllipse", "numberedEntryVert",
	"bendingArch", "three_corner_round_rect", "polygon", "patternFillRect",
}

// TestBasicShapesResolve checks that the shapes of mxBasic.js resolve to
// their registered painters, not to stencils or the rectangle fallback,
// and convert without warnings.
func TestBasicShapesResolve(t *testing.T) {
	var b strings.Builder
	b.WriteString(`<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>`)
	for i, n := range basicShapeNames {
		fmt.Fprintf(&b, `<mxCell id="v%d" style="shape=mxgraph.basic.%s;" vertex="1" parent="1"><mxGeometry x="%d" y="0" width="80" height="60" as="geometry"/></mxCell>`, i, n, i*100)
	}
	b.WriteString(`</root></mxGraphModel>`)
	root, err := parseXML([]byte(b.String()))
	if err != nil {
		t.Fatal(err)
	}
	m := parseModel(root)
	v := newView(m, nil, nil)
	c := &converter{warned: map[string]bool{}}
	for i, n := range basicShapeNames {
		name := "mxgraph.basic." + n
		d := shapeRegistry[name]
		if d == nil {
			t.Errorf("%s is not registered", name)
			continue
		}
		s := c.newShape(v.state(m.cells[fmt.Sprint("v", i)]))
		if s.def != d || s.stencil != nil {
			t.Errorf("%s: resolved to another shape", name)
		}
	}
	if len(c.warned) > 0 {
		t.Errorf("warnings: %v", c.warned)
	}
	res, err := Convert([]byte(b.String()), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) > 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	// the markers mxShapeBasicArc offers
	for _, mk := range []string{"classic", "classicThin", "open", "openThin", "block", "blockThin", "oval", "diamond", "diamondThin"} {
		if markers[mk] == nil {
			t.Errorf("marker %q is not registered", mk)
		}
	}
}

// basicStyles are parameter values of the basic shapes: none, zero,
// huge, negative and not numbers, and patterns that never end in
// JavaScript.
var basicStyles = []string{
	"",
	"dx=0;dy=0;size=0;dx1=0;dx2=0;dy1=0;dy2=0;isoAngle=0;startAngle=0;endAngle=0;arcWidth=0;step=0;fillStyle=diagGrid;",
	"dx=1e308;dy=1e308;size=1e308;dx1=1e308;dx2=-1e308;dy1=1e308;dy2=-1e308;isoAngle=1e308;startAngle=1e308;endAngle=-1e308;arcWidth=1e308;step=1e-300;fillStyle=grid;",
	"dx=-5;dy=-5;size=-5;startAngle=0.5;endAngle=0;arcWidth=-1;step=-5;fillStyle=diag;boundedLbl=1;",
	"dx=x;dy=y;size=z;startAngle=a;endAngle=b;arcWidth=c;step=d;fillStyle=diagRev;fillStrokeWidth=e;",
	"startAngle=0;endAngle=0.5;fillStyle=vert;step=1;top=0;right=0;bottom=0;left=0;",
	"startAngle=0.25;endAngle=0.75;fillStyle=hor;step=2;strokeColor=none;fillStrokeColor=none;",
	"startArrow=classic;endArrow=block;startSize=1e308;endSize=-5;startFill=0;endFillColor=#ff0000;",
	"startArrow=nosuchmarker;endArrow=oval;strokeWidth=1e308;",
	"polyCoords=[[0,0],[1,0],[0.5,1]];polyCurves=[[\"Q\",1,1],[\"Q\",0,0],[\"Q\",0.5,2]];",
	"polyCoords=[[0,0],[\"a\",1],null];polyline=1;",
	"polyCoords=[[1e308,-1e308],[1e308,1e308]];polyCurves=[[\"Q\",\"x\",1]];polyline=false;",
	"polyCoords=5;polyCurves={};",
	"polyCoords=[[0,0],[1,1]];polyCurves=[[\"Q\",null,\"0.5\"],\"Qab\",[]];",
	"polyCoords=[[0,0];polyCurves=[;",
	"rotation=45;direction=south;flipH=1;shadow=1;dashed=1;gradientColor=#ff0000;",
}

// TestBasicShapesPaint paints every basic shape with the parameter values
// and degenerate sizes of the shape tests; an internal error (a recovered
// panic) fails the test, a loop without end hangs it.
func TestBasicShapesPaint(t *testing.T) {
	var names []string
	for name := range shapeRegistry {
		if strings.HasPrefix(name, "mxgraph.basic.") {
			names = append(names, name)
		}
	}
	slices.Sort(names)
	var b strings.Builder
	b.WriteString("<mxfile>")
	id := 0
	for _, name := range names {
		fmt.Fprintf(&b, `<diagram id="%s" name="%s"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>`, name, name)
		for _, st := range basicStyles {
			for _, sz := range shapeSizes {
				id++
				style := html.EscapeString("shape=" + name + ";" + st)
				fmt.Fprintf(&b, `<mxCell id="v%d" value="x" style="%s" vertex="1" parent="1"><mxGeometry x="10" y="10" width="%g" height="%g" as="geometry"/></mxCell>`, id, style, sz[0], sz[1])
			}
		}
		b.WriteString(`</root></mxGraphModel></diagram>`)
	}
	b.WriteString("</mxfile>")
	res, err := Convert([]byte(b.String()), &Options{NoSystemFonts: true, NoTextIndex: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages != len(names) {
		t.Errorf("%d pages, want %d", res.Pages, len(names))
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "internal error") {
			t.Error(w)
		}
	}
}

// paintTestShape paints a shape of w by h and reports whether anything
// was drawn.
func paintTestShape(t *testing.T, name, st string, w, h float64) (*shape, bool) {
	t.Helper()
	s := testShape(t, name, st)
	s.conv.opts = &Options{}
	s.bounds = rect{0, 0, w, h}
	cv := s.conv.newCanvas()
	s.paint(newC2D(cv))
	return s, cv.drawn
}

// TestBasicPolygon checks which polyCoords and polyCurves draw: draw.io
// draws nothing where JavaScript throws, and nothing is drawn here where
// it would compute NaN.
func TestBasicPolygon(t *testing.T) {
	for _, tc := range []struct {
		style string
		drawn bool
	}{
		{"", false},
		{"polyCoords=[];", false},
		{"polyCoords=[[0,0],[1,0],[1,1]];", true},
		{"polyCoords=[[0,0],[1,0],[1,1]];polyline=1;", true},
		{`polyCoords=[[0,0],[1,0],[1,1]];polyCurves=[["Q",0.5,-0.5],null,["Q",0,1]];`, true},
		{`polyCoords=[[0,0],[1,0],[1,1]];polyCurves=[["Q",null,"0.5"]];`, true},
		{`polyCoords=[[0,0],[1,0],[1,1]];polyCurves={"a":1};`, true},
		{`polyCoords=[[0,0],[1,0],[1,1]];polyCurves=[["Q","x",1]];`, false},
		{`polyCoords=[[0,0],[1,0],[1,1]];polyCurves=[[],[],["Q","x",1]];polyline=1;`, true}, // the closing curve is not read
		{`polyCoords=[[0,0],[1,0],[1,1]];polyCurves=[[],[],["Q","x",1]];`, false},
		{`polyCoords=[[0,0],["a",1]];`, false},
		{`polyCoords=[[0,0],null];`, false},
		{`polyCoords=[[0,0],[1]];`, false},
		{`polyCoords=[[0,0],[true,"1"]];`, true},
		{"polyCoords=5;", false},
		{"polyCoords=[[0,0];", false},
		{"polyCoords=[[0,0],[1,1]];polyCurves=[;", false},
	} {
		if _, drawn := paintTestShape(t, "mxgraph.basic.polygon", tc.style, 100, 80); drawn != tc.drawn {
			t.Errorf("%q: drawn %v, want %v", tc.style, drawn, tc.drawn)
		}
	}
}

// TestBasicArcMarkers checks that the arc draws the markers draw.io knows
// and warns about others, which draw.io does not draw.
func TestBasicArcMarkers(t *testing.T) {
	s, drawn := paintTestShape(t, "mxgraph.basic.arc", "startArrow=diamond;endArrow=oval;", 100, 80)
	if !drawn || len(s.conv.warned) > 0 {
		t.Errorf("known markers: drawn %v, warnings %v", drawn, s.conv.warned)
	}
	s, drawn = paintTestShape(t, "mxgraph.basic.arc", "startArrow=nosuchmarker;", 100, 80)
	if !drawn || !s.conv.warned["marker:nosuchmarker"] {
		t.Errorf("unknown marker: drawn %v, warnings %v", drawn, s.conv.warned)
	}
}

// TestBasicPatternFillSteps checks that pattern fills whose loops would
// not end in JavaScript still paint the rectangle.
func TestBasicPatternFillSteps(t *testing.T) {
	for _, st := range []string{"step=0;fillStyle=diagGrid;", "step=-1;fillStyle=grid;", "step=1e-300;fillStyle=diagRev;"} {
		if _, drawn := paintTestShape(t, "mxgraph.basic.patternFillRect", st, 1e9, 1e9); !drawn {
			t.Errorf("%q: nothing drawn", st)
		}
	}
	if !patternSteps(100, 5) || patternSteps(100, 0) || patternSteps(100, -5) || patternSteps(1e9, 1) {
		t.Error("patternSteps")
	}
}

// TestBasicLabelMargins checks the label margins of the callouts and the
// layered rectangle against values worked out from the JavaScript.
func TestBasicLabelMargins(t *testing.T) {
	r := rect{0, 0, 120, 100}
	for _, tc := range []struct {
		name, style string
		want        *rect
	}{
		{"rectCallout", "", nil},
		{"rectCallout", "boundedLbl=1;", &rect{0, 0, 0, 0.5}},
		{"rectCallout", "boundedLbl=1;dy=15;", &rect{0, 0, 0, 15}},
		{"rectCallout", "boundedLbl=1;dy=500;", &rect{0, 0, 0, 500}}, // not clamped
		{"roundRectCallout", "boundedLbl=1;dy=20;", &rect{0, 0, 0, 20}},
		{"roundRectCallout", "boundedLbl=0;dy=20;", nil},
		{"layered_rect", "", nil},
		{"layered_rect", "boundedLbl=1;", &rect{0, 0, 0.5, 0.5}},
		{"layered_rect", "boundedLbl=1;dx=10;", &rect{0, 0, 10, 10}},
		{"layered_rect", "boundedLbl=1;dx=100;", &rect{0, 0, 50, 50}},
		{"layered_rect", "boundedLbl=1;dx=-3;", &rect{0, 0, 0, 0}},
	} {
		s := testShape(t, "mxgraph.basic."+tc.name, tc.style)
		got := s.def.labelMargins(s, r)
		if (got == nil) != (tc.want == nil) || got != nil && !rectNear(*got, *tc.want) {
			t.Errorf("%s %q: margins %v, want %v", tc.name, tc.style, got, tc.want)
		}
	}
	// turned for the direction and flips (mxUtils.getDirectedBounds)
	r = rect{10, 20, 120, 60}
	for _, tc := range []struct {
		name, style string
		want        rect
	}{
		{"rectCallout", "boundedLbl=1;dy=15;", rect{10, 20, 120, 45}},
		{"rectCallout", "boundedLbl=1;dy=15;direction=south;", rect{25, 20, 105, 60}},
		{"rectCallout", "boundedLbl=1;dy=15;flipV=1;", rect{10, 35, 120, 45}},
		{"layered_rect", "boundedLbl=1;dx=10;", rect{10, 20, 110, 50}},
		{"layered_rect", "boundedLbl=1;dx=10;flipH=1;", rect{20, 20, 110, 50}},
	} {
		s := testShape(t, "mxgraph.basic."+tc.name, tc.style)
		if got := shapeLabelBounds(s, r); !rectNear(got, tc.want) {
			t.Errorf("%s %q: label bounds %v, want %v", tc.name, tc.style, got, tc.want)
		}
	}
}

// TestJSONNumber checks the conversion of JSON values to numbers as
// JavaScript's arithmetic does it.
func TestJSONNumber(t *testing.T) {
	for _, tc := range []struct {
		v    any
		want float64
		ok   bool
	}{
		{1.5, 1.5, true}, {nil, 0, true}, {true, 1, true}, {false, 0, true},
		{"0.25", 0.25, true}, {" 2 ", 2, true}, {"", 0, true}, {"1px", 0, false},
		{"x", 0, false}, {map[string]any{}, 0, false},
		{[]any{}, 0, true}, {[]any{3.0}, 3, true}, {[]any{[]any{"4"}}, 4, true}, {[]any{nil}, 0, true},
		{[]any{true}, 0, false}, {[]any{1.0, 2.0}, 0, false},
	} {
		got, ok := jsonNumber(tc.v)
		if ok != tc.ok || ok && got != tc.want {
			t.Errorf("jsonNumber(%#v) = %v, %v; want %v, %v", tc.v, got, ok, tc.want, tc.ok)
		}
	}
}

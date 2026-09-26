package drawio

import (
	"fmt"
	"html"
	"math"
	"slices"
	"strings"
	"testing"
)

// shapeStyles are style variants every shape is painted with.
var shapeStyles = []string{
	"",
	"rounded=1;dashed=1;glass=1;shadow=1;",
	"direction=south;flipH=1;rotation=45;",
	"direction=north;flipV=1;fixedSize=1;",
	"direction=west;gradientColor=#ff0000;strokeWidth=5;",
	"size=1e308;dx=-1e308;dy=1e308;arcSize=1e308;absoluteArcSize=1;startSize=1e308;tabWidth=-5;tabHeight=1e308;",
	"size=-5;fixedSize=1;startSize=-10;footerSize=1e9;tipX=1e308;tipY=-1e308;base=-3;width=-1;height=1e308;",
	"fillColor=none;strokeColor=none;double=1;boundedLbl=1;lid=0;horizontal=0;swimlaneLine=0;separatorColor=#f00;",
	"participant=umlActor;lifelineMirror=1;umlStateConnection=connPointRefExit;symbol0=cube;symbol0Width=10;symbol0Height=10;symbol1=ext;",
}

// shapeSizes are degenerate and huge sizes (width, height).
var shapeSizes = [][2]float64{{0, 0}, {-10, -5}, {1, 1}, {0, 50}, {80, 0}, {120, 60}, {1e9, 1e9}, {1e300, 3}}

// TestShapesPaint paints every registered shape with the style variants
// and degenerate sizes, each shape on a page of its own; an internal error
// (a recovered panic) fails the test.
func TestShapesPaint(t *testing.T) {
	var names []string
	for name := range shapeRegistry {
		names = append(names, name)
	}
	slices.Sort(names)
	var b strings.Builder
	b.WriteString("<mxfile>")
	id := 0
	for _, name := range names {
		fmt.Fprintf(&b, `<diagram id="%s" name="%s"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>`, name, name)
		for _, st := range shapeStyles {
			for _, sz := range shapeSizes {
				id++
				style := html.EscapeString("shape=" + name + ";" + st)
				if shapeRegistry[name].paintEdge != nil {
					for _, pts := range [][]float64{{0, 0, 0, 0}, {0, 0, sz[0], sz[1]}, {0, 0, 5, 0, 5, 0, sz[0], sz[1]}, {10, 10, 50, 10, 50, 60, 10, 60}} {
						id++
						fmt.Fprintf(&b, `<mxCell id="e%d" style="%s;curved=%d" edge="1" parent="1"><mxGeometry relative="1" as="geometry">`, id, style, id%2)
						fmt.Fprintf(&b, `<mxPoint x="%g" y="%g" as="sourcePoint"/><mxPoint x="%g" y="%g" as="targetPoint"/><Array as="points">`, pts[0], pts[1], pts[len(pts)-2], pts[len(pts)-1])
						for i := 2; i+3 < len(pts); i += 2 {
							fmt.Fprintf(&b, `<mxPoint x="%g" y="%g"/>`, pts[i], pts[i+1])
						}
						b.WriteString(`</Array></mxGeometry></mxCell>`)
					}
					continue
				}
				fmt.Fprintf(&b, `<mxCell id="v%d" style="%s" vertex="1" parent="1"><mxGeometry x="10" y="10" width="%g" height="%g" as="geometry"/></mxCell>`, id, style, sz[0], sz[1])
			}
		}
		b.WriteString(`</root></mxGraphModel></diagram>`)
	}
	b.WriteString("</mxfile>")
	res, err := Convert([]byte(b.String()), &Options{NoSystemFonts: true, NoTextIndex: true})
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "internal error") {
			t.Error(w)
		}
	}
}

// testShape makes a shape of a registered kind with the given style.
func testShape(t *testing.T, name, st string) *shape {
	t.Helper()
	d := shapeRegistry[name]
	if d == nil {
		t.Fatalf("shape %q is not registered", name)
	}
	s := &shape{conv: &converter{warned: map[string]bool{}}, style: parseStyle("shape="+name+";"+st, false), def: d, name: name}
	s.apply()
	return s
}

// TestLabelMargins checks label margins against values worked out from
// the JavaScript.
func TestLabelMargins(t *testing.T) {
	r := rect{0, 0, 120, 100}
	for _, tc := range []struct {
		name, style string
		want        *rect
	}{
		{"document", "", nil},
		{"document", "boundedLbl=1;", &rect{0, 0, 0, 30}},
		{"document", "boundedLbl=1;size=0.5;", &rect{0, 0, 0, 50}},
		{"cube", "boundedLbl=1;", &rect{20, 20, 0, 0}},
		{"cube", "boundedLbl=0;", nil},
		{"cylinder", "boundedLbl=1;", &rect{0, 30, 0, 0}},
		{"cylinder", "boundedLbl=1;size=0.4;", &rect{0, 40, 0, 0}},
		{"note", "boundedLbl=1;", &rect{0, 30, 0, 0}},
		{"note2", "boundedLbl=1;", &rect{0, 15, 0, 15}},
		{"datastore", "", &rect{0, 2.5 * 13, 0, 0}}, // min(50, round(100/8) + 1 - 1)
		{"folder", "boundedLbl=1;", &rect{0, 15, 0, 0}},
		{"folder", "boundedLbl=1;labelInHeader=1;tabWidth=40;tabHeight=20;", &rect{80, 0, 0, 80}},
		{"folder", "boundedLbl=1;labelInHeader=1;tabPosition=left;rounded=1;", &rect{10, 0, 105, 85}},
		{"cylinder3", "boundedLbl=1;", &rect{0, 30, 0, 4.5}},
		{"cylinder3", "boundedLbl=1;lid=0;", &rect{0, 15, 0, 2.25}},
		{"manualInput", "boundedLbl=1;", &rect{0, 30, 0, 0}},
		{"callout", "", &rect{0, 0, 0, 30}},
		{"umlBoundary", "", &rect{20, 0, 0, 0}},
		{"umlFrame", "", &rect{0, 0, 60, 70}},
		{"umlState", "boundedLbl=1;umlStateConnection=connPointRefEntry;", &rect{10, 0, 0, 0}},
		{"component", "boundedLbl=1;", &rect{0, 30, 0, 0}},
	} {
		s := testShape(t, tc.name, tc.style)
		var got *rect
		if s.def.labelMargins != nil {
			got = s.def.labelMargins(s, r)
		}
		if (got == nil) != (tc.want == nil) || got != nil && !rectNear(*got, *tc.want) {
			t.Errorf("%s %q: margins %v, want %v", tc.name, tc.style, got, tc.want)
		}
	}
}

// TestLabelBounds checks the label rectangles of shapes that override
// getLabelBounds, and of margins turned for the direction.
func TestLabelBounds(t *testing.T) {
	r := rect{10, 20, 120, 60}
	for _, tc := range []struct {
		name, style string
		want        rect
	}{
		{"process", "", rect{22, 20, 96, 60}},
		{"process", "size=0.25;", rect{40, 20, 60, 60}},
		{"process", "fixedSize=1;size=20;", rect{30, 20, 80, 60}},
		{"process", "direction=south;", r},
		{"tape", "boundedLbl=1;", rect{10, 44, 120, 12}},
		{"tape", "boundedLbl=1;direction=south;", rect{58, 20, 24, 60}},
		{"swimlane", "", rect{10, 20, 120, 40}},
		{"swimlane", "startSize=23;", rect{10, 20, 120, 23}},
		{"swimlane", "startSize=100;", rect{10, 20, 120, 60}},
		{"swimlane", "flipV=1;", rect{10, 40, 120, 40}},
		{"swimlane", "direction=south;", rect{90, 20, 40, 60}},
		{"swimlane", "direction=north;", rect{10, 20, 40, 60}},
		{"swimlane", "startSize=0;fixedHeader=0;", r},
		{"table", "startSize=0;", r},
		{"umlLifeline", "", rect{10, 20, 120, 40}},
		{"umlLifeline", "lifelineMirror=1;", rect{10, 20, 120, 30}},
		{"umlControl", "", rect{10, 27.5, 120, 52.5}},
		{"rhombus", "double=1;", rect{14, 24, 112, 52}},
		{"ext", "double=1;strokeWidth=2;margin=1;", rect{14, 24, 112, 52}},
		{"endState", "", rect{14, 24, 112, 52}},
		{"doubleEllipse", "", rect{14, 24, 112, 52}},
		{"rectangle", "footerSize=10;", rect{10, 20, 120, 50}},
		// mxShape.getLabelBounds: margins turned for the direction and flips
		{"document", "boundedLbl=1;", rect{10, 20, 120, 42}},
		{"document", "boundedLbl=1;direction=south;", rect{28, 20, 102, 60}},
		{"document", "boundedLbl=1;direction=north;", rect{10, 20, 102, 60}},
		{"document", "boundedLbl=1;direction=west;", rect{10, 38, 120, 42}},
		{"document", "boundedLbl=1;flipV=1;", rect{10, 38, 120, 42}},
		{"cube", "boundedLbl=1;", rect{30, 40, 100, 40}},
		{"umlFrame", "", rect{10, 20, 60, 30}},
	} {
		s := testShape(t, tc.name, tc.style)
		var got rect
		if s.def.labelBounds != nil {
			got = s.def.labelBounds(s, r)
		} else {
			got = shapeLabelBounds(s, r)
		}
		if !rectNear(got, tc.want) {
			t.Errorf("%s %q: label bounds %v, want %v", tc.name, tc.style, got, tc.want)
		}
	}
}

func rectNear(a, b rect) bool {
	const eps = 1e-9
	return math.Abs(a.x-b.x) < eps && math.Abs(a.y-b.y) < eps && math.Abs(a.w-b.w) < eps && math.Abs(a.h-b.h) < eps
}

// TestStyleBool checks the JavaScript truthiness of style values.
func TestStyleBool(t *testing.T) {
	st := style{"a": "0", "b": "1", "c": "true", "d": "", "e": "0.0", "f": "false"}
	for k, want := range map[string]bool{"a": false, "b": true, "c": true, "d": false, "e": false, "f": true} {
		if got := styleBool(st, k, !want); got != want {
			t.Errorf("styleBool(%q=%q) = %v, want %v", k, st[k], got, want)
		}
	}
	if !styleBool(st, "missing", true) || styleBool(st, "missing", false) {
		t.Error("styleBool ignores the default")
	}
}

// TestMarkersRegistered checks that the markers draw.io adds are known.
func TestMarkersRegistered(t *testing.T) {
	for _, m := range []string{"dash", "box", "cross", "circle", "circlePlus", "halfCircle", "async", "openAsync",
		"ERone", "ERmandOne", "ERmany", "ERoneToMany", "ERzeroToMany", "ERzeroToOne", "sysMLx", "sysMLLost",
		"sysMLFound", "sysMLPackCont", "sysMLReqInt", "sysMLProvInt", "mermaidExtension", "mermaidDiamond"} {
		if markers[m] == nil {
			t.Errorf("marker %q is not registered", m)
		}
	}
}

// TestWaypoints checks that points closer than a unit to the point before
// are dropped, as mxShape.getWaypoints does.
func TestWaypoints(t *testing.T) {
	got := waypoints([]point{{0, 0}, {0.5, 0}, {10, 0}, {10, 0.2}, {10, 10}})
	want := []point{{0, 0}, {10, 0}, {10, 10}}
	if !slices.Equal(got, want) {
		t.Errorf("waypoints = %v, want %v", got, want)
	}
}

// TestResolveColors checks the inherit color keyword of table rows and cells.
func TestResolveColors(t *testing.T) {
	x := `<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>
<mxCell id="t" style="shape=table;childLayout=tableLayout;strokeColor=#ff0000;fillColor=#00ff00;" vertex="1" parent="1"><mxGeometry width="100" height="30" as="geometry"/></mxCell>
<mxCell id="r" style="shape=tableRow;strokeColor=inherit;fillColor=inherit;" vertex="1" parent="t"><mxGeometry width="100" height="30" as="geometry"/></mxCell>
<mxCell id="c" style="shape=partialRectangle;strokeColor=inherit;fillColor=swimlane;" vertex="1" parent="r"><mxGeometry width="50" height="30" as="geometry"/></mxCell>
<mxCell id="d" style="strokeColor=fillColor;fillColor=strokeColor;" vertex="1" parent="t"><mxGeometry width="50" height="30" as="geometry"/></mxCell>
</root></mxGraphModel>`
	root, err := parseXML([]byte(x))
	if err != nil {
		t.Fatal(err)
	}
	m := parseModel(root)
	v := newView(m, nil, nil)
	c := &converter{warned: map[string]bool{}}
	row := c.newShape(v.state(m.cells["r"]))
	if row.stroke != "#ff0000" || row.fill != "#00ff00" {
		t.Errorf("row: stroke %q fill %q", row.stroke, row.fill)
	}
	// inherit through the row; swimlane takes the fill of the nearest
	// swimlane shape (the row, which inherits it from the table)
	cl := c.newShape(v.state(m.cells["c"]))
	if cl.stroke != "#ff0000" || cl.fill != "#00ff00" {
		t.Errorf("cell: stroke %q fill %q", cl.stroke, cl.fill)
	}
	// fillColor and strokeColor read the parent's style values
	d := c.newShape(v.state(m.cells["d"]))
	if d.stroke != "#00ff00" || d.fill != "#ff0000" {
		t.Errorf("d: stroke %q fill %q", d.stroke, d.fill)
	}
}

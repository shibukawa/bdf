package drawio

import (
	"math"
	"testing"
)

// routeView builds the view of a one-page diagram given as mxGraphModel XML.
func routeView(t *testing.T, xml string) *graphView {
	t.Helper()
	f, err := readFile([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}
	return newView(parseModel(f.pages[0].model), nil, nil)
}

func checkRoute(t *testing.T, v *graphView, id string, want []point) {
	t.Helper()
	st := v.state(v.m.cells[id])
	if st == nil {
		t.Fatalf("edge %s has no state", id)
	}
	got := st.edgePoints()
	if _, ok := comparePoints(got, want); !ok {
		t.Errorf("edge %s:\n got  %s\n want %s", id, fmtPoints(got), fmtPoints(want))
	}
}

// The sequence message style is newer than the draw.io desktop release the
// reference exports come from, so these routes are worked out from
// mxEdgeStyle.SequenceMessage by hand.
func TestSequenceMessage(t *testing.T) {
	v := routeView(t, `<mxGraphModel><root>
		<mxCell id="0"/><mxCell id="1" parent="0"/>
		<mxCell id="A" style="shape=umlLifeline;perimeter=lifelinePerimeter;" vertex="1" parent="1">
			<mxGeometry x="100" y="100" width="100" height="300" as="geometry"/></mxCell>
		<mxCell id="bar" style="rounded=0;perimeter=orthogonalPerimeter;" vertex="1" parent="A">
			<mxGeometry x="45" y="100" width="10" height="80" as="geometry"/></mxCell>
		<mxCell id="B" style="shape=umlLifeline;perimeter=lifelinePerimeter;" vertex="1" parent="1">
			<mxGeometry x="400" y="100" width="100" height="300" as="geometry"/></mxCell>
		<mxCell id="m1" style="edgeStyle=sequenceEdgeStyle;endArrow=none;" edge="1" parent="1" source="A" target="B">
			<mxGeometry relative="1" as="geometry"><Array as="points"><mxPoint x="0" y="200"/></Array></mxGeometry></mxCell>
		<mxCell id="m2" style="edgeStyle=sequenceEdgeStyle;endArrow=none;" edge="1" parent="1" source="A" target="B">
			<mxGeometry relative="1" as="geometry"><Array as="points"><mxPoint x="0" y="20"/></Array></mxGeometry></mxCell>
		<mxCell id="m3" style="edgeStyle=sequenceEdgeStyle;endArrow=none;" edge="1" parent="1" source="A" target="bar">
			<mxGeometry relative="1" as="geometry"/></mxCell>
		<mxCell id="m4" style="edgeStyle=sequenceEdgeStyle;endArrow=none;" edge="1" parent="1" source="A" target="B">
			<mxGeometry relative="1" as="geometry"/></mxCell>
	</root></mxGraphModel>`)
	// at the waypoint's y, on the lifelines' dashed lines
	checkRoute(t, v, "m1", []point{{149.5, 200}, {300, 200}, {449.5, 200}})
	// clamped below the heads
	checkRoute(t, v, "m2", []point{{149.5, 140}, {300, 140}, {449.5, 140}})
	// a self-call to an activation bar on the same lifeline
	checkRoute(t, v, "m3", []point{{149.5, 250}, {179.5, 250}, {179.5, 280}, {155, 280}})
	// no waypoint: the middle of the terminals
	checkRoute(t, v, "m4", []point{{149.5, 250}, {300, 250}, {449.5, 250}})
}

func TestJSNumbers(t *testing.T) {
	for _, c := range []struct{ in, want float64 }{
		{-2.7, -2}, {2.7, 2}, {0, 0},
		// numbers printed in exponential notation give their first digit
		{3.552713678800501e-15, 3}, {-5e-7, -5}, {0.000001, 0}, {1e21, 1}, {2.5e22, 2},
	} {
		if got := jsParseInt(c.in); got != c.want {
			t.Errorf("jsParseInt(%v) = %v, want %v", c.in, got, c.want)
		}
	}
	for _, c := range []struct{ in, want float64 }{{2.5, 3}, {-2.5, -2}, {-2.6, -3}, {0.49999999999999994, 0}} {
		if got := jsRound(c.in); got != c.want {
			t.Errorf("jsRound(%v) = %v, want %v", c.in, got, c.want)
		}
	}
	s := style{"a": "0", "b": "1", "c": "true", "d": "", "e": "2", "f": " 1.0 "}
	for _, c := range []struct {
		key           string
		truthy, isOne bool
	}{{"a", false, false}, {"b", true, true}, {"c", true, false}, {"d", false, false}, {"e", true, false}, {"f", true, true}} {
		if got := styleTruthy(s, c.key, false); got != c.truthy {
			t.Errorf("styleTruthy(%q) = %v", c.key, got)
		}
		if got := styleIsOne(s, c.key, false); got != c.isOne {
			t.Errorf("styleIsOne(%q) = %v", c.key, got)
		}
	}
}

func TestPortConstraints(t *testing.T) {
	edge := &cellState{style: style{"sourcePortConstraint": "south"}}
	for _, c := range []struct {
		st   style
		want int
	}{
		{style{}, dirMaskSouth}, // the edge's sourcePortConstraint
		{style{"portConstraint": "eastwest"}, dirMaskEast | dirMaskWest},
		{style{"portConstraint": "north", "rotation": "90"}, dirMaskNorth},
		{style{"portConstraint": "north", "rotation": "90", "portConstraintRotation": "1"}, dirMaskEast},
		{style{"portConstraint": "north", "rotation": "-100", "portConstraintRotation": "1"}, dirMaskWest},
		{style{"portConstraint": "northeast", "rotation": "180", "portConstraintRotation": "1"}, dirMaskSouth | dirMaskWest},
	} {
		if got := getPortConstraints(&cellState{style: c.st}, edge, true, dirMaskAll); got != c.want {
			t.Errorf("getPortConstraints(%v) = %d, want %d", c.st, got, c.want)
		}
	}
	if got := reversePortConstraints(dirMaskNorth | dirMaskWest); got != dirMaskSouth|dirMaskEast {
		t.Errorf("reversePortConstraints = %d", got)
	}
}

func TestRotatedBounds(t *testing.T) {
	r := rotatedBounds(rect{0, 0, 100, 50}, 90)
	want := rect{25, -25, 50, 100}
	if math.Abs(r.x-want.x) > 1e-9 || math.Abs(r.y-want.y) > 1e-9 || math.Abs(r.w-want.w) > 1e-9 || math.Abs(r.h-want.h) > 1e-9 {
		t.Errorf("rotatedBounds = %v, want %v", r, want)
	}
}

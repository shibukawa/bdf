package drawio

import (
	"fmt"
	"html"
	"math"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// The shapes of mxBpmnShape2.js and mxArrows.js.
var (
	bpmnShapeNames = []string{"mxgraph.bpmn.shape", "mxgraph.bpmn.sendMarker", "mxgraph.bpmn.event", "mxgraph.bpmn.gateway2",
		"mxgraph.bpmn.task", "mxgraph.bpmn.task2", "mxgraph.bpmn.data", "mxgraph.bpmn.data2", "mxgraph.bpmn.swimlane",
		"mxgraph.bpmn.conversation", "mxgraph.bpmn.conversation2"}
	arrows2ShapeNames = []string{"mxgraph.arrows2.arrow", "mxgraph.arrows2.twoWayArrow", "mxgraph.arrows2.stylisedArrow",
		"mxgraph.arrows2.sharpArrow", "mxgraph.arrows2.sharpArrow2", "mxgraph.arrows2.calloutArrow", "mxgraph.arrows2.bendArrow",
		"mxgraph.arrows2.bendDoubleArrow", "mxgraph.arrows2.calloutDoubleArrow", "mxgraph.arrows2.calloutQuadArrow",
		"mxgraph.arrows2.calloutDouble90Arrow", "mxgraph.arrows2.quadArrow", "mxgraph.arrows2.triadArrow",
		"mxgraph.arrows2.tailedArrow", "mxgraph.arrows2.tailedNotchedArrow", "mxgraph.arrows2.stripedArrow",
		"mxgraph.arrows2.jumpInArrow", "mxgraph.arrows2.uTurnArrow", "mxgraph.arrows2.wedgeArrow",
		"mxgraph.arrows2.wedgeArrowDashed", "mxgraph.arrows2.wedgeArrowDashed2"}
)

// libStyles are style variants with the keys the BPMN and arrows2 shapes
// read, including values that are not numbers and absurd sizes.
var libStyles = func() []string {
	var out []string
	outlines := []string{"none", "standard", "eventInt", "eventNonint", "catching", "boundInt", "boundNonint", "throwing", "end", "foo"}
	symbols := []string{"general", "message", "timer", "escalation", "conditional", "link", "error", "cancel", "compensation",
		"signal", "multiple", "parallelMultiple", "terminate", "star", "exclusiveGw", "parallelGw", "complexGw", "standard"}
	for i, o := range outlines {
		for j, sym := range symbols {
			st := "outline=" + o + ";symbol=" + sym + ";"
			if (i+j)%3 == 0 {
				st += "background=gateway;fillColor=none;"
			}
			out = append(out, st)
		}
	}
	for _, gw := range []string{"event", "exclusive", "parallel", "complex", "foo"} {
		out = append(out, "gwType="+gw+";outline=end;symbol=message;")
	}
	for _, m := range []string{"abstract", "service", "send", "receive", "user", "manual", "businessRule", "script", "nime"} {
		out = append(out, "taskMarker="+m+";bpmnShapeType=transaction;isLoopStandard=1;isLoopMultiParallel=1;isLoopMultiSeq=1;isLoopComp=1;isLoopSub=1;isAdHoc=1;")
	}
	return append(out,
		"bpmnShapeType=call;strokeWidth=1e308;indent=-5;outline=end;",
		"bpmnShapeType=subprocess;strokeWidth=0;isAdHoc=1;fillColor=none;strokeColor=none;",
		"bpmnTransferType=input;isCollection=1;size=1e308;",
		"bpmnTransferType=output;isCollection=1;strokeWidth=5;",
		"bpmnConversationType=call;isLoopSub=1;",
		"isCollection=1;startSize=1e308;horizontal=0;",
		"dx=1e308;dy=-1e308;notch=1e308;arrowHead=-1;headCrossline=1;tailCrossline=1;boundedLbl=1;",
		"dx=abc;dy=;notch=NaN;feather=x;dy1=1e308;dx1=-1;dx2=1e308;dy2=0;dy3=abc;dx3=1;rounded=1;",
		"dy1=10;dx1=20;dx2=25;dy2=30;notch=20;arrowHead=20;",
		"startWidth=1e308;stepSize=1e-300;",
		"startWidth=-1;stepSize=abc;",
	)
}()

// TestLibShapesPaint paints the BPMN and arrows2 shapes with their style
// keys and degenerate sizes (TestShapesPaint paints them with the common
// style variants); an internal error (a recovered panic) fails the test.
func TestLibShapesPaint(t *testing.T) {
	var b strings.Builder
	b.WriteString("<mxfile>")
	id := 0
	for _, name := range append(append([]string(nil), bpmnShapeNames...), arrows2ShapeNames...) {
		fmt.Fprintf(&b, `<diagram id="%s" name="%s"><mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>`, name, name)
		for i, st := range libStyles {
			if !strings.HasPrefix(name, "mxgraph.bpmn.") && i < len(libStyles)-11 {
				continue // the BPMN keys
			}
			for _, sz := range shapeSizes {
				id++
				style := html.EscapeString("shape=" + name + ";" + st)
				if shapeRegistry[name].paintEdge != nil {
					fmt.Fprintf(&b, `<mxCell id="e%d" style="%s" edge="1" parent="1"><mxGeometry relative="1" as="geometry">`, id, style)
					fmt.Fprintf(&b, `<mxPoint x="0" y="0" as="sourcePoint"/><mxPoint x="%g" y="%g" as="targetPoint"/></mxGeometry></mxCell>`, sz[0], sz[1])
					continue
				}
				fmt.Fprintf(&b, `<mxCell id="v%d" value="label" style="%s" vertex="1" parent="1"><mxGeometry x="10" y="10" width="%g" height="%g" as="geometry"/></mxCell>`, id, style, sz[0], sz[1])
			}
		}
		b.WriteString(`</root></mxGraphModel></diagram>`)
	}
	b.WriteString("</mxfile>")
	res, err := Convert([]byte(b.String()), &Options{NoSystemFonts: true, NoTextIndex: true})
	if err != nil {
		t.Fatal(err)
	}
	if n := len(bpmnShapeNames) + len(arrows2ShapeNames); res.Pages != n {
		t.Errorf("%d pages, want %d", res.Pages, n)
	}
	for _, w := range res.Warnings {
		if strings.Contains(w, "internal error") || strings.Contains(w, "not supported") {
			t.Error(w)
		}
	}
}

// TestLibShapesRegistered checks that the shapes of both libraries are
// registered, and that the stencils of the task markers resolve.
func TestLibShapesRegistered(t *testing.T) {
	for _, name := range append(append([]string(nil), bpmnShapeNames...), arrows2ShapeNames...) {
		if shapeRegistry[name] == nil {
			t.Errorf("shape %q is not registered", name)
		}
	}
	var got []string
	for name := range shapeRegistry {
		if strings.HasPrefix(name, "mxgraph.bpmn.") || strings.HasPrefix(name, "mxgraph.arrows2.") {
			got = append(got, name)
		}
	}
	if len(got) != len(bpmnShapeNames)+len(arrows2ShapeNames) {
		slices.Sort(got)
		t.Errorf("registered library shapes %v", got)
	}
	for _, name := range []string{"mxgraph.bpmn.loop", "mxgraph.bpmn.compensation", "mxgraph.bpmn.ad_hoc", "mxgraph.bpmn.service_task",
		"mxgraph.bpmn.user_task", "mxgraph.bpmn.manual_task", "mxgraph.bpmn.business_rule_task", "mxgraph.bpmn.script_task"} {
		if lookupStencil(name) == nil {
			t.Errorf("stencil %q of the task markers is not found", name)
		}
	}
	// the library shapes are JavaScript shapes, not stencils
	for _, name := range []string{"mxgraph.bpmn.shape", "mxgraph.bpmn.task", "mxgraph.arrows2.arrow"} {
		if lookupStencil(name) != nil {
			t.Errorf("%q resolves to a stencil", name)
		}
	}
}

// TestLibLabelBounds checks the bounded labels of arrow and twoWayArrow
// and the label margins of data objects against values worked out from
// the JavaScript.
func TestLibLabelBounds(t *testing.T) {
	r := rect{10, 20, 110, 70}
	for _, tc := range []struct {
		name, style string
		want        rect
	}{
		{"mxgraph.arrows2.arrow", "dy=0.6;dx=40;", r},
		{"mxgraph.arrows2.arrow", "dy=0.6;dx=40;boundedLbl=1;", rect{10, 41, 70, 28}},
		{"mxgraph.arrows2.arrow", "dy=0.6;dx=40;boundedLbl=1;flipH=1;", rect{50, 41, 70, 28}},
		{"mxgraph.arrows2.arrow", "dy=0.6;dx=40;boundedLbl=1;direction=west;flipH=1;", rect{10, 41, 70, 28}},
		{"mxgraph.arrows2.arrow", "dy=0.6;dx=40;boundedLbl=1;direction=south;", rect{43, 20, 44, 30}},
		{"mxgraph.arrows2.arrow", "dy=0.6;dx=40;boundedLbl=1;direction=south;flipV=1;", rect{43, 60, 44, 30}},
		{"mxgraph.arrows2.arrow", "boundedLbl=1;", rect{10, 37.5, 109.5, 35}},
		{"mxgraph.arrows2.twoWayArrow", "dy=0.6;dx=35;boundedLbl=1;", rect{45, 41, 40, 28}},
		{"mxgraph.arrows2.twoWayArrow", "dy=0.6;dx=35;boundedLbl=1;direction=north;", rect{43, 55, 44, 0}},
		{"mxgraph.bpmn.swimlane", "startSize=20;", rect{10, 20, 110, 20}},
	} {
		s := testShape(t, tc.name, tc.style)
		if got := s.labelBounds(r); !rectNear(got, tc.want) {
			t.Errorf("%s %q: label bounds %v, want %v", tc.name, tc.style, got, tc.want)
		}
	}
	s := testShape(t, "mxgraph.bpmn.data2", "boundedLbl=1;size=15;")
	// mxCylinder.getLabelMargins: min(maxHeight, height * size * 2)
	if got := s.def.labelMargins(s, r); got == nil || !rectNear(*got, rect{0, 40, 0, 0}) {
		t.Errorf("data2 label margins %v, want the note's", got)
	}
}

// TestArrows2Undefined checks that the arrows whose sizes have no
// defaults draw nothing without them, as the path of NaN coordinates
// that draw.io writes into the SVG, and that sizes that are not numbers
// break the path the same way.
func TestArrows2Undefined(t *testing.T) {
	for _, tc := range []struct {
		name, style string
		drawn       bool
	}{
		{"mxgraph.arrows2.tailedArrow", "", false},
		{"mxgraph.arrows2.tailedNotchedArrow", "dy1=10;dx1=20;", false},
		{"mxgraph.arrows2.tailedArrow", "dy1=10;dx1=20;dx2=25;dy2=30;", true},
		{"mxgraph.arrows2.uTurnArrow", "dy=11;arrowHead=43;dx2=25;", true},
		{"mxgraph.arrows2.arrow", "", true},
		{"mxgraph.arrows2.arrow", "dy=abc;", false},
	} {
		c := &converter{warned: map[string]bool{}}
		s := testShape(t, tc.name, tc.style)
		s.conv = c
		cv := c.newCanvas()
		c2 := newC2D(cv)
		s.configureCanvas(c2, 0, 0, 100, 100)
		s.def.paintVertex(s, c2, 0, 0, 100, 100)
		if cv.drawn != tc.drawn {
			t.Errorf("%s %q: drawn %v, want %v", tc.name, tc.style, cv.drawn, tc.drawn)
		}
	}
	if v := jsStyleFloat(style{"a": "12px", "b": "x"}, "a", 0); v != 12 {
		t.Errorf("parseFloat(12px) = %v", v)
	}
	if v := jsStyleFloat(style{"b": "x"}, "b", 1); !math.IsNaN(v) {
		t.Errorf("parseFloat(x) = %v, want NaN", v)
	}
}

// TestLibShapesTestdata converts the test diagrams of testdata/shapes_libs
// (compared with draw.io's exports when the shapes were ported) without
// warnings.
func TestLibShapesTestdata(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join("testdata", "shapes_libs", "*.drawio"))
	if len(files) == 0 {
		t.Fatal("no test diagrams")
	}
	for _, f := range files {
		res, err := ConvertFile(f, testOptions())
		if err != nil {
			t.Errorf("%s: %v", f, err)
			continue
		}
		for _, w := range res.Warnings {
			t.Errorf("%s: %s", f, w)
		}
	}
}

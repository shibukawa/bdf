//go:build ignore

// Gen writes the test diagrams of the BPMN (mxgraph.bpmn.*) and arrows2
// (mxgraph.arrows2.*) shapes of this directory: every event outline with
// every symbol, gateways, tasks with their markers, data objects, lanes,
// conversations, and every arrows2 shape with the palette's style, the
// constructors' defaults and varied sizes, directions, flips and colors.
// Headers and labels name the styles (draw.io and the converter both draw
// them, which helps to compare an export shape by shape).
//
//	go run ./converter/drawio/testdata/shapes_libs/gen.go converter/drawio/testdata/shapes_libs
//	/Applications/draw.io.app/Contents/MacOS/draw.io -x -f png -s 2 -o ref.png converter/drawio/testdata/shapes_libs/bpmn_event.drawio
//	go run ./cmd/bdf generate -q converter/drawio/testdata/shapes_libs/bpmn_event.drawio out.bdf
//	node test/render.mjs out.bdf out/ 2.6667
package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

type doc struct {
	b  strings.Builder
	id int
}

func (d *doc) next() string { d.id++; return fmt.Sprint("c", d.id) }

func (d *doc) cell(style string, x, y, w, h float64, value string) {
	fmt.Fprintf(&d.b, `        <mxCell id="%s" value="%s" style="%s" vertex="1" parent="1">
          <mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/>
        </mxCell>
`, d.next(), html.EscapeString(value), html.EscapeString(style), x, y, w, h)
}

// text writes a small label centered in the box x, y, w, 12.
func (d *doc) text(x, y, w float64, s string) {
	d.cell("text;html=1;align=center;verticalAlign=top;fontSize=8;spacing=0;whiteSpace=wrap;", x, y, w, 12, s)
}

func (d *doc) edge(style string, pts ...[2]float64) {
	fmt.Fprintf(&d.b, `        <mxCell id="%s" value="" style="%s" edge="1" parent="1">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="%g" y="%g" as="sourcePoint"/>
            <mxPoint x="%g" y="%g" as="targetPoint"/>
`, d.next(), html.EscapeString(style), pts[0][0], pts[0][1], pts[len(pts)-1][0], pts[len(pts)-1][1])
	if len(pts) > 2 {
		d.b.WriteString("            <Array as=\"points\">\n")
		for _, p := range pts[1 : len(pts)-1] {
			fmt.Fprintf(&d.b, "              <mxPoint x=\"%g\" y=\"%g\"/>\n", p[0], p[1])
		}
		d.b.WriteString("            </Array>\n")
	}
	d.b.WriteString("          </mxGeometry>\n        </mxCell>\n")
}

// grid places a shape for each row and column, w by h in cells of cw by
// ch, with the column names above and the row names at the left; it
// returns the y below the grid.
func (d *doc) grid(y0 float64, rows, cols []string, cw, ch, w, h float64, style func(r, c string) string) float64 {
	const x0 = 80
	for j, c := range cols {
		d.text(x0+float64(j)*cw, y0, cw, c)
	}
	for i, r := range rows {
		y := y0 + 14 + float64(i)*ch
		d.text(0, y+(ch-12)/2, x0-4, r)
		for j, c := range cols {
			d.cell(style(r, c), x0+float64(j)*cw+(cw-w)/2, y+(ch-h)/2, w, h, "")
		}
	}
	return y0 + 14 + float64(len(rows))*ch + 16
}

// row places the shapes of styles in a row of cells cw wide with each
// style's label below; it returns the y below the row.
func (d *doc) row(y0 float64, cw, w, h float64, styles, labels []string) float64 {
	for i, s := range styles {
		x := 20 + float64(i)*cw
		d.cell(s, x+(cw-w)/2, y0, w, h, "")
		lbl := s
		if labels != nil {
			lbl = labels[i]
		}
		d.text(x, y0+h+4, cw, lbl)
	}
	return y0 + h + 40
}

func (d *doc) write(path string) {
	out := `<mxfile host="gen">
  <diagram id="p1" name="Page-1">
    <mxGraphModel grid="0" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="0" pageScale="1" pageWidth="827" pageHeight="1169" math="0" shadow="0">
      <root>
        <mxCell id="0"/>
        <mxCell id="1" parent="0"/>
` + d.b.String() + `      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
`
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		panic(err)
	}
}

var (
	outlines     = []string{"none", "standard", "eventInt", "eventNonint", "catching", "boundInt", "boundNonint", "throwing", "end"}
	eventSymbols = []string{"general", "message", "timer", "escalation", "conditional", "link", "error", "cancel", "compensation", "signal", "multiple", "parallelMultiple", "terminate", "star"}
	shapeSymbols = append(append([]string(nil), eventSymbols...), "exclusiveGw", "parallelGw", "complexGw")
)

func main() {
	dir := os.Args[1]

	// mxgraph.bpmn.shape: every outline with every symbol, and on a gateway
	d := &doc{}
	y := d.grid(0, outlines, shapeSymbols, 56, 52, 40, 40, func(r, c string) string {
		return "shape=mxgraph.bpmn.shape;html=1;perimeter=ellipsePerimeter;outline=" + r + ";symbol=" + c + ";"
	})
	y = d.grid(y, outlines, shapeSymbols, 56, 62, 54, 54, func(r, c string) string {
		return "shape=mxgraph.bpmn.shape;html=1;perimeter=rhombusPerimeter;background=gateway;outline=" + r + ";symbol=" + c + ";"
	})
	d.row(y, 90, 50, 50, []string{
		"shape=mxgraph.bpmn.shape;html=1;outline=throwing;symbol=message;fillColor=#dae8fc;strokeColor=#6c8ebf;",
		"shape=mxgraph.bpmn.shape;html=1;outline=end;symbol=error;fillColor=none;strokeColor=#b85450;",
		"shape=mxgraph.bpmn.shape;html=1;outline=catching;symbol=timer;strokeWidth=3;",
		"shape=mxgraph.bpmn.shape;html=1;outline=boundNonint;symbol=star;fillColor=#fff2cc;",
		"shape=mxgraph.bpmn.shape;html=1;outline=standard;symbol=link;rotation=30;",
		"shape=mxgraph.bpmn.shape;html=1;outline=standard;symbol=link;direction=south;",
		"shape=mxgraph.bpmn.shape;html=1;outline=standard;symbol=link;flipH=1;",
		"shape=mxgraph.bpmn.shape;html=1;background=gateway;outline=none;symbol=exclusiveGw;fillColor=#d5e8d4;gradientColor=#97d077;",
		"shape=mxgraph.bpmn.shape;html=1;background=foo;outline=standard;symbol=signal;",
		"shape=mxgraph.bpmn.shape;html=1;outline=end;symbol=message;fillColor=none;",
	}, []string{"throwing colors", "end fill none", "strokeWidth=3", "star fill", "rotation=30", "direction=south", "flipH=1", "gateway gradient", "background=foo", "end message none"})
	d.write(filepath.Join(dir, "bpmn_shape.drawio"))

	// mxgraph.bpmn.event: every outline with every symbol, and variants
	d = &doc{}
	y = d.grid(0, outlines, eventSymbols, 56, 52, 40, 40, func(r, c string) string {
		return "shape=mxgraph.bpmn.event;html=1;perimeter=ellipsePerimeter;outline=" + r + ";symbol=" + c + ";"
	})
	y = d.grid(y, outlines, eventSymbols, 56, 52, 40, 40, func(r, c string) string {
		return "shape=mxgraph.bpmn.event;html=1;outline=" + r + ";symbol=" + c + ";fillColor=#dae8fc;strokeColor=#6c8ebf;strokeWidth=2;"
	})
	d.row(y, 90, 60, 40, []string{
		"shape=mxgraph.bpmn.event;html=1;outline=end;symbol=terminate;fillColor=none;",
		"shape=mxgraph.bpmn.event;html=1;outline=throwing;symbol=message;fillColor=none;",
		"shape=mxgraph.bpmn.event;html=1;outline=standard;symbol=escalation;rotation=45;",
		"shape=mxgraph.bpmn.event;html=1;outline=catching;symbol=link;direction=north;",
		"shape=mxgraph.bpmn.event;html=1;outline=catching;symbol=link;flipV=1;flipH=1;",
		"shape=mxgraph.bpmn.event;html=1;outline=eventNonint;symbol=timer;dashed=1;dashPattern=1 2;",
		"shape=mxgraph.bpmn.event;html=1;symbol=signal;",
		"shape=mxgraph.bpmn.event;html=1;outline=boundInt;symbol=multiple;shadow=1;",
		"shape=mxgraph.bpmn.event;html=1;outline=end;symbol=cancel;strokeColor=#82b366;fillColor=#d5e8d4;gradientColor=#ffffff;",
	}, []string{"end terminate none", "throwing msg none", "rotation=45", "direction=north", "flipH flipV", "dashed", "no outline", "shadow", "gradient"})
	d.write(filepath.Join(dir, "bpmn_event.drawio"))

	// mxgraph.bpmn.gateway2
	d = &doc{}
	gw := "shape=mxgraph.bpmn.gateway2;html=1;perimeter=rhombusPerimeter;"
	y = d.grid(0, []string{"exclusive", "parallel", "complex", "event"}, []string{"default", "outline=none", "outline=catching", "outline=end;symbol=multiple", "strokeWidth=3", "fillColor=#dae8fc", "rotation=30", "direction=south", "flipH=1"}, 80, 70, 60, 60, func(r, c string) string {
		if c == "default" {
			c = ""
		}
		return gw + "gwType=" + r + ";" + c + ";"
	})
	y = d.grid(y, []string{"catching", "throwing", "boundNonint", "end", "standard"}, []string{"general", "message", "timer", "conditional", "multiple", "parallelMultiple", "signal", "star", "standard"}, 80, 70, 60, 60, func(r, c string) string {
		return gw + "outline=" + r + ";symbol=" + c + ";"
	})
	d.row(y, 90, 60, 60, []string{
		gw + "gwType=exclusive;fillColor=none;",
		gw + "gwType=parallel;strokeColor=#b85450;fillColor=#f8cecc;",
		gw + "gwType=foo;",
		gw + "outline=catching;symbol=message;gwType=parallel;",
		gw + "gwType=exclusive;shadow=1;",
	}, []string{"exclusive none", "parallel red", "gwType=foo", "parallel+outline", "shadow"})
	d.write(filepath.Join(dir, "bpmn_gateway.drawio"))

	// tasks: types, markers, loop markers and events
	d = &doc{}
	y = 0
	markers := []string{"abstract", "service", "send", "receive", "user", "manual", "businessRule", "script", "nime"}
	types := []string{"task", "transaction", "call", "subprocess"}
	for _, task := range []string{"task", "task2"} {
		y = d.grid(y, types, markers, 110, 80, 100, 70, func(r, c string) string {
			return "shape=mxgraph.bpmn." + task + ";html=1;rectStyle=rounded;size=10;bpmnShapeType=" + r + ";taskMarker=" + c + ";"
		})
		loops := []string{"isLoopStandard=1", "isLoopMultiParallel=1", "isLoopMultiSeq=1", "isLoopComp=1", "isLoopSub=1", "isAdHoc=1",
			"isLoopStandard=1;isLoopSub=1", "isLoopMultiSeq=1;isLoopComp=1;isAdHoc=1", "isLoopStandard=1;isLoopMultiParallel=1;isLoopMultiSeq=1;isLoopComp=1;isLoopSub=1;isAdHoc=1"}
		y = d.grid(y, []string{task, task + " sw=3", task + " call"}, loops, 110, 80, 100, 70, func(r, c string) string {
			extra := ""
			switch {
			case strings.HasSuffix(r, "sw=3"):
				extra = "strokeWidth=3;"
			case strings.HasSuffix(r, "call"):
				extra = "bpmnShapeType=call;strokeWidth=2;"
			}
			return "shape=mxgraph.bpmn." + task + ";html=1;rectStyle=rounded;size=10;taskMarker=abstract;" + extra + c + ";"
		})
		y = d.grid(y, []string{task + " events", task + " sw=2"}, []string{"outline=standard;symbol=message", "outline=catching;symbol=timer", "outline=throwing;symbol=signal", "outline=end;symbol=error", "outline=boundNonint;symbol=conditional", "outline=eventInt;symbol=star", "taskMarker=send;outline=end;symbol=general", "taskMarker=user;outline=catching;symbol=message", "outline=standard"}, 110, 80, 100, 70, func(r, c string) string {
			extra := ""
			if strings.HasSuffix(r, "sw=2") {
				extra = "strokeWidth=2;"
			}
			return "shape=mxgraph.bpmn." + task + ";html=1;rectStyle=rounded;size=10;taskMarker=abstract;" + extra + c + ";"
		})
		y = d.row(y, 110, 100, 70, []string{
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=service;fillColor=#dae8fc;strokeColor=#6c8ebf;isLoopStandard=1;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=send;fillColor=#fff2cc;strokeColor=#d6b656;isAdHoc=1;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=user;rotation=20;isLoopSub=1;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=manual;direction=south;isLoopMultiSeq=1;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=script;flipH=1;isLoopComp=1;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=businessRule;dashed=1;bpmnShapeType=transaction;indent=6;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=receive;rectStyle=snip;size=15;fillColor=none;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=abstract;shadow=1;gradientColor=#7ea6e0;fillColor=#dae8fc;",
			"shape=mxgraph.bpmn." + task + ";html=1;taskMarker=service;strokeWidth=4;bpmnShapeType=call;",
		}, []string{"service blue", "send yellow adhoc", "user rot20", "manual south", "script flipH", "rule transaction", "receive snip", "shadow gradient", "service call sw4"})
	}
	d.write(filepath.Join(dir, "bpmn_task.drawio"))

	// data objects, lanes, conversations and the send marker
	d = &doc{}
	y = 0
	for _, data := range []string{"data", "data2"} {
		s := "shape=mxgraph.bpmn." + data + ";html=1;size=15;"
		y = d.row(y, 70, 40, 60, []string{
			s, s + "isCollection=1;", s + "bpmnTransferType=input;", s + "bpmnTransferType=input;isCollection=1;",
			s + "bpmnTransferType=output;", s + "bpmnTransferType=output;isCollection=1;",
			s + "bpmnTransferType=output;strokeWidth=3;isCollection=1;", s + "bpmnTransferType=input;fillColor=#dae8fc;strokeColor=#6c8ebf;",
			s + "bpmnTransferType=output;direction=south;", s + "bpmnTransferType=output;flipH=1;isCollection=1;", s + "size=5;rotation=30;bpmnTransferType=input;",
		}, []string{data, "coll", "input", "input coll", "output", "output coll", "sw=3", "blue", "south", "flipH", "rot30"})
	}
	for _, conv := range []string{"conversation", "conversation2"} {
		s := "shape=mxgraph.bpmn." + conv + ";html=1;perimeter=hexagonPerimeter2;"
		y = d.row(y, 90, 70, 60, []string{
			s, s + "isLoopSub=1;", s + "bpmnConversationType=call;", s + "bpmnConversationType=call;isLoopSub=1;",
			s + "strokeWidth=3;isLoopSub=1;", s + "strokeWidth=2;bpmnConversationType=call;isLoopSub=1;",
			s + "fillColor=#d5e8d4;strokeColor=#82b366;direction=south;isLoopSub=1;",
		}, []string{conv, "loopSub", "call", "call loopSub", "sw3 loopSub", "sw2 call loopSub", "green south"})
	}
	y = d.row(y, 90, 60, 40, []string{
		"shape=mxgraph.bpmn.sendMarker;html=1;",
		"shape=mxgraph.bpmn.sendMarker;html=1;fillColor=#000000;strokeColor=#ffffff;",
		"shape=mxgraph.bpmn.sendMarker;html=1;rotation=90;",
	}, []string{"sendMarker", "inverse", "rotation=90"})
	lane := "shape=mxgraph.bpmn.swimlane;html=1;startSize=20;swimlaneLine=1;collapsible=0;fontStyle=0;swimlaneFillColor=#ffffff;whiteSpace=wrap;"
	d.cell(lane+"horizontal=0;strokeWidth=2;isCollection=1;", 20, y, 300, 100, "Lane")
	d.cell(lane+"horizontal=1;strokeWidth=2;isCollection=1;", 340, y, 150, 100, "Lane")
	d.cell(lane+"horizontal=1;isCollection=0;rounded=1;fillColor=#dae8fc;", 510, y, 150, 100, "Lane")
	d.cell(lane+"horizontal=1;isCollection=1;startSize=0;", 680, y, 150, 100, "Lane")
	d.write(filepath.Join(dir, "bpmn_misc.drawio"))

	// arrows2: the palette's styles, then the constructors' defaults
	d = &doc{}
	a := "html=1;shadow=0;dashed=0;align=center;verticalAlign=middle;shape=mxgraph.arrows2."
	palette := []struct {
		s    string
		w, h float64
	}{
		{"arrow;dy=0.6;dx=40;notch=0;", 100, 70}, {"arrow;dy=0.6;dx=40;flipH=1;notch=0;", 100, 70},
		{"arrow;dy=0.6;dx=40;direction=north;notch=0;", 70, 100}, {"arrow;dy=0.6;dx=40;direction=south;notch=0;", 70, 100},
		{"arrow;dy=0;dx=30;notch=30;", 100, 60}, {"arrow;dy=0.6;dx=40;notch=15;", 100, 70},
		{"arrow;dy=0;dx=10;notch=10;", 100, 30}, {"arrow;dy=0;dx=10;notch=0;", 100, 30},
		{"arrow;dy=0.67;dx=20;notch=0;", 100, 60}, {"twoWayArrow;dy=0.6;dx=35;", 100, 60},
		{"twoWayArrow;dy=0.65;dx=22;", 100, 60}, {"stylisedArrow;dy=0.6;dx=40;notch=15;feather=0.4;", 100, 60},
		{"sharpArrow;dy1=0.67;dx1=18;dx2=18;notch=0;", 100, 60}, {"sharpArrow2;dy1=0.67;dx1=18;dx2=18;dy3=0.15;dx3=27;notch=0;", 100, 60},
		{"calloutArrow;dy=10;dx=20;notch=60;arrowHead=10;", 100, 60}, {"bendArrow;dy=15;dx=38;notch=0;arrowHead=55;rounded=0;", 100, 100},
		{"bendArrow;dy=15;dx=38;notch=0;arrowHead=55;rounded=1;", 100, 100}, {"bendDoubleArrow;dy=15;dx=38;arrowHead=55;rounded=0;", 100, 100},
		{"bendDoubleArrow;dy=15;dx=38;arrowHead=55;rounded=1;", 100, 100}, {"calloutDoubleArrow;dy=10;dx=20;notch=24;arrowHead=10;", 100, 50},
		{"calloutQuadArrow;dy=10;dx=20;notch=24;arrowHead=10;", 100, 100}, {"calloutDouble90Arrow;dy1=10;dx1=20;dx2=70;dy2=70;arrowHead=10;", 100, 100},
		{"quadArrow;dy=10;dx=20;notch=24;arrowHead=10;", 100, 100}, {"triadArrow;dy=10;dx=20;arrowHead=40;", 100, 70},
		{"tailedArrow;dy1=10;dx1=20;notch=0;arrowHead=20;dx2=25;dy2=30;", 100, 60}, {"tailedNotchedArrow;dy1=10;dx1=20;notch=20;arrowHead=20;dx2=25;dy2=30;", 100, 60},
		{"stripedArrow;dy=0.6;dx=40;notch=25;", 100, 70}, {"jumpInArrow;dy=15;dx=38;arrowHead=55;", 100, 100},
		{"uTurnArrow;dy=11;arrowHead=43;dx2=25;", 100, 100},
	}
	for i, p := range palette {
		col, r := i%8, i/8
		x, y := 20+float64(col)*130, float64(r)*150
		d.cell(a+p.s, x+(120-p.w)/2, y, p.w, p.h, "")
		d.text(x, y+104, 120, p.s)
	}
	y = 4*150 + 20
	names := []string{"arrow", "twoWayArrow", "stylisedArrow", "sharpArrow", "sharpArrow2", "calloutArrow", "bendArrow", "bendDoubleArrow",
		"calloutDoubleArrow", "calloutQuadArrow", "calloutDouble90Arrow", "quadArrow", "triadArrow", "tailedArrow", "tailedNotchedArrow",
		"stripedArrow", "jumpInArrow", "uTurnArrow"}
	for i, n := range names {
		col, r := i%9, i/9
		x, yy := 20+float64(col)*115, y+float64(r)*110
		d.cell(a+n+";", x+10, yy, 90, 70, "")
		d.text(x, yy+74, 110, n+" (defaults)")
	}
	d.write(filepath.Join(dir, "arrows2.drawio"))

	// arrows2 with varied sizes, directions, flips and colors
	d = &doc{}
	variants := []string{"", "rotation=30;", "direction=south;", "direction=north;", "direction=west;", "flipH=1;", "flipV=1;",
		"strokeWidth=3;dashed=1;", "fillColor=#dae8fc;gradientColor=#7ea6e0;strokeColor=#6c8ebf;opacity=70;", "rounded=1;shadow=1;"}
	sizes := map[string]string{
		"arrow":                "dy=0.3;dx=25;notch=20;headCrossline=1;tailCrossline=1;",
		"twoWayArrow":          "dy=0.4;dx=15;",
		"stylisedArrow":        "dy=0.3;dx=20;notch=10;feather=0.8;",
		"sharpArrow":           "dy1=0.3;dx1=30;dx2=10;notch=15;",
		"sharpArrow2":          "dy1=0.4;dx1=30;dx2=12;dy3=0.3;dx3=40;notch=15;",
		"calloutArrow":         "dy=6;dx=15;notch=40;arrowHead=15;",
		"bendArrow":            "dy=8;dx=25;notch=12;arrowHead=40;rounded=1;",
		"bendDoubleArrow":      "dy=8;dx=25;arrowHead=40;rounded=1;",
		"calloutDoubleArrow":   "dy=6;dx=15;notch=15;arrowHead=12;",
		"calloutQuadArrow":     "dy=6;dx=15;notch=15;arrowHead=12;",
		"calloutDouble90Arrow": "dy1=6;dx1=15;dx2=40;dy2=40;arrowHead=12;",
		"quadArrow":            "dy=6;dx=15;arrowHead=12;",
		"triadArrow":           "dy=6;dx=15;arrowHead=30;",
		"tailedArrow":          "dy1=6;dx1=15;notch=10;arrowHead=10;dx2=15;dy2=20;",
		"tailedNotchedArrow":   "dy1=6;dx1=15;notch=15;arrowHead=10;dx2=15;dy2=20;",
		"stripedArrow":         "dy=0.4;dx=25;notch=40;",
		"jumpInArrow":          "dy=8;dx=25;arrowHead=35;",
		"uTurnArrow":           "dy=8;arrowHead=30;dx2=15;",
	}
	var cols []string
	for _, v := range variants {
		if v == "" {
			v = "sizes"
		}
		cols = append(cols, strings.TrimSuffix(v, ";"))
	}
	d.grid(0, names, cols, 90, 80, 70, 60, func(r, c string) string {
		v := c + ";"
		if c == "sizes" {
			v = ""
		}
		return a + r + ";" + sizes[r] + v
	})
	d.write(filepath.Join(dir, "arrows2_variants.drawio"))

	// the wedge edges, and bounded labels of arrow and twoWayArrow
	d = &doc{}
	for i, n := range []string{"wedgeArrow", "wedgeArrowDashed", "wedgeArrowDashed2"} {
		x := 20 + float64(i)*260
		base := "html=1;shape=mxgraph.arrows2." + n + ";bendable=0;"
		d.edge(base+"startWidth=50;fillColor=#dae8fc;", [2]float64{x, 80}, [2]float64{x + 200, 20})
		d.edge(base+"startWidth=20;stepSize=5;strokeWidth=2;strokeColor=#b85450;", [2]float64{x, 140}, [2]float64{x + 100, 200}, [2]float64{x + 220, 140})
		d.edge(base, [2]float64{x + 40, 220}, [2]float64{x + 40, 330})
		d.edge(base+"startWidth=10;stepSize=0;", [2]float64{x + 120, 220}, [2]float64{x + 200, 330})
		d.text(x, 340, 220, n)
	}
	y = 380
	lbl := "html=1;shape=mxgraph.arrows2."
	for i, s := range []string{
		"arrow;dy=0.6;dx=40;notch=0;boundedLbl=1;whiteSpace=wrap;",
		"arrow;dy=0.6;dx=40;notch=0;boundedLbl=1;flipH=1;whiteSpace=wrap;",
		"arrow;dy=0.6;dx=40;notch=0;boundedLbl=1;direction=south;whiteSpace=wrap;",
		"arrow;dy=0.6;dx=40;notch=0;boundedLbl=1;direction=north;flipV=1;whiteSpace=wrap;",
		"arrow;dy=0.6;dx=40;notch=0;boundedLbl=1;direction=west;whiteSpace=wrap;",
		"twoWayArrow;dy=0.6;dx=35;boundedLbl=1;whiteSpace=wrap;",
		"twoWayArrow;dy=0.6;dx=35;boundedLbl=1;direction=south;whiteSpace=wrap;",
		"arrow;dy=0.6;dx=40;notch=0;whiteSpace=wrap;",
	} {
		x := 20 + float64(i)*130
		w, h := 110.0, 70.0
		if strings.Contains(s, "south") || strings.Contains(s, "north") {
			w, h = 70, 110
		}
		d.cell(lbl+s, x+(110-w)/2, y, w, h, "A long label that wraps in the arrow")
		d.text(x, y+120, 120, strings.TrimSuffix(s, "whiteSpace=wrap;"))
	}
	d.write(filepath.Join(dir, "arrows2_edges.drawio"))
}

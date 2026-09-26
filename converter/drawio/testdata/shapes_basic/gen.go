//go:build ignore

// Gen writes the test diagrams of draw.io's basic shape library
// (mxgraph.basic.*, shapes_basic.go): defaults.drawio has every shape with
// the defaults of its constructor and with the style of draw.io's
// sidebar, params.drawio the parameters the shapes read, transforms.drawio
// rotations, directions, flips, fills and strokes. Each shape has a label
// below naming its style (draw.io draws the labels, which help to compare
// a draw.io export with the BDF rendering shape by shape).
//
//	go run ./converter/drawio/testdata/shapes_basic/gen.go converter/drawio/testdata/shapes_basic
//	/Applications/draw.io.app/Contents/MacOS/draw.io -x -f png -s 2 -o ref.png converter/drawio/testdata/shapes_basic/params.drawio
//	go run ./cmd/bdf generate -q converter/drawio/testdata/shapes_basic/params.drawio params.bdf
//	node test/render.mjs params.bdf out/ 2.6667
package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

type item struct {
	style string
	w, h  float64
	value string
}

type doc struct {
	b      strings.Builder
	id     int
	cols   int
	cw, ch float64
	n      int
}

func (d *doc) next() string { d.id++; return fmt.Sprint("c", d.id) }

func (d *doc) cell(style string, x, y, w, h float64, value string) {
	fmt.Fprintf(&d.b, `        <mxCell id="%s" value="%s" style="%s" vertex="1" parent="1">
          <mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/>
        </mxCell>
`, d.next(), html.EscapeString(value), html.EscapeString(style), x, y, w, h)
}

// add places a shape in the next grid slot with its style below.
func (d *doc) add(it item) {
	col, row := d.n%d.cols, d.n/d.cols
	d.n++
	x := 20 + float64(col)*d.cw
	y := 20 + float64(row)*d.ch
	w, h := it.w, it.h
	if w == 0 {
		w, h = d.cw-50, d.ch-60
	}
	d.cell(it.style, x+(d.cw-30-w)/2, y+(d.ch-60-h)/2, w, h, it.value)
	// spaces let the long styles wrap
	lbl := strings.ReplaceAll(strings.TrimPrefix(it.style, "shape=mxgraph.basic."), ";", "; ")
	d.cell("text;html=1;align=center;verticalAlign=top;fontSize=7;spacing=0;whiteSpace=wrap;", x-5, y+d.ch-52, d.cw-20, 24, lbl)
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

const p = "shape=mxgraph.basic."

// names are the shapes of mxBasic.js.
var names = []string{
	"cross2", "rectCallout", "roundRectCallout", "wave2", "octagon2", "isocube", "acute_triangle",
	"obtuse_triangle", "drop", "cone2", "pyramid", "4_point_star_2", "diag_snip_rect", "diag_round_rect",
	"corner_round_rect", "plaque", "frame", "plaque_frame", "rounded_frame", "frame_corner", "diag_stripe",
	"donut", "layered_rect", "button", "shaded_button", "pie", "arc", "partConcEllipse", "numberedEntryVert",
	"bendingArch", "three_corner_round_rect", "polygon", "patternFillRect",
}

// sidebar are the styles of draw.io's sidebar (Sidebar-Basic.js) and, for
// the shapes it does not have, typical values.
var sidebar = []string{
	"cross2;dx=10;", "rectCallout;dx=30;dy=15;boundedLbl=1;", "roundRectCallout;dx=30;dy=15;size=5;boundedLbl=1;",
	"wave2;dy=0.3;", "octagon2;dx=15;", "isocube;isoAngle=15;", "acute_triangle;dx=0.5;", "obtuse_triangle;dx=0.25;",
	"drop;", "cone2;dx=0.5;dy=0.9;", "pyramid;dx1=0.4;dx2=0.6;dy1=0.9;dy2=0.8;", "4_point_star_2;dx=0.8;",
	"diag_snip_rect;dx=6;", "diag_round_rect;dx=6;", "corner_round_rect;dx=6;", "plaque;dx=6;", "frame;dx=10;",
	"plaque_frame;dx=10;", "rounded_frame;dx=10;", "frame_corner;dx=10;", "diag_stripe;dx=10;", "donut;dx=25;",
	"layered_rect;dx=10;boundedLbl=1;", "button;dx=10;", "shaded_button;dx=10;fillColor=#E6E6E6;strokeColor=none;",
	"pie;startAngle=0.2;endAngle=0.9;", "arc;startAngle=0.3;endAngle=0.1;",
	"partConcEllipse;startAngle=0.25;endAngle=0.1;arcWidth=0.5;", "numberedEntryVert;dy=25;",
	"bendingArch;startAngle=0.6;endAngle=0.4;arcWidth=0.3;", "three_corner_round_rect;dx=6;",
	"polygon;polyCoords=[[0.25,0],[0.75,0],[1,0.25],[1,0.75],[0.75,1],[0.25,1],[0,0.75],[0,0.25]];polyline=0;",
	"patternFillRect;fillStyle=diag;step=5;fillStrokeWidth=0.2;fillStrokeColor=#dddddd;",
}

func main() {
	dir := os.Args[1]

	// every shape with its constructor's defaults, then the sidebar's style
	defaults := &doc{cols: 7, cw: 140, ch: 130}
	for _, n := range names {
		defaults.add(item{style: p + n + ";"})
	}
	for _, s := range sidebar {
		defaults.add(item{style: p + s + "whiteSpace=wrap;html=1;", value: "Text"})
	}
	defaults.write(filepath.Join(dir, "defaults.drawio"))

	// the parameters
	params := &doc{cols: 7, cw: 140, ch: 130}
	for _, s := range []string{
		"cross2;dx=30;", "cross2;dx=100;",
		"rectCallout;dx=80;dy=30;", "rectCallout;dx=5;dy=40;",
		"roundRectCallout;dx=30;dy=15;size=20;", "roundRectCallout;dx=0;dy=20;size=10;", "roundRectCallout;dx=200;dy=10;size=40;",
		"wave2;dy=0;", "wave2;dy=0.6;",
		"octagon2;dx=5;", "octagon2;dx=100;",
		"isocube;isoAngle=30;", "isocube;isoAngle=90;",
		"acute_triangle;dx=0.1;", "acute_triangle;dx=1;",
		"obtuse_triangle;dx=0.75;", "obtuse_triangle;dx=0;",
		"cone2;dx=0.2;dy=0.6;", "cone2;dx=0.8;dy=1;",
		"pyramid;dx1=0.2;dx2=0.8;dy1=0.5;dy2=0.3;",
		"4_point_star_2;dx=0.3;", "4_point_star_2;dx=1;",
		"diag_snip_rect;dx=20;", "diag_round_rect;dx=20;", "corner_round_rect;dx=20;", "plaque;dx=20;",
		"frame;dx=25;", "plaque_frame;dx=5;", "rounded_frame;dx=5;", "frame_corner;dx=30;",
		"diag_stripe;dx=3;", "diag_stripe;dx=40;",
		"donut;dx=10;", "donut;dx=100;",
		"layered_rect;dx=20;", "button;dx=5;fillColor=#dae8fc;", "shaded_button;dx=20;fillColor=#dae8fc;",
		"pie;startAngle=0;endAngle=0.5;", "pie;startAngle=0.5;endAngle=0;", "pie;startAngle=0.25;endAngle=0.75;", "pie;startAngle=0.9;endAngle=0.1;",
		"arc;startAngle=0;endAngle=0.5;", "arc;startAngle=0.5;endAngle=0;", "arc;startAngle=0.25;endAngle=0.75;",
		"arc;startAngle=0.3;endAngle=0.1;startArrow=classic;endArrow=block;",
		"arc;startAngle=0;endAngle=0.5;endArrow=oval;startArrow=diamond;startFill=0;",
		"arc;startAngle=0.9;endAngle=0.6;strokeWidth=3;endArrow=open;endSize=10;startArrow=classicThin;startFillColor=#ff0000;",
		"partConcEllipse;startAngle=0;endAngle=0.5;arcWidth=0.5;", "partConcEllipse;startAngle=0.6;endAngle=0.2;arcWidth=0.2;",
		"numberedEntryVert;dy=40;", "numberedEntryVert;dy=200;",
		"bendingArch;startAngle=0.25;endAngle=0.1;arcWidth=0.5;",
		"three_corner_round_rect;dx=20;",
		"polygon;polyCoords=[[0.25,0],[0.75,0],[1,0.25],[1,0.75],[0.75,1],[0.25,1],[0,0.75],[0,0.25]];polyline=1;fillColor=none;",
		`polygon;polyCoords=[[0,1],[0.5,0],[1,1]];polyCurves=[["Q",0.5,0.5],["Q",1,0],["Q",0.5,1.2]];`,
		`polygon;polyCoords=[[0,0],[1,0],[1,1],[0,1]];polyCurves=[[],["Q",1.3,0.5]];polyline=1;`,
		"patternFillRect;fillStyle=diagRev;step=5;fillStrokeWidth=0.2;fillStrokeColor=#dddddd;",
		"patternFillRect;fillStyle=diagGrid;step=8;fillStrokeWidth=1;fillStrokeColor=#999999;",
		"patternFillRect;fillStyle=grid;step=10;fillStrokeWidth=2;fillStrokeColor=#ff0000;",
		"patternFillRect;fillStyle=hor;step=5;top=0;bottom=0;",
		"patternFillRect;fillStyle=vert;step=5;left=0;right=0;fillColor=#fff2cc;",
		"patternFillRect;fillStyle=diag;step=5;strokeColor=none;",
	} {
		params.add(item{style: p + s})
	}
	// labels confined by boundedLbl
	for _, s := range []string{
		"rectCallout;dx=30;dy=30;boundedLbl=1;verticalAlign=bottom;",
		"roundRectCallout;dx=30;dy=30;size=10;boundedLbl=1;verticalAlign=bottom;",
		"layered_rect;dx=20;boundedLbl=1;align=right;verticalAlign=bottom;",
		"layered_rect;dx=20;boundedLbl=1;align=right;verticalAlign=bottom;flipH=1;",
	} {
		params.add(item{style: p + s + "whiteSpace=wrap;html=1;", value: "Label"})
	}
	params.write(filepath.Join(dir, "params.drawio"))

	// rotations, directions, flips, fills and strokes
	transforms := &doc{cols: 8, cw: 130, ch: 120}
	for _, s := range []string{
		"rectCallout;dx=30;dy=15;", "cone2;dx=0.2;dy=0.8;", "pie;startAngle=0.2;endAngle=0.9;",
		"arc;startAngle=0.3;endAngle=0.1;startArrow=classic;endArrow=block;", "layered_rect;dx=10;",
		`polygon;polyCoords=[[0,0],[1,0.3],[0.4,1]];polyCurves=[["Q",0.8,0]];`,
		"patternFillRect;fillStyle=diag;step=6;right=0;", "obtuse_triangle;dx=0.25;",
	} {
		for _, t := range []string{"", "rotation=30;", "direction=south;", "direction=north;", "direction=west;",
			"flipH=1;", "flipV=1;", "direction=south;flipH=1;rotation=20;"} {
			transforms.add(item{style: p + s + t})
		}
	}
	for _, s := range []string{"frame_corner;dx=10;", "donut;dx=15;", "button;dx=10;", "roundRectCallout;dx=30;dy=15;size=5;"} {
		for _, t := range []string{
			"fillColor=#dae8fc;gradientColor=#7ea6e0;", "fillColor=#f8cecc;gradientColor=#ea6b66;gradientDirection=east;",
			"strokeWidth=3;strokeColor=#b85450;", "dashed=1;strokeWidth=2;", "opacity=50;fillColor=#d5e8d4;",
			"shadow=1;", "fillColor=none;", "strokeColor=none;fillColor=#e1d5e7;",
		} {
			transforms.add(item{style: p + s + t})
		}
	}
	transforms.write(filepath.Join(dir, "transforms.drawio"))
}

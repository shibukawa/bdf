//go:build ignore

// Gen writes the shape test diagrams of this directory (except
// extent.drawio): every shape, style variants, UML shapes and tables,
// markers, edge shapes, mermaid shapes and mxgraph.basic.rect, each with
// a label naming its style (draw.io draws the labels, which help to
// compare a draw.io export with the BDF rendering shape by shape).
//
//	go run ./converter/drawio/testdata/shapes/gen.go converter/drawio/testdata/shapes
//	/Applications/draw.io.app/Contents/MacOS/draw.io -x -f png -s 1.5 -b 10 -o ref.png converter/drawio/testdata/shapes/basic.drawio
//	go run ./cmd/bdf generate converter/drawio/testdata/shapes/basic.drawio basic.bdf
//	node test/render.mjs basic.bdf out/ 2
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
	label string
}

type doc struct {
	b    strings.Builder
	id   int
	cols int
	cw   float64
	ch   float64
	n    int
}

func (d *doc) next() string { d.id++; return fmt.Sprint("c", d.id) }

func (d *doc) cell(style string, x, y, w, h float64, value string, parent string) string {
	id := d.next()
	if parent == "" {
		parent = "1"
	}
	fmt.Fprintf(&d.b, `        <mxCell id="%s" value="%s" style="%s" vertex="1" parent="%s">
          <mxGeometry x="%g" y="%g" width="%g" height="%g" as="geometry"/>
        </mxCell>
`, id, html.EscapeString(value), html.EscapeString(style), parent, x, y, w, h)
	return id
}

func (d *doc) label(x, y, w float64, text string) {
	d.cell("text;html=1;align=center;verticalAlign=top;fontSize=8;spacing=0;", x, y, w, 12, text, "")
}

// add places a shape in the next grid slot with its label below.
func (d *doc) add(it item) {
	col, row := d.n%d.cols, d.n/d.cols
	d.n++
	x := 20 + float64(col)*d.cw
	y := 20 + float64(row)*d.ch
	w, h := it.w, it.h
	if w == 0 {
		w, h = d.cw-40, d.ch-50
	}
	ox := (d.cw - 40 - w) / 2
	d.cell(it.style, x+ox, y, w, h, "", "")
	lbl := it.label
	if lbl == "" {
		lbl = it.style
	}
	d.label(x-10, y+d.ch-45, d.cw-20, lbl)
}

func (d *doc) edge(style string, pts ...[2]float64) {
	id := d.next()
	fmt.Fprintf(&d.b, `        <mxCell id="%s" value="" style="%s" edge="1" parent="1">
          <mxGeometry relative="1" as="geometry">
            <mxPoint x="%g" y="%g" as="sourcePoint"/>
            <mxPoint x="%g" y="%g" as="targetPoint"/>
`, id, html.EscapeString(style), pts[0][0], pts[0][1], pts[len(pts)-1][0], pts[len(pts)-1][1])
	if len(pts) > 2 {
		d.b.WriteString("            <Array as=\"points\">\n")
		for _, p := range pts[1 : len(pts)-1] {
			fmt.Fprintf(&d.b, "              <mxPoint x=\"%g\" y=\"%g\"/>\n", p[0], p[1])
		}
		d.b.WriteString("            </Array>\n")
	}
	d.b.WriteString("          </mxGeometry>\n        </mxCell>\n")
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

func newDoc(cols int, cw, ch float64) *doc { return &doc{cols: cols, cw: cw, ch: ch} }

const base = "whiteSpace=wrap;html=1;"

func main() {
	dir := os.Args[1]

	// every vertex shape with its default style
	basic := newDoc(8, 150, 120)
	for _, s := range []string{
		"rounded=0;", "rounded=1;", "ellipse;", "shape=doubleEllipse;", "rhombus;", "triangle;", "shape=cloud;", "shape=cylinder;",
		"shape=actor;", "shape=hexagon;perimeter=hexagonPerimeter2;", "shape=process;", "shape=process;rounded=1;", "shape=document;", "shape=parallelogram;", "shape=trapezoid;", "shape=step;",
		"shape=plus;", "shape=cube;", "shape=note;", "shape=note2;", "shape=card;", "shape=tape;", "shape=dataStorage;", "shape=internalStorage;",
		"shape=manualInput;", "shape=loopLimit;", "shape=offPageConnector;", "shape=delay;", "shape=display;", "shape=or;", "shape=xor;", "shape=sortShape;",
		"shape=collate;", "shape=singleArrow;", "shape=doubleArrow;", "shape=cross;", "shape=corner;", "shape=tee;", "shape=crossbar;", "shape=datastore;",
		"shape=cylinder2;", "shape=cylinder3;", "shape=folder;", "shape=callout;", "shape=wedgeCallout;", "shape=curlyBracket;", "shape=partialRectangle;", "shape=lineEllipse;",
		"shape=orEllipse;", "shape=sumEllipse;", "shape=switch;", "shape=message;", "shape=startState;", "shape=endState;", "shape=associativeEntity;", "shape=isoRectangle;",
		"shape=isoCube;", "shape=isoCube2;", "shape=dimension;", "shape=tapeData;", "shape=ext;double=1;", "shape=waypoint;", "shape=smileyFace;", "shape=zigzag;",
		"line;", "shape=parallelMarker;", "shape=transparent;", "shape=process2;", "shape=lineEllipse;line=vertical;", "shape=smileyFace;smileyType=sad;", "shape=partialRectangle;top=0;right=0;", "shape=zigzag;rounded=1;",
	} {
		basic.add(item{style: s + base})
	}
	basic.write(filepath.Join(dir, "basic.drawio"))

	// style variants of the shapes with size and rounded styles
	variants := newDoc(11, 110, 100)
	vs := []string{"", "rounded=1;", "direction=south;", "direction=north;", "direction=west;", "flipH=1;", "flipV=1;", "rotation=30;",
		"dashed=1;strokeWidth=3;", "fillColor=#dae8fc;gradientColor=#7ea6e0;opacity=70;strokeColor=#6c8ebf;"}
	sizeVariant := map[string]string{
		"process": "size=0.3;", "document": "size=0.5;", "parallelogram": "fixedSize=1;size=10;", "trapezoid": "size=0.4;",
		"step": "fixedSize=1;size=10;", "hexagon": "size=0.1;", "cube": "size=10;darkOpacity=0.3;darkOpacity2=-0.3;", "note": "size=10;darkOpacity=0.2;",
		"card": "size=10;", "tape": "size=0.2;", "callout": "size=10;position=0.2;position2=0.1;base=10;", "folder": "tabWidth=30;tabHeight=10;tabPosition=left;",
		"singleArrow": "arrowWidth=0.6;arrowSize=0.4;", "cylinder3": "size=5;lid=0;", "manualInput": "size=10;", "offPageConnector": "size=0.2;",
		"loopLimit": "size=10;", "corner": "dx=40;dy=10;", "tee": "dx=10;dy=40;", "internalStorage": "dx=40;dy=10;", "wedgeCallout": "tipX=0.6;tipY=0.3;base=10;",
		"datastore": "strokeWidth=2;", "cylinder": "size=0.3;", "isoCube2": "isoAngle=30;", "dataStorage": "fixedSize=1;size=5;", "curlyBracket": "size=0.2;rounded=1;",
		"umlFrame": "width=40;height=20;swimlaneFillColor=#ffe6cc;", "delay": "", "display": "size=0.1;", "doubleArrow": "arrowWidth=0.5;arrowSize=0.1;",
		"cross": "size=0.5;", "cylinder2": "size=5;", "swimlane": "startSize=10;", "rhombus": "double=1;", "ellipse": "glass=1;",
	}
	for _, name := range []string{"process", "document", "parallelogram", "trapezoid", "step", "hexagon", "cube", "note", "card", "tape",
		"callout", "folder", "singleArrow", "cylinder3", "manualInput", "offPageConnector", "loopLimit", "corner", "tee", "internalStorage",
		"wedgeCallout", "datastore", "cylinder", "isoCube2", "dataStorage", "curlyBracket", "umlFrame", "display", "doubleArrow", "cross",
		"cylinder2", "swimlane", "rhombus", "ellipse"} {
		for _, v := range vs {
			variants.add(item{style: "shape=" + name + ";" + v + base, label: name + " " + v})
		}
		variants.add(item{style: "shape=" + name + ";" + sizeVariant[name] + base, label: name + " " + sizeVariant[name]})
	}
	variants.write(filepath.Join(dir, "variants.drawio"))

	// UML shapes, swimlanes and a table
	uml := newDoc(6, 170, 160)
	for _, it := range []item{
		{style: "shape=umlActor;verticalLabelPosition=bottom;verticalAlign=top;", w: 30, h: 60},
		{style: "shape=umlBoundary;", w: 100, h: 80},
		{style: "ellipse;shape=umlEntity;", w: 80, h: 80},
		{style: "shape=umlDestroy;", w: 30, h: 30},
		{style: "ellipse;shape=umlControl;", w: 70, h: 80},
		{style: "shape=umlLifeline;perimeter=lifelinePerimeter;", w: 100, h: 110},
		{style: "shape=umlLifeline;participant=umlActor;perimeter=lifelinePerimeter;", w: 30, h: 110},
		{style: "shape=umlLifeline;participant=umlControl;", w: 40, h: 110},
		{style: "shape=umlLifeline;lifelineMirror=1;size=30;rounded=1;", w: 100, h: 110},
		{style: "shape=umlFrame;", w: 130, h: 100},
		{style: "shape=umlState;rounded=1;", w: 130, h: 80},
		{style: "shape=umlState;rounded=1;umlStateConnection=connPointRefExit;umlStateSymbol=collapseState;", w: 130, h: 80},
		{style: "shape=umlState;rounded=1;umlStateConnection=connPointRefEntry;absoluteArcSize=1;arcSize=10;", w: 130, h: 80},
		{style: "shape=component;", w: 120, h: 80},
		{style: "shape=module;jettyWidth=10;jettyHeight=8;", w: 120, h: 80},
		{style: "shape=lollipop;", w: 20, h: 60},
		{style: "shape=requires;", w: 30, h: 60},
		{style: "shape=requiredInterface;", w: 20, h: 60},
		{style: "shape=providedRequiredInterface;", w: 40, h: 40},
		{style: "swimlane;", w: 130, h: 110},
		{style: "swimlane;horizontal=0;", w: 130, h: 110},
		{style: "swimlane;rounded=1;fillColor=#dae8fc;swimlaneFillColor=#ffffff;", w: 130, h: 110},
		{style: "swimlane;swimlaneLine=0;separatorColor=#ff0000;startSize=30;", w: 130, h: 110},
		{style: "swimlane;swimlaneHead=0;fillColor=#f8cecc;", w: 130, h: 110},
		{style: "swimlane;swimlaneBody=0;fillColor=#d5e8d4;swimlaneFillColor=#fff2cc;", w: 130, h: 110},
		{style: "swimlane;horizontal=0;rounded=1;footerSize=20;fillColor=#e1d5e7;", w: 130, h: 110},
		{style: "swimlane;startSize=0;fixedHeader=0;", w: 130, h: 110},
		{style: "shape=seqQueue;", w: 100, h: 40},
		{style: "shape=seqCollections;", w: 100, h: 40},
		{style: "shape=seqDatabase;", w: 120, h: 80},
		{style: "shape=seqActorStick;", w: 60, h: 80},
		{style: "shape=seqBoundary;", w: 120, h: 60},
	} {
		uml.add(it)
	}
	// a 3x3 table as the sidebar makes it
	n := uml.n
	x0 := 20 + float64(n%uml.cols)*uml.cw
	y0 := 20 + float64(n/uml.cols)*uml.ch
	tbl := uml.cell("shape=table;startSize=0;container=1;collapsible=0;childLayout=tableLayout;", x0, y0, 150, 90, "", "")
	for r := 0; r < 3; r++ {
		row := uml.cell("shape=tableRow;horizontal=0;startSize=0;swimlaneHead=0;swimlaneBody=0;strokeColor=inherit;top=0;left=0;bottom=0;right=0;collapsible=0;dropTarget=0;fillColor=none;points=[[0,0.5],[1,0.5]];portConstraint=eastwest;", 0, float64(r)*30, 150, 30, "", tbl)
		for k := 0; k < 3; k++ {
			fill := "none"
			if r == 1 && k == 1 {
				fill = "#ffe6cc"
			}
			uml.cell("shape=partialRectangle;html=1;whiteSpace=wrap;connectable=0;strokeColor=inherit;overflow=hidden;fillColor="+fill+";top=0;left=0;bottom=0;right=0;pointerEvents=1;", float64(k)*50, 0, 50, 30, "", row)
		}
	}
	uml.n++
	uml.label(x0, y0+uml.ch-45, 150, "table")
	// a table with a title and rounded corners
	n = uml.n
	x0 = 20 + float64(n%uml.cols)*uml.cw
	y0 = 20 + float64(n/uml.cols)*uml.ch
	tbl = uml.cell("shape=table;startSize=20;container=1;collapsible=0;childLayout=tableLayout;fixedRows=1;rowLines=0;fontStyle=1;", x0, y0, 150, 110, "", "")
	for r := 0; r < 3; r++ {
		row := uml.cell("shape=tableRow;horizontal=0;startSize=0;swimlaneHead=0;swimlaneBody=0;fillColor=none;collapsible=0;dropTarget=0;points=[[0,0.5],[1,0.5]];portConstraint=eastwest;top=0;left=0;right=0;bottom=1;", 0, 20+float64(r)*30, 150, 30, "", tbl)
		for k := 0; k < 2; k++ {
			uml.cell("shape=partialRectangle;connectable=0;fillColor=none;top=0;left=0;bottom=0;right=0;editable=1;overflow=hidden;", float64(k)*30, 0, 30+float64(k)*90, 30, "", row)
		}
	}
	uml.n++
	uml.label(x0, y0+uml.ch-45, 150, "table title")
	uml.write(filepath.Join(dir, "uml.drawio"))

	// every marker, at both ends
	mk := newDoc(3, 200, 40)
	markers := []string{"classic", "classicThin", "block", "blockThin", "open", "openThin", "oval", "diamond", "diamondThin",
		"baseDash", "doubleBlock", "manyOptional", "dash", "box", "cross", "circle", "circlePlus", "halfCircle", "async", "openAsync",
		"ERone", "ERmandOne", "ERmany", "ERoneToMany", "ERzeroToMany", "ERzeroToOne", "sysMLx", "sysMLLost", "sysMLFound",
		"sysMLPackCont", "sysMLReqInt", "sysMLProvInt", "mermaidExtension", "mermaidDiamond"}
	for i, m := range markers {
		y := 20 + float64(i)*40
		mk.edge("endArrow="+m+";startArrow="+m+";html=1;", [2]float64{60, y}, [2]float64{190, y})
		mk.edge("endArrow="+m+";startArrow="+m+";startFill=0;endFill=0;html=1;", [2]float64{240, y}, [2]float64{370, y})
		mk.edge("endArrow="+m+";startArrow="+m+";startSize=12;endSize=12;strokeWidth=2;strokeColor=#0000ff;html=1;", [2]float64{420, y}, [2]float64{550, y})
		mk.edge("endArrow="+m+";html=1;endFill=1;", [2]float64{600, y + 15}, [2]float64{660, y - 10})
		mk.label(-40, y-6, 100, m)
	}
	mk.write(filepath.Join(dir, "markers.drawio"))

	// edge shapes
	es := newDoc(3, 200, 100)
	row := func(i int) float64 { return 30 + float64(i)*90 }
	for i, st := range []string{
		"shape=arrow;edgeStyle=none;html=1;fillColor=#dae8fc;",
		"shape=flexArrow;endArrow=classic;html=1;rounded=0;",
		"shape=flexArrow;endArrow=classic;startArrow=classic;html=1;rounded=0;fillColor=#d5e8d4;endWidth=30;endSize=5;width=6;",
		"shape=flexArrow;endArrow=classic;html=1;rounded=1;",
		"shape=flexArrow;endArrow=none;startArrow=block;html=1;curved=1;",
		"shape=link;html=1;rounded=0;",
		"shape=link;html=1;rounded=1;width=8;strokeWidth=2;",
		"shape=filledEdge;html=1;strokeWidth=8;fillColor=#ffcc00;endArrow=none;rounded=0;",
		"shape=pipe;html=1;strokeWidth=2;fillColor=#99ccff;endArrow=none;rounded=0;",
		"shape=wire;html=1;dashed=1;strokeWidth=4;fillColor=#ff0000;endArrow=block;rounded=0;",
		"shape=arrowConnector;html=1;",
	} {
		y := row(i)
		es.edge(st, [2]float64{40, y}, [2]float64{200, y})
		es.edge(st, [2]float64{260, y + 20}, [2]float64{330, y - 20}, [2]float64{420, y + 20}, [2]float64{500, y})
		es.label(520, y-6, 300, st)
	}
	es.write(filepath.Join(dir, "edges.drawio"))

	// shapes of mermaid imports
	mm := newDoc(6, 170, 130)
	for _, it := range []item{
		{style: "shape=gitTag;fillColor=#ffffde;", w: 100, h: 24},
		{style: "shape=gitTag;tabSize=16;tabInset=8;holeSize=3;holeColor=#ff0000;", w: 100, h: 40},
		{style: "shape=gitMergeCommit;fillColor=#9370db;", w: 30, h: 30},
		{style: "shape=gitCherryPick;fillColor=#333333;", w: 40, h: 40},
		{style: "shape=mindmapBang;", w: 120, h: 70},
		{style: "shape=ishikawaHead;", w: 60, h: 70},
		{style: "shape=mermaidOdd;", w: 120, h: 50},
		{style: "shape=mermaidBlockArrow;dirs=right;", w: 120, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=left;", w: 120, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=up;", w: 60, h: 80},
		{style: "shape=mermaidBlockArrow;dirs=down;", w: 60, h: 80},
		{style: "shape=mermaidBlockArrow;dirs=x;", w: 120, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=y;", w: 60, h: 80},
		{style: "shape=mermaidBlockArrow;dirs=right,up;", w: 80, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=x,y;", w: 100, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=left,down;", w: 80, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=right,left,up;", w: 120, h: 60},
		{style: "shape=mermaidBlockArrow;dirs=up,down,right;", w: 80, h: 60},
	} {
		mm.add(it)
	}
	mm.edge("shape=mermaidSankeyLink;fillColor=#dae8fc;gradientColor=#7ea6e0;gradientDirection=east;strokeColor=none;width=20;endArrow=none;", [2]float64{40, 450}, [2]float64{300, 520})
	mm.edge("shape=mermaidSankeyLink;fillColor=#f8cecc;strokeColor=#b85450;width=6;endArrow=none;", [2]float64{40, 560}, [2]float64{300, 500})
	mm.write(filepath.Join(dir, "mermaid.drawio"))

	// mxgraph.basic.rect: corner styles, sides and outlines
	r2 := newDoc(8, 130, 100)
	sides := []string{"", "top=0;", "right=0;", "bottom=0;", "left=0;", "top=0;bottom=0;", "left=0;right=0;",
		"top=0;left=0;", "right=0;bottom=0;", "top=0;right=0;bottom=0;", "left=0;top=0;right=0;", "left=0;bottom=0;right=0;",
		"top=0;left=0;bottom=0;", "right=0;left=0;top=0;", "top=0;right=0;", "bottom=0;left=0;"}
	for _, rs := range []string{"square", "rounded", "snip", "invRound", "fold"} {
		for _, outline := range []string{"single", "double", "frame"} {
			for i, sd := range sides {
				if outline != "double" && i%3 != 0 {
					continue
				}
				r2.add(item{style: "shape=mxgraph.basic.rect;rectStyle=" + rs + ";size=15;rectOutline=" + outline + ";indent=5;fillColor2=#ffe6cc;" + sd + base,
					label: rs + " " + outline + " " + sd, w: 90, h: 60})
			}
		}
	}
	for _, st := range []string{
		"topLeftStyle=rounded;topRightStyle=snip;bottomRightStyle=invRound;bottomLeftStyle=fold;size=20;",
		"topLeftStyle=rounded;topRightStyle=snip;bottomRightStyle=invRound;bottomLeftStyle=fold;size=20;rectOutline=double;indent=6;",
		"topLeftStyle=rounded;topRightStyle=snip;bottomRightStyle=invRound;bottomLeftStyle=fold;size=20;rectOutline=frame;indent=6;fillColor=#dae8fc;",
		"rectStyle=rounded;absoluteCornerSize=0;size=30;rectOutline=double;indent=10;",
		"rectStyle=snip;fillColor2=#d5e8d4;gradientColor2=#82b366;rectOutline=double;indent=8;size=12;",
		"rectStyle=fold;fillColor=none;size=12;",
		"rectStyle=rounded;dashed=1;size=12;strokeWidth=3;",
		"rectStyle=invRound;size=12;direction=south;rotation=10;",
	} {
		r2.add(item{style: "shape=mxgraph.basic.rect;" + st + base, label: st, w: 90, h: 60})
	}
	r2.write(filepath.Join(dir, "rect2.drawio"))
}

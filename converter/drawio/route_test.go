package drawio

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The edge routes are checked against draw.io's own rendering: every
// testdata/route/NAME.drawio has a NAME.svg exported with the draw.io
// desktop CLI (31.4.4):
//
//	for f in testdata/route/*.drawio; do
//		draw.io -x -f svg -o "${f%.drawio}.svg" "$f"
//	done
//
// The export draws each cell in a <g data-cell-id="…">; an edge's line is
// its first <path fill="none">. The export moves the drawing to its own
// origin, so every diagram has a plain rectangle with the id "ref" whose
// drawn position gives the translation. Edges have no arrows (markers
// shorten the line) and rounded=0.

// routeTolerance is the largest distance allowed between a computed point
// and draw.io's (which rounds to two decimals).
const routeTolerance = 0.5

// routeTestStencils are the sizes of the draw.io stencils the diagrams use
// (from draw.io's stencil libraries), standing in for the stencil lookup.
var routeTestStencils = map[string]*stencil{
	"mxgraph.electrical.transistors.n-channel_jfet_1":  {w0: 100, h0: 110, aspect: "fixed"},
	"mxgraph.electrical.electro-mechanical.relay_pole": {w0: 100, h0: 5, aspect: "fixed"},
	"mxgraph.basic.4_point_star":                       {w0: 92, h0: 92, aspect: "variable"},
}

// routeKnownDeviations are edges draw.io routes differently on purpose:
// the reason is logged instead of a failure.
var routeKnownDeviations = map[string]string{
	// mxEdgeStyle.Loop leaves y at 0 in view coordinates for north and
	// south loops through a control point beside the vertex, so draw.io's
	// result depends on the view's translation (which the export sets from
	// the drawing's bounds); here the translation is always 0.
	"loop/m1": "Loop: y depends on the view translation",
	"loop/m2": "Loop: y depends on the view translation",
}

func TestRouteReference(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "route", "*.drawio"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no test diagrams")
	}
	for _, f := range files {
		name := strings.TrimSuffix(filepath.Base(f), ".drawio")
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			file, err := readFile(data)
			if err != nil {
				t.Fatal(err)
			}
			m := parseModel(file.pages[0].model)
			v := newView(m, nil, func(name string) *stencil { return routeTestStencils[name] })

			svg, err := os.ReadFile(strings.TrimSuffix(f, ".drawio") + ".svg")
			if err != nil {
				t.Fatal(err)
			}
			ref, err := readRouteSVG(svg)
			if err != nil {
				t.Fatal(err)
			}
			refCell := m.cells["ref"]
			refBox, ok := ref.rects["ref"]
			if refCell == nil || refCell.geo == nil || !ok {
				t.Fatal(`the diagram needs a rectangle with id "ref"`)
			}
			off := point{refBox.x - refCell.geo.x, refBox.y - refCell.geo.y}

			edges := 0
			worst := 0.0
			for _, c := range m.ordered {
				if !c.edge {
					continue
				}
				edges++
				var got []point
				st := v.state(c)
				if st != nil {
					got = getWaypoints(st.edgePoints())
				}
				var want []point
				if pts, ok := ref.paths[c.id]; ok {
					want = []point{}
					for _, p := range pts {
						want = append(want, point{p.x - off.x, p.y - off.y})
					}
				}
				// end markers shorten the painted line: those ends are not compared
				if st != nil && len(got) == len(want) && len(got) > 0 {
					if st.style.has("startArrow") {
						got[0] = want[0]
					}
					if st.style.has("endArrow") {
						got[len(got)-1] = want[len(want)-1]
					}
				}
				d, ok := comparePoints(got, want)
				if reason, known := routeKnownDeviations[name+"/"+c.id]; known {
					t.Logf("edge %s: known deviation (%s):\n got  %s\n draw.io %s", c.id, reason, fmtPoints(got), fmtPoints(want))
					continue
				}
				worst = math.Max(worst, d)
				if !ok {
					t.Errorf("edge %s (%s):\n got  %s\n want %s", c.id, c.styleStr, fmtPoints(got), fmtPoints(want))
				}
			}
			t.Logf("%d edges, largest deviation %.3f", edges, worst)
		})
	}
}

// comparePoints reports the largest distance between matching points and
// whether the routes match.
func comparePoints(got, want []point) (float64, bool) {
	if len(got) != len(want) {
		return 0, false
	}
	worst := 0.0
	for i := range got {
		d := math.Hypot(got[i].x-want[i].x, got[i].y-want[i].y)
		if !(d <= routeTolerance) {
			return d, false
		}
		worst = math.Max(worst, d)
	}
	return worst, true
}

func fmtPoints(pts []point) string {
	if pts == nil {
		return "(not drawn)"
	}
	var b strings.Builder
	for i, p := range pts {
		if i > 0 {
			b.WriteString(" ")
		}
		fmt.Fprintf(&b, "(%.2f,%.2f)", p.x, p.y)
	}
	return b.String()
}

// routeSVG holds what the tests need from a draw.io SVG export: the
// points of each cell's first unfilled path and the position of each
// cell's first rectangle, in the export's coordinates.
type routeSVG struct {
	paths map[string][]point
	rects map[string]point
}

func readRouteSVG(data []byte) (*routeSVG, error) {
	out := &routeSVG{paths: map[string][]point{}, rects: map[string]point{}}
	type frame struct {
		id      string // innermost cell id
		dx, dy  float64
		rotated bool // under a transform other than a translation
	}
	stack := []frame{{}}
	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.Strict = false
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			top := stack[len(stack)-1]
			attr := map[string]string{}
			for _, a := range tok.Attr {
				attr[a.Name.Local] = a.Value
			}
			if id, ok := attr["data-cell-id"]; ok {
				top.id = id
			}
			if tr, ok := attr["transform"]; ok {
				if dx, dy, err := parseTranslate(tr); err == nil {
					top.dx += dx
					top.dy += dy
				} else {
					top.rotated = true
				}
			}
			stack = append(stack, top)
			if top.id == "" || top.rotated {
				break
			}
			switch tok.Name.Local {
			case "path":
				if _, done := out.paths[top.id]; done || attr["fill"] != "none" {
					break
				}
				pts, err := parsePathPoints(attr["d"])
				if err != nil {
					return nil, fmt.Errorf("cell %s: %v", top.id, err)
				}
				for i := range pts {
					pts[i].x += top.dx
					pts[i].y += top.dy
				}
				out.paths[top.id] = pts
			case "rect":
				if _, done := out.rects[top.id]; done {
					break
				}
				x, _ := strconv.ParseFloat(attr["x"], 64)
				y, _ := strconv.ParseFloat(attr["y"], 64)
				out.rects[top.id] = point{x + top.dx, y + top.dy}
			}
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		}
	}
	return out, nil
}

// parseTranslate reads a transform that is a translation, "translate(x,y)".
func parseTranslate(s string) (dx, dy float64, err error) {
	s = strings.TrimSpace(s)
	args, ok := strings.CutPrefix(s, "translate(")
	if !ok || !strings.HasSuffix(args, ")") {
		return 0, 0, fmt.Errorf("unsupported transform %q", s)
	}
	f := strings.FieldsFunc(strings.TrimSuffix(args, ")"), func(r rune) bool { return r == ',' || r == ' ' })
	if len(f) != 2 {
		return 0, 0, fmt.Errorf("unsupported transform %q", s)
	}
	dx, err = strconv.ParseFloat(f[0], 64)
	if err == nil {
		dy, err = strconv.ParseFloat(f[1], 64)
	}
	return dx, dy, err
}

// parsePathPoints reads the points of a path made of M and L commands.
func parsePathPoints(d string) ([]point, error) {
	f := strings.Fields(d)
	var pts []point
	for i := 0; i < len(f); {
		if f[i] != "M" && f[i] != "L" || i+2 >= len(f) {
			return nil, fmt.Errorf("unsupported path %q", d)
		}
		x, err1 := strconv.ParseFloat(f[i+1], 64)
		y, err2 := strconv.ParseFloat(f[i+2], 64)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("unsupported path %q", d)
		}
		pts = append(pts, point{x, y})
		i += 3
	}
	return pts, nil
}

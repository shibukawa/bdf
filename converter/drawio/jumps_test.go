package drawio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLineJumps checks the crossings found for edges with a jumpStyle: each
// horizontal edge crosses the three vertical ones drawn before it, and
// edges without a jump style get no routed points.
func TestLineJumps(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "edges", "jumps.drawio"))
	if err != nil {
		t.Fatal(err)
	}
	f, err := readFile(data)
	if err != nil {
		t.Fatal(err)
	}
	m := parseModel(f.pages[0].model)
	v := newView(m, nil, nil)
	jumps := func(id string) (n int, pts []routedPoint) {
		st := v.state(m.cells[id])
		for _, p := range st.routedPoints {
			if p.jump {
				n++
			}
		}
		return n, st.routedPoints
	}
	for _, id := range []string{"v1", "v2", "v3"} {
		if n, pts := jumps(id); n != 0 || pts != nil {
			t.Errorf("%s: %d jumps, routed %v", id, n, pts)
		}
	}
	for _, id := range []string{"h1", "h2", "h3"} {
		n, pts := jumps(id)
		if n != 3 || len(pts) != 5 {
			t.Fatalf("%s: %d jumps in %v", id, n, pts)
		}
		for i, x := range []float64{100, 200, 300} {
			if p := pts[i+1]; !p.jump || p.x != x {
				t.Errorf("%s: crossing %d at %v, want x=%v", id, i, p, x)
			}
		}
	}
	// the last edge turns before the third vertical edge: its second segment crosses it
	if n, pts := jumps("h4"); n != 3 || len(pts) != 7 || !pts[5].jump || pts[5].y != 190 {
		t.Errorf("h4: %d jumps in %v", n, pts)
	}
}

// TestLineJumpBudget lays out a grid of edges that all cross: the
// crossings grow with the square of the edges, and past the budget the
// edges are drawn without jumps.
func TestLineJumpBudget(t *testing.T) {
	const n = 1200 // 600 × 600 crossings, more than maxJumps
	var b strings.Builder
	for i := range n {
		x0, y0, x1, y1 := 0, 2*i, 4000, 2*i
		if i%2 == 1 {
			x0, y0, x1, y1 = 2*i, 0, 2*i, 4000
		}
		fmt.Fprintf(&b, `<mxCell id="e%d" style="jumpStyle=arc" edge="1" parent="1"><mxGeometry relative="1" as="geometry">`+
			`<mxPoint x="%d" y="%d" as="sourcePoint"/><mxPoint x="%d" y="%d" as="targetPoint"/></mxGeometry></mxCell>`, i, x0, y0, x1, y1)
	}
	m := parseModel(mustXML(t, limitsModel(b.String())))
	warned := false
	v := newView(m, func(key, _ string, _ ...any) { warned = warned || key == "jumps" }, nil)
	jumps, plain := 0, 0
	for _, st := range v.order {
		if !st.cell.edge {
			continue
		}
		if st.routedPoints == nil {
			plain++
		}
		for _, p := range st.routedPoints {
			if p.jump {
				jumps++
			}
		}
	}
	if !warned || plain == 0 || jumps <= maxJumps || jumps > maxJumps+n {
		t.Errorf("%d jumps, %d edges without, warned %v", jumps, plain, warned)
	}
}

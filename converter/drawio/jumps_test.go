package drawio

import (
	"os"
	"path/filepath"
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

package visio

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func TestSegmentCrossing(t *testing.T) {
	for _, c := range []struct {
		a0, a1, b0, b1 pt
		ok             bool
		at             pt
	}{
		{pt{0, 0}, pt{10, 0}, pt{5, -5}, pt{5, 5}, true, pt{5, 0}},
		{pt{0, 0}, pt{10, 0}, pt{10, -5}, pt{10, 5}, false, pt{}}, // at an end
		{pt{0, 0}, pt{10, 0}, pt{0, 1}, pt{10, 1}, false, pt{}},   // parallel
		{pt{0, 0}, pt{10, 0}, pt{5, 1}, pt{5, 5}, false, pt{}},    // apart
	} {
		at, ok := segmentCrossing(c.a0, c.a1, c.b0, c.b1)
		if ok != c.ok || ok && at.dist(c.at) > 1e-9 {
			t.Errorf("crossing of %v-%v and %v-%v = %v, %v", c.a0, c.a1, c.b0, c.b1, at, ok)
		}
	}
}

func TestJumper(t *testing.T) {
	h := connSeg{pt{0, 0}, pt{10, 0}, true}
	v := connSeg{pt{5, -5}, pt{5, 5}, true}
	conn := func(code, z int) *connPath { return &connPath{code: code, z: z} }
	for _, c := range []struct {
		name         string
		codeA, codeB int
		page         int
		want         int // 0 none, 1 a (horizontal), 2 b (vertical)
	}{
		{"horizontal lines", 0, 0, 1, 1},
		{"vertical lines", 0, 0, 2, 2},
		{"no jumps", 0, 0, 0, 0},
		{"last displayed", 0, 0, 4, 2},
		{"first displayed", 0, 0, 5, 1},
		{"always", 0, 2, 1, 2},
		{"never", 1, 0, 1, 0},
		{"other connector jumps", 3, 0, 2, 2},
		{"neither", 0, 4, 1, 0},
	} {
		a, b := conn(c.codeA, 1), conn(c.codeB, 2)
		got := 0
		switch jumper(a, b, h, v, c.page) {
		case a:
			got = 1
		case b:
			got = 2
		}
		if got != c.want {
			t.Errorf("%s: jumper = %d, want %d", c.name, got, c.want)
		}
	}
}

func TestApplyJumps(t *testing.T) {
	line := func() []*subpath { return []*subpath{{start: pt{0, 0}, segs: []seg{{p: pt{100, 0}}}}} }
	// a gap splits the figure
	out := applyJumps(line(), []jump{{at: pt{50, 0}, width: 6, style: 2, up: true}})
	if len(out) != 2 || out[0].end().dist(pt{47, 0}) > 1e-9 || out[1].start.dist(pt{53, 0}) > 1e-9 {
		t.Errorf("gap = %+v", out)
	}
	// an arc bulges upwards (y down) by half its width
	out = applyJumps(line(), []jump{{at: pt{50, 0}, width: 6, style: 1, up: true}})
	if len(out) != 1 || len(out[0].segs) != 4 || out[0].segs[1].p.dist(pt{50, -3}) > 1e-9 || out[0].end() != (pt{100, 0}) {
		t.Errorf("arc = %+v", out[0].segs)
	}
	// overlapping jumps merge into one
	out = applyJumps(line(), []jump{{at: pt{50, 0}, width: 6, style: 3, up: false}, {at: pt{54, 0}, width: 6, style: 3, up: false}})
	want := []pt{{47, 0}, {47, 3}, {57, 3}, {57, 0}, {100, 0}}
	if len(out[0].segs) != len(want) {
		t.Fatalf("squares = %+v", out[0].segs)
	}
	for i, p := range want {
		if out[0].segs[i].p.dist(p) > 1e-9 {
			t.Errorf("square point %d = %v, want %v", i, out[0].segs[i].p, p)
		}
	}
}

// TestJumpsOfDrawing checks the jumps of the text page of shapes.vsdx:
// four horizontal connectors cross two vertical ones in the page's style
// and their own, and a connector that never jumps crosses a third.
func TestJumpsOfDrawing(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "shapes.vsdx"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	st, _ := f.Stat()
	p, err := ooxml.Open(f, st.Size())
	if err != nil {
		t.Fatal(err)
	}
	d, err := readVSDX(p, func(k, m string) { t.Error(m) })
	if err != nil {
		t.Fatal(err)
	}
	c := &converter{opts: &Options{}, doc: bdf.NewDocument(), d: d, warned: map[string]bool{}, pageNum: map[int]int{}}
	pc := c.newPageCtx(d.pages[1], 612)
	pc.collectJumps()
	styles := map[int]int{}
	var widths []float64
	for s, js := range pc.jumps {
		if len(js) != 2 {
			t.Errorf("shape %d has %d jumps, want 2", s.id, len(js))
		}
		styles[js[0].style]++
		if js[0].style == 5 && js[0].up {
			t.Error("the jump set to go down goes up")
		}
		widths = append(widths, js[0].width)
	}
	if len(pc.jumps) != 4 || styles[1] != 1 || styles[2] != 1 || styles[3] != 1 || styles[5] != 1 {
		t.Errorf("%d connectors jump, styles %v", len(pc.jumps), styles)
	}
	// 2/3 of LineToLineX (1/8 in)
	for _, w := range widths {
		if math.Abs(w-6) > 1e-6 {
			t.Errorf("jump width %g, want 6", w)
		}
	}
}

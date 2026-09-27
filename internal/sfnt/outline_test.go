package sfnt

import (
	"math"
	"os"
	"testing"
)

// pointPen collects every point a glyph outline passes, control points
// included, and counts contours.
type pointPen struct {
	r        Rect
	contours int
	open     bool
	bad      bool
}

func newPointPen() *pointPen {
	return &pointPen{r: Rect{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}}
}

func (p *pointPen) pt(x, y float64) {
	if !p.open {
		p.bad = true
	}
	p.r.XMin, p.r.XMax = math.Min(p.r.XMin, x), math.Max(p.r.XMax, x)
	p.r.YMin, p.r.YMax = math.Min(p.r.YMin, y), math.Max(p.r.YMax, y)
}

func (p *pointPen) MoveTo(x, y float64) {
	if p.open {
		p.bad = true
	}
	p.open = true
	p.contours++
	p.pt(x, y)
}
func (p *pointPen) LineTo(x, y float64)         { p.pt(x, y) }
func (p *pointPen) QuadTo(cx, cy, x, y float64) { p.pt(cx, cy); p.pt(x, y) }
func (p *pointPen) CubeTo(ax, ay, bx, by, x, y float64) {
	p.pt(ax, ay)
	p.pt(bx, by)
	p.pt(x, y)
}
func (p *pointPen) Close() {
	if !p.open {
		p.bad = true
	}
	p.open = false
}

func TestTrueTypeOutlines(t *testing.T) {
	data, err := os.ReadFile("../../fixture/testdata/fonts/DejaVuSans-sub.ttf")
	if err != nil {
		t.Skip(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	o := NewOutlines(f)
	drawn := 0
	for gid := 0; gid < f.NumGlyphs; gid++ {
		p := newPointPen()
		if !o.Outline(uint16(gid), p) {
			t.Fatalf("glyph %d: not read", gid)
		}
		if p.bad || p.open {
			t.Fatalf("glyph %d: contours not opened and closed in pairs", gid)
		}
		want, ok := o.Bounds(uint16(gid))
		if !ok || p.contours == 0 {
			continue
		}
		drawn++
		// the glyph header holds the bounds of the points, control points included
		if p.r != want {
			t.Errorf("glyph %d: points span %v, header says %v", gid, p.r, want)
		}
	}
	if drawn == 0 {
		t.Fatal("no glyph drew anything")
	}
}

func TestCFFOutlines(t *testing.T) {
	data, err := os.ReadFile("testdata/STIXTwoMath-cff-subset.otf")
	if err != nil {
		t.Skip(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	o := NewOutlines(f)
	drawn := 0
	for gid := 0; gid < f.NumGlyphs; gid++ {
		p := newPointPen()
		if !o.Outline(uint16(gid), p) {
			t.Fatalf("glyph %d: not read", gid)
		}
		if p.bad || p.open {
			t.Fatalf("glyph %d: contours not opened and closed in pairs", gid)
		}
		want, ok := o.Bounds(uint16(gid))
		if !ok {
			continue
		}
		drawn++
		// the control points enclose the curves
		const eps = 1e-6
		if p.r.XMin > want.XMin+eps || p.r.YMin > want.YMin+eps || p.r.XMax < want.XMax-eps || p.r.YMax < want.YMax-eps {
			t.Errorf("glyph %d: points span %v, outline bounds %v", gid, p.r, want)
		}
	}
	if drawn == 0 {
		t.Fatal("no glyph drew anything")
	}
}

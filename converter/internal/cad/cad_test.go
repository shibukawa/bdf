package cad

import (
	"math"
	"slices"
	"testing"
)

func near(a, b Point, tol float64) bool { return math.Abs(a.X-b.X) <= tol && math.Abs(a.Y-b.Y) <= tol }

// bezier evaluates a cubic Bézier segment.
func bezier(p0, p1, p2, p3 Point, t float64) Point {
	u := 1 - t
	a, b, c, d := u*u*u, 3*u*u*t, 3*u*t*t, t*t*t
	return Point{a*p0.X + b*p1.X + c*p2.X + d*p3.X, a*p0.Y + b*p1.Y + c*p2.Y + d*p3.Y}
}

// TestBSplineBezier checks the Bézier segments of a B-spline against de
// Boor's algorithm.
func TestBSplineBezier(t *testing.T) {
	ctrl := []Point{{0, 0}, {1, 3}, {3, -1}, {4, 2}, {6, 0}, {7, 3}}
	for _, c := range []struct {
		degree int
		knots  []float64
	}{
		{3, []float64{0, 0, 0, 0, 1, 2, 3, 3, 3, 3}},
		{3, []float64{0, 0, 0, 0, 0.3, 0.3, 1, 1, 1, 1}}, // a double knot
		{2, []float64{0, 0, 0, 1, 2, 3, 4, 4, 4}},
	} {
		p := (&Path{}).BSpline(c.degree, c.knots, ctrl, nil)
		// collect the segments and the parameter span of each
		var segs [][4]Point
		var cur Point
		p.walk(func(v byte, pts []Point) {
			switch v {
			case moveTo:
				cur = pts[0]
			case cubicTo:
				segs = append(segs, [4]Point{cur, pts[0], pts[1], pts[2]})
				cur = pts[2]
			}
		})
		var spans [][2]float64
		for i := c.degree; i < len(ctrl); i++ {
			if c.knots[i+1] > c.knots[i] {
				spans = append(spans, [2]float64{c.knots[i], c.knots[i+1]})
			}
		}
		if len(segs) != len(spans) {
			t.Fatalf("degree %d: %d segments for %d spans", c.degree, len(segs), len(spans))
		}
		for i, s := range segs {
			for _, f := range []float64{0, 0.25, 0.5, 0.75, 1} {
				u := spans[i][0] + (spans[i][1]-spans[i][0])*f
				want := deBoor(c.degree, c.knots, ctrl, nil, u)
				if got := bezier(s[0], s[1], s[2], s[3], f); !near(got, want, 1e-9) {
					t.Errorf("degree %d segment %d at %v: %v, want %v", c.degree, i, f, got, want)
				}
			}
		}
	}
}

func TestRationalSpline(t *testing.T) {
	// a quarter circle as a rational quadratic
	w := math.Sqrt2 / 2
	p := (&Path{}).BSpline(2, []float64{0, 0, 0, 1, 1, 1}, []Point{{1, 0}, {1, 1}, {0, 1}}, []float64{1, w, 1})
	p.walk(func(v byte, pts []Point) {
		for _, q := range pts {
			if r := q.Len(); math.Abs(r-1) > 1e-6 {
				t.Errorf("point %v is %v from the centre", q, r)
			}
		}
	})
}

func TestArc(t *testing.T) {
	// cubic approximations of arcs stay within 0.03% of the radius
	p := (&Path{}).Arc(Point{1, 2}, 10, 0.3, 4)
	var cur Point
	p.walk(func(v byte, pts []Point) {
		switch v {
		case moveTo:
			cur = pts[0]
		case cubicTo:
			for f := 0.0; f <= 1; f += 0.1 {
				q := bezier(cur, pts[0], pts[1], pts[2], f)
				if r := q.Sub(Point{1, 2}).Len(); math.Abs(r-10) > 0.003 {
					t.Errorf("point at %v off the circle: %v", f, r)
				}
			}
			cur = pts[2]
		}
	})
	if !near(cur, Point{1 + 10*math.Cos(4), 2 + 10*math.Sin(4)}, 1e-9) {
		t.Errorf("arc ends at %v", cur)
	}
}

func TestBulgeArc(t *testing.T) {
	for _, c := range []struct {
		bulge  float64
		centre Point
		r      float64
	}{
		{1, Point{1, 0}, 1},         // a half circle, counter-clockwise: below the chord
		{-1, Point{1, 0}, 1},        // clockwise: above
		{0.5, Point{1, 0.75}, 1.25}, // less than half: the centre on the left
	} {
		ctr, r, a0, a1 := BulgeArc(Point{0, 0}, Point{2, 0}, c.bulge)
		if !near(ctr, c.centre, 1e-9) || math.Abs(r-c.r) > 1e-9 {
			t.Errorf("bulge %v: centre %v radius %v", c.bulge, ctr, r)
		}
		mid := (a0 + a1) / 2
		m := Point{ctr.X + r*math.Cos(mid), ctr.Y + r*math.Sin(mid)}
		if (m.Y < 0) != (c.bulge > 0) {
			t.Errorf("bulge %v: the arc passes %v", c.bulge, m)
		}
	}
}

func TestInterpolate(t *testing.T) {
	fit := []Point{{0, 0}, {1, 2}, {3, 1}, {4, 3}}
	p := (&Path{}).Interpolate(fit, nil, nil)
	var ends []Point
	p.walk(func(v byte, pts []Point) {
		switch v {
		case moveTo:
			ends = append(ends, pts[0])
		case cubicTo:
			ends = append(ends, pts[2])
		}
	})
	if !slices.Equal(ends, fit) {
		t.Errorf("the curve passes %v, want %v", ends, fit)
	}
}

func TestDashPattern(t *testing.T) {
	for _, c := range []struct {
		in     []float64
		dash   []float64
		offset float64
	}{
		{nil, nil, 0},
		{[]float64{1}, nil, 0},
		{[]float64{0.5, -0.25}, []float64{0.5, 0.25}, 0},
		{[]float64{1.25, -0.25, 0, -0.25}, []float64{1.25, 0.25, 0, 0.25}, 0},
		{[]float64{-1, 2, -3}, []float64{2, 4}, 5},
		{[]float64{1, 1, -1}, []float64{2, 1}, 0},
		{[]float64{0, -0.5}, []float64{0, 0.5}, 0},
		{[]float64{1, -2, 1}, []float64{2, 2}, 1}, // a dash at both ends
		{[]float64{-1}, nil, 0},
	} {
		dash, off := DashPattern(c.in)
		if !slices.Equal(dash, c.dash) || off != c.offset {
			t.Errorf("DashPattern(%v) = %v, %v; want %v, %v", c.in, dash, off, c.dash, c.offset)
		}
	}
}

func TestHatch(t *testing.T) {
	d := &Drawing{}
	square := (&Path{}).Polyline([]Point{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, true)
	// horizontal lines every 1: 11 lines cross the square (edges included)
	if !d.Hatch(square, []PatternLine{{Angle: 0, Offset: Point{0, 1}}}, Pen{}) {
		t.Fatal("hatch refused")
	}
	if len(d.Items) != 1 || d.Items[0].kind != kGroup || len(d.Items[0].items) != 1 {
		t.Fatalf("items %+v", d.Items)
	}
	lines := 0
	d.Items[0].items[0].path.walk(func(v byte, _ []Point) {
		if v == moveTo {
			lines++
		}
	})
	if lines != 11 {
		t.Errorf("%d lines", lines)
	}
	// too dense
	if d.Hatch(square, []PatternLine{{Angle: 0, Offset: Point{0, 1e-4}}}, Pen{}) {
		t.Error("a hatch of 100000 lines was drawn")
	}
}

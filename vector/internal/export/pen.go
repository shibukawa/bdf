// Package export holds what the writers of PDF and PostScript files share:
// the pen that writes the segments of Canvas paths as the operators of
// both, and the readings of a document's paints and instructions.
package export

import (
	"bytes"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/affine"
)

// Pen writes the segments of a path as the operators m, l, c and h (of a
// PDF content stream, or of the procedures of psbdf's prolog), with the
// rules of Canvas paths: a segment without a subpath starts one, and a
// closed subpath is followed by one that starts where it started. It
// keeps the box of the points it writes; with no buffer it only keeps the
// box. It is an sfnt.Pen, for the outlines of glyphs in font units under
// its transform.
type Pen struct {
	B *bytes.Buffer
	M affine.Matrix

	x, y   float64 // the current point
	sx, sy float64 // where the subpath started
	has    bool    // there is a current point
	moved  bool    // the subpath has its m
	// Drawn says the path has segments; Box is the box of their points,
	// when Boxed.
	Drawn bool
	Box   [4]float64
	Boxed bool
}

// NewPen makes a pen that writes to b (nil: keeps the box only) through
// the transform m.
func NewPen(b *bytes.Buffer, m affine.Matrix) *Pen { return &Pen{B: b, M: m} }

func (p *Pen) pt(x, y float64) string {
	x, y = p.M.Apply(x, y)
	if !math.IsNaN(x) && !math.IsNaN(y) && !math.IsInf(x, 0) && !math.IsInf(y, 0) {
		if !p.Boxed {
			p.Box, p.Boxed = [4]float64{x, y, x, y}, true
		} else {
			p.Box = [4]float64{math.Min(p.Box[0], x), math.Min(p.Box[1], y), math.Max(p.Box[2], x), math.Max(p.Box[3], y)}
		}
	}
	if p.B == nil {
		return ""
	}
	return Num(x) + " " + Num(y)
}

func (p *Pen) write(s string) {
	if p.B != nil {
		p.B.WriteString(s)
	}
}

// begin writes the m of a subpath whose first segment follows.
func (p *Pen) begin() {
	if !p.moved {
		p.write(p.pt(p.sx, p.sy) + " m\n")
		p.moved = true
	}
}

func (p *Pen) MoveTo(x, y float64) {
	p.x, p.y, p.sx, p.sy, p.has, p.moved = x, y, x, y, true, false
}

func (p *Pen) LineTo(x, y float64) {
	if !p.has {
		p.MoveTo(x, y)
		return
	}
	p.begin()
	p.write(p.pt(x, y) + " l\n")
	p.x, p.y, p.Drawn = x, y, true
}

func (p *Pen) CubeTo(c1x, c1y, c2x, c2y, x, y float64) {
	if !p.has {
		p.MoveTo(c1x, c1y)
	}
	p.begin()
	p.write(p.pt(c1x, c1y) + " " + p.pt(c2x, c2y) + " " + p.pt(x, y) + " c\n")
	p.x, p.y, p.Drawn = x, y, true
}

func (p *Pen) QuadTo(cx, cy, x, y float64) {
	if !p.has {
		p.MoveTo(cx, cy)
	}
	x0, y0 := p.x, p.y
	p.CubeTo(x0+2.0/3*(cx-x0), y0+2.0/3*(cy-y0), x+2.0/3*(cx-x), y+2.0/3*(cy-y), x, y)
}

func (p *Pen) Close() {
	if !p.has {
		return
	}
	if p.moved {
		p.write("h\n")
	}
	// the next subpath starts where this one started
	p.x, p.y, p.moved = p.sx, p.sy, false
}

func (p *Pen) Rect(x, y, w, h float64) {
	p.MoveTo(x, y)
	p.LineTo(x+w, y)
	p.LineTo(x+w, y+h)
	p.LineTo(x, y+h)
	p.Close()
	p.MoveTo(x, y)
}

// Ellipse adds an elliptical arc with Canvas semantics: a line from the
// current point to the start of the arc, angles measured clockwise from
// the x axis (y down) and rotated by rot.
func (p *Pen) Ellipse(cx, cy, rx, ry, rot, a0, a1 float64, ccw bool) {
	if rx < 0 || ry < 0 {
		return
	}
	const tau = 2 * math.Pi
	var sweep float64
	if !ccw {
		if a1-a0 >= tau {
			sweep = tau
		} else {
			sweep = math.Mod(a1-a0, tau)
			if sweep < 0 {
				sweep += tau
			}
		}
	} else {
		if a0-a1 >= tau {
			sweep = -tau
		} else {
			sweep = -math.Mod(a0-a1, tau)
			if sweep > 0 {
				sweep -= tau
			}
		}
	}
	sin, cos := math.Sincos(rot)
	at := func(a float64) (float64, float64) {
		s, c := math.Sincos(a)
		x, y := rx*c, ry*s
		return cx + x*cos - y*sin, cy + x*sin + y*cos
	}
	sx, sy := at(a0)
	if p.has {
		p.LineTo(sx, sy)
	} else {
		p.MoveTo(sx, sy)
	}
	if sweep == 0 || rx == 0 || ry == 0 || math.IsNaN(sweep) {
		ex, ey := at(a0 + sweep)
		p.LineTo(ex, ey)
		return
	}
	n := int(math.Ceil(math.Abs(sweep) / (math.Pi / 2)))
	step := sweep / float64(n)
	k := 4.0 / 3 * math.Tan(step/4)
	a := a0
	for i := 0; i < n; i++ {
		b := a + step
		sa, ca := math.Sincos(a)
		sb, cb := math.Sincos(b)
		// control points of the unit circle arc, then scaled and rotated
		tr := func(x, y float64) (float64, float64) {
			x, y = rx*x, ry*y
			return cx + x*cos - y*sin, cy + x*sin + y*cos
		}
		c1x, c1y := tr(ca-k*sa, sa+k*ca)
		c2x, c2y := tr(cb+k*sb, sb-k*cb)
		ex, ey := tr(cb, sb)
		p.CubeTo(c1x, c1y, c2x, c2y, ex, ey)
		a = b
	}
}

// ArcTo adds a line and an arc of radius r tangent to the lines from the
// current point to (x1, y1) and from there to (x2, y2).
func (p *Pen) ArcTo(x1, y1, x2, y2, r float64) {
	if r < 0 {
		return
	}
	if !p.has {
		p.MoveTo(x1, y1)
	}
	x0, y0 := p.x, p.y
	if (x0 == x1 && y0 == y1) || (x1 == x2 && y1 == y2) || r == 0 {
		p.LineTo(x1, y1)
		return
	}
	ax, ay := x0-x1, y0-y1
	bx, by := x2-x1, y2-y1
	la, lb := math.Hypot(ax, ay), math.Hypot(bx, by)
	cross := ax*by - ay*bx
	if math.Abs(cross) < 1e-12*la*lb {
		p.LineTo(x1, y1) // collinear
		return
	}
	cosT := (ax*bx + ay*by) / (la * lb)
	theta := math.Acos(math.Max(-1, math.Min(1, cosT)))
	d := r / math.Tan(theta/2)
	t0x, t0y := x1+ax/la*d, y1+ay/la*d
	t1x, t1y := x1+bx/lb*d, y1+by/lb*d
	// centre: along the bisector
	bisx, bisy := ax/la+bx/lb, ay/la+by/lb
	bl := math.Hypot(bisx, bisy)
	h := r / math.Sin(theta/2)
	cx, cy := x1+bisx/bl*h, y1+bisy/bl*h
	a0 := math.Atan2(t0y-cy, t0x-cx)
	a1 := math.Atan2(t1y-cy, t1x-cx)
	p.LineTo(t0x, t0y)
	p.Ellipse(cx, cy, r, r, 0, a0, a1, cross > 0)
}

// RoundRect adds a rectangle with corners rounded by r (Canvas roundRect
// with one radius).
func (p *Pen) RoundRect(x, y, w, h, r float64) {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	r = math.Max(0, math.Min(r, math.Min(w, h)/2))
	if r == 0 || math.IsNaN(r) {
		p.Rect(x, y, w, h)
		return
	}
	p.MoveTo(x+r, y)
	p.LineTo(x+w-r, y)
	p.Ellipse(x+w-r, y+r, r, r, 0, -math.Pi/2, 0, false)
	p.LineTo(x+w, y+h-r)
	p.Ellipse(x+w-r, y+h-r, r, r, 0, 0, math.Pi/2, false)
	p.LineTo(x+r, y+h)
	p.Ellipse(x+r, y+h-r, r, r, 0, math.Pi/2, math.Pi, false)
	p.LineTo(x, y+r)
	p.Ellipse(x+r, y+r, r, r, 0, math.Pi, 3*math.Pi/2, false)
	p.Close()
	p.MoveTo(x, y)
}

// Path writes a path of a document, moved by the transform.
func (p *Pen) Path(bp *bdf.Path) {
	a := bp.Args
	i := 0
	f := func(k int) float64 {
		if i+k >= len(a) {
			return 0
		}
		v := float64(a[i+k])
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0
		}
		return v
	}
	for _, v := range bp.Verbs {
		switch v {
		case bdf.VerbMove:
			p.MoveTo(f(0), f(1))
			i += 2
		case bdf.VerbLine:
			p.LineTo(f(0), f(1))
			i += 2
		case bdf.VerbQuad:
			p.QuadTo(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbCubic:
			p.CubeTo(f(0), f(1), f(2), f(3), f(4), f(5))
			i += 6
		case bdf.VerbClose:
			p.Close()
		case bdf.VerbRect:
			p.Rect(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbEllipse:
			p.Ellipse(f(0), f(1), f(2), f(3), f(4), f(5), f(6), f(7) != 0)
			i += 8
		case bdf.VerbArcTo:
			p.ArcTo(f(0), f(1), f(2), f(3), f(4))
			i += 5
		case bdf.VerbRoundRect:
			p.RoundRect(f(0), f(1), f(2), f(3), f(4))
			i += 5
		}
	}
}

package pdfbdf

import (
	"bytes"
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// matrix is an affine transform: x' = a·x + c·y + e, y' = b·x + d·y + f,
// as [a b c d e f].
type matrix [6]float64

var identity = matrix{1, 0, 0, 1, 0, 0}

// mul returns the transform that applies n, then m.
func (m matrix) mul(n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[2]*n[1],
		m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3],
		m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4],
		m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

func (m matrix) apply(x, y float64) (float64, float64) {
	return m[0]*x + m[2]*y + m[4], m[1]*x + m[3]*y + m[5]
}

func (m matrix) translate(x, y float64) matrix { return m.mul(matrix{1, 0, 0, 1, x, y}) }

// areaScale is how much the transform scales lengths, as the root of how
// much it scales areas.
func (m matrix) areaScale() float64 { return math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2])) }

func (m matrix) finite() bool {
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

// invert returns the inverse of the transform, when it has one.
func (m matrix) invert() (matrix, bool) {
	det := m[0]*m[3] - m[1]*m[2]
	if det == 0 || math.IsNaN(det) || math.IsInf(det, 0) || math.IsInf(1/det, 0) {
		return matrix{}, false
	}
	return matrix{
		m[3] / det, -m[1] / det, -m[2] / det, m[0] / det,
		(m[2]*m[5] - m[3]*m[4]) / det, (m[1]*m[4] - m[0]*m[5]) / det,
	}, true
}

// exact writes the matrix with the digits that a transform and its
// inverse, one after the other, need to be no transform.
func (m matrix) exact() string {
	var b strings.Builder
	for i, v := range m {
		if i > 0 {
			b.WriteByte(' ')
		}
		if math.IsNaN(v) || math.IsInf(v, 0) || v == 0 {
			v = 0 // (and no -0)
		}
		b.WriteString(strconv.FormatFloat(v, 'f', -1, 32))
	}
	return b.String()
}

// bounds returns the box around a rectangle under the transform.
func (m matrix) bounds(x0, y0, x1, y1 float64) [4]float64 {
	out := [4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}
	for _, c := range [][2]float64{{x0, y0}, {x1, y0}, {x0, y1}, {x1, y1}} {
		x, y := m.apply(c[0], c[1])
		out[0], out[1], out[2], out[3] = math.Min(out[0], x), math.Min(out[1], y), math.Max(out[2], x), math.Max(out[3], y)
	}
	return out
}

// cm writes the matrix as the operands of a cm, Tm or pattern matrix.
func (m matrix) String() string {
	return fine(m[0]) + " " + fine(m[1]) + " " + fine(m[2]) + " " + fine(m[3]) + " " + num(m[4]) + " " + num(m[5])
}

// pen writes the segments of a path as the operators of a content stream,
// with the rules of Canvas paths: a segment without a subpath starts one,
// and a closed subpath is followed by one that starts where it started.
type pen struct {
	b      *bytes.Buffer
	m      matrix
	x, y   float64 // the current point
	sx, sy float64 // where the subpath started
	has    bool    // there is a current point
	moved  bool    // the subpath has its m operator
	drawn  bool    // the path has segments
}

func (p *pen) pt(x, y float64) string {
	x, y = p.m.apply(x, y)
	return num(x) + " " + num(y)
}

// begin writes the m of a subpath whose first segment follows.
func (p *pen) begin() {
	if !p.moved {
		p.b.WriteString(p.pt(p.sx, p.sy) + " m\n")
		p.moved = true
	}
}

func (p *pen) moveTo(x, y float64) {
	p.x, p.y, p.sx, p.sy, p.has, p.moved = x, y, x, y, true, false
}

func (p *pen) lineTo(x, y float64) {
	if !p.has {
		p.moveTo(x, y)
		return
	}
	p.begin()
	p.b.WriteString(p.pt(x, y) + " l\n")
	p.x, p.y, p.drawn = x, y, true
}

func (p *pen) cubicTo(c1x, c1y, c2x, c2y, x, y float64) {
	if !p.has {
		p.moveTo(c1x, c1y)
	}
	p.begin()
	p.b.WriteString(p.pt(c1x, c1y) + " " + p.pt(c2x, c2y) + " " + p.pt(x, y) + " c\n")
	p.x, p.y, p.drawn = x, y, true
}

func (p *pen) quadTo(cx, cy, x, y float64) {
	if !p.has {
		p.moveTo(cx, cy)
	}
	x0, y0 := p.x, p.y
	p.cubicTo(x0+2.0/3*(cx-x0), y0+2.0/3*(cy-y0), x+2.0/3*(cx-x), y+2.0/3*(cy-y), x, y)
}

func (p *pen) close() {
	if !p.has {
		return
	}
	if p.moved {
		p.b.WriteString("h\n")
	}
	// the next subpath starts where this one started
	p.x, p.y, p.moved = p.sx, p.sy, false
}

func (p *pen) rect(x, y, w, h float64) {
	p.moveTo(x, y)
	p.lineTo(x+w, y)
	p.lineTo(x+w, y+h)
	p.lineTo(x, y+h)
	p.close()
	p.moveTo(x, y)
}

// ellipse adds an elliptical arc with Canvas semantics: a line from the
// current point to the start of the arc, angles measured clockwise from
// the x axis (y down) and rotated by rot.
func (p *pen) ellipse(cx, cy, rx, ry, rot, a0, a1 float64, ccw bool) {
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
		p.lineTo(sx, sy)
	} else {
		p.moveTo(sx, sy)
	}
	if sweep == 0 || rx == 0 || ry == 0 || math.IsNaN(sweep) {
		ex, ey := at(a0 + sweep)
		p.lineTo(ex, ey)
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
		p.cubicTo(c1x, c1y, c2x, c2y, ex, ey)
		a = b
	}
}

// arcTo adds a line and an arc of radius r tangent to the lines from the
// current point to (x1, y1) and from there to (x2, y2).
func (p *pen) arcTo(x1, y1, x2, y2, r float64) {
	if r < 0 {
		return
	}
	if !p.has {
		p.moveTo(x1, y1)
	}
	x0, y0 := p.x, p.y
	if (x0 == x1 && y0 == y1) || (x1 == x2 && y1 == y2) || r == 0 {
		p.lineTo(x1, y1)
		return
	}
	ax, ay := x0-x1, y0-y1
	bx, by := x2-x1, y2-y1
	la, lb := math.Hypot(ax, ay), math.Hypot(bx, by)
	cross := ax*by - ay*bx
	if math.Abs(cross) < 1e-12*la*lb {
		p.lineTo(x1, y1) // collinear
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
	p.lineTo(t0x, t0y)
	p.ellipse(cx, cy, r, r, 0, a0, a1, cross > 0)
}

// roundRect adds a rectangle with corners rounded by r (Canvas roundRect
// with one radius).
func (p *pen) roundRect(x, y, w, h, r float64) {
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	r = math.Max(0, math.Min(r, math.Min(w, h)/2))
	if r == 0 || math.IsNaN(r) {
		p.rect(x, y, w, h)
		return
	}
	p.moveTo(x+r, y)
	p.lineTo(x+w-r, y)
	p.ellipse(x+w-r, y+r, r, r, 0, -math.Pi/2, 0, false)
	p.lineTo(x+w, y+h-r)
	p.ellipse(x+w-r, y+h-r, r, r, 0, 0, math.Pi/2, false)
	p.lineTo(x+r, y+h)
	p.ellipse(x+r, y+h-r, r, r, 0, math.Pi/2, math.Pi, false)
	p.lineTo(x, y+r)
	p.ellipse(x+r, y+r, r, r, 0, math.Pi, 3*math.Pi/2, false)
	p.close()
	p.moveTo(x, y)
}

// path writes a path of a document, moved by m, and reports whether it
// has segments.
func (p *pen) path(bp *bdf.Path) {
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
			p.moveTo(f(0), f(1))
			i += 2
		case bdf.VerbLine:
			p.lineTo(f(0), f(1))
			i += 2
		case bdf.VerbQuad:
			p.quadTo(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbCubic:
			p.cubicTo(f(0), f(1), f(2), f(3), f(4), f(5))
			i += 6
		case bdf.VerbClose:
			p.close()
		case bdf.VerbRect:
			p.rect(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbEllipse:
			p.ellipse(f(0), f(1), f(2), f(3), f(4), f(5), f(6), f(7) != 0)
			i += 8
		case bdf.VerbArcTo:
			p.arcTo(f(0), f(1), f(2), f(3), f(4))
			i += 5
		case bdf.VerbRoundRect:
			p.roundRect(f(0), f(1), f(2), f(3), f(4))
			i += 5
		}
	}
}

// outlinePen is a pen for the outlines of glyphs (sfnt.Pen), in font
// units under a transform.
type outlinePen struct{ p *pen }

func (o outlinePen) MoveTo(x, y float64)                 { o.p.moveTo(x, y) }
func (o outlinePen) LineTo(x, y float64)                 { o.p.lineTo(x, y) }
func (o outlinePen) QuadTo(cx, cy, x, y float64)         { o.p.quadTo(cx, cy, x, y) }
func (o outlinePen) CubeTo(ax, ay, bx, by, x, y float64) { o.p.cubicTo(ax, ay, bx, by, x, y) }
func (o outlinePen) Close()                              { o.p.close() }

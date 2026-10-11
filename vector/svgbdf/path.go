package svgbdf

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// matrix is an affine transform as Canvas 2D keeps it: x' = a x + c y + e,
// y' = b x + d y + f.
type matrix [6]float64

var identity = matrix{1, 0, 0, 1, 0, 0}

// mul returns m × n: n applies first (ctx.transform(n) on a context at m).
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

func (m matrix) translate(x, y float64) matrix { return m.mul(matrix{1, 0, 0, 1, x, y}) }
func (m matrix) scale(x, y float64) matrix     { return m.mul(matrix{x, 0, 0, y, 0, 0}) }

// maxScale is the largest factor the transform stretches a length by.
func (m matrix) maxScale() float64 {
	return math.Max(math.Hypot(m[0], m[1]), math.Hypot(m[2], m[3]))
}

// areaScale is the square root of the area factor (the zoom of a uniform
// transform).
func (m matrix) areaScale() float64 { return math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2])) }

func (m matrix) finite() bool {
	for _, v := range m {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return false
		}
	}
	return true
}

// pageDigits is how many decimals a length in the units of the page is
// written with: a hundredth of a point shows at no zoom a screen has.
const pageDigits = 3

// digits is how many decimals the coordinates of what is drawn under m are
// written with, to be as exact on the page as pageDigits: the outline of a
// glyph in the thousandths of an em needs none, a unit square stretched
// over the page all of them.
func (m matrix) digits() int {
	s := m.maxScale()
	if !(s > 0) || math.IsInf(s, 0) {
		return pageDigits
	}
	return min(max(int(math.Ceil(math.Log10(s)))+pageDigits, 0), 9)
}

// appendNum writes a number with at most d decimals, as short as SVG reads
// it: no trailing zeros, no zero before the point.
func appendNum(b []byte, v float64, d int) []byte {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return append(b, '0')
	}
	from := len(b)
	b = strconv.AppendFloat(b, v, 'f', d, 64)
	if d > 0 {
		n := len(b)
		for b[n-1] == '0' {
			n--
		}
		if b[n-1] == '.' {
			n--
		}
		b = b[:n]
	}
	s := b[from:]
	switch {
	case string(s) == "-0":
		return append(b[:from], '0')
	case len(s) > 2 && s[0] == '0' && s[1] == '.':
		return append(b[:from], s[1:]...)
	case len(s) > 3 && s[0] == '-' && s[1] == '0' && s[2] == '.':
		copy(s[1:], s[2:])
		return b[:len(b)-1]
	}
	return b
}

func num(v float64, d int) string { return string(appendNum(nil, v, d)) }

// factor writes a number that scales others (a matrix entry, an opacity)
// with six significant digits.
func factor(v float64) string {
	if v == 0 || math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	d := 5 - int(math.Floor(math.Log10(math.Abs(v))))
	return num(v, min(max(d, 0), 12))
}

// transform is the transform attribute of what is drawn under m, with a
// space before it, or nothing for the identity.
func (m matrix) transform() string {
	if m == identity {
		return ""
	}
	if m[0] == 1 && m[1] == 0 && m[2] == 0 && m[3] == 1 {
		return ` transform="translate(` + num(m[4], pageDigits) + " " + num(m[5], pageDigits) + `)"`
	}
	return ` transform="matrix(` + factor(m[0]) + " " + factor(m[1]) + " " + factor(m[2]) + " " + factor(m[3]) + " " +
		num(m[4], pageDigits) + " " + num(m[5], pageDigits) + `)"`
}

// pathWriter writes the path data of a Path2D: the commands of Canvas,
// which start a sub-path where SVG would have an error.
type pathWriter struct {
	b      []byte
	d      int
	x, y   float64 // the current point
	sx, sy float64 // where the sub-path started
	open   bool    // there is a current point
	bad    bool    // a coordinate is no number
}

func (w *pathWriter) cmd(c byte, v ...float64) {
	w.b = append(w.b, c)
	for i, f := range v {
		if math.IsNaN(f) || math.IsInf(f, 0) {
			w.bad = true
		}
		n := len(w.b)
		w.b = appendNum(w.b, f, w.d)
		// a separator only where two numbers would read as one
		if i > 0 && w.b[n] != '-' {
			w.b = append(w.b, 0)
			copy(w.b[n+1:], w.b[n:])
			w.b[n] = ' '
		}
	}
}

func (w *pathWriter) moveTo(x, y float64) {
	w.cmd('M', x, y)
	w.x, w.y, w.sx, w.sy, w.open = x, y, x, y, true
}

func (w *pathWriter) lineTo(x, y float64) {
	if !w.open {
		w.moveTo(x, y)
		return
	}
	w.cmd('L', x, y)
	w.x, w.y = x, y
}

func (w *pathWriter) close() {
	if w.open {
		w.b = append(w.b, 'Z')
		w.x, w.y = w.sx, w.sy
	}
}

func (w *pathWriter) rect(x, y, width, height float64) {
	w.moveTo(x, y)
	w.cmd('h', width)
	w.cmd('v', height)
	w.cmd('h', -width)
	w.b = append(w.b, 'Z')
	if math.IsNaN(width+height) || math.IsInf(width+height, 0) {
		w.bad = true
	}
}

// ellipse adds an elliptical arc with Canvas semantics: a line from the
// current point to the start of the arc, angles measured clockwise from
// the x axis (y down) and rotated by rot.
func (w *pathWriter) ellipse(cx, cy, rx, ry, rot, a0, a1 float64, ccw bool) {
	if rx < 0 || ry < 0 {
		return
	}
	const tau = 2 * math.Pi
	var sweep float64
	if !ccw {
		if a1-a0 >= tau {
			sweep = tau
		} else if sweep = math.Mod(a1-a0, tau); sweep < 0 {
			sweep += tau
		}
	} else {
		if a0-a1 >= tau {
			sweep = -tau
		} else if sweep = -math.Mod(a0-a1, tau); sweep > 0 {
			sweep -= tau
		}
	}
	sin, cos := math.Sincos(rot)
	at := func(a float64) (float64, float64) {
		s, c := math.Sincos(a)
		x, y := rx*c, ry*s
		return cx + x*cos - y*sin, cy + x*sin + y*cos
	}
	x, y := at(a0)
	w.lineTo(x, y)
	if sweep == 0 || rx == 0 || ry == 0 || math.IsNaN(sweep) {
		x, y = at(a0 + sweep)
		w.lineTo(x, y)
		return
	}
	// arcs of a third of a turn at most: one of a half or more has two
	// readings
	n := int(math.Ceil(math.Abs(sweep) / (tau / 3)))
	flag := 0.0
	if sweep > 0 {
		flag = 1
	}
	deg := rot * 180 / math.Pi
	for i := 1; i <= n; i++ {
		x, y = at(a0 + sweep*float64(i)/float64(n))
		w.cmd('A', rx, ry)
		w.b = append(w.b, ' ')
		w.b = appendNum(w.b, deg, 3)
		w.b = append(w.b, " 0 "...)
		w.b = appendNum(w.b, flag, 0)
		w.b = append(w.b, ' ')
		w.b = appendNum(w.b, x, w.d)
		w.b = append(w.b, ' ')
		w.b = appendNum(w.b, y, w.d)
		if math.IsNaN(x+y+deg) || math.IsInf(x+y+deg, 0) {
			w.bad = true
		}
	}
	w.x, w.y = x, y
}

// arcTo adds a line and an arc of radius r tangent to the lines from the
// current point to (x1, y1) and from there to (x2, y2).
func (w *pathWriter) arcTo(x1, y1, x2, y2, r float64) {
	if r < 0 {
		return
	}
	if !w.open {
		w.moveTo(x1, y1)
	}
	x0, y0 := w.x, w.y
	if (x0 == x1 && y0 == y1) || (x1 == x2 && y1 == y2) || r == 0 {
		w.lineTo(x1, y1)
		return
	}
	ax, ay := x0-x1, y0-y1
	bx, by := x2-x1, y2-y1
	la, lb := math.Hypot(ax, ay), math.Hypot(bx, by)
	cross := ax*by - ay*bx
	if math.Abs(cross) < 1e-12*la*lb {
		w.lineTo(x1, y1) // collinear
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
	w.lineTo(t0x, t0y)
	w.ellipse(cx, cy, r, r, 0, math.Atan2(t0y-cy, t0x-cx), math.Atan2(t1y-cy, t1x-cx), cross > 0)
}

// roundRect adds a rectangle with corners rounded by r (Canvas roundRect
// with one radius).
func (w *pathWriter) roundRect(x, y, width, height, r float64) {
	if width < 0 {
		x, width = x+width, -width
	}
	if height < 0 {
		y, height = y+height, -height
	}
	r = math.Max(0, math.Min(r, math.Min(width, height)/2))
	if r == 0 || math.IsNaN(r) {
		w.rect(x, y, width, height)
		return
	}
	w.moveTo(x+r, y)
	w.lineTo(x+width-r, y)
	w.ellipse(x+width-r, y+r, r, r, 0, -math.Pi/2, 0, false)
	w.lineTo(x+width, y+height-r)
	w.ellipse(x+width-r, y+height-r, r, r, 0, 0, math.Pi/2, false)
	w.lineTo(x+r, y+height)
	w.ellipse(x+r, y+height-r, r, r, 0, math.Pi/2, math.Pi, false)
	w.lineTo(x, y+r)
	w.ellipse(x+r, y+r, r, r, 0, math.Pi, 3*math.Pi/2, false)
	w.b = append(w.b, 'Z')
	w.x, w.y, w.sx, w.sy = x, y, x, y
}

// pathData returns the d attribute of a path with d decimals, and false
// for a path with nothing to draw or with a coordinate that is no number.
func pathData(bp *bdf.Path, d int) (string, bool) {
	w := pathWriter{b: make([]byte, 0, 8*len(bp.Args)), d: d}
	a := bp.Args
	i := 0
	f := func(k int) float64 {
		if i+k >= len(a) {
			return 0
		}
		return float64(a[i+k])
	}
	for _, v := range bp.Verbs {
		switch v {
		case bdf.VerbMove:
			w.moveTo(f(0), f(1))
			i += 2
		case bdf.VerbLine:
			w.lineTo(f(0), f(1))
			i += 2
		case bdf.VerbQuad:
			if !w.open {
				w.moveTo(f(0), f(1))
			}
			w.cmd('Q', f(0), f(1), f(2), f(3))
			w.x, w.y = f(2), f(3)
			i += 4
		case bdf.VerbCubic:
			if !w.open {
				w.moveTo(f(0), f(1))
			}
			w.cmd('C', f(0), f(1), f(2), f(3), f(4), f(5))
			w.x, w.y = f(4), f(5)
			i += 6
		case bdf.VerbClose:
			w.close()
		case bdf.VerbRect:
			w.rect(f(0), f(1), f(2), f(3))
			i += 4
		case bdf.VerbEllipse:
			w.ellipse(f(0), f(1), f(2), f(3), f(4), f(5), f(6), f(7) != 0)
			i += 8
		case bdf.VerbArcTo:
			w.arcTo(f(0), f(1), f(2), f(3), f(4))
			i += 5
		case bdf.VerbRoundRect:
			w.roundRect(f(0), f(1), f(2), f(3), f(4))
			i += 5
		}
	}
	if w.bad || len(w.b) == 0 {
		return "", false
	}
	return string(w.b), true
}

// rectData returns the path data of a rectangle.
func rectData(x, y, width, height float64, d int) (string, bool) {
	w := pathWriter{d: d}
	w.rect(x, y, width, height)
	return string(w.b), !w.bad
}

// escaper writes text as the content of an element or the value of an
// attribute in double quotes.
var escaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")

// escape returns a string as XML holds it: what XML has no character for
// (control characters, and bytes that are no UTF-8) becomes U+FFFD, and the
// white space Canvas draws as a space is one.
func escape(s string) string {
	clean := true
	for _, r := range s {
		if r < 0x20 || r == 0xfffd || r == 0xfffe || r == 0xffff || r == '&' || r == '<' || r == '>' || r == '"' {
			clean = false
			break
		}
	}
	if clean {
		return s
	}
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '\t' || r == '\n' || r == '\r' || r == '\f' || r == '\v':
			b.WriteByte(' ')
		case r < 0x20 || r == 0xfffe || r == 0xffff:
			b.WriteRune(0xfffd)
		default:
			b.WriteRune(r)
		}
	}
	return escaper.Replace(b.String())
}

package gerber

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// macro is an aperture macro (AM): its statements, primitives and
// variable definitions, kept as text and evaluated for each aperture that
// instantiates it.
type macro struct {
	name  string
	stmts []string
}

// maxOutline bounds the vertices of an outline primitive.
const maxOutline = 1 << 16

// shape evaluates the macro with the parameters of an AD command and
// returns the shape in the file's units. warn reports what is not drawn.
func (m *macro) shape(params []float64, warn func(string)) (*shape, error) {
	vars := map[int]float64{}
	for i, v := range params {
		vars[i+1] = v
	}
	var prims []prim
	for _, st := range m.stmts {
		st = strings.TrimLeft(st, " \t\r\n")
		if st == "" || st[0] == '0' && (len(st) == 1 || st[1] < '0' || st[1] > '9') {
			continue // a comment: code 0 and its text
		}
		st = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\t' || r == '\r' || r == '\n' {
				return -1
			}
			return r
		}, st)
		if st == "" {
			continue
		}
		if st[0] == '$' {
			// $n=expression
			eq := strings.IndexByte(st, '=')
			if eq < 0 {
				return nil, fmt.Errorf("macro %s: bad statement %q", m.name, st)
			}
			n, err := strconv.Atoi(st[1:eq])
			if err != nil || n < 0 || n > 1<<16 {
				return nil, fmt.Errorf("macro %s: bad variable %q", m.name, st[:eq])
			}
			v, err := eval(st[eq+1:], vars)
			if err != nil {
				return nil, fmt.Errorf("macro %s: %w", m.name, err)
			}
			vars[n] = v
			continue
		}
		fields := strings.Split(st, ",")
		code, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("macro %s: bad primitive %q", m.name, st)
		}
		if code == 0 {
			continue
		}
		args := make([]float64, 0, len(fields)-1)
		for _, f := range fields[1:] {
			v, err := eval(f, vars)
			if err != nil {
				return nil, fmt.Errorf("macro %s: %w", m.name, err)
			}
			args = append(args, v)
		}
		p, ok, err := primitive(code, args, warn)
		if err != nil {
			return nil, fmt.Errorf("macro %s: %w", m.name, err)
		}
		if ok {
			prims = append(prims, p)
		}
	}
	return newShape(prims), nil
}

// arg returns args[i], or 0 when the primitive leaves it out.
func arg(args []float64, i int) float64 {
	if i < len(args) {
		return args[i]
	}
	return 0
}

// primitive returns the part a macro primitive draws; ok is false when it
// draws nothing.
func primitive(code int, a []float64, warn func(string)) (p prim, ok bool, err error) {
	for _, v := range a {
		if !finite(v) || math.Abs(v) > 1e7 {
			return prim{}, false, fmt.Errorf("primitive %d has a parameter out of range", code)
		}
	}
	exposure := func() bool {
		switch e := arg(a, 0); e {
		case 0:
			return false
		case 2:
			warn("macro primitives with exposure 2 (toggle) are drawn as exposure on")
		}
		return true
	}
	var cs []contour
	rot := 0.0
	on := true
	switch code {
	case 1: // circle: exposure, diameter, centre x, y [, rotation]
		on = exposure()
		d := math.Abs(arg(a, 1))
		if d == 0 {
			return prim{}, false, nil
		}
		cs = []contour{circle(vec{arg(a, 2), arg(a, 3)}, d/2)}
		rot = arg(a, 4)
	case 2, 20: // vector line: exposure, width, start x, y, end x, y, rotation
		on = exposure()
		w := math.Abs(arg(a, 1))
		s, e := vec{arg(a, 2), arg(a, 3)}, vec{arg(a, 4), arg(a, 5)}
		d := e.sub(s)
		l := d.len()
		if w == 0 || l == 0 {
			return prim{}, false, nil
		}
		n := vec{-d.Y, d.X}.mul(w / 2 / l)
		cs = []contour{polygon(s.sub(n), e.sub(n), e.add(n), s.add(n))}
		rot = arg(a, 6)
	case 21: // center line: exposure, width, height, centre x, y, rotation
		on = exposure()
		w, h := math.Abs(arg(a, 1)), math.Abs(arg(a, 2))
		if w == 0 || h == 0 {
			return prim{}, false, nil
		}
		cs = []contour{rect(vec{arg(a, 3), arg(a, 4)}, w, h)}
		rot = arg(a, 5)
	case 22: // lower left line: exposure, width, height, lower left x, y, rotation
		on = exposure()
		w, h := math.Abs(arg(a, 1)), math.Abs(arg(a, 2))
		if w == 0 || h == 0 {
			return prim{}, false, nil
		}
		cs = []contour{rect(vec{arg(a, 3) + w/2, arg(a, 4) + h/2}, w, h)}
		rot = arg(a, 5)
	case 4: // outline: exposure, n, n+1 points, rotation
		on = exposure()
		n := int(arg(a, 1))
		if n < 1 || n > maxOutline {
			return prim{}, false, fmt.Errorf("outline primitive with %d vertices", n)
		}
		pts := make([]vec, 0, n+1)
		for i := 0; i <= n && 3+2*i < len(a); i++ {
			pts = append(pts, vec{a[2+2*i], a[3+2*i]})
		}
		if len(pts) < 3 {
			return prim{}, false, nil
		}
		c := polygon(pts...)
		cs = []contour{c.orient(true)}
		rot = arg(a, 2+2*(n+1))
	case 5: // polygon: exposure, vertices, centre x, y, diameter, rotation
		on = exposure()
		n := int(arg(a, 1))
		d := math.Abs(arg(a, 4))
		if n < 3 || n > 12 {
			return prim{}, false, fmt.Errorf("polygon primitive with %d vertices", n)
		}
		if d == 0 {
			return prim{}, false, nil
		}
		cs = []contour{regular(vec{arg(a, 2), arg(a, 3)}, d, n, 0)}
		rot = arg(a, 5)
	case 6: // moiré: centre x, y, outer diameter, ring thickness, gap, rings, crosshair thickness, length, rotation
		c := vec{arg(a, 0), arg(a, 1)}
		od, t, gap := math.Abs(arg(a, 2)), math.Abs(arg(a, 3)), math.Abs(arg(a, 4))
		rings := min(int(arg(a, 5)), 1000)
		ct, cl := math.Abs(arg(a, 6)), math.Abs(arg(a, 7))
		for i := 0; i < rings; i++ {
			d := od - 2*float64(i)*(t+gap)
			if d <= 0 || t <= 0 {
				break
			}
			cs = append(cs, circle(c, d/2))
			if in := d - 2*t; in > 0 {
				h := circle(c, in/2)
				cs = append(cs, h.reverse())
			}
		}
		if ct > 0 && cl > 0 {
			cs = append(cs, rect(c, cl, ct), rect(c, ct, cl))
		}
		rot = arg(a, 8)
	case 7: // thermal: centre x, y, outer diameter, inner diameter, gap, rotation
		cs = thermal(vec{arg(a, 0), arg(a, 1)}, math.Abs(arg(a, 2)), math.Abs(arg(a, 3)), math.Abs(arg(a, 4)))
		rot = arg(a, 5)
	default:
		warn(fmt.Sprintf("aperture macro primitive %d is not drawn", code))
		return prim{}, false, nil
	}
	if len(cs) == 0 {
		return prim{}, false, nil
	}
	if rot != 0 {
		// primitives turn about the origin of the macro
		m := rotate(rot)
		for i := range cs {
			cs[i] = cs[i].transform(m)
		}
	}
	return prim{on: on, contours: cs}, true, nil
}

// thermal returns the four pieces of a thermal primitive: a ring cut by a
// cross of gap thickness.
func thermal(c vec, od, id, gap float64) []contour {
	R, r, g := od/2, id/2, gap/2
	if R <= r || g*math.Sqrt2 >= R {
		return nil
	}
	// the piece of the first quadrant
	ox := math.Sqrt(R*R - g*g)
	var q contour
	if r > g*math.Sqrt2 {
		ix := math.Sqrt(r*r - g*g)
		q = contour{start: vec{ix, g}}
		q.lineTo(vec{ox, g})
		q.arcTo(vec{g, ox}, vec{}, math.Atan2(ox, g)-math.Atan2(g, ox))
		q.lineTo(vec{g, ix})
		q.arcTo(vec{ix, g}, vec{}, math.Atan2(g, ix)-math.Atan2(ix, g))
	} else {
		q = contour{start: vec{g, g}}
		q.lineTo(vec{ox, g})
		q.arcTo(vec{g, ox}, vec{}, math.Atan2(ox, g)-math.Atan2(g, ox))
		q.lineTo(vec{g, g})
	}
	out := make([]contour, 4)
	for i := range out {
		out[i] = q.transform(translate(c).mul(rotate(90 * float64(i))))
	}
	return out
}

// eval evaluates an arithmetic expression of an aperture macro: numbers,
// variables $n, + - x / and parentheses.
func eval(s string, vars map[int]float64) (float64, error) {
	e := &expr{s: s, vars: vars}
	v, err := e.sum()
	if err != nil {
		return 0, err
	}
	if e.i != len(e.s) {
		return 0, fmt.Errorf("bad expression %q", s)
	}
	return v, nil
}

type expr struct {
	s     string
	i     int
	vars  map[int]float64
	depth int
}

func (e *expr) peek() byte {
	if e.i < len(e.s) {
		return e.s[e.i]
	}
	return 0
}

func (e *expr) sum() (float64, error) {
	v, err := e.product()
	if err != nil {
		return 0, err
	}
	for {
		switch e.peek() {
		case '+':
			e.i++
			w, err := e.product()
			if err != nil {
				return 0, err
			}
			v += w
		case '-':
			e.i++
			w, err := e.product()
			if err != nil {
				return 0, err
			}
			v -= w
		default:
			return v, nil
		}
	}
}

func (e *expr) product() (float64, error) {
	v, err := e.unary()
	if err != nil {
		return 0, err
	}
	for {
		switch e.peek() {
		case 'x', 'X':
			e.i++
			w, err := e.unary()
			if err != nil {
				return 0, err
			}
			v *= w
		case '/':
			e.i++
			w, err := e.unary()
			if err != nil {
				return 0, err
			}
			if w == 0 {
				v = 0
			} else {
				v /= w
			}
		default:
			return v, nil
		}
	}
}

func (e *expr) unary() (float64, error) {
	switch e.peek() {
	case '-':
		e.i++
		v, err := e.unary()
		return -v, err
	case '+':
		e.i++
		return e.unary()
	}
	return e.atom()
}

func (e *expr) atom() (float64, error) {
	switch c := e.peek(); {
	case c == '(':
		e.depth++
		if e.depth > 64 {
			return 0, fmt.Errorf("expression nested too deep")
		}
		e.i++
		v, err := e.sum()
		if err != nil {
			return 0, err
		}
		if e.peek() != ')' {
			return 0, fmt.Errorf("bad expression %q", e.s)
		}
		e.i++
		e.depth--
		return v, nil
	case c == '$':
		e.i++
		j := e.i
		for e.i < len(e.s) && e.s[e.i] >= '0' && e.s[e.i] <= '9' {
			e.i++
		}
		n, err := strconv.Atoi(e.s[j:e.i])
		if err != nil {
			return 0, fmt.Errorf("bad variable in %q", e.s)
		}
		return e.vars[n], nil
	case c >= '0' && c <= '9' || c == '.':
		j := e.i
		for e.i < len(e.s) && (e.s[e.i] >= '0' && e.s[e.i] <= '9' || e.s[e.i] == '.') {
			e.i++
		}
		v, err := strconv.ParseFloat(e.s[j:e.i], 64)
		if err != nil {
			return 0, fmt.Errorf("bad number in %q", e.s)
		}
		return v, nil
	}
	return 0, fmt.Errorf("bad expression %q", e.s)
}

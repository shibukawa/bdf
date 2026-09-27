package sxf

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// pval is a parameter of a STEP entity instance.
type pval struct {
	k    byte // 'n' number, 's' string, 'r' reference, 'e' enumeration, 'l' list, 't' typed, '$' null, '*' omitted
	num  float64
	str  string // a string, an enumeration or the name of a typed parameter
	ref  int
	list []pval // the items of a list, the value of a typed parameter
}

type record struct {
	name   string
	params []pval
}

// instance is a simple (one record) or complex (several) entity instance.
type instance struct {
	recs []record
}

func (in *instance) get(name string) *record {
	if in == nil {
		return nil
	}
	for i := range in.recs {
		if in.recs[i].name == name {
			return &in.recs[i]
		}
	}
	return nil
}

// own returns the record with its last n attributes: a record of a simple
// instance holds the attributes of its supertypes first (the name of a
// representation item, say), while a partial record of a complex instance
// holds only its own.
func (r *record) own(n int) *record {
	if r == nil || len(r.params) <= n {
		return r
	}
	return &record{name: r.name, params: r.params[len(r.params)-n:]}
}

func (r *record) p(i int) pval {
	if r == nil || i >= len(r.params) {
		return pval{k: '$'}
	}
	return r.params[i]
}

// p21file is the DATA section of a STEP file.
type p21file struct {
	inst map[int]*instance
	// scadec is set for files written by the SCADEC library
	scadec bool
	// the drawing
	invisible map[int]bool
	r         *renderer
	stack     []int // the symbol representation maps being drawn
	// curves keeps the curves made, by instance (nil while one is made)
	curves map[int]*cad.Path
	// ccw draws trimmed circles counter-clockwise whatever their sense:
	// the arcs of angular dimensions, which SXF fixes so but some writers
	// flag clockwise
	ccw bool
}

var errNotP21 = errors.New("not a STEP file")

// maxCurvePoints bounds the points of a composite curve, and maxList the
// outlines of a fill area and the styles of a fill.
const (
	maxCurvePoints = 1_000_000
	maxList        = 1000
)

// p21Reader reads the clear text encoding of STEP (ISO 10303-21).
type p21Reader struct {
	s   string
	i   int
	err error
}

func (r *p21Reader) skip() {
	for r.i < len(r.s) {
		c := r.s[r.i]
		switch {
		case c == ' ' || c == '\t' || c == '\r' || c == '\n':
			r.i++
		case c == '/' && r.i+1 < len(r.s) && r.s[r.i+1] == '*':
			j := strings.Index(r.s[r.i+2:], "*/")
			if j < 0 {
				r.i = len(r.s)
				return
			}
			r.i += j + 4
		default:
			return
		}
	}
}

func (r *p21Reader) keyword() string {
	j := r.i
	for j < len(r.s) {
		c := r.s[j]
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '!' {
			j++
			continue
		}
		break
	}
	k := strings.ToUpper(r.s[r.i:j])
	r.i = j
	return k
}

// params reads a parenthesized list of parameters.
func (r *p21Reader) params(depth int) []pval {
	r.skip()
	if r.i >= len(r.s) || r.s[r.i] != '(' {
		r.err = fmt.Errorf("a parameter list is missing at byte %d", r.i)
		return nil
	}
	if depth > 64 {
		r.err = errors.New("parameters nested too deeply")
		return nil
	}
	r.i++
	var out []pval
	for {
		r.skip()
		if r.i >= len(r.s) {
			r.err = errors.New("a parameter list is not closed")
			return out
		}
		c := r.s[r.i]
		if c == ')' {
			r.i++
			return out
		}
		if c == ',' {
			r.i++
			continue
		}
		out = append(out, r.param(depth))
		if r.err != nil {
			return out
		}
	}
}

func (r *p21Reader) param(depth int) pval {
	c := r.s[r.i]
	switch {
	case c == '$':
		r.i++
		return pval{k: '$'}
	case c == '*':
		r.i++
		return pval{k: '*'}
	case c == '#':
		r.i++
		j := r.i
		for j < len(r.s) && r.s[j] >= '0' && r.s[j] <= '9' {
			j++
		}
		n, _ := strconv.Atoi(r.s[r.i:j])
		r.i = j
		return pval{k: 'r', ref: n}
	case c == '\'':
		r.i++
		var b strings.Builder
		for r.i < len(r.s) {
			if r.s[r.i] == '\'' {
				if r.i+1 < len(r.s) && r.s[r.i+1] == '\'' {
					b.WriteByte('\'')
					r.i += 2
					continue
				}
				r.i++
				break
			}
			b.WriteByte(r.s[r.i])
			r.i++
		}
		return pval{k: 's', str: decodeStepString(b.String())}
	case c == '.':
		j := strings.IndexByte(r.s[r.i+1:], '.')
		if j < 0 {
			r.err = errors.New("an enumeration is not closed")
			return pval{}
		}
		v := pval{k: 'e', str: strings.ToUpper(r.s[r.i+1 : r.i+1+j])}
		r.i += j + 2
		return v
	case c == '"':
		j := strings.IndexByte(r.s[r.i+1:], '"')
		if j < 0 {
			r.err = errors.New("a binary value is not closed")
			return pval{}
		}
		r.i += j + 2
		return pval{k: '$'}
	case c == '(':
		return pval{k: 'l', list: r.params(depth + 1)}
	case c == '+' || c == '-' || c >= '0' && c <= '9':
		j := r.i + 1
		for j < len(r.s) && (r.s[j] >= '0' && r.s[j] <= '9' || r.s[j] == '.' || r.s[j] == 'E' || r.s[j] == 'e' || r.s[j] == '+' || r.s[j] == '-') {
			j++
		}
		v, err := strconv.ParseFloat(strings.TrimSuffix(r.s[r.i:j], "."), 64)
		if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
			v = 0
		}
		r.i = j
		return pval{k: 'n', num: v}
	case c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z':
		name := r.keyword()
		inner := r.params(depth + 1)
		return pval{k: 't', str: name, list: inner}
	}
	r.err = fmt.Errorf("unexpected %q at byte %d", c, r.i)
	return pval{}
}

// decodeStepString replaces the control directives of STEP strings:
// \X2\...\X0\ (UTF-16), \X4\...\X0\ (UCS-4), \X\hh (ISO 8859-1), \S\c and
// \P?\ (code page switches), \\.
func decodeStepString(s string) string {
	if !strings.Contains(s, "\\") {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		switch {
		case s[i+1] == '\\':
			b.WriteByte('\\')
			i++
		case strings.HasPrefix(s[i:], "\\X2\\") || strings.HasPrefix(s[i:], "\\X4\\"):
			wide := s[i+2] == '4'
			end := strings.Index(s[i+4:], "\\X0\\")
			if end < 0 {
				end = len(s) - i - 4
			}
			hex := s[i+4 : i+4+end]
			step := 4
			if wide {
				step = 8
			}
			var units []uint16
			for k := 0; k+step <= len(hex); k += step {
				v, err := strconv.ParseUint(hex[k:k+step], 16, 32)
				if err != nil {
					continue
				}
				if wide {
					b.WriteRune(rune(v))
				} else {
					units = append(units, uint16(v))
				}
			}
			if !wide {
				b.WriteString(string(utf16.Decode(units)))
			}
			i += 4 + end + 3
		case strings.HasPrefix(s[i:], "\\X\\") && i+5 <= len(s):
			if v, err := strconv.ParseUint(s[i+3:i+5], 16, 8); err == nil {
				b.WriteRune(rune(v))
			}
			i += 4
		case strings.HasPrefix(s[i:], "\\S\\") && i+3 < len(s):
			b.WriteRune(rune(s[i+3]) + 128)
			i += 3
		case strings.HasPrefix(s[i:], "\\P") && i+3 < len(s) && s[i+3] == '\\':
			i += 3
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// parseP21 reads the entity instances of the DATA section.
func parseP21(text string) (*p21file, error) {
	d := strings.Index(text, "DATA;")
	if d < 0 {
		return nil, errNotP21
	}
	r := &p21Reader{s: text, i: d + len("DATA;")}
	f := &p21file{inst: map[int]*instance{}, invisible: map[int]bool{}, curves: map[int]*cad.Path{},
		scadec: strings.Contains(text[:d], "SCADEC")}
	for {
		r.skip()
		if r.i >= len(r.s) {
			break
		}
		if strings.HasPrefix(r.s[r.i:], "ENDSEC") {
			break
		}
		if r.s[r.i] != '#' {
			return f, fmt.Errorf("an entity instance is expected at byte %d", r.i)
		}
		r.i++
		j := r.i
		for j < len(r.s) && r.s[j] >= '0' && r.s[j] <= '9' {
			j++
		}
		id, _ := strconv.Atoi(r.s[r.i:j])
		r.i = j
		r.skip()
		if r.i >= len(r.s) || r.s[r.i] != '=' {
			return f, fmt.Errorf("#%d: '=' is missing", id)
		}
		r.i++
		r.skip()
		in := &instance{}
		if r.i < len(r.s) && r.s[r.i] == '(' {
			// a complex instance: records in parentheses
			r.i++
			for {
				r.skip()
				if r.i >= len(r.s) || r.s[r.i] == ')' {
					r.i++
					break
				}
				name := r.keyword()
				if name == "" {
					return f, fmt.Errorf("#%d: a record name is missing", id)
				}
				in.recs = append(in.recs, record{name: name, params: r.params(0)})
				if r.err != nil {
					return f, fmt.Errorf("#%d: %w", id, r.err)
				}
			}
		} else {
			name := r.keyword()
			in.recs = append(in.recs, record{name: name, params: r.params(0)})
			if r.err != nil {
				return f, fmt.Errorf("#%d: %w", id, r.err)
			}
		}
		r.skip()
		if r.i < len(r.s) && r.s[r.i] == ';' {
			r.i++
		}
		f.inst[id] = in
	}
	if len(f.inst) == 0 {
		return nil, errNotP21
	}
	return f, nil
}

func (f *p21file) at(v pval) *instance {
	if v.k != 'r' {
		return nil
	}
	return f.inst[v.ref]
}

// refs returns the references of a list.
func refs(v pval) []int {
	var out []int
	for _, x := range v.list {
		if x.k == 'r' {
			out = append(out, x.ref)
		}
	}
	return out
}

// number returns a number, or the first number of a typed value or list
// (SXF writes PARAMETER_VALUE((0.0)) as well as PARAMETER_VALUE(0.0)).
func number(v pval) float64 {
	for range 8 {
		if v.k == 'n' {
			return v.num
		}
		if (v.k != 't' && v.k != 'l') || len(v.list) == 0 {
			return 0
		}
		v = v.list[0]
	}
	return 0
}

// point reads a CARTESIAN_POINT.
func (f *p21file) point(v pval) (cad.Point, bool) {
	r := f.at(v).get("CARTESIAN_POINT")
	if r == nil {
		return cad.Point{}, false
	}
	c := r.p(1).list
	if len(c) < 2 {
		return cad.Point{}, false
	}
	return cad.Point{X: c[0].num, Y: c[1].num}, true
}

// direction reads a DIRECTION (normalized; the x axis when missing).
func (f *p21file) direction(v pval) cad.Point {
	r := f.at(v).get("DIRECTION")
	if r == nil {
		return cad.Point{X: 1}
	}
	c := r.p(1).list
	if len(c) < 2 {
		return cad.Point{X: 1}
	}
	d := cad.Point{X: c[0].num, Y: c[1].num}
	if l := d.Len(); l > 0 {
		return d.Mul(1 / l)
	}
	return cad.Point{X: 1}
}

// placement reads an AXIS2_PLACEMENT_2D as a transform.
func (f *p21file) placement(v pval) canvas.Matrix {
	r := f.at(v).get("AXIS2_PLACEMENT_2D")
	if r == nil {
		return canvas.Identity
	}
	o, _ := f.point(r.p(1))
	d := f.direction(r.p(2))
	return canvas.Matrix{d.X, d.Y, -d.Y, d.X, o.X, o.Y}
}

func invert(m canvas.Matrix) canvas.Matrix {
	det := m[0]*m[3] - m[1]*m[2]
	if det == 0 {
		return canvas.Identity
	}
	a, b, c, d := m[3]/det, -m[1]/det, -m[2]/det, m[0]/det
	return canvas.Matrix{a, b, c, d, -(a*m[4] + c*m[5]), -(b*m[4] + d*m[5])}
}

// colour reads a colour.
func (f *p21file) colour(v pval) bdf.Color {
	in := f.at(v)
	if r := in.get("DRAUGHTING_PRE_DEFINED_COLOUR"); r != nil {
		for i, n := range colorNames {
			if strings.EqualFold(n, strings.TrimSpace(r.p(0).str)) {
				c := colorRGB[i]
				return bdf.RGB(c[0], c[1], c[2])
			}
		}
	}
	if r := in.get("COLOUR_RGB"); r != nil {
		c := func(x float64) uint8 { return uint8(math.Round(math.Min(math.Max(x, 0), 1) * 255)) }
		// the predefined colors that are not ISO's ones are written as RGB
		// with their SXF name
		name := strings.TrimPrefix(strings.TrimSpace(r.p(0).str), "$$SXF_")
		for i, n := range colorNames {
			if strings.EqualFold(n, name) {
				k := colorRGB[i]
				return bdf.RGB(k[0], k[1], k[2])
			}
		}
		return bdf.RGB(c(r.p(1).num), c(r.p(2).num), c(r.p(3).num))
	}
	return bdf.RGB(255, 255, 255)
}

// lengthValue reads a length: a number, a typed measure or a
// LENGTH_MEASURE_WITH_UNIT.
func (f *p21file) lengthValue(v pval) float64 {
	for range 8 {
		switch v.k {
		case 'n':
			return v.num
		case 't':
			if len(v.list) == 0 {
				return 0
			}
			v = v.list[0]
		case 'r':
			r := f.at(v).get("LENGTH_MEASURE_WITH_UNIT")
			if r == nil {
				return 0
			}
			v = r.p(0)
		default:
			return 0
		}
	}
	return 0
}

// curveStyle reads a CURVE_STYLE (name, font, width, colour).
func (f *p21file) curveStyle(in *instance) lineStyle {
	r := in.get("CURVE_STYLE")
	if r == nil {
		return lineStyle{color: bdf.RGB(255, 255, 255), width: 0.25}
	}
	w := f.lengthValue(r.p(2))
	if w <= 0 {
		w = 0.25
	}
	s := lineStyle{color: f.colour(r.p(3)), width: w}
	font := f.at(r.p(1))
	if fr := font.get("DRAUGHTING_PRE_DEFINED_CURVE_FONT"); fr != nil {
		name := strings.ToLower(strings.TrimSpace(fr.p(0).str))
		for i, n := range lineTypeNames {
			if n == name && i > 0 {
				scale := math.Max(w, 0.13) / 0.5
				for k, v := range lineTypePitch[i] {
					if k%2 == 1 {
						v = -v
					}
					s.dash = append(s.dash, v*scale)
				}
			}
		}
	} else if fr := font.get("CURVE_STYLE_FONT"); fr != nil {
		for _, p := range refs(fr.p(1)) {
			if pr := f.inst[p].get("CURVE_STYLE_FONT_PATTERN"); pr != nil {
				s.dash = append(s.dash, math.Abs(pr.p(0).num), -math.Abs(pr.p(1).num))
			}
		}
	}
	return s
}

// styles returns the styles of a styled item: the members of its
// presentation style assignments.
func (f *p21file) styles(r *record) []*instance {
	var out []*instance
	for _, a := range refs(r.p(0)) {
		if psa := f.inst[a].get("PRESENTATION_STYLE_ASSIGNMENT"); psa != nil {
			for _, s := range psa.p(0).list {
				if s.k == 'r' {
					out = append(out, f.inst[s.ref])
				}
			}
		}
	}
	return out
}

// styledItem returns the STYLED_ITEM record of an instance. A styled item
// has a name first in AP214 and later; AP202's has none.
func styledItem(in *instance) (*record, int) {
	r := in.get("STYLED_ITEM")
	if r == nil {
		return nil, 0
	}
	if len(r.params) >= 3 {
		return r, 1
	}
	return r, 0
}

// curve returns the path of a curve in its own coordinates.
func (f *p21file) curve(v pval, depth int) *cad.Path {
	// a curve used again (by several items or composite curves) is made
	// once; the arcs of angular dimensions are made apart
	if v.k != 'r' || f.ccw {
		return f.makeCurve(v, depth)
	}
	if c, ok := f.curves[v.ref]; ok {
		return c
	}
	f.curves[v.ref] = nil // a curve made of itself is none
	c := f.makeCurve(v, depth)
	f.curves[v.ref] = c
	return c
}

func (f *p21file) makeCurve(v pval, depth int) *cad.Path {
	in := f.at(v)
	if in == nil || depth > 32 {
		return nil
	}
	p := &cad.Path{}
	switch {
	case in.get("POLYLINE") != nil:
		var pts []cad.Point
		for _, x := range in.get("POLYLINE").p(1).list {
			if q, ok := f.point(x); ok {
				pts = append(pts, q)
			}
		}
		return p.Polyline(pts, false)
	case in.get("TRIMMED_CURVE") != nil:
		return f.trimmed(in.get("TRIMMED_CURVE"), depth)
	case in.get("CIRCLE") != nil:
		r := in.get("CIRCLE")
		m := f.placement(r.p(1))
		return p.Circle(cad.Point{}, r.p(2).num).Transform(m)
	case in.get("ELLIPSE") != nil:
		r := in.get("ELLIPSE")
		m := f.placement(r.p(1))
		a, b := r.p(2).num, r.p(3).num
		ratio := 1.0
		if a != 0 {
			ratio = b / a
		}
		return p.Ellipse(cad.Point{}, cad.Point{X: a}, ratio).Transform(m)
	case in.get("BEZIER_CURVE") != nil:
		var pts []cad.Point
		for _, x := range in.get("BEZIER_CURVE").p(2).list {
			if q, ok := f.point(x); ok {
				pts = append(pts, q)
			}
		}
		return bezierPath(pts)
	case in.get("B_SPLINE_CURVE_WITH_KNOTS") != nil:
		r := in.get("B_SPLINE_CURVE_WITH_KNOTS")
		var pts []cad.Point
		for _, x := range r.p(2).list {
			if q, ok := f.point(x); ok {
				pts = append(pts, q)
			}
		}
		var knots []float64
		mult := r.p(6).list
		for i, k := range r.p(7).list {
			n := 1
			if i < len(mult) {
				n = int(mult[i].num)
			}
			for range min(n, 64) {
				knots = append(knots, k.num)
			}
		}
		return p.BSpline(int(r.p(1).num), knots, pts, nil)
	case in.get("COMPOSITE_CURVE") != nil:
		// composite curves are joined as points, bounded
		for _, s := range in.get("COMPOSITE_CURVE").p(1).list {
			seg := f.at(s).get("COMPOSITE_CURVE_SEGMENT")
			if seg == nil {
				continue
			}
			c := f.curve(seg.p(2), depth+1)
			if c == nil {
				continue
			}
			pts := flatten(c)
			if p.Size()+len(pts) > maxCurvePoints {
				f.r.warn("longcurve", "composite curves of more than %d points are cut", maxCurvePoints)
				break
			}
			if seg.p(1).str == "F" {
				for i, j := 0, len(pts)-1; i < j; i, j = i+1, j-1 {
					pts[i], pts[j] = pts[j], pts[i]
				}
			}
			for i, q := range pts {
				if i == 0 && !p.Open() {
					p.MoveTo(q.X, q.Y)
				} else {
					p.LineTo(q.X, q.Y)
				}
			}
		}
		return p
	}
	return nil
}

// trimmed returns a trimmed line, circle or ellipse.
func (f *p21file) trimmed(r *record, depth int) *cad.Path {
	basis := f.at(r.p(1))
	sense := r.p(4).str != "F" || f.ccw
	param := func(v pval) (float64, cad.Point, bool, bool) {
		var t float64
		var pt cad.Point
		hasT, hasP := false, false
		for _, x := range v.list {
			switch {
			case x.k == 't' && x.str == "PARAMETER_VALUE":
				t, hasT = number(x), true
			case x.k == 'r':
				if q, ok := f.point(x); ok {
					pt, hasP = q, true
				}
			}
		}
		return t, pt, hasT, hasP
	}
	t1, p1, hasT1, hasP1 := param(r.p(2))
	t2, p2, hasT2, hasP2 := param(r.p(3))
	p := &cad.Path{}
	switch {
	case basis.get("LINE") != nil:
		lr := basis.get("LINE")
		o, _ := f.point(lr.p(1))
		vec := f.at(lr.p(2)).get("VECTOR")
		d := f.direction(vec.p(1)).Mul(vec.p(2).num)
		if !hasP1 && hasT1 {
			p1 = o.Add(d.Mul(t1))
		}
		if !hasP2 && hasT2 {
			p2 = o.Add(d.Mul(t2))
		}
		return p.MoveTo(p1.X, p1.Y).LineTo(p2.X, p2.Y)
	case basis.get("CIRCLE") != nil || basis.get("ELLIPSE") != nil:
		var m canvas.Matrix
		a, b := 0.0, 0.0
		if cr := basis.get("CIRCLE"); cr != nil {
			m = f.placement(cr.p(1))
			a, b = cr.p(2).num, cr.p(2).num
		} else {
			er := basis.get("ELLIPSE")
			m = f.placement(er.p(1))
			a, b = er.p(2).num, er.p(3).num
		}
		inv := invert(m)
		if !hasT1 && hasP1 {
			q := cad.Apply(inv, p1)
			t1 = math.Atan2(q.Y/math.Max(b, 1e-12), q.X/math.Max(a, 1e-12))
		}
		if !hasT2 && hasP2 {
			q := cad.Apply(inv, p2)
			t2 = math.Atan2(q.Y/math.Max(b, 1e-12), q.X/math.Max(a, 1e-12))
		}
		if math.IsNaN(t1) || math.IsInf(t1, 0) || math.IsNaN(t2) || math.IsInf(t2, 0) {
			return nil
		}
		t1, t2 = arcEnds(t1, t2, !sense)
		ratio := 1.0
		if a != 0 {
			ratio = b / a
		}
		return p.EllipseArc(cad.Point{}, cad.Point{X: a}, ratio, t1, t2).Transform(m)
	case basis.get("CLOTHOID") != nil:
		cr := basis.get("CLOTHOID")
		m := f.placement(cr.p(1))
		// SXF writes the lengths in the name of the trimmed curve and the
		// direction in the clothoid's
		var s0, s1 float64
		if parts := strings.Split(strings.TrimSpace(r.p(0).str), ","); len(parts) == 2 {
			s0, _ = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			s1, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		} else {
			s0, s1 = t1, t2
		}
		cw := strings.TrimSpace(cr.p(0).str) == "1"
		return p.Polyline(clothoid(cr.p(2).num, s0, s1, cw), false).Transform(m)
	}
	return nil
}

// draw draws the drawing sheet and returns its size and the title.
func (f *p21file) draw(r *renderer) (w, h float64, title string) {
	f.r = r
	ids := make([]int, 0, len(f.inst))
	for id := range f.inst {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	var sheet *record
	hidden := map[int]bool{}
	for _, id := range ids {
		in := f.inst[id]
		if t := in.get("DRAUGHTING_TITLE"); t != nil && title == "" {
			title = t.p(2).str
		}
		if s := in.get("DRAWING_SHEET_REVISION"); s != nil && sheet == nil {
			sheet = s
		}
		if iv := in.get("INVISIBILITY"); iv != nil {
			for _, id := range refs(iv.p(0)) {
				if hidden[id] {
					continue
				}
				hidden[id] = true
				if la := f.inst[id].get("PRESENTATION_LAYER_ASSIGNMENT"); la != nil {
					for _, item := range refs(la.p(2)) {
						f.invisible[item] = true
					}
				} else {
					f.invisible[id] = true
				}
			}
		}
		for _, rec := range in.recs {
			if strings.HasSuffix(rec.name, "_REPRESENTATION") {
				if c, ok := backgroundAttribute(stripSXFPrefix(rec.p(0).str)); ok && !r.light {
					r.bg = c
				}
			}
		}
	}
	if sheet == nil {
		r.warn("nosheet", "the file has no drawing sheet; nothing is drawn")
		return 0, 0, title
	}
	items := refs(sheet.p(1))
	for _, id := range items {
		if b := f.inst[id].get("PLANAR_BOX"); b != nil {
			w, h = b.p(1).num, b.p(2).num
		}
	}
	f.items(items, canvas.Identity, map[int]bool{}, 0)
	return w, h, title
}

// stripSXFPrefix removes the kind of a compound figure from its name
// ("$$SXF_FM_", "$$SXF_FG_", "$$SXF_G_", "$$SXF_P_").
func stripSXFPrefix(name string) string {
	for _, p := range []string{"$$SXF_FM_", "$$SXF_FG_", "$$SXF_G_", "$$SXF_P_"} {
		if strings.HasPrefix(name, p) {
			return name[len(p):]
		}
	}
	return name
}

// items draws annotation occurrences and callouts through m; seen keeps
// the instances drawn in this frame (a dimension's text is in two
// callouts).
// depth counts the callouts and compound figures the items are in.
func (f *p21file) items(ids []int, m canvas.Matrix, seen map[int]bool, depth int) {
	for _, id := range ids {
		f.item(id, m, seen, depth)
	}
}

// maxDepth bounds how deep callouts and compound figures are followed.
const maxDepth = 64

func (f *p21file) item(id int, m canvas.Matrix, seen map[int]bool, depth int) {
	if seen[id] || f.invisible[id] || !f.r.visit() {
		return
	}
	if depth > maxDepth {
		f.r.warn("depth", "callouts and compound figures nested deeper than %d levels are not drawn", maxDepth)
		return
	}
	seen[id] = true
	in := f.inst[id]
	if in == nil {
		return
	}
	if c := in.get("DRAUGHTING_CALLOUT"); c != nil {
		contents := refs(c.own(1).p(0))
		if in.get("ANGULAR_DIMENSION") != nil {
			for _, id := range contents {
				f.ccw = f.inst[id].get("DIMENSION_CURVE") != nil
				f.item(id, m, seen, depth+1)
			}
			f.ccw = false
			return
		}
		f.items(contents, m, seen, depth+1)
		return
	}
	si, k := styledItem(in)
	if si == nil {
		return
	}
	item := si.p(k + 1)
	styles := f.styles(&record{params: si.params[k:]})
	target := f.at(item)
	switch {
	case target.get("DRAUGHTING_CALLOUT") != nil:
		// a callout styled with NULL_STYLE; its members carry the styles
		f.item(item.ref, m, seen, depth+1)
	case in.get("ANNOTATION_FILL_AREA_OCCURRENCE") != nil:
		f.fill(target, styles, m)
	case target.get("MAPPED_ITEM") != nil:
		f.subfigure(target.get("MAPPED_ITEM").own(2), m, depth)
	case target.get("DEFINED_SYMBOL") != nil:
		f.symbol(in, target.get("DEFINED_SYMBOL"), styles, m)
	case target.get("TEXT_LITERAL_WITH_EXTENT") != nil || target.get("TEXT_LITERAL") != nil:
		f.text(target, styles, m)
	default:
		path := f.curve(item, 0)
		if path == nil {
			return
		}
		st := lineStyle{color: bdf.RGB(255, 255, 255), width: 0.25}
		for _, s := range styles {
			if s.get("CURVE_STYLE") != nil {
				st = f.curveStyle(s)
			}
		}
		f.r.stroke(path, st, m)
	}
}

func (f *p21file) subfigure(mi *record, m canvas.Matrix, depth int) {
	src := f.at(mi.p(0)).get("SYMBOL_REPRESENTATION_MAP")
	tgt := f.at(mi.p(1)).get("SYMBOL_TARGET")
	if src == nil {
		return
	}
	rep := f.at(src.p(1))
	if rep == nil || len(rep.recs) == 0 {
		return
	}
	rr := &rep.recs[0]
	for i := range rep.recs {
		if strings.HasSuffix(rep.recs[i].name, "REPRESENTATION") {
			rr = &rep.recs[i]
		}
	}
	name := rr.p(0).str
	if _, ok := backgroundAttribute(stripSXFPrefix(name)); ok {
		return // a group that only carries the background color attribute
	}
	if slices.Contains(f.stack, mi.p(0).ref) || len(f.stack) > 32 {
		f.r.warn("cycle", "compound figures that place themselves are drawn once")
		return
	}
	sm := m
	if tgt != nil {
		sm = sm.Mul(f.placement(tgt.p(1))).Mul(canvas.Scale(tgt.p(2).num, tgt.p(3).num))
	}
	sm = sm.Mul(invert(f.placement(src.p(0))))
	if strings.HasPrefix(name, "$$SXF_FG_") {
		sm = sm.Mul(canvas.Matrix{0, 1, 1, 0, 0, 0})
	}
	f.stack = append(f.stack, mi.p(0).ref)
	f.items(refs(rr.p(1)), sm, map[int]bool{}, depth+1)
	f.stack = f.stack[:len(f.stack)-1]
}

// capped returns the first maxList items of a list.
func capped[T any](list []T, r *renderer) []T {
	if len(list) > maxList {
		r.warn("longlist", "only the first %d outlines and styles of a fill are drawn", maxList)
		return list[:maxList]
	}
	return list
}

// leaderDirection returns the direction of the end of an annotation curve
// at the point at, out of the curve.
func (f *p21file) leaderDirection(v pval, at cad.Point) (cad.Point, bool) {
	si, k := styledItem(f.at(v))
	if si == nil {
		return cad.Point{}, false
	}
	c := f.curve(si.p(k+1), 0)
	if c == nil {
		return cad.Point{}, false
	}
	pts := flatten(c)
	if len(pts) < 2 {
		return cad.Point{}, false
	}
	a, b := pts[0], pts[1]
	if n := len(pts); pts[n-1].Sub(at).Len() < a.Sub(at).Len() {
		a, b = pts[n-1], pts[n-2]
	}
	d := a.Sub(b)
	if l := d.Len(); l > 0 {
		return d.Mul(1 / l), true
	}
	return cad.Point{}, false
}

func (f *p21file) symbolColour(styles []*instance) bdf.Color {
	for _, s := range styles {
		if ss := s.get("SYMBOL_STYLE"); ss != nil {
			if sc := f.at(ss.p(1)).get("SYMBOL_COLOUR"); sc != nil {
				return f.colour(sc.p(0))
			}
		}
	}
	return bdf.RGB(255, 255, 255)
}

func (f *p21file) symbol(in *instance, ds *record, styles []*instance, m canvas.Matrix) {
	def := f.at(ds.p(1))
	tgt := f.at(ds.p(2)).get("SYMBOL_TARGET")
	if tgt == nil {
		return
	}
	pl := f.at(tgt.p(1)).get("AXIS2_PLACEMENT_2D")
	if pl == nil {
		return
	}
	at, _ := f.point(pl.p(1))
	dir := f.direction(pl.p(2))
	if in.get("LEADER_TERMINATOR") != nil {
		// a leader's arrow points away from the leader (writers disagree on
		// the direction they give it)
		if d, ok := f.leaderDirection(in.get("TERMINATOR_SYMBOL").p(0), at); ok {
			dir = d
		}
	}
	scale := tgt.p(2).num
	col := f.symbolColour(styles)
	switch {
	case def.get("PRE_DEFINED_TERMINATOR_SYMBOL") != nil:
		// drawn with the width of the curve it ends
		w := 0.25
		if si, k := styledItem(f.at(in.get("TERMINATOR_SYMBOL").p(0))); si != nil {
			for _, s := range f.styles(&record{params: si.params[k:]}) {
				if s.get("CURVE_STYLE") != nil {
					w = f.curveStyle(s).width
				}
			}
		}
		f.r.arrow(arrowCode(def.get("PRE_DEFINED_TERMINATOR_SYMBOL").p(0).str), at, dir, scale, lineStyle{color: col, width: w}, m)
	case def.get("PRE_DEFINED_POINT_MARKER_SYMBOL") != nil:
		name := strings.ToLower(strings.TrimSpace(def.get("PRE_DEFINED_POINT_MARKER_SYMBOL").p(0).str))
		code := 7
		for i, n := range markerNames {
			if n == name {
				code = i + 1
			}
		}
		f.r.marker(code, at, math.Atan2(dir.Y, dir.X)*180/math.Pi, scale, col, m)
	default:
		f.r.warn("symbol", "predefined symbols (externally_defined_symbol) are not drawn")
	}
}

func (f *p21file) text(in *instance, styles []*instance, m canvas.Matrix) {
	r := in.get("TEXT_LITERAL_WITH_EXTENT")
	if r == nil {
		r = in.get("TEXT_LITERAL")
	}
	pl := f.at(r.p(2)).get("AXIS2_PLACEMENT_2D")
	if pl == nil {
		return
	}
	at, _ := f.point(pl.p(1))
	dir := f.direction(pl.p(2))
	t := textBox{s: r.p(1).str, x: at.X, y: at.Y, angle: math.Atan2(dir.Y, dir.X) * 180 / math.Pi, color: bdf.RGB(255, 255, 255)}
	// the anchor: SXF names it ("$$SXF_topline left", ...) and places the
	// text there; other files only align the text on its baseline
	row, col := -1, 0
	if name, ok := strings.CutPrefix(r.p(0).str, "$$SXF_"); ok {
		v, h, _ := strings.Cut(strings.ToLower(name), " ")
		row = slices.Index([]string{"baseline", "middleline", "topline"}, v)
		col = max(slices.Index([]string{"left", "centre", "right"}, h), 0)
	}
	if row < 0 {
		al := strings.ToLower(r.p(3).str)
		switch {
		case strings.Contains(al, "right"):
			col = 2
		case strings.Contains(al, "centre") || strings.Contains(al, "center"):
			col = 1
		}
		t.onBaseline = true
	}
	t.anchor = 1 + col + 3*max(row, 0)
	t.vertical = r.p(4).str == "DOWN"
	font := f.at(r.p(5))
	if fr := font.get("EXTERNALLY_DEFINED_TEXT_FONT"); fr != nil {
		if id := fr.p(0); id.k == 't' && len(id.list) > 0 {
			t.font = id.list[0].str
		} else {
			t.font = id.str
		}
	} else if fr := font.get("DRAUGHTING_PRE_DEFINED_TEXT_FONT"); fr != nil {
		t.font = fr.p(0).str
	}
	if ext := f.at(r.p(6)).get("PLANAR_EXTENT"); ext != nil {
		t.w, t.h = ext.p(1).num, ext.p(2).num
	}
	for _, s := range styles {
		// the appearance is TEXT_STYLE's, which a simple instance of a
		// subtype holds after the name
		for _, rec := range s.recs {
			if rec.name == "TEXT_STYLE" || strings.HasPrefix(rec.name, "TEXT_STYLE_WITH_") && len(rec.params) >= 3 {
				if fd := f.at(rec.p(1)).get("TEXT_STYLE_FOR_DEFINED_FONT"); fd != nil {
					t.color = f.colour(fd.p(0))
				}
			}
		}
		if bc := s.get("TEXT_STYLE_WITH_BOX_CHARACTERISTICS"); bc != nil {
			for _, c := range bc.own(1).p(0).list {
				if c.k != 't' || len(c.list) == 0 {
					continue
				}
				switch c.str {
				case "BOX_HEIGHT":
					if t.h == 0 {
						t.h = c.list[0].num
					}
				case "BOX_SLANT_ANGLE":
					t.slant = c.list[0].num * 180 / math.Pi
				}
			}
		}
		if sp := s.get("TEXT_STYLE_WITH_SPACING"); sp != nil {
			t.spacing = f.lengthValue(sp.own(1).p(0))
		}
	}
	if t.w == 0 {
		t.w = t.h * float64(len([]rune(t.s)))
	}
	if row >= 0 && f.scadec {
		// the SCADEC library, which most SXF writers use, moves the point
		// off the anchor: up half the height on the bottom row, down half
		// of it on the middle row and down all of it on the top row
		k := []float64{0.5, -0.5, -1}[row] * t.h
		a := t.angle * math.Pi / 180
		t.x += math.Sin(a) * k
		t.y -= math.Cos(a) * k
	}
	f.r.text(t, m)
}

func (f *p21file) fill(area *instance, styles []*instance, m canvas.Matrix) {
	a := area.get("ANNOTATION_FILL_AREA")
	if a == nil {
		return
	}
	bound := &cad.Path{}
	for _, b := range capped(a.p(1).list, f.r) {
		c := f.curve(b, 0)
		if c == nil {
			continue
		}
		pts := flatten(c)
		if len(pts) >= 3 {
			bound.Polyline(pts, true)
		}
	}
	if bound.Empty() {
		return
	}
	for _, s := range styles {
		fs := s.get("FILL_AREA_STYLE")
		if fs == nil {
			continue
		}
		var lines []hatchLine
		for _, st := range capped(fs.p(1).list, f.r) {
			x := f.at(st)
			switch {
			case x.get("FILL_AREA_STYLE_COLOUR") != nil:
				if f.r.budget(bound.Size()) {
					f.r.out.Fill(bound.Transform(m), cad.Fill{Color: f.r.paint(f.colour(x.get("FILL_AREA_STYLE_COLOUR").p(1)))}, true)
				}
			case x.get("FILL_AREA_STYLE_HATCHING") != nil:
				hr := x.get("FILL_AREA_STYLE_HATCHING")
				ls := f.curveStyle(f.at(hr.p(1)))
				var spacing, dirAngle float64
				if rf := f.at(hr.p(2)).get("ONE_DIRECTION_REPEAT_FACTOR"); rf != nil {
					if vec := f.at(rf.p(1)).get("VECTOR"); vec != nil {
						d := f.direction(vec.p(1)).Mul(vec.p(2).num)
						spacing = d.Len()
						dirAngle = math.Atan2(d.Y, d.X)
					}
				}
				ref, _ := f.point(hr.p(3))
				angle := hr.p(5).num
				// the spacing across the lines
				spacing *= math.Abs(math.Sin(dirAngle - angle))
				lines = append(lines, hatchLine{style: ls, x: ref.X, y: ref.Y, spacing: spacing, angle: angle * 180 / math.Pi})
			case x.get("EXTERNALLY_DEFINED_HATCH_STYLE") != nil:
				name := x.get("EXTERNALLY_DEFINED_HATCH_STYLE").p(0)
				n := name.str
				if name.k == 't' && len(name.list) > 0 {
					n = name.list[0].str
				}
				if n == "Area_control" {
					if f.r.budget(bound.Size()) {
						f.r.out.Fill(bound.Transform(m), cad.Fill{Color: f.r.bg}, true)
					}
				} else {
					f.r.warn("hatch:"+n, "the predefined hatch %q is not drawn", n)
				}
			case x.get("FILL_AREA_STYLE_TILES") != nil:
				f.r.warn("tiles", "hatches of tiled symbols are not drawn")
			}
		}
		if len(lines) > 0 {
			f.r.hatch(bound, lines, m)
		}
	}
}

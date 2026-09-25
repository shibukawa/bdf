package bdf

// Path verbs.
const (
	VerbMove      byte = 0
	VerbLine      byte = 1
	VerbQuad      byte = 2
	VerbCubic     byte = 3
	VerbClose     byte = 4
	VerbRect      byte = 5
	VerbEllipse   byte = 6
	VerbArcTo     byte = 7
	VerbRoundRect byte = 8
)

// verbArgs is the number of f32 operands per verb.
var verbArgs = [...]int{2, 2, 4, 6, 0, 4, 8, 5, 5}

// Path is a geometric path that maps to a Path2D.
type Path struct {
	Verbs []byte
	Args  []float32
}

func (p *Path) add(v byte, args ...float32) *Path {
	p.Verbs = append(p.Verbs, v)
	p.Args = append(p.Args, args...)
	return p
}

// MoveTo starts a new sub-path.
func (p *Path) MoveTo(x, y float32) *Path { return p.add(VerbMove, x, y) }

// LineTo adds a straight segment.
func (p *Path) LineTo(x, y float32) *Path { return p.add(VerbLine, x, y) }

// QuadTo adds a quadratic bezier.
func (p *Path) QuadTo(cx, cy, x, y float32) *Path { return p.add(VerbQuad, cx, cy, x, y) }

// CubicTo adds a cubic bezier.
func (p *Path) CubicTo(c1x, c1y, c2x, c2y, x, y float32) *Path {
	return p.add(VerbCubic, c1x, c1y, c2x, c2y, x, y)
}

// Close closes the current sub-path.
func (p *Path) Close() *Path { return p.add(VerbClose) }

// Rect adds a rectangle sub-path.
func (p *Path) Rect(x, y, w, h float32) *Path { return p.add(VerbRect, x, y, w, h) }

// Ellipse adds an elliptical arc (ctx.ellipse semantics, ccw as 0/1).
func (p *Path) Ellipse(x, y, rx, ry, rot, a0, a1 float32, ccw bool) *Path {
	var c float32
	if ccw {
		c = 1
	}
	return p.add(VerbEllipse, x, y, rx, ry, rot, a0, a1, c)
}

// Circle adds a full circle.
func (p *Path) Circle(x, y, r float32) *Path {
	return p.Ellipse(x, y, r, r, 0, 0, 6.2831855, false)
}

// ArcTo adds an arc tangent to two segments.
func (p *Path) ArcTo(x1, y1, x2, y2, r float32) *Path { return p.add(VerbArcTo, x1, y1, x2, y2, r) }

// RoundRect adds a rounded rectangle with a uniform radius.
func (p *Path) RoundRect(x, y, w, h, r float32) *Path { return p.add(VerbRoundRect, x, y, w, h, r) }

func (p *Path) encode(w *buf) {
	w.varuint(uint64(len(p.Verbs)))
	w.bytes(p.Verbs)
	w.f32s(p.Args...)
}

func decodePath(r *reader) *Path {
	n := int(r.varuint())
	if n > len(r.b) {
		r.fail("bad path verb count")
		return nil
	}
	p := &Path{Verbs: append([]byte(nil), r.bytes(n)...)}
	total := 0
	for _, v := range p.Verbs {
		if int(v) >= len(verbArgs) {
			r.fail("unknown path verb")
			return nil
		}
		total += verbArgs[v]
	}
	p.Args = make([]float32, total)
	for i := range p.Args {
		p.Args[i] = r.f32()
	}
	return p
}

// EncodePathCollection encodes paths as a Path collection part (type "path").
func EncodePathCollection(paths []*Path) []byte {
	var w buf
	w.varuint(uint64(len(paths)))
	for _, p := range paths {
		p.encode(&w)
	}
	return w.b
}

// DecodePathCollection decodes a Path collection part.
func DecodePathCollection(data []byte) ([]*Path, error) {
	r := &reader{b: data}
	n := int(r.varuint())
	if n > len(data) {
		return nil, &FormatError{Msg: "bad path collection count"}
	}
	out := make([]*Path, 0, n)
	for i := 0; i < n && r.err == nil; i++ {
		out = append(out, decodePath(r))
	}
	return out, r.err
}

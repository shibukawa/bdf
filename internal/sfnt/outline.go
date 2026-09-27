package sfnt

// Pen receives the outline of a glyph in font units (y up). Every contour
// starts with MoveTo and ends with Close.
type Pen interface {
	MoveTo(x, y float64)
	LineTo(x, y float64)
	QuadTo(cx, cy, x, y float64)
	CubeTo(c1x, c1y, c2x, c2y, x, y float64)
	Close()
}

// Outline draws the outline of a glyph with pen; ok is false when the
// glyph cannot be read (a glyph without contours, such as a space, draws
// nothing and is ok).
func (o *Outlines) Outline(gid uint16, pen Pen) (ok bool) {
	if o.cff != nil {
		return o.cff.run(int(gid), penSink{pen})
	}
	if o.f.IsCFF {
		return false
	}
	return o.glyf(gid, pen, affine{1, 0, 0, 1, 0, 0}, 0)
}

// penSink passes what a charstring draws to a Pen.
type penSink struct{ p Pen }

func (s penSink) moveTo(x, y float64) { s.p.MoveTo(x, y) }
func (s penSink) lineTo(x, y float64) { s.p.LineTo(x, y) }
func (s penSink) closeContour()       { s.p.Close() }
func (s penSink) cubeTo(x1, y1, x2, y2, x3, y3 float64) {
	s.p.CubeTo(x1, y1, x2, y2, x3, y3)
}

// affine is the transform of a composite glyph's component: x' = a x + c y + e.
type affine struct{ a, b, c, d, e, f float64 }

func (m affine) apply(x, y float64) (float64, float64) {
	return m.a*x + m.c*y + m.e, m.b*x + m.d*y + m.f
}

func (m affine) mul(n affine) affine {
	return affine{
		m.a*n.a + m.c*n.b, m.b*n.a + m.d*n.b,
		m.a*n.c + m.c*n.d, m.b*n.c + m.d*n.d,
		m.a*n.e + m.c*n.f + m.e, m.b*n.e + m.d*n.f + m.f,
	}
}

// glyfData returns the glyf entry of a glyph (nil for an empty glyph).
func (o *Outlines) glyfData(gid uint16) ([]byte, bool) {
	glyf, loca, head := o.f.Tables["glyf"], o.f.Tables["loca"], o.f.Tables["head"]
	if glyf == nil || loca == nil || len(head) < 54 || int(gid) >= o.f.NumGlyphs {
		return nil, false
	}
	var s, e int
	if be16(head, 50) != 0 {
		if int(gid)*4+8 > len(loca) {
			return nil, false
		}
		s, e = int(be32(loca, int(gid)*4)), int(be32(loca, int(gid)*4+4))
	} else {
		if int(gid)*2+4 > len(loca) {
			return nil, false
		}
		s, e = int(be16(loca, int(gid)*2))*2, int(be16(loca, int(gid)*2+2))*2
	}
	if s == e {
		return nil, true
	}
	if s > e || e > len(glyf) || e-s < 10 {
		return nil, false
	}
	return glyf[s:e], true
}

// Composite glyph flags.
const (
	argsAreWords   = 0x0001
	argsAreXY      = 0x0002
	haveScale      = 0x0008
	moreComponents = 0x0020
	haveXYScale    = 0x0040
	haveTwoByTwo   = 0x0080
)

func (o *Outlines) glyf(gid uint16, pen Pen, m affine, depth int) bool {
	g, ok := o.glyfData(gid)
	if !ok {
		return false
	}
	if g == nil {
		return true
	}
	n := int(int16(be16(g, 0)))
	if n >= 0 {
		return simpleGlyph(g, n, pen, m)
	}
	if depth > 8 {
		return false
	}
	pos := 10
	for {
		if pos+4 > len(g) {
			return false
		}
		flags, comp := be16(g, pos), be16(g, pos+2)
		pos += 4
		var dx, dy float64
		if flags&argsAreWords != 0 {
			if pos+4 > len(g) {
				return false
			}
			dx, dy = float64(int16(be16(g, pos))), float64(int16(be16(g, pos+2)))
			pos += 4
		} else {
			if pos+2 > len(g) {
				return false
			}
			dx, dy = float64(int8(g[pos])), float64(int8(g[pos+1]))
			pos += 2
		}
		if flags&argsAreXY == 0 {
			dx, dy = 0, 0 // anchored by point numbers: rare, placed at the origin
		}
		c := affine{1, 0, 0, 1, dx, dy}
		f2 := func(i int) float64 { return float64(int16(be16(g, i))) / 16384 }
		switch {
		case flags&haveScale != 0:
			if pos+2 > len(g) {
				return false
			}
			c.a, c.d = f2(pos), f2(pos)
			pos += 2
		case flags&haveXYScale != 0:
			if pos+4 > len(g) {
				return false
			}
			c.a, c.d = f2(pos), f2(pos+2)
			pos += 4
		case flags&haveTwoByTwo != 0:
			if pos+8 > len(g) {
				return false
			}
			c.a, c.b, c.c, c.d = f2(pos), f2(pos+2), f2(pos+4), f2(pos+6)
			pos += 8
		}
		if !o.glyf(comp, pen, m.mul(c), depth+1) {
			return false
		}
		if flags&moreComponents == 0 {
			return true
		}
	}
}

// Simple glyph point flags.
const (
	onCurve  = 0x01
	xShort   = 0x02
	yShort   = 0x04
	repeat   = 0x08
	xSame    = 0x10 // or positive short x
	ySame    = 0x20 // or positive short y
	maxPoint = 1 << 16
)

func simpleGlyph(g []byte, contours int, pen Pen, m affine) bool {
	pos := 10
	if pos+2*contours+2 > len(g) {
		return false
	}
	ends := make([]int, contours)
	points := 0
	for i := range ends {
		ends[i] = int(be16(g, pos+2*i))
		if i > 0 && ends[i] < ends[i-1] || ends[i] >= maxPoint {
			return false
		}
		points = ends[i] + 1
	}
	pos += 2 * contours
	pos += 2 + int(be16(g, pos)) // instructions
	flags := make([]byte, points)
	for i := 0; i < points; {
		if pos >= len(g) {
			return false
		}
		f := g[pos]
		pos++
		flags[i] = f
		i++
		if f&repeat != 0 {
			if pos >= len(g) {
				return false
			}
			r := int(g[pos])
			pos++
			for ; r > 0 && i < points; r-- {
				flags[i] = f
				i++
			}
		}
	}
	coords := func(short, same byte) ([]float64, bool) {
		out := make([]float64, points)
		v := 0
		for i, f := range flags {
			switch {
			case f&short != 0:
				if pos >= len(g) {
					return nil, false
				}
				d := int(g[pos])
				pos++
				if f&same == 0 {
					d = -d
				}
				v += d
			case f&same == 0:
				if pos+2 > len(g) {
					return nil, false
				}
				v += int(int16(be16(g, pos)))
				pos += 2
			}
			out[i] = float64(v)
		}
		return out, true
	}
	xs, ok := coords(xShort, xSame)
	if !ok {
		return false
	}
	ys, ok := coords(yShort, ySame)
	if !ok {
		return false
	}
	start := 0
	for _, end := range ends {
		contour(pen, m, xs[start:end+1], ys[start:end+1], flags[start:end+1])
		start = end + 1
	}
	return true
}

// contour draws a TrueType contour of on- and off-curve points: two
// off-curve points in a row imply the on-curve point between them.
func contour(pen Pen, m affine, xs, ys []float64, flags []byte) {
	n := len(xs)
	if n == 0 {
		return
	}
	on := func(i int) bool { return flags[i]&onCurve != 0 }
	// Start at an on-curve point, or between the first two off-curve ones.
	first := -1
	for i := 0; i < n; i++ {
		if on(i) {
			first = i
			break
		}
	}
	var sx, sy float64
	if first < 0 {
		// every point is a control point: start between the last and the first
		sx, sy = (xs[0]+xs[n-1])/2, (ys[0]+ys[n-1])/2
		px, py := m.apply(sx, sy)
		pen.MoveTo(px, py)
		quadLoop(pen, m, xs, ys, on, 0, n, sx, sy)
		pen.Close()
		return
	}
	sx, sy = xs[first], ys[first]
	px, py := m.apply(sx, sy)
	pen.MoveTo(px, py)
	quadLoop(pen, m, xs, ys, on, first+1, n, sx, sy)
	pen.Close()
}

// quadLoop draws the n points of a contour from index from, wrapping
// around, and closes the contour at (sx, sy).
func quadLoop(pen Pen, m affine, xs, ys []float64, on func(int) bool, from, n int, sx, sy float64) {
	var cx, cy float64
	pending := false
	for k := 0; k < n; k++ {
		i := (from + k) % len(xs)
		x, y := xs[i], ys[i]
		if k == n-1 && on(i) && x == sx && y == sy && !pending {
			break
		}
		if on(i) {
			if pending {
				ax, ay := m.apply(cx, cy)
				bx, by := m.apply(x, y)
				pen.QuadTo(ax, ay, bx, by)
				pending = false
			} else {
				px, py := m.apply(x, y)
				pen.LineTo(px, py)
			}
			continue
		}
		if pending {
			mx, my := (cx+x)/2, (cy+y)/2
			ax, ay := m.apply(cx, cy)
			bx, by := m.apply(mx, my)
			pen.QuadTo(ax, ay, bx, by)
		}
		cx, cy, pending = x, y, true
	}
	if pending {
		ax, ay := m.apply(cx, cy)
		bx, by := m.apply(sx, sy)
		pen.QuadTo(ax, ay, bx, by)
	}
}

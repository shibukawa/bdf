package sfnt

import (
	"errors"
	"math"
)

// Rect is the bounding box of a glyph's outline in font units (y up).
type Rect struct {
	XMin, YMin, XMax, YMax float64
}

// Outlines finds the bounding boxes of the glyphs of a font: from the glyph
// headers of TrueType outlines, by running the charstrings of CFF ones.
// It holds no mutable state, so one Outlines may serve several goroutines.
type Outlines struct {
	f   *Font
	cff *cffGlyphs
}

// NewOutlines prepares the outlines of f for Bounds.
func NewOutlines(f *Font) *Outlines {
	o := &Outlines{f: f}
	if f.IsCFF {
		if c, err := parseCFFGlyphs(f.Tables["CFF "]); err == nil {
			o.cff = c
		}
	}
	return o
}

// Bounds returns the bounding box of a glyph's outline; ok is false for
// glyphs without one (spaces, glyphs that cannot be read).
func (o *Outlines) Bounds(gid uint16) (r Rect, ok bool) {
	if o.cff != nil {
		return o.cff.bounds(int(gid))
	}
	glyf, loca, head := o.f.Tables["glyf"], o.f.Tables["loca"], o.f.Tables["head"]
	if glyf == nil || loca == nil || len(head) < 54 || int(gid) >= o.f.NumGlyphs {
		return r, false
	}
	var s, e int
	if be16(head, 50) != 0 {
		s, e = int(be32(loca, int(gid)*4)), int(be32(loca, int(gid)*4+4))
	} else {
		s, e = int(be16(loca, int(gid)*2))*2, int(be16(loca, int(gid)*2+2))*2
	}
	if e-s < 10 || e > len(glyf) {
		return r, false
	}
	g := glyf[s:e]
	return Rect{float64(int16(be16(g, 2))), float64(int16(be16(g, 4))), float64(int16(be16(g, 6))), float64(int16(be16(g, 8)))}, true
}

// cffGlyphs is what running the charstrings of a CFF font needs.
type cffGlyphs struct {
	charStrings [][]byte
	global      [][]byte
	local       [][][]byte // Subrs of each Private DICT
	fdSelect    []byte     // Private DICT of each glyph (CID-keyed fonts)
}

// cffIndex reads an INDEX at pos and returns its items and the position
// after it.
func cffIndex(b []byte, pos int) ([][]byte, int, error) {
	if pos < 0 || pos+2 > len(b) {
		return nil, pos, errors.New("cff: truncated INDEX")
	}
	count := int(be16(b, pos))
	if count == 0 {
		return nil, pos + 2, nil
	}
	if pos+3 > len(b) {
		return nil, pos, errors.New("cff: truncated INDEX")
	}
	offSize := int(b[pos+2])
	if offSize < 1 || offSize > 4 {
		return nil, pos, errors.New("cff: bad offSize")
	}
	offs := pos + 3
	data := offs + (count+1)*offSize - 1
	if data > len(b) {
		return nil, pos, errors.New("cff: truncated INDEX")
	}
	off := func(i int) int {
		v := 0
		for k := 0; k < offSize; k++ {
			v = v<<8 | int(b[offs+i*offSize+k])
		}
		return v
	}
	items := make([][]byte, count)
	for i := range items {
		s, e := data+off(i), data+off(i+1)
		if s > e || e > len(b) {
			return nil, pos, errors.New("cff: bad INDEX offsets")
		}
		items[i] = b[s:e]
	}
	return items, data + off(count), nil
}

// cffDict reads a DICT into operator → operands (12 x is 1200+x).
func cffDict(b []byte) map[int][]float64 {
	out := map[int][]float64{}
	var ops []float64
	for i := 0; i < len(b); {
		c := int(b[i])
		switch {
		case c <= 21:
			op := c
			i++
			if c == 12 && i < len(b) {
				op = 1200 + int(b[i])
				i++
			}
			out[op] = ops
			ops = nil
		case c == 28 && i+2 < len(b):
			ops = append(ops, float64(int16(be16(b, i+1))))
			i += 3
		case c == 29 && i+4 < len(b):
			ops = append(ops, float64(int32(be32(b, i+1))))
			i += 5
		case c == 30:
			// real number: skip its nibbles (no offsets are reals)
			i++
			for i < len(b) {
				v := b[i]
				i++
				if v&0x0f == 0x0f || v>>4 == 0x0f {
					break
				}
			}
			ops = append(ops, 0)
		case c >= 32 && c <= 246:
			ops = append(ops, float64(c-139))
			i++
		case c >= 247 && c <= 250 && i+1 < len(b):
			ops = append(ops, float64((c-247)*256+int(b[i+1])+108))
			i += 2
		case c >= 251 && c <= 254 && i+1 < len(b):
			ops = append(ops, float64(-(c-251)*256-int(b[i+1])-108))
			i += 2
		default:
			i++
		}
	}
	return out
}

func parseCFFGlyphs(b []byte) (*cffGlyphs, error) {
	if len(b) < 4 {
		return nil, errors.New("cff: too short")
	}
	_, pos, err := cffIndex(b, int(b[2])) // Name INDEX
	if err != nil {
		return nil, err
	}
	tops, pos, err := cffIndex(b, pos)
	if err != nil || len(tops) == 0 {
		return nil, errors.New("cff: no Top DICT")
	}
	if _, pos, err = cffIndex(b, pos); err != nil { // String INDEX
		return nil, err
	}
	c := &cffGlyphs{}
	if c.global, _, err = cffIndex(b, pos); err != nil {
		return nil, err
	}
	top := cffDict(tops[0])
	if v := top[17]; len(v) == 1 {
		if c.charStrings, _, err = cffIndex(b, int(v[0])); err != nil {
			return nil, err
		}
	}
	if len(c.charStrings) == 0 {
		return nil, errors.New("cff: no CharStrings")
	}
	private := func(v []float64) [][]byte {
		if len(v) != 2 {
			return nil
		}
		size, off := int(v[0]), int(v[1])
		if off < 0 || size < 0 || off+size > len(b) {
			return nil
		}
		if s := cffDict(b[off : off+size])[19]; len(s) == 1 {
			subrs, _, err := cffIndex(b, off+int(s[0]))
			if err == nil {
				return subrs
			}
		}
		return nil
	}
	fda, fds := top[1236], top[1237]
	if len(fda) == 1 && len(fds) == 1 {
		fds, _, err := cffIndex(b, int(fda[0]))
		if err != nil {
			return nil, err
		}
		if len(fds) == 0 {
			return nil, errors.New("cff: no Font DICT")
		}
		for _, fd := range fds {
			c.local = append(c.local, private(cffDict(fd)[18]))
		}
		c.fdSelect = cffFDSelect(b, int(top[1237][0]), len(c.charStrings))
	} else {
		c.local = [][][]byte{private(top[18])}
	}
	return c, nil
}

// cffFDSelect reads the Font DICT index of every glyph (formats 0 and 3).
func cffFDSelect(b []byte, pos, n int) []byte {
	out := make([]byte, n)
	if pos < 0 || pos >= len(b) {
		return out
	}
	switch b[pos] {
	case 0:
		if pos+1+n <= len(b) {
			copy(out, b[pos+1:pos+1+n])
		}
	case 3:
		nr := int(be16(b, pos+1))
		for i := 0; i < nr; i++ {
			r := pos + 3 + i*3
			if r+5 > len(b) {
				break
			}
			first, fd, next := int(be16(b, r)), b[r+2], int(be16(b, r+3))
			for g := first; g < next && g < n; g++ {
				out[g] = fd
			}
		}
	}
	return out
}

func subrBias(n int) int {
	switch {
	case n < 1240:
		return 107
	case n < 33900:
		return 1131
	}
	return 32768
}

// t2Sink receives what a Type 2 charstring draws, in font units (y up).
// moveTo starts a contour only once something is drawn from it.
type t2Sink interface {
	moveTo(x, y float64)
	lineTo(x, y float64)
	cubeTo(x1, y1, x2, y2, x3, y3 float64)
	closeContour()
}

// t2Run runs a Type 2 charstring, passing what it draws to a sink.
type t2Run struct {
	g          *cffGlyphs
	local      [][]byte
	sink       t2Sink
	stack      []float64
	x, y       float64
	sx, sy     float64 // the start of the contour, drawn once something follows it
	open       bool
	nStems     int
	widthSeen  bool
	done       bool
	callDepth  int
	operations int
}

// run runs the charstring of a glyph; false when the font has no such glyph.
func (c *cffGlyphs) run(gid int, sink t2Sink) bool {
	if gid < 0 || gid >= len(c.charStrings) {
		return false
	}
	fd := 0
	if c.fdSelect != nil && int(c.fdSelect[gid]) < len(c.local) {
		fd = int(c.fdSelect[gid])
	}
	t := &t2Run{g: c, local: c.local[fd], sink: sink}
	t.run(c.charStrings[gid])
	t.closeContour()
	return true
}

func (c *cffGlyphs) bounds(gid int) (Rect, bool) {
	b := &boundsSink{r: Rect{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}}
	if !c.run(gid, b) {
		return Rect{}, false
	}
	return b.r, b.any
}

// boundsSink collects the bounds of an outline, curve extrema included.
type boundsSink struct {
	r      Rect
	any    bool
	lx, ly float64 // the current point
}

func (b *boundsSink) point(x, y float64) {
	b.r.XMin, b.r.XMax = math.Min(b.r.XMin, x), math.Max(b.r.XMax, x)
	b.r.YMin, b.r.YMax = math.Min(b.r.YMin, y), math.Max(b.r.YMax, y)
	b.any = true
	b.lx, b.ly = x, y
}

func (b *boundsSink) moveTo(x, y float64) { b.point(x, y) }
func (b *boundsSink) lineTo(x, y float64) { b.point(x, y) }
func (b *boundsSink) closeContour()       {}

func (b *boundsSink) cubeTo(x1, y1, x2, y2, x3, y3 float64) {
	x0, y0 := b.lx, b.ly
	for _, s := range cubicExtrema(x0, x1, x2, x3) {
		b.point(cubicAt(x0, x1, x2, x3, s), cubicAt(y0, y1, y2, y3, s))
	}
	for _, s := range cubicExtrema(y0, y1, y2, y3) {
		b.point(cubicAt(x0, x1, x2, x3, s), cubicAt(y0, y1, y2, y3, s))
	}
	b.point(x3, y3)
}

func (t *t2Run) moveTo(dx, dy float64) {
	t.closeContour()
	t.x += dx
	t.y += dy
	t.sx, t.sy = t.x, t.y
}

func (t *t2Run) closeContour() {
	if t.open {
		t.sink.closeContour()
		t.open = false
	}
}

func (t *t2Run) start() {
	if !t.open {
		t.sink.moveTo(t.sx, t.sy)
		t.open = true
	}
}

func (t *t2Run) lineTo(dx, dy float64) {
	t.start()
	t.x += dx
	t.y += dy
	t.sink.lineTo(t.x, t.y)
}

// curveTo adds a cubic Bézier curve from the current point with the given
// relative control points.
func (t *t2Run) curveTo(dx1, dy1, dx2, dy2, dx3, dy3 float64) {
	t.start()
	x1, y1 := t.x+dx1, t.y+dy1
	x2, y2 := x1+dx2, y1+dy2
	x3, y3 := x2+dx3, y2+dy3
	t.sink.cubeTo(x1, y1, x2, y2, x3, y3)
	t.x, t.y = x3, y3
}

func cubicAt(p0, p1, p2, p3, s float64) float64 {
	u := 1 - s
	return u*u*u*p0 + 3*u*u*s*p1 + 3*u*s*s*p2 + s*s*s*p3
}

// cubicExtrema returns the parameters in (0, 1) where a cubic's derivative
// is zero.
func cubicExtrema(p0, p1, p2, p3 float64) []float64 {
	a := -p0 + 3*p1 - 3*p2 + p3
	b := 2 * (p0 - 2*p1 + p2)
	c := p1 - p0
	var out []float64
	add := func(s float64) {
		if s > 0 && s < 1 {
			out = append(out, s)
		}
	}
	if math.Abs(a) < 1e-9 {
		if math.Abs(b) > 1e-9 {
			add(-c / b)
		}
		return out
	}
	d := b*b - 4*a*c
	if d < 0 {
		return out
	}
	sq := math.Sqrt(d)
	add((-b + sq) / (2 * a))
	add((-b - sq) / (2 * a))
	return out
}

// takeWidth drops the advance width that the first stack-clearing
// operator of a charstring may carry: it is there when the operator has
// one operand more than it takes (evenArgs: it takes pairs).
func (t *t2Run) takeWidth(want int, evenArgs bool) {
	if t.widthSeen {
		return
	}
	t.widthSeen = true
	n := len(t.stack)
	if evenArgs && n%2 == 1 || !evenArgs && n > want {
		t.stack = t.stack[1:]
	}
}

func (t *t2Run) run(cs []byte) {
	if t.callDepth > 10 {
		t.done = true
		return
	}
	for i := 0; i < len(cs) && !t.done; {
		t.operations++
		if t.operations > 100000 {
			t.done = true
			return
		}
		v := int(cs[i])
		switch {
		case v == 28:
			if i+2 < len(cs) {
				t.stack = append(t.stack, float64(int16(be16(cs, i+1))))
			}
			i += 3
			continue
		case v >= 32 && v <= 246:
			t.stack = append(t.stack, float64(v-139))
			i++
			continue
		case v >= 247 && v <= 250:
			if i+1 < len(cs) {
				t.stack = append(t.stack, float64((v-247)*256+int(cs[i+1])+108))
			}
			i += 2
			continue
		case v >= 251 && v <= 254:
			if i+1 < len(cs) {
				t.stack = append(t.stack, float64(-(v-251)*256-int(cs[i+1])-108))
			}
			i += 2
			continue
		case v == 255:
			if i+4 < len(cs) {
				t.stack = append(t.stack, float64(int32(be32(cs, i+1)))/65536)
			}
			i += 5
			continue
		}
		i++
		s := t.stack
		switch v {
		case 1, 3, 18, 23: // hstem, vstem, hstemhm, vstemhm
			t.takeWidth(0, true)
			t.nStems += len(t.stack) / 2
		case 19, 20: // hintmask, cntrmask
			t.takeWidth(0, true)
			t.nStems += len(t.stack) / 2
			i += (t.nStems + 7) / 8
		case 21: // rmoveto
			t.takeWidth(2, false)
			if s = t.stack; len(s) >= 2 {
				t.moveTo(s[len(s)-2], s[len(s)-1])
			}
		case 22: // hmoveto
			t.takeWidth(1, false)
			if s = t.stack; len(s) >= 1 {
				t.moveTo(s[len(s)-1], 0)
			}
		case 4: // vmoveto
			t.takeWidth(1, false)
			if s = t.stack; len(s) >= 1 {
				t.moveTo(0, s[len(s)-1])
			}
		case 5: // rlineto
			for k := 0; k+1 < len(s); k += 2 {
				t.lineTo(s[k], s[k+1])
			}
		case 6, 7: // hlineto, vlineto: alternating
			horiz := v == 6
			for _, d := range s {
				if horiz {
					t.lineTo(d, 0)
				} else {
					t.lineTo(0, d)
				}
				horiz = !horiz
			}
		case 8: // rrcurveto
			for k := 0; k+5 < len(s); k += 6 {
				t.curveTo(s[k], s[k+1], s[k+2], s[k+3], s[k+4], s[k+5])
			}
		case 24: // rcurveline
			k := 0
			for ; k+5 < len(s)-2; k += 6 {
				t.curveTo(s[k], s[k+1], s[k+2], s[k+3], s[k+4], s[k+5])
			}
			if k+1 < len(s) {
				t.lineTo(s[k], s[k+1])
			}
		case 25: // rlinecurve
			k := 0
			for ; k+1 < len(s)-6; k += 2 {
				t.lineTo(s[k], s[k+1])
			}
			if k+5 < len(s) {
				t.curveTo(s[k], s[k+1], s[k+2], s[k+3], s[k+4], s[k+5])
			}
		case 26: // vvcurveto
			k := 0
			dx1 := 0.0
			if len(s)%4 == 1 {
				dx1 = s[0]
				k = 1
			}
			for ; k+3 < len(s); k += 4 {
				t.curveTo(dx1, s[k], s[k+1], s[k+2], 0, s[k+3])
				dx1 = 0
			}
		case 27: // hhcurveto
			k := 0
			dy1 := 0.0
			if len(s)%4 == 1 {
				dy1 = s[0]
				k = 1
			}
			for ; k+3 < len(s); k += 4 {
				t.curveTo(s[k], dy1, s[k+1], s[k+2], s[k+3], 0)
				dy1 = 0
			}
		case 30, 31: // vhcurveto, hvcurveto
			horiz := v == 31
			for k := 0; k+3 < len(s); k += 4 {
				last := 0.0
				if k+5 == len(s) {
					last = s[k+4]
				}
				if horiz {
					t.curveTo(s[k], 0, s[k+1], s[k+2], last, s[k+3])
				} else {
					t.curveTo(0, s[k], s[k+1], s[k+2], s[k+3], last)
				}
				horiz = !horiz
			}
		case 10, 29: // callsubr, callgsubr
			if len(s) == 0 {
				t.done = true
				return
			}
			subrs := t.local
			if v == 29 {
				subrs = t.g.global
			}
			n := int(s[len(s)-1]) + subrBias(len(subrs))
			t.stack = s[:len(s)-1]
			if n < 0 || n >= len(subrs) {
				t.done = true
				return
			}
			t.callDepth++
			t.run(subrs[n])
			t.callDepth--
			continue // the subroutine left the stack as it is
		case 11: // return
			return
		case 14: // endchar
			t.takeWidth(0, false)
			t.done = true
			return
		case 12:
			if i >= len(cs) {
				t.done = true
				return
			}
			op := cs[i]
			i++
			t.escape(op, s)
		}
		t.stack = t.stack[:0]
	}
}

// escape runs the two-byte operators: the flexes draw curves; the others
// (arithmetic, storage) are not used by the fonts this reads.
func (t *t2Run) escape(op byte, s []float64) {
	switch op {
	case 35: // flex
		if len(s) >= 12 {
			t.curveTo(s[0], s[1], s[2], s[3], s[4], s[5])
			t.curveTo(s[6], s[7], s[8], s[9], s[10], s[11])
		}
	case 34: // hflex
		if len(s) >= 7 {
			t.curveTo(s[0], 0, s[1], s[2], s[3], 0)
			t.curveTo(s[4], 0, s[5], -s[2], s[6], 0)
		}
	case 36: // hflex1
		if len(s) >= 9 {
			t.curveTo(s[0], s[1], s[2], s[3], s[4], 0)
			t.curveTo(s[5], 0, s[6], s[7], s[8], -(s[1] + s[3] + s[7]))
		}
	case 37: // flex1
		if len(s) >= 11 {
			dx := s[0] + s[2] + s[4] + s[6] + s[8]
			dy := s[1] + s[3] + s[5] + s[7] + s[9]
			x6, y6 := s[10], s[10]
			if math.Abs(dx) > math.Abs(dy) {
				y6 = -dy
			} else {
				x6 = -dx
			}
			t.curveTo(s[0], s[1], s[2], s[3], s[4], s[5])
			t.curveTo(s[6], s[7], s[8], s[9], x6, y6)
		}
	}
}

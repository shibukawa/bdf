package visio

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Fills, lines and shadows of shapes.

// fill kinds
const (
	fillNone = iota
	fillSolid
	fillPattern
	fillGradient
)

type fillStyle struct {
	kind    int
	fg, bg  bdf.Color
	pattern int // 2-24: hatches; 25-40: the gradients of older versions
	stops   []gradStop
	dir     int     // FillGradientDir
	angle   float64 // FillGradientAngle, radians counterclockwise
	rotate  bool    // the gradient turns with the shape
}

func (p *pageCtx) fillStyle(s *shape) fillStyle {
	pat := int(p.num(s, "FillPattern", 1))
	if pat == 0 {
		return fillStyle{}
	}
	fg, _ := p.color(s, "FillForegnd", "FillForegndTrans")
	bg, _ := p.color(s, "FillBkgnd", "FillBkgndTrans")
	f := fillStyle{kind: fillSolid, fg: fg, bg: bg, pattern: pat}
	if p.flag(s, "FillGradientEnabled") {
		for _, k := range s.rowKeys("FillGradient") {
			v, ok := p.rowVal(s, "FillGradient", k, "GradientStopColor")
			if !ok {
				continue
			}
			c, ok := p.c.d.color(v)
			if !ok {
				continue
			}
			c = withTrans(c, p.rowNum(s, "FillGradient", k, "GradientStopColorTrans", 0))
			f.stops = append(f.stops, gradStop{p.rowNum(s, "FillGradient", k, "GradientStopPosition", 0), c})
		}
		if len(f.stops) > 1 {
			f.kind = fillGradient
			f.dir = int(p.num(s, "FillGradientDir", 0))
			f.angle = p.num(s, "FillGradientAngle", 0)
			f.rotate = p.num(s, "RotateGradientWithShape", 1) != 0
			return f
		}
	}
	switch {
	case pat >= 2 && pat <= 24:
		f.kind = fillPattern
	case pat >= 25 && pat <= 40:
		f.kind = fillGradient
		f.stops = legacyGradient(pat, fg, bg)
		f.dir, f.angle = legacyGradientDir(pat)
		f.rotate = true
	}
	return f
}

// legacyGradient returns the stops of fill patterns 25-40, which run from
// the foreground color to the background color.
func legacyGradient(pat int, fg, bg bdf.Color) []gradStop {
	switch pat {
	case 26, 29: // the foreground in the middle
		return []gradStop{{0, bg}, {0.5, fg}, {1, bg}}
	case 35: // from the corners to the center
		return []gradStop{{0, fg}, {1, bg}}
	}
	return []gradStop{{0, fg}, {1, bg}}
}

// legacyGradientDir maps fill patterns 25-40 to a gradient direction and
// angle (see FillGradientDir).
func legacyGradientDir(pat int) (int, float64) {
	switch pat {
	case 25, 26:
		return 0, 0
	case 27:
		return 0, math.Pi
	case 28, 29:
		return 0, -math.Pi / 2
	case 30:
		return 0, math.Pi / 2
	case 31:
		return 0, -math.Pi / 4
	case 32:
		return 0, -3 * math.Pi / 4
	case 33:
		return 0, math.Pi / 4
	case 34:
		return 0, 3 * math.Pi / 4
	case 35:
		return 10, 0
	case 36:
		return 7, 0
	case 37:
		return 6, 0
	case 38:
		return 2, 0
	case 39:
		return 1, 0
	}
	return 3, 0
}

// setFill emits the fill paint for a shape whose local box is w×h and
// whose local coordinates m maps to the page.
func (p *pageCtx) setFill(cv *canvas.Canvas, f fillStyle, m canvas.Matrix, w, h float64) {
	switch f.kind {
	case fillSolid:
		cv.Obj.FillColor(f.fg)
	case fillPattern:
		ref := cv.Image(p.c.patternImage(f.pattern, f.fg, f.bg))
		pp := bdf.Pattern(ref, bdf.RepeatBoth)
		pp.Matrix = [6]float32{0.75, 0, 0, 0.75, 0, 0} // 96 dpi pixels
		cv.Obj.FillPaint(cv.Obj.AddPaint(pp))
	case fillGradient:
		cv.Obj.FillPaint(cv.Obj.AddPaint(gradientPaint(f, m, w, h)))
	}
}

// gradientPaint builds the Canvas gradient of a Visio gradient fill over
// the local box of a shape.
func gradientPaint(f fillStyle, m canvas.Matrix, w, h float64) bdf.Paint {
	stops := make([]bdf.Stop, len(f.stops))
	for i, s := range f.stops {
		stops[i] = bdf.Stop{Offset: float32(math.Max(0, math.Min(1, s.pos))), Color: s.c}
	}
	if f.dir == 0 {
		u := pt{math.Cos(f.angle), math.Sin(f.angle)}
		if !f.rotate {
			// the angle is the page's: undo the shape's rotation
			rot := math.Atan2(m[1], m[0])
			if m[0]*m[3]-m[1]*m[2] < 0 {
				rot = -rot
			}
			u = pt{math.Cos(f.angle + rot), math.Sin(f.angle + rot)}
		}
		t := (math.Abs(u.x)*w + math.Abs(u.y)*h) / 2
		c := pt{w / 2, h / 2}
		p0, p1 := apply(m, c.sub(u.mul(t))), apply(m, c.add(u.mul(t)))
		return bdf.LinearGradient(f32(p0.x), f32(p0.y), f32(p1.x), f32(p1.y), stops...)
	}
	// radial and rectangular gradients spread from a corner, an edge or
	// the center of the box (rectangles are approximated by circles)
	var c pt
	switch f.dir {
	case 1, 8:
		c = pt{w, 0}
	case 2, 9:
		c = pt{0, 0}
	case 4:
		c = pt{w / 2, 0}
	case 5:
		c = pt{w / 2, h}
	case 6, 11:
		c = pt{w, h}
	case 7, 12:
		c = pt{0, h}
	default:
		c = pt{w / 2, h / 2}
	}
	r := 0.0
	for _, q := range []pt{{0, 0}, {w, 0}, {0, h}, {w, h}} {
		r = math.Max(r, q.dist(c))
	}
	if f.dir >= 8 && f.dir <= 12 {
		r = math.Max(math.Max(c.x, w-c.x), math.Max(c.y, h-c.y))
	}
	pc := apply(m, c)
	scale := math.Sqrt(math.Abs(m[0]*m[3] - m[1]*m[2]))
	return bdf.RadialGradient(f32(pc.x), f32(pc.y), 0, f32(pc.x), f32(pc.y), f32(r*scale), stops...)
}

// patternBits are the 8×8 tiles of fill patterns 2-24 (set bit: the
// foreground color), as [MS-VSDX] pictures them.
var patternBits = [...][8]byte{
	2:  {0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80, 0x01},
	3:  {0x10, 0x10, 0x10, 0x10, 0xff, 0x10, 0x10, 0x10},
	4:  {0x04, 0x0a, 0x11, 0xa0, 0x40, 0xa0, 0x11, 0x0a},
	5:  {0x04, 0x02, 0x01, 0x80, 0x40, 0x20, 0x10, 0x08},
	6:  {0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00},
	7:  {0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20},
	8:  {0xbb, 0xee, 0xbb, 0xee, 0xbb, 0xee, 0xbb, 0xee},
	9:  {0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa},
	10: {0x22, 0x88, 0x22, 0x88, 0x22, 0x88, 0x22, 0x88},
	11: {0x22, 0x00, 0x88, 0x00, 0x22, 0x00, 0x88, 0x00},
	12: {0x02, 0x00, 0x20, 0x00, 0x02, 0x00, 0x20, 0x00},
	13: {0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00},
	14: {0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99, 0x99},
	15: {0x66, 0x33, 0x99, 0xcc, 0x66, 0x33, 0x99, 0xcc},
	16: {0x33, 0x66, 0xcc, 0x99, 0x33, 0x66, 0xcc, 0x99},
	17: {0xcc, 0x33, 0x33, 0xcc, 0xcc, 0x33, 0x33, 0xcc},
	18: {0xff, 0x99, 0xff, 0x66, 0xff, 0x99, 0xff, 0x66},
	19: {0x00, 0x00, 0xff, 0x00, 0x00, 0x00, 0xff, 0x00},
	20: {0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
	21: {0x44, 0x22, 0x11, 0x88, 0x44, 0x22, 0x11, 0x88},
	22: {0x44, 0x88, 0x11, 0x22, 0x44, 0x88, 0x11, 0x22},
	23: {0x22, 0xff, 0x22, 0x22, 0x22, 0xff, 0x22, 0x22},
	24: {0x55, 0xbb, 0x55, 0xee, 0x55, 0xbb, 0x55, 0xee},
}

// patternImage renders a fill pattern as an 8×8 PNG tile (cached).
func (c *converter) patternImage(pat int, fg, bg bdf.Color) bdf.Hash {
	type key struct {
		pat    int
		fg, bg bdf.Color
	}
	k := key{pat, fg, bg}
	if h, ok := c.patterns[k]; ok {
		return h
	}
	bits := patternBits[min(max(pat, 2), len(patternBits)-1)]
	nrgba := func(c bdf.Color) color.NRGBA {
		return color.NRGBA{uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), uint8(c)}
	}
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if bits[y]&(0x80>>x) != 0 {
				img.SetNRGBA(x, y, nrgba(fg))
			} else {
				img.SetNRGBA(x, y, nrgba(bg))
			}
		}
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	h := c.doc.AddImage(b.Bytes())
	c.patterns[k] = h
	return h
}

// lineStyle is a shape's resolved outline.
type lineStyle struct {
	width      float64 // points
	color      bdf.Color
	pattern    int
	cap        byte // BDF cap: 0 butt, 1 round, 2 square
	begin, end int  // arrowheads (BeginArrow, EndArrow)
	beginSize  int
	endSize    int
	rounding   float64 // corner radius, local units (inches)
}

func (p *pageCtx) lineStyle(s *shape) *lineStyle {
	pat := int(p.num(s, "LinePattern", 1))
	if pat == 0 {
		return nil
	}
	c, ok := p.color(s, "LineColor", "LineColorTrans")
	if !ok {
		c = 0x000000ff
	}
	l := &lineStyle{width: p.num(s, "LineWeight", 0.01) * 72, color: c, pattern: pat,
		begin: int(p.num(s, "BeginArrow", 0)), end: int(p.num(s, "EndArrow", 0)),
		beginSize: int(p.num(s, "BeginArrowSize", 2)), endSize: int(p.num(s, "EndArrowSize", 2)),
		rounding: p.num(s, "Rounding", 0)}
	switch int(p.num(s, "LineCap", 0)) {
	case 1:
		l.cap = 0
	case 2:
		l.cap = 2
	default:
		l.cap = 1
	}
	if l.width <= 0 {
		l.width = 0.25 // a hairline
	}
	return l
}

// dashes are the dash patterns of lines 2-23 in units of a tenth of a
// dash cycle as [MS-VSDX] pictures them (a dot is 0.25).
var dashes = [...][]float64{
	2:  {2, 1},
	3:  {0.25, 1},
	4:  {2, 1, 0.25, 1},
	5:  {2, 1, 0.25, 1, 0.25, 1},
	6:  {2, 1, 2, 1, 0.25, 1},
	7:  {5, 1, 2, 1},
	8:  {5, 1, 2, 1, 2, 1},
	9:  {1, 0.5},
	10: {0.25, 0.5},
	11: {1, 0.5, 0.25, 0.5},
	12: {1, 0.5, 0.25, 0.5, 0.25, 0.5},
	13: {1, 0.5, 1, 0.5, 0.25, 0.5},
	14: {2.5, 0.5, 1, 0.5},
	15: {2.5, 0.5, 1, 0.5, 1, 0.5},
	16: {4, 2},
	17: {0.25, 2},
	18: {4, 2, 0.25, 2},
	19: {4, 2, 0.25, 2, 0.25, 2},
	20: {4, 2, 4, 2, 0.25, 2},
	21: {10, 2, 4, 2},
	22: {10, 2, 4, 2, 4, 2},
	23: {0.5, 0.25},
}

// setLine emits the stroke style of a line.
func setLine(cv *canvas.Canvas, l *lineStyle) {
	cv.Obj.Line(f32(l.width), l.cap, 0, 10)
	cv.Obj.StrokeColor(l.color)
	if l.pattern >= 2 && l.pattern < len(dashes) {
		// the patterns scale with the line; round and square caps
		// lengthen the dashes by the width
		unit := math.Max(l.width, 1) * 10 / 3
		d := dashes[l.pattern]
		seg := make([]float32, len(d))
		for i, v := range d {
			v *= unit
			if l.cap != 0 {
				if i%2 == 0 {
					v = math.Max(v-l.width, 0.01)
				} else {
					v += l.width
				}
			}
			seg[i] = f32(v)
		}
		cv.Obj.Dash(seg, 0)
	}
}

// shadowStyle is a shape's resolved shadow.
type shadowStyle struct {
	color        bdf.Color
	dx, dy, blur float64 // points, y down
	pattern      int
}

func (p *pageCtx) shadowStyle(s *shape, topLevel, hasGeometry bool) *shadowStyle {
	pat := int(p.num(s, "ShdwPattern", 0))
	if pat == 0 {
		return nil
	}
	switch int(p.num(s, "ShapeShdwShow", 0)) {
	case 0:
		if !hasGeometry {
			return nil
		}
	case 1:
		if !hasGeometry || !topLevel {
			return nil
		}
	}
	c, ok := p.color(s, "ShdwForegnd", "ShdwForegndTrans")
	if !ok {
		return nil
	}
	sh := &shadowStyle{color: c, pattern: pat}
	var dx, dy float64
	if int(p.num(s, "ShapeShdwType", 0)) == 0 {
		// the page's default shadow
		dx, dy = num(p.pg.sheet.cells["ShdwOffsetX"].v), num(p.pg.sheet.cells["ShdwOffsetY"].v)
	} else {
		dx, dy = p.num(s, "ShapeShdwOffsetX", 0), p.num(s, "ShapeShdwOffsetY", 0)
	}
	sh.dx, sh.dy = dx*72, -dy*72
	sh.blur = p.num(s, "ShapeShdwBlur", 0) * 72
	return sh
}

// toPath converts figures into a BDF path (only the closed ones when
// closedOnly).
func toPath(subs []*subpath, closedOnly bool) *bdf.Path {
	p := &bdf.Path{}
	for _, sp := range subs {
		if closedOnly && !sp.closed() {
			continue
		}
		p.MoveTo(f32(sp.start.x), f32(sp.start.y))
		for _, s := range sp.segs {
			if s.cubic {
				p.CubicTo(f32(s.c1.x), f32(s.c1.y), f32(s.c2.x), f32(s.c2.y), f32(s.p.x), f32(s.p.y))
			} else {
				p.LineTo(f32(s.p.x), f32(s.p.y))
			}
		}
		if sp.closed() {
			p.Close()
		}
	}
	return p
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

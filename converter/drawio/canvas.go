package drawio

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// c2d is the drawing interface of draw.io's shapes (mxAbstractCanvas2D as
// mxSvgCanvas2D implements it) on top of a BDF object. Shapes and stencils
// are ported against it nearly line by line.
//
// As in mxAbstractCanvas2D, translate and scale are applied to the
// coordinates of paths here, and rotate becomes a TRANSFORM of the object;
// save and restore map to SAVE and RESTORE, so that the object's state
// stack follows the canvas's. Fill and stroke state is written before each
// fill or stroke when it differs from what the object has in effect.
type c2d struct {
	cv    *canvas.Canvas
	obj   *bdf.Object
	st    c2dState
	stack []c2dState

	path         *bdf.Path // current path, in canvas coordinates (dx, dy and scale applied)
	pathBox      extent
	lastX, lastY float64 // last point before dx, dy and scale (for arcTo)
}

// c2dState is the canvas state (mxAbstractCanvas2D.createState) plus what
// the object has in effect.
type c2dState struct {
	dx, dy, scale                    float64
	alpha, fillAlpha, strokeAlpha    float64
	fillColor                        string // "" = none
	gradientColor                    string // "" = none
	gradientFillAlpha, gradientAlpha float64
	gradientDirection                string
	strokeColor                      string // "" = none
	strokeWidth                      float64
	dashed, fixDash                  bool
	dashPattern                      string
	lineCap, lineJoin                string
	miterLimit                       float64
	rotation                         float64

	// emitted state of the object
	eFill, eStroke bdf.Color
	eLine          [4]float64 // width, cap, join, miter
	eDash          string
	eAlpha         float64
}

func newC2D(cv *canvas.Canvas) *c2d {
	return &c2d{cv: cv, obj: cv.Obj, st: c2dState{
		scale: 1, alpha: 1, fillAlpha: 1, strokeAlpha: 1, gradientFillAlpha: 1, gradientAlpha: 1,
		strokeWidth: 1, dashPattern: "3 3", lineCap: "flat", lineJoin: "miter", miterLimit: 10,
		eFill: bdf.RGB(0, 0, 0), eStroke: bdf.RGB(0, 0, 0), eLine: [4]float64{1, 0, 0, 10}, eAlpha: 1,
	}}
}

// --- state ---

func (c *c2d) save() {
	c.stack = append(c.stack, c.st)
	c.obj.Save()
}

func (c *c2d) restore() {
	if len(c.stack) == 0 {
		return
	}
	c.st = c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	c.obj.Restore()
}

func (c *c2d) translate(dx, dy float64) {
	c.st.dx += dx
	c.st.dy += dy
}

// rotate turns what follows by theta degrees about (cx, cy), after
// flipping (mxSvgCanvas2D.rotate).
func (c *c2d) rotate(theta float64, flipH, flipV bool, cx, cy float64) {
	if theta == 0 && !flipH && !flipV {
		return
	}
	s := &c.st
	cx = (cx + s.dx) * s.scale
	cy = (cy + s.dy) * s.scale
	m := identity
	if flipH && flipV {
		theta += 180
	} else if flipH != flipV {
		tx, sx, ty, sy := 0.0, 1.0, 0.0, 1.0
		if flipH {
			tx, sx = cx, -1
		}
		if flipV {
			ty, sy = cy, -1
		}
		m = translateM(tx, ty).mul(scaleM(sx, sy)).mul(translateM(-tx, -ty))
	}
	if flipH != flipV {
		theta = -theta
	}
	if theta != 0 {
		m = m.mul(rotateAbout(theta, cx, cy))
	}
	s.rotation += theta
	c.cv.Transform(canvas.Matrix(m))
}

func (c *c2d) setAlpha(v float64)       { c.st.alpha = v }
func (c *c2d) setFillAlpha(v float64)   { c.st.fillAlpha = v }
func (c *c2d) setStrokeAlpha(v float64) { c.st.strokeAlpha = v }

func (c *c2d) setFillColor(v string) {
	if v == "none" {
		v = ""
	}
	c.st.fillColor = v
	c.st.gradientColor = ""
}

// setGradient fills with a two-color gradient over the bounding box of
// what is filled; direction is south (default), north, east, west or radial.
func (c *c2d) setGradient(color1, color2 string, x, y, w, h float64, direction string, alpha1, alpha2 float64) {
	c.st.fillColor = color1
	c.st.gradientFillAlpha = alpha1
	c.st.gradientColor = color2
	c.st.gradientAlpha = alpha2
	c.st.gradientDirection = direction
}

func (c *c2d) setStrokeColor(v string) {
	if v == "none" {
		v = ""
	}
	c.st.strokeColor = v
}

func (c *c2d) setStrokeWidth(v float64)  { c.st.strokeWidth = v }
func (c *c2d) setDashed(v, fixDash bool) { c.st.dashed, c.st.fixDash = v, fixDash }
func (c *c2d) setDashPattern(v string)   { c.st.dashPattern = v }
func (c *c2d) setLineCap(v string)       { c.st.lineCap = v }
func (c *c2d) setLineJoin(v string)      { c.st.lineJoin = v }
func (c *c2d) setMiterLimit(v float64)   { c.st.miterLimit = v }

// setShadow is part of the interface; shadows are drawn for a whole shape
// (see paintShape), as draw.io's drop-shadow filter does.
func (c *c2d) setShadow(bool) {}

// --- paths ---

func (c *c2d) begin() {
	c.path = &bdf.Path{}
	c.pathBox = extent{}
	c.lastX, c.lastY = 0, 0
}

func (c *c2d) tx(x float64) float64 { return (x + c.st.dx) * c.st.scale }
func (c *c2d) ty(y float64) float64 { return (y + c.st.dy) * c.st.scale }

func (c *c2d) addPt(x, y float64) (float32, float32) {
	c.lastX, c.lastY = x, y
	px, py := c.tx(x), c.ty(y)
	c.pathBox.add(point{px, py})
	return f32(px), f32(py)
}

func (c *c2d) ensurePath() {
	if c.path == nil {
		c.begin()
	}
}

func (c *c2d) moveTo(x, y float64) {
	c.ensurePath()
	c.path.MoveTo(c.addPt(x, y))
}

func (c *c2d) lineTo(x, y float64) {
	c.ensurePath()
	c.path.LineTo(c.addPt(x, y))
}

func (c *c2d) quadTo(x1, y1, x2, y2 float64) {
	c.ensurePath()
	ax, ay := c.addPt(x1, y1)
	bx, by := c.addPt(x2, y2)
	c.path.QuadTo(ax, ay, bx, by)
}

func (c *c2d) curveTo(x1, y1, x2, y2, x3, y3 float64) {
	c.ensurePath()
	ax, ay := c.addPt(x1, y1)
	bx, by := c.addPt(x2, y2)
	ex, ey := c.addPt(x3, y3)
	c.path.CubicTo(ax, ay, bx, by, ex, ey)
}

// arcTo adds an SVG elliptical arc from the last point (mxAbstractCanvas2D.arcTo).
func (c *c2d) arcTo(rx, ry, angle float64, largeArc, sweep bool, x, y float64) {
	curves := arcToCurves(c.lastX, c.lastY, rx, ry, angle, largeArc, sweep, x, y)
	for i := 0; i+5 < len(curves); i += 6 {
		c.curveTo(curves[i], curves[i+1], curves[i+2], curves[i+3], curves[i+4], curves[i+5])
	}
}

func (c *c2d) close() {
	if c.path != nil {
		c.path.Close()
	}
}

func (c *c2d) end() {}

// rect, roundrect and ellipse replace the current path with a shape.
func (c *c2d) rect(x, y, w, h float64) {
	c.begin()
	s := c.st.scale
	px, py := c.tx(x), c.ty(y)
	c.path.Rect(f32(px), f32(py), f32(w*s), f32(h*s))
	c.pathBox = extent{rect{px, py, w * s, h * s}, true}
}

func (c *c2d) roundrect(x, y, w, h, dx, dy float64) {
	if dx <= 0 && dy <= 0 {
		c.rect(x, y, w, h)
		return
	}
	s := c.st.scale
	if dx == dy || dx <= 0 || dy <= 0 {
		r := math.Max(dx, dy)
		r = math.Min(r, math.Min(w, h)/2)
		c.begin()
		px, py := c.tx(x), c.ty(y)
		c.path.RoundRect(f32(px), f32(py), f32(w*s), f32(h*s), f32(r*s))
		c.pathBox = extent{rect{px, py, w * s, h * s}, true}
		return
	}
	// elliptical corners, as SVG's rx and ry (clamped to half the sides)
	dx, dy = math.Min(dx, w/2), math.Min(dy, h/2)
	const k = 0.5522847498
	c.begin()
	c.moveTo(x+dx, y)
	c.lineTo(x+w-dx, y)
	c.curveTo(x+w-dx+dx*k, y, x+w, y+dy-dy*k, x+w, y+dy)
	c.lineTo(x+w, y+h-dy)
	c.curveTo(x+w, y+h-dy+dy*k, x+w-dx+dx*k, y+h, x+w-dx, y+h)
	c.lineTo(x+dx, y+h)
	c.curveTo(x+dx-dx*k, y+h, x, y+h-dy+dy*k, x, y+h-dy)
	c.lineTo(x, y+dy)
	c.curveTo(x, y+dy-dy*k, x+dx-dx*k, y, x+dx, y)
	c.close()
}

func (c *c2d) ellipse(x, y, w, h float64) {
	c.begin()
	s := c.st.scale
	cx, cy := c.tx(x+w/2), c.ty(y+h/2)
	c.path.Ellipse(f32(cx), f32(cy), f32(math.Abs(w/2*s)), f32(math.Abs(h/2*s)), 0, 0, 2*math.Pi, false)
	c.pathBox = extent{rect{cx - math.Abs(w/2*s), cy - math.Abs(h/2*s), math.Abs(w * s), math.Abs(h * s)}, true}
}

// --- painting ---

func (c *c2d) fill()          { c.paint(true, false) }
func (c *c2d) stroke()        { c.paint(false, true) }
func (c *c2d) fillAndStroke() { c.paint(true, true) }

func (c *c2d) paint(fill, stroke bool) {
	if c.path == nil || len(c.path.Verbs) == 0 {
		c.path = nil
		return
	}
	path := c.path
	box := c.pathBox.r
	c.path = nil
	s := &c.st
	doFill := fill && s.fillColor != ""
	doStroke := stroke && s.strokeColor != ""
	if !doFill && !doStroke {
		return
	}
	ref := c.obj.AddPath(path)
	c.cv.Drawn = true
	if doFill {
		c.fillPath(ref, path, box)
	}
	if doStroke {
		col, ok := parseColor(s.strokeColor)
		if ok {
			col = col.withAlpha(s.alpha * s.strokeAlpha)
			if col.a > 0 {
				c.setStrokeState(col)
				c.obj.StrokePath(ref)
			}
		}
	}
}

func (c *c2d) fillPath(ref bdf.PathRef, path *bdf.Path, box rect) {
	s := &c.st
	alpha := s.alpha * s.fillAlpha
	c1, ok1 := parseColor(s.fillColor)
	if s.gradientColor != "" {
		c2, ok2 := parseColor(s.gradientColor)
		if ok1 || ok2 {
			if !ok1 {
				c1 = rgba{c2.r, c2.g, c2.b, 0}
			}
			if !ok2 {
				c2 = rgba{c1.r, c1.g, c1.b, 0}
			}
			c1 = c1.withAlpha(alpha * s.gradientFillAlpha)
			c2 = c2.withAlpha(alpha * s.gradientAlpha)
			c.fillGradient(ref, path, box, c1, c2)
		}
		return
	}
	if !ok1 {
		return
	}
	col := c1.withAlpha(alpha)
	if col.a == 0 {
		return
	}
	if col.bdf() != s.eFill {
		c.obj.FillColor(col.bdf())
		s.eFill = col.bdf()
	}
	c.obj.FillPath(ref, bdf.NonZero)
}

// fillGradient fills with a gradient over the path's bounding box, as SVG's
// objectBoundingBox gradients (mxSvgCanvas2D.createSvgGradient).
func (c *c2d) fillGradient(ref bdf.PathRef, path *bdf.Path, box rect, c1, c2 rgba) {
	dir := c.st.gradientDirection
	stops := []bdf.Stop{{Offset: 0, Color: c1.bdf()}, {Offset: 1, Color: c2.bdf()}}
	if dir == "radial" {
		if box.w <= 0 || box.h <= 0 {
			return
		}
		// A circle in the unit square of the bounding box: fill the path
		// mapped into that square under a transform back to the box.
		unit := mapPath(path, func(x, y float32) (float32, float32) {
			return f32((float64(x) - box.x) / box.w), f32((float64(y) - box.y) / box.h)
		})
		uref := c.obj.AddPath(unit)
		p := c.obj.AddPaint(bdf.RadialGradient(0.5, 0.5, 0, 0.5, 0.5, 0.5, stops...))
		c.obj.Save()
		c.obj.Transform(f32(box.w), 0, 0, f32(box.h), f32(box.x), f32(box.y))
		c.obj.FillPaint(p)
		c.obj.FillPath(uref, bdf.NonZero)
		c.obj.Restore()
		return
	}
	x0, y0, x1, y1 := box.x, box.y, box.x, box.y+box.h // south
	switch dir {
	case "east":
		x1, y1 = box.x+box.w, box.y
	case "north":
		y0, y1 = box.y+box.h, box.y
	case "west":
		x0, x1, y1 = box.x+box.w, box.x, box.y
	}
	p := c.obj.AddPaint(bdf.LinearGradient(f32(x0), f32(y0), f32(x1), f32(y1), stops...))
	c.obj.FillPaint(p)
	c.st.eFill = 0 // a paint is in effect
	c.obj.FillPath(ref, bdf.NonZero)
}

// mapPath returns a copy of a path with f applied to its coordinates
// (paths made here hold moves, lines, curves, rectangles and ellipses
// without rotation; rectangles and ellipses are mapped by their corners).
func mapPath(p *bdf.Path, f func(x, y float32) (float32, float32)) *bdf.Path {
	out := &bdf.Path{Verbs: append([]byte(nil), p.Verbs...), Args: make([]float32, len(p.Args))}
	ai := 0
	for _, v := range p.Verbs {
		switch v {
		case bdf.VerbMove, bdf.VerbLine:
			out.Args[ai], out.Args[ai+1] = f(p.Args[ai], p.Args[ai+1])
			ai += 2
		case bdf.VerbQuad:
			for k := 0; k < 4; k += 2 {
				out.Args[ai+k], out.Args[ai+k+1] = f(p.Args[ai+k], p.Args[ai+k+1])
			}
			ai += 4
		case bdf.VerbCubic:
			for k := 0; k < 6; k += 2 {
				out.Args[ai+k], out.Args[ai+k+1] = f(p.Args[ai+k], p.Args[ai+k+1])
			}
			ai += 6
		case bdf.VerbRect, bdf.VerbRoundRect:
			x0, y0 := f(p.Args[ai], p.Args[ai+1])
			x1, y1 := f(p.Args[ai]+p.Args[ai+2], p.Args[ai+1]+p.Args[ai+3])
			out.Args[ai], out.Args[ai+1], out.Args[ai+2], out.Args[ai+3] = x0, y0, x1-x0, y1-y0
			if v == bdf.VerbRoundRect {
				// radius scaled by the mean factor (only used for radial gradients)
				sx, sy := float64(x1-x0)/float64(p.Args[ai+2]), float64(y1-y0)/float64(p.Args[ai+3])
				out.Args[ai+4] = f32(float64(p.Args[ai+4]) * (sx + sy) / 2)
				ai += 5
			} else {
				ai += 4
			}
		case bdf.VerbEllipse:
			cx, cy := f(p.Args[ai], p.Args[ai+1])
			ex, ey := f(p.Args[ai]+p.Args[ai+2], p.Args[ai+1]+p.Args[ai+3])
			copy(out.Args[ai:ai+8], p.Args[ai:ai+8])
			out.Args[ai], out.Args[ai+1], out.Args[ai+2], out.Args[ai+3] = cx, cy, ex-cx, ey-cy
			ai += 8
		case bdf.VerbArcTo:
			out.Args[ai], out.Args[ai+1] = f(p.Args[ai], p.Args[ai+1])
			out.Args[ai+2], out.Args[ai+3] = f(p.Args[ai+2], p.Args[ai+3])
			out.Args[ai+4] = p.Args[ai+4]
			ai += 5
		case bdf.VerbClose:
		}
	}
	return out
}

// setStrokeState writes the stroke color, line and dash state
// (mxSvgCanvas2D.updateStroke).
func (c *c2d) setStrokeState(col rgba) {
	s := &c.st
	if col.bdf() != s.eStroke {
		c.obj.StrokeColor(col.bdf())
		s.eStroke = col.bdf()
	}
	w := math.Max(1, math.Max(0.01, s.strokeWidth*s.scale)) // minStrokeWidth
	var capV, joinV float64
	switch s.lineCap {
	case "round":
		capV = 1
	case "square":
		capV = 2
	}
	switch s.lineJoin {
	case "round":
		joinV = 1
	case "bevel":
		joinV = 2
	}
	line := [4]float64{w, capV, joinV, s.miterLimit}
	if line != s.eLine {
		c.obj.Line(f32(w), byte(capV), byte(joinV), f32(s.miterLimit))
		s.eLine = line
	}
	dash := ""
	if s.dashed {
		f := s.strokeWidth
		if s.fixDash {
			f = 1
		}
		dash = dashKey(s.dashPattern, f*s.scale)
	}
	if dash != s.eDash {
		var segs []float32
		if dash != "" {
			for _, t := range strings.Fields(dash) {
				v, _ := parseFloat(t)
				segs = append(segs, f32(v))
			}
		}
		c.obj.Dash(segs, 0)
		s.eDash = dash
	}
}

// dashKey scales a dash pattern ("3 3") and returns it as a string, the
// empty string for no dashes.
func dashKey(pattern string, scale float64) string {
	var b strings.Builder
	all0 := true
	for i, t := range strings.Fields(strings.ReplaceAll(pattern, ",", " ")) {
		v, ok := parseFloat(t)
		if !ok || v < 0 {
			return ""
		}
		if v != 0 {
			all0 = false
		}
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmtNum(math.Round(v*scale*100) / 100))
	}
	if all0 {
		return ""
	}
	return b.String()
}

// image draws an image part in the box x, y, w, h; aspect keeps its
// proportions (centered), flips mirror it.
func (c *c2d) image(x, y, w, h float64, img *imageRef, aspect, flipH, flipV bool) {
	if img == nil {
		return
	}
	s := c.st.scale
	x, y, w, h = c.tx(x), c.ty(y), w*s, h*s
	if aspect && img.w > 0 && img.h > 0 {
		sc := math.Min(w/img.w, h/img.h)
		nw, nh := img.w*sc, img.h*sc
		x += (w - nw) / 2
		y += (h - nh) / 2
		w, h = nw, nh
	}
	ref := c.cv.Image(img.hash)
	c.cv.Drawn = true
	alpha := c.st.alpha
	needSave := flipH || flipV || alpha < 1
	if needSave {
		c.obj.Save()
		if alpha < 1 {
			c.obj.Alpha(f32(alpha))
		}
		if flipH || flipV {
			sx, sy := 1.0, 1.0
			tx, ty := 0.0, 0.0
			if flipH {
				sx, tx = -1, 2*x+w
			}
			if flipV {
				sy, ty = -1, 2*y+h
			}
			c.cv.Transform(canvas.Matrix{sx, 0, 0, sy, tx, ty})
		}
	}
	c.obj.Image(ref, f32(x), f32(y), f32(w), f32(h))
	if needSave {
		c.obj.Restore()
	}
}

// arcToCurves converts an SVG arc to cubic Bézier curves (mxUtils.arcToCurves).
func arcToCurves(x0, y0, r1, r2, angle float64, largeArcFlag, sweepFlag bool, x, y float64) []float64 {
	x -= x0
	y -= y0
	if r1 == 0 || r2 == 0 {
		return nil
	}
	fS := 0.0
	if sweepFlag {
		fS = 1
	}
	fA := false
	if largeArcFlag {
		fA = true
	}
	psai := angle
	r1 = math.Abs(r1)
	r2 = math.Abs(r2)
	ctx := -x / 2
	cty := -y / 2
	cpsi := math.Cos(psai * math.Pi / 180)
	spsi := math.Sin(psai * math.Pi / 180)
	rxd := cpsi*ctx + spsi*cty
	ryd := -1*spsi*ctx + cpsi*cty
	rxdd := rxd * rxd
	rydd := ryd * ryd
	r1x := r1 * r1
	r2y := r2 * r2
	lamda := rxdd/r1x + rydd/r2y
	var sds float64
	if lamda > 1 {
		r1 = math.Sqrt(lamda) * r1
		r2 = math.Sqrt(lamda) * r2
		sds = 0
	} else {
		seif := 1.0
		if fA == (fS == 1) {
			seif = -1
		}
		// clamped for radii that only just span the chord (semicircles)
		sds = seif * math.Sqrt(math.Max(0, (r1x*r2y-r1x*rydd-r2y*rxdd)/(r1x*rydd+r2y*rxdd)))
	}
	txd := sds * r1 * ryd / r2
	tyd := -1 * sds * r2 * rxd / r1
	tx := cpsi*txd - spsi*tyd + x/2
	ty := spsi*txd + cpsi*tyd + y/2
	rad := math.Atan2((ryd-tyd)/r2, (rxd-txd)/r1) - math.Atan2(0, 1)
	s1 := rad
	if rad < 0 {
		s1 = 2*math.Pi + rad
	}
	rad = math.Atan2((-ryd-tyd)/r2, (-rxd-txd)/r1) - math.Atan2((ryd-tyd)/r2, (rxd-txd)/r1)
	dr := rad
	if rad < 0 {
		dr = 2*math.Pi + rad
	}
	if fS == 0 && dr > 0 {
		dr -= 2 * math.Pi
	} else if fS != 0 && dr < 0 {
		dr += 2 * math.Pi
	}
	sse := dr * 2 / math.Pi
	seg := int(math.Ceil(sse))
	if sse < 0 {
		seg = int(math.Ceil(-1 * sse))
	}
	segr := dr / float64(seg)
	t := 8 / 3.0 * math.Sin(segr/4) * math.Sin(segr/4) / math.Sin(segr/2)
	cpsir1 := cpsi * r1
	cpsir2 := cpsi * r2
	spsir1 := spsi * r1
	spsir2 := spsi * r2
	mc := math.Cos(s1)
	ms := math.Sin(s1)
	x2 := -t * (cpsir1*ms + spsir2*mc)
	y2 := -t * (spsir1*ms - cpsir2*mc)
	x3, y3 := 0.0, 0.0
	var result []float64
	for n := 0; n < seg; n++ {
		s1 += segr
		mc = math.Cos(s1)
		ms = math.Sin(s1)
		x3 = cpsir1*mc - spsir2*ms + tx
		y3 = spsir1*mc + cpsir2*ms + ty
		dx := -t * (cpsir1*ms + spsir2*mc)
		dy := -t * (spsir1*ms - cpsir2*mc)
		result = append(result, x2+x0, y2+y0, x3-dx+x0, y3-dy+y0, x3+x0, y3+y0)
		x2 = x3 + dx
		y2 = y3 + dy
	}
	return result
}

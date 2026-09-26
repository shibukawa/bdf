package drawio

import (
	"math"

	"github.com/shibukawa/bdf"
)

// Shapes follow mxShape: a shape reads its colors and flags from the style
// (mxShape.apply), sets up the canvas (configureCanvas), turns it for
// rotation, flips and direction (updateTransform), and paints the vertex
// or the edge. Each kind of shape is a shapeDef whose functions stand for
// the methods the JavaScript classes override; nil functions fall back to
// mxShape's behavior.

// shapeDef describes a kind of shape.
type shapeDef struct {
	// paintVertex paints a vertex (mxShape.paintVertexShape); nil paints
	// the background and then the foreground.
	paintVertex func(s *shape, c *c2d, x, y, w, h float64)
	// paintBackground and paintForeground are the two halves of the
	// default paintVertex.
	paintBackground func(s *shape, c *c2d, x, y, w, h float64)
	paintForeground func(s *shape, c *c2d, x, y, w, h float64)
	// paintEdge paints an edge through its points (mxShape.paintEdgeShape).
	paintEdge func(s *shape, c *c2d, pts []point)
	// labelBounds adjusts the label rectangle (mxShape.getLabelBounds,
	// after getLabelMargins); nil keeps it.
	labelBounds func(s *shape, r rect) rect
	// labelMargins returns the insets of the label rectangle
	// (mxShape.getLabelMargins, as left, top, right, bottom); nil for none.
	labelMargins func(s *shape, r rect) *rect
	// noInvert disables the swapped paint bounds of north and south
	// directions (isPaintBoundsInverted returns false).
	noInvert bool
	// noRotation paints without the style's rotation (edges, mxPolyline).
	noRotation bool
	// roundable marks shapes the rounded style applies to.
	roundable bool
}

var shapeRegistry = map[string]*shapeDef{}

func registerShape(name string, d *shapeDef) { shapeRegistry[name] = d }

// shape is a shape being painted for a cell state (mxShape after apply).
type shape struct {
	conv  *converter
	st    *cellState
	style style
	def   *shapeDef
	name  string

	// stencil is set for shapes drawn from a stencil library.
	stencil *stencil

	bounds rect
	edge   bool
	points []point // the painted points of an edge

	fill, gradient, stroke        string // "" = none
	gradientDirection             string
	opacity, fillOpacity          float64 // 0-100
	strokeOpacity                 float64
	strokewidth                   float64
	startSize, endSize            float64
	startArrow, endArrow          string
	rotation                      float64
	direction                     string
	flipH, flipV                  bool
	isShadow, isDashed, isRounded bool
	glass                         bool
	outline                       bool
}

// newShape creates the shape of a state (mxCellRenderer.createShape):
// a registered shape, a stencil, or the default rectangle or connector.
func (c *converter) newShape(st *cellState) *shape {
	s := &shape{conv: c, st: st, style: st.style}
	s.name = st.style.get("shape", "")
	if d, ok := shapeRegistry[s.name]; ok {
		s.def = d
	} else if sten := c.stencil(s.name); sten != nil {
		s.stencil = sten
		s.def = &shapeDef{}
	} else {
		if st.cell.edge {
			s.def = shapeRegistry["connector"]
		} else {
			s.def = shapeRegistry["rectangle"]
		}
		if s.name != "" {
			c.warnOnce("shape:"+s.name, "shape %q is not supported; drawn as a %s", s.name, map[bool]string{true: "line", false: "rectangle"}[st.cell.edge])
		}
	}
	s.apply()
	s.edge = st.cell.edge
	s.bounds = st.bounds()
	if s.edge {
		s.points = getWaypoints(st.edgePoints())
	}
	return s
}

// getWaypoints returns the points of a route that are painted: those at
// least 1 away from their predecessor in the route (mxShape.getWaypoints),
// or nil when fewer than two remain and the edge is not painted.
func getWaypoints(pts []point) []point {
	if len(pts) == 0 {
		return nil
	}
	out := []point{pts[0]}
	p0 := pts[0]
	for _, pe := range pts[1:] {
		if math.Abs(p0.x-pe.x) >= 1 || math.Abs(p0.y-pe.y) >= 1 {
			out = append(out, pe)
		}
		p0 = pe
	}
	if len(out) < 2 {
		return nil
	}
	return out
}

// apply reads the style (mxShape.apply).
func (s *shape) apply() {
	st := s.style
	s.fill = colorOrNone(st.get("fillColor", ""))
	s.gradient = colorOrNone(st.get("gradientColor", ""))
	s.gradientDirection = st.get("gradientDirection", "")
	s.opacity = st.num("opacity", 100)
	s.fillOpacity = st.num("fillOpacity", 100)
	s.strokeOpacity = st.num("strokeOpacity", 100)
	s.stroke = colorOrNone(st.get("strokeColor", ""))
	s.strokewidth = st.num("strokeWidth", 1)
	s.startSize = st.num("startSize", 0)
	s.endSize = st.num("endSize", 0)
	s.startArrow = st.get("startArrow", "")
	s.endArrow = st.get("endArrow", "")
	s.rotation = st.num("rotation", 0)
	s.direction = st.get("direction", "")
	s.flipH = st.is("flipH")
	s.flipV = st.is("flipV")
	if s.stencil != nil {
		s.flipH = st.is("stencilFlipH") || s.flipH
		s.flipV = st.is("stencilFlipV") || s.flipV
	}
	if s.direction == "north" || s.direction == "south" {
		s.flipH, s.flipV = s.flipV, s.flipH
	}
	s.isShadow = st.is("shadow") || s.conv.pageShadow
	s.isDashed = st.is("dashed")
	s.isRounded = st.is("rounded")
	s.glass = st.is("glass")
}

func colorOrNone(v string) string {
	if v == "none" {
		return ""
	}
	return v
}

// isPaintBoundsInverted reports whether the paint rectangle is turned for
// a north or south direction (stencils handle directions themselves).
func (s *shape) isPaintBoundsInverted() bool {
	return s.stencil == nil && !s.def.noInvert && (s.direction == "north" || s.direction == "south")
}

// shapeRotation is the rotation plus the direction's turn (mxShape.getShapeRotation).
func (s *shape) shapeRotation() float64 {
	if s.def.noRotation {
		return 0
	}
	rot := s.rotation
	switch s.direction {
	case "north":
		rot += 270
	case "west":
		rot += 180
	case "south":
		rot += 90
	}
	return rot
}

// arcSize returns the corner radius of rounded shapes (mxShape.getArcSize).
func (s *shape) arcSize(w, h float64) float64 {
	if s.style.get("absoluteArcSize", "0") == "1" {
		return math.Min(w/2, math.Min(h/2, s.style.num("arcSize", lineArcSize)/2))
	}
	f := s.style.num("arcSize", rectangleRoundingFactor*100) / 100
	return math.Min(w*f, h*f)
}

const (
	rectangleRoundingFactor = 0.15
	lineArcSize             = 20
	defaultMarkerSize       = 6
)

// paint draws the shape (mxShape.paint), with draw.io's drop shadow
// around the whole of it.
func (s *shape) paint(c *c2d) {
	x, y, w, h := s.bounds.x, s.bounds.y, s.bounds.w, s.bounds.h
	edge := s.edge
	if edge && s.points == nil {
		return
	}
	if !edge && (w <= 0 && h <= 0) && s.stencil == nil {
		// mxShape.checkBounds: nothing to draw
		if s.def.paintVertex == nil && s.def.paintBackground == nil {
			return
		}
	}
	shadow := s.isShadow && s.hasPaint()
	if shadow {
		s.beginShadow(c)
	}
	c.save()
	if s.isPaintBoundsInverted() {
		t := (w - h) / 2
		x += t
		y -= t
		w, h = h, w
	}
	// updateTransform
	c.rotate(s.shapeRotation(), s.flipH && !edge, s.flipV && !edge, x+w/2, y+h/2)
	s.configureCanvas(c, x, y, w, h)
	if s.stencil != nil {
		s.stencil.drawShape(c, s, x, y, w, h)
	} else {
		c.setStrokeWidth(s.strokewidth)
		if edge {
			if len(s.points) > 1 && s.def.paintEdge != nil {
				s.def.paintEdge(s, c, s.points)
			}
		} else {
			dx, dy := c.st.dx, c.st.dy
			s.paintVertexShape(c, x, y, w, h)
			c.st.dx, c.st.dy = dx, dy
		}
	}
	c.restore()
	if shadow {
		s.endShadow(c)
	}
}

// hasPaint reports whether the shape fills or strokes anything (a shadow
// is only drawn then).
func (s *shape) hasPaint() bool {
	return s.fill != "" || s.stroke != "" || s.stencil != nil || s.name == "image" || s.name == "label"
}

// paintVertexShape paints a vertex (mxShape.paintVertexShape).
func (s *shape) paintVertexShape(c *c2d, x, y, w, h float64) {
	if s.def.paintVertex != nil {
		s.def.paintVertex(s, c, x, y, w, h)
		return
	}
	s.paintBackground(c, x, y, w, h)
	s.paintForeground(c, x, y, w, h)
}

func (s *shape) paintBackground(c *c2d, x, y, w, h float64) {
	if s.def.paintBackground != nil {
		s.def.paintBackground(s, c, x, y, w, h)
	}
}

func (s *shape) paintForeground(c *c2d, x, y, w, h float64) {
	if s.def.paintForeground != nil {
		s.def.paintForeground(s, c, x, y, w, h)
	}
}

// configureCanvas sets up colors and line styles (mxShape.configureCanvas).
func (s *shape) configureCanvas(c *c2d, x, y, w, h float64) {
	c.setAlpha(s.opacity / 100)
	c.setFillAlpha(s.fillOpacity / 100)
	c.setStrokeAlpha(s.strokeOpacity / 100)
	c.setDashed(s.isDashed, s.style.is("fixDash"))
	if dp, ok := s.style["dashPattern"]; ok {
		c.setDashPattern(dp)
	}
	if s.fill != "" && s.gradient != "" {
		c.setGradient(s.fill, s.gradient, x, y, w, h, s.gradientDirection, 1, 1)
	} else {
		c.setFillColor(s.fill)
	}
	if v, ok := s.style["linecap"]; ok {
		c.setLineCap(v)
	}
	if v, ok := s.style["linejoin"]; ok {
		c.setLineJoin(v)
	}
	if v, ok := s.style["miterlimit"]; ok {
		if f, ok := parseFloat(v); ok {
			c.setMiterLimit(f)
		}
	}
	c.setStrokeColor(s.stroke)
}

// Shadows: draw.io puts a CSS drop-shadow filter on the shape's element
// (mxShape.createDropShadow). Here the shape is drawn into a group that is
// composited with a canvas shadow, which shadows the union of its fills
// and strokes the same way.

func (s *shape) shadowStyle() (dx, dy, blur float64, col rgba) {
	dx = s.style.num("shadowOffsetX", 2)
	dy = s.style.num("shadowOffsetY", 3)
	blur = s.style.num("shadowBlur", 2)
	col, ok := parseColor(s.style.get("shadowColor", "#000000"))
	if !ok {
		col = rgba{0, 0, 0, 255}
	}
	col = col.withAlpha(s.style.num("shadowOpacity", 25) / 100)
	return
}

func (s *shape) beginShadow(c *c2d) {
	dx, dy, blur, col := s.shadowStyle()
	r := s.paintExtent()
	c.obj.Save()
	c.obj.Shadow(col.bdf(), f32(blur), f32(dx), f32(dy))
	c.obj.GroupBegin(1, bdf.BlendSourceOver, f32(r.x), f32(r.y), f32(r.w), f32(r.h))
}

func (s *shape) endShadow(c *c2d) {
	c.obj.GroupEnd()
	c.obj.Restore()
}

// paintExtent is a rectangle that holds everything the shape paints:
// the (rotated) bounds or the points, grown by the stroke and markers.
func (s *shape) paintExtent() rect {
	var r rect
	if s.edge {
		r, _ = pointsExtent(s.points)
		m := math.Max(s.startSize, s.endSize)
		if m == 0 {
			m = defaultMarkerSize
		}
		return r.grow(2*m + 2*s.strokewidth + 20)
	}
	r = s.bounds
	if rot := s.shapeRotation(); rot != 0 {
		// the circle around the rectangle covers any rotation
		d := math.Hypot(r.w, r.h) / 2
		r = rect{r.cx() - d, r.cy() - d, 2 * d, 2 * d}
	}
	return r.grow(2*s.strokewidth + 10)
}

// addPoints adds a polyline with optional rounded corners to the path
// (mxShape.addPoints).
func (s *shape) addPoints(c *c2d, pts []point, rounded bool, arcSize float64, close bool, exclude []int, initialMove bool) {
	if len(pts) == 0 {
		return
	}
	pe := pts[len(pts)-1]
	if close && rounded {
		p0 := pts[0]
		wp := point{pe.x + (p0.x-pe.x)/2, pe.y + (p0.y-pe.y)/2}
		pts = append([]point{wp}, pts...)
	}
	pt := pts[0]
	i := 1
	if initialMove {
		c.moveTo(pt.x, pt.y)
	} else {
		c.lineTo(pt.x, pt.y)
	}
	limit := len(pts) - 1
	if close {
		limit = len(pts)
	}
	excluded := func(k int) bool {
		for _, e := range exclude {
			if e == k {
				return true
			}
		}
		return false
	}
	for i < limit {
		tmp := pts[mod(i, len(pts))]
		dx, dy := pt.x-tmp.x, pt.y-tmp.y
		if rounded && (dx != 0 || dy != 0) && !excluded(i-1) {
			dist := math.Sqrt(dx*dx + dy*dy)
			nx1 := dx * math.Min(arcSize, dist/2) / dist
			ny1 := dy * math.Min(arcSize, dist/2) / dist
			c.lineTo(tmp.x+nx1, tmp.y+ny1)
			next := pts[mod(i+1, len(pts))]
			for i < len(pts)-2 && math.Round(next.x-tmp.x) == 0 && math.Round(next.y-tmp.y) == 0 {
				next = pts[mod(i+2, len(pts))]
				i++
			}
			dx, dy = next.x-tmp.x, next.y-tmp.y
			dist = math.Max(1, math.Sqrt(dx*dx+dy*dy))
			nx2 := dx * math.Min(arcSize, dist/2) / dist
			ny2 := dy * math.Min(arcSize, dist/2) / dist
			x2, y2 := tmp.x+nx2, tmp.y+ny2
			c.quadTo(tmp.x, tmp.y, x2, y2)
			tmp = point{x2, y2}
		} else {
			c.lineTo(tmp.x, tmp.y)
		}
		pt = tmp
		i++
	}
	if close {
		c.close()
	} else {
		c.lineTo(pe.x, pe.y)
	}
}

// paintGlassEffect paints the glossy highlight of glass=1 (mxShape.paintGlassEffect).
func (s *shape) paintGlassEffect(c *c2d, x, y, w, h, arc float64) {
	sw := math.Ceil(s.strokewidth / 2)
	const size = 0.4
	c.setGradient("#ffffff", "#ffffff", x, y, w, h*0.6, "south", 0.9, 0.1)
	c.begin()
	arc += 2 * sw
	if s.isRounded {
		c.moveTo(x-sw+arc, y-sw)
		c.quadTo(x-sw, y-sw, x-sw, y-sw+arc)
		c.lineTo(x-sw, y+h*size)
		c.quadTo(x+w*0.5, y+h*0.7, x+w+sw, y+h*size)
		c.lineTo(x+w+sw, y-sw+arc)
		c.quadTo(x+w+sw, y-sw, x+w+sw-arc, y-sw)
	} else {
		c.moveTo(x-sw, y-sw)
		c.lineTo(x-sw, y+h*size)
		c.quadTo(x+w*0.5, y+h*0.7, x+w+sw, y+h*size)
		c.lineTo(x+w+sw, y-sw)
	}
	c.close()
	c.fill()
}

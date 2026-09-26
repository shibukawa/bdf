package drawio

import (
	"math"
)

// The shapes of mxGraph's core (mxCellRenderer.defaultShapes), with the
// changes draw.io makes to them in Shapes.js.

func init() {
	rectangle := &shapeDef{
		paintBackground: paintRectangleBackground,
		paintForeground: paintRectangleForeground,
		labelBounds:     rectangleLabelBounds,
		roundable:       true,
	}
	registerShape("rectangle", rectangle)
	registerShape("label", &shapeDef{
		paintBackground: paintRectangleBackground,
		paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
			paintLabelImage(s, c, x, y, w, h)
			paintRectangleForeground(s, c, x, y, w, h)
		},
		labelBounds: rectangleLabelBounds,
		roundable:   true,
	})
	registerShape("ellipse", &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.ellipse(x, y, w, h)
		c.fillAndStroke()
		paintEllipseCenter(s, c, x, y, w, h)
		if s.glass && s.fill != "" {
			s.paintGlassEffect(c, x, y, w, h, 0)
		}
	}})
	registerShape("doubleEllipse", &shapeDef{
		paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
			c.ellipse(x, y, w, h)
			c.fillAndStroke()
		},
		paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
			m := s.style.num("margin", math.Min(3+s.strokewidth, math.Min(w/5, h/5)))
			x, y, w, h = x+m, y+m, w-2*m, h-2*m
			if w > 0 && h > 0 {
				c.ellipse(x, y, w, h)
			}
			c.stroke()
		},
		labelBounds: func(s *shape, r rect) rect {
			m := s.style.num("margin", math.Min(3+s.strokewidth, math.Min(r.w/5, r.h/5)))
			return rect{r.x + m, r.y + m, r.w - 2*m, r.h - 2*m}
		},
	})
	registerShape("rhombus", &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			paintRhombus(s, c, x, y, w, h)
			if s.style.is("double") {
				m := math.Max(2, s.strokewidth+1)*2 + s.style.num("margin", 0)
				if w-2*m > 0 && h-2*m > 0 {
					paintRhombus(s, c, x+m, y+m, w-2*m, h-2*m)
				}
			}
			if s.glass && s.fill != "" {
				s.paintGlassEffect(c, x, y, w, h, 0)
			}
		},
		labelBounds: func(s *shape, r rect) rect {
			if s.style.is("double") {
				m := math.Max(2, s.strokewidth+1)*2 + s.style.num("margin", 0)
				return rect{r.x + m, r.y + m, r.w - 2*m, r.h - 2*m}
			}
			return r
		},
		roundable: true,
	})
	registerShape("line", &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.begin()
		if s.style.is("vertical") {
			mid := x + w/2
			c.moveTo(mid, y)
			c.lineTo(mid, y+h)
		} else {
			mid := y + h/2
			c.moveTo(x, mid)
			c.lineTo(x+w, mid)
		}
		c.stroke()
	}})
	registerShape("image", &shapeDef{paintVertex: paintImageShape, roundable: true})
	registerShape("actor", actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		width := w / 3
		c.moveTo(0, h)
		c.curveTo(0, 3*h/5, 0, 2*h/5, w/2, 2*h/5)
		c.curveTo(w/2-width, 2*h/5, w/2-width, 0, w/2, 0)
		c.curveTo(w/2+width, 0, w/2+width, 2*h/5, w/2, 2*h/5)
		c.curveTo(w, 2*h/5, w, 3*h/5, w, h)
		c.close()
	}))
	tri := actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		arc := s.style.num("arcSize", lineArcSize) / 2
		s.addPoints(c, []point{{0, 0}, {w, 0.5 * h}, {0, h}}, s.isRounded, arc, true, nil, true)
	})
	tri.roundable = true
	registerShape("triangle", tri)
	registerShape("cloud", actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		c.moveTo(0.25*w, 0.25*h)
		c.curveTo(0.05*w, 0.25*h, 0, 0.5*h, 0.16*w, 0.55*h)
		c.curveTo(0, 0.66*h, 0.18*w, 0.9*h, 0.31*w, 0.8*h)
		c.curveTo(0.4*w, h, 0.7*w, h, 0.8*w, 0.8*h)
		c.curveTo(w, 0.8*h, w, 0.6*h, 0.875*w, 0.5*h)
		c.curveTo(w, 0.3*h, 0.8*w, 0.1*h, 0.625*w, 0.2*h)
		c.curveTo(0.5*w, 0.05*h, 0.3*w, 0.05*h, 0.25*w, 0.25*h)
		c.close()
	}))
	registerShape("cylinder", &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			c.translate(x, y)
			c.begin()
			cylinderPath(s, c, w, h, false)
			c.fillAndStroke()
			c.begin()
			cylinderPath(s, c, w, h, true)
			c.stroke()
		},
		labelMargins: func(s *shape, r rect) *rect {
			if s.style.is("boundedLbl") {
				size := cylinderSize(s, r.h) * 2
				return &rect{0, size, 0, 0}
			}
			return nil
		},
	})
	registerShape("connector", &shapeDef{paintEdge: paintConnector, noRotation: true, noInvert: true, roundable: true})
}

// actorShape is an mxActor: a path in coordinates relative to the bounds.
func actorShape(path func(s *shape, c *c2d, x, y, w, h float64)) *shapeDef {
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		c.begin()
		path(s, c, x, y, w, h)
		c.fillAndStroke()
	}}
}

func paintRectangleBackground(s *shape, c *c2d, x, y, w, h float64) {
	if s.isRounded {
		r := s.arcSize(w, h)
		c.roundrect(x, y, w, h, r, r)
	} else {
		c.rect(x, y, w, h)
	}
	c.fillAndStroke()
}

func paintRectangleForeground(s *shape, c *c2d, x, y, w, h float64) {
	footer := math.Min(h, math.Max(0, s.style.num("footerSize", 0)))
	if footer > 0 {
		paintRectangleFooter(s, c, x, y, w, h, footer)
	}
	if s.glass && s.fill != "" {
		s.paintGlassEffect(c, x, y, w, h, s.arcSize(w+s.strokewidth, h+s.strokewidth))
	}
}

func paintRectangleFooter(s *shape, c *c2d, x, y, w, h, footer float64) {
	color := s.style.get("footerColor", s.stroke)
	r := 0.0
	if s.isRounded {
		r = s.arcSize(w, h)
	}
	rr := math.Min(w/2, r)
	if color != "" && color != "none" {
		c.setFillColor(color)
		c.setStrokeColor(color)
	}
	c.begin()
	if footer >= r {
		c.moveTo(x, y+h-footer)
		c.lineTo(x, y+h-r)
		c.quadTo(x, y+h, x+rr, y+h)
		c.lineTo(x+w-rr, y+h)
		c.quadTo(x+w, y+h, x+w, y+h-r)
		c.lineTo(x+w, y+h-footer)
	} else {
		t := 1 - math.Sqrt(footer/r)
		px := rr * t * t
		bx := rr * t
		c.moveTo(x+px, y+h-footer)
		c.quadTo(x+bx, y+h, x+rr, y+h)
		c.lineTo(x+w-rr, y+h)
		c.quadTo(x+w-bx, y+h, x+w-px, y+h-footer)
	}
	c.close()
	if color != "" && color != "none" {
		c.fillAndStroke()
	} else {
		c.stroke()
	}
}

func rectangleLabelBounds(s *shape, r rect) rect {
	if footer := math.Max(0, s.style.num("footerSize", 0)); footer > 0 {
		r.h -= math.Min(r.h, footer)
	}
	return r
}

// paintLabelImage draws the image of a label shape (mxLabel.paintImage).
func paintLabelImage(s *shape, c *c2d, x, y, w, h float64) {
	src := s.style.get("image", "")
	if src == "" {
		return
	}
	img := s.conv.imageFor(src)
	if img == nil {
		return
	}
	width := s.style.num("imageWidth", 24)
	height := s.style.num("imageHeight", 24)
	spacing := s.style.num("spacing", 2) + 5
	switch s.style.get("imageAlign", "left") {
	case "center":
		x += (w - width) / 2
	case "right":
		x += w - width - spacing
	default:
		x += spacing
	}
	switch s.style.get("imageVerticalAlign", "middle") {
	case "top":
		y += spacing
	case "bottom":
		y += h - height - spacing
	default:
		y += (h - height) / 2
	}
	c.image(x, y, width, height, img, false, false, false)
}

func paintEllipseCenter(s *shape, c *c2d, x, y, w, h float64) {
	r := math.Min(math.Max(0, s.style.num("centerRadius", 0)), math.Min(w/2, h/2))
	if r <= 0 {
		return
	}
	color := s.style.get("centerColor", s.stroke)
	if color == "" || color == "none" {
		return
	}
	c.setFillColor(color)
	c.ellipse(x+w/2-r, y+h/2-r, 2*r, 2*r)
	if s.style.is("centerStroke") && s.stroke != "" {
		c.setStrokeColor(s.stroke)
		c.fillAndStroke()
	} else {
		c.fill()
	}
}

func paintRhombus(s *shape, c *c2d, x, y, w, h float64) {
	hw, hh := w/2, h/2
	arc := s.style.num("arcSize", lineArcSize) / 2
	c.begin()
	s.addPoints(c, []point{{x + hw, y}, {x + w, y + hh}, {x + hw, y + h}, {x, y + hh}}, s.isRounded, arc, true, nil, true)
	c.fillAndStroke()
}

// paintImageShape draws an image shape (mxImageShape.paintVertexShape).
func paintImageShape(s *shape, c *c2d, x, y, w, h float64) {
	src := s.style.get("image", "")
	var img *imageRef
	if src != "" {
		img = s.conv.imageFor(src)
	}
	if img == nil {
		if src == "" {
			paintRectangleBackground(s, c, x, y, w, h)
		}
		return
	}
	r := 0.0
	if s.isRounded {
		r = s.arcSize(w, h)
	}
	bg := colorOrNone(s.style.get("imageBackground", ""))
	border := colorOrNone(s.style.get("imageBorder", ""))
	if bg != "" {
		c.setFillColor(bg)
		c.setStrokeColor(border)
		if r > 0 {
			c.roundrect(x, y, w, h, r, r)
		} else {
			c.rect(x, y, w, h)
		}
		c.fillAndStroke()
	}
	aspect := s.style.num("imageAspect", 1) == 1
	flipH := s.flipH || s.style.is("imageFlipH")
	flipV := s.flipV || s.style.is("imageFlipV")
	// the shape's flips are applied by the canvas already
	c.image(x, y, w, h, img, aspect, flipH != s.flipH, flipV != s.flipV)
	if border != "" {
		c.setStrokeColor(border)
		if r > 0 {
			c.roundrect(x, y, w, h, r, r)
		} else {
			c.rect(x, y, w, h)
		}
		c.stroke()
	}
}

// cylinderSize is the height of the cylinder's top ellipse (mxCylinder.getCylinderSize,
// as draw.io changes it: the size style, else a fifth of the height up to 40).
func cylinderSize(s *shape, h float64) float64 {
	if v, ok := s.style["size"]; ok {
		if f, ok := parseFloat(v); ok {
			return f
		}
	}
	return math.Min(40, math.Round(h/5))
}

func cylinderPath(s *shape, c *c2d, w, h float64, foreground bool) {
	dy := cylinderSize(s, h)
	if (foreground && s.fill != "") || (!foreground && s.fill == "") {
		c.moveTo(0, dy)
		c.curveTo(0, 2*dy, w, 2*dy, w, dy)
		if !foreground {
			c.stroke()
			c.begin()
		}
	}
	if !foreground {
		c.moveTo(0, dy)
		c.curveTo(0, -dy/3, w, -dy/3, w, dy)
		c.lineTo(w, h-dy)
		c.curveTo(w, h+dy/3, 0, h+dy/3, 0, h-dy)
		c.close()
	}
}

// --- edges ---

// paintConnector draws an edge with its end markers (mxConnector.paintEdgeShape).
func paintConnector(s *shape, c *c2d, pts []point) {
	pts = append([]point(nil), pts...)
	srcMarker := createMarker(s, c, pts, true)
	trgMarker := createMarker(s, c, pts, false)
	paintPolyline(s, c, pts)
	c.setDashed(false, false)
	if srcMarker != nil {
		c.setFillColor(s.style.get("startFillColor", s.stroke))
		srcMarker()
	}
	if trgMarker != nil {
		c.setFillColor(s.style.get("endFillColor", s.stroke))
		trgMarker()
	}
}

// paintPolyline strokes the line of an edge: straight, rounded, curved or
// Bézier (mxPolyline.paintEdgeShape).
func paintPolyline(s *shape, c *c2d, pts []point) {
	switch {
	case s.style.is("bezier"):
		n := len(pts)
		c.begin()
		c.moveTo(pts[0].x, pts[0].y)
		if n == 2 {
			c.lineTo(pts[1].x, pts[1].y)
		} else if (n-1)%3 == 0 {
			for i := 1; i+2 < n; i += 3 {
				c.curveTo(pts[i].x, pts[i].y, pts[i+1].x, pts[i+1].y, pts[i+2].x, pts[i+2].y)
			}
		} else {
			for i := 1; i < n-2; i++ {
				p0, p1 := pts[i], pts[i+1]
				c.quadTo(p0.x, p0.y, (p0.x+p1.x)/2, (p0.y+p1.y)/2)
			}
			c.quadTo(pts[n-2].x, pts[n-2].y, pts[n-1].x, pts[n-1].y)
		}
		c.stroke()
	case s.style.is("curved"):
		c.begin()
		n := len(pts)
		c.moveTo(pts[0].x, pts[0].y)
		for i := 1; i < n-2; i++ {
			p0, p1 := pts[i], pts[i+1]
			c.quadTo(p0.x, p0.y, (p0.x+p1.x)/2, (p0.y+p1.y)/2)
		}
		p0, p1 := pts[n-2], pts[n-1]
		c.quadTo(p0.x, p0.y, p1.x, p1.y)
		c.stroke()
	default:
		if paintLineJumps(s, c, pts, s.isRounded) {
			return
		}
		arc := s.style.num("arcSize", lineArcSize) / 2
		c.begin()
		s.addPoints(c, pts, s.isRounded, arc, false, nil, true)
		c.stroke()
	}
}

// createMarker prepares the marker at one end of an edge, moving that end
// back to where the line stops (mxConnector.createMarker). It returns nil
// when the end has no marker.
func createMarker(s *shape, c *c2d, pts []point, source bool) func() {
	n := len(pts)
	typ := s.endArrow
	p0, pe := &pts[n-2], &pts[n-1]
	sizeKey := "endSize"
	fillKey := "endFill"
	if source {
		typ = s.startArrow
		p0, pe = &pts[1], &pts[0]
		sizeKey = "startSize"
		fillKey = "startFill"
	}
	if typ == "" || typ == "none" {
		return nil
	}
	dx, dy := pe.x-p0.x, pe.y-p0.y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return nil
	}
	unitX, unitY := dx/dist, dy/dist
	size := s.style.num(sizeKey, defaultMarkerSize)
	filled := s.style.get(fillKey, "1") != "0"
	f, ok := markers[typ]
	if !ok {
		s.conv.warnOnce("marker:"+typ, "arrow %q is not supported; drawn as classic", typ)
		f = markers["classic"]
	}
	return f(c, s, typ, pe, unitX, unitY, size, source, s.strokewidth, filled)
}

// markerFunc prepares a marker at pe heading along (unitX, unitY) and may
// move pe back; the returned function paints it.
type markerFunc func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func()

var markers = map[string]markerFunc{}

func init() {
	markers["classic"] = createArrow(2)
	markers["classicThin"] = createArrow(3)
	markers["block"] = createArrow(2)
	markers["blockThin"] = createArrow(3)
	markers["open"] = createOpenArrow(2)
	markers["openThin"] = createOpenArrow(3)
	markers["oval"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		a := size / 2
		pt := *pe
		pe.x -= unitX * a
		pe.y -= unitY * a
		return func() {
			c.ellipse(pt.x-a, pt.y-a, size, size)
			fillOrStroke(c, filled)
		}
	}
	markers["baseDash"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		nx, ny := unitX*(size+sw+1), unitY*(size+sw+1)
		p := *pe
		return func() {
			c.begin()
			c.moveTo(p.x-ny/2, p.y+nx/2)
			c.lineTo(p.x+ny/2, p.y-nx/2)
			c.stroke()
		}
	}
	markers["doubleBlock"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		const widthFactor = 2
		endOffsetX, endOffsetY := unitX*sw*1.118, unitY*sw*1.118
		unitX *= size + sw
		unitY *= size + sw
		pt := point{pe.x - endOffsetX, pe.y - endOffsetY}
		pe.x += -unitX*2 - endOffsetX
		pe.y += -unitY*2 - endOffsetY
		return func() {
			c.begin()
			c.moveTo(pt.x, pt.y)
			c.lineTo(pt.x-unitX-unitY/widthFactor, pt.y-unitY+unitX/widthFactor)
			c.lineTo(pt.x+unitY/widthFactor-unitX, pt.y-unitY-unitX/widthFactor)
			c.close()
			c.moveTo(pt.x-unitX, pt.y-unitY)
			c.lineTo(pt.x-2*unitX-0.5*unitY, pt.y+0.5*unitX-2*unitY)
			c.lineTo(pt.x-2*unitX+0.5*unitY, pt.y-0.5*unitX-2*unitY)
			c.close()
			fillOrStroke(c, filled)
		}
	}
	markers["diamond"] = diamondMarker
	markers["diamondThin"] = diamondMarker
	markers["manyOptional"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		nx, ny := unitX*(size+sw+1), unitY*(size+sw+1)
		a := size / 2
		px, py := pe.x, pe.y
		pe.x -= 2*nx - unitX*sw/2
		pe.y -= 2*ny - unitY*sw/2
		return func() {
			c.begin()
			c.ellipse(px-1.5*nx-a, py-1.5*ny-a, 2*a, 2*a)
			fillOrStroke(c, filled)
			c.begin()
			c.moveTo(px, py)
			c.lineTo(px-nx, py-ny)
			c.moveTo(px+ny/2, py-nx/2)
			c.lineTo(px-nx, py-ny)
			c.lineTo(px-ny/2, py+nx/2)
			c.stroke()
		}
	}
}

func fillOrStroke(c *c2d, filled bool) {
	if filled {
		c.fillAndStroke()
	} else {
		c.stroke()
	}
}

func createArrow(widthFactor float64) markerFunc {
	return func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		endOffsetX, endOffsetY := unitX*sw*1.118, unitY*sw*1.118
		unitX *= size + sw
		unitY *= size + sw
		pt := point{pe.x - endOffsetX, pe.y - endOffsetY}
		f := 1.0
		if typ == "classic" || typ == "classicThin" {
			f = 3.0 / 4
		}
		pe.x += -unitX*f - endOffsetX
		pe.y += -unitY*f - endOffsetY
		return func() {
			c.begin()
			c.moveTo(pt.x, pt.y)
			c.lineTo(pt.x-unitX-unitY/widthFactor, pt.y-unitY+unitX/widthFactor)
			if typ == "classic" || typ == "classicThin" {
				c.lineTo(pt.x-unitX*3/4, pt.y-unitY*3/4)
			}
			c.lineTo(pt.x+unitY/widthFactor-unitX, pt.y-unitY-unitX/widthFactor)
			c.close()
			fillOrStroke(c, filled)
		}
	}
}

func createOpenArrow(widthFactor float64) markerFunc {
	return func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		endOffsetX, endOffsetY := unitX*sw*1.118, unitY*sw*1.118
		unitX *= size + sw
		unitY *= size + sw
		pt := point{pe.x - endOffsetX, pe.y - endOffsetY}
		pe.x += -endOffsetX * 2
		pe.y += -endOffsetY * 2
		return func() {
			c.begin()
			c.moveTo(pt.x-unitX-unitY/widthFactor, pt.y-unitY+unitX/widthFactor)
			c.lineTo(pt.x, pt.y)
			c.lineTo(pt.x+unitY/widthFactor-unitX, pt.y-unitY-unitX/widthFactor)
			c.stroke()
		}
	}
}

func diamondMarker(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
	swFactor := 0.9862
	tk := 3.4
	if typ == "diamond" {
		swFactor, tk = 0.7071, 2
	}
	endOffsetX, endOffsetY := unitX*sw*swFactor, unitY*sw*swFactor
	unitX *= size + sw
	unitY *= size + sw
	pt := point{pe.x - endOffsetX, pe.y - endOffsetY}
	pe.x += -unitX - endOffsetX
	pe.y += -unitY - endOffsetY
	return func() {
		c.begin()
		c.moveTo(pt.x, pt.y)
		c.lineTo(pt.x-unitX/2-unitY/tk, pt.y+unitX/tk-unitY/2)
		c.lineTo(pt.x-unitX, pt.y-unitY)
		c.lineTo(pt.x-unitX/2+unitY/tk, pt.y-unitY/2-unitX/tk)
		c.close()
		fillOrStroke(c, filled)
	}
}

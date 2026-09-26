package drawio

import "math"

// Edge shapes other than the connector: mxArrow, mxArrowConnector and
// draw.io's LinkShape, FlexArrowShape, FilledEdge, WireShape and
// PipeShape (Shapes.js).

func init() {
	registerShape("arrow", &shapeDef{paintEdge: paintArrow, augmentBounds: arrowBounds})
	registerShape("arrowConnector", arrowConnectorShape(arrowConnector{}))
	registerShape("link", arrowConnectorShape(arrowConnector{link: true}))
	registerShape("flexArrow", arrowConnectorShape(arrowConnector{flex: true}))
	registerShape("filledEdge", &shapeDef{paintEdge: paintFilledEdge, noRotation: true, noInvert: true, roundable: true})
	registerShape("wire", &shapeDef{paintEdge: paintWire, noRotation: true, noInvert: true, roundable: true})
	registerShape("pipe", &shapeDef{paintEdge: paintPipe, noRotation: true, noInvert: true, roundable: true,
		augmentBounds: func(s *shape, r rect) rect {
			return r.union(pointsBounds(s.points).grow(pipeWidth(s)/2 + s.strokewidth))
		}})
}

// mxConstants.ARROW_SPACING, ARROW_WIDTH and ARROW_SIZE.
const (
	arrowSpacing = 0
	arrowWidth   = 30
	arrowSize    = 30
)

// pointsBounds returns the bounding box of points (zero for none).
func pointsBounds(pts []point) rect {
	r, _ := pointsExtent(pts)
	return r
}

// arrowBounds is mxArrow.augmentBoundingBox.
func arrowBounds(s *shape, r rect) rect {
	return r.union(pointsBounds(s.points).grow(math.Max(arrowWidth, arrowSize)/2 + s.strokewidth))
}

// paintArrow draws a block arrow from the first to the last point
// (mxArrow.paintEdgeShape).
func paintArrow(s *shape, c *c2d, pts []point) {
	const spacing, width, arrow = arrowSpacing, arrowWidth, arrowSize
	p0, pe := pts[0], pts[len(pts)-1]
	dx, dy := pe.x-p0.x, pe.y-p0.y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return
	}
	length := dist - 2*spacing - arrow
	nx, ny := dx/dist, dy/dist
	basex, basey := length*nx, length*ny
	floorx, floory := width*ny/3, -width*nx/3
	p0x := p0.x - floorx/2 + spacing*nx
	p0y := p0.y - floory/2 + spacing*ny
	p1x, p1y := p0x+floorx, p0y+floory
	p2x, p2y := p1x+basex, p1y+basey
	p3x, p3y := p2x+floorx, p2y+floory
	// p4 is not needed
	p5x, p5y := p3x-3*floorx, p3y-3*floory
	c.begin()
	c.moveTo(p0x, p0y)
	c.lineTo(p1x, p1y)
	c.lineTo(p2x, p2y)
	c.lineTo(p3x, p3y)
	c.lineTo(pe.x-spacing*nx, pe.y-spacing*ny)
	c.lineTo(p5x, p5y)
	c.lineTo(p5x+floorx, p5y+floory)
	c.close()
	c.fillAndStroke()
}

// arrowConnector selects the variant of mxArrowConnector: LinkShape (link,
// an open double line) or FlexArrowShape (flex, sizes from the style).
type arrowConnector struct {
	link, flex bool
}

func (a arrowConnector) edgeWidth(s *shape) float64 {
	switch {
	case a.link:
		return s.style.num("width", 4) + math.Max(0, s.strokewidth-1)
	case a.flex:
		return s.style.num("width", 10) + math.Max(0, s.strokewidth-1)
	}
	return arrowWidth / 3
}

func (a arrowConnector) startArrowWidth(s *shape) float64 {
	if a.flex {
		return a.edgeWidth(s) + s.style.num("startWidth", 20)
	}
	return arrowWidth
}

func (a arrowConnector) endArrowWidth(s *shape) float64 {
	if a.flex {
		return a.edgeWidth(s) + s.style.num("endWidth", 20)
	}
	return arrowWidth
}

func isMarkerStart(s *shape) bool { return s.style.get("startArrow", "none") != "none" }
func isMarkerEnd(s *shape) bool   { return s.style.get("endArrow", "none") != "none" }

// arrowConnectorShape is mxArrowConnector with a's variations.
func arrowConnectorShape(a arrowConnector) *shapeDef {
	return &shapeDef{
		paintEdge: a.paint,
		augmentBounds: func(s *shape, r rect) rect {
			w := a.edgeWidth(s)
			if isMarkerStart(s) {
				w = math.Max(w, a.startArrowWidth(s))
			}
			if isMarkerEnd(s) {
				w = math.Max(w, a.endArrowWidth(s))
			}
			return r.union(pointsBounds(s.points).grow(w/2 + s.strokewidth))
		},
		roundable: true,
	}
}

// paint draws the outline of a wide arrow along the points
// (mxArrowConnector.paintEdgeShape).
func (a arrowConnector) paint(s *shape, c *c2d, pts []point) {
	// a curve is approximated with a fine polyline
	isCurved := s.style.num("curved", 0) == 1 && len(pts) > 2
	if isCurved {
		pts = curvePoints(pts)
	}
	strokeWidth := s.strokewidth
	startWidth := a.startArrowWidth(s) + strokeWidth
	endWidth := a.endArrowWidth(s) + strokeWidth
	edgeWidth := a.edgeWidth(s)
	openEnded := a.link
	markerStart := isMarkerStart(s)
	markerEnd := isMarkerEnd(s)
	spacing := 0.0
	if !openEnded {
		spacing = arrowSpacing + strokeWidth/2
	}
	// mxArrowConnector.apply
	startSize := s.style.num("startSize", arrowSize/5)*3 + strokeWidth
	endSize := s.style.num("endSize", arrowSize/5)*3 + strokeWidth
	// rounded joins overshoot the short segments of a curve
	isRounded := !isCurved && s.isRounded
	if isCurved && !openEnded {
		// the markers are straight: the curve is cut where they start
		if markerStart {
			pts = trimCurveForMarker(pts, spacing+startSize, true)
		}
		if markerEnd {
			pts = trimCurveForMarker(pts, spacing+endSize, false)
		}
	}
	pe := pts[len(pts)-1]
	dx, dy := pts[1].x-pts[0].x, pts[1].y-pts[0].y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return
	}
	nx, ny := dx/dist, dy/dist
	nx1, ny1 := nx, ny
	orthx, orthy := edgeWidth*ny, -edgeWidth*nx
	// the inbound path, painted in reverse order at the end
	var fns []func()
	if isRounded {
		c.setLineJoin("round")
	} else if len(pts) > 2 {
		// only mitre with waypoints
		c.setMiterLimit(1.42)
	}
	c.begin()
	startNx, startNy := nx, ny
	if markerStart && !openEnded {
		paintArrowMarker(c, pts[0].x, pts[0].y, nx, ny, startSize, startWidth, edgeWidth, spacing, true)
	} else {
		outStartX := pts[0].x + orthx/2 + spacing*nx
		outStartY := pts[0].y + orthy/2 + spacing*ny
		inEndX := pts[0].x - orthx/2 + spacing*nx
		inEndY := pts[0].y - orthy/2 + spacing*ny
		if openEnded {
			c.moveTo(outStartX, outStartY)
			fns = append(fns, func() { c.lineTo(inEndX, inEndY) })
		} else {
			c.moveTo(inEndX, inEndY)
			c.lineTo(outStartX, outStartY)
		}
	}
	for i := 0; i < len(pts)-2; i++ {
		// the direction the line bends
		pos := relativeCcw(pts[i].x, pts[i].y, pts[i+1].x, pts[i+1].y, pts[i+2].x, pts[i+2].y)
		dx1, dy1 := pts[i+2].x-pts[i+1].x, pts[i+2].y-pts[i+1].y
		dist1 := math.Sqrt(dx1*dx1 + dy1*dy1)
		if dist1 == 0 {
			continue
		}
		nx1, ny1 = dx1/dist1, dy1/dist1
		tmp1 := nx*nx1 + ny*ny1
		tmp := math.Max(math.Sqrt((tmp1+1)/2), 0.04)
		// the normal of the line through the control point and where the
		// sides of the edge intersect
		nx2, ny2 := nx+nx1, ny+ny1
		dist2 := math.Sqrt(nx2*nx2 + ny2*ny2)
		if dist2 == 0 {
			continue
		}
		nx2 /= dist2
		ny2 /= dist2
		// higher stroke widths need a larger minimum bend
		strokeWidthFactor := math.Max(tmp, math.Min(s.strokewidth/200+0.04, 0.35))
		angleFactor := math.Max(tmp, 0.06)
		if pos != 0 && isRounded {
			angleFactor = math.Max(0.1, strokeWidthFactor)
		}
		p := pts[i+1]
		outX := p.x + ny2*edgeWidth/2/angleFactor
		outY := p.y - nx2*edgeWidth/2/angleFactor
		inX := p.x - ny2*edgeWidth/2/angleFactor
		inY := p.y + nx2*edgeWidth/2/angleFactor
		switch {
		case pos == 0 || !isRounded:
			// straight to the intersection of the sides
			c.lineTo(outX, outY)
			fns = append(fns, func() { c.lineTo(inX, inY) })
		case pos == -1:
			c1x, c1y := inX+ny*edgeWidth, inY-nx*edgeWidth
			c2x, c2y := inX+ny1*edgeWidth, inY-nx1*edgeWidth
			c.lineTo(c1x, c1y)
			c.quadTo(outX, outY, c2x, c2y)
			fns = append(fns, func() { c.lineTo(inX, inY) })
		default:
			c.lineTo(outX, outY)
			c1x, c1y := outX-ny*edgeWidth, outY+nx*edgeWidth
			c2x, c2y := outX-ny1*edgeWidth, outY+nx1*edgeWidth
			fns = append(fns, func() { c.quadTo(inX, inY, c1x, c1y) })
			fns = append(fns, func() { c.lineTo(c2x, c2y) })
		}
		nx, ny = nx1, ny1
	}
	orthx, orthy = edgeWidth*ny1, -edgeWidth*nx1
	if markerEnd && !openEnded {
		paintArrowMarker(c, pe.x, pe.y, -nx, -ny, endSize, endWidth, edgeWidth, spacing, false)
	} else {
		c.lineTo(pe.x-spacing*nx1+orthx/2, pe.y-spacing*ny1+orthy/2)
		inStartX := pe.x - spacing*nx1 - orthx/2
		inStartY := pe.y - spacing*ny1 - orthy/2
		if !openEnded {
			c.lineTo(inStartX, inStartY)
		} else {
			c.moveTo(inStartX, inStartY)
			fns = append([]func(){func() { c.moveTo(inStartX, inStartY) }}, fns...)
		}
	}
	for i := len(fns) - 1; i >= 0; i-- {
		fns[i]()
	}
	if openEnded {
		c.stroke()
	} else {
		c.close()
		c.fillAndStroke()
	}
	c.setMiterLimit(4)
	if isRounded {
		c.setLineJoin("flat")
	}
	if len(pts) > 2 {
		// repaints the markers without the low miter limit
		if markerStart && !openEnded {
			c.begin()
			paintArrowMarker(c, pts[0].x, pts[0].y, startNx, startNy, startSize, startWidth, edgeWidth, spacing, true)
			c.stroke()
		}
		if markerEnd && !openEnded {
			c.begin()
			paintArrowMarker(c, pe.x, pe.y, -nx, -ny, endSize, endWidth, edgeWidth, spacing, true)
			c.stroke()
		}
	}
}

// paintArrowMarker adds the head of the arrow to the path
// (mxArrowConnector.paintMarker).
func paintArrowMarker(c *c2d, ptX, ptY, nx, ny, size, arrowWidth, edgeWidth, spacing float64, initialMove bool) {
	widthArrowRatio := edgeWidth / arrowWidth
	orthx, orthy := edgeWidth*ny/2, -edgeWidth*nx/2
	spaceX, spaceY := (spacing+size)*nx, (spacing+size)*ny
	if initialMove {
		c.moveTo(ptX-orthx+spaceX, ptY-orthy+spaceY)
	} else {
		c.lineTo(ptX-orthx+spaceX, ptY-orthy+spaceY)
	}
	c.lineTo(ptX-orthx/widthArrowRatio+spaceX, ptY-orthy/widthArrowRatio+spaceY)
	c.lineTo(ptX+spacing*nx, ptY+spacing*ny)
	c.lineTo(ptX+orthx/widthArrowRatio+spaceX, ptY+orthy/widthArrowRatio+spaceY)
	c.lineTo(ptX+orthx+spaceX, ptY+orthy+spaceY)
}

// relativeCcw is mxUtils.relativeCcw: the side of (px, py) relative to
// the line from (x1, y1) to (x2, y2) as -1, 0 or 1.
func relativeCcw(x1, y1, x2, y2, px, py float64) int {
	x2 -= x1
	y2 -= y1
	px -= x1
	py -= y1
	ccw := px*y2 - py*x2
	if ccw == 0 {
		ccw = px*x2 + py*y2
		if ccw > 0 {
			px -= x2
			py -= y2
			ccw = px*x2 + py*y2
			if ccw < 0 {
				ccw = 0
			}
		}
	}
	switch {
	case ccw < 0:
		return -1
	case ccw > 0:
		return 1
	}
	return 0
}

// curvePoints samples the curve mxPolyline.paintCurvedLine paints through
// pts (a quadratic spline through the midpoints) as a fine polyline
// (mxArrowConnector.getCurvePoints).
func curvePoints(pts []point) []point {
	result := []point{pts[0]}
	n := len(pts)
	p0 := pts[0]
	for i := 1; i < n-1; i++ {
		pc := pts[i]
		pe := pts[i+1]
		if i < n-2 {
			pe = point{(pts[i].x + pts[i+1].x) / 2, (pts[i].y + pts[i+1].y) / 2}
		}
		l := math.Hypot(pc.x-p0.x, pc.y-p0.y) + math.Hypot(pe.x-pc.x, pe.y-pc.y)
		steps := math.Max(4, math.Min(math.Ceil(l/8), 64))
		if !finite(steps) {
			steps = 4
		}
		for j := 1; j <= int(steps); j++ {
			u := float64(j) / steps
			iu := 1 - u
			result = append(result, point{
				iu*iu*p0.x + 2*iu*u*pc.x + u*u*pe.x,
				iu*iu*p0.y + 2*iu*u*pc.y + u*u*pe.y})
		}
		p0 = pe
	}
	return result
}

// trimCurveForMarker replaces the part of the polyline within dist of the
// source (or target) point by the chord to the point at that distance, on
// which the marker is painted (mxArrowConnector.trimCurveForMarker).
func trimCurveForMarker(pts []point, dist float64, source bool) []point {
	n := len(pts)
	pc := pts[n-1]
	if source {
		pc = pts[0]
	}
	// the point between the outer p0 and the inner p1 at dist from pc
	cutPoint := func(p0, p1 point) point {
		dx, dy := p1.x-p0.x, p1.y-p0.y
		fx, fy := p0.x-pc.x, p0.y-pc.y
		a := dx*dx + dy*dy
		b := 2 * (fx*dx + fy*dy)
		cc := fx*fx + fy*fy - dist*dist
		disc := b*b - 4*a*cc
		if a == 0 || disc < 0 {
			return p0
		}
		t := clamp((-b-math.Sqrt(disc))/(2*a), 0, 1)
		return point{p0.x + t*dx, p0.y + t*dy}
	}
	if source {
		for i := 1; i < n-1; i++ {
			dx, dy := pts[i].x-pc.x, pts[i].y-pc.y
			if dx*dx+dy*dy >= dist*dist {
				return append([]point{pc, cutPoint(pts[i], pts[i-1])}, pts[i:]...)
			}
		}
	} else {
		for i := n - 2; i > 0; i-- {
			dx, dy := pts[i].x-pc.x, pts[i].y-pc.y
			if dx*dx+dy*dy >= dist*dist {
				result := append([]point(nil), pts[:i+1]...)
				return append(result, cutPoint(pts[i], pts[i+1]), pc)
			}
		}
	}
	return pts
}

// paintFilledEdge draws a connector and, when the stroke is at least 3
// wide and a fillColor is set, a connector 2 narrower in the fill color
// on top of it (FilledEdge.paintEdgeShape).
func paintFilledEdge(s *shape, c *c2d, pts []point) {
	// paintConnector resets dashed
	dashed, fixDash := c.st.dashed, c.st.fixDash
	paintConnector(s, c, pts)
	if c.st.strokeWidth >= 3 {
		if fill, ok := s.style["fillColor"]; ok {
			c.setStrokeColor(fill)
			c.setStrokeWidth(c.st.strokeWidth - 2)
			c.setDashed(dashed, fixDash)
			paintConnector(s, c, pts)
		}
	}
}

// paintWire draws a line in the stroke color, a dashed line in the fill
// color over it and the markers in the stroke color (WireShape.paintEdgeShape).
func paintWire(s *shape, c *c2d, pts []point) {
	pts = append([]point(nil), pts...)
	srcMarker := createMarker(s, c, pts, true)
	trgMarker := createMarker(s, c, pts, false)
	c.setDashed(false, false)
	paintPolyline(s, c, pts)
	c.setDashed(s.isDashed, s.style.num("fixDash", 0) == 1)
	c.setStrokeColor(s.fill)
	paintPolyline(s, c, pts)
	c.setStrokeColor(s.stroke)
	c.setFillColor(s.stroke)
	c.setDashed(false, false)
	if srcMarker != nil {
		srcMarker()
	}
	if trgMarker != nil {
		trgMarker()
	}
}

// pipeWidth is PipeShape.getEdgeWidth.
func pipeWidth(s *shape) float64 { return math.Max(0, s.style.num("width", 4)) }

// paintPipe draws a connector as a casing in the stroke color and, with a
// fillColor, the inside of the pipe (PipeShape.paintEdgeShape).
func paintPipe(s *shape, c *c2d, pts []point) {
	dashed, fixDash := c.st.dashed, c.st.fixDash
	width := pipeWidth(s)
	c.setStrokeWidth(width + 2*s.strokewidth)
	c.setDashed(false, false)
	paintConnector(s, c, pts)
	if fill, ok := s.style["fillColor"]; ok {
		c.setStrokeWidth(width)
		c.setStrokeColor(fill)
		c.setDashed(dashed, fixDash)
		paintConnector(s, c, pts)
	}
}

package drawio

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// draw.io's basic shape library (mxgraph.basic.*, shapes/mxBasic.js),
// ported class by class. Every class there extends mxActor and overrides
// paintVertexShape, which starts with c.translate(x, y); basicShape does
// that and leaves the rest to the port. The parameters (dx, dy, size,
// startAngle, …) default to what the constructors set, not to the defVal
// of the customProperties that the format panel writes when they are
// edited. mxBasic.js registers no perimeters.
//
// mxgraph.basic.rect is not in mxBasic.js: it is mxShapeBasicRect2 of
// Shapes.js (shapes_rect2.go). The other mxgraph.basic names (cross,
// wave, octagon, cone, …) are stencils of stencils/basic.xml.

func init() {
	registerShape("mxgraph.basic.cross2", basicCrossShape)
	registerShape("mxgraph.basic.rectCallout", basicRectCalloutShape)
	registerShape("mxgraph.basic.roundRectCallout", basicRoundRectCalloutShape)
	registerShape("mxgraph.basic.wave2", basicWaveShape)
	registerShape("mxgraph.basic.octagon2", basicOctagonShape)
	// mxShapeBasicIsoCube paints exactly as IsoCubeShape2 of Shapes.js
	registerShape("mxgraph.basic.isocube", isoCube2Shape)
	registerShape("mxgraph.basic.acute_triangle", basicTriangleAcuteShape)
	registerShape("mxgraph.basic.obtuse_triangle", basicTriangleObtuseShape)
	registerShape("mxgraph.basic.drop", basicDropShape)
	registerShape("mxgraph.basic.cone2", basicCone2Shape)
	registerShape("mxgraph.basic.pyramid", basicPyramidShape)
	registerShape("mxgraph.basic.4_point_star_2", basic4PointStar2Shape)
	registerShape("mxgraph.basic.diag_snip_rect", basicDiagSnipRectShape)
	registerShape("mxgraph.basic.diag_round_rect", basicDiagRoundRectShape)
	registerShape("mxgraph.basic.corner_round_rect", basicCornerRoundRectShape)
	registerShape("mxgraph.basic.plaque", basicPlaqueShape)
	registerShape("mxgraph.basic.frame", basicFrameShape)
	registerShape("mxgraph.basic.plaque_frame", basicPlaqueFrameShape)
	registerShape("mxgraph.basic.rounded_frame", basicRoundedFrameShape)
	registerShape("mxgraph.basic.frame_corner", basicFrameCornerShape)
	registerShape("mxgraph.basic.diag_stripe", basicDiagStripeShape)
	registerShape("mxgraph.basic.donut", basicDonutShape)
	registerShape("mxgraph.basic.layered_rect", basicLayeredRectShape)
	registerShape("mxgraph.basic.button", basicButtonShape)
	registerShape("mxgraph.basic.shaded_button", basicShadedButtonShape)
	registerShape("mxgraph.basic.pie", basicPieShape)
	registerShape("mxgraph.basic.arc", basicArcShape)
	registerShape("mxgraph.basic.partConcEllipse", basicPartConcEllipseShape)
	registerShape("mxgraph.basic.numberedEntryVert", basicNumEntryVertShape)
	registerShape("mxgraph.basic.bendingArch", basicBendingArchShape)
	registerShape("mxgraph.basic.three_corner_round_rect", basicThreeCornerRoundRectShape)
	registerShape("mxgraph.basic.polygon", basicPolygonShape)
	registerShape("mxgraph.basic.patternFillRect", basicPatternFillRectShape)
}

// basicShape is a class of mxBasic.js: paint draws in coordinates
// relative to the bounds.
func basicShape(paint func(s *shape, c *c2d, w, h float64)) *shapeDef {
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		paint(s, c, w, h)
	}}
}

func min3(a, b, c float64) float64 { return math.Min(a, math.Min(b, c)) }

// basicCrossShape ports mxShapeBasicCross ("cross2"): a plus sign whose
// bars are 2*dx wide.
var basicCrossShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	c.begin()
	c.moveTo(w*0.5+dx, 0)
	c.lineTo(w*0.5+dx, h*0.5-dx)
	c.lineTo(w, h*0.5-dx)
	c.lineTo(w, h*0.5+dx)
	c.lineTo(w*0.5+dx, h*0.5+dx)
	c.lineTo(w*0.5+dx, h)
	c.lineTo(w*0.5-dx, h)
	c.lineTo(w*0.5-dx, h*0.5+dx)
	c.lineTo(0, h*0.5+dx)
	c.lineTo(0, h*0.5-dx)
	c.lineTo(w*0.5-dx, h*0.5-dx)
	c.lineTo(w*0.5-dx, 0)
	c.close()
	c.fillAndStroke()
})

// basicCalloutLabelMargins is mxShapeBasicRectCallout.getLabelMargins
// (shared by mxShapeBasicRoundRectCallout): boundedLbl keeps the label
// out of the tail (dy, unclamped).
func basicCalloutLabelMargins(s *shape, r rect) *rect {
	if boundedLbl(s) {
		return &rect{0, 0, 0, s.style.num("dy", 0.5)}
	}
	return nil
}

// basicRectCalloutShape ports mxShapeBasicRectCallout: a rectangle with a
// tail of height dy at dx on the bottom.
var basicRectCalloutShape = func() *shapeDef {
	d := basicShape(func(s *shape, c *c2d, w, h float64) {
		dx := clamp(s.style.num("dx", 0.5), 0, w)
		dy := clamp(s.style.num("dy", 0.5), 0, h)
		c.begin()
		c.moveTo(dx-dy*0.5, h-dy)
		c.lineTo(0, h-dy)
		c.lineTo(0, 0)
		c.lineTo(w, 0)
		c.lineTo(w, h-dy)
		c.lineTo(dx+dy*0.5, h-dy)
		c.lineTo(dx-dy, h)
		c.close()
		c.fillAndStroke()
	})
	d.labelMargins = basicCalloutLabelMargins
	return d
}()

// basicRoundRectCalloutShape ports mxShapeBasicRoundRectCallout: the
// rectangular callout with corners of radius size and a curved tail.
var basicRoundRectCalloutShape = func() *shapeDef {
	d := basicShape(func(s *shape, c *c2d, w, h float64) {
		dx := clamp(s.style.num("dx", 0.5), 0, w)
		dy := clamp(s.style.num("dy", 0.5), 0, h)
		r := clamp(s.style.num("size", 10), 0, h)
		r = min3((h-dy)/2, w/2, r)
		dx = math.Max(r+dy*0.5, dx)
		dx = math.Min(w-r-dy*0.5, dx)
		c.begin()
		c.moveTo(dx-dy*0.5, h-dy)
		c.lineTo(r, h-dy)
		c.arcTo(r, r, 0, false, true, 0, h-dy-r)
		c.lineTo(0, r)
		c.arcTo(r, r, 0, false, true, r, 0)
		c.lineTo(w-r, 0)
		c.arcTo(r, r, 0, false, true, w, r)
		c.lineTo(w, h-dy-r)
		c.arcTo(r, r, 0, false, true, w-r, h-dy)
		c.lineTo(dx+dy*0.5, h-dy)
		c.arcTo(1.9*dy, 1.4*dy, 0, false, true, dx-dy, h)
		c.arcTo(0.9*dy, 1.4*dy, 0, false, false, dx-dy*0.5, h-dy)
		c.close()
		c.fillAndStroke()
	})
	d.labelMargins = basicCalloutLabelMargins
	return d
}()

// basicWaveShape ports mxShapeBasicWave ("wave2"): a band with wavy top
// and bottom (dy is the wave height relative to the height).
var basicWaveShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dy := h * clamp(s.style.num("dy", 0.5), 0, h)
	const fy = 1.4
	c.begin()
	c.moveTo(0, dy/2)
	c.quadTo(w/6, dy*(1-fy), w/3, dy/2)
	c.quadTo(w/2, dy*fy, w*2/3, dy/2)
	c.quadTo(w*5/6, dy*(1-fy), w, dy/2)
	c.lineTo(w, h-dy/2)
	c.quadTo(w*5/6, h-dy*fy, w*2/3, h-dy/2)
	c.quadTo(w/2, h-dy*(1-fy), w/3, h-dy/2)
	c.quadTo(w/6, h-dy*fy, 0, h-dy/2)
	c.close()
	c.fillAndStroke()
})

// basicOctagonShape ports mxShapeBasicOctagon ("octagon2"): corners cut
// by 2*dx.
var basicOctagonShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w) * 2
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(dx, 0)
	c.lineTo(w-dx, 0)
	c.lineTo(w, dx)
	c.lineTo(w, h-dx)
	c.lineTo(w-dx, h)
	c.lineTo(dx, h)
	c.lineTo(0, h-dx)
	c.lineTo(0, dx)
	c.close()
	c.fillAndStroke()
})

// basicTriangleAcuteShape ports mxShapeBasicTriangleAcute: the top at dx
// (relative to the width).
var basicTriangleAcuteShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := w * clamp(s.style.num("dx", 0.5), 0, w)
	c.begin()
	c.moveTo(0, h)
	c.lineTo(dx, 0)
	c.lineTo(w, h)
	c.close()
	c.fillAndStroke()
})

// basicTriangleObtuseShape ports mxShapeBasicTriangleObtuse: the bottom
// left corner at dx (relative to the width).
var basicTriangleObtuseShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := w * clamp(s.style.num("dx", 0.5), 0, w)
	c.begin()
	c.moveTo(dx, h)
	c.lineTo(0, 0)
	c.lineTo(w, h)
	c.close()
	c.fillAndStroke()
})

// basicDropShape ports mxShapeBasicDrop: a circle at the bottom with a
// tip at the top center.
var basicDropShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	r := math.Min(h, w) * 0.5
	d := h - r
	a := math.Sqrt(d*d - r*r)
	angle := math.Atan(a / r)
	x1 := r * math.Sin(angle)
	y1 := r * math.Cos(angle)
	c.begin()
	c.moveTo(w*0.5, 0)
	c.lineTo(w*0.5+x1, h-r-y1)
	c.arcTo(r, r, 0, false, true, w*0.5+r, h-r)
	c.arcTo(r, r, 0, false, true, w*0.5, h)
	c.arcTo(r, r, 0, false, true, w*0.5-r, h-r)
	c.arcTo(r, r, 0, false, true, w*0.5-x1, h-r-y1)
	c.close()
	c.fillAndStroke()
})

// basicCone2Shape ports mxShapeBasicCone2: the tip at dx (relative to the
// width), the base an ellipse from dy (relative to the height) down.
var basicCone2Shape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := w * clamp(s.style.num("dx", 0.5), 0, w)
	dy := h * clamp(s.style.num("dy", 0.9), 0, h)
	ry := h - dy
	c.begin()
	c.moveTo(dx, 0)
	if ry > 0 {
		c.lineTo(w, h-ry)
		c.arcTo(w*0.5, ry, 0, false, true, w*0.5, h)
		c.arcTo(w*0.5, ry, 0, false, true, 0, h-ry)
	} else {
		c.lineTo(w, h)
		c.lineTo(0, h)
	}
	c.close()
	c.fillAndStroke()
})

// basicPyramidShape ports mxShapeBasicPyramid: a quadrilateral with a
// front edge from the top (dx1) to the bottom (dx2); dy1 and dy2 place
// the left and right corners (all relative).
var basicPyramidShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx1 := w * clamp(s.style.num("dx1", 0.5), 0, w)
	dx2 := w * clamp(s.style.num("dx2", 0.6), 0, w)
	dy1 := h * clamp(s.style.num("dy1", 0.9), 0, h)
	dy2 := h * clamp(s.style.num("dy2", 0.8), 0, h)
	c.begin()
	c.moveTo(dx1, 0)
	c.lineTo(w, dy2)
	c.lineTo(dx2, h)
	c.lineTo(0, dy1)
	c.close()
	c.fillAndStroke()
	c.setShadow(false)
	c.begin()
	c.moveTo(dx1, 0)
	c.lineTo(dx2, h)
	c.stroke()
})

// basic4PointStar2Shape ports mxShapeBasic4PointStar2: the inner corners
// at dx/2 (relative) from the sides.
var basic4PointStar2Shape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := 0.5 * clamp(s.style.num("dx", 0.8), 0, w)
	c.begin()
	c.moveTo(0, h*0.5)
	c.lineTo(dx*w, dx*h)
	c.lineTo(w*0.5, 0)
	c.lineTo(w-dx*w, dx*h)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx*w, h-dx*h)
	c.lineTo(w*0.5, h)
	c.lineTo(dx*w, h-dx*h)
	c.close()
	c.fillAndStroke()
})

// basicDiagSnipRectShape ports mxShapeBasicDiagSnipRect: the top left and
// bottom right corners cut by 2*dx.
var basicDiagSnipRectShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w) * 2
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(dx, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h-dx)
	c.lineTo(w-dx, h)
	c.lineTo(0, h)
	c.lineTo(0, dx)
	c.close()
	c.fillAndStroke()
})

// basicDiagRoundRectShape ports mxShapeBasicDiagRoundRect: the top left
// and bottom right corners rounded with radius 2*dx.
var basicDiagRoundRectShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w) * 2
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(dx, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h-dx)
	c.arcTo(dx, dx, 0, false, true, w-dx, h)
	c.lineTo(0, h)
	c.lineTo(0, dx)
	c.arcTo(dx, dx, 0, false, true, dx, 0)
	c.close()
	c.fillAndStroke()
})

// basicCornerRoundRectShape ports mxShapeBasicCornerRoundRect: the top
// left corner rounded with radius 2*dx.
var basicCornerRoundRectShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w) * 2
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(dx, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.lineTo(0, dx)
	c.arcTo(dx, dx, 0, false, true, dx, 0)
	c.close()
	c.fillAndStroke()
})

// basicPlaqueShape ports mxShapeBasicPlaque: a rectangle with concave
// corners of radius 2*dx.
var basicPlaqueShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w) * 2
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(w-dx, 0)
	c.arcTo(dx, dx, 0, false, false, w, dx)
	c.lineTo(w, h-dx)
	c.arcTo(dx, dx, 0, false, false, w-dx, h)
	c.lineTo(dx, h)
	c.arcTo(dx, dx, 0, false, false, 0, h-dx)
	c.lineTo(0, dx)
	c.arcTo(dx, dx, 0, false, false, dx, 0)
	c.close()
	c.fillAndStroke()
})

// basicFrameShape ports mxShapeBasicFrame: a rectangle with a hole, dx
// wide (the inner path runs the other way round).
var basicFrameShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.lineTo(0, 0)
	c.close()
	c.moveTo(dx, dx)
	c.lineTo(dx, h-dx)
	c.lineTo(w-dx, h-dx)
	c.lineTo(w-dx, dx)
	c.close()
	c.fillAndStroke()
})

// basicPlaqueFrameShape ports mxShapeBasicPlaqueFrame: a plaque with a
// hole of the same form, dx wide.
var basicPlaqueFrameShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.25, h*0.25, dx)
	c.begin()
	c.moveTo(w-dx, 0)
	c.arcTo(dx, dx, 0, false, false, w, dx)
	c.lineTo(w, h-dx)
	c.arcTo(dx, dx, 0, false, false, w-dx, h)
	c.lineTo(dx, h)
	c.arcTo(dx, dx, 0, false, false, 0, h-dx)
	c.lineTo(0, dx)
	c.arcTo(dx, dx, 0, false, false, dx, 0)
	c.close()
	c.moveTo(dx*2, dx)
	c.arcTo(dx*2, dx*2, 0, false, true, dx, dx*2)
	c.lineTo(dx, h-2*dx)
	c.arcTo(dx*2, dx*2, 0, false, true, dx*2, h-dx)
	c.lineTo(w-2*dx, h-dx)
	c.arcTo(dx*2, dx*2, 0, false, true, w-dx, h-2*dx)
	c.lineTo(w-dx, dx*2)
	c.arcTo(dx*2, dx*2, 0, false, true, w-2*dx, dx)
	c.close()
	c.fillAndStroke()
})

// basicRoundedFrameShape ports mxShapeBasicRoundedFrame: a rounded
// rectangle with a rounded hole, dx wide.
var basicRoundedFrameShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.25, h*0.25, dx)
	c.begin()
	c.moveTo(w-2*dx, 0)
	c.arcTo(dx*2, dx*2, 0, false, true, w, 2*dx)
	c.lineTo(w, h-2*dx)
	c.arcTo(dx*2, dx*2, 0, false, true, w-2*dx, h)
	c.lineTo(dx*2, h)
	c.arcTo(dx*2, dx*2, 0, false, true, 0, h-2*dx)
	c.lineTo(0, 2*dx)
	c.arcTo(dx*2, dx*2, 0, false, true, 2*dx, 0)
	c.close()
	c.moveTo(dx*2, dx)
	c.arcTo(dx, dx, 0, false, false, dx, dx*2)
	c.lineTo(dx, h-2*dx)
	c.arcTo(dx, dx, 0, false, false, dx*2, h-dx)
	c.lineTo(w-2*dx, h-dx)
	c.arcTo(dx, dx, 0, false, false, w-dx, h-2*dx)
	c.lineTo(w-dx, dx*2)
	c.arcTo(dx, dx, 0, false, false, w-2*dx, dx)
	c.close()
	c.fillAndStroke()
})

// basicFrameCornerShape ports mxShapeBasicFrameCorner: the top left
// corner of a frame, dx wide.
var basicFrameCornerShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w-dx, dx)
	c.lineTo(dx, dx)
	c.lineTo(dx, h-dx)
	c.lineTo(0, h)
	c.close()
	c.fillAndStroke()
})

// basicDiagStripeShape ports mxShapeBasicDiagStripe: a stripe along the
// diagonal from the bottom left to the top right, whose width grows with
// dx*100 over the side.
var basicDiagStripeShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w, h, dx)
	c.begin()
	c.moveTo(0, h)
	c.lineTo(w, 0)
	c.lineTo(w, math.Min(dx*100/w, h))
	c.lineTo(math.Min(dx*100/h, w), h)
	c.close()
	c.fillAndStroke()
})

// basicDonutShape ports mxShapeBasicDonut: an ellipse with an elliptic
// hole, dx wide.
var basicDonutShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(0, h*0.5)
	c.arcTo(w*0.5, h*0.5, 0, false, true, w*0.5, 0)
	c.arcTo(w*0.5, h*0.5, 0, false, true, w, h*0.5)
	c.arcTo(w*0.5, h*0.5, 0, false, true, w*0.5, h)
	c.arcTo(w*0.5, h*0.5, 0, false, true, 0, h*0.5)
	c.close()
	c.moveTo(w*0.5, dx)
	c.arcTo(w*0.5-dx, h*0.5-dx, 0, false, false, dx, h*0.5)
	c.arcTo(w*0.5-dx, h*0.5-dx, 0, false, false, w*0.5, h-dx)
	c.arcTo(w*0.5-dx, h*0.5-dx, 0, false, false, w-dx, h*0.5)
	c.arcTo(w*0.5-dx, h*0.5-dx, 0, false, false, w*0.5, dx)
	c.close()
	c.fillAndStroke()
})

// basicLayeredRectDx is the layer distance of mxShapeBasicLayeredRect.
func basicLayeredRectDx(s *shape, w, h float64) float64 {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	return min3(w*0.5, h*0.5, dx)
}

// basicLayeredRectShape ports mxShapeBasicLayeredRect: three rectangles
// dx/2 apart, the front one at the top left; boundedLbl confines the
// label to it.
var basicLayeredRectShape = func() *shapeDef {
	d := basicShape(func(s *shape, c *c2d, w, h float64) {
		dx := basicLayeredRectDx(s, w, h)
		c.begin()
		c.moveTo(dx, dx)
		c.lineTo(w, dx)
		c.lineTo(w, h)
		c.lineTo(dx, h)
		c.close()
		c.fillAndStroke()
		c.begin()
		c.moveTo(dx*0.5, dx*0.5)
		c.lineTo(w-dx*0.5, dx*0.5)
		c.lineTo(w-dx*0.5, h-dx*0.5)
		c.lineTo(dx*0.5, h-dx*0.5)
		c.close()
		c.fillAndStroke()
		c.begin()
		c.moveTo(0, 0)
		c.lineTo(w-dx, 0)
		c.lineTo(w-dx, h-dx)
		c.lineTo(0, h-dx)
		c.close()
		c.fillAndStroke()
	})
	d.labelMargins = func(s *shape, r rect) *rect {
		if boundedLbl(s) {
			dx := basicLayeredRectDx(s, r.w, r.h)
			return &rect{0, 0, dx, dx}
		}
		return nil
	}
	return d
}()

// basicButtonShape ports mxShapeBasicButton: a rectangle with a bevel dx
// wide, each side a filled and stroked trapezoid (the left one twice).
var basicButtonShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.close()
	c.fillAndStroke()
	c.setShadow(false)
	c.setLineJoin("round")
	c.begin()
	c.moveTo(0, h)
	c.lineTo(0, 0)
	c.lineTo(dx, dx)
	c.lineTo(dx, h-dx)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w-dx, dx)
	c.lineTo(dx, dx)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(w-dx, h-dx)
	c.lineTo(w-dx, dx)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(0, h)
	c.lineTo(dx, h-dx)
	c.lineTo(w-dx, h-dx)
	c.lineTo(w, h)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(0, h)
	c.lineTo(0, 0)
	c.lineTo(dx, dx)
	c.lineTo(dx, h-dx)
	c.close()
	c.fillAndStroke()
})

// basicShadedButtonShape ports mxShapeBasicShadedButton: an unstroked
// rectangle with a bevel dx wide, lit from the top left (white and black
// at an alpha that replaces the opacity).
var basicShadedButtonShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	c.setShadow(false)
	dx := clamp(s.style.num("dx", 0.5), 0, w)
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.close()
	c.fill()
	c.setFillColor("#ffffff")
	c.setAlpha(0.25)
	c.begin()
	c.moveTo(0, h)
	c.lineTo(0, 0)
	c.lineTo(dx, dx)
	c.lineTo(dx, h-dx)
	c.close()
	c.fill()
	c.setAlpha(0.5)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w-dx, dx)
	c.lineTo(dx, dx)
	c.close()
	c.fill()
	c.setFillColor("#000000")
	c.setAlpha(0.25)
	c.begin()
	c.moveTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(w-dx, h-dx)
	c.lineTo(w-dx, dx)
	c.close()
	c.fill()
	c.setAlpha(0.5)
	c.begin()
	c.moveTo(0, h)
	c.lineTo(dx, h-dx)
	c.lineTo(w-dx, h-dx)
	c.lineTo(w, h)
	c.close()
	c.fill()
})

// basicPieShape ports mxShapeBasicPie: a sector of the ellipse from
// startAngle to endAngle clockwise (fractions of a turn from the top).
var basicPieShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	startAngleSource := clamp(s.style.num("startAngle", 0.25), 0, 1)
	endAngleSource := clamp(s.style.num("endAngle", 0.75), 0, 1)
	startAngle := 2 * math.Pi * startAngleSource
	endAngle := 2 * math.Pi * endAngleSource
	rx := w * 0.5
	ry := h * 0.5
	startX := rx + math.Sin(startAngle)*rx
	startY := ry - math.Cos(startAngle)*ry
	endX := rx + math.Sin(endAngle)*rx
	endY := ry - math.Cos(endAngle)*ry
	angDiff := endAngle - startAngle
	if angDiff < 0 {
		angDiff += math.Pi * 2
	}
	bigArc := angDiff >= math.Pi
	c.begin()
	startAngleDiff := math.Mod(startAngleSource, 1)
	endAngleDiff := math.Mod(endAngleSource, 1)
	switch {
	case startAngleDiff == 0 && endAngleDiff == 0.5:
		c.moveTo(rx, ry)
		c.lineTo(startX, startY)
		c.arcTo(rx, ry, 0, false, true, w, h*0.5)
		c.arcTo(rx, ry, 0, false, true, w*0.5, h)
	case startAngleDiff == 0.5 && endAngleDiff == 0:
		c.moveTo(rx, ry)
		c.lineTo(startX, startY)
		c.arcTo(rx, ry, 0, false, true, 0, h*0.5)
		c.arcTo(rx, ry, 0, false, true, w*0.5, 0)
	default:
		c.moveTo(rx, ry)
		c.lineTo(startX, startY)
		c.arcTo(rx, ry, 0, bigArc, true, endX, endY)
	}
	c.close()
	c.fillAndStroke()
})

// basicArcShape ports mxShapeBasicArc: an elliptic arc from startAngle to
// endAngle clockwise, with the startArrow and endArrow markers of edges;
// the ellipse is inset so that the markers stay in the bounds.
var basicArcShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	startAngleSource := clamp(s.style.num("startAngle", 0.25), 0, 1)
	endAngleSource := clamp(s.style.num("endAngle", 0.75), 0, 1)
	startAngle := 2 * math.Pi * startAngleSource
	endAngle := 2 * math.Pi * endAngleSource

	// arrow markers
	startArrow := s.style.get("startArrow", "none")
	endArrow := s.style.get("endArrow", "none")
	startSize := s.style.num("startSize", defaultMarkerSize)
	endSize := s.style.num("endSize", defaultMarkerSize)

	// inset ellipse to keep markers within shape bounds
	markerPadding := 0.0
	if startArrow != "none" {
		markerPadding = math.Max(markerPadding, startSize)
	}
	if endArrow != "none" {
		markerPadding = math.Max(markerPadding, endSize)
	}
	if markerPadding > 0 {
		markerPadding = (markerPadding + s.strokewidth) / 2
	}
	cx := w * 0.5
	cy := h * 0.5
	rx := math.Max(0, w*0.5-markerPadding)
	ry := math.Max(0, h*0.5-markerPadding)
	startX := cx + math.Sin(startAngle)*rx
	startY := cy - math.Cos(startAngle)*ry
	endX := cx + math.Sin(endAngle)*rx
	endY := cy - math.Cos(endAngle)*ry

	var startMarker, endMarker func()
	if startArrow != "none" {
		startFilled := styleBool(s.style, "startFill", true)
		// direction pointing into the start (opposite of clockwise tangent)
		dx := -rx * math.Cos(startAngle)
		dy := -ry * math.Sin(startAngle)
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			pe := &point{startX, startY}
			startMarker = basicCreateMarker(c, s, startArrow, pe, dx/dist, dy/dist, startSize, true, s.strokewidth, startFilled)
			startX, startY = pe.x, pe.y
		}
	}
	if endArrow != "none" {
		endFilled := styleBool(s.style, "endFill", true)
		// direction along the clockwise tangent at end
		dx := rx * math.Cos(endAngle)
		dy := ry * math.Sin(endAngle)
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			pe := &point{endX, endY}
			endMarker = basicCreateMarker(c, s, endArrow, pe, dx/dist, dy/dist, endSize, false, s.strokewidth, endFilled)
			endX, endY = pe.x, pe.y
		}
	}

	// draw arc
	angDiff := endAngle - startAngle
	if angDiff < 0 {
		angDiff += math.Pi * 2
	}
	bigArc := angDiff > math.Pi
	c.begin()
	startAngleDiff := math.Mod(startAngleSource, 1)
	endAngleDiff := math.Mod(endAngleSource, 1)
	switch {
	case startMarker == nil && endMarker == nil && startAngleDiff == 0 && endAngleDiff == 0.5:
		c.moveTo(startX, startY)
		c.arcTo(rx, ry, 0, false, true, cx+rx, cy)
		c.arcTo(rx, ry, 0, false, true, cx, cy+ry)
	case startMarker == nil && endMarker == nil && startAngleDiff == 0.5 && endAngleDiff == 0:
		c.moveTo(startX, startY)
		c.arcTo(rx, ry, 0, false, true, cx-rx, cy)
		c.arcTo(rx, ry, 0, false, true, cx, cy-ry)
	default:
		c.moveTo(startX, startY)
		c.arcTo(rx, ry, 0, bigArc, true, endX, endY)
	}
	c.stroke()

	// paint markers
	c.setShadow(false)
	c.setDashed(false, false)
	if startMarker != nil {
		c.setFillColor(s.style.get("startFillColor", s.stroke))
		startMarker()
	}
	if endMarker != nil {
		c.setFillColor(s.style.get("endFillColor", s.stroke))
		endMarker()
	}
})

// basicCreateMarker is mxMarker.createMarker: nil for types it does not
// know, which draw.io does not draw either (connectors here draw them as
// classic).
func basicCreateMarker(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
	f := markers[typ]
	if f == nil {
		s.conv.warnOnce("marker:"+typ, "arrow %q is not supported; not drawn", typ)
		return nil
	}
	return f(c, s, typ, pe, unitX, unitY, size, source, sw, filled)
}

// basicPartConcEllipseShape ports mxShapeBasicPartConcEllipse: a part of
// a ring from startAngle to endAngle clockwise, arcWidth (relative to the
// radius) wide.
var basicPartConcEllipseShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	startAngle := 2 * math.Pi * clamp(s.style.num("startAngle", 0.25), 0, 1)
	endAngle := 2 * math.Pi * clamp(s.style.num("endAngle", 0.75), 0, 1)
	arcWidth := 1 - clamp(s.style.num("arcWidth", 0.5), 0, 1)
	rx := w * 0.5
	ry := h * 0.5
	rx2 := rx * arcWidth
	ry2 := ry * arcWidth
	angDiff := endAngle - startAngle
	if angDiff < 0 {
		angDiff += math.Pi * 2
	} else if angDiff == math.Pi {
		endAngle += 0.00001
	}
	startX := rx + math.Sin(startAngle)*rx
	startY := ry - math.Cos(startAngle)*ry
	innerStartX := rx + math.Sin(startAngle)*rx2
	innerStartY := ry - math.Cos(startAngle)*ry2
	endX := rx + math.Sin(endAngle)*rx
	endY := ry - math.Cos(endAngle)*ry
	innerEndX := rx + math.Sin(endAngle)*rx2
	innerEndY := ry - math.Cos(endAngle)*ry2
	bigArc := angDiff >= math.Pi
	c.begin()
	c.moveTo(startX, startY)
	c.arcTo(rx, ry, 0, bigArc, true, endX, endY)
	c.lineTo(innerEndX, innerEndY)
	c.arcTo(rx2, ry2, 0, bigArc, false, innerStartX, innerStartY)
	c.close()
	c.fillAndStroke()
})

// basicNumEntryVertShape ports mxShapeBasicNumEntryVert: a circle of
// diameter dy set into the top of a rectangle.
var basicNumEntryVertShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dy := clamp(s.style.num("dy", 0.5), 0, w)
	const inset = 5
	d := min3(dy, w-2*inset, h-inset)
	// SVG does not render ellipses with radii of 0 or less
	if d > 0 {
		c.ellipse(w*0.5-d*0.5, 0, d, d)
		c.fillAndStroke()
	}
	c.begin()
	c.moveTo(0, d*0.5)
	c.lineTo(w*0.5-d*0.5-inset, d*0.5)
	c.arcTo(d*0.5+inset, d*0.5+inset, 0, false, false, w*0.5+d*0.5+inset, d*0.5)
	c.lineTo(w, d*0.5)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.close()
	c.fillAndStroke()
})

// basicBendingArchShape ports mxShapeBasicBendingArch: a part of a ring
// (as partConcEllipse) around an ellipse 5 inside it.
var basicBendingArchShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	startAngle := 2 * math.Pi * clamp(s.style.num("startAngle", 0.25), 0, 1)
	endAngle := 2 * math.Pi * clamp(s.style.num("endAngle", 0.75), 0, 1)
	arcWidth := 1 - clamp(s.style.num("arcWidth", 0.5), 0, 1)
	rx := w * 0.5
	ry := h * 0.5
	rx2 := rx * arcWidth
	ry2 := ry * arcWidth
	startX := rx + math.Sin(startAngle)*rx
	startY := ry - math.Cos(startAngle)*ry
	innerStartX := rx + math.Sin(startAngle)*rx2
	innerStartY := ry - math.Cos(startAngle)*ry2
	endX := rx + math.Sin(endAngle)*rx
	endY := ry - math.Cos(endAngle)*ry
	innerEndX := rx + math.Sin(endAngle)*rx2
	innerEndY := ry - math.Cos(endAngle)*ry2
	angDiff := endAngle - startAngle
	if angDiff < 0 {
		angDiff += math.Pi * 2
	}
	bigArc := angDiff > math.Pi
	rx3 := rx2 - 5
	ry3 := ry2 - 5
	// SVG does not render ellipses with radii of 0 or less
	if rx3 > 0 && ry3 > 0 {
		c.ellipse(w*0.5-rx3, h*0.5-ry3, 2*rx3, 2*ry3)
		c.fillAndStroke()
	}
	c.begin()
	c.moveTo(startX, startY)
	c.arcTo(rx, ry, 0, bigArc, true, endX, endY)
	c.lineTo(innerEndX, innerEndY)
	c.arcTo(rx2, ry2, 0, bigArc, false, innerStartX, innerStartY)
	c.close()
	c.fillAndStroke()
})

// basicThreeCornerRoundRectShape ports mxShapeBasicThreeCornerRoundRect:
// all corners but the bottom left rounded with radius 2*dx.
var basicThreeCornerRoundRectShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	dx := clamp(s.style.num("dx", 0.5), 0, w) * 2
	dx = min3(w*0.5, h*0.5, dx)
	c.begin()
	c.moveTo(dx, 0)
	c.lineTo(w-dx, 0)
	c.arcTo(dx, dx, 0, false, true, w, dx)
	c.lineTo(w, h-dx)
	c.arcTo(dx, dx, 0, false, true, w-dx, h)
	c.lineTo(0, h)
	c.lineTo(0, dx)
	c.arcTo(dx, dx, 0, false, true, dx, 0)
	c.close()
	c.fillAndStroke()
})

// basicPolygonShape ports mxShapeBasicPolygon: the points of polyCoords
// (a JSON array of [x, y] relative to the size), joined by lines or, where
// polyCurves has ["Q", x, y] for the segment, quadratic curves; closed
// unless polyline is set. Where JavaScript throws (invalid JSON, a point
// that is not an array) nothing is drawn; so it is here where it would
// compute NaN (coordinates that are not numbers), which draw.io writes
// into the SVG path.
var basicPolygonShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	coords, ok := parseJSONArray(s.style.get("polyCoords", "[]"))
	if !ok {
		return
	}
	curves, ok := parseJSONArray(s.style.get("polyCurves", "[]"))
	if !ok {
		return
	}
	// polyline == false, as JavaScript compares
	polyline := styleBool(s.style, "polyline", false)
	n := len(coords)
	if n == 0 {
		return
	}
	pts := make([]point, n)
	for i, p := range coords {
		x, okx := jsonIndexNumber(p, 0)
		y, oky := jsonIndexNumber(p, 1)
		if !okx || !oky {
			return
		}
		pts[i] = point{x * w, y * h}
	}
	// ctrl[i] is the control point of the curve to pts[i] (to pts[0] for
	// i == n, the closing segment), nil for a line
	ctrl := make([]*point, n+1)
	for i := 1; i <= n; i++ {
		if i == n && polyline {
			break
		}
		if i-1 >= len(curves) {
			continue
		}
		ci, ok := curves[i-1].([]any)
		if !ok || len(ci) < 3 || ci[0] != "Q" {
			continue
		}
		x, okx := jsonNumber(ci[1])
		y, oky := jsonNumber(ci[2])
		if !okx || !oky {
			return
		}
		ctrl[i] = &point{x * w, y * h}
	}
	c.begin()
	c.moveTo(pts[0].x, pts[0].y)
	for i := 1; i < n; i++ {
		if cp := ctrl[i]; cp != nil {
			c.quadTo(cp.x, cp.y, pts[i].x, pts[i].y)
		} else {
			c.lineTo(pts[i].x, pts[i].y)
		}
	}
	if !polyline {
		if cp := ctrl[n]; cp != nil {
			c.quadTo(cp.x, cp.y, pts[0].x, pts[0].y)
		}
		c.close()
	}
	c.end()
	c.fillAndStroke()
})

// parseJSONArray parses a style value with JSON.parse; a value that is
// valid JSON but no array is an empty array (its length is undefined).
func parseJSONArray(v string) ([]any, bool) {
	var x any
	if err := json.Unmarshal([]byte(v), &x); err != nil {
		return nil, false
	}
	a, _ := x.([]any)
	return a, true
}

// jsonIndexNumber is v[i] of a JSON array as a number (see jsonNumber).
func jsonIndexNumber(v any, i int) (float64, bool) {
	a, ok := v.([]any)
	if !ok || i >= len(a) {
		return 0, false
	}
	return jsonNumber(a[i])
}

// jsonNumber converts a JSON value to a number as JavaScript's arithmetic
// does (null is 0, strings are parsed, arrays are joined to a string
// first); ok is false where that is NaN.
func jsonNumber(v any) (float64, bool) {
	switch v := v.(type) {
	case []any:
		switch len(v) {
		case 0:
			return 0, true
		case 1:
			if _, isBool := v[0].(bool); !isBool { // "true" is NaN
				return jsonNumber(v[0])
			}
		}
	case float64:
		return v, true
	case nil:
		return 0, true
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case string:
		t := strings.TrimSpace(v)
		if t == "" {
			return 0, true
		}
		f, err := strconv.ParseFloat(t, 64)
		return f, err == nil && finite(f)
	}
	return 0, false
}

// maxPatternLines bounds the lines of a pattern fill (steps of 0 or less
// never end in JavaScript; absurd sizes nearly never do).
const maxPatternLines = 100000

// patternSteps reports whether stepping from 0 up to limit by step ends
// within maxPatternLines steps.
func patternSteps(limit, step float64) bool {
	return step > 0 && limit/step < maxPatternLines
}

// basicPatternFillRectShape ports mxShapeBasicPatternFillRect: a filled
// rectangle hatched with lines step apart in fillStrokeColor and
// fillStrokeWidth (fillStyle=diag, diagRev, diagGrid, vert, hor or grid),
// whose sides are stroked where top, right, bottom and left are 1.
// draw.io's own fill patterns do not know these fillStyle values and fill
// the rectangle solid.
var basicPatternFillRectShape = basicShape(func(s *shape, c *c2d, w, h float64) {
	strokeColor := s.style.get("strokeColor", "#000000")
	strokeWidth := s.style.num("strokeWidth", 1)
	c.rect(0, 0, w, h)
	c.fill()
	fillStrokeColor := s.style.get("fillStrokeColor", "#cccccc")
	fillStrokeWidth := s.style.num("fillStrokeWidth", 1)
	c.setStrokeColor(fillStrokeColor)
	c.setStrokeWidth(fillStrokeWidth)
	step := s.style.num("step", 5)
	fillStyle := s.style.get("fillStyle", "none")
	if fillStyle == "diag" || fillStyle == "diagGrid" {
		step *= 1.41
		c.begin()
		for i := 0.0; i < h+w && patternSteps(h+w, step); i += step {
			if i <= h {
				c.moveTo(0, i)
				if i <= w {
					c.lineTo(i, 0)
				} else {
					c.lineTo(w, i-w)
				}
			} else {
				c.moveTo(i-h, h)
				if i <= w {
					c.lineTo(i, 0)
				} else {
					c.lineTo(w, i-w)
				}
			}
		}
		c.stroke()
	} else if fillStyle == "vert" || fillStyle == "grid" {
		c.begin()
		for i := 0.0; i <= w && patternSteps(w, step); i += step {
			c.moveTo(i, 0)
			c.lineTo(i, h)
		}
		c.stroke()
	}
	if fillStyle == "diagRev" || fillStyle == "diagGrid" {
		if fillStyle == "diagRev" {
			step *= 1.41
		}
		c.begin()
		for i := 0.0; i < h+w && patternSteps(h+w, step); i += step {
			if i <= h {
				c.moveTo(w, i)
				if i <= w {
					c.lineTo(w-i, 0)
				} else {
					c.lineTo(w-w, i-w)
				}
			} else {
				c.moveTo(w-i+h, h)
				if i <= w {
					c.lineTo(w-i, 0)
				} else {
					c.lineTo(0, i-w)
				}
			}
		}
		c.stroke()
	} else if fillStyle == "hor" || fillStyle == "grid" {
		c.begin()
		for i := 0.0; i <= h && patternSteps(h, step); i += step {
			c.moveTo(0, i)
			c.lineTo(w, i)
		}
		c.stroke()
	}
	c.setStrokeColor(strokeColor)
	c.setStrokeWidth(strokeWidth)
	c.begin()
	c.moveTo(0, 0)
	if s.style.get("top", "1") == "1" {
		c.lineTo(w, 0)
	} else {
		c.moveTo(w, 0)
	}
	if s.style.get("right", "1") == "1" {
		c.lineTo(w, h)
	} else {
		c.moveTo(w, h)
	}
	if s.style.get("bottom", "1") == "1" {
		c.lineTo(0, h)
	} else {
		c.moveTo(0, h)
	}
	if s.style.get("left", "1") == "1" {
		c.lineTo(0, 0)
	}
	c.end()
	c.stroke()
})

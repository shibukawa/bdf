package drawio

import "math"

// The arrows of draw.io's shapes/mxArrows.js (mxgraph.arrows2.*): block
// arrows drawn as mxActor paths in the bounds, and three wedge-shaped edge
// shapes. The sizes are read as the JavaScript reads them: most are
// absolute and clamped to the width or height, dy of the straight arrows
// is a fraction of half the height, and the constructors' defaults (0.5
// for dx and dy, and none at all for some keys) apply when the style has
// no value.

func init() {
	registerShape("mxgraph.arrows2.arrow", arrows2Arrow)
	registerShape("mxgraph.arrows2.twoWayArrow", arrows2TwoWayArrow)
	registerShape("mxgraph.arrows2.stylisedArrow", arrows2StylisedArrow)
	registerShape("mxgraph.arrows2.sharpArrow", arrows2SharpArrow)
	registerShape("mxgraph.arrows2.sharpArrow2", arrows2SharpArrow2)
	registerShape("mxgraph.arrows2.calloutArrow", arrows2CalloutArrow)
	registerShape("mxgraph.arrows2.bendArrow", arrows2BendArrow)
	registerShape("mxgraph.arrows2.bendDoubleArrow", arrows2BendDoubleArrow)
	registerShape("mxgraph.arrows2.calloutDoubleArrow", arrows2CalloutDoubleArrow)
	registerShape("mxgraph.arrows2.calloutQuadArrow", arrows2CalloutQuadArrow)
	registerShape("mxgraph.arrows2.calloutDouble90Arrow", arrows2CalloutDouble90Arrow)
	registerShape("mxgraph.arrows2.quadArrow", arrows2QuadArrow)
	registerShape("mxgraph.arrows2.triadArrow", arrows2TriadArrow)
	registerShape("mxgraph.arrows2.tailedArrow", arrows2TailedArrow(false))
	registerShape("mxgraph.arrows2.tailedNotchedArrow", arrows2TailedArrow(true))
	registerShape("mxgraph.arrows2.stripedArrow", arrows2StripedArrow)
	registerShape("mxgraph.arrows2.jumpInArrow", arrows2JumpInArrow)
	registerShape("mxgraph.arrows2.uTurnArrow", arrows2UTurnArrow)

	registerShape("mxgraph.arrows2.wedgeArrow", &shapeDef{paintEdge: paintWedgeArrow, augmentBounds: wedgeArrowBounds})
	registerShape("mxgraph.arrows2.wedgeArrowDashed", &shapeDef{paintEdge: paintWedgeArrowDashed(false), augmentBounds: wedgeArrowBounds})
	registerShape("mxgraph.arrows2.wedgeArrowDashed2", &shapeDef{paintEdge: paintWedgeArrowDashed(true), augmentBounds: wedgeArrowBounds})
}

// jsStyleFloat is parseFloat(mxUtils.getValue(style, key, def)): NaN for a
// value that is not a number; def NaN stands for a default that is
// undefined (a property the shape's constructor does not set).
func jsStyleFloat(st style, key string, def float64) float64 {
	v, ok := st[key]
	if !ok {
		return def
	}
	if f, ok := parseFloat(v); ok {
		return f
	}
	return math.NaN()
}

// arrowsVal is Math.max(0, Math.min(max, parseFloat(mxUtils.getValue(style,
// key, def)))), which keeps NaN.
func arrowsVal(s *shape, key string, def, max float64) float64 {
	return clamp(jsStyleFloat(s.style, key, def), 0, max)
}

// undefined is the value of a style key without a default.
var undefined = math.NaN()

// nanCanvas is the canvas with SVG's handling of broken path data:
// mxSvgCanvas2D writes a size that is not a number as NaN into the path,
// and SVG draws a path only up to its first error. An arc with such a
// size has no curves (mxUtils.arcToCurves makes none).
type nanCanvas struct {
	*c2d
	broken bool
}

func (c *nanCanvas) begin() {
	c.broken = false
	c.c2d.begin()
}

// ok reports whether the path goes on with the coordinates vs.
func (c *nanCanvas) ok(vs ...float64) bool {
	for _, v := range vs {
		if !finite(v) {
			c.broken = true
		}
	}
	return !c.broken
}

func (c *nanCanvas) moveTo(x, y float64) {
	if c.ok(x, y) {
		c.c2d.moveTo(x, y)
	}
}

func (c *nanCanvas) lineTo(x, y float64) {
	if c.ok(x, y) {
		c.c2d.lineTo(x, y)
	}
}

func (c *nanCanvas) arcTo(rx, ry, angle float64, largeArc, sweep bool, x, y float64) {
	if c.broken {
		return
	}
	for _, v := range []float64{rx, ry, angle, x, y} {
		if !finite(v) {
			return
		}
	}
	c.c2d.arcTo(rx, ry, angle, largeArc, sweep, x, y)
}

func (c *nanCanvas) close() {
	if !c.broken {
		c.c2d.close()
	}
}

// arrows2Shape is an arrows2 shape: paint draws in the bounds after
// paintVertexShape's translation to them.
func arrows2Shape(paint func(s *shape, c *nanCanvas, w, h float64)) *shapeDef {
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		paint(s, &nanCanvas{c2d: c}, w, h)
	}}
}

// arrows2Arrow ports mxShapeArrows2Arrow: an arrow with a notched tail
// (dx is the head length, dy the shaft width as a fraction of the height,
// headCrossline and tailCrossline draw lines across the shaft's ends).
var arrows2Arrow = func() *shapeDef {
	d := arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
		dy := h * 0.5 * arrowsVal(s, "dy", 0.5, 1)
		dx := arrowsVal(s, "dx", 0.5, w)
		notch := arrowsVal(s, "notch", 0, w)
		headCrossline := styleBool(s.style, "headCrossline", false)
		tailCrossline := styleBool(s.style, "tailCrossline", false)
		c.begin()
		c.moveTo(0, dy)
		c.lineTo(w-dx, dy)
		c.lineTo(w-dx, 0)
		c.lineTo(w, h*0.5)
		c.lineTo(w-dx, h)
		c.lineTo(w-dx, h-dy)
		c.lineTo(0, h-dy)
		c.lineTo(notch, h*0.5)
		c.close()
		c.fillAndStroke()
		if headCrossline {
			c.begin()
			c.moveTo(w-dx, dy)
			c.lineTo(w-dx, h-dy)
			c.stroke()
		}
		if tailCrossline {
			c.begin()
			c.moveTo(notch, dy)
			c.lineTo(notch, h-dy)
			c.stroke()
		}
	})
	// boundedLbl keeps the label on the shaft, the direction flipped by
	// the flips
	d.labelBounds = func(s *shape, r rect) rect {
		if !styleBool(s.style, "boundedLbl", false) {
			return r
		}
		w, h := r.w, r.h
		direction := s.direction
		if direction == "" {
			direction = "east"
		}
		if styleBool(s.style, "flipH", false) {
			switch direction {
			case "west":
				direction = "east"
			case "east":
				direction = "west"
			}
		}
		if styleBool(s.style, "flipV", false) {
			switch direction {
			case "north":
				direction = "south"
			case "south":
				direction = "north"
			}
		}
		var dy, dx float64
		if direction == "north" || direction == "south" {
			dy = w * 0.5 * arrowsVal(s, "dy", 0.5, 1)
			dx = arrowsVal(s, "dx", 0.5, h)
		} else {
			dy = h * 0.5 * arrowsVal(s, "dy", 0.5, 1)
			dx = arrowsVal(s, "dx", 0.5, w)
		}
		switch direction {
		case "east":
			return rect{r.x, r.y + dy, w - dx, h - 2*dy}
		case "west":
			return rect{r.x + dx, r.y + dy, w - dx, h - 2*dy}
		case "north":
			return rect{r.x + dy, r.y + dx, w - 2*dy, h - dx}
		}
		return rect{r.x + dy, r.y, w - 2*dy, h - dx}
	}
	return d
}()

// arrows2TwoWayArrow ports mxShapeArrows2TwoWayArrow: an arrow with heads
// at both ends.
var arrows2TwoWayArrow = func() *shapeDef {
	d := arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
		dy := h * 0.5 * arrowsVal(s, "dy", 0.5, 1)
		dx := arrowsVal(s, "dx", 0.5, w)
		c.begin()
		c.moveTo(dx, dy)
		c.lineTo(w-dx, dy)
		c.lineTo(w-dx, 0)
		c.lineTo(w, h*0.5)
		c.lineTo(w-dx, h)
		c.lineTo(w-dx, h-dy)
		c.lineTo(dx, h-dy)
		c.lineTo(dx, h)
		c.lineTo(0, h*0.5)
		c.lineTo(dx, 0)
		c.close()
		c.fillAndStroke()
	})
	d.labelBounds = func(s *shape, r rect) rect {
		if !styleBool(s.style, "boundedLbl", false) {
			return r
		}
		w, h := r.w, r.h
		if s.direction == "north" || s.direction == "south" {
			dy := w * 0.5 * arrowsVal(s, "dy", 0.5, 1)
			dx := arrowsVal(s, "dx", 0.5, h)
			return rect{r.x + dy, r.y + dx, w - 2*dy, h - 2*dx}
		}
		dy := h * 0.5 * arrowsVal(s, "dy", 0.5, 1)
		dx := arrowsVal(s, "dx", 0.5, w)
		return rect{r.x + dx, r.y + dy, w - 2*dx, h - 2*dy}
	}
	return d
}()

// arrows2StylisedArrow ports mxShapeArrows2StylisedArrow: an arrow whose
// shaft narrows to the tail (feather) and whose head reaches back 10
// beyond dx.
var arrows2StylisedArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := h * 0.5 * arrowsVal(s, "dy", 0.5, 1)
	dx := arrowsVal(s, "dx", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	feather := h * 0.5 * arrowsVal(s, "feather", 0.5, 1)
	c.begin()
	c.moveTo(0, feather)
	c.lineTo(w-dx, dy)
	c.lineTo(w-dx-10, 0)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx-10, h)
	c.lineTo(w-dx, h-dy)
	c.lineTo(0, h-feather)
	c.lineTo(notch, h*0.5)
	c.close()
	c.fillAndStroke()
})

// arrows2SharpArrow ports mxShapeArrows2SharpArrow: an arrow with a swept
// back head (dx1 where the head meets the shaft, dx2 its barbs).
var arrows2SharpArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy1 := h * 0.5 * arrowsVal(s, "dy1", 0.5, 1)
	dx1 := arrowsVal(s, "dx1", 0.5, w)
	dx2 := arrowsVal(s, "dx2", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	dx1a := arrowsVal(s, "dx1", 0.5, w)
	dy1a := h * 0.5 * arrowsVal(s, "dy1", 0.5, h)
	x2 := 0.0
	if h != 0 {
		x2 = dx1a + dx2*dy1a*2/h
	}
	c.begin()
	c.moveTo(0, dy1)
	c.lineTo(w-dx1, dy1)
	c.lineTo(w-x2, 0)
	c.lineTo(w-dx2, 0)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx2, h)
	c.lineTo(w-x2, h)
	c.lineTo(w-dx1, h-dy1)
	c.lineTo(0, h-dy1)
	c.lineTo(notch, h*0.5)
	c.close()
	c.fillAndStroke()
})

// arrows2SharpArrow2 ports mxShapeArrows2SharpArrow2: sharpArrow with the
// head's back corner at dx3, dy3.
var arrows2SharpArrow2 = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy1 := h * 0.5 * arrowsVal(s, "dy1", 0.5, 1)
	dx1 := arrowsVal(s, "dx1", 0.5, w)
	dx2 := arrowsVal(s, "dx2", 0.5, w)
	dy3 := h * 0.5 * arrowsVal(s, "dy3", 0.5, 1)
	dx3 := arrowsVal(s, "dx3", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	c.begin()
	c.moveTo(0, dy1)
	c.lineTo(w-dx1, dy1)
	c.lineTo(w-dx3, dy3)
	c.lineTo(w-dx2, 0)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx2, h)
	c.lineTo(w-dx3, h-dy3)
	c.lineTo(w-dx1, h-dy1)
	c.lineTo(0, h-dy1)
	c.lineTo(notch, h*0.5)
	c.close()
	c.fillAndStroke()
})

// arrows2CalloutArrow ports mxShapeArrows2CalloutArrow: a box (notch wide)
// with an arrow to the right (dy is half the shaft width, arrowHead how
// far the head reaches beyond it).
var arrows2CalloutArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	arrowHead := arrowsVal(s, "arrowHead", 0, h)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(notch, 0)
	c.lineTo(notch, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy-arrowHead)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx, h*0.5+dy+arrowHead)
	c.lineTo(w-dx, h*0.5+dy)
	c.lineTo(notch, h*0.5+dy)
	c.lineTo(notch, h)
	c.lineTo(0, h)
	c.close()
	c.fillAndStroke()
})

// arrows2BendArrow ports mxShapeArrows2BendArrow: an arrow up from the
// bottom left that turns right (rounded=1 rounds the turn).
var arrows2BendArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	notch := arrowsVal(s, "notch", 0, h)
	arrowHead := arrowsVal(s, "arrowHead", 40, h)
	rounded := styleIsOne(s.style, "rounded", false)
	c.begin()
	c.moveTo(w-dx, 0)
	c.lineTo(w, arrowHead*0.5)
	c.lineTo(w-dx, arrowHead)
	c.lineTo(w-dx, arrowHead/2+dy)
	if rounded {
		c.lineTo(dy*2.2, arrowHead/2+dy)
		c.arcTo(dy*0.2, dy*0.2, 0, false, false, dy*2, arrowHead/2+dy*1.2)
	} else {
		c.lineTo(dy*2, arrowHead/2+dy)
	}
	c.lineTo(dy*2, h)
	c.lineTo(dy, h-notch)
	c.lineTo(0, h)
	if rounded {
		c.lineTo(0, arrowHead/2+dy)
		c.arcTo(dy*2, dy*2, 0, false, true, dy*2, arrowHead/2-dy)
	} else {
		c.lineTo(0, arrowHead/2-dy)
	}
	c.lineTo(w-dx, arrowHead/2-dy)
	c.close()
	c.fillAndStroke()
})

// arrows2BendDoubleArrow ports mxShapeArrows2BendDoubleArrow: a corner
// with arrow heads to the right and to the bottom.
var arrows2BendDoubleArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	arrowHead := arrowsVal(s, "arrowHead", 40, w)
	rounded := styleIsOne(s.style, "rounded", false)
	c.begin()
	c.moveTo(w-dx, 0)
	c.lineTo(w, arrowHead*0.5)
	c.lineTo(w-dx, arrowHead)
	c.lineTo(w-dx, arrowHead/2+dy)
	if rounded {
		c.lineTo(arrowHead/2+dy*1.2, arrowHead/2+dy)
		c.arcTo(dy*0.2, dy*0.2, 0, false, false, arrowHead/2+dy, arrowHead/2+dy*1.2)
	} else {
		c.lineTo(arrowHead/2+dy, arrowHead/2+dy)
	}
	c.lineTo(arrowHead/2+dy, h-dx)
	c.lineTo(arrowHead, h-dx)
	c.lineTo(arrowHead/2, h)
	c.lineTo(0, h-dx)
	c.lineTo(arrowHead/2-dy, h-dx)
	if rounded {
		c.lineTo(arrowHead/2-dy, arrowHead/2+dy)
		c.arcTo(dy*2, dy*2, 0, false, true, arrowHead/2+dy, arrowHead/2-dy)
	} else {
		c.lineTo(arrowHead/2-dy, arrowHead/2-dy)
	}
	c.lineTo(w-dx, arrowHead/2-dy)
	c.close()
	c.fillAndStroke()
})

// arrows2CalloutDoubleArrow ports mxShapeArrows2CalloutDoubleArrow: a box
// in the middle (2 notch wide) with arrows to the left and the right.
var arrows2CalloutDoubleArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	arrowHead := arrowsVal(s, "arrowHead", 0, h)
	c.begin()
	c.moveTo(w/2-notch, 0)
	c.lineTo(w/2+notch, 0)
	c.lineTo(w/2+notch, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy-arrowHead)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx, h*0.5+dy+arrowHead)
	c.lineTo(w-dx, h*0.5+dy)
	c.lineTo(w/2+notch, h*0.5+dy)
	c.lineTo(w/2+notch, h)
	c.lineTo(w/2-notch, h)
	c.lineTo(w/2-notch, h*0.5+dy)
	c.lineTo(dx, h*0.5+dy)
	c.lineTo(dx, h*0.5+dy+arrowHead)
	c.lineTo(0, h*0.5)
	c.lineTo(dx, h*0.5-dy-arrowHead)
	c.lineTo(dx, h*0.5-dy)
	c.lineTo(w/2-notch, h*0.5-dy)
	c.close()
	c.fillAndStroke()
})

// arrows2CalloutQuadArrow ports mxShapeArrows2CalloutQuadArrow: a box in
// the middle with arrows to all four sides.
var arrows2CalloutQuadArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	arrowHead := arrowsVal(s, "arrowHead", 0, h)
	c.begin()
	c.moveTo(w*0.5+dy, h*0.5-notch)
	c.lineTo(w*0.5+notch, h*0.5-notch)
	c.lineTo(w*0.5+notch, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy-arrowHead)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx, h*0.5+dy+arrowHead)
	c.lineTo(w-dx, h*0.5+dy)
	c.lineTo(w*0.5+notch, h*0.5+dy)
	c.lineTo(w*0.5+notch, h*0.5+notch)
	c.lineTo(w*0.5+dy, h*0.5+notch)
	c.lineTo(w*0.5+dy, h-dx)
	c.lineTo(w*0.5+dy+arrowHead, h-dx)
	c.lineTo(w*0.5, h)
	c.lineTo(w*0.5-dy-arrowHead, h-dx)
	c.lineTo(w*0.5-dy, h-dx)
	c.lineTo(w*0.5-dy, h*0.5+notch)
	c.lineTo(w*0.5-notch, h*0.5+notch)
	c.lineTo(w*0.5-notch, h*0.5+dy)
	c.lineTo(dx, h*0.5+dy)
	c.lineTo(dx, h*0.5+dy+arrowHead)
	c.lineTo(0, h*0.5)
	c.lineTo(dx, h*0.5-dy-arrowHead)
	c.lineTo(dx, h*0.5-dy)
	c.lineTo(w*0.5-notch, h*0.5-dy)
	c.lineTo(w*0.5-notch, h*0.5-notch)
	c.lineTo(w*0.5-dy, h*0.5-notch)
	c.lineTo(w*0.5-dy, dx)
	c.lineTo(w*0.5-dy-arrowHead, dx)
	c.lineTo(w*0.5, 0)
	c.lineTo(w*0.5+dy+arrowHead, dx)
	c.lineTo(w*0.5+dy, dx)
	c.close()
	c.fillAndStroke()
})

// arrows2CalloutDouble90Arrow ports mxShapeArrows2CalloutDouble90Arrow: a
// box at the top left (dx2 by dy2) with arrows to the right and down.
var arrows2CalloutDouble90Arrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy1 := arrowsVal(s, "dy1", 0.5, h)
	dx1 := arrowsVal(s, "dx1", 0.5, w)
	dx2 := arrowsVal(s, "dx2", 0, w)
	dy2 := arrowsVal(s, "dy2", 0, w)
	arrowHead := arrowsVal(s, "arrowHead", 0, h)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(dx2, 0)
	c.lineTo(dx2, dy2*0.5-dy1)
	c.lineTo(w-dx1, dy2*0.5-dy1)
	c.lineTo(w-dx1, dy2*0.5-dy1-arrowHead)
	c.lineTo(w, dy2*0.5)
	c.lineTo(w-dx1, dy2*0.5+dy1+arrowHead)
	c.lineTo(w-dx1, dy2*0.5+dy1)
	c.lineTo(dx2, dy2*0.5+dy1)
	c.lineTo(dx2, dy2)
	c.lineTo(dx2/2+dy1, dy2)
	c.lineTo(dx2/2+dy1, h-dx1)
	c.lineTo(dx2/2+dy1+arrowHead, h-dx1)
	c.lineTo(dx2/2, h)
	c.lineTo(dx2/2-dy1-arrowHead, h-dx1)
	c.lineTo(dx2/2-dy1, h-dx1)
	c.lineTo(dx2/2-dy1, dy2)
	c.lineTo(0, dy2)
	c.close()
	c.fillAndStroke()
})

// arrows2QuadArrow ports mxShapeArrows2QuadArrow: a cross with arrow heads.
var arrows2QuadArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	arrowHead := arrowsVal(s, "arrowHead", 0, h)
	c.begin()
	c.moveTo(w*0.5+dy, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy)
	c.lineTo(w-dx, h*0.5-dy-arrowHead)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx, h*0.5+dy+arrowHead)
	c.lineTo(w-dx, h*0.5+dy)
	c.lineTo(w*0.5+dy, h*0.5+dy)
	c.lineTo(w*0.5+dy, h-dx)
	c.lineTo(w*0.5+dy+arrowHead, h-dx)
	c.lineTo(w*0.5, h)
	c.lineTo(w*0.5-dy-arrowHead, h-dx)
	c.lineTo(w*0.5-dy, h-dx)
	c.lineTo(w*0.5-dy, h*0.5+dy)
	c.lineTo(dx, h*0.5+dy)
	c.lineTo(dx, h*0.5+dy+arrowHead)
	c.lineTo(0, h*0.5)
	c.lineTo(dx, h*0.5-dy-arrowHead)
	c.lineTo(dx, h*0.5-dy)
	c.lineTo(w*0.5-dy, h*0.5-dy)
	c.lineTo(w*0.5-dy, dx)
	c.lineTo(w*0.5-dy-arrowHead, dx)
	c.lineTo(w*0.5, 0)
	c.lineTo(w*0.5+dy+arrowHead, dx)
	c.lineTo(w*0.5+dy, dx)
	c.close()
	c.fillAndStroke()
})

// arrows2TriadArrow ports mxShapeArrows2TriadArrow: a T with arrow heads
// to the left, the right and the top.
var arrows2TriadArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	arrowHead := arrowsVal(s, "arrowHead", 0, h)
	c.begin()
	c.moveTo(w*0.5+arrowHead*0.5-dy, h-arrowHead+dy)
	c.lineTo(w-dx, h-arrowHead+dy)
	c.lineTo(w-dx, h-arrowHead)
	c.lineTo(w, h-arrowHead*0.5)
	c.lineTo(w-dx, h)
	c.lineTo(w-dx, h-dy)
	c.lineTo(dx, h-dy)
	c.lineTo(dx, h)
	c.lineTo(0, h-arrowHead*0.5)
	c.lineTo(dx, h-arrowHead)
	c.lineTo(dx, h-arrowHead+dy)
	c.lineTo(w*0.5-arrowHead*0.5+dy, h-arrowHead+dy)
	c.lineTo(w*0.5-arrowHead*0.5+dy, dx)
	c.lineTo(w*0.5-arrowHead*0.5, dx)
	c.lineTo(w*0.5, 0)
	c.lineTo(w*0.5+arrowHead*0.5, dx)
	c.lineTo(w*0.5+arrowHead*0.5-dy, dx)
	c.close()
	c.fillAndStroke()
})

// arrows2TailedArrow ports mxShapeArrows2TailedArrow and, with notched,
// mxShapeArrows2TailedNotchedArrow, which slants the tail by notch instead
// of by its height: an arrow with a wide tail (dx2 long, dy2 half high).
// dx1, dy1, dx2 and dy2 have no defaults: without them nothing is drawn.
func arrows2TailedArrow(notched bool) *shapeDef {
	return arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
		dy1 := arrowsVal(s, "dy1", undefined, h)
		dx1 := arrowsVal(s, "dx1", undefined, w)
		dy2 := arrowsVal(s, "dy2", undefined, h)
		dx2 := arrowsVal(s, "dx2", undefined, w)
		notch := arrowsVal(s, "notch", 0, w)
		arrowHead := arrowsVal(s, "arrowHead", 0, h)
		x2 := 0.0
		if dy2 != 0 {
			if notched {
				x2 = dx2 + notch*(dy2-dy1)/dy2
			} else {
				x2 = dx2 + dy2*(dy2-dy1)/dy2
			}
		}
		c.begin()
		c.moveTo(0, h*0.5-dy2)
		c.lineTo(dx2, h*0.5-dy2)
		c.lineTo(x2, h*0.5-dy1)
		c.lineTo(w-dx1, h*0.5-dy1)
		c.lineTo(w-dx1, h*0.5-dy1-arrowHead)
		c.lineTo(w, h*0.5)
		c.lineTo(w-dx1, h*0.5+dy1+arrowHead)
		c.lineTo(w-dx1, h*0.5+dy1)
		c.lineTo(x2, h*0.5+dy1)
		c.lineTo(dx2, h*0.5+dy2)
		c.lineTo(0, h*0.5+dy2)
		c.lineTo(notch, h*0.5)
		c.close()
		c.fillAndStroke()
	})
}

// arrows2StripedArrow ports mxShapeArrows2StripedArrow: an arrow whose
// tail (notch long) is cut into two stripes.
var arrows2StripedArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := h * 0.5 * arrowsVal(s, "dy", 0.5, 1)
	dx := arrowsVal(s, "dx", 0.5, w)
	notch := arrowsVal(s, "notch", 0, w)
	c.begin()
	c.moveTo(notch, dy)
	c.lineTo(w-dx, dy)
	c.lineTo(w-dx, 0)
	c.lineTo(w, h*0.5)
	c.lineTo(w-dx, h)
	c.lineTo(w-dx, h-dy)
	c.lineTo(notch, h-dy)
	c.close()
	c.moveTo(0, h-dy)
	c.lineTo(notch*0.16, h-dy)
	c.lineTo(notch*0.16, dy)
	c.lineTo(0, dy)
	c.close()
	c.moveTo(notch*0.32, h-dy)
	c.lineTo(notch*0.8, h-dy)
	c.lineTo(notch*0.8, dy)
	c.lineTo(notch*0.32, dy)
	c.close()
	c.fillAndStroke()
})

// arrows2JumpInArrow ports mxShapeArrows2JumpInArrow: an arrow curving up
// from the bottom left to a head at the top right.
var arrows2JumpInArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	dx := arrowsVal(s, "dx", 0.5, w)
	arrowHead := arrowsVal(s, "arrowHead", 40, h)
	c.begin()
	c.moveTo(w-dx, 0)
	c.lineTo(w, arrowHead*0.5)
	c.lineTo(w-dx, arrowHead)
	c.lineTo(w-dx, arrowHead/2+dy)
	c.arcTo(w-dx, h-arrowHead/2-dy, 0, false, false, 0, h)
	c.arcTo(w-dx, h-arrowHead/2+dy, 0, false, true, w-dx, arrowHead/2-dy)
	c.close()
	c.fillAndStroke()
})

// arrows2UTurnArrow ports mxShapeArrows2UTurnArrow: a U turn from the
// bottom right to a head at the top (dx2 is the head's length; it has no
// default: without it nothing is drawn).
var arrows2UTurnArrow = arrows2Shape(func(s *shape, c *nanCanvas, w, h float64) {
	dy := arrowsVal(s, "dy", 0.5, h)
	arrowHead := arrowsVal(s, "arrowHead", 40, h)
	dx := (h - arrowHead/2 + dy) / 2
	dx2 := math.Max(0, jsStyleFloat(s.style, "dx2", undefined))
	c.begin()
	c.moveTo(dx, 0)
	c.lineTo(dx+dx2, arrowHead*0.5)
	c.lineTo(dx, arrowHead)
	c.lineTo(dx, arrowHead/2+dy)
	c.arcTo(dx-2*dy, dx-2*dy, 0, false, false, dx, h-2*dy)
	c.lineTo(math.Max(w, dx), h-2*dy)
	c.lineTo(math.Max(w, dx), h)
	c.lineTo(dx, h)
	c.arcTo(dx, dx, 0, false, true, dx, arrowHead/2-dy)
	c.close()
	c.fillAndStroke()
})

// --- edges ---

// wedgeArrowNormal returns the offset from the first point to the corners
// of a wedge's base (startWidth long, across the line from the first to
// the last point) and the vector of that line.
func wedgeArrowNormal(s *shape, pts []point) (nx, ny, dx, dy float64) {
	sw := math.Max(0, jsStyleFloat(s.style, "startWidth", 20))
	p0, pe := pts[0], pts[len(pts)-1]
	dx, dy = pe.x-p0.x, pe.y-p0.y
	dist := math.Sqrt(dx*dx + dy*dy)
	return dx * sw / dist, dy * sw / dist, dx, dy
}

// paintWedgeArrow ports mxShapeArrowsWedgeArrow ("mxgraph.arrows2.wedgeArrow",
// an mxArrow): a triangle from a base startWidth wide at each side of the
// first point to the last point.
func paintWedgeArrow(s *shape, c0 *c2d, pts []point) {
	c := &nanCanvas{c2d: c0}
	nx, ny, _, _ := wedgeArrowNormal(s, pts)
	p0, pe := pts[0], pts[len(pts)-1]
	c.begin()
	c.moveTo(p0.x+ny, p0.y-nx)
	c.lineTo(p0.x-ny, p0.y+nx)
	c.lineTo(pe.x, pe.y)
	c.close()
	c.fillAndStroke()
}

// maxWedgeSteps bounds the strokes of dashed wedges for absurd lengths.
const maxWedgeSteps = 100000

// paintWedgeArrowDashed ports mxShapeArrowsWedgeArrowDashed
// ("mxgraph.arrows2.wedgeArrowDashed": nine strokes across the line,
// narrowing to the last point) and, with v2, mxShapeArrowsWedgeArrowDashed2
// ("mxgraph.arrows2.wedgeArrowDashed2": a stroke every stepSize).
func paintWedgeArrowDashed(v2 bool) func(s *shape, c *c2d, pts []point) {
	return func(s *shape, c0 *c2d, pts []point) {
		c := &nanCanvas{c2d: c0}
		nx, ny, dx, dy := wedgeArrowNormal(s, pts)
		steps := 8.0
		if v2 {
			stepSize := math.Max(0, jsStyleFloat(s.style, "stepSize", 10))
			steps = 0
			if stepSize > 0 {
				steps = math.Floor(math.Sqrt(dx*dx+dy*dy) / stepSize)
			}
			if steps > maxWedgeSteps {
				return
			}
		}
		pcx, pcy := pts[0].x, pts[0].y // current point on the line
		c.begin()
		for i := 0.0; i <= steps; i++ {
			cnx := nx * (steps - i) / steps
			cny := ny * (steps - i) / steps
			if i == steps {
				cnx = nx * (steps - i*0.98) / steps
				cny = ny * (steps - i*0.98) / steps
			}
			c.moveTo(pcx+cny, pcy-cnx)
			c.lineTo(pcx-cny, pcy+cnx)
			pcx += dx / steps
			pcy += dy / steps
		}
		c.stroke()
	}
}

// wedgeArrowBounds is the bounding box of the wedges, which replaces the
// edge's (useSvgBoundingBox: the box of what the SVG paints, grown by half
// the stroke width; the strokes of the dashed wedges fill the same box).
func wedgeArrowBounds(s *shape, r rect) rect {
	if len(s.points) < 2 {
		return r
	}
	nx, ny, _, _ := wedgeArrowNormal(s, s.points)
	p0, pe := s.points[0], s.points[len(s.points)-1]
	if !finite(nx) || !finite(ny) {
		return r
	}
	b, _ := pointsExtent([]point{{p0.x + ny, p0.y - nx}, {p0.x - ny, p0.y + nx}, pe})
	return b.grow(s.strokewidth / 2)
}

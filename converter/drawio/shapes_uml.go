package drawio

import "math"

// The UML shapes of Shapes.js (actor, boundary, entity, control, destroy,
// lifeline, frame, state, component, module, interfaces) and the
// sequence-diagram participant icons of its mermaid support.

func init() {
	registerShape("umlActor", umlActorShape)
	registerShape("umlBoundary", umlBoundaryShape)
	registerShape("umlEntity", umlEntityShape)
	registerShape("umlDestroy", umlDestroyShape)
	registerShape("umlControl", umlControlShape)
	registerShape("umlLifeline", umlLifelineShape)
	registerShape("umlFrame", umlFrameShape)
	registerShape("umlState", umlStateShape)
	registerShape("component", componentShape)
	registerShape("module", moduleShape)
	registerShape("lollipop", lollipopShape)
	registerShape("requires", requiresShape)
	registerShape("requiredInterface", requiredInterfaceShape)
	registerShape("providedRequiredInterface", providedRequiredInterfaceShape)

	registerShape("seqBoundary", seqBoundaryShape)
	registerShape("seqControl", seqIconShape(func(s *shape, c *c2d, x, y, w, h float64) {
		paintUmlControlBackground(s, c, x, y, w, h)
		paintUmlControlForeground(s, c, x, y, w, h)
	}))
	registerShape("seqEntity", seqIconShape(func(s *shape, c *c2d, x, y, w, h float64) {
		// UmlEntityShape.paintVertexShape on a plain mxShape: no center
		// circle, the glass of mxShape
		c.ellipse(x, y, w, h)
		c.fillAndStroke()
		if s.glass && s.fill != "" {
			s.paintGlassEffect(c, x, y, w, h, 0)
		}
		paintUmlEntityLine(c, x, y, w, h)
	}))
	registerShape("seqQueue", seqQueueShape)
	registerShape("seqCollections", seqCollectionsShape)
	registerShape("seqDatabase", seqDatabaseShape)
	registerShape("seqActorStick", seqActorStickShape)
}

// umlActorShape ports UmlActorShape: a stick figure.
var umlActorShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	// head
	c.ellipse(w/4, 0, w/2, h/4)
	c.fillAndStroke()
	c.begin()
	c.moveTo(w/2, h/4)
	c.lineTo(w/2, 2*h/3)
	// arms
	c.moveTo(w/2, h/3)
	c.lineTo(0, h/3)
	c.moveTo(w/2, h/3)
	c.lineTo(w, h/3)
	// legs
	c.moveTo(w/2, 2*h/3)
	c.lineTo(0, h)
	c.moveTo(w/2, 2*h/3)
	c.lineTo(w, h)
	c.stroke()
}}

// umlBoundaryShape ports UmlBoundaryShape: a circle attached to a
// vertical line.
var umlBoundaryShape = &shapeDef{
	paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		// base line
		c.begin()
		c.moveTo(0, h/4)
		c.lineTo(0, h*3/4)
		c.stroke()
		// horizontal line
		c.begin()
		c.moveTo(0, h/2)
		c.lineTo(w/6, h/2)
		c.stroke()
		// circle
		c.ellipse(w/6, 0, w*5/6, h)
		c.fillAndStroke()
	},
	labelMargins: func(s *shape, r rect) *rect { return &rect{r.w / 6, 0, 0, 0} },
}

func paintUmlEntityLine(c *c2d, x, y, w, h float64) {
	c.begin()
	c.moveTo(x+w/8, y+h)
	c.lineTo(x+w*7/8, y+h)
	c.stroke()
}

// umlEntityShape ports UmlEntityShape: a circle on a line.
var umlEntityShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	paintEllipseShape(s, c, x, y, w, h)
	paintUmlEntityLine(c, x, y, w, h)
}}

// umlDestroyShape ports UmlDestroyShape: an X.
var umlDestroyShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	c.begin()
	c.moveTo(w, 0)
	c.lineTo(0, h)
	c.moveTo(0, 0)
	c.lineTo(w, h)
	c.stroke()
}}

func paintUmlControlBackground(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	// upper line
	c.begin()
	c.moveTo(w*3/8, h/8*1.1)
	c.lineTo(w*5/8, 0)
	c.stroke()
	// circle
	c.ellipse(0, h/8, w, h*7/8)
	c.fillAndStroke()
}

func paintUmlControlForeground(s *shape, c *c2d, x, y, w, h float64) {
	// lower line (the canvas is still translated)
	c.begin()
	c.moveTo(w*3/8, h/8*1.1)
	c.lineTo(w*5/8, h/4)
	c.stroke()
}

// umlControlShape ports UmlControlShape: a circle with an arrow head.
var umlControlShape = &shapeDef{
	paintBackground: paintUmlControlBackground,
	paintForeground: paintUmlControlForeground,
	labelBounds: func(s *shape, r rect) rect {
		return rect{r.x, r.y + r.h/8, r.w, r.h * 7 / 8}
	},
}

// lifelineMirrored is UmlLifeline.isMirrored: lifelineMirror=1 repeats the
// head at the foot.
func lifelineMirrored(s *shape) bool { return s.style.get("lifelineMirror", "0") == "1" }

// lifelineHeadSize is UmlLifeline.getHeadSize.
func lifelineHeadSize(s *shape, h float64) float64 {
	limit := h
	if lifelineMirrored(s) {
		limit = h / 2
	}
	return math.Max(0, math.Min(limit, s.style.num("size", 40)))
}

// umlLifelineShape ports UmlLifeline: a head (a rectangle, or the shape
// named by participant) over a dashed line (lifelineDashed=0 for solid).
// The copy of the label in a mirrored foot is not painted: labels are
// laid out elsewhere.
var umlLifelineShape = &shapeDef{
	paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
		size := lifelineHeadSize(s, h)
		mirror := lifelineMirrored(s)
		participant, ok := s.style["participant"]
		if !ok {
			paintRectangleBackground(s, c, x, y, w, size)
			if mirror {
				paintRectangleBackground(s, c, x, y+h-size, w, size)
			}
		} else if d := shapeRegistry[participant]; d != nil && participant != "umlLifeline" {
			// a shape of the participant's class applied to the same state
			p := *s
			p.def = d
			p.name = participant
			c.save()
			p.paintVertexShape(c, x, y, w, size)
			c.restore()
			if mirror {
				c.save()
				p.paintVertexShape(c, x, y+h-size, w, size)
				c.restore()
			}
		}
		lineEnd := h
		if mirror {
			lineEnd = h - size
		}
		if size < lineEnd {
			c.setDashed(s.style.get("lifelineDashed", "1") == "1", false)
			c.begin()
			c.moveTo(x+w/2, y+size)
			c.lineTo(x+w/2, y+lineEnd)
			c.stroke()
		}
	},
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		if _, ok := s.style["participant"]; ok {
			return // the participant shape owns the outline
		}
		size := lifelineHeadSize(s, h)
		paintRectangleForeground(s, c, x, y, w, math.Min(h, size))
		if lifelineMirrored(s) {
			paintRectangleForeground(s, c, x, y+h-size, w, size)
		}
	},
	labelBounds: func(s *shape, r rect) rect {
		return rect{r.x, r.y, r.w, lifelineHeadSize(s, r.h)}
	},
	roundable: true,
}

// umlFrameShape ports UmlFrame: a rectangle with a name tab in the top
// left corner (width, height), the rest filled with swimlaneFillColor.
var umlFrameShape = &shapeDef{
	paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
		const co = 10 // corner
		w0 := math.Min(w, math.Max(co, s.style.num("width", 60)))
		h0 := math.Min(h, math.Max(co*1.5, s.style.num("height", 30)))
		if bg := colorOrNone(s.style.get("swimlaneFillColor", "none")); bg != "" {
			c.setFillColor(bg)
			c.rect(x, y, w, h)
			c.fill()
		}
		if s.fill != "" && s.gradient != "" {
			c.setGradient(s.fill, s.gradient, x, y, w, h, s.gradientDirection, 1, 1)
		} else {
			c.setFillColor(s.fill)
		}
		c.begin()
		c.moveTo(x, y)
		c.lineTo(x+w0, y)
		c.lineTo(x+w0, y+math.Max(0, h0-co*1.5))
		c.lineTo(x+math.Max(0, w0-co), y+h0)
		c.lineTo(x, y+h0)
		c.close()
		c.fillAndStroke()
		c.begin()
		c.moveTo(x+w0, y)
		c.lineTo(x+w, y)
		c.lineTo(x+w, y+h)
		c.lineTo(x, y+h)
		c.lineTo(x, y+h0)
		c.stroke()
	},
	labelMargins: func(s *shape, r rect) *rect {
		return &rect{0, 0, r.w - s.style.num("width", 60), r.h - s.style.num("height", 30)}
	},
}

// umlStateShape ports UMLStateShape: a (rounded) rectangle with optional
// symbols (umlStateSymbol=collapseState) and connection point references
// (umlStateConnection=connPointRefEntry or connPointRefExit) on the left.
var umlStateShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		rounded := styleBool(s.style, "rounded", false)
		arc := s.style.num("arcSize", 0.1)
		if !styleBool(s.style, "absoluteArcSize", false) {
			arc = math.Min(w, h) * arc
		}
		arc = math.Min(arc, math.Min(w*0.5, h*0.5))
		if !rounded {
			arc = 0
		}
		connPoint, hasConn := s.style["umlStateConnection"]
		dx := 0.0
		if hasConn {
			dx = 10
		}
		c.begin()
		c.moveTo(dx, arc)
		c.arcTo(arc, arc, 0, false, true, dx+arc, 0)
		c.lineTo(w-arc, 0)
		c.arcTo(arc, arc, 0, false, true, w, arc)
		c.lineTo(w, h-arc)
		c.arcTo(arc, arc, 0, false, true, w-arc, h)
		c.lineTo(dx+arc, h)
		c.arcTo(arc, arc, 0, false, true, dx, h-arc)
		c.close()
		c.fillAndStroke()
		if s.style.get("umlStateSymbol", "") == "collapseState" {
			c.roundrect(w-40, h-20, 10, 10, 3, 3)
			c.stroke()
			c.roundrect(w-20, h-20, 10, 10, 3, 3)
			c.stroke()
			c.begin()
			c.moveTo(w-30, h-15)
			c.lineTo(w-20, h-15)
			c.stroke()
		}
		switch connPoint {
		case "connPointRefEntry":
			c.ellipse(0, h*0.5-10, 20, 20)
			c.fillAndStroke()
		case "connPointRefExit":
			c.ellipse(0, h*0.5-10, 20, 20)
			c.fillAndStroke()
			c.begin()
			c.moveTo(5, h*0.5-5)
			c.lineTo(15, h*0.5+5)
			c.moveTo(15, h*0.5-5)
			c.lineTo(5, h*0.5+5)
			c.stroke()
		}
	},
	labelMargins: func(s *shape, r rect) *rect {
		if boundedLbl(s) && s.style.has("umlStateConnection") {
			return &rect{10, 0, 0, 0}
		}
		return nil
	},
}

// jettyShape is ModuleShape (module=true) and ComponentShape: a rectangle
// with two jetties (jettyWidth, jettyHeight) on the left.
func jettyShape(module bool) *shapeDef {
	defW, defH := 32.0, 12.0
	if module {
		defW, defH = 20, 10
	}
	return cylinderShape(func(s *shape, c *c2d, w, h float64, fg bool) {
		dx := s.style.num("jettyWidth", defW)
		dy := s.style.num("jettyHeight", defH)
		x0 := dx / 2
		x1 := x0 + dx/2
		y0 := 0.3*h - dy/2
		y1 := 0.7*h - dy/2
		if module {
			y0 = math.Min(dy, h-dy)
			y1 = math.Min(y0+2*dy, h-dy)
		}
		if fg {
			c.moveTo(x0, y0)
			c.lineTo(x1, y0)
			c.lineTo(x1, y0+dy)
			c.lineTo(x0, y0+dy)
			c.moveTo(x0, y1)
			c.lineTo(x1, y1)
			c.lineTo(x1, y1+dy)
			c.lineTo(x0, y1+dy)
			return
		}
		c.moveTo(x0, 0)
		c.lineTo(w, 0)
		c.lineTo(w, h)
		c.lineTo(x0, h)
		c.lineTo(x0, y1+dy)
		c.lineTo(0, y1+dy)
		c.lineTo(0, y1)
		c.lineTo(x0, y1)
		c.lineTo(x0, y0+dy)
		c.lineTo(0, y0+dy)
		c.lineTo(0, y0)
		c.lineTo(x0, y0)
		c.close()
	})
}

var (
	componentShape = jettyShape(false)
	moduleShape    = jettyShape(true)
)

// lollipopShape ports LollipopShape: a provided interface (a circle on a
// stick).
var lollipopShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	sz := s.style.num("size", 10)
	c.translate(x, y)
	c.ellipse((w-sz)/2, 0, sz, sz)
	c.fillAndStroke()
	c.begin()
	c.moveTo(w/2, sz)
	c.lineTo(w/2, h)
	c.stroke()
}}

// requiresShape ports RequiresShape: a required interface (a cup on a
// stick).
var requiresShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	sz := s.style.num("size", 10)
	inset := s.style.num("inset", 2) + s.strokewidth
	c.translate(x, y)
	c.begin()
	c.moveTo(w/2, sz+inset)
	c.lineTo(w/2, h)
	c.stroke()
	c.begin()
	c.moveTo((w-sz)/2-inset, sz/2)
	c.quadTo((w-sz)/2-inset, sz+inset, w/2, sz+inset)
	c.quadTo((w+sz)/2+inset, sz+inset, (w+sz)/2+inset, sz/2)
	c.stroke()
}}

// requiredInterfaceShape ports RequiredInterfaceShape: a half circle.
var requiredInterfaceShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	c.begin()
	c.moveTo(0, 0)
	c.quadTo(w, 0, w, h/2)
	c.quadTo(w, h, 0, h)
	c.stroke()
}}

// providedRequiredInterfaceShape ports ProvidedRequiredInterfaceShape: a
// circle in a half circle.
var providedRequiredInterfaceShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	inset := s.style.num("inset", 2) + s.strokewidth
	c.translate(x, y)
	c.ellipse(0, inset, w-2*inset, h-2*inset)
	c.fillAndStroke()
	c.begin()
	c.moveTo(w/2, 0)
	c.quadTo(w, 0, w, h/2)
	c.quadTo(w, h, w/2, h)
	c.stroke()
}}

// seqIconShape is makeSeqIcon: an icon shape painted at 44 units at most,
// centered horizontally at the top.
func seqIconShape(paint func(s *shape, c *c2d, x, y, w, h float64)) *shapeDef {
	return &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
		size := math.Min(44, math.Min(w, h))
		paint(s, c, x+(w-size)/2, y, size, size)
	}}
}

// seqBoundaryShape ports SeqBoundaryShape: mermaid's boundary participant
// (a bar and a handle left of a circle centered on the lifeline).
var seqBoundaryShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	size := math.Min(44, math.Min(w, h))
	circleR := size / 2
	handlerW := size * 0.9
	barH := size * 0.45
	circleX := x + w/2 - circleR
	cy := y + size/2
	handlerEndX := circleX + handlerW*0.175
	leftX := handlerEndX - handlerW
	c.begin()
	c.moveTo(leftX, cy-barH/2)
	c.lineTo(leftX, cy+barH/2)
	c.stroke()
	c.begin()
	c.moveTo(leftX, cy)
	c.lineTo(handlerEndX, cy)
	c.stroke()
	c.ellipse(circleX, y, size, size)
	c.fillAndStroke()
}}

// seqQueueShape ports SeqQueueShape: a horizontal cylinder.
var seqQueueShape = &shapeDef{
	paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
		ry := h / 2
		rx := ry / (2.5 + h/50)
		c.save()
		c.translate(x, y)
		c.begin()
		c.moveTo(rx, 0)
		c.lineTo(w-rx, 0)
		c.arcTo(rx, ry, 0, false, true, w-rx, h)
		c.lineTo(rx, h)
		c.arcTo(rx, ry, 0, false, true, rx, 0)
		c.close()
		c.fillAndStroke()
		c.restore()
	},
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		ry := h / 2
		rx := ry / (2.5 + h/50)
		c.save()
		c.translate(x, y)
		c.begin()
		c.moveTo(w-rx, 0)
		c.arcTo(rx, ry, 0, false, false, w-rx, h)
		c.stroke()
		c.restore()
	},
}

// seqCollectionsShape ports SeqCollectionsShape: two stacked rectangles.
var seqCollectionsShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	const off = 6
	c.rect(x+off, y-off, w, h)
	c.fillAndStroke()
	c.rect(x, y, w, h)
	c.fillAndStroke()
}}

// seqDatabaseShape ports SeqDatabaseShape: a cylinder a third of the width.
var seqDatabaseShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	iconW, iconH := w/3, w/3
	rx := iconW / 2
	ry := rx / (2.5 + iconW/50)
	c.translate(x+(w-iconW)/2, y)
	c.begin()
	c.moveTo(0, ry)
	c.arcTo(rx, ry, 0, false, false, iconW, ry)
	c.arcTo(rx, ry, 0, false, false, 0, ry)
	c.lineTo(0, iconH-ry)
	c.arcTo(rx, ry, 0, false, false, iconW, iconH-ry)
	c.lineTo(iconW, ry)
	c.fillAndStroke()
}}

// seqActorStickShape ports SeqActorStickShape: a stick figure of fixed size.
var seqActorStickShape = &shapeDef{paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
	const headR, torsoH, armsW, legsH, legsW = 15, 20, 36, 15, 32
	cx := x + w/2
	headCY := y + headR
	torsoTop := headCY + headR
	torsoBot := torsoTop + torsoH
	armsY := torsoTop + torsoH/2
	legsBot := torsoBot + legsH
	c.ellipse(cx-headR, headCY-headR, headR*2, headR*2)
	c.fillAndStroke()
	c.begin()
	c.moveTo(cx, torsoTop)
	c.lineTo(cx, torsoBot)
	c.moveTo(cx-armsW/2, armsY)
	c.lineTo(cx+armsW/2, armsY)
	c.moveTo(cx, torsoBot)
	c.lineTo(cx-legsW/2, legsBot)
	c.moveTo(cx, torsoBot)
	c.lineTo(cx+legsW/2, legsBot)
	c.stroke()
}}

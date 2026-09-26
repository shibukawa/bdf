package drawio

// The BPMN shapes of draw.io's shapes/bpmn/mxBpmnShape2.js: events and
// gateways (the deprecated mxBpmnShape, mxShapeBpmnEvent and
// mxShapeBpmnGateway), tasks, data objects, the swimlane, conversations
// and the send marker. Their setShadow(false) calls are left out: the
// shadow is drawn for the whole shape (see shape.paint). The markers of
// tasks are stencils of the BPMN library (stencils/bpmn.xml).

func init() {
	registerShape("mxgraph.bpmn.shape", bpmnShape)
	registerShape("mxgraph.bpmn.sendMarker", bpmnSendMarkerShape)
	registerShape("mxgraph.bpmn.event", bpmnEventShape)
	registerShape("mxgraph.bpmn.gateway2", bpmnGatewayShape)
	registerShape("mxgraph.bpmn.task", bpmnTaskShape(false))
	registerShape("mxgraph.bpmn.task2", bpmnTaskShape(true))
	registerShape("mxgraph.bpmn.data", bpmnDataShape(false))
	registerShape("mxgraph.bpmn.data2", bpmnDataShape(true))
	registerShape("mxgraph.bpmn.swimlane", bpmnSwimlaneShape)
	registerShape("mxgraph.bpmn.conversation", bpmnConversationShape(false))
	registerShape("mxgraph.bpmn.conversation2", bpmnConversationShape(true))
}

// --- events ---

// bpmnOutlines are the outlines of events (mxBpmnShape.prototype.outlines
// and mxShapeBpmnEvent.prototype.outlines, which draw alike) in the box
// 0, 0, w, h; keepWidth leaves the stroke width of end events as it is
// (mxShapeBpmn2Task2 replaces setStrokeWidth with a function that does
// nothing while it draws the event).
var bpmnOutlines = map[string]func(c *c2d, w, h float64, keepWidth bool){
	"none": func(c *c2d, w, h float64, keepWidth bool) {},
	"standard": func(c *c2d, w, h float64, keepWidth bool) {
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
	},
	"eventInt": func(c *c2d, w, h float64, keepWidth bool) {
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
	},
	"eventNonint": func(c *c2d, w, h float64, keepWidth bool) {
		dashed := c.st.dashed
		c.setDashed(true, false)
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
		c.setDashed(dashed, false)
	},
	"catching": func(c *c2d, w, h float64, keepWidth bool) {
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
		const inset = 2
		c.ellipse(inset, inset, w-2*inset, h-2*inset)
		c.stroke()
	},
	"boundInt": func(c *c2d, w, h float64, keepWidth bool) {
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
		const inset = 2
		c.ellipse(inset, inset, w-2*inset, h-2*inset)
		c.stroke()
	},
	"boundNonint": func(c *c2d, w, h float64, keepWidth bool) {
		dashed := c.st.dashed
		c.setDashed(true, false)
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
		const inset = 2
		c.ellipse(inset, inset, w-2*inset, h-2*inset)
		c.stroke()
		c.setDashed(dashed, false)
	},
	"throwing": func(c *c2d, w, h float64, keepWidth bool) {
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
		const inset = 2
		c.ellipse(w*0.02+inset, h*0.02+inset, w*0.96-2*inset, h*0.96-2*inset)
		c.stroke()
	},
	"end": func(c *c2d, w, h float64, keepWidth bool) {
		sw := c.st.strokeWidth
		if !keepWidth {
			c.setStrokeWidth(sw * 3)
		}
		c.ellipse(0, 0, w, h)
		c.fillAndStroke()
		if !keepWidth {
			c.setStrokeWidth(sw)
		}
	},
}

// bpmnSymbolBoxes place the symbols of events in the event's box: the
// translation and the size, as fractions of the box (the if-chain of
// mxBpmnShape.redrawPath and mxShapeBpmnEvent.strictDrawShape; only
// mxBpmnShape has exclusiveGw).
var bpmnSymbolBoxes = map[string][4]float64{
	"message":          {0.15, 0.3, 0.7, 0.4},
	"timer":            {0.11, 0.11, 0.78, 0.78},
	"escalation":       {0.19, 0.15, 0.62, 0.57},
	"conditional":      {0.3, 0.16, 0.4, 0.68},
	"link":             {0.27, 0.33, 0.46, 0.34},
	"error":            {0.212, 0.243, 0.58, 0.507},
	"cancel":           {0.22, 0.22, 0.56, 0.56},
	"compensation":     {0.28, 0.35, 0.44, 0.3},
	"signal":           {0.19, 0.15, 0.62, 0.57},
	"multiple":         {0.2, 0.19, 0.6, 0.565},
	"parallelMultiple": {0.2, 0.2, 0.6, 0.6},
	"terminate":        {0.05, 0.05, 0.9, 0.9},
	"exclusiveGw":      {0.12, 0, 0.76, 1},
}

// bpmnSymbol paints an event symbol in the box 0, 0, w, h; inverse is set
// for throwing and end events, whose symbols are filled with the stroke
// color.
type bpmnSymbol func(s *shape, c *c2d, w, h float64, inverse bool)

// bpmnShapeSymbols are mxBpmnShape.prototype.symbols.
var bpmnShapeSymbols = map[string]bpmnSymbol{
	"general": func(s *shape, c *c2d, w, h float64, inverse bool) {},
	"message": bpmnMessageSymbol,
	"timer":   bpmnTimerSymbol,
	"escalation": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnEscalationPath(c, w, h)
		c.fillAndStroke()
	},
	"conditional": bpmnConditionalSymbol,
	"link": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnLinkPath(c, w, h)
		c.fillAndStroke()
	},
	"error": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnErrorPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"cancel": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnCancelPath(c, w, h)
		c.fillAndStroke()
	},
	"compensation": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnCompensationPath(c, w, h)
		c.fillAndStroke()
	},
	"signal": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnSignalPath(c, w, h)
		c.fillAndStroke()
	},
	"multiple": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnMultiplePath(c, w, h)
		c.fillAndStroke()
	},
	"parallelMultiple": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnCrossPath(c, w, h)
		c.fillAndStroke()
	},
	"terminate": bpmnTerminateSymbol,
	"exclusiveGw": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnInverted(c, func() { bpmnExclusivePath(c, w, h) })
	},
	"parallelGw": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnInverted(c, func() { bpmnCrossPath(c, w, h) })
	},
	"complexGw": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnInverted(c, func() { bpmnComplexPath(c, w, h) })
	},
	"star": bpmnStarSymbol,
}

// bpmnEventSymbols are mxShapeBpmnEvent.prototype.symbols, which set miter
// limits and only fill the inverse symbols.
var bpmnEventSymbols = map[string]bpmnSymbol{
	"general": func(s *shape, c *c2d, w, h float64, inverse bool) {},
	"message": bpmnMessageSymbol,
	"timer":   bpmnTimerSymbol,
	"escalation": func(s *shape, c *c2d, w, h float64, inverse bool) {
		c.setMiterLimit(6)
		bpmnEscalationPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"conditional": bpmnConditionalSymbol,
	"link": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnLinkPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"error": func(s *shape, c *c2d, w, h float64, inverse bool) {
		c.setMiterLimit(7)
		bpmnErrorPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"cancel": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnCancelPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"compensation": func(s *shape, c *c2d, w, h float64, inverse bool) {
		c.setMiterLimit(1)
		bpmnCompensationPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"signal": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnSignalPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"multiple": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnMultiplePath(c, w, h)
		bpmnFill(c, inverse)
	},
	"parallelMultiple": func(s *shape, c *c2d, w, h float64, inverse bool) {
		bpmnCrossPath(c, w, h)
		bpmnFill(c, inverse)
	},
	"terminate": bpmnTerminateSymbol,
	"star":      bpmnStarSymbol,
}

// bpmnFill fills an inverse symbol, and fills and strokes any other.
func bpmnFill(c *c2d, inverse bool) {
	if inverse {
		c.fill()
	} else {
		c.fillAndStroke()
	}
}

// bpmnInverted fills and strokes a path with the fill and stroke colors
// swapped (the gateway symbols of mxBpmnShape).
func bpmnInverted(c *c2d, path func()) {
	strokeColor, fillColor := c.st.strokeColor, c.st.fillColor
	c.setStrokeColor(fillColor)
	c.setFillColor(strokeColor)
	path()
	c.fillAndStroke()
	c.setStrokeColor(strokeColor)
	c.setFillColor(fillColor)
}

// paintBpmnSymbol paints the symbol sym of an event whose outline is
// outline in the box 0, 0, w, h, placed by bpmnSymbolBoxes, with the
// colors swapped for throwing and end events.
func paintBpmnSymbol(s *shape, c *c2d, symbols map[string]bpmnSymbol, sym, outline string, w, h float64) {
	f := symbols[sym]
	if f == nil {
		return
	}
	strokeColor, fillColor := c.st.strokeColor, c.st.fillColor
	if b, ok := bpmnSymbolBoxes[sym]; ok {
		c.translate(w*b[0], h*b[1])
		w *= b[2]
		h *= b[3]
	}
	inverse := false
	if sym == "star" {
		c.setFillColor(strokeColor)
	} else if outline == "throwing" || outline == "end" {
		c.setStrokeColor(fillColor)
		c.setFillColor(strokeColor)
		inverse = true
	}
	f(s, c, w, h, inverse)
	if sym == "star" {
		c.setFillColor(fillColor)
	} else if outline == "throwing" || outline == "end" {
		c.setStrokeColor(strokeColor)
		c.setFillColor(fillColor)
	}
}

func bpmnMessageSymbol(s *shape, c *c2d, w, h float64, inverse bool) {
	c.rect(0, 0, w, h)
	c.fillAndStroke()
	if s.style.get("fillColor", "none") == "none" && inverse {
		c.setStrokeColor("#ffffff")
	}
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w*0.5, h*0.5)
	c.lineTo(w, 0)
	c.stroke()
}

func bpmnTimerSymbol(s *shape, c *c2d, w, h float64, inverse bool) {
	c.ellipse(0, 0, w, h)
	c.fillAndStroke()
	c.begin()
	c.moveTo(w*0.5, 0)
	c.lineTo(w*0.5, h*0.0642)
	c.moveTo(w*0.7484, h*0.0654)
	c.lineTo(w*0.7126, h*0.1281)
	c.moveTo(w*0.93, h*0.2471)
	c.lineTo(w*0.8673, h*0.2854)
	c.moveTo(w, h*0.5)
	c.lineTo(w*0.9338, h*0.5)
	c.moveTo(w*0.93, h*0.7509)
	c.lineTo(w*0.8673, h*0.7126)
	c.moveTo(w*0.7484, h*0.9326)
	c.lineTo(w*0.7126, h*0.8699)
	c.moveTo(w*0.5, h*0.9338)
	c.lineTo(w*0.5, h)
	c.moveTo(w*0.2496, h*0.9325)
	c.lineTo(w*0.2854, h*0.8699)
	c.moveTo(w*0.068, h*0.7509)
	c.lineTo(w*0.1307, h*0.7126)
	c.moveTo(0, h*0.5)
	c.lineTo(w*0.0642, h*0.5)
	c.moveTo(w*0.068, h*0.2471)
	c.lineTo(w*0.1307, h*0.2854)
	c.moveTo(w*0.2496, h*0.0654)
	c.lineTo(w*0.2854, h*0.1281)
	c.moveTo(w*0.5246, h*0.0706)
	c.lineTo(w*0.5, h*0.5)
	c.lineTo(w*0.7804, h*0.5118)
	c.stroke()
}

func bpmnEscalationPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h)
	c.lineTo(w*0.5, 0)
	c.lineTo(w, h)
	c.lineTo(w*0.5, h*0.5)
	c.close()
}

func bpmnConditionalSymbol(s *shape, c *c2d, w, h float64, inverse bool) {
	c.rect(0, 0, w, h)
	c.fillAndStroke()
	c.begin()
	c.moveTo(0, h*0.1027)
	c.lineTo(w*0.798, h*0.1027)
	c.moveTo(0, h*0.3669)
	c.lineTo(w*0.798, h*0.3669)
	c.moveTo(0, h*0.6311)
	c.lineTo(w*0.798, h*0.6311)
	c.moveTo(0, h*0.8953)
	c.lineTo(w*0.798, h*0.8953)
	c.stroke()
}

func bpmnLinkPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h*0.76)
	c.lineTo(0, h*0.24)
	c.lineTo(w*0.63, h*0.24)
	c.lineTo(w*0.63, 0)
	c.lineTo(w, h*0.5)
	c.lineTo(w*0.63, h)
	c.lineTo(w*0.63, h*0.76)
	c.close()
}

func bpmnErrorPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h)
	c.lineTo(w*0.3287, h*0.123)
	c.lineTo(w*0.6194, h*0.6342)
	c.lineTo(w, 0)
	c.lineTo(w*0.6625, h*0.939)
	c.lineTo(w*0.3717, h*0.5064)
	c.close()
}

func bpmnCancelPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(w*0.1051, 0)
	c.lineTo(w*0.5, h*0.3738)
	c.lineTo(w*0.8909, 0)
	c.lineTo(w, h*0.1054)
	c.lineTo(w*0.623, h*0.5)
	c.lineTo(w, h*0.8926)
	c.lineTo(w*0.8909, h)
	c.lineTo(w*0.5, h*0.6242)
	c.lineTo(w*0.1051, h)
	c.lineTo(0, h*0.8926)
	c.lineTo(w*0.373, h*0.5)
	c.lineTo(0, h*0.1054)
	c.close()
}

func bpmnCompensationPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h*0.5)
	c.lineTo(w*0.5, 0)
	c.lineTo(w*0.5, h)
	c.close()
	c.moveTo(w*0.5, h*0.5)
	c.lineTo(w, 0)
	c.lineTo(w, h)
	c.close()
}

func bpmnSignalPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h)
	c.lineTo(w*0.5, 0)
	c.lineTo(w, h)
	c.close()
}

func bpmnMultiplePath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h*0.39)
	c.lineTo(w*0.5, 0)
	c.lineTo(w, h*0.39)
	c.lineTo(w*0.815, h)
	c.lineTo(w*0.185, h)
	c.close()
}

// bpmnCrossPath is the plus of parallelMultiple and parallelGw.
func bpmnCrossPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(w*0.38, 0)
	c.lineTo(w*0.62, 0)
	c.lineTo(w*0.62, h*0.38)
	c.lineTo(w, h*0.38)
	c.lineTo(w, h*0.62)
	c.lineTo(w*0.62, h*0.62)
	c.lineTo(w*0.62, h)
	c.lineTo(w*0.38, h)
	c.lineTo(w*0.38, h*0.62)
	c.lineTo(0, h*0.62)
	c.lineTo(0, h*0.38)
	c.lineTo(w*0.38, h*0.38)
	c.close()
}

func bpmnTerminateSymbol(s *shape, c *c2d, w, h float64, inverse bool) {
	c.ellipse(0, 0, w, h)
	c.fillAndStroke()
}

// bpmnExclusivePath is the x of exclusive gateways.
func bpmnExclusivePath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(w*0.105, 0)
	c.lineTo(w*0.5, h*0.38)
	c.lineTo(w*0.895, h*0)
	c.lineTo(w, h*0.11)
	c.lineTo(w*0.6172, h*0.5)
	c.lineTo(w, h*0.89)
	c.lineTo(w*0.895, h)
	c.lineTo(w*0.5, h*0.62)
	c.lineTo(w*0.105, h)
	c.lineTo(0, h*0.89)
	c.lineTo(w*0.3808, h*0.5)
	c.lineTo(0, h*0.11)
	c.close()
}

// bpmnComplexPath is the asterisk of complex gateways.
func bpmnComplexPath(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, h*0.44)
	c.lineTo(w*0.36, h*0.44)
	c.lineTo(w*0.1, h*0.18)
	c.lineTo(w*0.18, h*0.1)
	c.lineTo(w*0.44, h*0.36)
	c.lineTo(w*0.44, 0)
	c.lineTo(w*0.56, 0)
	c.lineTo(w*0.56, h*0.36)
	c.lineTo(w*0.82, h*0.1)
	c.lineTo(w*0.90, h*0.18)
	c.lineTo(w*0.64, h*0.44)
	c.lineTo(w, h*0.44)
	c.lineTo(w, h*0.56)
	c.lineTo(w*0.64, h*0.56)
	c.lineTo(w*0.9, h*0.82)
	c.lineTo(w*0.82, h*0.9)
	c.lineTo(w*0.56, h*0.64)
	c.lineTo(w*0.56, h)
	c.lineTo(w*0.44, h)
	c.lineTo(w*0.44, h*0.64)
	c.lineTo(w*0.18, h*0.9)
	c.lineTo(w*0.1, h*0.82)
	c.lineTo(w*0.36, h*0.56)
	c.lineTo(0, h*0.56)
	c.close()
}

func bpmnStarSymbol(s *shape, c *c2d, w, h float64, inverse bool) {
	c.translate(w/5, h/6)
	h *= 2.0 / 3
	w *= 3.0 / 5
	c.begin()
	c.moveTo(0, h/4)
	c.lineTo(w/3, h/4)
	c.lineTo(w/2, 0)
	c.lineTo(2*w/3, h/4)
	c.lineTo(w, h/4)
	c.lineTo(5*w/6, h/2)
	c.lineTo(w, 3*h/4)
	c.lineTo(2*w/3, 3*h/4)
	c.lineTo(w/2, h)
	c.lineTo(w/3, 3*h/4)
	c.lineTo(0, 3*h/4)
	c.lineTo(w/6, h/2)
	c.close()
	c.fillAndStroke()
}

// bpmnShape ports mxBpmnShape ("mxgraph.bpmn.shape", deprecated): an event
// outline (outline) with a symbol (symbol), in the middle of a rhombus
// for background=gateway.
var bpmnShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	bg := s.style.get("background", "none")
	isGateway := bg == "gateway"
	// the background layer: only known backgrounds move to the bounds
	switch bg {
	case "none":
		c.translate(x, y)
	case "gateway":
		c.translate(x, y)
		c.begin()
		c.moveTo(w/2, 0)
		c.lineTo(w, h/2)
		c.lineTo(w/2, h)
		c.lineTo(0, h/2)
		c.close()
		c.fillAndStroke()
	}
	// the outline layer
	outline := s.style.get("outline", "none")
	ow, oh := w, h
	if isGateway {
		c.translate(w/4, h/4)
		ow, oh = w/2, h/2
	}
	if f := bpmnOutlines[outline]; f != nil {
		f(c, ow, oh, false)
	}
	// the symbol layer
	if isGateway {
		w, h = w/2, h/2
	}
	paintBpmnSymbol(s, c, bpmnShapeSymbols, s.style.get("symbol", ""), outline, w, h)
}}

// paintBpmnEvent is mxShapeBpmnEvent.strictDrawShape: the outline and the
// symbol of an event in the box x, y, w, h. keepWidth is for
// mxShapeBpmn2Task2 (see bpmnOutlines).
func paintBpmnEvent(s *shape, c *c2d, x, y, w, h float64, outline, symbol string, keepWidth bool) {
	if f := bpmnOutlines[outline]; f != nil {
		c.translate(x, y)
		f(c, w, h, keepWidth)
	}
	paintBpmnSymbol(s, c, bpmnEventSymbols, symbol, outline, w, h)
}

// bpmnEventShape ports mxShapeBpmnEvent ("mxgraph.bpmn.event").
var bpmnEventShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	paintBpmnEvent(s, c, 0, 0, w, h, s.style.get("outline", "none"), s.style.get("symbol", ""), false)
}}

// bpmnGatewayShape ports mxShapeBpmnGateway ("mxgraph.bpmn.gateway2"): a
// rhombus with an event (gwType=event, the default, with an outline) or
// the symbol of an exclusive, parallel or complex gateway.
var bpmnGatewayShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	c.begin()
	c.moveTo(w*0.5, 0)
	c.lineTo(w, h*0.5)
	c.lineTo(w*0.5, h)
	c.lineTo(0, h*0.5)
	c.close()
	c.fillAndStroke()
	symbolW, symbolH := w*0.6, h*0.6
	gwType := s.style.get("gwType", "event")
	outline := s.style.get("outline", "none")
	symbol := s.style.get("symbol", "standard")
	tmpW, tmpH := w*0.5, h*0.5
	switch {
	case gwType == "event" && outline != "none":
		paintBpmnEvent(s, c, (w-symbolW)*0.5, (h-symbolH)*0.5, symbolW, symbolH, outline, symbol, false)
	case gwType == "exclusive":
		c.translate(w*0.31, h*0.25)
		tmpW *= 0.76
		fillColor := c.st.fillColor
		c.setFillColor(c.st.strokeColor)
		bpmnExclusivePath(c, tmpW, tmpH)
		c.fillAndStroke()
		c.setFillColor(fillColor)
		c.translate(-tmpW*0.12, 0)
	case gwType == "parallel":
		c.translate(w*0.25, h*0.25)
		fillColor := c.st.fillColor
		c.setFillColor(c.st.strokeColor)
		bpmnCrossPath(c, tmpW, tmpH)
		c.fillAndStroke()
		c.setFillColor(fillColor)
	case gwType == "complex":
		c.translate(w*0.25, h*0.25)
		fillColor := c.st.fillColor
		c.setFillColor(c.st.strokeColor)
		bpmnComplexPath(c, tmpW, tmpH)
		c.fillAndStroke()
		c.setFillColor(fillColor)
	}
}}

// paintBpmnSendMarker is mxShapeBpmn2SendMarker.paintVertexShape: an
// envelope. Tasks call it on the prototype, so its translation stays.
func paintBpmnSendMarker(c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	c.rect(0, 0, w, h)
	c.fillAndStroke()
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w*0.5, h*0.5)
	c.lineTo(w, 0)
	c.stroke()
}

// bpmnSendMarkerShape ports mxShapeBpmn2SendMarker ("mxgraph.bpmn.sendMarker").
var bpmnSendMarkerShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	paintBpmnSendMarker(c, x, y, w, h)
}}

// --- tasks ---

// bpmnStencil draws a shape of the BPMN stencil library in the box x, y,
// w, h. With keepWidth the stencil keeps the canvas's stroke width, as
// mxShapeBpmn2Task2 does by replacing setStrokeWidth with a function that
// does nothing: the markers inherit the stroke width of the style and set
// no other, so the shape is given the canvas's width as its style.
// It reports whether the library has the shape.
func bpmnStencil(s *shape, c *c2d, name string, x, y, w, h float64, keepWidth bool) bool {
	st := lookupStencil(name)
	if st == nil {
		return false
	}
	if keepWidth {
		t := *s
		t.style = s.style.clone()
		t.style["strokeWidth"] = fmtNum(c.st.strokeWidth)
		s = &t
	}
	st.drawShape(c, s, x, y, w, h)
	return true
}

// paintBpmnLines paints the three lines of multi-instance markers and
// collections in a 12 by 12 box at x, y (vertical lines, or horizontal for
// sequential multi-instances); the canvas stays translated to x, y, and
// the callers translate back as the JavaScript does.
func paintBpmnLines(c *c2d, x, y float64, horizontal bool) {
	c.translate(x, y)
	c.begin()
	if horizontal {
		c.moveTo(0, 2.4)
		c.lineTo(12, 2.4)
		c.moveTo(0, 6)
		c.lineTo(12, 6)
		c.moveTo(0, 9.6)
		c.lineTo(12, 9.6)
	} else {
		c.moveTo(2.4, 0)
		c.lineTo(2.4, 12)
		c.moveTo(6, 0)
		c.lineTo(6, 12)
		c.moveTo(9.6, 0)
		c.lineTo(9.6, 12)
	}
	c.stroke()
}

// paintBpmnSubprocessMarker paints the boxed plus of sub-processes at the
// canvas's origin.
func paintBpmnSubprocessMarker(c *c2d) {
	c.rect(0, 0, 14, 14)
	c.stroke()
	c.begin()
	c.moveTo(4, 7)
	c.lineTo(10, 7)
	c.moveTo(7, 4)
	c.lineTo(7, 10)
	c.stroke()
}

// bpmnTaskShape ports mxShapeBpmn2Task ("mxgraph.bpmn.task", legacy) and,
// with v2, mxShapeBpmn2Task2 ("mxgraph.bpmn.task2"), which insets the
// markers by the stroke width and draws them with a width of 1: an
// mxgraph.basic.rect (a double outline for bpmnShapeType=transaction,
// dashed for subprocess, thick for call) with loop markers at the bottom
// (isLoopStandard and the like), a task marker at the top left
// (taskMarker) and an event there (outline, symbol).
func bpmnTaskShape(v2 bool) *shapeDef {
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		st := s.style
		bpmnShapeType := st.get("bpmnShapeType", "task")
		taskMarker := st.get("taskMarker", "abstract")
		strokeWidth := st.num("strokeWidth", 1)
		inset := st.num("indent", 3)
		offsetY := 14.0
		if v2 {
			offsetY = 14 + strokeWidth*0.5
			if bpmnShapeType == "call" {
				offsetY = 14 + strokeWidth*2
			}
		}
		c.translate(x, y)

		// mxShapeBasicRect2.strictDrawShape with the override styles
		r := *s
		r.style = st.clone()
		switch bpmnShapeType {
		case "transaction":
			offsetY += inset
			r.style["rectOutline"] = "double"
			r.style["indent"] = "3"
		case "subprocess":
			r.style["dashed"] = "1"
		case "call":
			if v2 {
				r.style["strokeWidth"] = fmtNum(strokeWidth * 4)
			} else {
				r.style["strokeWidth"] = "4"
			}
		}
		paintRect2(&r, c, w, h)
		c.setStrokeWidth(strokeWidth)
		// setDashed(dashed) and then setDashed(false)
		c.setDashed(false, false)

		isLoopSub := styleBool(st, "isLoopSub", false)
		isLoopStandard := styleBool(st, "isLoopStandard", false)
		isLoopMultiParallel := styleBool(st, "isLoopMultiParallel", false)
		isLoopMultiSeq := styleBool(st, "isLoopMultiSeq", false)
		isLoopComp := styleBool(st, "isLoopComp", false)
		isAdHoc := styleBool(st, "isAdHoc", false)
		loopnum := 0
		for _, b := range []bool{isLoopStandard, isLoopMultiParallel, isLoopMultiSeq, isLoopComp, isLoopSub, isAdHoc} {
			if b {
				loopnum++
			}
		}
		const iconSpaceX = 14
		currXOffset := -iconSpaceX * float64(loopnum) * 0.5
		if v2 {
			c.setStrokeWidth(1)
		}
		if isLoopStandard && bpmnStencil(s, c, "mxgraph.bpmn.loop", w*0.5+currXOffset+1, h-offsetY+1, 12, 12, v2) {
			currXOffset += iconSpaceX
		}
		if isLoopMultiParallel {
			paintBpmnLines(c, w*0.5+currXOffset+1, h-offsetY+1, false)
			c.translate(-w*0.5-currXOffset-1, offsetY-1-h)
			currXOffset += iconSpaceX
		}
		if isLoopMultiSeq {
			paintBpmnLines(c, w*0.5+currXOffset+1, h-offsetY+1, true)
			c.translate(-w*0.5-currXOffset-1, offsetY-1-h)
			currXOffset += iconSpaceX
		}
		if isLoopComp && bpmnStencil(s, c, "mxgraph.bpmn.compensation", w*0.5+currXOffset, h-offsetY+1, 14, 12, v2) {
			currXOffset += iconSpaceX
		}
		if isLoopSub {
			c.translate(w*0.5+currXOffset, h-offsetY)
			paintBpmnSubprocessMarker(c)
			c.translate(-w*0.5-currXOffset, offsetY-h)
			currXOffset += iconSpaceX
		}
		if isAdHoc && lookupStencil("mxgraph.bpmn.ad_hoc") != nil {
			strokeColor := st.get("strokeColor", "#000000")
			fillColor := st.get("fillColor", "#ffffff")
			c.setStrokeColor("")
			c.setFillColor(strokeColor)
			bpmnStencil(s, c, "mxgraph.bpmn.ad_hoc", w*0.5+currXOffset+1, h-offsetY+4, 12, 6, v2)
			currXOffset += iconSpaceX
			c.setStrokeColor(strokeColor)
			c.setFillColor(fillColor)
			if v2 {
				c.setStrokeWidth(1)
			}
		}

		inlet := 0.0
		if v2 {
			inlet = strokeWidth * 0.5
			if bpmnShapeType == "call" {
				inlet = strokeWidth * 2
			}
		}
		switch taskMarker {
		case "service":
			c.setFillColor(st.get("fillColor", "#ffffff"))
			bpmnStencil(s, c, "mxgraph.bpmn.service_task", inlet+2, inlet+2, 16, 16, v2)
		case "send":
			strokeColor := st.get("strokeColor", "#000000")
			fillColor := st.get("fillColor", "#ffffff")
			c.setStrokeColor(fillColor)
			c.setFillColor(strokeColor)
			paintBpmnSendMarker(c, inlet+4, inlet+4, 18, 13)
		case "receive":
			paintBpmnSendMarker(c, inlet+4, inlet+4, 18, 13)
		case "user":
			bpmnStencil(s, c, "mxgraph.bpmn.user_task", inlet+2, inlet+2, 16, 16, v2)
		case "manual":
			bpmnStencil(s, c, "mxgraph.bpmn.manual_task", inlet+3, inlet+3, 18, 14, v2)
		case "businessRule":
			bpmnStencil(s, c, "mxgraph.bpmn.business_rule_task", inlet+4, inlet+4, 18, 14, v2)
		case "script":
			bpmnStencil(s, c, "mxgraph.bpmn.script_task", inlet+3, inlet+3, 19, 18, v2)
		}

		const symbolW, symbolH = 20, 20
		paintBpmnEvent(s, c, inlet, inlet, symbolW, symbolH, st.get("outline", "none"), st.get("symbol", "standard"), v2)
	}}
}

// --- data, swimlane, conversation ---

// bpmnDataShape ports mxShapeBpmn2Data ("mxgraph.bpmn.data", legacy) and,
// with v2, mxShapeBpmn2Data2 ("mxgraph.bpmn.data2"), which insets the
// markers by the stroke width and draws them with a width of 1: a note
// with an arrow for bpmnTransferType=input (outlined) or output (filled)
// and lines at the bottom for isCollection. They keep the label margins
// of the note.
func bpmnDataShape(v2 bool) *shapeDef {
	return &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			paintNote(s, c, x, y, w, h)
			st := s.style
			trType := st.get("bpmnTransferType", "none")
			isColl := styleBool(st, "isCollection", false)
			strokeWidth := st.num("strokeWidth", 1)
			if v2 {
				c.setStrokeWidth(1)
			}
			if trType == "input" || trType == "output" {
				arrX, arrY := 3.0, 3.0
				if v2 {
					arrX, arrY = 3+strokeWidth*0.5, 3+strokeWidth*0.5
				}
				const arrW, arrH = 14, 12
				c.translate(arrX, arrY)
				c.begin()
				c.moveTo(0, arrH*0.3)
				c.lineTo(arrW*0.55, arrH*0.3)
				c.lineTo(arrW*0.55, 0)
				c.lineTo(arrW, arrH*0.5)
				c.lineTo(arrW*0.55, arrH)
				c.lineTo(arrW*0.55, arrH*0.7)
				c.lineTo(0, arrH*0.7)
				c.close()
				c.translate(-arrX, -arrY)
				if trType == "input" {
					c.stroke()
				} else {
					fillColor := st.get("fillColor", "#ffffff")
					c.setFillColor(st.get("strokeColor", "#000000"))
					c.fillAndStroke()
					c.setFillColor(fillColor)
				}
			}
			if isColl {
				dy := 0.0
				if v2 {
					dy = strokeWidth * 0.5
				}
				paintBpmnLines(c, w*0.5-6, h-12-dy, false)
				c.translate(-w*0.5+6, -h+12)
			}
		},
		labelMargins: cylinderLabelMargins,
	}
}

// bpmnSwimlaneShape ports mxShapeBpmn2Swimlane ("mxgraph.bpmn.swimlane"):
// a swimlane with the lines of a collection (isCollection) at the bottom.
var bpmnSwimlaneShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		paintSwimlaneShape(s, c, x, y, w, h, true)
		if styleBool(s.style, "isCollection", false) {
			paintBpmnLines(c, w*0.5-6, h-12, false)
			c.translate(-w*0.5+6, -h+12)
		}
	},
	labelBounds: func(s *shape, r rect) rect { return swimlaneLabelBounds(s, r, true) },
	roundable:   true,
}

// bpmnConversationShape ports mxShapeBpmn2Conversation
// ("mxgraph.bpmn.conversation", legacy) and, with v2,
// mxShapeBpmn2Conversation2 ("mxgraph.bpmn.conversation2"), which insets
// the marker by the stroke width and draws it with a width of 1: a
// hexagon, thick for bpmnConversationType=call, with the sub-process
// marker for isLoopSub.
func bpmnConversationShape(v2 bool) *shapeDef {
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		call := s.style.get("bpmnConversationType", "conv") == "call"
		strokeWidth := s.style.num("strokeWidth", 1)
		if call {
			c.setStrokeWidth(strokeWidth * 4)
		}
		c.translate(x, y)
		c.begin()
		c.moveTo(0, h*0.5)
		c.lineTo(w*0.25, 0)
		c.lineTo(w*0.75, 0)
		c.lineTo(w, h*0.5)
		c.lineTo(w*0.75, h)
		c.lineTo(w*0.25, h)
		c.close()
		c.fillAndStroke()
		if v2 {
			c.setStrokeWidth(1)
		} else if call {
			c.setStrokeWidth(strokeWidth)
		}
		if styleBool(s.style, "isLoopSub", false) {
			switch {
			case !v2:
				c.translate(w*0.5-7, h-14)
			case call:
				c.translate(w*0.5-7, h-14-strokeWidth*2)
			default:
				c.translate(w*0.5-7, h-14-strokeWidth*0.5)
			}
			paintBpmnSubprocessMarker(c)
		}
	}}
}

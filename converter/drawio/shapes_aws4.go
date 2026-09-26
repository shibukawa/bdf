package drawio

import "math"

// The AWS shapes of draw.io's shapes/mxAWS4.js (mxgraph.aws4.*): the icons
// of the AWS palette are squares in the category color with a stencil of
// the aws4 library inside (resourceIcon, and the older productIcon with a
// frame and a label area below), and the group frames draw a stencil in
// their top left corner (group, and group2 on a square of the stroke
// color) or at the top center (groupCenter). The stencils are looked up by
// the names in resIcon, prIcon and grIcon and drawn with the shape itself,
// as mxStencil.drawShape(c, this, ...) does: filled with the stroke color
// (white for group2) and not stroked. The colors are read from the style
// as the JavaScript reads this.state.style, where "none" has removed a key
// and the default of the call applies. Their setShadow(false) calls are
// left out: the shadow is drawn for the whole shape (see shape.paint).
// The file registers no perimeters.

func init() {
	registerShape("mxgraph.aws4.productIcon", aws4ProductIconShape)
	registerShape("mxgraph.aws4.resourceIcon", aws4ResourceIconShape)
	registerShape("mxgraph.aws4.group", aws4GroupShape(false))
	registerShape("mxgraph.aws4.groupCenter", aws4GroupShape(true))
	registerShape("mxgraph.aws4.group2", aws4Group2Shape)
}

// aws4Rect adds the rectangle 0, 0, w, h as a path, as the shapes draw it.
func aws4Rect(c *c2d, w, h float64) {
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.close()
}

// aws4ProductIconShape ports mxShapeAws4ProductIcon
// ("mxgraph.aws4.productIcon"): the bounds filled with the stroke color
// (none by default), a square in the fill color (white by default) or a
// gradient to gradientColor, inset by 1 and as high as the bounds are
// wide, and the prIcon stencil filled with the stroke color in the middle
// of that square. The label is placed in the bounds as usual (the palette
// puts it below the square with verticalLabelPosition=middle and
// verticalAlign=bottom).
var aws4ProductIconShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	// op1 and op2 are meant to be 0 for fillColor=none and
	// gradientColor=none, but the JavaScript tests the variables before
	// it reads the colors, so both stops get the opacity; SVG clamps a
	// stop-opacity to 0..1 and ignores one that is not a number.
	op := jsStyleFloat(s.style, "opacity", 100)
	if math.IsNaN(op) {
		op = 1
	}
	op = clamp(op, 0, 1)
	const ind = 1
	strokeColor := s.style.get("strokeColor", "none")
	c.setFillColor(strokeColor)
	aws4Rect(c, w, h)
	c.fill()

	fillColor := s.style.get("fillColor", "#ffffff")
	gradientColor := s.style.get("gradientColor", fillColor)
	gradientDir := s.style.get("gradientDirection", "south")
	c.setFillColor(fillColor)
	if gradientColor != fillColor || op != 1 {
		// (a gradient between equal colors and opaque stops is the fill color)
		c.setGradient(fillColor, gradientColor, 0, 0, w, h, gradientDir, op, op)
	}
	c.begin()
	c.moveTo(ind, ind)
	c.lineTo(w-ind, ind)
	c.lineTo(w-ind, w-ind)
	c.lineTo(ind, w-ind)
	c.close()
	c.fill()

	if st := lookupStencil(s.style.get("prIcon", "")); st != nil {
		c.setFillColor(strokeColor)
		c.setStrokeColor("")
		st.drawShape(c, s, ind+w*0.15, ind+w*0.15, w*0.7-2*ind, w*0.7-2*ind)
	}
}}

// aws4ResourceIconShape ports mxShapeAws4ResourceIcon
// ("mxgraph.aws4.resourceIcon"): the bounds filled (not stroked) with the
// fill color or gradient, and the resIcon stencil filled with the stroke
// color (black by default) in the middle 80%.
var aws4ResourceIconShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	aws4Rect(c, w, h)
	c.fill()
	if st := lookupStencil(s.style.get("resIcon", "")); st != nil {
		c.setFillColor(s.style.get("strokeColor", "#000000"))
		c.setStrokeColor("")
		st.drawShape(c, s, w*0.1, h*0.1, w*0.8, h*0.8)
	}
}}

// aws4IconSize is the size of a group's icon, grIconSize (25 by default)
// as JavaScript's arithmetic reads the style value: NaN for a value that
// is not a number. An icon whose size is not a finite number is not
// drawn, as the NaN coordinates that draw.io writes into the SVG path.
func aws4IconSize(s *shape) float64 {
	if v, ok := s.style["grIconSize"]; ok {
		return numberJS(v)
	}
	return 25
}

// aws4GroupShape ports mxShapeAws4Group ("mxgraph.aws4.group") and, with
// center, mxShapeAws4GroupCenter ("mxgraph.aws4.groupCenter"): the frame
// of the bounds, stroked unless grStroke is not 1, and the grIcon stencil
// filled with the stroke color (black by default) in a grIconSize square
// at the top left corner, or at the top center. The pointer events that
// mxShapeAws4Group turns off for pointerEvents=0 draw nothing.
func aws4GroupShape(center bool) *shapeDef {
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		grStroke := s.style.get("grStroke", "1")
		aws4Rect(c, w, h)
		if grStroke == "1" || s.outline {
			c.fillAndStroke()
		} else {
			c.fill()
		}
		if st := lookupStencil(s.style.get("grIcon", "")); st != nil {
			c.setFillAlpha(s.strokeOpacity / 100)
			c.setFillColor(s.style.get("strokeColor", "#000000"))
			c.setStrokeColor("")
			size := aws4IconSize(s)
			if !finite(size) {
				return
			}
			x0 := 0.0
			if center {
				x0 = (w - size) * 0.5
			}
			st.drawShape(c, s, x0, 0, size, size)
		}
	}}
}

// aws4Group2Shape ports mxShapeAws4Group2 ("mxgraph.aws4.group2"): the
// frame of the bounds, a 25 by 25 square of the stroke color (black by
// default) at the top left corner, and the grIcon stencil filled with
// white in the middle 80% of a grIconSize square there.
var aws4Group2Shape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.translate(x, y)
	const size = 25
	aws4Rect(c, w, h)
	c.fillAndStroke()
	c.setFillColor(s.style.get("strokeColor", "#000000"))
	aws4Rect(c, size, size)
	c.fill()
	if st := lookupStencil(s.style.get("grIcon", "")); st != nil {
		c.setFillAlpha(s.strokeOpacity / 100)
		c.setFillColor("#ffffff")
		c.setStrokeColor("")
		size := aws4IconSize(s)
		if !finite(size) {
			return
		}
		st.drawShape(c, s, size*0.1, size*0.1, size*0.8, size*0.8)
	}
}}

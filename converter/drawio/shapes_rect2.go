package drawio

import "math"

// mxShapeBasicRect2 ("mxgraph.basic.rect", Shapes.js): a rectangle whose
// corners are square, rounded, snipped, inverted round or folded
// (rectStyle, or topLeftStyle and the like per corner), whose sides can be
// left out (top, right, bottom, left), with a single, double or frame
// outline (rectOutline, indent) and an inside fill (fillColor2,
// gradientColor2).

func init() {
	registerShape("mxgraph.basic.rect", &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		paintRect2(s, c, w, h)
	}})
}

// rect2 holds what the painting functions of mxShapeBasicRect2 share.
type rect2 struct {
	c                        *c2d
	w, h                     float64
	rectStyle                string
	tl, tr, br, bl           string // corner styles
	size, indent             float64
	top, right, bottom, left bool
}

// is reports whether a corner style (or the default rectStyle) is style.
func (r *rect2) is(corner, style string) bool {
	return corner == style || corner == "default" && r.rectStyle == style
}

// paintRect2 is mxShapeBasicRect2.strictDrawShape.
func paintRect2(s *shape, c *c2d, w, h float64) {
	st := s.style
	r := &rect2{c: c, w: w, h: h,
		rectStyle: st.get("rectStyle", "square"),
		tl:        st.get("topLeftStyle", "default"), tr: st.get("topRightStyle", "default"),
		br: st.get("bottomRightStyle", "default"), bl: st.get("bottomLeftStyle", "default"),
		top: styleBool(st, "top", true), right: styleBool(st, "right", true),
		bottom: styleBool(st, "bottom", true), left: styleBool(st, "left", true),
	}
	absoluteCornerSize := styleBool(st, "absoluteCornerSize", true)
	size := clamp(st.num("size", 10), 0, w)
	rectOutline := st.get("rectOutline", "single")
	indent := clamp(st.num("indent", 2), 0, w)
	relIndent := clamp(indent, 0, 50)
	fillColor := st.get("fillColor", "#ffffff")
	fillColor2 := st.get("fillColor2", "none")
	gradientColor2 := st.get("gradientColor2", "none")
	gdir2 := st.get("gradientDirection2", "south")
	opacity := st.num("opacity", 100)
	relSize := clamp(size, 0, 50)

	c.setDashed(styleBool(st, "dashed", false), false)
	if dp := st.get("dashPattern", ""); dp != "" {
		c.setDashPattern(dp)
	}
	c.setStrokeWidth(st.num("strokeWidth", 1))
	size = math.Min(h*0.5, math.Min(w*0.5, size))
	if !absoluteCornerSize {
		size = relSize * math.Min(w, h) / 100
	}
	size = math.Min(size, math.Min(w, h)*0.5)
	if !absoluteCornerSize {
		indent = relIndent * math.Min(w, h) / 100
	}
	indent = math.Min(indent, math.Min(w, h)*0.5-size)
	r.size, r.indent = size, indent
	top, right, bottom, left := r.top, r.right, r.bottom, r.left

	if (top || right || bottom || left) && rectOutline != "frame" {
		// outline fill
		c.begin()
		if !top {
			c.moveTo(0, 0)
		} else {
			r.moveNW()
		}
		if top {
			r.paintNW()
		}
		r.paintTop()
		if right {
			r.paintNE()
		}
		r.paintRight()
		if bottom {
			r.paintSE()
		}
		r.paintBottom()
		if left {
			r.paintSW()
		}
		r.paintLeft()
		c.close()
		c.fill()

		// inner fill: the gradient's stop opacities are the opacity style,
		// which SVG clamps to 1
		op1, op2 := clamp(opacity, 0, 1), clamp(opacity, 0, 1)
		if fillColor2 == "none" {
			op1 = 0
		}
		if gradientColor2 == "none" {
			op2 = 0
			c.setFillColor(fillColor2)
		} else {
			c.setGradient(fillColor2, gradientColor2, 0, 0, w, h, gdir2, op1, op2)
		}
		c.begin()
		if !top {
			c.moveTo(indent, 0)
		} else {
			r.moveNWInner()
		}
		r.paintLeftInner()
		if left && bottom {
			r.paintSWInner()
		}
		r.paintBottomInner()
		if bottom && right {
			r.paintSEInner()
		}
		r.paintRightInner()
		if right && top {
			r.paintNEInner()
		}
		r.paintTopInner()
		if top && left {
			r.paintNWInner()
		}
		c.fill()
		if fillColor == "none" {
			c.begin()
			r.paintFolds()
			c.stroke()
		}
	}

	// the outline of every combination of sides
	double := rectOutline == "double"
	frame := rectOutline == "frame"
	// finish strokes a single or double outline, or closes and fills a frame
	finish := func() {
		if frame {
			c.close()
			c.fillAndStroke()
		} else {
			c.stroke()
		}
	}
	switch {
	case !top && !right && !bottom && left:
		c.begin()
		r.moveSW()
		r.paintLeft()
		if frame {
			r.lineNWInner()
			r.paintLeftInner()
		} else if double {
			r.moveNWInner()
			r.paintLeftInner()
		}
		finish()
	case !top && !right && bottom && !left:
		c.begin()
		r.moveSE()
		r.paintBottom()
		if frame {
			r.lineSWInner()
			r.paintBottomInner()
		} else if double {
			r.moveSWInner()
			r.paintBottomInner()
		}
		finish()
	case !top && !right && bottom && left:
		c.begin()
		r.moveSE()
		r.paintBottom()
		r.paintSW()
		r.paintLeft()
		if frame {
			r.lineNWInner()
		} else if double {
			r.moveNWInner()
		}
		if frame || double {
			r.paintLeftInner()
			r.paintSWInner()
			r.paintBottomInner()
		}
		finish()
	case !top && right && !bottom && !left:
		c.begin()
		r.moveNE()
		r.paintRight()
		if frame {
			r.lineSEInner()
			r.paintRightInner()
		} else if double {
			r.moveSEInner()
			r.paintRightInner()
		}
		finish()
	case !top && right && !bottom && left:
		c.begin()
		r.moveSW()
		r.paintLeft()
		if frame {
			r.lineNWInner()
			r.paintLeftInner()
		} else if double {
			r.moveNWInner()
			r.paintLeftInner()
		}
		finish()
		c.begin()
		r.moveNE()
		r.paintRight()
		if frame {
			r.lineSEInner()
			r.paintRightInner()
		} else if double {
			r.moveSEInner()
			r.paintRightInner()
		}
		finish()
	case !top && right && bottom && !left:
		c.begin()
		r.moveNE()
		r.paintRight()
		r.paintSE()
		r.paintBottom()
		if frame {
			r.lineSWInner()
		} else if double {
			r.moveSWInner()
		}
		if frame || double {
			r.paintBottomInner()
			r.paintSEInner()
			r.paintRightInner()
		}
		finish()
	case !top && right && bottom && left:
		c.begin()
		r.moveNE()
		r.paintRight()
		r.paintSE()
		r.paintBottom()
		r.paintSW()
		r.paintLeft()
		if frame {
			r.lineNWInner()
		} else if double {
			r.moveNWInner()
		}
		if frame || double {
			r.paintLeftInner()
			r.paintSWInner()
			r.paintBottomInner()
			r.paintSEInner()
			r.paintRightInner()
		}
		finish()
	case top && !right && !bottom && !left:
		c.begin()
		r.moveNW()
		r.paintTop()
		if frame {
			r.lineNEInner()
			r.paintTopInner()
		} else if double {
			r.moveNEInner()
			r.paintTopInner()
		}
		finish()
	case top && !right && !bottom && left:
		c.begin()
		r.moveSW()
		r.paintLeft()
		r.paintNW()
		r.paintTop()
		if frame {
			r.lineNEInner()
		} else if double {
			r.moveNEInner()
		}
		if frame || double {
			r.paintTopInner()
			r.paintNWInner()
			r.paintLeftInner()
		}
		finish()
	case top && !right && bottom && !left:
		c.begin()
		r.moveNW()
		r.paintTop()
		if frame {
			r.lineNEInner()
			r.paintTopInner()
		} else if double {
			r.moveNEInner()
			r.paintTopInner()
		}
		finish()
		c.begin()
		r.moveSE()
		r.paintBottom()
		if frame {
			r.lineSWInner()
			r.paintBottomInner()
		} else if double {
			r.moveSWInner()
			r.paintBottomInner()
		}
		finish()
	case top && !right && bottom && left:
		c.begin()
		r.moveSE()
		r.paintBottom()
		r.paintSW()
		r.paintLeft()
		r.paintNW()
		r.paintTop()
		if frame {
			r.lineNEInner()
		} else if double {
			r.moveNEInner()
		}
		if frame || double {
			r.paintTopInner()
			r.paintNWInner()
			r.paintLeftInner()
			r.paintSWInner()
			r.paintBottomInner()
		}
		finish()
	case top && right && !bottom && !left:
		c.begin()
		r.moveNW()
		r.paintTop()
		r.paintNE()
		r.paintRight()
		if frame {
			r.lineSEInner()
		} else if double {
			r.moveSEInner()
		}
		if frame || double {
			r.paintRightInner()
			r.paintNEInner()
			r.paintTopInner()
		}
		finish()
	case top && right && !bottom && left:
		c.begin()
		r.moveSW()
		r.paintLeft()
		r.paintNW()
		r.paintTop()
		r.paintNE()
		r.paintRight()
		if frame {
			r.lineSEInner()
		} else if double {
			r.moveSEInner()
		}
		if frame || double {
			r.paintRightInner()
			r.paintNEInner()
			r.paintTopInner()
			r.paintNWInner()
			r.paintLeftInner()
		}
		finish()
	case top && right && bottom && !left:
		c.begin()
		r.moveNW()
		r.paintTop()
		r.paintNE()
		r.paintRight()
		r.paintSE()
		r.paintBottom()
		if frame {
			r.lineSWInner()
		} else if double {
			r.moveSWInner()
		}
		if frame || double {
			r.paintBottomInner()
			r.paintSEInner()
			r.paintRightInner()
			r.paintNEInner()
			r.paintTopInner()
		}
		finish()
	case top && right && bottom && left:
		c.begin()
		r.moveNW()
		r.paintNW()
		r.paintTop()
		r.paintNE()
		r.paintRight()
		r.paintSE()
		r.paintBottom()
		r.paintSW()
		r.paintLeft()
		c.close()
		if frame || double {
			r.moveSWInner()
			r.paintSWInner()
			r.paintBottomInner()
			r.paintSEInner()
			r.paintRightInner()
			r.paintNEInner()
			r.paintTopInner()
			r.paintNWInner()
			r.paintLeftInner()
			c.close()
		}
		if frame {
			c.fillAndStroke()
		} else {
			c.stroke()
		}
	}
	c.begin()
	r.paintFolds()
	c.stroke()
}

// The outline: moves to a corner and lines along a side or around a corner.

func (r *rect2) moveNW() {
	if r.is(r.tl, "square") || !r.left {
		r.c.moveTo(0, 0)
	} else {
		r.c.moveTo(0, r.size)
	}
}

func (r *rect2) moveNE() {
	if r.is(r.tr, "square") || !r.top {
		r.c.moveTo(r.w, 0)
	} else {
		r.c.moveTo(r.w-r.size, 0)
	}
}

func (r *rect2) moveSE() {
	if r.is(r.br, "square") || !r.right {
		r.c.moveTo(r.w, r.h)
	} else {
		r.c.moveTo(r.w, r.h-r.size)
	}
}

func (r *rect2) moveSW() {
	if r.is(r.bl, "square") || !r.bottom {
		r.c.moveTo(0, r.h)
	} else {
		r.c.moveTo(r.size, r.h)
	}
}

// corner adds the corner curve ending at (x, y) for a corner style.
func (r *rect2) corner(style string, x, y float64) {
	switch {
	case r.is(style, "rounded"):
		r.c.arcTo(r.size, r.size, 0, false, true, x, y)
	case r.is(style, "invRound"):
		r.c.arcTo(r.size, r.size, 0, false, false, x, y)
	case r.is(style, "snip"), r.is(style, "fold"):
		r.c.lineTo(x, y)
	}
}

func (r *rect2) paintNW() {
	if !r.left {
		r.c.lineTo(0, 0)
	} else {
		r.corner(r.tl, r.size, 0)
	}
}

func (r *rect2) paintTop() {
	if r.is(r.tr, "square") || !r.right {
		r.c.lineTo(r.w, 0)
	} else {
		r.c.lineTo(r.w-r.size, 0)
	}
}

func (r *rect2) paintNE() {
	if !r.top {
		r.c.lineTo(r.w, 0)
	} else {
		r.corner(r.tr, r.w, r.size)
	}
}

func (r *rect2) paintRight() {
	if r.is(r.br, "square") || !r.bottom {
		r.c.lineTo(r.w, r.h)
	} else {
		r.c.lineTo(r.w, r.h-r.size)
	}
}

func (r *rect2) paintLeft() {
	if r.is(r.tl, "square") || !r.top {
		r.c.lineTo(0, 0)
	} else {
		r.c.lineTo(0, r.size)
	}
}

func (r *rect2) paintSE() {
	if !r.right {
		r.c.lineTo(r.w, r.h)
	} else {
		r.corner(r.br, r.w-r.size, r.h)
	}
}

func (r *rect2) paintBottom() {
	if r.is(r.bl, "square") || !r.left {
		r.c.lineTo(0, r.h)
	} else {
		r.c.lineTo(r.size, r.h)
	}
}

func (r *rect2) paintSW() {
	if !r.bottom {
		r.c.lineTo(0, r.h)
	} else {
		r.corner(r.bl, 0, r.h-r.size)
	}
}

// The inner outline (double and frame) and the inside fill.

func (r *rect2) paintNWInner() {
	c, size, indent := r.c, r.size, r.indent
	switch {
	case r.is(r.tl, "rounded"):
		c.arcTo(size-indent*0.5, size-indent*0.5, 0, false, false, indent, indent*0.5+size)
	case r.is(r.tl, "invRound"):
		c.arcTo(size+indent, size+indent, 0, false, true, indent, indent+size)
	case r.is(r.tl, "snip"):
		c.lineTo(indent, indent*0.5+size)
	case r.is(r.tl, "fold"):
		c.lineTo(indent+size, indent+size)
		c.lineTo(indent, indent+size)
	}
}

func (r *rect2) paintTopInner() {
	c, size, indent := r.c, r.size, r.indent
	switch {
	case !r.left && !r.top:
		c.lineTo(0, 0)
	case !r.left && r.top:
		c.lineTo(0, indent)
	case r.left && !r.top:
		c.lineTo(indent, 0)
	case r.is(r.tl, "square"):
		c.lineTo(indent, indent)
	case r.is(r.tl, "rounded"), r.is(r.tl, "snip"):
		c.lineTo(size+indent*0.5, indent)
	default:
		c.lineTo(size+indent, indent)
	}
}

func (r *rect2) paintNEInner() {
	c, w, size, indent := r.c, r.w, r.size, r.indent
	switch {
	case r.is(r.tr, "rounded"):
		c.arcTo(size-indent*0.5, size-indent*0.5, 0, false, false, w-size-indent*0.5, indent)
	case r.is(r.tr, "invRound"):
		c.arcTo(size+indent, size+indent, 0, false, true, w-size-indent, indent)
	case r.is(r.tr, "snip"):
		c.lineTo(w-size-indent*0.5, indent)
	case r.is(r.tr, "fold"):
		c.lineTo(w-size-indent, size+indent)
		c.lineTo(w-size-indent, indent)
	}
}

func (r *rect2) paintRightInner() {
	c, w, size, indent := r.c, r.w, r.size, r.indent
	switch {
	case !r.top && !r.right:
		c.lineTo(w, 0)
	case !r.top && r.right:
		c.lineTo(w-indent, 0)
	case r.top && !r.right:
		c.lineTo(w, indent)
	case r.is(r.tr, "square"):
		c.lineTo(w-indent, indent)
	case r.is(r.tr, "rounded"), r.is(r.tr, "snip"):
		c.lineTo(w-indent, size+indent*0.5)
	default:
		c.lineTo(w-indent, size+indent)
	}
}

func (r *rect2) paintLeftInner() {
	c, h, size, indent := r.c, r.h, r.size, r.indent
	switch {
	case !r.bottom && !r.left:
		c.lineTo(0, h)
	case !r.bottom && r.left:
		c.lineTo(indent, h)
	case r.bottom && !r.left:
		c.lineTo(0, h-indent)
	case r.is(r.bl, "square"):
		c.lineTo(indent, h-indent)
	case r.is(r.bl, "rounded"), r.is(r.bl, "snip"):
		c.lineTo(indent, h-size-indent*0.5)
	default:
		c.lineTo(indent, h-size-indent)
	}
}

func (r *rect2) paintSEInner() {
	c, w, h, size, indent := r.c, r.w, r.h, r.size, r.indent
	switch {
	case r.is(r.br, "rounded"):
		c.arcTo(size-indent*0.5, size-indent*0.5, 0, false, false, w-indent, h-size-indent*0.5)
	case r.is(r.br, "invRound"):
		c.arcTo(size+indent, size+indent, 0, false, true, w-indent, h-size-indent)
	case r.is(r.br, "snip"):
		c.lineTo(w-indent, h-size-indent*0.5)
	case r.is(r.br, "fold"):
		c.lineTo(w-size-indent, h-size-indent)
		c.lineTo(w-indent, h-size-indent)
	}
}

func (r *rect2) paintBottomInner() {
	c, w, h, size, indent := r.c, r.w, r.h, r.size, r.indent
	switch {
	case !r.right && !r.bottom:
		c.lineTo(w, h)
	case !r.right && r.bottom:
		c.lineTo(w, h-indent)
	case r.right && !r.bottom:
		c.lineTo(w-indent, h)
	case r.is(r.br, "square"):
		c.lineTo(w-indent, h-indent)
	case r.is(r.br, "rounded"), r.is(r.br, "snip"):
		c.lineTo(w-size-indent*0.5, h-indent)
	default:
		c.lineTo(w-size-indent, h-indent)
	}
}

func (r *rect2) paintSWInner() {
	c, h, size, indent := r.c, r.h, r.size, r.indent
	switch {
	case !r.bottom:
		c.lineTo(indent, h)
	case r.is(r.bl, "square"):
		c.lineTo(indent, h-indent)
	case r.is(r.bl, "rounded"):
		c.arcTo(size-indent*0.5, size-indent*0.5, 0, false, false, size+indent*0.5, h-indent)
	case r.is(r.bl, "invRound"):
		c.arcTo(size+indent, size+indent, 0, false, true, size+indent, h-indent)
	case r.is(r.bl, "snip"):
		c.lineTo(size+indent*0.5, h-indent)
	case r.is(r.bl, "fold"):
		c.lineTo(indent+size, h-size-indent)
		c.lineTo(indent+size, h-indent)
	}
}

// innerSW, innerSE, innerNE and innerNW are the points where the inner
// outline turns at a corner; to moves or lines to them.

func (r *rect2) innerSW(to func(x, y float64)) {
	h, size, indent := r.h, r.size, r.indent
	switch {
	case !r.left:
		to(0, h-indent)
	case r.is(r.bl, "square"):
		to(indent, h-indent)
	case r.is(r.bl, "rounded"), r.is(r.bl, "snip"):
		to(indent, h-size-indent*0.5)
	case r.is(r.bl, "invRound"), r.is(r.bl, "fold"):
		to(indent, h-size-indent)
	}
}

func (r *rect2) innerSE(to func(x, y float64)) {
	w, h, size, indent := r.w, r.h, r.size, r.indent
	switch {
	case !r.bottom:
		to(w-indent, h)
	case r.is(r.br, "square"):
		to(w-indent, h-indent)
	case r.is(r.br, "rounded"), r.is(r.br, "snip"):
		to(w-indent, h-size-indent*0.5)
	case r.is(r.br, "invRound"), r.is(r.br, "fold"):
		to(w-indent, h-size-indent)
	}
}

func (r *rect2) innerNE(to func(x, y float64)) {
	// the square case also catches every right side
	if !r.right {
		to(r.w, r.indent)
	} else {
		to(r.w-r.indent, r.indent)
	}
}

func (r *rect2) innerNW(to func(x, y float64)) {
	size, indent := r.size, r.indent
	switch {
	case !r.top && !r.left:
		to(0, 0)
	case !r.top && r.left:
		to(indent, 0)
	case r.top && !r.left:
		to(0, indent)
	case r.is(r.tl, "square"):
		to(indent, indent)
	case r.is(r.tl, "rounded"), r.is(r.tl, "snip"):
		to(indent, size+indent*0.5)
	case r.is(r.tl, "invRound"), r.is(r.tl, "fold"):
		to(indent, size+indent)
	}
}

func (r *rect2) moveSWInner() { r.innerSW(r.c.moveTo) }
func (r *rect2) lineSWInner() { r.innerSW(r.c.lineTo) }
func (r *rect2) moveSEInner() { r.innerSE(r.c.moveTo) }
func (r *rect2) lineSEInner() { r.innerSE(r.c.lineTo) }
func (r *rect2) moveNEInner() { r.innerNE(r.c.moveTo) }
func (r *rect2) lineNEInner() { r.innerNE(r.c.lineTo) }
func (r *rect2) moveNWInner() { r.innerNW(r.c.moveTo) }
func (r *rect2) lineNWInner() { r.innerNW(r.c.lineTo) }

// paintFolds adds the creases of folded corners.
func (r *rect2) paintFolds() {
	if r.rectStyle != "fold" && r.tl != "fold" && r.tr != "fold" && r.br != "fold" && r.bl != "fold" {
		return
	}
	c, w, h, size := r.c, r.w, r.h, r.size
	if r.is(r.tl, "fold") && r.top && r.left {
		c.moveTo(0, size)
		c.lineTo(size, size)
		c.lineTo(size, 0)
	}
	if r.is(r.tr, "fold") && r.top && r.right {
		c.moveTo(w-size, 0)
		c.lineTo(w-size, size)
		c.lineTo(w, size)
	}
	if r.is(r.br, "fold") && r.bottom && r.right {
		c.moveTo(w-size, h)
		c.lineTo(w-size, h-size)
		c.lineTo(w, h-size)
	}
	if r.is(r.bl, "fold") && r.bottom && r.left {
		c.moveTo(0, h-size)
		c.lineTo(size, h-size)
		c.lineTo(size, h)
	}
}

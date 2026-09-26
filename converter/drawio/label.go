package drawio

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// Labels: where draw.io puts a cell's text (mxCellRenderer.getLabelBounds,
// rotateLabelBounds and mxText.getSpacing) and how the text box sits at
// that position (mxSvgCanvas2D.createCss for HTML labels, plainText for
// others).

// labelBox is a laid-out label.
type labelBox struct {
	lo       *tlayout
	trailing []paraMark // structure marks after the last paragraph
	// anchor is the label position (the rotation center); rotation is in degrees.
	anchor   point
	rotation float64
	// left and top place the content box; lines are relative to it.
	left, top float64
	boxW      float64 // width lines are aligned in
	clip      *rect
	bg        *rgba
	border    *rgba
	bgBox     rect
	alpha     float64
	plain     bool
	pad       []float64 // labelPadding: top, right, bottom, left
}

// textProps are the label properties of mxText (as it reads them from the style).
type textProps struct {
	align, valign              string
	spacing                    int
	spacingTop, spacingRight   int
	spacingBottom, spacingLeft int
	horizontal                 bool
	html, wrap, clipped        bool
	overflow                   string
}

func readTextProps(s style) textProps {
	t := textProps{
		align:      s.get("align", "center"),
		valign:     s.get("verticalAlign", "middle"),
		horizontal: s.get("horizontal", "1") != "0" && s.get("horizontal", "1") != "false",
		overflow:   s.get("overflow", "visible"),
	}
	t.spacing = parseInt(s.get("spacing", ""), 2)
	t.spacingTop = t.spacing + parseInt(s.get("spacingTop", ""), 0)
	t.spacingRight = t.spacing + parseInt(s.get("spacingRight", ""), 0)
	t.spacingBottom = t.spacing + parseInt(s.get("spacingBottom", ""), 0)
	t.spacingLeft = t.spacing + parseInt(s.get("spacingLeft", ""), 0)
	t.html = s.get("html", "") == "1" || s.get("whiteSpace", "") == "wrap"
	t.wrap = t.html && s.get("whiteSpace", "") == "wrap"
	t.clipped = t.overflow == "hidden"
	return t
}

// parseInt reads a number the way JavaScript's parseInt does (truncating);
// def when s is empty or not a number.
func parseInt(s string, def int) int {
	f, ok := parseFloat(s)
	if !ok {
		return def
	}
	return int(f)
}

// foreignObjectPadding is added to the width of HTML labels (mxSvgCanvas2D).
const foreignObjectPadding = 2

// Base spacings of draw.io's labels (Graph.js sets mxText.baseSpacingTop and baseSpacingBottom).
const (
	baseSpacingTop    = 5
	baseSpacingBottom = 1
)

// getSpacing returns the offset of the label from its anchor for the
// alignment (mxText.getSpacing); mx is the horizontal margin factor (-0.5
// center, -1 right, 0 left).
func (t textProps) getSpacing(noBase bool, mx float64) point {
	var dx, dy float64
	switch mx {
	case -0.5:
		dx = float64(t.spacingLeft-t.spacingRight) / 2
	case -1:
		dx = -float64(t.spacingRight)
	default:
		dx = float64(t.spacingLeft)
	}
	switch t.valign {
	case "middle":
		dy = float64(t.spacingTop-t.spacingBottom) / 2
	case "bottom":
		dy = -float64(t.spacingBottom)
		if !noBase {
			dy -= baseSpacingBottom
		}
	default:
		dy = float64(t.spacingTop)
		if !noBase {
			dy += baseSpacingTop
		}
	}
	return point{dx, dy}
}

// alignmentPoint is mxUtils.getAlignmentAsPoint.
func alignmentPoint(align, valign string) point {
	p := point{-0.5, -0.5}
	switch align {
	case "left":
		p.x = 0
	case "right":
		p.x = -1
	}
	switch valign {
	case "top":
		p.y = 0
	case "bottom":
		p.y = -1
	}
	return p
}

// textRotation is the rotation of a shape's label (mxShape.getTextRotation):
// the shape's rotation, turned -90° for vertical (horizontal=0) labels.
func (s *shape) textRotation(t textProps) float64 {
	rot := 0.0
	if !s.def.noRotation {
		rot = s.rotation
	}
	if !t.horizontal {
		rot -= 90
	}
	return rot
}

// labelBounds lets the shape adjust the label rectangle (mxShape.getLabelBounds).
func (s *shape) labelBounds(r rect, textInverted bool) rect {
	if s.def.labelBounds != nil {
		return s.def.labelBounds(s, r)
	}
	d := s.style.get("direction", "east")
	b := r
	if d != "south" && d != "north" && textInverted {
		b.w, b.h = b.h, b.w
	}
	var m *rect
	if s.stencil != nil {
		m = s.stencil.labelMargins(s, b)
	}
	if m == nil && s.def.labelMargins != nil {
		m = s.def.labelMargins(s, b)
	}
	if m == nil {
		return r
	}
	flipH, flipV := s.style.is("flipH"), s.style.is("flipV")
	if textInverted {
		m.x, m.y, m.w, m.h = m.h, m.x, m.y, m.w
		flipH, flipV = flipV, flipH
	}
	return directedBounds(r, *m, d, flipH, flipV)
}

// directedBounds applies margins (left, top, right, bottom as x, y, w, h)
// for a direction and flips (mxUtils.getDirectedBounds).
func directedBounds(r rect, m rect, d string, flipH, flipV bool) rect {
	m.x = math.Round(math.Max(0, math.Min(r.w, m.x)))
	m.y = math.Round(math.Max(0, math.Min(r.h, m.y)))
	m.w = math.Round(math.Max(0, math.Min(r.w, m.w)))
	m.h = math.Round(math.Max(0, math.Min(r.h, m.h)))
	vert := d == "south" || d == "north"
	horz := d == "east" || d == "west"
	if flipV && vert || flipH && horz {
		m.x, m.w = m.w, m.x
	}
	if flipH && vert || flipV && horz {
		m.y, m.h = m.h, m.y
	}
	m2 := m
	switch d {
	case "south":
		m2.y, m2.x, m2.w, m2.h = m.x, m.h, m.y, m.w
	case "west":
		m2.y, m2.x, m2.w, m2.h = m.h, m.w, m.x, m.y
	case "north":
		m2.y, m2.x, m2.w, m2.h = m.w, m.y, m.h, m.x
	}
	return rect{r.x + m2.x, r.y + m2.y, r.w - m2.w - m2.x, r.h - m2.h - m2.y}
}

// layoutLabel lays out the label of a cell; nil when it has none.
func (c *converter) layoutLabel(st *cellState, s *shape) *labelBox {
	if st.style.is("noLabel") {
		return nil
	}
	value := c.m.label(st.cell, st.style, c.page)
	if value == "" {
		return nil
	}
	fontSize := st.style.num("fontSize", 11)
	if fontSize <= 0 {
		return nil
	}
	t := readTextProps(st.style)
	if dir := st.style.get("textDirection", ""); strings.HasPrefix(dir, "vertical-") {
		c.warnOnce("verticaltext", "vertical text (textDirection=%s) is drawn horizontally", dir)
	}

	// mxCellRenderer.getLabelBounds
	b := rect{st.absoluteOffset.x, st.absoluteOffset.y, 0, 0}
	inverted := !t.horizontal && st.cell.vertex
	if st.cell.edge {
		sp := t.getSpacing(false, alignmentPoint(t.align, t.valign).x)
		b.x += sp.x
		b.y += sp.y
		if g := st.cell.geo; g != nil {
			b.w, b.h = math.Max(0, g.w), math.Max(0, g.h)
		}
	} else {
		if inverted {
			b.x, b.y = b.y, b.x
		}
		b.x += st.x
		b.y += st.y
		b.w, b.h = math.Max(1, st.w), math.Max(1, st.h)
	}
	if inverted {
		tt := (st.w - st.h) / 2
		b.x += tt
		b.y -= tt
		b.w, b.h = b.h, b.w
	}
	hpos := st.style.get("labelPosition", "center")
	vpos := st.style.get("verticalLabelPosition", "middle")
	if hpos == "center" && vpos == "middle" {
		b = s.labelBounds(b, inverted)
	}
	lw, hasLW := st.style["labelWidth"]
	if hasLW {
		if f, ok := parseFloat(lw); ok {
			b.w = f
		}
	}
	rotation := s.textRotation(t)
	if !st.cell.edge {
		// rotateLabelBounds
		m := alignmentPoint(t.align, t.valign)
		b.y -= m.y * b.h
		b.x -= m.x * b.w
		if t.overflow != "fill" && t.overflow != "width" && (t.overflow != "block" || st.style.get("blockSpacing", "") == "1") {
			sp := t.getSpacing(st.style.get("blockSpacing", "") == "1", m.x)
			b.x += sp.x
			b.y += sp.y
			if hpos == "center" && !hasLW {
				b.w = math.Max(0, b.w-float64(t.spacingLeft+t.spacingRight))
			} else {
				b.w = math.Max(0, b.w)
			}
			if vpos == "middle" {
				b.h = math.Max(0, b.h-float64(t.spacingTop+t.spacingBottom))
			} else {
				b.h = math.Max(0, b.h)
			}
		}
		if rotation != 0 {
			cx, cy := st.cx(), st.cy()
			if b.x != cx || b.y != cy {
				rad := toRadians(rotation)
				p := rotatePoint(point{b.x, b.y}, math.Cos(rad), math.Sin(rad), point{cx, cy})
				b.x, b.y = p.x, p.y
			}
		}
	}

	root := c.rootTextStyle(st.style, t)
	l := &labelBox{anchor: point{b.x, b.y}, rotation: rotation, alpha: clamp(st.style.num("textOpacity", 100)/100, 0, 1)}
	if col, ok := parseColor(st.style.get("labelBackgroundColor", "")); ok && col.a > 0 {
		l.bg = &col
	}
	if col, ok := parseColor(st.style.get("labelBorderColor", "")); ok && col.a > 0 {
		l.border = &col
	}
	l.pad = cssSpacing(st.style.get("labelPadding", ""))
	if t.html {
		c.layoutHTML(l, value, root, t, b)
	} else {
		c.layoutPlain(l, value, root, t, b)
	}
	if l.lo == nil || len(l.lo.lines) == 0 {
		return nil
	}
	return l
}

// rootTextStyle is the style of the label's text block (mxSvgCanvas2D.getTextCss).
func (c *converter) rootTextStyle(s style, t textProps) *tstyle {
	fs := int(s.num("fontStyle", 0))
	col, ok := parseColor(s.get("fontColor", "black"))
	if !ok {
		col = rgba{0, 0, 0, 0}
	}
	fam := parseFontFamilies(s.get("fontFamily", "Arial,Helvetica"))
	st := &tstyle{
		families:  fam,
		size:      s.num("fontSize", 11),
		bold:      fs&1 != 0,
		italic:    fs&2 != 0,
		underline: fs&4 != 0,
		strike:    fs&8 != 0,
		color:     col,
		lhFactor:  1.2,
		align:     map[string]string{"left": "left", "right": "right"}[t.align],
		nowrap:    !t.wrap,
	}
	c.setPrimary(st)
	return st
}

// layoutHTML lays out an HTML label and places it as mxSvgCanvas2D.createCss does.
func (c *converter) layoutHTML(l *labelBox, value string, root *tstyle, t textProps, b rect) {
	if root.align == "" {
		root.align = "center"
	}
	value = htmlLinefeeds(value)
	doc := parseHTML(value)
	tb := &textBuilder{c: c}
	tb.walk(doc, root)
	tb.endPara()
	l.trailing = tb.marks
	paras := tb.paras
	for _, p := range paras {
		markBreaks(p.items)
	}
	x, y, w, h := b.x, b.y, b.w, b.h
	// mxSvgCanvas2D passes the width plus its foreignObjectPadding
	w += foreignObjectPadding
	pt := alignmentPoint(t.align, t.valign)

	// the width lines wrap at, and the box width
	fw := 0.0 // 0: a 1px container (no wrapping)
	fixedItemW := -1.0
	switch {
	case t.clipped:
		fw = math.Round(w)
	case t.overflow == "fill":
		fw = math.Round(w)
		fixedItemW = math.Round(w - 2)
	case t.overflow == "width" || t.overflow == "block":
		fw = math.Round(w - 2)
		fixedItemW = fw
	}
	wrapW := 0.0
	if t.wrap && w > 0 {
		fw = math.Round(w)
		if fixedItemW >= 0 {
			wrapW = fixedItemW
		} else {
			wrapW = fw
		}
	} else if fixedItemW >= 0 {
		// no wrapping, but the item has a width
	}
	lo := layoutParagraphs(paras, wrapW)
	itemW := lo.width
	if lo.wrapped && wrapW > 0 {
		itemW = math.Max(itemW, wrapW)
	} else if wrapW > 0 {
		itemW = math.Min(itemW, wrapW)
	}
	if fixedItemW >= 0 {
		itemW = fixedItemW
	}
	itemH := lo.height
	var maxH float64 = -1
	if t.clipped || t.overflow == "width" && h > 0 || t.overflow == "block" && t.valign == "middle" {
		maxH = math.Round(h)
	}
	if t.overflow == "fill" {
		itemH = math.Round(h)
	} else if maxH >= 0 {
		itemH = math.Min(itemH, maxH)
	}

	// the flex container (dx, dy of createCss)
	ofl := t.clipped || t.overflow == "fill" || t.overflow == "width" || t.overflow == "block"
	dx, dy := pt.x*w, pt.y*h
	switch {
	case t.clipped:
		dy = 0
	case t.overflow == "width" || t.overflow == "block":
		dy = 0
	case t.overflow == "fill":
	default:
		dy = 0
	}
	if t.wrap && w > 0 {
		if ofl && t.overflow != "fill" {
			dy = 0
		}
	} else if !ofl {
		dx = 0
	}
	cw := fw
	if cw == 0 {
		cw = 1
	}
	ch := 1.0
	if t.overflow == "fill" {
		ch = math.Round(h)
	}
	cx0, cy0 := math.Round(x+dx), math.Round(y+dy)
	jf := map[string]float64{"left": 0, "right": 1}[t.align]
	if t.align != "left" && t.align != "right" {
		jf = 0.5
	}
	vf := map[string]float64{"top": 0, "bottom": 1}[t.valign]
	if t.valign != "top" && t.valign != "bottom" {
		vf = 0.5
	}
	l.left = cx0 + (cw-itemW)*jf
	l.top = cy0 + (ch-itemH)*vf
	l.boxW = itemW
	l.lo = lo
	// lines are aligned in the item by their block's text-align
	for _, ln := range lo.lines {
		avail := itemW - ln.x
		switch ln.para.align {
		case "center":
			ln.x += (avail - ln.width) / 2
		case "right":
			ln.x += avail - ln.width
		}
	}
	if ofl {
		l.clip = &rect{l.left, l.top, itemW, itemH}
		if t.overflow == "block" {
			l.clip = nil
		}
		l.bgBox = rect{l.left, l.top, itemW, itemH}
	} else {
		l.bgBox = rect{l.left, l.top, itemW, lo.height}
	}
	if (l.bg != nil || l.border != nil) && l.pad != nil {
		p := l.pad
		l.bgBox = rect{l.bgBox.x - p[3], l.bgBox.y - p[0], l.bgBox.w + p[1] + p[3], l.bgBox.h + p[0] + p[2]}
	}
}

// cssSpacing reads a CSS-style spacing shorthand (labelPadding) as top,
// right, bottom, left (mxUtils.parseCssSpacing); nil when it is empty or
// zero on all sides.
func cssSpacing(v string) []float64 {
	var vals []float64
	for _, t := range strings.Fields(v) {
		if len(vals) == 4 {
			break
		}
		f, ok := parseFloat(t)
		if !ok {
			f = 0
		}
		vals = append(vals, math.Max(0, f))
	}
	if len(vals) == 0 {
		return nil
	}
	top := vals[0]
	right, bottom := top, top
	if len(vals) > 1 {
		right = vals[1]
	}
	if len(vals) > 2 {
		bottom = vals[2]
	}
	left := right
	if len(vals) > 3 {
		left = vals[3]
	}
	if top == 0 && right == 0 && bottom == 0 && left == 0 {
		return nil
	}
	return []float64{top, right, bottom, left}
}

// htmlLinefeeds turns the newlines of an HTML label into line breaks as
// mxText does (trailing ones become empty lines), outside of tags.
func htmlLinefeeds(s string) string {
	n := 0
	for strings.HasSuffix(s, "\n") {
		s = s[:len(s)-1]
		n++
	}
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case r == '\n' && !inTag:
			b.WriteString("<br/>")
			continue
		}
		b.WriteRune(r)
	}
	for range n {
		b.WriteString("<div><br></div>")
	}
	return b.String()
}

// layoutPlain lays out a plain text label as mxSvgCanvas2D.plainText does:
// lines split at newlines, 1.2 line spacing, anchored by text-anchor.
func (c *converter) layoutPlain(l *labelBox, value string, root *tstyle, t textProps, b rect) {
	l.plain = true
	x, y, w, h := b.x, b.y, b.w, b.h
	size := root.size
	lines := strings.Split(value, "\n")
	lh := math.Round(size * 1.2)
	textHeight := size + float64(len(lines)-1)*lh
	cy := y + size - 1
	switch t.valign {
	case "middle":
		if t.overflow == "fill" {
			cy -= h / 2
		} else {
			cy -= textHeight / 2
		}
	case "bottom":
		if t.overflow == "fill" {
			cy -= h
		} else {
			cy -= textHeight + 1
		}
	}
	lo := &tlayout{}
	tb := &textBuilder{c: c}
	var maxW float64
	for i, s := range lines {
		tb.endPara()
		tb.text(s, root)
		p := tb.cur
		if p == nil {
			p = &tpara{st: root}
			tb.paras = append(tb.paras, p)
		}
		ln := &tline{para: p, first: true}
		for _, it := range p.items {
			ln.items = append(ln.items, it)
		}
		for len(ln.items) > 0 && ln.items[0].space {
			ln.items = ln.items[1:]
		}
		for len(ln.items) > 0 && ln.items[len(ln.items)-1].space {
			ln.items = ln.items[:len(ln.items)-1]
		}
		for _, it := range ln.items {
			ln.width += it.w
		}
		ln.asc = 0
		ln.top = cy + float64(i)*lh // baseline, see drawLabel (plain)
		switch t.align {
		case "center":
			ln.x = -ln.width / 2
		case "right":
			ln.x = -ln.width
		}
		maxW = math.Max(maxW, ln.width)
		lo.lines = append(lo.lines, ln)
		lo.paras = append(lo.paras, p)
	}
	lo.width, lo.height = maxW, textHeight
	l.lo = lo
	l.left, l.top = x, 0
	// the background box (mxSvgCanvas2D.addTextBackground)
	bx := x
	switch t.align {
	case "center":
		bx -= maxW / 2
	case "right":
		bx -= maxW
	}
	by := cy - size + 1
	l.bgBox = rect{bx, by, maxW, textHeight}
	if t.clipped && w > 0 && h > 0 {
		cx, cy2 := x, y
		switch t.align {
		case "center":
			cx -= w / 2
		case "right":
			cx -= w
		}
		if t.overflow != "fill" {
			switch t.valign {
			case "middle":
				cy2 -= h / 2
			case "bottom":
				cy2 -= h
			}
		}
		l.clip = &rect{cx - 2, cy2 - 2, w + 4, h + 4}
	}
}

// extent is the rectangle the label paints over (for the page size).
func (l *labelBox) extent() rect {
	r := l.bgBox
	if l.clip != nil {
		r = l.clip.intersect(r)
	}
	if r.w <= 0 && r.h <= 0 {
		return rect{}
	}
	if l.rotation != 0 {
		m := rotateAbout(l.rotation, l.anchor.x, l.anchor.y)
		var out rect
		for _, p := range []point{{r.x, r.y}, {r.x + r.w, r.y}, {r.x, r.y + r.h}, {r.x + r.w, r.y + r.h}} {
			x, y := m.apply(p.x, p.y)
			out = out.addPoint(point{x, y})
		}
		return out
	}
	return r
}

func (r rect) intersect(o rect) rect {
	x0, y0 := math.Max(r.x, o.x), math.Max(r.y, o.y)
	x1, y1 := math.Min(r.x+r.w, o.x+o.w), math.Min(r.y+r.h, o.y+o.h)
	if x1 < x0 || y1 < y0 {
		return rect{}
	}
	return rect{x0, y0, x1 - x0, y1 - y0}
}

// drawLabel draws a laid-out label with its structure marks.
func (c *converter) drawLabel(c2 *c2d, l *labelBox) {
	obj := c2.obj
	cv := c2.cv
	obj.Save()
	cv.drawn = true
	if l.rotation != 0 {
		cv.transform(rotateAbout(l.rotation, l.anchor.x, l.anchor.y))
	}
	if l.bg != nil {
		obj.FillColor(l.bg.withAlpha(l.alpha).bdf())
		obj.FillRect(f32(l.bgBox.x), f32(l.bgBox.y), f32(l.bgBox.w), f32(l.bgBox.h))
	}
	if l.border != nil {
		obj.StrokeColor(l.border.withAlpha(l.alpha).bdf())
		obj.StrokeRect(f32(l.bgBox.x-0.5), f32(l.bgBox.y-0.5), f32(l.bgBox.w+1), f32(l.bgBox.h+1))
	}
	if l.clip != nil {
		obj.ClipRect(f32(l.clip.x), f32(l.clip.y), f32(l.clip.w), f32(l.clip.h))
	}
	e := &textEmitter{c: c, cv: cv, alpha: l.alpha}
	obj.Mark(bdf.MarkBox, "")
	var lastPara *tpara
	paraIndex := 0
	for _, ln := range l.lo.lines {
		if ln.para != lastPara {
			p := ln.para
			startsItem := false
			for _, m := range p.marks {
				obj.Mark(m.kind, m.payload)
				if m.kind == markListItem {
					startsItem = true
				}
			}
			switch {
			case p.heading > 0:
				obj.Mark(bdf.MarkHeading, strconv.Itoa(p.heading))
			case startsItem:
			case paraIndex == 0:
				// the BOX mark stands for the first paragraph
			default:
				obj.Mark(bdf.MarkParagraph, "")
			}
			lastPara = p
			paraIndex++
		} else if ln.hardPrev {
			obj.Mark(bdf.MarkParagraph, "")
		} else if ln.cjkWrap {
			obj.Mark(bdf.MarkWrap, "")
		} else {
			obj.Mark(bdf.MarkLine, "")
		}
		if ln.para.hr {
			obj.FillColor(ln.para.hrColor.withAlpha(l.alpha).bdf())
			e.fill = 0
			obj.FillRect(f32(l.left+ln.x), f32(l.top+ln.top), f32(l.boxW-ln.x), 1)
			continue
		}
		var x, base float64
		if l.plain {
			x, base = l.left+ln.x, ln.top
		} else {
			x, base = l.left+ln.x, l.top+ln.top+ln.asc
		}
		if ln.first && len(ln.para.marker) > 0 {
			mw := 0.0
			for _, it := range ln.para.marker {
				mw += it.w
			}
			e.emitItems(obj, ln.para.marker, l.left+ln.para.indent-mw, base)
			obj.Mark(bdf.MarkLine, "")
		}
		if ln.first && ln.para.bullet != "" {
			e.bullet(obj, ln.para, l.left+ln.para.indent, base)
		}
		e.emitItems(obj, ln.items, x, base)
	}
	for _, m := range l.trailing {
		obj.Mark(m.kind, m.payload)
	}
	obj.Restore()
}

// textEmitter writes text instructions, avoiding repeated state.
type textEmitter struct {
	c     *converter
	cv    *canvas
	alpha float64
	font  bdf.FontRef
	size  float64
	fill  bdf.Color
	set   bool
}

// emitItems draws a line's items from x on the baseline base, as runs of
// one font and style.
func (e *textEmitter) emitItems(obj *bdf.Object, items []titem, x, base float64) {
	for i := 0; i < len(items); {
		j := i + 1
		for j < len(items) && sameRun(items[i], items[j]) {
			j++
		}
		run := items[i:j]
		w := 0.0
		var sb strings.Builder
		for _, it := range run {
			w += it.w
			sb.WriteRune(it.r)
		}
		st := run[0].st
		y := base + st.shift
		fc := run[0].fc
		asc, desc := cssMetrics(fc)
		if st.bg != nil {
			obj.FillColor(st.bg.withAlpha(e.alpha).bdf())
			e.fill = 0
			obj.FillRect(f32(x), f32(y-asc*st.size), f32(w), f32((asc+desc)*st.size))
		}
		ref := e.cv.font(fc.use)
		if !e.set || ref != e.font || st.size != e.size {
			obj.Font(ref, f32(st.size))
			e.font, e.size, e.set = ref, st.size, true
		}
		col := st.color.withAlpha(e.alpha).bdf()
		if col != e.fill {
			obj.FillColor(col)
			e.fill = col
		}
		text := sb.String()
		if strings.TrimSpace(text) != "" || len(run) > 0 {
			obj.FillText(text, f32(x), f32(y), f32(w))
		}
		th := math.Max(1, fc.ulTh*st.size)
		if st.underline {
			obj.FillRect(f32(x), f32(y-fc.ulPos*st.size-th/2), f32(w), f32(th))
		}
		if st.strike {
			obj.FillRect(f32(x), f32(y-0.3*st.size-th/2), f32(w), f32(th))
		}
		if st.link != "" {
			if u := e.c.linkURL(st.link); u != "" {
				obj.Link(f32(x), f32(y-asc*st.size), f32(w), f32((asc+desc)*st.size), u)
			}
		}
		x += w
		i = j
	}
}

// bullet draws the symbol marker of a list item whose content starts at
// x (Blink's ListMarker: a disc, circle or square a third of the ascent
// wide, 7px before the content, its top half the ascent above the baseline).
func (e *textEmitter) bullet(obj *bdf.Object, p *tpara, x, base float64) {
	st := p.st
	if len(p.items) > 0 {
		st = p.items[0].st
	}
	asc, _ := st.fontHeight()
	bw := math.Floor((asc*2/3+1)/2*64) / 64
	left := x - 7 - 1 - bw
	top := base - asc + asc/2
	col := st.color.withAlpha(e.alpha).bdf()
	switch p.bullet {
	case "square":
		obj.FillColor(col)
		obj.FillRect(f32(left), f32(top), f32(bw), f32(bw))
	case "circle":
		path := &bdf.Path{}
		path.Circle(f32(left+bw/2), f32(top+bw/2), f32(bw/2-0.5))
		ref := obj.AddPath(path)
		obj.StrokeColor(col)
		obj.StrokePath(ref)
	default:
		path := &bdf.Path{}
		path.Circle(f32(left+bw/2), f32(top+bw/2), f32(bw/2))
		ref := obj.AddPath(path)
		obj.FillColor(col)
		obj.FillPath(ref, bdf.NonZero)
	}
	e.fill = 0
}

// sameRun reports whether two adjacent items are drawn by one text instruction.
func sameRun(a, b titem) bool {
	return a.fc == b.fc && a.st.size == b.st.size && a.st.color == b.st.color && a.st.shift == b.st.shift &&
		a.st.underline == b.st.underline && a.st.strike == b.st.strike && a.st.link == b.st.link &&
		(a.st.bg == nil) == (b.st.bg == nil) && (a.st.bg == nil || *a.st.bg == *b.st.bg)
}

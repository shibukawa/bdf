package drawio

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
)

// Stencils: shapes defined in draw.io's XML stencil language (mxStencil).
// A stencil is a <shape> element with the size of its coordinate system
// (w, h), how it scales (aspect) and its drawing as a <background> and a
// <foreground> list of canvas instructions. The instructions are run
// against the canvas nearly one to one, with the coordinates scaled to the
// shape's bounds (computeAspect). The stencils of draw.io's libraries are
// looked up by name in stencil_registry.go.

// stencil is a parsed stencil shape (mxStencil).
type stencil struct {
	// w0, h0 are the stencil's own size (the w and h attributes).
	w0, h0 float64
	// aspect is "fixed" when the stencil keeps its proportions
	// ("variable" otherwise).
	aspect string

	// name is the name attribute of the shape (for warnings).
	name string
	// strokewidth is "inherit" (the cell's stroke width is used as it is)
	// or a factor that is scaled with the shape.
	strokewidth string
	// bg and fg are the background and foreground elements (nil when
	// missing).
	bg, fg *stencilNode
	// labelBounds are draw.io's conditional label bounds (Shapes.js).
	labelBounds []stencilLabelBounds
}

// stencilLabelBounds is a <labelBounds> element: the label rectangle in
// stencil coordinates, used when the style key named by if is 1 (or
// always, without if).
type stencilLabelBounds struct {
	cond       string
	hasCond    bool
	x, y, w, h float64
}

// stencilNode is an element of a stencil description.
type stencilNode struct {
	name  string
	attrs []xml.Attr
	kids  []*stencilNode
}

// attr returns the value of an attribute and whether it is set
// (getAttribute, which returns null for a missing attribute).
func (n *stencilNode) attr(key string) (string, bool) {
	for _, a := range n.attrs {
		if a.Name.Local == key {
			return a.Value, true
		}
	}
	return "", false
}

// get returns the value of an attribute, "" when it is not set.
func (n *stencilNode) get(key string) string {
	v, _ := n.attr(key)
	return v
}

// num returns an attribute as JavaScript's Number(getAttribute(key)) does.
func (n *stencilNode) num(key string) float64 { return numberJS(n.get(key)) }

// find returns the first element named name below n, in document order
// (getElementsByTagName(name)[0]).
func (n *stencilNode) find(name string) *stencilNode {
	for _, k := range n.kids {
		if k.name == name {
			return k
		}
		if r := k.find(name); r != nil {
			return r
		}
	}
	return nil
}

// findAll appends the elements named name below n, in document order.
func (n *stencilNode) findAll(name string, out []*stencilNode) []*stencilNode {
	for _, k := range n.kids {
		if k.name == name {
			out = append(out, k)
		}
		out = k.findAll(name, out)
	}
	return out
}

// parseStencilXML reads an XML document into a tree of elements (text,
// comments and processing instructions are dropped) and returns its root.
func parseStencilXML(data []byte) (*stencilNode, error) {
	d := xml.NewDecoder(bytes.NewReader(data))
	var root *stencilNode
	var stack []*stencilNode
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			n := &stencilNode{name: t.Name.Local, attrs: t.Copy().Attr}
			if len(stack) > 0 {
				p := stack[len(stack)-1]
				p.kids = append(p.kids, n)
			} else if root == nil {
				root = n
			}
			stack = append(stack, n)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		}
	}
	if root == nil {
		return nil, errors.New("no root element")
	}
	return root, nil
}

// newStencil parses a <shape> element (mxStencil.parseDescription and the
// labelBounds of Shapes.js; the connection constraints are not needed for
// drawing).
func newStencil(desc *stencilNode) *stencil {
	st := &stencil{name: desc.get("name"), fg: desc.find("foreground"), bg: desc.find("background")}
	// Number(getAttribute('w') || 100)
	st.w0, st.h0 = 100, 100
	if v := desc.get("w"); v != "" {
		st.w0 = numberJS(v)
	}
	if v := desc.get("h"); v != "" {
		st.h0 = numberJS(v)
	}
	st.aspect = "variable"
	if v, ok := desc.attr("aspect"); ok {
		st.aspect = v
	}
	st.strokewidth = "1"
	if v, ok := desc.attr("strokewidth"); ok {
		st.strokewidth = v
	}
	for _, n := range desc.findAll("labelBounds", nil) {
		lb := stencilLabelBounds{x: n.num("x"), y: n.num("y"), w: st.w0, h: st.h0}
		lb.cond, lb.hasCond = n.attr("if")
		if v := n.get("w"); v != "" {
			lb.w = numberJS(v)
		}
		if v := n.get("h"); v != "" {
			lb.h = numberJS(v)
		}
		st.labelBounds = append(st.labelBounds, lb)
	}
	return st
}

// numberJS converts a string to a number as JavaScript's Number() does:
// the empty (or blank) string is 0, anything that is not a number NaN.
func numberJS(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	if len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		if v, err := strconv.ParseUint(s[2:], 16, 64); err == nil {
			return float64(v)
		}
		return math.NaN()
	}
	switch s {
	case "Infinity", "+Infinity":
		return math.Inf(1)
	case "-Infinity":
		return math.Inf(-1)
	}
	if c := s[len(s)-1]; c != '.' && (c < '0' || c > '9') || strings.ContainsRune(s, '_') {
		return math.NaN() // "inf", "nan", "1e", "1_000" and the like
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		return math.NaN()
	}
	return v
}

// maxStencilDepth limits the nesting of include-shape elements (a stencil
// that includes itself).
const maxStencilDepth = 8

// stencilPainter holds what the painting of one stencil shares.
type stencilPainter struct {
	c *c2d
	s *shape
	// base is the depth of the canvas state stack when the stencil
	// started: restore elements do not pop states from before.
	base  int
	depth int
}

// drawShape paints the stencil in the box x, y, w, h (mxStencil.drawShape).
// The shape's paint has rotated and flipped the canvas and configured its
// colors; directions are applied here by computeAspect.
func (st *stencil) drawShape(c *c2d, s *shape, x, y, w, h float64) {
	st.draw(c, s, x, y, w, h, 0)
}

// draw is drawShape for a stencil nested depth include-shape elements deep.
func (st *stencil) draw(c *c2d, s *shape, x, y, w, h float64, depth int) {
	p := &stencilPainter{c: c, s: s, base: len(c.stack), depth: depth}
	direction := s.style.get("direction", "")
	aspect := st.computeAspect(x, y, w, h, direction)
	minScale := math.Min(aspect.w, aspect.h)
	var sw float64
	if st.strokewidth == "inherit" {
		sw = s.style.num("strokeWidth", 1)
	} else {
		sw = numberJS(st.strokewidth) * minScale
	}
	c.setStrokeWidth(sw)

	// the transparent rectangle that catches events draws nothing
	if s.style.get("pointerEvents", "0") == "1" {
		c.setStrokeColor("")
		c.rect(x, y, w, h)
		c.stroke()
		c.setStrokeColor(s.stroke)
	}

	st.drawChildren(p, x, y, w, h, st.bg, aspect, false, true)
	st.drawChildren(p, x, y, w, h, st.fg, aspect, true,
		!s.outline || s.style.num("backgroundOutline", 0) == 0)

	// restores the stack for unequal counts of save and restore
	for len(c.stack) > p.base {
		c.restore()
	}
}

// drawChildren draws the elements of node (mxStencil.drawChildren).
func (st *stencil) drawChildren(p *stencilPainter, x, y, w, h float64, node *stencilNode, aspect rect, disableShadow, paint bool) {
	if node == nil || w <= 0 || h <= 0 {
		return
	}
	for _, k := range node.kids {
		st.drawNode(p, k, aspect, disableShadow, paint)
	}
}

// computeAspect returns the offset (x, y) and the scales (w, h) that map
// stencil coordinates into the box x, y, w, h: stretched, or centered with
// the stencil's proportions for a fixed aspect, and with the axes swapped
// for north and south directions, where the canvas is turned by 90°
// (mxStencil.computeAspect).
func (st *stencil) computeAspect(x, y, w, h float64, direction string) rect {
	x0, y0 := x, y
	sx := w / st.w0
	sy := h / st.h0
	inverse := direction == "north" || direction == "south"
	if inverse {
		sy = w / st.h0
		sx = h / st.w0
		delta := (w - h) / 2
		x0 += delta
		y0 -= delta
	}
	if st.aspect == "fixed" {
		sy = math.Min(sx, sy)
		sx = sy
		// centers the shape inside the available space
		if inverse {
			x0 += (h - st.w0*sx) / 2
			y0 += (w - st.h0*sy) / 2
		} else {
			x0 += (w - st.w0*sx) / 2
			y0 += (h - st.h0*sy) / 2
		}
	}
	return rect{x0, y0, sx, sy}
}

// drawNode runs one element of a stencil (mxStencil.drawNode).
func (st *stencil) drawNode(p *stencilPainter, node *stencilNode, aspect rect, disableShadow, paint bool) {
	c, s := p.c, p.s
	name := node.name
	x0, y0 := aspect.x, aspect.y
	sx, sy := aspect.w, aspect.h
	minScale := math.Min(sx, sy)
	px := func(key string) float64 { return x0 + node.num(key)*sx }
	py := func(key string) float64 { return y0 + node.num(key)*sy }

	switch {
	case name == "save":
		c.save()
		return
	case name == "restore":
		if len(c.stack) > p.base {
			c.restore()
		}
		return
	case !paint:
		return
	}
	switch name {
	case "path":
		c.begin()
		regular := true
		if node.get("rounded") == "1" {
			regular = false
			arcSize := node.num("arcSize")
			pointCount := 0
			var segs [][]point
			for _, k := range node.kids {
				if k.name != "move" && k.name != "line" {
					// only moves and lines have rounded corners
					regular = true
					break
				}
				if k.name == "move" || len(segs) == 0 {
					segs = append(segs, nil)
				}
				segs[len(segs)-1] = append(segs[len(segs)-1],
					point{x0 + k.num("x")*sx, y0 + k.num("y")*sy})
				pointCount++
			}
			if !regular && pointCount > 0 {
				for _, seg := range segs {
					closed := false
					ps, pe := seg[0], seg[len(seg)-1]
					if ps.x == pe.x && ps.y == pe.y {
						seg = seg[:len(seg)-1]
						closed = true
					}
					s.addPoints(c, seg, true, arcSize, closed, nil, true)
				}
			} else {
				regular = true
			}
		}
		if regular {
			for _, k := range node.kids {
				st.drawNode(p, k, aspect, disableShadow, paint)
			}
		}
	case "close":
		c.close()
	case "move":
		c.moveTo(px("x"), py("y"))
	case "line":
		c.lineTo(px("x"), py("y"))
	case "quad":
		c.quadTo(px("x1"), py("y1"), px("x2"), py("y2"))
	case "curve":
		c.curveTo(px("x1"), py("y1"), px("x2"), py("y2"), px("x3"), py("y3"))
	case "arc":
		c.arcTo(node.num("rx")*sx, node.num("ry")*sy, node.num("x-axis-rotation"),
			node.num("large-arc-flag") != 0, node.num("sweep-flag") != 0, px("x"), py("y"))
	case "rect":
		c.rect(px("x"), py("y"), node.num("w")*sx, node.num("h")*sy)
	case "roundrect":
		arcsize := node.num("arcsize")
		if arcsize == 0 {
			arcsize = rectangleRoundingFactor * 100
		}
		w := node.num("w") * sx
		h := node.num("h") * sy
		factor := arcsize / 100
		r := math.Min(w*factor, h*factor)
		c.roundrect(px("x"), py("y"), w, h, r, r)
	case "ellipse":
		c.ellipse(px("x"), py("y"), node.num("w")*sx, node.num("h")*sy)
	case "image":
		if !s.outline && s.conv != nil {
			// scripts in the element's text (allowEval) are not run
			if src := node.get("src"); src != "" {
				img := s.conv.imageFor(src)
				c.image(px("x"), py("y"), node.num("w")*sx, node.num("h")*sy, img,
					false, node.get("flipH") == "1", node.get("flipV") == "1")
			}
		}
	case "text":
		// Text elements are not drawn yet (labels are not implemented).
		if !s.outline && s.conv != nil {
			s.conv.warnOnce("stencil-text", "text in stencil shapes (e.g. %q) is not drawn", st.name)
		}
	case "include-shape":
		if sub := lookupStencil(node.get("name")); sub != nil && p.depth < maxStencilDepth {
			sub.draw(c, s, px("x"), py("y"), node.num("w")*sx, node.num("h")*sy, p.depth+1)
		}
	case "fillstroke":
		c.fillAndStroke()
	case "fill":
		c.fill()
	case "stroke":
		c.stroke()
	case "strokewidth":
		f := minScale
		if node.get("fixed") == "1" {
			f = 1
		}
		c.setStrokeWidth(node.num("width") * f)
	case "dashed":
		c.setDashed(node.get("dashed") == "1", false)
	case "dashpattern":
		if v, ok := node.attr("pattern"); ok {
			var pat []string
			for _, t := range strings.Split(v, " ") {
				if t != "" {
					pat = append(pat, fmtNum(numberJS(t)*minScale))
				}
			}
			c.setDashPattern(strings.Join(pat, " "))
		}
	case "strokecolor":
		c.setStrokeColor(st.parseColor(s, node))
	case "linecap":
		c.setLineCap(node.get("cap"))
	case "linejoin":
		c.setLineJoin(node.get("join"))
	case "miterlimit":
		c.setMiterLimit(node.num("limit"))
	case "fillcolor":
		c.setFillColor(st.parseColor(s, node))
	case "alpha", "fillalpha", "strokealpha":
		// mxStencil sets the alpha of everything for all three
		c.setAlpha(node.num("alpha"))
	case "fontcolor", "fontstyle", "fontfamily", "fontsize":
		// for text elements, which are not drawn
	}
	if disableShadow && (name == "fillstroke" || name == "fill" || name == "stroke") {
		c.setShadow(false)
	}
}

// Default colors of stencils: draw.io's shapeBackgroundColor and
// shapeForegroundColor of the light theme.
const (
	stencilBackgroundColor = "#ffffff"
	stencilForegroundColor = "#000000"
)

// parseColor returns the color a strokecolor or fillcolor element sets
// (mxStencil.parseColor with draw.io's override for "default"): a color,
// "stroke" or "fill" for the shape's colors, or the name of a style key
// with an optional default attribute. "" means none.
func (st *stencil) parseColor(s *shape, node *stencilNode) string {
	value := node.get("color")
	switch {
	case value == "default":
		return defaultStencilColor(node)
	case value == "stroke":
		return s.stroke
	case value == "fill":
		return s.fill
	case value == "font":
		return s.style.get("fontColor", "black")
	case !validColor(value):
		value = stencilColorValue(s, node, value)
		if value == "default" {
			value = defaultStencilColor(node)
		}
	}
	return colorOrNone(value)
}

// defaultStencilColor is the color "default" stands for (mxStencil.getDefaultColorValue).
func defaultStencilColor(node *stencilNode) string {
	if node.name == "fillcolor" {
		return stencilBackgroundColor
	}
	return stencilForegroundColor
}

// stencilColorValue returns the color in the style key a color element
// names, or its default attribute (mxStencil.getColorValue). A key set to
// "none" in the cell's style wins over the default.
func stencilColorValue(s *shape, node *stencilNode, key string) string {
	def, hasDef := node.attr("default")
	result, ok := s.style[key]
	if !ok {
		temp, set := styleWithNone(s, key)
		if !hasDef || set {
			return temp
		}
		return def
	}
	if result == "default" && hasDef && def != "none" {
		return def
	}
	return result
}

// styleWithNone looks a key up in the style including the "none" values
// that the resolved style drops (mxStencil.getStyleWithNone).
func styleWithNone(s *shape, key string) (string, bool) {
	if v, ok := s.style[key]; ok {
		return v, true
	}
	if s.st == nil || s.st.cell == nil {
		return "", false
	}
	none := false
	for _, tok := range strings.Split(s.st.cell.styleStr, ";") {
		if k, v, ok := strings.Cut(tok, "="); ok {
			if k == key {
				none = v == "none"
			}
		} else if named, ok := namedStyles[tok]; ok {
			if v, ok := named[key]; ok {
				none = v == "none"
			}
		}
	}
	if none {
		return "none", true
	}
	return "", false
}

// validColor reports whether v is a CSS color (mxUtils.isValidColor).
func validColor(v string) bool {
	if strings.EqualFold(strings.TrimSpace(v), "transparent") {
		return true
	}
	_, ok := parseColor(v)
	return ok
}

// labelMargins returns the insets of the label rectangle r (left, top,
// right and bottom as x, y, w and h) given by the first <labelBounds>
// whose condition holds, or nil (draw.io's mxShape.getLabelMargins for
// stencils in Shapes.js). Directions and flips are left to the caller
// (mxShape.getLabelBounds).
func (st *stencil) labelMargins(s *shape, r rect) *rect {
	for _, lb := range st.labelBounds {
		if lb.hasCond && s.style.get(lb.cond, "0") != "1" {
			continue
		}
		a := st.computeAspect(r.x, r.y, r.w, r.h, "")
		x0 := a.x - r.x + lb.x*a.w
		y0 := a.y - r.y + lb.y*a.h
		return &rect{x0, y0, r.w - x0 - lb.w*a.w, r.h - y0 - lb.h*a.h}
	}
	return nil
}

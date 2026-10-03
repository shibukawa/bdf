package imagebdf

import (
	"encoding/base64"
	"image"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// svgImage is an SVG image part and the rasters drawn of it. The viewer
// draws SVG images with the browser; here they are drawn by a renderer of
// the SVG that previews meet: shapes and paths, fills and strokes, linear
// and radial gradients, transforms, use and symbol, nested svg, clip
// paths, group opacity, text in the fonts of the system, and embedded
// images (data: URLs). Masks, filters, patterns and markers are not drawn.
type svgImage struct {
	doc     *svgDoc
	w, h    float64 // natural size in CSS px (as the viewer measures it)
	rasters map[int]*picture
}

// svgRasterPixels bounds the size of a raster of an SVG image: the scale
// it is drawn at is no more than gives as many pixels, rounded up to a
// power of √2 (so a raster has up to twice as many).
const svgRasterPixels = 4096 * 4096

// newSVGImage parses an SVG image part.
func newSVGImage(data []byte) (*svgImage, bool) {
	doc, ok := parseSVG(data)
	if !ok {
		return nil, false
	}
	si := &svgImage{doc: doc, rasters: map[int]*picture{}}
	si.w, si.h = svgNaturalSize(doc.root)
	if math.IsInf(si.w, 0) || math.IsInf(si.h, 0) {
		return nil, false
	}
	return si, true
}

// svgNaturalSize is the size of an SVG image: the width and height of its
// root, a missing one following from the other and the view box's
// proportions, else the view box's size, else 300 × 150.
func svgNaturalSize(root *svgNode) (float64, float64) {
	abs := func(s string) (float64, bool) {
		if strings.HasSuffix(strings.TrimSpace(s), "%") {
			return 0, false
		}
		v, ok := svgLength(s, 0, 16)
		return v, ok && v > 0
	}
	w, okw := abs(root.attr["width"])
	h, okh := abs(root.attr["height"])
	vb := svgNumbers(root.attr["viewBox"])
	vw, vh := 0.0, 0.0
	if len(vb) == 4 && vb[2] > 0 && vb[3] > 0 {
		vw, vh = vb[2], vb[3]
	}
	if okw && okh {
		return w, h
	}
	switch {
	case okw && vw > 0:
		return w, w * vh / vw
	case okh && vw > 0:
		return h * vw / vh, h
	case vw > 0:
		return vw, vh
	}
	if !okw {
		w = 300
	}
	if !okh {
		h = 150
	}
	return w, h
}

// raster returns a raster of the image drawn with k device pixels per image
// pixel (rounded up to a power of √2, and bounded).
func (si *svgImage) raster(r *Renderer, k float64) (*picture, float64) {
	return si.rasterIn(r, k, 0)
}

// rasterIn is raster for an image that is inside as many SVG images as
// inside says. The picture is nil when the renderer has no room for it.
func (si *svgImage) rasterIn(r *Renderer, k float64, inside int) (*picture, float64) {
	step, k, w, h := si.rasterSize(k)
	if w == 0 {
		return nil, k
	}
	if p, ok := si.rasters[step]; ok {
		p.used = r.drawing
		return p, k
	}
	if !r.makeRoom(w * h) {
		r.warnf(tooManyPixels, r.room)
		return nil, k
	}
	s := newSurface(w, h)
	sr := &svgRenderer{r: r, doc: si.doc, budget: svgBudget, inside: inside}
	sr.render(s, matrix{k, 0, 0, k, 0, 0}, si.w, si.h)
	p := &picture{w: w, h: h, levels: []*image.RGBA{s.toRGBA()}}
	si.rasters[step] = p
	r.keep(p, func() { delete(si.rasters, step) })
	return p, float64(w) / si.w
}

// rasterSize returns the scale a raster for k device pixels per image pixel
// is drawn at, which is 2 to the power of step/2, and its size; no size for
// an image that has no raster.
func (si *svgImage) rasterSize(k float64) (step int, scale float64, w, h int) {
	if !(k > 0) {
		k = 1
	}
	if max := math.Sqrt(svgRasterPixels / (si.w * si.h)); k > max {
		k = max
	}
	// between 2^-32 and 2^32: no raster is scaled more, and the scale of
	// an image of no size is no number
	step = int(math.Min(math.Max(math.Ceil(2*math.Log2(k)), -64), 64))
	size := func() (w, h float64) {
		k = math.Pow(2, float64(step)/2)
		return math.Max(1, math.Ceil(si.w*k)), math.Max(1, math.Ceil(si.h*k))
	}
	fw, fh := size()
	for fw*fh > 2*svgRasterPixels && math.Min(si.w, si.h)*k < 1 && step > -64 {
		// less than a pixel wide or high, which is one all the same: the
		// other side gives the scale
		step--
		fw, fh = size()
	}
	if !(fw*fh <= 4*svgRasterPixels+2) {
		return step, k, 0, 0
	}
	return step, k, int(fw), int(fh)
}

// svgRenderer draws an SVG document.
type svgRenderer struct {
	r     *Renderer
	doc   *svgDoc
	depth int
	// budget is the number of elements left to draw: use elements can
	// repeat content exponentially, as they can in a browser
	budget int
	// layers counts the canvases of translucent elements in use, live
	// their pixels, and layerPixels the pixels of all that were
	layers, live, layerPixels int
	// inside is the number of SVG images this one is in
	inside int
}

// Limits of drawing an SVG image.
const (
	// svgBudget is the number of elements drawn of an SVG image at most.
	svgBudget = 200000
	// svgLayers is the most translucent elements drawn inside each other:
	// each is drawn on a canvas of its own, as large as the raster, and
	// svgLivePixels the most pixels of the canvases in use.
	svgLayers     = 8
	svgLivePixels = 4 * svgRasterPixels
	// svgLayerPixels is the most pixels of the canvases of the translucent
	// elements of a raster, which are blended in one after the other.
	svgLayerPixels = 1 << 30
	// svgInside is how deep SVG images may be in SVG images.
	svgInside = 4
)

// svgCtx is where drawing goes: a surface, the user space, the clip and
// the viewport that percentages refer to.
type svgCtx struct {
	s      *surface
	m      matrix
	clip   *mask
	vw, vh float64
}

// svgStyle is the computed style of an element.
type svgStyle struct {
	fill, stroke  string
	fillOpacity   float64
	strokeOpacity float64
	fillRule      byte
	clipRule      byte
	strokeWidth   string
	cap, join     byte
	miter         float64
	dash          string
	dashOffset    string
	color         svgColor
	fontFamily    string
	fontSize      float64
	fontWeight    int
	italic        bool
	anchor        string
	visible       bool
	letterSpacing float64
}

func defaultStyle() svgStyle {
	return svgStyle{
		fill: "black", stroke: "none", fillOpacity: 1, strokeOpacity: 1,
		strokeWidth: "1", miter: 4, color: svgColor{0, 0, 0, 1},
		fontFamily: "serif", fontSize: 16, fontWeight: 400, anchor: "start", visible: true,
	}
}

// declared returns the value an element declares for a property: in its
// style attribute, then the style sheets, then as a presentation attribute.
func (n *svgNode) declared(name string) (string, bool) {
	if st := n.attr["style"]; st != "" {
		for _, d := range parseDecls(st) {
			if d[0] == name {
				return d[1], true
			}
		}
	}
	if v, ok := n.sheet[name]; ok {
		return v, true
	}
	v, ok := n.attr[name]
	return strings.TrimSpace(v), ok
}

// compute resolves the inherited properties of an element.
func (sr *svgRenderer) compute(n *svgNode, parent svgStyle) svgStyle {
	st := parent
	get := func(name string) (string, bool) {
		v, ok := n.declared(name)
		if !ok || v == "inherit" || v == "" {
			return "", false
		}
		return v, true
	}
	if v, ok := get("color"); ok {
		if c, ok := parseColor(v, parent.color); ok {
			st.color = c
		}
	}
	if v, ok := get("fill"); ok {
		st.fill = v
	}
	if v, ok := get("stroke"); ok {
		st.stroke = v
	}
	num := func(name string, dst *float64) {
		if v, ok := get(name); ok {
			if f, rest, ok := svgNumber(v); ok {
				if strings.TrimSpace(rest) == "%" {
					f /= 100
				}
				*dst = math.Min(1, math.Max(0, f))
			}
		}
	}
	num("fill-opacity", &st.fillOpacity)
	num("stroke-opacity", &st.strokeOpacity)
	rule := func(name string, dst *byte) {
		if v, ok := get(name); ok {
			*dst = bdf.NonZero
			if v == "evenodd" {
				*dst = bdf.EvenOdd
			}
		}
	}
	rule("fill-rule", &st.fillRule)
	rule("clip-rule", &st.clipRule)
	if v, ok := get("stroke-width"); ok {
		st.strokeWidth = v
	}
	if v, ok := get("stroke-linecap"); ok {
		st.cap = map[string]byte{"round": bdf.CapRound, "square": bdf.CapSquare}[v]
	}
	if v, ok := get("stroke-linejoin"); ok {
		st.join = map[string]byte{"round": bdf.JoinRound, "bevel": bdf.JoinBevel}[v]
	}
	if v, ok := get("stroke-miterlimit"); ok {
		if f, _, ok := svgNumber(v); ok && f >= 1 {
			st.miter = f
		}
	}
	if v, ok := get("stroke-dasharray"); ok {
		st.dash = v
	}
	if v, ok := get("stroke-dashoffset"); ok {
		st.dashOffset = v
	}
	if v, ok := get("font-family"); ok {
		st.fontFamily = v
	}
	if v, ok := get("font-size"); ok {
		if f, ok := svgLength(v, parent.fontSize, parent.fontSize); ok && f > 0 {
			st.fontSize = f
		} else if kw, ok := map[string]float64{"xx-small": 9, "x-small": 10, "small": 13, "medium": 16, "large": 18, "x-large": 24, "xx-large": 32}[v]; ok {
			st.fontSize = kw
		}
	}
	if v, ok := get("font-weight"); ok {
		switch v {
		case "normal":
			st.fontWeight = 400
		case "bold":
			st.fontWeight = 700
		case "bolder":
			st.fontWeight = min(900, st.fontWeight+300)
		case "lighter":
			st.fontWeight = max(100, st.fontWeight-300)
		default:
			if w, err := strconv.Atoi(v); err == nil {
				st.fontWeight = w
			}
		}
	}
	if v, ok := get("font-style"); ok {
		st.italic = v == "italic" || strings.HasPrefix(v, "oblique")
	}
	if v, ok := get("text-anchor"); ok {
		st.anchor = v
	}
	if v, ok := get("visibility"); ok {
		st.visible = v == "visible"
	}
	if v, ok := get("letter-spacing"); ok {
		if f, ok := svgLength(v, 0, st.fontSize); ok {
			st.letterSpacing = f
		} else {
			st.letterSpacing = 0
		}
	}
	return st
}

// render draws the document on a surface whose user space is m, for a
// viewport of w × h CSS px.
func (sr *svgRenderer) render(s *surface, m matrix, w, h float64) {
	root := sr.doc.root
	ctx := &svgCtx{s: s, m: m, clip: &mask{r: s.bounds()}, vw: w, vh: h}
	st := sr.compute(root, defaultStyle())
	if vb, ok := viewBox(root.attr["viewBox"]); ok {
		ctx.m = ctx.m.mul(aspectTransform(vb, root.attr["preserveAspectRatio"], 0, 0, w, h))
		ctx.vw, ctx.vh = vb[2], vb[3]
	}
	sr.children(root, ctx, st)
}

func viewBox(s string) ([4]float64, bool) {
	v := svgNumbers(s)
	if len(v) != 4 || !(v[2] > 0) || !(v[3] > 0) {
		return [4]float64{}, false
	}
	return [4]float64{v[0], v[1], v[2], v[3]}, true
}

// aspectTransform maps a view box into the viewport (x, y, w, h) as
// preserveAspectRatio says.
func aspectTransform(vb [4]float64, par string, x, y, w, h float64) matrix {
	sx, sy := w/vb[2], h/vb[3]
	f := strings.Fields(par)
	if len(f) > 0 && f[0] == "defer" {
		f = f[1:]
	}
	align, slice := "xMidYMid", false
	if len(f) > 0 {
		align = f[0]
	}
	if len(f) > 1 && f[1] == "slice" {
		slice = true
	}
	if align != "none" {
		s := math.Min(sx, sy)
		if slice {
			s = math.Max(sx, sy)
		}
		sx, sy = s, s
	}
	tx, ty := x-vb[0]*sx, y-vb[1]*sy
	if align != "none" {
		switch {
		case strings.Contains(align, "xMid"):
			tx += (w - vb[2]*sx) / 2
		case strings.Contains(align, "xMax"):
			tx += w - vb[2]*sx
		}
		switch {
		case strings.Contains(align, "YMid"):
			ty += (h - vb[3]*sy) / 2
		case strings.Contains(align, "YMax"):
			ty += h - vb[3]*sy
		}
	}
	return matrix{sx, 0, 0, sy, tx, ty}
}

func (sr *svgRenderer) children(n *svgNode, ctx *svgCtx, st svgStyle) {
	for _, c := range n.children {
		if n.name == "switch" {
			// the first child that browsers draw: those without requirements
			if c.name != "#text" && c.attr["requiredExtensions"] == "" && c.attr["systemLanguage"] == "" {
				sr.node(c, ctx, st)
				return
			}
			continue
		}
		sr.node(c, ctx, st)
	}
}

// length resolves a length attribute: dir is 'x' or 'y' for percentages
// of the viewport's width or height, anything else for its diagonal.
func (ctx *svgCtx) length(s string, dir byte, em float64) float64 {
	ref := math.Hypot(ctx.vw, ctx.vh) / math.Sqrt2
	switch dir {
	case 'x':
		ref = ctx.vw
	case 'y':
		ref = ctx.vh
	}
	v, _ := svgLength(s, ref, em)
	return v
}

// node draws an element.
func (sr *svgRenderer) node(n *svgNode, ctx *svgCtx, parent svgStyle) {
	if n.name == "#text" {
		return
	}
	if d, _ := n.declared("display"); d == "none" {
		return
	}
	switch n.name {
	case "defs", "symbol", "clipPath", "mask", "linearGradient", "radialGradient", "pattern", "marker",
		"title", "desc", "metadata", "style", "script", "filter", "foreignObject":
		return
	}
	if sr.depth > 32 || sr.budget <= 0 {
		if sr.budget == 0 {
			sr.budget--
			sr.r.warnf("an SVG image has too many elements: the rest are not drawn")
		}
		return
	}
	sr.budget--
	sr.depth++
	defer func() { sr.depth-- }()
	st := sr.compute(n, parent)
	c := *ctx
	if t := n.attr["transform"]; t != "" {
		c.m = c.m.mul(parseTransform(t))
	}
	if cp, _ := n.declared("clip-path"); cp != "" && cp != "none" {
		if target := sr.ref(cp); target != nil && target.name == "clipPath" {
			c.clip = intersect(c.clip, sr.clipMask(target, n, &c, st))
		}
	}
	if c.clip.empty() {
		return
	}
	opacity := 1.0
	if v, ok := n.declared("opacity"); ok {
		if f, rest, ok := svgNumber(v); ok {
			if strings.TrimSpace(rest) == "%" {
				f /= 100
			}
			opacity = math.Min(1, math.Max(0, f))
		}
	}
	if opacity <= 0 {
		return
	}
	if m, _ := n.declared("mask"); m != "" && m != "none" {
		sr.r.warnf("SVG masks are not applied")
	}
	if opacity < 1 {
		// draw the element alone, then blend it in with its opacity
		pixels := c.s.w * c.s.h
		if sr.layers >= svgLayers || pixels > svgLivePixels-sr.live || pixels > svgLayerPixels-sr.layerPixels {
			sr.r.warnf("an SVG image has too many translucent elements: the rest of them are not drawn")
			return
		}
		sr.layers++
		sr.live += pixels
		sr.layerPixels += pixels
		layer := newSurface(c.s.w, c.s.h)
		lc := c
		lc.s = layer
		sr.element(n, &lc, st)
		sr.layers--
		sr.live -= pixels
		c.s.fill(&mask{r: c.clip.r}, surfaceShader{s: layer}, float32(opacity), bdf.BlendSourceOver, &mask{r: c.s.bounds()})
		return
	}
	sr.element(n, &c, st)
}

// element draws an element in its own user space.
func (sr *svgRenderer) element(n *svgNode, c *svgCtx, st svgStyle) {
	switch n.name {
	case "svg":
		sr.nested(n, false, c, st, c.length(n.attr["x"], 'x', st.fontSize), c.length(n.attr["y"], 'y', st.fontSize), n.attr["width"], n.attr["height"])
	case "g", "a", "switch":
		sr.children(n, c, st)
	case "use":
		target := sr.ref(n.attr["href"])
		if target == nil || sr.isAncestor(target, n) {
			return
		}
		x, y := c.length(n.attr["x"], 'x', st.fontSize), c.length(n.attr["y"], 'y', st.fontSize)
		switch target.name {
		case "symbol", "svg":
			sr.nested(target, true, c, st, x, y, n.attr["width"], n.attr["height"])
		default:
			uc := *c
			uc.m = uc.m.translate(x, y)
			sr.node(target, &uc, st)
		}
	case "text":
		sr.text(n, c, st)
	case "image":
		sr.image(n, c, st)
	default:
		if p := sr.shape(n, c, st); p != nil {
			sr.paint(p, c, st)
		}
	}
}

func (sr *svgRenderer) isAncestor(a, n *svgNode) bool {
	for p := n; p != nil; p = p.parent {
		if p == a {
			return true
		}
	}
	return false
}

// nested draws the content of an svg or symbol element in a new viewport
// at (x, y), ws × hs (the element's own width and height when empty, else
// 100%). viaUse says that a use element shows it, whose style it takes.
func (sr *svgRenderer) nested(content *svgNode, viaUse bool, c *svgCtx, st svgStyle, x, y float64, ws, hs string) {
	if ws == "" {
		ws = content.attr["width"]
	}
	if hs == "" {
		hs = content.attr["height"]
	}
	if ws == "" {
		ws = "100%"
	}
	if hs == "" {
		hs = "100%"
	}
	w, h := c.length(ws, 'x', st.fontSize), c.length(hs, 'y', st.fontSize)
	if !(w > 0 && h > 0) {
		return
	}
	nc := *c
	if viaUse {
		st = sr.compute(content, st)
	}
	if ov, _ := content.declared("overflow"); ov != "visible" && ov != "auto" {
		p := &path{}
		p.rect(x, y, w, h)
		nc.clip = intersect(nc.clip, sr.coverage(p, nc.m, bdf.NonZero, nc.s.bounds()))
	}
	if vb, ok := viewBox(content.attr["viewBox"]); ok {
		nc.m = nc.m.mul(aspectTransform(vb, content.attr["preserveAspectRatio"], x, y, w, h))
		nc.vw, nc.vh = vb[2], vb[3]
	} else {
		nc.m = nc.m.translate(x, y)
		nc.vw, nc.vh = w, h
	}
	sr.children(content, &nc, st)
}

// ref finds the element a reference names: "#id" or "url(#id)".
func (sr *svgRenderer) ref(s string) *svgNode {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "url(") {
		end := strings.IndexByte(s, ')')
		if end < 0 {
			return nil
		}
		s = strings.Trim(strings.TrimSpace(s[4:end]), `"'`)
	}
	if !strings.HasPrefix(s, "#") {
		return nil
	}
	return sr.doc.ids[s[1:]]
}

// shape returns the path of a basic shape or path element.
func (sr *svgRenderer) shape(n *svgNode, c *svgCtx, st svgStyle) *path {
	a := n.attr
	lx := func(k string) float64 { return c.length(a[k], 'x', st.fontSize) }
	ly := func(k string) float64 { return c.length(a[k], 'y', st.fontSize) }
	p := &path{}
	switch n.name {
	case "path":
		if n.path == nil {
			n.path = parsePathData(a["d"])
		}
		return n.path
	case "rect":
		x, y, w, h := lx("x"), ly("y"), lx("width"), ly("height")
		if !(w > 0 && h > 0) {
			return nil
		}
		rx, okx := svgLength(a["rx"], c.vw, st.fontSize)
		ry, oky := svgLength(a["ry"], c.vh, st.fontSize)
		switch {
		case okx && !oky:
			ry = rx
		case oky && !okx:
			rx = ry
		}
		rx, ry = math.Min(math.Max(rx, 0), w/2), math.Min(math.Max(ry, 0), h/2)
		if rx == 0 || ry == 0 {
			p.rect(x, y, w, h)
			return p
		}
		p.moveTo(x+rx, y)
		p.lineTo(x+w-rx, y)
		p.ellipse(x+w-rx, y+ry, rx, ry, 0, -math.Pi/2, 0, false)
		p.lineTo(x+w, y+h-ry)
		p.ellipse(x+w-rx, y+h-ry, rx, ry, 0, 0, math.Pi/2, false)
		p.lineTo(x+rx, y+h)
		p.ellipse(x+rx, y+h-ry, rx, ry, 0, math.Pi/2, math.Pi, false)
		p.lineTo(x, y+ry)
		p.ellipse(x+rx, y+ry, rx, ry, 0, math.Pi, 3*math.Pi/2, false)
		p.close()
	case "circle":
		r := c.length(a["r"], 'd', st.fontSize)
		if !(r > 0) {
			return nil
		}
		p.moveTo(lx("cx")+r, ly("cy"))
		p.ellipse(lx("cx"), ly("cy"), r, r, 0, 0, 2*math.Pi, false)
		p.close()
	case "ellipse":
		rx, ry := lx("rx"), ly("ry")
		if a["rx"] == "" || a["rx"] == "auto" {
			rx = ry
		}
		if a["ry"] == "" || a["ry"] == "auto" {
			ry = rx
		}
		if !(rx > 0 && ry > 0) {
			return nil
		}
		p.moveTo(lx("cx")+rx, ly("cy"))
		p.ellipse(lx("cx"), ly("cy"), rx, ry, 0, 0, 2*math.Pi, false)
		p.close()
	case "line":
		p.moveTo(lx("x1"), ly("y1"))
		p.lineTo(lx("x2"), ly("y2"))
	case "polyline", "polygon":
		v := svgNumbers(a["points"])
		if len(v) < 4 {
			return nil
		}
		p.moveTo(v[0], v[1])
		for i := 2; i+1 < len(v); i += 2 {
			p.lineTo(v[i], v[i+1])
		}
		if n.name == "polygon" {
			p.close()
		}
	default:
		return nil
	}
	return p
}

// bbox is the bounding box of a path's points in user space (control
// points included).
func (p *path) bbox() (x0, y0, x1, y1 float64) {
	x0, y0, x1, y1 = math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	add := func(x, y float64) {
		x0, y0, x1, y1 = math.Min(x0, x), math.Min(y0, y), math.Max(x1, x), math.Max(y1, y)
	}
	for _, s := range p.subs {
		if len(s.segs) == 0 {
			continue
		}
		add(s.start.x, s.start.y)
		for _, g := range s.segs {
			n := map[int]int{segLine: 1, segQuad: 2, segCubic: 3}[g.kind]
			for i := 0; i < n; i++ {
				add(g.p[2*i], g.p[2*i+1])
			}
		}
	}
	return
}

// paint fills and strokes a path with an element's style.
func (sr *svgRenderer) paint(p *path, c *svgCtx, st svgStyle) {
	if !st.visible {
		return
	}
	bx0, by0, bx1, by1 := p.bbox()
	bbox := [4]float64{bx0, by0, bx1 - bx0, by1 - by0}
	if sh := sr.paintServer(st.fill, st, bbox, c); sh != nil {
		cov := sr.coverage(p, c.m, st.fillRule, c.clip.r.Intersect(c.s.bounds()))
		c.s.fill(cov, sh, float32(st.fillOpacity), bdf.BlendSourceOver, c.clip)
	}
	if sh := sr.paintServer(st.stroke, st, bbox, c); sh != nil {
		ls := lineStyle{width: c.length(st.strokeWidth, 'd', st.fontSize), cap: st.cap, join: st.join, miter: st.miter}
		if !(ls.width > 0) {
			return
		}
		if st.dash != "" && st.dash != "none" {
			var segs []float32
			for _, f := range strings.FieldsFunc(st.dash, func(r rune) bool { return r == ',' || r == ' ' }) {
				segs = append(segs, float32(c.length(f, 'd', st.fontSize)))
			}
			ls.setDash(segs, float32(c.length(st.dashOffset, 'd', st.fontSize)))
		}
		cov, alpha, cut := sr.r.scratch().strokeMask(p, c.m, ls, c.clip.r.Intersect(c.s.bounds()))
		if cut {
			sr.r.warnf(strokeCut)
		}
		c.s.fill(cov, sh, alpha*float32(st.strokeOpacity), bdf.BlendSourceOver, c.clip)
	}
}

// coverage rasterizes a path in user space under m within bounds.
func (sr *svgRenderer) coverage(p *path, m matrix, rule byte, bounds image.Rectangle) *mask {
	if hidden(p, 0, m, bounds) {
		return &mask{}
	}
	sc := sr.r.scratch()
	sc.polys = p.flattenTo(sc, sc.polys, m, flatTol)
	return sc.rasterize(sc.polys, rule, bounds)
}

// paintServer returns the shader of a fill or stroke value, or nil for none.
func (sr *svgRenderer) paintServer(v string, st svgStyle, bbox [4]float64, c *svgCtx) shader {
	v = strings.TrimSpace(v)
	if v == "" || v == "none" {
		return nil
	}
	if strings.HasPrefix(v, "url(") {
		fallback := ""
		if end := strings.IndexByte(v, ')'); end >= 0 {
			fallback = strings.TrimSpace(v[end+1:])
		}
		if g := sr.ref(v); g != nil && (g.name == "linearGradient" || g.name == "radialGradient") {
			if sh := sr.gradient(g, st, bbox, c); sh != nil {
				return sh
			}
			return nil
		}
		if g := sr.ref(v); g != nil && g.name == "pattern" {
			sr.r.warnf("SVG patterns are not drawn")
		}
		if fallback == "" {
			return nil
		}
		v = fallback
	}
	col, ok := parseColor(v, st.color)
	if !ok {
		return nil
	}
	return solid{float32(col.r * col.a), float32(col.g * col.a), float32(col.b * col.a), float32(col.a)}
}

// gradient builds the shader of a gradient element for a shape with the
// bounding box bbox (x, y, w, h).
func (sr *svgRenderer) gradient(g *svgNode, st svgStyle, bbox [4]float64, c *svgCtx) shader {
	// attributes and stops, following href to the gradients it extends
	attr := map[string]string{}
	var stops []*svgNode
	for n, i := g, 0; n != nil && i < 8; n, i = sr.ref(n.attr["href"]), i+1 {
		if n.name != "linearGradient" && n.name != "radialGradient" {
			break
		}
		for k, v := range n.attr {
			if _, ok := attr[k]; !ok {
				attr[k] = v
			}
		}
		if stops == nil {
			for _, s := range n.children {
				if s.name == "stop" {
					stops = append(stops, s)
				}
			}
		}
	}
	if len(stops) == 0 {
		return nil
	}
	bboxUnits := attr["gradientUnits"] != "userSpaceOnUse"
	if bboxUnits && !(bbox[2] > 0 && bbox[3] > 0) {
		return nil // a gradient in the units of an empty box does not paint
	}
	coord := func(k, def string, dir byte) float64 {
		v := attr[k]
		if v == "" {
			v = def
		}
		if bboxUnits {
			f, rest, _ := svgNumber(v)
			if strings.TrimSpace(rest) == "%" {
				f /= 100
			}
			return f
		}
		return c.length(v, dir, st.fontSize)
	}
	m := c.m
	if bboxUnits {
		m = m.mul(matrix{bbox[2], 0, 0, bbox[3], bbox[0], bbox[1]})
	}
	if t := attr["gradientTransform"]; t != "" {
		m = m.mul(parseTransform(t))
	}
	p := &bdf.Paint{}
	if g.name == "linearGradient" {
		p.Kind = bdf.PaintLinear
		p.Coords = []float32{float32(coord("x1", "0%", 'x')), float32(coord("y1", "0%", 'y')), float32(coord("x2", "100%", 'x')), float32(coord("y2", "0%", 'y'))}
	} else {
		cx, cy, r := coord("cx", "50%", 'x'), coord("cy", "50%", 'y'), coord("r", "50%", 'd')
		fx, fy := cx, cy
		if attr["fx"] != "" {
			fx = coord("fx", "", 'x')
		}
		if attr["fy"] != "" {
			fy = coord("fy", "", 'y')
		}
		p.Kind = bdf.PaintRadial
		p.Coords = []float32{float32(fx), float32(fy), float32(coord("fr", "0%", 'd')), float32(cx), float32(cy), float32(r)}
	}
	last := 0.0
	for _, s := range stops {
		off := 0.0
		if v, rest, ok := svgNumber(s.attr["offset"]); ok {
			off = v
			if strings.TrimSpace(rest) == "%" {
				off /= 100
			}
		}
		off = math.Max(last, math.Min(1, math.Max(0, off)))
		last = off
		sst := sr.compute(s, st)
		col := svgColor{0, 0, 0, 1}
		if v, ok := s.declared("stop-color"); ok {
			if cc, ok := parseColor(v, sst.color); ok {
				col = cc
			}
		}
		if v, ok := s.declared("stop-opacity"); ok {
			if f, rest, ok := svgNumber(v); ok {
				if strings.TrimSpace(rest) == "%" {
					f /= 100
				}
				col.a *= math.Min(1, math.Max(0, f))
			}
		}
		p.Stops = append(p.Stops, bdf.Stop{Offset: float32(off), Color: packColor(col)})
	}
	return newGradientShader(p, m, nil)
}

func packColor(c svgColor) bdf.Color {
	b := func(v float64) uint32 { return uint32(math.Round(math.Min(1, math.Max(0, v)) * 255)) }
	return bdf.Color(b(c.r)<<24 | b(c.g)<<16 | b(c.b)<<8 | b(c.a))
}

// clipMask rasterizes a clipPath for an element: the union of its shapes.
func (sr *svgRenderer) clipMask(cp, el *svgNode, c *svgCtx, st svgStyle) *mask {
	m := c.m
	if cp.attr["clipPathUnits"] == "objectBoundingBox" {
		if p := sr.shape(el, c, st); p != nil {
			x0, y0, x1, y1 := p.bbox()
			m = m.mul(matrix{x1 - x0, 0, 0, y1 - y0, x0, y0})
		}
	}
	if t := cp.attr["transform"]; t != "" {
		m = m.mul(parseTransform(t))
	}
	var out *mask
	bounds := c.s.bounds()
	for _, ch := range cp.children {
		if ch.name == "#text" {
			continue
		}
		if d, _ := ch.declared("display"); d == "none" {
			continue
		}
		// the shapes of a clip path count as elements drawn
		if sr.budget <= 0 {
			break
		}
		sr.budget--
		cst := sr.compute(ch, st)
		cm := m
		if t := ch.attr["transform"]; t != "" {
			cm = cm.mul(parseTransform(t))
		}
		shape := ch
		if ch.name == "use" {
			if t := sr.ref(ch.attr["href"]); t != nil && t.name != "use" {
				shape = t
				cm = cm.translate(c.length(ch.attr["x"], 'x', cst.fontSize), c.length(ch.attr["y"], 'y', cst.fontSize))
				cst = sr.compute(t, cst)
			}
		}
		cc := *c
		cc.m = cm
		p := sr.shape(shape, &cc, cst)
		if p == nil {
			continue
		}
		cov := sr.coverage(p, cm, cst.clipRule, bounds)
		if out == nil {
			out = cov
		} else {
			out = union(out, cov)
		}
	}
	if out == nil {
		return &mask{}
	}
	return out
}

// text draws a text element with its tspans. A tspan with x or y starts a
// new chunk (text-anchor places each chunk); other runs of text continue
// from the end of the one before.
func (sr *svgRenderer) text(n *svgNode, c *svgCtx, st svgStyle) {
	type piece struct {
		s  string
		st svgStyle
	}
	type chunk struct {
		x, y   float64
		pieces []piece
	}
	var chunks []chunk
	x, y := 0.0, 0.0
	var walk func(n *svgNode, st svgStyle)
	walk = func(n *svgNode, st svgStyle) {
		xs, ys := svgNumbers(n.attr["x"]), svgNumbers(n.attr["y"])
		if len(xs) > 0 || len(ys) > 0 || len(chunks) == 0 {
			if len(xs) > 0 {
				x = xs[0]
			}
			if len(ys) > 0 {
				y = ys[0]
			}
			if dx := svgNumbers(n.attr["dx"]); len(dx) > 0 {
				x += dx[0]
			}
			if dy := svgNumbers(n.attr["dy"]); len(dy) > 0 {
				y += dy[0]
			}
			chunks = append(chunks, chunk{x: x, y: y})
		}
		for _, ch := range n.children {
			switch ch.name {
			case "#text":
				if s := collapseSpace(ch.text); s != "" {
					last := &chunks[len(chunks)-1]
					last.pieces = append(last.pieces, piece{s, st})
				}
			case "tspan", "a", "textPath":
				if d, _ := ch.declared("display"); d != "none" {
					walk(ch, sr.compute(ch, st))
				}
			}
		}
	}
	walk(n, st)
	for ci := range chunks {
		ch := &chunks[ci]
		if len(ch.pieces) == 0 {
			continue
		}
		if ci == 0 {
			ch.pieces[0].s = strings.TrimLeft(ch.pieces[0].s, " ")
		}
		if ci == len(chunks)-1 {
			last := &ch.pieces[len(ch.pieces)-1]
			last.s = strings.TrimRight(last.s, " ")
		}
		var glyphs [][]placed
		var widths []float64
		total := 0.0
		for _, p := range ch.pieces {
			f := bdf.Font{Kind: bdf.FontSystem, Family: p.st.fontFamily, Weight: uint16(p.st.fontWeight)}
			if p.st.italic {
				f.Style = bdf.StyleItalic
			}
			g, w := sr.r.fonts.chain(&f).layout(p.s, p.st.fontSize, p.st.letterSpacing, bdf.DirLTR)
			glyphs, widths = append(glyphs, g), append(widths, w)
			total += w
		}
		pen := ch.x
		switch ch.pieces[0].st.anchor {
		case "middle":
			pen -= total / 2
		case "end":
			pen -= total
		}
		for i, p := range ch.pieces {
			sr.glyphs(glyphs[i], pen, ch.y, p.st, c)
			pen += widths[i]
		}
	}
}

func collapseSpace(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return ' '
		}
		return r
	}, s)
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

// glyphs fills (and strokes) laid-out glyphs at (x, y) on the baseline.
func (sr *svgRenderer) glyphs(gs []placed, x, y float64, st svgStyle, c *svgCtx) {
	if !st.visible || len(gs) == 0 {
		return
	}
	p := &path{}
	for _, g := range gs {
		o := g.use.face.outline(g.gid)
		if o == nil {
			continue
		}
		k := st.fontSize / g.use.face.upem
		gm := matrix{1, 0, 0, 1, x + g.x, y}.scale(k, -k)
		if g.use.italic {
			gm = gm.mul(matrix{1, 0, fakeItalicSkew, 1, 0, 0})
		}
		for _, s := range o.subs {
			ns := subpath{closed: s.closed}
			ns.start.x, ns.start.y = gm.apply(s.start.x, s.start.y)
			for _, sg := range s.segs {
				var q segment
				q.kind = sg.kind
				for i := 0; i < 3; i++ {
					q.p[2*i], q.p[2*i+1] = gm.apply(sg.p[2*i], sg.p[2*i+1])
				}
				ns.segs = append(ns.segs, q)
			}
			p.subs = append(p.subs, ns)
		}
	}
	st.fillRule = bdf.NonZero
	sr.paint(p, c, st)
}

// image draws an image element: an embedded raster image or SVG document
// (data: URL) fitted into its box.
func (sr *svgRenderer) image(n *svgNode, c *svgCtx, st svgStyle) {
	x, y := c.length(n.attr["x"], 'x', st.fontSize), c.length(n.attr["y"], 'y', st.fontSize)
	var pic *picture
	var iw, ih float64
	if !n.tried {
		// the image of the element is kept for the next time it is drawn
		n.tried = true
		data, ok := dataURL(n.attr["href"])
		if !ok {
			if n.attr["href"] != "" {
				sr.r.warnf("SVG images linking to other files are not drawn")
			}
			return
		}
		if isSVG(data) {
			si, ok := newSVGImage(data)
			if !ok {
				return
			}
			for _, w := range si.doc.warnings {
				sr.r.warnf("%s", w)
			}
			n.image = si
		} else {
			var full bool
			if n.pic, full = sr.r.decode(data, func() { n.pic, n.tried = nil, false }); n.pic == nil {
				if !full {
					sr.r.warnf("an image in an SVG image cannot be decoded")
				}
				n.tried = !full // there may be room for it in the next image drawn
				return
			}
		}
	}
	if n.pic == nil && n.image == nil {
		return
	}
	if si := n.image; si != nil {
		if sr.depth > 8 {
			return
		}
		if sr.inside >= svgInside {
			sr.r.warnf("SVG images are more than %d deep in SVG images: the rest are not drawn", svgInside)
			return
		}
		iw, ih = si.w, si.h
		w, h := sizeOr(c, n, st, iw, ih)
		if pic, _ = si.rasterIn(sr.r, c.m.maxScale()*math.Max(w/iw, h/ih), sr.inside+1); pic == nil {
			return
		}
	} else {
		pic = n.pic
		pic.used = sr.r.drawing
		iw, ih = float64(pic.w), float64(pic.h)
	}
	w, h := sizeOr(c, n, st, iw, ih)
	if !(w > 0 && h > 0) {
		return
	}
	m := c.m.mul(aspectTransform([4]float64{0, 0, iw, ih}, n.attr["preserveAspectRatio"], x, y, w, h))
	// image pixels of the raster → device
	full := m.mul(matrix{iw / float64(pic.w), 0, 0, ih / float64(pic.h), 0, 0})
	inv, ok := full.invert()
	if !ok {
		return
	}
	box := &path{}
	box.rect(0, 0, iw, ih)
	cov := rasterize(box.flatten(m, flatTol), bdf.NonZero, c.clip.r.Intersect(c.s.bounds()))
	vp := &path{}
	vp.rect(x, y, w, h)
	cov = intersect(cov, rasterize(vp.flatten(c.m, flatTol), bdf.NonZero, c.s.bounds()))
	sh := &imageShader{inv: inv, img: pic, level: pic.levelFor(inv), smooth: true, sx1: float64(pic.w), sy1: float64(pic.h)}
	c.s.fill(cov, sh, 1, bdf.BlendSourceOver, c.clip)
}

// sizeOr returns the width and height of an image element: its attributes,
// or the image's own size for those missing (keeping its proportions).
func sizeOr(c *svgCtx, n *svgNode, st svgStyle, iw, ih float64) (float64, float64) {
	w, h := c.length(n.attr["width"], 'x', st.fontSize), c.length(n.attr["height"], 'y', st.fontSize)
	switch {
	case w > 0 && h > 0:
		return w, h
	case w > 0:
		return w, w * ih / iw
	case h > 0:
		return h * iw / ih, h
	}
	return iw, ih
}

// dataURL decodes a data: URL.
func dataURL(s string) ([]byte, bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "data:") {
		return nil, false
	}
	meta, payload, ok := strings.Cut(s[5:], ",")
	if !ok {
		return nil, false
	}
	if strings.HasSuffix(meta, ";base64") {
		payload = strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\r' || r == '\t' {
				return -1
			}
			return r
		}, payload)
		b, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			if b, err = base64.RawStdEncoding.DecodeString(strings.TrimRight(payload, "=")); err != nil {
				return nil, false
			}
		}
		return b, true
	}
	u, err := url.PathUnescape(payload)
	if err != nil {
		return nil, false
	}
	return []byte(u), true
}

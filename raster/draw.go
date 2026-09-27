package raster

import (
	"image"
	"math"

	"github.com/shibukawa/bdf"
)

// state is the drawing state of a context (what save and restore keep).
type state struct {
	m      matrix
	clip   *mask
	fill   style
	stroke style
	line   lineStyle
	alpha  float32
	blend  byte

	shadow             bdf.Color
	shadowBlur         float64 // device pixels
	shadowDX, shadowDY float64

	font     *bdf.Font
	size     float64
	align    byte
	baseline byte
	dir      byte
	spacing  float64
	smooth   bool
}

// initialState is the Canvas 2D initial state with a transform and clip.
func initialState(m matrix, clip *mask) state {
	return state{
		m: m, clip: clip,
		fill: style{color: 0x000000ff}, stroke: style{color: 0x000000ff},
		line:  lineStyle{width: 1, miter: 10},
		alpha: 1, size: 10, smooth: true,
	}
}

func (s *state) shadowActive() bool {
	return uint32(s.shadow)&0xff != 0 && (s.shadowBlur > 0 || s.shadowDX != 0 || s.shadowDY != 0)
}

// context is a canvas being drawn on: the page, or the temporary canvas of
// a group or soft mask.
type context struct {
	target *surface
	st     state
	stack  []state
}

type group struct {
	parent *context
	ctx    *context
	x, y   int // the group canvas's place on the parent's
	alpha  float32
	blend  byte
}

type softMask struct {
	parent   *context
	ctx      *context
	kind     byte
	transfer []byte
	group    *group
}

// drawer runs the instructions of objects on a context.
type drawer struct {
	r      *Renderer
	ctx    *context
	obj    *bdf.ObjectPart
	groups []*group
	masks  []*softMask
	depth  int
}

// drawTop draws a top-level object (a page layer or a sheet tile) from the
// initial state with a transform and a clip.
func (d *drawer) drawTop(o *bdf.ObjectPart, m matrix, clip *mask) {
	saved := d.ctx.st
	d.ctx.st = initialState(m, clip)
	stack := d.ctx.stack
	d.ctx.stack = nil
	d.run(o)
	// unwind groups and masks left open by a malformed stream
	d.masks = d.masks[:0]
	for len(d.groups) > 0 {
		d.groupEnd()
	}
	d.ctx.st = saved
	d.ctx.stack = stack
}

func (d *drawer) run(o *bdf.ObjectPart) {
	if d.depth > 64 {
		d.r.warnf("objects nested too deep")
		return
	}
	d.depth++
	prev := d.obj
	d.obj = o
	if err := o.Walk(d.instr); err != nil {
		d.r.warnf("object: %v", err)
	}
	d.obj = prev
	d.depth--
}

func f64(in bdf.Instr, i int) float64 {
	if v, ok := in.Args[i].(float32); ok {
		return float64(v)
	}
	return 0
}

func u64(in bdf.Instr, i int) uint64 {
	if v, ok := in.Args[i].(uint64); ok {
		return v
	}
	return 0
}

func (d *drawer) instr(in bdf.Instr) {
	ctx := d.ctx
	st := &ctx.st
	switch in.Op {
	case bdf.OpSave:
		ctx.stack = append(ctx.stack, *st)
	case bdf.OpRestore:
		if n := len(ctx.stack); n > 0 {
			*st = ctx.stack[n-1]
			ctx.stack = ctx.stack[:n-1]
		}
	case bdf.OpTransform:
		st.m = st.m.mul(matrix{f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4), f64(in, 5)})
	case bdf.OpTranslate:
		st.m = st.m.translate(f64(in, 0), f64(in, 1))
	case bdf.OpScale:
		st.m = st.m.scale(f64(in, 0), f64(in, 1))
	case bdf.OpClipPath:
		if p := d.path(int(u64(in, 0))); p != nil {
			st.clip = intersect(st.clip, d.coverage(p, byte(u64(in, 1)), st.m))
		}
	case bdf.OpClipRect:
		p := &path{}
		p.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3))
		st.clip = intersect(st.clip, d.coverage(p, bdf.NonZero, st.m))
	case bdf.OpFillColor:
		st.fill = style{color: bdf.Color(u64(in, 0))}
	case bdf.OpStrokeColor:
		st.stroke = style{color: bdf.Color(u64(in, 0))}
	case bdf.OpFillPaint:
		st.fill = d.paintStyle(int(u64(in, 0)))
	case bdf.OpStrokePaint:
		st.stroke = d.paintStyle(int(u64(in, 0)))
	case bdf.OpLine:
		if w := f64(in, 0); w > 0 && !math.IsInf(w, 0) {
			st.line.width = w
		}
		st.line.cap, st.line.join = byte(u64(in, 1)), byte(u64(in, 2))
		if m := f64(in, 3); m > 0 && !math.IsInf(m, 0) {
			st.line.miter = m
		}
	case bdf.OpDash:
		segs, _ := in.Args[0].([]float32)
		off, _ := in.Args[1].(float32)
		st.line.setDash(segs, off)
	case bdf.OpAlpha:
		if a := f64(in, 0); a >= 0 && a <= 1 {
			st.alpha = float32(a)
		}
	case bdf.OpBlend:
		if b := u64(in, 0); b < uint64(len(bdf.BlendNames)) {
			st.blend = byte(b)
		}
	case bdf.OpShadow:
		s := st.m.areaScale()
		if s == 0 || math.IsNaN(s) {
			s = 1
		}
		st.shadow = bdf.Color(u64(in, 0))
		st.shadowBlur, st.shadowDX, st.shadowDY = math.Max(0, f64(in, 1)*s), f64(in, 2)*s, f64(in, 3)*s
	case bdf.OpFilter:
		if s, _ := in.Args[0].(string); s != "" && s != "none" {
			d.r.warnf("FILTER %q is not applied", s)
		}
	case bdf.OpFont:
		if i := int(u64(in, 0)); i < len(d.obj.Fonts) {
			f := d.obj.Fonts[i]
			st.font = &f
		}
		st.size = f64(in, 1)
	case bdf.OpTextStyle:
		st.align, st.baseline, st.dir = byte(u64(in, 0)), byte(u64(in, 1)), byte(u64(in, 2))
		st.spacing = f64(in, 3)
	case bdf.OpFillRect:
		d.fillRect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3))
	case bdf.OpStrokeRect:
		p := &path{}
		p.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3))
		d.strokePath(p, st.m)
	case bdf.OpClearRect:
		p := &path{}
		p.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3))
		ctx.target.fill(d.coverage(p, bdf.NonZero, st.m), solid{0, 0, 0, 1}, 1, bdf.BlendDestinationOut, st.clip)
	case bdf.OpFillPath:
		if p := d.path(int(u64(in, 0))); p != nil {
			d.fillPath(p, byte(u64(in, 1)), st.m)
		}
	case bdf.OpStrokePath:
		if p := d.path(int(u64(in, 0))); p != nil {
			d.strokePath(p, st.m)
		}
	case bdf.OpFillPathAt:
		if p := d.path(int(u64(in, 0))); p != nil {
			d.fillPath(p, byte(u64(in, 1)), st.m.translate(f64(in, 2), f64(in, 3)))
		}
	case bdf.OpFillPathRun:
		rule := byte(u64(in, 0))
		glyphs, _ := in.Args[1].([]bdf.Glyph)
		var polys []polyline
		for _, g := range glyphs {
			if p := d.path(int(g.Path)); p != nil {
				m := st.m.translate(float64(g.X), float64(g.Y))
				polys = append(polys, p.flatten(m, flatTol)...)
			}
		}
		d.paint(rasterize(polys, rule, d.bounds()), d.shader(&st.fill, st.m), 1)
	case bdf.OpFillText, bdf.OpStrokeText:
		s, _ := in.Args[0].(string)
		d.text(s, f64(in, 1), f64(in, 2), f64(in, 3), in.Op == bdf.OpStrokeText)
	case bdf.OpImage:
		d.image(int(u64(in, 0)), nil, f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4))
	case bdf.OpImageSub:
		src := &[4]float64{f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4)}
		d.image(int(u64(in, 0)), src, f64(in, 5), f64(in, 6), f64(in, 7), f64(in, 8))
	case bdf.OpSmoothing:
		st.smooth = u64(in, 0) != 0
	case bdf.OpUse, bdf.OpUseAt:
		i := int(u64(in, 0))
		if i >= len(d.obj.Objects) {
			return
		}
		child := d.r.object(d.obj.Objects[i])
		if child == nil {
			return
		}
		saved := *st
		depth := len(ctx.stack)
		if in.Op == bdf.OpUseAt {
			st.m = st.m.translate(f64(in, 1), f64(in, 2))
		}
		d.run(child)
		// USE is wrapped in an implicit save and restore
		d.ctx.stack = d.ctx.stack[:min(depth, len(d.ctx.stack))]
		d.ctx.st = saved
	case bdf.OpGroupBegin:
		d.groupBegin(float32(f64(in, 0)), byte(u64(in, 1)), f64(in, 2), f64(in, 3), f64(in, 4), f64(in, 5))
	case bdf.OpGroupEnd:
		d.groupEnd()
	case bdf.OpMaskBegin:
		transfer, _ := in.Args[2].([]byte)
		d.maskBegin(byte(u64(in, 0)), bdf.Color(u64(in, 1)), transfer)
	case bdf.OpMaskEnd:
		d.maskEnd()
	}
}

// flatTol is how far flattened curves may stray from the true ones, in
// device pixels.
const flatTol = 0.1

// bounds is the area drawing can change: the target within the clip.
func (d *drawer) bounds() image.Rectangle {
	return d.ctx.st.clip.r.Intersect(d.ctx.target.bounds())
}

func (d *drawer) path(i int) *path {
	if i < 0 || i >= len(d.obj.Paths) {
		d.r.warnf("bad path reference")
		return nil
	}
	return d.r.path(d.obj, i)
}

// coverage rasterizes a path in user space under a transform.
func (d *drawer) coverage(p *path, rule byte, m matrix) *mask {
	if !m.finite() {
		return &mask{}
	}
	// an axis-aligned rectangle: exact coverage without scanning
	if len(p.subs) >= 1 && m[1] == 0 && m[2] == 0 {
		if x0, y0, x1, y1, ok := p.isRect(); ok {
			ax, ay := m.apply(x0, y0)
			bx, by := m.apply(x1, y1)
			return rectMask(ax, ay, bx, by, d.bounds())
		}
	}
	return rasterize(p.flatten(m, flatTol), rule, d.bounds())
}

func (d *drawer) fillRect(x, y, w, h float64) {
	p := &path{}
	p.rect(x, y, w, h)
	d.fillPath(p, bdf.NonZero, d.ctx.st.m)
}

func (d *drawer) fillPath(p *path, rule byte, m matrix) {
	st := &d.ctx.st
	d.paint(d.coverage(p, rule, m), d.shader(&st.fill, m), 1)
}

// strokePath strokes a path in user space under m. Strokes thinner than a
// device pixel are drawn a pixel wide with their width as alpha, as Skia
// draws hairlines.
func (d *drawer) strokePath(p *path, m matrix) {
	cov, alpha := strokeMask(p, m, d.ctx.st.line, d.bounds())
	d.paint(cov, d.shader(&d.ctx.st.stroke, m), alpha)
}

// strokeMask rasterizes the stroke of a path in user space under m within
// bounds, and returns the alpha to draw it with: strokes thinner than a
// device pixel are a pixel wide with their width as alpha.
func strokeMask(p *path, m matrix, ls lineStyle, bounds image.Rectangle) (*mask, float32) {
	if !m.finite() {
		return &mask{}, 0
	}
	scale := m.areaScale()
	if scale == 0 {
		return &mask{}, 0
	}
	alpha := float32(1)
	if dw := ls.width * scale; dw < 1 {
		alpha = float32(dw)
		ls.width = 1 / scale
	}
	ms := m.maxScale()
	lines := p.flatten(identity, flatTol/ms)
	outline := strokePolys(lines, &ls, ms)
	for i := range outline {
		for j, q := range outline[i].pts {
			x, y := m.apply(q.x, q.y)
			outline[i].pts[j] = point{x, y}
		}
	}
	return rasterize(outline, bdf.NonZero, bounds), alpha
}

// shader returns the shader of a fill or stroke style drawn under m.
func (d *drawer) shader(s *style, m matrix) shader {
	if s.paint == nil {
		return premul(s.color)
	}
	if s.paint.Kind == bdf.PaintPattern {
		return newPatternShader(d.r, s, m)
	}
	return newGradientShader(s.paint, m)
}

func (d *drawer) paintStyle(i int) style {
	if i < 0 || i >= len(d.obj.Paints) {
		d.r.warnf("bad paint reference")
		return style{}
	}
	p := &d.obj.Paints[i]
	s := style{paint: p}
	if p.Kind == bdf.PaintPattern {
		if int(p.Image) < len(d.obj.Images) {
			h := d.obj.Images[p.Image]
			if s.image = d.r.picture(h); s.image == nil {
				s.svg = d.r.svg(h)
			}
		}
	}
	return s
}

// paint composites a shader through a coverage with the state's alpha,
// composite operation, clip and shadow.
func (d *drawer) paint(cov *mask, sh shader, alpha float32) {
	st := &d.ctx.st
	if cov.empty() && !unbounded(st.blend) {
		return
	}
	if st.shadowActive() && !cov.empty() {
		d.shadow(cov, sh, alpha*st.alpha)
	}
	d.ctx.target.fill(cov, sh, alpha*st.alpha, st.blend, st.clip)
}

// shadow draws the shadow of what a draw covers: its alpha, blurred and
// offset, in the shadow colour.
func (d *drawer) shadow(cov *mask, sh shader, alpha float32) {
	st := &d.ctx.st
	sigma := st.shadowBlur / 2
	pad := int(math.Ceil(3 * sigma))
	r := cov.r
	w, h := r.Dx(), r.Dy()
	buf := make([]float32, 4*w)
	a := make([]float32, (w+2*pad)*(h+2*pad))
	stride := w + 2*pad
	for y := r.Min.Y; y < r.Max.Y; y++ {
		sh.shade(r.Min.X, y, w, buf)
		row := cov.row(y, r.Min.X, r.Max.X)
		for x := 0; x < w; x++ {
			c := float32(1)
			if row != nil {
				c = row[x]
			}
			a[(y-r.Min.Y+pad)*stride+x+pad] = buf[4*x+3] * c * alpha
		}
	}
	if sigma > 0 {
		blur(a, stride, h+2*pad, sigma)
	}
	dx, dy := int(math.Round(st.shadowDX)), int(math.Round(st.shadowDY))
	sm := &mask{r: image.Rect(r.Min.X-pad+dx, r.Min.Y-pad+dy, r.Max.X+pad+dx, r.Max.Y+pad+dy), a: a}
	d.ctx.target.fill(sm, premul(st.shadow), 1, st.blend, st.clip)
}

// text draws a FILL_TEXT or STROKE_TEXT.
func (d *drawer) text(s string, x, y, advance float64, stroke bool) {
	st := &d.ctx.st
	if st.size <= 0 || s == "" {
		return
	}
	ch := d.r.fonts.chain(st.font)
	glyphs, width := ch.layout(s, st.size, st.spacing, st.dir)
	if len(glyphs) == 0 {
		return
	}
	m := st.m.translate(x, y)
	if advance > 0 && width > 0 && math.Abs(width-advance)/advance > 0.005 {
		m = m.scale(advance/width, 1)
	}
	switch st.align {
	case bdf.AlignRight:
		m = m.translate(-width, 0)
	case bdf.AlignCenter:
		m = m.translate(-width/2, 0)
	case bdf.AlignEnd:
		if st.dir != bdf.DirRTL {
			m = m.translate(-width, 0)
		}
	case bdf.AlignStart:
		if st.dir == bdf.DirRTL {
			m = m.translate(-width, 0)
		}
	}
	m = m.translate(0, ch.baselineShift(st.baseline, st.size))
	var fill, bold []polyline
	var strokeLines []polyline
	boldW := fakeBoldWidth(st.size * m.areaScale())
	for _, g := range glyphs {
		p := g.use.face.outline(g.gid)
		if p == nil {
			continue
		}
		k := st.size / g.use.face.upem
		gm := m.translate(g.x, 0).scale(k, -k)
		if g.use.italic {
			gm = gm.mul(matrix{1, 0, fakeItalicSkew, 1, 0, 0})
		}
		if stroke {
			// stroke in user space: the line width is in units
			um := matrix{1, 0, 0, 1, 0, 0}.translate(g.x, 0).scale(k, -k)
			if g.use.italic {
				um = um.mul(matrix{1, 0, fakeItalicSkew, 1, 0, 0})
			}
			strokeLines = append(strokeLines, p.flatten(um, flatTol/math.Max(m.maxScale(), 1e-9))...)
			continue
		}
		polys := p.flatten(gm, flatTol)
		fill = append(fill, polys...)
		if g.use.bold && boldW > 0 {
			bold = append(bold, strokePolys(polys, &lineStyle{width: boldW, join: bdf.JoinRound, miter: 10}, 1)...)
		}
	}
	if stroke {
		ls := st.line
		scale := m.areaScale()
		alpha := float32(1)
		if dw := ls.width * scale; dw < 1 && scale > 0 {
			alpha, ls.width = float32(dw), 1/scale
		}
		outline := strokePolys(strokeLines, &ls, m.maxScale())
		for i := range outline {
			for j, q := range outline[i].pts {
				px, py := m.apply(q.x, q.y)
				outline[i].pts[j] = point{px, py}
			}
		}
		d.paint(rasterize(outline, bdf.NonZero, d.bounds()), d.shader(&st.stroke, m), alpha)
		return
	}
	cov := rasterize(fill, bdf.NonZero, d.bounds())
	if len(bold) > 0 {
		cov = union(cov, rasterize(bold, bdf.NonZero, d.bounds()))
	}
	d.paint(cov, d.shader(&st.fill, m), 1)
}

// union returns the larger coverage of two masks at each pixel.
func union(a, b *mask) *mask {
	if a.empty() {
		return b
	}
	if b.empty() {
		return a
	}
	r := a.r.Union(b.r)
	out := &mask{r: r, a: make([]float32, r.Dx()*r.Dy())}
	for _, m := range []*mask{a, b} {
		for y := m.r.Min.Y; y < m.r.Max.Y; y++ {
			row := m.row(y, m.r.Min.X, m.r.Max.X)
			o := (y-r.Min.Y)*r.Dx() + m.r.Min.X - r.Min.X
			for x := 0; x < m.r.Dx(); x++ {
				v := float32(1)
				if row != nil {
					v = row[x]
				}
				out.a[o+x] = max(out.a[o+x], v)
			}
		}
	}
	return out
}

// image draws an image, or the part src (x, y, w, h in image pixels) of
// it, into the rectangle (x, y, w, h).
func (d *drawer) image(i int, src *[4]float64, x, y, w, h float64) {
	if i < 0 || i >= len(d.obj.Images) {
		d.r.warnf("bad image reference")
		return
	}
	hash := d.obj.Images[i]
	pic := d.r.picture(hash)
	si := d.r.svg(hash)
	if pic == nil && si == nil || w == 0 || h == 0 {
		return
	}
	st := &d.ctx.st
	// the source rectangle is in image pixels: an SVG image's natural size
	nw, nh := 0.0, 0.0
	if pic != nil {
		nw, nh = float64(pic.w), float64(pic.h)
	} else {
		nw, nh = si.w, si.h
	}
	sx, sy, sw, sh := 0.0, 0.0, nw, nh
	if src != nil {
		sx, sy, sw, sh = src[0], src[1], src[2], src[3]
		if sw == 0 || sh == 0 {
			return
		}
	}
	if pic == nil {
		// drawn as sharp as it is shown: a raster at the device scale
		var k float64
		pic, k = si.raster(d.r, st.m.maxScale()*math.Max(math.Abs(w/sw), math.Abs(h/sh)))
		ky := float64(pic.h) / si.h
		sx, sy, sw, sh = sx*k, sy*ky, sw*k, sh*ky
	}
	// device → image pixels: invert (m × dest rect ← source rect)
	full := st.m.mul(matrix{w / sw, 0, 0, h / sh, x - sx*w/sw, y - sy*h/sh})
	inv, ok := full.invert()
	if !ok {
		return
	}
	p := &path{}
	p.rect(x, y, w, h)
	lo := func(a, b float64) float64 { return math.Min(a, a+b) }
	sh2 := &imageShader{
		inv: inv, img: pic, level: pic.levelFor(inv), smooth: st.smooth,
		sx0: math.Max(0, lo(sx, sw)), sy0: math.Max(0, lo(sy, sh)),
		sx1: math.Min(float64(pic.w), lo(sx, sw)+math.Abs(sw)), sy1: math.Min(float64(pic.h), lo(sy, sh)+math.Abs(sh)),
	}
	d.paint(d.coverage(p, bdf.NonZero, st.m), sh2, 1)
}

func (d *drawer) groupBegin(alpha float32, blend byte, x, y, w, h float64) {
	st := &d.ctx.st
	m := st.m
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, c := range [][2]float64{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}} {
		px, py := m.apply(c[0], c[1])
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	r := image.Rect(int(math.Floor(clampF(x0))), int(math.Floor(clampF(y0))), int(math.Ceil(clampF(x1))), int(math.Ceil(clampF(y1))))
	// nothing outside the parent canvas or its clip can show
	r = r.Intersect(d.bounds())
	if r.Empty() {
		r = image.Rect(0, 0, 1, 1)
	}
	gs := initialState(matrix{1, 0, 0, 1, -float64(r.Min.X), -float64(r.Min.Y)}.mul(m), nil)
	gs.fill, gs.stroke, gs.line.width, gs.font, gs.size = st.fill, st.stroke, st.line.width, st.font, st.size
	gc := &context{target: newSurface(r.Dx(), r.Dy()), st: gs}
	gc.st.clip = &mask{r: gc.target.bounds()}
	g := &group{parent: d.ctx, ctx: gc, x: r.Min.X, y: r.Min.Y, alpha: alpha, blend: blend}
	d.groups = append(d.groups, g)
	d.ctx = gc
}

func (d *drawer) groupEnd() {
	n := len(d.groups)
	if n == 0 {
		return
	}
	g := d.groups[n-1]
	d.groups = d.groups[:n-1]
	// masks begun inside the group end with it
	for len(d.masks) > 0 && d.masks[len(d.masks)-1].group == g {
		d.masks = d.masks[:len(d.masks)-1]
	}
	d.ctx = g.parent
	s := g.ctx.target
	r := image.Rect(g.x, g.y, g.x+s.w, g.y+s.h)
	blend := g.blend
	if int(blend) >= len(bdf.BlendNames) {
		blend = bdf.BlendSourceOver
	}
	d.ctx.target.fill(&mask{r: r}, surfaceShader{s: s, dx: g.x, dy: g.y}, g.alpha, blend, d.ctx.st.clip)
}

func (d *drawer) maskBegin(kind byte, backdrop bdf.Color, transfer []byte) {
	var g *group
	w, h := 1, 1
	if n := len(d.groups); n > 0 {
		g = d.groups[n-1]
		w, h = g.ctx.target.w, g.ctx.target.h
	}
	mc := &context{target: newSurface(w, h)}
	if kind == bdf.MaskLuminosity {
		c := premul(backdrop | 0xff)
		for i := 0; i < w*h; i++ {
			copy(mc.target.pix[4*i:], c[:])
		}
	}
	mc.st = initialState(d.ctx.st.m, &mask{r: mc.target.bounds()})
	if len(transfer) != 256 {
		transfer = nil
	}
	d.masks = append(d.masks, &softMask{parent: d.ctx, ctx: mc, kind: kind, transfer: transfer, group: g})
	d.ctx = mc
}

func (d *drawer) maskEnd() {
	n := len(d.masks)
	if n == 0 {
		return
	}
	m := d.masks[n-1]
	d.masks = d.masks[:n-1]
	d.ctx = m.parent
	g := m.group
	if g == nil || len(d.groups) == 0 || d.groups[len(d.groups)-1] != g {
		return
	}
	ms, gs := m.ctx.target, g.ctx.target
	for i := 0; i < ms.w*ms.h; i++ {
		p := ms.pix[4*i : 4*i+4]
		var v int
		if m.kind == bdf.MaskLuminosity {
			v = int(math.Round(float64(0.3*p[0]+0.59*p[1]+0.11*p[2]) * 255))
		} else {
			v = int(math.Round(float64(p[3]) * 255))
		}
		v = min(max(v, 0), 255)
		if m.transfer != nil {
			v = int(m.transfer[v])
		}
		k := float32(v) / 255
		q := gs.pix[4*i : 4*i+4]
		q[0], q[1], q[2], q[3] = q[0]*k, q[1]*k, q[2]*k, q[3]*k
	}
}

// blur applies an approximate Gaussian blur (three box blurs) to a plane.
func blur(a []float32, w, h int, sigma float64) {
	// box sizes for three passes (Kovesi)
	wIdeal := math.Sqrt(12*sigma*sigma/3 + 1)
	wl := int(math.Floor(wIdeal))
	if wl%2 == 0 {
		wl--
	}
	wu := wl + 2
	mIdeal := (12*sigma*sigma - 3*float64(wl*wl) - 12*float64(wl) - 9) / (-4*float64(wl) - 4)
	mm := int(math.Round(mIdeal))
	tmp := make([]float32, len(a))
	for i := 0; i < 3; i++ {
		size := wl
		if i >= mm {
			size = wu
		}
		r := (size - 1) / 2
		if r <= 0 {
			continue
		}
		boxH(a, tmp, w, h, r)
		boxV(tmp, a, w, h, r)
	}
}

func boxH(src, dst []float32, w, h, r int) {
	k := 1 / float32(2*r+1)
	for y := 0; y < h; y++ {
		row := src[y*w : (y+1)*w]
		out := dst[y*w : (y+1)*w]
		var sum float32
		for x := -r; x <= r; x++ {
			if x >= 0 && x < w {
				sum += row[x]
			}
		}
		for x := 0; x < w; x++ {
			out[x] = sum * k
			if x+r+1 < w {
				sum += row[x+r+1]
			}
			if x-r >= 0 {
				sum -= row[x-r]
			}
		}
	}
}

func boxV(src, dst []float32, w, h, r int) {
	k := 1 / float32(2*r+1)
	for x := 0; x < w; x++ {
		var sum float32
		for y := -r; y <= r; y++ {
			if y >= 0 && y < h {
				sum += src[y*w+x]
			}
		}
		for y := 0; y < h; y++ {
			dst[y*w+x] = sum * k
			if y+r+1 < h {
				sum += src[(y+r+1)*w+x]
			}
			if y-r >= 0 {
				sum -= src[(y-r)*w+x]
			}
		}
	}
}

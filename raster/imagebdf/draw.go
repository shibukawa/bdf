package imagebdf

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

	// made is the clip if it was made after the SAVE of this state: its
	// memory is let go of with the state (see drawer.live)
	made *mask
}

// Limits of what the drawing of an image takes. A document decides how
// often its objects are drawn, how deep its groups are and how many states
// it keeps, so a small one could otherwise ask for any time and memory.
const (
	// maxReusedInstructions bounds the instructions drawn from objects
	// that were drawn before in the same image. What is drawn once is
	// bounded by the size of the document; this bounds what a document
	// adds by drawing it again (bdf.MaxReusedInstructions is the limit
	// of reading the text, where an instruction costs far less).
	maxReusedInstructions = 1 << 24
	// maxLayers bounds the groups and soft masks that are open at a time.
	maxLayers = 64
	// maxStates bounds the states a canvas keeps for RESTORE.
	maxStates = 1024
	// maxLive bounds the memory of the canvases of the open groups and
	// soft masks and of the clips of the states kept, in canvases of the
	// image drawn (16 bytes a pixel), and minLive is what they may take of
	// a small image all the same, in bytes.
	maxLive = 8
	minLive = 64 << 20
	// maxShadowBlur bounds the blur of a shadow, in device pixels: the
	// shadow of a drawing takes a canvas three times as much larger on
	// every side.
	maxShadowBlur = 1024
)

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
	// unsaved counts the SAVEs that kept no state (past maxStates, or in
	// what is left out): as many RESTOREs restore none
	unsaved int
	// held is the memory counted in drawer.live for this canvas: its
	// pixels and the clips its states made
	held int
	// done says that the group or soft mask of the canvas has ended
	done bool
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

	// seen holds the objects drawn in this image, again says that the
	// object being drawn (or one that draws it) was drawn before, and
	// reused counts the instructions of such objects
	seen   map[bdf.Hash]struct{}
	again  bool
	reused int
	// skipped counts the groups and soft masks that are open but not
	// drawn: those past the limits, and those of objects past theirs
	skipped int
	// live is the memory of the canvases of the open groups and soft
	// masks and of the clips of the states kept, in bytes, and limit the
	// most it may be
	live, limit int
	// most is the most reused may be (maxReusedInstructions)
	most int
}

// enter reports whether an object was drawn before in this image, and
// notes that it is now.
func (d *drawer) enter(h bdf.Hash) bool {
	if _, ok := d.seen[h]; ok {
		return true
	}
	if d.seen == nil {
		d.seen = map[bdf.Hash]struct{}{}
	}
	d.seen[h] = struct{}{}
	return false
}

// hold counts n bytes of a canvas or a clip of ctx in the memory in use,
// and reports whether there is room for them.
func (d *drawer) hold(ctx *context, n int) bool {
	if n > d.limit-d.live {
		return false
	}
	ctx.held += n
	d.live += n
	return true
}

// drop lets go of the clip a state made.
func (d *drawer) drop(ctx *context, st *state) {
	if st.made != nil {
		n := st.made.bytes()
		ctx.held -= n
		d.live -= n
		st.made = nil
	}
}

// closed lets go of what a canvas that is done with held.
func (d *drawer) closed(ctx *context) {
	d.live -= ctx.held
	ctx.held, ctx.done = 0, true
}

// drawTop draws a top-level object (a page layer or a sheet tile) from the
// initial state with a transform and a clip.
func (d *drawer) drawTop(h bdf.Hash, m matrix, clip *mask) {
	o := d.r.object(h)
	if o == nil {
		return
	}
	top := d.ctx
	saved, stack, unsaved := top.st, top.stack, top.unsaved
	top.st = initialState(m, clip)
	top.stack, top.unsaved = nil, 0
	d.again = d.enter(h)
	d.run(o)
	// unwind groups and masks left open by a malformed stream
	d.masks = d.masks[:0]
	for len(d.groups) > 0 {
		d.groupEnd()
	}
	d.ctx = top
	d.skipped, d.live, top.held = 0, 0, 0
	top.st, top.stack, top.unsaved = saved, stack, unsaved
}

func (d *drawer) run(o *bdf.ObjectPart) {
	if d.depth > bdf.MaxUseDepth {
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

// leftOut takes an instruction of what is not drawn: a group or soft mask
// past the limits, or an object drawn again past maxReusedInstructions. It
// counts what begins and ends, so that what ends later ends what it began.
func (d *drawer) leftOut(in bdf.Instr) {
	switch in.Op {
	case bdf.OpSave:
		d.ctx.unsaved++
	case bdf.OpRestore:
		d.restore()
	case bdf.OpGroupBegin, bdf.OpMaskBegin:
		d.skipped++
	case bdf.OpGroupEnd:
		if d.skipped > 0 {
			d.skipped--
		} else {
			d.groupEnd()
		}
	case bdf.OpMaskEnd:
		if d.skipped > 0 {
			d.skipped--
		} else {
			d.maskEnd()
		}
	}
}

// restore takes a RESTORE.
func (d *drawer) restore() {
	ctx := d.ctx
	if ctx.unsaved > 0 {
		ctx.unsaved--
		return
	}
	if n := len(ctx.stack); n > 0 {
		d.drop(ctx, &ctx.st)
		ctx.st = ctx.stack[n-1]
		ctx.stack = ctx.stack[:n-1]
	}
}

// clip cuts the clip of the state to a coverage.
func (d *drawer) clip(cov *mask) {
	ctx := d.ctx
	st := &ctx.st
	old := st.clip
	d.drop(ctx, st)
	if !old.empty() && !cov.empty() && (old.a != nil || cov.a != nil) {
		// their product takes memory
		r := old.r.Intersect(cov.r)
		if !d.hold(ctx, 4*r.Dx()*r.Dy()) {
			d.r.warnf(tooMuchMemory)
			st.clip = &mask{}
			return
		}
		st.clip = intersect(old, cov)
		st.made = st.clip
		return
	}
	st.clip = intersect(old, cov)
}

func (d *drawer) instr(in bdf.Instr) {
	if d.again {
		if d.reused >= d.most {
			d.r.warnf("objects are drawn again too many times: the rest of them is not drawn")
			d.leftOut(in)
			return
		}
		d.reused++
	}
	if d.skipped > 0 {
		d.leftOut(in)
		return
	}
	ctx := d.ctx
	st := &ctx.st
	switch in.Op {
	case bdf.OpSave:
		if len(ctx.stack) >= maxStates {
			d.r.warnf("more than %d states are saved at a time: the rest are not kept", maxStates)
			ctx.unsaved++
			return
		}
		ctx.stack = append(ctx.stack, *st)
		st.made = nil
	case bdf.OpRestore:
		d.restore()
	case bdf.OpTransform:
		st.m = st.m.mul(matrix{f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4), f64(in, 5)})
	case bdf.OpTranslate:
		st.m = st.m.translate(f64(in, 0), f64(in, 1))
	case bdf.OpScale:
		st.m = st.m.scale(f64(in, 0), f64(in, 1))
	case bdf.OpClipPath:
		if p := d.path(int(u64(in, 0))); p != nil {
			d.clip(d.coverage(p, byte(u64(in, 1)), st.m))
		}
	case bdf.OpClipRect:
		p := &path{}
		p.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3))
		d.clip(d.coverage(p, bdf.NonZero, st.m))
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
		// a value that is no number is ignored, as Canvas ignores it
		if v := f64(in, 1) * s; !math.IsNaN(v) && !math.IsInf(v, 0) {
			st.shadowBlur = math.Min(math.Max(0, v), maxShadowBlur)
		}
		if v := f64(in, 2) * s; !math.IsNaN(v) && !math.IsInf(v, 0) {
			st.shadowDX = clampF(v)
		}
		if v := f64(in, 3) * s; !math.IsNaN(v) && !math.IsInf(v, 0) {
			st.shadowDY = clampF(v)
		}
	case bdf.OpFilter:
		if s, _ := in.Args[0].(string); s != "" && s != "none" {
			d.r.warnf("FILTER %q is not applied", s)
		}
	case bdf.OpFont:
		if i := u64(in, 0); i < uint64(len(d.obj.Fonts)) {
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
		sc, bounds := d.r.scratch(), d.bounds()
		polys := sc.polys
		for _, g := range glyphs {
			if p := d.path(int(g.Path)); p != nil {
				m := st.m.translate(float64(g.X), float64(g.Y))
				if hidden(p, 0, m, bounds) {
					continue
				}
				polys = p.flattenTo(sc, polys, m, flatTol)
				if sc.limited {
					return
				}
			}
		}
		sc.polys = polys
		d.paint(sc.rasterize(polys, rule, bounds), d.shader(&st.fill, st.m), 1)
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
		i := u64(in, 0)
		if i >= uint64(len(d.obj.Objects)) {
			return
		}
		child := d.r.object(d.obj.Objects[i])
		if child == nil {
			return
		}
		saved := *st
		depth, unsaved := len(ctx.stack), ctx.unsaved
		again := d.again
		d.again = d.enter(d.obj.Objects[i]) || again
		st.made = nil
		if in.Op == bdf.OpUseAt {
			st.m = st.m.translate(f64(in, 1), f64(in, 2))
		}
		d.run(child)
		d.again = again
		// USE is wrapped in an implicit save and restore, of the canvas it
		// draws on (which is gone if the object ended its group)
		if !ctx.done {
			d.drop(ctx, st)
			for i := depth; i < len(ctx.stack); i++ {
				d.drop(ctx, &ctx.stack[i])
			}
			ctx.stack = ctx.stack[:min(depth, len(ctx.stack))]
			ctx.unsaved = unsaved
			*st = saved
		}
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
	bounds := d.bounds()
	if hidden(p, 0, m, bounds) {
		return &mask{}
	}
	// an axis-aligned rectangle: exact coverage without scanning
	if len(p.subs) >= 1 && m[1] == 0 && m[2] == 0 {
		if x0, y0, x1, y1, ok := p.isRect(); ok {
			ax, ay := m.apply(x0, y0)
			bx, by := m.apply(x1, y1)
			return rectMask(ax, ay, bx, by, bounds)
		}
	}
	sc := d.r.scratch()
	sc.polys = p.flattenTo(sc, sc.polys, m, flatTol)
	return sc.rasterize(sc.polys, rule, bounds)
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
	cov, alpha, cut := d.r.scratch().strokeMask(p, m, d.ctx.st.line, d.bounds())
	if cut {
		d.r.warnf(strokeCut)
	}
	d.paint(cov, d.shader(&d.ctx.st.stroke, m), alpha)
}

// strokeCut warns of a stroke whose outline passed the limits.
const strokeCut = "a line has too many dashes, joins and caps: the rest of it is not drawn"

// strokeMask rasterizes the stroke of a path in user space under m within
// bounds, and returns the alpha to draw it with: strokes thinner than a
// device pixel are a pixel wide with their width as alpha.
func strokeMask(p *path, m matrix, ls lineStyle, bounds image.Rectangle) (*mask, float32) {
	cov, alpha, _ := (&scratch{}).strokeMask(p, m, ls, bounds)
	return cov, alpha
}

// strokeMask is the function of that name with the memory of a renderer;
// cut reports that the outline passed the limits of a stroke and is drawn
// without the rest.
func (sc *scratch) strokeMask(p *path, m matrix, ls lineStyle, bounds image.Rectangle) (cov *mask, alpha float32, cut bool) {
	if !m.finite() {
		return &mask{}, 0, false
	}
	scale := m.areaScale()
	if scale == 0 {
		return &mask{}, 0, false
	}
	alpha = 1
	if dw := ls.width * scale; dw < 1 {
		alpha = float32(dw)
		ls.width = 1 / scale
	}
	if hidden(p, ls.width/2*math.Max(math.Sqrt2, ls.miterReach()), m, bounds) {
		return &mask{}, alpha, false
	}
	ms := m.maxScale()
	sc.polys = p.flattenTo(sc, sc.polys, identity, flatTol/ms)
	if sc.limited {
		return &mask{}, alpha, false
	}
	s := newStroker(&ls, ms).within(m, bounds)
	s.stroke(sc.polys)
	outline := s.out
	for i := range outline {
		for j, q := range outline[i].pts {
			x, y := m.apply(q.x, q.y)
			outline[i].pts[j] = point{x, y}
		}
	}
	return sc.rasterize(outline, bdf.NonZero, bounds), alpha, s.cut
}

// shader returns the shader of a fill or stroke style drawn under m.
func (d *drawer) shader(s *style, m matrix) shader {
	if s.paint == nil {
		return premul(s.color)
	}
	if s.paint.Kind == bdf.PaintPattern {
		return newPatternShader(d.r, s, m)
	}
	return newGradientShader(s.paint, m, d.r.lut(s.paint))
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
	dx, dy := int(math.Round(st.shadowDX)), int(math.Round(st.shadowDY))
	if where := image.Rect(r.Min.X-pad+dx, r.Min.Y-pad+dy, r.Max.X+pad+dx, r.Max.Y+pad+dy); !where.Overlaps(d.bounds()) {
		// the shadow falls outside the canvas: nothing of it to blur
		if unbounded(st.blend) {
			d.ctx.target.fill(&mask{}, premul(st.shadow), 1, st.blend, st.clip)
		}
		return
	}
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
	if stroke {
		d.strokeText(glyphs, m)
		return
	}
	sc, bounds := d.r.scratch(), d.bounds()
	// synthesized bold outlines what is filled, in device space
	boldW := fakeBoldWidth(st.size * m.areaScale())
	bold := newStroker(&lineStyle{width: boldW, join: bdf.JoinRound, miter: 10}, 1).within(identity, bounds)
	fill := sc.polys
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
		thick := g.use.bold && boldW > 0
		if box, ok := p.extent(); ok && !plain {
			reach := 0.0
			if thick {
				reach = boldW
			}
			if outside(boxUnder(box, gm), reach, identity, bounds) {
				continue
			}
		}
		from := len(fill)
		fill = p.flattenTo(sc, fill, gm, flatTol)
		if sc.limited {
			return
		}
		if thick {
			bold.stroke(fill[from:])
		}
	}
	sc.polys = fill
	if bold.cut {
		d.r.warnf(strokeCut)
	}
	cov := sc.rasterize(fill, bdf.NonZero, bounds)
	if len(bold.out) > 0 {
		cov = union(cov, sc.rasterize(bold.out, bdf.NonZero, bounds))
	}
	if sc.limited {
		return
	}
	d.paint(cov, d.shader(&st.fill, m), 1)
}

// strokeText strokes glyphs laid out under m: in user space, where the
// line width is in units.
func (d *drawer) strokeText(glyphs []placed, m matrix) {
	st := &d.ctx.st
	sc, bounds := d.r.scratch(), d.bounds()
	ls := st.line
	scale := m.areaScale()
	alpha := float32(1)
	if dw := ls.width * scale; dw < 1 && scale > 0 {
		alpha, ls.width = float32(dw), 1/scale
	}
	reach := ls.width / 2 * math.Max(math.Sqrt2, ls.miterReach())
	lines := sc.polys
	for _, g := range glyphs {
		p := g.use.face.outline(g.gid)
		if p == nil {
			continue
		}
		k := st.size / g.use.face.upem
		um := matrix{1, 0, 0, 1, 0, 0}.translate(g.x, 0).scale(k, -k)
		if g.use.italic {
			um = um.mul(matrix{1, 0, fakeItalicSkew, 1, 0, 0})
		}
		if box, ok := p.extent(); ok && !plain && outside(boxUnder(box, um), reach, m, bounds) {
			continue
		}
		lines = p.flattenTo(sc, lines, um, flatTol/math.Max(m.maxScale(), 1e-9))
		if sc.limited {
			return
		}
	}
	sc.polys = lines
	outline := newStroker(&ls, m.maxScale()).within(m, bounds)
	outline.stroke(lines)
	if outline.cut {
		d.r.warnf(strokeCut)
	}
	for i := range outline.out {
		for j, q := range outline.out[i].pts {
			px, py := m.apply(q.x, q.y)
			outline.out[i].pts[j] = point{px, py}
		}
	}
	d.paint(sc.rasterize(outline.out, bdf.NonZero, bounds), d.shader(&st.stroke, m), alpha)
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
	st := &d.ctx.st
	p := &path{}
	p.rect(x, y, w, h)
	cov := d.coverage(p, bdf.NonZero, st.m)
	if cov.empty() && !unbounded(st.blend) && !plain {
		return // it shows nowhere: not worth decoding
	}
	hash := d.obj.Images[i]
	pic := d.r.picture(hash)
	si := d.r.svg(hash)
	if pic == nil && si == nil || w == 0 || h == 0 {
		return
	}
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
		if pic == nil {
			return
		}
		ky := float64(pic.h) / si.h
		sx, sy, sw, sh = sx*k, sy*ky, sw*k, sh*ky
	}
	// device → image pixels: invert (m × dest rect ← source rect)
	full := st.m.mul(matrix{w / sw, 0, 0, h / sh, x - sx*w/sw, y - sy*h/sh})
	inv, ok := full.invert()
	if !ok {
		return
	}
	lo := func(a, b float64) float64 { return math.Min(a, a+b) }
	sh2 := &imageShader{
		inv: inv, img: pic, level: pic.levelFor(inv), smooth: st.smooth,
		sx0: math.Max(0, lo(sx, sw)), sy0: math.Max(0, lo(sy, sh)),
		sx1: math.Min(float64(pic.w), lo(sx, sw)+math.Abs(sw)), sy1: math.Min(float64(pic.h), lo(sy, sh)+math.Abs(sh)),
	}
	d.paint(cov, sh2, 1)
}

// tooMuchMemory warns of groups, soft masks and clips past maxLive.
const tooMuchMemory = "the groups and clips of a page take too much memory: what is in the rest of them is not drawn"

// layer reports whether there is room for the canvas of one more group or
// soft mask, of w × h pixels, and counts it. One without room is left out
// with what is in it.
func (d *drawer) layer(ctx *context, w, h int) bool {
	switch {
	case len(d.groups)+len(d.masks) >= maxLayers:
		d.r.warnf("more than %d groups and soft masks are open: what is in the rest of them is not drawn", maxLayers)
	case !d.hold(ctx, 16*w*h):
		d.r.warnf(tooMuchMemory)
	default:
		return true
	}
	d.skipped++
	return false
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
	gc := &context{st: gs}
	if !d.layer(gc, r.Dx(), r.Dy()) {
		return
	}
	gc.target = newSurface(r.Dx(), r.Dy())
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
		d.closed(d.masks[len(d.masks)-1].ctx)
		d.masks = d.masks[:len(d.masks)-1]
	}
	d.closed(g.ctx)
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
	mc := &context{}
	if !d.layer(mc, w, h) {
		return
	}
	mc.target = newSurface(w, h)
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
	d.closed(m.ctx)
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

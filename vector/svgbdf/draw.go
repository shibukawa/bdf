package svgbdf

import (
	"bytes"
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// paint is a fill or stroke style: a color, or a gradient or pattern of an
// object.
type paint struct {
	color bdf.Color
	p     *bdf.Paint
	obj   *bdf.ObjectPart // whose images a pattern refers to
}

// state is the drawing state of a context (what save and restore keep).
type state struct {
	m    matrix
	clip int // the clip path (see drawer.clip), 0 for none

	fill, stroke paint
	width, miter float64
	cap, join    byte
	dash         []float32
	dashOffset   float64
	alpha        float64
	blend        byte

	shadow             bdf.Color
	shadowBlur         float64 // in units of the page, like the offsets
	shadowDX, shadowDY float64

	font     *bdf.Font
	size     float64
	align    byte
	baseline byte
	dir      byte
	spacing  float64
	smooth   bool
}

// Limits of what the drawing of a page takes. A document decides how often
// its objects are drawn, how deep its groups are and how many states it
// keeps, so a small one could otherwise ask for any time and memory.
const (
	// maxReusedInstructions bounds the instructions drawn from objects
	// that were drawn before in the same SVG. What is drawn once is
	// bounded by the size of the document; this bounds what a document
	// adds by drawing it again, each time written out in full.
	maxReusedInstructions = 1 << 21
	// maxLayers bounds the groups and soft masks that are open at a time.
	maxLayers = 64
	// maxStates bounds the states a context keeps for RESTORE.
	maxStates = 1024
)

// initialState is the Canvas 2D initial state with a transform and clip.
func initialState(m matrix, clip int) state {
	return state{
		m: m, clip: clip,
		fill: paint{color: 0x000000ff}, stroke: paint{color: 0x000000ff},
		width: 1, miter: 10, alpha: 1, size: 10, smooth: true,
	}
}

func (s *state) shadowActive() bool {
	return uint32(s.shadow)&0xff != 0 && (s.shadowBlur > 0 || s.shadowDX != 0 || s.shadowDY != 0)
}

// shape is a fill that waits for the stroke of the same outline, which
// makes one element of the two.
type shape struct {
	geom string // the element's name and geometry
	rule string
	tf   string
	clip int
	fill string
}

// context is what a canvas is to the other renderers: the page, or the
// content of a group or soft mask. Its elements have the transform they
// are drawn under each, so that the groups around them need none.
type context struct {
	st    state
	stack []state
	// unsaved counts the SAVEs that kept no state (past maxStates, or in
	// what is left out): as many RESTOREs restore none
	unsaved int
	buf     *bytes.Buffer
	// clip is the clip path of the group that is open for the elements
	// being written (0: none is), and attrs the attributes of the group
	// open inside it for elements that share them ("": none is)
	clip  int
	attrs string
	pend  *shape
	// blended says that an element of the context has a blend mode, which
	// must not reach what lies under the context
	blended bool
	// erase holds elements that erase what buf holds (destination-out),
	// until the next thing is drawn
	erase *bytes.Buffer
	// done says that the group or soft mask of the context has ended
	done bool
}

type group struct {
	parent *context
	ctx    *context
	alpha  float64
	blend  byte
	mask   string // the id of its soft mask
}

type softMask struct {
	parent   *context
	ctx      *context
	kind     byte
	backdrop bdf.Color
	transfer []byte
	group    *group
}

type glyphKey struct {
	obj   *bdf.ObjectPart // of an inline path
	hash  bdf.Hash        // of a path collection
	index uint32
}

// drawer runs the instructions of objects into one SVG.
type drawer struct {
	r    *Renderer
	view int
	// the part of the page the SVG shows, in units: its view box
	x, y, w, h float64

	ctx    *context
	obj    *bdf.ObjectPart
	groups []*group
	masks  []*softMask
	depth  int

	// seen holds the objects drawn in this SVG, again says that the
	// object being drawn (or one that draws it) was drawn before, and
	// reused counts the instructions of such objects
	seen   map[bdf.Hash]struct{}
	again  bool
	reused int
	// skipped counts the groups and soft masks that are open but not
	// drawn: those past the limits, and those of objects past theirs
	skipped int

	defs, css, links bytes.Buffer
	ids              int
	clips            map[string]int
	paints           map[*bdf.Paint]string
	glyphs           map[glyphKey]string
	filters          map[string]string
	fonts            map[bdf.Hash]string
	symbols          map[bdf.Hash]string
}

func (d *drawer) id(prefix string) string {
	d.ids++
	return prefix + strconv.Itoa(d.ids)
}

// region is the attributes of a mask or filter that reaches all the SVG
// shows.
func (d *drawer) region() string {
	return ` x="` + num(d.x, pageDigits) + `" y="` + num(d.y, pageDigits) + `" width="` + num(d.w, pageDigits) + `" height="` + num(d.h, pageDigits) + `"`
}

// enter reports whether an object was drawn before in this SVG, and notes
// that it is now.
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

// --- writing elements ---

// flush writes the fill that waited for a stroke.
func (d *drawer) flush(c *context) {
	if p := c.pend; p != nil {
		c.pend = nil
		d.write(c, p.clip, "<"+p.geom+p.rule+p.tf+p.fill+"/>")
	}
}

// settle applies the elements that erase to what the context holds: they
// become a mask of what was drawn before them.
func (d *drawer) settle(c *context) {
	e := c.erase
	if e == nil {
		return
	}
	c.erase = nil
	if c.buf.Len() == 0 {
		return
	}
	id := d.id("m")
	box, _ := rectData(d.x, d.y, d.w, d.h, pageDigits)
	d.defs.WriteString(`<mask id="` + id + `" maskUnits="userSpaceOnUse"` + d.region() + `><path d="` + box + `" fill="#fff"/>`)
	d.defs.Write(e.Bytes())
	d.defs.WriteString("</mask>\n")
	d.wrap(c, `<g mask="url(#`+id+`)">`)
}

// wrap puts what a context holds into a group.
func (d *drawer) wrap(c *context, open string) {
	b := bytes.NewBuffer(make([]byte, 0, c.buf.Len()+len(open)+8))
	b.WriteString(open)
	b.WriteByte('\n')
	b.Write(c.buf.Bytes())
	b.WriteString("</g>\n")
	c.buf = b
}

// closeAttrs closes the group of shared attributes.
func (d *drawer) closeAttrs(c *context) {
	if c.attrs != "" {
		c.attrs = ""
		c.buf.WriteString("</g>\n")
	}
}

// closeAll writes what waits and closes the groups open for the elements
// being written.
func (d *drawer) closeAll(c *context) {
	d.flush(c)
	d.closeAttrs(c)
	if c.clip != 0 {
		c.clip = 0
		c.buf.WriteString("</g>\n")
	}
}

func clipAttr(clip int) string {
	if clip == 0 {
		return ""
	}
	return ` clip-path="url(#c` + strconv.Itoa(clip) + `)"`
}

// within opens the group of a clip path for the elements that follow. A
// clip path is in the coordinates of what refers to it, which for an
// element with a transform are not those of the page: the group has none.
func (d *drawer) within(c *context, clip int) {
	d.settle(c)
	if c.clip == clip {
		return
	}
	d.closeAttrs(c)
	if c.clip != 0 {
		c.buf.WriteString("</g>\n")
	}
	if c.clip = clip; clip != 0 {
		c.buf.WriteString("<g" + clipAttr(clip) + ">\n")
	}
}

// write writes an element drawn over what is there.
func (d *drawer) write(c *context, clip int, el string) {
	d.within(c, clip)
	d.closeAttrs(c)
	c.buf.WriteString(el)
	c.buf.WriteByte('\n')
}

// writeIn writes an element into a group of the attributes it shares with
// its neighbours (the runs of a line of text, the glyphs of a run).
func (d *drawer) writeIn(c *context, clip int, attrs, el string) {
	d.flush(c)
	d.within(c, clip)
	if c.attrs != attrs {
		d.closeAttrs(c)
		// elements with nothing to share need no group
		if c.attrs = attrs; attrs != "" {
			c.buf.WriteString("<g" + attrs + ">\n")
		}
	}
	c.buf.WriteString(el)
	c.buf.WriteByte('\n')
}

// How a composite operation is drawn.
const (
	modeOver  = iota // drawn over what is there, with a blend mode or none
	modeErase        // destination-out: a mask of what is there
	modeKeep         // destination-in: a mask of what is there
	modeUnder        // destination-over: drawn first
)

// mode tells how a composite operation is drawn, and the blend mode of
// CSS for it ("" for none). SVG has no element for the operations that
// keep or leave out what is drawn by what is there: they are drawn as
// source-over, with a warning.
func (d *drawer) mode(blend byte) (int, string) {
	switch {
	case blend == bdf.BlendSourceOver:
		return modeOver, ""
	case blend <= bdf.BlendLuminosity:
		return modeOver, bdf.BlendNames[blend]
	case blend == bdf.BlendLighter:
		return modeOver, "plus-lighter"
	case blend == bdf.BlendDestinationOut:
		return modeErase, ""
	case blend == bdf.BlendDestinationIn:
		return modeKeep, ""
	case blend == bdf.BlendDestinationOver:
		return modeUnder, ""
	}
	if int(blend) < len(bdf.BlendNames) {
		d.r.warnf("the composite operation %s is drawn as source-over", bdf.BlendNames[blend])
	}
	return modeOver, ""
}

// compose draws an element as the state says: over what is there, with
// its blend mode and shadow, or as the composite operation makes of it
// and what is there. black says that the element is drawn in black (or,
// for destination-in, in any color) with the alpha it has.
func (d *drawer) compose(c *context, mode int, blend, el string, black bool) {
	st := &c.st
	switch mode {
	case modeErase:
		d.closeAll(c)
		if !black {
			el = `<g filter="url(#` + d.blackFilter() + `)">` + el + "</g>"
		}
		if c.erase == nil {
			c.erase = &bytes.Buffer{}
		}
		c.erase.WriteString("<g" + clipAttr(st.clip) + ">" + el + "</g>")
		return
	case modeKeep:
		d.closeAll(c)
		d.settle(c)
		if c.buf.Len() == 0 {
			return
		}
		id := d.id("m")
		d.defs.WriteString(`<mask id="` + id + `" maskUnits="userSpaceOnUse"` + d.region() + ` style="mask-type:alpha"><g` + clipAttr(st.clip) + ">" + el + "</g></mask>\n")
		d.wrap(c, `<g mask="url(#`+id+`)">`)
		return
	case modeUnder:
		d.closeAll(c)
		d.settle(c)
		b := bytes.NewBufferString("<g" + clipAttr(st.clip) + ">" + el + "</g>\n")
		b.Write(c.buf.Bytes())
		c.buf = b
		return
	}
	if st.shadowActive() {
		el = `<g filter="url(#` + d.shadowFilter(st) + `)">` + el + "</g>"
	}
	if blend == "" {
		d.flush(c)
		d.write(c, st.clip, el)
		return
	}
	// what blends is a child of the context: a group around it would be
	// what it blends with
	d.closeAll(c)
	d.settle(c)
	c.blended = true
	c.buf.WriteString("<g" + clipAttr(st.clip) + ` style="mix-blend-mode:` + blend + `">` + el + "</g>\n")
}

// blackFilter is the filter that keeps the alpha of what it is applied to
// and makes it black: what erases in a mask of luminance.
func (d *drawer) blackFilter() string {
	return d.filter("black", `<feColorMatrix values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0"/>`)
}

// filter defines a filter that reaches all the SVG shows, once.
func (d *drawer) filter(key, primitives string) string {
	if id, ok := d.filters[key]; ok {
		return id
	}
	id := d.id("f")
	d.filters[key] = id
	d.defs.WriteString(`<filter id="` + id + `" filterUnits="userSpaceOnUse"` + d.region() + ` color-interpolation-filters="sRGB">` + primitives + "</filter>\n")
	return id
}

// shadowFilter is the filter of the shadow of the state. It is applied to
// a group without a transform: the offset of a shadow has the directions
// of the page.
func (d *drawer) shadowFilter(st *state) string {
	c := st.shadow
	prim := `<feDropShadow dx="` + num(st.shadowDX, pageDigits) + `" dy="` + num(st.shadowDY, pageDigits) +
		`" stdDeviation="` + num(st.shadowBlur/2, pageDigits) + `" flood-color="` + hexColor(c) + `"`
	if a := uint32(c) & 0xff; a != 0xff {
		prim += ` flood-opacity="` + factor(float64(a)/255) + `"`
	}
	return d.filter(prim, prim+"/>")
}

func hexColor(c bdf.Color) string {
	const digits = "0123456789abcdef"
	v := uint32(c) >> 8
	if v>>20&0xf == v>>16&0xf && v>>12&0xf == v>>8&0xf && v>>4&0xf == v&0xf {
		return string([]byte{'#', digits[v>>16&0xf], digits[v>>8&0xf], digits[v&0xf]})
	}
	b := []byte("#000000")
	for i := 0; i < 6; i++ {
		b[6-i] = digits[v>>(4*i)&0xf]
	}
	return string(b)
}

// paintAttrs is the attributes of a fill or a stroke (name) in a style.
// Black is what SVG fills with and none what it strokes with, so those are
// left out.
func (d *drawer) paintAttrs(name string, p *paint, alpha float64, mode int) string {
	value, opacity := "", alpha
	switch {
	case mode == modeErase || mode == modeKeep:
		// only its alpha counts
		if p.p == nil {
			opacity *= float64(uint32(p.color)&0xff) / 255
		}
		value = "#000"
	case p.p != nil:
		value = d.paintServer(p)
	default:
		value = hexColor(p.color)
		opacity *= float64(uint32(p.color)&0xff) / 255
	}
	s := ""
	if !(name == "fill" && value == "#000") {
		s = " " + name + `="` + value + `"`
	}
	if opacity < 1 {
		s += " " + name + `-opacity="` + factor(math.Max(opacity, 0)) + `"`
	}
	return s
}

func (d *drawer) fillAttrs(st *state, mode int) string {
	return d.paintAttrs("fill", &st.fill, st.alpha, mode)
}

var (
	capNames  = [...]string{"", "round", "square"}
	joinNames = [...]string{"", "round", "bevel"}
)

// strokeAttrs is the attributes of a stroke in the state, for an element
// that is not filled or is filled by other attributes.
func (d *drawer) strokeAttrs(st *state, mode int) string {
	s := d.paintAttrs("stroke", &st.stroke, st.alpha, mode)
	if st.width != 1 {
		s += ` stroke-width="` + factor(st.width) + `"`
	}
	if int(st.cap) < len(capNames) && capNames[st.cap] != "" {
		s += ` stroke-linecap="` + capNames[st.cap] + `"`
	}
	if int(st.join) < len(joinNames) && joinNames[st.join] != "" {
		s += ` stroke-linejoin="` + joinNames[st.join] + `"`
	} else if st.miter != 4 {
		s += ` stroke-miterlimit="` + factor(st.miter) + `"`
	}
	if len(st.dash) > 0 {
		s += ` stroke-dasharray="`
		for i, v := range st.dash {
			if i > 0 {
				s += " "
			}
			s += factor(float64(v))
		}
		s += `"`
		if st.dashOffset != 0 {
			s += ` stroke-dashoffset="` + factor(st.dashOffset) + `"`
		}
	}
	return s
}

// shape fills or strokes an outline (geom: the name and the geometry of
// its element) under m. A stroke of the outline that was filled just
// before is the same element.
func (d *drawer) shape(geom string, rule byte, m matrix, stroke bool) {
	c := d.ctx
	st := &c.st
	if !m.finite() {
		return
	}
	mode, blend := d.mode(st.blend)
	tf := m.transform()
	if stroke {
		if p := c.pend; p != nil && mode == modeOver && blend == "" && !st.shadowActive() && p.geom == geom && p.tf == tf && p.clip == st.clip {
			c.pend = nil
			d.write(c, st.clip, "<"+geom+p.rule+tf+p.fill+d.strokeAttrs(st, mode)+"/>")
			return
		}
		d.compose(c, mode, blend, "<"+geom+tf+` fill="none"`+d.strokeAttrs(st, mode)+"/>", true)
		return
	}
	fr := ""
	if rule == bdf.EvenOdd {
		fr = ` fill-rule="evenodd"`
	}
	fill := d.fillAttrs(st, mode)
	if mode == modeOver && blend == "" && !st.shadowActive() {
		d.flush(c)
		c.pend = &shape{geom: geom, rule: fr, tf: tf, clip: st.clip, fill: fill}
		return
	}
	d.compose(c, mode, blend, "<"+geom+fr+tf+fill+"/>", true)
}

// path returns a path of the object being drawn, or nil.
func (d *drawer) path(i int) *bdf.Path {
	if i < 0 || i >= len(d.obj.Paths) {
		d.r.warnf("bad path reference")
		return nil
	}
	return d.r.path(d.obj, i)
}

// outline fills or strokes a path of the object.
func (d *drawer) outline(i int, rule byte, stroke bool) {
	p := d.path(i)
	if p == nil {
		return
	}
	m := d.ctx.st.m
	if data, ok := pathData(p, m.digits()); ok {
		d.shape(`path d="`+data+`"`, rule, m, stroke)
	}
}

func (d *drawer) rect(x, y, w, h float64, stroke bool) {
	m := d.ctx.st.m
	if data, ok := rectData(x, y, w, h, m.digits()); ok {
		d.shape(`path d="`+data+`"`, bdf.NonZero, m, stroke)
	}
}

// glyph returns the id of a path that is drawn in many places, defining it
// at its first.
func (d *drawer) glyph(i int, m matrix) (string, bool) {
	if i < 0 || i >= len(d.obj.Paths) {
		d.r.warnf("bad path reference")
		return "", false
	}
	e := d.obj.Paths[i]
	k := glyphKey{hash: e.Hash, index: e.Index}
	if e.Inline != nil {
		k = glyphKey{obj: d.obj, index: uint32(i)}
	}
	if id, ok := d.glyphs[k]; ok {
		return id, id != ""
	}
	id := ""
	if p := d.r.path(d.obj, i); p != nil {
		// a digit more than its first place needs: others may be larger
		if data, ok := pathData(p, min(m.digits()+1, 9)); ok {
			id = d.id("a")
			d.defs.WriteString(`<path id="` + id + `" d="` + data + `"/>` + "\n")
		}
	}
	d.glyphs[k] = id
	return id, id != ""
}

// placed fills paths at places: the outlines of glyphs, each defined once
// and referred to.
func (d *drawer) placed(rule byte, glyphs []bdf.Glyph) {
	c := d.ctx
	st := &c.st
	if !st.m.finite() {
		return
	}
	mode, blend := d.mode(st.blend)
	attrs := st.m.transform() + d.fillAttrs(st, mode)
	if rule == bdf.EvenOdd {
		attrs += ` fill-rule="evenodd"`
	}
	digits := st.m.digits()
	var b strings.Builder
	for _, g := range glyphs {
		id, ok := d.glyph(int(g.Path), st.m)
		if !ok || math.IsNaN(float64(g.X+g.Y)) || math.IsInf(float64(g.X+g.Y), 0) {
			continue
		}
		el := `<use href="#` + id + `"`
		if g.X != 0 {
			el += ` x="` + num(float64(g.X), digits) + `"`
		}
		if g.Y != 0 {
			el += ` y="` + num(float64(g.Y), digits) + `"`
		}
		el += "/>"
		if mode == modeOver && blend == "" && !st.shadowActive() {
			d.writeIn(c, st.clip, attrs, el)
			continue
		}
		b.WriteString(el)
	}
	if b.Len() > 0 {
		// one drawing, as the shadow and the blend mode of one are
		d.compose(c, mode, blend, "<g"+attrs+">"+b.String()+"</g>", true)
	}
}

// clip cuts the clip of the state to an outline under the transform of
// the state. A clip path refers to the one it narrows.
func (d *drawer) clip(data string, rule byte) {
	st := &d.ctx.st
	el := ""
	if data != "" && st.m.finite() {
		el = `<path d="` + data + `"` + st.m.transform()
		if rule == bdf.EvenOdd {
			el += ` clip-rule="evenodd"`
		}
		el += "/>"
	}
	key := strconv.Itoa(st.clip) + el
	id, ok := d.clips[key]
	if !ok {
		id = len(d.clips) + 1
		d.clips[key] = id
		d.defs.WriteString(`<clipPath id="c` + strconv.Itoa(id) + `"` + clipAttr(st.clip) + ">" + el + "</clipPath>\n")
	}
	st.clip = id
}

// --- text ---

var (
	baselineNames = [...]string{"", "text-before-edge", "central", "text-after-edge", "hanging", "ideographic"}
	styleNames    = [...]string{"", "italic", "oblique"}
)

// text draws a FILL_TEXT or STROKE_TEXT as a text element. The advance is
// the length the text is fitted to, as the viewer scales a run whose font
// measures otherwise.
func (d *drawer) text(s string, x, y, advance float64, stroke bool) {
	c := d.ctx
	st := &c.st
	if !(st.size > 0) || math.IsInf(st.size, 0) || s == "" || !st.m.finite() || math.IsNaN(x+y) || math.IsInf(x+y, 0) {
		return
	}
	mode, blend := d.mode(st.blend)
	f := st.font
	if f == nil {
		f = &bdf.Font{Kind: bdf.FontSystem, Family: "sans-serif", Weight: 400}
	}
	attrs := st.m.transform() + ` font-family="` + escape(d.family(f)) + `" font-size="` + factor(st.size) + `"`
	if f.Weight != 0 && f.Weight != 400 {
		attrs += ` font-weight="` + strconv.Itoa(int(f.Weight)) + `"`
	}
	if int(f.Style) < len(styleNames) && styleNames[f.Style] != "" {
		attrs += ` font-style="` + styleNames[f.Style] + `"`
	}
	if stroke {
		attrs += ` fill="none"` + d.strokeAttrs(st, mode)
	} else {
		attrs += d.fillAttrs(st, mode)
	}
	rtl := st.dir == bdf.DirRTL
	if rtl {
		attrs += ` direction="rtl"`
	}
	// start and end are the sides of the direction
	switch st.align {
	case bdf.AlignCenter:
		attrs += ` text-anchor="middle"`
	case bdf.AlignRight:
		if !rtl {
			attrs += ` text-anchor="end"`
		}
	case bdf.AlignLeft:
		if rtl {
			attrs += ` text-anchor="end"`
		}
	case bdf.AlignEnd:
		attrs += ` text-anchor="end"`
	}
	if int(st.baseline) < len(baselineNames) && baselineNames[st.baseline] != "" {
		attrs += ` dominant-baseline="` + baselineNames[st.baseline] + `"`
	}
	if st.spacing != 0 && !math.IsNaN(st.spacing) && !math.IsInf(st.spacing, 0) {
		attrs += ` letter-spacing="` + factor(st.spacing) + `"`
	}
	digits := st.m.digits()
	pos := ""
	if x != 0 {
		pos += ` x="` + num(x, digits) + `"`
	}
	if y != 0 {
		pos += ` y="` + num(y, digits) + `"`
	}
	if advance > 0 && !math.IsInf(advance, 0) {
		pos += ` textLength="` + num(advance, digits) + `" lengthAdjust="spacingAndGlyphs"`
	}
	text := escape(s)
	// Canvas draws every space. The attribute is read on the text element
	// itself: Chrome takes it from no element around it.
	if strings.HasPrefix(text, " ") || strings.HasSuffix(text, " ") || strings.Contains(text, "  ") {
		pos += ` xml:space="preserve"`
	}
	body := ">" + text + "</text>"
	if mode == modeOver && blend == "" && !st.shadowActive() {
		d.writeIn(c, st.clip, attrs, "<text"+pos+body)
		return
	}
	d.compose(c, mode, blend, "<text"+attrs+pos+body, true)
}

// --- images ---

// image draws an image, or the part src (x, y, w, h in image pixels) of
// it, into the rectangle (x, y, w, h).
func (d *drawer) image(i int, src *[4]float64, x, y, w, h float64) {
	if i < 0 || i >= len(d.obj.Images) {
		d.r.warnf("bad image reference")
		return
	}
	c := d.ctx
	st := &c.st
	if w == 0 || h == 0 || !st.m.finite() || math.IsNaN(x+y+w+h) || math.IsInf(x+y+w+h, 0) {
		return
	}
	info := d.r.image(d.obj.Images[i])
	if info == nil {
		return
	}
	// a rectangle of a negative size mirrors the image
	m := st.m
	if w < 0 || h < 0 {
		sx, sy := 1.0, 1.0
		if w < 0 {
			sx, x, w = -1, -x, -w
		}
		if h < 0 {
			sy, y, h = -1, -y, -h
		}
		m = m.scale(sx, sy)
	}
	attrs := ""
	if st.alpha < 1 {
		attrs += ` opacity="` + factor(st.alpha) + `"`
	}
	if !st.smooth {
		attrs += ` style="image-rendering:pixelated"`
	}
	digits := m.digits()
	var el string
	if src == nil {
		el = d.imageElement(info, x, y, w, h, digits, m.transform()+attrs)
	} else {
		sx, sy, sw, sh := src[0], src[1], src[2], src[3]
		if sw < 0 {
			sx, sw = sx+sw, -sw
		}
		if sh < 0 {
			sy, sh = sy+sh, -sh
		}
		if !(sw > 0) || !(sh > 0) || math.IsNaN(sx+sy) || math.IsInf(sx+sy+sw+sh, 0) {
			return
		}
		if info.w == 0 {
			d.r.warnf("a part of an image whose size cannot be told is not drawn")
			return
		}
		// the part is in the pixels of the image: a view box of them
		el = "<g" + m.transform() + attrs + "><svg"
		if x != 0 {
			el += ` x="` + num(x, digits) + `"`
		}
		if y != 0 {
			el += ` y="` + num(y, digits) + `"`
		}
		el += ` width="` + num(w, digits) + `" height="` + num(h, digits) +
			`" viewBox="` + num(sx, 3) + " " + num(sy, 3) + " " + num(sw, 3) + " " + num(sh, 3) + `" preserveAspectRatio="none">` +
			d.imageElement(info, 0, 0, info.w, info.h, 3, "") + "</svg></g>"
	}
	mode, blend := d.mode(st.blend)
	d.compose(c, mode, blend, el, false)
}

// imageElement is the element of an image stretched into a rectangle. An
// image that is in the SVG itself is there once, as a symbol.
func (d *drawer) imageElement(info *imageInfo, x, y, w, h float64, digits int, attrs string) string {
	box := ""
	if x != 0 {
		box += ` x="` + num(x, digits) + `"`
	}
	if y != 0 {
		box += ` y="` + num(y, digits) + `"`
	}
	box += ` width="` + num(w, digits) + `" height="` + num(h, digits) + `"`
	if !d.r.opts.Inline || info.w == 0 {
		return `<image href="` + escape(info.url) + `"` + box + ` preserveAspectRatio="none"` + attrs + "/>"
	}
	id, ok := d.symbols[info.hash]
	if !ok {
		id = d.id("i")
		d.symbols[info.hash] = id
		size := ` width="` + num(info.w, 3) + `" height="` + num(info.h, 3) + `"`
		d.defs.WriteString(`<symbol id="` + id + `" viewBox="0 0 ` + num(info.w, 3) + " " + num(info.h, 3) + `" preserveAspectRatio="none"><image href="` +
			escape(info.url) + `"` + size + ` preserveAspectRatio="none"/></symbol>` + "\n")
	}
	return `<use href="#` + id + `"` + box + attrs + "/>"
}

// --- groups and soft masks ---

// layer reports whether one more group or soft mask may be open. One past
// the limit is left out with what is in it.
func (d *drawer) layer() bool {
	if len(d.groups)+len(d.masks) >= maxLayers {
		d.r.warnf("more than %d groups and soft masks are open: what is in the rest of them is not drawn", maxLayers)
		d.skipped++
		return false
	}
	return true
}

func (d *drawer) groupBegin(alpha float64, blend byte) {
	if !d.layer() {
		return
	}
	c := d.ctx
	d.flush(c)
	st := &c.st
	gs := initialState(st.m, 0)
	gs.fill, gs.stroke, gs.width, gs.font, gs.size = st.fill, st.stroke, st.width, st.font, st.size
	g := &group{parent: c, ctx: &context{st: gs, buf: &bytes.Buffer{}}, alpha: alpha, blend: blend}
	d.groups = append(d.groups, g)
	d.ctx = g.ctx
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
		d.masks[len(d.masks)-1].ctx.done = true
		d.masks = d.masks[:len(d.masks)-1]
	}
	d.closeAll(g.ctx)
	d.settle(g.ctx)
	g.ctx.done = true
	c := g.parent
	d.ctx = c
	if g.ctx.buf.Len() == 0 || !(g.alpha > 0) {
		return
	}
	attrs := ""
	if g.alpha < 1 {
		attrs += ` opacity="` + factor(g.alpha) + `"`
	}
	if g.mask != "" {
		attrs += ` mask="url(#` + g.mask + `)"`
	}
	mode, blend := d.mode(g.blend)
	if mode == modeOver && blend == "" {
		// a group that changes nothing is its content, unless that
		// blends: with the group, not with what is under it
		d.flush(c)
		open, end := "", ""
		if attrs != "" || g.ctx.blended {
			if g.ctx.blended {
				attrs += ` style="isolation:isolate"`
			}
			open, end = "<g"+attrs+">\n", "</g>\n"
		}
		d.within(c, c.st.clip)
		d.closeAttrs(c)
		c.buf.WriteString(open)
		c.buf.Write(g.ctx.buf.Bytes())
		c.buf.WriteString(end)
		return
	}
	// a group is drawn without the shadow of the state
	shadow := c.st.shadow
	c.st.shadow = 0
	d.compose(c, mode, blend, "<g"+attrs+">\n"+g.ctx.buf.String()+"</g>", false)
	c.st.shadow = shadow
}

func (d *drawer) maskBegin(kind byte, backdrop bdf.Color, transfer []byte) {
	if !d.layer() {
		return
	}
	var g *group
	if n := len(d.groups); n > 0 {
		g = d.groups[n-1]
	}
	c := d.ctx
	d.flush(c)
	if len(transfer) != 256 {
		transfer = nil
	}
	mc := &context{st: initialState(c.st.m, 0), buf: &bytes.Buffer{}}
	d.masks = append(d.masks, &softMask{parent: c, ctx: mc, kind: kind, backdrop: backdrop, transfer: transfer, group: g})
	d.ctx = mc
}

func (d *drawer) maskEnd() {
	n := len(d.masks)
	if n == 0 {
		return
	}
	m := d.masks[n-1]
	d.masks = d.masks[:n-1]
	d.closeAll(m.ctx)
	d.settle(m.ctx)
	m.ctx.done = true
	d.ctx = m.parent
	g := m.group
	if g == nil || len(d.groups) == 0 || d.groups[len(d.groups)-1] != g {
		return
	}
	lum := m.kind == bdf.MaskLuminosity
	content := m.ctx.buf.String()
	if lum {
		box, _ := rectData(d.x, d.y, d.w, d.h, pageDigits)
		content = `<path d="` + box + `" fill="` + hexColor(m.backdrop) + `"/>` + "\n" + content
	}
	if m.transfer != nil {
		// the transfer function maps the value of the mask, which for a
		// mask of luminance is made a gray of the weights of BDF first
		var t strings.Builder
		for i, v := range m.transfer {
			if i > 0 {
				t.WriteByte(' ')
			}
			t.WriteString(num(float64(v)/255, 4))
		}
		prim := `<feComponentTransfer><feFuncA type="table" tableValues="` + t.String() + `"/></feComponentTransfer>`
		if lum {
			table := ` type="table" tableValues="` + t.String() + `"/>`
			prim = `<feColorMatrix values=".3 .59 .11 0 0 .3 .59 .11 0 0 .3 .59 .11 0 0 0 0 0 1 0"/>` +
				"<feComponentTransfer><feFuncR" + table + "<feFuncG" + table + "<feFuncB" + table + "</feComponentTransfer>"
		}
		id := d.id("f")
		d.defs.WriteString(`<filter id="` + id + `" filterUnits="userSpaceOnUse"` + d.region() + ` color-interpolation-filters="sRGB">` + prim + "</filter>\n")
		content = `<g filter="url(#` + id + `)">` + "\n" + content + "</g>\n"
	}
	id := d.id("m")
	kind := ""
	if !lum {
		kind = ` style="mask-type:alpha"`
	}
	d.defs.WriteString(`<mask id="` + id + `" maskUnits="userSpaceOnUse"` + d.region() + kind + ">\n" + content + "</mask>\n")
	g.mask = id
}

// --- instructions ---

// drawTop draws a top-level object (a page layer or a sheet tile) from the
// initial state with a transform and a clip.
func (d *drawer) drawTop(h bdf.Hash, m matrix, clip int) {
	o := d.r.object(h)
	if o == nil {
		return
	}
	top := d.ctx
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
	d.skipped = 0
	d.flush(top)
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
		ctx.st = ctx.stack[n-1]
		ctx.stack = ctx.stack[:n-1]
	}
}

func (d *drawer) paint(i int) paint {
	if i < 0 || i >= len(d.obj.Paints) {
		d.r.warnf("bad paint reference")
		return paint{}
	}
	return paint{p: &d.obj.Paints[i], obj: d.obj}
}

func (d *drawer) instr(in bdf.Instr) {
	if d.again {
		if d.reused >= maxReusedInstructions {
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
			data, _ := pathData(p, st.m.digits())
			d.clip(data, byte(u64(in, 1)))
		}
	case bdf.OpClipRect:
		data, _ := rectData(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), st.m.digits())
		d.clip(data, bdf.NonZero)
	case bdf.OpFillColor:
		st.fill = paint{color: bdf.Color(u64(in, 0))}
	case bdf.OpStrokeColor:
		st.stroke = paint{color: bdf.Color(u64(in, 0))}
	case bdf.OpFillPaint:
		st.fill = d.paint(int(u64(in, 0)))
	case bdf.OpStrokePaint:
		st.stroke = d.paint(int(u64(in, 0)))
	case bdf.OpLine:
		if w := f64(in, 0); w > 0 && !math.IsInf(w, 0) {
			st.width = w
		}
		st.cap, st.join = byte(u64(in, 1)), byte(u64(in, 2))
		if m := f64(in, 3); m > 0 && !math.IsInf(m, 0) {
			st.miter = m
		}
	case bdf.OpDash:
		segs, _ := in.Args[0].([]float32)
		off, _ := in.Args[1].(float32)
		// a list with a value that is no length is ignored, as Canvas
		// ignores it
		sum, valid := 0.0, true
		for _, v := range segs {
			valid = valid && v >= 0 && !math.IsInf(float64(v), 0)
			sum += float64(v)
		}
		if valid {
			st.dash = nil
			if sum > 0 {
				st.dash = segs
			}
		}
		if o := float64(off); !math.IsNaN(o) && !math.IsInf(o, 0) {
			st.dashOffset = o
		}
	case bdf.OpAlpha:
		if a := f64(in, 0); a >= 0 && a <= 1 {
			st.alpha = a
		}
	case bdf.OpBlend:
		if b := u64(in, 0); b < uint64(len(bdf.BlendNames)) {
			st.blend = byte(b)
		}
	case bdf.OpShadow:
		s := st.m.areaScale()
		if s == 0 || math.IsNaN(s) || math.IsInf(s, 0) {
			s = 1
		}
		st.shadow = bdf.Color(u64(in, 0))
		// a value that is no number is ignored, as Canvas ignores it
		if v := f64(in, 1) * s; !math.IsNaN(v) && !math.IsInf(v, 0) {
			st.shadowBlur = math.Max(0, v)
		}
		if v := f64(in, 2) * s; !math.IsNaN(v) && !math.IsInf(v, 0) {
			st.shadowDX = v
		}
		if v := f64(in, 3) * s; !math.IsNaN(v) && !math.IsInf(v, 0) {
			st.shadowDY = v
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
		d.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), false)
	case bdf.OpStrokeRect:
		d.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), true)
	case bdf.OpClearRect:
		// what erases, whatever the state draws with
		saved := *st
		st.fill, st.alpha, st.blend, st.shadow = paint{color: 0x000000ff}, 1, bdf.BlendDestinationOut, 0
		d.rect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), false)
		*st = saved
	case bdf.OpFillPath:
		d.outline(int(u64(in, 0)), byte(u64(in, 1)), false)
	case bdf.OpStrokePath:
		d.outline(int(u64(in, 0)), bdf.NonZero, true)
	case bdf.OpFillPathAt:
		d.placed(byte(u64(in, 1)), []bdf.Glyph{{Path: bdf.PathRef(u64(in, 0)), X: float32(f64(in, 2)), Y: float32(f64(in, 3))}})
	case bdf.OpFillPathRun:
		glyphs, _ := in.Args[1].([]bdf.Glyph)
		d.placed(byte(u64(in, 0)), glyphs)
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
		if in.Op == bdf.OpUseAt {
			st.m = st.m.translate(f64(in, 1), f64(in, 2))
		}
		d.run(child)
		d.again = again
		// USE is wrapped in an implicit save and restore, of the context
		// it draws on (which is gone if the object ended its group)
		if !ctx.done {
			ctx.stack = ctx.stack[:min(depth, len(ctx.stack))]
			ctx.unsaved = unsaved
			*st = saved
		}
	case bdf.OpGroupBegin:
		d.groupBegin(f64(in, 0), byte(u64(in, 1)))
	case bdf.OpGroupEnd:
		d.groupEnd()
	case bdf.OpMaskBegin:
		transfer, _ := in.Args[2].([]byte)
		d.maskBegin(byte(u64(in, 0)), bdf.Color(u64(in, 1)), transfer)
	case bdf.OpMaskEnd:
		d.maskEnd()
	case bdf.OpLink:
		if len(d.masks) == 0 {
			url, _ := in.Args[4].(string)
			d.link(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), url)
		}
	}
}

// link makes a rectangle of the page a link, over all that is drawn.
func (d *drawer) link(x, y, w, h float64, url string) {
	m := d.ctx.st.m
	href := d.r.href(d.view, url)
	if href == "" || !m.finite() {
		return
	}
	if data, ok := rectData(x, y, w, h, m.digits()); ok {
		d.links.WriteString(`<a href="` + escape(href) + `"><path d="` + data + `"` + m.transform() + ` fill="none" pointer-events="all"/></a>` + "\n")
	}
}

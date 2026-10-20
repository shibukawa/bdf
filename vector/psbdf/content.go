package psbdf

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/vector/internal/flatten"
)

// style is a fill or stroke style: a color, or a gradient or a pattern of
// an object.
type style struct {
	color bdf.Color
	paint *bdf.Paint
	obj   *bdf.ObjectPart
}

// gstate is the drawing state: what SAVE and RESTORE keep (the state of a
// Canvas), and the color PostScript has, which gsave and grestore keep
// too.
type gstate struct {
	ctm          matrix     // from the space of the page
	clip         [4]float64 // the box of the clip, on the page
	fill, stroke style
	alpha        float64
	blend        byte
	width        float64
	cap, join    byte
	miter        float64
	dash         bool
	font         *bdf.Font
	size         float64
	align        byte
	baseline     byte
	dir          byte
	spacing      float64
	smooth       bool
	filter       bool
	// the shadow: its color, and its blur and its offset in units of
	// the page
	shadow     bdf.Color
	shadowBlur float64
	shadowDX   float64
	shadowDY   float64

	// the color PostScript has ("" when not known)
	pColor string
}

func initialState() gstate {
	return gstate{ctm: identity, clip: everywhere, fill: style{color: 0x000000ff}, stroke: style{color: 0x000000ff}, alpha: 1, width: 1,
		miter: 10, size: 10, smooth: true}
}

// shadowOn reports whether what is drawn casts a shadow.
func (st *gstate) shadowOn() bool {
	return uint32(st.shadow)&0xff != 0 && (st.shadowBlur > 0 || st.shadowDX != 0 || st.shadowDY != 0)
}

// reach is how far a line reaches beyond its path, in units of the state.
func (st *gstate) reach() float64 {
	r := st.width / 2
	if st.join == bdf.JoinMiter {
		r *= math.Max(math.Sqrt2, st.miter)
	} else {
		r *= math.Sqrt2
	}
	return r
}

// frame is an open group.
type frame struct {
	// index is the group's among those the first pass decided
	index       int
	transparent bool
	box         [4]float64 // what is drawn in it, on the page
	rect        [4]float64 // its rectangle on the page
	clip        [4]float64
	// the shadow of the state it is drawn with, in units of the page
	shadowDX, shadowDY, shadowBlur float64
	// what ends with it: the state, and the depth of the saved states
	saved          gstate
	stack          []gstate
	floor, unsaved int
	maskFloor      int
	hidden         bool
}

// Limits of what a document may ask of the writer.
const (
	// maxReusedInstructions bounds the instructions written again from
	// objects that were written before.
	maxReusedInstructions = 1 << 22
	// maxLayers bounds the groups and soft masks that are open at a time.
	maxLayers = 64
	// maxStates bounds the states kept for RESTORE.
	maxStates = 1024
	// maxShadowBlur bounds the blur of a shadow, in units of the page.
	maxShadowBlur = 1024
)

type exporter struct {
	doc  *bdf.Document
	opts Options

	objects  map[bdf.Hash]*bdf.ObjectPart
	coll     map[bdf.Hash][]*bdf.Path
	images   map[bdf.Hash]*psImage
	shadings map[*bdf.Paint]*shading
	fonts    *fonts
	renderer *imagebdf.Renderer
	views    map[*bdf.View]*bdf.View
	pictures int
	warnings []string
	warned   map[string]bool

	// the pass: b is the page being written, nil in the first pass, which
	// finds the regions of what PostScript cannot draw
	b       *bytes.Buffer
	regions *flatten.Regions

	st     gstate
	stack  []gstate
	obj    *bdf.ObjectPart
	depth  int
	alt    bool
	floor  int
	unsave int
	groups []*frame
	// decisions are the groups drawn as pictures, in the order they
	// begin, which the first pass makes and the second reads (next)
	decisions []bool
	next      int
	// mask counts the soft masks being read (their instructions are not
	// drawn), hidden the open groups that are pictures (the second pass
	// writes nothing in them), skipped what is left out past the limits
	mask    int
	masks   []openMask
	hidden  int
	skipped int
	// the groups and soft masks open where the object being written
	// began (or the group), which it does not end
	groupFloor, maskFloor int
	// seen holds the objects written before, again says that the object
	// being written is one, and reused counts the instructions of such
	seen   map[bdf.Hash]struct{}
	again  bool
	reused int
}

func newExporter(doc *bdf.Document, opts *Options) *exporter {
	ex := &exporter{doc: doc, objects: map[bdf.Hash]*bdf.ObjectPart{}, coll: map[bdf.Hash][]*bdf.Path{}, images: map[bdf.Hash]*psImage{},
		views: map[*bdf.View]*bdf.View{}, warned: map[string]bool{}}
	if opts != nil {
		ex.opts = *opts
	}
	ex.fonts = newFonts(ex)
	return ex
}

func (ex *exporter) warnf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if !ex.warned[s] {
		ex.warned[s] = true
		ex.warnings = append(ex.warnings, s)
	}
}

// begin starts a pass over a page: the first (b nil) finds the regions,
// the second writes to b.
func (ex *exporter) begin(b *bytes.Buffer, regions *flatten.Regions) {
	ex.b, ex.regions = b, regions
	ex.st, ex.stack, ex.groups = initialState(), nil, nil
	ex.floor, ex.unsave, ex.mask, ex.hidden, ex.skipped = 0, 0, 0, 0, 0
	ex.masks, ex.groupFloor, ex.maskFloor = nil, 0, 0
	ex.seen, ex.again, ex.reused = map[bdf.Hash]struct{}{}, false, 0
	if b == nil {
		ex.decisions = nil
	}
	ex.next = 0
}

// writing reports whether PostScript is written now: in the second pass,
// outside soft masks and groups drawn as pictures.
func (ex *exporter) writing() bool { return ex.b != nil && ex.mask == 0 && ex.hidden == 0 }

// op writes a line of PostScript, when it is written now.
func (ex *exporter) op(format string, args ...any) {
	if !ex.writing() {
		return
	}
	if len(args) == 0 {
		ex.b.WriteString(format)
	} else {
		fmt.Fprintf(ex.b, format, args...)
	}
	ex.b.WriteByte('\n')
}

func (ex *exporter) object(h bdf.Hash) *bdf.ObjectPart {
	if o, ok := ex.objects[h]; ok {
		return o
	}
	ex.objects[h] = nil
	p := ex.doc.Part(h)
	if p == nil {
		ex.warnf("object %s is missing", h)
		return nil
	}
	o, err := bdf.DecodeObject(p.Data)
	if err != nil {
		ex.warnf("object %s: %v", h, err)
		return nil
	}
	ex.objects[h] = o
	return o
}

// run writes the instructions of an object.
func (ex *exporter) run(o *bdf.ObjectPart) {
	if ex.depth > bdf.MaxUseDepth {
		ex.warnf("objects nested too deep")
		return
	}
	ex.depth++
	prev := ex.obj
	ex.obj = o
	if err := o.Walk(ex.instr); err != nil {
		ex.warnf("object: %v", err)
	}
	ex.obj = prev
	ex.depth--
}

// top writes a layer of a page, from the initial state.
func (ex *exporter) top(h bdf.Hash) {
	o := ex.object(h)
	if o == nil {
		return
	}
	saved, stack, floor, unsave := ex.st, ex.stack, ex.floor, ex.unsave
	ex.op("q")
	st := initialState()
	st.ctm, st.clip = saved.ctm, saved.clip
	ex.st, ex.stack = st, nil
	ex.floor, ex.unsave, ex.skipped = 0, 0, 0
	groups, masks := len(ex.groups), len(ex.masks)
	again := ex.again
	ex.again = ex.enter(h)
	ex.run(o)
	ex.again = again
	ex.unwind(groups, masks)
	for range ex.stack {
		ex.op("Q")
	}
	ex.op("Q")
	ex.st, ex.stack, ex.floor, ex.unsave = saved, stack, floor, unsave
}

// unwind ends the groups and soft masks that an object left open.
func (ex *exporter) unwind(groups, masks int) {
	for len(ex.groups) > groups {
		ex.groupEnd()
	}
	for len(ex.masks) > masks {
		ex.maskEnd()
	}
}

func (ex *exporter) enter(h bdf.Hash) bool {
	if _, ok := ex.seen[h]; ok {
		return true
	}
	ex.seen[h] = struct{}{}
	return false
}

func f64(in bdf.Instr, i int) float64 {
	if v, ok := in.Args[i].(float32); ok {
		if f := float64(v); !math.IsNaN(f) && !math.IsInf(f, 0) {
			return f
		}
	}
	return 0
}

func u64(in bdf.Instr, i int) uint64 {
	if v, ok := in.Args[i].(uint64); ok {
		return v
	}
	return 0
}

// leftOut takes an instruction of what is not written, counting what
// begins and ends.
func (ex *exporter) leftOut(in bdf.Instr) {
	switch in.Op {
	case bdf.OpSave:
		ex.unsave++
	case bdf.OpRestore:
		ex.restore()
	case bdf.OpGroupBegin, bdf.OpMaskBegin:
		ex.skipped++
	case bdf.OpGroupEnd:
		if ex.skipped > 0 {
			ex.skipped--
		} else {
			ex.groupEnd()
		}
	case bdf.OpMaskEnd:
		if ex.skipped > 0 {
			ex.skipped--
		} else {
			ex.maskEnd()
		}
	}
}

func (ex *exporter) save() {
	if len(ex.stack) >= maxStates {
		ex.warnf("more than %d states are saved at a time: the rest are not kept", maxStates)
		ex.unsave++
		return
	}
	ex.stack = append(ex.stack, ex.st)
	ex.op("q")
}

func (ex *exporter) restore() {
	if ex.unsave > 0 {
		ex.unsave--
		return
	}
	if n := len(ex.stack); n > ex.floor {
		ex.st = ex.stack[n-1]
		ex.stack = ex.stack[:n-1]
		ex.op("Q")
	}
}

func (ex *exporter) instr(in bdf.Instr) {
	if ex.again {
		if ex.reused >= maxReusedInstructions {
			ex.warnf("objects are drawn again too many times: the rest of them is not written")
			ex.leftOut(in)
			return
		}
		ex.reused++
	}
	if ex.skipped > 0 {
		ex.leftOut(in)
		return
	}
	st := &ex.st
	switch in.Op {
	case bdf.OpSave:
		ex.save()
	case bdf.OpRestore:
		ex.restore()
	case bdf.OpTransform:
		ex.concat(matrix{f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4), f64(in, 5)})
	case bdf.OpTranslate:
		ex.concat(matrix{1, 0, 0, 1, f64(in, 0), f64(in, 1)})
	case bdf.OpScale:
		ex.concat(matrix{f64(in, 0), 0, 0, f64(in, 1), 0, 0})
	case bdf.OpClipPath:
		p := newPen(ex.buf(), identity)
		if bp := ex.path(int(u64(in, 0))); bp != nil {
			p.path(bp)
		}
		if !p.drawn {
			// nothing is inside a path without segments
			ex.op("0 0 0 0 re")
		}
		ex.clipBox(p)
		if byte(u64(in, 1)) == bdf.EvenOdd {
			ex.op("W* n")
		} else {
			ex.op("W n")
		}
	case bdf.OpClipRect:
		x, y, w, h := f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)
		st.clip = intersect(st.clip, st.ctm.bounds(x, y, x+w, y+h))
		ex.op("%s re W n", nums(x, y, w, h))
	case bdf.OpFillColor:
		st.fill = style{color: bdf.Color(u64(in, 0))}
	case bdf.OpStrokeColor:
		st.stroke = style{color: bdf.Color(u64(in, 0))}
	case bdf.OpFillPaint:
		st.fill = ex.paintStyle(int(u64(in, 0)))
	case bdf.OpStrokePaint:
		st.stroke = ex.paintStyle(int(u64(in, 0)))
	case bdf.OpLine:
		if w := f64(in, 0); w > 0 {
			st.width = w
			ex.op("%s w", num(w))
		}
		if c := u64(in, 1); c <= 2 {
			st.cap = byte(c)
			ex.op("%d J", c)
		}
		if j := u64(in, 2); j <= 2 {
			st.join = byte(j)
			ex.op("%d j", j)
		}
		if m := f64(in, 3); m >= 1 {
			st.miter = m
			ex.op("%s M", num(m))
		}
	case bdf.OpDash:
		segs, _ := in.Args[0].([]float32)
		off, _ := in.Args[1].(float32)
		ex.dash(segs, off)
	case bdf.OpAlpha:
		if a := f64(in, 0); a >= 0 && a <= 1 {
			st.alpha = a
		}
	case bdf.OpBlend:
		if b := u64(in, 0); b < uint64(len(bdf.BlendNames)) {
			st.blend = byte(b)
		}
	case bdf.OpShadow:
		// the blur and the offset are in the units of where they are set
		k := st.ctm.areaScale()
		if k == 0 || math.IsNaN(k) || math.IsInf(k, 0) {
			k = 1
		}
		st.shadow = bdf.Color(u64(in, 0))
		st.shadowBlur = math.Min(math.Max(f64(in, 1)*k, 0), maxShadowBlur)
		st.shadowDX, st.shadowDY = clampShadow(f64(in, 2)*k), clampShadow(f64(in, 3)*k)
	case bdf.OpFilter:
		f, _ := in.Args[0].(string)
		st.filter = f != "" && f != "none"
	case bdf.OpFont:
		if i := u64(in, 0); i < uint64(len(ex.obj.Fonts)) {
			f := ex.obj.Fonts[i]
			st.font = &f
		}
		st.size = f64(in, 1)
	case bdf.OpTextStyle:
		st.align, st.baseline, st.dir = byte(u64(in, 0)), byte(u64(in, 1)), byte(u64(in, 2))
		st.spacing = f64(in, 3)
	case bdf.OpFillRect:
		x, y, w, h := f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)
		ex.fill([4]float64{math.Min(x, x+w), math.Min(y, y+h), math.Max(x, x+w), math.Max(y, y+h)}, func(b *bytes.Buffer) bool {
			fmt.Fprintf(b, "%s re\n", nums(x, y, w, h))
			return true
		}, bdf.NonZero)
	case bdf.OpStrokeRect:
		x, y, w, h := f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)
		ex.strokePath([4]float64{math.Min(x, x+w), math.Min(y, y+h), math.Max(x, x+w), math.Max(y, y+h)}, func(b *bytes.Buffer) bool {
			fmt.Fprintf(b, "%s re\n", nums(x, y, w, h))
			return true
		})
	case bdf.OpClearRect:
		ex.clearRect(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3))
	case bdf.OpFillPath, bdf.OpFillPathAt:
		bp := ex.path(int(u64(in, 0)))
		if bp == nil {
			return
		}
		m := identity
		if in.Op == bdf.OpFillPathAt {
			m = m.translate(f64(in, 2), f64(in, 3))
		}
		probe := newPen(nil, m)
		probe.path(bp)
		ex.fill(probe.box, func(b *bytes.Buffer) bool {
			p := newPen(b, m)
			p.path(bp)
			return p.drawn
		}, byte(u64(in, 1)))
	case bdf.OpStrokePath:
		if bp := ex.path(int(u64(in, 0))); bp != nil {
			probe := newPen(nil, identity)
			probe.path(bp)
			ex.strokePath(probe.box, func(b *bytes.Buffer) bool {
				p := newPen(b, identity)
				p.path(bp)
				return p.drawn
			})
		}
	case bdf.OpFillPathRun:
		glyphs, _ := in.Args[1].([]bdf.Glyph)
		box := none
		for _, g := range glyphs {
			if bp := ex.path(int(g.Path)); bp != nil {
				probe := newPen(nil, identity.translate(float64(g.X), float64(g.Y)))
				probe.path(bp)
				if probe.boxed {
					box = union(box, probe.box)
				}
			}
		}
		ex.fill(box, func(b *bytes.Buffer) bool {
			drawn := false
			for _, g := range glyphs {
				if bp := ex.path(int(g.Path)); bp != nil {
					p := newPen(b, identity.translate(float64(g.X), float64(g.Y)))
					p.path(bp)
					drawn = drawn || p.drawn
				}
			}
			return drawn
		}, byte(u64(in, 0)))
	case bdf.OpFillText, bdf.OpStrokeText:
		str, _ := in.Args[0].(string)
		ex.text(str, f64(in, 1), f64(in, 2), f64(in, 3), in.Op == bdf.OpStrokeText)
	case bdf.OpImage:
		ex.image(int(u64(in, 0)), nil, f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4))
	case bdf.OpImageSub:
		src := &[4]float64{f64(in, 1), f64(in, 2), f64(in, 3), f64(in, 4)}
		ex.image(int(u64(in, 0)), src, f64(in, 5), f64(in, 6), f64(in, 7), f64(in, 8))
	case bdf.OpSmoothing:
		st.smooth = u64(in, 0) != 0
	case bdf.OpUse, bdf.OpUseAt:
		i := u64(in, 0)
		if i >= uint64(len(ex.obj.Objects)) {
			return
		}
		child := ex.object(ex.obj.Objects[i])
		if child == nil {
			return
		}
		// USE is wrapped in an implicit save and restore
		saved, depth, floor, unsave := *st, len(ex.stack), ex.floor, ex.unsave
		groups, masks, groupFloor, maskFloor := len(ex.groups), len(ex.masks), ex.groupFloor, ex.maskFloor
		ex.op("q")
		ex.floor, ex.unsave = len(ex.stack), 0
		ex.groupFloor, ex.maskFloor = groups, masks
		if in.Op == bdf.OpUseAt {
			ex.concat(matrix{1, 0, 0, 1, f64(in, 1), f64(in, 2)})
		}
		again := ex.again
		ex.again = ex.enter(ex.obj.Objects[i]) || again
		ex.run(child)
		ex.again = again
		ex.unwind(groups, masks)
		ex.groupFloor, ex.maskFloor = groupFloor, maskFloor
		for len(ex.stack) > depth {
			ex.stack = ex.stack[:len(ex.stack)-1]
			ex.op("Q")
		}
		ex.op("Q")
		ex.st = saved
		ex.floor, ex.unsave = floor, unsave
	case bdf.OpGroupBegin:
		ex.groupBegin(f64(in, 0), byte(u64(in, 1)), f64(in, 2), f64(in, 3), f64(in, 4), f64(in, 5))
	case bdf.OpGroupEnd:
		ex.groupEnd()
	case bdf.OpMaskBegin:
		ex.maskBegin()
	case bdf.OpMaskEnd:
		ex.maskEnd()
	}
}

func clampShadow(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return math.Min(math.Max(v, -1e6), 1e6)
}

// buf is the buffer of the second pass while it writes, nil otherwise.
func (ex *exporter) buf() *bytes.Buffer {
	if ex.writing() {
		return ex.b
	}
	return nil
}

// clipBox narrows the box of the clip by a path.
func (ex *exporter) clipBox(p *pen) {
	st := &ex.st
	if !p.boxed {
		st.clip = none
		return
	}
	st.clip = intersect(st.clip, st.ctm.bounds(p.box[0], p.box[1], p.box[2], p.box[3]))
}

// concat multiplies the transform of the state by m.
func (ex *exporter) concat(m matrix) {
	if !m.finite() {
		return
	}
	ex.st.ctm = ex.st.ctm.mul(m)
	ex.op("%s cm", m)
}

// dash sets the dashes of lines as Canvas does: a list with a value that
// is negative or no number is ignored, and one of zeros is a solid line.
func (ex *exporter) dash(segs []float32, off float32) {
	sum := 0.0
	for _, v := range segs {
		f := float64(v)
		if f < 0 || math.IsNaN(f) || math.IsInf(f, 0) {
			return
		}
		sum += f
	}
	st := &ex.st
	if sum == 0 {
		st.dash = false
		ex.op("[] 0 d")
		return
	}
	st.dash = true
	var b strings.Builder
	for i, v := range segs {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(num(float64(v)))
	}
	if len(segs)%2 == 1 {
		// an odd list is repeated
		b.WriteString(" " + b.String())
	}
	o := float64(off)
	if math.IsNaN(o) || math.IsInf(o, 0) {
		o = 0
	}
	if o < 0 {
		period := sum
		if len(segs)%2 == 1 {
			period *= 2
		}
		o = math.Mod(o, period) + period
	}
	ex.op("[%s] %s d", b.String(), num(o))
}

func (ex *exporter) path(i int) *bdf.Path {
	if i < 0 || i >= len(ex.obj.Paths) {
		ex.warnf("bad path reference")
		return nil
	}
	e := ex.obj.Paths[i]
	if e.Inline != nil {
		return e.Inline
	}
	ps, ok := ex.coll[e.Hash]
	if !ok {
		if part := ex.doc.Part(e.Hash); part != nil {
			var err error
			if ps, err = bdf.DecodePathCollection(part.Data); err != nil {
				ex.warnf("path collection %s: %v", e.Hash, err)
			}
		} else {
			ex.warnf("path collection %s is missing", e.Hash)
		}
		ex.coll[e.Hash] = ps
	}
	if int(e.Index) < len(ps) {
		return ps[e.Index]
	}
	return nil
}

func (ex *exporter) paintStyle(i int) style {
	if i < 0 || i >= len(ex.obj.Paints) {
		ex.warnf("bad paint reference")
		return style{}
	}
	return style{paint: &ex.obj.Paints[i], obj: ex.obj}
}

// look tells how something drawn in a style shows: visible that it shows
// at all, plain that PostScript draws it as it is (opaque, without a blend
// mode, a filter or a shadow, in a color or a gradient of opaque stops).
func (ex *exporter) look(sty *style) (visible, plain bool) {
	if ex.mask > 0 {
		return false, false
	}
	st := &ex.st
	alpha := st.alpha
	plainPaint := true
	switch {
	case sty == nil:
	case sty.paint == nil:
		alpha *= alphaOf(sty.color)
	default:
		sh := ex.shading(sty.paint)
		switch {
		case !sh.ok:
			return false, false
		case sh.solid:
			alpha *= alphaOf(sh.color)
		default:
			plainPaint = sh.opaque
		}
	}
	if alpha <= 0 {
		return false, false
	}
	return true, alpha >= 1 && plainPaint && st.blend == bdf.BlendSourceOver && !st.filter && !st.shadowOn()
}

// alphaOf is the alpha of a color, 0 to 1.
func alphaOf(c bdf.Color) float64 { return float64(uint32(c)&0xff) / 255 }

// drawn takes something drawn within box (in the space of the state,
// grown by reach): in the first pass, it is where it is drawn in the open
// groups, and a region when plain is not set; it reports whether the
// second pass writes it.
func (ex *exporter) drawn(box [4]float64, reach float64, plain bool) bool {
	st := &ex.st
	if !(box[2] >= box[0]) || !(box[3] >= box[1]) {
		return false
	}
	pageBox := st.ctm.bounds(box[0]-reach, box[1]-reach, box[2]+reach, box[3]+reach)
	shown := intersect(pageBox, st.clip)
	if !(shown[2] > shown[0]) || !(shown[3] > shown[1]) {
		if !st.shadowOn() {
			return false
		}
	}
	if ex.b == nil {
		region := shown
		if st.shadowOn() {
			// the shadow, beyond the clip of what casts it
			sb := grow([4]float64{pageBox[0] + st.shadowDX, pageBox[1] + st.shadowDY, pageBox[2] + st.shadowDX, pageBox[3] + st.shadowDY}, 2*st.shadowBlur)
			region = union(region, intersect(sb, st.clip))
		}
		for _, g := range ex.groups {
			g.box = union(g.box, region)
		}
		if !plain {
			ex.regions.Add(region)
		}
		return false
	}
	return plain && ex.writing()
}

// color makes the color PostScript has that of a solid color.
func (ex *exporter) color(c bdf.Color) {
	st := &ex.st
	r, g, b := float64(uint32(c)>>24)/255, float64(uint32(c)>>16&0xff)/255, float64(uint32(c)>>8&0xff)/255
	var key string
	if ex.opts.Gray {
		key = fine(luma(r, g, b)) + " g"
	} else if r == g && g == b {
		key = fine(r) + " g"
	} else {
		key = fine(r) + " " + fine(g) + " " + fine(b) + " rg"
	}
	if st.pColor != key && ex.writing() {
		ex.op("%s", key)
		st.pColor = key
	}
}

// luma is the luminance of Rec. 601, as raster/imagebdf makes gray.
func luma(r, g, b float64) float64 { return 0.299*r + 0.587*g + 0.114*b }

// paintWith writes what fills (or strokes) the path that build writes in
// a style: a color, or a gradient through the path as a clip.
func (ex *exporter) paintWith(sty style, build func(b *bytes.Buffer) bool, paintOp, clipOp string) {
	b := ex.b
	if sty.paint == nil {
		ex.color(sty.color)
		if build(b) {
			ex.op("%s", paintOp)
		}
		return
	}
	sh := ex.shading(sty.paint)
	if sh.solid {
		ex.color(sh.color)
		if build(b) {
			ex.op("%s", paintOp)
		}
		return
	}
	ex.op("q")
	if !build(b) {
		b.WriteString("0 0 0 0 re\n")
	}
	ex.op("%s n", clipOp)
	ex.op("%s shfill", ex.shadingDict(sh))
	ex.op("Q")
}

// fill fills the path that build writes (false: it has no segments) with
// the fill style; box is the box of the path.
func (ex *exporter) fill(box [4]float64, build func(b *bytes.Buffer) bool, rule byte) {
	visible, plain := ex.look(&ex.st.fill)
	if !visible || !ex.drawn(box, 0, plain) {
		return
	}
	op, clip := "f", "W"
	if rule == bdf.EvenOdd {
		op, clip = "f*", "W*"
	}
	ex.paintWith(ex.st.fill, build, op, clip)
}

// strokePath strokes the path that build writes with the stroke style.
func (ex *exporter) strokePath(box [4]float64, build func(b *bytes.Buffer) bool) {
	visible, plain := ex.look(&ex.st.stroke)
	if !visible || !ex.drawn(box, ex.st.reach(), plain) {
		return
	}
	ex.paintWith(ex.st.stroke, build, "S", "strokepath W")
}

// clearRect makes a rectangle of a page white again: a page has no pixels
// to clear. In a group it clears what the group drew, which only a picture
// shows.
func (ex *exporter) clearRect(x, y, w, h float64) {
	if ex.mask > 0 {
		return
	}
	if len(ex.groups) > 0 {
		if ex.b == nil {
			ex.groups[len(ex.groups)-1].transparent = true
		}
		return
	}
	if !ex.drawn([4]float64{math.Min(x, x+w), math.Min(y, y+h), math.Max(x, x+w), math.Max(y, y+h)}, 0, true) {
		return
	}
	ex.op("q 1 g %s re f Q", nums(x, y, w, h))
}

// layer reports whether one more group or soft mask may be open; one past
// the limit is left out with what is in it.
func (ex *exporter) layer() bool {
	if len(ex.groups)+ex.mask >= maxLayers {
		ex.warnf("more than %d groups and soft masks are open: what is in the rest of them is not written", maxLayers)
		ex.skipped++
		return false
	}
	return true
}

// groupBegin starts a group: drawn as it is when it is opaque and without
// a blend mode, a shadow, a filter or a soft mask; as a picture otherwise.
// The group starts from the initial state, but for the transform, the
// colors, the width of lines and the font.
func (ex *exporter) groupBegin(alpha float64, blend byte, x, y, w, h float64) {
	if !ex.layer() {
		return
	}
	st := &ex.st
	f := &frame{box: none, rect: everywhere, clip: st.clip, saved: *st, stack: ex.stack, floor: ex.floor, unsaved: ex.unsave, maskFloor: ex.maskFloor}
	if w > 0 && h > 0 {
		f.rect = st.ctm.bounds(x, y, x+w, y+h)
	}
	if st.shadowOn() {
		f.shadowDX, f.shadowDY, f.shadowBlur = st.shadowDX, st.shadowDY, st.shadowBlur
	}
	f.transparent = alpha < 1 || blend != bdf.BlendSourceOver || st.shadowOn() || st.filter
	if ex.b == nil {
		f.index = len(ex.decisions)
		ex.decisions = append(ex.decisions, f.transparent)
	} else {
		if ex.next < len(ex.decisions) {
			f.transparent = ex.decisions[ex.next]
		}
		ex.next++
		f.hidden = f.transparent
	}
	ex.groups = append(ex.groups, f)
	if f.hidden {
		ex.hidden++
	} else {
		ex.op("q")
		if st.dash {
			ex.op("[] 0 d")
		}
		if st.cap != bdf.CapButt {
			ex.op("0 J")
		}
		if st.join != bdf.JoinMiter {
			ex.op("0 j")
		}
		if st.miter != 10 {
			ex.op("10 M")
		}
	}
	g := initialState()
	g.ctm, g.clip, g.fill, g.stroke, g.width, g.font, g.size = st.ctm, st.clip, st.fill, st.stroke, st.width, st.font, st.size
	g.pColor = st.pColor
	ex.st, ex.stack, ex.floor, ex.unsave = g, nil, 0, 0
	ex.maskFloor = len(ex.masks)
}

// groupEnd ends the innermost group, when the object being written began
// it.
func (ex *exporter) groupEnd() {
	n := len(ex.groups)
	if n <= ex.groupFloor {
		return
	}
	f := ex.groups[n-1]
	// a soft mask the group did not end, ends with it
	for len(ex.masks) > ex.maskFloor {
		ex.maskEnd()
	}
	ex.groups = ex.groups[:n-1]
	ex.maskFloor = f.maskFloor
	if f.hidden {
		ex.hidden--
	} else {
		for range ex.stack {
			ex.op("Q")
		}
		ex.op("Q")
	}
	ex.st, ex.stack, ex.floor, ex.unsave = f.saved, f.stack, f.floor, f.unsaved
	if ex.b == nil {
		ex.decisions[f.index] = f.transparent
		if f.transparent {
			region := intersect(intersect(f.box, f.rect), f.clip)
			if f.shadowBlur > 0 || f.shadowDX != 0 || f.shadowDY != 0 {
				sb := grow([4]float64{region[0] + f.shadowDX, region[1] + f.shadowDY, region[2] + f.shadowDX, region[3] + f.shadowDY}, 2*f.shadowBlur)
				region = union(region, intersect(sb, f.clip))
			}
			ex.regions.Add(region)
			for _, g := range ex.groups {
				g.box = union(g.box, region)
			}
		}
	}
}

// openMask is the state a soft mask is read from, which its end returns
// to.
type openMask struct {
	saved         gstate
	stack         []gstate
	floor, unsave int
}

// maskBegin starts a soft mask: what it draws is not drawn, and the group
// it is the mask of is a picture. The mask is read from the initial state,
// with the transform.
func (ex *exporter) maskBegin() {
	if !ex.layer() {
		return
	}
	if n := len(ex.groups); n > 0 && ex.b == nil {
		ex.groups[n-1].transparent = true
	}
	st := &ex.st
	ex.masks = append(ex.masks, openMask{saved: *st, stack: ex.stack, floor: ex.floor, unsave: ex.unsave})
	m := initialState()
	m.ctm, m.pColor = st.ctm, st.pColor
	ex.st, ex.stack, ex.floor, ex.unsave = m, nil, 0, 0
	ex.mask++
}

// maskEnd ends the innermost soft mask, when the object or the group
// being written began it.
func (ex *exporter) maskEnd() {
	if len(ex.masks) <= ex.maskFloor {
		return
	}
	ex.mask--
	m := ex.masks[len(ex.masks)-1]
	ex.masks = ex.masks[:len(ex.masks)-1]
	ex.st, ex.stack, ex.floor, ex.unsave = m.saved, m.stack, m.floor, m.unsave
}

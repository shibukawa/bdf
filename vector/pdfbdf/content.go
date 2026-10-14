package pdfbdf

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/shibukawa/bdf"
)

// style is a fill or stroke style: a color, or a gradient or a pattern of
// an object.
type style struct {
	color bdf.Color
	paint *bdf.Paint
	obj   *bdf.ObjectPart
}

// gstate is the drawing state of a content stream: what SAVE and RESTORE
// keep (the state of a Canvas), and what the graphics state of the PDF
// holds at this point of the stream, which q and Q keep too.
type gstate struct {
	ctm          matrix // from the space the stream starts in
	fill, stroke style
	alpha        float64
	blend        byte
	width        float64
	cap, join    byte
	miter        float64
	dash         []float32
	dashOff      float32
	font         *bdf.Font
	size         float64
	align        byte
	baseline     byte
	dir          byte
	spacing      float64
	smooth       bool
	// the shadow: its color, and its blur and its offset in units of
	// the page (a shadow is not turned or stretched with what casts it)
	shadow     bdf.Color
	shadowBlur float64
	shadowDX   float64
	shadowDY   float64

	// the PDF graphics state: the colors ("" when not known), the alpha
	// constants (-1 when not known), the blend mode and the character
	// spacing
	pFill, pStroke string
	pCa, pCA       float64
	pBM            byte
	pTc            float64
}

func initialState() gstate {
	return gstate{ctm: identity, fill: style{color: 0x000000ff}, stroke: style{color: 0x000000ff}, alpha: 1, width: 1, miter: 10, size: 10, smooth: true,
		pFill: "c000000", pStroke: "c000000", pCa: 1, pCA: 1}
}

// shadowOn reports whether what is drawn casts a shadow.
func (st *gstate) shadowOn() bool {
	return uint32(st.shadow)&0xff != 0 && (st.shadowBlur > 0 || st.shadowDX != 0 || st.shadowDY != 0)
}

// link is a LINK of a stream: its rectangle in the space of the page.
type link struct {
	x0, y0, x1, y1 float64
	url            string
	// elem is its element in the structure tree of a tagged file.
	elem *structElem
}

// stream is a content stream being written: of a page, of the form of a
// layer, of a group or of a soft mask.
type stream struct {
	b     bytes.Buffer
	res   map[string]map[string]ref // by the kind of resource and its name
	st    gstate
	stack []gstate
	// toPage maps the space the stream starts in to that of the page of
	// the document, for the rectangles of links
	toPage matrix
	// base maps the space the stream starts in to the default space of
	// the page of the file. The forms of the file are all drawn in that
	// space, and hold their transforms themselves: the matrices of the
	// patterns and the soft masks in them are then those to the space of
	// the page, whichever of the two a viewer takes them for (Quartz
	// takes the matrix of a pattern in a form for one to the page, and
	// cuts a soft mask that is set in another space to the rectangle of
	// the page in that space).
	base  matrix
	links []link
	// done says that the group or the soft mask of the stream has ended
	done bool

	// of a tagged file (tags.go): the q that are open, the marked content
	// that is open, the elements of the marked content of the stream by
	// their numbers, the key of the stream in the parent tree (-1 without
	// marked content), and the form of the stream once it is written
	qd    int
	mc    mcState
	mcids []*structElem
	key   int
	ref   ref
	// page is set for the content stream of a page: its marked content is
	// that of the page, not of a form
	page bool
}

func newStream(toPage, base matrix) *stream {
	return &stream{res: map[string]map[string]ref{}, st: initialState(), toPage: toPage, base: base, key: -1}
}

func (s *stream) op(format string, args ...any) {
	if len(args) == 0 {
		s.b.WriteString(format)
	} else {
		fmt.Fprintf(&s.b, format, args...)
	}
	s.b.WriteByte('\n')
}

// use notes a resource of the stream and returns its name as a name
// object.
func (s *stream) use(kind, name string, r ref) string {
	m := s.res[kind]
	if m == nil {
		m = map[string]ref{}
		s.res[kind] = m
	}
	m[name] = r
	return "/" + name
}

// resources writes the resource dictionary of the stream.
func (s *stream) resources() string {
	var b strings.Builder
	b.WriteString("<<")
	kinds := make([]string, 0, len(s.res))
	for k := range s.res {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		names := make([]string, 0, len(s.res[k]))
		for n := range s.res[k] {
			names = append(names, n)
		}
		sort.Strings(names)
		b.WriteString(" /" + k + " <<")
		for _, n := range names {
			b.WriteString(" /" + n + " " + s.res[k][n].String())
		}
		b.WriteString(" >>")
	}
	b.WriteString(" >>")
	return b.String()
}

// group is an open GROUP_BEGIN: the stream of what is in it, and how it
// is drawn on the stream around it when it ends.
type group struct {
	parent     *stream
	s          *stream
	alpha      float64
	blend      byte
	x, y, w, h float64
	// to maps the space of the group to the default space of the page,
	// which its form is drawn in
	to    matrix
	masks []softMask
	// cast is set for a group in a state that has a shadow: what the
	// group draws, for the shadow of the group as a whole
	cast *groupCast
}

// groupCast holds the instructions of a group that casts a shadow: those
// of the object the group is in, between its GROUP_BEGIN and its GROUP_END.
type groupCast struct {
	obj   *bdf.ObjectPart
	depth int
	ins   []bdf.Instr
	nest  int  // the groups open inside it
	over  bool // more instructions than a shadow is drawn of
	done  bool // its GROUP_END was read
}

// record keeps an instruction for the groups it is in that cast a shadow.
func (ex *exporter) record(in bdf.Instr) {
	for _, g := range ex.groups {
		c := g.cast
		if c == nil || c.over || c.done || c.obj != ex.obj || c.depth != ex.depth {
			continue
		}
		switch in.Op {
		case bdf.OpGroupBegin:
			c.nest++
		case bdf.OpGroupEnd:
			if c.nest == 0 {
				c.done = true // its own
				continue
			}
			c.nest--
		}
		if len(c.ins) >= maxCastInstructions {
			c.over, c.ins = true, nil
			ex.warnf("a group of more than %d instructions casts no shadow", maxCastInstructions)
			continue
		}
		c.ins = append(c.ins, in)
	}
}

// softMask is a soft mask of a group: the ExtGState that sets it.
type softMask struct {
	state ref
}

// openMask is an open MASK_BEGIN.
type openMask struct {
	parent   *stream
	s        *stream
	group    *group
	kind     byte
	backdrop bdf.Color
	transfer []byte
}

// Limits of what a document may ask of the writer: a document decides how
// often its objects are drawn and how deep its groups are.
const (
	// maxReusedInstructions bounds the instructions written again from
	// objects that were written before (each is written as many times as
	// it is drawn: what is drawn once is bounded by the size of the
	// document).
	maxReusedInstructions = 1 << 22
	// maxCastInstructions bounds the instructions of a group that are kept
	// to draw its shadow of.
	maxCastInstructions = 1 << 20
	// maxLayers bounds the groups and soft masks that are open at a time.
	maxLayers = 64
	// maxStates bounds the states a stream keeps for RESTORE.
	maxStates = 1024
)

// run writes the instructions of an object on the current stream.
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

// top writes a top-level object (a layer of a page, a tile of a sheet) on
// the current stream, from the initial state.
func (ex *exporter) top(h bdf.Hash) {
	o := ex.object(h)
	if o == nil {
		return
	}
	s := ex.s
	saved, stack, floor, unsaved := s.st, s.stack, ex.floor, ex.unsaved
	s.q()
	if ex.tags != nil {
		ex.tags.begin()
	}
	st := initialState()
	st.ctm = saved.ctm
	s.st, s.stack = st, nil
	ex.floor, ex.unsaved, ex.skipped = 0, 0, 0
	groups, masks := len(ex.groups), len(ex.masks)
	again := ex.again
	ex.again = ex.enter(h)
	ex.run(o)
	ex.again = again
	ex.unwind(groups, masks)
	ex.s = s
	for range s.stack {
		s.Q()
	}
	s.Q()
	if ex.tags != nil {
		// (what the object left open ends with it)
		ex.tags.begin()
	}
	s.st, s.stack = saved, stack
	ex.floor, ex.unsaved, ex.skipped = floor, unsaved, 0
	ex.alt = nil
}

// unwind ends the groups and the soft masks that an object left open.
func (ex *exporter) unwind(groups, masks int) {
	for len(ex.masks) > masks || len(ex.groups) > groups {
		if n := len(ex.masks); n > masks && (len(ex.groups) == 0 || ex.masks[n-1].group == ex.groups[len(ex.groups)-1] || len(ex.groups) <= groups) {
			ex.maskEnd()
			continue
		}
		ex.groupEnd()
	}
}

// enter reports whether an object was written before, and notes that it
// is now.
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

// leftOut takes an instruction of what is not written: an object drawn
// again past maxReusedInstructions, or what is in a group past the limits.
// It counts what begins and ends, so that what ends later ends what it
// began.
func (ex *exporter) leftOut(in bdf.Instr) {
	switch in.Op {
	case bdf.OpSave:
		ex.unsaved++
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
	s := ex.s
	if len(s.stack) >= maxStates {
		ex.warnf("more than %d states are saved at a time: the rest are not kept", maxStates)
		ex.unsaved++
		return
	}
	s.stack = append(s.stack, s.st)
	s.q()
}

func (ex *exporter) restore() {
	if ex.unsaved > 0 {
		ex.unsaved--
		return
	}
	s := ex.s
	if n := len(s.stack); n > ex.floor {
		s.st = s.stack[n-1]
		s.stack = s.stack[:n-1]
		s.Q()
	}
}

func (ex *exporter) instr(in bdf.Instr) {
	if len(ex.groups) > 0 {
		ex.record(in)
	}
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
	s := ex.s
	st := &s.st
	// the alternative text of a MARK is that of the drawing after it
	alt := ex.alt
	if in.Op != bdf.OpMark {
		ex.alt = nil
	}
	if ex.tags != nil {
		// what is drawn is the content of an element of the structure
		// tree, or an artifact
		switch in.Op {
		case bdf.OpFillText, bdf.OpStrokeText:
			ex.tag(true)
		case bdf.OpFillPathAt, bdf.OpFillPathRun:
			ex.tag(alt != nil)
		case bdf.OpUse, bdf.OpUseAt:
			if alt != nil {
				ex.tag(true)
			}
		case bdf.OpFillRect, bdf.OpStrokeRect, bdf.OpFillPath, bdf.OpStrokePath, bdf.OpClearRect, bdf.OpImage, bdf.OpImageSub:
			ex.tag(false)
		case bdf.OpMark:
			payload, _ := in.Args[1].(string)
			ex.tags.mark(byte(u64(in, 0)), payload)
		}
	}
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
		p := &pen{b: &s.b, m: identity}
		if bp := ex.path(int(u64(in, 0))); bp != nil {
			p.path(bp)
		}
		if !p.drawn {
			// nothing is inside a path without segments
			s.op("0 0 0 0 re")
		}
		if byte(u64(in, 1)) == bdf.EvenOdd {
			s.op("W* n")
		} else {
			s.op("W n")
		}
	case bdf.OpClipRect:
		s.op("%s re W n", nums(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)))
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
			s.op("%s w", num(w))
		}
		if c := u64(in, 1); c <= 2 {
			st.cap = byte(c)
			s.op("%d J", c)
		}
		if j := u64(in, 2); j <= 2 {
			st.join = byte(j)
			s.op("%d j", j)
		}
		if m := f64(in, 3); m >= 1 {
			st.miter = m
			s.op("%s M", num(m))
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
		k := s.toPage.mul(st.ctm).areaScale()
		if k == 0 || math.IsNaN(k) || math.IsInf(k, 0) {
			k = 1
		}
		st.shadow = bdf.Color(u64(in, 0))
		st.shadowBlur = math.Min(math.Max(f64(in, 1)*k, 0), maxShadowBlur)
		st.shadowDX, st.shadowDY = clampShadow(f64(in, 2)*k), clampShadow(f64(in, 3)*k)
	case bdf.OpFilter:
		if f, _ := in.Args[0].(string); f != "" && f != "none" {
			ex.warnf("FILTER %q is not applied", f)
		}
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
		ex.castShadow([4]float64{x, y, x + w, y + h}, 0, func(o *bdf.Object) { o.FillRect(float32(x), float32(y), float32(w), float32(h)) })
		ex.fill(func() bool {
			s.op("%s re", nums(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)))
			return true
		}, bdf.NonZero)
	case bdf.OpStrokeRect:
		x, y, w, h := f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)
		ex.castShadow([4]float64{x, y, x + w, y + h}, st.reach(), func(o *bdf.Object) { o.StrokeRect(float32(x), float32(y), float32(w), float32(h)) })
		ex.strokePath(func() bool {
			s.op("%s re", nums(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3)))
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
		if box, ok := pathBounds(bp); ok && st.shadowOn() {
			ex.castShadow(m.bounds(box[0], box[1], box[2], box[3]), 0, func(o *bdf.Object) {
				o.FillPathAt(o.AddPath(bp), byte(u64(in, 1)), float32(m[4]), float32(m[5]))
			})
		}
		ex.marked(alt, func() {
			ex.fill(func() bool {
				p := &pen{b: &s.b, m: m}
				p.path(bp)
				return p.drawn
			}, byte(u64(in, 1)))
		})
	case bdf.OpStrokePath:
		if bp := ex.path(int(u64(in, 0))); bp != nil {
			if box, ok := pathBounds(bp); ok && st.shadowOn() {
				ex.castShadow(box, st.reach(), func(o *bdf.Object) { o.StrokePath(o.AddPath(bp)) })
			}
			ex.strokePath(func() bool {
				p := &pen{b: &s.b, m: identity}
				p.path(bp)
				return p.drawn
			})
		}
	case bdf.OpFillPathRun:
		glyphs, _ := in.Args[1].([]bdf.Glyph)
		if st.shadowOn() {
			ex.runShadow(byte(u64(in, 0)), glyphs)
		}
		ex.marked(alt, func() {
			ex.fill(func() bool {
				drawn := false
				for _, g := range glyphs {
					if bp := ex.path(int(g.Path)); bp != nil {
						p := &pen{b: &s.b, m: identity.translate(float64(g.X), float64(g.Y))}
						p.path(bp)
						drawn = drawn || p.drawn
					}
				}
				return drawn
			}, byte(u64(in, 0)))
		})
	case bdf.OpFillText, bdf.OpStrokeText:
		str, _ := in.Args[0].(string)
		ex.text(str, f64(in, 1), f64(in, 2), f64(in, 3), in.Op == bdf.OpStrokeText, alt)
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
		if t := ex.tags; t != nil && alt != nil && t.held == nil && s.mc.open && s.mc.elem != nil {
			// what the object draws stands for the text: all of it is the
			// content of the leaf of that text
			t.held = s.mc.elem
			defer func() { t.held = nil }()
		}
		ex.marked(alt, func() {
			// USE is wrapped in an implicit save and restore
			saved, depth, floor, unsaved := *st, len(s.stack), ex.floor, ex.unsaved
			groups, masks := len(ex.groups), len(ex.masks)
			s.q()
			ex.floor, ex.unsaved = len(s.stack), 0
			if in.Op == bdf.OpUseAt {
				ex.concat(matrix{1, 0, 0, 1, f64(in, 1), f64(in, 2)})
			}
			again := ex.again
			ex.again = ex.enter(ex.obj.Objects[i]) || again
			ex.run(child)
			ex.again = again
			ex.unwind(groups, masks)
			if s.done {
				// the object ended the group it was drawn in
				return
			}
			ex.s = s
			for len(s.stack) > depth {
				s.stack = s.stack[:len(s.stack)-1]
				s.Q()
			}
			s.Q()
			*st = saved
			ex.floor, ex.unsaved = floor, unsaved
		})
	case bdf.OpGroupBegin:
		ex.groupBegin(f64(in, 0), byte(u64(in, 1)), f64(in, 2), f64(in, 3), f64(in, 4), f64(in, 5))
	case bdf.OpGroupEnd:
		ex.groupEnd()
	case bdf.OpMaskBegin:
		transfer, _ := in.Args[2].([]byte)
		ex.maskBegin(byte(u64(in, 0)), bdf.Color(u64(in, 1)), transfer)
	case bdf.OpMaskEnd:
		ex.maskEnd()
	case bdf.OpLink:
		url, _ := in.Args[4].(string)
		ex.link(f64(in, 0), f64(in, 1), f64(in, 2), f64(in, 3), url)
	case bdf.OpMark:
		if byte(u64(in, 0)) == bdf.MarkAltText {
			payload, _ := in.Args[1].(string)
			ex.alt = &payload
		}
	}
}

// concat multiplies the transform of the state by m.
func (ex *exporter) concat(m matrix) {
	if !m.finite() {
		return
	}
	s := ex.s
	s.st.ctm = s.st.ctm.mul(m)
	s.op("%s cm", m)
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
	st := &ex.s.st
	st.dash, st.dashOff = nil, 0
	if sum == 0 {
		ex.s.op("[] 0 d")
		return
	}
	st.dash, st.dashOff = append([]float32(nil), segs...), off
	var b strings.Builder
	for i, v := range segs {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(num(float64(v)))
	}
	if len(segs)%2 == 1 {
		// an odd list is repeated, as in a PDF
		b.WriteString(" " + b.String())
	}
	o := float64(off)
	if math.IsNaN(o) || math.IsInf(o, 0) {
		o = 0
	}
	// (the phase of a PDF is not negative)
	if o < 0 {
		period := sum
		if len(segs)%2 == 1 {
			period *= 2
		}
		o = math.Mod(o, period) + period
	}
	ex.s.op("[%s] %s d", b.String(), num(o))
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

// marked writes what draw writes as marked content with an alternative
// text, when there is one: the text that the drawing stands for.
func (ex *exporter) marked(alt *string, draw func()) {
	if alt == nil {
		draw()
		return
	}
	s := ex.s
	s.op("/Span << /ActualText %s >> BDC", text(*alt))
	draw()
	// (what was drawn may have ended on another stream: a group left open)
	s.op("EMC")
}

// scoped writes what draw writes between q and Q: the PDF graphics state
// is as before after it.
func (ex *exporter) scoped(draw func()) {
	s := ex.s
	saved := s.st
	s.q()
	draw()
	s.Q()
	s.st.pFill, s.st.pStroke, s.st.pCa, s.st.pCA, s.st.pBM, s.st.pTc = saved.pFill, saved.pStroke, saved.pCa, saved.pCA, saved.pBM, saved.pTc
}

// blendMode makes the blend mode of the PDF graphics state that of the
// state. The composite operations of Canvas that are not blend modes have
// no counterpart: what they draw is drawn over what is there.
func (ex *exporter) blendMode() {
	s := ex.s
	b := s.st.blend
	if b > bdf.BlendLuminosity {
		ex.warnf("the composite operation %s is drawn as source-over", bdf.BlendNames[b])
		b = bdf.BlendSourceOver
	}
	if b != s.st.pBM {
		s.op("%s gs", ex.extGState("/BM /"+blendNames[b]))
		s.st.pBM = b
	}
}

var blendNames = [...]string{"Normal", "Multiply", "Screen", "Overlay", "Darken", "Lighten", "ColorDodge", "ColorBurn", "HardLight", "SoftLight",
	"Difference", "Exclusion", "Hue", "Saturation", "Color", "Luminosity"}

// alphaOf is the alpha of a color, 0 to 1.
func alphaOf(c bdf.Color) float64 { return float64(uint32(c)&0xff) / 255 }

func rgbOf(c bdf.Color) string {
	return nums(float64(uint32(c)>>24)/255, float64(uint32(c)>>16&0xff)/255, float64(uint32(c)>>8&0xff)/255)
}

// prepared is what a fill or a stroke is drawn with.
type prepared struct {
	ok bool
	// mask writes an ExtGState with a soft mask of the alphas of a
	// gradient whose stops have different ones (nil for none): the mask
	// under a transform to the space of the page
	mask func(toPage matrix) ref
	// clipImage is set for a pattern that does not repeat: its image is
	// drawn where the path is, which is the clip
	clipImage *pdfImage
	matrix    matrix
}

// prep makes the color and the alpha of the PDF graphics state those of
// the fill style (or of the stroke style, for a stroke).
func (ex *exporter) prep(stroke bool) prepared {
	s := ex.s
	st := &s.st
	ex.blendMode()
	sty := st.fill
	cur, curA := &st.pFill, &st.pCa
	setColor, setSpace, setPattern, alphaKey := "rg", "/Pattern cs", "scn", "/ca "
	if stroke {
		sty, cur, curA = st.stroke, &st.pStroke, &st.pCA
		setColor, setSpace, setPattern, alphaKey = "RG", "/Pattern CS", "SCN", "/CA "
	}
	alpha := st.alpha
	out := prepared{ok: true}
	if sty.paint == nil {
		key := fmt.Sprintf("c%06x", uint32(sty.color)>>8)
		if *cur != key {
			s.op("%s %s", rgbOf(sty.color), setColor)
			*cur = key
		}
		alpha *= alphaOf(sty.color)
	} else {
		p := ex.paint(sty, s.base.mul(st.ctm), stroke)
		switch {
		case !p.ok:
			return prepared{}
		case p.solid:
			key := fmt.Sprintf("c%06x", uint32(p.color)>>8)
			if *cur != key {
				s.op("%s %s", rgbOf(p.color), setColor)
				*cur = key
			}
			alpha *= alphaOf(p.color)
		case p.image != nil:
			out.clipImage, out.matrix = p.image, p.matrix
		default:
			if key := "p" + p.name; *cur != key {
				s.op("%s %s %s", setSpace, s.use("Pattern", p.name, p.ref), setPattern)
				*cur = key
			}
			alpha *= p.alpha
			out.mask = p.mask
		}
	}
	if alpha != *curA {
		s.op("%s gs", ex.extGState(alphaKey+fine(alpha)))
		*curA = alpha
	}
	return out
}

// fill fills the path that build writes (false: it has no segments) with
// the fill style.
func (ex *exporter) fill(build func() bool, rule byte) {
	s := ex.s
	op := "f"
	if rule == bdf.EvenOdd {
		op = "f*"
	}
	ex.withPrep(false, func(p prepared) {
		if p.clipImage != nil {
			// the image, where the path is
			ex.scoped(func() {
				if !build() {
					s.op("0 0 0 0 re")
				}
				if rule == bdf.EvenOdd {
					s.op("W* n")
				} else {
					s.op("W n")
				}
				s.op("%s cm %s Do", p.matrix, s.use("XObject", p.clipImage.name, p.clipImage.ref))
			})
			return
		}
		if build() {
			s.op("%s", op)
		}
	})
}

// strokePath strokes the path that build writes with the stroke style.
func (ex *exporter) strokePath(build func() bool) {
	s := ex.s
	ex.withPrep(true, func(p prepared) {
		if p.clipImage != nil {
			ex.warnf("a line in a pattern that does not repeat is not drawn")
			return
		}
		if build() {
			s.op("S")
		}
	})
}

// withPrep draws with the fill or the stroke style prepared: under the
// soft mask of the style, when it has one.
func (ex *exporter) withPrep(stroke bool, draw func(p prepared)) {
	// (the color is set before the q of a soft mask, so that it stays)
	p := ex.prep(stroke)
	if !p.ok {
		return
	}
	if p.mask == nil {
		draw(p)
		return
	}
	ex.scoped(func() {
		if ex.softMask(p.mask) {
			draw(p)
		}
	})
}

// softMask sets a soft mask in the PDF graphics state: an ExtGState that
// mask writes for the transform from the space of the drawing to the
// default space of the page. The mask is set in that space, and its form
// holds the transform (see stream.base). It reports whether the space has
// a way back.
func (ex *exporter) softMask(mask func(toPage matrix) ref) bool {
	s := ex.s
	t := s.base.mul(s.st.ctm)
	inv, ok := t.invert()
	if !ok || !t.finite() {
		return false
	}
	s.op("%s cm %s gs %s cm", inv.exact(), s.use("ExtGState", ex.newName("GS"), mask(t)), t.exact())
	return true
}

// clearRect makes a rectangle of a page white again: a page of a PDF has
// no pixels to clear. In a group, where what is under the group would show,
// it is not done.
func (ex *exporter) clearRect(x, y, w, h float64) {
	if len(ex.groups) > 0 || len(ex.masks) > 0 {
		ex.warnf("CLEAR_RECT in a group is not applied")
		return
	}
	ex.scoped(func() {
		s := ex.s
		s.op("1 1 1 rg")
		if s.st.pCa != 1 {
			s.op("%s gs", ex.extGState("/ca 1"))
		}
		if s.st.pBM != bdf.BlendSourceOver {
			s.op("%s gs", ex.extGState("/BM /Normal"))
		}
		s.op("%s re f", nums(x, y, w, h))
	})
}

// extGState returns the name of an ExtGState with the entries, as a
// resource of the current stream.
func (ex *exporter) extGState(entries string) string {
	g, ok := ex.gstates[entries]
	if !ok {
		g = named{name: ex.newName("GS"), ref: ex.pw.add("<< /Type /ExtGState " + entries + " >>")}
		ex.gstates[entries] = g
	}
	return ex.s.use("ExtGState", g.name, g.ref)
}

// named is an object of the file with its name as a resource.
type named struct {
	name string
	ref  ref
}

func (ex *exporter) newName(prefix string) string {
	ex.names++
	return fmt.Sprintf("%s%d", prefix, ex.names)
}

// layer reports whether one more group or soft mask may be open; one past
// the limit is left out with what is in it.
func (ex *exporter) layer() bool {
	if len(ex.groups)+len(ex.masks) >= maxLayers {
		ex.warnf("more than %d groups and soft masks are open: what is in the rest of them is not written", maxLayers)
		ex.skipped++
		return false
	}
	return true
}

// groupBegin starts the stream of a group: a form that is drawn with the
// alpha and the blend mode of the group when it ends. The space of the
// form is that of the stream around it where the group begins.
func (ex *exporter) groupBegin(alpha float64, blend byte, x, y, w, h float64) {
	if !ex.layer() {
		return
	}
	parent := ex.s
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	if !(alpha >= 0 && alpha <= 1) {
		alpha = 1
	}
	g := &group{parent: parent, alpha: alpha, blend: blend, x: x, y: y, w: w, h: h, to: parent.base.mul(parent.st.ctm)}
	g.s = newStream(parent.toPage.mul(parent.st.ctm), g.to)
	// a group starts from the initial state, with the styles, the line
	// width and the font of the state around it; the form of a group
	// starts with the graphics state it is drawn in, but for the alpha
	// and the blend mode
	st := initialState()
	p := &parent.st
	st.fill, st.stroke, st.width, st.font, st.size = p.fill, p.stroke, p.width, p.font, p.size
	st.pFill, st.pStroke = "", ""
	g.s.st = st
	// the form of the group is drawn in the space of the page: it holds
	// the transform of the group, and is cut to it
	g.s.op("%s cm %s re W n", g.to.exact(), nums(x, y, w, h))
	g.s.op("%s w 0 J 0 j 10 M [] 0 d", num(st.width))
	if p.shadowOn() && ex.obj != nil {
		g.cast = &groupCast{obj: ex.obj, depth: ex.depth}
	}
	ex.groups = append(ex.groups, g)
	ex.s = g.s
	ex.floors = append(ex.floors, ex.floor, ex.unsaved)
	ex.floor, ex.unsaved = 0, 0
}

// finish closes what a stream left open.
func (s *stream) finish() {
	for range s.stack {
		s.Q()
	}
	s.endMC()
	s.stack, s.done = nil, true
}

// Kinds of forms: one that is drawn as what is in it, a transparency
// group (drawn as a whole, with an alpha, a blend mode or a soft mask, or
// as a soft mask of alphas), and a transparency group with the colors
// that the luminosity of a soft mask is taken from.
const (
	formPlain = iota
	formGroup
	formGroupRGB
)

// form writes a stream as a form XObject and returns it.
func (ex *exporter) form(s *stream, bbox [4]float64, kind int) ref {
	s.finish()
	dict := fmt.Sprintf(" /Type /XObject /Subtype /Form /BBox [%s] /Resources %s", nums(bbox[0], bbox[1], bbox[2], bbox[3]), s.resources())
	switch kind {
	case formGroup:
		dict += " /Group << /Type /Group /S /Transparency /I true >>"
	case formGroupRGB:
		dict += " /Group << /Type /Group /S /Transparency /I true /CS /DeviceRGB >>"
	}
	if s.key >= 0 {
		// (the form has marked content of the structure tree)
		dict += fmt.Sprintf(" /StructParents %d", s.key)
	}
	s.ref = ex.pw.addStream(dict, s.b.Bytes(), true)
	return s.ref
}

func (ex *exporter) groupEnd() {
	n := len(ex.groups)
	if n == 0 {
		return
	}
	g := ex.groups[n-1]
	// masks begun inside the group end with it
	for len(ex.masks) > 0 && ex.masks[len(ex.masks)-1].group == g {
		ex.maskEnd()
	}
	ex.groups = ex.groups[:n-1]
	if k := len(ex.floors); k >= 2 {
		ex.floor, ex.unsaved = ex.floors[k-2], ex.floors[k-1]
		ex.floors = ex.floors[:k-2]
	}
	ex.s = g.parent
	s := ex.s
	s.links = append(s.links, g.s.links...)
	if !(g.w > 0) || !(g.h > 0) {
		return // nothing of it shows
	}
	back, ok := g.to.invert()
	if !ok || !g.to.finite() {
		return
	}
	bbox := g.to.bounds(g.x, g.y, g.x+g.w, g.y+g.h)
	bbox = [4]float64{bbox[0] - 1, bbox[1] - 1, bbox[2] + 1, bbox[3] + 1}
	// the shadow of the state around a group is that of the group as a
	// whole: of what it draws, drawn again by itself
	if c := g.cast; c != nil && !c.over && len(c.ins) > 0 {
		ex.castShadow([4]float64{g.x, g.y, g.x + g.w, g.y + g.h}, 0, func(o *bdf.Object) {
			ex.shadows.deps(ex, c.obj)
			o.GroupBegin(float32(g.alpha), bdf.BlendSourceOver, float32(g.x), float32(g.y), float32(g.w), float32(g.h))
			rp := bdf.NewReplay(c.obj, o)
			for _, in := range c.ins {
				if rp.Emit(in) != nil {
					break
				}
			}
			o.GroupEnd()
		})
	}
	form := ex.form(g.s, bbox, formGroupRGB)
	// every soft mask but the last is that of a form around the one before
	for _, m := range g.masks[:max(len(g.masks)-1, 0)] {
		in := newStream(identity, identity)
		in.op("%s gs %s Do", in.use("ExtGState", ex.newName("GS"), m.state), in.use("XObject", ex.newName("Fm"), form))
		form = ex.form(in, bbox, formGroupRGB)
	}
	blend := g.blend
	if blend > bdf.BlendLuminosity {
		if int(blend) < len(bdf.BlendNames) {
			ex.warnf("the composite operation %s is drawn as source-over", bdf.BlendNames[blend])
		}
		blend = bdf.BlendSourceOver
	}
	// (the content of the form is marked in the form: the form is not
	// drawn as the content of what is open here)
	s.endMC()
	ex.scoped(func() {
		// the alpha and the blend mode are those of the group alone
		if s.st.pCa != g.alpha {
			s.op("%s gs", ex.extGState("/ca "+fine(g.alpha)))
		}
		if s.st.pBM != blend {
			s.op("%s gs", ex.extGState("/BM /"+blendNames[blend]))
		}
		// in the space of the page, where the form and its soft mask are
		s.op("%s cm", back.exact())
		if k := len(g.masks); k > 0 {
			s.op("%s gs", s.use("ExtGState", ex.newName("GS"), g.masks[k-1].state))
		}
		s.op("%s Do", s.use("XObject", ex.newName("Fm"), form))
	})
}

// maskBegin starts the stream of a soft mask of the innermost group: what
// is drawn until MASK_END is the mask of the group, in the space of the
// group where the mask begins.
func (ex *exporter) maskBegin(kind byte, backdrop bdf.Color, transfer []byte) {
	if !ex.layer() {
		return
	}
	parent := ex.s
	var g *group
	if n := len(ex.groups); n > 0 {
		g = ex.groups[n-1]
	}
	m := &openMask{parent: parent, group: g, kind: kind, backdrop: backdrop, transfer: transfer}
	m.s = newStream(parent.toPage, parent.base)
	st := initialState()
	st.ctm = identity
	st.pFill, st.pStroke = "", ""
	m.s.st = st
	m.s.op("1 w 0 J 0 j 10 M [] 0 d")
	ex.masks = append(ex.masks, m)
	ex.s = m.s
	if ex.tags != nil {
		ex.tags.off++ // a mask is not content
	}
	ex.floors = append(ex.floors, ex.floor, ex.unsaved)
	ex.floor, ex.unsaved = 0, 0
	// the mask is drawn in the space of the group: with the transform of
	// the group where it begins
	if parent.st.ctm != identity {
		ex.concat(parent.st.ctm)
	}
}

func (ex *exporter) maskEnd() {
	n := len(ex.masks)
	if n == 0 {
		return
	}
	m := ex.masks[n-1]
	ex.masks = ex.masks[:n-1]
	if ex.tags != nil {
		ex.tags.off--
	}
	if k := len(ex.floors); k >= 2 {
		ex.floor, ex.unsaved = ex.floors[k-2], ex.floors[k-1]
		ex.floors = ex.floors[:k-2]
	}
	ex.s = m.parent
	g := m.group
	if g == nil || len(ex.groups) == 0 || ex.groups[len(ex.groups)-1] != g || !(g.w > 0) || !(g.h > 0) {
		return // a mask outside a group masks nothing
	}
	// the form of the mask is in the space of the page (see stream.base):
	// what was drawn in the space of the group, under the transform of
	// the group and cut to it
	to := g.to
	if !to.finite() {
		return
	}
	box := to.bounds(g.x, g.y, g.x+g.w, g.y+g.h)
	in := m.s
	in.finish()
	var body bytes.Buffer
	fmt.Fprintf(&body, "%s cm %s re W n\n", to.exact(), nums(g.x, g.y, g.w, g.h))
	body.Write(in.b.Bytes())
	kind, extra := "/Alpha", ""
	group := " /Group << /Type /Group /S /Transparency /I true >>"
	if m.kind == bdf.MaskLuminosity {
		kind, extra = "/Luminosity", " /BC ["+rgbOf(m.backdrop)+"]"
		group = " /Group << /Type /Group /S /Transparency /I true /CS /DeviceRGB >>"
	}
	form := ex.pw.addStream(fmt.Sprintf(" /Type /XObject /Subtype /Form /BBox [%s] /Resources %s%s", nums(box[0]-1, box[1]-1, box[2]+1, box[3]+1), in.resources(), group), body.Bytes(), true)
	if len(m.transfer) == 256 {
		tr := ex.pw.addStream(" /FunctionType 0 /Domain [0 1] /Range [0 1] /Size [256] /BitsPerSample 8", m.transfer, true)
		extra += " /TR " + tr.String()
	}
	state := ex.pw.add("<< /Type /ExtGState /SMask << /Type /Mask /S " + kind + " /G " + form.String() + extra + " >> >>")
	g.masks = append(g.masks, softMask{state: state})
}

// link notes a LINK: its rectangle in the space of the page.
func (ex *exporter) link(x, y, w, h float64, url string) {
	if url == "" || len(ex.masks) > 0 {
		return
	}
	s := ex.s
	m := s.toPage.mul(s.st.ctm)
	x0, y0, x1, y1 := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, c := range [][2]float64{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}} {
		px, py := m.apply(c[0], c[1])
		x0, y0, x1, y1 = math.Min(x0, px), math.Min(y0, py), math.Max(x1, px), math.Max(y1, py)
	}
	if x1 > x0 && y1 > y0 && len(s.links) < maxLinks {
		l := link{x0: x0, y0: y0, x1: x1, y1: y1, url: url}
		if ex.tags != nil {
			l.elem = ex.tags.link()
		}
		s.links = append(s.links, l)
	}
}

// maxLinks bounds the links of a stream.
const maxLinks = 1 << 16

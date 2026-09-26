package bdf

import (
	"crypto/sha256"
	"encoding"
	"fmt"
	"hash"
	"sort"
)

// Prefix sharing: several objects that begin with the same instructions
// (pages exported from one master, tiles with the same background) get
// that prefix as one shared child, so the container stores it once and the
// renderer can cache its bitmap.
//
// A cut is legal only at an instruction boundary where the graphics state
// stack is at depth 0 (SAVE/RESTORE, GROUP_BEGIN/GROUP_END and
// MASK_BEGIN/MASK_END balanced), no clip, shadow or filter has been set at
// depth 0 before it, and the previous instruction is not a MARK (an
// ALT_TEXT belongs to the instruction after it). USE wraps the child in an implicit save/restore,
// so the rest of the object starts with the state the prefix established:
// its net transform and the last fill, stroke, line, dash, alpha, blend,
// font, text style and smoothing set at depth 0 are re-emitted first.

// splitPoint is a legal cut: the index of the first instruction of the rest.
type splitPoint struct {
	instr int
	off   int // byte offset in the op stream
}

// instrSpan is one decoded instruction with its byte range.
type instrSpan struct {
	Instr
	start, end int
}

// decodeSpans decodes the builder's op stream with byte offsets.
func (o *Object) decodeSpans() ([]instrSpan, *ObjectPart, error) {
	part, err := DecodeObject(o.Encode())
	if err != nil {
		return nil, nil, err
	}
	var out []instrSpan
	r := &reader{b: part.Ops}
	for !r.eof() {
		start := r.pos
		in, err := part.readInstr(r)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, instrSpan{Instr: in, start: start, end: r.pos})
	}
	return out, part, nil
}

// stateReplay is the depth-0 state a prefix leaves behind.
type stateReplay struct {
	ctm       [6]float32
	hasCTM    bool
	fill      *Instr // last FILL_COLOR or FILL_PAINT
	stroke    *Instr
	line      *Instr
	dash      *Instr
	alpha     *Instr
	blend     *Instr
	font      *Instr
	textStyle *Instr
	smoothing *Instr
	unsafe    bool // a clip, shadow or filter at depth 0: no later cut is legal
}

func (s *stateReplay) apply(in *Instr, depth int) {
	if depth != 0 {
		return
	}
	switch in.Op {
	case OpTransform:
		s.mulCTM([6]float32{f32arg(in, 0), f32arg(in, 1), f32arg(in, 2), f32arg(in, 3), f32arg(in, 4), f32arg(in, 5)})
	case OpTranslate:
		s.mulCTM([6]float32{1, 0, 0, 1, f32arg(in, 0), f32arg(in, 1)})
	case OpScale:
		s.mulCTM([6]float32{f32arg(in, 0), 0, 0, f32arg(in, 1), 0, 0})
	case OpFillColor, OpFillPaint:
		s.fill = in
	case OpStrokeColor, OpStrokePaint:
		s.stroke = in
	case OpLine:
		s.line = in
	case OpDash:
		s.dash = in
	case OpAlpha:
		s.alpha = in
	case OpBlend:
		s.blend = in
	case OpFont:
		s.font = in
	case OpTextStyle:
		s.textStyle = in
	case OpSmoothing:
		s.smoothing = in
	case OpClipPath, OpClipRect, OpShadow, OpFilter:
		s.unsafe = true
	}
}

func (s *stateReplay) mulCTM(m [6]float32) {
	if !s.hasCTM {
		s.ctm = [6]float32{1, 0, 0, 1, 0, 0}
		s.hasCTM = true
	}
	c := s.ctm
	s.ctm = [6]float32{
		c[0]*m[0] + c[2]*m[1], c[1]*m[0] + c[3]*m[1],
		c[0]*m[2] + c[2]*m[3], c[1]*m[2] + c[3]*m[3],
		c[0]*m[4] + c[2]*m[5] + c[4], c[1]*m[4] + c[3]*m[5] + c[5],
	}
}

// instrs returns the instructions that re-establish the state.
func (s *stateReplay) instrs() []Instr {
	var out []Instr
	if s.hasCTM && s.ctm != [6]float32{1, 0, 0, 1, 0, 0} {
		out = append(out, Instr{Op: OpTransform, Args: []any{s.ctm[0], s.ctm[1], s.ctm[2], s.ctm[3], s.ctm[4], s.ctm[5]}})
	}
	for _, in := range []*Instr{s.fill, s.stroke, s.line, s.dash, s.alpha, s.blend, s.font, s.textStyle, s.smoothing} {
		if in != nil {
			out = append(out, *in)
		}
	}
	return out
}

func f32arg(in *Instr, i int) float32 { return in.Args[i].(float32) }
func uarg(in *Instr, i int) uint64    { return in.Args[i].(uint64) }

// refMap renumbers the resources an instruction stream references into a
// destination object, adding each on first use.
type refMap struct {
	src     *Object
	part    *ObjectPart
	dst     *Object
	paths   map[uint64]PathRef
	paints  map[uint64]PaintRef
	fonts   map[uint64]FontRef
	images  map[uint64]ImageRef
	objects map[uint64]ObjRef
	// Fonts / Objects list, per new index, the source index it came from
	// (-1 for a placeholder the caller fills in).
	Fonts   []FontRef
	Objects []int
}

func newRefMap(src *Object, part *ObjectPart, dst *Object) *refMap {
	return &refMap{src: src, part: part, dst: dst, paths: map[uint64]PathRef{}, paints: map[uint64]PaintRef{}, fonts: map[uint64]FontRef{}, images: map[uint64]ImageRef{}, objects: map[uint64]ObjRef{}}
}

func (m *refMap) path(i uint64) PathRef {
	if r, ok := m.paths[i]; ok {
		return r
	}
	e := m.part.Paths[i]
	var r PathRef
	if e.Inline != nil {
		r = m.dst.AddPath(e.Inline)
	} else {
		r = m.dst.AddExtPath(e.Hash, e.Index)
	}
	m.paths[i] = r
	return r
}

func (m *refMap) image(i uint64) ImageRef {
	if r, ok := m.images[i]; ok {
		return r
	}
	r := m.dst.AddImage(m.part.Images[i])
	m.images[i] = r
	return r
}

func (m *refMap) paint(i uint64) PaintRef {
	if r, ok := m.paints[i]; ok {
		return r
	}
	p := m.part.Paints[i]
	if p.Kind == PaintPattern {
		p.Image = m.image(uint64(p.Image))
	}
	r := m.dst.AddPaint(p)
	m.paints[i] = r
	return r
}

func (m *refMap) font(i uint64) FontRef {
	if r, ok := m.fonts[i]; ok {
		return r
	}
	r := m.dst.AddFont(m.part.Fonts[i])
	m.fonts[i] = r
	m.Fonts = append(m.Fonts, FontRef(i))
	return r
}

func (m *refMap) object(i uint64) ObjRef {
	if r, ok := m.objects[i]; ok {
		return r
	}
	r := m.dst.AddObject(m.part.Objects[i], m.src.objBBox[i])
	m.objects[i] = r
	m.Objects = append(m.Objects, int(i))
	return r
}

// emit replays one instruction into the destination.
func (m *refMap) emit(in Instr) error {
	d := m.dst
	f := func(i int) float32 { return in.Args[i].(float32) }
	u := func(i int) uint64 { return in.Args[i].(uint64) }
	s := func(i int) string { return in.Args[i].(string) }
	switch in.Op {
	case OpSave:
		d.Save()
	case OpRestore:
		d.Restore()
	case OpTransform:
		d.Transform(f(0), f(1), f(2), f(3), f(4), f(5))
	case OpTranslate:
		d.Translate(f(0), f(1))
	case OpScale:
		d.Scale(f(0), f(1))
	case OpClipPath:
		d.ClipPath(m.path(u(0)), byte(u(1)))
	case OpClipRect:
		d.ClipRect(f(0), f(1), f(2), f(3))
	case OpFillColor:
		d.FillColor(Color(u(0)))
	case OpFillPaint:
		d.FillPaint(m.paint(u(0)))
	case OpStrokeColor:
		d.StrokeColor(Color(u(0)))
	case OpStrokePaint:
		d.StrokePaint(m.paint(u(0)))
	case OpLine:
		d.Line(f(0), byte(u(1)), byte(u(2)), f(3))
	case OpDash:
		d.Dash(in.Args[0].([]float32), f(1))
	case OpAlpha:
		d.Alpha(f(0))
	case OpBlend:
		d.Blend(byte(u(0)))
	case OpShadow:
		d.Shadow(Color(u(0)), f(1), f(2), f(3))
	case OpFilter:
		d.Filter(s(0))
	case OpFont:
		d.Font(m.font(u(0)), f(1))
	case OpTextStyle:
		d.TextStyle(byte(u(0)), byte(u(1)), byte(u(2)), f(3))
	case OpFillRect:
		d.FillRect(f(0), f(1), f(2), f(3))
	case OpStrokeRect:
		d.StrokeRect(f(0), f(1), f(2), f(3))
	case OpFillPath:
		d.FillPath(m.path(u(0)), byte(u(1)))
	case OpStrokePath:
		d.StrokePath(m.path(u(0)))
	case OpFillPathAt:
		d.FillPathAt(m.path(u(0)), byte(u(1)), f(2), f(3))
	case OpFillPathRun:
		glyphs := append([]Glyph(nil), in.Args[1].([]Glyph)...)
		for i := range glyphs {
			glyphs[i].Path = m.path(uint64(glyphs[i].Path))
		}
		d.FillPathRun(byte(u(0)), glyphs)
	case OpClearRect:
		d.ClearRect(f(0), f(1), f(2), f(3))
	case OpFillText:
		d.FillText(s(0), f(1), f(2), f(3))
	case OpStrokeText:
		d.StrokeText(s(0), f(1), f(2), f(3))
	case OpImage:
		d.Image(m.image(u(0)), f(1), f(2), f(3), f(4))
	case OpImageSub:
		d.ImageSub(m.image(u(0)), f(1), f(2), f(3), f(4), f(5), f(6), f(7), f(8))
	case OpSmoothing:
		d.Smoothing(u(0) != 0, byte(u(1)))
	case OpUse:
		d.Use(m.object(u(0)))
	case OpUseAt:
		d.UseAt(m.object(u(0)), f(1), f(2))
	case OpGroupBegin:
		d.GroupBegin(f(0), byte(u(1)), f(2), f(3), f(4), f(5))
	case OpGroupEnd:
		d.GroupEnd()
	case OpMaskBegin:
		d.MaskBegin(byte(u(0)), Color(u(1)), in.Args[2].([]byte))
	case OpMaskEnd:
		d.MaskEnd()
	case OpLink:
		d.Link(f(0), f(1), f(2), f(3), s(4))
	case OpMark:
		d.Mark(byte(u(0)), s(1))
	case OpExt:
		d.Ext(in.Args[0].([]byte))
	default:
		return fmt.Errorf("bdf: cannot re-emit opcode 0x%02x", in.Op)
	}
	return nil
}

// Split is an object cut in two: Prefix draws the first instructions,
// Rest begins with USE of it (its object reference 0, to be set by the
// caller once the prefix has a hash), re-establishes the state and draws
// the remaining instructions. Both keep the original bounding box.
//
// PrefixFonts / RestFonts and PrefixObjects / RestObjects map each font and
// child-object reference of the new objects back to the original one, in
// new-reference order (RestObjects[0] is the prefix placeholder, -1).
type Split struct {
	Prefix, Rest               *Object
	PrefixFonts, RestFonts     []FontRef
	PrefixObjects, RestObjects []int
}

// splitAt cuts the object before instruction index i (which must be a
// legal cut; see cutPoints).
func (o *Object) splitAt(spans []instrSpan, part *ObjectPart, i int) (*Split, error) {
	prefix := NewObject().SetBBox(o.BBox.X, o.BBox.Y, o.BBox.W, o.BBox.H)
	rest := NewObject().SetBBox(o.BBox.X, o.BBox.Y, o.BBox.W, o.BBox.H)
	pm := newRefMap(o, part, prefix)
	var state stateReplay
	depth := 0
	for _, sp := range spans[:i] {
		state.apply(&sp.Instr, depth)
		switch sp.Op {
		case OpSave, OpGroupBegin, OpMaskBegin:
			depth++
		case OpRestore, OpGroupEnd, OpMaskEnd:
			depth--
		}
		if err := pm.emit(sp.Instr); err != nil {
			return nil, err
		}
	}
	rm := newRefMap(o, part, rest)
	rm.objects[^uint64(0)] = rest.AddObject(Hash{}, o.BBox)
	rm.Objects = append(rm.Objects, -1)
	rest.Use(0)
	for _, in := range state.instrs() {
		if err := rm.emit(in); err != nil {
			return nil, err
		}
	}
	for _, sp := range spans[i:] {
		if err := rm.emit(sp.Instr); err != nil {
			return nil, err
		}
	}
	return &Split{Prefix: prefix, Rest: rest, PrefixFonts: pm.Fonts, RestFonts: rm.Fonts, PrefixObjects: pm.Objects, RestObjects: rm.Objects}, nil
}

// prefixKeys walks the instructions and returns, for every legal cut, the
// cut and a digest of everything before it: the instruction bytes plus
// the content of each resource they reference, so two objects with equal
// keys draw the same thing up to that point whatever their tables look like.
func (o *Object) prefixKeys(spans []instrSpan, part *ObjectPart) (points []splitPoint, keys [][32]byte) {
	h := sha256.New()
	var state stateReplay
	depth := 0
	feed := func(kind byte, i uint64) {
		var w buf
		w.u8(kind)
		switch kind {
		case 'p':
			e := part.Paths[i]
			if e.Inline != nil {
				e.Inline.encode(&w)
			} else {
				w.hash(e.Hash)
				w.varuint(uint64(e.Index))
			}
		case 'g':
			p := part.Paints[i]
			p.encode(&w)
			if p.Kind == PaintPattern {
				w.hash(part.Images[p.Image])
			}
		case 'f':
			part.Fonts[i].encode(&w)
		case 'i':
			w.hash(part.Images[i])
		case 'o':
			w.hash(part.Objects[i])
		}
		h.Write(w.b)
	}
	for idx, sp := range spans {
		if idx > 0 {
			prev := spans[idx-1]
			if depth == 0 && !state.unsafe && prev.Op != OpMark {
				points = append(points, splitPoint{instr: idx, off: sp.start})
				keys = append(keys, digestSnapshot(h))
			}
		}
		h.Write(part.Ops[sp.start:sp.end])
		in := sp.Instr
		switch in.Op {
		case OpClipPath, OpFillPath, OpFillPathAt:
			feed('p', uarg(&in, 0))
		case OpStrokePath:
			feed('p', uarg(&in, 0))
		case OpFillPathRun:
			for _, g := range in.Args[1].([]Glyph) {
				feed('p', uint64(g.Path))
			}
		case OpFillPaint, OpStrokePaint:
			feed('g', uarg(&in, 0))
		case OpFont:
			feed('f', uarg(&in, 0))
		case OpImage, OpImageSub:
			feed('i', uarg(&in, 0))
		case OpUse, OpUseAt:
			feed('o', uarg(&in, 0))
		}
		state.apply(&in, depth)
		switch in.Op {
		case OpSave, OpGroupBegin, OpMaskBegin:
			depth++
		case OpRestore, OpGroupEnd, OpMaskEnd:
			depth--
		}
	}
	return points, keys
}

// digestSnapshot returns the running digest without disturbing it.
func digestSnapshot(h hash.Hash) [32]byte {
	var out [32]byte
	m, ok := h.(encoding.BinaryMarshaler)
	if !ok {
		copy(out[:], h.Sum(nil))
		return out
	}
	state, err := m.MarshalBinary()
	if err != nil {
		copy(out[:], h.Sum(nil))
		return out
	}
	h2 := sha256.New()
	if err := h2.(encoding.BinaryUnmarshaler).UnmarshalBinary(state); err != nil {
		copy(out[:], h.Sum(nil))
		return out
	}
	copy(out[:], h2.Sum(nil))
	return out
}

// PrefixShare is one shared prefix and the objects that were cut to use it.
type PrefixShare struct {
	Prefix *Object
	// Fonts and child objects of Prefix, mapped back to the references of
	// the object at Members[0] (see Split).
	PrefixFonts   []FontRef
	PrefixObjects []int
	Members       []PrefixMember
}

// PrefixMember is a rewritten object: index into the input slice and the
// rest object, which references the prefix as its child object 0.
type PrefixMember struct {
	Index       int
	Rest        *Object
	RestFonts   []FontRef
	RestObjects []int
	// Bytes of the op stream the member no longer carries.
	Saved int
}

// SharePrefixes finds the longest legal instruction prefix that two or more
// of the objects have in common (at least minBytes of op stream) and
// returns, per distinct prefix, the shared object and the rewritten
// members. Objects that share nothing are absent from the result. The
// objects themselves are not modified.
//
// Equality is by content: instruction bytes plus the paths, paints, fonts,
// images and child hashes they reference. Callers whose resource tables hold
// placeholders (fonts or children filled in later) must make the
// placeholders distinct per resource before calling.
func SharePrefixes(objs []*Object, minBytes int) ([]PrefixShare, error) {
	type cand struct {
		points []splitPoint
		keys   [][32]byte
		spans  []instrSpan
		part   *ObjectPart
	}
	cands := make([]cand, len(objs))
	count := map[[32]byte]int{}
	for i, o := range objs {
		spans, part, err := o.decodeSpans()
		if err != nil {
			return nil, err
		}
		points, keys := o.prefixKeys(spans, part)
		cands[i] = cand{points, keys, spans, part}
		seen := map[[32]byte]bool{}
		for _, k := range keys {
			if !seen[k] {
				seen[k] = true
				count[k]++
			}
		}
	}
	// Each shared key is a candidate group; take the groups greedily by the
	// bytes they save, each object joining at most one group. A three-page
	// master at a shorter cut beats two pages sharing a slightly longer one.
	type cut struct {
		obj   int
		point splitPoint
	}
	byKey := map[[32]byte][]cut{}
	var keys [][32]byte
	for i, c := range cands {
		for j, k := range c.keys {
			if c.points[j].off < minBytes || count[k] < 2 {
				continue
			}
			if _, ok := byKey[k]; !ok {
				keys = append(keys, k)
			}
			byKey[k] = append(byKey[k], cut{i, c.points[j]})
		}
	}
	taken := make([]bool, len(objs))
	var groups [][]cut
	for {
		var best []cut
		bestSaved := 0
		for _, k := range keys {
			var free []cut
			for _, c := range byKey[k] {
				if !taken[c.obj] {
					free = append(free, c)
				}
			}
			if len(free) < 2 {
				continue
			}
			if saved := (len(free) - 1) * free[0].point.off; saved > bestSaved {
				best, bestSaved = free, saved
			}
		}
		if best == nil {
			break
		}
		for _, c := range best {
			taken[c.obj] = true
		}
		groups = append(groups, best)
	}
	sort.SliceStable(groups, func(a, b int) bool { return groups[a][0].obj < groups[b][0].obj })
	var out []PrefixShare
	for _, members := range groups {
		share := PrefixShare{}
		for n, c := range members {
			i := c.obj
			sp, err := objs[i].splitAt(cands[i].spans, cands[i].part, c.point.instr)
			if err != nil {
				return nil, err
			}
			if n == 0 {
				share.Prefix, share.PrefixFonts, share.PrefixObjects = sp.Prefix, sp.PrefixFonts, sp.PrefixObjects
			}
			share.Members = append(share.Members, PrefixMember{Index: i, Rest: sp.Rest, RestFonts: sp.RestFonts, RestObjects: sp.RestObjects, Saved: c.point.off})
		}
		out = append(out, share)
	}
	return out, nil
}

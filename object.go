package bdf

// Local resource references inside an object.
type (
	PathRef  uint32
	PaintRef uint32
	FontRef  uint32
	ImageRef uint32
	ObjRef   uint32
)

// Rect is an axis-aligned rectangle.
type Rect struct {
	X, Y, W, H float32
}

// pathEntry is either an inline path or a reference into a Path collection part.
type pathEntry struct {
	Inline *Path
	Hash   Hash
	Index  uint32
}

// Object builds and encodes an Object part: a display list plus its resource tables.
//
// Coordinates are local to the object. The builder tracks the bounding box of
// drawing operations approximately; set BBox explicitly to override.
type Object struct {
	BBox     Rect
	autoBBox bool
	hasBBox  bool

	strings []string
	strIdx  map[string]uint32
	paths   []pathEntry
	paints  []Paint
	fonts   []Font
	images  []Hash
	objects []Hash
	objBBox []Rect

	ops   buf
	depth int
}

// NewObject creates an empty object whose bbox is computed from its content.
func NewObject() *Object {
	return &Object{strIdx: map[string]uint32{}, autoBBox: true}
}

// SetBBox fixes the bounding box instead of computing it.
func (o *Object) SetBBox(x, y, w, h float32) *Object {
	o.BBox = Rect{x, y, w, h}
	o.autoBBox = false
	o.hasBBox = true
	return o
}

func (o *Object) grow(x, y, w, h float32) {
	if !o.autoBBox {
		return
	}
	if w < 0 {
		x, w = x+w, -w
	}
	if h < 0 {
		y, h = y+h, -h
	}
	if !o.hasBBox {
		o.BBox = Rect{x, y, w, h}
		o.hasBBox = true
		return
	}
	x0 := min(o.BBox.X, x)
	y0 := min(o.BBox.Y, y)
	x1 := max(o.BBox.X+o.BBox.W, x+w)
	y1 := max(o.BBox.Y+o.BBox.H, y+h)
	o.BBox = Rect{x0, y0, x1 - x0, y1 - y0}
}

func (o *Object) growPath(p *Path) {
	if !o.autoBBox || p == nil {
		return
	}
	ai := 0
	for _, v := range p.Verbs {
		n := verbArgs[v]
		a := p.Args[ai : ai+n]
		ai += n
		switch v {
		case VerbMove, VerbLine:
			o.grow(a[0], a[1], 0, 0)
		case VerbQuad:
			o.grow(a[0], a[1], 0, 0)
			o.grow(a[2], a[3], 0, 0)
		case VerbCubic:
			o.grow(a[0], a[1], 0, 0)
			o.grow(a[2], a[3], 0, 0)
			o.grow(a[4], a[5], 0, 0)
		case VerbRect, VerbRoundRect:
			o.grow(a[0], a[1], a[2], a[3])
		case VerbEllipse:
			o.grow(a[0]-a[2], a[1]-a[3], 2*a[2], 2*a[3])
		case VerbArcTo:
			o.grow(a[0], a[1], 0, 0)
			o.grow(a[2], a[3], 0, 0)
		}
	}
}

// --- resources ---

func (o *Object) str(s string) uint32 {
	if i, ok := o.strIdx[s]; ok {
		return i
	}
	i := uint32(len(o.strings))
	o.strings = append(o.strings, s)
	o.strIdx[s] = i
	return i
}

// AddPath registers an inline path and returns its reference.
func (o *Object) AddPath(p *Path) PathRef {
	o.paths = append(o.paths, pathEntry{Inline: p})
	return PathRef(len(o.paths) - 1)
}

// AddExtPath registers a reference to path index in a Path collection part.
func (o *Object) AddExtPath(h Hash, index uint32) PathRef {
	o.paths = append(o.paths, pathEntry{Hash: h, Index: index})
	return PathRef(len(o.paths) - 1)
}

// AddPaint registers a gradient or pattern.
func (o *Object) AddPaint(p Paint) PaintRef {
	o.paints = append(o.paints, p)
	return PaintRef(len(o.paints) - 1)
}

// AddFont registers a font.
func (o *Object) AddFont(f Font) FontRef {
	o.fonts = append(o.fonts, f)
	return FontRef(len(o.fonts) - 1)
}

// UpdateFont replaces a registered font (used by converters that finalize
// font programs after the object has been built).
func (o *Object) UpdateFont(ref FontRef, f Font) { o.fonts[ref] = f }

// UpdateObject replaces the hash of a registered child object.
func (o *Object) UpdateObject(ref ObjRef, h Hash) { o.objects[ref] = h }

// AddImage registers an image part reference.
func (o *Object) AddImage(h Hash) ImageRef {
	o.images = append(o.images, h)
	return ImageRef(len(o.images) - 1)
}

// AddObject registers a child object reference. bbox is the child's bounding
// box, used to grow this object's bbox on USE.
func (o *Object) AddObject(h Hash, bbox Rect) ObjRef {
	o.objects = append(o.objects, h)
	o.objBBox = append(o.objBBox, bbox)
	return ObjRef(len(o.objects) - 1)
}

// --- ops ---

func (o *Object) op(code byte) *buf {
	o.ops.u8(code)
	return &o.ops
}

func (o *Object) Save() *Object    { o.op(OpSave); o.depth++; return o }
func (o *Object) Restore() *Object { o.op(OpRestore); o.depth--; return o }
func (o *Object) Transform(a, b, c, d, e, f float32) *Object {
	o.op(OpTransform).f32s(a, b, c, d, e, f)
	return o
}
func (o *Object) Translate(x, y float32) *Object { o.op(OpTranslate).f32s(x, y); return o }
func (o *Object) Scale(x, y float32) *Object     { o.op(OpScale).f32s(x, y); return o }
func (o *Object) ClipPath(p PathRef, rule byte) *Object {
	w := o.op(OpClipPath)
	w.varuint(uint64(p))
	w.u8(rule)
	return o
}
func (o *Object) ClipRect(x, y, w, h float32) *Object { o.op(OpClipRect).f32s(x, y, w, h); return o }

func (o *Object) FillColor(c Color) *Object    { o.op(OpFillColor).u32(uint32(c)); return o }
func (o *Object) FillPaint(p PaintRef) *Object { o.op(OpFillPaint).varuint(uint64(p)); return o }
func (o *Object) StrokeColor(c Color) *Object  { o.op(OpStrokeColor).u32(uint32(c)); return o }
func (o *Object) StrokePaint(p PaintRef) *Object {
	o.op(OpStrokePaint).varuint(uint64(p))
	return o
}
func (o *Object) Line(width float32, cap, join byte, miter float32) *Object {
	w := o.op(OpLine)
	w.f32(width)
	w.u8(cap)
	w.u8(join)
	w.f32(miter)
	return o
}
func (o *Object) Dash(segments []float32, offset float32) *Object {
	w := o.op(OpDash)
	w.varuint(uint64(len(segments)))
	w.f32s(segments...)
	w.f32(offset)
	return o
}
func (o *Object) Alpha(a float32) *Object { o.op(OpAlpha).f32(a); return o }
func (o *Object) Blend(b byte) *Object    { o.op(OpBlend).u8(b); return o }
func (o *Object) Shadow(c Color, blur, dx, dy float32) *Object {
	w := o.op(OpShadow)
	w.u32(uint32(c))
	w.f32s(blur, dx, dy)
	return o
}
func (o *Object) Filter(css string) *Object { o.op(OpFilter).varuint(uint64(o.str(css))); return o }
func (o *Object) Font(f FontRef, size float32) *Object {
	w := o.op(OpFont)
	w.varuint(uint64(f))
	w.f32(size)
	return o
}
func (o *Object) TextStyle(align, baseline, dir byte, letterSpacing float32) *Object {
	w := o.op(OpTextStyle)
	w.u8(align)
	w.u8(baseline)
	w.u8(dir)
	w.f32(letterSpacing)
	return o
}

func (o *Object) FillRect(x, y, w, h float32) *Object {
	o.grow(x, y, w, h)
	o.op(OpFillRect).f32s(x, y, w, h)
	return o
}
func (o *Object) StrokeRect(x, y, w, h float32) *Object {
	o.grow(x, y, w, h)
	o.op(OpStrokeRect).f32s(x, y, w, h)
	return o
}
func (o *Object) FillPath(p PathRef, rule byte) *Object {
	o.growPath(o.paths[p].Inline)
	w := o.op(OpFillPath)
	w.varuint(uint64(p))
	w.u8(rule)
	return o
}
func (o *Object) StrokePath(p PathRef) *Object {
	o.growPath(o.paths[p].Inline)
	o.op(OpStrokePath).varuint(uint64(p))
	return o
}
func (o *Object) FillPathAt(p PathRef, rule byte, x, y float32) *Object {
	if ip := o.paths[p].Inline; ip != nil && o.autoBBox {
		tmp := &Object{autoBBox: true}
		tmp.growPath(ip)
		o.grow(tmp.BBox.X+x, tmp.BBox.Y+y, tmp.BBox.W, tmp.BBox.H)
	}
	w := o.op(OpFillPathAt)
	w.varuint(uint64(p))
	w.u8(rule)
	w.f32s(x, y)
	return o
}

// Glyph is one placed path in a FILL_PATH_RUN.
type Glyph struct {
	Path PathRef
	X, Y float32
}

func (o *Object) FillPathRun(rule byte, glyphs []Glyph) *Object {
	w := o.op(OpFillPathRun)
	w.u8(rule)
	w.varuint(uint64(len(glyphs)))
	for _, g := range glyphs {
		if ip := o.paths[g.Path].Inline; ip != nil && o.autoBBox {
			tmp := &Object{autoBBox: true}
			tmp.growPath(ip)
			o.grow(tmp.BBox.X+g.X, tmp.BBox.Y+g.Y, tmp.BBox.W, tmp.BBox.H)
		}
		w.varuint(uint64(g.Path))
		w.f32s(g.X, g.Y)
	}
	return o
}
func (o *Object) ClearRect(x, y, w, h float32) *Object { o.op(OpClearRect).f32s(x, y, w, h); return o }

// FillText draws s at (x, y). advance is the expected width in units (0 = no correction).
func (o *Object) FillText(s string, x, y, advance float32) *Object {
	o.grow(x, y, advance, 0)
	w := o.op(OpFillText)
	w.varuint(uint64(o.str(s)))
	w.f32s(x, y, advance)
	return o
}
func (o *Object) StrokeText(s string, x, y, advance float32) *Object {
	o.grow(x, y, advance, 0)
	w := o.op(OpStrokeText)
	w.varuint(uint64(o.str(s)))
	w.f32s(x, y, advance)
	return o
}

func (o *Object) Image(img ImageRef, x, y, w, h float32) *Object {
	o.grow(x, y, w, h)
	b := o.op(OpImage)
	b.varuint(uint64(img))
	b.f32s(x, y, w, h)
	return o
}
func (o *Object) ImageSub(img ImageRef, sx, sy, sw, sh, dx, dy, dw, dh float32) *Object {
	o.grow(dx, dy, dw, dh)
	b := o.op(OpImageSub)
	b.varuint(uint64(img))
	b.f32s(sx, sy, sw, sh, dx, dy, dw, dh)
	return o
}
func (o *Object) Smoothing(enabled bool, quality byte) *Object {
	w := o.op(OpSmoothing)
	if enabled {
		w.u8(1)
	} else {
		w.u8(0)
	}
	w.u8(quality)
	return o
}

func (o *Object) Use(obj ObjRef) *Object {
	bb := o.objBBox[obj]
	o.grow(bb.X, bb.Y, bb.W, bb.H)
	o.op(OpUse).varuint(uint64(obj))
	return o
}
func (o *Object) UseAt(obj ObjRef, x, y float32) *Object {
	bb := o.objBBox[obj]
	o.grow(bb.X+x, bb.Y+y, bb.W, bb.H)
	w := o.op(OpUseAt)
	w.varuint(uint64(obj))
	w.f32s(x, y)
	return o
}
func (o *Object) GroupBegin(alpha float32, blend byte, x, y, w, h float32) *Object {
	b := o.op(OpGroupBegin)
	b.f32(alpha)
	b.u8(blend)
	b.f32s(x, y, w, h)
	return o
}
func (o *Object) GroupEnd() *Object { o.op(OpGroupEnd); return o }

// MaskBegin starts drawing a soft mask for the innermost group (kind is
// MaskAlpha or MaskLuminosity; backdrop is the luminosity backdrop colour;
// transfer is nil or 256 entries mapping each mask value).
func (o *Object) MaskBegin(kind byte, backdrop Color, transfer []byte) *Object {
	b := o.op(OpMaskBegin)
	b.u8(kind)
	b.u32(uint32(backdrop))
	b.varuint(uint64(len(transfer)))
	b.bytes(transfer)
	return o
}

// MaskEnd multiplies the alpha of what the innermost group holds by the mask.
func (o *Object) MaskEnd() *Object { o.op(OpMaskEnd); return o }

func (o *Object) Link(x, y, w, h float32, url string) *Object {
	b := o.op(OpLink)
	b.f32s(x, y, w, h)
	b.varuint(uint64(o.str(url)))
	return o
}
func (o *Object) Mark(kind byte, payload string) *Object {
	b := o.op(OpMark)
	b.u8(kind)
	b.varuint(uint64(o.str(payload)))
	return o
}
func (o *Object) Ext(payload []byte) *Object {
	b := o.op(OpExt)
	b.u32(uint32(len(payload)))
	b.bytes(payload)
	return o
}

// Encode serializes the object as an Object part.
func (o *Object) Encode() []byte {
	var w buf
	w.bytes([]byte("BOBJ"))
	w.u16(OpsetVersion)
	w.u16(0)
	w.f32s(o.BBox.X, o.BBox.Y, o.BBox.W, o.BBox.H)
	w.varuint(uint64(len(o.strings)))
	for _, s := range o.strings {
		w.str(s)
	}
	w.varuint(uint64(len(o.paths)))
	for _, p := range o.paths {
		if p.Inline != nil {
			w.u8(0)
			p.Inline.encode(&w)
		} else {
			w.u8(1)
			w.hash(p.Hash)
			w.varuint(uint64(p.Index))
		}
	}
	w.varuint(uint64(len(o.paints)))
	for i := range o.paints {
		o.paints[i].encode(&w)
	}
	w.varuint(uint64(len(o.fonts)))
	for i := range o.fonts {
		o.fonts[i].encode(&w)
	}
	w.varuint(uint64(len(o.images)))
	for _, h := range o.images {
		w.hash(h)
	}
	w.varuint(uint64(len(o.objects)))
	for _, h := range o.objects {
		w.hash(h)
	}
	w.varuint(uint64(len(o.ops.b)))
	w.bytes(o.ops.b)
	return w.b
}

// Deps returns the hashes of parts this object references (fonts, images, paths, objects).
func (o *Object) Deps() []Hash {
	var out []Hash
	for _, f := range o.fonts {
		if f.Kind == FontEmbedded {
			out = append(out, f.Hash)
		}
	}
	out = append(out, o.images...)
	for _, p := range o.paths {
		if p.Inline == nil {
			out = append(out, p.Hash)
		}
	}
	return append(out, o.objects...)
}

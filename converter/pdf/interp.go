package pdf

import (
	"math"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

// pending is a BDF object under construction whose font and child hashes are
// patched in once every font program has been finalized.
type pending struct {
	obj      *bdf.Object
	fonts    map[bdf.FontRef]*pdfFont
	fontRefs map[*pdfFont]bdf.FontRef
	children map[bdf.ObjRef]*pending
	hash     bdf.Hash
	encoded  bool
	building bool
	bbox     bdf.Rect
	st       *formStruct // structure state a form was converted under
}

func newPending(bbox bdf.Rect) *pending {
	p := &pending{obj: bdf.NewObject(), fonts: map[bdf.FontRef]*pdfFont{}, fontRefs: map[*pdfFont]bdf.FontRef{}, children: map[bdf.ObjRef]*pending{}, bbox: bbox}
	p.obj.SetBBox(bbox.X, bbox.Y, bbox.W, bbox.H)
	return p
}

func (p *pending) fontRef(f *pdfFont) bdf.FontRef {
	if r, ok := p.fontRefs[f]; ok {
		return r
	}
	r := p.obj.AddFont(bdf.SystemFont(f.family, f.weight, f.style))
	p.fontRefs[f] = r
	p.fonts[r] = f
	return r
}

func (p *pending) childRef(child *pending) bdf.ObjRef {
	for r, c := range p.children {
		if c == child {
			return r
		}
	}
	r := p.obj.AddObject(bdf.Hash{}, child.bbox)
	p.children[r] = child
	return r
}

// gstate mirrors the PDF graphics state that affects emitted instructions.
type gstate struct {
	ctm  matrix // current user space → object space
	clip rect   // object space; nil-equivalent when infinite

	fillCS, strokeCS       *colorSpace
	fillComps, strokeComps []float64
	fillPattern            types.Object
	strokePattern          types.Object
	fillColor, strokeColor bdf.Color

	lineWidth   float64
	lineCap     byte
	lineJoin    byte
	miterLimit  float64
	dash        []float64
	dashPhase   float64
	fillAlpha   float64
	strokeAlpha float64
	blend       byte
	softMask    *softMask

	font    *pdfFont
	size    float64
	charSp  float64
	wordSp  float64
	hscale  float64
	leading float64
	rise    float64
	render  int
}

// emitted tracks what the object's instruction stream currently has set, so
// state is only written when it changes.
type emitted struct {
	fillColor   bdf.Color
	fillIsPaint bool
	strokeColor bdf.Color
	strokeIsPnt bool
	lineWidth   float64
	lineCap     byte
	lineJoin    byte
	miterLimit  float64
	dashKey     string
	alpha       float64
	blend       byte
	font        *pdfFont
	fontSize    float64
	letterSp    float64
	smoothing   bool
	fillSet     bool // FILL_COLOR has been written at least once
	strokeSet   bool
	lineSet     bool
}

func initialEmitted() emitted {
	return emitted{fillColor: bdf.RGB(0, 0, 0), strokeColor: bdf.RGB(0, 0, 0), lineWidth: 1, miterLimit: 10, alpha: 1, smoothing: true}
}

var infRect = rect{math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1)}

func initialGState(ctm matrix) gstate {
	return gstate{ctm: ctm, clip: infRect, fillCS: csGray, strokeCS: csGray, fillComps: []float64{0}, strokeComps: []float64{0},
		fillColor: bdf.RGB(0, 0, 0), strokeColor: bdf.RGB(0, 0, 0), lineWidth: 1, miterLimit: 10, fillAlpha: 1, strokeAlpha: 1, hscale: 1}
}

// textState is the text object state (BT … ET).
type textState struct {
	tm, tlm  matrix
	base     matrix  // origin of tx: the line matrix at which the current block/run started
	tx       float64 // displacement along the baseline since base was set (text space)
	ty       float64 // vertical writing: displacement along the text space y axis
	open     bool    // a SAVE/TRANSFORM block for the current line matrix is open
	openTlm  matrix
	openH    float64
	openRise float64
	openFlip bool
}

// actualText is an /ActualText replacement in effect for a marked-content sequence.
type actualText struct {
	text string
	used bool
}

// interp executes one content stream into a pending object.
type interp struct {
	c        *converter
	p        *pending
	obj      *bdf.Object
	res      types.Dict
	gs       gstate
	stack    []gstate
	em       emitted
	emStack  []emitted
	ctmStack [][2]any
	base     matrix // initial ctm (pattern space reference)
	depth    int
	masks    []maskGroup // groups open for soft masks (softmask.go)

	path       *bdf.Path
	pathBBox   rect
	pathEmpty  bool
	cur, start [2]float64
	pendClip   int // -1 none, 0 nonzero, 1 evenodd

	text       textState
	curRun     textRun
	vrun       vertRun // vertical text being accumulated (vertical.go)
	actual     *actualText
	mcStack    []mcEntry
	inText     bool
	compat     int
	type3Glyph bool

	// Structure (structure.go). st is nil where no structure or language
	// MARKs are emitted (glyph procedures, pattern cells, annotations).
	st         *structState
	mcids      *mcidTable
	inherit    *mcTarget // target of content outside marked content (the caller's for a form)
	structUsed bool
}

func (c *converter) newInterp(p *pending, res types.Dict, base matrix, depth int) *interp {
	in := &interp{c: c, p: p, obj: p.obj, res: res, gs: initialGState(base), em: initialEmitted(), base: base, depth: depth, pendClip: -1, curRun: textRun{empty: true}}
	in.resetPath()
	return in
}

func (in *interp) resetPath() {
	in.path = &bdf.Path{}
	in.pathBBox = rect{}
	in.pathEmpty = true
	in.pendClip = -1
}

// --- state emission helpers ---

func (in *interp) save() {
	in.stack = append(in.stack, in.gs)
	in.emStack = append(in.emStack, in.em)
	in.obj.Save()
	in.gs.dash = append([]float64(nil), in.gs.dash...)
	in.gs.fillComps = append([]float64(nil), in.gs.fillComps...)
	in.gs.strokeComps = append([]float64(nil), in.gs.strokeComps...)
}

func (in *interp) restore() {
	if len(in.stack) == 0 {
		return
	}
	for n := len(in.masks); n > 0 && in.masks[n-1].depth >= len(in.stack); n = len(in.masks) {
		in.closeMask(true)
	}
	in.gs = in.stack[len(in.stack)-1]
	in.stack = in.stack[:len(in.stack)-1]
	in.em = in.emStack[len(in.emStack)-1]
	in.emStack = in.emStack[:len(in.emStack)-1]
	in.obj.Restore()
}

// textSave/textRestore bracket a text block: they save the emitted state and
// the transform, but leave the PDF graphics state (which outlives the block) alone.
func (in *interp) textSave() {
	in.emStack = append(in.emStack, in.em)
	in.ctmStack = append(in.ctmStack, [2]any{in.gs.ctm, in.gs.clip})
	in.obj.Save()
}

func (in *interp) textRestore() {
	if len(in.ctmStack) == 0 {
		return
	}
	top := in.ctmStack[len(in.ctmStack)-1]
	in.ctmStack = in.ctmStack[:len(in.ctmStack)-1]
	in.gs.ctm = top[0].(matrix)
	in.gs.clip = top[1].(rect)
	in.em = in.emStack[len(in.emStack)-1]
	in.emStack = in.emStack[:len(in.emStack)-1]
	in.obj.Restore()
}

func (in *interp) transform(m matrix) {
	in.gs.ctm = m.mul(in.gs.ctm)
	in.obj.Transform(float32(m[0]), float32(m[1]), float32(m[2]), float32(m[3]), float32(m[4]), float32(m[5]))
}

func (in *interp) syncFill() {
	if !in.em.fillSet || in.em.fillIsPaint || in.em.fillColor != in.gs.fillColor {
		in.obj.FillColor(in.gs.fillColor)
		in.em.fillColor = in.gs.fillColor
		in.em.fillIsPaint = false
		in.em.fillSet = true
	}
}

func (in *interp) syncStroke() {
	if !in.em.strokeSet || in.em.strokeIsPnt || in.em.strokeColor != in.gs.strokeColor {
		in.obj.StrokeColor(in.gs.strokeColor)
		in.em.strokeColor = in.gs.strokeColor
		in.em.strokeIsPnt = false
		in.em.strokeSet = true
	}
}

func (in *interp) syncLine() {
	lw := in.gs.lineWidth
	if lw <= 0 {
		// Zero-width lines are the thinnest renderable line; approximate one device pixel.
		lw = 0.8 / math.Max(in.gs.ctm.scaleFactor(), 1e-6)
	}
	if !in.em.lineSet || in.em.lineWidth != lw || in.em.lineCap != in.gs.lineCap || in.em.lineJoin != in.gs.lineJoin || in.em.miterLimit != in.gs.miterLimit {
		in.obj.Line(float32(lw), in.gs.lineCap, in.gs.lineJoin, float32(in.gs.miterLimit))
		in.em.lineWidth, in.em.lineCap, in.em.lineJoin, in.em.miterLimit = lw, in.gs.lineCap, in.gs.lineJoin, in.gs.miterLimit
		in.em.lineSet = true
	}
	key := dashKey(in.gs.dash, in.gs.dashPhase)
	if key != in.em.dashKey {
		segs := make([]float32, 0, len(in.gs.dash))
		allZero := true
		for _, d := range in.gs.dash {
			segs = append(segs, float32(d))
			if d > 0 {
				allZero = false
			}
		}
		if allZero {
			segs = nil
		}
		in.obj.Dash(segs, float32(in.gs.dashPhase))
		in.em.dashKey = key
	}
}

func dashKey(d []float64, phase float64) string {
	if len(d) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, v := range d {
		sb.WriteString(ftoa(v))
		sb.WriteByte(',')
	}
	sb.WriteString(ftoa(phase))
	return sb.String()
}

func ftoa(v float64) string { return strings.TrimRight(strings.TrimRight(strconvFormat(v), "0"), ".") }

func (in *interp) syncAlpha(stroke bool) {
	a := in.gs.fillAlpha
	if stroke {
		a = in.gs.strokeAlpha
	}
	if in.em.alpha != a {
		in.obj.Alpha(float32(a))
		in.em.alpha = a
	}
	if in.em.blend != in.gs.blend {
		in.obj.Blend(in.gs.blend)
		in.em.blend = in.gs.blend
	}
}

// --- path construction ---

func (in *interp) moveTo(x, y float64) {
	in.path.MoveTo(float32(x), float32(y))
	in.cur, in.start = [2]float64{x, y}, [2]float64{x, y}
	in.growPath(x, y)
}

func (in *interp) lineTo(x, y float64) {
	if in.pathEmpty {
		in.moveTo(x, y)
		return
	}
	in.path.LineTo(float32(x), float32(y))
	in.cur = [2]float64{x, y}
	in.growPath(x, y)
}

func (in *interp) curveTo(x1, y1, x2, y2, x3, y3 float64) {
	if in.pathEmpty {
		in.moveTo(x1, y1)
	}
	in.path.CubicTo(float32(x1), float32(y1), float32(x2), float32(y2), float32(x3), float32(y3))
	in.growPath(x1, y1)
	in.growPath(x2, y2)
	in.growPath(x3, y3)
	in.cur = [2]float64{x3, y3}
}

func (in *interp) closePath() {
	if !in.pathEmpty {
		in.path.Close()
		in.cur = in.start
	}
}

func (in *interp) rectPath(x, y, w, h float64) {
	in.path.Rect(float32(x), float32(y), float32(w), float32(h))
	in.growPath(x, y)
	in.growPath(x+w, y+h)
	in.cur, in.start = [2]float64{x, y}, [2]float64{x, y}
}

func (in *interp) growPath(x, y float64) {
	if in.pathEmpty {
		in.pathBBox = rect{x, y, x, y}
		in.pathEmpty = false
		return
	}
	in.pathBBox.x0 = math.Min(in.pathBBox.x0, x)
	in.pathBBox.y0 = math.Min(in.pathBBox.y0, y)
	in.pathBBox.x1 = math.Max(in.pathBBox.x1, x)
	in.pathBBox.y1 = math.Max(in.pathBBox.y1, y)
}

// endPath paints the current path as requested and applies a pending clip.
func (in *interp) endPath(fill bool, rule byte, stroke bool) {
	var ref bdf.PathRef
	added := false
	addPath := func() bdf.PathRef {
		if !added {
			ref = in.obj.AddPath(in.path)
			added = true
		}
		return ref
	}
	if !in.pathEmpty {
		if fill || stroke {
			in.syncDraw(in.curTarget())
			in.enterMask()
		}
		if fill {
			if in.gs.fillCS != nil && in.gs.fillCS.family == "Pattern" {
				in.patternFill(addPath(), rule)
			} else {
				in.syncFill()
				in.syncAlpha(false)
				in.obj.FillPath(addPath(), rule)
			}
		}
		if stroke {
			if in.gs.strokeCS != nil && in.gs.strokeCS.family == "Pattern" {
				// Stroke with a pattern: approximate with the pattern's average colour.
				in.gs.strokeColor = in.patternAverage(in.gs.strokePattern)
			}
			in.syncStroke()
			in.syncLine()
			in.syncAlpha(true)
			in.obj.StrokePath(addPath())
		}
		if in.pendClip >= 0 {
			in.obj.ClipPath(addPath(), byte(in.pendClip))
			in.gs.clip = in.gs.clip.intersect(in.gs.ctm.transformRect(in.pathBBox))
		}
	} else if in.pendClip >= 0 {
		// Clipping with an empty path clips everything.
		in.obj.ClipRect(0, 0, 0, 0)
		in.gs.clip = rect{}
	}
	in.resetPath()
}

// --- colour ---

func (in *interp) setColor(stroke bool, comps []float64, patName string) {
	gs := &in.gs
	cs := gs.fillCS
	if stroke {
		cs = gs.strokeCS
	}
	if cs == nil {
		cs = csGray
	}
	var col bdf.Color
	var pat types.Object
	if cs.family == "Pattern" {
		if patName != "" {
			if pd := in.c.pdf.dict(in.res["Pattern"]); pd != nil {
				pat = pd[patName]
			}
		}
		col = in.patternAverage(pat)
		if cs.under != nil && len(comps) >= cs.under.n {
			col = cs.under.rgb(comps)
		}
	} else {
		if len(comps) < cs.n {
			comps = append(comps, cs.initial()[len(comps):]...)
		}
		col = cs.rgb(comps)
	}
	if stroke {
		gs.strokeComps, gs.strokeColor, gs.strokePattern = comps, col, pat
	} else {
		gs.fillComps, gs.fillColor, gs.fillPattern = comps, col, pat
	}
}

func (in *interp) setColorSpace(stroke bool, cs *colorSpace) {
	if stroke {
		in.gs.strokeCS = cs
		in.gs.strokeComps = cs.initial()
		in.gs.strokeColor = cs.rgb(in.gs.strokeComps)
		in.gs.strokePattern = nil
	} else {
		in.gs.fillCS = cs
		in.gs.fillComps = cs.initial()
		in.gs.fillColor = cs.rgb(in.gs.fillComps)
		in.gs.fillPattern = nil
	}
}

// --- ExtGState ---

func (in *interp) applyExtGState(name string) {
	p := in.c.pdf
	egs := p.dict(in.res["ExtGState"])
	if egs == nil {
		return
	}
	d := p.dict(egs[name])
	if d == nil {
		return
	}
	for k, v := range d {
		switch k {
		case "LW":
			in.gs.lineWidth = p.numOr(v, in.gs.lineWidth)
		case "LC":
			in.gs.lineCap = byte(p.intOr(v, 0))
		case "LJ":
			in.gs.lineJoin = byte(p.intOr(v, 0))
		case "ML":
			in.gs.miterLimit = p.numOr(v, 10)
		case "D":
			if a := p.array(v); len(a) == 2 {
				in.gs.dash = p.nums(a[0])
				in.gs.dashPhase = p.numOr(a[1], 0)
			}
		case "CA":
			in.gs.strokeAlpha = clamp(p.numOr(v, 1), 0, 1)
		case "ca":
			in.gs.fillAlpha = clamp(p.numOr(v, 1), 0, 1)
		case "BM":
			bm := p.name(v)
			if bm == "" {
				if a := p.array(v); len(a) > 0 {
					bm = p.name(a[0])
				}
			}
			in.gs.blend = blendMode(bm)
		case "Font":
			if a := p.array(v); len(a) == 2 {
				if f := in.c.loadFont(a[0], in.res); f != nil {
					in.gs.font = f
					in.gs.size = p.numOr(a[1], in.gs.size)
				}
			}
		case "SMask":
			if p.name(v) == "None" || p.deref(v) == nil {
				in.setSoftMask(nil)
			} else {
				in.setSoftMask(in.loadSoftMask(v))
			}
		}
	}
}

func blendMode(name string) byte {
	switch name {
	case "Multiply":
		return bdf.BlendMultiply
	case "Screen":
		return bdf.BlendScreen
	case "Overlay":
		return bdf.BlendOverlay
	case "Darken":
		return bdf.BlendDarken
	case "Lighten":
		return bdf.BlendLighten
	case "ColorDodge":
		return bdf.BlendColorDodge
	case "ColorBurn":
		return bdf.BlendColorBurn
	case "HardLight":
		return bdf.BlendHardLight
	case "SoftLight":
		return bdf.BlendSoftLight
	case "Difference":
		return bdf.BlendDifference
	case "Exclusion":
		return bdf.BlendExclusion
	case "Hue":
		return bdf.BlendHue
	case "Saturation":
		return bdf.BlendSaturation
	case "Color":
		return bdf.BlendColor
	case "Luminosity":
		return bdf.BlendLuminosity
	}
	return bdf.BlendSourceOver
}

// --- main loop ---

// run executes a content stream.
func (in *interp) run(content []byte) {
	l := &lexer{b: content}
	var args []types.Object
	for {
		kind, obj, op := l.next()
		if kind == tokEOF {
			break
		}
		if kind == tokOperand {
			args = append(args, obj)
			if len(args) > 64 {
				args = args[len(args)-64:]
			}
			continue
		}
		if op == "BI" {
			in.inlineImage(l)
			args = args[:0]
			continue
		}
		in.exec(op, args)
		args = args[:0]
	}
	in.endText()
	for len(in.stack) > 0 {
		in.restore()
	}
	for len(in.masks) > 0 {
		in.closeMask(false)
	}
	in.gs.softMask = nil
}

func fnum(p *pdf, args []types.Object, i int) float64 {
	if i < 0 || i >= len(args) {
		return 0
	}
	v, _ := p.num(args[i])
	return v
}

func (in *interp) exec(op string, args []types.Object) {
	p := in.c.pdf
	n := len(args)
	f := func(i int) float64 { return fnum(p, args, i) }
	switch op {
	// graphics state
	case "q":
		in.save()
	case "Q":
		in.restore()
	case "cm":
		if n >= 6 {
			in.transform(matrix{f(n - 6), f(n - 5), f(n - 4), f(n - 3), f(n - 2), f(n - 1)})
		}
	case "w":
		in.gs.lineWidth = f(n - 1)
	case "J":
		in.gs.lineCap = byte(clamp(f(n-1), 0, 2))
	case "j":
		in.gs.lineJoin = byte(clamp(f(n-1), 0, 2))
	case "M":
		in.gs.miterLimit = f(n - 1)
	case "d":
		if n >= 2 {
			in.gs.dash = p.nums(args[n-2])
			in.gs.dashPhase = f(n - 1)
		}
	case "gs":
		in.flushRun()
		if n >= 1 {
			in.applyExtGState(p.name(args[n-1]))
		}
	case "ri", "i":
	// path construction
	case "m":
		if n >= 2 {
			in.moveTo(f(n-2), f(n-1))
		}
	case "l":
		if n >= 2 {
			in.lineTo(f(n-2), f(n-1))
		}
	case "c":
		if n >= 6 {
			in.curveTo(f(n-6), f(n-5), f(n-4), f(n-3), f(n-2), f(n-1))
		}
	case "v":
		if n >= 4 {
			in.curveTo(in.cur[0], in.cur[1], f(n-4), f(n-3), f(n-2), f(n-1))
		}
	case "y":
		if n >= 4 {
			in.curveTo(f(n-4), f(n-3), f(n-2), f(n-1), f(n-2), f(n-1))
		}
	case "h":
		in.closePath()
	case "re":
		if n >= 4 {
			in.rectPath(f(n-4), f(n-3), f(n-2), f(n-1))
		}
	// path painting
	case "S":
		in.endPath(false, 0, true)
	case "s":
		in.closePath()
		in.endPath(false, 0, true)
	case "f", "F":
		in.endPath(true, bdf.NonZero, false)
	case "f*":
		in.endPath(true, bdf.EvenOdd, false)
	case "B":
		in.endPath(true, bdf.NonZero, true)
	case "B*":
		in.endPath(true, bdf.EvenOdd, true)
	case "b":
		in.closePath()
		in.endPath(true, bdf.NonZero, true)
	case "b*":
		in.closePath()
		in.endPath(true, bdf.EvenOdd, true)
	case "n":
		in.endPath(false, 0, false)
	case "W":
		in.pendClip = int(bdf.NonZero)
	case "W*":
		in.pendClip = int(bdf.EvenOdd)
	// colour
	case "CS", "cs":
		if n >= 1 {
			in.setColorSpace(op == "CS", in.c.loadColorSpace(args[n-1], in.res))
		}
	case "SC", "SCN", "sc", "scn":
		in.flushRun()
		stroke := op == "SC" || op == "SCN"
		var comps []float64
		pat := ""
		for _, a := range args {
			if v, ok := p.num(a); ok {
				comps = append(comps, v)
			} else if nm, ok := a.(types.Name); ok {
				pat = nm.Value()
			}
		}
		in.setColor(stroke, comps, pat)
	case "G", "g":
		in.flushRun()
		in.setColorSpace(op == "G", csGray)
		in.setColor(op == "G", []float64{f(n - 1)}, "")
	case "RG", "rg":
		in.flushRun()
		if n >= 3 {
			in.setColorSpace(op == "RG", csRGB)
			in.setColor(op == "RG", []float64{f(n - 3), f(n - 2), f(n - 1)}, "")
		}
	case "K", "k":
		in.flushRun()
		if n >= 4 {
			in.setColorSpace(op == "K", csCMYK)
			in.setColor(op == "K", []float64{f(n - 4), f(n - 3), f(n - 2), f(n - 1)}, "")
		}
	// text
	case "BT":
		in.beginText()
	case "ET":
		in.endText()
	case "Tc":
		in.flushRun()
		in.gs.charSp = f(n - 1)
	case "Tw":
		in.flushRun()
		in.gs.wordSp = f(n - 1)
	case "Tz":
		in.flushRun()
		in.gs.hscale = f(n-1) / 100
	case "TL":
		in.gs.leading = f(n - 1)
	case "Ts":
		in.flushRun()
		in.gs.rise = f(n - 1)
	case "Tr":
		in.flushRun()
		in.gs.render = int(f(n - 1))
	case "Tf":
		in.flushRun()
		if n >= 2 {
			in.gs.size = f(n - 1)
			in.gs.font = nil
			if fd := p.dict(in.res["Font"]); fd != nil {
				if ref, ok := fd[p.name(args[n-2])]; ok {
					in.gs.font = in.c.loadFont(ref, in.res)
				}
			}
			if in.gs.font == nil {
				in.c.warnOnce("font-missing", "font resource %s not found; using a default font", p.name(args[n-2]))
				in.gs.font = in.c.defaultFont()
			}
		}
	case "Td":
		if n >= 2 {
			in.setLineMatrix(matrix{1, 0, 0, 1, f(n - 2), f(n - 1)}.mul(in.text.tlm))
		}
	case "TD":
		if n >= 2 {
			in.gs.leading = -f(n - 1)
			in.setLineMatrix(matrix{1, 0, 0, 1, f(n - 2), f(n - 1)}.mul(in.text.tlm))
		}
	case "Tm":
		if n >= 6 {
			in.setLineMatrix(matrix{f(n - 6), f(n - 5), f(n - 4), f(n - 3), f(n - 2), f(n - 1)})
		}
	case "T*":
		in.setLineMatrix(matrix{1, 0, 0, 1, 0, -in.gs.leading}.mul(in.text.tlm))
	case "Tj":
		if n >= 1 {
			in.showText(literalBytes(args[n-1]))
		}
	case "'":
		if n >= 1 {
			in.setLineMatrix(matrix{1, 0, 0, 1, 0, -in.gs.leading}.mul(in.text.tlm))
			in.showText(literalBytes(args[n-1]))
		}
	case "\"":
		if n >= 3 {
			in.gs.wordSp = f(n - 3)
			in.gs.charSp = f(n - 2)
			in.setLineMatrix(matrix{1, 0, 0, 1, 0, -in.gs.leading}.mul(in.text.tlm))
			in.showText(literalBytes(args[n-1]))
		}
	case "TJ":
		if n >= 1 {
			for _, e := range p.array(args[n-1]) {
				if v, ok := p.num(e); ok {
					in.adjustText(v)
				} else {
					in.showText(literalBytes(e))
				}
			}
		}
	case "d0", "d1":
	// XObjects, images, shadings
	case "Do":
		if n >= 1 {
			in.doXObject(p.name(args[n-1]))
		}
	case "sh":
		if n >= 1 {
			in.doShading(p.name(args[n-1]))
		}
	// marked content
	case "BDC":
		if n >= 2 {
			in.markedContent(p.name(args[n-2]), args[n-1])
		}
	case "BMC":
		if n >= 1 {
			in.markedContent(p.name(args[n-1]), nil)
		}
	case "EMC":
		if n := len(in.mcStack); n > 0 {
			prev := in.curTarget()
			in.mcStack = in.mcStack[:n-1]
			in.actual = nil
			for i := n - 2; i >= 0; i-- {
				if in.mcStack[i].actual != nil {
					in.actual = in.mcStack[i].actual
					break
				}
			}
			if !sameTarget(prev, in.curTarget()) {
				in.flushRun()
			}
		}
	case "MP", "DP":
	case "BX":
		in.compat++
	case "EX":
		in.compat--
	default:
		if in.compat == 0 {
			in.c.warnOnce("op-"+op, "unknown operator %q ignored", op)
		}
	}
}

var blockTags = map[string]bool{"P": true, "H": true, "H1": true, "H2": true, "H3": true, "H4": true, "H5": true, "H6": true,
	"LI": true, "LBody": true, "Lbl": true, "TD": true, "TH": true, "Caption": true, "Div": true, "BlockQuote": true,
	"Note": true, "Figure": true, "TOCI": true, "Title": true, "Formula": true, "Sect": true, "Art": true}

func (in *interp) markedContent(tag string, props types.Object) {
	p := in.c.pdf
	// Inline properties come from the content lexer, which has already
	// unescaped their strings; named ones are resources parsed by pdfcpu.
	d := p.dict(props)
	text := func(o types.Object) string { return decodeText(literalBytes(o)) }
	if nm, ok := props.(types.Name); ok {
		d, text = p.dict(p.dict(in.res["Properties"])[nm.Value()]), p.text
	}
	if in.c.tree == nil && blockTags[tag] {
		// Untagged document: guess paragraphs from the tag names.
		in.flushRun()
		in.obj.Mark(bdf.MarkParagraph, tag)
	}
	e := mcEntry{}
	if v, ok := d["ActualText"]; ok {
		e.actual = &actualText{text: text(v)}
		// The replacement applies to the glyphs that follow, which may span runs.
		in.flushRun()
	}
	if in.st != nil {
		e.tgt, e.set = in.markedTarget(tag, d, text)
		if e.set && !sameTarget(e.tgt, in.curTarget()) {
			in.flushRun()
		}
	}
	in.mcStack = append(in.mcStack, e)
	if e.actual != nil {
		in.actual = e.actual
	}
}

// --- XObjects ---

func (in *interp) doXObject(name string) {
	p := in.c.pdf
	xd := p.dict(in.res["XObject"])
	if xd == nil {
		return
	}
	ref := xd[name]
	sd := p.stream(ref)
	if sd == nil {
		return
	}
	switch p.name(sd.Dict["Subtype"]) {
	case "Image":
		in.drawImage(objKey(ref), sd)
	case "Form":
		in.drawForm(objKey(ref), sd, identity)
	case "PS":
	default:
		in.c.warnOnce("xobj-"+name, "XObject %s has unsupported subtype", name)
	}
}

func (in *interp) drawImage(key string, sd *types.StreamDict) {
	p := in.c.pdf
	isMask := p.boolOr(sd.Dict["ImageMask"], false)
	if key != "" && isMask {
		key += "/" + in.gs.fillColor.CSS()
	}
	var h bdf.Hash
	var ok bool
	var size [2]int
	if key != "" {
		if e, found := in.c.images[key]; found {
			h, ok, size = e.hash, e.ok, e.size
		}
	}
	if !ok && (key == "" || in.c.images[key] == nil) {
		img, err := in.c.loadImage(sd.Dict, sd.Raw, sd.FilterPipeline, in.res, in.gs.fillColor)
		if err != nil {
			in.c.warnf("image: %v", err)
			if key != "" {
				in.c.images[key] = &imageEntry{}
			}
			return
		}
		h = in.c.doc.AddImage(img.data)
		ok = true
		size = [2]int{img.w, img.h}
		if key != "" {
			in.c.images[key] = &imageEntry{hash: h, ok: true, size: size}
		}
	}
	if !ok {
		return
	}
	in.syncDraw(in.targetOf(sd.Dict))
	in.enterMask()
	ref := in.obj.AddImage(h)
	in.syncAlpha(false)
	in.save()
	in.transform(matrix{1, 0, 0, -1, 0, 1})
	// Disable smoothing when a small image is scaled up and does not ask for interpolation.
	if !p.boolOr(sd.Dict["Interpolate"], false) {
		dev := in.gs.ctm.transformRect(rect{0, 0, 1, 1})
		if size[0] > 0 && size[1] > 0 && (dev.x1-dev.x0 > float64(size[0])*2 || dev.y1-dev.y0 > float64(size[1])*2) {
			in.obj.Smoothing(false, 0)
		}
	}
	in.obj.Image(ref, 0, 0, 1, 1)
	in.restore()
}

// drawForm draws a form XObject with an extra matrix applied first (annotations).
func (in *interp) drawForm(key string, sd *types.StreamDict, extra matrix) {
	p := in.c.pdf
	if in.depth > 12 {
		in.c.warnOnce("form-depth", "form XObjects nested too deeply; skipping")
		return
	}
	var fs *formStruct
	if in.st != nil {
		tgt := in.targetOf(sd.Dict)
		in.syncDraw(tgt)
		// The form continues the caller's walk: it starts in the caller's
		// state and leaves its own for what follows the USE.
		fs = &formStruct{mcids: in.mcids, inherit: tgt, entry: in.st.clone()}
		if sp, ok := p.num(sd.Dict["StructParents"]); ok && in.c.tree != nil {
			fs.mcids = in.c.tree.contentTable(int(sp), true, 0)
		}
	}
	child := in.c.formObject(key, sd, in.res, in.depth+1, fs)
	if child == nil {
		return
	}
	if fs != nil && child.st != nil && child.st.used {
		*in.st = child.st.exit.clone()
		in.structUsed = true
	}
	mat := p.matrixOr(sd.Dict["Matrix"], identity)
	bbox := p.rectOr(sd.Dict["BBox"], rect{})
	group := p.dict(sd.Dict["Group"]) != nil && (in.gs.fillAlpha < 1 || in.gs.blend != bdf.BlendSourceOver)
	in.enterMask()
	in.syncAlpha(false)
	in.save()
	if !extra.isIdentity() {
		in.transform(extra)
	}
	if !mat.isIdentity() {
		in.transform(mat)
	}
	if !bbox.empty() {
		in.obj.ClipRect(float32(bbox.x0), float32(bbox.y0), float32(bbox.x1-bbox.x0), float32(bbox.y1-bbox.y0))
		in.gs.clip = in.gs.clip.intersect(in.gs.ctm.transformRect(bbox))
	}
	ref := in.p.childRef(child)
	if group {
		in.obj.GroupBegin(float32(in.gs.fillAlpha), in.gs.blend, float32(bbox.x0), float32(bbox.y0), float32(bbox.x1-bbox.x0), float32(bbox.y1-bbox.y0))
		in.obj.Use(ref)
		in.obj.GroupEnd()
	} else {
		in.obj.Use(ref)
	}
	in.restore()
}

// --- shadings and patterns ---

func (in *interp) doShading(name string) {
	p := in.c.pdf
	shd := p.dict(in.res["Shading"])
	if shd == nil {
		return
	}
	sh := in.c.shadingFor(objKey(shd[name]), shd[name], in.res)
	if sh == nil {
		return
	}
	// Paint the clip region (in current user space).
	inv, ok := in.gs.ctm.inverse()
	if !ok {
		return
	}
	area := in.gs.clip
	if math.IsInf(area.x0, 0) || math.IsInf(area.y0, 0) {
		area = in.c.pageBox
	}
	r := inv.transformRect(area)
	if sh.bbox != nil {
		r = r.intersect(*sh.bbox)
	}
	if r.empty() {
		return
	}
	in.syncDraw(in.curTarget())
	in.enterMask()
	in.syncAlpha(false)
	in.paintShading(sh, r)
}

func (in *interp) paintShading(sh *shading, r rect) {
	if sh.hasAvg {
		in.obj.FillColor(sh.avg)
		in.em.fillColor = sh.avg
		in.em.fillIsPaint = false
	} else {
		in.obj.FillPaint(in.obj.AddPaint(sh.paint))
		in.em.fillIsPaint = true
	}
	in.obj.FillRect(float32(r.x0), float32(r.y0), float32(r.x1-r.x0), float32(r.y1-r.y0))
}

func (in *interp) patternAverage(pat types.Object) bdf.Color {
	p := in.c.pdf
	d := p.dict(pat)
	if d == nil {
		return bdf.RGB(128, 128, 128)
	}
	if p.intOr(d["PatternType"], 1) == 2 {
		if sh := in.c.shadingFor(objKey(d["Shading"]), d["Shading"], in.res); sh != nil {
			if sh.hasAvg {
				return sh.avg
			}
			if len(sh.paint.Stops) > 0 {
				return sh.paint.Stops[len(sh.paint.Stops)/2].Color
			}
		}
	}
	return bdf.RGB(128, 128, 128)
}

// patternFill fills the path with the current fill pattern.
func (in *interp) patternFill(path bdf.PathRef, rule byte) {
	p := in.c.pdf
	d := p.dict(in.gs.fillPattern)
	if d == nil {
		in.syncFill()
		in.syncAlpha(false)
		in.obj.FillPath(path, rule)
		return
	}
	patMat := p.matrixOr(d["Matrix"], identity).mul(in.base) // pattern space → object space
	inv, ok := in.gs.ctm.inverse()
	if !ok {
		return
	}
	patInv, ok := patMat.inverse()
	if !ok {
		return
	}
	rel := patMat.mul(inv) // pattern space → current user space
	// Area to cover, in pattern space.
	area := in.gs.ctm.transformRect(in.pathBBox).intersect(in.gs.clip)
	if area.empty() {
		return
	}
	pr := patInv.transformRect(area)
	in.syncAlpha(false)
	in.save()
	in.obj.ClipPath(path, rule)
	in.transform(rel)
	switch p.intOr(d["PatternType"], 1) {
	case 2:
		if sh := in.c.shadingFor(objKey(d["Shading"]), d["Shading"], in.res); sh != nil {
			if sh.bbox != nil {
				pr = pr.intersect(*sh.bbox)
			}
			in.paintShading(sh, pr)
		}
	case 1:
		in.tilingFill(d, in.gs.fillPattern, pr)
	}
	in.restore()
}

func (in *interp) tilingFill(d types.Dict, patObj types.Object, area rect) {
	p := in.c.pdf
	sd := p.stream(patObj)
	if sd == nil {
		return
	}
	cell := in.c.formObject(objKey(patObj), sd, in.res, in.depth+1, nil)
	if cell == nil {
		return
	}
	bbox := p.rectOr(d["BBox"], rect{0, 0, 1, 1})
	xs, ys := p.numOr(d["XStep"], bbox.x1-bbox.x0), p.numOr(d["YStep"], bbox.y1-bbox.y0)
	if xs == 0 || math.IsNaN(xs) {
		xs = bbox.x1 - bbox.x0
	}
	if ys == 0 || math.IsNaN(ys) {
		ys = bbox.y1 - bbox.y0
	}
	xs, ys = math.Abs(xs), math.Abs(ys)
	if xs == 0 || ys == 0 {
		return
	}
	if p.intOr(d["PaintType"], 1) == 2 {
		// Uncoloured pattern: the cell inherits the colour given with the pattern.
		in.obj.FillColor(in.gs.fillColor)
		in.obj.StrokeColor(in.gs.fillColor)
		in.em.fillColor, in.em.fillIsPaint = in.gs.fillColor, false
		in.em.strokeColor, in.em.strokeIsPnt = in.gs.fillColor, false
	}
	i0 := math.Floor((area.x0 - bbox.x1) / xs)
	i1 := math.Ceil((area.x1 - bbox.x0) / xs)
	j0 := math.Floor((area.y0 - bbox.y1) / ys)
	j1 := math.Ceil((area.y1 - bbox.y0) / ys)
	count := (i1 - i0 + 1) * (j1 - j0 + 1)
	if count > 4096 || count <= 0 || math.IsNaN(count) {
		in.c.warnOnce("tiling-big", "tiling pattern with %d cells approximated by a flat colour", int(count))
		in.obj.FillColor(bdf.RGB(160, 160, 160))
		in.em.fillColor, in.em.fillIsPaint = bdf.RGB(160, 160, 160), false
		in.obj.FillRect(float32(area.x0), float32(area.y0), float32(area.x1-area.x0), float32(area.y1-area.y0))
		return
	}
	ref := in.p.childRef(cell)
	for j := j0; j <= j1; j++ {
		for i := i0; i <= i1; i++ {
			in.obj.UseAt(ref, float32(i*xs), float32(j*ys))
		}
	}
}

// --- inline images ---

var inlineAbbrev = map[string]string{"BPC": "BitsPerComponent", "CS": "ColorSpace", "D": "Decode", "DP": "DecodeParms", "F": "Filter", "H": "Height", "W": "Width", "IM": "ImageMask", "I": "Interpolate", "L": "Length"}
var filterAbbrev = map[string]string{"AHx": "ASCIIHexDecode", "A85": "ASCII85Decode", "LZW": "LZWDecode", "Fl": "FlateDecode", "RL": "RunLengthDecode", "CCF": "CCITTFaxDecode", "DCT": "DCTDecode"}

func (in *interp) inlineImage(l *lexer) {
	p := in.c.pdf
	d := types.NewDict()
	for {
		kind, obj, op := l.next()
		if kind == tokEOF {
			return
		}
		if kind == tokOperator {
			if op == "ID" {
				break
			}
			continue
		}
		key, ok := obj.(types.Name)
		if !ok {
			continue
		}
		kind, val, _ := l.next()
		if kind != tokOperand {
			return
		}
		k := key.Value()
		if full, ok := inlineAbbrev[k]; ok {
			k = full
		}
		d[k] = val
	}
	// Filters
	var filters []types.PDFFilter
	var names []string
	switch fv := d["Filter"].(type) {
	case types.Name:
		names = []string{fv.Value()}
	case types.Array:
		for _, e := range fv {
			names = append(names, p.name(e))
		}
	}
	var parms []types.Dict
	switch pv := d["DecodeParms"].(type) {
	case types.Dict:
		parms = []types.Dict{pv}
	case types.Array:
		for _, e := range pv {
			parms = append(parms, p.dict(e))
		}
	}
	for i, nm := range names {
		if full, ok := filterAbbrev[nm]; ok {
			nm = full
		}
		f := types.PDFFilter{Name: nm}
		if i < len(parms) {
			f.DecodeParms = parms[i]
		}
		filters = append(filters, f)
	}
	// Expected length for unfiltered data.
	length := -1
	if len(filters) == 0 {
		w, h := p.intOr(d["Width"], 0), p.intOr(d["Height"], 0)
		bpc := p.intOr(d["BitsPerComponent"], 8)
		ncomp := 1
		if p.boolOr(d["ImageMask"], false) {
			bpc = 1
		} else if cso, ok := d["ColorSpace"]; ok {
			cs := in.c.loadColorSpace(expandCSAbbrev(cso), in.res)
			ncomp = cs.n
		}
		length = (w*bpc*ncomp + 7) / 8 * h
	} else if lv := p.intOr(d["Length"], -1); lv >= 0 {
		length = lv
	}
	data := l.readInlineImageData(length)
	if cso, ok := d["ColorSpace"]; ok {
		d["ColorSpace"] = expandCSAbbrev(cso)
	}
	img, err := in.c.loadImage(d, data, filters, in.res, in.gs.fillColor)
	if err != nil {
		in.c.warnf("inline image: %v", err)
		return
	}
	h := in.c.doc.AddImage(img.data)
	in.syncDraw(in.curTarget())
	in.enterMask()
	ref := in.obj.AddImage(h)
	in.syncAlpha(false)
	in.save()
	in.transform(matrix{1, 0, 0, -1, 0, 1})
	in.obj.Image(ref, 0, 0, 1, 1)
	in.restore()
}

func expandCSAbbrev(o types.Object) types.Object {
	switch v := o.(type) {
	case types.Name:
		switch v.Value() {
		case "G":
			return types.Name("DeviceGray")
		case "RGB":
			return types.Name("DeviceRGB")
		case "CMYK":
			return types.Name("DeviceCMYK")
		case "I":
			return types.Name("Indexed")
		}
	case types.Array:
		out := make(types.Array, len(v))
		for i, e := range v {
			out[i] = expandCSAbbrev(e)
		}
		return out
	}
	return o
}

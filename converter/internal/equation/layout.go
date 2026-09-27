package equation

import (
	"math"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/sfnt"
)

// Style is how a whole formula is set.
type Style struct {
	Size    float64 // font size in points
	Display bool    // display style: a formula on a line of its own
	Color   bdf.Color
	Bold    bool
}

// Layout lays a formula out.
func (e *Engine) Layout(n Node, st Style) *Box {
	v := env{base: st.Size, display: st.Display, color: st.Color, bold: st.Bold}
	if v.base <= 0 {
		v.base = 12
	}
	if v.color == 0 {
		v.color = bdf.RGBA(0, 0, 0, 255)
	}
	return e.node(n, v, nil)
}

// vtarget is the height and depth a vertically stretchy operator grows to
// cover.
type vtarget struct{ asc, desc float64 }

func (e *Engine) node(n Node, v env, t *vtarget) *Box {
	switch n := n.(type) {
	case *Row:
		if n == nil {
			break
		}
		return e.row(n, v, t)
	case *Atom:
		if n == nil {
			break
		}
		return e.atom(n, v, t)
	case *Frac:
		return e.frac(n, v)
	case *Radical:
		return e.radical(n, v)
	case *Scripts:
		return e.scripts(n, v, t)
	case *UnderOver:
		return e.underOver(n, v, t)
	case *Table:
		return e.table(n, v)
	case *Space:
		b := newBox()
		b.W, b.class = e.em(v, n.Width), None
		return b
	case *Styled:
		return e.styled(n, v, t)
	case *Enclose:
		return e.enclose(n, v)
	case *Phantom:
		return e.phantom(n, v)
	case *Bar:
		return e.bar(n, v)
	}
	b := newBox()
	b.class = None
	return b
}

// isStretchyV reports whether a node is an operator (possibly with
// scripts) that grows vertically to the height of its row.
func isStretchyV(n Node) bool {
	switch n := n.(type) {
	case *Atom:
		if n == nil || n.Kind != Op || n.LargeOp {
			return false
		}
		return n.Stretchy && !lookupOp(normalizeOp(n.Text)).horizontal
	case *Styled:
		return isStretchyV(n.Kid)
	case *Scripts:
		return isStretchyV(n.Base)
	case *Row:
		return n != nil && len(n.Kids) == 1 && isStretchyV(n.Kids[0])
	}
	return false
}

func (e *Engine) row(r *Row, v env, t *vtarget) *Box {
	boxes := make([]*Box, len(r.Kids))
	var asc, desc float64
	have := false
	for i, k := range r.Kids {
		if isStretchyV(k) {
			continue
		}
		b := e.node(k, v, nil)
		boxes[i] = b
		if b.class != None || b.H > 0 || b.D > 0 {
			asc, desc = math.Max(asc, b.H), math.Max(desc, b.D)
			have = true
		}
	}
	for i, k := range r.Kids {
		if boxes[i] != nil {
			continue
		}
		tt := t
		if have {
			tt = &vtarget{asc, desc}
		}
		boxes[i] = e.node(k, v, tt)
	}
	return e.hlist(boxes, v, r.Class)
}

// spacing is the space between atoms of two classes in TeX (TeXbook p.
// 170): 1 thin, 2 medium, 3 thick; negative only in display and text
// style.
var spacing = [8][8]int8{
	//        Ord Op Bin Rel Open Close Punct Inner
	/*Ord*/ {0, 1, -2, -3, 0, 0, 0, -1},
	/*Op*/ {1, 1, 0, -3, 0, 0, 0, -1},
	/*Bin*/ {-2, -2, 0, 0, -2, 0, 0, -2},
	/*Rel*/ {-3, -3, 0, 0, -3, 0, 0, -3},
	/*Open*/ {0, 0, 0, 0, 0, 0, 0, 0},
	/*Close*/ {0, 1, -2, -3, 0, 0, 0, -1},
	/*Punct*/ {-1, -1, 0, -1, -1, -1, -1, -1},
	/*Inner*/ {-1, 1, -2, -3, -1, 0, -1, -1},
}

func classIndex(c Class) int {
	switch c {
	case LargeOp:
		return 1
	case Bin:
		return 2
	case Rel:
		return 3
	case Open:
		return 4
	case Close:
		return 5
	case Punct:
		return 6
	case Inner:
		return 7
	}
	return 0
}

// space returns the space between atoms of classes l and r.
func (e *Engine) space(l, r Class, v env) float64 {
	s := spacing[classIndex(l)][classIndex(r)]
	if s < 0 {
		if v.level > 0 {
			return 0
		}
		s = -s
	}
	mu := e.size(v) / 18
	switch s {
	case 1:
		return 3 * mu
	case 2:
		return 4 * mu
	case 3:
		return 5 * mu
	}
	return 0
}

// hlist sets boxes side by side with the spacing of their classes.
func (e *Engine) hlist(boxes []*Box, v env, class Class) *Box {
	cls := make([]Class, len(boxes))
	prev := -1
	for i, b := range boxes {
		c := b.class
		if c == None || c == ClassAuto && b.W == 0 {
			cls[i] = None
			continue
		}
		if c == ClassAuto {
			c = Ord
		}
		// TeX's rules 5 and 6: a binary operator with nothing to operate
		// on is an ordinary atom
		if c == Bin && (prev < 0 || cls[prev] == Bin || cls[prev] == LargeOp || cls[prev] == Rel || cls[prev] == Open || cls[prev] == Punct) {
			c = Ord
		}
		if (c == Rel || c == Close || c == Punct) && prev >= 0 && cls[prev] == Bin {
			cls[prev] = Ord
		}
		cls[i] = c
		prev = i
	}
	if prev >= 0 && cls[prev] == Bin {
		cls[prev] = Ord
	}
	out := newBox()
	out.class = class
	if class == ClassAuto {
		out.class = Ord
	}
	x := 0.0
	last := -1
	for i, b := range boxes {
		if cls[i] != None && last >= 0 {
			x += e.space(cls[last], cls[i], v)
		}
		out.add(b, x, 0)
		out.grow(b, 0)
		x += b.W
		if b.single && b.Italic > 0 && i < len(boxes)-1 {
			x += b.Italic
		}
		if cls[i] != None {
			last = i
		}
	}
	out.W = x
	if len(boxes) == 1 {
		b := boxes[0]
		out.single, out.accentX, out.largeOp, out.movable, out.stretchH, out.accentLow = b.single, b.accentX, b.largeOp, b.movable, b.stretchH, b.accentLow
		out.Italic = b.Italic
		if class == ClassAuto {
			out.class = b.class
		}
	}
	return out
}

func (e *Engine) styled(s *Styled, v env, t *vtarget) *Box {
	switch s.Display {
	case 1:
		v.display = true
	case 2:
		v.display = false
	}
	if s.Level > 0 {
		v.level = int(s.Level) - 1
	}
	if s.LevelUp {
		v.level++
	}
	if s.Color != nil {
		v.color = *s.Color
	}
	if s.Size > 0 {
		v.base = s.Size * v.base / e.size(v)
	}
	if s.Base > 0 {
		v.base = s.Base
	}
	if s.Scale > 0 {
		v.base *= s.Scale
	}
	if s.Variant != VarAuto {
		v.variant = s.Variant
	}
	if s.Bold {
		v.bold = true
	}
	return e.node(s.Kid, v, t)
}

// atom lays a token out.
func (e *Engine) atom(a *Atom, v env, t *vtarget) *Box {
	text := a.Text
	class := a.Class
	if a.Kind == Op {
		text = normalizeOp(text)
		if class == ClassAuto {
			class = lookupOp(text).class
		}
		if r, ok := singleRune(text); ok {
			switch {
			case a.LargeOp:
				b := e.largeOp(r, v)
				b.class, b.movable = class, a.MovableLimits
				return b
			case a.Stretchy && (t != nil || a.MinSize > 0) && !lookupOp(text).horizontal:
				b := e.fence(a, r, v, t)
				b.class = class
				return b
			}
		}
	}
	if class == ClassAuto {
		class = Ord
		if a.Kind == Ident && len([]rune(text)) > 1 && IsFunctionName(text) {
			class = LargeOp
		}
	}
	b := e.token(a, text, v)
	b.class = class
	if a.Kind == Op && a.Stretchy && lookupOp(text).horizontal {
		b.stretchH = true
	}
	if class == LargeOp && a.MovableLimits {
		// lim, max: limits under them in display style, at their side
		// otherwise
		b.largeOp, b.movable = true, true
	}
	if a.Kind == Op && a.Spaced {
		return e.spaced(b, a, v)
	}
	return b
}

// spaced surrounds an operator with the space it asks for instead of the
// space of its class.
func (e *Engine) spaced(b *Box, a *Atom, v env) *Box {
	l, r := a.LSpace, a.RSpace
	out := newBox()
	out.add(b, e.em(v, l), 0)
	out.grow(b, 0)
	out.W = b.W + e.em(v, l+r)
	out.class = None
	return out
}

// token lays out the characters of a token side by side.
func (e *Engine) token(a *Atom, text string, v env) *Box {
	runes := []rune(text)
	variant := a.Variant
	if variant == VarAuto {
		variant = v.variant
	}
	if variant == VarAuto {
		variant = VarNormal
		if a.Kind == Ident && len(runes) == 1 && isLetter(runes[0]) && !upperGreek(runes[0]) {
			variant = VarItalic
		}
		if v.bold {
			switch {
			case variant == VarItalic:
				variant = VarBoldItalic
			case a.Kind != Op:
				variant = VarBold
			}
		}
	}
	b := newBox()
	size := e.size(v)
	x := 0.0
	var last *glyphMetrics
	for _, r := range runes {
		var gl glyph
		switch {
		case a.Kind == Text && a.Font != "":
			gl = e.resolve(r, a.Font, variant.isBold(), variant.isItalic())
		case a.Kind == Text:
			gl = e.resolve(r, "", variant.isBold(), variant.isItalic())
		default:
			gl = e.char(r, variant, "")
		}
		adv := e.record(gl) * size
		m := e.metrics(gl.fc, gl.gid)
		if !e.hasGlyph(gl) && !unicode.IsSpace(r) {
			// no font has it: the viewer draws it with a font of its own
			m = &glyphMetrics{adv: adv / size, ink: sfnt.Rect{XMax: adv / size, YMin: -0.15, YMax: 0.7}, hasInk: true}
		}
		it := item{kind: iGlyph, x: x, fc: gl.fc, size: size, text: string(gl.r), adv: adv, sx: 1, sy: 1, color: v.color,
			ink: [4]float64{0, 0, adv, 0}}
		if m.hasInk {
			it.ink = [4]float64{m.ink.XMin * size, -m.ink.YMax * size, m.ink.XMax * size, -m.ink.YMin * size}
			b.H = math.Max(b.H, m.ink.YMax*size)
			b.D = math.Max(b.D, -m.ink.YMin*size)
		}
		b.items = append(b.items, it)
		x += adv
		last = m
	}
	b.W = x
	if last != nil {
		b.Italic = last.italic * size
		if len(runes) == 1 {
			b.single = true
			if last.hasAccent {
				b.accentX = last.topAccent * size
			}
			b.accentLow = last.hasInk && last.ink.YMin > 0.3
		}
	}
	return b
}

// glyphBox is a box of one glyph of the formula font (a variant or part
// of the character r), drawn at its baseline.
func (e *Engine) glyphBox(r rune, gid uint16, v env) *Box {
	size := e.size(v)
	m := e.metrics(e.face, gid)
	pr := e.glyphRune(r, gid)
	it := item{kind: iGlyph, fc: e.face, size: size, text: string(pr), adv: m.adv * size, sx: 1, sy: 1, color: v.color,
		ink: [4]float64{0, 0, m.adv * size, 0}}
	if pr != r {
		it.alt = string(r)
	}
	b := newBox()
	if m.hasInk {
		it.ink = [4]float64{m.ink.XMin * size, -m.ink.YMax * size, m.ink.XMax * size, -m.ink.YMin * size}
		b.H, b.D = m.ink.YMax*size, -m.ink.YMin*size
	}
	b.items = []item{it}
	b.W = m.adv * size
	b.Italic = m.italic * size
	if m.hasAccent {
		b.accentX = m.topAccent * size
	}
	b.accentLow = m.hasInk && m.ink.YMin > 0.3
	return b
}

// shiftBox moves what b draws down by dy (up for negative dy).
func shiftBox(b *Box, dy float64) *Box {
	if dy == 0 {
		return b
	}
	for i := range b.items {
		b.items[i].y += dy
	}
	b.H -= dy
	b.D += dy
	return b
}

// largeOp lays a large operator out: in display style a variant at least
// DisplayOperatorMinHeight tall, centered on the math axis.
func (e *Engine) largeOp(r rune, v env) *Box {
	var b *Box
	gl := e.resolve(r, "", false, false)
	if gl.fc == e.face && e.face != nil {
		gid := gl.gid
		if con := e.mt.Vertical(gid); v.display && con != nil {
			for _, vr := range con.Variants {
				gid = vr.Glyph
				if float64(vr.Advance)/e.upem >= e.c.displayOpMin {
					break
				}
			}
		}
		b = e.glyphBox(r, gid, v)
	} else {
		a := &Atom{Kind: Op, Text: string(r)}
		b = e.token(a, string(r), v)
		if v.display {
			b = scaleBox(b, 1.4, 1.4)
		}
	}
	axis := e.em(v, e.c.axis)
	b = shiftBox(b, (b.H-b.D)/2-axis)
	b.largeOp = true
	b.single = false
	return b
}

// scaleBox scales a box of glyphs about its origin (text fonts, which have
// no larger variants).
func scaleBox(b *Box, sx, sy float64) *Box {
	for i := range b.items {
		it := &b.items[i]
		it.x *= sx
		it.y *= sy
		it.sx *= sx
		it.sy *= sy
		it.ink = [4]float64{it.ink[0] * sx, it.ink[1] * sy, it.ink[2] * sx, it.ink[3] * sy}
	}
	b.W *= sx
	b.H *= sy
	b.D *= sy
	b.Italic *= sx
	return b
}

// fence lays out a delimiter grown to cover t (or to its MinSize).
func (e *Engine) fence(a *Atom, r rune, v env, t *vtarget) *Box {
	size := e.size(v)
	axis := e.em(v, e.c.axis)
	var asc, desc float64
	if t != nil {
		asc, desc = t.asc, t.desc
	}
	if a.Symmetric || t == nil {
		h := math.Max(asc-axis, desc+axis)
		asc, desc = axis+h, h-axis
	}
	total := asc + desc
	if a.Shortfall {
		total = math.Max(total*0.901, total-0.5*size)
	}
	if a.MinSize > 0 {
		m := a.MinSize * v.base
		if m > total {
			total = m
			if t == nil || a.Symmetric {
				asc, desc = axis+m/2, m/2-axis
			}
		}
	}
	b, grown := e.stretchV(r, total, v)
	if grown {
		s := b.H + b.D
		want := s/2 + (asc-desc)/2
		b = shiftBox(b, b.H-want)
	}
	return b
}

func (e *Engine) frac(f *Frac, v env) *Box {
	if f.Kind == FracLinear {
		num, den := e.node(f.Num, v, nil), e.node(f.Den, v, nil)
		return e.hlist([]*Box{num, e.atom(&Atom{Kind: Op, Text: "/", Class: Ord}, v, nil), den}, v, Ord)
	}
	num, den := e.node(f.Num, v.num(), nil), e.node(f.Den, v.den(), nil)
	if f.Kind == FracSkewed {
		return e.skewed(num, den, v)
	}
	c := e.c
	size := e.size(v)
	axis := c.axis * size
	t := c.rule * size
	if f.Thickness > 0 {
		t = f.Thickness * size
	}
	if f.NoBar {
		t = 0
	}
	var up, down float64
	if t > 0 {
		gapN, gapD := c.numGap, c.denGap
		up, down = c.numShift, c.denShift
		if v.display {
			gapN, gapD = c.numGapDisplay, c.denGapDisplay
			up, down = c.numShiftDisplay, c.denShiftDisplay
		}
		up = math.Max(up*size, axis+t/2+gapN*size+num.D)
		down = math.Max(down*size, den.H+gapD*size+t/2-axis)
	} else {
		gap := c.stackGap
		up, down = c.stackTop, c.stackBottom
		if v.display {
			gap = c.stackGapDisplay
			up, down = c.stackTopDisplay, c.stackBottomDisp
		}
		up, down, gap = up*size, down*size, gap*size
		if actual := (up - num.D) - (den.H - down); actual < gap {
			up += (gap - actual) / 2
			down += (gap - actual) / 2
		}
	}
	pad := 0.12 * size
	inner := math.Max(num.W, den.W)
	b := newBox()
	b.W = inner + 2*pad
	b.add(num, pad+alignOffset(inner, num.W, f.NumAlign), -up)
	b.grow(num, -up)
	b.add(den, pad+alignOffset(inner, den.W, f.DenAlign), down)
	b.grow(den, down)
	if t > 0 {
		b.rule(pad, -(axis + t/2), inner, t, v.color)
		b.H = math.Max(b.H, axis+t/2)
		b.D = math.Max(b.D, t/2-axis)
	}
	return b
}

func alignOffset(width, w float64, a Align) float64 {
	switch a {
	case Left:
		return 0
	case Right:
		return width - w
	}
	return (width - w) / 2
}

// skewed sets a fraction as a raised numerator, a slash and a lowered
// denominator.
func (e *Engine) skewed(num, den *Box, v env) *Box {
	size := e.size(v)
	axis := e.em(v, e.c.axis)
	vgap, hgap := e.c.skewV*size, e.c.skewH*size
	up := axis + vgap/2 + num.D
	down := den.H - axis + vgap/2
	asc, desc := up+num.H, down+den.D
	slash, _ := e.stretchV('/', asc+desc, v)
	slash = shiftBox(slash, slash.H-((slash.H+slash.D)/2+(asc-desc)/2))
	b := newBox()
	b.add(num, 0, -up)
	b.grow(num, -up)
	sx := num.W + (hgap-slash.W)/2
	b.add(slash, sx, 0)
	b.grow(slash, 0)
	b.add(den, num.W+hgap, down)
	b.grow(den, down)
	b.W = num.W + hgap + den.W
	return b
}

func (e *Engine) radical(r *Radical, v env) *Box {
	base := e.node(r.Base, v.crampedEnv(), nil)
	c := e.c
	size := e.size(v)
	gap := c.radGap * size
	if v.display {
		gap = c.radGapDisplay * size
	}
	t := c.radRule * size
	target := base.H + base.D + gap + t
	var sign *Box
	var top float64 // the top of the rule above the baseline
	if e.face == nil {
		top = base.H + gap + t
		sign = e.pathRadical(top, math.Max(base.D, 0.1*size), v)
	} else {
		sign, _ = e.stretchV('√', target, v)
		if extra := sign.H + sign.D - target; extra > 0 {
			gap += extra / 2
		}
		top = base.H + gap + t
		sign = shiftBox(sign, sign.H-top)
	}
	b := newBox()
	x := 0.0
	if r.Degree != nil {
		dv := v
		dv.level += 2
		dv.display = false
		deg := e.node(r.Degree, dv, nil)
		before, after := c.radKernBefore*size, c.radKernAfter*size
		raise := c.radDegreeRaise*(sign.H+sign.D) - sign.D
		b.add(deg, before, -raise-deg.D)
		b.grow(deg, -raise-deg.D)
		x = math.Max(0, before+deg.W+after)
	}
	b.add(sign, x, 0)
	b.grow(sign, 0)
	x += sign.W
	b.add(base, x, 0)
	b.grow(base, 0)
	b.rule(x, -top, base.W, t, v.color)
	b.H = math.Max(b.H, top+c.radExtra*size)
	b.W = x + base.W
	return b
}

// scripts attaches sub- and superscripts to a base, after it and before it.
func (e *Engine) scripts(s *Scripts, v env, t *vtarget) *Box {
	base := e.node(s.Base, v, t)
	return e.attach(base, s.Sub, s.Sup, s.PreSub, s.PreSup, v)
}

func (e *Engine) attach(base *Box, subN, supN, preSubN, preSupN Node, v env) *Box {
	c := e.c
	size := e.size(v)
	var sub, sup, psub, psup *Box
	if subN != nil {
		sub = e.node(subN, v.sub(), nil)
	}
	if supN != nil {
		sup = e.node(supN, v.script(), nil)
	}
	if preSubN != nil {
		psub = e.node(preSubN, v.sub(), nil)
	}
	if preSupN != nil {
		psup = e.node(preSupN, v.script(), nil)
	}
	if sub == nil && sup == nil && psub == nil && psup == nil {
		return base
	}
	subH, _ := maxDims(sub, psub)
	_, supD := maxDims(sup, psup)
	hasSub, hasSup := sub != nil || psub != nil, sup != nil || psup != nil
	baseH, baseD := base.H, base.D
	if base.single && !base.largeOp {
		// a character: its scripts hang from the font's positions, not
		// from its own height (TeX's rule 18a)
		baseH, baseD = 0, 0
	}
	var subShift, supShift float64
	if hasSub {
		subShift = math.Max(c.subShift*size, baseD+c.subDrop*size)
		if !hasSup {
			subShift = math.Max(subShift, subH-c.subTopMax*size)
		}
	}
	if hasSup {
		up := c.supShift
		if v.cramped {
			up = c.supShiftCramped
		}
		supShift = math.Max(up*size, baseH-c.supDrop*size)
		supShift = math.Max(supShift, c.supBottomMin*size+supD)
	}
	if hasSub && hasSup {
		if gap := (subShift - subH) + (supShift - supD); gap < c.subSupGap*size {
			subShift += c.subSupGap*size - gap
		}
		if d := c.supBottomMaxWithSub*size - (supShift - supD); d > 0 {
			supShift += d
			subShift -= d
		}
	}
	after := c.spaceAfterScript * size
	b := newBox()
	x := 0.0
	if psub != nil || psup != nil {
		w := math.Max(width(psub), width(psup))
		if psup != nil {
			b.add(psup, x+w-psup.W, -supShift)
			b.grow(psup, -supShift)
		}
		if psub != nil {
			b.add(psub, x+w-psub.W, subShift)
			b.grow(psub, subShift)
		}
		x += w + after
	}
	b.add(base, x, 0)
	b.grow(base, 0)
	x += base.W
	italic := base.Italic
	subX, supX := x, x+italic
	if base.largeOp {
		// an integral's subscript tucks under it
		subX, supX = x-italic, x
	}
	end := x
	if sup != nil {
		b.add(sup, supX, -supShift)
		b.grow(sup, -supShift)
		end = math.Max(end, supX+sup.W)
	}
	if sub != nil {
		b.add(sub, subX, subShift)
		b.grow(sub, subShift)
		end = math.Max(end, subX+sub.W)
	}
	if sub != nil || sup != nil {
		end += after
	} else {
		end += italic
	}
	b.W = end
	b.class = base.class
	b.largeOp = base.largeOp
	return b
}

func maxDims(a, b *Box) (h, d float64) {
	for _, x := range []*Box{a, b} {
		if x != nil {
			h, d = math.Max(h, x.H), math.Max(d, x.D)
		}
	}
	return h, d
}

func width(b *Box) float64 {
	if b == nil {
		return 0
	}
	return b.W
}

// underOver places limits, accents and braces under and over a base.
func (e *Engine) underOver(u *UnderOver, v env, t *vtarget) *Box {
	bv := v
	if u.AccentOver {
		bv.cramped = true
	}
	base := e.node(u.Base, bv, t)
	if base.largeOp && base.movable && !v.display && !u.AccentOver && !u.AccentUnder {
		return e.attach(base, u.Under, u.Over, nil, nil, v)
	}
	c := e.c
	size := e.size(v)
	ov, uv := v.script(), v.sub()
	if u.AccentOver {
		ov = v
		ov.display = false
	}
	if u.AccentUnder {
		uv = v
		uv.display = false
	}
	over := e.horizontalScript(u.Over, ov, base.W)
	under := e.horizontalScript(u.Under, uv, base.W)
	if base.stretchH {
		// a brace or an arrow grows to what is over and under it
		w := math.Max(width(over), width(under))
		if w > base.W {
			if a := coreAtom(u.Base); a != nil {
				if r, ok := singleRune(normalizeOp(a.Text)); ok {
					base = e.stretchH(r, w, v)
				}
			}
		}
	}
	var overGap, overShift, underGap, underShift, extraAsc, extraDesc float64
	switch {
	case base.largeOp:
		overGap, overShift = c.upperGap*size, c.upperRise*size
		underGap, underShift = c.lowerGap*size, c.lowerDrop*size
	case base.stretchH:
		overGap, overShift = c.stretchGapAbove*size, c.stretchTop*size
		underGap, underShift = c.stretchGapBel*size, c.stretchBottom*size
	default:
		overGap, extraAsc = c.overGap*size, c.overExtra*size
		underGap, extraDesc = c.underGap*size, c.underExt*size
	}
	w := base.W
	if over != nil {
		w = math.Max(w, over.W)
	}
	if under != nil {
		w = math.Max(w, under.W)
	}
	italic := 0.0
	if base.largeOp {
		italic = base.Italic
	}
	b := newBox()
	bx := (w - base.W) / 2
	b.add(base, bx, 0)
	b.grow(base, 0)
	if over != nil {
		var shift, x float64
		switch {
		case u.AccentOver:
			if over.accentLow && e.face != nil {
				// an accent glyph made to sit over a lower-case letter
				shift = math.Max(0, base.H-c.accentBase*size)
			} else {
				shift = base.H + over.D + 0.08*size
			}
			at := base.W / 2
			if base.single {
				at = base.accentAt()
			}
			oat := inkCenter(over)
			if !math.IsNaN(over.accentX) {
				oat = over.accentX
			}
			x = bx + at - oat
		default:
			shift = math.Max(overShift, base.H+overGap+over.D)
			x = (w-over.W)/2 + italic/2
		}
		b.add(over, x, -shift)
		b.grow(over, -shift)
		b.H += extraAsc
	}
	if under != nil {
		var shift float64
		if u.AccentUnder {
			shift = base.D + under.H + 0.08*size
		} else {
			shift = math.Max(underShift, base.D+underGap+under.H)
		}
		b.add(under, (w-under.W)/2-italic/2, shift)
		b.grow(under, shift)
		b.D += extraDesc
	}
	b.W = w
	b.class = base.class
	return b
}

// horizontalScript lays out what goes over or under a base: operators that
// grow horizontally (braces, arrows, wide accents) grow to its width, and
// the spacing forms of accents become the combining ones, which a formula
// font places for accents.
func (e *Engine) horizontalScript(n Node, v env, w float64) *Box {
	if n == nil {
		return nil
	}
	kid, kv := unwrap(n, v)
	if a, ok := kid.(*Atom); ok && a.Kind == Op {
		if r, ok := singleRune(normalizeOp(a.Text)); ok {
			if a.Accent && e.face != nil {
				if cr := combiningAccent(r); cr != r && e.face.Loaded.Has(cr) {
					r = cr
				}
			}
			if a.Stretchy && (lookupOp(string(r)).horizontal || e.growsH(r)) {
				return e.stretchH(r, w, kv)
			}
			return e.token(&Atom{Kind: Op, Text: string(r)}, string(r), kv)
		}
	}
	return e.node(n, v, nil)
}

// growsH reports whether the formula font has wider forms of r.
func (e *Engine) growsH(r rune) bool {
	if e.face == nil {
		return false
	}
	g, ok := e.face.Loaded.Glyph(r)
	return ok && e.mt.Horizontal(g) != nil
}

// unwrap applies the styles around a node to v.
func unwrap(n Node, v env) (Node, env) {
	for {
		switch s := n.(type) {
		case *Styled:
			if s.Color != nil {
				v.color = *s.Color
			}
			if s.Bold {
				v.bold = true
			}
			if s.Variant != VarAuto {
				v.variant = s.Variant
			}
			if s.Display != 0 || s.Level != 0 || s.LevelUp || s.Size > 0 || s.Base > 0 || s.Scale > 0 {
				return n, v
			}
			n = s.Kid
		case *Row:
			if s == nil || len(s.Kids) != 1 {
				return n, v
			}
			n = s.Kids[0]
		default:
			return n, v
		}
	}
}

// coreAtom returns the operator a node is (inside styles and rows of one).
func coreAtom(n Node) *Atom {
	switch n := n.(type) {
	case *Atom:
		return n
	case *Styled:
		return coreAtom(n.Kid)
	case *Row:
		if n != nil && len(n.Kids) == 1 {
			return coreAtom(n.Kids[0])
		}
	}
	return nil
}

func (e *Engine) table(t *Table, v env) *Box {
	size := e.size(v)
	cv := v
	cv.display = t.Display
	nCols := 0
	for _, row := range t.Rows {
		nCols = max(nCols, len(row))
	}
	if nCols == 0 {
		return newBox()
	}
	cells := make([][]*Box, len(t.Rows))
	colW := make([]float64, nCols)
	rowH := make([]float64, len(t.Rows))
	rowD := make([]float64, len(t.Rows))
	strutH, strutD := 0.8*size, 0.3*size
	for i, row := range t.Rows {
		cells[i] = make([]*Box, len(row))
		rowH[i], rowD[i] = strutH, strutD
		for j, cell := range row {
			if t.Aligned && j%2 == 1 {
				// the second half of a pair starts with an empty atom, so
				// that the relation after the alignment point is spaced
				// ("&=" as amsmath sets it)
				cell = leadingOrd(cell)
			}
			b := e.node(cell, cv, nil)
			cells[i][j] = b
			colW[j] = math.Max(colW[j], b.W)
			rowH[i] = math.Max(rowH[i], b.H)
			rowD[i] = math.Max(rowD[i], b.D)
		}
	}
	rowGap := t.RowGap * size
	if t.RowGap == 0 {
		rowGap = 0.2 * size
		if t.Display {
			rowGap = 0.3 * size
		}
	}
	colGap := t.ColGap * size
	if t.ColGap == 0 {
		colGap = size
		if t.Aligned {
			colGap = 2 * size
		}
	}
	align := func(j int) Align {
		if t.Aligned {
			if j%2 == 0 {
				return Right
			}
			return Left
		}
		if len(t.ColAlign) == 0 {
			return Center
		}
		if j < len(t.ColAlign) {
			return t.ColAlign[j]
		}
		return t.ColAlign[len(t.ColAlign)-1]
	}
	colX := make([]float64, nCols)
	x := 0.0
	for j := range colW {
		if j > 0 {
			if t.Aligned && j%2 == 1 {
				// the halves of a pair meet at the alignment point
			} else {
				x += colGap
			}
		}
		colX[j] = x
		x += colW[j]
	}
	total := 0.0
	for i := range t.Rows {
		total += rowH[i] + rowD[i]
		if i > 0 {
			total += rowGap
		}
	}
	axis := e.em(v, e.c.axis)
	top := axis + total/2
	b := newBox()
	y := -top
	for i, row := range cells {
		if i > 0 {
			y += rowGap
		}
		y += rowH[i]
		for j, cell := range row {
			b.add(cell, colX[j]+alignOffset(colW[j], cell.W, align(j)), y)
		}
		y += rowD[i]
	}
	b.W = x
	b.H, b.D = top, total-top
	return b
}

func (e *Engine) enclose(n *Enclose, v env) *Box {
	kid := e.node(n.Kid, v, nil)
	size := e.size(v)
	t := e.c.rule * size
	pad := 0.2 * size
	if !(n.Top || n.Bottom || n.Left || n.Right || n.Round || n.Circle) {
		pad = 0
	}
	b := newBox()
	b.add(kid, pad, 0)
	b.grow(kid, 0)
	b.W = kid.W + 2*pad
	x0, x1 := 0.0, b.W
	y0, y1 := -(kid.H + pad), kid.D+pad
	if pad == 0 {
		y0, y1 = -kid.H, kid.D
	}
	col := v.color
	if n.Round || n.Circle {
		p := &bdf.Path{}
		if n.Circle {
			p.Ellipse(f32(b.W/2), f32((y0+y1)/2), f32(b.W/2), f32((y1-y0)/2), 0, 0, 2*math.Pi, false)
		} else {
			p.RoundRect(f32(x0), f32(y0), f32(x1-x0), f32(y1-y0), f32(pad))
		}
		b.stroke(p, 0, 0, t, col, [4]float64{x0, y0, x1, y1})
	}
	if n.Top {
		b.rule(x0, y0, x1-x0, t, col)
	}
	if n.Bottom {
		b.rule(x0, y1-t, x1-x0, t, col)
	}
	if n.Left {
		b.rule(x0, y0, t, y1-y0, col)
	}
	if n.Right {
		b.rule(x1-t, y0, t, y1-y0, col)
	}
	mid := -e.em(v, e.c.axis)
	if n.StrikeH {
		b.rule(x0, mid-t/2, x1-x0, t, col)
	}
	if n.StrikeV {
		b.rule((x0+x1)/2-t/2, y0, t, y1-y0, col)
	}
	line := func(ax, ay, bx, by float64) {
		p := (&bdf.Path{}).MoveTo(f32(ax), f32(ay)).LineTo(f32(bx), f32(by))
		b.stroke(p, 0, 0, t, col, [4]float64{math.Min(ax, bx), math.Min(ay, by), math.Max(ax, bx), math.Max(ay, by)})
	}
	if n.StrikeUp {
		line(x0, y1, x1, y0)
	}
	if n.StrikeDown {
		line(x0, y0, x1, y1)
	}
	b.H, b.D = -y0, y1
	if pad == 0 {
		b.H, b.D = kid.H, kid.D
	}
	return b
}

// bar draws a rule over or under what it holds.
func (e *Engine) bar(n *Bar, v env) *Box {
	c := e.c
	size := e.size(v)
	kv := v
	if !n.Under {
		kv.cramped = true
	}
	kid := e.node(n.Kid, kv, nil)
	b := newBox()
	b.add(kid, 0, 0)
	b.grow(kid, 0)
	b.W, b.class = kid.W, kid.class
	if n.Under {
		t := c.underRule * size
		y := kid.D + c.underGap*size
		b.rule(0, y, kid.W, t, v.color)
		b.D = math.Max(b.D, y+t+c.underExt*size)
		return b
	}
	t := c.overRule * size
	y := kid.H + c.overGap*size + t
	b.rule(0, -y, kid.W, t, v.color)
	b.H = math.Max(b.H, y+c.overExtra*size)
	return b
}

func (e *Engine) phantom(p *Phantom, v env) *Box {
	kid := e.node(p.Kid, v, nil)
	b := newBox()
	if p.Show {
		b.add(kid, 0, 0)
	}
	b.W, b.H, b.D = kid.W, kid.H, kid.D
	if p.ZeroWidth {
		b.W = 0
	}
	if p.ZeroAsc {
		b.H = 0
	}
	if p.ZeroDesc {
		b.D = 0
	}
	b.class = kid.class
	return b
}

func f32(v float64) float32 { return float32(v) }

// leadingOrd puts an empty ordinary atom before what a row holds.
func leadingOrd(n Node) Node {
	kids := []Node{n}
	if r, ok := n.(*Row); ok && r != nil && r.Class == ClassAuto {
		kids = r.Kids
	}
	return &Row{Kids: append([]Node{&Row{Class: Ord}}, kids...)}
}

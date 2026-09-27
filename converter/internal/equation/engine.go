package equation

import (
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/sfnt"
)

// Fonts are what formulas are laid out with.
type Fonts struct {
	Set *fontset.Set
	// Math is the formula font asked for ("" for the default): a font with
	// a MATH table resolves it (see fontdb.DB.ResolveMath).
	Math string
	// Text returns the face that draws r in a text font: for characters
	// the formula font lacks (East Asian text) and for text asked in a
	// family (family "" is the text around the formula). It is also what
	// formulas are drawn with when no formula font is available.
	Text func(family string, r rune, bold, italic bool) *fontset.Choice
	// Warn receives problems, such as a missing formula font.
	Warn func(msg string)
}

// Engine lays formulas out with one set of fonts. It is not safe for
// concurrent use.
type Engine struct {
	f      Fonts
	face   *fontset.Choice // the formula font; nil without one
	mt     *sfnt.Math
	upem   float64
	c      consts
	glyphs map[glyphKey]*glyphMetrics
}

// consts are the layout parameters in em.
type consts struct {
	scriptScale, scriptScriptScale float64
	delimitedMin, displayOpMin     float64
	axis, accentBase               float64
	subShift, subTopMax, subDrop   float64
	supShift, supShiftCramped      float64
	supBottomMin, supDrop          float64
	subSupGap, supBottomMaxWithSub float64
	spaceAfterScript               float64
	upperGap, upperRise            float64
	lowerGap, lowerDrop            float64
	stackTop, stackTopDisplay      float64
	stackBottom, stackBottomDisp   float64
	stackGap, stackGapDisplay      float64
	stretchTop, stretchBottom      float64
	stretchGapAbove, stretchGapBel float64
	numShift, numShiftDisplay      float64
	denShift, denShiftDisplay      float64
	numGap, numGapDisplay          float64
	rule                           float64
	denGap, denGapDisplay          float64
	skewH, skewV                   float64
	overGap, overRule, overExtra   float64
	underGap, underRule, underExt  float64
	radGap, radGapDisplay          float64
	radRule, radExtra              float64
	radKernBefore, radKernAfter    float64
	radDegreeRaise                 float64
	minOverlap                     float64
}

// defaultConsts are TeX's parameters (those of Computer Modern at 10 pt,
// in em) for text fonts, which have no MATH table.
var defaultConsts = consts{
	scriptScale: 0.7, scriptScriptScale: 0.5, delimitedMin: 1.5, displayOpMin: 1.3,
	axis: 0.25, accentBase: 0.45,
	subShift: 0.15, subTopMax: 0.344, subDrop: 0.05,
	supShift: 0.413, supShiftCramped: 0.289, supBottomMin: 0.108, supDrop: 0.386,
	subSupGap: 0.16, supBottomMaxWithSub: 0.344, spaceAfterScript: 0.05,
	upperGap: 0.111, upperRise: 0.2, lowerGap: 0.167, lowerDrop: 0.6,
	stackTop: 0.444, stackTopDisplay: 0.677, stackBottom: 0.345, stackBottomDisp: 0.686,
	stackGap: 0.12, stackGapDisplay: 0.28,
	stretchTop: 0.111, stretchBottom: 0.6, stretchGapAbove: 0.2, stretchGapBel: 0.167,
	numShift: 0.394, numShiftDisplay: 0.677, denShift: 0.345, denShiftDisplay: 0.686,
	numGap: 0.04, numGapDisplay: 0.12, rule: 0.04, denGap: 0.04, denGapDisplay: 0.12,
	skewH: 0.35, skewV: 0.1,
	overGap: 0.12, overRule: 0.04, overExtra: 0.04, underGap: 0.12, underRule: 0.04, underExt: 0.04,
	radGap: 0.05, radGapDisplay: 0.15, radRule: 0.04, radExtra: 0.04,
	radKernBefore: 0.278, radKernAfter: -0.556, radDegreeRaise: 0.6, minOverlap: 0.02,
}

type glyphKey struct {
	face *fontdb.Face
	gid  uint16
}

// glyphMetrics are a glyph's metrics in em.
type glyphMetrics struct {
	adv       float64
	ink       sfnt.Rect // y up; zero for glyphs without an outline
	hasInk    bool
	italic    float64
	topAccent float64
	hasAccent bool
}

// New returns an engine that lays formulas out with f.
func New(f Fonts) *Engine {
	e := &Engine{f: f, c: defaultConsts, glyphs: map[glyphKey]*glyphMetrics{}}
	if f.Set != nil {
		e.face = f.Set.ChooseMath(f.Math)
	}
	if e.face != nil {
		f.Set.Pin(e.face)
		e.mt = e.face.Loaded.Math()
		e.upem = e.face.Loaded.UnitsPerEm()
		e.readConsts()
	} else if f.Warn != nil {
		f.Warn("no font with a MATH table (such as STIX Two Math or Cambria Math) is available; formulas are laid out with the text font")
	}
	return e
}

func (e *Engine) readConsts() {
	m := e.mt.Constants
	u := func(v int) float64 { return float64(v) / e.upem }
	pct := func(v int, d float64) float64 {
		if v <= 0 {
			return d
		}
		return float64(v) / 100
	}
	e.c = consts{
		scriptScale: pct(m.ScriptPercentScaleDown, 0.7), scriptScriptScale: pct(m.ScriptScriptPercentScaleDown, 0.5),
		delimitedMin: u(m.DelimitedSubFormulaMinHeight), displayOpMin: u(m.DisplayOperatorMinHeight),
		axis: u(m.AxisHeight), accentBase: u(m.AccentBaseHeight),
		subShift: u(m.SubscriptShiftDown), subTopMax: u(m.SubscriptTopMax), subDrop: u(m.SubscriptBaselineDropMin),
		supShift: u(m.SuperscriptShiftUp), supShiftCramped: u(m.SuperscriptShiftUpCramped),
		supBottomMin: u(m.SuperscriptBottomMin), supDrop: u(m.SuperscriptBaselineDropMax),
		subSupGap: u(m.SubSuperscriptGapMin), supBottomMaxWithSub: u(m.SuperscriptBottomMaxWithSubscript),
		spaceAfterScript: u(m.SpaceAfterScript),
		upperGap:         u(m.UpperLimitGapMin), upperRise: u(m.UpperLimitBaselineRiseMin),
		lowerGap: u(m.LowerLimitGapMin), lowerDrop: u(m.LowerLimitBaselineDropMin),
		stackTop: u(m.StackTopShiftUp), stackTopDisplay: u(m.StackTopDisplayStyleShiftUp),
		stackBottom: u(m.StackBottomShiftDown), stackBottomDisp: u(m.StackBottomDisplayStyleShiftDown),
		stackGap: u(m.StackGapMin), stackGapDisplay: u(m.StackDisplayStyleGapMin),
		stretchTop: u(m.StretchStackTopShiftUp), stretchBottom: u(m.StretchStackBottomShiftDown),
		stretchGapAbove: u(m.StretchStackGapAboveMin), stretchGapBel: u(m.StretchStackGapBelowMin),
		numShift: u(m.FractionNumeratorShiftUp), numShiftDisplay: u(m.FractionNumeratorDisplayStyleShiftUp),
		denShift: u(m.FractionDenominatorShiftDown), denShiftDisplay: u(m.FractionDenominatorDisplayStyleShiftDown),
		numGap: u(m.FractionNumeratorGapMin), numGapDisplay: u(m.FractionNumDisplayStyleGapMin),
		rule: u(m.FractionRuleThickness), denGap: u(m.FractionDenominatorGapMin), denGapDisplay: u(m.FractionDenomDisplayStyleGapMin),
		skewH: u(m.SkewedFractionHorizontalGap), skewV: u(m.SkewedFractionVerticalGap),
		overGap: u(m.OverbarVerticalGap), overRule: u(m.OverbarRuleThickness), overExtra: u(m.OverbarExtraAscender),
		underGap: u(m.UnderbarVerticalGap), underRule: u(m.UnderbarRuleThickness), underExt: u(m.UnderbarExtraDescender),
		radGap: u(m.RadicalVerticalGap), radGapDisplay: u(m.RadicalDisplayStyleVerticalGap),
		radRule: u(m.RadicalRuleThickness), radExtra: u(m.RadicalExtraAscender),
		radKernBefore: u(m.RadicalKernBeforeDegree), radKernAfter: u(m.RadicalKernAfterDegree),
		radDegreeRaise: pct(m.RadicalDegreeBottomRaisePercent, 0.6),
		minOverlap:     u(e.mt.MinConnectorOverlap),
	}
	if e.c.rule <= 0 {
		e.c.rule = defaultConsts.rule
	}
	if e.c.radRule <= 0 {
		e.c.radRule = e.c.rule
	}
	if e.c.overRule <= 0 {
		e.c.overRule = e.c.rule
	}
	if e.c.underRule <= 0 {
		e.c.underRule = e.c.rule
	}
}

// HasMathFont reports whether formulas are laid out with a formula font.
func (e *Engine) HasMathFont() bool { return e.face != nil }

// XHeight returns the height of the formula font's x in em (0 without a
// formula font).
func (e *Engine) XHeight() float64 {
	if e.face == nil {
		return 0
	}
	return XHeight(e.face)
}

// XHeight returns the height of the x of a face in em (0 when it has
// none).
func XHeight(fc *fontset.Choice) float64 {
	if fc == nil || fc.Loaded == nil {
		return 0
	}
	g, ok := fc.Loaded.Glyph('x')
	if !ok {
		return 0
	}
	r, ok := fc.Loaded.Bounds(g)
	if !ok {
		return 0
	}
	return r.YMax
}

// env is the style a node is laid out in.
type env struct {
	base    float64 // font size at script level 0 (points)
	level   int     // 0 display and text, 1 script, 2 script-script …
	display bool
	cramped bool
	color   bdf.Color
	variant Variant
	bold    bool
}

func (e *Engine) size(v env) float64 {
	switch {
	case v.level <= 0:
		return v.base
	case v.level == 1:
		return v.base * e.c.scriptScale
	}
	return v.base * e.c.scriptScriptScale
}

// em converts a length in em at v's size to points.
func (e *Engine) em(v env, x float64) float64 { return x * e.size(v) }

// script is the style of superscripts; sub that of subscripts.
func (v env) script() env {
	v.level++
	v.display = false
	return v
}

func (v env) sub() env {
	v = v.script()
	v.cramped = true
	return v
}

// num and den are the styles of a fraction's numerator and denominator.
func (v env) num() env {
	if v.display {
		v.display = false
	} else {
		v.level++
	}
	return v
}

func (v env) den() env {
	v = v.num()
	v.cramped = true
	return v
}

func (v env) crampedEnv() env {
	v.cramped = true
	return v
}

// metrics returns the metrics of glyph g of fc in em (cached).
func (e *Engine) metrics(fc *fontset.Choice, g uint16) *glyphMetrics {
	if fc == nil || fc.Loaded == nil {
		return &glyphMetrics{adv: 0.5}
	}
	k := glyphKey{fc.Loaded.Face, g}
	if m, ok := e.glyphs[k]; ok {
		return m
	}
	l := fc.Loaded
	m := &glyphMetrics{adv: l.Advance(g)}
	m.ink, m.hasInk = l.Bounds(g)
	if fc == e.face && e.mt != nil {
		if v, ok := e.mt.Italic(g); ok {
			m.italic = float64(v) / e.upem
		}
		if v, ok := e.mt.TopAccent(g); ok {
			m.topAccent, m.hasAccent = float64(v)/e.upem, true
		}
	}
	e.glyphs[k] = m
	return m
}

// glyph is a character resolved to the face that draws it.
type glyph struct {
	r   rune
	fc  *fontset.Choice
	gid uint16
}

// resolve picks the face that draws r: the formula font when it has r,
// else a text font (bold or italic as the variant asks, which the text
// font draws for itself).
func (e *Engine) resolve(r rune, family string, bold, italic bool) glyph {
	if family == "" && e.face != nil {
		if g, ok := e.face.Loaded.Glyph(r); ok {
			return glyph{r: r, fc: e.face, gid: g}
		}
	}
	var fc *fontset.Choice
	if e.f.Text != nil {
		fc = e.f.Text(family, r, bold, italic)
	}
	gl := glyph{r: r, fc: fc}
	if fc != nil && fc.Loaded != nil {
		gl.gid, _ = fc.Loaded.Glyph(r)
	}
	return gl
}

// char resolves a character of a token in variant v: the styled form
// (𝑥 for an italic x) when the formula font has it, else the plain
// character in a text font of that style.
func (e *Engine) char(r rune, v Variant, family string) glyph {
	if family == "" {
		if s := styled(r, v); s != r {
			if e.face != nil && e.face.Loaded.Has(s) {
				return glyph{r: s, fc: e.face, gid: mustGlyph(e.face, s)}
			}
		}
		if e.face != nil && e.face.Loaded.Has(r) && !v.isItalic() && !v.isBold() {
			return e.resolve(r, "", false, false)
		}
	}
	return e.resolve(r, family, v.isBold(), v.isItalic())
}

func mustGlyph(fc *fontset.Choice, r rune) uint16 {
	g, _ := fc.Loaded.Glyph(r)
	return g
}

// record measures a character for the embedded subset and returns its
// advance in em.
func (e *Engine) record(gl glyph) float64 {
	if e.f.Set == nil || gl.fc == nil {
		return 0.5
	}
	return e.f.Set.Advance(gl.fc, gl.r)
}

// glyphRune returns the character that draws glyph g of the formula font:
// the character itself when it maps to g, else a private use one.
func (e *Engine) glyphRune(r rune, g uint16) rune {
	if e.face != nil {
		if cg, ok := e.face.Loaded.Glyph(r); ok && cg == g {
			e.f.Set.Advance(e.face, r)
			return r
		}
		if pr, ok := e.f.Set.GlyphRune(e.face, g); ok {
			return pr
		}
	}
	return r
}

// singleRune returns the only character of s.
func singleRune(s string) (rune, bool) {
	r, n := utf8.DecodeRuneInString(s)
	return r, n > 0 && n == len(s)
}

// isLetter reports whether r takes an identifier's italic.
func isLetter(r rune) bool {
	return unicode.IsLetter(r) && r < 0x2000 || r == 'ı' || r == 'ȷ'
}

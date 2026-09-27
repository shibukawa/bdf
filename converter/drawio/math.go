package drawio

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/equation"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Mathematical typesetting: a diagram whose model has math="1" has the
// LaTeX in its labels typeset by MathJax, between $$ … $$ and \[ … \]
// (display) and \( … \) (inline). The formula engine
// (converter/internal/equation) lays each formula out when the label's
// text is built, and it becomes an item of the label's line as wide as the
// formula, which reaches its height and depth; display formulas take lines
// of their own in HTML labels. AsciiMath between backquotes stays text.

// formula is a laid out formula in a label's line.
type formula struct {
	box  *equation.Box
	text string // the linear notation, for extraction
}

// mathDelims are the delimiters of formulas, display ones first.
var mathDelims = []struct {
	open, close string
	display     bool
}{
	{"$$", "$$", true},
	{`\[`, `\]`, true},
	{`\(`, `\)`, false},
}

// nextFormula finds the first formula in s: where it starts and ends (its
// delimiters included), whether it is a display formula, and its LaTeX.
// start is -1 when s holds none.
func nextFormula(s string) (start, end int, display bool, tex string) {
	start = -1
	for _, d := range mathDelims {
		i := strings.Index(s, d.open)
		if i < 0 || start >= 0 && i >= start {
			continue
		}
		j := strings.Index(s[i+len(d.open):], d.close)
		if j < 0 {
			continue
		}
		start, end, display = i, i+len(d.open)+j+len(d.close), d.display
		tex = s[i+len(d.open) : i+len(d.open)+j]
	}
	return start, end, display, tex
}

// formula adds a formula in the style of the text around it.
func (b *textBuilder) formula(tex string, display bool, st *tstyle) {
	c := b.c
	n := equation.ParseTeX(tex)
	eng := c.mathEngine()
	size := st.size
	// MathJax sets formulas as large as they must be for their x to be as
	// tall as the text's
	if mx, tx := eng.XHeight(), equation.XHeight(st.primary); mx > 0 && tx > 0 {
		size *= math.Min(math.Max(tx/mx, 1), 1.25)
	}
	c.mathStyle = st
	box := eng.Layout(n, equation.Style{Size: size, Display: display, Color: st.color.withAlpha(b.alpha).bdf(), Bold: st.bold})
	c.mathStyle = nil
	ownLine := display && !b.inlineOnly
	if ownLine && b.cur != nil && lineHasContent(b.cur.items) {
		b.hardBreak(st)
	}
	p := b.para(st)
	fc := c.fonts.FaceForFamilies(st.families, st.bold, st.italic, 'x')
	p.items = append(p.items, titem{r: '￼', st: st, fc: fc, w: box.Width(), eq: &formula{box: box, text: equation.Linear(n)}})
	b.space = false
	if ownLine {
		b.hardBreak(st)
	}
}

// lineHasContent reports whether items hold something visible after their
// last forced break.
func lineHasContent(items []titem) bool {
	for i := len(items) - 1; i >= 0; i-- {
		switch {
		case items[i].hard:
			return false
		case !items[i].space:
			return true
		}
	}
	return false
}

// mathEngine returns the engine that lays the diagram's formulas out.
func (c *converter) mathEngine() *equation.Engine {
	if c.eq == nil {
		c.eq = equation.New(equation.Fonts{Set: c.fonts, Text: c.mathText, Warn: func(m string) { c.warnf("%s", m) }})
	}
	return c.eq
}

// mathText picks the face for the characters of a formula the formula font
// lacks, and for text in a family: the fonts of the label's text.
func (c *converter) mathText(family string, r rune, bold, italic bool) *fontset.Choice {
	families := []string{"Helvetica"}
	if st := c.mathStyle; st != nil && len(st.families) > 0 {
		families = st.families
	}
	if family != "" {
		families = append([]string{family}, families...)
	}
	return c.fonts.FaceForFamilies(families, bold, italic, r)
}

// emitFormula draws a formula of a line whose baseline is at base.
func (e *textEmitter) emitFormula(obj *bdf.Object, it titem, x, base float64) {
	// the text layer sizes the formula's run with the font in effect
	if ref := e.cv.Font(it.fc.Use); !e.set || ref != e.font || it.st.size != e.size {
		obj.Font(ref, f32(it.st.size))
		e.font, e.size, e.set = ref, it.st.size, true
	}
	y := base + it.st.shift
	equation.Place(e.cv, it.eq.box, x, y, it.eq.text)
	if it.st.link != "" {
		if u := e.c.linkURL(it.st.link); u != "" {
			obj.Link(f32(x), f32(y-it.eq.box.Height()), f32(it.w), f32(it.eq.box.Height()+it.eq.box.Depth()), u)
		}
	}
}

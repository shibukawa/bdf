package drawingml

import (
	"strings"

	"github.com/shibukawa/bdf/converter/internal/equation"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// PowerPoint and Excel write equations as Office Math in a14:m elements
// among the runs of a paragraph (in the mc:Choice that ooxml.MathChoice
// picks over the picture of the fallback). The formula engine
// (converter/internal/equation) lays them out when the paragraph is read,
// and each becomes an item of the line as wide as the formula that reaches
// its height and depth; m:oMathPara sets its formulas on lines of their
// own.

// formula is a laid out formula in a line.
type formula struct {
	box  *equation.Box
	text string // the linear notation, for extraction
}

// mathEngine returns the engine that lays the document's formulas out.
func (c *Renderer) mathEngine() *equation.Engine {
	if c.eq == nil {
		c.eq = equation.New(equation.Fonts{Set: c.fonts, Math: "Cambria Math", Text: c.mathText, Warn: c.warn})
	}
	return c.eq
}

// mathText picks the face for the characters of a formula the formula font
// lacks, and for normal text in a family: the fonts of the formula's run.
func (c *Renderer) mathText(family string, r rune, bold, italic bool) *fontset.Choice {
	var latin, ea string
	if st := c.mathRun; st != nil {
		latin, ea = st.latin, st.ea
	}
	if family != "" {
		latin = family
	}
	return c.fonts.FaceFor(latin, ea, bold, italic, r)
}

// addMath reads an a14:m element of a paragraph whose paragraph properties
// are pc; rest is the content of the paragraph after it.
func (s *Drawing) addMath(tf *textFrame, pa *para, pc chain, fontScale float64, m *ooxml.Node, rest []*ooxml.Node) {
	reader := &equation.OMML{RunStyle: func(rPr *ooxml.Node) (*equation.Styled, string, bool, bool) {
		st := s.runStyle(tf, tf.rChain(rPr, pc), fontScale, rPr)
		sty := &equation.Styled{Base: st.size}
		if c, ok := st.color(); ok {
			col := c.bdf()
			sty.Color = &col
		}
		return sty, st.latin, st.bold, st.italic
	}}
	// the formula is as large as its first run; what no run colors (the
	// delimiters, fraction bars) takes the paragraph's color
	rPr := firstMathRPr(m)
	st := s.runStyle(tf, tf.rChain(rPr, pc), fontScale, rPr)
	if def := s.runStyle(tf, tf.rChain(nil, pc), fontScale, nil); rPr.Child("solidFill") != nil || rPr.Child("gradFill") != nil {
		c := *st
		c.fill = def.fill
		c.setKey()
		st = &c
	}
	for _, k := range m.Elements() {
		switch k.Name {
		case "oMath":
			s.addFormula(pa, reader.ParseOMML(k), st, false)
		case "oMathPara":
			lines, jc := reader.ParseOMMLPara(k)
			if len(lines) == 0 {
				continue
			}
			more := hasContent(rest)
			if !more && !pa.hasItems() {
				// the paragraph is the formula: it is aligned as the
				// formula says
				pa.algn = map[string]string{"left": "l", "right": "r"}[jc]
				if pa.algn == "" {
					pa.algn = "ctr"
				}
			} else if len(pa.items) > 0 {
				s.addBreak(pa, st)
			}
			for i, l := range lines {
				if i > 0 {
					s.addBreak(pa, st)
				}
				s.addFormula(pa, l, st, true)
			}
			if more {
				s.addBreak(pa, st)
			}
		}
	}
}

func (s *Drawing) addBreak(pa *para, st *runStyle) {
	pa.items = append(pa.items, item{r: '\n', st: st, fc: s.c.faceFor(st, ' '), kind: itemBreak})
}

// addFormula lays a formula out in the style of the run st and adds it to
// the paragraph.
func (s *Drawing) addFormula(pa *para, n equation.Node, st *runStyle, display bool) {
	c := s.c
	eng := c.mathEngine()
	style := equation.Style{Size: st.size, Display: display}
	if col, ok := st.color(); ok {
		style.Color = col.bdf()
	}
	c.mathRun = st
	b := eng.Layout(n, style)
	c.mathRun = nil
	pa.items = append(pa.items, item{kind: itemMath, st: st, fc: c.faceFor(st, 'x'), w: b.Width(), brk: true,
		eq: &formula{box: b, text: equation.Linear(n)}})
}

// hasItems reports whether a paragraph holds anything but white space.
func (pa *para) hasItems() bool {
	for _, it := range pa.items {
		if it.kind != itemChar || !isBreakSpace(it.r) {
			return true
		}
	}
	return false
}

// hasContent reports whether the rest of a paragraph holds text or
// formulas.
func hasContent(rest []*ooxml.Node) bool {
	for _, k := range rest {
		switch k.Name {
		case "r", "fld":
			if strings.TrimSpace(k.Child("t").Content()) != "" || k.Name == "fld" {
				return true
			}
		case "br", "m":
			return true
		}
	}
	return false
}

// firstMathRPr returns the run properties (a:rPr) of the first run of a
// formula, which set its size and color; nil when it has none.
func firstMathRPr(n *ooxml.Node) *ooxml.Node {
	for _, k := range n.Elements() {
		if k.Name == "r" && strings.HasSuffix(k.Space, "/math") {
			for _, r := range k.Elements() {
				if r.Name == "rPr" && strings.Contains(r.Space, "drawingml") {
					return r
				}
			}
			continue
		}
		if r := firstMathRPr(k); r != nil {
			return r
		}
	}
	return nil
}

// emitMath draws a formula of a line whose baseline is at base.
func (e *textEmitter) emitMath(it item, dx, base float64) {
	x := it.x + dx
	// the text layer sizes the formula's run with the font in effect
	e.setFont(it.fc, it.st.size)
	equation.Place(e.cv, it.eq.box, x, base, it.eq.text)
	if it.st.link != "" {
		e.link(x, base-it.eq.box.Height(), x+it.w, base+it.eq.box.Depth(), it.st.link)
	}
}

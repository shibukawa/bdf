package wordproc

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/equation"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// Formulas (Office Math in Word documents, MathML and LaTeX in HTML and
// Markdown) are laid out by the formula engine (converter/internal/equation)
// when they are read, and become inline objects of their paragraphs: an
// object as wide as the formula that reaches its height above the baseline
// and its depth below.

// mathEngine returns the engine that lays the document's formulas out.
func (c *converter) mathEngine() *equation.Engine {
	if c.eq == nil {
		c.eq = equation.New(equation.Fonts{Set: c.fonts, Math: c.mathFont, Text: c.mathText, Warn: c.warn})
	}
	return c.eq
}

// mathText picks the face for the characters of a formula that the formula
// font lacks (East Asian text) and for text in a text font: the fonts of
// the text around the formula.
func (c *converter) mathText(family string, r rune, bold, italic bool) *fontset.Choice {
	var latin, ea string
	if st := c.mathRun; st != nil {
		latin, ea = st.fonts[0], st.fonts[2]
		if r >= 0x80 {
			latin = st.fonts[1]
		}
	}
	if family != "" {
		latin = family
	}
	return c.fonts.FaceFor(latin, ea, bold, italic, r)
}

// addFormula lays a formula out in the style of the run st and adds it to
// p as an inline object.
func (c *converter) addFormula(p *para, n equation.Node, st *runStyle, display bool) {
	c.mathRun = st
	eng := c.mathEngine()
	size := st.size
	if c.css {
		// a web page's formulas are as large as they must be for their
		// x to be as tall as the text's (as MathJax sets them)
		if mx, tx := eng.XHeight(), equation.XHeight(c.face(st, 'x')); mx > 0 && tx > 0 {
			size *= math.Min(math.Max(tx/mx, 1), 1.25)
		}
	}
	b := eng.Layout(n, equation.Style{Size: size, Display: display, Color: st.color})
	c.mathRun = nil
	text := equation.Linear(n)
	h := b.Height()
	o := &inlineObj{w: b.Width(), h: h + b.Depth(), desc: b.Depth(), central: c.css, paint: func(e *emitter, box drawingml.Box) {
		equation.Place(e.cv, b, box.X, box.Y+h, text)
	}}
	p.items = append(p.items, item{kind: kObject, obj: o, st: st, w: o.w})
	if st.link != "" {
		o.link = st.link
	}
}

// middle returns how far above the baseline the middle of the paragraph's
// character squares is, where upright characters of vertical text are
// centered.
func (p *para) middle() float64 {
	if p.markFace == nil {
		return 0
	}
	return (p.markFace.Asc - p.markFace.Desc) * p.mark.size / 2
}

// omml returns the reader of the document's Office Math.
func (c *converter) omml() *equation.OMML {
	if c.ommlReader == nil {
		c.ommlReader = &equation.OMML{}
		if s := c.settings(); s != nil {
			pr := s.Child("mathPr")
			c.ommlReader.NaryLim = pr.Path("naryLim").AttrStr("val", "")
			c.ommlReader.IntLim = pr.Path("intLim").AttrStr("val", "")
		}
	}
	return c.ommlReader
}

// math reads an Office Math zone in a paragraph (m:oMath, inline) or a
// display math paragraph (m:oMathPara, a line of its own for each of its
// formulas).
func (w *walker) math(p *para, n *ooxml.Node, link string, more bool) {
	st := w.mathStyle(p, n, link)
	if n.Name == "oMath" {
		w.c.addFormula(p, w.c.omml().ParseOMML(n), st, false)
		return
	}
	lines, jc := w.c.omml().ParseOMMLPara(n)
	if len(lines) == 0 {
		return
	}
	alone := !more
	for _, it := range p.items {
		if it.kind != kChar || !lbSpace(it.r) {
			alone = false
		}
	}
	if alone {
		// the paragraph is the formula: it is justified as the formula
		// says, whatever the paragraph says
		pp := *p.pp
		pp.jc = jc
		p.pp = &pp
		p.items = p.items[:0]
	} else if len(p.items) > 0 {
		p.items = append(p.items, item{kind: kBreak, r: '\n', st: st})
	}
	for i, l := range lines {
		if i > 0 {
			p.items = append(p.items, item{kind: kBreak, r: '\n', st: st})
		}
		w.c.addFormula(p, l, st, true)
	}
	if more {
		p.items = append(p.items, item{kind: kBreak, r: '\n', st: st})
	}
}

// mathStyle is the style of a formula's text: that of its first run.
func (w *walker) mathStyle(p *para, n *ooxml.Node, link string) *runStyle {
	var rPr *ooxml.Node
	var find func(n *ooxml.Node) bool
	find = func(n *ooxml.Node) bool {
		for _, k := range n.Elements() {
			if k.Name == "r" {
				for _, r := range k.Elements() {
					if r.Name == "rPr" && strings.Contains(r.Space, "wordprocessingml") {
						rPr = r
						return true
					}
				}
				continue
			}
			if find(k) {
				return true
			}
		}
		return false
	}
	find(n)
	rp := w.c.runProps(p.style, val(rPr.Child("rStyle")), w.tl, rPr)
	return w.c.runStyle(rp, w.linkOf(link), w.dark)
}

// moreContent reports whether elements after the i-th of a paragraph's
// content hold text or drawings.
func moreContent(kids []*ooxml.Node, i int) bool {
	for _, k := range kids[i+1:] {
		switch k.Name {
		case "r", "hyperlink", "fldSimple", "smartTag", "customXml", "ins", "sdt", "oMath", "oMathPara":
			return true
		}
	}
	return false
}

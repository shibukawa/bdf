package font

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/otlayout"
)

// The glyph grid.
const (
	glyphCols  = 10
	glyphCellW = contentW / glyphCols
	glyphCellH = 82.0
	glyphBoxH  = 54.0 // the part of a cell the glyph is drawn in
)

// glyphView adds the rows of the glyphs view.
func (c *converter) glyphView(s *scroll) {
	fc := c.fc
	n := fc.f.NumGlyphs
	s.add(30, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "1")
		p.label("Glyphs", margin, y+22, 18, true, colText, 0)
	})
	named := 0
	for _, nm := range fc.glyphs {
		if nm != "" {
			named++
		}
	}
	info := num(n) + " glyphs"
	switch {
	case named == 0:
		info += ", without names"
	case named < n:
		info += fmt.Sprintf(", %s of them named", num(named))
	}
	mapped := 0
	for _, cs := range fc.chars {
		if len(cs) > 0 {
			mapped++
		}
	}
	info += fmt.Sprintf("; %s mapped from characters", num(mapped))
	s.add(18, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkParagraph, "")
		p.label(info, margin, y+12, bodySize, false, colGray, 0)
	})
	s.space(10)
	for first := 0; first < n; first += glyphCols {
		s.add(glyphCellH, func(p *pen, y float64) {
			for k := range glyphCols {
				g := first + k
				if g >= n {
					break
				}
				c.glyphCell(p, uint16(g), margin+float64(k)*glyphCellW, y)
			}
			p.rect(margin, y, contentW, 0.5, colRule)
			p.rect(margin, y+glyphCellH-0.5, contentW, 0.5, colRule)
			for k := 0; k <= glyphCols; k++ {
				p.rect(margin+float64(k)*glyphCellW-0.25, y, 0.5, glyphCellH, colRule)
			}
		})
	}
}

// glyphCell draws a cell of the glyph grid: the glyph between its origin
// and advance, its ID and GDEF class above, its name and characters below.
func (c *converter) glyphCell(p *pen, g uint16, x, y float64) {
	fc := c.fc
	p.obj().Mark(bdf.MarkParagraph, "")
	p.label(strconv.Itoa(int(g)), x+3, y+9, tinySize, true, colGray, 0)
	if fc.gdef != nil {
		cls := ""
		switch fc.gdef.Classes[g] {
		case otlayout.GlyphLigature:
			cls = "ligature"
		case otlayout.GlyphMark:
			cls = "mark"
		case otlayout.GlyphComponent:
			cls = "component"
		}
		if cls != "" {
			p.label(cls, x+glyphCellW-3, y+9, tinySize-0.5, false, colAccent, 2)
		}
	}
	// the advance, as guides on the baseline
	asc, desc := fc.vertical()
	size := (glyphBoxH - 14) / (asc + desc)
	b := fc.bounds[g]
	if b.ok {
		size = min(size, (glyphCellW-10)/max(b.x1-b.x0, fc.advances[g], 0.01))
	} else if a := fc.advances[g]; a > 0 {
		size = min(size, (glyphCellW-10)/a)
	}
	top := y + 14
	base := top + 3 + asc*size
	adv := fc.advances[g] * size
	ox := x + glyphCellW/2 - adv/2
	if b.ok && (b.x0*size < -(glyphCellW-10)/2+adv/2 || b.x1*size > adv/2+(glyphCellW-10)/2) {
		// ink far outside the advance (marks of zero width): center the ink
		ox = x + glyphCellW/2 - (b.x0+b.x1)/2*size
	}
	p.rect(ox, base, adv, 0.5, colGuide)
	p.rect(ox-0.25, base-asc*size, 0.5, (asc+desc)*size, colGuide)
	p.rect(ox+adv-0.25, base-asc*size, 0.5, (asc+desc)*size, colGuide)
	alt := ""
	if cs := fc.chars[g]; len(cs) > 0 && !unicode.IsControl(cs[0]) {
		alt = string(cs[0])
	}
	p.glyph(g, ox, base, size, alt)
	// name and characters
	ly := y + glyphBoxH + 15
	if nm := fc.glyphs[g]; nm != "" {
		p.label(c.ui.fit(nm, tinySize, false, glyphCellW-6), x+glyphCellW/2, ly, tinySize, false, colText, 1)
	}
	if cl := charLabel(append([]rune(nil), fc.chars[g]...)); cl != "" {
		p.label(cl, x+glyphCellW/2, ly+9, tinySize, false, colAccent, 1)
	}
}

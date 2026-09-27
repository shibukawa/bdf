package font

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
)

type block struct {
	first, last rune
	name        string
}

// blockCoverage is how much of a block the font maps.
type blockCoverage struct {
	block
	mapped, assigned int
}

// assigned reports whether a code point is assigned a character a font
// may draw (in Go's Unicode version, which the block table follows):
// control characters and surrogates are left out.
func assigned(r rune) bool {
	return unicode.In(r, unicode.L, unicode.M, unicode.N, unicode.P, unicode.S, unicode.Z, unicode.Cf, unicode.Co)
}

// coverage returns the blocks the font maps characters of, in order, with
// a block for mapped code points outside every block.
func (fc *face) coverage() []blockCoverage {
	var out []blockCoverage
	i := 0
	var stray []rune
	for _, b := range blocks {
		for i < len(fc.runes) && fc.runes[i] < b.first {
			stray = append(stray, fc.runes[i])
			i++
		}
		n := 0
		for i < len(fc.runes) && fc.runes[i] <= b.last {
			n++
			i++
		}
		if n == 0 {
			continue
		}
		bc := blockCoverage{block: b, mapped: n}
		for r := b.first; r <= b.last; r++ {
			if assigned(r) {
				bc.assigned++
			}
		}
		out = append(out, bc)
	}
	stray = append(stray, fc.runes[i:]...)
	if len(stray) > 0 {
		out = append(out, blockCoverage{block: block{stray[0], stray[len(stray)-1], "No block"}, mapped: len(stray)})
	}
	return out
}

// hex formats a code point as at least four hex digits.
func hex(r rune) string {
	s := strings.ToUpper(strconv.FormatInt(int64(r), 16))
	for len(s) < 4 {
		s = "0" + s
	}
	return s
}

// The code charts: a column of row labels and 16 cells.
const (
	rowLabelW  = 48.0
	chartCellW = (contentW - rowLabelW) / 16
	chartCellH = 44.0
)

// characters adds the rows of the characters view.
func (c *converter) characters(s *scroll) {
	fc := c.fc
	cov := fc.coverage()
	s.add(30, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "1")
		p.label("Characters", margin, y+22, 18, true, colText, 0)
	})
	summary := fmt.Sprintf("%s characters in %d Unicode blocks", num(len(fc.runes)), len(cov))
	if n := variationSequences(fc.f.Tables["cmap"]); n > 0 {
		summary += fmt.Sprintf(", and %s variation sequences", num(n))
	}
	s.add(18, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkParagraph, "")
		p.label(summary, margin, y+12, bodySize, false, colGray, 0)
	})
	// the blocks, with links to their charts
	c.heading(s, "Unicode blocks")
	const barW = 110.0
	for i, b := range cov {
		s.add(15, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			anchor := "block" + strconv.Itoa(i)
			w := p.label(c.ui.fit(b.name, bodySize, false, 250), margin, y+11, bodySize, false, colAccent, 0)
			p.linkTo(margin, y, w, 15, anchor)
			rng := "U+" + hex(b.first) + "–" + hex(b.last)
			p.label(rng, margin+260, y+11, smallSize+0.5, false, colGray, 0)
			frac := 1.0
			of := num(b.mapped)
			if b.assigned > 0 {
				frac = min(1, float64(b.mapped)/float64(b.assigned))
				of = num(b.mapped) + " of " + num(b.assigned)
			}
			p.label(of, margin+contentW-barW-10, y+11, smallSize+0.5, false, colText, 2)
			bx := margin + contentW - barW
			p.rect(bx, y+4, barW, 7, colEmpty)
			p.rect(bx, y+4, barW*frac, 7, colAccent)
		})
	}
	// the charts, up to maxChartChars characters
	shown := 0
	for i, b := range cov {
		if shown+b.mapped > maxChartChars && shown > 0 {
			more := fmt.Sprintf("The charts of the %d blocks after these are left out: the font maps more than %s characters.", len(cov)-i, num(maxChartChars))
			s.space(16)
			s.add(16, func(p *pen, y float64) {
				p.obj().Mark(bdf.MarkParagraph, "")
				p.label(more, margin, y+11, bodySize, false, colGray, 0)
			})
			break
		}
		c.chart(s, b, "block"+strconv.Itoa(i))
		shown += b.mapped
	}
}

// maxChartChars bounds the characters of the code charts: more than any
// font of real text maps, fewer than a font that maps every code point
// (the Last Resort font) would make of the document.
const maxChartChars = 150000

// chart adds the code chart of a block: the rows of 16 code points that
// hold a mapped character.
func (c *converter) chart(s *scroll, b blockCoverage, anchor string) {
	fc := c.fc
	s.space(16)
	s.anchor(anchor)
	title := b.name
	info := "U+" + hex(b.first) + "–" + hex(b.last) + " · " + num(b.mapped)
	if b.assigned > 0 {
		info += " of " + num(b.assigned)
	}
	h := s.add(34, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "2")
		w := p.label(title, margin, y+16, 11.5, true, colText, 0)
		p.obj().Mark(bdf.MarkParagraph, "")
		p.label(info, margin+w+10, y+16, smallSize+0.5, false, colGray, 0)
		// column headings
		for k := range 16 {
			x := margin + rowLabelW + (float64(k)+0.5)*chartCellW
			p.label(strings.ToUpper(strconv.FormatInt(int64(k), 16)), x, y+30, smallSize, false, colGray, 1)
		}
	})
	h.keep = true
	first := b.first &^ 15
	gap := false
	for row := first; row <= b.last; row += 16 {
		any := false
		for r := row; r < row+16; r++ {
			if fc.glyph(r) != 0 && r >= b.first && r <= b.last {
				any = true
				break
			}
		}
		if !any {
			if !gap {
				gap = true
				s.add(10, func(p *pen, y float64) {
					p.label("⋯", margin+rowLabelW/2, y+8, smallSize, false, colGray, 1)
				})
			}
			continue
		}
		gap = false
		s.add(chartCellH, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			p.label(hex(row)[:len(hex(row))-1]+"x", margin+rowLabelW-8, y+chartCellH/2+3, smallSize, false, colGray, 2)
			for k := range 16 {
				r := row + rune(k)
				x := margin + rowLabelW + float64(k)*chartCellW
				switch g := fc.glyph(r); {
				case r < b.first || r > b.last:
				case g != 0:
					c.cellGlyph(p, g, x, y, chartCellW, chartCellH, string(r))
				case assigned(r):
					p.rect(x, y, chartCellW, chartCellH, colEmpty)
				default:
					p.rect(x, y, chartCellW, chartCellH, colVoid)
				}
			}
			// grid
			p.rect(margin+rowLabelW, y, contentW-rowLabelW, 0.5, colRule)
			p.rect(margin+rowLabelW, y+chartCellH-0.5, contentW-rowLabelW, 0.5, colRule)
			for k := 0; k <= 16; k++ {
				p.rect(margin+rowLabelW+float64(k)*chartCellW-0.25, y, 0.5, chartCellH, colRule)
			}
		})
	}
}

// cellGlyph draws a glyph in the middle of a cell at the size the cells
// share, smaller when its ink would not fit in the cell.
func (c *converter) cellGlyph(p *pen, g uint16, x, y, w, h float64, alt string) {
	fc := c.fc
	asc, desc := fc.vertical()
	size := h * 0.62 / (asc + desc)
	base := y + (h-(asc+desc)*size)/2 + asc*size
	b := fc.bounds[g]
	ox := x + w/2 - fc.advances[g]*size/2
	if b.ok {
		if fw, fh := w*0.9/max(b.x1-b.x0, 1e-3), h*0.92/max(b.y1-b.y0, 1e-3); fw < size || fh < size {
			// too large for the cell: shrunk and centered on its ink
			size = min(fw, fh)
			base = y + h/2 + (b.y1+b.y0)/2*size
		} else if top, bottom := base-b.y1*size, base-b.y0*size; top < y+1 || bottom > y+h-1 {
			// beyond the cell above or below: moved in
			base += max(y+1-top, 0) - max(bottom-(y+h-1), 0)
		}
		ox = x + w/2 - (b.x0+b.x1)/2*size
	}
	p.glyph(g, ox, base, size, alt)
}

// vertical returns the ascent and descent of the font in em, both
// positive.
func (fc *face) vertical() (asc, desc float64) {
	asc, desc = float64(fc.f.Ascent)/fc.upem, -float64(fc.f.Descent)/fc.upem
	if os2 := fc.f.Tables["OS/2"]; len(os2) >= 72 && (asc+desc <= 0 || asc <= 0) {
		asc, desc = float64(s16(os2, 68))/fc.upem, -float64(s16(os2, 70))/fc.upem
	}
	if asc+desc <= 0.2 || asc <= 0 {
		asc, desc = 0.8, 0.2
	}
	return asc, desc
}

// charLabel names the characters a glyph is mapped from.
func charLabel(rs []rune) string {
	if len(rs) == 0 {
		return ""
	}
	sort.Slice(rs, func(i, j int) bool { return rs[i] < rs[j] })
	s := "U+" + hex(rs[0])
	if len(rs) > 1 {
		s += " +" + strconv.Itoa(len(rs)-1)
	}
	return s
}

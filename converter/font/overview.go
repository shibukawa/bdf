package font

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/otlayout"
)

// Sizes of the views' own text, in pt.
const (
	titleSize   = 22.0
	headingSize = 13.0
	bodySize    = 9.5
	smallSize   = 7.5
	tinySize    = 6.5
)

// waterfallSizes are the sizes the sample text is shown at.
var waterfallSizes = []float64{8, 9, 10, 12, 14, 18, 24, 30, 36, 48, 60, 72}

// heading adds a section heading.
func (c *converter) heading(s *scroll, text string) {
	s.space(18)
	b := s.add(26, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "2")
		p.label(text, margin, y+16, headingSize, true, colText, 0)
		p.rect(margin, y+22, contentW, 0.75, colText)
	})
	b.keep = true
	s.space(4)
}

// overview adds the rows of the overview.
func (c *converter) overview(s *scroll) {
	fc := c.fc
	name := fc.fullName()
	if name == "" {
		name = "Untitled font"
	}
	s.add(34, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "1")
		p.label(name, margin, y+24, titleSize, true, colText, 0)
	})
	var facts []string
	if v := fc.name(5); v != "" {
		facts = append(facts, v)
	}
	facts = append(facts, num(fc.f.NumGlyphs)+" glyphs", num(len(fc.runes))+" characters", fc.outlineFormat())
	if c.fl.container != "" {
		facts = append(facts, c.fl.container)
	}
	sub := strings.Join(facts, " · ")
	s.add(18, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkParagraph, "")
		p.label(c.ui.fit(sub, bodySize, false, contentW), margin, y+12, bodySize, false, colGray, 0)
	})
	if c.paint.glyphFont == nil {
		s.add(16, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			p.label("The font's license does not allow embedding it: its glyphs are drawn as outlines, without hinting or shaping.", margin, y+11, smallSize, false, colMissing, 0)
		})
	}
	s.space(10)
	c.specimen(s)
	c.waterfall(s)
	c.samples(s)
	for _, sec := range fc.sections(c.fl, c.paint.glyphFont == nil) {
		c.table(s, sec)
	}
	c.scriptsTable(s)
	c.paletteRows(s)
	c.tablesGrid(s)
}

// tablesGrid adds the tables of the font in columns.
func (c *converter) tablesGrid(s *scroll) {
	rows := c.fc.tableList()
	if len(rows) == 0 {
		return
	}
	c.heading(s, "Tables")
	const cols = 4
	colW := contentW / cols
	for i := 0; i < len(rows); i += cols {
		line := rows[i:min(i+cols, len(rows))]
		s.add(15, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			for k, r := range line {
				if k > 0 {
					p.obj().Mark(bdf.MarkLine, "")
				}
				x := margin + float64(k)*colW
				p.label(r.label, x, y+11, bodySize, true, colText, 0)
				p.label(r.value, x+colW-16, y+11, smallSize+0.5, false, colGray, 2)
			}
			p.rect(margin, y+14.5, contentW, 0.5, colRule)
		})
	}
}

// specimen adds the characters shown large, with lines of the character
// set beside them.
func (c *converter) specimen(s *scroll) {
	pl := c.plan
	if pl.specimen == "" {
		return
	}
	const big = 120.0
	fc := c.fc
	w := 0.0
	for _, r := range pl.specimen {
		w += fc.advances[fc.glyph(r)] * big
	}
	w = min(w, contentW*0.55)
	const h = 160.0
	s.add(h, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkFigure, "")
		p.obj().Save()
		p.obj().ClipRect(f32(margin), f32(y), f32(contentW), f32(h))
		base := y + 0.72*h
		if pl.rtl {
			p.sample(pl.specimen, margin+w, base, big, true, pl.lang)
		} else {
			p.sample(pl.specimen, margin, base, big, false, pl.lang)
		}
		// beside it, lines of the character set
		x := margin + w + 28
		size := 20.0
		if n := len(pl.charset); n > 0 {
			size = min(22, (h-10)/float64(n)/1.35)
		}
		ly := y + 8
		for _, l := range pl.charset {
			ly += size * 1.3
			if pl.rtl {
				p.sample(l, margin+contentW, ly, size, true, pl.lang)
			} else {
				p.sample(l, x, ly, size, false, pl.lang)
			}
		}
		p.obj().Restore()
		p.set, p.fset = false, false
		p.obj().Mark(bdf.MarkEnd, "")
	})
}

// waterfall adds the sample text at each size.
func (c *converter) waterfall(s *scroll) {
	pl := c.plan
	if pl.waterfall == "" {
		return
	}
	c.heading(s, "Sizes (pt)")
	for _, size := range waterfallSizes {
		h := size*1.3 + 6
		s.add(h, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			base := y + size*1.02 + 3
			p.label(strconv.FormatFloat(size, 'f', -1, 64), margin, base, smallSize, false, colGray, 0)
			p.obj().Save()
			p.obj().ClipRect(f32(margin+30), f32(y), f32(contentW-30), f32(h))
			if pl.rtl {
				p.sample(pl.waterfall, margin+contentW, base, size, true, pl.lang)
			} else {
				p.sample(pl.waterfall, margin+30, base, size, false, pl.lang)
			}
			p.obj().Restore()
			p.set, p.fset = false, false
		})
	}
}

// samples adds the sample text of each script the font covers.
func (c *converter) samples(s *scroll) {
	if len(c.plan.scripts) == 0 {
		return
	}
	c.heading(s, "Sample text")
	fc := c.fc
	const size = 18.0
	for _, sc := range c.plan.scripts {
		s.add(15, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			p.label(sc.name, margin, y+11, smallSize, false, colGray, 0)
		})
		// greedy lines at spaces from the advances, with room for shaping
		var lines []string
		words := strings.Fields(sc.text)
		cur, w := "", 0.0
		for _, wd := range words {
			ww := 0.0
			for _, r := range " " + wd {
				ww += fc.advances[fc.glyph(r)] * size
			}
			if cur != "" && w+ww > contentW*0.9 {
				lines = append(lines, cur)
				cur, w = "", 0
			}
			if cur != "" {
				cur += " "
			}
			cur += wd
			w += ww
		}
		if cur != "" {
			lines = append(lines, cur)
		}
		for _, l := range lines {
			s.add(size*1.5, func(p *pen, y float64) {
				p.obj().Mark(bdf.MarkParagraph, "")
				base := y + size*1.1
				p.obj().Save()
				p.obj().ClipRect(f32(margin), f32(y), f32(contentW), f32(size*1.5))
				if sc.rtl {
					p.sample(l, margin+contentW, base, size, true, sc.lang)
				} else {
					p.sample(l, margin, base, size, false, sc.lang)
				}
				p.obj().Restore()
				p.set, p.fset = false, false
			})
		}
		s.space(6)
	}
}

// labelW is the width of the label column of the overview's tables.
const labelW = 190.0

// table adds a section of rows of a label and a value.
func (c *converter) table(s *scroll, sec section) {
	if len(sec.rows) == 0 {
		return
	}
	c.heading(s, sec.title)
	for _, r := range sec.rows {
		lines := c.ui.wrap(r.value, bodySize, false, contentW-labelW-8)
		labels := c.ui.wrap(r.label, smallSize+0.5, false, labelW-10)
		n := max(len(lines), len(labels))
		h := float64(n)*13 + 6
		s.add(h, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			for k, l := range labels {
				p.label(l, margin, y+12+float64(k)*13, smallSize+0.5, false, colGray, 0)
			}
			for k, l := range lines {
				p.obj().Mark(bdf.MarkLine, "")
				t := c.ui.line(l, bodySize, false)
				for j := range t.runs {
					if r.lang != "" && t.runs[j].lang == "" {
						t.runs[j].lang = r.lang
					}
				}
				p.text(t, margin+labelW, y+12+float64(k)*13, colText)
			}
			p.rect(margin, y+h-0.5, contentW, 0.5, colRule)
		})
	}
}

// scriptsTable adds the scripts and language systems of the layout
// tables, the features, and the characters by Unicode script.
func (c *converter) scriptsTable(s *scroll) {
	fc := c.fc
	sec := section{title: "Scripts and features"}
	for _, t := range []struct {
		name string
		t    *otlayout.Table
	}{{"GSUB", fc.gsub}, {"GPOS", fc.gpos}} {
		if t.t == nil {
			continue
		}
		var scripts []string
		for _, sc := range t.t.Scripts {
			n := sc.Tag
			if nm := otlayout.ScriptName(sc.Tag); nm != "" {
				n = nm + " (" + strings.TrimSpace(sc.Tag) + ")"
			}
			var langs []string
			for _, l := range sc.Langs {
				langs = append(langs, strings.TrimSpace(l.Tag))
			}
			if len(langs) > 0 {
				n += ": " + strings.Join(langs, " ")
			}
			scripts = append(scripts, n)
		}
		sec.add(t.name+" scripts", strings.Join(scripts, "; "))
		var tags []string
		seen := map[string]bool{}
		for _, f := range t.t.Features {
			if !seen[f.Tag] {
				seen[f.Tag] = true
				tags = append(tags, f.Tag)
			}
		}
		sec.add(t.name+" features", strings.Join(tags, " "))
	}
	var counts []string
	for _, sc := range fc.scriptCounts() {
		counts = append(counts, fmt.Sprintf("%s %s", strings.ReplaceAll(sc.name, "_", " "), num(sc.n)))
	}
	sec.add("Characters by script", strings.Join(counts, ", "))
	c.table(s, sec)
}

// paletteRows adds the colors of the CPAL palettes.
func (c *converter) paletteRows(s *scroll) {
	pals := c.fc.palettes()
	if len(pals) == 0 {
		return
	}
	c.heading(s, "Color palettes")
	const sw = 14.0
	perRow := int((contentW - 70) / (sw + 3))
	for i, pal := range pals {
		rows := (len(pal) + perRow - 1) / perRow
		h := float64(max(rows, 1))*(sw+3) + 6
		s.add(h, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			p.label("Palette "+strconv.Itoa(i), margin, y+11, smallSize, false, colGray, 0)
			for k, col := range pal {
				x := margin + 70 + float64(k%perRow)*(sw+3)
				yy := y + float64(k/perRow)*(sw+3)
				p.rect(x, yy, sw, sw, colVoid)
				p.rect(x+0.5, yy+0.5, sw-1, sw-1, bdf.Color(col))
			}
		})
	}
}

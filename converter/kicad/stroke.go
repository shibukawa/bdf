package kicad

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"math"
	"strings"
	"sync"
	"unicode"
)

// The stroke font KiCad draws text with by default: NewStroke, one line per
// code point from U+0020 (see tools/gen-newstroke and NEWSTROKE.txt).
//
//go:embed newstroke.txt.gz
var newstrokeGz []byte

var (
	strokeOnce  sync.Once
	strokeLines []string
	strokeMu    sync.Mutex
	strokeCache = map[rune]*glyph{}
	missingLine string // the glyph of U+007F, which stands for glyphs the font lacks
)

func loadStroke() {
	zr, err := gzip.NewReader(bytes.NewReader(newstrokeGz))
	if err != nil {
		return
	}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 64<<10), 64<<10)
	for sc.Scan() {
		strokeLines = append(strokeLines, sc.Text())
	}
	if len(strokeLines) > 0x7F-0x20 {
		missingLine = strokeLines[0x7F-0x20]
	}
}

// The metrics of KiCad's stroke font (STROKE_FONT, KiCad 10).
const (
	fontScale          = 1.0 / 21
	fontOffset         = -8 // moves the glyphs' origin to the baseline
	italicTilt         = 1.0 / 8
	interlinePitch     = 1.68 * 0.9583 // of the text height, with KiCad's legacy factor
	overbarHeight      = 1.23
	underlineOffset    = -0.16
	interChar          = 0.2
	tabWidth           = 4 // in character cells
	superSubSize       = 0.8
	subscriptOffset    = 0.15  // down, of the reduced height
	superscriptOffset  = -0.35 // down, of the reduced height
	firstLineHeight    = 1.17
	barTrim            = 0.1 // of the text width, off each end of over- and underlines
	strokeOffsetX      = 1 / 1.52
	strokeOffsetY      = 0.052
	boxFudge           = 0.17
	overbarBoxFraction = 1.0 / 6
)

type pt struct{ X, Y float64 }

// glyph is a decoded glyph: strokes in ems (x from the left bearing, y down
// with the baseline near 0), its advance and the lowest point.
type glyph struct {
	strokes [][]pt
	width   float64
	maxY    float64
}

// strokeGlyph returns the glyph of a character, or nil when the font lacks
// it.
func strokeGlyph(r rune) *glyph {
	strokeOnce.Do(loadStroke)
	i := int(r) - 0x20
	if i < 0 || i >= len(strokeLines) {
		return nil
	}
	strokeMu.Lock()
	defer strokeMu.Unlock()
	if g, ok := strokeCache[r]; ok {
		return g
	}
	line := strokeLines[i]
	var g *glyph
	if line != missingLine || r == 0x7F {
		g = decodeGlyph(line)
	}
	strokeCache[r] = g
	return g
}

// decodeGlyph decodes a glyph of the Hershey encoding: the left and right
// bearing, then points (each coordinate a character, 'R' being 0), " R"
// lifting the pen.
func decodeGlyph(s string) *glyph {
	if len(s) < 2 {
		return nil
	}
	startX := float64(int(s[0])-'R') * fontScale
	endX := float64(int(s[1])-'R') * fontScale
	g := &glyph{width: endX - startX}
	var cur []pt
	for i := 2; i+1 < len(s); i += 2 {
		if s[i] == ' ' && s[i+1] == 'R' {
			if len(cur) > 0 {
				g.strokes = append(g.strokes, cur)
			}
			cur = nil
			continue
		}
		p := pt{float64(int(s[i])-'R')*fontScale - startX, float64(int(s[i+1])-'R'+fontOffset) * fontScale}
		g.maxY = math.Max(g.maxY, p.Y)
		cur = append(cur, p)
	}
	if len(cur) > 0 {
		g.strokes = append(g.strokes, cur)
	}
	return g
}

// fallbackAdvance returns the advance in ems (of the text width) that KiCad
// gives a character its stroke font draws since KiCad 6 with CJK glyphs
// this font lacks, and whether it is one of those: the converter draws
// them with TrueType fonts in the same room.
func fallbackAdvance(r rune) float64 {
	switch {
	case r >= 0x3000 && r <= 0x30FF: // CJK punctuation, hiragana, katakana
		return 22.0 / 21
	case r >= 0x4E00 && r <= 0x9FFF: // CJK unified ideographs
		return 31.0 / 21
	case r >= 0xFF61 && r <= 0xFFDF: // halfwidth forms
		return 11.0 / 21
	case r >= 0xFF00 && r <= 0xFFEF: // fullwidth forms
		return 22.0 / 21
	case unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hangul, r):
		return 31.0 / 21
	}
	return 22.0 / 21
}

// markup is a node of KiCad's text markup: ^{superscript}, _{subscript} and
// ~{overbar}, nested.
type markup struct {
	text              string
	sub, sup, overbar bool
	children          []*markup
	root              bool
}

// parseMarkup parses marked-up text as KiCad's MARKUP parser does.
func parseMarkup(s string) *markup {
	type tok struct {
		text    string
		isText  bool
		open    bool
		close   bool
		control byte
	}
	var toks []tok
	start := 0
	var control byte
	var opens []int // the tokens of the groups open
	rs := []byte(s)
	for i := 0; i <= len(rs); i++ {
		if i == len(rs) {
			toks = append(toks, tok{text: string(rs[start:i]), isText: true})
			break
		}
		switch c := rs[i]; c {
		case '_', '^', '~':
			control = c
		case '{':
			if control != 0 {
				toks = append(toks, tok{text: string(rs[start : i-1]), isText: true}, tok{open: true, control: control})
				opens = append(opens, len(toks)-1)
				control = 0
				start = i + 1
			}
		case '}':
			if len(opens) > 0 {
				toks = append(toks, tok{text: string(rs[start:i]), isText: true}, tok{close: true})
				start = i + 1
				opens = opens[:len(opens)-1]
			}
			control = 0
		default:
			control = 0
		}
	}
	// a group never closed is text
	for _, i := range opens {
		toks[i] = tok{text: string([]byte{toks[i].control, '{'}), isText: true}
	}
	pos := 0
	var parse func() *markup
	parse = func() *markup {
		n := &markup{}
		for pos < len(toks) {
			t := toks[pos]
			pos++
			switch {
			case t.isText:
				if t.text != "" {
					n.children = append(n.children, &markup{text: t.text})
				}
			case t.open:
				c := parse()
				switch t.control {
				case '^':
					c.sup = true
				case '_':
					c.sub = true
				case '~':
					c.overbar = true
				}
				n.children = append(n.children, c)
			case t.close:
				return n
			}
		}
		return n
	}
	root := parse()
	root.root = true
	return root
}

// plain returns the text of marked-up text without the markup.
func (m *markup) plain() string {
	var sb strings.Builder
	var walk func(n *markup)
	walk = func(n *markup) {
		sb.WriteString(n.text)
		for _, c := range n.children {
			walk(c)
		}
	}
	walk(m)
	return sb.String()
}

// cjkPart is a character the stroke font lacks, drawn with a TrueType
// font in the middle of its cell: the cell starts at x along the line and
// is cell wide (mm), the font's em square is em tall, and the baseline is
// lifted by dy (mm, up) for superscripts.
type cjkPart struct {
	r           rune
	x, cell, dy float64
	em          float64
}

// cjkEm is the em size of the TrueType characters in heights of the text:
// KiCad's CJK strokes stand about as tall as capitals (kana) or a little
// taller (kanji), as CJK fonts' glyphs do in 1.25 em.
const cjkEm = 1.25

// lineGlyphs is one line of text laid out in its own space: x along the
// baseline from the line's start, y down, in mm.
type lineGlyphs struct {
	text    string // the line without markup
	strokes [][]pt
	parts   []cjkPart
	advance float64
	// top and bottom of the glyph cells (y down, mm), for the text box
	top, bottom float64
}

// strokeLayout lays out text as KiCad's stroke font does (FONT and
// STROKE_FONT), rounding as KiCad does to the file's internal unit (iu
// units per mm).
type strokeLayout struct {
	sizeX, sizeY float64 // glyph width and height, mm
	italic       bool
	underline    bool
	pen          float64 // the stroke width, mm
	iu           float64
}

// round and trunc round to the internal unit as KiROUND and an int
// conversion do.
func (l *strokeLayout) round(v float64) float64 { return math.Round(v*l.iu) / l.iu }
func (l *strokeLayout) trunc(v float64) float64 {
	return math.Trunc(v*l.iu+1e-9*math.Copysign(1, v)) / l.iu
}

// spaceAdvance is the advance of a space in text widths.
func spaceAdvance() float64 {
	if g := strokeGlyph(' '); g != nil {
		return g.width
	}
	return 0.6
}

// textStyle is the style of a run of text.
type textStyle struct {
	italic, sub, sup, underline bool
}

// line lays out one line of marked-up text.
func (l *strokeLayout) line(s string) *lineGlyphs {
	m := parseMarkup(s)
	out := &lineGlyphs{text: m.plain(), top: -l.sizeY}
	out.advance = l.node(out, m, 0, textStyle{italic: l.italic, underline: l.underline})
	return out
}

// node lays out a node of markup (drawMarkup) and returns the cursor after
// it.
func (l *strokeLayout) node(out *lineGlyphs, n *markup, x float64, st textStyle) float64 {
	next := x
	ns := st
	drawOver, drawUnder := false, false
	if !n.root {
		if n.sub {
			ns.sub = true
		} else if n.sup {
			ns.sup = true
		}
		drawOver = n.overbar
		drawUnder = st.underline
		if n.text != "" {
			next = l.run(out, n.text, next, ns)
		}
	}
	for _, c := range n.children {
		next = l.node(out, c, next, ns)
	}
	trim := l.sizeX * barTrim
	if drawUnder && next-x > 2*trim {
		y := -l.sizeY * underlineOffset
		out.strokes = append(out.strokes, []pt{{x + trim, y}, {next - trim, y}})
	}
	if drawOver && next-x > 2*trim {
		y := -l.sizeY * overbarHeight
		out.strokes = append(out.strokes, []pt{{x + trim, y}, {next - trim, y}})
	}
	return next
}

// run lays out a run of text in one style starting at x0 and returns the
// cursor after it (STROKE_FONT::GetTextAsGlyphs).
func (l *strokeLayout) run(out *lineGlyphs, s string, x0 float64, st textStyle) float64 {
	sx, sy := l.sizeX, l.sizeY
	dy := 0.0
	if st.sub || st.sup {
		sx *= superSubSize
		sy *= superSubSize
		if st.sub {
			dy = sy * subscriptOffset
		} else {
			dy = sy * superscriptOffset
		}
	}
	out.top = math.Min(out.top, dy-sy)
	out.bottom = math.Max(out.bottom, dy)
	tilt := 0.0
	if st.italic {
		tilt = italicTilt
	}
	space := spaceAdvance()
	x := x0
	count := 0
	for _, r := range s {
		switch r {
		case '\t':
			// to the next column of four character cells
			count = (count/tabWidth+1)*tabWidth - 1
			nx := l.trunc(x0 + l.sizeX*float64(count) + l.sizeX*space)
			for nx <= x {
				count += tabWidth
				nx += l.sizeX * tabWidth
			}
			x = nx
			count++
			continue
		case ' ':
			x += l.round(sx * space)
			count++
			continue
		}
		g := strokeGlyph(r)
		if g == nil {
			adv := l.round(sx * fallbackAdvance(r))
			if unicode.IsGraphic(r) {
				out.parts = append(out.parts, cjkPart{r: r, x: x, cell: adv, dy: -dy, em: sy * cjkEm})
			}
			x += adv
			count++
			continue
		}
		for _, stroke := range g.strokes {
			pts := make([]pt, len(stroke))
			for i, p := range stroke {
				px, py := p.X*sx, p.Y*sy
				if tilt != 0 {
					px -= py * tilt
				}
				pts[i] = pt{px + x, py + dy}
			}
			out.strokes = append(out.strokes, pts)
		}
		x += l.round(g.width * sx)
		count++
	}
	return x
}

// textBlock is a text of one or more lines placed as KiCad places it.
type textBlock struct {
	lines []*lineGlyphs
	// pos is the position of each line's start, before the text is turned
	// or mirrored about its anchor (mm, y down)
	pos []pt
}

// hAlign and vAlign are the justification of a text.
type hAlign int8
type vAlign int8

const (
	alignLeft   hAlign = -1
	alignCenter hAlign = 0
	alignRight  hAlign = 1
	alignTop    vAlign = -1
	alignMiddle vAlign = 0
	alignBottom vAlign = 1
)

// interline returns the distance between lines.
func (l *strokeLayout) interline(lineSpacing float64) float64 {
	return l.trunc(l.sizeY * interlinePitch * lineSpacing)
}

// splitLines splits a text into its lines as KiCad does (wxStringSplit):
// a newline at the end starts no line.
func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	if len(lines) > 1 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// layoutBlock lays out the lines of a text anchored at (0, 0)
// (FONT::getLinePositions).
func (l *strokeLayout) layoutBlock(text string, h hAlign, v vAlign, lineSpacing float64) *textBlock {
	if lineSpacing == 0 {
		lineSpacing = 1
	}
	lines := splitLines(text)
	b := &textBlock{}
	interline := l.interline(lineSpacing)
	height := 0.0
	for i, s := range lines {
		b.lines = append(b.lines, l.line(s))
		if i == 0 {
			height = l.trunc(height + l.sizeY*firstLineHeight)
		} else {
			height += interline
		}
	}
	offX := l.trunc(l.pen * strokeOffsetX)
	offY := l.sizeY - l.trunc(l.pen*strokeOffsetY)
	switch v {
	case alignMiddle:
		offY -= l.trunc(height / 2)
	case alignBottom:
		offY -= height
	}
	for i, lg := range b.lines {
		x := offX
		switch h {
		case alignCenter:
			x = -l.trunc(lg.advance / 2)
		case alignRight:
			x = -(lg.advance + offX)
		}
		b.pos = append(b.pos, pt{x, offY + float64(i)*interline})
	}
	return b
}

// boundary returns the width and height of a line as KiCad measures it
// (FONT::StringBoundaryLimits): the glyph cells without the space after
// the last character, grown by the pen.
func (l *strokeLayout) boundary(lg *lineGlyphs) (w, h float64) {
	w = lg.advance - l.round(l.sizeX*interChar)
	h = lg.bottom - lg.top
	grow := 2 * l.round(l.pen*1.5)
	return math.Max(w, 0) + grow, h + grow
}

// textBox returns the box of a text at (0, 0) as KiCad computes it
// (EDA_TEXT::GetTextBox), before the text is turned: x0, y0 is its top
// left corner (y down).
func (l *strokeLayout) textBox(text string, h hAlign, v vAlign, lineSpacing float64, mirror bool) (x0, y0, w, hh float64) {
	if lineSpacing == 0 {
		lineSpacing = 1
	}
	lines := splitLines(text)
	first := l.line(lines[0])
	w, ext := l.boundary(first)
	fudge := l.round(ext * boxFudge)
	hh = ext + fudge
	for _, s := range lines[1:] {
		lw, _ := l.boundary(l.line(s))
		w = math.Max(w, lw)
	}
	if n := len(lines); n > 1 {
		hh += l.round(float64(n-1) * l.sizeY * interlinePitch * lineSpacing)
	}
	if strings.Contains(lines[0], "~{") {
		hh += l.trunc(ext * overbarBoxFraction)
	}
	italic := 0.0
	if l.italic {
		italic = l.round(l.sizeY * italicTilt)
	}
	switch h {
	case alignLeft:
		if mirror {
			x0 -= w - italic
		}
	case alignCenter:
		x0 -= l.trunc((w - italic) / 2)
	case alignRight:
		if !mirror {
			x0 -= w - italic
		}
	}
	switch v {
	case alignTop:
		y0 -= fudge
	case alignMiddle:
		y0 -= l.trunc(hh / 2)
	case alignBottom:
		y0 -= hh
		y0 += fudge
	}
	return x0, y0, w, hh
}

// Pen widths of text whose thickness is not given (EDA_TEXT), from the
// smaller of its width and height.
func boldThickness(size float64) float64   { return size / 5 }
func normalThickness(size float64) float64 { return size / 8 }

// clampThickness keeps text from being so thick its strokes merge
// (ClampTextPenSize): to a quarter of its size, or less when strict (the
// texts of pins).
func clampThickness(t, size float64, strict bool) float64 {
	f := 0.25
	if strict {
		f = 0.18
	}
	return math.Min(t, math.Round(size*f*1e6)/1e6)
}

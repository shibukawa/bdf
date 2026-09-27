package hpgl

import (
	"math"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

// Labels are drawn as text in the fonts of the document. A plotter draws
// them in its stick font, a fixed-pitch font of strokes, or in the font an
// SD or AD instruction designates; the text keeps their sizes and places:
// the cap height of the characters, their pitch (each run of characters is
// stretched to its cells), the label direction, slant and text path, and
// the label origin.

// fontDef is a font designated by SD or AD.
type fontDef struct {
	symbolSet int
	spacing   int     // 0 fixed, 1 proportional
	pitch     float64 // characters per inch
	height    float64 // points
	posture   int
	weight    int
	typeface  int
}

// the default font: the stick font in Roman-8, 9 characters per inch,
// 11.5 points
var defaultFont = fontDef{symbolSet: 277, pitch: 9, height: 11.5, typeface: 48}

// stick reports whether a font is one of the plotter's stroke fonts (stick,
// fixed arc and variable arc).
func (f fontDef) stick() bool { return f.typeface >= 48 && f.typeface <= 50 }

type labelState struct {
	std, alt  fontDef
	useAlt    bool
	sizeMode  int // 0 the font's size, 1 SI (cm), 2 SR (percentages of P1–P2)
	sizeW     float64
	sizeH     float64
	relDir    bool // DR (DI otherwise)
	run, rise float64
	path      int // DV: 0 right, 1 down, 2 left, 3 up
	reverseLF bool
	origin    int // LO
	extraW    float64
	extraH    float64
	slant     float64
	term      byte
	termPrint bool
	td        bool
	mode16    bool
	row       int
}

func (l *labelState) defaults() {
	*l = labelState{std: defaultFont, alt: defaultFont, run: 1, origin: 1, term: 0x03}
}

func (l *labelState) font() fontDef {
	if l.useAlt {
		return l.alt
	}
	return l.std
}

// labelText reads the text of a label up to its terminator. In 16-bit
// mode the terminator is a byte pair with a zero first byte.
func (c *converter) labelText() []byte {
	b := c.b
	l := &c.gl.label
	start := c.i
	if l.mode16 {
		for i := c.i; i+1 < len(b); i += 2 {
			if b[i] == 0 && b[i+1] == l.term {
				text := b[start:i]
				if l.termPrint {
					text = b[start : i+2]
				}
				c.i = i + 2
				return text
			}
		}
	} else if k := indexByte(b[c.i:], l.term); k >= 0 {
		end := c.i + k
		c.i = end + 1
		if l.termPrint {
			return b[start:c.i]
		}
		return b[start:end]
	}
	// no terminator: a plotter would label the rest of the file
	end := min(len(b), c.i+4096)
	c.i = end
	c.warnOnce("unterminated", "a label has no terminator; its text ends with the file")
	return b[start:end]
}

func indexByte(b []byte, c byte) int {
	for i, v := range b {
		if v == c {
			return i
		}
	}
	return -1
}

func (g *gl) defineTerminator() {
	c := g.c
	b := c.b
	l := &g.label
	if c.i >= len(b) {
		return
	}
	t := b[c.i]
	if t == ';' {
		c.i++
		l.term, l.termPrint = 0x03, false
		return
	}
	c.i++
	a := c.params()
	if t == 0 || t == 5 || t == 0x1b {
		return
	}
	l.term = t
	l.termPrint = len(a) > 0 && a[0] == 0
}

func (g *gl) symbolMode() {
	c := g.c
	b := c.b
	g.symbol = nil
	if c.i >= len(b) {
		return
	}
	n := 1
	if g.label.mode16 {
		n = 2
	}
	s := b[c.i:min(len(b), c.i+n)]
	if s[0] == ';' || s[0] < 0x20 && !g.label.mode16 {
		if s[0] == ';' {
			c.i++
		}
		return
	}
	g.symbol = append([]byte(nil), s...)
	c.i += len(s)
}

// characterInstruction carries out the instructions of the character
// group that take numeric parameters.
func (g *gl) characterInstruction(m string, a []float64) {
	l := &g.label
	if !valid(a...) {
		return
	}
	switch m {
	case "SI", "SR":
		switch {
		case len(a) == 0 && m == "SI":
			l.sizeMode = 0
		case len(a) == 0:
			l.sizeMode, l.sizeW, l.sizeH = 2, 0.75, 1.5
		case len(a) >= 2 && a[0] != 0 && a[1] != 0:
			l.sizeW, l.sizeH = a[0], a[1]
			l.sizeMode = 1
			if m == "SR" {
				l.sizeMode = 2
			}
		}
	case "SD", "AD":
		f := &l.std
		if m == "AD" {
			f = &l.alt
		}
		if len(a) == 0 {
			*f = defaultFont
			return
		}
		for i := 0; i+1 < len(a); i += 2 {
			v := a[i+1]
			switch int(a[i]) {
			case 1:
				f.symbolSet = int(v)
			case 2:
				f.spacing = int(v)
			case 3:
				if v > 0 {
					f.pitch = v
				}
			case 4:
				if v > 0 {
					f.height = v
				}
			case 5:
				f.posture = int(v)
			case 6:
				f.weight = int(v)
			case 7:
				f.typeface = int(v)
			}
		}
	case "SS":
		l.useAlt = false
	case "SA":
		l.useAlt = true
	case "CS", "CA":
		// the character sets of HP-GL: set 0 is ASCII
		f := &l.std
		if m == "CA" {
			f = &l.alt
		}
		f.symbolSet = 277
	case "FI", "FN":
		g.c.warnOnce("FI", "fonts selected by ID (FI, FN) are drawn in the default font")
	case "DI", "DR":
		l.relDir = m == "DR"
		l.run, l.rise = 1, 0
		if len(a) >= 2 && (a[0] != 0 || a[1] != 0) {
			l.run, l.rise = a[0], a[1]
		}
		g.crPoint = g.pos
	case "DV":
		l.path, l.reverseLF = 0, false
		if len(a) > 0 && a[0] >= 0 && a[0] <= 3 {
			l.path = int(a[0])
		}
		if len(a) > 1 {
			l.reverseLF = a[1] == 1
		}
		g.crPoint = g.pos
	case "LO":
		n := 1
		if len(a) > 0 {
			n = int(a[0])
		}
		if n >= 1 && n <= 9 || n >= 11 && n <= 19 || n == 21 {
			l.origin = n
		}
		g.crPoint = g.pos
	case "ES":
		l.extraW, l.extraH = 0, 0
		if len(a) > 0 {
			l.extraW = a[0]
		}
		if len(a) > 1 {
			l.extraH = a[1]
		}
	case "SL":
		l.slant = 0
		if len(a) > 0 {
			l.slant = math.Max(-100, math.Min(100, a[0]))
		}
	case "CP":
		g.characterPlot(a)
	case "TD":
		l.td = len(a) > 0 && a[0] == 1
	case "LM":
		l.mode16 = len(a) > 0 && a[0] == 1
		l.row = 0
		if len(a) > 1 && a[1] >= 0 && a[1] <= 255 {
			l.row = int(a[1])
		}
		g.symbol = nil
	case "CF", "SB":
	}
}

// metrics are the sizes of the characters of the font in effect, in
// plotter units: the cap height, the cell of a fixed-pitch font (0 for a
// proportional one), the horizontal size of an em for proportional
// advances, the line feed, the advance of the vertical text paths and the
// offset of label origins 11 to 19. A negative size (SI, SR) mirrors the
// characters.
type metrics struct {
	capH, cell, em, lf, vadv, offset float64
	fixed, stick                     bool
	mirrorX, mirrorY                 bool
	spec                             cad.FontSpec
}

func (g *gl) metrics() metrics { return g.metricsSized(g.label.sizeMode) }

// metricsSized returns the metrics with a size mode: 0 for the font's own
// size, which SO and SI select for the rest of a label.
func (g *gl) metricsSized(sizeMode int) metrics {
	l := &g.label
	f := l.font()
	m := metrics{fixed: f.spacing == 0, stick: f.stick()}
	const pt = 1016.0 / 72
	var w, h float64 // nominal character width and cap height
	switch sizeMode {
	case 1:
		w, h = l.sizeW*400, l.sizeH*400
	case 2:
		w, h = l.sizeW/100*(g.p2.X-g.p1.X), l.sizeH/100*(g.p2.Y-g.p1.Y)
	}
	m.mirrorX, m.mirrorY = w < 0, h < 0
	w, h = math.Abs(w), math.Abs(h)
	if m.stick {
		if sizeMode == 0 {
			m.cell = 1016 / f.pitch
			w = m.cell / 1.5
			h = 0.67 * f.height * pt
		} else {
			m.cell = 1.5 * w
		}
		// the stick font's point size is 1.5 times its cap height
		size := h / 0.67
		m.capH = h
		m.lf = 2 * h
		m.em = 2 * w
		m.offset = size / 3
		m.vadv = 0.96 * size
	} else {
		size := f.height * pt
		if m.fixed {
			size = 120 / f.pitch * pt
		}
		if sizeMode == 0 {
			m.capH = 0.7 * size
			m.em = size
			m.cell = 0.6 * size
		} else {
			m.capH = h
			m.em = 2 * w
			m.cell = 1.2 * w
			size = h / 0.7
		}
		m.lf = 1.2 * size
		m.offset = size / 4
		m.vadv = 0.9 * size
	}
	if !m.fixed {
		m.cell = 0
	}
	m.lf *= 1 + l.extraH
	m.spec = labelFont(f)
	return m
}

// labelFont returns the font that draws a designated font: its typeface
// and spacing choose a sans-serif, serif or fixed-pitch font, its stroke
// weight and posture bold and italic.
func labelFont(f fontDef) cad.FontSpec {
	spec := cad.FontSpec{Bold: f.weight >= 3 && f.weight != 9999, Italic: f.posture == 1 || f.posture == 2}
	switch f.typeface {
	case 0, 3, 6, 4099, 4102, 48, 49:
		// line printer, Courier, Letter Gothic and the fixed-pitch stick fonts
		spec.Family = "monospace"
	case 5, 8, 18, 22, 23, 24, 4101, 4140, 4197, 4297, 16901:
		// Times, Prestige, Garamond, Bodoni, Century Schoolbook, University
		// Roman, CG Times, Clarendon, Marigold, Times New Roman
		spec.Family = "serif"
	}
	if f.spacing == 0 && spec.Family == "" {
		spec.Family = "monospace"
	}
	return spec
}

// labelSpace returns the transform of label units (plotter units along the
// label direction, and up from it) into the drawing, mirrored as the
// character size says.
func (g *gl) labelSpace(m metrics) canvas.Matrix {
	l := &g.label
	run, rise := l.run, l.rise
	if l.relDir {
		run = run / 100 * (g.p2.X - g.p1.X)
		rise = rise / 100 * (g.p2.Y - g.p1.Y)
	}
	n := math.Hypot(run, rise)
	if n == 0 || !valid(n) {
		run, rise, n = 1, 0, 1
	}
	c, s := run/n, rise/n
	sx, sy := 1.0, 1.0
	if m.mirrorX {
		sx = -1
	}
	if m.mirrorY {
		sy = -1
	}
	d := g.dev
	return canvas.Matrix{d[0], d[1], d[2], d[3], 0, 0}.Mul(canvas.Matrix{c, s, -s, c, 0, 0}).Mul(canvas.Scale(sx, sy))
}

// textPath returns the direction of the text path and of line feeds in
// label units.
func (l *labelState) textPath() (path, lf cad.Point) {
	path = [4]cad.Point{{X: 1}, {Y: -1}, {X: -1}, {Y: 1}}[l.path]
	lf = cad.Point{X: path.Y, Y: -path.X} // −90° from the path
	if l.reverseLF {
		lf = lf.Mul(-1)
	}
	return path, lf
}

// glyph is a character of a label placed on the page.
type glyph struct {
	r   rune
	at  cad.Point     // origin (drawing)
	adv float64       // advance along the text path (plotter units)
	m   metrics       // the size and font
	lm  canvas.Matrix // label units to the drawing
}

// labelInstruction draws a label (LB) from the pen location.
func (g *gl) labelInstruction(text []byte) {
	if g.inPoly {
		return
	}
	g.flush()
	g.drawLabel(g.decode(text))
}

// labelChar is a decoded character of a label: a printing character or a
// control code.
type labelChar struct {
	r       rune
	control byte // the control code, when ctrl
	ctrl    bool
	wide    bool
}

// decode turns the bytes of a label into characters, by the symbol set
// of the font in effect (switched by SO and SI inside the label).
func (g *gl) decode(text []byte) []labelChar {
	l := &g.label
	var out []labelChar
	useAlt := l.useAlt
	set := func() int {
		if useAlt {
			return l.alt.symbolSet
		}
		return l.std.symbolSet
	}
	if l.mode16 {
		for i := 0; i+1 < len(text); i += 2 {
			b1, b2 := text[i], text[i+1]
			if b1 == 0 {
				if b2 < 0x20 && !l.td {
					out = append(out, labelChar{control: b2, ctrl: true})
					if b2 == 0x0e {
						useAlt = true
					} else if b2 == 0x0f {
						useAlt = false
					}
					continue
				}
				out = append(out, labelChar{r: decodeByte(b2, set())})
				continue
			}
			out = append(out, labelChar{r: decodeJIS(b1, b2), wide: true})
		}
		return out
	}
	sjis := japaneseSet(set()) || looksShiftJIS(text)
	for i := 0; i < len(text); i++ {
		b := text[i]
		if (b < 0x20 || b == 0x7f) && !l.td {
			out = append(out, labelChar{control: b, ctrl: true})
			if b == 0x0e {
				useAlt = true
			} else if b == 0x0f {
				useAlt = false
			}
			continue
		}
		if l.row != 0 && (set()%32 == 11 || set() >= 1600) {
			// a row of a 16-bit set in 8-bit mode
			out = append(out, labelChar{r: decodeJIS(byte(l.row), b), wide: true})
			continue
		}
		if sjis && isSJISLead(b) && i+1 < len(text) {
			if r, ok := decodeSJIS(text[i], text[i+1]); ok {
				out = append(out, labelChar{r: r, wide: true})
				i++
				continue
			}
		}
		if sjis && b >= 0xa1 && b <= 0xdf {
			out = append(out, labelChar{r: rune(0xff61 + int(b) - 0xa1)})
			continue
		}
		out = append(out, labelChar{r: decodeByte(b, set())})
	}
	return out
}

// japaneseSet reports whether a symbol set is Japanese (JIS, Kana-8, the
// Kanji sets).
func japaneseSet(s int) bool { return s%32 == 11 || s == 1611 || s == 1643 }

func isSJISLead(b byte) bool { return b >= 0x81 && b <= 0x9f || b >= 0xe0 && b <= 0xfc }

var sjisDecoder = japanese.ShiftJIS.NewDecoder()

func decodeSJIS(b1, b2 byte) (rune, bool) {
	if b2 < 0x40 || b2 == 0x7f || b2 > 0xfc {
		return 0, false
	}
	s, err := sjisDecoder.Bytes([]byte{b1, b2})
	if err != nil {
		return 0, false
	}
	r, _ := utf8.DecodeRune(s)
	if r == utf8.RuneError {
		return 0, false
	}
	return r, true
}

// looksShiftJIS reports whether the high bytes of a label are Shift_JIS
// text: every one of them in a double-byte character, at least one of
// which is kana or a kanji of JIS level 1 (accented Roman-8 letters
// followed by ASCII letters make valid pairs of level 2 kanji and
// half-width katakana, but not those).
func looksShiftJIS(text []byte) bool {
	common := 0
	for i := 0; i < len(text); i++ {
		b := text[i]
		if b < 0x80 {
			continue
		}
		if !isSJISLead(b) || i+1 >= len(text) {
			return false
		}
		if _, ok := decodeSJIS(b, text[i+1]); !ok {
			return false
		}
		if b >= 0x81 && b <= 0x9f {
			common++
		}
		i++
	}
	return common > 0
}

var eucDecoder = japanese.EUCJP.NewDecoder()

// decodeJIS decodes a character of JIS X 0208 (row and cell).
func decodeJIS(b1, b2 byte) rune {
	if b1 >= 0x21 && b1 <= 0x7e && b2 >= 0x21 && b2 <= 0x7e {
		s, err := eucDecoder.Bytes([]byte{b1 | 0x80, b2 | 0x80})
		if err == nil {
			if r, _ := utf8.DecodeRune(s); r != utf8.RuneError {
				return r
			}
		}
	}
	return ' '
}

// decodeByte decodes a character by a symbol set.
func decodeByte(b byte, set int) rune {
	if b >= 0x20 && b < 0x7f {
		return rune(b)
	}
	switch set {
	case 14: // 0N ISO 8859-1
		return rune(b)
	case 341: // 10U PC-8
		return charmap.CodePage437.DecodeByte(b)
	case 405: // 12U PC-850
		return charmap.CodePage850.DecodeByte(b)
	case 629, 309: // 19U Windows Latin 1, 9U
		return charmap.Windows1252.DecodeByte(b)
	}
	if japaneseSet(set) && b >= 0xa1 && b <= 0xdf {
		return rune(0xff61 + int(b) - 0xa1)
	}
	if b >= 0xa0 {
		return roman8[b-0xa0]
	}
	return ' '
}

// roman8 is the upper half of HP Roman-8.
var roman8 = [96]rune{
	' ', 'À', 'Â', 'È', 'Ê', 'Ë', 'Î', 'Ï', '´', 'ˋ', 'ˆ', '¨', '˜', 'Ù', 'Û', '₤',
	'¯', 'Ý', 'ý', '°', 'Ç', 'ç', 'Ñ', 'ñ', '¡', '¿', '¤', '£', '¥', '§', 'ƒ', '¢',
	'â', 'ê', 'ô', 'û', 'á', 'é', 'ó', 'ú', 'à', 'è', 'ò', 'ù', 'ä', 'ë', 'ö', 'ü',
	'Å', 'î', 'Ø', 'Æ', 'å', 'í', 'ø', 'æ', 'Ä', 'ì', 'Ö', 'Ü', 'É', 'ï', 'ß', 'Ô',
	'Á', 'Ã', 'ã', 'Ð', 'ð', 'Í', 'Ì', 'Ó', 'Ò', 'Õ', 'õ', 'Š', 'š', 'Ú', 'Ÿ', 'ÿ',
	'Þ', 'þ', '·', 'µ', '¶', '¾', '—', '¼', '½', 'ª', 'º', '«', '■', '»', '±', ' ',
}

// labelLine is a line of a label being laid out.
type labelLine struct {
	start  cad.Point // the reference point of the line (drawing)
	glyphs []glyph
}

// drawLabel lays out and draws the characters of a label from the pen
// location, and moves the pen after them when the label origin and text
// path let labels continue one another.
func (g *gl) drawLabel(chars []labelChar) {
	l := &g.label
	pathV, lfV := l.textPath()
	m := g.metrics()
	lm := g.labelSpace(m)
	start := g.pos
	pen := g.pos // the pen, drawing
	var lines []labelLine
	line := labelLine{start: pen}
	var last glyph
	haveLast := false
	flushLine := func() {
		if len(line.glyphs) > 0 {
			lines = append(lines, line)
		}
		line = labelLine{start: pen}
	}
	for _, ch := range chars {
		if ch.ctrl {
			switch ch.control {
			case 0x0d: // carriage return
				flushLine()
				pen = g.crPoint
				line.start = pen
			case 0x0a: // line feed
				flushLine()
				d := linear(lm, lfV.Mul(m.lf))
				pen = pen.Add(d)
				g.crPoint = g.crPoint.Add(d)
				line.start = pen
			case 0x08: // backspace
				if haveLast {
					pen = pen.Sub(linear(lm, pathV.Mul(last.adv)))
					haveLast = false
				}
			case 0x09: // tab: every eighth column from the carriage return point
				if col := g.spaceAdvance(m); col > 0 {
					rel := linear(invert(lm), pen.Sub(g.crPoint))
					along := rel.X*pathV.X + rel.Y*pathV.Y
					next := (math.Floor(along/(8*col)+1e-9) + 1) * 8 * col
					pen = pen.Add(linear(lm, pathV.Mul(next-along)))
				}
			case 0x0e, 0x0f: // shift out and in: the alternate and standard fonts
				// in the size of the font for the rest of the label, as
				// GhostPCL draws them (SI and SR hold for the next label)
				l.useAlt = ch.control == 0x0e
				m = g.metricsSized(0)
				lm = g.labelSpace(m)
			}
			continue
		}
		adv := m.cell
		switch {
		case l.path%2 == 1:
			adv = m.vadv
		case adv != 0 && ch.wide && !l.mode16:
			adv *= 2 // two bytes, two cells
		case adv == 0:
			adv = g.c.fonts.RuneAdvance(m.spec, ch.r) * m.em
		}
		adv *= 1 + l.extraW
		gl := glyph{r: ch.r, at: pen, adv: adv, m: m, lm: lm}
		line.glyphs = append(line.glyphs, gl)
		pen = pen.Add(linear(lm, pathV.Mul(adv)))
		last, haveLast = gl, true
	}
	flushLine()
	for i, ln := range lines {
		g.drawLine(ln, i > 0)
	}
	if concatenates(l.path, l.origin) {
		g.pos = pen
	} else {
		g.pos = start
	}
}

// spaceAdvance returns the width of a column (a space) in plotter units.
func (g *gl) spaceAdvance(m metrics) float64 {
	if m.cell != 0 {
		return m.cell * (1 + g.label.extraW)
	}
	return g.c.fonts.RuneAdvance(m.spec, ' ') * m.em * (1 + g.label.extraW)
}

// concatenates reports whether labels with a text path and a label origin
// continue one another.
func concatenates(path, origin int) bool {
	o := origin
	if o > 10 && o < 20 {
		o -= 10
	}
	switch path {
	case 0:
		return o == 1 || o == 2 || o == 3 || origin == 21
	case 1:
		return o == 3 || o == 6 || o == 9
	case 2:
		return o == 7 || o == 8 || o == 9
	default:
		return o == 1 || o == 4 || o == 7 || origin == 21
	}
}

// drawLine draws a line of a label, placed by the label origin: the box
// around the line (its advances, and the cap height) goes to the
// reference point by its corner, the middle of a side or its centre.
func (g *gl) drawLine(ln labelLine, cont bool) {
	l := &g.label
	g0 := ln.glyphs[0]
	m0, inv := g0.m, invert(g0.lm)
	var box cad.Rect
	for _, gl := range ln.glyphs {
		w := gl.adv
		if l.path%2 == 1 {
			// a column as wide as a character
			w = m0.em / 2
			if m0.cell != 0 {
				w = m0.cell
			}
		}
		p := linear(inv, gl.at.Sub(ln.start))
		box = box.Add(p).Add(p.Add(cad.Point{X: w, Y: gl.m.capH}))
	}
	o := l.origin
	if o == 21 {
		o = 1
	}
	base := o
	if base > 10 {
		base -= 10
	}
	var shift cad.Point
	switch (base - 1) / 3 {
	case 1:
		shift.X = -(box.Min.X + box.Max.X) / 2
	case 2:
		shift.X = -box.Max.X
	default:
		shift.X = -box.Min.X
	}
	switch (base - 1) % 3 {
	case 1:
		shift.Y = -(box.Min.Y + box.Max.Y) / 2
	case 2:
		shift.Y = -box.Max.Y
	default:
		shift.Y = -box.Min.Y
	}
	if o > 10 {
		// away from the reference point
		switch (base - 1) / 3 {
		case 0:
			shift.X += m0.offset
		case 2:
			shift.X -= m0.offset
		}
		switch (base - 1) % 3 {
		case 0:
			shift.Y += m0.offset
		case 2:
			shift.Y -= m0.offset
		}
	}
	col := g.penColorOf(g.pens.current())
	brk := cad.BreakBox
	if cont {
		brk = cad.BreakLine
	}
	switch l.path {
	case 0:
		g.drawRuns(ln.glyphs, shift, col, brk)
	case 1:
		g.drawColumn(ln.glyphs, shift, col, brk)
	default:
		for _, gl := range ln.glyphs {
			g.drawRuns([]glyph{gl}, shift, col, brk)
			brk = cad.BreakNone
		}
	}
}

// drawRuns draws glyphs along the baseline, a text for each run of one
// font and size; each run is stretched to the advances of its characters.
func (g *gl) drawRuns(glyphs []glyph, shift cad.Point, col bdf.Color, brk cad.Break) {
	for len(glyphs) > 0 {
		n := 1
		for n < len(glyphs) && glyphs[n].m == glyphs[0].m && glyphs[n].lm == glyphs[0].lm {
			n++
		}
		run := glyphs[:n]
		glyphs = glyphs[n:]
		var sb strings.Builder
		for _, gl := range run {
			sb.WriteRune(gl.r)
		}
		s := sb.String()
		if strings.TrimSpace(s) == "" {
			continue
		}
		m, lm := run[0].m, run[0].lm
		em := g.renderEm(m)
		cells := make([]float64, len(run))
		for i, gl := range run {
			cells[i] = gl.adv / em
		}
		q := run[0].at.Add(linear(lm, shift))
		tm := canvas.Matrix{lm[0], lm[1], lm[2], lm[3], q.X, q.Y}.Mul(canvas.Matrix{1, 0, g.label.slant, 1, 0, 0}).Mul(canvas.Scale(em, em))
		t := g.c.fonts.NewCellText(m.spec, s, cells, 0, col, tm)
		t.Break = brk
		brk = cad.BreakNone
		g.emitText(t)
	}
}

// drawColumn draws glyphs down a column (text path 1), upright.
func (g *gl) drawColumn(glyphs []glyph, shift cad.Point, col bdf.Color, brk cad.Break) {
	var sb strings.Builder
	cells := make([]float64, len(glyphs))
	m, lm := glyphs[0].m, glyphs[0].lm
	em := g.renderEm(m)
	for i, gl := range glyphs {
		sb.WriteRune(gl.r)
		cells[i] = gl.adv / em
	}
	s := sb.String()
	if strings.TrimSpace(s) == "" {
		return
	}
	w := m.cell / 1.5
	if m.cell == 0 {
		w = m.em / 2
	}
	// the column runs down from above the first character, through the
	// middle of the characters
	q := glyphs[0].at.Add(linear(lm, shift.Add(cad.Point{X: w / 2, Y: (glyphs[0].adv + m.capH) / 2})))
	tm := canvas.Matrix{lm[0], lm[1], lm[2], lm[3], q.X, q.Y}.Mul(canvas.Rotate(-90)).Mul(canvas.Scale(em, em))
	t := g.c.fonts.NewCellText(m.spec, s, cells, 0, col, tm)
	t.Vertical = true
	t.Break = brk
	g.emitText(t)
}

// renderEm returns the size of the em that gives the font drawing a label
// its cap height (plotter units).
func (g *gl) renderEm(m metrics) float64 {
	_, _, capH := g.c.fonts.Metrics(m.spec)
	if capH <= 0 {
		capH = 0.7
	}
	return m.capH / capH
}

func (g *gl) emitText(t *cad.Text) {
	if !g.c.budget() {
		return
	}
	g.c.hasText = true
	g.c.cur.d.Text(t)
	g.c.cur.marked = true
}

// characterPlot is CP: the pen moves by spaces and lines; without them, a
// carriage return and a line feed.
func (g *gl) characterPlot(a []float64) {
	l := &g.label
	m := g.metrics()
	lm := g.labelSpace(m)
	pathV, lfV := l.textPath()
	if len(a) < 2 {
		d := linear(lm, lfV.Mul(m.lf))
		g.crPoint = g.crPoint.Add(d)
		g.pos = g.crPoint
		return
	}
	spaces, lines := a[0], a[1]
	adv := g.spaceAdvance(m)
	if l.path%2 == 1 {
		adv = m.vadv
	}
	g.pos = g.pos.Add(linear(lm, pathV.Mul(spaces*adv)))
	d := linear(lm, lfV.Mul(-lines*m.lf))
	g.pos = g.pos.Add(d)
	g.crPoint = g.crPoint.Add(d)
}

// drawSymbol draws the symbol of symbol mode centred on the pen.
func (g *gl) drawSymbol() {
	if len(g.symbol) == 0 || g.inPoly {
		return
	}
	g.flush()
	chars := g.decode(g.symbol)
	if len(chars) == 0 || chars[0].ctrl {
		return
	}
	m := g.metrics()
	lm := g.labelSpace(m)
	adv := m.cell
	if adv == 0 {
		adv = g.c.fonts.RuneAdvance(m.spec, chars[0].r) * m.em
	}
	w := adv
	if m.stick {
		w = adv / 1.5 // the character, without the space after it
	}
	gl := glyph{r: chars[0].r, at: g.pos, adv: adv, m: m, lm: lm}
	g.drawRuns([]glyph{gl}, cad.Point{X: -w / 2, Y: -m.capH / 2}, g.penColorOf(g.pens.current()), cad.BreakBox)
}

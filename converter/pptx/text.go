package pptx

import (
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// Text properties are inherited along a chain: the paragraph's own a:pPr,
// the shape's list style, the list styles of the placeholders it inherits
// from, the master's text style for the placeholder type (or the
// presentation's default text style for other shapes). A property is taken
// from the first element in the chain that sets it.

type chain []*node

func (c chain) attr(name string) (string, bool) {
	for _, n := range c {
		if v, ok := n.attr(name); ok {
			return v, true
		}
	}
	return "", false
}

func (c chain) child(name string) *node {
	for _, n := range c {
		if k := n.child(name); k != nil {
			return k
		}
	}
	return nil
}

// first returns the first child among names found in the chain (the
// first element that has any of them wins).
func (c chain) first(names ...string) *node {
	for _, n := range c {
		if n == nil {
			continue
		}
		for _, k := range n.Kids {
			for _, name := range names {
				if k.Name == name {
					return k
				}
			}
		}
	}
	return nil
}

func (c chain) boolAttr(name string, def bool) bool {
	for _, n := range c {
		if _, ok := n.attr(name); ok {
			return n.attrBool(name, def)
		}
	}
	return def
}

// textFrame is a text body and the styles it inherits.
type textFrame struct {
	body    *node
	bodyPrs chain   // own bodyPr, then inherited ones
	own     *node   // the shape's own list style (may be nil)
	lists   []*node // inherited list styles, most specific first
	// extra run properties ranked right after the shape's own list style
	// (the shape style's font reference, a table style's text properties)
	extra []*node
	part  string
	cc    *colorCtx
}

var lvlNames = [...]string{"lvl1pPr", "lvl2pPr", "lvl3pPr", "lvl4pPr", "lvl5pPr", "lvl6pPr", "lvl7pPr", "lvl8pPr", "lvl9pPr"}

// pChain returns the paragraph property chain for a level. Its first three
// entries are always the paragraph's pPr and the own list style's level
// and default entries (any of them may be nil).
func (tf *textFrame) pChain(pPr *node, lvl int) chain {
	ch := chain{pPr, tf.own.child(lvlNames[lvl]), tf.own.child("defPPr")}
	for _, l := range tf.lists {
		if l != nil {
			ch = append(ch, l.child(lvlNames[lvl]), l.child("defPPr"))
		}
	}
	return ch
}

// rChain returns the run property chain: the run's rPr, then the default
// run properties of each paragraph chain entry, with the extra properties
// after the shape's own list style.
func (tf *textFrame) rChain(rPr *node, pc chain) chain {
	ch := chain{rPr}
	for i, p := range pc {
		ch = append(ch, p.child("defRPr"))
		if i == 2 {
			ch = append(ch, tf.extra...)
		}
	}
	return ch
}

// runStyle is a resolved set of character properties.
type runStyle struct {
	size         float64 // points
	bold, italic bool
	latin, ea    string
	sym          string
	fill         fill
	textFill     fill // fill before hyperlink styling (bullets follow it)
	outline      *line
	underline    string
	uColor       *rgba
	strike       string
	baseline     float64 // fraction of the size, positive raises
	spacing      float64 // points
	caps         string
	highlight    *rgba
	link         string
	lang         string
}

func (s *slideCtx) runStyle(tf *textFrame, rc chain, scale float64, rPr *node) *runStyle {
	st := &runStyle{size: 18, underline: "none", strike: "noStrike", caps: "none"}
	if v, ok := rc.attr("sz"); ok {
		st.size = atof(v, 1800) / 100
	}
	st.size *= scale
	st.bold = rc.boolAttr("b", false)
	st.italic = rc.boolAttr("i", false)
	st.underline, _ = rc.attr("u")
	if st.underline == "" {
		st.underline = "none"
	}
	if v, ok := rc.attr("strike"); ok {
		st.strike = v
	}
	if v, ok := rc.attr("baseline"); ok {
		st.baseline = parsePct(v)
	}
	if v, ok := rc.attr("spc"); ok {
		st.spacing = atof(v, 0) / 100 * scale
	}
	if v, ok := rc.attr("cap"); ok {
		st.caps = v
	}
	st.lang, _ = rc.attr("lang")
	altLang, _ := rc.attr("altLang")
	script := scriptOf(st.lang)
	if script == "" {
		script = scriptOf(altLang)
	}
	if script == "" {
		script = "Jpan"
	}
	face := func(name string) string {
		if k := rc.child(name); k != nil {
			return s.th.fontFor(k.attrStr("typeface", ""), script)
		}
		return ""
	}
	st.latin = face("latin")
	st.ea = face("ea")
	st.sym = face("sym")
	if st.latin == "" {
		st.latin = s.th.fontFor("+mn-lt", script)
	}
	if st.ea == "" {
		st.ea = s.th.fontFor("+mn-ea", script)
	}
	if e := rc.first("noFill", "solidFill", "gradFill", "blipFill", "pattFill", "grpFill"); e != nil {
		st.fill = s.resolveFill(e, tf.part, tf.cc, nil)
	} else {
		c, _ := tf.cc.scheme1("tx1")
		st.fill = fill{kind: fillSolid, color: c}
	}
	if ln := rc.child("ln"); ln != nil {
		st.outline = s.resolveLine([]*node{ln}, []string{tf.part}, nil, tf.cc)
	}
	if u := rc.child("uFill"); u != nil {
		if c, ok := tf.cc.color(u.child("solidFill")); ok {
			st.uColor = &c
		}
	}
	if h := rc.child("highlight"); h != nil {
		if c, ok := tf.cc.color(h); ok {
			st.highlight = &c
		}
	}
	st.textFill = st.fill
	if h := rPr.child("hlinkClick"); h != nil {
		st.link = s.linkTarget(h, tf.part)
		// Hyperlinks take the theme's hyperlink color and an underline
		// unless the run sets its own.
		if rPr.child("solidFill") == nil {
			if c, ok := tf.cc.scheme1("hlink"); ok {
				st.fill = fill{kind: fillSolid, color: c}
			}
		}
		if _, ok := rPr.attr("u"); !ok {
			st.underline = "sng"
		}
	}
	return st
}

func scriptOf(lang string) string {
	l := strings.ToLower(lang)
	switch {
	case strings.HasPrefix(l, "ja"):
		return "Jpan"
	case strings.HasPrefix(l, "zh-tw"), strings.HasPrefix(l, "zh-hk"), strings.HasPrefix(l, "zh-mo"), strings.HasPrefix(l, "zh-hant"):
		return "Hant"
	case strings.HasPrefix(l, "zh"):
		return "Hans"
	case strings.HasPrefix(l, "ko"):
		return "Hang"
	}
	return ""
}

// color returns the solid color text is drawn with (gradients use their
// average), or false for invisible text.
func (st *runStyle) color() (rgba, bool) {
	if st.fill.kind == fillNone {
		return rgba{}, false
	}
	return st.fill.average()
}

// para is a paragraph ready for line breaking.
type para struct {
	pc      chain
	lvl     int
	items   []item
	algn    string
	marL    float64
	indent  float64
	lnPct   float64 // line spacing as a fraction of the natural height (0 with lnPts)
	lnPts   float64
	befPct  float64
	befPts  float64
	aftPct  float64
	aftPts  float64
	defTab  float64
	tabs    []float64
	bullet  *bullet
	end     *runStyle // paragraph end properties (for empty paragraphs)
	endFace *faceChoice
	hasText bool
}

// item is one character (or tab / line break) with its style and face.
type item struct {
	r    rune
	st   *runStyle
	fc   *faceChoice
	w    float64 // advance in points, letter spacing included
	brk  bool    // a line may break after this item
	kind byte    // itemChar, itemTab, itemBreak
	x    float64 // position in the line (set by layout)
	ls   float64 // extra letter spacing from justification
}

const (
	itemChar = iota
	itemTab
	itemBreak
)

type bullet struct {
	text string
	st   *runStyle
	fc   *faceChoice
	w    float64
}

// paragraphs builds the paragraphs of a text body.
func (s *slideCtx) paragraphs(tf *textFrame, fontScale, lnReduce float64) []*para {
	var out []*para
	counters := [9]int{}
	schemes := [9]string{}
	for _, p := range tf.body.children("p") {
		pPr := p.child("pPr")
		lvl := int(pPr.attrInt("lvl", 0))
		if lvl < 0 || lvl > 8 {
			lvl = 0
		}
		pc := tf.pChain(pPr, lvl)
		pa := &para{pc: pc, lvl: lvl}
		pa.algn, _ = pc.attr("algn")
		if pa.algn == "" {
			pa.algn = "l"
		}
		pa.marL = emuStr(pc, "marL")
		pa.indent = emuStr(pc, "indent")
		pa.lnPct = 1
		if ls := pc.child("lnSpc"); ls != nil {
			if v := ls.child("spcPct"); v != nil {
				pa.lnPct = pct(v, "val", 1)
			} else if v := ls.child("spcPts"); v != nil {
				pa.lnPct, pa.lnPts = 0, float64(v.attrInt("val", 0))/100
			}
		}
		if pa.lnPct > 0 {
			pa.lnPct = max(pa.lnPct-lnReduce, 0.1)
		}
		if sp := pc.child("spcBef"); sp != nil {
			pa.befPct, pa.befPts = spacingOf(sp)
		}
		if sp := pc.child("spcAft"); sp != nil {
			pa.aftPct, pa.aftPts = spacingOf(sp)
		}
		pa.defTab = 72
		if v, ok := pc.attr("defTabSz"); ok {
			if t := atof(v, 914400) / emuPerPt; t > 0 {
				pa.defTab = t
			}
		}
		if tl := pc.child("tabLst"); tl != nil {
			for _, t := range tl.children("tab") {
				pa.tabs = append(pa.tabs, t.emuAttr("pos", 0))
			}
		}
		// runs
		for _, k := range p.Kids {
			switch k.Name {
			case "r", "fld":
				rPr := k.child("rPr")
				st := s.runStyle(tf, tf.rChain(rPr, pc), fontScale, rPr)
				text := k.child("t").Text
				if k.Name == "fld" && strings.HasPrefix(k.attrStr("type", ""), "slidenum") {
					text = itoa(s.num)
				}
				s.addText(pa, st, text)
			case "br":
				rPr := k.child("rPr")
				st := s.runStyle(tf, tf.rChain(rPr, pc), fontScale, rPr)
				pa.items = append(pa.items, item{r: '\n', st: st, fc: s.c.faceFor(st, ' '), kind: itemBreak})
			}
		}
		endPr := p.child("endParaRPr")
		pa.end = s.runStyle(tf, tf.rChain(endPr, pc), fontScale, endPr)
		pa.endFace = s.c.faceFor(pa.end, 'x')
		for _, it := range pa.items {
			if it.kind == itemChar {
				pa.hasText = true
				break
			}
		}
		// bullets and numbering
		bu := pc.first("buNone", "buAutoNum", "buChar", "buBlip")
		if pa.hasText {
			for l := lvl + 1; l < 9; l++ {
				counters[l], schemes[l] = 0, ""
			}
			if bu != nil && bu.Name == "buAutoNum" {
				scheme := bu.attrStr("type", "arabicPeriod")
				start := int(bu.attrInt("startAt", 1))
				if schemes[lvl] != scheme || counters[lvl] == 0 {
					counters[lvl] = start
				} else {
					counters[lvl]++
				}
				schemes[lvl] = scheme
				pa.bullet = s.makeBullet(tf, pa, autoNumber(scheme, counters[lvl]), "", fontScale)
			} else {
				counters[lvl], schemes[lvl] = 0, ""
				if bu != nil && (bu.Name == "buChar" || bu.Name == "buBlip") {
					ch := bu.attrStr("char", "•")
					if bu.Name == "buBlip" {
						ch = "•"
					}
					var font string
					if f := pc.first("buFont", "buFontTx"); f != nil && f.Name == "buFont" {
						font = f.attrStr("typeface", "")
					}
					pa.bullet = s.makeBullet(tf, pa, ch, font, fontScale)
				}
			}
		}
		out = append(out, pa)
	}
	return out
}

func spacingOf(sp *node) (p, pts float64) {
	if v := sp.child("spcPct"); v != nil {
		return pct(v, "val", 0), 0
	}
	if v := sp.child("spcPts"); v != nil {
		return 0, float64(v.attrInt("val", 0)) / 100
	}
	return 0, 0
}

func emuStr(c chain, name string) float64 {
	if v, ok := c.attr(name); ok {
		return atof(v, 0) / emuPerPt
	}
	return 0
}

// addText appends the characters of a run.
func (s *slideCtx) addText(pa *para, st *runStyle, text string) {
	switch st.caps {
	case "all":
		text = strings.ToUpper(text)
	case "small":
		// small capitals: lower-case letters become smaller capitals
		small := *st
		small.size = st.size * 0.8
		var buf []rune
		bufSmall := false
		flush := func() {
			if len(buf) == 0 {
				return
			}
			if bufSmall {
				s.addChars(pa, &small, string(buf))
			} else {
				s.addChars(pa, st, string(buf))
			}
			buf = buf[:0]
		}
		for _, r := range text {
			lower := unicode.IsLower(r)
			if lower != bufSmall {
				flush()
				bufSmall = lower
			}
			if lower {
				r = unicode.ToUpper(r)
			}
			buf = append(buf, r)
		}
		flush()
		return
	}
	s.addChars(pa, st, text)
}

func (s *slideCtx) addChars(pa *para, st *runStyle, text string) {
	size := st.size
	if st.baseline != 0 {
		size *= 2.0 / 3
	}
	for _, r := range text {
		switch r {
		case '\r':
			continue
		case '\v':
			pa.items = append(pa.items, item{r: '\n', st: st, fc: s.c.faceFor(st, ' '), kind: itemBreak})
			continue
		case '\t':
			pa.items = append(pa.items, item{r: '\t', st: st, fc: s.c.faceFor(st, ' '), kind: itemTab, brk: true})
			continue
		}
		fc := s.c.faceFor(st, r)
		w := s.c.advance(fc, r)*size + st.spacing
		pa.items = append(pa.items, item{r: r, st: st, fc: fc, w: w})
	}
}

// makeBullet builds the bullet of a paragraph from its first run.
func (s *slideCtx) makeBullet(tf *textFrame, pa *para, text, font string, fontScale float64) *bullet {
	var first *runStyle
	for _, it := range pa.items {
		if it.kind == itemChar {
			first = it.st
			break
		}
	}
	if first == nil {
		return nil
	}
	st := *first
	st.fill = first.textFill
	st.underline, st.strike, st.link, st.highlight, st.spacing, st.baseline = "none", "noStrike", "", nil, 0, 0
	pc := pa.pc
	if f := pc.first("buFont", "buFontTx"); f != nil && f.Name == "buFont" {
		font = s.th.fontFor(f.attrStr("typeface", ""), "Jpan")
	}
	if font != "" {
		text = mapSymbolFont(font, text)
		if isSymbolFont(font) {
			font = ""
		}
	}
	if font != "" {
		st.latin, st.ea = font, font
	}
	if sz := pc.first("buSzPct", "buSzPts", "buSzTx"); sz != nil {
		switch sz.Name {
		case "buSzPct":
			st.size = first.size * pct(sz, "val", 1)
		case "buSzPts":
			st.size = float64(sz.attrInt("val", 1800)) / 100 * fontScale
		}
	}
	if cl := pc.first("buClr", "buClrTx"); cl != nil && cl.Name == "buClr" {
		if c, ok := tf.cc.color(cl); ok {
			st.fill = fill{kind: fillSolid, color: c}
		}
	}
	b := &bullet{text: text, st: &st}
	for _, r := range text {
		b.fc = s.c.faceFor(&st, r)
		b.w += s.c.advance(b.fc, r) * st.size
	}
	return b
}

func isSymbolFont(name string) bool {
	switch fontdb.Normalize(name) {
	case "wingdings", "wingdings2", "wingdings3", "symbol", "webdings":
		return true
	}
	return false
}

// mapSymbolFont maps the characters of the symbol fonts that bullets
// commonly use to their Unicode equivalents.
func mapSymbolFont(font, text string) string {
	n := fontdb.Normalize(font)
	var table map[rune]rune
	switch n {
	case "wingdings":
		table = wingdings
	case "symbol":
		table = symbolFont
	default:
		return text
	}
	var out []rune
	for _, r := range text {
		if r >= 0xf000 && r <= 0xf0ff {
			r -= 0xf000
		}
		if m, ok := table[r]; ok {
			r = m
		}
		out = append(out, r)
	}
	return string(out)
}

var wingdings = map[rune]rune{
	0x6c: '●', 0x6d: '❍', 0x6e: '■', 0x6f: '□', 0x70: '◻', 0x71: '❑', 0x72: '❒', 0x73: '⬧', 0x74: '⧫',
	0x75: '◆', 0x76: '❖', 0x77: '⬥', 0x78: '⌧', 0xa7: '▪', 0xa8: '◻', 0x9f: '•', 0xa1: '○', 0xa2: '⭕',
	0xd8: '➢', 0xe0: '➔', 0xe8: '➔', 0xfc: '✔', 0xfb: '✘', 0xfd: '☒', 0xfe: '☑', 0xf0: '⇨', 0xef: '⇦',
	0xd9: '⮙', 0xda: '⮛', 0x46: '☞', 0x4a: '☺', 0x4c: '☹', 0xab: '★', 0xb2: '✧', 0x2a: '✉', 0x3e: '✇',
	0x3f: '✍', 0x41: '✌', 0x43: '👍', 0x44: '👎', 0x51: '✈', 0x52: '☼', 0x54: '❄',
}

var symbolFont = map[rune]rune{
	0xb7: '•', 0xa8: '♣', 0xa9: '♦', 0xaa: '♥', 0xab: '♠', 0xae: '→', 0xde: '⇒', 0xd8: '¬', 0x2d: '−',
	0xb0: '°', 0xb1: '±', 0xb4: '×', 0xb8: '÷', 0xa5: '∞', 0x61: 'α', 0x62: 'β', 0x67: 'γ', 0x64: 'δ',
	0x70: 'π', 0x6d: 'μ', 0x53: 'Σ', 0x57: 'Ω', 0xe0: '◊',
}

// autoNumber formats a bullet number in an ST_TextAutonumberScheme.
func autoNumber(scheme string, n int) string {
	var num string
	switch {
	case strings.HasPrefix(scheme, "alphaLc"):
		num = strings.ToLower(alpha(n))
	case strings.HasPrefix(scheme, "alphaUc"):
		num = alpha(n)
	case strings.HasPrefix(scheme, "romanLc"):
		num = strings.ToLower(roman(n))
	case strings.HasPrefix(scheme, "romanUc"):
		num = roman(n)
	case strings.HasPrefix(scheme, "circleNum"):
		if n >= 1 && n <= 20 {
			if strings.Contains(scheme, "WdBlack") {
				if n <= 10 {
					return string(rune(0x2776 + n - 1))
				}
				return string(rune(0x24eb + n - 11))
			}
			return string(rune(0x2460 + n - 1))
		}
		num = itoa(n)
	case strings.HasPrefix(scheme, "arabicDb"):
		num = toFullWidth(itoa(n))
	case strings.HasPrefix(scheme, "ea1Jpn") || strings.HasPrefix(scheme, "ea1Chs") || strings.HasPrefix(scheme, "ea1Cht"):
		num = kanjiNumber(n)
	default:
		num = itoa(n)
	}
	switch {
	case strings.HasSuffix(scheme, "ParenBoth"):
		return "(" + num + ")"
	case strings.HasSuffix(scheme, "ParenR"):
		return num + ")"
	case strings.HasSuffix(scheme, "Period"):
		if strings.HasPrefix(scheme, "arabicDb") || strings.HasPrefix(scheme, "ea1") {
			return num + "．"
		}
		return num + "."
	case strings.HasSuffix(scheme, "Minus"):
		return "- " + num + " -"
	case strings.HasSuffix(scheme, "Comma"):
		return num + "、"
	}
	return num
}

func alpha(n int) string {
	if n < 1 {
		return itoa(n)
	}
	// A..Z, AA..ZZ, AAA… (repeated letters, as Office numbers them)
	letter := rune('A' + (n-1)%26)
	return strings.Repeat(string(letter), (n-1)/26+1)
}

func roman(n int) string {
	if n < 1 || n > 3999 {
		return itoa(n)
	}
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	var b strings.Builder
	for i, v := range vals {
		for n >= v {
			b.WriteString(syms[i])
			n -= v
		}
	}
	return b.String()
}

func toFullWidth(s string) string {
	var out []rune
	for _, r := range s {
		if r >= '0' && r <= '9' {
			r += 0xff10 - '0'
		}
		out = append(out, r)
	}
	return string(out)
}

func kanjiNumber(n int) string {
	digits := []rune("〇一二三四五六七八九")
	if n < 0 || n > 99 {
		return itoa(n)
	}
	if n < 10 {
		return string(digits[n])
	}
	s := ""
	if n/10 > 1 {
		s = string(digits[n/10])
	}
	s += "十"
	if n%10 != 0 {
		s += string(digits[n%10])
	}
	return s
}

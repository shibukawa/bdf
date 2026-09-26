package drawingml

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Text properties are inherited along a chain: the paragraph's own a:pPr,
// the shape's list style, the list styles of the placeholders it inherits
// from, and the host's text style (Host.TextStyle: in a presentation the
// master's text style for the placeholder type, or the default text style
// for other shapes). A property is taken from the first element in the
// chain that sets it.

type chain []*ooxml.Node

func (c chain) attr(name string) (string, bool) {
	for _, n := range c {
		if v, ok := n.Attr(name); ok {
			return v, true
		}
	}
	return "", false
}

func (c chain) child(name string) *ooxml.Node {
	for _, n := range c {
		if k := n.Child(name); k != nil {
			return k
		}
	}
	return nil
}

// first returns the first child among names found in the chain (the
// first element that has any of them wins).
func (c chain) first(names ...string) *ooxml.Node {
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
		if _, ok := n.Attr(name); ok {
			return n.AttrBool(name, def)
		}
	}
	return def
}

// textFrame is a text body and the styles it inherits.
type textFrame struct {
	body    *ooxml.Node
	phType  string        // placeholder type ("" for other shapes)
	bodyPrs chain         // own bodyPr, then inherited ones
	own     *ooxml.Node   // the shape's own list style (may be nil)
	lists   []*ooxml.Node // inherited list styles, most specific first
	// extra run properties ranked right after the shape's own list style
	// (the shape style's font reference, a table style's text properties)
	extra []*ooxml.Node
	part  string
	cc    *colorCtx
}

var lvlNames = [...]string{"lvl1pPr", "lvl2pPr", "lvl3pPr", "lvl4pPr", "lvl5pPr", "lvl6pPr", "lvl7pPr", "lvl8pPr", "lvl9pPr"}

// pChain returns the paragraph property chain for a level. Its first three
// entries are always the paragraph's pPr and the own list style's level
// and default entries (any of them may be nil).
func (tf *textFrame) pChain(pPr *ooxml.Node, lvl int) chain {
	ch := chain{pPr, tf.own.Child(lvlNames[lvl]), tf.own.Child("defPPr")}
	for _, l := range tf.lists {
		if l != nil {
			ch = append(ch, l.Child(lvlNames[lvl]), l.Child("defPPr"))
		}
	}
	return ch
}

// rChain returns the run property chain: the run's rPr, then the default
// run properties of each paragraph chain entry, with the extra properties
// after the shape's own list style.
func (tf *textFrame) rChain(rPr *ooxml.Node, pc chain) chain {
	ch := chain{rPr}
	for i, p := range pc {
		ch = append(ch, p.Child("defRPr"))
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
	lang         string // not part of key: runs split by language are drawn as one
	altLang      string
	key          string // identity of the visible properties, for merging runs
}

func (s *Drawing) runStyle(tf *textFrame, rc chain, scale float64, rPr *ooxml.Node) *runStyle {
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
	st.altLang, _ = rc.attr("altLang")
	script := fontset.Script(st.lang)
	if script == "" {
		script = fontset.Script(st.altLang)
	}
	if script == "" {
		script = "Jpan"
	}
	face := func(name string) string {
		if k := rc.child(name); k != nil {
			return s.th.fontFor(k.AttrStr("typeface", ""), script)
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
		st.outline = s.resolveLine([]*ooxml.Node{ln}, []string{tf.part}, nil, tf.cc)
	}
	if u := rc.child("uFill"); u != nil {
		if c, ok := tf.cc.color(u.Child("solidFill")); ok {
			st.uColor = &c
		}
	}
	if h := rc.child("highlight"); h != nil {
		if c, ok := tf.cc.color(h); ok {
			st.highlight = &c
		}
	}
	defer st.setKey()
	st.textFill = st.fill
	if h := rPr.Child("hlinkClick"); h != nil {
		st.link = s.linkTarget(h, tf.part)
		// Hyperlinks take the theme's hyperlink color and an underline
		// unless the run sets its own.
		if rPr.Child("solidFill") == nil {
			if c, ok := tf.cc.scheme1("hlink"); ok {
				st.fill = fill{kind: fillSolid, color: c}
			}
		}
		if _, ok := rPr.Attr("u"); !ok {
			st.underline = "sng"
		}
	}
	return st
}

// setKey records the properties that affect drawing, so that runs split in
// the markup (spell checking, language tags) are drawn as one.
func (st *runStyle) setKey() {
	c, _ := st.fill.average()
	var b strings.Builder
	fmt.Fprintf(&b, "%g|%v|%v|%s|%s|%s|%d|%v|%s|%s|%g|%g|%s|%s|%v", st.size, st.bold, st.italic, st.latin, st.ea, st.sym,
		st.fill.kind, c, st.underline, st.strike, st.baseline, st.spacing, st.caps, st.link, st.outline != nil)
	if st.uColor != nil {
		fmt.Fprintf(&b, "|u%v", *st.uColor)
	}
	if st.highlight != nil {
		fmt.Fprintf(&b, "|h%v", *st.highlight)
	}
	if st.outline != nil {
		fmt.Fprintf(&b, "|o%p", st.outline)
	}
	if st.fill.kind == fillGrad {
		fmt.Fprintf(&b, "|g%p", st.fill.grad)
	}
	st.key = b.String()
}

// faceFor picks the face that draws r in a run (see fontset.Set.FaceFor).
func (c *Renderer) faceFor(st *runStyle, r rune) *fontset.Choice {
	return c.fonts.FaceFor(st.latin, st.ea, st.bold, st.italic, r)
}

// runLang returns the language of a run of text: its lang attribute, and
// for East Asian text an East Asian lang or altLang, or else the language
// that its kana or hangul, or those of its paragraph pa (may be nil), imply.
// It returns false for East Asian text that nothing assigns a language:
// Han characters alone could be Chinese as well as Japanese or Korean.
func runLang(run []item, pa *para) (string, bool) {
	var st *runStyle
	for _, it := range run {
		if it.kind == itemChar && fontdb.IsCJK(it.r) {
			st = it.st
			break
		}
	}
	switch {
	case st == nil:
		return run[0].st.lang, true
	case fontset.Script(st.lang) != "":
		return st.lang, true
	case fontset.Script(st.altLang) != "":
		return st.altLang, true
	}
	if l := scriptLang(run); l != "" {
		return l, true
	}
	if pa != nil && pa.eaLang != "" {
		return pa.eaLang, true
	}
	return "", false
}

// scriptLang returns "ja" for text with kana and "ko" for text with hangul.
func scriptLang(items []item) string {
	for _, it := range items {
		if it.kind != itemChar {
			continue
		}
		if l := fontset.RuneLang(it.r); l != "" {
			return l
		}
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
	endFace *fontset.Choice
	hasText bool
	eaLang  string // language its kana or hangul imply (see scriptLang)
}

// item is one character (or tab / line break) with its style and face.
type item struct {
	r    rune
	st   *runStyle
	fc   *fontset.Choice
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
	fc   *fontset.Choice
	w    float64
}

// paragraphs builds the paragraphs of a text body.
func (s *Drawing) paragraphs(tf *textFrame, fontScale, lnReduce float64) []*para {
	var out []*para
	counters := [9]int{}
	schemes := [9]string{}
	for _, p := range tf.body.Children("p") {
		pPr := p.Child("pPr")
		lvl := int(pPr.AttrInt("lvl", 0))
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
			if v := ls.Child("spcPct"); v != nil {
				pa.lnPct = v.AttrPct("val", 1)
			} else if v := ls.Child("spcPts"); v != nil {
				pa.lnPct, pa.lnPts = 0, float64(v.AttrInt("val", 0))/100
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
			if t := atof(v, 914400) / ooxml.EMUPerPoint; t > 0 {
				pa.defTab = t
			}
		}
		if tl := pc.child("tabLst"); tl != nil {
			for _, t := range tl.Children("tab") {
				pa.tabs = append(pa.tabs, t.AttrEMU("pos", 0))
			}
		}
		// runs
		for _, k := range p.Kids {
			switch k.Name {
			case "r", "fld":
				rPr := k.Child("rPr")
				st := s.runStyle(tf, tf.rChain(rPr, pc), fontScale, rPr)
				text := k.Child("t").Content()
				if k.Name == "fld" {
					if t, ok := s.host.Field(k.AttrStr("type", "")); ok {
						text = t
					}
				}
				s.addText(pa, st, text)
			case "br":
				rPr := k.Child("rPr")
				st := s.runStyle(tf, tf.rChain(rPr, pc), fontScale, rPr)
				pa.items = append(pa.items, item{r: '\n', st: st, fc: s.c.faceFor(st, ' '), kind: itemBreak})
			}
		}
		endPr := p.Child("endParaRPr")
		pa.end = s.runStyle(tf, tf.rChain(endPr, pc), fontScale, endPr)
		pa.endFace = s.c.faceFor(pa.end, 'x')
		for _, it := range pa.items {
			if it.kind == itemChar {
				pa.hasText = true
				break
			}
		}
		pa.eaLang = scriptLang(pa.items)
		// bullets and numbering
		bu := pc.first("buNone", "buAutoNum", "buChar", "buBlip")
		if pa.hasText {
			for l := lvl + 1; l < 9; l++ {
				counters[l], schemes[l] = 0, ""
			}
			if bu != nil && bu.Name == "buAutoNum" {
				scheme := bu.AttrStr("type", "arabicPeriod")
				start := int(bu.AttrInt("startAt", 1))
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
					ch := bu.AttrStr("char", "•")
					if bu.Name == "buBlip" {
						ch = "•"
					}
					var font string
					if f := pc.first("buFont", "buFontTx"); f != nil && f.Name == "buFont" {
						font = f.AttrStr("typeface", "")
					}
					pa.bullet = s.makeBullet(tf, pa, ch, font, fontScale)
				}
			}
		}
		out = append(out, pa)
	}
	return out
}

func spacingOf(sp *ooxml.Node) (p, pts float64) {
	if v := sp.Child("spcPct"); v != nil {
		return v.AttrPct("val", 0), 0
	}
	if v := sp.Child("spcPts"); v != nil {
		return 0, float64(v.AttrInt("val", 0)) / 100
	}
	return 0, 0
}

func emuStr(c chain, name string) float64 {
	if v, ok := c.attr(name); ok {
		return atof(v, 0) / ooxml.EMUPerPoint
	}
	return 0
}

// addText appends the characters of a run.
func (s *Drawing) addText(pa *para, st *runStyle, text string) {
	switch st.caps {
	case "all":
		text = strings.ToUpper(text)
	case "small":
		// small capitals: lower-case letters become smaller capitals
		small := *st
		small.size = st.size * 0.8
		small.setKey()
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

func (s *Drawing) addChars(pa *para, st *runStyle, text string) {
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
		w := s.c.fonts.Advance(fc, r)*size + st.spacing
		pa.items = append(pa.items, item{r: r, st: st, fc: fc, w: w})
	}
}

// makeBullet builds the bullet of a paragraph from its first run.
func (s *Drawing) makeBullet(tf *textFrame, pa *para, text, font string, fontScale float64) *bullet {
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
		font = s.th.fontFor(f.AttrStr("typeface", ""), "Jpan")
	}
	if font != "" {
		text = fontset.MapSymbol(font, text)
		if fontset.IsSymbol(font) {
			font = ""
		}
	}
	if font != "" {
		st.latin, st.ea = font, font
	}
	if sz := pc.first("buSzPct", "buSzPts", "buSzTx"); sz != nil {
		switch sz.Name {
		case "buSzPct":
			st.size = first.size * sz.AttrPct("val", 1)
		case "buSzPts":
			st.size = float64(sz.AttrInt("val", 1800)) / 100 * fontScale
		}
	}
	if cl := pc.first("buClr", "buClrTx"); cl != nil && cl.Name == "buClr" {
		if c, ok := tf.cc.color(cl); ok {
			st.fill = fill{kind: fillSolid, color: c}
		}
	}
	st.setKey()
	b := &bullet{text: text, st: &st}
	for _, r := range text {
		b.fc = s.c.faceFor(&st, r)
		b.w += s.c.fonts.Advance(b.fc, r) * st.size
	}
	return b
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

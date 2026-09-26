package drawio

import (
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// Label text layout. draw.io renders HTML labels as HTML in a
// foreignObject: an inline-block of line-height 1.2 inside a flex box
// (mxSvgCanvas2D.createCss). The layout here reproduces the parts of CSS
// that labels use: blocks with collapsing margins, list indents and
// markers, inline styles, white-space collapsing, line breaking at spaces
// and between East Asian characters, and line boxes sized from the fonts'
// ascent and descent with the half-leading of the line height.

// tstyle is the computed style of inline text (CSS values in px).
type tstyle struct {
	families          []string
	size              float64
	bold, italic      bool
	underline, strike bool
	color             rgba
	bg                *rgba   // background-color of an inline element
	shift             float64 // baseline shift (positive: down), for sub and sup
	link              string
	pre               bool        // white-space: pre / pre-wrap / break-spaces
	nowrap            bool        // white-space: nowrap / pre
	lhFactor          float64     // line-height as a multiple of the font size (0 with lhPx)
	lhPx              float64     // line-height in px
	align             string      // text-align of the block (inherited)
	primary           *faceChoice // the face of the first family, for line metrics
}

func (s *tstyle) clone() *tstyle { c := *s; return &c }

// lineHeight returns the computed line height for text of this style.
func (s *tstyle) lineHeight() float64 {
	if s.lhPx > 0 {
		return s.lhPx
	}
	f := s.lhFactor
	if f <= 0 {
		f = 1.2
	}
	return f * s.size
}

// titem is one character, or a forced line break, of a paragraph.
type titem struct {
	r     rune
	st    *tstyle
	fc    *faceChoice
	w     float64 // advance in px
	brk   bool    // a line may break after this item
	hard  bool    // a forced line break (<br>)
	space bool    // collapsible white space
}

// tpara is a block's inline content.
type tpara struct {
	items   []titem
	st      *tstyle // the block's style (its strut)
	align   string
	indent  float64 // left padding in px (lists, blockquotes)
	top     float64 // space before: the collapsed margins
	marker  []titem // outside list marker text (ordered lists)
	bullet  string  // outside list marker symbol: disc, circle or square
	heading int
	marks   []paraMark // structure MARKs before the paragraph
	hr      bool       // a horizontal rule, no text
	hrColor rgba
}

// paraMark is a structure mark to emit before a paragraph.
type paraMark struct {
	kind    byte
	payload string
}

// tline is a laid-out line.
type tline struct {
	para     *tpara
	items    []titem
	first    bool    // first line of its paragraph
	hardPrev bool    // the previous line ended with a forced break
	cjkWrap  bool    // the previous line wrapped between East Asian characters (no separator)
	x        float64 // left of the line in the content box
	width    float64
	top      float64
	asc      float64 // from the top of the line to the baseline
	height   float64
}

// tlayout is laid-out label text.
type tlayout struct {
	paras   []*tpara
	lines   []*tline
	width   float64 // shrink-to-fit width of the content
	height  float64
	wrapped bool
}

// textBuilder turns label HTML (or plain text) into paragraphs.
type textBuilder struct {
	c       *converter
	paras   []*tpara
	cur     *tpara
	pending float64 // collapsed margin before the next paragraph
	lists   []listState
	space   bool // the last item added is a collapsible space
	marks   []paraMark
	heading int
	indent  float64
	liNext  bool // the next paragraph starts a list item
}

type listState struct {
	ordered bool
	n       int
	typ     string
}

func (b *textBuilder) para(st *tstyle) *tpara {
	if b.cur == nil {
		b.cur = &tpara{st: st, align: st.align, indent: b.indent, top: b.pending, heading: b.heading, marks: b.marks}
		b.marks = nil
		b.pending = 0
		b.space = true // leading spaces of a block are dropped
		if b.liNext {
			b.cur.marker, b.cur.bullet = b.makeMarker(st)
			b.liNext = false
		}
		b.paras = append(b.paras, b.cur)
	}
	return b.cur
}

// endPara closes the current paragraph (a block boundary).
func (b *textBuilder) endPara() {
	if b.cur != nil {
		b.cur = nil
	}
	b.space = true
}

// margin collapses a vertical margin into the pending space.
func (b *textBuilder) margin(m float64) {
	b.endPara()
	b.pending = math.Max(b.pending, m)
}

// text adds text in a style, collapsing white space unless it is preformatted.
func (b *textBuilder) text(s string, st *tstyle) {
	for _, r := range s {
		if !st.pre && (r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f') {
			if b.space {
				continue
			}
			b.add(' ', st, true)
			b.space = true
			continue
		}
		if st.pre && r == '\n' {
			b.hardBreak(st)
			continue
		}
		if r == '\r' {
			continue
		}
		if r == '\t' {
			r = ' '
		}
		b.add(r, st, false)
		b.space = false
	}
}

func (b *textBuilder) add(r rune, st *tstyle, space bool) {
	p := b.para(st)
	fc := b.c.faceFor(st.families, st.bold, st.italic, r)
	w := b.c.advance(fc, r) * st.size
	p.items = append(p.items, titem{r: r, st: st, fc: fc, w: w, space: space})
}

func (b *textBuilder) hardBreak(st *tstyle) {
	p := b.para(st)
	fc := b.c.faceFor(st.families, st.bold, st.italic, ' ')
	p.items = append(p.items, titem{r: '\n', st: st, fc: fc, hard: true})
	b.space = true
}

// makeMarker builds the marker of a list item: text for ordered lists, a
// symbol (drawn, as Blink does) for unordered ones.
func (b *textBuilder) makeMarker(st *tstyle) ([]titem, string) {
	if len(b.lists) == 0 {
		return nil, ""
	}
	l := &b.lists[len(b.lists)-1]
	l.n++
	switch {
	case l.typ == "none":
		return nil, ""
	case !l.ordered:
		switch {
		case l.typ == "circle" || l.typ == "" && len(b.lists) == 2:
			return nil, "circle"
		case l.typ == "square" || l.typ == "" && len(b.lists) >= 3:
			return nil, "square"
		}
		return nil, "disc"
	}
	var out []titem
	for _, r := range listNumber(l.typ, l.n) + ". " {
		fc := b.c.faceFor(st.families, st.bold, st.italic, r)
		out = append(out, titem{r: r, st: st, fc: fc, w: b.c.advance(fc, r) * st.size})
	}
	return out, ""
}

// listNumber formats an ordered list number for a list-style-type.
func listNumber(typ string, n int) string {
	switch typ {
	case "a", "lower-alpha", "lower-latin":
		return alphaNumber(n, 'a')
	case "A", "upper-alpha", "upper-latin":
		return alphaNumber(n, 'A')
	case "i", "lower-roman":
		return strings.ToLower(romanNumber(n))
	case "I", "upper-roman":
		return romanNumber(n)
	}
	return strconv.Itoa(n)
}

func alphaNumber(n int, base rune) string {
	var out []rune
	for n > 0 {
		n--
		out = append([]rune{base + rune(n%26)}, out...)
		n /= 26
	}
	return string(out)
}

func romanNumber(n int) string {
	if n <= 0 || n >= 4000 {
		return strconv.Itoa(n)
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

// blockTags are the elements laid out as blocks.
var blockTags = map[string]bool{
	"div": true, "p": true, "ul": true, "ol": true, "li": true, "h1": true, "h2": true, "h3": true, "h4": true,
	"h5": true, "h6": true, "blockquote": true, "pre": true, "table": true, "tr": true, "center": true,
	"section": true, "article": true, "header": true, "footer": true, "nav": true, "aside": true, "main": true,
	"dl": true, "dt": true, "dd": true, "figure": true, "figcaption": true, "address": true, "form": true,
	"fieldset": true, "tbody": true, "thead": true, "tfoot": true, "caption": true, "hr": true,
}

var headingSizes = map[string]struct{ size, margin float64 }{
	"h1": {2, 0.67}, "h2": {1.5, 0.83}, "h3": {1.17, 1}, "h4": {1, 1.33}, "h5": {0.83, 1.67}, "h6": {0.67, 2.33},
}

// htmlFontSizes are the sizes of <font size=1..7>.
var htmlFontSizes = []float64{10, 13, 16, 18, 24, 32, 48}

// walk lays out an element's children.
func (b *textBuilder) walk(n *hnode, st *tstyle) {
	for _, k := range n.kids {
		if k.tag == "" {
			b.text(k.text, st)
			continue
		}
		b.element(k, st)
	}
}

func (b *textBuilder) element(n *hnode, parent *tstyle) {
	switch n.tag {
	case "script", "style", "head", "title", "template":
		return
	case "br":
		b.hardBreak(parent)
		return
	case "img":
		if !b.c.warned["htmlimg"] {
			b.c.warnOnce("htmlimg", "images inside HTML labels are not drawn")
		}
		return
	case "wbr":
		if b.cur != nil && len(b.cur.items) > 0 {
			b.cur.items[len(b.cur.items)-1].brk = true
		}
		return
	}
	st := parent.clone()
	st.bg = nil // backgrounds are not inherited
	st.shift = parent.shift
	css := parseCSS(n.attrs["style"])
	block := blockTags[n.tag]
	var marginTop, marginBottom, padLeft float64
	em := parent.size
	switch n.tag {
	case "b", "strong":
		st.bold = true
	case "i", "em", "cite", "var", "dfn", "address":
		st.italic = true
	case "u", "ins":
		st.underline = true
	case "s", "strike", "del":
		st.strike = true
	case "sub":
		st.size = parent.size * 0.83
		st.shift += parent.size * 0.2
	case "sup":
		st.size = parent.size * 0.83
		st.shift -= parent.size * 0.35
	case "small":
		st.size = parent.size * 0.83
	case "big":
		st.size = parent.size * 1.2
	case "code", "kbd", "samp", "tt", "pre":
		st.families = []string{fontdb.Mono}
		st.size = parent.size * 13 / 16
		if n.tag == "pre" {
			st.pre = true
			st.nowrap = true
			marginTop, marginBottom = em, em
		}
	case "td", "th":
		// table cells: separated by a space on the row's line
		if b.cur != nil && len(b.cur.items) > 0 && !b.space {
			b.add(' ', parent, true)
			b.space = true
		}
		st.bold = st.bold || n.tag == "th"
	case "a":
		if href := n.attrs["href"]; href != "" {
			st.link = href
			st.color = rgba{0, 0, 0xee, 255}
			st.underline = true
		}
	case "font":
		if v, ok := n.attrs["color"]; ok {
			if col, ok := parseColor(v); ok {
				st.color = col
			}
		}
		if v, ok := n.attrs["face"]; ok {
			if f := parseFontFamilies(v); len(f) > 0 {
				st.families = f
			}
		}
		if v, ok := n.attrs["size"]; ok {
			if sz, ok := fontSizeAttr(v); ok {
				st.size = sz
			}
		}
	case "p":
		marginTop, marginBottom = em, em
	case "blockquote":
		marginTop, marginBottom = em, em
		padLeft = 40
	case "ul", "ol":
		if len(b.lists) == 0 {
			marginTop, marginBottom = em, em
		}
		padLeft = 40
	case "dd":
		padLeft = 40
	case "dl":
		marginTop, marginBottom = em, em
	case "center":
		st.align = "center"
	case "h1", "h2", "h3", "h4", "h5", "h6":
		h := headingSizes[n.tag]
		st.size = parent.size * h.size
		st.bold = true
		marginTop, marginBottom = st.size*h.margin, st.size*h.margin
	case "hr":
		b.margin(em * 0.5)
		col := rgba{0xee, 0xee, 0xee, 255} // Chrome's inset border
		b.paras = append(b.paras, &tpara{st: st, hr: true, hrColor: col, top: b.pending, indent: b.indent})
		b.pending = 0
		b.margin(em * 0.5)
		return
	}
	if v, ok := n.attrs["align"]; ok && block {
		st.align = strings.ToLower(v)
	}
	b.applyCSS(css, st, parent, &marginTop, &marginBottom, &padLeft)
	b.c.setPrimary(st)

	if !block {
		b.walk(n, st)
		return
	}
	// block element
	switch n.tag {
	case "td", "th":
	case "tr":
		b.endPara()
	}
	b.margin(marginTop)
	b.indent += padLeft
	heading := b.heading
	if strings.HasPrefix(n.tag, "h") && len(n.tag) == 2 && n.tag[1] >= '1' && n.tag[1] <= '6' {
		b.heading = int(n.tag[1] - '0')
	}
	switch n.tag {
	case "ul", "ol":
		typ := strings.ToLower(n.attrs["type"])
		if v, ok := css["list-style-type"]; ok {
			typ = strings.ToLower(v)
		} else if v, ok := css["list-style"]; ok {
			typ = strings.ToLower(strings.Fields(v + " ")[0])
		}
		if typ == "disc" || typ == "decimal" {
			typ = ""
		}
		ls := listState{ordered: n.tag == "ol", typ: typ}
		if s, ok := n.attrs["start"]; ok {
			if v, err := strconv.Atoi(s); err == nil {
				ls.n = v - 1
			}
		}
		b.lists = append(b.lists, ls)
		b.marks = append(b.marks, paraMark{kind: markList})
		b.walk(n, st)
		b.endPara()
		b.lists = b.lists[:len(b.lists)-1]
		b.marks = append(b.marks, paraMark{kind: markEnd})
	case "li":
		b.endPara()
		b.marks = append(b.marks, paraMark{kind: markListItem})
		b.liNext = true
		b.walk(n, st)
		b.endPara()
		b.liNext = false
	default:
		b.walk(n, st)
		b.endPara()
	}
	b.heading = heading
	b.indent -= padLeft
	b.margin(marginBottom)
}

// Structure mark kinds (bdf.Mark*), repeated here for brevity.
const (
	markList     = 7
	markListItem = 8
	markEnd      = 11
)

// setPrimary resolves the style's primary font (for line metrics).
func (c *converter) setPrimary(st *tstyle) {
	fam := "Helvetica"
	if len(st.families) > 0 {
		fam = st.families[0]
	}
	st.primary = c.choice(fam, st.bold, st.italic, false)
}

// fontSizeAttr reads <font size>: 1-7, or relative +n/-n from 3.
func fontSizeAttr(v string) (float64, bool) {
	v = strings.TrimSpace(v)
	n, err := strconv.Atoi(strings.TrimPrefix(v, "+"))
	if err != nil {
		return 0, false
	}
	if strings.HasPrefix(v, "+") || strings.HasPrefix(v, "-") {
		n += 3
	}
	n = max(1, min(7, n))
	return htmlFontSizes[n-1], true
}

// applyCSS applies an inline style attribute.
func (b *textBuilder) applyCSS(css map[string]string, st, parent *tstyle, marginTop, marginBottom, padLeft *float64) {
	em := parent.size
	for k, v := range css {
		lv := strings.ToLower(v)
		switch k {
		case "color":
			if col, ok := parseColor(v); ok {
				st.color = col
			}
		case "background-color", "background":
			if col, ok := parseColor(strings.Fields(v + " ")[0]); ok && col.a > 0 {
				st.bg = &col
			}
		case "font-size":
			if sz, ok := cssLength(lv, parent.size, parent.size); ok && sz > 0 {
				st.size = sz
			} else if sz, ok := cssFontKeyword(lv, parent.size); ok {
				st.size = sz
			}
		case "font-family":
			if f := parseFontFamilies(v); len(f) > 0 {
				st.families = f
			}
		case "font-weight":
			switch lv {
			case "bold", "bolder":
				st.bold = true
			case "normal", "lighter":
				st.bold = false
			default:
				if w, err := strconv.Atoi(lv); err == nil {
					st.bold = w >= 600
				}
			}
		case "font-style":
			st.italic = lv == "italic" || lv == "oblique" || strings.HasPrefix(lv, "oblique ")
		case "text-decoration", "text-decoration-line":
			if strings.Contains(lv, "none") {
				st.underline, st.strike = false, false
			}
			if strings.Contains(lv, "underline") {
				st.underline = true
			}
			if strings.Contains(lv, "line-through") {
				st.strike = true
			}
		case "text-align":
			switch lv {
			case "left", "start":
				st.align = "left"
			case "right", "end":
				st.align = "right"
			case "center", "-webkit-center":
				st.align = "center"
			case "justify":
				st.align = "justify"
			}
		case "line-height":
			if f, err := strconv.ParseFloat(lv, 64); err == nil {
				st.lhFactor, st.lhPx = f, 0
			} else if lv == "normal" {
				st.lhFactor, st.lhPx = 1.2, 0
			} else if px, ok := cssLength(lv, st.size, st.size); ok {
				st.lhFactor, st.lhPx = 0, px
			}
		case "white-space":
			st.pre = lv == "pre" || lv == "pre-wrap" || lv == "break-spaces" || lv == "pre-line"
			st.nowrap = lv == "nowrap" || lv == "pre"
		case "margin-top":
			if px, ok := cssLength(lv, em, 0); ok {
				*marginTop = px
			}
		case "margin-bottom":
			if px, ok := cssLength(lv, em, 0); ok {
				*marginBottom = px
			}
		case "margin":
			f := strings.Fields(lv)
			if len(f) > 0 {
				if px, ok := cssLength(f[0], em, 0); ok {
					*marginTop, *marginBottom = px, px
				}
				if len(f) >= 3 {
					if px, ok := cssLength(f[2], em, 0); ok {
						*marginBottom = px
					}
				}
			}
		case "padding-left", "margin-left":
			if px, ok := cssLength(lv, em, 0); ok {
				*padLeft = px
			}
		case "vertical-align":
			switch lv {
			case "sub":
				st.shift = parent.shift + parent.size*0.2
			case "super":
				st.shift = parent.shift - parent.size*0.35
			case "baseline":
				st.shift = parent.shift
			}
		}
	}
}

// cssLength reads a CSS length in px; em is the font size ems refer to,
// pct what 100% is.
func cssLength(v string, em, pct float64) (float64, bool) {
	v = strings.TrimSpace(v)
	units := []struct {
		suffix string
		f      float64
	}{{"px", 1}, {"pt", 4.0 / 3}, {"em", em}, {"rem", 16}, {"%", pct / 100}, {"pc", 16}, {"in", 96}, {"cm", 96 / 2.54}, {"mm", 96 / 25.4}}
	for _, u := range units {
		if strings.HasSuffix(v, u.suffix) {
			if f, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(v, u.suffix)), 64); err == nil {
				return f * u.f, true
			}
			return 0, false
		}
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil && f == 0 {
		return 0, true
	}
	return 0, false
}

func cssFontKeyword(v string, parent float64) (float64, bool) {
	switch v {
	case "xx-small":
		return 9, true
	case "x-small":
		return 10, true
	case "small":
		return 13, true
	case "medium":
		return 16, true
	case "large":
		return 18, true
	case "x-large":
		return 24, true
	case "xx-large":
		return 32, true
	case "xxx-large":
		return 48, true
	case "smaller":
		return parent / 1.2, true
	case "larger":
		return parent * 1.2, true
	}
	return 0, false
}

// Line breaking follows the browser's rules closely enough for labels:
// breaks after spaces, between East Asian characters (except before
// closing punctuation and small kana, and after opening brackets:
// kinsoku), between East Asian and other text, and after hyphens inside
// words. word-wrap is normal, so long words overflow.

const noStart = "、。，．,.:;?!)]}」』】〕〉》）］｝〙〗〟”’ゝゞヽヾーぁぃぅぇぉっゃゅょゎゕゖァィゥェォッャュョヮヵヶ々〻‐゠–〜？！：；・…‥%％°℃"
const noEnd = "([{「『【〔〈《（［｛〘〖〝“‘¥$£＄￥"

func canBreak(a, b rune) bool {
	if a == ' ' || a == '　' {
		return b != ' '
	}
	if b == ' ' || a == ' ' || b == ' ' {
		return false
	}
	if strings.ContainsRune(noStart, b) || strings.ContainsRune(noEnd, a) {
		return false
	}
	if fontdb.IsCJK(a) || fontdb.IsCJK(b) {
		return true
	}
	if (a == '-' || a == '‐' || a == '/') && unicode.IsLetter(b) {
		return true
	}
	return false
}

// markBreaks sets the break opportunities of a paragraph's items.
func markBreaks(items []titem) {
	for i := range items {
		if items[i].hard {
			continue
		}
		if i+1 >= len(items) {
			items[i].brk = true
			continue
		}
		if items[i+1].hard {
			items[i].brk = true
			continue
		}
		if items[i].st.nowrap && items[i+1].st.nowrap {
			if !items[i].brk {
				items[i].brk = false
			}
			continue
		}
		items[i].brk = items[i].brk || canBreak(items[i].r, items[i+1].r)
	}
}

// layoutParagraphs lays out paragraphs in a column of the given width
// (0: no limit).
func layoutParagraphs(paras []*tpara, width float64) *tlayout {
	lo := &tlayout{paras: paras}
	y := 0.0
	for _, p := range paras {
		y += p.top
		if p.hr {
			ln := &tline{para: p, first: true, top: y, height: 2, x: p.indent}
			lo.lines = append(lo.lines, ln)
			y += 2
			continue
		}
		avail := 0.0
		if width > 0 {
			avail = math.Max(1, width-p.indent)
		}
		lines := breakLines(p, avail)
		for i, ln := range lines {
			ln.para = p
			ln.first = i == 0
			ln.x = p.indent
			ln.metrics()
			ln.top = y
			y += ln.height
			lo.lines = append(lo.lines, ln)
			if ln.width > avail && avail > 0 {
				lo.wrapped = true
			}
		}
	}
	lo.height = y
	for _, ln := range lo.lines {
		lo.width = math.Max(lo.width, ln.x+ln.width)
	}
	return lo
}

// breakLines breaks a paragraph into lines of at most width px (no limit
// for 0). Collapsible spaces at the start and end of lines are dropped.
func breakLines(p *tpara, width float64) []*tline {
	var lines []*tline
	items := p.items
	start := 0
	hardPrev, cjkPrev := false, false
	// emit ends a line before end; hard and cjk say how the next line
	// follows it (after a forced break, or a wrap between East Asian
	// characters).
	emit := func(end int, cjk bool, hard bool) {
		seg := items[start:end]
		// drop leading and trailing collapsible spaces and the break itself
		for len(seg) > 0 && seg[0].space {
			seg = seg[1:]
		}
		for len(seg) > 0 && (seg[len(seg)-1].space || seg[len(seg)-1].hard) {
			seg = seg[:len(seg)-1]
		}
		ln := &tline{items: seg, cjkWrap: cjkPrev, hardPrev: hardPrev}
		for _, it := range seg {
			ln.width += it.w
		}
		lines = append(lines, ln)
		hardPrev, cjkPrev = hard, cjk
	}
	for start < len(items) {
		x := 0.0
		lastBreak := -1
		i := start
		broke := false
		for ; i < len(items); i++ {
			it := items[i]
			if it.hard {
				emit(i+1, false, true)
				start = i + 1
				broke = true
				break
			}
			// leading spaces of a line take no room
			if !(it.space && x == 0) {
				x += it.w
			}
			if width > 0 && x > width+0.01 && !it.space && lastBreak >= start {
				cjk := fontdb.IsCJK(items[lastBreak].r) && lastBreak+1 < len(items) && fontdb.IsCJK(items[lastBreak+1].r)
				emit(lastBreak+1, cjk, false)
				start = lastBreak + 1
				broke = true
				break
			}
			if it.brk {
				lastBreak = i
			}
		}
		if !broke {
			emit(len(items), false, false)
			start = len(items)
		}
	}
	if len(lines) == 0 {
		// an empty paragraph (only a <br>) still has a line
		lines = append(lines, &tline{})
	}
	return lines
}

// metrics sets the line's height and baseline from its inline boxes and
// the strut of its block, as Blink does for an explicit line-height: each
// box's font ascent and descent (whole pixels) plus the half-leading,
// floored, above the baseline and the rest below.
func (ln *tline) metrics() {
	asc, desc := math.Inf(-1), math.Inf(-1)
	box := func(st *tstyle) {
		a, d := st.fontHeight()
		lh := layoutUnit(st.lineHeight())
		hl := math.Floor((lh - (a + d)) / 2)
		a += hl
		d = lh - a
		asc = math.Max(asc, a-st.shift)
		desc = math.Max(desc, d+st.shift)
	}
	if ln.para != nil && ln.para.st != nil {
		box(ln.para.st)
	}
	for _, it := range ln.items {
		box(it.st)
	}
	ln.asc = asc
	ln.height = asc + desc
}

// layoutUnit rounds a length down to Blink's layout precision (1/64 px).
func layoutUnit(v float64) float64 { return math.Floor(v*64) / 64 }

// fontHeight returns the ascent and descent in px of the style's primary
// font (the first family; fallback fonts do not change line boxes with an
// explicit line-height), as Blink computes them: the font's ascent and
// descent rounded to whole pixels, with the ascent of Helvetica, Times and
// Courier raised by 15% of the height as Chrome does on macOS (to match
// Arial, Times New Roman and Courier New elsewhere).
func (st *tstyle) fontHeight() (float64, float64) {
	if st.primary == nil {
		return math.Round(0.9 * st.size), math.Round(0.25 * st.size)
	}
	a, d := cssMetrics(st.primary)
	asc, desc := math.Round(a*st.size), math.Round(d*st.size)
	switch fontdb.Normalize(st.primary.use.requested) {
	case "helvetica", "times", "courier":
		if st.primary.l != nil {
			switch fontdb.Normalize(st.primary.l.Face.Family) {
			case "helvetica", "times", "courier":
				asc += math.Floor((asc+desc)*0.15 + 0.5)
			}
		}
	}
	return asc, desc
}

// cssMetrics returns the ascent and descent (em) browsers use for the
// content area of text in a face: the typographic metrics when the font
// asks for them (USE_TYPO_METRICS), else the hhea ones.
func cssMetrics(fc *faceChoice) (float64, float64) {
	if fc == nil {
		return 0.9, 0.25
	}
	if fc.l == nil {
		return fc.asc, fc.desc
	}
	f := fc.l.Font
	upem := float64(f.UnitsPerEm)
	if os2, ok := f.OS2(); ok && os2.FsSelection&0x80 != 0 && os2.TypoAscent-os2.TypoDescent > 0 {
		return float64(os2.TypoAscent) / upem, float64(-os2.TypoDescent) / upem
	}
	if f.Ascent != 0 || f.Descent != 0 {
		return float64(f.Ascent) / upem, float64(-f.Descent) / upem
	}
	return fc.asc, fc.desc
}

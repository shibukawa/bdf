package visio

import (
	"encoding/xml"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// The text of a shape is a Text element whose character data is formatted
// by the rows of the Character (cp), Paragraph (pp) and Tabs (tp) sections
// that markers among the characters select; fields (fld) hold the text
// they were last displayed with. It is laid out in the text block, which
// the TextXForm cells place in the shape. The text is described as a
// DrawingML text body and laid out by the DrawingML renderer the Office
// converters share, so that line breaking, fonts and the structure of the
// text layer are the same.

// textRun is characters with one Character row.
type textRun struct {
	text string
	char string // Character row key
	br   bool   // a line break (U+2028) within the paragraph
}

type textPara struct {
	pp   string // Paragraph row key
	runs []textRun
}

// splitText splits the content of a Text element into paragraphs.
func splitText(t *ooxml.Node) []*textPara {
	cur := &textPara{pp: "0"}
	paras := []*textPara{cur}
	cp, pp := "0", "0"
	started := false // the paragraph has characters
	addText := func(s string) {
		for s != "" {
			i := strings.IndexAny(s, "\n \r")
			chunk := s
			if i >= 0 {
				chunk = s[:i]
			}
			if chunk != "" {
				if !started {
					cur.pp = pp
					started = true
				}
				cur.runs = append(cur.runs, textRun{text: chunk, char: cp})
			}
			if i < 0 {
				return
			}
			r, size := rune(s[i]), 1
			if strings.HasPrefix(s[i:], " ") {
				r, size = ' ', len(" ")
			}
			switch r {
			case '\n':
				cur = &textPara{pp: pp}
				paras = append(paras, cur)
				started = false
			case ' ':
				if !started {
					cur.pp = pp
					started = true
				}
				cur.runs = append(cur.runs, textRun{br: true, char: cp})
			}
			s = s[i+size:]
		}
	}
	for _, seg := range t.Segments() {
		if seg.Elem == nil {
			addText(seg.Text)
			continue
		}
		switch seg.Elem.Name {
		case "cp":
			cp = seg.Elem.AttrStr("IX", "0")
		case "pp":
			pp = seg.Elem.AttrStr("IX", "0")
			if !started {
				cur.pp = pp
			}
		case "fld":
			var b strings.Builder
			for _, s := range seg.Elem.Segments() {
				if s.Elem == nil {
					b.WriteString(s.Text)
				}
			}
			addText(b.String())
		}
	}
	// the paragraph mark of the last paragraph ends it
	if len(paras) > 1 && len(paras[len(paras)-1].runs) == 0 {
		paras = paras[:len(paras)-1]
	}
	return paras
}

func elem(name string, kv ...string) *ooxml.Node {
	n := &ooxml.Node{Name: name}
	for i := 0; i+1 < len(kv); i += 2 {
		n.Attrs = append(n.Attrs, xml.Attr{Name: xml.Name{Local: kv[i]}, Value: kv[i+1]})
	}
	return n
}

func setAttr(n *ooxml.Node, k, v string) {
	for i, a := range n.Attrs {
		if a.Name.Local == k {
			n.Attrs[i].Value = v
			return
		}
	}
	n.Attrs = append(n.Attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v})
}

func add(n *ooxml.Node, kids ...*ooxml.Node) *ooxml.Node {
	n.Kids = append(n.Kids, kids...)
	return n
}

// emu converts inches to EMU for DrawingML attributes.
func emu(in float64) string { return strconv.Itoa(int(math.Round(in * 914400))) }

// centipoints writes points in hundredths, as DrawingML sizes and spacing
// are.
func centipoints(pt float64) string { return strconv.Itoa(int(math.Round(pt * 100))) }

// textBody describes the text of a shape as a DrawingML text body; it
// returns nil when the shape has no text.
func (p *pageCtx) textBody(s *shape) *ooxml.Node {
	t := s.textNode()
	if t == nil {
		return nil
	}
	paras := splitText(t)
	hasText := false
	for _, pa := range paras {
		for _, r := range pa.runs {
			if strings.TrimSpace(r.text) != "" {
				hasText = true
			}
		}
	}
	if !hasText {
		return nil
	}
	anchor := "t"
	switch int(p.num(s, "VerticalAlign", 1)) {
	case 1:
		anchor = "ctr"
	case 2:
		anchor = "b"
	}
	bodyPr := elem("bodyPr", "lIns", emu(p.num(s, "LeftMargin", 0)), "tIns", emu(p.num(s, "TopMargin", 0)),
		"rIns", emu(p.num(s, "RightMargin", 0)), "bIns", emu(p.num(s, "BottomMargin", 0)), "anchor", anchor, "wrap", "square")
	if p.num(s, "TextDirection", 0) == 1 {
		setAttr(bodyPr, "vert", "eaVert")
	}
	body := add(elem("txBody"), bodyPr, elem("lstStyle"))
	defTab := p.num(s, "DefaultTabStop", 0.5)
	for _, pa := range paras {
		body.Kids = append(body.Kids, p.paragraph(s, pa, defTab))
	}
	return body
}

// paragraph builds an a:p.
func (p *pageCtx) paragraph(s *shape, pa *textPara, defTab float64) *ooxml.Node {
	num := func(name string, def float64) float64 { return p.rowNum(s, "Paragraph", pa.pp, name, def) }
	algn := "l"
	switch int(num("HorzAlign", 1)) {
	case 1:
		algn = "ctr"
	case 2:
		algn = "r"
	case 3:
		algn = "just"
	case 4:
		algn = "dist"
	}
	indFirst, indLeft := num("IndFirst", 0), num("IndLeft", 0)
	pPr := elem("pPr", "algn", algn)
	if defTab > 0 {
		setAttr(pPr, "defTabSz", emu(defTab))
	}
	// runs first: the line height depends on their sizes
	var runs []*ooxml.Node
	var endPr *ooxml.Node
	maxSize := 0.0
	for _, r := range pa.runs {
		rPr, size := p.runProps(s, r.char)
		maxSize = math.Max(maxSize, size)
		if r.br {
			runs = append(runs, add(elem("br"), rPr))
			continue
		}
		text := r.text
		switch int(p.rowNum(s, "Character", r.char, "Case", 0)) {
		case 1:
			text = strings.ToUpper(text)
		case 2:
			text = titleCase(text)
		}
		t := elem("t")
		t.Text = text
		runs = append(runs, add(elem("r"), rPr, t))
		endPr = rPr
	}
	if endPr == nil {
		var size float64
		endPr, size = p.runProps(s, "0")
		maxSize = size
	}
	// SpLine: a height, or a multiple of the largest font of the line
	switch sp := num("SpLine", -1.2); {
	case sp < 0:
		add(pPr, add(elem("lnSpc"), elem("spcPts", "val", centipoints(-sp*maxSize))))
	case sp > 0:
		add(pPr, add(elem("lnSpc"), elem("spcPts", "val", centipoints(sp*72))))
	}
	if v := num("SpBefore", 0); v > 0 {
		add(pPr, add(elem("spcBef"), elem("spcPts", "val", centipoints(v*72))))
	}
	if v := num("SpAfter", 0); v > 0 {
		add(pPr, add(elem("spcAft"), elem("spcPts", "val", centipoints(v*72))))
	}
	marL, indent := indLeft, indFirst
	if bu := int(num("Bullet", 0)); bu > 0 && bu <= len(bulletChars) {
		ch := bulletChars[bu-1]
		if v, ok := p.rowVal(s, "Paragraph", pa.pp, "BulletStr"); ok && v != "" {
			ch = v
		}
		// the text starts TextPosAfterBullet after the bullet
		after := num("TextPosAfterBullet", 0)
		if after <= 0 {
			after = 0.25
		}
		marL, indent = indLeft+indFirst+after, -after
		add(pPr, elem("buChar", "char", ch))
	} else {
		add(pPr, elem("buNone"))
	}
	setAttr(pPr, "marL", emu(marL))
	setAttr(pPr, "indent", emu(indent))
	if tabs := p.tabs(s); len(tabs) > 0 {
		add(pPr, add(elem("tabLst"), tabs...))
	}
	para := add(elem("p"), pPr)
	para.Kids = append(para.Kids, runs...)
	end := *endPr
	end.Name = "endParaRPr"
	return add(para, &end)
}

// bulletChars are the bullets of the Bullet cell's values 1-7.
var bulletChars = []string{"•", "◆", "▪", "□", "❖", "➢", "✓"}

func titleCase(s string) string {
	prev := ' '
	return strings.Map(func(r rune) rune {
		defer func() { prev = r }()
		if unicode.IsSpace(prev) {
			return unicode.ToUpper(r)
		}
		return r
	}, s)
}

// tabs returns the tab stops of the first Tabs row.
func (p *pageCtx) tabs(s *shape) []*ooxml.Node {
	var out []*ooxml.Node
	for i := 1; i <= 160; i++ {
		v, ok := p.rowVal(s, "Tabs", "0", "Position"+strconv.Itoa(i))
		if !ok {
			break
		}
		algn := "l"
		switch int(p.rowNum(s, "Tabs", "0", "Alignment"+strconv.Itoa(i), 0)) {
		case 1:
			algn = "ctr"
		case 2:
			algn = "r"
		case 3:
			algn = "dec"
		}
		out = append(out, elem("tab", "pos", emu(num(v)), "algn", algn))
	}
	return out
}

// runProps builds the a:rPr of a Character row, and returns the font size
// in points.
func (p *pageCtx) runProps(s *shape, key string) (*ooxml.Node, float64) {
	val := func(name string) string { v, _ := p.rowVal(s, "Character", key, name); return v }
	rnum := func(name string, def float64) float64 { return p.rowNum(s, "Character", key, name, def) }
	size := rnum("Size", 12.0/72) * 72
	if size <= 0 {
		size = 12
	}
	style := int(rnum("Style", 0))
	rPr := elem("rPr", "sz", centipoints(size))
	set := func(k, v string) { rPr.Attrs = append(rPr.Attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v}) }
	if style&1 != 0 {
		set("b", "1")
	}
	if style&2 != 0 {
		set("i", "1")
	}
	switch {
	case rnum("DblUnderline", 0) != 0:
		set("u", "dbl")
	case style&4 != 0:
		set("u", "sng")
	}
	switch {
	case rnum("DoubleStrikethrough", 0) != 0:
		set("strike", "dblStrike")
	case rnum("Strikethru", 0) != 0:
		set("strike", "sngStrike")
	}
	if style&8 != 0 {
		set("cap", "small")
	}
	switch int(rnum("Pos", 0)) {
	case 1:
		set("baseline", "30000")
	case 2:
		set("baseline", "-25000")
	}
	if ls := rnum("Letterspace", 0); ls != 0 {
		set("spc", centipoints(ls*72))
	}
	if lang := langTag(val("LangID")); lang != "" {
		set("lang", lang)
	}
	c, ok := p.c.d.color(val("Color"))
	if !ok {
		c = 0x000000ff
	}
	c = withTrans(c, rnum("ColorTrans", 0))
	clr := elem("srgbClr", "val", hexOf(c)[1:])
	if a := alphaOf(c); a < 1 {
		clr.Kids = append(clr.Kids, elem("alpha", "val", strconv.Itoa(int(math.Round(a*100000)))))
	}
	rPr.Kids = append(rPr.Kids, add(elem("solidFill"), clr))
	latin := p.c.d.fontName(val("Font"))
	if latin == "" {
		if p.th != nil {
			latin = p.th.minor.latin
		} else {
			latin = "Arial"
		}
	}
	rPr.Kids = append(rPr.Kids, elem("latin", "typeface", latin))
	if ea := p.c.d.fontName(val("AsianFont")); ea != "" {
		rPr.Kids = append(rPr.Kids, elem("ea", "typeface", ea))
	} else {
		rPr.Kids = append(rPr.Kids, elem("ea", "typeface", latin))
	}
	return rPr, size
}

// langTag turns a LangID value (a tag, or a Windows locale ID in older
// documents) into a language tag.
func langTag(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if id, err := strconv.Atoi(v); err == nil {
		if t, ok := lcids[id]; ok {
			return t
		}
		return ""
	}
	return v
}

var lcids = map[int]string{
	1025: "ar-SA", 1028: "zh-TW", 1029: "cs-CZ", 1030: "da-DK", 1031: "de-DE", 1032: "el-GR", 1033: "en-US",
	1034: "es-ES", 1035: "fi-FI", 1036: "fr-FR", 1037: "he-IL", 1038: "hu-HU", 1040: "it-IT", 1041: "ja-JP",
	1042: "ko-KR", 1043: "nl-NL", 1044: "nb-NO", 1045: "pl-PL", 1046: "pt-BR", 1049: "ru-RU", 1053: "sv-SE",
	1054: "th-TH", 1055: "tr-TR", 2052: "zh-CN", 2057: "en-GB", 2070: "pt-PT", 3082: "es-ES", 3076: "zh-HK",
}

// textMatrix maps the text block of a shape (points, y down, from its top
// left corner) to the page; it returns the block's size in points. Text is
// never mirrored: a flipped shape turns its text upside down instead.
func (p *pageCtx) textMatrix(s *shape, m canvas.Matrix) (canvas.Matrix, float64, float64) {
	w, h := p.num(s, "Width", 0), p.num(s, "Height", 0)
	tw, th := p.dim(s, "TxtWidth", w), p.dim(s, "TxtHeight", h)
	px, py := p.dim(s, "TxtPinX", w/2), p.dim(s, "TxtPinY", h/2)
	lx, ly := p.dim(s, "TxtLocPinX", tw/2), p.dim(s, "TxtLocPinY", th/2)
	angle := p.num(s, "TxtAngle", 0)
	k := 72 * p.scale
	t := m.Mul(canvas.Translate(px, py)).Mul(canvas.Rotate(angle * 180 / math.Pi)).Mul(canvas.Translate(-lx, -ly)).
		Mul(canvas.Translate(0, th)).Mul(canvas.Scale(1/k, -1/k))
	bw, bh := tw*k, th*k
	if t[0]*t[3]-t[1]*t[2] < 0 {
		t = t.Mul(canvas.Translate(bw, 0)).Mul(canvas.Scale(-1, 1))
	}
	return t, math.Abs(bw), math.Abs(bh)
}

// drawText lays out and draws the text of a shape whose local coordinates
// m maps to the page.
func (p *pageCtx) drawText(cv *canvas.Canvas, s *shape, m canvas.Matrix) {
	if p.flag(s, "HideText") {
		return
	}
	body := p.textBody(s)
	if body == nil {
		return
	}
	t, bw, bh := p.textMatrix(s, m)
	if bw < 1 || fitsText(s) {
		// no room, or a block as wide as its text (whose width Visio
		// measured with its fonts): the text runs on without wrapping
		setAttr(body.Child("bodyPr"), "wrap", "none")
	}
	tb := p.dm.LayoutText(body, s.part, 0, 0, bw, bh, t)
	if tb == nil {
		return
	}
	if bg, ok := p.textBackground(s); ok {
		x0, y0, x1, y1 := tb.Bounds()
		cv.Obj.Save()
		cv.Transform(t)
		cv.Obj.FillColor(bg)
		cv.Obj.FillRect(f32(x0), f32(y0), f32(x1-x0), f32(y1-y0))
		cv.Obj.Restore()
	}
	tb.Draw(cv)
	cv.Drawn = true
}

// fitsText reports whether a shape's text block is as wide as its text
// (TxtWidth is TEXTWIDTH(TheText), or at least that), rather than a width
// to wrap the text in.
func fitsText(s *shape) bool {
	f := strings.ToUpper(s.formula("TxtWidth"))
	return strings.Contains(f, "TEXTWIDTH(") && !strings.Contains(f, "MIN(")
}

// textBackground returns the color behind a shape's text: a color, or in
// older documents 0 for none and a color index plus one.
func (p *pageCtx) textBackground(s *shape) (bdf.Color, bool) {
	v, ok := p.val(s, "TextBkgnd")
	if !ok {
		return 0, false
	}
	v = strings.TrimSpace(v)
	if !strings.HasPrefix(v, "#") {
		i := atoi(v, 0)
		if i <= 0 {
			return 0, false
		}
		v = strconv.Itoa(i - 1)
	}
	c, ok := p.c.d.color(v)
	if !ok {
		return 0, false
	}
	c = withTrans(c, p.num(s, "TextBkgndTrans", 0))
	return c, c&0xff != 0
}

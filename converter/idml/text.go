package idml

import (
	"encoding/xml"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// A story (Stories/Story_*.xml) is the text of a chain of text frames:
// ParagraphStyleRange elements holding CharacterStyleRange elements, whose
// Content elements hold the characters and whose Br elements end
// paragraphs. Properties come from the range's own attributes, then its
// applied style and the styles that one is based on, then (for character
// properties) the paragraph's. The text is described as a DrawingML text
// body and laid out by the DrawingML renderer the Office converters
// share, flowing through the frames of the chain (drawingml.TextFlow), so
// that line breaking, fonts and the structure of the text layer are the
// same as for PowerPoint and Visio.

// Characters with a meaning of their own in InDesign text.
const (
	charPageNumber   = '\u0018' // the current page number
	charSectionName  = '\u0019' // the section marker
	charIndentToHere = '\u0007'
	charRightTab     = '\u0008'
	charEndNested    = '\u0003' // end of nested style
	charSoftHyphen   = '­'
	charLineBreak    = ' '      // forced line break
	charColumnBreak  = '\u000e' // column, frame, page breaks and the like
)

// textStyle is the resolved character properties of a run.
type textStyle struct {
	size      float64
	font      string
	bold      bool
	italic    bool
	underline bool
	strike    bool
	caps      string
	position  string
	shift     float64 // baseline shift in points
	tracking  float64 // 1/1000 em
	lang      string
	fill      string
	fillTint  float64
	leading   float64 // points; 0 for auto
	autoLead  float64 // percent
	hasText   bool
}

// resolveText reads the character properties along a chain of elements
// (ranges and styles, nearest first).
func (c *converter) resolveText(chain []*ooxml.Node) textStyle {
	get := func(name, def string) string {
		if v, ok := lookup(chain, name); ok {
			return v
		}
		return def
	}
	num := func(name string, def float64) float64 {
		if v, ok := lookup(chain, name); ok {
			if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
				return f
			}
		}
		return def
	}
	st := textStyle{size: num("PointSize", 12), font: get("AppliedFont", "Minion Pro"), caps: get("Capitalization", "Normal"),
		position: get("Position", "Normal"), shift: num("BaselineShift", 0), tracking: num("Tracking", 0),
		fill: get("FillColor", "Color/Black"), fillTint: num("FillTint", -1), autoLead: num("AutoLeading", 120)}
	if st.size <= 0 {
		st.size = 12
	}
	style := strings.ToLower(get("FontStyle", "Regular"))
	for _, w := range []string{"bold", "semibold", "black", "heavy", "extrabold", "ultra"} {
		if strings.Contains(style, w) {
			st.bold = true
		}
	}
	st.italic = strings.Contains(style, "italic") || strings.Contains(style, "oblique")
	st.underline = get("Underline", "false") == "true"
	st.strike = get("StrikeThru", "false") == "true"
	if v, ok := lookup(chain, "Leading"); ok {
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil && f > 0 {
			st.leading = f
		}
	}
	st.lang = langTag(get("AppliedLanguage", ""))
	return st
}

// langTag turns an AppliedLanguage value ("$ID/Japanese", "$ID/English:
// USA") into a language tag.
func langTag(v string) string {
	v = strings.TrimPrefix(strings.TrimSpace(v), "$ID/")
	if u, err := url.PathUnescape(v); err == nil {
		v = u
	}
	if v == "" || v == "[No Language]" {
		return ""
	}
	if t, ok := languages[v]; ok {
		return t
	}
	base, _, _ := strings.Cut(v, ":")
	if t, ok := languages[strings.TrimSpace(base)]; ok {
		return t
	}
	return ""
}

var languages = map[string]string{
	"Japanese": "ja-JP", "English": "en", "English: USA": "en-US", "English: UK": "en-GB", "English: Canadian": "en-CA",
	"German": "de", "German: Standard": "de-DE", "German: Swiss": "de-CH", "German: Austrian": "de-AT",
	"French": "fr", "French: Standard": "fr-FR", "French: Canadian": "fr-CA", "Spanish": "es", "Spanish: Castilian": "es-ES",
	"Italian": "it", "Portuguese": "pt-PT", "Portuguese: Brazilian": "pt-BR", "Dutch": "nl", "Danish": "da", "Swedish": "sv",
	"Norwegian": "no", "Finnish": "fi", "Russian": "ru", "Polish": "pl", "Czech": "cs", "Hungarian": "hu", "Turkish": "tr",
	"Greek": "el", "Korean": "ko", "Simplified Chinese": "zh-CN", "Traditional Chinese": "zh-TW", "Chinese": "zh",
	"Arabic": "ar", "Hebrew": "he", "Thai": "th", "Vietnamese": "vi", "Hindi": "hi",
}

// storyText is a story described as a DrawingML text body.
type storyText struct {
	body     *ooxml.Node
	vertical bool
	hasText  bool
}

// storyBody builds the text body of a story; pageNum is the number of the
// page the story's first frame is on, for page number markers.
func (c *converter) storyBody(story *ooxml.Node, pageNum string) *storyText {
	out := &storyText{body: elem("txBody")}
	out.body.Kids = append(out.body.Kids, elem("bodyPr"), elem("lstStyle"))
	out.vertical = story.Child("StoryPreference").AttrStr("StoryOrientation", "Horizontal") == "Vertical"
	b := &storyBuilder{c: c, st: out, pageNum: pageNum}
	b.walk(story, nil, nil)
	b.endParagraph(true)
	return out
}

// storyBuilder turns the ranges of a story into paragraphs.
type storyBuilder struct {
	c         *converter
	st        *storyText
	pageNum   string
	para      *ooxml.Node   // the paragraph being built (nil between paragraphs)
	pchain    []*ooxml.Node // the paragraph's style chain
	maxSize   float64       // the largest run of the paragraph
	runs      int
	lastStyle textStyle // the style of the last run
	depth     int
}

const maxStoryDepth = 64

// walk reads the elements of a story. psr and csr are the paragraph and
// character style ranges in effect.
func (b *storyBuilder) walk(n *ooxml.Node, psr, csr *ooxml.Node) {
	if b.depth > maxStoryDepth {
		return
	}
	b.depth++
	defer func() { b.depth-- }()
	for _, k := range n.Elements() {
		switch k.Name {
		case "ParagraphStyleRange":
			b.walk(k, k, csr)
		case "CharacterStyleRange":
			b.walk(k, psr, k)
		case "Content":
			b.content(k.Content(), psr, csr)
		case "Br":
			b.ensureParagraph(psr)
			b.endParagraph(false)
		case "HyperlinkTextSource", "XMLElement", "Change", "HiddenText", "TextVariableInstance", "Cell", "Row":
			// wrappers around ranges (hidden text is drawn: conditions are not read)
			b.walk(k, psr, csr)
		case "Table":
			b.c.warnOnce("table", "tables in text are not drawn")
		case "Footnote":
			b.c.warnOnce("footnote", "footnotes are not drawn")
		case "Note":
			// an editorial note: not part of the text
		case "Rectangle", "Oval", "Polygon", "GraphicLine", "TextFrame", "Group":
			b.c.warnOnce("anchored", "objects anchored in text are not drawn")
		}
	}
}

// ensureParagraph starts a paragraph when none is open.
func (b *storyBuilder) ensureParagraph(psr *ooxml.Node) {
	if b.para != nil {
		return
	}
	b.pchain = append([]*ooxml.Node{psr}, styleChain(psr, "AppliedParagraphStyle", b.c.d.pstyles)...)
	if psr == nil {
		b.pchain = b.pchain[1:]
	}
	b.para = elem("p")
	b.maxSize, b.runs = 0, 0
}

// endParagraph closes the paragraph being built; last says that the story
// ends here (an empty trailing paragraph is kept only when the story has
// nothing else).
func (b *storyBuilder) endParagraph(last bool) {
	if b.para == nil {
		if last && len(b.st.body.Children("p")) == 0 {
			b.ensureParagraph(nil)
		} else {
			return
		}
	}
	if last && b.runs == 0 && len(b.st.body.Children("p")) > 0 {
		b.para = nil
		return
	}
	pPr, end := b.paragraphProps()
	b.para.Kids = append([]*ooxml.Node{pPr}, b.para.Kids...)
	b.para.Kids = append(b.para.Kids, end)
	b.st.body.Kids = append(b.st.body.Kids, b.para)
	b.para = nil
}

// content adds the characters of a Content element to the paragraph.
func (b *storyBuilder) content(text string, psr, csr *ooxml.Node) {
	if text == "" {
		return
	}
	b.ensureParagraph(psr)
	chain := append([]*ooxml.Node{csr}, styleChain(csr, "AppliedCharacterStyle", b.c.d.cstyles)...)
	if csr == nil {
		chain = chain[1:]
	}
	chain = append(chain, b.pchain...)
	st := b.c.resolveText(chain)
	rPr := b.runProps(st)
	var sb strings.Builder
	flush := func() {
		if sb.Len() == 0 {
			return
		}
		t := elem("t")
		t.Text = sb.String()
		b.para.Kids = append(b.para.Kids, add(elem("r"), rPr, t))
		sb.Reset()
		b.runs++
	}
	for _, r := range text {
		switch {
		case r == charLineBreak:
			flush()
			b.para.Kids = append(b.para.Kids, add(elem("br"), rPr))
		case r == charPageNumber:
			sb.WriteString(b.pageNum)
		case r == charSectionName, r == charSoftHyphen, r == charIndentToHere, r == charEndNested, r == 0x200b, r == 0xfeff:
		case r == charRightTab:
			sb.WriteByte('\t')
		case r == charColumnBreak:
			b.c.warnOnce("break", "column, frame and page breaks in text are ignored")
		case r < ' ' && r != '\t':
		default:
			sb.WriteRune(r)
		}
	}
	flush()
	b.maxSize = math.Max(b.maxSize, st.size)
	b.lastStyle = st
}

// runProps builds the a:rPr of a run.
func (b *storyBuilder) runProps(st textStyle) *ooxml.Node {
	rPr := elem("rPr", "sz", centipoints(st.size))
	set := func(k, v string) { rPr.Attrs = append(rPr.Attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v}) }
	if st.bold {
		set("b", "1")
	}
	if st.italic {
		set("i", "1")
	}
	if st.underline {
		set("u", "sng")
	}
	if st.strike {
		set("strike", "sngStrike")
	}
	switch st.caps {
	case "AllCaps":
		set("cap", "all")
	case "SmallCaps", "CapToSmallCap":
		set("cap", "small")
	}
	switch st.position {
	case "Superscript", "OTSuperscript", "OTNumerator":
		set("baseline", "30000")
	case "Subscript", "OTSubscript", "OTDenominator":
		set("baseline", "-25000")
	default:
		if st.shift != 0 && st.size > 0 {
			set("baseline", strconv.Itoa(int(math.Round(st.shift/st.size*100000))))
		}
	}
	if st.tracking != 0 {
		set("spc", centipoints(st.tracking/1000*st.size))
	}
	if st.lang != "" {
		set("lang", st.lang)
	}
	if col, ok := b.c.swatch(st.fill); ok {
		if st.fillTint >= 0 {
			col = col.tint(st.fillTint)
		}
		rPr.Kids = append(rPr.Kids, add(elem("solidFill"), colorNode(col)))
	} else if g := b.c.gradient(st.fill); g != nil {
		rPr.Kids = append(rPr.Kids, add(elem("solidFill"), colorNode(b.c.average(paintOf{grad: g, tint: st.fillTint, alpha: 1}))))
	} else {
		rPr.Kids = append(rPr.Kids, elem("noFill"))
	}
	rPr.Kids = append(rPr.Kids, elem("latin", "typeface", st.font), elem("ea", "typeface", st.font))
	return rPr
}

// colorNode makes an a:srgbClr with the color's alpha.
func colorNode(col rgba) *ooxml.Node {
	c := col.bdf()
	clr := elem("srgbClr", "val", strings.ToUpper(strconv.FormatUint(uint64(c>>8), 16)))
	for len(clr.Attrs[0].Value) < 6 {
		clr.Attrs[0].Value = "0" + clr.Attrs[0].Value
	}
	if a := uint8(c); a < 255 {
		clr.Kids = append(clr.Kids, elem("alpha", "val", strconv.Itoa(int(math.Round(float64(a)/255*100000)))))
	}
	return clr
}

// paragraphProps builds the a:pPr and a:endParaRPr of the paragraph being
// built from its style chain.
func (b *storyBuilder) paragraphProps() (pPr, end *ooxml.Node) {
	chain := b.pchain
	get := func(name, def string) string {
		if v, ok := lookup(chain, name); ok {
			return v
		}
		return def
	}
	num := func(name string, def float64) float64 {
		if v, ok := lookup(chain, name); ok {
			if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
				return f
			}
		}
		return def
	}
	algn := "l"
	switch get("Justification", "LeftAlign") {
	case "RightAlign":
		algn = "r"
	case "CenterAlign":
		algn = "ctr"
	case "LeftJustified", "RightJustified", "CenterJustified":
		algn = "just"
	case "FullyJustified":
		algn = "dist"
	}
	pPr = elem("pPr", "algn", algn, "marL", emu(num("LeftIndent", 0)), "indent", emu(num("FirstLineIndent", 0)))
	st := b.c.resolveText(chain)
	if b.runs > 0 {
		st = b.lastStyle
	}
	size := math.Max(b.maxSize, st.size)
	// leading: a height, or automatic as a percentage of the largest size
	leading := st.leading
	if leading <= 0 {
		leading = st.autoLead / 100 * size
	}
	add(pPr, add(elem("lnSpc"), elem("spcPts", "val", centipoints(leading))))
	if v := num("SpaceBefore", 0); v > 0 {
		add(pPr, add(elem("spcBef"), elem("spcPts", "val", centipoints(v))))
	}
	if v := num("SpaceAfter", 0); v > 0 {
		add(pPr, add(elem("spcAft"), elem("spcPts", "val", centipoints(v))))
	}
	switch get("BulletsAndNumberingListType", "NoList") {
	case "BulletList":
		ch := "•"
		for _, n := range chain {
			if bc := n.Path("Properties", "BulletChar"); bc != nil {
				if v := bc.AttrInt("BulletCharacterValue", 0); v > 0 && v < 0x110000 {
					ch = string(rune(v))
				}
				break
			}
		}
		add(pPr, elem("buChar", "char", ch))
	case "NumberedList":
		scheme := "arabicPeriod"
		switch f := get("NumberingFormat", "^#."); {
		case strings.Contains(f, "A"):
			scheme = "alphaUcPeriod"
		case strings.Contains(f, "a"):
			scheme = "alphaLcPeriod"
		case strings.Contains(f, "I"):
			scheme = "romanUcPeriod"
		case strings.Contains(f, "i"):
			scheme = "romanLcPeriod"
		case strings.HasSuffix(f, ")"):
			scheme = "arabicParenR"
		}
		add(pPr, elem("buAutoNum", "type", scheme))
	default:
		add(pPr, elem("buNone"))
	}
	var tabs []*ooxml.Node
	for _, n := range chain {
		tl := n.Path("Properties", "TabList")
		if tl == nil {
			continue
		}
		for _, it := range tl.Children("ListItem") {
			pos, err := strconv.ParseFloat(strings.TrimSpace(it.Child("Position").Content()), 64)
			if err != nil {
				continue
			}
			algn := "l"
			switch strings.TrimSpace(it.Child("Alignment").Content()) {
			case "CenterAlign":
				algn = "ctr"
			case "RightAlign":
				algn = "r"
			case "CharacterAlign":
				algn = "dec"
			}
			tabs = append(tabs, elem("tab", "pos", emu(pos), "algn", algn))
		}
		break
	}
	if len(tabs) > 0 {
		add(pPr, add(elem("tabLst"), tabs...))
	}
	end = b.runProps(st)
	end.Name = "endParaRPr"
	return pPr, end
}

func elem(name string, kv ...string) *ooxml.Node {
	n := &ooxml.Node{Name: name}
	for i := 0; i+1 < len(kv); i += 2 {
		n.Attrs = append(n.Attrs, xml.Attr{Name: xml.Name{Local: kv[i]}, Value: kv[i+1]})
	}
	return n
}

func add(n *ooxml.Node, kids ...*ooxml.Node) *ooxml.Node {
	n.Kids = append(n.Kids, kids...)
	return n
}

// emu formats points as EMUs, the unit of DrawingML lengths.
func emu(pt float64) string { return strconv.Itoa(int(math.Round(pt * 12700))) }

// centipoints formats points as hundredths, the unit of DrawingML sizes.
func centipoints(pt float64) string { return strconv.Itoa(int(math.Round(pt * 100))) }

// frameBodyPr builds the a:bodyPr of a text frame from its
// TextFramePreference: insets, columns and vertical alignment, and the
// story's orientation.
func frameBodyPr(pref *ooxml.Node, vertical bool) *ooxml.Node {
	ins := [4]float64{}
	if v := listValues(pref, "InsetSpacing"); len(v) == 4 {
		for i := range ins {
			ins[i], _ = strconv.ParseFloat(v[i], 64)
		}
	} else if len(v) == 1 {
		f, _ := strconv.ParseFloat(v[0], 64)
		ins = [4]float64{f, f, f, f}
	}
	anchor := "t"
	switch pref.AttrStr("VerticalJustification", "TopAlign") {
	case "CenterAlign":
		anchor = "ctr"
	case "BottomAlign":
		anchor = "b"
	}
	// InDesign's insets are left top right bottom (right to left for a
	// vertical frame, where the text starts at the right)
	bp := elem("bodyPr", "lIns", emu(ins[0]), "tIns", emu(ins[1]), "rIns", emu(ins[2]), "bIns", emu(ins[3]), "anchor", anchor, "wrap", "square")
	if cols := pref.AttrInt("TextColumnCount", 1); cols > 1 {
		bp.Attrs = append(bp.Attrs, xml.Attr{Name: xml.Name{Local: "numCol"}, Value: strconv.Itoa(int(min(cols, 40)))},
			xml.Attr{Name: xml.Name{Local: "spcCol"}, Value: emu(pref.AttrFloat("TextColumnGutter", 12))})
	}
	if vertical {
		bp.Attrs = append(bp.Attrs, xml.Attr{Name: xml.Name{Local: "vert"}, Value: "eaVert"})
	}
	return bp
}

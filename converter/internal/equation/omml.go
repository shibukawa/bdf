package equation

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// OMML is what reading Office Math needs to know about its document.
type OMML struct {
	// NaryLim and IntLim are where the limits of n-ary operators and of
	// integrals go when an operator does not say: "undOvr" or "subSup"
	// (m:mathPr's naryLim and intLim; "" for Word's defaults).
	NaryLim, IntLim string
	// RunStyle reads the formatting of the runs of Office Math in DrawingML
	// text (PowerPoint, Excel), whose a:rPr inherits from its text body:
	// the style to set a run in (nil for none), and its font and weight
	// for normal text. The runs of Word documents have a w:rPr, which the
	// reader reads itself.
	RunStyle func(rPr *ooxml.Node) (style *Styled, font string, bold, italic bool)
}

// alignMark is an alignment point of an equation array (& in a run).
type alignMark struct{}

func (*alignMark) isNode() {}

// ParseOMMLPara reads a display math paragraph (m:oMathPara): the formula
// on each of its lines, and their justification ("left", "right" or
// "center").
func (o *OMML) ParseOMMLPara(n *ooxml.Node) (lines []Node, jc string) {
	jc = "center"
	switch n.Path("oMathParaPr", "jc").AttrStr("val", "") {
	case "left":
		jc = "left"
	case "right":
		jc = "right"
	}
	for _, k := range n.Elements() {
		if k.Name == "oMath" {
			lines = append(lines, o.ParseOMML(k))
		}
	}
	return lines, jc
}

// ParseOMML reads an Office Math zone (m:oMath) or any of its arguments.
func (o *OMML) ParseOMML(n *ooxml.Node) Node {
	return row(o.content(n))
}

// content reads the elements of an argument.
func (o *OMML) content(n *ooxml.Node) []Node {
	var out []Node
	for _, k := range n.Elements() {
		switch k.Name {
		case "r":
			out = append(out, o.run(k)...)
		case "ins", "smartTag", "customXml", "sdt", "sdtContent", "oMath", "e", "fldSimple", "hyperlink":
			out = append(out, o.content(k)...)
		default:
			if node := o.element(k); node != nil {
				out = append(out, node)
			}
		}
	}
	return out
}

// arg reads an argument element; nil when it is missing or empty.
func (o *OMML) arg(n *ooxml.Node) Node {
	if n == nil {
		return nil
	}
	nodes := o.content(n)
	if len(nodes) == 0 {
		return nil
	}
	return row(nodes)
}

// argOrEmpty is arg, with an empty row for a missing argument.
func (o *OMML) argOrEmpty(n *ooxml.Node) Node {
	if a := o.arg(n); a != nil {
		return a
	}
	return &Row{}
}

// on reads an on/off property: present and not switched off.
func on(n *ooxml.Node) bool {
	return n != nil && n.AttrBool("val", true)
}

// chr reads a character property; def when it is missing, "" when empty.
func chr(pr *ooxml.Node, name, def string) string {
	c := pr.Child(name)
	if c == nil {
		return def
	}
	v, ok := c.Attr("val")
	if !ok {
		return def
	}
	return v
}

func (o *OMML) element(n *ooxml.Node) Node {
	switch n.Name {
	case "acc":
		pr := n.Child("accPr")
		a := NewOp(chr(pr, "chr", "̂"))
		a.Accent, a.Stretchy = true, true
		return &UnderOver{Base: o.argOrEmpty(n.Child("e")), Over: a, AccentOver: true}
	case "bar":
		return &Bar{Kid: o.argOrEmpty(n.Child("e")), Under: n.Path("barPr", "pos").AttrStr("val", "bot") != "top"}
	case "borderBox":
		pr := n.Child("borderBoxPr")
		return &Enclose{Kid: o.argOrEmpty(n.Child("e")),
			Top: !on(pr.Child("hideTop")), Bottom: !on(pr.Child("hideBot")),
			Left: !on(pr.Child("hideLeft")), Right: !on(pr.Child("hideRight")),
			StrikeH: on(pr.Child("strikeH")), StrikeV: on(pr.Child("strikeV")),
			StrikeUp: on(pr.Child("strikeBLTR")), StrikeDown: on(pr.Child("strikeTLBR"))}
	case "box", "groupChr", "limLow", "limUpp", "phant":
		return o.boxLike(n)
	case "d":
		return o.delimiters(n)
	case "eqArr":
		return o.eqArr(n)
	case "f":
		f := &Frac{Num: o.argOrEmpty(n.Child("num")), Den: o.argOrEmpty(n.Child("den"))}
		switch n.Path("fPr", "type").AttrStr("val", "bar") {
		case "skw":
			f.Kind = FracSkewed
		case "lin":
			f.Kind = FracLinear
		case "noBar":
			f.NoBar = true
		}
		return f
	case "func":
		name := o.argOrEmpty(n.Child("fName"))
		return &Row{Kids: []Node{&Row{Kids: []Node{name}, Class: LargeOp}, o.argOrEmpty(n.Child("e"))}}
	case "m":
		return o.matrix(n)
	case "nary":
		return o.nary(n)
	case "rad":
		r := &Radical{Base: o.argOrEmpty(n.Child("e"))}
		if !on(n.Path("radPr", "degHide")) {
			r.Degree = o.arg(n.Child("deg"))
		}
		return r
	case "sPre":
		return &Scripts{Base: o.argOrEmpty(n.Child("e")), PreSub: o.arg(n.Child("sub")), PreSup: o.arg(n.Child("sup"))}
	case "sSub", "sSup", "sSubSup":
		return &Scripts{Base: o.argOrEmpty(n.Child("e")), Sub: o.arg(n.Child("sub")), Sup: o.arg(n.Child("sup"))}
	}
	return nil
}

func (o *OMML) boxLike(n *ooxml.Node) Node {
	e := o.argOrEmpty(n.Child("e"))
	switch n.Name {
	case "groupChr":
		pr := n.Child("groupChrPr")
		a := NewOp(chr(pr, "chr", "⏟"))
		a.Stretchy, a.Accent = true, true
		if pr.Path("pos").AttrStr("val", "bot") == "top" {
			return &UnderOver{Base: e, Over: a, AccentOver: true}
		}
		return &UnderOver{Base: e, Under: a, AccentUnder: true}
	case "limLow":
		return &UnderOver{Base: e, Under: o.arg(n.Child("lim"))}
	case "limUpp":
		return &UnderOver{Base: e, Over: o.arg(n.Child("lim"))}
	case "phant":
		pr := n.Child("phantPr")
		show := true
		if s := pr.Child("show"); s != nil {
			show = on(s)
		}
		return &Phantom{Kid: e, Show: show && !on(pr.Child("transp")), ZeroWidth: on(pr.Child("zeroWid")),
			ZeroAsc: on(pr.Child("zeroAsc")), ZeroDesc: on(pr.Child("zeroDesc"))}
	}
	if on(n.Path("boxPr", "opEmu")) {
		return &Row{Kids: []Node{e}, Class: Rel}
	}
	return e
}

// delimiters reads m:d: arguments between growing delimiters, separated by
// a separator.
func (o *OMML) delimiters(n *ooxml.Node) Node {
	pr := n.Child("dPr")
	beg, sep, end := chr(pr, "begChr", "("), chr(pr, "sepChr", "|"), chr(pr, "endChr", ")")
	grow := true
	if g := pr.Child("grow"); g != nil {
		grow = on(g)
	}
	symmetric := pr.Path("shp").AttrStr("val", "centered") != "match"
	fence := func(s string, class Class) *Atom {
		a := NewOp(s)
		a.Stretchy, a.Symmetric, a.Class = grow, symmetric, class
		return a
	}
	// spaced as an ordinary atom: Word sets f(x) without the thin space
	// TeX puts before \left
	r := &Row{Class: Ord}
	if beg != "" {
		r.Kids = append(r.Kids, fence(beg, Open))
	}
	for i, e := range n.Children("e") {
		if i > 0 && sep != "" {
			class := Ord
			if sep == "," || sep == ";" {
				class = Punct
			}
			r.Kids = append(r.Kids, fence(sep, class))
		}
		r.Kids = append(r.Kids, o.argOrEmpty(e))
	}
	if end != "" {
		r.Kids = append(r.Kids, fence(end, Close))
	}
	return r
}

// eqArr reads an equation array: rows aligned at their & marks.
func (o *OMML) eqArr(n *ooxml.Node) Node {
	t := &Table{Display: true, RowGap: 0.2}
	for _, e := range n.Children("e") {
		cells := [][]Node{nil}
		for _, k := range flatten(o.content(e)) {
			if _, ok := k.(*alignMark); ok {
				cells = append(cells, nil)
				t.Aligned = true
				continue
			}
			cells[len(cells)-1] = append(cells[len(cells)-1], k)
		}
		var rowNodes []Node
		for _, c := range cells {
			rowNodes = append(rowNodes, row(c))
		}
		t.Rows = append(t.Rows, rowNodes)
	}
	if !t.Aligned {
		t.ColAlign = []Align{Center}
	}
	return t
}

// flatten spreads the rows a list holds into it, so that alignment marks in
// them are seen.
func flatten(nodes []Node) []Node {
	var out []Node
	for _, n := range nodes {
		if r, ok := n.(*Row); ok && r.Class == ClassAuto {
			out = append(out, flatten(r.Kids)...)
			continue
		}
		out = append(out, n)
	}
	return out
}

func (o *OMML) matrix(n *ooxml.Node) Node {
	t := &Table{}
	for _, mc := range n.Path("mPr", "mcs").Children("mc") {
		count := int(mc.Path("mcPr", "count").AttrInt("val", 1))
		a := Center
		switch mc.Path("mcPr", "mcJc").AttrStr("val", "center") {
		case "left":
			a = Left
		case "right":
			a = Right
		}
		for i := 0; i < count && i < 256; i++ {
			t.ColAlign = append(t.ColAlign, a)
		}
	}
	for _, mr := range n.Children("mr") {
		var cells []Node
		for _, e := range mr.Children("e") {
			cells = append(cells, o.argOrEmpty(e))
		}
		t.Rows = append(t.Rows, cells)
	}
	return t
}

func (o *OMML) nary(n *ooxml.Node) Node {
	pr := n.Child("naryPr")
	c := chr(pr, "chr", "∫")
	op := NewOp(c)
	op.LargeOp, op.Class = true, LargeOp
	integral := strings.ContainsRune(integrals, []rune(c + " ")[0])
	loc := pr.Path("limLoc").AttrStr("val", "")
	if loc == "" {
		loc = o.NaryLim
		if integral {
			loc = o.IntLim
		}
	}
	if loc == "" {
		loc = "undOvr"
		if integral {
			loc = "subSup"
		}
	}
	var sub, sup Node
	if !on(pr.Child("subHide")) {
		sub = o.arg(n.Child("sub"))
	}
	if !on(pr.Child("supHide")) {
		sup = o.arg(n.Child("sup"))
	}
	var base Node = op
	switch {
	case sub == nil && sup == nil:
	case loc == "undOvr":
		op.MovableLimits = true
		base = &UnderOver{Base: op, Under: sub, Over: sup}
	default:
		base = &Scripts{Base: op, Sub: sub, Sup: sup}
	}
	return &Row{Kids: []Node{base, o.argOrEmpty(n.Child("e"))}}
}

// run reads the text of a math run into tokens.
func (o *OMML) run(n *ooxml.Node) []Node {
	var text strings.Builder
	for _, k := range n.Elements() {
		switch k.Name {
		case "t":
			text.WriteString(k.Content())
		case "tab":
			text.WriteString(" ")
		}
	}
	s := text.String()
	if s == "" {
		return nil
	}
	mrp, wrp := n.Child("rPr"), (*ooxml.Node)(nil)
	for _, k := range n.Elements() {
		if k.Name == "rPr" && k != mrp {
			wrp = k
		}
	}
	// m:rPr and w:rPr (a:rPr) have the same local name: tell them apart
	// by their namespace
	if mrp != nil && isWordRPr(mrp) {
		mrp, wrp = wrp, mrp
	}
	drawingML := wrp != nil && strings.Contains(wrp.Space, "drawingml")
	variant := VarAuto
	sty := mrp.Path("sty").AttrStr("val", "")
	plain := sty == "p" || sty == "b"
	switch sty {
	case "p":
		variant = VarNormal
	case "b":
		variant = VarBold
	case "i":
		variant = VarItalic
	case "bi":
		variant = VarBoldItalic
	}
	switch mrp.Path("scr").AttrStr("val", "roman") {
	case "script":
		variant = pick(variant.isBold(), VarBoldScript, VarScript)
	case "fraktur":
		variant = pick(variant.isBold(), VarBoldFraktur, VarFraktur)
	case "double-struck":
		variant = VarDoubleStruck
	case "sans-serif":
		switch sty {
		case "p":
			variant = VarSansSerif
		case "b":
			variant = VarSansSerifBold
		case "bi":
			variant = VarSansSerifBoldItalic
		default:
			variant = VarSansSerifItalic
		}
	case "monospace":
		variant = VarMonospace
	}
	var style *Styled
	var dmFont string
	var dmBold, dmItalic bool
	switch {
	case drawingML && o.RunStyle != nil:
		style, dmFont, dmBold, dmItalic = o.RunStyle(wrp)
	case !drawingML:
		style = runStyle(wrp)
	}
	wrap := func(n Node) Node {
		if style == nil {
			return n
		}
		s := *style
		s.Kid = n
		return &s
	}
	if on(mrp.Child("nor")) {
		font, bold, italic := wrp.Path("rFonts").AttrStr("ascii", ""), on(wrp.Child("b")), on(wrp.Child("i"))
		if drawingML {
			font, bold, italic = dmFont, dmBold, dmItalic
		}
		if strings.EqualFold(font, "Cambria Math") {
			font = ""
		}
		v := VarNormal
		if bold {
			v = VarBold
		}
		if italic {
			v = pick(v == VarBold, VarBoldItalic, VarItalic)
		}
		return []Node{wrap(&Atom{Kind: Text, Text: s, Font: font, Variant: v})}
	}
	var out []Node
	rs := []rune(s)
	// PowerPoint stores the letters of formulas styled (𝑥 for an italic
	// x): back to the letter and its style
	var styles []Variant
	for i, r := range rs {
		if p, v, ok := plainLetter(r); ok {
			if styles == nil {
				styles = make([]Variant, len(rs))
			}
			rs[i], styles[i] = p, v
		}
	}
	styleAt := func(i int) Variant {
		if styles == nil || styles[i] == VarAuto {
			return variant
		}
		if styles[i] == VarItalic && variant == VarAuto && !upperGreek(rs[i]) {
			// the default style of a letter
			return VarAuto
		}
		return styles[i]
	}
	for i := 0; i < len(rs); {
		r := rs[i]
		switch {
		case r == '&':
			out = append(out, &alignMark{})
			i++
		case r == ' ' || r == ' ':
			out = append(out, wrap(&Space{Width: 0.2}))
			i++
		case unicode.IsDigit(r):
			j := i + 1
			for j < len(rs) && (unicode.IsDigit(rs[j]) || rs[j] == '.' && j+1 < len(rs) && unicode.IsDigit(rs[j+1])) {
				j++
			}
			out = append(out, wrap(&Atom{Kind: Number, Text: string(rs[i:j]), Variant: numberVariant(styleAt(i))}))
			i = j
		case isLetter(r) || unicode.Is(unicode.Han, r) || unicode.In(r, unicode.Hiragana, unicode.Katakana):
			j := i + 1
			if plain {
				for j < len(rs) && isLetter(rs[j]) && styleAt(j) == styleAt(i) {
					j++
				}
			}
			if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana) {
				for j < len(rs) && unicode.In(rs[j], unicode.Han, unicode.Hiragana, unicode.Katakana) {
					j++
				}
				out = append(out, wrap(&Atom{Kind: Text, Text: string(rs[i:j])}))
				i = j
				continue
			}
			out = append(out, wrap(&Atom{Kind: Ident, Text: string(rs[i:j]), Variant: styleAt(i)}))
			i = j
		default:
			a := NewOp(string(r))
			a.Stretchy = false
			if variant.isBold() {
				out = append(out, wrap(&Styled{Kid: a, Bold: true}))
			} else {
				out = append(out, wrap(a))
			}
			i++
		}
	}
	return out
}

// numberVariant is the variant digits take in a run of variant v.
func numberVariant(v Variant) Variant {
	switch v {
	case VarBold, VarBoldItalic:
		return VarBold
	case VarDoubleStruck, VarSansSerif, VarSansSerifBold, VarMonospace:
		return v
	}
	return VarAuto
}

func pick(c bool, a, b Variant) Variant {
	if c {
		return a
	}
	return b
}

// isWordRPr reports whether an rPr is WordprocessingML's (fonts, size,
// color) rather than Office Math's (style, script, normal text).
func isWordRPr(n *ooxml.Node) bool {
	return strings.Contains(n.Space, "wordprocessingml") || strings.Contains(n.Space, "drawingml")
}

// runStyle reads the formatting of a run that changes how a formula is
// drawn: its size and color.
func runStyle(rpr *ooxml.Node) *Styled {
	if rpr == nil {
		return nil
	}
	var s Styled
	set := false
	if sz := rpr.Path("sz").AttrStr("val", ""); sz != "" {
		if v, err := strconv.ParseFloat(sz, 64); err == nil && v > 0 {
			s.Base = v / 2
			set = true
		}
	}
	if c := rpr.Path("color").AttrStr("val", ""); c != "" && c != "auto" {
		if col, ok := ParseColor("#" + c); ok {
			s.Color = &col
			set = true
		}
	}
	if !set {
		return nil
	}
	return &s
}

package equation

import (
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// ParseMathML reads a MathML math element as golang.org/x/net/html parses
// it (presentation MathML, with the deprecated elements pages still use:
// mfenced, mstyle's attributes). display reports display="block".
func ParseMathML(n *html.Node) (formula Node, display bool) {
	display = strings.EqualFold(attrOf(n, "display"), "block") || strings.EqualFold(attrOf(n, "mode"), "display")
	var f Node = mrow(n)
	if v := attrOf(n, "displaystyle"); v != "" {
		f = &Styled{Kid: f, Display: boolStyle(v)}
	}
	return withStyle(n, f), display
}

func attrOf(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, name) {
			return strings.TrimSpace(a.Val)
		}
	}
	return ""
}

func localName(n *html.Node) string {
	name := strings.ToLower(n.Data)
	if i := strings.IndexByte(name, ':'); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// elems returns the element children of n.
func elems(n *html.Node) []*html.Node {
	var out []*html.Node
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if k.Type == html.ElementNode {
			out = append(out, k)
		}
	}
	return out
}

// mrow reads the children of n as a row.
func mrow(n *html.Node) Node {
	var kids []Node
	for _, k := range elems(n) {
		if m := mathml(k); m != nil {
			kids = append(kids, m)
		}
	}
	return row(kids)
}

// arg returns the i-th child of n as a node (an empty row when missing).
func mathArg(kids []*html.Node, i int) Node {
	if i >= len(kids) {
		return &Row{}
	}
	if m := mathml(kids[i]); m != nil {
		return m
	}
	return &Row{}
}

// optArg is mathArg, with nil for a missing child or <none/>.
func optArg(kids []*html.Node, i int) Node {
	if i >= len(kids) || localName(kids[i]) == "none" {
		return nil
	}
	return mathml(kids[i])
}

func mathml(n *html.Node) Node {
	kids := elems(n)
	var out Node
	switch localName(n) {
	case "mi":
		t := tokenText(n)
		a := &Atom{Kind: Ident, Text: t, Variant: mathVariant(attrOf(n, "mathvariant"))}
		if a.Variant == VarAuto {
			if r, ok := singleRune(t); ok && isLetter(r) {
				a.Variant = VarItalic
			} else {
				a.Variant = VarNormal
			}
		}
		out = a
	case "mn":
		out = &Atom{Kind: Number, Text: tokenText(n), Variant: mathVariant(attrOf(n, "mathvariant"))}
	case "mo":
		out = mathOp(n)
	case "mtext", "ms":
		t := tokenText(n)
		if localName(n) == "ms" {
			t = "\"" + t + "\""
		}
		out = &Atom{Kind: Text, Text: t, Variant: mathVariant(attrOf(n, "mathvariant"))}
	case "mspace":
		out = &Space{Width: mathLength(attrOf(n, "width"), 0)}
	case "mrow", "merror", "mpadded", "mstyle", "math":
		out = mrow(n)
		if localName(n) == "mstyle" {
			s := &Styled{Kid: out}
			if v := attrOf(n, "displaystyle"); v != "" {
				s.Display = boolStyle(v)
			}
			if v := attrOf(n, "scriptlevel"); v != "" {
				if strings.HasPrefix(v, "+") {
					s.LevelUp = true
				} else if l, err := strconv.Atoi(v); err == nil && l >= 0 {
					s.Level = int8(min(l, 3) + 1)
				}
			}
			if v := mathVariant(attrOf(n, "mathvariant")); v != VarAuto {
				s.Variant = v
			}
			out = s
		}
	case "mphantom":
		out = &Phantom{Kid: mrow(n)}
	case "mfrac":
		f := &Frac{Num: mathArg(kids, 0), Den: mathArg(kids, 1)}
		switch lt := attrOf(n, "linethickness"); lt {
		case "", "medium":
		case "thin":
			f.Thickness = 0.03
		case "thick":
			f.Thickness = 0.08
		default:
			if v := mathLength(lt, -1); v == 0 {
				f.NoBar = true
			} else if v > 0 {
				f.Thickness = v
			}
		}
		if strings.EqualFold(attrOf(n, "bevelled"), "true") {
			f.Kind = FracSkewed
		}
		f.NumAlign, f.DenAlign = mathAlign(attrOf(n, "numalign")), mathAlign(attrOf(n, "denomalign"))
		out = f
	case "msqrt":
		out = &Radical{Base: mrow(n)}
	case "mroot":
		out = &Radical{Base: mathArg(kids, 0), Degree: mathArg(kids, 1)}
	case "msub":
		out = &Scripts{Base: mathArg(kids, 0), Sub: optArg(kids, 1)}
	case "msup":
		out = &Scripts{Base: mathArg(kids, 0), Sup: optArg(kids, 1)}
	case "msubsup":
		out = &Scripts{Base: mathArg(kids, 0), Sub: optArg(kids, 1), Sup: optArg(kids, 2)}
	case "munder", "mover", "munderover":
		u := &UnderOver{Base: mathArg(kids, 0)}
		switch localName(n) {
		case "munder":
			u.Under = optArg(kids, 1)
		case "mover":
			u.Over = optArg(kids, 1)
		default:
			u.Under, u.Over = optArg(kids, 1), optArg(kids, 2)
		}
		u.AccentOver = accentOf(n, "accent", u.Over)
		u.AccentUnder = accentOf(n, "accentunder", u.Under)
		out = u
	case "mmultiscripts":
		out = multiscripts(kids)
	case "mtable":
		out = mathTable(n)
	case "menclose":
		out = enclose(n)
	case "mfenced":
		out = mfenced(n, kids)
	case "semantics":
		for _, k := range kids {
			switch localName(k) {
			case "annotation", "annotation-xml":
				continue
			}
			return withStyle(n, mathml(k))
		}
		for _, k := range kids {
			if localName(k) == "annotation" && strings.Contains(strings.ToLower(attrOf(k, "encoding")), "tex") {
				return ParseTeX(textOf(k))
			}
		}
		return nil
	case "maction":
		if len(kids) > 0 {
			out = mathml(kids[0])
		}
	case "annotation", "annotation-xml", "none", "mprescripts", "mglyph", "malignmark", "maligngroup":
		return nil
	default:
		out = mrow(n)
	}
	return withStyle(n, out)
}

// withStyle applies an element's mathcolor, mathsize and style color.
func withStyle(n *html.Node, out Node) Node {
	if out == nil {
		return nil
	}
	s := &Styled{Kid: out}
	set := false
	col := attrOf(n, "mathcolor")
	if col == "" {
		col = attrOf(n, "color")
	}
	if col == "" {
		for _, d := range strings.Split(attrOf(n, "style"), ";") {
			if k, v, ok := strings.Cut(d, ":"); ok && strings.EqualFold(strings.TrimSpace(k), "color") {
				col = strings.TrimSpace(v)
			}
		}
	}
	if c, ok := ParseColor(col); ok {
		s.Color, set = &c, true
	}
	switch sz := strings.ToLower(attrOf(n, "mathsize")); {
	case sz == "":
	case sz == "small":
		s.Scale, set = 0.8, true
	case sz == "big":
		s.Scale, set = 1.2, true
	case strings.HasSuffix(sz, "%"):
		if v, err := strconv.ParseFloat(strings.TrimSuffix(sz, "%"), 64); err == nil && v > 0 {
			s.Scale, set = v/100, true
		}
	case strings.HasSuffix(sz, "px"), strings.HasSuffix(sz, "pt"):
		if v, err := strconv.ParseFloat(sz[:len(sz)-2], 64); err == nil && v > 0 {
			if strings.HasSuffix(sz, "px") {
				v *= 0.75
			}
			s.Size, set = v, true
		}
	default:
		if v := mathLength(sz, 0); v > 0 {
			s.Scale, set = v, true
		}
	}
	if !set {
		return out
	}
	return s
}

func boolStyle(v string) int8 {
	if strings.EqualFold(v, "true") {
		return 1
	}
	return 2
}

// tokenText is the text of a token element, its white space trimmed and
// collapsed.
func tokenText(n *html.Node) string {
	return strings.Join(strings.FieldsFunc(textOf(n), unicode.IsSpace), " ")
}

func textOf(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			if k.Type == html.TextNode {
				b.WriteString(k.Data)
			} else {
				walk(k)
			}
		}
	}
	walk(n)
	return b.String()
}

func mathOp(n *html.Node) Node {
	t := tokenText(n)
	a := NewOp(t)
	if _, ok := singleRune(normalizeOp(t)); !ok {
		if IsFunctionName(t) {
			a.Kind, a.Variant, a.Class = Ident, VarNormal, LargeOp
		}
	}
	flag := func(name string, v *bool) {
		switch strings.ToLower(attrOf(n, name)) {
		case "true":
			*v = true
		case "false":
			*v = false
		}
	}
	flag("stretchy", &a.Stretchy)
	flag("symmetric", &a.Symmetric)
	flag("largeop", &a.LargeOp)
	flag("movablelimits", &a.MovableLimits)
	flag("accent", &a.Accent)
	if a.LargeOp {
		a.Class = LargeOp
	}
	switch attrOf(n, "form") {
	case "prefix":
		if a.Class == Bin {
			a.Class = Ord
		}
	}
	ls, rs := attrOf(n, "lspace"), attrOf(n, "rspace")
	if ls != "" || rs != "" {
		a.Spaced = true
		a.LSpace, a.RSpace = mathLength(ls, 0), mathLength(rs, 0)
	}
	if v := attrOf(n, "minsize"); v != "" {
		a.MinSize = mathLength(v, 0)
	}
	if strings.EqualFold(attrOf(n, "separator"), "true") {
		a.Class = Punct
	}
	a.Variant = mathVariant(attrOf(n, "mathvariant"))
	return a
}

// accentOf reads an accent attribute of munder/mover, which defaults to
// what the operator dictionary says of the script.
func accentOf(n *html.Node, name string, script Node) bool {
	switch strings.ToLower(attrOf(n, name)) {
	case "true":
		return true
	case "false":
		return false
	}
	a := coreAtom(script)
	return a != nil && a.Kind == Op && a.Accent
}

func multiscripts(kids []*html.Node) Node {
	s := &Scripts{Base: mathArg(kids, 0)}
	var post, pre []*html.Node
	cur := &post
	for _, k := range kids[min(1, len(kids)):] {
		if localName(k) == "mprescripts" {
			cur = &pre
			continue
		}
		*cur = append(*cur, k)
	}
	pair := func(list []*html.Node) (sub, sup Node) {
		var subs, sups []Node
		for i := 0; i+1 < len(list); i += 2 {
			if a := optArg(list, i); a != nil {
				subs = append(subs, a)
			}
			if a := optArg(list, i+1); a != nil {
				sups = append(sups, a)
			}
		}
		if len(subs) > 0 {
			sub = row(subs)
		}
		if len(sups) > 0 {
			sup = row(sups)
		}
		return sub, sup
	}
	s.Sub, s.Sup = pair(post)
	s.PreSub, s.PreSup = pair(pre)
	return s
}

func mathTable(n *html.Node) Node {
	t := &Table{}
	for _, a := range strings.Fields(attrOf(n, "columnalign")) {
		t.ColAlign = append(t.ColAlign, mathAlign(a))
	}
	if v := strings.Fields(attrOf(n, "rowspacing")); len(v) > 0 {
		t.RowGap = mathLength(v[0], 0)
	}
	if v := strings.Fields(attrOf(n, "columnspacing")); len(v) > 0 {
		t.ColGap = mathLength(v[0], 0)
	}
	if strings.EqualFold(attrOf(n, "displaystyle"), "true") {
		t.Display = true
	}
	for _, tr := range elems(n) {
		name := localName(tr)
		if name != "mtr" && name != "mlabeledtr" {
			continue
		}
		var cells []Node
		for i, td := range elems(tr) {
			if name == "mlabeledtr" && i == 0 {
				continue
			}
			cells = append(cells, withStyle(td, mrow(td)))
		}
		t.Rows = append(t.Rows, cells)
	}
	return t
}

func enclose(n *html.Node) Node {
	kid := mrow(n)
	e := &Enclose{Kid: kid}
	notation := strings.Fields(attrOf(n, "notation"))
	if len(notation) == 0 {
		notation = []string{"longdiv"}
	}
	for _, v := range notation {
		switch v {
		case "box":
			e.Top, e.Bottom, e.Left, e.Right = true, true, true, true
		case "roundedbox":
			e.Round = true
		case "circle":
			e.Circle = true
		case "top":
			e.Top = true
		case "bottom":
			e.Bottom = true
		case "left":
			e.Left = true
		case "right":
			e.Right = true
		case "actuarial":
			e.Top, e.Right = true, true
		case "madruwb":
			e.Bottom, e.Right = true, true
		case "longdiv":
			e.Top, e.Left = true, true
		case "horizontalstrike":
			e.StrikeH = true
		case "verticalstrike":
			e.StrikeV = true
		case "updiagonalstrike", "updiagonalarrow":
			e.StrikeUp = true
		case "downdiagonalstrike":
			e.StrikeDown = true
		case "radical":
			return &Radical{Base: kid}
		}
	}
	return e
}

func mfenced(n *html.Node, kids []*html.Node) Node {
	open, close := "(", ")"
	if v, ok := attrValue(n, "open"); ok {
		open = strings.TrimSpace(v)
	}
	if v, ok := attrValue(n, "close"); ok {
		close = strings.TrimSpace(v)
	}
	seps := []rune(",")
	if v, ok := attrValue(n, "separators"); ok {
		seps = []rune(strings.Join(strings.Fields(v), ""))
	}
	var inner []Node
	for i, k := range kids {
		if i > 0 && len(seps) > 0 {
			s := NewOp(string(seps[min(i-1, len(seps)-1)]))
			s.Class = Punct
			inner = append(inner, s)
		}
		if m := mathml(k); m != nil {
			inner = append(inner, m)
		}
	}
	if open == "" {
		open = "."
	}
	if close == "" {
		close = "."
	}
	return fenced(open, close, row(inner))
}

func attrValue(n *html.Node, name string) (string, bool) {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, name) {
			return a.Val, true
		}
	}
	return "", false
}

func mathAlign(s string) Align {
	switch strings.ToLower(s) {
	case "left":
		return Left
	case "right":
		return Right
	}
	return Center
}

// mathVariant reads a mathvariant attribute.
func mathVariant(s string) Variant {
	switch strings.ToLower(s) {
	case "normal":
		return VarNormal
	case "bold":
		return VarBold
	case "italic":
		return VarItalic
	case "bold-italic":
		return VarBoldItalic
	case "double-struck":
		return VarDoubleStruck
	case "script":
		return VarScript
	case "bold-script":
		return VarBoldScript
	case "fraktur":
		return VarFraktur
	case "bold-fraktur":
		return VarBoldFraktur
	case "sans-serif":
		return VarSansSerif
	case "bold-sans-serif":
		return VarSansSerifBold
	case "sans-serif-italic":
		return VarSansSerifItalic
	case "sans-serif-bold-italic":
		return VarSansSerifBoldItalic
	case "monospace":
		return VarMonospace
	}
	return VarAuto
}

// mathLength reads a MathML length in em (def when it cannot be read).
func mathLength(s string, def float64) float64 {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "":
		return def
	case "veryverythinmathspace":
		return 1.0 / 18
	case "verythinmathspace":
		return 2.0 / 18
	case "thinmathspace":
		return 3.0 / 18
	case "mediummathspace":
		return 4.0 / 18
	case "thickmathspace":
		return 5.0 / 18
	case "verythickmathspace":
		return 6.0 / 18
	case "veryverythickmathspace":
		return 7.0 / 18
	}
	i := len(s)
	for i > 0 && (s[i-1] >= 'a' && s[i-1] <= 'z' || s[i-1] == '%') {
		i--
	}
	v, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return def
	}
	switch s[i:] {
	case "", "em":
		return v
	case "ex":
		return v * 0.43
	case "px":
		return v / 16
	case "pt":
		return v / 12
	case "mu":
		return v / 18
	case "%":
		return v / 100
	case "cm":
		return v * 28.35 / 12
	case "mm":
		return v * 2.835 / 12
	case "in":
		return v * 6
	}
	return def
}

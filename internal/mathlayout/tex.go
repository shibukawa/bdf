package mathlayout

import (
	"strconv"
	"strings"
	"unicode"
)

// maxDepth is how deeply the groups and the arguments of a formula may
// nest, and the elements of MathML and of Office Math: what is nested
// deeper is read as a part of what it is in. The readers and the layout
// call themselves for what is nested, and the layout looks into a row of
// one for what it holds, which takes the time of the depth at each level.
const maxDepth = 200

// ParseTeX reads a formula written in LaTeX's math mode, with the commands
// of amsmath and amssymb that Markdown documents and MathJax and KaTeX
// pages use. What it does not know it keeps as text. Groups and arguments
// nested deeper than maxDepth are read without their braces and commands.
func ParseTeX(src string) Node { return parseTeX(src, 0) }

// parseTeX reads a formula that is nested in another at a depth (the
// optional argument of a command, which is a formula).
func parseTeX(src string, depth int) Node {
	if depth >= maxDepth {
		return &Atom{Kind: Text, Text: src}
	}
	p := &texParser{s: []rune(src), depth: depth}
	nodes, _ := p.list(texTop)
	if rows := p.rows; rows != nil {
		// line breaks at the top level: the lines are centered below each
		// other, as in gather
		rows = append(rows, []Node{row(nodes)})
		return &Table{Rows: rows, Display: true, RowGap: 0.3}
	}
	return row(nodes)
}

type texParser struct {
	s     []rune
	i     int
	rows  [][]Node // lines ended by \\ at the top level
	depth int      // how deeply the list being read is nested
	flat  int      // open braces that were read over, their groups being too deep
}

// deep reports whether what is read is nested as deeply as it may be:
// what would nest deeper is read as a part of it.
func (p *texParser) deep() bool { return p.depth >= maxDepth }

// what ended a list
const (
	texTop   = iota // the end of the input
	texGroup        // }
	texCell         // & or \\ in an environment
	texRight        // \right
	texOpt          // ] of an optional argument
)

type texEnd struct {
	kind  int    // what the list was waiting for
	what  string // "}", "&", "\\", "\\end", "\\right", "\\middle", "]", ""
	delim string // the delimiter after \right or \middle; the name after \end
}

func row(nodes []Node) Node {
	if len(nodes) == 1 {
		return nodes[0]
	}
	return &Row{Kids: nodes}
}

func (p *texParser) eof() bool { return p.i >= len(p.s) }

func (p *texParser) peek() rune {
	if p.eof() {
		return 0
	}
	return p.s[p.i]
}

func (p *texParser) skipSpace() {
	for !p.eof() {
		switch r := p.s[p.i]; {
		case r == '%':
			for !p.eof() && p.s[p.i] != '\n' {
				p.i++
			}
		case unicode.IsSpace(r):
			p.i++
		default:
			return
		}
	}
}

// command reads a control sequence after its backslash.
func (p *texParser) command() string {
	start := p.i
	for !p.eof() && isASCIILetter(p.s[p.i]) {
		p.i++
	}
	if p.i == start && !p.eof() {
		p.i++
	}
	return string(p.s[start:p.i])
}

func isASCIILetter(r rune) bool { return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' }

// peekCommand returns the control sequence at the current position without
// reading it ("" when there is none).
func (p *texParser) peekCommand() string {
	if p.peek() != '\\' {
		return ""
	}
	save := p.i
	p.i++
	c := p.command()
	p.i = save
	return c
}

// list reads nodes until what kind waits for. Its callers ask deep before
// they call it.
func (p *texParser) list(kind int) ([]Node, texEnd) {
	p.depth++
	defer func() { p.depth-- }()
	var nodes []Node
	for {
		p.skipSpace()
		if p.eof() {
			return nodes, texEnd{kind: kind}
		}
		r := p.peek()
		switch {
		case r == '}':
			p.i++
			if p.flat > 0 {
				p.flat-- // of a group read over
				continue
			}
			if kind == texGroup || kind == texTop {
				return nodes, texEnd{kind: kind, what: "}"}
			}
			continue
		case r == ']' && kind == texOpt:
			p.i++
			return nodes, texEnd{kind: kind, what: "]"}
		case r == '&':
			p.i++
			if kind == texCell {
				return nodes, texEnd{kind: kind, what: "&"}
			}
			continue
		case r == '^' || r == '_':
			nodes = p.scriptsOf(nodes, &Row{})
			continue
		case r == '\\':
			switch c := p.peekCommand(); c {
			case "\\", "cr", "newline":
				p.i += 1 + len(c)
				p.optional() // \\[2pt]
				if kind == texCell {
					return nodes, texEnd{kind: kind, what: "\\\\"}
				}
				if kind == texTop {
					p.rows = append(p.rows, []Node{row(nodes)})
					nodes = nil
				}
				continue
			case "end":
				p.i += 4
				name := p.textArg()
				if kind == texCell {
					return nodes, texEnd{kind: kind, what: "\\end", delim: name}
				}
				continue
			case "right", "middle":
				if kind == texRight {
					p.i += 1 + len(c)
					return nodes, texEnd{kind: kind, what: "\\" + c, delim: p.delimiter()}
				}
				p.i += 1 + len(c)
				p.delimiter()
				continue
			case "over", "choose", "atop", "brace", "brack":
				p.i += 1 + len(c)
				if p.deep() {
					continue
				}
				den, end := p.list(kind)
				f := &Frac{Num: row(nodes), Den: row(den)}
				var n Node = f
				switch c {
				case "choose":
					f.NoBar = true
					n = fenced("(", ")", f)
				case "brace":
					f.NoBar = true
					n = fenced("{", "}", f)
				case "brack":
					f.NoBar = true
					n = fenced("[", "]", f)
				case "atop":
					f.NoBar = true
				}
				return []Node{n}, end
			case "displaystyle", "textstyle", "scriptstyle", "scriptscriptstyle":
				p.i += 1 + len(c)
				if p.deep() {
					continue
				}
				rest, end := p.list(kind)
				s := &Styled{Kid: row(rest)}
				switch c {
				case "displaystyle":
					s.Display, s.Level = 1, 1
				case "textstyle":
					s.Display, s.Level = 2, 1
				case "scriptstyle":
					s.Display, s.Level = 2, 2
				default:
					s.Display, s.Level = 2, 3
				}
				return append(nodes, s), end
			case "color":
				p.i += 6
				col, ok := ParseColor(p.colorArg())
				if p.deep() {
					continue
				}
				rest, end := p.list(kind)
				var n Node = row(rest)
				if ok {
					n = &Styled{Kid: n, Color: &col}
				}
				return append(nodes, n), end
			case "limits", "nolimits", "displaylimits":
				p.i += 1 + len(c)
				continue
			}
			if v, ok := texOldFonts[p.peekCommand()]; ok {
				p.i += 1 + len(p.peekCommand())
				if p.deep() {
					continue
				}
				rest, end := p.list(kind)
				return append(nodes, &Styled{Kid: row(rest), Variant: v}), end
			}
		case r == '\'':
			nodes = p.scriptsOf(nodes, &Row{})
			continue
		}
		n := p.primary()
		if n == nil {
			continue
		}
		nodes = p.scriptsOf(append(nodes, n), nil)
	}
}

// scriptsOf reads the scripts, primes and limits after the last node (or
// after an empty base when base is set) and attaches them.
func (p *texParser) scriptsOf(nodes []Node, empty Node) []Node {
	var base Node
	if empty != nil {
		base = empty
	} else {
		base = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
	}
	var sub, sup Node
	limits := 0 // 1 \limits, -1 \nolimits
	primes := 0
	for {
		p.skipSpace()
		switch p.peek() {
		case '^':
			p.i++
			if sup != nil {
				// a double superscript: TeX refuses it; keep both
				sup = &Row{Kids: []Node{sup, p.arg()}}
			} else {
				sup = p.arg()
			}
			continue
		case '_':
			p.i++
			sub = p.arg()
			continue
		case '\'':
			p.i++
			primes++
			continue
		case '\\':
			switch p.peekCommand() {
			case "limits":
				p.i += 7
				limits = 1
				continue
			case "nolimits":
				p.i += 9
				limits = -1
				continue
			case "displaylimits":
				p.i += 14
				continue
			}
		}
		break
	}
	if primes > 0 {
		pr := &Atom{Kind: Op, Text: strings.Repeat("′", primes), Class: Ord}
		if primes <= 4 {
			pr.Text = string([]rune("′″‴⁗")[primes-1])
		}
		if sup != nil {
			sup = &Row{Kids: []Node{pr, sup}}
		} else {
			sup = pr
		}
	}
	if sub == nil && sup == nil {
		return append(nodes, base)
	}
	if limits == 1 || limits == 0 && hasLimits(base) {
		if a := coreAtom(base); a != nil && limits == 1 {
			a.MovableLimits = false
		}
		return append(nodes, &UnderOver{Base: base, Under: sub, Over: sup})
	}
	return append(nodes, &Scripts{Base: base, Sub: sub, Sup: sup})
}

// hasLimits reports whether scripts go under and over a base (in display
// style): large operators other than integrals, functions like lim, and
// \mathop groups.
func hasLimits(n Node) bool {
	if r, ok := n.(*Row); ok && r.Class == LargeOp {
		return true
	}
	if u, ok := n.(*UnderOver); ok {
		// \overbrace and \underbrace
		for _, k := range []Node{u.Over, u.Under} {
			if a := coreAtom(k); a != nil && a.Stretchy && strings.ContainsRune("⏞⏟⏜⏝⎴⎵", []rune(a.Text)[0]) {
				return true
			}
		}
	}
	a := coreAtom(n)
	if a == nil {
		return false
	}
	if a.LargeOp {
		return a.MovableLimits
	}
	return a.Kind == Ident && functionNames[a.Text]
}

// arg reads the argument of a command or script: a group or one token.
// Nested too deeply, the argument is empty, and what would have been it is
// read after the command.
func (p *texParser) arg() Node {
	p.skipSpace()
	switch r := p.peek(); {
	case r == '{':
		p.i++
		if p.deep() {
			p.flat++
			return &Row{}
		}
		nodes, _ := p.list(texGroup)
		return row(nodes)
	case r == 0:
		return &Row{}
	case r >= '0' && r <= '9':
		p.i++
		return &Atom{Kind: Number, Text: string(r)}
	}
	if n := p.nested(); n != nil {
		return n
	}
	return &Row{}
}

// nested reads one token or construct for a command that takes it (the
// token after \sqrt, which may be a command that takes one again); nil
// when it is nested too deeply, and the token is read after the command.
func (p *texParser) nested() Node {
	if p.deep() {
		return nil
	}
	p.depth++
	defer func() { p.depth-- }()
	return p.primary()
}

// textArg reads a group as plain text.
func (p *texParser) textArg() string {
	p.skipSpace()
	if p.peek() != '{' {
		if p.eof() {
			return ""
		}
		p.i++
		return string(p.s[p.i-1])
	}
	p.i++
	start, depth := p.i, 1
	for !p.eof() {
		switch p.s[p.i] {
		case '\\':
			p.i++
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				s := string(p.s[start:p.i])
				p.i++
				return s
			}
		}
		p.i++
	}
	return string(p.s[start:])
}

// optional reads an optional argument in brackets as text.
func (p *texParser) optional() (string, bool) {
	p.skipSpace()
	if p.peek() != '[' {
		return "", false
	}
	p.i++
	start, depth := p.i, 0
	for !p.eof() {
		switch p.s[p.i] {
		case '{':
			depth++
		case '}':
			depth--
		case ']':
			if depth == 0 {
				s := string(p.s[start:p.i])
				p.i++
				return s, true
			}
		}
		p.i++
	}
	return string(p.s[start:]), true
}

// colorArg reads the color of \color and \textcolor: a name, or a model in
// brackets and its values.
func (p *texParser) colorArg() string {
	model, hasModel := p.optional()
	v := strings.TrimSpace(p.textArg())
	if !hasModel {
		return v
	}
	parts := strings.Split(v, ",")
	num := func(i int, scale float64) int {
		if i >= len(parts) {
			return 0
		}
		f, _ := strconv.ParseFloat(strings.TrimSpace(parts[i]), 64)
		return int(f*scale + 0.5)
	}
	switch strings.TrimSpace(model) {
	case "rgb":
		return "rgb(" + strconv.Itoa(num(0, 255)) + "," + strconv.Itoa(num(1, 255)) + "," + strconv.Itoa(num(2, 255)) + ")"
	case "RGB":
		return "rgb(" + strconv.Itoa(num(0, 1)) + "," + strconv.Itoa(num(1, 1)) + "," + strconv.Itoa(num(2, 1)) + ")"
	case "HTML":
		return "#" + v
	case "gray":
		g := strconv.Itoa(num(0, 255))
		return "rgb(" + g + "," + g + "," + g + ")"
	}
	return v
}

// delimiter reads the delimiter after \left, \right, \middle or \big.
func (p *texParser) delimiter() string {
	p.skipSpace()
	if p.eof() {
		return "."
	}
	r := p.s[p.i]
	p.i++
	if r != '\\' {
		if r == '<' {
			return "⟨"
		}
		if r == '>' {
			return "⟩"
		}
		return string(r)
	}
	c := p.command()
	if d, ok := texDelims[c]; ok {
		return d
	}
	if d, ok := texOrd[c]; ok {
		return d
	}
	if d, ok := texOps[c]; ok {
		return d
	}
	return "."
}

// fenced surrounds n with delimiters that grow to it.
func fenced(open, close string, kids ...Node) Node {
	r := &Row{Class: Inner}
	if open != "." && open != "" {
		o := NewOp(open)
		o.Stretchy, o.Symmetric, o.Class = true, true, Open
		r.Kids = append(r.Kids, o)
	}
	r.Kids = append(r.Kids, kids...)
	if close != "." && close != "" {
		c := NewOp(close)
		c.Stretchy, c.Symmetric, c.Class = true, true, Close
		r.Kids = append(r.Kids, c)
	}
	return r
}

// texDelim is a delimiter of \left, \middle or \right.
func texDelim(d string, class Class) Node {
	if d == "." || d == "" {
		return &Space{Width: 0.12}
	}
	a := NewOp(d)
	a.Stretchy, a.Symmetric, a.Shortfall, a.Class = true, true, true, class
	return a
}

// primary reads one token or construct.
func (p *texParser) primary() Node {
	r := p.peek()
	switch {
	case r == '{':
		p.i++
		if p.deep() {
			p.flat++ // read over, with the brace that closes it
			return nil
		}
		nodes, _ := p.list(texGroup)
		return &Row{Kids: nodes}
	case r == '\\':
		p.i++
		return p.control(p.command())
	case r >= '0' && r <= '9' || r == '.' && p.i+1 < len(p.s) && p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9':
		start := p.i
		for !p.eof() && (p.s[p.i] >= '0' && p.s[p.i] <= '9' || p.s[p.i] == '.' && p.i+1 < len(p.s) && p.s[p.i+1] >= '0' && p.s[p.i+1] <= '9') {
			p.i++
		}
		return &Atom{Kind: Number, Text: string(p.s[start:p.i])}
	case r == '~':
		p.i++
		return &Space{Width: 0.25}
	case r == '$' || r == '#':
		p.i++
		return nil
	}
	p.i++
	if unicode.IsLetter(r) {
		return &Atom{Kind: Ident, Text: string(r)}
	}
	if unicode.IsDigit(r) {
		return &Atom{Kind: Number, Text: string(r)}
	}
	a := NewOp(string(r))
	// only \left and \right grow delimiters in TeX
	a.Stretchy = false
	if r == '/' {
		a.Class = Ord
	}
	if r == '|' {
		a.Class = Ord
	}
	return a
}

// control reads what a control sequence stands for.
func (p *texParser) control(c string) Node {
	if s, ok := texSpaces[c]; ok {
		return &Space{Width: s}
	}
	if s, ok := texOrd[c]; ok {
		kind := Ident
		if r, _ := singleRune(s); !unicode.IsLetter(r) {
			kind = Op
		}
		if kind == Op {
			a := NewOp(s)
			a.Stretchy = false
			if c == "|" || c == "vert" || c == "Vert" {
				a.Class = Ord
			}
			if c == "lvert" || c == "lVert" {
				a.Class = Open
			}
			if c == "rvert" || c == "rVert" {
				a.Class = Close
			}
			if strings.HasSuffix(c, "dots") || c == "ldots" || c == "dots" {
				a.Class = Inner
			}
			return a
		}
		return &Atom{Kind: Ident, Text: s}
	}
	if s, ok := texVarGreek[c]; ok {
		return &Atom{Kind: Ident, Text: s, Variant: VarItalic}
	}
	if s, ok := texOps[c]; ok {
		a := NewOp(s)
		a.Stretchy = a.Stretchy && lookupOp(s).horizontal
		if c == "colon" {
			a.Class = Punct
		}
		return a
	}
	if s, ok := texLargeOps[c]; ok {
		return NewOp(s)
	}
	if s, ok := texDelims[c]; ok {
		a := NewOp(s)
		a.Stretchy = false
		return a
	}
	if movable, ok := functionNames[c]; ok {
		return &Atom{Kind: Ident, Text: c, Variant: VarNormal, Class: LargeOp, MovableLimits: movable}
	}
	if v, ok := texVariants[c]; ok {
		return &Styled{Kid: p.arg(), Variant: v}
	}
	if acc, ok := texAccents[c]; ok {
		a := NewOp(acc.char)
		a.Accent, a.Stretchy = true, acc.stretch
		return &UnderOver{Base: p.arg(), Over: a, AccentOver: true}
	}
	if s, ok := texUnderAccents[c]; ok {
		a := NewOp(s)
		a.Accent, a.Stretchy = true, true
		return &UnderOver{Base: p.arg(), Under: a, AccentUnder: true}
	}
	if size, ok := texBig[strings.TrimRight(c, "lrm")]; ok {
		a := NewOp(p.delimiter())
		a.Stretchy, a.Symmetric, a.MinSize = true, true, size
		switch c[len(c)-1] {
		case 'l':
			a.Class = Open
		case 'r':
			a.Class = Close
		case 'm':
			a.Class = Rel
		default:
			a.Class = Ord
		}
		return a
	}
	switch c {
	case "frac", "dfrac", "tfrac", "cfrac":
		f := &Frac{Num: p.arg(), Den: p.arg()}
		switch c {
		case "dfrac", "cfrac":
			return &Styled{Kid: f, Display: 1, Level: 1}
		case "tfrac":
			return &Styled{Kid: f, Display: 2, Level: 1}
		}
		return f
	case "binom", "dbinom", "tbinom":
		n := fenced("(", ")", &Frac{Num: p.arg(), Den: p.arg(), NoBar: true})
		switch c {
		case "dbinom":
			return &Styled{Kid: n, Display: 1, Level: 1}
		case "tbinom":
			return &Styled{Kid: n, Display: 2, Level: 1}
		}
		return n
	case "sqrt":
		if deg, ok := p.optional(); ok {
			return &Radical{Base: p.arg(), Degree: parseTeX(deg, p.depth)}
		}
		return &Radical{Base: p.arg()}
	case "root":
		// \root n \of x
		if p.deep() {
			return nil
		}
		var deg []Node
		for !p.eof() && p.peekCommand() != "of" {
			if n := p.nested(); n != nil {
				deg = append(deg, n)
			}
			p.skipSpace()
		}
		p.i += 3
		return &Radical{Base: p.arg(), Degree: row(deg)}
	case "left":
		open := p.delimiter()
		if p.deep() {
			// the delimiter as it is; the one of \right is read over
			if open == "." {
				return nil
			}
			a := NewOp(open)
			a.Stretchy = false
			return a
		}
		var kids []Node
		kids = append(kids, texDelim(open, Open))
		for {
			nodes, end := p.list(texRight)
			kids = append(kids, nodes...)
			if end.what == "\\middle" {
				kids = append(kids, texDelim(end.delim, Rel))
				continue
			}
			kids = append(kids, texDelim(end.delim, Close))
			break
		}
		return &Row{Kids: kids, Class: Inner}
	case "overline":
		return &Bar{Kid: p.arg()}
	case "underline":
		return &Bar{Kid: p.arg(), Under: true}
	case "overbrace", "underbrace", "overparen", "underparen", "overbracket", "underbracket":
		ch := map[string]string{"overbrace": "⏞", "underbrace": "⏟", "overparen": "⏜", "underparen": "⏝",
			"overbracket": "⎴", "underbracket": "⎵"}[c]
		a := NewOp(ch)
		a.Stretchy, a.Accent = true, true
		base := p.arg()
		var n Node
		if strings.HasPrefix(c, "over") {
			n = &UnderOver{Base: base, Over: a, AccentOver: true}
		} else {
			n = &UnderOver{Base: base, Under: a, AccentUnder: true}
		}
		// its label is a limit (see hasLimits)
		return n
	case "overset", "stackrel":
		over := p.arg()
		base := p.arg()
		r := &Row{Kids: []Node{&UnderOver{Base: base, Over: over}}}
		if c == "stackrel" {
			r.Class = Rel
		}
		return r
	case "underset":
		under := p.arg()
		return &Row{Kids: []Node{&UnderOver{Base: p.arg(), Under: under}}}
	case "xrightarrow", "xleftarrow", "xleftrightarrow", "xRightarrow", "xLeftarrow", "xLeftrightarrow",
		"xmapsto", "xhookrightarrow", "xhookleftarrow", "xlongequal", "xtwoheadrightarrow", "xrightleftharpoons":
		ch := map[string]string{"xrightarrow": "→", "xleftarrow": "←", "xleftrightarrow": "↔", "xRightarrow": "⇒",
			"xLeftarrow": "⇐", "xLeftrightarrow": "⇔", "xmapsto": "↦", "xhookrightarrow": "↪",
			"xhookleftarrow": "↩", "xlongequal": "=", "xtwoheadrightarrow": "↠", "xrightleftharpoons": "⇌"}[c]
		var under Node
		if s, ok := p.optional(); ok {
			under = padded(parseTeX(s, p.depth))
		}
		over := padded(p.arg())
		a := NewOp(ch)
		a.Stretchy = true
		return &Row{Kids: []Node{&UnderOver{Base: a, Over: over, Under: under}}, Class: Rel}
	case "text", "textrm", "textnormal", "mbox", "hbox", "textup", "textmd", "textsf", "texttt", "emph", "textit", "textbf", "textsl":
		t := &Atom{Kind: Text, Text: texText(p.textArg())}
		switch c {
		case "textit", "emph", "textsl":
			t.Variant = VarItalic
		case "textbf":
			t.Variant = VarBold
		}
		return t
	case "operatorname", "operatornamewithlimits":
		star := false
		if p.peek() == '*' {
			p.i++
			star = true
		}
		name := &Atom{Kind: Ident, Text: texText(p.textArg()), Variant: VarNormal, Class: LargeOp}
		if star || c == "operatornamewithlimits" {
			name.MovableLimits = true
			return &Row{Kids: []Node{name}, Class: LargeOp}
		}
		return name
	case "mathop":
		return &Row{Kids: []Node{p.arg()}, Class: LargeOp}
	case "mathbin":
		return &Row{Kids: []Node{p.arg()}, Class: Bin}
	case "mathrel":
		return &Row{Kids: []Node{p.arg()}, Class: Rel}
	case "mathopen":
		return &Row{Kids: []Node{p.arg()}, Class: Open}
	case "mathclose":
		return &Row{Kids: []Node{p.arg()}, Class: Close}
	case "mathpunct":
		return &Row{Kids: []Node{p.arg()}, Class: Punct}
	case "mathord":
		return &Row{Kids: []Node{p.arg()}, Class: Ord}
	case "mathinner":
		return &Row{Kids: []Node{p.arg()}, Class: Inner}
	case "textcolor":
		col, ok := ParseColor(p.colorArg())
		n := p.arg()
		if ok {
			return &Styled{Kid: n, Color: &col}
		}
		return n
	case "colorbox", "fcolorbox":
		if c == "fcolorbox" {
			p.textArg()
		}
		p.textArg()
		return p.arg()
	case "boxed", "fbox", "framebox":
		return &Enclose{Kid: p.arg(), Top: true, Bottom: true, Left: true, Right: true}
	case "cancel":
		return &Enclose{Kid: p.arg(), StrikeUp: true}
	case "bcancel":
		return &Enclose{Kid: p.arg(), StrikeDown: true}
	case "xcancel":
		return &Enclose{Kid: p.arg(), StrikeUp: true, StrikeDown: true}
	case "sout":
		return &Enclose{Kid: p.arg(), StrikeH: true}
	case "phantom":
		return &Phantom{Kid: p.arg()}
	case "hphantom":
		return &Phantom{Kid: p.arg(), ZeroAsc: true, ZeroDesc: true}
	case "vphantom":
		return &Phantom{Kid: p.arg(), ZeroWidth: true}
	case "smash":
		p.optional()
		return &Phantom{Kid: p.arg(), Show: true, ZeroAsc: true, ZeroDesc: true}
	case "not":
		p.skipSpace()
		n := p.nested()
		if a, ok := n.(*Atom); ok {
			if neg, ok := texNegated[normalizeOp(a.Text)]; ok {
				a.Text = neg
				return a
			}
			a.Text += "̸"
			return a
		}
		return n
	case "pmod":
		return &Row{Kids: []Node{&Space{Width: 1}, NewOp("("), &Atom{Kind: Ident, Text: "mod", Variant: VarNormal, Class: Ord},
			&Space{Width: 1.0 / 3}, p.arg(), NewOp(")")}}
	case "bmod":
		return &Atom{Kind: Ident, Text: "mod", Variant: VarNormal, Class: Bin}
	case "mod":
		return &Row{Kids: []Node{&Space{Width: 1}, &Atom{Kind: Ident, Text: "mod", Variant: VarNormal, Class: Ord}, &Space{Width: 1.0 / 3}}}
	case "hspace", "kern", "mkern", "mskip", "hskip":
		var s string
		if c == "hspace" {
			p.skipSpace()
			if p.peek() == '*' {
				p.i++
			}
			s = p.textArg()
		} else {
			s = p.dimension()
		}
		return &Space{Width: texLength(s)}
	case "begin":
		name := p.textArg()
		if p.deep() {
			return nil // its cells are read one after the other
		}
		return p.environment(name)
	case "tag", "label", "ref", "eqref", "notag", "nonumber", "hline", "hdashline", "strut", "mathstrut",
		"allowbreak", "nolimits", "limits", "displaylimits", "relax", "vphantomstrut":
		if c == "tag" || c == "label" || c == "ref" || c == "eqref" {
			p.skipSpace()
			if p.peek() == '*' {
				p.i++
			}
			p.textArg()
		}
		return nil
	}
	// unknown: keep its name as text
	return &Atom{Kind: Text, Text: "\\" + c}
}

// padded puts space around the label of an extensible arrow.
func padded(n Node) Node {
	return &Row{Kids: []Node{&Space{Width: 0.3}, n, &Space{Width: 0.3}}}
}

// texText turns the text argument of \text into its characters.
func texText(s string) string {
	r := strings.NewReplacer("\\,", " ", "\\ ", " ", "~", " ", "\\&", "&", "\\%", "%", "\\$", "$", "\\#", "#",
		"\\_", "_", "\\{", "{", "\\}", "}", "{", "", "}", "", "$", "")
	return r.Replace(s)
}

// dimension reads the length after \kern or \mkern.
func (p *texParser) dimension() string {
	p.skipSpace()
	start := p.i
	for !p.eof() && (p.s[p.i] == '-' || p.s[p.i] == '.' || p.s[p.i] >= '0' && p.s[p.i] <= '9') {
		p.i++
	}
	for k := 0; k < 2 && !p.eof() && isASCIILetter(p.s[p.i]); k++ {
		p.i++
	}
	return string(p.s[start:p.i])
}

// texLength converts a TeX length to em.
func texLength(s string) float64 {
	s = strings.TrimSpace(s)
	i := len(s)
	for i > 0 && isASCIILetter(rune(s[i-1])) {
		i--
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(s[:i]), 64)
	if err != nil {
		return 0
	}
	switch s[i:] {
	case "em":
		return v
	case "ex":
		return v * 0.43
	case "mu":
		return v / 18
	case "pt":
		return v / 10
	case "px":
		return v / 13.33
	case "cm":
		return v * 2.845
	case "mm":
		return v * 0.2845
	case "in":
		return v * 7.227
	}
	return v / 10
}

// environment reads \begin{name} … \end{name}.
func (p *texParser) environment(name string) Node {
	base := strings.TrimSuffix(name, "*")
	var colAlign []Align
	if base == "array" || base == "subarray" || base == "alignat" || base == "alignedat" {
		spec := p.textArg()
		for _, r := range spec {
			switch r {
			case 'l':
				colAlign = append(colAlign, Left)
			case 'c':
				colAlign = append(colAlign, Center)
			case 'r':
				colAlign = append(colAlign, Right)
			}
		}
	}
	var rows [][]Node
	var cur []Node
	for {
		nodes, end := p.list(texCell)
		cur = append(cur, row(nodes))
		if end.what == "&" {
			continue
		}
		if !(len(cur) == 1 && len(nodes) == 0 && end.what != "&") || len(rows) == 0 && end.what != "\\end" {
			rows = append(rows, cur)
		}
		cur = nil
		if end.what != "\\\\" {
			break
		}
	}
	t := &Table{Rows: rows, ColAlign: colAlign}
	switch base {
	case "matrix", "array", "subarray", "smallmatrix":
		if base == "smallmatrix" || base == "subarray" {
			return &Styled{Kid: t, Level: 2, Display: 2}
		}
		return t
	case "pmatrix":
		return fenced("(", ")", t)
	case "bmatrix":
		return fenced("[", "]", t)
	case "Bmatrix":
		return fenced("{", "}", t)
	case "vmatrix":
		return fenced("|", "|", t)
	case "Vmatrix":
		return fenced("‖", "‖", t)
	case "cases", "dcases":
		t.ColAlign, t.ColGap = []Align{Left, Left}, 1
		if base == "dcases" {
			t.Display = true
		}
		return fenced("{", ".", t)
	case "rcases":
		t.ColAlign, t.ColGap = []Align{Left, Left}, 1
		return fenced(".", "}", t)
	case "aligned", "align", "alignat", "alignedat", "split", "eqnarray", "flalign":
		t.Aligned, t.Display = true, true
		return t
	case "gathered", "gather", "multline", "equation", "displaymath", "math":
		if len(rows) == 1 && len(rows[0]) == 1 {
			return rows[0][0]
		}
		t.Display, t.RowGap = true, 0.3
		return t
	}
	return t
}

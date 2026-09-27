package epub

import (
	"slices"
	"strings"

	"golang.org/x/net/html"
)

// The book's style sheets are not applied (the text takes the reader
// style), but a few properties that carry meaning in books are read from
// them: the writing mode (Japanese books are vertical), tate-chu-yoko,
// emphasis marks, alignment, hidden elements and the size of pictures
// (external characters, gaiji, are pictures one em wide). The rules read
// are those whose selectors are one compound selector: a type, classes, an
// id or :root (".vrtl", "span.tcy", "html", "img.gaiji"); rules with
// combinators or attribute and pseudo-class selectors are left out.

// honored are the properties read from the style sheets.
var honored = map[string]bool{
	"writing-mode": true, "-epub-writing-mode": true, "-webkit-writing-mode": true,
	"text-combine-upright": true, "-epub-text-combine": true, "-webkit-text-combine": true,
	"text-emphasis-style": true, "text-emphasis": true, "-epub-text-emphasis-style": true,
	"-webkit-text-emphasis-style": true, "-webkit-text-emphasis": true,
	"text-align": true, "display": true, "width": true, "height": true,
}

// pictureOnly are the honored properties that apply to pictures only.
var pictureOnly = map[string]bool{"width": true, "height": true}

// maxImports is how deeply style sheets may import others.
const maxImports = 8

// cssRule is a rule of a style sheet with the honored declarations.
type cssRule struct {
	sel   selector
	spec  int // specificity: ids, classes, types
	order int
	decls [][2]string
}

// selector is a compound selector.
type selector struct {
	tag     string // "" for any element
	id      string
	classes []string
}

func (s *selector) matches(n *html.Node, classes []string) bool {
	if s.tag != "" && s.tag != n.Data {
		return false
	}
	if s.id != "" && attrVal(n, "id") != s.id {
		return false
	}
	for _, c := range s.classes {
		if !slices.Contains(classes, c) {
			return false
		}
	}
	return true
}

// parseCSS reads the rules of a style sheet; load returns the text of an
// imported one by its URL (relative to the sheet), or "" when there is none.
func parseCSS(src string, load func(url string) string, depth int, rules *[]cssRule) {
	src = stripComments(src)
	p := &cssParser{s: src}
	p.sheet(load, depth, rules)
}

type cssParser struct {
	s string
	i int
}

// sheet reads rules up to the end of the text or the '}' that closes a
// block.
func (p *cssParser) sheet(load func(string) string, depth int, rules *[]cssRule) {
	for {
		p.space()
		if p.i >= len(p.s) {
			return
		}
		if p.s[p.i] == '}' {
			p.i++
			return
		}
		prelude, end := p.until("{;")
		prelude = strings.TrimSpace(prelude)
		if end == ';' || end == 0 {
			// an at-rule without a block: @import, @charset, @namespace
			if url, ok := importURL(prelude); ok && depth < maxImports && load != nil {
				if text := load(url); text != "" {
					parseCSS(text, func(u string) string { return load(joinURL(url, u)) }, depth+1, rules)
				}
			}
			continue
		}
		if strings.HasPrefix(prelude, "@") {
			name := strings.ToLower(strings.Fields(prelude)[0])
			switch name {
			case "@media", "@supports", "@layer", "@document", "@container":
				// the conditions are not tested: the rules inside apply
				p.sheet(load, depth, rules)
			default:
				p.skipBlock() // @font-face, @page, @keyframes …
			}
			continue
		}
		body, _ := p.until("}")
		decls := declarations(body)
		if len(decls) == 0 {
			continue
		}
		for _, s := range strings.Split(prelude, ",") {
			sel, spec, ok := parseSelector(strings.TrimSpace(s))
			if !ok {
				continue
			}
			*rules = append(*rules, cssRule{sel: sel, spec: spec, order: len(*rules), decls: decls})
		}
	}
}

func (p *cssParser) space() {
	for p.i < len(p.s) && isCSSSpace(p.s[p.i]) {
		p.i++
	}
}

func isCSSSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' }

// until reads up to one of the stop characters outside strings and
// parentheses, and consumes it; end is 0 at the end of the text.
func (p *cssParser) until(stops string) (text string, end byte) {
	start := p.i
	depth := 0
	for p.i < len(p.s) {
		c := p.s[p.i]
		switch {
		case c == '"' || c == '\'':
			p.i++
			for p.i < len(p.s) && p.s[p.i] != c {
				if p.s[p.i] == '\\' {
					p.i++
				}
				p.i++
			}
		case c == '\\':
			p.i++
		case c == '(':
			depth++
		case c == ')':
			depth = max(depth-1, 0)
		case depth == 0 && strings.IndexByte(stops, c) >= 0:
			text = p.s[start:p.i]
			p.i++
			return text, c
		}
		p.i++
	}
	return p.s[start:min(p.i, len(p.s))], 0
}

// skipBlock skips the rest of a block whose '{' was read.
func (p *cssParser) skipBlock() {
	depth := 1
	for p.i < len(p.s) && depth > 0 {
		_, end := p.until("{}")
		switch end {
		case '{':
			depth++
		case '}':
			depth--
		default:
			return
		}
	}
}

func stripComments(s string) string {
	var b strings.Builder
	for {
		i := strings.Index(s, "/*")
		if i < 0 {
			b.WriteString(s)
			return b.String()
		}
		b.WriteString(s[:i])
		j := strings.Index(s[i+2:], "*/")
		if j < 0 {
			return b.String()
		}
		b.WriteByte(' ')
		s = s[i+2+j+2:]
	}
}

// importURL reads the URL of an @import rule.
func importURL(prelude string) (string, bool) {
	rest, ok := strings.CutPrefix(prelude, "@import")
	if !ok {
		return "", false
	}
	rest = strings.TrimSpace(rest)
	if r, ok := strings.CutPrefix(rest, "url("); ok {
		end := strings.IndexByte(r, ')')
		if end < 0 {
			return "", false
		}
		return strings.Trim(strings.TrimSpace(r[:end]), `"'`), true
	}
	if len(rest) > 1 && (rest[0] == '"' || rest[0] == '\'') {
		if end := strings.IndexByte(rest[1:], rest[0]); end >= 0 {
			return rest[1 : end+1], true
		}
	}
	return "", false
}

// declarations reads the honored declarations of a block, in order.
func declarations(body string) [][2]string {
	var out [][2]string
	p := &cssParser{s: body}
	for p.i < len(p.s) {
		d, _ := p.until(";")
		k, v, ok := strings.Cut(d, ":")
		if !ok {
			continue
		}
		k = strings.ToLower(strings.TrimSpace(k))
		if !honored[k] {
			continue
		}
		v = strings.TrimSpace(v)
		v = strings.TrimSpace(strings.TrimSuffix(strings.ToLower(v), "!important"))
		if v == "" || strings.ContainsAny(v, ";{}\"") {
			continue
		}
		out = append(out, [2]string{k, v})
	}
	return out
}

// parseSelector reads a compound selector and its specificity.
func parseSelector(s string) (sel selector, spec int, ok bool) {
	if s == "" {
		return sel, 0, false
	}
	i := 0
	ident := func() string {
		start := i
		for i < len(s) {
			c := s[i]
			if c == '-' || c == '_' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80 {
				i++
				continue
			}
			if c == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			break
		}
		return strings.ReplaceAll(s[start:i], `\`, "")
	}
	if s[0] == '*' {
		i = 1
	} else if t := ident(); t != "" {
		sel.tag = strings.ToLower(t)
		spec++
	}
	for i < len(s) {
		c := s[i]
		i++
		switch c {
		case '.':
			name := ident()
			if name == "" {
				return sel, 0, false
			}
			sel.classes = append(sel.classes, name)
			spec += 100
		case '#':
			name := ident()
			if name == "" {
				return sel, 0, false
			}
			sel.id = name
			spec += 10000
		case ':':
			if name := strings.ToLower(ident()); name != "root" || sel.tag != "" && sel.tag != "html" {
				return sel, 0, false
			}
			sel.tag = "html"
			spec += 100
		default:
			// combinators, attribute selectors, pseudo-elements
			return sel, 0, false
		}
	}
	return sel, spec, true
}

// joinURL resolves a relative URL against the URL of the sheet it is in.
func joinURL(base, ref string) string {
	if strings.Contains(ref, ":") || strings.HasPrefix(ref, "/") {
		return ref
	}
	if i := strings.LastIndexByte(base, '/'); i >= 0 {
		return base[:i+1] + ref
	}
	return ref
}

// styler applies the honored declarations of a chapter's rules to its
// elements.
type styler struct {
	rules   []cssRule
	byClass map[string][]int
	byTag   map[string][]int
	byID    map[string][]int
	any     []int
}

func newStyler(rules []cssRule) *styler {
	s := &styler{rules: rules, byClass: map[string][]int{}, byTag: map[string][]int{}, byID: map[string][]int{}}
	for i, r := range rules {
		switch {
		case r.sel.id != "":
			s.byID[r.sel.id] = append(s.byID[r.sel.id], i)
		case len(r.sel.classes) > 0:
			s.byClass[r.sel.classes[0]] = append(s.byClass[r.sel.classes[0]], i)
		case r.sel.tag != "":
			s.byTag[r.sel.tag] = append(s.byTag[r.sel.tag], i)
		default:
			s.any = append(s.any, i)
		}
	}
	return s
}

// apply puts the declarations that the rules give an element in its style
// attribute, before its own declarations (which win).
func (s *styler) apply(n *html.Node) {
	if len(s.rules) == 0 || n.Namespace != "" {
		return
	}
	classes := strings.Fields(attrVal(n, "class"))
	var cand []int
	cand = append(cand, s.any...)
	cand = append(cand, s.byTag[n.Data]...)
	if id := attrVal(n, "id"); id != "" {
		cand = append(cand, s.byID[id]...)
	}
	for _, c := range classes {
		cand = append(cand, s.byClass[c]...)
	}
	if len(cand) == 0 {
		return
	}
	slices.Sort(cand)
	cand = slices.Compact(cand)
	type winner struct {
		v           string
		spec, order int
	}
	won := map[string]winner{}
	var keys []string
	for _, i := range cand {
		r := &s.rules[i]
		if !r.sel.matches(n, classes) {
			continue
		}
		for _, d := range r.decls {
			if pictureOnly[d[0]] && n.Data != "img" {
				continue
			}
			if d[0] == "display" && d[1] != "none" {
				continue
			}
			w, ok := won[d[0]]
			if ok && (w.spec > r.spec || w.spec == r.spec && w.order > r.order) {
				continue
			}
			if !ok {
				keys = append(keys, d[0])
			}
			won[d[0]] = winner{d[1], r.spec, r.order}
		}
	}
	if len(keys) == 0 {
		return
	}
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(won[k].v)
		b.WriteString("; ")
	}
	for i, a := range n.Attr {
		if a.Namespace == "" && a.Key == "style" {
			n.Attr[i].Val = b.String() + a.Val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: "style", Val: strings.TrimSuffix(b.String(), " ")})
}

// inlineStyle reads a declaration of an element's style attribute (the
// last one, as CSS does).
func inlineStyle(n *html.Node, prop string) string {
	v := ""
	for _, d := range strings.Split(attrVal(n, "style"), ";") {
		k, val, ok := strings.Cut(d, ":")
		if ok && strings.EqualFold(strings.TrimSpace(k), prop) {
			v = strings.ToLower(strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(val), "!important")))
		}
	}
	return v
}

// Writing modes of a chapter.
const (
	modeUnset = iota
	modeHorizontal
	modeVertical
	modeVerticalLR // vertical-lr: laid out horizontally
)

// writingMode reads the writing mode of an element's style.
func writingMode(n *html.Node) int {
	mode := modeUnset
	for _, p := range []string{"-epub-writing-mode", "-webkit-writing-mode", "writing-mode"} {
		switch inlineStyle(n, p) {
		case "vertical-rl", "tb-rl", "tb":
			mode = modeVertical
		case "horizontal-tb", "lr-tb", "lr", "rl-tb", "rl":
			mode = modeHorizontal
		case "vertical-lr", "tb-lr":
			mode = modeVerticalLR
		}
	}
	return mode
}

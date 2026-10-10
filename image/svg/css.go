package svg

import (
	"slices"
	"sort"
	"strings"
)

// cssRule is a rule of a style sheet with one selector.
type cssRule struct {
	sel   []cssCompound // descendant combinators between them
	spec  int
	order int
	decls [][2]string
}

// cssCompound is a simple selector sequence: tag, #id and .classes.
type cssCompound struct {
	tag, id string
	classes []string
}

// parseCSS reads the rules of a style sheet that previews need: compound
// selectors of type, id and class, and descendant combinators. At-rules
// are skipped.
func parseCSS(src string) []cssRule {
	if strings.Contains(src, "/*") {
		// without the comments, in one pass
		var b strings.Builder
		for {
			i := strings.Index(src, "/*")
			if i < 0 {
				b.WriteString(src)
				break
			}
			b.WriteString(src[:i])
			j := strings.Index(src[i+2:], "*/")
			if j < 0 {
				break
			}
			src = src[i+2+j+2:]
		}
		src = b.String()
	}
	var rules []cssRule
	order := 0
	for {
		open := strings.IndexByte(src, '{')
		if open < 0 {
			break
		}
		prelude := strings.TrimSpace(src[:open])
		// find the matching close brace
		depth, end := 0, -1
		for k := open; k < len(src); k++ {
			switch src[k] {
			case '{':
				depth++
			case '}':
				depth--
			}
			if depth == 0 {
				end = k
				break
			}
		}
		if end < 0 {
			break
		}
		body := src[open+1 : end]
		src = src[end+1:]
		if strings.HasPrefix(prelude, "@") {
			continue
		}
		decls := parseDecls(body)
		for _, s := range strings.Split(prelude, ",") {
			sel, spec, ok := parseSelector(strings.TrimSpace(s))
			if !ok {
				continue
			}
			rules = append(rules, cssRule{sel: sel, spec: spec, order: order, decls: decls})
			order++
		}
	}
	return rules
}

func parseSelector(s string) ([]cssCompound, int, bool) {
	if s == "" {
		return nil, 0, false
	}
	var out []cssCompound
	spec := 0
	for _, part := range strings.Fields(strings.ReplaceAll(s, ">", " ")) {
		var c cssCompound
		i := 0
		read := func() string {
			j := i
			for j < len(part) && part[j] != '.' && part[j] != '#' && part[j] != '[' && part[j] != ':' {
				j++
			}
			v := part[i:j]
			i = j
			return v
		}
		if part[0] != '.' && part[0] != '#' {
			c.tag = read()
			if c.tag != "*" {
				spec++
			} else {
				c.tag = ""
			}
		}
		for i < len(part) {
			switch part[i] {
			case '.':
				i++
				c.classes = append(c.classes, read())
				spec += 100
			case '#':
				i++
				c.id = read()
				spec += 10000
			default:
				return nil, 0, false // attribute selectors and pseudo-classes
			}
		}
		out = append(out, c)
	}
	return out, spec, true
}

func (c *cssCompound) matches(n *Node) bool {
	if c.tag != "" && c.tag != n.Name || c.id != "" && n.Attr["id"] != c.id {
		return false
	}
	if len(c.classes) > 0 {
		have := strings.Fields(n.Attr["class"])
		for _, want := range c.classes {
			found := false
			for _, h := range have {
				if h == want {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

func (r *cssRule) matches(n *Node) bool {
	last := len(r.sel) - 1
	if !r.sel[last].matches(n) {
		return false
	}
	k := last - 1
	for p := n.Parent; p != nil && k >= 0; p = p.Parent {
		if r.sel[k].matches(p) {
			k--
		}
	}
	return k < 0
}

// styleSheet is the rules of the style sheets of a document, by what the
// last part of their selector asks for: an element is matched against the
// rules that ask for its id, one of its classes or its name, and those that
// ask for none.
type styleSheet struct {
	byID, byClass, byTag map[string][]*cssRule
	others               []*cssRule
	// matches counts the rules matched against elements, and most is the
	// most there may be (svgMatches)
	matches, most int
}

func newStyleSheet(rules []cssRule, most int) *styleSheet {
	s := &styleSheet{byID: map[string][]*cssRule{}, byClass: map[string][]*cssRule{}, byTag: map[string][]*cssRule{}, most: most}
	for i := range rules {
		r := &rules[i]
		switch last := r.sel[len(r.sel)-1]; {
		case last.id != "":
			s.byID[last.id] = append(s.byID[last.id], r)
		case len(last.classes) > 0:
			s.byClass[last.classes[0]] = append(s.byClass[last.classes[0]], r)
		case last.tag != "":
			s.byTag[last.tag] = append(s.byTag[last.tag], r)
		default:
			s.others = append(s.others, r)
		}
	}
	return s
}

// apply stores on each element the declarations of the rules that match
// it, the more specific and later ones winning.
func (s *styleSheet) apply(n *Node) {
	var hit []*cssRule
	match := func(rules []*cssRule) {
		for _, r := range rules {
			if s.matches++; s.matches > s.most {
				return
			}
			if r.matches(n) {
				hit = append(hit, r)
			}
		}
	}
	if id := n.Attr["id"]; id != "" {
		match(s.byID[id])
	}
	if class := n.Attr["class"]; class != "" {
		classes := strings.Fields(class)
		for i, c := range classes {
			if !slices.Contains(classes[:i], c) {
				match(s.byClass[c])
			}
		}
	}
	match(s.byTag[n.Name])
	match(s.others)
	if len(hit) > 0 {
		sort.SliceStable(hit, func(i, j int) bool {
			if hit[i].spec != hit[j].spec {
				return hit[i].spec < hit[j].spec
			}
			return hit[i].order < hit[j].order
		})
		n.Sheet = map[string]string{}
		for _, r := range hit {
			for _, d := range r.decls {
				n.Sheet[d[0]] = d[1]
			}
		}
	}
	for _, c := range n.Children {
		if c.Name != TextNode && s.matches <= s.most {
			s.apply(c)
		}
	}
}

// parseDecls reads CSS declarations ("fill: red; stroke: blue").
func parseDecls(s string) [][2]string {
	var out [][2]string
	for _, d := range strings.Split(s, ";") {
		k, v, ok := strings.Cut(d, ":")
		if !ok {
			continue
		}
		v = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(v), "!important"))
		out = append(out, [2]string{strings.ToLower(strings.TrimSpace(k)), strings.TrimSpace(v)})
	}
	return out
}

// Declared returns the value an element declares for a property: in its
// style attribute, then in the style sheets of the document, then as a
// presentation attribute. Properties are not inherited here: an element
// that declares none has the value of its parent, or the initial one.
func (n *Node) Declared(name string) (string, bool) {
	if st := n.Attr["style"]; st != "" {
		for _, d := range parseDecls(st) {
			if d[0] == name {
				return d[1], true
			}
		}
	}
	if v, ok := n.Sheet[name]; ok {
		return v, true
	}
	v, ok := n.Attr[name]
	return strings.TrimSpace(v), ok
}

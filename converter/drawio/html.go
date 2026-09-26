package drawio

import (
	"html"
	"strings"
)

// A tolerant parser for the HTML of draw.io labels (html=1): the markup
// its rich text editor writes (b, i, u, font, span and div with inline
// styles, br, lists, headings, links, tables) and whatever users paste.

// hnode is an element or a text node.
type hnode struct {
	tag   string // lowercase element name; "" for text nodes
	attrs map[string]string
	kids  []*hnode
	text  string // text nodes: the decoded text
}

var voidTags = map[string]bool{
	"br": true, "hr": true, "img": true, "input": true, "meta": true, "link": true, "wbr": true,
	"area": true, "base": true, "col": true, "embed": true, "source": true, "track": true, "param": true,
}

// autoClose lists, for elements whose end tag may be left out, the
// elements whose start closes them.
var autoClose = map[string][]string{
	"p":  {"p", "div", "ul", "ol", "table", "h1", "h2", "h3", "h4", "h5", "h6", "blockquote", "pre", "hr", "li"},
	"li": {"li"},
	"td": {"td", "th", "tr"},
	"th": {"td", "th", "tr"},
	"tr": {"tr"},
	"dt": {"dt", "dd"},
	"dd": {"dt", "dd"},
}

// parseHTML parses a label into a tree under a root element.
func parseHTML(s string) *hnode {
	root := &hnode{tag: "#root"}
	stack := []*hnode{root}
	top := func() *hnode { return stack[len(stack)-1] }
	appendText := func(t string) {
		if t == "" {
			return
		}
		p := top()
		if n := len(p.kids); n > 0 && p.kids[n-1].tag == "" {
			p.kids[n-1].text += t
			return
		}
		p.kids = append(p.kids, &hnode{text: t})
	}
	for len(s) > 0 {
		i := strings.IndexByte(s, '<')
		if i < 0 {
			appendText(html.UnescapeString(s))
			break
		}
		if i > 0 {
			appendText(html.UnescapeString(s[:i]))
			s = s[i:]
		}
		// comments, doctype, CDATA
		if strings.HasPrefix(s, "<!--") {
			end := strings.Index(s, "-->")
			if end < 0 {
				break
			}
			s = s[end+3:]
			continue
		}
		if strings.HasPrefix(s, "<!") || strings.HasPrefix(s, "<?") {
			end := strings.IndexByte(s, '>')
			if end < 0 {
				break
			}
			s = s[end+1:]
			continue
		}
		end := tagEnd(s)
		if end < 0 || len(s) < 2 || !(isLetter(s[1]) || s[1] == '/' && len(s) > 2 && isLetter(s[2])) {
			// a lone "<": text
			appendText("<")
			s = s[1:]
			continue
		}
		tag := s[1:end]
		s = s[end+1:]
		if strings.HasPrefix(tag, "/") {
			name := strings.ToLower(strings.TrimSpace(tag[1:]))
			if sp := strings.IndexAny(name, " \t\r\n"); sp >= 0 {
				name = name[:sp]
			}
			for k := len(stack) - 1; k > 0; k-- {
				if stack[k].tag == name {
					stack = stack[:k]
					break
				}
			}
			continue
		}
		selfClosing := strings.HasSuffix(tag, "/")
		tag = strings.TrimSuffix(tag, "/")
		name, attrs := parseTag(tag)
		if name == "" {
			continue
		}
		if closers, ok := autoCloseBy[name]; ok {
			for k := len(stack) - 1; k > 0; k-- {
				if closers[stack[k].tag] {
					stack = stack[:k]
					break
				}
				if blocksAutoClose[stack[k].tag] {
					break
				}
			}
		}
		n := &hnode{tag: name, attrs: attrs}
		top().kids = append(top().kids, n)
		if !selfClosing && !voidTags[name] {
			stack = append(stack, n)
			if name == "script" || name == "style" {
				// raw text up to the end tag, dropped
				endTag := "</" + name
				k := strings.Index(strings.ToLower(s), endTag)
				if k < 0 {
					s = ""
				} else {
					s = s[k:]
				}
			}
		}
	}
	return root
}

// autoCloseBy maps an element to the open elements its start tag closes.
var autoCloseBy = func() map[string]map[string]bool {
	m := map[string]map[string]bool{}
	for open, closers := range autoClose {
		for _, c := range closers {
			if m[c] == nil {
				m[c] = map[string]bool{}
			}
			m[c][open] = true
		}
	}
	return m
}()

// blocksAutoClose stops the search for an element to auto-close.
var blocksAutoClose = map[string]bool{"ul": true, "ol": true, "table": true, "div": true, "blockquote": true, "td": true, "th": true}

func isLetter(b byte) bool { return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' }

// tagEnd returns the index of the ">" that ends the tag starting s,
// skipping quoted attribute values; -1 when there is none.
func tagEnd(s string) int {
	var q byte
	for i := 1; i < len(s); i++ {
		c := s[i]
		switch {
		case q != 0:
			if c == q {
				q = 0
			}
		case c == '"' || c == '\'':
			q = c
		case c == '>':
			return i
		case c == '<':
			return -1
		}
	}
	return -1
}

// parseTag splits a tag's content into its lowercase name and attributes.
func parseTag(s string) (string, map[string]string) {
	i := 0
	for i < len(s) && s[i] != ' ' && s[i] != '\t' && s[i] != '\n' && s[i] != '\r' {
		i++
	}
	name := strings.ToLower(s[:i])
	attrs := map[string]string{}
	s = s[i:]
	for {
		s = strings.TrimLeft(s, " \t\r\n")
		if s == "" {
			break
		}
		j := 0
		for j < len(s) && s[j] != '=' && s[j] != ' ' && s[j] != '\t' && s[j] != '\n' && s[j] != '\r' {
			j++
		}
		key := strings.ToLower(s[:j])
		s = strings.TrimLeft(s[j:], " \t\r\n")
		val := ""
		if strings.HasPrefix(s, "=") {
			s = strings.TrimLeft(s[1:], " \t\r\n")
			if s != "" && (s[0] == '"' || s[0] == '\'') {
				q := s[0]
				k := strings.IndexByte(s[1:], q)
				if k < 0 {
					val, s = s[1:], ""
				} else {
					val, s = s[1:1+k], s[2+k:]
				}
			} else {
				k := 0
				for k < len(s) && s[k] != ' ' && s[k] != '\t' && s[k] != '\n' && s[k] != '\r' {
					k++
				}
				val, s = s[:k], s[k:]
			}
		}
		if key != "" {
			attrs[key] = html.UnescapeString(val)
		}
		if j == 0 && val == "" {
			s = s[1:] // skip a stray character
		}
	}
	return name, attrs
}

// parseCSS reads a style attribute into property → value (lowercase
// property names, values trimmed; !important dropped).
func parseCSS(s string) map[string]string {
	out := map[string]string{}
	for _, decl := range splitDecls(s) {
		k, v, ok := strings.Cut(decl, ":")
		if !ok {
			continue
		}
		k = strings.ToLower(strings.TrimSpace(k))
		v = strings.TrimSpace(v)
		v = strings.TrimSpace(strings.TrimSuffix(v, "!important"))
		if k != "" {
			out[k] = v
		}
	}
	return out
}

// splitDecls splits CSS declarations at semicolons outside parentheses and quotes.
func splitDecls(s string) []string {
	var out []string
	depth, start := 0, 0
	var q byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case q != 0:
			if c == q {
				q = 0
			}
		case c == '"' || c == '\'':
			q = c
		case c == '(':
			depth++
		case c == ')':
			depth--
		case c == ';' && depth == 0:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

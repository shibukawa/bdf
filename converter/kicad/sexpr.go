package kicad

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// node is a list of the S-expression files: (name arg ...). An argument is
// an atom or a string (item.list nil) or a list.
type node struct {
	name  string
	items []item
}

type item struct {
	s      string // atom or string
	quoted bool   // a string: never a keyword or a number
	list   *node
}

// maxDepth bounds the nesting of lists (KiCad files nest about a dozen
// deep).
const maxDepth = 256

var errSyntax = errors.New("kicad: malformed S-expression")

// parse reads the first list of a file.
func parse(b []byte) (*node, error) {
	p := &parser{b: b}
	p.space()
	if p.i >= len(p.b) || p.b[p.i] != '(' {
		return nil, errSyntax
	}
	return p.list(0)
}

type parser struct {
	b []byte
	i int
}

func (p *parser) space() {
	for p.i < len(p.b) {
		switch p.b[p.i] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			p.i++
		default:
			return
		}
	}
}

// list reads a list; p.i is at its "(".
func (p *parser) list(depth int) (*node, error) {
	if depth > maxDepth {
		return nil, fmt.Errorf("%w: nested too deep", errSyntax)
	}
	p.i++ // (
	n := &node{}
	first := true
	for {
		p.space()
		if p.i >= len(p.b) {
			return nil, fmt.Errorf("%w: unexpected end", errSyntax)
		}
		switch c := p.b[p.i]; c {
		case ')':
			p.i++
			return n, nil
		case '(':
			child, err := p.list(depth + 1)
			if err != nil {
				return nil, err
			}
			if first {
				// a list without a name, as in ((a b)): keep it as an item
				first = false
			}
			n.items = append(n.items, item{list: child})
		case '"':
			s, err := p.str()
			if err != nil {
				return nil, err
			}
			if first {
				n.name, first = s, false
				continue
			}
			n.items = append(n.items, item{s: s, quoted: true})
		default:
			start := p.i
			for p.i < len(p.b) {
				c := p.b[p.i]
				if c == '(' || c == ')' || c == '"' || c == ' ' || c == '\t' || c == '\n' || c == '\r' {
					break
				}
				p.i++
			}
			s := string(p.b[start:p.i])
			if first {
				n.name, first = s, false
				continue
			}
			n.items = append(n.items, item{s: s})
		}
	}
}

// str reads a quoted string; p.i is at its opening quote.
func (p *parser) str() (string, error) {
	p.i++
	var sb strings.Builder
	for p.i < len(p.b) {
		c := p.b[p.i]
		switch c {
		case '"':
			p.i++
			return sb.String(), nil
		case '\\':
			if p.i+1 >= len(p.b) {
				return "", fmt.Errorf("%w: unterminated string", errSyntax)
			}
			p.i += 2
			switch e := p.b[p.i-1]; e {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			default:
				sb.WriteByte(e)
			}
		default:
			sb.WriteByte(c)
			p.i++
		}
	}
	return "", fmt.Errorf("%w: unterminated string", errSyntax)
}

// child returns the first child list with the name, or nil.
func (n *node) child(name string) *node {
	if n == nil {
		return nil
	}
	for _, it := range n.items {
		if it.list != nil && it.list.name == name {
			return it.list
		}
	}
	return nil
}

// children returns the child lists with the name.
func (n *node) children(name string) []*node {
	if n == nil {
		return nil
	}
	var out []*node
	for _, it := range n.items {
		if it.list != nil && it.list.name == name {
			out = append(out, it.list)
		}
	}
	return out
}

// lists returns every child list.
func (n *node) lists() []*node {
	if n == nil {
		return nil
	}
	var out []*node
	for _, it := range n.items {
		if it.list != nil {
			out = append(out, it.list)
		}
	}
	return out
}

// arg returns the i-th argument that is not a list ("" when there is
// none).
func (n *node) arg(i int) string {
	if n == nil {
		return ""
	}
	for _, it := range n.items {
		if it.list != nil {
			continue
		}
		if i == 0 {
			return it.s
		}
		i--
	}
	return ""
}

// num returns the i-th argument as a number (0 when it is not one).
func (n *node) num(i int) float64 {
	v, err := strconv.ParseFloat(n.arg(i), 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return v
}

// int returns the i-th argument as an integer.
func (n *node) int(i int) int {
	v := n.num(i)
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0
	}
	return int(v)
}

// str returns the first argument of the named child ("" when absent).
func (n *node) str(name string) string { return n.child(name).arg(0) }

// numOf returns the first argument of the named child as a number, or def
// when the child is absent.
func (n *node) numOf(name string, def float64) float64 {
	c := n.child(name)
	if c == nil || c.arg(0) == "" {
		return def
	}
	return c.num(0)
}

// flag reports a boolean property, written either as a bare keyword among
// the arguments (hide, KiCad 6 and 7) or as a child list (hide yes,
// KiCad 8 and later; a child without a value is true).
func (n *node) flag(name string) bool {
	if n == nil {
		return false
	}
	for _, it := range n.items {
		if it.list == nil && !it.quoted && it.s == name {
			return true
		}
		if it.list != nil && it.list.name == name {
			v := it.list.arg(0)
			return v == "" || v == "yes" || v == "true"
		}
	}
	return false
}

// flagOr is flag with a default for a child that is absent.
func (n *node) flagOr(name string, def bool) bool {
	if n == nil {
		return def
	}
	for _, it := range n.items {
		if it.list == nil && !it.quoted && it.s == name {
			return true
		}
		if it.list != nil && it.list.name == name {
			v := it.list.arg(0)
			return v == "" || v == "yes" || v == "true"
		}
	}
	return def
}

// hasAtom reports whether an unquoted argument is the keyword.
func (n *node) hasAtom(s string) bool {
	if n == nil {
		return false
	}
	for _, it := range n.items {
		if it.list == nil && !it.quoted && it.s == s {
			return true
		}
	}
	return false
}

// xy reads (at x y [angle]) or (xy x y) style positions.
func (n *node) xy() (x, y float64) {
	if n == nil {
		return 0, 0
	}
	return n.num(0), n.num(1)
}

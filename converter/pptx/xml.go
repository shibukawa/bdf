package pptx

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

// node is a generic XML element. DrawingML is deeply optional and
// inherited from several places, so the converter walks element trees
// rather than decoding into structs.
type node struct {
	Space string // namespace URI
	Name  string // local name
	Attr  []xml.Attr
	Kids  []*node
	Text  string // character data directly inside the element
}

const (
	nsRel = "relationships" // suffix of the relationships namespaces (transitional and strict)
	nsMC  = "markup-compatibility/2006"
)

// parseXML reads a document into a node tree. mc:AlternateContent is
// replaced by its mc:Fallback (or its first mc:Choice when there is no
// fallback): fallbacks carry the pictures and plain shapes that stand in
// for features this converter does not implement.
func parseXML(data []byte) (*node, error) {
	d := xml.NewDecoder(bytes.NewReader(data))
	d.Strict = false
	var stack []*node
	var root *node
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			n := &node{Space: t.Name.Space, Name: t.Name.Local, Attr: t.Attr}
			if len(stack) > 0 {
				p := stack[len(stack)-1]
				p.Kids = append(p.Kids, n)
			} else if root == nil {
				root = n
			}
			stack = append(stack, n)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += string(t)
			}
		}
	}
	if root == nil {
		return nil, io.ErrUnexpectedEOF
	}
	root.resolveAlternates()
	return root, nil
}

func (n *node) resolveAlternates() {
	var kids []*node
	changed := false
	for _, k := range n.Kids {
		if k.Name == "AlternateContent" && strings.HasSuffix(k.Space, nsMC) {
			changed = true
			var pick *node
			for _, c := range k.Kids {
				if c.Name == "Fallback" {
					pick = c
				}
			}
			if pick == nil {
				for _, c := range k.Kids {
					if c.Name == "Choice" {
						pick = c
						break
					}
				}
			}
			if pick != nil {
				for _, c := range pick.Kids {
					c.resolveAlternates()
					kids = append(kids, c)
				}
			}
			continue
		}
		k.resolveAlternates()
		kids = append(kids, k)
	}
	if changed {
		n.Kids = kids
	}
}

// text returns the character data of an element ("" for nil).
func (n *node) text() string {
	if n == nil {
		return ""
	}
	return n.Text
}

// kids returns the child elements (nil for nil).
func (n *node) kids() []*node {
	if n == nil {
		return nil
	}
	return n.Kids
}

// attrs returns the attributes (nil for nil).
func (n *node) attrs() []xml.Attr {
	if n == nil {
		return nil
	}
	return n.Attr
}

// child returns the first child element with the local name, or nil.
func (n *node) child(name string) *node {
	if n == nil {
		return nil
	}
	for _, k := range n.Kids {
		if k.Name == name {
			return k
		}
	}
	return nil
}

// children returns the child elements with the local name.
func (n *node) children(name string) []*node {
	if n == nil {
		return nil
	}
	var out []*node
	for _, k := range n.Kids {
		if k.Name == name {
			out = append(out, k)
		}
	}
	return out
}

// path follows child elements by local name.
func (n *node) path(names ...string) *node {
	for _, name := range names {
		n = n.child(name)
		if n == nil {
			return nil
		}
	}
	return n
}

// attr returns an unqualified (or any-namespace) attribute by local name.
func (n *node) attr(name string) (string, bool) {
	if n == nil {
		return "", false
	}
	for _, a := range n.Attr {
		if a.Name.Local == name && !strings.HasSuffix(a.Name.Space, nsRel) {
			return a.Value, true
		}
	}
	return "", false
}

func (n *node) attrStr(name, def string) string {
	if v, ok := n.attr(name); ok {
		return v
	}
	return def
}

func (n *node) attrInt(name string, def int64) int64 {
	if v, ok := n.attr(name); ok {
		if i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return int64(f)
		}
	}
	return def
}

func (n *node) attrFloat(name string, def float64) float64 {
	if v, ok := n.attr(name); ok {
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return f
		}
	}
	return def
}

// attrBool reads xsd:boolean ("1", "true", "on" in older producers).
func (n *node) attrBool(name string, def bool) bool {
	if v, ok := n.attr(name); ok {
		switch strings.TrimSpace(v) {
		case "1", "true", "on":
			return true
		case "0", "false", "off":
			return false
		}
	}
	return def
}

// rid returns a relationship attribute (r:id, r:embed, r:link …).
func (n *node) rid(name string) string {
	if n == nil {
		return ""
	}
	for _, a := range n.Attr {
		if a.Name.Local == name && strings.HasSuffix(a.Name.Space, nsRel) {
			return a.Value
		}
	}
	return ""
}

// pct parses an ST_Percentage value: 1000ths of a percent ("50000") or,
// in strict documents, "50%". It returns a fraction.
func pct(n *node, name string, def float64) float64 {
	v, ok := n.attr(name)
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	if strings.HasSuffix(v, "%") {
		if f, err := strconv.ParseFloat(strings.TrimSuffix(v, "%"), 64); err == nil {
			return f / 100
		}
		return def
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f / 100000
	}
	return def
}

// emuPt converts EMU to points.
const emuPerPt = 12700

func emu(v int64) float64 { return float64(v) / emuPerPt }

func (n *node) emuAttr(name string, def float64) float64 {
	v, ok := n.attr(name)
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	// Strict documents may use universal measures ("1in", "2.5cm", "12pt").
	for _, u := range []struct {
		suffix string
		pt     float64
	}{{"pt", 1}, {"in", 72}, {"cm", 72 / 2.54}, {"mm", 72 / 25.4}, {"pc", 12}, {"pi", 12}} {
		if strings.HasSuffix(v, u.suffix) {
			if f, err := strconv.ParseFloat(strings.TrimSuffix(v, u.suffix), 64); err == nil {
				return f * u.pt
			}
			return def
		}
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f / emuPerPt
	}
	return def
}

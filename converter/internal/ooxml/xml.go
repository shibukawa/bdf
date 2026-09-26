package ooxml

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

// Node is a generic XML element. Office markup (DrawingML above all) is
// deeply optional and inherited from several places, so converters walk
// element trees rather than decoding into structs. The accessors are
// nil-safe: a missing element reads as empty.
type Node struct {
	Space string // namespace URI
	Name  string // local name
	Attrs []xml.Attr
	Kids  []*Node
	Text  string // character data directly inside the element
}

const (
	nsRel = "relationships" // suffix of the relationships namespaces (transitional and strict)
	nsMC  = "markup-compatibility/2006"
)

// Parse reads a document into a node tree. mc:AlternateContent is
// replaced by its mc:Fallback (or its first mc:Choice when there is no
// fallback): fallbacks carry the pictures and plain shapes that stand in
// for features the converters do not implement.
func Parse(data []byte) (*Node, error) {
	d := xml.NewDecoder(bytes.NewReader(data))
	d.Strict = false
	for {
		tok, err := d.Token()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		if err != nil {
			return nil, err
		}
		if start, ok := tok.(xml.StartElement); ok {
			return ReadElement(d, start)
		}
	}
}

// ReadElement reads the element that start opens from d, up to its end,
// into a node tree (with mc:AlternateContent resolved as Parse does). It
// lets a reader stream the bulk of a large part and keep the rest as trees.
func ReadElement(d *xml.Decoder, start xml.StartElement) (*Node, error) {
	root := &Node{Space: start.Name.Space, Name: start.Name.Local, Attrs: start.Attr}
	stack := []*Node{root}
	for len(stack) > 0 {
		tok, err := d.Token()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			n := &Node{Space: t.Name.Space, Name: t.Name.Local, Attrs: t.Attr}
			p := stack[len(stack)-1]
			p.Kids = append(p.Kids, n)
			stack = append(stack, n)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			stack[len(stack)-1].Text += string(t)
		}
	}
	root.resolveAlternates()
	return root, nil
}

func (n *Node) resolveAlternates() {
	var kids []*Node
	changed := false
	for _, k := range n.Kids {
		if k.Name == "AlternateContent" && strings.HasSuffix(k.Space, nsMC) {
			changed = true
			var pick *Node
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

// Content returns the character data of an element ("" for nil).
func (n *Node) Content() string {
	if n == nil {
		return ""
	}
	return n.Text
}

// Elements returns the child elements (nil for nil).
func (n *Node) Elements() []*Node {
	if n == nil {
		return nil
	}
	return n.Kids
}

// Child returns the first child element with the local name, or nil.
func (n *Node) Child(name string) *Node {
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

// Children returns the child elements with the local name.
func (n *Node) Children(name string) []*Node {
	if n == nil {
		return nil
	}
	var out []*Node
	for _, k := range n.Kids {
		if k.Name == name {
			out = append(out, k)
		}
	}
	return out
}

// Path follows child elements by local name.
func (n *Node) Path(names ...string) *Node {
	for _, name := range names {
		n = n.Child(name)
		if n == nil {
			return nil
		}
	}
	return n
}

// Attr returns an unqualified (or any-namespace) attribute by local name.
func (n *Node) Attr(name string) (string, bool) {
	if n == nil {
		return "", false
	}
	for _, a := range n.Attrs {
		if a.Name.Local == name && !strings.HasSuffix(a.Name.Space, nsRel) {
			return a.Value, true
		}
	}
	return "", false
}

// AttrStr returns an attribute (see Attr), or def when it is absent.
func (n *Node) AttrStr(name, def string) string {
	if v, ok := n.Attr(name); ok {
		return v
	}
	return def
}

// AttrInt reads an integer attribute (a decimal one is truncated).
func (n *Node) AttrInt(name string, def int64) int64 {
	if v, ok := n.Attr(name); ok {
		if i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return int64(f)
		}
	}
	return def
}

// AttrFloat reads a decimal attribute.
func (n *Node) AttrFloat(name string, def float64) float64 {
	if v, ok := n.Attr(name); ok {
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return f
		}
	}
	return def
}

// AttrBool reads xsd:boolean ("1", "true", "on" in older producers).
func (n *Node) AttrBool(name string, def bool) bool {
	if v, ok := n.Attr(name); ok {
		switch strings.TrimSpace(v) {
		case "1", "true", "on":
			return true
		case "0", "false", "off":
			return false
		}
	}
	return def
}

// RelID returns a relationship attribute (r:id, r:embed, r:link …).
func (n *Node) RelID(name string) string {
	if n == nil {
		return ""
	}
	for _, a := range n.Attrs {
		if a.Name.Local == name && strings.HasSuffix(a.Name.Space, nsRel) {
			return a.Value
		}
	}
	return ""
}

// AttrPct parses an ST_Percentage value: 1000ths of a percent ("50000") or,
// in strict documents, "50%". It returns a fraction.
func (n *Node) AttrPct(name string, def float64) float64 {
	v, ok := n.Attr(name)
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

// EMUPerPoint is the number of English Metric Units in a point.
const EMUPerPoint = 12700

// AttrEMU reads a DrawingML coordinate (ST_Coordinate: EMU, or a universal
// measure in strict documents) in points.
func (n *Node) AttrEMU(name string, def float64) float64 {
	v, ok := n.Attr(name)
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
		return f / EMUPerPoint
	}
	return def
}

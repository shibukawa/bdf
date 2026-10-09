// Package xmltree reads the XML of Office documents into element trees:
// what the converters of Office formats (converter/internal/ooxml) and the
// reader of Office Math (internal/mathlayout) walk.
package xmltree

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
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
	// runs places the character data of an element with child elements
	// among them, for mixed content (see Segments); nil for others.
	runs []textRun
}

// textRun is character data that comes before child element Kids[before]
// (after the last one when before is len(Kids)).
type textRun struct {
	before int
	text   string
}

// addText appends character data to the element's content.
func (n *Node) addText(s string) {
	n.Text += s
	if k := len(n.runs) - 1; k >= 0 && n.runs[k].before == len(n.Kids) {
		n.runs[k].text += s
	} else {
		n.runs = append(n.runs, textRun{len(n.Kids), s})
	}
}

// Segment is a piece of the content of an element: a child element or
// character data between child elements. Exactly one of Elem and Text is
// set.
type Segment struct {
	Elem *Node
	Text string
}

const (
	nsRel = "relationships" // suffix of the relationships namespaces (transitional and strict)
	nsMC  = "markup-compatibility/2006"
)

// MaxDepth bounds how deeply elements may nest. The tree is never built
// deeper, so every later recursion over it (resolveAlternates here, and the
// converters' walks over nested shapes, groups, tables, cells and text
// boxes) is bounded and cannot overflow the stack; a deeper part is
// rejected as malformed. The deepest document known is a test file of
// Apache POI with 5,000 tables one inside the other, which nest 15,000
// elements deep; the time to lay tables out grows with the square of
// their depth, so the limit is not far above it.
const MaxDepth = 20000

// MaxElements bounds the elements of one tree: a part with more is
// rejected as malformed, as a deeper one is. An element of a tree takes a
// few hundred bytes whatever its markup took (four, for "<a/>"), so a part
// could otherwise ask for memory out of all proportion with its size: a
// gigabyte of markup, which a megabyte of a package inflates to, would
// take tens. A reader that keeps several trees gives them a budget of its
// own (ParseCounting, ReadFrom).
const MaxElements = 8 << 20

// MaxTokenBytes bounds one token of a part read from a stream (NewReader):
// the text of an element, or the value of an attribute, which holds a
// picture in VML's gfxdata (some megabytes of base64). A part parsed from
// memory has no such bound.
const MaxTokenBytes = 256 << 20

// ErrTooManyElements is the error of a tree of more than MaxElements
// elements, or more than the budget its reader gave it.
var ErrTooManyElements = errors.New("ooxml: too many elements")

// The trees are read with xmlro, which allocates nothing for a token: the
// bytes of a name or a value are the reader's until the next token, and the
// tree copies only what it keeps. Names and short values repeat throughout
// a part, so they are interned: one string for every "p" and "val". The
// namespaces are resolved as encoding/xml resolves them, and the trees are
// the same (ReadElementCounting reads one from an encoding/xml decoder for
// a reader that still streams with it; the tests compare them).

// Parse reads a document into a node tree. mc:AlternateContent is
// replaced by its mc:Fallback (or its first mc:Choice when there is no
// fallback): fallbacks carry the pictures and plain shapes that stand in
// for features the converters do not implement.
func Parse(data []byte) (*Node, error) { return ParseChoosing(data, nil) }

// ParseChoosing is Parse that replaces mc:AlternateContent with its first
// mc:Choice whose required namespaces (the prefixes of its Requires
// attribute) all pass supported, and with the fallback when none does.
// Word processing documents need it: they keep their DrawingML shapes in
// a choice and legacy VML in the fallback.
func ParseChoosing(data []byte, supported func(prefix string) bool) (*Node, error) {
	return ParsePicking(data, SupportedChoice(supported))
}

// ParsePicking is Parse that replaces mc:AlternateContent with its first
// mc:Choice that pick accepts, and with the fallback when it accepts none.
func ParsePicking(data []byte, pick func(choice *Node) bool) (*Node, error) {
	budget := MaxElements
	return parseFrom(xmlro.NewBytesReader(data, ReaderOptions()), pick, &budget)
}

// ParsePickingReader parses XML from a reader without retaining the source
// bytes alongside the resulting node tree.
func ParsePickingReader(r io.Reader, pick func(choice *Node) bool) (*Node, error) {
	budget := MaxElements
	return ParseCounting(r, pick, &budget)
}

// ParseCounting is ParsePickingReader for a reader that keeps several
// trees: budget is how many elements they may have together, which the
// tree's elements are taken from (ErrTooManyElements when they run out).
func ParseCounting(r io.Reader, pick func(choice *Node) bool, budget *int) (*Node, error) {
	return parseFrom(NewReader(r), pick, budget)
}

// NewReader returns a reader of the XML tokens of src with the bounds of
// this package (ReaderOptions): what reads a part an element at a time
// streams with it, and hands the elements it keeps to ReadFrom.
func NewReader(src io.Reader) *xmlro.Reader { return xmlro.NewReader(src, ReaderOptions()) }

// ReaderOptions bounds a token reader: a token of MaxTokenBytes at most,
// and nesting a little past MaxDepth, which the trees check themselves. A
// DOCTYPE is read past, as encoding/xml reads it. The buffer starts small,
// as most parts of a package are (the relationships of a part, a layout),
// and grows to the largest token of a part.
func ReaderOptions() xmlro.Options {
	return xmlro.Options{BufferSize: 16 << 10, MaxBufferBytes: MaxTokenBytes, MaxDepth: MaxDepth + 64, AllowDoctype: true}
}

// ParseFrom reads the document a reader is at the start of into a tree, as
// ParseCounting does: for a reader that is reused from one part to the
// next (Reset), with the buffer it grew.
func ParseFrom(r *xmlro.Reader, pick func(choice *Node) bool, budget *int) (*Node, error) {
	return parseFrom(r, pick, budget)
}

// parseFrom reads the root element of a document into a tree.
func parseFrom(r *xmlro.Reader, pick func(choice *Node) bool, budget *int) (*Node, error) {
	for {
		k, err := r.Next()
		if err != nil {
			return nil, err
		}
		switch k {
		case xmlro.StartElement:
			return ReadFrom(r, pick, budget)
		case xmlro.EOF:
			return nil, io.ErrUnexpectedEOF
		}
	}
}

// ReadFrom reads the element the reader is on into a tree, with
// mc:AlternateContent replaced by its first mc:Choice that pick accepts
// (see ParsePicking), and leaves the reader on the element's EndElement.
// budget is how many elements the trees of the reader's caller may have
// together (ErrTooManyElements when they run out). It lets a reader stream
// the bulk of a large part and keep the rest as trees.
func ReadFrom(r *xmlro.Reader, pick func(choice *Node) bool, budget *int) (*Node, error) {
	if r.Kind() != xmlro.StartElement {
		return nil, xmlro.ErrNotStart
	}
	b := builder{r: r, budget: budget}
	root, err := b.start()
	if err != nil {
		return nil, err
	}
	stack := []*Node{root}
	var text []byte
	for len(stack) > 0 {
		k, err := r.Next()
		if err != nil {
			return nil, err
		}
		switch k {
		case xmlro.StartElement:
			if len(stack) >= MaxDepth {
				return nil, fmt.Errorf("ooxml: element nesting deeper than %d", MaxDepth)
			}
			n, err := b.start()
			if err != nil {
				return nil, err
			}
			p := stack[len(stack)-1]
			p.Kids = append(p.Kids, n)
			stack = append(stack, n)
		case xmlro.EndElement:
			if n := stack[len(stack)-1]; len(n.Kids) == 0 {
				n.runs = nil // Text alone says it all
			}
			stack = stack[:len(stack)-1]
		case xmlro.Text, xmlro.CData:
			if k == xmlro.Text {
				text = r.Text().AppendTo(text[:0]) // entities decoded
			} else {
				text = append(text[:0], r.Text()...) // a CDATA section holds none
			}
			stack[len(stack)-1].addText(string(text))
		case xmlro.ProcInst:
			// InDesign writes the characters XML cannot hold (its page
			// number marker U+0018 and the like) as <?ACE 18?>
			if string(r.Name()) == "ACE" {
				if code, err := strconv.ParseUint(strings.TrimSpace(string(r.Text())), 16, 32); err == nil && code < 0x110000 {
					stack[len(stack)-1].addText(string(rune(code)))
				}
			}
		case xmlro.EOF:
			return nil, io.ErrUnexpectedEOF
		}
	}
	root.resolveAlternates(pick)
	return root, nil
}

// builder makes the nodes of one tree.
type builder struct {
	r      *xmlro.Reader
	budget *int
	names  interner
}

// maxInternedValue is the longest attribute value that is interned: the
// values of Office markup are mostly short and repeated (a size, a colour,
// a style id), and a long one is a picture or a formula.
const maxInternedValue = 32

// start makes the node of the StartElement the reader is on.
func (b *builder) start() (*Node, error) {
	if *b.budget <= 0 {
		return nil, ErrTooManyElements
	}
	*b.budget--
	r := b.r
	n := &Node{Name: b.names.str(r.LocalName())}
	if p := r.Prefix(); p != nil {
		if uri, ok := r.LookupNamespace(p); ok {
			n.Space = uri
		} else {
			n.Space = b.names.str(p) // a prefix bound to nothing stays, as encoding/xml keeps it
		}
	} else {
		n.Space = r.Namespace()
	}
	for {
		name, val, ok := r.NextAttr()
		if !ok {
			break
		}
		var a xml.Attr
		if i := bytes.IndexByte(name, ':'); i >= 0 {
			prefix, local := name[:i], name[i+1:]
			switch uri, ok := r.LookupNamespace(prefix); {
			case xmlro.Equal(prefix, "xmlns"):
				a.Name.Space = "xmlns" // a declaration, as encoding/xml reports it
			case ok:
				a.Name.Space = uri
			default:
				a.Name.Space = b.names.str(prefix)
			}
			a.Name.Local = b.names.str(local)
		} else {
			a.Name.Local = b.names.str(name) // no default namespace for attributes
		}
		if len(val) <= maxInternedValue && !val.HasEntities() {
			a.Value = b.names.str(val)
		} else {
			a.Value = val.String()
		}
		n.Attrs = append(n.Attrs, a)
	}
	return n, nil
}

// interner keeps one string for the names and values that repeat in a
// part: an open addressing table that compares bytes to strings without
// making a string of the bytes (which TinyGo would do for a map index).
type interner struct {
	table []string // "" is an empty slot
	n     int
}

func (t *interner) str(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if t.table == nil {
		t.table = make([]string, 64)
	}
	mask := len(t.table) - 1
	for i := int(fnv(b)) & mask; ; i = (i + 1) & mask {
		s := t.table[i]
		if s == "" {
			if 2*(t.n+1) > len(t.table) {
				t.grow()
				return t.str(b)
			}
			s = string(b)
			t.table[i] = s
			t.n++
			return s
		}
		if xmlro.Equal(b, s) {
			return s
		}
	}
}

func (t *interner) grow() {
	old := t.table
	t.table = make([]string, 2*len(old))
	t.n = 0
	mask := len(t.table) - 1
	for _, s := range old {
		if s == "" {
			continue
		}
		i := int(fnv(s)) & mask
		for t.table[i] != "" {
			i = (i + 1) & mask
		}
		t.table[i] = s
		t.n++
	}
}

// fnv is the FNV-1a hash of a name.
func fnv[S ~string | ~[]byte](s S) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

// ReadElement reads the element that start opens from d, up to its end,
// into a node tree (with mc:AlternateContent resolved as Parse does), for
// a reader that streams a part with encoding/xml; one that streams with
// NewReader uses ReadFrom.
func ReadElement(d *xml.Decoder, start xml.StartElement) (*Node, error) {
	return ReadElementPicking(d, start, nil)
}

// SupportedChoice picks the choices whose required namespaces (the
// prefixes of their Requires attribute) all pass supported.
func SupportedChoice(supported func(prefix string) bool) func(choice *Node) bool {
	if supported == nil {
		return nil
	}
	return func(c *Node) bool { return requiresAll(c.AttrStr("Requires", ""), supported) }
}

// MathChoice picks the choices that hold Office Math in DrawingML text
// (a14:m, PowerPoint's and Excel's equations), whose fallback is a picture
// or the formula as plain text, and no other choice.
func MathChoice(c *Node) bool {
	if strings.TrimSpace(c.AttrStr("Requires", "")) != "a14" {
		return false
	}
	var find func(n *Node) bool
	find = func(n *Node) bool {
		for _, k := range n.Kids {
			if k.Name == "m" && k.Space == NSA14 || find(k) {
				return true
			}
		}
		return false
	}
	return find(c)
}

// NSA14 is the namespace of the Office 2010 DrawingML extensions.
const NSA14 = "http://schemas.microsoft.com/office/drawing/2010/main"

// ReadElementPicking is ReadElement that replaces mc:AlternateContent with
// its first mc:Choice that pick accepts.
func ReadElementPicking(d *xml.Decoder, start xml.StartElement, pick func(choice *Node) bool) (*Node, error) {
	budget := MaxElements
	return ReadElementCounting(d, start, pick, &budget)
}

// ReadElementCounting is ReadElementPicking with a budget of elements, as
// ReadFrom has.
func ReadElementCounting(d *xml.Decoder, start xml.StartElement, pick func(choice *Node) bool, budget *int) (*Node, error) {
	if *budget <= 0 {
		return nil, ErrTooManyElements
	}
	*budget--
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
			if len(stack) >= MaxDepth {
				return nil, fmt.Errorf("ooxml: element nesting deeper than %d", MaxDepth)
			}
			if *budget <= 0 {
				return nil, ErrTooManyElements
			}
			*budget--
			n := &Node{Space: t.Name.Space, Name: t.Name.Local, Attrs: t.Attr}
			p := stack[len(stack)-1]
			p.Kids = append(p.Kids, n)
			stack = append(stack, n)
		case xml.EndElement:
			if n := stack[len(stack)-1]; len(n.Kids) == 0 {
				n.runs = nil // Text alone says it all
			}
			stack = stack[:len(stack)-1]
		case xml.CharData:
			stack[len(stack)-1].addText(string(t))
		case xml.ProcInst:
			// InDesign writes the characters XML cannot hold (its page
			// number marker U+0018 and the like) as <?ACE 18?>
			if t.Target == "ACE" {
				if code, err := strconv.ParseUint(strings.TrimSpace(string(t.Inst)), 16, 32); err == nil && code < 0x110000 {
					stack[len(stack)-1].addText(string(rune(code)))
				}
			}
		}
	}
	root.resolveAlternates(pick)
	return root, nil
}

func (n *Node) resolveAlternates(pick func(choice *Node) bool) {
	// Most OOXML nodes have no alternate child. Avoid allocating a replacement
	// child list and index map for every paragraph, run and text node.
	hasAlternate := false
	for _, k := range n.Kids {
		if k.Name == "AlternateContent" && strings.HasSuffix(k.Space, nsMC) {
			hasAlternate = true
			break
		}
	}
	if !hasAlternate {
		for _, k := range n.Kids {
			k.resolveAlternates(pick)
		}
		return
	}
	var kids []*Node
	changed := false
	// start[i] is where old child i starts among the new children, so that
	// character data stays in place
	start := make([]int, len(n.Kids)+1)
	for i, k := range n.Kids {
		start[i] = len(kids)
		if k.Name == "AlternateContent" && strings.HasSuffix(k.Space, nsMC) {
			changed = true
			var picked *Node
			if pick != nil {
				for _, c := range k.Kids {
					if c.Name == "Choice" && pick(c) {
						picked = c
						break
					}
				}
			}
			for _, c := range k.Kids {
				if picked == nil && c.Name == "Fallback" {
					picked = c
				}
			}
			if picked == nil {
				for _, c := range k.Kids {
					if c.Name == "Choice" {
						picked = c
						break
					}
				}
			}
			if picked != nil {
				for _, c := range picked.Kids {
					c.resolveAlternates(pick)
					kids = append(kids, c)
				}
			}
			continue
		}
		k.resolveAlternates(pick)
		kids = append(kids, k)
	}
	if changed {
		start[len(n.Kids)] = len(kids)
		for i := range n.runs {
			n.runs[i].before = start[n.runs[i].before]
		}
		n.Kids = kids
	}
}

func requiresAll(requires string, supported func(prefix string) bool) bool {
	fields := strings.Fields(requires)
	for _, p := range fields {
		if !supported(p) {
			return false
		}
	}
	return len(fields) > 0
}

// Segments returns the content of an element in document order: its child
// elements and the character data around them (mixed content, as in
// "<a>one<b/>two</a>"). It returns nil for nil.
func (n *Node) Segments() []Segment {
	if n == nil {
		return nil
	}
	if len(n.Kids) == 0 {
		if n.Text == "" {
			return nil
		}
		return []Segment{{Text: n.Text}}
	}
	out := make([]Segment, 0, len(n.Kids)+len(n.runs))
	r := 0
	for i, k := range n.Kids {
		for ; r < len(n.runs) && n.runs[r].before <= i; r++ {
			out = append(out, Segment{Text: n.runs[r].text})
		}
		out = append(out, Segment{Elem: k})
	}
	for ; r < len(n.runs); r++ {
		out = append(out, Segment{Text: n.runs[r].text})
	}
	return out
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

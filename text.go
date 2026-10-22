package bdf

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MARK kinds (docs/spec.md §7.8).
const (
	MarkParagraph byte = 0
	MarkLine      byte = 1
	MarkCell      byte = 2
	MarkBox       byte = 3
	MarkAltText   byte = 4
	MarkWrap      byte = 5
	MarkHeading   byte = 6  // payload: level "1"-"6"
	MarkList      byte = 7  // opens a list, closed by MarkEnd
	MarkListItem  byte = 8  // starts an item of the innermost list
	MarkTable     byte = 9  // opens a table, closed by MarkEnd
	MarkFigure    byte = 10 // payload: alternative text; opens a figure, closed by MarkEnd
	MarkEnd       byte = 11 // closes the innermost list, table or figure
	MarkLang      byte = 12 // payload: BCP 47 language of the runs that follow ("" = document default)
	MarkBarcode   byte = 13 // payload: JSON of the symbol the next drawing draws (image/barcode.Mark)
)

// Separators between consecutive text runs (docs/spec.md §7.9).
const (
	SepNone  byte = 0
	SepSpace byte = 1
	SepBreak byte = 2
)

// TextRun is a text drawing instruction (or ALT_TEXT mark) located in the
// coordinate space of the object being walked.
type TextRun struct {
	Text    string
	X, Y    float32 // anchor after applying the transform
	Advance float32
	Size    float32
	Font    *Font
	Align   byte
	Matrix  [6]float32
	Sep     byte // separator from the previous run, from MARKs or position heuristics
	Ordinal int  // position in the extraction order of the top-level object
	AltText bool
	Node    int    // index of the structure node the run belongs to (TextContent.Nodes)
	Lang    string // BCP 47 language set by a LANG mark; "" for the document's default
}

// NodeKind is the kind of a structure node (docs/spec.md §7.8). Leaves
// (paragraph, heading, caption, a cell outside a table) hold runs;
// containers (list, item, table, a cell inside a table, figure) hold nodes.
type NodeKind byte

// Kinds of structure nodes.
const (
	NodeParagraph NodeKind = iota
	NodeHeading
	NodeCaption
	NodeList
	NodeItem
	NodeTable
	NodeCell
	NodeFigure
)

var nodeKindNames = [...]string{"paragraph", "heading", "caption", "list", "item", "table", "cell", "figure"}

func (k NodeKind) String() string {
	if int(k) < len(nodeKindNames) {
		return nodeKindNames[k]
	}
	return fmt.Sprintf("node(%d)", byte(k))
}

// TextNode is a node of the structure of an object: what its MARKs say the
// text is (docs/spec.md §7.8). Nodes point at their parent by index.
type TextNode struct {
	Kind   NodeKind
	Parent int // index of the enclosing node, -1 at the top level
	// Level is the level of a heading, 1 to 6, or 0 when the payload is
	// not a level.
	Level int
	// Cell is the place of a cell: relative to its table inside a TABLE,
	// that of the sheet outside one. Nil when the payload is not a
	// reference.
	Cell *CellRef
	// Alt is the alternative text of a figure ("" when it is decorative).
	Alt string
	// Code says a paragraph is a line of code (the word "code" of its
	// payload; docs/spec.md §7.8).
	Code bool
	// Quote is how deep in quotations a paragraph or a list is: the count
	// of the word "quote" in its payload.
	Quote int
	// List is the kind of the markers of a list's items: "bullet",
	// "ordered", or "" when the payload does not say.
	List string
}

// words reads the payload of a PARAGRAPH, BOX or LIST mark into a node:
// the words it knows (docs/spec.md §7.8), the others ignored.
func (n *TextNode) words(payload string) {
	for _, w := range strings.Fields(payload) {
		switch w {
		case "code":
			n.Code = true
		case "quote":
			n.Quote++
		case "bullet", "ordered":
			n.List = w
		}
	}
}

// CellRef is the payload of a CELL mark: a reference or a range, of a
// header cell or not ("B2", "A4:B4 col").
type CellRef struct {
	Row, Col   int    // of the top-left cell, from 0
	Rows, Cols int    // the cells spanned, 1 or more
	Scope      string // "col" (a column header), "row" (a row header) or ""
}

// ParseCellRef reads the payload of a CELL mark.
func ParseCellRef(payload string) (CellRef, bool) {
	s := strings.TrimSpace(payload)
	scope := ""
	if i := strings.LastIndexAny(s, " \t"); i >= 0 {
		switch tail := s[i+1:]; tail {
		case "col", "row":
			scope, s = tail, strings.TrimSpace(s[:i])
		default:
			return CellRef{}, false
		}
	}
	one := func(s string) (row, col int, ok bool) {
		i := 0
		for i < len(s) && (s[i] >= 'A' && s[i] <= 'Z' || s[i] >= 'a' && s[i] <= 'z') {
			col = col*26 + int(s[i]|0x20-'a') + 1
			if i++; i > 4 {
				return 0, 0, false
			}
		}
		if i == 0 || i == len(s) || len(s)-i > 8 {
			return 0, 0, false
		}
		n, err := strconv.Atoi(s[i:])
		if err != nil || n < 1 {
			return 0, 0, false
		}
		return n - 1, col - 1, true
	}
	a, b, isRange := strings.Cut(s, ":")
	r0, c0, ok := one(a)
	if !ok {
		return CellRef{}, false
	}
	out := CellRef{Row: r0, Col: c0, Rows: 1, Cols: 1, Scope: scope}
	if isRange {
		r1, c1, ok := one(b)
		if !ok {
			return CellRef{}, false
		}
		out.Rows, out.Cols = max(1, r1-r0+1), max(1, c1-c0+1)
	}
	return out, true
}

// TextLink is a LINK: the bounding box of its rectangle in the extraction
// space, and where it came in the text.
type TextLink struct {
	X, Y, W, H float32
	URL        string
	After      int // index of the last run before the LINK (-1 before the first)
	Node       int // the leaf, else the container, open at the LINK (-1 at the top level)
}

// TextImage is an IMAGE or IMAGE_SUB: the part it draws, the bounding box
// it is drawn in, and where it came in the text.
type TextImage struct {
	Hash       Hash
	X, Y, W, H float32
	After      int // index of the last run before the image (-1 before the first)
	Node       int // the leaf, else the container, open at the image (-1 at the top level)
	Before     int // index of the first node made after the image (len(Nodes) at the image)
}

// TextContent is the text of an object with its structure: the runs, the
// nodes they belong to, the links and the images, in walk order.
type TextContent struct {
	Runs   []TextRun
	Nodes  []TextNode
	Links  []TextLink
	Images []TextImage
}

type matrix [6]float32

func mul(m, n matrix) matrix {
	return matrix{
		m[0]*n[0] + m[2]*n[1],
		m[1]*n[0] + m[3]*n[1],
		m[0]*n[2] + m[2]*n[3],
		m[1]*n[2] + m[3]*n[3],
		m[0]*n[4] + m[2]*n[5] + m[4],
		m[1]*n[4] + m[3]*n[5] + m[5],
	}
}

type textState struct {
	m     matrix
	font  *Font
	size  float32
	align byte
}

// Limits of walking the objects an object draws with USE. A document decides
// how deep objects draw objects and how often, so a small one could otherwise
// keep a reader busy without end (ten objects that each draw the next one ten
// times draw the last one 10^10 times).
const (
	// MaxUseDepth bounds how deep objects draw objects.
	MaxUseDepth = 64
	// MaxReusedInstructions bounds the instructions read from objects that
	// were walked before. What is read once is bounded by the size of the
	// document; this bounds what a document adds by drawing it again.
	MaxReusedInstructions = 1 << 27
	// maxTextRuns bounds the runs of a walk, which are kept in memory.
	maxTextRuns = 1 << 22
)

// useLimits counts a walk of objects against the limits.
type useLimits struct {
	seen   map[Hash]struct{}
	reused int
	err    error // the limit that was passed; the walk stops
}

// enter reports whether the instructions of a child object count as reused:
// whether it was walked before.
func (l *useLimits) enter(h Hash) bool {
	if _, ok := l.seen[h]; ok {
		return true
	}
	if l.seen == nil {
		l.seen = map[Hash]struct{}{}
	}
	l.seen[h] = struct{}{}
	return false
}

func (l *useLimits) fail(msg string) {
	if l.err == nil {
		l.err = &FormatError{Msg: msg}
	}
}

type textExtractor struct {
	runs    []TextRun
	resolve func(Hash) *ObjectPart
	pending byte // separator requested by a MARK for the next run
	hasMark bool
	alt     *string // ALT_TEXT waiting for the drawing op it describes
	limits  *useLimits
	depth   int // of the object being walked: 0 is the top-level one

	// the structure (docs/spec.md §7.8), which runs on through USE in
	// walk order
	nodes    []TextNode
	open     []int     // open containers, outermost first
	leaf     int       // the leaf the next run joins (-1: start one)
	reserved *TextNode // the leaf a leaf MARK reserved for the next run
	lang     string
	links    []TextLink
	images   []TextImage
}

func (t *textExtractor) top() int {
	if n := len(t.open); n > 0 {
		return t.open[n-1]
	}
	return -1
}

func (t *textExtractor) topKind() (NodeKind, bool) {
	if i := t.top(); i >= 0 {
		return t.nodes[i].Kind, true
	}
	return 0, false
}

func (t *textExtractor) topIs(k NodeKind) bool {
	got, ok := t.topKind()
	return ok && got == k
}

// here is the node that what comes next is in: the leaf, else the
// innermost container, else -1.
func (t *textExtractor) here() int {
	if t.leaf >= 0 && t.reserved == nil {
		return t.leaf
	}
	return t.top()
}

func (t *textExtractor) add(n TextNode) int {
	n.Parent = t.top()
	t.nodes = append(t.nodes, n)
	return len(t.nodes) - 1
}

// leafNode is the leaf of the next run: the one the last leaf MARK
// reserved, or an implicit one (text right in a list is an item of it,
// right in a table its caption).
func (t *textExtractor) leafNode() int {
	if t.leaf >= 0 && t.reserved == nil {
		return t.leaf
	}
	r := t.reserved
	t.reserved = nil
	if t.topIs(NodeList) {
		t.openNode(TextNode{Kind: NodeItem})
	}
	switch {
	case t.topIs(NodeTable):
		t.leaf = t.add(TextNode{Kind: NodeCaption})
	case r != nil:
		t.leaf = t.add(*r)
	default:
		t.leaf = t.add(TextNode{Kind: NodeParagraph})
	}
	return t.leaf
}

func (t *textExtractor) openNode(n TextNode) {
	t.open = append(t.open, t.add(n))
	t.leaf, t.reserved = -1, nil
}

func (t *textExtractor) closeTop() {
	if n := len(t.open); n > 0 {
		t.open = t.open[:n-1]
	}
	t.leaf, t.reserved = -1, nil
}

// structure applies a structure MARK (docs/spec.md §7.8): a leaf MARK
// reserves the kind of the next run's leaf, a container MARK opens or
// closes a container.
func (t *textExtractor) structure(kind byte, payload string) {
	switch kind {
	case MarkParagraph, MarkBox:
		n := &TextNode{Kind: NodeParagraph}
		n.words(payload)
		t.reserved = n
	case MarkHeading:
		n := &TextNode{Kind: NodeHeading}
		if f := strings.Fields(payload); len(f) > 0 {
			if l, err := strconv.Atoi(f[0]); err == nil && l >= 1 && l <= 6 {
				n.Level = l
			}
		}
		t.reserved = n
	case MarkList:
		n := TextNode{Kind: NodeList}
		n.words(payload)
		t.openNode(n)
	case MarkListItem:
		if t.topIs(NodeItem) {
			t.closeTop()
		}
		if t.topIs(NodeList) {
			t.openNode(TextNode{Kind: NodeItem})
		} else {
			t.reserved = &TextNode{Kind: NodeParagraph}
		}
	case MarkTable:
		t.openNode(TextNode{Kind: NodeTable})
	case MarkCell:
		var ref *CellRef
		if c, ok := ParseCellRef(payload); ok {
			ref = &c
		}
		if t.topIs(NodeCell) {
			t.closeTop()
		}
		if t.topIs(NodeTable) {
			t.openNode(TextNode{Kind: NodeCell, Cell: ref})
		} else {
			t.reserved = &TextNode{Kind: NodeCell, Cell: ref}
		}
	case MarkFigure:
		t.openNode(TextNode{Kind: NodeFigure, Alt: payload})
	case MarkEnd:
		if t.topIs(NodeItem) || t.topIs(NodeCell) {
			t.closeTop()
		}
		if len(t.open) > 0 {
			t.closeTop()
		}
	}
}

// box is the bounding box of a rectangle moved by m.
func box(m matrix, x, y, w, h float32) (bx, by, bw, bh float32) {
	x0, y0 := float32(math.Inf(1)), float32(math.Inf(1))
	x1, y1 := float32(math.Inf(-1)), float32(math.Inf(-1))
	for _, p := range [4][2]float32{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}} {
		tx := m[0]*p[0] + m[2]*p[1] + m[4]
		ty := m[1]*p[0] + m[3]*p[1] + m[5]
		x0, y0, x1, y1 = min(x0, tx), min(y0, ty), max(x1, tx), max(y1, ty)
	}
	return x0, y0, x1 - x0, y1 - y0
}

// child returns the child object an instruction draws, nil for one that is
// missing, and whether it was walked before.
func (t *textExtractor) child(o *ObjectPart, in Instr) (*ObjectPart, bool) {
	i := in.Args[0].(uint64)
	if i >= uint64(len(o.Objects)) {
		return nil, false
	}
	c := t.resolve(o.Objects[i])
	if c == nil {
		return nil, false
	}
	return c, t.limits.enter(o.Objects[i])
}

func (t *textExtractor) walk(o *ObjectPart, m matrix, reused bool) error {
	st := textState{m: m, size: 10}
	var stack []textState
	masking := 0 // inside MASK_BEGIN … MASK_END: a soft mask, not content
	// want counts the instructions and picks those that say something
	// about text: the others are most of a page, and are not decoded
	want := func(op byte) bool {
		if t.limits.err != nil {
			return false
		}
		if reused {
			if t.limits.reused++; t.limits.reused > MaxReusedInstructions {
				t.limits.fail("objects are drawn too many times")
				return false
			}
		}
		switch op {
		case OpFillPathAt, OpFillPathRun:
			return t.alt != nil // places of an ALT_TEXT
		}
		return textOps[op]
	}
	err := o.walk(want, func(in Instr) {
		switch in.Op {
		case OpMaskBegin:
			masking++
			return
		case OpMaskEnd:
			masking--
			return
		}
		if masking > 0 {
			return
		}
		switch in.Op {
		case OpSave:
			stack = append(stack, st)
		case OpRestore:
			if n := len(stack); n > 0 {
				st = stack[n-1]
				stack = stack[:n-1]
			}
		case OpTransform:
			st.m = mul(st.m, matrix{f(in, 0), f(in, 1), f(in, 2), f(in, 3), f(in, 4), f(in, 5)})
		case OpTranslate:
			st.m = mul(st.m, matrix{1, 0, 0, 1, f(in, 0), f(in, 1)})
		case OpScale:
			st.m = mul(st.m, matrix{f(in, 0), 0, 0, f(in, 1), 0, 0})
		case OpFont:
			if i := in.Args[0].(uint64); i < uint64(len(o.Fonts)) {
				st.font = &o.Fonts[i]
			}
			st.size = f(in, 1)
		case OpTextStyle:
			st.align = byte(in.Args[0].(uint64))
		case OpMark:
			kind := byte(in.Args[0].(uint64))
			payload := in.Args[1].(string)
			switch kind {
			case MarkLine:
				t.flushAlt(st)
				t.mark(SepSpace)
			case MarkParagraph, MarkCell, MarkBox, MarkHeading, MarkList, MarkListItem, MarkTable, MarkFigure, MarkEnd:
				t.flushAlt(st)
				t.mark(SepBreak)
				t.structure(kind, payload)
			case MarkLang:
				t.lang = payload
			case MarkWrap:
				t.flushAlt(st)
				t.mark(SepNone)
			case MarkAltText:
				t.flushAlt(st)
				p := payload
				t.alt = &p
			}
			// LANG and unknown kinds leave the separator to the next MARK or
			// to the position guess.
		case OpFillText, OpStrokeText:
			if t.alt != nil {
				// The drawing op right after ALT_TEXT renders that text.
				t.emit(*t.alt, f(in, 1), f(in, 2), f(in, 3), st, true)
				t.alt = nil
				return
			}
			t.emit(in.Args[0].(string), f(in, 1), f(in, 2), f(in, 3), st, false)
		case OpFillPathAt:
			if t.alt != nil {
				t.emit(*t.alt, f(in, 2), f(in, 3), 0, st, true)
				t.alt = nil
			}
		case OpFillPathRun:
			if t.alt != nil {
				if gl := in.Args[1].([]Glyph); len(gl) > 0 {
					t.emit(*t.alt, gl[0].X, gl[0].Y, 0, st, true)
				} else {
					t.emit(*t.alt, 0, 0, 0, st, true)
				}
				t.alt = nil
			}
		case OpLink:
			x, y, w, h := box(st.m, f(in, 0), f(in, 1), f(in, 2), f(in, 3))
			t.links = append(t.links, TextLink{X: x, Y: y, W: w, H: h, URL: in.Args[4].(string), After: len(t.runs) - 1, Node: t.here()})
		case OpImage, OpImageSub:
			i := in.Args[0].(uint64)
			if i >= uint64(len(o.Images)) {
				return
			}
			at := 1 // the destination rectangle
			if in.Op == OpImageSub {
				at = 5
			}
			x, y, w, h := box(st.m, f(in, at), f(in, at+1), f(in, at+2), f(in, at+3))
			t.images = append(t.images, TextImage{Hash: o.Images[i], X: x, Y: y, W: w, H: h, After: len(t.runs) - 1, Node: t.here(), Before: len(t.nodes)})
		case OpUse, OpUseAt:
			if t.alt != nil {
				x, y := float32(0), float32(0)
				if in.Op == OpUseAt {
					x, y = f(in, 1), f(in, 2)
				}
				// The text the child draws spans its bbox along the baseline.
				var advance float32
				if child, _ := t.child(o, in); child != nil {
					x += child.BBox.X
					advance = child.BBox.W
				}
				t.emit(*t.alt, x, y, advance, st, true)
				t.alt = nil
				return
			}
			child, again := t.child(o, in)
			if child == nil {
				return
			}
			if t.depth >= MaxUseDepth {
				t.limits.fail("objects draw objects too deep")
				return
			}
			cm := st.m
			if in.Op == OpUseAt {
				cm = mul(cm, matrix{1, 0, 0, 1, f(in, 1), f(in, 2)})
			}
			t.depth++
			_ = t.walk(child, cm, reused || again)
			t.depth--
		}
	})
	if t.limits.err != nil {
		return t.limits.err
	}
	return err
}

func f(in Instr, i int) float32 { return in.Args[i].(float32) }

// textOps are the instructions the extraction of text reads, besides the
// paths that an ALT_TEXT describes.
var textOps = [256]bool{
	OpSave: true, OpRestore: true, OpTransform: true, OpTranslate: true, OpScale: true,
	OpFont: true, OpTextStyle: true, OpFillText: true, OpStrokeText: true,
	OpUse: true, OpUseAt: true, OpMaskBegin: true, OpMaskEnd: true, OpMark: true,
	OpLink: true, OpImage: true, OpImageSub: true,
}

// flushAlt emits an ALT_TEXT that no drawing op consumed, without a position.
func (t *textExtractor) flushAlt(st textState) {
	if t.alt != nil {
		t.emit(*t.alt, 0, 0, 0, st, true)
		t.alt = nil
	}
}

func (t *textExtractor) mark(sep byte) {
	if !t.hasMark || sep > t.pending {
		t.pending = sep
	}
	t.hasMark = true
}

func (t *textExtractor) emit(text string, x, y, advance float32, st textState, alt bool) {
	m := st.m
	run := TextRun{
		Text: text, Advance: advance, Size: st.size, Font: st.font, Align: st.align, Matrix: m,
		X: m[0]*x + m[2]*y + m[4], Y: m[1]*x + m[3]*y + m[5], Ordinal: len(t.runs), AltText: alt,
	}
	if t.hasMark {
		run.Sep = t.pending
	} else if len(t.runs) > 0 {
		run.Sep = guessSep(t.runs[len(t.runs)-1], run)
	}
	t.pending, t.hasMark = SepNone, false
	if len(t.runs) >= maxTextRuns {
		t.limits.fail("too many text runs")
		return
	}
	run.Node, run.Lang = t.leafNode(), t.lang
	t.runs = append(t.runs, run)
}

// guessSep estimates the separator between runs that carry no MARK information.
func guessSep(prev, cur TextRun) byte {
	size := prev.Size
	if size <= 0 {
		size = 10
	}
	if math.Abs(float64(cur.Y-prev.Y)) > float64(size)*0.5 {
		return SepSpace
	}
	end := prev.X + prev.Advance
	if prev.Advance == 0 || cur.X-end > float64To32(0.2)*size {
		return SepSpace
	}
	return SepNone
}

func float64To32(v float64) float32 { return float32(v) }

// ExtractText returns the text runs of an object and the objects it USEs, in
// walk order. resolve must return decoded child objects (nil skips the child).
//
// The objects may not draw each other deeper than MaxUseDepth, nor read more
// than MaxReusedInstructions from objects that are drawn again.
func ExtractText(o *ObjectPart, resolve func(Hash) *ObjectPart) ([]TextRun, error) {
	return extractText(o, resolve, &useLimits{})
}

// ExtractContent returns the runs of an object together with its structure
// (headings, lists, tables, figures: the MARKs of docs/spec.md §7.8), its
// links and its images, for text layers and for writing the document as
// text. The runs are the same as ExtractText's, and so are the limits.
func ExtractContent(o *ObjectPart, resolve func(Hash) *ObjectPart) (*TextContent, error) {
	t, err := extract(o, resolve, &useLimits{})
	return &TextContent{Runs: t.runs, Nodes: t.nodes, Links: t.links, Images: t.images}, err
}

// extractText is ExtractText with limits that count across the top-level
// objects of a document.
func extractText(o *ObjectPart, resolve func(Hash) *ObjectPart, limits *useLimits) ([]TextRun, error) {
	t, err := extract(o, resolve, limits)
	return t.runs, err
}

func extract(o *ObjectPart, resolve func(Hash) *ObjectPart, limits *useLimits) (*textExtractor, error) {
	t := &textExtractor{resolve: resolve, limits: limits, leaf: -1}
	err := t.walk(o, matrix{1, 0, 0, 1, 0, 0}, false)
	t.flushAlt(textState{m: matrix{1, 0, 0, 1, 0, 0}, size: 10})
	return t, err
}

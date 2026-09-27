package parquet

import (
	"strings"
	"unicode/utf8"
)

// A nested field (a group, or a repeated field) is put together from the
// entries of its leaves: the levels of each entry tell which fields on its
// path are there and where lists get a new element. The leaves of a group
// fill the same values, element by element. The value is then written the
// way JSON writes it, with the annotations applied: lists (and repeated
// fields) as arrays, maps and groups as objects, Variants as what they
// hold.

// group is a value of a group: the values of its fields, by position.
type group struct{ fields []any }

// repeated holds the elements of a repeated field.
type repeated struct{ items []any }

// maxItems bounds the elements of a repeated field in a row.
const maxItems = 1 << 20

// assembler puts the entries of a leaf into the values of rows.
type assembler struct {
	path []*node // from the row's field down to the leaf
	last []*repeated
	idx  []int
}

func newAssembler(top, leaf *node) *assembler {
	p := pathTo(top, leaf)
	return &assembler{path: p, last: make([]*repeated, len(p)), idx: make([]int, len(p))}
}

// put puts an entry into the value of its row.
func (a *assembler) put(slot *any, rep, def int, v scalar) bool {
	for i, n := range a.path {
		if n.defLevel > def {
			// this field is null, or an empty list
			if n.rep == repRepeated {
				if _, ok := (*slot).(*repeated); !ok {
					*slot = &repeated{}
				}
			}
			return true
		}
		if n.rep == repRepeated {
			lst, ok := (*slot).(*repeated)
			if !ok {
				lst = &repeated{}
				*slot = lst
			}
			switch {
			case n.repLevel >= rep: // a new element
				if a.last[i] == lst {
					a.idx[i]++
				} else {
					a.last[i], a.idx[i] = lst, 0
				}
			case a.last[i] != lst:
				a.last[i], a.idx[i] = lst, max(len(lst.items)-1, 0)
			}
			j := a.idx[i]
			if j >= maxItems {
				return false
			}
			for len(lst.items) <= j {
				lst.items = append(lst.items, nil)
			}
			slot = &lst.items[j]
		}
		if n.leaf {
			*slot = v
			return true
		}
		g, ok := (*slot).(*group)
		if !ok {
			g = &group{fields: make([]any, len(n.children))}
			*slot = g
		}
		slot = &g.fields[a.path[i+1].index]
	}
	return true
}

// maxCellText bounds the text of a cell, as Excel does.
const maxCellText = 32767

// maxNesting bounds the nesting of the values written.
const maxNesting = 128

// writer writes nested values as JSON writes them.
type writer struct {
	b     strings.Builder
	fmts  []format // of the leaves, by column
	depth int
}

func (w *writer) full() bool { return w.b.Len() > maxCellText }

func (w *writer) raw(s string) { w.b.WriteString(s) }

// str writes a JSON string.
func (w *writer) str(s string) {
	w.b.WriteByte('"')
	for i := 0; i < len(s); {
		c := s[i]
		switch {
		case c == '"' || c == '\\':
			w.b.WriteByte('\\')
			w.b.WriteByte(c)
		case c == '\n':
			w.b.WriteString(`\n`)
		case c == '\r':
			w.b.WriteString(`\r`)
		case c == '\t':
			w.b.WriteString(`\t`)
		case c < 0x20 || c == 0x7f:
			w.b.WriteString(`\u00`)
			w.b.WriteByte("0123456789abcdef"[c>>4])
			w.b.WriteByte("0123456789abcdef"[c&15])
		case c < utf8.RuneSelf:
			w.b.WriteByte(c)
		default:
			r, n := utf8.DecodeRuneInString(s[i:])
			w.b.WriteRune(r)
			i += n
			continue
		}
		i++
	}
	w.b.WriteByte('"')
}

// scalar writes a leaf's value within a nested value.
func (w *writer) scalar(n *node, v scalar) {
	f := w.fmts[n.col]
	if f.kind == fNull {
		w.raw("null")
		return
	}
	s := f.text(v, -1)
	if f.quoted() {
		w.str(s)
	} else {
		w.raw(s)
	}
}

// value writes the value of a field; inst tells that raw is an element of
// a repeated field rather than all of them.
func (w *writer) value(n *node, raw any, inst bool) {
	if w.full() {
		return
	}
	if w.depth++; w.depth > maxNesting {
		w.raw("…")
		w.depth--
		return
	}
	defer func() { w.depth-- }()
	if n.rep == repRepeated && !inst {
		lst, _ := raw.(*repeated)
		w.raw("[")
		if lst != nil {
			for i, it := range lst.items {
				if i > 0 {
					w.raw(", ")
				}
				if w.full() {
					break
				}
				w.value(n, it, true)
			}
		}
		w.raw("]")
		return
	}
	if raw == nil {
		w.raw("null")
		return
	}
	if n.leaf {
		v, _ := raw.(scalar)
		w.scalar(n, v)
		return
	}
	g, _ := raw.(*group)
	if g == nil {
		w.raw("null")
		return
	}
	switch {
	case n.isVariant():
		w.variantGroup(n, g)
	case n.isList():
		w.list(n, g)
	case n.isMap():
		w.mapValue(n, g)
	default:
		w.raw("{")
		for i, c := range n.children {
			if i > 0 {
				w.raw(", ")
			}
			if w.full() {
				break
			}
			w.str(c.name)
			w.raw(": ")
			w.value(c, g.fields[i], false)
		}
		w.raw("}")
	}
}

func (w *writer) list(n *node, g *group) {
	el, inner := n.listElement()
	lst, _ := g.fields[0].(*repeated)
	w.raw("[")
	if lst != nil {
		for i, it := range lst.items {
			if i > 0 {
				w.raw(", ")
			}
			if w.full() {
				break
			}
			if !inner {
				w.value(el, it, true)
				continue
			}
			var v any
			if ig, _ := it.(*group); ig != nil {
				v = ig.fields[0]
			}
			w.value(el, v, false)
		}
	}
	w.raw("]")
}

func (w *writer) mapValue(n *node, g *group) {
	kv := n.children[0]
	lst, _ := g.fields[0].(*repeated)
	w.raw("{")
	if lst != nil {
		for i, it := range lst.items {
			if i > 0 {
				w.raw(", ")
			}
			if w.full() {
				break
			}
			ig, _ := it.(*group)
			if ig == nil {
				w.raw("null: null")
				continue
			}
			w.value(kv.children[0], ig.fields[0], false)
			w.raw(": ")
			w.value(kv.children[1], ig.fields[1], false)
		}
	}
	w.raw("}")
}

// cellText cuts a cell's text to the length a cell holds.
func cellText(s string) string {
	if len(s) <= maxCellText {
		return s
	}
	cut := maxCellText
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut] + "…"
}

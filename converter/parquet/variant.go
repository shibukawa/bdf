package parquet

import (
	"encoding/binary"
	"errors"
	"math"
	"slices"
	"strconv"
)

// A Variant (VariantEncoding.md in parquet-format) is a value of any type
// in two binaries: the metadata, a dictionary of the names of the fields
// of its objects, and the value, a header byte (its basic type: a
// primitive, a short string, an object or an array) and its data. A
// shredded Variant (VariantShredding.md) stores parts of the value as
// typed columns: typed_value holds the value when it is of the column's
// type, a group of fields for an object (each a value and a typed_value),
// or a list of such elements for an array; value holds the rest.

var errVariant = errors.New("malformed Variant")

// variantMeta is the dictionary of the field names of a Variant.
type variantMeta struct {
	b       []byte
	offSize int
	n       int
	offsets int // where the offsets start
	strings int // where the strings start
}

func parseVariantMeta(b []byte) (*variantMeta, error) {
	if len(b) < 1 || b[0]&0x0f != 1 {
		return nil, errVariant
	}
	m := &variantMeta{b: b, offSize: int(b[0]>>6) + 1}
	if len(b) < 1+m.offSize {
		return nil, errVariant
	}
	m.n = int(leUint(b[1:], m.offSize))
	m.offsets = 1 + m.offSize
	if uint64(m.n+1)*uint64(m.offSize) > uint64(len(b)-m.offsets) {
		return nil, errVariant
	}
	m.strings = m.offsets + (m.n+1)*m.offSize
	return m, nil
}

// name returns the field name of an id.
func (m *variantMeta) name(id int) (string, bool) {
	if id < 0 || id >= m.n {
		return "", false
	}
	start := leUint(m.b[m.offsets+id*m.offSize:], m.offSize)
	end := leUint(m.b[m.offsets+(id+1)*m.offSize:], m.offSize)
	if start > end || end > uint64(len(m.b)-m.strings) {
		return "", false
	}
	return string(m.b[m.strings+int(start) : m.strings+int(end)]), true
}

// leUint reads an unsigned little-endian integer of n bytes (n <= 4).
func leUint(b []byte, n int) uint64 {
	var v uint64
	for i := range n {
		v |= uint64(b[i]) << (8 * i)
	}
	return v
}

// variantField is a field of a Variant object.
type variantField struct {
	name  string
	value []byte
}

// variantItems reads the header of an object or an array: the offsets of
// its elements, and for objects their field ids.
func variantItems(b []byte) (ids []int, vals [][]byte, err error) {
	h := b[0]
	basic, vh := h&3, h>>2
	offSize, idSize, large := int(vh&3)+1, 0, false
	if basic == 2 {
		idSize, large = int(vh>>2&3)+1, vh>>4&1 != 0
	} else {
		large = vh>>2&1 != 0
	}
	p := 1
	var n int
	if large {
		if len(b) < p+4 {
			return nil, nil, errVariant
		}
		n = int(binary.LittleEndian.Uint32(b[p:]))
		p += 4
	} else {
		if len(b) < p+1 {
			return nil, nil, errVariant
		}
		n = int(b[p])
		p++
	}
	if n < 0 || uint64(n)*uint64(idSize)+uint64(n+1)*uint64(offSize) > uint64(len(b)-p) {
		return nil, nil, errVariant
	}
	if idSize > 0 {
		ids = make([]int, n)
		for i := range n {
			ids[i] = int(leUint(b[p+i*idSize:], idSize))
		}
		p += n * idSize
	}
	offs := p
	data := p + (n+1)*offSize
	total := leUint(b[offs+n*offSize:], offSize)
	if total > uint64(len(b)-data) {
		return nil, nil, errVariant
	}
	fields := b[data : data+int(total)]
	vals = make([][]byte, n)
	for i := range n {
		start := leUint(b[offs+i*offSize:], offSize)
		if start >= total {
			return nil, nil, errVariant
		}
		end := total
		if basic == 3 {
			end = leUint(b[offs+(i+1)*offSize:], offSize)
			if end < start || end > total {
				return nil, nil, errVariant
			}
		}
		vals[i] = fields[start:end]
	}
	return ids, vals, nil
}

// variantObject returns the fields of a Variant object.
func variantObject(m *variantMeta, b []byte) ([]variantField, error) {
	if len(b) == 0 || b[0]&3 != 2 {
		return nil, errVariant
	}
	ids, vals, err := variantItems(b)
	if err != nil {
		return nil, err
	}
	out := make([]variantField, len(ids))
	for i, id := range ids {
		name, ok := m.name(id)
		if !ok {
			return nil, errVariant
		}
		out[i] = variantField{name, vals[i]}
	}
	return out, nil
}

// variant writes a Variant value.
func (w *writer) variant(m *variantMeta, b []byte) {
	if w.full() {
		return
	}
	if w.depth++; w.depth > maxNesting {
		w.raw("…")
		w.depth--
		return
	}
	defer func() { w.depth-- }()
	if len(b) == 0 {
		w.raw("null")
		return
	}
	h := b[0]
	switch vh := int(h >> 2); h & 3 {
	case 0:
		w.variantPrimitive(vh, b[1:])
	case 1:
		if vh > len(b)-1 {
			w.raw("null")
			return
		}
		w.str(validText(b[1 : 1+vh]))
	case 2:
		fields, err := variantObject(m, b)
		if err != nil {
			w.raw("null")
			return
		}
		w.raw("{")
		for i, f := range fields {
			if i > 0 {
				w.raw(", ")
			}
			w.str(f.name)
			w.raw(": ")
			w.variant(m, f.value)
		}
		w.raw("}")
	case 3:
		_, vals, err := variantItems(b)
		if err != nil {
			w.raw("null")
			return
		}
		w.raw("[")
		for i, v := range vals {
			if i > 0 {
				w.raw(", ")
			}
			if w.full() {
				break
			}
			w.variant(m, v)
		}
		w.raw("]")
	}
}

// variantSizes are the sizes of the data of the primitive types of fixed
// size.
var variantSizes = [...]int{0, 0, 0, 1, 2, 4, 8, 8, 5, 9, 17, 4, 8, 8, 4, -1, -1, 8, 8, 8, 16}

func (w *writer) variantPrimitive(typ int, b []byte) {
	if typ >= len(variantSizes) {
		w.raw("null")
		return
	}
	size := variantSizes[typ]
	if size < 0 {
		if len(b) < 4 || uint64(binary.LittleEndian.Uint32(b)) > uint64(len(b)-4) {
			w.raw("null")
			return
		}
		b = b[4 : 4+binary.LittleEndian.Uint32(b)]
	} else if len(b) < size {
		w.raw("null")
		return
	}
	le := binary.LittleEndian
	switch typ {
	case 0:
		w.raw("null")
	case 1:
		w.raw("true")
	case 2:
		w.raw("false")
	case 3:
		w.raw(strconv.Itoa(int(int8(b[0]))))
	case 4:
		w.raw(strconv.Itoa(int(int16(le.Uint16(b)))))
	case 5:
		w.raw(strconv.Itoa(int(int32(le.Uint32(b)))))
	case 6:
		w.raw(strconv.FormatInt(int64(le.Uint64(b)), 10))
	case 7:
		w.raw(floatText(math.Float64frombits(le.Uint64(b)), 64))
	case 8, 9, 10:
		// the scale, then the unscaled value, little endian
		be := slices.Clone(b[1:size])
		slices.Reverse(be)
		w.raw(decimalText(bigEndianInt(be), int32(b[0])))
	case 11:
		w.str(dateText(int64(int32(le.Uint32(b)))))
	case 12, 13, 18, 19:
		f := format{kind: fTimestamp, unit: unitMicros}
		if typ >= 18 {
			f.unit = unitNanos
		}
		w.str(f.text(scalar{u: le.Uint64(b)}, -1))
	case 14:
		w.raw(floatText(float64(math.Float32frombits(le.Uint32(b))), 32))
	case 15:
		w.str(hexText(b))
	case 16:
		w.str(validText(b))
	case 17:
		f := format{kind: fTime, unit: unitMicros}
		w.str(f.text(scalar{u: le.Uint64(b)}, -1))
	case 20:
		w.str(uuidText(b[:16]))
	}
}

// variantGroup writes a Variant column's value: its metadata, and the
// value, shredded or not.
func (w *writer) variantGroup(n *node, g *group) {
	metaN, valueN, typedN := n.variantFields()
	mv, _ := g.fields[metaN.index].(scalar)
	m, err := parseVariantMeta(mv.b)
	if err != nil {
		w.raw("null")
		return
	}
	w.shredded(m, g, valueN, typedN)
}

// shredded writes a Variant value from the value and typed_value fields
// of a group (a Variant column, a field of a shredded object or an
// element of a shredded array).
func (w *writer) shredded(m *variantMeta, g *group, valueN, typedN *node) {
	var value []byte
	hasValue := false
	if valueN != nil {
		if v, ok := g.fields[valueN.index].(scalar); ok {
			value, hasValue = v.b, true
		}
	}
	var typed any
	if typedN != nil {
		typed = g.fields[typedN.index]
	}
	switch {
	case typed == nil && hasValue:
		w.variant(m, value)
	case typed == nil:
		w.raw("null")
	case typedN.leaf:
		w.value(typedN, typed, false)
	case typedN.isList():
		tg, _ := typed.(*group)
		el, inner := typedN.listElement()
		var items []any
		if lst, _ := tg.fields[0].(*repeated); lst != nil {
			items = lst.items
		}
		w.raw("[")
		for i, it := range items {
			if i > 0 {
				w.raw(", ")
			}
			if w.full() {
				break
			}
			eg, _ := it.(*group)
			en := el
			if inner && eg != nil {
				// the element group of a three-level list
				eg, _ = eg.fields[0].(*group)
			}
			if eg == nil || en.leaf {
				w.raw("null")
				continue
			}
			_, ev, et := en.variantFields()
			w.shredded(m, eg, ev, et)
		}
		w.raw("]")
	default:
		// an object: the shredded fields that are there, and the fields of
		// value (a partially shredded object)
		tg, _ := typed.(*group)
		type entry struct {
			name  string
			write func()
		}
		var entries []entry
		for i, c := range typedN.children {
			fg, _ := tg.fields[i].(*group)
			if c.leaf || fg == nil {
				continue
			}
			_, fv, ft := c.variantFields()
			if (fv == nil || fg.fields[fv.index] == nil) && (ft == nil || fg.fields[ft.index] == nil) {
				continue // missing
			}
			entries = append(entries, entry{c.name, func() { w.shredded(m, fg, fv, ft) }})
		}
		if hasValue {
			// the fields that are shredded are not taken from value (the
			// names must differ; where they do not, the shredded field wins)
			if fields, err := variantObject(m, value); err == nil {
				for _, f := range fields {
					if !slices.ContainsFunc(typedN.children, func(c *node) bool { return c.name == f.name }) {
						entries = append(entries, entry{f.name, func() { w.variant(m, f.value) }})
					}
				}
			}
		}
		slices.SortStableFunc(entries, func(a, b entry) int {
			switch {
			case a.name < b.name:
				return -1
			case a.name > b.name:
				return 1
			}
			return 0
		})
		w.raw("{")
		for i, e := range entries {
			if i > 0 {
				w.raw(", ")
			}
			if w.full() {
				break
			}
			w.str(e.name)
			w.raw(": ")
			e.write()
		}
		w.raw("}")
	}
}

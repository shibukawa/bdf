package parquet

import (
	"encoding/binary"
	"math"
)

// The file metadata and the page headers of a Parquet file are Thrift
// structures in the compact protocol: a struct is a list of fields, each
// a header byte (the field id as a delta from the previous one, and the
// type) and its value, ended by a zero byte. Integers are zigzag varints,
// binaries and strings a varint length and their bytes, and a list a
// header byte (its length, and the type of its elements) and its
// elements.

// Types of the compact protocol.
const (
	tStop   = 0
	tTrue   = 1
	tFalse  = 2
	tByte   = 3
	tI16    = 4
	tI32    = 5
	tI64    = 6
	tDouble = 7
	tBinary = 8
	tList   = 9
	tSet    = 10
	tMap    = 11
	tStruct = 12
	tUUID   = 13
)

const maxThriftDepth = 64

// tenc writes the compact protocol.
type tenc struct {
	b    []byte
	last []int16
}

func newTenc() *tenc { return &tenc{last: []int16{0}} }

func (e *tenc) field(id int16, typ byte) {
	if d := id - e.last[len(e.last)-1]; d > 0 && d <= 15 {
		e.b = append(e.b, byte(d)<<4|typ)
	} else {
		e.b = append(e.b, typ)
		e.b = binary.AppendVarint(e.b, int64(id))
	}
	e.last[len(e.last)-1] = id
}

func (e *tenc) i32(id int16, v int32) { e.field(id, tI32); e.b = binary.AppendVarint(e.b, int64(v)) }
func (e *tenc) i64(id int16, v int64) { e.field(id, tI64); e.b = binary.AppendVarint(e.b, v) }
func (e *tenc) bin(id int16, v string) {
	e.field(id, tBinary)
	e.b = binary.AppendUvarint(e.b, uint64(len(v)))
	e.b = append(e.b, v...)
}

func (e *tenc) boolean(id int16, v bool) {
	if v {
		e.field(id, tTrue)
	} else {
		e.field(id, tFalse)
	}
}

func (e *tenc) begin(id int16) { e.field(id, tStruct); e.last = append(e.last, 0) }
func (e *tenc) elem()          { e.last = append(e.last, 0) } // a struct in a list
func (e *tenc) end()           { e.b = append(e.b, tStop); e.last = e.last[:len(e.last)-1] }

func (e *tenc) list(id int16, typ byte, n int) {
	e.field(id, tList)
	e.listHeader(typ, n)
}

func (e *tenc) listHeader(typ byte, n int) {
	if n < 15 {
		e.b = append(e.b, byte(n)<<4|typ)
	} else {
		e.b = append(e.b, 0xf0|typ)
		e.b = binary.AppendUvarint(e.b, uint64(n))
	}
}

func (e *tenc) bytes() []byte { return append(e.b, tStop) }

// tdec reads the compact protocol. The first error stops it: later reads
// return zero values, and err tells.
type tdec struct {
	b     []byte
	pos   int
	err   error
	depth int
}

func (t *tdec) fail() {
	if t.err == nil {
		t.err = errCorrupt
	}
	t.pos = len(t.b)
}

func (t *tdec) byte() byte {
	if t.pos >= len(t.b) {
		t.fail()
		return 0
	}
	b := t.b[t.pos]
	t.pos++
	return b
}

func (t *tdec) uvarint() uint64 {
	v, n := binary.Uvarint(t.b[t.pos:])
	if n <= 0 {
		t.fail()
		return 0
	}
	t.pos += n
	return v
}

func (t *tdec) varint() int64 {
	v, n := binary.Varint(t.b[t.pos:])
	if n <= 0 {
		t.fail()
		return 0
	}
	t.pos += n
	return v
}

func (t *tdec) i32() int32 {
	v := t.varint()
	if v < math.MinInt32 || v > math.MaxInt32 {
		t.fail()
		return 0
	}
	return int32(v)
}

func (t *tdec) i64() int64 { return t.varint() }

func (t *tdec) binary() []byte {
	n := t.uvarint()
	if n > uint64(len(t.b)-t.pos) {
		t.fail()
		return nil
	}
	b := t.b[t.pos : t.pos+int(n)]
	t.pos += int(n)
	return b
}

func (t *tdec) string() string { return string(t.binary()) }

// fields reads a struct, calling f with the id and the type of each
// field; f reads the value, or skips it.
func (t *tdec) fields(f func(id int16, typ byte)) {
	if t.depth++; t.depth > maxThriftDepth {
		t.fail()
	}
	defer func() { t.depth-- }()
	var id int16
	for t.err == nil {
		h := t.byte()
		if h == tStop {
			return
		}
		typ := h & 0x0f
		if d := h >> 4; d != 0 {
			id += int16(d)
		} else {
			v := t.varint()
			if v < math.MinInt16 || v > math.MaxInt16 {
				t.fail()
				return
			}
			id = int16(v)
		}
		f(id, typ)
	}
}

// list reads the header of a list or a set: the type of its elements and
// their number.
func (t *tdec) list() (typ byte, n int) {
	h := t.byte()
	typ = h & 0x0f
	size := uint64(h >> 4)
	if size == 15 {
		size = t.uvarint()
	}
	if size > uint64(len(t.b)-t.pos) {
		t.fail()
		return 0, 0
	}
	return typ, int(size)
}

// listOf reads a list whose elements are of type want, calling f to read
// each; the elements of a list of another type are skipped.
func (t *tdec) listOf(typ, want byte, f func()) {
	if typ != tList && typ != tSet {
		t.skip(typ)
		return
	}
	et, n := t.list()
	for i := 0; i < n && t.err == nil; i++ {
		if et == want {
			f()
		} else {
			t.skipElem(et)
		}
	}
}

func (t *tdec) skipElem(typ byte) {
	if typ == tTrue || typ == tFalse {
		t.byte()
		return
	}
	t.skip(typ)
}

func (t *tdec) skip(typ byte) {
	switch typ {
	case tTrue, tFalse:
	case tByte:
		t.byte()
	case tI16, tI32, tI64:
		t.uvarint()
	case tDouble:
		t.bytes(8)
	case tUUID:
		t.bytes(16)
	case tBinary:
		t.binary()
	case tList, tSet:
		if t.depth++; t.depth > maxThriftDepth {
			t.fail()
		}
		et, n := t.list()
		for i := 0; i < n && t.err == nil; i++ {
			t.skipElem(et)
		}
		t.depth--
	case tMap:
		if t.depth++; t.depth > maxThriftDepth {
			t.fail()
		}
		n := t.uvarint()
		if n > uint64(len(t.b)-t.pos) {
			t.fail()
		}
		if n > 0 {
			kv := t.byte()
			for i := uint64(0); i < n && t.err == nil; i++ {
				t.skipElem(kv >> 4)
				t.skipElem(kv & 0x0f)
			}
		}
		t.depth--
	case tStruct:
		t.fields(func(_ int16, typ byte) { t.skip(typ) })
	default:
		t.fail()
	}
}

func (t *tdec) bytes(n int) []byte {
	if n > len(t.b)-t.pos {
		t.fail()
		return nil
	}
	b := t.b[t.pos : t.pos+n]
	t.pos += n
	return b
}

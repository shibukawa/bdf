package parquet

import (
	"encoding/binary"
	"errors"
	"math"
)

// The file metadata and the page headers of a Parquet file are Thrift
// structures in the compact protocol: a struct is a list of fields, each a
// header byte (the field id as a delta from the previous one, and the
// type) and its value, ended by a zero byte. Integers are zigzag varints,
// binaries and strings are a varint length and their bytes, and a list is
// a header byte (its length, and the type of its elements) and its
// elements. The reader decodes the fields it knows and skips the others.

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

// maxThriftDepth bounds the nesting of structures and containers.
const maxThriftDepth = 64

var errThrift = errors.New("malformed Thrift structure")

// thrift reads the compact protocol from a buffer. The first error stops
// it: later reads return zero values, and err tells.
type thrift struct {
	b     []byte
	pos   int
	err   error
	depth int
}

func (t *thrift) fail() {
	if t.err == nil {
		t.err = errThrift
	}
	t.pos = len(t.b)
}

func (t *thrift) byte() byte {
	if t.pos >= len(t.b) {
		t.fail()
		return 0
	}
	c := t.b[t.pos]
	t.pos++
	return c
}

func (t *thrift) uvarint() uint64 {
	v, n := binary.Uvarint(t.b[t.pos:])
	if n <= 0 {
		t.fail()
		return 0
	}
	t.pos += n
	return v
}

func (t *thrift) varint() int64 {
	u := t.uvarint()
	return int64(u>>1) ^ -int64(u&1)
}

func (t *thrift) i32() int32 {
	v := t.varint()
	if v < math.MinInt32 || v > math.MaxInt32 {
		t.fail()
		return 0
	}
	return int32(v)
}

func (t *thrift) i64() int64 { return t.varint() }

func (t *thrift) binary() []byte {
	n := t.uvarint()
	if n > uint64(len(t.b)-t.pos) {
		t.fail()
		return nil
	}
	b := t.b[t.pos : t.pos+int(n)]
	t.pos += int(n)
	return b
}

func (t *thrift) string() string { return string(t.binary()) }

// fields reads a struct, calling f with the id and the type of each field;
// f reads the value, or skips it.
func (t *thrift) fields(f func(id int16, typ byte)) {
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
func (t *thrift) list() (typ byte, n int) {
	h := t.byte()
	typ = h & 0x0f
	size := uint64(h >> 4)
	if size == 15 {
		size = t.uvarint()
	}
	// an element takes a byte at least
	if size > uint64(len(t.b)-t.pos) {
		t.fail()
		return 0, 0
	}
	return typ, int(size)
}

// listOf reads a list whose elements are of type want, calling f to read
// each; the elements of a list of another type are skipped.
func (t *thrift) listOf(typ, want byte, f func()) {
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

// strings reads a list of strings.
func (t *thrift) strings(typ byte) []string {
	var out []string
	t.listOf(typ, tBinary, func() { out = append(out, t.string()) })
	return out
}

// skipElem skips an element of a container, where booleans take a byte.
func (t *thrift) skipElem(typ byte) {
	if typ == tTrue || typ == tFalse {
		t.byte()
		return
	}
	t.skip(typ)
}

// skip skips a field value of a type.
func (t *thrift) skip(typ byte) {
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

func (t *thrift) bytes(n int) []byte {
	if n > len(t.b)-t.pos {
		t.fail()
		return nil
	}
	b := t.b[t.pos : t.pos+n]
	t.pos += n
	return b
}

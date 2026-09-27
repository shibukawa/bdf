package parquet

import (
	"bytes"
	"encoding/binary"
	"math/bits"
	"strings"
	"testing"
)

// A writer of small Parquet files, for the structures that today's
// writers no longer make: two-level lists, maps annotated MAP_KEY_VALUE,
// repeated fields without annotations, BIT_PACKED levels, Hadoop's LZ4
// framing.

// tenc writes the Thrift compact protocol.
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
	if n < 15 {
		e.b = append(e.b, byte(n)<<4|typ)
	} else {
		e.b = append(e.b, 0xf0|typ)
		e.b = binary.AppendUvarint(e.b, uint64(n))
	}
}
func (e *tenc) bytes() []byte { return append(e.b, tStop) }

// entry is a leaf's entry: its levels, and its value (int32, string or nil).
type entry struct {
	rep, def int
	v        any
}

// pageOpts tells how the page of a column is written.
type pageOpts struct {
	bitPacked bool  // BIT_PACKED levels (data pages v1)
	v2        bool  // a data page v2 whose values are not compressed
	codec     int32 // codecNone or codecLZ4 (Hadoop's framing)
}

// testColumn is a leaf of the schema with its entries.
type testColumn struct {
	path    []string
	phys    int32
	maxRep  int
	maxDef  int
	entries []entry
}

// levels encodes levels in the hybrid encoding (a run each) or
// BIT_PACKED.
func levels(ls []int, max int, bitPacked bool) []byte {
	width := bits.Len(uint(max))
	if bitPacked {
		out := make([]byte, (len(ls)*width+7)/8)
		for i, l := range ls {
			for k := range width {
				if l>>(width-1-k)&1 != 0 {
					p := i*width + k
					out[p/8] |= 0x80 >> (p % 8)
				}
			}
		}
		return out
	}
	var out []byte
	for _, l := range ls {
		out = append(out, 2) // a run of one
		for k := 0; k < (width+7)/8; k++ {
			out = append(out, byte(l>>(8*k)))
		}
	}
	return out
}

func plainValues(c *testColumn) []byte {
	var out []byte
	for _, e := range c.entries {
		switch v := e.v.(type) {
		case int32:
			out = binary.LittleEndian.AppendUint32(out, uint32(v))
		case string:
			out = binary.LittleEndian.AppendUint32(out, uint32(len(v)))
			out = append(out, v...)
		}
	}
	return out
}

// hadoopLZ4 frames data as Hadoop's LZ4 codec does, in one block of
// literals.
func hadoopLZ4(data []byte) []byte {
	var block []byte
	n := len(data)
	if n < 15 {
		block = append(block, byte(n)<<4)
	} else {
		block = append(block, 0xf0)
		for n -= 15; n >= 255; n -= 255 {
			block = append(block, 255)
		}
		block = append(block, byte(n))
	}
	block = append(block, data...)
	out := binary.BigEndian.AppendUint32(nil, uint32(len(data)))
	out = binary.BigEndian.AppendUint32(out, uint32(len(block)))
	return append(out, block...)
}

func (c *testColumn) page(o pageOpts) []byte {
	var reps, defs []int
	for _, e := range c.entries {
		reps, defs = append(reps, e.rep), append(defs, e.def)
	}
	vals := plainValues(c)
	h := newTenc()
	if o.v2 {
		var rl, dl []byte
		if c.maxRep > 0 {
			rl = levels(reps, c.maxRep, false)
		}
		if c.maxDef > 0 {
			dl = levels(defs, c.maxDef, false)
		}
		body := append(append(rl, dl...), vals...)
		h.i32(1, pageDataV2)
		h.i32(2, int32(len(body)))
		h.i32(3, int32(len(body)))
		h.begin(8)
		h.i32(1, int32(len(c.entries)))
		h.i32(2, 0)
		h.i32(3, 0)
		h.i32(4, encPlain)
		h.i32(5, int32(len(dl)))
		h.i32(6, int32(len(rl)))
		h.boolean(7, false)
		h.end()
		return append(h.bytes(), body...)
	}
	var body []byte
	enc := int32(encRLE)
	if o.bitPacked {
		enc = encBitPacked
	}
	for _, l := range []struct {
		ls  []int
		max int
	}{{reps, c.maxRep}, {defs, c.maxDef}} {
		if l.max == 0 {
			continue
		}
		b := levels(l.ls, l.max, o.bitPacked)
		if !o.bitPacked {
			body = binary.LittleEndian.AppendUint32(body, uint32(len(b)))
		}
		body = append(body, b...)
	}
	body = append(body, vals...)
	data := body
	if o.codec == codecLZ4 {
		data = hadoopLZ4(body)
	}
	h.i32(1, pageData)
	h.i32(2, int32(len(body)))
	h.i32(3, int32(len(data)))
	h.begin(5)
	h.i32(1, int32(len(c.entries)))
	h.i32(2, encPlain)
	h.i32(3, enc)
	h.i32(4, enc)
	h.end()
	return append(h.bytes(), data...)
}

// buildFile writes a file of one row group.
func buildFile(schema []schemaElement, rows int64, cols []*testColumn, o pageOpts) []byte {
	b := []byte("PAR1")
	var offsets, sizes []int64
	for _, c := range cols {
		p := c.page(o)
		offsets, sizes = append(offsets, int64(len(b))), append(sizes, int64(len(p)))
		b = append(b, p...)
	}
	e := newTenc()
	e.i32(1, 1)
	e.list(2, tStruct, len(schema))
	for _, s := range schema {
		e.elem()
		if s.typ >= 0 {
			e.i32(1, s.typ)
		}
		if s.repetition >= 0 {
			e.i32(3, s.repetition)
		}
		e.bin(4, s.name)
		if s.numChildren > 0 {
			e.i32(5, s.numChildren)
		}
		if s.converted >= 0 {
			e.i32(6, s.converted)
		}
		if s.logical.kind != logNone {
			e.begin(10)
			e.begin(s.logical.kind)
			e.end()
			e.end()
		}
		e.end()
	}
	e.i64(3, rows)
	e.list(4, tStruct, 1)
	e.elem()
	e.list(1, tStruct, len(cols))
	for i, c := range cols {
		e.elem()
		e.i64(2, offsets[i])
		e.begin(3)
		e.i32(1, c.phys)
		e.list(2, tI32, 1)
		e.b = binary.AppendVarint(e.b, encPlain)
		e.list(3, tBinary, len(c.path))
		for _, p := range c.path {
			e.b = binary.AppendUvarint(e.b, uint64(len(p)))
			e.b = append(e.b, p...)
		}
		e.i32(4, o.codec)
		e.i64(5, int64(len(c.entries)))
		e.i64(6, sizes[i])
		e.i64(7, sizes[i])
		e.i64(9, offsets[i])
		e.end()
		e.end()
	}
	e.i64(2, 0)
	e.i64(3, rows)
	e.end()
	footer := e.bytes()
	b = append(b, footer...)
	b = binary.LittleEndian.AppendUint32(b, uint32(len(footer)))
	return append(b, "PAR1"...)
}

// el makes schema elements: groups when typ is -1.
func el(name string, typ, rep, children, conv int32, logical int16) schemaElement {
	return schemaElement{typ: typ, repetition: rep, name: name, numChildren: children, converted: conv, logical: logicalType{kind: logical}}
}

func root(children int32) schemaElement { return el("schema", -1, -1, children, convNone, logNone) }

func TestLegacyStructures(t *testing.T) {
	for _, c := range []struct {
		name   string
		schema []schemaElement
		rows   int64
		cols   []*testColumn
		want   []string // the type, then the rows
	}{
		{"two-level list (rule 1)",
			[]schemaElement{root(1), el("my_list", -1, repOptional, 1, convList, logNone), el("element", typeInt32, repRepeated, 0, convNone, logNone)},
			4, []*testColumn{{[]string{"my_list", "element"}, typeInt32, 1, 2,
				[]entry{{0, 2, int32(1)}, {1, 2, int32(2)}, {0, 0, nil}, {0, 1, nil}, {0, 2, int32(3)}}}},
			[]string{"list<int32>", "[1, 2]", "", "[]", "[3]"}},
		{"repeated group of fields (rule 2)",
			[]schemaElement{root(1), el("pairs", -1, repOptional, 1, convList, logNone), el("element", -1, repRepeated, 2, convNone, logNone),
				el("str", typeBinary, repRequired, 0, convUTF8, logNone), el("num", typeInt32, repRequired, 0, convNone, logNone)},
			2, []*testColumn{
				{[]string{"pairs", "element", "str"}, typeBinary, 1, 2, []entry{{0, 2, "a"}, {1, 2, "b"}, {0, 0, nil}}},
				{[]string{"pairs", "element", "num"}, typeInt32, 1, 2, []entry{{0, 2, int32(1)}, {1, 2, int32(2)}, {0, 0, nil}}}},
			[]string{"list<struct<str: string, num: int32>>", `[{"str": "a", "num": 1}, {"str": "b", "num": 2}]`, ""}},
		{"list of two-level lists (rule 3)",
			[]schemaElement{root(1), el("lol", -1, repOptional, 1, convList, logNone), el("array", -1, repRepeated, 1, convList, logNone),
				el("array", typeInt32, repRepeated, 0, convNone, logNone)},
			2, []*testColumn{{[]string{"lol", "array", "array"}, typeInt32, 2, 3,
				[]entry{{0, 3, int32(1)}, {2, 3, int32(2)}, {1, 3, int32(3)}, {0, 2, nil}}}},
			[]string{"list<list<int32>>", "[[1, 2], [3]]", "[[]]"}},
		{"one-field tuples (rule 4)",
			[]schemaElement{root(1), el("tl", -1, repOptional, 1, convList, logNone), el("tl_tuple", -1, repRepeated, 1, convNone, logNone),
				el("str", typeBinary, repRequired, 0, convUTF8, logNone)},
			1, []*testColumn{{[]string{"tl", "tl_tuple", "str"}, typeBinary, 1, 2, []entry{{0, 2, "x"}, {1, 2, "y"}}}},
			[]string{"list<struct<str: string>>", `[{"str": "x"}, {"str": "y"}]`}},
		{"misnamed three-level list (rule 5)",
			[]schemaElement{root(1), el("names", -1, repOptional, 1, convNone, logList), el("element", -1, repRepeated, 1, convNone, logNone),
				el("str", typeBinary, repOptional, 0, convNone, logString)},
			1, []*testColumn{{[]string{"names", "element", "str"}, typeBinary, 1, 3, []entry{{0, 3, "a"}, {1, 2, nil}}}},
			[]string{"list<string>", `["a", null]`}},
		{"MAP_KEY_VALUE map",
			[]schemaElement{root(1), el("m", -1, repOptional, 1, convMapKeyValue, logNone), el("map", -1, repRepeated, 2, convNone, logNone),
				el("key", typeBinary, repRequired, 0, convUTF8, logNone), el("value", typeInt32, repOptional, 0, convNone, logNone)},
			1, []*testColumn{
				{[]string{"m", "map", "key"}, typeBinary, 1, 2, []entry{{0, 2, "k1"}, {1, 2, "k2"}}},
				{[]string{"m", "map", "value"}, typeInt32, 1, 3, []entry{{0, 3, int32(1)}, {1, 2, nil}}}},
			[]string{"map<string, int32>", `{"k1": 1, "k2": null}`}},
		{"repeated field without annotation",
			[]schemaElement{root(1), el("xs", typeInt32, repRepeated, 0, convNone, logNone)},
			2, []*testColumn{{[]string{"xs"}, typeInt32, 1, 1, []entry{{0, 1, int32(5)}, {1, 1, int32(6)}, {0, 0, nil}}}},
			[]string{"list<int32>", "[5, 6]", "[]"}},
	} {
		for _, o := range []pageOpts{{}, {bitPacked: true}, {v2: true}, {codec: codecLZ4}} {
			b := buildFile(c.schema, c.rows, c.cols, o)
			g, _, err := table(bytes.NewReader(b), int64(len(b)), &Options{}, func(msg string) { t.Errorf("%s %+v: warning: %s", c.name, o, msg) })
			if err != nil {
				t.Fatalf("%s %+v: %v", c.name, o, err)
			}
			var got []string
			for _, row := range g.Rows[1:] {
				got = append(got, row[0].Text)
			}
			if strings.Join(got, "\n") != strings.Join(c.want, "\n") {
				t.Errorf("%s %+v:\n got  %q\n want %q", c.name, o, got, c.want)
			}
		}
	}
}

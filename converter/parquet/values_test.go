package parquet

import (
	"encoding/binary"
	"math"
	"math/big"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	ts := func(unit int16, v int64) string {
		return format{kind: fTimestamp, unit: unit}.text(scalar{u: uint64(v)}, -1)
	}
	for _, c := range []struct{ got, want string }{
		{decimalText(big.NewInt(12345), 2), "123.45"},
		{decimalText(big.NewInt(-5), 3), "-0.005"},
		{decimalText(big.NewInt(0), 2), "0.00"},
		{decimalText(big.NewInt(7), 0), "7"},
		{decimalText(big.NewInt(7), -2), "700"},
		{decimalText(bigEndianInt([]byte{0xff, 0x85}), 1), "-12.3"},
		{decimalText(bigEndianInt([]byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0}), 0), "-2361183241434822606848"},
		{floatText(0.1, 64), "0.1"},
		{floatText(-1e-7, 64), "-1e-07"},
		{floatText(123456789012345680000, 64), "123456789012345680000"},
		{floatText(1e21, 64), "1e+21"},
		{floatText(float64(float32(0.1)), 32), "0.1"},
		{floatText(math.NaN(), 64), "NaN"},
		{floatText(math.Inf(-1), 64), "-Infinity"},
		{halfText(0x3c00), "1"},
		{halfText(0x3555), "0.3333"},
		{halfText(0x7bff), "65500"},
		{halfText(0x0001), "6e-08"},
		{halfText(0xc000), "-2"},
		{halfText(0x7c00), "Infinity"},
		{halfText(0x7e00), "NaN"},
		{halfText(0x8000), "-0"},
		{dateText(-1), "1969-12-31"},
		{ts(unitMillis, -1), "1969-12-31 23:59:59.999"},
		{ts(unitMicros, 1_500_000), "1970-01-01 00:00:01.500"},
		{ts(unitNanos, 1_000_000_001), "1970-01-01 00:00:01.000000001"},
		{format{kind: fTime, unit: unitMillis}.text(scalar{u: 45296789}, -1), "12:34:56.789"},
		{format{kind: fTime, unit: unitMicros}.text(scalar{u: 45296000000}, 6), "12:34:56.000000"},
		{format{kind: fInt96}.text(scalar{u: 3_600_000_000_001, x: 2440589}, -1), "1970-01-02 01:00:00.000000001"},
		{format{kind: fUint, phys: typeInt32}.text(scalar{u: math.MaxUint64}, -1), "4294967295"},
		{intervalText([]byte{14, 0, 0, 0, 3, 0, 0, 0, 0xec, 0xd0, 0x38, 0}), "P1Y2M3DT1H2M3.5S"},
		{intervalText(make([]byte, 12)), "PT0S"},
		{intervalText([]byte{0, 0, 0, 0, 0, 0, 0, 0, 10, 0, 0, 0}), "PT0.01S"},
		{uuidText([]byte{0xf2, 0x4f, 0x9b, 0x64, 0x81, 0xfa, 0x49, 0xd1, 0xb7, 0x4e, 0x8c, 0x09, 0xa6, 0xe3, 0x1c, 0x56}), "f24f9b64-81fa-49d1-b74e-8c09a6e31c56"},
		{binaryText([]byte("tab\tand\nlines")), "tab\tand\nlines"},
		{binaryText([]byte("bell\a")), "0x62656c6c07"},
		{binaryText([]byte{0xe3, 0x81}), "0xe381"},
		{binaryText(nil), "0x"},
		{validText([]byte("ok\xff")), "ok�"},
	} {
		if c.got != c.want {
			t.Errorf("%q, want %q", c.got, c.want)
		}
	}
	if s := hexText(make([]byte, 100)); !strings.HasSuffix(s, "… (100 bytes)") || len(s) != 2+64+len("… (100 bytes)") {
		t.Errorf("long binary: %q", s)
	}
	if s := cellText(strings.Repeat("あ", 20000)); len(s) > maxCellText+len("…") || !strings.HasSuffix(s, "あ…") {
		t.Errorf("long text: %d bytes", len(s))
	}
}

func TestWKT(t *testing.T) {
	le := func(vals ...any) []byte {
		var b []byte
		for _, v := range vals {
			switch v := v.(type) {
			case byte:
				b = append(b, v)
			case uint32:
				b = binary.LittleEndian.AppendUint32(b, v)
			case float64:
				b = binary.LittleEndian.AppendUint64(b, math.Float64bits(v))
			}
		}
		return b
	}
	be := binary.BigEndian.AppendUint64(binary.BigEndian.AppendUint64(append([]byte{0}, 0, 0, 0, 1), math.Float64bits(1)), math.Float64bits(2))
	for _, c := range []struct {
		b    []byte
		want string
	}{
		{be, "POINT (1 2)"},
		{le(byte(1), uint32(0xa0000001), uint32(4326), 1.5, 2.5, 3.5), "POINT Z (1.5 2.5 3.5)"},             // EWKB: Z and an SRID
		{le(byte(1), uint32(2002), uint32(2), 0.0, 0.0, 9.0, 1.0, 1.0, 8.0), "LINESTRING M (0 0 9, 1 1 8)"}, // ISO M
		{le(byte(1), uint32(3001), 1.0, 2.0, 3.0, 4.0), "POINT ZM (1 2 3 4)"},
		{le(byte(1), uint32(2), uint32(0)), "LINESTRING EMPTY"},
		{le(byte(1), uint32(3), uint32(0)), "POLYGON EMPTY"},
		{le(byte(1), uint32(5), uint32(1), byte(1), uint32(2), uint32(2), 0.0, 0.0, 1.0, 1.0), "MULTILINESTRING ((0 0, 1 1))"},
		{le(byte(1), uint32(7), uint32(0)), "GEOMETRYCOLLECTION EMPTY"},
		{le(byte(1), uint32(7), uint32(1), byte(1), uint32(7), uint32(1), byte(1), uint32(1), 1.0, 1.0), "GEOMETRYCOLLECTION (GEOMETRYCOLLECTION (POINT (1 1)))"},
	} {
		if got, ok := wkt(c.b); !ok || got != c.want {
			t.Errorf("%x: %q (%v), want %q", c.b, got, ok, c.want)
		}
	}
	for _, bad := range [][]byte{
		nil,
		{2, 1, 0, 0, 0},             // byte order
		le(byte(1), uint32(9)),      // type
		le(byte(1), uint32(1), 1.0), // truncated
		le(byte(1), uint32(2), uint32(0x7fffffff)),  // count
		append(le(byte(1), uint32(1), 1.0, 2.0), 0), // trailing bytes
	} {
		if s, ok := wkt(bad); ok {
			t.Errorf("%x: %q", bad, s)
		}
	}
}

// variantMetadata makes the metadata of a Variant with field names.
func variantMetadata(names ...string) []byte {
	b := []byte{0x01, byte(len(names))}
	off := 0
	b = append(b, 0)
	for _, n := range names {
		off += len(n)
		b = append(b, byte(off))
	}
	for _, n := range names {
		b = append(b, n...)
	}
	return b
}

func TestVariant(t *testing.T) {
	prim := func(typ byte, data ...byte) []byte { return append([]byte{typ << 2}, data...) }
	le32 := func(v uint32) []byte { return binary.LittleEndian.AppendUint32(nil, v) }
	le64 := func(v uint64) []byte { return binary.LittleEndian.AppendUint64(nil, v) }
	short := func(s string) []byte { return append([]byte{byte(len(s))<<2 | 1}, s...) }
	meta := variantMetadata("a", "b")
	// an object {"a": 1, "b": [true, "x"]}: ids and offsets of one byte
	arr := append([]byte{3, 2, 0, 1, 3}, append(prim(1), short("x")...)...)
	obj := append([]byte{2, 2, 0, 1, 0, 2, byte(2 + len(arr))}, append(prim(3, 1), arr...)...)
	for _, c := range []struct {
		v    []byte
		want string
	}{
		{prim(0), "null"},
		{prim(2), "false"},
		{prim(3, 0xff), "-1"},
		{prim(4, 0x00, 0x80), "-32768"},
		{prim(5, le32(123456)...), "123456"},
		{prim(6, le64(1<<62)...), "4611686018427387904"},
		{prim(7, le64(math.Float64bits(1.25))...), "1.25"},
		{prim(8, append([]byte{2}, le32(1234)...)...), "12.34"},
		{prim(9, append([]byte{1}, le64(uint64(1<<63))...)...), "-922337203685477580.8"},
		{prim(10, append([]byte{0}, append(le64(1), le64(0)...)...)...), "1"},
		{prim(10, append([]byte{0}, append(le64(^uint64(0)), le64(^uint64(0))...)...)...), "-1"},
		{prim(11, le32(uint32(20000))...), `"2024-10-04"`},
		{prim(12, le64(1_500_000)...), `"1970-01-01 00:00:01.500"`},
		{prim(19, le64(1)...), `"1970-01-01 00:00:00.000000001"`},
		{prim(14, le32(math.Float32bits(0.1))...), "0.1"},
		{prim(15, append(le32(2), 0xca, 0xfe)...), `"0xcafe"`},
		{prim(16, append(le32(4), "long"...)...), `"long"`},
		{prim(17, le64(45296000001)...), `"12:34:56.000001"`},
		{prim(20, 0xf2, 0x4f, 0x9b, 0x64, 0x81, 0xfa, 0x49, 0xd1, 0xb7, 0x4e, 0x8c, 0x09, 0xa6, 0xe3, 0x1c, 0x56), `"f24f9b64-81fa-49d1-b74e-8c09a6e31c56"`},
		{short("quote\"d"), `"quote\"d"`},
		{obj, `{"a": 1, "b": [true, "x"]}`},
		{[]byte{3, 0, 0}, "[]"},
		// malformed values are null
		{prim(5, 1, 2), "null"},
		{prim(16, le32(100)...), "null"},
		{[]byte{2, 1, 9, 0, 1, 0}, "null"}, // a field id out of the dictionary
		{[]byte{3, 2, 0, 5, 1}, "null"},    // offsets out of order
		{nil, "null"},
	} {
		m, err := parseVariantMeta(meta)
		if err != nil {
			t.Fatal(err)
		}
		w := &writer{}
		w.variant(m, c.v)
		if got := w.b.String(); got != c.want {
			t.Errorf("%x: %s, want %s", c.v, got, c.want)
		}
	}
	for _, bad := range [][]byte{nil, {0x02, 0, 0}, {0x01, 5, 0}, {0xc1, 1, 0, 0}} {
		if _, err := parseVariantMeta(bad); err == nil {
			t.Errorf("metadata %x: no error", bad)
		}
	}
	// deep nesting ends
	deep := []byte{}
	for range 500 {
		deep = append([]byte{3, 1, 0, byte(len(deep))}, deep...)
		if len(deep) > 250 {
			break
		}
	}
	w := &writer{}
	m, _ := parseVariantMeta(meta)
	w.variant(m, deep)
	if w.b.Len() == 0 {
		t.Error("deep array: nothing written")
	}
}

func TestSchemaErrors(t *testing.T) {
	for _, elems := range [][]schemaElement{
		nil,
		{el("schema", -1, -1, 3, convNone, logNone), el("a", typeInt32, repRequired, 0, convNone, logNone)},
		{el("a", typeInt32, repRequired, 0, convNone, logNone)},
		{root(1), el("a", 42, repRequired, 0, convNone, logNone)},
		{root(1), el("f", typeFixed, repRequired, 0, convNone, logNone)},
		{root(1), el("a", typeInt32, repRequired, 0, convNone, logNone), el("extra", typeInt32, repRequired, 0, convNone, logNone)},
	} {
		if _, _, err := buildSchema(elems); err == nil {
			t.Errorf("%+v: no error", elems)
		}
	}
	// logical types that do not apply to the physical type are ignored
	root, leaves, err := buildSchema([]schemaElement{root(2),
		el("u", typeInt32, repRequired, 0, convNone, logUUID), el("s", typeInt64, repRequired, 0, convNone, logString)})
	if err != nil {
		t.Fatal(err)
	}
	if typeName(root.children[0]) != "int32" || typeName(root.children[1]) != "int64" || formatOf(leaves[0]).kind != fInt {
		t.Errorf("types %q %q", typeName(root.children[0]), typeName(root.children[1]))
	}
}

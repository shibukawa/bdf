package parquet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/bits"
	"time"

	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/snappy"
)

// leaf is a leaf of the schema: a column chunk in every row group.
type leaf struct {
	name    string
	top     int // the column at the top of the schema it belongs to
	flat    bool
	phys    int32
	typeLen int32
	maxDef  int
	maxRep  int
	conv    int32
	logical int16
	unit    int16
	scale   int32
	typ     Type
}

type chunkMeta struct {
	codec     int32
	numValues int64
	size      int64 // compressed, with the page headers
	dataOff   int64
	dictOff   int64
}

type rowGroup struct {
	numRows int64
	chunks  []chunkMeta
}

type fileMeta struct {
	columns   []Column
	leaves    []*leaf
	numRows   int64
	rowGroups []rowGroup
}

// schemaElement is an element of the flat list the file holds its schema
// as, in depth-first order.
type schemaElement struct {
	typ         int32 // -1 for groups
	typeLength  int32
	repetition  int32
	name        string
	numChildren int32
	converted   int32
	scale       int32
	logical     int16
	unit        int16
}

// Schema reads the columns of a file and its row count from its footer.
// A column of a type the package does not take (a group, a repeated field,
// a decimal of bytes) is of type Other.
func Schema(data []byte) ([]Column, int64, error) {
	m, err := readMeta(data)
	if err != nil {
		return nil, 0, err
	}
	return m.columns, m.numRows, nil
}

// readMeta reads the footer of a file.
func readMeta(data []byte) (*fileMeta, error) {
	if len(data) < 12 || string(data[:4]) != "PAR1" || string(data[len(data)-4:]) != "PAR1" {
		if len(data) >= 4 && string(data[:4]) == "PARE" {
			return nil, errors.New("parquet: the file is encrypted")
		}
		return nil, ErrNotParquet
	}
	n := int64(binary.LittleEndian.Uint32(data[len(data)-8:]))
	if n <= 0 || n > maxFooter || n > int64(len(data))-12 {
		return nil, errCorrupt
	}
	footer := data[int64(len(data))-8-n : len(data)-8]
	t := &tdec{b: footer}
	m := &fileMeta{}
	var elems []schemaElement
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 2 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() { elems = append(elems, readSchemaElement(t)) })
		case id == 3 && typ == tI64:
			m.numRows = t.i64()
		case id == 4 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() { m.rowGroups = append(m.rowGroups, readRowGroup(t)) })
		default:
			t.skip(typ)
		}
	})
	if t.err != nil {
		return nil, t.err
	}
	if m.numRows < 0 || len(elems) == 0 {
		return nil, errCorrupt
	}
	if err := m.schema(elems); err != nil {
		return nil, err
	}
	for _, g := range m.rowGroups {
		if g.numRows < 0 || len(g.chunks) != len(m.leaves) {
			return nil, errCorrupt
		}
	}
	return m, nil
}

func readSchemaElement(t *tdec) schemaElement {
	e := schemaElement{typ: -1, repetition: -1, converted: convNone}
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && typ == tI32:
			e.typ = t.i32()
		case id == 2 && typ == tI32:
			e.typeLength = t.i32()
		case id == 3 && typ == tI32:
			e.repetition = t.i32()
		case id == 4 && typ == tBinary:
			e.name = t.string()
		case id == 5 && typ == tI32:
			e.numChildren = t.i32()
		case id == 6 && typ == tI32:
			e.converted = t.i32()
		case id == 7 && typ == tI32:
			e.scale = t.i32()
		case id == 10 && typ == tStruct:
			t.fields(func(lid int16, typ byte) {
				if typ != tStruct {
					t.skip(typ)
					return
				}
				e.logical = lid
				t.fields(func(fid int16, typ byte) {
					switch {
					case lid == logDecimal && fid == 1 && typ == tI32:
						e.scale = t.i32()
					case lid == logTimestamp && fid == 2 && typ == tStruct:
						t.fields(func(uid int16, typ byte) {
							e.unit = uid
							t.skip(typ)
						})
					default:
						t.skip(typ)
					}
				})
			})
		default:
			t.skip(typ)
		}
	})
	return e
}

func readRowGroup(t *tdec) rowGroup {
	var g rowGroup
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && (typ == tList || typ == tSet):
			t.listOf(typ, tStruct, func() { g.chunks = append(g.chunks, readColumnChunk(t)) })
		case id == 3 && typ == tI64:
			g.numRows = t.i64()
		default:
			t.skip(typ)
		}
	})
	return g
}

func readColumnChunk(t *tdec) chunkMeta {
	var c chunkMeta
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 3 && typ == tStruct:
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 4 && typ == tI32:
					c.codec = t.i32()
				case id == 5 && typ == tI64:
					c.numValues = t.i64()
				case id == 7 && typ == tI64:
					c.size = t.i64()
				case id == 9 && typ == tI64:
					c.dataOff = t.i64()
				case id == 11 && typ == tI64:
					c.dictOff = t.i64()
				default:
					t.skip(typ)
				}
			})
		default:
			t.skip(typ)
		}
	})
	return c
}

// schema builds the columns and the leaves from the flat list of elements.
func (m *fileMeta) schema(elems []schemaElement) error {
	root := elems[0]
	if root.numChildren < 0 || int(root.numChildren) > len(elems)-1 {
		return errCorrupt
	}
	if root.numChildren > MaxColumns {
		return fmt.Errorf("parquet: %d columns, more than %d", root.numChildren, MaxColumns)
	}
	pos := 1
	// walk reads the element at pos and its children, a leaf each path
	var walk func(top, def, rep int, inGroup bool) error
	walk = func(top, def, rep int, inGroup bool) error {
		if pos >= len(elems) {
			return errCorrupt
		}
		e := elems[pos]
		pos++
		switch e.repetition {
		case repOptional:
			def++
		case repRepeated:
			def++
			rep++
		}
		if e.numChildren > 0 || e.typ < 0 {
			if int(e.numChildren) > len(elems)-pos || len(m.leaves) > MaxColumns*4 {
				return errCorrupt
			}
			for range int(e.numChildren) {
				if err := walk(top, def, rep, true); err != nil {
					return err
				}
			}
			return nil
		}
		l := &leaf{name: e.name, top: top, flat: !inGroup && e.repetition != repRepeated, phys: e.typ, typeLen: e.typeLength,
			maxDef: def, maxRep: rep, conv: e.converted, logical: e.logical, unit: e.unit, scale: e.scale}
		l.typ = typeOf(l)
		m.leaves = append(m.leaves, l)
		return nil
	}
	for i := range int(root.numChildren) {
		e := elems[pos]
		start := len(m.leaves)
		if err := walk(i, 0, 0, false); err != nil {
			return err
		}
		c := Column{Name: e.name, Type: Other}
		if len(m.leaves) == start+1 && m.leaves[start].flat {
			c.Type = m.leaves[start].typ
		}
		m.columns = append(m.columns, c)
	}
	return nil
}

// typeOf is the type of the column of a leaf.
func typeOf(l *leaf) Type {
	switch l.phys {
	case typeBoolean:
		return Boolean
	case typeInt32:
		switch {
		case l.logical == logDate || l.conv == convDate:
			return Date
		case l.logical == logDecimal || l.conv == convDecimal:
			return Number
		}
		return Integer
	case typeInt64:
		switch {
		case l.logical == logTimestamp || l.conv == convTSMillis || l.conv == convTSMicros:
			return DateTime
		case l.logical == logDecimal || l.conv == convDecimal:
			return Number
		}
		return Integer
	case typeInt96:
		return DateTime
	case typeFloat, typeDouble:
		return Number
	case typeBinary:
		switch {
		case l.logical == logString || l.logical == logEnum || l.logical == logJSON,
			l.conv == convUTF8 || l.conv == convEnum || l.conv == convJSON:
			return String
		}
	}
	return Other
}

// Read reads the rows of a file, up to limits. Columns of type Other hold
// nil, except binaries, which hold their bytes.
func Read(data []byte, limits Limits) (*Table, error) {
	m, err := readMeta(data)
	if err != nil {
		return nil, err
	}
	maxRows := limits.MaxRows
	if maxRows <= 0 {
		maxRows = DefaultMaxRows
	}
	n := int(min(m.numRows, int64(maxRows)))
	t := &Table{Columns: m.columns, Rows: make([][]any, n)}
	cells := make([]any, n*len(m.columns))
	for i := range t.Rows {
		t.Rows[i] = cells[i*len(m.columns) : (i+1)*len(m.columns) : (i+1)*len(m.columns)]
	}
	base := 0
	for gi := range m.rowGroups {
		g := &m.rowGroups[gi]
		rows := int(min(g.numRows, int64(n-base)))
		if rows <= 0 {
			break
		}
		for li, l := range m.leaves {
			if !l.flat {
				continue
			}
			vals, err := readChunk(data, l, &g.chunks[li], rows)
			if err != nil {
				return nil, fmt.Errorf("parquet: column %q: %w", l.name, err)
			}
			for i, v := range vals {
				t.Rows[base+i][l.top] = v
			}
		}
		base += rows
	}
	t.Rows = t.Rows[:base]
	return t, nil
}

// readChunk reads the first n values of a column chunk: a value for each
// row, nil for a null.
func readChunk(file []byte, l *leaf, c *chunkMeta, n int) ([]any, error) {
	if l.maxRep > 0 {
		return nil, errors.New("a repeated column")
	}
	pos := c.dataOff
	if c.dictOff > 0 && c.dictOff < pos {
		pos = c.dictOff
	}
	if pos < 4 || pos > int64(len(file)) {
		return nil, errCorrupt
	}
	out := make([]any, 0, n)
	var dict []any
	var left int64 = c.numValues
	for len(out) < n && left > 0 {
		h, hl, err := readPageHeader(file[pos:])
		if err != nil {
			return nil, err
		}
		pos += int64(hl)
		if h.compressed < 0 || int64(h.compressed) > int64(len(file))-pos || h.compressed > maxPage || h.uncompressed < 0 || h.uncompressed > maxPage {
			return nil, errCorrupt
		}
		page := file[pos : pos+int64(h.compressed)]
		pos += int64(h.compressed)
		switch h.typ {
		case pageDictionary:
			buf, err := decompress(c.codec, page, int(h.uncompressed))
			if err != nil {
				return nil, err
			}
			if h.encoding != encPlain && h.encoding != encPlainDict {
				return nil, fmt.Errorf("a dictionary page encoded %d", h.encoding)
			}
			d := &plain{b: buf, phys: l.phys, size: int(l.typeLen)}
			dict = make([]any, 0, max(h.numValues, 0))
			for range int(h.numValues) {
				v, err := d.next()
				if err != nil {
					return nil, err
				}
				dict = append(dict, v)
			}
		case pageData, pageDataV2:
			if h.numValues < 0 {
				return nil, errCorrupt
			}
			var defs *hybrid
			var vals []byte
			if h.typ == pageData {
				buf, err := decompress(c.codec, page, int(h.uncompressed))
				if err != nil {
					return nil, err
				}
				p := 0
				if l.maxDef > 0 {
					if h.defEncoding != encRLE {
						return nil, errors.New("definition levels that are not RLE")
					}
					if len(buf) < 4 {
						return nil, errCorrupt
					}
					ln := int(binary.LittleEndian.Uint32(buf))
					if ln < 0 || ln > len(buf)-4 {
						return nil, errCorrupt
					}
					defs = &hybrid{b: buf[4 : 4+ln], width: uint(bits.Len(uint(l.maxDef)))}
					p = 4 + ln
				}
				vals = buf[p:]
			} else {
				rl, dl := int(h.repLength), int(h.defLength)
				if rl < 0 || dl < 0 || rl+dl > len(page) {
					return nil, errCorrupt
				}
				if l.maxDef > 0 {
					defs = &hybrid{b: page[rl : rl+dl], width: uint(bits.Len(uint(l.maxDef)))}
				}
				vals = page[rl+dl:]
				if size := int(h.uncompressed) - rl - dl; size == 0 {
					vals = nil
				} else if !h.uncompressV2 && c.codec != codecNone {
					if vals, err = decompress(c.codec, vals, size); err != nil {
						return nil, err
					}
				}
			}
			var dec func() (any, error)
			switch h.encoding {
			case encPlain:
				d := &plain{b: vals, phys: l.phys, size: int(l.typeLen)}
				dec = d.next
			case encPlainDict, encRLEDict:
				if dict == nil {
					return nil, errors.New("dictionary-encoded values without a dictionary page")
				}
				if len(vals) == 0 {
					return nil, errCorrupt
				}
				idx := &hybrid{b: vals[1:], width: uint(vals[0])}
				if idx.width > 32 {
					return nil, errCorrupt
				}
				dec = func() (any, error) {
					i, err := idx.next()
					if err != nil {
						return nil, err
					}
					if i >= uint64(len(dict)) {
						return nil, errCorrupt
					}
					return dict[i], nil
				}
			case encRLE:
				if l.phys != typeBoolean {
					return nil, errors.New("RLE values of a column that is not boolean")
				}
				if len(vals) < 4 {
					return nil, errCorrupt
				}
				rle := &hybrid{b: vals[4:], width: 1}
				dec = func() (any, error) {
					v, err := rle.next()
					return v != 0, err
				}
			default:
				return nil, fmt.Errorf("the encoding %d is not read here (converter/parquet reads it)", h.encoding)
			}
			for range int(h.numValues) {
				if left--; left < 0 {
					break
				}
				def := l.maxDef
				if defs != nil {
					d, err := defs.next()
					if err != nil {
						return nil, err
					}
					def = int(d)
				}
				var v any
				if def == l.maxDef {
					if v, err = dec(); err != nil {
						return nil, err
					}
					v = convert(l, v)
				}
				if len(out) < n {
					out = append(out, v)
				}
			}
		default:
			// an index page, or one of a kind not known: skipped
		}
	}
	for len(out) < n {
		out = append(out, nil)
	}
	return out, nil
}

// pageHeader is what the reader needs of a page header.
type pageHeader struct {
	typ          int32
	uncompressed int32
	compressed   int32
	numValues    int32
	encoding     int32
	defEncoding  int32
	repLength    int32
	defLength    int32
	uncompressV2 bool
}

func readPageHeader(b []byte) (*pageHeader, int, error) {
	t := &tdec{b: b[:min(len(b), 1<<20)]}
	h := &pageHeader{}
	t.fields(func(id int16, typ byte) {
		switch {
		case id == 1 && typ == tI32:
			h.typ = t.i32()
		case id == 2 && typ == tI32:
			h.uncompressed = t.i32()
		case id == 3 && typ == tI32:
			h.compressed = t.i32()
		case id == 5 && typ == tStruct:
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 1 && typ == tI32:
					h.numValues = t.i32()
				case id == 2 && typ == tI32:
					h.encoding = t.i32()
				case id == 3 && typ == tI32:
					h.defEncoding = t.i32()
				default:
					t.skip(typ)
				}
			})
		case id == 7 && typ == tStruct:
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 1 && typ == tI32:
					h.numValues = t.i32()
				case id == 2 && typ == tI32:
					h.encoding = t.i32()
				default:
					t.skip(typ)
				}
			})
		case id == 8 && typ == tStruct:
			h.uncompressV2 = false
			compressed := true
			t.fields(func(id int16, typ byte) {
				switch {
				case id == 1 && typ == tI32:
					h.numValues = t.i32()
				case id == 4 && typ == tI32:
					h.encoding = t.i32()
				case id == 5 && typ == tI32:
					h.defLength = t.i32()
				case id == 6 && typ == tI32:
					h.repLength = t.i32()
				case id == 7 && (typ == tTrue || typ == tFalse):
					compressed = typ == tTrue
				default:
					t.skip(typ)
				}
			})
			h.uncompressV2 = !compressed
		default:
			t.skip(typ)
		}
	})
	if t.err != nil {
		return nil, 0, errCorrupt
	}
	return h, t.pos, nil
}

// decompress unpacks the bytes of a page.
func decompress(codec int32, data []byte, size int) ([]byte, error) {
	if size < 0 || size > maxPage {
		return nil, errCorrupt
	}
	switch codec {
	case codecNone:
		return data, nil
	case codecSnappy:
		out, err := snappy.Decode(make([]byte, 0, size), data)
		if err != nil {
			return nil, fmt.Errorf("snappy: %w", err)
		}
		return out, nil
	case codecGzip:
		r, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		out := make([]byte, size)
		if _, err := io.ReadFull(r, out); err != nil {
			return nil, err
		}
		return out, nil
	}
	names := map[int32]string{codecLZO: "LZO", codecBrotli: "Brotli", codecLZ4: "LZ4", codecZstd: "Zstandard", codecLZ4Raw: "LZ4"}
	if n := names[codec]; n != "" {
		return nil, fmt.Errorf("the %s codec is not read here (converter/parquet reads it)", n)
	}
	return nil, fmt.Errorf("unknown codec %d", codec)
}

// hybrid decodes the RLE/bit-packed hybrid of levels and dictionary
// indices: runs of a repeated value, and runs of values packed in width
// bits each.
type hybrid struct {
	b      []byte
	pos    int
	width  uint
	left   int
	packed bool
	value  uint64
	start  int
	bit    uint64
}

func (h *hybrid) next() (uint64, error) {
	for h.left == 0 {
		if h.pos >= len(h.b) {
			return 0, errCorrupt
		}
		hdr, n := binary.Uvarint(h.b[h.pos:])
		if n <= 0 {
			return 0, errCorrupt
		}
		h.pos += n
		if hdr&1 == 1 {
			groups := hdr >> 1
			size := groups * uint64(h.width)
			avail := uint64(len(h.b) - h.pos)
			count := groups * 8
			if size > avail {
				size = avail
				if h.width > 0 {
					count = avail * 8 / uint64(h.width)
				}
			}
			if count > 1<<30 {
				count = 1 << 30
			}
			h.packed, h.left, h.start, h.bit = true, int(count), h.pos, 0
			h.pos += int(size)
		} else {
			count := hdr >> 1
			if count > 1<<30 {
				return 0, errCorrupt
			}
			nb := int(h.width+7) / 8
			if nb > len(h.b)-h.pos {
				return 0, errCorrupt
			}
			var v uint64
			for i := range nb {
				v |= uint64(h.b[h.pos+i]) << (8 * i)
			}
			h.pos += nb
			h.packed, h.left, h.value = false, int(count), v
		}
	}
	h.left--
	if !h.packed {
		return h.value, nil
	}
	v := bitsLE(h.b[h.start:h.pos], h.bit, h.width)
	h.bit += uint64(h.width)
	return v, nil
}

// bitsLE reads width bits at bit offset pos of b, least significant bit
// first; bits past the end are zero.
func bitsLE(b []byte, pos uint64, width uint) uint64 {
	var v uint64
	for i := range width {
		p := pos + uint64(i)
		if p>>3 < uint64(len(b)) && b[p>>3]>>(p&7)&1 != 0 {
			v |= 1 << i
		}
	}
	return v
}

// plain decodes PLAIN values of a physical type into Go values: bool,
// int64 (INT32 and INT64), float64, []byte, and the twelve bytes of an
// INT96.
type plain struct {
	b    []byte
	pos  int
	bit  uint64
	phys int32
	size int
}

func (d *plain) next() (any, error) {
	switch d.phys {
	case typeBoolean:
		if d.bit >= uint64(len(d.b))*8 {
			return nil, errCorrupt
		}
		v := d.b[d.bit>>3]>>(d.bit&7)&1 != 0
		d.bit++
		return v, nil
	case typeInt32:
		if len(d.b)-d.pos < 4 {
			return nil, errCorrupt
		}
		v := int64(int32(binary.LittleEndian.Uint32(d.b[d.pos:])))
		d.pos += 4
		return v, nil
	case typeFloat:
		if len(d.b)-d.pos < 4 {
			return nil, errCorrupt
		}
		v := float64(math.Float32frombits(binary.LittleEndian.Uint32(d.b[d.pos:])))
		d.pos += 4
		return v, nil
	case typeInt64:
		if len(d.b)-d.pos < 8 {
			return nil, errCorrupt
		}
		v := int64(binary.LittleEndian.Uint64(d.b[d.pos:]))
		d.pos += 8
		return v, nil
	case typeDouble:
		if len(d.b)-d.pos < 8 {
			return nil, errCorrupt
		}
		v := math.Float64frombits(binary.LittleEndian.Uint64(d.b[d.pos:]))
		d.pos += 8
		return v, nil
	case typeInt96:
		if len(d.b)-d.pos < 12 {
			return nil, errCorrupt
		}
		v := [12]byte(d.b[d.pos : d.pos+12])
		d.pos += 12
		return v, nil
	case typeBinary:
		if len(d.b)-d.pos < 4 {
			return nil, errCorrupt
		}
		n := binary.LittleEndian.Uint32(d.b[d.pos:])
		d.pos += 4
		if uint64(n) > uint64(len(d.b)-d.pos) {
			return nil, errCorrupt
		}
		b := d.b[d.pos : d.pos+int(n) : d.pos+int(n)]
		d.pos += int(n)
		return b, nil
	default:
		if d.size < 0 || len(d.b)-d.pos < d.size {
			return nil, errCorrupt
		}
		b := d.b[d.pos : d.pos+d.size : d.pos+d.size]
		d.pos += d.size
		return b, nil
	}
}

// convert makes the value of a column's type of a decoded value.
func convert(l *leaf, v any) any {
	switch l.typ {
	case String:
		if b, ok := v.([]byte); ok {
			return string(b)
		}
	case Number:
		switch x := v.(type) {
		case int64:
			if l.logical == logDecimal || l.conv == convDecimal {
				return float64(x) / math.Pow10(int(l.scale))
			}
			return float64(x)
		}
	case Date:
		if d, ok := v.(int64); ok {
			return unixDay.AddDate(0, 0, int(d))
		}
	case DateTime:
		switch x := v.(type) {
		case int64:
			switch {
			case l.logical == logTimestamp && l.unit == unitMicros, l.conv == convTSMicros:
				return time.UnixMicro(x).UTC()
			case l.logical == logTimestamp && l.unit == unitNanos:
				return time.Unix(0, x).UTC()
			}
			return time.UnixMilli(x).UTC()
		case [12]byte:
			// the nanoseconds of the day, then the Julian day
			nanos := int64(binary.LittleEndian.Uint64(x[:8]))
			day := int64(binary.LittleEndian.Uint32(x[8:]))
			return unixDay.AddDate(0, 0, int(day-2440588)).Add(time.Duration(nanos)).UTC()
		}
	case Other:
		if b, ok := v.([]byte); ok && (l.logical == logDecimal || l.conv == convDecimal) {
			// a decimal of bytes: a big-endian two's complement integer
			n := new(big.Int).SetBytes(b)
			if len(b) > 0 && b[0]&0x80 != 0 {
				n.Sub(n, new(big.Int).Lsh(big.NewInt(1), uint(len(b)*8)))
			}
			f, _ := new(big.Float).SetInt(n).Float64()
			return f / math.Pow10(int(l.scale))
		}
	}
	return v
}

package parquet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// A column chunk is a dictionary page (when the values are dictionary
// encoded) and data pages, each a Thrift header and the page's bytes,
// compressed with the chunk's codec. A data page holds the repetition
// levels, the definition levels and the values that are not null; in the
// first version of data pages they are compressed together, in the second
// only the values are, and the levels are in the hybrid encoding without
// a length.

// maxHeaderSize bounds the size of a page header (with its statistics).
const maxHeaderSize = 16 << 20

// column reads the entries of a column chunk: the levels of each value,
// and the value when it is not null.
type column struct {
	r          io.ReaderAt
	leaf       *node
	codec      int32
	pos, end   int64
	valuesLeft int64
	dict       []scalar
	// the current page
	left       int
	reps, defs levelDecoder
	vals       valueDecoder
}

var errColumnEnd = errors.New("the column chunk ends early")

func newColumn(r io.ReaderAt, dataEnd int64, leaf *node, cc *columnChunk) (*column, error) {
	switch {
	case cc.encrypted:
		return nil, errors.New("the column is encrypted")
	case cc.filePath != "":
		return nil, fmt.Errorf("the column is in another file (%s)", cc.filePath)
	case cc.meta == nil:
		return nil, errors.New("the column chunk has no metadata")
	case cc.meta.typ != leaf.phys:
		return nil, errors.New("the column chunk's type is not the schema's")
	}
	m := cc.meta
	start := m.dataPageOffset
	if m.dictPageOffset > 0 && m.dictPageOffset < start {
		start = m.dictPageOffset
	}
	if start < 4 || start >= dataEnd || m.totalCompressed < 0 {
		return nil, errors.New("the column chunk is out of the file")
	}
	// parquet-mr 1.2.8 and earlier left the dictionary page's header out
	// of the chunk's size: allow for it (the values' count ends the chunk)
	end := min(start+min(m.totalCompressed, dataEnd)+100, dataEnd)
	return &column{r: r, leaf: leaf, codec: m.codec, pos: start, end: end, valuesLeft: m.numValues}, nil
}

// next returns the next entry: its repetition and definition levels, and
// its value when the definition level is the leaf's (it is not null).
func (c *column) next() (rep, def int, v scalar, err error) {
	for c.left == 0 {
		if c.valuesLeft <= 0 {
			return 0, 0, v, io.EOF
		}
		if err := c.readPage(); err != nil {
			return 0, 0, v, err
		}
	}
	c.left--
	if c.reps != nil {
		r, err := c.reps.next()
		if err != nil || r > uint64(c.leaf.repLevel) {
			return 0, 0, v, errEncoding
		}
		rep = int(r)
	}
	if c.defs != nil {
		d, err := c.defs.next()
		if err != nil || d > uint64(c.leaf.defLevel) {
			return 0, 0, v, errEncoding
		}
		def = int(d)
	}
	if def == c.leaf.defLevel {
		v, err = c.vals.next()
	}
	return rep, def, v, err
}

// readPage reads the next page: the dictionary, or a data page to read
// entries from.
func (c *column) readPage() error {
	h, data, err := c.readRaw()
	if err != nil {
		return err
	}
	switch h.typ {
	case pageDictionary:
		if c.dict != nil {
			return errors.New("a second dictionary page")
		}
		if h.encoding != encPlain && h.encoding != encPlainDict {
			return fmt.Errorf("a dictionary page of encoding %d", h.encoding)
		}
		buf, err := decompress(c.codec, data, int(h.uncompressed))
		if err != nil {
			return err
		}
		n := int(h.dictValues)
		if n < 0 || n > len(buf)*8+1 {
			return errEncoding
		}
		dec := &plainDecoder{b: buf, phys: c.leaf.phys, size: int(c.leaf.typeLen)}
		c.dict = make([]scalar, 0, min(n, len(buf)+1))
		for range n {
			s, err := dec.next()
			if err != nil {
				return err
			}
			c.dict = append(c.dict, s)
		}
		return nil
	case pageData:
		return c.dataPage(h, data)
	case pageDataV2:
		return c.dataPageV2(h, data)
	}
	return nil // index pages and pages of types to come
}

// readRaw reads a page's header and its bytes.
func (c *column) readRaw() (*pageHeader, []byte, error) {
	if c.pos >= c.end {
		return nil, nil, errColumnEnd
	}
	n := min(c.end-c.pos, 256)
	var h *pageHeader
	var hl int
	for {
		buf := make([]byte, n)
		if _, err := c.r.ReadAt(buf, c.pos); err != nil && err != io.EOF {
			return nil, nil, err
		}
		var err error
		if h, hl, err = readPageHeader(buf); err == nil {
			break
		}
		if n == c.end-c.pos || n >= maxHeaderSize {
			return nil, nil, fmt.Errorf("page header: %w", err)
		}
		n = min(n*8, c.end-c.pos, maxHeaderSize)
	}
	c.pos += int64(hl)
	if h.compressed < 0 || h.compressed > maxPageSize || int64(h.compressed) > c.end-c.pos {
		return nil, nil, errColumnEnd
	}
	data := make([]byte, h.compressed)
	if _, err := c.r.ReadAt(data, c.pos); err != nil && !(err == io.EOF && len(data) == 0) {
		return nil, nil, err
	}
	c.pos += int64(h.compressed)
	return h, data, nil
}

func (c *column) dataPage(h *pageHeader, data []byte) error {
	buf, err := decompress(c.codec, data, int(h.uncompressed))
	if err != nil {
		return err
	}
	n := int(h.numValues)
	if n < 0 {
		return errEncoding
	}
	p := 0
	c.reps, c.defs = nil, nil
	if c.leaf.repLevel > 0 {
		if c.reps, p, err = levelsV1(buf, p, h.repEncoding, c.leaf.repLevel, n); err != nil {
			return err
		}
	}
	if c.leaf.defLevel > 0 {
		if c.defs, p, err = levelsV1(buf, p, h.defEncoding, c.leaf.defLevel, n); err != nil {
			return err
		}
	}
	if c.vals, err = c.values(h.encoding, buf[p:]); err != nil {
		return err
	}
	c.left = n
	c.valuesLeft -= int64(n)
	return nil
}

func levelsV1(buf []byte, p int, enc int32, max, n int) (levelDecoder, int, error) {
	width := levelWidth(max)
	switch enc {
	case encRLE:
		if len(buf)-p < 4 {
			return nil, 0, errEncoding
		}
		l := binary.LittleEndian.Uint32(buf[p:])
		if uint64(l) > uint64(len(buf)-p-4) {
			return nil, 0, errEncoding
		}
		return newHybrid(buf[p+4:p+4+int(l)], width), p + 4 + int(l), nil
	case encBitPacked:
		l := (uint64(n)*uint64(width) + 7) / 8
		if l > uint64(len(buf)-p) {
			return nil, 0, errEncoding
		}
		return &bitPacked{b: buf[p : p+int(l)], width: width}, p + int(l), nil
	}
	return nil, 0, fmt.Errorf("levels of encoding %d", enc)
}

func (c *column) dataPageV2(h *pageHeader, data []byte) error {
	rl, dl := int64(h.repLength), int64(h.defLength)
	if rl < 0 || dl < 0 || rl+dl > int64(len(data)) || h.numValues < 0 {
		return errEncoding
	}
	c.reps, c.defs = nil, nil
	if c.leaf.repLevel > 0 {
		c.reps = newHybrid(data[:rl], levelWidth(c.leaf.repLevel))
	}
	if c.leaf.defLevel > 0 {
		c.defs = newHybrid(data[rl:rl+dl], levelWidth(c.leaf.defLevel))
	}
	vals := data[rl+dl:]
	// a page of nulls may have no values, and no bytes for them either
	if size := int64(h.uncompressed) - rl - dl; size == 0 {
		vals = nil
	} else if !h.uncompressV2 && c.codec != codecNone {
		var err error
		if vals, err = decompress(c.codec, vals, int(size)); err != nil {
			return err
		}
	}
	var err error
	if c.vals, err = c.values(h.encoding, vals); err != nil {
		return err
	}
	c.left = int(h.numValues)
	c.valuesLeft -= int64(h.numValues)
	return nil
}

var encodingNames = map[int32]string{encPlain: "PLAIN", encPlainDict: "PLAIN_DICTIONARY", encRLE: "RLE", encBitPacked: "BIT_PACKED",
	encDeltaBinary: "DELTA_BINARY_PACKED", encDeltaLength: "DELTA_LENGTH_BYTE_ARRAY", encDeltaByteArray: "DELTA_BYTE_ARRAY",
	encRLEDict: "RLE_DICTIONARY", encByteStreamSplt: "BYTE_STREAM_SPLIT", 10: "ALP"}

// values returns the decoder of a page's values.
func (c *column) values(enc int32, b []byte) (valueDecoder, error) {
	phys := c.leaf.phys
	switch {
	case enc == encPlain:
		return &plainDecoder{b: b, phys: phys, size: int(c.leaf.typeLen)}, nil
	case enc == encPlainDict || enc == encRLEDict:
		if c.dict == nil {
			return nil, errors.New("dictionary-encoded values without a dictionary page")
		}
		return newDictDecoder(b, c.dict)
	case enc == encRLE && phys == typeBoolean:
		return newBoolRLE(b)
	case enc == encDeltaBinary && (phys == typeInt32 || phys == typeInt64):
		return newDeltaBinary(b, phys == typeInt32)
	case enc == encDeltaLength && phys == typeBinary:
		d, _, err := newDeltaLength(b)
		return d, err
	case enc == encDeltaByteArray && (phys == typeBinary || phys == typeFixed):
		return newDeltaByteArray(b)
	case enc == encByteStreamSplt:
		return newByteStreamSplit(b, phys, int(c.leaf.typeLen))
	}
	if name, ok := encodingNames[enc]; ok {
		return nil, fmt.Errorf("the %s encoding is not supported for this column", name)
	}
	return nil, fmt.Errorf("unknown encoding %d", enc)
}

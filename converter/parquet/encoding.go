package parquet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
)

// The encodings of levels and values (Encodings.md in parquet-format).
// Each decoder streams its values, so that a page is read only as far as
// the rows shown need, and a small page that says it holds billions of
// values costs nothing until they are read.

var errEncoding = errors.New("malformed encoded data")

// scalar is a value of a column: the bits of booleans, integers and
// floating-point numbers in u (int32 values sign-extended), the twelve
// bytes of an INT96 in u (the low eight) and x, and the bytes of byte
// arrays in b.
type scalar struct {
	u uint64
	x uint32
	b []byte
}

type valueDecoder interface {
	next() (scalar, error)
}

// bitsLE reads width bits (up to 64) at bit offset pos of b, least
// significant bit first; bits past the end of b are zero.
func bitsLE(b []byte, pos uint64, width uint) uint64 {
	if width == 0 {
		return 0
	}
	i := pos >> 3
	shift := uint(pos & 7)
	var lo uint64
	if i+8 <= uint64(len(b)) {
		lo = binary.LittleEndian.Uint64(b[i:])
	} else {
		for k := uint64(0); k < 8 && i+k < uint64(len(b)); k++ {
			lo |= uint64(b[i+k]) << (8 * k)
		}
	}
	v := lo >> shift
	if shift+width > 64 && i+8 < uint64(len(b)) {
		v |= uint64(b[i+8]) << (64 - shift)
	}
	return v & mask(width)
}

func mask(width uint) uint64 {
	if width >= 64 {
		return ^uint64(0)
	}
	return 1<<width - 1
}

// hybrid decodes the RLE/bit-packed hybrid of levels, dictionary indices
// and booleans: runs of a repeated value, and runs of values packed in
// width bits each.
type hybrid struct {
	b      []byte
	pos    int
	width  uint
	left   int  // values left in the run
	packed bool // a bit-packed run
	value  uint64
	start  int    // of the packed run's bytes
	bit    uint64 // of the next packed value
}

func newHybrid(b []byte, width uint) *hybrid { return &hybrid{b: b, width: width} }

func (h *hybrid) next() (uint64, error) {
	for h.left == 0 {
		if h.pos >= len(h.b) {
			return 0, errEncoding
		}
		hdr, n := binary.Uvarint(h.b[h.pos:])
		if n <= 0 {
			return 0, errEncoding
		}
		h.pos += n
		if hdr&1 == 1 {
			groups := hdr >> 1
			size := groups * uint64(h.width) // bytes: 8 values per group
			avail := uint64(len(h.b) - h.pos)
			count := groups * 8
			if size > avail {
				// a truncated last run: the values its bytes hold
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
				return 0, errEncoding
			}
			nb := int(h.width+7) / 8
			if nb > len(h.b)-h.pos {
				return 0, errEncoding
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

// bitPacked decodes the deprecated BIT_PACKED encoding of levels: values
// packed most significant bit first.
type bitPacked struct {
	b     []byte
	width uint
	bit   uint64
}

func (p *bitPacked) next() (uint64, error) {
	if (p.bit+uint64(p.width)+7)/8 > uint64(len(p.b)) {
		return 0, errEncoding
	}
	var v uint64
	for range p.width {
		v = v<<1 | uint64(p.b[p.bit>>3]>>(7-p.bit&7)&1)
		p.bit++
	}
	return v, nil
}

type levelDecoder interface {
	next() (uint64, error)
}

// levelWidth is the bit width of levels up to max.
func levelWidth(max int) uint { return uint(bits.Len(uint(max))) }

// plainDecoder decodes PLAIN values: booleans a bit each, numbers little
// endian, byte arrays after their length, fixed-length byte arrays as
// they are.
type plainDecoder struct {
	b    []byte
	pos  int
	bit  uint64 // of booleans
	phys int32
	size int // of fixed-length byte arrays
}

func (d *plainDecoder) next() (scalar, error) {
	switch d.phys {
	case typeBoolean:
		if d.bit >= uint64(len(d.b))*8 {
			return scalar{}, errEncoding
		}
		v := uint64(d.b[d.bit>>3]>>(d.bit&7)) & 1
		d.bit++
		return scalar{u: v}, nil
	case typeInt32, typeFloat:
		if len(d.b)-d.pos < 4 {
			return scalar{}, errEncoding
		}
		v := binary.LittleEndian.Uint32(d.b[d.pos:])
		d.pos += 4
		if d.phys == typeInt32 {
			return scalar{u: uint64(int64(int32(v)))}, nil
		}
		return scalar{u: uint64(v)}, nil
	case typeInt64, typeDouble:
		if len(d.b)-d.pos < 8 {
			return scalar{}, errEncoding
		}
		v := binary.LittleEndian.Uint64(d.b[d.pos:])
		d.pos += 8
		return scalar{u: v}, nil
	case typeInt96:
		if len(d.b)-d.pos < 12 {
			return scalar{}, errEncoding
		}
		s := scalar{u: binary.LittleEndian.Uint64(d.b[d.pos:]), x: binary.LittleEndian.Uint32(d.b[d.pos+8:])}
		d.pos += 12
		return s, nil
	case typeBinary:
		if len(d.b)-d.pos < 4 {
			return scalar{}, errEncoding
		}
		n := binary.LittleEndian.Uint32(d.b[d.pos:])
		d.pos += 4
		if uint64(n) > uint64(len(d.b)-d.pos) {
			return scalar{}, errEncoding
		}
		b := d.b[d.pos : d.pos+int(n) : d.pos+int(n)]
		d.pos += int(n)
		return scalar{b: b}, nil
	default:
		if len(d.b)-d.pos < d.size {
			return scalar{}, errEncoding
		}
		b := d.b[d.pos : d.pos+d.size : d.pos+d.size]
		d.pos += d.size
		return scalar{b: b}, nil
	}
}

// dictDecoder decodes dictionary indices (PLAIN_DICTIONARY and
// RLE_DICTIONARY): a byte of their bit width, then the hybrid encoding.
type dictDecoder struct {
	idx  *hybrid
	dict []scalar
}

func newDictDecoder(b []byte, dict []scalar) (*dictDecoder, error) {
	if len(b) == 0 {
		// a page of nulls has no indices
		return &dictDecoder{idx: newHybrid(nil, 0), dict: dict}, nil
	}
	if b[0] > 32 {
		return nil, errEncoding
	}
	return &dictDecoder{idx: newHybrid(b[1:], uint(b[0])), dict: dict}, nil
}

func (d *dictDecoder) next() (scalar, error) {
	i, err := d.idx.next()
	if err != nil {
		return scalar{}, err
	}
	if i >= uint64(len(d.dict)) {
		return scalar{}, fmt.Errorf("dictionary index %d out of %d values", i, len(d.dict))
	}
	return d.dict[i], nil
}

// boolRLE decodes booleans of the RLE encoding: the hybrid encoding in
// one bit, after its length.
type boolRLE struct{ h *hybrid }

func newBoolRLE(b []byte) (*boolRLE, error) {
	if len(b) < 4 {
		return nil, errEncoding
	}
	n := binary.LittleEndian.Uint32(b)
	if uint64(n) > uint64(len(b)-4) {
		return nil, errEncoding
	}
	return &boolRLE{newHybrid(b[4:4+n], 1)}, nil
}

func (d *boolRLE) next() (scalar, error) {
	v, err := d.h.next()
	return scalar{u: v}, err
}

// deltaBinary decodes DELTA_BINARY_PACKED integers: a header (the values
// in a block, the miniblocks in a block, the number of values and the
// first value), then blocks of the smallest delta and miniblocks of the
// deltas above it, each packed in its own bit width.
type deltaBinary struct {
	b          []byte
	pos        int
	blockSize  uint64
	miniblocks uint64
	perMini    uint64
	total      uint64
	read       uint64
	last       uint64
	i32        bool
	minDelta   uint64
	widths     []byte
	mini       uint64 // the current miniblock of the block
	miniLeft   uint64 // values left in it
	miniStart  int
	bit        uint64
}

func newDeltaBinary(b []byte, i32 bool) (*deltaBinary, error) {
	d := &deltaBinary{b: b, i32: i32}
	var ok bool
	if d.blockSize, ok = d.uvarint(); !ok {
		return nil, errEncoding
	}
	if d.miniblocks, ok = d.uvarint(); !ok {
		return nil, errEncoding
	}
	if d.total, ok = d.uvarint(); !ok {
		return nil, errEncoding
	}
	first, ok := d.uvarint()
	if !ok || d.blockSize == 0 || d.blockSize%128 != 0 || d.blockSize > 1<<20 ||
		d.miniblocks == 0 || d.blockSize%d.miniblocks != 0 || (d.blockSize/d.miniblocks)%32 != 0 {
		return nil, errEncoding
	}
	d.perMini = d.blockSize / d.miniblocks
	d.last = uint64(int64(first>>1) ^ -int64(first&1))
	return d, nil
}

func (d *deltaBinary) uvarint() (uint64, bool) {
	v, n := binary.Uvarint(d.b[d.pos:])
	if n <= 0 {
		return 0, false
	}
	d.pos += n
	return v, true
}

func (d *deltaBinary) nextInt() (uint64, error) {
	if d.read >= d.total {
		return 0, errEncoding
	}
	d.read++
	if d.read == 1 {
		return d.last, nil
	}
	for d.miniLeft == 0 {
		if d.widths == nil || d.mini+1 >= d.miniblocks {
			// a new block
			md, ok := d.uvarint()
			if !ok || uint64(len(d.b)-d.pos) < d.miniblocks {
				return 0, errEncoding
			}
			d.minDelta = uint64(int64(md>>1) ^ -int64(md&1))
			d.widths = d.b[d.pos : d.pos+int(d.miniblocks)]
			d.pos += int(d.miniblocks)
			d.mini = 0
		} else {
			d.mini++
		}
		w := uint64(d.widths[d.mini])
		if w > 64 {
			return 0, errEncoding
		}
		d.miniStart, d.bit, d.miniLeft = d.pos, 0, d.perMini
		size := d.perMini * w / 8
		if size > uint64(len(d.b)-d.pos) {
			// the last miniblock may lack its padding
			size = uint64(len(d.b) - d.pos)
		}
		d.pos += int(size)
	}
	w := uint(d.widths[d.mini])
	delta := bitsLE(d.b[d.miniStart:d.pos], d.bit, w)
	d.bit += uint64(w)
	d.miniLeft--
	d.last += d.minDelta + delta
	if d.i32 {
		d.last = uint64(int64(int32(d.last)))
	}
	return d.last, nil
}

func (d *deltaBinary) next() (scalar, error) {
	v, err := d.nextInt()
	return scalar{u: v}, err
}

// deltaEnd returns where the DELTA_BINARY_PACKED values at the start of b
// end, without decoding them.
func deltaEnd(b []byte) (int, error) {
	d, err := newDeltaBinary(b, false)
	if err != nil {
		return 0, err
	}
	left := uint64(0)
	if d.total > 0 {
		left = d.total - 1
	}
	for left > 0 {
		if _, ok := d.uvarint(); !ok || uint64(len(b)-d.pos) < d.miniblocks {
			return 0, errEncoding
		}
		widths := b[d.pos : d.pos+int(d.miniblocks)]
		d.pos += int(d.miniblocks)
		for _, w := range widths {
			if left == 0 {
				break
			}
			if w > 64 {
				return 0, errEncoding
			}
			size := d.perMini * uint64(w) / 8
			if size > uint64(len(b)-d.pos) {
				if left > d.perMini {
					return 0, errEncoding
				}
				size = uint64(len(b) - d.pos)
			}
			d.pos += int(size)
			left -= min(left, d.perMini)
		}
	}
	return d.pos, nil
}

// deltaLength decodes DELTA_LENGTH_BYTE_ARRAY: the lengths of the byte
// arrays (DELTA_BINARY_PACKED), then their bytes.
type deltaLength struct {
	lengths *deltaBinary
	data    []byte
	pos     int
}

func newDeltaLength(b []byte) (*deltaLength, int, error) {
	end, err := deltaEnd(b)
	if err != nil {
		return nil, 0, err
	}
	lengths, err := newDeltaBinary(b[:end], true)
	if err != nil {
		return nil, 0, err
	}
	return &deltaLength{lengths: lengths, data: b[end:]}, end, nil
}

func (d *deltaLength) next() (scalar, error) {
	n, err := d.lengths.nextInt()
	if err != nil {
		return scalar{}, err
	}
	if int64(n) < 0 || int64(n) > int64(len(d.data)-d.pos) {
		return scalar{}, errEncoding
	}
	b := d.data[d.pos : d.pos+int(n) : d.pos+int(n)]
	d.pos += int(n)
	return scalar{b: b}, nil
}

// deltaByteArray decodes DELTA_BYTE_ARRAY: the lengths of the prefixes
// each value shares with the one before (DELTA_BINARY_PACKED), then the
// rest of the values (DELTA_LENGTH_BYTE_ARRAY).
type deltaByteArray struct {
	prefixes *deltaBinary
	suffixes *deltaLength
	prev     []byte
}

func newDeltaByteArray(b []byte) (*deltaByteArray, error) {
	end, err := deltaEnd(b)
	if err != nil {
		return nil, err
	}
	prefixes, err := newDeltaBinary(b[:end], true)
	if err != nil {
		return nil, err
	}
	suffixes, _, err := newDeltaLength(b[end:])
	if err != nil {
		return nil, err
	}
	return &deltaByteArray{prefixes: prefixes, suffixes: suffixes}, nil
}

func (d *deltaByteArray) next() (scalar, error) {
	p, err := d.prefixes.nextInt()
	if err != nil {
		return scalar{}, err
	}
	s, err := d.suffixes.next()
	if err != nil {
		return scalar{}, err
	}
	if int64(p) < 0 || int64(p) > int64(len(d.prev)) {
		return scalar{}, errEncoding
	}
	v := make([]byte, int(p)+len(s.b))
	copy(v, d.prev[:p])
	copy(v[p:], s.b)
	d.prev = v
	return scalar{b: v}, nil
}

// byteStreamSplit decodes BYTE_STREAM_SPLIT: the first bytes of all the
// values, then their second bytes, and so on.
type byteStreamSplit struct {
	b    []byte
	n, k int // values, bytes per value
	i    int
	phys int32
}

func newByteStreamSplit(b []byte, phys int32, size int) (*byteStreamSplit, error) {
	k := size
	switch phys {
	case typeInt32, typeFloat:
		k = 4
	case typeInt64, typeDouble:
		k = 8
	case typeFixed:
	default:
		return nil, errEncoding
	}
	if k <= 0 || len(b)%k != 0 {
		return nil, errEncoding
	}
	return &byteStreamSplit{b: b, n: len(b) / k, k: k, phys: phys}, nil
}

func (d *byteStreamSplit) next() (scalar, error) {
	if d.i >= d.n {
		return scalar{}, errEncoding
	}
	v := make([]byte, d.k)
	for j := range d.k {
		v[j] = d.b[j*d.n+d.i]
	}
	d.i++
	switch d.phys {
	case typeInt32:
		return scalar{u: uint64(int64(int32(binary.LittleEndian.Uint32(v))))}, nil
	case typeFloat:
		return scalar{u: uint64(binary.LittleEndian.Uint32(v))}, nil
	case typeInt64, typeDouble:
		return scalar{u: binary.LittleEndian.Uint64(v)}, nil
	}
	return scalar{b: v}, nil
}

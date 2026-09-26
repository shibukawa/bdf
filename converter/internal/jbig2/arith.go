package jbig2

// qeRow is a row of the probability estimation table (T.88 Table E.1).
type qeRow struct {
	qe         uint32
	nmps, nlps uint8
	switchMPS  bool
}

var qeTable = [47]qeRow{
	{0x5601, 1, 1, true}, {0x3401, 2, 6, false}, {0x1801, 3, 9, false}, {0x0ac1, 4, 12, false},
	{0x0521, 5, 29, false}, {0x0221, 38, 33, false}, {0x5601, 7, 6, true}, {0x5401, 8, 14, false},
	{0x4801, 9, 14, false}, {0x3801, 10, 14, false}, {0x3001, 11, 17, false}, {0x2401, 12, 18, false},
	{0x1c01, 13, 20, false}, {0x1601, 29, 21, false}, {0x5601, 15, 14, true}, {0x5401, 16, 14, false},
	{0x5101, 17, 15, false}, {0x4801, 18, 16, false}, {0x3801, 19, 17, false}, {0x3401, 20, 18, false},
	{0x3001, 21, 19, false}, {0x2801, 22, 19, false}, {0x2401, 23, 20, false}, {0x2201, 24, 21, false},
	{0x1c01, 25, 22, false}, {0x1801, 26, 23, false}, {0x1601, 27, 24, false}, {0x1401, 28, 25, false},
	{0x1201, 29, 26, false}, {0x1101, 30, 27, false}, {0x0ac1, 31, 28, false}, {0x09c1, 32, 29, false},
	{0x08a1, 33, 30, false}, {0x0521, 34, 31, false}, {0x0441, 35, 32, false}, {0x02a1, 36, 33, false},
	{0x0221, 37, 34, false}, {0x0141, 38, 35, false}, {0x0111, 39, 36, false}, {0x0085, 40, 37, false},
	{0x0049, 41, 38, false}, {0x0025, 42, 39, false}, {0x0015, 43, 40, false}, {0x0009, 44, 41, false},
	{0x0005, 45, 42, false}, {0x0001, 45, 43, false}, {0x5601, 46, 46, false},
}

// arithDecoder is the MQ arithmetic decoder (T.88 Annex E). A context is a byte
// holding its state index shifted left by one and its more probable symbol in
// bit 0; the zero value is the initial state. Past the end of the data it reads
// 0xFF bytes, as the standard prescribes.
type arithDecoder struct {
	data []byte
	pos  int
	c, a uint32
	ct   int
	over int // bytes read past the end of the data or at a marker
}

// maxOverrun is how many bytes past the end of its data, or past the marker
// that ends it, the decoder reads before the data counts as exhausted. Coded
// data needs no more than a few; reading on decodes garbage from corrupt or
// truncated data.
const maxOverrun = 16

// exhausted reports whether the decoder has run far past the end of its data.
func (d *arithDecoder) exhausted() bool { return d.over > maxOverrun }

// newArithDecoder starts decoding data (INITDEC).
func newArithDecoder(data []byte) *arithDecoder {
	d := &arithDecoder{data: data}
	d.c = d.byteAt(0) << 16
	d.byteIn()
	d.c <<= 7
	d.ct -= 7
	d.a = 0x8000
	return d
}

func (d *arithDecoder) byteAt(i int) uint32 {
	if i < len(d.data) {
		return uint32(d.data[i])
	}
	return 0xff
}

// byteIn reads the next byte into the code register (BYTEIN).
func (d *arithDecoder) byteIn() {
	if d.byteAt(d.pos) == 0xff {
		if b := d.byteAt(d.pos + 1); b > 0x8f {
			// A marker, or the end of the data: feed 1 bits without advancing.
			d.c += 0xff00
			d.ct = 8
			d.over++
		} else {
			d.pos++
			d.c += b << 9
			d.ct = 7
		}
		return
	}
	d.pos++
	if d.pos >= len(d.data) {
		d.over++
	}
	d.c += d.byteAt(d.pos) << 8
	d.ct = 8
}

// decode decodes a binary decision in context cx (DECODE).
func (d *arithDecoder) decode(cx *uint8) int {
	st := *cx
	row := &qeTable[st>>1]
	mps := int(st & 1)
	qe := row.qe
	d.a -= qe
	var bit int
	if d.c>>16 < qe {
		// The lower subinterval, with conditional exchange (LPS_EXCHANGE).
		if d.a < qe {
			bit = mps
			*cx = row.nmps<<1 | uint8(mps)
		} else {
			bit = 1 - mps
			if row.switchMPS {
				mps = 1 - mps
			}
			*cx = row.nlps<<1 | uint8(mps)
		}
		d.a = qe
	} else {
		d.c -= qe << 16
		if d.a&0x8000 != 0 {
			return mps
		}
		// MPS_EXCHANGE.
		if d.a < qe {
			bit = 1 - mps
			if row.switchMPS {
				mps = 1 - mps
			}
			*cx = row.nlps<<1 | uint8(mps)
		} else {
			bit = mps
			*cx = row.nmps<<1 | uint8(mps)
		}
	}
	for { // RENORMD
		if d.ct == 0 {
			d.byteIn()
		}
		d.a <<= 1
		d.c <<= 1
		d.ct--
		if d.a&0x8000 != 0 {
			return bit
		}
	}
}

// intContext is the context of an integer decoding procedure (IAx).
type intContext [512]uint8

// decodeInt decodes an integer (T.88 A.2); ok is false for the out-of-band value.
// Values beyond 32 bits are clamped.
func (d *arithDecoder) decodeInt(cx *intContext) (v int, ok bool) {
	prev := 1
	bit := func() int {
		b := d.decode(&cx[prev])
		if prev < 256 {
			prev = prev<<1 | b
		} else {
			prev = (prev<<1|b)&511 | 256
		}
		return b
	}
	bits := func(n int) int64 {
		var v int64
		for i := 0; i < n; i++ {
			v = v<<1 | int64(bit())
		}
		return v
	}
	sign := bit()
	var val int64
	switch {
	case bit() == 0:
		val = bits(2)
	case bit() == 0:
		val = bits(4) + 4
	case bit() == 0:
		val = bits(6) + 20
	case bit() == 0:
		val = bits(8) + 84
	case bit() == 0:
		val = bits(12) + 340
	default:
		val = bits(32) + 4436
	}
	if val > 1<<31-1 {
		val = 1<<31 - 1
	}
	if sign == 1 {
		if val == 0 {
			return 0, false
		}
		val = -val
	}
	return int(val), true
}

// decodeIAID decodes a symbol ID of n bits (T.88 A.3); cx has 1<<n entries.
func (d *arithDecoder) decodeIAID(cx []uint8, n int) int {
	prev := 1
	for i := 0; i < n; i++ {
		prev = prev<<1 | d.decode(&cx[prev])
	}
	return prev - 1<<n
}

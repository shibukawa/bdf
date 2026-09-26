package jpx

import "math/bits"

// The MQ arithmetic decoder and the raw (bypass) bit reader of Annex C and
// D.6. Both read a codeword segment as if two 0xFF bytes followed it: the
// pair acts as a marker, so a decoder that runs past the segment is fed 1
// bits, as the standard asks.

// mqEntry is a state of Table C.2 combined with the sense of the MPS; the
// index of an entry is state<<1 | mps, and nmps and nlps are the indices of
// the next entries, with the MPS switch already applied.
type mqEntry struct {
	qe   uint32
	nmps uint8
	nlps uint8
}

var mqTable = buildMQTable()

// buildMQTable fills 94 entries; the rest pad the table so that any uint8
// indexes it without a bounds check.
func buildMQTable() (t [256]mqEntry) {
	states := [47]struct {
		qe         uint16
		nmps, nlps uint8
		sw         bool
	}{
		{0x5601, 1, 1, true}, {0x3401, 2, 6, false}, {0x1801, 3, 9, false},
		{0x0ac1, 4, 12, false}, {0x0521, 5, 29, false}, {0x0221, 38, 33, false},
		{0x5601, 7, 6, true}, {0x5401, 8, 14, false}, {0x4801, 9, 14, false},
		{0x3801, 10, 14, false}, {0x3001, 11, 17, false}, {0x2401, 12, 18, false},
		{0x1c01, 13, 20, false}, {0x1601, 29, 21, false}, {0x5601, 15, 14, true},
		{0x5401, 16, 14, false}, {0x5101, 17, 15, false}, {0x4801, 18, 16, false},
		{0x3801, 19, 17, false}, {0x3401, 20, 18, false}, {0x3001, 21, 19, false},
		{0x2801, 22, 19, false}, {0x2401, 23, 20, false}, {0x2201, 24, 21, false},
		{0x1c01, 25, 22, false}, {0x1801, 26, 23, false}, {0x1601, 27, 24, false},
		{0x1401, 28, 25, false}, {0x1201, 29, 26, false}, {0x1101, 30, 27, false},
		{0x0ac1, 31, 28, false}, {0x09c1, 32, 29, false}, {0x08a1, 33, 30, false},
		{0x0521, 34, 31, false}, {0x0441, 35, 32, false}, {0x02a1, 36, 33, false},
		{0x0221, 37, 34, false}, {0x0141, 38, 35, false}, {0x0111, 39, 36, false},
		{0x0085, 40, 37, false}, {0x0049, 41, 38, false}, {0x0025, 42, 39, false},
		{0x0015, 43, 40, false}, {0x0009, 44, 41, false}, {0x0005, 45, 42, false},
		{0x0001, 45, 43, false}, {0x5601, 46, 46, false},
	}
	for s, st := range states {
		for mps := uint8(0); mps < 2; mps++ {
			lps := mps
			if st.sw {
				lps ^= 1
			}
			t[s<<1|int(mps)] = mqEntry{qe: uint32(st.qe), nmps: st.nmps<<1 | mps, nlps: st.nlps<<1 | lps}
		}
	}
	return t
}

// Context labels of the code-block coder (Annex D): 0-8 significance, 9-13
// sign, 14-16 magnitude refinement, then run-length and uniform.
const (
	ctxRL  = 17
	ctxUni = 18
	numCtx = 19
)

// mqDecoder is the MQ decoder of C.3 with the register conventions of
// Figure C.14: c holds Chigh in its upper 16 bits.
type mqDecoder struct {
	data []byte
	bp   int
	a, c uint32
	ct   uint32
	cx   [32]uint8 // numCtx used; indexed with cx&31
}

// resetContexts puts every context in its initial state (Table D.7).
func (m *mqDecoder) resetContexts() {
	m.cx = [32]uint8{}
	m.cx[0] = 4 << 1
	m.cx[ctxRL] = 3 << 1
	m.cx[ctxUni] = 46 << 1
}

// at returns byte i of the segment, or 0xFF past its end.
func at(data []byte, i int) byte {
	if i < len(data) {
		return data[i]
	}
	return 0xff
}

// init starts decoding a segment (INITDEC).
func (m *mqDecoder) init(data []byte) {
	m.data = data
	m.bp = 0
	m.c = uint32(at(data, 0)) << 16
	m.bytein()
	m.c <<= 7
	m.ct -= 7
	m.a = 0x8000
}

// bytein is BYTEIN of Figure C.18; byte bp is the last byte read.
func (m *mqDecoder) bytein() {
	if at(m.data, m.bp) == 0xff {
		if b := at(m.data, m.bp+1); b > 0x8f {
			m.c += 0xff00
			m.ct = 8
		} else {
			m.bp++
			m.c += uint32(b) << 9
			m.ct = 7
		}
	} else {
		m.bp++
		m.c += uint32(at(m.data, m.bp)) << 8
		m.ct = 8
	}
}

// decode decodes a decision in context cx (DECODE of Figure C.15).
func (m *mqDecoder) decode(cx int) uint32 {
	cx &= 31
	st := m.cx[cx]
	e := &mqTable[st]
	qe := e.qe
	a := m.a - qe
	var d uint32
	if m.c>>16 < qe {
		// LPS_EXCHANGE
		if a < qe {
			d = uint32(st & 1)
			m.cx[cx] = e.nmps
		} else {
			d = uint32(st&1) ^ 1
			m.cx[cx] = e.nlps
		}
		a = qe
	} else {
		m.c -= qe << 16
		if a&0x8000 != 0 {
			m.a = a
			return uint32(st & 1)
		}
		// MPS_EXCHANGE
		if a < qe {
			d = uint32(st&1) ^ 1
			m.cx[cx] = e.nlps
		} else {
			d = uint32(st & 1)
			m.cx[cx] = e.nmps
		}
	}
	// RENORMD, shifting as many bits at a time as the byte buffer holds.
	n := uint32(bits.LeadingZeros16(uint16(a)))
	a <<= n
	c, ct := m.c, m.ct
	for n > ct {
		c <<= ct
		n -= ct
		m.c = c
		m.bytein()
		c, ct = m.c, m.ct
	}
	c <<= n
	ct -= n
	m.a, m.c, m.ct = a, c, ct
	return d
}

// rawDecoder reads the raw bits of the passes that the selective arithmetic
// coding bypass leaves uncoded; a 0 bit is stuffed after every 0xFF byte.
type rawDecoder struct {
	data []byte
	pos  int
	c    uint32
	ct   uint32
}

func (r *rawDecoder) init(data []byte) {
	r.data = data
	r.pos = 0
	r.c = 0
	r.ct = 0
}

func (r *rawDecoder) bit() uint32 {
	if r.ct == 0 {
		if r.c == 0xff {
			if b := at(r.data, r.pos); b > 0x8f {
				r.ct = 8 // at the terminator: keep returning 1 bits
			} else {
				r.c = uint32(b)
				r.pos++
				r.ct = 7
			}
		} else {
			r.c = uint32(at(r.data, r.pos))
			r.pos = min(r.pos+1, len(r.data)+1)
			r.ct = 8
		}
	}
	r.ct--
	return r.c >> r.ct & 1
}

package jbig2

// Segment types (T.88 7.3).
const (
	typeSymbolDict                  = 0
	typeIntermediateText            = 4
	typeImmediateText               = 6
	typeImmediateLosslessText       = 7
	typePatternDict                 = 16
	typeIntermediateHalftone        = 20
	typeImmediateHalftone           = 22
	typeImmediateLosslessHalftone   = 23
	typeIntermediateGeneric         = 36
	typeImmediateGeneric            = 38
	typeImmediateLosslessGeneric    = 39
	typeIntermediateRefinement      = 40
	typeImmediateRefinement         = 42
	typeImmediateLosslessRefinement = 43
	typePageInfo                    = 48
	typeEndOfPage                   = 49
	typeEndOfStripe                 = 50
	typeEndOfFile                   = 51
	typeProfiles                    = 52
	typeTables                      = 53
	typeExtension                   = 62
)

// unknownLength is the data length of a segment whose end is found by scanning (7.2.7).
const unknownLength = 0xffffffff

func isIntermediate(typ int) bool {
	return typ == typeIntermediateText || typ == typeIntermediateHalftone ||
		typ == typeIntermediateGeneric || typ == typeIntermediateRefinement
}

func isGeneric(typ int) bool {
	return typ == typeIntermediateGeneric || typ == typeImmediateGeneric || typ == typeImmediateLosslessGeneric
}

// segmentHeader is a segment header (7.2).
type segmentHeader struct {
	number uint32
	typ    int
	page   uint32   // page association
	refs   []uint32 // referred-to segments
	length uint32   // data length, or unknownLength
}

// parseSegmentHeader parses the segment header at the start of b and returns
// it with its size.
func parseSegmentHeader(b []byte) (h segmentHeader, n int, err error) {
	r := reader{b: b}
	h.number = r.u32()
	flags := r.u8()
	h.typ = int(flags & 0x3f)
	count := int(r.u8() >> 5)
	switch {
	case count == 7:
		r.pos--
		c := r.u32() & 0x1fffffff
		if int64(c) > int64(len(b)) {
			return h, 0, errTruncated // each reference takes at least a byte
		}
		count = int(c)
		r.skip((count + 8) / 8) // retention flags
	case count > 4:
		return h, 0, errorf("segment %d: invalid referred-to segment count", h.number)
	}
	size := 1
	if h.number > 65536 {
		size = 4
	} else if h.number > 256 {
		size = 2
	}
	if r.err == nil && int64(count)*int64(size) > int64(len(b)-r.pos) {
		return h, 0, errTruncated
	}
	if count > 0 {
		h.refs = make([]uint32, count)
	}
	for i := range h.refs {
		switch size {
		case 1:
			h.refs[i] = uint32(r.u8())
		case 2:
			h.refs[i] = uint32(r.u16())
		default:
			h.refs[i] = r.u32()
		}
	}
	if flags&0x40 != 0 {
		h.page = r.u32()
	} else {
		h.page = uint32(r.u8())
	}
	h.length = r.u32()
	if r.err != nil {
		return h, 0, r.err
	}
	return h, r.pos, nil
}

// genericRegionLength finds the end of an immediate generic region segment of
// unknown length (7.2.7): its coded data ends with 0xFF 0xAC (arithmetic) or
// 0x00 0x00 (MMR) followed by a 4-byte row count.
func genericRegionLength(data []byte) (int, bool) {
	if len(data) < regionInfoSize+1 {
		return 0, false
	}
	flags := data[regionInfoSize]
	start := regionInfoSize + 1
	m0, m1 := byte(0xff), byte(0xac)
	if flags&1 != 0 {
		m0, m1 = 0, 0
	} else if flags>>1&3 == 0 {
		start += 8
	} else {
		start += 2
	}
	for i := start; i+6 <= len(data); i++ {
		if data[i] == m0 && data[i+1] == m1 {
			return i + 6, true
		}
	}
	return 0, false
}

// regionInfoSize is the size of the region segment information field.
const regionInfoSize = 17

// regionInfo is the region segment information field (7.4.1).
type regionInfo struct {
	w, h uint32
	x, y uint32
	op   int // external combination operator
}

func readRegionInfo(r *reader) (regionInfo, error) {
	ri := regionInfo{w: r.u32(), h: r.u32(), x: r.u32(), y: r.u32()}
	ri.op = int(r.u8() & 7)
	if r.err != nil {
		return ri, r.err
	}
	if ri.op > opReplace {
		return ri, errorf("invalid combination operator %d", ri.op)
	}
	return ri, nil
}

// size returns the region's width and height, checked against the limits.
func (ri regionInfo) size() (w, h int, err error) {
	if w, err = dim(ri.w); err != nil {
		return
	}
	h, err = dim(ri.h)
	return
}

// point is a pixel offset, such as an adaptive template pixel.
type point struct{ x, y int }

// reader reads big-endian fields. Reading past the end sets err and yields zeros.
type reader struct {
	b   []byte
	pos int
	err error
}

func (r *reader) need(n int) bool {
	if r.err != nil || n > len(r.b)-r.pos {
		r.err = errTruncated
		return false
	}
	return true
}

func (r *reader) skip(n int) {
	if r.need(n) {
		r.pos += n
	}
}

func (r *reader) u8() byte {
	if !r.need(1) {
		return 0
	}
	r.pos++
	return r.b[r.pos-1]
}

func (r *reader) i8() int { return int(int8(r.u8())) }

func (r *reader) u16() uint16 {
	if !r.need(2) {
		return 0
	}
	r.pos += 2
	return uint16(r.b[r.pos-2])<<8 | uint16(r.b[r.pos-1])
}

func (r *reader) u32() uint32 {
	if !r.need(4) {
		return 0
	}
	r.pos += 4
	b := r.b[r.pos-4:]
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// at reads n adaptive template pixels.
func (r *reader) at(n int) []point {
	p := make([]point, n)
	for i := range p {
		p[i].x = r.i8()
		p[i].y = r.i8()
	}
	return p
}

// rest returns the unread bytes.
func (r *reader) rest() []byte {
	if r.err != nil {
		return nil
	}
	return r.b[r.pos:]
}

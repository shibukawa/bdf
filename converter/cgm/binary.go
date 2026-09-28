package cgm

import (
	"encoding/binary"
	"math"

	"github.com/shibukawa/bdf/converter/internal/cad"
)

// The binary encoding (ISO/IEC 8632-3). Every element starts on a 16-bit
// boundary with a word holding its class (4 bits), id (7 bits) and the
// length of its parameters (5 bits); a length of 31 is followed by
// partitions, each with a word of its own length (15 bits) and a flag for
// more partitions. Parameters of odd length are padded to the next word.

// realFormat is how reals are stored: fixed point (a signed whole part and
// an unsigned fraction of half the bits each) or IEEE floating point, 32
// or 64 bits.
type realFormat struct {
	fixed bool
	bits  int
}

// precisions are the sizes of the parameters of the binary encoding, as
// the metafile descriptor and control elements set them.
type precisions struct {
	intBits, indexBits, colourBits, colourIndexBits, nameBits int
	real                                                      realFormat
	vdcReal                                                   bool // VDC TYPE real
	vdcIntBits                                                int
	vdcRealFormat                                             realFormat
	components                                                int // components of a direct colour: 3, or 4 for CMYK
}

func defaultPrecisions() precisions {
	return precisions{intBits: 16, indexBits: 16, colourBits: 8, colourIndexBits: 8, nameBits: 16,
		real: realFormat{fixed: true, bits: 32}, vdcIntBits: 16, vdcRealFormat: realFormat{fixed: true, bits: 32}, components: 3}
}

// errShort is raised (as a panic) by a parameter read past the end of an
// element's parameters; the interpreter recovers it for the element.
type shortError struct{}

func (shortError) Error() string { return "the element is shorter than its parameters" }

var errShort = shortError{}

// binReader splits a binary metafile into elements.
type binReader struct {
	data []byte
	pos  int
	// truncated reports that the last element ran past the end of the file
	truncated bool
}

// next returns the next element; false at the end of the data.
func (r *binReader) next() (element, bool) {
	for {
		if r.pos%2 == 1 {
			r.pos++
		}
		if r.pos+2 > len(r.data) {
			return element{}, false
		}
		h := binary.BigEndian.Uint16(r.data[r.pos:])
		r.pos += 2
		c := code{int(h >> 12), int(h>>5) & 0x7f}
		n := int(h & 0x1f)
		var data []byte
		if n < 31 {
			data = r.take(n)
		} else {
			for {
				if r.pos+2 > len(r.data) {
					r.truncated = true
					break
				}
				w := binary.BigEndian.Uint16(r.data[r.pos:])
				r.pos += 2
				part := r.take(int(w & 0x7fff))
				if data == nil {
					data = part
				} else {
					// the first partition is a piece of the file without
					// room after it, so the first append copies it
					data = append(data, part...)
				}
				if w&0x8000 == 0 || r.truncated {
					break
				}
				if r.pos%2 == 1 {
					r.pos++
				}
			}
		}
		if c == eNoOp {
			continue
		}
		return element{code: c, data: data}, true
	}
}

// take returns the next n bytes, or those left.
func (r *binReader) take(n int) []byte {
	if r.pos+n > len(r.data) {
		r.truncated = true
		n = len(r.data) - r.pos
	}
	b := r.data[r.pos : r.pos+n : r.pos+n]
	r.pos += n
	return b
}

// binParams reads the parameters of a binary element.
type binParams struct {
	b   []byte
	pos int
	pr  *precisions
}

func (p *binParams) more() bool { return p.pos < len(p.b) }

// bytes returns the next n bytes.
func (p *binParams) bytes(n int) []byte {
	if n < 0 || p.pos+n > len(p.b) {
		panic(errShort)
	}
	b := p.b[p.pos : p.pos+n]
	p.pos += n
	return b
}

// uint reads an unsigned integer of bits (a multiple of 8).
func (p *binParams) uint(bits int) uint64 {
	n := bits / 8
	if n < 1 || n > 8 {
		n = 2
	}
	var v uint64
	for _, c := range p.bytes(n) {
		v = v<<8 | uint64(c)
	}
	return v
}

// sint reads a signed integer of bits.
func (p *binParams) sint(bits int) int64 {
	n := bits / 8
	if n < 1 || n > 8 {
		n = 2
	}
	v := p.uint(n * 8)
	shift := 64 - uint(n*8)
	return int64(v<<shift) >> shift
}

func (p *binParams) realOf(f realFormat) float64 {
	switch {
	case f.fixed && f.bits == 64:
		whole := p.sint(32)
		frac := p.uint(32)
		return float64(whole) + float64(frac)/(1<<32)
	case f.fixed:
		whole := p.sint(16)
		frac := p.uint(16)
		return float64(whole) + float64(frac)/(1<<16)
	case f.bits == 64:
		return finite(math.Float64frombits(p.uint(64)))
	default:
		return finite(float64(math.Float32frombits(uint32(p.uint(32)))))
	}
}

func finite(v float64) float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return v
}

func (p *binParams) int() int           { return clampInt(p.sint(p.pr.intBits)) }
func (p *binParams) index() int         { return clampInt(p.sint(p.pr.indexBits)) }
func (p *binParams) name() int          { return clampInt(p.sint(p.pr.nameBits)) }
func (p *binParams) enum(string) int    { return int(p.sint(16)) }
func (p *binParams) real() float64      { return p.realOf(p.pr.real) }
func (p *binParams) colourIndex() int   { return clampInt(int64(p.uint(p.pr.colourIndexBits))) }
func (p *binParams) str() []byte        { return p.string() }
func (p *binParams) rest() []byte       { return p.bytes(len(p.b) - p.pos) }
func (p *binParams) components() int    { return p.pr.components }
func (p *binParams) colourMax() float64 { return math.Exp2(float64(p.pr.colourBits)) - 1 }

func clampInt(v int64) int {
	return int(max(min(v, math.MaxInt32), math.MinInt32))
}

func (p *binParams) vdc() float64 {
	if p.pr.vdcReal {
		return p.realOf(p.pr.vdcRealFormat)
	}
	return float64(p.sint(p.pr.vdcIntBits))
}

func (p *binParams) point() cad.Point {
	x := p.vdc()
	return cad.Point{X: x, Y: p.vdc()}
}

// component reads a component of a direct colour.
func (p *binParams) component() float64 { return float64(p.uint(p.pr.colourBits)) }

// scaleFactor reads the metric scaling factor of SCALING MODE. Part 3 says
// it is a floating point number whatever the real precision, but some
// writers use the real precision: the reading that gives a plausible
// factor wins.
func (p *binParams) scaleFactor() float64 {
	left := len(p.b) - p.pos
	if left < 4 {
		return 0
	}
	bits := p.pr.real.bits
	if left < 8 {
		bits = 32
	}
	start := p.pos
	fl := p.realOf(realFormat{bits: bits})
	p.pos = start
	fx := p.realOf(realFormat{fixed: true, bits: bits})
	plausible := func(v float64) bool { return v > 1e-9 && v < 1e9 }
	if plausible(fl) || !plausible(fx) {
		return fl
	}
	return fx
}

// string reads a string: a length octet, or 255 and words of 15-bit
// lengths with a continuation flag.
func (p *binParams) string() []byte {
	n := int(p.uint(8))
	if n < 255 {
		return p.bytes(n)
	}
	var out []byte
	for {
		w := p.uint(16)
		out = append(out, p.bytes(int(w&0x7fff))...)
		if w&0x8000 == 0 {
			return out
		}
	}
}

// sdr reads a structured data record (the parameters of APPLICATION
// STRUCTURE ATTRIBUTE): encoded as a string, whose content is a list of
// data type, count and values.
func (p *binParams) sdr() []sdrMember {
	start := p.pos
	defer func() {
		if r := recover(); r != nil && r != errShort {
			panic(r)
		}
	}()
	body := p.string()
	if p.pos != len(p.b) {
		// not a string: the members follow as they are
		p.pos = start
		body = p.rest()
	}
	q := &binParams{b: body, pr: p.pr}
	var out []sdrMember
	for q.more() {
		m := sdrMember{typ: q.index()}
		n := q.int()
		for i := 0; i < n && q.more(); i++ {
			m.values = append(m.values, q.sdrValue(m.typ))
		}
		out = append(out, m)
	}
	return out
}

// sdrValue reads a value of an SDR data type (ISO/IEC 8632-1 table 7).
func (p *binParams) sdrValue(typ int) any {
	switch typ {
	case 2: // CI
		return p.colourIndex()
	case 3: // CD
		c := make([]float64, p.pr.components)
		for i := range c {
			c[i] = p.component()
		}
		return c
	case 4: // N
		return p.name()
	case 5: // E
		return p.enum("")
	case 6: // I
		return p.int()
	case 8: // IF8
		return int(p.sint(8))
	case 9: // IF16
		return int(p.sint(16))
	case 10: // IF32
		return clampInt(p.sint(32))
	case 11: // IX
		return p.index()
	case 12: // R
		return p.real()
	case 13, 14: // S, SF
		return p.string()
	case 15: // VC
		return p.real()
	case 16: // VDC
		return p.vdc()
	case 17: // CCO
		return p.component()
	case 18: // UI8
		return int(p.uint(8))
	case 19: // UI32
		return clampInt(int64(p.uint(32)))
	case 22: // UI16
		return int(p.uint(16))
	}
	// data types without a size of their own end the record
	p.pos = len(p.b)
	return nil
}

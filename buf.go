package bdf

import (
	"encoding/binary"
	"math"
)

// buf is a little-endian byte buffer with the primitive encoders used by the format.
type buf struct {
	b []byte
}

func (w *buf) u8(v byte)    { w.b = append(w.b, v) }
func (w *buf) u16(v uint16) { w.b = binary.LittleEndian.AppendUint16(w.b, v) }
func (w *buf) u32(v uint32) { w.b = binary.LittleEndian.AppendUint32(w.b, v) }
func (w *buf) u64(v uint64) { w.b = binary.LittleEndian.AppendUint64(w.b, v) }
func (w *buf) f32(v float32) {
	w.b = binary.LittleEndian.AppendUint32(w.b, math.Float32bits(v))
}
func (w *buf) f32s(vs ...float32) {
	for _, v := range vs {
		w.f32(v)
	}
}
func (w *buf) varuint(v uint64) { w.b = binary.AppendUvarint(w.b, v) }
func (w *buf) bytes(p []byte)   { w.b = append(w.b, p...) }
func (w *buf) str(s string) {
	w.varuint(uint64(len(s)))
	w.b = append(w.b, s...)
}
func (w *buf) hash(h Hash) { w.b = append(w.b, h[:]...) }

// reader is the decoding counterpart of buf.
type reader struct {
	b   []byte
	pos int
	err error
}

func (r *reader) fail(msg string) {
	if r.err == nil {
		r.err = &FormatError{Offset: r.pos, Msg: msg}
	}
}

func (r *reader) need(n int) bool {
	if r.err != nil {
		return false
	}
	if r.pos+n > len(r.b) {
		r.fail("unexpected end of data")
		return false
	}
	return true
}

func (r *reader) eof() bool { return r.err != nil || r.pos >= len(r.b) }

func (r *reader) u8() byte {
	if !r.need(1) {
		return 0
	}
	v := r.b[r.pos]
	r.pos++
	return v
}

func (r *reader) u16() uint16 {
	if !r.need(2) {
		return 0
	}
	v := binary.LittleEndian.Uint16(r.b[r.pos:])
	r.pos += 2
	return v
}

func (r *reader) u32() uint32 {
	if !r.need(4) {
		return 0
	}
	v := binary.LittleEndian.Uint32(r.b[r.pos:])
	r.pos += 4
	return v
}

func (r *reader) u64() uint64 {
	if !r.need(8) {
		return 0
	}
	v := binary.LittleEndian.Uint64(r.b[r.pos:])
	r.pos += 8
	return v
}

func (r *reader) f32() float32 { return math.Float32frombits(r.u32()) }

func (r *reader) varuint() uint64 {
	if r.err != nil {
		return 0
	}
	v, n := binary.Uvarint(r.b[r.pos:])
	if n <= 0 {
		r.fail("bad varuint")
		return 0
	}
	r.pos += n
	return v
}

func (r *reader) bytes(n int) []byte {
	if !r.need(n) {
		return nil
	}
	v := r.b[r.pos : r.pos+n]
	r.pos += n
	return v
}

func (r *reader) str() string {
	n := r.varuint()
	if n > uint64(len(r.b)) {
		r.fail("bad string length")
		return ""
	}
	return string(r.bytes(int(n)))
}

func (r *reader) hash() Hash {
	var h Hash
	copy(h[:], r.bytes(16))
	return h
}

// FormatError reports malformed data.
type FormatError struct {
	Offset int
	Msg    string
}

func (e *FormatError) Error() string { return "bdf: " + e.Msg }

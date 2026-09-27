package jww

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"golang.org/x/text/encoding/japanese"
)

// archive reads what an MFC CArchive wrote (MFC technical note TN002):
// little-endian numbers, length-prefixed strings in Shift_JIS, and objects
// preceded by tags that name their class the first time it appears and
// refer back to it later.
type archive struct {
	b   []byte
	i   int
	err error
	// classes and objects share one index; 0 stands for null
	classes map[uint32]string
	objects map[uint32]any
	next    uint32
}

var errShort = errors.New("the file ends in the middle of the data")

func newArchive(b []byte) *archive {
	return &archive{b: b, classes: map[uint32]string{}, objects: map[uint32]any{}, next: 1}
}

func (a *archive) need(n int) bool {
	if a.err != nil {
		return false
	}
	if n < 0 || a.i+n > len(a.b) {
		a.err = errShort
		return false
	}
	return true
}

func (a *archive) u8() uint8 {
	if !a.need(1) {
		return 0
	}
	a.i++
	return a.b[a.i-1]
}

func (a *archive) u16() uint16 {
	if !a.need(2) {
		return 0
	}
	a.i += 2
	return binary.LittleEndian.Uint16(a.b[a.i-2:])
}

func (a *archive) u32() uint32 {
	if !a.need(4) {
		return 0
	}
	a.i += 4
	return binary.LittleEndian.Uint32(a.b[a.i-4:])
}

func (a *archive) i32() int32 { return int32(a.u32()) }

func (a *archive) f64() float64 {
	if !a.need(8) {
		return 0
	}
	a.i += 8
	v := math.Float64frombits(binary.LittleEndian.Uint64(a.b[a.i-8:]))
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return v
}

func (a *archive) skip(n int) {
	if a.need(n) {
		a.i += n
	}
}

var sjis = japanese.ShiftJIS.NewDecoder()

// str reads a CString: a byte length (0xff: a 16-bit one follows, 0xffff:
// a 32-bit one), in Shift_JIS; 0xfffe marks a UTF-16 string.
func (a *archive) str() string {
	n := int(a.u8())
	unicode := false
	if n == 0xff {
		n = int(a.u16())
		if n == 0xfffe {
			unicode = true
			n = int(a.u8())
			if n == 0xff {
				n = int(a.u16())
				if n == 0xffff {
					n = int(a.u32())
				}
			}
		} else if n == 0xffff {
			n = int(a.u32())
		}
	}
	if unicode {
		n *= 2
	}
	if n > len(a.b) || !a.need(n) {
		return ""
	}
	raw := a.b[a.i : a.i+n]
	a.i += n
	if unicode {
		u := make([]uint16, n/2)
		for k := range u {
			u[k] = binary.LittleEndian.Uint16(raw[2*k:])
		}
		return string(utf16Decode(u))
	}
	s, err := sjis.String(string(raw))
	if err != nil {
		return string(raw)
	}
	return s
}

func utf16Decode(u []uint16) []rune {
	var out []rune
	for i := 0; i < len(u); i++ {
		r := rune(u[i])
		if r >= 0xd800 && r < 0xdc00 && i+1 < len(u) && u[i+1] >= 0xdc00 && u[i+1] < 0xe000 {
			r = (r-0xd800)<<10 + rune(u[i+1]-0xdc00) + 0x10000
			i++
		}
		out = append(out, r)
	}
	return out
}

// count reads the element count of a collection.
func (a *archive) count() int {
	n := uint32(a.u16())
	if n == 0xffff {
		n = a.u32()
	}
	if int(n) > len(a.b) {
		// more elements than bytes: the data is broken
		a.err = fmt.Errorf("a list of %d elements does not fit in the file", n)
		return 0
	}
	return int(n)
}

// objectTag reads the tag before an object. It returns the class of a new
// object (registering the class and the object's index), or an object
// written before, or neither for a null pointer.
func (a *archive) objectTag() (class string, index uint32, old any) {
	tag := uint32(a.u16())
	if a.err != nil || tag == 0 {
		return "", 0, nil
	}
	switch {
	case tag == 0xffff:
		// a new class: schema, name length and name
		a.u16()
		n := int(a.u16())
		if !a.need(n) {
			return "", 0, nil
		}
		class = string(a.b[a.i : a.i+n])
		a.i += n
		a.classes[a.next] = class
		a.next++
	case tag&0x8000 != 0:
		idx := tag & 0x7fff
		if idx == 0x7fff {
			idx = a.u32() & 0x7fffffff
		}
		c, ok := a.classes[idx]
		if !ok {
			a.err = fmt.Errorf("unknown class reference %d", idx)
			return "", 0, nil
		}
		class = c
	default:
		idx := tag
		if idx == 0x7fff {
			idx = a.u32()
		}
		o, ok := a.objects[idx]
		if !ok {
			a.err = fmt.Errorf("unknown object reference %d", idx)
		}
		return "", 0, o
	}
	index = a.next
	a.next++
	return class, index, nil
}

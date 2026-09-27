// Package cff reads and subsets CFF font programs: the Type 2 outlines of
// OpenType fonts with a "CFF " table, and those PDF files embed.
package cff

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

// Font holds what the converters need from a CFF font program (bare, as
// PDF embeds it, or the CFF table of an OpenType font): glyph count, glyph
// names (or CIDs) per GID, the built-in encoding and the font matrix.
type Font struct {
	Data       []byte
	NumGlyphs  int
	IsCID      bool
	Charset    []int // gid → SID (or CID)
	NameToGID  map[string]int
	CIDToGID   map[int]int
	Encoding   map[int]int // code → gid (built-in encoding)
	FontMatrix [6]float64
	Strings    []string
	// Notice and Copyright strings of the Top DICT (the only license text a
	// bare CFF program carries).
	Notice, Copyright string
}

// ReadIndex reads the INDEX at pos and returns its items and the position
// after it.
func ReadIndex(b []byte, pos int) ([][]byte, int, error) {
	if pos+2 > len(b) {
		return nil, pos, errors.New("cff: truncated index")
	}
	count := int(be16(b, pos))
	pos += 2
	if count == 0 {
		return nil, pos, nil
	}
	if pos >= len(b) {
		return nil, pos, errors.New("cff: truncated index")
	}
	offSize := int(b[pos])
	pos++
	if offSize < 1 || offSize > 4 {
		return nil, pos, errors.New("cff: bad offSize")
	}
	readOff := func(i int) int {
		p := pos + i*offSize
		v := 0
		for k := 0; k < offSize; k++ {
			if p+k < len(b) {
				v = v<<8 | int(b[p+k])
			}
		}
		return v
	}
	dataStart := pos + (count+1)*offSize - 1
	items := make([][]byte, count)
	for i := 0; i < count; i++ {
		s, e := dataStart+readOff(i), dataStart+readOff(i+1)
		if s < 0 || e > len(b) || s > e {
			return nil, pos, errors.New("cff: bad index offsets")
		}
		items[i] = b[s:e]
	}
	return items, dataStart + readOff(count), nil
}

// cffParseDict parses a CFF DICT into operator → operands.
// ParseDict reads a DICT into operator → operands (12 x is 1200+x).
func ParseDict(b []byte) map[int][]float64 {
	out := map[int][]float64{}
	var operands []float64
	for i := 0; i < len(b); {
		c := int(b[i])
		switch {
		case c <= 21:
			op := c
			i++
			if c == 12 && i < len(b) {
				op = 1200 + int(b[i])
				i++
			}
			out[op] = append([]float64(nil), operands...)
			operands = operands[:0]
		case c == 28:
			operands = append(operands, float64(int16(be16(b, i+1))))
			i += 3
		case c == 29:
			operands = append(operands, float64(int32(be32(b, i+1))))
			i += 5
		case c == 30:
			// real number, nibble encoded
			var sb strings.Builder
			i++
		loop:
			for i < len(b) {
				v := b[i]
				i++
				for _, nib := range []byte{v >> 4, v & 15} {
					switch {
					case nib <= 9:
						sb.WriteByte('0' + nib)
					case nib == 0xa:
						sb.WriteByte('.')
					case nib == 0xb:
						sb.WriteByte('E')
					case nib == 0xc:
						sb.WriteString("E-")
					case nib == 0xe:
						sb.WriteByte('-')
					case nib == 0xf:
						break loop
					}
				}
			}
			f, _ := strconv.ParseFloat(sb.String(), 64)
			operands = append(operands, f)
		case c >= 32 && c <= 246:
			operands = append(operands, float64(c-139))
			i++
		case c >= 247 && c <= 250:
			if i+1 < len(b) {
				operands = append(operands, float64((c-247)*256+int(b[i+1])+108))
			}
			i += 2
		case c >= 251 && c <= 254:
			if i+1 < len(b) {
				operands = append(operands, float64(-(c-251)*256-int(b[i+1])-108))
			}
			i += 2
		default:
			i++
		}
	}
	return out
}

// Parse reads a CFF font program.
func Parse(data []byte) (*Font, error) {
	if len(data) < 4 {
		return nil, errors.New("cff: too short")
	}
	f := &Font{Data: data, FontMatrix: [6]float64{0.001, 0, 0, 0.001, 0, 0}}
	pos := int(data[2])
	var err error
	if _, pos, err = ReadIndex(data, pos); err != nil { // Name INDEX
		return nil, err
	}
	topDicts, pos, err := ReadIndex(data, pos)
	if err != nil || len(topDicts) == 0 {
		return nil, errors.New("cff: no top dict")
	}
	strIdx, _, err := ReadIndex(data, pos)
	if err != nil {
		return nil, err
	}
	f.Strings = make([]string, len(strIdx))
	for i, s := range strIdx {
		f.Strings[i] = string(s)
	}
	top := ParseDict(topDicts[0])
	if fm := top[1207]; len(fm) == 6 {
		f.FontMatrix = [6]float64{fm[0], fm[1], fm[2], fm[3], fm[4], fm[5]}
	}
	if _, ok := top[1230]; ok {
		f.IsCID = true
	}
	if v := top[1]; len(v) == 1 {
		f.Notice = f.SIDName(int(v[0]))
	}
	if v := top[1200]; len(v) == 1 {
		f.Copyright = f.SIDName(int(v[0]))
	}
	csOff := top[17]
	if len(csOff) == 0 {
		return nil, errors.New("cff: no CharStrings")
	}
	charStrings, _, err := ReadIndex(data, int(csOff[0]))
	if err != nil {
		return nil, err
	}
	f.NumGlyphs = len(charStrings)
	if f.NumGlyphs == 0 {
		return nil, errors.New("cff: no glyphs")
	}

	// charset: gid → SID/CID
	f.Charset = make([]int, f.NumGlyphs)
	charsetOff := 0
	if v := top[15]; len(v) > 0 {
		charsetOff = int(v[0])
	}
	switch charsetOff {
	case 0, 1, 2: // ISOAdobe / Expert / ExpertSubset: identity SIDs are a good enough approximation
		for i := range f.Charset {
			f.Charset[i] = i
		}
	default:
		p := charsetOff
		if p >= 0 && p < len(data) {
			format := data[p]
			p++
			f.Charset[0] = 0
			gid := 1
			switch format {
			case 0:
				for gid < f.NumGlyphs && p+1 < len(data) {
					f.Charset[gid] = int(be16(data, p))
					p += 2
					gid++
				}
			case 1, 2:
				for gid < f.NumGlyphs && p+2 < len(data) {
					first := int(be16(data, p))
					var nLeft int
					if format == 1 {
						nLeft = int(data[p+2])
						p += 3
					} else {
						nLeft = int(be16(data, p+2))
						p += 4
					}
					for k := 0; k <= nLeft && gid < f.NumGlyphs; k++ {
						f.Charset[gid] = first + k
						gid++
					}
				}
			}
		}
	}
	if f.IsCID {
		f.CIDToGID = make(map[int]int, f.NumGlyphs)
		for gid, cid := range f.Charset {
			if _, dup := f.CIDToGID[cid]; !dup {
				f.CIDToGID[cid] = gid
			}
		}
	} else {
		f.NameToGID = make(map[string]int, f.NumGlyphs)
		for gid, sid := range f.Charset {
			name := f.SIDName(sid)
			if _, dup := f.NameToGID[name]; !dup && name != "" {
				f.NameToGID[name] = gid
			}
		}
		// built-in encoding: code → gid
		f.Encoding = map[int]int{}
		encOff := 0
		if v := top[16]; len(v) > 0 {
			encOff = int(v[0])
		}
		switch encOff {
		case 0, 1: // standard (expert treated as standard)
			for code, name := range StandardEncoding {
				if name == "" {
					continue
				}
				if gid, ok := f.NameToGID[name]; ok {
					f.Encoding[code] = gid
				}
			}
		default:
			p := encOff
			if p >= 0 && p+1 < len(data) {
				format := data[p]
				p++
				switch format & 0x7f {
				case 0:
					n := int(data[p])
					p++
					for i := 1; i <= n && p < len(data); i++ {
						f.Encoding[int(data[p])] = i
						p++
					}
				case 1:
					nRanges := int(data[p])
					p++
					gid := 1
					for i := 0; i < nRanges && p+1 < len(data); i++ {
						first, nLeft := int(data[p]), int(data[p+1])
						p += 2
						for k := 0; k <= nLeft; k++ {
							f.Encoding[first+k] = gid
							gid++
						}
					}
				}
				if format&0x80 != 0 && p < len(data) {
					nSups := int(data[p])
					p++
					for i := 0; i < nSups && p+2 < len(data); i++ {
						code, sid := int(data[p]), int(be16(data, p+1))
						p += 3
						if gid, ok := f.NameToGID[f.SIDName(sid)]; ok {
							f.Encoding[code] = gid
						}
					}
				}
			}
		}
	}
	return f, nil
}

// SIDName returns the string of a string ID.
func (f *Font) SIDName(sid int) string {
	if sid < 0 {
		return ""
	}
	if sid < len(StandardStrings) {
		return StandardStrings[sid]
	}
	if i := sid - len(StandardStrings); i < len(f.Strings) {
		return f.Strings[i]
	}
	return ""
}

func be32(b []byte, i int) uint32 {
	if i < 0 || i+4 > len(b) {
		return 0
	}
	return binary.BigEndian.Uint32(b[i:])
}

func be16(b []byte, i int) uint16 {
	if i < 0 || i+2 > len(b) {
		return 0
	}
	return binary.BigEndian.Uint16(b[i:])
}

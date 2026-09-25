package pdf

import (
	"errors"
	"strconv"
	"strings"
)

// cffFont holds what we need from a bare CFF font program: glyph count,
// glyph names (or CIDs) per GID, the built-in encoding and the font matrix.
type cffFont struct {
	data       []byte
	numGlyphs  int
	isCID      bool
	charset    []int // gid → SID (or CID)
	nameToGID  map[string]int
	cidToGID   map[int]int
	encoding   map[int]int // code → gid (built-in encoding)
	fontMatrix matrix
	strings    []string
}

func cffReadIndex(b []byte, pos int) ([][]byte, int, error) {
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
func cffParseDict(b []byte) map[int][]float64 {
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

func parseCFF(data []byte) (*cffFont, error) {
	if len(data) < 4 {
		return nil, errors.New("cff: too short")
	}
	f := &cffFont{data: data, fontMatrix: matrix{0.001, 0, 0, 0.001, 0, 0}}
	pos := int(data[2])
	var err error
	if _, pos, err = cffReadIndex(data, pos); err != nil { // Name INDEX
		return nil, err
	}
	topDicts, pos, err := cffReadIndex(data, pos)
	if err != nil || len(topDicts) == 0 {
		return nil, errors.New("cff: no top dict")
	}
	strIdx, _, err := cffReadIndex(data, pos)
	if err != nil {
		return nil, err
	}
	f.strings = make([]string, len(strIdx))
	for i, s := range strIdx {
		f.strings[i] = string(s)
	}
	top := cffParseDict(topDicts[0])
	if fm := top[1207]; len(fm) == 6 {
		f.fontMatrix = matrix{fm[0], fm[1], fm[2], fm[3], fm[4], fm[5]}
	}
	if _, ok := top[1230]; ok {
		f.isCID = true
	}
	csOff := top[17]
	if len(csOff) == 0 {
		return nil, errors.New("cff: no CharStrings")
	}
	charStrings, _, err := cffReadIndex(data, int(csOff[0]))
	if err != nil {
		return nil, err
	}
	f.numGlyphs = len(charStrings)

	// charset: gid → SID/CID
	f.charset = make([]int, f.numGlyphs)
	charsetOff := 0
	if v := top[15]; len(v) > 0 {
		charsetOff = int(v[0])
	}
	switch charsetOff {
	case 0, 1, 2: // ISOAdobe / Expert / ExpertSubset: identity SIDs are a good enough approximation
		for i := range f.charset {
			f.charset[i] = i
		}
	default:
		p := charsetOff
		if p < len(data) {
			format := data[p]
			p++
			f.charset[0] = 0
			gid := 1
			switch format {
			case 0:
				for gid < f.numGlyphs && p+1 < len(data) {
					f.charset[gid] = int(be16(data, p))
					p += 2
					gid++
				}
			case 1, 2:
				for gid < f.numGlyphs && p+2 < len(data) {
					first := int(be16(data, p))
					var nLeft int
					if format == 1 {
						nLeft = int(data[p+2])
						p += 3
					} else {
						nLeft = int(be16(data, p+2))
						p += 4
					}
					for k := 0; k <= nLeft && gid < f.numGlyphs; k++ {
						f.charset[gid] = first + k
						gid++
					}
				}
			}
		}
	}
	if f.isCID {
		f.cidToGID = make(map[int]int, f.numGlyphs)
		for gid, cid := range f.charset {
			if _, dup := f.cidToGID[cid]; !dup {
				f.cidToGID[cid] = gid
			}
		}
	} else {
		f.nameToGID = make(map[string]int, f.numGlyphs)
		for gid, sid := range f.charset {
			name := f.sidName(sid)
			if _, dup := f.nameToGID[name]; !dup && name != "" {
				f.nameToGID[name] = gid
			}
		}
		// built-in encoding: code → gid
		f.encoding = map[int]int{}
		encOff := 0
		if v := top[16]; len(v) > 0 {
			encOff = int(v[0])
		}
		switch encOff {
		case 0, 1: // standard (expert treated as standard)
			for code, name := range standardEncoding {
				if name == "" {
					continue
				}
				if gid, ok := f.nameToGID[name]; ok {
					f.encoding[code] = gid
				}
			}
		default:
			p := encOff
			if p < len(data) {
				format := data[p]
				p++
				switch format & 0x7f {
				case 0:
					n := int(data[p])
					p++
					for i := 1; i <= n && p < len(data); i++ {
						f.encoding[int(data[p])] = i
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
							f.encoding[first+k] = gid
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
						if gid, ok := f.nameToGID[f.sidName(sid)]; ok {
							f.encoding[code] = gid
						}
					}
				}
			}
		}
	}
	return f, nil
}

func (f *cffFont) sidName(sid int) string {
	if sid < len(cffStandardStrings) {
		return cffStandardStrings[sid]
	}
	if i := sid - len(cffStandardStrings); i < len(f.strings) {
		return f.strings[i]
	}
	return ""
}

// glyphNameToRune maps a glyph name to a Unicode code point using the AGL and
// the uniXXXX / uXXXX[XX] conventions.
func glyphNameToRune(name string) (rune, bool) {
	if name == "" {
		return 0, false
	}
	if i := strings.IndexByte(name, '.'); i > 0 {
		name = name[:i]
	}
	if r, ok := glyphList[name]; ok {
		return r, true
	}
	if strings.HasPrefix(name, "uni") && len(name) >= 7 {
		if v, err := strconv.ParseUint(name[3:7], 16, 32); err == nil {
			return rune(v), true
		}
	}
	if strings.HasPrefix(name, "u") && len(name) >= 5 && len(name) <= 7 {
		if v, err := strconv.ParseUint(name[1:], 16, 32); err == nil {
			return rune(v), true
		}
	}
	// Names like "Cxx"/"Gxx" carry no Unicode.
	return 0, false
}

// glyphNameIndex parses names of the form gNN, glyphNN, cidNN, GNN, index NN.
func glyphNameIndex(name string) (int, bool) {
	for _, prefix := range []string{"glyph", "cid", "index", "g", "G", "c", "C"} {
		if strings.HasPrefix(name, prefix) {
			if v, err := strconv.Atoi(name[len(prefix):]); err == nil {
				return v, true
			}
		}
	}
	return 0, false
}

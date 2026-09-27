package woff2

import (
	"encoding/binary"
	"errors"
	"sort"
)

// ErrDecodeNotAvailable is returned by Decode when the Brotli decoder is not
// compiled in (builds tagged bdf_noconv).
var ErrDecodeNotAvailable = errors.New("woff2: decoder not available in this build")

var errFormat = errors.New("woff2: malformed font")

// Decode unpacks a WOFF2 font into an sfnt font (TrueType or CFF-flavored
// OpenType), rebuilding the transformed glyf, loca and hmtx tables. Font
// collections are not supported.
func Decode(data []byte) ([]byte, error) {
	if len(data) < 48 || string(data[:4]) != "wOF2" {
		return nil, errors.New("woff2: not a WOFF2 font")
	}
	flavor := binary.BigEndian.Uint32(data[4:])
	if flavor == 0x74746366 { // 'ttcf'
		return nil, errors.New("woff2: font collections are not supported")
	}
	numTables := int(binary.BigEndian.Uint16(data[12:]))
	compressedLen := int(binary.BigEndian.Uint32(data[20:]))
	type entry struct {
		tag                string
		version            int
		origLen, transLen  int
		transformed        bool
		offset, streamSize int
	}
	pos := 48
	entries := make([]entry, numTables)
	for i := range entries {
		if pos >= len(data) {
			return nil, errFormat
		}
		flags := data[pos]
		pos++
		e := entry{version: int(flags >> 6)}
		if idx := int(flags & 0x3f); idx == 0x3f {
			if pos+4 > len(data) {
				return nil, errFormat
			}
			e.tag = string(data[pos : pos+4])
			pos += 4
		} else if idx < len(knownTags) {
			e.tag = knownTags[idx]
		} else {
			return nil, errFormat
		}
		var err error
		if e.origLen, pos, err = readBase128(data, pos); err != nil {
			return nil, err
		}
		switch e.tag {
		case "glyf", "loca":
			e.transformed = e.version == 0
		default:
			e.transformed = e.version != 0
		}
		e.streamSize = e.origLen
		if e.transformed {
			if e.transLen, pos, err = readBase128(data, pos); err != nil {
				return nil, err
			}
			e.streamSize = e.transLen
		}
		entries[i] = e
	}
	if compressedLen < 0 || pos+compressedLen > len(data) {
		return nil, errFormat
	}
	stream, err := decompress(data[pos : pos+compressedLen])
	if err != nil {
		return nil, err
	}
	tables := map[string][]byte{}
	off := 0
	for i := range entries {
		e := &entries[i]
		if e.streamSize < 0 || off+e.streamSize > len(stream) {
			return nil, errFormat
		}
		e.offset = off
		off += e.streamSize
		if !e.transformed {
			tables[e.tag] = stream[e.offset : e.offset+e.streamSize]
		}
	}
	for _, e := range entries {
		if !e.transformed {
			continue
		}
		switch e.tag {
		case "glyf":
			glyf, loca, err := rebuildGlyf(stream[e.offset : e.offset+e.streamSize])
			if err != nil {
				return nil, err
			}
			tables["glyf"], tables["loca"] = glyf, loca
			if head := tables["head"]; len(head) >= 54 {
				h := append([]byte(nil), head...)
				binary.BigEndian.PutUint16(h[50:], 1) // rebuilt in the long format
				tables["head"] = h
			}
		case "loca":
			// rebuilt with glyf
		case "hmtx":
			// handled below, once glyf is known
		default:
			return nil, errors.New("woff2: unknown transform of table " + e.tag)
		}
	}
	for _, e := range entries {
		if e.transformed && e.tag == "hmtx" {
			h, err := rebuildHmtx(stream[e.offset:e.offset+e.streamSize], tables)
			if err != nil {
				return nil, err
			}
			tables["hmtx"] = h
		}
	}
	return buildSFNT(flavor, tables), nil
}

// readBase128 reads a UIntBase128.
func readBase128(b []byte, pos int) (int, int, error) {
	var v uint32
	for i := 0; i < 5; i++ {
		if pos >= len(b) {
			return 0, pos, errFormat
		}
		c := b[pos]
		pos++
		if i == 0 && c == 0x80 || v&0xfe000000 != 0 {
			return 0, pos, errFormat
		}
		v = v<<7 | uint32(c&0x7f)
		if c&0x80 == 0 {
			if v > 1<<30 {
				return 0, pos, errFormat
			}
			return int(v), pos, nil
		}
	}
	return 0, pos, errFormat
}

// stream is a cursor over one of the streams of a transformed glyf table.
type stream struct {
	b   []byte
	pos int
	err bool
}

func (s *stream) u8() int {
	if s.pos >= len(s.b) {
		s.err = true
		return 0
	}
	v := s.b[s.pos]
	s.pos++
	return int(v)
}

func (s *stream) u16() int {
	if s.pos+2 > len(s.b) {
		s.err = true
		return 0
	}
	v := binary.BigEndian.Uint16(s.b[s.pos:])
	s.pos += 2
	return int(v)
}

func (s *stream) bytes(n int) []byte {
	if n < 0 || s.pos+n > len(s.b) {
		s.err = true
		return nil
	}
	v := s.b[s.pos : s.pos+n]
	s.pos += n
	return v
}

// u255 reads a 255UInt16.
func (s *stream) u255() int {
	switch c := s.u8(); c {
	case 253:
		return s.u16()
	case 255:
		return s.u8() + 253
	case 254:
		return s.u8() + 506
	default:
		return c
	}
}

// rebuildGlyf reverses the glyf transform (WOFF2 §5.1): it returns the glyf
// table and a long-format loca table.
func rebuildGlyf(t []byte) (glyf, loca []byte, err error) {
	if len(t) < 36 {
		return nil, nil, errFormat
	}
	optionFlags := binary.BigEndian.Uint16(t[2:])
	numGlyphs := int(binary.BigEndian.Uint16(t[4:]))
	pos := 36
	var ss [7]*stream
	for i := range ss {
		n := int(binary.BigEndian.Uint32(t[8+4*i:]))
		if n < 0 || pos+n > len(t) {
			return nil, nil, errFormat
		}
		ss[i] = &stream{b: t[pos : pos+n]}
		pos += n
	}
	nContour, nPoints, flagS, glyphS, compS, bboxS, instrS := ss[0], ss[1], ss[2], ss[3], ss[4], ss[5], ss[6]
	bitmapLen := ((numGlyphs + 31) >> 5) << 2
	bboxBitmap := bboxS.bytes(bitmapLen)
	var overlap []byte
	if optionFlags&1 != 0 {
		if pos+(numGlyphs+7)>>3 > len(t) {
			return nil, nil, errFormat
		}
		overlap = t[pos : pos+(numGlyphs+7)>>3]
	}
	if bboxS.err {
		return nil, nil, errFormat
	}
	loca = make([]byte, 0, (numGlyphs+1)*4)
	for g := 0; g < numGlyphs; g++ {
		loca = binary.BigEndian.AppendUint32(loca, uint32(len(glyf)))
		hasBBox := bboxBitmap[g>>3]&(0x80>>(g&7)) != 0
		n := int(int16(nContour.u16()))
		switch {
		case n == 0:
			if hasBBox {
				return nil, nil, errFormat
			}
		case n < 0:
			if !hasBBox {
				return nil, nil, errFormat
			}
			glyf = binary.BigEndian.AppendUint16(glyf, 0xffff)
			glyf = append(glyf, bboxS.bytes(8)...)
			// components: copy their records, noting whether any has instructions
			haveInstr := false
			for {
				flags := compS.u16()
				size := 4 // glyph index and two byte arguments
				if flags&compArgsAreWords != 0 {
					size += 2
				}
				switch {
				case flags&compHaveScale != 0:
					size += 2
				case flags&compHaveXYScale != 0:
					size += 4
				case flags&compHave2x2 != 0:
					size += 8
				}
				glyf = binary.BigEndian.AppendUint16(glyf, uint16(flags))
				glyf = append(glyf, compS.bytes(size)...)
				if flags&compHaveInstr != 0 {
					haveInstr = true
				}
				if compS.err || flags&compMore == 0 {
					break
				}
			}
			if haveInstr {
				ilen := glyphS.u255()
				glyf = binary.BigEndian.AppendUint16(glyf, uint16(ilen))
				glyf = append(glyf, instrS.bytes(ilen)...)
			}
		default:
			ends := make([]int, n)
			total := 0
			for i := range ends {
				total += nPoints.u255()
				ends[i] = total - 1
			}
			if total > 0xffff {
				return nil, nil, errFormat
			}
			xs, ys := make([]int, total), make([]int, total)
			on := make([]bool, total)
			x, y := 0, 0
			for i := 0; i < total; i++ {
				f := flagS.u8()
				on[i] = f&0x80 == 0
				dx, dy := decodeTriplet(f&0x7f, glyphS)
				x += dx
				y += dy
				xs[i], ys[i] = x, y
			}
			ilen := glyphS.u255()
			instr := instrS.bytes(ilen)
			var xMin, yMin, xMax, yMax int
			if hasBBox {
				b := bboxS.bytes(8)
				if len(b) == 8 {
					xMin, yMin = int(int16(binary.BigEndian.Uint16(b))), int(int16(binary.BigEndian.Uint16(b[2:])))
					xMax, yMax = int(int16(binary.BigEndian.Uint16(b[4:]))), int(int16(binary.BigEndian.Uint16(b[6:])))
				}
			} else if total > 0 {
				xMin, yMin, xMax, yMax = xs[0], ys[0], xs[0], ys[0]
				for i := 1; i < total; i++ {
					xMin, xMax = min(xMin, xs[i]), max(xMax, xs[i])
					yMin, yMax = min(yMin, ys[i]), max(yMax, ys[i])
				}
			}
			glyf = binary.BigEndian.AppendUint16(glyf, uint16(n))
			for _, v := range []int{xMin, yMin, xMax, yMax} {
				glyf = binary.BigEndian.AppendUint16(glyf, uint16(int16(v)))
			}
			for _, e := range ends {
				glyf = binary.BigEndian.AppendUint16(glyf, uint16(e))
			}
			glyf = binary.BigEndian.AppendUint16(glyf, uint16(ilen))
			glyf = append(glyf, instr...)
			// flags, then coordinates as 16-bit deltas
			overlapBit := overlap != nil && overlap[g>>3]&(0x80>>(g&7)) != 0
			for i := 0; i < total; i++ {
				var f byte
				if on[i] {
					f = flagOnCurve
				}
				if i == 0 && overlapBit {
					f |= flagOverlap
				}
				glyf = append(glyf, f)
			}
			px := 0
			for _, v := range xs {
				glyf = binary.BigEndian.AppendUint16(glyf, uint16(int16(v-px)))
				px = v
			}
			py := 0
			for _, v := range ys {
				glyf = binary.BigEndian.AppendUint16(glyf, uint16(int16(v-py)))
				py = v
			}
		}
		for len(glyf)%4 != 0 {
			glyf = append(glyf, 0)
		}
	}
	for _, s := range ss {
		if s.err {
			return nil, nil, errFormat
		}
	}
	loca = binary.BigEndian.AppendUint32(loca, uint32(len(glyf)))
	return glyf, loca, nil
}

// decodeTriplet reads the data bytes of a point delta for a flag byte
// without its on-curve bit (WOFF2 §5.2).
func decodeTriplet(flag int, s *stream) (dx, dy int) {
	sign := func(f, v int) int {
		if f&1 != 0 {
			return v
		}
		return -v
	}
	switch {
	case flag < 10:
		dy = sign(flag, (flag&14)<<7+s.u8())
	case flag < 20:
		dx = sign(flag, ((flag-10)&14)<<7+s.u8())
	case flag < 84:
		b0 := flag - 20
		b1 := s.u8()
		dx = sign(flag, 1+(b0&0x30)+b1>>4)
		dy = sign(flag>>1, 1+(b0&0x0c)<<2+b1&0x0f)
	case flag < 120:
		b0 := flag - 84
		dx = sign(flag, 1+(b0/12)<<8+s.u8())
		dy = sign(flag>>1, 1+((b0%12)>>2)<<8+s.u8())
	case flag < 124:
		b1, b2, b3 := s.u8(), s.u8(), s.u8()
		dx = sign(flag, b1<<4+b2>>4)
		dy = sign(flag>>1, (b2&0x0f)<<8+b3)
	default:
		b1, b2, b3, b4 := s.u8(), s.u8(), s.u8(), s.u8()
		dx = sign(flag, b1<<8+b2)
		dy = sign(flag>>1, b3<<8+b4)
	}
	return dx, dy
}

// rebuildHmtx reverses the hmtx transform (WOFF2 §5.4): left side bearings
// left out equal the glyphs' xMin.
func rebuildHmtx(t []byte, tables map[string][]byte) ([]byte, error) {
	hhea, maxp, glyf, loca := tables["hhea"], tables["maxp"], tables["glyf"], tables["loca"]
	if len(t) < 1 || len(hhea) < 36 || len(maxp) < 6 {
		return nil, errFormat
	}
	numH := int(binary.BigEndian.Uint16(hhea[34:]))
	numGlyphs := int(binary.BigEndian.Uint16(maxp[4:]))
	if numH < 1 || numH > numGlyphs {
		return nil, errFormat
	}
	s := &stream{b: t, pos: 1}
	flags := t[0]
	adv := make([]int, numH)
	for i := range adv {
		adv[i] = s.u16()
	}
	xMin := func(g int) int {
		if len(loca) < (g+2)*4 {
			return 0
		}
		o := int(binary.BigEndian.Uint32(loca[g*4:]))
		if o+10 > len(glyf) || o == int(binary.BigEndian.Uint32(loca[g*4+4:])) {
			return 0
		}
		return int(int16(binary.BigEndian.Uint16(glyf[o+2:])))
	}
	lsb := make([]int, numGlyphs)
	for g := 0; g < numGlyphs; g++ {
		switch {
		case g < numH && flags&1 == 0, g >= numH && flags&2 == 0:
			lsb[g] = int(int16(s.u16()))
		default:
			lsb[g] = xMin(g)
		}
	}
	if s.err {
		return nil, errFormat
	}
	out := make([]byte, 0, numH*4+(numGlyphs-numH)*2)
	for g := 0; g < numGlyphs; g++ {
		if g < numH {
			out = binary.BigEndian.AppendUint16(out, uint16(adv[g]))
		}
		out = binary.BigEndian.AppendUint16(out, uint16(int16(lsb[g])))
	}
	return out, nil
}

// buildSFNT writes an sfnt font with its tables in tag order.
func buildSFNT(flavor uint32, tables map[string][]byte) []byte {
	tags := make([]string, 0, len(tables))
	for t := range tables {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	n := len(tags)
	es, sr := 0, 1
	for sr*2 <= n {
		sr *= 2
		es++
	}
	out := make([]byte, 12+16*n)
	binary.BigEndian.PutUint32(out, flavor)
	binary.BigEndian.PutUint16(out[4:], uint16(n))
	binary.BigEndian.PutUint16(out[6:], uint16(sr*16))
	binary.BigEndian.PutUint16(out[8:], uint16(es))
	binary.BigEndian.PutUint16(out[10:], uint16(n*16-sr*16))
	for i, tag := range tags {
		data := tables[tag]
		rec := out[12+16*i:]
		copy(rec, tag)
		binary.BigEndian.PutUint32(rec[4:], checksum(data))
		binary.BigEndian.PutUint32(rec[8:], uint32(len(out)))
		binary.BigEndian.PutUint32(rec[12:], uint32(len(data)))
		out = append(out, data...)
		for len(out)%4 != 0 {
			out = append(out, 0)
		}
	}
	return out
}

func checksum(b []byte) uint32 {
	var sum uint32
	for i := 0; i < len(b); i += 4 {
		var w [4]byte
		copy(w[:], b[i:])
		sum += binary.BigEndian.Uint32(w[:])
	}
	return sum
}

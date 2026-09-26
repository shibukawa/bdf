package woff2

import "encoding/binary"

// glyf flags (simple glyphs).
const (
	flagOnCurve      = 0x01
	flagXShort       = 0x02
	flagYShort       = 0x04
	flagRepeat       = 0x08
	flagXSame        = 0x10
	flagYSame        = 0x20
	flagOverlap      = 0x40
	flagReserved     = 0x80 // cubic outlines in newer glyf; not representable in WOFF2
	compArgsAreWords = 0x0001
	compHaveScale    = 0x0008
	compMore         = 0x0020
	compHaveXYScale  = 0x0040
	compHave2x2      = 0x0080
	compHaveInstr    = 0x0100
)

// transformGlyf produces the WOFF2 transformed glyf table (spec §5.1). It
// reports false when a glyph cannot be represented, and the caller then stores
// glyf and loca untransformed.
func transformGlyf(glyf, loca []byte, numGlyphs int, longLoca bool) ([]byte, bool) {
	offset := func(i int) (int, bool) {
		if longLoca {
			if i*4+4 > len(loca) {
				return 0, false
			}
			return int(binary.BigEndian.Uint32(loca[i*4:])), true
		}
		if i*2+2 > len(loca) {
			return 0, false
		}
		return int(binary.BigEndian.Uint16(loca[i*2:])) * 2, true
	}
	var nContours, nPoints, flagS, glyphS, compS, bboxS, instrS []byte
	bboxBitmap := make([]byte, ((numGlyphs+31)>>5)<<2)
	overlap := make([]byte, (numGlyphs+7)>>3)
	hasOverlap := false
	start, ok := offset(0)
	if !ok {
		return nil, false
	}
	for g := 0; g < numGlyphs; g++ {
		end, ok := offset(g + 1)
		if !ok || end < start || end > len(glyf) {
			return nil, false
		}
		d := glyf[start:end]
		start = end
		if len(d) == 0 {
			nContours = binary.BigEndian.AppendUint16(nContours, 0)
			continue
		}
		if len(d) < 10 {
			return nil, false
		}
		nc := int16(binary.BigEndian.Uint16(d))
		switch {
		case nc == 0:
			// No outline: the decoder emits an empty glyph and forbids a bbox.
			nContours = binary.BigEndian.AppendUint16(nContours, 0)
		case nc > 0:
			var bbox [4]int16
			var overlapped bool
			p := 10
			ends := make([]int, nc)
			prev := -1
			for i := range ends {
				if p+2 > len(d) {
					return nil, false
				}
				ends[i] = int(binary.BigEndian.Uint16(d[p:]))
				if ends[i] <= prev {
					return nil, false
				}
				nPoints = append255(nPoints, ends[i]-prev)
				prev = ends[i]
				p += 2
			}
			np := prev + 1
			if p+2 > len(d) {
				return nil, false
			}
			instrLen := int(binary.BigEndian.Uint16(d[p:]))
			p += 2
			if p+instrLen > len(d) {
				return nil, false
			}
			instr := d[p : p+instrLen]
			p += instrLen
			flags := make([]byte, 0, np)
			for len(flags) < np {
				if p >= len(d) {
					return nil, false
				}
				f := d[p]
				p++
				if f&flagReserved != 0 {
					return nil, false
				}
				flags = append(flags, f)
				if f&flagRepeat != 0 {
					if p >= len(d) {
						return nil, false
					}
					for r := int(d[p]); r > 0 && len(flags) < np; r-- {
						flags = append(flags, f)
					}
					p++
				}
			}
			overlapped = flags[0]&flagOverlap != 0
			coord := func(short, same byte) ([]int, bool) {
				out := make([]int, np)
				for i, f := range flags {
					switch {
					case f&short != 0:
						if p >= len(d) {
							return nil, false
						}
						v := int(d[p])
						p++
						if f&same == 0 {
							v = -v
						}
						out[i] = v
					case f&same == 0:
						if p+2 > len(d) {
							return nil, false
						}
						out[i] = int(int16(binary.BigEndian.Uint16(d[p:])))
						p += 2
					}
				}
				return out, true
			}
			dx, ok := coord(flagXShort, flagXSame)
			if !ok {
				return nil, false
			}
			dy, ok := coord(flagYShort, flagYSame)
			if !ok {
				return nil, false
			}
			x, y := 0, 0
			for i := range flags {
				x += dx[i]
				y += dy[i]
				if i == 0 || x < int(bbox[0]) {
					bbox[0] = int16(x)
				}
				if i == 0 || y < int(bbox[1]) {
					bbox[1] = int16(y)
				}
				if i == 0 || x > int(bbox[2]) {
					bbox[2] = int16(x)
				}
				if i == 0 || y > int(bbox[3]) {
					bbox[3] = int16(y)
				}
				var on byte = 128
				if flags[i]&flagOnCurve != 0 {
					on = 0
				}
				flagS, glyphS = appendTriplet(flagS, glyphS, on, dx[i], dy[i])
			}
			glyphS = append255(glyphS, instrLen)
			instrS = append(instrS, instr...)
			nContours = binary.BigEndian.AppendUint16(nContours, uint16(nc))
			// The bbox is implied when it matches the points.
			for i := 0; i < 4; i++ {
				if int16(binary.BigEndian.Uint16(d[2+i*2:])) != bbox[i] {
					bboxBitmap[g>>3] |= 0x80 >> (g & 7)
					bboxS = append(bboxS, d[2:10]...)
					break
				}
			}
			if overlapped {
				overlap[g>>3] |= 0x80 >> (g & 7)
				hasOverlap = true
			}
		default: // composite
			p := 10
			haveInstr := false
			for {
				if p+4 > len(d) {
					return nil, false
				}
				f := binary.BigEndian.Uint16(d[p:])
				haveInstr = haveInstr || f&compHaveInstr != 0
				p += 4
				if f&compArgsAreWords != 0 {
					p += 4
				} else {
					p += 2
				}
				switch {
				case f&compHaveScale != 0:
					p += 2
				case f&compHaveXYScale != 0:
					p += 4
				case f&compHave2x2 != 0:
					p += 8
				}
				if f&compMore == 0 {
					break
				}
			}
			if p > len(d) {
				return nil, false
			}
			compS = append(compS, d[10:p]...)
			if haveInstr {
				if p+2 > len(d) {
					return nil, false
				}
				n := int(binary.BigEndian.Uint16(d[p:]))
				if p+2+n > len(d) {
					return nil, false
				}
				glyphS = append255(glyphS, n)
				instrS = append(instrS, d[p+2:p+2+n]...)
			}
			nContours = binary.BigEndian.AppendUint16(nContours, 0xffff)
			bboxBitmap[g>>3] |= 0x80 >> (g & 7)
			bboxS = append(bboxS, d[2:10]...)
		}
	}
	bbox := append(bboxBitmap, bboxS...)
	out := make([]byte, 36)
	var optionFlags uint16
	if hasOverlap {
		optionFlags = 1
	}
	binary.BigEndian.PutUint16(out[2:], optionFlags)
	binary.BigEndian.PutUint16(out[4:], uint16(numGlyphs))
	binary.BigEndian.PutUint16(out[6:], 1) // loca is rebuilt in the long format
	for i, s := range [][]byte{nContours, nPoints, flagS, glyphS, compS, bbox, instrS} {
		binary.BigEndian.PutUint32(out[8+i*4:], uint32(len(s)))
	}
	for _, s := range [][]byte{nContours, nPoints, flagS, glyphS, compS, bbox, instrS} {
		out = append(out, s...)
	}
	if hasOverlap {
		out = append(out, overlap...)
	}
	return out, true
}

// appendTriplet encodes one point delta as a flag byte and 1–4 data bytes (spec §5.2).
func appendTriplet(flags, data []byte, onCurveBit byte, x, y int) ([]byte, []byte) {
	absX, absY := x, y
	if absX < 0 {
		absX = -absX
	}
	if absY < 0 {
		absY = -absY
	}
	var xSign, ySign byte
	if x >= 0 {
		xSign = 1
	}
	if y >= 0 {
		ySign = 1
	}
	xySign := xSign + 2*ySign
	switch {
	case x == 0 && absY < 1280:
		flags = append(flags, onCurveBit+byte((absY&0xf00)>>7)+ySign)
		data = append(data, byte(absY))
	case y == 0 && absX < 1280:
		flags = append(flags, onCurveBit+10+byte((absX&0xf00)>>7)+xSign)
		data = append(data, byte(absX))
	case absX < 65 && absY < 65:
		flags = append(flags, onCurveBit+20+byte((absX-1)&0x30)+byte(((absY-1)&0x30)>>2)+xySign)
		data = append(data, byte(((absX-1)&0xf)<<4|(absY-1)&0xf))
	case absX < 769 && absY < 769:
		flags = append(flags, onCurveBit+84+12*byte(((absX-1)&0x300)>>8)+byte(((absY-1)&0x300)>>6)+xySign)
		data = append(data, byte(absX-1), byte(absY-1))
	case absX < 4096 && absY < 4096:
		flags = append(flags, onCurveBit+120+xySign)
		data = append(data, byte(absX>>4), byte((absX&0xf)<<4|absY>>8), byte(absY))
	default:
		flags = append(flags, onCurveBit+124+xySign)
		data = append(data, byte(absX>>8), byte(absX), byte(absY>>8), byte(absY))
	}
	return flags, data
}

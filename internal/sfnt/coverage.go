package sfnt

// CmapCoverage is a compact view of the Unicode cmap subtable used to answer
// whether a font contains a character, without loading its glyph data.
type CmapCoverage struct {
	format uint16
	data   []byte
}

// ParseCmapCoverage reads the best Unicode cmap subtable (the same platform
// and encoding preference as Font.Cmap) and keeps only that subtable. It
// returns nil when the table has no supported, well-formed Unicode subtable.
func ParseCmapCoverage(table []byte) *CmapCoverage {
	if len(table) < 4 {
		return nil
	}
	n := int(be16(table, 2))
	if n > (len(table)-4)/8 {
		n = (len(table) - 4) / 8
	}
	best := -1
	var coverage *CmapCoverage
	for i := 0; i < n; i++ {
		rec := 4 + i*8
		pid, eid, off := be16(table, rec), be16(table, rec+2), int(be32(table, rec+4))
		if off < 0 || off+2 > len(table) {
			continue
		}
		format := be16(table, off)
		length, ok := cmapSubtableLength(table[off:])
		if !ok || off+length > len(table) || !validCoverageSubtable(table[off:off+length], format) {
			continue
		}
		score := -1
		switch {
		case pid == 3 && eid == 10:
			score = 3
		case pid == 3 && eid == 1:
			score = 2
		case pid == 0:
			score = 1
		}
		if score > best {
			best = score
			coverage = &CmapCoverage{format: format, data: append([]byte(nil), table[off:off+length]...)}
		}
	}
	return coverage
}

func cmapSubtableLength(b []byte) (int, bool) {
	if len(b) < 4 {
		return 0, false
	}
	switch be16(b, 0) {
	case 0, 4, 6:
		return int(be16(b, 2)), true
	case 12:
		if len(b) < 8 {
			return 0, false
		}
		length := be32(b, 4)
		if uint64(length) > uint64(^uint(0)>>1) {
			return 0, false
		}
		return int(length), true
	default:
		return 0, false
	}
}

func validCoverageSubtable(b []byte, format uint16) bool {
	switch format {
	case 0:
		return len(b) >= 262
	case 4:
		if len(b) < 16 {
			return false
		}
		segX2 := int(be16(b, 6))
		if segX2 == 0 || segX2&1 != 0 {
			return false
		}
		count := segX2 / 2
		ends, starts, rangeOffs := 14, 16+segX2, 16+3*segX2
		if rangeOffs+segX2 > len(b) {
			return false
		}
		var prev uint32
		for i := 0; i < count; i++ {
			end, start := uint32(be16(b, ends+i*2)), uint32(be16(b, starts+i*2))
			if start > end || (i > 0 && start <= prev) {
				return false
			}
			prev = end
		}
		return true
	case 6:
		if len(b) < 10 {
			return false
		}
		count := int(be16(b, 8))
		return 10+count*2 <= len(b)
	case 12:
		if len(b) < 16 {
			return false
		}
		count := int(be32(b, 12))
		if count > (len(b)-16)/12 || count > 0x110000 {
			return false
		}
		var prev uint32
		for i := 0; i < count; i++ {
			rec := 16 + i*12
			start, end := be32(b, rec), be32(b, rec+4)
			if start > end || (i > 0 && start <= prev) {
				return false
			}
			prev = end
		}
		return true
	default:
		return false
	}
}

// Has reports whether the selected Unicode subtable maps r to a nonzero
// glyph. It does not inspect glyph outlines or metrics.
func (c *CmapCoverage) Has(r rune) bool {
	if c == nil || r < 0 {
		return false
	}
	b := c.data
	code := uint32(r)
	switch c.format {
	case 0:
		return code < 256 && 6+int(code) < len(b) && b[6+int(code)] != 0
	case 4:
		if code > 0xfffe {
			return false
		}
		segX2 := int(be16(b, 6))
		count := segX2 / 2
		ends, starts, deltas, rangeOffs := 14, 16+segX2, 16+2*segX2, 16+3*segX2
		lo, hi := 0, count
		for lo < hi {
			mid := lo + (hi-lo)/2
			if uint32(be16(b, ends+mid*2)) < code {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo == count {
			return false
		}
		start := uint32(be16(b, starts+lo*2))
		if code < start {
			return false
		}
		delta, ro := be16(b, deltas+lo*2), int(be16(b, rangeOffs+lo*2))
		var glyph uint16
		if ro == 0 {
			glyph = uint16(code) + delta
		} else {
			addr := rangeOffs + lo*2 + ro + int(code-start)*2
			if addr < 0 || addr+2 > len(b) {
				return false
			}
			glyph = be16(b, addr)
			if glyph != 0 {
				glyph += delta
			}
		}
		return glyph != 0
	case 6:
		first, count := uint32(be16(b, 6)), int(be16(b, 8))
		if code < first || uint64(code-first) >= uint64(count) {
			return false
		}
		return be16(b, 10+int(code-first)*2) != 0
	case 12:
		if code > 0x10ffff {
			return false
		}
		count := int(be32(b, 12))
		lo, hi := 0, count
		for lo < hi {
			mid := lo + (hi-lo)/2
			rec := 16 + mid*12
			if be32(b, rec) <= code {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo == 0 {
			return false
		}
		rec := 16 + (lo-1)*12
		start, end, glyph := be32(b, rec), be32(b, rec+4), be32(b, rec+8)
		if code > end || code-start > 0xffff {
			return false
		}
		return uint16(glyph+(code-start)) != 0
	}
	return false
}

package pdf2bdf

import (
	"encoding/binary"
	"errors"
	"sort"
)

// Minimal TrueType/OpenType table reader and writer, enough to rebuild an
// embedded font program with a synthesized cmap so browsers can use it.

type sfnt struct {
	tables     map[string][]byte
	isCFF      bool
	numGlyphs  int
	unitsPerEm int
	ascent     int16
	descent    int16
	cmap       map[uint32]uint16 // best Unicode cmap (3,1)/(3,10)/(0,x)
	cmapMac    map[uint32]uint16 // (1,0)
	cmapSymbol map[uint32]uint16 // (3,0)
	postNames  map[string]uint16
	advances   []uint16
}

func be16(b []byte, i int) uint16 {
	if i+1 >= len(b) {
		return 0
	}
	return binary.BigEndian.Uint16(b[i:])
}

func be32(b []byte, i int) uint32 {
	if i+3 >= len(b) {
		return 0
	}
	return binary.BigEndian.Uint32(b[i:])
}

// parseSFNT reads the table directory and the tables needed for glyph lookup.
func parseSFNT(data []byte) (*sfnt, error) {
	if len(data) < 12 {
		return nil, errors.New("font too short")
	}
	tag := be32(data, 0)
	off := 0
	if tag == 0x74746366 { // 'ttcf': use the first font
		off = int(be32(data, 12))
		if off+12 > len(data) {
			return nil, errors.New("bad TTC")
		}
		tag = be32(data, off)
	}
	f := &sfnt{tables: map[string][]byte{}, unitsPerEm: 1000}
	switch tag {
	case 0x00010000, 0x74727565: // 1.0, 'true'
	case 0x4f54544f: // 'OTTO'
		f.isCFF = true
	default:
		return nil, errors.New("not an sfnt font")
	}
	n := int(be16(data, off+4))
	for i := 0; i < n; i++ {
		rec := off + 12 + i*16
		if rec+16 > len(data) {
			break
		}
		name := string(data[rec : rec+4])
		toff, tlen := int(be32(data, rec+8)), int(be32(data, rec+12))
		if toff < 0 || toff > len(data) {
			continue
		}
		if toff+tlen > len(data) {
			tlen = len(data) - toff
		}
		f.tables[name] = data[toff : toff+tlen]
	}
	if cff, ok := f.tables["CFF "]; ok && f.tables["glyf"] == nil {
		f.isCFF = true
		_ = cff
	}
	if maxp := f.tables["maxp"]; len(maxp) >= 6 {
		f.numGlyphs = int(be16(maxp, 4))
	}
	if head := f.tables["head"]; len(head) >= 54 {
		if u := int(be16(head, 18)); u > 0 {
			f.unitsPerEm = u
		}
	}
	if hhea := f.tables["hhea"]; len(hhea) >= 36 {
		f.ascent = int16(be16(hhea, 4))
		f.descent = int16(be16(hhea, 6))
		nhm := int(be16(hhea, 34))
		if hmtx := f.tables["hmtx"]; hmtx != nil {
			f.advances = make([]uint16, nhm)
			for i := 0; i < nhm; i++ {
				f.advances[i] = be16(hmtx, i*4)
			}
		}
	}
	f.parseCmap()
	f.parsePost()
	return f, nil
}

func (f *sfnt) parseCmap() {
	cm := f.tables["cmap"]
	if len(cm) < 4 {
		return
	}
	n := int(be16(cm, 2))
	best := -1
	for i := 0; i < n; i++ {
		rec := 4 + i*8
		pid, eid, off := be16(cm, rec), be16(cm, rec+2), int(be32(cm, rec+4))
		if off >= len(cm) {
			continue
		}
		m := parseCmapSubtable(cm[off:])
		if m == nil {
			continue
		}
		switch {
		case pid == 3 && eid == 0:
			f.cmapSymbol = m
		case pid == 1 && eid == 0:
			f.cmapMac = m
		case pid == 3 && (eid == 1 || eid == 10), pid == 0:
			score := 1
			if pid == 3 && eid == 10 {
				score = 3
			} else if pid == 3 {
				score = 2
			}
			if score > best {
				best = score
				f.cmap = m
			}
		}
	}
}

func parseCmapSubtable(b []byte) map[uint32]uint16 {
	format := be16(b, 0)
	m := map[uint32]uint16{}
	switch format {
	case 0:
		for c := 0; c < 256 && 6+c < len(b); c++ {
			if g := uint16(b[6+c]); g != 0 {
				m[uint32(c)] = g
			}
		}
	case 4:
		segX2 := int(be16(b, 6))
		ends, starts, deltas, rangeOffs := 14, 16+segX2, 16+2*segX2, 16+3*segX2
		for s := 0; s < segX2/2; s++ {
			end, start := uint32(be16(b, ends+s*2)), uint32(be16(b, starts+s*2))
			delta, ro := be16(b, deltas+s*2), int(be16(b, rangeOffs+s*2))
			if start > end || end-start > 0xffff {
				continue
			}
			for c := start; c <= end && c != 0xffff; c++ {
				var g uint16
				if ro == 0 {
					g = uint16(c) + delta
				} else {
					addr := rangeOffs + s*2 + ro + int(c-start)*2
					if addr+1 >= len(b) {
						continue
					}
					g = be16(b, addr)
					if g != 0 {
						g += delta
					}
				}
				if g != 0 {
					m[c] = g
				}
			}
		}
	case 6:
		first, count := uint32(be16(b, 6)), int(be16(b, 8))
		for i := 0; i < count; i++ {
			if g := be16(b, 10+i*2); g != 0 {
				m[first+uint32(i)] = g
			}
		}
	case 12:
		ngroups := int(be32(b, 12))
		for i := 0; i < ngroups && 16+i*12+12 <= len(b); i++ {
			rec := 16 + i*12
			start, end, gid := be32(b, rec), be32(b, rec+4), be32(b, rec+8)
			if end-start > 0xffff {
				end = start + 0xffff
			}
			for c := start; c <= end; c++ {
				m[c] = uint16(gid + (c - start))
			}
		}
	default:
		return nil
	}
	return m
}

func (f *sfnt) parsePost() {
	post := f.tables["post"]
	if len(post) < 34 || be32(post, 0) != 0x00020000 {
		return
	}
	n := int(be16(post, 32))
	idx := make([]uint16, n)
	for i := 0; i < n; i++ {
		idx[i] = be16(post, 34+i*2)
	}
	var names []string
	p := 34 + n*2
	for p < len(post) {
		l := int(post[p])
		p++
		if p+l > len(post) {
			break
		}
		names = append(names, string(post[p:p+l]))
		p += l
	}
	f.postNames = map[string]uint16{}
	for g, ix := range idx {
		var name string
		if ix >= 258 {
			if int(ix)-258 < len(names) {
				name = names[ix-258]
			}
		} else if int(ix) < len(macGlyphNames) {
			name = macGlyphNames[ix]
		}
		if name != "" {
			if _, dup := f.postNames[name]; !dup {
				f.postNames[name] = uint16(g)
			}
		}
	}
}

// --- writing ---

func checksum(b []byte) uint32 {
	var sum uint32
	for i := 0; i+3 < len(b); i += 4 {
		sum += be32(b, i)
	}
	rem := len(b) % 4
	if rem != 0 {
		var last [4]byte
		copy(last[:], b[len(b)-rem:])
		sum += be32(last[:], 0)
	}
	return sum
}

// buildSFNT assembles tables into a font file.
func buildSFNT(tables map[string][]byte, isCFF bool) []byte {
	tags := make([]string, 0, len(tables))
	for t := range tables {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	n := len(tags)
	// searchRange etc.
	entrySelector := 0
	for (1 << uint(entrySelector+1)) <= n {
		entrySelector++
	}
	searchRange := (1 << uint(entrySelector)) * 16
	out := make([]byte, 12+16*n)
	if isCFF {
		copy(out, "OTTO")
	} else {
		binary.BigEndian.PutUint32(out, 0x00010000)
	}
	binary.BigEndian.PutUint16(out[4:], uint16(n))
	binary.BigEndian.PutUint16(out[6:], uint16(searchRange))
	binary.BigEndian.PutUint16(out[8:], uint16(entrySelector))
	binary.BigEndian.PutUint16(out[10:], uint16(n*16-searchRange))
	headOff := -1
	for i, t := range tags {
		data := tables[t]
		off := len(out)
		if t == "head" {
			headOff = off
		}
		rec := 12 + i*16
		copy(out[rec:], t)
		binary.BigEndian.PutUint32(out[rec+4:], checksum(data))
		binary.BigEndian.PutUint32(out[rec+8:], uint32(off))
		binary.BigEndian.PutUint32(out[rec+12:], uint32(len(data)))
		out = append(out, data...)
		for len(out)%4 != 0 {
			out = append(out, 0)
		}
	}
	if headOff >= 0 && headOff+12 <= len(out) {
		binary.BigEndian.PutUint32(out[headOff+8:], 0)
		adj := 0xB1B0AFBA - checksum(out)
		binary.BigEndian.PutUint32(out[headOff+8:], adj)
	}
	return out
}

// buildCmapTable writes a cmap with a (3,1) format 4 subtable and, when needed, a (3,10) format 12 one.
func buildCmapTable(m map[uint32]uint16) []byte {
	codes := make([]uint32, 0, len(m))
	needs12 := false
	for c := range m {
		codes = append(codes, c)
		if c > 0xffff {
			needs12 = true
		}
	}
	sort.Slice(codes, func(i, j int) bool { return codes[i] < codes[j] })

	// format 4 segments (BMP only)
	type seg struct {
		start, end uint32
		gids       []uint16
	}
	var segs []seg
	for _, c := range codes {
		if c > 0xfffe {
			continue
		}
		if len(segs) > 0 && segs[len(segs)-1].end+1 == c {
			s := &segs[len(segs)-1]
			s.end = c
			s.gids = append(s.gids, m[c])
		} else {
			segs = append(segs, seg{c, c, []uint16{m[c]}})
		}
	}
	segs = append(segs, seg{0xffff, 0xffff, []uint16{0}})
	segCount := len(segs)
	var ends, starts, deltas, rangeOffs []byte
	var glyphIds []byte
	for i, s := range segs {
		ends = binary.BigEndian.AppendUint16(ends, uint16(s.end))
		starts = binary.BigEndian.AppendUint16(starts, uint16(s.start))
		// Use delta when gids are consecutive, else glyphIdArray.
		consecutive := true
		for k := 1; k < len(s.gids); k++ {
			if s.gids[k] != s.gids[0]+uint16(k) {
				consecutive = false
				break
			}
		}
		if consecutive {
			deltas = binary.BigEndian.AppendUint16(deltas, uint16(int(s.gids[0])-int(s.start)))
			rangeOffs = binary.BigEndian.AppendUint16(rangeOffs, 0)
		} else {
			deltas = binary.BigEndian.AppendUint16(deltas, 0)
			// offset from this rangeOffset entry to the glyphIdArray position
			ro := (segCount-i)*2 + len(glyphIds)
			rangeOffs = binary.BigEndian.AppendUint16(rangeOffs, uint16(ro))
			for _, g := range s.gids {
				glyphIds = binary.BigEndian.AppendUint16(glyphIds, g)
			}
		}
	}
	es := 0
	for (1 << uint(es+1)) <= segCount {
		es++
	}
	sr := (1 << uint(es)) * 2
	f4 := make([]byte, 0, 16+segCount*8+len(glyphIds))
	f4 = binary.BigEndian.AppendUint16(f4, 4)
	f4 = binary.BigEndian.AppendUint16(f4, 0) // length placeholder
	f4 = binary.BigEndian.AppendUint16(f4, 0)
	f4 = binary.BigEndian.AppendUint16(f4, uint16(segCount*2))
	f4 = binary.BigEndian.AppendUint16(f4, uint16(sr))
	f4 = binary.BigEndian.AppendUint16(f4, uint16(es))
	f4 = binary.BigEndian.AppendUint16(f4, uint16(segCount*2-sr))
	f4 = append(f4, ends...)
	f4 = binary.BigEndian.AppendUint16(f4, 0)
	f4 = append(f4, starts...)
	f4 = append(f4, deltas...)
	f4 = append(f4, rangeOffs...)
	f4 = append(f4, glyphIds...)
	binary.BigEndian.PutUint16(f4[2:], uint16(len(f4)))

	var f12 []byte
	if needs12 {
		type group struct{ start, end, gid uint32 }
		var groups []group
		for _, c := range codes {
			g := uint32(m[c])
			if len(groups) > 0 && groups[len(groups)-1].end+1 == c && groups[len(groups)-1].gid+(c-groups[len(groups)-1].start) == g {
				groups[len(groups)-1].end = c
			} else {
				groups = append(groups, group{c, c, g})
			}
		}
		f12 = binary.BigEndian.AppendUint16(f12, 12)
		f12 = binary.BigEndian.AppendUint16(f12, 0)
		f12 = binary.BigEndian.AppendUint32(f12, uint32(16+len(groups)*12))
		f12 = binary.BigEndian.AppendUint32(f12, 0)
		f12 = binary.BigEndian.AppendUint32(f12, uint32(len(groups)))
		for _, g := range groups {
			f12 = binary.BigEndian.AppendUint32(f12, g.start)
			f12 = binary.BigEndian.AppendUint32(f12, g.end)
			f12 = binary.BigEndian.AppendUint32(f12, g.gid)
		}
	}
	nsub := 1
	if f12 != nil {
		nsub = 2
	}
	out := make([]byte, 0)
	out = binary.BigEndian.AppendUint16(out, 0)
	out = binary.BigEndian.AppendUint16(out, uint16(nsub))
	hdr := 4 + nsub*8
	out = binary.BigEndian.AppendUint16(out, 3)
	out = binary.BigEndian.AppendUint16(out, 1)
	out = binary.BigEndian.AppendUint32(out, uint32(hdr))
	if f12 != nil {
		out = binary.BigEndian.AppendUint16(out, 3)
		out = binary.BigEndian.AppendUint16(out, 10)
		out = binary.BigEndian.AppendUint32(out, uint32(hdr+len(f4)))
	}
	out = append(out, f4...)
	out = append(out, f12...)
	return out
}

// buildNameTable writes a name table with the given family and style.
func buildNameTable(family, style string) []byte {
	full := family
	if style != "Regular" {
		full += " " + style
	}
	recs := []struct {
		id  uint16
		val string
	}{{1, family}, {2, style}, {3, family + ";" + style}, {4, full}, {5, "Version 1.0"}, {6, sanitizePS(family + "-" + style)}}
	var strs []byte
	out := make([]byte, 0)
	out = binary.BigEndian.AppendUint16(out, 0)
	out = binary.BigEndian.AppendUint16(out, uint16(len(recs)))
	out = binary.BigEndian.AppendUint16(out, uint16(6+len(recs)*12))
	for _, r := range recs {
		var u16 []byte
		for _, ch := range r.val {
			if ch > 0xffff {
				ch = '?'
			}
			u16 = binary.BigEndian.AppendUint16(u16, uint16(ch))
		}
		out = binary.BigEndian.AppendUint16(out, 3)
		out = binary.BigEndian.AppendUint16(out, 1)
		out = binary.BigEndian.AppendUint16(out, 0x0409)
		out = binary.BigEndian.AppendUint16(out, r.id)
		out = binary.BigEndian.AppendUint16(out, uint16(len(u16)))
		out = binary.BigEndian.AppendUint16(out, uint16(len(strs)))
		strs = append(strs, u16...)
	}
	return append(out, strs...)
}

func sanitizePS(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c > 32 && c < 127 && c != '[' && c != ']' && c != '(' && c != ')' && c != '{' && c != '}' && c != '<' && c != '>' && c != '/' && c != '%' {
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		return "Font"
	}
	return string(out)
}

// buildOS2Table writes a version 3 OS/2 table.
func buildOS2Table(weight int, italic bool, ascent, descent int16, unitsPerEm int) []byte {
	b := make([]byte, 96)
	binary.BigEndian.PutUint16(b[0:], 3)
	binary.BigEndian.PutUint16(b[2:], uint16(unitsPerEm/2)) // xAvgCharWidth
	binary.BigEndian.PutUint16(b[4:], uint16(weight))
	binary.BigEndian.PutUint16(b[6:], 5) // usWidthClass normal
	// subscript/superscript metrics
	em := uint16(unitsPerEm)
	for i, v := range []uint16{em * 65 / 100, em * 60 / 100, 0, em * 75 / 1000 * 10 / 10, em * 65 / 100, em * 60 / 100, 0, em * 35 / 100, em * 5 / 100, em * 26 / 100} {
		binary.BigEndian.PutUint16(b[10+i*2:], v)
	}
	copy(b[58:], "BDF ")
	var sel uint16 = 1 << 6 // REGULAR
	if italic {
		sel = 1 << 0
	}
	if weight >= 600 {
		sel |= 1 << 5
	}
	binary.BigEndian.PutUint16(b[62:], sel)
	binary.BigEndian.PutUint16(b[64:], 0x20)
	binary.BigEndian.PutUint16(b[66:], 0xffff)
	binary.BigEndian.PutUint16(b[68:], uint16(ascent))
	binary.BigEndian.PutUint16(b[70:], uint16(descent))
	binary.BigEndian.PutUint16(b[72:], 0)
	binary.BigEndian.PutUint16(b[74:], uint16(ascent))
	binary.BigEndian.PutUint16(b[76:], uint16(-descent))
	binary.BigEndian.PutUint32(b[78:], 1)
	binary.BigEndian.PutUint32(b[82:], 0)
	binary.BigEndian.PutUint16(b[86:], uint16(ascent*2/3)) // sxHeight
	binary.BigEndian.PutUint16(b[88:], uint16(ascent*3/4)) // sCapHeight
	binary.BigEndian.PutUint16(b[90:], 0)
	binary.BigEndian.PutUint16(b[92:], 0x20)
	binary.BigEndian.PutUint16(b[94:], 0)
	return b
}

func buildPostTable() []byte {
	b := make([]byte, 32)
	binary.BigEndian.PutUint32(b, 0x00030000)
	return b
}

func buildHeadTable(unitsPerEm int, indexToLoc int16, bbox [4]int16) []byte {
	b := make([]byte, 54)
	binary.BigEndian.PutUint32(b[0:], 0x00010000)
	binary.BigEndian.PutUint32(b[4:], 0x00010000)
	binary.BigEndian.PutUint32(b[12:], 0x5F0F3CF5)
	binary.BigEndian.PutUint16(b[16:], 0x000b)
	binary.BigEndian.PutUint16(b[18:], uint16(unitsPerEm))
	for i, v := range bbox {
		binary.BigEndian.PutUint16(b[36+i*2:], uint16(v))
	}
	binary.BigEndian.PutUint16(b[44:], 0)
	binary.BigEndian.PutUint16(b[46:], 8)
	binary.BigEndian.PutUint16(b[48:], 2)
	binary.BigEndian.PutUint16(b[50:], uint16(indexToLoc))
	return b
}

func buildHheaTable(ascent, descent int16, numHMetrics int, advanceMax uint16) []byte {
	b := make([]byte, 36)
	binary.BigEndian.PutUint32(b[0:], 0x00010000)
	binary.BigEndian.PutUint16(b[4:], uint16(ascent))
	binary.BigEndian.PutUint16(b[6:], uint16(descent))
	binary.BigEndian.PutUint16(b[10:], advanceMax)
	binary.BigEndian.PutUint16(b[18:], 1)
	binary.BigEndian.PutUint16(b[34:], uint16(numHMetrics))
	return b
}

func buildHmtxTable(advances []uint16) []byte {
	b := make([]byte, 0, len(advances)*4)
	for _, a := range advances {
		b = binary.BigEndian.AppendUint16(b, a)
		b = binary.BigEndian.AppendUint16(b, 0)
	}
	return b
}

func buildMaxpTable(numGlyphs int, cff bool) []byte {
	if cff {
		b := make([]byte, 6)
		binary.BigEndian.PutUint32(b, 0x00005000)
		binary.BigEndian.PutUint16(b[4:], uint16(numGlyphs))
		return b
	}
	b := make([]byte, 32)
	binary.BigEndian.PutUint32(b, 0x00010000)
	binary.BigEndian.PutUint16(b[4:], uint16(numGlyphs))
	binary.BigEndian.PutUint16(b[6:], 512)  // maxPoints
	binary.BigEndian.PutUint16(b[8:], 64)   // maxContours
	binary.BigEndian.PutUint16(b[10:], 512) // maxCompositePoints
	binary.BigEndian.PutUint16(b[12:], 64)
	binary.BigEndian.PutUint16(b[14:], 2) // maxZones
	binary.BigEndian.PutUint16(b[16:], 64)
	binary.BigEndian.PutUint16(b[18:], 64)
	binary.BigEndian.PutUint16(b[20:], 64)
	binary.BigEndian.PutUint16(b[22:], 64)
	binary.BigEndian.PutUint16(b[24:], 64)
	binary.BigEndian.PutUint16(b[26:], 64)
	binary.BigEndian.PutUint16(b[28:], 64)
	binary.BigEndian.PutUint16(b[30:], 64)
	return b
}

// rebuild produces a font program with a synthesized cmap (unicode → gid).
// Existing glyph data, metrics and hinting tables are kept; cmap, name, OS/2
// and post are replaced so the result is sanitizer-friendly.
func (f *sfnt) rebuild(cmap map[uint32]uint16, family string, weight int, italic bool) []byte {
	tables := map[string][]byte{}
	keep := []string{"glyf", "loca", "head", "hhea", "hmtx", "maxp", "cvt ", "fpgm", "prep", "CFF ", "gasp", "vhea", "vmtx"}
	for _, t := range keep {
		if b, ok := f.tables[t]; ok {
			tables[t] = append([]byte(nil), b...)
		}
	}
	if f.isCFF {
		delete(tables, "glyf")
		delete(tables, "loca")
	}
	ascent, descent := f.ascent, f.descent
	if ascent == 0 && descent == 0 {
		ascent, descent = int16(f.unitsPerEm*8/10), int16(-f.unitsPerEm*2/10)
	}
	if tables["head"] == nil {
		tables["head"] = buildHeadTable(f.unitsPerEm, 0, [4]int16{0, int16(-f.unitsPerEm / 5), int16(f.unitsPerEm), int16(f.unitsPerEm)})
	} else if len(tables["head"]) >= 54 {
		// Clear flags bit 11.. and set macStyle for consistency with OS/2.
		var ms uint16
		if weight >= 600 {
			ms |= 1
		}
		if italic {
			ms |= 2
		}
		binary.BigEndian.PutUint16(tables["head"][44:], ms)
	}
	if tables["hhea"] == nil || tables["hmtx"] == nil {
		adv := make([]uint16, f.numGlyphs)
		for i := range adv {
			if i < len(f.advances) {
				adv[i] = f.advances[i]
			} else {
				adv[i] = uint16(f.unitsPerEm / 2)
			}
		}
		tables["hhea"] = buildHheaTable(ascent, descent, len(adv), uint16(f.unitsPerEm))
		tables["hmtx"] = buildHmtxTable(adv)
	}
	if tables["maxp"] == nil {
		tables["maxp"] = buildMaxpTable(f.numGlyphs, f.isCFF)
	}
	style := "Regular"
	switch {
	case weight >= 600 && italic:
		style = "Bold Italic"
	case weight >= 600:
		style = "Bold"
	case italic:
		style = "Italic"
	}
	tables["cmap"] = buildCmapTable(cmap)
	tables["name"] = buildNameTable(family, style)
	tables["OS/2"] = buildOS2Table(weight, italic, ascent, descent, f.unitsPerEm)
	tables["post"] = buildPostTable()
	return buildSFNT(tables, f.isCFF)
}

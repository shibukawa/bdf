package sfnt

import (
	"encoding/binary"
	"sort"
)

// Advance returns the advance width of a glyph in font units. Glyphs past
// numberOfHMetrics share the last advance.
func (f *Font) Advance(gid uint16) int {
	n := len(f.Advances)
	if n == 0 {
		return f.UnitsPerEm / 2
	}
	if int(gid) >= n {
		return int(f.Advances[n-1])
	}
	return int(f.Advances[gid])
}

// OS2 is the part of the OS/2 table that font selection and layout use.
type OS2 struct {
	Weight      int
	FsType      uint16
	FsSelection uint16
	WinAscent   int
	WinDescent  int
	TypoAscent  int
	TypoDescent int
	TypoLineGap int
}

// OS2 parses the OS/2 table; ok is false when it is missing.
func (f *Font) OS2() (o OS2, ok bool) {
	b := f.Tables["OS/2"]
	if len(b) < 78 {
		return o, false
	}
	o.Weight = int(be16(b, 4))
	o.FsType = be16(b, 8)
	o.FsSelection = be16(b, 62)
	o.TypoAscent = int(int16(be16(b, 68)))
	o.TypoDescent = int(int16(be16(b, 70)))
	o.TypoLineGap = int(int16(be16(b, 72)))
	o.WinAscent = int(be16(b, 74))
	o.WinDescent = int(be16(b, 76))
	return o, true
}

// Embeddable reports whether the OS/2 embedding permissions (fsType) allow
// embedding the outlines in a document, and whether they allow subsetting.
func (o OS2) Embeddable() (embed, subset bool) {
	if o.FsType&0x000f == 0x0002 || o.FsType&0x0200 != 0 {
		return false, false
	}
	return true, o.FsType&0x0100 == 0
}

// Underline returns the post table underline position (negative is below
// the baseline) and thickness in font units.
func (f *Font) Underline() (pos, thickness int) {
	b := f.Tables["post"]
	if len(b) < 12 {
		return -f.UnitsPerEm / 10, f.UnitsPerEm / 20
	}
	pos, thickness = int(int16(be16(b, 8))), int(int16(be16(b, 10)))
	if thickness <= 0 {
		thickness = f.UnitsPerEm / 20
	}
	return pos, thickness
}

// Subset returns a TrueType font holding only the given glyphs (plus glyph
// 0 and the components of composite glyphs), renumbered in ascending order
// of their old index, and the old → new glyph index map. Tables that index
// glyphs and are not rewritten here (kerning, layout, vertical metrics,
// device metrics) are dropped; Rebuild then adds cmap, name, OS/2 and post.
// It returns nil for CFF-based fonts.
func (f *Font) Subset(gids []uint16) (*Font, map[uint16]uint16) {
	glyf, loca, head := f.Tables["glyf"], f.Tables["loca"], f.Tables["head"]
	if f.IsCFF || glyf == nil || loca == nil || len(head) < 54 || f.NumGlyphs == 0 {
		return nil, nil
	}
	long := be16(head, 50) != 0
	offset := func(g int) uint32 {
		if long {
			return be32(loca, g*4)
		}
		return uint32(be16(loca, g*2)) * 2
	}
	glyphData := func(g int) []byte {
		if g < 0 || g >= f.NumGlyphs {
			return nil
		}
		s, e := offset(g), offset(g+1)
		if s > e || int(e) > len(glyf) {
			return nil
		}
		return glyf[s:e]
	}
	needed := map[uint16]bool{}
	var visit func(g uint16)
	visit = func(g uint16) {
		if needed[g] || int(g) >= f.NumGlyphs {
			return
		}
		needed[g] = true
		forComponents(glyphData(int(g)), func(p int, idx uint16) { visit(idx) })
	}
	visit(0)
	for _, g := range gids {
		visit(g)
	}
	order := make([]uint16, 0, len(needed))
	for g := range needed {
		order = append(order, g)
	}
	sort.Slice(order, func(i, j int) bool { return order[i] < order[j] })
	remap := make(map[uint16]uint16, len(order))
	for i, g := range order {
		remap[g] = uint16(i)
	}

	var newGlyf []byte
	offsets := make([]uint32, 0, len(order)+1)
	for _, g := range order {
		offsets = append(offsets, uint32(len(newGlyf)))
		d := append([]byte(nil), glyphData(int(g))...)
		forComponents(d, func(p int, idx uint16) {
			binary.BigEndian.PutUint16(d[p:], remap[idx])
		})
		newGlyf = append(newGlyf, d...)
		for len(newGlyf)%4 != 0 {
			newGlyf = append(newGlyf, 0)
		}
	}
	offsets = append(offsets, uint32(len(newGlyf)))
	useLong := len(newGlyf) >= 0x20000
	var newLoca []byte
	for _, o := range offsets {
		if useLong {
			newLoca = binary.BigEndian.AppendUint32(newLoca, o)
		} else {
			newLoca = binary.BigEndian.AppendUint16(newLoca, uint16(o/2))
		}
	}
	newHead := append([]byte(nil), head...)
	if useLong {
		binary.BigEndian.PutUint16(newHead[50:], 1)
	} else {
		binary.BigEndian.PutUint16(newHead[50:], 0)
	}

	// hmtx: one full metric per glyph.
	hmtx := f.Tables["hmtx"]
	nhm := len(f.Advances)
	lsb := func(g int) uint16 {
		if g < nhm {
			return be16(hmtx, g*4+2)
		}
		return be16(hmtx, nhm*4+(g-nhm)*2)
	}
	var newHmtx []byte
	var advMax uint16
	advances := make([]uint16, len(order))
	for i, g := range order {
		a := uint16(f.Advance(g))
		advances[i] = a
		advMax = max(advMax, a)
		newHmtx = binary.BigEndian.AppendUint16(newHmtx, a)
		newHmtx = binary.BigEndian.AppendUint16(newHmtx, lsb(int(g)))
	}
	newHhea := append([]byte(nil), f.Tables["hhea"]...)
	if len(newHhea) >= 36 {
		binary.BigEndian.PutUint16(newHhea[10:], advMax)
		binary.BigEndian.PutUint16(newHhea[34:], uint16(len(order)))
	}
	newMaxp := append([]byte(nil), f.Tables["maxp"]...)
	if len(newMaxp) >= 6 {
		binary.BigEndian.PutUint16(newMaxp[4:], uint16(len(order)))
	}

	out := &Font{Tables: map[string][]byte{}, NumGlyphs: len(order), UnitsPerEm: f.UnitsPerEm,
		Ascent: f.Ascent, Descent: f.Descent, LineGap: f.LineGap, Advances: advances}
	for _, t := range []string{"cvt ", "fpgm", "prep", "gasp", "OS/2", "post"} {
		if b, ok := f.Tables[t]; ok {
			out.Tables[t] = b
		}
	}
	out.Tables["glyf"], out.Tables["loca"], out.Tables["head"] = newGlyf, newLoca, newHead
	out.Tables["hmtx"], out.Tables["hhea"], out.Tables["maxp"] = newHmtx, newHhea, newMaxp
	return out, remap
}

// forComponents calls fn with the byte offset and glyph index of each
// component of a composite glyph.
func forComponents(d []byte, fn func(p int, idx uint16)) {
	if len(d) < 10 || int16(be16(d, 0)) >= 0 {
		return
	}
	p := 10
	for p+4 <= len(d) {
		flags := be16(d, p)
		fn(p+2, be16(d, p+2))
		p += 4
		if flags&1 != 0 { // ARG_1_AND_2_ARE_WORDS
			p += 4
		} else {
			p += 2
		}
		switch {
		case flags&8 != 0: // WE_HAVE_A_SCALE
			p += 2
		case flags&0x40 != 0: // WE_HAVE_AN_X_AND_Y_SCALE
			p += 4
		case flags&0x80 != 0: // WE_HAVE_A_TWO_BY_TWO
			p += 8
		}
		if flags&0x20 == 0 { // MORE_COMPONENTS
			return
		}
	}
}

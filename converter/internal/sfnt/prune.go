package sfnt

import "encoding/binary"

// PruneGlyphs empties every TrueType glyph that is not in keep (plus the
// components of kept composite glyphs and glyph 0). Glyph indices are left
// unchanged so cmap, hmtx and post stay valid; the glyf table shrinks and
// loca is rewritten.
func (f *Font) PruneGlyphs(keep map[uint16]bool) {
	glyf, loca, head := f.Tables["glyf"], f.Tables["loca"], f.Tables["head"]
	if glyf == nil || loca == nil || len(head) < 54 || f.NumGlyphs == 0 {
		return
	}
	long := be16(head, 50) != 0
	offsets := make([]uint32, f.NumGlyphs+1)
	for i := range offsets {
		if long {
			offsets[i] = be32(loca, i*4)
		} else {
			offsets[i] = uint32(be16(loca, i*2)) * 2
		}
	}
	glyphData := func(g int) []byte {
		if g < 0 || g >= f.NumGlyphs {
			return nil
		}
		s, e := offsets[g], offsets[g+1]
		if s > e || int(e) > len(glyf) {
			return nil
		}
		return glyf[s:e]
	}
	// Close the set under composite glyph references.
	needed := map[uint16]bool{0: true}
	var visit func(g uint16)
	visit = func(g uint16) {
		if needed[g] || int(g) >= f.NumGlyphs {
			return
		}
		needed[g] = true
		d := glyphData(int(g))
		if len(d) < 10 || int16(be16(d, 0)) >= 0 {
			return
		}
		p := 10
		for {
			if p+4 > len(d) {
				return
			}
			flags, idx := be16(d, p), be16(d, p+2)
			visit(idx)
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
	for g := range keep {
		visit(g)
	}
	newGlyf := make([]byte, 0, len(glyf)/2)
	newOffsets := make([]uint32, f.NumGlyphs+1)
	for g := 0; g < f.NumGlyphs; g++ {
		newOffsets[g] = uint32(len(newGlyf))
		if needed[uint16(g)] {
			newGlyf = append(newGlyf, glyphData(g)...)
			for len(newGlyf)%4 != 0 {
				newGlyf = append(newGlyf, 0)
			}
		}
	}
	newOffsets[f.NumGlyphs] = uint32(len(newGlyf))
	// Rewrite loca in the long format when needed, otherwise keep the short one.
	useLong := long || len(newGlyf) >= 0x20000
	var newLoca []byte
	for _, o := range newOffsets {
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
	f.Tables["glyf"], f.Tables["loca"], f.Tables["head"] = newGlyf, newLoca, newHead
}

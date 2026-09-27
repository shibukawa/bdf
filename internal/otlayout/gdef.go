package otlayout

// GDEF is the glyph definition table: which glyphs are bases, ligatures,
// marks and components, the classes marks attach by, and the glyph sets
// lookups can restrict themselves to.
type GDEF struct {
	Classes    map[uint16]int // glyph → GlyphBase, GlyphLigature, GlyphMark or GlyphComponent
	MarkAttach map[uint16]int // mark → mark attachment class
	MarkSets   [][]uint16
}

// Glyph classes of the GDEF table.
const (
	GlyphBase      = 1
	GlyphLigature  = 2
	GlyphMark      = 3
	GlyphComponent = 4
)

// ParseGDEF reads a GDEF table; it returns nil when b is not one.
func ParseGDEF(b []byte) *GDEF {
	d := data(b)
	if len(d) < 12 || d.u16(0) != 1 {
		return nil
	}
	g := &GDEF{Classes: classDef(d.at(d.u16(4))), MarkAttach: classDef(d.at(d.u16(10)))}
	if d.u16(2) >= 2 && len(d) >= 14 {
		if ms := d.at(d.u16(12)); ms != nil && ms.u16(0) == 1 {
			n := ms.count(ms.u16(2), 4, 4)
			for i := range n {
				g.MarkSets = append(g.MarkSets, coverage(ms.at(ms.u32(4+i*4))))
			}
		}
	}
	return g
}

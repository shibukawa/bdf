package otlayout

// GDEF is the glyph definition table: which glyphs are bases, ligatures,
// marks and components, the classes marks attach by, and the glyph sets
// lookups can restrict themselves to.
type GDEF struct {
	Classes    map[uint16]int // glyph → GlyphBase, GlyphLigature, GlyphMark or GlyphComponent
	MarkAttach map[uint16]int // mark → mark attachment class
	MarkSets   [][]uint16     // the first sets, when they hold more than a million glyphs together
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
			// sets may share a coverage table: those past maxElements
			// glyphs together are left out
			n, left := ms.count(ms.u16(2), 4, 4), maxElements
			for i := range n {
				cov := coverage(ms.at(ms.u32(4 + i*4)))
				if !take(&left, 1+len(cov)) {
					break
				}
				g.MarkSets = append(g.MarkSets, cov)
			}
		}
	}
	return g
}

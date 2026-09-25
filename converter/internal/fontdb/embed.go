package fontdb

import (
	"sort"
)

// Program builds the font program to embed for the given characters: a
// subset holding just their glyphs (TrueType outlines), or the whole font
// (CFF outlines, fonts whose license forbids subsetting, or all when
// noSubset is set), with a cmap for exactly those characters. ok is false
// when the license forbids embedding.
func (l *Loaded) Program(runes []rune, noSubset bool) (data []byte, ok bool) {
	if !l.embed {
		return nil, false
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	cmap := map[uint32]uint16{}
	var gids []uint16
	for _, r := range runes {
		if g, ok := l.Glyph(r); ok {
			cmap[uint32(r)] = g
			gids = append(gids, g)
		}
	}
	f := l.Font
	weight, italic := l.Face.Weight, l.Face.Italic
	if !noSubset && l.subset {
		if sub, remap := f.Subset(gids); sub != nil {
			for r, g := range cmap {
				cmap[r] = remap[g]
			}
			return sub.Rebuild(cmap, l.Face.Family, weight, italic), true
		}
	}
	if noSubset {
		for r, g := range f.Cmap {
			cmap[r] = g
		}
	}
	return f.Rebuild(cmap, l.Face.Family, weight, italic), true
}

// Size is the size of the font file the face comes from.
func (l *Loaded) Size() int { return len(l.Data) }

// CanSubset reports whether Program can drop unused glyphs.
func (l *Loaded) CanSubset() bool { return l.subset && !l.Font.IsCFF && l.Font.Tables["glyf"] != nil }

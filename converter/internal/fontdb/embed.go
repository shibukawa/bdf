package fontdb

import (
	"sort"

	"github.com/shibukawa/bdf/converter/internal/sfnt"
)

// fsPreviewPrint is the OS/2 fsType written when the original font does not
// have one: a BDF document is only viewed and printed.
const fsPreviewPrint = 0x0004

// Program builds the font program to embed for the given characters: a
// subset holding just their glyphs (TrueType outlines), or the whole font
// (CFF outlines, fonts whose license forbids subsetting, or all when
// noSubset is set), with a cmap for exactly those characters. The OS/2
// fsType and the copyright and license strings of the font carry over. ok is
// false when the license forbids embedding; ignoreFSType embeds (and
// subsets) such fonts anyway, for callers that hold the rights to.
func (l *Loaded) Program(runes []rune, noSubset, ignoreFSType bool) (data []byte, ok bool) {
	if !l.embed && !ignoreFSType {
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
	info := sfnt.FontInfo{Family: l.Face.Family, Weight: l.Face.Weight, Italic: l.Face.Italic, FSType: f.FSType, Notices: f.Notices}
	if !f.HasFSType {
		info.FSType = fsPreviewPrint
	}
	if !noSubset && (l.subset || ignoreFSType) {
		if sub, remap := f.Subset(gids); sub != nil {
			for r, g := range cmap {
				cmap[r] = remap[g]
			}
			return sub.Rebuild(cmap, info), true
		}
	}
	if noSubset {
		for r, g := range f.Cmap {
			cmap[r] = g
		}
	}
	return f.Rebuild(cmap, info), true
}

// Size is the size of the font file the face comes from.
func (l *Loaded) Size() int { return len(l.Data) }

// CanSubset reports whether Program can drop unused glyphs (ignoreFSType as
// for Program).
func (l *Loaded) CanSubset(ignoreFSType bool) bool {
	return (l.subset || ignoreFSType) && !l.Font.IsCFF && l.Font.Tables["glyf"] != nil
}

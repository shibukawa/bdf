package fontdb

import (
	"sort"

	"github.com/shibukawa/bdf/converter/internal/sfnt"
)

// fsPreviewPrint is the OS/2 fsType written when the original font does not
// have one: a BDF document is only viewed and printed.
const fsPreviewPrint = 0x0004

// Program builds the font program to embed for the given characters: a
// subset holding just their glyphs, or the whole font (fonts whose license
// forbids subsetting or whose outlines cannot be subset, or all when
// noSubset is set), with a cmap for exactly those characters. glyphs maps
// more characters (private use ones) to glyphs that no character maps to,
// such as the larger variants of a math font's delimiters. The OS/2
// fsType and the copyright and license strings of the font carry over. ok is
// false when the license forbids embedding; ignoreFSType embeds (and
// subsets) such fonts anyway, for callers that hold the rights to.
func (l *Loaded) Program(runes []rune, glyphs map[rune]uint16, noSubset, ignoreFSType bool) (data []byte, ok bool) {
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
	extra := make([]rune, 0, len(glyphs))
	for r := range glyphs {
		extra = append(extra, r)
	}
	sort.Slice(extra, func(i, j int) bool { return extra[i] < extra[j] })
	for _, r := range extra {
		if g := glyphs[r]; int(g) < l.Font.NumGlyphs {
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
		sub, remap := f.Subset(gids)
		if sub == nil {
			sub, remap = f.SubsetCFF(gids)
		}
		if sub != nil {
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
	f := l.Font
	return (l.subset || ignoreFSType) && (!f.IsCFF && f.Tables["glyf"] != nil || f.IsCFF && f.Tables["CFF "] != nil)
}

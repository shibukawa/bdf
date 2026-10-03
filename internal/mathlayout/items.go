package mathlayout

import "github.com/shibukawa/bdf"

// ItemKind is what an Item draws.
type ItemKind uint8

const (
	ItemGlyph  ItemKind = iota // a glyph (or, without a formula font, a character) of a face
	ItemRule                   // a filled rectangle: a fraction bar, the rule of a radical
	ItemStroke                 // a stroked path: an enclosure, a strike, a delimiter drawn without a formula font
)

// Item is something a laid out formula draws, relative to the left end of
// its baseline (y grows downwards, points).
type Item struct {
	Kind ItemKind
	// X and Y are where a glyph's baseline starts, the top left corner of
	// a rule, and the origin of a stroke's path.
	X, Y  float64
	Color bdf.Color

	// a glyph
	Face *Face
	// Glyph is the glyph in the face; HasGlyph is false for a character
	// the face lacks, which a host draws with a font of its own.
	Glyph    uint16
	HasGlyph bool
	Size     float64 // the font size
	// Text is the character that draws the glyph: the character itself,
	// or the one Fonts.GlyphRune gave a glyph no character maps to (then
	// Alt is the character the glyph is a form of).
	Text, Alt string
	Advance   float64
	// ScaleX and ScaleY stretch a glyph that is grown without a formula
	// font (1 otherwise).
	ScaleX, ScaleY float64

	// W and H are the size of a rule; W is the line width of a stroke.
	W, H float64
	Path *bdf.Path // of a stroke
}

// Items returns what b draws, in drawing order.
func (b *Box) Items() []Item {
	out := make([]Item, len(b.items))
	for i := range b.items {
		it := &b.items[i]
		out[i] = Item{Kind: ItemKind(it.kind), X: it.x, Y: it.y, Color: it.color,
			Face: it.fc, Glyph: it.gid, HasGlyph: it.hasGID, Size: it.size, Text: it.text, Alt: it.alt, Advance: it.adv,
			ScaleX: it.sx, ScaleY: it.sy, W: it.w, H: it.h, Path: it.path}
	}
	return out
}

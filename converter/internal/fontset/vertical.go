package fontset

import (
	"strings"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// East Asian vertical text: the converters lay a vertical line out as a
// horizontal one turned by 90°, then turn the upright characters back one
// by one.

// Upright reports whether a character stands upright in East Asian
// vertical text; Latin text, long vowel marks, dashes, brackets and
// ellipses are turned with the line instead.
func Upright(r rune) bool {
	if !fontdb.IsCJK(r) {
		return false
	}
	return !strings.ContainsRune("ー－―‐〜～…‥（）「」『』【】〔〕［］｛｝〈〉《》〘〙〖〗＝｜＿", r)
}

// verticalForms are the vertical presentation forms of East Asian
// punctuation; without them the horizontal glyph moves to the top right of
// its square, where vertical text puts it.
var verticalForms = map[rune]rune{'、': '︑', '。': '︒', '，': '︐', '．': '︒', '：': '︓', '；': '︔', '！': '︕', '？': '︖'}

// VerticalForm returns the vertical presentation form of a punctuation
// mark (、 → ︑); false for characters that have none.
func VerticalForm(r rune) (rune, bool) {
	v, ok := verticalForms[r]
	return v, ok
}

// Package linebreak holds the line breaking rules of the converters that
// lay text out themselves (Office documents): the simple rules Office
// applies. Lines break after runs of spaces, between East Asian characters
// (except before closing punctuation and small kana, and after opening
// brackets: kinsoku), between East Asian and other text, and after hyphens
// inside words.
package linebreak

import (
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// noStart are the characters that may not start a line, noEnd those that
// may not end one.
const (
	noStart = "、。，．,.:;?!)]}」』】〕〉》）］｝〙〗〟”’ゝゞヽヾーぁぃぅぇぉっゃゅょゎゕゖァィゥェォッャュョヮヵヶ々〻‐゠–〜？！：；・…‥%％°℃"
	noEnd   = "([{「『【〔〈《（［｛〘〖〝“‘¥$£＄￥"
)

// IsSpace reports whether r is a space that lines break after (and that is
// not drawn at the end of a line).
func IsSpace(r rune) bool {
	return r == ' ' || r == '\u3000' || r == '\u2009' || r == '\u2002' || r == '\u2003'
}

// NoStart reports whether r may not start a line.
func NoStart(r rune) bool { return strings.ContainsRune(noStart, r) }

// NoEnd reports whether r may not end a line.
func NoEnd(r rune) bool { return strings.ContainsRune(noEnd, r) }

// Allowed reports whether a line may break between a and b.
func Allowed(a, b rune) bool {
	if IsSpace(a) {
		return !IsSpace(b)
	}
	if IsSpace(b) || a == '\u00a0' || b == '\u00a0' {
		return false
	}
	if NoStart(b) || NoEnd(a) {
		return false
	}
	if fontdb.IsCJK(a) || fontdb.IsCJK(b) {
		return true
	}
	if (a == '-' || a == '‐') && unicode.IsLetter(b) {
		return true
	}
	return false
}

// Joins reports whether the text on both sides of a break between a and b
// joins without a space when the lines are read as one (East Asian text).
func Joins(a, b rune) bool {
	return !IsSpace(a) && (fontdb.IsCJK(a) || fontdb.IsCJK(b))
}

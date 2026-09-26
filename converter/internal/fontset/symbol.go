package fontset

import "github.com/shibukawa/bdf/converter/internal/fontdb"

// IsSymbol reports whether a family is a symbol font (Wingdings, Symbol,
// Webdings), whose character codes are not Unicode.
func IsSymbol(name string) bool {
	switch fontdb.Normalize(name) {
	case "wingdings", "wingdings2", "wingdings3", "symbol", "webdings":
		return true
	}
	return false
}

// MapSymbol maps the characters of the symbol fonts that bullets and
// Windows metafiles commonly use to their Unicode equivalents; codes in
// the private use area U+F000–U+F0FF stand for the font's codes 0x00–0xFF.
func MapSymbol(font, text string) string {
	n := fontdb.Normalize(font)
	var table map[rune]rune
	switch n {
	case "wingdings":
		table = wingdings
	case "symbol":
		table = symbolFont
	default:
		return text
	}
	var out []rune
	for _, r := range text {
		if r >= 0xf000 && r <= 0xf0ff {
			r -= 0xf000
		}
		if m, ok := table[r]; ok {
			r = m
		}
		out = append(out, r)
	}
	return string(out)
}

var wingdings = map[rune]rune{
	0x6c: '●', 0x6d: '❍', 0x6e: '■', 0x6f: '□', 0x70: '◻', 0x71: '❑', 0x72: '❒', 0x73: '⬧', 0x74: '⧫',
	0x75: '◆', 0x76: '❖', 0x77: '⬥', 0x78: '⌧', 0xa7: '▪', 0xa8: '◻', 0x9f: '•', 0xa1: '○', 0xa2: '⭕',
	0xd8: '➢', 0xe0: '➔', 0xe8: '➔', 0xfc: '✔', 0xfb: '✘', 0xfd: '☒', 0xfe: '☑', 0xf0: '⇨', 0xef: '⇦',
	0xd9: '⮙', 0xda: '⮛', 0x46: '☞', 0x4a: '☺', 0x4c: '☹', 0xab: '★', 0xb2: '✧', 0x2a: '✉', 0x3e: '✇',
	0x3f: '✍', 0x41: '✌', 0x43: '👍', 0x44: '👎', 0x51: '✈', 0x52: '☼', 0x54: '❄',
}

var symbolFont = map[rune]rune{
	0xb7: '•', 0xa8: '♣', 0xa9: '♦', 0xaa: '♥', 0xab: '♠', 0xae: '→', 0xde: '⇒', 0xd8: '¬', 0x2d: '−',
	0xb0: '°', 0xb1: '±', 0xb4: '×', 0xb8: '÷', 0xa5: '∞', 0x61: 'α', 0x62: 'β', 0x67: 'γ', 0x64: 'δ',
	0x70: 'π', 0x6d: 'μ', 0x53: 'Σ', 0x57: 'Ω', 0xe0: '◊',
}

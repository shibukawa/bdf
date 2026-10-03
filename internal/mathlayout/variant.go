package mathlayout

import "sync"

// The Mathematical Alphanumeric Symbols draw the letters and digits of a
// formula in their styles (italic variables, bold vectors, double-struck
// sets). Some letters of a style were encoded earlier, in the Letterlike
// Symbols block, and are holes in theirs.

// latinStart is the first letter (A) of each style.
var latinStart = map[Variant]rune{
	VarBold: 0x1D400, VarItalic: 0x1D434, VarBoldItalic: 0x1D468, VarScript: 0x1D49C, VarBoldScript: 0x1D4D0,
	VarFraktur: 0x1D504, VarDoubleStruck: 0x1D538, VarBoldFraktur: 0x1D56C, VarSansSerif: 0x1D5A0,
	VarSansSerifBold: 0x1D5D4, VarSansSerifItalic: 0x1D608, VarSansSerifBoldItalic: 0x1D63C, VarMonospace: 0x1D670,
}

// greekStart is the first Greek letter (Α) of each style that has Greek.
var greekStart = map[Variant]rune{
	VarBold: 0x1D6A8, VarItalic: 0x1D6E2, VarBoldItalic: 0x1D71C, VarSansSerifBold: 0x1D756, VarSansSerifBoldItalic: 0x1D790,
}

// digitStart is the digit zero of each style that has digits.
var digitStart = map[Variant]rune{
	VarBold: 0x1D7CE, VarDoubleStruck: 0x1D7D8, VarSansSerif: 0x1D7E2, VarSansSerifBold: 0x1D7EC, VarMonospace: 0x1D7F6,
}

// holes are the letters of a style encoded in the Letterlike Symbols.
var holes = map[Variant]map[rune]rune{
	VarItalic: {'h': 0x210E},
	VarScript: {'B': 0x212C, 'E': 0x2130, 'F': 0x2131, 'H': 0x210B, 'I': 0x2110, 'L': 0x2112, 'M': 0x2133, 'R': 0x211B,
		'e': 0x212F, 'g': 0x210A, 'o': 0x2134},
	VarFraktur:      {'C': 0x212D, 'H': 0x210C, 'I': 0x2111, 'R': 0x211C, 'Z': 0x2128},
	VarDoubleStruck: {'C': 0x2102, 'H': 0x210D, 'N': 0x2115, 'P': 0x2119, 'Q': 0x211A, 'R': 0x211D, 'Z': 0x2124},
}

// greekIndex is the position of a Greek letter (or of the symbols that
// follow the Greek letters) in a style's block.
func greekIndex(r rune) (int, bool) {
	switch {
	case r >= 0x391 && r <= 0x3A1:
		return int(r - 0x391), true
	case r == 0x3F4: // ϴ
		return 17, true
	case r >= 0x3A3 && r <= 0x3A9:
		return int(r-0x3A3) + 18, true
	case r == '∇':
		return 25, true
	case r >= 0x3B1 && r <= 0x3C9:
		return int(r-0x3B1) + 26, true
	}
	switch r {
	case '∂':
		return 51, true
	case 0x3F5: // ϵ
		return 52, true
	case 0x3D1: // ϑ
		return 53, true
	case 0x3F0: // ϰ
		return 54, true
	case 0x3D5: // ϕ
		return 55, true
	case 0x3F1: // ϱ
		return 56, true
	case 0x3D6: // ϖ
		return 57, true
	}
	return 0, false
}

// styled returns the character that draws r in variant v, or r when the
// style has no form of it.
func styled(r rune, v Variant) rune {
	switch {
	case r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z':
		if h, ok := holes[v][r]; ok {
			return h
		}
		start, ok := latinStart[v]
		if !ok {
			return r
		}
		if r >= 'a' {
			return start + 26 + (r - 'a')
		}
		return start + (r - 'A')
	case r >= '0' && r <= '9':
		if start, ok := digitStart[v]; ok {
			return start + (r - '0')
		}
		if v == VarBoldItalic || v == VarBoldScript || v == VarBoldFraktur {
			return 0x1D7CE + (r - '0')
		}
		return r
	case r == 'ı' && v == VarItalic:
		return 0x1D6A4
	case r == 'ȷ' && v == VarItalic:
		return 0x1D6A5
	}
	if i, ok := greekIndex(r); ok {
		if start, ok := greekStart[v]; ok {
			return start + rune(i)
		}
	}
	return r
}

// upperGreek reports whether r is a capital Greek letter, which formulas
// set upright unless asked otherwise (as TeX and Office do).
func upperGreek(r rune) bool {
	return r >= 0x391 && r <= 0x3A9 || r == 0x3F4
}

// isBold reports whether a variant is one of the bold ones.
func (v Variant) isBold() bool {
	switch v {
	case VarBold, VarBoldItalic, VarBoldScript, VarBoldFraktur, VarSansSerifBold, VarSansSerifBoldItalic:
		return true
	}
	return false
}

// isItalic reports whether a variant slants its letters.
func (v Variant) isItalic() bool {
	switch v {
	case VarItalic, VarBoldItalic, VarSansSerifItalic, VarSansSerifBoldItalic:
		return true
	}
	return false
}

// plainLetter maps a letter or digit of the Mathematical Alphanumeric
// Symbols (or one of their holes) back to the plain character and its
// variant.
func plainLetter(r rune) (rune, Variant, bool) {
	if r < 0x2100 {
		return 0, 0, false
	}
	plainOnce.Do(func() {
		plain = map[rune][2]rune{}
		add := func(v Variant, base rune) {
			if s := styled(base, v); s != base {
				plain[s] = [2]rune{base, rune(v)}
			}
		}
		for v := VarItalic; v <= VarMonospace; v++ {
			for c := 'A'; c <= 'Z'; c++ {
				add(v, c)
				add(v, c+32)
			}
			for c := '0'; c <= '9'; c++ {
				add(v, c)
			}
			for _, g := range "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρςστυφχψω∂ϵϑϰϕϱϖ∇ϴ" {
				add(v, g)
			}
		}
		add(VarItalic, 'ı')
		add(VarItalic, 'ȷ')
	})
	p, ok := plain[r]
	return p[0], Variant(p[1]), ok
}

var (
	plainOnce sync.Once
	plain     map[rune][2]rune
)

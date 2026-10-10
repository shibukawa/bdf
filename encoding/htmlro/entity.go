package htmlro

import "unsafe"

// lookupEntity resolves a named character reference by its name, with the
// semicolon when it has one, in the table of table.go, and allocates
// nothing.
func lookupEntity(name []byte) (string, bool) {
	// A string view of the bytes, for comparisons that allocate on neither
	// compiler; it does not outlive the call.
	s := unsafe.String(unsafe.SliceData(name), len(name))
	lo, hi := 0, entityCount
	for lo < hi {
		m := int(uint(lo+hi) >> 1)
		switch n := entityNames[entityNameOffsets[m]:entityNameOffsets[m+1]]; {
		case n < s:
			lo = m + 1
		case n > s:
			hi = m
		default:
			return entityValues[entityValueOffsets[m]:entityValueOffsets[m+1]], true
		}
	}
	return "", false
}

// replacementTable maps the numeric references 0x80 to 0x9F to the
// characters Windows-1252 had there, as the specification requires.
// https://html.spec.whatwg.org/multipage/syntax.html#consume-a-character-reference
var replacementTable = [...]rune{
	'€', // 0x80
	'\u0081',
	'‚',
	'ƒ',
	'„',
	'…',
	'†',
	'‡',
	'ˆ',
	'‰',
	'Š',
	'‹',
	'Œ',
	'\u008D',
	'Ž',
	'\u008F',
	'\u0090',
	'‘',
	'’',
	'“',
	'”',
	'•',
	'–',
	'—',
	'˜',
	'™',
	'š',
	'›',
	'œ',
	'\u009D',
	'ž',
	'Ÿ', // 0x9F
}

// entityAt consumes a character reference at the start of s, whose first
// byte is '&': a named one as the value of its table row, a numeric one as
// the rune, and n bytes of s. n is 1 when there is no reference, and the
// '&' stands for itself. attribute is whether s is an attribute value,
// where a legacy reference followed by '=' or by more name characters is
// not one.
func entityAt(s []byte, attribute bool) (val string, r rune, n int) {
	// https://html.spec.whatwg.org/multipage/syntax.html#consume-a-character-reference
	i := 1
	if len(s) <= 1 {
		return "", '&', 1
	}

	if s[i] == '#' {
		if len(s) <= 2 { // We need to have at least "&#".
			return "", '&', 1
		}
		i++
		c := s[i]
		hex := false
		if c == 'x' || c == 'X' {
			hex = true
			i++
		}

		i0 := i
		x := '\x00'
		for i < len(s) {
			c = s[i]
			var d rune
			var mult rune
			if hex {
				mult = 16
				if '0' <= c && c <= '9' {
					d = rune(c) - '0'
				} else if 'a' <= c && c <= 'f' {
					d = rune(c) - 'a' + 10
				} else if 'A' <= c && c <= 'F' {
					d = rune(c) - 'A' + 10
				} else {
					break
				}
			} else {
				mult = 10
				if '0' <= c && c <= '9' {
					d = rune(c) - '0'
				} else {
					break
				}
			}
			if x <= 0x10FFFF {
				x = mult*x + d
			}
			i++
		}

		if i == i0 { // No characters matched.
			return "", '&', 1
		}

		if i < len(s) && s[i] == ';' {
			i++
		}

		if 0x80 <= x && x <= 0x9F {
			// Replace characters from Windows-1252 with UTF-8 equivalents.
			x = replacementTable[x-0x80]
		} else if x == 0 || (0xD800 <= x && x <= 0xDFFF) || x > 0x10FFFF {
			// Replace invalid characters with the replacement character.
			x = '�'
		}

		return "", x, i
	}

	// Consume the maximum number of characters possible, with the
	// consumed characters matching one of the named references.
	for i < len(s) {
		c := s[i]
		i++
		// Lower-cased characters are more common in entities, so we check for them first.
		if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
			continue
		}
		if c != ';' {
			i--
		}
		break
	}

	name := s[1:i]
	if len(name) == 0 {
		// No-op.
	} else if attribute && name[len(name)-1] != ';' && len(s) > i && s[i] == '=' {
		// No-op.
	} else if v, ok := lookupEntity(name); ok {
		return v, 0, i
	} else if !attribute {
		maxLen := min(len(name)-1, maxEntityWithoutSemicolon)
		for j := maxLen; j > 1; j-- {
			if v, ok := lookupEntity(name[:j]); ok {
				return v, 0, j + 1
			}
		}
	}

	return "", '&', 1
}

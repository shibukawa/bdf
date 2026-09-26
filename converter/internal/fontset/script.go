package fontset

import (
	"strings"
	"unicode"
)

// Script returns the East Asian script a language tag stands for ("Jpan",
// "Hans", "Hant" or "Hang"), or "" for other languages. Documents pick
// their East Asian fonts by it.
func Script(lang string) string {
	l := strings.ToLower(lang)
	switch {
	case strings.HasPrefix(l, "ja"):
		return "Jpan"
	case strings.HasPrefix(l, "zh-tw"), strings.HasPrefix(l, "zh-hk"), strings.HasPrefix(l, "zh-mo"), strings.HasPrefix(l, "zh-hant"):
		return "Hant"
	case strings.HasPrefix(l, "zh"):
		return "Hans"
	case strings.HasPrefix(l, "ko"):
		return "Hang"
	}
	return ""
}

// RuneLang returns the language a character implies: "ja" for kana and
// "ko" for hangul; "" for others (Han characters alone could be Chinese
// as well as Japanese or Korean).
func RuneLang(r rune) string {
	switch {
	case unicode.In(r, unicode.Hiragana, unicode.Katakana):
		return "ja"
	case unicode.Is(unicode.Hangul, r):
		return "ko"
	}
	return ""
}

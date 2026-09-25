package fontdb

import (
	"strings"
	"unicode"
)

// Generic families.
const (
	Sans  = "sans-serif"
	Serif = "serif"
	Mono  = "monospace"
)

// substitutes lists metric-compatible (or at least look-alike) families to
// try when a requested family is not installed, by normalized name.
var substitutes = map[string][]string{
	"calibri":           {"Carlito"},
	"calibrilight":      {"Carlito"},
	"cambria":           {"Caladea"},
	"cambriamath":       {"Caladea"},
	"arial":             {"Liberation Sans", "Arimo", "Helvetica", "Helvetica Neue"},
	"helvetica":         {"Liberation Sans", "Arimo", "Arial", "Helvetica Neue"},
	"helveticaneue":     {"Liberation Sans", "Arimo", "Helvetica", "Arial"},
	"arialnarrow":       {"Liberation Sans Narrow"},
	"timesnewroman":     {"Liberation Serif", "Tinos", "Times"},
	"times":             {"Liberation Serif", "Tinos", "Times New Roman"},
	"couriernew":        {"Liberation Mono", "Cousine", "Courier"},
	"courier":           {"Liberation Mono", "Cousine", "Courier New"},
	"georgia":           {"Gelasio"},
	"symbol":            {"OpenSymbol", "Standard Symbols PS"},
	"wingdings":         {"OpenSymbol"},
	"segoeuisymbol":     {"Noto Sans Symbols", "Noto Sans Symbols2", "DejaVu Sans"},
	"segoeuiemoji":      {"Noto Color Emoji"},
	"meiryo":            jaGothic,
	"メイリオ":              jaGothic,
	"meiryoui":          jaGothic,
	"yugothic":          jaGothic,
	"游ゴシック":             jaGothic,
	"yugothicui":        jaGothic,
	"yugothicmedium":    jaGothic,
	"游ゴシックmedium":       jaGothic,
	"游ゴシックlight":        jaGothic,
	"mspgothic":         jaGothic,
	"mspゴシック":           jaGothic,
	"msuigothic":        jaGothic,
	"hiraginosans":      jaGothic,
	"ヒラギノ角ゴpron":        jaGothic,
	"ヒラギノ角ゴシック":         jaGothic,
	"msgothic":          jaGothicMono,
	"msゴシック":            jaGothicMono,
	"yumincho":          jaMincho,
	"游明朝":               jaMincho,
	"mspmincho":         jaMincho,
	"msp明朝":             jaMincho,
	"msmincho":          jaMinchoMono,
	"ms明朝":              jaMinchoMono,
	"hiraginomincho":    jaMincho,
	"ヒラギノ明朝pron":        jaMincho,
	"microsoftyahei":    zhHans,
	"微软雅黑":              zhHans,
	"simsun":            zhHansSerif,
	"宋体":                zhHansSerif,
	"simhei":            zhHans,
	"黑体":                zhHans,
	"dengxian":          zhHans,
	"等线":                zhHans,
	"microsoftjhenghei": zhHant,
	"微軟正黑體":             zhHant,
	"pmingliu":          zhHantSerif,
	"新細明體":              zhHantSerif,
	"mingliu":           zhHantSerif,
	"malgungothic":      koGothic,
	"맑은고딕":              koGothic,
	"gulim":             koGothic,
	"굴림":                koGothic,
	"dotum":             koGothic,
	"돋움":                koGothic,
	"batang":            koSerif,
	"바탕":                koSerif,
}

var (
	jaGothic     = []string{"Noto Sans CJK JP", "Noto Sans JP", "Source Han Sans JP", "IPAexGothic", "IPAPGothic", "IPAGothic", "TakaoPGothic", "VL PGothic", "M PLUS 1p", "Hiragino Sans", "Yu Gothic", "Meiryo", "WenQuanYi Zen Hei"}
	jaGothicMono = []string{"IPAGothic", "TakaoGothic", "VL Gothic", "Noto Sans Mono CJK JP", "MS Gothic", "Noto Sans CJK JP", "IPAexGothic", "WenQuanYi Zen Hei Mono"}
	jaMincho     = []string{"Noto Serif CJK JP", "Noto Serif JP", "Source Han Serif JP", "IPAexMincho", "IPAPMincho", "IPAMincho", "TakaoPMincho", "Hiragino Mincho ProN", "Yu Mincho"}
	jaMinchoMono = []string{"IPAMincho", "TakaoMincho", "MS Mincho", "Noto Serif CJK JP", "IPAexMincho"}
	zhHans       = []string{"Noto Sans CJK SC", "Noto Sans SC", "Source Han Sans SC", "WenQuanYi Zen Hei", "WenQuanYi Micro Hei", "Microsoft YaHei", "PingFang SC"}
	zhHansSerif  = []string{"Noto Serif CJK SC", "Noto Serif SC", "Source Han Serif SC", "AR PL UMing CN", "SimSun"}
	zhHant       = []string{"Noto Sans CJK TC", "Noto Sans TC", "Source Han Sans TC", "WenQuanYi Zen Hei", "Microsoft JhengHei", "PingFang TC"}
	zhHantSerif  = []string{"Noto Serif CJK TC", "Noto Serif TC", "AR PL UMing TW", "PMingLiU"}
	koGothic     = []string{"Noto Sans CJK KR", "Noto Sans KR", "Source Han Sans KR", "NanumGothic", "UnDotum", "Malgun Gothic"}
	koSerif      = []string{"Noto Serif CJK KR", "Noto Serif KR", "NanumMyeongjo", "UnBatang", "Batang"}

	genericLatin = map[string][]string{
		Sans:  {"Liberation Sans", "Arimo", "Arial", "Helvetica", "Carlito", "DejaVu Sans", "FreeSans", "Noto Sans", "Open Sans", "Roboto"},
		Serif: {"Liberation Serif", "Tinos", "Times New Roman", "Times", "Caladea", "DejaVu Serif", "FreeSerif", "Noto Serif"},
		Mono:  {"Liberation Mono", "Cousine", "Courier New", "DejaVu Sans Mono", "FreeMono", "Noto Sans Mono"},
	}
	// genericSymbols have wide coverage of symbols, arrows and dingbats.
	genericSymbols = []string{"DejaVu Sans", "Noto Sans Symbols", "Noto Sans Symbols2", "Segoe UI Symbol", "OpenSymbol", "FreeSerif", "FreeSans", "Arial Unicode MS", "Noto Color Emoji", "Unifont"}
)

// Classify guesses the generic family of a family name.
func Classify(name string) string {
	n := Normalize(name)
	switch {
	case strings.Contains(n, "mono") || strings.Contains(n, "courier") || strings.Contains(n, "consola") || strings.Contains(n, "code") ||
		n == "msgothic" || n == "msゴシック" || n == "msmincho" || n == "ms明朝" || strings.Contains(n, "lucidaconsole"):
		return Mono
	case strings.Contains(n, "sans") || strings.Contains(n, "gothic") || strings.Contains(n, "ゴシック") || strings.Contains(n, "黑") || strings.Contains(n, "hei"):
		return Sans
	case strings.Contains(n, "serif") || strings.Contains(n, "times") || strings.Contains(n, "mincho") || strings.Contains(n, "明朝") ||
		strings.Contains(n, "cambria") || strings.Contains(n, "georgia") || strings.Contains(n, "garamond") || strings.Contains(n, "century") ||
		strings.Contains(n, "palatino") || strings.Contains(n, "bodoni") || strings.Contains(n, "bookantiqua") || strings.Contains(n, "roman") ||
		strings.Contains(n, "song") || strings.Contains(n, "宋") || strings.Contains(n, "ming") || strings.Contains(n, "明") || strings.Contains(n, "batang") || strings.Contains(n, "바탕"):
		return Serif
	}
	return Sans
}

// Resolved is a family request mapped to an available face.
type Resolved struct {
	Face        *Face // nil when no font is available at all
	SynthBold   bool
	SynthItalic bool
	Exact       bool   // the requested family itself was found
	Generic     string // generic family of the request
}

// Resolve maps a family request to a face: the family itself, a known
// substitute, a generic family of the same kind (serif, sans-serif,
// monospace; CJK chains when cjk is set), and finally any face at all.
func (db *DB) Resolve(family string, bold, italic, cjk bool) Resolved {
	res := Resolved{Generic: Classify(family)}
	try := func(names ...string) bool {
		for _, n := range names {
			if faces := db.Family(n); len(faces) > 0 {
				res.Face, res.SynthBold, res.SynthItalic = Match(faces, bold, italic)
				return true
			}
		}
		return false
	}
	if family != "" && try(family) {
		res.Exact = true
		return res
	}
	if try(substitutes[Normalize(family)]...) {
		return res
	}
	if cjk {
		chain := jaGothic
		if res.Generic == Serif {
			chain = jaMincho
		}
		if try(chain...) || try(zhHans...) || try(koGothic...) {
			return res
		}
	}
	if try(genericLatin[res.Generic]...) || try(genericLatin[Sans]...) {
		return res
	}
	if !cjk && (try(jaGothic...) || try(zhHans...)) {
		return res
	}
	if len(db.Faces) > 0 {
		res.Face, res.SynthBold, res.SynthItalic = Match(db.firstFamily(), bold, italic)
	}
	return res
}

// firstFamily returns the faces of the alphabetically first family, the
// deterministic last resort.
func (db *DB) firstFamily() []*Face {
	var best string
	for _, f := range db.Faces {
		if best == "" || f.Family < best {
			best = f.Family
		}
	}
	var out []*Face
	for _, f := range db.Faces {
		if f.Family == best {
			out = append(out, f)
		}
	}
	return out
}

// Fallbacks returns faces to try, in order, for characters the resolved
// face lacks: CJK chains, the generic family, symbol fonts, then every
// other face so that a small private font directory is used completely.
func (db *DB) Fallbacks(generic string, bold, italic bool) []Resolved {
	var out []Resolved
	seen := map[*Face]bool{}
	add := func(names []string) {
		for _, n := range names {
			faces := db.Family(n)
			if len(faces) == 0 {
				continue
			}
			f, sb, si := Match(faces, bold, italic)
			if !seen[f] {
				seen[f] = true
				out = append(out, Resolved{Face: f, SynthBold: sb, SynthItalic: si, Generic: generic})
			}
		}
	}
	add(genericLatin[generic])
	if generic == Serif {
		add(jaMincho)
	}
	add(jaGothic)
	add(zhHans)
	add(zhHant)
	add(koGothic)
	add(genericLatin[Sans])
	add(genericSymbols)
	if len(db.Faces) <= 64 {
		fams := map[string]bool{}
		var names []string
		for _, f := range db.Faces {
			if !fams[f.Family] {
				fams[f.Family] = true
				names = append(names, f.Family)
			}
		}
		add(names)
	}
	return out
}

// IsCJK reports whether r belongs to the East Asian scripts that Office
// lays out with the "ea" font: Han, kana, Hangul, Bopomofo, CJK symbols
// and punctuation, and the full-width forms.
func IsCJK(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x11FF, // Hangul Jamo
		r >= 0x2E80 && r <= 0x2FDF, // radicals
		r >= 0x2FF0 && r <= 0x303F, // ideographic description, CJK symbols and punctuation
		r >= 0x3040 && r <= 0x31FF, // kana, bopomofo, Hangul compatibility, kanbun, katakana ext
		r >= 0x3200 && r <= 0x33FF, // enclosed CJK, compatibility
		r >= 0x3400 && r <= 0x4DBF,
		r >= 0x4E00 && r <= 0x9FFF,
		r >= 0xA960 && r <= 0xA97F,
		r >= 0xAC00 && r <= 0xD7FF,
		r >= 0xF900 && r <= 0xFAFF,
		r >= 0xFE30 && r <= 0xFE4F,
		r >= 0xFF00 && r <= 0xFFEF,
		r >= 0x20000 && r <= 0x3FFFF:
		return true
	}
	return unicode.Is(unicode.Han, r)
}

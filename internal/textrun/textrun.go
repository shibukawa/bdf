// Package textrun holds what the renderers of BDF text share: how a string
// of a FILL_TEXT is put in visual order, where its baseline lies, which
// families a font reference names, and the bold and italic that are
// synthesized for a face that has none. The renderers (raster/imagebdf,
// vector/pdfbdf) draw text the same way with it.
package textrun

import (
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/fontdb"
	"golang.org/x/text/unicode/bidi"
)

// GenericFamilies maps CSS generic family keywords to the generic families
// of fontdb.
var GenericFamilies = map[string]string{
	"sans-serif": fontdb.Sans, "serif": fontdb.Serif, "monospace": fontdb.Mono,
	"cursive": fontdb.Sans, "fantasy": fontdb.Sans, "system-ui": fontdb.Sans, "math": fontdb.Serif,
	"ui-sans-serif": fontdb.Sans, "ui-serif": fontdb.Serif, "ui-monospace": fontdb.Mono, "ui-rounded": fontdb.Sans,
}

// SplitFamilies splits a CSS font-family list into names without quotes.
func SplitFamilies(list string) []string {
	var out []string
	for _, f := range strings.Split(list, ",") {
		f = strings.TrimSpace(f)
		if len(f) >= 2 && (f[0] == '"' || f[0] == '\'') && f[len(f)-1] == f[0] {
			f = strings.ReplaceAll(f[1:len(f)-1], `\`+string(f[0]), string(f[0]))
		}
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}

// VisualOrder reorders a string for display: runs of right-to-left
// characters are reversed, and so is the order of the runs in a
// right-to-left paragraph. It covers what previews need, not the whole
// Unicode bidirectional algorithm (no mirroring, no Arabic shaping).
func VisualOrder(s string, rtlBase bool) []rune {
	runes := []rune(s)
	if !rtlBase && !HasRTL(s) {
		return runes
	}
	const (
		ltr = iota
		rtl
		neutral
	)
	dirs := make([]int, len(runes))
	for i, r := range runes {
		p, _ := bidi.LookupRune(r)
		switch p.Class() {
		case bidi.R, bidi.AL:
			dirs[i] = rtl
		case bidi.L, bidi.EN, bidi.AN:
			dirs[i] = ltr
		default:
			dirs[i] = neutral
		}
	}
	base := ltr
	if rtlBase {
		base = rtl
	}
	// neutrals take the direction around them when both sides agree
	for i := 0; i < len(dirs); {
		if dirs[i] != neutral {
			i++
			continue
		}
		j := i
		for j < len(dirs) && dirs[j] == neutral {
			j++
		}
		before, after := base, base
		if i > 0 {
			before = dirs[i-1]
		}
		if j < len(dirs) {
			after = dirs[j]
		}
		d := base
		if before == after {
			d = before
		}
		for k := i; k < j; k++ {
			dirs[k] = d
		}
		i = j
	}
	type run struct{ from, to, dir int }
	var runs []run
	for i := 0; i < len(runes); {
		j := i
		for j < len(runes) && dirs[j] == dirs[i] {
			j++
		}
		runs = append(runs, run{i, j, dirs[i]})
		i = j
	}
	if base == rtl {
		for i, j := 0, len(runs)-1; i < j; i, j = i+1, j-1 {
			runs[i], runs[j] = runs[j], runs[i]
		}
	}
	out := make([]rune, 0, len(runes))
	for _, r := range runs {
		if r.dir == rtl {
			for k := r.to - 1; k >= r.from; k-- {
				out = append(out, runes[k])
			}
		} else {
			out = append(out, runes[r.from:r.to]...)
		}
	}
	return out
}

// HasRTL reports whether a string has right-to-left characters.
func HasRTL(s string) bool {
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		i += n
		if r < 0x590 {
			continue
		}
		p, _ := bidi.LookupRune(r)
		if c := p.Class(); c == bidi.R || c == bidi.AL {
			return true
		}
	}
	return false
}

// BaselineShift is how far below the anchor the alphabetic baseline lies
// for a textBaseline, from the em box of a face with the ascent asc and
// the descent dsc in em, both positive (as Blink does).
func BaselineShift(baseline byte, size, asc, dsc float64) float64 {
	emAsc := size * asc / (asc + dsc)
	emDsc := size - emAsc
	switch baseline {
	case bdf.BaselineTop:
		return emAsc
	case bdf.BaselineHanging:
		return asc * size * 0.8
	case bdf.BaselineMiddle:
		return (emAsc - emDsc) / 2
	case bdf.BaselineIdeographic:
		return -dsc * size
	case bdf.BaselineBottom:
		return -emDsc
	}
	return 0
}

// FakeBoldWidth is the outline growth of synthesized bold for a size in
// device pixels (Skia's interpolation between 1/24 at 9 px and 1/32 at 36 px).
func FakeBoldWidth(size float64) float64 {
	switch {
	case size <= 9:
		return size / 24
	case size >= 36:
		return size / 32
	}
	k := 1.0/24 + (size-9)/(36-9)*(1.0/32-1.0/24)
	return size * k
}

// FakeItalicSkew is the slant of synthesized italic (x per unit of height).
const FakeItalicSkew = 0.25

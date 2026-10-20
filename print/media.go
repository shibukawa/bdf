package print

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Media is a size of paper by its PWG name (PWG 5101.1), such as
// "iso_a4_210x297mm": a class, a name, and the size, which the name
// carries itself, so that any name a printer reports can be read.
type Media struct {
	Name string
	// Width and Height are in points (1/72 inch).
	Width, Height float64
}

// mm is a point in millimetres.
const mm = 72 / 25.4

// ParseMedia reads the size of a PWG self-describing media name: its
// last part, "<w>x<h>mm" or "<w>x<h>in".
func ParseMedia(name string) (Media, error) {
	i := strings.LastIndexByte(name, '_')
	if i < 0 {
		return Media{}, fmt.Errorf("print: %q is not a PWG media name", name)
	}
	dims := name[i+1:]
	var unit float64
	switch {
	case strings.HasSuffix(dims, "mm"):
		unit, dims = mm, dims[:len(dims)-2]
	case strings.HasSuffix(dims, "in"):
		unit, dims = 72, dims[:len(dims)-2]
	default:
		return Media{}, fmt.Errorf("print: %q has no size in mm or in", name)
	}
	ws, hs, ok := strings.Cut(dims, "x")
	w, err1 := strconv.ParseFloat(ws, 64)
	h, err2 := strconv.ParseFloat(hs, 64)
	if !ok || err1 != nil || err2 != nil || !(w > 0) || !(h > 0) || w > 10000 || h > 10000 {
		return Media{}, fmt.Errorf("print: %q has no size in mm or in", name)
	}
	return Media{Name: name, Width: w * unit, Height: h * unit}, nil
}

// aliases are the short names that -media takes for the sizes people
// name: B4 and B5 are the JIS sizes, as in Japan.
var aliases = map[string]string{
	"a0": "iso_a0_841x1189mm", "a1": "iso_a1_594x841mm", "a2": "iso_a2_420x594mm", "a3": "iso_a3_297x420mm",
	"a4": "iso_a4_210x297mm", "a5": "iso_a5_148x210mm", "a6": "iso_a6_105x148mm",
	"b4": "jis_b4_257x364mm", "b5": "jis_b5_182x257mm", "b6": "jis_b6_128x182mm",
	"jis-b4": "jis_b4_257x364mm", "jis-b5": "jis_b5_182x257mm", "iso-b4": "iso_b4_250x353mm", "iso-b5": "iso_b5_176x250mm",
	"letter": "na_letter_8.5x11in", "legal": "na_legal_8.5x14in", "ledger": "na_ledger_11x17in", "tabloid": "na_ledger_11x17in",
	"executive": "na_executive_7.25x10.5in",
	"hagaki":    "jpn_hagaki_100x148mm", "postcard": "jpn_hagaki_100x148mm", "oufuku": "jpn_oufuku_148x200mm",
	"chou3": "jpn_chou3_120x235mm", "chou4": "jpn_chou4_90x205mm", "kaku2": "jpn_kaku2_240x332mm", "you4": "jpn_you4_105x235mm",
	"dl": "iso_dl_110x220mm", "c4": "iso_c4_229x324mm", "c5": "iso_c5_162x229mm", "c6": "iso_c6_114x162mm",
	"l": "oe_photo-l_3.5x5in", "2l": "na_5x7_5x7in", "4x6": "na_index-4x6_4x6in",
}

// ResolveMedia turns a -media value into a PWG name: a short name
// ("a4", "b5", "hagaki", "letter"), a size alone ("104x150mm", "4x6in",
// a custom_ name), or a PWG name as it is.
func ResolveMedia(s string) (Media, error) {
	if n, ok := aliases[strings.ToLower(s)]; ok {
		s = n
	}
	if !strings.Contains(s, "_") {
		// a size alone, as a label's: 104x150mm, 4x6in
		if m, err := ParseMedia("custom_" + s); err == nil {
			return m, nil
		}
	}
	return ParseMedia(s)
}

// tolerance is how far a page may be from a size of paper and still be
// on it: 2 mm, the rounding of sizes given in inches or points.
const tolerance = 2 * mm

// Matches reports whether a page of w × h points is the size of the
// media, upright or turned.
func (m Media) Matches(w, h float64) bool {
	near := func(a, b float64) bool { return math.Abs(a-b) <= tolerance }
	return near(m.Width, w) && near(m.Height, h) || near(m.Width, h) && near(m.Height, w)
}

// ChooseMedia picks the media of a page of w × h points among a
// printer's: the first of ready (media-ready) that matches the page,
// else the first of supported (media-supported) that does; "" when none
// does. Names that are not self-describing are passed over.
func ChooseMedia(w, h float64, ready, supported []string) string {
	for _, list := range [][]string{ready, supported} {
		for _, n := range list {
			if m, err := ParseMedia(n); err == nil && m.Matches(w, h) {
				return n
			}
		}
	}
	return ""
}

package font

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/internal/sfnt"
)

// section is a table of the overview: rows of a label and a value.
type section struct {
	title string
	rows  []row
}

type row struct {
	label, value string
	lang         string // the language of the value, when it is a localized name
}

func (s *section) add(label, value string) {
	if strings.TrimSpace(value) != "" {
		s.rows = append(s.rows, row{label: label, value: value})
	}
}

// be16, be32 and the signed forms read a table, 0 past its end.
func be16(b []byte, i int) int { return int(sfnt.BE16(b, i)) }
func be32(b []byte, i int) int { return int(sfnt.BE32(b, i)) }
func s16(b []byte, i int) int  { return int(int16(sfnt.BE16(b, i))) }

// fixed reads a 16.16 fixed-point number.
func fixed(b []byte, i int) float64 { return float64(int32(sfnt.BE32(b, i))) / 65536 }

// num formats a number with thousands separators.
func num(n int) string {
	s := strconv.Itoa(max(n, -n))
	var b strings.Builder
	if n < 0 {
		b.WriteByte('-')
	}
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	return b.String()
}

// decimal formats v with up to 3 decimals.
func decimal(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'f', -1, 64)
}

func bytesize(n int) string {
	switch {
	case n >= 1<<20:
		return decimal(math.Round(float64(n)/(1<<20)*10)/10) + " MB"
	case n >= 1<<10:
		return decimal(math.Round(float64(n)/(1<<10)*10)/10) + " KB"
	}
	return num(n) + " bytes"
}

// outlineFormat names the kind of outlines of a font.
func (fc *face) outlineFormat() string {
	t := fc.f.Tables
	switch {
	case t["CFF2"] != nil:
		return "OpenType with CFF2 outlines"
	case t["CFF "] != nil:
		return "OpenType with CFF outlines"
	case t["glyf"] != nil:
		return "TrueType outlines"
	case t["sbix"] != nil || t["CBDT"] != nil || t["EBDT"] != nil:
		return "Bitmaps only"
	}
	return "No outlines"
}

// weightNames names the OS/2 weight classes.
var weightNames = map[int]string{100: "Thin", 200: "Extra Light", 300: "Light", 400: "Regular", 500: "Medium",
	600: "Semi Bold", 700: "Bold", 800: "Extra Bold", 900: "Black"}

var widthNames = []string{"", "Ultra-condensed", "Extra-condensed", "Condensed", "Semi-condensed", "Normal",
	"Semi-expanded", "Expanded", "Extra-expanded", "Ultra-expanded"}

// embedding describes the OS/2 fsType embedding permissions.
func embedding(fs int) string {
	var parts []string
	switch fs & 0xF {
	case 0:
		parts = append(parts, "Installable")
	case 2:
		parts = append(parts, "Restricted (may not be embedded)")
	case 4:
		parts = append(parts, "Preview & print")
	case 8:
		parts = append(parts, "Editable")
	default:
		parts = append(parts, "0x"+strconv.FormatInt(int64(fs), 16))
	}
	if fs&0x100 != 0 {
		parts = append(parts, "no subsetting")
	}
	if fs&0x200 != 0 {
		parts = append(parts, "bitmap embedding only")
	}
	return strings.Join(parts, ", ")
}

// headTime reads a head date: seconds since 1904-01-01; ok is false for
// none (0) or one past 2104.
func headTime(b []byte, i int) (t time.Time, ok bool) {
	hi, lo := int64(be32(b, i)), int64(be32(b, i+4))
	secs := hi<<32 | lo
	if secs <= 0 || secs > 200*365*86400 {
		return t, false
	}
	return time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(secs) * time.Second), true
}

// longDateTime formats a head date.
func longDateTime(b []byte, i int) string {
	if t, ok := headTime(b, i); ok {
		return t.Format("2006-01-02 15:04:05 UTC")
	}
	return ""
}

// sections returns the tables of the overview.
func (fc *face) sections(fl *file, outlines bool) []section {
	f, t := fc.f, fc.f.Tables
	names := section{title: "Names"}
	// the family and full names in each language the font names itself in
	addLocal := func(label string, ids ...uint16) {
		for _, id := range ids {
			if v := fc.name(id); v != "" {
				names.add(label, v)
				for _, r := range fc.localNames(id) {
					names.rows = append(names.rows, row{label: label + " (" + r.lang + ")", value: r.text, lang: r.lang})
				}
				return
			}
		}
	}
	addName := func(label string, ids ...uint16) {
		for _, id := range ids {
			if v := fc.name(id); v != "" {
				names.add(label, v)
				return
			}
		}
	}
	addLocal("Family", 16, 1)
	addName("Style", 17, 2)
	addLocal("Full name", 4)
	addName("PostScript name", 6)
	addName("Version", 5)
	addName("Unique ID", 3)
	addName("Designer", 9)
	addName("Designer URL", 12)
	addName("Manufacturer", 8)
	addName("Vendor URL", 11)
	if os2 := t["OS/2"]; len(os2) >= 62 {
		names.add("Vendor ID", strings.TrimRight(string(os2[58:62]), " \x00"))
	}
	addName("Description", 10)
	addName("Sample text", 19)
	addName("WWS family", 21)
	addName("WWS style", 22)
	addName("Variations PostScript prefix", 25)

	lic := section{title: "License"}
	lic.add("Copyright", fc.name(0))
	lic.add("Trademark", fc.name(7))
	lic.add("License", fc.name(13))
	lic.add("License URL", fc.name(14))
	if f.HasFSType {
		e := embedding(int(f.FSType))
		if outlines {
			e += " — drawn as outlines here, the font is not embedded"
		}
		lic.add("Embedding", e)
	}

	info := section{title: "Font"}
	format := fc.outlineFormat()
	switch fl.container {
	case "WOFF":
		format += ", in a WOFF web font"
	case "WOFF2":
		format += ", in a WOFF2 web font"
	}
	info.add("Format", format)
	if fl.collection {
		info.add("Collection", fmt.Sprintf("font %d of %d", fc.index+1, len(fl.faces)))
	}
	info.add("File size", bytesize(fl.size))
	info.add("Glyphs", num(f.NumGlyphs))
	info.add("Characters", num(len(fc.runes)))
	if n, more := variationSequences(t["cmap"]); n > 0 {
		info.add("Variation sequences", moreThan(more)+num(n))
	}
	info.add("Units per em", num(f.UnitsPerEm))
	if head := t["head"]; len(head) >= 54 {
		info.add("Revision", decimal(fixed(head, 4)))
		info.add("Created", longDateTime(head, 20))
		info.add("Modified", longDateTime(head, 28))
	}
	if os2 := t["OS/2"]; len(os2) >= 8 {
		w := be16(os2, 4)
		ws := strconv.Itoa(w)
		if n := weightNames[w]; n != "" {
			ws += " (" + n + ")"
		}
		info.add("Weight class", ws)
		if wd := be16(os2, 6); wd >= 1 && wd <= 9 {
			info.add("Width class", strconv.Itoa(wd)+" ("+widthNames[wd]+")")
		}
		if len(os2) >= 64 {
			info.add("Style flags", fsSelection(be16(os2, 62)))
		}
		if len(os2) >= 42 {
			var p []string
			for i := range 10 {
				p = append(p, strconv.Itoa(int(os2[32+i])))
			}
			if strings.Join(p, "") != "0000000000" {
				info.add("PANOSE", strings.Join(p, " "))
			}
		}
		if be16(os2, 0) >= 5 && len(os2) >= 100 {
			lo, hi := float64(be16(os2, 96))/20, float64(be16(os2, 98))/20
			info.add("Optical sizes", decimal(lo)+"–"+decimal(hi)+" pt")
		}
	}
	if post := t["post"]; len(post) >= 16 {
		if a := fixed(post, 4); a != 0 {
			info.add("Italic angle", decimal(a)+"°")
		}
		if be32(post, 12) != 0 {
			info.add("Monospaced", "yes")
		}
	}
	info.add("Hinting", hinting(t))
	info.add("Vertical metrics", present(t, "vhea", "vmtx", "VORG"))
	info.add("Color", colorFormats(t))
	if s := strikes(t["EBLC"], 48); s != "" {
		info.add("Embedded bitmaps", s+" ppem")
	}
	if t["MATH"] != nil {
		info.add("Math", "MATH table (formula layout)")
	}
	info.add("Kerning", fc.kerning())
	var other []string
	if t["morx"] != nil || t["mort"] != nil {
		other = append(other, "AAT (morx)")
	}
	if t["Silf"] != nil {
		other = append(other, "Graphite (Silf)")
	}
	info.add("Other layout", strings.Join(other, ", "))
	design, supported := metaLangs(t["meta"])
	info.add("Design languages", design)
	info.add("Supported languages", supported)

	metrics := section{title: "Metrics (font units)"}
	metrics.add("Ascender, descender, line gap (hhea)", fmt.Sprintf("%d, %d, %d", f.Ascent, f.Descent, f.LineGap))
	if os2 := t["OS/2"]; len(os2) >= 78 {
		metrics.add("Typographic ascender, descender, line gap (OS/2)", fmt.Sprintf("%d, %d, %d", s16(os2, 68), s16(os2, 70), s16(os2, 72)))
		metrics.add("Windows ascent, descent (OS/2)", fmt.Sprintf("%d, %d", be16(os2, 74), be16(os2, 76)))
		if be16(os2, 0) >= 2 && len(os2) >= 90 {
			metrics.add("x-height, cap height", fmt.Sprintf("%d, %d", s16(os2, 86), s16(os2, 88)))
		}
	}
	if head := t["head"]; len(head) >= 54 {
		metrics.add("Bounding box", fmt.Sprintf("%d, %d – %d, %d", s16(head, 36), s16(head, 38), s16(head, 40), s16(head, 42)))
	}
	if post := t["post"]; len(post) >= 12 {
		metrics.add("Underline position, thickness", fmt.Sprintf("%d, %d", s16(post, 8), s16(post, 10)))
	}

	out := []section{names, lic, info, metrics}
	if v := fc.variations(); len(v.rows) > 0 {
		out = append(out, v)
	}
	return out
}

// tableList returns the tables of the font and their sizes, by tag.
func (fc *face) tableList() []row {
	t := fc.f.Tables
	tags := make([]string, 0, len(t))
	for tag := range t {
		tags = append(tags, tag)
	}
	slices.Sort(tags)
	var out []row
	for _, tag := range tags {
		out = append(out, row{label: tag, value: bytesize(len(t[tag]))})
	}
	return out
}

func fsSelection(fs int) string {
	var s []string
	for bit, name := range []string{"Italic", "Underscore", "Negative", "Outlined", "Strikeout", "Bold", "Regular", "Use typographic metrics", "WWS", "Oblique"} {
		if fs&(1<<bit) != 0 {
			s = append(s, name)
		}
	}
	return strings.Join(s, ", ")
}

func present(t map[string][]byte, tags ...string) string {
	var s []string
	for _, tag := range tags {
		if t[tag] != nil {
			s = append(s, tag)
		}
	}
	return strings.Join(s, ", ")
}

func hinting(t map[string][]byte) string {
	var s []string
	if t["fpgm"] != nil || t["prep"] != nil {
		s = append(s, "TrueType instructions")
	}
	if t["gasp"] != nil {
		s = append(s, "gasp")
	}
	if t["hdmx"] != nil || t["VDMX"] != nil || t["LTSH"] != nil {
		s = append(s, "device metrics")
	}
	return strings.Join(s, ", ")
}

// colorFormats lists the color glyph tables of a font.
func colorFormats(t map[string][]byte) string {
	var s []string
	if c := t["COLR"]; len(c) >= 14 {
		n := be16(c, 2)
		if be16(c, 0) >= 1 && len(c) >= 34 {
			if bl := be32(c, 14); bl > 0 && bl+4 <= len(c) {
				n = max(n, be32(c, bl))
			}
			s = append(s, "COLRv1 ("+num(n)+" glyphs)")
		} else {
			s = append(s, "COLRv0 ("+num(n)+" glyphs)")
		}
	}
	if c := t["CPAL"]; len(c) >= 12 {
		s = append(s, fmt.Sprintf("CPAL (%d palettes of %d colors)", be16(c, 4), be16(c, 2)))
	}
	if b := t["sbix"]; len(b) >= 8 {
		var sizes []string
		n := min(be32(b, 4), (len(b)-8)/4)
		for i := range n {
			if o := be32(b, 8+i*4); o+2 <= len(b) {
				sizes = append(sizes, strconv.Itoa(be16(b, o)))
			}
		}
		s = append(s, "sbix bitmaps ("+strings.Join(sizes, ", ")+" ppem)")
	}
	if sz := strikes(t["CBLC"], 48); sz != "" {
		s = append(s, "CBDT bitmaps ("+sz+" ppem)")
	}
	if b := t["SVG "]; len(b) >= 10 {
		if o := be32(b, 2); o+2 <= len(b) {
			s = append(s, "SVG ("+num(be16(b, o))+" documents)")
		}
	}
	return strings.Join(s, ", ")
}

// strikes lists the ppem sizes of an EBLC or CBLC table.
func strikes(b []byte, size int) string {
	if len(b) < 8 {
		return ""
	}
	n := min(be32(b, 4), (len(b)-8)/size)
	var s []string
	for i := range n {
		s = append(s, strconv.Itoa(int(b[8+i*size+44])))
	}
	return strings.Join(s, ", ")
}

// maxRanges bounds the ranges of default variation sequences that are
// counted: the variation selectors of a subtable may share their ranges, or
// have ranges that overlap.
const maxRanges = 1 << 22

// variationSequences counts the variation sequences of a cmap's format 14
// subtable; more is true when it has more ranges than are counted.
func variationSequences(cm []byte) (count int, more bool) {
	n := be16(cm, 2)
	for i := range n {
		rec := 4 + i*8
		if be16(cm, rec) != 0 || be16(cm, rec+2) != 5 {
			continue
		}
		st := cm[min(be32(cm, rec+4), len(cm)):]
		if be16(st, 0) != 14 {
			continue
		}
		left := maxRanges
		defaults := map[int]int{} // the sequences of a table of ranges, by its offset
		m := min(be32(st, 6), max(len(st)-10, 0)/11)
		for j := range m {
			r := 10 + j*11
			if d := be32(st, r+3); d > 0 && d+4 <= len(st) {
				sum, ok := defaults[d]
				if !ok {
					ranges := min(be32(st, d), (len(st)-d-4)/4)
					if left -= ranges; left < 0 {
						return count, true
					}
					for k := range ranges {
						sum += int(st[d+4+k*4+3]) + 1
					}
					defaults[d] = sum
				}
				count += sum
			}
			if u := be32(st, r+7); u > 0 && u+4 <= len(st) {
				count += min(be32(st, u), 1<<20)
			}
		}
		return count, false
	}
	return 0, false
}

// moreThan starts a number that is a lower bound.
func moreThan(more bool) string {
	if more {
		return "more than "
	}
	return ""
}

// metaLangs reads the design and supported languages of the meta table.
func metaLangs(b []byte) (design, supported string) {
	if len(b) < 16 || be32(b, 0) != 1 {
		return "", ""
	}
	n := min(be32(b, 12), (len(b)-16)/12)
	for i := range n {
		rec := 16 + i*12
		tag, off, l := string(b[rec:rec+4]), be32(b, rec+4), be32(b, rec+8)
		if off < 0 || l < 0 || off+l > len(b) {
			continue
		}
		v := strings.ReplaceAll(string(b[off:off+l]), ",", ", ")
		switch tag {
		case "dlng":
			design = v
		case "slng":
			supported = v
		}
	}
	return
}

// kerning describes the kerning of a font.
func (fc *face) kerning() string {
	var s []string
	if fc.gpos != nil {
		pairs := 0
		for _, f := range fc.gpos.Features {
			if f.Tag == "kern" {
				for _, l := range f.Lookups {
					pairs += fc.gpos.PairCount(l)
				}
			}
		}
		if pairs > 0 {
			s = append(s, "GPOS kern (about "+num(pairs)+" pairs)")
		}
	}
	if k := fc.f.Tables["kern"]; len(k) >= 18 && be16(k, 0) == 0 {
		s = append(s, "kern table ("+num(be16(k, 10))+" pairs)")
	}
	return strings.Join(s, ", ")
}

// axis is a variation axis of fvar.
type axis struct {
	tag           string
	name          string
	min, def, max float64
	hidden        bool
}

// instance is a named instance of fvar.
type instance struct {
	name   string
	ps     string
	coords []float64
}

// fvar reads the axes and named instances of a variable font.
func (fc *face) fvar() ([]axis, []instance) {
	b := fc.f.Tables["fvar"]
	if len(b) < 16 || be16(b, 0) != 1 {
		return nil, nil
	}
	off, n, size, ni, isize := be16(b, 4), be16(b, 8), be16(b, 10), be16(b, 12), be16(b, 14)
	if size < 20 || off+n*size > len(b) {
		return nil, nil
	}
	var axes []axis
	for i := range n {
		a := off + i*size
		ax := axis{tag: string(b[a : a+4]), min: fixed(b, a+4), def: fixed(b, a+8), max: fixed(b, a+12), hidden: be16(b, a+16)&1 != 0}
		ax.name = fc.name(uint16(be16(b, a+18)))
		axes = append(axes, ax)
	}
	var insts []instance
	p := off + n*size
	if isize < 4+4*n {
		return axes, nil
	}
	for i := range ni {
		r := p + i*isize
		if r+isize > len(b) {
			break
		}
		in := instance{name: fc.name(uint16(be16(b, r)))}
		for k := range n {
			in.coords = append(in.coords, fixed(b, r+4+k*4))
		}
		if isize >= 6+4*n {
			if id := be16(b, r+4+4*n); id != 0xFFFF {
				in.ps = fc.name(uint16(id))
			}
		}
		insts = append(insts, in)
	}
	return axes, insts
}

// variations is the section of a variable font's axes and instances.
func (fc *face) variations() section {
	s := section{title: "Variations"}
	axes, insts := fc.fvar()
	for _, a := range axes {
		v := fmt.Sprintf("%s to %s, default %s", decimal(a.min), decimal(a.max), decimal(a.def))
		if a.hidden {
			v += " (hidden)"
		}
		label := a.tag
		if a.name != "" {
			label += " " + a.name
		}
		s.add("Axis "+label, v)
	}
	for _, in := range insts {
		var c []string
		for k, v := range in.coords {
			if k < len(axes) {
				c = append(c, axes[k].tag+" "+decimal(v))
			}
		}
		name := in.name
		if name == "" {
			name = "(unnamed)"
		}
		v := strings.Join(c, ", ")
		if in.ps != "" {
			v += " — " + in.ps
		}
		s.add("Instance "+name, v)
	}
	if len(axes) > 0 {
		s.add("Drawn at", "the default instance")
	}
	return s
}

// palettes returns the colors of the CPAL palettes (0xRRGGBBAA).
func (fc *face) palettes() [][]uint32 {
	b := fc.f.Tables["CPAL"]
	if len(b) < 12 {
		return nil
	}
	entries, n, records, off := be16(b, 2), be16(b, 4), be16(b, 6), be32(b, 8)
	var out [][]uint32
	for i := range min(n, 64) {
		first := be16(b, 12+i*2)
		var pal []uint32
		for k := range min(entries, 256) {
			c := off + (first+k)*4
			if first+k >= records || c+4 > len(b) {
				break
			}
			// BGRA
			pal = append(pal, uint32(b[c+2])<<24|uint32(b[c+1])<<16|uint32(b[c])<<8|uint32(b[c+3]))
		}
		out = append(out, pal)
	}
	return out
}

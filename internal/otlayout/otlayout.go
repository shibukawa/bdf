// Package otlayout reads the OpenType layout tables GSUB, GPOS and GDEF:
// the scripts and language systems a font declares, its features and the
// lookups they apply, and what each lookup substitutes or positions. It
// lists what a font can do, as a font viewer shows it; it does not shape
// text.
//
// Fonts are untrusted input: offsets and counts are checked against the
// table, and what a table expands to (the glyphs of a coverage or class
// definition, the pairs of a class-based kerning lookup) is bounded.
// Records may share what they point at, so that a small table says it holds
// millions of features or subtables: a table is read into maxElements of
// them at most (ErrTooLarge), and its lookups list so many entries in all
// the calls made of them (maxSteps), then no more (Truncated).
package otlayout

import (
	"encoding/binary"
	"errors"
	"slices"
)

// Table is a GSUB or GPOS table.
type Table struct {
	Scripts  []Script
	Features []Feature
	Lookups  []Lookup
	// FeatureVariations is the number of feature variation records: the
	// regions of a variable font's design space where other lookups apply.
	FeatureVariations int

	gpos      bool
	pairCache map[int][]*pairSub
	// steps is how many entries the lookups may still list (see spend).
	steps int
}

// Script is a script the table has features for.
type Script struct {
	Tag string
	// Default is the language system used for languages the script does
	// not list; nil when there is none.
	Default *LangSys
	Langs   []LangSys
}

// LangSys is a language system: the features that apply to text of a
// script in a language.
type LangSys struct {
	Tag string // "" for a script's default
	// Required is the index of the feature that always applies, -1 when
	// there is none.
	Required int
	Features []int
}

// Feature is a feature of the table.
type Feature struct {
	Tag     string
	Lookups []int
	// UIName is the name ID of the name a stylistic set (ss01–ss20) or a
	// character variant (cv01–cv99) gives itself, 0 when it has none.
	UIName uint16
	// Chars are the characters a character variant applies to.
	Chars []rune
}

// Lookup is a lookup of the table.
type Lookup struct {
	// Type is the lookup type. An extension lookup has the type of the
	// subtables it points to.
	Type int
	Flag uint16
	// MarkSet is the index of the mark glyph set the lookup is restricted
	// to (see GDEF.MarkSets), -1 when there is none.
	MarkSet int

	subtables []data
}

// Lookup flags.
const (
	RightToLeft         = 0x0001
	IgnoreBaseGlyphs    = 0x0002
	IgnoreLigatures     = 0x0004
	IgnoreMarks         = 0x0008
	UseMarkFilteringSet = 0x0010
	MarkAttachmentType  = 0xFF00
)

// GSUB lookup types.
const (
	SubstSingle       = 1
	SubstMultiple     = 2
	SubstAlternate    = 3
	SubstLigature     = 4
	SubstContext      = 5
	SubstChainContext = 6
	substExtension    = 7
	SubstReverseChain = 8
)

// GPOS lookup types.
const (
	PosSingle       = 1
	PosPair         = 2
	PosCursive      = 3
	PosMarkToBase   = 4
	PosMarkToLig    = 5
	PosMarkToMark   = 6
	PosContext      = 7
	PosChainContext = 8
	posExtension    = 9
)

// maxGlyphs bounds the glyphs a coverage or class definition expands to:
// a font has at most 65,536 glyphs.
const maxGlyphs = 1 << 16

// maxElements bounds what a table is read into: the scripts and their
// language systems, the features of each, the lookups of the features and
// the subtables of the lookups, counted together.
const maxElements = 1 << 20

// maxSteps bounds the entries the lookups of a table of n bytes list
// (glyphs covered, substitutions, pairs, rules, anchors) in all the calls
// made of the table. Listing every lookup of a font a few times takes
// less than four entries for a byte of the table, and less than a million
// in all but the largest fonts.
func maxSteps(n int) int { return min(1<<22+32*n, 1<<28) }

// data is a table or part of one; reads past its end return 0.
type data []byte

func (d data) u16(i int) int {
	if i < 0 || i+2 > len(d) {
		return 0
	}
	return int(binary.BigEndian.Uint16(d[i:]))
}

func (d data) i16(i int) int { return int(int16(d.u16(i))) }

func (d data) u32(i int) int {
	if i < 0 || i+4 > len(d) {
		return 0
	}
	return int(binary.BigEndian.Uint32(d[i:]))
}

func (d data) tag(i int) string {
	if i < 0 || i+4 > len(d) {
		return ""
	}
	return string(d[i : i+4])
}

// at returns d from off, or nil when off is 0 or past the end.
func (d data) at(off int) data {
	if off <= 0 || off >= len(d) {
		return nil
	}
	return d[off:]
}

// count returns n, or fewer when the table cannot hold n records of size
// bytes from off: a damaged count must not make a large allocation.
func (d data) count(n, off, size int) int {
	if n < 0 || off < 0 {
		return 0
	}
	if size > 0 && off+n*size > len(d) {
		n = max(0, (len(d)-off)/size)
	}
	return n
}

// ErrNotLayout is returned for data that is not a GSUB or GPOS table.
var ErrNotLayout = errors.New("otlayout: not a layout table")

// ErrTooLarge is returned for a table that holds more scripts, features,
// lookups and subtables than a table is read into.
var ErrTooLarge = errors.New("otlayout: the table holds too many features or lookups")

// take counts n elements against what a table may still be read into; it
// returns false when they are too many.
func take(left *int, n int) bool {
	*left -= n
	return *left >= 0
}

// spend counts n entries a lookup lists against what the table may still
// list (maxSteps); it returns false when they are too many, and the lookup
// stops there.
func (t *Table) spend(n int) bool {
	t.steps -= n
	return t.steps >= 0
}

// Truncated reports whether lookups were cut short: the table lists more
// than any font does, and what was read of it is not all it says.
func (t *Table) Truncated() bool { return t.steps < 0 }

// ParseGSUB reads a GSUB table.
func ParseGSUB(b []byte) (*Table, error) { return parse(b, false) }

// ParseGPOS reads a GPOS table.
func ParseGPOS(b []byte) (*Table, error) { return parse(b, true) }

func parse(b []byte, gpos bool) (*Table, error) {
	d := data(b)
	if len(d) < 10 || d.u16(0) != 1 {
		return nil, ErrNotLayout
	}
	t := &Table{gpos: gpos, steps: maxSteps(len(b))}
	left := maxElements
	t.Scripts = parseScripts(d.at(d.u16(4)), &left)
	t.Features = parseFeatures(d.at(d.u16(6)), &left)
	t.Lookups = t.parseLookups(d.at(d.u16(8)), &left)
	if left < 0 {
		return nil, ErrTooLarge
	}
	if d.u16(2) >= 1 && len(d) >= 14 {
		if fv := d.at(d.u32(10)); fv != nil {
			t.FeatureVariations = fv.count(fv.u32(4), 8, 8)
		}
	}
	// feature indices outside the list are dropped
	for i := range t.Scripts {
		s := &t.Scripts[i]
		if s.Default != nil {
			t.clean(s.Default)
		}
		for j := range s.Langs {
			t.clean(&s.Langs[j])
		}
	}
	for i := range t.Features {
		f := &t.Features[i]
		keep := f.Lookups[:0]
		for _, l := range f.Lookups {
			if l < len(t.Lookups) {
				keep = append(keep, l)
			}
		}
		f.Lookups = keep
	}
	return t, nil
}

func (t *Table) clean(l *LangSys) {
	if l.Required >= len(t.Features) {
		l.Required = -1
	}
	keep := l.Features[:0]
	for _, f := range l.Features {
		if f < len(t.Features) {
			keep = append(keep, f)
		}
	}
	l.Features = keep
}

// parseScripts reads the script list; left is how many elements the table
// may still be read into, below zero when the list holds more.
func parseScripts(d data, left *int) []Script {
	if d == nil {
		return nil
	}
	n := d.count(d.u16(0), 2, 6)
	if !take(left, n) {
		return nil
	}
	out := make([]Script, 0, n)
	for i := range n {
		rec := 2 + i*6
		s := d.at(d.u16(rec + 4))
		if s == nil {
			continue
		}
		sc := Script{Tag: d.tag(rec)}
		if ls := s.at(s.u16(0)); ls != nil {
			l := parseLangSys(ls, "", left)
			sc.Default = &l
		}
		m := s.count(s.u16(2), 4, 6)
		if !take(left, m) {
			return nil
		}
		for j := range m {
			r := 4 + j*6
			if ls := s.at(s.u16(r + 4)); ls != nil {
				sc.Langs = append(sc.Langs, parseLangSys(ls, s.tag(r), left))
			}
		}
		if *left < 0 {
			return nil
		}
		out = append(out, sc)
	}
	return out
}

func parseLangSys(d data, tag string, left *int) LangSys {
	l := LangSys{Tag: tag, Required: -1}
	if r := d.u16(2); r != 0xFFFF {
		l.Required = r
	}
	n := d.count(d.u16(4), 6, 2)
	if !take(left, n) {
		return l
	}
	l.Features = make([]int, n)
	for i := range n {
		l.Features[i] = d.u16(6 + i*2)
	}
	return l
}

func parseFeatures(d data, left *int) []Feature {
	if d == nil {
		return nil
	}
	n := d.count(d.u16(0), 2, 6)
	if !take(left, n) {
		return nil
	}
	out := make([]Feature, n)
	for i := range n {
		rec := 2 + i*6
		f := Feature{Tag: d.tag(rec)}
		if ft := d.at(d.u16(rec + 4)); ft != nil {
			m := ft.count(ft.u16(2), 4, 2)
			if !take(left, m) {
				return nil
			}
			f.Lookups = make([]int, m)
			for j := range m {
				f.Lookups[j] = ft.u16(4 + j*2)
			}
			if p := ft.at(ft.u16(0)); p != nil {
				f.params(p, left)
			}
		}
		if *left < 0 {
			return nil
		}
		out[i] = f
	}
	return out
}

// params reads the feature parameters of stylistic sets and character
// variants.
func (f *Feature) params(p data, left *int) {
	switch {
	case len(f.Tag) == 4 && f.Tag[:2] == "ss" && digits(f.Tag[2:]):
		if p.u16(0) == 0 && len(p) >= 4 {
			f.UIName = uint16(p.u16(2))
		}
	case len(f.Tag) == 4 && f.Tag[:2] == "cv" && digits(f.Tag[2:]):
		if p.u16(0) != 0 || len(p) < 14 {
			return
		}
		f.UIName = uint16(p.u16(2))
		n := p.count(p.u16(12), 14, 3)
		if !take(left, n) {
			return
		}
		for i := range n {
			o := 14 + i*3
			f.Chars = append(f.Chars, rune(int(p[o])<<16|int(p[o+1])<<8|int(p[o+2])))
		}
	}
}

func digits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return s != ""
}

func (t *Table) parseLookups(d data, left *int) []Lookup {
	if d == nil {
		return nil
	}
	ext := substExtension
	if t.gpos {
		ext = posExtension
	}
	n := d.count(d.u16(0), 2, 2)
	if !take(left, n) {
		return nil
	}
	out := make([]Lookup, n)
	for i := range n {
		l := Lookup{MarkSet: -1}
		lt := d.at(d.u16(2 + i*2))
		if lt == nil {
			out[i] = l
			continue
		}
		l.Type, l.Flag = lt.u16(0), uint16(lt.u16(2))
		m := lt.count(lt.u16(4), 6, 2)
		if !take(left, m) {
			return nil
		}
		if l.Flag&UseMarkFilteringSet != 0 {
			l.MarkSet = lt.u16(6 + m*2)
		}
		isExt := l.Type == ext
		if isExt {
			l.Type = 0 // until an extension subtable says
		}
		for j := range m {
			st := lt.at(lt.u16(6 + j*2))
			if st == nil {
				continue
			}
			if isExt {
				// an extension subtable: the type and a 32-bit offset
				typ := st.u16(2)
				if st.u16(0) != 1 || typ == ext || l.Type != 0 && typ != l.Type {
					continue
				}
				if s := st.at(st.u32(4)); s != nil {
					l.Type = typ
					l.subtables = append(l.subtables, s)
				}
				continue
			}
			l.subtables = append(l.subtables, st)
		}
		out[i] = l
	}
	return out
}

// ScriptFeatures returns, for each feature, the scripts and language
// systems that use it, as "script" (the script's default) or
// "script/LANG" in the order of the script list.
func (t *Table) ScriptFeatures() [][]string {
	out := make([][]string, len(t.Features))
	add := func(l *LangSys, name string) {
		seen := map[int]bool{}
		use := func(f int) {
			if f >= 0 && f < len(out) && !seen[f] {
				seen[f] = true
				out[f] = append(out[f], name)
			}
		}
		use(l.Required)
		for _, f := range l.Features {
			use(f)
		}
	}
	for _, s := range t.Scripts {
		if s.Default != nil {
			add(s.Default, s.Tag)
		}
		for i := range s.Langs {
			add(&s.Langs[i], s.Tag+"/"+s.Langs[i].Tag)
		}
	}
	return out
}

// coverage returns the glyphs of a coverage table in coverage index order.
func coverage(d data) []uint16 {
	if d == nil {
		return nil
	}
	switch d.u16(0) {
	case 1:
		n := d.count(d.u16(2), 4, 2)
		out := make([]uint16, n)
		for i := range n {
			out[i] = uint16(d.u16(4 + i*2))
		}
		return out
	case 2:
		n := d.count(d.u16(2), 4, 6)
		var out []uint16
		for i := range n {
			r := 4 + i*6
			start, end := d.u16(r), d.u16(r+2)
			for g := start; g <= end && len(out) < maxGlyphs; g++ {
				out = append(out, uint16(g))
			}
		}
		return out
	}
	return nil
}

// classDef returns the classes of a class definition table: glyph →
// class, without the glyphs of class 0.
func classDef(d data) map[uint16]int {
	out := map[uint16]int{}
	if d == nil {
		return out
	}
	switch d.u16(0) {
	case 1:
		start := d.u16(2)
		n := d.count(d.u16(4), 6, 2)
		for i := range n {
			if c := d.u16(6 + i*2); c != 0 && start+i < maxGlyphs {
				out[uint16(start+i)] = c
			}
		}
	case 2:
		n := d.count(d.u16(2), 4, 6)
		total := 0
		for i := range n {
			r := 4 + i*6
			start, end, c := d.u16(r), d.u16(r+2), d.u16(r+4)
			if c == 0 {
				continue
			}
			for g := start; g <= end && total < maxGlyphs; g++ {
				out[uint16(g)] = c
				total++
			}
		}
	}
	return out
}

// classGlyphs inverts a class definition: class → its glyphs, ascending.
func classGlyphs(cd map[uint16]int) map[int][]uint16 {
	out := map[int][]uint16{}
	for g, c := range cd {
		out[c] = append(out[c], g)
	}
	for _, gs := range out {
		slices.Sort(gs)
	}
	return out
}

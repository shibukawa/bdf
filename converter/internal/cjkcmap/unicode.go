package cjkcmap

import (
	"strings"
	"sync"
)

// The UTF-32 CMaps that give CID → Unicode, in order of preference. A CID
// that only a vertical CMap reaches is the vertical form of its character.
var unicodeSources = map[string][]string{
	"Adobe-Japan1": {"UniJIS-UTF32-H", "UniJIS2004-UTF32-H", "UniJISX0213-UTF32-H", "UniJISX02132004-UTF32-H",
		"UniJIS-UCS2-HW-H", // half-width Latin (CIDs 231–)
		"UniJIS-UTF32-V", "UniJIS2004-UTF32-V", "UniJISX0213-UTF32-V", "UniJISX02132004-UTF32-V", "UniJIS-UCS2-HW-V"},
	"Adobe-Japan2": {"UniHojo-UTF32-H", "UniHojo-UTF32-V"},
	"Adobe-GB1":    {"UniGB-UTF32-H", "UniGB-UTF32-V"},
	"Adobe-CNS1":   {"UniCNS-UTF32-H", "UniCNS-UTF32-V"},
	"Adobe-Korea1": {"UniKS-UTF32-H", "UniKS-UTF32-V"},
}

type unicodeTable struct {
	once  sync.Once
	tab   []uint32 // CID → code point | verticalBit; 0 when unmapped
	vform map[uint32]rune
}

const verticalBit = 1 << 31

var unicodeTables = map[string]*unicodeTable{}

func init() {
	for c := range unicodeSources {
		unicodeTables[c] = &unicodeTable{}
	}
}

func table(collection string) *unicodeTable {
	t := unicodeTables[collection]
	if t != nil {
		t.once.Do(func() { t.tab, t.vform = buildUnicode(unicodeSources[collection]) })
	}
	return t
}

// Unicode returns the character of a CID in a character collection
// ("Adobe-Japan1", …) as text, and whether the CID is the vertical form
// of that character (the glyph a vertical CMap maps it to).
func Unicode(collection string, cid uint32) (r rune, vertical, ok bool) {
	t := table(collection)
	if t == nil || int(cid) >= len(t.tab) || t.tab[cid] == 0 {
		return 0, false, false
	}
	v := t.tab[cid]
	return rune(v &^ verticalBit), v&verticalBit != 0, true
}

// VerticalForm returns the vertical presentation form (U+FE10–FE19,
// U+FE30–FE4F) that stands for a CID, such as ︑ for the vertical 、.
func VerticalForm(collection string, cid uint32) (rune, bool) {
	t := table(collection)
	if t == nil {
		return 0, false
	}
	r, ok := t.vform[cid]
	return r, ok
}

// buildUnicode inverts the sources. Among the characters mapped to a CID,
// printable ASCII wins, then ordinary characters, then presentation forms,
// then compatibility ideographs and radicals, then private use; then the
// earlier source, then the lower code point.
func buildUnicode(sources []string) ([]uint32, map[uint32]rune) {
	var tab []uint32
	var best []uint64
	vform := map[uint32]rune{}
	for prio, name := range sources {
		m := Lookup(name)
		if m == nil {
			continue
		}
		vertical := strings.HasSuffix(name, "-V")
		for _, r := range m.ranges {
			for k := uint64(0); k <= uint64(r.hi-r.lo); k++ {
				cp, cid := r.lo+uint32(k), r.cid+uint32(k)
				if cid >= 1<<17 || cp > 0x10ffff {
					continue
				}
				for int(cid) >= len(tab) {
					tab = append(tab, 0)
					best = append(best, ^uint64(0))
				}
				class := charClass(cp)
				if class == classVerticalForm {
					if _, ok := vform[cid]; !ok {
						vform[cid] = rune(cp)
					}
				}
				rank := uint64(class)<<48 | uint64(prio)<<32 | uint64(cp)
				if rank < best[cid] {
					best[cid] = rank
					tab[cid] = cp
					if vertical {
						tab[cid] |= verticalBit
					}
				}
			}
		}
	}
	return tab, vform
}

const (
	classASCII = iota
	classOrdinary
	classVerticalForm
	classCompatibility
	classPrivate
)

func charClass(cp uint32) int {
	switch {
	case cp >= 0x20 && cp < 0x7f:
		return classASCII
	case cp >= 0xFE10 && cp <= 0xFE19, cp >= 0xFE30 && cp <= 0xFE4F:
		return classVerticalForm
	case cp >= 0xF900 && cp <= 0xFAFF, cp >= 0x2F800 && cp <= 0x2FA1F, cp >= 0x2E80 && cp <= 0x2FDF:
		return classCompatibility // compatibility ideographs, radicals
	case cp >= 0xE000 && cp <= 0xF8FF, cp >= 0xF0000:
		return classPrivate
	}
	return classOrdinary
}

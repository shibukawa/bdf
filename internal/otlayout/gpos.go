package otlayout

import "math/bits"

// Value is how a positioning lookup moves a glyph, in font units.
type Value struct {
	XPlacement, YPlacement, XAdvance, YAdvance int
}

// Zero reports whether the value moves nothing.
func (v Value) Zero() bool { return v == Value{} }

// value reads a value record of the given format at off; the device and
// variation offsets are skipped.
func value(d data, off, format int) Value {
	var v Value
	for bit := 0; bit < 4; bit++ {
		if format&(1<<bit) == 0 {
			continue
		}
		n := d.i16(off)
		off += 2
		switch bit {
		case 0:
			v.XPlacement = n
		case 1:
			v.YPlacement = n
		case 2:
			v.XAdvance = n
		case 3:
			v.YAdvance = n
		}
	}
	return v
}

// valueSize is the size of a value record of a format.
func valueSize(format int) int { return 2 * bits.OnesCount16(uint16(format&0xFF)) }

// Singles calls fn with the glyphs a single adjustment lookup (GPOS type
// 1) moves and how, until fn returns false.
func (t *Table) Singles(i int, fn func(g uint16, v Value) bool) {
	if !t.gpos || i < 0 || i >= len(t.Lookups) || t.Lookups[i].Type != PosSingle {
		return
	}
	seen := map[uint16]bool{}
	for _, st := range t.Lookups[i].subtables {
		cov := coverage(st.at(st.u16(2)))
		f := st.u16(4)
		format, n := st.u16(0), st.count(st.u16(6), 8, valueSize(f))
		for k, g := range cov {
			var v Value
			switch {
			case format == 1:
				v = value(st, 6, f)
			case format == 2 && k < n:
				v = value(st, 8+k*valueSize(f), f)
			default:
				continue
			}
			if seen[g] {
				continue
			}
			seen[g] = true
			if !fn(g, v) {
				return
			}
		}
	}
}

// Pair is two glyphs a pair adjustment lookup moves when they meet.
type Pair struct {
	First, Second uint16
	V1, V2        Value
}

// pairSub is a pair adjustment subtable prepared for looking pairs up.
type pairSub struct {
	d        data
	format   int
	vf1, vf2 int
	cov      map[uint16]int // first glyph → coverage index
	class1   map[uint16]int
	class2   map[uint16]int
	n1, n2   int // class counts
}

func (t *Table) pairSubs(i int) []*pairSub {
	if !t.gpos || i < 0 || i >= len(t.Lookups) || t.Lookups[i].Type != PosPair {
		return nil
	}
	if t.pairCache == nil {
		t.pairCache = map[int][]*pairSub{}
	}
	if ps, ok := t.pairCache[i]; ok {
		return ps
	}
	var out []*pairSub
	for _, st := range t.Lookups[i].subtables {
		p := &pairSub{d: st, format: st.u16(0), vf1: st.u16(4), vf2: st.u16(6), cov: map[uint16]int{}}
		for k, g := range coverage(st.at(st.u16(2))) {
			if _, dup := p.cov[g]; !dup {
				p.cov[g] = k
			}
		}
		switch p.format {
		case 1:
		case 2:
			p.class1 = classDef(st.at(st.u16(8)))
			p.class2 = classDef(st.at(st.u16(10)))
			p.n1, p.n2 = st.u16(12), st.u16(14)
			rec := valueSize(p.vf1) + valueSize(p.vf2)
			if rec == 0 || p.n1*p.n2 > (len(st)-16)/rec {
				continue // no adjustments, or records that do not fit
			}
		default:
			continue
		}
		out = append(out, p)
	}
	t.pairCache[i] = out
	return out
}

// PairValue returns how pair adjustment lookup i moves first and second
// when they meet: the first subtable that has the pair decides.
func (t *Table) PairValue(i int, first, second uint16) (v1, v2 Value, ok bool) {
	for _, p := range t.pairSubs(i) {
		k, covered := p.cov[first]
		if !covered {
			continue
		}
		s1, s2 := valueSize(p.vf1), valueSize(p.vf2)
		switch p.format {
		case 1:
			set := p.d.at(p.d.u16(10 + k*2))
			if set == nil || k >= p.d.u16(8) {
				continue
			}
			n := set.count(set.u16(0), 2, 2+s1+s2)
			// the records are sorted by the second glyph
			lo, hi := 0, n
			for lo < hi {
				m := (lo + hi) / 2
				if uint16(set.u16(2+m*(2+s1+s2))) < second {
					lo = m + 1
				} else {
					hi = m
				}
			}
			if lo < n {
				r := 2 + lo*(2+s1+s2)
				if uint16(set.u16(r)) == second {
					return value(set, r+2, p.vf1), value(set, r+2+s1, p.vf2), true
				}
			}
		case 2:
			c1, c2 := p.class1[first], p.class2[second]
			if c1 >= p.n1 || c2 >= p.n2 {
				continue
			}
			r := 16 + (c1*p.n2+c2)*(s1+s2)
			return value(p.d, r, p.vf1), value(p.d, r+s1, p.vf2), true
		}
	}
	return
}

// PairCount returns about how many pairs lookup i adjusts: the pairs of
// each subtable with an adjustment that is not zero, counted per subtable
// (a pair two subtables list counts twice). For class-based subtables it
// is the product of the sizes of the classes, without the class of second
// glyphs no class lists.
func (t *Table) PairCount(i int) int {
	total := 0
	for _, p := range t.pairSubs(i) {
		s1, s2 := valueSize(p.vf1), valueSize(p.vf2)
		switch p.format {
		case 1:
			n := p.d.count(p.d.u16(8), 10, 2)
			for k := range n {
				set := p.d.at(p.d.u16(10 + k*2))
				m := set.count(set.u16(0), 2, 2+s1+s2)
				for j := range m {
					r := 2 + j*(2+s1+s2)
					if !value(set, r+2, p.vf1).Zero() || !value(set, r+2+s1, p.vf2).Zero() {
						total++
					}
				}
			}
		case 2:
			size1 := make([]int, p.n1)
			for g := range p.cov {
				if c := p.class1[g]; c < p.n1 {
					size1[c]++
				}
			}
			size2 := make([]int, p.n2)
			for _, c := range p.class2 {
				if c < p.n2 {
					size2[c]++
				}
			}
			for c1 := range p.n1 {
				for c2 := 1; c2 < p.n2; c2++ {
					r := 16 + (c1*p.n2+c2)*(s1+s2)
					if !value(p.d, r, p.vf1).Zero() || !value(p.d, r+s1, p.vf2).Zero() {
						total += size1[c1] * size2[c2]
					}
				}
			}
		}
	}
	return total
}

// Pairs calls fn with pairs of lookup i whose adjustment is not zero,
// subtable by subtable, until fn returns false. Class-based subtables are
// expanded to every pair of their classes (without the class of second
// glyphs no class lists), so there may be very many.
func (t *Table) Pairs(i int, fn func(Pair) bool) {
	for _, p := range t.pairSubs(i) {
		s1, s2 := valueSize(p.vf1), valueSize(p.vf2)
		switch p.format {
		case 1:
			for _, first := range coverage(p.d.at(p.d.u16(2))) {
				k := p.cov[first]
				if k >= p.d.u16(8) {
					continue
				}
				set := p.d.at(p.d.u16(10 + k*2))
				m := set.count(set.u16(0), 2, 2+s1+s2)
				for j := range m {
					r := 2 + j*(2+s1+s2)
					pr := Pair{First: first, Second: uint16(set.u16(r)), V1: value(set, r+2, p.vf1), V2: value(set, r+2+s1, p.vf2)}
					if (!pr.V1.Zero() || !pr.V2.Zero()) && !fn(pr) {
						return
					}
				}
			}
		case 2:
			// the class pairs that adjust something, by first class, read
			// once: the records hold n1 × n2 of them
			type adjust struct {
				c2     int
				v1, v2 Value
			}
			glyphs2 := classGlyphs(p.class2)
			byClass := make([][]adjust, p.n1)
			for c1 := range p.n1 {
				for c2 := 1; c2 < p.n2; c2++ {
					r := 16 + (c1*p.n2+c2)*(s1+s2)
					v1, v2 := value(p.d, r, p.vf1), value(p.d, r+s1, p.vf2)
					if (!v1.Zero() || !v2.Zero()) && len(glyphs2[c2]) > 0 {
						byClass[c1] = append(byClass[c1], adjust{c2, v1, v2})
					}
				}
			}
			for _, first := range coverage(p.d.at(p.d.u16(2))) {
				c1 := p.class1[first]
				if c1 >= p.n1 {
					continue
				}
				for _, a := range byClass[c1] {
					for _, second := range glyphs2[a.c2] {
						if !fn(Pair{First: first, Second: second, V1: a.v1, V2: a.v2}) {
							return
						}
					}
				}
			}
		}
	}
}

// Anchor is a point a glyph attaches to, in font units.
type Anchor struct{ X, Y int }

func anchor(d data) (Anchor, bool) {
	if d == nil || d.u16(0) < 1 || d.u16(0) > 3 {
		return Anchor{}, false
	}
	return Anchor{d.i16(2), d.i16(4)}, true
}

// Attachment is a subtable of a mark to base or mark to mark lookup: the
// class and anchor of each mark, and for each glyph the marks attach to,
// its anchor for each class.
type Attachment struct {
	Marks map[uint16]MarkAnchor
	Bases map[uint16][]*Anchor // indexed by mark class; nil when the base has none for the class
}

// MarkAnchor is the class of a mark and its anchor.
type MarkAnchor struct {
	Class  int
	Anchor Anchor
}

// Attachments returns the subtables of a mark to base (GPOS type 4) or
// mark to mark (type 6) lookup.
func (t *Table) Attachments(i int) []Attachment {
	if !t.gpos || i < 0 || i >= len(t.Lookups) {
		return nil
	}
	l := &t.Lookups[i]
	if l.Type != PosMarkToBase && l.Type != PosMarkToMark {
		return nil
	}
	var out []Attachment
	for _, st := range l.subtables {
		if st.u16(0) != 1 {
			continue
		}
		marks := coverage(st.at(st.u16(2)))
		bases := coverage(st.at(st.u16(4)))
		classes := st.u16(6)
		ma, ba := st.at(st.u16(8)), st.at(st.u16(10))
		if ma == nil || ba == nil || classes == 0 {
			continue
		}
		a := Attachment{Marks: map[uint16]MarkAnchor{}, Bases: map[uint16][]*Anchor{}}
		n := ma.count(ma.u16(0), 2, 4)
		for k, g := range marks {
			if k >= n {
				break
			}
			c := ma.u16(2 + k*4)
			if an, ok := anchor(ma.at(ma.u16(2 + k*4 + 2))); ok && c < classes {
				a.Marks[g] = MarkAnchor{c, an}
			}
		}
		n = ba.count(ba.u16(0), 2, 2*classes)
		for k, g := range bases {
			if k >= n {
				break
			}
			anchors := make([]*Anchor, classes)
			has := false
			for c := range classes {
				// base anchor offsets are from the base array
				if an, ok := anchor(ba.at(ba.u16(2 + (k*classes+c)*2))); ok {
					anchors[c] = &an
					has = true
				}
			}
			if has {
				a.Bases[g] = anchors
			}
		}
		out = append(out, a)
	}
	return out
}

// MarkCounts returns how many marks a mark attachment lookup (GPOS types
// 4, 5 and 6) attaches, and to how many glyphs; for a cursive attachment
// lookup (type 3), the glyphs it joins.
func (t *Table) MarkCounts(i int) (marks, bases int) {
	if !t.gpos || i < 0 || i >= len(t.Lookups) {
		return 0, 0
	}
	l := &t.Lookups[i]
	m, b := map[uint16]bool{}, map[uint16]bool{}
	for _, st := range l.subtables {
		if st.u16(0) != 1 {
			continue
		}
		switch l.Type {
		case PosCursive:
			for _, g := range coverage(st.at(st.u16(2))) {
				b[g] = true
			}
		case PosMarkToBase, PosMarkToLig, PosMarkToMark:
			for _, g := range coverage(st.at(st.u16(2))) {
				m[g] = true
			}
			for _, g := range coverage(st.at(st.u16(4))) {
				b[g] = true
			}
		}
	}
	return len(m), len(b)
}

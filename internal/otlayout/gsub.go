package otlayout

import "strconv"

// Subst is a substitution of a GSUB lookup.
type Subst struct {
	// In is the glyph replaced, or the components of a ligature.
	In []uint16
	// Out is what replaces it: a glyph, a sequence (a multiple substitution;
	// empty deletes the glyph), or the glyphs an alternate substitution
	// chooses from.
	Out []uint16
}

// Substs calls fn with the substitutions of GSUB lookup i (a single,
// multiple, alternate, ligature or reverse chaining lookup) in the order
// of its subtables and their coverage, until fn returns false. What an
// earlier subtable substitutes is not listed again from a later one, since
// the first subtable that applies wins; the substitutions of a reverse
// chaining lookup depend on their context and are all listed.
func (t *Table) Substs(i int, fn func(Subst) bool) {
	if t.gpos || i < 0 || i >= len(t.Lookups) {
		return
	}
	l := &t.Lookups[i]
	seen := map[string]bool{}
	emit := func(s Subst) bool {
		if l.Type != SubstReverseChain {
			k := key(s.In)
			if seen[k] {
				return true
			}
			seen[k] = true
		}
		return fn(s)
	}
	for _, st := range l.subtables {
		if !substs(l.Type, st, emit) {
			return
		}
	}
}

// key identifies a glyph sequence.
func key(gs []uint16) string {
	b := make([]byte, 0, len(gs)*2)
	for _, g := range gs {
		b = append(b, byte(g>>8), byte(g))
	}
	return string(b)
}

// substs lists the substitutions of one subtable; it returns false when
// fn asked to stop.
func substs(typ int, d data, fn func(Subst) bool) bool {
	switch typ {
	case SubstSingle:
		cov := coverage(d.at(d.u16(2)))
		switch d.u16(0) {
		case 1:
			delta := d.i16(4)
			for _, g := range cov {
				if !fn(Subst{In: []uint16{g}, Out: []uint16{uint16(int(g) + delta)}}) {
					return false
				}
			}
		case 2:
			n := d.count(d.u16(4), 6, 2)
			for k, g := range cov {
				if k >= n {
					break
				}
				if !fn(Subst{In: []uint16{g}, Out: []uint16{uint16(d.u16(6 + k*2))}}) {
					return false
				}
			}
		}
	case SubstMultiple, SubstAlternate:
		if d.u16(0) != 1 {
			return true
		}
		cov := coverage(d.at(d.u16(2)))
		n := d.count(d.u16(4), 6, 2)
		for k, g := range cov {
			if k >= n {
				break
			}
			seq := d.at(d.u16(6 + k*2))
			if seq == nil {
				continue
			}
			if !fn(Subst{In: []uint16{g}, Out: glyphArray(seq, 0)}) {
				return false
			}
		}
	case SubstLigature:
		if d.u16(0) != 1 {
			return true
		}
		cov := coverage(d.at(d.u16(2)))
		n := d.count(d.u16(4), 6, 2)
		for k, g := range cov {
			if k >= n {
				break
			}
			set := d.at(d.u16(6 + k*2))
			if set == nil {
				continue
			}
			m := set.count(set.u16(0), 2, 2)
			for j := range m {
				lig := set.at(set.u16(2 + j*2))
				if lig == nil {
					continue
				}
				comps := lig.count(max(lig.u16(2)-1, 0), 4, 2)
				in := make([]uint16, 1, comps+1)
				in[0] = g
				for c := range comps {
					in = append(in, uint16(lig.u16(4+c*2)))
				}
				if !fn(Subst{In: in, Out: []uint16{uint16(lig.u16(0))}}) {
					return false
				}
			}
		}
	case SubstReverseChain:
		if d.u16(0) != 1 {
			return true
		}
		cov := coverage(d.at(d.u16(2)))
		p := 4
		p += 2 + d.u16(p)*2 // backtrack coverages
		p += 2 + d.u16(p)*2 // lookahead coverages
		n := d.count(d.u16(p), p+2, 2)
		for k, g := range cov {
			if k >= n {
				break
			}
			if !fn(Subst{In: []uint16{g}, Out: []uint16{uint16(d.u16(p + 2 + k*2))}}) {
				return false
			}
		}
	}
	return true
}

// glyphArray reads a count at off followed by that many glyph IDs.
func glyphArray(d data, off int) []uint16 {
	n := d.count(d.u16(off), off+2, 2)
	out := make([]uint16, n)
	for i := range n {
		out[i] = uint16(d.u16(off + 2 + i*2))
	}
	return out
}

// Context is what a contextual lookup does: how many rules it has and the
// lookups its rules apply to the glyphs they match.
type Context struct {
	Rules int
	// Nested are the lookups the rules apply, in the order they first
	// appear, each once.
	Nested []int
}

// Context returns the rules of a contextual lookup (GSUB types 5, 6 and
// 8, GPOS types 7 and 8); ok is false for other lookups. A reverse
// chaining lookup has a rule per subtable and applies no other lookup.
func (t *Table) Context(i int) (c Context, ok bool) {
	if i < 0 || i >= len(t.Lookups) {
		return c, false
	}
	l := &t.Lookups[i]
	ctx, chained := SubstContext, SubstChainContext
	if t.gpos {
		ctx, chained = PosContext, PosChainContext
	}
	switch {
	case l.Type == ctx || l.Type == chained:
		seen := map[int]bool{}
		for _, st := range l.subtables {
			seqContext(st, l.Type == chained, &c, seen, len(t.Lookups))
		}
		return c, true
	case !t.gpos && l.Type == SubstReverseChain:
		c.Rules = len(l.subtables)
		return c, true
	}
	return c, false
}

// seqContext reads the rules of a (chained) sequence context subtable.
func seqContext(d data, chained bool, c *Context, seen map[int]bool, lookups int) {
	records := func(r data, off, n int) {
		n = r.count(n, off, 4)
		for k := range n {
			li := r.u16(off + k*4 + 2)
			if li < lookups && !seen[li] {
				seen[li] = true
				c.Nested = append(c.Nested, li)
			}
		}
	}
	// rule reads one rule of formats 1 and 2
	rule := func(r data) {
		if r == nil {
			return
		}
		c.Rules++
		if !chained {
			glyphs, n := r.u16(0), r.u16(2)
			records(r, 4+max(glyphs-1, 0)*2, n)
			return
		}
		p := 0
		p += 2 + r.u16(p)*2           // backtrack
		p += 2 + max(r.u16(p)-1, 0)*2 // input, less the first glyph
		p += 2 + r.u16(p)*2           // lookahead
		records(r, p+2, r.u16(p))
	}
	ruleSets := func(off int) {
		n := d.count(d.u16(off), off+2, 2)
		for k := range n {
			set := d.at(d.u16(off + 2 + k*2))
			if set == nil {
				continue
			}
			m := set.count(set.u16(0), 2, 2)
			for j := range m {
				rule(set.at(set.u16(2 + j*2)))
			}
		}
	}
	switch d.u16(0) {
	case 1:
		ruleSets(4)
	case 2:
		if chained {
			ruleSets(10)
		} else {
			ruleSets(6)
		}
	case 3:
		c.Rules++
		if !chained {
			glyphs, n := d.u16(2), d.u16(4)
			records(d, 6+glyphs*2, n)
			return
		}
		p := 2
		p += 2 + d.u16(p)*2 // backtrack coverages
		p += 2 + d.u16(p)*2 // input coverages
		p += 2 + d.u16(p)*2 // lookahead coverages
		records(d, p+2, d.u16(p))
	}
}

// Covered returns the glyphs lookup i applies to: those its subtables
// cover (for a contextual lookup, the first glyph of its rules).
func (t *Table) Covered(i int) []uint16 {
	if i < 0 || i >= len(t.Lookups) {
		return nil
	}
	l := &t.Lookups[i]
	seen := map[uint16]bool{}
	var out []uint16
	for _, st := range l.subtables {
		var cov []uint16
		switch {
		case st.u16(0) == 3 && (l.Type == SubstContext && !t.gpos || l.Type == PosContext && t.gpos):
			cov = coverage(st.at(st.u16(6)))
		case st.u16(0) == 3 && (l.Type == SubstChainContext && !t.gpos || l.Type == PosChainContext && t.gpos):
			p := 2 + 2 + st.u16(2)*2 // past the backtrack coverages, at the input count
			cov = coverage(st.at(st.u16(p + 2)))
		case t.gpos && (l.Type == PosMarkToBase || l.Type == PosMarkToLig || l.Type == PosMarkToMark):
			cov = coverage(st.at(st.u16(2))) // the marks
		default:
			cov = coverage(st.at(st.u16(2)))
		}
		for _, g := range cov {
			if !seen[g] {
				seen[g] = true
				out = append(out, g)
			}
		}
	}
	return out
}

// TypeName names a lookup type of the table.
func (t *Table) TypeName(typ int) string {
	names := gsubTypes
	if t.gpos {
		names = gposTypes
	}
	if typ > 0 && typ < len(names) && names[typ] != "" {
		return names[typ]
	}
	return "type " + strconv.Itoa(typ)
}

var gsubTypes = []string{"", "single", "multiple", "alternate", "ligature", "contextual", "chained contextual", "extension", "reverse chaining"}

var gposTypes = []string{"", "single adjustment", "pair adjustment", "cursive attachment", "mark to base", "mark to ligature", "mark to mark", "contextual", "chained contextual", "extension"}

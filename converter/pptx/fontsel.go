package pptx

import (
	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// faceChoice is a font face picked for some text, with its metrics.
type faceChoice struct {
	l           *fontdb.Loaded // nil when no font file is available
	use         fontUse
	asc, desc   float64 // em
	ulPos, ulTh float64 // em, underline position (negative: below) and thickness
}

type resolveKey struct {
	name              string
	bold, italic, cjk bool
}

type fallbackKey struct {
	primary *faceChoice
	r       rune
}

// choice resolves a family request to a face (cached).
func (c *converter) choice(name string, bold, italic, cjk bool) *faceChoice {
	k := resolveKey{name, bold, italic, cjk}
	if fc, ok := c.choices[k]; ok {
		return fc
	}
	res := c.db.Resolve(name, bold, italic, cjk)
	fc := c.fromResolved(res, name, bold, italic)
	c.choices[k] = fc
	return fc
}

func (c *converter) fromResolved(res fontdb.Resolved, requested string, bold, italic bool) *faceChoice {
	fc := &faceChoice{
		use: fontUse{face: res.Face, requested: requested, generic: res.Generic, bold: bold, italic: italic,
			synthBold: res.SynthBold, synthItalic: res.SynthItalic},
		asc: 0.9, desc: 0.25, ulPos: -0.1, ulTh: 0.05,
	}
	if fc.use.generic == "" {
		fc.use.generic = fontdb.Classify(requested)
	}
	if res.Face != nil {
		if l, err := res.Face.Load(); err == nil {
			fc.l = l
			fc.asc, fc.desc = l.Ascent, l.Descent
			fc.ulPos, fc.ulTh = l.UnderlinePos, l.UnderTh
		} else {
			c.warnOnce("fontload:"+res.Face.Path, "font %s: %v", res.Face.Path, err)
			fc.use.face = nil
		}
	}
	return fc
}

// faceFor picks the face that draws r: the run's East Asian font for CJK
// characters and its Latin font otherwise, then the other one, then the
// generic fallbacks, so that no character is left to the viewer's fonts
// when some available font has it.
func (c *converter) faceFor(st *runStyle, r rune) *faceChoice {
	cjk := fontdb.IsCJK(r)
	primary := st.latin
	other := st.ea
	if cjk {
		primary, other = st.ea, st.latin
		if primary == "" {
			primary = other
		}
	}
	fc := c.choice(primary, st.bold, st.italic, cjk)
	if fc.l == nil || fc.l.Has(r) || r == '\t' || r == '\n' {
		return fc
	}
	k := fallbackKey{fc, r}
	if f, ok := c.fallback[k]; ok {
		return f
	}
	best := fc
	if other != "" && other != primary {
		if o := c.choice(other, st.bold, st.italic, !cjk); o.l != nil && o.l.Has(r) {
			best = o
		}
	}
	if best == fc {
		fk := fallbackListKey{fc.use.generic, st.bold, st.italic}
		list, ok := c.fallbackLists[fk]
		if !ok {
			list = c.db.Fallbacks(fc.use.generic, st.bold, st.italic)
			c.fallbackLists[fk] = list
		}
		for _, res := range list {
			l, err := res.Face.Load()
			if err != nil || !l.Has(r) {
				continue
			}
			best = c.fromResolvedCached(res, primary, st.bold, st.italic)
			break
		}
	}
	c.fallback[k] = best
	return best
}

type fallbackListKey struct {
	generic      string
	bold, italic bool
}

func (c *converter) fromResolvedCached(res fontdb.Resolved, requested string, bold, italic bool) *faceChoice {
	k := resolveKey{"\x00" + res.Face.Path + "#" + itoa(res.Face.Index) + "#" + requested, bold, italic, false}
	if fc, ok := c.choices[k]; ok {
		return fc
	}
	fc := c.fromResolved(res, requested, bold, italic)
	c.choices[k] = fc
	return fc
}

// advance measures r in em and records it for the embedded subset.
func (c *converter) advance(fc *faceChoice, r rune) float64 {
	if fc.l != nil {
		if g, ok := fc.l.Glyph(r); ok {
			m := c.faceRunes[fc.l.Face]
			if m == nil {
				m = map[rune]bool{}
				c.faceRunes[fc.l.Face] = m
			}
			m[r] = true
			return fc.l.Advance(g)
		}
	}
	switch {
	case r == ' ' || r == ' ':
		return 0.25
	case fontdb.IsCJK(r):
		return 1
	case r < 0x20 || r == 0x200b || r == 0xfeff:
		return 0
	}
	return 0.55
}

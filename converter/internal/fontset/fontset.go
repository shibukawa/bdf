// Package fontset picks the fonts that lay out the text of a document being
// converted, measures the characters with them, and embeds the faces in use
// as subsets of those characters. BDF has no layout engine, so converters
// of documents that are laid out when they are opened (Office documents,
// draw.io diagrams) measure text here with the fonts that are then
// embedded; the requested families resolve through fontdb, and characters
// the requested fonts lack fall back to other available fonts.
package fontset

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// Choice is a font face picked for some text, with its metrics.
type Choice struct {
	Loaded         *fontdb.Loaded // nil when no font file is available
	Use            Use
	Asc, Desc      float64 // em
	ULPos, ULThick float64 // em, underline position (negative: below) and thickness
	CapHeight      float64 // em, the height of capital letters
}

// Use is what a FONT instruction asks for: a resolved face (nil when no
// font file is available) and how the request looked, for the system-font
// fallback and synthetic styles. It is comparable, so that objects can
// collect the fonts they use in a map.
type Use struct {
	Face        *fontdb.Face
	Requested   string
	Generic     string
	Bold        bool
	Italic      bool
	SynthBold   bool
	SynthItalic bool
}

// Set is the fonts of one document: the faces its text resolved to, the
// characters measured with each, and after Embed, the embedded fonts.
type Set struct {
	db     *fontdb.DB
	warn   func(msg string)
	warned map[string]bool

	choices       map[resolveKey]*Choice
	fallback      map[fallbackKey]*Choice
	fallbackLists map[fallbackListKey][]fontdb.Resolved
	runes         map[*fontdb.Face]map[rune]bool // characters measured per face
	missing       map[rune]bool                  // characters no available font has
	embedded      map[*fontdb.Face]bdf.Hash
}

type resolveKey struct {
	name              string
	bold, italic, cjk bool
}

type fallbackKey struct {
	primary *Choice
	r       rune
}

type fallbackListKey struct {
	generic      string
	bold, italic bool
}

// New returns an empty set that resolves fonts in db. warn receives
// non-fatal problems (fonts that cannot be loaded or embedded).
func New(db *fontdb.DB, warn func(msg string)) *Set {
	return &Set{db: db, warn: warn, warned: map[string]bool{},
		choices: map[resolveKey]*Choice{}, fallback: map[fallbackKey]*Choice{},
		fallbackLists: map[fallbackListKey][]fontdb.Resolved{},
		runes:         map[*fontdb.Face]map[rune]bool{}, missing: map[rune]bool{},
		embedded: map[*fontdb.Face]bdf.Hash{}}
}

func (s *Set) warnf(format string, args ...any) {
	if s.warn != nil {
		s.warn(fmt.Sprintf(format, args...))
	}
}

// Choose resolves a family request to a face (cached). cjk asks for a face
// for East Asian text, which a missing family then falls back to.
func (s *Set) Choose(name string, bold, italic, cjk bool) *Choice {
	k := resolveKey{name, bold, italic, cjk}
	if fc, ok := s.choices[k]; ok {
		return fc
	}
	res := s.db.Resolve(name, bold, italic, cjk)
	fc := s.fromResolved(res, name, bold, italic)
	s.choices[k] = fc
	return fc
}

func (s *Set) fromResolved(res fontdb.Resolved, requested string, bold, italic bool) *Choice {
	fc := &Choice{
		Use: Use{Face: res.Face, Requested: requested, Generic: res.Generic, Bold: bold, Italic: italic,
			SynthBold: res.SynthBold, SynthItalic: res.SynthItalic},
		Asc: 0.9, Desc: 0.25, ULPos: -0.1, ULThick: 0.05, CapHeight: 0.7,
	}
	if fc.Use.Generic == "" {
		fc.Use.Generic = fontdb.Classify(requested)
	}
	if res.Face != nil {
		if l, err := res.Face.Load(); err == nil {
			fc.Loaded = l
			fc.Asc, fc.Desc = l.Ascent, l.Descent
			fc.ULPos, fc.ULThick = l.UnderlinePos, l.UnderTh
			fc.CapHeight = l.CapHeight
		} else {
			if key := res.Face.Path; !s.warned[key] {
				s.warned[key] = true
				s.warnf("font %s: %v", res.Face.Path, err)
			}
			fc.Use.Face = nil
		}
	}
	return fc
}

// FaceFor picks the face that draws r in text whose Latin font is latin
// and East Asian font is ea: the East Asian font for CJK characters and
// the Latin font otherwise, then the other one, then the generic
// fallbacks, so that no character is left to the viewer's fonts when some
// available font has it.
func (s *Set) FaceFor(latin, ea string, bold, italic bool, r rune) *Choice {
	cjk := fontdb.IsCJK(r)
	primary := latin
	other := ea
	if cjk {
		primary, other = ea, latin
		if primary == "" {
			primary = other
		}
	}
	fc := s.Choose(primary, bold, italic, cjk)
	if fc.Loaded == nil || fc.Loaded.Has(r) || r == '\t' || r == '\n' {
		return fc
	}
	k := fallbackKey{fc, r}
	if f, ok := s.fallback[k]; ok {
		return f
	}
	best := fc
	if other != "" && other != primary {
		if o := s.Choose(other, bold, italic, !cjk); o.Loaded != nil && o.Loaded.Has(r) {
			best = o
		}
	}
	if best == fc {
		best = s.genericFallback(fc, primary, bold, italic, r)
	}
	s.fallback[k] = best
	return best
}

// FaceForFamilies picks the face that draws r in text whose font is a CSS
// font-family list (draw.io labels): the first family, then the first of
// the others that has r, then the generic fallbacks of the first. Unlike
// FaceFor, CJK characters are not sent to a separate East Asian font; the
// list says which family draws them.
func (s *Set) FaceForFamilies(families []string, bold, italic bool, r rune) *Choice {
	cjk := fontdb.IsCJK(r)
	if len(families) == 0 {
		families = []string{"Helvetica"}
	}
	fc := s.Choose(families[0], bold, italic, cjk)
	if fc.Loaded == nil || fc.Loaded.Has(r) || r == '\t' || r == '\n' {
		return fc
	}
	k := fallbackKey{fc, r}
	if f, ok := s.fallback[k]; ok {
		return f
	}
	best := fc
	for _, fam := range families[1:] {
		if o := s.Choose(fam, bold, italic, cjk); o.Loaded != nil && o.Loaded.Has(r) {
			best = o
			break
		}
	}
	if best == fc {
		best = s.genericFallback(fc, families[0], bold, italic, r)
	}
	s.fallback[k] = best
	return best
}

// genericFallback returns the first of the fallbacks of fc's generic family
// that has r, or fc when none has it.
func (s *Set) genericFallback(fc *Choice, requested string, bold, italic bool, r rune) *Choice {
	fk := fallbackListKey{fc.Use.Generic, bold, italic}
	list, ok := s.fallbackLists[fk]
	if !ok {
		list = s.db.Fallbacks(fc.Use.Generic, bold, italic)
		s.fallbackLists[fk] = list
	}
	for _, res := range list {
		l, err := res.Face.Load()
		if err != nil || !l.Has(r) {
			continue
		}
		return s.fromResolvedCached(res, requested, bold, italic)
	}
	return fc
}

func (s *Set) fromResolvedCached(res fontdb.Resolved, requested string, bold, italic bool) *Choice {
	k := resolveKey{"\x00" + res.Face.Path + "#" + strconv.Itoa(res.Face.Index) + "#" + requested, bold, italic, false}
	if fc, ok := s.choices[k]; ok {
		return fc
	}
	fc := s.fromResolved(res, requested, bold, italic)
	s.choices[k] = fc
	return fc
}

// Advance measures r in em and records it for the embedded subset. A
// character the face lacks gets an estimated advance.
func (s *Set) Advance(fc *Choice, r rune) float64 {
	if fc.Loaded != nil {
		if g, ok := fc.Loaded.Glyph(r); ok {
			m := s.runes[fc.Loaded.Face]
			if m == nil {
				m = map[rune]bool{}
				s.runes[fc.Loaded.Face] = m
			}
			m[r] = true
			return fc.Loaded.Advance(g)
		}
	}
	if fc.Loaded != nil && !unicode.IsSpace(r) && r >= 0x20 && r != 0x200b && r != 0xfeff {
		s.missing[r] = true
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

// ReportMissing warns about the characters that no available font has.
func (s *Set) ReportMissing() {
	if len(s.missing) == 0 {
		return
	}
	var rs []rune
	for r := range s.missing {
		rs = append(rs, r)
	}
	slices.Sort(rs)
	var b strings.Builder
	for i, r := range rs {
		if i == 20 {
			fmt.Fprintf(&b, " … (%d more)", len(rs)-20)
			break
		}
		fmt.Fprintf(&b, " %c U+%04X", r, r)
	}
	s.warnf("no available font has glyphs for:%s; viewers draw them with their own fonts", b.String())
}

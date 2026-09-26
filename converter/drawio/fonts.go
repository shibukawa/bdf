package drawio

import (
	"sort"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/woff2"
)

// Fonts: labels are laid out here with the metrics of real fonts, which
// are then embedded as subsets of the characters used, as the PowerPoint
// converter does (docs/design.md §3.4).

// fontUse is what a FONT instruction asks for: a resolved face (nil when no
// font file is available) and how the request looked, for the system-font
// fallback and synthetic styles.
type fontUse struct {
	face        *fontdb.Face
	requested   string
	generic     string
	bold        bool
	italic      bool
	synthBold   bool
	synthItalic bool
}

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

type fallbackListKey struct {
	generic      string
	bold, italic bool
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

// faceFor picks the face that draws r in text of the given families (a CSS
// font-family list): the first family that has it, then the generic
// fallbacks, so that no character is left to the viewer's fonts when some
// available font has it.
func (c *converter) faceFor(families []string, bold, italic bool, r rune) *faceChoice {
	cjk := fontdb.IsCJK(r)
	if len(families) == 0 {
		families = []string{"Helvetica"}
	}
	fc := c.choice(families[0], bold, italic, cjk)
	if fc.l == nil || fc.l.Has(r) || r == '\t' || r == '\n' {
		return fc
	}
	k := fallbackKey{fc, r}
	if f, ok := c.fallback[k]; ok {
		return f
	}
	best := fc
	for _, fam := range families[1:] {
		if o := c.choice(fam, bold, italic, cjk); o.l != nil && o.l.Has(r) {
			best = o
			break
		}
	}
	if best == fc {
		fk := fallbackListKey{fc.use.generic, bold, italic}
		list, ok := c.fallbackLists[fk]
		if !ok {
			list = c.db.Fallbacks(fc.use.generic, bold, italic)
			c.fallbackLists[fk] = list
		}
		for _, res := range list {
			l, err := res.Face.Load()
			if err != nil || !l.Has(r) {
				continue
			}
			best = c.fromResolvedCached(res, families[0], bold, italic)
			break
		}
	}
	c.fallback[k] = best
	return best
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
	if fc.l != nil && !unicode.IsSpace(r) && r >= 0x20 && r != 0x200b && r != 0xfeff {
		c.missing[r] = true
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

// finalize builds the embedded fonts and encodes every canvas.
func (c *converter) finalize() {
	type faceKey = *fontdb.Face
	hashes := map[faceKey]bdf.Hash{}
	embedded := map[faceKey]bool{}
	if !c.opts.SystemFonts {
		faces := make([]faceKey, 0, len(c.faceRunes))
		for f := range c.faceRunes {
			faces = append(faces, f)
		}
		sort.Slice(faces, func(i, j int) bool {
			if faces[i].Path != faces[j].Path {
				return faces[i].Path < faces[j].Path
			}
			return faces[i].Index < faces[j].Index
		})
		for _, f := range faces {
			l, err := f.Load()
			if err != nil {
				continue
			}
			if !l.CanSubset(c.opts.IgnoreFSType) && !c.opts.NoSubset && l.Size() > maxWholeFont {
				c.warnf("font %s is not embedded: its outlines cannot be subset and the file is %d KB", f.Family, l.Size()/1024)
				continue
			}
			runes := make([]rune, 0, len(c.faceRunes[f]))
			for r := range c.faceRunes[f] {
				runes = append(runes, r)
			}
			data, ok := l.Program(runes, c.opts.NoSubset, c.opts.IgnoreFSType)
			if !ok {
				c.warnf("font %s is not embedded: its license does not allow embedding", f.Family)
				continue
			}
			if !c.opts.NoWOFF2 {
				if w, err := woff2.Encode(data); err == nil {
					data = w
				} else if err != woff2.ErrNotAvailable {
					c.warnf("font %s: not stored as WOFF2: %v", f.Family, err)
				}
			}
			hashes[f] = c.doc.AddFont(data)
			embedded[f] = true
		}
		c.embeddedFonts = len(hashes)
	}
	for _, cv := range c.canvases {
		if cv.encoded {
			continue
		}
		cv.encoded = true
		for i, u := range cv.fonts {
			cv.obj.UpdateFont(bdf.FontRef(i), fontRecord(u, hashes, embedded))
		}
		cv.hash, cv.bbox = c.doc.AddObject(cv.obj)
	}
}

// maxWholeFont is the largest font file embedded whole when it cannot be
// subset (CFF outlines).
const maxWholeFont = 2 << 20

func fontRecord(u fontUse, hashes map[*fontdb.Face]bdf.Hash, embedded map[*fontdb.Face]bool) bdf.Font {
	if u.face != nil && embedded[u.face] {
		// The face registered in the browser is the one that was measured;
		// only styles it lacks are synthesized.
		f := bdf.EmbeddedFont(hashes[u.face], 400, bdf.StyleNormal)
		if u.synthBold {
			f.Weight = 700
		}
		if u.synthItalic {
			f.Style = bdf.StyleItalic
		}
		f.Family = u.generic
		return f
	}
	family := cssQuote(u.requested)
	if u.face != nil && u.face.Family != u.requested {
		family += ", " + cssQuote(u.face.Family)
	}
	if family != "" {
		family += ", "
	}
	family += u.generic
	var w uint16 = 400
	if u.bold {
		w = 700
	}
	st := bdf.StyleNormal
	if u.italic {
		st = bdf.StyleItalic
	}
	return bdf.SystemFont(family, w, st)
}

func cssQuote(s string) string {
	if s == "" {
		return ""
	}
	out := []rune{'"'}
	for _, r := range s {
		if r == '"' || r == '\\' {
			out = append(out, '\\')
		}
		out = append(out, r)
	}
	return string(append(out, '"'))
}

// parseFontFamilies splits a CSS font-family list ("Helvetica, 'Noto Sans
// JP', sans-serif") into family names; generic names map to the fonts
// fontdb resolves them with.
func parseFontFamilies(s string) []string {
	var out []string
	for _, f := range splitArgs(s) {
		f = strings.TrimSpace(f)
		f = strings.Trim(f, `"'`)
		if f == "" {
			continue
		}
		switch strings.ToLower(f) {
		case "sans-serif", "system-ui", "-apple-system", "blinkmacsystemfont", "ui-sans-serif":
			f = fontdb.Sans
		case "serif", "ui-serif":
			f = fontdb.Serif
		case "monospace", "ui-monospace":
			f = fontdb.Mono
		}
		out = append(out, f)
	}
	return out
}

// Package fontchain finds the fonts text is drawn with, as a browser does:
// the font parts of a document, the families of the system that fonts
// refer to, and fallbacks for the characters those lack. The rasterizer
// and the vector exporters each keep a face of their own kind for a font
// program; a Program holds what they share of it.
package fontchain

import (
	"io/fs"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/internal/textrun"
)

// Program is a parsed font program with the metrics text is laid out by.
// A face embeds it.
type Program struct {
	Font     *sfnt.Font
	Upem     float64
	Asc, Dsc float64 // hhea ascent and descent in em, both positive
}

// NewProgram reads the metrics of a font program: 1000 units per em and
// an em box of 0.8 above and 0.2 below the baseline when it does not say.
func NewProgram(f *sfnt.Font) Program {
	p := Program{Font: f, Upem: float64(f.UnitsPerEm)}
	if p.Upem <= 0 {
		p.Upem = 1000
	}
	p.Asc, p.Dsc = float64(f.Ascent)/p.Upem, -float64(f.Descent)/p.Upem
	if p.Asc+p.Dsc <= 0 {
		p.Asc, p.Dsc = 0.8, 0.2
	}
	return p
}

// Glyph is the glyph of a character, when the font has one.
func (p *Program) Glyph(r rune) (uint16, bool) {
	g, ok := p.Font.Cmap[uint32(r)]
	return g, ok && g != 0
}

// Advance is the advance of a glyph in em.
func (p *Program) Advance(g uint16) float64 { return float64(p.Font.Advance(g)) / p.Upem }

// Metrics are the ascent and descent of the em box, in em.
func (p *Program) Metrics() (asc, dsc float64) { return p.Asc, p.Dsc }

// Face is a face of a user of the package: a pointer to its kind of face,
// which embeds a Program.
type Face interface {
	comparable
	Glyph(r rune) (uint16, bool)
	Advance(g uint16) float64
	Metrics() (asc, dsc float64)
}

// Use is a face with the styles a browser would synthesize for it.
type Use[F Face] struct {
	Face         F
	Bold, Italic bool
}

// Placed is a glyph positioned along the baseline.
type Placed[F Face] struct {
	Use Use[F]
	GID uint16
	R   rune
	X   float64 // pen position in units
}

// Chain is the faces a FONT instruction draws with, in the order
// characters are looked up.
type Chain[F Face] struct {
	Faces []Use[F]
	// the families of the system are resolved when the embedded face
	// lacks a glyph
	init      func()
	fonts     *Fonts[F]
	fallbacks []fontdb.Resolved
}

// Fonts loads the faces of a document: embedded font parts and, for the
// families fonts are referred to by, the fonts of the system.
type Fonts[F Face] struct {
	// Doc is the document whose font parts are loaded.
	Doc *bdf.Document
	// FontFS, FontDirs and NoSystemFonts say where the fonts of the system
	// are (see fontdb.New).
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// New makes a face of a font program; system tells a font of the
	// system from one of the document.
	New func(f *sfnt.Font, system bool) F
	// Warnf takes the warnings: a font that cannot be read, text without
	// a font.
	Warnf func(format string, args ...any)
	// NotDone says what is not done with what has no font: "not written"
	// or "not drawn".
	NotDone string
	// All are the faces loaded, in that order.
	All []F

	db     *fontdb.DB
	parts  map[bdf.Hash]F // the zero value: not usable
	files  map[*fontdb.Face]F
	chains map[bdf.Font]*Chain[F]
}

func (fs *Fonts[F]) ready() {
	if fs.parts == nil {
		fs.parts, fs.files, fs.chains = map[bdf.Hash]F{}, map[*fontdb.Face]F{}, map[bdf.Font]*Chain[F]{}
	}
}

// DB is the fonts of the system, found when first asked for.
func (fs *Fonts[F]) DB() *fontdb.DB {
	if fs.db == nil {
		fs.db = fontdb.New(fs.FontFS, fs.FontDirs, !fs.NoSystemFonts)
	}
	return fs.db
}

// Part is the face of a font part of the document, or the zero value for
// one that cannot be read (with a warning).
func (fs *Fonts[F]) Part(h bdf.Hash) F {
	fs.ready()
	var none F
	if fc, ok := fs.parts[h]; ok {
		return fc
	}
	fs.parts[h] = none
	p := fs.Doc.Part(h)
	if p == nil {
		fs.Warnf("font %s is missing", h)
		return none
	}
	data := p.Data
	if woff2.IsWOFF(data) {
		if string(data[:4]) != "wOF2" {
			fs.Warnf("WOFF 1 fonts are %s", fs.NotDone)
			return none
		}
		var err error
		if data, err = woff2.Decode(data); err != nil {
			fs.Warnf("font %s: %v", h, err)
			return none
		}
	}
	f, err := sfnt.Parse(data)
	if err != nil || f.Cmap == nil {
		fs.Warnf("font %s cannot be read", h)
		return none
	}
	fc := fs.New(f, false)
	fs.parts[h] = fc
	fs.All = append(fs.All, fc)
	return fc
}

// File is the face of a font of the system, or the zero value for one
// that cannot be loaded.
func (fs *Fonts[F]) File(face *fontdb.Face) F {
	fs.ready()
	var none F
	if fc, ok := fs.files[face]; ok {
		return fc
	}
	fs.files[face] = none
	l, err := face.Load()
	if err != nil {
		return none
	}
	fc := fs.New(l.Font, true)
	fs.files[face] = fc
	fs.All = append(fs.All, fc)
	return fc
}

// Chain resolves a font reference as a browser would: the embedded face
// (registered without weight or style, so those are synthesized), the
// families of its list that are installed, and fallbacks for the rest.
func (fs *Fonts[F]) Chain(f *bdf.Font) *Chain[F] {
	fs.ready()
	if f == nil {
		f = &bdf.Font{Kind: bdf.FontSystem, Family: "sans-serif", Weight: 400}
	}
	if c, ok := fs.chains[*f]; ok {
		return c
	}
	// (the caller may reuse its FONT instruction before a missing glyph
	// makes the chain resolve the remaining families)
	ref := *f
	f = &ref
	c := &Chain[F]{fonts: fs}
	fs.chains[*f] = c
	bold, italic := f.Weight >= 600, f.Style != bdf.StyleNormal
	var none F
	seen := map[F]bool{}
	add := func(fc F, b, i bool) {
		if fc != none && !seen[fc] {
			seen[fc] = true
			c.Faces = append(c.Faces, Use[F]{fc, b, i})
		}
	}
	if f.Kind == bdf.FontEmbedded {
		add(fs.Part(f.Hash), bold, italic)
	}
	c.init = func() {
		generic := ""
		db := fs.DB()
		for _, name := range textrun.SplitFamilies(f.Family) {
			if g, ok := textrun.GenericFamilies[strings.ToLower(name)]; ok {
				if generic == "" {
					generic = g
				}
				res := db.Resolve(g, bold, italic, false)
				if res.Face != nil {
					add(fs.File(res.Face), res.SynthBold, res.SynthItalic)
				}
				continue
			}
			if faces := db.Family(name); len(faces) > 0 {
				fa, sb, si := fontdb.Match(faces, bold, italic)
				add(fs.File(fa), sb, si)
			}
		}
		if generic == "" {
			generic = fontdb.Sans
			if len(c.Faces) == 0 {
				// a family that is not installed: its look-alike or generic kind
				for _, name := range textrun.SplitFamilies(f.Family) {
					if res := db.Resolve(name, bold, italic, false); res.Face != nil {
						generic = res.Generic
						add(fs.File(res.Face), res.SynthBold, res.SynthItalic)
						break
					}
				}
			}
		}
		c.fallbacks = db.Fallbacks(generic, bold, italic)
		// even text with no supported glyph needs the first usable face's
		// missing-glyph advance and baseline metrics
		if len(c.Faces) == 0 {
			for _, res := range c.fallbacks {
				add(fs.File(res.Face), res.SynthBold, res.SynthItalic)
				if len(c.Faces) > 0 {
					break
				}
			}
		}
		if len(c.Faces) == 0 {
			fs.Warnf("no font for %q: its text is %s", f.Family, fs.NotDone)
		}
	}
	return c
}

// InitFallbacks resolves the families of the system the chain refers to,
// once.
func (c *Chain[F]) InitFallbacks() {
	if c.init != nil {
		init := c.init
		c.init = nil
		init()
	}
}

// Glyph finds the face of the chain that has a glyph of a character, else
// a fallback that has, else the first face with its missing glyph (0).
func (c *Chain[F]) Glyph(r rune) (Use[F], uint16) {
	for _, fu := range c.Faces {
		if g, ok := fu.Face.Glyph(r); ok {
			return fu, g
		}
	}
	if c.init != nil {
		c.InitFallbacks()
		return c.Glyph(r)
	}
	var none F
	for _, res := range c.fallbacks {
		if !res.Face.HasRune(r) {
			continue
		}
		if fc := c.fonts.File(res.Face); fc != none {
			if g, ok := fc.Glyph(r); ok {
				return Use[F]{fc, res.SynthBold, res.SynthItalic}, g
			}
		}
	}
	return c.Faces[0], 0
}

// Layout places the glyphs of a string at size with letter spacing, in
// visual order, and returns them with the advance of the whole string.
func (c *Chain[F]) Layout(str string, size, spacing float64, dir byte) ([]Placed[F], float64) {
	if len(c.Faces) == 0 {
		c.InitFallbacks()
		if len(c.Faces) == 0 {
			return nil, 0
		}
	}
	runes := textrun.VisualOrder(str, dir == bdf.DirRTL)
	out := make([]Placed[F], 0, len(runes))
	x := 0.0
	for _, r := range runes {
		switch r {
		case '\t', '\n', '\f', '\r', '\v':
			r = ' ' // Canvas draws these as spaces
		}
		use, gid := c.Glyph(r)
		out = append(out, Placed[F]{Use: use, GID: gid, R: r, X: x})
		x += use.Face.Advance(gid)*size + spacing
	}
	return out, x
}

// BaselineShift is how far below the anchor the alphabetic baseline lies
// for a textBaseline, from the em box of the first face (as Blink does).
func (c *Chain[F]) BaselineShift(baseline byte, size float64) float64 {
	asc, dsc := 0.8, 0.2
	if len(c.Faces) > 0 {
		asc, dsc = c.Faces[0].Face.Metrics()
	}
	return textrun.BaselineShift(baseline, size, asc, dsc)
}

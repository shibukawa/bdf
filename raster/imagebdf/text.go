package imagebdf

import (
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/sfnt"
	"golang.org/x/text/unicode/bidi"
)

// fontFace is a parsed font program.
type fontFace struct {
	f        *sfnt.Font
	out      *sfnt.Outlines
	upem     float64
	asc, dsc float64 // hhea ascent and descent in em, both positive
	glyphs   map[uint16]*path
}

func newFace(f *sfnt.Font) *fontFace {
	ff := &fontFace{f: f, out: sfnt.NewOutlines(f), upem: float64(f.UnitsPerEm), glyphs: map[uint16]*path{}}
	if ff.upem <= 0 {
		ff.upem = 1000
	}
	ff.asc, ff.dsc = float64(f.Ascent)/ff.upem, -float64(f.Descent)/ff.upem
	if ff.asc+ff.dsc <= 0 {
		ff.asc, ff.dsc = 0.8, 0.2
	}
	return ff
}

func (ff *fontFace) glyph(r rune) (uint16, bool) {
	g, ok := ff.f.Cmap[uint32(r)]
	return g, ok && g != 0
}

// advance is the advance of a glyph in em.
func (ff *fontFace) advance(g uint16) float64 { return float64(ff.f.Advance(g)) / ff.upem }

// outline returns the outline of a glyph in font units (y up).
func (ff *fontFace) outline(g uint16) *path {
	if p, ok := ff.glyphs[g]; ok {
		return p
	}
	pp := &pathPen{p: &path{}}
	if !ff.out.Outline(g, pp) {
		pp.p = nil
	}
	ff.glyphs[g] = pp.p
	return pp.p
}

// pathPen builds a path from a glyph outline.
type pathPen struct{ p *path }

func (pp *pathPen) MoveTo(x, y float64)                 { pp.p.moveTo(x, y) }
func (pp *pathPen) LineTo(x, y float64)                 { pp.p.lineTo(x, y) }
func (pp *pathPen) QuadTo(cx, cy, x, y float64)         { pp.p.quadTo(cx, cy, x, y) }
func (pp *pathPen) CubeTo(ax, ay, bx, by, x, y float64) { pp.p.cubicTo(ax, ay, bx, by, x, y) }
func (pp *pathPen) Close()                              { pp.p.close() }

// faceUse is a face with the styles a browser would synthesize for it.
type faceUse struct {
	face         *fontFace
	bold, italic bool
}

// fontChain is the faces a FONT instruction draws with, in the order
// characters are looked up.
type fontChain struct {
	faces []faceUse
	// Resolve system families only when the embedded face lacks a glyph.
	init      func()
	fs        *fonts
	fallbacks []fontdb.Resolved
}

// fonts loads the faces of a document: embedded font parts and, for the
// families fonts are referred to by, the fonts of the system.
type fonts struct {
	r      *Renderer
	db     *fontdb.DB
	parts  map[bdf.Hash]*fontFace // nil: not usable
	files  map[*fontdb.Face]*fontFace
	chains map[bdf.Font]*fontChain
}

func (fs *fonts) fontDB() *fontdb.DB {
	if fs.db == nil {
		o := fs.r.opts
		fs.db = fontdb.New(o.FontFS, o.FontDirs, !o.NoSystemFonts)
	}
	return fs.db
}

func (fs *fonts) part(h bdf.Hash) *fontFace {
	if ff, ok := fs.parts[h]; ok {
		return ff
	}
	fs.parts[h] = nil
	p := fs.r.doc.Part(h)
	if p == nil {
		fs.r.warnf("font %s is missing", h)
		return nil
	}
	data := p.Data
	if woff2.IsWOFF(data) {
		if string(data[:4]) != "wOF2" {
			fs.r.warnf("WOFF 1 fonts are not drawn")
			return nil
		}
		var err error
		if data, err = woff2.Decode(data); err != nil {
			fs.r.warnf("font %s: %v", h, err)
			return nil
		}
	}
	f, err := sfnt.Parse(data)
	if err != nil || f.Cmap == nil {
		fs.r.warnf("font %s cannot be read", h)
		return nil
	}
	ff := newFace(f)
	fs.parts[h] = ff
	return ff
}

func (fs *fonts) file(face *fontdb.Face) *fontFace {
	if ff, ok := fs.files[face]; ok {
		return ff
	}
	fs.files[face] = nil
	l, err := face.Load()
	if err != nil {
		return nil
	}
	ff := newFace(l.Font)
	fs.files[face] = ff
	return ff
}

// genericFamilies maps CSS generic family keywords to the generic families
// of fontdb.
var genericFamilies = map[string]string{
	"sans-serif": fontdb.Sans, "serif": fontdb.Serif, "monospace": fontdb.Mono,
	"cursive": fontdb.Sans, "fantasy": fontdb.Sans, "system-ui": fontdb.Sans, "math": fontdb.Serif,
	"ui-sans-serif": fontdb.Sans, "ui-serif": fontdb.Serif, "ui-monospace": fontdb.Mono, "ui-rounded": fontdb.Sans,
}

// splitFamilies splits a CSS font-family list into names without quotes.
func splitFamilies(list string) []string {
	var out []string
	for _, f := range strings.Split(list, ",") {
		f = strings.TrimSpace(f)
		if len(f) >= 2 && (f[0] == '"' || f[0] == '\'') && f[len(f)-1] == f[0] {
			f = strings.ReplaceAll(f[1:len(f)-1], `\`+string(f[0]), string(f[0]))
		}
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}

// chain resolves a font reference as a browser would: the embedded face
// (registered without weight or style, so those are synthesized), the
// families of its list that are installed, and fallbacks for the rest.
func (fs *fonts) chain(f *bdf.Font) *fontChain {
	if f == nil {
		f = &bdf.Font{Kind: bdf.FontSystem, Family: "sans-serif", Weight: 400}
	}
	if c, ok := fs.chains[*f]; ok {
		return c
	}
	// The caller can reuse its FONT instruction before a missing glyph
	// makes us resolve the remaining families.
	ref := *f
	f = &ref
	c := &fontChain{fs: fs}
	fs.chains[*f] = c
	bold, italic := f.Weight >= 600, f.Style != bdf.StyleNormal
	seen := map[*fontFace]bool{}
	add := func(ff *fontFace, b, i bool) {
		if ff != nil && !seen[ff] {
			seen[ff] = true
			c.faces = append(c.faces, faceUse{ff, b, i})
		}
	}
	if f.Kind == bdf.FontEmbedded {
		add(fs.part(f.Hash), bold, italic)
	}
	c.init = func() {
		generic := ""
		db := fs.fontDB()
		for _, name := range splitFamilies(f.Family) {
			if g, ok := genericFamilies[strings.ToLower(name)]; ok {
				if generic == "" {
					generic = g
				}
				res := db.Resolve(g, bold, italic, false)
				if res.Face != nil {
					add(fs.file(res.Face), res.SynthBold, res.SynthItalic)
				}
				continue
			}
			if faces := db.Family(name); len(faces) > 0 {
				face, sb, si := fontdb.Match(faces, bold, italic)
				add(fs.file(face), sb, si)
			}
		}
		if generic == "" {
			generic = fontdb.Sans
			if len(c.faces) == 0 {
				// a family that is not installed: its look-alike or generic kind
				for _, name := range splitFamilies(f.Family) {
					if res := db.Resolve(name, bold, italic, false); res.Face != nil {
						generic = res.Generic
						add(fs.file(res.Face), res.SynthBold, res.SynthItalic)
						break
					}
				}
			}
		}
		c.fallbacks = db.Fallbacks(generic, bold, italic)
		// Even text with no supported glyph needs the first usable face's
		// missing-glyph advance and baseline metrics.
		if len(c.faces) == 0 {
			for _, res := range c.fallbacks {
				add(fs.file(res.Face), res.SynthBold, res.SynthItalic)
				if len(c.faces) > 0 {
					break
				}
			}
		}
		if len(c.faces) == 0 {
			fs.r.warnf("no font for %q: its text is not drawn", f.Family)
		}
	}
	return c
}

func (c *fontChain) initFallbacks() {
	if c.init != nil {
		init := c.init
		c.init = nil
		init()
	}
}

func (c *fontChain) glyph(r rune) (faceUse, uint16) {
	for _, fu := range c.faces {
		if g, ok := fu.face.glyph(r); ok {
			return fu, g
		}
	}
	if c.init != nil {
		c.initFallbacks()
		return c.glyph(r)
	}
	for _, res := range c.fallbacks {
		if !res.Face.HasRune(r) {
			continue
		}
		if ff := c.fs.file(res.Face); ff != nil {
			if g, ok := ff.glyph(r); ok {
				return faceUse{ff, res.SynthBold, res.SynthItalic}, g
			}
		}
	}
	return c.faces[0], 0
}

// placed is a glyph positioned along the baseline.
type placed struct {
	use faceUse
	gid uint16
	x   float64 // pen position in units
}

// layout places the glyphs of a string at size with letter spacing, in
// visual order, and returns them with the advance of the whole string.
func (c *fontChain) layout(text string, size, spacing float64, dir byte) ([]placed, float64) {
	if len(c.faces) == 0 {
		c.initFallbacks()
		if len(c.faces) == 0 {
			return nil, 0
		}
	}
	runes := visualOrder(text, dir == bdf.DirRTL)
	out := make([]placed, 0, len(runes))
	x := 0.0
	for _, r := range runes {
		switch r {
		case '\t', '\n', '\f', '\r', '\v':
			r = ' ' // Canvas draws these as spaces
		}
		use, gid := c.glyph(r)
		out = append(out, placed{use: use, gid: gid, x: x})
		x += use.face.advance(gid)*size + spacing
	}
	return out, x
}

// visualOrder reorders a string for display: runs of right-to-left
// characters are reversed, and so is the order of the runs in a
// right-to-left paragraph. It covers what previews need, not the whole
// Unicode bidirectional algorithm (no mirroring, no Arabic shaping).
func visualOrder(s string, rtlBase bool) []rune {
	runes := []rune(s)
	if !rtlBase && !hasRTL(s) {
		return runes
	}
	const (
		ltr = iota
		rtl
		neutral
	)
	dirs := make([]int, len(runes))
	for i, r := range runes {
		p, _ := bidi.LookupRune(r)
		switch p.Class() {
		case bidi.R, bidi.AL:
			dirs[i] = rtl
		case bidi.L, bidi.EN, bidi.AN:
			dirs[i] = ltr
		default:
			dirs[i] = neutral
		}
	}
	base := ltr
	if rtlBase {
		base = rtl
	}
	// neutrals take the direction around them when both sides agree
	for i := 0; i < len(dirs); {
		if dirs[i] != neutral {
			i++
			continue
		}
		j := i
		for j < len(dirs) && dirs[j] == neutral {
			j++
		}
		before, after := base, base
		if i > 0 {
			before = dirs[i-1]
		}
		if j < len(dirs) {
			after = dirs[j]
		}
		d := base
		if before == after {
			d = before
		}
		for k := i; k < j; k++ {
			dirs[k] = d
		}
		i = j
	}
	type run struct{ from, to, dir int }
	var runs []run
	for i := 0; i < len(runes); {
		j := i
		for j < len(runes) && dirs[j] == dirs[i] {
			j++
		}
		runs = append(runs, run{i, j, dirs[i]})
		i = j
	}
	if base == rtl {
		for i, j := 0, len(runs)-1; i < j; i, j = i+1, j-1 {
			runs[i], runs[j] = runs[j], runs[i]
		}
	}
	out := make([]rune, 0, len(runes))
	for _, r := range runs {
		if r.dir == rtl {
			for k := r.to - 1; k >= r.from; k-- {
				out = append(out, runes[k])
			}
		} else {
			out = append(out, runes[r.from:r.to]...)
		}
	}
	return out
}

func hasRTL(s string) bool {
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		i += n
		if r < 0x590 {
			continue
		}
		p, _ := bidi.LookupRune(r)
		if c := p.Class(); c == bidi.R || c == bidi.AL {
			return true
		}
	}
	return false
}

// baselineShift is how far below the anchor the alphabetic baseline lies
// for a textBaseline, from the em box of the first face (as Blink does).
func (c *fontChain) baselineShift(baseline byte, size float64) float64 {
	asc, dsc := 0.8, 0.2
	if len(c.faces) > 0 {
		asc, dsc = c.faces[0].face.asc, c.faces[0].face.dsc
	}
	emAsc := size * asc / (asc + dsc)
	emDsc := size - emAsc
	switch baseline {
	case bdf.BaselineTop:
		return emAsc
	case bdf.BaselineHanging:
		return asc * size * 0.8
	case bdf.BaselineMiddle:
		return (emAsc - emDsc) / 2
	case bdf.BaselineIdeographic:
		return -dsc * size
	case bdf.BaselineBottom:
		return -emDsc
	}
	return 0
}

// fakeBoldWidth is the outline growth of synthesized bold for a size in
// device pixels (Skia's interpolation between 1/24 at 9 px and 1/32 at 36 px).
func fakeBoldWidth(size float64) float64 {
	switch {
	case size <= 9:
		return size / 24
	case size >= 36:
		return size / 32
	}
	k := 1.0/24 + (size-9)/(36-9)*(1.0/32-1.0/24)
	return size * k
}

// fakeItalicSkew is the slant of synthesized italic (x per unit of height).
const fakeItalicSkew = 0.25

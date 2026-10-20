package psbdf

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/internal/textrun"
)

// This file is vector/pdfbdf's text.go, writing Type 3 fonts: the fonts
// of a document are found and text is laid out as pdfbdf does it, and the
// glyphs in use become procedures of their outlines.

// face is a font program text is drawn with, as Type 3 fonts of the
// outlines of its glyphs in use, of 256 glyphs each.
type face struct {
	f        *sfnt.Font
	upem     float64
	asc, dsc float64 // hhea ascent and descent in em, both positive
	out      *sfnt.Outlines

	// the glyphs in use: their font among subs and their code in it
	codes map[uint16]code
	subs  []*subfont
}

// code is where a glyph is: its font among the subfonts of its face and
// its code in that font.
type code struct {
	sub int
	c   byte
}

// subfont is a Type 3 font of up to 256 glyphs of a face, in the order of
// their codes.
type subfont struct {
	name   string
	glyphs []uint16
}

func newFace(f *sfnt.Font) *face {
	fc := &face{f: f, upem: float64(f.UnitsPerEm), out: sfnt.NewOutlines(f), codes: map[uint16]code{}}
	if fc.upem <= 0 {
		fc.upem = 1000
	}
	fc.asc, fc.dsc = float64(f.Ascent)/fc.upem, -float64(f.Descent)/fc.upem
	if fc.asc+fc.dsc <= 0 {
		fc.asc, fc.dsc = 0.8, 0.2
	}
	return fc
}

func (fc *face) glyph(r rune) (uint16, bool) {
	g, ok := fc.f.Cmap[uint32(r)]
	return g, ok && g != 0
}

// advance is the advance of a glyph in em.
func (fc *face) advance(g uint16) float64 { return float64(fc.f.Advance(g)) / fc.upem }

// use returns the code of a glyph, giving it one when it has none.
func (fc *face) use(fs *fonts, g uint16) code {
	if c, ok := fc.codes[g]; ok {
		return c
	}
	n := len(fc.subs)
	if n == 0 || len(fc.subs[n-1].glyphs) == 256 {
		fc.subs = append(fc.subs, &subfont{name: fmt.Sprintf("BDF%d", fs.made)})
		fs.made++
		n++
	}
	sf := fc.subs[n-1]
	c := code{n - 1, byte(len(sf.glyphs))}
	sf.glyphs = append(sf.glyphs, g)
	fc.codes[g] = c
	return c
}

// faceUse is a face with the styles a browser would synthesize for it.
type faceUse struct {
	face         *face
	bold, italic bool
}

// fontChain is the faces a FONT instruction draws with, in the order
// characters are looked up.
type fontChain struct {
	faces []faceUse
	// the families of the system are resolved when the embedded face
	// lacks a glyph
	init      func()
	fs        *fonts
	fallbacks []fontdb.Resolved
}

// fonts loads the faces of a document: embedded font parts and, for the
// families fonts are referred to by, the fonts of the system.
type fonts struct {
	ex     *exporter
	db     *fontdb.DB
	parts  map[bdf.Hash]*face // nil: not usable
	files  map[*fontdb.Face]*face
	chains map[bdf.Font]*fontChain
	all    []*face
	// made counts the Type 3 fonts
	made int
}

func newFonts(ex *exporter) *fonts {
	return &fonts{ex: ex, parts: map[bdf.Hash]*face{}, files: map[*fontdb.Face]*face{}, chains: map[bdf.Font]*fontChain{}}
}

// count is the number of Type 3 fonts made.
func (fs *fonts) count() int { return fs.made }

func (fs *fonts) fontDB() *fontdb.DB {
	if fs.db == nil {
		o := fs.ex.opts
		fs.db = fontdb.New(o.FontFS, o.FontDirs, !o.NoSystemFonts)
	}
	return fs.db
}

func (fs *fonts) part(h bdf.Hash) *face {
	if fc, ok := fs.parts[h]; ok {
		return fc
	}
	fs.parts[h] = nil
	p := fs.ex.doc.Part(h)
	if p == nil {
		fs.ex.warnf("font %s is missing", h)
		return nil
	}
	data := p.Data
	if woff2.IsWOFF(data) {
		if string(data[:4]) != "wOF2" {
			fs.ex.warnf("WOFF 1 fonts are not written")
			return nil
		}
		var err error
		if data, err = woff2.Decode(data); err != nil {
			fs.ex.warnf("font %s: %v", h, err)
			return nil
		}
	}
	f, err := sfnt.Parse(data)
	if err != nil || f.Cmap == nil {
		fs.ex.warnf("font %s cannot be read", h)
		return nil
	}
	fc := newFace(f)
	fs.parts[h] = fc
	fs.all = append(fs.all, fc)
	return fc
}

func (fs *fonts) file(face *fontdb.Face) *face {
	if fc, ok := fs.files[face]; ok {
		return fc
	}
	fs.files[face] = nil
	l, err := face.Load()
	if err != nil {
		return nil
	}
	fc := newFace(l.Font)
	fs.files[face] = fc
	fs.all = append(fs.all, fc)
	return fc
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
	ref := *f
	f = &ref
	c := &fontChain{fs: fs}
	fs.chains[*f] = c
	bold, italic := f.Weight >= 600, f.Style != bdf.StyleNormal
	seen := map[*face]bool{}
	add := func(fc *face, b, i bool) {
		if fc != nil && !seen[fc] {
			seen[fc] = true
			c.faces = append(c.faces, faceUse{fc, b, i})
		}
	}
	if f.Kind == bdf.FontEmbedded {
		add(fs.part(f.Hash), bold, italic)
	}
	c.init = func() {
		generic := ""
		db := fs.fontDB()
		for _, name := range textrun.SplitFamilies(f.Family) {
			if g, ok := textrun.GenericFamilies[strings.ToLower(name)]; ok {
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
				fa, sb, si := fontdb.Match(faces, bold, italic)
				add(fs.file(fa), sb, si)
			}
		}
		if generic == "" {
			generic = fontdb.Sans
			if len(c.faces) == 0 {
				// a family that is not installed: its look-alike or generic kind
				for _, name := range textrun.SplitFamilies(f.Family) {
					if res := db.Resolve(name, bold, italic, false); res.Face != nil {
						generic = res.Generic
						add(fs.file(res.Face), res.SynthBold, res.SynthItalic)
						break
					}
				}
			}
		}
		c.fallbacks = db.Fallbacks(generic, bold, italic)
		// even text with no supported glyph needs the first usable face's
		// missing-glyph advance and baseline metrics
		if len(c.faces) == 0 {
			for _, res := range c.fallbacks {
				add(fs.file(res.Face), res.SynthBold, res.SynthItalic)
				if len(c.faces) > 0 {
					break
				}
			}
		}
		if len(c.faces) == 0 {
			fs.ex.warnf("no font for %q: its text is not written", f.Family)
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
		if fc := c.fs.file(res.Face); fc != nil {
			if g, ok := fc.glyph(r); ok {
				return faceUse{fc, res.SynthBold, res.SynthItalic}, g
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
func (c *fontChain) layout(str string, size, spacing float64, dir byte) ([]placed, float64) {
	if len(c.faces) == 0 {
		c.initFallbacks()
		if len(c.faces) == 0 {
			return nil, 0
		}
	}
	runes := textrun.VisualOrder(str, dir == bdf.DirRTL)
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

func (c *fontChain) baselineShift(baseline byte, size float64) float64 {
	asc, dsc := 0.8, 0.2
	if len(c.faces) > 0 {
		asc, dsc = c.faces[0].face.asc, c.faces[0].face.dsc
	}
	return textrun.BaselineShift(baseline, size, asc, dsc)
}

// text writes a FILL_TEXT or STROKE_TEXT: filled text of a color with the
// Type 3 fonts, other text (strokes, gradients, synthesized bold) as the
// paths of the outlines of its glyphs.
func (ex *exporter) text(str string, x, y, advance float64, stroke bool) {
	st := &ex.st
	if !(st.size > 0) || math.IsInf(st.size, 0) || str == "" {
		return
	}
	ch := ex.fonts.chain(st.font)
	glyphs, width := ch.layout(str, st.size, st.spacing, st.dir)
	if len(glyphs) == 0 {
		return
	}
	// the string is as wide as the document says it is
	hs := 1.0
	if advance > 0 && width > 0 && math.Abs(width-advance)/advance > 0.005 {
		hs = advance / width
	}
	dx := 0.0
	switch st.align {
	case bdf.AlignRight:
		dx = -width
	case bdf.AlignCenter:
		dx = -width / 2
	case bdf.AlignEnd:
		if st.dir != bdf.DirRTL {
			dx = -width
		}
	case bdf.AlignStart:
		if st.dir == bdf.DirRTL {
			dx = -width
		}
	}
	by := y + ch.baselineShift(st.baseline, st.size)
	sty := &st.fill
	if stroke {
		sty = &st.stroke
	}
	visible, plain := ex.look(sty)
	if !visible {
		return
	}
	// around the lines of the font, and the slant of a synthesized italic
	left, right := x+dx*hs, x+(dx+width)*hs
	reach := 0.3 * st.size
	if stroke {
		reach += st.reach()
	}
	if !ex.drawn([4]float64{math.Min(left, right), by - 1.2*st.size, math.Max(left, right), by + 0.5*st.size}, reach, plain) {
		return
	}
	// the runs of glyphs of the same face
	type run struct {
		use    faceUse
		glyphs []placed
	}
	var runs []run
	for _, g := range glyphs {
		if n := len(runs); n > 0 && runs[n-1].use == g.use {
			runs[n-1].glyphs = append(runs[n-1].glyphs, g)
			continue
		}
		runs = append(runs, run{use: g.use, glyphs: []placed{g}})
	}
	for _, r := range runs {
		skew := 0.0
		if r.use.italic {
			skew = textrun.FakeItalicSkew
		}
		build := func(b *bytes.Buffer) bool {
			return ex.glyphPaths(b, r.use.face, r.glyphs, x, by, dx, hs, skew)
		}
		switch {
		case stroke:
			ex.paintWith(st.stroke, build, "S", "strokepath W")
		case r.use.bold:
			// synthesized bold: the fill, and a line around it
			ex.paintWith(st.fill, build, "f", "W")
			saved := st.pColor
			ex.op("q %s w 1 j [] 0 d", num(textrun.FakeBoldWidth(st.size)))
			ex.paintWith(st.fill, build, "S", "strokepath W")
			ex.op("Q")
			st.pColor = saved
		case st.fill.paint == nil || ex.shading(st.fill.paint).solid:
			ex.glyphText(r.use.face, r.glyphs, x, by, dx, hs, skew)
		default:
			ex.paintWith(st.fill, build, "f", "W")
		}
	}
}

// glyphText writes glyphs of a face with its Type 3 fonts, at the places
// of the layout.
func (ex *exporter) glyphText(fc *face, glyphs []placed, x, y, dx, hs, skew float64) {
	st := &ex.st
	if st.fill.paint == nil {
		ex.color(st.fill.color)
	} else {
		ex.color(ex.shading(st.fill.paint).color)
	}
	size := st.size
	// the glyph space has its y axis up
	fm := matrix{size * hs, 0, size * hs * skew, -size, 0, 0}
	for i := 0; i < len(glyphs); {
		sub := fc.use(ex.fonts, glyphs[i].gid).sub
		j := i
		var hex, dxs strings.Builder
		for ; j < len(glyphs); j++ {
			c := fc.use(ex.fonts, glyphs[j].gid)
			if c.sub != sub {
				break
			}
			if j > i && (j-i)%16 == 0 {
				// (lines of the file are short)
				hex.WriteByte('\n')
				dxs.WriteByte('\n')
			} else if j > i {
				dxs.WriteByte(' ')
			}
			fmt.Fprintf(&hex, "%02X", c.c)
			if j+1 < len(glyphs) {
				dxs.WriteString(num((glyphs[j+1].x - glyphs[j].x) * hs))
			} else {
				dxs.WriteString(num(fc.advance(glyphs[j].gid) * size * hs))
			}
		}
		ex.op("/%s [%s] Tf %s m <%s>\n[%s] xshow n", fc.subs[sub].name, fm, nums(x+(dx+glyphs[i].x)*hs, y), hex.String(), dxs.String())
		i = j
	}
}

// glyphPaths writes glyphs of a face as the paths of their outlines, and
// reports whether they have segments.
func (ex *exporter) glyphPaths(b *bytes.Buffer, fc *face, glyphs []placed, x, y, dx, hs, skew float64) bool {
	k := ex.st.size / fc.upem
	drawn := false
	for _, g := range glyphs {
		m := matrix{hs, 0, 0, 1, x + (dx+g.x)*hs, y}.mul(matrix{k, 0, 0, -k, 0, 0}).mul(matrix{1, 0, skew, 1, 0, 0})
		p := newPen(b, m)
		fc.out.Outline(g.gid, outlinePen{p})
		drawn = drawn || p.drawn
	}
	return drawn
}

// define writes the Type 3 fonts of the glyphs in use, deflated.
func (fs *fonts) define(b *bytes.Buffer) {
	var d bytes.Buffer
	for _, fc := range fs.all {
		for _, sf := range fc.subs {
			fc.define(&d, sf)
		}
	}
	if d.Len() == 0 {
		return
	}
	b.WriteString("%%BeginResource: procset bdf-fonts 1 0\nbdfrun\n")
	writeA85(b, deflate(d.Bytes()))
	b.WriteString("%%EndResource\n")
}

// define writes a Type 3 font of glyphs of a face: a procedure of the
// outline of each, in the units of the font.
func (fc *face) define(b *bytes.Buffer, sf *subfont) {
	bbox := none
	var procs bytes.Buffer
	for i, g := range sf.glyphs {
		var path bytes.Buffer
		p := newPen(&path, identity)
		fc.out.Outline(g, outlinePen{p})
		fmt.Fprintf(&procs, "/g%d {%d 0 ", i, fc.f.Advance(g))
		if p.drawn && p.boxed {
			bbox = union(bbox, p.box)
			fmt.Fprintf(&procs, "%s setcachedevice\n%sf}\n", nums(p.box[0], p.box[1], p.box[2], p.box[3]), path.Bytes())
		} else {
			procs.WriteString("0 0 0 0 setcachedevice}\n")
		}
	}
	if !(bbox[2] >= bbox[0]) {
		bbox = [4]float64{0, 0, 0, 0}
	}
	upem := num(fc.upem)
	fmt.Fprintf(b, "/%s <<\n/FontType 3 /FontMatrix [1 %s div 0 0 1 %s div 0 0] /FontBBox [%s]\n/Encoding [", sf.name, upem, upem,
		nums(math.Floor(bbox[0]), math.Floor(bbox[1]), math.Ceil(bbox[2]), math.Ceil(bbox[3])))
	for i := 0; i < 256; i++ {
		if i%16 == 0 {
			b.WriteByte('\n')
		}
		if i < len(sf.glyphs) {
			fmt.Fprintf(b, "/g%d ", i)
		} else {
			b.WriteString("/.notdef ")
		}
	}
	b.WriteString("]\n/CharProcs <<\n/.notdef {0 0 0 0 0 0 setcachedevice}\n")
	b.Write(procs.Bytes())
	b.WriteString(">>\n/BuildGlyph {exch /CharProcs get exch 2 copy known not {pop /.notdef} if get bdf begin exec end} bind\n" +
		"/BuildChar {1 index /Encoding get exch get 1 index /BuildGlyph get exec} bind\n>> definefont pop\n")
}

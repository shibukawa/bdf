package pdfbdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"maps"
	"math"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/cff"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/internal/textrun"
)

// face is a font program text is drawn with: embedded in the file as a
// font of its glyphs in use, or drawn as the outlines of its glyphs when
// it has no program a PDF holds or may not be embedded.
type face struct {
	f        *sfnt.Font
	upem     float64
	asc, dsc float64 // hhea ascent and descent in em, both positive
	// outlines is set for a face that is drawn as paths
	outlines bool
	out      *sfnt.Outlines
	// subset says that the glyphs that are not in use may be left out
	subset bool
	// cids are the CIDs of the glyphs of a CID-keyed CFF program, which
	// the codes of its text are (nil: the codes are the glyph numbers)
	cids []int

	// in the file: its name as a resource, its font dictionary, and the
	// glyphs in use with the characters they stand for
	name string
	ref  ref
	used map[uint16]rune
}

func (ex *exporter) newFace(f *sfnt.Font, system bool) *face {
	fc := &face{f: f, upem: float64(f.UnitsPerEm), subset: true, used: map[uint16]rune{}}
	if fc.upem <= 0 {
		fc.upem = 1000
	}
	fc.asc, fc.dsc = float64(f.Ascent)/fc.upem, -float64(f.Descent)/fc.upem
	if fc.asc+fc.dsc <= 0 {
		fc.asc, fc.dsc = 0.8, 0.2
	}
	switch {
	case !f.IsCFF && f.Tables["glyf"] != nil && f.Tables["loca"] != nil:
	case f.IsCFF && f.Tables["CFF "] != nil:
		cf, err := cff.Parse(f.Tables["CFF "])
		if err != nil {
			fc.outlines = true
		} else if cf.IsCID {
			fc.cids = cf.Charset
		}
	default:
		fc.outlines = true // a CFF2 program, or none
	}
	if o, ok := f.OS2(); ok && !ex.opts.IgnoreFSType {
		embed, subset := o.Embeddable()
		// (a font of the document was embedded in it already: only a
		// font of the system is asked whether it may be)
		if system && !embed {
			fc.outlines = true
		}
		fc.subset = subset
	}
	if fc.outlines {
		fc.out = sfnt.NewOutlines(f)
	}
	return fc
}

func (fc *face) glyph(r rune) (uint16, bool) {
	g, ok := fc.f.Cmap[uint32(r)]
	return g, ok && g != 0
}

// advance is the advance of a glyph in em.
func (fc *face) advance(g uint16) float64 { return float64(fc.f.Advance(g)) / fc.upem }

// code is the code of a glyph in the text of the file.
func (fc *face) code(g uint16) uint16 {
	if fc.cids != nil && int(g) < len(fc.cids) {
		return uint16(fc.cids[g])
	}
	return g
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
}

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
	fc := fs.ex.newFace(f, false)
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
	fc := fs.ex.newFace(l.Font, true)
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
	r   rune
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
		out = append(out, placed{use: use, gid: gid, r: r, x: x})
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

// text writes a FILL_TEXT or STROKE_TEXT: text of the embedded fonts where
// the faces can be embedded, the outlines of the glyphs where not. alt is
// the text that the string stands for, when a MARK told one.
func (ex *exporter) text(str string, x, y, advance float64, stroke bool, alt *string) {
	s := ex.s
	st := &s.st
	if st.size <= 0 || str == "" {
		ex.marked(alt, func() {})
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
	// the runs of glyphs of the same face
	type run struct {
		use    faceUse
		glyphs []placed
	}
	var runs []run
	paths, gap := false, false
	for _, g := range glyphs {
		if g.gid == 0 && ex.pdfa != nil {
			// a PDF/A file has no glyph that stands for a character its
			// font lacks: the character is left out (the glyphs after it
			// are where they would be)
			ex.warnf("characters that no font has are left out of the PDF/A file")
			gap = true
			continue
		}
		if n := len(runs); n > 0 && runs[n-1].use == g.use && !gap {
			runs[n-1].glyphs = append(runs[n-1].glyphs, g)
			continue
		}
		gap = false
		runs = append(runs, run{use: g.use, glyphs: []placed{g}})
		paths = paths || g.use.face.outlines
	}
	if len(runs) == 0 {
		return
	}
	if st.shadowOn() {
		// (around the lines of the font, and the slant of a synthesized italic)
		left, right := x+dx*hs, x+(dx+width)*hs
		reach := 0.3 * st.size
		if stroke {
			reach += st.reach()
		}
		ex.castShadow([4]float64{math.Min(left, right), by - 1.2*st.size, math.Max(left, right), by + 0.5*st.size}, reach, func(o *bdf.Object) {
			if stroke {
				o.StrokeText(str, float32(x), float32(y), float32(advance))
			} else {
				o.FillText(str, float32(x), float32(y), float32(advance))
			}
		})
	}
	// what is searched and copied: the string as it is read, when the
	// glyphs are not in its order or are not text
	if alt == nil && (paths || textrun.HasRTL(str) || st.dir == bdf.DirRTL) {
		alt = &str
	}
	if a := ex.pdfa; a != nil && a.conf == "A" && alt == nil && strings.ContainsFunc(str, privateUse) {
		// (level a asks for the text that such characters stand for,
		// which the document does not tell)
		ex.warnf("text in characters of the private use areas has no replacement text: the file is not a PDF/A-%da file for it (it is one of level u)", a.part)
	}
	ex.marked(alt, func() {
		ex.withPrep(stroke, func(p prepared) {
			if p.clipImage != nil {
				ex.warnf("text in a pattern that does not repeat is not drawn")
				return
			}
			for _, r := range runs {
				fc := r.use.face
				bold := 0.0
				if r.use.bold && !stroke {
					bold = textrun.FakeBoldWidth(st.size)
				}
				skew := 0.0
				if r.use.italic {
					skew = textrun.FakeItalicSkew
				}
				draw := func() {
					if fc.outlines {
						ex.glyphPaths(fc, r.glyphs, x, by, dx, hs, skew, stroke, bold > 0)
						return
					}
					ex.glyphText(fc, r.glyphs, x, by, dx, hs, skew, stroke, bold > 0)
				}
				if bold == 0 {
					draw()
					continue
				}
				// synthesized bold: the fill, and a line around it
				ex.scoped(func() {
					ex.strokeLikeFill()
					s.op("%s w 1 j [] 0 d", num(bold))
					draw()
				})
			}
		})
	})
}

// privateUse reports whether a character is of a private use area: one
// that only its font gives a meaning.
func privateUse(r rune) bool {
	return r >= 0xe000 && r <= 0xf8ff || r >= 0xf0000 && r <= 0xffffd || r >= 0x100000 && r <= 0x10fffd
}

// strokeLikeFill makes the stroke of the PDF graphics state what its fill
// is: for the line of synthesized bold.
func (ex *exporter) strokeLikeFill() {
	s := ex.s
	st := &s.st
	if strings.HasPrefix(st.pFill, "p") {
		s.op("/Pattern CS /%s SCN", st.pFill[1:])
	} else if st.fill.paint == nil {
		s.op("%s RG", rgbOf(st.fill.color))
	} else if sh := ex.shading(st.fill.paint); sh.solid {
		s.op("%s RG", rgbOf(sh.color))
	}
	if st.pCA != st.pCa {
		s.op("%s gs", ex.extGState("/CA "+fine(st.pCa)))
	}
}

// glyphText writes glyphs of a face as a text object.
func (ex *exporter) glyphText(fc *face, glyphs []placed, x, y, dx, hs, skew float64, stroke, both bool) {
	s := ex.s
	st := &s.st
	if fc.name == "" {
		fc.name, fc.ref = ex.newName("F"), ex.pw.alloc()
	}
	var hex strings.Builder
	for _, g := range glyphs {
		if _, ok := fc.used[g.gid]; !ok {
			fc.used[g.gid] = g.r
		}
		fmt.Fprintf(&hex, "%04X", fc.code(g.gid))
	}
	mode := ""
	switch {
	case both:
		mode = "2 Tr "
	case stroke:
		mode = "1 Tr "
	}
	body := func() {
		tc := ""
		if st.spacing != st.pTc {
			tc = num(st.spacing) + " Tc "
			st.pTc = st.spacing
		}
		// the text space has its y axis up
		tm := matrix{hs, 0, hs * skew, -1, x + (dx+glyphs[0].x)*hs, y}
		s.op("BT %s %s Tf %s%s%s Tm <%s> Tj ET", s.use("Font", fc.name, fc.ref), num(st.size), tc, mode, tm, hex.String())
	}
	if mode == "" || both {
		// (the text rendering mode of synthesized bold is in its q and Q)
		body()
		return
	}
	ex.scoped(body)
}

// glyphPaths writes glyphs of a face as the paths of their outlines.
func (ex *exporter) glyphPaths(fc *face, glyphs []placed, x, y, dx, hs, skew float64, stroke, both bool) {
	s := ex.s
	st := &s.st
	k := st.size / fc.upem
	drawn := false
	for _, g := range glyphs {
		m := matrix{hs, 0, 0, 1, x + (dx+g.x)*hs, y}.mul(matrix{k, 0, 0, -k, 0, 0}).mul(matrix{1, 0, skew, 1, 0, 0})
		p := &pen{b: &s.b, m: m}
		fc.out.Outline(g.gid, outlinePen{p})
		drawn = drawn || p.drawn
	}
	if !drawn {
		return
	}
	switch {
	case both:
		s.op("B")
	case stroke:
		s.op("S")
	default:
		s.op("f")
	}
}

// fontTables are the tables of a font program that a PDF reads: of a
// TrueType program, and of an OpenType program with CFF outlines.
var (
	glyfTables = []string{"head", "hhea", "maxp", "hmtx", "loca", "glyf", "cvt ", "fpgm", "prep"}
	cffTables  = []string{"CFF ", "cmap", "head", "hhea", "maxp", "hmtx", "OS/2", "name", "post"}
)

// writeFonts writes the fonts in use: each with the glyphs that are drawn.
func (ex *exporter) writeFonts() int {
	n := 0
	for _, fc := range ex.fonts.all {
		if fc.name == "" || len(fc.used) == 0 {
			continue
		}
		n++
		ex.writeFont(fc)
	}
	return n
}

func (ex *exporter) writeFont(fc *face) {
	pw := ex.pw
	gids := make([]uint16, 0, len(fc.used))
	keep := map[uint16]bool{0: true}
	for g := range fc.used {
		gids = append(gids, g)
		keep[g] = true
	}
	sort.Slice(gids, func(i, j int) bool { return gids[i] < gids[j] })
	// the program with the glyphs in use: the font of the document is not
	// changed
	f := *fc.f
	f.Tables = maps.Clone(fc.f.Tables)
	names := glyfTables
	if f.IsCFF {
		names = cffTables
		if fc.subset {
			f.PruneCFF(keep)
		}
	} else if fc.subset {
		f.PruneGlyphs(keep)
	}
	tables := map[string][]byte{}
	for _, t := range names {
		if b := f.Tables[t]; b != nil {
			tables[t] = b
		}
	}
	program := sfnt.Build(tables, f.IsCFF)

	// a name of its own for a part of a font
	base := psName(fc.f)
	if fc.subset {
		tag := []byte("AAAAAA")
		h := uint32(2166136261)
		for _, g := range gids {
			h = (h ^ uint32(g)) * 16777619
		}
		h ^= uint32(len(program))
		for i := range tag {
			tag[i] += byte(h % 26)
			h /= 26
		}
		base = string(tag) + "+" + base
	}
	k := 1000 / fc.upem
	bbox := [4]float64{0, -200, 1000, 900}
	if head := fc.f.Tables["head"]; len(head) >= 44 {
		for i := range bbox {
			bbox[i] = float64(int16(binary.BigEndian.Uint16(head[36+2*i:]))) * k
		}
	}
	var file string
	if f.IsCFF {
		file = "/FontFile3 " + pw.addStream(" /Subtype /OpenType", program, true).String()
	} else {
		file = "/FontFile2 " + pw.addStream(fmt.Sprintf(" /Length1 %d", len(program)), program, true).String()
	}
	desc := pw.add(fmt.Sprintf("<< /Type /FontDescriptor /FontName %s /Flags 4 /FontBBox [%s] /ItalicAngle 0 /Ascent %s /Descent %s /CapHeight %s /StemV 80 %s >>",
		name(base), nums(bbox[0], bbox[1], bbox[2], bbox[3]), num(fc.asc*1000), num(-fc.dsc*1000), num(fc.asc*700), file))

	// the widths of the glyphs in use, by their codes
	type cw struct {
		code  uint16
		width float64
	}
	widths := make([]cw, 0, len(gids))
	seen := map[uint16]bool{}
	for _, g := range gids {
		if c := fc.code(g); !seen[c] {
			seen[c] = true
			widths = append(widths, cw{c, float64(fc.f.Advance(g)) * k})
		}
	}
	sort.Slice(widths, func(i, j int) bool { return widths[i].code < widths[j].code })
	var w strings.Builder
	for i := 0; i < len(widths); {
		j := i
		fmt.Fprintf(&w, "%d [", widths[i].code)
		for ; j < len(widths) && int(widths[j].code) == int(widths[i].code)+j-i; j++ {
			w.WriteString(fine(widths[j].width) + " ")
		}
		w.WriteString("] ")
		i = j
	}
	subtype, toGID := "/CIDFontType2", " /CIDToGIDMap /Identity"
	if f.IsCFF {
		subtype, toGID = "/CIDFontType0", ""
	}
	cid := pw.add(fmt.Sprintf("<< /Type /Font /Subtype %s /BaseFont %s /CIDSystemInfo << /Registry (Adobe) /Ordering (Identity) /Supplement 0 >> /FontDescriptor %s /DW 1000 /W [%s]%s >>",
		subtype, name(base), desc, w.String(), toGID))

	// the characters of the codes, for searching and copying
	var cmap bytes.Buffer
	cmap.WriteString("/CIDInit /ProcSet findresource begin\n12 dict begin\nbegincmap\n/CIDSystemInfo << /Registry (Adobe) /Ordering (UCS) /Supplement 0 >> def\n" +
		"/CMapName /Adobe-Identity-UCS def\n/CMapType 2 def\n1 begincodespacerange\n<0000> <FFFF>\nendcodespacerange\n")
	type cu struct {
		code uint16
		r    rune
	}
	chars := make([]cu, 0, len(gids))
	seen = map[uint16]bool{}
	for _, g := range gids {
		if c := fc.code(g); !seen[c] && g != 0 {
			seen[c] = true
			chars = append(chars, cu{c, fc.used[g]})
		}
	}
	sort.Slice(chars, func(i, j int) bool { return chars[i].code < chars[j].code })
	for i := 0; i < len(chars); i += 100 {
		part := chars[i:min(i+100, len(chars))]
		fmt.Fprintf(&cmap, "%d beginbfchar\n", len(part))
		for _, c := range part {
			fmt.Fprintf(&cmap, "<%04X> <", c.code)
			for _, u := range utf16.Encode([]rune{c.r}) {
				fmt.Fprintf(&cmap, "%04X", u)
			}
			cmap.WriteString(">\n")
		}
		cmap.WriteString("endbfchar\n")
	}
	cmap.WriteString("endcmap\nCMapName currentdict /CMap defineresource pop\nend\nend\n")
	toUnicode := pw.addStream("", cmap.Bytes(), true)
	pw.put(fc.ref, fmt.Sprintf("<< /Type /Font /Subtype /Type0 /BaseFont %s /Encoding /Identity-H /DescendantFonts [%s] /ToUnicode %s >>", name(base), cid, toUnicode))
}

// psName returns the PostScript name of a font program, or one made of
// its family name.
func psName(f *sfnt.Font) string {
	best, rank := "", 0
	t := f.Tables["name"]
	if len(t) >= 6 {
		count, strs := int(binary.BigEndian.Uint16(t[2:])), int(binary.BigEndian.Uint16(t[4:]))
		for i := 0; i < count && 6+12*(i+1) <= len(t); i++ {
			rec := t[6+12*i:]
			platform, id := binary.BigEndian.Uint16(rec), binary.BigEndian.Uint16(rec[6:])
			n, off := int(binary.BigEndian.Uint16(rec[8:])), strs+int(binary.BigEndian.Uint16(rec[10:]))
			r := map[uint16]int{6: 3, 4: 2, 1: 1}[id]
			if r <= rank || off+n > len(t) || n == 0 {
				continue
			}
			raw := t[off : off+n]
			var s string
			if platform == 1 {
				s = string(raw)
			} else {
				u := make([]uint16, n/2)
				for k := range u {
					u[k] = binary.BigEndian.Uint16(raw[2*k:])
				}
				s = string(utf16.Decode(u))
			}
			if s = sfnt.SanitizePS(s); s != "" {
				best, rank = s, r
			}
		}
	}
	if best == "" {
		return "Font"
	}
	return best
}

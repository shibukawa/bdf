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
	"github.com/shibukawa/bdf/internal/cff"
	"github.com/shibukawa/bdf/internal/fontchain"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/internal/textrun"
	"github.com/shibukawa/bdf/vector/internal/export"
)

// face is a font program text is drawn with: embedded in the file as a
// font of its glyphs in use, or drawn as the outlines of its glyphs when
// it has no program a PDF holds or may not be embedded.
type face struct {
	fontchain.Program
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
	fc := &face{Program: fontchain.NewProgram(f), subset: true, used: map[uint16]rune{}}
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

// code is the code of a glyph in the text of the file.
func (fc *face) code(g uint16) uint16 {
	if fc.cids != nil && int(g) < len(fc.cids) {
		return uint16(fc.cids[g])
	}
	return g
}

// The faces of a document and the chains of a font reference are found by
// internal/fontchain, which lays text out on them too.
type (
	fonts   = fontchain.Fonts[*face]
	faceUse = fontchain.Use[*face]
	placed  = fontchain.Placed[*face]
)

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
	ch := ex.fonts.Chain(st.font)
	glyphs, width := ch.Layout(str, st.size, st.spacing, st.dir)
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
	by := y + ch.BaselineShift(st.baseline, st.size)
	// the runs of glyphs of the same face
	type run struct {
		use    faceUse
		glyphs []placed
	}
	var runs []run
	paths, gap := false, false
	for _, g := range glyphs {
		if g.GID == 0 && ex.pdfa != nil {
			// a PDF/A file has no glyph that stands for a character its
			// font lacks: the character is left out (the glyphs after it
			// are where they would be)
			ex.Warnf("characters that no font has are left out of the PDF/A file")
			gap = true
			continue
		}
		if n := len(runs); n > 0 && runs[n-1].use == g.Use && !gap {
			runs[n-1].glyphs = append(runs[n-1].glyphs, g)
			continue
		}
		gap = false
		runs = append(runs, run{use: g.Use, glyphs: []placed{g}})
		paths = paths || g.Use.Face.outlines
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
		ex.Warnf("text in characters of the private use areas has no replacement text: the file is not a PDF/A-%da file for it (it is one of level u)", a.part)
	}
	ex.marked(alt, func() {
		ex.withPrep(stroke, func(p prepared) {
			if p.clipImage != nil {
				ex.Warnf("text in a pattern that does not repeat is not drawn")
				return
			}
			for _, r := range runs {
				fc := r.use.Face
				bold := 0.0
				if r.use.Bold && !stroke {
					bold = textrun.FakeBoldWidth(st.size)
				}
				skew := 0.0
				if r.use.Italic {
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
					s.op("%s w 1 j [] 0 d", export.Num(bold))
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
		s.op("%s gs", ex.extGState("/CA "+export.Fine(st.pCa)))
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
		if _, ok := fc.used[g.GID]; !ok {
			fc.used[g.GID] = g.R
		}
		fmt.Fprintf(&hex, "%04X", fc.code(g.GID))
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
			tc = export.Num(st.spacing) + " Tc "
			st.pTc = st.spacing
		}
		// the text space has its y axis up
		tm := matrix{hs, 0, hs * skew, -1, x + (dx+glyphs[0].X)*hs, y}
		s.op("BT %s %s Tf %s%s%s Tm <%s> Tj ET", s.use("Font", fc.name, fc.ref), export.Num(st.size), tc, mode, export.Mat(tm), hex.String())
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
	k := st.size / fc.Upem
	drawn := false
	for _, g := range glyphs {
		m := matrix{hs, 0, 0, 1, x + (dx+g.X)*hs, y}.Mul(matrix{k, 0, 0, -k, 0, 0}).Mul(matrix{1, 0, skew, 1, 0, 0})
		p := export.NewPen(&s.b, m)
		fc.out.Outline(g.GID, p)
		drawn = drawn || p.Drawn
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
	for _, fc := range ex.fonts.All {
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
	f := *fc.Font
	f.Tables = maps.Clone(fc.Font.Tables)
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
	base := psName(fc.Font)
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
	k := 1000 / fc.Upem
	bbox := [4]float64{0, -200, 1000, 900}
	if head := fc.Font.Tables["head"]; len(head) >= 44 {
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
		name(base), export.Nums(bbox[0], bbox[1], bbox[2], bbox[3]), export.Num(fc.Asc*1000), export.Num(-fc.Dsc*1000), export.Num(fc.Asc*700), file))

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
			widths = append(widths, cw{c, float64(fc.Font.Advance(g)) * k})
		}
	}
	sort.Slice(widths, func(i, j int) bool { return widths[i].code < widths[j].code })
	var w strings.Builder
	for i := 0; i < len(widths); {
		j := i
		fmt.Fprintf(&w, "%d [", widths[i].code)
		for ; j < len(widths) && int(widths[j].code) == int(widths[i].code)+j-i; j++ {
			w.WriteString(export.Fine(widths[j].width) + " ")
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

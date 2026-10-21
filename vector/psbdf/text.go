package psbdf

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/fontchain"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/internal/textrun"
	"github.com/shibukawa/bdf/vector/internal/export"
)

// The fonts of a document are found and text is laid out by
// internal/fontchain, as vector/pdfbdf does; the glyphs in use become
// procedures of their outlines, in Type 3 fonts.

// face is a font program text is drawn with, as Type 3 fonts of the
// outlines of its glyphs in use, of 256 glyphs each.
type face struct {
	fontchain.Program
	out *sfnt.Outlines

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

func newFace(f *sfnt.Font, _ bool) *face {
	return &face{Program: fontchain.NewProgram(f), out: sfnt.NewOutlines(f), codes: map[uint16]code{}}
}

// The faces of a document and the chains of a font reference are found by
// internal/fontchain, which lays text out on them too.
type (
	faceUse = fontchain.Use[*face]
	placed  = fontchain.Placed[*face]
)

// fonts are the faces of the document, with the count of the Type 3
// fonts made of them.
type fonts struct {
	*fontchain.Fonts[*face]
	made int
}

func newFonts(ex *exporter) *fonts {
	return &fonts{Fonts: &fontchain.Fonts[*face]{Doc: ex.Doc, FontFS: ex.opts.FontFS, FontDirs: ex.opts.FontDirs, NoSystemFonts: ex.opts.NoSystemFonts,
		New: newFace, Warnf: ex.Warnf, NotDone: "not written"}}
}

// count is the number of Type 3 fonts made.
func (fs *fonts) count() int { return fs.made }

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

// text writes a FILL_TEXT or STROKE_TEXT: filled text of a color with the
// Type 3 fonts, other text (strokes, gradients, synthesized bold) as the
// paths of the outlines of its glyphs.
func (ex *exporter) text(str string, x, y, advance float64, stroke bool) {
	st := &ex.st
	if !(st.size > 0) || math.IsInf(st.size, 0) || str == "" {
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
		if n := len(runs); n > 0 && runs[n-1].use == g.Use {
			runs[n-1].glyphs = append(runs[n-1].glyphs, g)
			continue
		}
		runs = append(runs, run{use: g.Use, glyphs: []placed{g}})
	}
	for _, r := range runs {
		skew := 0.0
		if r.use.Italic {
			skew = textrun.FakeItalicSkew
		}
		build := func(b *bytes.Buffer) bool {
			return ex.glyphPaths(b, r.use.Face, r.glyphs, x, by, dx, hs, skew)
		}
		switch {
		case stroke:
			ex.paintWith(st.stroke, build, "S", "strokepath W")
		case r.use.Bold:
			// synthesized bold: the fill, and a line around it
			ex.paintWith(st.fill, build, "f", "W")
			saved := st.pColor
			ex.op("q %s w 1 j [] 0 d", export.Num(textrun.FakeBoldWidth(st.size)))
			ex.paintWith(st.fill, build, "S", "strokepath W")
			ex.op("Q")
			st.pColor = saved
		case st.fill.paint == nil || ex.shading(st.fill.paint).solid:
			ex.glyphText(r.use.Face, r.glyphs, x, by, dx, hs, skew)
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
		sub := fc.use(ex.fonts, glyphs[i].GID).sub
		j := i
		var hex, dxs strings.Builder
		for ; j < len(glyphs); j++ {
			c := fc.use(ex.fonts, glyphs[j].GID)
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
				dxs.WriteString(export.Num((glyphs[j+1].X - glyphs[j].X) * hs))
			} else {
				dxs.WriteString(export.Num(fc.Advance(glyphs[j].GID) * size * hs))
			}
		}
		ex.op("/%s [%s] Tf %s m <%s>\n[%s] xshow n", fc.subs[sub].name, export.Mat(fm), export.Nums(x+(dx+glyphs[i].X)*hs, y), hex.String(), dxs.String())
		i = j
	}
}

// glyphPaths writes glyphs of a face as the paths of their outlines, and
// reports whether they have segments.
func (ex *exporter) glyphPaths(b *bytes.Buffer, fc *face, glyphs []placed, x, y, dx, hs, skew float64) bool {
	k := ex.st.size / fc.Upem
	drawn := false
	for _, g := range glyphs {
		m := matrix{hs, 0, 0, 1, x + (dx+g.X)*hs, y}.Mul(matrix{k, 0, 0, -k, 0, 0}).Mul(matrix{1, 0, skew, 1, 0, 0})
		p := export.NewPen(b, m)
		fc.out.Outline(g.GID, p)
		drawn = drawn || p.Drawn
	}
	return drawn
}

// define writes the Type 3 fonts of the glyphs in use, deflated.
func (fs *fonts) define(b *bytes.Buffer) {
	var d bytes.Buffer
	for _, fc := range fs.All {
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
		p := export.NewPen(&path, identity)
		fc.out.Outline(g, p)
		fmt.Fprintf(&procs, "/g%d {%d 0 ", i, fc.Font.Advance(g))
		if p.Drawn && p.Boxed {
			bbox = union(bbox, p.Box)
			fmt.Fprintf(&procs, "%s setcachedevice\n%sf}\n", export.Nums(p.Box[0], p.Box[1], p.Box[2], p.Box[3]), path.Bytes())
		} else {
			procs.WriteString("0 0 0 0 setcachedevice}\n")
		}
	}
	if !(bbox[2] >= bbox[0]) {
		bbox = [4]float64{0, 0, 0, 0}
	}
	upem := export.Num(fc.Upem)
	fmt.Fprintf(b, "/%s <<\n/FontType 3 /FontMatrix [1 %s div 0 0 1 %s div 0 0] /FontBBox [%s]\n/Encoding [", sf.name, upem, upem,
		export.Nums(math.Floor(bbox[0]), math.Floor(bbox[1]), math.Ceil(bbox[2]), math.Ceil(bbox[3])))
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

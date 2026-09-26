// Package jww converts the drawings of Jw_cad (.jww) into BDF documents.
//
// A drawing becomes one page of its sheet (A0 to A4, 2A to 5A, and the
// larger sheets), grown to hold figures drawn outside it. Coordinates in
// a JWW file are millimetres on the sheet with the origin at its centre;
// the scales of the layer groups only matter for the values of
// dimensions, which the file stores as text. Lines, arcs, ellipses,
// points, text, dimensions, solids and blocks are drawn in the colors and
// line types of the file, as Jw_cad shows them on screen (or as it prints
// them, with the colors parameter), with the line widths it prints. Text
// keeps Jw_cad's fixed pitch: full-width characters are as wide as the
// text size, half-width ones half as wide. See docs/design.md §3.13.
package jww

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/width"
)

// Colors selects the colors a drawing is drawn in.
type Colors int

const (
	// Screen draws in the screen colors of the file, on its background.
	Screen Colors = iota
	// Print draws in the printer colors of the file, on white paper.
	Print
	// Mono draws black on white paper, as Jw_cad prints by default.
	Mono
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title.
	Title string
	// Colors selects the colors: the screen's (the default), the
	// printer's, or black.
	Colors Colors
	// FontFS holds fonts that are not in the local file system; it is
	// searched before FontDirs (see converter.Options.FontFS).
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index part.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	W, H          float64 // the size of the page in mm
	EmbeddedFonts int
}

// sheets are the sizes of Jw_cad's sheets in mm, by number.
var sheets = map[int][2]float64{
	0: {1189, 841}, 1: {841, 594}, 2: {594, 420}, 3: {420, 297}, 4: {297, 210},
	8: {1682, 1189}, 9: {2378, 1682}, 10: {3364, 2378}, 11: {4756, 3364},
	12: {10000, 7071}, 13: {50000, 35355}, 14: {100000, 70711},
}

type converter struct {
	opts     *Options
	h        *header
	d        *data
	fonts    *cad.Fonts
	warnings []string
	warned   map[string]bool
	bg       bdf.Color
	count    int   // figures drawn
	blocks   []int // the block definitions being drawn
}

// maxFigures bounds the figures drawn, blocks included.
const maxFigures = 5_000_000

// ConvertFile converts a JWW file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Convert(f, st.Size(), opts)
}

const maxSize = 1 << 30

// Convert converts a drawing read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("jww: the file is larger than %d bytes", maxSize)
	}
	buf := make([]byte, size)
	if _, err := r.ReadAt(buf, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("jww: %w", err)
	}
	c := &converter{opts: opts, warned: map[string]bool{}}
	a := newArchive(buf)
	h, err := readHeader(a)
	if err != nil {
		return nil, fmt.Errorf("jww: %w", err)
	}
	c.h = h
	start := a.i
	// CTime takes 4 bytes in files of older builds of Jw_cad, 8 in newer
	// ones: try both
	var d *data
	for _, ts := range []int{4, 8} {
		b := newArchive(buf)
		b.i = start
		dd, err := readData(b, h, ts)
		if err == nil {
			d = dd
			break
		}
		if d == nil || len(dd.figures) > len(d.figures) {
			d = dd
		}
		if ts == 8 {
			c.warnf("the figures are read only in part: %v", err)
		}
	}
	c.d = d

	doc := bdf.NewDocument()
	doc.Meta.Source = "jww"
	doc.Meta.DC.Language = bdf.DCValues{"ja"}
	if m := strings.TrimSpace(h.memo); m != "" {
		doc.Meta.DC.Description = bdf.DCValues{m}
	}
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	set := fontset.New(db, func(msg string) { c.warnf("%s", msg) })
	c.fonts = &cad.Fonts{Set: set}
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	c.bg = bdf.RGB(255, 255, 255)
	if opts.Colors == Screen {
		c.bg = colorref(h.screen[0])
	}

	dr := &cad.Drawing{}
	c.figures(dr, d.figures, canvas.Identity, 0)

	// the sheet, grown to hold what is drawn outside it
	sw, sh := 420.0, 297.0
	if s, ok := sheets[h.paper]; ok {
		sw, sh = s[0], s[1]
	}
	r0 := cad.Rect{}.Add(cad.Point{X: -sw / 2, Y: -sh / 2}).Add(cad.Point{X: sw / 2, Y: sh / 2})
	if b := dr.Bounds(); b.Valid() {
		r0 = r0.Union(b)
	}
	const mm = 72 / 25.4
	pw, ph := r0.W()*mm, r0.H()*mm
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	page := view.AddPage(float32(pw), float32(ph))
	cvs := canvas.NewBuilder(doc, set)
	cv := cvs.New()
	cv.Obj.SetBBox(0, 0, page.W, page.H)
	cv.Obj.FillColor(c.bg)
	cv.Obj.FillRect(0, 0, page.W, page.H)
	cv.Drawn = true
	m := canvas.Matrix{mm, 0, 0, -mm, -r0.Min.X * mm, r0.Max.Y * mm}
	pl := &cad.Plotter{Fonts: c.fonts, Thin: 0.1 * mm}
	pl.Plot(cv, dr, m, cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: pw, Y: ph}))

	res := &Result{Doc: doc, W: r0.W(), H: r0.H()}
	if !opts.SystemFonts {
		res.EmbeddedFonts = set.Embed(doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	cvs.Encode()
	set.ReportMissing()
	page.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: cv.Hash()}}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	res.Warnings = c.warnings
	return res, nil
}

func colorref(v uint32) bdf.Color { return bdf.RGB(uint8(v), uint8(v>>8), uint8(v>>16)) }

// color returns the color of a line color number.
func (c *converter) color(n int) bdf.Color {
	if c.opts.Colors == Mono {
		if n == 0 {
			return c.bg
		}
		return bdf.RGB(0, 0, 0)
	}
	switch {
	case n >= 0 && n <= 9:
		if c.opts.Colors == Screen {
			return colorref(c.h.screen[n])
		}
		return colorref(c.h.prtPen[n].color)
	case n >= 100 && n <= 356 && c.h.hasSXF:
		return colorref(c.h.sxfPens[n-100].color)
	}
	return colorref(c.h.prtPen[2].color)
}

func (c *converter) pen(n int) pen {
	switch {
	case n >= 0 && n <= 9:
		return c.h.prtPen[n]
	case n >= 100 && n <= 356 && c.h.hasSXF:
		return c.h.sxfPens[n-100]
	}
	return c.h.prtPen[2]
}

// widthMM converts a line width of the file into mm: hundredths of a mm
// when the drawing says so, printer dots otherwise.
func (c *converter) widthMM(w int) float64 {
	if w <= 0 {
		return 0
	}
	if c.h.maxWidth <= -101 {
		return float64(w) / 100
	}
	return float64(w) * 25.4 / float64(c.h.dpi)
}

// visible reports whether a figure's layer is shown.
func (c *converter) visible(b *base) bool {
	g := b.glayer & 15
	l := b.layer & 15
	lg := &c.h.groups[g]
	return lg.state != 0 && lg.layers[l] != 0
}

const mm = 72 / 25.4

func (c *converter) linePen(b *base) cad.Pen {
	p := c.pen(b.color)
	w := p.width
	if b.width > 0 {
		w = b.width
	}
	pen := cad.Pen{Color: c.color(b.color), Width: c.widthMM(w) * mm, Cap: cad.CapRound, Join: cad.JoinRound}
	if w <= 1 && c.h.maxWidth > -101 {
		pen.Width = 0 // the thinnest line
	}
	pen.Dash, pen.DashOffset = c.dash(b.style)
	return pen
}

// dash returns the dash pattern of a line type in mm.
func (c *converter) dash(style int) ([]float64, float64) {
	if style <= 1 || style >= 11 && style <= 15 {
		return nil, 0 // continuous, or a random line drawn straight
	}
	if style >= 30 && style <= 62 && c.h.hasSXF {
		t := c.h.sxfTypes[style-30]
		if t.segments > 0 {
			var elems []float64
			for i := 0; i < t.segments && i < len(t.pitch); i++ {
				v := math.Abs(t.pitch[i])
				if i%2 == 1 {
					v = -v
				}
				elems = append(elems, v)
			}
			return cad.DashPattern(elems)
		}
	}
	lt, ok := c.h.lineTypes[style]
	if !ok || lt.unit <= 0 || lt.unit > 32 {
		return nil, 0
	}
	bit := float64(max(lt.prtPitch, 1)) * 25.4 / float64(c.h.dpi)
	var elems []float64
	for i := 0; i < lt.unit; i++ {
		on := lt.mask>>uint(i)&1 != 0
		v := bit
		if !on {
			v = -bit
		}
		if n := len(elems); n > 0 && (elems[n-1] > 0) == on {
			elems[n-1] += v
		} else {
			elems = append(elems, v)
		}
	}
	return cad.DashPattern(elems)
}

func pt(p point) cad.Point { return cad.Point{X: p.x, Y: p.y} }

// figures draws a list of figures through m.
func (c *converter) figures(out *cad.Drawing, list []any, m canvas.Matrix, depth int) {
	for _, o := range list {
		c.figure(out, o, m, depth)
	}
}

func (c *converter) figure(out *cad.Drawing, o any, m canvas.Matrix, depth int) {
	if c.count++; c.count > maxFigures {
		c.warnOnce("budget", "only the first %d figures are drawn", maxFigures)
		return
	}
	switch f := o.(type) {
	case *sen:
		if !c.visible(&f.base) || c.auxiliary(&f.base) {
			return
		}
		path := (&cad.Path{}).MoveTo(f.a.x, f.a.y).LineTo(f.b.x, f.b.y)
		out.Stroke(path.Transform(m), c.scaled(c.linePen(&f.base), m))
	case *enko:
		if !c.visible(&f.base) || c.auxiliary(&f.base) {
			return
		}
		out.Stroke(c.arcPath(f).Transform(m), c.scaled(c.linePen(&f.base), m))
	case *ten:
		c.point(out, f, m)
	case *moji:
		if !c.visible(&f.base) {
			return
		}
		c.text(out, f, m)
	case *sunpou:
		if !c.visible(&f.base) {
			return
		}
		c.figure(out, &f.line, m, depth)
		c.figure(out, &f.text, m, depth)
		for _, s := range f.aux {
			if s != nil {
				c.figure(out, s, m, depth)
			}
		}
		for _, p := range f.points {
			if p != nil {
				c.point(out, p, m)
			}
		}
	case *solid:
		if !c.visible(&f.base) {
			return
		}
		c.solid(out, f, m)
	case *block:
		if !c.visible(&f.base) {
			return
		}
		def := c.d.blocks[f.number]
		if def == nil {
			c.warnOnce("block", "a block refers to a missing definition")
			return
		}
		if depth > 32 {
			c.warnOnce("depth", "blocks nested deeper than 32 levels are not drawn")
			return
		}
		if slices.Contains(c.blocks, f.number) {
			c.warnOnce("cycle", "a block that inserts itself is drawn once")
			return
		}
		bm := m.Mul(canvas.Translate(f.p.x, f.p.y)).Mul(canvas.Rotate(f.rot * 180 / math.Pi)).Mul(canvas.Scale(f.sx, f.sy))
		c.blocks = append(c.blocks, f.number)
		c.figures(out, def.list, bm, depth+1)
		c.blocks = c.blocks[:len(c.blocks)-1]
	}
}

// auxiliary reports whether a figure is an auxiliary line (補助線), which
// Jw_cad does not print.
func (c *converter) auxiliary(b *base) bool {
	return b.style == 9 || b.color == 9
}

// scaled adapts a pen to the transform of a block: dashes scale with it.
func (c *converter) scaled(p cad.Pen, m canvas.Matrix) cad.Pen {
	s := cad.Scale(m)
	if s != 1 && p.Dash != nil {
		d := make([]float64, len(p.Dash))
		for i, v := range p.Dash {
			d[i] = v * s
		}
		p.Dash, p.DashOffset = d, p.DashOffset*s
	}
	return p
}

func (c *converter) arcPath(e *enko) *cad.Path {
	flat := e.flat
	if flat <= 0 {
		flat = 1
	}
	major := cad.Point{X: e.r * math.Cos(e.tilt), Y: e.r * math.Sin(e.tilt)}
	p := &cad.Path{}
	if e.full || math.Abs(e.arc) >= 2*math.Pi-1e-9 {
		return p.Ellipse(pt(e.c), major, flat)
	}
	return p.EllipseArc(pt(e.c), major, flat, e.start, e.start+e.arc)
}

func (c *converter) point(out *cad.Drawing, t *ten, m canvas.Matrix) {
	if t.kari || !c.visible(&t.base) {
		return
	}
	r := 0.1
	if c.h.drawPrtTen {
		if pr := c.pen(t.color).radius; pr > 0 {
			r = pr
		}
	}
	col := c.color(t.color)
	p := cad.Apply(m, pt(t.p))
	out.Fill((&cad.Path{}).Circle(p, r), cad.Fill{Color: col}, false)
}

func (c *converter) solid(out *cad.Drawing, s *solid, m canvas.Matrix) {
	col := c.color(s.color)
	if s.color == 10 {
		col = colorref(s.rgb)
	}
	path := &cad.Path{}
	if s.style >= 101 {
		// a solid of a circle: centre, radius, flatness, tilt, start and
		// arc angles, and the kind
		ctr := pt(s.p1)
		r, flat := s.p4.x, s.p4.y
		tilt, start, arc := s.p2.x, s.p2.y, s.p3.x
		kind := s.p3.y
		if flat <= 0 {
			flat = 1
		}
		major := cad.Point{X: r * math.Cos(tilt), Y: r * math.Sin(tilt)}
		switch {
		case s.style == 101 && kind == 100, s.style == 111 && kind == 100:
			path.Ellipse(ctr, major, flat)
		case s.style == 101 && kind == 0:
			// a sector
			path.MoveTo(ctr.X, ctr.Y)
			path.EllipseArc(ctr, major, flat, start, start+arc).Close()
		case s.style == 101 && kind == 5, s.style == 101 && kind == -1:
			// a segment (the outer arc solid is the same shape)
			path.EllipseArc(ctr, major, flat, start, start+arc).Close()
		case s.style == 105 || s.style == 106:
			// a ring: outer arc and inner arc back
			inner := kind
			ir := cad.Point{X: inner * math.Cos(tilt), Y: inner * math.Sin(tilt)}
			iflat := flat
			if s.style == 106 && inner > 0 {
				// the same difference along both axes
				iflat = (r*flat - (r - inner)) / inner
			}
			if math.Abs(arc) >= 2*math.Pi-1e-9 || arc == 0 {
				path.Ellipse(ctr, major, flat)
				if inner > 0 {
					path.Ellipse(ctr, ir, iflat)
				}
			} else {
				path.EllipseArc(ctr, major, flat, start, start+arc)
				path.EllipseArc(ctr, ir, iflat, start+arc, start).Close()
			}
		default:
			path.Ellipse(ctr, major, flat)
		}
		out.Fill(path.Transform(m), cad.Fill{Color: col}, true)
		return
	}
	path.Polyline([]cad.Point{pt(s.p1), pt(s.p2), pt(s.p3), pt(s.p4)}, true)
	out.Fill(path.Transform(m), cad.Fill{Color: col}, false)
}

var sjisEnc = japanese.ShiftJIS.NewEncoder()

// fullWidth reports whether Jw_cad gives a character a full cell: two
// bytes in Shift_JIS.
func fullWidth(r rune) bool {
	if r < 0x80 {
		return false
	}
	if b, err := sjisEnc.String(string(r)); err == nil {
		return len(b) == 2
	}
	k := width.LookupRune(r).Kind()
	return k == width.EastAsianWide || k == width.EastAsianFullwidth
}

func (c *converter) text(out *cad.Drawing, t *moji, m canvas.Matrix) {
	s := t.s
	if strings.HasPrefix(s, "^@") {
		c.warnOnce("image", "images and other special text (^@) are not drawn")
		return
	}
	if strings.TrimSpace(s) == "" || t.h <= 0 {
		return
	}
	kind := t.kind
	spec := cad.FontSpec{Family: t.font, Italic: kind/10000%2 == 1, Bold: kind/20000 >= 1}
	em := t.h
	var cells []float64
	for _, r := range s {
		w := t.w / em
		if !fullWidth(r) {
			w /= 2
		}
		cells = append(cells, w)
	}
	spacing := t.spacing / em
	angle := t.angle * math.Pi / 180
	var tm canvas.Matrix
	if t.flag&0x0020 == 0 {
		// the reference point is the bottom left corner of the cells
		tm = m.Mul(canvas.Translate(t.a.x, t.a.y)).Mul(rotation(angle)).Mul(canvas.Translate(0, 0.14*em)).Mul(canvas.Scale(em, em))
	} else {
		// vertical text: the text box turned a quarter turn clockwise about
		// the reference point, its top left corner; the characters stand
		// upright down the centre of the column
		tm = m.Mul(canvas.Translate(t.a.x, t.a.y)).Mul(rotation(angle)).Mul(canvas.Translate(t.w/2, 0)).Mul(rotation(-math.Pi / 2)).Mul(canvas.Scale(em, em))
	}
	tx := c.fonts.NewCellText(spec, s, cells, spacing, c.color(t.color), tm)
	tx.Vertical = t.flag&0x0020 != 0
	out.Text(tx)
}

func rotation(a float64) canvas.Matrix {
	s, co := math.Sincos(a)
	return canvas.Matrix{co, s, -s, co, 0, 0}
}

func (c *converter) warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

func (c *converter) warnOnce(key, format string, args ...any) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warnf(format, args...)
}

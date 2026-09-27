// Package hpgl converts HP-GL/2 plot files (.plt) into BDF documents: the
// vector language of HP plotters and printers, with the HP-GL of older
// pen plotters and cutters, the PJL job and PCL or HP RTL escape sequences
// that files printed to a large-format plotter or a PCL printer carry
// around it, and the raster images of HP RTL.
//
// Each plotted page (a page ends with PG, BP, a reset or a form feed)
// becomes a page of one fixed view. A plot's coordinates are plotter units
// (0.025 mm) on the paper, so the page is the plot size (PS) the file sets,
// grown to hold what is drawn outside it, or the size of what is drawn
// when it sets none; HP-GL/2 in a PCL job is placed in the picture frame of
// the PCL page. Pens draw in the colors and widths of the file's palette
// (PC, PW), with its line types, line ends and joins, fills (solid,
// hatching, shading, raster patterns) and clipping window; labels are text
// in the fonts of the document, sized, turned, slanted and placed as the
// plotter places them. See docs/design.md §3.19.
package hpgl

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/hpglsniff"
	"github.com/shibukawa/bdf/imgconv"
)

// Colors selects the colors the pens draw in.
type Colors int

const (
	// Pens draws in the colors of the file's palette (the default).
	Pens Colors = iota
	// Mono draws every pen but pen 0 (white) in black, as a monochrome
	// plotter does.
	Mono
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title.
	Title string
	// Colors selects the colors of the pens.
	Colors Colors
	// Pages selects 1-based pages (see converter.Pages); nil converts all
	// of them. Numbers past the last page are reported as warnings.
	Pages conv.Pages
	// Images controls how the raster images of HP RTL are stored and the
	// resolution cap they are scaled down to (see imgconv.Options).
	Images imgconv.Options
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
	Doc      *bdf.Document
	Warnings []string
	// Pages is the number of pages converted; Sizes are their sizes in mm.
	Pages int
	Sizes [][2]float64
	// Images is the number of raster images stored.
	Images        int
	EmbeddedFonts int
}

// ConvertFile converts a plot file.
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

// maxSize bounds the input read.
const maxSize = 1 << 30

// Budgets that keep hostile files from running away.
const (
	maxItems = 5_000_000 // items drawn in the whole document
	maxPages = 10_000
)

// plu is the size of a plotter unit in pt.
const plu = 72.0 / 1016

// Convert converts a plot read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("hpgl: the file is larger than %d bytes", maxSize)
	}
	buf := make([]byte, size)
	if _, err := r.ReadAt(buf, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("hpgl: %w", err)
	}
	c := newConverter(opts)
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	c.set = fontset.New(db, func(msg string) { c.warnf("%s", msg) })
	c.fonts = &cad.Fonts{Set: c.set}
	c.doc = bdf.NewDocument()
	c.doc.Meta.Source = "hpgl"
	c.run(buf)
	c.rtl.end()
	c.endPage()
	if len(c.pages) == 0 && !hpglsniff.Is(buf[:min(len(buf), 1024)]) || c.instructions == 0 && c.rasters == 0 {
		// text that is not a plot (a gnuplot script named .plt)
		return nil, fmt.Errorf("hpgl: the input is not an HP-GL/2 plot")
	}
	if len(c.pages) == 0 {
		// a plot that draws nothing is a blank page of its plot size
		c.pages = append(c.pages, c.blankPage())
	}
	return c.finish(db)
}

// finish lays out the pages that were plotted into the document.
func (c *converter) finish(db *fontdb.DB) (*Result, error) {
	opts := c.opts
	doc := c.doc
	switch {
	case opts.Title != "":
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	case c.title != "":
		doc.Meta.DC.Title = bdf.DCValues{c.title}
	case c.jobName != "":
		doc.Meta.DC.Title = bdf.DCValues{c.jobName}
	}
	if c.hasText && len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	pages := c.pages
	if opts.Pages != nil {
		var sel []*page
		for _, n := range opts.Pages.Numbers(len(pages)) {
			if n < 1 || n > len(pages) {
				c.warnf("page %d does not exist (the plot has %d)", n, len(pages))
				continue
			}
			sel = append(sel, pages[n-1])
		}
		pages = sel
	}
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	cvs := canvas.NewBuilder(doc, c.set)
	pl := &cad.Plotter{Fonts: c.fonts, Thin: 0.1 * 72 / 25.4}
	type laid struct {
		pg *bdf.Page
		cv *canvas.Canvas
	}
	var out []laid
	res := &Result{Doc: doc, Images: c.images}
	for _, p := range pages {
		m, w, h := c.pageTransform(p)
		pg := view.AddPage(float32(w), float32(h))
		cv := cvs.New()
		cv.Obj.SetBBox(0, 0, pg.W, pg.H)
		cv.Obj.FillColor(bdf.RGB(255, 255, 255))
		cv.Obj.FillRect(0, 0, pg.W, pg.H)
		cv.Drawn = true
		pl.Plot(cv, &p.d, m, cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: w, Y: h}))
		out = append(out, laid{pg, cv})
		res.Sizes = append(res.Sizes, [2]float64{w * 25.4 / 72, h * 25.4 / 72})
	}
	res.Pages = len(out)
	if !opts.SystemFonts {
		res.EmbeddedFonts = c.set.Embed(doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	cvs.Encode()
	c.set.ReportMissing()
	for _, l := range out {
		l.pg.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: l.cv.Hash()}}
	}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	res.Warnings = c.warnings
	return res, nil
}

// pageMargin is the white space around a drawing whose paper is not known
// (plotter units: 5 mm).
const pageMargin = 200

// pageTransform returns the transform of a page's drawing (plotter units,
// y up) onto the page (pt, y down) and the page's size in pt.
func (c *converter) pageTransform(p *page) (canvas.Matrix, float64, float64) {
	if p.pcl != nil {
		// the physical page of a PCL printer; the drawing is in its
		// coordinates already (see pcl.go)
		return canvas.Scale(plu, plu), p.pcl.w * plu, p.pcl.h * plu
	}
	b := p.d.Bounds()
	view := p.view // drawing to view coordinates (y up)
	if view == (canvas.Matrix{}) {
		view = canvas.Identity
	}
	var r cad.Rect
	if p.hard.Valid() {
		r = p.hard.Transform(view)
	}
	if b.Valid() {
		vb := b.Transform(view)
		if !r.Valid() {
			vb = vb.Add(cad.Point{X: vb.Min.X - pageMargin, Y: vb.Min.Y - pageMargin}).Add(cad.Point{X: vb.Max.X + pageMargin, Y: vb.Max.Y + pageMargin})
		}
		r = r.Union(vb)
	}
	if !r.Valid() {
		r = cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: defaultHardW, Y: defaultHardH})
	}
	// at least a millimetre, at most 200 m
	w := math.Min(math.Max(r.W(), 40), 8e6)
	h := math.Min(math.Max(r.H(), 40), 8e6)
	m := canvas.Matrix{plu, 0, 0, -plu, -r.Min.X * plu, r.Max.Y * plu}.Mul(view)
	return m, w * plu, h * plu
}

// blankPage is the page of a plot that draws nothing.
func (c *converter) blankPage() *page {
	p := &page{}
	if c.gl != nil && c.gl.psSet {
		p.hard = cad.Rect{}.Add(cad.Point{}).Add(cad.Point{X: c.gl.hardW, Y: c.gl.hardH})
	}
	return p
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

// summary describes a result in a line.
func summary(res *Result) string {
	var parts []string
	if res.Pages == 1 && len(res.Sizes) == 1 {
		parts = append(parts, fmt.Sprintf("1 page (%.0f × %.0f mm)", res.Sizes[0][0], res.Sizes[0][1]))
	} else {
		parts = append(parts, fmt.Sprintf("%d page(s)", res.Pages))
	}
	if res.Images > 0 {
		parts = append(parts, fmt.Sprintf("%d raster image(s)", res.Images))
	}
	parts = append(parts, fmt.Sprintf("%d embedded font(s)", res.EmbeddedFonts))
	return strings.Join(parts, ", ")
}

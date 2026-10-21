// Package psbdf writes BDF documents as PostScript (Level 3), for the
// printers that take PostScript and not PDF: office printers of the years
// before PDF was common, the PostScript of other makers (KPDL, BR-Script),
// and large-format PostScript printers.
//
// The instructions of the objects (docs/spec.md §7) become the operators
// of PostScript, in the same order, as vector/pdfbdf writes them as those
// of a PDF: paths, transforms, clips, lines and dashes, colors, linear and
// radial gradients (the smooth shadings of Level 3), images (JPEG as it
// is, others as their pixels) and the objects they use. Text is drawn
// with Type 3 fonts made of the outlines of the glyphs in use, which every
// interpreter draws the same, and laid out as raster/imagebdf lays it out.
//
// PostScript has no transparency. What it cannot draw (alpha, blend modes,
// groups with an alpha, a blend mode or a soft mask, shadows, filters,
// images with alpha, gradients with alpha, patterns and conic gradients)
// is found in a first pass over each page, and the rectangles where it is
// drawn are drawn as pictures (vector/internal/flatten): the rest of the
// page stays vectors, under a clip that leaves the rectangles out.
//
// The file follows the Document Structuring Conventions, so that spoolers
// can count and reorder its pages, and asks the printer for the paper,
// the copies, the sides and the tray with setpagedevice, each as a feature
// that a printer without it skips. Images and pictures are written in
// ASCII85, so that the file is text that any connection passes.
package psbdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/vector/internal/export"
	"github.com/shibukawa/bdf/vector/internal/flatten"
)

// Options controls the writing of a file.
type Options struct {
	// Views are the ids of the views to write. nil writes the views of
	// pages (fixed and flow).
	Views []string
	// Pages selects the pages written, by their zero-based number among
	// the pages of the views, in that order; -1 writes a blank page (the
	// sheet between copies of a two-sided document). nil writes all of
	// them.
	Pages []int
	// SkipRoles are the roles of the layers that are left out
	// (bdf.RoleBackground …): for paper that has them printed already.
	SkipRoles []string
	// FontFS, FontDirs and NoSystemFonts find the fonts of text of fonts
	// referred to by name, and of characters an embedded font lacks.
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	// Place places a page of w × h points on a sheet: the size of the
	// sheet in points, and the transform from the space of the page (its
	// y axis down from its top left corner) to that of the sheet (y up
	// from its bottom left corner). nil prints a page on a sheet of its
	// size.
	Place func(w, h float64) (sheetW, sheetH float64, m [6]float64)
	// Copies, Collate, Sides ("one-sided", "two-sided-long-edge",
	// "two-sided-short-edge"), MediaType ("stationery") and MediaPosition
	// (the input tray) are asked of the printer; what is not set (0, "")
	// is left to it.
	Copies        int
	Collate       bool
	Sides         string
	MediaType     string
	MediaPosition int
	// Gray writes colors as gray (luminance): for black and white
	// printing.
	Gray bool
	// RasterDPI is the resolution of what is drawn as pictures: 300 when
	// 0.
	RasterDPI int
	// Title is the %%Title of the file; the title of the document when "".
	Title string
}

// Result tells what was written.
type Result struct {
	Pages int
	// Fonts is the number of Type 3 fonts made.
	Fonts int
	// Pictures counts the rectangles drawn as pictures.
	Pictures int
	Warnings []string
}

// ErrNoPages is returned for a document without pages to write.
var ErrNoPages = errors.New("psbdf: the document has no pages")

// pageOut is a page to write.
type pageOut struct {
	view  *bdf.View
	index int // in its view; -1 for a blank page
	w, h  float64
}

// Write writes a document as a PostScript file.
func Write(w io.Writer, doc *bdf.Document, opts *Options) (res *Result, err error) {
	defer func() {
		if p := recover(); p != nil {
			res, err = nil, fmt.Errorf("psbdf: the writing failed: %v", p)
		}
	}()
	if doc == nil {
		return nil, errors.New("psbdf: no document")
	}
	ex := newExporter(doc, opts)
	pages, err := ex.pages()
	if err != nil {
		return nil, err
	}
	// the pages first: the fonts they use are defined before them
	var bodies [][]byte
	sheets := make([][2]float64, len(pages))
	for i, pg := range pages {
		sw, sh, _ := ex.place(pg)
		sheets[i] = [2]float64{sw, sh}
		bodies = append(bodies, ex.page(pg, i))
	}
	bw := &bytes.Buffer{}
	ex.header(bw, len(pages), sheets)
	if _, err := w.Write(bw.Bytes()); err != nil {
		return nil, err
	}
	for i, b := range bodies {
		fmt.Fprintf(w, "%%%%Page: %d %d\n", i+1, i+1)
		if i > 0 && sheets[i] != sheets[i-1] {
			fmt.Fprintf(w, "%%%%BeginPageSetup\n%s%%%%EndPageSetup\n", feature("PageSize", fmt.Sprintf("/PageSize [%s]", export.Nums(sheets[i][0], sheets[i][1]))))
		}
		if _, err := w.Write(b); err != nil {
			return nil, err
		}
	}
	if _, err := io.WriteString(w, "%%Trailer\nend\n%%EOF\n"); err != nil {
		return nil, err
	}
	return &Result{Pages: len(pages), Fonts: ex.fonts.count(), Pictures: ex.pictures, Warnings: ex.Warnings}, nil
}

// pages returns the pages to write.
func (ex *exporter) pages() ([]pageOut, error) {
	doc := ex.Doc
	var views []*bdf.View
	if len(ex.opts.Views) > 0 {
		for _, id := range ex.opts.Views {
			found := false
			for _, v := range doc.Views {
				if v != nil && v.ID == id {
					views, found = append(views, v), true
				}
			}
			if !found {
				return nil, fmt.Errorf("psbdf: no view %q", id)
			}
		}
	} else {
		for _, v := range doc.Views {
			if v != nil && (v.Kind == bdf.ViewFixed || v.Kind == bdf.ViewFlow) {
				views = append(views, v)
			}
		}
	}
	var all []pageOut
	for _, v := range views {
		if v.Kind != bdf.ViewFixed && v.Kind != bdf.ViewFlow {
			return nil, fmt.Errorf("psbdf: view %q has no pages (a %s view)", v.ID, v.Kind)
		}
		for i, p := range v.Pages {
			if p == nil || !(p.W > 0) || !(p.H > 0) || math.IsInf(float64(p.W), 0) || math.IsInf(float64(p.H), 0) {
				continue
			}
			all = append(all, pageOut{view: v, index: i, w: float64(p.W), h: float64(p.H)})
		}
	}
	if ex.opts.Pages == nil {
		if len(all) == 0 {
			return nil, ErrNoPages
		}
		return all, nil
	}
	var out []pageOut
	for _, n := range ex.opts.Pages {
		switch {
		case n == -1 && len(out) > 0:
			out = append(out, pageOut{view: out[len(out)-1].view, index: -1, w: out[len(out)-1].w, h: out[len(out)-1].h})
		case n >= 0 && n < len(all):
			out = append(out, all[n])
		}
	}
	if len(out) == 0 {
		return nil, ErrNoPages
	}
	return out, nil
}

// place returns the sheet of a page and the transform from the page to it.
func (ex *exporter) place(pg pageOut) (sheetW, sheetH float64, m matrix) {
	if ex.opts.Place != nil {
		sw, sh, mm := ex.opts.Place(pg.w, pg.h)
		if sw > 0 && sh > 0 && matrix(mm).Finite() {
			return sw, sh, matrix(mm)
		}
	}
	return pg.w, pg.h, matrix{1, 0, 0, -1, 0, pg.h}
}

// prolog defines the procedures of the pages: the operators of a PDF's
// content, which the pages are written with as pdfbdf writes them.
const prolog = `/bdf 64 dict def bdf begin
/m /moveto load def /l /lineto load def /c /curveto load def /h /closepath load def
/re { 4 2 roll moveto 1 index 0 rlineto 0 exch rlineto neg 0 rlineto closepath } bind def
/f /fill load def /f* /eofill load def /S /stroke load def /n /newpath load def
/W /clip load def /W* /eoclip load def
/q /gsave load def /Q /grestore load def /cm { 6 array astore concat } bind def
/w /setlinewidth load def /J /setlinecap load def /j /setlinejoin load def /M /setmiterlimit load def
/d /setdash load def /rg /setrgbcolor load def /g /setgray load def
/B { gsave fill grestore stroke } bind def
/Tf { exch findfont exch makefont setfont } bind def
% what follows in the file, deflated in ASCII85, is run
/bdfrun { currentfile /ASCII85Decode filter dup /FlateDecode filter cvx exec flushfile } bind def
% dict colorspace filter: an image of the data that follows in ASCII85
/bdfimage { exch setcolorspace currentfile /ASCII85Decode filter dup 3 -1 roll filter
  3 -1 roll dup /DataSource 4 -1 roll put image flushfile } bind def
end
`

// header writes what comes before the pages: the comments, the prolog,
// the setup of the device and the fonts.
func (ex *exporter) header(b *bytes.Buffer, pages int, sheets [][2]float64) {
	maxW, maxH := 0.0, 0.0
	for _, s := range sheets {
		maxW, maxH = math.Max(maxW, s[0]), math.Max(maxH, s[1])
	}
	title := ex.opts.Title
	if title == "" {
		title = ex.Doc.Meta.DC.Title.First()
	}
	b.WriteString("%!PS-Adobe-3.0\n%%Creator: bdf (vector/psbdf)\n")
	if title != "" {
		fmt.Fprintf(b, "%%%%Title: %s\n", psText(title))
	}
	fmt.Fprintf(b, "%%%%LanguageLevel: 3\n%%%%Pages: %d\n%%%%BoundingBox: 0 0 %d %d\n%%%%HiResBoundingBox: 0 0 %s\n%%%%DocumentData: Clean7Bit\n%%%%EndComments\n",
		pages, int(math.Ceil(maxW)), int(math.Ceil(maxH)), export.Nums(maxW, maxH))
	b.WriteString("%%BeginProlog\n%%BeginResource: procset bdf 1 0\n")
	b.WriteString(prolog)
	b.WriteString("%%EndResource\n%%EndProlog\n%%BeginSetup\n")
	o := ex.opts
	b.WriteString(feature("PageSize", fmt.Sprintf("/PageSize [%s]", export.Nums(sheets[0][0], sheets[0][1]))))
	if o.Copies > 1 {
		b.WriteString(feature("NumCopies", fmt.Sprintf("/NumCopies %d", o.Copies)))
		if o.Collate {
			b.WriteString(feature("Collate", "/Collate true"))
		}
	}
	switch o.Sides {
	case "one-sided":
		b.WriteString(feature("Duplex", "/Duplex false"))
	case "two-sided-long-edge":
		b.WriteString(feature("Duplex", "/Duplex true /Tumble false"))
	case "two-sided-short-edge":
		b.WriteString(feature("Duplex", "/Duplex true /Tumble true"))
	}
	if o.MediaType != "" {
		b.WriteString(feature("MediaType", "/MediaType "+psText(o.MediaType)))
	}
	if o.MediaPosition > 0 {
		b.WriteString(feature("MediaPosition", fmt.Sprintf("/MediaPosition %d", o.MediaPosition)))
	}
	b.WriteString("bdf begin\n")
	ex.fonts.define(b)
	b.WriteString("%%EndSetup\n")
}

// feature asks the device for something, as a feature that a printer
// without it skips.
func feature(name, entries string) string {
	return fmt.Sprintf("%%%%BeginFeature: *%s\nmark { << %s >> setpagedevice } stopped cleartomark\n%%%%EndFeature\n", name, entries)
}

// psText writes a string of PostScript: in parentheses, with what is not
// printable ASCII escaped.
func psText(s string) string {
	var b strings.Builder
	b.WriteByte('(')
	for _, c := range []byte(s) {
		switch {
		case c == '(' || c == ')' || c == '\\':
			b.WriteByte('\\')
			b.WriteByte(c)
		case c < 0x20 || c > 0x7e:
			fmt.Fprintf(&b, "\\%03o", c)
		default:
			b.WriteByte(c)
		}
	}
	b.WriteByte(')')
	return b.String()
}

// page writes a page: its vectors, under a clip that leaves out the
// rectangles of what is drawn as pictures, then the pictures.
func (ex *exporter) page(pg pageOut, n int) []byte {
	b := &bytes.Buffer{}
	_, _, place := ex.place(pg)
	b.WriteString("save\n")
	if pg.index < 0 {
		b.WriteString("restore showpage\n")
		return b.Bytes()
	}
	p := pg.view.Pages[pg.index]
	// the first pass: what PostScript cannot draw
	var regions flatten.Regions
	ex.begin(nil, &regions)
	ex.layers(p)
	boxes := regions.Merge(pg.w, pg.h, 1)
	// the second: the vectors
	ex.begin(b, nil)
	fmt.Fprintf(b, "%s cm\n", export.Mat(place))
	if len(boxes) > 0 {
		// the page, with holes where the pictures go (the boxes do not overlap)
		fmt.Fprintf(b, "q %s re", export.Nums(-huge, -huge, 2*huge, 2*huge))
		for _, bx := range boxes {
			fmt.Fprintf(b, " %s re", export.Nums(bx[0], bx[1], bx[2]-bx[0], bx[3]-bx[1]))
		}
		b.WriteString(" W* n\n")
	}
	if len(boxes) != 1 || boxes[0] != [4]float64{0, 0, pg.w, pg.h} {
		ex.layers(p)
	}
	if len(boxes) > 0 {
		b.WriteString("Q\n")
		ex.picturesOf(b, pg, boxes)
	}
	b.WriteString("restore showpage\n")
	return b.Bytes()
}

// huge is a rectangle that holds every drawing.
const huge = 1 << 15

// layers writes the layers of a page that are not skipped.
func (ex *exporter) layers(p *bdf.Page) {
	for _, l := range p.Layers {
		if ex.Skip(l.Role) {
			continue
		}
		ex.top(l.Obj)
	}
}

// picturesOf draws the rectangles of a page as pictures of what the image
// renderer draws there.
func (ex *exporter) picturesOf(b *bytes.Buffer, pg pageOut, boxes [][4]float64) {
	if ex.renderer == nil {
		ex.renderer = imagebdf.New(ex.Doc, &imagebdf.Options{FontFS: ex.opts.FontFS, FontDirs: ex.opts.FontDirs, NoSystemFonts: ex.opts.NoSystemFonts,
			ImageCacheBytes: -1})
	}
	dpi := float64(ex.opts.RasterDPI)
	if dpi <= 0 {
		dpi = 300
	}
	view := ex.withoutSkipped(pg.view)
	for _, bx := range boxes {
		img, err := flatten.Picture(ex.renderer, view, pg.index, bx, dpi)
		if err != nil {
			ex.Warnf("a part of a page is not drawn: %v", err)
			continue
		}
		ex.pictures++
		writeImage(b, ex.picture(img), matrix{bx[2] - bx[0], 0, 0, bx[3] - bx[1], bx[0], bx[1]}, false)
	}
	for _, w := range ex.renderer.Warnings() {
		ex.Warnf("%s", w)
	}
}

// withoutSkipped returns a view without the layers of the roles left out, for the
// image renderer.
func (ex *exporter) withoutSkipped(v *bdf.View) *bdf.View {
	if len(ex.opts.SkipRoles) == 0 {
		return v
	}
	if s, ok := ex.views[v]; ok {
		return s
	}
	c := *v
	c.Pages = make([]*bdf.Page, len(v.Pages))
	for i, p := range v.Pages {
		if p == nil {
			continue
		}
		cp := *p
		cp.Layers = nil
		for _, l := range p.Layers {
			if !ex.Skip(l.Role) {
				cp.Layers = append(cp.Layers, l)
			}
		}
		c.Pages[i] = &cp
	}
	ex.views[v] = &c
	return &c
}

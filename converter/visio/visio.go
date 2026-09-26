// Package visio converts Visio drawings into BDF documents: the packages of
// Visio 2013 and later (.vsdx, .vsdm, .vstx, .vstm) and the XML drawings of
// Visio 2003 to 2010 (.vdx, .vtx). The binary .vsd format is not read.
//
// Each foreground page becomes a page of the size of its paper, with its
// background pages as background layers (identical on every page that
// shares them, so the container stores them once) and its own shapes as
// the body layer. Shapes are drawn from their ShapeSheets: the values
// Visio computed for the cells, inherited from masters and styles, and the
// dynamic theme for the cells it sets (see sheet.go and theme.go). Text is
// laid out by the DrawingML renderer the Office converters share, with the
// metrics of the fonts that are then embedded as subsets. See
// docs/design.md §3.8.
package visio

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based foreground pages (see converter.Pages); nil
	// converts every one.
	Pages conv.Pages
	// Title overrides the document title.
	Title string
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontDirs.
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
	Pages         int
	EmbeddedFonts int
}

type converter struct {
	opts     *Options
	doc      *bdf.Document
	d        *document
	fonts    *fontset.Set
	cvs      *canvas.Builder
	dm       *drawingml.Drawing
	warnings []string
	warned   map[string]bool
	pageNum  map[int]int // page ID → output page number
	patterns map[any]bdf.Hash
	pictures map[string]*picture
}

// ConvertFile converts a Visio drawing file.
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

// maxXMLSize bounds the XML drawings read.
const maxXMLSize = 1 << 30

// Convert converts a drawing read from r: a package, or an XML drawing.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	c := &converter{opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, pageNum: map[int]int{},
		patterns: map[any]bdf.Hash{}, pictures: map[string]*picture{}}
	warnKey := func(key, msg string) { c.warnOnce(key, "%s", msg) }
	head := make([]byte, 4)
	r.ReadAt(head, 0)
	var err error
	if bytes.Equal(head, []byte("PK\x03\x04")) {
		p, err := ooxml.Open(r, size)
		if err != nil {
			return nil, fmt.Errorf("visio: %w", err)
		}
		if c.d, err = readVSDX(p, warnKey); err != nil {
			return nil, fmt.Errorf("visio: %w", err)
		}
		c.doc.Meta.Source = "vsdx"
	} else {
		if size > maxXMLSize {
			return nil, fmt.Errorf("visio: the drawing is larger than %d bytes", maxXMLSize)
		}
		data := make([]byte, size)
		if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
			return nil, fmt.Errorf("visio: %w", err)
		}
		if c.d, err = readVDX(data, warnKey); err != nil {
			return nil, fmt.Errorf("visio: %w", err)
		}
		c.doc.Meta.Source = "vdx"
	}
	d := c.d
	warn := func(msg string) { c.warnf("%s", msg) }
	db := fontdb.New(opts.FontDirs, !opts.NoSystemFonts)
	c.fonts = fontset.New(db, warn)
	c.cvs = canvas.NewBuilder(c.doc, c.fonts)
	rd := drawingml.New(drawingml.Config{Package: d.pkg, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: warn})
	c.dm = rd.NewDrawing("", nil, textHost{})
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	c.doc.Meta.DC = d.dc
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	if len(c.doc.Meta.DC.Language) == 0 {
		// the language of the text of the root style sheet
		if root := d.styles[0]; root != nil {
			if x := root.sections[secKey{"Character", -1}]; x != nil && x.rows["0"] != nil {
				if l := langTag(x.rows["0"].cells["LangID"].v); l != "" {
					c.doc.Meta.DC.Language = bdf.DCValues{l}
				}
			}
		}
	}
	view := c.doc.NewView("pages", bdf.ViewFixed, c.doc.Meta.DC.Title.First())

	var fg []*page
	for _, pg := range d.pages {
		if !pg.background {
			fg = append(fg, pg)
		}
	}
	sel := opts.Pages.Numbers(len(fg))
	if sel == nil {
		for i := range fg {
			sel = append(sel, i+1)
		}
	}
	for i, n := range sel {
		if n < 1 || n > len(fg) {
			return nil, fmt.Errorf("visio: page %d out of range (1-%d)", n, len(fg))
		}
		c.pageNum[fg[n-1].id] = i + 1
	}
	type pageRef struct {
		page   *bdf.Page
		layers []*canvas.Canvas
	}
	var pages []pageRef
	for _, n := range sel {
		pg := fg[n-1]
		w, h := pageSize(pg)
		page := view.AddPage(f32(w), f32(h))
		pr := pageRef{page: page}
		for _, bg := range c.backgrounds(pg) {
			if cv := c.renderSafe(bg, h, "background page "+bg.name); cv.Drawn {
				page.Layers = append(page.Layers, bdf.Layer{Role: bdf.RoleBackground})
				pr.layers = append(pr.layers, cv)
			}
		}
		if cv := c.renderSafe(pg, h, "page "+pg.name); cv.Drawn {
			page.Layers = append(page.Layers, bdf.Layer{Role: bdf.RoleBody})
			pr.layers = append(pr.layers, cv)
		}
		pages = append(pages, pr)
	}
	res := &Result{Doc: c.doc, Pages: len(sel)}
	if !opts.SystemFonts {
		res.EmbeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	c.fonts.ReportMissing()
	for _, pr := range pages {
		for i, cv := range pr.layers {
			pr.page.Layers[i].Obj = cv.Hash()
		}
		if len(pr.page.Layers) == 0 {
			// an empty page still needs something to draw
			o := bdf.NewObject()
			o.SetBBox(0, 0, pr.page.W, pr.page.H)
			hash, _ := c.doc.AddObject(o)
			pr.page.Layers = append(pr.page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
		}
	}
	if !opts.NoTextIndex {
		if _, err := c.doc.BuildTextIndex(view); err != nil {
			c.warnf("text index: %v", err)
		}
	}
	res.Warnings = c.warnings
	return res, nil
}

// backgrounds lists the background pages of a page, the farthest first.
func (c *converter) backgrounds(pg *page) []*page {
	var out []*page
	seen := map[int]bool{pg.id: true}
	for id := pg.backPage; id >= 0 && !seen[id]; {
		seen[id] = true
		bg := c.d.pageByID(id)
		if bg == nil {
			break
		}
		out = append([]*page{bg}, out...)
		id = bg.backPage
	}
	return out
}

// renderSafe draws the shapes of a page into a new layer; a page whose
// malformed content trips the renderer is left empty with a warning.
func (c *converter) renderSafe(pg *page, height float64, what string) (cv *canvas.Canvas) {
	cv = c.cvs.New()
	defer func() {
		if r := recover(); r != nil {
			c.warnf("%s: internal error: %v", what, r)
			cv = c.cvs.New()
		}
	}()
	c.newPageCtx(pg, height).drawPage(cv)
	return cv
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

// textHost is the DrawingML host of Visio text: nothing is inherited from
// placeholders or document text styles, and fields keep their text.
type textHost struct{}

func (textHost) Placeholder(*ooxml.Node, string) ([]drawingml.Inherited, bool) { return nil, true }
func (textHost) TextStyle(string) *ooxml.Node                                  { return nil }
func (textHost) Field(string) (string, bool)                                   { return "", false }
func (textHost) Link(string, ooxml.Rel) string                                 { return "" }

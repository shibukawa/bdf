// Package wordproc lays out word-processing documents for the converters of
// Word documents (converter/docx) and of HTML and Markdown
// (converter/html, converter/markdown).
//
// BDF has no layout engine, so documents are laid out here, with the
// metrics of the fonts that are then embedded as subsets: paragraphs are
// broken into lines (tab stops, East Asian line breaking rules and
// spacing, justification, the document grid), tables into rows and cells,
// and the lines and rows into columns and pages (keeping lines and
// paragraphs together, widows and orphans, floating objects, footnotes).
//
// Both inputs are read into the same model: blocks (paragraphs and
// tables) whose text is already measured. WordprocessingML (docx.go) is
// read with Word's formatting rules; drawings are rendered by the
// DrawingML renderer the Office converters share
// (converter/internal/ooxml/drawingml), and their text boxes hold
// paragraphs laid out here. HTML (html.go) is read in reader mode: the
// elements take the formatting of a fixed style sheet, not the page's CSS.
//
// A document converts into up to two views (docs/spec.md §4.1): a flow
// view of pages, with header, body and footer layers and the body
// rectangle of each page, and a scroll view laid out once more without
// pages, as one long column of the text width (what Word's draft and web
// layouts show), stored as strips cut between lines. The views share fonts
// and images. Sections of East Asian vertical text are laid out as pages
// turned by 90° in the page view, and horizontally in the scroll view. See
// docs/design.md §3.9 and §3.13.
package wordproc

import (
	"fmt"
	"io/fs"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/imgconv"
)

// View selections.
const (
	ViewsBoth   = "both"
	ViewsPages  = "pages"
	ViewsScroll = "scroll"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages of the page view; nil keeps every page.
	// The scroll view always holds the whole document.
	Pages []int
	// Views selects the views to make: ViewsBoth, ViewsPages or
	// ViewsScroll. "" takes the input's default: both views for Word
	// documents, the scroll view for HTML.
	Views string
	// Title overrides the document title.
	Title string
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// FontFS holds fonts that are not in the local file system; it is
	// searched before FontDirs (see converter.Options.FontFS).
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout. Viewers then substitute their own fonts; the
	// text advances keep lines at the widths computed here.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2
	// (builds tagged bdf_noconv never produce WOFF2).
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index parts.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Pages         int // pages of the page view
	Strips        int // strips of the scroll view
	EmbeddedFonts int
}

type converter struct {
	pkg      *ooxml.Package
	opts     *Options
	doc      *bdf.Document
	db       *fontdb.DB
	fonts    *fontset.Set
	cvs      *canvas.Builder
	r        *drawingml.Renderer
	dr       *drawingml.Drawing
	warnings []string
	warned   map[string]bool

	main        string
	st          *styles
	num         *numbering
	sections    []*section
	lastSection *section
	clrMap      map[string]string
	eaScript    string
	baseSize    float64 // the default font size: the character grid's cell less its extra pitch

	// settings
	defTab                 float64
	evenOdd                bool
	compatMode             int
	noExpandShiftReturn    bool
	suppressSpBfAfterPgBrk bool
	background             bdf.Color
	hasBackground          bool
	footFmt, endFmt        string

	notes      map[string]*note
	noteOrder  []*note // referenced notes in document order
	notesUsed  map[string]int
	hf         map[string][]block // header and footer blocks by part
	boxes      map[*ooxml.Node][]block
	noteCache  map[noteKey]*flow
	hfCache    map[hfKey]*flow
	bmPage     map[string][2]int // bookmark → physical page index, printed number
	bmY        map[string]float64
	curEmitter *emitter
	cjk, latin int
	tcyGroups  int

	// css lays blocks out with CSS rules instead of Word's: adjoining
	// margins collapse, and the margin above the first block of an area
	// (a column, a page, a cell) is dropped.
	css bool
	// scrollWidth is the text width of the scroll view (0: the widest
	// section's).
	scrollWidth float64
}

// newConverter sets up the state of a conversion; views defaults to
// defViews.
func newConverter(opts *Options, defViews string) (*converter, string, error) {
	if opts == nil {
		opts = &Options{}
	}
	views := opts.Views
	switch views {
	case "":
		views = defViews
	case ViewsBoth, ViewsPages, ViewsScroll:
	default:
		return nil, "", fmt.Errorf("unknown views %q (want both, pages or scroll)", views)
	}
	c := &converter{opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, clrMap: map[string]string{},
		notes: map[string]*note{}, notesUsed: map[string]int{}, hf: map[string][]block{}, boxes: map[*ooxml.Node][]block{},
		noteCache: map[noteKey]*flow{}, hfCache: map[hfKey]*flow{}, bmPage: map[string][2]int{}, bmY: map[string]float64{},
		defTab: 36, compatMode: 12, footFmt: "decimal", endFmt: "lowerRoman"}
	c.db = fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	c.fonts = fontset.New(c.db, c.warn)
	c.cvs = canvas.NewBuilder(c.doc, c.fonts)
	return c, views, nil
}

func (c *converter) warn(msg string) { c.warnf("%s", msg) }

// checkFonts warns when there is no font to lay text out with.
func (c *converter) checkFonts() {
	if len(c.db.Faces) == 0 && !c.opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
}

// finish lays out the views of the document read into the sections,
// embeds the fonts and builds the text indexes.
func (c *converter) finish(views string) *Result {
	opts := c.opts
	res := &Result{Doc: c.doc}
	title := c.doc.Meta.DC.Title.First()
	var jobs []func()
	if views != ViewsScroll {
		v := c.doc.NewView("pages", bdf.ViewFlow, title)
		v.Continuous = &bdf.Continuous{Gap: 24}
		jobs = append(jobs, c.pageView(v, &res.Pages))
	}
	if views != ViewsPages {
		v := c.doc.NewView("scroll", bdf.ViewScroll, title)
		jobs = append(jobs, c.scrollView(v, &res.Strips))
	}
	if !opts.SystemFonts {
		res.EmbeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	c.fonts.ReportMissing()
	for _, job := range jobs {
		job()
	}
	if !opts.NoTextIndex {
		for _, v := range c.doc.Views {
			if _, err := c.doc.BuildTextIndex(v); err != nil {
				c.warnf("text index: %v", err)
			}
		}
	}
	res.Warnings = c.warnings
	return res
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

// Package docx converts Word documents (.docx) into BDF documents.
//
// BDF has no layout engine, so the document is laid out by the converter,
// with the metrics of the fonts that are then embedded as subsets:
// paragraphs are broken into lines (tab stops, East Asian line breaking
// rules and spacing, justification, the document grid), tables into rows
// and cells, and the lines and rows into columns and pages (keeping lines
// and paragraphs together, widows and orphans, floating objects,
// footnotes). Drawings are rendered by the DrawingML renderer the Office
// converters share (converter/internal/ooxml/drawingml); their text boxes
// hold paragraphs laid out by the converter. The layout engine is shared
// with the HTML and Markdown converters (converter/internal/wordproc).
//
// A document converts into two views (docs/spec.md §4.1): a flow view of
// pages, with header, body and footer layers and the body rectangle of
// each page, and a scroll view laid out once more without pages, as one
// long column of the text width (what Word's draft and web layouts show),
// stored as strips cut between lines. The views share fonts and images.
// Sections of East Asian vertical text are laid out as pages turned by 90°
// in the page view, and horizontally in the scroll view. See
// docs/design.md §3.9.
package docx

import (
	"io"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/wordproc"
	"github.com/shibukawa/bdf/imgconv"
)

// View selections.
const (
	ViewsBoth   = wordproc.ViewsBoth
	ViewsPages  = wordproc.ViewsPages
	ViewsScroll = wordproc.ViewsScroll
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages of the page view; nil keeps every page.
	// The scroll view always holds the whole document.
	Pages []int
	// Views selects the views to make: ViewsBoth (the default when ""),
	// ViewsPages or ViewsScroll.
	Views string
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

// ConvertFile converts a .docx file.
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

// Convert converts a document read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	var wo *wordproc.Options
	if opts != nil {
		o := wordproc.Options(*opts)
		wo = &o
	}
	res, err := wordproc.ConvertDOCX(r, size, wo)
	if err != nil {
		return nil, err
	}
	out := Result(*res)
	return &out, nil
}

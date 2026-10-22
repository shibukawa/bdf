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
	"io/fs"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/wordproc"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/imgconv"
)

// View selections.
const (
	ViewsBoth   = wordproc.ViewsBoth
	ViewsPages  = wordproc.ViewsPages
	ViewsScroll = wordproc.ViewsScroll
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages of the page view (see converter.Pages);
	// nil keeps every page. The scroll view always holds the whole document.
	Pages conv.Pages
	// Views selects the views to make: ViewsBoth (the default when ""),
	// ViewsPages or ViewsScroll.
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
	// Data says whether the values of charts are stored as tables the
	// document holds beside its drawing (docs/spec.md §4.6): "none" not
	// at all; "" and "all" alike, a chart's values being few.
	Data string
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Pages         int // pages of the page view
	Strips        int // strips of the scroll view
	EmbeddedFonts int
}

// TemplateData is what a template is filled with (see Fill).
type TemplateData struct {
	// Data is the data: the values of JSON (nil, bool, float64, string,
	// []any, map[string]any), and time.Time for dates.
	Data map[string]any
	// Strict makes a tag whose name the data does not have an error;
	// without it the tag is replaced by nothing, with a warning.
	Strict bool
	// Files holds the pictures that the data names by their paths.
	Files fs.FS
	// Format writes values in the number formats of tags
	// ({{total | fmt:"#,##0"}}); converter/xlsx has one (Formatter).
	// Without it the formats of tags are not applied.
	Format tag.Formatter
}

// Filled is a template filled with data.
type Filled struct {
	Result
	// Fields are the fields of the template: the names that its tags take
	// from the data.
	Fields []tag.Field
}

// Fill converts a document as a template filled with data: its text holds
// tags (package form/tag) that the data replaces. A tag is replaced by its
// value in the format of the text it is in; a section repeats a part of
// the text of a paragraph, the paragraphs and tables from its one tag to
// the other, or the rows of a table (the row, for tags in a row; for tags
// in different rows the rows between them are slots, a row for an item).
// "| rows" in the tag of a section repeats the row or the paragraph it is
// in, "| pages" its paragraphs with each item on a new page. A picture
// whose alternative text has a picture tag ({{%photo}}) is replaced by the
// picture the data names.
//
// What the data makes longer flows on onto further pages as text does. The
// body, the headers and footers and the text boxes are filled; footnotes
// are not. The page view alone is made unless the options ask for more.
func Fill(r io.ReaderAt, size int64, t *TemplateData, opts *Options) (*Filled, error) {
	var wo *wordproc.Options
	if opts != nil {
		o := wordproc.Options(*opts)
		wo = &o
	}
	var wt *wordproc.TemplateData
	if t != nil {
		d := wordproc.TemplateData(*t)
		wt = &d
	}
	res, fields, err := wordproc.FillDOCX(r, size, wt, wo)
	if err != nil {
		return nil, err
	}
	return &Filled{Result: Result(*res), Fields: fields}, nil
}

// Inspect reads the tags of a document without filling it: the fields that
// data for it has, or the first mistake in its tags.
func Inspect(r io.ReaderAt, size int64) ([]tag.Field, error) {
	return wordproc.InspectDOCX(r, size)
}

// ConvertFile converts a .docx file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	return conv.OpenFile(path, func(r io.ReaderAt, size int64) (*Result, error) { return Convert(r, size, opts) })
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

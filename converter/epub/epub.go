// Package epub converts EPUB publications (EPUB 3 and EPUB 2) into BDF
// documents.
//
// A reflowable book is laid out like the HTML converter's reader mode
// (converter/html): the content documents of the spine, in reading order,
// are laid out one after the other by the layout engine of the Word
// converter (converter/internal/wordproc) with a fixed style sheet, each
// starting a page. The book's CSS is not applied, but the properties that
// carry meaning in books are read from it (see css.go): the writing mode,
// tate-chu-yoko, emphasis marks, alignment, hidden elements and the size
// of pictures. The default view is one of book pages (A5, with page
// numbers); chapters in East Asian vertical text (writing-mode:
// vertical-rl, as most Japanese books are) are laid out vertically there.
// A scroll view, horizontal throughout, can be added or made instead. SVG
// is drawn as the HTML converter draws it: SVG files as they are, svg
// elements and SVG content documents as SVG documents of their own. MathML
// formulas are laid out by the formula engine (converter/internal/equation),
// upright on the middle of the line in vertical text; an epub:switch takes
// its MathML case.
//
// A fixed-layout book (rendition:layout pre-paginated) whose pages are
// pictures (comics, photo books) converts into a fixed view of those
// pictures, a page each, at the size of the page's viewport: raster
// images, SVG files, and svg elements (drawn as the image they fit into
// their view box, when that is all they hold). Books bound on the right
// (page-progression-direction rtl) have views whose direction is rtl.
// Metadata comes from the package document. Encrypted publications (DRM)
// are refused; fonts obfuscated as the EPUB specification describes are
// left out like other embedded fonts. See docs/design.md §3.24.
package epub

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/webdoc"
	"github.com/shibukawa/bdf/converter/internal/wordproc"
	"github.com/shibukawa/bdf/imgconv"
	"golang.org/x/net/html/charset"
)

var charsetReader = charset.NewReaderLabel

// View selections.
const (
	ViewsPages  = wordproc.ViewsPages
	ViewsScroll = wordproc.ViewsScroll
	ViewsBoth   = wordproc.ViewsBoth
)

// Paper sizes in points.
var (
	A4 = Paper{595.3, 841.9}
	A5 = Paper{419.5, 595.3}
	A6 = Paper{297.6, 419.5}
	B5 = Paper{515.9, 728.5} // JIS
	B6 = Paper{362.8, 515.9} // JIS
	// Letter is US Letter.
	Letter = Paper{612, 792}
)

// Paper is the size of the pages of the page view in points.
type Paper struct{ Width, Height float64 }

// ParsePaper reads a paper size: a name (a4, a5, a6, b5, b6, letter) or
// the width and height in millimeters ("148x210").
func ParsePaper(s string) (Paper, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "a4":
		return A4, nil
	case "a5":
		return A5, nil
	case "a6", "bunko":
		return A6, nil
	case "b5":
		return B5, nil
	case "b6":
		return B6, nil
	case "letter":
		return Letter, nil
	}
	w, h, ok := strings.Cut(strings.TrimSuffix(strings.ToLower(strings.TrimSpace(s)), "mm"), "x")
	if ok {
		fw, err1 := strconv.ParseFloat(strings.TrimSpace(w), 64)
		fh, err2 := strconv.ParseFloat(strings.TrimSpace(h), 64)
		if err1 == nil && err2 == nil && fw >= 30 && fh >= 30 && fw <= 2000 && fh <= 2000 {
			return Paper{fw * 72 / 25.4, fh * 72 / 25.4}, nil
		}
	}
	return Paper{}, fmt.Errorf("unknown paper %q (want a4, a5, a6, b5, b6, letter or WIDTHxHEIGHT in millimeters)", s)
}

// Options controls the conversion.
type Options struct {
	// Views selects the views of a reflowable book: ViewsPages (the default
	// when ""), ViewsScroll or ViewsBoth. A fixed-layout book of pictures
	// has one view of its pages.
	Views string
	// Pages selects 1-based pages of the page view (see converter.Pages);
	// nil keeps every page.
	Pages conv.Pages
	// Title overrides the document title.
	Title string
	// Paper is the size of the pages (zero: A5). The margins are a tenth
	// of the width at the sides and 9 % of the height above and below.
	Paper Paper

	// The reader style:

	// Width is the text width of the scroll view in points (0: 36 ems).
	Width float64
	// FontSize is the size of the text in points (0: 12).
	FontSize float64
	// Font and MonoFont are the font families of the text and of code
	// ("": installed ones are picked).
	Font, MonoFont string

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
	// EmbedFonts embeds the fonts the text is laid out with, as subsets.
	// By default they are referred to by name, as the HTML converter does,
	// and viewers draw the text with their own fonts at the advances
	// measured here.
	EmbedFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2.
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
	Pages         int // pages of the page view (or of the fixed layout)
	Strips        int // strips of the scroll view
	EmbeddedFonts int
	Chapters      int  // content documents laid out
	Images        int  // pictures drawn
	Vertical      bool // some chapters are vertical text
	FixedLayout   bool // a fixed-layout book of pictures
}

// ConvertFile converts an EPUB file.
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

// Convert converts a publication read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	switch opts.Views {
	case "", ViewsPages, ViewsScroll, ViewsBoth:
	default:
		return nil, fmt.Errorf("epub: unknown views %q (want pages, scroll or both)", opts.Views)
	}
	pub, err := open(r, size)
	if err != nil {
		return nil, err
	}
	c := &converter{pub: pub, opts: opts, byPath: map[string]*chapter{}, sheets: map[string]string{}, warned: map[string]bool{},
		parsed: map[string][]cssRule{}, stylers: map[string]*styler{}}
	return c.convert()
}

// converter is the state of a conversion.
type converter struct {
	pub      *publication
	opts     *Options
	dc       bdf.DublinCore
	props    map[string]string
	chapters []*chapter
	byPath   map[string]*chapter
	sheets   map[string]string    // style sheets by path
	parsed   map[string][]cssRule // their rules
	stylers  map[string]*styler   // by the style sheets of a chapter (see styler)
	modes    bool                 // the style sheets say writing modes
	svgBytes int                  // of the SVG documents made of the svg elements of fixed-layout pages
	tags     int                  // of the content documents read
	warnings []string
	warned   map[string]bool
}

func (c *converter) warn(msg string) {
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

func (c *converter) warnOnce(key, msg string) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warn(msg)
}

// warnLeftOut warns of a file that is left out for an error: with msg, or
// once for the publication when it is left out because the files read
// before it have reached what may be read (errInflated, errTags).
func (c *converter) warnLeftOut(err error, msg string) {
	switch {
	case errors.Is(err, errInflated):
		c.warnOnce("inflated", fmt.Sprintf("the files of the publication are more than %d MiB larger than it, decompressed; the files read after them are left out",
			maxInflated>>20))
	case errors.Is(err, errTags):
		c.warnOnce("tags", fmt.Sprintf("the content documents of the publication have more than %d tags; the content documents after them are left out", maxTags))
	default:
		c.warn(msg)
	}
}

func (c *converter) convert() (*Result, error) {
	c.dc, c.props = c.pub.metadata()
	if c.opts.Title != "" {
		c.dc.Title = bdf.DCValues{c.opts.Title}
	}
	if err := c.spine(); err != nil {
		return nil, err
	}
	for _, ch := range c.chapters {
		if err := c.load(ch); err != nil {
			if errors.Is(err, ErrDRM) {
				return nil, err
			}
			c.warnLeftOut(err, fmt.Sprintf("%s: %v; left out", ch.path, err))
		}
	}
	var loaded []*chapter
	for _, ch := range c.chapters {
		if ch.body != nil {
			loaded = append(loaded, ch)
		}
	}
	if len(loaded) == 0 {
		return nil, errors.New("epub: no content document could be read")
	}
	c.chapters = loaded
	rtl := strings.EqualFold(c.pub.pkg.Spine.PPD, "rtl")
	if c.fixedLayout() {
		res, err := c.pictures()
		if err != nil {
			return nil, err
		}
		direction(res.Doc, rtl)
		return res, nil
	}
	res, err := c.reflow()
	if err != nil {
		return nil, err
	}
	direction(res.Doc, rtl)
	return res, nil
}

// direction marks the page views of a book bound on the right.
func direction(doc *bdf.Document, rtl bool) {
	if !rtl {
		return
	}
	for _, v := range doc.Views {
		if v.Kind == bdf.ViewFixed || v.Kind == bdf.ViewFlow {
			v.Direction = bdf.DirectionRTL
		}
	}
}

// image returns the bytes of a picture by the reference a chapter was
// rewritten with: a path in the container or a data: URL.
func (c *converter) image(src string) ([]byte, error) {
	if webdoc.IsDataURL(src) {
		return webdoc.DataURL(src)
	}
	if i := strings.IndexByte(src, ':'); i > 0 && !strings.ContainsAny(src[:i], "/.") {
		return nil, errors.New("pictures outside the publication are not fetched")
	}
	return c.pub.picture(src)
}

// link keeps the absolute http:, https: and mailto: links.
func link(href string) string {
	s := strings.TrimSpace(href)
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "mailto:") {
		return s
	}
	return ""
}

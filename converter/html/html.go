// Package html converts HTML documents (.html, .xhtml, and web pages saved
// as MHTML) into BDF documents in reader mode.
//
// Like a browser's reader view, the converter keeps the content and leaves
// the page's design out: the article is picked out of the page (with
// go-readability, a port of Mozilla's Readability) and its elements take
// the formatting of a fixed style sheet, not the page's CSS. Headings,
// paragraphs, lists, quotations, code, tables, figures, links and images
// are laid out by the layout engine of the Word converter
// (converter/internal/wordproc), with the fonts embedded as subsets. BDF
// does not lay text out again, so the document converts into a scroll view
// of a fixed width (docs/spec.md §4.1); a view of A4 pages can be added.
//
// Images are read from the files beside the document, from an MHTML
// archive, from data: URLs and, unless Options.NoRemote is set, from the
// network. See docs/design.md §3.10.
package html

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"codeberg.org/readeck/go-readability/v2"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/wordproc"
	"github.com/shibukawa/bdf/imgconv"
	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"
)

// View selections.
const (
	ViewsScroll = wordproc.ViewsScroll
	ViewsPages  = wordproc.ViewsPages
	ViewsBoth   = wordproc.ViewsBoth
)

// Article extraction.
const (
	// ExtractAuto picks the article out of pages that look like they have
	// one (Readability's isProbablyReaderable), and keeps the whole of
	// other documents.
	ExtractAuto = "auto"
	// ExtractArticle always picks the article out.
	ExtractArticle = "article"
	// ExtractNone keeps the whole document.
	ExtractNone = "none"
)

// Options controls the conversion.
type Options struct {
	// Views selects the views to make: ViewsScroll (the default when ""),
	// ViewsPages (A4 pages) or ViewsBoth.
	Views string
	// Pages selects 1-based pages of the page view; nil keeps every page.
	Pages []int
	// Title overrides the document title.
	Title string
	// Extract controls how the article is picked out of the page:
	// ExtractAuto (the default when ""), ExtractArticle or ExtractNone.
	Extract string

	// Dir is the directory relative references (images) resolve in; ""
	// reads no local files. ConvertFile sets it to the file's directory.
	Dir string
	// BaseURL is the address of the document: relative references resolve
	// against it, and links become absolute. A <base> element or the
	// address an MHTML archive was saved from takes precedence.
	BaseURL string
	// NoRemote leaves out the images on the network instead of fetching
	// them. A server that converts documents it does not trust should set
	// it, or give a Fetch that limits what can be reached.
	NoRemote bool
	// Fetch loads an http: or https: URL; nil uses an HTTP client with a
	// timeout and a size limit.
	Fetch func(url string) ([]byte, error)

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
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontDirs.
	NoSystemFonts bool
	// EmbedFonts embeds the fonts the text is laid out with, as subsets.
	// By default the fonts are referred to by name (the text families as
	// the generic sans-serif and monospace), and viewers draw the text with
	// their own fonts, keeping the advances measured here.
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
	Pages         int  // pages of the page view
	Strips        int  // strips of the scroll view
	EmbeddedFonts int  //
	Extracted     bool // the article was picked out of the page
	Images        int  // images drawn
}

// ConvertFile converts an HTML or MHTML file; relative references resolve
// in its directory (unless opts.Dir says otherwise).
func ConvertFile(path string, opts *Options) (*Result, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.Dir == "" {
		o.Dir = filepath.Dir(path)
	}
	return ConvertBytes(data, &o)
}

// Convert converts a document read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	data, err := io.ReadAll(io.NewSectionReader(r, 0, size))
	if err != nil {
		return nil, err
	}
	return ConvertBytes(data, opts)
}

// ConvertBytes converts an HTML document or an MHTML archive.
func ConvertBytes(data []byte, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	res := newResources(opts)
	contentType := ""
	if IsMHTML(data) {
		m, err := readMHTML(data)
		if err != nil {
			return nil, fmt.Errorf("html: %w", err)
		}
		data, contentType = m.html, m.contentType
		res.parts = m.parts
		if m.location != "" {
			res.setBase(m.location)
		}
	}
	doc, err := parse(data, contentType)
	if err != nil {
		return nil, fmt.Errorf("html: %w", err)
	}
	return convert(doc, opts, res)
}

// parse decodes a document to UTF-8 by its byte order mark, the charset of
// its content type or of its meta element, and parses it.
func parse(data []byte, contentType string) (*xhtml.Node, error) {
	enc, _, certain := charset.DetermineEncoding(data, contentType)
	if !certain && utf8.Valid(data) {
		enc = nil
	}
	var r io.Reader = bytes.NewReader(data)
	if enc != nil {
		r = enc.NewDecoder().Reader(r)
	}
	return xhtml.Parse(r)
}

// ConvertNode converts a parsed document.
func ConvertNode(doc *xhtml.Node, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	return convert(doc, opts, newResources(opts))
}

// convert converts a parsed document whose references res resolves.
func convert(doc *xhtml.Node, opts *Options, res *resources) (*Result, error) {
	if b := findElement(doc, atom.Base); b != nil {
		if href := attr(b, "href"); href != "" {
			res.setBase(href)
		}
	}
	meta := readMeta(doc)
	body := doc
	extracted := false
	switch opts.Extract {
	case "", ExtractAuto, ExtractArticle:
		if opts.Extract == ExtractArticle || readability.CheckDocument(doc) {
			if a, err := readability.FromDocument(doc, res.base); err == nil && a.Node != nil {
				body = articleBody(a, meta)
				meta.fromArticle(a)
				extracted = true
			} else if opts.Extract == ExtractArticle {
				warn(opts, "no article found; the whole document is kept")
			}
		}
	case ExtractNone:
	default:
		return nil, fmt.Errorf("html: unknown extraction %q (want auto, article or none)", opts.Extract)
	}
	dc := meta.dublinCore(body)
	if opts.Title != "" {
		dc.Title = bdf.DCValues{opts.Title}
	}
	res.prefetch(body)
	images := 0
	d := &wordproc.HTMLDocument{Body: body, DC: dc, Source: "html", Width: opts.Width, FontSize: opts.FontSize,
		Font: opts.Font, MonoFont: opts.MonoFont, Link: res.link,
		Image: func(src string) ([]byte, error) {
			b, err := res.image(src)
			if err == nil {
				images++
			}
			return b, err
		}}
	wo := &wordproc.Options{Pages: opts.Pages, Views: opts.Views, Title: opts.Title, Images: opts.Images,
		FontDirs: opts.FontDirs, NoSystemFonts: opts.NoSystemFonts, SystemFonts: !opts.EmbedFonts, NoSubset: opts.NoSubset,
		NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType, NoTextIndex: opts.NoTextIndex, Warn: opts.Warn}
	r, err := wordproc.ConvertHTML(d, wo)
	if err != nil {
		return nil, err
	}
	out := &Result{Doc: r.Doc, Warnings: r.Warnings, Pages: r.Pages, Strips: r.Strips, EmbeddedFonts: r.EmbeddedFonts,
		Extracted: extracted, Images: images}
	if opts.Warn == nil {
		out.Warnings = append(res.warnings, out.Warnings...)
	}
	return out, nil
}

func warn(opts *Options, msg string) {
	if opts.Warn != nil {
		opts.Warn(msg)
	}
}

// articleBody puts the article that Readability picked out under the
// heading that a reader view shows: the title, and a line with the byline,
// the site and the date.
func articleBody(a readability.Article, m *metadata) *xhtml.Node {
	body := &xhtml.Node{Type: xhtml.ElementNode, Data: "body", DataAtom: atom.Body}
	if l := a.Language(); l != "" {
		body.Attr = append(body.Attr, xhtml.Attribute{Key: "lang", Val: l})
	}
	title := strings.TrimSpace(a.Title())
	if title == "" {
		title = m.title
	}
	if title != "" {
		h := element("h1", atom.H1)
		h.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: title})
		body.AppendChild(h)
	}
	var parts []string
	for _, s := range []string{a.Byline(), a.SiteName()} {
		if s = strings.Join(strings.Fields(s), " "); s != "" && !containsFold(parts, s) {
			parts = append(parts, s)
		}
	}
	if t, err := a.PublishedTime(); err == nil && !t.IsZero() {
		parts = append(parts, t.Format("2006-01-02"))
	}
	if len(parts) > 0 {
		p, small := element("p", atom.P), element("small", atom.Small)
		small.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: strings.Join(parts, " · ")})
		p.AppendChild(small)
		body.AppendChild(p)
	}
	if title != "" || len(parts) > 0 {
		body.AppendChild(element("hr", atom.Hr))
	}
	a.Node.Parent, a.Node.PrevSibling, a.Node.NextSibling = nil, nil, nil
	body.AppendChild(a.Node)
	return body
}

func containsFold(list []string, s string) bool {
	for _, x := range list {
		if strings.EqualFold(x, s) {
			return true
		}
	}
	return false
}

func element(tag string, a atom.Atom) *xhtml.Node {
	return &xhtml.Node{Type: xhtml.ElementNode, Data: tag, DataAtom: a}
}

func attr(n *xhtml.Node, name string) string {
	for _, a := range n.Attr {
		if a.Namespace == "" && a.Key == name {
			return strings.TrimSpace(a.Val)
		}
	}
	return ""
}

// findElement returns the first element of a kind, in document order.
func findElement(n *xhtml.Node, a atom.Atom) *xhtml.Node {
	if n.Type == xhtml.ElementNode && n.DataAtom == a {
		return n
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if f := findElement(k, a); f != nil {
			return f
		}
	}
	return nil
}

// resolve resolves a reference against a base URL (nil: kept as it is).
func resolve(base *url.URL, ref string) (*url.URL, error) {
	u, err := url.Parse(strings.TrimSpace(ref))
	if err != nil {
		return nil, err
	}
	if base != nil {
		u = base.ResolveReference(u)
	}
	return u, nil
}

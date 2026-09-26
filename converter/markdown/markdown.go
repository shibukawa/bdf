// Package markdown converts Markdown documents into BDF documents.
//
// The document is rendered to HTML with goldmark (CommonMark with GitHub's
// extensions: tables, task lists, strikethrough, autolinks, and
// footnotes and definition lists), and the HTML is laid out in reader mode
// by the HTML converter (converter/html): the raw HTML a Markdown document
// holds (<p align="center">, <details>, <img width>) is laid out with the
// rest. Headings get the ids GitHub gives them, so that links to them
// ("#見出し") work. YAML or TOML front matter becomes the document's
// metadata. See docs/design.md §3.13.
package markdown

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf/converter/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

// Options are those of HTML documents; Extract is not used (the whole
// document is kept).
type Options = html.Options

// Result is the outcome of a conversion.
type Result = html.Result

// ConvertFile converts a Markdown file; relative references resolve in its
// directory (unless opts.Dir says otherwise).
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

// ConvertBytes converts a Markdown document.
func ConvertBytes(src []byte, opts *Options) (*Result, error) {
	o := Options{}
	if opts != nil {
		o = *opts
	}
	o.Extract = html.ExtractNone
	doc, err := ToHTML(src)
	if err != nil {
		return nil, err
	}
	res, err := html.ConvertBytes(doc, &o)
	if err != nil {
		return nil, err
	}
	res.Doc.Meta.Source = "markdown"
	return res, nil
}

// ToHTML renders a Markdown document as an HTML document, with its front
// matter in the head.
func ToHTML(src []byte) ([]byte, error) {
	src = toUTF8(src)
	fm, body := splitFrontMatter(src)
	var out bytes.Buffer
	out.WriteString("<!DOCTYPE html>\n")
	fm.writeHead(&out)
	out.WriteString("<body>\n")
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.DefinitionList,
			extension.NewFootnote(extension.WithFootnoteBacklinkHTML("↑"))),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(gmhtml.WithUnsafe()),
	)
	ctx := parser.NewContext(parser.WithIDs(&githubIDs{used: map[string]bool{}}))
	if err := md.Convert(body, &out, parser.WithContext(ctx)); err != nil {
		return nil, err
	}
	out.WriteString("</body>\n</html>\n")
	return out.Bytes(), nil
}

// toUTF8 decodes a document that is not UTF-8: Shift_JIS or EUC-JP when
// it decodes as one, else Windows-1252.
func toUTF8(b []byte) []byte {
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	if utf8.Valid(b) {
		return b
	}
	for _, dec := range []func([]byte) ([]byte, error){japanese.ShiftJIS.NewDecoder().Bytes, japanese.EUCJP.NewDecoder().Bytes} {
		if d, err := dec(b); err == nil && !bytes.ContainsRune(d, utf8.RuneError) {
			return d
		}
	}
	d, _ := charmap.Windows1252.NewDecoder().Bytes(b)
	return d
}

// IsText reports whether data looks like text: no NUL bytes and few other
// control characters.
func IsText(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	ctl := 0
	for _, c := range data {
		switch {
		case c == 0:
			return false
		case c < 0x20 && c != '\n' && c != '\r' && c != '\t' && c != '\f' && c != 0x1b:
			ctl++
		}
	}
	return ctl*100 < len(data)
}

// githubIDs gives headings the ids GitHub does: the text in lower case,
// without punctuation, with hyphens for spaces; repeated ids get -1, -2 …
type githubIDs struct {
	used map[string]bool
}

func (g *githubIDs) Generate(value []byte, kind ast.NodeKind) []byte {
	var b strings.Builder
	for _, r := range strings.ToLower(string(value)) {
		switch {
		case r == ' ':
			b.WriteByte('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.Is(unicode.M, r):
			b.WriteRune(r)
		}
	}
	slug := b.String()
	if slug == "" {
		slug = "heading"
	}
	id := slug
	for n := 1; g.used[id]; n++ {
		id = slug + "-" + strconv.Itoa(n)
	}
	g.used[id] = true
	return []byte(id)
}

func (g *githubIDs) Put(value []byte) { g.used[string(value)] = true }

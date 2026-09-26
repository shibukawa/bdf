package html

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	conv "github.com/shibukawa/bdf/converter"
)

// Params are the format-specific options of HTML and Markdown documents.
var Params = []conv.Param{
	{Name: "views", Usage: "scroll (default): one long column; pages: A4 pages; both"},
	{Name: "width", Usage: "text width of the scroll view in points (default: 36 ems of the text)"},
	{Name: "size", Usage: "text size in points (default 12)"},
	{Name: "font", Usage: "font family of the text (default: an installed sans-serif family)"},
	{Name: "mono", Usage: "font family of code (default: an installed monospaced family)"},
	{Name: "remote", Usage: "true (default): fetch images from the network; false: leave them out"},
	{Name: "base", Usage: "URL the document came from: relative links and images resolve against it"},
}

func init() {
	conv.Register(&conv.Format{
		Name:        "html",
		Description: "HTML document or web archive, in reader mode",
		Extensions:  []string{".html", ".htm", ".xhtml", ".mhtml", ".mht"},
		Params: append([]conv.Param{{Name: "extract", Usage: "auto (default): pick the article out of web pages; article: always; none: keep the whole page"}},
			Params...),
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return IsHTML(head) || IsMHTML(head) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts, err := FromConverter(o)
			if err != nil {
				return nil, err
			}
			opts.Extract = o.Param("extract")
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			return Summary(res), nil
		},
	})
}

// FromConverter maps the registry's options (and Params) onto Options.
func FromConverter(o *conv.Options) (*Options, error) {
	opts := &Options{Views: o.Param("views"), Pages: o.Pages, Title: o.Title, Dir: o.Dir, BaseURL: o.Param("base"),
		Font: o.Param("font"), MonoFont: o.Param("mono"), Images: o.Images, FontDirs: o.FontDirs,
		NoSystemFonts: o.NoSystemFonts, EmbedFonts: o.EmbedFonts && !o.SystemFonts, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2,
		IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn}
	for _, p := range []struct {
		name string
		dst  *float64
	}{{"width", &opts.Width}, {"size", &opts.FontSize}} {
		if v := o.Param(p.name); v != "" {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil || f <= 0 {
				return nil, fmt.Errorf("parameter %s: %q is not a positive number", p.name, v)
			}
			*p.dst = f
		}
	}
	if o.Param("remote") != "" {
		remote, err := o.BoolParam("remote")
		if err != nil {
			return nil, err
		}
		opts.NoRemote = !remote
	}
	return opts, nil
}

// Summary makes a registry result of a result.
func Summary(res *Result) *conv.Result {
	var s string
	switch {
	case res.Pages > 0 && res.Strips > 0:
		s = fmt.Sprintf("%d page(s), %d scroll strip(s)", res.Pages, res.Strips)
	case res.Pages > 0:
		s = fmt.Sprintf("%d page(s)", res.Pages)
	default:
		s = fmt.Sprintf("%d scroll strip(s)", res.Strips)
	}
	s += fmt.Sprintf(", %d image(s), %d embedded font(s)", res.Images, res.EmbeddedFonts)
	if res.Extracted {
		s += ", article extracted"
	}
	return &conv.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: s}
}

// IsHTML reports whether data starts like a whole HTML document: a
// doctype or an html element, after white space, comments and an XML
// declaration. Fragments (a README that starts with <p align="center">)
// are left to the Markdown converter, whose raw HTML they are.
func IsHTML(data []byte) bool {
	b := bytes.TrimPrefix(data[:min(len(data), 4096)], []byte("\xef\xbb\xbf"))
	for {
		b = bytes.TrimLeft(b, " \t\r\n\f")
		lower := bytes.ToLower(b[:min(len(b), 16)])
		switch {
		case bytes.HasPrefix(lower, []byte("<?xml")):
			i := bytes.Index(b, []byte("?>"))
			if i < 0 {
				return false
			}
			b = b[i+2:]
		case bytes.HasPrefix(lower, []byte("<!--")):
			i := bytes.Index(b, []byte("-->"))
			if i < 0 {
				return false
			}
			b = b[i+3:]
		case bytes.HasPrefix(lower, []byte("<!doctype html")):
			return true
		case bytes.HasPrefix(lower, []byte("<html")) || bytes.HasPrefix(lower, []byte("<head")):
			return len(lower) > 5 && (lower[5] == '>' || lower[5] == ' ' || lower[5] == '\n' || lower[5] == '\t' || lower[5] == '\r')
		default:
			return false
		}
	}
}

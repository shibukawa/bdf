package markdown

import (
	"bytes"
	"io"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/html"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "markdown",
		Description: "Markdown document (CommonMark, GitHub Flavored Markdown), in reader mode",
		Extensions:  []string{".md", ".markdown", ".mdown", ".mkd", ".mdx"},
		Params:      html.Params,
		// any text: Markdown is what is left when no other format claims it
		Fallback: true,
		Detect:   func(head []byte, r io.ReaderAt, size int64) bool { return IsText(head) && !isXML(head) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts, err := html.FromConverter(o)
			if err != nil {
				return nil, err
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			return html.Summary(res), nil
		},
	})
}

// isXML reports whether text starts like an XML document other than HTML
// (an SVG picture, a data file), which is not Markdown.
func isXML(head []byte) bool {
	b := bytes.TrimLeft(bytes.TrimPrefix(head, []byte("\xef\xbb\xbf")), " \t\r\n")
	lower := bytes.ToLower(b[:min(len(b), 16)])
	return bytes.HasPrefix(lower, []byte("<?xml")) || bytes.HasPrefix(lower, []byte("<svg")) ||
		bytes.HasPrefix(lower, []byte("<!doctype")) && !bytes.HasPrefix(lower, []byte("<!doctype html"))
}

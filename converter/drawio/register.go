package drawio

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
)

func init() {
	conv.Register(&conv.Format{
		Name:        "drawio",
		Description: "draw.io diagram",
		// .drawio.svg and .drawio.png exports are told apart by their content
		Extensions: []string{".drawio", ".dio"},
		Params:     []conv.Param{{Name: "border", Usage: "margin around each page's drawing in pixels (default 10; 0 for none)"}},
		Detect:     isDrawio,
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts := &Options{Pages: o.Pages, Title: o.Title, Images: o.Images,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn}
			if v := o.Param("border"); v != "" {
				b, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return nil, fmt.Errorf("drawio: parameter border: %q is not a number", v)
				}
				if b <= 0 {
					b = -1 // Options.Border: negative for none, zero for the default
				}
				opts.Border = b
			}
			data := make([]byte, size)
			if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
				return nil, fmt.Errorf("drawio: %w", err)
			}
			res, err := Convert(data, opts)
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d page(s), %d embedded font(s)", res.Pages, res.EmbeddedFonts)}, nil
		},
	})
}

// maxSniff is the most read to find a diagram in an SVG or PNG export.
const maxSniff = 64 << 20

// isDrawio reports whether an input is a draw.io file: XML with an mxfile
// or mxGraphModel root, an SVG export whose root element carries the
// diagram in its content attribute, or a PNG export with an mxfile (or
// mxGraphModel) text chunk.
func isDrawio(head []byte, r io.ReaderAt, size int64) bool {
	t := bytes.TrimLeft(head, "\xef\xbb\xbf \t\r\n")
	switch {
	case bytes.HasPrefix(t, []byte("<")):
		if bytes.Contains(t, []byte("<mxfile")) || bytes.Contains(t, []byte("<mxGraphModel")) {
			return true
		}
		i := bytes.Index(t, []byte("<svg"))
		if i < 0 {
			return false
		}
		// the root tag may be longer than the head: read up to its end
		tag := t[i:]
		if !bytes.Contains(tag, []byte(">")) && size <= maxSniff {
			all := make([]byte, size)
			n, _ := r.ReadAt(all, 0)
			all = all[:n]
			if j := bytes.Index(all, []byte("<svg")); j >= 0 {
				tag = all[j:]
			}
		}
		if end := bytes.IndexByte(tag, '>'); end >= 0 {
			tag = tag[:end]
		}
		return bytes.Contains(tag, []byte(" content=\"")) &&
			(bytes.Contains(tag, []byte("&lt;mxfile")) || bytes.Contains(tag, []byte("&lt;mxGraphModel")))
	case bytes.HasPrefix(head, []byte("\x89PNG\r\n\x1a\n")):
		// text chunks come before the image data; look through the chunk list
		off := int64(8)
		var hdr [8]byte
		for off+8 <= size {
			if _, err := r.ReadAt(hdr[:], off); err != nil {
				return false
			}
			n := int64(uint32(hdr[0])<<24 | uint32(hdr[1])<<16 | uint32(hdr[2])<<8 | uint32(hdr[3]))
			typ := string(hdr[4:8])
			if typ == "tEXt" || typ == "zTXt" {
				key := make([]byte, 13)
				m, _ := r.ReadAt(key, off+8)
				key = key[:m]
				if bytes.HasPrefix(key, []byte("mxfile\x00")) || bytes.HasPrefix(key, []byte("mxGraphModel\x00")) {
					return true
				}
			}
			if typ == "IEND" || typ == "IDAT" {
				return false
			}
			off += 12 + n
		}
	}
	return false
}

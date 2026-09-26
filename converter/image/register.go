package image

import (
	"bytes"
	"io"

	"github.com/shibukawa/bdf/converter"
)

func init() {
	converter.Register(&converter.Format{
		Name:        "image",
		Description: "Image stored as it is: PNG, JPEG, GIF, WebP, AVIF, BMP, ICO or SVG",
		Extensions:  []string{".png", ".apng", ".jpg", ".jpeg", ".jpe", ".jfif", ".gif", ".webp", ".avif", ".bmp", ".ico", ".svg"},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			return sniff(head, r, size) != ""
		},
		// draw.io's PNG and SVG exports are images too
		Fallback: true,
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			res, err := Convert(r, size, &Options{Title: o.Title, FileName: o.FileName, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: res.Summary()}, nil
		},
	})
}

// maxSniff bounds what is read to find the root element of an SVG file
// that starts with long comments or a document type.
const maxSniff = 64 << 10

// sniff returns the format of an image from its first bytes, or "".
func sniff(head []byte, r io.ReaderAt, size int64) string {
	switch {
	case bytes.HasPrefix(head, []byte("\x89PNG\r\n\x1a\n")):
		return "png"
	case bytes.HasPrefix(head, []byte("\xff\xd8\xff")):
		return "jpeg"
	case bytes.HasPrefix(head, []byte("GIF87a")), bytes.HasPrefix(head, []byte("GIF89a")):
		return "gif"
	case len(head) >= 12 && string(head[:4]) == "RIFF" && string(head[8:12]) == "WEBP":
		return "webp"
	case len(head) >= 16 && string(head[4:8]) == "ftyp":
		// the major brand, then the compatible ones
		n := min(len(head), int(be.Uint32(head)))
		for i := 8; i+4 <= n; i += 4 {
			if b := string(head[i : i+4]); i != 12 && (b == "avif" || b == "avis") {
				return "avif"
			}
		}
		return ""
	case bmpHeader(head):
		return "bmp"
	case icoHeader(head, size):
		return "ico"
	}
	isSVG, known := svgRoot(head)
	if !known && size > int64(len(head)) {
		more := make([]byte, min(size, maxSniff))
		n, _ := r.ReadAt(more, 0)
		isSVG, _ = svgRoot(more[:n])
	}
	if isSVG {
		return "svg"
	}
	return ""
}

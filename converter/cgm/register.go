package cgm

import (
	"bytes"
	"fmt"
	"io"

	"github.com/klauspost/compress/gzip"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
)

// Detect reports whether an input is a metafile in the binary or clear
// text encoding, or one compressed with gzip.
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	if bytes.HasPrefix(head, []byte{0x1f, 0x8b}) {
		zr, err := gzip.NewReader(io.NewSectionReader(r, 0, size))
		if err != nil {
			return false
		}
		b := make([]byte, 1024)
		n, _ := io.ReadFull(zr, b)
		head = b[:n]
	}
	return encodingOf(head) != ""
}

func init() {
	conv.Register(&conv.Format{
		Name:        "cgm",
		Description: "Computer Graphics Metafile (CGM, binary or clear text)",
		Extensions:  []string{".cgm", ".cgz"},
		Detect:      Detect,
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			res, err := Convert(r, size, &Options{Title: o.Title, Pages: o.Pages, Images: o.Images, FontFS: o.FontFS,
				FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d page(s) of %d picture(s), %s encoding, %d embedded font(s)", res.Pages, res.Pictures, res.Encoding, res.EmbeddedFonts)}, nil
		},
	})
}

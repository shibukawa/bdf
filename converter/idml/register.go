package idml

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "idml",
		Description: "InDesign document (IDML)",
		Extensions:  []string{".idml"},
		Files:       "the images the document places as links, by file name",
		// a zip package whose mimetype names IDML, or that has a design map
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
				return false
			}
			if bytes.Contains(head, []byte(mimetype)) {
				return true
			}
			p, err := ooxml.Open(r, size)
			return err == nil && p.Has("designmap.xml")
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			res, err := Convert(r, size, &Options{Pages: o.Pages, Title: o.Title, Images: o.Images,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn, Files: o.Files, Dir: o.Dir})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d page(s), %d embedded font(s)", res.Pages, res.EmbeddedFonts)}, nil
		},
	})
}

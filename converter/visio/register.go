package visio

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "visio",
		Description: "Visio drawing",
		Extensions:  []string{".vsdx", ".vsdm", ".vstx", ".vstm", ".vdx", ".vtx"},
		// a package with a Visio document part, or an XML drawing
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			if bytes.HasPrefix(head, []byte("PK\x03\x04")) {
				p, err := ooxml.Open(r, size)
				return err == nil && p.Has("visio/document.xml")
			}
			return bytes.Contains(head, []byte("<VisioDocument")) && bytes.Contains(head, []byte("schemas.microsoft.com/visio/2003/core"))
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			res, err := Convert(r, size, &Options{Pages: o.Pages, Title: o.Title, Images: o.Images,
				FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d page(s), %d embedded font(s)", res.Pages, res.EmbeddedFonts)}, nil
		},
	})
}

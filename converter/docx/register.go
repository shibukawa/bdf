package docx

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "docx",
		Description: "Word document",
		Extensions:  []string{".docx", ".docm", ".dotx", ".dotm"},
		Params:      []conv.Param{{Name: "views", Usage: "both (default): pages and scroll views; pages: only the pages; scroll: only the one long column"}},
		// an Office Open XML package with a word processing document part
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
				return false
			}
			p, err := ooxml.Open(r, size)
			if err != nil {
				return false
			}
			if rel, ok := p.RelOfType("", "/officeDocument"); ok && p.Has(rel.Target) {
				n, err := p.XML(rel.Target)
				return err == nil && n.Name == "document"
			}
			return p.Has("word/document.xml")
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			res, err := Convert(r, size, &Options{Pages: o.Pages, Views: o.Param("views"), Title: o.Title, Images: o.Images,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d page(s), %d scroll strip(s), %d embedded font(s)", res.Pages, res.Strips, res.EmbeddedFonts)}, nil
		},
	})
}

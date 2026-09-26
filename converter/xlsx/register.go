package xlsx

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "xlsx",
		Description: "Excel workbook",
		Extensions:  []string{".xlsx", ".xlsm", ".xltx", ".xltm"},
		Params:      []conv.Param{{Name: "hidden", Usage: "true: include hidden sheets"}},
		// an Office Open XML package with a workbook part
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
				return false
			}
			p, err := ooxml.Open(r, size)
			return err == nil && p.Has("xl/workbook.xml")
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			hidden, err := o.BoolParam("hidden")
			if err != nil {
				return nil, fmt.Errorf("xlsx: %w", err)
			}
			res, err := Convert(r, size, &Options{Sheets: o.Pages, Hidden: hidden, Title: o.Title, Images: o.Images,
				FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d sheet(s), %d embedded font(s)", res.Sheets, res.EmbeddedFonts)}, nil
		},
	})
}

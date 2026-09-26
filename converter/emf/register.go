package emf

import (
	"fmt"
	"io"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/metafile"
)

func init() {
	converter.Register(&converter.Format{
		Name:        "emf",
		Description: "Windows metafile",
		Extensions:  []string{".emf", ".wmf"},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			return metafile.Kind(head) != ""
		},
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			res, err := Convert(r, size, &Options{Title: o.Title, Images: o.Images, FontDirs: o.FontDirs,
				NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2,
				IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("1 page (%.0f × %.0f pt), %d embedded font(s)", res.W, res.H, res.EmbeddedFonts)}, nil
		},
	})
}

package jww

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
)

func init() {
	conv.Register(&conv.Format{
		Name:        "jww",
		Description: "Jw_cad drawing",
		Extensions:  []string{".jww"},
		Params: []conv.Param{
			{Name: "colors", Usage: "screen (default: the screen colors of the file), print (its printer colors on white) or mono (black on white)"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return bytes.HasPrefix(head, []byte("JwwData.")) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			var colors Colors
			switch v := o.Param("colors"); v {
			case "", "screen":
			case "print":
				colors = Print
			case "mono":
				colors = Mono
			default:
				return nil, fmt.Errorf("parameter colors: %q is not screen, print or mono", v)
			}
			res, err := Convert(r, size, &Options{Title: o.Title, Colors: colors, FontFS: o.FontFS, FontDirs: o.FontDirs,
				NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2,
				IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("1 page (%.0f × %.0f mm), %d embedded font(s)", res.W, res.H, res.EmbeddedFonts)}, nil
		},
	})
}

package dxf

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
)

func init() {
	conv.Register(&conv.Format{
		Name:        "dxf",
		Description: "AutoCAD DXF drawing",
		Extensions:  []string{".dxf"},
		Params: []conv.Param{
			{Name: "views", Usage: "all (default), model or layouts"},
			{Name: "background", Usage: "the background of model space: dark (default) or light"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return Detect(head) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			views := o.Param("views")
			switch views {
			case "", "all", "model", "layouts":
			default:
				return nil, fmt.Errorf("parameter views: %q is not all, model or layouts", views)
			}
			light := false
			switch bg := o.Param("background"); bg {
			case "", "dark":
			case "light":
				light = true
			default:
				return nil, fmt.Errorf("parameter background: %q is not dark or light", bg)
			}
			res, err := Convert(r, size, &Options{Pages: o.Pages, Title: o.Title, Views: views, Light: light,
				FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
				NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d view(s), %d embedded font(s)", res.Views, res.EmbeddedFonts)}, nil
		},
	})
}

// Detect reports whether the first bytes of a file are those of a DXF
// file: the binary sentinel, or a text file that starts with a comment or
// a section.
func Detect(head []byte) bool {
	if bytes.HasPrefix(head, binarySentinel) {
		return true
	}
	head = bytes.TrimPrefix(head, []byte("\xef\xbb\xbf"))
	lines := strings.SplitN(string(head), "\n", 5)
	for i := 0; i+1 < len(lines) && i < 4; i += 2 {
		code := strings.TrimSpace(lines[i])
		val := strings.TrimSpace(lines[i+1])
		switch code {
		case "999":
			continue
		case "0":
			return val == "SECTION"
		}
		return false
	}
	return false
}

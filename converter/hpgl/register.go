package hpgl

import (
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/hpglsniff"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "hpgl",
		Description: "HP-GL/2 plot file (HP-GL, HP RTL)",
		Extensions:  []string{".plt", ".hpgl", ".hgl", ".hpg", ".hp2", ".rtl"},
		Params: []conv.Param{
			{Name: "colors", Usage: "pens (default: the pen colors of the file) or mono (black on white)"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return hpglsniff.Is(head) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			var colors Colors
			switch v := o.Param("colors"); v {
			case "", "pens":
			case "mono":
				colors = Mono
			default:
				return nil, fmt.Errorf("parameter colors: %q is not pens or mono", v)
			}
			res, err := Convert(r, size, &Options{Title: o.Title, Colors: colors, Pages: o.Pages, Images: o.Images,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts,
				NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: summary(res)}, nil
		},
	})
}

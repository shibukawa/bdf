package musicxml

import (
	"io"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/music"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "musicxml",
		Description: "MusicXML score",
		// .xml is too common to name MusicXML; Detect recognizes such files
		// by their content
		Extensions: []string{".musicxml", ".mxl"},
		Params:     []conv.Param{{Name: "tab", Usage: "add guitar TAB views (true by default; false to omit)"}},
		Detect:     Detect,
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			tab, err := music.TabEnabled(o)
			if err != nil {
				return nil, err
			}
			score, warnings, err := Parse(r, size, nil)
			if err != nil {
				return nil, err
			}
			if o.Warn != nil {
				for _, w := range warnings {
					o.Warn(w)
				}
				warnings = nil
			}
			res, err := music.Build(score, music.Options{Source: "musicxml", Title: o.Title, Pages: o.Pages,
				GuitarTAB: tab,
				FontFS:    o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts,
				NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex,
				Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: append(warnings, res.Warnings...), Summary: res.Summary()}, nil
		},
	})
}

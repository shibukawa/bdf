package pptx

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "pptx",
		Description: "PowerPoint presentation",
		Extensions:  []string{".pptx", ".pptm", ".ppsx", ".ppsm", ".potx", ".potm"},
		Params: []conv.Param{
			{Name: "hidden", Usage: "true: include hidden slides"},
			{Name: "data", Usage: conv.ChartDataParamUsage},
		},
		// an Office Open XML package with a presentation part
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
				return false
			}
			p, err := ooxml.Open(r, size)
			return err == nil && p.Has("ppt/presentation.xml")
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			po, err := OptionsOf(o)
			if err != nil {
				return nil, err
			}
			res, err := Convert(r, size, po)
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d slide(s), %d embedded font(s)", res.Slides, res.EmbeddedFonts)}, nil
		},
	})
}

// OptionsOf returns the options of this package for the options of a
// conversion by the registry (nil for nil), with the parameter "hidden".
func OptionsOf(o *conv.Options) (*Options, error) {
	if o == nil {
		return nil, nil
	}
	hidden, err := o.BoolParam("hidden")
	if err != nil {
		return nil, fmt.Errorf("pptx: %w", err)
	}
	data, err := conv.DataParam(o.Param("data"))
	if err != nil {
		return nil, fmt.Errorf("pptx: %w", err)
	}
	return &Options{Slides: o.Pages, Hidden: hidden, Title: o.Title, Images: o.Images, Data: data,
		FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
		NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn}, nil
}

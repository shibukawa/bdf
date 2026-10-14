package xlsx

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

func init() {
	conv.Register(&conv.Format{
		Name:        "xlsx",
		Description: "Excel workbook",
		Extensions:  []string{".xlsx", ".xlsm", ".xltx", ".xltm"},
		Params: []conv.Param{
			{Name: "hidden", Usage: "true: include hidden sheets"},
			{Name: "views", Usage: "sheet (default): the sheets as Excel shows them; print: the pages Excel prints (print areas, page setup, print titles, header and footer); both"},
			{Name: "paper", Usage: "the paper of printed pages whose sheet names none: " + strings.Join(PaperNames(), ", ") + " (default: a4 for a Japanese workbook, else letter)"},
		},
		// an Office Open XML package with a workbook part
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
				return false
			}
			p, err := ooxml.Open(r, size)
			return err == nil && p.Has("xl/workbook.xml")
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			opts, err := OptionsOf(o)
			if err != nil {
				return nil, err
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d sheet(s), %d embedded font(s)", res.Sheets, res.EmbeddedFonts)}, nil
		},
	})
}

// OptionsOf returns the options of this package for the options of a
// conversion by the registry (nil for nil): those that apply to workbooks,
// with the parameters "hidden", "views" and "paper".
func OptionsOf(o *conv.Options) (*Options, error) {
	if o == nil {
		return nil, nil
	}
	hidden, err := o.BoolParam("hidden")
	if err != nil {
		return nil, fmt.Errorf("xlsx: %w", err)
	}
	return &Options{Sheets: o.Pages, Hidden: hidden, Title: o.Title, Images: o.Images,
		FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset,
		NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn,
		Views: o.Params["views"], Paper: o.Params["paper"], FileName: filepath.Base(o.FileName)}, nil
}

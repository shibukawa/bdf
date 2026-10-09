package kicad

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ziputil"

	"github.com/shibukawa/bdf/converter"
)

// Detect reports whether an input is a KiCad schematic or board of KiCad 6
// or later, or a zip archive holding one.
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	h := bytes.TrimLeft(bytes.TrimPrefix(head, []byte("\xef\xbb\xbf")), " \t\r\n")
	if bytes.HasPrefix(h, []byte("(kicad_sch")) || bytes.HasPrefix(h, []byte("(kicad_pcb")) {
		return true
	}
	if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		return false
	}
	zr, err := ziputil.NewReader(r, size)
	if err != nil {
		return false
	}
	for _, f := range zr.File {
		switch strings.ToLower(path.Ext(f.Name)) {
		case ".kicad_sch", ".kicad_pcb", ".kicad_pro":
			if !strings.HasPrefix(f.Name, "__MACOSX/") {
				return true
			}
		}
	}
	return false
}

func init() {
	converter.Register(&converter.Format{
		Name:        "kicad",
		Description: "KiCad schematic or board (KiCad 6 and later), or a zip archive of a project",
		Extensions:  []string{".kicad_sch", ".kicad_pcb", ".kicad_pro"},
		Refines:     "",
		Files:       "the sheets of a hierarchical schematic and the project file (.kicad_pro), by their paths relative to the input",
		Params: []converter.Param{
			{Name: "views", Usage: "a project: all (default: the schematic and the board), schematic or board"},
			{Name: "layers", Usage: "a board: all (default: its front, its back and a view of each layer), board or layers"},
		},
		Detect: Detect,
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			layers := o.Param("layers")
			switch layers {
			case "", "all", "board", "layers":
			default:
				return nil, fmt.Errorf("parameter layers: %q is not all, board or layers", layers)
			}
			opts := &Options{Pages: o.Pages, Title: o.Title, Views: o.Param("views"), Layers: layers, FileName: o.FileName,
				FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts,
				NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn}
			if o.Files != nil || o.Dir != "" {
				opts.Refs = o.ReadRef
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			var parts []string
			if res.Sheets > 0 {
				parts = append(parts, fmt.Sprintf("%d sheet(s)", res.Sheets))
			}
			if res.Boards > 0 {
				parts = append(parts, "a board")
			}
			parts = append(parts, fmt.Sprintf("%d embedded font(s)", res.EmbeddedFonts))
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: strings.Join(parts, ", ")}, nil
		},
	})
}

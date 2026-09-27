package gerber

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter"
)

// Detect reports whether an input is a Gerber file, an Excellon drill file,
// or a zip archive holding them (and not an Office document or an SXF
// drawing, which are zip archives too).
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	if bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		zr, err := zip.NewReader(r, size)
		if err != nil {
			return false
		}
		checked := 0
		for _, f := range zr.File {
			name := strings.ToLower(f.Name)
			switch {
			case name == "[content_types].xml" || strings.HasSuffix(name, ".p21") || strings.HasSuffix(name, ".sfc"):
				return false
			}
		}
		for _, f := range zr.File {
			if f.FileInfo().IsDir() || strings.HasPrefix(f.Name, "__MACOSX/") || checked >= 64 {
				continue
			}
			checked++
			rc, err := f.Open()
			if err != nil {
				continue
			}
			b, _ := io.ReadAll(io.LimitReader(rc, 64<<10))
			rc.Close()
			if s := sniff(f.Name, b); s == "gerber" || s == "excellon" {
				return true
			}
		}
		return false
	}
	if size > int64(len(head)) && !IsExcellon(head) {
		// the format specification may come after a long header of comments
		head = make([]byte, min(size, 64<<10))
		n, _ := r.ReadAt(head, 0)
		head = head[:n]
	}
	return IsGerber(head) || IsExcellon(head)
}

func init() {
	converter.Register(&converter.Format{
		Name:        "gerber",
		Description: "PCB fabrication data (Gerber RS-274X, Excellon drill, and zip archives of them)",
		Extensions: []string{".gbr", ".ger", ".gtl", ".gbl", ".gto", ".gbo", ".gts", ".gbs", ".gtp", ".gbp", ".gko", ".gml", ".gm1",
			".g1", ".g2", ".g3", ".g4", ".gp1", ".gp2", ".pho", ".art", ".cmp", ".sol", ".plc", ".pls", ".stc", ".sts", ".crc", ".crs",
			".drl", ".drd", ".xln", ".exc", ".zip"},
		Params: []converter.Param{
			{Name: "views", Usage: "all (default: the top and bottom of the board and a view for each file), board or layers"},
			{Name: "mask", Usage: "solder mask color: green, red, blue, black, white, yellow, purple or #rrggbb (default: the job file's, else green)"},
			{Name: "silkscreen", Usage: "silkscreen color: white, black, yellow or #rrggbb (default: the job file's, else white)"},
			{Name: "finish", Usage: "exposed copper: gold, silver or copper (default: the job file's finish, else gold)"},
		},
		Detect: Detect,
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			opts := &Options{Title: o.Title, FileName: path.Base(o.FileName), Mask: o.Param("mask"),
				Silkscreen: o.Param("silkscreen"), Finish: o.Param("finish"), Warn: o.Warn}
			if o.FileName == "" {
				opts.FileName = ""
			}
			switch v := o.Param("views"); v {
			case "", "all":
			case "board":
				opts.Views = ViewsBoard
			case "layers":
				opts.Views = ViewsLayers
			default:
				return nil, fmt.Errorf("parameter views: %q is not all, board or layers", v)
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			views := fmt.Sprintf("%d layer(s)", res.Layers)
			if res.Board {
				views = "the top and bottom of the board, " + views
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%s (%.1f × %.1f mm, drawn at %s:1)", views, res.W, res.H, strconv.FormatFloat(res.Scale, 'g', 3, 64))}, nil
		},
	})
}

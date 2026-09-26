package sxf

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"

	conv "github.com/shibukawa/bdf/converter"
)

// Detect reports whether an input is an SXF drawing: a STEP file of the
// associative draughting schema (P21 or SFC), or a zip archive holding a
// P21 or SFC file (P2Z).
func Detect(head []byte, r io.ReaderAt, size int64) bool {
	if bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		zr, err := zip.NewReader(r, size)
		if err != nil {
			return false
		}
		for _, f := range zr.File {
			if name := strings.ToLower(f.Name); strings.HasSuffix(name, ".p21") || strings.HasSuffix(name, ".sfc") {
				return true
			}
		}
		return false
	}
	text := bytes.TrimLeft(bytes.TrimPrefix(head, []byte("\xef\xbb\xbf")), " \t\r\n")
	if !bytes.HasPrefix(text, []byte("ISO-10303-21")) {
		return false
	}
	// the schema is in the header, which may be longer than head
	if size > int64(len(head)) {
		head = make([]byte, min(size, 64<<10))
		n, _ := r.ReadAt(head, 0)
		head = head[:n]
	}
	return bytes.Contains(bytes.ToUpper(head), []byte("ASSOCIATIVE_DRAUGHTING"))
}

func init() {
	conv.Register(&conv.Format{
		Name:        "sxf",
		Description: "SXF drawing (P21, P2Z, SFC)",
		Extensions:  []string{".p21", ".p2z", ".sfc"},
		Params: []conv.Param{
			{Name: "background", Usage: "file (default: the background the drawing names, black when it does not) or light (white paper)"},
		},
		Detect: Detect,
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			light := false
			switch v := o.Param("background"); v {
			case "", "file":
			case "light":
				light = true
			default:
				return nil, fmt.Errorf("parameter background: %q is not file or light", v)
			}
			res, err := Convert(r, size, &Options{Title: o.Title, Light: light, FontFS: o.FontFS, FontDirs: o.FontDirs,
				NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2,
				IgnoreFSType: o.IgnoreFSType, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("1 page (%.0f × %.0f mm, %s), %d embedded font(s)", res.W, res.H, strings.ToUpper(res.Format), res.EmbeddedFonts)}, nil
		},
	})
}

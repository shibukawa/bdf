package ai

import (
	"fmt"
	"io"

	"github.com/shibukawa/bdf/converter"
)

func init() {
	converter.Register(&converter.Format{
		Name:        "ai",
		Description: "Adobe Illustrator artwork",
		Extensions:  []string{".ai"},
		Refines:     "pdf",
		Params: []converter.Param{
			{Name: "box", Usage: "page boundary: trim (default, the artboard), bleed, media, crop or art"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			return Kind(r, size) != ""
		},
		CheckPassword: func(r io.ReaderAt, size int64, password string) (bool, error) {
			if f := converter.Lookup("pdf"); f != nil && Kind(r, size) == "pdf" {
				return f.CheckPassword(r, size, password)
			}
			return false, nil
		},
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			res, err := Convert(r, size, &Options{Pages: o.Pages, Title: o.Title, Box: o.Param("box"),
				NoTextIndex: o.NoTextIndex, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType,
				Images: o.Images, Password: o.Password, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings, Protected: res.Protected,
				Summary: fmt.Sprintf("%d artboard(s)", res.Artboards)}, nil
		},
	})
}

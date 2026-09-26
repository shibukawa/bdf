package pdf

import (
	"bytes"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
)

func init() {
	conv.Register(&conv.Format{
		Name:        "pdf",
		Description: "PDF",
		Extensions:  []string{".pdf"},
		Params: []conv.Param{
			{Name: "kind", Usage: "view kind: fixed (default) or flow"},
			{Name: "no-share", Usage: "true: do not move the instruction prefix pages have in common into a shared object"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			return bytes.Contains(head, []byte("%PDF-"))
		},
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			noShare, err := o.BoolParam("no-share")
			if err != nil {
				return nil, fmt.Errorf("pdf: %w", err)
			}
			res, err := Convert(io.NewSectionReader(r, 0, size), &Options{Pages: o.Pages, Title: o.Title, Kind: o.Param("kind"),
				NoTextIndex: o.NoTextIndex, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType,
				Images: o.Images, NoSharePrefix: noShare, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			return &conv.Result{Doc: res.Doc, Warnings: res.Warnings,
				Summary: fmt.Sprintf("%d page(s), %d shared prefix(es) saving %d bytes", res.Pages, res.SharedPrefixes, res.SharedBytes)}, nil
		},
	})
}

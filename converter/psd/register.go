package psd

import (
	"fmt"
	"io"

	"github.com/shibukawa/bdf/converter"
)

func init() {
	converter.Register(&converter.Format{
		Name:        "psd",
		Description: "Photoshop document",
		Extensions:  []string{".psd", ".psb"},
		Params: []converter.Param{
			{Name: "artboards", Usage: "false: the whole canvas on one page instead of a page per artboard"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			return IsPSD(head)
		},
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			noArtboards := false
			if o.Param("artboards") != "" {
				on, err := o.BoolParam("artboards")
				if err != nil {
					return nil, fmt.Errorf("psd: %w", err)
				}
				noArtboards = !on
			}
			res, err := Convert(r, size, &Options{Pages: o.Pages, Title: o.Title, NoArtboards: noArtboards,
				Images: o.Images, NoTextIndex: o.NoTextIndex, Warn: o.Warn})
			if err != nil {
				return nil, err
			}
			summary := fmt.Sprintf("%d page(s)", res.Pages)
			if res.Artboards > 0 {
				summary = fmt.Sprintf("%d artboard(s)", res.Artboards)
			}
			if res.Composited {
				summary += ", layers composited by the converter"
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: summary}, nil
		},
	})
}

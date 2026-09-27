package tiff

import (
	"fmt"
	"io"
	"strconv"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/tiff"
)

func init() {
	converter.Register(&converter.Format{
		Name:        "tiff",
		Description: "TIFF image, one page or more",
		Extensions:  []string{".tif", ".tiff"},
		Params: []converter.Param{
			{Name: "dpi", Usage: fmt.Sprintf("resolution in pixels per inch of pages that do not give one (default %d)", DefaultDPI)},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool {
			// Canon CR2 raw files are TIFF files too.
			return tiff.Sniff(head) && !(len(head) >= 10 && string(head[8:10]) == "CR")
		},
		Convert: func(r io.ReaderAt, size int64, o *converter.Options) (*converter.Result, error) {
			opts := &Options{Title: o.Title, Pages: o.Pages, Images: o.Images, Warn: o.Warn}
			if v := o.Param("dpi"); v != "" {
				dpi, err := strconv.ParseFloat(v, 64)
				if err != nil || !(dpi > 0) {
					return nil, fmt.Errorf("parameter dpi: %q is not a positive number", v)
				}
				opts.DPI = dpi
			}
			res, err := Convert(r, size, opts)
			if err != nil {
				return nil, err
			}
			summary := fmt.Sprintf("%d page(s)", res.Pages)
			if res.Scaled > 0 {
				summary += fmt.Sprintf(", %d scaled down to the resolution cap", res.Scaled)
			}
			return &converter.Result{Doc: res.Doc, Warnings: res.Warnings, Summary: summary}, nil
		},
	})
}

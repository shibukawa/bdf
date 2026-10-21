// Package print prints BDF documents on printers: it picks what to send
// (the pages as a PDF file or a PostScript file, or as PWG Raster, Apple
// Raster or PCLm drawn by raster/imagebdf),
// turns a Ticket of what the job asks into the attributes of IPP, and
// sends the job to a printer or a queue with package print/ipp, waiting
// for it to be done.
//
// The name of the package is that of the builtin print function, which
// it hides only where it is imported.
package print

import (
	"errors"
	"fmt"

	"github.com/shibukawa/bdf"
)

// ErrNoPages is returned for a document without pages to print.
var ErrNoPages = errors.New("print: the document has no pages to print")

// PrintView returns the view of a document that is printed: the view of
// the id, or the first view of pages (fixed or flow), or the first
// scroll view, whose strips are pages too.
func PrintView(doc *bdf.Document, id string) (*bdf.View, error) {
	if id != "" {
		for _, v := range doc.Views {
			if v.ID == id {
				if v.Kind == bdf.ViewSheet {
					return nil, fmt.Errorf("print: view %q is a sheet: print a view of pages (an Excel workbook converted with -param views=print has one)", id)
				}
				if len(v.Pages) == 0 {
					return nil, fmt.Errorf("print: view %q has no pages", id)
				}
				return v, nil
			}
		}
		return nil, fmt.Errorf("print: no view %q", id)
	}
	sheets := false
	for _, kinds := range [][]string{{bdf.ViewFixed, bdf.ViewFlow}, {bdf.ViewScroll}} {
		for _, v := range doc.Views {
			sheets = sheets || v.Kind == bdf.ViewSheet
			for _, k := range kinds {
				if v.Kind == k && len(v.Pages) > 0 {
					return v, nil
				}
			}
		}
	}
	if sheets {
		return nil, fmt.Errorf("%w: its views are sheets, which have no pages (an Excel workbook has its print view with -param views=print)", ErrNoPages)
	}
	return nil, ErrNoPages
}

// pageNumbers returns the zero-based pages of a view that a selection
// of one-based pages names (nil for all), dropping the numbers past the
// view's pages.
func pageNumbers(v *bdf.View, pages []int) []int {
	var out []int
	if pages == nil {
		for i := range v.Pages {
			out = append(out, i)
		}
		return out
	}
	for _, n := range pages {
		if n >= 1 && n <= len(v.Pages) {
			out = append(out, n-1)
		}
	}
	return out
}

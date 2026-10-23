// Package render is the render node of docflow: it draws the .bdf
// document that reaches it as a PDF file (vector/pdfbdf), or a page of it
// as a PNG image (raster/imagebdf), for an exit whose far side has no
// viewer of BDF (a mail's reader, a web service). The file made is the
// job's payload from there on, of the kind any. A program that imports
// the package has the kind.
package render

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"path"
	"strings"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/vector/pdfbdf"
)

// KindRender is the kind of the render node.
const KindRender = "render"

// The formats a render node makes.
const (
	FormatPDF = "pdf"
	FormatPNG = "png"
)

// DPI is the resolution of a PNG image of a page.
const DPI = 150

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindRender, Role: docflow.RoleStep, Takes: []string{docflow.PayloadBDF}, Makes: docflow.PayloadAny,
		Check: check, Step: step})
}

func check(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.To = docflow.CleanNames(in.To)
	n.Format = strings.ToLower(strings.TrimSpace(in.Format))
	if n.Format == "" {
		n.Format = FormatPDF
	}
	switch n.Format {
	case FormatPDF:
	case FormatPNG:
		n.Page = in.Page
		if n.Page < 0 {
			return fmt.Errorf("%w: page %d: the first is 1", docflow.ErrInvalid, in.Page)
		}
	case "svg":
		return fmt.Errorf("%w: format svg: not yet (vector/svgbdf is not on main); pdf or png", docflow.ErrInvalid)
	default:
		return fmt.Errorf("%w: format %q: pdf or png", docflow.ErrInvalid, in.Format)
	}
	return nil
}

// step draws the document.
func step(ctx context.Context, o *docflow.Out) (*docflow.Made, error) {
	doc, err := o.Document(ctx)
	if err != nil {
		return nil, err
	}
	mo := o.Manager().Options()
	name := stem(o.Payload.Name) + "." + o.Node.Format
	var buf bytes.Buffer
	switch o.Node.Format {
	case FormatPNG:
		v, err := docflow.PagesView(doc)
		if err != nil {
			return nil, docflow.Final(err)
		}
		page := max(o.Node.Page, 1) - 1
		if page >= len(v.Pages) {
			return nil, docflow.Final(fmt.Errorf("page %d of %d", page+1, len(v.Pages)))
		}
		img, err := imagebdf.New(doc, &imagebdf.Options{FontDirs: mo.FontDirs, NoSystemFonts: mo.NoSystemFonts}).Page(v, page, DPI/72.0)
		if err != nil {
			return nil, docflow.Final(err)
		}
		if err := png.Encode(&buf, img); err != nil {
			return nil, err
		}
		return &docflow.Made{Kind: docflow.PayloadAny, Format: FormatPNG, Name: name, Pages: 1, Body: &buf}, nil
	default:
		res, err := pdfbdf.Write(&buf, doc, &pdfbdf.Options{FontDirs: mo.FontDirs, NoSystemFonts: mo.NoSystemFonts})
		if err != nil {
			return nil, docflow.Final(err)
		}
		pages := o.Payload.Pages
		if res != nil && res.Pages > 0 {
			pages = res.Pages
		}
		return &docflow.Made{Kind: docflow.PayloadAny, Format: FormatPDF, Name: name, Pages: pages, Body: &buf}, nil
	}
}

// stem is a file's name without its extension: "document" for none.
func stem(name string) string {
	base := path.Base(strings.ReplaceAll(name, "\\", "/"))
	base = strings.TrimSuffix(base, path.Ext(base))
	if base == "" || base == "." || base == "/" {
		return "document"
	}
	return base
}

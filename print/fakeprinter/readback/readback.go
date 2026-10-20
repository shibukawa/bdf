// Package readback turns what a fake printer received back into pages:
// the images of PWG Raster and Apple Raster as they are, PDF and PCLm
// drawn (converter/pdf, then raster/imagebdf), PostScript drawn by
// Ghostscript where it is, JPEG and PNG decoded, and the label, receipt
// and dot-matrix languages read as their printers print them (print/pdl):
// their dots, and the commands of the printer's own barcodes.
//
// It is the half of a fake printer that a real one does with paper: the
// tests compare what comes back with what bdf meant to print, and the
// farm (tools/fakeprinter) shows it.
package readback

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/pdl"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Page is a page a printer printed.
type Page struct {
	Image image.Image
	// DPI is the resolution of the image.
	DPI int
	// Symbols are the commands of the printer's own barcodes on the page
	// (labels): see pdl.PrintedPage.
	Symbols []string
	// Cut is "partial" when a receipt printer cut the paper after it.
	Cut string
}

// Result is what a job printed.
type Result struct {
	// Format is what the data was read as: a MIME type or a language of
	// print/pdl.
	Format string
	Pages  []Page
	// Commands are the commands of a label or dot-matrix language that
	// are not on a page (pdl.Printed's), and the lines of a PJL header.
	Commands []string
	// Drawer is the pin of the cash drawer a receipt printer opened.
	Drawer int
}

// Options say how to read.
type Options struct {
	// Format is what the data is; "" tells it from its first bytes.
	Format string
	// Width and DPI are the dots across and the resolution of a printer of
	// a label, receipt or dot-matrix language (which those do not say).
	Width, DPI int
	// RenderDPI is the resolution PDF, PCLm and PostScript are drawn at:
	// 150 when 0.
	RenderDPI int
	// FontDirs and NoSystemFonts are those of the conversion of a PDF.
	FontDirs      []string
	NoSystemFonts bool
	// Ghostscript is the gs that draws PostScript; "" looks for it.
	Ghostscript string
	// MaxPixels bounds a raster page (64 Mpx when 0).
	MaxPixels int
}

// ErrNoGhostscript is returned for PostScript where there is no gs.
var ErrNoGhostscript = errors.New("readback: PostScript is drawn by Ghostscript, which is not here")

// Job reads a job of a fake printer, as the printer's language, width and
// resolution say.
func Job(p *fakeprinter.Printer, j *fakeprinter.Job, o *Options) (*Result, error) {
	var opts Options
	if o != nil {
		opts = *o
	}
	if opts.Format == "" {
		opts.Format = j.Format
		if opts.Format == "" || opts.Format == "application/octet-stream" || opts.Format == "application/vnd.cups-raw" {
			opts.Format = p.Language
		}
	}
	if opts.Width == 0 {
		opts.Width, opts.DPI = p.Width, p.DPI
	}
	return Read(j.Data(), &opts)
}

// Read reads a document a printer was sent.
func Read(data []byte, o *Options) (*Result, error) {
	var opts Options
	if o != nil {
		opts = *o
	}
	res := &Result{}
	if lines, doc, ok := fakeprinter.PJL(data); ok {
		res.Commands, data = lines, doc
	}
	format := opts.Format
	if sniffed := fakeprinter.Sniff(data); format == "" || format == "application/octet-stream" || format == "application/vnd.cups-raw" ||
		format == "postscript" || format == "pdf" {
		format = sniffed
	}
	res.Format = format
	dpi := opts.RenderDPI
	if dpi <= 0 {
		dpi = 150
	}
	maxPixels := opts.MaxPixels
	if maxPixels <= 0 {
		maxPixels = 64 << 20
	}
	switch format {
	case "image/pwg-raster", "image/urf":
		read := raster.Read
		if format == "image/urf" {
			read = raster.ReadURF
		}
		pages, err := read(bytes.NewReader(data), maxPixels)
		if err != nil {
			return nil, err
		}
		for _, pg := range pages {
			res.Pages = append(res.Pages, Page{Image: pg.Image, DPI: int(pg.XDPI)})
		}
	case "application/pdf", "application/PCLm":
		pages, err := drawPDF(data, dpi, &opts)
		if err != nil {
			return nil, err
		}
		res.Pages = pages
	case "application/postscript":
		pages, err := drawPS(data, dpi, opts.Ghostscript)
		if err != nil {
			return nil, err
		}
		res.Pages = pages
	case "image/jpeg", "image/png":
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		res.Pages = []Page{{Image: img}}
	case pdl.ZPL, pdl.SBPL, pdl.TPCL, pdl.ESCPOS, pdl.ESCP:
		pr, err := pdl.Read(format, data, opts.Width, opts.DPI)
		if err != nil {
			return nil, err
		}
		for _, pg := range pr.Pages {
			res.Pages = append(res.Pages, Page{Image: Dots(pg.Dots), DPI: opts.DPI, Symbols: pg.Symbols, Cut: pg.Cut})
		}
		res.Commands = append(res.Commands, pr.Commands...)
		res.Drawer = pr.Drawer
	default:
		return nil, fmt.Errorf("readback: no reader of %q", format)
	}
	return res, nil
}

// Dots turns the dots of a label into an image: black dots black on
// white.
func Dots(b *pdl.Bitmap) *image.Gray {
	img := image.NewGray(image.Rect(0, 0, b.W, b.H))
	for y := 0; y < b.H; y++ {
		row := b.Row(y)
		for x := 0; x < b.W; x++ {
			if row[x>>3]&(0x80>>(x&7)) == 0 {
				img.Pix[y*img.Stride+x] = 255
			}
		}
	}
	return img
}

// drawPDF draws the pages of a PDF (or PCLm, a PDF of strips of images).
func drawPDF(data []byte, dpi int, o *Options) ([]Page, error) {
	res, err := pdf.Convert(bytes.NewReader(data), &pdf.Options{NoTextIndex: true})
	if err != nil {
		return nil, err
	}
	if len(res.Doc.Views) == 0 {
		return nil, errors.New("readback: a PDF of no pages")
	}
	v := res.Doc.Views[0]
	r := imagebdf.New(res.Doc, &imagebdf.Options{FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts})
	var pages []Page
	for i := range v.Pages {
		img, err := r.PageDPI(v, i, float64(dpi))
		if err != nil {
			return nil, err
		}
		pages = append(pages, Page{Image: img, DPI: dpi})
	}
	return pages, nil
}

// drawPS draws PostScript with Ghostscript.
func drawPS(data []byte, dpi int, gs string) ([]Page, error) {
	if gs == "" {
		var err error
		if gs, err = exec.LookPath("gs"); err != nil {
			return nil, ErrNoGhostscript
		}
	}
	dir, err := os.MkdirTemp("", "readback-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	in := filepath.Join(dir, "in.ps")
	if err := os.WriteFile(in, data, 0o600); err != nil {
		return nil, err
	}
	out, err := exec.Command(gs, "-q", "-dSAFER", "-dBATCH", "-dNOPAUSE", "-sDEVICE=png16m", fmt.Sprintf("-r%d", dpi),
		"-sOutputFile="+filepath.Join(dir, "p%04d.png"), in).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("readback: gs: %v: %s", err, strings.TrimSpace(string(out)))
	}
	names, _ := filepath.Glob(filepath.Join(dir, "p*.png"))
	sort.Strings(names)
	var pages []Page
	for _, n := range names {
		b, err := os.ReadFile(n)
		if err != nil {
			return nil, err
		}
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		pages = append(pages, Page{Image: img, DPI: dpi})
	}
	return pages, nil
}

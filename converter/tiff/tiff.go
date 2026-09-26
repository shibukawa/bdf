// Package tiff converts TIFF images into BDF documents: each page of the
// file (each image of the IFD chain that is not a reduced-resolution copy
// or a mask) becomes a page of a fixed view, drawn by one image the size
// of the page (see docs/design.md §3.10).
//
// The page size comes from the image's resolution. Pages finer than the
// resolution cap (imgconv.Options.MaxDPI and MaxPixels) are scaled down to
// it; the others keep their pixels. A JPEG page is stored as one JPEG
// stream made from its strips without re-encoding when it needs no
// scaling, and every other page is decoded and stored as PNG, JPEG or
// WebP (see imgconv.EncodePixels). The Orientation tag turns or mirrors
// the image on the page.
package tiff

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// DefaultDPI is the resolution assumed for pages that do not give one: a
// CSS pixel per image pixel at the page's real size.
const DefaultDPI = 96

// Options controls the conversion.
type Options struct {
	// Title overrides the document title (by default the DocumentName tag).
	Title string
	// Pages selects 1-based pages (see converter.Pages); nil converts all of
	// them.
	Pages converter.Pages
	// Images controls how pages are stored and the resolution cap they are
	// scaled down to (see imgconv). The zero value stores lossless pages as
	// PNG and JPEG pages as JPEG, capped at imgconv.DefaultMaxDPI and
	// imgconv.DefaultMaxPixels.
	Images imgconv.Options
	// DPI is the resolution assumed for pages that do not give one
	// (default DefaultDPI).
	DPI float64
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Pages is the number of pages converted, Scaled the number of those
	// scaled down to the resolution cap.
	Pages, Scaled int
}

// ConvertFile converts a TIFF file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Convert(f, st.Size(), opts)
}

// Convert converts a TIFF file read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	f, err := tiff.Open(r, size)
	if err != nil {
		return nil, err
	}
	if f.IFDs[0].Has(tiff.TagDNGVersion) {
		return nil, errors.New("tiff: DNG camera raw files are not supported")
	}
	pages := f.Pages()
	if len(pages) == 0 {
		return nil, errors.New("tiff: the file has no pages, only reduced-resolution images or masks")
	}
	res := &Result{Doc: bdf.NewDocument()}
	warned := map[string]bool{}
	warn := func(msg string) {
		if warned[msg] {
			return
		}
		warned[msg] = true
		if opts.Warn != nil {
			opts.Warn(msg)
		} else {
			res.Warnings = append(res.Warnings, msg)
		}
	}
	if f.ChainErr != nil {
		warn(fmt.Sprintf("%v: the pages after page %d are missing", f.ChainErr, len(pages)))
	}
	sel := opts.Pages.Numbers(len(pages))
	if sel == nil {
		for i := range pages {
			sel = append(sel, i+1)
		}
	}
	for _, n := range sel {
		if n < 1 || n > len(pages) {
			return nil, fmt.Errorf("tiff: page %d out of range (1-%d)", n, len(pages))
		}
	}
	doc := res.Doc
	doc.Meta.Source = "tiff"
	readMeta(&doc.Meta.DC, pages[0])
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	dpi := opts.DPI
	if !(dpi > 0) {
		dpi = DefaultDPI
	}
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	done := map[int]bool{}
	for _, n := range sel {
		if done[n] {
			continue
		}
		done[n] = true
		d := pages[n-1]
		p := newPage(d, dpi)
		if p.w <= 0 || p.h <= 0 {
			warn(fmt.Sprintf("page %d has no size; it is left out", n))
			continue
		}
		if d.Photometric() == tiff.PhotometricSeparated {
			warn("CMYK pages are converted to RGB without colour management")
		}
		img, scaled, err := storePage(doc, d, p, opts.Images, func(msg string) { warn(fmt.Sprintf("page %d: %s", n, msg)) })
		if err != nil {
			warn(fmt.Sprintf("page %d: %v; the page is left blank", n, err))
			view.AddPage(p.pageW, p.pageH)
			continue
		}
		if scaled {
			res.Scaled++
		}
		o := bdf.NewObject()
		ref := o.AddImage(img)
		o.SetBBox(0, 0, p.pageW, p.pageH)
		if t := p.transform(); t != [6]float32{1, 0, 0, 1, 0, 0} {
			o.Save().Transform(t[0], t[1], t[2], t[3], t[4], t[5]).Image(ref, 0, 0, p.imgW, p.imgH).Restore()
		} else {
			o.Image(ref, 0, 0, p.imgW, p.imgH)
		}
		h, _ := doc.AddObject(o)
		view.AddPage(p.pageW, p.pageH, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	res.Pages = len(view.Pages)
	return res, nil
}

// page is the geometry of a page.
type page struct {
	w, h         int     // pixels
	wPt, hPt     float64 // the image's size in points
	imgW, imgH   float32 // the same, as drawn
	pageW, pageH float32 // the page's size in points (the image's, turned by the orientation)
	orientation  int
}

func newPage(d *tiff.IFD, def float64) *page {
	p := &page{orientation: d.Orientation()}
	p.w, p.h = d.Size()
	xdpi, ydpi := resolution(d, def)
	p.wPt, p.hPt = float64(float64(p.w)/xdpi*72), float64(float64(p.h)/ydpi*72)
	p.imgW, p.imgH = float32(p.wPt), float32(p.hPt)
	p.pageW, p.pageH = p.imgW, p.imgH
	if p.orientation >= 5 {
		p.pageW, p.pageH = p.imgH, p.imgW
	}
	return p
}

// resolution returns the pixels per inch along x and y. Without a unit
// (ResolutionUnit 1) values that look like pixels per inch (ImageMagick
// writes a 300 dpi scan so) are taken as such, as tiff2pdf does, and
// smaller ones give only the pixels' aspect ratio.
func resolution(d *tiff.IFD, def float64) (x, y float64) {
	x, y, unit := d.Resolution()
	if x == 0 {
		x = y
	}
	if y == 0 {
		y = x
	}
	switch {
	case x == 0:
		return def, def
	case unit == tiff.ResUnitCentimeter:
		x, y = x*2.54, y*2.54
	case unit == tiff.ResUnitNone && min(x, y) < 50:
		return def, float64(def * y / x)
	}
	if x < 1 || y < 1 {
		return def, def
	}
	return x, y
}

// transform returns the matrix that puts the image, drawn at (0, 0) as it
// is stored, on the page as the Orientation tag says: the tag names where
// the stored image's first row and first column belong.
func (p *page) transform() [6]float32 {
	w, h := p.imgW, p.imgH
	switch p.orientation {
	case 2: // mirrored left to right
		return [6]float32{-1, 0, 0, 1, w, 0}
	case 3: // turned 180°
		return [6]float32{-1, 0, 0, -1, w, h}
	case 4: // mirrored top to bottom
		return [6]float32{1, 0, 0, -1, 0, h}
	case 5: // transposed
		return [6]float32{0, 1, 1, 0, 0, 0}
	case 6: // turned 90° clockwise to be upright
		return [6]float32{0, 1, -1, 0, h, 0}
	case 7: // transposed the other way
		return [6]float32{0, -1, -1, 0, h, w}
	case 8: // turned 90° counterclockwise to be upright
		return [6]float32{0, -1, 1, 0, 0, w}
	}
	return [6]float32{1, 0, 0, 1, 0, 0}
}

// storePage stores the image of a page, scaled down to the resolution cap
// when it is finer, and reports whether it was.
func storePage(doc *bdf.Document, d *tiff.IFD, p *page, img imgconv.Options, warn func(string)) (bdf.Hash, bool, error) {
	tw, th := img.FitSize(p.w, p.h, p.wPt, p.hPt)
	scale := tw != p.w || th != p.h
	if !scale {
		if b, ok := d.JPEG(); ok {
			r, _ := imgconv.Optimize(b, img) // on error r keeps the JPEG
			return doc.AddImage(r.Data), false, nil
		}
	}
	pix, damaged, err := d.Decode()
	if err != nil {
		return bdf.Hash{}, false, err
	}
	if damaged {
		warn("the pixel data is damaged; what is missing is left blank")
	}
	if scale {
		pix = imgconv.Resize(pix, tw, th)
	}
	r, _ := imgconv.EncodePixels(pix, d.Compression() != tiff.CompressionJPEG, img) // on error r is the PNG or JPEG
	return doc.AddImage(r.Data), scale, nil
}

// readMeta maps the descriptive tags of the first page to Dublin Core as
// XMP maps them (DocumentName, which XMP leaves out, is the title).
func readMeta(dc *bdf.DublinCore, d *tiff.IFD) {
	if s := d.Strings(tiff.TagDocumentName); len(s) > 0 {
		dc.Title = bdf.DCValues{s[0]}
	}
	if s := d.Strings(tiff.TagImageDescription); len(s) > 0 {
		dc.Description = bdf.DCValues{s[0]}
	}
	for _, s := range d.Strings(tiff.TagArtist) {
		for a := range strings.SplitSeq(s, ";") {
			if a = strings.TrimSpace(a); a != "" && !slices.Contains(dc.Creator, a) {
				dc.Creator = append(dc.Creator, a)
			}
		}
	}
	for _, s := range d.Strings(tiff.TagCopyright) {
		dc.Rights = append(dc.Rights, s)
	}
	if s := d.Strings(tiff.TagDateTime); len(s) > 0 {
		if m := dateTime.FindStringSubmatch(s[0]); m != nil && m[1] != "0000" {
			dc.Modified = bdf.DCValues{m[1] + "-" + m[2] + "-" + m[3] + "T" + m[4] + ":" + m[5] + ":" + m[6]}
		}
	}
}

// dateTime matches the DateTime tag, "YYYY:MM:DD HH:MM:SS".
var dateTime = regexp.MustCompile(`^(\d{4}):(\d{2}):(\d{2}) (\d{2}):(\d{2}):(\d{2})$`)

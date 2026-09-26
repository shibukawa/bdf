// Package psd converts Photoshop documents (.psd, and .psb, the large
// document format) into BDF documents: the image as Photoshop composited
// it, one page per artboard.
//
// A document without artboards becomes one page, the size of its canvas at
// its resolution (a 72 ppi image is one point per pixel). A document with
// artboards becomes a page per visible artboard, cut from the composite, in
// the order of the Layers panel from the bottom up: the order they were
// added in, unless they were moved.
//
// The composite is the image Photoshop stores beside the layers. A file
// saved without "Maximize Compatibility" has none, and the converter
// composites the layers itself (see composite.go for what it draws).
// Colour profiles are not applied.
package psd

import (
	"errors"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/xmp"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages (artboards); nil converts all of them.
	Pages []int
	// Title overrides the document title.
	Title string
	// NoArtboards puts the whole canvas on one page even when the document
	// has artboards.
	NoArtboards bool
	// Images controls how the page images are stored (see imgconv). The
	// zero value stores them as PNG.
	Images imgconv.Options
	// NoTextIndex skips building the (empty) text index part.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	Pages    int
	// Artboards counts the pages that are artboards (0: the canvas is the page).
	Artboards int
	// Composited reports that the file had no composite image and the
	// converter composited the layers.
	Composited bool
}

// IsPSD reports whether data starts like a Photoshop document.
func IsPSD(head []byte) bool {
	return len(head) >= 6 && string(head[:4]) == "8BPS" && head[4] == 0 && (head[5] == 1 || head[5] == 2)
}

// ConvertFile converts a .psd or .psb file.
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

// region is a page: a rectangle of the canvas.
type region struct {
	r        image.Rectangle
	artboard bool
}

// Convert converts a Photoshop document read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
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
	f, err := readFile(r, size)
	if errors.Is(err, errFormat) {
		return nil, fmt.Errorf("psd: not a Photoshop document")
	}
	if f == nil {
		return nil, fmt.Errorf("psd: %w", err)
	}
	if err != nil {
		warn(err.Error())
	}

	img, err := f.composite()
	if err != nil || img == nil {
		if len(f.layers) == 0 {
			if err == nil {
				err = fmt.Errorf("the file has no composite image and no layers")
			}
			return nil, fmt.Errorf("psd: %w", err)
		}
		if err != nil {
			warn(fmt.Sprintf("composite image: %v; compositing the layers instead", err))
		} else {
			warn(`the file was saved without "Maximize Compatibility"; the layers are composited by the converter`)
		}
		img = f.compositeLayers(warn)
		res.Composited = true
	}

	canvas := image.Rect(0, 0, f.hdr.w, f.hdr.h)
	var regions []region
	if !opts.NoArtboards {
		for _, n := range layerTree(f.layers) {
			if a := n.l.artboard; n.group && a != nil && !n.l.hidden {
				if rr := image.Rect(a.left, a.top, a.right, a.bottom).Intersect(canvas); !rr.Empty() {
					regions = append(regions, region{r: rr, artboard: true})
				}
			}
		}
	}
	if len(regions) == 0 {
		regions = []region{{r: canvas}}
	}
	pages := opts.Pages
	if pages == nil {
		for i := range regions {
			pages = append(pages, i+1)
		}
	}

	doc := res.Doc
	doc.Meta.Source = "psd"
	if b := f.resources[resXMP]; b != nil {
		doc.Meta.DC = xmp.DublinCore(b)
		doc.Meta.DC.Format = nil // the format of the source, not of this document
	}
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	xres, yres := f.resolution()
	for _, n := range pages {
		if n < 1 || n > len(regions) {
			return nil, fmt.Errorf("psd: page %d out of range (1-%d)", n, len(regions))
		}
		rg := regions[n-1]
		crop := image.NewNRGBA(image.Rect(0, 0, rg.r.Dx(), rg.r.Dy()))
		for y := 0; y < rg.r.Dy(); y++ {
			copy(crop.Pix[y*crop.Stride:(y+1)*crop.Stride], img.Pix[img.PixOffset(rg.r.Min.X, rg.r.Min.Y+y):])
		}
		enc, err := imgconv.EncodeImage(crop, true, opts.Images)
		if err != nil && enc.Data == nil {
			return nil, fmt.Errorf("psd: page %d: %w", n, err)
		}
		w, h := float32(float64(rg.r.Dx())*72/xres), float32(float64(rg.r.Dy())*72/yres)
		obj := bdf.NewObject()
		obj.SetBBox(0, 0, w, h)
		obj.Image(obj.AddImage(doc.AddImage(enc.Data)), 0, 0, w, h)
		hash, _ := doc.AddObject(obj)
		view.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
		if rg.artboard {
			res.Artboards++
		}
	}
	res.Pages = len(pages)
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			warn(fmt.Sprintf("text index: %v", err))
		}
	}
	return res, nil
}

// composite returns the composite image of the file; nil without an error
// when the file says it has none.
func (f *file) composite() (*image.NRGBA, error) {
	if !f.hasComposite() {
		return nil, nil
	}
	n := f.hdr.baseChannels()
	alpha := f.globalAlpha && f.hdr.channels > n && f.hdr.mode != modeMultichannel
	if alpha {
		n++
	}
	planes, err := f.readPlanes(n)
	if err != nil {
		return nil, err
	}
	var a []uint8
	if alpha {
		planes, a = planes[:n-1], planes[n-1]
		f.unmatte(planes, a)
	}
	if f.hdr.mode == modeMultichannel || f.hdr.mode == modeDuotone {
		planes = planes[:1]
	}
	return f.toNRGBA(planes, a, f.hdr.w, f.hdr.h), nil
}

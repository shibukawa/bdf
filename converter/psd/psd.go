// Package psd converts Photoshop documents (.psd, and .psb, the large
// document format) into BDF documents: the image as Photoshop composited
// it, one page per artboard.
//
// A document without artboards becomes one page, the size of its canvas at
// its resolution (a 72 ppi image is one point per pixel). Pages finer than
// the resolution cap of image inputs (imgconv.Options.MaxDPI and
// MaxPixels) are scaled down to it. A document with
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
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/xmp"
	"github.com/shibukawa/bdf/image/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Pages selects 1-based pages (artboards; see converter.Pages); nil
	// converts all of them.
	Pages converter.Pages
	// Title overrides the document title.
	Title string
	// NoArtboards puts the whole canvas on one page even when the document
	// has artboards.
	NoArtboards bool
	// Images controls how the page images are stored and the resolution
	// cap they are scaled down to (see imgconv). The zero value stores them
	// as PNG, capped at imgconv.DefaultMaxDPI and imgconv.DefaultMaxPixels.
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
	// Scaled counts the pages scaled down to the resolution cap.
	Scaled int
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

// maxPageCanvases bounds the pixels of the pages of a document together,
// in canvases: the artboards of a document lie beside each other and stay
// below one. A variable, for the tests.
var maxPageCanvases int64 = 4

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
	composited := err != nil || img == nil // by the converter, once the pages are known
	if composited {
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
	pages := opts.Pages.Numbers(len(regions))
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
	// The pages converted and how each is stored. Each is cropped, scaled
	// and encoded: artboards that lie on each other would have the canvas
	// converted over and over.
	room := maxPageCanvases * int64(canvas.Dx()) * int64(canvas.Dy())
	counted := map[int]bool{}
	var fits []fit
	for i, n := range pages {
		if n < 1 || n > len(regions) {
			return nil, fmt.Errorf("psd: page %d out of range (1-%d)", n, len(regions))
		}
		rg := regions[n-1]
		if !counted[n] {
			counted[n] = true
			if room -= int64(rg.r.Dx()) * int64(rg.r.Dy()); room < 0 {
				warn(fmt.Sprintf("the artboards cover more than %d times the canvas: artboard %d and those after it are left out", maxPageCanvases, n))
				pages = pages[:i]
				break
			}
		}
		fits = append(fits, rg.fit(opts.Images, xres, yres))
	}
	// The layers are composited at the resolution the pages are stored at
	// (scale.go); a composite in the file has that of the document.
	on := scale{full: canvas.Size(), to: canvas.Size()}
	if composited {
		on.to = compositeSize(canvas.Size(), pages, regions, fits)
		if img, err = f.compositeLayers(on.to, warn); err != nil {
			return nil, fmt.Errorf("psd: %w", err)
		}
		res.Composited = true
	}
	for i, n := range pages {
		rg, ft := regions[n-1], fits[i]
		at := on.place(rg.r)
		crop := image.NewNRGBA(image.Rect(0, 0, at.Dx(), at.Dy()))
		for y := 0; y < at.Dy(); y++ {
			copy(crop.Pix[y*crop.Stride:(y+1)*crop.Stride], img.Pix[img.PixOffset(at.Min.X, at.Min.Y+y):])
		}
		w, h := ft.w, ft.h
		var pix image.Image = crop
		if ft.tw != rg.r.Dx() || ft.th != rg.r.Dy() {
			pix = imgconv.Resize(crop, ft.tw, ft.th)
			res.Scaled++
		}
		enc, err := imgconv.EncodePixels(pix, true, opts.Images)
		if err != nil && enc.Data == nil {
			return nil, fmt.Errorf("psd: page %d: %w", n, err)
		}
		obj := bdf.NewObject()
		obj.SetBBox(0, 0, w, h)
		obj.Image(obj.AddImage(doc.AddImage(enc.Data)), 0, 0, w, h)
		hash, _ := doc.AddObject(obj)
		view.AddPage(w, h, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
		if rg.artboard {
			res.Artboards++
		}
	}
	res.Pages = len(view.Pages)
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

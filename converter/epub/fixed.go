package epub

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/webdoc"
	"github.com/shibukawa/bdf/imgconv"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"golang.org/x/net/html"
)

// pxToPt converts CSS pixels to points.
const pxToPt = 0.75

// fixedLayout reports whether the book converts into pages of pictures:
// every chapter of its reading order is a fixed-layout page that holds one
// picture and nothing else (a comic, a photo book).
func (c *converter) fixedLayout() bool {
	n := 0
	for _, ch := range c.chapters {
		if !ch.linear {
			continue
		}
		if _, ok := ch.singlePicture(); !ok || !ch.fixed {
			return false
		}
		n++
	}
	return n > 0
}

// storedImage is a picture stored in the document.
type storedImage struct {
	hash bdf.Hash
	w, h float64 // intrinsic size in pixels (0: unknown)
}

// pictures converts a fixed-layout book of pictures: each page of its
// reading order becomes a page of a fixed view, at the size of its viewport
// (else of its picture), with the picture fitted into it and centered (as
// SVG's xMidYMid meet does). Chapters that are not in the reading order
// are left out.
func (c *converter) pictures() (*Result, error) {
	doc := bdf.NewDocument()
	doc.Meta.Source = "epub"
	doc.Meta.DC = c.dc
	res := &Result{Doc: doc, FixedLayout: true}
	var pages []*chapter
	for _, ch := range c.chapters {
		if ch.linear {
			pages = append(pages, ch)
		} else {
			c.warnOnce("nonlinear", "the fixed-layout pages outside the reading order are left out")
		}
	}
	sel := c.opts.Pages.Numbers(len(pages))
	for _, n := range sel {
		if n < 1 || n > len(pages) {
			return nil, fmt.Errorf("epub: page %d out of range (1-%d)", n, len(pages))
		}
	}
	if c.opts.Views == ViewsScroll {
		c.warn("a fixed-layout book has no scroll view; its pages are converted")
	}
	view := doc.NewView("pages", bdf.ViewFixed, c.dc.Title.First())
	stored := map[string]*storedImage{}
	var last [2]float64
	for i, ch := range pages {
		if sel != nil && !slices.Contains(sel, i+1) {
			continue
		}
		pic, _ := ch.singlePicture()
		var im *storedImage
		var err error
		if pic.svg != nil {
			im = c.storeSVG(doc, pic.svg)
		} else {
			im, err = c.store(doc, pic.src, stored)
		}
		if err != nil {
			c.warn(fmt.Sprintf("%s: %v; the page is left blank", ch.path, err))
		}
		w, h := ch.viewport[0], ch.viewport[1]
		if w <= 0 || h <= 0 {
			w, h = pic.w, pic.h
		}
		if (w <= 0 || h <= 0) && im != nil {
			w, h = im.w, im.h
		}
		if w <= 0 || h <= 0 {
			w, h = last[0], last[1]
		}
		if w <= 0 || h <= 0 {
			w, h = 600, 800
		}
		last = [2]float64{w, h}
		pw, ph := float32(w*pxToPt), float32(h*pxToPt)
		if im == nil {
			view.AddPage(pw, ph)
			continue
		}
		// the picture fitted into the page, centered
		x, y, dw, dh := float32(0), float32(0), pw, ph
		if im.w > 0 && im.h > 0 {
			s := min(float64(pw)/im.w, float64(ph)/im.h)
			dw, dh = float32(im.w*s), float32(im.h*s)
			x, y = (pw-dw)/2, (ph-dh)/2
		}
		o := bdf.NewObject()
		o.SetBBox(0, 0, pw, ph)
		ref := o.AddImage(im.hash)
		alt := strings.Join(strings.Fields(pic.alt), " ")
		if alt != "" {
			o.Mark(bdf.MarkFigure, alt)
		}
		o.Image(ref, x, y, dw, dh)
		if alt != "" {
			o.Mark(bdf.MarkEnd, "")
		}
		hash, _ := doc.AddObject(o)
		view.AddPage(pw, ph, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
		res.Images++
	}
	res.Pages = len(view.Pages)
	res.Chapters = len(pages)
	res.Warnings = c.warnings
	return res, nil
}

// storeSVG stores an svg element of a page as an SVG document (see
// webdoc.SVGDocument), with the elements it uses from elsewhere in its
// document and its images as data: URLs.
func (c *converter) storeSVG(doc *bdf.Document, n *html.Node) *storedImage {
	top := n
	for top.Parent != nil {
		top = top.Parent
	}
	ids := map[string]*html.Node{}
	webdoc.WalkElements(top, func(e *html.Node) {
		if id := attrVal(e, "id"); id != "" && ids[id] == nil {
			ids[id] = e
		}
	})
	data := webdoc.SVGDocument(n, "", func(id string) *html.Node { return ids[id] }, c.embed)
	im := &storedImage{hash: doc.AddImage(data)}
	im.w, im.h = imgconv.ParseSVGSize(attrVal(n, "width"), attrVal(n, "height"), attrVal(n, "viewBox")).Pixels()
	return im
}

// embed returns a picture of the publication as a data: URL, for an SVG
// document drawn as an image, which loads nothing ("" when it cannot be
// read).
func (c *converter) embed(href string) string {
	if href == "" || strings.HasPrefix(href, "#") || webdoc.IsDataURL(href) {
		return ""
	}
	data, err := c.image(href)
	if err != nil {
		c.warn(fmt.Sprintf("image %s in an SVG: %v", href, err))
		return ""
	}
	format := imgconv.Sniff(data)
	if format == "" {
		return ""
	}
	return "data:" + imgconv.MIME(format) + ";base64," + base64.StdEncoding.EncodeToString(data)
}

// store stores a picture once.
func (c *converter) store(doc *bdf.Document, src string, stored map[string]*storedImage) (*storedImage, error) {
	if im, ok := stored[src]; ok {
		return im, nil
	}
	data, err := c.image(src)
	if err != nil {
		return nil, err
	}
	switch imgconv.Sniff(data) {
	case "svg":
		// stored as it is: the viewer draws it (spec §6.2)
		width, height, viewBox, _ := imgconv.SVGRoot(data)
		im := &storedImage{hash: doc.AddImage(data)}
		im.w, im.h = imgconv.ParseSVGSize(width, height, viewBox).Pixels()
		stored[src] = im
		return im, nil
	case "":
		return nil, fmt.Errorf("%s: unsupported picture format", src)
	}
	im := &storedImage{}
	if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		im.w, im.h = float64(cfg.Width), float64(cfg.Height)
	}
	if r, err := imgconv.Optimize(data, c.opts.Images); err == nil || r.Data != nil {
		data = r.Data
	}
	im.hash = doc.AddImage(data)
	stored[src] = im
	return im, nil
}

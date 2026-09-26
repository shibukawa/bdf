// Package image converts the images browsers display by themselves (PNG,
// JPEG, GIF, WebP, AVIF, BMP, ICO and SVG) into BDF documents: one page the
// size of the image, which draws it. The image passes through: it is stored
// as it is, neither decoded nor re-encoded, and the browser decodes it as it
// would the file, so the page shows what the browser shows for the file.
//
// What the converter reads is the image's size and its metadata, which
// become the document's Dublin Core as other formats' document properties
// do (docs/spec.md §4.3): XMP, EXIF (with the Windows tags) and IPTC, the
// text chunks of PNG, the comments of GIF, and the title, description and
// RDF metadata of SVG. See docs/design.md §3.13.
package image

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shibukawa/bdf"
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title (by default the image's title).
	Title string
	// FileName is the image's file name, the alternative text of an image
	// that has no description or title.
	FileName string
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Format is the image format: "png", "jpeg", "gif", "webp", "avif",
	// "bmp", "ico" or "svg".
	Format string
	// Width and Height are the image's size in CSS pixels, as it is drawn
	// (after its orientation). The page is as large, at 96 pixels per inch.
	Width, Height float64
	// Animated reports an image with several frames, of which the page
	// shows the first.
	Animated bool
}

// maxSize bounds the images read.
const maxSize = 1 << 30

// ptPerPx converts CSS pixels to points.
const ptPerPx = 0.75

// formatNames name the formats for people.
var formatNames = map[string]string{
	"png": "PNG", "jpeg": "JPEG", "gif": "GIF", "webp": "WebP", "avif": "AVIF", "bmp": "BMP", "ico": "ICO", "svg": "SVG",
}

// ConvertFile converts an image file.
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
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.FileName == "" {
		o.FileName = filepath.Base(path)
	}
	return Convert(f, st.Size(), &o)
}

// Convert converts an image read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("image: the file is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("image: %w", err)
	}
	format := sniff(data[:min(len(data), 1024)], bytes.NewReader(data), size)
	if format == "" {
		return nil, fmt.Errorf("image: not a PNG, JPEG, GIF, WebP, AVIF, BMP, ICO or SVG image")
	}
	var p *picture
	var err error
	if format == "svg" {
		p, err = readSVG(data)
	} else {
		p, err = readRaster(format, data)
	}
	if err != nil {
		return nil, fmt.Errorf("image: %s: %w", formatNames[format], err)
	}
	res := &Result{Doc: bdf.NewDocument(), Format: format, Animated: p.animated}
	warn := func(msg string) {
		if opts.Warn != nil {
			opts.Warn(msg)
		} else {
			res.Warnings = append(res.Warnings, msg)
		}
	}
	for _, w := range p.warnings {
		warn(w)
	}

	// metadata, in order of precedence: XMP, EXIF, IPTC, the format's own
	var dc bdf.DublinCore
	for _, x := range p.xmp {
		merge(&dc, parseXMP(x))
	}
	orientation := 0
	for _, e := range p.exif {
		if info, ok := parseEXIF(e); ok {
			merge(&dc, info.dc)
			if orientation == 0 {
				orientation = info.orientation
			}
		}
	}
	if p.iptc != nil {
		merge(&dc, parseIPTC(p.iptc))
	}
	merge(&dc, p.native)

	// Browsers turn JPEG and PNG images by their EXIF orientation (5 to 8
	// swap the sides), but not WebP images; AVIF has its own (irot).
	w, h := p.w, p.h
	switch {
	case orientation >= 5 && (format == "jpeg" || format == "png"):
		w, h = h, w
	case orientation > 1 && format == "webp":
		warn("the EXIF orientation of the WebP image is not applied: browsers draw WebP images as they are stored")
	}
	res.Width, res.Height = w, h
	if p.animated {
		warn("animated image: the first frame is drawn")
	}

	doc := res.Doc
	doc.Meta.Source = format
	doc.Meta.DC = dc
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	pw, ph := float32(w*ptPerPx), float32(h*ptPerPx)
	o := bdf.NewObject()
	img := o.AddImage(doc.AddImage(data))
	// the image as a figure, for screen readers
	alt := dc.Description.First()
	if alt == "" {
		alt = doc.Meta.DC.Title.First()
	}
	if alt == "" {
		alt = opts.FileName
	}
	o.Mark(bdf.MarkFigure, alt)
	if format != "svg" {
		o.Smoothing(true, 2) // high: an image scaled down looks as it does in an image element
	}
	o.Image(img, 0, 0, pw, ph)
	o.Mark(bdf.MarkEnd, "")
	hash, _ := doc.AddObject(o)
	view := doc.NewView("pages", bdf.ViewFixed, doc.Meta.DC.Title.First())
	view.AddPage(pw, ph, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	return res, nil
}

// Summary describes a result in a line ("1 page (640 × 480 px, JPEG)").
func (r *Result) Summary() string {
	s := fmt.Sprintf("1 page (%s × %s px, %s", px(r.Width), px(r.Height), formatNames[r.Format])
	if r.Animated {
		s += ", animated"
	}
	return s + ")"
}

// px writes a size in pixels with at most two decimals.
func px(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

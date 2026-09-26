package pptx

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"regexp"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
	_ "golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// imageEntry is an image part stored in the document.
type imageEntry struct {
	target string
	hash   bdf.Hash
	w, h   float64 // intrinsic size in pixels (0 when unknown)
	ok     bool
	mf     []byte // a Windows metafile, replayed instead of stored
	mfKind string
}

// image loads the image a blip relationship points to.
func (c *converter) image(part, rid string) *imageEntry {
	r, ok := c.pkg.target(part, rid)
	if !ok || r.External {
		if ok {
			c.warnOnce("extimg", "linked (external) images are not embedded")
		}
		return &imageEntry{}
	}
	if e, ok := c.images[r.Target]; ok {
		return e
	}
	e := &imageEntry{target: r.Target}
	c.images[r.Target] = e
	data, err := c.pkg.read(r.Target)
	if err != nil {
		c.warnf("image %s: %v", r.Target, err)
		return e
	}
	format := imgconv.Sniff(data)
	switch {
	case format != "":
	case bytes.HasPrefix(data, []byte("II*\x00")) || bytes.HasPrefix(data, []byte("MM\x00*")):
		img, err := tiff.Decode(bytes.NewReader(data))
		if err != nil {
			c.warnf("image %s: %v", r.Target, err)
			return e
		}
		res, _ := imgconv.EncodeImage(toNRGBA(img), true, c.opts.Images)
		data, format = res.Data, res.Format
	case metafileKind(data) != "":
		e.mf, e.mfKind, e.ok = data, metafileKind(data), true
		return e
	default:
		c.warnOnce("imgfmt:"+r.Target, "unsupported image format (%s)", r.Target)
		return e
	}
	if format == "svg" {
		e.w, e.h = svgSize(data)
	} else if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		e.w, e.h = float64(cfg.Width), float64(cfg.Height)
	}
	if format != "svg" {
		if res, err := imgconv.Optimize(data, c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
	}
	e.hash = c.doc.AddImage(data)
	e.ok = true
	return e
}

func toNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Src)
	return out
}

var (
	svgTag     = regexp.MustCompile(`(?s)<svg\b[^>]*>`)
	svgAttr    = regexp.MustCompile(`\b(width|height|viewBox)\s*=\s*["']([^"']*)["']`)
	svgNumUnit = regexp.MustCompile(`^\s*([0-9.]+)\s*(px|pt)?\s*$`)
)

// svgSize reads the intrinsic size of an SVG document.
func svgSize(data []byte) (w, h float64) {
	tag := svgTag.Find(data)
	var vb []float64
	for _, m := range svgAttr.FindAllSubmatch(tag, -1) {
		v := string(m[2])
		switch string(m[1]) {
		case "width", "height":
			if mm := svgNumUnit.FindStringSubmatch(v); mm != nil {
				f, _ := strconv.ParseFloat(mm[1], 64)
				if mm[2] == "pt" {
					f *= 96.0 / 72
				}
				if string(m[1]) == "width" {
					w = f
				} else {
					h = f
				}
			}
		case "viewBox":
			for _, s := range strings.FieldsFunc(v, func(r rune) bool { return r == ' ' || r == ',' }) {
				f, _ := strconv.ParseFloat(s, 64)
				vb = append(vb, f)
			}
		}
	}
	if (w == 0 || h == 0) && len(vb) == 4 {
		w, h = vb[2], vb[3]
	}
	return w, h
}

package drawingml

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/metafile"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/image/imgconv"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// imageEntry is an image part stored in the document.
type imageEntry struct {
	target string
	hash   bdf.Hash
	w, h   float64 // intrinsic size in pixels (0 when unknown)
	ok     bool
	mf     []byte // a Windows metafile, replayed instead of stored
	// paint draws what stands for the image in its box (InjectDrawing)
	paint func(o *bdf.Object, x, y, w, h float32)
}

// injectedPrefix starts the relationship ids of the images given to
// InjectImage: no relationship of a package has such an id.
const injectedPrefix = "\x00image:"

// InjectImage stores a picture that is not a part of the package (the
// picture of a filled template) and returns the relationship id that a
// blip refers to it by (r:embed), with its size in pixels. ok is false
// when the data is not a picture that is stored.
func (c *Renderer) InjectImage(data []byte) (rid string, w, h float64, ok bool) {
	format := imgconv.Sniff(data)
	if format == "svg" {
		// (as it is: a signature drawn on a pad is one)
		s, ok := imgconv.SVGDocSize(data)
		if !ok {
			return "", 0, 0, false
		}
		e := &imageEntry{hash: c.doc.AddImage(data), ok: true}
		e.w, e.h = s.Pixels()
		rid = injectedPrefix + e.hash.String()
		c.images[rid] = e
		return rid, e.w, e.h, true
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if format == "" || format == "svg" || err != nil || cfg.Width <= 0 || cfg.Height <= 0 {
		return "", 0, 0, false
	}
	if res, err := imgconv.Optimize(data, c.imgOpts); err == nil || res.Data != nil {
		data = res.Data
	}
	e := &imageEntry{hash: c.doc.AddImage(data), w: float64(cfg.Width), h: float64(cfg.Height), ok: true}
	rid = injectedPrefix + e.hash.String()
	c.images[rid] = e
	return rid, e.w, e.h, true
}

// InjectDrawing registers a drawing that stands for a picture (the QR code
// or the barcode of a filled template, which is drawn as vectors) and
// returns the relationship id that a blip refers to it by: paint draws it
// in the box of the picture.
func (c *Renderer) InjectDrawing(paint func(o *bdf.Object, x, y, w, h float32)) (rid string) {
	rid = injectedPrefix + "drawing:" + strconv.Itoa(len(c.images))
	c.images[rid] = &imageEntry{ok: true, paint: paint}
	return rid
}

// image loads the image a blip relationship points to.
func (c *Renderer) image(part, rid string) *imageEntry {
	if strings.HasPrefix(rid, injectedPrefix) {
		if e, ok := c.images[rid]; ok {
			return e
		}
		return &imageEntry{}
	}
	r, ok := c.pkg.Target(part, rid)
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
	data, err := c.pkg.Read(r.Target)
	if err != nil {
		c.warnf("image %s: %v", r.Target, err)
		return e
	}
	format := imgconv.Sniff(data)
	switch {
	case format != "":
	case tiff.Sniff(data):
		// The first page, as JPEG or PNG; Optimize below converts it.
		pic, damaged, err := tiff.Picture(data, c.imgOpts)
		if err != nil {
			c.warnf("image %s: %v", r.Target, err)
			return e
		}
		if damaged {
			c.warnf("image %s: the pixel data is damaged; what is missing is left blank", r.Target)
		}
		data, format = pic, imgconv.Sniff(pic)
	case metafile.Kind(data) != "":
		e.mf, e.ok = data, true
		return e
	default:
		c.warnOnce("imgfmt:"+r.Target, "unsupported image format (%s)", r.Target)
		return e
	}
	if format == "svg" {
		if s, ok := imgconv.SVGDocSize(data); ok {
			e.w, e.h = s.Natural()
		}
	} else if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		e.w, e.h = float64(cfg.Width), float64(cfg.Height)
	}
	if format != "svg" {
		if res, err := imgconv.Optimize(data, c.imgOpts); err == nil || res.Data != nil {
			data = res.Data
		}
	}
	e.hash = c.doc.AddImage(data)
	e.ok = true
	return e
}

var ()

package drawingml

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
		e.w, e.h = svgSize(data)
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

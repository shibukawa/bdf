package drawingml

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// recolor is the chain of per-pixel effects of a blip (a:clrChange,
// a:clrRepl, a:duotone, a:biLevel, a:grayscl, a:lum), applied in document
// order. Raster
// pictures are re-encoded with the effects baked in; metafiles apply them
// to the colors of their records.
type recolor struct {
	key string
	fx  []recolorFx
}

type recolorFx struct {
	kind     string
	from, to rgba // clrChange; clrRepl and duotone use to (and from)
	useA     bool
	thresh   float64 // biLevel
	bright   float64 // lum, in 0-255 units
	contrast float64 // lum slope
}

// blipRecolor returns the pixel effects of a:blip, or nil when it has none.
func blipRecolor(blip *ooxml.Node, cc *colorCtx) *recolor {
	rc := &recolor{}
	var key strings.Builder
	for _, k := range blip.Elements() {
		fx := recolorFx{kind: k.Name}
		switch k.Name {
		case "clrChange":
			from, ok1 := cc.color(k.Child("clrFrom"))
			to, ok2 := cc.color(k.Child("clrTo"))
			if !ok1 || !ok2 || from == to {
				continue
			}
			fx.from, fx.to, fx.useA = from, to, k.AttrBool("useA", true)
		case "clrRepl":
			c, ok := cc.color(k)
			if !ok {
				continue
			}
			fx.to = c
		case "duotone":
			var cs []rgba
			for _, c := range k.Elements() {
				if isColor(c) {
					if v, ok := cc.resolve(c); ok {
						cs = append(cs, v)
					}
				}
			}
			if len(cs) != 2 {
				continue
			}
			fx.from, fx.to = cs[0], cs[1]
		case "biLevel":
			fx.thresh = k.AttrPct("thresh", 0.5)
		case "grayscl":
		case "lum":
			b, c := k.AttrPct("bright", 0), k.AttrPct("contrast", 0)
			if b == 0 && c == 0 {
				continue
			}
			fx.bright = math.Max(-1, math.Min(b, 1)) * 255
			c = math.Max(-1, math.Min(c, 1)) * 100
			if c >= 0 {
				fx.contrast = 128 / math.Max(128-1.27*c, 1)
			} else {
				fx.contrast = (128 + 1.27*c) / 128
			}
		default:
			continue
		}
		rc.fx = append(rc.fx, fx)
		key.WriteString(fx.kind)
		for _, c := range []rgba{fx.from, fx.to} {
			key.WriteString(" " + ftoa(c.R) + "," + ftoa(c.G) + "," + ftoa(c.B) + "," + ftoa(c.A))
		}
		key.WriteString(" " + ftoa(fx.thresh) + " " + ftoa(fx.bright) + " " + ftoa(fx.contrast))
		if fx.useA {
			key.WriteString(" a")
		}
		key.WriteString(";")
	}
	if len(rc.fx) == 0 {
		return nil
	}
	rc.key = key.String()
	return rc
}

func to8(v float64) int { return int(math.Round(clamp01(v) * 255)) }

func luma(r, g, b int) float64 { return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255 }

// pixel applies the effects to one non-premultiplied pixel. tol is the
// per-channel tolerance of a:clrChange (lossy sources need some).
func (rc *recolor) pixel(c color.NRGBA, tol int) color.NRGBA {
	for _, fx := range rc.fx {
		switch fx.kind {
		case "clrChange":
			if abs(int(c.R)-to8(fx.from.R)) <= tol && abs(int(c.G)-to8(fx.from.G)) <= tol &&
				abs(int(c.B)-to8(fx.from.B)) <= tol && (!fx.useA || abs(int(c.A)-to8(fx.from.A)) <= tol) {
				c = color.NRGBA{uint8(to8(fx.to.R)), uint8(to8(fx.to.G)), uint8(to8(fx.to.B)), uint8(to8(fx.to.A))}
			}
		case "clrRepl":
			c = color.NRGBA{uint8(to8(fx.to.R)), uint8(to8(fx.to.G)), uint8(to8(fx.to.B)), c.A}
		case "duotone":
			l := luma(int(c.R), int(c.G), int(c.B))
			mix := func(a, b float64) uint8 { return uint8(to8(a*(1-l) + b*l)) }
			c = color.NRGBA{mix(fx.from.R, fx.to.R), mix(fx.from.G, fx.to.G), mix(fx.from.B, fx.to.B), c.A}
		case "biLevel":
			v := uint8(0)
			if luma(int(c.R), int(c.G), int(c.B)) >= fx.thresh {
				v = 255
			}
			c = color.NRGBA{v, v, v, c.A}
		case "grayscl":
			v := uint8(to8(luma(int(c.R), int(c.G), int(c.B))))
			c = color.NRGBA{v, v, v, c.A}
		case "lum":
			// PowerPoint applies half of the brightness before the
			// contrast and half after it
			f := func(v uint8) uint8 {
				return uint8(math.Round(math.Max(0, math.Min(255, (float64(v)+fx.bright/2-128)*fx.contrast+128+fx.bright/2))))
			}
			c = color.NRGBA{f(c.R), f(c.G), f(c.B), c.A}
		}
	}
	return c
}

// Color applies the effects to a color of a metafile record
// (metafile.Recolor).
func (rc *recolor) Color(c color.NRGBA) color.NRGBA { return rc.pixel(c, 9) }

// Image applies the effects to a metafile bitmap, whose pixels are exact
// (metafile.Recolor).
func (rc *recolor) Image(img *image.NRGBA) { rc.pixels(img, 0) }

// pixels applies the effects to decoded pixels in place.
func (rc *recolor) pixels(img *image.NRGBA, tol int) {
	b := img.Bounds()
	for y := 0; y < b.Dy(); y++ {
		row := img.Pix[y*img.Stride : y*img.Stride+b.Dx()*4]
		for i := 0; i < len(row); i += 4 {
			c := rc.pixel(color.NRGBA{row[i], row[i+1], row[i+2], row[i+3]}, tol)
			row[i], row[i+1], row[i+2], row[i+3] = c.R, c.G, c.B, c.A
		}
	}
}

func opaque(img *image.NRGBA) bool {
	for i := 3; i < len(img.Pix); i += 4 {
		if img.Pix[i] != 255 {
			return false
		}
	}
	return true
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// recolored returns the raster picture of e with the effects baked in.
func (c *Renderer) recolored(e *imageEntry, rc *recolor) *imageEntry {
	key := e.target + "\x00" + rc.key
	if r, ok := c.images[key]; ok {
		return r
	}
	r := &imageEntry{target: e.target, w: e.w, h: e.h}
	c.images[key] = r
	data, err := c.pkg.Read(e.target)
	if err != nil {
		return e
	}
	var img image.Image
	format := imgconv.Sniff(data)
	// A photograph (a JPEG, or a JPEG page of a TIFF file) is stored as JPEG again.
	photo := format == "jpeg"
	tol := 9 // same guesses as LibreOffice (tdf#149670)
	switch format {
	case "jpeg":
		tol = 15
	case "png":
		tol = 1
	case "bmp":
		tol = 0
	}
	switch {
	case format == "svg" || format == "avif" || format == "":
		if tiff.Sniff(data) {
			var d *tiff.IFD
			if d, err = tiff.FirstPage(data); err == nil {
				img, _, err = d.Decode()
				photo = d.Compression() == tiff.CompressionJPEG
			}
			tol = 1
			break
		}
		c.warnOnce("recolor:"+e.target, "picture recoloring is not supported for %s", e.target)
		return e
	default:
		img, err = imgconv.Decode(data)
	}
	if err != nil {
		c.warnf("image %s: %v", e.target, err)
		return e
	}
	n := toNRGBA(img)
	if n == img {
		n = image.NewNRGBA(n.Rect)
		copy(n.Pix, img.(*image.NRGBA).Pix)
	}
	rc.pixels(n, tol)
	var out []byte
	if photo && opaque(n) {
		// a photograph stays a JPEG: much smaller and faster than PNG
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, n, &jpeg.Options{Quality: 90}); err == nil {
			out = buf.Bytes()
			if res, err := imgconv.Optimize(out, c.imgOpts); err == nil || res.Data != nil {
				out = res.Data
			}
		}
	}
	if out == nil {
		res, _ := imgconv.EncodeImage(n, !photo, c.imgOpts)
		out = res.Data
	}
	r.hash = c.doc.AddImage(out)
	r.ok = true
	return r
}

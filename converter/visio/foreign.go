package visio

import (
	"image"
	"image/draw"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/metafile"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// Pictures and embedded objects are foreign shapes: their ForeignData is a
// picture (in an image part), or an OLE
// object whose preview picture is drawn. The picture fills the rectangle
// the ImgOffsetX, ImgOffsetY, ImgWidth and ImgHeight cells give, cropped
// by the shape's box. Windows metafiles are replayed.

// picture is a foreign data item ready to draw.
type picture struct {
	hash bdf.Hash // a stored image
	mf   []byte   // or a metafile
	ok   bool
}

// picture loads the picture of a ForeignData element of a shape in part.
func (c *converter) picture(fd *ooxml.Node, part string) *picture {
	key := part + "\x00" + fd.Child("Rel").RelID("id")
	if pc, ok := c.pictures[key]; ok {
		return pc
	}
	pc := &picture{}
	c.pictures[key] = pc
	r, ok := c.d.pkg.Target(part, fd.Child("Rel").RelID("id"))
	if !ok {
		return pc
	}
	data, err := c.d.pkg.Read(r.Target)
	if err != nil {
		c.warnf("picture %s: %v", r.Target, err)
		return pc
	}
	if imgconv.Sniff(data) == "" && metafile.Kind(data) == "" && !tiff.Sniff(data) {
		// an embedded object: its preview picture
		data = nil
		for _, pr := range c.d.pkg.Rels(r.Target) {
			if strings.HasSuffix(pr.Type, "/image") && !pr.External {
				data, _ = c.d.pkg.Read(pr.Target)
				break
			}
		}
		if data == nil {
			c.warnOnce("ole", "embedded objects without a preview picture are not drawn")
			return pc
		}
	}
	switch {
	case metafile.Kind(data) != "":
		pc.mf, pc.ok = data, true
		return pc
	case tiff.Sniff(data):
		// The first page, as JPEG or PNG, then converted like the others.
		pic, damaged, err := tiff.Picture(data, c.opts.Images)
		if err != nil {
			c.warnf("picture: %v", err)
			return pc
		}
		if damaged {
			c.warnOnce("tiffdamaged", "a TIFF picture's pixel data is damaged; what is missing is left blank")
		}
		data = pic
		if res, err := imgconv.Optimize(data, c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
	case imgconv.Sniff(data) == "bmp":
		img, err := imgconv.Decode(data)
		if err != nil {
			c.warnf("picture: %v", err)
			return pc
		}
		res, _ := imgconv.EncodeImage(toNRGBA(img), true, c.opts.Images)
		data = res.Data
	case imgconv.Sniff(data) == "":
		c.warnOnce("picfmt", "a picture in an unsupported format is not drawn")
		return pc
	default:
		if res, err := imgconv.Optimize(data, c.opts.Images); err == nil || res.Data != nil {
			data = res.Data
		}
	}
	pc.hash, pc.ok = c.doc.AddImage(data), true
	return pc
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

// drawForeign draws the picture of a foreign shape.
func (p *pageCtx) drawForeign(cv *canvas.Canvas, s *shape, m canvas.Matrix) {
	fd, owner := s.foreignNode()
	if t := fd.AttrStr("ForeignType", ""); t == "Ink" {
		p.c.warnOnce("ink", "ink is not drawn")
		return
	}
	pc := p.c.picture(fd, owner.part)
	if !pc.ok {
		return
	}
	w, h := p.num(s, "Width", 0), p.num(s, "Height", 0)
	ix, iy := p.dim(s, "ImgOffsetX", 0), p.dim(s, "ImgOffsetY", 0)
	iw, ih := p.dim(s, "ImgWidth", w), p.dim(s, "ImgHeight", h)
	k := 72 * p.scale
	// the shape's box in points, y down
	box := m.Mul(canvas.Translate(0, h)).Mul(canvas.Scale(1/k, -1/k))
	x, y := ix*k, (h-iy-ih)*k
	cv.Obj.Save()
	cv.Transform(box)
	cv.Obj.ClipRect(0, 0, f32(w*k), f32(h*k))
	if pc.mf != nil {
		metafile.Draw(cv, pc.mf, x, y, iw*k, ih*k, &metafile.Options{Doc: p.c.doc, Fonts: p.c.fonts,
			Images: p.c.opts.Images, Warn: func(msg string) { p.c.warnOnce("mf:"+msg, "%s", msg) }})
	} else {
		cv.Obj.Image(cv.Image(pc.hash), f32(x), f32(y), f32(iw*k), f32(ih*k))
	}
	cv.Obj.Restore()
	cv.Drawn = true
}

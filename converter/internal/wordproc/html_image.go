package wordproc

import (
	"bytes"
	"cmp"
	"encoding/binary"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/imgconv"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
	"golang.org/x/net/html"
)

// htmlImage is an image the document refers to, stored in the BDF.
type htmlImage struct {
	hash  bdf.Hash
	w, h  float64 // natural size in points (0: none)
	ratio float64 // natural aspect ratio, width / height (0: none)
}

// img reads an img: a picture in the line, shrunk to the line when it is
// wider, or floating left or right with the text beside it.
func (r *htmlReader) img(n *html.Node, css map[string]string, st *hstyle) {
	alt, hasAlt := attr(n, "alt")
	alt = strings.Join(strings.Fields(alt), " ")
	im := r.image(imageSource(n))
	if im == nil {
		// what browsers show for a picture they cannot draw
		if alt != "" {
			s := st.chars()
			s.rp.color = mutedColor
			r.text(alt, s)
		}
		return
	}
	if !hasAlt && st.figAlt != "" {
		alt = st.figAlt
	}
	aw, ah := r.length(attrStr(n, "width")), r.length(attrStr(n, "height"))
	if v := r.length(css["width"]); v > 0 {
		aw = v
	}
	if v := r.length(css["height"]); v > 0 {
		ah = v
	}
	r.picture(n, css, st, im, aw, ah, alt)
}

// picture places an image in the text: in the line, shrunk to the line
// when it is wider, or floating left or right with the text beside it. aw
// and ah are the width and height its attributes and style ask for (0:
// not said).
func (r *htmlReader) picture(n *html.Node, css map[string]string, st *hstyle, im *htmlImage, aw, ah float64, alt string) {
	w, h, fill := replacedSize(im.w, im.h, im.ratio, aw, ah)
	hash := im.hash
	o := &inlineObj{w: w, h: h, alt: alt, fit: true, fill: fill, paint: func(e *emitter, box drawingml.Box) {
		e.cv.Obj.Image(e.cv.Image(hash), f32(box.X), f32(box.Y), f32(box.W), f32(box.H))
	}}
	float := css["float"]
	if float == "" {
		float = strings.ToLower(attrStr(n, "align"))
	}
	if float == "left" || float == "right" {
		p := r.para(st)
		f := &floatObj{inlineObj: *o, hRel: "column", vRel: "paragraph", hAlign: float, wrap: "square",
			dist: [4]float64{0, r.em(0.5), 0, 0}}
		if float == "right" {
			f.side, f.dist[2] = "left", r.em(1)
		} else {
			f.side, f.dist[3] = "right", r.em(1)
		}
		p.anchors = append(p.anchors, f)
		return
	}
	p := r.para(st)
	if r.space {
		if n := len(p.items); n > 0 && p.items[n-1].kind != kBreak {
			r.c.addChar(p, r.style(r.spaceSt), ' ')
		}
		r.space, r.spaceNL = false, false
	}
	p.items = append(p.items, item{kind: kObject, obj: o, st: r.style(st), w: o.w})
}

// replacedSize sizes a picture as CSS sizes replaced elements: the width
// and height asked for (aw, ah; 0 when not said), a missing one following
// from the natural aspect ratio, else the natural size (w, h), else 300 ×
// 150 CSS pixels. fill reports a picture with a ratio but no size at all
// (an SVG image with only a view box), which browsers make as wide as the
// line; its size is then a nominal 300 pixels wide.
func replacedSize(w, h, ratio, aw, ah float64) (float64, float64, bool) {
	dw, dh := 300*pxToPt, 150*pxToPt
	switch {
	case aw > 0 && ah > 0:
		return aw, ah, false
	case aw > 0 && ratio > 0:
		return aw, aw / ratio, false
	case ah > 0 && ratio > 0:
		return ah * ratio, ah, false
	case aw > 0:
		return aw, cmp.Or(h, dh), false
	case ah > 0:
		return cmp.Or(w, dw), ah, false
	case w > 0 && h > 0:
		return w, h, false
	case w > 0:
		return w, dh, false
	case h > 0:
		return dw, h, false
	case ratio > 0:
		return dw, dw / ratio, true
	}
	return dw, dh, false
}

// imageSource returns the address of an img's picture: src, else the
// first candidate of srcset.
func imageSource(n *html.Node) string {
	if s := attrStr(n, "src"); s != "" {
		return s
	}
	if f := strings.Fields(attrStr(n, "srcset")); len(f) > 0 {
		return strings.TrimSuffix(f[0], ",")
	}
	return ""
}

// length reads a CSS length or a size attribute (CSS pixels without a
// unit) in points; 0 for percentages and what it cannot read.
func (r *htmlReader) length(v string) float64 {
	v = strings.TrimSpace(strings.ToLower(v))
	for _, u := range []struct {
		suffix string
		pt     float64
	}{{"px", pxToPt}, {"pt", 1}, {"rem", r.size}, {"em", r.size}, {"in", 72}, {"cm", 72 / 2.54}, {"mm", 72 / 25.4}, {"", pxToPt}} {
		if s, ok := strings.CutSuffix(v, u.suffix); ok {
			f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
			if err != nil || f <= 0 {
				return 0
			}
			return f * u.pt
		}
	}
	return 0
}

// image loads, measures and stores an image once per address; nil when it
// cannot be drawn.
func (r *htmlReader) image(src string) *htmlImage {
	if src == "" || r.d.Image == nil {
		return nil
	}
	if im, ok := r.images[src]; ok {
		return im
	}
	r.images[src] = nil
	name := src
	if len(name) > 80 {
		name = name[:77] + "…"
	}
	data, err := r.d.Image(src)
	if err != nil {
		r.c.warnf("image %s: %v", name, err)
		return nil
	}
	format := imgconv.Sniff(data)
	switch format {
	case "svg":
		// stored as it is: the viewer draws it (spec §6.2)
		width, height, viewBox, _ := imgconv.SVGRoot(data)
		s := imgconv.ParseSVGSize(width, height, viewBox)
		im := &htmlImage{hash: r.c.doc.AddImage(data), w: s.W * pxToPt, h: s.H * pxToPt, ratio: s.Ratio()}
		r.images[src] = im
		return im
	case "":
		r.c.warnf("image %s: unsupported format", name)
		return nil
	}
	im := &htmlImage{}
	if format == "avif" {
		im.w, im.h = avifSize(data)
	} else if cfg, _, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		im.w, im.h = float64(cfg.Width), float64(cfg.Height)
	}
	im.w, im.h = im.w*pxToPt, im.h*pxToPt
	if im.w > 0 && im.h > 0 {
		im.ratio = im.w / im.h
	}
	if res, err := imgconv.Optimize(data, r.c.opts.Images); err == nil || res.Data != nil {
		data = res.Data
	}
	im.hash = r.c.doc.AddImage(data)
	r.images[src] = im
	return im
}

// avifSize reads the size of an AVIF image from its first image spatial
// extents property (ispe).
func avifSize(b []byte) (w, h float64) {
	i := bytes.Index(b, []byte("ispe"))
	if i < 4 || i+16 > len(b) {
		return 0, 0
	}
	return float64(binary.BigEndian.Uint32(b[i+8:])), float64(binary.BigEndian.Uint32(b[i+12:]))
}

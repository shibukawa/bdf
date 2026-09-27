package wordproc

import (
	"bytes"
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
	hash bdf.Hash
	w, h float64 // intrinsic size in points (0: unknown)
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
	w, h := im.w, im.h
	aw, ah := r.length(attrStr(n, "width")), r.length(attrStr(n, "height"))
	if v := r.length(css["width"]); v > 0 {
		aw = v
	}
	if v := r.length(css["height"]); v > 0 {
		ah = v
	}
	switch {
	case aw > 0 && ah > 0:
		w, h = aw, ah
	case aw > 0:
		if w > 0 {
			h *= aw / w
		}
		w = aw
	case ah > 0:
		if h > 0 {
			w *= ah / h
		}
		h = ah
	}
	if w <= 0 || h <= 0 {
		// the size of a replaced element that says none
		w, h = 300*pxToPt, 150*pxToPt
	}
	if !hasAlt && st.figAlt != "" {
		alt = st.figAlt
	}
	hash := im.hash
	o := &inlineObj{w: w, h: h, alt: alt, fit: true, paint: func(e *emitter, box drawingml.Box) {
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
		r.c.warnOnce("svgimg", "SVG images are not drawn (%s)", name)
		return nil
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

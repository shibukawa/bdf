// Package jpx decodes JPEG 2000 images (ITU-T T.800 | ISO/IEC 15444-1):
// bare codestreams and JP2/JPX files, the data of the JPXDecode filter of
// PDF.
//
// It implements Part 1 apart from reduced-resolution and layer-limited
// decoding: tiles and tile-parts, image and tile offsets, component
// subsampling, all progression orders and progression order changes,
// precincts, packed packet headers, SOP and EPH markers, every code-block
// style, both wavelet filters, scalar quantization, region of interest
// (maxshift) and the multiple component transforms. A truncated codestream
// decodes to the extent its data allows.
package jpx

import (
	"bytes"
	"fmt"
	"slices"
)

// Image is a decoded JPEG 2000 image.
type Image struct {
	// Width and Height are the size of the image area on the reference
	// grid (Xsiz−XOsiz, Ysiz−YOsiz).
	Width, Height int
	// Components are the channels, each resampled to Width×Height;
	// subsampled components are upsampled by pixel replication. With a
	// palette they are the palette's output channels. With a channel
	// definition box the colour channels come first in association order,
	// then the opacity channel, then channels of no stated use.
	Components []Component
	// ColorSpace is the enumerated colour space of the JP2 colr box (16
	// sRGB, 17 greyscale, 18 sYCC, 12 CMYK, …); 0 when there is none or the
	// colour space is an ICC profile.
	ColorSpace int
	// ICC is the ICC profile of the colr box, if any.
	ICC []byte
	// Alpha is the index of the opacity channel in Components, or -1.
	Alpha int
	// Premultiplied reports whether the colour channels are premultiplied
	// by the opacity.
	Premultiplied bool
}

// Component is a channel of an Image.
type Component struct {
	Precision int
	Signed    bool
	// Data holds Width*Height samples in row-major order, in the range
	// 0..2^Precision−1, or −2^(Precision−1)..2^(Precision−1)−1 when Signed.
	Data []int32
}

func errorf(format string, args ...any) error { return fmt.Errorf("jpx: "+format, args...) }

// Decode decodes a JPEG 2000 codestream (starting with the SOC marker) or a
// JP2/JPX file (starting with the signature box).
func Decode(data []byte) (img *Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			img, err = nil, errorf("internal error: %v", r)
		}
	}()
	var h *jp2Header
	cs := data
	if bytes.HasPrefix(data, jp2Signature) {
		if h, cs, err = parseJP2(data); err != nil {
			return nil, err
		}
	}
	c, err := parseCodestream(cs)
	if err != nil {
		return nil, err
	}
	chans, err := channels(c, h)
	if err != nil {
		return nil, err
	}
	planes, err := c.decode()
	if err != nil {
		return nil, err
	}
	img = &Image{Width: c.x1 - c.x0, Height: c.y1 - c.y0, Alpha: -1}
	if h != nil {
		img.ColorSpace, img.ICC = h.colorSpace, h.icc
	}
	// A plane that feeds a single channel directly becomes its data.
	uses := make([]int, len(planes))
	for _, ch := range chans {
		uses[ch.comp]++
	}
	for i, ch := range chans {
		if ch.alpha {
			img.Alpha, img.Premultiplied = i, ch.premult
		}
		img.Components = append(img.Components, ch.render(c, planes, uses[ch.comp] == 1))
	}
	return img, nil
}

// channel is an output channel: a component, possibly mapped through a
// palette column.
type channel struct {
	comp    int
	pal     *palette
	col     int
	alpha   bool
	premult bool
}

// channels applies the palette, component mapping and channel definition
// boxes.
func channels(c *codestream, h *jp2Header) ([]channel, error) {
	var chans []channel
	if h != nil && h.pal != nil && h.cmap != nil {
		for _, m := range h.cmap {
			if m.comp >= len(c.comps) || m.typ > 1 || m.typ == 1 && m.col >= len(h.pal.cols) {
				return nil, errorf("invalid cmap box")
			}
			ch := channel{comp: m.comp}
			if m.typ == 1 {
				ch.pal, ch.col = h.pal, m.col
			}
			chans = append(chans, ch)
		}
	} else {
		for i := range c.comps {
			chans = append(chans, channel{comp: i})
		}
	}
	var total int64
	for range chans {
		total += int64(c.x1-c.x0) * int64(c.y1-c.y0)
	}
	if total > maxSamples {
		return nil, errorf("image too large")
	}
	if h == nil || h.cdef == nil {
		return chans, nil
	}
	// Colour channels by association, then opacity, then the rest.
	type colour struct{ ch, asoc int }
	var colours []colour
	var rest []int
	opacity := -1
	seen := make([]bool, len(chans))
	for _, d := range h.cdef {
		if d.ch >= len(chans) || seen[d.ch] {
			continue
		}
		seen[d.ch] = true
		switch {
		case d.typ == 0 && d.asoc > 0 && d.asoc < 0xffff:
			colours = append(colours, colour{d.ch, d.asoc})
		case (d.typ == 1 || d.typ == 2) && opacity < 0:
			opacity = d.ch
			chans[d.ch].alpha, chans[d.ch].premult = true, d.typ == 2
		default:
			rest = append(rest, d.ch)
		}
	}
	for i, s := range seen {
		if !s {
			rest = append(rest, i)
		}
	}
	slices.SortStableFunc(colours, func(a, b colour) int { return a.asoc - b.asoc })
	out := make([]channel, 0, len(chans))
	for _, cl := range colours {
		out = append(out, chans[cl.ch])
	}
	if opacity >= 0 {
		out = append(out, chans[opacity])
	}
	for _, i := range rest {
		out = append(out, chans[i])
	}
	return out, nil
}

// render produces the samples of a channel at the full image size. If own
// is set, the plane is not needed elsewhere and its data may be reused.
func (ch channel) render(c *codestream, planes []plane, own bool) Component {
	cp := c.comps[ch.comp]
	pl := &planes[ch.comp]
	out := Component{Precision: cp.prec, Signed: cp.signed}
	src := pl.data
	if ch.pal != nil {
		pc := ch.pal.cols[ch.col]
		out.Precision, out.Signed = pc.prec, pc.signed
		n := len(ch.pal.cols)
		src = make([]int32, len(pl.data))
		for i, v := range pl.data {
			e := min(max(int(v), 0), ch.pal.entries-1)
			src[i] = ch.pal.values[e*n+ch.col]
		}
	}
	w, h := c.x1-c.x0, c.y1-c.y0
	if pl.w == 0 || pl.h == 0 {
		// Subsampling leaves the component without samples.
		out.Data = make([]int32, w*h)
		return out
	}
	if cp.dx == 1 && cp.dy == 1 {
		// The plane covers the image area exactly.
		if ch.pal == nil && !own {
			src = slices.Clone(src)
		}
		out.Data = src
		return out
	}
	cols := make([]int, w)
	for x := range cols {
		cols[x] = min(max((c.x0+x)/cp.dx-pl.x0, 0), pl.w-1)
	}
	out.Data = make([]int32, w*h)
	for y := range h {
		sy := min(max((c.y0+y)/cp.dy-pl.y0, 0), pl.h-1)
		row := src[sy*pl.w : sy*pl.w+pl.w]
		dst := out.Data[y*w : y*w+w]
		for x, sx := range cols {
			dst[x] = row[sx]
		}
	}
	return out
}

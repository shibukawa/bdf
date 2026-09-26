package psd

import (
	"image"
	"math"
	"slices"
)

// The composite of a file saved without "Maximize Compatibility" is not in
// the file: the converter composites the layers itself. It draws pixel
// layers (text and smart object layers keep their rendered pixels) and
// solid colour fill layers with their opacity, fill opacity, blend mode,
// layer masks and clipping masks, groups (isolated or pass-through) and
// artboards with their backgrounds. It does not apply adjustment layers,
// layer effects, gradient and pattern fills or vector masks that were not
// rendered into a mask channel; the result then differs from Photoshop's,
// and a warning says so.

// node is a layer or a group of the layer tree.
type node struct {
	l        *layer
	group    bool
	children []*node // bottom to top
}

// layerTree nests the layers (stored bottom to top, a group between a
// bounding divider below its layers and the group's own record above them).
func layerTree(layers []*layer) []*node {
	stack := [][]*node{nil}
	for _, l := range layers {
		top := len(stack) - 1
		switch l.section {
		case 3:
			stack = append(stack, nil)
		case 1, 2:
			var children []*node
			if top > 0 {
				children = stack[top]
				stack = stack[:top]
				top--
			}
			stack[top] = append(stack[top], &node{l: l, group: true, children: children})
		default:
			stack[top] = append(stack[top], &node{l: l})
		}
	}
	for len(stack) > 1 { // dividers without a group: keep their layers
		top := len(stack) - 1
		stack[top-1] = append(stack[top-1], stack[top]...)
		stack = stack[:top]
	}
	return stack[0]
}

// rgba is premultiplied RGBA over a rectangle of the canvas; outside it
// everything is transparent.
type rgba struct {
	r   image.Rectangle
	pix []float32
}

func newRGBA(r image.Rectangle) *rgba {
	return &rgba{r: r, pix: make([]float32, 4*r.Dx()*r.Dy())}
}

func (b *rgba) at(x, y int) int { return 4 * ((y-b.r.Min.Y)*b.r.Dx() + x - b.r.Min.X) }

type compositor struct {
	f      *file
	canvas image.Rectangle
	warn   func(string)
}

// adjustmentKeys are the additional information keys of adjustment layers.
var adjustmentKeys = []string{"levl", "curv", "brit", "blnc", "hue ", "hue2", "selc", "mixr", "grdm", "phfl", "expA",
	"vibA", "thrs", "nvrt", "post", "CgEd", "clrL", "blwh"}

// compositeLayers draws the layer tree into an image the size of the canvas.
func (f *file) compositeLayers(warn func(string)) *image.NRGBA {
	cp := &compositor{f: f, canvas: image.Rect(0, 0, f.hdr.w, f.hdr.h), warn: warn}
	dst := newRGBA(cp.canvas)
	cp.drawList(dst, layerTree(f.layers))
	img := image.NewNRGBA(cp.canvas)
	for i := 0; i < len(dst.pix); i += 4 {
		a := dst.pix[i+3]
		if a <= 0 {
			continue
		}
		for c := 0; c < 3; c++ {
			img.Pix[i+c] = clamp8(float64(dst.pix[i+c]/a) * 255)
		}
		img.Pix[i+3] = clamp8(float64(a) * 255)
	}
	return img
}

// drawList draws layers bottom to top. A layer with clipping set is clipped
// to the nearest layer below without it, its base.
func (cp *compositor) drawList(dst *rgba, nodes []*node) {
	for i := 0; i < len(nodes); {
		j := i + 1
		for j < len(nodes) && nodes[j].l.clipping {
			j++
		}
		base, clips := nodes[i], nodes[i+1:j]
		i = j
		if base.l.hidden {
			continue // a hidden base hides what is clipped to it
		}
		shown := slices.ContainsFunc(clips, func(n *node) bool { return !n.l.hidden })
		if !shown {
			cp.drawNode(dst, base)
			continue
		}
		// The base and the layers clipped to it are composited alone, the
		// clipped ones inside the base's pixels, then drawn with the base's
		// blend mode and opacity.
		buf := cp.render(base)
		if buf == nil {
			continue
		}
		for _, c := range clips {
			if !c.l.hidden {
				if src := cp.render(c); src != nil {
					blend(buf, src, c.l.blend, opacity(c.l), true)
				}
			}
		}
		blend(dst, buf, base.l.blend, opacity(base.l), false)
	}
}

func opacity(l *layer) float32 { return float32(l.opacity) / 255 }

// drawNode draws a layer or a group (not clipped) onto dst.
func (cp *compositor) drawNode(dst *rgba, n *node) {
	if n.group && n.l.artboard == nil && (n.l.sectionBlend == "pass" || n.l.blend == "pass") {
		// A pass-through group's layers blend with what is below it; its
		// opacity and mask then fade the result into what was there.
		m := cp.mask(n.l)
		if n.l.opacity == 255 && m == nil {
			cp.drawList(dst, n.children)
			return
		}
		t := &rgba{r: dst.r, pix: slices.Clone(dst.pix)}
		cp.drawList(t, n.children)
		k := opacity(n.l)
		for y := dst.r.Min.Y; y < dst.r.Max.Y; y++ {
			for x := dst.r.Min.X; x < dst.r.Max.X; x++ {
				kk := k * m.value(x, y)
				i := dst.at(x, y)
				for c := 0; c < 4; c++ {
					dst.pix[i+c] += (t.pix[i+c] - dst.pix[i+c]) * kk
				}
			}
		}
		return
	}
	if src := cp.render(n); src != nil {
		blend(dst, src, n.l.blend, opacity(n.l), false)
	}
}

// render draws a layer or group alone, with its fill opacity and masks but
// not its opacity or blend mode.
func (cp *compositor) render(n *node) *rgba {
	l := n.l
	var buf *rgba
	if n.group {
		buf = newRGBA(cp.canvas)
		if a := l.artboard; a != nil {
			r := image.Rect(a.left, a.top, a.right, a.bottom).Intersect(cp.canvas)
			var bg [4]float32
			switch a.background {
			case 1:
				bg = [4]float32{1, 1, 1, 1}
			case 2:
				bg = [4]float32{0, 0, 0, 1}
			case 4:
				bg = [4]float32{float32(a.color[0] / 255), float32(a.color[1] / 255), float32(a.color[2] / 255), 1}
			}
			for y := r.Min.Y; y < r.Max.Y; y++ {
				for x := r.Min.X; x < r.Max.X; x++ {
					copy(buf.pix[buf.at(x, y):], bg[:])
				}
			}
			cp.drawList(buf, n.children)
			// An artboard clips its layers.
			for y := buf.r.Min.Y; y < buf.r.Max.Y; y++ {
				for x := buf.r.Min.X; x < buf.r.Max.X; x++ {
					if !(image.Point{x, y}).In(r) {
						clear(buf.pix[buf.at(x, y) : buf.at(x, y)+4])
					}
				}
			}
		} else {
			cp.drawList(buf, n.children)
		}
	} else {
		buf = cp.renderLayer(l)
		if buf == nil {
			return nil
		}
		if l.fill < 255 {
			k := float32(l.fill) / 255
			for i := range buf.pix {
				buf.pix[i] *= k
			}
		}
	}
	if m := cp.mask(l); m != nil {
		for y := buf.r.Min.Y; y < buf.r.Max.Y; y++ {
			for x := buf.r.Min.X; x < buf.r.Max.X; x++ {
				k := m.value(x, y)
				i := buf.at(x, y)
				for c := 0; c < 4; c++ {
					buf.pix[i+c] *= k
				}
			}
		}
	}
	return buf
}

// renderLayer returns the pixels of a layer, or nil when it draws none.
func (cp *compositor) renderLayer(l *layer) *rgba {
	f := cp.f
	for _, k := range adjustmentKeys {
		if slices.Contains(l.keys, k) {
			cp.warn("adjustment layers are not applied")
			return nil
		}
	}
	if l.effects {
		cp.warn("layer effects (shadows, strokes, glows...) are not drawn")
	}
	lr := image.Rect(l.left, l.top, l.right, l.bottom)
	if !cp.plausible(lr) {
		cp.warn("layer " + l.name + ": bad bounds")
		return nil
	}
	planes := make([][]uint8, f.hdr.baseChannels())
	var alpha []uint8
	var err error
	have := false
	for _, ch := range l.channels {
		switch {
		case ch.id >= 0 && ch.id < len(planes):
			planes[ch.id], err = f.layerPlane(ch, lr.Dx(), lr.Dy(), true)
			have = have || planes[ch.id] != nil
		case ch.id == chAlpha:
			alpha, err = f.layerPlane(ch, lr.Dx(), lr.Dy(), false)
		}
		if err != nil {
			cp.warn("layer " + l.name + ": " + err.Error())
			return nil
		}
	}
	if !have {
		if l.fillColor == nil {
			if slices.Contains(l.keys, "GdFl") || slices.Contains(l.keys, "PtFl") {
				cp.warn("gradient and pattern fill layers are not drawn")
			}
			return nil
		}
		// A solid colour fill covers the canvas; its masks shape it.
		if l.mask == nil && (slices.Contains(l.keys, "vmsk") || slices.Contains(l.keys, "vsms")) {
			cp.warn("shape layers whose vector mask is not rendered into a mask channel are not drawn")
			return nil
		}
		buf := newRGBA(cp.canvas)
		c := [4]float32{float32(l.fillColor[0] / 255), float32(l.fillColor[1] / 255), float32(l.fillColor[2] / 255), 1}
		for i := 0; i < len(buf.pix); i += 4 {
			copy(buf.pix[i:], c[:])
		}
		return buf
	}
	for i := range planes {
		if planes[i] == nil {
			planes[i] = make([]uint8, lr.Dx()*lr.Dy())
		}
	}
	img := f.toNRGBA(planes, alpha, lr.Dx(), lr.Dy())
	r := lr.Intersect(cp.canvas)
	if r.Empty() {
		return nil
	}
	buf := newRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s := img.PixOffset(x-lr.Min.X, y-lr.Min.Y)
			a := float32(img.Pix[s+3]) / 255
			i := buf.at(x, y)
			buf.pix[i] = float32(img.Pix[s]) / 255 * a
			buf.pix[i+1] = float32(img.Pix[s+1]) / 255 * a
			buf.pix[i+2] = float32(img.Pix[s+2]) / 255 * a
			buf.pix[i+3] = a
		}
	}
	return buf
}

// plausible reports whether the bounds of a layer or mask are within reach
// of the canvas: layers may extend past it, not endlessly.
func (cp *compositor) plausible(r image.Rectangle) bool {
	w, h := cp.canvas.Dx(), cp.canvas.Dy()
	return r.Min.X >= -4*w && r.Max.X <= 5*w && r.Min.Y >= -4*h && r.Max.Y <= 5*h
}

// maskPlane is a decoded layer mask.
type maskPlane struct {
	r     image.Rectangle
	pix   []uint8
	deflt float32
	next  *maskPlane // a second mask multiplied in
}

func (m *maskPlane) value(x, y int) float32 {
	if m == nil {
		return 1
	}
	v := m.deflt
	if (image.Point{x, y}).In(m.r) {
		v = float32(m.pix[(y-m.r.Min.Y)*m.r.Dx()+x-m.r.Min.X]) / 255
	}
	return v * m.next.value(x, y)
}

// mask returns the enabled masks of a layer, nil when it has none.
func (cp *compositor) mask(l *layer) *maskPlane {
	var out *maskPlane
	for _, ch := range l.channels {
		m := l.mask
		if ch.id == chRealMask {
			m = l.realMask
		} else if ch.id != chMask {
			continue
		}
		if m == nil || m.disabled() {
			continue
		}
		r := image.Rect(m.left, m.top, m.right, m.bottom)
		if !cp.plausible(r) {
			cp.warn("layer " + l.name + " mask: bad bounds")
			continue
		}
		pix, err := cp.f.layerPlane(ch, r.Dx(), r.Dy(), false)
		if err != nil {
			cp.warn("layer " + l.name + " mask: " + err.Error())
			continue
		}
		if pix == nil {
			r = image.Rectangle{}
		}
		out = &maskPlane{r: r, pix: pix, deflt: float32(m.defaultColor) / 255, next: out}
	}
	return out
}

// blend composites src onto dst with a blend mode and an opacity (W3C
// Compositing and Blending). atop keeps dst's alpha: the source is clipped
// to what is there.
func blend(dst, src *rgba, mode string, opacity float32, atop bool) {
	r := src.r.Intersect(dst.r)
	fn := blendFuncs[mode]
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			s, d := src.at(x, y), dst.at(x, y)
			sa := src.pix[s+3] * opacity
			if sa <= 0 {
				continue
			}
			ba := dst.pix[d+3]
			var cs, cb [3]float32
			for c := 0; c < 3; c++ {
				cs[c] = src.pix[s+c] / src.pix[s+3]
				if ba > 0 {
					cb[c] = dst.pix[d+c] / ba
				}
			}
			mixed := cs
			if fn != nil && ba > 0 {
				mixed = fn(cb, cs)
			}
			for c := 0; c < 3; c++ {
				v := sa*ba*mixed[c] + (1-sa)*ba*cb[c]
				if !atop {
					v += sa * (1 - ba) * cs[c]
				}
				dst.pix[d+c] = v
			}
			if !atop {
				dst.pix[d+3] = sa + ba*(1-sa)
			}
		}
	}
}

type blendFunc func(cb, cs [3]float32) [3]float32

func separable(f func(cb, cs float32) float32) blendFunc {
	return func(cb, cs [3]float32) [3]float32 {
		return [3]float32{f(cb[0], cs[0]), f(cb[1], cs[1]), f(cb[2], cs[2])}
	}
}

func colorDodge(cb, cs float32) float32 {
	switch {
	case cb <= 0:
		return 0
	case cs >= 1:
		return 1
	}
	return min(1, cb/(1-cs))
}

func colorBurn(cb, cs float32) float32 {
	switch {
	case cb >= 1:
		return 1
	case cs <= 0:
		return 0
	}
	return 1 - min(1, (1-cb)/cs)
}

func hardLight(cb, cs float32) float32 {
	if cs <= 0.5 {
		return cb * 2 * cs
	}
	s := 2*cs - 1
	return cb + s - cb*s
}

func vividLight(cb, cs float32) float32 {
	if cs <= 0.5 {
		return colorBurn(cb, 2*cs)
	}
	return colorDodge(cb, 2*cs-1)
}

func clamp01(v float32) float32 { return max(0, min(1, v)) }

// blendFuncs maps Photoshop's blend mode keys to their functions; normal,
// dissolve and unknown modes draw the source colour.
var blendFuncs = map[string]blendFunc{
	"mul ": separable(func(cb, cs float32) float32 { return cb * cs }),
	"scrn": separable(func(cb, cs float32) float32 { return cb + cs - cb*cs }),
	"over": separable(func(cb, cs float32) float32 { return hardLight(cs, cb) }),
	"dark": separable(func(cb, cs float32) float32 { return min(cb, cs) }),
	"lite": separable(func(cb, cs float32) float32 { return max(cb, cs) }),
	"div ": separable(colorDodge),
	"idiv": separable(colorBurn),
	"hLit": separable(hardLight),
	"sLit": separable(func(cb, cs float32) float32 {
		if cs <= 0.5 {
			return cb - (1-2*cs)*cb*(1-cb)
		}
		return cb + (2*cs-1)*(float32(math.Sqrt(float64(cb)))-cb)
	}),
	"diff": separable(func(cb, cs float32) float32 { return float32(math.Abs(float64(cb - cs))) }),
	"smud": separable(func(cb, cs float32) float32 { return cb + cs - 2*cb*cs }),
	"lbrn": separable(func(cb, cs float32) float32 { return clamp01(cb + cs - 1) }),
	"lddg": separable(func(cb, cs float32) float32 { return clamp01(cb + cs) }),
	"vLit": separable(vividLight),
	"lLit": separable(func(cb, cs float32) float32 { return clamp01(cb + 2*cs - 1) }),
	"pLit": separable(func(cb, cs float32) float32 {
		if cs <= 0.5 {
			return min(cb, 2*cs)
		}
		return max(cb, 2*cs-1)
	}),
	"hMix": separable(func(cb, cs float32) float32 {
		if cb+cs >= 1 {
			return 1
		}
		return 0
	}),
	"fsub": separable(func(cb, cs float32) float32 { return clamp01(cb - cs) }),
	"fdiv": separable(func(cb, cs float32) float32 {
		if cs <= 0 {
			if cb <= 0 {
				return 0
			}
			return 1
		}
		return min(1, cb/cs)
	}),
	"hue ": func(cb, cs [3]float32) [3]float32 { return setLum(setSat(cs, sat(cb)), lum(cb)) },
	"sat ": func(cb, cs [3]float32) [3]float32 { return setLum(setSat(cb, sat(cs)), lum(cb)) },
	"colr": func(cb, cs [3]float32) [3]float32 { return setLum(cs, lum(cb)) },
	"lum ": func(cb, cs [3]float32) [3]float32 { return setLum(cb, lum(cs)) },
	"dkCl": func(cb, cs [3]float32) [3]float32 {
		if lum(cs) < lum(cb) {
			return cs
		}
		return cb
	},
	"lgCl": func(cb, cs [3]float32) [3]float32 {
		if lum(cs) > lum(cb) {
			return cs
		}
		return cb
	},
}

func lum(c [3]float32) float32 { return 0.3*c[0] + 0.59*c[1] + 0.11*c[2] }

func setLum(c [3]float32, l float32) [3]float32 {
	d := l - lum(c)
	for i := range c {
		c[i] += d
	}
	l = lum(c)
	n, x := min(c[0], c[1], c[2]), max(c[0], c[1], c[2])
	for i := range c {
		if n < 0 {
			c[i] = l + (c[i]-l)*l/(l-n)
		}
		if x > 1 {
			c[i] = l + (c[i]-l)*(1-l)/(x-l)
		}
	}
	return c
}

func sat(c [3]float32) float32 { return max(c[0], c[1], c[2]) - min(c[0], c[1], c[2]) }

func setSat(c [3]float32, s float32) [3]float32 {
	idx := []int{0, 1, 2}
	slices.SortFunc(idx, func(a, b int) int {
		switch {
		case c[a] < c[b]:
			return -1
		case c[a] > c[b]:
			return 1
		}
		return 0
	})
	lo, mid, hi := idx[0], idx[1], idx[2]
	var out [3]float32
	if c[hi] > c[lo] {
		out[mid] = (c[mid] - c[lo]) * s / (c[hi] - c[lo])
		out[hi] = s
	}
	return out
}

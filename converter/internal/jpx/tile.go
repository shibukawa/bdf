package jpx

import (
	"errors"
	"math"
	"slices"
)

// Tile decoding: the geometry of tile-components, resolution levels,
// sub-bands, precincts and code-blocks (Annex B), then Tier-2, Tier-1,
// dequantization (Annex E), the inverse wavelet transform, the inverse
// component transform (Annex G) and the DC level shift.

// Sub-band types; they are also the order of the bands in a packet.
const (
	bandLL = iota
	bandHL
	bandLH
	bandHH
)

type band struct {
	typ            int
	x0, y0, x1, y1 int // band coordinates
	offX, offY     int // origin in the tile-component buffer
	cbw, cbh       int // code-block size exponents
	mb             int // magnitude bit-planes (E-2)
	delta          float32
}

type codeBlock struct {
	x0, y0, x1, y1 int
	zbp            int // missing most significant bit-planes
	lblock         int
	segs           []segment // nil until the block is first included
	data           []byte
}

type precBand struct {
	blocks    []codeBlock
	incl, zbp tagTree
}

type precinct struct {
	done  int // layers decoded
	bands []precBand
}

type resolution struct {
	x0, y0, x1, y1 int
	ppx, ppy       int
	npx, npy       int
	bands          []band
	precincts      []precinct
}

type tileComp struct {
	x0, y0, x1, y1 int
	reversible     bool
	style          byte
	roi            int
	res            []resolution
	ints           []int32
	flts           []float32
}

type tileDecoder struct {
	c              *codestream
	index          int
	x0, y0, x1, y1 int
	dx, dy         []int
	comps          []tileComp
	layers         int
	sop, eph       bool
	pocs           []poc

	body    []byte // packet data
	bpos    int
	hdr     []byte // packet headers, when they come from PPM or PPT
	hpos    int
	pending []pendingData
	blocks  int

	t1 t1Decoder
}

// plane is a decoded component at its own resolution.
type plane struct {
	x0, y0, w, h int
	data         []int32
}

// decode decodes all tiles into component planes. Samples of tiles without
// data are left zero.
func (c *codestream) decode() ([]plane, error) {
	planes := make([]plane, len(c.comps))
	for i, cp := range c.comps {
		x0, y0 := ceilDiv(c.x0, cp.dx), ceilDiv(c.y0, cp.dy)
		planes[i] = plane{x0: x0, y0: y0, w: ceilDiv(c.x1, cp.dx) - x0, h: ceilDiv(c.y1, cp.dy) - y0}
	}
	found := false
	for i, t := range c.tiles {
		if t == nil || len(t.parts) == 0 {
			continue
		}
		found = true
		if err := c.decodeTile(i, t, planes); err != nil {
			return nil, err
		}
		c.tiles[i] = nil // release the tile data
	}
	if !found {
		return nil, errorf("no tile data")
	}
	for i := range planes {
		if planes[i].data == nil {
			planes[i].data = make([]int32, planes[i].w*planes[i].h)
		}
	}
	return planes, nil
}

func (c *codestream) decodeTile(ti int, t *tile, planes []plane) error {
	td := &tileDecoder{c: c, index: ti}
	p, q := ti%c.ntx, ti/c.ntx
	td.x0 = max(c.tx0+p*c.tw, c.x0)
	td.x1 = min(c.tx0+(p+1)*c.tw, c.x1)
	td.y0 = max(c.ty0+q*c.th, c.y0)
	td.y1 = min(c.ty0+(q+1)*c.th, c.y1)

	cd := c.main.cod
	if t.hdr.cod != nil {
		cd = t.hdr.cod
	}
	td.layers, td.sop, td.eph = cd.layers, cd.sop, cd.eph
	td.comps = make([]tileComp, len(c.comps))
	td.dx = make([]int, len(c.comps))
	td.dy = make([]int, len(c.comps))
	maxRes := 0
	precincts := 0
	for i, cp := range c.comps {
		td.dx[i], td.dy[i] = cp.dx, cp.dy
		cs := &c.main.cod.cs
		switch {
		case t.hdr.coc != nil && t.hdr.coc[i] != nil:
			cs = t.hdr.coc[i]
		case t.hdr.cod != nil:
			cs = &t.hdr.cod.cs
		case c.main.coc != nil && c.main.coc[i] != nil:
			cs = c.main.coc[i]
		}
		qt := c.main.qcd
		switch {
		case t.hdr.qcc != nil && t.hdr.qcc[i] != nil:
			qt = t.hdr.qcc[i]
		case t.hdr.qcd != nil:
			qt = t.hdr.qcd
		case c.main.qcc != nil && c.main.qcc[i] != nil:
			qt = c.main.qcc[i]
		}
		roi := 0
		if t.hdr.rgn != nil && t.hdr.rgn[i] >= 0 {
			roi = t.hdr.rgn[i]
		} else if c.main.rgn != nil && c.main.rgn[i] >= 0 {
			roi = c.main.rgn[i]
		}
		n, err := td.setupComp(&td.comps[i], cp, cs, qt, roi)
		if err != nil {
			return err
		}
		if precincts += n; precincts > maxPrecincts {
			return errorf("too many precincts")
		}
		maxRes = max(maxRes, cs.levels+1)
	}
	td.pocs = t.hdr.pocs
	if td.pocs == nil {
		td.pocs = c.main.pocs
	}
	if td.pocs == nil {
		td.pocs = []poc{{rs: 0, cs: 0, lye: td.layers, re: maxRes, ce: len(c.comps), prog: cd.prog}}
	}

	td.body = t.parts[0]
	if len(t.parts) > 1 {
		td.body = slices.Concat(t.parts...)
	}
	switch {
	case c.ppm:
		td.hdr = slices.Concat(t.ppm...)
		if td.hdr == nil {
			td.hdr = []byte{} // no headers left for this tile
		}
	case t.ppt != nil:
		slices.SortStableFunc(t.ppt, func(a, b zseg) int { return a.z - b.z })
		td.hdr = []byte{}
		for _, s := range t.ppt {
			td.hdr = append(td.hdr, s.data...)
		}
	}
	if err := td.readPackets(); err != nil && !errors.Is(err, errStop) {
		return err
	}
	td.body, td.hdr = nil, nil

	for i := range td.comps {
		td.decodeComp(&td.comps[i])
	}
	if cd.mct == 1 && len(td.comps) >= 3 {
		td.inverseMCT()
	}
	for i := range td.comps {
		td.store(&td.comps[i], c.comps[i], &planes[i])
	}
	return nil
}

// setupComp lays out the resolution levels, sub-bands, precincts and
// code-blocks of a tile-component and returns its number of precincts.
func (td *tileDecoder) setupComp(tc *tileComp, cp component, cs *codingStyle, qt *quant, roi int) (int, error) {
	tc.x0, tc.y0 = ceilDiv(td.x0, cp.dx), ceilDiv(td.y0, cp.dy)
	tc.x1, tc.y1 = ceilDiv(td.x1, cp.dx), ceilDiv(td.y1, cp.dy)
	tc.reversible, tc.style, tc.roi = cs.reversible, cs.style, roi
	nl := cs.levels
	tc.res = make([]resolution, nl+1)
	count := 0
	for r := range tc.res {
		res := &tc.res[r]
		lev := nl - r
		res.x0, res.y0 = ceilDivPow2(tc.x0, lev), ceilDivPow2(tc.y0, lev)
		res.x1, res.y1 = ceilDivPow2(tc.x1, lev), ceilDivPow2(tc.y1, lev)
		res.ppx, res.ppy = 15, 15
		if cs.ppx != nil {
			res.ppx, res.ppy = int(cs.ppx[r]), int(cs.ppy[r])
		}
		if res.x1 > res.x0 && res.y1 > res.y0 {
			res.npx = ceilDivPow2(res.x1, res.ppx) - res.x0>>res.ppx
			res.npy = ceilDivPow2(res.y1, res.ppy) - res.y0>>res.ppy
		}
		n := int64(res.npx) * int64(res.npy)
		if n > maxPrecincts {
			return 0, errorf("too many precincts")
		}
		count += int(n)

		// Sub-bands.
		types := []int{bandLL}
		if r > 0 {
			types = []int{bandHL, bandLH, bandHH}
		}
		res.bands = make([]band, len(types))
		for bi, typ := range types {
			b := &res.bands[bi]
			b.typ = typ
			nb := nl - r + 1 // decomposition level of the band
			if r == 0 {
				nb = nl
			}
			xo, yo := typ&1, typ>>1
			b.x0, b.x1 = bandCoord(tc.x0, nb, xo), bandCoord(tc.x1, nb, xo)
			b.y0, b.y1 = bandCoord(tc.y0, nb, yo), bandCoord(tc.y1, nb, yo)
			if r > 0 {
				low := &tc.res[r-1]
				b.offX = xo * (low.x1 - low.x0)
				b.offY = yo * (low.y1 - low.y0)
			}
			bpx, bpy := res.ppx, res.ppy
			if r > 0 {
				bpx, bpy = bpx-1, bpy-1
			}
			b.cbw, b.cbh = min(cs.cbw, bpx), min(cs.cbh, bpy)

			si := 0
			if r > 0 {
				si = 3*(r-1) + bi + 1
			}
			st := qt.step(si, r)
			b.mb = qt.guard + st.exp - 1
			gain := [4]int{0, 1, 1, 2}[typ]
			b.delta = float32(math.Ldexp(1+float64(st.mant)/2048, cp.prec+gain-st.exp))
		}
	}
	if count > maxPrecincts {
		return 0, errorf("too many precincts")
	}
	// Precincts and code-blocks.
	for r := range tc.res {
		res := &tc.res[r]
		res.precincts = make([]precinct, res.npx*res.npy)
		for py := range res.npy {
			for px := range res.npx {
				p := &res.precincts[py*res.npx+px]
				p.bands = make([]precBand, len(res.bands))
				// Precinct origin in the resolution level, then in the bands.
				x0 := (res.x0>>res.ppx + px) << res.ppx
				y0 := (res.y0>>res.ppy + py) << res.ppy
				bpx, bpy := res.ppx, res.ppy
				if r > 0 {
					x0, y0, bpx, bpy = x0>>1, y0>>1, bpx-1, bpy-1
				}
				for bi := range res.bands {
					if err := td.setupPrecBand(&p.bands[bi], &res.bands[bi], x0, y0, bpx, bpy); err != nil {
						return 0, err
					}
				}
			}
		}
	}
	return count, nil
}

// bandCoord maps a tile-component coordinate to sub-band coordinates:
// ceil((t - o*2^(nb-1)) / 2^nb) (B-15).
func bandCoord(t, nb, o int) int {
	if o == 0 {
		return ceilDivPow2(t, nb)
	}
	if nb >= 31 {
		return 0
	}
	return (t + 1<<(nb-1) - 1) >> nb
}

func (td *tileDecoder) setupPrecBand(pb *precBand, b *band, px0, py0, bpx, bpy int) error {
	x0, x1 := max(px0, b.x0), min(px0+1<<bpx, b.x1)
	y0, y1 := max(py0, b.y0), min(py0+1<<bpy, b.y1)
	if x0 >= x1 || y0 >= y1 {
		return nil
	}
	cx0, cx1 := x0>>b.cbw, ceilDivPow2(x1, b.cbw)
	cy0, cy1 := y0>>b.cbh, ceilDivPow2(y1, b.cbh)
	cw, ch := cx1-cx0, cy1-cy0
	if td.blocks += cw * ch; td.blocks > maxBlocks {
		return errorf("too many code-blocks")
	}
	pb.blocks = make([]codeBlock, cw*ch)
	for j := range ch {
		for i := range cw {
			cb := &pb.blocks[j*cw+i]
			cb.x0 = max((cx0+i)<<b.cbw, x0)
			cb.x1 = min((cx0+i+1)<<b.cbw, x1)
			cb.y0 = max((cy0+j)<<b.cbh, y0)
			cb.y1 = min((cy0+j+1)<<b.cbh, y1)
		}
	}
	pb.incl = newTagTree(cw, ch)
	pb.zbp = newTagTree(cw, ch)
	return nil
}

// decodeComp decodes the code-blocks of a tile-component into its buffer and
// applies the inverse wavelet transform. A tile-component without coded
// data keeps no buffer: its coefficients are all zero.
func (td *tileDecoder) decodeComp(tc *tileComp) {
	w, h := tc.x1-tc.x0, tc.y1-tc.y0
	if w <= 0 || h <= 0 {
		return
	}
	for r := range tc.res {
		res := &tc.res[r]
		for pi := range res.precincts {
			for bi := range res.precincts[pi].bands {
				pb := &res.precincts[pi].bands[bi]
				for k := range pb.blocks {
					cb := &pb.blocks[k]
					if len(cb.segs) > 0 {
						if tc.ints == nil && tc.flts == nil {
							tc.alloc()
						}
						td.decodeBlock(tc, &res.bands[bi], cb, w)
					}
					cb.data, cb.segs = nil, nil
				}
			}
		}
	}
	if tc.ints == nil && tc.flts == nil {
		return
	}
	levels := make([]level, len(tc.res)-1)
	for r := 1; r < len(tc.res); r++ {
		res, low := &tc.res[r], &tc.res[r-1]
		levels[r-1] = level{
			w: res.x1 - res.x0, h: res.y1 - res.y0,
			lw: low.x1 - low.x0, lh: low.y1 - low.y0,
			px: res.x0 & 1, py: res.y0 & 1,
		}
	}
	if tc.reversible {
		idwt53(tc.ints, w, levels, make([]int32, scratchSize(w, h)))
	} else {
		idwt97(tc.flts, w, levels, make([]float32, scratchSize(w, h)))
	}
}

// alloc allocates the coefficient buffer of a tile-component.
func (tc *tileComp) alloc() {
	n := (tc.x1 - tc.x0) * (tc.y1 - tc.y0)
	if tc.reversible {
		tc.ints = make([]int32, n)
	} else {
		tc.flts = make([]float32, n)
	}
}

// decodeBlock runs Tier-1 on a code-block and writes its dequantized
// coefficients into the tile-component buffer of width stride.
func (td *tileDecoder) decodeBlock(tc *tileComp, b *band, cb *codeBlock, stride int) {
	p0 := b.mb + tc.roi - 1 - cb.zbp
	if p0 < 0 || p0 > 30 {
		return // magnitudes at twice their scale must fit 32 bits
	}
	w, h := cb.x1-cb.x0, cb.y1-cb.y0
	orient := orientLL
	switch b.typ {
	case bandHL:
		orient = orientHL
	case bandHH:
		orient = orientHH
	}
	t := &td.t1
	t.decode(w, h, orient, tc.style, p0, cb.segs, cb.data)
	s := t.stride
	thresh := uint32(0)
	if tc.roi > 0 {
		thresh = 2 << tc.roi
	}
	ox := b.offX + cb.x0 - b.x0
	oy := b.offY + cb.y0 - b.y0
	half := b.delta / 2
	for y := range h {
		src := t.data[y*w : y*w+w]
		flags := t.flags[(y+1)*s+1 : (y+1)*s+1+w]
		o := (oy+y)*stride + ox
		for x, m := range src {
			if m == 0 {
				continue
			}
			if thresh != 0 && m >= thresh {
				m >>= tc.roi
			}
			neg := flags[x]&fNeg != 0
			if tc.reversible {
				v := int32(m >> 1)
				if neg {
					v = -v
				}
				tc.ints[o+x] = v
			} else {
				v := float32(m) * half
				if neg {
					v = -v
				}
				tc.flts[o+x] = v
			}
		}
	}
}

// inverseMCT undoes the component transform of the first three components
// (G.2, G.3) when they share size and filter.
func (td *tileDecoder) inverseMCT() {
	a, b, c := &td.comps[0], &td.comps[1], &td.comps[2]
	for _, o := range []*tileComp{b, c} {
		if o.x0 != a.x0 || o.x1 != a.x1 || o.y0 != a.y0 || o.y1 != a.y1 || o.reversible != a.reversible {
			return
		}
	}
	if a.ints == nil && b.ints == nil && c.ints == nil && a.flts == nil && b.flts == nil && c.flts == nil {
		return // zero stays zero
	}
	for _, o := range []*tileComp{a, b, c} {
		if o.ints == nil && o.flts == nil {
			o.alloc()
		}
	}
	if a.reversible {
		y0, y1, y2 := a.ints, b.ints[:len(a.ints)], c.ints[:len(a.ints)]
		for i := range y0 {
			g := y0[i] - (y1[i]+y2[i])>>2
			y0[i], y1[i], y2[i] = y2[i]+g, g, y1[i]+g
		}
		return
	}
	y0, y1, y2 := a.flts, b.flts[:len(a.flts)], c.flts[:len(a.flts)]
	for i := range y0 {
		// float32 conversions prevent fused multiply-adds (see dwt.go).
		y, cb, cr := y0[i], y1[i], y2[i]
		y0[i] = y + float32(1.402*cr)
		y1[i] = y - float32(0.34413*cb) - float32(0.71414*cr)
		y2[i] = y + float32(1.772*cb)
	}
}

// store applies the DC level shift, clamps, and copies the tile-component
// into its plane. A reversible tile-component that covers its plane
// becomes the plane's data.
func (td *tileDecoder) store(tc *tileComp, cp component, pl *plane) {
	w, h := tc.x1-tc.x0, tc.y1-tc.y0
	if w <= 0 || h <= 0 {
		return
	}
	var shift, lo, hi int32
	if cp.signed {
		lo, hi = -1<<(cp.prec-1), 1<<(cp.prec-1)-1
	} else {
		shift, hi = 1<<(cp.prec-1), 1<<cp.prec-1
	}
	whole := pl.data == nil && tc.x0 == pl.x0 && tc.y0 == pl.y0 && w == pl.w && h == pl.h
	if whole && tc.ints != nil {
		for i, v := range tc.ints {
			tc.ints[i] = min(max(v+shift, lo), hi)
		}
		pl.data, tc.ints = tc.ints, nil
		return
	}
	if pl.data == nil {
		pl.data = make([]int32, pl.w*pl.h)
	}
	ox, oy := tc.x0-pl.x0, tc.y0-pl.y0
	flo, fhi := float32(lo-shift), float32(hi-shift)
	for y := range h {
		dst := pl.data[(oy+y)*pl.w+ox : (oy+y)*pl.w+ox+w]
		switch {
		case tc.ints != nil:
			for x, v := range tc.ints[y*w : y*w+w] {
				dst[x] = min(max(v+shift, lo), hi)
			}
		case tc.flts != nil:
			for x, v := range tc.flts[y*w : y*w+w] {
				if !(v >= flo) { // also NaN
					v = flo
				} else if v > fhi {
					v = fhi
				}
				var n int32
				if v >= 0 {
					n = int32(v + 0.5)
				} else {
					n = -int32(0.5 - v)
				}
				dst[x] = n + shift
			}
		default:
			v := min(max(shift, lo), hi)
			for x := range dst {
				dst[x] = v
			}
		}
	}
	tc.ints, tc.flts = nil, nil
}

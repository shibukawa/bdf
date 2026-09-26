package jbig2

// Context bits of the generic refinement templates (6.3.5.3). Both take bit 0
// from pixel x−1 of the current row and further bits from the row above and
// from the 3×3 neighbourhood of the corresponding reference pixel.
var refinementBits = [2]int{13, 10}

// refinementParams are the parameters of the generic refinement region
// decoding procedure.
type refinementParams struct {
	template int
	at       []point // two AT pixels for template 0: one in the new bitmap, one in the reference
	ref      *Bitmap
	dx, dy   int
	tpgron   bool
}

func newRefinementStats(template int) []uint8 { return make([]uint8, 1<<refinementBits[template]) }

// decodeRefinement decodes a w×h bitmap with the generic refinement region
// decoding procedure (6.3.5). If the coded data runs out, the rows not yet
// decoded stay white.
func (d *decoder) decodeRefinement(ad *arithDecoder, stats []uint8, w, h int, p *refinementParams) (*Bitmap, error) {
	bm, err := d.newBitmap(w, h)
	if err != nil {
		return nil, err
	}
	ref, dx, dy := p.ref, clampCoord(int64(p.dx)), clampCoord(int64(p.dy))
	var at0, at1 point
	nominal := true
	if p.template == 0 {
		at0, at1 = p.at[0], p.at[1]
		nominal = at0 == point{-1, -1} && at1 == point{-1, -1}
	}
	// Typical prediction decodes this context, whose only 1 is the reference
	// pixel corresponding to the current one (6.3.5.6).
	sltp := 0x100
	if p.template == 1 {
		sltp = 0x080
	}
	ltp := 0
	for y := 0; y < h && !ad.exhausted(); y++ {
		if p.tpgron {
			ltp ^= ad.decode(&stats[sltp])
		}
		row, up := bm.row(y), bm.row(y-1)
		ry := y - dy
		rp, r0, rm := ref.row(ry+1), ref.row(ry), ref.row(ry-1)
		// Three-pixel windows with bit 0 at the right: the row above at x+1,
		// x, x−1, and the reference rows around x−dx.
		u := window{1, -1, 0}.start(up, w)
		rw := window{1 - dx, -1 - dx, 0}
		vp, v0, vm := rw.start(rp, ref.Width), rw.start(r0, ref.Width), rw.start(rm, ref.Width)
		c := 0
		for x := 0; x < w; x++ {
			b := -1
			if ltp != 0 {
				if vp|v0|vm == 0 {
					b = 0
				} else if vp&v0&vm == 7 {
					b = 1
				}
			}
			if b < 0 {
				var cx int
				if p.template == 0 {
					cx = c | (u&3)<<1 | (vp&7)<<4 | (v0&7)<<7 | (vm&3)<<10
					if nominal {
						cx |= (u>>2&1)<<3 | (vm>>2&1)<<12
					} else {
						cx |= bm.pixel(x+at0.x, y+at0.y)<<3 | ref.pixel(x-dx+at1.x, ry+at1.y)<<12
					}
				} else {
					cx = c | (u&7)<<1 | (vp&3)<<4 | (v0&7)<<6 | (vm>>1&1)<<9
				}
				b = ad.decode(&stats[cx])
			}
			if b != 0 {
				row[x>>3] |= 0x80 >> uint(x&7)
			}
			c = b
			u = (u<<1 | bit(up, w, x+2)) & 7
			rx := x - dx + 2
			vp = (vp<<1 | bit(rp, ref.Width, rx)) & 7
			v0 = (v0<<1 | bit(r0, ref.Width, rx)) & 7
			vm = (vm<<1 | bit(rm, ref.Width, rx)) & 7
		}
	}
	return bm, nil
}

// refinementRegionSegment decodes a generic refinement region segment (7.4.7).
// It refines the intermediate region it refers to, or else the part of the
// page it covers.
func (d *decoder) refinementRegionSegment(s *segment, data []byte) error {
	r := reader{b: data}
	ri, err := readRegionInfo(&r)
	if err != nil {
		return err
	}
	flags := r.u8()
	p := &refinementParams{template: int(flags & 1), tpgron: flags&2 != 0}
	if p.template == 0 {
		p.at = r.at(2)
	}
	body := r.rest()
	if r.err != nil {
		return r.err
	}
	w, h, err := ri.size()
	if err != nil {
		return err
	}
	for _, v := range d.referred(s) {
		if bm, ok := v.(*Bitmap); ok {
			p.ref = bm
		}
	}
	if p.ref == nil {
		if d.page == nil {
			return errorf("refinement region segment %d before the page information", s.number)
		}
		if p.ref, err = d.newBitmap(w, h); err != nil {
			return err
		}
		p.ref.compose(d.page.bm, -clampCoord(int64(ri.x)), -clampCoord(int64(ri.y)), opReplace)
	}
	bm, err := d.decodeRefinement(newArithDecoder(body), newRefinementStats(p.template), w, h, p)
	if err != nil {
		return err
	}
	return d.placeRegion(s, bm, ri)
}

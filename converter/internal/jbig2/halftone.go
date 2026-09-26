package jbig2

// patternDict is the result of a pattern dictionary segment.
type patternDict struct {
	patterns []*Bitmap
}

// patternDictSegment decodes a pattern dictionary segment (7.4.4, 6.7): a
// collective bitmap of GRAYMAX+1 patterns side by side.
func (d *decoder) patternDictSegment(s *segment, data []byte) error {
	r := reader{b: data}
	flags := r.u8()
	pw, ph := int(r.u8()), int(r.u8())
	grayMax := int64(r.u32())
	body := r.rest()
	if r.err != nil {
		return r.err
	}
	if pw == 0 || ph == 0 {
		return errorf("empty patterns")
	}
	n := grayMax + 1
	if n*int64(pw) > maxDimension {
		return errorf("too many patterns")
	}
	w := int(n) * pw
	var bm *Bitmap
	var err error
	if flags&1 != 0 {
		bm, _, err = d.decodeMMR(body, w, ph)
	} else {
		template := int(flags >> 1 & 3)
		at := []point{{-pw, 0}, {-3, -1}, {2, -2}, {-2, -2}}
		bm, err = d.decodeGeneric(newArithDecoder(body), newGenericStats(template), w, ph,
			&genericParams{template: template, at: at[:len(genericTemplates[template].nominal)]})
	}
	if err != nil {
		return err
	}
	pd := &patternDict{patterns: make([]*Bitmap, n)}
	for i := range pd.patterns {
		pd.patterns[i] = bm.sub(i*pw, 0, pw, ph)
	}
	s.result = pd
	return nil
}

// halftoneRegionSegment decodes a halftone region segment (7.4.5, 6.6).
func (d *decoder) halftoneRegionSegment(s *segment, data []byte) error {
	r := reader{b: data}
	ri, err := readRegionInfo(&r)
	if err != nil {
		return err
	}
	flags := r.u8()
	mmr, template, enableSkip := flags&1 != 0, int(flags>>1&3), flags&8 != 0
	op, defPixel := int(flags>>4&7), int(flags>>7)
	gw, gh := r.u32(), r.u32()
	gx, gy := int64(int32(r.u32())), int64(int32(r.u32()))
	rx, ry := int64(r.u16()), int64(r.u16())
	body := r.rest()
	if r.err != nil {
		return r.err
	}
	if op > opReplace {
		return errorf("invalid combination operator %d", op)
	}
	w, h, err := ri.size()
	if err != nil {
		return err
	}
	var pats []*Bitmap
	for _, v := range d.referred(s) {
		if pd, ok := v.(*patternDict); ok {
			pats = pd.patterns
		}
	}
	if len(pats) == 0 {
		return errorf("halftone region %d without patterns", s.number)
	}
	gridW, err := dim(gw)
	if err != nil {
		return err
	}
	gridH, err := dim(gh)
	if err != nil {
		return err
	}
	// Grid cells are normally as large as the patterns; more cells than twice
	// the region's pixels only draw over each other.
	if cells := int64(gridW) * int64(gridH); cells > 2*int64(w)*int64(h)+1<<16 {
		return errorf("halftone grid too large")
	}
	bm, err := d.newBitmap(w, h)
	if err != nil {
		return err
	}
	bm.fill(defPixel)
	pw, ph := int64(pats[0].Width), int64(pats[0].Height)
	// Grid positions (6.6.5.2) are in 1/256 pixel.
	pos := func(mg, ng int) (x, y int64) {
		return (gx + int64(mg)*ry + int64(ng)*rx) >> 8, (gy + int64(mg)*rx - int64(ng)*ry) >> 8
	}
	var skip *Bitmap
	if enableSkip && !mmr {
		if skip, err = d.newBitmap(gridW, gridH); err != nil {
			return err
		}
		for mg := 0; mg < gridH; mg++ {
			for ng := 0; ng < gridW; ng++ {
				x, y := pos(mg, ng)
				if x+pw <= 0 || x >= int64(w) || y+ph <= 0 || y >= int64(h) {
					skip.set(ng, mg)
				}
			}
		}
	}
	planes, err := d.decodeGrayScale(body, gridW, gridH, ceilLog2(len(pats)), mmr, template, skip)
	if err != nil {
		return err
	}
	for mg := 0; mg < gridH; mg++ {
		if err := d.spend(int64(gridW) * int64(len(planes)+1)); err != nil {
			return err
		}
		for ng := 0; ng < gridW; ng++ {
			x, y := pos(mg, ng)
			if x+pw <= 0 || x >= int64(w) || y+ph <= 0 || y >= int64(h) {
				continue // off the region
			}
			v := 0
			for j, p := range planes {
				v |= p.pixel(ng, mg) << uint(j)
			}
			if err := d.draw(bm, pats[min(v, len(pats)-1)], x, y, op); err != nil {
				return err
			}
		}
	}
	return d.placeRegion(s, bm, ri)
}

// decodeGrayScale decodes a gray-scale image as bits-per-pixel bitplanes,
// most significant first, each Gray coded against the one above (Annex C.5).
// It returns the planes with plane j holding bit j of the values.
func (d *decoder) decodeGrayScale(data []byte, w, h, bpp int, mmr bool, template int, skip *Bitmap) ([]*Bitmap, error) {
	planes := make([]*Bitmap, bpp)
	var ad *arithDecoder
	var stats []uint8
	params := &genericParams{template: template, skip: skip}
	if !mmr {
		ad, stats = newArithDecoder(data), newGenericStats(template)
		at := []point{{3, -1}, {-3, -1}, {2, -2}, {-2, -2}}
		if template >= 2 {
			at[0].x = 2
		}
		params.at = at[:len(genericTemplates[template].nominal)]
	}
	for j := bpp - 1; j >= 0; j-- {
		var err error
		if mmr {
			var n int
			planes[j], n, err = d.decodeMMR(data, w, h)
			data = data[n:]
		} else {
			planes[j], err = d.decodeGeneric(ad, stats, w, h, params)
		}
		if err != nil {
			return nil, err
		}
		if j+1 < bpp {
			p, q := planes[j].Data, planes[j+1].Data
			for i := range p {
				p[i] ^= q[i]
			}
		}
	}
	return planes, nil
}

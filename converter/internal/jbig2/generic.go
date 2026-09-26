package jbig2

// window describes the pixels a row above the current one contributes to a
// context: pixels x+hi down to x+lo land at context bits shift upwards. A
// window with lo > hi is empty.
type window struct{ hi, lo, shift int }

func (w window) mask() int { return 1<<(w.hi-w.lo+1) - 1 }

// start returns the window's value at x = 0.
func (w window) start(row []byte, width int) int {
	v := 0
	for k := w.lo; k <= w.hi; k++ {
		v = v<<1 | bit(row, width, k)
	}
	return v
}

// genericTemplate is a generic region template (6.2.5.3) in the context bit
// order of the standard: the current row supplies bits 0 to cur−1 (x−1, x−2,
// ...), rows y−1 and y−2 the windows rows, and the adaptive template (AT)
// pixels the bits atBits. When the AT pixels are at their nominal positions,
// the wider windows fast take them in instead.
type genericTemplate struct {
	bits    int
	cur     int
	rows    [2]window
	atBits  []int
	nominal []point
	fast    [2]window
	sltp    int // the context of the typical prediction bit (6.2.5.7)
}

var genericTemplates = [4]genericTemplate{
	{16, 4, [2]window{{2, -2, 5}, {1, -1, 12}}, []int{4, 10, 11, 15},
		[]point{{3, -1}, {-3, -1}, {2, -2}, {-2, -2}}, [2]window{{3, -3, 4}, {2, -2, 11}}, 0x9b25},
	{13, 3, [2]window{{2, -2, 4}, {2, -1, 9}}, []int{3},
		[]point{{3, -1}}, [2]window{{3, -2, 3}, {2, -1, 9}}, 0x0795},
	{10, 2, [2]window{{1, -2, 3}, {1, -1, 7}}, []int{2},
		[]point{{2, -1}}, [2]window{{2, -2, 2}, {1, -1, 7}}, 0x00e5},
	{10, 4, [2]window{{1, -3, 5}, {0, 1, 0}}, []int{4},
		[]point{{2, -1}}, [2]window{{2, -3, 4}, {0, 1, 0}}, 0x0195},
}

// genericParams are the parameters of the generic region decoding procedure.
type genericParams struct {
	template int
	tpgdon   bool
	at       []point // as many AT pixels as the template has
	skip     *Bitmap // pixels not to decode (USESKIP), or nil
}

// newGenericStats returns fresh contexts for a generic region template.
func newGenericStats(template int) []uint8 {
	return make([]uint8, 1<<genericTemplates[template].bits)
}

// decodeGeneric decodes a w×h bitmap with the arithmetic generic region
// decoding procedure (6.2.5). If the coded data runs out, the rows not yet
// decoded stay white.
func (d *decoder) decodeGeneric(ad *arithDecoder, stats []uint8, w, h int, p *genericParams) (*Bitmap, error) {
	bm, err := d.newBitmap(w, h)
	if err != nil {
		return nil, err
	}
	t := &genericTemplates[p.template]
	rows, atBits, at := t.rows, t.atBits, p.at
	nominal := len(at) == len(t.nominal)
	for i := range t.nominal {
		nominal = nominal && at[i] == t.nominal[i]
	}
	if nominal {
		rows, atBits, at = t.fast, nil, nil
	}
	m0, m1 := rows[0].mask(), rows[1].mask()
	s0, s1 := uint(rows[0].shift), uint(rows[1].shift)
	h0, h1 := rows[0].hi+1, rows[1].hi+1
	cm := 1<<t.cur - 1
	ltp := 0
	for y := 0; y < h && !ad.exhausted(); y++ {
		row := bm.row(y)
		if p.tpgdon {
			ltp ^= ad.decode(&stats[t.sltp])
			if ltp != 0 {
				if y > 0 {
					copy(row, bm.row(y-1))
				}
				continue
			}
		}
		r1, r2 := bm.row(y-1), bm.row(y-2)
		w1, w2 := rows[0].start(r1, w), rows[1].start(r2, w)
		c := 0
		for x := 0; x < w; x++ {
			b := 0
			if p.skip == nil || p.skip.pixel(x, y) == 0 {
				cx := c | w1<<s0 | w2<<s1
				for i, a := range at {
					cx |= bm.pixel(x+a.x, y+a.y) << uint(atBits[i])
				}
				if b = ad.decode(&stats[cx]); b != 0 {
					row[x>>3] |= 0x80 >> uint(x&7)
				}
			}
			c = (c<<1 | b) & cm
			w1 = (w1<<1 | bit(r1, w, x+h0)) & m0
			w2 = (w2<<1 | bit(r2, w, x+h1)) & m1
		}
	}
	return bm, nil
}

// genericRegionSegment decodes a generic region segment (7.4.6).
func (d *decoder) genericRegionSegment(s *segment, data []byte) error {
	r := reader{b: data}
	ri, err := readRegionInfo(&r)
	if err != nil {
		return err
	}
	flags := r.u8()
	mmr, template, tpgdon := flags&1 != 0, int(flags>>1&3), flags&8 != 0
	if flags&0x10 != 0 {
		return errorf("extended generic region templates are not supported")
	}
	var at []point
	if !mmr {
		at = r.at(len(genericTemplates[template].nominal))
	}
	body := r.rest()
	if r.err != nil {
		return r.err
	}
	if s.length == unknownLength && len(data) >= 4 {
		// The data ends with the row count (7.2.7).
		rows := uint32(data[len(data)-4])<<24 | uint32(data[len(data)-3])<<16 | uint32(data[len(data)-2])<<8 | uint32(data[len(data)-1])
		if rows > ri.h {
			return errorf("generic region row count %d exceeds its height", rows)
		}
		ri.h = rows
	}
	w, h, err := ri.size()
	if err != nil {
		return err
	}
	var bm *Bitmap
	if mmr {
		bm, _, err = d.decodeMMR(body, w, h)
	} else {
		bm, err = d.decodeGeneric(newArithDecoder(body), newGenericStats(template), w, h,
			&genericParams{template: template, tpgdon: tpgdon, at: at})
	}
	if err != nil {
		return err
	}
	return d.placeRegion(s, bm, ri)
}

package jbig2

// Reference corners of symbol instances (7.4.3.1.1).
const (
	cornerBottomLeft = iota
	cornerTopLeft
	cornerBottomRight
	cornerTopRight
)

// textContexts are the arithmetic integer contexts of the text region
// decoding procedure, which symbol dictionaries share for aggregation.
type textContexts struct {
	iadt, iafs, iads, iait, iari, iardw, iardh, iardx, iardy intContext
	iaid                                                     []uint8
}

func newTextContexts(codeLen int) *textContexts {
	return &textContexts{iaid: make([]uint8, 1<<codeLen)}
}

// textParams are the parameters of the text region decoding procedure (6.4.2).
type textParams struct {
	huff, refine bool
	w, h         int
	instances    int64
	logStrips    int
	syms         []*Bitmap
	codeLen      int        // symbol ID code length
	codes        *huffTable // Huffman symbol ID codes; nil for fixed-length codes
	defPixel     int
	op           int
	transposed   bool
	corner       int
	dsOffset     int
	// Huffman tables
	fs, ds, dt, rdw, rdh, rdx, rdy, rsize *huffTable
	// refinement
	rtemplate int
	rat       []point
}

// textDecoder reads the values of a text region from either coder.
type textDecoder struct {
	ad *arithDecoder
	br *bitReader
	cx *textContexts
	gr []uint8 // refinement contexts
}

func (t *textDecoder) int(huff *huffTable, cx *intContext) (int, bool, error) {
	if t.br != nil {
		return huff.decode(t.br)
	}
	v, ok := t.ad.decodeInt(cx)
	return v, ok, nil
}

// value reads an integer that may not be out of band.
func (t *textDecoder) value(huff *huffTable, cx *intContext) (int, error) {
	v, ok, err := t.int(huff, cx)
	if err == nil && !ok {
		err = errorf("unexpected out-of-band value")
	}
	return v, err
}

// decodeText decodes a text region (6.4.5). If the arithmetic coded data runs
// out, the instances decoded so far make up the region.
func (d *decoder) decodeText(p *textParams, t *textDecoder) (*Bitmap, error) {
	bm, err := d.newBitmap(p.w, p.h)
	if err != nil {
		return nil, err
	}
	if p.defPixel != 0 {
		bm.fill(1)
	}
	strips := int64(1) << p.logStrips
	dt, err := t.value(p.dt, &t.cx.iadt)
	if err != nil {
		return nil, err
	}
	stript := -int64(dt) * strips
	firsts := int64(0)
	for n := int64(0); n < p.instances; {
		if t.ad != nil && t.ad.exhausted() {
			return bm, nil // the coded data ran out: keep what was drawn
		}
		dt, err := t.value(p.dt, &t.cx.iadt)
		if err != nil {
			return nil, err
		}
		stript += int64(dt) * strips
		var curs int64
		for first := true; ; first = false {
			if first {
				dfs, err := t.value(p.fs, &t.cx.iafs)
				if err != nil {
					return nil, err
				}
				firsts += int64(dfs)
				curs = firsts
			} else {
				ids, ok, err := t.int(p.ds, &t.cx.iads)
				if err != nil {
					return nil, err
				}
				if !ok {
					break // the end of the strip
				}
				if n >= p.instances {
					return bm, nil // more instances than announced
				}
				curs += int64(ids) + int64(p.dsOffset)
			}
			curt := int64(0)
			if strips > 1 {
				if t.br != nil {
					v, err := t.br.readBits(p.logStrips)
					if err != nil {
						return nil, err
					}
					curt = int64(v)
				} else {
					v, err := t.value(nil, &t.cx.iait)
					if err != nil {
						return nil, err
					}
					curt = int64(v)
				}
			}
			ti := stript + curt
			if t.ad != nil && t.ad.exhausted() {
				return bm, nil
			}
			ib, err := d.textSymbol(p, t)
			if err != nil {
				return nil, err
			}
			w, h := int64(ib.Width), int64(ib.Height)
			if !p.transposed && (p.corner == cornerTopRight || p.corner == cornerBottomRight) {
				curs += w - 1
			} else if p.transposed && (p.corner == cornerBottomLeft || p.corner == cornerBottomRight) {
				curs += h - 1
			}
			x, y := curs, ti
			if p.transposed {
				x, y = ti, curs
			}
			if p.corner == cornerTopRight || p.corner == cornerBottomRight {
				x -= w - 1
			}
			if p.corner == cornerBottomLeft || p.corner == cornerBottomRight {
				y -= h - 1
			}
			if err := d.draw(bm, ib, x, y, p.op); err != nil {
				return nil, err
			}
			if !p.transposed && (p.corner == cornerTopLeft || p.corner == cornerBottomLeft) {
				curs += w - 1
			} else if p.transposed && (p.corner == cornerTopLeft || p.corner == cornerTopRight) {
				curs += h - 1
			}
			curs = max(min(curs, 1<<40), -1<<40)
			n++
		}
		firsts = max(min(firsts, 1<<40), -1<<40)
		stript = max(min(stript, 1<<40), -1<<40)
	}
	return bm, nil
}

// textSymbol decodes the symbol ID of an instance and, if it is refined, the
// refinement, returning the instance's bitmap (6.4.10, 6.4.11).
func (d *decoder) textSymbol(p *textParams, t *textDecoder) (*Bitmap, error) {
	var id int
	switch {
	case t.br == nil:
		id = t.ad.decodeIAID(t.cx.iaid, p.codeLen)
	case p.codes != nil:
		v, err := t.value(p.codes, nil)
		if err != nil {
			return nil, err
		}
		id = v
	default:
		v, err := t.br.readBits(p.codeLen)
		if err != nil {
			return nil, err
		}
		id = int(v)
	}
	if id < 0 || id >= len(p.syms) || p.syms[id] == nil {
		return nil, errorf("invalid symbol ID %d", id)
	}
	sym := p.syms[id]
	if !p.refine {
		return sym, nil
	}
	var ri int
	if t.br != nil {
		b, err := t.br.readBit()
		if err != nil {
			return nil, err
		}
		ri = b
	} else {
		v, err := t.value(nil, &t.cx.iari)
		if err != nil {
			return nil, err
		}
		ri = v
	}
	if ri == 0 {
		return sym, nil
	}
	if err := d.spend(64); err != nil {
		return nil, err
	}
	var v [4]int
	for i, c := range []struct {
		huff *huffTable
		cx   *intContext
	}{{p.rdw, &t.cx.iardw}, {p.rdh, &t.cx.iardh}, {p.rdx, &t.cx.iardx}, {p.rdy, &t.cx.iardy}} {
		var err error
		if v[i], err = t.value(c.huff, c.cx); err != nil {
			return nil, err
		}
	}
	rdw, rdh, rdx, rdy := int64(v[0]), int64(v[1]), int64(v[2]), int64(v[3])
	ad := t.ad
	if t.br != nil {
		size, err := t.value(p.rsize, nil)
		if err != nil {
			return nil, err
		}
		t.br.align()
		data, err := t.br.bytes(size)
		if err != nil {
			return nil, err
		}
		ad = newArithDecoder(data)
	}
	w, h := int64(sym.Width)+rdw, int64(sym.Height)+rdh
	if w < 0 || h < 0 || w > maxDimension || h > maxDimension {
		return nil, errorf("refined symbol size out of range")
	}
	return d.decodeRefinement(ad, t.gr, int(w), int(h), &refinementParams{
		template: p.rtemplate, at: p.rat, ref: sym,
		dx: clampCoord(rdw>>1 + rdx), dy: clampCoord(rdh>>1 + rdy),
	})
}

// textRegionSegment decodes a text region segment (7.4.3).
func (d *decoder) textRegionSegment(s *segment, data []byte) error {
	r := reader{b: data}
	ri, err := readRegionInfo(&r)
	if err != nil {
		return err
	}
	flags := int(r.u16())
	p := &textParams{
		huff:       flags&1 != 0,
		refine:     flags&2 != 0,
		logStrips:  flags >> 2 & 3,
		corner:     flags >> 4 & 3,
		transposed: flags&0x40 != 0,
		op:         flags >> 7 & 3,
		defPixel:   flags >> 9 & 1,
		dsOffset:   flags >> 10 & 0x1f,
		rtemplate:  flags >> 15 & 1,
	}
	if p.dsOffset >= 16 {
		p.dsOffset -= 32
	}
	var huffFlags int
	if p.huff {
		huffFlags = int(r.u16())
	}
	if p.refine && p.rtemplate == 0 {
		p.rat = r.at(2)
	}
	p.instances = int64(r.u32())
	if r.err != nil {
		return r.err
	}
	if p.w, p.h, err = ri.size(); err != nil {
		return err
	}
	sel := tableSelector{}
	for _, v := range d.referred(s) {
		switch v := v.(type) {
		case *symbolDict:
			p.syms = append(p.syms, v.exported...)
			if int64(len(p.syms)) > maxSymbols {
				return errorf("too many symbols")
			}
		case *huffTable:
			sel.custom = append(sel.custom, v)
		}
	}
	if p.instances == 0 {
		bm, err := d.newBitmap(p.w, p.h)
		if err != nil {
			return err
		}
		bm.fill(p.defPixel)
		return d.placeRegion(s, bm, ri)
	}
	p.codeLen = ceilLog2(len(p.syms))
	t := &textDecoder{cx: newTextContexts(p.codeLen)}
	if p.refine {
		t.gr = newRefinementStats(p.rtemplate)
	}
	if p.huff {
		p.fs = sel.pick(huffFlags&3, 3, 6, 7)
		p.ds = sel.pick(huffFlags>>2&3, 3, 8, 9, 10)
		p.dt = sel.pick(huffFlags>>4&3, 3, 11, 12, 13)
		p.rdw = sel.pick(huffFlags>>6&3, 3, 14, 15)
		p.rdh = sel.pick(huffFlags>>8&3, 3, 14, 15)
		p.rdx = sel.pick(huffFlags>>10&3, 3, 14, 15)
		p.rdy = sel.pick(huffFlags>>12&3, 3, 14, 15)
		p.rsize = sel.pick(huffFlags>>14&1, 1, 1)
		if sel.err != nil {
			return sel.err
		}
		t.br = &bitReader{data: r.rest()}
		if p.codes, err = readSymbolIDTable(t.br, len(p.syms)); err != nil {
			return err
		}
	} else {
		t.ad = newArithDecoder(r.rest())
	}
	bm, err := d.decodeText(p, t)
	if err != nil {
		return err
	}
	return d.placeRegion(s, bm, ri)
}

// readSymbolIDTable reads the Huffman table of symbol IDs of a text region
// (7.4.3.1.7): 35 run code lengths, then the code lengths of the n symbols
// coded with the run codes.
func readSymbolIDTable(br *bitReader, n int) (*huffTable, error) {
	runLines := make([]huffLine, 35)
	for i := range runLines {
		v, err := br.readBits(4)
		if err != nil {
			return nil, err
		}
		runLines[i] = huffLine{prefLen: int(v), low: int64(i)}
	}
	runCodes, err := newHuffTable(runLines)
	if err != nil {
		return nil, err
	}
	lines := make([]huffLine, n)
	for i := 0; i < n; {
		c, _, err := runCodes.decode(br)
		if err != nil {
			return nil, err
		}
		length, repeat := c, 1
		var extra uint32
		switch c {
		case 32:
			if i == 0 {
				return nil, errorf("invalid symbol ID table")
			}
			extra, err = br.readBits(2)
			length, repeat = lines[i-1].prefLen, 3+int(extra)
		case 33:
			extra, err = br.readBits(3)
			length, repeat = 0, 3+int(extra)
		case 34:
			extra, err = br.readBits(7)
			length, repeat = 0, 11+int(extra)
		}
		if err != nil {
			return nil, err
		}
		if repeat > n-i {
			return nil, errorf("invalid symbol ID table")
		}
		for ; repeat > 0; repeat-- {
			lines[i] = huffLine{prefLen: length, low: int64(i)}
			i++
		}
	}
	br.align()
	return newHuffTable(lines)
}

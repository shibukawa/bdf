package jbig2

// symbolDict is the result of a symbol dictionary segment.
type symbolDict struct {
	exported []*Bitmap
	// Contexts retained for a later dictionary (7.4.2.1.1), or nil.
	gb, gr []uint8
}

// symbolDictSegment decodes a symbol dictionary segment (7.4.2, 6.5).
func (d *decoder) symbolDictSegment(s *segment, data []byte) error {
	r := reader{b: data}
	flags := int(r.u16())
	huff, refagg := flags&1 != 0, flags&2 != 0
	template, rtemplate := flags>>10&3, flags>>12&1
	var at, rat []point
	if !huff {
		at = r.at(len(genericTemplates[template].nominal))
	}
	if refagg && rtemplate == 0 {
		rat = r.at(2)
	}
	numEx, numNew := int64(r.u32()), int64(r.u32())
	if r.err != nil {
		return r.err
	}

	var in []*Bitmap
	var last *symbolDict
	sel := tableSelector{}
	for _, v := range d.referred(s) {
		switch v := v.(type) {
		case *symbolDict:
			in = append(in, v.exported...)
			last = v
		case *huffTable:
			sel.custom = append(sel.custom, v)
		}
	}
	total := int64(len(in)) + numNew
	if total > maxSymbols || numEx > total {
		return errorf("too many symbols")
	}
	if err := d.spend(64 * total); err != nil {
		return err
	}

	var tDH, tDW, tBM, tAI *huffTable
	if huff {
		tDH = sel.pick(flags>>2&3, 3, 4, 5)
		tDW = sel.pick(flags>>4&3, 3, 2, 3)
		tBM = sel.pick(flags>>6&1, 1, 1)
		tAI = sel.pick(flags>>7&1, 1, 1)
		if sel.err != nil {
			return sel.err
		}
	}

	codeLen := ceilLog2(int(total))
	t := &textDecoder{cx: newTextContexts(codeLen)}
	var iadh, iadw, iaex, iaai intContext
	var gb []uint8
	if huff {
		t.br = &bitReader{data: r.rest()}
	} else {
		t.ad = newArithDecoder(r.rest())
		gb = newGenericStats(template)
	}
	if refagg {
		t.gr = newRefinementStats(rtemplate)
	}
	if flags&0x100 != 0 { // bitmap coding context used
		if last == nil || (gb != nil && len(last.gb) != len(gb)) || (t.gr != nil && len(last.gr) != len(t.gr)) {
			return errorf("symbol dictionary %d: no matching retained contexts", s.number)
		}
		copy(gb, last.gb)
		copy(t.gr, last.gr)
	}

	syms := make([]*Bitmap, len(in), total)
	copy(syms, in)
	genericParams := &genericParams{template: template, at: at}
	height := int64(0)
	for empty := 0; int64(len(syms)) < total; {
		dh, err := t.value(tDH, &iadh)
		if err != nil {
			return err
		}
		height += int64(dh)
		if height < 0 || height > maxDimension {
			return errorf("symbol height out of range")
		}
		first := len(syms)
		width, totWidth := int64(0), int64(0)
		for {
			dw, ok, err := t.int(tDW, &iadw)
			if err != nil {
				return err
			}
			if !ok {
				break // the end of the height class
			}
			if int64(len(syms)) >= total {
				return errorf("too many symbols")
			}
			if t.ad != nil && t.ad.exhausted() {
				return errorf("symbol dictionary %d: coded data ends early", s.number)
			}
			width += int64(dw)
			totWidth += width
			if width < 0 || width > maxDimension || totWidth > maxDimension {
				return errorf("symbol width out of range")
			}
			var bm *Bitmap
			switch {
			case huff && !refagg:
				// Collected into the height class bitmap below; keep the width.
				bm = &Bitmap{Width: int(width)}
			case !refagg:
				bm, err = d.decodeGeneric(t.ad, gb, int(width), int(height), genericParams)
			default:
				bm, err = d.aggregateSymbol(t, syms, int(width), int(height), codeLen, rtemplate, rat, tAI, &iaai)
			}
			if err != nil {
				return err
			}
			syms = append(syms, bm)
		}
		if huff && !refagg {
			if err := d.heightClassBitmap(t, tBM, syms[first:], int(totWidth), int(height)); err != nil {
				return err
			}
		}
		if len(syms) == first {
			if empty++; empty > 1<<16 || (t.ad != nil && t.ad.exhausted()) {
				return errorf("symbol dictionary %d: no symbols decoded", s.number)
			}
		}
	}

	// Export flags (6.5.10).
	sd := &symbolDict{}
	export := false
	for i, runs := int64(0), 0; i < total; runs++ {
		var n int
		var err error
		if huff {
			n, err = t.value(standardTables[1], nil)
		} else {
			n, err = t.value(nil, &iaex)
		}
		if err != nil {
			return err
		}
		if n < 0 || int64(n) > total-i || (n == 0 && runs > int(total)+2) || (t.ad != nil && t.ad.exhausted()) {
			return errorf("invalid symbol export flags")
		}
		if export {
			if int64(len(sd.exported)+n) > numEx {
				return errorf("too many exported symbols")
			}
			sd.exported = append(sd.exported, syms[i:i+int64(n)]...)
		}
		i += int64(n)
		export = !export
	}
	if flags&0x200 != 0 { // bitmap coding context retained
		sd.gb, sd.gr = gb, t.gr
	}
	s.result = sd
	return nil
}

// aggregateSymbol decodes a refinement/aggregate coded symbol (6.5.8.2).
func (d *decoder) aggregateSymbol(t *textDecoder, syms []*Bitmap, w, h, codeLen, rtemplate int, rat []point, tAI *huffTable, iaai *intContext) (*Bitmap, error) {
	n, err := t.value(tAI, iaai)
	if err != nil {
		return nil, err
	}
	if n < 1 {
		return nil, errorf("invalid aggregate instance count %d", n)
	}
	if n > 1 {
		// Several symbols, drawn by the text region decoding procedure.
		p := &textParams{
			huff: t.br != nil, refine: true, w: w, h: h, instances: int64(n),
			syms: syms, codeLen: codeLen, corner: cornerTopLeft,
			fs: standardTables[6], ds: standardTables[8], dt: standardTables[11],
			rdw: standardTables[15], rdh: standardTables[15], rdx: standardTables[15], rdy: standardTables[15],
			rsize: standardTables[1], rtemplate: rtemplate, rat: rat,
		}
		return d.decodeText(p, t)
	}
	// One symbol, refined.
	var id, rdx, rdy int
	ad := t.ad
	if t.br != nil {
		v, err := t.br.readBits(codeLen)
		if err != nil {
			return nil, err
		}
		id = int(v)
		if rdx, err = t.value(standardTables[15], nil); err != nil {
			return nil, err
		}
		if rdy, err = t.value(standardTables[15], nil); err != nil {
			return nil, err
		}
		size, err := t.value(standardTables[1], nil)
		if err != nil {
			return nil, err
		}
		t.br.align()
		data, err := t.br.bytes(size)
		if err != nil {
			return nil, err
		}
		ad = newArithDecoder(data)
	} else {
		id = t.ad.decodeIAID(t.cx.iaid, codeLen)
		if rdx, err = t.value(nil, &t.cx.iardx); err != nil {
			return nil, err
		}
		if rdy, err = t.value(nil, &t.cx.iardy); err != nil {
			return nil, err
		}
	}
	if id < 0 || id >= len(syms) || syms[id] == nil {
		return nil, errorf("invalid symbol ID %d", id)
	}
	return d.decodeRefinement(ad, t.gr, w, h, &refinementParams{
		template: rtemplate, at: rat, ref: syms[id], dx: rdx, dy: rdy,
	})
}

// heightClassBitmap decodes the collective bitmap of a Huffman coded height
// class, uncompressed or MMR coded, and cuts it into the symbols (6.5.9).
func (d *decoder) heightClassBitmap(t *textDecoder, tBM *huffTable, syms []*Bitmap, w, h int) error {
	size, err := t.value(tBM, nil)
	if err != nil {
		return err
	}
	t.br.align()
	var bm *Bitmap
	if size == 0 {
		if bm, err = d.newBitmap(w, h); err != nil {
			return err
		}
		data, err := t.br.bytes(len(bm.Data))
		if err != nil {
			return err
		}
		copy(bm.Data, data)
		bm.clearPadding()
	} else {
		data, err := t.br.bytes(size)
		if err != nil {
			return err
		}
		if bm, _, err = d.decodeMMR(data, w, h); err != nil {
			return err
		}
	}
	x := 0
	for i, s := range syms {
		syms[i] = bm.sub(x, 0, s.Width, h)
		x += s.Width
	}
	return nil
}

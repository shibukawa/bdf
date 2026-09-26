package jbig2

// Segment level encoders for tests, and a reference renderer that computes the
// page a stream should decode to.

import "encoding/binary"

// stream collects segments.
type stream struct {
	buf  []byte
	next uint32 // number of the next segment
}

// segOpts are the less common segment header options.
type segOpts struct {
	refs    []uint32
	page    uint32 // page association; 1 when zero unless global
	global  bool   // page association 0
	page4   bool   // 4-byte page association
	unknown bool   // unknown data length (immediate generic regions)
	number  uint32 // segment number, if not the next one
}

// add appends a segment and returns its number.
func (s *stream) add(typ int, data []byte, o segOpts) uint32 {
	num := s.next
	if o.number != 0 {
		num = o.number
	}
	s.next = num + 1
	b := binary.BigEndian.AppendUint32(nil, num)
	flags := byte(typ)
	page := o.page
	if page == 0 && !o.global {
		page = 1
	}
	if o.page4 || page > 255 {
		flags |= 0x40
	}
	b = append(b, flags)
	n := len(o.refs)
	if n <= 4 {
		b = append(b, byte(n<<5))
	} else {
		b = binary.BigEndian.AppendUint32(b, 0xe0000000|uint32(n))
		b = append(b, make([]byte, (n+8)/8)...)
	}
	for _, r := range o.refs {
		switch {
		case num > 65536:
			b = binary.BigEndian.AppendUint32(b, r)
		case num > 256:
			b = binary.BigEndian.AppendUint16(b, uint16(r))
		default:
			b = append(b, byte(r))
		}
	}
	if flags&0x40 != 0 {
		b = binary.BigEndian.AppendUint32(b, page)
	} else {
		b = append(b, byte(page))
	}
	if o.unknown {
		b = binary.BigEndian.AppendUint32(b, unknownLength)
	} else {
		b = binary.BigEndian.AppendUint32(b, uint32(len(data)))
	}
	s.buf = append(append(s.buf, b...), data...)
	return num
}

func pageInfo(w, h uint32, flags byte, striping uint16) []byte {
	b := binary.BigEndian.AppendUint32(nil, w)
	b = binary.BigEndian.AppendUint32(b, h)
	b = append(b, make([]byte, 8)...)
	b = append(b, flags)
	return binary.BigEndian.AppendUint16(b, striping)
}

func regionInfoBytes(w, h, x, y, op int) []byte {
	b := binary.BigEndian.AppendUint32(nil, uint32(w))
	b = binary.BigEndian.AppendUint32(b, uint32(h))
	b = binary.BigEndian.AppendUint32(b, uint32(x))
	b = binary.BigEndian.AppendUint32(b, uint32(y))
	return append(b, byte(op))
}

func atBytes(at []point) []byte {
	var b []byte
	for _, p := range at {
		b = append(b, byte(int8(p.x)), byte(int8(p.y)))
	}
	return b
}

// genericRegion codes a generic region segment's data.
func genericRegion(bm *Bitmap, x, y, op int, mmr bool, template int, at []point, tpgdon bool) []byte {
	b := regionInfoBytes(bm.Width, bm.Height, x, y, op)
	if mmr {
		return append(append(b, 1), encodeMMR(bm, true)...)
	}
	flags := byte(template << 1)
	if tpgdon {
		flags |= 8
	}
	b = append(append(b, flags), atBytes(at)...)
	e := newMQEncoder()
	encodeGeneric(e, newGenericStats(template), bm, template, at, tpgdon, nil)
	return append(b, e.flush()...)
}

// refinementRegion codes a refinement region segment's data.
func refinementRegion(bm, ref *Bitmap, x, y, op, template int, at []point, tpgron bool) []byte {
	b := regionInfoBytes(bm.Width, bm.Height, x, y, op)
	flags := byte(template)
	if tpgron {
		flags |= 2
	}
	b = append(b, flags)
	if template == 0 {
		b = append(b, atBytes(at)...)
	}
	e := newMQEncoder()
	encodeRefinement(e, newRefinementStats(template), bm, ref, 0, 0, template, at, tpgron)
	return append(b, e.flush()...)
}

// tableSegment codes a custom Huffman table (B.2). Lines are listed as
// {prefix length, range length}; the lower, upper and OOB prefix lengths follow.
func tableSegment(low, high int32, lines [][2]int, lowerLen, upperLen, oobLen int) ([]byte, *huffTable) {
	ps, rs := 1, 1
	for _, l := range append(lines, [2]int{lowerLen, 0}, [2]int{upperLen, 0}, [2]int{oobLen, 0}) {
		for 1<<ps <= l[0] {
			ps++
		}
		for 1<<rs <= l[1] {
			rs++
		}
	}
	flags := byte((ps-1)<<1 | (rs-1)<<4)
	if oobLen > 0 {
		flags |= 1
	}
	b := append([]byte{flags}, binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32(nil, uint32(low)), uint32(high))...)
	var w bitWriter
	var hl []huffLine
	cur := int64(low)
	for _, l := range lines {
		w.write(uint64(l[0]), ps)
		w.write(uint64(l[1]), rs)
		hl = append(hl, huffLine{prefLen: l[0], rangeLen: l[1], low: cur})
		cur += 1 << uint(l[1])
	}
	if cur != int64(high) {
		panic("table lines do not end at high")
	}
	w.write(uint64(lowerLen), ps)
	w.write(uint64(upperLen), ps)
	hl = append(hl, huffLine{prefLen: lowerLen, rangeLen: 32, low: int64(low) - 1, kind: lineLower},
		huffLine{prefLen: upperLen, rangeLen: 32, low: int64(high), kind: lineUpper})
	if oobLen > 0 {
		w.write(uint64(oobLen), ps)
		hl = append(hl, huffLine{prefLen: oobLen, kind: lineOOB})
	}
	t, err := newHuffTable(hl)
	if err != nil {
		panic(err)
	}
	return append(b, w.out...), t
}

// coder writes values with the arithmetic coder or Huffman tables.
type coder struct {
	e  *mqEncoder
	w  *bitWriter
	cx *textContexts
	gr []uint8
}

func (c *coder) int(t *huffTable, cx *intContext, v int) {
	if c.w != nil {
		encodeHuff(c.w, t, v, false)
	} else {
		c.e.encodeInt(cx, v, false)
	}
}

func (c *coder) oob(t *huffTable, cx *intContext) {
	if c.w != nil {
		encodeHuff(c.w, t, 0, true)
	} else {
		c.e.encodeInt(cx, 0, true)
	}
}

// refinement codes a refined bitmap: inline with the arithmetic coder, or as
// a size and a separately coded block with Huffman coding.
func (c *coder) refinement(sizeTable *huffTable, target, ref *Bitmap, dx, dy, template int, at []point) {
	if c.w == nil {
		encodeRefinement(c.e, c.gr, target, ref, dx, dy, template, at, false)
		return
	}
	e := newMQEncoder()
	encodeRefinement(e, c.gr, target, ref, dx, dy, template, at, false)
	data := e.flush()
	encodeHuff(c.w, sizeTable, len(data), false)
	c.w.bytes(data)
}

// textInst is a symbol instance of a text region: symbol id with its reference
// corner at (s, t), optionally refined to target.
type textInst struct {
	id                 int
	s, t               int
	refine             bool
	rdw, rdh, rdx, rdy int
	target             *Bitmap
}

func (in *textInst) bitmap(syms []*Bitmap) *Bitmap {
	if in.refine {
		return in.target
	}
	return syms[in.id]
}

// textConfig holds the parameters of a text region.
type textConfig struct {
	huff, refine bool
	logStrips    int
	corner       int
	transposed   bool
	op, defPixel int
	dsOffset     int
	rtemplate    int
	rat          []point
	syms         []*Bitmap
	codeLen      int
	// Huffman
	fsSel, dsSel, dtSel, rdwSel, rdhSel, rdxSel, rdySel, rsizeSel int
	fs, ds, dt, rdw, rdh, rdx, rdy, rsize                         *huffTable
	codeLens                                                      []int // symbol ID code lengths; nil for fixed-length IDs
	codes                                                         *huffTable
}

// encodeTextData codes the instances (6.4), in their order, grouped into
// strips by consecutive instances of the same strip.
func encodeTextData(c *coder, tc *textConfig, insts []textInst) {
	strips := 1 << tc.logStrips
	stripOf := func(t int) int {
		if t >= 0 {
			return t / strips * strips
		}
		return -((-t + strips - 1) / strips * strips)
	}
	stript := -strips
	c.int(tc.dt, &c.cx.iadt, 1)
	firsts, curs := 0, 0
	for i := 0; i < len(insts); {
		st := stripOf(insts[i].t)
		c.int(tc.dt, &c.cx.iadt, (st-stript)/strips)
		stript = st
		for first := true; i < len(insts) && stripOf(insts[i].t) == st; first = false {
			in := &insts[i]
			bm := in.bitmap(tc.syms)
			w, h := bm.Width, bm.Height
			pre, post := 0, 0
			right := tc.corner == cornerTopRight || tc.corner == cornerBottomRight
			bottom := tc.corner == cornerBottomLeft || tc.corner == cornerBottomRight
			switch {
			case !tc.transposed && right:
				pre = w - 1
			case !tc.transposed:
				post = w - 1
			case bottom:
				pre = h - 1
			default:
				post = h - 1
			}
			target := in.s - pre
			if first {
				c.int(tc.fs, &c.cx.iafs, target-firsts)
				firsts = target
			} else {
				c.int(tc.ds, &c.cx.iads, target-curs-tc.dsOffset)
			}
			curs = in.s + post
			if strips > 1 {
				if c.w != nil {
					c.w.write(uint64(in.t-st), tc.logStrips)
				} else {
					c.e.encodeInt(&c.cx.iait, in.t-st, false)
				}
			}
			switch {
			case c.w == nil:
				c.e.encodeIAID(c.cx.iaid, tc.codeLen, in.id)
			case tc.codes != nil:
				encodeHuff(c.w, tc.codes, in.id, false)
			default:
				c.w.write(uint64(in.id), tc.codeLen)
			}
			if tc.refine {
				ri := 0
				if in.refine {
					ri = 1
				}
				if c.w != nil {
					c.w.write(uint64(ri), 1)
				} else {
					c.e.encodeInt(&c.cx.iari, ri, false)
				}
				if in.refine {
					c.int(tc.rdw, &c.cx.iardw, in.rdw)
					c.int(tc.rdh, &c.cx.iardh, in.rdh)
					c.int(tc.rdx, &c.cx.iardx, in.rdx)
					c.int(tc.rdy, &c.cx.iardy, in.rdy)
					c.refinement(tc.rsize, in.target, tc.syms[in.id], in.rdw>>1+in.rdx, in.rdh>>1+in.rdy, tc.rtemplate, tc.rat)
				}
			}
			i++
		}
		c.oob(tc.ds, &c.cx.iads)
	}
}

// symbolIDTable codes the symbol ID Huffman table of a text region (7.4.3.1.7)
// and returns the table.
func symbolIDTable(w *bitWriter, lens []int) *huffTable {
	// Run code sequence.
	type run struct{ code, extra, bits int }
	var runs []run
	for i := 0; i < len(lens); {
		j := i
		for j < len(lens) && lens[j] == lens[i] {
			j++
		}
		n := j - i
		switch {
		case lens[i] == 0 && n >= 11:
			n = min(n, 138)
			runs = append(runs, run{34, n - 11, 7})
		case lens[i] == 0 && n >= 3:
			n = min(n, 10)
			runs = append(runs, run{33, n - 3, 3})
		case lens[i] != 0 && n >= 4:
			runs = append(runs, run{lens[i], 0, 0})
			n = min(n-1, 6)
			runs = append(runs, run{32, n - 3, 2})
			n++
		default:
			n = 1
			runs = append(runs, run{lens[i], 0, 0})
		}
		i += n
	}
	freq := make([]int, 35)
	for _, r := range runs {
		freq[r.code]++
	}
	runLens := huffLengths(freq, 15)
	var runLines []huffLine
	for i, l := range runLens {
		w.write(uint64(l), 4)
		runLines = append(runLines, huffLine{prefLen: l, low: int64(i)})
	}
	runTable, _ := newHuffTable(runLines)
	for _, r := range runs {
		encodeHuff(w, runTable, r.code, false)
		w.write(uint64(r.extra), r.bits)
	}
	w.align()
	lines := make([]huffLine, len(lens))
	for i, l := range lens {
		lines[i] = huffLine{prefLen: l, low: int64(i)}
	}
	t, _ := newHuffTable(lines)
	return t
}

// textRegion codes a text region segment's data.
func textRegion(w, h, x, y, op int, tc *textConfig, insts []textInst) []byte {
	b := regionInfoBytes(w, h, x, y, op)
	flags := tc.logStrips<<2 | tc.corner<<4 | tc.op<<7 | tc.defPixel<<9 | (tc.dsOffset&0x1f)<<10 | tc.rtemplate<<15
	if tc.huff {
		flags |= 1
	}
	if tc.refine {
		flags |= 2
	}
	if tc.transposed {
		flags |= 0x40
	}
	b = binary.BigEndian.AppendUint16(b, uint16(flags))
	if tc.huff {
		hf := tc.fsSel | tc.dsSel<<2 | tc.dtSel<<4 | tc.rdwSel<<6 | tc.rdhSel<<8 | tc.rdxSel<<10 | tc.rdySel<<12 | tc.rsizeSel<<14
		b = binary.BigEndian.AppendUint16(b, uint16(hf))
		pickStd := func(t **huffTable, sel int, std ...int) {
			if *t == nil {
				*t = standardTables[std[sel]]
			}
		}
		pickStd(&tc.fs, tc.fsSel, 6, 7)
		pickStd(&tc.ds, tc.dsSel, 8, 9, 10)
		pickStd(&tc.dt, tc.dtSel, 11, 12, 13)
		pickStd(&tc.rdw, tc.rdwSel, 14, 15)
		pickStd(&tc.rdh, tc.rdhSel, 14, 15)
		pickStd(&tc.rdx, tc.rdxSel, 14, 15)
		pickStd(&tc.rdy, tc.rdySel, 14, 15)
		pickStd(&tc.rsize, tc.rsizeSel, 1)
	}
	if tc.refine && tc.rtemplate == 0 {
		b = append(b, atBytes(tc.rat)...)
	}
	b = binary.BigEndian.AppendUint32(b, uint32(len(insts)))
	tc.codeLen = ceilLog2(len(tc.syms))
	c := &coder{cx: newTextContexts(tc.codeLen)}
	if tc.refine {
		c.gr = newRefinementStats(tc.rtemplate)
	}
	if tc.huff {
		c.w = &bitWriter{}
		lens := tc.codeLens
		if lens == nil {
			freq := make([]int, len(tc.syms))
			for _, in := range insts {
				freq[in.id]++
			}
			lens = huffLengths(freq, 31)
		}
		tc.codes = symbolIDTable(c.w, lens)
		encodeTextData(c, tc, insts)
		return append(b, c.w.out...)
	}
	c.e = newMQEncoder()
	encodeTextData(c, tc, insts)
	return append(b, c.e.flush()...)
}

// sdSymbol is a new symbol of a symbol dictionary. With refinement/aggregate
// coding it is either one refined symbol (ref, rdx, rdy) or several instances
// (agg) drawn with the text region procedure.
type sdSymbol struct {
	bm       *Bitmap
	ref      int
	rdx, rdy int
	agg      []textInst
}

// sdConfig holds the parameters of a symbol dictionary.
type sdConfig struct {
	huff, refagg        bool
	template, rtemplate int
	at, rat             []point
	dhSel, dwSel        int
	bmSel, aiSel        int
	dh, dw, bmsize, ai  *huffTable // custom tables where selected
	mmr                 bool       // MMR coded height class bitmaps
	in                  []*Bitmap
	syms                []sdSymbol // grouped into height classes by consecutive heights
	export              []bool     // for input and new symbols; nil exports the new ones
	ctxUsed, ctxRetain  bool
	gb, gr              []uint8 // contexts to start from, and the retained ones afterwards
}

// symbolDictData codes a symbol dictionary segment's data.
func symbolDictData(c *sdConfig) []byte {
	flags := c.dhSel<<2 | c.dwSel<<4 | c.bmSel<<6 | c.aiSel<<7 | c.template<<10 | c.rtemplate<<12
	if c.huff {
		flags |= 1
	}
	if c.refagg {
		flags |= 2
	}
	if c.ctxUsed {
		flags |= 0x100
	}
	if c.ctxRetain {
		flags |= 0x200
	}
	b := binary.BigEndian.AppendUint16(nil, uint16(flags))
	if !c.huff {
		b = append(b, atBytes(c.at)...)
	}
	if c.refagg && c.rtemplate == 0 {
		b = append(b, atBytes(c.rat)...)
	}
	total := len(c.in) + len(c.syms)
	export := c.export
	if export == nil {
		export = make([]bool, total)
		for i := len(c.in); i < total; i++ {
			export[i] = true
		}
	}
	numEx := 0
	for _, e := range export {
		if e {
			numEx++
		}
	}
	b = binary.BigEndian.AppendUint32(b, uint32(numEx))
	b = binary.BigEndian.AppendUint32(b, uint32(len(c.syms)))

	codeLen := ceilLog2(total)
	cd := &coder{cx: newTextContexts(codeLen)}
	var iadh, iadw, iaex, iaai intContext
	if c.huff {
		cd.w = &bitWriter{}
		pick := func(t **huffTable, sel int, std ...int) {
			if *t == nil {
				*t = standardTables[std[sel]]
			}
		}
		pick(&c.dh, c.dhSel, 4, 5)
		pick(&c.dw, c.dwSel, 2, 3)
		pick(&c.bmsize, c.bmSel, 1)
		pick(&c.ai, c.aiSel, 1)
	} else {
		cd.e = newMQEncoder()
		if !c.ctxUsed || c.gb == nil {
			c.gb = newGenericStats(c.template)
		}
	}
	if c.refagg && (!c.ctxUsed || c.gr == nil) {
		c.gr = newRefinementStats(c.rtemplate)
	}
	cd.gr = c.gr
	all := append([]*Bitmap{}, c.in...)
	height := 0
	for i := 0; i < len(c.syms); {
		j := i
		for j < len(c.syms) && c.syms[j].bm.Height == c.syms[i].bm.Height {
			j++
		}
		h := c.syms[i].bm.Height
		cd.int(c.dh, &iadh, h-height)
		height = h
		width, totWidth := 0, 0
		for k := i; k < j; k++ {
			s := &c.syms[k]
			cd.int(c.dw, &iadw, s.bm.Width-width)
			width = s.bm.Width
			totWidth += width
			switch {
			case c.huff && !c.refagg:
			case !c.refagg:
				encodeGeneric(cd.e, c.gb, s.bm, c.template, c.at, false, nil)
			case len(s.agg) > 1:
				cd.int(c.ai, &iaai, len(s.agg))
				tc := &textConfig{huff: c.huff, refine: true, corner: cornerTopLeft, rtemplate: c.rtemplate, rat: c.rat,
					syms: all, codeLen: codeLen,
					fs: standardTables[6], ds: standardTables[8], dt: standardTables[11],
					rdw: standardTables[15], rdh: standardTables[15], rdx: standardTables[15], rdy: standardTables[15],
					rsize: standardTables[1]}
				encodeTextData(cd, tc, s.agg)
			default:
				cd.int(c.ai, &iaai, 1)
				if c.huff {
					cd.w.write(uint64(s.ref), codeLen)
				} else {
					cd.e.encodeIAID(cd.cx.iaid, codeLen, s.ref)
				}
				cd.int(standardTables[15], &cd.cx.iardx, s.rdx)
				cd.int(standardTables[15], &cd.cx.iardy, s.rdy)
				cd.refinement(standardTables[1], s.bm, all[s.ref], s.rdx, s.rdy, c.rtemplate, c.rat)
			}
			all = append(all, s.bm)
		}
		cd.oob(c.dw, &iadw)
		if c.huff && !c.refagg {
			coll := newBitmap(totWidth, h)
			x := 0
			for k := i; k < j; k++ {
				coll.compose(c.syms[k].bm, x, 0, opOr)
				x += c.syms[k].bm.Width
			}
			if c.mmr {
				data := encodeMMR(coll, i%2 == 0) // the EOFB is optional
				encodeHuff(cd.w, c.bmsize, len(data), false)
				cd.w.bytes(data)
			} else {
				encodeHuff(cd.w, c.bmsize, 0, false)
				cd.w.bytes(coll.Data)
			}
		}
		i = j
	}
	// Export flags as alternating run lengths, starting with non-exported.
	cur := false
	for i := 0; i < total; {
		j := i
		for j < total && export[j] == cur {
			j++
		}
		if c.huff {
			encodeHuff(cd.w, standardTables[1], j-i, false)
		} else {
			cd.e.encodeInt(&iaex, j-i, false)
		}
		i = j
		cur = !cur
	}
	if c.huff {
		return append(b, cd.w.out...)
	}
	return append(b, cd.e.flush()...)
}

// exported returns the symbols a dictionary exports.
func (c *sdConfig) exported() []*Bitmap {
	var all []*Bitmap
	all = append(all, c.in...)
	for _, s := range c.syms {
		all = append(all, s.bm)
	}
	if c.export == nil {
		return all[len(c.in):]
	}
	var ex []*Bitmap
	for i, e := range c.export {
		if e {
			ex = append(ex, all[i])
		}
	}
	return ex
}

// patternDictData codes a pattern dictionary segment's data.
func patternDictData(pats []*Bitmap, mmr bool, template int) []byte {
	pw, ph := pats[0].Width, pats[0].Height
	coll := newBitmap(pw*len(pats), ph)
	for i, p := range pats {
		coll.compose(p, i*pw, 0, opOr)
	}
	flags := byte(template << 1)
	if mmr {
		flags = 1
	}
	b := append([]byte{flags, byte(pw), byte(ph)}, binary.BigEndian.AppendUint32(nil, uint32(len(pats)-1))...)
	if mmr {
		return append(b, encodeMMR(coll, false)...)
	}
	at := []point{{-pw, 0}, {-3, -1}, {2, -2}, {-2, -2}}[:len(genericTemplates[template].nominal)]
	e := newMQEncoder()
	encodeGeneric(e, newGenericStats(template), coll, template, at, false, nil)
	return append(b, e.flush()...)
}

// htConfig holds the parameters of a halftone region.
type htConfig struct {
	x, y, w, h, op int // region
	mmr            bool
	template       int
	skip           bool
	hop, defPixel  int
	gw, gh         int
	gx, gy         int32
	rx, ry         uint16
	values         []int // gray values, gw×gh
	pats           []*Bitmap
	bitplaneEOFB   bool
}

func (c *htConfig) pos(mg, ng int) (int, int) {
	x := (int64(c.gx) + int64(mg)*int64(c.ry) + int64(ng)*int64(c.rx)) >> 8
	y := (int64(c.gy) + int64(mg)*int64(c.rx) - int64(ng)*int64(c.ry)) >> 8
	return int(x), int(y)
}

// skipped reports whether grid cell (mg, ng) lies off the region.
func (c *htConfig) skipped(mg, ng int) bool {
	x, y := c.pos(mg, ng)
	pw, ph := c.pats[0].Width, c.pats[0].Height
	return x+pw <= 0 || x >= c.w || y+ph <= 0 || y >= c.h
}

// halftoneRegionData codes a halftone region segment's data. Skipped cells get
// gray value 0.
func halftoneRegionData(c *htConfig) []byte {
	b := regionInfoBytes(c.w, c.h, c.x, c.y, c.op)
	flags := byte(c.template<<1 | c.hop<<4 | c.defPixel<<7)
	if c.mmr {
		flags |= 1
	}
	if c.skip {
		flags |= 8
	}
	b = append(b, flags)
	b = binary.BigEndian.AppendUint32(b, uint32(c.gw))
	b = binary.BigEndian.AppendUint32(b, uint32(c.gh))
	b = binary.BigEndian.AppendUint32(b, uint32(c.gx))
	b = binary.BigEndian.AppendUint32(b, uint32(c.gy))
	b = binary.BigEndian.AppendUint16(b, c.rx)
	b = binary.BigEndian.AppendUint16(b, c.ry)
	var skip *Bitmap
	if c.skip && !c.mmr {
		skip = newBitmap(c.gw, c.gh)
		for mg := 0; mg < c.gh; mg++ {
			for ng := 0; ng < c.gw; ng++ {
				if c.skipped(mg, ng) {
					skip.set(ng, mg)
					c.values[mg*c.gw+ng] = 0
				}
			}
		}
	}
	bpp := ceilLog2(len(c.pats))
	e := newMQEncoder()
	stats := newGenericStats(c.template)
	at := []point{{3, -1}, {-3, -1}, {2, -2}, {-2, -2}}
	if c.template >= 2 {
		at[0].x = 2
	}
	at = at[:len(genericTemplates[c.template].nominal)]
	for j := bpp - 1; j >= 0; j-- {
		plane := newBitmap(c.gw, c.gh)
		for mg := 0; mg < c.gh; mg++ {
			for ng := 0; ng < c.gw; ng++ {
				v := c.values[mg*c.gw+ng]
				if (v>>uint(j)^v>>uint(j+1))&1 != 0 {
					plane.set(ng, mg)
				}
			}
		}
		if c.mmr {
			b = append(b, encodeMMR(plane, c.bitplaneEOFB)...)
		} else {
			encodeGeneric(e, stats, plane, c.template, at, false, skip)
		}
	}
	if !c.mmr {
		b = append(b, e.flush()...)
	}
	return b
}

// Reference rendering.

// refCompose combines src into dst pixel by pixel.
func refCompose(dst, src *Bitmap, x, y, op int) {
	for j := 0; j < src.Height; j++ {
		for i := 0; i < src.Width; i++ {
			dx, dy := x+i, y+j
			if dx < 0 || dy < 0 || dx >= dst.Width || dy >= dst.Height {
				continue
			}
			a, b := dst.pixel(dx, dy), src.pixel(i, j)
			var v int
			switch op {
			case opOr:
				v = a | b
			case opAnd:
				v = a & b
			case opXor:
				v = a ^ b
			case opXnor:
				v = 1 ^ a ^ b
			default:
				v = b
			}
			i0 := dy*dst.Stride + dx>>3
			if v != 0 {
				dst.Data[i0] |= 0x80 >> uint(dx&7)
			} else {
				dst.Data[i0] &^= 0x80 >> uint(dx&7)
			}
		}
	}
}

// renderText draws text region instances as the decoder should.
func renderText(w, h int, tc *textConfig, insts []textInst) *Bitmap {
	bm := newBitmap(w, h)
	if tc.defPixel != 0 {
		bm.fill(1)
	}
	for i := range insts {
		in := &insts[i]
		sym := in.bitmap(tc.syms)
		x, y := in.s, in.t
		if tc.transposed {
			x, y = in.t, in.s
		}
		if tc.corner == cornerTopRight || tc.corner == cornerBottomRight {
			x -= sym.Width - 1
		}
		if tc.corner == cornerBottomLeft || tc.corner == cornerBottomRight {
			y -= sym.Height - 1
		}
		refCompose(bm, sym, x, y, tc.op)
	}
	return bm
}

// renderHalftone draws a halftone region as the decoder should.
func renderHalftone(c *htConfig) *Bitmap {
	bm := newBitmap(c.w, c.h)
	if c.defPixel != 0 {
		bm.fill(1)
	}
	for mg := 0; mg < c.gh; mg++ {
		for ng := 0; ng < c.gw; ng++ {
			x, y := c.pos(mg, ng)
			refCompose(bm, c.pats[c.values[mg*c.gw+ng]], x, y, c.hop)
		}
	}
	return bm
}

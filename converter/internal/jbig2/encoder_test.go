package jbig2

// A JBIG2 encoder for tests: the coding procedures that mirror the decoder's.
// It forms contexts pixel by pixel from T.88's definitions, independently of the
// decoder's sliding windows.

import "sort"

// mqEncoder is the MQ arithmetic encoder (T.88 E.2).
type mqEncoder struct {
	a, c uint32
	ct   int
	out  []byte // out[0] is the byte before the coded data
}

func newMQEncoder() *mqEncoder { return &mqEncoder{a: 0x8000, ct: 12, out: []byte{0}} }

func (e *mqEncoder) encode(cx *uint8, d int) {
	st := *cx
	row := &qeTable[st>>1]
	mps := int(st & 1)
	qe := row.qe
	e.a -= qe
	if d == mps {
		if e.a&0x8000 != 0 {
			e.c += qe
			return
		}
		if e.a < qe {
			e.a = qe
		} else {
			e.c += qe
		}
		*cx = row.nmps<<1 | uint8(mps)
	} else {
		if e.a < qe {
			e.c += qe
		} else {
			e.a = qe
		}
		if row.switchMPS {
			mps = 1 - mps
		}
		*cx = row.nlps<<1 | uint8(mps)
	}
	for {
		e.a <<= 1
		e.c <<= 1
		e.ct--
		if e.ct == 0 {
			e.byteOut()
		}
		if e.a&0x8000 != 0 {
			return
		}
	}
}

func (e *mqEncoder) byteOut() {
	last := len(e.out) - 1
	if e.out[last] == 0xff {
		e.out = append(e.out, byte(e.c>>20))
		e.c &= 0xfffff
		e.ct = 7
		return
	}
	if e.c >= 0x8000000 {
		e.out[last]++
		if e.out[last] == 0xff {
			e.c &= 0x7ffffff
			e.out = append(e.out, byte(e.c>>20))
			e.c &= 0xfffff
			e.ct = 7
			return
		}
	}
	e.out = append(e.out, byte(e.c>>19))
	e.c &= 0x7ffff
	e.ct = 8
}

// flush terminates the code and appends the 0xFF 0xAC marker.
func (e *mqEncoder) flush() []byte {
	t := e.c + e.a
	e.c |= 0xffff
	if e.c >= t {
		e.c -= 0x8000
	}
	e.c <<= uint(e.ct)
	e.byteOut()
	e.c <<= uint(e.ct)
	e.byteOut()
	if e.out[len(e.out)-1] != 0xff {
		e.out = append(e.out, 0xff)
	}
	e.out = append(e.out, 0xac)
	return e.out[1:]
}

// encodeInt encodes v, or OOB, with an integer procedure (A.2).
func (e *mqEncoder) encodeInt(cx *intContext, v int, oob bool) {
	prev := 1
	put := func(b int) {
		e.encode(&cx[prev], b)
		if prev < 256 {
			prev = prev<<1 | b
		} else {
			prev = (prev<<1|b)&511 | 256
		}
	}
	putBits := func(v int64, n int) {
		for i := n - 1; i >= 0; i-- {
			put(int(v >> uint(i) & 1))
		}
	}
	if oob {
		put(1)
		put(0)
		putBits(0, 2)
		return
	}
	a := int64(v)
	if v < 0 {
		put(1)
		a = -a
	} else {
		put(0)
	}
	ranges := []struct {
		n   int
		low int64
	}{{2, 0}, {4, 4}, {6, 20}, {8, 84}, {12, 340}, {32, 4436}}
	for i, r := range ranges {
		if i == len(ranges)-1 || a < r.low+1<<uint(r.n) {
			if i < len(ranges)-1 {
				put(0)
			}
			putBits(a-r.low, r.n)
			return
		}
		put(1)
	}
}

// encodeIAID encodes an n-bit symbol ID (A.3).
func (e *mqEncoder) encodeIAID(cx []uint8, n, v int) {
	prev := 1
	for i := n - 1; i >= 0; i-- {
		b := v >> uint(i) & 1
		e.encode(&cx[prev], b)
		prev = prev<<1 | b
	}
}

// genericContext forms the context of pixel (x, y) of a generic region in the
// bit order of T.88 6.2.5.3.
func genericContext(b *Bitmap, x, y, template int, at []point) int {
	p := func(dx, dy int) int { return b.pixel(x+dx, y+dy) }
	a := func(i int) int { return b.pixel(x+at[i].x, y+at[i].y) }
	switch template {
	case 0:
		return p(-1, 0) | p(-2, 0)<<1 | p(-3, 0)<<2 | p(-4, 0)<<3 | a(0)<<4 |
			p(2, -1)<<5 | p(1, -1)<<6 | p(0, -1)<<7 | p(-1, -1)<<8 | p(-2, -1)<<9 |
			a(1)<<10 | a(2)<<11 | p(1, -2)<<12 | p(0, -2)<<13 | p(-1, -2)<<14 | a(3)<<15
	case 1:
		return p(-1, 0) | p(-2, 0)<<1 | p(-3, 0)<<2 | a(0)<<3 |
			p(2, -1)<<4 | p(1, -1)<<5 | p(0, -1)<<6 | p(-1, -1)<<7 | p(-2, -1)<<8 |
			p(2, -2)<<9 | p(1, -2)<<10 | p(0, -2)<<11 | p(-1, -2)<<12
	case 2:
		return p(-1, 0) | p(-2, 0)<<1 | a(0)<<2 |
			p(1, -1)<<3 | p(0, -1)<<4 | p(-1, -1)<<5 | p(-2, -1)<<6 |
			p(1, -2)<<7 | p(0, -2)<<8 | p(-1, -2)<<9
	default:
		return p(-1, 0) | p(-2, 0)<<1 | p(-3, 0)<<2 | p(-4, 0)<<3 | a(0)<<4 |
			p(1, -1)<<5 | p(0, -1)<<6 | p(-1, -1)<<7 | p(-2, -1)<<8 | p(-3, -1)<<9
	}
}

// encodeGeneric codes src with the arithmetic generic region procedure. Pixels
// set in skip are not coded and must be 0 in src.
func encodeGeneric(e *mqEncoder, stats []uint8, src *Bitmap, template int, at []point, tpgdon bool, skip *Bitmap) {
	dec := newBitmap(src.Width, src.Height)
	ltp := 0
	for y := 0; y < src.Height; y++ {
		if tpgdon {
			same := 1
			for x := 0; x < src.Width; x++ {
				if src.pixel(x, y) != src.pixel(x, y-1) {
					same = 0
					break
				}
			}
			e.encode(&stats[genericTemplates[template].sltp], same^ltp)
			ltp = same
			if same == 1 {
				if y > 0 {
					copy(dec.row(y), dec.row(y-1))
				}
				continue
			}
		}
		for x := 0; x < src.Width; x++ {
			if skip != nil && skip.pixel(x, y) != 0 {
				continue
			}
			v := src.pixel(x, y)
			e.encode(&stats[genericContext(dec, x, y, template, at)], v)
			if v != 0 {
				dec.set(x, y)
			}
		}
	}
}

// refinementContext forms the context of pixel (x, y) of a refinement region.
func refinementContext(b, ref *Bitmap, x, y, dx, dy, template int, at []point) int {
	p := func(ax, ay int) int { return b.pixel(x+ax, y+ay) }
	r := func(ax, ay int) int { return ref.pixel(x-dx+ax, y-dy+ay) }
	if template == 0 {
		return p(-1, 0) | p(1, -1)<<1 | p(0, -1)<<2 | p(at[0].x, at[0].y)<<3 |
			r(1, 1)<<4 | r(0, 1)<<5 | r(-1, 1)<<6 | r(1, 0)<<7 | r(0, 0)<<8 | r(-1, 0)<<9 |
			r(1, -1)<<10 | r(0, -1)<<11 | r(at[1].x, at[1].y)<<12
	}
	return p(-1, 0) | p(1, -1)<<1 | p(0, -1)<<2 | p(-1, -1)<<3 |
		r(1, 1)<<4 | r(0, 1)<<5 | r(1, 0)<<6 | r(0, 0)<<7 | r(-1, 0)<<8 | r(0, -1)<<9
}

// typicalRef returns the value of the 3×3 reference neighbourhood of (x, y) if
// it is uniform, or −1.
func typicalRef(ref *Bitmap, x, y int) int {
	s := 0
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			s += ref.pixel(x+i, y+j)
		}
	}
	switch s {
	case 0:
		return 0
	case 9:
		return 1
	}
	return -1
}

// encodeRefinement codes src against ref with the generic refinement procedure.
func encodeRefinement(e *mqEncoder, stats []uint8, src, ref *Bitmap, dx, dy, template int, at []point, tpgron bool) {
	dec := newBitmap(src.Width, src.Height)
	sltp := 0x100
	if template == 1 {
		sltp = 0x80
	}
	ltp := 0
	for y := 0; y < src.Height; y++ {
		typical := 0
		if tpgron {
			typical = 1
			for x := 0; x < src.Width && typical == 1; x++ {
				if v := typicalRef(ref, x-dx, y-dy); v >= 0 && v != src.pixel(x, y) {
					typical = 0
				}
			}
			e.encode(&stats[sltp], typical^ltp)
			ltp = typical
		}
		for x := 0; x < src.Width; x++ {
			v := src.pixel(x, y)
			if typical == 1 && typicalRef(ref, x-dx, y-dy) >= 0 {
				// Predicted, not coded.
			} else {
				e.encode(&stats[refinementContext(dec, ref, x, y, dx, dy, template, at)], v)
			}
			if v != 0 {
				dec.set(x, y)
			}
		}
	}
}

// bitWriter packs bits most significant first.
type bitWriter struct {
	out []byte
	n   int // bits used in the last byte
}

func (w *bitWriter) write(v uint64, n int) {
	for i := n - 1; i >= 0; i-- {
		if w.n == 0 {
			w.out = append(w.out, 0)
		}
		if v>>uint(i)&1 != 0 {
			w.out[len(w.out)-1] |= 0x80 >> uint(w.n)
		}
		w.n = (w.n + 1) & 7
	}
}

func (w *bitWriter) align() { w.n = 0 }

func (w *bitWriter) bytes(b []byte) {
	w.align()
	w.out = append(w.out, b...)
}

// changes returns the changing elements of row y of b: the positions whose
// colour differs from the pixel before, starting from white.
func changes(b *Bitmap, y int) []int {
	var c []int
	prev := 0
	for x := 0; x < b.Width; x++ {
		if v := b.pixel(x, y); v != prev {
			c = append(c, x)
			prev = v
		}
	}
	return c
}

// encodeMMR codes b with T.6, followed by an end-of-facsimile-block code if
// withEOFB.
func encodeMMR(b *Bitmap, withEOFB bool) []byte {
	var w bitWriter
	width := b.Width
	ref := []int{width, width, width}
	find := func(line []int, i *int, a0, color int) int {
		for *i > 0 && line[*i-1] > a0 {
			*i--
		}
		for line[*i] <= a0 {
			*i++
		}
		if *i&1 != color {
			return *i + 1
		}
		return *i
	}
	run := func(color, n int) {
		codes := whiteRunCodes
		if color == 1 {
			codes = blackRunCodes
		}
		for n >= 2560 {
			c := extendedRunCodes[12]
			w.write(uint64(c[0]), c[1])
			n -= 2560
		}
		if n >= 64 {
			m := n / 64
			var c [2]int
			if m <= 27 {
				c = codes[63+m]
			} else {
				c = extendedRunCodes[m-28]
			}
			w.write(uint64(c[0]), c[1])
			n -= m * 64
		}
		w.write(uint64(codes[n][0]), codes[n][1])
	}
	mode := func(m int) {
		for _, c := range modeCodes {
			if c[0] == m {
				w.write(uint64(c[1]), c[2])
			}
		}
	}
	for y := 0; y < b.Height; y++ {
		cur := append(changes(b, y), width, width, width)
		a0, color, ic, ir := -1, 0, 0, 0
		for a0 < width {
			i := find(cur, &ic, a0, color)
			a1, a2 := cur[i], cur[i+1]
			j := find(ref, &ir, a0, color)
			b1, b2 := ref[j], ref[j+1]
			switch {
			case b2 < a1:
				mode(modePass)
				a0 = b2
			case a1-b1 >= -3 && a1-b1 <= 3:
				mode([]int{modeVL3, modeVL2, modeVL1, modeV0, modeVR1, modeVR2, modeVR3}[a1-b1+3])
				a0 = a1
				color = 1 - color
			default:
				mode(modeHorizontal)
				run(color, a1-max(a0, 0))
				run(1-color, a2-a1)
				a0 = a2
			}
		}
		ref = cur
	}
	if withEOFB {
		w.write(eofb, 24)
	}
	return w.out
}

// huffCodes assigns the prefix codes of a table's lines (B.3).
func huffCodes(lines []huffLine) []int {
	var count [34]int
	maxLen := 0
	for _, l := range lines {
		if l.prefLen > 0 {
			count[l.prefLen]++
			maxLen = max(maxLen, l.prefLen)
		}
	}
	codes := make([]int, len(lines))
	first := 0
	for n := 1; n <= maxLen; n++ {
		first = (first + count[n-1]) << 1
		if n == 1 {
			first = 0
		}
		c := first
		for i, l := range lines {
			if l.prefLen == n {
				codes[i] = c
				c++
			}
		}
	}
	return codes
}

// encodeHuff writes v, or OOB, with a Huffman table.
func encodeHuff(w *bitWriter, t *huffTable, v int, oob bool) {
	codes := huffCodes(t.lines)
	for i, l := range t.lines {
		if l.prefLen == 0 {
			continue
		}
		switch {
		case oob && l.kind == lineOOB:
			w.write(uint64(codes[i]), l.prefLen)
			return
		case oob:
		case l.kind == lineNormal && int64(v) >= l.low && int64(v) < l.low+1<<uint(l.rangeLen):
			w.write(uint64(codes[i]), l.prefLen)
			w.write(uint64(int64(v)-l.low), l.rangeLen)
			return
		case l.kind == lineLower && int64(v) <= l.low:
			w.write(uint64(codes[i]), l.prefLen)
			w.write(uint64(l.low-int64(v)), 32)
			return
		case l.kind == lineUpper && int64(v) >= l.low:
			w.write(uint64(codes[i]), l.prefLen)
			w.write(uint64(int64(v)-l.low), 32)
			return
		}
	}
	panic("value not in Huffman table")
}

// huffLengths returns Huffman code lengths for symbol frequencies (0 for
// unused symbols), at most maxLen.
func huffLengths(freq []int, maxLen int) []int {
	type node struct {
		f    int
		syms []int
	}
	var nodes []node
	for i, f := range freq {
		if f > 0 {
			nodes = append(nodes, node{f, []int{i}})
		}
	}
	lens := make([]int, len(freq))
	if len(nodes) == 1 {
		lens[nodes[0].syms[0]] = 1
		return lens
	}
	for len(nodes) > 1 {
		sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].f < nodes[j].f })
		a, b := nodes[0], nodes[1]
		for _, s := range append(append([]int{}, a.syms...), b.syms...) {
			lens[s]++
		}
		nodes = append([]node{{a.f + b.f, append(append([]int{}, a.syms...), b.syms...)}}, nodes[2:]...)
	}
	for _, l := range lens {
		if l > maxLen {
			panic("Huffman code too long")
		}
	}
	return lens
}

package jbig2

// The streams of the test fixtures: each case builds a stream with the test
// encoder and the page it should decode to with the reference renderer.

import (
	"encoding/binary"
	"math"
	"math/rand"
)

type fixture struct {
	name          string
	page, globals []byte
	want          *Bitmap
	// reference, if set, says why Ghostscript cannot check the fixture; its
	// expected page is then the reference rendering.
	reference string
}

// Drawing helpers.

func fillRect(b *Bitmap, x, y, w, h int) {
	for j := max(y, 0); j < min(y+h, b.Height); j++ {
		for i := max(x, 0); i < min(x+w, b.Width); i++ {
			b.set(i, j)
		}
	}
}

func fillEllipse(b *Bitmap, cx, cy, rx, ry float64) {
	for j := 0; j < b.Height; j++ {
		for i := 0; i < b.Width; i++ {
			dx, dy := (float64(i)+0.5-cx)/rx, (float64(j)+0.5-cy)/ry
			if dx*dx+dy*dy <= 1 {
				b.set(i, j)
			}
		}
	}
}

// stroke draws a thick line.
func stroke(b *Bitmap, x0, y0, x1, y1, r float64) {
	for j := 0; j < b.Height; j++ {
		for i := 0; i < b.Width; i++ {
			px, py := float64(i)+0.5, float64(j)+0.5
			vx, vy := x1-x0, y1-y0
			t := 0.0
			if l := vx*vx + vy*vy; l > 0 {
				t = math.Max(0, math.Min(1, ((px-x0)*vx+(py-y0)*vy)/l))
			}
			dx, dy := px-(x0+t*vx), py-(y0+t*vy)
			if dx*dx+dy*dy <= r*r {
				b.set(i, j)
			}
		}
	}
}

// blobImage draws random ellipses, bars and strokes, leaving some rows blank
// and some rows repeated, as typical prediction likes.
func blobImage(w, h int, seed int64) *Bitmap {
	r := rand.New(rand.NewSource(seed))
	b := newBitmap(w, h)
	for i := 0; i < 6; i++ {
		fillEllipse(b, r.Float64()*float64(w), r.Float64()*float64(h)*0.8, 4+r.Float64()*float64(w)/6, 3+r.Float64()*float64(h)/6)
	}
	for i := 0; i < 4; i++ {
		fillRect(b, r.Intn(w), r.Intn(h), 2+r.Intn(w/3), 1+r.Intn(8))
	}
	for i := 0; i < 5; i++ {
		stroke(b, r.Float64()*float64(w), r.Float64()*float64(h), r.Float64()*float64(w), r.Float64()*float64(h), 0.8+r.Float64()*2)
	}
	for i := 0; i < w*h/60; i++ { // speckles
		b.set(r.Intn(w), r.Intn(h))
	}
	for y := h * 7 / 8; y < h; y++ { // blank rows at the bottom
		for i := range b.row(y) {
			b.row(y)[i] = 0
		}
	}
	return b
}

// glyph draws a letter-like shape.
func glyph(r *rand.Rand, w, h int) *Bitmap {
	b := newBitmap(w, h)
	fw, fh := float64(w), float64(h)
	pt := func() (float64, float64) { return 1 + r.Float64()*(fw-2), 1 + r.Float64()*(fh-2) }
	for i := 0; i < 2+r.Intn(3); i++ {
		x0, y0 := pt()
		x1, y1 := pt()
		stroke(b, x0, y0, x1, y1, 0.7+r.Float64()*0.9)
	}
	if r.Intn(2) == 0 {
		fillEllipse(b, fw/2, fh/2, fw/3, fh/4)
		inner := newBitmap(w, h)
		fillEllipse(inner, fw/2, fh/2, fw/3-1.5, fh/4-1.5)
		refCompose(b, inner, 0, 0, opXor)
	}
	return b
}

// glyphs draws n glyphs with the given heights (cycled) and random widths.
func glyphs(seed int64, n int, heights ...int) []*Bitmap {
	r := rand.New(rand.NewSource(seed))
	g := make([]*Bitmap, n)
	for i := range g {
		h := heights[i%len(heights)]
		g[i] = glyph(r, 4+r.Intn(h), h)
	}
	return g
}

// sortByHeight orders symbols into height classes, by width within a class.
func sortByHeight(g []*Bitmap) {
	for i := 1; i < len(g); i++ {
		for j := i; j > 0 && (g[j].Height < g[j-1].Height || g[j].Height == g[j-1].Height && g[j].Width < g[j-1].Width); j-- {
			g[j], g[j-1] = g[j-1], g[j]
		}
	}
}

// flip returns a copy of b with n pixels inverted.
func flip(b *Bitmap, r *rand.Rand, n int) *Bitmap {
	c := newBitmap(b.Width, b.Height)
	copy(c.Data, b.Data)
	for i := 0; i < n; i++ {
		x, y := r.Intn(b.Width), r.Intn(b.Height)
		c.Data[y*c.Stride+x>>3] ^= 0x80 >> uint(x&7)
	}
	return c
}

// resized returns b drawn at (dx, dy) on a w×h bitmap with a few pixels flipped.
func resized(b *Bitmap, r *rand.Rand, w, h, dx, dy int) *Bitmap {
	c := newBitmap(w, h)
	refCompose(c, b, dx, dy, opOr)
	return flip(c, r, 3)
}

func plainSD(syms []*Bitmap, template int, at []point) *sdConfig {
	c := &sdConfig{template: template, at: at}
	for _, s := range syms {
		c.syms = append(c.syms, sdSymbol{bm: s})
	}
	return c
}

func nominalAT(template int) []point { return genericTemplates[template].nominal }

// lineOfText lays out instances left to right on baseline rows.
func lineOfText(r *rand.Rand, syms []*Bitmap, x0, y0, width, lines, lineH int, corner int) []textInst {
	var insts []textInst
	for l := 0; l < lines; l++ {
		x := x0 + r.Intn(3)
		t := y0 + l*lineH
		for x < x0+width {
			id := r.Intn(len(syms))
			s := syms[id]
			in := textInst{id: id, s: x, t: t}
			switch corner {
			case cornerTopRight:
				in.s = x + s.Width - 1
			case cornerBottomLeft:
				in.t = t + s.Height - 1
			case cornerBottomRight:
				in.s, in.t = x+s.Width-1, t+s.Height-1
			}
			insts = append(insts, in)
			x += s.Width + r.Intn(4) - 1
		}
	}
	return insts
}

// sortStrips orders instances by strip, keeping their order within a strip.
func sortStrips(insts []textInst, logStrips int) {
	strip := func(t int) int { return t >> uint(logStrips) }
	for i := 1; i < len(insts); i++ {
		for j := i; j > 0 && strip(insts[j].t) < strip(insts[j-1].t); j-- {
			insts[j], insts[j-1] = insts[j-1], insts[j]
		}
	}
}

// page starts a stream with a page information segment and its model.
func newPage(w, h int, flags byte) (*stream, *Bitmap) {
	s := &stream{}
	s.add(typePageInfo, pageInfo(uint32(w), uint32(h), flags, 0), segOpts{})
	want := newBitmap(w, h)
	if flags&4 != 0 {
		want.fill(1)
	}
	return s, want
}

func genericCase(name string, template int, at []point, tpgdon, mmr bool) fixture {
	s, want := newPage(203, 131, 0)
	img := blobImage(181, 113, int64(len(name)))
	s.add(typeImmediateGeneric, genericRegion(img, 11, 7, opOr, mmr, template, at, tpgdon), segOpts{})
	refCompose(want, img, 11, 7, opOr)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: name, page: s.buf, want: want}
}

// unknownLengthRegion codes an immediate generic region whose segment data
// length is unknown: the data ends with a marker and the row count.
func unknownLengthRegion(img *Bitmap, x, y, op int, mmr bool) []byte {
	b := genericRegion(img, x, y, op, mmr, 0, nominalAT(0), true)
	if mmr {
		b = append(b, 0, 0)
	}
	return binary.BigEndian.AppendUint32(b, uint32(img.Height))
}

func caseUnknownLength() fixture {
	s, want := newPage(160, 150, 0)
	a, b := blobImage(150, 60, 3), blobImage(140, 70, 4)
	s.add(typeImmediateGeneric, unknownLengthRegion(a, 5, 5, opOr, false), segOpts{unknown: true})
	refCompose(want, a, 5, 5, opOr)
	s.add(typeImmediateGeneric, unknownLengthRegion(b, 12, 70, opXor, true), segOpts{unknown: true})
	refCompose(want, b, 12, 70, opXor)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "generic_unknown_length", page: s.buf, want: want}
}

func caseRefinement() fixture {
	r := rand.New(rand.NewSource(7))
	s, want := newPage(220, 170, 2|0x20)
	// A generic region, then a refinement of the page. (jbig2dec refines as
	// T.88 does only a region covering the whole page.)
	a := blobImage(120, 90, 8)
	s.add(typeImmediateGeneric, genericRegion(a, 4, 4, opOr, false, 1, nominalAT(1), true), segOpts{})
	refCompose(want, a, 4, 4, opOr)
	b := flip(want, r, 400)
	fillRect(b, 10, 30, 20, 3)
	s.add(typeImmediateRefinement, refinementRegion(b, want, 0, 0, opReplace, 0, []point{{-2, -1}, {-2, 0}}, true), segOpts{})
	want = b
	// An intermediate text region, refined twice.
	g := glyphs(10, 12, 12, 15)
	sd := plainSD(g, 2, nominalAT(2))
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{})
	tc := &textConfig{syms: sd.exported(), corner: cornerTopLeft}
	in := lineOfText(r, tc.syms, 1, 1, 80, 3, 18, cornerTopLeft)
	n = s.add(typeIntermediateText, textRegion(90, 60, 125, 100, opOr, tc, in), segOpts{refs: []uint32{n}})
	c := renderText(90, 60, tc, in)
	d := flip(c, r, 200)
	n = s.add(typeIntermediateRefinement, refinementRegion(d, c, 125, 100, opOr, 1, nil, false), segOpts{refs: []uint32{n}})
	e := flip(d, r, 100)
	s.add(typeImmediateLosslessRefinement, refinementRegion(e, d, 125, 100, opOr, 0, []point{{-1, -1}, {-1, -1}}, true), segOpts{refs: []uint32{n}})
	refCompose(want, e, 125, 100, opOr)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "refinement", page: s.buf, want: want}
}

// textPage builds a page with a symbol dictionary and one text region.
func textPage(name string, sd *sdConfig, tc *textConfig, insts func(syms []*Bitmap) []textInst) fixture {
	const w, h = 240, 160
	s, want := newPage(w, h, 0)
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{})
	tc.syms = sd.exported()
	in := insts(tc.syms)
	s.add(typeImmediateText, textRegion(w-10, h-10, 5, 5, opOr, tc, in), segOpts{refs: []uint32{n}})
	refCompose(want, renderText(w-10, h-10, tc, in), 5, 5, opOr)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: name, page: s.buf, want: want}
}

func caseTextArith() fixture {
	g := glyphs(11, 24, 12, 16, 20)
	sortByHeight(g)
	r := rand.New(rand.NewSource(12))
	return textPage("text_arith", plainSD(g, 0, nominalAT(0)), &textConfig{corner: cornerTopLeft},
		func(syms []*Bitmap) []textInst { return lineOfText(r, syms, 2, 2, 220, 6, 24, cornerTopLeft) })
}

func caseTextStrips() fixture {
	g := glyphs(13, 20, 10, 14, 15)
	sortByHeight(g)
	r := rand.New(rand.NewSource(14))
	sd := plainSD(g, 1, []point{{-3, -1}})
	tc := &textConfig{corner: cornerBottomLeft, logStrips: 2, dsOffset: -3}
	return textPage("text_strips", sd, tc, func(syms []*Bitmap) []textInst {
		in := lineOfText(r, syms, 2, 2, 220, 7, 21, cornerBottomLeft)
		for i := range in { // move instances within their 4-row strips
			in[i].t += r.Intn(3)
		}
		return in
	})
}

// caseTextCorners draws eight text regions: each reference corner, transposed or not.
func caseTextCorners() fixture {
	const w, h = 300, 260
	s, want := newPage(w, h, 0)
	g := glyphs(15, 16, 9, 13)
	sortByHeight(g)
	sd := plainSD(g, 2, []point{{-3, -1}})
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{})
	syms := sd.exported()
	r := rand.New(rand.NewSource(16))
	for i := 0; i < 8; i++ {
		tc := &textConfig{syms: syms, corner: i % 4, transposed: i >= 4, logStrips: i % 3, dsOffset: i - 4}
		var in []textInst
		t := 3
		for k := 0; k < 4; k++ { // strips
			sPos := 2 + r.Intn(4)
			for j := 0; j < 5; j++ {
				id := r.Intn(len(syms))
				sym := syms[id]
				along, across := sym.Width, sym.Height
				if tc.transposed {
					along, across = across, along
				}
				inst := textInst{id: id, s: sPos, t: t + r.Intn(1<<tc.logStrips)}
				// Place the reference corner so the symbol starts at sPos.
				right := tc.corner == cornerTopRight || tc.corner == cornerBottomRight
				bottom := tc.corner == cornerBottomLeft || tc.corner == cornerBottomRight
				if !tc.transposed && right || tc.transposed && bottom {
					inst.s += along - 1
				}
				if !tc.transposed && bottom || tc.transposed && right {
					inst.t += across - 1
				}
				in = append(in, inst)
				sPos += along + 1
			}
			t += 16
		}
		x, y := 5+(i%4)*74, 5+(i/4)*128
		s.add(typeImmediateText, textRegion(72, 124, x, y, opOr, tc, in), segOpts{refs: []uint32{n}})
		refCompose(want, renderText(72, 124, tc, in), x, y, opOr)
	}
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "text_corners", page: s.buf, want: want}
}

// caseCombination composes regions with each combination operator, onto a
// page with default pixel 1, and text regions with each text combination
// operator and default pixel.
func caseCombination() fixture {
	const w, h = 260, 200
	s, want := newPage(w, h, 4|0x40)
	base := blobImage(200, 150, 21)
	s.add(typeImmediateGeneric, genericRegion(base, 30, 25, opReplace, false, 0, nominalAT(0), true), segOpts{})
	refCompose(want, base, 30, 25, opReplace)
	for op := opOr; op <= opReplace; op++ {
		img := blobImage(70, 60, int64(30+op))
		x, y := 10+op*45, 10+op*25
		s.add(typeImmediateGeneric, genericRegion(img, x, y, op, false, 3, nominalAT(3), false), segOpts{})
		refCompose(want, img, x, y, op)
	}
	g := glyphs(22, 10, 11)
	sd := plainSD(g, 0, []point{{2, -2}, {-4, -1}, {1, -3}, {-3, -2}})
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{})
	syms := sd.exported()
	r := rand.New(rand.NewSource(23))
	for op := opOr; op <= opXnor; op++ { // text regions have no REPLACE
		tc := &textConfig{syms: syms, corner: cornerTopLeft, op: op, defPixel: op & 1}
		in := lineOfText(r, syms, 1, 1, 40, 3, 9, cornerTopLeft)
		x, y := 5+op*50, 140
		s.add(typeImmediateText, textRegion(48, 30, x, y, (op+2)%5, tc, in), segOpts{refs: []uint32{n}})
		refCompose(want, renderText(48, 30, tc, in), x, y, (op+2)%5)
	}
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "combination", page: s.buf, want: want}
}

func caseTextHuffman() fixture {
	g := glyphs(31, 30, 11, 15, 18, 22)
	sortByHeight(g)
	r := rand.New(rand.NewSource(32))
	sd := plainSD(g, 0, nil)
	sd.huff = true
	// Tables B.5, B.3 and B.12; the globals case has the defaults.
	sd.dhSel, sd.dwSel = 1, 1
	tc := &textConfig{huff: true, corner: cornerTopLeft, dtSel: 1}
	return textPage("text_huffman", sd, tc, func(syms []*Bitmap) []textInst {
		return lineOfText(r, syms[:25], 2, 2, 220, 5, 28, cornerTopLeft) // some symbols unused
	})
}

// customTables returns custom Huffman tables and their segment data.
func customTables() (segs [][]byte, fs, ds, dt, dh, dw, bm *huffTable) {
	add := func(low, high int32, lines [][2]int, lower, upper, oob int) *huffTable {
		b, t := tableSegment(low, high, lines, lower, upper, oob)
		segs = append(segs, b)
		return t
	}
	fs = add(-64, 64, [][2]int{{2, 5}, {2, 5}, {2, 5}, {3, 5}}, 4, 4, 0)
	ds = add(-8, 8, [][2]int{{2, 2}, {2, 2}, {3, 2}, {3, 2}}, 4, 4, 3)
	dt = add(1, 17, [][2]int{{1, 0}, {2, 0}, {3, 1}, {4, 2}, {5, 3}}, 6, 6, 0)
	dh = add(0, 64, [][2]int{{1, 5}, {2, 5}}, 3, 3, 0)
	dw = add(-4, 64, [][2]int{{2, 2}, {2, 4}, {3, 5}, {4, 4}}, 5, 5, 3)
	bm = add(0, 3328, [][2]int{{1, 8}, {2, 10}, {3, 11}}, 4, 4, 0)
	return
}

// caseTextHuffmanCustom uses custom tables and MMR coded height class bitmaps
// in the dictionary, and custom and other standard tables with strips in two
// text regions.
func caseTextHuffmanCustom() fixture {
	const w, h = 240, 200
	s, want := newPage(w, h, 0)
	segs, fs, ds, dt, dh, dw, bmt := customTables()
	var tables []uint32
	for _, b := range segs {
		tables = append(tables, s.add(typeTables, b, segOpts{}))
	}
	g := glyphs(41, 24, 12, 17, 13)
	sortByHeight(g)
	// Decreasing widths within some height classes need negative deltas (B.3 or custom).
	for i := 0; i+1 < len(g); i += 3 {
		if g[i].Height == g[i+1].Height {
			g[i], g[i+1] = g[i+1], g[i]
		}
	}
	sd := plainSD(g, 0, nil)
	sd.huff, sd.mmr = true, true
	sd.dhSel, sd.dwSel, sd.bmSel = 3, 3, 1
	sd.dh, sd.dw, sd.bmsize = dh, dw, bmt
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{refs: []uint32{tables[3], tables[4], tables[5]}})
	syms := sd.exported()
	r := rand.New(rand.NewSource(42))

	tc := &textConfig{huff: true, syms: syms, corner: cornerBottomLeft, logStrips: 1, fsSel: 3, dsSel: 3, dtSel: 3, fs: fs, ds: ds, dt: dt}
	in := lineOfText(r, syms, 2, 2, 220, 4, 22, cornerBottomLeft)
	for i := range in {
		in[i].t += r.Intn(2)
	}
	s.add(typeImmediateText, textRegion(230, 95, 5, 5, opOr, tc, in), segOpts{refs: []uint32{n, tables[0], tables[1], tables[2]}})
	refCompose(want, renderText(230, 95, tc, in), 5, 5, opOr)

	tc2 := &textConfig{huff: true, syms: syms, corner: cornerTopRight, logStrips: 3, fsSel: 1, dsSel: 1, dtSel: 2, dsOffset: 2}
	in2 := lineOfText(r, syms, 2, 2, 220, 4, 22, cornerTopRight)
	for i := range in2 {
		in2[i].t += r.Intn(8)
	}
	sortStrips(in2, tc2.logStrips) // B.13 has no deltas below 1
	s.add(typeImmediateText, textRegion(230, 95, 5, 102, opXor, tc2, in2), segOpts{refs: []uint32{n}})
	refCompose(want, renderText(230, 95, tc2, in2), 5, 102, opXor)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "text_huffman_custom", page: s.buf, want: want}
}

// refaggOpts selects the features of a refinement/aggregate case.
type refaggOpts struct {
	huff       bool
	aggregates bool // symbols aggregated from several instances
	aggRefine  bool // refined instances within aggregates
	rtemplate  int
	reference  string
}

// refaggCase builds a dictionary of refined and aggregated symbols on top of a
// plain one, and a text region with refined instances.
func refaggCase(name string, o refaggOpts) fixture {
	const w, h = 260, 180
	s, want := newPage(w, h, 2)
	r := rand.New(rand.NewSource(51))
	base := glyphs(52, 8, 14, 18)
	sortByHeight(base)
	sdA := plainSD(base, 0, nominalAT(0))
	sdA.huff = o.huff
	nA := s.add(typeSymbolDict, symbolDictData(sdA), segOpts{})
	in := sdA.exported()

	const hc = 20 // one height class
	sdB := &sdConfig{huff: o.huff, refagg: true, template: 1, at: nominalAT(1), rtemplate: o.rtemplate, in: in}
	if o.rtemplate == 0 {
		sdB.rat = []point{{-2, -1}, {-2, 0}}
	}
	if o.huff {
		sdB.at = nil
		sdB.dwSel = 1 // B.3, for negative width deltas
	}
	add := func(sym sdSymbol) { sdB.syms = append(sdB.syms, sym) }
	// Single refinements of input symbols.
	add(sdSymbol{bm: resized(in[2], r, in[2].Width, hc, 0, 1), ref: 2, rdx: 0, rdy: 1})
	add(sdSymbol{bm: resized(in[5], r, in[5].Width+2, hc, 1, 2), ref: 5, rdx: 1, rdy: 2})
	// Aggregates of input symbols, some refined.
	agg := func(ids []int, refineLast bool, all []*Bitmap) sdSymbol {
		var insts []textInst
		x := 0
		for k, id := range ids {
			inst := textInst{id: id, s: x, t: k % 3}
			if refineLast && k == len(ids)-1 {
				inst.refine = true
				inst.rdw, inst.rdh, inst.rdx, inst.rdy = 2, -1, 1, 0
				inst.target = resized(all[id], r, all[id].Width+2, all[id].Height-1, 1, 0)
			}
			insts = append(insts, inst)
			x += inst.bitmap(all).Width + 1
		}
		bm := renderText(x, hc, &textConfig{corner: cornerTopLeft, syms: all}, insts)
		return sdSymbol{bm: bm, agg: insts}
	}
	all := func() []*Bitmap {
		a := append([]*Bitmap{}, in...)
		for _, sym := range sdB.syms {
			a = append(a, sym.bm)
		}
		return a
	}
	if o.aggregates {
		add(agg([]int{0, 1, 3}, false, all()))
		add(agg([]int{4, 6}, o.aggRefine, all()))
		// Symbols built on new symbols of the same dictionary.
		add(agg([]int{len(in), len(in) + 1}, false, all()))
	}
	last := len(in) + len(sdB.syms) - 1
	add(sdSymbol{bm: resized(all()[last], r, all()[last].Width-1, hc, 0, 0), ref: last})
	sdB.export = make([]bool, len(in)+len(sdB.syms))
	for i := range sdB.export {
		sdB.export[i] = i%3 != 1 // export some input symbols again, skip some new ones
	}
	nB := s.add(typeSymbolDict, symbolDictData(sdB), segOpts{refs: []uint32{nA}})
	syms := append(append([]*Bitmap{}, in...), sdB.exported()...)

	tc := &textConfig{huff: o.huff, refine: true, syms: syms, corner: cornerTopLeft, rtemplate: 1 - o.rtemplate}
	if tc.rtemplate == 0 {
		tc.rat = []point{{-3, -1}, {1, -1}}
	}
	if o.huff {
		tc.rdwSel, tc.rdhSel, tc.rdxSel, tc.rdySel = 0, 1, 1, 0
	}
	insts := lineOfText(r, syms, 2, 2, 240, 5, 30, cornerTopLeft)
	for i := range insts {
		if i%4 == 1 {
			in := &insts[i]
			sym := syms[in.id]
			in.refine = true
			in.rdw, in.rdh, in.rdx, in.rdy = r.Intn(3)-1, r.Intn(3)-1, r.Intn(3)-1, r.Intn(3)-1
			in.target = resized(sym, r, sym.Width+in.rdw, sym.Height+in.rdh, in.rdw>>1+in.rdx, in.rdh>>1+in.rdy)
		}
	}
	s.add(typeImmediateText, textRegion(w-10, h-10, 5, 5, opOr, tc, insts), segOpts{refs: []uint32{nA, nB}})
	refCompose(want, renderText(w-10, h-10, tc, insts), 5, 5, opOr)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: name, page: s.buf, want: want, reference: o.reference}
}

// caseBeyondJbig2dec covers what jbig2dec 0.20 decodes differently from
// T.88: refinements of intermediate generic and halftone regions and of part
// of the page, typical prediction with refinement template 1, and a halftone
// region with default pixel 1.
func caseBeyondJbig2dec() fixture {
	s, want := newPage(220, 170, 2|0x20)
	r := rand.New(rand.NewSource(111))
	c := blobImage(90, 60, 112)
	n := s.add(typeIntermediateGeneric, genericRegion(c, 20, 10, opOr, false, 2, nominalAT(2), false), segOpts{})
	d := flip(c, r, 200)
	s.add(typeImmediateRefinement, refinementRegion(d, c, 20, 10, opOr, 1, nil, true), segOpts{refs: []uint32{n}})
	refCompose(want, d, 20, 10, opOr)
	pats := ditherPatterns(8, 6)
	pd := s.add(typePatternDict, patternDictData(pats, false, 3), segOpts{})
	ht := &htConfig{x: 100, y: 90, w: 110, h: 70, op: opXor, template: 3, hop: opOr,
		gw: 20, gh: 13, rx: 6 << 8, pats: pats, values: grayValues(20, 13, 8)}
	n = s.add(typeIntermediateHalftone, halftoneRegionData(ht), segOpts{refs: []uint32{pd}})
	img := renderHalftone(ht)
	e := flip(img, r, 60)
	s.add(typeImmediateRefinement, refinementRegion(e, img, ht.x, ht.y, ht.op, 0, []point{{-1, -1}, {-1, -1}}, true), segOpts{refs: []uint32{n}})
	refCompose(want, e, ht.x, ht.y, ht.op)
	// Part of the page, refined.
	part := newBitmap(80, 50)
	refCompose(part, want, -60, -40, opReplace)
	f := flip(part, r, 120)
	s.add(typeImmediateRefinement, refinementRegion(f, part, 60, 40, opReplace, 0, []point{{-2, -1}, {-2, 0}}, true), segOpts{})
	refCompose(want, f, 60, 40, opReplace)
	// Default pixel 1 in a halftone region the patterns do not cover.
	ht2 := &htConfig{x: 5, y: 110, w: 80, h: 55, op: opOr, template: 1, hop: opXor, defPixel: 1,
		gw: 12, gh: 8, gx: 3 << 8, gy: 2 << 8, rx: 6 << 8, pats: pats, values: grayValues(12, 8, 8)}
	s.add(typeImmediateHalftone, halftoneRegionData(ht2), segOpts{refs: []uint32{pd}})
	refCompose(want, renderHalftone(ht2), ht2.x, ht2.y, ht2.op)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "beyond_jbig2dec", page: s.buf, want: want,
		reference: "jbig2dec does not keep intermediate generic and halftone regions, refines the whole page unshifted instead of the part under the region, " +
			"decodes the typical prediction bit of refinement template 1 in the context of reference pixel (x+1, y) instead of (x, y), " +
			"and fills halftone regions with default pixel 1 as memset(1)"}
}

// ditherPatterns returns n patterns of increasing darkness.
func ditherPatterns(n, size int) []*Bitmap {
	pats := make([]*Bitmap, n)
	order := rand.New(rand.NewSource(int64(n))).Perm(size * size)
	for i := range pats {
		p := newBitmap(size, size)
		for k := 0; k < i*size*size/(n-1); k++ {
			p.set(order[k]%size, order[k]/size)
		}
		pats[i] = p
	}
	return pats
}

func grayValues(gw, gh, n int) []int {
	v := make([]int, gw*gh)
	for y := 0; y < gh; y++ {
		for x := 0; x < gw; x++ {
			d := math.Hypot(float64(x-gw/3), float64(y-gh/2)) / math.Hypot(float64(gw), float64(gh)) * 2.5
			v[y*gw+x] = min(n-1, int(math.Abs(math.Sin(d*3))*float64(n)))
		}
	}
	return v
}

func caseHalftone(name string, mmr bool) fixture {
	const w, h = 230, 190
	s, want := newPage(w, h, 0x20)
	var c *htConfig
	if mmr {
		pats := ditherPatterns(8, 4)
		n := s.add(typePatternDict, patternDictData(pats, true, 0), segOpts{})
		c = &htConfig{x: 10, y: 12, w: 200, h: 160, op: opXor, mmr: true, hop: opReplace, defPixel: 1,
			gw: 50, gh: 40, rx: 4 << 8, pats: pats, values: grayValues(50, 40, 8), bitplaneEOFB: true}
		base := blobImage(w, h, 61)
		s.add(typeImmediateGeneric, genericRegion(base, 0, 0, opOr, false, 0, nominalAT(0), true), segOpts{})
		refCompose(want, base, 0, 0, opOr)
		s.add(typeImmediateHalftone, halftoneRegionData(c), segOpts{refs: []uint32{n}})
		refCompose(want, renderHalftone(c), c.x, c.y, c.op)
	} else {
		pats := ditherPatterns(16, 6)
		n := s.add(typePatternDict, patternDictData(pats, false, 0), segOpts{})
		c = &htConfig{x: 5, y: 5, w: 220, h: 180, op: opOr, template: 0, skip: true, hop: opOr,
			gw: 48, gh: 40, gx: -3 << 8, gy: 2 << 8, rx: 5 << 8, ry: 1 << 8, pats: pats, values: grayValues(48, 40, 16)}
		s.add(typeImmediateHalftone, halftoneRegionData(c), segOpts{refs: []uint32{n}})
		refCompose(want, renderHalftone(c), c.x, c.y, c.op)
		// Templates 1 to 3 in pattern dictionaries.
		for t := 1; t <= 3; t++ {
			pats := ditherPatterns(5+t, 3+t)
			n := s.add(typePatternDict, patternDictData(pats, false, t), segOpts{})
			c := &htConfig{x: 60 * (t - 1), y: 140, w: 55, h: 45, op: opXor, template: t, hop: opXor,
				gw: 12, gh: 10, gx: 1 << 7, gy: -1 << 7, rx: 0x3c0, ry: 0x80, pats: pats, values: grayValues(12, 10, len(pats))}
			s.add(typeImmediateLosslessHalftone, halftoneRegionData(c), segOpts{refs: []uint32{n}})
			refCompose(want, renderHalftone(c), c.x, c.y, c.op)
		}
	}
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: name, page: s.buf, want: want}
}

// caseStriped builds a striped page of unknown height with default pixel 1.
func caseStriped() fixture {
	const w, stripe = 180, 40
	s := &stream{}
	s.add(typePageInfo, pageInfo(w, 0xffffffff, 4|8, 0x8000|stripe), segOpts{})
	want := newBitmap(w, 120)
	want.fill(1)
	a := blobImage(170, 34, 71)
	s.add(typeImmediateGeneric, genericRegion(a, 5, 3, opXor, false, 0, nominalAT(0), false), segOpts{})
	refCompose(want, a, 5, 3, opXor)
	s.add(typeEndOfStripe, binary.BigEndian.AppendUint32(nil, 39), segOpts{})
	b := blobImage(160, 38, 72)
	s.add(typeImmediateGeneric, unknownLengthRegion(b, 10, 41, opReplace, false), segOpts{unknown: true})
	refCompose(want, b, 10, 41, opReplace)
	s.add(typeEndOfStripe, binary.BigEndian.AppendUint32(nil, 79), segOpts{})
	g := glyphs(73, 8, 10)
	sd := plainSD(g, 3, nominalAT(3))
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{})
	tc := &textConfig{syms: sd.exported(), corner: cornerTopLeft, defPixel: 1, op: opXor}
	in := lineOfText(rand.New(rand.NewSource(74)), tc.syms, 1, 1, 150, 3, 12, cornerTopLeft)
	s.add(typeImmediateText, textRegion(160, 40, 8, 80, opAnd, tc, in), segOpts{refs: []uint32{n}})
	refCompose(want, renderText(160, 40, tc, in), 8, 80, opAnd)
	s.add(typeEndOfStripe, binary.BigEndian.AppendUint32(nil, 119), segOpts{})
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "striped", page: s.buf, want: want}
}

// caseGlobals keeps a table and a symbol dictionary in the global segments.
func caseGlobals() fixture {
	var gs stream
	segs, fs, _, _, _, _, _ := customTables()
	tFS := gs.add(typeTables, segs[0], segOpts{global: true})
	g := glyphs(81, 18, 13, 16)
	sortByHeight(g)
	sd := plainSD(g, 0, nominalAT(0))
	sd.huff = true
	nSD := gs.add(typeSymbolDict, symbolDictData(sd), segOpts{global: true})

	s := &stream{next: gs.next}
	const w, h = 200, 120
	s.add(typePageInfo, pageInfo(w, h, 0, 0), segOpts{})
	want := newBitmap(w, h)
	tc := &textConfig{huff: true, syms: sd.exported(), corner: cornerTopLeft, fsSel: 3, fs: fs, dsSel: 2}
	in := lineOfText(rand.New(rand.NewSource(82)), tc.syms, 2, 2, 180, 4, 25, cornerTopLeft)
	s.add(typeImmediateText, textRegion(190, 110, 5, 5, opOr, tc, in), segOpts{refs: []uint32{nSD, tFS}})
	refCompose(want, renderText(190, 110, tc, in), 5, 5, opOr)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "globals", page: s.buf, globals: gs.buf, want: want}
}

// caseSegments exercises the segment header: 2- and 4-byte segment numbers,
// the long form of the referred-to segments, 4-byte page associations, and
// segments to ignore.
func caseSegments() fixture {
	const w, h = 200, 100
	s := &stream{next: 250}
	o := segOpts{page4: true}
	s.add(typePageInfo, pageInfo(w, h, 0, 0), o)
	s.add(typeProfiles, []byte{0, 0, 0, 1, 0, 0, 0, 1}, o)
	s.add(typeExtension, []byte{0x10, 0, 0, 0, 1, 2, 3}, o)
	want := newBitmap(w, h)
	var dicts []uint32
	var syms []*Bitmap
	// Seven referred-to segments: jbig2dec 0.20 reads ⌊(n+1)/8⌋ bytes of
	// retention flags where T.88 has ⌈(n+1)/8⌉, which agree for n = 7.
	for i := 0; i < 7; i++ {
		sd := plainSD(glyphs(int64(90+i), 3, 12+i), i%4, nominalAT(i%4))
		dicts = append(dicts, s.add(typeSymbolDict, symbolDictData(sd), o))
		syms = append(syms, sd.exported()...)
	}
	s.next = 70000
	tc := &textConfig{syms: syms, corner: cornerTopLeft}
	in := lineOfText(rand.New(rand.NewSource(97)), syms, 2, 2, 180, 4, 22, cornerTopLeft)
	o.refs = dicts
	s.add(typeImmediateText, textRegion(190, 95, 5, 3, opOr, tc, in), o)
	refCompose(want, renderText(190, 95, tc, in), 5, 3, opOr)
	s.add(typeEndOfPage, nil, segOpts{page4: true})
	s.add(typeEndOfFile, nil, segOpts{global: true})
	return fixture{name: "segments", page: s.buf, want: want}
}

// caseRetainedContexts codes a dictionary with the contexts another one retained.
func caseRetainedContexts() fixture {
	const w, h = 200, 110
	s, want := newPage(w, h, 0)
	g := glyphs(101, 24, 11, 14)
	sortByHeight(g)
	a := plainSD(g[:12], 1, []point{{-2, -2}})
	a.ctxRetain = true
	nA := s.add(typeSymbolDict, symbolDictData(a), segOpts{})
	b := plainSD(g[12:], 1, []point{{-2, -2}})
	b.ctxUsed = true
	b.gb = append([]uint8{}, a.gb...)
	b.in = a.exported()
	nB := s.add(typeSymbolDict, symbolDictData(b), segOpts{refs: []uint32{nA}})
	tc := &textConfig{syms: append(a.exported(), b.exported()...), corner: cornerTopLeft}
	in := lineOfText(rand.New(rand.NewSource(102)), tc.syms, 2, 2, 180, 4, 24, cornerTopLeft)
	s.add(typeImmediateText, textRegion(190, 100, 5, 5, opOr, tc, in), segOpts{refs: []uint32{nA, nB}})
	refCompose(want, renderText(190, 100, tc, in), 5, 5, opOr)
	s.add(typeEndOfPage, nil, segOpts{})
	return fixture{name: "retained_contexts", page: s.buf, want: want,
		reference: "jbig2dec does not implement retained bitmap coding contexts; they follow xpdf, which retains the generic and refinement contexts"}
}

// fixtures returns all test fixtures.
func fixtures() []fixture {
	return []fixture{
		genericCase("generic_t0", 0, nominalAT(0), false, false),
		genericCase("generic_t0_at", 0, []point{{4, -1}, {-4, -1}, {3, -2}, {-3, -2}}, true, false),
		genericCase("generic_t1", 1, []point{{-3, -2}}, true, false),
		genericCase("generic_t2", 2, nominalAT(2), false, false),
		genericCase("generic_t3", 3, []point{{-5, 0}}, true, false),
		genericCase("generic_mmr", 0, nil, false, true),
		caseUnknownLength(),
		caseRefinement(),
		caseTextArith(),
		caseTextStrips(),
		caseTextCorners(),
		caseCombination(),
		caseTextHuffman(),
		caseTextHuffmanCustom(),
		refaggCase("refagg_arith", refaggOpts{aggregates: true}),
		refaggCase("refagg_arith_t1", refaggOpts{aggregates: true, aggRefine: true, rtemplate: 1}),
		refaggCase("refagg_huffman", refaggOpts{huff: true, aggregates: true, aggRefine: true, rtemplate: 1,
			reference: "jbig2dec decodes Huffman coded refinements from other bytes than the BMSIZE bytes that hold them, and reads a symbol ID table for the text region of an aggregate, which T.88 codes with fixed-length IDs (6.5.8.2.3)"}),
		caseBeyondJbig2dec(),
		caseHalftone("halftone_arith", false),
		caseHalftone("halftone_mmr", true),
		caseStriped(),
		caseGlobals(),
		caseSegments(),
		caseRetainedContexts(),
	}
}

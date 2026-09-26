package jpx

import (
	"cmp"
	"errors"
	"slices"
)

// Tier-2 decoding (Annex B): packet headers, and the order of the packets
// in a tile.

// errStop ends the packets of a tile early: the data ran out, or a packet
// header does not parse. What was read so far is still decoded.
var errStop = errors.New("jpx: end of tile data")

// Progression orders.
const (
	progLRCP = iota
	progRLCP
	progRPCL
	progPCRL
	progCPRL
)

// bitReader reads packet header bits (B.10.1): after a 0xFF byte, the next
// byte holds only seven bits.
type bitReader struct {
	data []byte
	pos  int // next byte
	buf  uint32
	ct   int
}

func (b *bitReader) bit() uint32 {
	if b.ct == 0 {
		b.ct = 8
		if b.buf == 0xff {
			b.ct = 7
		}
		b.buf = 0
		if b.pos < len(b.data) {
			b.buf = uint32(b.data[b.pos])
		}
		b.pos++
	}
	b.ct--
	return b.buf >> b.ct & 1
}

func (b *bitReader) bits(n int) uint32 {
	var v uint32
	for range n {
		v = v<<1 | b.bit()
	}
	return v
}

// align skips to the end of the header; a header never ends with 0xFF, so
// the byte after one belongs to it.
func (b *bitReader) align() {
	if b.buf == 0xff {
		b.pos++
		b.buf = 0
	}
	b.ct = 0
}

// overrun reports whether the reader went past its data.
func (b *bitReader) overrun() bool { return b.pos > len(b.data) }

// tagTree is a tag tree (B.10.2) over the code-blocks of a precinct band.
type tagTree struct {
	nodes  []tagNode
	parent []int32
}

type tagNode struct {
	value int32
	low   int32
}

const tagInf = 1 << 30

func newTagTree(w, h int) tagTree {
	n := 0
	for lw, lh := w, h; ; lw, lh = (lw+1)/2, (lh+1)/2 {
		n += lw * lh
		if lw*lh <= 1 {
			break
		}
	}
	t := tagTree{nodes: make([]tagNode, n), parent: make([]int32, n)}
	for i := range t.nodes {
		t.nodes[i].value = tagInf
	}
	base := 0
	for lw, lh := w, h; lw*lh > 1; {
		pw, ph := (lw+1)/2, (lh+1)/2
		next := base + lw*lh
		for y := range lh {
			for x := range lw {
				t.parent[base+y*lw+x] = int32(next + (y/2)*pw + x/2)
			}
		}
		base, lw, lh = next, pw, ph
	}
	t.parent[n-1] = -1
	return t
}

// decode reports whether the value of leaf is below threshold, reading as
// many bits as needed.
func (t *tagTree) decode(br *bitReader, leaf int, threshold int32) bool {
	var stack [40]int32
	sp := 0
	i := int32(leaf)
	for t.parent[i] >= 0 && sp < len(stack) {
		stack[sp] = i
		sp++
		i = t.parent[i]
	}
	var low int32
	for {
		n := &t.nodes[i]
		if low > n.low {
			n.low = low
		} else {
			low = n.low
		}
		for low < threshold && low < n.value {
			if br.bit() != 0 {
				n.value = low
			} else {
				low++
			}
		}
		n.low = low
		if sp == 0 {
			return n.value < threshold
		}
		sp--
		i = stack[sp]
	}
}

// value decodes the full value of leaf, up to limit.
func (t *tagTree) value(br *bitReader, leaf int, limit int32) int32 {
	var th int32 = 1
	for !t.decode(br, leaf, th) {
		if th >= limit || br.overrun() {
			return limit
		}
		th++
	}
	return t.nodes[leaf].value
}

// maxSegPasses is the number of coding passes that segment k of a
// code-block can hold under the code-block style.
func maxSegPasses(style byte, k int) int {
	switch {
	case style&cbTermAll != 0:
		return 1
	case style&cbBypass != 0:
		if k == 0 {
			return 10
		}
		if k%2 == 1 {
			return 2 // significance and refinement passes, raw
		}
		return 1 // cleanup pass, arithmetic coded
	}
	return 1 << 30
}

// pendingData is the number of bytes a code-block has in a packet body.
type pendingData struct {
	cb *codeBlock
	n  int
}

// packet decodes the packet of a precinct for a layer and appends the
// code-block data it carries.
func (td *tileDecoder) packet(tc *tileComp, p *precinct, layer int) error {
	body := td.body
	start := td.bpos
	if td.sop && td.bpos+6 <= len(body) && body[td.bpos] == 0xff && body[td.bpos+1] == 0x91 {
		td.bpos += 6
	}
	hdrStart := td.bpos
	br := bitReader{data: body, pos: td.bpos}
	if td.hdr != nil {
		br = bitReader{data: td.hdr, pos: td.hpos}
	}
	if br.pos >= len(br.data) {
		return errStop
	}
	pend := td.pending[:0]
	if br.bit() != 0 {
		for bi := range p.bands {
			pb := &p.bands[bi]
			for i := range pb.blocks {
				cb := &pb.blocks[i]
				var in bool
				if cb.segs == nil {
					in = pb.incl.decode(&br, i, int32(layer)+1)
				} else {
					in = br.bit() != 0
				}
				if !in {
					continue
				}
				if cb.segs == nil {
					cb.zbp = int(pb.zbp.value(&br, i, 1024))
					cb.lblock = 3
					cb.segs = make([]segment, 0, 1)
				}
				n := readPassCount(&br)
				for br.bit() != 0 {
					if cb.lblock++; cb.lblock > 32 {
						return errStop
					}
				}
				total := 0
				for n > 0 {
					k := len(cb.segs) - 1
					if k < 0 || cb.segs[k].passes >= maxSegPasses(tc.style, k) {
						if len(cb.segs) >= 256 {
							return errStop
						}
						cb.segs = append(cb.segs, segment{})
						k++
					}
					s := &cb.segs[k]
					take := min(maxSegPasses(tc.style, k)-s.passes, n)
					nbits := cb.lblock + floorLog2(take)
					if nbits > 32 {
						return errStop
					}
					l := min(int(br.bits(nbits)), maxSegLen)
					s.passes += take
					s.length = min(s.length+l, maxSegLen)
					total = min(total+l, maxSegLen)
					n -= take
				}
				pend = append(pend, pendingData{cb, total})
				if br.overrun() {
					return errStop
				}
			}
		}
	}
	br.align()
	if br.overrun() {
		return errStop
	}
	hp := br.pos
	if td.eph && hp+2 <= len(br.data) && br.data[hp] == 0xff && br.data[hp+1] == 0x92 {
		hp += 2
	}
	if td.hdr != nil {
		td.hpos = hp
	} else {
		td.bpos = hp
	}
	stop := false
	for _, pd := range pend {
		n := pd.n
		if avail := len(body) - td.bpos; n > avail {
			n, stop = avail, true
		}
		chunk := body[td.bpos : td.bpos+n : td.bpos+n]
		if pd.cb.data == nil {
			pd.cb.data = chunk // no copy for a first contribution
		} else {
			pd.cb.data = append(pd.cb.data, chunk...)
		}
		td.bpos += n
	}
	td.pending = pend[:0]
	if stop {
		return errStop
	}
	if td.c.packetHook != nil {
		td.c.packetHook(td.index, start, hdrStart, hp, td.bpos)
	}
	return nil
}

// maxSegLen bounds recorded segment lengths; longer ones exceed any data.
const maxSegLen = 1 << 30

// readPassCount decodes the number of new coding passes (Table B.4).
func readPassCount(br *bitReader) int {
	if br.bit() == 0 {
		return 1
	}
	if br.bit() == 0 {
		return 2
	}
	if v := br.bits(2); v != 3 {
		return 3 + int(v)
	}
	if v := br.bits(5); v != 31 {
		return 6 + int(v)
	}
	return 37 + int(br.bits(7))
}

func floorLog2(n int) int {
	l := 0
	for n > 1 {
		n >>= 1
		l++
	}
	return l
}

// maxIdle bounds, over a codestream, the steps of the packet iteration
// that read nothing: packets that progression order changes revisit, and
// resolution levels without precincts.
const maxIdle = 1 << 26

// idle counts n such steps and reports whether the budget is spent.
func (td *tileDecoder) idle(n int) bool {
	td.c.idle += n
	return td.c.idle > maxIdle
}

// visit decodes the packet of precinct p for layer l unless an earlier
// progression already did.
func (td *tileDecoder) visit(tc *tileComp, p *precinct, l int) error {
	if l < p.done {
		if td.idle(1) {
			return errStop
		}
		return nil
	}
	p.done = l + 1
	return td.packet(tc, p, l)
}

// readPackets decodes the packets of the tile in progression order.
func (td *tileDecoder) readPackets() error {
	for _, pc := range td.pocs {
		lye := min(pc.lye, td.layers)
		ce := min(pc.ce, len(td.comps))
		re := pc.re
		var err error
		switch pc.prog {
		case progLRCP:
			for l := 0; l < lye && err == nil; l++ {
				for r := pc.rs; r < re && err == nil; r++ {
					for c := pc.cs; c < ce && err == nil; c++ {
						err = td.visitRes(c, r, l)
					}
				}
			}
		case progRLCP:
			for r := pc.rs; r < re && err == nil; r++ {
				for l := 0; l < lye && err == nil; l++ {
					for c := pc.cs; c < ce && err == nil; c++ {
						err = td.visitRes(c, r, l)
					}
				}
			}
		default:
			err = td.visitPositions(pc, lye, ce)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (td *tileDecoder) visitRes(c, r, l int) error {
	tc := &td.comps[c]
	if td.idle(1) {
		return errStop
	}
	if r >= len(tc.res) {
		return nil
	}
	res := &tc.res[r]
	for i := range res.precincts {
		if err := td.visit(tc, &res.precincts[i], l); err != nil {
			return err
		}
	}
	return nil
}

// precinctPos is a precinct with the reference grid position at which the
// position-driven progressions (B.12.1.3-5) reach it.
type precinctPos struct {
	x, y    int64
	c, r, p int
}

func (td *tileDecoder) visitPositions(pc poc, lye, ce int) error {
	var list []precinctPos
	for c := pc.cs; c < ce; c++ {
		tc := &td.comps[c]
		for r := pc.rs; r < pc.re && r < len(tc.res); r++ {
			res := &tc.res[r]
			if td.idle(1 + len(res.precincts)) {
				return errStop
			}
			if len(res.precincts) == 0 {
				continue
			}
			lev := len(tc.res) - 1 - r
			dx := int64(td.dx[c]) << lev
			dy := int64(td.dy[c]) << lev
			for py := range res.npy {
				y := int64((res.y0>>res.ppy)+py) << res.ppy * dy
				if py == 0 && res.y0&(1<<res.ppy-1) != 0 {
					y = int64(td.y0)
				}
				for px := range res.npx {
					x := int64((res.x0>>res.ppx)+px) << res.ppx * dx
					if px == 0 && res.x0&(1<<res.ppx-1) != 0 {
						x = int64(td.x0)
					}
					list = append(list, precinctPos{x, y, c, r, py*res.npx + px})
				}
			}
		}
	}
	var order func(a, b precinctPos) int
	switch pc.prog {
	case progRPCL:
		order = func(a, b precinctPos) int {
			return cmpChain(cmp.Compare(a.r, b.r), cmp.Compare(a.y, b.y), cmp.Compare(a.x, b.x), cmp.Compare(a.c, b.c))
		}
	case progPCRL:
		order = func(a, b precinctPos) int {
			return cmpChain(cmp.Compare(a.y, b.y), cmp.Compare(a.x, b.x), cmp.Compare(a.c, b.c), cmp.Compare(a.r, b.r))
		}
	default:
		order = func(a, b precinctPos) int {
			return cmpChain(cmp.Compare(a.c, b.c), cmp.Compare(a.y, b.y), cmp.Compare(a.x, b.x), cmp.Compare(a.r, b.r))
		}
	}
	slices.SortFunc(list, order)
	for _, e := range list {
		tc := &td.comps[e.c]
		p := &tc.res[e.r].precincts[e.p]
		for l := p.done; l < lye; l++ {
			if err := td.visit(tc, p, l); err != nil {
				return err
			}
		}
	}
	return nil
}

func cmpChain(c ...int) int {
	for _, v := range c {
		if v != 0 {
			return v
		}
	}
	return 0
}

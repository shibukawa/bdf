package jpx

import "slices"

// Marker codes (Table A.2).
const (
	mSOC = 0xff4f
	mSIZ = 0xff51
	mCOD = 0xff52
	mCOC = 0xff53
	mQCD = 0xff5c
	mQCC = 0xff5d
	mRGN = 0xff5e
	mPOC = 0xff5f
	mPPM = 0xff60
	mPPT = 0xff61
	mSOT = 0xff90
	mSOD = 0xff93
	mEOC = 0xffd9
)

// Limits against malformed or hostile input.
const (
	maxGrid      = 1 << 30 // reference grid coordinates and tile sizes
	maxSamples   = 1 << 28 // samples of all components together
	maxTiles     = 65535
	maxComps     = 16384
	maxPrecision = 30
	maxBlocks    = 1 << 22 // code-blocks in a tile
	maxPrecincts = 1 << 22 // precincts in a tile
)

type component struct {
	prec   int
	signed bool
	dx, dy int
}

// codingStyle holds the parameters of COD or COC that apply to one
// component (SPcod, SPcoc).
type codingStyle struct {
	levels     int // number of decomposition levels
	cbw, cbh   int // code-block size exponents
	style      byte
	reversible bool    // 5/3 filter
	ppx, ppy   []uint8 // precinct size exponents per resolution; nil for 15
}

// cod is a COD marker segment.
type cod struct {
	sop, eph bool
	prog     int
	layers   int
	mct      int
	cs       codingStyle
}

type qstep struct{ exp, mant int }

// quant is a QCD or QCC marker segment.
type quant struct {
	style int // 0 no quantization, 1 scalar derived, 2 scalar expounded
	guard int
	steps []qstep
}

// step returns the exponent and mantissa of band i (0 for LL, then three per
// resolution level) of resolution r.
func (q *quant) step(i, r int) qstep {
	if q.style == 1 {
		s := q.steps[0]
		if r > 0 {
			s.exp = max(0, s.exp-r+1)
		}
		return s
	}
	return q.steps[min(i, len(q.steps)-1)]
}

// poc is a progression order change: packets with resolutions in
// [rs, re), components in [cs, ce) and layers below lye, in order prog.
type poc struct {
	rs, cs, lye, re, ce, prog int
}

// header gathers the marker segments of the main header, or of the tile-part
// headers of a tile.
type header struct {
	cod  *cod
	coc  []*codingStyle
	qcd  *quant
	qcc  []*quant
	rgn  []int // ROI shift per component, -1 when not given
	pocs []poc
}

// tile collects the tile-parts of a tile.
type tile struct {
	hdr   header
	parts [][]byte // tile-part bodies
	ppt   []zseg
	ppm   [][]byte // packet headers of the tile-parts, from PPM
}

type zseg struct {
	z    int
	data []byte
}

// codestream is a parsed codestream: the main header and the tile-parts.
type codestream struct {
	x0, y0, x1, y1   int // image area on the reference grid
	tx0, ty0, tw, th int // tile grid
	ntx, nty         int
	comps            []component
	main             header
	ppm              bool // packet headers are in PPM marker segments
	tiles            []*tile
	idle             int // see maxIdle

	// packetHook, when set, is told where each packet of a tile lies in
	// its packet data: the packet starts at start, its header spans
	// [hdr, hdrEnd) unless headers are packed elsewhere, and it ends at end.
	packetHook func(tile, start, hdr, hdrEnd, end int)
}

type byteReader struct {
	b   []byte
	err bool
}

func (r *byteReader) u8() int {
	if len(r.b) < 1 {
		r.err = true
		return 0
	}
	v := r.b[0]
	r.b = r.b[1:]
	return int(v)
}

func (r *byteReader) u16() int {
	if len(r.b) < 2 {
		r.err = true
		return 0
	}
	v := int(r.b[0])<<8 | int(r.b[1])
	r.b = r.b[2:]
	return v
}

func (r *byteReader) u32() uint32 {
	if len(r.b) < 4 {
		r.err = true
		return 0
	}
	v := uint32(r.b[0])<<24 | uint32(r.b[1])<<16 | uint32(r.b[2])<<8 | uint32(r.b[3])
	r.b = r.b[4:]
	return v
}

func be16(b []byte) int { return int(b[0])<<8 | int(b[1]) }

func be32(b []byte) uint32 {
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// parseCodestream reads the main header and splits the tile-parts. A
// codestream truncated after the main header keeps the tile-parts it has.
func parseCodestream(data []byte) (*codestream, error) {
	if len(data) < 4 || be16(data) != mSOC || be16(data[2:]) != mSIZ {
		return nil, errorf("not a JPEG 2000 codestream")
	}
	c := &codestream{}
	var ppm []zseg
	pos := 2
	for {
		if pos+4 > len(data) {
			return nil, errorf("truncated main header")
		}
		m := be16(data[pos:])
		if m == mSOT || m == mEOC {
			break
		}
		if m < 0xff00 {
			return nil, errorf("invalid marker %04X in main header", m)
		}
		l := be16(data[pos+2:])
		if l < 2 || pos+2+l > len(data) {
			return nil, errorf("truncated marker segment %04X", m)
		}
		seg := data[pos+4 : pos+2+l]
		pos += 2 + l
		var err error
		switch m {
		case mSIZ:
			if c.comps != nil {
				return nil, errorf("duplicate SIZ marker")
			}
			err = c.parseSIZ(seg)
		case mPPM:
			if len(seg) < 1 {
				return nil, errorf("invalid PPM marker")
			}
			ppm = append(ppm, zseg{int(seg[0]), seg[1:]})
		default:
			err = c.parseCommon(&c.main, m, seg)
		}
		if err != nil {
			return nil, err
		}
	}
	if c.main.cod == nil {
		return nil, errorf("missing COD marker")
	}
	if c.main.qcd == nil {
		return nil, errorf("missing QCD marker")
	}
	c.tiles = make([]*tile, c.ntx*c.nty)
	var ppmData []byte
	if ppm != nil {
		c.ppm = true
		slices.SortStableFunc(ppm, func(a, b zseg) int { return a.z - b.z })
		for _, s := range ppm {
			ppmData = append(ppmData, s.data...)
		}
	}
	// Tile-parts. Anything unexpected ends the codestream.
	for pos+12 <= len(data) && be16(data[pos:]) == mSOT {
		r := byteReader{b: data[pos+2 : pos+12]}
		if r.u16() != 10 {
			return nil, errorf("invalid SOT marker")
		}
		isot, psot := r.u16(), int64(r.u32())
		if isot >= len(c.tiles) {
			return nil, errorf("invalid tile index %d", isot)
		}
		end := int64(len(data))
		if psot != 0 {
			if psot < 14 {
				return nil, errorf("invalid tile-part length")
			}
			end = min(end, int64(pos)+psot)
		} else if be16(data[len(data)-2:]) == mEOC {
			end -= 2
		}
		t := c.tiles[isot]
		if t == nil {
			t = &tile{}
			c.tiles[isot] = t
		}
		p := pos + 12
		sod := false
		for p+2 <= int(end) {
			m := be16(data[p:])
			if m == mSOD {
				p += 2
				sod = true
				break
			}
			if m < 0xff00 || p+4 > int(end) {
				break
			}
			l := be16(data[p+2:])
			if l < 2 || p+2+l > int(end) {
				break
			}
			seg := data[p+4 : p+2+l]
			p += 2 + l
			var err error
			if m == mPPT {
				if len(seg) < 1 {
					return nil, errorf("invalid PPT marker")
				}
				t.ppt = append(t.ppt, zseg{int(seg[0]), seg[1:]})
			} else {
				err = c.parseCommon(&t.hdr, m, seg)
			}
			if err != nil {
				return nil, err
			}
		}
		if !sod {
			break // truncated tile-part header
		}
		t.parts = append(t.parts, data[p:end])
		if c.ppm && len(ppmData) >= 4 {
			n := int(min(be32(ppmData), uint32(len(ppmData)-4)))
			t.ppm = append(t.ppm, ppmData[4:4+n])
			ppmData = ppmData[4+n:]
		}
		pos = int(end)
	}
	return c, nil
}

func (c *codestream) parseSIZ(b []byte) error {
	r := byteReader{b: b}
	r.u16() // capabilities
	x1, y1 := r.u32(), r.u32()
	x0, y0 := r.u32(), r.u32()
	tw, th := r.u32(), r.u32()
	tx0, ty0 := r.u32(), r.u32()
	n := r.u16()
	if r.err || len(r.b) < 3*n {
		return errorf("truncated SIZ marker")
	}
	if n < 1 || n > maxComps {
		return errorf("invalid number of components %d", n)
	}
	for _, v := range []uint32{x1, y1, tw, th} {
		if v > maxGrid {
			return errorf("image too large")
		}
	}
	if x0 >= x1 || y0 >= y1 || tw == 0 || th == 0 || tx0 > x0 || ty0 > y0 ||
		tx0+tw <= x0 || ty0+th <= y0 {
		return errorf("invalid image or tile size")
	}
	c.x0, c.y0, c.x1, c.y1 = int(x0), int(y0), int(x1), int(y1)
	c.tx0, c.ty0, c.tw, c.th = int(tx0), int(ty0), int(tw), int(th)
	c.ntx = ceilDiv(c.x1-c.tx0, c.tw)
	c.nty = ceilDiv(c.y1-c.ty0, c.th)
	if int64(c.ntx)*int64(c.nty) > maxTiles {
		return errorf("too many tiles")
	}
	c.comps = make([]component, n)
	var total int64
	for i := range c.comps {
		s, dx, dy := r.u8(), r.u8(), r.u8()
		if dx == 0 || dy == 0 {
			return errorf("invalid component subsampling")
		}
		cp := component{prec: s&0x7f + 1, signed: s&0x80 != 0, dx: dx, dy: dy}
		if cp.prec > maxPrecision {
			return errorf("unsupported component precision %d", cp.prec)
		}
		c.comps[i] = cp
		total += int64(ceilDiv(c.x1, dx)-ceilDiv(c.x0, dx)) * int64(ceilDiv(c.y1, dy)-ceilDiv(c.y0, dy))
	}
	if total > maxSamples {
		return errorf("image too large")
	}
	return nil
}

// parseCommon parses the marker segments allowed in both the main header
// and tile-part headers. Unknown segments are skipped.
func (c *codestream) parseCommon(h *header, m int, b []byte) error {
	if c.comps == nil {
		return errorf("missing SIZ marker")
	}
	r := byteReader{b: b}
	comp := func() int {
		if len(c.comps) < 257 {
			return r.u8()
		}
		return r.u16()
	}
	switch m {
	case mCOD:
		scod := r.u8()
		d := &cod{sop: scod&2 != 0, eph: scod&4 != 0}
		d.prog, d.layers, d.mct = r.u8(), r.u16(), r.u8()
		if r.err {
			return errorf("truncated COD marker")
		}
		if d.prog > progCPRL || d.layers == 0 {
			return errorf("invalid COD marker")
		}
		if d.mct > 1 {
			return errorf("unsupported multiple component transform %d", d.mct)
		}
		if err := parseSPcod(&r, &d.cs, scod&1 != 0); err != nil {
			return err
		}
		h.cod = d
	case mCOC:
		i := comp()
		scoc := r.u8()
		if r.err || i >= len(c.comps) {
			return errorf("invalid COC marker")
		}
		cs := &codingStyle{}
		if err := parseSPcod(&r, cs, scoc&1 != 0); err != nil {
			return err
		}
		if h.coc == nil {
			h.coc = make([]*codingStyle, len(c.comps))
		}
		h.coc[i] = cs
	case mQCD:
		q, err := parseQuant(r.b)
		if err != nil {
			return err
		}
		h.qcd = q
	case mQCC:
		i := comp()
		if r.err || i >= len(c.comps) {
			return errorf("invalid QCC marker")
		}
		q, err := parseQuant(r.b)
		if err != nil {
			return err
		}
		if h.qcc == nil {
			h.qcc = make([]*quant, len(c.comps))
		}
		h.qcc[i] = q
	case mRGN:
		i := comp()
		style, shift := r.u8(), r.u8()
		if r.err || i >= len(c.comps) || style != 0 {
			return errorf("invalid RGN marker")
		}
		if h.rgn == nil {
			h.rgn = make([]int, len(c.comps))
			for k := range h.rgn {
				h.rgn[k] = -1
			}
		}
		h.rgn[i] = shift
	case mPOC:
		size := 7
		if len(c.comps) >= 257 {
			size = 9
		}
		if len(b) < size || len(b)%size != 0 {
			return errorf("invalid POC marker")
		}
		for len(r.b) >= size {
			var p poc
			p.rs, p.cs = r.u8(), comp()
			p.lye, p.re, p.ce = r.u16(), r.u8(), comp()
			p.prog = r.u8()
			if p.ce == 0 {
				p.ce = 256
				if size == 9 {
					p.ce = 16384
				}
			}
			if p.prog > progCPRL {
				return errorf("invalid POC marker")
			}
			h.pocs = append(h.pocs, p)
		}
	}
	return nil
}

// parseSPcod parses the coding style parameters shared by COD and COC.
func parseSPcod(r *byteReader, cs *codingStyle, precincts bool) error {
	levels, cbw, cbh, style, xform := r.u8(), r.u8(), r.u8(), r.u8(), r.u8()
	if r.err {
		return errorf("truncated coding style")
	}
	if levels > 32 || cbw > 8 || cbh > 8 || cbw+cbh > 8 {
		return errorf("invalid coding style")
	}
	if style&0xc0 != 0 {
		return errorf("unsupported code-block style %#x", style)
	}
	if xform > 1 {
		return errorf("unsupported wavelet transform %d", xform)
	}
	*cs = codingStyle{levels: levels, cbw: cbw + 2, cbh: cbh + 2, style: byte(style), reversible: xform == 1}
	if precincts {
		if len(r.b) < levels+1 {
			return errorf("truncated precinct sizes")
		}
		cs.ppx = make([]uint8, levels+1)
		cs.ppy = make([]uint8, levels+1)
		for i := range levels + 1 {
			v := r.b[i]
			cs.ppx[i], cs.ppy[i] = v&0xf, v>>4
			if i > 0 && (cs.ppx[i] == 0 || cs.ppy[i] == 0) {
				return errorf("invalid precinct size")
			}
		}
	}
	return nil
}

func parseQuant(b []byte) (*quant, error) {
	if len(b) < 1 {
		return nil, errorf("invalid quantization marker")
	}
	q := &quant{style: int(b[0] & 0x1f), guard: int(b[0] >> 5)}
	b = b[1:]
	switch q.style {
	case 0:
		for _, v := range b {
			q.steps = append(q.steps, qstep{exp: int(v >> 3)})
		}
	case 1, 2:
		for ; len(b) >= 2; b = b[2:] {
			v := be16(b)
			q.steps = append(q.steps, qstep{exp: v >> 11, mant: v & 0x7ff})
			if q.style == 1 {
				break
			}
		}
	default:
		return nil, errorf("unsupported quantization style %d", q.style)
	}
	if len(q.steps) == 0 {
		return nil, errorf("invalid quantization marker")
	}
	return q, nil
}

func ceilDiv(a, b int) int { return (a + b - 1) / b }

// ceilDivPow2 is ceil(a / 2^n) for 0 <= a < 2^30.
func ceilDivPow2(a, n int) int {
	if n >= 30 {
		return b2i(a > 0)
	}
	return (a + 1<<n - 1) >> n
}

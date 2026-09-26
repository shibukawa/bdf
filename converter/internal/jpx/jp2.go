package jpx

import "bytes"

// JP2 and JPX file format boxes (ITU-T T.800 Annex I, T.801 Annex M).

var jp2Signature = []byte{0, 0, 0, 12, 'j', 'P', ' ', ' ', 0x0d, 0x0a, 0x87, 0x0a}

const (
	boxJP2H = 0x6a703268 // jp2h
	boxJP2C = 0x6a703263 // jp2c
	boxJPCH = 0x6a706368 // jpch, JPX codestream header
	boxJPLH = 0x6a706c68 // jplh, JPX compositing layer header
	boxCOLR = 0x636f6c72 // colr
	boxPCLR = 0x70636c72 // pclr
	boxCMAP = 0x636d6170 // cmap
	boxCDEF = 0x63646566 // cdef
)

// jp2Header is what the JP2 header boxes say about the codestream.
type jp2Header struct {
	colr       bool
	colorSpace int
	icc        []byte
	pal        *palette
	cmap       []cmapEntry
	cdef       []cdefEntry
}

// palette is a pclr box: entries rows of len(cols) values.
type palette struct {
	entries int
	cols    []component
	values  []int32
}

type cmapEntry struct{ comp, typ, col int }

type cdefEntry struct{ ch, typ, asoc int }

// box returns the type and content of the box at the start of b, and the
// rest of b. A box that claims more than there is is cut short.
func box(b []byte) (typ uint32, content, rest []byte, err error) {
	if len(b) < 8 {
		return 0, nil, nil, errorf("truncated box")
	}
	l, typ, hl := uint64(be32(b)), be32(b[4:]), uint64(8)
	switch l {
	case 0:
		l = uint64(len(b))
	case 1:
		if len(b) < 16 {
			return 0, nil, nil, errorf("truncated box")
		}
		l, hl = uint64(be32(b[8:]))<<32|uint64(be32(b[12:])), 16
	}
	if l < hl {
		return 0, nil, nil, errorf("invalid box length")
	}
	l = min(l, uint64(len(b)))
	return typ, b[hl:l], b[l:], nil
}

// parseJP2 returns the first codestream of a JP2 or JPX file and its
// header.
func parseJP2(data []byte) (*jp2Header, []byte, error) {
	h := &jp2Header{}
	var cs []byte
	for b := data; len(b) > 0; {
		typ, content, rest, err := box(b)
		if err != nil {
			if cs != nil {
				break
			}
			return nil, nil, err
		}
		b = rest
		switch typ {
		case boxJP2H, boxJPCH, boxJPLH:
			if err := h.parse(content); err != nil {
				return nil, nil, err
			}
		case boxJP2C:
			if cs == nil {
				cs = content
			}
		}
	}
	if cs == nil {
		return nil, nil, errorf("no codestream in JP2 file")
	}
	return h, cs, nil
}

// parse reads the boxes of a header superbox; values already found in an
// earlier header are kept.
func (h *jp2Header) parse(b []byte) error {
	for len(b) > 0 {
		typ, c, rest, err := box(b)
		if err != nil {
			return nil // ignore trailing garbage in a header
		}
		b = rest
		switch typ {
		case boxCOLR:
			if h.colr || len(c) < 3 {
				continue
			}
			switch c[0] {
			case 1:
				if len(c) >= 7 {
					h.colr = true
					h.colorSpace = int(be32(c[3:]))
				}
			case 2, 3:
				h.colr = true
				h.icc = bytes.Clone(c[3:])
			}
		case boxPCLR:
			if h.pal == nil {
				p, err := parsePalette(c)
				if err != nil {
					return err
				}
				h.pal = p
			}
		case boxCMAP:
			if h.cmap == nil {
				for ; len(c) >= 4; c = c[4:] {
					h.cmap = append(h.cmap, cmapEntry{be16(c), int(c[2]), int(c[3])})
				}
			}
		case boxCDEF:
			if h.cdef == nil && len(c) >= 2 {
				n := be16(c)
				c = c[2:]
				for ; n > 0 && len(c) >= 6; n, c = n-1, c[6:] {
					h.cdef = append(h.cdef, cdefEntry{be16(c), be16(c[2:]), be16(c[4:])})
				}
			}
		}
	}
	return nil
}

func parsePalette(b []byte) (*palette, error) {
	if len(b) < 3 {
		return nil, errorf("invalid pclr box")
	}
	p := &palette{entries: be16(b)}
	n := int(b[2])
	if p.entries == 0 || p.entries > 1024 || n == 0 || len(b) < 3+n {
		return nil, errorf("invalid pclr box")
	}
	size := 0
	for _, v := range b[3 : 3+n] {
		c := component{prec: int(v&0x7f) + 1, signed: v&0x80 != 0}
		if c.prec > maxPrecision {
			return nil, errorf("unsupported palette precision %d", c.prec)
		}
		p.cols = append(p.cols, c)
		size += (c.prec + 7) / 8
	}
	b = b[3+n:]
	if len(b) < p.entries*size {
		return nil, errorf("truncated pclr box")
	}
	p.values = make([]int32, p.entries*n)
	for i := range p.values {
		c := p.cols[i%n]
		var v uint32
		for range (c.prec + 7) / 8 {
			v = v<<8 | uint32(b[0])
			b = b[1:]
		}
		v &= 1<<c.prec - 1
		if c.signed && v&(1<<(c.prec-1)) != 0 {
			p.values[i] = int32(v) - 1<<c.prec
		} else {
			p.values[i] = int32(v)
		}
	}
	return p, nil
}

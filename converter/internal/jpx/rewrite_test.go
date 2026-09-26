package jpx

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Codestream rewriting for tests: moving packet headers into PPM or PPT
// marker segments and moving marker segments between headers exercise
// syntax that OpenJPEG's encoder does not write.

type rawSeg struct {
	m int
	b []byte
}

type rawPart struct {
	isot, tpsot, tnsot int
	hdr                []rawSeg
	body               []byte
}

type rawStream struct {
	main  []rawSeg // after SIZ
	siz   []byte
	parts []rawPart
}

func splitStream(data []byte) (rawStream, error) {
	var s rawStream
	pos := 2
	for pos+4 <= len(data) {
		m := be16(data[pos:])
		if m == mSOT {
			break
		}
		l := be16(data[pos+2:])
		seg := data[pos+4 : pos+2+l]
		if m == mSIZ {
			s.siz = seg
		} else {
			s.main = append(s.main, rawSeg{m, seg})
		}
		pos += 2 + l
	}
	for pos+12 <= len(data) && be16(data[pos:]) == mSOT {
		p := rawPart{isot: be16(data[pos+4:]), tpsot: int(data[pos+10]), tnsot: int(data[pos+11])}
		end := pos + int(be32(data[pos+6:]))
		if be32(data[pos+6:]) == 0 {
			end = len(data) - 2
		}
		q := pos + 12
		for be16(data[q:]) != mSOD {
			l := be16(data[q+2:])
			p.hdr = append(p.hdr, rawSeg{be16(data[q:]), data[q+4 : q+2+l]})
			q += 2 + l
		}
		p.body = data[q+2 : end]
		s.parts = append(s.parts, p)
		pos = end
	}
	if s.siz == nil || len(s.parts) == 0 {
		return s, fmt.Errorf("cannot split codestream")
	}
	return s, nil
}

func (s rawStream) bytes() []byte {
	var b bytes.Buffer
	seg := func(m int, d []byte) {
		binary.Write(&b, binary.BigEndian, [2]uint16{uint16(m), uint16(len(d) + 2)})
		b.Write(d)
	}
	b.Write([]byte{0xff, 0x4f})
	seg(mSIZ, s.siz)
	for _, m := range s.main {
		seg(m.m, m.b)
	}
	for _, p := range s.parts {
		var h bytes.Buffer
		for _, m := range p.hdr {
			binary.Write(&h, binary.BigEndian, [2]uint16{uint16(m.m), uint16(len(m.b) + 2)})
			h.Write(m.b)
		}
		psot := 12 + h.Len() + 2 + len(p.body)
		binary.Write(&b, binary.BigEndian, struct {
			M, L, I uint16
			P       uint32
			TP, TN  uint8
		}{mSOT, 10, uint16(p.isot), uint32(psot), uint8(p.tpsot), uint8(p.tnsot)})
		b.Write(h.Bytes())
		b.Write([]byte{0xff, 0x93})
		b.Write(p.body)
	}
	b.Write([]byte{0xff, 0xd9})
	return b.Bytes()
}

// without returns the segments other than those with marker m.
func without(segs []rawSeg, m int) []rawSeg {
	var out []rawSeg
	for _, s := range segs {
		if s.m != m {
			out = append(out, s)
		}
	}
	return out
}

func find(segs []rawSeg, m int) []byte {
	for _, s := range segs {
		if s.m == m {
			return s.b
		}
	}
	return nil
}

// chunks splits data into marker segments with an index byte first.
func chunks(m int, data []byte, size int) []rawSeg {
	var out []rawSeg
	for z := 0; len(data) > 0 || z == 0; z++ {
		n := min(size, len(data))
		out = append(out, rawSeg{m, append([]byte{byte(z)}, data[:n]...)})
		data = data[n:]
	}
	return out
}

// packHeaders moves the packet headers of a codestream into PPM (in the
// main header) or PPT (in the tile-part headers) marker segments. Small
// segments make the packet headers of a tile-part span several of them.
func packHeaders(data []byte, ppm bool, segSize int) ([]byte, error) {
	c, err := parseCodestream(data)
	if err != nil {
		return nil, err
	}
	type pkt struct{ start, hdr, hdrEnd, end int }
	pkts := map[int][]pkt{}
	c.packetHook = func(t, start, hdr, hdrEnd, end int) {
		pkts[t] = append(pkts[t], pkt{start, hdr, hdrEnd, end})
	}
	if _, err := c.decode(); err != nil {
		return nil, err
	}
	s, err := splitStream(data)
	if err != nil {
		return nil, err
	}
	var ppmSegs []rawSeg
	offset := map[int]int{}
	zppt := map[int]int{}
	for i := range s.parts {
		p := &s.parts[i]
		off := offset[p.isot]
		offset[p.isot] += len(p.body)
		var hdrs, body []byte
		last := 0
		for _, k := range pkts[p.isot] {
			if k.start < off || k.start >= off+len(p.body) {
				continue
			}
			b := p.body
			hdrs = append(hdrs, b[k.hdr-off:k.hdrEnd-off]...)
			body = append(body, b[k.start-off:k.hdr-off]...)
			body = append(body, b[k.hdrEnd-off:k.end-off]...)
			last = k.end - off
		}
		body = append(body, p.body[last:]...)
		p.body = body
		if ppm {
			// Each tile-part starts a segment; its headers may continue
			// in the next ones.
			d := binary.BigEndian.AppendUint32(nil, uint32(len(hdrs)))
			for _, sg := range chunks(mPPM, append(d, hdrs...), segSize) {
				sg.b[0] = byte(len(ppmSegs))
				ppmSegs = append(ppmSegs, sg)
			}
			continue
		}
		for _, sg := range chunks(mPPT, hdrs, segSize) {
			sg.b[0] = byte(zppt[p.isot])
			zppt[p.isot]++
			p.hdr = append(p.hdr, sg)
		}
	}
	s.main = append(s.main, ppmSegs...)
	return s.bytes(), nil
}

// spcod returns the coding style part of a COD segment, with the precinct
// sizes, and the Scod flags.
func spcod(cod []byte) (scod byte, sp []byte) {
	return cod[0], cod[5:]
}

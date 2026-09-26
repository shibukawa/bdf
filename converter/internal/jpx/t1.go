package jpx

// Tier-1 decoding (Annex D): the three coding passes of every bit-plane of a
// code-block.
//
// Magnitudes are kept at twice their scale so that the reconstruction point
// half-way between two quantization levels stays representable: a
// coefficient that becomes significant in bit-plane p is set to 3<<p (the
// bit and one half below it), and each refinement bit moves it up or down by
// the new half.

// Per-sample state flags. The low byte holds the significance of the eight
// neighbours, bits 8-11 the signs of the four direct neighbours.
const (
	fNE uint32 = 1 << iota
	fSE
	fSW
	fNW
	fN
	fE
	fS
	fW
	fNegN
	fNegE
	fNegS
	fNegW
	fSig   // significant
	fRef   // refined at least once
	fVisit // coded in the significance pass of the current bit-plane
	fNeg   // negative

	fNbr = fNE | fSE | fSW | fNW | fN | fE | fS | fW
	// fCausal clears the samples of the next stripe for vertically causal
	// context formation.
	fCausal = ^(fS | fSE | fSW | fNegS)
)

// Code-block style flags (Table A.19).
const (
	cbBypass  = 0x01
	cbReset   = 0x02
	cbTermAll = 0x04
	cbCausal  = 0x08
	cbPredict = 0x10
	cbSegSym  = 0x20
)

// Band orientations for the significance contexts (Table D.1).
const (
	orientLL = iota // LL and LH: horizontal neighbours weigh most
	orientHL        // vertical neighbours weigh most
	orientHH        // diagonal neighbours weigh most
)

var (
	zcLUT = buildZC()
	scLUT = buildSC()
)

// buildZC indexes the significance context by the neighbour byte of the
// flags.
func buildZC() (t [3][256]uint8) {
	for f := range 256 {
		f := uint32(f)
		h := b2i(f&fE != 0) + b2i(f&fW != 0)
		v := b2i(f&fN != 0) + b2i(f&fS != 0)
		d := b2i(f&fNE != 0) + b2i(f&fSE != 0) + b2i(f&fSW != 0) + b2i(f&fNW != 0)
		t[orientLL][f] = zcLLHL(h, v, d)
		t[orientHL][f] = zcLLHL(v, h, d)
		var n uint8
		switch hv := h + v; {
		case d >= 3:
			n = 8
		case d == 2 && hv >= 1:
			n = 7
		case d == 2:
			n = 6
		case d == 1 && hv >= 2:
			n = 5
		case d == 1 && hv == 1:
			n = 4
		case d == 1:
			n = 3
		case hv >= 2:
			n = 2
		case hv == 1:
			n = 1
		}
		t[orientHH][f] = n
	}
	return t
}

func zcLLHL(h, v, d int) uint8 {
	switch {
	case h == 2:
		return 8
	case h == 1 && v >= 1:
		return 7
	case h == 1 && d >= 1:
		return 6
	case h == 1:
		return 5
	case v == 2:
		return 4
	case v == 1:
		return 3
	case d >= 2:
		return 2
	case d == 1:
		return 1
	}
	return 0
}

// buildSC indexes the sign context by flag bits 4-11 (Table D.3); bit 7 of
// an entry is the bit to XOR with the decoded sign.
func buildSC() (t [256]uint8) {
	contrib := func(sig, neg bool) int {
		switch {
		case !sig:
			return 0
		case neg:
			return -1
		}
		return 1
	}
	for i := range 256 {
		f := uint32(i) << 4
		h := contrib(f&fE != 0, f&fNegE != 0) + contrib(f&fW != 0, f&fNegW != 0)
		v := contrib(f&fN != 0, f&fNegN != 0) + contrib(f&fS != 0, f&fNegS != 0)
		h = max(-1, min(1, h))
		v = max(-1, min(1, v))
		var cx, x int
		switch {
		case h == 1:
			cx = 12 + v
		case h == 0 && v == 0:
			cx = 9
		case h == 0:
			cx, x = 10, b2i(v < 0)
		default:
			cx, x = 12-v, 1
		}
		t[i] = uint8(cx | x<<7)
	}
	return t
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// segment is a codeword segment of a code-block: a number of coding passes
// and the length of their data.
type segment struct {
	passes int
	length int
}

// t1Decoder decodes code-blocks; its buffers are reused between blocks.
type t1Decoder struct {
	w, h   int
	stride int      // w+2: flags have a border of one sample
	flags  []uint32 // (w+2)×(h+2)
	data   []uint32 // w×h magnitudes at twice their scale
	mq     mqDecoder
	raw    rawDecoder
}

// decode runs the coding passes of a code-block. p0 is the bit-plane of the
// first (cleanup) pass; data holds the concatenated segments.
func (t *t1Decoder) decode(w, h int, orient int, style byte, p0 int, segs []segment, data []byte) {
	t.w, t.h, t.stride = w, h, w+2
	nf := (w + 2) * (h + 2)
	if cap(t.flags) < nf {
		t.flags = make([]uint32, nf)
	} else {
		t.flags = t.flags[:nf]
		clear(t.flags)
	}
	if cap(t.data) < w*h {
		t.data = make([]uint32, w*h)
	} else {
		t.data = t.data[:w*h]
		clear(t.data)
	}
	t.mq.resetContexts()
	zc := &zcLUT[orient]
	causal := style&cbCausal != 0
	plane, kind, pass := p0, 2, 0 // kind: 0 significance, 1 refinement, 2 cleanup
	off := 0
	for _, seg := range segs {
		end := off + seg.length
		if end > len(data) || end < off {
			end = len(data)
		}
		b := data[off:end]
		off = end
		raw := style&cbBypass != 0 && pass >= 10 && kind != 2
		if raw {
			t.raw.init(b)
		} else {
			t.mq.init(b)
		}
		for range seg.passes {
			if plane < 0 {
				return
			}
			if (style&cbBypass != 0 && pass >= 10 && kind != 2) != raw {
				return // passes do not match the segment kind
			}
			switch {
			case kind == 0 && raw:
				t.sigPassRaw(plane, causal)
			case kind == 0:
				t.sigPass(plane, zc, causal)
			case kind == 1 && raw:
				t.refPassRaw(plane)
			case kind == 1:
				t.refPass(plane, causal)
			default:
				t.cleanupPass(plane, zc, causal)
				if style&cbSegSym != 0 {
					for range 4 {
						t.mq.decode(ctxUni)
					}
				}
			}
			if style&cbReset != 0 && !raw {
				t.mq.resetContexts()
			}
			pass++
			if kind++; kind == 3 {
				kind = 0
				plane--
			}
		}
	}
}

// setSig makes the sample at flag index i significant with the given sign
// and records it in the flags of its neighbours.
func (t *t1Decoder) setSig(i int, neg uint32) {
	f := t.flags
	s := t.stride
	if neg != 0 {
		f[i] |= fSig | fNeg
		f[i-s] |= fS | fNegS
		f[i+s] |= fN | fNegN
		f[i-1] |= fE | fNegE
		f[i+1] |= fW | fNegW
	} else {
		f[i] |= fSig
		f[i-s] |= fS
		f[i+s] |= fN
		f[i-1] |= fE
		f[i+1] |= fW
	}
	f[i-s-1] |= fSE
	f[i-s+1] |= fSW
	f[i+s-1] |= fNE
	f[i+s+1] |= fNW
}

// decodeSign decodes the sign of a sample whose (possibly causal-masked)
// flags are f.
func (t *t1Decoder) decodeSign(f uint32) uint32 {
	sc := scLUT[f>>4&0xff]
	return t.mq.decode(int(sc&0x1f)) ^ uint32(sc>>7)
}

func (t *t1Decoder) sigPass(p int, zc *[256]uint8, causal bool) {
	w, h, s := t.w, t.h, t.stride
	flags, data := t.flags, t.data
	v := uint32(3) << p
	for y0 := 0; y0 < h; y0 += 4 {
		ye := min(y0+4, h)
		for x := range w {
			fi := (y0+1)*s + x + 1
			di := y0*w + x
			for y := y0; y < ye; y++ {
				f := flags[fi]
				if f&(fSig|fVisit) == 0 && f&fNbr != 0 {
					if causal && y == y0+3 {
						f &= fCausal
					}
					if f&fNbr != 0 {
						if t.mq.decode(int(zc[f&fNbr])) != 0 {
							neg := t.decodeSign(f)
							data[di] = v
							t.setSig(fi, neg)
						}
						flags[fi] |= fVisit
					}
				}
				fi += s
				di += w
			}
		}
	}
}

func (t *t1Decoder) sigPassRaw(p int, causal bool) {
	w, h, s := t.w, t.h, t.stride
	flags, data := t.flags, t.data
	v := uint32(3) << p
	for y0 := 0; y0 < h; y0 += 4 {
		ye := min(y0+4, h)
		for x := range w {
			fi := (y0+1)*s + x + 1
			di := y0*w + x
			for y := y0; y < ye; y++ {
				f := flags[fi]
				if f&(fSig|fVisit) == 0 && f&fNbr != 0 {
					if causal && y == y0+3 {
						f &= fCausal
					}
					if f&fNbr != 0 {
						if t.raw.bit() != 0 {
							neg := t.raw.bit()
							data[di] = v
							t.setSig(fi, neg)
						}
						flags[fi] |= fVisit
					}
				}
				fi += s
				di += w
			}
		}
	}
}

func (t *t1Decoder) refPass(p int, causal bool) {
	w, h, s := t.w, t.h, t.stride
	flags, data := t.flags, t.data
	half := uint32(1) << p
	for y0 := 0; y0 < h; y0 += 4 {
		ye := min(y0+4, h)
		for x := range w {
			fi := (y0+1)*s + x + 1
			di := y0*w + x
			for y := y0; y < ye; y++ {
				f := flags[fi]
				if f&(fSig|fVisit) == fSig {
					cx := 16
					if f&fRef == 0 {
						if causal && y == y0+3 {
							f &= fCausal
						}
						cx = 14
						if f&fNbr != 0 {
							cx = 15
						}
					}
					if t.mq.decode(cx) != 0 {
						data[di] += half
					} else {
						data[di] -= half
					}
					flags[fi] |= fRef
				}
				fi += s
				di += w
			}
		}
	}
}

func (t *t1Decoder) refPassRaw(p int) {
	w, h, s := t.w, t.h, t.stride
	flags, data := t.flags, t.data
	half := uint32(1) << p
	for y0 := 0; y0 < h; y0 += 4 {
		ye := min(y0+4, h)
		for x := range w {
			fi := (y0+1)*s + x + 1
			di := y0*w + x
			for y := y0; y < ye; y++ {
				if flags[fi]&(fSig|fVisit) == fSig {
					if t.raw.bit() != 0 {
						data[di] += half
					} else {
						data[di] -= half
					}
					flags[fi] |= fRef
				}
				fi += s
				di += w
			}
		}
	}
}

func (t *t1Decoder) cleanupPass(p int, zc *[256]uint8, causal bool) {
	w, h, s := t.w, t.h, t.stride
	flags, data := t.flags, t.data
	v := uint32(3) << p
	for y0 := 0; y0 < h; y0 += 4 {
		ye := min(y0+4, h)
		for x := range w {
			fi := (y0+1)*s + x + 1
			di := y0*w + x
			y := y0
			if ye-y0 == 4 {
				f3 := flags[fi+3*s]
				if causal {
					f3 &= fCausal
				}
				if (flags[fi]|flags[fi+s]|flags[fi+2*s]|f3)&(fSig|fVisit|fNbr) == 0 {
					// Run-length mode.
					if t.mq.decode(ctxRL) == 0 {
						continue
					}
					r := int(t.mq.decode(ctxUni)) << 1
					r |= int(t.mq.decode(ctxUni))
					y += r
					fi += r * s
					di += r * w
					f := flags[fi]
					if causal && y == y0+3 {
						f &= fCausal
					}
					neg := t.decodeSign(f)
					data[di] = v
					t.setSig(fi, neg)
					y++
					fi += s
					di += w
				}
			}
			for ; y < ye; y++ {
				f := flags[fi]
				if f&(fSig|fVisit) != 0 {
					flags[fi] = f &^ fVisit
				} else {
					if causal && y == y0+3 {
						f &= fCausal
					}
					if t.mq.decode(int(zc[f&fNbr])) != 0 {
						neg := t.decodeSign(f)
						data[di] = v
						t.setSig(fi, neg)
					}
				}
				fi += s
				di += w
			}
		}
	}
}

// The code tables below are from golang.org/x/image/ccitt (gen.go):
//
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file of golang.org/x/image.

package tiff

import (
	"errors"
	"math/bits"
	"strconv"
	"strings"
)

// CCITT fax decoding (ITU-T T.4 and T.6), for TIFF compressions 2
// (modified Huffman, rows byte-aligned), 3 (T.4: one-dimensional or, with
// T4Options bit 0, two-dimensional coding; EOL codes, with fill bits
// before them under T4Options bit 2) and 4 (T.6).
//
// golang.org/x/image/ccitt reads neither two-dimensional T.4 data nor fill
// bits, both of which fax software writes by default. This decoder takes
// both, finds EOL codes wherever they are, and after a damaged row of T.4
// data carries on at the next EOL, as fax receivers do.

// ccittMode selects the coding.
type ccittMode int

const (
	ccittRLE ccittMode = iota // modified Huffman, rows start on byte boundaries, no EOL
	ccittT4                   // T.4, rows start with EOL codes (sometimes left out)
	ccittT6                   // T.6, two-dimensional rows without EOL codes
)

// Run-length codes: "run:code" pairs, the terminating codes (runs of 0 to
// 63) followed by the make-up codes (multiples of 64).
const (
	whiteCodeList = "0:00110101 1:000111 2:0111 3:1000 4:1011 5:1100 6:1110 7:1111 " +
		"8:10011 9:10100 10:00111 11:01000 12:001000 13:000011 14:110100 15:110101 " +
		"16:101010 17:101011 18:0100111 19:0001100 20:0001000 21:0010111 22:0000011 23:0000100 " +
		"24:0101000 25:0101011 26:0010011 27:0100100 28:0011000 29:00000010 30:00000011 31:00011010 " +
		"32:00011011 33:00010010 34:00010011 35:00010100 36:00010101 37:00010110 38:00010111 39:00101000 " +
		"40:00101001 41:00101010 42:00101011 43:00101100 44:00101101 45:00000100 46:00000101 47:00001010 " +
		"48:00001011 49:01010010 50:01010011 51:01010100 52:01010101 53:00100100 54:00100101 55:01011000 " +
		"56:01011001 57:01011010 58:01011011 59:01001010 60:01001011 61:00110010 62:00110011 63:00110100 " +
		"64:11011 128:10010 192:010111 256:0110111 320:00110110 384:00110111 448:01100100 512:01100101 " +
		"576:01101000 640:01100111 704:011001100 768:011001101 832:011010010 896:011010011 960:011010100 1024:011010101 " +
		"1088:011010110 1152:011010111 1216:011011000 1280:011011001 1344:011011010 1408:011011011 1472:010011000 1536:010011001 " +
		"1600:010011010 1664:011000 1728:010011011 1792:00000001000 1856:00000001100 1920:00000001101 1984:000000010010 2048:000000010011 " +
		"2112:000000010100 2176:000000010101 2240:000000010110 2304:000000010111 2368:000000011100 2432:000000011101 2496:000000011110 2560:000000011111 "
	blackCodeList = "0:0000110111 1:010 2:11 3:10 4:011 5:0011 6:0010 7:00011 " +
		"8:000101 9:000100 10:0000100 11:0000101 12:0000111 13:00000100 14:00000111 15:000011000 " +
		"16:0000010111 17:0000011000 18:0000001000 19:00001100111 20:00001101000 21:00001101100 22:00000110111 23:00000101000 " +
		"24:00000010111 25:00000011000 26:000011001010 27:000011001011 28:000011001100 29:000011001101 30:000001101000 31:000001101001 " +
		"32:000001101010 33:000001101011 34:000011010010 35:000011010011 36:000011010100 37:000011010101 38:000011010110 39:000011010111 " +
		"40:000001101100 41:000001101101 42:000011011010 43:000011011011 44:000001010100 45:000001010101 46:000001010110 47:000001010111 " +
		"48:000001100100 49:000001100101 50:000001010010 51:000001010011 52:000000100100 53:000000110111 54:000000111000 55:000000100111 " +
		"56:000000101000 57:000001011000 58:000001011001 59:000000101011 60:000000101100 61:000001011010 62:000001100110 63:000001100111 " +
		"64:0000001111 128:000011001000 192:000011001001 256:000001011011 320:000000110011 384:000000110100 448:000000110101 512:0000001101100 " +
		"576:0000001101101 640:0000001001010 704:0000001001011 768:0000001001100 832:0000001001101 896:0000001110010 960:0000001110011 1024:0000001110100 " +
		"1088:0000001110101 1152:0000001110110 1216:0000001110111 1280:0000001010010 1344:0000001010011 1408:0000001010100 1472:0000001010101 1536:0000001011010 " +
		"1600:0000001011011 1664:0000001100100 1728:0000001100101 1792:00000001000 1856:00000001100 1920:00000001101 1984:000000010010 2048:000000010011 " +
		"2112:000000010100 2176:000000010101 2240:000000010110 2304:000000010111 2368:000000011100 2432:000000011101 2496:000000011110 2560:000000011111 "
)

// codeEntry is a decoded code: its value and its length in bits (0 for no
// code).
type codeEntry struct {
	val int16
	n   uint8
}

const lutBits = 13 // the longest run-length code

var whiteLUT, blackLUT [1 << lutBits]codeEntry

// Two-dimensional mode codes.
const (
	modePass = iota
	modeHorizontal
	modeV0
	modeVR1
	modeVR2
	modeVR3
	modeVL1
	modeVL2
	modeVL3
	modeExt
)

var modeLUT [1 << 7]codeEntry

func init() {
	fill := func(lut []codeEntry, bitsLen int, code string, val int) {
		c, _ := strconv.ParseUint(code, 2, 32)
		shift := bitsLen - len(code)
		for i := range 1 << shift {
			lut[int(c)<<shift|i] = codeEntry{int16(val), uint8(len(code))}
		}
	}
	for _, t := range []struct {
		list string
		lut  []codeEntry
	}{{whiteCodeList, whiteLUT[:]}, {blackCodeList, blackLUT[:]}} {
		for _, kv := range strings.Fields(t.list) {
			k, code, _ := strings.Cut(kv, ":")
			v, _ := strconv.Atoi(k)
			fill(t.lut, lutBits, code, v)
		}
	}
	for code, m := range map[string]int{"0001": modePass, "001": modeHorizontal, "1": modeV0,
		"011": modeVR1, "000011": modeVR2, "0000011": modeVR3, "010": modeVL1, "000010": modeVL2, "0000010": modeVL3, "0000001": modeExt} {
		fill(modeLUT[:], 7, code, m)
	}
}

// bitReader reads bits most significant first.
type bitReader struct {
	b   []byte
	pos int // in bits
}

// peek returns the next n (≤ 24) bits, with zeros past the end.
func (r *bitReader) peek(n int) uint32 {
	var v uint32
	i, o := r.pos>>3, r.pos&7
	for k := range 4 {
		v <<= 8
		if i+k < len(r.b) {
			v |= uint32(r.b[i+k])
		}
	}
	return v << o >> (32 - n)
}

func (r *bitReader) skip(n int)  { r.pos += n }
func (r *bitReader) done() bool  { return r.pos >= 8*len(r.b) }
func (r *bitReader) alignByte()  { r.pos = (r.pos + 7) &^ 7 }
func (r *bitReader) bit() uint32 { v := r.peek(1); r.pos++; return v }
func (r *bitReader) zeros() int {
	n := 0
	for !r.done() && r.peek(1) == 0 {
		r.pos++
		n++
	}
	return n
}

var errCCITT = errors.New("tiff: bad CCITT data")

// eol consumes an EOL code (at least 11 zero bits, which fill bits
// lengthen, then a 1) if one is next, and reports whether it was.
func (r *bitReader) eol() bool {
	start := r.pos
	if r.zeros() >= 11 && !r.done() {
		r.pos++
		return true
	}
	r.pos = start
	return false
}

// nextEOL moves past the next EOL code, for resynchronizing after a
// damaged row; it reports whether there was one.
func (r *bitReader) nextEOL() bool {
	for !r.done() {
		if r.eol() {
			return true
		}
		r.pos++
	}
	return false
}

// run decodes a run length of the given colour: make-up codes, then a
// terminating code.
func (r *bitReader) run(white bool) (int, error) {
	lut := &blackLUT
	if white {
		lut = &whiteLUT
	}
	total := 0
	for {
		if r.done() {
			return 0, errCCITT
		}
		e := lut[r.peek(lutBits)]
		if e.n == 0 {
			return 0, errCCITT
		}
		r.skip(int(e.n))
		total += int(e.val)
		if e.val < 64 {
			return total, nil
		}
	}
}

// decodeCCITT decodes rows of width pixels into packed rows (1 bits for
// black, rows padded to bytes). lsb reverses the bits of each byte first
// (FillOrder 2). damaged reports rows that did not decode, which are left
// white.
func decodeCCITT(src []byte, mode ccittMode, twoD, lsb bool, width, rows int) (out []byte, damaged bool) {
	if lsb {
		b := make([]byte, len(src))
		for i, v := range src {
			b[i] = bits.Reverse8(v)
		}
		src = b
	}
	rb := (width + 7) / 8
	out = make([]byte, rb*rows)
	r := &bitReader{b: src}
	ref := []int{width, width} // changing elements of the reference row
	cur := make([]int, 0, 64)
	for y := range rows {
		if r.done() {
			return out, true
		}
		row2D := mode == ccittT6
		if mode == ccittT4 {
			r.eol()
			if twoD {
				row2D = r.bit() == 0
			}
		}
		var err error
		cur = cur[:0]
		if row2D {
			cur, err = decode2D(r, ref, cur, width)
		} else {
			cur, err = decode1D(r, cur, width)
		}
		if err != nil {
			damaged = true
			if mode != ccittT4 || !r.nextEOL() {
				return out, true
			}
			r.pos -= 12 // let the next row read the EOL (and tag bit) again
			if twoD {
				// The row after a damaged one must not refer to it.
				ref = append(ref[:0], width, width)
			}
			continue
		}
		setRuns(out[y*rb:(y+1)*rb], cur, width)
		ref = append(append(ref[:0], cur...), width, width)
		if mode == ccittRLE {
			r.alignByte()
		}
	}
	return out, damaged
}

// decode1D decodes a row of alternating white and black runs, appending
// the changing elements (the positions where the colour changes) to cur.
func decode1D(r *bitReader, cur []int, width int) ([]int, error) {
	a0, white := 0, true
	for a0 < width {
		n, err := r.run(white)
		if err != nil {
			return cur, err
		}
		a0 = min(a0+n, width)
		cur = append(cur, a0)
		white = !white
	}
	return cur, nil
}

// decode2D decodes a row coded against the reference row ref (its
// changing elements, ending with width twice).
func decode2D(r *bitReader, ref, cur []int, width int) ([]int, error) {
	a0, white := -1, true
	k := 0 // the first changing element of ref right of a0
	for a0 < width {
		for k < len(ref) && ref[k] <= a0 {
			k++
		}
		// b1 is the first changing element right of a0 of the opposite
		// colour to a0's: elements at even indexes turn white to black.
		i := k
		if i%2 == 0 != white {
			i++
		}
		b1, b2 := width, width
		if i < len(ref) {
			b1 = ref[i]
		}
		if i+1 < len(ref) {
			b2 = ref[i+1]
		}
		if r.done() {
			return cur, errCCITT
		}
		e := modeLUT[r.peek(7)]
		if e.n == 0 {
			return cur, errCCITT
		}
		r.skip(int(e.n))
		switch e.val {
		case modePass:
			a0 = b2
		case modeHorizontal:
			start := max(a0, 0)
			n1, err := r.run(white)
			if err != nil {
				return cur, err
			}
			n2, err := r.run(!white)
			if err != nil {
				return cur, err
			}
			a1 := min(start+n1, width)
			a0 = min(a1+n2, width)
			cur = append(cur, a1, a0)
		case modeExt:
			return cur, errCCITT // uncompressed mode
		default:
			a1 := b1 + []int{0, 1, 2, 3, -1, -2, -3}[e.val-modeV0]
			if a1 < max(a0, 0) || a1 > width {
				return cur, errCCITT
			}
			cur = append(cur, a1)
			a0 = a1
			white = !white
		}
	}
	return cur, nil
}

// setRuns sets the bits of the black runs of a row given by its changing
// elements.
func setRuns(row []byte, changes []int, width int) {
	for i := 0; i < len(changes); i += 2 {
		from := changes[i]
		to := width
		if i+1 < len(changes) {
			to = changes[i+1]
		}
		for x := from; x < min(to, width); x++ {
			row[x>>3] |= 0x80 >> (x & 7)
		}
	}
}

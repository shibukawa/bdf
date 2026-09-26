package jbig2

// MMR is the two-dimensional coding of ITU-T T.6 (Group 4 facsimile), with
// 1 = black. The decoder reports how many bytes it used, since symbol
// dictionaries and halftone regions place more data after an MMR bitmap.

// Coding modes (T.4 Table 4).
const (
	modeInvalid = iota
	modePass
	modeHorizontal
	modeV0
	modeVR1
	modeVR2
	modeVR3
	modeVL1
	modeVL2
	modeVL3
	modeEOL // the prefix of an end-of-line code
)

// modeCodes lists the mode codes as {mode, code, length}.
var modeCodes = [][3]int{
	{modeV0, 0b1, 1}, {modeVR1, 0b011, 3}, {modeVL1, 0b010, 3}, {modeHorizontal, 0b001, 3},
	{modePass, 0b0001, 4}, {modeVR2, 0b000011, 6}, {modeVL2, 0b000010, 6},
	{modeVR3, 0b0000011, 7}, {modeVL3, 0b0000010, 7}, {modeEOL, 0b0000000, 7},
}

// Run length codes (T.4 Tables 2 and 3) as {code, length}: terminating codes
// for 0 to 63, then make-up codes for 64 to 1728 in steps of 64.
var (
	whiteRunCodes = [][2]int{
		{0x35, 8}, {0x07, 6}, {0x07, 4}, {0x08, 4}, {0x0b, 4}, {0x0c, 4}, {0x0e, 4}, {0x0f, 4},
		{0x13, 5}, {0x14, 5}, {0x07, 5}, {0x08, 5}, {0x08, 6}, {0x03, 6}, {0x34, 6}, {0x35, 6},
		{0x2a, 6}, {0x2b, 6}, {0x27, 7}, {0x0c, 7}, {0x08, 7}, {0x17, 7}, {0x03, 7}, {0x04, 7},
		{0x28, 7}, {0x2b, 7}, {0x13, 7}, {0x24, 7}, {0x18, 7}, {0x02, 8}, {0x03, 8}, {0x1a, 8},
		{0x1b, 8}, {0x12, 8}, {0x13, 8}, {0x14, 8}, {0x15, 8}, {0x16, 8}, {0x17, 8}, {0x28, 8},
		{0x29, 8}, {0x2a, 8}, {0x2b, 8}, {0x2c, 8}, {0x2d, 8}, {0x04, 8}, {0x05, 8}, {0x0a, 8},
		{0x0b, 8}, {0x52, 8}, {0x53, 8}, {0x54, 8}, {0x55, 8}, {0x24, 8}, {0x25, 8}, {0x58, 8},
		{0x59, 8}, {0x5a, 8}, {0x5b, 8}, {0x4a, 8}, {0x4b, 8}, {0x32, 8}, {0x33, 8}, {0x34, 8},
		{0x1b, 5}, {0x12, 5}, {0x17, 6}, {0x37, 7}, {0x36, 8}, {0x37, 8}, {0x64, 8}, {0x65, 8},
		{0x68, 8}, {0x67, 8}, {0xcc, 9}, {0xcd, 9}, {0xd2, 9}, {0xd3, 9}, {0xd4, 9}, {0xd5, 9},
		{0xd6, 9}, {0xd7, 9}, {0xd8, 9}, {0xd9, 9}, {0xda, 9}, {0xdb, 9}, {0x98, 9}, {0x99, 9},
		{0x9a, 9}, {0x18, 6}, {0x9b, 9},
	}
	blackRunCodes = [][2]int{
		{0x37, 10}, {0x02, 3}, {0x03, 2}, {0x02, 2}, {0x03, 3}, {0x03, 4}, {0x02, 4}, {0x03, 5},
		{0x05, 6}, {0x04, 6}, {0x04, 7}, {0x05, 7}, {0x07, 7}, {0x04, 8}, {0x07, 8}, {0x18, 9},
		{0x17, 10}, {0x18, 10}, {0x08, 10}, {0x67, 11}, {0x68, 11}, {0x6c, 11}, {0x37, 11}, {0x28, 11},
		{0x17, 11}, {0x18, 11}, {0xca, 12}, {0xcb, 12}, {0xcc, 12}, {0xcd, 12}, {0x68, 12}, {0x69, 12},
		{0x6a, 12}, {0x6b, 12}, {0xd2, 12}, {0xd3, 12}, {0xd4, 12}, {0xd5, 12}, {0xd6, 12}, {0xd7, 12},
		{0x6c, 12}, {0x6d, 12}, {0xda, 12}, {0xdb, 12}, {0x54, 12}, {0x55, 12}, {0x56, 12}, {0x57, 12},
		{0x64, 12}, {0x65, 12}, {0x52, 12}, {0x53, 12}, {0x24, 12}, {0x37, 12}, {0x38, 12}, {0x27, 12},
		{0x28, 12}, {0x58, 12}, {0x59, 12}, {0x2b, 12}, {0x2c, 12}, {0x5a, 12}, {0x66, 12}, {0x67, 12},
		{0x0f, 10}, {0xc8, 12}, {0xc9, 12}, {0x5b, 12}, {0x33, 12}, {0x34, 12}, {0x35, 12}, {0x6c, 13},
		{0x6d, 13}, {0x4a, 13}, {0x4b, 13}, {0x4c, 13}, {0x4d, 13}, {0x72, 13}, {0x73, 13}, {0x74, 13},
		{0x75, 13}, {0x76, 13}, {0x77, 13}, {0x52, 13}, {0x53, 13}, {0x54, 13}, {0x55, 13}, {0x5a, 13},
		{0x5b, 13}, {0x64, 13}, {0x65, 13},
	}
	// extendedRunCodes are the make-up codes for 1792 to 2560, common to both colours.
	extendedRunCodes = [][2]int{
		{0x08, 11}, {0x0c, 11}, {0x0d, 11}, {0x12, 12}, {0x13, 12}, {0x14, 12}, {0x15, 12},
		{0x16, 12}, {0x17, 12}, {0x1c, 12}, {0x1d, 12}, {0x1e, 12}, {0x1f, 12},
	}
)

// runLength returns the run length a code in the run code lists stands for.
func runLength(i int) int {
	if i < 64 {
		return i
	}
	return (i - 63) * 64
}

// codeEntry is a lookup table entry: a value and its code length (0 = invalid).
type codeEntry struct {
	v uint16
	n uint8
}

const runLookupBits = 13

var (
	modeLookup [1 << 7]codeEntry
	runLookup  [2][1 << runLookupBits]codeEntry // white, black
)

func init() {
	fill := func(t []codeEntry, bits int, code, n, v int) {
		shift := bits - n
		for i := 0; i < 1<<shift; i++ {
			t[code<<shift|i] = codeEntry{uint16(v), uint8(n)}
		}
	}
	for _, c := range modeCodes {
		fill(modeLookup[:], 7, c[1], c[2], c[0])
	}
	for color, codes := range [2][][2]int{whiteRunCodes, blackRunCodes} {
		for i, c := range codes {
			fill(runLookup[color][:], runLookupBits, c[0], c[1], runLength(i))
		}
		for i, c := range extendedRunCodes {
			fill(runLookup[color][:], runLookupBits, c[0], c[1], 1792+64*i)
		}
	}
}

// eofb is the end-of-facsimile-block code: two end-of-line codes.
const eofb = 0x001001

// mmrReader reads MMR codes; past the end of the data it reads 0 bits.
type mmrReader struct {
	data []byte
	pos  int // in bits
}

// peek returns the next n bits, n ≤ 24.
func (m *mmrReader) peek(n int) uint32 {
	i := m.pos >> 3
	var v uint32
	for k := 0; k < 4; k++ {
		v <<= 8
		if i+k < len(m.data) {
			v |= uint32(m.data[i+k])
		}
	}
	return v << uint(m.pos&7) >> uint(32-n)
}

// run decodes a run length of the given colour (0 = white).
func (m *mmrReader) run(color int) (int, error) {
	total := 0
	for {
		e := runLookup[color][m.peek(runLookupBits)]
		if e.n == 0 {
			return 0, errorf("invalid MMR run length code")
		}
		m.pos += int(e.n)
		total += int(e.v)
		if e.v < 64 {
			return total, nil
		}
		if total > 1<<26 {
			return 0, errorf("invalid MMR run length")
		}
	}
}

// line decodes one coding line against the reference line ref, a list of
// changing elements followed by at least three entries equal to the width w.
// It appends the changing elements of the new line to cur.
func (m *mmrReader) line(ref, cur []int, w int) ([]int, error) {
	a0, color, i := -1, 0, 0
	for steps := 0; a0 < w; steps++ {
		if steps > 2*w+16 || m.pos>>3 > len(m.data)+8 {
			return nil, errorf("invalid MMR data")
		}
		// b1 is the first changing element of the reference line after a0 of
		// the colour opposite to a0's; b2 the one after it.
		for i > 0 && ref[i-1] > a0 {
			i--
		}
		for ref[i] <= a0 {
			i++
		}
		if i&1 != color {
			i++
		}
		b1, b2 := ref[i], ref[i+1]
		e := modeLookup[m.peek(7)]
		m.pos += int(e.n)
		start := max(a0, 0)
		switch e.v {
		case modePass:
			a0 = b2
		case modeHorizontal:
			r1, err := m.run(color)
			if err != nil {
				return nil, err
			}
			r2, err := m.run(1 - color)
			if err != nil {
				return nil, err
			}
			a1 := min(start+r1, w)
			a0 = min(a1+r2, w)
			cur = append(cur, a1, a0)
		case modeV0, modeVR1, modeVR2, modeVR3, modeVL1, modeVL2, modeVL3:
			a1 := b1 + vertical[e.v]
			if a1 < start {
				return nil, errorf("invalid MMR data")
			}
			a1 = min(a1, w)
			cur = append(cur, a1)
			a0 = a1
			color = 1 - color
		default:
			return nil, errorf("invalid or unsupported MMR code")
		}
	}
	return cur, nil
}

// vertical holds the offset of a1 from b1 in each vertical mode.
var vertical = [...]int{modeV0: 0, modeVR1: 1, modeVR2: 2, modeVR3: 3, modeVL1: -1, modeVL2: -2, modeVL3: -3}

// decodeMMR decodes a w×h MMR coded bitmap. An end-of-facsimile-block code
// before the last row leaves the remaining rows white; one right after the last
// row is consumed. It returns the bitmap and the number of bytes used.
func (d *decoder) decodeMMR(data []byte, w, h int) (*Bitmap, int, error) {
	bm, err := d.newBitmap(w, h)
	if err != nil {
		return nil, 0, err
	}
	m := mmrReader{data: data}
	ref := []int{w, w, w}
	var cur []int
	for y := 0; y < h; y++ {
		if m.peek(24) == eofb {
			m.pos += 24
			return bm, (m.pos + 7) >> 3, nil
		}
		if cur, err = m.line(ref, cur[:0], w); err != nil {
			return nil, 0, err
		}
		row := bm.row(y)
		for k := 0; k < len(cur); k += 2 {
			end := w
			if k+1 < len(cur) {
				end = cur[k+1]
			}
			fillSpan(row, cur[k], end)
		}
		// The new line becomes the reference line, without empty runs.
		ref = ref[:0]
		for _, c := range cur {
			if n := len(ref); n > 0 && ref[n-1] == c {
				ref = ref[:n-1]
			} else {
				ref = append(ref, c)
			}
		}
		ref = append(ref, w, w, w)
	}
	if m.peek(24) == eofb {
		m.pos += 24
	}
	return bm, min((m.pos+7)>>3, len(data)), nil
}

// fillSpan sets pixels x0 to x1−1 of a row.
func fillSpan(row []byte, x0, x1 int) {
	for x0 < x1 && x0&7 != 0 {
		row[x0>>3] |= 0x80 >> uint(x0&7)
		x0++
	}
	for ; x0+8 <= x1; x0 += 8 {
		row[x0>>3] = 0xff
	}
	for ; x0 < x1; x0++ {
		row[x0>>3] |= 0x80 >> uint(x0&7)
	}
}

package barcode

import (
	"fmt"
	"strings"
)

// QR codes (ISO/IEC 18004, model 2). A text is one segment in the smallest
// of the numeric, alphanumeric and byte modes that holds it (bytes are the
// UTF-8 of the text), in the smallest version that has room for it at the
// level of error correction, with the mask whose penalty is the lowest.

// The error correction codewords of a block and the number of blocks, by
// level (L, M, Q, H) and version (1 to 40).
var (
	qrEccLen = [4][41]int8{
		{-1, 7, 10, 15, 20, 26, 18, 20, 24, 30, 18, 20, 24, 26, 30, 22, 24, 28, 30, 28, 28, 28, 28, 30, 30, 26, 28, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
		{-1, 10, 16, 26, 18, 24, 16, 18, 22, 22, 26, 30, 22, 22, 24, 24, 28, 28, 26, 26, 26, 26, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28},
		{-1, 13, 22, 18, 26, 18, 24, 18, 22, 20, 24, 28, 26, 24, 20, 30, 24, 28, 28, 26, 30, 28, 30, 30, 30, 30, 28, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
		{-1, 17, 28, 22, 16, 22, 28, 26, 26, 24, 28, 24, 28, 22, 24, 24, 30, 28, 28, 26, 28, 30, 24, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30},
	}
	qrBlocks = [4][41]int8{
		{-1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 4, 4, 4, 4, 4, 6, 6, 6, 6, 7, 8, 8, 9, 9, 10, 12, 12, 12, 13, 14, 15, 16, 17, 18, 19, 19, 20, 21, 22, 24, 25},
		{-1, 1, 1, 1, 2, 2, 4, 4, 4, 5, 5, 5, 8, 9, 9, 10, 10, 11, 13, 14, 16, 17, 17, 18, 20, 21, 23, 25, 26, 28, 29, 31, 33, 35, 37, 38, 40, 43, 45, 47, 49},
		{-1, 1, 1, 2, 2, 4, 4, 6, 6, 8, 8, 8, 10, 12, 16, 12, 17, 16, 18, 21, 20, 23, 23, 25, 27, 29, 34, 34, 35, 38, 40, 43, 45, 48, 51, 53, 56, 59, 62, 65, 68},
		{-1, 1, 1, 2, 4, 4, 4, 5, 6, 8, 8, 11, 11, 16, 16, 18, 16, 19, 21, 25, 25, 25, 34, 30, 32, 35, 37, 40, 42, 45, 48, 51, 54, 57, 60, 63, 66, 70, 74, 77, 81},
	}
	// the bits of the levels in the format information
	qrLevelBits = [4]int{1, 0, 3, 2}
)

// Modes of a segment.
const (
	qrNumeric = iota
	qrAlnum
	qrByte
)

const qrAlnumChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

// qrRawModules is the number of modules of a version that hold data and
// error correction, the function patterns left out.
func qrRawModules(ver int) int {
	n := (16*ver+128)*ver + 64
	if ver >= 2 {
		align := ver/7 + 2
		n -= (25*align-10)*align - 55
		if ver >= 7 {
			n -= 36
		}
	}
	return n
}

// qrDataCodewords is the number of codewords of data of a version at a
// level.
func qrDataCodewords(ver, level int) int {
	return qrRawModules(ver)/8 - int(qrEccLen[level][ver])*int(qrBlocks[level][ver])
}

// qrCountBits is the width of the character count of a mode in a version.
func qrCountBits(mode, ver int) int {
	i := 0
	if ver >= 27 {
		i = 2
	} else if ver >= 10 {
		i = 1
	}
	return [3][3]int{{10, 12, 14}, {9, 11, 13}, {8, 16, 16}}[mode][i]
}

// bits is a string of bits.
type bits []bool

func (b *bits) put(v, n int) {
	for i := n - 1; i >= 0; i-- {
		*b = append(*b, v>>i&1 != 0)
	}
}

// encodeQR makes the QR code of a text. mask is the mask to use (-1: the
// one with the lowest penalty).
func encodeQR(text, level string, quiet, mask int) (*Symbol, error) {
	lv := strings.Index("LMQH", strings.ToUpper(level))
	switch {
	case level == "":
		lv = 1
	case lv < 0 || len(level) != 1:
		return nil, fmt.Errorf("barcode: %q is not a level of error correction of a QR code (L, M, Q or H)", level)
	}
	// the mode: the smallest that holds every character
	mode := qrNumeric
	for i := 0; i < len(text) && mode != qrByte; i++ {
		switch c := text[i]; {
		case c >= '0' && c <= '9':
		case strings.IndexByte(qrAlnumChars, c) >= 0:
			mode = qrAlnum
		default:
			mode = qrByte
		}
	}
	var data bits
	count := len(text)
	switch mode {
	case qrNumeric:
		for i := 0; i < len(text); i += 3 {
			n, v := min(3, len(text)-i), 0
			for _, c := range []byte(text[i : i+n]) {
				v = v*10 + int(c-'0')
			}
			data.put(v, 3*n+1)
		}
	case qrAlnum:
		for i := 0; i+1 < len(text); i += 2 {
			data.put(strings.IndexByte(qrAlnumChars, text[i])*45+strings.IndexByte(qrAlnumChars, text[i+1]), 11)
		}
		if len(text)%2 == 1 {
			data.put(strings.IndexByte(qrAlnumChars, text[len(text)-1]), 6)
		}
	default:
		for i := 0; i < len(text); i++ {
			data.put(int(text[i]), 8)
		}
	}
	// the version: the smallest with room
	ver := 0
	for v := 1; v <= 40; v++ {
		if cb := qrCountBits(mode, v); count < 1<<cb && 4+cb+len(data) <= qrDataCodewords(v, lv)*8 {
			ver = v
			break
		}
	}
	if ver == 0 {
		return nil, fmt.Errorf("barcode: a text of %d bytes is more than a QR code of level %c holds", len(text), "LMQH"[lv])
	}
	var all bits
	all.put([3]int{1, 2, 4}[mode], 4)
	all.put(count, qrCountBits(mode, ver))
	all = append(all, data...)
	capacity := qrDataCodewords(ver, lv) * 8
	all.put(0, min(4, capacity-len(all)))
	all.put(0, (8-len(all)%8)%8)
	for pad := 0xec; len(all) < capacity; pad ^= 0xec ^ 0x11 {
		all.put(pad, 8)
	}
	words := make([]byte, len(all)/8)
	for i, b := range all {
		if b {
			words[i/8] |= 0x80 >> (i % 8)
		}
	}
	q := &qrMatrix{ver: ver, size: ver*4 + 17}
	q.mods = make([]bool, q.size*q.size)
	q.fn = make([]bool, q.size*q.size)
	q.function()
	q.data(qrInterleave(words, ver, lv))
	if mask < 0 {
		best := -1
		for m := 0; m < 8; m++ {
			q.mask(m)
			q.format(lv, m)
			if p := q.penalty(); best < 0 || p < best {
				best, mask = p, m
			}
			q.mask(m) // (and off again)
		}
	}
	q.mask(mask & 7)
	q.format(lv, mask&7)
	if quiet < 0 {
		quiet = 4
	}
	s := &Symbol{Text: text, W: q.size + 2*quiet, H: q.size + 2*quiet}
	s.mods = make([]bool, s.W*s.H)
	for y := 0; y < q.size; y++ {
		copy(s.mods[(y+quiet)*s.W+quiet:], q.mods[y*q.size:(y+1)*q.size])
	}
	return s, nil
}

// qrInterleave adds the error correction to the codewords of data: the
// codewords of the blocks, then their error correction, one of each block
// after the other.
func qrInterleave(data []byte, ver, level int) []byte {
	blocks, ecc := int(qrBlocks[level][ver]), int(qrEccLen[level][ver])
	raw := qrRawModules(ver) / 8
	short := blocks - raw%blocks // the blocks that are a codeword shorter
	shortLen := raw/blocks - ecc // the codewords of data of a short block
	div := rsDivisor(ecc)
	var parts, eccs [][]byte
	for i, k := 0, 0; i < blocks; i++ {
		n := shortLen
		if i >= short {
			n++
		}
		parts = append(parts, data[k:k+n])
		eccs = append(eccs, rsRemainder(data[k:k+n], div))
		k += n
	}
	out := make([]byte, 0, raw)
	for i := 0; i <= shortLen; i++ {
		for _, p := range parts {
			if i < len(p) {
				out = append(out, p[i])
			}
		}
	}
	for i := 0; i < ecc; i++ {
		for _, e := range eccs {
			out = append(out, e[i])
		}
	}
	return out
}

// rsMul multiplies two elements of the field of the Reed–Solomon code of
// QR codes: GF(2^8) with the polynomial x^8 + x^4 + x^3 + x^2 + 1.
func rsMul(x, y byte) byte {
	z := 0
	for i := 7; i >= 0; i-- {
		z = z<<1 ^ z>>7*0x11d
		z ^= int(y>>i&1) * int(x)
	}
	return byte(z)
}

// rsDivisor returns the generator polynomial of a degree, without its
// leading 1.
func rsDivisor(degree int) []byte {
	out := make([]byte, degree)
	out[degree-1] = 1
	root := byte(1)
	for i := 0; i < degree; i++ {
		for j := range out {
			out[j] = rsMul(out[j], root)
			if j+1 < len(out) {
				out[j] ^= out[j+1]
			}
		}
		root = rsMul(root, 2)
	}
	return out
}

// rsRemainder returns the error correction codewords of data: the
// remainder of its division by the generator.
func rsRemainder(data, divisor []byte) []byte {
	out := make([]byte, len(divisor))
	for _, b := range data {
		f := b ^ out[0]
		copy(out, out[1:])
		out[len(out)-1] = 0
		for i := range out {
			out[i] ^= rsMul(divisor[i], f)
		}
	}
	return out
}

// qrMatrix is the modules of a QR code being made: fn marks those of the
// function patterns, which hold no data and take no mask.
type qrMatrix struct {
	ver, size int
	mods, fn  []bool
}

func (q *qrMatrix) set(x, y int, dark bool) {
	q.mods[y*q.size+x] = dark
	q.fn[y*q.size+x] = true
}

// function draws the function patterns: the timing patterns, the finder
// patterns with their separators, the alignment patterns, and the places
// of the format and version information.
func (q *qrMatrix) function() {
	n := q.size
	for i := 0; i < n; i++ {
		q.set(6, i, i%2 == 0)
		q.set(i, 6, i%2 == 0)
	}
	for _, c := range [][2]int{{3, 3}, {n - 4, 3}, {3, n - 4}} {
		for dy := -4; dy <= 4; dy++ {
			for dx := -4; dx <= 4; dx++ {
				x, y := c[0]+dx, c[1]+dy
				if x < 0 || y < 0 || x >= n || y >= n {
					continue
				}
				d := max(abs(dx), abs(dy))
				q.set(x, y, d != 2 && d != 4)
			}
		}
	}
	pos := qrAlignment(q.ver)
	for i, cx := range pos {
		for j, cy := range pos {
			last := len(pos) - 1
			if i == 0 && j == 0 || i == 0 && j == last || i == last && j == 0 {
				continue // the finder patterns are there
			}
			for dy := -2; dy <= 2; dy++ {
				for dx := -2; dx <= 2; dx++ {
					q.set(cx+dx, cy+dy, max(abs(dx), abs(dy)) != 1)
				}
			}
		}
	}
	q.format(0, 0)
	if q.ver >= 7 {
		rem := q.ver
		for i := 0; i < 12; i++ {
			rem = rem<<1 ^ rem>>11*0x1f25
		}
		v := q.ver<<12 | rem
		for i := 0; i < 18; i++ {
			a, b := n-11+i%3, i/3
			q.set(a, b, v>>i&1 != 0)
			q.set(b, a, v>>i&1 != 0)
		}
	}
}

// qrAlignment returns the places of the centres of the alignment patterns
// of a version, along a side.
func qrAlignment(ver int) []int {
	if ver == 1 {
		return nil
	}
	n := ver/7 + 2
	step := 26
	if ver != 32 {
		step = (ver*4 + n*2 + 1) / (n*2 - 2) * 2
	}
	out := make([]int, n)
	out[0] = 6
	for i, p := n-1, ver*4+10; i >= 1; i, p = i-1, p-step {
		out[i] = p
	}
	return out
}

// format writes the format information: the level and the mask, with
// their error correction, at its two places.
func (q *qrMatrix) format(level, mask int) {
	data := qrLevelBits[level]<<3 | mask
	rem := data
	for i := 0; i < 10; i++ {
		rem = rem<<1 ^ rem>>9*0x537
	}
	v := (data<<10 | rem) ^ 0x5412
	bit := func(i int) bool { return v>>i&1 != 0 }
	n := q.size
	for i := 0; i <= 5; i++ {
		q.set(8, i, bit(i))
	}
	q.set(8, 7, bit(6))
	q.set(8, 8, bit(7))
	q.set(7, 8, bit(8))
	for i := 9; i < 15; i++ {
		q.set(14-i, 8, bit(i))
	}
	for i := 0; i < 8; i++ {
		q.set(n-1-i, 8, bit(i))
	}
	for i := 8; i < 15; i++ {
		q.set(8, n-15+i, bit(i))
	}
	q.set(8, n-8, true) // always dark
}

// data puts the codewords on the modules that are free: in pairs of
// columns from the right, up and down in turn.
func (q *qrMatrix) data(words []byte) {
	n, i := q.size, 0
	for right := n - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5 // the column of the timing pattern
		}
		for v := 0; v < n; v++ {
			for j := 0; j < 2; j++ {
				x, y := right-j, v
				if (right+1)&2 == 0 {
					y = n - 1 - v
				}
				if !q.fn[y*n+x] && i < len(words)*8 {
					q.mods[y*n+x] = words[i>>3]>>(7-i&7)&1 != 0
					i++
				}
			}
		}
	}
}

// mask turns the modules of data over where a mask says to (and back,
// when it is applied again).
func (q *qrMatrix) mask(m int) {
	n := q.size
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if q.fn[y*n+x] {
				continue
			}
			var on bool
			switch m {
			case 0:
				on = (x+y)%2 == 0
			case 1:
				on = y%2 == 0
			case 2:
				on = x%3 == 0
			case 3:
				on = (x+y)%3 == 0
			case 4:
				on = (x/3+y/2)%2 == 0
			case 5:
				on = x*y%2+x*y%3 == 0
			case 6:
				on = (x*y%2+x*y%3)%2 == 0
			case 7:
				on = ((x+y)%2+x*y%3)%2 == 0
			}
			if on {
				q.mods[y*n+x] = !q.mods[y*n+x]
			}
		}
	}
}

// penalty rates the modules as the standard rates a mask: runs of five and
// more modules of a color, blocks of two by two, what looks like a finder
// pattern, and how far the dark modules are from half of all.
func (q *qrMatrix) penalty() int {
	n := q.size
	p := 0
	at := func(x, y int, turned bool) bool {
		if turned {
			x, y = y, x
		}
		return q.mods[y*n+x]
	}
	// the finder-like pattern, dark and light 1:1:3:1:1 with four light
	// modules at a side
	finder := [11]bool{true, false, true, true, true, false, true, false, false, false, false}
	for _, turned := range []bool{false, true} {
		for y := 0; y < n; y++ {
			run := 1
			for x := 1; x <= n; x++ {
				if x < n && at(x, y, turned) == at(x-1, y, turned) {
					run++
					continue
				}
				if run >= 5 {
					p += run - 2
				}
				run = 1
			}
			for x := 0; x+11 <= n; x++ {
				fwd, rev := true, true
				for k := 0; k < 11; k++ {
					v := at(x+k, y, turned)
					fwd = fwd && v == finder[k]
					rev = rev && v == finder[10-k]
				}
				if fwd {
					p += 40
				}
				if rev {
					p += 40
				}
			}
		}
	}
	dark := 0
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			v := q.mods[y*n+x]
			if v {
				dark++
			}
			if x+1 < n && y+1 < n && v == q.mods[y*n+x+1] && v == q.mods[(y+1)*n+x] && v == q.mods[(y+1)*n+x+1] {
				p += 3
			}
		}
	}
	total := n * n
	return p + ((abs(dark*20-total*10)+total-1)/total-1)*10
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

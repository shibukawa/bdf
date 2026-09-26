package jbig2

// bitReader reads bits most significant first, as Huffman coded data is packed.
type bitReader struct {
	data []byte
	pos  int // in bits
}

func (r *bitReader) readBit() (int, error) {
	i := r.pos >> 3
	if i >= len(r.data) {
		return 0, errTruncated
	}
	b := int(r.data[i]>>(7-uint(r.pos&7))) & 1
	r.pos++
	return b, nil
}

// readBits reads an n-bit unsigned number, n ≤ 32.
func (r *bitReader) readBits(n int) (uint32, error) {
	var v uint32
	for ; n > 0; n-- {
		b, err := r.readBit()
		if err != nil {
			return 0, err
		}
		v = v<<1 | uint32(b)
	}
	return v, nil
}

// align skips to the next byte boundary.
func (r *bitReader) align() { r.pos = (r.pos + 7) &^ 7 }

// bytes returns the n bytes at the (aligned) position and skips them.
func (r *bitReader) bytes(n int) ([]byte, error) {
	i := r.pos >> 3
	if n < 0 || n > len(r.data)-i {
		return nil, errTruncated
	}
	r.pos += n << 3
	return r.data[i : i+n], nil
}

// rest returns the bytes from the (aligned) position on.
func (r *bitReader) rest() []byte {
	i := r.pos >> 3
	if i >= len(r.data) {
		return nil
	}
	return r.data[i:]
}

// Kinds of Huffman table lines (B.2).
const (
	lineNormal = iota
	lineLower  // values below the range: low − offset
	lineUpper  // values above the range: low + offset
	lineOOB    // the out-of-band value
)

// huffLine is a line of a Huffman table: values low to low+2^rangeLen−1 are
// coded as a prefix of prefLen bits followed by rangeLen bits of offset.
type huffLine struct {
	prefLen  int
	rangeLen int
	low      int64
	kind     int
}

// huffTable is a Huffman table with prefix codes assigned as in B.3.
type huffTable struct {
	lines  []huffLine
	count  [33]int // codes of each length
	sorted []int   // line indices by code length, in table order
	maxLen int
}

func newHuffTable(lines []huffLine) (*huffTable, error) {
	t := &huffTable{lines: lines}
	for _, l := range lines {
		if l.prefLen < 0 || l.prefLen > 32 || l.rangeLen < 0 || l.rangeLen > 32 {
			return nil, errorf("invalid Huffman table")
		}
		if l.prefLen > 0 {
			t.count[l.prefLen]++
			t.maxLen = max(t.maxLen, l.prefLen)
		}
	}
	var start [34]int
	for n := 1; n <= 32; n++ {
		start[n+1] = start[n] + t.count[n]
	}
	t.sorted = make([]int, start[33])
	for i, l := range lines {
		if l.prefLen > 0 {
			t.sorted[start[l.prefLen]] = i
			start[l.prefLen]++
		}
	}
	return t, nil
}

// decode decodes a value; ok is false for the out-of-band value.
func (t *huffTable) decode(r *bitReader) (v int, ok bool, err error) {
	code, first, index := 0, 0, 0
	for n := 1; n <= t.maxLen; n++ {
		b, err := r.readBit()
		if err != nil {
			return 0, false, err
		}
		code = code<<1 | b
		c := t.count[n]
		if code >= first && code-first < c {
			return t.value(r, &t.lines[t.sorted[index+code-first]])
		}
		index += c
		first = (first + c) << 1
	}
	return 0, false, errorf("invalid Huffman code")
}

func (t *huffTable) value(r *bitReader, l *huffLine) (int, bool, error) {
	if l.kind == lineOOB {
		return 0, false, nil
	}
	off, err := r.readBits(l.rangeLen)
	if err != nil {
		return 0, false, err
	}
	v := l.low + int64(off)
	if l.kind == lineLower {
		v = l.low - int64(off)
	}
	if v < -(1<<31-1) || v > 1<<31-1 {
		return 0, false, errorf("Huffman coded value out of range")
	}
	return int(v), true, nil
}

// standardTables holds the standard Huffman tables B.1 to B.15 (index 1 to 15).
var standardTables [16]*huffTable

// Lines of the standard tables as {prefix length, range length, range low},
// with the lower range, upper range and OOB lines last.
var standardLines = [16][][3]int{
	1: {{1, 4, 0}, {2, 8, 16}, {3, 16, 272}, {3, 32, 65808}},
	2: {{1, 0, 0}, {2, 0, 1}, {3, 0, 2}, {4, 3, 3}, {5, 6, 11}, {6, 32, 75}, {6, 0, 0}},
	3: {{8, 8, -256}, {1, 0, 0}, {2, 0, 1}, {3, 0, 2}, {4, 3, 3}, {5, 6, 11}, {8, 32, -257}, {7, 32, 75}, {6, 0, 0}},
	4: {{1, 0, 1}, {2, 0, 2}, {3, 0, 3}, {4, 3, 4}, {5, 6, 12}, {5, 32, 76}},
	5: {{7, 8, -255}, {1, 0, 1}, {2, 0, 2}, {3, 0, 3}, {4, 3, 4}, {5, 6, 12}, {7, 32, -256}, {6, 32, 76}},
	6: {{5, 10, -2048}, {4, 9, -1024}, {4, 8, -512}, {4, 7, -256}, {5, 6, -128}, {5, 5, -64}, {4, 5, -32},
		{2, 7, 0}, {3, 7, 128}, {3, 8, 256}, {4, 9, 512}, {4, 10, 1024}, {6, 32, -2049}, {6, 32, 2048}},
	7: {{4, 9, -1024}, {3, 8, -512}, {4, 7, -256}, {5, 6, -128}, {5, 5, -64}, {4, 5, -32}, {4, 5, 0},
		{5, 5, 32}, {5, 6, 64}, {4, 7, 128}, {3, 8, 256}, {3, 9, 512}, {3, 10, 1024}, {5, 32, -1025}, {5, 32, 2048}},
	8: {{8, 3, -15}, {9, 1, -7}, {8, 1, -5}, {9, 0, -3}, {7, 0, -2}, {4, 0, -1}, {2, 1, 0}, {5, 0, 2},
		{6, 0, 3}, {3, 4, 4}, {6, 1, 20}, {4, 4, 22}, {4, 5, 38}, {5, 6, 70}, {5, 7, 134}, {6, 7, 262},
		{7, 8, 390}, {6, 10, 646}, {9, 32, -16}, {9, 32, 1670}, {2, 0, 0}},
	9: {{8, 4, -31}, {9, 2, -15}, {8, 2, -11}, {9, 1, -7}, {7, 1, -5}, {4, 1, -3}, {3, 1, -1}, {3, 1, 1},
		{5, 1, 3}, {6, 1, 5}, {3, 5, 7}, {6, 2, 39}, {4, 5, 43}, {4, 6, 75}, {5, 7, 139}, {5, 8, 267},
		{6, 8, 523}, {7, 9, 779}, {6, 11, 1291}, {9, 32, -32}, {9, 32, 3339}, {2, 0, 0}},
	10: {{7, 4, -21}, {8, 0, -5}, {7, 0, -4}, {5, 0, -3}, {2, 2, -2}, {5, 0, 2}, {6, 0, 3}, {7, 0, 4},
		{8, 0, 5}, {2, 6, 6}, {5, 5, 70}, {6, 5, 102}, {6, 6, 134}, {6, 7, 198}, {6, 8, 326}, {6, 9, 582},
		{6, 10, 1094}, {7, 11, 2118}, {8, 32, -22}, {8, 32, 4166}, {2, 0, 0}},
	11: {{1, 0, 1}, {2, 1, 2}, {4, 0, 4}, {4, 1, 5}, {5, 1, 7}, {5, 2, 9}, {6, 2, 13}, {7, 2, 17},
		{7, 3, 21}, {7, 4, 29}, {7, 5, 45}, {7, 6, 77}, {7, 32, 141}},
	12: {{1, 0, 1}, {2, 0, 2}, {3, 1, 3}, {5, 0, 5}, {5, 1, 6}, {6, 1, 8}, {7, 0, 10}, {7, 1, 11},
		{7, 2, 13}, {7, 3, 17}, {7, 4, 25}, {8, 5, 41}, {8, 32, 73}},
	13: {{1, 0, 1}, {3, 0, 2}, {4, 0, 3}, {5, 0, 4}, {4, 1, 5}, {3, 3, 7}, {6, 1, 15}, {6, 2, 17},
		{6, 3, 21}, {6, 4, 29}, {6, 5, 45}, {7, 6, 77}, {7, 32, 141}},
	14: {{3, 0, -2}, {3, 0, -1}, {1, 0, 0}, {3, 0, 1}, {3, 0, 2}},
	15: {{7, 4, -24}, {6, 2, -8}, {5, 1, -4}, {4, 0, -2}, {3, 0, -1}, {1, 0, 0}, {3, 0, 1}, {4, 0, 2},
		{5, 1, 3}, {6, 2, 5}, {7, 4, 9}, {7, 32, -25}, {7, 32, 25}},
}

// Which standard tables end with lower range, upper range and OOB lines.
var (
	standardHasLower = [16]bool{3: true, 5: true, 6: true, 7: true, 8: true, 9: true, 10: true, 15: true}
	standardHasOOB   = [16]bool{2: true, 3: true, 8: true, 9: true, 10: true}
)

func init() {
	for i := 1; i <= 15; i++ {
		src := standardLines[i]
		lines := make([]huffLine, len(src))
		for j, l := range src {
			lines[j] = huffLine{prefLen: l[0], rangeLen: l[1], low: int64(l[2])}
		}
		n := len(lines)
		if standardHasOOB[i] {
			n--
			lines[n].kind = lineOOB
		}
		if i != 14 { // B.14 has no range lines
			lines[n-1].kind = lineUpper
			if standardHasLower[i] {
				lines[n-2].kind = lineLower
			}
		}
		t, err := newHuffTable(lines)
		if err != nil {
			panic(err)
		}
		standardTables[i] = t
	}
}

// parseTableSegment parses a code table segment (7.4.13, B.2).
func parseTableSegment(data []byte) (*huffTable, error) {
	r := reader{b: data}
	flags := r.u8()
	low, high := int64(int32(r.u32())), int64(int32(r.u32()))
	if r.err != nil {
		return nil, r.err
	}
	ps, rs := int(flags>>1&7)+1, int(flags>>4&7)+1
	br := bitReader{data: r.rest()}
	var lines []huffLine
	read := func(n int) int {
		v, e := br.readBits(n)
		if e != nil && r.err == nil {
			r.err = e
		}
		return int(v)
	}
	for cur := low; cur < high && r.err == nil; {
		l := huffLine{prefLen: read(ps), rangeLen: read(rs), low: cur}
		if l.rangeLen > 32 {
			return nil, errorf("invalid Huffman table")
		}
		lines = append(lines, l)
		cur += 1 << l.rangeLen
	}
	lines = append(lines,
		huffLine{prefLen: read(ps), rangeLen: 32, low: low - 1, kind: lineLower},
		huffLine{prefLen: read(ps), rangeLen: 32, low: high, kind: lineUpper})
	if flags&1 != 0 {
		lines = append(lines, huffLine{prefLen: read(ps), kind: lineOOB})
	}
	if r.err != nil {
		return nil, r.err
	}
	return newHuffTable(lines)
}

// tableSelector picks standard or custom Huffman tables for a segment; custom
// tables come from the referred-to table segments in order.
type tableSelector struct {
	custom []*huffTable
	err    error
}

// pick returns standard table std[sel], or the next custom table if sel is
// custom; std has 0 for invalid selections.
func (s *tableSelector) pick(sel int, custom int, std ...int) *huffTable {
	if s.err != nil {
		return nil
	}
	if sel == custom {
		if len(s.custom) == 0 {
			s.err = errorf("missing custom Huffman table")
			return nil
		}
		t := s.custom[0]
		s.custom = s.custom[1:]
		return t
	}
	if sel >= len(std) || std[sel] == 0 {
		s.err = errorf("invalid Huffman table selection")
		return nil
	}
	return standardTables[std[sel]]
}

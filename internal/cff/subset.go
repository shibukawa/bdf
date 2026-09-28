package cff

import (
	"encoding/binary"
	"errors"
	"math"
)

// Subset returns a CFF font program that holds only the glyphs in keep,
// glyph 0 and the components of seac accents, renumbered in their original
// order; order[newGID] is the old GID. Subroutines the kept glyphs do not call
// become a bare return (their numbers must not change). The charset keeps
// each glyph's name or CID, FDSelect is rebuilt, and a custom encoding is
// replaced by StandardEncoding: its codes point at old GIDs, and browsers
// map characters through the cmap anyway.
//
// Dropping glyphs instead of emptying them matters for CFF: its CharStrings
// INDEX has an offset per glyph, and in a CJK font those offsets alone are
// far larger than the few glyphs a document uses (TrueType's loca has no
// such cost in WOFF2, which rebuilds it).
//
// When a charstring cannot be followed (arithmetic operators, a subroutine
// number out of range) every subroutine is kept and only the glyphs are dropped.
// Charstrings that run longer than a font's size explains (subroutines that
// call subroutines many times over) are an error: ErrTooLong.
func Subset(data []byte, keep map[int]bool) (out []byte, order []int, err error) {
	return subset(data, keep, false)
}

// Prune is Subset that keeps the glyph numbers: the glyphs not in keep (nor
// glyph 0 or a component of a kept seac accent) stay in the font, empty,
// so that tables which refer to glyphs by number (GSUB, GPOS, hmtx) stay
// valid. The subroutines only they call become a bare return.
func Prune(data []byte, keep map[int]bool) ([]byte, error) {
	out, _, err := subset(data, keep, true)
	return out, err
}

func subset(data []byte, keep map[int]bool, keepNumbers bool) (out []byte, order []int, err error) {
	cf, err := Parse(data)
	if err != nil {
		return nil, nil, err
	}
	hdrSize := int(data[2])
	if hdrSize < 4 || hdrSize > len(data) {
		return nil, nil, errors.New("cff: bad header")
	}
	names, pos, err := ReadIndex(data, hdrSize)
	if err != nil || len(names) != 1 {
		return nil, nil, errors.New("cff: only single-font CFF is subset")
	}
	nameIdx := data[hdrSize:pos]
	topDicts, pos, err := ReadIndex(data, pos)
	if err != nil || len(topDicts) != 1 {
		return nil, nil, errors.New("cff: no top dict")
	}
	strStart := pos
	if _, pos, err = ReadIndex(data, pos); err != nil {
		return nil, nil, err
	}
	strIdx := data[strStart:pos]
	gsubrs, _, err := ReadIndex(data, pos)
	if err != nil {
		return nil, nil, err
	}
	top, err := DictEntries(topDicts[0])
	if err != nil {
		return nil, nil, err
	}
	if cf.IsCID {
		// A CID-keyed font's Private DICTs hang off its Font DICTs; a stray
		// Private in the Top DICT would point into the old layout.
		kept := top[:0:0]
		for _, e := range top {
			if e.Op != 18 {
				kept = append(kept, e)
			}
		}
		top = kept
	}
	topVal := func(op int) []float64 {
		for _, e := range top {
			if e.Op == op {
				return e.Args
			}
		}
		return nil
	}
	off := func(op, i int) (int, bool) {
		v := topVal(op)
		if len(v) <= i || v[i] < 0 || operand(v[i]) > len(data) {
			return 0, false
		}
		return operand(v[i]), true
	}
	csOff, ok := off(17, 0)
	if !ok {
		return nil, nil, errors.New("cff: no CharStrings")
	}
	charStrings, _, err := ReadIndex(data, csOff)
	if err != nil {
		return nil, nil, err
	}
	n := len(charStrings)
	if n != cf.NumGlyphs || len(cf.Charset) != n {
		return nil, nil, errors.New("cff: inconsistent glyph count")
	}
	// The renumbered font always has a custom charset.
	if topVal(15) == nil {
		top = append(top, DictEntry{Op: 15, Args: []float64{0}})
	}
	// Private DICTs with their local subroutines: one for a name-keyed font,
	// one per Font DICT for a CID-keyed font.
	type private struct {
		dict  []DictEntry
		subrs *subrSet
	}
	readPrivate := func(sizeOff []float64) (*private, error) {
		if len(sizeOff) != 2 {
			return nil, errors.New("cff: bad Private operands")
		}
		size, o := operand(sizeOff[0]), operand(sizeOff[1])
		if size < 0 || o < 0 || o+size > len(data) {
			return nil, errors.New("cff: Private out of range")
		}
		d, err := DictEntries(data[o : o+size])
		if err != nil {
			return nil, err
		}
		p := &private{dict: d, subrs: &subrSet{}}
		for _, e := range d {
			if e.Op == 19 && len(e.Args) == 1 {
				items, _, err := ReadIndex(data, o+operand(e.Args[0]))
				if err != nil {
					return nil, err
				}
				p.subrs = newSubrSet(items)
			}
		}
		return p, nil
	}
	var privates []*private
	var fdDicts [][]DictEntry
	fdOf := func(int) int { return 0 }
	if cf.IsCID {
		fdaOff, ok1 := off(1236, 0)
		fdsOff, ok2 := off(1237, 0)
		if !ok1 || !ok2 {
			return nil, nil, errors.New("cff: CID font without FDArray/FDSelect")
		}
		fds, _, err := ReadIndex(data, fdaOff)
		if err != nil || len(fds) == 0 || len(fds) > 255 {
			return nil, nil, errors.New("cff: bad FDArray")
		}
		for _, fd := range fds {
			d, err := DictEntries(fd)
			if err != nil {
				return nil, nil, err
			}
			var p *private
			for _, e := range d {
				if e.Op == 18 {
					if p, err = readPrivate(e.Args); err != nil {
						return nil, nil, err
					}
				}
			}
			if p == nil {
				p = &private{subrs: &subrSet{}}
			}
			fdDicts = append(fdDicts, d)
			privates = append(privates, p)
		}
		sel, err := fdSelect(data, fdsOff, n)
		if err != nil {
			return nil, nil, err
		}
		fdOf = func(gid int) int {
			if fd := sel[gid]; fd < len(privates) {
				return fd
			}
			return 0
		}
	} else if v := topVal(18); v != nil {
		p, err := readPrivate(v)
		if err != nil {
			return nil, nil, err
		}
		privates = []*private{p}
	} else {
		privates = []*private{{subrs: &subrSet{}}}
	}

	// Walk the kept glyphs, following subroutine calls and seac components.
	// Every byte of a charstring or subroutine is run once for each glyph
	// that reaches it, which subroutines shared by many glyphs make a few
	// times the size of the font, not more.
	left := 32*len(data) + 1<<16
	global := newSubrSet(gsubrs)
	kept := map[int]bool{0: true}
	queue := []int{0}
	for gid := range keep {
		if gid > 0 && gid < n && !kept[gid] {
			kept[gid] = true
			queue = append(queue, gid)
		}
	}
	followed := true
	for len(queue) > 0 {
		gid := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		st := &t2State{seac: [2]int{-1, -1}, left: &left}
		if _, err := st.exec(charStrings[gid], privates[fdOf(gid)].subrs, global, 0); err == ErrTooLong {
			return nil, nil, err
		} else if err != nil {
			followed = false
		}
		if st.seac[0] >= 0 && !cf.IsCID {
			for _, code := range st.seac {
				if code >= 0 && code < 256 {
					if g, ok := cf.NameToGID[StandardEncoding[code]]; ok && !kept[g] {
						kept[g] = true
						queue = append(queue, g)
					}
				}
			}
		}
	}

	// Renumber: kept glyphs in their original order (all of them, the
	// others emptied, when the numbers are kept).
	for gid := 0; gid < n; gid++ {
		if kept[gid] || keepNumbers {
			order = append(order, gid)
		}
	}
	newCS := make([][]byte, len(order))
	charset := []byte{0} // format 0: SID (or CID) of every glyph after .notdef
	for i, gid := range order {
		newCS[i] = charStrings[gid]
		if !kept[gid] {
			newCS[i] = []byte{14} // endchar
		}
		if i > 0 {
			charset = binary.BigEndian.AppendUint16(charset, uint16(cf.Charset[gid]))
		}
	}
	var fdSelect []byte // format 3: ranges of glyphs sharing a Font DICT
	if cf.IsCID {
		var ranges []byte
		nRanges := 0
		for i, gid := range order {
			if fd := byte(fdOf(gid)); i == 0 || ranges[len(ranges)-1] != fd {
				ranges = binary.BigEndian.AppendUint16(ranges, uint16(i))
				ranges = append(ranges, fd)
				nRanges++
			}
		}
		fdSelect = binary.BigEndian.AppendUint16([]byte{3}, uint16(nRanges))
		fdSelect = append(fdSelect, ranges...)
		fdSelect = binary.BigEndian.AppendUint16(fdSelect, uint16(len(order)))
	}
	prune := func(s *subrSet) [][]byte {
		if followed {
			return s.pruned()
		}
		return s.items
	}
	gsubrIdx := Index(prune(global))
	csIdx := Index(newCS)
	privBlocks := make([][]byte, len(privates)) // Private DICT followed by its Subrs INDEX
	privSizes := make([]int, len(privates))
	for i, p := range privates {
		if p.dict == nil {
			continue
		}
		// Subrs is relative to the Private DICT: the INDEX follows the DICT directly.
		hasSubrs := false
		for _, e := range p.dict {
			hasSubrs = hasSubrs || e.Op == 19
		}
		size := len(EncodeDict(p.dict, func(e DictEntry) ([]int, bool) { return nil, e.Op == 19 }))
		d := EncodeDict(p.dict, func(e DictEntry) ([]int, bool) {
			if e.Op == 19 {
				return []int{size}, true
			}
			return nil, false
		})
		privBlocks[i] = d
		if hasSubrs {
			privBlocks[i] = append(d, Index(prune(p.subrs))...)
		}
		privSizes[i] = len(d)
	}

	// Layout: header, Name, Top DICT, String, Global Subr INDEXes, then
	// charset, encoding, FDSelect, CharStrings, FDArray and Private blocks.
	// Offset operands are written as 5-byte integers, so the DICT sizes are
	// known before the offsets are.
	offsets := map[int][]int{}
	encodeTop := func() []byte {
		return EncodeDict(top, func(e DictEntry) ([]int, bool) {
			switch e.Op {
			case 16:
				if len(e.Args) == 1 && e.Args[0] <= 1 {
					return nil, false // predefined encoding: codes map to names, still valid
				}
				return []int{0}, true
			case 15, 17, 18, 1236, 1237:
			default:
				return nil, false
			}
			if v, ok := offsets[e.Op]; ok {
				return v, true
			}
			return make([]int, len(e.Args)), true
		})
	}
	encodeFDs := func(privOff []int) []byte {
		items := make([][]byte, len(fdDicts))
		for i, d := range fdDicts {
			items[i] = EncodeDict(d, func(e DictEntry) ([]int, bool) {
				if e.Op == 18 {
					return []int{privSizes[i], privOff[i]}, true
				}
				return nil, false
			})
		}
		return Index(items)
	}
	topLen := len(Index([][]byte{encodeTop()}))
	p := hdrSize + len(nameIdx) + topLen + len(strIdx) + len(gsubrIdx)
	offsets[15] = []int{p}
	p += len(charset)
	if fdSelect != nil {
		offsets[1237] = []int{p}
		p += len(fdSelect)
	}
	offsets[17] = []int{p}
	p += len(csIdx)
	privOff := make([]int, len(privates))
	var fdArray []byte
	if cf.IsCID {
		offsets[1236] = []int{p}
		p += len(encodeFDs(privOff))
	}
	for i, b := range privBlocks {
		privOff[i] = p
		p += len(b)
	}
	if !cf.IsCID && privates[0].dict != nil {
		offsets[18] = []int{privSizes[0], privOff[0]}
	}
	if cf.IsCID {
		fdArray = encodeFDs(privOff)
	}
	out = make([]byte, 0, p)
	out = append(out, data[:hdrSize]...)
	out = append(out, nameIdx...)
	out = append(out, Index([][]byte{encodeTop()})...)
	out = append(out, strIdx...)
	out = append(out, gsubrIdx...)
	out = append(out, charset...)
	out = append(out, fdSelect...)
	out = append(out, csIdx...)
	out = append(out, fdArray...)
	for _, b := range privBlocks {
		out = append(out, b...)
	}
	if len(out) != p {
		return nil, nil, errors.New("cff: layout mismatch")
	}
	if check, err := Parse(out); err != nil || check.NumGlyphs != len(order) || check.IsCID != cf.IsCID {
		return nil, nil, errors.New("cff: subset does not parse")
	}
	return out, order, nil
}

// DictEntry is one operator of a DICT with its operands, kept as written.
type DictEntry struct {
	Op   int // 12 x is 1200+x
	Args []float64
	Raw  []byte // operand bytes
}

func DictEntries(b []byte) ([]DictEntry, error) {
	var out []DictEntry
	start := 0
	for i := 0; i < len(b); {
		c := int(b[i])
		switch {
		case c <= 21:
			op := c
			end := i + 1
			if c == 12 {
				if end >= len(b) {
					return nil, errors.New("cff: truncated DICT")
				}
				op = 1200 + int(b[end])
				end++
			}
			raw := b[start:i]
			args := ParseDict(append(append([]byte(nil), raw...), 0))[0] // operands only
			out = append(out, DictEntry{Op: op, Args: args, Raw: raw})
			i = end
			start = end
		case c == 28:
			i += 3
		case c == 29:
			i += 5
		case c == 30:
			i++
			for i < len(b) {
				v := b[i]
				i++
				if v&0x0f == 0x0f || v&0xf0 == 0xf0 {
					break
				}
			}
		case c >= 32 && c <= 246:
			i++
		case c >= 247 && c <= 254:
			i += 2
		default:
			return nil, errors.New("cff: reserved byte in DICT")
		}
		if i > len(b) {
			return nil, errors.New("cff: truncated DICT")
		}
	}
	return out, nil
}

// cffEncodeDict writes entries back; replace returns new integer operands
// for an entry (written as 5-byte integers) or false to keep it as it was.
func EncodeDict(entries []DictEntry, replace func(DictEntry) ([]int, bool)) []byte {
	var out []byte
	for _, e := range entries {
		if v, ok := replace(e); ok {
			if v == nil {
				v = []int{0}
			}
			for _, x := range v {
				out = append(out, 29)
				out = binary.BigEndian.AppendUint32(out, uint32(int32(x)))
			}
		} else {
			out = append(out, e.Raw...)
		}
		if e.Op >= 1200 {
			out = append(out, 12, byte(e.Op-1200))
		} else {
			out = append(out, byte(e.Op))
		}
	}
	return out
}

// cffIndex writes an INDEX.
func Index(items [][]byte) []byte {
	if len(items) == 0 {
		return []byte{0, 0}
	}
	total := 1
	for _, it := range items {
		total += len(it)
	}
	offSize := 1
	for total >= 1<<(8*offSize) {
		offSize++
	}
	out := make([]byte, 0, 3+(len(items)+1)*offSize+total)
	out = binary.BigEndian.AppendUint16(out, uint16(len(items)))
	out = append(out, byte(offSize))
	o := 1
	put := func(v int) {
		for k := offSize - 1; k >= 0; k-- {
			out = append(out, byte(v>>(8*k)))
		}
	}
	put(o)
	for _, it := range items {
		o += len(it)
		put(o)
	}
	for _, it := range items {
		out = append(out, it...)
	}
	return out
}

// cffFDSelect returns the font DICT index of every glyph.
func fdSelect(data []byte, off, numGlyphs int) ([]int, error) {
	if off >= len(data) {
		return nil, errors.New("cff: FDSelect out of range")
	}
	sel := make([]int, numGlyphs)
	switch data[off] {
	case 0:
		end := off + 1 + numGlyphs
		if end > len(data) {
			return nil, errors.New("cff: truncated FDSelect")
		}
		for g := range sel {
			sel[g] = int(data[off+1+g])
		}
		return sel, nil
	case 3:
		nRanges := int(be16(data, off+1))
		end := off + 3 + 3*nRanges + 2
		if end > len(data) {
			return nil, errors.New("cff: truncated FDSelect")
		}
		for r := 0; r < nRanges; r++ {
			rec := off + 3 + 3*r
			first, fd, next := int(be16(data, rec)), int(data[rec+2]), int(be16(data, rec+3))
			for g := first; g < next && g < numGlyphs; g++ {
				sel[g] = fd
			}
		}
		return sel, nil
	}
	return nil, errors.New("cff: unknown FDSelect format")
}

// subrSet is a subroutine INDEX with the entries charstrings were seen to call.
type subrSet struct {
	items [][]byte
	used  []bool
	bias  int
}

func newSubrSet(items [][]byte) *subrSet {
	s := &subrSet{items: items, used: make([]bool, len(items)), bias: 107}
	switch {
	case len(items) >= 33900:
		s.bias = 32768
	case len(items) >= 1240:
		s.bias = 1131
	}
	return s
}

// pruned returns the INDEX items with unused subroutines replaced by return.
func (s *subrSet) pruned() [][]byte {
	out := make([][]byte, len(s.items))
	for i, it := range s.items {
		if s.used[i] {
			out[i] = it
		} else {
			out[i] = []byte{11}
		}
	}
	return out
}

// t2State is the part of the Type 2 charstring machine that decides which
// bytes are operators: the operand stack, the stem count (hintmask length)
// and the seac components of an accented endchar.
type t2State struct {
	stack  []float64
	nStems int
	seac   [2]int
	// left is how many operands and operators may still be run, shared by
	// the glyphs of a font; nil for no limit.
	left *int
}

var errT2 = errors.New("cff: charstring cannot be followed")

// ErrTooLong is returned by Subset and Prune for a font whose charstrings
// run too long to follow, or leave too many operands on the stack.
var ErrTooLong = errors.New("cff: charstrings run too long")

// maxStack bounds the operands on the stack: sixteen times the 513 of
// CFF2, the larger of the limits of the formats. A subroutine call keeps
// the stack, so that calls in a row would grow it without end.
const maxStack = 16 * 513

// exec interprets cs, marking the subroutines it calls. It reports true when
// the glyph ended (endchar) and an error for charstrings it cannot follow.
func (st *t2State) exec(cs []byte, local, global *subrSet, depth int) (bool, error) {
	if depth > 10 {
		return false, errT2
	}
	for i := 0; i < len(cs); {
		if st.left != nil {
			if *st.left--; *st.left < 0 {
				return false, ErrTooLong
			}
		}
		if len(st.stack) > maxStack {
			return false, ErrTooLong
		}
		b := cs[i]
		switch {
		case b == 28:
			if i+3 > len(cs) {
				return false, errT2
			}
			st.stack = append(st.stack, float64(int16(be16(cs, i+1))))
			i += 3
			continue
		case b >= 32 && b <= 246:
			st.stack = append(st.stack, float64(int(b)-139))
			i++
			continue
		case b >= 247 && b <= 254:
			if i+2 > len(cs) {
				return false, errT2
			}
			if b <= 250 {
				st.stack = append(st.stack, float64((int(b)-247)*256+int(cs[i+1])+108))
			} else {
				st.stack = append(st.stack, float64(-(int(b)-251)*256-int(cs[i+1])-108))
			}
			i += 2
			continue
		case b == 255:
			if i+5 > len(cs) {
				return false, errT2
			}
			st.stack = append(st.stack, float64(int32(be32(cs, i+1)))/65536)
			i += 5
			continue
		}
		op := int(b)
		i++
		if b == 12 {
			if i >= len(cs) {
				return false, errT2
			}
			op = 1200 + int(cs[i])
			i++
		}
		switch op {
		case 1, 3, 18, 23: // hstem, vstem, hstemhm, vstemhm
			st.nStems += len(st.stack) / 2
		case 19, 20: // hintmask, cntrmask (operands are an implied vstem)
			st.nStems += len(st.stack) / 2
			i += (st.nStems + 7) / 8
		case 10, 29: // callsubr, callgsubr
			set := local
			if op == 29 {
				set = global
			}
			if len(st.stack) == 0 {
				return false, errT2
			}
			v := st.stack[len(st.stack)-1]
			st.stack = st.stack[:len(st.stack)-1]
			k := int(v) + set.bias
			if v != math.Trunc(v) || k < 0 || k >= len(set.items) {
				return false, errT2
			}
			set.used[k] = true
			done, err := st.exec(set.items[k], local, global, depth+1)
			if err != nil || done {
				return done, err
			}
			continue
		case 11: // return
			return false, nil
		case 14: // endchar; with four operands it is seac
			if k := len(st.stack); k >= 4 {
				st.seac = [2]int{int(st.stack[k-2]), int(st.stack[k-1])}
			}
			return true, nil
		case 1200, 1234, 1235, 1236, 1237: // dotsection, hflex, flex, hflex1, flex1
		default:
			if op >= 1200 {
				// Arithmetic and storage operators can compute subroutine numbers.
				return false, errT2
			}
		}
		st.stack = st.stack[:0]
	}
	return false, nil
}

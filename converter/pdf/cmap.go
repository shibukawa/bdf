package pdf

import (
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf/converter/internal/cjkcmap"
)

// cmap is a parsed CMap: byte-code → CID (encoding CMaps) or byte-code →
// Unicode string (ToUnicode CMaps). An encoding CMap falls back on its
// parent (usecmap) for codes it does not map.
type cmap struct {
	ranges   []codespace // sorted by byte length
	single   map[uint32]uint32
	cidRange []cidRange // sorted by lo once parsed
	overlap  bool       // cidRange has overlapping ranges: scan them in order
	unicode  map[uint32]string
	uniRange []bfRange
	identity bool // unmapped codes are CIDs (Identity-H/V, or usecmap of them)
	vertical bool
	wmodeSet bool // the CMap defines /WMode (else it takes its parent's)
	seen1    bool
	seen2    bool

	useName string        // the usecmap operand
	parent  *cmap         // an embedded parent CMap (/UseCMap stream)
	pre     *cjkcmap.CMap // a predefined CMap: the encoding itself or the parent
}

type codespace struct {
	nbytes  int
	low, hi uint32
}

type cidRange struct {
	nbytes int
	lo, hi uint32
	cid    uint32
}

type bfRange struct {
	lo, hi uint32
	dst    []rune // first code's mapping; incremented for the rest
}

func identityCMap() *cmap {
	return &cmap{identity: true, ranges: []codespace{{2, 0, 0xffff}}, single: map[uint32]uint32{}, unicode: map[uint32]string{}}
}

func newCMap() *cmap {
	return &cmap{single: map[uint32]uint32{}, unicode: map[uint32]string{}}
}

// glyphCode is one code extracted from a string operand.
type glyphCode struct {
	code   uint32
	nbytes int
	cid    uint32
}

// decode splits s into codes according to the codespace ranges.
func (cm *cmap) decode(s []byte) []glyphCode {
	var out []glyphCode
	i := 0
	for i < len(s) {
		matched := false
		var code uint32
		for n := 1; n <= 4 && i+n <= len(s); n++ {
			code = code<<8 | uint32(s[i+n-1])
			for _, r := range cm.ranges {
				if r.nbytes == n && code >= r.low && code <= r.hi {
					out = append(out, glyphCode{code: code, nbytes: n, cid: cm.cid(code)})
					i += n
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			// Use the shortest codespace length (or 1) for undefined codes.
			n := 1
			if len(cm.ranges) > 0 {
				n = cm.ranges[0].nbytes
				for _, r := range cm.ranges {
					if r.nbytes < n {
						n = r.nbytes
					}
				}
			}
			if i+n > len(s) {
				n = len(s) - i
			}
			var c uint32
			for k := 0; k < n; k++ {
				c = c<<8 | uint32(s[i+k])
			}
			out = append(out, glyphCode{code: c, nbytes: n, cid: cm.cid(c)})
			i += n
		}
	}
	return out
}

func (cm *cmap) cid(code uint32) uint32 {
	cid, _ := cm.lookup(code, 0)
	return cid
}

// lookup maps a code with the CMap's own entries, then its parents.
func (cm *cmap) lookup(code uint32, depth int) (uint32, bool) {
	if v, ok := cm.single[code]; ok {
		return v, true
	}
	if cm.overlap {
		for _, r := range cm.cidRange {
			if code >= r.lo && code <= r.hi {
				return r.cid + (code - r.lo), true
			}
		}
	} else if i := sort.Search(len(cm.cidRange), func(i int) bool { return cm.cidRange[i].hi >= code }); i < len(cm.cidRange) && cm.cidRange[i].lo <= code {
		r := cm.cidRange[i]
		return r.cid + (code - r.lo), true
	}
	if cm.parent != nil && depth < 8 {
		if v, ok := cm.parent.lookup(code, depth+1); ok {
			return v, true
		}
	}
	if cm.pre != nil {
		if v, ok := cm.pre.CID(code); ok {
			return v, true
		}
	}
	if cm.identity {
		return code, true
	}
	return 0, false
}

// predefinedCMap returns an encoding CMap for a predefined name: Identity-H
// and Identity-V, or one of Adobe's CJK CMaps; nil for unknown names.
func predefinedCMap(name string) *cmap {
	switch name {
	case "Identity-H":
		return identityCMap()
	case "Identity-V":
		cm := identityCMap()
		cm.vertical = true
		return cm
	}
	pm := cjkcmap.Lookup(name)
	if pm == nil {
		return nil
	}
	cm := newCMap()
	cm.pre = pm
	cm.vertical = pm.Vertical
	cm.inherit(pm.Codespace)
	return cm
}

// finish resolves the parent named by usecmap (or given as a stream) and
// settles the code space ranges.
func (cm *cmap) finish(parent *cmap) {
	if parent == nil && cm.useName != "" {
		parent = predefinedCMap(cm.useName)
	}
	if parent != nil {
		if parent.pre != nil && len(parent.single) == 0 && len(parent.cidRange) == 0 {
			cm.pre = parent.pre // a predefined CMap: look it up directly
		} else if parent.identity && len(parent.single) == 0 && len(parent.cidRange) == 0 {
			cm.identity = true
		} else {
			cm.parent = parent
		}
		if len(cm.ranges) == 0 {
			cm.ranges = append(cm.ranges, parent.ranges...)
		}
		if !cm.wmodeSet {
			cm.vertical = parent.vertical
		}
	}
	sort.Slice(cm.cidRange, func(i, j int) bool { return cm.cidRange[i].lo < cm.cidRange[j].lo })
	for i := 1; i < len(cm.cidRange); i++ {
		if cm.cidRange[i].lo <= cm.cidRange[i-1].hi {
			cm.overlap = true
		}
	}
	if len(cm.ranges) == 0 {
		// Infer from the code lengths seen, defaulting to 2 bytes for CID CMaps.
		if cm.seen2 || (!cm.seen1 && len(cm.unicode) == 0) {
			cm.ranges = append(cm.ranges, codespace{2, 0, 0xffff})
		}
		if cm.seen1 {
			cm.ranges = append(cm.ranges, codespace{1, 0, 0xff})
		}
	}
}

func (cm *cmap) inherit(cs []cjkcmap.Codespace) {
	for _, r := range cs {
		cm.ranges = append(cm.ranges, codespace{nbytes: r.Bytes, low: r.Lo, hi: r.Hi})
	}
}

// toUnicode returns the Unicode string for a code, or "" when unmapped.
func (cm *cmap) toUnicode(code uint32) string {
	if cm == nil {
		return ""
	}
	if s, ok := cm.unicode[code]; ok {
		return s
	}
	for _, r := range cm.uniRange {
		if code >= r.lo && code <= r.hi && len(r.dst) > 0 {
			d := append([]rune(nil), r.dst...)
			d[len(d)-1] += rune(code - r.lo)
			return string(d)
		}
	}
	return ""
}

func bytesToCode(b []byte) uint32 {
	var c uint32
	for _, x := range b {
		c = c<<8 | uint32(x)
	}
	return c
}

func utf16Runes(b []byte) []rune {
	var out []rune
	for i := 0; i+1 < len(b); i += 2 {
		u := rune(b[i])<<8 | rune(b[i+1])
		if u >= 0xd800 && u < 0xdc00 && i+3 < len(b) {
			lo := rune(b[i+2])<<8 | rune(b[i+3])
			u = 0x10000 + (u-0xd800)<<10 + (lo - 0xdc00)
			i += 2
		}
		out = append(out, u)
	}
	if len(b) == 1 {
		out = append(out, rune(b[0]))
	}
	return out
}

// parseCMap parses an embedded CMap or ToUnicode stream.
func parseCMap(data []byte) *cmap { return parseCMapParent(data, nil) }

// parseCMapParent parses an embedded CMap whose parent (/UseCMap) may be
// another embedded one.
func parseCMapParent(data []byte, parent *cmap) *cmap {
	cm := newCMap()
	l := &lexer{b: data}
	var stack []types.Object
	for {
		kind, obj, op := l.next()
		if kind == tokEOF {
			break
		}
		if kind == tokOperand {
			stack = append(stack, obj)
			if len(stack) > 64 {
				stack = stack[len(stack)-64:]
			}
			continue
		}
		switch op {
		case "begincodespacerange":
			items := readUntil(l, "endcodespacerange")
			for i := 0; i+1 < len(items); i += 2 {
				lo, hi := literalBytes(items[i]), literalBytes(items[i+1])
				if len(lo) == 0 {
					continue
				}
				cm.ranges = append(cm.ranges, codespace{nbytes: len(lo), low: bytesToCode(lo), hi: bytesToCode(hi)})
			}
		case "begincidrange":
			items := readUntil(l, "endcidrange")
			for i := 0; i+2 < len(items); i += 3 {
				lo, hi := literalBytes(items[i]), literalBytes(items[i+1])
				cid, _ := numOf(items[i+2])
				cm.cidRange = append(cm.cidRange, cidRange{nbytes: len(lo), lo: bytesToCode(lo), hi: bytesToCode(hi), cid: uint32(cid)})
				cm.noteLen(len(lo))
			}
		case "begincidchar":
			items := readUntil(l, "endcidchar")
			for i := 0; i+1 < len(items); i += 2 {
				code := literalBytes(items[i])
				cid, _ := numOf(items[i+1])
				cm.single[bytesToCode(code)] = uint32(cid)
				cm.noteLen(len(code))
			}
		case "beginbfchar":
			items := readUntil(l, "endbfchar")
			for i := 0; i+1 < len(items); i += 2 {
				code := literalBytes(items[i])
				var dst string
				if n, ok := items[i+1].(types.Name); ok {
					if r, ok := glyphNameToRune(n.Value()); ok {
						dst = string(r)
					}
				} else {
					dst = string(utf16Runes(literalBytes(items[i+1])))
				}
				cm.unicode[bytesToCode(code)] = dst
				cm.noteLen(len(code))
			}
		case "beginbfrange":
			items := readUntil(l, "endbfrange")
			for i := 0; i+2 < len(items); i += 3 {
				lo, hi := bytesToCode(literalBytes(items[i])), bytesToCode(literalBytes(items[i+1]))
				cm.noteLen(len(literalBytes(items[i])))
				switch dst := items[i+2].(type) {
				case types.Array:
					for k, e := range dst {
						cm.unicode[lo+uint32(k)] = string(utf16Runes(literalBytes(e)))
					}
				default:
					cm.uniRange = append(cm.uniRange, bfRange{lo: lo, hi: hi, dst: utf16Runes(literalBytes(dst))})
				}
			}
		case "usecmap":
			if len(stack) > 0 {
				if n, ok := stack[len(stack)-1].(types.Name); ok {
					cm.useName = n.Value()
				}
			}
		case "def":
			if len(stack) >= 2 {
				if n, ok := stack[len(stack)-2].(types.Name); ok && n.Value() == "WMode" {
					if v, ok := numOf(stack[len(stack)-1]); ok {
						cm.vertical = v == 1
						cm.wmodeSet = true
					}
				}
			}
		}
		stack = stack[:0]
	}
	cm.finish(parent)
	return cm
}

func (cm *cmap) noteLen(n int) {
	switch n {
	case 1:
		cm.seen1 = true
	case 2:
		cm.seen2 = true
	}
}

func readUntil(l *lexer, end string) []types.Object {
	var items []types.Object
	for {
		kind, obj, op := l.next()
		if kind == tokEOF || (kind == tokOperator && op == end) {
			return items
		}
		if kind == tokOperand {
			items = append(items, obj)
		}
	}
}

func numOf(o types.Object) (float64, bool) {
	switch v := o.(type) {
	case types.Integer:
		return float64(v), true
	case types.Float:
		return float64(v), true
	}
	return 0, false
}

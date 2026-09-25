package pdf2bdf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// Type 1 font programs (FontFile) are converted to a name-keyed CFF so they
// take the same path as embedded CFF fonts (subsetting, OpenType wrapper,
// WOFF2). The charstrings are interpreted, not translated operator by
// operator: subroutines are inlined, flex becomes two curves, seac accents
// are composed into one outline and hints are dropped (as pdf.js does; the
// glyphs are drawn antialiased at screen sizes).

// type1Font is what a Type 1 font program provides.
type type1Font struct {
	name       string
	fontMatrix []float64
	bbox       []float64
	encoding   map[int]string // built-in encoding; nil for StandardEncoding
	info       map[string]string
	italic     float64
	fixedPitch bool
	lenIV      int
	subrs      [][]byte
	names      []string          // glyph names in program order
	charString map[string][]byte // decrypted, lenIV bytes removed
}

// type1ToCFF converts a Type 1 font program into a CFF font program. It
// returns the number of glyphs whose outlines could not be converted (they
// are left empty).
func type1ToCFF(data []byte) ([]byte, *cffFont, int, error) {
	t, err := parseType1(data)
	if err != nil {
		return nil, nil, 0, err
	}
	out, failed, err := t.cff()
	if err != nil {
		return nil, nil, 0, err
	}
	cf, err := parseCFF(out)
	if err != nil {
		return nil, nil, 0, err
	}
	// The CFF keeps StandardEncoding; the program's own encoding is what a
	// PDF font without /Encoding (or with Differences on top of it) uses.
	enc := t.encoding
	if enc == nil {
		enc = map[int]string{}
		for code, name := range standardEncoding {
			if name != "" {
				enc[code] = name
			}
		}
	}
	cf.encoding = map[int]int{}
	for code, name := range enc {
		if gid, ok := cf.nameToGID[name]; ok {
			cf.encoding[code] = gid
		}
	}
	return out, cf, failed, nil
}

// --- parsing ---

// psScanner reads PostScript tokens from a Type 1 program.
type psScanner struct {
	b []byte
	p int
}

func isPSSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\f' || c == 0
}

func isPSDelim(c byte) bool {
	return c == '(' || c == ')' || c == '<' || c == '>' || c == '[' || c == ']' || c == '{' || c == '}' || c == '/' || c == '%'
}

// next returns the next token: names keep their leading '/', strings are
// returned with their parentheses, delimiters are one byte long.
func (s *psScanner) next() []byte {
	for s.p < len(s.b) {
		c := s.b[s.p]
		switch {
		case isPSSpace(c):
			s.p++
		case c == '%':
			for s.p < len(s.b) && s.b[s.p] != '\n' && s.b[s.p] != '\r' {
				s.p++
			}
		default:
			start := s.p
			switch c {
			case '(':
				depth := 0
				for s.p < len(s.b) {
					switch s.b[s.p] {
					case '\\':
						s.p++
					case '(':
						depth++
					case ')':
						depth--
					}
					s.p++
					if depth == 0 {
						break
					}
				}
				return s.b[start:min(s.p, len(s.b))]
			case '[', ']', '{', '}':
				s.p++
				return s.b[start:s.p]
			case '<', '>':
				s.p++
				if s.p < len(s.b) && s.b[s.p] == c {
					s.p++
				}
				return s.b[start:s.p]
			case '/':
				s.p++
			}
			for s.p < len(s.b) && !isPSSpace(s.b[s.p]) && !isPSDelim(s.b[s.p]) {
				s.p++
			}
			if s.p == start { // a stray ')' (or binary data): one byte
				s.p++
			}
			return s.b[start:s.p]
		}
	}
	return nil
}

func (s *psScanner) number() (float64, bool) {
	v, err := strconv.ParseFloat(string(s.next()), 64)
	return v, err == nil
}

// psString decodes a (string) token.
func psString(tok []byte) string {
	if len(tok) < 2 || tok[0] != '(' {
		return ""
	}
	tok = tok[1 : len(tok)-1]
	var out []byte
	for i := 0; i < len(tok); i++ {
		c := tok[i]
		if c != '\\' || i+1 >= len(tok) {
			out = append(out, c)
			continue
		}
		i++
		switch e := tok[i]; e {
		case 'n':
			out = append(out, '\n')
		case 'r':
			out = append(out, '\r')
		case 't':
			out = append(out, '\t')
		case 'b':
			out = append(out, '\b')
		case 'f':
			out = append(out, '\f')
		case '\r', '\n':
		default:
			if e >= '0' && e <= '7' {
				v := 0
				for k := 0; k < 3 && i < len(tok) && tok[i] >= '0' && tok[i] <= '7'; k++ {
					v = v*8 + int(tok[i]-'0')
					i++
				}
				i--
				out = append(out, byte(v))
			} else {
				out = append(out, e)
			}
		}
	}
	// Latin-1 to UTF-8.
	r := make([]rune, len(out))
	for i, c := range out {
		r[i] = rune(c)
	}
	return string(r)
}

// numbers reads a [...] or {...} array of numbers.
func (s *psScanner) numbers() []float64 {
	open := s.next()
	if len(open) != 1 || (open[0] != '[' && open[0] != '{') {
		return nil
	}
	var out []float64
	for {
		tok := s.next()
		if tok == nil || len(tok) == 1 && (tok[0] == ']' || tok[0] == '}') {
			return out
		}
		if v, err := strconv.ParseFloat(string(tok), 64); err == nil {
			out = append(out, v)
		}
	}
}

func decrypt(b []byte, r uint16, skip int) []byte {
	if skip < 0 { // lenIV -1: charstrings are not encrypted
		return b
	}
	out := make([]byte, len(b))
	for i, c := range b {
		out[i] = c ^ byte(r>>8)
		r = (uint16(c)+r)*52845 + 22719
	}
	if skip > len(out) {
		return nil
	}
	return out[skip:]
}

func isHex(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}

// parseType1 reads a Type 1 font program (PFA, PFB or the PDF FontFile form).
func parseType1(data []byte) (*type1Font, error) {
	if len(data) > 6 && data[0] == 0x80 {
		// PFB: segments of ASCII (1) and binary (2) data.
		var joined []byte
		for p := 0; p+6 <= len(data) && data[p] == 0x80 && data[p+1] != 3; {
			n := int(binary.LittleEndian.Uint32(data[p+2:]))
			if n < 0 || p+6+n > len(data) {
				return nil, errors.New("type1: truncated PFB segment")
			}
			joined = append(joined, data[p+6:p+6+n]...)
			p += 6 + n
		}
		data = joined
	}
	t := &type1Font{fontMatrix: []float64{0.001, 0, 0, 0.001, 0, 0}, info: map[string]string{}, lenIV: 4, charString: map[string][]byte{}}
	s := &psScanner{b: data}
	eexec := -1
	for eexec < 0 {
		tok := s.next()
		if tok == nil {
			return nil, errors.New("type1: no eexec section")
		}
		switch string(tok) {
		case "eexec":
			eexec = s.p
		case "/FontName":
			t.name = string(bytes.TrimPrefix(s.next(), []byte("/")))
		case "/FontMatrix":
			if m := s.numbers(); len(m) == 6 {
				t.fontMatrix = m
			}
		case "/FontBBox":
			if b := s.numbers(); len(b) == 4 {
				t.bbox = b
			}
		case "/Encoding":
			if tok := s.next(); string(tok) != "StandardEncoding" {
				t.encoding = map[int]string{}
				// dup <code> /<name> put ... def
				for {
					tok := s.next()
					if tok == nil || string(tok) == "def" {
						break
					}
					if string(tok) != "dup" {
						continue
					}
					save := s.p
					code, ok := s.number()
					name := s.next()
					if !ok || len(name) < 2 || name[0] != '/' || string(s.next()) != "put" {
						s.p = save
						continue
					}
					if code >= 0 && code < 256 {
						t.encoding[int(code)] = string(name[1:])
					}
				}
			}
		case "/Notice", "/Copyright", "/FullName", "/FamilyName", "/Weight":
			save := s.p
			if v := s.next(); len(v) > 0 && v[0] == '(' {
				t.info[string(tok[1:])] = psString(v)
			} else {
				s.p = save
			}
		case "/ItalicAngle":
			t.italic, _ = s.number()
		case "/isFixedPitch":
			t.fixedPitch = string(s.next()) == "true"
		}
	}
	// The encrypted part follows eexec and one whitespace run; PFA fonts
	// have it in hex.
	for eexec < len(data) && isPSSpace(data[eexec]) {
		eexec++
	}
	enc := data[eexec:]
	if len(enc) >= 4 && isHex(enc[0]) && isHex(enc[1]) && isHex(enc[2]) && isHex(enc[3]) {
		var raw []byte
		var hi byte
		half := false
		for _, c := range enc {
			var v byte
			switch {
			case c >= '0' && c <= '9':
				v = c - '0'
			case c >= 'a' && c <= 'f':
				v = c - 'a' + 10
			case c >= 'A' && c <= 'F':
				v = c - 'A' + 10
			case isPSSpace(c):
				continue
			default:
				goto done
			}
			if half {
				raw = append(raw, hi<<4|v)
			} else {
				hi = v
			}
			half = !half
		}
	done:
		enc = raw
	}
	priv := decrypt(enc, 55665, 4)

	// Private dict: lenIV, Subrs, CharStrings. The charstrings are read raw
	// (they are binary) and decrypted once lenIV is known.
	s = &psScanner{b: priv}
	var subrsRaw [][]byte
	csRaw := map[string][]byte{}
	readBinary := func() ([]byte, bool) {
		n, ok := s.number()
		if !ok || n < 0 {
			return nil, false
		}
		s.next() // RD or -|
		s.p++    // the single space before the data
		if s.p+int(n) > len(priv) {
			return nil, false
		}
		b := priv[s.p : s.p+int(n)]
		s.p += int(n)
		return b, true
	}
	for {
		tok := s.next()
		if tok == nil {
			break
		}
		switch string(tok) {
		case "/lenIV":
			if v, ok := s.number(); ok {
				t.lenIV = int(v)
			}
		case "/Subrs":
			count, ok := s.number()
			if !ok || count < 0 || count > 65535 {
				return nil, errors.New("type1: bad Subrs")
			}
			subrsRaw = make([][]byte, int(count))
			// Subsetters may leave out unused entries: stop at the first
			// "dup" that is not followed by an index (dup /CharStrings …).
			for k := 0; k < int(count); k++ {
				for tok := s.next(); tok != nil && string(tok) != "dup"; tok = s.next() {
				}
				save := s.p
				idx, ok := s.number()
				if !ok {
					s.p = save
					break
				}
				b, ok := readBinary()
				if !ok {
					return nil, errors.New("type1: truncated Subrs")
				}
				if idx >= 0 && int(idx) < len(subrsRaw) {
					subrsRaw[int(idx)] = b
				}
			}
		case "/CharStrings":
			count, ok := s.number()
			if !ok || count < 0 {
				return nil, errors.New("type1: bad CharStrings")
			}
			for len(t.names) < int(count) {
				tok := s.next()
				if tok == nil || string(tok) == "end" {
					break
				}
				if len(tok) < 2 || tok[0] != '/' {
					continue
				}
				b, ok := readBinary()
				if !ok {
					return nil, errors.New("type1: truncated CharStrings")
				}
				name := string(tok[1:])
				if _, dup := csRaw[name]; !dup {
					csRaw[name] = b
					t.names = append(t.names, name)
				}
			}
		}
	}
	if len(t.names) == 0 {
		return nil, errors.New("type1: no CharStrings")
	}
	t.subrs = make([][]byte, len(subrsRaw))
	for i, b := range subrsRaw {
		t.subrs[i] = decrypt(b, 4330, t.lenIV)
	}
	for name, b := range csRaw {
		t.charString[name] = decrypt(b, 4330, t.lenIV)
	}
	return t, nil
}

// --- charstring conversion ---

var errT1 = errors.New("type1: charstring cannot be converted")

// t2Writer emits a Type 2 charstring from absolute coordinates, merging runs
// of lines and curves into one operator.
type t2Writer struct {
	out       []byte
	cx, cy    float64 // current point of the emitted path
	op        byte
	args      []float64
	width     float64
	widthDone bool
}

func appendT2Number(b []byte, v float64) []byte {
	if v == math.Trunc(v) && v >= -32768 && v <= 32767 {
		i := int(v)
		switch {
		case i >= -107 && i <= 107:
			return append(b, byte(i+139))
		case i >= 108 && i <= 1131:
			i -= 108
			return append(b, byte(i>>8+247), byte(i))
		case i >= -1131 && i <= -108:
			i = -i - 108
			return append(b, byte(i>>8+251), byte(i))
		default:
			return append(b, 28, byte(i>>8), byte(i))
		}
	}
	f := int32(math.Round(v * 65536))
	return append(b, 255, byte(f>>24), byte(f>>16), byte(f>>8), byte(f))
}

func (w *t2Writer) flush() {
	if w.op == 0 {
		return
	}
	for _, a := range w.args {
		w.out = appendT2Number(w.out, a)
	}
	w.out = append(w.out, w.op)
	w.op, w.args = 0, w.args[:0]
}

// delta returns the offset from the emitted point to (x, y) as it will be
// written (16.16 when fractional) and advances the emitted point by it.
func (w *t2Writer) delta(x, y float64) (float64, float64) {
	dx, dy := x-w.cx, y-w.cy
	if dx != math.Trunc(dx) {
		dx = math.Round(dx*65536) / 65536
	}
	if dy != math.Trunc(dy) {
		dy = math.Round(dy*65536) / 65536
	}
	w.cx += dx
	w.cy += dy
	return dx, dy
}

func (w *t2Writer) prefixWidth() {
	if !w.widthDone {
		w.widthDone = true
		if w.width != 0 { // nominalWidthX is 0
			w.out = appendT2Number(w.out, w.width)
		}
	}
}

func (w *t2Writer) moveTo(x, y float64) {
	w.flush()
	dx, dy := w.delta(x, y)
	w.prefixWidth()
	w.out = appendT2Number(appendT2Number(w.out, dx), dy)
	w.out = append(w.out, 21) // rmoveto
}

func (w *t2Writer) segment(op byte, pts ...float64) {
	n := len(pts)
	if w.op != op || len(w.args)+n > 48 {
		w.flush()
		w.op = op
	}
	for i := 0; i < n; i += 2 {
		dx, dy := w.delta(pts[i], pts[i+1])
		w.args = append(w.args, dx, dy)
	}
}

func (w *t2Writer) end() []byte {
	w.flush()
	w.prefixWidth()
	return append(w.out, 14) // endchar
}

// t1Converter interprets one Type 1 glyph (and the glyphs a seac uses).
type t1Converter struct {
	font      *type1Font
	w         *t2Writer
	stack     []float64
	ps        []float64 // results of othersubrs, read back by pop
	x, y      float64   // current point, absolute
	ox, oy    float64   // origin of the glyph being read (seac accent)
	sbx       float64   // left sidebearing of the glyph (hsbw), where a seac accent is measured from
	moved     bool      // a subpath has been started
	component bool      // reading a seac base or accent: keep the outer width
	flex      bool
	flexPts   [][2]float64
	ops       int
}

func (c *t1Converter) pop(n int) ([]float64, error) {
	if len(c.stack) < n {
		return nil, errT1
	}
	v := c.stack[len(c.stack)-n:]
	c.stack = c.stack[:len(c.stack)-n]
	return v, nil
}

func (c *t1Converter) moveTo(x, y float64) {
	c.x, c.y = x, y
	if c.flex {
		c.flexPts = append(c.flexPts, [2]float64{x, y})
		return
	}
	c.w.moveTo(x, y)
	c.moved = true
}

func (c *t1Converter) lineTo(x, y float64) {
	if !c.moved {
		c.w.moveTo(c.x, c.y)
		c.moved = true
	}
	c.x, c.y = x, y
	c.w.segment(5, x, y) // rlineto
}

func (c *t1Converter) curveTo(p ...float64) {
	if !c.moved {
		c.w.moveTo(c.x, c.y)
		c.moved = true
	}
	c.x, c.y = p[4], p[5]
	c.w.segment(8, p...) // rrcurveto
}

// run interprets a charstring. It returns true at endchar (or seac).
func (c *t1Converter) run(cs []byte, depth int) (bool, error) {
	if depth > 10 {
		return false, errT1
	}
	for i := 0; i < len(cs); {
		if c.ops++; c.ops > 100000 || len(c.stack) > 64 {
			return false, errT1
		}
		v := cs[i]
		switch {
		case v >= 32 && v <= 246:
			c.stack = append(c.stack, float64(int(v)-139))
			i++
			continue
		case v >= 247 && v <= 254:
			if i+1 >= len(cs) {
				return false, errT1
			}
			if v <= 250 {
				c.stack = append(c.stack, float64((int(v)-247)*256+int(cs[i+1])+108))
			} else {
				c.stack = append(c.stack, float64(-(int(v)-251)*256-int(cs[i+1])-108))
			}
			i += 2
			continue
		case v == 255:
			if i+5 > len(cs) {
				return false, errT1
			}
			c.stack = append(c.stack, float64(int32(binary.BigEndian.Uint32(cs[i+1:]))))
			i += 5
			continue
		}
		op := int(v)
		i++
		if v == 12 {
			if i >= len(cs) {
				return false, errT1
			}
			op = 1200 + int(cs[i])
			i++
		}
		switch op {
		case 1, 3, 9, 1200, 1201, 1202: // hints, closepath, dotsection
		case 13, 1207: // hsbw, sbw
			n := 2
			if op == 1207 {
				n = 4
			}
			a, err := c.pop(n)
			if err != nil {
				return false, err
			}
			sby, wx := 0.0, a[1]
			if n == 4 {
				sby, wx = a[1], a[2]
			}
			if !c.component {
				c.w.width, c.sbx = wx, a[0]
			}
			c.x, c.y = c.ox+a[0], c.oy+sby
		case 21:
			a, err := c.pop(2)
			if err != nil {
				return false, err
			}
			c.moveTo(c.x+a[0], c.y+a[1])
		case 22:
			a, err := c.pop(1)
			if err != nil {
				return false, err
			}
			c.moveTo(c.x+a[0], c.y)
		case 4:
			a, err := c.pop(1)
			if err != nil {
				return false, err
			}
			c.moveTo(c.x, c.y+a[0])
		case 5:
			a, err := c.pop(2)
			if err != nil {
				return false, err
			}
			c.lineTo(c.x+a[0], c.y+a[1])
		case 6:
			a, err := c.pop(1)
			if err != nil {
				return false, err
			}
			c.lineTo(c.x+a[0], c.y)
		case 7:
			a, err := c.pop(1)
			if err != nil {
				return false, err
			}
			c.lineTo(c.x, c.y+a[0])
		case 8, 30, 31: // rrcurveto, vhcurveto, hvcurveto
			var d [6]float64
			switch op {
			case 8:
				a, err := c.pop(6)
				if err != nil {
					return false, err
				}
				copy(d[:], a)
			case 30:
				a, err := c.pop(4)
				if err != nil {
					return false, err
				}
				d = [6]float64{0, a[0], a[1], a[2], a[3], 0}
			case 31:
				a, err := c.pop(4)
				if err != nil {
					return false, err
				}
				d = [6]float64{a[0], 0, a[1], a[2], 0, a[3]}
			}
			x1, y1 := c.x+d[0], c.y+d[1]
			x2, y2 := x1+d[2], y1+d[3]
			c.curveTo(x1, y1, x2, y2, x2+d[4], y2+d[5])
		case 10: // callsubr
			a, err := c.pop(1)
			if err != nil {
				return false, err
			}
			k := int(a[0])
			if k < 0 || k >= len(c.font.subrs) || c.font.subrs[k] == nil {
				return false, errT1
			}
			if done, err := c.run(c.font.subrs[k], depth+1); err != nil || done {
				return done, err
			}
			continue
		case 11: // return
			return false, nil
		case 14: // endchar
			return true, nil
		case 1206: // seac: compose the base and the accent
			a, err := c.pop(5)
			if err != nil || c.component {
				return false, errT1
			}
			asb, adx, ady, bchar, achar := a[0], a[1], a[2], int(a[3]), int(a[4])
			if bchar < 0 || bchar > 255 || achar < 0 || achar > 255 {
				return false, errT1
			}
			base, accent := c.font.charString[standardEncoding[bchar]], c.font.charString[standardEncoding[achar]]
			if base == nil || accent == nil {
				return false, errT1
			}
			c.component = true
			for k, part := range [][]byte{base, accent} {
				c.stack, c.moved = c.stack[:0], false
				c.ox, c.oy = 0, 0
				if k == 1 {
					c.ox, c.oy = c.sbx+adx-asb, ady
				}
				if _, err := c.run(part, depth+1); err != nil {
					return false, err
				}
			}
			return true, nil
		case 1212: // div
			a, err := c.pop(2)
			if err != nil || a[1] == 0 {
				return false, errT1
			}
			c.stack = append(c.stack, a[0]/a[1])
			continue
		case 1216: // callothersubr
			a, err := c.pop(2)
			if err != nil {
				return false, err
			}
			n, which := int(a[0]), int(a[1])
			if n < 0 {
				return false, errT1
			}
			args, err := c.pop(n)
			if err != nil {
				return false, err
			}
			switch {
			case which == 0 && n == 3: // end of flex: two curves through the collected points
				c.flex = false
				if pts := c.flexPts; len(pts) == 7 {
					c.curveTo(pts[1][0], pts[1][1], pts[2][0], pts[2][1], pts[3][0], pts[3][1])
					c.curveTo(pts[4][0], pts[4][1], pts[5][0], pts[5][1], pts[6][0], pts[6][1])
				} else if len(pts) > 0 {
					c.lineTo(pts[len(pts)-1][0], pts[len(pts)-1][1])
				}
				c.ps = []float64{args[2], args[1]} // pop pop → x y
			case which == 1 && n == 0: // start of flex
				c.flex, c.flexPts = true, nil
			case which == 2 && n == 0: // flex point (recorded by the moveto)
			case which == 3 && n == 1: // hint replacement: returns the subr number
				c.ps = []float64{args[0]}
			case which >= 14 && which <= 18: // multiple master blends
				return false, errT1
			default:
				c.ps = c.ps[:0]
				for k := len(args) - 1; k >= 0; k-- {
					c.ps = append(c.ps, args[k])
				}
			}
			continue
		case 1217: // pop
			if len(c.ps) == 0 {
				return false, errT1
			}
			c.stack = append(c.stack, c.ps[len(c.ps)-1])
			c.ps = c.ps[:len(c.ps)-1]
			continue
		case 1233: // setcurrentpoint
			a, err := c.pop(2)
			if err != nil {
				return false, err
			}
			c.x, c.y = c.ox+a[0], c.oy+a[1]
		default:
			return false, errT1
		}
		c.stack = c.stack[:0]
	}
	return false, nil
}

// convert returns the Type 2 charstring of a glyph.
func (t *type1Font) convert(name string) ([]byte, error) {
	cs, ok := t.charString[name]
	if !ok {
		return nil, errT1
	}
	c := &t1Converter{font: t, w: &t2Writer{}}
	if _, err := c.run(cs, 0); err != nil {
		return nil, err
	}
	return c.w.end(), nil
}

// --- CFF output ---

// appendDictNumber writes a DICT operand (integer or real).
func appendDictNumber(b []byte, v float64) []byte {
	if v == math.Trunc(v) && math.Abs(v) < 1<<31 {
		i := int32(v)
		switch {
		case i >= -107 && i <= 107:
			return append(b, byte(i+139))
		case i >= 108 && i <= 1131:
			i -= 108
			return append(b, byte(i>>8+247), byte(i))
		case i >= -1131 && i <= -108:
			i = -i - 108
			return append(b, byte(i>>8+251), byte(i))
		case i >= -32768 && i <= 32767:
			return append(b, 28, byte(i>>8), byte(i))
		default:
			return binary.BigEndian.AppendUint32(append(b, 29), uint32(i))
		}
	}
	s := strconv.FormatFloat(v, 'g', -1, 64)
	var nibbles []byte
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c >= '0' && c <= '9':
			nibbles = append(nibbles, c-'0')
		case c == '.':
			nibbles = append(nibbles, 0xa)
		case c == '-':
			nibbles = append(nibbles, 0xe)
		case c == 'e':
			if i+1 < len(s) && s[i+1] == '-' {
				nibbles = append(nibbles, 0xc)
				i++
			} else {
				nibbles = append(nibbles, 0xb)
				if i+1 < len(s) && s[i+1] == '+' {
					i++
				}
			}
		}
	}
	nibbles = append(nibbles, 0xf)
	if len(nibbles)%2 == 1 {
		nibbles = append(nibbles, 0xf)
	}
	b = append(b, 30)
	for i := 0; i < len(nibbles); i += 2 {
		b = append(b, nibbles[i]<<4|nibbles[i+1])
	}
	return b
}

// cff builds a name-keyed CFF font from the converted glyphs. Glyphs whose
// charstrings cannot be converted are left empty and counted.
func (t *type1Font) cff() ([]byte, int, error) {
	names := []string{".notdef"}
	for _, n := range t.names {
		if n != ".notdef" {
			names = append(names, n)
		}
	}
	failed := 0
	charStrings := make([][]byte, len(names))
	for i, n := range names {
		cs, err := t.convert(n)
		if err != nil {
			if _, ok := t.charString[n]; ok {
				failed++
			}
			cs = []byte{14}
		}
		charStrings[i] = cs
	}

	std := make(map[string]int, len(cffStandardStrings))
	for i, s := range cffStandardStrings {
		std[s] = i
	}
	var strs [][]byte
	sid := func(s string) int {
		if i, ok := std[s]; ok {
			return i
		}
		strs = append(strs, []byte(s))
		std[s] = len(cffStandardStrings) + len(strs) - 1
		return std[s]
	}
	charset := []byte{0}
	for _, n := range names[1:] {
		charset = binary.BigEndian.AppendUint16(charset, uint16(sid(n)))
	}
	var top []cffDictEntry
	add := func(op int, vals ...float64) {
		var raw []byte
		for _, v := range vals {
			raw = appendDictNumber(raw, v)
		}
		top = append(top, cffDictEntry{op: op, args: vals, raw: raw})
	}
	for _, k := range []struct {
		key string
		op  int
	}{{"Notice", 1}, {"Copyright", 1200}, {"FullName", 2}, {"FamilyName", 3}, {"Weight", 4}} {
		if v := t.info[k.key]; v != "" {
			add(k.op, float64(sid(v)))
		}
	}
	if t.fixedPitch {
		add(1201, 1)
	}
	if t.italic != 0 {
		add(1202, t.italic)
	}
	if m := t.fontMatrix; m[0] != 0.001 || m[1] != 0 || m[2] != 0 || m[3] != 0.001 || m[4] != 0 || m[5] != 0 {
		add(1207, m...)
	}
	if len(t.bbox) == 4 {
		add(5, t.bbox...)
	}
	add(15, 0)
	add(17, 0)
	add(18, 0, 0)
	private := []byte{139, 20} // defaultWidthX 0 (widths are written against nominalWidthX 0)

	name := t.name
	if name == "" {
		name = "Type1"
	}
	name = sanitizePS(name)
	offsets := map[int][]int{}
	encodeTop := func() []byte {
		return cffEncodeDict(top, func(e cffDictEntry) ([]int, bool) {
			switch e.op {
			case 15, 17, 18:
				if v, ok := offsets[e.op]; ok {
					return v, true
				}
				return make([]int, len(e.args)), true
			}
			return nil, false
		})
	}
	head := []byte{1, 0, 4, 4}
	nameIdx := cffIndex([][]byte{[]byte(name)})
	strIdx := cffIndex(strs)
	gsubrIdx := cffIndex(nil)
	csIdx := cffIndex(charStrings)
	p := len(head) + len(nameIdx) + len(cffIndex([][]byte{encodeTop()})) + len(strIdx) + len(gsubrIdx)
	offsets[15] = []int{p}
	p += len(charset)
	offsets[17] = []int{p}
	p += len(csIdx)
	offsets[18] = []int{len(private), p}
	p += len(private)
	var out []byte
	for _, part := range [][]byte{head, nameIdx, cffIndex([][]byte{encodeTop()}), strIdx, gsubrIdx, charset, csIdx, private} {
		out = append(out, part...)
	}
	if len(out) != p {
		return nil, 0, fmt.Errorf("type1: CFF layout mismatch")
	}
	return out, failed, nil
}

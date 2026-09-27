package cgm

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/cad"
)

// The clear text encoding (ISO/IEC 8632-4): an element is its name and
// its parameters up to a semicolon or a slash. Parameters are numbers
// (with an optional base, 16#FF), strings in single or double quotes
// (a doubled quote stands for itself), and words for enumerated values.
// Commas, parentheses and white space separate them; comments are
// enclosed in percent signs.

// token is a parameter of the clear text encoding.
type token struct {
	kind byte // 'n' number, 's' string, 'w' word
	s    string
	num  float64
}

// params reads the parameters of an element, in either encoding.
type params interface {
	more() bool
	int() int
	index() int
	name() int
	// enum reads an enumerated value; names lists the clear text words of
	// its values in order, separated by spaces (synonyms by |).
	enum(names string) int
	real() float64
	vdc() float64
	point() cad.Point
	colourIndex() int
	component() float64
	// components is the number of components of a direct colour, and
	// colourMax their largest value without a COLOUR VALUE EXTENT.
	components() int
	colourMax() float64
	str() []byte
	scaleFactor() float64
	sdr() []sdrMember
}

// sdrMember is a member of a structured data record: its data type and
// values.
type sdrMember struct {
	typ    int
	values []any
}

// element is an element of a metafile. Its parameters are read when it is
// interpreted, with the precisions then in effect.
type element struct {
	code
	data []byte    // binary
	toks []token   // clear text
	sub  []element // clear text: the elements of BEGMFDEFAULTS
	incr bool      // clear text: an INCR element
	text bool      // clear text
	word string    // clear text: the element's name
}

func (e *element) params(pr *precisions) params {
	if e.text {
		return &textParams{toks: e.toks, pr: pr, incr: e.incr}
	}
	return &binParams{b: e.data, pr: pr}
}

// textReader splits a clear text metafile into elements.
type textReader struct {
	s   string
	pos int
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == '\f' || c == '\v' || c == ',' || c == '(' || c == ')'
}

// skip skips separators and comments.
func (r *textReader) skip() {
	for r.pos < len(r.s) {
		c := r.s[r.pos]
		switch {
		case isSpace(c):
			r.pos++
		case c == '%':
			end := strings.IndexByte(r.s[r.pos+1:], '%')
			if end < 0 {
				r.pos = len(r.s)
			} else {
				r.pos += end + 2
			}
		default:
			return
		}
	}
}

// token reads a token; terminators come back as words ";".
func (r *textReader) token() (token, bool) {
	r.skip()
	if r.pos >= len(r.s) {
		return token{}, false
	}
	c := r.s[r.pos]
	switch {
	case c == ';' || c == '/':
		r.pos++
		return token{kind: 'w', s: ";"}, true
	case c == '\'' || c == '"':
		var b strings.Builder
		r.pos++
		for r.pos < len(r.s) {
			d := r.s[r.pos]
			r.pos++
			if d == c {
				if r.pos < len(r.s) && r.s[r.pos] == c {
					b.WriteByte(c)
					r.pos++
					continue
				}
				break
			}
			b.WriteByte(d)
		}
		return token{kind: 's', s: b.String()}, true
	case c == '+' || c == '-' || c == '.' || c >= '0' && c <= '9':
		start := r.pos
		r.pos++
		for r.pos < len(r.s) {
			d := r.s[r.pos]
			if d >= '0' && d <= '9' || d == '.' || d == 'e' || d == 'E' || d == '#' || d >= 'a' && d <= 'z' || d >= 'A' && d <= 'Z' ||
				(d == '+' || d == '-') && (r.s[r.pos-1] == 'e' || r.s[r.pos-1] == 'E') {
				r.pos++
				continue
			}
			break
		}
		return token{kind: 'n', s: r.s[start:r.pos], num: parseNumber(r.s[start:r.pos])}, true
	default:
		start := r.pos
		for r.pos < len(r.s) {
			d := r.s[r.pos]
			if isSpace(d) || d == ';' || d == '/' || d == '\'' || d == '"' || d == '%' {
				break
			}
			r.pos++
		}
		if r.pos == start {
			r.pos++
		}
		return token{kind: 'w', s: r.s[start:r.pos]}, true
	}
}

// parseNumber parses a number: decimal, with an exponent, or based
// (base#digits, with an optional sign before the base).
func parseNumber(s string) float64 {
	if i := strings.IndexByte(s, '#'); i > 0 {
		neg := false
		b := s[:i]
		if b[0] == '-' || b[0] == '+' {
			neg = b[0] == '-'
			b = b[1:]
		}
		base, err := strconv.Atoi(b)
		if err != nil || base < 2 || base > 36 {
			return 0
		}
		v, err := strconv.ParseInt(s[i+1:], base, 64)
		if err != nil {
			return 0
		}
		if neg {
			v = -v
		}
		return float64(v)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return finite(v)
}

// keyword normalizes an element or enumeration name: upper case, without
// the underscores and dollar signs the encoding ignores.
func keyword(s string) string {
	s = strings.ToUpper(s)
	if strings.ContainsAny(s, "_$") {
		s = strings.NewReplacer("_", "", "$", "").Replace(s)
	}
	return s
}

// next returns the next element; false at the end of the text.
func (r *textReader) next() (element, bool) { return r.element(true) }

// element reads an element; nested BEGMFDEFAULTS (when defaults is false)
// are read as plain elements.
func (r *textReader) element(defaults bool) (element, bool) {
	for {
		t, ok := r.token()
		if !ok {
			return element{}, false
		}
		if t.kind != 'w' || t.s == ";" {
			continue // stray parameters
		}
		name := keyword(t.s)
		e := element{text: true, word: name}
		if c, ok := clearText[name]; ok {
			e.code = c
		} else if c, ok := incremental[name]; ok {
			e.code, e.incr = c, true
		} else {
			e.code = eUnknownKeyword
		}
		for {
			t, ok := r.token()
			if !ok || t.kind == 'w' && t.s == ";" {
				break
			}
			e.toks = append(e.toks, t)
		}
		if e.code == eDefaultsReplacement && defaults {
			for {
				sub, ok := r.element(false)
				if !ok || sub.code == eEndMFDefaults {
					break
				}
				e.sub = append(e.sub, sub)
			}
		}
		return e, true
	}
}

// textParams reads the parameters of a clear text element.
type textParams struct {
	toks []token
	pos  int
	pr   *precisions
	incr bool
	// last is the last point read, for incremental point lists
	last    cad.Point
	npoints int
}

func (p *textParams) more() bool { return p.pos < len(p.toks) }

func (p *textParams) next() token {
	if p.pos >= len(p.toks) {
		panic(errShort)
	}
	t := p.toks[p.pos]
	p.pos++
	return t
}

func (p *textParams) number() float64 {
	t := p.next()
	if t.kind != 'n' {
		panic(errShort)
	}
	return t.num
}

func (p *textParams) int() int {
	v := p.number()
	return clampInt(int64(max(min(v, 1<<40), -1<<40)))
}

func (p *textParams) index() int         { return p.int() }
func (p *textParams) name() int          { return p.int() }
func (p *textParams) real() float64      { return p.number() }
func (p *textParams) vdc() float64       { return p.number() }
func (p *textParams) colourIndex() int   { return p.int() }
func (p *textParams) component() float64 { return p.number() }
func (p *textParams) components() int    { return p.pr.components }
func (p *textParams) colourMax() float64 {
	return float64(max(p.pr.colourBits, 1))
}
func (p *textParams) scaleFactor() float64 { return p.number() }

func (p *textParams) point() cad.Point {
	q := cad.Point{X: p.number()}
	q.Y = p.number()
	if p.incr && p.npoints > 0 {
		q = q.Add(p.last)
	}
	p.last = q
	p.npoints++
	return q
}

func (p *textParams) enum(names string) int {
	t := p.next()
	if t.kind == 'n' {
		return int(t.num)
	}
	w := keyword(t.s)
	for i, n := range strings.Fields(names) {
		for _, alt := range strings.Split(n, "|") {
			if alt == w {
				return i
			}
		}
	}
	panic(errShort)
}

func (p *textParams) str() []byte {
	t := p.next()
	if t.kind == 'n' {
		return []byte(t.s)
	}
	return []byte(t.s)
}

// sdr reads a structured data record: the rest of the parameters, or the
// tokens of a string holding them, as data type, count and values.
func (p *textParams) sdr() []sdrMember {
	toks := p.toks[p.pos:]
	p.pos = len(p.toks)
	if len(toks) == 1 && toks[0].kind == 's' {
		r := &textReader{s: toks[0].s}
		toks = nil
		for {
			t, ok := r.token()
			if !ok {
				break
			}
			toks = append(toks, t)
		}
	}
	var out []sdrMember
	for i := 0; i+1 < len(toks); {
		if toks[i].kind != 'n' || toks[i+1].kind != 'n' {
			break
		}
		m := sdrMember{typ: int(toks[i].num)}
		n := int(toks[i+1].num)
		i += 2
		for j := 0; j < n && i < len(toks); j++ {
			t := toks[i]
			i++
			switch t.kind {
			case 'n':
				m.values = append(m.values, t.num)
			case 's':
				m.values = append(m.values, []byte(t.s))
			default:
				m.values = append(m.values, keyword(t.s))
			}
		}
		out = append(out, m)
	}
	return out
}

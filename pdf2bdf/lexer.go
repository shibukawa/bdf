package pdf2bdf

import (
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Content stream lexer: produces operands (as pdfcpu types) and operators.

type tokenKind int

const (
	tokEOF tokenKind = iota
	tokOperand
	tokOperator
)

type lexer struct {
	b   []byte
	pos int
}

func isWhite(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t' || c == '\f' || c == 0
}

func isDelim(c byte) bool {
	switch c {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	}
	return false
}

func (l *lexer) skipSpace() {
	for l.pos < len(l.b) {
		c := l.b[l.pos]
		if isWhite(c) {
			l.pos++
		} else if c == '%' {
			for l.pos < len(l.b) && l.b[l.pos] != '\n' && l.b[l.pos] != '\r' {
				l.pos++
			}
		} else {
			return
		}
	}
}

// next returns the next token. Operators are returned as strings.
func (l *lexer) next() (tokenKind, types.Object, string) {
	l.skipSpace()
	if l.pos >= len(l.b) {
		return tokEOF, nil, ""
	}
	c := l.b[l.pos]
	switch {
	case c == '/':
		l.pos++
		return tokOperand, types.Name(l.readName()), ""
	case c == '(':
		l.pos++
		return tokOperand, types.StringLiteral(l.readLiteral()), ""
	case c == '<':
		if l.pos+1 < len(l.b) && l.b[l.pos+1] == '<' {
			l.pos += 2
			return tokOperand, l.readDict(), ""
		}
		l.pos++
		return tokOperand, types.HexLiteral(l.readHex()), ""
	case c == '[':
		l.pos++
		return tokOperand, l.readArray(), ""
	case c == ']' || c == '>' || c == ')' || c == '{' || c == '}':
		l.pos++ // stray delimiter: skip
		return l.next()
	case (c >= '0' && c <= '9') || c == '+' || c == '-' || c == '.':
		return tokOperand, l.readNumber(), ""
	}
	start := l.pos
	for l.pos < len(l.b) && !isWhite(l.b[l.pos]) && !isDelim(l.b[l.pos]) {
		l.pos++
	}
	if l.pos == start {
		l.pos++
		return l.next()
	}
	word := string(l.b[start:l.pos])
	switch word {
	case "true":
		return tokOperand, types.Boolean(true), ""
	case "false":
		return tokOperand, types.Boolean(false), ""
	case "null":
		return tokOperand, nil, ""
	}
	return tokOperator, nil, word
}

func (l *lexer) readName() string {
	start := l.pos
	for l.pos < len(l.b) && !isWhite(l.b[l.pos]) && !isDelim(l.b[l.pos]) {
		l.pos++
	}
	raw := l.b[start:l.pos]
	if bytesIndexByte(raw, '#') < 0 {
		return string(raw)
	}
	out := make([]byte, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		if raw[i] == '#' && i+2 < len(raw) {
			if v, err := strconv.ParseUint(string(raw[i+1:i+3]), 16, 8); err == nil {
				out = append(out, byte(v))
				i += 2
				continue
			}
		}
		out = append(out, raw[i])
	}
	return string(out)
}

func bytesIndexByte(b []byte, c byte) int {
	for i, x := range b {
		if x == c {
			return i
		}
	}
	return -1
}

func (l *lexer) readNumber() types.Object {
	start := l.pos
	l.pos++
	for l.pos < len(l.b) {
		c := l.b[l.pos]
		if (c >= '0' && c <= '9') || c == '.' || c == '-' || c == '+' || c == 'e' || c == 'E' {
			l.pos++
		} else {
			break
		}
	}
	s := string(l.b[start:l.pos])
	if i, err := strconv.Atoi(s); err == nil {
		return types.Integer(i)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// Things like "--5" or "3.4.5" appear in the wild; salvage a prefix.
		for n := len(s) - 1; n > 0; n-- {
			if f, err = strconv.ParseFloat(s[:n], 64); err == nil {
				break
			}
		}
	}
	return types.Float(f)
}

// readLiteral reads a (...) string and returns its unescaped bytes as a
// StringLiteral whose Value is the raw (unescaped) content.
func (l *lexer) readLiteral() string {
	var out []byte
	depth := 1
	for l.pos < len(l.b) {
		c := l.b[l.pos]
		l.pos++
		switch c {
		case '\\':
			if l.pos >= len(l.b) {
				break
			}
			e := l.b[l.pos]
			l.pos++
			switch e {
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
			case '\r':
				if l.pos < len(l.b) && l.b[l.pos] == '\n' {
					l.pos++
				}
			case '\n':
			default:
				if e >= '0' && e <= '7' {
					v := int(e - '0')
					for k := 0; k < 2 && l.pos < len(l.b) && l.b[l.pos] >= '0' && l.b[l.pos] <= '7'; k++ {
						v = v*8 + int(l.b[l.pos]-'0')
						l.pos++
					}
					out = append(out, byte(v))
				} else {
					out = append(out, e)
				}
			}
		case '(':
			depth++
			out = append(out, c)
		case ')':
			depth--
			if depth == 0 {
				return string(out)
			}
			out = append(out, c)
		default:
			out = append(out, c)
		}
	}
	return string(out)
}

func (l *lexer) readHex() string {
	start := l.pos
	for l.pos < len(l.b) && l.b[l.pos] != '>' {
		l.pos++
	}
	s := l.b[start:l.pos]
	if l.pos < len(l.b) {
		l.pos++
	}
	out := make([]byte, 0, len(s))
	for _, c := range s {
		if !isWhite(c) {
			out = append(out, c)
		}
	}
	return string(out)
}

func (l *lexer) readArray() types.Array {
	var a types.Array
	for {
		l.skipSpace()
		if l.pos >= len(l.b) {
			return a
		}
		if l.b[l.pos] == ']' {
			l.pos++
			return a
		}
		kind, obj, op := l.next()
		switch kind {
		case tokEOF:
			return a
		case tokOperand:
			a = append(a, obj)
		case tokOperator:
			// Operators inside arrays are invalid; keep them as names so TJ still works.
			_ = op
		}
	}
}

func (l *lexer) readDict() types.Dict {
	d := types.NewDict()
	for {
		l.skipSpace()
		if l.pos >= len(l.b) {
			return d
		}
		if l.b[l.pos] == '>' {
			l.pos++
			if l.pos < len(l.b) && l.b[l.pos] == '>' {
				l.pos++
			}
			return d
		}
		kind, obj, _ := l.next()
		if kind == tokEOF {
			return d
		}
		key, ok := obj.(types.Name)
		if !ok {
			continue
		}
		kind, val, _ := l.next()
		if kind != tokOperand {
			return d
		}
		d[key.Value()] = val
	}
}

// literalBytes returns the raw bytes of a string operand produced by this lexer.
func literalBytes(o types.Object) []byte {
	switch v := o.(type) {
	case types.StringLiteral:
		return []byte(string(v)) // already unescaped by readLiteral
	case types.HexLiteral:
		b, err := v.Bytes()
		if err != nil {
			return nil
		}
		return b
	}
	return nil
}

// readInlineImageData reads the binary data after the ID operator up to EI.
// length is the expected data length for unfiltered images (or -1).
func (l *lexer) readInlineImageData(length int) []byte {
	if l.pos < len(l.b) && isWhite(l.b[l.pos]) {
		l.pos++
	}
	start := l.pos
	if length >= 0 && start+length <= len(l.b) {
		end := start + length
		// Verify that EI follows (allowing whitespace); otherwise fall back to a search.
		p := end
		for p < len(l.b) && isWhite(l.b[p]) {
			p++
		}
		if p+1 < len(l.b) && l.b[p] == 'E' && l.b[p+1] == 'I' && (p+2 >= len(l.b) || isWhite(l.b[p+2]) || isDelim(l.b[p+2])) {
			l.pos = p + 2
			return l.b[start:end]
		}
	}
	// Search for whitespace + EI + (whitespace|EOF).
	for p := start; p+1 < len(l.b); p++ {
		if l.b[p] == 'E' && l.b[p+1] == 'I' && (p == 0 || isWhite(l.b[p-1])) && (p+2 >= len(l.b) || isWhite(l.b[p+2]) || isDelim(l.b[p+2])) {
			end := p - 1
			l.pos = p + 2
			return l.b[start:end]
		}
	}
	l.pos = len(l.b)
	return l.b[start:]
}

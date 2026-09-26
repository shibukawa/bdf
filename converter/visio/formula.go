package visio

import (
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Visio writes the values of the cells whose formulas it evaluated, so
// formulas are not evaluated in general. Drawings made by other programs
// sometimes resize instances of masters without the geometry that follows
// from the new size; for them, the formulas of cells an instance inherits
// that only scale its Width and Height ("Width*0.5") are evaluated again
// with the instance's size. Visio's own drawings hold these values in the
// instance, so they are not affected.

// evalSize evaluates a formula made of numbers (with units), Width,
// Height, + - * / and parentheses; ok is false for anything else.
func evalSize(f string, w, h float64) (v float64, ok bool) {
	p := &exprParser{s: strings.TrimSpace(f), w: w, h: h}
	if p.s == "" {
		return 0, false
	}
	defer func() {
		if recover() != nil {
			v, ok = 0, false
		}
	}()
	v = p.expr()
	p.space()
	if p.i != len(p.s) || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, false
	}
	return v, true
}

type exprParser struct {
	s    string
	i    int
	w, h float64
}

type badExpr struct{}

func (p *exprParser) space() {
	for p.i < len(p.s) && p.s[p.i] == ' ' {
		p.i++
	}
}

func (p *exprParser) peek() byte {
	p.space()
	if p.i < len(p.s) {
		return p.s[p.i]
	}
	return 0
}

func (p *exprParser) expr() float64 {
	v := p.term()
	for {
		switch p.peek() {
		case '+':
			p.i++
			v += p.term()
		case '-':
			p.i++
			v -= p.term()
		default:
			return v
		}
	}
}

func (p *exprParser) term() float64 {
	v := p.factor()
	for {
		switch p.peek() {
		case '*':
			p.i++
			v *= p.factor()
		case '/':
			p.i++
			d := p.factor()
			if d == 0 {
				panic(badExpr{})
			}
			v /= d
		default:
			return v
		}
	}
}

func (p *exprParser) factor() float64 {
	switch c := p.peek(); {
	case c == '-':
		p.i++
		return -p.factor()
	case c == '(':
		p.i++
		v := p.expr()
		if p.peek() != ')' {
			panic(badExpr{})
		}
		p.i++
		return v
	case c >= '0' && c <= '9' || c == '.':
		start := p.i
		digits := func() {
			for p.i < len(p.s) && (p.s[p.i] >= '0' && p.s[p.i] <= '9' || p.s[p.i] == '.') {
				p.i++
			}
		}
		digits()
		// an exponent (1E-3), not a unit
		if p.i+1 < len(p.s) && (p.s[p.i] == 'E' || p.s[p.i] == 'e') && strings.ContainsRune("+-0123456789", rune(p.s[p.i+1])) {
			p.i += 2
			digits()
		}
		v, err := strconv.ParseFloat(p.s[start:p.i], 64)
		if err != nil {
			panic(badExpr{})
		}
		return v * p.unit()
	case unicode.IsLetter(rune(c)):
		start := p.i
		for p.i < len(p.s) && (unicode.IsLetter(rune(p.s[p.i])) || p.s[p.i] == '.' || p.s[p.i] >= '0' && p.s[p.i] <= '9') {
			p.i++
		}
		switch strings.ToLower(p.s[start:p.i]) {
		case "width":
			return p.w
		case "height":
			return p.h
		}
	}
	panic(badExpr{})
}

// unit reads the unit after a number and returns its size in internal
// units (inches, radians).
func (p *exprParser) unit() float64 {
	p.space()
	start := p.i
	for p.i < len(p.s) && unicode.IsLetter(rune(p.s[p.i])) {
		p.i++
	}
	switch strings.ToLower(p.s[start:p.i]) {
	case "":
		return 1
	case "in", "dl", "il":
		return 1
	case "pt", "dp":
		return 1.0 / 72
	case "mm":
		return 1 / 25.4
	case "cm":
		return 1 / 2.54
	case "ft":
		return 12
	case "deg":
		return math.Pi / 180
	case "rad":
		return 1
	}
	p.i = start
	panic(badExpr{})
}

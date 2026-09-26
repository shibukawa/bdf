package xlsx

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// A small evaluator of Excel formulas for conditional formats: the rules'
// expressions ("MOD(ROW(),2)=0", "$C2>AVERAGE($C:$C)", "COUNTIF($A:$A,A1)>1")
// work on the cached values of the cells. It knows the operators, cell and
// range references (relative ones move with the cell the rule is evaluated
// for, from the first cell of the rule's range; other sheets by name) and
// the common logical, information, math, text, statistical and date
// functions. Functions whose value depends on when the file is opened
// (TODAY, NOW, RAND) and anything else make the formula unsupported.

type fkind byte

const (
	fBlank fkind = iota
	fNum
	fStr
	fBool
	fErr
	fRange
)

type fval struct {
	kind fkind
	num  float64
	str  string
	rng  *frange
}

type frange struct {
	ws *worksheet
	cellRange
}

func fnum(v float64) fval { return fval{kind: fNum, num: v} }
func fstr(s string) fval  { return fval{kind: fStr, str: s} }
func ferr(s string) fval  { return fval{kind: fErr, str: s} }
func fbool(b bool) fval {
	if b {
		return fval{kind: fBool, num: 1}
	}
	return fval{kind: fBool}
}

// expression tree
type fnode interface{}

type (
	fLit struct{ v fval }
	fRef struct {
		sheet          string
		r0, c0, r1, c1 int
		absR0, absC0   bool
		absR1, absC1   bool
		wholeCols      bool
		wholeRows      bool
		isRange        bool
	}
	fUnary struct {
		op string
		x  fnode
	}
	fBinary struct {
		op   string
		a, b fnode
	}
	fCall struct {
		name string
		args []fnode
	}
)

var errUnsupported = fmt.Errorf("unsupported formula")

type fparser struct {
	s   string
	pos int
}

// parseFormula parses a formula (with or without its leading "=").
func parseFormula(s string) (fnode, error) {
	s = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(s), "="))
	p := &fparser{s: s}
	n, err := p.expr(0)
	if err != nil {
		return nil, err
	}
	p.space()
	if p.pos != len(p.s) {
		return nil, fmt.Errorf("unexpected %q", p.s[p.pos:])
	}
	return n, nil
}

var precedence = map[string]int{"=": 1, "<>": 1, "<": 1, ">": 1, "<=": 1, ">=": 1, "&": 2, "+": 3, "-": 3, "*": 4, "/": 4, "^": 5}

func (p *fparser) space() {
	for p.pos < len(p.s) && (p.s[p.pos] == ' ' || p.s[p.pos] == '\n' || p.s[p.pos] == '\r' || p.s[p.pos] == '\t') {
		p.pos++
	}
}

func (p *fparser) binop() string {
	p.space()
	for _, op := range []string{"<>", "<=", ">=", "=", "<", ">", "&", "+", "-", "*", "/", "^"} {
		if strings.HasPrefix(p.s[p.pos:], op) {
			return op
		}
	}
	return ""
}

func (p *fparser) expr(minPrec int) (fnode, error) {
	left, err := p.unary()
	if err != nil {
		return nil, err
	}
	for {
		op := p.binop()
		prec, ok := precedence[op]
		if !ok || prec < minPrec {
			return left, nil
		}
		p.pos += len(op)
		next := prec + 1
		if op == "^" {
			next = prec
		}
		right, err := p.expr(next)
		if err != nil {
			return nil, err
		}
		left = fBinary{op, left, right}
	}
}

func (p *fparser) unary() (fnode, error) {
	p.space()
	if p.pos < len(p.s) && (p.s[p.pos] == '-' || p.s[p.pos] == '+') {
		op := p.s[p.pos : p.pos+1]
		p.pos++
		x, err := p.unary()
		if err != nil {
			return nil, err
		}
		return fUnary{op, x}, nil
	}
	x, err := p.primary()
	if err != nil {
		return nil, err
	}
	p.space()
	for p.pos < len(p.s) && p.s[p.pos] == '%' {
		p.pos++
		x = fBinary{"/", x, fLit{fnum(100)}}
	}
	return x, nil
}

func (p *fparser) primary() (fnode, error) {
	p.space()
	if p.pos >= len(p.s) {
		return nil, fmt.Errorf("unexpected end")
	}
	c := p.s[p.pos]
	switch {
	case c == '(':
		p.pos++
		x, err := p.expr(0)
		if err != nil {
			return nil, err
		}
		p.space()
		if p.pos >= len(p.s) || p.s[p.pos] != ')' {
			return nil, fmt.Errorf("missing )")
		}
		p.pos++
		return x, nil
	case c == '"':
		var b strings.Builder
		p.pos++
		for p.pos < len(p.s) {
			if p.s[p.pos] == '"' {
				if p.pos+1 < len(p.s) && p.s[p.pos+1] == '"' {
					b.WriteByte('"')
					p.pos += 2
					continue
				}
				p.pos++
				return fLit{fstr(b.String())}, nil
			}
			b.WriteByte(p.s[p.pos])
			p.pos++
		}
		return nil, fmt.Errorf("unterminated string")
	case c >= '0' && c <= '9' || c == '.':
		start := p.pos
		for p.pos < len(p.s) && (p.s[p.pos] >= '0' && p.s[p.pos] <= '9' || p.s[p.pos] == '.') {
			p.pos++
		}
		if p.pos < len(p.s) && (p.s[p.pos] == 'E' || p.s[p.pos] == 'e') {
			p.pos++
			if p.pos < len(p.s) && (p.s[p.pos] == '+' || p.s[p.pos] == '-') {
				p.pos++
			}
			for p.pos < len(p.s) && p.s[p.pos] >= '0' && p.s[p.pos] <= '9' {
				p.pos++
			}
		}
		// "1:3" (whole rows) is a reference
		if p.pos < len(p.s) && p.s[p.pos] == ':' {
			p.pos = start
			return p.reference("")
		}
		v, err := strconv.ParseFloat(p.s[start:p.pos], 64)
		if err != nil {
			return nil, err
		}
		return fLit{fnum(v)}, nil
	case c == '#':
		start := p.pos
		for p.pos < len(p.s) && !strings.ContainsRune(",)", rune(p.s[p.pos])) && p.s[p.pos] != ' ' {
			p.pos++
			if p.s[p.pos-1] == '!' || p.s[p.pos-1] == '?' {
				break
			}
		}
		return fLit{ferr(p.s[start:p.pos])}, nil
	case c == '\'':
		// 'Sheet name'!A1
		end := p.pos + 1
		for end < len(p.s) {
			if p.s[end] == '\'' {
				if end+1 < len(p.s) && p.s[end+1] == '\'' {
					end += 2
					continue
				}
				break
			}
			end++
		}
		if end+1 >= len(p.s) || p.s[end+1] != '!' {
			return nil, fmt.Errorf("bad sheet reference")
		}
		sheet := strings.ReplaceAll(p.s[p.pos+1:end], "''", "'")
		p.pos = end + 2
		return p.reference(sheet)
	case c == '$' || unicode.IsLetter(rune(c)) || c == '_':
		start := p.pos
		for p.pos < len(p.s) && (p.s[p.pos] == '$' || p.s[p.pos] == '.' || p.s[p.pos] == '_' || unicode.IsLetter(rune(p.s[p.pos])) || unicode.IsDigit(rune(p.s[p.pos]))) {
			p.pos++
		}
		word := p.s[start:p.pos]
		p.space()
		if p.pos < len(p.s) && p.s[p.pos] == '(' {
			return p.call(strings.ToUpper(strings.TrimPrefix(word, "_xlfn.")))
		}
		if p.pos < len(p.s) && p.s[p.pos] == '!' {
			p.pos++
			return p.reference(word)
		}
		switch strings.ToUpper(word) {
		case "TRUE":
			return fLit{fbool(true)}, nil
		case "FALSE":
			return fLit{fbool(false)}, nil
		}
		p.pos = start
		return p.reference("")
	}
	return nil, fmt.Errorf("unexpected %q", p.s[p.pos:])
}

func (p *fparser) call(name string) (fnode, error) {
	p.pos++ // (
	var args []fnode
	p.space()
	if p.pos < len(p.s) && p.s[p.pos] == ')' {
		p.pos++
		return fCall{name, nil}, nil
	}
	for {
		p.space()
		if p.pos < len(p.s) && (p.s[p.pos] == ',' || p.s[p.pos] == ')') {
			args = append(args, fLit{}) // an omitted argument
		} else {
			a, err := p.expr(0)
			if err != nil {
				return nil, err
			}
			args = append(args, a)
		}
		p.space()
		if p.pos >= len(p.s) {
			return nil, fmt.Errorf("missing )")
		}
		switch p.s[p.pos] {
		case ',':
			p.pos++
		case ')':
			p.pos++
			return fCall{name, args}, nil
		default:
			return nil, fmt.Errorf("unexpected %q", p.s[p.pos:])
		}
	}
}

// reference parses A1, $A$1, A1:B2, A:B, 1:3.
func (p *fparser) reference(sheet string) (fnode, error) {
	part := func() (col, row int, absC, absR, hasC, hasR bool, err error) {
		if p.pos < len(p.s) && p.s[p.pos] == '$' {
			absC = true
			p.pos++
		}
		col = 0
		for p.pos < len(p.s) && unicode.IsLetter(rune(p.s[p.pos])) {
			ch := unicode.ToUpper(rune(p.s[p.pos]))
			col = col*26 + int(ch-'A'+1)
			hasC = true
			p.pos++
		}
		if p.pos < len(p.s) && p.s[p.pos] == '$' {
			if hasC {
				absR = true
			} else {
				absR, absC = true, false
			}
			p.pos++
		}
		start := p.pos
		for p.pos < len(p.s) && p.s[p.pos] >= '0' && p.s[p.pos] <= '9' {
			p.pos++
		}
		if p.pos > start {
			row, _ = strconv.Atoi(p.s[start:p.pos])
			hasR = true
		}
		if !hasC && !hasR || col > maxCols || row > maxRows {
			err = fmt.Errorf("bad reference")
		}
		return col - 1, row - 1, absC, absR, hasC, hasR, err
	}
	c0, r0, ac0, ar0, hc0, hr0, err := part()
	if err != nil {
		return nil, err
	}
	ref := fRef{sheet: sheet, r0: r0, c0: c0, absR0: ar0, absC0: ac0}
	if p.pos < len(p.s) && p.s[p.pos] == ':' {
		p.pos++
		c1, r1, ac1, ar1, hc1, hr1, err := part()
		if err != nil {
			return nil, err
		}
		ref.r1, ref.c1, ref.absR1, ref.absC1, ref.isRange = r1, c1, ar1, ac1, true
		switch {
		case hc0 && hc1 && !hr0 && !hr1:
			ref.wholeCols = true
		case hr0 && hr1 && !hc0 && !hc1:
			ref.wholeRows = true
		case !(hc0 && hr0 && hc1 && hr1):
			return nil, fmt.Errorf("bad range")
		}
		return ref, nil
	}
	if !hc0 || !hr0 {
		return nil, fmt.Errorf("bad reference")
	}
	return ref, nil
}

// fenv is where a formula is evaluated: the sheet and the offset of the
// cell from the rule's first cell.
type fenv struct {
	s      *sheetCtx
	dr, dc int
	row    int // the evaluated cell
	col    int
}

func (e *fenv) sheet(name string) *worksheet {
	if name == "" || name == e.s.ws.name {
		return e.s.ws
	}
	return e.s.c.sheetNamed(name)
}

func (e *fenv) ref(n fRef) (fval, error) {
	ws := e.sheet(n.sheet)
	if ws == nil {
		return fval{}, errUnsupported
	}
	move := func(v int, abs bool, d int) int {
		if abs {
			return v
		}
		return v + d
	}
	var rg cellRange
	switch {
	case n.wholeCols:
		rg = cellRange{0, move(n.c0, n.absC0, e.dc), maxRows - 1, move(n.c1, n.absC1, e.dc)}
	case n.wholeRows:
		rg = cellRange{move(n.r0, n.absR0, e.dr), 0, move(n.r1, n.absR1, e.dr), maxCols - 1}
	default:
		rg = cellRange{move(n.r0, n.absR0, e.dr), move(n.c0, n.absC0, e.dc), move(n.r0, n.absR0, e.dr), move(n.c0, n.absC0, e.dc)}
		if n.isRange {
			rg.r1, rg.c1 = move(n.r1, n.absR1, e.dr), move(n.c1, n.absC1, e.dc)
		}
	}
	if rg.r0 > rg.r1 {
		rg.r0, rg.r1 = rg.r1, rg.r0
	}
	if rg.c0 > rg.c1 {
		rg.c0, rg.c1 = rg.c1, rg.c0
	}
	if rg.r0 < 0 || rg.c0 < 0 || rg.r1 >= maxRows || rg.c1 >= maxCols {
		return ferr("#REF!"), nil
	}
	if !n.isRange {
		return cellFval(ws.cellAt(rg.r0, rg.c0)), nil
	}
	return fval{kind: fRange, rng: &frange{ws: ws, cellRange: rg}}, nil
}

func cellFval(cl *cell) fval {
	if cl == nil {
		return fval{}
	}
	switch cl.kind {
	case cellNum:
		return fnum(cl.num)
	case cellStr:
		return fstr(cl.text.plain)
	case cellBool:
		return fbool(cl.num != 0)
	case cellErr:
		return ferr(cl.text.plain)
	}
	return fval{}
}

// values lists the cells of a range that have values, row by row.
func (r *frange) values(fn func(v fval)) {
	ws := r.ws
	for i := range ws.rows {
		rw := &ws.rows[i]
		if rw.idx < r.r0 {
			continue
		}
		if rw.idx > r.r1 {
			break
		}
		for j := range rw.cells {
			cl := &rw.cells[j]
			if cl.col >= r.c0 && cl.col <= r.c1 && cl.kind != cellBlank {
				fn(cellFval(cl))
			}
		}
	}
}

func (r *frange) size() int { return (r.r1 - r.r0 + 1) * (r.c1 - r.c0 + 1) }

// scalar takes the value a range stands for where one value is expected:
// its cell in the evaluated row or column (implicit intersection).
func (e *fenv) scalar(v fval) fval {
	if v.kind != fRange {
		return v
	}
	r := v.rng
	switch {
	case r.r0 == r.r1 && r.c0 == r.c1:
		return cellFval(r.ws.cellAt(r.r0, r.c0))
	case r.c0 == r.c1 && e.row >= r.r0 && e.row <= r.r1:
		return cellFval(r.ws.cellAt(e.row, r.c0))
	case r.r0 == r.r1 && e.col >= r.c0 && e.col <= r.c1:
		return cellFval(r.ws.cellAt(r.r0, e.col))
	}
	return ferr("#VALUE!")
}

func toNum(v fval) (float64, bool) {
	switch v.kind {
	case fNum, fBool:
		return v.num, true
	case fBlank:
		return 0, true
	case fStr:
		f, err := strconv.ParseFloat(strings.TrimSpace(v.str), 64)
		return f, err == nil
	}
	return 0, false
}

func toStr(v fval) string {
	switch v.kind {
	case fNum:
		return formatGeneral(v.num, 15)
	case fBool:
		if v.num != 0 {
			return "TRUE"
		}
		return "FALSE"
	case fStr, fErr:
		return v.str
	}
	return ""
}

func truthy(v fval) (bool, bool) {
	switch v.kind {
	case fNum, fBool:
		return v.num != 0, true
	case fBlank:
		return false, true
	case fStr:
		switch strings.ToUpper(v.str) {
		case "TRUE":
			return true, true
		case "FALSE":
			return false, true
		}
	}
	return false, false
}

// cmpVals compares values the way Excel does: numbers < text < booleans,
// text case-insensitively; blanks are zero or empty text.
func cmpVals(a, b fval) int {
	rank := func(v fval) int {
		switch v.kind {
		case fStr:
			return 1
		case fBool:
			return 2
		}
		return 0
	}
	if a.kind == fBlank {
		if b.kind == fStr {
			a = fstr("")
		} else if b.kind == fBool {
			a = fbool(false)
		}
	}
	if b.kind == fBlank {
		if a.kind == fStr {
			b = fstr("")
		} else if a.kind == fBool {
			b = fbool(false)
		}
	}
	ra, rb := rank(a), rank(b)
	if ra != rb {
		return ra - rb
	}
	if ra == 1 {
		return strings.Compare(strings.ToLower(a.str), strings.ToLower(b.str))
	}
	return sign(a.num - b.num)
}

func (e *fenv) eval(n fnode) (fval, error) {
	switch n := n.(type) {
	case fLit:
		return n.v, nil
	case fRef:
		return e.ref(n)
	case fUnary:
		x, err := e.eval(n.x)
		if err != nil {
			return x, err
		}
		x = e.scalar(x)
		if x.kind == fErr {
			return x, nil
		}
		f, ok := toNum(x)
		if !ok {
			return ferr("#VALUE!"), nil
		}
		if n.op == "-" {
			f = -f
		}
		return fnum(f), nil
	case fBinary:
		a, err := e.eval(n.a)
		if err != nil {
			return a, err
		}
		b, err := e.eval(n.b)
		if err != nil {
			return b, err
		}
		a, b = e.scalar(a), e.scalar(b)
		if a.kind == fErr {
			return a, nil
		}
		if b.kind == fErr {
			return b, nil
		}
		switch n.op {
		case "&":
			return fstr(toStr(a) + toStr(b)), nil
		case "=", "<>", "<", ">", "<=", ">=":
			c := cmpVals(a, b)
			switch n.op {
			case "=":
				return fbool(c == 0), nil
			case "<>":
				return fbool(c != 0), nil
			case "<":
				return fbool(c < 0), nil
			case ">":
				return fbool(c > 0), nil
			case "<=":
				return fbool(c <= 0), nil
			}
			return fbool(c >= 0), nil
		}
		x, ok1 := toNum(a)
		y, ok2 := toNum(b)
		if !ok1 || !ok2 {
			return ferr("#VALUE!"), nil
		}
		switch n.op {
		case "+":
			return fnum(x + y), nil
		case "-":
			return fnum(x - y), nil
		case "*":
			return fnum(x * y), nil
		case "/":
			if y == 0 {
				return ferr("#DIV/0!"), nil
			}
			return fnum(x / y), nil
		case "^":
			return fnum(math.Pow(x, y)), nil
		}
	case fCall:
		return e.call(n)
	}
	return fval{}, errUnsupported
}

// nums collects the numbers of arguments: those of ranges, and arguments
// that are numbers or read as numbers.
func (e *fenv) nums(args []fnode, fn func(float64)) (fval, error) {
	for _, a := range args {
		v, err := e.eval(a)
		if err != nil {
			return v, err
		}
		switch v.kind {
		case fRange:
			if v.rng.size() > 1<<22 && v.rng.r1-v.rng.r0 < maxRows-1 && v.rng.c1-v.rng.c0 < maxCols-1 {
				return fval{}, errUnsupported
			}
			v.rng.values(func(x fval) {
				if x.kind == fNum {
					fn(x.num)
				}
			})
		case fErr:
			return v, nil
		default:
			if f, ok := toNum(v); ok {
				fn(f)
			} else {
				return ferr("#VALUE!"), nil
			}
		}
	}
	return fval{}, nil
}

func (e *fenv) call(n fCall) (fval, error) {
	arg := func(i int) (fval, error) {
		if i >= len(n.args) {
			return fval{}, nil
		}
		v, err := e.eval(n.args[i])
		return e.scalar(v), err
	}
	numArg := func(i int) (float64, fval, error) {
		v, err := arg(i)
		if err != nil || v.kind == fErr {
			return 0, v, err
		}
		f, ok := toNum(v)
		if !ok {
			return 0, ferr("#VALUE!"), nil
		}
		return f, fval{kind: fNum}, nil
	}
	strArg := func(i int) (string, fval, error) {
		v, err := arg(i)
		if err != nil || v.kind == fErr {
			return "", v, err
		}
		return toStr(v), fval{kind: fStr}, nil
	}
	switch n.name {
	case "AND", "OR":
		res := n.name == "AND"
		for _, a := range n.args {
			v, err := e.eval(a)
			if err != nil {
				return v, err
			}
			v = e.scalar(v)
			if v.kind == fErr {
				return v, nil
			}
			b, ok := truthy(v)
			if !ok {
				return ferr("#VALUE!"), nil
			}
			if n.name == "AND" {
				res = res && b
			} else {
				res = res || b
			}
		}
		return fbool(res), nil
	case "NOT":
		v, err := arg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		b, ok := truthy(v)
		if !ok {
			return ferr("#VALUE!"), nil
		}
		return fbool(!b), nil
	case "IF":
		c, err := arg(0)
		if err != nil || c.kind == fErr {
			return c, err
		}
		b, _ := truthy(c)
		if b {
			return arg(1)
		}
		if len(n.args) < 3 {
			return fbool(false), nil
		}
		return arg(2)
	case "IFERROR":
		v, err := arg(0)
		if err != nil {
			return v, err
		}
		if v.kind == fErr {
			return arg(1)
		}
		return v, nil
	case "ISBLANK", "ISNUMBER", "ISTEXT", "ISNONTEXT", "ISERROR", "ISERR", "ISNA", "ISLOGICAL":
		v, err := arg(0)
		if err != nil {
			return v, err
		}
		switch n.name {
		case "ISBLANK":
			return fbool(v.kind == fBlank), nil
		case "ISNUMBER":
			return fbool(v.kind == fNum), nil
		case "ISTEXT":
			return fbool(v.kind == fStr), nil
		case "ISNONTEXT":
			return fbool(v.kind != fStr), nil
		case "ISERROR":
			return fbool(v.kind == fErr), nil
		case "ISERR":
			return fbool(v.kind == fErr && v.str != "#N/A"), nil
		case "ISNA":
			return fbool(v.kind == fErr && v.str == "#N/A"), nil
		}
		return fbool(v.kind == fBool), nil
	case "ISEVEN", "ISODD":
		f, v, err := numArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		even := int64(math.Trunc(f))%2 == 0
		return fbool(even == (n.name == "ISEVEN")), nil
	case "ROW", "COLUMN":
		if len(n.args) == 0 {
			if n.name == "ROW" {
				return fnum(float64(e.row + 1)), nil
			}
			return fnum(float64(e.col + 1)), nil
		}
		v, err := e.eval(n.args[0])
		if err != nil {
			return v, err
		}
		ref, ok := n.args[0].(fRef)
		if !ok {
			return ferr("#VALUE!"), nil
		}
		r, c := ref.r0, ref.c0
		if !ref.absR0 {
			r += e.dr
		}
		if !ref.absC0 {
			c += e.dc
		}
		if n.name == "ROW" {
			return fnum(float64(r + 1)), nil
		}
		return fnum(float64(c + 1)), nil
	case "MOD":
		a, v, err := numArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		b, v, err := numArg(1)
		if err != nil || v.kind == fErr {
			return v, err
		}
		if b == 0 {
			return ferr("#DIV/0!"), nil
		}
		return fnum(a - b*math.Floor(a/b)), nil
	case "ABS", "INT", "ROUND", "ROUNDUP", "ROUNDDOWN", "SIGN":
		a, v, err := numArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		switch n.name {
		case "ABS":
			return fnum(math.Abs(a)), nil
		case "INT":
			return fnum(math.Floor(a)), nil
		case "SIGN":
			return fnum(float64(sign(a))), nil
		}
		d, v, err := numArg(1)
		if err != nil || v.kind == fErr {
			return v, err
		}
		p := math.Pow(10, math.Trunc(d))
		switch n.name {
		case "ROUNDUP":
			return fnum(math.Copysign(math.Ceil(math.Abs(a)*p-1e-9), a) / p), nil
		case "ROUNDDOWN":
			return fnum(math.Trunc(a*p) / p), nil
		}
		return fnum(roundHalfAway(a, int(d))), nil
	case "LEN", "UPPER", "LOWER", "TRIM", "LEFT", "RIGHT", "MID", "VALUE":
		s, v, err := strArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		rs := []rune(s)
		switch n.name {
		case "LEN":
			return fnum(float64(len(rs))), nil
		case "UPPER":
			return fstr(strings.ToUpper(s)), nil
		case "LOWER":
			return fstr(strings.ToLower(s)), nil
		case "TRIM":
			return fstr(strings.Join(strings.Fields(s), " ")), nil
		case "VALUE":
			if f, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
				return fnum(f), nil
			}
			return ferr("#VALUE!"), nil
		case "MID":
			start, v, err := numArg(1)
			if err != nil || v.kind == fErr {
				return v, err
			}
			cnt, v, err := numArg(2)
			if err != nil || v.kind == fErr {
				return v, err
			}
			i := max(int(start)-1, 0)
			if i >= len(rs) {
				return fstr(""), nil
			}
			return fstr(string(rs[i:min(len(rs), i+max(int(cnt), 0))])), nil
		}
		cnt := 1.0
		if len(n.args) > 1 {
			c, v, err := numArg(1)
			if err != nil || v.kind == fErr {
				return v, err
			}
			cnt = c
		}
		k := min(max(int(cnt), 0), len(rs))
		if n.name == "LEFT" {
			return fstr(string(rs[:k])), nil
		}
		return fstr(string(rs[len(rs)-k:])), nil
	case "EXACT":
		a, v, err := strArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		b, v, err := strArg(1)
		if err != nil || v.kind == fErr {
			return v, err
		}
		return fbool(a == b), nil
	case "SEARCH", "FIND":
		needle, v, err := strArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		hay, v, err := strArg(1)
		if err != nil || v.kind == fErr {
			return v, err
		}
		if n.name == "SEARCH" {
			needle, hay = strings.ToLower(needle), strings.ToLower(hay)
		}
		i := strings.Index(hay, needle)
		if i < 0 {
			return ferr("#VALUE!"), nil
		}
		return fnum(float64(len([]rune(hay[:i])) + 1)), nil
	case "SUM", "AVERAGE", "MIN", "MAX", "COUNT":
		sum, cnt := 0.0, 0
		lo, hi := math.Inf(1), math.Inf(-1)
		v, err := e.nums(n.args, func(f float64) {
			sum += f
			cnt++
			lo, hi = math.Min(lo, f), math.Max(hi, f)
		})
		if err != nil || v.kind == fErr {
			return v, err
		}
		switch n.name {
		case "SUM":
			return fnum(sum), nil
		case "COUNT":
			return fnum(float64(cnt)), nil
		case "AVERAGE":
			if cnt == 0 {
				return ferr("#DIV/0!"), nil
			}
			return fnum(sum / float64(cnt)), nil
		}
		if cnt == 0 {
			return fnum(0), nil
		}
		if n.name == "MIN" {
			return fnum(lo), nil
		}
		return fnum(hi), nil
	case "COUNTA", "COUNTBLANK":
		cnt, total := 0, 0
		for _, a := range n.args {
			v, err := e.eval(a)
			if err != nil {
				return v, err
			}
			if v.kind == fRange {
				total += v.rng.size()
				v.rng.values(func(x fval) {
					if x.kind != fStr || x.str != "" {
						cnt++
					}
				})
				continue
			}
			total++
			if v.kind != fBlank {
				cnt++
			}
		}
		if n.name == "COUNTA" {
			return fnum(float64(cnt)), nil
		}
		return fnum(float64(total - cnt)), nil
	case "COUNTIF":
		if len(n.args) != 2 {
			return fval{}, errUnsupported
		}
		rv, err := e.eval(n.args[0])
		if err != nil {
			return rv, err
		}
		crit, err := arg(1)
		if err != nil {
			return crit, err
		}
		match := criterion(crit)
		cnt := 0
		switch rv.kind {
		case fRange:
			if rv.rng.size() > 1<<22 && rv.rng.r1-rv.rng.r0 < maxRows-1 && rv.rng.c1-rv.rng.c0 < maxCols-1 {
				return fval{}, errUnsupported
			}
			rv.rng.values(func(x fval) {
				if match(x) {
					cnt++
				}
			})
		default:
			if match(rv) {
				cnt++
			}
		}
		return fnum(float64(cnt)), nil
	case "YEAR", "MONTH", "DAY", "WEEKDAY", "HOUR", "MINUTE":
		f, v, err := numArg(0)
		if err != nil || v.kind == fErr {
			return v, err
		}
		t, leap, ok := serialTime(f, e.s.c.date1904)
		if !ok {
			return ferr("#NUM!"), nil
		}
		switch n.name {
		case "YEAR":
			return fnum(float64(t.Year())), nil
		case "MONTH":
			if leap {
				return fnum(2), nil
			}
			return fnum(float64(t.Month())), nil
		case "DAY":
			if leap {
				return fnum(29), nil
			}
			return fnum(float64(t.Day())), nil
		case "HOUR":
			return fnum(float64(t.Hour())), nil
		case "MINUTE":
			return fnum(float64(t.Minute())), nil
		}
		return fnum(float64(t.Weekday()) + 1), nil
	}
	return fval{}, errUnsupported
}

// criterion makes a COUNTIF criterion: a value, or a comparison ("<5",
// "<>done") with * and ? wildcards in text.
func criterion(c fval) func(fval) bool {
	if c.kind != fStr {
		return func(v fval) bool { return v.kind != fBlank && cmpVals(v, c) == 0 }
	}
	s := c.str
	op := "="
	for _, o := range []string{"<>", "<=", ">=", "=", "<", ">"} {
		if strings.HasPrefix(s, o) {
			op, s = o, s[len(o):]
			break
		}
	}
	target := fstr(s)
	if f, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
		target = fnum(f)
	}
	wild := target.kind == fStr && strings.ContainsAny(s, "*?")
	return func(v fval) bool {
		if v.kind == fBlank {
			return op == "<>" && s != "" || op == "=" && s == ""
		}
		var c int
		if wild && v.kind == fStr {
			if wildMatch(strings.ToLower(s), strings.ToLower(v.str)) {
				c = 0
			} else {
				c = 1
			}
		} else if (target.kind == fNum) != (v.kind == fNum) {
			return op == "<>"
		} else {
			c = cmpVals(v, target)
		}
		switch op {
		case "=":
			return c == 0
		case "<>":
			return c != 0
		case "<":
			return c < 0
		case ">":
			return c > 0
		case "<=":
			return c <= 0
		}
		return c >= 0
	}
}

func wildMatch(pat, s string) bool {
	p, t := []rune(pat), []rune(s)
	var match func(i, j int) bool
	match = func(i, j int) bool {
		for i < len(p) {
			switch p[i] {
			case '*':
				for k := j; k <= len(t); k++ {
					if match(i+1, k) {
						return true
					}
				}
				return false
			case '?':
				if j >= len(t) {
					return false
				}
			default:
				if j >= len(t) || p[i] != t[j] {
					return false
				}
			}
			i++
			j++
		}
		return j == len(t)
	}
	return match(0, 0)
}

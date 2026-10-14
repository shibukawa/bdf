package pdfme

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// The text of a read-only field of a template of pdfme may hold
// expressions in braces ("{total}", "{Number(price) * qty}", "{currentPage}
// / {totalPages}"): JavaScript expressions over the values of the data.
// pdfme evaluates them with a parser of JavaScript and a list of what they
// may use. Here they are evaluated by an interpreter of that part of the
// language: literals, the operators, member access, calls of the functions
// and methods listed below, array and object literals and arrow
// functions. An expression of anything else (or one that fails) is left as
// its text, as pdfme leaves an expression it cannot evaluate.

// Limits of an expression: its length, nesting, evaluation steps, and
// temporary values (a template is not trusted).
const (
	maxExprLen   = 4096
	maxExprDepth = 64
	maxExprSteps = 1 << 20
	maxExprText  = 1 << 20
	maxExprItems = 1 << 16
	maxExprAlloc = 16 << 20
	maxEvalDepth = 256
)

// node is a node of the tree of an expression.
type node struct {
	kind  byte
	op    string
	val   any
	name  string
	kids  []*node
	parms []string
}

const (
	nLit = iota
	nIdent
	nUnary
	nBinary
	nCond
	nMember
	nIndex
	nCall
	nArray
	nObject
	nArrow
)

type exprParser struct {
	s     string
	pos   int
	depth int
}

func (p *exprParser) fail(format string, args ...any) {
	panic(exprError(fmt.Sprintf(format, args...)))
}

type exprError string

func (e exprError) Error() string { return string(e) }

func (p *exprParser) ws() {
	for p.pos < len(p.s) && (p.s[p.pos] == ' ' || p.s[p.pos] == '\t' || p.s[p.pos] == '\n' || p.s[p.pos] == '\r') {
		p.pos++
	}
}

func (p *exprParser) peek(tok string) bool {
	p.ws()
	return strings.HasPrefix(p.s[p.pos:], tok)
}

func (p *exprParser) eat(tok string) bool {
	if p.peek(tok) {
		p.pos += len(tok)
		return true
	}
	return false
}

func (p *exprParser) want(tok string) {
	if !p.eat(tok) {
		p.fail("%q is expected at %d", tok, p.pos)
	}
}

func identStart(c rune) bool { return c == '_' || c == '$' || unicode.IsLetter(c) }
func identPart(c rune) bool  { return identStart(c) || unicode.IsDigit(c) }

func (p *exprParser) ident() (string, bool) {
	p.ws()
	start := p.pos
	for p.pos < len(p.s) {
		c, n := utf8.DecodeRuneInString(p.s[p.pos:])
		if p.pos == start && !identStart(c) || p.pos > start && !identPart(c) {
			break
		}
		p.pos += n
	}
	return p.s[start:p.pos], p.pos > start
}

// parseExpr parses an expression whole.
func parseExpr(s string) (n *node, err error) {
	if len(s) > maxExprLen {
		return nil, exprError("the expression is too long")
	}
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(exprError)
			if !ok {
				panic(r)
			}
			n, err = nil, e
		}
	}()
	p := &exprParser{s: s}
	n = p.expr()
	p.ws()
	if p.pos != len(p.s) {
		p.fail("%q is not expected at %d", p.s[p.pos:], p.pos)
	}
	return n, nil
}

func (p *exprParser) expr() *node {
	if p.depth++; p.depth > maxExprDepth {
		p.fail("the expression is nested too deep")
	}
	defer func() { p.depth-- }()
	if n := p.arrow(); n != nil {
		return n
	}
	cond := p.binary(0)
	if p.eat("?") {
		a := p.expr()
		p.want(":")
		return &node{kind: nCond, kids: []*node{cond, a, p.expr()}}
	}
	return cond
}

// arrow parses an arrow function, when one starts here: "x => …" or
// "(a, b) => …".
func (p *exprParser) arrow() *node {
	save := p.pos
	var parms []string
	if p.eat("(") {
		for !p.eat(")") {
			name, ok := p.ident()
			if !ok {
				p.pos = save
				return nil
			}
			parms = append(parms, name)
			if !p.eat(",") && !p.peek(")") {
				p.pos = save
				return nil
			}
		}
	} else if name, ok := p.ident(); ok {
		parms = []string{name}
	} else {
		p.pos = save
		return nil
	}
	if !p.eat("=>") {
		p.pos = save
		return nil
	}
	return &node{kind: nArrow, parms: parms, kids: []*node{p.expr()}}
}

// The binary operators, by how tightly they bind.
var binaryOps = [][]string{{"||"}, {"&&"}, {"===", "!==", "==", "!="}, {"<=", ">=", "<", ">"}, {"+", "-"}, {"*", "/", "%"}}

func (p *exprParser) binary(level int) *node {
	if level == len(binaryOps) {
		return p.power()
	}
	left := p.binary(level + 1)
	for {
		found := ""
		for _, op := range binaryOps[level] {
			// (not an arrow, nor the start of a longer operator)
			if p.peek(op) && !p.peek("=>") && !(op == "*" && p.peek("**")) {
				found = op
				break
			}
		}
		if found == "" {
			return left
		}
		p.pos += len(found)
		left = &node{kind: nBinary, op: found, kids: []*node{left, p.binary(level + 1)}}
	}
}

func (p *exprParser) power() *node {
	base := p.unary()
	if p.eat("**") {
		return &node{kind: nBinary, op: "**", kids: []*node{base, p.power()}}
	}
	return base
}

func (p *exprParser) unary() *node {
	for _, op := range []string{"!", "-", "+"} {
		if p.eat(op) {
			if p.depth++; p.depth > maxExprDepth {
				p.fail("the expression is nested too deep")
			}
			n := &node{kind: nUnary, op: op, kids: []*node{p.unary()}}
			p.depth--
			return n
		}
	}
	return p.postfix()
}

func (p *exprParser) postfix() *node {
	n := p.primary()
	for {
		switch {
		case p.eat("."):
			name, ok := p.ident()
			if !ok {
				p.fail("a name is expected at %d", p.pos)
			}
			n = &node{kind: nMember, name: name, kids: []*node{n}}
		case p.eat("["):
			n = &node{kind: nIndex, kids: []*node{n, p.expr()}}
			p.want("]")
		case p.eat("("):
			call := &node{kind: nCall, kids: []*node{n}}
			for !p.eat(")") {
				call.kids = append(call.kids, p.expr())
				if !p.eat(",") && !p.peek(")") {
					p.fail("\",\" or \")\" is expected at %d", p.pos)
				}
			}
			n = call
		default:
			return n
		}
	}
}

func (p *exprParser) primary() *node {
	p.ws()
	if p.pos >= len(p.s) {
		p.fail("the expression ends early")
	}
	c := p.s[p.pos]
	switch {
	case c == '(':
		p.pos++
		n := p.expr()
		p.want(")")
		return n
	case c == '[':
		p.pos++
		n := &node{kind: nArray}
		for !p.eat("]") {
			n.kids = append(n.kids, p.expr())
			if !p.eat(",") && !p.peek("]") {
				p.fail("\",\" or \"]\" is expected at %d", p.pos)
			}
		}
		return n
	case c == '{':
		p.pos++
		n := &node{kind: nObject}
		for !p.eat("}") {
			p.ws()
			var key string
			if p.pos < len(p.s) && (p.s[p.pos] == '"' || p.s[p.pos] == '\'') {
				key = p.str()
			} else if name, ok := p.ident(); ok {
				key = name
			} else {
				p.fail("a name is expected at %d", p.pos)
			}
			p.want(":")
			n.parms = append(n.parms, key)
			n.kids = append(n.kids, p.expr())
			if !p.eat(",") && !p.peek("}") {
				p.fail("\",\" or \"}\" is expected at %d", p.pos)
			}
		}
		return n
	case c == '"' || c == '\'':
		return &node{kind: nLit, val: p.str()}
	case c >= '0' && c <= '9' || c == '.':
		start := p.pos
		for p.pos < len(p.s) && (p.s[p.pos] >= '0' && p.s[p.pos] <= '9' || p.s[p.pos] == '.' ||
			(p.s[p.pos] == 'e' || p.s[p.pos] == 'E') && p.pos > start || (p.s[p.pos] == '+' || p.s[p.pos] == '-') && (p.s[p.pos-1] == 'e' || p.s[p.pos-1] == 'E')) {
			p.pos++
		}
		v, err := strconv.ParseFloat(p.s[start:p.pos], 64)
		if err != nil {
			p.fail("%q is not a number", p.s[start:p.pos])
		}
		return &node{kind: nLit, val: v}
	}
	name, ok := p.ident()
	if !ok {
		p.fail("%q is not expected at %d", p.s[p.pos:], p.pos)
	}
	switch name {
	case "true":
		return &node{kind: nLit, val: true}
	case "false":
		return &node{kind: nLit, val: false}
	case "null", "undefined":
		return &node{kind: nLit, val: nil}
	case "new", "typeof", "function", "this", "delete", "void", "in", "instanceof", "await", "yield", "class":
		p.fail("%q is not of the expressions that are evaluated", name)
	}
	return &node{kind: nIdent, name: name}
}

func (p *exprParser) str() string {
	quote := p.s[p.pos]
	p.pos++
	var b strings.Builder
	for p.pos < len(p.s) && p.s[p.pos] != quote {
		c := p.s[p.pos]
		if c == '\\' && p.pos+1 < len(p.s) {
			p.pos++
			switch e := p.s[p.pos]; e {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			default:
				b.WriteByte(e)
			}
			p.pos++
			continue
		}
		b.WriteByte(c)
		p.pos++
	}
	if p.pos >= len(p.s) {
		p.fail("a text is not closed")
	}
	p.pos++
	return b.String()
}

// --- evaluation ---

// closure is an arrow function with the names around it.
type closure struct {
	n   *node
	env *exprEnv
}

// builtin is a function of the language.
type builtin func(args []any) any

type exprEnv struct {
	vars   map[string]any
	parent *exprEnv
}

func (e *exprEnv) lookup(name string) (any, bool) {
	for ; e != nil; e = e.parent {
		if v, ok := e.vars[name]; ok {
			return v, true
		}
	}
	return nil, false
}

type evaluator struct {
	steps, depth, allocated int
}

// allocate bounds cumulative temporary results as well as each result.
// A step limit alone does not stop a few concat or replaceAll calls from
// allocating exponentially more than the template holds.
func (ev *evaluator) allocate(v any) {
	n := 0
	switch x := v.(type) {
	case string:
		checkTextLen(len(x))
		n = len(x)
	case []any:
		checkItems(len(x))
		n = 16 * len(x)
	case map[string]any:
		checkItems(len(x))
		n = 32 * len(x)
	}
	if n > maxExprAlloc-ev.allocated {
		panic(exprError("the expression allocates too much"))
	}
	ev.allocated += n
}

func checkTextLen(n int) {
	if n > maxExprText {
		panic(exprError("a text of the expression is too long"))
	}
}

func checkItems(n int) {
	if n > maxExprItems {
		panic(exprError("an array or object of the expression is too large"))
	}
}

func (ev *evaluator) step() {
	if ev.steps++; ev.steps > maxExprSteps {
		panic(exprError("the expression takes too long"))
	}
}

// scan charges linear work inside builtins to the same budget as AST
// evaluation. A builtin returning an empty text can still scan a large
// input on every iteration, without allocating anything.
func (ev *evaluator) scan(n int) {
	if n > maxExprSteps-ev.steps {
		panic(exprError("the expression takes too long"))
	}
	ev.steps += n
}

// Builtins coerce nested values too (Number([text]), for example). Count
// each occurrence, so wrapping a large string cannot hide its scan cost.
func (ev *evaluator) scanValue(v any, depth int) {
	if depth > maxExprDepth {
		panic(exprError("a value of the expression is nested too deep"))
	}
	ev.scan(1)
	switch x := v.(type) {
	case string:
		ev.scan(len(x))
	case []any:
		for _, child := range x {
			ev.scanValue(child, depth+1)
		}
	case map[string]any:
		for key, child := range x {
			ev.scan(len(key))
			ev.scanValue(child, depth+1)
		}
	}
}

// evalExpr evaluates an expression with the values of names.
func evalExpr(n *node, vars map[string]any) (v any, err error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(exprError)
			if !ok {
				panic(r)
			}
			v, err = nil, e
		}
	}()
	ev := &evaluator{}
	v = ev.eval(n, &exprEnv{vars: vars, parent: &exprEnv{vars: exprGlobals}})
	checkValue(v)
	return v, nil
}

// checkValue bounds a value's expanded representation before JSON encoding
// or returning it. Counting each occurrence also bounds shared subarrays,
// which can otherwise describe an enormous output with very little memory.
func checkValue(v any) {
	remaining := maxExprText
	var visit func(any, int)
	spend := func(n int) {
		if n > remaining {
			panic(exprError("a value of the expression is too large"))
		}
		remaining -= n
	}
	visit = func(v any, depth int) {
		if depth > maxExprDepth {
			panic(exprError("a value of the expression is nested too deep"))
		}
		switch x := v.(type) {
		case string:
			spend(len(x))
		case []any:
			checkItems(len(x))
			spend(len(x) + 2)
			for _, e := range x {
				visit(e, depth+1)
			}
		case map[string]any:
			checkItems(len(x))
			spend(len(x) + 2)
			for k, e := range x {
				spend(len(k))
				visit(e, depth+1)
			}
		default:
			spend(1)
		}
	}
	visit(v, 0)
}

func truthy(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case float64:
		return x != 0 && !math.IsNaN(x)
	case string:
		return x != ""
	}
	return true
}

// toNumber is Number(v).
func toNumber(v any) float64 { return numberAt(v, 0) }

func numberAt(v any, depth int) float64 {
	if depth > maxExprDepth {
		panic(exprError("a value of the expression is nested too deep"))
	}
	switch x := v.(type) {
	case nil:
		return 0
	case bool:
		if x {
			return 1
		}
		return 0
	case float64:
		return x
	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return 0
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return math.NaN()
		}
		return f
	case []any:
		if len(x) == 0 {
			return 0
		}
		if len(x) == 1 {
			return numberAt(x[0], depth+1)
		}
	}
	return math.NaN()
}

// toText is String(v).
func toText(v any) string {
	if s, ok := v.(string); ok {
		checkTextLen(len(s))
		return s
	}
	var b strings.Builder
	appendText(&b, v, 0)
	return b.String()
}

// appendText checks before growing the buffer, including nested arrays
// whose entries can refer to the same large value many times.
func appendText(b *strings.Builder, v any, depth int) {
	if depth > maxExprDepth {
		panic(exprError("a value of the expression is nested too deep"))
	}
	write := func(s string) {
		if len(s) > maxExprText-b.Len() {
			panic(exprError("a text of the expression is too long"))
		}
		b.WriteString(s)
	}
	switch x := v.(type) {
	case nil:
		return
	case bool:
		write(strconv.FormatBool(x))
	case float64:
		switch {
		case math.IsNaN(x):
			write("NaN")
		case math.IsInf(x, 1):
			write("Infinity")
		case math.IsInf(x, -1):
			write("-Infinity")
		case x == math.Trunc(x) && math.Abs(x) < 1e21:
			write(strconv.FormatFloat(x, 'f', 0, 64))
		default:
			write(strconv.FormatFloat(x, 'g', -1, 64))
		}
	case string:
		write(x)
	case []any:
		checkItems(len(x))
		for i, e := range x {
			if i > 0 {
				write(",")
			}
			appendText(b, e, depth+1)
		}
	case map[string]any:
		write("[object Object]")
	}
}

func looseEqual(a, b any) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case string, bool:
			return x == toNumber(y)
		}
	case string:
		switch y := b.(type) {
		case string:
			return x == y
		case float64, bool:
			return toNumber(x) == toNumber(y)
		}
	case bool:
		if y, ok := b.(bool); ok {
			return x == y
		}
		return b != nil && toNumber(x) == toNumber(b)
	}
	return false
}

func strictEqual(a, b any) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case float64:
		y, ok := b.(float64)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	}
	return false
}

func (ev *evaluator) eval(n *node, env *exprEnv) any {
	ev.step()
	if ev.depth++; ev.depth > maxEvalDepth {
		panic(exprError("the expression evaluation is nested too deep"))
	}
	defer func() { ev.depth-- }()
	switch n.kind {
	case nLit:
		return n.val
	case nIdent:
		v, ok := env.lookup(n.name)
		if !ok {
			panic(exprError(n.name + " is not defined"))
		}
		return v
	case nUnary:
		v := ev.eval(n.kids[0], env)
		switch n.op {
		case "!":
			return !truthy(v)
		case "-":
			return -toNumber(v)
		}
		return toNumber(v)
	case nCond:
		if truthy(ev.eval(n.kids[0], env)) {
			return ev.eval(n.kids[1], env)
		}
		return ev.eval(n.kids[2], env)
	case nBinary:
		a := ev.eval(n.kids[0], env)
		switch n.op {
		case "&&":
			if !truthy(a) {
				return a
			}
			return ev.eval(n.kids[1], env)
		case "||":
			if truthy(a) {
				return a
			}
			return ev.eval(n.kids[1], env)
		}
		b := ev.eval(n.kids[1], env)
		switch n.op {
		case "+":
			_, as := a.(string)
			_, bs := b.(string)
			_, aa := a.([]any)
			_, ba := b.([]any)
			if as || bs || aa || ba {
				var bld strings.Builder
				appendText(&bld, a, 0)
				appendText(&bld, b, 0)
				s := bld.String()
				ev.allocate(s)
				return s
			}
			return toNumber(a) + toNumber(b)
		case "-":
			return toNumber(a) - toNumber(b)
		case "*":
			return toNumber(a) * toNumber(b)
		case "/":
			return toNumber(a) / toNumber(b)
		case "%":
			return math.Mod(toNumber(a), toNumber(b))
		case "**":
			return math.Pow(toNumber(a), toNumber(b))
		case "==":
			return looseEqual(a, b)
		case "!=":
			return !looseEqual(a, b)
		case "===":
			return strictEqual(a, b)
		case "!==":
			return !strictEqual(a, b)
		}
		// the comparisons: of texts as texts, of the rest as numbers
		if x, ok := a.(string); ok {
			if y, ok := b.(string); ok {
				switch n.op {
				case "<":
					return x < y
				case "<=":
					return x <= y
				case ">":
					return x > y
				}
				return x >= y
			}
		}
		x, y := toNumber(a), toNumber(b)
		switch n.op {
		case "<":
			return x < y
		case "<=":
			return x <= y
		case ">":
			return x > y
		}
		return x >= y
	case nArray:
		checkItems(len(n.kids))
		out := make([]any, len(n.kids))
		ev.allocate(out)
		for i, k := range n.kids {
			out[i] = ev.eval(k, env)
		}
		return out
	case nObject:
		out := map[string]any{}
		for i, k := range n.kids {
			out[n.parms[i]] = ev.eval(k, env)
		}
		ev.allocate(out)
		return out
	case nArrow:
		return &closure{n: n, env: env}
	case nMember:
		return ev.member(ev.eval(n.kids[0], env), n.name)
	case nIndex:
		obj, key := ev.eval(n.kids[0], env), ev.eval(n.kids[1], env)
		switch x := obj.(type) {
		case []any:
			if i, ok := key.(float64); ok {
				if i >= 0 && int(i) < len(x) && i == math.Trunc(i) {
					return x[int(i)]
				}
				return nil
			}
		case string:
			if i, ok := key.(float64); ok {
				checkTextLen(len(x))
				if i >= 0 && i < float64(len(x)) && i == math.Trunc(i) {
					ev.scan(len(x))
					if s := runeAt(x, int(i)); s != "" {
						return s
					}
				}
				return nil
			}
		}
		return ev.member(obj, toText(key))
	case nCall:
		fn := ev.eval(n.kids[0], env)
		args := make([]any, len(n.kids)-1)
		for i, k := range n.kids[1:] {
			args[i] = ev.eval(k, env)
		}
		return ev.call(fn, args)
	}
	panic(exprError("the expression is not understood"))
}

func (ev *evaluator) call(fn any, args []any) any {
	ev.step()
	switch f := fn.(type) {
	case builtin:
		// Global functions also scan their arguments: Number/parseFloat,
		// String and JSON can do substantial work without producing a
		// large result. Include nested values used by coercion and JSON.
		for _, a := range args {
			ev.scanValue(a, 0)
		}
		v := f(args)
		ev.allocate(v)
		return v
	case *closure:
		vars := make(map[string]any, len(f.n.parms))
		for i, p := range f.n.parms {
			if i < len(args) {
				vars[p] = args[i]
			} else {
				vars[p] = nil
			}
		}
		return ev.eval(f.n.kids[0], &exprEnv{vars: vars, parent: f.env})
	}
	panic(exprError("what is called is not a function"))
}

func arg(args []any, i int) any {
	if i < len(args) {
		return args[i]
	}
	return nil
}

func intArg(args []any, i, def int) int {
	if i >= len(args) || args[i] == nil {
		return def
	}
	f := toNumber(args[i])
	if math.IsNaN(f) {
		return 0
	}
	return int(math.Max(math.Min(f, 1<<30), -(1 << 30)))
}

// span clamps the start and the end of a slice of n items, where a
// negative index counts from the end.
func span(args []any, n int) (int, int) {
	clamp := func(i int) int {
		if i < 0 {
			i += n
		}
		return max(0, min(i, n))
	}
	a, b := clamp(intArg(args, 0, 0)), clamp(intArg(args, 1, n))
	if b < a {
		b = a
	}
	return a, b
}

// member returns a property of a value, or one of its methods.
func (ev *evaluator) member(obj any, name string) (value any) {
	switch name {
	case "constructor", "__proto__", "prototype", "toLocaleString", "valueOf":
		panic(exprError(name + " is not of the expressions that are evaluated"))
	}
	switch x := obj.(type) {
	case map[string]any:
		return x[name]
	case string:
		checkTextLen(len(x))
		defer func() {
			fn, ok := value.(builtin)
			if !ok || name == "toString" || name == "repeat" {
				return
			}
			value = builtin(func(a []any) any {
				cost := len(x)
				if name == "charAt" {
					if i := intArg(a, 0, 0); i < 0 || i >= len(x) {
						cost = 0 // (runeAt rejects this without scanning)
					}
				}
				ev.scan(cost)
				return fn(a)
			})
		}()
		switch name {
		case "length":
			ev.scan(len(x))
			return float64(utf16Len(x))
		case "toUpperCase":
			return builtin(func([]any) any { return strings.ToUpper(x) })
		case "toLowerCase":
			return builtin(func([]any) any { return strings.ToLower(x) })
		case "trim":
			return builtin(func([]any) any { return strings.TrimSpace(x) })
		case "toString":
			return builtin(func([]any) any { return x })
		case "slice", "substring":
			return builtin(func(a []any) any {
				rs := []rune(x)
				s, e := span(a, len(rs))
				return string(rs[s:e])
			})
		case "charAt":
			return builtin(func(a []any) any {
				return runeAt(x, intArg(a, 0, 0))
			})
		case "includes":
			return builtin(func(a []any) any { return strings.Contains(x, toText(arg(a, 0))) })
		case "startsWith":
			return builtin(func(a []any) any { return strings.HasPrefix(x, toText(arg(a, 0))) })
		case "endsWith":
			return builtin(func(a []any) any { return strings.HasSuffix(x, toText(arg(a, 0))) })
		case "indexOf":
			return builtin(func(a []any) any { return float64(strings.Index(x, toText(arg(a, 0)))) })
		case "replace":
			return builtin(func(a []any) any { return replaceText(x, toText(arg(a, 0)), toText(arg(a, 1)), 1) })
		case "replaceAll":
			return builtin(func(a []any) any { return replaceText(x, toText(arg(a, 0)), toText(arg(a, 1)), -1) })
		case "split":
			return builtin(func(a []any) any {
				var parts []string
				if sep := toText(arg(a, 0)); arg(a, 0) == nil {
					parts = []string{x}
				} else {
					n := strings.Count(x, sep) + 1
					if sep == "" {
						n = utf8.RuneCountInString(x)
					}
					checkItems(n)
					parts = strings.Split(x, sep)
				}
				out := make([]any, len(parts))
				for i, p := range parts {
					out[i] = p
				}
				return out
			})
		case "padStart", "padEnd":
			return builtin(func(a []any) any {
				n, fill := intArg(a, 0, 0), " "
				if arg(a, 1) != nil {
					fill = toText(arg(a, 1))
				}
				length := utf8.RuneCountInString(x)
				if n > 1<<16 || fill == "" || n <= length {
					return x
				}
				pad := make([]rune, n-length)
				for i := 0; i < len(pad); {
					for _, r := range fill {
						pad[i] = r
						i++
						if i == len(pad) {
							break
						}
					}
				}
				if name == "padStart" {
					return string(pad) + x
				}
				return x + string(pad)
			})
		case "repeat":
			return builtin(func(a []any) any {
				n := intArg(a, 0, 0)
				if n < 0 || len(x) > 0 && n > maxExprText/len(x) {
					panic(exprError("a text of the expression is too long"))
				}
				return strings.Repeat(x, n)
			})
		}
	case float64:
		switch name {
		case "toFixed":
			return builtin(func(a []any) any { return strconv.FormatFloat(x, 'f', max(0, min(intArg(a, 0, 0), 100)), 64) })
		case "toString":
			return builtin(func([]any) any { return toText(x) })
		}
	case []any:
		checkItems(len(x))
		defer func() {
			if fn, ok := value.(builtin); ok && name != "at" {
				value = builtin(func(a []any) any {
					ev.scan(len(x))
					return fn(a)
				})
			}
		}()
		switch name {
		case "length":
			return float64(len(x))
		case "join":
			return builtin(func(a []any) any {
				sep := ","
				if arg(a, 0) != nil {
					sep = toText(arg(a, 0))
				}
				var b strings.Builder
				for i, e := range x {
					if i > 0 {
						appendText(&b, sep, 0)
					}
					appendText(&b, e, 0)
				}
				return b.String()
			})
		case "slice":
			return builtin(func(a []any) any { s, e := span(a, len(x)); return append([]any(nil), x[s:e]...) })
		case "concat":
			return builtin(func(a []any) any {
				n := len(x)
				for _, e := range a {
					if l, ok := e.([]any); ok {
						checkItems(len(l))
						n += len(l)
					} else {
						n++
					}
					checkItems(n)
				}
				out := make([]any, len(x), n)
				copy(out, x)
				for _, e := range a {
					if l, ok := e.([]any); ok {
						out = append(out, l...)
					} else {
						out = append(out, e)
					}
				}
				return out
			})
		case "includes":
			return builtin(func(a []any) any {
				for _, e := range x {
					if strictEqual(e, arg(a, 0)) {
						return true
					}
				}
				return false
			})
		case "indexOf":
			return builtin(func(a []any) any {
				for i, e := range x {
					if strictEqual(e, arg(a, 0)) {
						return float64(i)
					}
				}
				return -1.0
			})
		case "at":
			return builtin(func(a []any) any {
				i := intArg(a, 0, 0)
				if i < 0 {
					i += len(x)
				}
				if i >= 0 && i < len(x) {
					return x[i]
				}
				return nil
			})
		case "map", "filter", "some", "every", "find", "reduce":
			return builtin(func(a []any) any {
				fn := arg(a, 0)
				each := func(i int, extra ...any) any { return ev.call(fn, append(extra, x[i], float64(i), x)) }
				switch name {
				case "map":
					out := make([]any, len(x))
					for i := range x {
						out[i] = each(i)
					}
					return out
				case "filter":
					out := []any{}
					for i := range x {
						if truthy(each(i)) {
							out = append(out, x[i])
						}
					}
					return out
				case "some":
					for i := range x {
						if truthy(each(i)) {
							return true
						}
					}
					return false
				case "every":
					for i := range x {
						if !truthy(each(i)) {
							return false
						}
					}
					return true
				case "find":
					for i := range x {
						if truthy(each(i)) {
							return x[i]
						}
					}
					return nil
				}
				// reduce: from the value given, or from the first item
				i, acc := 0, arg(a, 1)
				if len(a) < 2 {
					if len(x) == 0 {
						panic(exprError("reduce of nothing, with no value to start from"))
					}
					i, acc = 1, x[0]
				}
				for ; i < len(x); i++ {
					acc = each(i, acc)
				}
				return acc
			})
		}
	}
	if obj == nil {
		panic(exprError("a property of nothing (" + name + ")"))
	}
	return nil
}

func runeAt(s string, at int) string {
	if at < 0 || at >= len(s) {
		return ""
	}
	for _, r := range s {
		if at == 0 {
			return string(r)
		}
		at--
	}
	return ""
}

// utf16Len is the length of a text as JavaScript counts it.
func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		if n++; r > 0xffff {
			n++
		}
	}
	return n
}

func numFn(f func(float64) float64) builtin {
	return func(a []any) any { return f(toNumber(arg(a, 0))) }
}

// exprGlobals are the names that every expression has.
var exprGlobals = map[string]any{
	"Math": map[string]any{
		"PI": math.Pi, "E": math.E,
		"round": numFn(func(x float64) float64 { return math.Floor(x + 0.5) }), "floor": numFn(math.Floor), "ceil": numFn(math.Ceil), "abs": numFn(math.Abs),
		"trunc": numFn(math.Trunc), "sqrt": numFn(math.Sqrt), "sign": numFn(func(x float64) float64 {
			switch {
			case x > 0:
				return 1
			case x < 0:
				return -1
			}
			return x
		}),
		"pow": builtin(func(a []any) any { return math.Pow(toNumber(arg(a, 0)), toNumber(arg(a, 1))) }),
		"min": builtin(func(a []any) any {
			m := math.Inf(1)
			for _, v := range a {
				m = math.Min(m, toNumber(v))
			}
			return m
		}),
		"max": builtin(func(a []any) any {
			m := math.Inf(-1)
			for _, v := range a {
				m = math.Max(m, toNumber(v))
			}
			return m
		}),
	},
	"Number":  builtin(func(a []any) any { return toNumber(arg(a, 0)) }),
	"String":  builtin(func(a []any) any { return toText(arg(a, 0)) }),
	"Boolean": builtin(func(a []any) any { return truthy(arg(a, 0)) }),
	"isNaN":   builtin(func(a []any) any { return math.IsNaN(toNumber(arg(a, 0))) }),
	"parseFloat": builtin(func(a []any) any {
		return parseFloatPrefix(strings.TrimSpace(toText(arg(a, 0))))
	}),
	"parseInt": builtin(func(a []any) any {
		s := strings.TrimSpace(toText(arg(a, 0)))
		end := 0
		for end < len(s) && (s[end] >= '0' && s[end] <= '9' || end == 0 && (s[end] == '-' || s[end] == '+')) {
			end++
		}
		f, err := strconv.ParseFloat(s[:end], 64)
		if err != nil {
			return math.NaN()
		}
		return f
	}),
	"Array": map[string]any{"isArray": builtin(func(a []any) any { _, ok := arg(a, 0).([]any); return ok })},
	"Object": map[string]any{
		"keys": builtin(func(a []any) any {
			out := []any{}
			for _, k := range sortedKeys(arg(a, 0)) {
				out = append(out, k)
			}
			return out
		}),
		"values": builtin(func(a []any) any {
			out := []any{}
			m, _ := arg(a, 0).(map[string]any)
			for _, k := range sortedKeys(arg(a, 0)) {
				out = append(out, m[k])
			}
			return out
		}),
		"entries": builtin(func(a []any) any {
			out := []any{}
			m, _ := arg(a, 0).(map[string]any)
			for _, k := range sortedKeys(arg(a, 0)) {
				out = append(out, []any{k, m[k]})
			}
			return out
		}),
	},
	"JSON": map[string]any{
		"stringify": builtin(func(a []any) any {
			checkValue(arg(a, 0))
			b, err := json.Marshal(arg(a, 0))
			if err != nil {
				panic(exprError("a value is not written as JSON"))
			}
			checkTextLen(len(b))
			return string(b)
		}),
		"parse": builtin(func(a []any) any {
			var v any
			if err := json.Unmarshal([]byte(toText(arg(a, 0))), &v); err != nil {
				panic(exprError("a text is not JSON"))
			}
			checkValue(v)
			return v
		}),
	},
}

func sortedKeys(v any) []string {
	m, _ := v.(map[string]any)
	checkItems(len(m))
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// replaceText measures the replacement before strings.Replace allocates it.
func replaceText(s, old, replacement string, count int) string {
	n := strings.Count(s, old)
	if count >= 0 {
		n = min(n, count)
	}
	if growth := len(replacement) - len(old); growth > 0 && n > (maxExprText-len(s))/growth {
		panic(exprError("a text of the expression is too long"))
	}
	return strings.Replace(s, old, replacement, count)
}

// parseFloatPrefix scans once for the numeric prefix. Trying ParseFloat on
// every prefix makes a long number take quadratic time.
func parseFloatPrefix(s string) float64 {
	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	if strings.HasPrefix(s[i:], "Infinity") {
		if len(s) > 0 && s[0] == '-' {
			return math.Inf(-1)
		}
		return math.Inf(1)
	}
	start := i
	digits := func() {
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
	}
	digits()
	n := i - start
	if i < len(s) && s[i] == '.' {
		i++
		start = i
		digits()
		n += i - start
	}
	if n == 0 {
		return math.NaN()
	}
	end := i
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		i++
		if i < len(s) && (s[i] == '+' || s[i] == '-') {
			i++
		}
		start = i
		digits()
		if i > start {
			end = i
		}
	}
	f, _ := strconv.ParseFloat(s[:end], 64)
	return f
}

func evalText(n *node, vars map[string]any) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(exprError); ok {
				out, err = "", e
			} else {
				panic(r)
			}
		}
	}()
	v, err := evalExpr(n, vars)
	if err != nil {
		return "", err
	}
	return toText(v), nil
}

// expand replaces the expressions in braces of a text with their values.
// An expression that is not evaluated stays as it is, with its braces;
// failed reports whether one did.
func expand(text string, vars map[string]any) (out string, failed bool) {
	if !strings.Contains(text, "{") || !strings.Contains(text, "}") {
		return text, false
	}
	if len(text) > maxExprText {
		return text, true
	}
	var b strings.Builder
	write := func(s string) bool {
		if len(s) > maxExprText-b.Len() {
			return false
		}
		b.WriteString(s)
		return true
	}
	for i := 0; i < len(text); {
		if text[i] != '{' {
			if b.Len() == maxExprText {
				return text, true
			}
			b.WriteByte(text[i])
			i++
			continue
		}
		// the brace that closes this one
		depth, end := 0, -1
		for j := i; j < len(text); j++ {
			if text[j] == '{' {
				depth++
			} else if text[j] == '}' {
				if depth--; depth == 0 {
					end = j
					break
				}
			}
		}
		if end < 0 {
			if !write(text[i:]) {
				return text, true
			}
			break
		}
		src := text[i+1 : end]
		n, err := parseExpr(strings.TrimSpace(src))
		var value string
		if err == nil {
			value, err = evalText(n, vars)
		}
		if err != nil {
			failed = true
			value = text[i : end+1]
		}
		if !write(value) {
			return text, true
		}
		i = end + 1
	}
	return b.String(), failed
}

// names collects the names that an expression reads and does not bind
// itself (as the parameters of its functions).
func (n *node) names(bound map[string]bool, out map[string]bool) {
	if n == nil {
		return
	}
	switch n.kind {
	case nIdent:
		if !bound[n.name] {
			out[n.name] = true
		}
		return
	case nArrow:
		inner := map[string]bool{}
		for k := range bound {
			inner[k] = true
		}
		for _, p := range n.parms {
			inner[p] = true
		}
		bound = inner
	}
	for _, k := range n.kids {
		k.names(bound, out)
	}
}

// exprNames returns the names that the expressions in braces of a text
// read, without those that every expression has, in their order.
func exprNames(text string) []string {
	if !strings.Contains(text, "{") || !strings.Contains(text, "}") {
		return nil
	}
	var out []string
	seen := map[string]bool{}
	for i := 0; i < len(text); i++ {
		if text[i] != '{' {
			continue
		}
		depth, end := 0, -1
		for j := i; j < len(text); j++ {
			if text[j] == '{' {
				depth++
			} else if text[j] == '}' {
				if depth--; depth == 0 {
					end = j
					break
				}
			}
		}
		if end < 0 {
			break
		}
		if n, err := parseExpr(strings.TrimSpace(text[i+1 : end])); err == nil {
			found := map[string]bool{}
			n.names(nil, found)
			for _, name := range sortedNames(found) {
				if _, global := exprGlobals[name]; !global && !seen[name] {
					seen[name] = true
					out = append(out, name)
				}
			}
		}
		i = end
	}
	return out
}

func sortedNames(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

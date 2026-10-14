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

// Limits of an expression: its length, the depth of its nesting, and the
// steps of its evaluation (a template is not trusted).
const (
	maxExprLen   = 4096
	maxExprDepth = 64
	maxExprSteps = 1 << 20
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
	steps int
}

func (ev *evaluator) step() {
	if ev.steps++; ev.steps > maxExprSteps {
		panic(exprError("the expression takes too long"))
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
	return ev.eval(n, &exprEnv{vars: vars, parent: &exprEnv{vars: exprGlobals}}), nil
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
func toNumber(v any) float64 {
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
			return toNumber(x[0])
		}
	}
	return math.NaN()
}

// toText is String(v).
func toText(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case bool:
		return strconv.FormatBool(x)
	case float64:
		switch {
		case math.IsNaN(x):
			return "NaN"
		case math.IsInf(x, 1):
			return "Infinity"
		case math.IsInf(x, -1):
			return "-Infinity"
		case x == math.Trunc(x) && math.Abs(x) < 1e21:
			return strconv.FormatFloat(x, 'f', 0, 64)
		}
		return strconv.FormatFloat(x, 'g', -1, 64)
	case string:
		return x
	case []any:
		parts := make([]string, len(x))
		for i, e := range x {
			parts[i] = toText(e)
		}
		return strings.Join(parts, ",")
	case map[string]any:
		return "[object Object]"
	}
	return ""
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
				s := toText(a) + toText(b)
				if len(s) > 1<<20 {
					panic(exprError("a text of the expression is too long"))
				}
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
		out := make([]any, len(n.kids))
		for i, k := range n.kids {
			out[i] = ev.eval(k, env)
		}
		return out
	case nObject:
		out := map[string]any{}
		for i, k := range n.kids {
			out[n.parms[i]] = ev.eval(k, env)
		}
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
				rs := []rune(x)
				if i >= 0 && int(i) < len(rs) && i == math.Trunc(i) {
					return string(rs[int(i)])
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
		return f(args)
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
func (ev *evaluator) member(obj any, name string) any {
	switch name {
	case "constructor", "__proto__", "prototype", "toLocaleString", "valueOf":
		panic(exprError(name + " is not of the expressions that are evaluated"))
	}
	switch x := obj.(type) {
	case map[string]any:
		return x[name]
	case string:
		rs := []rune(x)
		switch name {
		case "length":
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
			return builtin(func(a []any) any { s, e := span(a, len(rs)); return string(rs[s:e]) })
		case "charAt":
			return builtin(func(a []any) any {
				if i := intArg(a, 0, 0); i >= 0 && i < len(rs) {
					return string(rs[i])
				}
				return ""
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
			return builtin(func(a []any) any { return strings.Replace(x, toText(arg(a, 0)), toText(arg(a, 1)), 1) })
		case "replaceAll":
			return builtin(func(a []any) any { return strings.ReplaceAll(x, toText(arg(a, 0)), toText(arg(a, 1))) })
		case "split":
			return builtin(func(a []any) any {
				var parts []string
				if sep := toText(arg(a, 0)); arg(a, 0) == nil {
					parts = []string{x}
				} else {
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
				if n > 1<<16 || fill == "" || n <= len(rs) {
					return x
				}
				pad := []rune(strings.Repeat(fill, (n-len(rs))/len([]rune(fill))+1))[:n-len(rs)]
				if name == "padStart" {
					return string(pad) + x
				}
				return x + string(pad)
			})
		case "repeat":
			return builtin(func(a []any) any {
				n := intArg(a, 0, 0)
				if n < 0 || n*len(x) > 1<<20 {
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
		switch name {
		case "length":
			return float64(len(x))
		case "join":
			return builtin(func(a []any) any {
				sep := ","
				if arg(a, 0) != nil {
					sep = toText(arg(a, 0))
				}
				parts := make([]string, len(x))
				for i, e := range x {
					parts[i] = toText(e)
				}
				return strings.Join(parts, sep)
			})
		case "slice":
			return builtin(func(a []any) any { s, e := span(a, len(x)); return append([]any(nil), x[s:e]...) })
		case "concat":
			return builtin(func(a []any) any {
				out := append([]any(nil), x...)
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
		// the number that the text starts with
		s := strings.TrimSpace(toText(arg(a, 0)))
		end := 0
		for i := 1; i <= len(s); i++ {
			if _, err := strconv.ParseFloat(s[:i], 64); err == nil {
				end = i
			} else if c := s[i-1]; !(c >= '0' && c <= '9' || c == '.' || c == '-' || c == '+' || c == 'e' || c == 'E') {
				break
			}
		}
		if end == 0 {
			return math.NaN()
		}
		f, _ := strconv.ParseFloat(s[:end], 64)
		return f
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
			b, err := json.Marshal(arg(a, 0))
			if err != nil {
				panic(exprError("a value is not written as JSON"))
			}
			return string(b)
		}),
		"parse": builtin(func(a []any) any {
			var v any
			if err := json.Unmarshal([]byte(toText(arg(a, 0))), &v); err != nil {
				panic(exprError("a text is not JSON"))
			}
			return v
		}),
	},
}

func sortedKeys(v any) []string {
	m, _ := v.(map[string]any)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// expand replaces the expressions in braces of a text with their values.
// An expression that is not evaluated stays as it is, with its braces;
// failed reports whether one did.
func expand(text string, vars map[string]any) (out string, failed bool) {
	if !strings.Contains(text, "{") || !strings.Contains(text, "}") {
		return text, false
	}
	var b strings.Builder
	for i := 0; i < len(text); {
		if text[i] != '{' {
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
			b.WriteString(text[i:])
			break
		}
		src := text[i+1 : end]
		n, err := parseExpr(strings.TrimSpace(src))
		var v any
		if err == nil {
			v, err = evalExpr(n, vars)
		}
		if err != nil {
			failed = true
			b.WriteString(text[i : end+1])
		} else {
			b.WriteString(toText(v))
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

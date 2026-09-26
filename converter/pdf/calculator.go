package pdf

import (
	"math"
	"strconv"
)

// A PostScript calculator function (type 4) is compiled into a flat list
// of instructions: if and ifelse become conditional jumps.
type psOp struct {
	code psCode
	num  float64 // psPush
	to   int     // psJumpFalse, psJump
}

type psCode uint8

const (
	psPush psCode = iota
	psJumpFalse
	psJump
	psAbs
	psAdd
	psAtan
	psCeiling
	psCos
	psCvi
	psCvr
	psDiv
	psExp
	psFloor
	psIdiv
	psLn
	psLog
	psMod
	psMul
	psNeg
	psRound
	psSin
	psSqrt
	psSub
	psTruncate
	psAnd
	psBitshift
	psEq
	psFalse
	psGe
	psGt
	psLe
	psLt
	psNe
	psNot
	psOr
	psTrue
	psXor
	psCopy
	psDup
	psExch
	psIndex
	psPop
	psRoll
)

var psOps = map[string]psCode{
	"abs": psAbs, "add": psAdd, "atan": psAtan, "ceiling": psCeiling, "cos": psCos, "cvi": psCvi, "cvr": psCvr,
	"div": psDiv, "exp": psExp, "floor": psFloor, "idiv": psIdiv, "ln": psLn, "log": psLog, "mod": psMod,
	"mul": psMul, "neg": psNeg, "round": psRound, "sin": psSin, "sqrt": psSqrt, "sub": psSub, "truncate": psTruncate,
	"and": psAnd, "bitshift": psBitshift, "eq": psEq, "false": psFalse, "ge": psGe, "gt": psGt, "le": psLe,
	"lt": psLt, "ne": psNe, "not": psNot, "or": psOr, "true": psTrue, "xor": psXor,
	"copy": psCopy, "dup": psDup, "exch": psExch, "index": psIndex, "pop": psPop, "roll": psRoll,
}

// psNode is a parsed token: a number, an operator or a procedure.
type psNode struct {
	op   string
	num  float64
	isN  bool
	proc []psNode
}

// compilePS compiles a calculator program ({ … }). It returns nil when the
// program cannot be read.
func compilePS(src []byte) []psOp {
	toks := psTokens(src)
	i := 0
	for i < len(toks) && toks[i] != "{" {
		i++
	}
	if i == len(toks) {
		return nil
	}
	i++
	proc, ok := psParse(toks, &i, 0)
	if !ok {
		return nil
	}
	var out []psOp
	if !psCompile(proc, &out) {
		return nil
	}
	return out
}

func psTokens(src []byte) []string {
	var toks []string
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == '%':
			for i < len(src) && src[i] != '\n' && src[i] != '\r' {
				i++
			}
		case c == '{' || c == '}':
			toks = append(toks, string(c))
			i++
		case c <= ' ':
			i++
		default:
			j := i
			for j < len(src) && src[j] > ' ' && src[j] != '{' && src[j] != '}' && src[j] != '%' {
				j++
			}
			toks = append(toks, string(src[i:j]))
			i = j
		}
	}
	return toks
}

func psParse(toks []string, i *int, depth int) ([]psNode, bool) {
	if depth > 64 {
		return nil, false
	}
	var out []psNode
	for *i < len(toks) {
		t := toks[*i]
		*i++
		switch t {
		case "{":
			p, ok := psParse(toks, i, depth+1)
			if !ok {
				return nil, false
			}
			out = append(out, psNode{proc: p})
		case "}":
			return out, true
		default:
			if v, err := strconv.ParseFloat(t, 64); err == nil {
				out = append(out, psNode{num: v, isN: true})
			} else {
				out = append(out, psNode{op: t})
			}
		}
	}
	return nil, false
}

func psCompile(proc []psNode, out *[]psOp) bool {
	for k := 0; k < len(proc); k++ {
		n := proc[k]
		switch {
		case n.isN:
			*out = append(*out, psOp{code: psPush, num: n.num})
		case n.proc != nil:
			// Procedures are only operands of if / ifelse.
			if k+1 < len(proc) && proc[k+1].op == "if" {
				j := len(*out)
				*out = append(*out, psOp{code: psJumpFalse})
				if !psCompile(n.proc, out) {
					return false
				}
				(*out)[j].to = len(*out)
				k++
			} else if k+2 < len(proc) && proc[k+1].proc != nil && proc[k+2].op == "ifelse" {
				j := len(*out)
				*out = append(*out, psOp{code: psJumpFalse})
				if !psCompile(n.proc, out) {
					return false
				}
				e := len(*out)
				*out = append(*out, psOp{code: psJump})
				(*out)[j].to = len(*out)
				if !psCompile(proc[k+1].proc, out) {
					return false
				}
				(*out)[e].to = len(*out)
				k += 2
			} else {
				return false
			}
		default:
			c, ok := psOps[n.op]
			if !ok {
				return false
			}
			*out = append(*out, psOp{code: c})
		}
	}
	return true
}

// psVal is an operand: a number, or a boolean (b set).
type psVal struct {
	v float64
	b bool
}

func psBool(t bool) psVal {
	if t {
		return psVal{v: 1, b: true}
	}
	return psVal{b: true}
}

// runPS evaluates a compiled program on the inputs and returns the stack
// (nil on an error such as a stack underflow).
func runPS(prog []psOp, in []float64) []float64 {
	const maxStack = 100
	st := make([]psVal, 0, 16)
	for _, v := range in {
		st = append(st, psVal{v: v})
	}
	pop := func() (psVal, bool) {
		if len(st) == 0 {
			return psVal{}, false
		}
		v := st[len(st)-1]
		st = st[:len(st)-1]
		return v, true
	}
	for pc := 0; pc < len(prog); pc++ {
		op := prog[pc]
		if len(st) > maxStack {
			return nil
		}
		switch op.code {
		case psPush:
			st = append(st, psVal{v: op.num})
			continue
		case psJumpFalse:
			c, ok := pop()
			if !ok {
				return nil
			}
			if c.v == 0 {
				pc = op.to - 1
			}
			continue
		case psJump:
			pc = op.to - 1
			continue
		case psTrue:
			st = append(st, psBool(true))
			continue
		case psFalse:
			st = append(st, psBool(false))
			continue
		}
		// Stack operators.
		switch op.code {
		case psDup:
			if len(st) == 0 {
				return nil
			}
			st = append(st, st[len(st)-1])
			continue
		case psPop:
			if _, ok := pop(); !ok {
				return nil
			}
			continue
		case psExch:
			if len(st) < 2 {
				return nil
			}
			n := len(st)
			st[n-1], st[n-2] = st[n-2], st[n-1]
			continue
		case psCopy:
			c, ok := pop()
			n := int(c.v)
			if !ok || n < 0 || n > len(st) || len(st)+n > maxStack {
				return nil
			}
			st = append(st, st[len(st)-n:]...)
			continue
		case psIndex:
			c, ok := pop()
			n := int(c.v)
			if !ok || n < 0 || n >= len(st) {
				return nil
			}
			st = append(st, st[len(st)-1-n])
			continue
		case psRoll:
			jv, ok1 := pop()
			nv, ok2 := pop()
			n, j := int(nv.v), int(jv.v)
			if !ok1 || !ok2 || n < 0 || n > len(st) {
				return nil
			}
			if n > 0 {
				s := st[len(st)-n:]
				j = ((j % n) + n) % n
				r := append(append([]psVal(nil), s[n-j:]...), s[:n-j]...)
				copy(s, r)
			}
			continue
		}
		// One-operand operators.
		a, ok := pop()
		if !ok {
			return nil
		}
		var r psVal
		switch op.code {
		case psAbs:
			r.v = math.Abs(a.v)
		case psCeiling:
			r.v = math.Ceil(a.v)
		case psCos:
			r.v = math.Cos(a.v * math.Pi / 180)
		case psCvi, psTruncate:
			r.v = math.Trunc(a.v)
		case psCvr:
			r.v = a.v
		case psFloor:
			r.v = math.Floor(a.v)
		case psLn:
			r.v = math.Log(a.v)
		case psLog:
			r.v = math.Log10(a.v)
		case psNeg:
			r.v = -a.v
		case psRound:
			r.v = math.Floor(a.v + 0.5)
		case psSin:
			r.v = math.Sin(a.v * math.Pi / 180)
		case psSqrt:
			r.v = math.Sqrt(a.v)
		case psNot:
			if a.b {
				r = psBool(a.v == 0)
			} else {
				r.v = float64(^int64(a.v))
			}
		default:
			// Two-operand operators: a is the top of the stack.
			b, ok := pop()
			if !ok {
				return nil
			}
			switch op.code {
			case psAdd:
				r.v = b.v + a.v
			case psSub:
				r.v = b.v - a.v
			case psMul:
				r.v = b.v * a.v
			case psDiv:
				if a.v == 0 {
					return nil
				}
				r.v = b.v / a.v
			case psIdiv:
				if int64(a.v) == 0 {
					return nil
				}
				r.v = float64(int64(b.v) / int64(a.v))
			case psMod:
				if int64(a.v) == 0 {
					return nil
				}
				r.v = float64(int64(b.v) % int64(a.v))
			case psAtan:
				deg := math.Atan2(b.v, a.v) * 180 / math.Pi
				if deg < 0 {
					deg += 360
				}
				r.v = deg
			case psExp:
				r.v = math.Pow(b.v, a.v)
			case psEq:
				r = psBool(b.v == a.v)
			case psNe:
				r = psBool(b.v != a.v)
			case psGe:
				r = psBool(b.v >= a.v)
			case psGt:
				r = psBool(b.v > a.v)
			case psLe:
				r = psBool(b.v <= a.v)
			case psLt:
				r = psBool(b.v < a.v)
			case psAnd, psOr, psXor:
				if a.b && b.b {
					x, y := b.v != 0, a.v != 0
					switch op.code {
					case psAnd:
						r = psBool(x && y)
					case psOr:
						r = psBool(x || y)
					default:
						r = psBool(x != y)
					}
				} else {
					x, y := int64(b.v), int64(a.v)
					switch op.code {
					case psAnd:
						r.v = float64(x & y)
					case psOr:
						r.v = float64(x | y)
					default:
						r.v = float64(x ^ y)
					}
				}
			case psBitshift:
				x, s := int64(b.v), int64(a.v)
				if s >= 0 {
					r.v = float64(x << uint(min(s, 63)))
				} else {
					r.v = float64(x >> uint(min(-s, 63)))
				}
			default:
				return nil
			}
		}
		if math.IsNaN(r.v) || math.IsInf(r.v, 0) {
			r.v = 0
		}
		st = append(st, r)
	}
	out := make([]float64, len(st))
	for i, v := range st {
		out[i] = v.v
	}
	return out
}

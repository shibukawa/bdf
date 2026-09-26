package main

import (
	"fmt"
	"strconv"
	"strings"
)

// The sidebar scripts build their palette entries from string and number
// variables: a style prefix, a scale, a package name. Instead of running
// them, the generator reads them with a small evaluator for the subset of
// JavaScript they are written in: var declarations of string and number
// expressions (literals, variables, mxConstants, + - * /), and the calls
// that create entries (createVertexTemplateEntry, and new mxCell inside
// addEntry functions, titled by createVertexTemplateFromCells). Variables
// are one global scope in which the latest declaration wins, which is
// what each palette function sees, as they declare their locals before
// using them.

// paletteEntry is an entry of a sidebar palette.
type paletteEntry struct {
	palette string // the palette's id (addPaletteFunctions), e.g. "aws3Compute"
	style   string
	w, h    float64
	title   string
	tags    string // the stencil's own tags (getTagsForStencil's second argument)
	group   bool   // a cell of a composite entry (addEntry)
}

// jsToken is a token of a script: kind 's' string, 'n' number, 'i'
// identifier, 'p' punctuation.
type jsToken struct {
	kind byte
	text string
}

// tokenize splits a script into tokens, dropping comments and whitespace.
// Regular expression literals are not recognized (the sidebar scripts
// have none).
func tokenize(src string) ([]jsToken, error) {
	var toks []jsToken
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			i++
		case strings.HasPrefix(src[i:], "//"):
			for i < len(src) && src[i] != '\n' {
				i++
			}
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i+2:], "*/")
			if end < 0 {
				return nil, fmt.Errorf("unterminated comment")
			}
			i += end + 4
		case c == '\'' || c == '"':
			var b strings.Builder
			j := i + 1
			for ; j < len(src) && src[j] != c; j++ {
				if src[j] == '\\' && j+1 < len(src) {
					j++
					switch src[j] {
					case 'n':
						b.WriteByte('\n')
					case 't':
						b.WriteByte('\t')
					default:
						b.WriteByte(src[j])
					}
					continue
				}
				b.WriteByte(src[j])
			}
			if j >= len(src) {
				return nil, fmt.Errorf("unterminated string")
			}
			toks = append(toks, jsToken{'s', b.String()})
			i = j + 1
		case c >= '0' && c <= '9' || c == '.' && i+1 < len(src) && src[i+1] >= '0' && src[i+1] <= '9':
			j := i
			for j < len(src) && (src[j] >= '0' && src[j] <= '9' || src[j] == '.') {
				j++
			}
			toks = append(toks, jsToken{'n', src[i:j]})
			i = j
		case c == '_' || c == '$' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			j := i
			for j < len(src) && (src[j] == '_' || src[j] == '$' || src[j] >= 'a' && src[j] <= 'z' ||
				src[j] >= 'A' && src[j] <= 'Z' || src[j] >= '0' && src[j] <= '9') {
				j++
			}
			toks = append(toks, jsToken{'i', src[i:j]})
			i = j
		default:
			toks = append(toks, jsToken{'p', string(c)})
			i++
		}
	}
	return toks, nil
}

// mxConstants are the style constants the scripts use.
var mxConstants = map[string]string{
	"STYLE_SHAPE":                   "shape",
	"STYLE_STROKEWIDTH":             "strokeWidth",
	"STYLE_VERTICAL_ALIGN":          "verticalAlign",
	"STYLE_VERTICAL_LABEL_POSITION": "verticalLabelPosition",
	"STYLE_FILLCOLOR":               "fillColor",
	"STYLE_STROKECOLOR":             "strokeColor",
}

// jsReader walks the tokens of a script and collects its palette entries.
type jsReader struct {
	toks    []jsToken
	pos     int
	env     map[string]any // string or float64
	palette string
	entries []paletteEntry
	// cells and entryTags are the cells (indexes in entries) and the
	// tags of the composite entry (addEntry) being read
	cells     []int
	entryTags string
}

// readPalettes returns the palette entries of a sidebar script.
func readPalettes(src string) ([]paletteEntry, error) {
	toks, err := tokenize(src)
	if err != nil {
		return nil, err
	}
	r := &jsReader{toks: toks, env: map[string]any{}}
	for r.pos < len(r.toks) {
		t := r.toks[r.pos]
		switch {
		case t.kind == 'i' && t.text == "var" && r.peek(1).kind == 'i' && r.peek(2).text == "=":
			name := r.peek(1).text
			r.pos += 3
			if v, ok := r.expr(); ok {
				r.env[name] = v
			} else {
				delete(r.env, name)
			}
		case t.kind == 'i' && r.peek(1).text == "(" && r.prevDot():
			r.call(t.text)
		case t.kind == 'i' && t.text == "new" && r.peek(1).text == "mxCell" && r.peek(2).text == "(":
			r.pos += 3
			args := r.args()
			if len(args) == 3 {
				if style, ok := args[2].(string); ok {
					e := paletteEntry{palette: r.palette, style: style, group: true}
					if g, ok := args[1].([]float64); ok && len(g) == 4 {
						e.w, e.h = g[2], g[3]
					}
					r.cells = append(r.cells, len(r.entries))
					r.entries = append(r.entries, e)
				}
			}
		default:
			r.pos++
		}
	}
	return r.entries, nil
}

func (r *jsReader) peek(n int) jsToken {
	if r.pos+n < len(r.toks) {
		return r.toks[r.pos+n]
	}
	return jsToken{}
}

// prevDot reports whether the token before the current one is a ".".
func (r *jsReader) prevDot() bool {
	return r.pos > 0 && r.toks[r.pos-1].text == "."
}

// call reads a method call at the current token (its name) that makes
// entries or names the palette; other calls are stepped over.
func (r *jsReader) call(name string) {
	switch name {
	case "createVertexTemplateEntry":
		r.pos += 2
		args := r.args()
		e := paletteEntry{palette: r.palette}
		e.style, _ = arg(args, 0).(string)
		e.w, _ = arg(args, 1).(float64)
		e.h, _ = arg(args, 2).(float64)
		e.title, _ = arg(args, 4).(string)
		if e.title == "" {
			e.title, _ = arg(args, 3).(string) // the value when the title is null
		}
		e.tags, _ = arg(args, 7).(string)
		if e.style != "" {
			r.entries = append(r.entries, e)
		}
	case "addPaletteFunctions":
		// only the id: the entries follow in the last argument
		r.pos += 2
		if v, ok := r.expr(); ok {
			r.palette, _ = v.(string)
		}
	case "addEntry":
		r.pos += 2
		r.cells, r.entryTags = nil, ""
		if v, ok := r.expr(); ok {
			r.entryTags, _ = v.(string)
		}
	case "createVertexTemplateFromCells":
		r.pos += 2
		args := r.args()
		title, _ := arg(args, 3).(string)
		for _, i := range r.cells {
			r.entries[i].title = title
			r.entries[i].tags = r.entryTags
		}
		r.cells = nil
	default:
		r.pos++
	}
}

func arg(args []any, i int) any {
	if i < len(args) {
		return args[i]
	}
	return nil
}

// args reads the arguments of a call after its "(" up to the matching
// ")"; arguments that cannot be evaluated are nil.
func (r *jsReader) args() []any {
	var out []any
	if r.peek(0).text == ")" {
		r.pos++
		return nil
	}
	for r.pos < len(r.toks) {
		start := r.pos
		v, ok := r.expr()
		if !ok {
			v = nil
			r.pos = start
			r.skipArg()
		}
		out = append(out, v)
		switch r.peek(0).text {
		case ",":
			r.pos++
		case ")":
			r.pos++
			return out
		default:
			r.skipArg()
			if r.peek(0).text == "," {
				r.pos++
				continue
			}
			r.pos++
			return out
		}
	}
	return out
}

// skipArg steps to the "," or ")" that ends the current argument.
func (r *jsReader) skipArg() {
	depth := 0
	for r.pos < len(r.toks) {
		switch r.peek(0).text {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			if depth == 0 {
				return
			}
			depth--
		case ",":
			if depth == 0 {
				return
			}
		}
		r.pos++
	}
}

// expr evaluates an additive expression.
func (r *jsReader) expr() (any, bool) {
	v, ok := r.term()
	for ok {
		op := r.peek(0)
		if op.kind != 'p' || op.text != "+" && op.text != "-" {
			break
		}
		r.pos++
		var w any
		w, ok = r.term()
		if !ok {
			break
		}
		if op.text == "+" {
			vs, vIsStr := v.(string)
			ws, wIsStr := w.(string)
			if vIsStr || wIsStr {
				if !vIsStr {
					vs = jsString(v)
				}
				if !wIsStr {
					ws = jsString(w)
				}
				v = vs + ws
				continue
			}
		}
		a, aok := v.(float64)
		b, bok := w.(float64)
		if !aok || !bok {
			return nil, false
		}
		if op.text == "+" {
			v = a + b
		} else {
			v = a - b
		}
	}
	return v, ok
}

// term evaluates a multiplicative expression.
func (r *jsReader) term() (any, bool) {
	v, ok := r.unary()
	for ok {
		op := r.peek(0)
		if op.kind != 'p' || op.text != "*" && op.text != "/" {
			break
		}
		r.pos++
		var w any
		w, ok = r.unary()
		a, aok := v.(float64)
		b, bok := w.(float64)
		if !ok || !aok || !bok {
			return nil, false
		}
		if op.text == "*" {
			v = a * b
		} else {
			v = a / b
		}
	}
	return v, ok
}

func (r *jsReader) unary() (any, bool) {
	if t := r.peek(0); t.kind == 'p' && t.text == "-" {
		r.pos++
		v, ok := r.unary()
		f, fok := v.(float64)
		return -f, ok && fok
	}
	return r.primary()
}

// primary evaluates literals, variables, mxConstants, parentheses,
// new mxGeometry(...) and this.getTagsForStencil(...).join(...).
func (r *jsReader) primary() (any, bool) {
	t := r.peek(0)
	switch {
	case t.kind == 's':
		r.pos++
		return t.text, true
	case t.kind == 'n':
		r.pos++
		f, err := strconv.ParseFloat(t.text, 64)
		return f, err == nil
	case t.text == "(":
		r.pos++
		v, ok := r.expr()
		if !ok || r.peek(0).text != ")" {
			return nil, false
		}
		r.pos++
		return v, true
	case t.text == "new" && r.peek(1).text == "mxGeometry" && r.peek(2).text == "(":
		r.pos += 3
		args := r.args()
		g := make([]float64, 4)
		for i := range g {
			g[i], _ = arg(args, i).(float64)
		}
		return g, true
	case t.text == "mxConstants" && r.peek(1).text == "." && r.peek(2).kind == 'i':
		v, ok := mxConstants[r.peek(2).text]
		r.pos += 3
		return v, ok
	case t.text == "this" && r.peek(1).text == "." && r.peek(2).text == "getTagsForStencil" && r.peek(3).text == "(":
		r.pos += 4
		args := r.args()
		tags, _ := arg(args, 1).(string)
		// .join(' ')
		if r.peek(0).text == "." && r.peek(1).text == "join" && r.peek(2).text == "(" {
			r.pos += 3
			r.args()
		}
		return strings.ReplaceAll(tags, "_", " "), true
	case t.kind == 'i':
		switch t.text {
		case "null", "undefined":
			r.pos++
			return nil, true
		case "true", "false":
			r.pos++
			return t.text == "true", true
		}
		v, ok := r.env[t.text]
		if !ok || r.peek(1).text == "." || r.peek(1).text == "(" {
			return nil, false
		}
		r.pos++
		return v, true
	}
	return nil, false
}

// jsString converts a value to a string as JavaScript's + does.
func jsString(v any) string {
	switch v := v.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "null"
	}
	return fmt.Sprint(v)
}

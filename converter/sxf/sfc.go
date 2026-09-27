package sxf

import (
	"fmt"
	"strconv"
	"strings"
)

// feature is one feature of an SFC file: its name and its arguments, the
// strings unquoted and the numbers and lists as they are written.
type feature struct {
	name string
	args []arg
}

type arg struct {
	s     string
	isStr bool // a string (\'...\'), not a number or a list ('...')
}

func (f *feature) str(i int) string {
	if i < len(f.args) {
		return f.args[i].s
	}
	return ""
}

func (f *feature) num(i int) float64 {
	if i >= len(f.args) {
		return 0
	}
	v, err := strconv.ParseFloat(strings.TrimSpace(f.args[i].s), 64)
	if err != nil {
		return 0
	}
	return v
}

func (f *feature) int(i int) int { return int(f.num(i)) }

// list reads a list argument: "(1.0,2.0,3.0)".
func (f *feature) list(i int) []float64 {
	if i >= len(f.args) {
		return nil
	}
	return parseList(f.args[i].s)
}

func parseList(s string) []float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(strings.TrimPrefix(s, "("), ")")
	var out []float64
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			v = 0
		}
		out = append(out, v)
	}
	return out
}

// readSFC splits an SFC file (decoded) into its features in order: the
// "#n = name(args)" inside each /*SXF ... SXF*/ comment (and the /*SXF3
// and /*SXF3.1 ones of later features).
func readSFC(text string) ([]*feature, error) {
	var out []*feature
	rest := text
	for {
		i := strings.Index(rest, "/*SXF")
		if i < 0 {
			break
		}
		rest = rest[i+len("/*SXF"):]
		// the version after the keyword ("3", "3.1")
		for len(rest) > 0 && (rest[0] >= '0' && rest[0] <= '9' || rest[0] == '.') {
			rest = rest[1:]
		}
		j := strings.Index(rest, "*/")
		if j < 0 {
			return out, fmt.Errorf("a feature comment is not closed")
		}
		body := rest[:j]
		rest = rest[j+2:]
		// the closing keyword ("SXF", "SXF3") before */
		if k := strings.LastIndex(body, "SXF"); k >= 0 {
			body = body[:k]
		}
		if f := parseFeature(body); f != nil {
			out = append(out, f)
		}
	}
	return out, nil
}

// parseFeature reads "#n = name(args)".
func parseFeature(s string) *feature {
	eq := strings.IndexByte(s, '=')
	open := strings.IndexByte(s, '(')
	if eq < 0 || open < eq {
		return nil
	}
	f := &feature{name: strings.ToLower(strings.TrimSpace(s[eq+1 : open]))}
	body := s[open+1:]
	for i := 0; i < len(body); {
		c := body[i]
		switch {
		case c == ' ' || c == '\t' || c == '\r' || c == '\n' || c == ',':
			i++
		case c == ')':
			return f
		case c == '\\' && i+1 < len(body) && body[i+1] == '\'':
			// a string, up to the next \' (a doubled backslash is one)
			i += 2
			var b strings.Builder
			for i < len(body) {
				if body[i] == '\\' && i+1 < len(body) && body[i+1] == '\'' {
					i += 2
					break
				}
				if body[i] == '\\' && i+1 < len(body) && body[i+1] == '\\' {
					b.WriteByte('\\')
					i += 2
					continue
				}
				b.WriteByte(body[i])
				i++
			}
			f.args = append(f.args, arg{s: b.String(), isStr: true})
		case c == '\'':
			j := strings.IndexByte(body[i+1:], '\'')
			if j < 0 {
				f.args = append(f.args, arg{s: body[i+1:]})
				return f
			}
			f.args = append(f.args, arg{s: body[i+1 : i+1+j]})
			i += j + 2
		default:
			// an unquoted value, up to the next comma or parenthesis
			j := i
			depth := 0
			for j < len(body) {
				if body[j] == '(' {
					depth++
				} else if body[j] == ')' {
					if depth == 0 {
						break
					}
					depth--
				} else if body[j] == ',' && depth == 0 {
					break
				}
				j++
			}
			f.args = append(f.args, arg{s: strings.TrimSpace(body[i:j])})
			i = j
		}
	}
	return f
}

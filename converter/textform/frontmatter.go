package textform

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// grid is the paper of a template: what its front matter says.
type grid struct {
	cpi, lpi       float64
	columns, lines int
	left, top      int
	skip, title    int
	font           string
	ambiguousWide  bool
}

// defaultGrid is the paper of a template without front matter: fanfold
// paper of 8 by 11 inches at 10 CPI and 6 LPI.
func defaultGrid() grid {
	return grid{cpi: 10, lpi: 6, columns: 80, lines: 66}
}

// rows is the lines a page prints (0: as many as it has).
func (g *grid) rows() int {
	if g.lines == 0 {
		return 0
	}
	return g.lines - g.top - g.skip
}

// width is the columns of a line after the left margin.
func (g *grid) width() int { return g.columns - g.left }

// splitFrontMatter cuts the front matter off a template: the lines between
// a first line "---" and the next line "---" (or "..."). body is what
// follows, and first the number of the line of the file it starts on.
func splitFrontMatter(s string) (front []string, body string, first int, has bool) {
	if !strings.HasPrefix(s, "---\n") && s != "---" {
		return nil, s, 1, false
	}
	lines := strings.SplitAfter(s, "\n")
	for i := 1; i < len(lines); i++ {
		l := strings.TrimRight(lines[i], " \t\n")
		if l == "---" || l == "..." {
			for _, fl := range lines[1:i] {
				front = append(front, strings.TrimSuffix(fl, "\n"))
			}
			return front, strings.Join(lines[i+1:], ""), i + 2, true
		}
	}
	// (a first line "---" without its end: a rule of the body)
	return nil, s, 1, false
}

// parseFrontMatter reads the front matter: the keys of the grid, in the
// part of YAML that they need (keys and plain or quoted scalars, the
// margin as a flow mapping {left: 2, top: 1} or as an indented block, and
// comments). Unknown keys are told in warnings.
func parseFrontMatter(lines []string, warn func(format string, args ...any)) (grid, error) {
	g := defaultGrid()
	set := map[string]string{}
	var order []string
	block := "" // the key of an indented block being read
	for i, raw := range lines {
		line := strings.TrimRight(stripComment(raw), " \t")
		if strings.TrimSpace(line) == "" {
			continue
		}
		at := func(format string, args ...any) error {
			return fmt.Errorf("text: the front matter, line %d: %s", i+2, fmt.Sprintf(format, args...))
		}
		indented := line[0] == ' ' || line[0] == '\t'
		key, value, ok := strings.Cut(strings.TrimSpace(line), ":")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" {
			return g, at("%q is not a key and its value (key: value)", strings.TrimSpace(line))
		}
		if indented {
			if block == "" {
				return g, at("%q is indented under no key", key)
			}
			key = block + "." + key
		} else {
			block = ""
			if value == "" {
				block = key
				continue
			}
		}
		if strings.HasPrefix(value, "{") {
			if !strings.HasSuffix(value, "}") {
				return g, at("the mapping of %s is not closed with }", key)
			}
			for _, kv := range strings.Split(value[1:len(value)-1], ",") {
				if strings.TrimSpace(kv) == "" {
					continue
				}
				k, v, ok := strings.Cut(kv, ":")
				if !ok {
					return g, at("%q in the mapping of %s is not a key and its value", strings.TrimSpace(kv), key)
				}
				k = key + "." + strings.TrimSpace(k)
				if _, dup := set[k]; !dup {
					order = append(order, k)
				}
				set[k] = unquote(strings.TrimSpace(v))
			}
			continue
		}
		if _, dup := set[key]; !dup {
			order = append(order, key)
		}
		set[key] = unquote(value)
	}
	num := func(key, v string, lo, hi float64) (float64, error) {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || math.IsNaN(f) || f < lo || f > hi {
			return 0, fmt.Errorf("text: the front matter: %s is %q; it takes a number from %g to %g", key, v, lo, hi)
		}
		return f, nil
	}
	whole := func(key, v string, lo, hi int) (int, error) {
		n, err := strconv.Atoi(v)
		if err != nil || n < lo || n > hi {
			return 0, fmt.Errorf("text: the front matter: %s is %q; it takes a whole number from %d to %d", key, v, lo, hi)
		}
		return n, nil
	}
	var err error
	for _, k := range order {
		v := set[k]
		switch k {
		case "cpi":
			g.cpi, err = num(k, v, 1, 30)
			if err == nil && !pitches[g.cpi] {
				warn("the front matter: %g CPI is not a pitch of ESC/P (10, 12, 15, 17.1 or 20); the text mode of a printer prints at the nearest one", g.cpi)
			}
		case "lpi":
			g.lpi, err = num(k, v, 1, 30)
		case "columns":
			g.columns, err = whole(k, v, 1, maxColumns)
		case "lines":
			g.lines, err = whole(k, v, 0, maxLines)
		case "margin.left":
			g.left, err = whole(k, v, 0, maxColumns)
		case "margin.top":
			g.top, err = whole(k, v, 0, maxLines)
		case "skip":
			g.skip, err = whole(k, v, 0, maxLines)
		case "title":
			g.title, err = whole(k, v, 0, maxLines)
		case "font":
			g.font = v
		case "ambiguous":
			switch strings.ToLower(v) {
			case "narrow", "1":
				g.ambiguousWide = false
			case "wide", "2":
				g.ambiguousWide = true
			default:
				err = fmt.Errorf("text: the front matter: ambiguous is %q; it is narrow or wide", v)
			}
		default:
			warn("the front matter has %s, which a template of text does not read", k)
		}
		if err != nil {
			return g, err
		}
	}
	switch {
	case g.left >= g.columns:
		return g, fmt.Errorf("text: the front matter: the left margin of %d columns leaves no column of the %d of the paper", g.left, g.columns)
	case g.lines > 0 && g.top+g.skip >= g.lines:
		return g, fmt.Errorf("text: the front matter: the top margin and the skip (%d and %d lines) leave no line of the %d of the paper", g.top, g.skip, g.lines)
	case g.lines > 0 && g.title >= g.rows():
		return g, fmt.Errorf("text: the front matter: the %d lines of the title leave no line of the %d of a page", g.title, g.rows())
	}
	return g, nil
}

// pitches are the pitches of ESC/P.
var pitches = map[float64]bool{10: true, 12: true, 15: true, 17: true, 17.1: true, 17.14: true, 20: true}

// stripComment drops a comment: from a # that starts the line or follows a
// space, outside quotes.
func stripComment(s string) string {
	var quote byte
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '"' || c == '\'':
			quote = c
		case c == '#' && (i == 0 || s[i-1] == ' ' || s[i-1] == '\t'):
			return s[:i]
		}
	}
	return s
}

// unquote reads a scalar: quoted or plain.
func unquote(s string) string {
	if len(s) >= 2 && (s[0] == '"' && s[len(s)-1] == '"' || s[0] == '\'' && s[len(s)-1] == '\'') {
		if s[0] == '"' {
			if u, err := strconv.Unquote(s); err == nil {
				return u
			}
		}
		return strings.ReplaceAll(s[1:len(s)-1], "''", "'")
	}
	return s
}

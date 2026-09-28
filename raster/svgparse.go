package raster

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

// svgNode is an element of an SVG document, or a run of text ("#text").
type svgNode struct {
	name     string
	attr     map[string]string
	children []*svgNode
	parent   *svgNode
	sheet    map[string]string // declarations of the style sheets that match
	text     string

	// what drawing the element found, kept for the next time it is drawn
	// (a use element draws it again): the path of its data, and the
	// picture or the SVG image of its data: URL
	path  *path
	pic   *picture
	image *svgImage
	tried bool // the data: URL was read, with an image or without
}

// svgDoc is a parsed SVG document.
type svgDoc struct {
	root *svgNode
	ids  map[string]*svgNode
	// warnings say what was left out of a document past the limits
	warnings []string
}

// Limits of an SVG document. Its elements are kept in memory, each many
// times larger than its markup, and the functions that go through them
// call themselves for the elements inside.
const (
	// svgDepth is how deep elements may be in elements (the limit of
	// libxml2, the parser of browsers, is 256 too).
	svgDepth = 256
	// svgNodes is the most elements and runs of text.
	svgNodes = 1 << 20
	// svgMatches is the most elements that the rules of the style sheets
	// are matched against, one rule and one element at a time.
	svgMatches = 1 << 22
)

// parseSVG reads an SVG document; it is lenient, as browsers are with the
// SVG images of web pages.
func parseSVG(data []byte) (*svgDoc, bool) {
	return parseSVGWithin(data, svgDepth, svgNodes, svgMatches)
}

// parseSVGWithin is parseSVG with its limits: the depth of the elements,
// their number and the matches of the rules of the style sheets.
func parseSVGWithin(data []byte, depth, most, matches int) (*svgDoc, bool) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	dec.Strict = false
	dec.AutoClose = xml.HTMLAutoClose
	dec.Entity = xml.HTMLEntity
	dec.CharsetReader = charset.NewReaderLabel
	doc := &svgDoc{ids: map[string]*svgNode{}}
	var stack []*svgNode
	var styles []string
	nodes := 0
read:
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			if doc.root == nil {
				return nil, false
			}
			break // keep what was read, as browsers draw up to an error
		}
		switch t := tok.(type) {
		case xml.StartElement:
			// what passes a limit ends the reading as an error does
			if len(stack) >= depth {
				doc.warnings = append(doc.warnings, fmt.Sprintf("the elements of an SVG image are more than %d deep: the rest are not read", depth))
				break read
			}
			if nodes++; nodes > most {
				doc.warnings = append(doc.warnings, fmt.Sprintf("an SVG image has more than %d elements: the rest are not read", most))
				break read
			}
			n := &svgNode{name: t.Name.Local, attr: make(map[string]string, len(t.Attr))}
			for _, a := range t.Attr {
				if a.Name.Space == "xmlns" || a.Name.Local == "xmlns" {
					continue
				}
				if _, ok := n.attr[a.Name.Local]; ok && a.Name.Local == "href" && a.Name.Space != "" {
					continue // href wins over xlink:href
				}
				n.attr[a.Name.Local] = a.Value
			}
			if len(stack) > 0 {
				p := stack[len(stack)-1]
				n.parent = p
				p.children = append(p.children, n)
			} else if doc.root == nil {
				doc.root = n
			}
			if id := n.attr["id"]; id != "" {
				if _, dup := doc.ids[id]; !dup {
					doc.ids[id] = n
				}
			}
			stack = append(stack, n)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			p := stack[len(stack)-1]
			switch p.name {
			case "style":
				styles = append(styles, string(t))
			case "text", "tspan", "textPath", "a":
				if nodes++; nodes > most {
					continue
				}
				p.children = append(p.children, &svgNode{name: "#text", text: string(t), parent: p})
			}
		}
	}
	if doc.root == nil || doc.root.name != "svg" {
		return nil, false
	}
	if len(styles) > 0 {
		sheet := newStyleSheet(parseCSS(strings.Join(styles, "\n")), matches)
		sheet.apply(doc.root)
		if sheet.matches > sheet.most {
			doc.warnings = append(doc.warnings, "the style sheets of an SVG image have too many rules for its elements: the rest of the elements are drawn without them")
		}
	}
	return doc, true
}

// cssRule is a rule of a style sheet with one selector.
type cssRule struct {
	sel   []cssCompound // descendant combinators between them
	spec  int
	order int
	decls [][2]string
}

// cssCompound is a simple selector sequence: tag, #id and .classes.
type cssCompound struct {
	tag, id string
	classes []string
}

// parseCSS reads the rules of a style sheet that previews need: compound
// selectors of type, id and class, and descendant combinators. At-rules
// are skipped.
func parseCSS(src string) []cssRule {
	if strings.Contains(src, "/*") {
		// without the comments, in one pass
		var b strings.Builder
		for {
			i := strings.Index(src, "/*")
			if i < 0 {
				b.WriteString(src)
				break
			}
			b.WriteString(src[:i])
			j := strings.Index(src[i+2:], "*/")
			if j < 0 {
				break
			}
			src = src[i+2+j+2:]
		}
		src = b.String()
	}
	var rules []cssRule
	order := 0
	for {
		open := strings.IndexByte(src, '{')
		if open < 0 {
			break
		}
		prelude := strings.TrimSpace(src[:open])
		// find the matching close brace
		depth, end := 0, -1
		for k := open; k < len(src); k++ {
			switch src[k] {
			case '{':
				depth++
			case '}':
				depth--
			}
			if depth == 0 {
				end = k
				break
			}
		}
		if end < 0 {
			break
		}
		body := src[open+1 : end]
		src = src[end+1:]
		if strings.HasPrefix(prelude, "@") {
			continue
		}
		decls := parseDecls(body)
		for _, s := range strings.Split(prelude, ",") {
			sel, spec, ok := parseSelector(strings.TrimSpace(s))
			if !ok {
				continue
			}
			rules = append(rules, cssRule{sel: sel, spec: spec, order: order, decls: decls})
			order++
		}
	}
	return rules
}

func parseSelector(s string) ([]cssCompound, int, bool) {
	if s == "" {
		return nil, 0, false
	}
	var out []cssCompound
	spec := 0
	for _, part := range strings.Fields(strings.ReplaceAll(s, ">", " ")) {
		var c cssCompound
		i := 0
		read := func() string {
			j := i
			for j < len(part) && part[j] != '.' && part[j] != '#' && part[j] != '[' && part[j] != ':' {
				j++
			}
			v := part[i:j]
			i = j
			return v
		}
		if part[0] != '.' && part[0] != '#' {
			c.tag = read()
			if c.tag != "*" {
				spec++
			} else {
				c.tag = ""
			}
		}
		for i < len(part) {
			switch part[i] {
			case '.':
				i++
				c.classes = append(c.classes, read())
				spec += 100
			case '#':
				i++
				c.id = read()
				spec += 10000
			default:
				return nil, 0, false // attribute selectors and pseudo-classes
			}
		}
		out = append(out, c)
	}
	return out, spec, true
}

func (c *cssCompound) matches(n *svgNode) bool {
	if c.tag != "" && c.tag != n.name || c.id != "" && n.attr["id"] != c.id {
		return false
	}
	if len(c.classes) > 0 {
		have := strings.Fields(n.attr["class"])
		for _, want := range c.classes {
			found := false
			for _, h := range have {
				if h == want {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

func (r *cssRule) matches(n *svgNode) bool {
	last := len(r.sel) - 1
	if !r.sel[last].matches(n) {
		return false
	}
	k := last - 1
	for p := n.parent; p != nil && k >= 0; p = p.parent {
		if r.sel[k].matches(p) {
			k--
		}
	}
	return k < 0
}

// styleSheet is the rules of the style sheets of a document, by what the
// last part of their selector asks for: an element is matched against the
// rules that ask for its id, one of its classes or its name, and those that
// ask for none.
type styleSheet struct {
	byID, byClass, byTag map[string][]*cssRule
	others               []*cssRule
	// matches counts the rules matched against elements, and most is the
	// most there may be (svgMatches)
	matches, most int
}

func newStyleSheet(rules []cssRule, most int) *styleSheet {
	s := &styleSheet{byID: map[string][]*cssRule{}, byClass: map[string][]*cssRule{}, byTag: map[string][]*cssRule{}, most: most}
	for i := range rules {
		r := &rules[i]
		switch last := r.sel[len(r.sel)-1]; {
		case last.id != "":
			s.byID[last.id] = append(s.byID[last.id], r)
		case len(last.classes) > 0:
			s.byClass[last.classes[0]] = append(s.byClass[last.classes[0]], r)
		case last.tag != "":
			s.byTag[last.tag] = append(s.byTag[last.tag], r)
		default:
			s.others = append(s.others, r)
		}
	}
	return s
}

// apply stores on each element the declarations of the rules that match
// it, the more specific and later ones winning.
func (s *styleSheet) apply(n *svgNode) {
	var hit []*cssRule
	match := func(rules []*cssRule) {
		for _, r := range rules {
			if s.matches++; s.matches > s.most {
				return
			}
			if r.matches(n) {
				hit = append(hit, r)
			}
		}
	}
	if id := n.attr["id"]; id != "" {
		match(s.byID[id])
	}
	if class := n.attr["class"]; class != "" {
		classes := strings.Fields(class)
		for i, c := range classes {
			if !slices.Contains(classes[:i], c) {
				match(s.byClass[c])
			}
		}
	}
	match(s.byTag[n.name])
	match(s.others)
	if len(hit) > 0 {
		sort.SliceStable(hit, func(i, j int) bool {
			if hit[i].spec != hit[j].spec {
				return hit[i].spec < hit[j].spec
			}
			return hit[i].order < hit[j].order
		})
		n.sheet = map[string]string{}
		for _, r := range hit {
			for _, d := range r.decls {
				n.sheet[d[0]] = d[1]
			}
		}
	}
	for _, c := range n.children {
		if c.name != "#text" && s.matches <= s.most {
			s.apply(c)
		}
	}
}

// parseDecls reads CSS declarations ("fill: red; stroke: blue").
func parseDecls(s string) [][2]string {
	var out [][2]string
	for _, d := range strings.Split(s, ";") {
		k, v, ok := strings.Cut(d, ":")
		if !ok {
			continue
		}
		v = strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(v), "!important"))
		out = append(out, [2]string{strings.ToLower(strings.TrimSpace(k)), strings.TrimSpace(v)})
	}
	return out
}

// svgColor is a parsed colour: RGBA from 0 to 1 (not premultiplied).
type svgColor struct {
	r, g, b, a float64
}

// parseColor reads a CSS colour; current is currentColor.
func parseColor(s string, current svgColor) (svgColor, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch {
	case s == "currentcolor":
		return current, true
	case s == "transparent":
		return svgColor{}, true
	case strings.HasPrefix(s, "#"):
		h := s[1:]
		nib := func(i int) float64 {
			v, _ := strconv.ParseUint(h[i:i+1], 16, 8)
			return float64(v*17) / 255
		}
		byt := func(i int) float64 {
			v, _ := strconv.ParseUint(h[i:i+2], 16, 8)
			return float64(v) / 255
		}
		if _, err := strconv.ParseUint(h, 16, 64); err != nil {
			return svgColor{}, false
		}
		switch len(h) {
		case 3:
			return svgColor{nib(0), nib(1), nib(2), 1}, true
		case 4:
			return svgColor{nib(0), nib(1), nib(2), nib(3)}, true
		case 6:
			return svgColor{byt(0), byt(2), byt(4), 1}, true
		case 8:
			return svgColor{byt(0), byt(2), byt(4), byt(6)}, true
		}
		return svgColor{}, false
	case strings.HasPrefix(s, "rgb"), strings.HasPrefix(s, "hsl"):
		open, close := strings.IndexByte(s, '('), strings.LastIndexByte(s, ')')
		if open < 0 || close < open {
			return svgColor{}, false
		}
		parts := strings.FieldsFunc(s[open+1:close], func(r rune) bool { return r == ',' || r == ' ' || r == '/' })
		if len(parts) < 3 {
			return svgColor{}, false
		}
		num := func(p string, scale float64) float64 {
			if strings.HasSuffix(p, "%") {
				v, _ := strconv.ParseFloat(strings.TrimSuffix(p, "%"), 64)
				return v / 100
			}
			v, _ := strconv.ParseFloat(strings.TrimSuffix(p, "deg"), 64)
			return v / scale
		}
		a := 1.0
		if len(parts) >= 4 {
			a = num(parts[3], 1)
		}
		clamp := func(v float64) float64 { return math.Min(1, math.Max(0, v)) }
		if strings.HasPrefix(s, "rgb") {
			return svgColor{clamp(num(parts[0], 255)), clamp(num(parts[1], 255)), clamp(num(parts[2], 255)), clamp(a)}, true
		}
		hh := math.Mod(num(parts[0], 1), 360) / 360
		if hh < 0 {
			hh++
		}
		sat, l := clamp(num(parts[1], 100)), clamp(num(parts[2], 100))
		r, g, b := hslToRGB(hh, sat, l)
		return svgColor{r, g, b, clamp(a)}, true
	}
	if c, ok := namedColors[s]; ok {
		return svgColor{float64(c>>16) / 255, float64(c>>8&0xff) / 255, float64(c&0xff) / 255, 1}, true
	}
	return svgColor{}, false
}

func hslToRGB(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		return l, l, l
	}
	q := l * (1 + s)
	if l >= 0.5 {
		q = l + s - l*s
	}
	p := 2*l - q
	hue := func(t float64) float64 {
		if t < 0 {
			t++
		}
		if t > 1 {
			t--
		}
		switch {
		case t < 1.0/6:
			return p + (q-p)*6*t
		case t < 0.5:
			return q
		case t < 2.0/3:
			return p + (q-p)*(2.0/3-t)*6
		}
		return p
	}
	return hue(h + 1.0/3), hue(h), hue(h - 1.0/3)
}

// svgNumber reads a number at the start of s: its value and the rest.
func svgNumber(s string) (float64, string, bool) {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == ',' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	s = s[i:]
	j := 0
	if j < len(s) && (s[j] == '+' || s[j] == '-') {
		j++
	}
	digits, dot := false, false
	for j < len(s) {
		c := s[j]
		if c >= '0' && c <= '9' {
			digits = true
		} else if c == '.' && !dot {
			dot = true
		} else {
			break
		}
		j++
	}
	if !digits {
		return 0, s, false
	}
	if j < len(s) && (s[j] == 'e' || s[j] == 'E') {
		k := j + 1
		if k < len(s) && (s[k] == '+' || s[k] == '-') {
			k++
		}
		if k < len(s) && s[k] >= '0' && s[k] <= '9' {
			for k < len(s) && s[k] >= '0' && s[k] <= '9' {
				k++
			}
			j = k
		}
	}
	v, err := strconv.ParseFloat(s[:j], 64)
	return v, s[j:], err == nil
}

// svgNumbers reads a list of numbers.
func svgNumbers(s string) []float64 {
	var out []float64
	for {
		v, rest, ok := svgNumber(s)
		if !ok {
			return out
		}
		out = append(out, v)
		s = rest
	}
}

// unit sizes in px
var svgUnits = map[string]float64{"": 1, "px": 1, "pt": 96.0 / 72, "pc": 16, "in": 96, "cm": 96 / 2.54, "mm": 96 / 25.4, "q": 96 / 101.6}

// svgLength reads a length in px; ref resolves percentages and em is the
// font size.
func svgLength(s string, ref, em float64) (float64, bool) {
	v, rest, ok := svgNumber(s)
	if !ok {
		return 0, false
	}
	u := strings.ToLower(strings.TrimSpace(rest))
	switch u {
	case "%":
		return v * ref / 100, true
	case "em":
		return v * em, true
	case "ex":
		return v * em / 2, true
	}
	if k, ok := svgUnits[u]; ok {
		return v * k, true
	}
	return 0, false
}

// parseTransform reads an SVG transform list.
func parseTransform(s string) matrix {
	m := identity
	for {
		s = strings.TrimLeft(s, " ,\t\n\r")
		open := strings.IndexByte(s, '(')
		if open < 0 {
			return m
		}
		name := strings.TrimSpace(s[:open])
		close := strings.IndexByte(s[open:], ')')
		if close < 0 {
			return m
		}
		a := svgNumbers(s[open+1 : open+close])
		s = s[open+close+1:]
		arg := func(i int, d float64) float64 {
			if i < len(a) {
				return a[i]
			}
			return d
		}
		var t matrix
		switch name {
		case "matrix":
			if len(a) != 6 {
				return m
			}
			t = matrix{a[0], a[1], a[2], a[3], a[4], a[5]}
		case "translate":
			t = matrix{1, 0, 0, 1, arg(0, 0), arg(1, 0)}
		case "scale":
			sx := arg(0, 1)
			t = matrix{sx, 0, 0, arg(1, sx), 0, 0}
		case "rotate":
			sn, cs := math.Sincos(arg(0, 0) * math.Pi / 180)
			cx, cy := arg(1, 0), arg(2, 0)
			t = matrix{1, 0, 0, 1, cx, cy}.mul(matrix{cs, sn, -sn, cs, 0, 0}).mul(matrix{1, 0, 0, 1, -cx, -cy})
		case "skewX":
			t = matrix{1, 0, math.Tan(arg(0, 0) * math.Pi / 180), 1, 0, 0}
		case "skewY":
			t = matrix{1, math.Tan(arg(0, 0) * math.Pi / 180), 0, 1, 0, 0}
		default:
			return m
		}
		m = m.mul(t)
	}
}

// parsePathData reads SVG path data into a path, up to the first error.
func parsePathData(d string) *path {
	p := &path{}
	var cmd byte
	var cx, cy, sx, sy float64 // current point, sub-path start
	var lcx, lcy float64       // last control point
	var lastCmd byte
	s := d
	nums := func(n int) ([]float64, bool) {
		out := make([]float64, n)
		for i := range out {
			v, rest, ok := svgNumber(s)
			if !ok {
				return nil, false
			}
			out[i] = v
			s = rest
		}
		return out, true
	}
	flag := func() (bool, bool) {
		s = strings.TrimLeft(s, " ,\t\n\r")
		if s == "" || (s[0] != '0' && s[0] != '1') {
			return false, false
		}
		f := s[0] == '1'
		s = s[1:]
		return f, true
	}
	for {
		s = strings.TrimLeft(s, " ,\t\n\r")
		if s == "" {
			return p
		}
		c := s[0]
		if strings.IndexByte("MmLlHhVvCcSsQqTtAaZz", c) >= 0 {
			cmd = c
			s = s[1:]
		} else if cmd == 0 {
			return p
		}
		rel := cmd >= 'a'
		ox, oy := 0.0, 0.0
		if rel {
			ox, oy = cx, cy
		}
		switch cmd {
		case 'M', 'm':
			a, ok := nums(2)
			if !ok {
				return p
			}
			cx, cy = a[0]+ox, a[1]+oy
			sx, sy = cx, cy
			p.moveTo(cx, cy)
			// further pairs are line-tos
			if cmd == 'M' {
				cmd = 'L'
			} else {
				cmd = 'l'
			}
			lastCmd = 'M'
			continue
		case 'L', 'l':
			a, ok := nums(2)
			if !ok {
				return p
			}
			cx, cy = a[0]+ox, a[1]+oy
			p.lineTo(cx, cy)
		case 'H', 'h':
			a, ok := nums(1)
			if !ok {
				return p
			}
			cx = a[0] + ox
			p.lineTo(cx, cy)
		case 'V', 'v':
			a, ok := nums(1)
			if !ok {
				return p
			}
			cy = a[0] + oy
			p.lineTo(cx, cy)
		case 'C', 'c':
			a, ok := nums(6)
			if !ok {
				return p
			}
			p.cubicTo(a[0]+ox, a[1]+oy, a[2]+ox, a[3]+oy, a[4]+ox, a[5]+oy)
			lcx, lcy = a[2]+ox, a[3]+oy
			cx, cy = a[4]+ox, a[5]+oy
		case 'S', 's':
			a, ok := nums(4)
			if !ok {
				return p
			}
			x1, y1 := cx, cy
			if lastCmd == 'C' || lastCmd == 'S' {
				x1, y1 = 2*cx-lcx, 2*cy-lcy
			}
			p.cubicTo(x1, y1, a[0]+ox, a[1]+oy, a[2]+ox, a[3]+oy)
			lcx, lcy = a[0]+ox, a[1]+oy
			cx, cy = a[2]+ox, a[3]+oy
		case 'Q', 'q':
			a, ok := nums(4)
			if !ok {
				return p
			}
			p.quadTo(a[0]+ox, a[1]+oy, a[2]+ox, a[3]+oy)
			lcx, lcy = a[0]+ox, a[1]+oy
			cx, cy = a[2]+ox, a[3]+oy
		case 'T', 't':
			a, ok := nums(2)
			if !ok {
				return p
			}
			x1, y1 := cx, cy
			if lastCmd == 'Q' || lastCmd == 'T' {
				x1, y1 = 2*cx-lcx, 2*cy-lcy
			}
			p.quadTo(x1, y1, a[0]+ox, a[1]+oy)
			lcx, lcy = x1, y1
			cx, cy = a[0]+ox, a[1]+oy
		case 'A', 'a':
			a, ok := nums(3)
			if !ok {
				return p
			}
			large, ok1 := flag()
			sweep, ok2 := flag()
			e, ok3 := nums(2)
			if !ok1 || !ok2 || !ok3 {
				return p
			}
			x, y := e[0]+ox, e[1]+oy
			arcTo(p, cx, cy, a[0], a[1], a[2], large, sweep, x, y)
			cx, cy = x, y
		case 'Z', 'z':
			p.close()
			cx, cy = sx, sy
		}
		lastCmd = cmd &^ 0x20 // upper case
	}
}

// arcTo adds an SVG elliptical arc from (x0, y0) to (x, y).
func arcTo(p *path, x0, y0, rx, ry, rotDeg float64, large, sweep bool, x, y float64) {
	if x0 == x && y0 == y {
		return
	}
	rx, ry = math.Abs(rx), math.Abs(ry)
	if rx == 0 || ry == 0 {
		p.lineTo(x, y)
		return
	}
	phi := rotDeg * math.Pi / 180
	sn, cs := math.Sincos(phi)
	dx, dy := (x0-x)/2, (y0-y)/2
	x1p, y1p := cs*dx+sn*dy, -sn*dx+cs*dy
	lambda := x1p*x1p/(rx*rx) + y1p*y1p/(ry*ry)
	if lambda > 1 {
		k := math.Sqrt(lambda)
		rx, ry = rx*k, ry*k
	}
	num := rx*rx*ry*ry - rx*rx*y1p*y1p - ry*ry*x1p*x1p
	den := rx*rx*y1p*y1p + ry*ry*x1p*x1p
	co := 0.0
	if den > 0 && num > 0 {
		co = math.Sqrt(num / den)
	}
	if large == sweep {
		co = -co
	}
	cxp, cyp := co*rx*y1p/ry, -co*ry*x1p/rx
	cx := cs*cxp - sn*cyp + (x0+x)/2
	cy := sn*cxp + cs*cyp + (y0+y)/2
	ang := func(ux, uy, vx, vy float64) float64 {
		a := math.Atan2(ux*vy-uy*vx, ux*vx+uy*vy)
		return a
	}
	t1 := ang(1, 0, (x1p-cxp)/rx, (y1p-cyp)/ry)
	dt := ang((x1p-cxp)/rx, (y1p-cyp)/ry, (-x1p-cxp)/rx, (-y1p-cyp)/ry)
	if !sweep && dt > 0 {
		dt -= 2 * math.Pi
	} else if sweep && dt < 0 {
		dt += 2 * math.Pi
	}
	p.ellipse(cx, cy, rx, ry, phi, t1, t1+dt, !sweep)
}

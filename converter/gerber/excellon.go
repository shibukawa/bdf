package gerber

import (
	"bytes"
	"math"
	"strconv"
	"strings"
)

// excellonParser reads an Excellon drill file: tools and the holes, slots
// and routed paths drilled with them.
type excellonParser struct {
	warnf func(string, ...any)
	img   *image
	attrs map[string]string

	unit float64 // mm per unit
	// the digits of coordinates without a decimal point, and whether their
	// leading zeros are kept (LZ: the trailing ones are left out) or their
	// trailing ones (TZ, the default)
	intDigits, decDigits int
	lz                   bool
	digitsSet            bool
	// decimal: numbers without a point are whole units (KiCad's decimal
	// format)
	decimal     bool
	incremental bool

	header bool
	tools  map[int]float64 // diameters in mm
	shapes map[float64]*shape
	tool   float64 // the diameter of the current tool
	cur    vec
	// routing: G00 moves, G01 to G03 draw while the tool is down (M15)
	mode int // 5 drill, 0 move, 1 linear, 2 clockwise, 3 counterclockwise
	down bool
	// upDown: the file lifts and lowers the tool (M15, M16) when it routes
	upDown   bool
	plated   string
	toolSeen bool
}

// parseExcellon reads an Excellon drill file.
func parseExcellon(data []byte, count *int, warnf func(string, ...any)) *layer {
	p := &excellonParser{
		warnf: warnf, img: newImage(count), attrs: map[string]string{},
		unit: 25.4, intDigits: 2, decDigits: 4, tools: map[int]float64{}, shapes: map[float64]*shape{}, mode: 5,
	}
	// lines end with LF, CR LF or CR
	for _, line := range bytes.FieldsFunc(data, func(r rune) bool { return r == '\n' || r == '\r' }) {
		s := strings.TrimSpace(string(line))
		if s == "" {
			continue
		}
		if s[0] == ';' || s[0] == '(' {
			p.comment(s)
			continue
		}
		if i := strings.IndexByte(s, ';'); i >= 0 {
			s = strings.TrimSpace(s[:i])
		}
		if p.line(s) {
			break
		}
		if p.img.over() {
			p.warnf("%v", errTooLarge)
			break
		}
	}
	return &layer{img: p.img, drill: true, attrs: p.attrs, plated: p.plated}
}

// comment reads what comments tell: X2 attributes (KiCad), the plating and
// the number format (Altium, KiCad).
func (p *excellonParser) comment(s string) {
	c := strings.TrimSpace(strings.TrimLeft(s, ";("))
	switch {
	case strings.HasPrefix(c, "#@!"):
		a := strings.TrimSpace(c[3:])
		if strings.HasPrefix(a, "TF.") {
			name, val, _ := strings.Cut(a[3:], ",")
			p.attrs[name] = val
		}
	case strings.HasPrefix(c, "TYPE="):
		switch strings.ToUpper(c[5:]) {
		case "PLATED":
			p.plated = "plated"
		case "NON_PLATED", "NONPLATED":
			p.plated = "nonplated"
		}
	case strings.HasPrefix(c, "FILE_FORMAT="):
		// Altium: FILE_FORMAT=2:5
		if a, b, ok := strings.Cut(c[12:], ":"); ok {
			p.setDigits(a, b)
		}
	case strings.HasPrefix(c, "FORMAT={"):
		// KiCad: FORMAT={3:3/ absolute / metric / suppress trailing zeros}
		f := c[8:]
		if a, b, ok := strings.Cut(f, ":"); ok {
			if i := strings.IndexByte(b, '/'); i >= 0 {
				b = b[:i]
			}
			p.setDigits(a, b)
		}
		switch {
		case strings.Contains(f, "decimal"):
			p.decimal = true
		case strings.Contains(f, "suppress trailing"):
			p.lz = true
		case strings.Contains(f, "suppress leading"):
			p.lz = false
		}
	}
}

func (p *excellonParser) setDigits(a, b string) {
	n, err1 := strconv.Atoi(strings.TrimSpace(a))
	m, err2 := strconv.Atoi(strings.TrimSpace(b))
	if err1 == nil && err2 == nil && n >= 0 && n <= 9 && m >= 0 && m <= 9 {
		p.intDigits, p.decDigits, p.digitsSet = n, m, true
	}
}

// line handles a line; it reports the end of the program.
func (p *excellonParser) line(s string) bool {
	u := strings.ToUpper(s)
	switch {
	case u == "M48":
		p.header = true
		return false
	case u == "%" || u == "M95":
		p.header = false
		return false
	case u == "M30" || u == "M00":
		return true
	case strings.HasPrefix(u, "METRIC") || strings.HasPrefix(u, "INCH"):
		p.units(u)
		return false
	case u == "M71":
		p.setUnit(1)
		return false
	case u == "M72":
		p.setUnit(25.4)
		return false
	case strings.HasPrefix(u, "ICI"):
		p.incremental = strings.Contains(u, "ON")
		return false
	case strings.HasPrefix(u, "FMAT") || strings.HasPrefix(u, "VER") || strings.HasPrefix(u, "ATC") ||
		strings.HasPrefix(u, "DETECT") || strings.HasPrefix(u, "BLKD") || strings.HasPrefix(u, "SBK") ||
		strings.HasPrefix(u, "TCST") || strings.HasPrefix(u, "OM") || strings.HasPrefix(u, "CP,") || u == "R,T" || u == "R,C":
		return false
	case strings.HasPrefix(u, "M97") || strings.HasPrefix(u, "M98"):
		p.warnf("text drilled with M97 and M98 is left out")
		return false
	}
	return p.words(u)
}

// units reads METRIC or INCH and what follows it: LZ or TZ, and a
// number format ("000.000").
func (p *excellonParser) units(u string) {
	parts := strings.Split(u, ",")
	if parts[0] == "METRIC" {
		p.setUnit(1)
	} else {
		p.setUnit(25.4)
	}
	for _, f := range parts[1:] {
		switch f = strings.TrimSpace(f); {
		case f == "LZ":
			p.lz = true
		case f == "TZ":
			p.lz = false
		case strings.Contains(f, "."):
			a, b, _ := strings.Cut(f, ".")
			p.intDigits, p.decDigits, p.digitsSet = len(a), len(b), true
		}
	}
}

func (p *excellonParser) setUnit(u float64) {
	p.unit = u
	if !p.digitsSet {
		if u == 1 {
			p.intDigits, p.decDigits = 3, 3
		} else {
			p.intDigits, p.decDigits = 2, 4
		}
	}
}

// field is a letter of a line and the text of its value.
type field struct {
	c byte
	v string
}

func splitFields(s string) []field {
	var out []field
	for i := 0; i < len(s); {
		c := s[i]
		j := i + 1
		for j < len(s) && (isDigit(s[j]) || s[j] == '.' || s[j] == '-' || s[j] == '+') {
			j++
		}
		if c >= 'A' && c <= 'Z' {
			out = append(out, field{c, s[i+1 : j]})
		}
		i = j
	}
	return out
}

// number reads a number of a line (in the file's unit, as it is).
func number(v string) float64 {
	f, _ := strconv.ParseFloat(v, 64)
	return f
}

// coord reads a coordinate in mm.
func (p *excellonParser) coord(v string) float64 {
	neg := strings.HasPrefix(v, "-")
	v = strings.TrimLeft(v, "+-")
	var f float64
	if strings.Contains(v, ".") || p.decimal {
		f, _ = strconv.ParseFloat(v, 64)
	} else {
		if p.lz && len(v) < p.intDigits+p.decDigits {
			v += strings.Repeat("0", p.intDigits+p.decDigits-len(v))
		}
		n, _ := strconv.ParseFloat(v, 64)
		f = n / math.Pow10(p.decDigits)
	}
	if neg {
		f = -f
	}
	return f * p.unit
}

// words handles a line of codes and coordinates.
func (p *excellonParser) words(s string) bool {
	fs := splitFields(s)
	if len(fs) == 0 {
		return false
	}
	if fs[0].c == 'T' {
		p.toolLine(fs)
		return false
	}
	if p.header {
		return false
	}
	g, slot, repeat := -1, -1, 0
	for k, f := range fs {
		switch f.c {
		case 'G':
			n, _ := strconv.Atoi(f.v)
			if n == 85 {
				slot = k
			} else if slot < 0 {
				g = n
			}
		case 'M':
			n, _ := strconv.Atoi(f.v)
			switch n {
			case 15:
				p.down, p.upDown = true, true
			case 16, 17:
				p.down, p.upDown = false, true
			case 30, 0:
				return true
			case 71:
				p.setUnit(1)
			case 72:
				p.setUnit(25.4)
			}
		case 'R':
			repeat, _ = strconv.Atoi(f.v)
		}
	}
	switch g {
	case 0, 1, 2, 3, 5:
		p.mode = g
		if g == 5 {
			p.down = false
		}
	case 90:
		p.incremental = false
	case 91:
		p.incremental = true
	}
	if repeat > 0 {
		// R: the hole repeated n times, each moved by the offsets
		var d vec
		for _, f := range fs {
			switch f.c {
			case 'X':
				d.X = p.coord(f.v)
			case 'Y':
				d.Y = p.coord(f.v)
			}
		}
		for range min(repeat, maxRepeat) {
			p.cur = p.cur.add(d)
			p.hole(p.cur)
		}
		return false
	}
	end := len(fs)
	if slot >= 0 {
		end = slot
	}
	to, moved, arc := p.point(fs[:end])
	if slot >= 0 {
		// X Y G85 X Y: a slot between the two points
		b, _, _ := p.point(fs[slot+1:])
		p.cur = b
		p.route(to, seg{to: b})
		return false
	}
	if !moved {
		return false
	}
	from := p.cur
	p.cur = to
	switch {
	case p.mode == 5:
		p.hole(to)
	case p.mode == 0:
		// a move with the tool up
	case p.down || !p.upDown:
		// routing (the tool is down; files without M15 and M16 cut
		// whenever they route)
		sg := seg{to: to}
		if p.mode == 2 || p.mode == 3 {
			if c, sw, ok := excellonArc(from, to, arc, p.mode == 2); ok {
				sg = seg{to: to, arc: true, c: c, sweep: sw}
			}
		}
		p.route(from, sg)
	}
	return false
}

// arcSpec is how a routed arc gives its centre: offsets I and J from its
// start, or its radius A.
type arcSpec struct {
	i, j, r     float64
	hasIJ, hasA bool
}

// point reads the coordinates of fields (the current point for an axis
// they leave out) and the arc's centre or radius.
func (p *excellonParser) point(fs []field) (pt vec, moved bool, arc arcSpec) {
	pt = p.cur
	for _, f := range fs {
		switch f.c {
		case 'X':
			pt.X, moved = p.axis(f.v, p.cur.X), true
		case 'Y':
			pt.Y, moved = p.axis(f.v, p.cur.Y), true
		case 'I':
			arc.i, arc.hasIJ = p.coord(f.v), true
		case 'J':
			arc.j, arc.hasIJ = p.coord(f.v), true
		case 'A':
			arc.r, arc.hasA = p.coord(f.v), true
		}
	}
	return pt, moved, arc
}

// axis reads a coordinate, added to cur in incremental mode.
func (p *excellonParser) axis(v string, cur float64) float64 {
	c := p.coord(v)
	if p.incremental {
		return cur + c
	}
	return c
}

// toolLine defines or selects a tool: T01C0.8 (with feeds and speeds), T01.
func (p *excellonParser) toolLine(fs []field) {
	n, err := strconv.Atoi(fs[0].v)
	if err != nil {
		return
	}
	for _, f := range fs[1:] {
		if f.c == 'C' {
			d := p.coord2(f.v)
			if finite(d) && d >= 0 && d < 1000 {
				p.tools[n] = d
			}
		}
	}
	if p.header {
		return
	}
	d, ok := p.tools[n]
	if !ok && n != 0 {
		p.warnf("tool T%d has no diameter; its holes are drawn 1 mm wide", n)
		d = 1
		p.tools[n] = d
	}
	p.tool = d
	p.toolSeen = true
}

// coord2 reads a tool diameter: a decimal number in the file's unit.
func (p *excellonParser) coord2(v string) float64 {
	if !strings.Contains(v, ".") {
		// diameters without a point are in the coordinate format
		return p.coord(v)
	}
	return number(v) * p.unit
}

// hole drills a hole at pt with the current tool.
func (p *excellonParser) hole(pt vec) {
	if p.tool <= 0 {
		if !p.toolSeen {
			p.warnf("holes drilled before a tool is chosen are left out")
		}
		return
	}
	s, ok := p.shapes[p.tool]
	if !ok {
		s, _ = standardShape("C", []float64{p.tool})
		p.shapes[p.tool] = s
	}
	p.img.flash(false, s, pt)
}

// route draws a slot or a routed path from a with the current tool.
func (p *excellonParser) route(a vec, s seg) {
	if p.tool <= 0 {
		return
	}
	if a == s.to && !s.arc {
		p.hole(a)
		return
	}
	p.img.line(false, p.tool, a, s)
}

// excellonArc returns the centre and sweep of an arc of a routed path from
// a to b: about a centre at the offsets I, J from a, or of radius r (the
// shorter of the two arcs).
func excellonArc(a, b vec, arc arcSpec, cw bool) (vec, float64, bool) {
	var c vec
	r := arc.r
	switch {
	case arc.hasIJ:
		c = a.add(vec{arc.i, arc.j})
	case arc.hasA && r > 0:
		m := a.add(b).mul(0.5)
		d := b.sub(a)
		l := d.len()
		if l == 0 || l > 2*r+1e-9 {
			return vec{}, 0, false
		}
		h := math.Sqrt(math.Max(r*r-l*l/4, 0))
		n := vec{-d.Y, d.X}.mul(1 / l)
		// the centre on the side that makes the shorter arc in the direction
		if cw {
			c = m.add(n.mul(-h))
		} else {
			c = m.add(n.mul(h))
		}
	default:
		return vec{}, 0, false
	}
	a0 := math.Atan2(a.Y-c.Y, a.X-c.X)
	sw := math.Atan2(b.Y-c.Y, b.X-c.X) - a0
	if cw {
		for sw >= 0 {
			sw -= 2 * math.Pi
		}
	} else {
		for sw <= 0 {
			sw += 2 * math.Pi
		}
	}
	return c, sw, true
}

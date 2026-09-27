package gerber

import (
	"bytes"
	"math"
	"strconv"
	"strings"
)

// layer is what a file holds: its image in mm and what the file says
// about itself.
type layer struct {
	name string // the file's name (in an archive, its path)
	img  *image
	// drill: an Excellon file
	drill bool
	// negative: the image is the absence of material (%IPNEG, or
	// TF.FilePolarity,Negative in a copper layer)
	negative bool
	// attrs holds the file attributes (TF), by name without the dot
	// ("FileFunction"), and the comments of drill files that carry them
	attrs map[string]string
	// plated tells a drill file's holes: "plated", "nonplated" or "" (not
	// said)
	plated string
	fn     function
}

// gerberParser reads a Gerber file (RS-274X, and the deprecated commands
// of older files) into an image.
type gerberParser struct {
	warnf func(format string, args ...any)
	root  *image
	// blocks are the step and repeat and aperture blocks being read
	blocks []block

	// the coordinate format (FS): digits of the integer and decimal parts
	xi, xd, yi, yd int
	trailing       bool // trailing zeros left out (deprecated)
	incremental    bool // incremental coordinates (deprecated)
	fsSet          bool
	unit           float64 // mm per unit of the file
	unitSet        bool

	apertures map[int]*aperture
	macros    map[string]*macro
	ap        *aperture
	cur       vec
	interp    int  // 1 linear, 2 clockwise, 3 counterclockwise
	multi     bool // multi quadrant mode (G75)
	lastOp    int  // D01, D02 or D03: the operation of a coordinate block without one (deprecated)
	clear     bool // clear polarity (LPC)

	// region mode (G36): the contours so far and the one being drawn
	region   bool
	contours []contour
	cont     *contour

	// aperture transformations (LM, LR, LS)
	mirrorX, mirrorY bool
	rot, scale       float64
	// image transformations of older files (IR, MI, OF, SF, AS)
	imgRot           float64
	imgMirA, imgMirB bool
	imgOff, imgScale vec
	swapAxes         bool
	imgM             affine

	negative bool
	attrs    map[string]string
	done     bool
}

// block is a step and repeat block (SR) or an aperture block (AB) being
// read.
type block struct {
	aperture bool
	code     int // the D code of an aperture block
	nx, ny   int
	d        vec // step and repeat distances
	img      *image
}

// sink returns the image objects go to.
func (p *gerberParser) sink() *image {
	if n := len(p.blocks); n > 0 {
		return p.blocks[n-1].img
	}
	return p.root
}

// maxRepeat bounds the copies of a step and repeat block.
const maxRepeat = 100_000

// parseGerber reads a Gerber file.
func parseGerber(data []byte, count *int, warnf func(string, ...any)) (*layer, error) {
	p := &gerberParser{
		warnf: warnf, root: newImage(count), unit: 25.4, xi: 2, xd: 4, yi: 2, yd: 4,
		apertures: map[int]*aperture{}, macros: map[string]*macro{}, interp: 1,
		scale: 1, imgScale: vec{1, 1}, imgM: identity, attrs: map[string]string{},
	}
	i := 0
	for i < len(data) && !p.done {
		c := data[i]
		switch {
		case c == '%':
			j := bytes.IndexByte(data[i+1:], '%')
			if j < 0 {
				p.warnf("the file ends inside an extended command")
				i = len(data)
				break
			}
			p.extended(clean(data[i+1 : i+1+j]))
			i += j + 2
		case c <= ' ':
			i++
		default:
			j := bytes.IndexByte(data[i:], '*')
			if j < 0 {
				if s := strings.TrimSpace(string(clean(data[i:]))); s != "" && s != "M02" && s != "M00" && s != "M01" {
					p.warnf("the file ends inside a block")
				}
				i = len(data)
				break
			}
			p.word(clean(data[i : i+j]))
			i += j + 1
		}
		if p.root.over() {
			p.warnf("%v", errTooLarge)
			break
		}
	}
	if p.region {
		p.endRegion()
	}
	for len(p.blocks) > 0 {
		b := p.blocks[len(p.blocks)-1]
		if b.aperture {
			p.warnf("an aperture block is not closed")
			p.blocks = p.blocks[:len(p.blocks)-1]
		} else {
			p.endRepeat()
		}
	}
	if !p.fsSet {
		p.warnf("the file has no format specification (FS); coordinates are read as 2.4 digits with leading zeros left out")
	}
	if !p.unitSet {
		p.warnf("the file gives no unit (MO); inches are assumed")
	}
	l := &layer{img: p.root, negative: p.negative, attrs: p.attrs}
	if strings.EqualFold(p.attrs["FilePolarity"], "Negative") && strings.HasPrefix(p.attrs["FileFunction"], "Copper") {
		l.negative = true
	}
	return l, nil
}

// clean returns b without line breaks.
func clean(b []byte) []byte {
	if bytes.IndexByte(b, '\n') < 0 && bytes.IndexByte(b, '\r') < 0 {
		return b
	}
	out := make([]byte, 0, len(b))
	for _, c := range b {
		if c != '\n' && c != '\r' {
			out = append(out, c)
		}
	}
	return out
}

// extended handles the content of an extended command (%...%): one
// command, or several in older files, each ending with '*'.
func (p *gerberParser) extended(b []byte) {
	s := strings.TrimSpace(string(b))
	if strings.HasPrefix(s, "AM") {
		stmts := strings.Split(s, "*")
		name := strings.TrimSpace(stmts[0][2:])
		var body []string
		for _, st := range stmts[1:] {
			if strings.TrimSpace(st) != "" {
				body = append(body, st)
			}
		}
		p.macros[name] = &macro{name: name, stmts: body}
		return
	}
	for _, cmd := range strings.Split(s, "*") {
		if cmd = strings.TrimSpace(cmd); cmd != "" {
			p.command(cmd)
		}
	}
}

// command handles one extended command.
func (p *gerberParser) command(s string) {
	code, rest := s, ""
	if len(s) >= 2 {
		code, rest = s[:2], s[2:]
	}
	switch code {
	case "FS":
		p.formatSpec(rest)
	case "MO":
		switch strings.ToUpper(strings.TrimSpace(rest)) {
		case "MM":
			p.unit, p.unitSet = 1, true
		case "IN":
			p.unit, p.unitSet = 25.4, true
		default:
			p.warnf("unknown unit %q", rest)
		}
	case "AD":
		p.defineAperture(rest)
	case "AB":
		p.apertureBlock(rest)
	case "LP":
		switch rest {
		case "D":
			p.clear = false
		case "C":
			p.clear = true
		default:
			p.warnf("unknown polarity %q", rest)
		}
	case "LM":
		switch rest {
		case "N":
			p.mirrorX, p.mirrorY = false, false
		case "X":
			p.mirrorX, p.mirrorY = true, false
		case "Y":
			p.mirrorX, p.mirrorY = false, true
		case "XY":
			p.mirrorX, p.mirrorY = true, true
		default:
			p.warnf("unknown mirroring %q", rest)
		}
	case "LR":
		if v, err := strconv.ParseFloat(rest, 64); err == nil && finite(v) {
			p.rot = v
		}
	case "LS":
		if v, err := strconv.ParseFloat(rest, 64); err == nil && finite(v) && v > 0 && v < 1e6 {
			p.scale = v
		}
	case "SR":
		p.stepRepeat(rest)
	case "TF":
		name, val, _ := strings.Cut(strings.TrimPrefix(rest, "."), ",")
		p.attrs[name] = val
	case "TA", "TO", "TD":
		// aperture and object attributes: not needed to draw
	case "IP":
		switch rest {
		case "POS":
			p.negative = false
		case "NEG":
			p.negative = true
		}
	case "IR":
		if v, err := strconv.Atoi(rest); err == nil {
			p.imgRot = float64(v)
			p.imageTransform()
		}
	case "MI":
		a, b := axisValues(rest)
		p.imgMirA, p.imgMirB = a.X != 0, b.X != 0
		p.imageTransform()
	case "OF":
		a, b := axisValues(rest)
		p.imgOff = vec{a.X, b.X}
		p.imageTransform()
	case "SF":
		a, b := axisValues(rest)
		p.imgScale = vec{a.X, b.X}
		if a.Y == 0 {
			p.imgScale.X = 1
		}
		if b.Y == 0 {
			p.imgScale.Y = 1
		}
		p.imageTransform()
	case "AS":
		p.swapAxes = rest == "AYBX"
		p.imageTransform()
	case "IN", "LN", "IC", "IJ", "KO", "IO", "PF", "RO", "PK", "PM":
		// names, and settings of film plotters that do not change the image
	default:
		p.warnf("unknown extended command %q", code)
	}
}

// axisValues reads the A and B values of an image command ("A1B0",
// "A1.5B-2"); the Y of each is 1 when the value is present.
func axisValues(s string) (a, b vec) {
	for len(s) > 0 {
		c := s[0]
		j := 1
		for j < len(s) && (s[j] >= '0' && s[j] <= '9' || s[j] == '.' || s[j] == '-' || s[j] == '+') {
			j++
		}
		v, err := strconv.ParseFloat(s[1:j], 64)
		if err == nil && finite(v) {
			switch c {
			case 'A':
				a = vec{v, 1}
			case 'B':
				b = vec{v, 1}
			}
		}
		s = s[j:]
	}
	return a, b
}

// imageTransform makes the transformation of the image commands of older
// files: axes swapped, scaled, mirrored, turned and moved, in that order.
func (p *gerberParser) imageTransform() {
	m := identity
	if p.swapAxes {
		m = affine{0, 1, 1, 0, 0, 0}
	}
	sx, sy := p.imgScale.X, p.imgScale.Y
	if sx != sy || sx <= 0 {
		// shapes keep their circles only under uniform scales
		p.warnf("image scale factors that differ between the axes are drawn as their mean")
		sx = math.Sqrt(math.Abs(sx * sy))
		if sx == 0 {
			sx = 1
		}
		sy = sx
	}
	m = affine{sx, 0, 0, sy, 0, 0}.mul(m)
	mx, my := 1.0, 1.0
	if p.imgMirA {
		mx = -1
	}
	if p.imgMirB {
		my = -1
	}
	m = affine{mx, 0, 0, my, 0, 0}.mul(m)
	m = rotate(p.imgRot).mul(m)
	m = translate(p.imgOff.mul(p.unit)).mul(m)
	p.imgM = m
}

// formatSpec reads FS: the zero omission, the notation and the digits of X
// and Y ("LAX26Y26").
func (p *gerberParser) formatSpec(s string) {
	p.fsSet = true
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case 'L', 'D':
			p.trailing = false
		case 'T':
			p.trailing = true
		case 'A':
			p.incremental = false
		case 'I':
			p.incremental = true
		case 'X', 'Y':
			if i+2 < len(s) && isDigit(s[i+1]) && isDigit(s[i+2]) {
				n, d := int(s[i+1]-'0'), int(s[i+2]-'0')
				if c == 'X' {
					p.xi, p.xd = n, d
				} else {
					p.yi, p.yd = n, d
				}
				i += 2
			}
		case 'N', 'G', 'M':
			for i+1 < len(s) && isDigit(s[i+1]) {
				i++
			}
		}
	}
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

// defineAperture reads AD: "D10C,0.5", "D11RoundRect,0.25X...".
func (p *gerberParser) defineAperture(s string) {
	if !strings.HasPrefix(s, "D") {
		p.warnf("bad aperture definition %q", s)
		return
	}
	j := 1
	for j < len(s) && isDigit(s[j]) {
		j++
	}
	code, err := strconv.Atoi(s[1:j])
	if err != nil {
		p.warnf("bad aperture definition %q", s)
		return
	}
	name, params, _ := strings.Cut(s[j:], ",")
	values, err := parseParams(params)
	if err != nil {
		p.warnf("aperture D%d: %v", code, err)
		return
	}
	var sh *shape
	switch name {
	case "C", "R", "O", "P":
		sh, err = standardShape(name, values)
	default:
		m := p.macros[name]
		if m == nil {
			p.warnf("aperture D%d uses the undefined macro %q", code, name)
			return
		}
		sh, err = m.shape(values, func(msg string) { p.warnf("%s", msg) })
	}
	if err != nil {
		p.warnf("aperture D%d: %v", code, err)
		return
	}
	u := p.unit
	scaled := sh.transform(affine{u, 0, 0, u, 0, 0})
	p.apertures[code] = &aperture{shape: scaled}
}

// apertureBlock reads AB: "D12" opens a block aperture, "" closes it.
func (p *gerberParser) apertureBlock(s string) {
	if s == "" {
		n := len(p.blocks)
		if n == 0 || !p.blocks[n-1].aperture {
			p.warnf("AB closes no aperture block")
			return
		}
		b := p.blocks[n-1]
		p.blocks = p.blocks[:n-1]
		p.apertures[b.code] = &aperture{shape: &shape{block: b.img, bounds: b.img.bounds}}
		return
	}
	code, err := strconv.Atoi(strings.TrimPrefix(s, "D"))
	if err != nil || len(p.blocks) >= 32 {
		p.warnf("bad aperture block %q", s)
		return
	}
	p.blocks = append(p.blocks, block{aperture: true, code: code, img: newImage(p.root.count)})
}

// stepRepeat reads SR: "X3Y2I5.0J4.0" opens a block repeated 3 × 2
// times, "" (or X1Y1) closes the one open.
func (p *gerberParser) stepRepeat(s string) {
	if n := len(p.blocks); n > 0 && !p.blocks[n-1].aperture {
		p.endRepeat()
	}
	nx, ny := 1, 1
	var d vec
	for i := 0; i < len(s); {
		c := s[i]
		j := i + 1
		for j < len(s) && (isDigit(s[j]) || s[j] == '.' || s[j] == '-' || s[j] == '+') {
			j++
		}
		v, _ := strconv.ParseFloat(s[i+1:j], 64)
		switch c {
		case 'X':
			nx = int(v)
		case 'Y':
			ny = int(v)
		case 'I':
			d.X = v * p.unit
		case 'J':
			d.Y = v * p.unit
		}
		i = j
	}
	if nx <= 1 && ny <= 1 {
		return
	}
	if nx < 1 || ny < 1 || nx*ny > maxRepeat || !finite(d.X) || !finite(d.Y) {
		p.warnf("step and repeat %q is not drawn repeated", s)
		nx, ny = 1, 1
	}
	p.blocks = append(p.blocks, block{nx: nx, ny: ny, d: d, img: newImage(p.root.count)})
}

// endRepeat closes the step and repeat block on top and copies it.
func (p *gerberParser) endRepeat() {
	n := len(p.blocks)
	b := p.blocks[n-1]
	p.blocks = p.blocks[:n-1]
	lin := p.imgM
	lin[4], lin[5] = 0, 0
	dst := p.sink()
	for iy := range b.ny {
		for ix := range b.nx {
			off := lin.apply(vec{float64(ix) * b.d.X, float64(iy) * b.d.Y})
			dst.insert(b.img, translate(off), false)
			if dst.over() {
				return
			}
		}
	}
}

// word handles a word command: G, D and M codes and coordinates.
func (p *gerberParser) word(b []byte) {
	s := strings.TrimSpace(string(b))
	if s == "" {
		return
	}
	if isComment(s) {
		// KiCad writes attributes into comments: G04 #@! TF.FileFunction,...
		c := strings.TrimPrefix(s, "G0")
		c = strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(c, "G"), "4"))
		if strings.HasPrefix(c, "#@!") {
			a := strings.TrimSpace(c[3:])
			if strings.HasPrefix(a, "TF.") {
				name, val, _ := strings.Cut(a[3:], ",")
				p.attrs[name] = val
			}
		}
		return
	}
	var x, y, i, j string
	var hasX, hasY bool
	d, m := -1, -1
	var gs []int
	for k := 0; k < len(s); {
		c := s[k]
		switch c {
		case ' ':
			k++
		case 'G', 'D', 'M', 'N':
			e := k + 1
			for e < len(s) && isDigit(s[e]) {
				e++
			}
			n, err := strconv.Atoi(s[k+1 : e])
			if err != nil {
				p.warnf("bad command %q", s)
				return
			}
			switch c {
			case 'G':
				if n == 4 {
					return // a comment after other codes
				}
				gs = append(gs, n)
			case 'D':
				d = n
			case 'M':
				m = n
			}
			k = e
		case 'X', 'Y', 'I', 'J':
			e := k + 1
			for e < len(s) && (isDigit(s[e]) || s[e] == '-' || s[e] == '+' || s[e] == '.') {
				e++
			}
			v := s[k+1 : e]
			switch c {
			case 'X':
				x, hasX = v, true
			case 'Y':
				y, hasY = v, true
			case 'I':
				i = v
			case 'J':
				j = v
			}
			k = e
		default:
			p.warnf("unknown command %q", s)
			return
		}
	}
	for _, g := range gs {
		p.gcode(g)
	}
	if d >= 10 {
		p.ap = p.apertures[d]
		if p.ap == nil {
			p.warnf("aperture D%d is not defined", d)
		}
	}
	op := 0
	switch {
	case d >= 1 && d <= 3:
		op, p.lastOp = d, d
	case hasX || hasY || i != "" || j != "":
		op = p.lastOp
		if op == 0 {
			op = 1
		}
	}
	if op != 0 {
		to := p.cur
		if hasX {
			to.X = p.coord(x, p.xi, p.xd, p.cur.X)
		}
		if hasY {
			to.Y = p.coord(y, p.yi, p.yd, p.cur.Y)
		}
		// offsets are relative in either notation
		var off vec
		if i != "" {
			off.X = p.coord(i, p.xi, p.xd, 0)
		}
		if j != "" {
			off.Y = p.coord(j, p.yi, p.yd, 0)
		}
		p.operate(op, to, off)
	}
	switch m {
	case 0, 1, 2:
		p.done = true
	}
}

func isComment(s string) bool {
	return strings.HasPrefix(s, "G04") || strings.HasPrefix(s, "G4") && (len(s) == 2 || !isDigit(s[2]))
}

// coord reads a coordinate in the file's format and unit; cur is the
// current value, which incremental coordinates add to.
func (p *gerberParser) coord(s string, intDigits, decDigits int, cur float64) float64 {
	neg := false
	switch {
	case strings.HasPrefix(s, "-"):
		neg, s = true, s[1:]
	case strings.HasPrefix(s, "+"):
		s = s[1:]
	}
	var v float64
	if strings.Contains(s, ".") {
		v, _ = strconv.ParseFloat(s, 64)
	} else {
		if p.trailing && len(s) < intDigits+decDigits {
			s += strings.Repeat("0", intDigits+decDigits-len(s))
		}
		n, _ := strconv.ParseFloat(s, 64)
		v = n / math.Pow10(decDigits)
	}
	if neg {
		v = -v
	}
	v *= p.unit
	if p.incremental {
		return cur + v
	}
	return v
}

// gcode handles a G code.
func (p *gerberParser) gcode(g int) {
	switch g {
	case 1, 10, 11, 12:
		p.interp = 1
	case 2:
		p.interp = 2
	case 3:
		p.interp = 3
	case 36:
		if p.region {
			p.endRegion()
		}
		p.region = true
		p.contours, p.cont = nil, nil
	case 37:
		if p.region {
			p.endRegion()
		}
	case 54, 55:
		// select an aperture (with the D code) / prepare a flash
	case 70:
		p.unit, p.unitSet = 25.4, true
	case 71:
		p.unit, p.unitSet = 1, true
	case 74:
		p.multi = false
	case 75:
		p.multi = true
	case 90:
		p.incremental = false
	case 91:
		p.incremental = true
	default:
		p.warnf("unknown code G%02d", g)
	}
}

// apertureTransform returns the transformation of flashed and drawn
// apertures: mirrored, then turned, then scaled (LM, LR, LS), then the
// linear part of the image transformation.
func (p *gerberParser) apertureTransform() affine {
	m := identity
	if p.mirrorX {
		m[0] = -1
	}
	if p.mirrorY {
		m[3] = -1
	}
	m = rotate(p.rot).mul(m)
	if p.scale != 1 {
		m = affine{p.scale, 0, 0, p.scale, 0, 0}.mul(m)
	}
	lin := p.imgM
	lin[4], lin[5] = 0, 0
	return lin.mul(m)
}

// operate carries out D01 (interpolate), D02 (move) or D03 (flash) to a
// point; off holds I and J.
func (p *gerberParser) operate(op int, to, off vec) {
	from := p.cur
	p.cur = to
	switch op {
	case 1:
		s := seg{to: to}
		if p.interp != 1 {
			var ok bool
			if s, ok = p.arc(from, to, off, p.interp == 2); !ok {
				s = seg{to: to}
			}
		}
		if p.region {
			if p.cont == nil {
				p.cont = &contour{start: from}
			}
			p.cont.segs = append(p.cont.segs, s)
			return
		}
		p.draw(from, s)
	case 2:
		if p.region && p.cont != nil {
			p.contours = append(p.contours, *p.cont)
			p.cont = nil
		}
	case 3:
		if p.region {
			p.warnf("flashes in regions are left out")
			return
		}
		if p.ap == nil {
			p.warnf("a flash without an aperture is left out")
			return
		}
		p.sink().flash(p.clear, p.ap.at(p.apertureTransform()), p.imgM.apply(to))
	}
}

// arc returns the segment of a circular interpolation from a to b with the
// centre offset off, clockwise or not.
func (p *gerberParser) arc(a, b, off vec, cw bool) (seg, bool) {
	norm := func(sw float64) float64 {
		// into (0, 2π) counterclockwise, (-2π, 0) clockwise
		if cw {
			for sw >= 0 {
				sw -= 2 * math.Pi
			}
			for sw <= -2*math.Pi {
				sw += 2 * math.Pi
			}
		} else {
			for sw <= 0 {
				sw += 2 * math.Pi
			}
			for sw >= 2*math.Pi {
				sw -= 2 * math.Pi
			}
		}
		return sw
	}
	full := 2 * math.Pi
	if cw {
		full = -full
	}
	if p.multi {
		c := a.add(off)
		if a.sub(c).len() == 0 {
			return seg{}, false
		}
		sw := full
		if a != b {
			sw = norm(math.Atan2(b.Y-c.Y, b.X-c.X) - math.Atan2(a.Y-c.Y, a.X-c.X))
		}
		return seg{to: b, arc: true, c: c, sweep: sw}, true
	}
	// single quadrant: the signs of the offsets are those that make an arc
	// of at most 90°
	if a == b {
		return seg{}, false
	}
	best, bestErr := seg{}, math.Inf(1)
	for _, sx := range []float64{1, -1} {
		for _, sy := range []float64{1, -1} {
			c := a.add(vec{sx * math.Abs(off.X), sy * math.Abs(off.Y)})
			r0, r1 := a.sub(c).len(), b.sub(c).len()
			if r0 == 0 {
				continue
			}
			sw := norm(math.Atan2(b.Y-c.Y, b.X-c.X) - math.Atan2(a.Y-c.Y, a.X-c.X))
			e := math.Abs(r0 - r1)
			if math.Abs(sw) > math.Pi/2+1e-6 {
				e += 1e6
			}
			if e < bestErr {
				best, bestErr = seg{to: b, arc: true, c: c, sweep: sw}, e
			}
		}
	}
	return best, bestErr < math.Inf(1)
}

// draw draws a segment from a with the current aperture.
func (p *gerberParser) draw(a vec, s seg) {
	if p.ap == nil {
		p.warnf("a line drawn without an aperture is left out")
		return
	}
	sh := p.ap.at(p.apertureTransform())
	if sh.block != nil {
		p.warnf("lines drawn with block apertures are left out")
		return
	}
	im := p.sink()
	m := p.imgM
	ta := m.apply(a)
	ts := seg{to: m.apply(s.to), arc: s.arc}
	if s.arc {
		ts.c = m.apply(s.c)
		ts.sweep = s.sweep
		if m.det() < 0 {
			ts.sweep = -s.sweep
		}
	}
	w := sh.d
	if !sh.round {
		if s.arc {
			p.warnf("arcs drawn with apertures that are not circles are drawn as round lines")
			w = math.Min(sh.bounds.w(), sh.bounds.h())
		} else {
			if c, ok := sweep(sh.outline(), ta, ts.to); ok {
				im.fill(p.clear, c)
			}
			return
		}
	}
	if ta == ts.to && (!ts.arc || math.Abs(ts.sweep) < 1e-12) {
		// a draw of zero length is a dot
		if w > 0 {
			im.fill(p.clear, dot(ta, w))
		}
		return
	}
	im.line(p.clear, w, ta, ts)
}

// endRegion closes the region being read and fills its contours.
func (p *gerberParser) endRegion() {
	if p.cont != nil {
		p.contours = append(p.contours, *p.cont)
	}
	p.region, p.cont = false, nil
	im := p.sink()
	for _, c := range p.contours {
		if len(c.segs) == 0 {
			continue
		}
		c = c.transform(p.imgM)
		if c.end() != c.start {
			c.lineTo(c.start)
		}
		if len(c.segs) < 2 && !c.segs[0].arc {
			continue
		}
		im.fill(p.clear, c.orient(true))
	}
	p.contours = nil
}

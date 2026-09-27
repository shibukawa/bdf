package kicad

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
)

// transform is KiCad's symbol TRANSFORM: a library point (x, y), y down,
// goes to (x1·x + y1·y, x2·x + y2·y) about the symbol's position.
type transform struct{ x1, y1, x2, y2 float64 }

func (t transform) apply(x, y float64) (float64, float64) {
	return t.x1*x + t.y1*y, t.x2*x + t.y2*y
}

// then returns t modified by an incremental transform, as
// SCH_SYMBOL::SetOrientation composes them.
func (t transform) then(m transform) transform {
	return transform{
		x1: t.x1*m.x1 + t.x2*m.y1,
		y1: t.y1*m.x1 + t.y2*m.y1,
		x2: t.x1*m.x2 + t.x2*m.y2,
		y2: t.y1*m.x2 + t.y2*m.y2,
	}
}

var (
	tIdentity = transform{1, 0, 0, 1}
	tCCW      = transform{0, 1, -1, 0}
	tCW       = transform{0, -1, 1, 0}
	tMirrorY  = transform{-1, 0, 0, 1}
	tMirrorX  = transform{1, 0, 0, -1}
)

// symbolTransform returns the transform of a symbol placed at an angle
// (0, 90, 180 or 270) and mirrored ("x", "y" or "").
func symbolTransform(angle float64, mirror string) transform {
	t := tIdentity
	switch (int(math.Round(angle))%360 + 360) % 360 {
	case 90:
		t = t.then(tCCW)
	case 180:
		t = t.then(tCCW).then(tCCW)
	case 270:
		t = t.then(tCW)
	}
	switch mirror {
	case "x":
		t = t.then(tMirrorX)
	case "y":
		t = t.then(tMirrorY)
	}
	return t
}

// symInst is a symbol placed on a sheet.
type symInst struct {
	n           *node
	lib         *node
	x, y        float64
	t           transform
	unit, style int
	units       int
	ref         string
	dnp         bool
	power       bool
	nameOffset  float64
	hideNames   bool
	hideNumbers bool
	alternates  map[string]string // pin number → alternate name
	fields      map[string]string // field values by name, for text variables
}

// point maps a library point (as the file writes it, y up) onto the sheet.
func (si *symInst) point(x, y float64) cad.Point {
	wx, wy := si.t.apply(x, -y)
	return cad.Point{X: si.x + wx, Y: si.y + wy}
}

// symbolInstance prepares a placed symbol for drawing.
func (p *schPage) symbolInstance(sym *node) *symInst {
	name := sym.str("lib_name")
	if name == "" {
		name = sym.str("lib_id")
	}
	lib := p.s.file.libs[name]
	if lib == nil {
		p.c.warnOnce("lib:"+name, "symbol %s is not in the schematic's library cache: only its fields are drawn", name)
	}
	at := sym.child("at")
	x, y := at.xy()
	si := &symInst{n: sym, lib: lib, x: x, y: y, t: symbolTransform(at.num(2), sym.str("mirror")),
		style: 1, nameOffset: schPinNameOffset, alternates: map[string]string{}, fields: map[string]string{}}
	if s := sym.child("body_style"); s != nil {
		si.style = s.int(0)
	} else if s := sym.child("convert"); s != nil {
		si.style = s.int(0)
	}
	si.ref, si.unit = p.c.symbolInstance(sym, p.s, nil, p.v6refs)
	if si.unit <= 0 {
		si.unit = 1
	}
	si.dnp = sym.flag("dnp")
	for _, pin := range sym.children("pin") {
		if alt := pin.str("alternate"); alt != "" {
			si.alternates[pin.arg(0)] = alt
		}
	}
	for _, f := range sym.children("property") {
		si.fields[f.arg(0)] = f.arg(1)
	}
	if lib != nil {
		si.power = lib.child("power") != nil || lib.hasAtom("power")
		if pn := lib.child("pin_names"); pn != nil {
			if pn.child("offset") != nil {
				si.nameOffset = pn.numOf("offset", schPinNameOffset)
			}
			si.hideNames = pn.flag("hide")
		}
		si.hideNumbers = lib.child("pin_numbers").flag("hide")
		units := map[int]bool{}
		for _, u := range lib.children("symbol") {
			if n, _, ok := unitOf(u.arg(0)); ok && n > 0 {
				units[n] = true
			}
		}
		si.units = len(units)
	}
	return si
}

// unitOf reads the unit and body style a library unit's name ends with
// ("R_1_1": unit 1, style 1; 0 is common to all).
func unitOf(name string) (unit, style int, ok bool) {
	i := strings.LastIndexByte(name, '_')
	if i < 0 {
		return 0, 0, false
	}
	j := strings.LastIndexByte(name[:i], '_')
	if j < 0 {
		return 0, 0, false
	}
	u, e1 := strconv.Atoi(name[j+1 : i])
	s, e2 := strconv.Atoi(name[i+1:])
	if e1 != nil || e2 != nil {
		return 0, 0, false
	}
	return u, s, true
}

// items returns the drawing items of the unit and body style a symbol
// shows: the library symbol's own and those of its matching units.
func (si *symInst) items() []*node {
	if si.lib == nil {
		return nil
	}
	var out []*node
	add := func(n *node) {
		for _, it := range n.lists() {
			switch it.name {
			case "arc", "bezier", "circle", "polyline", "rectangle", "text", "text_box", "textbox", "pin":
				out = append(out, it)
			}
		}
	}
	add(si.lib)
	for _, u := range si.lib.children("symbol") {
		unit, style, ok := unitOf(u.arg(0))
		if !ok {
			continue
		}
		if (unit == 0 || unit == si.unit) && (style == 0 || style == si.style) {
			add(u)
		}
	}
	return out
}

// symbolBody draws a symbol's body and pins: the fills of its background
// (background), or its outlines, other fills, pins and texts.
func (p *schPage) symbolBody(si *symInst, background bool) {
	outline := dimIf(schOutline, si.dnp)
	for _, it := range si.items() {
		switch it.name {
		case "pin":
			if !background {
				p.pin(si, it)
			}
			continue
		case "text":
			if !background {
				p.libText(si, it)
			}
			continue
		case "text_box", "textbox":
			if !background {
				p.textBox(it, si)
			}
			continue
		}
		path, _ := shapePath(it, func(q cad.Point) cad.Point { return si.point(q.X, q.Y) })
		if path == nil {
			continue
		}
		f := parseFill(it.child("fill"))
		st := parseStroke(it.child("stroke"))
		color := outline
		if st.hasColor {
			color = dimIf(st.color, si.dnp)
		}
		if background {
			if f.kind == "background" {
				p.fillShape(path, f, color, si.dnp, st.width)
			}
			continue
		}
		if f.kind != "background" {
			p.fillShape(path, f, color, si.dnp, st.width)
		}
		if st.width < 0 || st.dash == "none" {
			continue
		}
		w := st.width
		if w <= 0 {
			w = p.c.sch.line
		}
		p.d.Stroke(path, p.pen(w, st.dash, color))
	}
}

// libText draws a text of a symbol: its box moved with the symbol, the text
// kept readable (horizontal, or vertical reading upwards) inside it.
func (p *schPage) libText(si *symInst, it *node) {
	e := parseEffects(it.child("effects"))
	if e.hide || it.flag("private") {
		return
	}
	at := it.child("at")
	angle := at.num(2) / 10 // library text angles are in tenths of a degree
	s := p.symbolVars(si, it.arg(0))
	p.readableText(s, e, at.num(0), -at.num(1), angle, si, dimIf(schOutline, si.dnp), p.c.sch.line)
}

// readableText draws a text given in a symbol's library space (y down)
// where KiCad 10 plots it (SCH_TEXT::Plot): at its position moved by the
// symbol's transform, horizontal or reading upwards, its justification
// flipped where the transform reverses it.
func (p *schPage) readableText(s string, e effects, lx, ly, angle float64, si *symInst, color bdf.Color, defPen float64) {
	t := si.t
	origHoriz := keepUpright(angle) == 0
	screenHoriz := (t.x1 != 0) != !origHoriz
	var flipH bool
	switch {
	case origHoriz && screenHoriz:
		flipH = t.x1 < 0
	case origHoriz:
		flipH = t.x2 > 0
	case screenHoriz:
		flipH = t.y1 > 0
	default:
		flipH = t.y2 < 0
	}
	if flipH {
		e.h = -e.h
	}
	// a mirror would reverse the order of the lines
	if det := t.x1*t.y2 - t.x2*t.y1; det < 0 && origHoriz == (t.x1 > 0) {
		e.v = -e.v
	}
	drawAngle := 0.0
	if !screenHoriz {
		drawAngle = 90
	}
	e.mirror = false
	if e.hasColor {
		color = dimIf(e.color, si.dnp)
	}
	wx, wy := t.apply(lx, ly)
	p.c.text(p.d, placedText{s: s, x: si.x + wx, y: si.y + wy, angle: drawAngle, e: e, color: color, thickness: penWidth(e, defPen), brk: cad.BreakBox})
}

// pinDir is the direction from a pin's connection point to its root on the
// body, by its orientation.
func pinDir(angle float64) (float64, float64) {
	switch (int(math.Round(angle))%360 + 360) % 360 {
	case 90:
		return 0, -1
	case 180:
		return -1, 0
	case 270:
		return 0, 1
	}
	return 1, 0
}

type pinOrient int

const (
	pinRight pinOrient = iota
	pinLeft
	pinUp
	pinDown
)

// pin draws a pin: its line and decoration, and its name and number.
func (p *schPage) pin(si *symInst, it *node) {
	if it.flag("hide") {
		return
	}
	at := it.child("at")
	length := it.numOf("length", 2.54)
	cx, cy := si.t.apply(at.num(0), -at.num(1))
	conn := cad.Point{X: si.x + cx, Y: si.y + cy}
	ldx, ldy := pinDir(at.num(2))
	wdx, wdy := si.t.apply(ldx, ldy)
	wdx, wdy = math.Round(wdx), math.Round(wdy)
	root := cad.Point{X: conn.X + wdx*length, Y: conn.Y + wdy*length}
	orient := pinRight
	switch {
	case wdx < 0:
		orient = pinLeft
	case wdy < 0:
		orient = pinUp
	case wdy > 0:
		orient = pinDown
	}
	color := dimIf(schPin, si.dnp)
	pn := p.pen(p.c.sch.line, "solid", color)
	nameE := parseEffects(it.child("name").child("effects"))
	numE := parseEffects(it.child("number").child("effects"))
	// decorations as KiCad plots them (SCH_PIN::PlotPinType): of the
	// project's size, else those outside the body by the number's size and
	// those inside by the name's
	ext, in := p.c.sch.pinSymbol, p.c.sch.pinSymbol
	if ext <= 0 {
		ext, in = p.c.trunc(numE.sizeY/2), p.c.trunc(nameE.sizeY/2)
		if nameE.sizeY == 0 {
			in = ext
		}
	}
	// m points from the root back to the connection point
	m := cad.Point{X: -wdx, Y: -wdy}
	x1, y1 := root.X, root.Y
	pt := func(x, y float64) cad.Point { return cad.Point{X: x, Y: y} }
	line := func(pts ...cad.Point) {
		p.d.Stroke((&cad.Path{}).Polyline(pts, false), pn)
	}
	shape := it.arg(1)
	switch shape {
	case "inverted", "inverted_clock":
		p.d.Stroke((&cad.Path{}).Circle(pt(x1+m.X*ext, y1+m.Y*ext), ext), pn)
		line(pt(x1+m.X*ext*2, y1+m.Y*ext*2), conn)
	case "edge_clock_high":
		// the falling edge: a notch outside the body
		if m.Y == 0 {
			line(pt(x1, y1+in), pt(x1+m.X*in*2, y1), pt(x1, y1-in))
		} else {
			line(pt(x1+in, y1), pt(x1, y1+m.Y*in*2), pt(x1-in, y1))
		}
		line(pt(x1+m.X*in*2, y1+m.Y*in*2), conn)
	default:
		line(root, conn)
	}
	switch shape {
	case "clock", "inverted_clock", "clock_low":
		if m.Y == 0 {
			line(pt(x1, y1+in), pt(x1-m.X*in*2, y1), pt(x1, y1-in))
		} else {
			line(pt(x1+in, y1), pt(x1, y1-m.Y*in*2), pt(x1-in, y1))
		}
	}
	switch shape {
	case "input_low", "clock_low":
		if m.Y == 0 {
			line(pt(x1+m.X*ext*2, y1), pt(x1+m.X*ext*2, y1-ext*2), root)
		} else {
			line(pt(x1, y1+m.Y*ext*2), pt(x1-ext*2, y1+m.Y*ext*2), root)
		}
	case "output_low":
		if m.Y == 0 {
			line(pt(x1, y1-ext*2), pt(x1+m.X*ext*2, y1))
		} else {
			line(pt(x1-ext*2, y1), pt(x1, y1+m.Y*ext*2))
		}
	case "non_logic":
		line(pt(x1-(m.X+m.Y)*ext, y1-(m.Y-m.X)*ext), pt(x1+(m.X+m.Y)*ext, y1+(m.Y-m.X)*ext))
		line(pt(x1-(m.X-m.Y)*ext, y1-(m.Y+m.X)*ext), pt(x1+(m.X-m.Y)*ext, y1+(m.Y+m.X)*ext))
	}
	if it.arg(0) == "no_connect" {
		t := schTargetPinR
		line(pt(conn.X-t, conn.Y-t), pt(conn.X+t, conn.Y+t))
		line(pt(conn.X+t, conn.Y-t), pt(conn.X-t, conn.Y+t))
	}
	p.pinTexts(si, it, conn, root, orient, nameE, numE)
}

// pinTexts draws a pin's name and number where KiCad 10 plots them
// (SCH_PIN::PlotPinTexts): in the default font, with the default pen.
func (p *schPage) pinTexts(si *symInst, it *node, conn, root cad.Point, orient pinOrient, nameE, numE effects) {
	name := it.child("name").arg(0)
	if alt, ok := si.alternates[it.child("number").arg(0)]; ok {
		name = alt
	}
	if name == "~" {
		name = ""
	}
	number := it.child("number").arg(0)
	if number == "~" {
		number = ""
	}
	// pin names and numbers hold no spaces (SCH_PIN::SetName)
	name = strings.ReplaceAll(name, " ", "_")
	number = strings.ReplaceAll(number, " ", "_")
	showName := name != "" && nameE.sizeY > 0 && !si.hideNames
	showNumber := number != "" && numE.sizeY > 0 && !si.hideNumbers
	if !showName && !showNumber {
		return
	}
	if strings.HasPrefix(number, "[") && strings.HasSuffix(number, "]") && strings.Contains(number, "\n") {
		// a stack of pins, one number a line
		number = number[1 : len(number)-1]
	}
	pen := p.c.sch.line
	off := p.c.sch.pinTextOffset() + schPinTextMargin + pen
	vertical := orient == pinUp || orient == pinDown
	draw := func(s string, size float64, at cad.Point, h hAlign, v vAlign, color bdf.Color) {
		e := effects{sizeX: size, sizeY: size, lineSpacing: 1, h: h, v: v}
		angle := 0.0
		if vertical {
			angle = 90
		}
		p.c.text(p.d, placedText{s: p.expand(s), x: at.X, y: at.Y, angle: angle, e: e, color: dimIf(color, si.dnp), thickness: pen, brk: cad.BreakBox})
	}
	// the middle of the pin, as KiCad halves integers
	mid := cad.Point{X: p.c.half(root.X + conn.X), Y: p.c.half(root.Y + conn.Y)}
	// beside the pin: above a horizontal one, left of a vertical one
	beside := func(d float64) cad.Point {
		if vertical {
			return cad.Point{X: root.X - d, Y: mid.Y}
		}
		return cad.Point{X: mid.X, Y: root.Y - d}
	}
	if inside := si.nameOffset; inside > 0 {
		if showName {
			at, h := root, alignLeft
			switch orient {
			case pinRight:
				at.X += inside
			case pinLeft:
				at.X, h = at.X-inside, alignRight
			case pinDown:
				at.Y, h = at.Y+inside, alignRight
			case pinUp:
				at.Y -= inside
			}
			draw(name, nameE.sizeY, at, h, alignMiddle, schPinName)
		}
		if showNumber {
			draw(number, numE.sizeY, beside(off), alignCenter, alignBottom, schPinNumber)
		}
		return
	}
	switch {
	case showName && showNumber:
		draw(name, nameE.sizeY, beside(off), alignCenter, alignBottom, schPinName)
		draw(number, numE.sizeY, beside(-off), alignCenter, alignTop, schPinNumber)
	case showName:
		draw(name, nameE.sizeY, beside(off), alignCenter, alignBottom, schPinName)
	default:
		draw(number, numE.sizeY, beside(off), alignCenter, alignBottom, schPinNumber)
	}
}

// half halves a sum of coordinates as KiCad's integer division does.
func (c *conv) half(v float64) float64 {
	return math.Trunc(math.Round(v*c.iu)/2) / c.iu
}

// symbolVars resolves the fields of a symbol in its texts: ${REFERENCE},
// ${VALUE} and the others by name.
func (p *schPage) symbolVars(si *symInst, s string) string {
	return p.c.expand(s, func(name string) (string, bool) {
		switch strings.ToUpper(name) {
		case "REFERENCE":
			return si.shownRef(), true
		case "VALUE":
			return si.fields["Value"], true
		case "FOOTPRINT":
			return si.fields["Footprint"], true
		case "DATASHEET":
			return si.fields["Datasheet"], true
		case "DESCRIPTION":
			return si.fields["Description"], true
		case "UNIT":
			return unitLetter(si.unit), true
		case "DNP":
			if si.dnp {
				return "DNP", true
			}
			return "", true
		}
		if v, ok := si.fields[name]; ok {
			return v, true
		}
		return p.titleVars()(name)
	})
}

// unitLetter is the letter of a unit: A, B, … Z, AA, AB, …
func unitLetter(u int) string {
	if u <= 0 {
		return ""
	}
	s := ""
	for u > 0 {
		u--
		s = string(rune('A'+u%26)) + s
		u /= 26
	}
	return s
}

// shownRef is the reference shown: with the unit's letter when the symbol
// has several units.
func (si *symInst) shownRef() string {
	if si.units > 1 {
		return si.ref + unitLetter(si.unit)
	}
	return si.ref
}

// symbolFields draws the visible fields of a symbol, centered on their
// boxes as KiCad draws them.
func (p *schPage) symbolFields(si *symInst) {
	for _, f := range si.n.children("property") {
		e := parseEffects(f.child("effects"))
		if e.hide || f.flag("hide") || f.flag("private") {
			continue
		}
		name, value := f.arg(0), f.arg(1)
		if name == "Reference" {
			value = si.shownRef()
		}
		if value == "~" {
			value = ""
		}
		s := p.symbolVars(si, value)
		if strings.TrimSpace(s) == "" {
			continue
		}
		if f.flag("show_name") {
			s = name + ": " + s
		}
		color := schFields
		switch name {
		case "Reference":
			color = schReference
		case "Value":
			color = schValue
		}
		if e.hasColor {
			color = e.color
		}
		color = dimIf(color, si.dnp)
		at := f.child("at")
		fx, fy := at.xy()
		angle := at.num(2)
		penW := penWidth(e, p.c.sch.line)
		// the field's position in the symbol's own space
		ix, iy := si.t.invApply(fx-si.x, fy-si.y)
		x0, y0, w, h := p.c.textBox(s, e)
		var box cad.Rect
		for _, q := range [][2]float64{{x0, y0}, {x0 + w, y0 + h}} {
			rx, ry := rotatePt(q[0], q[1], angle)
			wx, wy := si.t.apply(ix+rx, iy+ry)
			box = box.Add(cad.Point{X: si.x + wx, Y: si.y + wy})
		}
		drawAngle := 0.0
		if int(math.Round(angle))%180 != 0 {
			drawAngle = 90
		}
		if si.t.y1 != 0 {
			drawAngle = 90 - drawAngle
		}
		ce := e
		ce.h, ce.v, ce.mirror = alignCenter, alignMiddle, false
		p.c.text(p.d, placedText{s: s, x: (box.Min.X + box.Max.X) / 2, y: (box.Min.Y + box.Max.Y) / 2, angle: drawAngle,
			e: ce, color: color, thickness: penW, brk: cad.BreakBox})
	}
}

// invApply is the inverse of apply (the transforms are orthogonal).
func (t transform) invApply(x, y float64) (float64, float64) {
	det := t.x1*t.y2 - t.y1*t.x2
	if det == 0 {
		return x, y
	}
	return (t.y2*x - t.y1*y) / det, (-t.x2*x + t.x1*y) / det
}

// dnpMark crosses out a symbol that is not placed.
func (p *schPage) dnpMark(si *symInst) {
	var body, pins cad.Rect
	for _, it := range si.items() {
		switch it.name {
		case "pin":
			if it.flag("hide") {
				continue
			}
			at := it.child("at")
			length := it.numOf("length", 2.54)
			dx, dy := pinDir(at.num(2))
			pins = pins.Add(si.point(at.num(0), at.num(1)))
			pins = pins.Add(si.point(at.num(0)+dx*length, at.num(1)-dy*length))
		case "text":
		default:
			if path, _ := shapePath(it, func(q cad.Point) cad.Point { return si.point(q.X, q.Y) }); path != nil {
				body = body.Union(path.Bounds())
			}
		}
	}
	if !body.Valid() {
		body = pins
	}
	if !body.Valid() {
		return
	}
	all := body.Union(pins)
	mx := math.Max(body.Min.X-all.Min.X, all.Max.X-body.Max.X)
	my := math.Max(body.Min.Y-all.Min.Y, all.Max.Y-body.Max.Y)
	mx, my = math.Max(mx*0.6, my*0.3), math.Max(my*0.6, mx*0.3)
	x0, y0, x1, y1 := body.Min.X-mx, body.Min.Y-my, body.Max.X+mx, body.Max.Y+my
	pn := p.pen(schDNPWidth, "solid", schDNPMarker)
	p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x0, Y: y0}, {X: x1, Y: y1}}, false), pn)
	p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{{X: x1, Y: y0}, {X: x0, Y: y1}}, false), pn)
}

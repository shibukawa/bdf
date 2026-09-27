package kicad

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
)

// sheetBox draws a hierarchical sheet: its fill (background), or its
// outline and pins, and a link to its page.
func (p *schPage) sheetBox(sh *node, background bool) {
	x, y := sh.child("at").xy()
	sz := sh.child("size")
	w, h := sz.num(0), sz.num(1)
	if w <= 0 || h <= 0 {
		return
	}
	rect := (&cad.Path{}).Polyline([]cad.Point{{X: x, Y: y}, {X: x + w, Y: y}, {X: x + w, Y: y + h}, {X: x, Y: y + h}}, true)
	if background {
		if f := sh.child("fill"); f != nil {
			if c, ok := parseColor(f.child("color")); ok {
				p.d.Fill(rect, cad.Fill{Color: c}, false)
			}
		}
		return
	}
	st := parseStroke(sh.child("stroke"))
	color := schSheet
	if st.hasColor {
		color = st.color
	}
	width := st.width
	if width <= 0 {
		width = p.c.sch.line
	}
	p.d.Stroke(rect, p.pen(width, st.dash, color))
	for _, pin := range sh.children("pin") {
		p.sheetPin(pin)
	}
	// the sheet's page: the instance of this sheet under the page's own
	uuid := sh.str("uuid")
	for _, s := range p.all {
		if s.parent == p.s && s.node != nil && s.node.str("uuid") == uuid {
			if n, ok := p.pageOf[s]; ok {
				p.links = append(p.links, pageLink{x: x, y: y, w: w, h: h, page: n})
			}
		}
	}
}

// sheetPin draws a pin of a sheet as a hierarchical label inside the box.
func (p *schPage) sheetPin(pin *node) {
	e := parseEffects(pin.child("effects"))
	if e.hide {
		return
	}
	at := pin.child("at")
	x, y := at.xy()
	var sp spin
	switch (int(math.Round(at.num(2)))%360 + 360) % 360 {
	case 0: // right edge
		sp = spinLeft
	case 90: // top edge
		sp = spinBottom
	case 270: // bottom edge
		sp = spinUp
	default: // left edge
		sp = spinRight
	}
	shape := pin.arg(1)
	switch shape {
	case "input":
		shape = "output"
	case "output":
		shape = "input"
	}
	angle, h := spinText(sp)
	e.h, e.v, e.mirror = h, alignMiddle, false
	penW := penWidth(e, p.c.sch.line)
	color := schSheetLabel
	if e.hasColor {
		color = e.color
	}
	p.hierShape(x, y, e.sizeY, shape, sp, penW, color)
	dist := p.c.round(p.c.sch.textOffset*e.sizeY) + e.sizeX
	var off cad.Point
	switch sp {
	case spinLeft:
		off.X = -dist
	case spinUp:
		off.Y = -dist
	case spinRight:
		off.X = dist
	case spinBottom:
		off.Y = dist
	}
	p.c.text(p.d, placedText{s: p.expand(pin.arg(0)), x: x + off.X, y: y + off.Y, angle: angle, e: e, color: color, thickness: penW, brk: cad.BreakBox})
}

// sheetFields draws the name and file name of a sheet, and its other
// visible fields.
func (p *schPage) sheetFields(sh *node) {
	for _, f := range sh.children("property") {
		e := parseEffects(f.child("effects"))
		if e.hide || f.flag("hide") {
			continue
		}
		name, value := f.arg(0), f.arg(1)
		s := p.expand(value)
		color := schSheetFields
		switch name {
		case "Sheetname", "Sheet name":
			color = schSheetName
		case "Sheetfile", "Sheet file":
			color = schSheetFileName
			if !f.flag("show_name") {
				s = "File: " + s
			}
		}
		if strings.TrimSpace(s) == "" {
			continue
		}
		if f.flag("show_name") {
			s = name + ": " + s
		}
		if e.hasColor {
			color = e.color
		}
		// plotted where it is, as it is justified (SCH_FIELD::Plot)
		at := f.child("at")
		x, y := at.xy()
		e.mirror = false
		p.c.text(p.d, placedText{s: s, x: x, y: y, angle: keepUpright(at.num(2)), e: e, color: color,
			thickness: penWidth(e, p.c.sch.line), brk: cad.BreakBox})
	}
}

// textBox draws a text box of the sheet, or of a symbol (si not nil): its
// frame and its text wrapped inside the margins.
func (p *schPage) textBox(it *node, si *symInst) {
	e := parseEffects(it.child("effects"))
	var x0, y0, x1, y1 float64
	if s, en := it.child("start"), it.child("end"); s != nil && en != nil {
		x0, y0, x1, y1 = s.num(0), s.num(1), en.num(0), en.num(1)
	} else {
		x, y := it.child("at").xy()
		sz := it.child("size")
		x0, y0, x1, y1 = x, y, x+sz.num(0), y+sz.num(1)
		if si != nil {
			// library text boxes grow down in the file's y-up space
			y1 = y - sz.num(1)
		}
	}
	tr := func(q cad.Point) cad.Point { return q }
	if si != nil {
		tr = func(q cad.Point) cad.Point { return si.point(q.X, q.Y) }
	}
	a, b := tr(cad.Point{X: x0, Y: y0}), tr(cad.Point{X: x1, Y: y1})
	left, right := math.Min(a.X, b.X), math.Max(a.X, b.X)
	top, bottom := math.Min(a.Y, b.Y), math.Max(a.Y, b.Y)
	rect := (&cad.Path{}).Polyline([]cad.Point{{X: left, Y: top}, {X: right, Y: top}, {X: right, Y: bottom}, {X: left, Y: bottom}}, true)
	dnp := si != nil && si.dnp
	outline := schNote
	if si != nil {
		outline = schOutline
	}
	st := parseStroke(it.child("stroke"))
	color := outline
	if st.hasColor {
		color = st.color
	}
	p.fillShape(rect, parseFill(it.child("fill")), color, dnp, st.width)
	if st.width >= 0 && st.dash != "none" {
		w := st.width
		if w <= 0 {
			w = p.c.sch.line
		}
		p.d.Stroke(rect, p.pen(w, st.dash, dimIf(color, dnp)))
	}
	var s string
	if si != nil {
		s = p.symbolVars(si, it.arg(0))
	} else {
		s = p.expand(it.arg(0))
	}
	p.boxText(s, e, left, top, right, bottom, it.child("margins"), dimIf(ifColor(e, outline), dnp))
}

func ifColor(e effects, def bdf.Color) bdf.Color {
	if e.hasColor {
		return e.color
	}
	return def
}

// boxText draws the text of a text box or table cell, its variables
// expanded, inside its margins, wrapped at spaces to their width.
func (p *schPage) boxText(s string, e effects, left, top, right, bottom float64, margins *node, color bdf.Color) {
	if strings.TrimSpace(s) == "" {
		return
	}
	ml, mt, mr, mb := e.sizeY*0.75, e.sizeY*0.75, e.sizeY*0.75, e.sizeY*0.75
	if margins != nil {
		ml, mt, mr, mb = margins.num(0), margins.num(1), margins.num(2), margins.num(3)
	}
	penW := penWidth(e, p.c.sch.line)
	s = p.wrap(s, e, right-left-ml-mr, penW)
	var x, y float64
	switch e.h {
	case alignLeft:
		x = left + ml
	case alignRight:
		x = right - mr
	default:
		x = (left + right) / 2
	}
	switch e.v {
	case alignTop:
		y = top + mt
	case alignBottom:
		y = bottom - mb
	default:
		y = (top + bottom) / 2
	}
	e.mirror = false
	p.c.text(p.d, placedText{s: s, x: x, y: y, e: e, color: color, thickness: penW, brk: cad.BreakBox})
}

// wrap breaks the lines of a text at spaces so that each fits the width
// (FONT::LinebreakText).
func (p *schPage) wrap(s string, e effects, width, penW float64) string {
	if width <= 0 || e.face != "" {
		return s
	}
	l := &strokeLayout{sizeX: e.sizeX, sizeY: e.sizeY, italic: e.italic, pen: penW, iu: p.c.iu}
	space := l.line(" ").advance
	var out []string
	for _, line := range splitLines(s) {
		words := strings.Split(line, " ")
		cur, curW := "", 0.0
		for i, w := range words {
			ww := l.line(w).advance
			if i == 0 {
				cur, curW = w, ww
				continue
			}
			if curW+space+ww < width-penW {
				cur += " " + w
				curW += space + ww
				continue
			}
			out = append(out, strings.TrimRight(cur, " "))
			cur, curW = w, ww
		}
		out = append(out, strings.TrimRight(cur, " "))
	}
	return strings.Join(out, "\n")
}

// table draws a table: its cells' fills and texts, then its borders and
// separators.
func (p *schPage) table(t *node) {
	border := t.child("border")
	sep := t.child("separators")
	borderSt := parseStroke(border.child("stroke"))
	sepSt := parseStroke(sep.child("stroke"))
	var all cad.Rect
	var cells []cad.Rect
	for _, cell := range t.child("cells").children("table_cell") {
		x, y := cell.child("at").xy()
		sz := cell.child("size")
		r := cad.Rect{}.Add(cad.Point{X: x, Y: y}).Add(cad.Point{X: x + sz.num(0), Y: y + sz.num(1)})
		if sz.num(0) <= 0 && sz.num(1) <= 0 {
			// covered by a merged cell
			continue
		}
		all = all.Union(r)
		cells = append(cells, r)
		rect := (&cad.Path{}).Polyline([]cad.Point{r.Min, {X: r.Max.X, Y: r.Min.Y}, r.Max, {X: r.Min.X, Y: r.Max.Y}}, true)
		p.fillShape(rect, parseFill(cell.child("fill")), schNote, false, 0)
		e := parseEffects(cell.child("effects"))
		p.boxText(p.expand(cell.arg(0)), e, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, cell.child("margins"), ifColor(e, schNote))
	}
	if !all.Valid() {
		return
	}
	strokeOf := func(st stroke) cad.Pen {
		w := st.width
		if w <= 0 {
			w = p.c.sch.line
		}
		c := schNote
		if st.hasColor {
			c = st.color
		}
		return p.pen(w, st.dash, c)
	}
	if sep == nil || sep.flagOr("rows", true) || sep.flagOr("cols", true) {
		sp := strokeOf(sepSt)
		rows, cols := sep == nil || sep.flagOr("rows", true), sep == nil || sep.flagOr("cols", true)
		for _, r := range cells {
			if rows && r.Min.Y > all.Min.Y+1e-6 {
				p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{r.Min, {X: r.Max.X, Y: r.Min.Y}}, false), sp)
			}
			if cols && r.Min.X > all.Min.X+1e-6 {
				p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{r.Min, {X: r.Min.X, Y: r.Max.Y}}, false), sp)
			}
		}
	}
	if border == nil || border.flagOr("external", true) {
		p.d.Stroke((&cad.Path{}).Polyline([]cad.Point{all.Min, {X: all.Max.X, Y: all.Min.Y}, all.Max, {X: all.Min.X, Y: all.Max.Y}}, true), strokeOf(borderSt))
	}
}

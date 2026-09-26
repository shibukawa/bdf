package wordproc

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/linebreak"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// emitter writes ops into one object (a page layer, a strip of the scroll
// view), keeping the drawing state it has set and the structure it has
// opened.
type emitter struct {
	c        *converter
	cv       *canvas.Canvas
	pg       *page // the page drawn (nil in the scroll view)
	view     *viewInfo
	font     bdf.FontRef
	size     float64
	hasFont  bool
	color    bdf.Color
	hasColor bool
	ls       float64
	stack    []*frame
	justOpen bool // a structure mark that starts a paragraph was written, and no text since

	// vertical draws lines of East Asian vertical text: each line in a child
	// object (inLine) where the upright characters are turned back one by
	// one, and an ALT_TEXT with the line's text before its USE_AT, so that
	// extraction sees one run per line. Inline objects and links go in the
	// parent (links collects those of a child), which extraction walks.
	vertical bool
	inLine   bool
	links    *[]linkRect
}

// linkRect is a link area of a child object, for its parent.
type linkRect struct {
	x, y, w, h float64
	url        string
}

// viewInfo resolves what depends on the whole view: links to bookmarks.
type viewInfo struct {
	link     func(bookmark string) string
	pageNum  func(bookmark string) (int, bool)
	numPages int
	secPages map[int]int
	paged    bool
}

func (e *emitter) invalidate() {
	e.hasFont, e.hasColor = false, false
	e.ls = 0
}

// op draws an op, opening and closing the structures its text is in.
func (e *emitter) op(o *op) {
	if o.top {
		e.close(0)
	}
	if o.text {
		e.enter(o.ctx, o.leaf, o.level)
	}
	o.fn(e, o.dx, o.dy)
}

// finish closes the structures still open.
func (e *emitter) finish() { e.close(0) }

// close ends the open structures after the first n.
func (e *emitter) close(n int) {
	for len(e.stack) > n {
		f := e.stack[len(e.stack)-1]
		e.stack = e.stack[:len(e.stack)-1]
		switch f.kind {
		case fTable, fList, fFigure:
			e.cv.Obj.Mark(bdf.MarkEnd, "")
		}
	}
}

// enter moves into the structure ctx and writes the mark that starts the
// text: a paragraph or heading, a line of a paragraph.
func (e *emitter) enter(ctx []*frame, leaf byte, level int) {
	n := 0
	for n < len(e.stack) && n < len(ctx) && e.stack[n] == ctx[n] {
		n++
	}
	e.close(n)
	opened := false
	for _, f := range ctx[n:] {
		switch f.kind {
		case fTable:
			e.cv.Obj.Mark(bdf.MarkTable, "")
		case fCell:
			e.cv.Obj.Mark(bdf.MarkCell, f.payload)
			opened = true
		case fList:
			e.cv.Obj.Mark(bdf.MarkList, "")
		case fItem:
			e.cv.Obj.Mark(bdf.MarkListItem, "")
			opened = true
		case fBox:
			e.cv.Obj.Mark(bdf.MarkBox, "")
			opened = true
		case fFigure:
			e.cv.Obj.Mark(bdf.MarkFigure, f.payload)
		}
		e.stack = append(e.stack, f)
	}
	opened = opened || e.justOpen
	e.justOpen = opened && leaf == 0
	switch leaf {
	case 'H':
		e.cv.Obj.Mark(bdf.MarkHeading, strconv.Itoa(level))
	case 'P':
		// a cell, list item or text box starts a paragraph already
		if !opened {
			e.cv.Obj.Mark(bdf.MarkParagraph, "")
		}
	case 'L':
		e.cv.Obj.Mark(bdf.MarkLine, "")
	case 'W':
		e.cv.Obj.Mark(bdf.MarkWrap, "")
	}
}

func (e *emitter) setFont(fc *fontset.Choice, size float64) {
	ref := e.cv.Font(fc.Use)
	if !e.hasFont || ref != e.font || size != e.size {
		e.cv.Obj.Font(ref, f32(size))
		e.font, e.size, e.hasFont = ref, size, true
	}
}

func (e *emitter) setColor(c bdf.Color) {
	if !e.hasColor || c != e.color {
		e.cv.Obj.FillColor(c)
		e.color, e.hasColor = c, true
	}
}

func (e *emitter) setLetterSpacing(ls float64) {
	if ls != e.ls {
		e.cv.Obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, f32(ls))
		e.ls = ls
	}
}

func (e *emitter) fillRect(c bdf.Color, x, y, w, h float64) {
	e.setColor(c)
	e.cv.Obj.FillRect(f32(x), f32(y), f32(w), f32(h))
	e.cv.Drawn = true
}

// rule draws a horizontal or vertical border line from (x0,y0) to (x1,y1),
// centered on it.
func (e *emitter) rule(b *border, x0, y0, x1, y1 float64) {
	if b.none() {
		return
	}
	w := b.width
	lines := []float64{0}
	if b.style == "double" {
		lines = []float64{-w, w}
	}
	for _, off := range lines {
		ax0, ay0, ax1, ay1 := x0, y0, x1, y1
		if y0 == y1 {
			ay0, ay1 = y0+off, y1+off
		} else {
			ax0, ax1 = x0+off, x1+off
		}
		switch b.style {
		case "dotted", "dashed", "dashSmallGap", "dotDash", "dotDotDash":
			dash := []float32{f32(w), f32(2 * w)}
			if b.style != "dotted" {
				dash = []float32{f32(4 * w), f32(2 * w)}
			}
			e.cv.Obj.Save()
			e.cv.Obj.StrokeColor(b.color)
			e.cv.Obj.Line(f32(w), bdf.CapButt, bdf.JoinMiter, 10)
			e.cv.Obj.Dash(dash, 0)
			p := &bdf.Path{}
			p.MoveTo(f32(ax0), f32(ay0)).LineTo(f32(ax1), f32(ay1))
			e.cv.Obj.StrokePath(e.cv.Obj.AddPath(p))
			e.cv.Obj.Restore()
			e.invalidate()
		default:
			e.setColor(b.color)
			if ay0 == ay1 {
				e.cv.Obj.FillRect(f32(ax0-w/2), f32(ay0-w/2), f32(ax1-ax0+w), f32(w))
			} else {
				e.cv.Obj.FillRect(f32(ax0-w/2), f32(ay0-w/2), f32(w), f32(ay1-ay0+w))
			}
		}
	}
	e.cv.Drawn = true
}

func (e *emitter) cellBorders(bs [4]*border, x, y, w, h float64) {
	e.rule(bs[sTop], x, y, x+w, y)
	e.rule(bs[sBottom], x, y+h, x+w, y+h)
	e.rule(bs[sLeft], x, y, x, y+h)
	e.rule(bs[sRight], x+w, y, x+w, y+h)
}

// line draws a laid-out line whose area starts at x and whose top is at y.
func (e *emitter) line(ln *line, x, y float64) {
	if e.vertical {
		e.verticalLine(ln, x, y)
		return
	}
	base := y + ln.baseline
	if len(ln.label) > 0 {
		e.runs(ln, ln.label, x, base)
		e.cv.Obj.Mark(bdf.MarkLine, "")
	}
	e.runs(ln, ln.items, x, base)
}

// runs draws items as runs of text of one style, with their tabs, inline
// objects and fields.
func (e *emitter) runs(ln *line, items []item, x, base float64) {
	var link struct {
		url            string
		x0, x1, y0, y1 float64
	}
	flushLink := func() {
		if link.url != "" && link.x1 > link.x0 {
			url := link.url
			if strings.HasPrefix(url, "#bm:") {
				url = ""
				if e.view != nil && e.view.link != nil {
					url = e.view.link(strings.TrimPrefix(link.url, "#bm:"))
				}
			}
			switch {
			case url == "":
			case e.links != nil:
				*e.links = append(*e.links, linkRect{link.x0, link.y0, link.x1 - link.x0, link.y1 - link.y0, url})
			default:
				e.cv.Obj.Link(f32(link.x0), f32(link.y0), f32(link.x1-link.x0), f32(link.y1-link.y0), url)
			}
		}
		link.url = ""
	}
	addLink := func(it *item, x0, x1 float64) {
		url := ""
		if it.st != nil {
			url = it.st.link
		}
		if url != link.url {
			flushLink()
			if url == "" {
				return
			}
			link.url, link.x0, link.x1 = url, x0, x1
			link.y0, link.y1 = base-ln.asc, base+ln.desc
			return
		}
		link.x0, link.x1 = math.Min(link.x0, x0), math.Max(link.x1, x1)
	}
	for i := 0; i < len(items); {
		it := &items[i]
		switch it.kind {
		case kTab:
			e.leader(it, x, base)
			i++
			continue
		case kObject:
			if e.inLine {
				i++ // drawn by the parent
				continue
			}
			o := it.obj
			ox, oy := x+it.x+o.ext[0], base-o.h-o.ext[3]
			e.drawObject(o, drawingml.Box{X: ox, Y: oy, W: o.w, H: o.h})
			addLink(it, ox, ox+o.w)
			i++
			continue
		case kField:
			text := e.fieldText(it.fld)
			if text != "" {
				st := it.st
				e.cv.SetLang(runLang(st, text))
				e.text(text, it, x+it.x, base-it.shift, it.w, 0, st)
				addLink(it, x+it.x, x+it.x+it.w)
			}
			i++
			continue
		case kChar:
		default:
			i++
			continue
		}
		if it.ls < 0 || (lbSpace(it.r) && i == len(items)-1) {
			i++
			continue
		}
		if e.inLine && it.tcy != 0 {
			j := i + 1
			for j < len(items) && items[j].tcy == it.tcy {
				j++
			}
			e.horizontalInVertical(items[i:j], x, base-it.shift)
			addLink(it, x+it.x, x+it.x+it.w+it.pad)
			i = j
			continue
		}
		// a run: characters of the same style, face and spacing that follow
		// each other without a gap
		j := i + 1
		for j < len(items) {
			n := &items[j]
			p := &items[j-1]
			if n.kind != kChar || n.ls < 0 || n.st != it.st && n.st.key != it.st.key || n.fc != it.fc || n.ls+n.pad != it.ls+it.pad ||
				n.size != it.size || p.gap != 0 || math.Abs(p.x+p.w+max(p.ls, 0)+p.pad-n.x) > 0.01 ||
				e.inLine && fontset.Upright(n.r) != fontset.Upright(it.r) {
				break
			}
			j++
		}
		k := j
		if j == len(items) || items[j].kind == kBreak || items[j].kind == kPage || items[j].kind == kColumn {
			for k > i && lbSpace(items[k-1].r) {
				k--
			}
		}
		if k > i {
			var b strings.Builder
			adv := 0.0
			for _, r := range items[i:k] {
				b.WriteRune(r.r)
				adv += r.w + max(r.ls, 0) + r.pad
			}
			text := b.String()
			if e.inLine && fontset.Upright(it.r) {
				e.upright(items[i:k], x, base-it.shift, adv)
			} else {
				e.cv.SetLang(runLang(it.st, text))
				e.text(text, it, x+it.x, base-it.shift, adv, max(it.ls, 0)+it.pad, it.st)
			}
			addLink(it, x+it.x, x+it.x+adv)
		}
		i = j
	}
	flushLink()
}

// runLang returns the language of a run: the East Asian language for
// East Asian text, the Latin one otherwise.
func runLang(st *runStyle, text string) string {
	for _, r := range text {
		if fontdb.IsCJK(r) {
			if st.eaLang != "" {
				return st.eaLang
			}
			if l := fontset.RuneLang(r); l != "" {
				return l
			}
			return ""
		}
	}
	return st.lang
}

// text draws a run of text at (x, y) with its decorations; ls is the
// spacing added after each character (justification, the character grid).
func (e *emitter) text(text string, it *item, x, y, adv, ls float64, st *runStyle) {
	rp := st.rp
	size := it.size
	fc := it.fc
	e.background(it, x, y, adv)
	e.setFont(fc, size)
	e.setColor(st.color)
	e.setLetterSpacing(rp.spacing + math.Max(ls, 0))
	if rp.scale != 1 && rp.scale > 0 {
		e.cv.Obj.Save()
		e.cv.Obj.Transform(f32(rp.scale), 0, 0, 1, f32(x), f32(y))
		e.cv.Obj.FillText(text, 0, 0, f32(adv/rp.scale))
		e.cv.Obj.Restore()
		e.invalidate()
	} else {
		e.cv.Obj.FillText(text, f32(x), f32(y), f32(adv))
	}
	e.cv.Drawn = true
	e.decorate(text, it, x, y, adv)
}

// background draws the highlight or shading of a run behind its text.
func (e *emitter) background(it *item, x, y, adv float64) {
	rp := it.st.rp
	asc, desc := it.fc.Asc*it.size, it.fc.Desc*it.size
	if rp.highlight != "" && rp.highlight != "none" {
		if c, ok := namedColors[strings.ToLower(rp.highlight)]; ok {
			e.fillRect(c, x, y-asc, adv, asc+desc)
		}
	} else if rp.hasShd {
		e.fillRect(rp.shd, x, y-asc, adv, asc+desc)
	}
}

// decorate draws the underline, strikethrough and emphasis marks of a run.
func (e *emitter) decorate(text string, it *item, x, y, adv float64) {
	st, fc, size := it.st, it.fc, it.size
	rp := st.rp
	th := math.Max(fc.ULThick*size, 0.5)
	lineColor := st.color
	if rp.hasUColor {
		lineColor = rp.uColor
	}
	if u := rp.u; u != "none" && u != "" {
		uy := y - fc.ULPos*size
		if e.inLine {
			// vertical text is underlined on its right
			uy = y - fc.Asc*size - 2*th
		}
		w := adv - rp.spacing
		switch u {
		case "double", "wavyDouble":
			e.fillRect(lineColor, x, uy, w, th)
			e.fillRect(lineColor, x, uy+2*th, w, th)
		case "thick", "dottedHeavy", "dashedHeavy", "dashLongHeavy", "dashDotHeavy", "dashDotDotHeavy", "wavyHeavy":
			e.fillRect(lineColor, x, uy, w, 2*th)
		case "dotted", "dash", "dashLong", "dotDash", "dotDotDash", "wave":
			dash := []float32{f32(th), f32(2 * th)}
			if u != "dotted" {
				dash = []float32{f32(4 * th), f32(2 * th)}
			}
			e.cv.Obj.Save()
			e.cv.Obj.StrokeColor(lineColor)
			e.cv.Obj.Line(f32(th), bdf.CapButt, bdf.JoinMiter, 10)
			e.cv.Obj.Dash(dash, 0)
			p := &bdf.Path{}
			p.MoveTo(f32(x), f32(uy+th/2)).LineTo(f32(x+w), f32(uy+th/2))
			e.cv.Obj.StrokePath(e.cv.Obj.AddPath(p))
			e.cv.Obj.Restore()
			e.invalidate()
		default:
			e.fillRect(lineColor, x, uy, w, th)
		}
	}
	if rp.has(tStrike) || rp.has(tDStrike) {
		sy := y - 0.3*size
		e.fillRect(st.color, x, sy, adv, th)
		if rp.has(tDStrike) {
			e.fillRect(st.color, x, sy-2*th, adv, th)
		}
	}
	if rp.em != "" && rp.em != "none" {
		e.emphasis(text, it, x, y, adv, st)
	}
}

// upright draws East Asian characters of vertical text standing upright:
// each character's em box is turned back by 90° around its center, in the
// middle of its cell (with the letter spacing of justification and the
// character grid). Punctuation takes its vertical presentation form, and
// small kana move to the top right of their squares.
func (e *emitter) upright(run []item, x, y, adv float64) {
	first := run[0]
	e.background(&first, x+first.x, y, adv)
	e.setColor(first.st.color)
	e.setLetterSpacing(0)
	size := first.size
	half := (first.fc.Asc - first.fc.Desc) * size / 2
	for _, r := range run {
		cell := r.w + max(r.ls, 0) + r.pad
		cx := x + r.x + cell/2
		glyph, gfc := r.r, r.fc
		dx, dy := 0.0, 0.0
		if v, ok := fontset.VerticalForm(r.r); ok {
			if vf := e.c.face(r.st, v); vf.Loaded != nil && vf.Loaded.Has(v) {
				glyph, gfc = v, vf
				e.c.fonts.Advance(vf, v)
			} else {
				dx, dy = 0.55*size, -0.55*size
			}
		} else if smallKana(r.r) {
			dx, dy = 0.1*size, -0.1*size
		}
		e.setFont(gfc, size)
		e.cv.Obj.Save()
		e.cv.Obj.Transform(0, -1, 1, 0, f32(cx), f32(y-half))
		e.cv.Obj.FillText(string(glyph), f32(-r.w/2+dx), f32(half+dy), f32(r.w-r.st.rp.spacing))
		e.cv.Obj.Restore()
	}
	e.cv.Drawn = true
	var b strings.Builder
	for _, r := range run {
		b.WriteRune(r.r)
	}
	e.decorate(b.String(), &first, x+first.x, y, adv)
}

// horizontalInVertical draws characters set horizontally in one square of
// vertical text (tate-chu-yoko), upright and centered in the square,
// compressed to the width of the line when the run asks for it.
func (e *emitter) horizontalInVertical(run []item, x, y float64) {
	first := run[0]
	st, fc, size := first.st, first.fc, first.size
	var b strings.Builder
	width := 0.0
	for _, r := range run {
		b.WriteRune(r.r)
		width += e.c.fonts.Advance(r.fc, r.r) * size * st.rp.scale
	}
	cell := first.w + first.pad + max(first.ls, 0)
	e.background(&first, x+first.x, y, cell)
	e.setFont(fc, size)
	e.setColor(st.color)
	e.setLetterSpacing(0)
	half := (fc.Asc - fc.Desc) * size / 2
	sx := 1.0
	if st.rp.tcyFit && width > size {
		sx = size / width
	}
	e.cv.Obj.Save()
	// upright (turned back by 90°), then compressed along the text
	e.cv.Obj.Transform(0, -f32(sx), 1, 0, f32(x+first.x+cell/2), f32(y-half))
	e.cv.Obj.FillText(b.String(), f32(-width/2), f32(half), f32(width))
	e.cv.Obj.Restore()
	e.cv.Drawn = true
}

// smallKana reports whether r is a small kana, which vertical text sets at
// the top right of its square.
func smallKana(r rune) bool {
	return strings.ContainsRune("ぁぃぅぇぉっゃゅょゎゕゖァィゥェォッャュョヮヵヶ", r)
}

// verticalLine draws a line of vertical text (see emitter.vertical) whose
// area starts at x and whose top is at y, in the turned coordinates.
func (e *emitter) verticalLine(ln *line, x, y float64) {
	base := y + ln.baseline
	if len(ln.label) > 0 {
		if e.verticalRun(ln, ln.label, x, base) {
			e.cv.Obj.Mark(bdf.MarkLine, "")
		}
	}
	e.verticalRun(ln, ln.items, x, base)
	// inline objects stand upright in their cells
	for i := range ln.items {
		it := &ln.items[i]
		if it.kind != kObject {
			continue
		}
		o := it.obj
		cx := x + it.x + (o.ext[1]+o.h+o.ext[3])/2
		cy := base - (o.ext[0]+o.w+o.ext[2])/2
		e.cv.Obj.Save()
		e.cv.Obj.Transform(0, -1, 1, 0, f32(cx), f32(cy))
		e.drawObject(o, drawingml.Box{X: -o.w / 2, Y: -o.h / 2, W: o.w, H: o.h})
		e.cv.Obj.Restore()
		e.invalidate()
	}
}

// verticalRun draws the text of items in a child object with the line's
// text as its ALT_TEXT; false when there is no text.
func (e *emitter) verticalRun(ln *line, items []item, x, base float64) bool {
	var b strings.Builder
	first, last := -1, -1
	for i, it := range items {
		switch it.kind {
		case kChar:
			b.WriteRune(it.r)
		case kTab:
			b.WriteRune(' ')
		case kField:
			b.WriteString(e.fieldText(it.fld))
		default:
			continue
		}
		if first < 0 {
			first = i
		}
		last = i
	}
	text := strings.TrimRight(b.String(), " \u3000")
	if first < 0 || text == "" {
		return false
	}
	x0 := x + items[first].x
	end := &items[last]
	width := end.x + end.w + end.pad + max(end.ls, 0) - items[first].x
	size := ln.asc + ln.desc
	// The bbox spans the line along its length, which is the extent of the
	// ALT_TEXT run (spec §7.8); across it there is room for turned glyphs,
	// emphasis marks and underlines.
	ch, ref := e.cv.Child(bdf.Rect{X: 0, Y: f32(-ln.asc), W: f32(width), H: f32(size)})
	ch.Obj.SetBBox(0, f32(-ln.asc-size), f32(width), f32(3*size))
	var links []linkRect
	che := &emitter{c: e.c, cv: ch, pg: e.pg, view: e.view, inLine: true, links: &links}
	che.runs(ln, items, x-x0, 0)
	// the font the text is measured with for search highlights
	e.setFont(items[first].fc, items[first].size)
	e.cv.SetLang(runLang(items[first].st, text))
	e.cv.Obj.Mark(bdf.MarkAltText, text)
	e.cv.Obj.UseAt(ref, f32(x0), f32(base))
	e.cv.Drawn = true
	for _, l := range links {
		e.cv.Obj.Link(f32(l.x+x0), f32(l.y+base), f32(l.w), f32(l.h), l.url)
	}
	return true
}

// emphasis draws East Asian emphasis marks (bouten) over the characters.
func (e *emitter) emphasis(text string, it *item, x, y, adv float64, st *runStyle) {
	size := it.size
	n := 0
	for range text {
		n++
	}
	if n == 0 {
		return
	}
	step := adv / float64(n)
	r := size * 0.08
	cy := y - it.fc.Asc*size - r*1.5
	if st.rp.em == "underDot" {
		cy = y + it.fc.Desc*size + r*1.5
	}
	p := &bdf.Path{}
	i := 0
	for _, ch := range text {
		if !linebreak.IsSpace(ch) {
			cx := x + step*(float64(i)+0.5)
			p.Circle(f32(cx), f32(cy), f32(r))
		}
		i++
	}
	e.setColor(st.color)
	e.cv.Obj.FillPath(e.cv.Obj.AddPath(p), 0)
}

// leader draws the leader of a tab (dots, a line) across its width.
func (e *emitter) leader(it *item, x, base float64) {
	if it.r == '\t' || it.w <= 0 || it.st == nil || it.fc == nil {
		return
	}
	size := it.size
	x0, x1 := x+it.x, x+it.x+it.w
	e.setColor(it.st.color)
	switch it.r {
	case '_':
		th := math.Max(it.fc.ULThick*size, 0.5)
		e.fillRect(it.st.color, x0, base-it.fc.ULPos*size, x1-x0, th)
	case '-':
		e.fillRect(it.st.color, x0+size*0.1, base-0.3*size, x1-x0-size*0.2, math.Max(size*0.06, 0.5))
	default:
		// dots on a grid of the font's period width
		step := e.c.fonts.Advance(it.fc, '.') * size
		if step <= 0 {
			step = size * 0.28
		}
		r := size * 0.055
		p := &bdf.Path{}
		cy := base - r
		if it.r == '·' {
			cy = base - 0.3*size
		}
		for cx := math.Ceil((x0+step*0.5)/step) * step; cx+r < x1-step*0.3; cx += step {
			p.Circle(f32(cx), f32(cy), f32(r))
		}
		e.cv.Obj.FillPath(e.cv.Obj.AddPath(p), 0)
	}
	e.cv.Drawn = true
}

// fieldText returns the text of a page field on the page drawn.
func (e *emitter) fieldText(f *field) string {
	if e.pg == nil || e.view == nil || !e.view.paged {
		if f.typ == "PAGEREF" && e.view != nil && e.view.pageNum != nil {
			if _, ok := e.view.pageNum(f.arg); !ok {
				return f.cached
			}
		}
		return f.cached
	}
	n := 0
	fmtName := e.pg.sec.pgFmt
	switch f.typ {
	case "PAGE":
		n = e.pg.num
	case "NUMPAGES":
		n = e.view.numPages
		fmtName = "decimal"
	case "SECTIONPAGES":
		n = e.view.secPages[e.pg.secNo]
		fmtName = "decimal"
	case "PAGEREF":
		p, ok := e.view.pageNum(f.arg)
		if !ok {
			return f.cached
		}
		n = p
	}
	switch f.format {
	case "roman":
		return strings.ToLower(roman(n))
	case "ROMAN", "Roman":
		return roman(n)
	case "alphabetic":
		return strings.ToLower(letters(n))
	case "ALPHABETIC", "Alphabetic":
		return letters(n)
	case "Arabic", "arabic":
		return strconv.Itoa(n)
	}
	return formatNumber(n, fmtName)
}

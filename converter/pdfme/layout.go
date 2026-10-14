package pdfme

import (
	"math"
	"sort"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/formbox"
)

// Tables and lists, and the layout of a blank base: what grows with the
// data (a table with its rows, a list with its items, a text that makes
// its field taller) moves what is under it on its page, and goes on over
// further pages. The rules are those of pdfme.

// tableCell is a cell of a table that is laid out.
type tableCell struct {
	text string
	ts   *textStyle
}

// tableRow is a row of a table, with its height in millimetres.
type tableRow struct {
	cells []tableCell
	h     float64
	head  bool
}

// tableLayout is a table that is laid out: its columns and rows.
type tableLayout struct {
	cols     []float64 // mm
	head     *tableRow // nil for a table without the names of its columns
	body     []tableRow
	showHead bool
}

// cellStyle is the style of the cells of a part of a table: what the
// template says of them over what cells are without it.
func cellStyle(styles map[string]any) *textStyle {
	ts := &textStyle{size: 10, lineHeight: 1, align: "left", valign: "middle", color: "#000000", borderColor: "#000000", pad: sides{5, 5, 5, 5}}
	if styles == nil {
		return ts
	}
	s := &schema{m: styles}
	ts.font = s.str("fontName", ts.font)
	ts.align, ts.valign = s.str("alignment", ts.align), s.str("verticalAlignment", ts.valign)
	ts.size = clamp(s.num("fontSize", ts.size), 0, 2000)
	ts.lineHeight = clamp(s.num("lineHeight", ts.lineHeight), 0, 100)
	ts.spacing = clamp(s.num("characterSpacing", 0), -100, 1000)
	ts.bg, ts.color = s.str("backgroundColor", ""), s.str("fontColor", ts.color)
	ts.borderColor = s.str("borderColor", ts.borderColor)
	if v, ok := styles["borderWidth"]; ok {
		ts.border = sidesOf(v, 0)
	}
	if v, ok := styles["padding"]; ok {
		ts.pad = sidesOf(v, 0)
	}
	return ts
}

// shows reports whether a table shows the names of its columns.
func (s *schema) showsHead() bool { return s.m["showHead"] != false }

// tableOf lays a table out: the rows from and to of its body (to < 0:
// all), each as tall as the lines of its cells.
func (f *filler) tableOf(s *schema, v any, from, to int, split bool) *tableLayout {
	heads := s.heads()
	body := f.rowsOf(s, v)
	if to < 0 || to > len(body) {
		to = len(body)
	}
	from = min(max(from, 0), to)
	body = body[from:to]
	n := len(heads)
	if n == 0 && len(body) > 0 {
		n = len(body[0])
	}
	tl := &tableLayout{showHead: s.showsHead() && (!split || s.flag("repeatHead")) && len(heads) > 0}
	if n == 0 {
		return tl
	}
	// the columns: the shares of the width that the template gives them,
	// and what is left in equal shares for those without one
	tl.cols = make([]float64, n)
	pct, _ := s.m["headWidthPercentages"].([]any)
	left, without := s.w, 0
	for i := range tl.cols {
		if i < len(pct) {
			if p, ok := pct[i].(float64); ok && p > 0 && p <= 100 {
				tl.cols[i] = s.w * p / 100
				left -= tl.cols[i]
				continue
			}
		}
		without++
	}
	for i := range tl.cols {
		if tl.cols[i] == 0 && without > 0 {
			tl.cols[i] = math.Max(left, 0) / float64(without)
		}
	}
	headStyle, bodyStyles := s.obj("headStyles"), s.obj("bodyStyles")
	alt, _ := bodyStyles["alternateBackgroundColor"].(string)
	plain, _ := bodyStyles["backgroundColor"].(string)
	if from%2 == 1 {
		alt, plain = plain, alt // (a part that starts at an odd row keeps its stripes)
	}
	columns := s.obj("columnStyles")
	aligns, _ := columns["alignment"].(map[string]any)
	valigns, _ := columns["verticalAlignment"].(map[string]any)
	row := func(cells []string, styles map[string]any, bg string, head bool) tableRow {
		r := tableRow{head: head, cells: make([]tableCell, n)}
		for c := range r.cells {
			ts := cellStyle(styles)
			if !head {
				ts.bg = bg
			}
			key := strconv.Itoa(c)
			if a, ok := aligns[key].(string); ok {
				ts.align = a
			}
			if a, ok := valigns[key].(string); ok {
				ts.valign = a
			}
			text := ""
			if c < len(cells) {
				text = cells[c]
			}
			r.cells[c] = tableCell{text, ts}
			st := &formbox.Style{Size: ts.size, Spacing: ts.spacing}
			st.Font, st.Bold, st.Italic = fontOf(ts.font)
			cw := math.Max(0, tl.cols[c]-ts.border[1]-ts.border[3]-ts.pad[1]-ts.pad[3]) * mm
			lines := len(f.wrap(text, st, cw))
			h := float64(lines)*ts.size*pt2mm*ts.lineHeight + ts.border[0] + ts.border[2] + ts.pad[0] + ts.pad[2]
			r.h = math.Max(r.h, h)
		}
		return r
	}
	if len(heads) > 0 {
		h := row(heads, headStyle, "", true)
		tl.head = &h
	}
	tl.body = make([]tableRow, len(body))
	for i, cells := range body {
		bg := plain
		if i%2 == 0 {
			bg = alt
		}
		tl.body[i] = row(cells, bodyStyles, bg, false)
	}
	return tl
}

// cellRef names a cell of a table ("B2").
func cellRef(row, col int) string {
	name := ""
	for n := col + 1; n > 0; n = (n - 1) / 26 {
		name = string(rune('A'+(n-1)%26)) + name
	}
	return name + strconv.Itoa(row+1)
}

// table draws a table: its cells with their backgrounds, borders and
// texts, and the border of the table over them.
func (f *filler) table(l *formbox.Layer, s *schema, v any) {
	from, to := 0, -1
	if s.part {
		from, to = s.from, s.to
	}
	tl := f.tableOf(s, v, from, to, s.split)
	var rows []*tableRow
	if tl.showHead && tl.head != nil {
		rows = append(rows, tl.head)
	}
	for i := range tl.body {
		rows = append(rows, &tl.body[i])
	}
	if len(rows) == 0 {
		return
	}
	l.Mark(bdf.MarkTable, "")
	y := s.y
	for ri, r := range rows {
		x := s.x
		for ci := range r.cells {
			cell, w := &r.cells[ci], tl.cols[ci]
			ts := cell.ts
			b := box{x * mm, y * mm, w * mm, r.h * mm}
			if c, ok := colorOf(ts.bg); ok && b.w > 0 && b.h > 0 {
				l.Rect(b.x, b.y, b.w, b.h, c)
			}
			if c, ok := colorOf(ts.borderColor); ok {
				rect := func(x, y, w, h float64) {
					if w > 0 && h > 0 {
						l.Rect(x, y, w, h, c)
					}
				}
				bt, br, bb, bl := ts.border[0]*mm, ts.border[1]*mm, ts.border[2]*mm, ts.border[3]*mm
				rect(b.x, b.y, b.w, bt)
				rect(b.x+b.w-br, b.y, br, b.h)
				rect(b.x, b.y+b.h-bb, b.w, bb)
				rect(b.x, b.y, bl, b.h)
			}
			ref := cellRef(ri, ci)
			if r.head {
				ref += " col"
			}
			l.Mark(bdf.MarkCell, ref)
			inner := *ts
			inner.bg, inner.borderColor, inner.border, inner.pad = "", "", sides{}, sides{}
			f.text(l, cell.text, box{
				b.x + (ts.border[3]+ts.pad[3])*mm, b.y + (ts.border[0]+ts.pad[0])*mm,
				math.Max(0, w-ts.border[1]-ts.border[3]-ts.pad[1]-ts.pad[3]) * mm,
				math.Max(0, r.h-ts.border[0]-ts.border[2]-ts.pad[0]-ts.pad[2]) * mm,
			}, &inner, 0, -1)
			x += w
		}
		y += r.h
	}
	l.Mark(bdf.MarkEnd, "")
	styles := s.obj("tableStyles")
	bw, _ := styles["borderWidth"].(float64)
	color, _ := styles["borderColor"].(string)
	if bw > 0 && bw < 1e3 {
		f.shape(l, false, box{s.x * mm, s.y * mm, s.w * mm, (y - s.y) * mm}, bw*mm, 0, "", color, false)
	}
}

// listItem is an item of a list: its text, and how deep it is (by the
// tabs before it).
type listItem struct {
	text  string
	level int
}

const maxListLevel = 8

func parseListItem(raw string) listItem {
	level := 0
	for level < len(raw) && raw[level] == '\t' {
		level++
	}
	return listItem{raw[level:], min(level, maxListLevel)}
}

// listLayout is a list that is laid out: its items with their markers and
// heights (millimetres).
type listLayout struct {
	items   []listItem
	markers []string
	heights []float64
	// the width of a marker, what is between it and its item, and how far
	// a level is in
	markerW, gap, indent float64
}

// listOf lays out the items from and to of a list (to < 0: all).
func (f *filler) listOf(s *schema, v any, from, to int) *listLayout {
	all := f.itemsOf(s, v)
	if to < 0 || to > len(all) {
		to = len(all)
	}
	from = min(max(from, 0), to)
	ll := &listLayout{markerW: clamp(s.num("markerWidth", 6), 0, 1e3), gap: clamp(s.num("markerGap", 2), 0, 1e3), indent: clamp(s.num("indentSize", 6), 0, 1e3)}
	// the markers count the items of the whole list
	ordered := s.str("listStyle", "bullet") == "ordered"
	var counters [maxListLevel + 1]int
	markers := make([]string, len(all))
	for i, raw := range all {
		markers[i] = "•"
		if ordered {
			level := parseListItem(raw).level
			counters[level]++
			for j := level + 1; j <= maxListLevel; j++ {
				counters[j] = 0
			}
			markers[i] = strconv.Itoa(counters[level]) + "."
		}
	}
	ts := s.textStyle()
	st := &formbox.Style{Size: ts.size, Spacing: ts.spacing}
	st.Font, st.Bold, st.Italic = fontOf(ts.font)
	lineH := ts.size * ts.lineHeight * pt2mm
	spacing := clamp(s.num("itemSpacing", 1), 0, 1e3)
	for i, raw := range all[from:to] {
		it := parseListItem(raw)
		bodyW := math.Max(s.w-float64(it.level)*ll.indent-ll.markerW-ll.gap, 0)
		lines := len(f.wrap(it.text, st, math.Max(bodyW, 0.1)*mm))
		h := float64(max(lines, 1)) * lineH
		if i != to-from-1 {
			h += spacing
		}
		ll.items, ll.markers, ll.heights = append(ll.items, it), append(ll.markers, markers[from+i]), append(ll.heights, h)
	}
	return ll
}

// list draws a list: each item with its marker before it.
func (f *filler) list(l *formbox.Layer, s *schema, v any) {
	from, to := 0, -1
	if s.part {
		from, to = s.from, s.to
	}
	ll := f.listOf(s, v, from, to)
	if len(ll.items) == 0 {
		return
	}
	ts := s.textStyle()
	if c, ok := colorOf(ts.bg); ok && s.w > 0 && s.h > 0 {
		b := s.box()
		l.Rect(b.x, b.y, b.w, b.h, c)
	}
	ts.bg, ts.dynamic, ts.valign = "", nil, "top"
	l.Mark(bdf.MarkList, "")
	y := s.y
	for i, it := range ll.items {
		markerX := float64(it.level) * ll.indent
		bodyX := markerX + ll.markerW + ll.gap
		l.Mark(bdf.MarkListItem, "")
		marker := *ts
		marker.align = "right"
		f.text(l, ll.markers[i], box{(s.x + markerX) * mm, y * mm, ll.markerW * mm, ll.heights[i] * mm}, &marker, 0, -1)
		f.text(l, it.text, box{(s.x + bodyX) * mm, y * mm, math.Max(s.w-bodyX, 0) * mm, ll.heights[i] * mm}, ts, 0, -1)
		y += ll.heights[i]
	}
	l.Mark(bdf.MarkEnd, "")
}

// expands reports whether a text makes its field taller when it has more
// lines than the field has room for (on a blank base alone).
func (f *filler) expands(s *schema) bool {
	return f.t.blank != nil && (s.typ == "text" || s.typ == "multiVariableText") && s.str("overflow", "") == "expand"
}

// units are the parts that a field is as tall as, in millimetres: the
// rows of a table, the items of a list, the lines of a text that makes its
// field taller. A field of another kind is one part, as tall as it is.
type units struct {
	heights []float64
	// avoidFirst keeps the first part from being alone at the foot of a
	// page (the names of the columns of a table)
	avoidFirst bool
	// patch makes the field of the parts start to end of a field
	patch func(s *schema, start, end int)
}

// unitsOf measures the parts of a field with its value.
func (f *filler) unitsOf(s *schema) units {
	whole := units{heights: []float64{s.h}}
	v := f.content(s, false)
	switch {
	case s.typ == "table":
		tl := f.tableOf(s, v, 0, -1, false)
		u := units{heights: f.tableUnits(s, tl), avoidFirst: true}
		u.patch = func(c *schema, start, end int) {
			// (the first part is the names of the columns)
			c.from, c.to = max(start-1, 0), end-1
		}
		return u
	case s.typ == "list":
		ll := f.listOf(s, v, 0, -1)
		if len(ll.items) == 0 {
			return units{heights: []float64{0}}
		}
		return units{heights: ll.heights, patch: func(c *schema, start, end int) { c.from, c.to = start, end }}
	case f.expands(s):
		text := f.textOf(v)
		if s.typ == "multiVariableText" {
			text = f.variableText(s, v)
		}
		if text == "" {
			return whole
		}
		ts := s.textStyle()
		st := &formbox.Style{Size: ts.size, Spacing: ts.spacing}
		st.Font, st.Bold, st.Italic = fontOf(ts.font)
		cw := math.Max(0, s.w-ts.border[1]-ts.border[3]-ts.pad[1]-ts.pad[3]) * mm
		lines := f.wrap(text, st, cw)
		heights := f.lineHeights(len(lines), st, ts.lineHeight)
		total := 0.0
		for i := range heights {
			heights[i] *= pt2mm
			if i == 0 {
				heights[i] += ts.border[0] + ts.pad[0]
			}
			if i == len(heights)-1 {
				heights[i] += ts.border[2] + ts.pad[2]
			}
			total += heights[i]
		}
		n := len(heights)
		switch {
		case total <= s.h || n == 0:
			return units{heights: []float64{s.h}, patch: func(c *schema, start, end int) { c.total = 1 }}
		case n == 1:
			return units{heights: []float64{math.Max(s.h, total)}, patch: func(c *schema, start, end int) { c.total = 1 }}
		}
		return units{heights: heights, patch: func(c *schema, start, end int) { c.from, c.to, c.total = start, end, n }}
	}
	return whole
}

// tableUnits returns the heights of the parts of a table: the names of
// its columns and its rows. For a table that says the names again on each
// page, the first row of a page is as tall as the names and the row.
func (f *filler) tableUnits(s *schema, tl *tableLayout) []float64 {
	headH := 0.0
	base := []float64{0}
	if s.showsHead() && tl.head != nil {
		headH = tl.head.h
		base[0] = headH
	}
	for i := range tl.body {
		base = append(base, tl.body[i].h)
	}
	bp := f.t.blank
	if !s.flag("repeatHead") || bp == nil || headH <= 0 {
		return base
	}
	padTop := bp.pad[0]
	content := bp.h - bp.pad[0] - bp.pad[2]
	if content <= 0 {
		return base
	}
	pageStart := func(page int) float64 { return float64(page)*content + padTop }
	initial := int(math.Max(0, math.Floor((s.y-padTop)/content)))
	const headRows, safety, eps = 1, 0.5, 0.01
	page, y := initial, s.y
	onPage, travels := 0, false
	var out []float64
	for i, rowH := range base {
		body := i >= headRows
		for guard := 0; guard < 8; guard++ {
			start := pageStart(page)
			remaining := start + content - y
			needsHead := body && onPage == 0 && page > initial && !travels
			total := rowH
			if needsHead {
				total += headH
			}
			if total > remaining-safety+eps {
				if onPage == 0 && math.Abs(y-start) < safety {
					out = append(out, total)
					y += total
					onPage++
					travels = false
					break
				}
				alone := !travels && page == initial && body && len(out) == headRows && onPage == headRows
				page++
				if alone {
					travels, y, onPage = true, pageStart(page)+headH, headRows
				} else {
					travels, y, onPage = false, pageStart(page), 0
				}
				continue
			}
			out = append(out, total)
			y += total
			onPage++
			if body {
				travels = false
			}
			if y >= start+content-safety {
				page++
				y, onPage, travels = pageStart(page), 0, false
			}
			break
		}
		if len(out) <= i {
			out = append(out, rowH)
		}
	}
	return out
}

// layout places the fields of a blank base on its pages: as the template
// has them when nothing grows, and when something does, each page of the
// template as the pages that its fields take.
func (f *filler) layout() [][]*schema {
	t := f.t
	grows := false
	for _, page := range t.pages {
		for _, s := range page {
			grows = grows || s.typ == "table" || s.typ == "list" || f.expands(s)
		}
	}
	bp := t.blank
	content, padTop := bp.h-bp.pad[0]-bp.pad[2], bp.pad[0]
	if !grows || content <= 0 {
		return t.pages
	}
	type item struct {
		s     *schema
		baseY float64
		order int
		u     units
	}
	const eps = 0.01
	var result [][]*schema
	for _, list := range t.pages {
		if len(list) == 0 {
			result = append(result, nil)
			continue
		}
		items := make([]item, len(list))
		for i, s := range list {
			items[i] = item{s: s, baseY: math.Max(0, s.y-padTop), order: i, u: f.unitsOf(s)}
		}
		sort.SliceStable(items, func(a, b int) bool {
			if math.Abs(items[a].baseY-items[b].baseY) > eps {
				return items[a].baseY < items[b].baseY
			}
			return items[a].order < items[b].order
		})
		type placedSchema struct {
			s     *schema
			order int
		}
		var pages [][]placedSchema
		offset := 0.0
		for _, it := range items {
			heights := it.u.heights
			startY := it.baseY + offset
			unit := 0
			page := int(math.Max(0, math.Floor(startY/content)))
			y := math.Mod(startY, content)
			if y < 0 {
				y = 0
			}
			end := 0.0
			splittable := len(heights) > 1
			for unit < len(heights) {
				if page >= maxPages*2 {
					break
				}
				for len(pages) <= page {
					pages = append(pages, nil)
				}
				left := content - y
				if heights[unit] > left+eps && math.Abs(left-content) > eps {
					page++
					y = 0
					continue
				}
				chunk := 0.0
				first := unit
				for unit < len(heights) && y+chunk+heights[unit] <= content+eps {
					chunk += heights[unit]
					unit++
				}
				if it.u.avoidFirst && splittable && first == 0 && unit == 1 && y > eps {
					// (the names of the columns are not left alone at the foot of a page)
					unit = 0
					if len(pages[page]) > 0 {
						page++
					}
					y = 0
					continue
				}
				if unit == first {
					chunk += heights[unit]
					unit++
				}
				c := *it.s
				c.y, c.h, c.origin = y+padTop, chunk, it.s
				if it.u.patch != nil {
					c.part, c.split = true, first > 0
					it.u.patch(&c, first, unit)
				}
				pages[page] = append(pages[page], placedSchema{&c, it.order})
				y += chunk
				if y >= content-eps {
					page++
					y = 0
				}
				end = float64(page)*content + y
			}
			offset = end - (it.baseY + it.s.h)
		}
		for len(pages) > 1 && len(pages[len(pages)-1]) == 0 {
			pages = pages[:len(pages)-1]
		}
		for _, p := range pages {
			sort.SliceStable(p, func(a, b int) bool { return p[a].order < p[b].order })
			out := make([]*schema, len(p))
			for i := range p {
				out[i] = p[i].s
			}
			result = append(result, out)
			if len(result) > maxPages {
				return result
			}
		}
	}
	return result
}

package pdfme

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/formbox"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/barcode"
)

// The drawing of the fields, as the generator of pdfme draws them: the
// places and the sizes of the template are in millimetres, the sizes of
// text in points, and a page is drawn in points from its top left.

// pt2mm is the millimetres of a point, as pdfme has them.
const pt2mm = 0.3528

// maxLines bounds the lines of a text that are drawn.
const maxLines = 1 << 14

// box is a place on a page, in points.
type box struct{ x, y, w, h float64 }

func (s *schema) box() box { return box{s.x * mm, s.y * mm, s.w * mm, s.h * mm} }

// sides are the widths at the top, the right, the bottom and the left of a
// box, in millimetres.
type sides [4]float64

func sidesOf(v any, def float64) sides {
	ok := func(n float64) bool { return n > 0 && n < 1e4 }
	switch x := v.(type) {
	case float64:
		if ok(x) {
			return sides{x, x, x, x}
		}
		return sides{}
	case map[string]any:
		out := sides{}
		for i, k := range []string{"top", "right", "bottom", "left"} {
			if n, isNum := x[k].(float64); isNum && ok(n) {
				out[i] = n
			}
		}
		return out
	}
	return sides{def, def, def, def}
}

func (s sides) any() bool { return s[0] > 0 || s[1] > 0 || s[2] > 0 || s[3] > 0 }

// colorOf reads a color of a template ("#RGB", "#RGBA", "#RRGGBB",
// "#RRGGBBAA", with or without the #). ok is false for none, and for one
// that is clear.
func colorOf(s string) (c bdf.Color, ok bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	var v [8]byte
	n := len(s)
	if n != 3 && n != 4 && n != 6 && n != 8 {
		return 0, false
	}
	for i := 0; i < n; i++ {
		switch ch := s[i]; {
		case ch >= '0' && ch <= '9':
			v[i] = ch - '0'
		case ch >= 'a' && ch <= 'f':
			v[i] = ch - 'a' + 10
		case ch >= 'A' && ch <= 'F':
			v[i] = ch - 'A' + 10
		default:
			return 0, false
		}
	}
	var r, g, b, a uint8 = 0, 0, 0, 255
	if n <= 4 {
		r, g, b = v[0]*17, v[1]*17, v[2]*17
		if n == 4 {
			a = v[3] * 17
		}
	} else {
		r, g, b = v[0]<<4|v[1], v[2]<<4|v[3], v[4]<<4|v[5]
		if n == 8 {
			a = v[6]<<4 | v[7]
		}
	}
	if a == 0 {
		return 0, false
	}
	return bdf.RGBA(r, g, b, a), true
}

// fontOf reads the name of a font of a template as a family and its
// style: "NotoSansJP-Bold" is Noto Sans JP, bold. Without a name it is the
// font of pdfme, Roboto.
func fontOf(name string) (family string, bold, italic bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "Roboto", false, false
	}
	family = name
	if i := strings.LastIndexAny(name, "-_"); i > 0 {
		switch style := strings.ToLower(name[i+1:]); style {
		case "regular", "normal", "medium", "book", "light", "bold", "italic", "bolditalic", "oblique", "boldoblique", "black", "heavy", "semibold":
			family = name[:i]
			bold = strings.Contains(style, "bold") || style == "black" || style == "heavy"
			italic = strings.Contains(style, "italic") || strings.Contains(style, "oblique")
		}
	}
	return family, bold, italic
}

// textStyle is how the text of a field is drawn.
type textStyle struct {
	font                      string
	size, lineHeight, spacing float64
	align, valign             string
	color, bg, borderColor    string
	border, pad               sides
	// dynamic is the size that fits ({min, max, fit}), when the text has
	// one
	dynamic           map[string]any
	underline, strike bool
}

func clamp(v, lo, hi float64) float64 { return math.Min(math.Max(v, lo), hi) }

func (s *schema) textStyle() *textStyle {
	ts := &textStyle{
		font: s.str("fontName", ""), size: clamp(s.num("fontSize", 13), 0, 2000), lineHeight: clamp(s.num("lineHeight", 1), 0, 100),
		spacing: clamp(s.num("characterSpacing", 0), -100, 1000), align: s.str("alignment", "left"), valign: s.str("verticalAlignment", "top"),
		color: s.str("fontColor", ""), bg: s.str("backgroundColor", ""), borderColor: s.str("borderColor", ""),
		border: sidesOf(s.m["borderWidth"], 0), pad: sidesOf(s.m["padding"], 0),
		dynamic: s.obj("dynamicFontSize"), underline: s.flag("underline"), strike: s.flag("strikethrough"),
	}
	return ts
}

// line is a line of a text that was broken into lines; hard is set for the
// last of a paragraph.
type line struct {
	text string
	hard bool
}

// wrap breaks a text into the lines of a width.
func (f *filler) wrap(value string, st *formbox.Style, width float64) []line {
	value = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(value)
	var out []line
	for _, para := range strings.Split(value, "\n") {
		ls := f.r.Wrap(para, st, math.Max(width, 0.01))
		for i, t := range ls {
			out = append(out, line{t, i == len(ls)-1})
		}
		if len(out) > maxLines {
			f.warnf("a text has more than %d lines: the rest are left out", maxLines)
			return out[:maxLines]
		}
	}
	return out
}

// lineHeights returns the heights of the lines of a text in points: the
// first is as tall as its font reaches over the line, the others as the
// size of the font.
func (f *filler) lineHeights(n int, st *formbox.Style, lineHeight float64) []float64 {
	out := make([]float64, n)
	asc, _ := f.r.Metrics(st)
	for i := range out {
		out[i] = st.Size * lineHeight
		if i == 0 {
			out[i] = asc * lineHeight
		}
	}
	return out
}

// fitSize finds the size of a text that fits its box: grown and shrunk in
// steps of a quarter of a point, as pdfme does it.
func (f *filler) fitSize(value string, ts *textStyle, st formbox.Style, cw, ch float64) float64 {
	num := func(key string, def float64) float64 {
		if v, ok := ts.dynamic[key].(float64); ok && !math.IsNaN(v) {
			return v
		}
		return def
	}
	lo, hi := clamp(num("min", 4), 0, 2000), clamp(num("max", 72), 0, 2000)
	size := st.Size
	if hi < lo {
		return size
	}
	horizontal := ts.dynamic["fit"] == "horizontal"
	size = clamp(size, lo, hi)
	if utf8.RuneCountInString(value) > 20000 {
		return size // (too long a text to try sizes on)
	}
	paras := strings.Split(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(value), "\n")
	calc := func(size float64) (tw, th float64) {
		st.Size = size
		lines := f.wrap(value, &st, cw)
		for _, h := range f.lineHeights(len(lines), &st, ts.lineHeight) {
			th += h
		}
		if horizontal {
			for _, p := range paras {
				tw = math.Max(tw, f.r.Width(p, &st))
			}
		} else {
			for _, ln := range lines {
				tw = math.Max(tw, f.r.Width(ln.text, &st))
			}
		}
		return tw, th
	}
	const step = 0.25
	tw, th := calc(size)
	steps := 0
	for ; size < hi && (horizontal && tw < cw || !horizontal && th < ch) && steps < 2000; steps++ {
		size += step
		nw, nh := calc(size)
		if nh >= ch {
			size -= step
			break
		}
		tw, th = nw, nh
	}
	for ; size > lo && size > 0 && (tw > cw || th > ch) && steps < 4000; steps++ {
		size -= step
		tw, th = calc(size)
	}
	return size
}

// wide reports whether a character is one that lines break at without a
// space (those of Chinese, Japanese and Korean).
func wide(r rune) bool {
	return r >= 0x2E80 && r <= 0x9FFF || r >= 0xAC00 && r <= 0xD7AF || r >= 0xF900 && r <= 0xFAFF || r >= 0xFF00 && r <= 0xFFEF || r >= 0x20000 && r <= 0x3FFFF
}

// text draws a text in a box: its background and its border, and its
// lines where the alignment of the field puts them. from and to select
// the lines of a part of a text that goes on over pages (to < 0: all).
func (f *filler) text(l *formbox.Layer, value string, b box, ts *textStyle, from, to int) {
	if value == "" {
		return
	}
	if len(value) > maxText {
		value = strings.ToValidUTF8(value[:maxText], "")
		f.warnf("a text is longer than %d bytes: the rest is left out", maxText)
	}
	if c, ok := colorOf(ts.bg); ok && b.w > 0 && b.h > 0 {
		l.Rect(b.x, b.y, b.w, b.h, c)
	}
	if c, ok := colorOf(ts.borderColor); ok && ts.border.any() {
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
	cx, cy := b.x+(ts.border[3]+ts.pad[3])*mm, b.y+(ts.border[0]+ts.pad[0])*mm
	cw := math.Max(0, b.w-(ts.border[1]+ts.border[3]+ts.pad[1]+ts.pad[3])*mm)
	ch := math.Max(0, b.h-(ts.border[0]+ts.border[2]+ts.pad[0]+ts.pad[2])*mm)
	col, ok := colorOf(ts.color)
	if !ok {
		col = bdf.RGB(0, 0, 0)
	}
	st := &formbox.Style{Size: ts.size, Spacing: ts.spacing, Color: col}
	st.Font, st.Bold, st.Italic = fontOf(ts.font)
	if ts.dynamic != nil {
		st.Size = f.fitSize(value, ts, *st, cw, ch)
	}
	if st.Size <= 0 {
		return
	}
	lines := f.wrap(value, st, cw)
	if from > 0 || to >= 0 {
		if to < 0 || to > len(lines) {
			to = len(lines)
		}
		lines = lines[min(max(from, 0), to):to]
	}
	if len(lines) == 0 {
		return
	}
	asc, desc := f.r.Metrics(st)
	lh := ts.lineHeight
	adj := (lh - 1) * st.Size / 2
	if lh == 0 {
		adj = 0
	}
	other := lh * st.Size * float64(len(lines)-1)
	var first float64 // the line of the first line, under the top of the content
	switch ts.valign {
	case "bottom":
		first = ch - other - desc - adj
	case "middle":
		first = (ch-other-asc-desc)/2 + asc
	default:
		first = asc + adj
	}
	l.Mark(bdf.MarkParagraph, "")
	for i, ln := range lines {
		if i > 0 {
			// the lines of a paragraph join with a space, and without one
			// where the text broke between characters that have none
			kind := byte(bdf.MarkLine)
			if prev := lines[i-1]; !prev.hard && prev.text != "" && ln.text != "" {
				a, _ := utf8.DecodeLastRuneInString(prev.text)
				c, _ := utf8.DecodeRuneInString(ln.text)
				if wide(a) || wide(c) {
					kind = bdf.MarkWrap
				}
			}
			l.Mark(kind, "")
		}
		width := f.r.Width(ln.text, st)
		lst := *st
		x := 0.0
		switch ts.align {
		case "center":
			x = (cw - width) / 2
		case "right":
			x = cw - width
		case "justify":
			if n := utf8.RuneCountInString(ln.text); !ln.hard && n > 0 {
				lst.Spacing += (cw - width) / float64(n)
			}
		}
		y := cy + first + lh*st.Size*float64(i)
		if ln.text != "" {
			l.Text(ln.text, &lst, cx+x, y)
		}
		if width <= 0 {
			continue
		}
		if ts.strike {
			l.Line(cx+x, y-asc/3, cx+x+width+1, y-asc/3, st.Size/12, col)
		}
		if ts.underline {
			l.Line(cx+x, y+asc/12, cx+x+width+1, y+asc/12, st.Size/12, col)
		}
	}
}

// within draws a field turned and as clear as the template has it: turned
// clockwise about the middle of its box.
func (f *filler) within(l *formbox.Layer, s *schema, draw func()) {
	rot, alpha := s.num("rotate", 0), s.num("opacity", 1)
	if rot == 0 && alpha >= 1 {
		draw()
		return
	}
	o := l.Obj
	o.Save()
	if alpha < 1 {
		o.Alpha(float32(math.Max(alpha, 0)))
	}
	if rot != 0 {
		cx, cy := (s.x+s.w/2)*mm, (s.y+s.h/2)*mm
		sn, cs := math.Sincos(rot * math.Pi / 180)
		l.Transform([6]float64{cs, sn, -sn, cs, cx - cs*cx + sn*cy, cy - sn*cx - cs*cy})
	}
	draw()
	o.Restore()
}

// content returns what a field shows: the value that the data has for it,
// or the content that the template has, with its expressions.
func (f *filler) content(s *schema, static bool) any {
	text := s.str("content", "")
	switch {
	case s.readOnly && s.typ == "table":
		if v, ok := f.data[s.name]; ok && v != nil {
			return v
		}
		return text
	case s.readOnly && s.typ == "multiVariableText":
		return text
	case s.readOnly:
		if !strings.Contains(text, "{") {
			return text
		}
		out, failed := expand(text, f.vars)
		if failed && !f.measuring {
			f.warnf("an expression of %q is not evaluated; it is left as it is", s.name)
		}
		return out
	case static:
		return text
	}
	v, _ := f.value(s)
	return v
}

// textOf writes a value of the data as the text of a field.
func (f *filler) textOf(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case bool:
		if x {
			return "true"
		}
		return "false"
	case []any, map[string]any:
		b, _ := json.Marshal(x)
		return string(b)
	}
	return tag.ValueText(v, "", f.td.Format)
}

// variableText fills the variables of a text with variables ("{name}")
// with their values.
func (f *filler) variableText(s *schema, v any) string {
	text := s.str("text", "")
	if str, ok := v.(string); ok {
		var m map[string]any
		if json.Unmarshal([]byte(str), &m) != nil {
			// (a text that is no object of variables is the text itself)
			if s.readOnly && str != "" && len(s.variables()) > 0 {
				return str
			}
			m = nil
		}
		v, _ = tag.Normalize(m).(map[string]any)
	}
	vars, _ := v.(map[string]any)
	if len(s.variables()) == 0 && s.readOnly {
		return text
	}
	for name, val := range vars {
		text = strings.ReplaceAll(text, "{"+name+"}", f.textOf(val))
	}
	return text
}

// draw draws a field on a page. static is set for a field of the base
// page, which every page has.
func (f *filler) draw(l *formbox.Layer, s *schema, static bool) {
	if !static && s.name == "" {
		return // (pdfme draws no field without a name)
	}
	v := f.content(s, static)
	b := s.box()
	f.within(l, s, func() {
		switch s.typ {
		case "text", "select", "multiVariableText":
			ts := s.textStyle()
			from, to := 0, -1
			if s.part {
				// a part of a text that goes on: without a size that fits,
				// and with the border and the room of its ends alone
				ts.dynamic = nil
				if s.total > 1 {
					from, to = s.from, s.to
					if s.from > 0 {
						ts.border[0], ts.pad[0] = 0, 0
					}
					if s.to < s.total {
						ts.border[2], ts.pad[2] = 0, 0
					}
				}
			} else if f.expands(s) {
				ts.dynamic = nil
			}
			text := f.textOf(v)
			if s.typ == "multiVariableText" {
				text = f.variableText(s, v)
			}
			f.text(l, text, b, ts, from, to)
		case "date", "time", "dateTime":
			ts := s.textStyle()
			ts.valign, ts.lineHeight = "middle", 1
			f.text(l, f.dateText(s, v), b, ts, 0, -1)
		case "image", "signature", "svg":
			f.picture(l, s, v, b)
		case "checkbox", "radioGroup":
			f.icon(l, s, v, b)
		case "circleMark":
			f.circleMark(l, s, v, b)
		case "line":
			if c, ok := colorOf(s.str("color", "")); ok && b.w != 0 && b.h != 0 {
				l.Rect(b.x, b.y, b.w, b.h, c)
			}
		case "rectangle", "ellipse":
			f.shape(l, s.typ == "ellipse", b, s.num("borderWidth", 0)*mm, s.num("radius", 0)*mm, s.str("color", ""), s.str("borderColor", ""),
				s.num("opacity", 1) < 1)
		case "table":
			f.table(l, s, v)
		case "list":
			f.list(l, s, v)
		default:
			if _, isCode := codeKinds[s.typ]; isCode {
				f.code(l, s, v, b)
				return
			}
			f.warnf("a field of the type %q is not drawn (%q)", s.typ, s.name)
		}
	})
}

// shape draws a rectangle or an ellipse: filled, and with a border that
// is inside its box. clear is set for a shape that is not opaque: its
// border is then not over its fill (pdfme fills and draws the border at
// once, and a PDF does not paint the one over the other).
func (f *filler) shape(l *formbox.Layer, ellipse bool, b box, bw, radius float64, fill, border string, clear bool) {
	fc, filled := colorOf(fill)
	bc, bordered := colorOf(border)
	bordered = bordered && bw > 0
	if !filled && !bordered || b.w <= 0 || b.h <= 0 {
		return
	}
	if !bordered {
		bw = 0
	}
	if b.w-bw <= 0 || b.h-bw <= 0 {
		return
	}
	// the path that is in from the box by a width
	path := func(in float64) *bdf.Path {
		w, h := b.w-2*in, b.h-2*in
		if w <= 0 || h <= 0 {
			return nil
		}
		p := &bdf.Path{}
		switch r := radius - (in - bw/2); {
		case ellipse:
			p.Ellipse(float32(b.x+b.w/2), float32(b.y+b.h/2), float32(w/2), float32(h/2), 0, 0, float32(2*math.Pi), false).Close()
		case radius > 0 && r > 0:
			p.RoundRect(float32(b.x+in), float32(b.y+in), float32(w), float32(h), float32(math.Min(r, math.Min(w, h)/2)))
		default:
			p.Rect(float32(b.x+in), float32(b.y+in), float32(w), float32(h))
		}
		return p
	}
	if filled {
		in := bw / 2
		if bordered && (clear || uint8(bc) != 255) {
			in = bw
		}
		if p := path(in); p != nil {
			l.Fill(p, fc)
		}
	}
	if bordered {
		l.Stroke(path(bw/2), bw, bc, false)
	}
}

// marked reports whether a mark is on: true, or the text of it.
func marked(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		return x == "true"
	}
	return false
}

// icon draws a checkbox or a radio button: the icons of pdfme (a square
// with round corners and a check, a circle and a dot), as large as fits
// the box.
func (f *filler) icon(l *formbox.Layer, s *schema, v any, b box) {
	col, ok := colorOf(s.str("color", "#000000"))
	if !ok {
		col = bdf.RGB(0, 0, 0)
	}
	k := math.Min(b.w, b.h) / 24
	if k <= 0 {
		return
	}
	ox, oy := b.x+(b.w-24*k)/2, b.y+(b.h-24*k)/2
	at := func(x, y float64) (float32, float32) { return float32(ox + x*k), float32(oy + y*k) }
	on := marked(v)
	if s.typ == "radioGroup" {
		cx, cy := at(12, 12)
		l.Stroke((&bdf.Path{}).Circle(cx, cy, float32(10*k)), 2*k, col, true)
		if on {
			l.Stroke((&bdf.Path{}).Circle(cx, cy, float32(k)), 2*k, col, true)
		}
		return
	}
	x, y := at(3, 3)
	l.Stroke((&bdf.Path{}).RoundRect(x, y, float32(18*k), float32(18*k), float32(2*k)), 2*k, col, true)
	if on {
		x0, y0 := at(16, 9)
		x1, y1 := at(10.5, 14.5)
		x2, y2 := at(8, 12)
		l.Stroke((&bdf.Path{}).MoveTo(x0, y0).LineTo(x1, y1).LineTo(x2, y2), 2*k, col, true)
	}
}

// circleMark draws a circle around what is under it, when it is on.
func (f *filler) circleMark(l *formbox.Layer, s *schema, v any, b box) {
	bw := s.num("borderWidth", 0) * mm
	col, ok := colorOf(s.str("color", ""))
	if !marked(v) || !ok || bw <= 0 || b.w-bw <= 0 || b.h-bw <= 0 {
		return
	}
	p := (&bdf.Path{}).Ellipse(float32(b.x+b.w/2), float32(b.y+b.h/2), float32((b.w-bw)/2), float32((b.h-bw)/2), 0, 0, float32(2*math.Pi), false).Close()
	l.Stroke(p, bw, col, false)
}

// svgSize reads the size of an SVG image: of its viewBox, or its width
// and height.
func svgSize(data []byte) (w, h float64) {
	i := bytes.Index(data, []byte("<svg"))
	if i < 0 {
		return 0, 0
	}
	tagEnd := bytes.IndexByte(data[i:], '>')
	if tagEnd < 0 {
		return 0, 0
	}
	head := string(data[i : i+tagEnd])
	attr := func(name string) string {
		for rest := head; ; {
			j := strings.Index(rest, name+"=")
			if j < 0 {
				return ""
			}
			before := byte(' ')
			if j > 0 {
				before = rest[j-1]
			}
			rest = rest[j+len(name)+1:]
			if before != ' ' && before != '\t' && before != '\n' && before != '\r' || rest == "" {
				continue
			}
			if q := rest[0]; q == '"' || q == '\'' {
				if e := strings.IndexByte(rest[1:], q); e >= 0 {
					return rest[1 : 1+e]
				}
			}
			return ""
		}
	}
	num := func(s string) float64 {
		s = strings.TrimSpace(s)
		end := 0
		for end < len(s) && (s[end] >= '0' && s[end] <= '9' || s[end] == '.' || s[end] == '-' || s[end] == '+' || s[end] == 'e') {
			end++
		}
		v, _ := strconv.ParseFloat(s[:end], 64)
		return v
	}
	if vb := strings.FieldsFunc(attr("viewBox"), func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\n' }); len(vb) == 4 {
		if w, h = num(vb[2]), num(vb[3]); w > 0 && h > 0 {
			return w, h
		}
	}
	if w, h = num(attr("width")), num(attr("height")); w > 0 && h > 0 {
		return w, h
	}
	return 0, 0
}

// loaded is a picture that the data named, as it was read: the picture of
// the document (its size is 0 × 0 for an SVG image that says none), or
// why it is left out.
type loaded struct {
	pic formbox.Picture
	ok  bool
	why string
}

// load reads a picture that the data names: the path of a file, a data:
// URL, or the markup of an SVG image.
func (f *filler) load(name string) loaded {
	if ld, ok := f.pictures[name]; ok {
		return ld
	}
	var ld loaded
	var data []byte
	if t := strings.TrimSpace(name); strings.HasPrefix(t, "<") {
		data = []byte(t) // (the markup of an SVG image)
	} else {
		var err error
		if data, err = tag.PictureData(f.td.Files, name); err != nil {
			ld.why = err.Error()
		}
	}
	if ld.why == "" {
		if pic, ok := f.r.Picture(data); !ok {
			ld.why = "it is not a picture that is read"
		} else {
			ld.pic, ld.ok = pic, true
			if pic.W <= 0 || pic.H <= 0 {
				ld.pic.W, ld.pic.H = svgSize(data)
			}
		}
	}
	if f.pictures == nil {
		f.pictures = map[string]loaded{}
	}
	f.pictures[name] = ld
	return ld
}

// picture draws a picture in the box of its field: whole, or covering
// the box, as the field says.
func (f *filler) picture(l *formbox.Layer, s *schema, v any, b box) {
	name, _ := v.(string)
	if name == "" || b.w <= 0 || b.h <= 0 {
		return
	}
	ld := f.load(name)
	if !ld.ok {
		f.warnf("the picture of %q is left out: %s", s.name, ld.why)
		return
	}
	pic := ld.pic
	sw, sh := pic.W, pic.H
	if sw <= 0 || sh <= 0 {
		sw, sh = b.w, b.h
	}
	cover := s.typ == "image" && s.str("objectFit", "contain") == "cover"
	scale := math.Min(b.w/sw, b.h/sh)
	hpos, vpos := "center", "center"
	if cover {
		scale = math.Max(b.w/sw, b.h/sh)
	} else if parts := strings.Fields(s.str("objectPosition", "")); len(parts) == 2 {
		hpos, vpos = parts[0], parts[1]
	}
	w, h := sw*scale, sh*scale
	offset := func(room float64, pos string) float64 {
		switch pos {
		case "left", "top":
			return 0
		case "right", "bottom":
			return room
		}
		return room / 2
	}
	alt := ""
	if !s.readOnly {
		alt = s.name
	}
	o := l.Obj
	l.Mark(bdf.MarkFigure, alt)
	if cover {
		o.Save()
		o.ClipRect(float32(b.x), float32(b.y), float32(b.w), float32(b.h))
	}
	l.Image(pic, b.x+offset(b.w-w, hpos), b.y+offset(b.h-h, vpos), w, h)
	if cover {
		o.Restore()
	}
	l.Mark(bdf.MarkEnd, "")
}

// code draws a QR code or a barcode in the box of its field, stretched to
// it; a barcode has what it says under its bars, when the field asks.
func (f *filler) code(l *formbox.Layer, s *schema, v any, b box) {
	var text string
	switch x := v.(type) {
	case string:
		text = x
	case float64:
		text = strconv.FormatFloat(x, 'f', -1, 64)
	}
	if text == "" || b.w <= 0 || b.h <= 0 {
		return
	}
	kind := codeKinds[s.typ]
	if kind == "" {
		f.warnf("the %s code of %q is not drawn: this kind of code is not made", s.typ, s.name)
		return
	}
	// pdfme places and styles its own human-readable line below the bars.
	sym, err := barcode.Encode(kind, text, &barcode.Options{Quiet: barcode.NoQuiet, NoText: true})
	if err != nil {
		f.warnf("the code of %q is left out: %v", s.name, err)
		return
	}
	if c, ok := colorOf(s.str("backgroundColor", "")); ok {
		l.Rect(b.x, b.y, b.w, b.h, c)
	}
	col, ok := colorOf(s.str("barColor", ""))
	if !ok {
		col = bdf.RGB(0, 0, 0)
	}
	bars := b
	labelled := sym.Linear && s.m["includetext"] != false
	// what a barcode says is under its bars, in the room of a line of
	// 10 points under bars as tall as the box, all made as tall as the box
	const label, labelSize = 9.0, 10.0
	scale := 1.0
	if labelled {
		scale = b.h / (b.h + label)
		bars.h = b.h * scale
	}
	if s.typ == "itf14" {
		// the frame around the bars of ITF-14, and the room at their sides:
		// 4 and 15 of the widths of a narrow bar, as pdfme has them
		const frame, quiet = 4.0, 15.0
		fw := bars.w * frame / (float64(sym.W) + 2*(frame+quiet))
		fh := bars.h * frame / (b.h + 2*frame)
		if fw > 0 && fh > 0 && bars.w > 2*fw && bars.h > 2*fh {
			l.Rect(bars.x, bars.y, bars.w, fh, col)
			l.Rect(bars.x, bars.y+bars.h-fh, bars.w, fh, col)
			l.Rect(bars.x, bars.y, fw, bars.h, col)
			l.Rect(bars.x+bars.w-fw, bars.y, fw, bars.h, col)
			side := fw * (frame + quiet) / frame
			bars = box{bars.x + side, bars.y + fh, bars.w - 2*side, bars.h - 2*fh}
		}
	}
	l.Mark(bdf.MarkFigure, sym.Text)
	l.Code(sym, bars.x, bars.y, bars.w, bars.h, col)
	l.Mark(bdf.MarkEnd, "")
	if !labelled {
		return
	}
	tc, ok := colorOf(s.str("textColor", ""))
	if !ok {
		tc = col
	}
	text = sym.Text
	if s.typ == "code39" {
		text = "*" + text + "*" // (the marks of its start and end)
	}
	st := &formbox.Style{Font: "Helvetica", Size: labelSize * scale, Spacing: labelSize * scale * 0.15, Color: tc}
	width := f.r.Width(text, st)
	if width > b.w && width > 0 {
		st.Size *= b.w / width
		st.Spacing *= b.w / width
		width = b.w
	}
	l.Mark(bdf.MarkParagraph, "")
	l.Text(text, st, b.x+(b.w-width)/2, b.y+b.h-0.1*st.Size)
}

// Months and days of the week, for the formats of dates.
var (
	monthsEn = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	daysEn   = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	daysJa   = []string{"日", "月", "火", "水", "木", "金", "土"}
)

// formatDate writes a date in a format of date-fns ("yyyy/MM/dd",
// "MMM d, yyyy HH:mm"), in English or in Japanese.
func formatDate(t time.Time, format, locale string) string {
	ja := strings.HasPrefix(locale, "ja")
	var out strings.Builder
	pad := func(v, n int) string {
		s := strconv.Itoa(v)
		for len(s) < n {
			s = "0" + s
		}
		return s
	}
	rs := []rune(format)
	for i := 0; i < len(rs); {
		c := rs[i]
		if c == '\'' {
			// a text in quotes ('' is a quote)
			j := i + 1
			for ; j < len(rs); j++ {
				if rs[j] == '\'' {
					if j+1 < len(rs) && rs[j+1] == '\'' {
						out.WriteRune('\'')
						j++
						continue
					}
					break
				}
				out.WriteRune(rs[j])
			}
			if j == i+1 {
				out.WriteRune('\'')
			}
			i = j + 1
			continue
		}
		n := 1
		for i+n < len(rs) && rs[i+n] == c {
			n++
		}
		if c == 'd' && n == 1 && i+1 < len(rs) && rs[i+1] == 'o' {
			// the day with its ending ("1st")
			d := t.Day()
			end := "th"
			if ja {
				end = "日"
			} else if d%100 < 11 || d%100 > 13 {
				switch d % 10 {
				case 1:
					end = "st"
				case 2:
					end = "nd"
				case 3:
					end = "rd"
				}
			}
			out.WriteString(strconv.Itoa(d) + end)
			i += 2
			continue
		}
		h12 := (t.Hour()+11)%12 + 1
		switch c {
		case 'y':
			if n == 2 {
				out.WriteString(pad(t.Year()%100, 2))
			} else {
				out.WriteString(pad(t.Year(), n))
			}
		case 'M', 'L':
			switch {
			case n >= 4 && ja:
				out.WriteString(strconv.Itoa(int(t.Month())) + "月")
			case n >= 4:
				out.WriteString(monthsEn[t.Month()-1])
			case n == 3 && ja:
				out.WriteString(strconv.Itoa(int(t.Month())) + "月")
			case n == 3:
				out.WriteString(monthsEn[t.Month()-1][:3])
			default:
				out.WriteString(pad(int(t.Month()), n))
			}
		case 'd':
			out.WriteString(pad(t.Day(), n))
		case 'E', 'e', 'c', 'i':
			switch {
			case ja && n >= 4:
				out.WriteString(daysJa[t.Weekday()] + "曜日")
			case ja:
				out.WriteString(daysJa[t.Weekday()])
			case n >= 4:
				out.WriteString(daysEn[t.Weekday()])
			default:
				out.WriteString(daysEn[t.Weekday()][:3])
			}
		case 'H':
			out.WriteString(pad(t.Hour(), n))
		case 'h':
			out.WriteString(pad(h12, n))
		case 'm':
			out.WriteString(pad(t.Minute(), n))
		case 's':
			out.WriteString(pad(t.Second(), n))
		case 'a', 'b', 'B':
			switch pm := t.Hour() >= 12; {
			case ja && pm:
				out.WriteString("午後")
			case ja:
				out.WriteString("午前")
			case pm:
				out.WriteString("PM")
			default:
				out.WriteString("AM")
			}
		default:
			out.WriteString(string(rs[i : i+n]))
		}
		i += n
	}
	return out.String()
}

// dateText writes the value of a date, a time, or a date and a time in
// the format of its field.
func (f *filler) dateText(s *schema, v any) string {
	var t time.Time
	switch x := v.(type) {
	case time.Time:
		t = x
	case string:
		x = strings.TrimSpace(x)
		if x == "" {
			return ""
		}
		var ok bool
		if t, ok = tag.ParseDate(x); !ok {
			for _, layout := range []string{"15:04", "15:04:05"} {
				if tm, err := time.Parse(layout, x); err == nil {
					t, ok = time.Date(2021, 1, 1, tm.Hour(), tm.Minute(), tm.Second(), 0, time.UTC), true
					break
				}
			}
		}
		if !ok {
			f.warnf("%q is not a date that is read (%q); it is drawn as it is", x, s.name)
			return x
		}
	default:
		return f.textOf(v)
	}
	locale := s.str("locale", "en")
	format := strings.TrimSpace(s.str("format", ""))
	if format == "" || format == "undefined" {
		date, clock := "MM/dd/yyyy", "hh:mm aa"
		if strings.HasPrefix(locale, "ja") {
			date, clock = "yyyy/MM/dd", "HH:mm"
		}
		switch s.typ {
		case "date":
			format = date
		case "time":
			format = "HH:mm"
		default:
			format = date + " " + clock
		}
	}
	return formatDate(t, format, locale)
}

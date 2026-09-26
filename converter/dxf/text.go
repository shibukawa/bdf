package dxf

import (
	"math"
	"path"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// ttfFamilies maps the file names of common TrueType fonts to their
// families.
var ttfFamilies = map[string]string{
	"arial": "Arial", "arialbd": "Arial", "ariali": "Arial", "arialbi": "Arial", "arialn": "Arial Narrow",
	"arialuni": "Arial Unicode MS", "times": "Times New Roman", "timesbd": "Times New Roman", "timesi": "Times New Roman",
	"timesbi": "Times New Roman", "cour": "Courier New", "courbd": "Courier New", "couri": "Courier New",
	"calibri": "Calibri", "cambria": "Cambria", "verdana": "Verdana", "tahoma": "Tahoma", "segoeui": "Segoe UI",
	"century": "Century", "georgia": "Georgia", "consola": "Consolas",
	"msgothic": "MS Gothic", "msmincho": "MS Mincho", "meiryo": "Meiryo", "meiryob": "Meiryo",
	"yugothm": "Yu Gothic", "yugothr": "Yu Gothic", "yugothb": "Yu Gothic", "yugothl": "Yu Gothic",
	"yumin": "Yu Mincho", "yuminl": "Yu Mincho", "yumindb": "Yu Mincho",
	"simsun": "SimSun", "simhei": "SimHei", "msyh": "Microsoft YaHei", "msjh": "Microsoft JhengHei",
	"mingliu": "MingLiU", "malgun": "Malgun Gothic", "gulim": "Gulim", "batang": "Batang",
}

// fontSpec returns the font that stands in for a text style's font.
// AutoCAD's own SHX fonts are drawn with a sans-serif font, and big fonts
// (the SHX fonts of East Asian characters) with an East Asian one.
func fontSpec(st *style) cad.FontSpec {
	var f cad.FontSpec
	if st == nil {
		return f
	}
	f.Bold, f.Italic = st.bold, st.italic
	file := strings.TrimSpace(st.font)
	base := strings.ToLower(path.Base(strings.ReplaceAll(file, "\\", "/")))
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)
	switch {
	case st.family != "":
		f.Family = st.family
	case ext == ".ttf" || ext == ".ttc" || ext == ".otf":
		if fam, ok := ttfFamilies[name]; ok {
			f.Family = fam
			f.Bold = f.Bold || strings.HasSuffix(name, "bd") || strings.HasSuffix(name, "bi") || name == "meiryob" || name == "yugothb"
			f.Italic = f.Italic || strings.HasSuffix(name, "i") && name != "yumin" && name != "msyh"
		} else {
			f.Family = strings.TrimSuffix(path.Base(strings.ReplaceAll(file, "\\", "/")), path.Ext(file))
		}
	case ext == ".shx" || ext == "" && isSHXName(name):
		// the default sans-serif
		f.Italic = f.Italic || strings.Contains(name, "italic") || strings.Contains(name, "script")
	case file != "":
		f.Family = file
	}
	if big := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(st.bigFont), "@")); big != "" {
		switch {
		case strings.Contains(big, "gbc") || strings.Contains(big, "hz") || strings.Contains(big, "chinese"):
			f.EastAsian = "SimSun"
		case strings.Contains(big, "whg") || strings.Contains(big, "kor"):
			f.EastAsian = "Gulim"
		default:
			f.EastAsian = "MS Gothic"
		}
	}
	return f
}

func isSHXName(n string) bool {
	switch n {
	case "txt", "simplex", "romans", "romand", "romanc", "romant", "complex", "italic", "italicc", "italict",
		"isocp", "isocp2", "isocp3", "isoct", "isoct2", "isoct3", "gothice", "gothicg", "gothici", "greekc", "greeks",
		"monotxt", "scriptc", "scripts", "syastro", "symap", "symath", "symeteo", "symusic", "extfont", "extfont2":
		return true
	}
	return false
}

func (c *converter) style(name string) *style {
	if s := c.d.styles[key(name)]; s != nil {
		return s
	}
	if s := c.d.styles["STANDARD"]; s != nil {
		return s
	}
	return &style{name: "Standard", font: "txt", width: 1}
}

// segment is a part of a text with the same decorations.
type segment struct {
	s                           string
	underline, overline, strike bool
}

// textCodes replaces the control codes of single-line text: %%d (degree),
// %%p (plus-minus), %%c (diameter), %%nnn (a character code) and %%%, and
// cuts the text where %%u, %%o and %%k switch underlining, overlining and
// striking through.
func textCodes(s string) []segment {
	var out []segment
	cur := segment{}
	var b strings.Builder
	flush := func() {
		if b.Len() > 0 {
			cur.s = b.String()
			out = append(out, cur)
			b.Reset()
		}
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) && s[i+1] == '%' {
			ch := s[i+2]
			switch ch | 0x20 {
			case 'd':
				b.WriteString("°")
			case 'p':
				b.WriteString("±")
			case 'c':
				b.WriteString("Ø")
			case 'u':
				flush()
				cur.underline = !cur.underline
			case 'o':
				flush()
				cur.overline = !cur.overline
			case 'k':
				flush()
				cur.strike = !cur.strike
			default:
				if ch == '%' {
					b.WriteByte('%')
				} else if ch >= '0' && ch <= '9' {
					j := i + 2
					for j < len(s) && j < i+5 && s[j] >= '0' && s[j] <= '9' {
						j++
					}
					v, _ := strconv.Atoi(s[i+2 : j])
					b.WriteRune(rune(v))
					i = j - 1
					continue
				} else {
					b.WriteString(s[i : i+3])
				}
			}
			i += 2
			continue
		}
		b.WriteByte(s[i])
	}
	flush()
	return out
}

// text draws a TEXT, ATTRIB or ATTDEF.
func (c *converter) text(e *entity, x *ctx, p props) {
	if e.typ != "TEXT" && e.int(70, 0)&1 != 0 {
		return // an invisible attribute
	}
	raw := e.str(1)
	if e.typ == "ATTRIB" && e.has(101) {
		// an attribute holding multiline text (AutoCAD 2018)
		if m := embeddedMText(e); m != nil {
			c.mtext(m, x, p)
			return
		}
	}
	segs := textCodes(raw)
	if len(segs) == 0 {
		return
	}
	st := c.style(e.str(7))
	spec := fontSpec(st)
	fonts := c.fonts
	_, desc, capH := fonts.Metrics(spec)
	h := e.num(40, 0)
	if h <= 0 {
		h = st.height
	}
	if h <= 0 {
		h = 2.5
	}
	wf := e.num(41, st.width)
	if wf <= 0 {
		wf = 1
	}
	rot := e.num(50, 0) * math.Pi / 180
	obl := e.num(51, st.oblique) * math.Pi / 180
	gen := e.int(71, 0)
	if !e.has(71) {
		if st.backward {
			gen |= 2
		}
		if st.upset {
			gen |= 4
		}
	}
	ha := e.int(72, 0)
	vcode := 73
	if e.typ != "TEXT" {
		vcode = 74
	}
	va := e.int(vcode, 0)
	p1, p2 := e.vec3(10, [3]float64{}), e.vec3(11, [3]float64{})
	m := x.m.Mul(ocs(e, p1[2]))
	em := h / capH
	adv := 0.0
	for _, s := range segs {
		adv += fonts.Measure(spec, s.s)
	}
	sx := wf
	anchor := pt2(p1)
	var dx, dy float64
	switch {
	case ha == 3 || ha == 5:
		// aligned (the height follows the width) or fit (the width
		// stretches) between the two points
		d := pt2(p2).Sub(pt2(p1))
		if l := d.Len(); l > 0 && adv > 0 {
			rot = math.Atan2(d.Y, d.X)
			if ha == 3 {
				em = l / (adv * wf)
			} else {
				sx = l / (adv * em)
			}
		}
	default:
		if ha != 0 || va != 0 {
			anchor = pt2(p2)
		}
		width := adv * em * sx
		switch ha {
		case 1, 4:
			dx = -width / 2
		case 2:
			dx = -width
		}
		switch {
		case ha == 4 && va == 0:
			dy = -em * capH / 2
		case va == 1:
			dy = em * desc
		case va == 2:
			dy = -em * capH / 2
		case va == 3:
			dy = -em * capH
		}
	}
	mirror := canvas.Identity
	if gen&2 != 0 {
		mirror = canvas.Scale(-1, 1)
	}
	if gen&4 != 0 {
		mirror = mirror.Mul(canvas.Scale(1, -1))
	}
	tm := m.Mul(canvas.Translate(anchor.X, anchor.Y)).Mul(canvas.Matrix{math.Cos(rot), math.Sin(rot), -math.Sin(rot), math.Cos(rot), 0, 0}).
		Mul(mirror).Mul(canvas.Translate(dx, dy)).Mul(canvas.Matrix{1, 0, math.Tan(obl), 1, 0, 0}).Mul(canvas.Scale(em*sx, em))
	c.drawSegments(x.out, segs, spec, p.color, tm, cad.BreakBox)
}

// drawSegments draws the segments of a line one after the other from the
// origin of their em space.
func (c *converter) drawSegments(out *cad.Drawing, segs []segment, spec cad.FontSpec, color bdf.Color, m canvas.Matrix, brk cad.Break) {
	ex := 0.0
	for i, s := range segs {
		t := c.fonts.NewText(spec, s.s, color, m.Mul(canvas.Translate(ex, 0)))
		t.Underline, t.Overline, t.Strike = s.underline, s.overline, s.strike
		t.Break = brk
		if i > 0 {
			t.Break = cad.BreakNone
		}
		out.Text(t)
		ex += t.Advance
	}
}

// embeddedMText returns the MTEXT that an ATTRIB embeds after its 101
// marker.
func embeddedMText(e *entity) *entity {
	m := &entity{typ: "MTEXT"}
	in := false
	for _, t := range e.tags {
		if t.code == 101 {
			in = true
			continue
		}
		if in {
			if t.code >= 1000 {
				break
			}
			m.tags = append(m.tags, t)
		}
	}
	if !in {
		return nil
	}
	// the common properties come from the ATTRIB
	for _, code := range []int{8, 6, 62, 420, 370, 440, 7} {
		if t, ok := e.get(code); ok && !m.has(code) {
			m.tags = append(m.tags, t)
		}
	}
	return m
}

package dxf

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/linebreak"
)

// mstyle is the character format in effect in multiline text.
type mstyle struct {
	font                        cad.FontSpec
	height                      float64 // cap height in drawing units
	width                       float64 // width factor
	oblique                     float64 // degrees
	tracking                    float64
	color                       bdf.Color
	underline, overline, strike bool
}

// mtoken is a piece of a paragraph: a word (no break inside), a space, a
// tab or a stacked fraction.
type mtoken struct {
	kind     byte // 'w' word, ' ' space, '\t' tab, 's' stack
	s        string
	num, den string
	stack    byte // '^', '/' or '#'
	st       mstyle
	w        float64 // width in drawing units
	cjkBreak bool    // a line may break before this word without a space
}

type mpara struct {
	tokens []mtoken
	align  byte // 'l', 'c', 'r', 'j', 'd'; 0 for the attachment's
	indent float64
	left   float64
}

// parseMText splits the content of an MTEXT into paragraphs of tokens.
func (c *converter) parseMText(s string, base mstyle, baseStyle *style, dark bool, entColor bdf.Color) []mpara {
	s = caretDecode(s)
	paras := []mpara{{}}
	st := base
	var stack []mstyle
	var word strings.Builder
	para := func() *mpara { return &paras[len(paras)-1] }
	flush := func() {
		if word.Len() == 0 {
			return
		}
		c.addWord(para(), word.String(), st)
		word.Reset()
	}
	// value reads up to the terminating semicolon
	value := func(i int) (string, int) {
		j := strings.IndexByte(s[i:], ';')
		if j < 0 {
			return s[i:], len(s)
		}
		return s[i : i+j], i + j + 1
	}
	number := func(v string, cur float64) float64 {
		v = strings.TrimSpace(v)
		rel := strings.HasSuffix(v, "x") || strings.HasSuffix(v, "X")
		v = strings.TrimRight(v, "xX")
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return cur
		}
		if rel {
			return cur * math.Abs(f)
		}
		return math.Abs(f)
	}
	for i := 0; i < len(s); {
		ch := s[i]
		switch {
		case ch == '\\' && i+1 < len(s):
			cmd := s[i+1]
			i += 2
			switch cmd {
			case '\\', '{', '}':
				word.WriteByte(cmd)
			case '~':
				word.WriteString(" ")
			case 'P', 'N', 'X':
				flush()
				paras = append(paras, mpara{align: para().align})
			case 'L':
				flush()
				st.underline = true
			case 'l':
				flush()
				st.underline = false
			case 'O':
				flush()
				st.overline = true
			case 'o':
				flush()
				st.overline = false
			case 'K':
				flush()
				st.strike = true
			case 'k':
				flush()
				st.strike = false
			case 'f', 'F':
				flush()
				var v string
				v, i = value(i)
				st.font = mtextFont(v, baseStyle)
			case 'H':
				flush()
				var v string
				v, i = numValue(s, i)
				st.height = number(v, st.height)
			case 'W':
				flush()
				var v string
				v, i = numValue(s, i)
				st.width = number(v, st.width)
			case 'T':
				flush()
				var v string
				v, i = numValue(s, i)
				st.tracking = number(v, st.tracking)
			case 'Q':
				flush()
				var v string
				v, i = numValue(s, i)
				if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
					st.oblique = f
				}
			case 'A':
				_, i = numValue(s, i)
			case 'C':
				flush()
				var v string
				v, i = numValue(s, i)
				if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
					switch {
					case n == 0 || n == 256:
						st.color = entColor
					case n > 0 && n < 256:
						st.color = c.aci(n, dark)
					}
				}
			case 'c':
				flush()
				var v string
				v, i = numValue(s, i)
				if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
					// the bytes are in the order blue, green, red
					st.color = bdf.RGB(uint8(n), uint8(n>>8), uint8(n>>16))
				}
			case 'S':
				flush()
				var v string
				v, i = stackValue(s, i)
				num, den, typ := splitStack(v)
				tok := mtoken{kind: 's', num: num, den: den, stack: typ, st: st}
				c.measureStack(&tok)
				p := para()
				p.tokens = append(p.tokens, tok)
			case 'p':
				flush()
				var v string
				v, i = value(i)
				p := para()
				for _, f := range strings.Split(v, ",") {
					if f == "" {
						continue
					}
					switch f[0] {
					case 'q':
						if len(f) > 1 {
							p.align = f[1]
						}
					case 'i':
						p.indent, _ = strconv.ParseFloat(strings.TrimPrefix(f[1:], "x"), 64)
					case 'l':
						p.left, _ = strconv.ParseFloat(f[1:], 64)
					case 'x':
						if len(f) > 1 {
							switch f[1] {
							case 'q':
								if len(f) > 2 {
									p.align = f[2]
								}
							case 'i':
								p.indent, _ = strconv.ParseFloat(f[2:], 64)
							}
						}
					}
				}
			default:
				// an unknown code is shown as it is
				word.WriteByte('\\')
				word.WriteByte(cmd)
			}
		case ch == '{':
			flush()
			stack = append(stack, st)
			i++
		case ch == '}':
			flush()
			if n := len(stack); n > 0 {
				st = stack[n-1]
				stack = stack[:n-1]
			}
			i++
		case ch == '\n':
			flush()
			paras = append(paras, mpara{align: para().align})
			i++
		case ch == '\t':
			flush()
			p := para()
			p.tokens = append(p.tokens, mtoken{kind: '\t', st: st})
			i++
		case ch == ' ' || ch < ' ':
			flush()
			p := para()
			p.tokens = append(p.tokens, mtoken{kind: ' ', s: " ", st: st, w: c.fonts.RuneAdvance(st.font, ' ') * st.height / capOf(c, st.font) * st.width})
			i++
		case ch == '%' && i+2 < len(s) && s[i+1] == '%':
			switch s[i+2] | 0x20 {
			case 'd':
				word.WriteString("°")
			case 'p':
				word.WriteString("±")
			case 'c':
				word.WriteString("Ø")
			default:
				word.WriteString(s[i : i+3])
			}
			i += 3
		default:
			r, n := utf8.DecodeRuneInString(s[i:])
			word.WriteRune(r)
			i += n
		}
	}
	flush()
	return paras
}

// numValue reads a numeric value that a semicolon may or may not end.
func numValue(s string, i int) (string, int) {
	j := i
	for j < len(s) && (s[j] >= '0' && s[j] <= '9' || s[j] == '.' || s[j] == '-' || s[j] == '+' || s[j] == 'x' || s[j] == 'X' || s[j] == 'e' || s[j] == 'E') {
		j++
	}
	v := s[i:j]
	if j < len(s) && s[j] == ';' {
		j++
	}
	return v, j
}

// stackValue reads the expression of a \S code: up to an unescaped
// semicolon.
func stackValue(s string, i int) (string, int) {
	var b strings.Builder
	for i < len(s) {
		switch {
		case s[i] == '\\' && i+1 < len(s):
			b.WriteByte(s[i])
			b.WriteByte(s[i+1])
			i += 2
			continue
		case s[i] == ';':
			return b.String(), i + 1
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String(), i
}

func splitStack(v string) (num, den string, typ byte) {
	var b strings.Builder
	for i := 0; i < len(v); i++ {
		ch := v[i]
		if ch == '\\' && i+1 < len(v) {
			b.WriteByte(v[i+1])
			i++
			continue
		}
		if typ == 0 && (ch == '^' || ch == '/' || ch == '#') {
			num = b.String()
			b.Reset()
			typ = ch
			continue
		}
		b.WriteByte(ch)
	}
	if typ == 0 {
		return b.String(), "", 0
	}
	return num, strings.TrimPrefix(b.String(), " "), typ
}

// caretDecode replaces the caret notation of control characters (^J, ^I,
// "^ " for a caret).
func caretDecode(s string) string {
	if !strings.Contains(s, "^") {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '^' && i+1 < len(s) {
			b.WriteByte(byte((int(s[i+1]) - 64 + 126) % 126))
			i++
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// mtextFont reads a font code: "Arial|b1|i0|c0|p34" or "romans.shx".
func mtextFont(v string, base *style) cad.FontSpec {
	parts := strings.Split(v, "|")
	st := &style{font: strings.TrimSpace(parts[0])}
	if base != nil {
		st.bigFont = base.bigFont
	}
	for _, p := range parts[1:] {
		if len(p) < 2 {
			continue
		}
		switch p[0] {
		case 'b':
			st.bold = p[1] == '1'
		case 'i':
			st.italic = p[1] == '1'
		}
	}
	name := strings.ToLower(st.font)
	if !strings.HasSuffix(name, ".shx") && !strings.HasSuffix(name, ".ttf") && !strings.HasSuffix(name, ".ttc") && !isSHXName(name) {
		st.family = st.font
	}
	return fontSpec(st)
}

func capOf(c *converter, f cad.FontSpec) float64 {
	_, _, capH := c.fonts.Metrics(f)
	return capH
}

// breaksAround reports whether a line may break before or after r without
// a space (East Asian characters).
func breaksAround(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul) ||
		r >= 0x3000 && r < 0x3040 || r >= 0xff00 && r < 0xfff0
}

// addWord adds a word, cut where lines may break inside it (between East
// Asian characters).
func (c *converter) addWord(p *mpara, s string, st mstyle) {
	em := st.height / capOf(c, st.font)
	var cur []rune
	emit := func(cjkBreak bool) {
		if len(cur) == 0 {
			return
		}
		w := c.fonts.Measure(st.font, string(cur)) * em * st.width
		if st.tracking > 0 && st.tracking != 1 {
			w *= st.tracking
		}
		p.tokens = append(p.tokens, mtoken{kind: 'w', s: string(cur), st: st, w: w, cjkBreak: cjkBreak})
		cur = nil
	}
	var prev rune
	first := true
	brk := false
	for _, r := range s {
		if !first && (breaksAround(prev) || breaksAround(r)) && linebreak.Allowed(prev, r) {
			emit(brk)
			brk = true
		}
		cur = append(cur, r)
		prev, first = r, false
	}
	emit(brk)
}

func (c *converter) measureStack(t *mtoken) {
	sub := t.st
	sub.height *= 0.7
	em := sub.height / capOf(c, sub.font)
	nw := c.fonts.Measure(sub.font, t.num) * em * sub.width
	dw := c.fonts.Measure(sub.font, t.den) * em * sub.width
	switch t.stack {
	case '#':
		t.w = nw + dw + c.fonts.RuneAdvance(sub.font, '/')*em*sub.width
	default:
		t.w = math.Max(nw, dw)
	}
}

// mtext draws an MTEXT.
func (c *converter) mtext(e *entity, x *ctx, p props) {
	st := c.style(e.str(7))
	spec := fontSpec(st)
	h := e.num(40, 0)
	if h <= 0 {
		h = st.height
	}
	if h <= 0 {
		h = 2.5
	}
	base := mstyle{font: spec, height: h, width: 1, tracking: 1, color: p.color}
	if st.width > 0 {
		base.width = st.width
	}
	var content strings.Builder
	for _, t := range e.tags {
		if t.code == 3 {
			content.WriteString(t.s)
		}
	}
	content.WriteString(e.str(1))
	if e.int(72, 1) == 3 {
		c.warnOnce("mtextvertical", "vertical multiline text is drawn horizontally")
	}
	paras := c.parseMText(content.String(), base, st, x.dark, p.color)

	// the frame: origin at the insertion point, x along the text
	ins := e.vec3(10, [3]float64{})
	n := norm(e.vec3(210, [3]float64{0, 0, 1}))
	var xdir [3]float64
	if e.has(11) {
		xdir = norm(e.vec3(11, [3]float64{1, 0, 0}))
	} else {
		a := e.num(50, 0) * math.Pi / 180
		o := ocsMatrix(n, 0)
		v := cad.Apply(o, cad.Point{X: math.Cos(a), Y: math.Sin(a)})
		xdir = [3]float64{v.X, v.Y, 0}
	}
	ydir := norm(cross(n, xdir))
	frame := x.m.Mul(canvas.Matrix{xdir[0], xdir[1], ydir[0], ydir[1], ins[0], ins[1]})

	refW := e.num(41, 0)
	attach := e.int(71, 1)
	if attach < 1 || attach > 9 {
		attach = 1
	}
	spacing := e.num(44, 1)
	if spacing <= 0 {
		spacing = 1
	}
	exact := e.int(73, 1) == 2

	lines := layoutMText(paras, refW)
	if len(lines) == 0 {
		return
	}
	// line positions: the first baseline is a cap height below the top
	blockW := refW
	for _, l := range lines {
		if refW <= 0 {
			blockW = math.Max(blockW, l.width)
		}
	}
	y := 0.0
	for i := range lines {
		lh := lines[i].height
		if exact || lh == 0 {
			lh = h
		}
		if i == 0 {
			y = -lh
		} else {
			y -= lh * 5 / 3 * spacing
		}
		lines[i].y = y
	}
	blockH := -lines[len(lines)-1].y
	var ox, oy float64
	switch (attach - 1) % 3 {
	case 1:
		ox = -blockW / 2
	case 2:
		ox = -blockW
	}
	switch (attach - 1) / 3 {
	case 1:
		oy = blockH / 2
	case 2:
		oy = blockH
	}
	if fill := e.int(90, 0); fill&3 != 0 {
		scale := e.num(45, 1.5)
		margin := math.Max(scale-1, 0) * h / 2
		bg := bdf.RGB(255, 255, 255)
		switch {
		case fill&2 != 0:
			if x.dark {
				bg = darkBackground
			}
		case e.has(421):
			bg = rgb(e.int(421, 0))
		case e.has(63):
			bg = c.aci(e.int(63, 7), x.dark)
		}
		r := []cad.Point{{X: ox - margin, Y: oy + margin}, {X: ox + blockW + margin, Y: oy + margin},
			{X: ox + blockW + margin, Y: oy - blockH - margin}, {X: ox - margin, Y: oy - blockH - margin}}
		x.out.Fill((&cad.Path{}).Polyline(r, true).Transform(frame), cad.Fill{Color: bg}, false)
	}
	defAlign := [3]byte{'l', 'c', 'r'}[(attach-1)%3]
	for li, l := range lines {
		align := l.align
		if align == 0 {
			align = defAlign
		}
		lx := ox + l.left
		avail := blockW - l.left
		switch align {
		case 'c':
			lx += (avail - l.width) / 2
		case 'r':
			lx += avail - l.width
		}
		brk := cad.BreakBox
		if li > 0 {
			switch {
			case l.newPara:
				brk = cad.BreakParagraph
			case l.wrapCJK:
				brk = cad.BreakWrap
			default:
				brk = cad.BreakLine
			}
		}
		c.drawLine(x.out, l, frame.Mul(canvas.Translate(lx, oy+l.y)), brk)
	}
}

type mline struct {
	tokens  []mtoken
	width   float64
	height  float64 // the largest cap height
	y       float64
	left    float64
	align   byte
	newPara bool
	wrapCJK bool
}

// layoutMText breaks paragraphs into lines no wider than width (0: no
// wrapping).
func layoutMText(paras []mpara, width float64) []mline {
	var lines []mline
	for _, p := range paras {
		cur := mline{align: p.align, newPara: true}
		indent := p.indent + p.left
		if indent < 0 {
			indent = 0
		}
		cur.left = indent
		finish := func() {
			// trailing spaces do not count
			for len(cur.tokens) > 0 && cur.tokens[len(cur.tokens)-1].kind == ' ' {
				cur.width -= cur.tokens[len(cur.tokens)-1].w
				cur.tokens = cur.tokens[:len(cur.tokens)-1]
			}
			lines = append(lines, cur)
		}
		for _, t := range p.tokens {
			if t.kind == '\t' {
				// to the next multiple of four text heights
				step := t.st.height * 4
				t.w = step - math.Mod(cur.width, step)
			}
			if width > 0 && t.kind != ' ' && len(cur.tokens) > 0 && cur.left+cur.width+t.w > width+1e-9 {
				wrapCJK := cur.tokens[len(cur.tokens)-1].kind != ' ' && t.cjkBreak
				if t.kind == 'w' && !t.cjkBreak && cur.tokens[len(cur.tokens)-1].kind != ' ' {
					// no break opportunity here
				} else {
					finish()
					cur = mline{align: p.align, left: p.left, wrapCJK: wrapCJK}
					if cur.left < 0 {
						cur.left = 0
					}
				}
			}
			if t.kind == ' ' && len(cur.tokens) == 0 && !cur.newPara {
				continue // a space at the start of a wrapped line
			}
			cur.tokens = append(cur.tokens, t)
			cur.width += t.w
			if t.kind != ' ' {
				cur.height = math.Max(cur.height, t.st.height)
			}
		}
		finish()
	}
	return lines
}

// drawLine draws the tokens of a line from the origin of m (its baseline
// start), grouping words and spaces of the same format into texts.
func (c *converter) drawLine(out *cad.Drawing, l mline, m canvas.Matrix, brk cad.Break) {
	x := 0.0
	first := true
	var run []mtoken
	var runX float64
	flush := func() {
		if len(run) == 0 {
			return
		}
		st := run[0].st
		var b strings.Builder
		for _, t := range run {
			b.WriteString(t.s)
		}
		s := b.String()
		if strings.TrimSpace(s) == "" {
			run = nil
			return
		}
		em := st.height / capOf(c, st.font)
		sx := st.width
		if st.tracking > 0 {
			sx *= st.tracking
		}
		tm := m.Mul(canvas.Translate(runX, 0)).Mul(canvas.Matrix{1, 0, math.Tan(st.oblique * math.Pi / 180), 1, 0, 0}).Mul(canvas.Scale(em*sx, em))
		t := c.fonts.NewText(st.font, s, st.color, tm)
		t.Underline, t.Overline, t.Strike = st.underline, st.overline, st.strike
		t.Break = cad.BreakNone
		if first {
			t.Break = brk
			first = false
		}
		out.Text(t)
		run = nil
	}
	for _, t := range l.tokens {
		switch t.kind {
		case 'w', ' ':
			if len(run) > 0 && run[0].st != t.st {
				flush()
			}
			if len(run) == 0 {
				runX = x
			}
			run = append(run, t)
		case '\t':
			flush()
		case 's':
			flush()
			c.drawStack(out, t, m.Mul(canvas.Translate(x, 0)), &first, brk)
		}
		x += t.w
	}
	flush()
}

func (c *converter) drawStack(out *cad.Drawing, t mtoken, m canvas.Matrix, first *bool, brk cad.Break) {
	sub := t.st
	sub.height *= 0.7
	em := sub.height / capOf(c, sub.font)
	put := func(s string, x, y float64) float64 {
		if s == "" {
			return 0
		}
		tm := m.Mul(canvas.Translate(x, y)).Mul(canvas.Scale(em*sub.width, em))
		tx := c.fonts.NewText(sub.font, s, sub.color, tm)
		tx.Break = cad.BreakNone
		if *first {
			tx.Break = brk
			*first = false
		}
		out.Text(tx)
		return tx.Advance * em * sub.width
	}
	nw := c.fonts.Measure(sub.font, t.num) * em * sub.width
	dw := c.fonts.Measure(sub.font, t.den) * em * sub.width
	h := t.st.height
	switch t.stack {
	case '#':
		x := put(t.num, 0, h*0.35)
		x += put("/", x, 0)
		put(t.den, x, 0)
	case '/', '^':
		put(t.num, (t.w-nw)/2, h*0.6)
		put(t.den, (t.w-dw)/2, -h*0.45)
		if t.stack == '/' {
			line := (&cad.Path{}).MoveTo(0, h*0.4).LineTo(t.w, h*0.4).Transform(m)
			out.Stroke(line, cad.Pen{Color: t.st.color, Width: 0.25 * 72 / 25.4, Cap: cad.CapButt})
		}
	default:
		put(t.num, 0, 0)
	}
}

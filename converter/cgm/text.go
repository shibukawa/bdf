package cgm

import (
	"bytes"
	"math"
	"slices"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// Character set types of CHARACTER SET LIST.
const (
	set94 = iota
	set96
	set94Multi
	set96Multi
	setComplete
)

// charset is an entry of CHARACTER SET LIST: its type and the final
// bytes of its ISO 2022 designation.
type charset struct {
	typ  int
	tail string
}

// final returns the designation's final bytes, without the escape and
// intermediate bytes some writers include.
func (c charset) final() string {
	t := strings.TrimPrefix(c.tail, "\x1b")
	if c.typ == setComplete {
		return strings.TrimPrefix(t, "%")
	}
	return strings.TrimLeft(t, " !\"#$%&'()*+,-./")
}

func (c charset) utf8() bool {
	if c.typ != setComplete {
		return false
	}
	switch c.final() {
	case "G", "/G", "/H", "/I":
		return true
	}
	return false
}

func (c charset) utf16() bool {
	if c.typ != setComplete {
		return false
	}
	switch c.final() {
	case "/@", "/A", "/E", "/J", "/K", "/L", "@":
		return true
	}
	return false
}

// plain reports whether a set is ASCII or Latin-1, which writers declare
// whatever their strings hold.
func (c charset) plain() bool {
	f := c.final()
	return c.typ == set94 && (f == "B" || f == "@") || c.typ == set96 && f == "A"
}

var latin96 = map[string]*charmap.Charmap{
	"A": charmap.ISO8859_1, "B": charmap.ISO8859_2, "C": charmap.ISO8859_3, "D": charmap.ISO8859_4,
	"L": charmap.ISO8859_5, "G": charmap.ISO8859_6, "F": charmap.ISO8859_7, "H": charmap.ISO8859_8,
	"M": charmap.ISO8859_9, "V": charmap.ISO8859_10, "b": charmap.ISO8859_15,
}

// charsetAt returns the character set of an index of the list; the
// default set 1 is ASCII (G0) and Latin-1 (G1).
func (in *interp) charsetAt(i int, g1 bool) charset {
	if i >= 1 && i <= len(in.charsets) {
		return in.charsets[i-1]
	}
	if g1 {
		return charset{typ: set96, tail: "A"}
	}
	return charset{typ: set94, tail: "B"}
}

// decodeText converts the bytes of a string in the character sets of the
// text attributes: the designated set (G0) and the alternate one (G1),
// which shift out and shift in select in 7-bit codes and the high bytes
// select in 8-bit ones. Undeclared multibyte text (UTF-8 or Shift_JIS, as
// the whole metafile shows) is read as such.
func (in *interp) decodeText(b []byte) string {
	g0 := in.charsetAt(in.st.a.charSet, false)
	g1 := in.charsetAt(in.st.a.altCharSet, true)
	switch {
	case g0.utf8():
		return strings.ToValidUTF8(string(b), "�")
	case g0.utf16():
		return decodeUTF16(b)
	}
	ascii := true
	for _, c := range b {
		if c >= 0x80 || c < 0x20 {
			ascii = false
			break
		}
	}
	if ascii && g0.plain() {
		return string(b)
	}
	if in.fallback != nil && g0.plain() && (g1.plain() || g1.typ == set94) {
		if s, err := in.fallback.NewDecoder().Bytes(b); err == nil {
			return strings.ToValidUTF8(string(s), "�")
		}
	}
	var out strings.Builder
	shift := false
	for i := 0; i < len(b); i++ {
		c := b[i]
		switch {
		case c == 0x0e:
			shift = true
		case c == 0x0f:
			shift = false
		case c == 0x1b:
			// an escape sequence: intermediate bytes and a final one
			for i+1 < len(b) && b[i+1] >= 0x20 && b[i+1] <= 0x2f {
				i++
			}
			i++
		case c == '\t':
			out.WriteByte(' ')
		case c < 0x20 || c == 0x7f:
		case c < 0x80:
			set := g0
			if shift {
				set = g1
			}
			i += decodeIn(&out, set, c, b[i+1:], shift || set.typ == set96 || set.typ == set96Multi)
		case c < 0xa0:
			// C1 controls: writers mean Windows-1252
			r, _ := utf8.DecodeRune(mustDecode(charmap.Windows1252, []byte{c}))
			if r != 0 && r != utf8.RuneError {
				out.WriteRune(r)
			}
		default:
			set := g1
			if set.typ == set94 && set.plain() {
				set = charset{typ: set96, tail: "A"}
			}
			i += decodeIn(&out, set, c&0x7f, b[i+1:], true)
		}
	}
	return out.String()
}

// decodeIn writes the character of code c (7 bits) in a set; rest holds
// the bytes after it, of which it returns how many the character took.
func decodeIn(out *strings.Builder, set charset, c byte, rest []byte, _ bool) int {
	f := set.final()
	switch set.typ {
	case set94, setComplete:
		switch {
		case f == "J" && c == 0x5c:
			out.WriteRune('¥')
		case f == "J" && c == 0x7e:
			out.WriteRune('‾')
		case f == "I" && c >= 0x21 && c <= 0x5f:
			out.WriteRune(rune(0xff61 + int(c) - 0x21))
		case f == "A" && c == 0x23:
			out.WriteRune('£')
		default:
			out.WriteByte(c)
		}
		return 0
	case set96:
		cm := latin96[f]
		if cm == nil {
			cm = charmap.ISO8859_1
		}
		if c < 0x20 {
			return 0
		}
		out.Write(mustDecode(cm, []byte{c | 0x80}))
		return 0
	}
	// multibyte sets: two bytes of 7 bits
	if len(rest) == 0 {
		return 0
	}
	d := rest[0] & 0x7f
	var enc encoding.Encoding
	pair := []byte{c | 0x80, d | 0x80}
	switch f {
	case "@", "B", "Q", "O":
		enc = japanese.EUCJP
	case "D":
		enc = japanese.EUCJP
		pair = []byte{0x8f, c | 0x80, d | 0x80}
	case "A":
		enc = simplifiedchinese.GBK
	case "C":
		enc = korean.EUCKR
	}
	if enc == nil {
		out.WriteRune('�')
		return 1
	}
	s := mustDecode(enc, pair)
	if len(s) == 0 {
		s = []byte("�")
	}
	out.Write(s)
	return 1
}

func mustDecode(e encoding.Encoding, b []byte) []byte {
	s, err := e.NewDecoder().Bytes(b)
	if err != nil {
		return []byte("�")
	}
	return s
}

// decodeUTF16 decodes big-endian UTF-16 (little-endian with its byte
// order mark).
func decodeUTF16(b []byte) string {
	le := false
	if len(b) >= 2 {
		switch {
		case b[0] == 0xfe && b[1] == 0xff:
			b = b[2:]
		case b[0] == 0xff && b[1] == 0xfe:
			b, le = b[2:], true
		}
	}
	u := make([]uint16, len(b)/2)
	for i := range u {
		if le {
			u[i] = uint16(b[2*i]) | uint16(b[2*i+1])<<8
		} else {
			u[i] = uint16(b[2*i])<<8 | uint16(b[2*i+1])
		}
	}
	return string(utf16.Decode(u))
}

// decodeName decodes a name (a metafile, picture or font name): UTF-8 when
// it is valid, Shift_JIS when it reads as Japanese, Latin-1 otherwise.
func (in *interp) decodeName(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	if japaneseSJIS([][]byte{b}, 1) {
		return string(mustDecode(japanese.ShiftJIS, b))
	}
	return string(mustDecode(charmap.ISO8859_1, b))
}

// guessEncoding decides how to read the non-ASCII bytes of strings that
// the character sets do not say are multibyte: UTF-8 when all are valid
// UTF-8, Shift_JIS when they read as Japanese text (kana or common kanji,
// which Latin-1 letters do not make), nil for Latin-1.
func guessEncoding(strs [][]byte) (encoding.Encoding, string) {
	multi := false
	valid := true
	for _, s := range strs {
		if isASCII(s) {
			continue
		}
		multi = true
		if !utf8.Valid(s) {
			valid = false
		}
	}
	switch {
	case !multi:
		return nil, ""
	case valid:
		return encoding.Nop, "utf-8"
	case japaneseSJIS(strs, 2):
		return japanese.ShiftJIS, "shift_jis"
	}
	return nil, ""
}

func isASCII(b []byte) bool {
	for _, c := range b {
		if c >= 0x80 {
			return false
		}
	}
	return true
}

// japaneseSJIS reports whether the non-ASCII strings decode as Shift_JIS
// into at least n kana or JIS level 1 kanji (as the DXF reader does).
func japaneseSJIS(strs [][]byte, n int) bool {
	found := 0
	for _, s := range strs {
		if isASCII(s) {
			continue
		}
		out, err := japanese.ShiftJIS.NewDecoder().Bytes(s)
		if err != nil || bytes.ContainsRune(out, utf8.RuneError) {
			return false
		}
		for i := 0; i < len(s); i++ {
			b := s[i]
			if b < 0x81 || b >= 0xa0 && b < 0xe0 {
				continue
			}
			if b == 0x82 || b == 0x83 || b >= 0x88 && b <= 0x9f {
				found++
			}
			i++
		}
	}
	return found >= n
}

// fontSpec returns the font that stands in for a font of FONT LIST: the
// standard PostScript fonts by their families, the Hershey and other
// stroke fonts of CGM by the default sans-serif, Japanese fonts as the
// East Asian font.
func fontSpec(name string) cad.FontSpec {
	n := strings.ToLower(strings.TrimSpace(name))
	var f cad.FontSpec
	f.Bold = strings.Contains(n, "bold") || strings.Contains(n, "black") || strings.Contains(n, "heavy") || strings.Contains(n, "demi")
	f.Italic = strings.Contains(n, "italic") || strings.Contains(n, "oblique") || strings.Contains(n, "slant")
	base := n
	if i := strings.IndexAny(base, "-,:"); i > 0 && !strings.HasPrefix(base, "hershey") {
		base = base[:i]
	}
	base = strings.TrimSpace(base)
	switch {
	case n == "":
	case strings.Contains(n, "hershey") || strings.Contains(n, "simplex") || strings.Contains(n, "duplex") ||
		strings.Contains(n, "complex") || strings.Contains(n, "triplex") || strings.Contains(n, "cartographic") ||
		strings.Contains(n, "stroke") || strings.Contains(n, "iso_3098") || strings.Contains(n, "iso 3098"):
		// the default sans-serif
	case strings.Contains(n, "明朝") || strings.Contains(n, "mincho") || strings.Contains(n, "ming"):
		f.EastAsian = "MS Mincho"
	case strings.Contains(n, "ゴシック") || strings.Contains(n, "gothic"):
		f.EastAsian = "MS Gothic"
	case strings.HasPrefix(base, "helvetica") || strings.HasPrefix(base, "arial") || base == "swiss" || base == "sans" ||
		base == "sans-serif" || base == "sansserif":
		f.Family = "Helvetica"
	case strings.HasPrefix(base, "times") || base == "roman" || base == "serif" || base == "dutch":
		f.Family = "Times"
	case strings.HasPrefix(base, "courier") || base == "mono" || base == "monospace" || base == "typewriter":
		f.Family = "Courier"
	case strings.HasPrefix(base, "symbol"):
		f.Family = "Symbol"
	default:
		f.Family = strings.TrimSpace(name)
	}
	return f
}

// textPiece is a string of a text with its own attributes (a TEXT and the
// APPEND TEXT that continue it).
type textPiece struct {
	s                  string
	font               cad.FontSpec
	colour             bdf.Color
	expansion, spacing float64
}

// pendingText is a text whose TEXT or RESTRICTED TEXT is not final.
type pendingText struct {
	at         cad.Point
	restricted bool
	dx, dy     float64
	pieces     []textPiece
	a          attrs
}

// textAttrs returns the font, expansion, spacing and colour of text, with
// those that come from the text bundle.
func (in *interp) textAttrs() (spec cad.FontSpec, expansion, spacing float64, col bdf.Color) {
	a := &in.st.a
	font, c := a.font, a.textColour
	expansion, spacing = a.expansion, a.spacing
	if r, ok := in.st.textReps[a.textBundle]; ok {
		if a.asf[asfFont] {
			font = r.font
		}
		if a.asf[asfExpansion] {
			expansion = r.expansion
		}
		if a.asf[asfSpacing] {
			spacing = r.spacing
		}
		if a.asf[asfTextColour] {
			c = r.colour
		}
	}
	if font >= 1 && font <= len(in.fontSpecs) {
		spec = in.fontSpecs[font-1]
	} else if font != 1 {
		in.warnOnce("font:"+itoa(font), "text font %d is not in the font list; the default font is used", font)
	}
	if !(expansion > 0) || !isFinite(expansion) {
		expansion = 1
	}
	return spec, expansion, spacing, in.colourOf(c)
}

// textElement interprets TEXT, RESTRICTED TEXT and APPEND TEXT.
func (in *interp) textElement(c code, p params) {
	switch c {
	case eText, eRestrictedText:
		in.flushText()
		t := &pendingText{a: in.st.a}
		if c == eRestrictedText {
			t.restricted = true
			t.dx, t.dy = p.vdc(), p.vdc()
		}
		t.at = p.point()
		final := p.enum("NOTFINAL FINAL") == 1
		in.pending = t
		in.appendPiece(p.str())
		if final {
			in.flushText()
		}
	case eAppendText:
		final := p.enum("NOTFINAL FINAL") == 1
		if in.pending == nil {
			in.warnOnce("append", "APPEND TEXT without a text to continue is ignored")
			return
		}
		in.appendPiece(p.str())
		if final {
			in.flushText()
		}
	}
}

func (in *interp) appendPiece(b []byte) {
	s := in.decodeText(b)
	in.points += len(s)
	spec, exp, sp, col := in.textAttrs()
	in.pending.pieces = append(in.pending.pieces, textPiece{s: s, font: spec, colour: col, expansion: exp, spacing: sp})
}

// flushText draws the pending text.
func (in *interp) flushText() {
	t := in.pending
	if t == nil {
		return
	}
	in.pending = nil
	if !in.drawing() {
		return
	}
	in.drawText(t)
}

// drawText lays out a text (ISO/IEC 8632-1 §7.7.5): characters as tall as
// the character height from the baseline to the capline, along the
// character base vector and upright along the up vector (whose lengths
// set the width of characters relative to their height), in the text
// path, aligned on the text point.
func (in *interp) drawText(t *pendingText) {
	a := &t.a
	h := a.charHeight
	if h < 0 {
		h = in.pic.long / 100
	}
	up, base := a.up, a.base
	if up.Len() == 0 || base.Len() == 0 || !isFinite(up.Len()) || !isFinite(base.Len()) {
		up, base = cad.Point{Y: 1}, cad.Point{X: 1}
	}
	ratio := base.Len() / up.Len()
	u, b := up.Mul(1/up.Len()), base.Mul(1/base.Len())
	var pieces []textPiece
	for _, pc := range t.pieces {
		if pc.s != "" {
			pieces = append(pieces, pc)
		}
	}
	if len(pieces) == 0 || h <= 0 || !isFinite(h) {
		return
	}
	fonts := in.fonts
	path := a.path
	if path == 1 || path == 2 {
		// left and up: the characters in the reverse order
		slices.Reverse(pieces)
		for i := range pieces {
			r := []rune(pieces[i].s)
			slices.Reverse(r)
			pieces[i].s = string(r)
		}
	}
	asc, desc, capH := fonts.Metrics(pieces[0].font)
	if capH <= 0 {
		capH = 0.7
	}
	em0 := h / capH
	if path >= 2 {
		in.verticalText(t, pieces, u, b, h, em0)
		return
	}
	// widths along the baseline, in VDC
	type laid struct {
		pc      textPiece
		em, sx  float64
		spacing float64 // ems along x
		adv     float64 // ems along x
	}
	var ls []laid
	total, count := 0.0, 0
	for _, pc := range pieces {
		_, _, c := fonts.Metrics(pc.font)
		if c <= 0 {
			c = 0.7
		}
		em := h / c
		sx := em * pc.expansion * ratio
		l := laid{pc: pc, em: em, sx: sx, spacing: pc.spacing * h / sx}
		n := utf8.RuneCountInString(pc.s)
		l.adv = fonts.Measure(pc.font, pc.s) + l.spacing*float64(n)
		total += l.adv * sx
		count += n
		ls = append(ls, l)
	}
	// the spacing after the last character is not part of the text
	last := ls[len(ls)-1]
	width := total - last.spacing*last.sx
	fx, fy := 1.0, 1.0
	if t.restricted && width > 0 {
		dx, dy := math.Abs(t.dx), math.Abs(t.dy)
		capVDC, all := h, (asc+desc)*em0
		switch a.restrictedType {
		case 2: // boxed-cap
			fx, fy = dx/width, dy/capVDC
		case 3: // boxed-all
			fx, fy = dx/width, dy/all
		case 4: // isotropic-cap
			fx = math.Min(dx/width, dy/capVDC)
			fy = fx
		case 5: // isotropic-all
			fx = math.Min(dx/width, dy/all)
			fy = fx
		case 6: // justified: the spacing fills the box
			if count > 1 {
				extra := (dx - width) / float64(count-1)
				for i := range ls {
					ls[i].spacing += extra / ls[i].sx
					ls[i].adv += extra / ls[i].sx * float64(utf8.RuneCountInString(ls[i].pc.s))
				}
				width = dx
			} else if width > dx {
				fx = dx / width
			}
		default: // basic: shrunk to fit
			if s := math.Min(dx/width, dy/capVDC); s < 1 {
				fx, fy = s, s
			}
		}
		if !(fx > 0) || !isFinite(fx) {
			fx = 1
		}
		if !(fy > 0) || !isFinite(fy) {
			fy = 1
		}
		width *= fx
	}
	ascV, descV, capV := asc*em0*fy, desc*em0*fy, h*fy
	var ox, oy float64
	switch a.hAlign {
	case 0:
		if path == 1 {
			ox = -width
		}
	case 2:
		ox = -width / 2
	case 3:
		ox = -width
	case 4:
		ox = -width * a.contH
	}
	switch a.vAlign {
	case 1:
		oy = -ascV
	case 2:
		oy = -capV
	case 3:
		oy = -capV / 2
	case 5:
		oy = descV
	case 6:
		oy = descV - a.contV*(ascV+descV)
	}
	o := t.at.Add(b.Mul(ox)).Add(u.Mul(oy))
	for i, l := range ls {
		sx, em := l.sx*fx, l.em*fy
		m := canvas.Matrix{b.X * sx, b.Y * sx, u.X * em, u.Y * em, o.X, o.Y}
		tx := fonts.NewText(l.pc.font, l.pc.s, l.pc.colour, m)
		tx.Spacing = l.spacing
		tx.Advance = l.adv
		if i > 0 {
			tx.Break = cad.BreakNone
		}
		in.out.Text(tx)
		o = o.Add(b.Mul(l.adv * sx))
	}
}

// verticalText lays out a text of the up or down paths: a column of
// upright characters one em apart (plus the spacing), which is how the
// renderer draws vertical text. Up is drawn as down in reverse.
func (in *interp) verticalText(t *pendingText, pieces []textPiece, u, b cad.Point, h, em float64) {
	a := &t.a
	fonts := in.fonts
	var cells []float64
	var s strings.Builder
	spacing := pieces[0].spacing * h / em
	for _, pc := range pieces {
		for range pc.s {
			cells = append(cells, 1)
		}
		s.WriteString(pc.s)
	}
	length := (float64(len(cells)) + spacing*float64(len(cells)-1)) * em
	var ox, oy float64
	switch a.hAlign {
	case 1:
		ox = em / 2
	case 3:
		ox = -em / 2
	case 4:
		ox = em/2 - a.contH*em
	}
	bottom := a.path == 2 // up: normally aligned at the bottom
	switch a.vAlign {
	case 0:
		if bottom {
			oy = length
		}
	case 3:
		oy = length / 2
	case 4, 5:
		oy = length
	case 6:
		oy = length * (1 - a.contV)
	}
	o := t.at.Add(b.Mul(ox)).Add(u.Mul(oy))
	m := canvas.Matrix{-u.X * em, -u.Y * em, b.X * em, b.Y * em, o.X, o.Y}
	tx := fonts.NewCellText(pieces[0].font, s.String(), cells, spacing, pieces[0].colour, m)
	tx.Vertical = true
	in.out.Text(tx)
}

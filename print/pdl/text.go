package pdl

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/textpart"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/unicode/norm"
)

// The text mode: a form filled from a template of text (converter/textform)
// keeps its grid of columns and lines in its print image (form/textpart),
// and ESC/P and ESC/POS print its characters with the printer's own fonts,
// a column and a line of the grid at a time, instead of the dots of its
// pages. A dot-matrix printer prints text much faster than dots, and darker
// through the copies of a carbon-copy slip.
//
// Characters are sent as Japanese printers take them: ASCII, the katakana
// of JIS X 0201 in a byte, and the rest in JIS X 0208 in the Kanji mode.

// Print modes (Settings.Mode).
const (
	// ModeRaster prints the dots of the pages.
	ModeRaster = "raster"
	// ModeText prints the characters of the print image of a form of a
	// template of text in the printer's fonts (ESC/P, ESC/POS).
	ModeText = "text"
)

// HasTextMode reports whether a language prints in the text mode.
func HasTextMode(lang string) bool { return lang == ESCP || lang == ESCPOS }

// textImage returns the print image to print in the text mode, nil to
// print the dots. ESC/P prints a document with a print image in the text
// mode unless the settings ask the dots; ESC/POS only when they ask the
// text mode (a receipt of dots keeps its logo and its fonts).
func textImage(doc *bdf.Document, view *bdf.View, o *Options, res *Result) (*textpart.Image, error) {
	mode, lang := o.Settings.Mode, o.Profile.Language
	switch {
	case mode == ModeRaster:
		return nil, nil
	case !HasTextMode(lang):
		if mode == ModeText {
			return nil, fmt.Errorf("pdl: %s has no text mode (ESC/P and ESC/POS have)", Title(lang))
		}
		return nil, nil
	case mode == "" && lang == ESCPOS:
		return nil, nil
	}
	im, err := textpart.FromDocument(doc)
	if errors.Is(err, textpart.ErrNone) {
		if mode == ModeText {
			return nil, errors.New("pdl: the text mode prints the print image of a form filled from a template of text, which the document has not")
		}
		return nil, nil
	}
	if err == nil && (len(doc.Views) == 0 || doc.Views[0].ID != view.ID || len(im.Pages) != len(view.Pages)) {
		err = errors.New("pdl: the print image is of the pages of the document's first view, not of this one")
	}
	if err != nil {
		if mode == ModeText {
			return nil, err
		}
		res.Warnings = append(res.Warnings, err.Error()+": printed as dots")
		return nil, nil
	}
	return im, nil
}

// writeText writes a job in the text mode. A form that the text mode
// cannot print (its pitch, its width) is an error when the settings ask the
// text mode, and is printed as dots in the language's default mode: done
// is false then, and nothing is written.
func writeText(w *errWriter, im *textpart.Image, o *Options, res *Result) (done bool, err error) {
	var buf bytes.Buffer
	tr := &Result{}
	bw := &errWriter{w: &buf}
	switch o.Profile.Language {
	case ESCP:
		err = writeESCPText(bw, im, o, tr)
	case ESCPOS:
		err = writeESCPOSText(bw, im, o, tr)
	}
	if err != nil {
		if o.Settings.Mode == ModeText {
			return true, err
		}
		res.Warnings = append(res.Warnings, err.Error()+": printed as dots")
		return false, nil
	}
	res.Pages += tr.Pages
	res.Text = true
	res.Unsupported = append(res.Unsupported, tr.Unsupported...)
	res.Warnings = append(res.Warnings, tr.Warnings...)
	w.Write(buf.Bytes())
	return true, w.err
}

// Kinds of glyphs.
const (
	gASCII     = iota // a byte of ASCII
	gKana             // a byte of the katakana of JIS X 0201 (0xA1–0xDF)
	gYen              // ¥: 0x5C of the Japanese set of ESC R 8
	gPC437            // a byte of code page 437 (ESC/POS: box drawing)
	gPC850            // a byte of code page 850 (ESC/POS: ×, ÷ and Latin-1)
	gKanji            // JIS X 0208 in two columns
	gHalfKanji        // JIS X 0208 in one column (ESC/P: FS SI)
	gBlank            // columns left empty (the preprinted paper)
)

// glyph is a character of a line as the printer is sent it.
type glyph struct {
	kind int
	b    [2]byte
	// col and cols are the columns of the grid (from the left margin), as
	// printed: twice the character's in a range of double width.
	col, cols int
	size      string
}

// wide reports whether the glyph is of double width.
func (g *glyph) wide() bool { return g.size == textpart.SizeWide || g.size == textpart.SizeBig }

// tall reports whether the glyph is of double height.
func (g *glyph) tall() bool { return g.size == textpart.SizeTall || g.size == textpart.SizeBig }

// textEncoder makes the glyphs of the lines of a print image for a
// language, and notes the characters it could not send.
type textEncoder struct {
	im      *textpart.Image
	lang    string
	noPaper bool
	warned  map[string]bool
	res     *Result
	sjis    *encoding.Encoder
	// kanji and half are set once a glyph of each kind was made.
	kanji, half bool
}

func newTextEncoder(im *textpart.Image, o *Options, res *Result) *textEncoder {
	return &textEncoder{im: im, lang: o.Profile.Language, noPaper: slices.Contains(o.SkipRoles, bdf.RoleBackground),
		warned: map[string]bool{}, res: res, sjis: japanese.ShiftJIS.NewEncoder()}
}

func (e *textEncoder) warn(key, format string, args ...any) {
	if !e.warned[key] {
		e.warned[key] = true
		e.res.Warnings = append(e.res.Warnings, fmt.Sprintf("pdl: "+format, args...))
	}
}

// line returns the glyphs of a line, without the spaces at its end.
func (e *textEncoder) line(l textpart.Line) []glyph {
	var out []glyph
	for _, c := range e.im.Cells(l) {
		g := glyph{col: c.Col, cols: c.Cols, size: c.Size}
		switch {
		case c.Paper && e.noPaper:
			g.kind, g.size = gBlank, ""
		default:
			e.encode(&g, c.Text)
		}
		out = append(out, g)
	}
	for len(out) > 0 {
		g := out[len(out)-1]
		if g.kind != gBlank && !(g.kind == gASCII && g.b[0] == ' ') {
			break
		}
		out = out[:len(out)-1]
	}
	return out
}

// encode makes the bytes of a character.
func (e *textEncoder) encode(g *glyph, text string) {
	s := norm.NFC.String(text)
	r, n := utf8.DecodeRuneInString(s)
	if n < len(s) {
		e.warn("mark "+s, "%q: the marks on a character are left out", s)
	}
	w := textpart.Width(r, e.im.AmbiguousWide)
	switch {
	case r >= 0x20 && r < 0x7f:
		g.kind, g.b[0] = gASCII, byte(r)
		return
	case r >= 0xff61 && r <= 0xff9f:
		g.kind, g.b[0] = gKana, byte(r-0xff61+0xa1)
		return
	case r == 0xa5:
		g.kind = gYen
		return
	case w == 2:
		if b, ok := jisX0208(e.sjis, r); ok {
			g.kind, g.b = gKanji, b
			e.kanji = true
			return
		}
		e.warn("rune "+string(r), "%q (U+%04X) is not in JIS X 0208: printed as 〓", r, r)
		g.kind, g.b = gKanji, [2]byte{0x22, 0x2e}
		e.kanji = true
		return
	}
	// a character of one column besides ASCII and katakana: the box
	// drawing and the other characters of ambiguous width of a template of
	// "ambiguous: narrow"
	switch e.lang {
	case ESCP:
		if b, ok := jisX0208(e.sjis, r); ok {
			g.kind, g.b = gHalfKanji, b
			e.kanji, e.half = true, true
			return
		}
	case ESCPOS:
		lr := lightBox(r)
		if lr != r {
			e.warn("heavy", "the heavy lines of box drawing are printed light (code page 437 has none)")
		}
		if b, ok := charmap.CodePage437.EncodeRune(lr); ok && b >= 0x80 {
			g.kind, g.b[0] = gPC437, b
			return
		}
		if b, ok := charmap.CodePage850.EncodeRune(lr); ok && b >= 0x80 {
			g.kind, g.b[0] = gPC850, b
			return
		}
		if _, ok := jisX0208(e.sjis, r); ok {
			e.warn("rune "+string(r), "%q (U+%04X) takes one column of the form but two of a receipt printer's Kanji: printed as ?; a template of \"ambiguous: wide\" gives it two", r, r)
			g.kind, g.b[0] = gASCII, '?'
			return
		}
	}
	e.warn("rune "+string(r), "%q (U+%04X) is not among the printer's characters: printed as ?", r, r)
	g.kind, g.b[0] = gASCII, '?'
}

// jisX0208 returns the two bytes of a character of JIS X 0208 (its rows 1
// to 88, NEC's special characters of row 13 among them; not the extensions
// of IBM and NEC past them).
func jisX0208(sjis *encoding.Encoder, r rune) ([2]byte, bool) {
	b, err := sjis.Bytes([]byte(string(r)))
	if err != nil || len(b) != 2 || b[0] >= 0xed || b[1] < 0x40 || b[1] == 0x7f || b[1] > 0xfc {
		return [2]byte{}, false
	}
	s1, s2 := int(b[0]), int(b[1])
	if s1 >= 0xe0 {
		s1 -= 0x40
	}
	j1 := (s1-0x81)*2 + 0x21
	j2 := 0
	switch {
	case s2 >= 0x9f:
		j1++
		j2 = s2 - 0x9f + 0x21
	case s2 >= 0x80:
		j2 = s2 - 0x40 + 0x20
	default:
		j2 = s2 - 0x40 + 0x21
	}
	return [2]byte{byte(j1), byte(j2)}, true
}

// lightBox returns the light line of box drawing for a heavy one.
func lightBox(r rune) rune {
	switch r {
	case '━':
		return '─'
	case '┃':
		return '│'
	case '┏':
		return '┌'
	case '┓':
		return '┐'
	case '┗':
		return '└'
	case '┛':
		return '┘'
	case '┣':
		return '├'
	case '┫':
		return '┤'
	case '┳':
		return '┬'
	case '┻':
		return '┴'
	case '╋':
		return '┼'
	}
	return r
}

// pages returns the glyphs of the pages to print.
func (e *textEncoder) pages(seq []int) [][][]glyph {
	cache := map[int][][]glyph{}
	out := make([][][]glyph, len(seq))
	for i, p := range seq {
		lines, ok := cache[p]
		if !ok {
			for _, l := range e.im.Pages[p].Lines {
				lines = append(lines, e.line(l))
			}
			cache[p] = lines
		}
		out[i] = lines
	}
	return out
}

// textEnd returns the end column of the glyphs of the pages: the
// widest line, the left margin not included.
func textEnd(pages [][][]glyph) int {
	end := 0
	for _, p := range pages {
		for _, l := range p {
			if len(l) > 0 {
				g := l[len(l)-1]
				end = max(end, g.col+g.cols)
			}
		}
	}
	return end
}

// tallLine reports whether a line has glyphs of double height.
func tallLine(l []glyph) bool {
	return slices.ContainsFunc(l, func(g glyph) bool { return g.tall() && g.kind != gBlank })
}

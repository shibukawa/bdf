package pdl

import (
	"fmt"
	"math"

	"github.com/shibukawa/bdf/form/textpart"
)

// ESC/POS in the text mode: a receipt printer's characters are of a width
// of their own (Font A: 12 dots, Kanji 24), so a column of the grid is
// their width with the spacing that makes it as wide as the grid's at the
// printer's dots (ESC SP, FS S), and a character whose place the spacing
// leaves short of the grid's is moved there (ESC $). A line is printed
// and the paper fed to the next with ESC J, in dots (GS P), so lines of
// any pitch stay on the grid; a line of double height is fed over the two
// lines it takes. Kanji are sent in JIS codes in the Kanji mode (FS C 0,
// FS & … FS .), the katakana of a byte from the table of katakana (ESC t
// 1), the box drawing of one column from code page 437 (ESC t 0) and the
// other characters of one column from code page 850 (ESC t 2); double
// width and height are GS !, and FS ! for Kanji too. The JIS codes, rather
// than Shift_JIS, are those that the printers like the TM (Toshiba TEC's
// B-EP2DL) take too.

// escposFontDots and escposKanjiDots are the widths of Font A and of its
// Kanji, in dots, on the TM printers of 180 and 203 dpi.
const (
	escposFontDots  = 12
	escposKanjiDots = 24
)

func writeESCPOSText(w *errWriter, im *textpart.Image, o *Options, res *Result) error {
	s := o.Settings
	cut := escposUnsupported(s, res)
	dpi := o.Profile.DPI
	pitch := float64(dpi) / im.CPI
	adv := int(pitch)
	if adv < escposFontDots {
		return fmt.Errorf("pdl: a column of %g CPI is %.1f dots at %d dpi, narrower than the printer's characters (%d dots)", im.CPI, pitch, dpi, escposFontDots)
	}
	if adv-escposFontDots > 255 {
		return fmt.Errorf("pdl: a column of %g CPI is wider than ESC/POS spaces characters", im.CPI)
	}
	lineDots := float64(dpi) / im.LPI
	e := newTextEncoder(im, o, res)
	seq, _ := o.sequence(false)
	pages := e.pages(seq)
	if need := int(math.Round(float64(im.Margin.Left+textEnd(pages)) * pitch)); need > o.Profile.Width {
		return fmt.Errorf("pdl: the form is %d dots wide at %g CPI; the printer prints %d", need, im.CPI, o.Profile.Width)
	}
	w.str("\x1b@")
	w.Write([]byte{0x1d, 'P', byte(dpi), byte(dpi)}) // GS P: motion units of the dots
	w.str("\x1bR\x00")                               // ESC R 0: ASCII's \ and ~
	w.str("\x1bM\x00")                               // ESC M 0: Font A
	w.Write([]byte{0x1b, ' ', byte(adv - escposFontDots)})
	if e.kanji {
		w.str("\x1cC\x00") // FS C 0: JIS codes in the Kanji mode
		w.Write([]byte{0x1c, 'S', 0, byte(2*adv - escposKanjiDots)})
	}
	st := &escposState{w: w, pitch: pitch, left: im.Margin.Left, adv: adv, kadv: 2 * adv, table: -1, kanjiSize: e.kanji}
	for i, lines := range pages {
		y := func(line int) int { return int(math.Round(float64(im.Margin.Top+line) * lineDots)) }
		cur := 0
		for li, l := range lines {
			if len(l) == 0 {
				continue
			}
			escposFeed(w, y(li)-cur)
			st.line(l)
			next := li + 1
			if tallLine(l) {
				next++
			}
			escposFeed(w, max(y(next)-y(li), 1)) // ESC J prints the line
			cur = y(next)
		}
		end := y(len(lines))
		if im.Lines > 0 {
			end = int(math.Round(float64(im.Lines) * lineDots))
		}
		escposFeed(w, end-cur)
		res.Pages++
		escposCut(w, cut, i, len(pages))
	}
	escposDrawer(w, s)
	return w.err
}

// escposState is the modes an ESC/POS printer is in, and where it prints
// the next character of the line.
type escposState struct {
	w     *errWriter
	pitch float64
	left  int
	// adv and kadv are the dots a character of a byte and a Kanji take
	// with their spacing, at the normal size
	adv, kadv int
	x         int
	kanji     bool
	table     int
	size      byte
	// kanjiSize sends FS ! with GS !: Epson's GS ! sizes Kanji too, but
	// printers like it may size them only by FS ! (the TM-T90 is
	// documented with it)
	kanjiSize bool
}

func (s *escposState) ank() {
	if s.kanji {
		s.w.str("\x1c.")
		s.kanji = false
	}
}

// line sends the characters of a line (the printer starts it at the left
// of its print area).
func (s *escposState) line(l []glyph) {
	w := s.w
	s.x = 0
	for _, g := range l {
		if g.kind == gBlank || g.kind == gASCII && g.b[0] == ' ' {
			continue // (the next character is moved to its place)
		}
		if x := int(math.Round(float64(s.left+g.col) * s.pitch)); x != s.x {
			s.ank()
			w.Write([]byte{0x1b, '$', byte(x), byte(x >> 8)}) // ESC $: the place in the line
			s.x = x
		}
		size := byte(0)
		if g.wide() {
			size |= 0x10
		}
		if g.tall() {
			size |= 0x01
		}
		if size != s.size {
			s.ank()
			w.Write([]byte{0x1d, '!', size})
			if s.kanjiSize {
				w.Write([]byte{0x1c, '!', size>>2&4 | size<<3&8}) // FS !: 4 double width, 8 double height
			}
			s.size = size
		}
		mul := 1
		if g.wide() {
			mul = 2
		}
		switch g.kind {
		case gKanji:
			if !s.kanji {
				w.str("\x1c&")
				s.kanji = true
			}
			w.Write(g.b[:])
			s.x += s.kadv * mul
			continue
		case gKana:
			s.ank()
			s.setTable(1)
			w.Write(g.b[:1])
		case gPC437, gPC850:
			s.ank()
			if g.kind == gPC437 {
				s.setTable(0)
			} else {
				s.setTable(2)
			}
			w.Write(g.b[:1])
		case gYen:
			s.ank()
			w.str("\x1bR\x08\\\x1bR\x00") // ESC R 8: the Japanese set's ¥
		default:
			s.ank()
			w.Write(g.b[:1])
		}
		s.x += s.adv * mul
	}
	s.ank()
}

// setTable selects a table of characters of a byte with ESC t: 0 code
// page 437, 1 katakana, 2 code page 850.
func (s *escposState) setTable(n int) {
	if s.table != n {
		s.w.Write([]byte{0x1b, 't', byte(n)})
		s.table = n
	}
}

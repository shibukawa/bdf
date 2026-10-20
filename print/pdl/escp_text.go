package pdl

import (
	"fmt"
	"math"

	"github.com/shibukawa/bdf/form/textpart"
)

// ESC/P in the text mode: the pitch of the grid is the printer's pitch
// (ESC P, ESC M, ESC g; SI condensed), its line the line spacing (ESC 2,
// ESC 0, ESC 3), its paper the page length (ESC C) with the skip over the
// perforation (ESC N), and its left margin the printer's (ESC l). A line is
// its characters and CR LF, a page ends with a form feed. Kanji are sent in
// the Kanji mode (FS & … FS .) in the codes of JIS X 0208, spaced to take
// two columns (FS S, in 1/180 inch), and the characters of one column of
// JIS X 0208 (the box drawing of a template of "ambiguous: narrow") at half
// width (FS !, bit 1 as FS SI; FS T spaces them).
//
// The sizes are those of ESC/P-J84 as Canon's and Brother's emulations
// document them: ESC W doubles the width of all characters, and is bit 2
// of FS !; FS ! sets the half width (bit 1) and the double size of Kanji
// (bit 3, as FS W) at once, so a form with Kanji sets the three with FS !
// alone. ESC w doubles the height of the characters of a byte upward from
// their base line, without changing the line feed: a line of double height
// is printed on the line under it, its characters standing on its base
// line as the print image has them. Kanji of double size print their upper
// halves, feed 24/180 inch and print their lower halves with the other
// characters: the line is lowered so that these are on the line under it.
// Kanji of double height alone have no command and are printed at their
// size, with a warning.

// escpDots is the dots of ESC/P's units: 1/180 inch.
const escpDots = 180

// escpKanjiDots is the width of a Kanji of a 24-pin printer, in dots.
const escpKanjiDots = 24

// escpPitch returns the commands of a pitch and the dots of its column.
func escpPitch(cpi float64) (string, float64, error) {
	switch {
	case cpi == 10:
		return "\x1bP\x12", 18, nil // ESC P: 10 CPI; DC2: not condensed
	case cpi == 12:
		return "\x1bM\x12", 15, nil
	case cpi == 15:
		return "\x1bg\x12", 12, nil
	case cpi >= 17 && cpi <= 17.2:
		return "\x1bP\x0f", 180 / (120.0 / 7), nil // SI: 10 CPI condensed to 17.1
	case cpi == 20:
		return "\x1bM\x0f", 9, nil // SI: 12 CPI condensed to 20
	}
	return "", 0, fmt.Errorf("pdl: ESC/P has no pitch of %g CPI (10, 12, 15, 17.1 or 20)", cpi)
}

// escpSpacing returns the command of a line spacing.
func escpSpacing(lpi float64) (string, error) {
	switch lpi {
	case 6:
		return "\x1b2", nil
	case 8:
		return "\x1b0", nil
	}
	n := escpDots / lpi
	if n != math.Trunc(n) || n < 1 || n > 255 {
		return "", fmt.Errorf("pdl: ESC/P has no line spacing of %g LPI (6, 8, or 180/n)", lpi)
	}
	return string([]byte{0x1b, '3', byte(n)}), nil
}

// escpSpaces returns the left and right spacing that make a character of
// a width (dots) as wide as a number of columns.
func escpSpaces(col float64, cols, dots int) (byte, byte, bool) {
	extra := col*float64(cols) - float64(dots)
	if extra < 0 || extra != math.Trunc(extra) || extra > 2*255 {
		return 0, 0, false
	}
	l := int(extra) / 2
	return byte(l), byte(int(extra) - l), true
}

func writeESCPText(w *errWriter, im *textpart.Image, o *Options, res *Result) error {
	s := o.Settings
	escpUnsupported(s, res)
	if s.PerforationSkip != 0 {
		// the print image's skip is the form's
		res.unsupported("bdf-perforation-skip", s.PerforationSkip)
	}
	pitch, col, err := escpPitch(im.CPI)
	if err != nil {
		return err
	}
	spacing, err := escpSpacing(im.LPI)
	if err != nil {
		return err
	}
	e := newTextEncoder(im, o, res)
	seq, _ := o.sequence(false)
	pages := e.pages(seq)
	kl, kr, ok := escpSpaces(col, 2, escpKanjiDots)
	if e.kanji && !ok {
		return fmt.Errorf("pdl: the Kanji of ESC/P are wider than two columns of %g CPI", im.CPI)
	}
	hl, hr, ok := escpSpaces(col, 1, escpKanjiDots/2)
	if e.half && !ok {
		return fmt.Errorf("pdl: the half-width Kanji of ESC/P are wider than a column of %g CPI", im.CPI)
	}
	if im.Margin.Left > 255 {
		return fmt.Errorf("pdl: a left margin of %d columns", im.Margin.Left)
	}
	width := float64(o.Profile.Width) * escpDots / float64(o.Profile.DPI)
	if need := float64(im.Margin.Left+textEnd(pages)) * col; need > width+0.5 {
		return fmt.Errorf("pdl: the form is %.1f inches wide at %g CPI; the printer prints %.1f", need/escpDots, im.CPI, width/escpDots)
	}
	w.str("\x1b@")
	w.str("\x1bR\x00") // ESC R 0: ASCII's \ and ~ (not ¥ and ‾)
	w.str(pitch)
	w.str(spacing)
	if im.Lines > 0 {
		switch {
		case im.Lines <= 127:
			w.Write([]byte{0x1b, 'C', byte(im.Lines)})
		case math.Mod(float64(im.Lines), im.LPI) == 0 && float64(im.Lines)/im.LPI <= 22:
			w.Write([]byte{0x1b, 'C', 0, byte(float64(im.Lines) / im.LPI)})
		default:
			return fmt.Errorf("pdl: ESC/P has no page of %d lines", im.Lines)
		}
		if im.Skip > 0 {
			w.Write([]byte{0x1b, 'N', byte(min(im.Skip, 127))})
		} else {
			w.str("\x1bO") // ESC O: no skip but the printer's
		}
	}
	w.Write([]byte{0x1b, 'l', byte(im.Margin.Left)})
	if e.kanji {
		w.Write([]byte{0x1c, 'S', kl, kr})
	}
	if e.half {
		w.Write([]byte{0x1c, 'T', hl, hr})
	}
	lineDots := escpDots / im.LPI
	escpKanjiSizes(pages, lineDots, e)
	st := &escpState{w: w, fs: e.kanji}
	for _, lines := range pages {
		for range im.Margin.Top {
			w.str("\n")
		}
		for li := 0; li < len(lines); li++ {
			l := lines[li]
			if tallLine(l) && li+1 < len(lines) && len(lines[li+1]) == 0 {
				if escpBigKanji(l) {
					// the lower halves of the Kanji on the next line
					if d := int(math.Round(lineDots)) - escpKanjiDots; d > 0 {
						w.Write([]byte{0x1b, 'J', byte(d)})
					}
				} else {
					// enlarged upward from the next line's base line
					w.str("\r\n")
				}
				li++
			}
			for _, g := range l {
				st.glyph(g)
			}
			st.ank()
			w.str("\r\n")
		}
		if im.Lines > 0 {
			w.str("\x0c")
		}
		res.Pages++
	}
	return w.err
}

// escpKanjiSizes makes the sizes of Kanji those ESC/P prints: Kanji of
// double height alone at their size, and of double size at double width
// when a line is lower than their upper half.
func escpKanjiSizes(pages [][][]glyph, lineDots float64, e *textEncoder) {
	for _, p := range pages {
		for _, l := range p {
			for i := range l {
				g := &l[i]
				if g.kind != gKanji && g.kind != gHalfKanji {
					continue
				}
				switch {
				case g.size == textpart.SizeTall:
					e.warn("kanji tall", "ESC/P has no Kanji of double height alone: they are printed at their size")
					g.size = ""
				case g.size == textpart.SizeBig && lineDots < escpKanjiDots:
					e.warn("kanji big", "ESC/P prints Kanji of double size on lines of 24/180 inch or more: they are printed at double width")
					g.size = textpart.SizeWide
				}
			}
		}
	}
}

// escpBigKanji reports whether a line has Kanji of double size.
func escpBigKanji(l []glyph) bool {
	for _, g := range l {
		if (g.kind == gKanji || g.kind == gHalfKanji) && g.size == textpart.SizeBig {
			return true
		}
	}
	return false
}

// escpState is the modes an ESC/P printer is in.
type escpState struct {
	w *errWriter
	// fs sets the width with FS ! (a form with Kanji), else with ESC W
	fs    bool
	kanji bool
	// wide is ESC W (FS ! bit 2), tall ESC w, half and big FS ! bits 1
	// and 3
	wide, tall, half, big bool
}

// ank leaves the Kanji mode.
func (s *escpState) ank() {
	if s.kanji {
		s.w.str("\x1c.")
		s.kanji = false
	}
}

func (s *escpState) glyph(g glyph) {
	w := s.w
	switch g.kind {
	case gBlank:
		s.ank()
		s.set(false, false, s.half, false)
		for range g.cols {
			w.str(" ")
		}
	case gASCII, gKana, gYen:
		s.ank()
		s.set(g.wide(), g.tall(), s.half, false)
		switch g.kind {
		case gYen:
			w.str("\x1bR\x08\\\x1bR\x00") // ESC R 8: the Japanese set's ¥
		default:
			w.Write(g.b[:1])
		}
	case gKanji, gHalfKanji:
		big := g.size == textpart.SizeBig
		s.set(g.wide() && !big, s.tall, g.kind == gHalfKanji, big)
		if !s.kanji {
			w.str("\x1c&")
			s.kanji = true
		}
		w.Write(g.b[:])
	}
}

// set sets the sizes: out of the Kanji mode for ESC W and ESC w.
func (s *escpState) set(wide, tall, half, big bool) {
	w := s.w
	if s.fs && (wide != s.wide || half != s.half || big != s.big) {
		w.Write([]byte{0x1c, '!', b2(half)<<1 | b2(wide)<<2 | b2(big)<<3})
		s.wide, s.half, s.big = wide, half, big
	}
	if wide != s.wide {
		s.ank()
		w.Write([]byte{0x1b, 'W', b2(wide)})
		s.wide = wide
	}
	if tall != s.tall {
		s.ank()
		w.Write([]byte{0x1b, 'w', b2(tall)})
		s.tall = tall
	}
}

func b2(b bool) byte {
	if b {
		return 1
	}
	return 0
}

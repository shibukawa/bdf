package pdl

import (
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf/image/barcode"
)

// ZPL: a label is one ^XA … ^XZ format: its width (^PW) and length (^LL),
// its dots as one graphic field (^GFA, in the ASCII compression of ZPL
// II), and the QR codes and barcodes the printer prints as barcode fields
// over them. The copies of a label are ^PQ.

func writeZPL(w *errWriter, rd *render, res *Result) error {
	s := rd.o.Settings
	pages, each := rd.sequence(true)
	cut := s.Cut
	for i, page := range pages {
		l, err := rd.labelOf(page)
		if err != nil {
			return err
		}
		bm, syms, err := rd.page(page, l, rd.p.Barcodes && !s.NoPrinterBarcodes)
		if err != nil {
			return err
		}
		w.str("^XA\n")
		if i == 0 {
			zplSetup(w, rd, res)
		}
		w.printf("^PW%d^LL%d^LH0,0^LS0\n", l.W, l.H)
		switch cut {
		case CutPage:
			w.str("^MMC\n")
		case CutJob:
			w.str("^MMD\n") // a delayed cut, which ~JK makes at the end
		case CutNone:
			w.str("^MMT\n")
		}
		if !bm.blank(0, bm.H) {
			w.str("^FO0,0")
			zplGraphic(w, bm)
			w.str("^FS\n")
		}
		for _, sym := range syms {
			zplSymbol(w, sym)
			res.Symbols++
		}
		w.printf("^PQ%d\n^XZ\n", each)
		res.Pages += each
	}
	if cut == CutJob {
		w.str("~JK\n")
	}
	return w.err
}

// zplSetup writes what the first label sets for the job: the media's
// tracking, the darkness and the speed.
func zplSetup(w *errWriter, rd *render, res *Result) {
	s := rd.o.Settings
	switch s.Tracking {
	case "continuous":
		w.str("^MNN\n")
	case "gap", "web":
		w.str("^MNY\n")
	case "mark":
		w.str("^MNM\n")
	}
	if s.Darkness != 0 {
		// ^MD is relative to the printer's darkness, -30 to 30
		w.printf("^MD%d\n", int(math.Round(float64(s.Darkness)*30/100)))
	}
	if s.Speed != 0 {
		// ^PR is in inches a second, 2 to 14
		ips := int(math.Round(float64(s.Speed) / 100 / 25.4))
		w.printf("^PR%d\n", min(max(ips, 2), 14))
	}
	if s.Drawer != 0 {
		res.unsupported("bdf-cash-drawer", s.Drawer)
	}
	if s.Gap != 0 {
		res.unsupported("bdf-label-gap", s.Gap)
	}
	if s.PerforationSkip != 0 {
		res.unsupported("bdf-perforation-skip", s.PerforationSkip)
	}
}

// zplGraphic writes a bitmap as ^GFA in the ASCII compression of ZPL II:
// the hexadecimal digits of each row, a run of one digit as a count
// before it (G to Y for 1 to 19, g to z for 20 to 400), the rest of a row
// of 0 as ",", of F as "!", and a row as the one before as ":".
func zplGraphic(w *errWriter, bm *Bitmap) {
	n := bm.Stride * bm.H
	w.printf("^GFA,%d,%d,%d,", n, n, bm.Stride)
	var prev []byte
	var b strings.Builder
	for y := 0; y < bm.H; y++ {
		row := bm.Row(y)
		if prev != nil && string(row) == string(prev) {
			w.str(":")
			continue
		}
		prev = row
		b.Reset()
		zplRow(&b, row)
		w.str(b.String())
	}
}

const hexDigits = "0123456789ABCDEF"

// zplRow writes the digits of a row in the ASCII compression.
func zplRow(b *strings.Builder, row []byte) {
	digits := make([]byte, 0, 2*len(row))
	for _, c := range row {
		digits = append(digits, hexDigits[c>>4], hexDigits[c&15])
	}
	// the rest of the row of zeros or of F's
	end := len(digits)
	for end > 0 && digits[end-1] == '0' {
		end--
	}
	fill := byte(0)
	if end < len(digits) {
		fill = ','
	} else {
		for end > 0 && digits[end-1] == 'F' {
			end--
		}
		if end < len(digits) {
			fill = '!'
		}
	}
	for i := 0; i < end; {
		j := i + 1
		for j < end && digits[j] == digits[i] {
			j++
		}
		zplRun(b, j-i, digits[i])
		i = j
	}
	if fill != 0 {
		b.WriteByte(fill)
	}
}

// zplRun writes a run of n of a digit: each part of up to 419 a count
// (none for one) and the digit.
func zplRun(b *strings.Builder, n int, digit byte) {
	for n > 0 {
		k := min(n, 419)
		if k > 1 {
			if k >= 20 {
				b.WriteByte(byte('g' + k/20 - 1))
			}
			if k%20 > 0 {
				b.WriteByte(byte('G' + k%20 - 1))
			}
		}
		b.WriteByte(digit)
		n -= k
	}
}

// zplOrient is the orientation of a field of a symbol of quarter turns
// clockwise: N, R (90°), I (180°), B (270°).
func zplOrient(turns int) string { return string("NRIB"[turns&3]) }

// zplSymbol writes a symbol as a barcode field at its place.
func zplSymbol(w *errWriter, s *Symbol) {
	x, y, _, _ := s.Origin()
	m := s.Module()
	o := zplOrient(s.Turns)
	sym := s.Symbol
	text := "N"
	if sym.HasText() {
		text = "Y"
	}
	h := int(math.Round(float64(sym.BarHeight(float32(s.W), float32(s.H)))))
	data := s.Mark.Text
	w.printf("^FO%d,%d", x, y)
	switch sym.Kind {
	case barcode.QR:
		w.printf("^BQ%s,2,%d,%s", o, min(m, 100), s.Mark.Level)
		if ascii(data) {
			data = s.Mark.Level + "A," + data
		} else {
			// the automatic mode refuses bytes of 0x80 to 0x9F and 0xE0
			// to 0xFF: the bytes of UTF-8 as bytes
			data = fmt.Sprintf("%sM,B%04d%s", s.Mark.Level, len(data), data)
		}
	case barcode.EAN13:
		w.printf("^BY%d^BE%s,%d,%s,N", m, o, h, text)
		data = data[:12]
	case barcode.EAN8:
		w.printf("^BY%d^B8%s,%d,%s,N", m, o, h, text)
		data = data[:7]
	case barcode.UPCA:
		w.printf("^BY%d^BU%s,%d,%s,N,Y", m, o, h, text)
		data = data[:11]
	case barcode.Code39:
		w.printf("^BY%d,3.0^B3%s,N,%d,N,N", m, o, h)
	case barcode.Code128:
		w.printf("^BY%d^BC%s,%d,N,N,N,A", m, o, h)
	case barcode.GS1128:
		w.printf("^BY%d^BC%s,%d,N,N,N,D", m, o, h)
	case barcode.ITF, barcode.ITF14:
		// a narrow element of ITF is two modules, and a wide one five
		w.printf("^BY%d,2.5^B2%s,%d,N,N,N", 2*m, o, h)
	case barcode.NW7:
		start, stop := data[:1], data[len(data)-1:]
		w.printf("^BY%d,3.0^BK%s,N,%d,N,N,%s,%s", m, o, h, start, stop)
		data = data[1 : len(data)-1]
	}
	w.printf("^FH_^FD%s^FS\n", zplEscape(data))
}

// ascii reports whether a text is of ASCII alone.
func ascii(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// zplEscape writes the bytes of a field's data that ZPL would read as
// commands (^, ~, _) or that are not printable as _ and their two hex
// digits (^FH_).
func zplEscape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c > 0x7e || c == '^' || c == '~' || c == '_' {
			fmt.Fprintf(&b, "_%02X", c)
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

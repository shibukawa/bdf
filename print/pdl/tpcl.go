package pdl

import (
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf/image/barcode"
)

// TPCL (Toshiba TEC): a label is its size ({D, in tenths of a millimetre:
// the pitch of label and gap, the width, the length}), the image buffer
// cleared ({C}), its dots as one graphic ({SG; in nibble mode, which
// passes the framing of { | } whole), the QR codes and barcodes the
// printer prints ({XB} and their data {RB}), and the issue of its copies
// ({XS;I}). The commands are those of Toshiba TEC's external equipment
// interface specifications (B-EX4, B-FV4, BX400).

func writeTPCL(w *errWriter, rd *render, res *Result) error {
	s := rd.o.Settings
	if s.Drawer != 0 {
		res.unsupported("bdf-cash-drawer", s.Drawer)
	}
	if s.PerforationSkip != 0 {
		res.unsupported("bdf-perforation-skip", s.PerforationSkip)
	}
	gap := s.Gap
	if gap == 0 {
		gap = 3 * 72 / 25.4
	}
	if s.Tracking == "continuous" {
		gap = 0
	}
	ribbon := 0 // direct thermal
	if rd.p.Ribbon {
		ribbon = 1
	}
	// the sensor: 0 none, 1 reflective (a black mark), 2 transmissive (a gap)
	sensor := 2
	switch s.Tracking {
	case "continuous":
		sensor = 0
	case "mark":
		sensor = 1
	}
	speed := byte('3')
	if s.Speed != 0 {
		speed = tpclSpeed(float64(s.Speed) / 100 / 25.4)
	}
	pages, each := rd.sequence(true)
	size := ""
	for i, page := range pages {
		l, err := rd.labelOf(page)
		if err != nil {
			return err
		}
		bm, syms, err := rd.page(page, l, rd.p.Barcodes && !s.NoPrinterBarcodes)
		if err != nil {
			return err
		}
		dpi := rd.p.DPI
		length := tenths(float64(l.H), dpi)
		pitch := length + int(math.Round(gap/72*254))
		// (the width rounded down: the head is no wider)
		width := int(float64(l.W) * 10 / dotsPerMM(dpi))
		sz := fmt.Sprintf("{D%s,%04d,%s|}", tpclNum(pitch), min(width, 9999), tpclNum(length))
		if sz != size {
			// a new size clears the image buffer
			w.str(sz + "\n")
			size = sz
		}
		if i == 0 && s.Darkness != 0 {
			d := int(math.Round(float64(s.Darkness) / 10))
			sign := "+"
			if d < 0 {
				sign, d = "-", -d
			}
			w.printf("{AY;%s%02d,%d|}\n", sign, d, 1-ribbon)
		}
		w.str("{C|}\n")
		if !bm.blank(0, bm.H) {
			tpclGraphic(w, bm)
		}
		for n, sym := range syms {
			tpclSymbol(w, sym, n, dpi)
			res.Symbols++
		}
		cut := 0
		switch {
		case s.Cut == CutPage:
			cut = 1
		case s.Cut == CutJob && i == len(pages)-1:
			cut = min(each, 100)
		}
		w.printf("{XS;I,%04d,%03d%dC%c%d00|}\n", min(each, 9999), cut, sensor, speed, ribbon)
		res.Pages += each
	}
	return w.err
}

// tpclNum writes a length of 4 digits, or 5 past 9999.
func tpclNum(n int) string {
	if n > 9999 {
		return fmt.Sprintf("%05d", n)
	}
	return fmt.Sprintf("%04d", n)
}

// tpclSpeed is the speed of a B-EX4 of the speeds it takes (3, 6, 10, 12
// and 14 inches a second): the nearest one under ips.
func tpclSpeed(ips float64) byte {
	c := byte('3')
	for _, s := range []struct {
		ips  float64
		code byte
	}{{6, '6'}, {10, 'A'}, {12, 'C'}, {14, 'E'}} {
		if ips >= s.ips {
			c = s.code
		}
	}
	return c
}

// tpclGraphic writes a bitmap as {SG; at the label's corner, in nibble
// mode: a byte of 0x30 and four dots a nibble, the high one first.
func tpclGraphic(w *errWriter, bm *Bitmap) {
	w.printf("{SG;0000,0000,%04d,%s,0,", bm.W, tpclNum(bm.H))
	buf := make([]byte, 0, 2*bm.Stride)
	for y := 0; y < bm.H; y++ {
		buf = buf[:0]
		for _, c := range bm.Row(y) {
			buf = append(buf, 0x30|c>>4, 0x30|c&15)
		}
		w.Write(buf)
	}
	w.str("|}\n")
}

// tpclSymbol writes a QR code or a barcode as the field n ({XB}) and its
// data ({RB}): its place in tenths of a millimetre, its rotation, its
// modules, and its height in tenths of a millimetre.
func tpclSymbol(w *errWriter, s *Symbol, n, dpi int) {
	x, y, _, _ := s.Origin()
	m := s.Module()
	sym := s.Symbol
	rot := s.Turns & 3
	h := tenths(float64(sym.BarHeight(float32(s.W), float32(s.H))), dpi)
	text := 0
	if sym.HasText() {
		text = 1
	}
	data := s.Mark.Text
	w.printf("{XB%02d;%04d,%s,", n, tenths(float64(x), dpi), tpclNum(tenths(float64(y), dpi)))
	ean := func(typ string, digits int) {
		// the check digit attached by the printer
		w.printf("%s,3,%02d,%d,%04d,+0000000000,000,%d,00|}\n", typ, min(m, 15), rot, min(h, 1000), text)
		data = data[:digits]
	}
	bars := func(typ string, narrow, wide, gap int) {
		w.printf("%s,1,%02d,%02d,%02d,%02d,%02d,%d,%04d,+0000000000,0,00|}\n", typ, narrow, narrow, wide, wide, gap, rot, min(h, 1000))
	}
	switch sym.Kind {
	case barcode.QR:
		w.printf("T,%s,%02d,A,%d,M2|}\n", s.Mark.Level, min(m, 52), rot)
	case barcode.EAN13:
		ean("5", 12)
	case barcode.EAN8:
		ean("0", 7)
	case barcode.UPCA:
		ean("K", 11)
	case barcode.Code39:
		bars("3", m, 3*m, m)
	case barcode.ITF, barcode.ITF14:
		bars("2", 2*m, 5*m, 0)
	case barcode.NW7:
		bars("4", m, 3*m, m)
	}
	w.printf("{RB%02d;%s|}\n", n, tpclEscape(data))
}

// tpclEscape writes the control bytes of data as TPCL's escapes: > and
// the letter of 0x40 more, and > itself as >0.
func tpclEscape(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '>':
			b.WriteString(">0")
		case c < 0x20:
			b.WriteByte('>')
			b.WriteByte(0x40 + c)
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// tpclPrintable reports whether TPCL prints a symbol: not Code 128 nor
// GS1-128, whose check digits and code sets the specifications leave
// unclear, nor data of the framing's characters or of more than 2048
// bytes, nor more symbols than its 32 fields.
func tpclPrintable(s *Symbol) bool {
	switch s.Symbol.Kind {
	case barcode.Code128, barcode.GS1128:
		return false
	}
	return len(s.Mark.Text) <= 2048 && !strings.ContainsAny(s.Mark.Text, "{|}")
}

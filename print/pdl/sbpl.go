package pdl

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf/image/barcode"
)

// SBPL (SATO): a label is one job of <A> … <Z> between STX and ETX: its
// size (<A1>), its dots as graphics (<G>H, in hexadecimal, of a height in
// rows of 8 dots), the QR codes and barcodes the printer prints, and the
// copies (<Q>). The commands are those of SATO's programming references
// (CL4NX Plus, CT4-LX); see docs/print-server.md for what was not tried
// on a printer.

const (
	stx = "\x02"
	etx = "\x03"
	esc = "\x1b"
)

// sbplMaxRows is the rows of dots one <G> draws at most: 400 rows of 8
// dots, the bound of a printer of 203 dpi.
const sbplMaxRows = 400 * 8

func writeSBPL(w *errWriter, rd *render, res *Result) error {
	s := rd.o.Settings
	if s.Tracking != "" {
		res.unsupported("media-tracking", s.Tracking)
	}
	if s.Gap != 0 {
		res.unsupported("bdf-label-gap", s.Gap)
	}
	if s.Drawer != 0 {
		res.unsupported("bdf-cash-drawer", s.Drawer)
	}
	if s.PerforationSkip != 0 {
		res.unsupported("bdf-perforation-skip", s.PerforationSkip)
	}
	pages, each := rd.sequence(true)
	for i, page := range pages {
		l, err := rd.labelOf(page)
		if err != nil {
			return err
		}
		bm, syms, err := rd.page(page, l, rd.p.Barcodes && !s.NoPrinterBarcodes)
		if err != nil {
			return err
		}
		w.str(stx + esc + "A")
		w.printf(esc+"A1V%05dH%04d", l.H, l.W)
		if s.Darkness != 0 {
			// 1 (lightest) to 5 (darkest), 3 in the middle
			w.printf(esc+"#E%dA", min(max(int(math.Round(3+float64(s.Darkness)/50)), 1), 5))
		}
		if s.Speed != 0 {
			w.printf(esc+"CS%02d", min(max(int(math.Round(float64(s.Speed)/100/25.4)), 2), 14))
		}
		for y0 := 0; y0 < bm.H; y0 += sbplMaxRows {
			y1 := min(y0+sbplMaxRows, bm.H)
			if bm.blank(y0, y1) {
				continue
			}
			sbplGraphic(w, bm, y0, y1)
		}
		for _, sym := range syms {
			sbplSymbol(w, sym)
			res.Symbols++
		}
		switch {
		case s.Cut == CutPage:
			w.str(esc + "CT0001")
		case s.Cut == CutJob && i == len(pages)-1:
			w.printf(esc+"CT%04d", min(each, 9999))
		}
		w.printf(esc+"Q%d"+esc+"Z"+etx, each)
		res.Pages += each
	}
	return w.err
}

// sbplPos writes the position of what follows: SBPL counts dots from 1.
func sbplPos(w *errWriter, x, y int) { w.printf(esc+"V%d"+esc+"H%d", y+1, x+1) }

// sbplGraphic writes rows y0 to y1 of a bitmap as <G>H at their place, in
// rows of 8 dots (white below the last).
func sbplGraphic(w *errWriter, bm *Bitmap, y0, y1 int) {
	sbplPos(w, 0, y0)
	blocks := (y1 - y0 + 7) / 8
	w.printf(esc+"GH%03d%03d", bm.Stride, blocks)
	var b strings.Builder
	for y := y0; y < y0+8*blocks; y++ {
		b.Reset()
		if y < y1 {
			for _, c := range bm.Row(y) {
				b.WriteByte(hexDigits[c>>4])
				b.WriteByte(hexDigits[c&15])
			}
		} else {
			b.WriteString(strings.Repeat("0", 2*bm.Stride))
		}
		w.str(b.String())
	}
}

// sbplSymbol writes a QR code or a barcode at its place: <%> turns it
// (counterclockwise), <B> and <BD> are the bars of a ratio of 1:3 and
// 2:5, <BG> Code 128, <2D30> a QR code of model 2.
func sbplSymbol(w *errWriter, s *Symbol) {
	x, y, _, _ := s.Origin()
	m := s.Module()
	sym := s.Symbol
	h := int(math.Round(float64(sym.BarHeight(float32(s.W), float32(s.H)))))
	data := s.Mark.Text
	sbplPos(w, x, y)
	if t := s.Turns & 3; t != 0 {
		w.printf(esc+"%%%d", (4-t)&3)
	}
	bars := func(cmd, typ string, narrow int, data string) {
		w.printf(esc+"%s%s%02d%03d%s", cmd, typ, min(narrow, 36), min(max(h, 1), 999), data)
	}
	ean := "B"
	if sym.HasText() {
		ean = "BD" // with the digits under the bars
	}
	switch sym.Kind {
	case barcode.QR:
		w.printf(esc+"2D30,%s,%02d,1,0"+esc+"DN%04d,%s", s.Mark.Level, min(m, 99), len(data), data)
	case barcode.EAN13:
		bars(ean, "3", m, data[:12])
	case barcode.EAN8:
		bars(ean, "4", m, data[:7])
	case barcode.UPCA:
		bars(ean, "H", m, data[:11])
	case barcode.Code39:
		bars("B", "1", m, "*"+data+"*")
	case barcode.ITF, barcode.ITF14:
		// a narrow element of ITF is two modules, a wide one five
		bars("BD", "2", 2*m, data)
	case barcode.NW7:
		bars("B", "0", m, data)
	case barcode.Code128:
		w.printf(esc+"BG%02d%03d>H%s", min(m, 36), min(max(h, 1), 999), strings.ReplaceAll(data, ">", ">J"))
	}
}

// sbplPrintable reports whether SBPL prints a symbol: not GS1-128, whose
// application identifiers SBPL is not given, nor data that would end a
// command.
func sbplPrintable(s *Symbol) bool {
	if s.Symbol.Kind == barcode.GS1128 {
		return false
	}
	return !strings.ContainsAny(s.Mark.Text, stx+etx+esc)
}

// tenths is a length of dots at a resolution in tenths of a millimetre.
func tenths(dots float64, dpi int) int { return int(math.Round(dots * 10 / dotsPerMM(dpi))) }

// dotsPerMM is the dots a millimetre of a head: heads of "203 dpi" and
// "305 dpi" have 8 and 12 dots a millimetre.
func dotsPerMM(dpi int) float64 {
	switch dpi {
	case 203:
		return 8
	case 305:
		return 12
	case 609:
		return 24
	}
	return float64(dpi) / 25.4
}

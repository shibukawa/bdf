package pdl

// ESC/POS: a receipt is its rows of dots, sent in bands with GS v 0 (print
// raster bit image) as they are drawn, the paper fed over the rows with
// nothing on them, then cut. The printer prints a band once it has it, so
// a receipt as long as a page needs is never held whole.

// escposBand is the rows a GS v 0 command sends: within the buffers of
// the small printers.
const escposBand = 128

func writeESCPOS(w *errWriter, rd *render, res *Result) error {
	s := rd.o.Settings
	if s.Darkness != 0 {
		res.unsupported("print-darkness", s.Darkness)
	}
	if s.Speed != 0 {
		res.unsupported("print-speed", s.Speed)
	}
	if s.Tracking != "" && s.Tracking != "continuous" {
		res.unsupported("media-tracking", s.Tracking)
	}
	if s.PerforationSkip != 0 {
		res.unsupported("bdf-perforation-skip", s.PerforationSkip)
	}
	cut := s.Cut
	if cut == "" {
		cut = CutPage
	}
	pages, _ := rd.sequence(false)
	w.str("\x1b@") // ESC @: initialize
	// GS P: the motion units of the dots, so that ESC J feeds whole dots
	// (the TM-T88, T20 and m30 feed half-dots by default)
	w.Write([]byte{0x1d, 'P', byte(rd.p.DPI), byte(rd.p.DPI)})
	for i, page := range pages {
		l, err := rd.labelOf(page)
		if err != nil {
			return err
		}
		feed := 0 // blank rows not sent yet
		err = rd.rows(page, l, func(_ int, rows *Bitmap) error {
			for y0 := 0; y0 < rows.H; {
				if rows.blank(y0, y0+1) {
					feed++
					y0++
					continue
				}
				y1 := y0 + 1
				for y1 < rows.H && y1-y0 < escposBand && !rows.blank(y1, min(y1+8, rows.H)) {
					y1++
				}
				escposFeed(w, feed)
				feed = 0
				escposRaster(w, rows.Rows(y0, y1))
				y0 = y1
			}
			return w.err
		})
		if err != nil {
			return err
		}
		escposFeed(w, feed)
		res.Pages++
		if cut == CutPage || cut == CutJob && i == len(pages)-1 {
			w.str("\x1dVB\x00") // GS V 66 0: feed to the cutter, partial cut
		}
	}
	if s.Drawer > 0 {
		// ESC p m t1 t2: a pulse of 50 ms on, 500 ms off, on pin 2 or 5
		w.printf("\x1bp%c\x19\xfa", byte(s.Drawer-1))
	}
	return w.err
}

// escposRaster writes rows with GS v 0 in the normal mode.
func escposRaster(w *errWriter, b *Bitmap) {
	x, y := b.Stride, b.H
	w.Write([]byte{0x1d, 'v', '0', 0, byte(x), byte(x >> 8), byte(y), byte(y >> 8)})
	w.Write(b.Bits[:b.Stride*b.H])
}

// escposFeed feeds the paper over n rows of dots with ESC J, whose motion
// unit is the dot as GS P set it.
func escposFeed(w *errWriter, n int) {
	for n > 0 {
		k := min(n, 255)
		w.Write([]byte{0x1b, 'J', byte(k)})
		n -= k
	}
}

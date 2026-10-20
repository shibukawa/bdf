package pdl

import "math"

// ESC/P: a page of a continuous form is its rows of dots in bands of 24,
// each a bit image of 24 pins (ESC * 39: 180 dots an inch across and
// along, three bytes a column, the top dot in the high bit), the paper fed
// 24/180 inch after each and over the bands with nothing on them, and a
// form feed at its end. The page length (ESC C) makes the form feed reach
// the top of the next form, and the perforation skip (ESC N) keeps the
// print off the perforation. Epson's VP printers of today take ESC/P-J84,
// which has ESC * and not the raster graphics of ESC/P2.

func writeESCP(w *errWriter, rd *render, res *Result) error {
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
	if s.Cut != "" && s.Cut != CutNone {
		res.unsupported("bdf-cut", s.Cut)
	}
	if s.Drawer != 0 {
		res.unsupported("bdf-cash-drawer", s.Drawer)
	}
	if s.Gap != 0 {
		res.unsupported("bdf-label-gap", s.Gap)
	}
	if rd.p.DPI != 180 {
		// the dots of ESC * 39
		p := rd.p
		p.Width = int(math.Round(float64(p.Width) * 180 / float64(p.DPI)))
		p.DPI = 180
		rd.p = p
	}
	pages, _ := rd.sequence(false)
	w.str("\x1b@")
	length := ""
	for _, page := range pages {
		l, err := rd.labelOf(page)
		if err != nil {
			return err
		}
		if ln := escpLength(l.H, rd.o.Settings.PerforationSkip); ln != length {
			w.str(ln)
			length = ln
		}
		w.str("\x1b3\x18") // ESC 3 24: a line of 24/180 inch
		feed := 0
		var band *Bitmap
		next := 0 // the rows of the band filled
		flush := func(rows int) {
			if band.blank(0, rows) {
				feed += 24
				return
			}
			escpFeed(w, feed)
			feed = 0
			escpBand(w, band, rows)
			w.str("\r")
			feed = 24
		}
		err = rd.rows(page, l, func(_ int, rows *Bitmap) error {
			for y := 0; y < rows.H; y++ {
				if band == nil {
					band = NewBitmap(l.W, 24)
				}
				copy(band.Row(next), rows.Row(y))
				if next++; next == 24 {
					flush(24)
					clear(band.Bits)
					next = 0
				}
			}
			return w.err
		})
		if err != nil {
			return err
		}
		if next > 0 {
			flush(next)
		}
		w.str("\x0c") // FF: to the top of the next form
		res.Pages++
	}
	return w.err
}

// escpLength writes the page length of a form of n rows of 1/180 inch,
// and its perforation skip in points: in inches when it is whole ones,
// else in lines of 1/6 inch (ESC 2), as ESC/P-J84 has them.
func escpLength(rows int, skip float64) string {
	b := []byte{}
	if rows%180 == 0 && rows/180 <= 22 {
		b = append(b, 0x1b, 'C', 0, byte(rows/180))
	} else {
		lines := min(max(int(math.Round(float64(rows)/30)), 1), 127)
		b = append(b, 0x1b, '2', 0x1b, 'C', byte(lines))
	}
	if skip > 0 {
		lines := min(max(int(math.Round(skip/12)), 1), 127)
		b = append(b, 0x1b, '2', 0x1b, 'N', byte(lines))
	}
	return string(b)
}

// escpFeed feeds the paper n/180 inch with ESC J.
func escpFeed(w *errWriter, n int) {
	for n > 0 {
		k := min(n, 255)
		w.Write([]byte{0x1b, 'J', byte(k)})
		n -= k
	}
}

// escpBand writes the first rows of a band of 24 as ESC * 39: the columns
// up to the last with a dot, three bytes each, the top dot in the high
// bit of the first.
func escpBand(w *errWriter, band *Bitmap, rows int) {
	cols := 0
	for y := 0; y < rows; y++ {
		r := band.Row(y)
		for i := len(r) - 1; i >= 0; i-- {
			if r[i] != 0 {
				for b := 0; b < 8; b++ {
					if r[i]&(1<<b) != 0 {
						cols = max(cols, 8*i+8-b)
						break
					}
				}
				break
			}
		}
	}
	cols = min(cols, band.W)
	data := make([]byte, 0, 5+3*cols)
	data = append(data, 0x1b, '*', 39, byte(cols), byte(cols>>8))
	for x := 0; x < cols; x++ {
		var c [3]byte
		for y := 0; y < rows; y++ {
			if band.Black(x, y) {
				c[y/8] |= 0x80 >> (y & 7)
			}
		}
		data = append(data, c[:]...)
	}
	w.Write(data)
}

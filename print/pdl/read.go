package pdl

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Printed is what a printer of a language prints from what it is sent,
// as Read reads it: the subset of each language the writers of this
// package write. The tests of the writers and the fake printers
// (print/fakeprinter) read with it.
type Printed struct {
	// Pages are the labels, receipts or forms, in order.
	Pages []PrintedPage
	// Commands are the other commands, in order: those of ZPL without
	// the caret, of SBPL without the ESC, the issues of TPCL ({XS, {AY),
	// the setup of ESC/P ("1/6", "length 33", "skip 2").
	Commands []string
	// Drawer is the pin of the cash drawer ESC/POS opened, from 1; 0 for
	// none.
	Drawer int
}

// PrintedPage is a page of what a printer printed.
type PrintedPage struct {
	Dots *Bitmap
	// Symbols are the commands of the printer's own barcodes, with their
	// data: ZPL "^FO…^B…^FD…" (the data unescaped), SBPL "@x,y B…" or
	// "@x,y 2D30… DN…", TPCL "{XB…" and "{RB…".
	Symbols []string
	// Cut is "partial" when ESC/POS cut the paper after the page.
	Cut string
}

// Read reads what a printer of a language prints. width and dpi are the
// printer's dots across and resolution, which ESC/POS, ESC/P and TPCL do
// not say in what they are sent.
func Read(lang string, data []byte, width, dpi int) (*Printed, error) {
	switch lang {
	case ZPL:
		return ReadZPL(string(data))
	case SBPL:
		return ReadSBPL(string(data))
	case TPCL:
		return ReadTPCL(string(data), dpi)
	case ESCPOS:
		return ReadESCPOS(data, width, dpi)
	case ESCP:
		return ReadESCP(data, width)
	}
	return nil, fmt.Errorf("pdl: no reader of %q", lang)
}

// Sniff tells the language of what a printer is sent from its first
// bytes, or "" (ESC/POS and ESC/P both begin with ESC @; ESC/POS is told
// by GS P or GS v 0 in its first bytes).
func Sniff(head []byte) string {
	switch {
	case bytes.HasPrefix(head, []byte("^XA")), bytes.HasPrefix(head, []byte("~")) && bytes.Contains(head, []byte("^XA")):
		return ZPL
	case bytes.HasPrefix(head, []byte(stx+esc+"A")):
		return SBPL
	case bytes.HasPrefix(head, []byte("{D")):
		return TPCL
	case bytes.HasPrefix(head, []byte("\x1b@")):
		if bytes.Contains(head, []byte("\x1dP")) || bytes.Contains(head, []byte("\x1dv0")) {
			return ESCPOS
		}
		return ESCP
	}
	return ""
}

// ReadZPL reads the formats writeZPL writes: ^PW and ^LL, the graphic
// field at 0,0 in the ASCII compression, the barcode fields.
func ReadZPL(data string) (*Printed, error) {
	out := &Printed{}
	for _, f := range strings.Split(data, "^XZ") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !strings.HasPrefix(f, "^XA") {
			return nil, fmt.Errorf("pdl: a ZPL format of %.20q", f)
		}
		var l PrintedPage
		w, h := 0, 0
		var field strings.Builder
		for _, c := range strings.Split(f[3:], "^")[1:] {
			c = strings.TrimRight(c, "\n")
			switch {
			case strings.HasPrefix(c, "PW"):
				w, _ = strconv.Atoi(c[2:])
			case strings.HasPrefix(c, "LL"):
				h, _ = strconv.Atoi(c[2:])
			case strings.HasPrefix(c, "GFA,"):
				bm, err := ReadGFA(c[4:], w, h)
				if err != nil {
					return nil, err
				}
				l.Dots = bm
			case strings.HasPrefix(c, "FO"):
				field.Reset()
				field.WriteString("^" + c)
			case strings.HasPrefix(c, "FS"):
				if s := field.String(); s != "" && !strings.HasPrefix(s, "^FO0,0") || strings.Contains(s, "^B") {
					l.Symbols = append(l.Symbols, s)
				}
				field.Reset()
			case strings.HasPrefix(c, "FD"):
				field.WriteString("^FD" + unescapeFH(c[2:]))
			case field.Len() > 0:
				field.WriteString("^" + c)
			default:
				out.Commands = append(out.Commands, c)
			}
		}
		if w <= 0 || h <= 0 || w > 1<<15 || h > 1<<16 {
			return nil, fmt.Errorf("pdl: a ZPL label of %d × %d dots", w, h)
		}
		if l.Dots == nil {
			l.Dots = NewBitmap(w, h)
		}
		out.Pages = append(out.Pages, l)
	}
	return out, nil
}

func unescapeFH(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '_' && i+2 < len(s) {
			n, _ := strconv.ParseUint(s[i+1:i+3], 16, 8)
			b.WriteByte(byte(n))
			i += 2
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// ReadGFA reads the data of ^GFA (what follows "^GFA,") in the ASCII
// compression of ZPL II, for a label of w × h dots.
func ReadGFA(s string, w, h int) (*Bitmap, error) {
	parts := strings.SplitN(s, ",", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("pdl: ^GFA,%.20s", s)
	}
	total, _ := strconv.Atoi(parts[0])
	stride, _ := strconv.Atoi(parts[2])
	if stride != (w+7)/8 || total != stride*h {
		return nil, fmt.Errorf("pdl: ^GFA of %d bytes, %d a row, for %d × %d", total, stride, w, h)
	}
	bm := NewBitmap(w, h)
	var digits []byte
	var prev []byte
	flush := func() error {
		if len(digits) > 2*stride {
			return fmt.Errorf("pdl: a ^GFA row of %d digits", len(digits))
		}
		for len(digits) < 2*stride {
			digits = append(digits, '0')
		}
		row := make([]byte, stride)
		for i := range row {
			hi := strings.IndexByte(hexDigits, digits[2*i])
			lo := strings.IndexByte(hexDigits, digits[2*i+1])
			row[i] = byte(hi<<4 | lo)
		}
		prev = row
		digits = digits[:0]
		return nil
	}
	y := 0
	emit := func() error {
		if y >= h {
			return fmt.Errorf("pdl: ^GFA of more than %d rows", h)
		}
		copy(bm.Row(y), prev)
		y++
		return nil
	}
	count := 0
	for i := 0; i < len(parts[3]); i++ {
		c := parts[3][i]
		switch {
		case c >= 'G' && c <= 'Y':
			count += int(c-'G') + 1
		case c >= 'g' && c <= 'z':
			count += (int(c-'g') + 1) * 20
		case c == ',' || c == '!':
			fill := byte('0')
			if c == '!' {
				fill = 'F'
			}
			for len(digits) < 2*stride {
				digits = append(digits, fill)
			}
		case c == ':':
			if prev == nil {
				return nil, fmt.Errorf("pdl: ':' before a ^GFA row")
			}
			if err := emit(); err != nil {
				return nil, err
			}
			continue
		case strings.IndexByte(hexDigits, c) >= 0:
			n := max(count, 1)
			count = 0
			if len(digits)+n > 2*stride {
				return nil, fmt.Errorf("pdl: a ^GFA row of more than %d digits", 2*stride)
			}
			for range n {
				digits = append(digits, c)
			}
		default:
			return nil, fmt.Errorf("pdl: %q in ^GFA", c)
		}
		if len(digits) == 2*stride {
			if err := flush(); err != nil {
				return nil, err
			}
			if err := emit(); err != nil {
				return nil, err
			}
		}
	}
	if y != h {
		return nil, fmt.Errorf("pdl: ^GFA of %d rows of %d", y, h)
	}
	return bm, nil
}

// ReadSBPL reads the jobs writeSBPL writes: the size of <A1>, the <G>H
// graphics at their places, the commands of the symbols.
func ReadSBPL(data string) (*Printed, error) {
	out := &Printed{}
	for _, job := range strings.Split(data, stx) {
		if job == "" {
			continue
		}
		if !strings.HasPrefix(job, esc+"A") || !strings.HasSuffix(job, esc+"Z"+etx) {
			return nil, fmt.Errorf("pdl: an SBPL job of %.30q", job)
		}
		var bm *Bitmap
		var syms []string
		x, y := 0, 0
		cmds := strings.Split(strings.TrimSuffix(job, esc+"Z"+etx), esc)[1:]
		for i := 0; i < len(cmds); i++ {
			c := cmds[i]
			switch {
			case c == "A":
			case strings.HasPrefix(c, "A1V"):
				var w, h int
				fmt.Sscanf(c, "A1V%05dH%04d", &h, &w)
				if w <= 0 || h <= 0 || w > 1<<15 || h > 1<<16 {
					return nil, fmt.Errorf("pdl: an SBPL label of %d × %d dots", w, h)
				}
				bm = NewBitmap(w, h)
			case strings.HasPrefix(c, "V"):
				y, _ = strconv.Atoi(c[1:])
				y--
			case strings.HasPrefix(c, "H"):
				x, _ = strconv.Atoi(c[1:])
				x--
			case strings.HasPrefix(c, "GH"):
				if bm == nil || len(c) < 8 {
					return nil, fmt.Errorf("pdl: <G>H before <A1>")
				}
				bb, _ := strconv.Atoi(c[2:5])
				cc, _ := strconv.Atoi(c[5:8])
				hex := c[8:]
				if x != 0 || bb != bm.Stride || len(hex) != 2*bb*cc*8 {
					return nil, fmt.Errorf("pdl: <G>H%03d%03d of %d digits at %d", bb, cc, len(hex), x)
				}
				for r := 0; r < cc*8; r++ {
					for i := 0; i < bb; i++ {
						v, _ := strconv.ParseUint(hex[2*(r*bb+i):2*(r*bb+i)+2], 16, 8)
						if y+r >= 0 && y+r < bm.H {
							bm.Row(y + r)[i] = byte(v)
						} else if v != 0 {
							return nil, fmt.Errorf("pdl: SBPL dots below the label")
						}
					}
				}
			case strings.HasPrefix(c, "2D30"), strings.HasPrefix(c, "B"), strings.HasPrefix(c, "%"):
				s := fmt.Sprintf("@%d,%d ", x, y) + c
				if strings.HasPrefix(c, "2D30") && i+1 < len(cmds) {
					i++
					s += " " + cmds[i]
				}
				syms = append(syms, s)
			default:
				out.Commands = append(out.Commands, c)
			}
		}
		if bm == nil {
			return nil, fmt.Errorf("pdl: an SBPL job without <A1>")
		}
		out.Pages = append(out.Pages, PrintedPage{Dots: bm, Symbols: syms})
	}
	return out, nil
}

// ReadTPCL reads what writeTPCL writes: the label's size, the graphic in
// nibble mode, the fields of the symbols, the issue.
func ReadTPCL(data string, dpi int) (*Printed, error) {
	out := &Printed{}
	var w, h int
	var bm *Bitmap
	var syms []string
	for _, c := range strings.Split(strings.TrimSpace(data), "|}\n") {
		c = strings.TrimSuffix(c, "|}")
		switch {
		case strings.HasPrefix(c, "{D"):
			var pitch, width, length int
			fmt.Sscanf(c, "{D%d,%d,%d", &pitch, &width, &length)
			w, h = int(float64(width)*dotsPerMM(dpi)/10+0.5), int(float64(length)*dotsPerMM(dpi)/10+0.5)
			if w <= 0 || h <= 0 || w > 1<<15 || h > 1<<16 {
				return nil, fmt.Errorf("pdl: a TPCL label of %d × %d dots", w, h)
			}
		case c == "{C":
			if w == 0 {
				return nil, fmt.Errorf("pdl: {C before {D")
			}
			bm, syms = NewBitmap(w, h), nil
		case strings.HasPrefix(c, "{SG;0000,0000,"):
			parts := strings.SplitN(c[len("{SG;0000,0000,"):], ",", 4)
			if bm == nil || len(parts) != 4 {
				return nil, fmt.Errorf("pdl: {SG before {C")
			}
			gw, _ := strconv.Atoi(parts[0])
			gh, _ := strconv.Atoi(parts[1])
			nib := parts[3]
			stride := (gw + 7) / 8
			if parts[2] != "0" || gw != bm.W || gh != bm.H || len(nib) != 2*stride*gh {
				return nil, fmt.Errorf("pdl: {SG of %s × %s type %s, %d bytes, for %d × %d", parts[0], parts[1], parts[2], len(nib), bm.W, bm.H)
			}
			for i := 0; i < stride*gh; i++ {
				hi, lo := nib[2*i], nib[2*i+1]
				if hi&0xf0 != 0x30 || lo&0xf0 != 0x30 {
					return nil, fmt.Errorf("pdl: a nibble of %#x %#x", hi, lo)
				}
				bm.Bits[i] = hi<<4 | lo&15
			}
		case strings.HasPrefix(c, "{XB"), strings.HasPrefix(c, "{RB"):
			syms = append(syms, c)
		case strings.HasPrefix(c, "{XS;I,"):
			if bm == nil {
				return nil, fmt.Errorf("pdl: {XS before {C")
			}
			out.Pages = append(out.Pages, PrintedPage{Dots: bm, Symbols: syms})
			out.Commands = append(out.Commands, c)
		case strings.HasPrefix(c, "{AY;"):
			out.Commands = append(out.Commands, c)
		default:
			return nil, fmt.Errorf("pdl: a TPCL command %.30q", c)
		}
	}
	return out, nil
}

// ReadESCPOS prints the commands that writeESCPOS writes: the rows of GS
// v 0 and the feeds of ESC J, a receipt ending at each cut, and the
// drawer of ESC p.
func ReadESCPOS(data []byte, width, dpi int) (*Printed, error) {
	if width <= 0 || width > 1<<15 {
		return nil, fmt.Errorf("pdl: ESC/POS %d dots wide", width)
	}
	out := &Printed{}
	var rows [][]byte
	stride := (width + 7) / 8
	end := func(cut string) {
		bm := NewBitmap(width, len(rows))
		for y, r := range rows {
			copy(bm.Row(y), r)
		}
		out.Pages = append(out.Pages, PrintedPage{Dots: bm, Cut: cut})
		rows = nil
	}
	for i := 0; i < len(data); {
		switch {
		case bytes.HasPrefix(data[i:], []byte("\x1b@")):
			i += 2
		case bytes.HasPrefix(data[i:], []byte("\x1dP")) && i+3 < len(data):
			if dpi != 0 && (int(data[i+2]) != dpi || int(data[i+3]) != dpi) {
				return nil, fmt.Errorf("pdl: GS P %d %d at %d dpi", data[i+2], data[i+3], dpi)
			}
			i += 4
		case bytes.HasPrefix(data[i:], []byte("\x1bJ")) && i+2 < len(data):
			for range int(data[i+2]) {
				rows = append(rows, make([]byte, stride))
			}
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1dv0\x00")) && i+8 <= len(data):
			x := int(data[i+4]) | int(data[i+5])<<8
			y := int(data[i+6]) | int(data[i+7])<<8
			i += 8
			if x > stride || i+x*y > len(data) {
				return nil, fmt.Errorf("pdl: GS v 0 of %d × %d bytes at %d", x, y, i)
			}
			for range y {
				r := make([]byte, stride)
				copy(r, data[i:i+x])
				rows = append(rows, r)
				i += x
			}
		case bytes.HasPrefix(data[i:], []byte("\x1dVB")) && i+3 < len(data):
			end("partial")
			i += 4
		case bytes.HasPrefix(data[i:], []byte("\x1bp")) && i+4 < len(data):
			out.Drawer = int(data[i+2]) + 1
			i += 5
		default:
			return nil, fmt.Errorf("pdl: an ESC/POS byte 0x%02x at %d", data[i], i)
		}
	}
	if len(rows) > 0 {
		end("")
	}
	return out, nil
}

// ReadESCP prints what writeESCP writes: the bit images of 24 pins at the
// place the feeds take the paper to, a page from form feed to form feed,
// and the page length and skip.
func ReadESCP(data []byte, width int) (*Printed, error) {
	if width <= 0 || width > 1<<15 {
		return nil, fmt.Errorf("pdl: ESC/P %d dots wide", width)
	}
	out := &Printed{}
	var rows [][]byte
	stride := (width + 7) / 8
	y := 0
	row := func(r int) []byte {
		for len(rows) <= r {
			rows = append(rows, make([]byte, stride))
		}
		return rows[r]
	}
	for i := 0; i < len(data); {
		switch {
		case bytes.HasPrefix(data[i:], []byte("\x1b@")):
			i += 2
		case bytes.HasPrefix(data[i:], []byte("\x1b3\x18")):
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1b2")):
			out.Commands = append(out.Commands, "1/6")
			i += 2
		case bytes.HasPrefix(data[i:], []byte("\x1bC\x00")) && i+3 < len(data):
			out.Commands = append(out.Commands, fmt.Sprintf("length %din", data[i+3]))
			i += 4
		case bytes.HasPrefix(data[i:], []byte("\x1bC")) && i+2 < len(data):
			out.Commands = append(out.Commands, fmt.Sprintf("length %d", data[i+2]))
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1bN")) && i+2 < len(data):
			out.Commands = append(out.Commands, fmt.Sprintf("skip %d", data[i+2]))
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1bJ")) && i+2 < len(data):
			y += int(data[i+2])
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1b*\x27")) && i+4 < len(data):
			cols := int(data[i+3]) | int(data[i+4])<<8
			i += 5
			if cols > width || i+3*cols > len(data) || y > 1<<16 {
				return nil, fmt.Errorf("pdl: ESC * of %d columns", cols)
			}
			for x := 0; x < cols; x++ {
				for b := 0; b < 24; b++ {
					if data[i+3*x+b/8]&(0x80>>(b&7)) != 0 {
						row(y + b)[x/8] |= 0x80 >> (x & 7)
					}
				}
			}
			i += 3 * cols
		case data[i] == '\r':
			i++
		case data[i] == '\x0c':
			bm := NewBitmap(width, len(rows))
			for r, b := range rows {
				copy(bm.Row(r), b)
			}
			out.Pages = append(out.Pages, PrintedPage{Dots: bm})
			rows, y = nil, 0
			i++
		default:
			return nil, fmt.Errorf("pdl: an ESC/P byte 0x%02x at %d", data[i], i)
		}
	}
	return out, nil
}

package pdl

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// zplLabel is a label as a reader of ZPL prints it: its dots, the barcode
// fields over them, and its commands.
type zplLabel struct {
	dots     *Bitmap
	fields   []string // "^FO… ^B… ^FD…" with the data unescaped
	commands []string
}

// readZPL reads the formats writeZPL writes: ^PW and ^LL, the graphic
// field at 0,0 in the ASCII compression, the barcode fields.
func readZPL(data string) ([]zplLabel, error) {
	var out []zplLabel
	for _, f := range strings.Split(data, "^XZ") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if !strings.HasPrefix(f, "^XA") {
			return nil, fmt.Errorf("a format of %.20q", f)
		}
		var l zplLabel
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
				bm, err := readGFA(c[4:], w, h)
				if err != nil {
					return nil, err
				}
				l.dots = bm
			case strings.HasPrefix(c, "FO"):
				field.Reset()
				field.WriteString("^" + c)
			case strings.HasPrefix(c, "FS"):
				if s := field.String(); s != "" && !strings.HasPrefix(s, "^FO0,0") || strings.Contains(s, "^B") {
					l.fields = append(l.fields, s)
				}
				field.Reset()
			case strings.HasPrefix(c, "FD"):
				field.WriteString("^FD" + unescapeFH(c[2:]))
			case field.Len() > 0:
				field.WriteString("^" + c)
			default:
				l.commands = append(l.commands, c)
			}
		}
		if l.dots == nil {
			l.dots = NewBitmap(w, h)
		}
		out = append(out, l)
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

// readGFA reads the data of ^GFA in the ASCII compression of ZPL II.
func readGFA(s string, w, h int) (*Bitmap, error) {
	parts := strings.SplitN(s, ",", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("^GFA,%.20s", s)
	}
	total, _ := strconv.Atoi(parts[0])
	stride, _ := strconv.Atoi(parts[2])
	if stride != (w+7)/8 || total != stride*h {
		return nil, fmt.Errorf("^GFA of %d bytes, %d a row, for %d × %d", total, stride, w, h)
	}
	bm := NewBitmap(w, h)
	var digits []byte
	var prev []byte
	flush := func() error {
		if len(digits) > 2*stride {
			return fmt.Errorf("a row of %d digits", len(digits))
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
			return fmt.Errorf("more than %d rows", h)
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
				return nil, fmt.Errorf("':' before a row")
			}
			if err := emit(); err != nil {
				return nil, err
			}
			continue
		case strings.IndexByte(hexDigits, c) >= 0:
			n := max(count, 1)
			count = 0
			for range n {
				digits = append(digits, c)
			}
		default:
			return nil, fmt.Errorf("%q in ^GFA", c)
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
		return nil, fmt.Errorf("%d rows of %d", y, h)
	}
	return bm, nil
}

func TestZPLRows(t *testing.T) {
	for _, row := range [][]byte{
		{0, 0, 0, 0}, {0xff, 0xff, 0xff}, {0x12, 0x34, 0, 0}, {0xf0, 0xff, 0xff},
		bytes.Repeat([]byte{0xaa}, 300), append(bytes.Repeat([]byte{0x11}, 250), 0x80),
	} {
		var b strings.Builder
		zplRow(&b, row)
		bm, err := readGFA(fmt.Sprintf("%d,%d,%d,%s", len(row), len(row), len(row), b.String()), len(row)*8, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bm.Row(0), row) {
			t.Errorf("%x: %s reads as %x", row, b.String(), bm.Row(0))
		}
	}
}

func TestZPL(t *testing.T) {
	d, v := testDoc(t, 4*72, 3*72, 2)
	for _, c := range []struct {
		name    string
		profile string
		o       Options
		labels  int
		symbols bool
	}{
		{"printer barcodes", "zd421-203", Options{Copies: 3}, 2, true},
		{"collated", "zd421-300", Options{Copies: 2, Collate: true, Settings: Settings{NoPrinterBarcodes: true, Darkness: 50, Speed: 10160, Tracking: "gap", Cut: CutPage}}, 4, false},
	} {
		p, _ := LookupProfile(c.profile)
		o := c.o
		o.Profile = p
		var buf bytes.Buffer
		res, err := Write(&buf, d, v, &o)
		if err != nil {
			t.Fatal(c.name, err)
		}
		labels, err := readZPL(buf.String())
		if err != nil {
			t.Fatal(c.name, err)
		}
		if len(labels) != c.labels || res.Pages != 2*o.Copies {
			t.Fatalf("%s: %d labels, %+v", c.name, len(labels), res)
		}
		rd := newRender(d, v, &o, p)
		for i, l := range labels {
			page := i % 2
			lb, _ := rd.labelOf(page)
			sameDots(t, fmt.Sprintf("%s label %d", c.name, i+1), l.dots, wantDots(t, d, v, page, lb, map[bool][]string{true: {"code128", "qr"}}[c.symbols]))
			if c.symbols {
				if len(l.fields) != 2 || !strings.Contains(l.fields[0], "^BCN,") || !strings.Contains(l.fields[0], "^FDBDF-0042") ||
					!strings.Contains(l.fields[1], "^BQN,2,") && !strings.Contains(l.fields[1], ",M") || !strings.Contains(l.fields[1], "^FDMA,https://example.com/r/42") {
					t.Errorf("%s label %d: fields %q", c.name, i+1, l.fields)
				}
			} else if len(l.fields) != 0 {
				t.Errorf("%s label %d: fields %q", c.name, i+1, l.fields)
			}
		}
		if c.symbols && res.Symbols != 4 {
			t.Errorf("%s: %d symbols", c.name, res.Symbols)
		}
		out := buf.String()
		for _, want := range map[bool][]string{true: {"^PQ3\n"}, false: {"^PQ1\n", "^MD15", "^PR4", "^MNY", "^MMC"}}[c.symbols] {
			if !strings.Contains(out, want) {
				t.Errorf("%s: no %s", c.name, want)
			}
		}
	}
}

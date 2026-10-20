package pdl

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// readSBPL reads the jobs writeSBPL writes: the size of <A1>, the <G>H
// graphics at their places, the commands of the symbols.
func readSBPL(data string) (labels []*Bitmap, symbols [][]string, rest []string, err error) {
	for _, job := range strings.Split(data, stx) {
		if job == "" {
			continue
		}
		if !strings.HasPrefix(job, esc+"A") || !strings.HasSuffix(job, esc+"Z"+etx) {
			return nil, nil, nil, fmt.Errorf("a job of %.30q", job)
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
				bm = NewBitmap(w, h)
			case strings.HasPrefix(c, "V"):
				y, _ = strconv.Atoi(c[1:])
				y--
			case strings.HasPrefix(c, "H"):
				x, _ = strconv.Atoi(c[1:])
				x--
			case strings.HasPrefix(c, "GH"):
				bb, _ := strconv.Atoi(c[2:5])
				cc, _ := strconv.Atoi(c[5:8])
				hex := c[8:]
				if x != 0 || bb != bm.Stride || len(hex) != 2*bb*cc*8 {
					return nil, nil, nil, fmt.Errorf("<G>H%03d%03d of %d digits at %d", bb, cc, len(hex), x)
				}
				for r := 0; r < cc*8; r++ {
					for i := 0; i < bb; i++ {
						v, _ := strconv.ParseUint(hex[2*(r*bb+i):2*(r*bb+i)+2], 16, 8)
						if y+r < bm.H {
							bm.Row(y + r)[i] = byte(v)
						} else if v != 0 {
							return nil, nil, nil, fmt.Errorf("dots below the label")
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
				rest = append(rest, c)
			}
		}
		labels = append(labels, bm)
		symbols = append(symbols, syms)
	}
	return labels, symbols, rest, nil
}

func TestSBPL(t *testing.T) {
	d, v := testDoc(t, 4*72, 3*72, 2)
	p, _ := LookupProfile("cl4nx-203")
	var buf bytes.Buffer
	res, err := Write(&buf, d, v, &Options{Profile: p, Copies: 2, Settings: Settings{Darkness: 100, Speed: 15240, Cut: CutJob}})
	if err != nil {
		t.Fatal(err)
	}
	labels, syms, rest, err := readSBPL(buf.String())
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 2 || res.Pages != 4 || res.Symbols != 4 {
		t.Fatalf("%d labels, %+v", len(labels), res)
	}
	rd := newRender(d, v, &Options{Profile: p}, p)
	for i, bm := range labels {
		l, _ := rd.labelOf(i)
		sameDots(t, fmt.Sprintf("label %d", i+1), bm, wantDots(t, d, v, i, l, []string{"code128", "qr"}))
		if len(syms[i]) != 2 || !regexp.MustCompile(`^@\d+,\d+ BG\d\d\d\d\d>HBDF-0042$`).MatchString(syms[i][0]) ||
			!regexp.MustCompile(`^@\d+,\d+ 2D30,M,\d\d,1,0 DN0024,https://example.com/r/42$`).MatchString(syms[i][1]) {
			t.Errorf("label %d: symbols %q", i+1, syms[i])
		}
	}
	want := []string{"#E5A", "CS06", "Q2", "#E5A", "CS06", "CT0002", "Q2"}
	if fmt.Sprint(rest) != fmt.Sprint(want) {
		t.Errorf("commands %q, want %q", rest, want)
	}
}

// readTPCL reads what writeTPCL writes: the label's size, the graphic in
// nibble mode, the fields of the symbols, the issue.
func readTPCL(data string, dpi int) (labels []*Bitmap, fields [][]string, issues []string, err error) {
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
		case c == "{C":
			bm, syms = NewBitmap(w, h), nil
		case strings.HasPrefix(c, "{SG;0000,0000,"):
			parts := strings.SplitN(c[len("{SG;0000,0000,"):], ",", 4)
			gw, _ := strconv.Atoi(parts[0])
			gh, _ := strconv.Atoi(parts[1])
			nib := parts[3]
			stride := (gw + 7) / 8
			if parts[2] != "0" || gw != bm.W || gh != bm.H || len(nib) != 2*stride*gh {
				return nil, nil, nil, fmt.Errorf("{SG of %s × %s type %s, %d bytes, for %d × %d", parts[0], parts[1], parts[2], len(nib), bm.W, bm.H)
			}
			for i := 0; i < stride*gh; i++ {
				hi, lo := nib[2*i], nib[2*i+1]
				if hi&0xf0 != 0x30 || lo&0xf0 != 0x30 {
					return nil, nil, nil, fmt.Errorf("a nibble of %#x %#x", hi, lo)
				}
				bm.Bits[i] = hi<<4 | lo&15
			}
		case strings.HasPrefix(c, "{XB"), strings.HasPrefix(c, "{RB"):
			syms = append(syms, c)
		case strings.HasPrefix(c, "{XS;I,"):
			labels, fields = append(labels, bm), append(fields, syms)
			issues = append(issues, c)
		case strings.HasPrefix(c, "{AY;"):
			issues = append(issues, c)
		default:
			return nil, nil, nil, fmt.Errorf("a command %.30q", c)
		}
	}
	return labels, fields, issues, nil
}

func TestTPCL(t *testing.T) {
	d, v := testDoc(t, 4*72, 3*72, 2)
	p, _ := LookupProfile("b-ex4t1-203")
	var buf bytes.Buffer
	res, err := Write(&buf, d, v, &Options{Profile: p, Copies: 3, Settings: Settings{Darkness: -30, Speed: 25400, Cut: CutPage, Tracking: "mark"}})
	if err != nil {
		t.Fatal(err)
	}
	labels, fields, issues, err := readTPCL(buf.String(), p.DPI)
	if err != nil {
		t.Fatal(err)
	}
	// TPCL prints the QR code and draws Code 128
	if len(labels) != 2 || res.Pages != 6 || res.Symbols != 2 {
		t.Fatalf("%d labels, %+v", len(labels), res)
	}
	rd := newRender(d, v, &Options{Profile: p}, p)
	for i, bm := range labels {
		l, _ := rd.labelOf(i)
		// the size in tenths of a millimetre comes back within a dot
		if bm.W != l.W || bm.H != l.H {
			t.Fatalf("label %d: %d × %d, want %d × %d", i+1, bm.W, bm.H, l.W, l.H)
		}
		sameDots(t, fmt.Sprintf("label %d", i+1), bm, wantDots(t, d, v, i, l, []string{"qr"}))
		if len(fields[i]) != 2 || !regexp.MustCompile(`^\{XB00;\d{4},\d{4},T,M,\d\d,A,0,M2$`).MatchString(fields[i][0]) || fields[i][1] != "{RB00;https://example.com/r/42" {
			t.Errorf("label %d: fields %q", i+1, fields[i])
		}
	}
	want := []string{"{AY;-03,0", "{XS;I,0003,0011CA100", "{XS;I,0003,0011CA100"}
	if fmt.Sprint(issues) != fmt.Sprint(want) {
		t.Errorf("%q, want %q", issues, want)
	}
	if !strings.HasPrefix(buf.String(), "{D0791,1040,0761|}\n") {
		t.Errorf("the size: %.30q", buf.String())
	}
}

// readESCP prints what writeESCP writes: the bit images of 24 pins at
// the place the feeds take the paper to, a page from form feed to form
// feed, and the page length and skip.
func readESCP(data []byte, width int) (pages []*Bitmap, setup []string, err error) {
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
			setup = append(setup, "1/6")
			i += 2
		case bytes.HasPrefix(data[i:], []byte("\x1bC\x00")):
			setup = append(setup, fmt.Sprintf("length %din", data[i+3]))
			i += 4
		case bytes.HasPrefix(data[i:], []byte("\x1bC")):
			setup = append(setup, fmt.Sprintf("length %d", data[i+2]))
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1bN")):
			setup = append(setup, fmt.Sprintf("skip %d", data[i+2]))
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1bJ")):
			y += int(data[i+2])
			i += 3
		case bytes.HasPrefix(data[i:], []byte("\x1b*\x27")):
			cols := int(data[i+3]) | int(data[i+4])<<8
			i += 5
			if cols > width || i+3*cols > len(data) {
				return nil, nil, fmt.Errorf("ESC * of %d columns", cols)
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
			pages = append(pages, bm)
			rows, y = nil, 0
			i++
		default:
			return nil, nil, fmt.Errorf("a byte 0x%02x at %d", data[i], i)
		}
	}
	return pages, setup, nil
}

func TestESCP(t *testing.T) {
	// a form of 10 × 5.5 inches, as wide as 80 columns
	d, v := testDoc(t, 8*72, 5.5*72, 2)
	p, _ := LookupProfile("escp-80")
	var buf bytes.Buffer
	res, err := Write(&buf, d, v, &Options{Profile: p, Copies: 2, Settings: Settings{PerforationSkip: 24, Cut: CutPage}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Pages != 4 || fmt.Sprint(res.Unsupported) != "[bdf-cut=page]" {
		t.Errorf("%+v", res)
	}
	pages, setup, err := readESCP(buf.Bytes(), p.Width)
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 4 {
		t.Fatalf("%d pages", len(pages))
	}
	if fmt.Sprint(setup) != "[1/6 length 33 1/6 skip 2]" {
		t.Errorf("setup %q", setup)
	}
	rd := newRender(d, v, &Options{Profile: p}, p)
	for i, pg := range pages {
		l, _ := rd.labelOf(i / 2)
		want := wantDots(t, d, v, i/2, l, nil)
		// the rows past the last band with dots are not sent
		got := NewBitmap(want.W, want.H)
		copy(got.Bits, pg.Bits)
		if pg.H > want.H+23 {
			t.Fatalf("page %d: %d rows of %d", i+1, pg.H, want.H)
		}
		sameDots(t, fmt.Sprintf("page %d", i+1), got, want)
	}
	if escpLength(11*180, 0) != "\x1bC\x00\x0b" {
		t.Error("a form of 11 inches")
	}
}

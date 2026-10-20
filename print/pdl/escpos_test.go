package pdl

import (
	"bytes"
	"fmt"
	"testing"
)

// escposPage is a receipt as a reader of ESC/POS prints it: its rows of
// dots, and whether it was cut after it.
type escposPage struct {
	dots *Bitmap
	cut  string
}

// readESCPOS prints the commands that writeESCPOS writes: the rows of GS v
// 0 and the feeds of ESC J, a receipt ending at each cut. drawer is the
// pin of ESC p (0 for none).
func readESCPOS(data []byte, width, dpi int) (pages []escposPage, drawer int, err error) {
	var rows [][]byte
	stride := (width + 7) / 8
	end := func(cut string) {
		bm := NewBitmap(width, len(rows))
		for y, r := range rows {
			copy(bm.Row(y), r)
		}
		pages = append(pages, escposPage{bm, cut})
		rows = nil
	}
	for i := 0; i < len(data); {
		switch {
		case bytes.HasPrefix(data[i:], []byte("\x1b@")):
			i += 2
		case bytes.HasPrefix(data[i:], []byte("\x1dP")) && i+3 < len(data):
			if int(data[i+2]) != dpi || int(data[i+3]) != dpi {
				return nil, 0, fmt.Errorf("GS P %d %d at %d dpi", data[i+2], data[i+3], dpi)
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
				return nil, 0, fmt.Errorf("GS v 0 of %d × %d bytes at %d", x, y, i)
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
			drawer = int(data[i+2]) + 1
			i += 5
		default:
			return nil, 0, fmt.Errorf("a byte 0x%02x at %d", data[i], i)
		}
	}
	if len(rows) > 0 {
		end("")
	}
	return pages, drawer, nil
}

func TestESCPOS(t *testing.T) {
	d, v := testDoc(t, 72*80/25.4, 260, 2)
	for _, c := range []struct {
		profile string
		s       Settings
		copies  int
		cuts    []string
	}{
		{"pos-80", Settings{}, 1, []string{"partial", "partial"}},
		{"tm-t88", Settings{Cut: CutJob, Drawer: 2}, 2, []string{"partial"}},
		{"pos-58", Settings{Cut: CutNone, Drawer: 1}, 1, []string{""}},
	} {
		p, _ := LookupProfile(c.profile)
		var buf bytes.Buffer
		res, err := Write(&buf, d, v, &Options{Profile: p, Settings: c.s, Copies: c.copies, Collate: true})
		if err != nil {
			t.Fatal(c.profile, err)
		}
		if res.Pages != 2*c.copies || res.Symbols != 0 || len(res.Unsupported) != 0 {
			t.Errorf("%s: %+v", c.profile, res)
		}
		pages, drawer, err := readESCPOS(buf.Bytes(), p.Width, p.DPI)
		if err != nil {
			t.Fatal(c.profile, err)
		}
		if drawer != c.s.Drawer {
			t.Errorf("%s: the drawer of pin %d, want %d", c.profile, drawer, c.s.Drawer)
		}
		// receipts without a cut between them run into one
		var cuts []string
		for _, pg := range pages {
			cuts = append(cuts, pg.cut)
		}
		if fmt.Sprint(cuts) != fmt.Sprint(c.cuts) {
			t.Fatalf("%s: cuts %v, want %v", c.profile, cuts, c.cuts)
		}
		rd := newRender(d, v, &Options{Profile: p}, p)
		l, err := rd.labelOf(0)
		if err != nil {
			t.Fatal(err)
		}
		if l.W != p.Width {
			t.Errorf("%s: %d dots wide, want %d", c.profile, l.W, p.Width)
		}
		if len(c.cuts) == 2*c.copies {
			for i, pg := range pages {
				sameDots(t, fmt.Sprintf("%s receipt %d", c.profile, i+1), pg.dots, wantDots(t, d, v, i%2, l, nil))
			}
		} else if total := pages[0].dots.H; total != 2*c.copies*l.H {
			t.Errorf("%s: %d rows, want %d", c.profile, total, 2*c.copies*l.H)
		}
	}
}

func TestESCPOSUnsupported(t *testing.T) {
	d, v := testDoc(t, 200, 200, 1)
	p, _ := LookupProfile("pos-80")
	res, err := Write(new(bytes.Buffer), d, v, &Options{Profile: p, Settings: Settings{Darkness: 10, Tracking: "gap", PerforationSkip: 12}})
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(res.Unsupported) != "[print-darkness=10 media-tracking=gap bdf-perforation-skip=12]" {
		t.Errorf("%v", res.Unsupported)
	}
}

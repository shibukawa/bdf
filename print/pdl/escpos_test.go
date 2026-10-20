package pdl

import (
	"bytes"
	"fmt"
	"testing"
)

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
		printed, err := ReadESCPOS(buf.Bytes(), p.Width, p.DPI)
		if err != nil {
			t.Fatal(c.profile, err)
		}
		pages, drawer := printed.Pages, printed.Drawer
		if drawer != c.s.Drawer {
			t.Errorf("%s: the drawer of pin %d, want %d", c.profile, drawer, c.s.Drawer)
		}
		// receipts without a cut between them run into one
		var cuts []string
		for _, pg := range pages {
			cuts = append(cuts, pg.Cut)
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
				sameDots(t, fmt.Sprintf("%s receipt %d", c.profile, i+1), pg.Dots, wantDots(t, d, v, i%2, l, nil))
			}
		} else if total := pages[0].Dots.H; total != 2*c.copies*l.H {
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

package pdl

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestZPLRows(t *testing.T) {
	for _, row := range [][]byte{
		{0, 0, 0, 0}, {0xff, 0xff, 0xff}, {0x12, 0x34, 0, 0}, {0xf0, 0xff, 0xff},
		bytes.Repeat([]byte{0xaa}, 300), append(bytes.Repeat([]byte{0x11}, 250), 0x80),
	} {
		var b strings.Builder
		zplRow(&b, row)
		bm, err := ReadGFA(fmt.Sprintf("%d,%d,%d,%s", len(row), len(row), len(row), b.String()), len(row)*8, 1)
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
		printed, err := ReadZPL(buf.String())
		if err != nil {
			t.Fatal(c.name, err)
		}
		labels := printed.Pages
		if len(labels) != c.labels || res.Pages != 2*o.Copies {
			t.Fatalf("%s: %d labels, %+v", c.name, len(labels), res)
		}
		rd := newRender(d, v, &o, p)
		for i, l := range labels {
			page := i % 2
			lb, _ := rd.labelOf(page)
			sameDots(t, fmt.Sprintf("%s label %d", c.name, i+1), l.Dots, wantDots(t, d, v, page, lb, map[bool][]string{true: {"code128", "qr"}}[c.symbols]))
			if c.symbols {
				if len(l.Symbols) != 2 || !strings.Contains(l.Symbols[0], "^BCN,") || !strings.Contains(l.Symbols[0], "^FDBDF-0042") ||
					!strings.Contains(l.Symbols[1], "^BQN,2,") && !strings.Contains(l.Symbols[1], ",M") || !strings.Contains(l.Symbols[1], "^FDMA,https://example.com/r/42") {
					t.Errorf("%s label %d: fields %q", c.name, i+1, l.Symbols)
				}
			} else if len(l.Symbols) != 0 {
				t.Errorf("%s label %d: fields %q", c.name, i+1, l.Symbols)
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

package pdl

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestSBPL(t *testing.T) {
	d, v := testDoc(t, 4*72, 3*72, 2)
	p, _ := LookupProfile("cl4nx-203")
	var buf bytes.Buffer
	res, err := Write(&buf, d, v, &Options{Profile: p, Copies: 2, Settings: Settings{Darkness: 100, Speed: 15240, Cut: CutJob}})
	if err != nil {
		t.Fatal(err)
	}
	printed, err := ReadSBPL(buf.String())
	if err != nil {
		t.Fatal(err)
	}
	labels, rest := printed.Pages, printed.Commands
	if len(labels) != 2 || res.Pages != 4 || res.Symbols != 4 {
		t.Fatalf("%d labels, %+v", len(labels), res)
	}
	rd := newRender(d, v, &Options{Profile: p}, p)
	for i, pg := range labels {
		l, _ := rd.labelOf(i)
		sameDots(t, fmt.Sprintf("label %d", i+1), pg.Dots, wantDots(t, d, v, i, l, []string{"code128", "qr"}))
		if len(pg.Symbols) != 2 || !regexp.MustCompile(`^@\d+,\d+ BG\d\d\d\d\d>HBDF-0042$`).MatchString(pg.Symbols[0]) ||
			!regexp.MustCompile(`^@\d+,\d+ 2D30,M,\d\d,1,0 DN0024,https://example.com/r/42$`).MatchString(pg.Symbols[1]) {
			t.Errorf("label %d: symbols %q", i+1, pg.Symbols)
		}
	}
	want := []string{"#E5A", "CS06", "Q2", "#E5A", "CS06", "CT0002", "Q2"}
	if fmt.Sprint(rest) != fmt.Sprint(want) {
		t.Errorf("commands %q, want %q", rest, want)
	}
}

func TestTPCL(t *testing.T) {
	d, v := testDoc(t, 4*72, 3*72, 2)
	p, _ := LookupProfile("b-ex4t1-203")
	var buf bytes.Buffer
	res, err := Write(&buf, d, v, &Options{Profile: p, Copies: 3, Settings: Settings{Darkness: -30, Speed: 25400, Cut: CutPage, Tracking: "mark"}})
	if err != nil {
		t.Fatal(err)
	}
	printed, err := ReadTPCL(buf.String(), p.DPI)
	if err != nil {
		t.Fatal(err)
	}
	labels, issues := printed.Pages, printed.Commands
	// TPCL prints the QR code and draws Code 128
	if len(labels) != 2 || res.Pages != 6 || res.Symbols != 2 {
		t.Fatalf("%d labels, %+v", len(labels), res)
	}
	rd := newRender(d, v, &Options{Profile: p}, p)
	for i, pg := range labels {
		bm := pg.Dots
		l, _ := rd.labelOf(i)
		// the size in tenths of a millimetre comes back within a dot
		if bm.W != l.W || bm.H != l.H {
			t.Fatalf("label %d: %d × %d, want %d × %d", i+1, bm.W, bm.H, l.W, l.H)
		}
		sameDots(t, fmt.Sprintf("label %d", i+1), bm, wantDots(t, d, v, i, l, []string{"qr"}))
		if len(pg.Symbols) != 2 || !regexp.MustCompile(`^\{XB00;\d{4},\d{4},T,M,\d\d,A,0,M2$`).MatchString(pg.Symbols[0]) || pg.Symbols[1] != "{RB00;https://example.com/r/42" {
			t.Errorf("label %d: fields %q", i+1, pg.Symbols)
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
	printed, err := ReadESCP(buf.Bytes(), p.Width)
	if err != nil {
		t.Fatal(err)
	}
	pages, setup := printed.Pages, printed.Commands
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
		copy(got.Bits, pg.Dots.Bits)
		if pg.Dots.H > want.H+23 {
			t.Fatalf("page %d: %d rows of %d", i+1, pg.Dots.H, want.H)
		}
		sameDots(t, fmt.Sprintf("page %d", i+1), got, want)
	}
	if escpLength(11*180, 0) != "\x1bC\x00\x0b" {
		t.Error("a form of 11 inches")
	}
}

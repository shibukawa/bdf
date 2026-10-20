package print

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
	"github.com/shibukawa/bdf/print/pdl"
	"github.com/shibukawa/bdf/print/raw"
	"github.com/shibukawa/bdf/print/raw/rawtest"
)

// TestPrintRawLabels prints labels of ZPL over 9100: the profile of the
// URL, the copies and the settings of the ticket in the file, what ZPL
// has no command for left out.
func TestPrintRawLabels(t *testing.T) {
	s := rawtest.NewSocket()
	defer s.Close()
	p, err := raw.Parse(s.URL + "?lang=zpl&profile=zd421-300")
	if err != nil {
		t.Fatal(err)
	}
	tk := Ticket{Copies: 2, Collate: "uncollated", Media: "custom_100x150mm", Darkness: 20, MediaTracking: "gap", Cut: "page", CashDrawer: 1, Sides: "two-sided-long-edge", JobName: "labels"}
	var calls []string
	res, err := PrintRaw(context.Background(), p, testDoc(2), &Options{Ticket: tk,
		OnCreated: func(*ipp.JobStatus) error { calls = append(calls, "begin"); return nil },
		OnSent:    func(*ipp.JobStatus) error { calls = append(calls, "sent"); return nil }})
	if err != nil {
		t.Fatal(err)
	}
	jobs := s.Jobs()
	if len(jobs) != 1 || res.Format != "ZPL" || res.Sheets != 4 || strings.Join(calls, " ") != "begin sent" {
		t.Fatalf("%d jobs, %+v, %v", len(jobs), res, calls)
	}
	zpl := string(jobs[0].Data)
	// 100 mm of 300 dpi, 150 mm long; the A4 page fitted to it
	for _, want := range []string{"^XA\n^MNY\n^MD6\n^PW1181^LL1771", "^MMC\n", "^GFA,", "^PQ2\n^XZ\n"} {
		if !strings.Contains(zpl, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	if strings.Count(zpl, "^XA") != 2 || strings.Contains(zpl, "@PJL") {
		t.Errorf("%d formats", strings.Count(zpl, "^XA"))
	}
	if strings.Join(res.Unsupported, ";") != "sides=two-sided-long-edge;bdf-cash-drawer=1" {
		t.Errorf("unsupported %q", res.Unsupported)
	}
	// the language and pdf are not mixed
	if _, err := PrintRaw(context.Background(), p, testDoc(1), &Options{Format: "pdf"}); err == nil {
		t.Error("a ZPL printer sent PDF")
	}
	// the options' language over a printer of PostScript
	ps, _ := raw.Parse(s.URL)
	if res, err := PrintRaw(context.Background(), ps, testDoc(1), &Options{Format: "escpos"}); err != nil || res.Format != "ESC/POS" {
		t.Errorf("escpos: %+v %v", res, err)
	}
	if !bytes.HasPrefix(s.Jobs()[1].Data, []byte("\x1b@\x1dP")) {
		t.Errorf("ESC/POS begins %q", s.Jobs()[1].Data[:8])
	}
	// the language of the URL's profile
	tm, _ := raw.Parse(s.URL + "?profile=tm-t88")
	if res, err := PrintRaw(context.Background(), tm, testDoc(1), nil); err != nil || res.Format != "ESC/POS" {
		t.Errorf("tm-t88: %+v %v", res, err)
	}
	if _, err := PrintRaw(context.Background(), p, testDoc(1), &Options{Ticket: Ticket{PrintMode: "text"}}); err == nil {
		t.Error("the text mode")
	}
}

// TestPrintRawText prints a form of a template of text over 9100: ESC/P in
// the text mode by default, ESC/POS when the ticket asks it.
func TestPrintRawText(t *testing.T) {
	f, err := os.Open("../testdata/form/text-slip.bdf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fi, _ := f.Stat()
	r, err := bdf.OpenSingle(f, fi.Size())
	if err != nil {
		t.Fatal(err)
	}
	doc, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	s := rawtest.NewSocket()
	defer s.Close()
	vp, _ := raw.Parse(s.URL + "?profile=vp-f2000")
	res, err := PrintRaw(context.Background(), vp, doc, &Options{SkipRoles: []string{bdf.RoleBackground}})
	if err != nil || res.Format != "ESC/P (text)" || res.Sheets != 1 {
		t.Fatalf("%+v %v", res, err)
	}
	data := s.Jobs()[0].Data
	if !bytes.Contains(data, []byte("\x1bC\x21\x1bN\x02")) || !bytes.Contains(data, []byte("\x1c&")) || bytes.Contains(data, []byte("\x1b*")) {
		t.Errorf("ESC/P of the slip: %q", data[:min(len(data), 64)])
	}
	// the names printed on the paper are left out: 得意先 in JIS
	if bytes.Contains(data, []byte("\x46\x40\x30\x55\x40\x68")) {
		t.Error("the preprinted paper was printed")
	}
	if res, err := PrintRaw(context.Background(), vp, doc, &Options{Ticket: Ticket{PrintMode: "raster"}}); err != nil || res.Format != "ESC/P" {
		t.Errorf("the raster mode: %+v %v", res, err)
	}
	pos, _ := raw.Parse(s.URL + "?profile=tm-t88")
	if _, err := PrintRaw(context.Background(), pos, doc, &Options{Ticket: Ticket{PrintMode: "text"}}); err == nil || !strings.Contains(err.Error(), "wide") {
		t.Errorf("80 columns on a receipt: %v", err)
	}
}

// TestPrintCUPSRaw sends ESC/POS to an IPP queue that takes data of the
// printer's own language, as a raw queue of CUPS does.
func TestPrintCUPSRaw(t *testing.T) {
	p := ipptest.New()
	p.Formats = []string{"application/vnd.cups-raw", "application/octet-stream"}
	c := serve(t, p)
	prof, _ := pdl.LookupProfile("tm-m30")
	res, err := Print(context.Background(), c, testDoc(1), &Options{Format: "escpos", Profile: &prof, Ticket: Ticket{CashDrawer: 2, Darkness: 10}})
	if err != nil {
		t.Fatal(err)
	}
	jobs := p.Jobs()
	if len(jobs) != 1 || jobs[0].Format != FormatCUPSRaw || res.Format != "ESC/POS" || res.Sheets != 1 {
		t.Fatalf("%d jobs, %+v", len(jobs), res)
	}
	data := jobs[0].Documents[0]
	if !bytes.HasPrefix(data, []byte("\x1b@\x1dP\xcb\xcb")) || !bytes.HasSuffix(data, []byte("\x1dVB\x00\x1bp\x01\x19\xfa")) {
		t.Errorf("the receipt: %q … %q", data[:6], data[len(data)-10:])
	}
	if strings.Join(res.Unsupported, ";") != "print-darkness=10" {
		t.Errorf("unsupported %q", res.Unsupported)
	}
	// a printer of PDF alone does not take it
	p.Formats = []string{"application/pdf"}
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Format: "zpl"}); err == nil || !strings.Contains(err.Error(), "raw queue") {
		t.Errorf("a printer of PDF: %v", err)
	}
	// nor a profile of another language
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Format: "zpl", Profile: &prof}); err == nil {
		t.Error("a profile of ESC/POS for ZPL")
	}
}

// TestDeviceTicket checks the settings of a ticket and that an IPP
// printer of pages is told what it supports of them.
func TestDeviceTicket(t *testing.T) {
	for _, bad := range []Ticket{{Darkness: 200}, {Cut: "half"}, {CashDrawer: 3}, {MediaTracking: "holes"}, {PrintMode: "graphics"}, {LabelGap: -1}} {
		if err := bad.Check(); err == nil {
			t.Errorf("%+v: no error", bad)
		}
	}
	m := Ticket{Darkness: 5, Cut: "job"}.Merge(Ticket{Darkness: 9, Speed: 5080, CashDrawer: 1})
	if m.Darkness != 5 || m.Speed != 5080 || m.Cut != "job" || m.CashDrawer != 1 {
		t.Errorf("merged %+v", m)
	}
	pr := &ipp.Printer{}
	a := jobAttributes(&Ticket{Darkness: 5, Speed: 5080, Cut: "page"}, pr, FormatPDF, "", true)
	if strings.Join(a.unsupported, ";") != "print-darkness=5;print-speed=5080;bdf-cut=page" {
		t.Errorf("unsupported %q", a.unsupported)
	}
	if !IsDeviceFormat("ZPL") || IsDeviceFormat("pdf") || deviceLanguage("esc/pos") != pdl.ESCPOS {
		t.Error("IsDeviceFormat")
	}
	if _, err := ResolveMedia("104x150mm"); err != nil {
		t.Error(err)
	}
}

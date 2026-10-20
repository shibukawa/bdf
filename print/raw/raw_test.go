package raw_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/raw"
	"github.com/shibukawa/bdf/print/raw/rawtest"
)

func TestParse(t *testing.T) {
	for _, c := range []struct {
		url  string
		want raw.Printer
	}{
		{"socket://printer.local", raw.Printer{URL: "socket://printer.local", Scheme: "socket", Addr: "printer.local:9100", Format: raw.FormatPostScript}},
		{"SOCKET://10.0.0.5:9101?format=pdf&pjl=false", raw.Printer{URL: "socket://10.0.0.5:9101", Scheme: "socket", Addr: "10.0.0.5:9101", Format: raw.FormatPDF, PJL: "false"}},
		{"lpd://[::1]/lp1?user=alice&format=postscript", raw.Printer{URL: "lpd://[::1]/lp1", Scheme: "lpd", Addr: "[::1]:515", Queue: "lp1", Format: raw.FormatPostScript, User: "alice"}},
	} {
		p, err := raw.Parse(c.url)
		if err != nil || *p != c.want {
			t.Errorf("%s: %+v %v, want %+v", c.url, p, err, c.want)
		}
	}
	for _, bad := range []string{"ipp://host/", "socket://", "socket://host/path", "lpd://host", "lpd://host/", "socket://host?format=pcl", "socket://host?pjl=yes", "socket://host?copies=2"} {
		if p, err := raw.Parse(bad); err == nil {
			t.Errorf("%s: %+v", bad, p)
		}
	}
	if !raw.IsURL("Socket://x") || !raw.IsURL("lpd://x/q") || raw.IsURL("ipp://x") {
		t.Error("IsURL")
	}
	p, _ := raw.Parse("socket://h")
	if p.UsesPJL(raw.FormatPostScript) || !p.UsesPJL(raw.FormatPDF) {
		t.Error("PJL is for PDF by default")
	}
	p, _ = raw.Parse("socket://h?pjl=true")
	if !p.UsesPJL(raw.FormatPostScript) {
		t.Error("pjl=true")
	}
}

func TestSocket(t *testing.T) {
	s := rawtest.NewSocket()
	defer s.Close()
	s.Back("@PJL USTATUS JOB\r\n")
	p, err := raw.Parse(s.URL + "?format=pdf")
	if err != nil {
		t.Fatal(err)
	}
	begun := 0
	doc := []byte("%PDF-1.7\n…\n%%EOF\n")
	n, err := p.Send(context.Background(), &raw.Job{Name: "Report \"Q3\"", Format: raw.FormatPDF, Data: doc, PJL: []string{"@PJL SET DUPLEX=ON"},
		Begin: func() error { begun++; return nil }})
	if err != nil || begun != 1 {
		t.Fatalf("%v, begun %d", err, begun)
	}
	jobs := s.Jobs()
	if len(jobs) != 1 || int64(len(jobs[0].Data)) != n {
		t.Fatalf("%d jobs, %d bytes sent", len(jobs), n)
	}
	lines, got, ok := rawtest.PJL(jobs[0].Data)
	if !ok || !bytes.Equal(got, doc) || strings.Join(lines, "|") != `@PJL|@PJL JOB NAME="Report _Q3_"|@PJL SET DUPLEX=ON|@PJL ENTER LANGUAGE=PDF` {
		t.Errorf("PJL %q, document %q", lines, got)
	}
	if !bytes.HasSuffix(jobs[0].Data, []byte("\x1b%-12345X@PJL EOJ NAME=\"Report _Q3_\"\r\n\x1b%-12345X")) {
		t.Errorf("the end of the job: %q", jobs[0].Data[len(jobs[0].Data)-40:])
	}

	// PostScript as it is, and a check that sends nothing
	p, _ = raw.Parse(s.URL)
	if _, err := p.Send(context.Background(), &raw.Job{Name: "x", Format: raw.FormatPostScript, Data: []byte("%!PS\nshowpage\n")}); err != nil {
		t.Fatal(err)
	}
	if err := p.Check(context.Background()); err != nil {
		t.Fatal(err)
	}
	if jobs := s.Jobs(); len(jobs) != 2 || string(jobs[1].Data) != "%!PS\nshowpage\n" {
		t.Fatalf("%d jobs: %q", len(jobs), jobs[len(jobs)-1].Data)
	}

	// a connection cut in the document: the printer may have printed a part
	s.CutAfter(1000)
	_, err = p.Send(context.Background(), &raw.Job{Format: raw.FormatPostScript, Data: make([]byte, 32<<20)})
	if !errors.Is(err, raw.ErrPartial) {
		t.Errorf("cut: %v", err)
	}
	// a printer that is not there: nothing was sent
	s.Close()
	begun = 0
	_, err = p.Send(context.Background(), &raw.Job{Format: raw.FormatPostScript, Data: []byte("x"), Begin: func() error { begun++; return nil }})
	if err == nil || errors.Is(err, raw.ErrPartial) || begun != 0 {
		t.Errorf("no printer: %v, begun %d", err, begun)
	}
	if err := p.Check(context.Background()); err == nil {
		t.Error("a printer that is not there is checked")
	}
}

func TestLPD(t *testing.T) {
	s := rawtest.NewLPD("lp1")
	defer s.Close()
	p, err := raw.Parse(s.URL + "?user=alice")
	if err != nil {
		t.Fatal(err)
	}
	begun := 0
	doc := []byte("%!PS-Adobe-3.0\nshowpage\n")
	n, err := p.Send(context.Background(), &raw.Job{Name: "a b.ps", Format: raw.FormatPostScript, Data: doc, Begin: func() error { begun++; return nil }})
	if err != nil || begun != 1 || n != int64(len(doc)) {
		t.Fatalf("%v, begun %d, %d bytes", err, begun, n)
	}
	jobs := s.Jobs()
	if len(jobs) != 1 || jobs[0].Queue != "lp1" || !bytes.Equal(jobs[0].Data, doc) {
		t.Fatalf("%+v", jobs)
	}
	c := jobs[0].Control
	for _, want := range []string{"\nPalice\n", "\nJa_b.ps\n", "\nNa_b.ps\n", "\nldfA", "\nUdfA"} {
		if !strings.Contains("\n"+c, want) {
			t.Errorf("the control file lacks %q: %q", want, c)
		}
	}

	// the queue refused: nothing was sent
	s.Refuse(true)
	begun = 0
	if _, err := p.Send(context.Background(), &raw.Job{Format: raw.FormatPostScript, Data: doc, Begin: func() error { begun++; return nil }}); !errors.Is(err, raw.ErrRefused) || begun != 0 {
		t.Errorf("refused: %v, begun %d", err, begun)
	}
	s.Refuse(false)
	// cut before the control file: the printer prints nothing, and the job
	// may be sent again
	s.CutAfter(1)
	if _, err := p.Send(context.Background(), &raw.Job{Format: raw.FormatPostScript, Data: doc, Begin: func() error { begun++; return nil }}); err == nil || errors.Is(err, raw.ErrPartial) || begun != 0 {
		t.Errorf("cut before the control file: %v, begun %d", err, begun)
	}
	// cut in the control file: the printer may print it
	s.CutControl()
	if _, err := p.Send(context.Background(), &raw.Job{Format: raw.FormatPostScript, Data: doc, Begin: func() error { begun++; return nil }}); !errors.Is(err, raw.ErrPartial) || begun != 1 {
		t.Errorf("cut in the control file: %v, begun %d", err, begun)
	}
	if len(s.Jobs()) != 1 {
		t.Errorf("%d jobs", len(s.Jobs()))
	}
}

func TestCanceled(t *testing.T) {
	s := rawtest.NewSocket()
	defer s.Close()
	p, _ := raw.Parse(s.URL)
	ctx, cancel := context.WithCancel(context.Background())
	_, err := p.Send(ctx, &raw.Job{Format: raw.FormatPostScript, Data: make([]byte, 64<<20), Begin: func() error { cancel(); return nil }})
	if !errors.Is(err, raw.ErrPartial) || !errors.Is(err, context.Canceled) {
		t.Errorf("%v", err)
	}
}

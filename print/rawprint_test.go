package print

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
	"github.com/shibukawa/bdf/print/raw/rawtest"
)

// TestPrintRawPostScript prints PostScript over 9100: the ticket is in the
// file.
func TestPrintRawPostScript(t *testing.T) {
	s := rawtest.NewSocket()
	defer s.Close()
	p, err := raw.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	var calls []string
	tk := Ticket{Copies: 2, Sides: "two-sided-short-edge", Media: "iso_a5_148x210mm", ColorMode: "monochrome", MediaType: "labels", OutputBin: "top", JobName: "doc"}
	res, err := PrintRaw(context.Background(), p, testDoc(3), &Options{Ticket: tk,
		Progress:  func(e Event) { calls = append(calls, e.Kind) },
		OnCreated: func(*ipp.JobStatus) error { calls = append(calls, "begin"); return nil },
		OnSent:    func(*ipp.JobStatus) error { calls = append(calls, "on-sent"); return nil }})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(calls, " ") != "begin sent on-sent" {
		t.Errorf("calls %q", calls)
	}
	jobs := s.Jobs()
	if len(jobs) != 1 || res.Format != FormatPostScript || res.Sheets != 3 || res.Bytes != int64(len(jobs[0].Data)) || res.Job != nil {
		t.Fatalf("%d jobs, %+v", len(jobs), res)
	}
	ps := string(jobs[0].Data)
	for _, want := range []string{"%!PS-Adobe-3.0\n", "%%Title: (doc)\n", "%%Pages: 3\n", "<< /PageSize [419.5276 595.2756] >>", "<< /NumCopies 2 >>", "<< /Collate true >>",
		"<< /Duplex true /Tumble true >>", "<< /MediaType (labels) >>", "\n0 g\n"} {
		if !strings.Contains(ps, want) {
			t.Errorf("the file lacks %q", want)
		}
	}
	if strings.Join(res.Unsupported, ";") != "output-bin=top" {
		t.Errorf("unsupported %q", res.Unsupported)
	}
	if _, err := PrintRaw(context.Background(), p, testDoc(1), &Options{Ticket: tk, Strict: true}); !errors.Is(err, ErrUnsupported) {
		t.Errorf("strict: %v", err)
	}
	if _, err := PrintRaw(context.Background(), p, testDoc(1), &Options{Format: "pwg-raster"}); err == nil || !strings.Contains(err.Error(), "PostScript or PDF, not PWG Raster") {
		t.Errorf("raster: %v", err)
	}
	// a page at its own size without a media
	if _, err := PrintRaw(context.Background(), p, testDoc(1), nil); err != nil {
		t.Fatal(err)
	}
	if ps := string(s.Jobs()[1].Data); !strings.Contains(ps, "<< /PageSize [595.276 841.89] >>") || strings.Contains(ps, "Duplex") || !strings.Contains(ps, "\nsave\n1 0 0 -1 0 841.89 cm\n") {
		t.Errorf("the page at its size: %q", ps[strings.Index(ps, "%%BeginSetup"):][:300])
	}
	// cut while it is sent
	s.CutAfter(100)
	if _, err := PrintRaw(context.Background(), p, testDoc(40), &Options{Ticket: Ticket{Copies: 30}}); !errors.Is(err, raw.ErrPartial) {
		t.Errorf("cut: %v", err)
	}
}

// TestPrintRawPDF prints PDF over LPD, with a PJL header.
func TestPrintRawPDF(t *testing.T) {
	s := rawtest.NewLPD("pdfq")
	defer s.Close()
	p, err := raw.Parse(s.URL + "?format=pdf")
	if err != nil {
		t.Fatal(err)
	}
	tk := Ticket{Copies: 2, Sides: "two-sided-long-edge", Media: "na_letter_8.5x11in", ColorMode: "monochrome", MediaType: "labels", Scaling: "fit"}
	res, err := PrintRaw(context.Background(), p, testDoc(2), &Options{Ticket: tk})
	if err != nil {
		t.Fatal(err)
	}
	jobs := s.Jobs()
	if len(jobs) != 1 || res.Format != FormatPDF || res.Sheets != 4 {
		t.Fatalf("%d jobs, %+v", len(jobs), res)
	}
	lines, pdf, ok := rawtest.PJL(jobs[0].Data)
	if !ok || !strings.HasPrefix(string(pdf), "%PDF-1.7\n") {
		t.Fatalf("PJL %q", lines)
	}
	if got := strings.Join(lines, "|"); got != `@PJL|@PJL JOB NAME="document"|@PJL SET DUPLEX=ON|@PJL SET BINDING=LONGEDGE|@PJL SET PAPER=LETTER|@PJL SET RENDERMODE=GRAYSCALE|@PJL ENTER LANGUAGE=PDF` {
		t.Errorf("PJL %s", got)
	}
	if got := strings.Join(res.Unsupported, ";"); got != "media-type=labels;print-scaling=fit" {
		t.Errorf("unsupported %q", got)
	}
	// without PJL, the sides, the paper and gray are not told
	p.PJL = "false"
	res, err = PrintRaw(context.Background(), p, testDoc(1), &Options{Ticket: Ticket{Sides: "one-sided", Media: "iso_a4_210x297mm", ColorMode: "monochrome"}})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(res.Unsupported, ";"); got != "sides=one-sided;media=iso_a4_210x297mm;print-color-mode=monochrome" || !strings.HasPrefix(string(s.Jobs()[1].Data), "%PDF-1.7\n") {
		t.Errorf("unsupported %q", got)
	}
	// a refused queue
	s.Refuse(true)
	if _, err := PrintRaw(context.Background(), p, testDoc(1), nil); !errors.Is(err, raw.ErrRefused) {
		t.Errorf("refused: %v", err)
	}
}

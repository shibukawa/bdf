package spool

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

// rawFixture is a spool of a 9100 printer (queue office) and an LPD printer
// of PDF (queue lpd).
func rawFixture(t *testing.T) (*fixture, *fakeprinter.Raw, *fakeprinter.Raw) {
	t.Helper()
	f := newFixture(t)
	sock, lpd := fakeprinter.NewSocket(t), fakeprinter.NewLPD(t, "pdf")
	t.Cleanup(sock.Close)
	t.Cleanup(lpd.Close)
	ps, _ := raw.Parse(sock.URL)
	pl, _ := raw.Parse(lpd.URL + "?format=pdf")
	f.o.Printers = []Printer{{Name: "socket", Raw: ps}, {Name: "lpd", Raw: pl}}
	f.o.Queues = []Queue{{Name: "office", Printers: []string{"socket"}}, {Name: "lpd", Printers: []string{"lpd"}}}
	return f, sock, lpd
}

// TestRunRaw prints on printers of 9100 and LPD: a job is done once the
// printer has it.
func TestRunRaw(t *testing.T) {
	fast(t)
	f, sock, lpd := rawFixture(t)
	s := f.running()
	doc := bdfBytes(t, testDoc(2))
	a := submit(t, s, Submission{Name: "a.bdf", Document: bytes.NewReader(doc), Ticket: print.Ticket{Sides: "two-sided-long-edge", OutputBin: "top"}})
	b := submit(t, s, Submission{Queue: "lpd", Name: "b.bdf", Document: bytes.NewReader(doc)})
	for _, id := range []int32{a.ID, b.ID} {
		waitFor(t, "the jobs", func() bool { return jobState(s, id).State == ipp.JobCompleted })
	}
	ja, jb := jobState(s, a.ID), jobState(s, b.ID)
	if ja.Format != print.FormatPostScript || ja.Reasons[0] != "job-completed-with-warnings" || !strings.Contains(strings.Join(ja.Unsupported, ";"), "output-bin=top") ||
		ja.Printer.Name != "socket" || !ja.Printer.Begun || !ja.Printer.Sent || ja.Sheets != 2 {
		t.Errorf("9100: %+v", ja)
	}
	if jb.Format != print.FormatPDF || jb.Reasons[0] != "job-completed-successfully" || jb.Printer.URL != strings.Split(lpd.URL, "?")[0] {
		t.Errorf("LPD: %+v", jb)
	}
	// (what each printer is sent, print's tests of PrintRaw check)
	if js := sock.Jobs(); len(js) != 1 {
		t.Errorf("the 9100 printer has %d jobs", len(js))
	}
	if js := lpd.Jobs(); len(js) != 1 || !strings.Contains(js[0].Control, "J#2.1_b.bdf") {
		t.Errorf("the LPD printer has %d jobs", len(js))
	}
	info := s.Printers()
	if info[0].Name != "lpd" || !info[0].Up || info[0].State != "idle" || info[1].URL != sock.URL {
		t.Errorf("printers %+v", info)
	}
}

// TestRunRawDown sends again to a printer that was not reached, and holds a
// job that was cut while it was sent: the printer may have printed a part
// of it.
func TestRunRawDown(t *testing.T) {
	fast(t)
	f, sock, _ := rawFixture(t)
	sock.Close()
	s := f.running()
	doc := bdfBytes(t, testDoc(1))
	a := submit(t, s, Submission{Document: bytes.NewReader(doc)})
	time.Sleep(200 * time.Millisecond)
	if j := jobState(s, a.ID); j.State != ipp.JobPending || len(sock.Jobs()) != 0 {
		t.Fatalf("a job sent to a printer that is down: %+v", j)
	}
	if info := s.Printers(); info[1].Up {
		t.Errorf("printer up: %+v", info[1])
	}
	if err := sock.Reopen(); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the job once the printer is back", func() bool { return jobState(s, a.ID).State == ipp.JobCompleted })

	sock.CutAfter(100)
	b := submit(t, s, Submission{Document: bytes.NewReader(bdfBytes(t, testDoc(30)))})
	waitFor(t, "the cut job held", func() bool {
		j := jobState(s, b.ID)
		return j.State == ipp.JobPendingHeld && j.Reasons[0] == ReasonDuplicateCheck
	})
	if j := jobState(s, b.ID); !strings.Contains(j.Message, "may have printed a part") || j.Printer != nil || len(sock.Jobs()) != 1 {
		t.Fatalf("held: %+v", j)
	}
	// a person sends it again
	if _, err := s.Release(context.Background(), b.ID, ""); err != nil {
		t.Fatal(err)
	}
	waitFor(t, "the job released", func() bool { return jobState(s, b.ID).State == ipp.JobCompleted })
	if j := jobState(s, b.ID); j.Attempts != 2 || len(sock.Jobs()) != 2 {
		t.Fatalf("released: %d attempts, %d jobs", j.Attempts, len(sock.Jobs()))
	}
}

// TestCrashRaw crashes a spool at the points of a send to a 9100 printer:
// a job the printer may have begun to print is held for a person.
func TestCrashRaw(t *testing.T) {
	fast(t)
	for _, c := range []struct {
		point  string
		state  ipp.JobState
		reason string
		jobs   int
	}{
		{"attempt-recorded", ipp.JobCompleted, "job-completed-successfully", 1},
		{"created-recorded", ipp.JobPendingHeld, ReasonDuplicateCheck, 0},
		{"sent-recorded", ipp.JobCompleted, "job-completed-with-warnings", 1},
	} {
		t.Run(c.point, func(t *testing.T) {
			p := fakeprinter.NewSocket(t)
			defer p.Close()
			path := filepath.Join(t.TempDir(), "spool.db")
			cmd := exec.Command(os.Args[0], "-test.run=^TestCrashChild$")
			cmd.Env = append(os.Environ(), "SPOOL_CRASH_POINT="+c.point, "SPOOL_CRASH_PATH="+path, "SPOOL_CRASH_PRINTER="+p.URL)
			out, err := cmd.CombinedOutput()
			if ee, ok := err.(*exec.ExitError); !ok || ee.ExitCode() != 3 {
				t.Fatalf("the child did not crash at %s: %v\n%s", c.point, err, out)
			}
			s, err := Open(context.Background(), Options{Path: path, Queues: []Queue{{Name: "office", Printers: []string{"mfp1"}}},
				Printers: []Printer{crashPrinter(t, p.URL)}})
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			s.Start()
			waitFor(t, "the job after the crash", func() bool { return jobState(s, 1).State == c.state })
			time.Sleep(100 * time.Millisecond)
			if j := jobState(s, 1); j.Reasons[0] != c.reason || len(p.Jobs()) != c.jobs {
				t.Fatalf("after %s: job %s %v, %d jobs at the printer", c.point, j.State, j.Reasons, len(p.Jobs()))
			}
		})
	}
}

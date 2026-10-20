package print

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
)

func TestHooks(t *testing.T) {
	p := fakeprinter.New()
	c := serve(t, p)
	pr, err := c.GetPrinterAttributes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var created, sent *ipp.JobStatus
	res, err := Print(context.Background(), c, testDoc(1), &Options{
		Printer:   pr,
		OnCreated: func(js *ipp.JobStatus) error { created = js; return nil },
		OnSent:    func(js *ipp.JobStatus) error { sent = js; return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if created == nil || sent == nil || created.ID != res.Job.ID || sent.ID != res.Job.ID {
		t.Fatalf("created %v sent %v result %v", created, sent, res.Job)
	}
	if n := p.Requests["Get-Printer-Attributes"]; n != 1 {
		t.Errorf("Get-Printer-Attributes asked %d times with Options.Printer", n)
	}
	// A record that cannot be written cancels the job before its document.
	bad := errors.New("the spool is full")
	_, err = Print(context.Background(), c, testDoc(1), &Options{
		Printer:   pr,
		OnCreated: func(*ipp.JobStatus) error { return bad },
		OnSent:    func(*ipp.JobStatus) error { t.Error("sent after a failed record"); return nil },
	})
	if !errors.Is(err, bad) {
		t.Fatalf("OnCreated error: %v", err)
	}
	last := p.Jobs()[len(p.Jobs())-1]
	if p.JobState(last) != ipp.JobCanceled || len(last.Documents) != 0 {
		t.Fatalf("job after a failed record: %v with %d documents", p.JobState(last), len(last.Documents))
	}
	// Options.Printer still refuses a stopped printer.
	stopped := *pr
	stopped.State = ipp.PrinterStopped
	stopped.Reasons = []string{"media-empty-error"}
	if _, err := Print(context.Background(), c, testDoc(1), &Options{Printer: &stopped}); !errors.Is(err, ErrPrinterStopped) {
		t.Fatalf("stopped printer: %v", err)
	}
}

func TestHooksPrintJob(t *testing.T) {
	// A printer without Create-Job: the job is known once it has the document.
	p := fakeprinter.New()
	c := serve(t, p)
	pr, err := c.GetPrinterAttributes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	pr.Operations = []ipp.Op{ipp.OpPrintJob, ipp.OpGetJobAttributes, ipp.OpGetPrinterAttributes}
	var order []string
	_, err = Print(context.Background(), c, testDoc(1), &Options{
		Printer:   pr,
		OnCreated: func(*ipp.JobStatus) error { order = append(order, "created"); return nil },
		OnSent:    func(*ipp.JobStatus) error { order = append(order, "sent"); return nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(order, ",") != "created,sent" || p.Requests["Print-Job"] != 1 {
		t.Fatalf("order %v, requests %v", order, p.Requests)
	}
}

func TestCheck(t *testing.T) {
	p := fakeprinter.New()
	p.Processing = time.Hour
	c := serve(t, p)
	res, err := Print(context.Background(), c, testDoc(1), nil)
	if err != nil {
		t.Fatal(err)
	}
	js, err := Check(context.Background(), c, res.Job)
	if err != nil || js.State != ipp.JobProcessing || js.ID != res.Job.ID {
		t.Fatalf("processing: %v %v", js, err)
	}
	p.Abort(res.Job.ID)
	if js, err = Check(context.Background(), c, js); !errors.Is(err, ErrAborted) || js.State != ipp.JobAborted {
		t.Fatalf("aborted: %v %v", js, err)
	}
	p.Forget()
	if js, err = Check(context.Background(), c, js); js != nil || !ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
		t.Fatalf("forgotten: %v %v", js, err)
	}
}

func TestTicketJSON(t *testing.T) {
	tk := Ticket{JobName: "a.bdf", Copies: 2, Sides: "two-sided-long-edge", MediaSource: "tray-1", ColorMode: "monochrome",
		Pages: []int{1, 3}, Scaling: raster.ScaleFit, Orientation: raster.Landscape, Finishings: []int{4}, Priority: 60}
	b, err := json.Marshal(tk)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"job-name":"a.bdf","copies":2,"sides":"two-sided-long-edge","media-source":"tray-1","print-color-mode":"monochrome","pages":[1,3],"print-scaling":"fit","orientation-requested":4,"finishings":[4],"job-priority":60}`
	if string(b) != want {
		t.Fatalf("JSON:\n%s\nwant\n%s", b, want)
	}
	var back Ticket
	if err := json.Unmarshal(b, &back); err != nil || !reflect.DeepEqual(back, tk) {
		t.Fatalf("back %+v %v", back, err)
	}
	q := Ticket{Sides: "one-sided", ColorMode: "color", Media: "iso_a4_210x297mm", Copies: 1}
	m := Ticket{ColorMode: "monochrome"}.Merge(q)
	if m.Sides != "one-sided" || m.ColorMode != "monochrome" || m.Media != "iso_a4_210x297mm" || m.Copies != 1 {
		t.Fatalf("merged %+v", m)
	}
}

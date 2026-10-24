package ipp_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	dipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/printconv"
	"github.com/shibukawa/bdf/print/raster"
)

// pclm is a document of PCLm, the raster of Mopria, of pages A6 pages.
func pclm(t *testing.T, pages int) []byte {
	t.Helper()
	var b bytes.Buffer
	w, err := raster.NewPCLmWriter(&b, raster.PCLmOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for range pages {
		pg := raster.Page{Width: 310, Height: 437, XDPI: 75, YDPI: 75, Type: raster.SGray8, PageWidth: 298, PageHeight: 420}
		if err := w.StartPage(pg); err != nil {
			t.Fatal(err)
		}
		img := image.NewGray(image.Rect(0, 0, pg.Width, pg.Height))
		for i := range img.Pix {
			img.Pix[i] = 255
		}
		for x := range pg.Width {
			img.SetGray(x, 200, color.Gray{})
		}
		if err := w.WriteRows(img); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// On a CUPS server node's own port a job comes in two requests, as IPP
// Everywhere asks and as the clients of Mopria and Windows send:
// Create-Job tells the job's id, Send-Document brings its document (here
// PCLm, Mopria's raster), and the job is that id. Until then the job
// waits, listed and told; its sender alone sends its document or closes
// it; one closed, canceled or timed out without a document ends.
func TestCreateJob(t *testing.T) {
	e := flowtest.New(t, func(o *docflow.Options) { o.Converter = printconv.Converter{NoSystemFonts: true} })
	ctx := context.Background()
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office", To: []string{"keep"}})
	cups := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		docflow.PortHandler(e.M(), "office").ServeHTTP(w, r)
	}))
	defer cups.Close()
	uri := "ipp://" + strings.TrimPrefix(cups.URL, "http://") + "/printers/office"
	alice := &ipp.Client{URL: uri, User: "alice"}
	bob := &ipp.Client{URL: uri, User: "bob"}
	status := func(_ *ipp.Message, err error) ipp.Status {
		var se *ipp.StatusError
		if errors.As(err, &se) {
			return se.Status
		}
		if err != nil {
			t.Fatal(err)
		}
		return ipp.StatusOK
	}
	create := func(c *ipp.Client, name string) int32 {
		t.Helper()
		m := c.NewRequest(ipp.OpCreateJob)
		m.Add(ipp.TagOperation, ipp.Attr("job-name", ipp.Name(name)))
		m.Add(ipp.TagJob, ipp.Attr("job-account-id", ipp.Name("dept-7")))
		res, err := c.Do(ctx, m, nil)
		if err != nil {
			t.Fatal(err)
		}
		job := res.Group(ipp.TagJob)
		id, _ := job.Get("job-id").Int()
		if st, _ := job.Get("job-state").Int(); id <= 0 || ipp.JobState(st) != ipp.JobPending || job.Get("job-state-reasons").Text() != "job-incoming" {
			t.Fatalf("Create-Job: %v", job.Attrs)
		}
		return id
	}
	send := func(c *ipp.Client, id int32, last *bool, format string, doc []byte) (*ipp.Message, error) {
		m := c.NewRequest(ipp.OpSendDocument)
		m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))))
		if last != nil {
			m.Add(ipp.TagOperation, ipp.Attr("last-document", ipp.Boolean(*last)))
		}
		if format != "" {
			m.Add(ipp.TagOperation, ipp.Attr("document-format", ipp.MimeType(format)))
		}
		if doc != nil {
			return c.Do(ctx, m, bytes.NewReader(doc))
		}
		return c.Do(ctx, m, nil)
	}
	attrs := func(c *ipp.Client, op ipp.Op, id int32) (*ipp.Group, error) {
		m := c.NewRequest(op)
		m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))))
		res, err := c.Do(ctx, m, nil)
		if err != nil {
			return nil, err
		}
		return res.Group(ipp.TagJob), nil
	}
	state := func(g *ipp.Group) (ipp.JobState, string) {
		n, _ := g.Get("job-state").Int()
		return ipp.JobState(n), g.Get("job-state-reasons").Text()
	}
	yes, no := true, false

	id := create(alice, "statement")
	g, err := attrs(alice, ipp.OpGetJobAttributes, id)
	if st, why := state(g); err != nil || st != ipp.JobPending || why != "job-incoming" || g.Get("job-name").Text() != "statement" ||
		g.Get("job-originating-user-name").Text() != "alice" {
		t.Fatalf("waiting: %v %v", g, err)
	}
	// listed with the jobs not completed
	m := alice.NewRequest(ipp.OpGetJobs)
	m.Add(ipp.TagOperation, ipp.Attr("requested-attributes", ipp.Keywords([]string{"job-id", "job-state"})...))
	res, err := alice.Do(ctx, m, nil)
	if err != nil || len(res.Groups) != 2 {
		t.Fatalf("Get-Jobs: %v %v", res, err)
	}
	// its sender alone sends its document; last-document is required
	if st := status(send(bob, id, &yes, "application/PCLm", pclm(t, 2))); st != ipp.StatusClientErrorNotAuthorized {
		t.Errorf("another's Send-Document: %v", st)
	}
	if st := status(send(alice, id, nil, "application/PCLm", pclm(t, 2))); st != ipp.StatusClientErrorBadRequest {
		t.Errorf("Send-Document without last-document: %v", st)
	}
	res, err = send(alice, id, &yes, "application/PCLm", pclm(t, 2))
	if err != nil {
		t.Fatal(err)
	}
	if n, _ := res.Group(ipp.TagJob).Get("job-id").Int(); n != id {
		t.Fatalf("Send-Document answered job %d, not %d", n, id)
	}
	j := e.Job(id)
	if j.Dest != "office" || j.Name != "statement" || j.Source != "ipp" || j.Pages != 2 || j.Meta["ipp.job-account-id"] != "dept-7" ||
		j.Meta["ipp.document-format"] != "application/PCLm" {
		t.Fatalf("the job %+v", j)
	}
	if g, err := attrs(alice, ipp.OpGetJobAttributes, id); err != nil || g.Get("job-impressions-completed") == nil {
		t.Fatalf("the job's attributes: %v %v", g, err)
	}
	// a job has one document; closing it tells it
	if st := status(send(alice, id, &yes, "application/pdf", pclm(t, 1))); st != ipp.StatusClientErrorNotPossible {
		t.Errorf("a second document: %v", st)
	}
	if g, err := attrs(alice, ipp.OpCloseJob, id); err != nil || g.Get("job-id") == nil {
		t.Errorf("Close-Job of the job: %v %v", g, err)
	}

	// a document that is not the last is the job's still
	id = create(alice, "first of more")
	if _, err := send(alice, id, &no, "", pclm(t, 1)); err != nil {
		t.Fatal(err)
	}
	if j := e.Job(id); j.Pages != 1 {
		t.Errorf("not the last: %+v", j)
	}

	// closed without a document: aborted, and told so for a while
	id = create(alice, "empty")
	if _, err := attrs(bob, ipp.OpCloseJob, id); status(nil, err) != ipp.StatusClientErrorNotAuthorized {
		t.Errorf("another's Close-Job: %v", err)
	}
	g, err = attrs(alice, ipp.OpCloseJob, id)
	if st, why := state(g); err != nil || st != ipp.JobAborted || why != "job-data-insufficient" {
		t.Fatalf("closed without a document: %v %v", g, err)
	}
	if st := status(send(alice, id, &yes, "application/PCLm", pclm(t, 1))); st != ipp.StatusClientErrorNotPossible {
		t.Errorf("Send-Document of an ended job: %v", st)
	}
	// so is one whose last Send-Document has none
	id = create(alice, "nothing")
	if res, err := send(alice, id, &yes, "", nil); err != nil {
		t.Fatal(err)
	} else if st, why := state(res.Group(ipp.TagJob)); st != ipp.JobAborted || why != "job-data-insufficient" {
		t.Errorf("the last Send-Document empty: %v %s", st, why)
	}

	// canceled while it waits
	id = create(alice, "canceled")
	if _, err := attrs(alice, ipp.OpCancelJob, id); err != nil {
		t.Fatal(err)
	}
	if g, _ := attrs(alice, ipp.OpGetJobAttributes, id); g == nil {
		t.Fatal("the canceled job is not told")
	} else if st, _ := state(g); st != ipp.JobCanceled {
		t.Errorf("canceled: %v", st)
	}

	// Cancel-My-Jobs cancels the sender's, not another's
	mine, theirs := create(alice, "mine"), create(bob, "theirs")
	if _, err := alice.Do(ctx, alice.NewRequest(ipp.OpCancelMyJobs), nil); err != nil {
		t.Fatal(err)
	}
	if g, _ := attrs(alice, ipp.OpGetJobAttributes, mine); g == nil {
		t.Fatal("no job")
	} else if st, _ := state(g); st != ipp.JobCanceled {
		t.Errorf("alice's job after Cancel-My-Jobs: %v", st)
	}
	if g, _ := attrs(bob, ipp.OpGetJobAttributes, theirs); g == nil {
		t.Fatal("no job")
	} else if st, _ := state(g); st != ipp.JobPending {
		t.Errorf("bob's job after alice's Cancel-My-Jobs: %v", st)
	}

	// one whose document does not come in time is aborted
	dipp.SetCreatedTimeout(t, time.Millisecond)
	id = create(bob, "late")
	time.Sleep(5 * time.Millisecond)
	if g, _ := attrs(bob, ipp.OpGetJobAttributes, id); g == nil {
		t.Fatal("no job")
	} else if st, why := state(g); st != ipp.JobAborted || why != "aborted-by-system" {
		t.Errorf("timed out: %v %s", st, why)
	}

	// Identify-Printer is answered
	if _, err := alice.Do(ctx, alice.NewRequest(ipp.OpIdentifyPrinter), nil); err != nil {
		t.Errorf("Identify-Printer: %v", err)
	}
}

package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// TestMeta: what a job says about itself (job-account-id,
// job-accounting-user-id, attributes of its own, document-metadata) is
// kept with the job, told by Get-Job-Attributes and goes to the printer
// with it, by IPP and by the query of the API.
func TestMeta(t *testing.T) {
	e := newEnv(t, nil)
	p := e.printer("laser")
	c := e.ippClient("laser")
	ctx := context.Background()
	job := []ipp.Attribute{ipp.Attr("job-account-id", ipp.Name("dept-7")), ipp.Attr("x-cost-center", ipp.Keyword("cc-42")), ipp.Attr("copies", ipp.Integer(2))}
	meta := ipp.Attr("document-metadata", ipp.OctetString([]byte("title=Invoice 42")))
	js, err := c.PrintJob(ctx, "invoice", "application/x-bdf", job, bytes.NewReader(testDoc(t, 1)), meta)
	if err != nil {
		t.Fatal(err)
	}
	e.tick()
	for i := 0; len(p.Jobs()) == 0 && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	j := e.job(js.ID)
	if j.Meta == nil || j.Meta.AccountID != "dept-7" || strings.Join(j.Meta.Custom["x-cost-center"], ",") != "cc-42" || strings.Join(j.Meta.Metadata, ";") != "title=Invoice 42" {
		t.Fatalf("job %+v meta %+v", j, j.Meta)
	}
	if len(p.Jobs()) != 1 {
		t.Fatalf("the printer has %d jobs", len(p.Jobs()))
	}
	pa := p.Jobs()[0].Attrs
	if ipp.Get(pa, "job-account-id").Text() != "dept-7" || ipp.Get(pa, "x-cost-center").Text() != "cc-42" ||
		strings.Join(ipp.Get(pa, "document-metadata").Strings(), ";") != "title=Invoice 42" || ipp.Get(pa, "copies") != nil {
		t.Errorf("the printer got %v", pa)
	}
	st, err := c.GetJobAttributes(ctx, js.ID)
	if err != nil || ipp.Get(st.Attrs, "job-account-id").Text() != "dept-7" || ipp.Get(st.Attrs, "document-metadata") == nil {
		t.Errorf("Get-Job-Attributes: %+v %v", st, err)
	}
	// what cannot be taken is told unsupported, and the job printed without it
	m := c.NewRequest(ipp.OpPrintJob)
	m.Add(ipp.TagOperation, ipp.Attr("job-name", ipp.Name("bad")), ipp.Attr("document-format", ipp.MimeType("application/x-bdf")),
		ipp.Attr("document-metadata", ipp.OctetString([]byte("nokey"))))
	m.Add(ipp.TagJob, ipp.Attr("x-Bad", ipp.Name("1")))
	res, err := c.Do(ctx, m, bytes.NewReader(testDoc(t, 1)))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status() != ipp.StatusOKIgnoredOrSubstituted || res.Attr(ipp.TagUnsupportedGroup, "document-metadata") == nil || res.Attr(ipp.TagUnsupportedGroup, "x-Bad") == nil {
		t.Errorf("bad: %v %v", res.Status(), res.Groups)
	}
	if id, _ := res.Attr(ipp.TagJob, "job-id").Int(); e.job(id).Meta != nil {
		t.Errorf("kept %+v", e.job(id).Meta)
	}
	// by the API, in the query
	e2 := newEnv(t, func(o *Options) { o.SubmitTokens = map[string]string{"wms": "s3cret"} })
	p2 := e2.printer("laser")
	r, out := e2.do(nil, "POST", "/api/destinations/laser/jobs?name=api&job-account-id=dept-8&document-metadata=title%3DT&x-order=A-1&x-order=A-2", "wms:s3cret", testDoc(t, 1))
	if r.StatusCode != http.StatusCreated {
		t.Fatalf("submit: %d %s", r.StatusCode, out)
	}
	var j2 Job
	json.Unmarshal(out, &j2)
	if j2.Meta == nil || j2.Meta.AccountID != "dept-8" || strings.Join(j2.Meta.Custom["x-order"], ",") != "A-1,A-2" || strings.Join(j2.Meta.Metadata, ";") != "title=T" {
		t.Fatalf("job %s", out)
	}
	e2.tick()
	for i := 0; len(p2.Jobs()) == 0 && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if len(p2.Jobs()) != 1 || ipp.Get(p2.Jobs()[0].Attrs, "job-account-id").Text() != "dept-8" || strings.Join(ipp.Get(p2.Jobs()[0].Attrs, "x-order").Strings(), ",") != "A-1,A-2" {
		t.Errorf("the printer got %v", p2.Jobs())
	}
	if r, out := e2.do(nil, "POST", "/api/destinations/laser/jobs?name=bad&x-Bad=1", "wms:s3cret", testDoc(t, 1)); r.StatusCode != http.StatusBadRequest {
		t.Errorf("a bad name: %d %s", r.StatusCode, out)
	}
}

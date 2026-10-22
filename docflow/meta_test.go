package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// TestIPPMeta: the job attributes a printer is asked again
// (job-account-id, job-accounting-user-id, x-…, document-metadata) are
// kept as the job's metadata under ipp.…, told by Get-Job-Attributes and
// go to the printer with the job, whether they came as IPP attributes
// or as metadata of the API.
func TestIPPMeta(t *testing.T) {
	e := newEnv(t, nil)
	p := e.printer("laser")
	c := e.ippClient("laser")
	ctx := context.Background()
	job := []ipp.Attribute{ipp.Attr("job-account-id", ipp.Name("dept-7")), ipp.Attr("x-cost-center", ipp.Keyword("cc-42")),
		ipp.Attr("x-tags", ipp.Text("a"), ipp.Text("b")), ipp.Attr("copies", ipp.Integer(2))}
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
	if j.Meta["ipp.job-account-id"] != "dept-7" || j.Meta["ipp.x-cost-center"] != "cc-42" || fmt.Sprint(j.Meta["ipp.x-tags"]) != "[a b]" ||
		fmt.Sprint(j.Meta["ipp.document-metadata"]) != "[title=Invoice 42]" || fmt.Sprint(j.Meta["ipp.copies"]) != "2" {
		t.Fatalf("meta %v", j.Meta)
	}
	if len(p.Jobs()) != 1 {
		t.Fatalf("the printer has %d jobs", len(p.Jobs()))
	}
	pa := p.Jobs()[0].Attrs
	if ipp.Get(pa, "job-account-id").Text() != "dept-7" || ipp.Get(pa, "x-cost-center").Text() != "cc-42" || strings.Join(ipp.Get(pa, "x-tags").Strings(), ",") != "a,b" ||
		strings.Join(ipp.Get(pa, "document-metadata").Strings(), ";") != "title=Invoice 42" || ipp.Get(pa, "copies") != nil {
		t.Errorf("the printer got %v", pa)
	}
	st, err := c.GetJobAttributes(ctx, js.ID)
	if err != nil || ipp.Get(st.Attrs, "job-account-id").Text() != "dept-7" || ipp.Get(st.Attrs, "document-metadata") == nil {
		t.Errorf("Get-Job-Attributes: %+v %v", st, err)
	}
	// what cannot be taken is told unsupported, and not kept
	m := c.NewRequest(ipp.OpPrintJob)
	m.Add(ipp.TagOperation, ipp.Attr("job-name", ipp.Name("bad")), ipp.Attr("document-format", ipp.MimeType("application/x-bdf")),
		ipp.Attr("document-metadata", ipp.OctetString([]byte("nokey"))))
	m.Add(ipp.TagJob, ipp.Attr("x-Bad", ipp.Name("1")), ipp.Attr("media", ipp.Keyword("iso_a4_210x297mm")))
	res, err := c.Do(ctx, m, bytes.NewReader(testDoc(t, 1)))
	if err != nil {
		t.Fatal(err)
	}
	if res.Status() != ipp.StatusOKIgnoredOrSubstituted || res.Attr(ipp.TagUnsupportedGroup, "document-metadata") == nil || res.Attr(ipp.TagUnsupportedGroup, "x-Bad") == nil {
		t.Errorf("bad: %v %v", res.Status(), res.Groups)
	}
	id, _ := res.Attr(ipp.TagJob, "job-id").Int()
	if bad := e.job(id); bad.Meta["ipp.document-metadata"] != nil || bad.Meta["ipp.x-Bad"] != nil || bad.Meta["ipp.media"] != "iso_a4_210x297mm" {
		t.Errorf("kept %v", bad.Meta)
	}
	// as metadata of the API: the ipp. keys go to the printer, one that
	// is not an attribute's value stays a claim
	e2 := newEnv(t, func(o *Options) { o.SubmitTokens = map[string]string{"wms": "s3cret"} })
	p2 := e2.printer("laser")
	r, out := e2.do(nil, "POST", "/api/destinations/laser/jobs?name=api&meta.ipp.job-account-id=dept-8&meta.ipp.x-order=A-1&meta.ipp.x-Bad=1&meta.dept=sales", "wms:s3cret", testDoc(t, 1))
	if r.StatusCode != http.StatusCreated {
		t.Fatalf("submit: %d %s", r.StatusCode, out)
	}
	var j2 Job
	json.Unmarshal(out, &j2)
	if j2.Meta["ipp.job-account-id"] != "dept-8" || j2.Meta["dept"] != "sales" {
		t.Fatalf("job %s", out)
	}
	e2.tick()
	for i := 0; len(p2.Jobs()) == 0 && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if len(p2.Jobs()) != 1 {
		t.Fatalf("the printer has %d jobs", len(p2.Jobs()))
	}
	pa = p2.Jobs()[0].Attrs
	if ipp.Get(pa, "job-account-id").Text() != "dept-8" || ipp.Get(pa, "x-order").Text() != "A-1" || ipp.Get(pa, "x-Bad") != nil || ipp.Get(pa, "dept") != nil {
		t.Errorf("the printer got %v", pa)
	}
}

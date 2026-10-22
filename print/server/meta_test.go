package server

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/ipp"
)

// TestIPPMeta: what a job says about itself reaches the spool, the
// printer behind and Get-Job-Attributes: job-account-id,
// job-accounting-user-id, attributes of one's own (x-…) and
// document-metadata (with the document, of Print-Job or Send-Document).
// What cannot be taken is told unsupported.
func TestIPPMeta(t *testing.T) {
	f := newIPPFixture(t, ippConfigYAML)
	c := f.client("office", "alice")
	ctx := context.Background()
	job := []ipp.Attribute{
		ipp.Attr("job-account-id", ipp.Name("dept-7")),
		ipp.Attr("job-accounting-user-id", ipp.Name("alice")),
		ipp.Attr("x-cost-center", ipp.Keyword("cc-42")),
		ipp.Attr("x-tags", ipp.Text("a"), ipp.Text("b")),
	}
	meta := ipp.Attr("document-metadata", ipp.OctetString([]byte("title=Invoice 42")), ipp.OctetString([]byte("x-order=A-1")))
	js, err := c.PrintJob(ctx, "invoice.bdf", "application/x-bdf", job, bytes.NewReader(docBytes(t)), meta)
	if err != nil {
		t.Fatal(err)
	}
	j := f.waitJob(js.ID, ipp.JobCompleted)
	m := j.Ticket.Meta
	if m.AccountID != "dept-7" || m.AccountingUserID != "alice" || strings.Join(m.Metadata, ";") != "title=Invoice 42;x-order=A-1" ||
		strings.Join(m.Custom["x-cost-center"], ",") != "cc-42" || strings.Join(m.Custom["x-tags"], ",") != "a,b" {
		t.Fatalf("meta %+v", m)
	}
	// the printer behind was asked them
	got := f.printer.Jobs()
	if len(got) != 1 {
		t.Fatalf("the printer has %d jobs", len(got))
	}
	pa := got[0].Attrs
	if ipp.Get(pa, "job-account-id").Text() != "dept-7" || ipp.Get(pa, "x-cost-center").Text() != "cc-42" ||
		strings.Join(ipp.Get(pa, "x-tags").Strings(), ",") != "a,b" ||
		strings.Join(ipp.Get(pa, "document-metadata").Strings(), ";") != "title=Invoice 42;x-order=A-1" {
		t.Errorf("the printer got %v", pa)
	}
	// Get-Job-Attributes tells them: the account and the attributes of
	// one's own among the job template attributes, document-metadata
	// among the job description
	get := func(id int32, req ...string) *ipp.Message {
		t.Helper()
		m := c.NewRequest(ipp.OpGetJobAttributes)
		m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))), ipp.Attr("requested-attributes", ipp.Keywords(req)...))
		res, err := c.Do(ctx, m, nil)
		if err != nil {
			t.Fatal(err)
		}
		return res
	}
	res := get(js.ID, "job-template")
	if res.Attr(ipp.TagJob, "job-account-id").Text() != "dept-7" || strings.Join(res.Attr(ipp.TagJob, "x-tags").Strings(), ",") != "a,b" ||
		res.Attr(ipp.TagJob, "document-metadata") != nil {
		t.Errorf("job-template: %v", res.Groups)
	}
	res = get(js.ID, "job-description")
	if a := res.Attr(ipp.TagJob, "document-metadata"); a == nil || a.Values[0].Tag != ipp.TagOctetString || strings.Join(a.Strings(), ";") != "title=Invoice 42;x-order=A-1" ||
		res.Attr(ipp.TagJob, "x-tags") != nil {
		t.Errorf("job-description: %v", res.Groups)
	}
	if res = get(js.ID, "all"); res.Attr(ipp.TagJob, "document-metadata") == nil || res.Attr(ipp.TagJob, "x-cost-center") == nil {
		t.Errorf("all: %v", res.Groups)
	}
	// Create-Job, then document-metadata with Send-Document
	js2, err := c.CreateJob(ctx, "two.bdf", job[:1])
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.SendDocument(ctx, js2.ID, "two.bdf", "application/x-bdf", bytes.NewReader(docBytes(t)), true, meta); err != nil {
		t.Fatal(err)
	}
	if m := f.waitJob(js2.ID, ipp.JobCompleted).Ticket.Meta; m.AccountID != "dept-7" || len(m.Metadata) != 2 || len(m.Custom) != 0 {
		t.Errorf("by Send-Document: %+v", m)
	}
	// what cannot be taken is unsupported: the job is printed without it
	bad := map[string]ipp.Attribute{
		"account":  ipp.Attr("job-account-id", ipp.Integer(7)),
		"name":     ipp.Attr("x-Bad", ipp.Name("1")),
		"metadata": ipp.Attr("document-metadata", ipp.OctetString([]byte("nokey"))),
	}
	for name, a := range bad {
		m := c.NewRequest(ipp.OpPrintJob)
		m.Add(ipp.TagOperation, ipp.Attr("job-name", ipp.Name(name)), ipp.Attr("document-format", ipp.MimeType("application/x-bdf")))
		if a.Name == "document-metadata" {
			m.Add(ipp.TagOperation, a)
		} else {
			m.Add(ipp.TagJob, a)
		}
		res, err := c.Do(ctx, m, bytes.NewReader(docBytes(t)))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if res.Status() != ipp.StatusOKIgnoredOrSubstituted || res.Attr(ipp.TagUnsupportedGroup, a.Name) == nil {
			t.Errorf("%s: %v %v", name, res.Status(), res.Groups)
		}
		id, _ := res.Attr(ipp.TagJob, "job-id").Int()
		if j := f.waitJob(id, ipp.JobCompleted); !j.Ticket.Meta.Empty() {
			t.Errorf("%s: kept %+v", name, j.Ticket.Meta)
		}
		// with fidelity the job is refused
		m = c.NewRequest(ipp.OpValidateJob)
		m.Add(ipp.TagOperation, ipp.Attr("ipp-attribute-fidelity", ipp.Boolean(true)))
		if a.Name == "document-metadata" {
			m.Add(ipp.TagOperation, a)
		} else {
			m.Add(ipp.TagJob, a)
		}
		if _, err := c.Do(ctx, m, nil); statusOf(err) != ipp.StatusClientErrorAttributesOrValuesNotSupported {
			t.Errorf("%s with fidelity: %v", name, err)
		}
	}
	// the printer says it takes the account
	p, err := c.GetPrinterAttributes(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"job-account-id-supported", "job-accounting-user-id-supported"} {
		if a, ok := p.Has(name); !ok || !slices.Equal(a.Ints(), []int32{1}) {
			t.Errorf("%s: %v", name, a)
		}
	}
	if a, _ := p.Has("job-creation-attributes-supported"); !a.Has("job-account-id") || !a.Has("document-metadata") {
		t.Errorf("job-creation-attributes-supported: %v", a)
	}
}

package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/ipp"
)

func ptr(x float64) *float64 { return &x }

// An assert right after an entry is the sender's gate: the request is
// answered once its checks pass, and refused (422) with what failed when
// they do not; the job is kept, aborted, and goes to the assert's
// on_error nodes.
func TestAssert(t *testing.T) {
	e := newEnv(t, nil)
	e.node("seen", Node{Kind: KindSink}, http.StatusOK)
	e.node("errors", Node{Kind: KindSink}, http.StatusOK)
	for _, bad := range []Check{
		{Key: ""},
		{Key: "a b"},
		{Key: "doc", Pattern: "("},
		{Key: "amount", Min: ptr(10), Max: ptr(1)},
	} {
		e.node("check", Node{Kind: KindAssert, Checks: []Check{bad}}, http.StatusBadRequest)
	}
	// put before the nodes after it; none checked passes everything
	e.node("check", Node{Kind: KindAssert}, http.StatusOK)
	e.node("check", Node{Kind: KindAssert, To: []string{"seen"}, OnError: []string{"errors"}, Checks: []Check{
		{Key: "doc-type", Enum: []string{"invoice", "order", " invoice "}},
		{Key: "amount", Optional: true, Min: ptr(0), Max: ptr(10000)},
		{Key: "zip", Optional: true, Pattern: `^\d{3}-\d{4}$`, Message: "zip is a postal code: 123-4567"},
		{Key: "$pages", Max: ptr(3)},
	}}, http.StatusOK)
	var n Node
	json.Unmarshal(e.admin("GET", "/api/nodes/check", nil, http.StatusOK), &n)
	if len(n.Checks) != 4 || strings.Join(n.Checks[0].Enum, ",") != "invoice,order" {
		t.Fatalf("checks %+v", n.Checks)
	}
	var made struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("in", Node{Kind: KindUpload, To: []string{"check"}}, http.StatusOK), &made)

	post := func(query string, pages, status int) (*Job, apiRefusal) {
		t.Helper()
		res, out := e.do(nil, "POST", "/api/upload/in?name=x&"+query, made.Key, testDoc(t, pages))
		if res.StatusCode != status {
			t.Fatalf("%s: %d %s; want %d", query, res.StatusCode, out, status)
		}
		var ref apiRefusal
		if status == http.StatusUnprocessableEntity {
			json.Unmarshal(out, &ref)
			return e.job(ref.Job), ref
		}
		var j Job
		json.Unmarshal(out, &j)
		return &j, ref
	}
	j, _ := post("meta.doc-type=invoice&meta.amount=12.5", 1, http.StatusCreated)
	if legs(j) != "in:routed check:routed seen:kept" || j.State != JobCompleted {
		t.Fatalf("passed %s: %s", j.State, legs(j))
	}
	for _, c := range []struct {
		query string
		pages int
		want  []string
	}{
		{"", 1, []string{"doc-type: required"}},
		{"meta.doc-type=", 1, []string{"doc-type: required"}},
		{"meta.doc-type=memo", 1, []string{`doc-type: "memo" is not one of invoice, order`}},
		{"meta.doc-type=order&meta.amount=12000", 1, []string{"amount: 12000 is more than 10000"}},
		{"meta.doc-type=order&meta.amount=-1", 1, []string{"amount: -1 is less than 0"}},
		{"meta.doc-type=order&meta.amount=lots", 1, []string{`amount: "lots" is not a number`}},
		{"meta.doc-type=order&meta.zip=1234567", 1, []string{"zip is a postal code: 123-4567"}},
		{"meta.amount=NaN", 4, []string{"doc-type: required", `amount: "NaN" is not a number`, "$pages: 4 is more than 3"}},
	} {
		j, ref := post(c.query, c.pages, http.StatusUnprocessableEntity)
		var got []string
		for _, f := range ref.Failures {
			if f.Node != "check" {
				t.Errorf("%s: failed at %s", c.query, f.Node)
			}
			got = append(got, f.Error)
		}
		if strings.Join(got, "|") != strings.Join(c.want, "|") || !strings.HasPrefix(ref.Error, "refused: check: ") {
			t.Errorf("%q: %q %v; want %v", c.query, ref.Error, got, c.want)
		}
		if j.State != JobAborted || legs(j) != "in:routed check:refused errors:kept" || j.Legs[1].Error != strings.Join(c.want, "; ") {
			t.Errorf("%q: %s %s", c.query, j.State, legs(j))
		}
	}
	// what reached the refusing assert went to its on_error node, with
	// the failure; nothing reached the node after it
	if kept := e.m().sinkJobs(context.Background(), "seen"); len(kept) != 1 {
		t.Errorf("seen keeps %d", len(kept))
	}

	// asserts side by side are all checked; a queue sent to over the
	// API gates its job as an entry does
	e.node("site", Node{Kind: KindAssert, To: []string{"seen"}, Checks: []Check{{Key: "site", Enum: []string{"tokyo"}}}}, http.StatusOK)
	e.node("fan", Node{Kind: KindQueue, To: []string{"check", "site", "seen"}}, http.StatusOK)
	res, out := e.do(nil, "POST", "/api/destinations/fan/jobs?name=q&meta.site=osaka", "", testDoc(t, 1))
	var ref apiRefusal
	json.Unmarshal(out, &ref)
	if res.StatusCode != http.StatusUnprocessableEntity || len(ref.Failures) != 2 || ref.Failures[0].Node != "check" || ref.Failures[1].Node != "site" {
		t.Fatalf("fan: %d %s", res.StatusCode, out)
	}
	// site has no on_error and the manager no default: it stops there
	if j := e.job(ref.Job); legs(j) != "fan:routed check:refused errors:kept site:refused" {
		t.Errorf("fan: %s", legs(j))
	}
	if res, out := e.do(nil, "POST", "/api/destinations/fan/jobs?name=q&meta.site=tokyo&meta.doc-type=order", "", testDoc(t, 1)); res.StatusCode != http.StatusCreated {
		t.Fatalf("fan passes: %d %s", res.StatusCode, out)
	}
	// sent to the assert itself
	if res, _ := e.do(nil, "POST", "/api/destinations/site/jobs?name=q", "", testDoc(t, 1)); res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("to the assert: %d", res.StatusCode)
	}

	// an assert further on checks the job as it reaches it: the sender's
	// request was taken, and the job is aborted there
	e.node("late", Node{Kind: KindAssert, To: []string{"seen"}, OnError: []string{"errors"}, Checks: []Check{{Key: "doc-type"}}}, http.StatusOK)
	e.node("hop", Node{Kind: KindQueue, To: []string{"late"}}, http.StatusOK)
	e.node("in2", Node{Kind: KindUpload, To: []string{"hop"}}, http.StatusOK)
	j = e.printJob("in2", "late", 1)
	if legs(j) != "in2:routed hop:routed late:refused errors:kept" || j.State != JobAborted {
		t.Errorf("late: %s %s", j.State, legs(j))
	}

	// the live view counts what an assert refused
	var live Live
	json.Unmarshal(e.admin("GET", "/api/topology/live", nil, http.StatusOK), &live)
	if live.Nodes["check"].Refused != 9 || live.Nodes["site"].Refused != 2 {
		t.Errorf("live %+v", live.Nodes)
	}
}

// With no on_error of its own, a refusing assert sends the job to the
// manager's default.
func TestAssertDefault(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.OnError = []string{"ops"} })
	e.node("ops", Node{Kind: KindSink}, http.StatusOK)
	e.node("seen", Node{Kind: KindSink}, http.StatusOK)
	e.node("check", Node{Kind: KindAssert, To: []string{"seen"}, Checks: []Check{{Key: "doc-type"}}}, http.StatusOK)
	e.node("in", Node{Kind: KindUpload, To: []string{"check"}}, http.StatusOK)
	res, out := e.do(nil, "POST", "/api/destinations/in/jobs?name=x", "", testDoc(t, 1))
	var ref apiRefusal
	json.Unmarshal(out, &ref)
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("%d %s", res.StatusCode, out)
	}
	if j := e.job(ref.Job); legs(j) != "in:routed check:refused ops:kept" || j.Legs[2].From != "check" {
		t.Errorf("legs %s", legs(j))
	}
}

// A CUPS server answers a refusal as IPP: client-error-bad-request, with
// what failed.
func TestAssertIPP(t *testing.T) {
	e := newEnv(t, nil)
	e.node("seen", Node{Kind: KindSink}, http.StatusOK)
	e.node("check", Node{Kind: KindAssert, To: []string{"seen"}, Checks: []Check{{Key: "ipp.job-account-id", Pattern: "^dept-"}}}, http.StatusOK)
	c := e.cups("office", "check")
	ctx := context.Background()
	print := func(account string) *ipp.Message {
		t.Helper()
		m := c.NewRequest(ipp.OpPrintJob)
		m.Add(ipp.TagOperation, ipp.Attr("job-name", ipp.Name("slip")), ipp.Attr("document-format", ipp.MimeType("application/x-bdf")))
		m.Add(ipp.TagJob, ipp.Attr("job-account-id", ipp.Name(account)))
		res, err := c.Do(ctx, m, bytes.NewReader(testDoc(t, 1)))
		if err != nil && !ipp.IsStatus(err, ipp.StatusClientErrorBadRequest) {
			t.Fatal(err)
		}
		return res
	}
	if res := print("dept-7"); res.Status() != ipp.StatusOK {
		t.Fatalf("passes: %v", res.Status())
	}
	res := print("sales")
	if msg := res.Attr(ipp.TagOperation, "status-message").Text(); res.Status() != ipp.StatusClientErrorBadRequest || !strings.Contains(msg, `ipp.job-account-id: "sales" does not match ^dept-`) {
		t.Fatalf("refused: %v %q", res.Status(), msg)
	}
	if res.Attr(ipp.TagJob, "job-id") != nil {
		t.Error("a refused job has an IPP job")
	}
	if j := e.job(2); legs(j) != "office:routed check:refused" || j.State != JobAborted {
		t.Errorf("job %s %s", j.State, legs(j))
	}
}

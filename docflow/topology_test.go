package docflow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/yuruidb"
)

// node puts a routing node.
func (e *env) node(name string, n Node, status int) []byte {
	e.t.Helper()
	return e.admin("PUT", "/api/nodes/"+name, n, status)
}

// tick runs a tick through the API.
func (e *env) tick() TickReport {
	e.t.Helper()
	var rep TickReport
	json.Unmarshal(e.admin("POST", "/api/tick", nil, http.StatusOK), &rep)
	return rep
}

// job reads a job through the API.
func (e *env) job(id any) *Job {
	e.t.Helper()
	var j Job
	json.Unmarshal(e.admin("GET", fmt.Sprint("/api/jobs/", id), nil, http.StatusOK), &j)
	return &j
}

// legs are the states of a job's legs, "node:state", in order.
func legs(j *Job) string {
	var out []string
	for _, l := range j.Legs {
		s := l.Node + ":" + l.State
		if l.Member != "" {
			s += "(" + l.Member + ")"
		}
		out = append(out, s)
	}
	return strings.Join(out, " ")
}

// A queue fans a job out to two locations: one job, a slip on each
// board, one document, canceled together, removed together.
func TestQueueFanout(t *testing.T) {
	e := newEnv(t, nil)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Title: "A"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/wh-b", Node{Kind: KindDisplay, Title: "B"}, http.StatusOK)
	e.node("both", Node{Kind: KindQueue, Title: "両方の倉庫", Description: "  Every slip goes to both warehouses.\nThe night shift reads wh-b. " + strings.Repeat("x", 1000), To: []string{"wh-a", "wh-b"}}, http.StatusOK)
	if n, err := e.m().Node(context.Background(), "both"); err != nil || !strings.HasPrefix(n.Description, "Every slip goes to both warehouses.\nThe night") || len(n.Description) != 1000 {
		t.Fatalf("description %+v %v", n, err)
	}
	a, b := e.enroll("wh-a", "a"), e.enroll("wh-b", "b")

	js := e.printJob("both", "picking", 2)
	if js.Number != 1 || js.State != JobCompleted {
		t.Fatalf("job %+v", js)
	}
	sa, sb := a.board().Open[0], b.board().Open[0]
	if sa.Job != 1 || sb.Job != 1 || sa.ID == sb.ID || sa.Name != "picking" || sb.Pages != 2 {
		t.Fatalf("slips %+v %+v", sa, sb)
	}
	j := e.job(1)
	if j.Dest != "both" || j.State != JobCompleted || j.Live != 2 || legs(j) != "both:routed wh-a:shown wh-b:shown" {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	if res, out := e.do(nil, "GET", "/api/jobs/1/document", adminToken, nil); res.StatusCode != http.StatusOK || !bytes.HasPrefix(out, bdf.Magic[:]) {
		t.Fatalf("document: %d", res.StatusCode)
	}
	// each board works its slip on its own
	a.action(sa.ID, "done", http.StatusOK)
	if got := e.job(1); got.State != JobCompleted {
		t.Fatalf("after one done: %s", got.State)
	}
	// it is canceled: the slip still open leaves its board
	if res, out := e.do(nil, "POST", "/api/jobs/1/cancel", adminToken, nil); res.StatusCode != http.StatusOK {
		t.Fatalf("cancel: %d %s", res.StatusCode, out)
	}
	if bb := b.board(); bb == nil || len(bb.Open) != 0 || bb.Ended[0].State != SlipCanceled {
		t.Fatalf("board b after cancel %+v", bb)
	}
	if ab := a.board(); ab == nil || ab.Ended[0].State != SlipDone {
		t.Fatalf("board a after cancel %+v", ab)
	}
	if got := e.job(1); got.State != JobCanceled || got.Canceled == nil {
		t.Fatalf("job after cancel %+v", got)
	}
	if res, _ := e.do(nil, "POST", "/api/jobs/1/cancel", adminToken, nil); res.StatusCode != http.StatusConflict {
		t.Fatalf("cancel twice: %d", res.StatusCode)
	}
	// the job is not the locations' (their jobs are none), it is the queue's
	var jobs []Job
	json.Unmarshal(e.admin("GET", "/api/destinations/both/jobs", nil, http.StatusOK), &jobs)
	if len(jobs) != 1 || jobs[0].Number != 1 {
		t.Fatalf("jobs of both %+v", jobs)
	}
	if json.Unmarshal(e.admin("GET", "/api/destinations/wh-a/jobs", nil, http.StatusOK), &jobs); len(jobs) != 0 {
		t.Fatalf("jobs of wh-a %+v", jobs)
	}
	// after the keep, the slips go, then the job and its document
	e.clock.add(14*24*time.Hour + time.Minute)
	rep := e.tick()
	if rep.Removed != 2 || rep.Jobs != 1 {
		t.Fatalf("tick %+v", rep)
	}
	if res, _ := e.do(nil, "GET", "/api/jobs/1", adminToken, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("the job is still there: %d", res.StatusCode)
	}
	if _, err := e.st.docs.Open(context.Background(), jobID(1)); err == nil {
		t.Fatal("the document is still there")
	}
}

// A class sends jobs to its members in turn, and skips a member that
// does not take one.
func TestClass(t *testing.T) {
	e := newEnv(t, nil)
	p1, p2 := e.printer("p1"), e.printer("p2")
	e.node("tokyo", Node{Kind: KindClass, Members: []string{"p1", "p2"}}, http.StatusOK)
	for i := range 4 {
		js := e.printJob("tokyo", fmt.Sprint("job ", i), 1)
		if js.State != JobProcessing {
			t.Fatalf("job %+v", js)
		}
	}
	var pending []Delivery
	json.Unmarshal(e.admin("GET", "/api/deliveries", nil, http.StatusOK), &pending)
	if len(pending) != 4 || pending[0].Printer != "p1" || pending[1].Printer != "p2" || pending[2].Printer != "p1" {
		t.Fatalf("pending %+v", pending)
	}
	rep := e.tick()
	if rep.Delivered != 4 || len(p1.Jobs()) != 2 || len(p2.Jobs()) != 2 {
		t.Fatalf("tick %+v: %d %d", rep, len(p1.Jobs()), len(p2.Jobs()))
	}
	j := e.job(3)
	if j.State != JobCompleted || legs(j) != "tokyo:done(p1)" || j.Live != 0 {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	var jobs []Job
	json.Unmarshal(e.admin("GET", "/api/destinations/tokyo/jobs", nil, http.StatusOK), &jobs)
	if len(jobs) != 4 || jobs[1].Number != 3 || jobs[1].State != JobCompleted {
		t.Fatalf("jobs of tokyo %+v", jobs)
	}
	if rep := e.tick(); rep.Delivered != 0 {
		t.Fatalf("delivered again %+v", rep)
	}

	// failover: the first member is not there; the next takes it
	e.admin("PUT", "/api/nodes/gone", Node{Kind: KindPrinter, URL: "ipp://127.0.0.1:1/ipp/print"}, http.StatusOK)
	e.node("osaka", Node{Kind: KindClass, Members: []string{"gone", "p2"}, Policy: PolicyFailover}, http.StatusOK)
	e.printJob("osaka", "failover", 1)
	rep = e.tick()
	if rep.Undelivered != 1 || rep.Delivered != 0 {
		t.Fatalf("tick %+v", rep)
	}
	d, _ := e.m().Delivery(context.Background(), e.job(5).Legs[0].Ref)
	if d.State != DeliveryPending || d.Attempts != 1 || !slices.Equal(d.Tried, []string{"gone"}) || !d.RetryAt.After(e.clock.Now()) {
		t.Fatalf("delivery %+v", d)
	}
	if rep := e.tick(); rep.Undelivered != 0 || rep.Delivered != 0 {
		t.Fatalf("tried before its time %+v", rep)
	}
	e.clock.add(time.Minute)
	if rep := e.tick(); rep.Delivered != 1 || len(p2.Jobs()) != 3 {
		t.Fatalf("not failed over %+v", rep)
	}
	if j := e.job(5); j.State != JobCompleted || legs(j) != "osaka:done(p2)" {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	// an administrator retries at once what waits
	e.printJob("osaka", "retry", 1)
	e.tick()
	j = e.job(6)
	e.admin("POST", "/api/deliveries/"+j.Legs[0].Ref+"/retry", nil, http.StatusOK)
	if rep := e.tick(); rep.Delivered != 1 {
		t.Fatalf("retry %+v", rep)
	}
}

// A job printed to a printer itself goes to that printer; in the
// background of a process it is sent without a tick.
func TestPrinterDirect(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.Background = true })
	p := e.printer("laser")
	js := e.printJob("laser", "direct", 1)
	if js.State != JobProcessing {
		t.Fatalf("job %+v", js)
	}
	for i := 0; len(p.Jobs()) == 0 && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if len(p.Jobs()) != 1 || p.Jobs()[0].Name != "direct" {
		t.Fatalf("not printed in the background: %d", len(p.Jobs()))
	}
	for i := 0; e.job(1).State != JobCompleted && i < 100; i++ {
		time.Sleep(50 * time.Millisecond)
	}
	if j := e.job(1); j.State != JobCompleted || legs(j) != "laser:done" {
		t.Fatalf("job %+v", j)
	}
}

// An inbox holds jobs until a person releases them to a node of its
// choice, or discards them.
func TestInbox(t *testing.T) {
	e := newEnv(t, nil)
	p := e.printer("laser")
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
	e.node("review", Node{Kind: KindInbox, Title: "確認", To: []string{"wh-a", "laser"}, Keep: Duration(time.Hour)}, http.StatusOK)
	e.node("in", Node{Kind: KindQueue, To: []string{"review"}}, http.StatusOK)
	d := e.enroll("wh-a", "tablet")
	d.board()

	js := e.printJob("in", "to review", 1)
	if js.State != JobHeld {
		t.Fatalf("job %+v", js)
	}
	var v InboxView
	json.Unmarshal(e.admin("GET", "/api/inboxes/review", nil, http.StatusOK), &v)
	if len(v.Held) != 1 || v.Held[0].Job != 1 || v.Held[0].Name != "to review" || v.Held[0].Dest != "in" || v.Title != "確認" {
		t.Fatalf("inbox %+v", v)
	}
	if d.board() != nil {
		t.Fatal("on the board before its release")
	}
	item := v.Held[0].ID
	// not to a node the inbox does not name
	e.admin("POST", "/api/items/"+item+"/release", map[string]any{"to": []string{"review"}}, http.StatusBadRequest)
	// to the location, by Sato
	var it Item
	json.Unmarshal(e.admin("POST", "/api/items/"+item+"/release", map[string]any{"to": []string{"wh-a"}, "worker": "Sato"}, http.StatusOK), &it)
	if it.State != ItemReleased || it.Released.Worker != "Sato" || !slices.Equal(it.Released.To, []string{"wh-a"}) {
		t.Fatalf("item %+v", it)
	}
	if b := d.board(); b == nil || len(b.Open) != 1 || b.Open[0].Job != 1 {
		t.Fatalf("board %+v", b)
	}
	j := e.job(1)
	if j.State != JobCompleted || legs(j) != "in:routed review:released wh-a:shown" || j.Live != 2 {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	e.admin("POST", "/api/items/"+item+"/release", nil, http.StatusConflict)
	json.Unmarshal(e.admin("GET", "/api/inboxes/review", nil, http.StatusOK), &v)
	if len(v.Held) != 0 || len(v.Ended) != 1 {
		t.Fatalf("inbox after release %+v", v)
	}

	// released to everything the inbox names, when no node is named
	e.printJob("in", "both", 1)
	json.Unmarshal(e.admin("GET", "/api/inboxes/review", nil, http.StatusOK), &v)
	e.admin("POST", "/api/items/"+v.Held[0].ID+"/release", nil, http.StatusOK)
	if rep := e.tick(); rep.Delivered != 1 || len(p.Jobs()) != 1 {
		t.Fatalf("tick %+v", rep)
	}
	if j := e.job(2); j.State != JobCompleted || legs(j) != "in:routed review:released wh-a:shown laser:done" {
		t.Fatalf("job %+v: %s", j, legs(j))
	}

	// discarded: the job ends here
	e.printJob("in", "discard", 1)
	json.Unmarshal(e.admin("GET", "/api/inboxes/review", nil, http.StatusOK), &v)
	json.Unmarshal(e.admin("POST", "/api/items/"+v.Held[0].ID+"/discard", map[string]any{"worker": "Sato"}, http.StatusOK), &it)
	if it.State != ItemDiscarded || it.Discarded.Worker != "Sato" {
		t.Fatalf("item %+v", it)
	}
	if j := e.job(3); j.State != JobCompleted || legs(j) != "in:routed review:discarded" {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	// canceled while held: it leaves the inbox
	e.printJob("in", "cancel", 1)
	if res, out := e.do(nil, "POST", "/api/jobs/4/cancel", adminToken, nil); res.StatusCode != http.StatusOK {
		t.Fatalf("cancel: %d %s", res.StatusCode, out)
	}
	json.Unmarshal(e.admin("GET", "/api/inboxes/review", nil, http.StatusOK), &v)
	if len(v.Held) != 0 || len(v.Ended) != 4 || v.Ended[0].State != ItemDiscarded {
		t.Fatalf("inbox after cancel %+v", v)
	}
	if j := e.job(4); j.State != JobCanceled {
		t.Fatalf("job %+v", j)
	}
	// an inbox with items held is not deleted; the ended ones go after its keep
	e.printJob("in", "held", 1)
	e.admin("DELETE", "/api/nodes/review", nil, http.StatusConflict)
	e.clock.add(2 * time.Hour)
	if rep := e.tick(); rep.Items != 4 {
		t.Fatalf("tick %+v", rep)
	}
	json.Unmarshal(e.admin("GET", "/api/inboxes/review", nil, http.StatusOK), &v)
	if len(v.Held) != 1 || len(v.Ended) != 0 {
		t.Fatalf("inbox after keep %+v", v)
	}
	if j := e.job(3); j.Live != 0 {
		t.Fatalf("job 3 still needed: %+v", j)
	}
}

// An archive keeps a job's document and a record of it, written once.
func TestArchive(t *testing.T) {
	e := newEnv(t, nil)
	store := "bbolt:" + filepath.ToSlash(filepath.Join(t.TempDir(), "audit.db"))
	e.node("audit", Node{Kind: KindArchive, Store: "nowhere://x"}, http.StatusBadRequest)
	e.node("audit", Node{Kind: KindArchive, Store: store, Collection: "no/slash"}, http.StatusBadRequest)
	e.node("audit", Node{Kind: KindArchive, Title: "監査", Store: store}, http.StatusOK)
	e.node("audit-in", Node{Kind: KindQueue, To: []string{"audit"}}, http.StatusOK)
	js := e.printJob("audit-in", "invoice", 3)
	if js.State != JobProcessing {
		t.Fatalf("job %+v", js)
	}
	if rep := e.tick(); rep.Delivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	j := e.job(1)
	if j.State != JobCompleted || legs(j) != "audit-in:routed audit:done" {
		t.Fatalf("job %+v: %s", j, legs(j))
	}
	ctx := context.Background()
	st, err := e.st.archive(ctx, store, "archive-audit")
	if err != nil {
		t.Fatal(err)
	}
	page, err := st.List(ctx, yuruidb.ListOptions{Limit: 10})
	if err != nil || len(page.Blobs) != 2 {
		t.Fatalf("archived %+v %v", page, err)
	}
	key := archiveKey(j, j.received())
	if page.Blobs[0].ID != key+".bdf" || page.Blobs[1].ID != key+".json" || !strings.HasPrefix(key, "2026-10-20.0000000001.") {
		t.Fatalf("keys %+v", page.Blobs)
	}
	r, err := st.Open(ctx, key+".json")
	if err != nil {
		t.Fatal(err)
	}
	var rec ArchiveRecord
	if err := json.NewDecoder(r).Decode(&rec); err != nil {
		t.Fatal(err)
	}
	r.Close()
	if rec.Job != 1 || rec.Name != "invoice" || rec.Pages != 3 || rec.SHA256 != j.SHA256 || rec.Archive != "audit" || rec.Dest != "audit-in" ||
		rec.Kind != PayloadBDF || rec.Format != "bdf" || rec.Payload != j.ID || rec.Object != key+".bdf" || rec.Meta["$entry"] != "audit-in" || len(rec.Payloads) != 1 {
		t.Fatalf("record %+v", rec)
	}
	r, err = st.Open(ctx, key+".bdf")
	if err != nil {
		t.Fatal(err)
	}
	head := make([]byte, len(bdf.Magic))
	r.Read(head)
	r.Close()
	if !bytes.Equal(head, bdf.Magic[:]) || page.Blobs[0].Size != j.Size {
		t.Fatalf("document %q %d", head, page.Blobs[0].Size)
	}
	// written once: a second delivery of the same job would not change it
	if _, err := e.m().archive(ctx, &Node{Name: "audit", Kind: KindArchive, Store: store}, j, j.received(), nil); err != nil {
		t.Fatal(err)
	}
	if page, _ := st.List(ctx, yuruidb.ListOptions{Limit: 10}); len(page.Blobs) != 2 {
		t.Fatalf("written twice %+v", page)
	}
}

// Nodes are checked when put: their kind, their edges, their names in
// the one namespace, and cycles; a node sent to is not deleted.
func TestNodes(t *testing.T) {
	e := newEnv(t, nil)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
	e.printer("laser")
	e.admin("PUT", "/api/nodes/wh-b", Node{Kind: KindDisplay, Fallback: &Fallback{Printer: "laser"}}, http.StatusOK)
	e.node("x", Node{Kind: "tee", To: []string{"wh-a"}}, http.StatusBadRequest)
	// a node is put before the nodes after it: none yet is fine
	e.node("x", Node{Kind: KindQueue}, http.StatusOK)
	e.node("x", Node{Kind: KindClass}, http.StatusOK)
	e.admin("DELETE", "/api/nodes/x", nil, http.StatusNoContent)
	e.node("x", Node{Kind: KindQueue, To: []string{"nowhere"}}, http.StatusBadRequest)
	e.node("x", Node{Kind: KindClass, Members: []string{"laser"}, Policy: "random"}, http.StatusBadRequest)
	e.node("wh-a", Node{Kind: KindQueue, To: []string{"laser"}}, http.StatusConflict)
	e.node("x", Node{Kind: KindQueue, To: []string{"x"}}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/wh-c", Node{Kind: KindDisplay, Fallback: &Fallback{Printer: "wh-a"}}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/p", Node{Kind: KindPrinter}, http.StatusBadRequest)
	e.node("x", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Accept: "maybe"}, http.StatusBadRequest)
	e.node("x", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Formats: []string{"PDF/A"}}, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/laser", Node{Kind: KindDisplay}, http.StatusConflict)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindPrinter, URL: "ipp://127.0.0.1:1/ipp/print"}, http.StatusConflict)

	e.node("a", Node{Kind: KindQueue, To: []string{"wh-a"}}, http.StatusOK)
	e.node("b", Node{Kind: KindQueue, To: []string{"a"}}, http.StatusOK)
	e.node("c", Node{Kind: KindInbox, To: []string{"b"}}, http.StatusOK)
	// a → c would close a cycle a → c → b → a
	if out := e.node("a", Node{Kind: KindQueue, To: []string{"wh-a", "c"}}, http.StatusBadRequest); !strings.Contains(string(out), "cycle") {
		t.Fatalf("cycle: %s", out)
	}
	// a node changes kind, keeping its name
	e.node("a", Node{Kind: KindClass, Members: []string{"laser", "wh-a"}}, http.StatusOK)
	var n Node
	json.Unmarshal(e.admin("GET", "/api/nodes/a", nil, http.StatusOK), &n)
	if n.Kind != KindClass || n.Policy != PolicyRoundRobin || n.To != nil {
		t.Fatalf("node %+v", n)
	}
	// sent to: not deleted
	e.admin("DELETE", "/api/nodes/a", nil, http.StatusConflict)
	e.admin("DELETE", "/api/nodes/wh-a", nil, http.StatusConflict)
	e.admin("DELETE", "/api/nodes/laser", nil, http.StatusConflict)
	e.admin("DELETE", "/api/nodes/c", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/b", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/a", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/wh-a", nil, http.StatusNoContent)

	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	if len(topo) != 2 || topo[0].Name != "laser" || topo[0].Kind != KindPrinter || topo[1].Name != "wh-b" || topo[1].Fallback.Printer != "laser" {
		t.Fatalf("topology %+v", topo)
	}
	// a job to a node that is not there
	if res, _ := e.do(nil, "POST", "/api/destinations/nowhere/jobs?name=x", "", testDoc(t, 1)); res.StatusCode != http.StatusNotFound {
		t.Fatalf("printed to nowhere: %d", res.StatusCode)
	}
}

// CUPS clients list the CUPS server nodes as the destinations of the
// server; the other nodes are reached through them.
func TestCupsGetPrinters(t *testing.T) {
	e := newEnv(t, nil)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Title: "第 1 倉庫"}, http.StatusOK)
	e.printer("laser")
	e.node("tokyo", Node{Kind: KindClass, Members: []string{"laser"}}, http.StatusOK)
	e.node("review", Node{Kind: KindInbox}, http.StatusOK)
	e.node("front", Node{Kind: KindCups, Title: "第 1 倉庫の受付", To: []string{"wh-a"}}, http.StatusOK)
	e.cups("office", "tokyo")
	c := &ipp.Client{URL: e.srv.URL + IPPPath, User: "lpstat"}
	ps, err := c.CupsGetPrinters(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, p := range ps {
		names = append(names, p.Name)
		if p.URI != strings.Replace(e.srv.URL, "http://", "ipp://", 1)+IPPPath+"/"+p.Name || !p.Accepting || p.State != ipp.PrinterIdle || p.IsClass {
			t.Errorf("printer %+v", p)
		}
		switch p.Name {
		case "front":
			if p.Info != "第 1 倉庫の受付" || p.DeviceURI != "cups://front" {
				t.Errorf("front %+v", p)
			}
		case "office":
			if p.Info != "office" || p.DeviceURI != "cups://office" {
				t.Errorf("office %+v", p)
			}
		}
	}
	if !slices.Equal(names, []string{"front", "office"}) {
		t.Fatalf("names %v", names)
	}
	// a class is no IPP printer: CUPS-Get-Classes lists none
	if res, err := c.Do(context.Background(), c.NewRequest(ipp.OpCupsGetClasses), nil); !ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
		t.Fatalf("CUPS-Get-Classes: %+v %v", res, err)
	}
	// the server itself is not a printer; a CUPS server node is one, and
	// a class none
	if _, err := c.GetPrinterAttributes(context.Background()); err == nil {
		t.Fatal("the server answered as a printer")
	}
	pa, err := e.ippClient("office").GetPrinterAttributes(context.Background())
	if err != nil || pa.Name != "office" {
		t.Fatalf("CUPS server attributes %+v %v", pa, err)
	}
	if _, err := e.ippClient("tokyo").GetPrinterAttributes(context.Background()); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("class attributes: %v", err)
	}
}

// An endpoint takes documents and filled templates from an application
// with its own key, and the job flows from it.
func TestEndpoint(t *testing.T) {
	e := newEnv(t, nil)
	ctx := context.Background()
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
	d := e.enroll("wh-a", "tablet")
	d.board()

	// the key is shown once
	var made struct {
		Node
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("intake", Node{Kind: KindEndpoint, Title: "WMS", To: []string{"wh-a"}}, http.StatusOK), &made)
	if !strings.HasPrefix(made.Key, "bdfk_") || made.KeyHash != "" || made.Kind != KindEndpoint {
		t.Fatalf("endpoint %+v", made)
	}
	var again struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("intake", Node{Kind: KindEndpoint, Title: "WMS 2", To: []string{"wh-a"}}, http.StatusOK), &again)
	if again.Key != "" {
		t.Fatal("the key was shown twice")
	}
	var topo []TopoNode
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &topo)
	if topo[0].Name != "intake" || topo[0].In != "/in/intake" || topo[0].Title != "WMS 2" {
		t.Fatalf("topology %+v", topo)
	}

	// a document, with the key; not without, nor with another
	doc := testDoc(t, 2)
	if res, _ := e.do(nil, "POST", "/in/intake?name=pick", "", doc); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("no key: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/in/intake?name=pick", "bdfk_wrong", doc); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("a wrong key: %d", res.StatusCode)
	}
	res, out := e.do(nil, "POST", "/in/intake?name=pick", made.Key, doc)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("in: %d %s", res.StatusCode, out)
	}
	var j Job
	json.Unmarshal(out, &j)
	if j.Dest != "intake" || j.Source != "endpoint" || j.SubmittedBy.Kind != PrincipalEndpoint || j.SubmittedBy.ID != "intake" || j.Pages != 2 || legs(&j) != "intake:routed wh-a:shown" {
		t.Fatalf("job %+v: %s", j, legs(&j))
	}
	if b := d.board(); b == nil || len(b.Open) != 1 || b.Open[0].Name != "pick" {
		t.Fatalf("board %+v", b)
	}

	// a template node fills its templates with the data an endpoint takes
	tpl, err := os.ReadFile("../testdata/form/pages.bdft")
	if err != nil {
		t.Fatal(err)
	}
	// its templates are uploaded once it is there: none named yet
	e.node("forms", Node{Kind: KindTemplate, To: []string{"wh-a"}}, http.StatusOK)
	e.node("forms", Node{Kind: KindTemplate, Template: "slip", TemplateFrom: "doc_type", To: []string{"wh-a"}}, http.StatusBadRequest)
	e.node("forms", Node{Kind: KindTemplate, Template: "no.dot", To: []string{"wh-a"}}, http.StatusBadRequest)
	e.node("forms", Node{Kind: KindTemplate, Template: "slip", To: []string{"wh-a"}}, http.StatusOK)
	var tp Template
	json.Unmarshal(e.admin("PUT", "/api/nodes/forms/templates/slip?name=pages.bdft", tpl, http.StatusOK), &tp)
	if tp.Format != "xlsx" || tp.Node != "forms" || len(tp.Fields) == 0 || tp.Fields[0].Name != "pages" || tp.Size != int64(len(tpl)) {
		t.Fatalf("template %+v", tp)
	}
	var tps []Template
	json.Unmarshal(e.admin("GET", "/api/nodes/forms/templates", nil, http.StatusOK), &tps)
	if len(tps) != 1 || tps[0].ID != "slip" {
		t.Fatalf("templates %+v", tps)
	}
	if res, out := e.do(nil, "GET", "/api/nodes/forms/templates/slip/file", adminToken, nil); res.StatusCode != http.StatusOK || !bytes.Equal(out, tpl) {
		t.Fatalf("file: %d", res.StatusCode)
	}
	e.admin("PUT", "/api/nodes/forms/templates/no.dot?name=pages.bdft", tpl, http.StatusBadRequest)
	e.admin("PUT", "/api/nodes/wh-a/templates/slip?name=pages.bdft", tpl, http.StatusNotFound)
	if res, _ := e.do(nil, "PUT", "/api/nodes/forms/templates/bad", adminToken, []byte("not a template")); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("not a template: %d", res.StatusCode)
	}
	var din struct {
		Key string `json:"key"`
	}
	json.Unmarshal(e.node("data-in", Node{Kind: KindEndpoint, Accept: AcceptData, To: []string{"forms"}}, http.StatusOK), &din)
	var typed []TopoNode // a fresh slice: Unmarshal keeps the old elements' fields
	json.Unmarshal(e.admin("GET", "/api/topology", nil, http.StatusOK), &typed)
	byName := map[string]TopoNode{}
	for _, n := range typed {
		byName[n.Name] = n
	}
	if f := byName["forms"]; !slices.Equal(f.Templates, []string{"slip"}) || f.Takes != PayloadData || !slices.Equal(f.Emits, []string{PayloadBDF}) {
		t.Fatalf("topology forms %+v", f)
	}
	if d := byName["data-in"]; d.Takes != PayloadAny || !slices.Equal(d.Emits, []string{PayloadData}) || d.Output != "" {
		t.Fatalf("topology data-in %+v", d)
	}
	data := map[string]any{"pages": []map[string]any{{"slip": "No. 7", "customer": "ACME", "carry_in": 10}, {"slip": "No. 7", "customer": "ACME", "carry_in": 20}}}
	res, out = e.do(nil, "POST", "/in/data-in?name=slip-7.json", din.Key, data)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("data: %d %s", res.StatusCode, out)
	}
	json.Unmarshal(out, &j)
	if j.Kind != PayloadData || j.Format != "json" || j.Pages != 0 || j.State != JobProcessing || legs(&j) != "data-in:routed forms:pending" || len(j.Payloads) != 1 || j.Payloads[0].Kind != PayloadData {
		t.Fatalf("data job %+v: %s", j, legs(&j))
	}
	if rep := e.tick(); rep.Delivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	filled := e.job(j.Number)
	if filled.State != JobCompleted || legs(filled) != "data-in:routed forms:done wh-a:shown" || len(filled.Payloads) != 2 {
		t.Fatalf("filled job %+v: %s", filled, legs(filled))
	}
	if made := filled.Payloads[1]; made.Kind != PayloadBDF || made.Format != "bdf" || made.Pages != 2 || made.Made != "forms" || made.Name != "slip-7.bdf" || made.ID != filled.Legs[1].Ref ||
		filled.Legs[1].Made != made.ID || filled.Legs[2].Payload != made.ID || filled.Legs[0].Payload != filled.ID {
		t.Fatalf("payloads %+v legs %+v", filled.Payloads, filled.Legs)
	}
	if b := d.board(); b == nil || len(b.Open) != 2 || b.Open[1].Name != "slip-7.json" || b.Open[1].Pages != 2 || b.Open[1].Payload != filled.Payloads[1].ID {
		t.Fatalf("board %+v", b)
	}
	// the document to show is the one made; the data is what was received
	if res, out := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/document", j.Number), adminToken, nil); res.StatusCode != http.StatusOK || !bytes.HasPrefix(out, bdf.Magic[:]) || res.Header.Get("Content-Type") != "application/x-bdf" {
		t.Fatalf("document: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if res, out := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/document?payload=%s", j.Number, j.ID), adminToken, nil); res.StatusCode != http.StatusOK || !bytes.HasPrefix(out, []byte("{")) || res.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("data: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	if res, _ := e.do(nil, "GET", fmt.Sprintf("/api/jobs/%d/pages/1.png?payload=%s", j.Number, j.ID), adminToken, nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("a page of data: %d", res.StatusCode)
	}
	// a document is refused at a data endpoint, and data at a document one
	if res, _ := e.do(nil, "POST", "/in/data-in?name=pick", din.Key, doc); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("a document at a data endpoint: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/in/intake", made.Key, data); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("data at a document endpoint: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/api/destinations/forms/jobs?name=x", adminToken, doc); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("a document to a template node: %d", res.StatusCode)
	}
	// no template of the id: the step fails for good, at once
	e.node("forms", Node{Kind: KindTemplate, Template: "nowhere", To: []string{"wh-a"}}, http.StatusOK)
	_, out = e.do(nil, "POST", "/in/data-in", din.Key, data)
	json.Unmarshal(out, &j)
	if rep := e.tick(); rep.Undelivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	if got := e.job(j.Number); got.State != JobAborted || !strings.Contains(got.Legs[1].Error, "no template") {
		t.Fatalf("no template %+v", got)
	}
	// the metadata names the template (a header's key is in lower case)
	e.node("forms", Node{Kind: KindTemplate, TemplateFrom: "doc-type", To: []string{"wh-a"}}, http.StatusOK)
	res, out = e.doWith(func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+din.Key)
		r.Header.Set("X-Meta-Doc-Type", "slip")
	}, "POST", "/in/data-in?name=by-meta", data)
	json.Unmarshal(out, &j)
	if res.StatusCode != http.StatusCreated || j.Meta["doc-type"] != "slip" || j.Meta["$entry"] != "data-in" || j.Meta["$sender"] != "data-in" || j.Meta["$sender_kind"] != PrincipalEndpoint || j.Meta["$file"] != "by-meta" || j.Meta["$source"] != "endpoint" {
		t.Fatalf("meta: %d %+v", res.StatusCode, j.Meta)
	}
	if rep := e.tick(); rep.Delivered != 1 {
		t.Fatalf("tick %+v", rep)
	}
	if got := e.job(j.Number); got.State != JobCompleted || legs(got) != "data-in:routed forms:done wh-a:shown" {
		t.Fatalf("by meta %+v: %s", got, legs(got))
	}
	_, out = e.do(nil, "POST", "/in/data-in", din.Key, data)
	json.Unmarshal(out, &j)
	e.tick()
	if got := e.job(j.Number); got.State != JobAborted || !strings.Contains(got.Legs[1].Error, "doc-type") {
		t.Fatalf("no doc-type %+v", got)
	}
	// the templates go with their node
	e.admin("DELETE", "/api/nodes/forms", nil, http.StatusConflict) // data-in sends to it
	e.admin("DELETE", "/api/nodes/data-in", nil, http.StatusNoContent)
	e.admin("DELETE", "/api/nodes/forms", nil, http.StatusNoContent)
	e.admin("GET", "/api/nodes/forms/templates", nil, http.StatusNotFound)
	if _, err := e.m().Template(ctx, "forms", "slip"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("forms' template stayed: %v", err)
	}

	// what an endpoint takes: some formats
	e.node("intake", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Accept: AcceptFile, Formats: []string{"PDF"}}, http.StatusOK)
	if res, _ := e.do(nil, "POST", "/in/intake?name=pick", made.Key, doc); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("bdf at a pdf endpoint: %d", res.StatusCode)
	}
	e.node("intake", Node{Kind: KindEndpoint, To: []string{"wh-a"}, Formats: []string{"bdf", "pdf"}}, http.StatusOK)
	if res, _ := e.do(nil, "POST", "/in/intake?name=pick", made.Key, doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("bdf at a bdf endpoint: %d", res.StatusCode)
	}
	var got Node
	json.Unmarshal(e.admin("GET", "/api/nodes/intake", nil, http.StatusOK), &got)
	if got.Accept != AcceptFile || got.Output != OutputBDF || !slices.Equal(got.Formats, []string{"bdf", "pdf"}) || got.KeyHash != "" {
		t.Fatalf("node %+v", got)
	}
	// the refusals are in the audit log
	var log []AuditEntry
	json.Unmarshal(e.admin("GET", "/api/audit?limit=50", nil, http.StatusOK), &log)
	refused := 0
	for _, en := range log {
		if en.Action == "refuse-document" && en.Subject == "intake" && en.Who == "intake" {
			refused++
		}
	}
	if refused != 2 {
		t.Fatalf("refusals audited: %d", refused)
	}

	// the key rotated: the old one is refused
	var rot map[string]string
	json.Unmarshal(e.admin("POST", "/api/nodes/intake/key", nil, http.StatusOK), &rot)
	if res, _ := e.do(nil, "POST", "/in/intake", made.Key, doc); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("the old key: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/in/intake", rot["key"], doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("the new key: %d", res.StatusCode)
	}
	e.admin("POST", "/api/nodes/wh-a/key", nil, http.StatusBadRequest)
	// the endpoint is a destination of the API too; IPP reaches it
	// through a CUPS server alone
	if js := e.printJob("intake", "api", 1); js.State != JobCompleted || legs(js) != "intake:routed wh-a:shown" {
		t.Fatalf("api %+v: %s", js, legs(js))
	}
	if _, err := e.ippClient("intake").PrintJob(ctx, "ipp", "application/x-bdf", nil, bytes.NewReader(doc)); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("ipp to the endpoint: %v", err)
	}
	if js, err := e.cups("cups1", "intake").PrintJob(ctx, "ipp", "application/x-bdf", nil, bytes.NewReader(doc)); err != nil || js.State != ipp.JobCompleted {
		t.Fatalf("ipp through a CUPS server: %+v %v", js, err)
	} else if j := e.job(js.ID); legs(j) != "cups1:routed intake:routed wh-a:shown" {
		t.Fatalf("job %+v: %s", j, legs(j))
	}

	// the layout of the canvas
	var l Layout
	json.Unmarshal(e.admin("PUT", "/api/layout", Layout{Positions: map[string]Point{"intake": {X: 40, Y: 80}, "wh-a": {X: 300, Y: 80}, "gone": {X: 1, Y: 1}}}, http.StatusOK), &l)
	if len(l.Positions) != 2 || l.Positions["wh-a"].X != 300 {
		t.Fatalf("layout %+v", l)
	}
	json.Unmarshal(e.admin("PUT", "/api/layout", Layout{Positions: map[string]Point{"wh-a": {X: 320, Y: 90}}}, http.StatusOK), &l)
	json.Unmarshal(e.admin("GET", "/api/layout", nil, http.StatusOK), &l)
	if l.Positions["intake"].X != 40 || l.Positions["wh-a"].X != 320 {
		t.Fatalf("layout %+v", l)
	}
}

func TestKeys(t *testing.T) {
	key := newKey()
	h, err := hashKey(key)
	if err != nil || !strings.HasPrefix(h, "$argon2id$v=19$m=19456,t=2,p=1$") {
		t.Fatalf("%s %v", h, err)
	}
	if !checkKey(h, key) || checkKey(h, key+"x") || checkKey("", key) || checkKey("$argon2id$v=19$m=0,t=0,p=0$AA$AA", key) {
		t.Fatal("checkKey")
	}
}

// Scripts use the admin token or a submit token: a token's user sees and
// cancels the jobs it sent and prints; an administrator sees every job
// and the audit log. (People sign in to the pages with OIDC: scope_test.)
func TestTokens(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.SubmitTokens = map[string]string{"wms": "s3cret", "erp": "pw"} })
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
	if res, _ := e.do(nil, "GET", "/api/me", "nobody:x", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("me with a wrong token: %d", res.StatusCode)
	}
	var me struct {
		Name  string
		Admin bool
	}
	_, out := e.do(nil, "GET", "/api/me", "wms:s3cret", nil)
	json.Unmarshal(out, &me)
	if me.Name != "wms" || me.Admin {
		t.Fatalf("user %+v", me)
	}
	_, out = e.do(nil, "GET", "/api/me", adminToken, nil)
	json.Unmarshal(out, &me)
	if me.Name != "admin" || !me.Admin {
		t.Fatalf("admin %+v", me)
	}

	// the user prints with the token, and sees its jobs only
	doc := testDoc(t, 1)
	if res, out := e.do(nil, "POST", "/api/destinations/wh-a/jobs?name=mine", "wms:s3cret", doc); res.StatusCode != http.StatusCreated {
		t.Fatalf("print as a user: %d %s", res.StatusCode, out)
	}
	e.do(nil, "POST", "/api/destinations/wh-a/jobs?name=erps", "erp:pw", doc)
	var jobs []Job
	_, out = e.do(nil, "GET", "/api/my/jobs", "wms:s3cret", nil)
	json.Unmarshal(out, &jobs)
	if len(jobs) != 1 || jobs[0].Name != "mine" || *jobs[0].SubmittedBy != (Principal{Kind: PrincipalToken, ID: "wms", Name: "wms", Via: ViaSubmitToken}) {
		t.Fatalf("my jobs %+v", jobs)
	}
	_, out = e.do(nil, "GET", "/api/my/jobs", adminToken, nil)
	json.Unmarshal(out, &jobs)
	if len(jobs) != 2 || jobs[0].Name != "erps" {
		t.Fatalf("admin's jobs %+v", jobs)
	}
	// another's job is not there for a user; their own is
	if res, _ := e.do(nil, "GET", "/api/jobs/2", "wms:s3cret", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("another's job: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "GET", "/api/jobs/1/document", "wms:s3cret", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("own document: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/api/jobs/1/cancel", "wms:s3cret", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("cancel own job: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "GET", "/api/destinations", "wms:s3cret", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("destinations: %d", res.StatusCode)
	}
	// the admin's pages' API is not the user's
	for _, path := range []string{"/api/deliveries", "/api/audit", "/api/nodes"} {
		if res, _ := e.do(nil, "GET", path, "wms:s3cret", nil); res.StatusCode != http.StatusForbidden {
			t.Fatalf("%s as a user: %d", path, res.StatusCode)
		}
	}

	// the audit log, newest first
	var log []AuditEntry
	e.admin("PUT", "/api/nodes/q", Node{Kind: KindQueue, To: []string{"wh-a"}}, http.StatusOK)
	json.Unmarshal(e.admin("GET", "/api/audit?limit=100", nil, http.StatusOK), &log)
	var actions []string
	for _, en := range log {
		actions = append(actions, en.Who+":"+en.Action)
	}
	want := []string{"admin:put-node", "wms:cancel-job", "admin:put-node"}
	if !slices.Equal(actions, want) {
		t.Fatalf("audit %v", actions)
	}
	if log[0].Subject != "q" || log[1].Subject != jobID(1) {
		t.Fatalf("audit subjects %+v", log[:2])
	}
}

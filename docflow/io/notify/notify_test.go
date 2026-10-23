package notify_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	_ "github.com/shibukawa/bdf/docflow/io/http"
	"github.com/shibukawa/bdf/docflow/io/notify"
	_ "github.com/shibukawa/bdf/docflow/render"
)

// A send that fails for good goes along the node's on_error edge to a
// notify node, which tells Slack of the error and sends the job on with
// the document that reached the failing node.
func TestOnError(t *testing.T) {
	slack := flowtest.NewHook(t)
	refuse := flowtest.NewHook(t, nethttp.StatusBadRequest)
	log := flowtest.NewHook(t)
	t.Setenv(docflow.CredentialEnv("slack-ops"), slack.URL+"/services/T000/B000/xyz")
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "log", Kind: "http", URL: log.URL},
		docflow.Node{Name: "alert", Kind: notify.KindNotify, Via: notify.ViaSlack, Credential: "slack-ops", To: []string{"log"}},
		docflow.Node{Name: "erp", Kind: "http", URL: refuse.URL, OnError: []string{"alert"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"erp"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "order.bdf", doc, nil)
	e.Tick()
	got := slack.Got()
	if len(got) != 1 {
		t.Fatalf("%d notifications: %+v", len(got), e.Job(j.Number).Legs)
	}
	var msg struct{ Text string }
	json.Unmarshal(got[0].Body, &msg)
	if !strings.HasPrefix(msg.Text, "Job 1 (order.bdf) failed at erp: ") || !strings.Contains(msg.Text, "400") {
		t.Errorf("text %q", msg.Text)
	}
	if l := log.Got(); len(l) != 1 || !bytes.Equal(l[0].Body, doc) {
		t.Errorf("sent on: %d", len(l))
	}
	job := e.Job(j.Number)
	if job.State != docflow.JobAborted {
		t.Errorf("job %s", job.State)
	}
	d := e.Delivery(j, "alert")
	if d.Cause == nil || d.Cause.Node != "erp" || d.State != docflow.DeliveryDone {
		t.Errorf("alert %+v", d)
	}
	var froms []string
	for _, l := range job.Legs {
		froms = append(froms, l.From+">"+l.Node+":"+l.State)
	}
	if strings.Join(froms, " ") != ">in:routed in>erp:failed erp>alert:done alert>log:done" {
		t.Errorf("legs %v", froms)
	}
}

// A node with no on_error edges sends its failures to the manager's
// default; a notify node on an error's way that fails does not tell of
// itself; with neither, a failure stops the job as before.
func TestOnErrorDefault(t *testing.T) {
	hook := flowtest.NewHook(t)
	refuse := flowtest.NewHook(t, nethttp.StatusBadRequest)
	e := flowtest.New(t, func(o *docflow.Options) { o.OnError = []string{"ops"} })
	e.Put(docflow.Node{Name: "ops", Kind: notify.KindNotify, Via: notify.ViaHTTP, URL: hook.URL, Text: "{{meta.$error_node}} failed for {{meta.customer}}"},
		docflow.Node{Name: "broken", Kind: notify.KindNotify, Via: notify.ViaHTTP, URL: refuse.URL},
		docflow.Node{Name: "erp", Kind: "http", URL: refuse.URL},
		docflow.Node{Name: "crm", Kind: "http", URL: refuse.URL, OnError: []string{"broken"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"erp", "crm"}})
	j := e.Submit("in", "x", flowtest.Doc(t, 1), map[string]any{"customer": "acme"})
	e.Tick()
	got := hook.Got()
	if len(got) != 1 {
		t.Fatalf("%d notifications: %+v", len(got), e.Job(j.Number).Legs)
	}
	var body map[string]any
	json.Unmarshal(got[0].Body, &body)
	if body["text"] != "erp failed for acme" || body["error_node"] != "erp" || body["node"] != "ops" || body["job"] != 1.0 {
		t.Errorf("body %v", body)
	}
	if got[0].Header.Get("Idempotency-Key") == "" {
		t.Error("no idempotency key")
	}
	// crm → broken (fails) → not the default again
	if d := e.Delivery(j, "broken"); d.State != docflow.DeliveryFailed || d.Cause.Node != "crm" {
		t.Errorf("broken %s %+v", d.State, d.Cause)
	}
	if n := len(refuse.Got()); n != 3 {
		t.Errorf("%d refused", n)
	}

	plain := flowtest.New(t, nil)
	plain.Put(docflow.Node{Name: "erp", Kind: "http", URL: refuse.URL}, docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"erp"}})
	j = plain.Submit("in", "x", flowtest.Doc(t, 1), nil)
	plain.Tick()
	if job := plain.Job(j.Number); job.State != docflow.JobAborted || len(job.Legs) != 2 {
		t.Errorf("%s %+v", job.State, job.Legs)
	}
}

// A job an assert right after its entry refuses goes to the assert's
// on_error notify node, which tells what failed; the sender is refused.
func TestNotifyRefused(t *testing.T) {
	slack := flowtest.NewHook(t)
	t.Setenv(docflow.CredentialEnv("slack-ops"), slack.URL+"/services/T000/B000/xyz")
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "seen", Kind: docflow.KindSink},
		docflow.Node{Name: "alert", Kind: notify.KindNotify, Via: notify.ViaSlack, Credential: "slack-ops"},
		docflow.Node{Name: "check", Kind: docflow.KindAssert, To: []string{"seen"}, OnError: []string{"alert"}, Checks: []docflow.Check{{Key: "doc-type"}}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"check"}})
	doc := flowtest.Doc(t, 1)
	_, err := e.M().Submit(context.Background(), "in", docflow.Input{Name: "order.bdf", Source: "api", Doc: bytes.NewReader(doc), Size: int64(len(doc))})
	var ref *docflow.Refusal
	if !errors.As(err, &ref) || ref.Job == nil || len(ref.Failures) != 1 || ref.Failures[0].Key != "doc-type" {
		t.Fatalf("refused: %v", err)
	}
	e.Tick()
	got := slack.Got()
	if len(got) != 1 {
		t.Fatalf("%d notifications: %+v", len(got), e.Job(ref.Job.Number).Legs)
	}
	var msg struct{ Text string }
	json.Unmarshal(got[0].Body, &msg)
	if msg.Text != "Job 1 (order.bdf) failed at check: doc-type: required" {
		t.Errorf("text %q", msg.Text)
	}
	if j := e.Job(ref.Job.Number); j.State != docflow.JobAborted {
		t.Errorf("job %s", j.State)
	}
}

// A notify node passes on what reaches it, telling that it did.
func TestNotifyPasses(t *testing.T) {
	slack := flowtest.NewHook(t)
	out := flowtest.NewHook(t)
	t.Setenv(docflow.CredentialEnv("slack"), slack.URL)
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "out", Kind: "http", URL: out.URL},
		docflow.Node{Name: "told", Kind: notify.KindNotify, Via: notify.ViaSlack, Credential: "slack", To: []string{"out"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, Accept: docflow.AcceptData, To: []string{"told"}})
	j := e.Submit("in", "rows.json", []byte(`{"a": 1}`), nil)
	e.Tick()
	var msg struct{ Text string }
	if g := slack.Got(); len(g) == 1 {
		json.Unmarshal(g[0].Body, &msg)
	}
	if msg.Text != "Job 1 (rows.json) passed told" {
		t.Errorf("text %q", msg.Text)
	}
	if g := out.Got(); len(g) != 1 || string(g[0].Body) != `{"a": 1}` {
		t.Errorf("out %d", len(g))
	}
	if job := e.Job(j.Number); job.State != docflow.JobCompleted {
		t.Errorf("job %s", job.State)
	}
}

// An on_error edge carries what reaches its node: data from a data entry
// cannot go along it to a render node, which takes .bdf. The settings of
// a notify node are checked.
func TestOnErrorTyped(t *testing.T) {
	hook := flowtest.NewHook(t)
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "out", Kind: "http", URL: hook.URL},
		docflow.Node{Name: "pdf", Kind: "render", To: []string{"out"}},
		docflow.Node{Name: "erp", Kind: "http", URL: hook.URL},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, Accept: docflow.AcceptData, To: []string{"erp"}})
	ctx := context.Background()
	_, err := e.M().PutNode(ctx, docflow.Node{Name: "erp", Kind: "http", URL: hook.URL, OnError: []string{"pdf"}})
	if !errors.Is(err, docflow.ErrInvalid) || !strings.Contains(err.Error(), "in → erp → pdf") {
		t.Errorf("data along on_error to render: %v", err)
	}
	for _, n := range []docflow.Node{
		{Name: "n", Kind: notify.KindNotify, Via: "fax"},
		{Name: "n", Kind: notify.KindNotify, Via: notify.ViaSlack},
		{Name: "n", Kind: notify.KindNotify, Via: notify.ViaHTTP, URL: "mailto:x"},
		{Name: "n", Kind: notify.KindNotify, Via: notify.ViaMail, URL: "smtp://localhost"},
		{Name: "n", Kind: notify.KindNotify, Via: notify.ViaSlack, Credential: "s", OnError: []string{"nowhere"}},
	} {
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", n, err)
		}
	}
}

// A notify node by mail sends the text, with the first line as the
// subject when it has none.
func TestNotifyMail(t *testing.T) {
	srv := flowtest.NewSMTP(t)
	refuse := flowtest.NewHook(t, nethttp.StatusBadRequest)
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "mail", Kind: notify.KindNotify, Via: notify.ViaMail, URL: "smtp://" + srv.Addr, From: "docflow@example.com", Recipients: []string{"ops@example.com"}},
		docflow.Node{Name: "erp", Kind: "http", URL: refuse.URL, OnError: []string{"mail"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"erp"}})
	j := e.Submit("in", "x", flowtest.Doc(t, 1), nil)
	e.Tick()
	mails := srv.Mails()
	if len(mails) != 1 {
		t.Fatalf("%d mails", len(mails))
	}
	if s := flowtest.Header(mails[0].Data, "Subject"); !strings.Contains(s, "failed") || flowtest.Header(mails[0].Data, "Message-Id") != "<"+e.Delivery(j, "mail").ID+"@bdf-docflow>" {
		t.Errorf("subject %q, message id %q", s, flowtest.Header(mails[0].Data, "Message-Id"))
	}
}

// The settings the editor marks as required are those Check refuses a
// notify node without: how it notifies, and what that way needs.
func TestRequired(t *testing.T) {
	e := flowtest.New(t, nil)
	e.Required(docflow.Node{Name: "slack", Kind: notify.KindNotify, Via: notify.ViaSlack, Credential: "ops-webhook"})
	e.Required(docflow.Node{Name: "mail", Kind: notify.KindNotify, Via: notify.ViaMail, URL: "smtp://127.0.0.1:2525", From: "d@example.com", Recipients: []string{"a@example.com"}})
	e.Required(docflow.Node{Name: "web", Kind: notify.KindNotify, Via: notify.ViaHTTP, URL: "https://example.com/hook"})
}

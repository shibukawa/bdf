package http_test

import (
	"bytes"
	"context"
	"errors"
	nethttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	"github.com/shibukawa/bdf/docflow/io/http"
)

// A document sent to a web service: the URL and a header expanded, the
// credential (a file) as a bearer token, the sender's metadata as
// X-Meta-*, the delivery's id as the Idempotency-Key, the same on the try
// after a 503.
func TestHTTPOut(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "erp"), []byte("tok-123\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	hook := flowtest.NewHook(t, nethttp.StatusServiceUnavailable, nethttp.StatusCreated)
	e := flowtest.New(t, func(o *docflow.Options) { o.CredentialDir = dir })
	e.Put(docflow.Node{Name: "erp", Kind: http.KindHTTP, URL: hook.URL + "/in/remote?name={{meta.$name}}&customer={{meta.customer}}", Method: "put", Credential: "erp",
		Headers: map[string]string{"X-Order": "{{meta.order}}"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"erp"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "a b.bdf", doc, map[string]any{"customer": "a&b", "order": "7\r\nX-Evil: 1", "doc/type": "x"})
	e.Tick()
	d := e.Delivery(j, "erp")
	if d.State != docflow.DeliveryPending || !strings.Contains(d.LastError, "503") {
		t.Fatalf("after 503: %s %s", d.State, d.LastError)
	}
	e.Advance(time.Minute)
	e.Tick()
	if d = e.Delivery(j, "erp"); d.State != docflow.DeliveryDone || d.Result != "201 Created" {
		t.Fatalf("%s %s %s", d.State, d.Result, d.LastError)
	}
	got := hook.Got()
	if len(got) != 2 {
		t.Fatalf("%d requests", len(got))
	}
	r := got[1]
	if r.Method != "PUT" || r.Path != "/in/remote" || r.Query != "name=a+b.bdf&customer=a%26b" {
		t.Errorf("%s %s?%s", r.Method, r.Path, r.Query)
	}
	if r.Header.Get("Authorization") != "Bearer tok-123" || r.Header.Get("Idempotency-Key") != d.ID || got[0].Header.Get("Idempotency-Key") != d.ID {
		t.Errorf("headers %v", r.Header)
	}
	if r.Header.Get("X-Order") != "7  X-Evil: 1" || r.Header.Get("X-Evil") != "" || r.Header.Get("X-Meta-Customer") != "a&b" || r.Header.Get("Content-Type") != "application/vnd.bdf" {
		t.Errorf("headers %v", r.Header)
	}
	if !bytes.Equal(r.Body, doc) {
		t.Error("body")
	}
}

// A refusal (4xx) and a redirect fail for good; the settings are
// checked.
func TestHTTPRefused(t *testing.T) {
	hook := flowtest.NewHook(t, nethttp.StatusBadRequest)
	moved := flowtest.NewHook(t, nethttp.StatusFound)
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "bad", Kind: http.KindHTTP, URL: hook.URL},
		docflow.Node{Name: "moved", Kind: http.KindHTTP, URL: moved.URL},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"bad", "moved"}})
	j := e.Submit("in", "x", flowtest.Doc(t, 1), nil)
	e.Tick()
	for _, node := range []string{"bad", "moved"} {
		if d := e.Delivery(j, node); d.State != docflow.DeliveryFailed || d.Attempts != 1 {
			t.Errorf("%s: %s after %d: %s", node, d.State, d.Attempts, d.LastError)
		}
	}
	if job := e.Job(j.Number); job.State != docflow.JobAborted {
		t.Errorf("job %s", job.State)
	}
	ctx := context.Background()
	for _, n := range []docflow.Node{
		{Name: "x", Kind: http.KindHTTP, URL: "ftp://x"},
		{Name: "x", Kind: http.KindHTTP, URL: hook.URL, Method: "DELETE"},
		{Name: "x", Kind: http.KindHTTP, URL: hook.URL, Headers: map[string]string{"Authorization": "x"}},
		{Name: "x", Kind: http.KindHTTP, URL: hook.URL, Headers: map[string]string{"Bad Header": "x"}},
	} {
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", n, err)
		}
	}
}

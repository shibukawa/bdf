package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	nethttp "net/http"
	"os"
	"path/filepath"
	"reflect"
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

// A multipart body: the metadata the sender gave as JSON in the part
// meta (its values as given, none of the manager's), the document in the
// part file with its name and type; no X-Meta-* headers, and the node's
// own headers still sent.
func TestHTTPOutMultipart(t *testing.T) {
	hook := flowtest.NewHook(t, nethttp.StatusOK)
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "erp", Kind: http.KindHTTP, URL: hook.URL + "/upload", Body: "Multipart", Headers: map[string]string{"X-Order": "{{meta.order}}"}},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"erp"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", `a "b".bdf`, doc, map[string]any{"customer": "a&b", "order": "7\r\nX-Evil: 1", "count": 3, "tags": []any{"x", "y"}})
	e.Tick()
	if d := e.Delivery(j, "erp"); d.State != docflow.DeliveryDone {
		t.Fatalf("%s %s", d.State, d.LastError)
	}
	got := hook.Got()
	if len(got) != 1 {
		t.Fatalf("%d requests", len(got))
	}
	r := got[0]
	for k := range r.Header {
		if strings.HasPrefix(k, "X-Meta-") || k == "Content-Disposition" {
			t.Errorf("header %s", k)
		}
	}
	if r.Header.Get("X-Order") != "7  X-Evil: 1" || r.Header.Get("Idempotency-Key") == "" {
		t.Errorf("headers %v", r.Header)
	}
	typ, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || typ != "multipart/form-data" {
		t.Fatalf("Content-Type %q: %v", r.Header.Get("Content-Type"), err)
	}
	mr := multipart.NewReader(bytes.NewReader(r.Body), params["boundary"])
	p, err := mr.NextPart()
	if err != nil || p.FormName() != "meta" || p.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("first part %v: %v", p, err)
	}
	var meta map[string]any
	if err := json.NewDecoder(p).Decode(&meta); err != nil {
		t.Fatal(err)
	}
	if want := map[string]any{"customer": "a&b", "order": "7\r\nX-Evil: 1", "count": 3.0, "tags": []any{"x", "y"}}; !reflect.DeepEqual(meta, want) {
		t.Errorf("meta %v", meta)
	}
	if p, err = mr.NextPart(); err != nil || p.FormName() != "file" || p.FileName() != `a "b".bdf` || p.Header.Get("Content-Type") != "application/vnd.bdf" {
		t.Fatalf("second part %v: %v", p, err)
	}
	if b, _ := io.ReadAll(p); !bytes.Equal(b, doc) {
		t.Error("file")
	}
	if _, err := mr.NextPart(); err != io.EOF {
		t.Errorf("a third part: %v", err)
	}
}

// To another Docflow's endpoint, with the endpoint's key: either body
// arrives as the same document, with its name and the sender's
// metadata (a number stays one in multipart, a header is text).
func TestHTTPToDocflow(t *testing.T) {
	dir := t.TempDir()
	e := flowtest.New(t, func(o *docflow.Options) { o.CredentialDir = dir })
	key := e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink}, docflow.Node{Name: "remote", Kind: docflow.KindEndpoint, To: []string{"keep"}})
	if err := os.WriteFile(filepath.Join(dir, "remote"), []byte(key), 0o600); err != nil {
		t.Fatal(err)
	}
	e.Put(docflow.Node{Name: "direct", Kind: http.KindHTTP, URL: e.Srv.URL + "/in/remote", Credential: "remote"},
		docflow.Node{Name: "multi", Kind: http.KindHTTP, URL: e.Srv.URL + "/in/remote", Credential: "remote", Body: http.BodyMultipart},
		docflow.Node{Name: "in", Kind: docflow.KindEndpoint, To: []string{"direct", "multi"}})
	doc := flowtest.Doc(t, 2)
	j := e.Submit("in", "請求書.bdf", doc, map[string]any{"customer": "a&b", "count": 3})
	e.Tick()
	want := map[string]any{"direct": "3", "multi": 3.0}
	for _, node := range []string{"direct", "multi"} {
		if d := e.Delivery(j, node); d.State != docflow.DeliveryDone || d.Result != "201 Created" {
			t.Fatalf("%s: %s %s %s", node, d.State, d.Result, d.LastError)
		}
	}
	got := map[string]any{}
	for n := j.Number + 1; n <= j.Number+2; n++ {
		r := e.Job(n)
		if r.Dest != "remote" || r.Name != "請求書.bdf" || r.Meta["customer"] != "a&b" || r.Meta["$sender_kind"] != docflow.PrincipalEndpoint || !bytes.Equal(e.Payload(r.Payloads[0].ID), doc) {
			t.Fatalf("job %d: %s %q %+v", n, r.Dest, r.Name, r.Meta)
		}
		if _, ok := r.Meta["count"].(string); ok {
			got["direct"] = r.Meta["count"]
		} else {
			got["multi"] = r.Meta["count"]
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("count %v", got)
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
		{Name: "x", Kind: http.KindHTTP, URL: hook.URL, Body: "json"},
		{Name: "x", Kind: http.KindHTTP, URL: hook.URL, Body: http.BodyMultipart, Headers: map[string]string{"content-type": "text/plain"}},
	} {
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", n, err)
		}
	}
}

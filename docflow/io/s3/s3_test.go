package s3_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	"github.com/shibukawa/bdf/docflow/io/s3"
)

// fakeS3 is a bucket's objects behind the S3 API, path-style: PUT (with
// If-None-Match: *) and GET, of requests signed with the access key.
type fakeS3 struct {
	*httptest.Server
	mu      sync.Mutex
	objects map[string][]byte
	puts    int
}

func newS3(t *testing.T) *fakeS3 {
	f := &fakeS3{objects: map[string][]byte{}}
	f.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 Credential=AKIATEST/") || r.Header.Get("X-Amz-Content-Sha256") == "" {
			http.Error(w, "AccessDenied", http.StatusForbidden)
			return
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		key := strings.TrimPrefix(r.URL.Path, "/")
		switch r.Method {
		case http.MethodPut:
			f.puts++
			if _, ok := f.objects[key]; ok && r.Header.Get("If-None-Match") == "*" {
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
			b, _ := io.ReadAll(r.Body)
			f.objects[key] = b
		case http.MethodGet:
			b, ok := f.objects[key]
			if !ok {
				http.Error(w, "NoSuchKey", http.StatusNotFound)
				return
			}
			w.Write(b)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(f.Close)
	return f
}

func (f *fakeS3) object(key string) ([]byte, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, ok := f.objects[key]
	return b, ok
}

func (f *fakeS3) url(bucket string) string {
	return "s3://" + bucket + "?region=us-east-1&endpoint=" + f.URL
}

// An exit writes the document under the key its path expands to, once:
// a value cannot leave the prefix, and an object there already (a write
// whose outcome was not known) is taken as written.
func TestS3Out(t *testing.T) {
	f := newS3(t)
	t.Setenv(docflow.CredentialEnv("s3-test"), "AKIATEST:secret")
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "bucket", Kind: s3.KindS3, URL: f.url("out"), Credential: "s3-test", Path: "orders/{{meta.customer}}/{{meta.$name}}"},
		docflow.Node{Name: "plain", Kind: s3.KindS3, URL: f.url("out"), Credential: "s3-test"},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"bucket", "plain"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "a.bdf", doc, map[string]any{"customer": "../../acme"})
	e.Tick()
	d := e.Delivery(j, "bucket")
	if d.State != docflow.DeliveryDone || d.Result != "s3://out/orders/acme/a.bdf" {
		t.Fatalf("%s %s: %s", d.State, d.Result, d.LastError)
	}
	if b, ok := f.object("out/orders/acme/a.bdf"); !ok || !bytes.Equal(b, doc) {
		t.Errorf("object: %v", ok)
	}
	if d := e.Delivery(j, "plain"); d.Result != "s3://out/1/a.bdf" {
		t.Errorf("default path: %s %s", d.Result, d.LastError)
	}
	// the same key again: 412, done, not overwritten
	j2 := e.Submit("in", "a.bdf", flowtest.Doc(t, 2), map[string]any{"customer": "acme"})
	e.Tick()
	if d := e.Delivery(j2, "bucket"); d.State != docflow.DeliveryDone {
		t.Errorf("412: %s %s", d.State, d.LastError)
	}
	if b, _ := f.object("out/orders/acme/a.bdf"); !bytes.Equal(b, doc) {
		t.Error("overwritten")
	}
	// a refusal (a wrong key: 403) fails for good
	t.Setenv(docflow.CredentialEnv("s3-test"), "WRONG:secret")
	j3 := e.Submit("in", "b.bdf", doc, map[string]any{"customer": "x"})
	e.Tick()
	if d := e.Delivery(j3, "bucket"); d.State != docflow.DeliveryFailed || !strings.Contains(d.LastError, "403") || strings.Contains(d.LastError, "secret") {
		t.Errorf("403: %s %s", d.State, d.LastError)
	}
	ctx := context.Background()
	for _, n := range []docflow.Node{
		{Name: "x", Kind: s3.KindS3, URL: "https://out.s3.amazonaws.com"},
		{Name: "x", Kind: s3.KindS3, URL: "s3://out/prefix?region=us-east-1"},
		{Name: "x", Kind: s3.KindS3, URL: "s3://out?endpoint=ftp://x"},
	} {
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", n, err)
		}
	}
}

// An S3 Trigger takes the object a reference names, in the three shapes,
// with the node's key; it refuses another bucket, keys outside its
// prefix, requests without its key, and objects not there.
func TestS3Trigger(t *testing.T) {
	f := newS3(t)
	t.Setenv(docflow.CredentialEnv("s3-test"), "AKIATEST:secret")
	e := flowtest.New(t, nil)
	doc := flowtest.Doc(t, 1)
	f.objects["in/scans/2026/a b.bdf"] = doc
	f.objects["in/other.bdf"] = doc
	key := e.Put(docflow.Node{Name: "copy", Kind: s3.KindS3, URL: f.url("archive"), Credential: "s3-test", Path: "{{meta.s3.key}}"},
		docflow.Node{Name: "trigger", Kind: s3.KindTrigger, URL: f.url("in"), Credential: "s3-test", Path: "scans", To: []string{"copy"}})
	if key == "" {
		t.Fatal("no key")
	}
	post := func(body, key string) int {
		req, _ := http.NewRequest(http.MethodPost, e.Srv.URL+docflow.NodePath(s3.KindTrigger, "trigger"), strings.NewReader(body))
		if key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		return res.StatusCode
	}
	for i, body := range []string{
		`{"bucket": "in", "key": "scans/2026/a b.bdf"}`,
		`{"Records": [{"s3": {"bucket": {"name": "in"}, "object": {"key": "scans/2026/a+b.bdf"}}}]}`,
		`{"detail-type": "Object Created", "detail": {"bucket": {"name": "in"}, "object": {"key": "scans/2026/a b.bdf"}}}`,
	} {
		if s := post(body, key); s != http.StatusCreated {
			t.Fatalf("%d: %d", i, s)
		}
	}
	e.Tick()
	if b, ok := f.object("archive/scans/2026/a b.bdf"); !ok || !bytes.Equal(b, doc) {
		t.Error("not copied")
	}
	j := e.Job(1)
	if j.Meta["s3.bucket"] != "in" || j.Meta["s3.key"] != "scans/2026/a b.bdf" || j.Name != "a b.bdf" || j.Dest != "trigger" {
		t.Errorf("job %+v", j)
	}
	for body, want := range map[string]int{
		`{"bucket": "other", "key": "scans/x.bdf"}`:     http.StatusForbidden,
		`{"bucket": "in", "key": "other.bdf"}`:          http.StatusForbidden,
		`{"bucket": "in", "key": "scans/none.bdf"}`:     http.StatusNotFound,
		`{"bucket": "in", "key": "scans/../other.bdf"}`: http.StatusBadRequest,
		`{"Records": []}`: http.StatusBadRequest,
		`not json`:        http.StatusBadRequest,
	} {
		if s := post(body, key); s != want {
			t.Errorf("%s: %d, want %d", body, s, want)
		}
	}
	if s := post(`{"bucket": "in", "key": "scans/2026/a b.bdf"}`, "bdfk_wrong"); s != http.StatusUnauthorized {
		t.Errorf("wrong key: %d", s)
	}
	// the key rotates as an uploader's
	res, _ := adminPost(e, "/api/nodes/trigger/key")
	if res != http.StatusOK {
		t.Errorf("rotate: %d", res)
	}
}

func adminPost(e *flowtest.Env, path string) (int, error) {
	req, _ := http.NewRequest(http.MethodPost, e.Srv.URL+path, nil)
	req.Header.Set("Authorization", "Bearer "+flowtest.AdminToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}

// The settings the editor marks as required are those Check refuses a
// node without: the bucket, of either kind.
func TestRequired(t *testing.T) {
	e := flowtest.New(t, nil)
	e.Required(docflow.Node{Name: "copy", Kind: s3.KindS3, URL: "s3://archive?region=ap-northeast-1"})
	e.Required(docflow.Node{Name: "trigger", Kind: s3.KindTrigger, URL: "s3://in?region=ap-northeast-1"})
}

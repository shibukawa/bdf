package segment

import (
	"bytes"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shibukawa/bdf"
)

// book is a document of 25 pages, each with an object of its own.
func book(t *testing.T) *bdf.Reader {
	t.Helper()
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	for i := range 25 {
		o := bdf.NewObject()
		o.FillRect(float32(i), 0, 1, 1)
		h, _ := d.AddObject(o)
		v.AddPage(100, 100, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func post(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	j, _ := json.Marshal(body)
	res, err := http.Post(url, "application/json", bytes.NewReader(j))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { res.Body.Close() })
	return res
}

func TestHandler(t *testing.T) {
	doc := book(t)
	h := &Handler{
		Open: func(r *http.Request) (*bdf.Reader, error) {
			switch r.URL.Path {
			case "/book":
				return doc, nil
			case "/locked":
				return nil, &StatusError{Status: http.StatusUnauthorized, Err: errors.New("log in")}
			case "/theirs":
				return nil, fs.ErrPermission
			}
			return nil, fs.ErrNotExist
		},
		// a sample: the first 10 pages only
		Allow: func(r *http.Request, s bdf.Segment) error {
			if s.From >= 10 && r.URL.Query().Get("owner") == "" {
				return fs.ErrPermission
			}
			return nil
		},
	}
	srv := httptest.NewServer(h)
	defer srv.Close()
	priv, _ := ecdh.P256().GenerateKey(rand.Reader)
	key := priv.PublicKey().Bytes()

	res := post(t, srv.URL+"/book?owner=1", Request{Key: key, Page: 13, Have: []bdf.Segment{{View: "pages", From: 0, To: 10}}})
	if res.StatusCode != http.StatusOK || res.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("status %d, Cache-Control %q", res.StatusCode, res.Header.Get("Cache-Control"))
	}
	b, _ := io.ReadAll(res.Body)
	seg, err := bdf.OpenSingle(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	if err := seg.UnlockECDH(priv); err != nil {
		t.Fatal(err)
	}
	if s := seg.Manifest.Segment; s == nil || *s != (bdf.Segment{View: "pages", From: 10, To: 20}) {
		t.Fatalf("segment %+v", s)
	}

	for _, c := range []struct {
		path   string
		body   any
		status int
	}{
		{"/book", Request{Key: key, Page: 13}, http.StatusForbidden}, // past the sample
		{"/book", Request{Key: key, Page: 25}, http.StatusBadRequest},
		{"/book", Request{Key: key, View: "none"}, http.StatusBadRequest},
		{"/book", Request{Key: key, Have: []bdf.Segment{{View: "pages", From: 20, To: 30}}}, http.StatusBadRequest},
		{"/book", Request{Key: key[:33]}, http.StatusBadRequest},
		{"/locked", Request{Key: key}, http.StatusUnauthorized},
		{"/theirs", Request{Key: key}, http.StatusForbidden},
		{"/nothing", Request{Key: key}, http.StatusNotFound},
	} {
		if res := post(t, srv.URL+c.path, c.body); res.StatusCode != c.status {
			t.Errorf("%s %+v: status %d, want %d", c.path, c.body, res.StatusCode, c.status)
		}
	}
	if res, _ := http.Get(srv.URL + "/book"); res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET: status %d", res.StatusCode)
	}
	// a form, which another site can send without asking
	res, _ = http.Post(srv.URL+"/book", "text/plain", bytes.NewReader([]byte("{}")))
	if res.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("text/plain: status %d", res.StatusCode)
	}
}

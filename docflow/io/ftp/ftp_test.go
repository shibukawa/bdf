package ftp

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/flowtest"
)

func settings(t *testing.T, v map[string]any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// An exit uploads under the path its settings expand to, over TLS
// (AUTH TLS, the certificate pinned, the data connections resuming the
// session), through a temporary name, once.
func TestFTPOut(t *testing.T) {
	srv := newFTPServer(t, ftpOpts{TLS: true, MLSD: true})
	t.Setenv(docflow.CredentialEnv("ftp-pw"), testPassword)
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "out", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://tester@" + srv.Addr + "/", "credential": "ftp-pw", "tls_fingerprint": srv.Fingerprint, "path": "orders/{{meta.customer}}/{{meta.$name}}"})},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"out"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "a.bdf", doc, map[string]any{"customer": "../../acme"})
	e.Tick()
	d := e.Delivery(j, "out")
	if d.State != docflow.DeliveryDone || d.Result != "ftp://"+srv.Addr+"/orders/acme/a.bdf" {
		t.Fatalf("%s %s: %s", d.State, d.Result, d.LastError)
	}
	if got, err := os.ReadFile(filepath.Join(srv.root, "orders", "acme", "a.bdf")); err != nil || !bytes.Equal(got, doc) {
		t.Fatalf("the file uploaded: %d bytes, %v", len(got), err)
	}
	if names := srv.names(t); !slices.Equal(names, []string{"orders/acme/a.bdf"}) {
		t.Fatalf("the server has %v", names)
	}
	srv.mu.Lock()
	dataTLS, resumed := srv.dataTLS, srv.dataResumed
	srv.mu.Unlock()
	if dataTLS == 0 || resumed == 0 {
		t.Errorf("data connections over TLS: %d, resuming the session: %d", dataTLS, resumed)
	}
	// a file there already is left
	srv.writeFile(t, "orders/acme/b.bdf", []byte("before"), time.Time{})
	j = e.Submit("in", "b.bdf", doc, map[string]any{"customer": "acme"})
	e.Tick()
	if d := e.Delivery(j, "out"); d.State != docflow.DeliveryDone {
		t.Fatalf("a file there: %s: %s", d.State, d.LastError)
	}
	if got, _ := os.ReadFile(filepath.Join(srv.root, "orders", "acme", "b.bdf")); string(got) != "before" {
		t.Fatalf("wrote over the file there: %d bytes", len(got))
	}
}

// The ways of TLS: none to a server without it, TLS required to a
// server without it fails for good, TLS from the start (ftps://), a
// wrong fingerprint and a wrong password fail for good, anonymous signs
// in with no user.
func TestFTPTLS(t *testing.T) {
	plain := newFTPServer(t, ftpOpts{MLSD: true})
	implicit := newFTPServer(t, ftpOpts{Implicit: true, MLSD: true})
	secure := newFTPServer(t, ftpOpts{TLS: true, MLSD: true})
	t.Setenv(docflow.CredentialEnv("ftp-pw"), testPassword)
	t.Setenv(docflow.CredentialEnv("ftp-bad"), "nope")
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "clear", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://tester@" + plain.Addr + "/", "credential": "ftp-pw", "tls": "none"})},
		docflow.Node{Name: "anon", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://" + plain.Addr + "/anon", "tls": "none"})},
		docflow.Node{Name: "notls", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://tester@" + plain.Addr + "/", "credential": "ftp-pw"})},
		docflow.Node{Name: "ftps", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftps://tester@" + implicit.Addr + "/", "credential": "ftp-pw", "tls_fingerprint": strings.ToUpper(implicit.Fingerprint)})},
		docflow.Node{Name: "wrongfp", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://tester@" + secure.Addr + "/", "credential": "ftp-pw", "tls_fingerprint": implicit.Fingerprint})},
		docflow.Node{Name: "wrongpw", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://tester@" + secure.Addr + "/", "credential": "ftp-bad", "tls_fingerprint": secure.Fingerprint})},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"clear", "anon", "notls", "ftps", "wrongfp", "wrongpw"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "a.bdf", doc, nil)
	e.Tick()
	for _, name := range []string{"clear", "anon", "ftps"} {
		if d := e.Delivery(j, name); d.State != docflow.DeliveryDone {
			t.Errorf("%s: %s: %s", name, d.State, d.LastError)
		}
	}
	if names := plain.names(t); !slices.Equal(names, []string{"1/a.bdf", "anon/1/a.bdf"}) {
		t.Errorf("the plain server has %v", names)
	}
	if names := implicit.names(t); !slices.Equal(names, []string{"1/a.bdf"}) {
		t.Errorf("the implicit server has %v", names)
	}
	for name, want := range map[string]string{"notls": "offers no TLS", "wrongfp": "certificate", "wrongpw": "530"} {
		if d := e.Delivery(j, name); d.State != docflow.DeliveryFailed || !strings.Contains(d.LastError, want) {
			t.Errorf("%s: %s: %s", name, d.State, d.LastError)
		}
	}
}

// A polling entry over a server without MLSD (NLST, SIZE and MDTM)
// takes the files that fit and moves them to the done directory.
func TestFTPPoll(t *testing.T) {
	srv := newFTPServer(t, ftpOpts{TLS: true})
	t.Setenv(docflow.CredentialEnv("ftp-pw"), testPassword)
	e := flowtest.New(t, nil)
	old := e.Now().Add(-time.Hour)
	doc := flowtest.Doc(t, 1)
	srv.writeFile(t, "in/a.bdf", doc, old)
	srv.writeFile(t, "in/b.bdf", doc, old)
	srv.writeFile(t, "in/c.txt", []byte("text"), old)
	srv.writeFile(t, "in/sub/d.bdf", doc, old)
	e.Put(docflow.Node{Name: "seen", Kind: docflow.KindSink},
		docflow.Node{Name: "poll", Kind: KindPoll, To: []string{"seen"}, Settings: settings(t, map[string]any{
			"url": "ftp://tester@" + srv.Addr + "/in", "credential": "ftp-pw", "tls_fingerprint": srv.Fingerprint,
			"pattern": "*.bdf", "after": "move", "done_dir": "/done"})})
	rep, err := e.M().Tick(context.Background())
	if err != nil || rep.Polled != 2 {
		t.Fatalf("tick: %+v, %v", rep, err)
	}
	j := e.Job(1)
	if j.Name != "a.bdf" || j.Source != KindPoll || j.Meta["ftp.path"] != "/in/a.bdf" || j.Meta["ftp.time"] != old.UTC().Format(time.RFC3339) {
		t.Fatalf("job 1: %+v", j)
	}
	if names := srv.names(t); !slices.Equal(names, []string{"done/a.bdf", "done/b.bdf", "in/c.txt", "in/sub/d.bdf"}) {
		t.Fatalf("the server has %v", names)
	}
	// a second look, a minute later: nothing left that fits
	e.Advance(time.Minute)
	if rep, _ := e.M().Tick(context.Background()); rep.Polled != 0 {
		t.Fatalf("tick 2: %+v", rep)
	}
}

// The settings: the schemes and their ports, the TLS, the fingerprint,
// no password in the URL.
func TestFTPSettings(t *testing.T) {
	e := flowtest.New(t, nil)
	m := e.M()
	ctx := context.Background()
	for _, c := range []struct{ url, tls, fp, want string }{
		{"ftp://u:pw@h/", "", "", "password"},
		{"sftp://u@h/", "", "", "ftp|ftps"},
		{"ftp://u@h/", "maybe", "", "tls"},
		{"ftp://u@h/", "", "abc", "fingerprint"},
	} {
		s := map[string]any{"url": c.url}
		if c.tls != "" {
			s["tls"] = c.tls
		}
		if c.fp != "" {
			s["tls_fingerprint"] = c.fp
		}
		if _, err := m.PutNode(ctx, docflow.Node{Name: "x", Kind: KindFTP, Settings: settings(t, s)}); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%+v: %v", c, err)
		}
	}
	e.Required(docflow.Node{Name: "x", Kind: KindFTP, Settings: settings(t, map[string]any{"url": "ftp://u@h/"})})
	e.Required(docflow.Node{Name: "y", Kind: KindPoll, Settings: settings(t, map[string]any{"url": "ftp://u@h/", "after": "move", "done_dir": "done"})})
}

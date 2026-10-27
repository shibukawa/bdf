package ftp

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/flowtest"
	"github.com/shibukawa/bdf/docflow/io/transfer"
)

// TestFTPReal sends to, and polls, a real FTP server the environment
// names, with a directory in and a directory done under the URL's:
// BDF_FTP_REAL_URL (ftp://tester@127.0.0.1:2121/), BDF_FTP_REAL_PASSWORD,
// BDF_FTP_REAL_FINGERPRINT (the certificate's SHA-256; empty: the
// roots) and BDF_FTP_REAL_TLS (explicit, the default, or none). It is
// skipped without them. A vsftpd with require_ssl_reuse is what it was
// written against: testdata/vsftpd has one, and how to run it.
func TestFTPReal(t *testing.T) {
	url := os.Getenv("BDF_FTP_REAL_URL")
	if url == "" {
		t.Skip("BDF_FTP_REAL_URL is not set")
	}
	t.Setenv(docflow.CredentialEnv("ftp-real"), os.Getenv("BDF_FTP_REAL_PASSWORD"))
	e := flowtest.New(t, nil)
	common := map[string]any{"credential": "ftp-real", "tls_fingerprint": os.Getenv("BDF_FTP_REAL_FINGERPRINT")}
	if v := os.Getenv("BDF_FTP_REAL_TLS"); v != "" {
		common["tls"] = v
	}
	with := func(kv ...any) json.RawMessage {
		m := map[string]any{}
		for k, v := range common {
			m[k] = v
		}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		b, _ := json.Marshal(m)
		return b
	}
	name := "docflow-real-" + time.Now().UTC().Format("150405") + ".bdf"
	e.Put(docflow.Node{Name: "out", Kind: KindFTP, Settings: with("url", url, "path", "in/{{meta.$name}}")},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"out"}},
		docflow.Node{Name: "seen", Kind: docflow.KindSink},
		docflow.Node{Name: "poll", Kind: KindPoll, To: []string{"seen"}, Settings: with("url", url+"in", "pattern", name, "after", "move", "done_dir", "done")})
	doc := flowtest.Doc(t, 2)
	j := e.Submit("in", name, doc, nil)
	e.Tick()
	if d := e.Delivery(j, "out"); d.State != docflow.DeliveryDone {
		t.Fatalf("out: %s: %s", d.State, d.LastError)
	} else {
		t.Logf("sent: %s", d.Result)
	}
	// the first look leaves the file (the server's clock is not the
	// test's: the file is not settled), the second, a minute later,
	// takes it, seen the same
	polled := 0
	for range 2 {
		rep, err := e.M().Tick(context.Background())
		if err != nil {
			t.Fatalf("poll: %+v, %v", rep, err)
		}
		polled += rep.Polled
		e.Advance(time.Minute)
	}
	if polled != 1 {
		t.Fatalf("polled %d files", polled)
	}
	if pj := e.Job(2); pj.Name != name || pj.Size != int64(len(doc)) {
		t.Fatalf("the job polled: %+v", pj)
	}
	// what the server has now: nothing in in, the file in done (removed after)
	s, mode, err := site(url, os.Getenv("BDF_FTP_REAL_TLS"), os.Getenv("BDF_FTP_REAL_FINGERPRINT"))
	if err != nil {
		t.Fatal(err)
	}
	c, err := Dial(context.Background(), Config{Addr: s.Address(), User: s.User, Password: os.Getenv("BDF_FTP_REAL_PASSWORD"), TLS: mode, Fingerprint: os.Getenv("BDF_FTP_REAL_FINGERPRINT")})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	t.Logf("features: %v", c.Features())
	names := func(dir string) []string {
		es, err := c.List(context.Background(), s.Path(dir))
		if err != nil {
			t.Fatal(err)
		}
		var out []string
		for _, e := range es {
			out = append(out, e.Name)
		}
		return out
	}
	if slices.Contains(names("in"), name) {
		t.Errorf("in still has %s", name)
	}
	if !slices.Contains(names("in/done"), name) {
		t.Errorf("done has not %s: %v", name, names("in/done"))
	}
	if err := c.Remove(context.Background(), s.Path("in/done/"+name)); err != nil {
		t.Error(err)
	}
	var _ transfer.FS = c
}

package sftp

import (
	"context"
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/flowtest"
)

// TestSFTPReal sends to, and polls, a real SFTP server the environment
// names, with a directory in and a directory done under the URL's:
// BDF_SFTP_REAL_URL (sftp://tester@127.0.0.1:2222/~/) and
// BDF_SFTP_REAL_PASSWORD (or BDF_SFTP_REAL_KEY, a PEM file). The host
// key is fetched from the server first, as the editor's Discover does.
// It is skipped without them. OpenSSH's sshd with internal-sftp is what
// it was written against: testdata/sshd has one, and how to run it.
func TestSFTPReal(t *testing.T) {
	url := os.Getenv("BDF_SFTP_REAL_URL")
	if url == "" {
		t.Skip("BDF_SFTP_REAL_URL is not set")
	}
	dir := t.TempDir()
	if key := os.Getenv("BDF_SFTP_REAL_KEY"); key != "" {
		b, err := os.ReadFile(key)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dir+"/sftp-real", b, 0o600); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Setenv(docflow.CredentialEnv("sftp-real"), os.Getenv("BDF_SFTP_REAL_PASSWORD"))
	}
	e := flowtest.New(t, func(o *docflow.Options) { o.CredentialDir = dir })
	s, err := docflow.SettingsOf[Settings](&docflow.Node{Kind: KindSFTP, Settings: json.RawMessage(`{"url": ` + mustJSON(url) + `}`)})
	if err != nil {
		t.Fatal(err)
	}
	st, err := site(s.URL, "SHA256:placeholder-for-the-parse")
	if err != nil {
		t.Fatal(err)
	}
	info, err := HostKey(context.Background(), st.Address())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("host key: %s %s", info.Type, info.Fingerprint)
	with := func(kv ...any) json.RawMessage {
		m := map[string]any{"credential": "sftp-real", "host_key": info.HostKey}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
		b, _ := json.Marshal(m)
		return b
	}
	name := "docflow-real-" + time.Now().UTC().Format("150405") + ".bdf"
	e.Put(docflow.Node{Name: "out", Kind: KindSFTP, Settings: with("url", url, "path", "in/{{meta.$name}}")},
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
	check, _ := ParseHostKey(info.HostKey)
	c, err := Dial(context.Background(), Config{Addr: st.Address(), User: st.User, Password: os.Getenv("BDF_SFTP_REAL_PASSWORD"), HostKey: check})
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	names := func(dir string) []string {
		es, err := c.List(context.Background(), st.Path(dir))
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
	if err := c.Remove(context.Background(), st.Path("in/done/"+name)); err != nil {
		t.Error(err)
	}
}

func mustJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

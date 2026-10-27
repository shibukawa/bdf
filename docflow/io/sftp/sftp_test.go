package sftp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/flowtest"
	"golang.org/x/crypto/ssh"
)

func settings(t *testing.T, v map[string]any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func hostKeyLine(s *sshServer) string {
	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(s.HostKey)))
}

// An exit uploads under the path its settings expand to, through a
// temporary name, once: a value cannot leave the base, a file there
// already (a write whose outcome was not known) is taken as written,
// and an overwriting node writes over it. The result is the file's URL.
func TestSFTPOut(t *testing.T) {
	srv := newSSHServer(t, true)
	t.Setenv(docflow.CredentialEnv("sftp-pw"), testPassword)
	e := flowtest.New(t, nil)
	base := map[string]any{"url": "sftp://tester@" + srv.Addr + "/", "credential": "sftp-pw", "host_key": hostKeyLine(srv)}
	with := func(kv ...any) map[string]any {
		out := map[string]any{}
		for k, v := range base {
			out[k] = v
		}
		for i := 0; i < len(kv); i += 2 {
			out[kv[i].(string)] = kv[i+1]
		}
		return out
	}
	e.Put(docflow.Node{Name: "out", Kind: KindSFTP, Settings: settings(t, with("path", "orders/{{meta.customer}}/{{meta.$name}}"))},
		docflow.Node{Name: "over", Kind: KindSFTP, Settings: settings(t, with("path", "orders/{{meta.customer}}/{{meta.$name}}", "overwrite", true))},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"out"}},
		docflow.Node{Name: "in2", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"over"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "a.bdf", doc, map[string]any{"customer": "../../acme"})
	e.Tick()
	d := e.Delivery(j, "out")
	if d.State != docflow.DeliveryDone || d.Result != "sftp://"+srv.Addr+"/orders/acme/a.bdf" {
		t.Fatalf("%s %s: %s", d.State, d.Result, d.LastError)
	}
	if got, err := os.ReadFile(filepath.Join(srv.root, "orders", "acme", "a.bdf")); err != nil || !bytes.Equal(got, doc) {
		t.Fatalf("the file uploaded: %d bytes, %v", len(got), err)
	}
	if names := srv.names(t); !slices.Equal(names, []string{"orders/acme/a.bdf"}) {
		t.Fatalf("the server has %v", names)
	}
	// a file there already: left by out, written over by over
	srv.writeFile(t, "orders/acme/b.bdf", []byte("before"), time.Time{})
	j = e.Submit("in", "b.bdf", doc, map[string]any{"customer": "acme"})
	e.Tick()
	if d := e.Delivery(j, "out"); d.State != docflow.DeliveryDone {
		t.Fatalf("out, a file there: %s: %s", d.State, d.LastError)
	}
	if got, _ := os.ReadFile(filepath.Join(srv.root, "orders", "acme", "b.bdf")); string(got) != "before" {
		t.Fatalf("out wrote over the file there: %d bytes", len(got))
	}
	// the overwriting node writes over it
	j = e.Submit("in2", "b.bdf", doc, map[string]any{"customer": "acme"})
	e.Tick()
	if d := e.Delivery(j, "over"); d.State != docflow.DeliveryDone {
		t.Fatalf("over: %s: %s", d.State, d.LastError)
	}
	if got, _ := os.ReadFile(filepath.Join(srv.root, "orders", "acme", "b.bdf")); !bytes.Equal(got, doc) {
		t.Fatalf("over did not write over the file there: %d bytes", len(got))
	}
	for _, n := range srv.names(t) {
		if strings.Contains(n, ".part") {
			t.Errorf("a temporary file left: %s", n)
		}
	}
}

// A private key signs in as a password does; a host key that is not
// the server's fails for good; a credential not there is tried again.
func TestSFTPAuth(t *testing.T) {
	srv := newSSHServer(t, true)
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sftp-key"), srv.KeyPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(docflow.CredentialEnv("sftp-pw"), testPassword)
	e := flowtest.New(t, func(o *docflow.Options) { o.CredentialDir = dir })
	other := newSSHServer(t, true) // another host key
	e.Put(docflow.Node{Name: "key", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/~/", "credential": "sftp-key", "host_key": hostKeyLine(srv), "path": "k/{{meta.$name}}"})},
		docflow.Node{Name: "fp", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/~/", "credential": "sftp-pw", "host_key": ssh.FingerprintSHA256(srv.HostKey), "path": "f/{{meta.$name}}"})},
		docflow.Node{Name: "wrong", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/", "credential": "sftp-pw", "host_key": hostKeyLine(other)})},
		docflow.Node{Name: "nocred", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/", "credential": "sftp-none", "host_key": hostKeyLine(srv)})},
		docflow.Node{Name: "in", Kind: docflow.KindUpload, Output: docflow.OutputFile, To: []string{"key", "fp", "wrong", "nocred"}})
	doc := flowtest.Doc(t, 1)
	j := e.Submit("in", "a.bdf", doc, nil)
	e.Tick()
	for _, name := range []string{"key", "fp"} {
		if d := e.Delivery(j, name); d.State != docflow.DeliveryDone {
			t.Errorf("%s: %s: %s", name, d.State, d.LastError)
		}
	}
	if names := srv.names(t); !slices.Equal(names, []string{"f/a.bdf", "k/a.bdf"}) {
		t.Errorf("the server has %v", names)
	}
	if d := e.Delivery(j, "wrong"); d.State != docflow.DeliveryFailed || !strings.Contains(d.LastError, "host key") {
		t.Errorf("a wrong host key: %s: %s", d.State, d.LastError)
	}
	if d := e.Delivery(j, "nocred"); d.State != docflow.DeliveryPending || !strings.Contains(d.LastError, "sftp-none") || d.Attempts != 1 {
		t.Errorf("a credential not there: %s (%d): %s", d.State, d.Attempts, d.LastError)
	}
}

// Discover fetches the server's host key, for the editor; a node needs
// a host key, and a password in the URL is refused.
func TestSFTPDiscover(t *testing.T) {
	srv := newSSHServer(t, true)
	e := flowtest.New(t, nil)
	body, _ := json.Marshal(map[string]any{"settings": map[string]any{"url": "sftp://tester@" + srv.Addr + "/"}})
	req, _ := http.NewRequest(http.MethodPost, e.Srv.URL+"/api/kinds/sftp/discover", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+flowtest.AdminToken)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var info HostKeyInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("discover: %d, %v", res.StatusCode, err)
	}
	if info.HostKey != hostKeyLine(srv) || info.Fingerprint != ssh.FingerprintSHA256(srv.HostKey) || info.Type != "ssh-ed25519" {
		t.Fatalf("discover: %+v", info)
	}
	m := e.M()
	ctx := context.Background()
	if _, err := m.PutNode(ctx, docflow.Node{Name: "x", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/", "credential": "c"})}); err == nil || !strings.Contains(err.Error(), "host_key") {
		t.Errorf("without a host key: %v", err)
	}
	if _, err := m.PutNode(ctx, docflow.Node{Name: "x", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester:pw@" + srv.Addr + "/", "credential": "c", "host_key": info.HostKey})}); err == nil || !strings.Contains(err.Error(), "password") {
		t.Errorf("with a password in the URL: %v", err)
	}
	if _, err := m.PutNode(ctx, docflow.Node{Name: "x", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://" + srv.Addr + "/", "credential": "c", "host_key": info.HostKey})}); err == nil || !strings.Contains(err.Error(), "user") {
		t.Errorf("without a user: %v", err)
	}
	if _, err := m.PutNode(ctx, docflow.Node{Name: "x", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/", "credential": "c", "host_key": "not a key"})}); err == nil || !strings.Contains(err.Error(), "host key") {
		t.Errorf("with a bad host key: %v", err)
	}
	e.Required(docflow.Node{Name: "x", Kind: KindSFTP, Settings: settings(t, map[string]any{"url": "sftp://tester@" + srv.Addr + "/", "credential": "c", "host_key": info.HostKey})})
}

// A polling entry takes the files that fit, in order of name, as many
// as its settings say per tick and no more often than its every,
// leaving those still being written for the next look; what it takes
// becomes a job with the file's path, size and time, and is deleted or
// moved; what the engine refuses is left renamed, and so is a file too
// large, with a note.
func TestSFTPPoll(t *testing.T) {
	srv := newSSHServer(t, true)
	t.Setenv(docflow.CredentialEnv("sftp-pw"), testPassword)
	e := flowtest.New(t, func(o *docflow.Options) { o.MaxDocument = 4000 })
	old := e.Now().Add(-time.Hour)
	doc := flowtest.Doc(t, 1)
	srv.writeFile(t, "in/b.bdf", doc, old)
	srv.writeFile(t, "in/a.bdf", doc, old)
	srv.writeFile(t, "in/c.txt", []byte("text"), old)
	srv.writeFile(t, "in/fresh.bdf", doc, e.Now())
	srv.writeFile(t, "in/.hidden.bdf", doc, old)
	srv.writeFile(t, "in/d.bdf.part", doc, old)
	srv.writeFile(t, "in/z.bdf.docflow-other.taking", doc, old)
	srv.writeFile(t, "in/old.bdf.docflow-poll.taking", doc, old)
	srv.writeFile(t, "in/bad.bdf", []byte("not a document at all"), old)
	srv.writeFile(t, "in/big.bdf", bytes.Repeat([]byte("x"), 5000), old)
	e.Put(docflow.Node{Name: "seen", Kind: docflow.KindSink},
		docflow.Node{Name: "poll", Kind: KindPoll, To: []string{"seen"}, Settings: settings(t, map[string]any{
			"url": "sftp://tester@" + srv.Addr + "/in", "credential": "sftp-pw", "host_key": hostKeyLine(srv),
			"pattern": "*.bdf", "every": "1m", "settle": "10s", "max_per_tick": 2})})
	ctx := context.Background()
	tick := func() *docflow.TickReport {
		t.Helper()
		rep, _ := e.M().Tick(ctx)
		return rep
	}
	// 1: the claim left by an earlier tick first, then a.bdf; b.bdf
	// waits (two per tick); bad.bdf and big.bdf are noted
	rep := tick()
	if rep.Ticked != 1 || rep.Polled != 2 {
		t.Fatalf("tick 1: %+v", rep)
	}
	j := e.Job(1)
	if j.Name != "old.bdf" || j.Source != KindPoll || j.ClaimedUser != "tester" || j.Meta["sftp.path"] != "/in/old.bdf" || j.Meta["sftp.size"] != float64(len(doc)) || j.Meta["sftp.time"] != old.UTC().Format(time.RFC3339) {
		t.Fatalf("job 1: %+v", j)
	}
	if j := e.Job(2); j.Name != "a.bdf" {
		t.Fatalf("job 2: %s", j.Name)
	}
	want := []string{"in/.hidden.bdf", "in/b.bdf", "in/bad.bdf", "in/big.bdf", "in/c.txt", "in/d.bdf.part", "in/fresh.bdf", "in/z.bdf.docflow-other.taking"}
	if names := srv.names(t); !slices.Equal(names, want) {
		t.Fatalf("after tick 1 the server has %v", names)
	}
	if len(rep.Notes) != 1 || !strings.Contains(rep.Notes[0], "big.bdf") {
		t.Fatalf("tick 1 notes: %v", rep.Notes)
	}
	// 2: at once again: not due
	if rep := tick(); rep.Polled != 0 || rep.Ticked != 1 {
		t.Fatalf("tick 2 (not due): %+v", rep)
	}
	// 3: a minute later: b.bdf, and bad.bdf refused; fresh.bdf, seen
	// the same as before, is settled now and taken too? no: two per
	// tick, and bad.bdf took a turn
	e.Advance(time.Minute)
	rep = tick()
	if rep.Polled != 1 || len(rep.Notes) != 2 {
		t.Fatalf("tick 3: %+v", rep)
	}
	if j := e.Job(3); j.Name != "b.bdf" {
		t.Fatalf("job 3: %s", j.Name)
	}
	names := srv.names(t)
	if !slices.Contains(names, "in/bad.bdf.docflow-poll.refused") || slices.Contains(names, "in/bad.bdf") {
		t.Fatalf("after tick 3 the server has %v", names)
	}
	// 4: another minute: fresh.bdf (unchanged since the last look)
	e.Advance(time.Minute)
	if rep := tick(); rep.Polled != 1 {
		t.Fatalf("tick 4: %+v", rep)
	}
	if j := e.Job(4); j.Name != "fresh.bdf" {
		t.Fatalf("job 4: %s", j.Name)
	}
	// 5: nothing left that fits
	e.Advance(time.Minute)
	if rep := tick(); rep.Polled != 0 {
		t.Fatalf("tick 5: %+v", rep)
	}
	want = []string{"in/.hidden.bdf", "in/bad.bdf.docflow-poll.refused", "in/big.bdf", "in/c.txt", "in/d.bdf.part", "in/z.bdf.docflow-other.taking"}
	if names := srv.names(t); !slices.Equal(names, want) {
		t.Fatalf("at the end the server has %v", names)
	}
}

// After taking, a file is moved to the done directory, made if it is
// not there, under the base.
func TestSFTPPollMove(t *testing.T) {
	srv := newSSHServer(t, false) // no posix-rename: rename over removes first
	t.Setenv(docflow.CredentialEnv("sftp-pw"), testPassword)
	e := flowtest.New(t, nil)
	doc := flowtest.Doc(t, 1)
	srv.writeFile(t, "in/a.bdf", doc, e.Now().Add(-time.Hour))
	srv.writeFile(t, "in/done/a.bdf", []byte("earlier"), e.Now().Add(-time.Hour))
	e.Put(docflow.Node{Name: "seen", Kind: docflow.KindSink},
		docflow.Node{Name: "poll", Kind: KindPoll, To: []string{"seen"}, Settings: settings(t, map[string]any{
			"url": "sftp://tester@" + srv.Addr + "/in", "credential": "sftp-pw", "host_key": hostKeyLine(srv), "after": "move", "done_dir": "done"})})
	rep, err := e.M().Tick(context.Background())
	if err != nil || rep.Polled != 1 {
		t.Fatalf("tick: %+v, %v", rep, err)
	}
	if names := srv.names(t); !slices.Equal(names, []string{"in/done/a.bdf"}) {
		t.Fatalf("the server has %v", names)
	}
	if got, _ := os.ReadFile(filepath.Join(srv.root, "in", "done", "a.bdf")); !bytes.Equal(got, doc) {
		t.Fatalf("the file moved: %d bytes", len(got))
	}
	// the settings of a move need a done directory
	if _, err := e.M().PutNode(context.Background(), docflow.Node{Name: "poll2", Kind: KindPoll, To: []string{"seen"}, Settings: settings(t, map[string]any{
		"url": "sftp://tester@" + srv.Addr + "/in", "credential": "sftp-pw", "host_key": hostKeyLine(srv), "after": "move"})}); err == nil || !strings.Contains(err.Error(), "done_dir") {
		t.Fatalf("move without a done directory: %v", err)
	}
}

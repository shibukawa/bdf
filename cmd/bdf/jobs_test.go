package main

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
	"github.com/shibukawa/bdf/print/server"
	"github.com/shibukawa/bdf/print/server/client"
)

// printServer serves a print server of a fake printer, with the
// converters of bdf printd, and returns its URL.
func printServer(t *testing.T, extra string) (string, *server.Server, *ipptest.Printer) {
	t.Helper()
	p, u := fakePrinter(t)
	cfg, err := server.ParseConfig([]byte(fmt.Sprintf("printers: { mfp1: { url: %q } }\nqueues: { office: { printers: [mfp1] }, spare: { printers: [mfp1] } }\n%s", u, extra)))
	if err != nil {
		t.Fatal(err)
	}
	cfg.Spool = filepath.Join(t.TempDir(), "spool.db")
	cfg.TempDir = t.TempDir()
	fonts := filepath.Join("..", "..", "converter", "pptx", "testdata", "fonts")
	srv, err := server.New(context.Background(), cfg, server.Options{Converter: printConverter{fontDirs: []string{fonts}, noSystem: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { srv.Close() })
	hs := httptest.NewServer(srv.Handler())
	t.Cleanup(hs.Close)
	srv.Spool().Start()
	return hs.URL, srv, p
}

func TestPrintToServer(t *testing.T) {
	url, srv, p := printServer(t, "")
	doc := filepath.Join("..", "..", "testdata", "form", "rows.bdf")
	t.Setenv(serverEnv, "")

	// queued, then followed until printed
	code, out, errOut := runOut(t, "print", "-server", url, "-queue", "office", "-sides", "long", "-wait", "30s", doc)
	if code != 0 || strings.TrimSpace(out) != "1" || !strings.Contains(errOut, "job 1 completed on mfp1") {
		t.Fatalf("exit %d: %q %s", code, out, errOut)
	}
	pj := p.Jobs()[0]
	if pj.Name != "#1.1 rows.bdf" || !strings.HasPrefix(string(pj.Document(0)), "%PDF-") {
		t.Fatalf("at the printer: %q", pj.Name)
	}

	// an input the server converts; $BDF_SERVER names the server
	t.Setenv(serverEnv, url)
	pdf := filepath.Join("..", "..", "converter", "pdf", "testdata", "cairo-cff.pdf")
	if code, out, errOut := runOut(t, "print", "-queue", "office", "-wait", "60s", "-pages", "1", pdf); code != 0 || strings.TrimSpace(out) != "2" {
		t.Fatalf("pdf: exit %d: %q %s", code, out, errOut)
	}
	// Markdown is told by its name, and printed as pages
	md := filepath.Join(t.TempDir(), "notes.md")
	os.WriteFile(md, []byte("# Notes\n\nA line.\n"), 0o644)
	if code, out, errOut := runOut(t, "print", "-queue", "office", "-wait", "60s", md); code != 0 || strings.TrimSpace(out) != "3" {
		t.Fatalf("markdown: exit %d: %q %s", code, out, errOut)
	}

	// an encrypted document with its password (the server is on this
	// machine: no TLS is needed)
	r, err := bdf.OpenSingleFile(doc)
	if err != nil {
		t.Fatal(err)
	}
	d, err := r.ToDocument()
	r.Close()
	if err != nil {
		t.Fatal(err)
	}
	if d.Lock, err = bdf.NewPasswordLock("open sesame", 1000); err != nil {
		t.Fatal(err)
	}
	enc := filepath.Join(t.TempDir(), "secret.bdf")
	f, _ := os.Create(enc)
	d.WriteSingle(f)
	f.Close()
	t.Setenv(passwordEnv, "open sesame")
	if code, _, errOut := runOut(t, "print", "-queue", "office", "-wait", "30s", enc); code != 0 {
		t.Fatalf("encrypted: exit %d: %s", code, errOut)
	}
	t.Setenv(passwordEnv, "")

	// a queue that does not send: the job stays, and is changed
	c := &client.Client{URL: url}
	ctx := context.Background()
	c.QueueOp(ctx, "spare", "pause")
	if code, out, errOut := runOut(t, "print", "-queue", "spare", "-wait", "1s", doc); code != exitTimeout || strings.TrimSpace(out) != "5" {
		t.Fatalf("not done in time: exit %d: %q %s", code, out, errOut)
	}
	if code, out, _ := runOut(t, "print", "-queue", "spare", "-priority", "20", "-hold", "indefinite", "-q", doc); code != 0 || strings.TrimSpace(out) != "6" {
		t.Fatalf("held: exit %d: %q", code, out)
	}
	j, _ := c.Job(ctx, 6)
	if j.State != "pending-held" || j.Priority != 20 || strings.Contains(string(j.Ticket), "priority") || strings.Contains(string(j.Ticket), "hold") {
		t.Fatalf("held job %+v %s", j, j.Ticket)
	}

	for _, step := range []struct {
		args []string
		want string
	}{
		{[]string{"jobs"}, "spare"},
		{[]string{"jobs", "-all"}, "completed"},
		{[]string{"jobs", "show", "1"}, "printer:    mfp1 (job 1)"},
		{[]string{"jobs", "priority", "5", "90"}, "job 5 pending"},
		{[]string{"jobs", "release", "6"}, "job 6 pending"},
		{[]string{"jobs", "hold", "6"}, "job 6 pending-held"},
		{[]string{"jobs", "move", "6", "office"}, "job 6 pending-held"},
		{[]string{"jobs", "cancel", "5"}, "job 5 canceled"},
		{[]string{"jobs", "retry", "5"}, "job 5 pending"},
		{[]string{"jobs", "accounting", "-csv"}, "id,uuid,queue"},
		{[]string{"printers", "-server", url}, "mfp1"},
	} {
		code, out, errOut := runOut(t, step.args...)
		if code != 0 || !strings.Contains(out, step.want) {
			t.Errorf("%v: exit %d: %s %s", step.args, code, out, errOut)
		}
	}
	if j, _ := c.Job(ctx, 6); j.Queue != "office" {
		t.Errorf("moved %+v", j)
	}

	// what is refused
	c.QueueOp(ctx, "spare", "reject")
	for _, bad := range []struct {
		args []string
		code int
	}{
		{[]string{"print", "-queue", "spare", doc}, exitStopped},
		{[]string{"print", "-queue", "nowhere", doc}, exitUnsupported},
		{[]string{"print", doc}, 2}, // two queues
		{[]string{"print", "-server", url, "-printer", "ipp://h/", doc}, 2},
		{[]string{"print", "-queue", "office", "-param", "a=b", doc}, 2},
		{[]string{"print", "-server", "http://127.0.0.1:1", "-queue", "office", doc}, exitConnection},
		{[]string{"jobs", "show", "99"}, exitUnsupported},
	} {
		if code, _, errOut := runOut(t, bad.args...); code != bad.code {
			t.Errorf("%v: exit %d, not %d: %s", bad.args, code, bad.code, errOut)
		}
	}

	// an encrypted deck waits for its password, then prints
	deck := filepath.Join("..", "..", "converter", "internal", "offcrypto", "testdata", "agile.pptx")
	code, out, errOut = runOut(t, "print", "-queue", "office", "-q", deck)
	if code != 0 || strings.TrimSpace(out) != "7" {
		t.Fatalf("deck: exit %d: %q %s", code, out, errOut)
	}
	waitState := func(want string) *client.Job {
		t.Helper()
		var j *client.Job
		for i := 0; i < 300; i++ {
			if j, _ = c.Job(ctx, 7); j != nil && j.State == want {
				return j
			}
			time.Sleep(50 * time.Millisecond)
		}
		t.Fatalf("job 7 is %+v, not %s", j, want)
		return nil
	}
	if j := waitState("pending-held"); j.Reasons[0] != "job-password-wait" {
		t.Fatalf("deck without its password: %+v", j)
	}
	t.Setenv(passwordEnv, "パスワード🔑bdf")
	if code, out, errOut := runOut(t, "jobs", "password", "7"); code != 0 || !strings.Contains(out, "job 7 pending") {
		t.Fatalf("password: exit %d: %s %s", code, out, errOut)
	}
	waitState("completed")
	_ = srv
}

func TestPrintToServerToken(t *testing.T) {
	const token = "erp-0123456789abcdef"
	url, _, _ := printServer(t, "tokens: [ { name: erp, value: "+token+", act-as: true } ]\n")
	doc := filepath.Join("..", "..", "testdata", "form", "rows.bdf")
	t.Setenv(serverEnv, url)
	t.Setenv(tokenEnv, "")
	if code, _, errOut := runOut(t, "print", "-queue", "office", doc); code != exitConnection || !strings.Contains(errOut, "token") {
		t.Fatalf("no token: exit %d: %s", code, errOut)
	}
	tf := filepath.Join(t.TempDir(), "token")
	os.WriteFile(tf, []byte(token+"\n"), 0o600)
	if code, out, errOut := runOut(t, "print", "-queue", "office", "-token-file", tf, "-user", "suzuki", doc); code != 0 || strings.TrimSpace(out) != "1" {
		t.Fatalf("token file: exit %d: %s", code, errOut)
	}
	// the token acts for suzuki, who sees the job; the user of the
	// process has none
	t.Setenv(tokenEnv, token)
	if code, out, _ := runOut(t, "jobs", "-all", "-json", "-user", "suzuki"); code != 0 || !strings.Contains(out, `"user": "suzuki"`) {
		t.Fatalf("jobs: exit %d: %s", code, out)
	}
	if code, out, _ := runOut(t, "jobs", "-all", "-json"); code != 0 || strings.TrimSpace(out) != "[]" {
		t.Fatalf("jobs of the user of the process: exit %d: %s", code, out)
	}
}

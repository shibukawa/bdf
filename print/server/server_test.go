package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/ipp/ipptest"
	"github.com/shibukawa/bdf/print/spool"
)

const (
	erpToken   = "erp-0123456789abcdef"
	aliceToken = "alice-0123456789abcdef"
	adminToken = "admin-0123456789abcdef"
)

func docBytes(t *testing.T) []byte {
	t.Helper()
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	o := bdf.NewObject()
	o.FillColor(0x000000ff).FillRect(50, 100, 100, 40)
	h, _ := d.AddObject(o)
	v.AddPage(595.28, 841.89, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// fixture is a server on a fake printer.
type fixture struct {
	t       *testing.T
	printer *ipptest.Printer
	srv     *Server
	http    *httptest.Server
	cfg     *Config
}

// configYAML is a configuration with a printer at %s and three tokens.
const configYAML = `
printers:
  mfp1: { url: "%s" }
queues:
  office: { printers: [mfp1], defaults: { sides: two-sided-long-edge } }
  warehouse: { printers: [mfp1] }
tokens:
  - { name: erp, value: ` + erpToken + `, queues: [office], roles: [submit], act-as: true }
  - { name: alice, value: ` + aliceToken + `, roles: [submit] }
  - { name: admin, value: ` + adminToken + `, roles: [submit, manage, admin] }
`

func newFixture(t *testing.T, yaml string) *fixture {
	t.Helper()
	p := ipptest.New()
	ps := httptest.NewServer(p)
	t.Cleanup(ps.Close)
	u, _ := url.Parse(ps.URL)
	cfg, err := ParseConfig([]byte(fmt.Sprintf(yaml, "ipp://"+u.Host+"/ipp/print")))
	if err != nil {
		t.Fatal(err)
	}
	cfg.dir = t.TempDir()
	cfg.TempDir = t.TempDir()
	if err := cfg.resolve(); err != nil {
		t.Fatal(err)
	}
	srv, err := New(context.Background(), cfg, Options{Grace: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { srv.Close() })
	hs := httptest.NewServer(srv.Handler())
	t.Cleanup(hs.Close)
	return &fixture{t: t, printer: p, srv: srv, http: hs, cfg: cfg}
}

// do sends a request with a token ("" none) and decodes a JSON answer.
func (f *fixture) do(method, path, token string, body io.Reader, header map[string]string, out any) *http.Response {
	f.t.Helper()
	req, err := http.NewRequest(method, f.http.URL+path, body)
	if err != nil {
		f.t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		f.t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	if out != nil && len(b) > 0 {
		if err := json.Unmarshal(b, out); err != nil {
			f.t.Fatalf("%s %s: %v: %s", method, path, err, b)
		}
	}
	return res
}

// form is a multipart body: the fields in order, "document" a file.
func form(t *testing.T, fields ...string) (io.Reader, string) {
	t.Helper()
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for i := 0; i < len(fields); i += 2 {
		if fields[i] == "document" {
			fw, _ := w.CreateFormFile("document", "invoice.bdf")
			fw.Write([]byte(fields[i+1]))
			continue
		}
		w.WriteField(fields[i], fields[i+1])
	}
	w.Close()
	return &b, w.FormDataContentType()
}

// submit submits a document and returns the status and the job (or the
// error).
func (f *fixture) submit(queue, token string, header map[string]string, fields ...string) (int, spool.Job, apiError) {
	f.t.Helper()
	body, typ := form(f.t, fields...)
	h := map[string]string{"Content-Type": typ}
	for k, v := range header {
		h[k] = v
	}
	var raw json.RawMessage
	res := f.do("POST", "/api/v1/queues/"+queue+"/jobs", token, body, h, &raw)
	var j spool.Job
	var e apiError
	if res.StatusCode == http.StatusCreated {
		json.Unmarshal(raw, &j)
	} else {
		json.Unmarshal(raw, &e)
	}
	return res.StatusCode, j, e
}

func TestConfig(t *testing.T) {
	c, err := ParseConfig([]byte(`
printers: { p: { url: "ipp://192.0.2.1/ipp/print", fingerprint: "sha256:ab" } }
queues: { q: { printers: [p], defaults: { sides: one-sided, copies: 2 } } }
retention: { documents: 7d, records: 400d }
limits: { max-document: 1MiB, max-bytes: 2GiB, retry-for: 90m }
`))
	if err != nil {
		t.Fatal(err)
	}
	if c.Listen != "127.0.0.1:8640" || c.User != "bdf" || time.Duration(c.Retention.Documents) != 7*24*time.Hour ||
		c.Limits.MaxDocument != 1<<20 || c.Limits.MaxBytes != 2<<30 || time.Duration(c.Limits.RetryFor) != 90*time.Minute {
		t.Fatalf("config %+v", c)
	}
	if tk, _ := c.Queues["q"].ticket(); tk.Sides != "one-sided" || tk.Copies != 2 {
		t.Fatalf("defaults %+v", tk)
	}
	if c.pinned()["192.0.2.1"] != "sha256:ab" {
		t.Fatalf("pinned %v", c.pinned())
	}
	bad := map[string]string{
		"no printers":   `queues: { q: { printers: [p] } }`,
		"unknown field": "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }\nlisten2: x",
		"no printer p2": "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p2] } }",
		"defaults":      "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p], defaults: { side: x } } }",
		"role":          "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }\ntokens: [ { name: a, value: x, roles: [root] } ]",
		"no queue":      "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }\ntokens: [ { name: a, value: x, queues: [r] } ]",
		"not loopback":  "listen: 0.0.0.0:8640\nprinters: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }",
		"tls":           "tls: { cert: c.pem }\nprinters: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }",
		"duration":      "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }\nlimits: { retry-for: soon }",
		"size":          "printers: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }\nlimits: { max-document: big }",
	}
	for name, y := range bad {
		if _, err := ParseConfig([]byte(y)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
	if _, err := ParseConfig([]byte("listen: 0.0.0.0:8640\ninsecure: true\nprinters: { p: { url: ipp://h/ } }\nqueues: { q: { printers: [p] } }")); err != nil {
		t.Errorf("insecure: %v", err)
	}

	// secrets: from a file, the environment, never shorter than 16
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "erp.token"), []byte(erpToken+"\n"), 0o600)
	os.WriteFile(filepath.Join(dir, "printd.yaml"), []byte(`
spool: data/spool.db
printers: { p: { url: ipp://h/, password-file: mfp.pw } }
queues: { q: { printers: [p] } }
tokens:
  - { name: erp, file: erp.token }
  - { name: ops, env: BDF_TEST_OPS_TOKEN }
`), 0o600)
	t.Setenv("BDF_TEST_OPS_TOKEN", adminToken)
	c, err = LoadConfig(filepath.Join(dir, "printd.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Spool != filepath.Join(dir, "data", "spool.db") || c.Tokens[0].secret != erpToken || c.Tokens[1].secret != adminToken ||
		c.Printers["p"].PasswordFile != filepath.Join(dir, "mfp.pw") {
		t.Fatalf("resolved %+v %+v", c, c.Tokens)
	}
	t.Setenv("BDF_TEST_OPS_TOKEN", "short")
	if _, err := LoadConfig(filepath.Join(dir, "printd.yaml")); err == nil || !strings.Contains(err.Error(), "16") {
		t.Fatalf("a short secret: %v", err)
	}
}

func TestAuthAndPermissions(t *testing.T) {
	f := newFixture(t, configYAML)
	doc := string(docBytes(t))
	f.srv.Spool().PauseQueue(context.Background(), "office") // the jobs stay

	var e apiError
	res := f.do("GET", "/api/v1/jobs", "", nil, nil, &e)
	if res.StatusCode != 401 || e.Error.Code != "unauthorized" || !strings.HasPrefix(res.Header.Get("WWW-Authenticate"), "Bearer") {
		t.Fatalf("no token: %d %+v", res.StatusCode, e)
	}
	if res := f.do("GET", "/api/v1/jobs", erpToken+"x", nil, nil, nil); res.StatusCode != 401 {
		t.Fatalf("a wrong token: %d", res.StatusCode)
	}
	if st, _, e := f.submit("warehouse", erpToken, nil, "document", doc); st != 403 || e.Error.Code != "forbidden" {
		t.Fatalf("a queue the token may not submit to: %d %+v", st, e)
	}

	// erp submits for its users (act-as); alice may not name another user
	st, mine, _ := f.submit("office", erpToken, map[string]string{userHeader: "suzuki"}, "document", doc)
	if st != 201 || mine.User != "suzuki" || mine.Ticket.Sides != "" {
		t.Fatalf("erp: %d %+v", st, mine)
	}
	st, alices, _ := f.submit("office", aliceToken, map[string]string{userHeader: "suzuki"}, "document", doc, "name", "a.bdf")
	if st != 201 || alices.User != "alice" || alices.Name != "a.bdf" {
		t.Fatalf("alice: %d %+v", st, alices)
	}
	if res := f.do("GET", "/api/v1/jobs", aliceToken, nil, map[string]string{userHeader: strings.Repeat("x", 256)}, nil); res.StatusCode != 401 {
		t.Fatalf("a user name too long: %d", res.StatusCode)
	}

	// alice sees her own job only; admin (manage) sees both
	var list jobList
	f.do("GET", "/api/v1/jobs", aliceToken, nil, nil, &list)
	if len(list.Jobs) != 1 || list.Jobs[0].ID != alices.ID {
		t.Fatalf("alice's list: %+v", list)
	}
	f.do("GET", "/api/v1/jobs?user=suzuki", aliceToken, nil, nil, &list)
	if len(list.Jobs) != 0 {
		t.Fatalf("alice lists suzuki's: %+v", list)
	}
	f.do("GET", "/api/v1/jobs", adminToken, nil, nil, &list)
	if len(list.Jobs) != 2 {
		t.Fatalf("admin's list: %+v", list)
	}
	path := fmt.Sprintf("/api/v1/jobs/%d", mine.ID)
	if res := f.do("GET", path, aliceToken, nil, nil, &e); res.StatusCode != 404 || e.Error.Code != "not-found" {
		t.Fatalf("alice reads suzuki's job: %d", res.StatusCode)
	}
	if res := f.do("POST", path+"/cancel", aliceToken, nil, nil, nil); res.StatusCode != 404 {
		t.Fatalf("alice cancels suzuki's job: %d", res.StatusCode)
	}

	// erp acts on the jobs it submitted for suzuki
	var j spool.Job
	if res := f.do("POST", path+"/hold", erpToken, nil, map[string]string{userHeader: "suzuki"}, &j); res.StatusCode != 200 || j.State != ipp.JobPendingHeld {
		t.Fatalf("hold: %d %+v", res.StatusCode, j)
	}
	if res := f.do("PATCH", path, erpToken, strings.NewReader(`{"priority": 90}`), map[string]string{userHeader: "suzuki"}, &j); res.StatusCode != 200 || j.Priority != 90 {
		t.Fatalf("priority: %d %+v", res.StatusCode, j)
	}
	if res := f.do("POST", path+"/move", erpToken, strings.NewReader(`{"queue": "warehouse"}`), map[string]string{userHeader: "suzuki"}, &e); res.StatusCode != 403 {
		t.Fatalf("move to a queue the token may not use: %d", res.StatusCode)
	}
	if res := f.do("POST", path+"/frobnicate", erpToken, nil, map[string]string{userHeader: "suzuki"}, &e); res.StatusCode != 404 {
		t.Fatalf("an unknown operation: %d", res.StatusCode)
	}

	// admin cancels alice's job, as the operator
	if res := f.do("POST", fmt.Sprintf("/api/v1/jobs/%d/cancel", alices.ID), adminToken, nil, nil, &j); res.StatusCode != 200 ||
		j.State != ipp.JobCanceled || j.Reasons[0] != "job-canceled-by-operator" {
		t.Fatalf("cancel: %d %+v", res.StatusCode, j)
	}
	if res := f.do("POST", fmt.Sprintf("/api/v1/jobs/%d/cancel", alices.ID), adminToken, nil, nil, &e); res.StatusCode != 409 || e.Error.Code != "not-possible" {
		t.Fatalf("cancel twice: %d %+v", res.StatusCode, e)
	}

	// the admin's paths
	for _, p := range []string{"/api/v1/backup", "/api/v1/discover"} {
		if res := f.do("GET", p, aliceToken, nil, nil, &e); res.StatusCode != 403 {
			t.Fatalf("%s for alice: %d", p, res.StatusCode)
		}
	}
	if res := f.do("POST", "/api/v1/queues/office/reject", aliceToken, nil, nil, nil); res.StatusCode != 403 {
		t.Fatalf("reject for alice: %d", res.StatusCode)
	}
	if res := f.do("GET", "/api/v1/accounting", aliceToken, nil, nil, nil); res.StatusCode != 403 {
		t.Fatalf("accounting for alice: %d", res.StatusCode)
	}
	var q spool.QueueInfo
	if res := f.do("POST", "/api/v1/queues/office/reject", adminToken, nil, nil, &q); res.StatusCode != 200 || !q.Rejecting {
		t.Fatalf("reject: %d %+v", res.StatusCode, q)
	}
	if st, _, e := f.submit("office", erpToken, nil, "document", doc); st != 503 || e.Error.Code != "not-accepting" {
		t.Fatalf("a queue that rejects: %d %+v", st, e)
	}
}

func TestNoTokens(t *testing.T) {
	f := newFixture(t, `
printers: { mfp1: { url: "%s" } }
queues: { office: { printers: [mfp1] } }
`)
	f.srv.Spool().PauseQueue(context.Background(), "office")
	st, j, _ := f.submit("office", "", map[string]string{userHeader: "sato"}, "document", string(docBytes(t)))
	if st != 201 || j.User != "sato" {
		t.Fatalf("%d %+v", st, j)
	}
	st, j, _ = f.submit("office", "", nil, "document", string(docBytes(t)))
	if st != 201 || j.User != "anonymous" {
		t.Fatalf("%d %+v", st, j)
	}
	if res := f.do("POST", "/api/v1/queues/office/accept", "", nil, nil, nil); res.StatusCode != 200 {
		t.Fatalf("admin without tokens: %d", res.StatusCode)
	}
}

func TestSubmitErrors(t *testing.T) {
	f := newFixture(t, configYAML+"limits: { max-document: 4KiB, max-jobs: 1 }\n")
	f.srv.Spool().PauseQueue(context.Background(), "office")
	doc := string(docBytes(t))
	cases := []struct {
		name   string
		fields []string
		status int
		code   string
	}{
		{"no document", []string{"name", "x"}, 400, "invalid"},
		{"too large", []string{"document", strings.Repeat("x", 5000)}, 413, "too-large"},
		{"unknown format", []string{"document", "%PDF-1.7"}, 415, "unknown-format"},
		{"empty", []string{"document", ""}, 400, "bad-document"},
		{"bad ticket", []string{"ticket", `{"sidez": "one-sided"}`, "document", doc}, 400, "invalid"},
		{"ticket after", []string{"document", doc, "ticket", `{"copies": -1}`}, 400, "invalid"},
		{"unknown field", []string{"document", doc, "colour", "red"}, 400, "invalid"},
		{"priority", []string{"priority", "101", "document", doc}, 400, "invalid"},
		{"hold", []string{"hold-until", "lunch", "document", doc}, 400, "invalid"},
	}
	for _, c := range cases {
		st, _, e := f.submit("office", adminToken, nil, c.fields...)
		if st != c.status || e.Error.Code != c.code || e.Error.Message == "" {
			t.Errorf("%s: %d %+v", c.name, st, e)
		}
	}
	if st, _, e := f.submit("nowhere", adminToken, nil, "document", doc); st != 404 || e.Error.Code != "no-queue" {
		t.Errorf("no queue: %d %+v", st, e)
	}
	// the ticket after the document is taken
	st, j, _ := f.submit("office", adminToken, nil, "document", doc, "ticket", `{"copies": 2, "sides": "one-sided"}`, "priority", "70", "hold-until", "indefinite")
	if st != 201 || j.Ticket.Copies != 2 || j.Ticket.Sides != "one-sided" || j.Priority != 70 || j.State != ipp.JobPendingHeld || j.Name != "invoice.bdf" {
		t.Fatalf("%d %+v", st, j)
	}
	if st, _, e := f.submit("office", adminToken, nil, "document", doc); st != 507 || e.Error.Code != "full" {
		t.Errorf("full: %d %+v", st, e)
	}
	// nothing is left of the documents refused
	left, _ := os.ReadDir(f.cfg.TempDir)
	if len(left) != 0 {
		t.Errorf("left in the temporary directory: %v", left)
	}
}

func TestPrintAndFollow(t *testing.T) {
	f := newFixture(t, configYAML)
	f.srv.Spool().Start()

	// the events of the jobs alice may see
	req, _ := http.NewRequest("GET", f.http.URL+"/api/v1/events", nil)
	req.Header.Set("Authorization", "Bearer "+aliceToken)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 || res.Header.Get("Content-Type") != "text/event-stream" {
		t.Fatalf("events: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
	}
	events := make(chan string, 100)
	go func() {
		sc := bufio.NewScanner(res.Body)
		var ev string
		for sc.Scan() {
			line := sc.Text()
			switch {
			case strings.HasPrefix(line, "event: "):
				ev = line[7:]
			case strings.HasPrefix(line, "data: "):
				events <- ev + " " + line[6:]
			}
		}
		close(events)
	}()

	doc := string(docBytes(t))
	_, other, _ := f.submit("office", erpToken, nil, "document", doc)
	st, j, _ := f.submit("office", aliceToken, nil, "document", doc)
	if st != 201 || j.Ticket.Sides != "" {
		t.Fatalf("%d %+v", st, j)
	}
	deadline := time.After(15 * time.Second)
	for done := false; !done; {
		select {
		case ev, ok := <-events:
			if !ok {
				t.Fatal("the stream ended")
			}
			if strings.HasPrefix(ev, "job ") {
				var got spool.Job
				json.Unmarshal([]byte(ev[4:]), &got)
				if got.ID == other.ID {
					t.Fatal("alice is told of erp's job")
				}
				done = got.ID == j.ID && got.State == ipp.JobCompleted
			}
		case <-deadline:
			t.Fatal("no event of the job done")
		}
	}
	var got spool.Job
	f.do("GET", fmt.Sprintf("/api/v1/jobs/%d", j.ID), aliceToken, nil, nil, &got)
	if got.State != ipp.JobCompleted || got.Printer == nil || got.Printer.Name != "mfp1" {
		t.Fatalf("job %+v", got)
	}
	// the queue's defaults went to the printer
	pj := f.printer.Job(got.Printer.JobID)
	if pj == nil {
		t.Fatal("the printer has not the job")
	}
	for _, a := range pj.Attrs {
		if a.Name == "sides" && a.Values[0].String() != "two-sided-long-edge" {
			t.Fatalf("sides %v", a.Values)
		}
	}

	// the lists, the queues, the printers
	var list jobList
	f.do("GET", "/api/v1/jobs?state=completed", aliceToken, nil, nil, &list)
	if len(list.Jobs) != 1 || list.Jobs[0].ID != j.ID {
		t.Fatalf("completed: %+v", list)
	}
	var queues struct{ Queues []spool.QueueInfo }
	f.do("GET", "/api/v1/queues", aliceToken, nil, nil, &queues)
	if len(queues.Queues) != 2 || queues.Queues[0].Name != "office" {
		t.Fatalf("queues %+v", queues)
	}
	var pr struct {
		spool.PrinterInfo
		Capabilities *capabilities
	}
	f.do("GET", "/api/v1/printers/mfp1", aliceToken, nil, nil, &pr)
	if pr.Name != "mfp1" || pr.Capabilities == nil || len(pr.Capabilities.Formats) == 0 {
		t.Fatalf("printer %+v", pr)
	}

	// metrics and accounting
	res2, err := http.Get(f.http.URL + "/metrics")
	if err != nil || res2.StatusCode != 401 {
		t.Fatalf("metrics without token: %v %v", err, res2.StatusCode)
	}
	req, _ = http.NewRequest("GET", f.http.URL+"/metrics", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res2, _ = http.DefaultClient.Do(req)
	b, _ := io.ReadAll(res2.Body)
	res2.Body.Close()
	for _, want := range []string{
		`bdf_print_jobs_total{queue="office",state="completed"} 2`,
		`bdf_print_sheets_total{queue="office"} 2`,
		`bdf_print_printer_up{printer="mfp1"} 1`,
		`bdf_print_printer_state{printer="mfp1"} 3`,
		"# TYPE bdf_print_spool_file_bytes gauge",
	} {
		if !bytes.Contains(b, []byte(want)) {
			t.Errorf("metrics lack %s:\n%s", want, b)
		}
	}
	req, _ = http.NewRequest("GET", f.http.URL+"/api/v1/accounting?format=csv&from=2000-01-01", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res2, _ = http.DefaultClient.Do(req)
	b, _ = io.ReadAll(res2.Body)
	res2.Body.Close()
	if lines := strings.Split(strings.TrimSpace(string(b)), "\n"); res2.Header.Get("Content-Type") != "text/csv; charset=utf-8" || len(lines) != 3 || !strings.HasPrefix(lines[0], "id,") {
		t.Fatalf("accounting: %s", b)
	}
}

func TestListPages(t *testing.T) {
	f := newFixture(t, configYAML)
	f.srv.Spool().PauseQueue(context.Background(), "office")
	doc := string(docBytes(t))
	for i := 0; i < 5; i++ {
		f.submit("office", adminToken, nil, "document", doc)
	}
	var ids []int32
	after := int32(0)
	for {
		var list jobList
		f.do("GET", fmt.Sprintf("/api/v1/jobs?limit=2&after=%d", after), adminToken, nil, nil, &list)
		for _, j := range list.Jobs {
			ids = append(ids, j.ID)
		}
		if list.Next == 0 {
			break
		}
		after = list.Next
	}
	if fmt.Sprint(ids) != "[1 2 3 4 5]" {
		t.Fatalf("pages: %v", ids)
	}
	var e apiError
	for _, q := range []string{"limit=0", "limit=1001", "after=x", "state=printing"} {
		if res := f.do("GET", "/api/v1/jobs?"+q, adminToken, nil, nil, &e); res.StatusCode != 400 {
			t.Errorf("%s: %d", q, res.StatusCode)
		}
	}
}

func TestBackup(t *testing.T) {
	f := newFixture(t, configYAML)
	f.srv.Spool().PauseQueue(context.Background(), "office")
	_, j, _ := f.submit("office", adminToken, nil, "document", string(docBytes(t)))
	req, _ := http.NewRequest("GET", f.http.URL+"/api/v1/backup", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != 200 || !strings.HasPrefix(res.Header.Get("Content-Disposition"), "attachment") || res.ContentLength != int64(len(b)) {
		t.Fatalf("backup: %d %v %d", res.StatusCode, res.Header, len(b))
	}
	// the copy is a spool with the job
	path := filepath.Join(t.TempDir(), "copy.db")
	os.WriteFile(path, b, 0o600)
	o, _ := f.cfg.spoolOptions()
	o.Path = path
	sp, err := spool.Open(context.Background(), o)
	if err != nil {
		t.Fatal(err)
	}
	defer sp.Close()
	if got, err := sp.Get(context.Background(), j.ID); err != nil || got.UUID != j.UUID {
		t.Fatalf("the copy: %v %+v", err, got)
	}
}

func TestPageAndHealth(t *testing.T) {
	f := newFixture(t, configYAML)
	for path, typ := range map[string]string{"/": "text/html; charset=utf-8", "/app.js": "text/javascript; charset=utf-8"} {
		res, err := http.Get(f.http.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != 200 || res.Header.Get("Content-Type") != typ || !strings.Contains(res.Header.Get("Content-Security-Policy"), "script-src 'self'") {
			t.Fatalf("%s: %d %v", path, res.StatusCode, res.Header)
		}
	}
	var h map[string]string
	if res := f.do("GET", "/healthz", "", nil, nil, &h); res.StatusCode != 200 || h["status"] != "ok" {
		t.Fatalf("healthz %d %v", res.StatusCode, h)
	}
	var e apiError
	if res := f.do("GET", "/api/v1/nothing", adminToken, nil, nil, &e); res.StatusCode != 404 || e.Error.Code != "not-found" {
		t.Fatalf("unknown path %d %+v", res.StatusCode, e)
	}
}

func TestServeStops(t *testing.T) {
	rt := readTimeout
	readTimeout = 300 * time.Millisecond
	t.Cleanup(func() { readTimeout = rt })
	p := ipptest.New()
	ps := httptest.NewServer(p)
	defer ps.Close()
	u, _ := url.Parse(ps.URL)
	cfg, err := ParseConfig([]byte(fmt.Sprintf("printers: { mfp1: { url: \"ipp://%s/ipp/print\" } }\nqueues: { office: { printers: [mfp1] } }\n", u.Host)))
	if err != nil {
		t.Fatal(err)
	}
	cfg.dir = t.TempDir()
	cfg.resolve()
	srv, err := New(context.Background(), cfg, Options{Grace: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	// a second server of the same spool is refused
	if _, err := New(context.Background(), cfg, Options{}); err == nil {
		t.Fatal("a second server of the spool")
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ctx, ln) }()
	base := "http://" + ln.Addr().String()
	res, err := http.Get(base + "/api/v1/events")
	if err != nil || res.StatusCode != 200 {
		t.Fatalf("events: %v", err)
	}
	lines := make(chan string, 100)
	ended := make(chan struct{})
	go func() {
		sc := bufio.NewScanner(res.Body)
		for sc.Scan() {
			lines <- sc.Text()
		}
		close(ended)
	}()
	// the stream outlives the time a request has to arrive
	time.Sleep(3 * readTimeout)
	body, typ := form(t, "document", string(docBytes(t)))
	if res, err := http.Post(base+"/api/v1/queues/office/jobs", typ, body); err != nil || res.StatusCode != 201 {
		t.Fatalf("submit: %v", err)
	}
	for got := false; !got; {
		select {
		case l := <-lines:
			got = l == "event: job"
		case <-ended:
			t.Fatal("the event stream ended")
		case <-time.After(5 * time.Second):
			t.Fatal("no event")
		}
	}
	cancel()
	select {
	case err := <-served:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("Serve did not return")
	}
	select {
	case <-ended:
	case <-time.After(5 * time.Second):
		t.Fatal("the event stream did not end")
	}
	// the spool is closed: another server opens it
	srv2, err := New(context.Background(), cfg, Options{})
	if err != nil {
		t.Fatal(err)
	}
	srv2.Close()
}

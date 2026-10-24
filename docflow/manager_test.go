package docflow

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/ipp"
)

const a4w, a4h = 595.276, 841.89

// testDoc is a .bdf file of n pages.
func testDoc(t *testing.T, n int) []byte {
	t.Helper()
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	for i := range n {
		o := bdf.NewObject()
		o.FillColor(0x000000ff).FillRect(50, 100, float32(50*(i+1)), 40)
		h, _ := d.AddObject(o)
		v.AddPage(a4w, a4h, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// clock is a test's time.
type clock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *clock) add(d time.Duration) {
	c.mu.Lock()
	c.t = c.t.Add(d)
	c.mu.Unlock()
}

const adminToken = "admin-token-for-tests"

// env is a manager's store served over HTTP. Every request makes a
// manager of its own, as a function of a serverless platform would: what
// works here keeps no state in an instance.
type env struct {
	t     *testing.T
	st    *Store
	clock *clock
	o     Options
	srv   *httptest.Server
}

func newEnv(t *testing.T, edit func(*Options)) *env {
	t.Helper()
	return newEnvAt(t, "bbolt:"+filepath.ToSlash(filepath.Join(t.TempDir(), "pm.db")), "", edit)
}

// newEnvAt is an env of a store given by its URLs.
func newEnvAt(t *testing.T, records, documents string, edit func(*Options)) *env {
	t.Helper()
	st, err := OpenStore(context.Background(), records, documents)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	e := &env{t: t, st: st, clock: &clock{t: time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)}}
	e.o = Options{Secret: bytes.Repeat([]byte{7}, 32), AdminToken: adminToken, TickToken: "tick-token", Now: e.clock.Now, TempDir: t.TempDir(),
		NoSystemFonts: true}
	if edit != nil {
		edit(&e.o)
	}
	e.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Handler(e.m()).ServeHTTP(w, r)
	}))
	t.Cleanup(e.srv.Close)
	return e
}

func (e *env) m() *Manager {
	m, err := New(e.st, e.o)
	if err != nil {
		e.t.Fatal(err)
	}
	return m
}

// do sends a request; body is JSON unless it is bytes.
func (e *env) do(c *http.Client, method, path, token string, body any) (*http.Response, []byte) {
	e.t.Helper()
	return e.send(c, func(r *http.Request) {
		if token != "" {
			r.Header.Set("Authorization", "Bearer "+token)
		}
	}, method, path, body)
}

// doWith sends a request that edit changes first.
func (e *env) doWith(edit func(*http.Request), method, path string, body any) (*http.Response, []byte) {
	e.t.Helper()
	return e.send(nil, edit, method, path, body)
}

func (e *env) send(c *http.Client, edit func(*http.Request), method, path string, body any) (*http.Response, []byte) {
	e.t.Helper()
	var r io.Reader
	switch b := body.(type) {
	case nil:
	case []byte:
		r = bytes.NewReader(b)
	default:
		j, _ := json.Marshal(b)
		r = bytes.NewReader(j)
	}
	req, err := http.NewRequest(method, e.srv.URL+path, r)
	if err != nil {
		e.t.Fatal(err)
	}
	if _, raw := body.([]byte); body != nil && !raw {
		req.Header.Set("Content-Type", "application/json")
	}
	edit(req)
	if c == nil {
		c = http.DefaultClient
	}
	res, err := c.Do(req)
	if err != nil {
		e.t.Fatal(err)
	}
	defer res.Body.Close()
	out, _ := io.ReadAll(res.Body)
	return res, out
}

// admin sends a request of the API and wants status.
func (e *env) admin(method, path string, body any, status int) []byte {
	e.t.Helper()
	res, out := e.do(nil, method, path, adminToken, body)
	if res.StatusCode != status {
		e.t.Fatalf("%s %s: %d %s; want %d", method, path, res.StatusCode, out, status)
	}
	return out
}

// device is a tablet of a test: its key, and a client with its cookie.
type device struct {
	e   *env
	id  string
	key *ecdsa.PrivateKey
	c   *http.Client
	// etag is the board's revision the device has.
	etag string
}

// enroll adds a device to a location, as its page does.
func (e *env) enroll(loc, name string) *device {
	e.t.Helper()
	var en Enrollment
	json.Unmarshal(e.admin("POST", "/api/nodes/"+loc+"/enrollments", nil, http.StatusCreated), &en)
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	jar, _ := cookiejar.New(nil)
	d := &device{e: e, key: key, c: &http.Client{Jar: jar}}
	res, out := e.do(d.c, "POST", "/device/enroll", "", map[string]string{"code": strings.ToLower(en.Code), "name": name, "key": base64.StdEncoding.EncodeToString(der)})
	if res.StatusCode != http.StatusCreated {
		e.t.Fatalf("enroll: %d %s", res.StatusCode, out)
	}
	var got map[string]string
	json.Unmarshal(out, &got)
	d.id = got["device"]
	if got["location"] != loc {
		e.t.Fatalf("enrolled at %q", got["location"])
	}
	d.login()
	return d
}

// sign signs a challenge as WebCrypto does: r and s, 32 bytes each.
func (d *device) sign(challenge string) string {
	sum := sha256.Sum256([]byte(challenge))
	r, s, err := ecdsa.Sign(rand.Reader, d.key, sum[:])
	if err != nil {
		d.e.t.Fatal(err)
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return base64.StdEncoding.EncodeToString(sig)
}

func (d *device) login() {
	d.e.t.Helper()
	_, out := d.e.do(d.c, "GET", "/device/challenge", "", nil)
	var ch map[string]string
	json.Unmarshal(out, &ch)
	res, out := d.e.do(d.c, "POST", "/device/session", "", map[string]string{"device": d.id, "challenge": ch["challenge"], "signature": d.sign(ch["challenge"])})
	if res.StatusCode != http.StatusOK {
		d.e.t.Fatalf("session: %d %s", res.StatusCode, out)
	}
}

// board polls the board: nil when it did not change.
func (d *device) board() *Board {
	d.e.t.Helper()
	req, _ := http.NewRequest("GET", d.e.srv.URL+"/device/board", nil)
	if d.etag != "" {
		req.Header.Set("If-None-Match", d.etag)
	}
	res, err := d.c.Do(req)
	if err != nil {
		d.e.t.Fatal(err)
	}
	defer res.Body.Close()
	out, _ := io.ReadAll(res.Body)
	switch res.StatusCode {
	case http.StatusNotModified:
		return nil
	case http.StatusOK:
	default:
		d.e.t.Fatalf("board: %d %s", res.StatusCode, out)
	}
	d.etag = res.Header.Get("ETag")
	var b Board
	if err := json.Unmarshal(out, &b); err != nil {
		d.e.t.Fatal(err)
	}
	return &b
}

func (d *device) action(id, action string, status int) *Slip {
	d.e.t.Helper()
	res, out := d.e.do(d.c, "POST", "/device/slips/"+id+"/"+action, "", map[string]string{"worker": "Sato"})
	if res.StatusCode != status {
		d.e.t.Fatalf("%s %s: %d %s; want %d", action, id, res.StatusCode, out, status)
	}
	var s Slip
	json.Unmarshal(out, &s)
	return &s
}

// printer is a fake printer the manager prints on.
func (e *env) printer(name string) *fakeprinter.Printer {
	p := fakeprinter.New()
	e.admin("PUT", "/api/nodes/"+name, Node{Kind: KindPrinter, URL: fakeprinter.Start(e.t, p)}, http.StatusOK)
	return p
}

// display puts a display node.
func (e *env) display(name string, n Node) {
	e.t.Helper()
	n.Kind = KindDisplay
	e.admin("PUT", "/api/nodes/"+name, n, http.StatusOK)
}

// ippClient prints to a CUPS server node's printer on the manager's own
// server (/ipp/print/<name>): IPP enters the topology there alone.
func (e *env) ippClient(name string) *ipp.Client {
	return &ipp.Client{URL: e.srv.URL + IPPPath + "/" + name, User: "wms"}
}

// cups puts a CUPS server node that forwards to the nodes of to, and is
// the IPP printer computers print to in front of them.
func (e *env) cups(name string, to ...string) *ipp.Client {
	e.t.Helper()
	e.node(name, Node{Kind: KindCups, To: to}, http.StatusOK)
	return e.ippClient(name)
}

// printJob sends a .bdf of pages to a destination as an application
// does, over the API: as the user wms, anonymous (no submit token is
// set).
func (e *env) printJob(dest, name string, pages int) *Job {
	e.t.Helper()
	res, out := e.do(nil, "POST", "/api/destinations/"+dest+"/jobs?name="+url.QueryEscape(name)+"&user=wms", "", testDoc(e.t, pages))
	if res.StatusCode != http.StatusCreated {
		e.t.Fatalf("print %q to %s: %d %s", name, dest, res.StatusCode, out)
	}
	var j Job
	if err := json.Unmarshal(out, &j); err != nil {
		e.t.Fatal(err)
	}
	return &j
}

func TestDisplay(t *testing.T) {
	e := newEnv(t, nil)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Title: "第 1 倉庫"}, http.StatusOK)
	d := e.enroll("wh-a", "tablet 1")

	b := d.board()
	if b == nil || len(b.Open) != 0 || b.Title != "第 1 倉庫" {
		t.Fatalf("empty board %+v", b)
	}
	if d.board() != nil {
		t.Fatal("the board changed by itself")
	}

	// the WMS sends it, as an application does
	j := e.printJob("wh-a", "picking list 1", 2)
	if j.Number != 1 || j.State != JobCompleted {
		t.Fatalf("job %+v", j)
	}
	b = d.board()
	if b == nil || len(b.Open) != 1 {
		t.Fatalf("board %+v", b)
	}
	s := b.Open[0]
	if s.Name != "picking list 1" || s.Pages != 2 || s.State != SlipNew || s.Source != "api" ||
		s.ClaimedUser != "wms" || s.SubmittedBy == nil || s.SubmittedBy.Kind != PrincipalAnonymous {
		t.Fatalf("slip %+v", s)
	}
	if !s.Expires.Equal(s.Received.Add(14 * 24 * time.Hour)) {
		t.Fatalf("expires %v, received %v", s.Expires, s.Received)
	}

	// its pages, as images
	res, out := e.do(d.c, "GET", "/device/slips/"+s.ID+"/pages/2.png?w=400", "", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("page: %d %s", res.StatusCode, out)
	}
	img, err := png.Decode(bytes.NewReader(out))
	if err != nil || img.Bounds().Dx() != 400 {
		t.Fatalf("page image %v %v", err, img.Bounds())
	}
	if res, _ := e.do(d.c, "GET", "/device/slips/"+s.ID+"/pages/3.png", "", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("page 3 of 2: %d", res.StatusCode)
	}

	if got := d.action(s.ID, "open", http.StatusOK); got.State != SlipOpened || got.Opened.Device != d.id || got.Opened.Worker != "Sato" {
		t.Fatalf("opened %+v", got)
	}
	if b = d.board(); b == nil || b.Open[0].State != SlipOpened {
		t.Fatalf("board after open %+v", b)
	}
	d.action(s.ID, "done", http.StatusOK)
	if b = d.board(); b == nil || len(b.Open) != 0 || len(b.Ended) != 1 || b.Ended[0].State != SlipDone {
		t.Fatalf("board after done %+v", b)
	}
	d.action(s.ID, "done", http.StatusConflict)

	// done, it opens again as often as asked, until it expires
	for range 2 {
		if res, _ := e.do(d.c, "GET", "/device/slips/"+s.ID+"/pages/1.png", "", nil); res.StatusCode != http.StatusOK {
			t.Fatalf("page of a done slip: %d", res.StatusCode)
		}
		if got := d.action(s.ID, "open", http.StatusOK); got.State != SlipDone || got.Opened == nil {
			t.Fatalf("opened again %+v", got)
		}
	}

	// a computer prints through a CUPS server, as to any IPP printer: it
	// sees its job there; one it cancels leaves the board
	ctx := context.Background()
	c := e.cups("cups1", "wh-a")
	js, err := c.PrintJob(ctx, "picking list 2", "application/x-bdf", nil, bytes.NewReader(testDoc(t, 1)))
	if err != nil || js.ID != 2 || js.State != ipp.JobCompleted {
		t.Fatalf("job 2: %+v %v", js, err)
	}
	got, err := c.GetJobAttributes(ctx, 2)
	if err != nil || got.State != ipp.JobCompleted || got.Name != "picking list 2" {
		t.Fatalf("job 2: %+v %v", got, err)
	}
	if _, err := c.GetJobAttributes(ctx, 1); !ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
		t.Fatalf("job 1 is not the CUPS server's: %v", err)
	}
	if err := c.CancelJob(ctx, 2); err != nil {
		t.Fatal(err)
	}
	if err := c.CancelJob(ctx, 2); !ipp.IsStatus(err, ipp.StatusClientErrorNotPossible) {
		t.Fatalf("cancel twice: %v", err)
	}
	if res, out := e.do(nil, "POST", "/api/jobs/1/cancel", adminToken, nil); res.StatusCode != http.StatusConflict {
		t.Fatalf("cancel a done job: %d %s", res.StatusCode, out)
	}
	b = d.board()
	if b == nil || len(b.Open) != 0 || len(b.Ended) != 2 || b.Ended[0].State != SlipCanceled {
		t.Fatalf("board after cancel %+v", b)
	}
	// the sender took it back: it is not shown again
	if res, _ := e.do(d.c, "GET", "/device/slips/"+b.Ended[0].ID+"/pages/1.png", "", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("page of a canceled slip: %d", res.StatusCode)
	}
	jobs, err := c.GetJobs(ctx, "completed", false, 10)
	if err != nil || len(jobs) != 1 || jobs[0].ID != 2 || jobs[0].State != ipp.JobCanceled {
		t.Fatalf("jobs %+v %v", jobs, err)
	}
	var sent []Job
	json.Unmarshal(e.admin("GET", "/api/destinations/wh-a/jobs", nil, http.StatusOK), &sent)
	if len(sent) != 1 || sent[0].Number != 1 {
		t.Fatalf("jobs of wh-a %+v", sent)
	}

	// one that came a day later is opened, not done
	e.clock.add(24 * time.Hour)
	d.login() // the session ended meanwhile
	e.printJob("wh-a", "picking list 3", 1)
	b = d.board()
	if b == nil || len(b.Open) != 1 {
		t.Fatalf("board with job 3 %+v", b)
	}
	s3 := b.Open[0]
	d.action(s3.ID, "open", http.StatusOK)

	// the slips leave when they expire, the location's keep after they
	// came, whatever was done to them: on the board at once, and at the
	// next tick
	e.clock.add(13*24*time.Hour + time.Minute)
	d.login()
	if b = d.board(); b == nil || len(b.Open) != 1 || len(b.Ended) != 0 {
		t.Fatalf("board once 1 and 2 expired %+v", b)
	}
	if res, _ := e.do(d.c, "GET", "/device/slips/"+s.ID+"/pages/1.png", "", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("page of an expired slip: %d", res.StatusCode)
	}
	var rep TickReport
	json.Unmarshal(e.admin("POST", "/api/tick", nil, http.StatusOK), &rep)
	if rep.Removed != 2 {
		t.Fatalf("tick %+v", rep)
	}
	if _, err := e.m().Slip(context.Background(), s.ID); err == nil {
		t.Fatal("the slip is still there")
	}
	if res, _ := e.do(d.c, "GET", "/device/slips/"+s3.ID+"/pages/1.png", "", nil); res.StatusCode != http.StatusOK {
		t.Fatalf("page of slip 3: %d", res.StatusCode)
	}
	e.clock.add(24 * time.Hour)
	json.Unmarshal(e.admin("POST", "/api/tick", nil, http.StatusOK), &rep)
	if rep.Removed != 1 {
		t.Fatalf("tick %+v", rep)
	}
	d.login()
	if b = d.board(); b == nil || len(b.Open) != 0 || len(b.Ended) != 0 {
		t.Fatalf("board after keep %+v", b)
	}
}

func TestSubmitAPI(t *testing.T) {
	e := newEnv(t, func(o *Options) { o.SubmitTokens = map[string]string{"wms": "s3cret"} })
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
	doc := testDoc(t, 1)
	if res, _ := e.do(nil, "POST", "/api/locations/wh-a/slips?name=x", "", doc); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("no token: %d", res.StatusCode)
	}
	res, out := e.do(nil, "POST", "/api/locations/wh-a/slips?name=invoice", "wms:s3cret", doc)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("submit: %d %s", res.StatusCode, out)
	}
	var s Slip
	json.Unmarshal(out, &s)
	if s.SubmittedBy == nil || *s.SubmittedBy != (Principal{Kind: PrincipalToken, ID: "wms", Name: "wms", Via: ViaSubmitToken}) || s.Source != "api" || s.Pages != 1 {
		t.Fatalf("slip %+v", s)
	}
	if res, _ := e.do(nil, "POST", "/api/locations/wh-a/slips", "wms:s3cret", []byte("not a document")); res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("not a document: %d", res.StatusCode)
	}
	if res, _ := e.do(nil, "POST", "/api/locations/nowhere/slips", "wms:s3cret", doc); res.StatusCode != http.StatusNotFound {
		t.Fatalf("no location: %d", res.StatusCode)
	}
	// IPP, at a CUPS server, asks for the token too
	c := e.cups("cups1", "wh-a")
	if _, err := c.PrintJob(context.Background(), "x", "", nil, bytes.NewReader(doc)); err == nil {
		t.Fatal("IPP without a token")
	}
	c.Username, c.Password = "wms", "s3cret"
	if js, err := c.PrintJob(context.Background(), "x", "", nil, bytes.NewReader(doc)); err != nil || js.State != ipp.JobCompleted {
		t.Fatalf("IPP with the token: %+v %v", js, err)
	}
	got, out := e.do(nil, "GET", "/api/slips/"+s.ID+"/document", adminToken, nil)
	if got.StatusCode != http.StatusOK || !bytes.HasPrefix(out, bdf.Magic[:]) {
		t.Fatalf("document: %d", got.StatusCode)
	}
	for _, path := range []string{"/api/nodes", "/api/topology", "/api/slips/" + s.ID} {
		if res, _ := e.do(nil, "GET", path, "", nil); res.StatusCode != http.StatusUnauthorized {
			t.Fatalf("%s without the admin token: %d", path, res.StatusCode)
		}
	}
}

func TestFallback(t *testing.T) {
	e := newEnv(t, nil)
	p := e.printer("wh-a-laser")
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Fallback: &Fallback{Printer: "wh-a-laser", After: Duration(10 * time.Minute)}}, http.StatusOK)

	// no device is connected: the slip is printed at the next tick
	e.printJob("wh-a", "offline", 1)
	var rep TickReport
	json.Unmarshal(e.admin("POST", "/api/tick", nil, http.StatusOK), &rep)
	if rep.Printed != 1 || len(p.Jobs()) != 1 {
		t.Fatalf("offline: %+v, %d jobs", rep, len(p.Jobs()))
	}
	if j := p.Jobs()[0]; j.Name != "offline" {
		t.Fatalf("printed %q", j.Name)
	}
	s, _ := e.m().SlipOfJob(context.Background(), 1)
	if len(s.Printed) != 1 || s.Printed[0].Why != whyFallback || s.State != SlipNew {
		t.Fatalf("slip %+v", s)
	}
	// once
	json.Unmarshal(e.admin("POST", "/api/tick", nil, http.StatusOK), &rep)
	if rep.Printed != 0 || len(p.Jobs()) != 1 {
		t.Fatalf("again: %+v", rep)
	}

	// a device is connected: the slip waits for it, then is printed
	d := e.enroll("wh-a", "tablet")
	d.board()
	e.printJob("wh-a", "online", 1)
	e.clock.add(9 * time.Minute)
	d.login()
	d.board()
	e.m().Tick(context.Background())
	if len(p.Jobs()) != 1 {
		t.Fatal("printed before its time")
	}
	e.clock.add(2 * time.Minute)
	d.board()
	r, err := e.m().Tick(context.Background())
	if err != nil || r.Printed != 1 || len(p.Jobs()) != 2 {
		t.Fatalf("after 11 minutes: %+v %v", r, err)
	}

	// an opened slip is not printed; the device prints one on demand
	e.printJob("wh-a", "opened", 1)
	b := d.board()
	opened := b.Open[len(b.Open)-1]
	d.action(opened.ID, "open", http.StatusOK)
	e.clock.add(time.Hour)
	d.login()
	d.board()
	e.m().Tick(context.Background())
	if len(p.Jobs()) != 2 {
		t.Fatal("an opened slip was printed")
	}
	got := d.action(opened.ID, "print", http.StatusOK)
	if len(p.Jobs()) != 3 || len(got.Printed) != 1 || got.Printed[0].Why != whyDevice || got.Printed[0].Device != d.id {
		t.Fatalf("printed on demand: %+v", got)
	}
}

func TestFallbackFails(t *testing.T) {
	e := newEnv(t, nil)
	e.admin("PUT", "/api/nodes/gone", Node{Kind: KindPrinter, URL: "ipp://127.0.0.1:1/ipp/print"}, http.StatusOK)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Fallback: &Fallback{Printer: "gone"}}, http.StatusOK)
	e.printJob("wh-a", "x", 1)
	r, err := e.m().Tick(context.Background())
	if err == nil || r.Failed != 1 {
		t.Fatalf("tick %+v %v", r, err)
	}
	s, _ := e.m().SlipOfJob(context.Background(), 1)
	if s.Attempts != 1 || s.LastError == "" || s.Lease != nil || !s.RetryAt.After(e.clock.Now()) {
		t.Fatalf("slip %+v", s)
	}
	// not again before RetryAt
	if r, _ := e.m().Tick(context.Background()); r.Failed != 0 {
		t.Fatalf("retried at once: %+v", r)
	}
	e.clock.add(time.Minute)
	if r, _ := e.m().Tick(context.Background()); r.Failed != 1 {
		t.Fatalf("not retried: %+v", r)
	}
	// a device asking is told why
	d := e.enroll("wh-a", "tablet")
	res, out := e.do(d.c, "POST", "/device/slips/"+s.ID+"/print", "", nil)
	if res.StatusCode != http.StatusBadGateway || !strings.Contains(string(out), "gone") {
		t.Fatalf("print on demand: %d %s", res.StatusCode, out)
	}
}

// Ticks at once (instances of a function, a timer and a scheduler) print
// a slip once.
func TestTicksAtOnce(t *testing.T) {
	e := newEnv(t, nil)
	p := e.printer("laser")
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay, Fallback: &Fallback{Printer: "laser"}}, http.StatusOK)
	for i := range 3 {
		e.printJob("wh-a", fmt.Sprint("slip ", i), 1)
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			res, out := e.do(nil, "POST", "/api/tick", "tick-token", nil)
			if res.StatusCode != http.StatusOK {
				t.Errorf("tick: %d %s", res.StatusCode, out)
			}
		})
	}
	wg.Wait()
	if n := len(p.Jobs()); n != 3 {
		t.Fatalf("%d jobs printed; want 3", n)
	}
}

func TestDevices(t *testing.T) {
	e := newEnv(t, nil)
	e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
	e.admin("PUT", "/api/nodes/wh-b", Node{Kind: KindDisplay}, http.StatusOK)
	a := e.enroll("wh-a", "a")
	b := e.enroll("wh-b", "b")
	e.printJob("wh-a", "for a", 1)
	slip := a.board().Open[0]

	// a slip of another location is not there for a device
	if res, _ := e.do(b.c, "GET", "/device/slips/"+slip.ID+"/pages/1.png", "", nil); res.StatusCode != http.StatusNotFound {
		t.Fatalf("another location's slip: %d", res.StatusCode)
	}
	b.action(slip.ID, "done", http.StatusNotFound)

	// a code is used once, and not after it expires
	var en Enrollment
	json.Unmarshal(e.admin("POST", "/api/nodes/wh-a/enrollments?ttl=5m", nil, http.StatusCreated), &en)
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	der, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	body := map[string]string{"code": en.Code, "key": base64.StdEncoding.EncodeToString(der)}
	if res, out := e.do(nil, "POST", "/device/enroll", "", body); res.StatusCode != http.StatusCreated {
		t.Fatalf("enroll: %d %s", res.StatusCode, out)
	}
	if res, _ := e.do(nil, "POST", "/device/enroll", "", body); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("a code used twice: %d", res.StatusCode)
	}
	json.Unmarshal(e.admin("POST", "/api/nodes/wh-a/enrollments?ttl=5m", nil, http.StatusCreated), &en)
	e.clock.add(6 * time.Minute)
	body["code"] = en.Code
	if res, _ := e.do(nil, "POST", "/device/enroll", "", body); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("an expired code: %d", res.StatusCode)
	}
	// a key that is not ECDSA P-256
	json.Unmarshal(e.admin("POST", "/api/nodes/wh-a/enrollments", nil, http.StatusCreated), &en)
	if res, _ := e.do(nil, "POST", "/device/enroll", "", map[string]string{"code": en.Code, "key": "AAAA"}); res.StatusCode != http.StatusBadRequest {
		t.Fatalf("a bad key: %d", res.StatusCode)
	}

	// a signature of another key, and an old challenge, do not open a session
	_, out := e.do(nil, "GET", "/device/challenge", "", nil)
	var ch map[string]string
	json.Unmarshal(out, &ch)
	if res, _ := e.do(nil, "POST", "/device/session", "", map[string]string{"device": a.id, "challenge": ch["challenge"], "signature": b.sign(ch["challenge"])}); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("another key's signature: %d", res.StatusCode)
	}
	e.clock.add(ChallengeFor + time.Second)
	if res, _ := e.do(nil, "POST", "/device/session", "", map[string]string{"device": a.id, "challenge": ch["challenge"], "signature": a.sign(ch["challenge"])}); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("an expired challenge: %d", res.StatusCode)
	}

	// a cookie changed, or of a session ended, or of a device removed
	a.login()
	m := e.m()
	u, _ := http.NewRequest("GET", e.srv.URL, nil)
	cookie := a.c.Jar.Cookies(u.URL)[0].Value
	if _, err := m.Session(context.Background(), cookie); err != nil {
		t.Fatal(err)
	}
	payload, mac, _ := strings.Cut(cookie, ".")
	raw, _ := b64.DecodeString(payload)
	forged := b64.EncodeToString(bytes.Replace(raw, []byte("wh-a"), []byte("wh-b"), 1)) + "." + mac
	if _, err := m.Session(context.Background(), forged); err == nil {
		t.Fatal("a forged session")
	}
	e.clock.add(11 * time.Minute)
	if res, _ := e.do(a.c, "GET", "/device/board", "", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("an ended session: %d", res.StatusCode)
	}
	a.login()
	e.admin("DELETE", "/api/devices/"+a.id, nil, http.StatusNoContent)
	if res, _ := e.do(a.c, "GET", "/device/board", "", nil); res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("a removed device: %d", res.StatusCode)
	}
	var ds []Device
	json.Unmarshal(e.admin("GET", "/api/nodes/wh-a/devices", nil, http.StatusOK), &ds)
	if len(ds) != 1 {
		t.Fatalf("devices %+v", ds)
	}
}

// TestEnrollmentURL: an enrollment's URL is at the request's origin, or
// at DeviceOrigin's.
func TestEnrollmentURL(t *testing.T) {
	for _, tc := range []struct {
		origin func(*http.Request) string
		want   string
	}{
		{nil, "http://127.0.0.1:"},
		{func(*http.Request) string { return "https://192.168.1.7:8443" }, "https://192.168.1.7:8443/tablet?code="},
	} {
		e := newEnv(t, func(o *Options) { o.DeviceOrigin = tc.origin })
		e.admin("PUT", "/api/nodes/wh-a", Node{Kind: KindDisplay}, http.StatusOK)
		var en struct{ Code, URL, QR string }
		json.Unmarshal(e.admin("POST", "/api/nodes/wh-a/enrollments", nil, http.StatusCreated), &en)
		if !strings.HasPrefix(en.URL, tc.want) || !strings.HasSuffix(en.URL, "/tablet?code="+en.Code) || en.QR == "" {
			t.Errorf("url %q, want %s…", en.URL, tc.want)
		}
	}
}

func TestStore(t *testing.T) {
	ctx := context.Background()
	for _, u := range []string{"", "cfkv://x", "file:///tmp/x", "bbolt:" + filepath.Join(t.TempDir(), "x.db") + "\x00"} {
		if st, err := OpenStore(ctx, u, ""); err == nil {
			st.Close()
			t.Errorf("%q opened", u)
		}
	}
	if _, err := OpenStore(ctx, "bbolt:"+filepath.Join(t.TempDir(), "x.db"), "dynamodb://t"); err == nil {
		t.Error("documents in DynamoDB")
	}
	st, err := OpenStore(ctx, "bbolt:"+filepath.Join(t.TempDir(), "x.db"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if st.Shared {
		t.Error("a bbolt file is shared")
	}
	if _, err := New(st, Options{Secret: []byte("short")}); err == nil {
		t.Error("a short secret")
	}
}

func TestDuration(t *testing.T) {
	for in, want := range map[string]time.Duration{"14d": 14 * 24 * time.Hour, "10m": 10 * time.Minute, "90s": 90 * time.Second} {
		d, err := ParseDuration(in)
		if err != nil || d.d() != want {
			t.Errorf("%s: %v %v", in, d, err)
		}
		if d.String() != in && !(in == "90s" && d.String() == "1m30s") {
			t.Errorf("%s written %s", in, d)
		}
	}
	for _, bad := range []string{"", "-1m", "xd", "-2d"} {
		if _, err := ParseDuration(bad); err == nil {
			t.Errorf("%q read", bad)
		}
	}
}

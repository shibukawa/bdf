// Package flowtest is what the tests of the adapters of docflow share: a
// manager over a store of a temporary directory, served over HTTP, with
// a clock of its own, and the far sides an exit sends to, in the
// process (an SMTP server, a web hook).
package flowtest

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/docflow"
)

// AdminToken authorizes the API of an Env.
const AdminToken = "admin-token-for-tests"

// Env is a manager's store served over HTTP; every request makes a
// manager of its own, as a function of a serverless platform would.
type Env struct {
	T    *testing.T
	St   *docflow.Store
	Opts docflow.Options
	Srv  *httptest.Server
	mu   sync.Mutex
	now  time.Time
}

// New makes an Env; edit changes its options.
func New(t *testing.T, edit func(*docflow.Options)) *Env {
	t.Helper()
	st, err := docflow.OpenStore(context.Background(), "bbolt:"+filepath.ToSlash(filepath.Join(t.TempDir(), "flow.db")), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	e := &Env{T: t, St: st, now: time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)}
	e.Opts = docflow.Options{Secret: bytes.Repeat([]byte{7}, 32), AdminToken: AdminToken, Now: e.Now, TempDir: t.TempDir(), NoSystemFonts: true}
	if edit != nil {
		edit(&e.Opts)
	}
	e.Srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		docflow.Handler(e.M()).ServeHTTP(w, r)
	}))
	t.Cleanup(e.Srv.Close)
	return e
}

// Now is the Env's clock.
func (e *Env) Now() time.Time {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.now
}

// Advance moves the clock.
func (e *Env) Advance(d time.Duration) {
	e.mu.Lock()
	e.now = e.now.Add(d)
	e.mu.Unlock()
}

// M is a manager of the store.
func (e *Env) M() *docflow.Manager {
	m, err := docflow.New(e.St, e.Opts)
	if err != nil {
		e.T.Fatal(err)
	}
	return m
}

// Put puts nodes, in order, and returns the key of the last keyed one.
func (e *Env) Put(nodes ...docflow.Node) string {
	e.T.Helper()
	var key string
	for _, n := range nodes {
		k, err := e.M().PutNode(context.Background(), n)
		if err != nil {
			e.T.Fatalf("put %s: %v", n.Name, err)
		}
		if k != "" {
			key = k
		}
	}
	return key
}

// Required checks the settings a kind says are required (Kind.Required,
// as the editor marks them) against its Check: valid is put, and none of
// those settings is left out of it (those whose If its settings fit) but
// the node is refused.
func (e *Env) Required(valid docflow.Node) {
	e.T.Helper()
	ctx := context.Background()
	if _, err := e.M().PutNode(ctx, valid); err != nil {
		e.T.Fatalf("put %s: %v", valid.Name, err)
	}
	raw, err := json.Marshal(valid)
	if err != nil {
		e.T.Fatal(err)
	}
	var settings map[string]any
	json.Unmarshal(raw, &settings)
	var required []docflow.Required
	var fields []docflow.Field
	for _, k := range docflow.KindInfos() {
		if k.Name == valid.Kind {
			required, fields = k.Required, k.Fields
		}
	}
	if len(required) == 0 {
		e.T.Fatalf("%s requires nothing", valid.Kind)
	}
	// the settings of a kind's own fields are under "settings"
	own, _ := settings["settings"].(map[string]any)
	isField := func(key string) bool {
		for _, f := range fields {
			if f.Key == key {
				return true
			}
		}
		return false
	}
	valueOf := func(key string) (any, bool) {
		if isField(key) {
			v, ok := own[key]
			return v, ok
		}
		v, ok := settings[key]
		return v, ok
	}
	checked := 0
	for _, r := range required {
		fits := true
		for k, v := range r.If {
			x, _ := valueOf(k)
			fits = fits && fmt.Sprint(x) == v
		}
		if !fits {
			continue
		}
		if _, ok := valueOf(r.Key); !ok {
			e.T.Errorf("%s: %s requires %q, which the node has no value of", valid.Name, valid.Kind, r.Key)
			continue
		}
		without := map[string]any{}
		for k, v := range settings {
			if k != r.Key {
				without[k] = v
			}
		}
		if isField(r.Key) {
			rest := map[string]any{}
			for k, v := range own {
				if k != r.Key {
					rest[k] = v
				}
			}
			without["settings"] = rest
		}
		b, _ := json.Marshal(without)
		var n docflow.Node
		if err := json.Unmarshal(b, &n); err != nil {
			e.T.Fatal(err)
		}
		if _, err := e.M().PutNode(ctx, n); !errors.Is(err, docflow.ErrInvalid) {
			e.T.Errorf("%s without %s: %v, not refused", valid.Name, r.Key, err)
		}
		checked++
	}
	if checked == 0 {
		e.T.Errorf("%s: no required setting fits the node", valid.Name)
	}
}

// Doc is a .bdf document of pages.
func Doc(t *testing.T, pages int) []byte {
	t.Helper()
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	for i := range pages {
		o := bdf.NewObject()
		o.FillColor(0x000000ff).FillRect(50, 100, float32(50*(i+1)), 40)
		h, _ := d.AddObject(o)
		v.AddPage(595, 842, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	}
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// Submit sends a document to a destination.
func (e *Env) Submit(dest, name string, doc []byte, meta map[string]any) *docflow.Job {
	e.T.Helper()
	j, err := e.M().Submit(context.Background(), dest, docflow.Input{Name: name, Source: "api", Doc: bytes.NewReader(doc), Size: int64(len(doc)), Meta: meta})
	if err != nil {
		e.T.Fatalf("submit to %s: %v", dest, err)
	}
	return j
}

// Tick ticks until a tick sends nothing more (the steps make deliveries
// the next tick sends), five times at most.
func (e *Env) Tick() {
	e.T.Helper()
	for range 5 {
		rep, _ := e.M().Tick(context.Background())
		if rep == nil || rep.Delivered == 0 && rep.Undelivered == 0 {
			return
		}
	}
}

// Job is a job as it is now.
func (e *Env) Job(number int32) *docflow.Job {
	e.T.Helper()
	j, err := e.M().JobOfNumber(context.Background(), number)
	if err != nil {
		e.T.Fatal(err)
	}
	return j
}

// Delivery is the delivery of the job's leg at a node (the first).
func (e *Env) Delivery(j *docflow.Job, node string) *docflow.Delivery {
	e.T.Helper()
	for _, l := range e.Job(j.Number).Legs {
		if l.Node == node && l.Ref != "" {
			d, err := e.M().Delivery(context.Background(), l.Ref)
			if err != nil {
				e.T.Fatal(err)
			}
			return d
		}
	}
	e.T.Fatalf("job %d has no delivery at %s: %+v", j.Number, node, e.Job(j.Number).Legs)
	return nil
}

// Payload reads a payload's bytes.
func (e *Env) Payload(id string) []byte {
	e.T.Helper()
	var b bytes.Buffer
	if err := e.M().WritePayload(context.Background(), id, &b); err != nil {
		e.T.Fatal(err)
	}
	return b.Bytes()
}

// Request is a request a Hook got.
type Request struct {
	Method, Path, Query string
	Header              http.Header
	Body                []byte
}

// Hook is a web service that records what it gets and answers the
// statuses given, in turn (the last one after).
type Hook struct {
	*httptest.Server
	mu       sync.Mutex
	got      []Request
	statuses []int
}

// NewHook starts a Hook.
func NewHook(t *testing.T, statuses ...int) *Hook {
	h := &Hook{statuses: statuses}
	h.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		h.mu.Lock()
		h.got = append(h.got, Request{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery, Header: r.Header.Clone(), Body: b})
		status := http.StatusOK
		if len(h.statuses) > 0 {
			status = h.statuses[0]
			if len(h.statuses) > 1 {
				h.statuses = h.statuses[1:]
			}
		}
		h.mu.Unlock()
		w.WriteHeader(status)
	}))
	t.Cleanup(h.Close)
	return h
}

// Got are the requests the Hook got.
func (h *Hook) Got() []Request {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]Request(nil), h.got...)
}

// Mail is a mail an SMTP server got.
type Mail struct {
	From string
	To   []string
	Auth string
	Data string
}

// SMTP is an SMTP server of one connection at a time, that takes AUTH
// PLAIN and refuses the recipients of Refuse with 550.
type SMTP struct {
	Addr   string
	Refuse string
	mu     sync.Mutex
	mails  []Mail
}

// NewSMTP starts an SMTP server on the loopback.
func NewSMTP(t *testing.T) *SMTP {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	s := &SMTP{Addr: l.Addr().String()}
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			go s.serve(c)
		}
	}()
	return s
}

// Mails are the mails the server got.
func (s *SMTP) Mails() []Mail {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Mail(nil), s.mails...)
}

func (s *SMTP) serve(c net.Conn) {
	defer c.Close()
	tp := textproto.NewConn(c)
	tp.PrintfLine("220 fake ESMTP")
	var m Mail
	for {
		line, err := tp.ReadLine()
		if err != nil {
			return
		}
		cmd, arg, _ := strings.Cut(line, " ")
		switch strings.ToUpper(cmd) {
		case "EHLO", "HELO":
			tp.PrintfLine("250-fake")
			tp.PrintfLine("250 AUTH PLAIN")
		case "AUTH":
			if p, ok := strings.CutPrefix(arg, "PLAIN "); ok {
				b, _ := base64.StdEncoding.DecodeString(p)
				m.Auth = string(b)
			}
			tp.PrintfLine("235 ok")
		case "MAIL":
			m.From = strings.Trim(strings.TrimPrefix(arg, "FROM:"), "<>")
			tp.PrintfLine("250 ok")
		case "RCPT":
			to := strings.Trim(strings.TrimPrefix(arg, "TO:"), "<>")
			if s.Refuse != "" && to == s.Refuse {
				tp.PrintfLine("550 no such user")
				continue
			}
			m.To = append(m.To, to)
			tp.PrintfLine("250 ok")
		case "DATA":
			tp.PrintfLine("354 go on")
			b, err := tp.ReadDotBytes()
			if err != nil {
				return
			}
			m.Data = string(b)
			s.mu.Lock()
			s.mails = append(s.mails, m)
			s.mu.Unlock()
			m = Mail{}
			tp.PrintfLine("250 queued")
		case "QUIT":
			tp.PrintfLine("221 bye")
			return
		default:
			tp.PrintfLine("250 ok")
		}
	}
}

// Header is the value of a header of a mail.
func Header(data, name string) string {
	r := textproto.NewReader(bufio.NewReader(strings.NewReader(data)))
	h, err := r.ReadMIMEHeader()
	if err != nil && len(h) == 0 {
		return fmt.Sprint(err)
	}
	return h.Get(name)
}

package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/server"
	"github.com/shibukawa/bdf/print/server/client"
	"github.com/shibukawa/bdf/print/spool"
)

// fill sets every field of a value to something not zero.
func fill(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			v.Set(reflect.ValueOf(time.Date(2026, 10, 20, 9, 30, 0, 0, time.UTC)))
			return
		}
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).IsExported() {
				fill(v.Field(i))
			}
		}
	case reflect.Pointer:
		v.Set(reflect.New(v.Type().Elem()))
		fill(v.Elem())
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		fill(v.Index(0))
	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		k, e := reflect.New(v.Type().Key()).Elem(), reflect.New(v.Type().Elem()).Elem()
		fill(k)
		fill(e)
		v.SetMapIndex(k, e)
	case reflect.String:
		v.SetString("x")
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(3)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(3)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(3)
	}
}

// sameShape checks that what the server writes of a value of the spool
// reads into the client's type and writes back the same.
func sameShape(t *testing.T, server, client any) {
	t.Helper()
	fill(reflect.ValueOf(server).Elem())
	b, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(client); err != nil {
		t.Fatalf("%T: %v\n%s", client, err, b)
	}
	c, _ := json.Marshal(client)
	var want, got any
	json.Unmarshal(b, &want)
	json.Unmarshal(c, &got)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("%T:\nserver %s\nclient %s", client, b, c)
	}
}

func TestShapes(t *testing.T) {
	sameShape(t, &spool.Job{}, &client.Job{})
	sameShape(t, &spool.QueueInfo{}, &client.Queue{})
	sameShape(t, &spool.PrinterInfo{}, &client.Printer{})
}

func testDoc(t *testing.T) []byte {
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

// newServer serves a server of a fake printer and returns a client of it.
func newServer(t *testing.T, token string) (*client.Client, *server.Server, *fakeprinter.Printer) {
	t.Helper()
	p := fakeprinter.New()
	cfg, err := server.ParseConfig([]byte(fmt.Sprintf(`
printers: { mfp1: { url: "%s" } }
queues: { office: { printers: [mfp1] }, spare: { printers: [mfp1] } }
tokens: [ { name: ops, value: %s, roles: [submit, manage, admin] } ]
`, fakeprinter.Start(t, p), token)))
	if err != nil {
		t.Fatal(err)
	}
	cfg.Spool = filepath.Join(t.TempDir(), "spool.db")
	cfg.TempDir = t.TempDir()
	srv, err := server.New(context.Background(), cfg, server.Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { srv.Close() })
	hs := httptest.NewServer(srv.Handler())
	t.Cleanup(hs.Close)
	return &client.Client{URL: hs.URL, Token: token}, srv, p
}

func TestClient(t *testing.T) {
	const token = "ops-0123456789abcdef"
	c, srv, p := newServer(t, token)
	ctx := context.Background()

	if _, err := (&client.Client{URL: c.URL}).Queues(ctx); !client.IsCode(err, "unauthorized") {
		t.Fatalf("no token: %v", err)
	}
	if _, err := c.Job(ctx, 99); !client.IsCode(err, "not-found") {
		t.Fatalf("no job: %v", err)
	}

	// held in a paused queue, changed, then released
	if _, err := c.QueueOp(ctx, "office", "pause"); err != nil {
		t.Fatal(err)
	}
	j, err := c.Submit(ctx, client.Submission{Queue: "office", Name: "invoice.bdf", Document: bytes.NewReader(testDoc(t)),
		Ticket: print.Ticket{Copies: 2, Sides: "two-sided-long-edge"}, Options: &client.JobOptions{Format: "pdf"}, Priority: 70, HoldUntil: "indefinite"})
	if err != nil {
		t.Fatal(err)
	}
	if j.State != "pending-held" || j.Name != "invoice.bdf" || j.Priority != 70 || j.Options.Format != "pdf" || j.User != "ops" ||
		!strings.Contains(string(j.Ticket), `"copies":2`) {
		t.Fatalf("submitted %+v", j)
	}
	if j, err = c.SetPriority(ctx, j.ID, 80); err != nil || j.Priority != 80 {
		t.Fatalf("priority: %v %+v", err, j)
	}
	if j, err = c.Move(ctx, j.ID, "spare"); err != nil || j.Queue != "spare" {
		t.Fatalf("move: %v %+v", err, j)
	}
	if j, err = c.Release(ctx, j.ID, ""); err != nil || j.State != "pending" {
		t.Fatalf("release: %v %+v", err, j)
	}

	// waited for until printed
	srv.Spool().Start()
	var mu sync.Mutex
	var seen []string
	done, err := c.Wait(ctx, j.ID, func(j *client.Job) {
		mu.Lock()
		seen = append(seen, j.State)
		mu.Unlock()
	})
	if err != nil || done.State != "completed" || done.Printer == nil || done.Printer.Name != "mfp1" {
		t.Fatalf("wait: %v %+v", err, done)
	}
	if seen[0] != "pending" || seen[len(seen)-1] != "completed" {
		t.Fatalf("changes %v", seen)
	}
	if len(p.Jobs()) != 1 {
		t.Fatalf("the printer has %d jobs", len(p.Jobs()))
	}

	// canceled, then sent again
	if _, err := c.QueueOp(ctx, "office", "resume"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.QueueOp(ctx, "office", "pause"); err != nil {
		t.Fatal(err)
	}
	k, err := c.Submit(ctx, client.Submission{Queue: "office", Document: bytes.NewReader(testDoc(t))})
	if err != nil {
		t.Fatal(err)
	}
	if k, err = c.Cancel(ctx, k.ID); err != nil || k.State != "canceled" {
		t.Fatalf("cancel: %v %+v", err, k)
	}
	if _, err = c.Cancel(ctx, k.ID); !client.IsCode(err, "not-possible") {
		t.Fatalf("cancel twice: %v", err)
	}
	if k, err = c.Retry(ctx, k.ID); err != nil || k.State != "pending" {
		t.Fatalf("retry: %v %+v", err, k)
	}

	list, err := c.Jobs(ctx, client.JobFilter{State: "all"})
	if err != nil || len(list.Jobs) != 2 {
		t.Fatalf("jobs: %v %+v", err, list)
	}
	qs, err := c.Queues(ctx)
	if err != nil || len(qs) != 2 || !qs[0].Paused || qs[0].Jobs["pending"] != 1 {
		t.Fatalf("queues: %v %+v", err, qs)
	}
	pr, err := c.Printer(ctx, "mfp1")
	if err != nil || !pr.Up || pr.Capabilities == nil || len(pr.Capabilities.Formats) == 0 {
		t.Fatalf("printer: %v %+v", err, pr)
	}
	var acc bytes.Buffer
	if err := c.Accounting(ctx, time.Time{}, time.Time{}, "csv", &acc); err != nil || strings.Count(acc.String(), "\n") != 2 {
		t.Fatalf("accounting: %v %s", err, acc.String())
	}
	var backup bytes.Buffer
	if err := c.Backup(ctx, &backup); err != nil || backup.Len() == 0 {
		t.Fatalf("backup: %v %d", err, backup.Len())
	}

	// the events of a queue
	ectx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		time.Sleep(200 * time.Millisecond)
		c.QueueOp(ctx, "office", "resume")
	}()
	err = c.Events(ectx, 0, func(ev client.Event) error {
		if ev.Queue != nil && ev.Queue.Name == "office" && !ev.Queue.Paused {
			return fmt.Errorf("resumed")
		}
		return nil
	})
	if err == nil || err.Error() != "resumed" {
		t.Fatalf("events: %v", err)
	}
}

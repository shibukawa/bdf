package ipp_test

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	dipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/printconv"
)

// What Bonjour finds is listed for the editor: every IPP printer for a
// printer (AirPrint, IPP Everywhere and Mopria printers, told apart), the
// queues of CUPS servers for a CUPS server (sink) asked with no URL; the
// CUPS servers of the manager itself are left out, and a serverless
// deployment has no local network to look on.
func TestDiscoverBonjour(t *testing.T) {
	e := flowtest.New(t, nil)
	ctx := context.Background()
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office", To: []string{"keep"}})
	office, err := e.M().Node(ctx, "office")
	if err != nil {
		t.Fatal(err)
	}
	dipp.SetBrowse(t, func(context.Context) ([]ipp.Service, error) {
		return []ipp.Service{
			{Name: "Canon MF", Host: "canon.local", Port: 631, URL: "ipp://192.0.2.10:631/ipp/print", TXT: map[string]string{"rp": "ipp/print", "ty": "Canon MF644C", "URF": "V1.4", "note": "2F",
				"pdl": "application/pdf,image/urf,image/pwg-raster", "Mopria-Certified": "2.0"}},
			{Name: "Laser @ mac", Host: "mac.local", Port: 631, URL: "ipps://192.0.2.11:631/printers/Laser", Secure: true, TXT: map[string]string{"rp": "printers/Laser", "ty": "HP LaserJet"}},
			{Name: "Inkjet", Host: "inkjet.local", Port: 631, URL: "ipp://192.0.2.13:631/ipp/print", TXT: map[string]string{"rp": "ipp/print", "pdl": "application/PCLm,image/pwg-raster", "mopria-certified": "1.3"}},
			{Name: "Office @ here", Host: "here-docflow.local", Port: office.Port, URL: "ipp://192.0.2.12:" + strconv.Itoa(office.Port) + "/printers/office",
				TXT: map[string]string{"rp": "printers/office", "UUID": strings.TrimPrefix(dipp.PrinterUUID("office"), "urn:uuid:")}},
		}, nil
	})
	found := func(k string, in docflow.Node) []dipp.Found {
		t.Helper()
		got, err := kind(t, k).Discover(ctx, e.M(), in)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(got)
		var fs []dipp.Found
		json.Unmarshal(b, &fs)
		return fs
	}
	printers := found(dipp.KindPrinter, docflow.Node{})
	if len(printers) != 3 {
		t.Fatalf("printers %+v", printers)
	}
	// what each speaks: AirPrint (URF), IPP Everywhere (PWG raster),
	// Mopria (certified, in a key of any case)
	if c := printers[0]; c.Name != "Canon MF" || c.URL != "ipp://192.0.2.10:631/ipp/print" || c.Model != "Canon MF644C" || c.Note != "2F" || !c.AirPrint || !c.Everywhere ||
		!c.Mopria || c.Queue || c.Secure {
		t.Errorf("canon %+v", c)
	}
	if i := printers[2]; i.Name != "Inkjet" || i.AirPrint || !i.Everywhere || !i.Mopria {
		t.Errorf("inkjet %+v", i)
	}
	if l := printers[1]; l.Name != "Laser @ mac" || !l.Queue || !l.Secure || l.AirPrint || l.Everywhere || l.Mopria {
		t.Errorf("laser %+v", l)
	}
	if qs := found(dipp.KindCupsOut, docflow.Node{}); len(qs) != 1 || qs[0].Name != "Laser @ mac" || qs[0].URL != "ipps://192.0.2.11:631/printers/Laser" {
		t.Errorf("queues %+v", qs)
	}

	s := flowtest.New(t, func(o *docflow.Options) { o.NoPorts = true })
	if _, err := kind(t, dipp.KindPrinter).Discover(ctx, s.M(), docflow.Node{}); !errors.Is(err, docflow.ErrUnsupported) {
		t.Errorf("serverless: %v", err)
	}
}

// A printer picked from Bonjour keeps its instance name; when its address
// takes no connection, it is looked for again by the name and the job
// printed where it is now, and the next job goes there first. Picked by
// hand, it is not looked for.
func TestBonjourMoved(t *testing.T) {
	e := flowtest.New(t, func(o *docflow.Options) { o.Converter = printconv.Converter{NoSystemFonts: true} })
	ctx := context.Background()
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office", To: []string{"keep"}})
	// the printer now: the CUPS server (source)'s port, as cmd/docflow serves it
	cups := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		docflow.PortHandler(e.M(), "office").ServeHTTP(w, r)
	}))
	defer cups.Close()
	now := strings.Replace(cups.URL, "http://", "ipp://", 1) + "/printers/office"
	// where it was: an address nothing listens on
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	was := "ipp://" + ln.Addr().String() + "/printers/office"
	ln.Close()
	asked := 0
	dipp.SetBrowse(t, func(context.Context) ([]ipp.Service, error) {
		asked++
		return []ipp.Service{{Name: "Office @ mac", Host: "mac.local", URL: now, TXT: map[string]string{"rp": "printers/office"}}}, nil
	})

	for _, bad := range []docflow.Node{
		{Name: "p", Kind: dipp.KindPrinter, URL: "socket://192.0.2.1:9100", Bonjour: "Office @ mac"},
		{Name: "p", Kind: dipp.KindPrinter, URL: was, Bonjour: strings.Repeat("x", 64)},
	} {
		if _, err := e.M().PutNode(ctx, bad); !errors.Is(err, docflow.ErrInvalid) {
			t.Errorf("%+v: %v", bad, err)
		}
	}
	e.Put(docflow.Node{Name: "picked", Kind: dipp.KindPrinter, URL: was, Bonjour: " Office @ mac "},
		docflow.Node{Name: "queue", Kind: dipp.KindCupsOut, URL: was, Bonjour: "office @ MAC"},
		docflow.Node{Name: "typed", Kind: dipp.KindPrinter, URL: was})
	if n, _ := e.M().Node(ctx, "picked"); n.Bonjour != "Office @ mac" || n.URL != was {
		t.Fatalf("picked %+v", n)
	}
	for i, dest := range []string{"picked", "queue", "picked"} {
		j := e.Submit(dest, "statement", flowtest.Doc(t, 1), nil)
		e.Tick()
		if d := e.Delivery(j, dest); d.State != docflow.DeliveryDone {
			t.Fatalf("%d %s: delivery %+v", i, dest, d)
		}
	}
	// found once for each node: the second job of picked went there first
	if asked != 2 {
		t.Errorf("Bonjour asked %d times", asked)
	}
	j := e.Submit("typed", "statement", flowtest.Doc(t, 1), nil)
	e.Tick()
	if d := e.Delivery(j, "typed"); d.State == docflow.DeliveryDone || asked != 2 {
		t.Fatalf("typed by hand: delivery %+v, Bonjour asked %d times", d, asked)
	}
}

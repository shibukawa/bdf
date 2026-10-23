package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
	docflowipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/ipp"
	"golang.org/x/net/dns/dnsmessage"
)

// The process listens on the ports of the CUPS server nodes, and stops
// when a node goes; an uploader has none (it is at the API).
func TestPorts(t *testing.T) {
	ctx := context.Background()
	a, err := Open(ctx, DocflowConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1}, OIDCConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	free := func() int {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		port := ln.Addr().(*net.TCPAddr).Port
		ln.Close()
		return port
	}
	cupsPort, apiPort := free(), free()
	if _, err := a.m.PutNode(ctx, docflow.Node{Name: "wh-a", Kind: docflow.KindDisplay}); err != nil {
		t.Fatal(err)
	}
	if _, err := a.m.PutNode(ctx, docflow.Node{Name: "cups1", Kind: docflowipp.KindCups, To: []string{"wh-a"}, Port: cupsPort}); err != nil {
		t.Fatal(err)
	}
	key, err := a.m.PutNode(ctx, docflow.Node{Name: "api", Kind: docflow.KindUpload, To: []string{"wh-a"}, Port: apiPort})
	if err != nil {
		t.Fatal(err)
	}
	a.ports.reconcile(ctx)
	if open := a.ports.open(); len(open) != 1 || open[0] != cupsPort {
		t.Fatalf("ports open %v", open)
	}
	// a CUPS client prints to the CUPS server's one printer, the node, at
	// /printers/<node> of its port; the job goes on to the display
	c := &ipp.Client{URL: fmt.Sprintf("http://127.0.0.1:%d/printers/cups1", cupsPort), User: "lp"}
	js, err := c.PrintJob(ctx, "from cups", "application/x-bdf", nil, bytes.NewReader(doc(t)))
	if err != nil || js.State != ipp.JobCompleted {
		t.Fatalf("print through the CUPS port: %+v %v", js, err)
	}
	if j, err := a.m.JobOfNumber(ctx, js.ID); err != nil || j.Dest != "cups1" || len(j.Legs) != 2 || j.Legs[1].Node != "wh-a" {
		t.Fatalf("the job through the CUPS port: %+v %v", j, err)
	}

	// an application posts to the uploader at the API of the manager's
	// own server
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()
	req, _ := http.NewRequest("POST", srv.URL+docflow.NodePath(docflow.KindUpload, "api")+"?name=x", bytes.NewReader(doc(t)))
	req.Header.Set("Authorization", "Bearer "+key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("post to the uploader: %d", res.StatusCode)
	}
	if res, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/", apiPort), "application/octet-stream", nil); err == nil {
		res.Body.Close()
		t.Fatal("something listens on the uploader's port")
	}
	// the CUPS server goes: the listener closes
	if err := a.m.DeleteNode(ctx, "cups1"); err != nil {
		t.Fatal(err)
	}
	a.ports.reconcile(ctx)
	if open := a.ports.open(); len(open) != 0 {
		t.Fatalf("ports open %v", open)
	}
	if res, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/", cupsPort), "application/ipp", nil); err == nil {
		res.Body.Close()
		t.Fatal("the CUPS server's port is still open")
	}
}

// A CUPS server node that says so is advertised with DNS-SD while its
// port is open, as "<title> @ <computer>" on a host name of the
// process's own; it stops when the node no longer says so. The
// responder is asked on the loopback, at a port of the test's own.
func TestPortsAdvertise(t *testing.T) {
	ctx := context.Background()
	a, err := Open(ctx, DocflowConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1}, OIDCConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	cupsPort := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	uc, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	mdnsPort := uc.LocalAddr().(*net.UDPAddr).Port
	uc.Close()
	a.ports.mdns = ipp.ResponderOptions{Port: mdnsPort, Interfaces: []net.Interface{}, Loopback: true}
	put := func(dnssd bool) {
		t.Helper()
		if _, err := a.m.PutNode(ctx, docflow.Node{Name: "wh-a", Kind: docflow.KindDisplay}); err != nil {
			t.Fatal(err)
		}
		if _, err := a.m.PutNode(ctx, docflow.Node{Name: "office", Kind: docflowipp.KindCups, Title: "Office", To: []string{"wh-a"}, Port: cupsPort, DNSSD: dnssd}); err != nil {
			t.Fatal(err)
		}
		a.ports.reconcile(ctx)
	}
	// ask asks for the IPP printers; "" when nothing answers in time
	ask := func(wait time.Duration) string {
		t.Helper()
		c, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		b := dnsmessage.NewBuilder(nil, dnsmessage.Header{ID: 9})
		b.StartQuestions()
		b.Question(dnsmessage.Question{Name: dnsmessage.MustNewName("_ipp._tcp.local."), Type: dnsmessage.TypePTR, Class: dnsmessage.ClassINET})
		q, _ := b.Finish()
		c.WriteTo(q, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: mdnsPort})
		c.SetReadDeadline(time.Now().Add(wait))
		buf := make([]byte, 9000)
		n, _, err := c.ReadFrom(buf)
		if err != nil {
			return ""
		}
		var p dnsmessage.Parser
		if _, err := p.Start(buf[:n]); err != nil {
			t.Fatal(err)
		}
		p.SkipAllQuestions()
		ans, _ := p.AllAnswers()
		if len(ans) != 1 {
			t.Fatalf("answers %v", ans)
		}
		return ans[0].Body.(*dnsmessage.PTRResource).PTR.String()
	}
	put(true)
	want := "Office @ " + machineName() + "._ipp._tcp.local."
	got := ""
	for i := 0; i < 50 && got == ""; i++ {
		// the names are probed first (a second, more on a busy machine)
		got = ask(200 * time.Millisecond)
	}
	if got != want {
		t.Fatalf("advertised %q, want %q", got, want)
	}
	// reconciled again unchanged, it goes on as it is
	a.ports.reconcile(ctx)
	if got := ask(time.Second); got != want {
		t.Fatalf("advertised %q after a reconcile", got)
	}
	// the node no longer says so: goodbye
	put(false)
	if got := ask(500 * time.Millisecond); got != "" {
		t.Fatalf("still advertised: %q", got)
	}
}

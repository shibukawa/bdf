package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/manager"
)

// The process listens on the ports of the entry nodes, and stops when a
// node loses its port.
func TestPorts(t *testing.T) {
	ctx := context.Background()
	a, err := Open(ctx, ManagerConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1}, OIDCConfig{})
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
	if _, err := a.m.PutNode(ctx, manager.Node{Name: "wh-a", Kind: manager.KindDisplay}); err != nil {
		t.Fatal(err)
	}
	if _, err := a.m.PutNode(ctx, manager.Node{Name: "cups1", Kind: manager.KindCups, Port: cupsPort}); err != nil {
		t.Fatal(err)
	}
	key, err := a.m.PutNode(ctx, manager.Node{Name: "api", Kind: manager.KindEndpoint, To: []string{"wh-a"}, Port: apiPort})
	if err != nil {
		t.Fatal(err)
	}
	a.ports.reconcile(ctx)
	if open := a.ports.open(); len(open) != 2 {
		t.Fatalf("ports open %v", open)
	}
	// a CUPS client prints to a node at /printers/<node> of the CUPS server's port
	c := &ipp.Client{URL: fmt.Sprintf("http://127.0.0.1:%d/printers/wh-a", cupsPort), User: "lp"}
	js, err := c.PrintJob(ctx, "from cups", "application/x-bdf", nil, bytes.NewReader(doc(t)))
	if err != nil || js.State != ipp.JobCompleted {
		t.Fatalf("print through the CUPS port: %+v %v", js, err)
	}
	// an application posts to the endpoint's port
	req, _ := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:%d/?name=x", apiPort), bytes.NewReader(doc(t)))
	req.Header.Set("Authorization", "Bearer "+key)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("post to the endpoint's port: %d", res.StatusCode)
	}
	// the endpoint loses its port: the listener closes
	if _, err := a.m.PutNode(ctx, manager.Node{Name: "api", Kind: manager.KindEndpoint, To: []string{"wh-a"}}); err != nil {
		t.Fatal(err)
	}
	a.ports.reconcile(ctx)
	if open := a.ports.open(); len(open) != 1 || open[0] != cupsPort {
		t.Fatalf("ports open %v", open)
	}
	if res, err := http.Post(fmt.Sprintf("http://127.0.0.1:%d/", apiPort), "application/octet-stream", nil); err == nil {
		res.Body.Close()
		t.Fatal("the endpoint's port is still open")
	}
}

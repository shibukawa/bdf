package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestPrintdCheck(t *testing.T) {
	_, u := fakePrinter(t)
	dir := t.TempDir()
	cfg := filepath.Join(dir, "printd.yaml")
	os.WriteFile(cfg, []byte(fmt.Sprintf("printers:\n  mfp1: { url: %q }\n  gone: { url: \"ipp://127.0.0.1:1/ipp/print\" }\nqueues:\n  office: { printers: [mfp1, gone] }\n", u)), 0o600)
	code, out, errOut := runOut(t, "printd", "-config", cfg, "-check")
	if code != 1 || !strings.Contains(out, "mfp1") || !strings.Contains(out, "idle") || !strings.Contains(out, "BDF Test Printer") ||
		!strings.Contains(out, "gone") || !strings.Contains(out, "error:") {
		t.Fatalf("exit %d\n%s\n%s", code, out, errOut)
	}
	if _, err := os.Stat(filepath.Join(dir, "spool.db")); !os.IsNotExist(err) {
		t.Fatal("-check opened the spool")
	}
	os.WriteFile(cfg, []byte("printers: {}\n"), 0o600)
	if code, _, errOut := runOut(t, "printd", "-config", cfg, "-check"); code != 1 || !strings.Contains(errOut, "no printers") {
		t.Fatalf("a bad configuration: exit %d: %s", code, errOut)
	}
}

func TestPrintdRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test stops the server with SIGINT")
	}
	p, u := fakePrinter(t)
	dir := t.TempDir()
	cfg := filepath.Join(dir, "printd.yaml")
	os.WriteFile(cfg, []byte(fmt.Sprintf("printers:\n  mfp1: { url: %q }\nqueues:\n  office: { printers: [mfp1] }\n", u)), 0o600)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	cmd := exec.Command(os.Args[0], "printd", "-config", cfg, "-listen", addr, "-log", "text")
	cmd.Env = append(os.Environ(), "BDF_TEST_MAIN=1")
	var log bytes.Buffer
	cmd.Stderr = &log
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()
	defer func() {
		cmd.Process.Kill()
		<-exited
	}()
	base := "http://" + addr
	wait := func(what string, ok func() bool) {
		t.Helper()
		for deadline := time.Now().Add(15 * time.Second); !ok(); time.Sleep(50 * time.Millisecond) {
			if time.Now().After(deadline) {
				t.Fatalf("timed out waiting for %s; log:\n%s", what, log.String())
			}
		}
	}
	wait("the server", func() bool {
		res, err := http.Get(base + "/healthz")
		if err == nil {
			res.Body.Close()
		}
		return err == nil && res.StatusCode == 200
	})

	// a second server of the same spool is refused
	if code, _, errOut := runOut(t, "printd", "-config", cfg, "-listen", "127.0.0.1:0"); code != 1 || !strings.Contains(errOut, "spool") {
		t.Fatalf("a second server: exit %d: %s", code, errOut)
	}

	doc, err := os.ReadFile(filepath.Join("..", "..", "testdata", "form", "rows.bdf"))
	if err != nil {
		t.Fatal(err)
	}
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, _ := mw.CreateFormFile("document", "rows.bdf")
	fw.Write(doc)
	mw.WriteField("ticket", `{"copies": 2}`)
	mw.Close()
	req, _ := http.NewRequest("POST", base+"/api/v1/queues/office/jobs", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("X-BDF-User", "sato")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var job struct {
		ID    int32  `json:"id"`
		State string `json:"state"`
		User  string `json:"user"`
	}
	json.NewDecoder(res.Body).Decode(&job)
	res.Body.Close()
	if res.StatusCode != 201 || job.User != "sato" {
		t.Fatalf("submit: %d %+v", res.StatusCode, job)
	}
	wait("the job printed", func() bool {
		res, err := http.Get(fmt.Sprintf("%s/api/v1/jobs/%d", base, job.ID))
		if err != nil {
			return false
		}
		defer res.Body.Close()
		json.NewDecoder(res.Body).Decode(&job)
		return job.State == "completed"
	})
	if len(p.Jobs()) != 1 {
		t.Fatalf("the printer has %d jobs", len(p.Jobs()))
	}

	cmd.Process.Signal(os.Interrupt)
	select {
	case err := <-exited:
		exited <- err
		if err != nil {
			t.Fatalf("exit: %v; log:\n%s", err, log.String())
		}
	case <-time.After(15 * time.Second):
		t.Fatalf("printd did not stop; log:\n%s", log.String())
	}
	for _, want := range []string{"msg=serving", "msg=\"job submitted\"", "msg=stopped"} {
		if !strings.Contains(log.String(), want) {
			t.Errorf("the log lacks %s:\n%s", want, log.String())
		}
	}

	code, out, errOut := runOut(t, "printd", "-config", cfg, "-compact")
	if code != 0 || !strings.Contains(out, "spool.db:") {
		t.Fatalf("-compact: exit %d: %s %s", code, out, errOut)
	}
}

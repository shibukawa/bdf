package ipp_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	_ "github.com/shibukawa/bdf/converter/image"
	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	dipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/printconv"
)

var failLine = regexp.MustCompile(`(?m)^.*\[FAIL\]\n(?:\s{8}.*\n)*`)

// A CUPS server node's own port passes CUPS's ipptool checks of IPP/1.1
// and IPP/2.0 (with Create-Job and Send-Document), and those of the
// attributes an IPP Everywhere printer has (PWG 5100.14), but
// page-ranges and overrides: a document goes on whole, and clients pick
// the pages they send. ipptool prints a PDF of its own there, which is
// made .bdf. Skipped where CUPS's ipptool is not (the CI's print job
// installs it, macOS has it).
func TestIPPTool(t *testing.T) {
	tool, err := exec.LookPath("ipptool")
	if err != nil {
		t.Skip("no ipptool")
	}
	var dir string
	for _, d := range []string{"/usr/share/cups/ipptool", "/usr/local/share/cups/ipptool", "/opt/homebrew/share/cups/ipptool"} {
		if _, err := os.Stat(filepath.Join(d, "ipp-2.0.test")); err == nil {
			dir = d
		}
	}
	if dir == "" {
		t.Skip("no tests of ipptool")
	}
	e := flowtest.New(t, func(o *docflow.Options) { o.Converter = printconv.Converter{NoSystemFonts: true} })
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office", To: []string{"keep"}})
	cups := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		docflow.PortHandler(e.M(), "office").ServeHTTP(w, r)
	}))
	defer cups.Close()
	// ipptool reaches the printer by name, as clients do
	uri := "ipp://" + strings.Replace(strings.TrimPrefix(cups.URL, "http://"), "127.0.0.1", "localhost", 1) + "/printers/office"
	everywhere, _ := filepath.Abs("../../../test/print/ipp-everywhere-server.test")
	doc := filepath.Join(dir, "onepage-a4.pdf")
	for _, test := range []string{filepath.Join(dir, "ipp-2.0.test"), everywhere} {
		out, err := exec.Command(tool, "-t", "-V", "2.0", "-f", doc, "-d", "NO_PAGE_RANGES=1", uri, test).CombinedOutput()
		if err != nil || strings.Contains(string(out), "[FAIL]") {
			if fails := failLine.FindAllString(string(out), -1); len(fails) > 0 {
				t.Errorf("ipptool %s:\n%s", filepath.Base(test), strings.Join(fails, ""))
			} else {
				t.Errorf("ipptool %s: %v\n%s", filepath.Base(test), err, out)
			}
			continue
		}
		if !strings.Contains(string(out), "[PASS]") {
			t.Errorf("ipptool %s: nothing passed\n%s", filepath.Base(test), out)
		}
	}
	// what ipptool printed came in as jobs, made .bdf
	jobs, err := e.M().JobsOf(t.Context(), "office", 0)
	if err != nil || len(jobs) == 0 {
		t.Fatalf("jobs %d %v", len(jobs), err)
	}
	for _, j := range jobs {
		if j.Pages == 0 || j.Format != "pdf" && j.Format != "image" {
			t.Errorf("job %d: %s, %d pages", j.Number, j.Format, j.Pages)
		}
	}
}

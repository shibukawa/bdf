package main

import (
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/ipp/ipptest"
)

// runOut runs the bdf command and returns its exit status, output and
// error output.
func runOut(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "BDF_TEST_MAIN=1")
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	if exit, ok := err.(*exec.ExitError); ok {
		return exit.ExitCode(), stdout.String(), stderr.String()
	}
	if err != nil {
		t.Fatal(err)
	}
	return 0, stdout.String(), stderr.String()
}

// fakePrinter serves a printer that is not there and returns its URL.
func fakePrinter(t *testing.T) (*ipptest.Printer, string) {
	t.Helper()
	p := ipptest.New()
	srv := httptest.NewServer(p)
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	return p, "ipp://" + u.Host + "/ipp/print"
}

func TestPrintersURL(t *testing.T) {
	_, u := fakePrinter(t)
	code, out, errOut := runOut(t, "printers", "-url", u)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{"name: ipptest", "model: BDF Test Printer", "state: idle", "formats: application/pdf, image/pwg-raster",
		"sides: one-sided, two-sided-long-edge", "resolution: 300x300dpi, 600x600dpi", "copies: 1-99"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in\n%s", want, out)
		}
	}
	code, out, _ = runOut(t, "printers", "-url", u, "-json")
	if code != 0 || !strings.Contains(out, `"MakeModel": "BDF Test Printer"`) {
		t.Fatalf("-json: exit %d\n%s", code, out)
	}
	code, out, _ = runOut(t, "printers", "-cups", u)
	if code != 0 || !strings.HasPrefix(out, "ipptest\tipp://") {
		t.Fatalf("-cups: exit %d\n%s", code, out)
	}
	if code, _, errOut := runOut(t, "printers", "-url", "ipp://127.0.0.1:1/ipp/print"); code != 1 || !strings.Contains(errOut, "bdf: ipp:") {
		t.Fatalf("no printer: exit %d: %s", code, errOut)
	}
}

func TestPrint(t *testing.T) {
	p, u := fakePrinter(t)
	doc := filepath.Join("..", "..", "testdata", "form", "rows.bdf")
	code, _, errOut := runOut(t, append([]string{"print", "-printer", u, "-copies", "2", "-sides", "long", "-media", "a4", "-mono"}, append(testFontFlags, doc)...)...)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if !strings.Contains(errOut, "job 1 completed") || !strings.Contains(errOut, "as PDF on iso_a4_210x297mm") {
		t.Errorf("output: %s", errOut)
	}
	j := p.Jobs()[0]
	if j.Name != "rows.bdf" || !strings.HasPrefix(string(j.Document(0)), "%PDF-") {
		t.Errorf("job %q: %q", j.Name, string(j.Document(0))[:8])
	}

	// a printer of PWG Raster only, and an input converted first
	p.Formats = []string{"image/pwg-raster"}
	tiff, _ := filepath.Glob(filepath.Join("..", "..", "converter", "tiff", "testdata", "*.tif"))
	if len(tiff) == 0 {
		t.Fatal("no TIFF in converter/tiff/testdata")
	}
	code, _, errOut = runOut(t, append([]string{"print", "-printer", u, "-q", "-pages", "1"}, append(testFontFlags, tiff[0])...)...)
	if code != 0 || errOut != "" {
		t.Fatalf("a TIFF converted: exit %d: %s", code, errOut)
	}
	// a sheet has no pages to print without a print view
	csv := filepath.Join(t.TempDir(), "list.csv")
	os.WriteFile(csv, []byte("name,qty\napple,3\npear,5\n"), 0o644)
	if code, _, errOut := runOut(t, "print", "-printer", u, csv); code != 1 || !strings.Contains(errOut, "no pages") {
		t.Errorf("a sheet: exit %d: %s", code, errOut)
	}
	if j := p.Jobs()[1]; j.Format != "image/pwg-raster" || !strings.HasPrefix(string(j.Document(0)), "RaS2") {
		t.Errorf("raster job: %s", j.Format)
	}

	// what a printer refuses, and the exit statuses
	p.Formats = []string{"application/pdf"}
	code, _, errOut = runOut(t, append([]string{"print", "-printer", u, "-strict", "-source", "tray-9"}, append(testFontFlags, doc)...)...)
	if code != exitUnsupported || !strings.Contains(errOut, "media-source=tray-9") {
		t.Errorf("strict: exit %d: %s", code, errOut)
	}
	p.State, p.Reasons = 5, []string{"media-empty-error"}
	code, _, errOut = runOut(t, append([]string{"print", "-printer", u}, append(testFontFlags, doc)...)...)
	if code != exitStopped || !strings.Contains(errOut, "media-empty-error") {
		t.Errorf("stopped: exit %d: %s", code, errOut)
	}
	code, _, errOut = runOut(t, "print", "-printer", "ipp://127.0.0.1:1/ipp/print", doc)
	if code != exitConnection {
		t.Errorf("unreachable: exit %d: %s", code, errOut)
	}
	if code, _, _ := runOut(t, "print", doc); code != 2 {
		t.Errorf("no printer: exit %d", code)
	}
	if code, _, _ := runOut(t, "print", "-printer", u, "-sides", "both", doc); code != 2 {
		t.Errorf("bad sides: exit %d", code)
	}
}

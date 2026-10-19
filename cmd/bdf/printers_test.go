package main

import (
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
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

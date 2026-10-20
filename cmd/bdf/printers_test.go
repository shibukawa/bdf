package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/fakeprinter"
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
func fakePrinter(t *testing.T) (*fakeprinter.Printer, string) {
	t.Helper()
	p := fakeprinter.New()
	return p, fakeprinter.Start(t, p)
}

func TestPrintersURL(t *testing.T) {
	_, u := fakePrinter(t)
	code, out, errOut := runOut(t, "printers", "-url", u)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	for _, want := range []string{"name: fakeprinter", "model: BDF Fake Printer", "state: idle", "formats: application/pdf, image/pwg-raster",
		"sides: one-sided, two-sided-long-edge", "resolution: 300x300dpi, 600x600dpi", "copies: 1-99"} {
		if !strings.Contains(out, want) {
			t.Errorf("no %q in\n%s", want, out)
		}
	}
	code, out, _ = runOut(t, "printers", "-url", u, "-json")
	if code != 0 || !strings.Contains(out, `"MakeModel": "BDF Fake Printer"`) {
		t.Fatalf("-json: exit %d\n%s", code, out)
	}
	code, out, _ = runOut(t, "printers", "-cups", u)
	if code != 0 || !strings.HasPrefix(out, "fakeprinter\tipp://") {
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

// TestPrintRaw prints on printers of 9100 and LPD.
func TestPrintRaw(t *testing.T) {
	sock, lpd := fakeprinter.NewSocket(t), fakeprinter.NewLPD(t, "q1")
	defer sock.Close()
	defer lpd.Close()
	doc := filepath.Join("..", "..", "testdata", "form", "rows.bdf")
	code, _, errOut := runOut(t, append([]string{"print", "-printer", sock.URL, "-copies", "2", "-media", "a4", "-output-bin", "top"}, append(testFontFlags, doc)...)...)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if !strings.Contains(errOut, "page(s) sent as PostScript on iso_a4_210x297mm") || !strings.Contains(errOut, "does not support output-bin=top") {
		t.Errorf("output: %s", errOut)
	}
	if js := sock.Jobs(); len(js) != 1 || !strings.Contains(string(js[0].Data()), "<< /NumCopies 2 >>") {
		t.Fatalf("%d jobs", len(js))
	}
	code, _, errOut = runOut(t, append([]string{"print", "-printer", lpd.URL + "?format=pdf", "-sides", "short", "-q"}, append(testFontFlags, doc)...)...)
	if code != 0 || errOut != "" {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if js := lpd.Jobs(); len(js) != 1 || !strings.Contains(string(js[0].Data()), "@PJL SET BINDING=SHORTEDGE") || !strings.Contains(js[0].Control, "Jrows.bdf") {
		t.Fatalf("%d jobs", len(lpd.Jobs()))
	}
	lpd.Refuse(true)
	if code, _, errOut := runOut(t, append([]string{"print", "-printer", lpd.URL}, append(testFontFlags, doc)...)...); code != 4 || !strings.Contains(errOut, "refused") {
		t.Errorf("refused: exit %d: %s", code, errOut)
	}
	sock.Close()
	if code, _, errOut := runOut(t, append([]string{"print", "-printer", sock.URL}, append(testFontFlags, doc)...)...); code != 3 {
		t.Errorf("no printer: exit %d: %s", code, errOut)
	}
	if code, _, errOut := runOut(t, "print", "-printer", "lpd://host", doc); code != 2 || !strings.Contains(errOut, "lpd://host/queue") {
		t.Errorf("bad URL: exit %d: %s", code, errOut)
	}
}

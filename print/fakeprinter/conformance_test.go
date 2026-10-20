package fakeprinter_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/fakeprinter"
)

// ipptool returns the path of CUPS's ipptool and the directory of its
// tests, or skips: the fakes are checked against them where CUPS is (the
// CI's print job installs it, macOS has it).
func ipptool(t *testing.T) (tool, dir string) {
	t.Helper()
	tool, err := exec.LookPath("ipptool")
	if err != nil {
		t.Skip("no ipptool")
	}
	for _, d := range []string{"/usr/share/cups/ipptool", "/usr/local/share/cups/ipptool", "/opt/homebrew/share/cups/ipptool"} {
		if _, err := os.Stat(filepath.Join(d, "ipp-1.1.test")); err == nil {
			return tool, d
		}
	}
	t.Skip("no tests of ipptool")
	return "", ""
}

// document writes a document ipptool sends as the format its extension
// says (the fakes keep documents without reading them).
func document(t *testing.T, ext string) string {
	t.Helper()
	head := map[string]string{
		"pdf": "%PDF-1.4\n%fake\n", "ps": "%!PS-Adobe-3.0\nshowpage\n", "jpg": "\xff\xd8\xff\xe0fake",
		"pwg": "RaS2fake", "urf": "UNIRAST\x00fake",
	}[ext]
	name := filepath.Join(t.TempDir(), "doc."+ext)
	if err := os.WriteFile(name, []byte(head+strings.Repeat("x", 4096)), 0o644); err != nil {
		t.Fatal(err)
	}
	return name
}

var failLine = regexp.MustCompile(`(?m)^.*\[FAIL\]\n(?:\s{8}.*\n)*`)

// runIPPTool runs a test of ipptool on a printer and fails the test with
// what failed.
func runIPPTool(t *testing.T, version, file, uri, test string, flags ...string) {
	t.Helper()
	tool, _ := ipptool(t)
	args := append(append([]string{"-t", "-V", version, "-f", file}, flags...), uri, test)
	out, err := exec.Command(tool, args...).CombinedOutput()
	if err != nil || strings.Contains(string(out), "[FAIL]") {
		fails := failLine.FindAllString(string(out), -1)
		if len(fails) == 0 {
			t.Fatalf("ipptool %s: %v\n%s", filepath.Base(test), err, out)
		}
		t.Fatalf("ipptool %s on %s:\n%s", filepath.Base(test), uri, strings.Join(fails, ""))
	}
	if n := strings.Count(string(out), "[PASS]"); n == 0 {
		t.Fatalf("ipptool %s: nothing passed\n%s", filepath.Base(test), out)
	}
}

// TestIPPToolPresets checks every preset reached over IPP with ipptool's
// tests of IPP/1.1 and IPP/2.0, and the IPP Everywhere printer with the
// attributes PWG 5100.14 requires.
func TestIPPToolPresets(t *testing.T) {
	_, dir := ipptool(t)
	everywhere, _ := filepath.Abs("../../test/print/ipp-everywhere-server.test")
	for _, c := range []struct {
		preset, ext, version string
		tests                []string
		flags                []string
	}{
		{"office-mfp", "pdf", "2.0", []string{filepath.Join(dir, "ipp-2.0.test"), everywhere}, nil},
		{"airprint-only", "jpg", "2.0", []string{filepath.Join(dir, "ipp-2.0.test")}, nil},
		{"pwg-only", "pwg", "2.0", []string{filepath.Join(dir, "ipp-2.0.test")}, nil},
		{"ps-office", "ps", "2.0", []string{filepath.Join(dir, "ipp-2.0.test")}, nil},
		{"pclm-inkjet", "jpg", "2.0", []string{filepath.Join(dir, "ipp-2.0.test")}, nil},
		// an old printer takes requests of a known length only (-L)
		{"legacy-ipp11", "pdf", "1.1", []string{filepath.Join(dir, "ipp-1.1.test")}, []string{"-L"}},
	} {
		t.Run(c.preset, func(t *testing.T) {
			ps, ok := fakeprinter.LookupPreset(c.preset)
			if !ok {
				t.Fatal("no preset")
			}
			p := ps.New(c.preset)
			uri := fakeprinter.Start(t, p)
			// ipptool reaches the printer by name, as clients do
			uri = strings.Replace(uri, "127.0.0.1", "localhost", 1)
			file := document(t, c.ext)
			for _, test := range c.tests {
				runIPPTool(t, c.version, file, uri, test, c.flags...)
			}
			if len(p.Jobs()) == 0 {
				t.Error("ipptool's jobs did not reach the printer")
			}
		})
	}
}

// TestIPPToolCUPS checks the CUPS server's printers with ipptool, and its
// list of printers with CUPS's own test of CUPS-Get-Printers.
func TestIPPToolCUPS(t *testing.T) {
	_, dir := ipptool(t)
	c := fakeprinter.CUPSPreset()
	base := strings.Replace(fakeprinter.StartCUPS(t, c), "127.0.0.1", "localhost", 1)
	runIPPTool(t, "2.0", document(t, "pdf"), base+"/printers/office", filepath.Join(dir, "ipp-2.0.test"))
	runIPPTool(t, "2.0", document(t, "pdf"), base+"/classes/floor", filepath.Join(dir, "print-job.test"))
	runIPPTool(t, "2.0", document(t, "pdf"), base+"/", filepath.Join(dir, "get-printers.test"))
}

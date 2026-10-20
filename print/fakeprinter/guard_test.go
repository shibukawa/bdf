package fakeprinter_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// printerRole are the marks of a printer written by hand in a test: an
// IPP answer made or written, an LPD acknowledgement, a raw port read to
// its end by a listener of the test's own.
var printerRole = regexp.MustCompile(`ipp\.NewResponse\(|Header\(\)\.Set\("Content-Type", "application/ipp"\)|print/ipp/ipptest|print/raw/rawtest|ServeSocket\(|ServeLPD\(`)

// TestNoPrinterByHand keeps the tests of printing on the fake printers:
// none of the tests of print/ and of the bdf command plays a printer of
// its own (the fakes are checked against CUPS's tools; a printer written
// in a test is not).
func TestNoPrinterByHand(t *testing.T) {
	for _, root := range []string{"..", "../../cmd/bdf"} {
		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && filepath.Base(path) == "fakeprinter" {
				return filepath.SkipDir
			}
			if d.IsDir() || !strings.HasSuffix(path, "_test.go") {
				return nil
			}
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			for _, m := range printerRole.FindAllIndex(b, -1) {
				line := 1 + strings.Count(string(b[:m[0]]), "\n")
				t.Errorf("%s:%d: a printer of the test's own (%s); use print/fakeprinter", path, line, b[m[0]:m[1]])
			}
			return nil
		})
	}
}

package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// TestMain runs the test binary as the bdf command when the tests start it
// with BDF_TEST_MAIN set, so that they can check what the command does.
func TestMain(m *testing.M) {
	if os.Getenv("BDF_TEST_MAIN") != "" {
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// run runs the bdf command and returns its exit status and error output.
func run(t *testing.T, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "BDF_TEST_MAIN=1")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode(), stderr.String()
	}
	if err != nil {
		t.Fatal(err)
	}
	return 0, stderr.String()
}

// TestGeneratePages converts pages of a three-page TIFF file with -pages.
// An open range ("2-") used to be expanded to 2^30 page numbers before the
// page count was known.
func TestGeneratePages(t *testing.T) {
	in := filepath.Join("..", "..", "converter", "tiff", "testdata", "scan.tif")
	for _, c := range []struct {
		pages  string
		widths []float32 // of the pages converted (1000, 300 and 200 pixels at 300, 150 and 96 dpi)
		status int
		err    string
	}{
		{"2-", []float32{144, 150}, 0, ""},
		{"-2", []float32{240, 144}, 0, ""},
		{"3,1", []float32{150, 240}, 0, ""},
		{"1-3", []float32{240, 144, 150}, 0, ""},
		{"4-", nil, 1, "page 4 out of range (1-3)"},
		{"2-9", nil, 1, "page 4 out of range (1-3)"},
		{"0", nil, 2, "pages are numbered from 1"},
		{"3-2", nil, 2, "ends before it starts"},
	} {
		out := filepath.Join(t.TempDir(), "out.bdf")
		status, stderr := run(t, "generate", "-q", "-pages", c.pages, in, out)
		if status != c.status || !strings.Contains(stderr, c.err) {
			t.Errorf("-pages %s: exit status %d, %q; want %d and %q", c.pages, status, stderr, c.status, c.err)
			continue
		}
		if status != 0 {
			continue
		}
		r, err := bdf.OpenSingleFile(out)
		if err != nil {
			t.Fatal(err)
		}
		var widths []float32
		for _, p := range r.Manifest.Views[0].Pages {
			widths = append(widths, p.W)
		}
		if len(widths) != len(c.widths) {
			t.Errorf("-pages %s: page widths %v, want %v", c.pages, widths, c.widths)
			continue
		}
		for i := range widths {
			if d := widths[i] - c.widths[i]; d > 0.01 || d < -0.01 {
				t.Errorf("-pages %s: page widths %v, want %v", c.pages, widths, c.widths)
				break
			}
		}
	}
}

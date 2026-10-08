package all_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/all"
)

// FuzzConvert converts an input of any registered format, as a server that
// converts what is uploaded to it does: the format is told by the content,
// or by the file name's extension. Malformed input is an error (or a
// document with warnings), never a panic, a hang or memory out of
// proportion with the input.
func FuzzConvert(f *testing.F) {
	exts := map[string]bool{}
	for _, fm := range converter.Formats() {
		for _, e := range fm.Extensions {
			exts[e] = true
		}
	}
	dirs, _ := filepath.Glob("../*/testdata")
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !exts[strings.ToLower(filepath.Ext(name))] {
				continue
			}
			b, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil || len(b) > 512<<10 {
				continue // large seeds make every iteration slow
			}
			f.Add(name, b)
		}
	}
	opts := converter.Options{
		FontDirs:      []string{"../pptx/testdata/fonts", "../docx/testdata/fonts"},
		NoSystemFonts: true,
		Params:        map[string]string{"remote": "false"},
	}
	f.Fuzz(func(t *testing.T, name string, data []byte) {
		defer slowInput(t, name, data)()
		o := opts
		o.FileName = name
		o.Warn = func(string) {}
		res, err := converter.Convert(bytes.NewReader(data), int64(len(data)), "", &o)
		if err != nil {
			return
		}
		var out bytes.Buffer
		if err := res.Doc.WriteSingle(&out); err != nil {
			t.Fatalf("write: %v", err)
		}
	})
}

// slowInput writes an input that takes longer than a second to convert to
// the directory BDF_FUZZ_SLOW_DIR names, when it is set: the fuzzer keeps
// the inputs that reach new code, not the slow ones.
func slowInput(t *testing.T, name string, data []byte) func() {
	dir := os.Getenv("BDF_FUZZ_SLOW_DIR")
	if dir == "" {
		return func() {}
	}
	start := time.Now()
	return func() {
		if d := time.Since(start); d > time.Second {
			ext := strings.ToLower(filepath.Ext(name))
			if len(ext) > 12 || strings.ContainsAny(ext, "/\\") {
				ext = ""
			}
			out := filepath.Join(dir, fmt.Sprintf("%.1fs-%x%s", d.Seconds(), sha256.Sum256(data), ext))
			os.WriteFile(out, data, 0o600)
			t.Logf("slow input (%v): %s", d, out)
		}
	}
}

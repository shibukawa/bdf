package all_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

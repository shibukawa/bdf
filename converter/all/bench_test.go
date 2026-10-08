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

// BenchmarkConvert converts the test documents of several formats, whole,
// with the repository's fonts: what a server spends on an upload, after
// the first conversion has loaded the fonts.
func BenchmarkConvert(b *testing.B) {
	for _, name := range []string{"pdf/testdata/chrome-doc.pdf", "pptx/testdata/features.pptx", "xlsx/testdata/basic.xlsx",
		"docx/testdata/basic.docx", "epub/testdata/basic.epub", "drawio/testdata/showcase.drawio", "markdown/testdata/basic.md"} {
		data, err := os.ReadFile("../" + name)
		if err != nil {
			b.Fatal(err)
		}
		opts := converter.Options{FontDirs: []string{"../pptx/testdata/fonts", "../docx/testdata/fonts"}, NoSystemFonts: true,
			FileName: name, Params: map[string]string{"remote": "false"}, Warn: func(string) {}}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res, err := converter.Convert(bytes.NewReader(data), int64(len(data)), "", &opts)
				if err != nil {
					b.Fatal(err)
				}
				if err := res.Doc.WriteSingle(&bytes.Buffer{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertFiles converts the files BDF_BENCH_FILES names (paths
// separated by commas), for documents outside the repository.
func BenchmarkConvertFiles(b *testing.B) {
	files := os.Getenv("BDF_BENCH_FILES")
	if files == "" {
		b.Skip("BDF_BENCH_FILES is not set")
	}
	for _, name := range strings.Split(files, ",") {
		data, err := os.ReadFile(name)
		if err != nil {
			b.Fatal(err)
		}
		opts := converter.Options{FontDirs: []string{"../pptx/testdata/fonts", "../docx/testdata/fonts"}, NoSystemFonts: true,
			FileName: filepath.Base(name), Params: map[string]string{"remote": "false"}, Warn: func(string) {}}
		b.Run(filepath.Base(name), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				res, err := converter.Convert(bytes.NewReader(data), int64(len(data)), "", &opts)
				if err != nil {
					b.Fatal(err)
				}
				if err := res.Doc.WriteSingle(&bytes.Buffer{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

package bdf

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// These benchmarks include all parts of real documents, so changes to
// compression are measured with their fonts, paths and page objects.
func BenchmarkDocumentIO(b *testing.B) {
	for _, name := range []string{"pdf/chrome-doc", "drawio/showcase", "epub/basic"} {
		data, err := os.ReadFile("testdata/" + name + ".bdf")
		if err != nil {
			b.Fatal(err)
		}
		r, err := OpenSingle(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			b.Fatal(err)
		}
		doc, err := r.ToDocument()
		if err != nil {
			b.Fatal(err)
		}
		b.Run(name+"/read", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := r.ToDocument(); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(name+"/write", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if err := doc.WriteSingle(io.Discard); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

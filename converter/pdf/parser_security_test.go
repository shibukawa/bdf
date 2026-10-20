package pdf

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
)

func TestPDFParserResourceLimits(t *testing.T) {
	objects := []any{
		`<< /Type /Catalog /Pages 2 0 R >>`,
		`<< /Type /Pages /Kids [3 0 R] /Count 1 >>`,
		`<< /Type /Page /Parent 2 0 R /MediaBox [0 0 8 8] /Contents 4 0 R >>`,
		[2]string{"", "0 0 4 4 re f"},
	}
	deep := append([]any(nil), objects...)
	deep[0] = `<< /Type /Catalog /Pages 2 0 R /Extra ` + strings.Repeat("[", 1000) + "0" + strings.Repeat("]", 1000) + " >>"
	stream := append([]any(nil), objects...)
	stream[3] = "<< /Length 1099511627776 >>\nstream\na\nendstream"
	xref := bytes.Replace(buildPDF(objects, "/Root 1 0 R"), []byte("/Size 5"), []byte("/Size 2147483647"), 1)
	for name, c := range map[string]struct {
		data []byte
		// mayRead is set for input a reader may repair rather than refuse:
		// since pdfcpu v0.16.3, a trailer /Size past MaxObjectCount is left
		// unset in relaxed validation and the object walk supplies it. What
		// matters is that it costs no more than a small input should.
		mayRead bool
	}{
		"deep objects":  {data: buildPDF(deep, "/Root 1 0 R")},
		"stream length": {data: buildPDF(stream, "/Root 1 0 R")},
		"xref size":     {data: xref, mayRead: true},
	} {
		t.Run(name, func(t *testing.T) {
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			_, err := NewStream(bytes.NewReader(c.data), nil)
			runtime.ReadMemStats(&after)
			if err == nil && !c.mayRead {
				t.Fatal("hostile parser input was accepted")
			}
			if n := after.TotalAlloc - before.TotalAlloc; n > 16<<20 {
				t.Errorf("small hostile input allocated %d bytes", n)
			}
		})
	}
}

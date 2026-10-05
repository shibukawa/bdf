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
	for name, data := range map[string][]byte{
		"deep objects":  buildPDF(deep, "/Root 1 0 R"),
		"stream length": buildPDF(stream, "/Root 1 0 R"),
		"xref size":     xref,
	} {
		t.Run(name, func(t *testing.T) {
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			_, err := NewStream(bytes.NewReader(data), nil)
			runtime.ReadMemStats(&after)
			if err == nil {
				t.Fatal("hostile parser input was accepted")
			}
			if n := after.TotalAlloc - before.TotalAlloc; n > 16<<20 {
				t.Errorf("small hostile input allocated %d bytes", n)
			}
		})
	}
}

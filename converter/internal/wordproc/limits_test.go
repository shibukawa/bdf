package wordproc

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

const limNS = `xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" ` +
	`xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"`

// limDocx builds a .docx with the given body and converts it.
func limDocx(t testing.TB, body, views string) *Result {
	t.Helper()
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	add := func(name, data string) {
		w, err := z.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(data))
	}
	add("_rels/.rels", `<?xml version="1.0"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`+
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`)
	add("word/document.xml", `<?xml version="1.0"?><w:document `+limNS+`><w:body>`+body+`</w:body></w:document>`)
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	opts := &Options{FontDirs: []string{"../../docx/testdata/fonts"}, NoSystemFonts: true, Views: views, NoTextIndex: true}
	res, err := ConvertDOCX(bytes.NewReader(buf.Bytes()), int64(buf.Len()), opts)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// A cell whose gridSpan states a huge column count does not allocate a grid
// of that size: it is capped with a warning.
func TestGridSpanCapped(t *testing.T) {
	body := `<w:tbl><w:tr><w:tc><w:tcPr><w:gridSpan w:val="2000000000"/></w:tcPr><w:p/></w:tc></w:tr></w:tbl>`
	res := limDocx(t, body, ViewsPages)
	found := false
	for _, w := range res.Warnings {
		if strings.Contains(w, "more than") && strings.Contains(w, "columns") {
			found = true
		}
	}
	if !found {
		t.Errorf("a gridSpan of 2000000000 should warn about the column count; warnings: %v", res.Warnings)
	}
}

// Benchmark the long keepNext chain. Its page-placement behavior is covered
// by the DOCX tests; wall-clock regression limits are too sensitive to CI load.
func BenchmarkKeepNextParagraphs(b *testing.B) {
	for _, size := range []int{10000, 40000} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			body := strings.Repeat(`<w:p><w:pPr><w:keepNext/></w:pPr></w:p>`, size)
			b.ResetTimer()
			for range b.N {
				limDocx(b, body, ViewsPages)
			}
		})
	}
}

// A paragraph's custom tab stops are capped so resolving them (quadratic in
// their number) stays bounded regardless of how many the markup lists.
func TestTabStopsCapped(t *testing.T) {
	c, _, err := newConverter(&Options{}, ViewsBoth)
	if err != nil {
		t.Fatal(err)
	}
	tabsPPr := func(n int) *ooxml.Node {
		var b strings.Builder
		b.WriteString(`<pPr><tabs>`)
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&b, `<tab val="left" pos="%d"/>`, i*10)
		}
		b.WriteString(`</tabs></pPr>`)
		root, perr := ooxml.Parse([]byte(b.String()))
		if perr != nil {
			t.Fatal(perr)
		}
		return root
	}
	// under the cap the stops are all kept
	p := defaultPProps()
	c.applyPPr(p, tabsPPr(5))
	if len(p.tabs) != 5 {
		t.Errorf("5 tab stops: kept %d, want 5", len(p.tabs))
	}
	// far above the cap the count is bounded
	p = defaultPProps()
	c.applyPPr(p, tabsPPr(maxTabs*4))
	if len(p.tabs) != maxTabs {
		t.Errorf("%d tab stops: kept %d, want the cap %d", maxTabs*4, len(p.tabs), maxTabs)
	}
}

// A list or page number falls back to a decimal number when a hostile start
// value would repeat a letter far more than any real list.
func TestLettersCapped(t *testing.T) {
	if got := letters(28); got != "BB" {
		t.Errorf("letters(28) = %q, want BB", got)
	}
	if got := letters(2000000000); got != strconv.Itoa(2000000000) {
		t.Errorf("letters(2000000000) = %q, want the decimal fallback", got)
	}
}

package htmlro

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// benchDoc is an attribute-heavy page of n rows: 2n elements, 5n
// attributes, n text nodes with a reference each.
func benchDoc(n int) []byte {
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html><head><title>bench</title><meta charset="utf-8"><style>p { margin: 0 }</style></head><body>`)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, `<div class="c%d row" id="i%d" data-index="%d" title="row %d"><span class="s">text %d &amp; more</span></div>`+"\n", i%7, i, i, i, i)
	}
	sb.WriteString(`<script>var x = "<p>";</script></body></html>`)
	return []byte(sb.String())
}

// benchInput is benchDoc(5000), or the file BENCH_HTML names.
func benchInput(b *testing.B) []byte {
	if path := os.Getenv("BENCH_HTML"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			b.Fatal(err)
		}
		return data
	}
	return benchDoc(5000)
}

func BenchmarkReader(b *testing.B) {
	data := benchInput(b)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	r := NewBytesReader(data, Options{})
	var buf []byte
	for i := 0; i < b.N; i++ {
		r.ResetBytes(data)
		for {
			k, err := r.Next()
			if err != nil || k == EOF {
				break
			}
			switch k {
			case StartTag, SelfClosingTag:
				r.Name()
				for {
					_, v, ok := r.NextAttr()
					if !ok {
						break
					}
					buf = v.AppendTo(buf[:0])
				}
			case Text:
				buf = r.Text().AppendTo(buf[:0])
			}
		}
	}
}

func BenchmarkXNetTokenizer(b *testing.B) {
	data := benchInput(b)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		z := html.NewTokenizer(bytes.NewReader(data))
		for {
			tt := z.Next()
			if tt == html.ErrorToken {
				break
			}
			switch tt {
			case html.StartTagToken, html.SelfClosingTagToken:
				_, more := z.TagName()
				for more {
					_, _, more = z.TagAttr()
				}
			case html.TextToken:
				z.Text()
			}
		}
	}
}

func BenchmarkParseBytes(b *testing.B) {
	data := benchInput(b)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := ParseBytes(data, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	data := benchInput(b)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Parse(bytes.NewReader(data), Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkXNetParse(b *testing.B) {
	data := benchInput(b)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := html.Parse(bytes.NewReader(data)); err != nil {
			b.Fatal(err)
		}
	}
}

//go:build !tinygo

package htmlro

import (
	"testing"

	"golang.org/x/net/html"
)

func FuzzReaderMatchesTokenizer(f *testing.F) {
	for _, tt := range tokenTests {
		f.Add(tt.html, true, "")
	}
	f.Add("<script><!--<script </script>x</script>", false, "")
	f.Add("<title>&amp;</title>", false, "title")
	f.Fuzz(func(t *testing.T, input string, cdata bool, fragment string) {
		if len(fragment) > 16 {
			return
		}
		compareTokens(t, input, cdata, fragment)
	})
}

func FuzzParseMatchesXNet(f *testing.F) {
	for _, tt := range tokenTests {
		f.Add(tt.html, true)
	}
	f.Add("<table><tr><td>a<b>b</td>c</table>d", true)
	f.Add("<svg><![CDATA[x]]><foreignObject><p>y", false)
	f.Fuzz(func(t *testing.T, input string, scripting bool) {
		compareParse(t, "fuzz", []byte(input), scripting)
	})
}

func FuzzUnescapeString(f *testing.F) {
	for _, tt := range tokenTests {
		f.Add(tt.html)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if got, want := UnescapeString(s), html.UnescapeString(s); got != want {
			t.Errorf("%q: got %q, want %q", s, got, want)
		}
	})
}

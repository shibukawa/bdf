package webdoc

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/net/html/charset"
)

// DetermineEncoding finds what golang.org/x/net/html/charset finds.
func TestDetermineEncodingMatchesXNet(t *testing.T) {
	cases := []struct{ content, contentType string }{
		{"", ""},
		{"\xef\xbb\xbf<p>x", ""},
		{"\xfe\xff\x00<", ""},
		{"\xff\xfe<\x00", "text/html; charset=shift_jis"},
		{"<p>x", "text/html; charset=shift_jis"},
		{"<p>x", "text/html; charset=nonsense"},
		{"<p>x", "nonsense"},
		{`<meta charset="EUC-JP">`, ""},
		{`<META CHARSET=euc-jp>`, ""},
		{`<meta charset=utf-16>`, ""},
		{`<meta charset="UTF-16LE">`, ""},
		{`<meta http-equiv="Content-Type" content="text/html; charset=windows-1252">`, ""},
		{`<meta content="text/html; charset=big5" http-equiv=content-type>`, ""},
		{`<meta content="text/html; charset=big5">`, ""},
		{`<meta http-equiv="refresh" content="charset=big5">`, ""},
		{`<meta http-equiv="content-type" content="text/html; charset='gb18030'">`, ""},
		{`<meta http-equiv="content-type" content="text/html; charset = iso-8859-2 ; x">`, ""},
		{`<meta http-equiv="content-type" content="text/html; charset=">`, ""},
		{`<meta http-equiv="content-type" content="text/html; charset='unterminated">`, ""},
		{`<meta charset=x charset=euc-jp>`, ""},
		{`<meta charset="euc-jp" charset="x">`, ""},
		{`<meta name=x><meta charset=koi8-r>`, ""},
		{`<!-- <meta charset=euc-jp> --><meta charset=koi8-r>`, ""},
		{`<script><meta charset=euc-jp></script><meta charset=koi8-r>`, ""},
		{`<meta charset=&quot;euc-jp&quot;>`, ""},
		{`<meta charset="euc&#45;jp">`, ""},
		{"<p>\xe3\x81\x82", ""},
		{"<p>\xe3\x81", ""},
		{"<p>\xff\xfe\xfd", ""},
		{"<p>plain ascii only", ""},
		{"<p>x" + string(make([]byte, 1100)) + `<meta charset=euc-jp>`, ""},
		{string(make([]byte, 1020)) + `<meta charset=euc-jp>`, ""},
	}
	for _, c := range cases {
		e1, n1, c1 := DetermineEncoding([]byte(c.content), c.contentType)
		e2, n2, c2 := charset.DetermineEncoding([]byte(c.content), c.contentType)
		if n1 != n2 || c1 != c2 || (e1 == nil) != (e2 == nil) {
			t.Errorf("%q %q: got %q %v (nil %v), want %q %v (nil %v)", c.content, c.contentType, n1, c1, e1 == nil, n2, c2, e2 == nil)
		}
	}
	// And on the real pages of a corpus, when there is one.
	dir := os.Getenv("HTMLRO_CORPUS")
	if dir == "" {
		return
	}
	n := 0
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		_, n1, c1 := DetermineEncoding(data, "")
		_, n2, c2 := charset.DetermineEncoding(data, "")
		if n1 != n2 || c1 != c2 {
			t.Errorf("%s: got %q %v, want %q %v", path, n1, c1, n2, c2)
		}
		n++
		return nil
	})
	t.Logf("%d documents compared", n)
}

func FuzzDetermineEncodingMatchesXNet(f *testing.F) {
	f.Add(`<meta charset="EUC-JP">`, "")
	f.Add(`<meta http-equiv="Content-Type" content="text/html; charset=windows-1252">`, "text/html")
	f.Add("\xef\xbb\xbf<p>x", "")
	f.Fuzz(func(t *testing.T, content, contentType string) {
		_, n1, c1 := DetermineEncoding([]byte(content), contentType)
		_, n2, c2 := charset.DetermineEncoding([]byte(content), contentType)
		if n1 != n2 || c1 != c2 {
			t.Errorf("%q %q: got %q %v, want %q %v", content, contentType, n1, c1, n2, c2)
		}
	})
}

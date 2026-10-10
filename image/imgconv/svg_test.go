package imgconv

import "testing"

func TestSVGSize(t *testing.T) {
	for _, c := range []struct {
		width, height, viewBox string
		w, h                   float64
	}{
		{"200", "100", "", 200, 100},
		{"200px", "1in", "0 0 10 10", 200, 96},
		{"72pt", "2.54cm", "", 96, 96},
		{"1e2", "5em", "", 100, 80},
		{"", "", "0 0 40 30", 40, 30},
		{"80", "", "0,0,40,30", 80, 60},
		{"100%", "50", "0 0 40 20", 100, 50},
		{"", "", "", 300, 150},
		{"50", "", "", 50, 150},
		{"-5", "abc", "0 0 0 0", 300, 150},
	} {
		if w, h := ParseSVGSize(c.width, c.height, c.viewBox).Pixels(); w != c.w || h != c.h {
			t.Errorf("ParseSVGSize(%q, %q, %q).Pixels() = %v × %v, want %v × %v", c.width, c.height, c.viewBox, w, h, c.w, c.h)
		}
	}
}

func TestSVGRoot(t *testing.T) {
	for _, c := range []struct {
		data, width, height, viewBox string
		ok                           bool
	}{
		{`<?xml version="1.0"?><!-- c --><svg xmlns="http://www.w3.org/2000/svg" width="4em" viewBox="0 0 8 4"><rect/></svg>`, "4em", "", "0 0 8 4", true},
		{"\xef\xbb\xbf<svg:svg xmlns:svg=\"http://www.w3.org/2000/svg\" height='10'/>", "", "10", "", true},
		{`<html><svg width="1"/></html>`, "", "", "", false},
		{`not xml`, "", "", "", false},
	} {
		w, h, vb, ok := SVGRoot([]byte(c.data))
		if w != c.width || h != c.height || vb != c.viewBox || ok != c.ok {
			t.Errorf("SVGRoot(%q) = %q %q %q %v", c.data, w, h, vb, ok)
		}
	}
	if r := ParseSVGSize("", "", "0 0 800 400").Ratio(); r != 2 {
		t.Errorf("ratio %v", r)
	}
	if s := ParseSVGSize("30", "", ""); s.Ratio() != 0 || s.H != 0 {
		t.Errorf("width only %+v", s)
	}
}

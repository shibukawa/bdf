package wordproc

import (
	"math"
	"testing"
)

func TestNumberFormats(t *testing.T) {
	for _, c := range []struct {
		n      int
		format string
		want   string
	}{
		{4, "upperRoman", "IV"}, {14, "lowerRoman", "xiv"}, {28, "upperLetter", "BB"}, {3, "lowerLetter", "c"},
		{12, "decimalFullWidth", "１２"}, {3, "decimalEnclosedCircle", "③"}, {21, "japaneseCounting", "二十一"},
		{105, "ideographDigital", "一〇五"}, {2, "aiueoFullWidth", "イ"}, {3, "iroha", "ﾊ"}, {7, "decimalZero", "07"},
		{22, "ordinal", "22nd"}, {1234, "japaneseCounting", "千二百三十四"}, {12000, "japaneseCounting", "一万二千"},
	} {
		if got := formatNumber(c.n, c.format); got != c.want {
			t.Errorf("formatNumber(%d, %s) = %q, want %q", c.n, c.format, got, c.want)
		}
	}
}

func TestAutoColumns(t *testing.T) {
	a := &autoColumns{min: []float64{10, 20, 30}, max: []float64{40, 20, 100}}
	for _, c := range []struct {
		avail float64
		fill  bool
		want  []float64
	}{
		{200, false, []float64{40, 20, 100}}, // the content fits: its widths
		{320, true, []float64{80, 40, 200}},  // filling: in proportion to the content
		{120, false, []float64{28, 20, 72}},  // between: the slack shared by the difference
		{30, false, []float64{5, 10, 15}},    // too narrow: the minimums shrink
		{60, false, []float64{10, 20, 30}},   // exactly the minimums
	} {
		a.fill = c.fill
		got := a.widths(c.avail)
		for i := range got {
			if math.Abs(got[i]-c.want[i]) > 1e-9 {
				t.Errorf("widths(%v, fill %v) = %v, want %v", c.avail, c.fill, got, c.want)
				break
			}
		}
	}
	a.width, a.fill = 90, false
	if got := a.widths(1000); math.Abs(got[0]+got[1]+got[2]-90) > 1e-9 {
		t.Errorf("width 90: %v", got)
	}
}

func TestParaWidths(t *testing.T) {
	pp := defaultPProps()
	p := &para{pp: pp, mark: &runStyle{size: 10}}
	word := func(s string, w float64) {
		for i, r := range s {
			p.items = append(p.items, item{kind: kChar, r: r, w: w, brk: i == len(s)-1 && r != ' '})
		}
	}
	word("ab", 5)
	word(" ", 3)
	word("cdef", 5)
	p.items = append(p.items, item{kind: kBreak})
	word("g", 5)
	p.items = append(p.items, item{kind: kObject, obj: &inlineObj{w: 500, fit: true}, w: 500})
	mn, mx := paraWidths(p)
	// the widest word; the widest line (a shrinking picture needs no width)
	if mn != 20 || mx != 505 {
		t.Errorf("widths %v, %v", mn, mx)
	}
}

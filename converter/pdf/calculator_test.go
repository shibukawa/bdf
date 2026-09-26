package pdf

import (
	"slices"
	"testing"
)

func TestCalculator(t *testing.T) {
	// The stripe function Chrome writes for repeating-linear-gradient.
	skia := `{pop
dup truncate sub
dup 0 le {1 add} if
dup 0 le {pop .149 .5451 .8235 0} if
dup dup 0 gt exch 1 le and {
dup .5 le {
0 sub pop .149 .5451 .8235
}{
.5 sub pop .8275 .2118 .5098
} ifelse
0} if
0 gt {.8275 .2118 .5098} if
}`
	for _, tc := range []struct {
		src  string
		in   []float64
		want []float64
	}{
		{"{1 exch sub}", []float64{0.25}, []float64{0.75}},
		{"{dup 0.5 gt {pop 1} {pop 0} ifelse}", []float64{0.7}, []float64{1}},
		{"{dup 0.5 gt {pop 1} {pop 0} ifelse}", []float64{0.3}, []float64{0}},
		{"{2 copy add 3 1 roll mul}", []float64{2, 3}, []float64{5, 6}},
		{"{1 index 2 index}", []float64{7, 8}, []float64{7, 8, 7, 7}},
		{"{3 -1 roll}", []float64{1, 2, 3}, []float64{2, 3, 1}},
		{"{7 2 idiv 7 2 mod -7 2 mod 1 3 bitshift 16 -2 bitshift}", nil, []float64{3, 1, -1, 8, 4}},
		{"{2.5 round -2.5 round 2.7 cvi 90 sin 0 cos 1 0 atan 4 sqrt 2 3 exp}", nil, []float64{3, -2, 2, 1, 1, 90, 2, 8}},
		{"{true false or 5 3 and 1 2 lt not}", nil, []float64{1, 1, 0}},
		{skia, []float64{1.25, 0}, []float64{.149, .5451, .8235}},
		{skia, []float64{2.75, 0}, []float64{.8275, .2118, .5098}},
		{"{pop}", nil, nil},     // stack underflow
		{"{1 0 div}", nil, nil}, // division by zero
		{"{1 {2} 3}", nil, nil}, // a procedure outside if/ifelse
		{"{1 foo}", nil, nil},   // unknown operator
		{"{1 2 add", nil, nil},  // unterminated
	} {
		prog := compilePS([]byte(tc.src))
		got := runPS(prog, tc.in)
		if prog == nil {
			got = nil
		}
		if len(tc.want) == 0 && len(got) == 0 {
			continue
		}
		if !slices.EqualFunc(got, tc.want, func(a, b float64) bool { return a-b < 1e-9 && b-a < 1e-9 }) {
			t.Errorf("%s %v = %v, want %v", tc.src, tc.in, got, tc.want)
		}
	}
	f := &pdfFunction{kind: 4, domain: []float64{0, 1}, rng: []float64{0, 1}, prog: compilePS([]byte("{1 exch sub}"))}
	if out := f.eval(0.2); len(out) != 1 || out[0] != 0.8 {
		t.Errorf("type 4 eval = %v", out)
	}
	bad := &pdfFunction{kind: 4, domain: []float64{0, 1}, rng: []float64{0, 1, 0, 2}}
	if out := bad.eval(0.2); len(out) != 2 || out[0] != 0.5 || out[1] != 1 {
		t.Errorf("unreadable type 4 eval = %v", out)
	}
}

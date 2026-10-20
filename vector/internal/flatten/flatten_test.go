package flatten

import (
	"math"
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	var r Regions
	if !r.Empty() || r.Merge(100, 100, 1) != nil {
		t.Fatal("no boxes are no rectangles")
	}
	r.Add([4]float64{10.5, 10.5, 20, 20})
	r.Add([4]float64{21, 10, 30, 15})   // within the margin of the first
	r.Add([4]float64{60, 60, 70, 70})   // apart
	r.Add([4]float64{95, 95, 120, 120}) // past the page
	r.Add([4]float64{5, 5, 5, 9})       // empty
	r.Add([4]float64{math.NaN(), 0, 1, 1})
	got := r.Merge(100, 100, 1)
	want := [][4]float64{{9, 9, 31, 21}, {59, 59, 71, 71}, {94, 94, 100, 100}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%v, want %v", got, want)
	}
	// merging makes rectangles that overlap others
	r = Regions{}
	r.Add([4]float64{0, 0, 10, 10})
	r.Add([4]float64{20, 20, 30, 30})
	r.Add([4]float64{5, 25, 6, 26}) // in the box of the two
	r.Add([4]float64{8, 8, 22, 22})
	if got := r.Merge(1000, 1000, 0); !reflect.DeepEqual(got, [][4]float64{{0, 0, 30, 30}}) {
		t.Errorf("%v", got)
	}
	// most of the page is the page
	r = Regions{}
	r.Add([4]float64{0, 0, 80, 80})
	if got := r.Merge(100, 100, 1); !reflect.DeepEqual(got, [][4]float64{{0, 0, 100, 100}}) {
		t.Errorf("%v", got)
	}
	// too many rectangles are the page
	r = Regions{}
	for i := 0; i < 100; i++ {
		r.Add([4]float64{float64(i * 10), 0, float64(i*10 + 1), 1})
	}
	if got := r.Merge(1000, 1000, 0); len(got) != 1 || got[0] != [4]float64{0, 0, 1000, 1000} {
		t.Errorf("%d rectangles", len(got))
	}
	r = Regions{}
	for i := 0; i <= maxBoxes; i++ {
		r.Add([4]float64{0, 0, 1, 1})
	}
	if got := r.Merge(1000, 1000, 0); len(got) != 1 || got[0] != [4]float64{0, 0, 1000, 1000} {
		t.Errorf("%d rectangles", len(got))
	}
}

package music

import (
	"slices"
	"testing"
)

func TestPlaceAccidentals(t *testing.T) {
	c4 := &Note{Pitch: Pitch{Step: 0, Octave: 4}}
	d4 := &Note{Pitch: Pitch{Step: 1, Octave: 4}}
	c5 := &Note{Pitch: Pitch{Step: 0, Octave: 5}}
	without := &Note{Pitch: Pitch{Step: 2, Octave: 4}}
	sm := &staffMeasure{clef: Clef{Sign: ClefG}, acc: map[*Note]Accidental{
		c4: AccSharp, d4: AccNatural, c5: AccFlat, without: AccNone,
	}}
	ev := &Event{Notes: []*Note{c4, c5, without, d4}}
	got := placeAccidentals(sm, ev)
	if len(got) != 3 {
		t.Fatalf("%d accidentals, want 3", len(got))
	}
	if got[0].note != c5 || got[0].col != 0 || got[0].acc != AccFlat ||
		got[1].note != d4 || got[1].col != 0 || got[1].acc != AccNatural ||
		got[2].note != c4 || got[2].col != 1 || got[2].acc != AccSharp {
		t.Errorf("accidental order and columns: %+v", got)
	}
}

func TestColumnLefts(t *testing.T) {
	if got, want := columnLefts(100, map[int]float64{0: 2, 1: 3, 2: 4}), []float64{98, 95, 91}; !slices.Equal(got, want) {
		t.Errorf("column lefts %v, want %v", got, want)
	}
}

func TestPrepareLayoutCaches(t *testing.T) {
	time := &TimeSig{Beats: 3, BeatType: 4}
	s := &Score{
		Measures: []*Measure{
			{Length: 4 * PPQ},
			{Length: 4 * PPQ, Number: "7", Time: time},
			{Length: 4 * PPQ, Implicit: true},
			{Length: 4 * PPQ},
		},
		Parts: []*Part{{Staves: 2}, {Staves: 1}},
	}
	e := &engraver{s: s}
	e.prepare()
	if got := e.timeAt(0); *got != (TimeSig{Beats: 4, BeatType: 4}) {
		t.Errorf("default time %+v", got)
	}
	for i := 1; i < len(s.Measures); i++ {
		if got := e.timeAt(i); got != time {
			t.Errorf("measure %d time %+v, want the explicit time", i, got)
		}
	}
	for i, want := range []string{"1", "7", "7", "8"} {
		if got := e.measureNumber(i); got != want {
			t.Errorf("measure %d number %q, want %q", i, got, want)
		}
	}
	if got := e.bracketWidth(); got != 1.2 {
		t.Errorf("bracket width %v, want 1.2", got)
	}
	sys := &system{y: []float64{0, 10.5, 21}}
	tops := e.staffTops(sys)
	if len(tops) != 2 || !slices.Equal(tops[0], []float64{0, 10.5}) || !slices.Equal(tops[1], []float64{21}) {
		t.Errorf("staff tops %v", tops)
	}
}

var benchmarkAccidentals []placedAccidental

func BenchmarkPlaceAccidentals(b *testing.B) {
	sm := &staffMeasure{clef: Clef{Sign: ClefG}, acc: map[*Note]Accidental{}}
	ev := &Event{Type: Quarter}
	for range 6000 {
		n := &Note{Pitch: Pitch{Step: 0, Alter: 1, Octave: 4}}
		ev.Notes = append(ev.Notes, n)
		sm.acc[n] = AccSharp
	}
	b.ResetTimer()
	for range b.N {
		benchmarkAccidentals = placeAccidentals(sm, ev)
	}
}

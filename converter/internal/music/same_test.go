package music

import (
	"math/rand"
	"slices"
	"sort"
	"strconv"
	"testing"
	"time"
)

// The engraver keeps what it used to work out again for every measure and
// every chord. These tests compare it with the code it had, kept here.

// placeAccidentalsBefore looks at every accidental placed so far for each
// column.
func placeAccidentalsBefore(sm *staffMeasure, ev *Event) []placedAccidental {
	var out []placedAccidental
	for _, n := range ev.Notes {
		if a, ok := sm.acc[n]; ok && a != AccNone {
			out = append(out, placedAccidental{note: n, acc: a, pos: staffPos(sm.clef, n.Pitch)})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].pos < out[j].pos })
	for i := range out {
		for c := 0; ; c++ {
			ok := true
			for j := 0; j < i; j++ {
				if out[j].col == c && out[i].pos-out[j].pos < 6 {
					ok = false
					break
				}
			}
			if ok {
				out[i].col = c
				break
			}
		}
	}
	return out
}

// randChord makes a chord of n notes within span staff positions, some with
// accidentals.
func randChord(r *rand.Rand, n, span int) (*staffMeasure, *Event) {
	sm := &staffMeasure{clef: Clef{Sign: ClefSign(r.Intn(3))}, acc: map[*Note]Accidental{}}
	ev := &Event{Type: Quarter}
	for i := 0; i < n; i++ {
		p := 28 + r.Intn(span) // from C4
		note := &Note{Pitch: Pitch{Step: p % 7, Octave: p / 7}}
		ev.Notes = append(ev.Notes, note)
		if r.Intn(4) > 0 {
			sm.acc[note] = Accidental(1 + r.Intn(6)) // AccNone too
		}
	}
	return sm, ev
}

func TestPlaceAccidentalsSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 2000; i++ {
		sm, ev := randChord(r, 1+r.Intn(40), 1+r.Intn(30))
		got, want := placeAccidentals(sm, ev), placeAccidentalsBefore(sm, ev)
		if !slices.Equal(got, want) {
			t.Fatalf("chord %d: %+v, want %+v", i, got, want)
		}
	}
}

func TestColumnLeftsSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 500; i++ {
		w := map[int]float64{}
		for k := r.Intn(12); k >= 0; k-- {
			w[k] = (0.7 + r.Float64()) * []float64{1, graceScale}[r.Intn(2)]
		}
		minX := 100 * r.Float64()
		left := columnLefts(minX, w)
		for col := range w {
			// as it was: the widths taken off for each accidental
			x := minX
			for k := 0; k <= col; k++ {
				x -= w[k]
			}
			if left[col] != x {
				t.Fatalf("run %d column %d at %v, want %v", i, col, left[col], x)
			}
		}
	}
}

func TestPlaceAccidentalsTime(t *testing.T) {
	// 6000 accidentals on one line take 6000 columns: 18 million steps,
	// where they were 36 billion
	sm := &staffMeasure{clef: Clef{Sign: ClefG}, acc: map[*Note]Accidental{}}
	ev := &Event{Type: Quarter}
	for i := 0; i < 6000; i++ {
		n := &Note{Pitch: Pitch{Step: 0, Alter: 1, Octave: 4}}
		ev.Notes = append(ev.Notes, n)
		sm.acc[n] = AccSharp
	}
	start := time.Now()
	pa := placeAccidentals(sm, ev)
	if d := time.Since(start); d > 2*time.Second {
		t.Errorf("took %v", d)
	}
	for i, a := range pa {
		if a.col != i {
			t.Fatalf("accidental %d in column %d", i, a.col)
		}
	}
}

// timeAtBefore goes back to the last time signature.
func timeAtBefore(s *Score, mi int) *TimeSig {
	for i := mi; i >= 0; i-- {
		if t := s.Measures[i].Time; t != nil {
			return t
		}
	}
	return &TimeSig{Beats: 4, BeatType: 4}
}

// measureNumberBefore counts from the first measure.
func measureNumberBefore(s *Score, mi int) string {
	if n := s.Measures[mi].Number; n != "" {
		return n
	}
	n := 0
	for i := 0; i <= mi; i++ {
		m := s.Measures[i]
		if m.Number != "" {
			if v, err := strconv.Atoi(m.Number); err == nil {
				n = v
				continue
			}
		}
		if !m.Implicit {
			n++
		}
	}
	return strconv.Itoa(n)
}

// staffTopsBefore goes through the staves for a part.
func staffTopsBefore(e *engraver, sys *system, part int) []float64 {
	var out []float64
	for si, ref := range e.staves {
		if ref.part == part {
			out = append(out, sys.y[si])
		}
	}
	return out
}

// bracketWidthBefore looks for a grand staff among the parts.
func bracketWidthBefore(s *Score) float64 {
	grand := slices.ContainsFunc(s.Parts, func(p *Part) bool { return p.Staves > 1 })
	w := 0.0
	if len(s.Parts) > 1 && !grand {
		w = 1.6
	}
	if grand {
		w = 1.2
	}
	return w
}

func TestPrepareSame(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	numbers := []string{"", "", "", "1", "7", "12", "0", "-3", "X1", "5a", " 9", "007"}
	for i := 0; i < 300; i++ {
		s := &Score{}
		for k := r.Intn(40); k >= 0; k-- {
			m := &Measure{Length: 4 * PPQ, Number: numbers[r.Intn(len(numbers))], Implicit: r.Intn(5) == 0}
			if r.Intn(4) == 0 {
				m.Time = &TimeSig{Beats: 1 + r.Intn(12), BeatType: 1 << r.Intn(5)}
			}
			s.Measures = append(s.Measures, m)
		}
		for k := r.Intn(6); k >= 0; k-- {
			s.Parts = append(s.Parts, &Part{Staves: r.Intn(4) * r.Intn(2)}) // 0 is one staff
		}
		e := &engraver{s: s}
		e.prepare()
		for mi := range s.Measures {
			if got, want := e.timeAt(mi), timeAtBefore(s, mi); *got != *want || s.Measures[mi].Time != nil && got != want {
				t.Fatalf("score %d measure %d: time %+v, want %+v", i, mi, got, want)
			}
			if got, want := e.measureNumber(mi), measureNumberBefore(s, mi); got != want {
				t.Fatalf("score %d measure %d: number %q, want %q", i, mi, got, want)
			}
		}
		if got, want := e.bracketWidth(), bracketWidthBefore(s); got != want {
			t.Fatalf("score %d: bracket width %g, want %g", i, got, want)
		}
		sys := &system{y: make([]float64, len(e.staves))}
		for si := range sys.y {
			sys.y[si] = 10.5 * float64(si)
		}
		tops := e.staffTops(sys)
		for pi := range s.Parts {
			if want := staffTopsBefore(e, sys, pi); !slices.Equal(tops[pi], want) {
				t.Fatalf("score %d part %d: staves at %v, want %v", i, pi, tops[pi], want)
			}
		}
	}
}

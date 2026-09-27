package music

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

// dump writes the events of a part's staff as text: one line per measure,
// events as pitch/rest, note value, dots, tuplet and ties.
func dump(p *Part, staff int) string {
	var b strings.Builder
	for i, m := range p.Measures {
		if i > 0 {
			b.WriteString(" | ")
		}
		first := true
		evs := slices.Clone(m.Events)
		slices.SortStableFunc(evs, func(a, b *Event) int {
			if a.Offset != b.Offset {
				return a.Offset - b.Offset
			}
			return a.Voice - b.Voice
		})
		for _, e := range evs {
			if e.Staff != staff {
				continue
			}
			if !first {
				b.WriteByte(' ')
			}
			first = false
			if e.Voice != 1 {
				fmt.Fprintf(&b, "v%d:", e.Voice)
			}
			switch {
			case e.MeasureRest:
				b.WriteString("R")
			case e.Rest:
				b.WriteString("r")
			default:
				for k, n := range e.Notes {
					if k > 0 {
						b.WriteByte('+')
					}
					if n.TieStop {
						b.WriteByte('~')
					}
					b.WriteString(pitchName(n.Pitch))
				}
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", 1<<max(int(e.Type)-1, 0))
			b.WriteString(strings.Repeat(".", e.Dots))
			if e.Tuplet != nil {
				fmt.Fprintf(&b, "/%d", e.Tuplet.Actual)
			}
			if len(e.Notes) > 0 && e.Notes[0].TieStart {
				b.WriteByte('~')
			}
		}
	}
	return b.String()
}

func pitchName(p Pitch) string {
	s := string("cdefgab"[p.Step])
	switch p.Alter {
	case 1:
		s += "#"
	case -1:
		s += "b"
	case 2:
		s += "x"
	case -2:
		s += "bb"
	}
	return fmt.Sprintf("%s%d", s, p.Octave)
}

// melody makes a one-track performance from (key, length in ticks)
// pairs; key -1 is a rest.
func melody(notes ...int) *Performance {
	t := &Track{Name: "Melody", Volume: -1, Pan: -1}
	tick := 0
	for i := 0; i+1 < len(notes); i += 2 {
		if notes[i] >= 0 {
			t.Notes = append(t.Notes, PlayNote{Tick: tick, Dur: notes[i+1], Key: notes[i], Vel: 100})
		}
		tick += notes[i+1]
	}
	return &Performance{Tracks: []*Track{t}, End: tick}
}

const (
	q  = PPQ
	h  = 2 * PPQ
	e8 = PPQ / 2
	s  = PPQ / 4
	t3 = PPQ / 3
)

func TestNotate(t *testing.T) {
	cases := []struct {
		name string
		perf *Performance
		opts NotateOptions
		want string
	}{
		{"twinkle", melody(60, q, 60, q, 67, q, 67, q, 69, q, 69, q, 67, h),
			NotateOptions{}, "c4 4 c4 4 g4 4 g4 4 | a4 4 a4 4 g4 2"},
		{"across the bar line", melody(60, 3*q, 62, h, 64, 3*q),
			NotateOptions{}, "c4 2. d4 4~ | ~d4 4 e4 2."},
		{"syncopation in a beat", melody(60, s, 62, e8, 64, s, 65, q, -1, h),
			NotateOptions{}, "c4 16 d4 8 e4 16 f4 4 r2"},
		{"off the beat", melody(-1, e8, 60, q, 62, e8, 64, h),
			NotateOptions{}, "r8 c4 8~ ~c4 8 d4 8 e4 2"},
		{"triplets", melody(60, t3, 62, t3, 64, t3, 65, q, 67, h),
			NotateOptions{}, "c4 8/3 d4 8/3 e4 8/3 f4 4 g4 2"},
		{"rests", melody(-1, q, 60, q, -1, h),
			NotateOptions{}, "r4 c4 4 r2"},
		{"rest across beats", melody(60, q, -1, h, 62, q),
			NotateOptions{}, "c4 4 r4 r4 d4 4"},
		{"empty measure", melody(60, 4*q, -1, 4*q, 60, 4*q),
			NotateOptions{}, "c4 1 | R1 | c4 1"},
		{"six eight", melody(60, 3*e8, 62, e8, 64, e8, 65, e8, 67, 6*e8),
			NotateOptions{Time: &TimeSig{6, 8, ""}}, "c4 4. d4 8 e4 8 f4 8 | g4 2."},
		{"sharps in G", melody(67, q, 69, q, 71, q, 66, q),
			NotateOptions{Key: &KeySig{Fifths: 1}}, "g4 4 a4 4 b4 4 f#4 4"},
		{"flats in F", melody(65, q, 70, q, 69, q, 68, q),
			NotateOptions{Key: &KeySig{Fifths: -1}}, "f4 4 bb4 4 a4 4 ab4 4"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sc := Notate(c.perf, c.opts)
			if len(sc.Parts) != 1 {
				t.Fatalf("%d parts", len(sc.Parts))
			}
			if got := dump(sc.Parts[0], 0); got != c.want {
				t.Errorf("got  %s\nwant %s", got, c.want)
			}
		})
	}
}

func TestNotateVoices(t *testing.T) {
	// a held note under a moving line
	tr := &Track{Name: "Piano", Volume: -1, Pan: -1, Notes: []PlayNote{
		{Tick: 0, Dur: 4 * q, Key: 64, Vel: 90},
		{Tick: 0, Dur: q, Key: 72, Vel: 90},
		{Tick: q, Dur: q, Key: 74, Vel: 90},
		{Tick: 2 * q, Dur: h, Key: 76, Vel: 90},
	}}
	sc := Notate(&Performance{Tracks: []*Track{tr}}, NotateOptions{Key: &KeySig{}})
	if got, want := dump(sc.Parts[0], 0), "c5 4 v2:e4 1 d5 4 e5 2"; got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

func TestNotateGrandStaff(t *testing.T) {
	tr := &Track{Name: "Piano", Volume: -1, Pan: -1, Notes: []PlayNote{
		{Tick: 0, Dur: h, Key: 72, Vel: 90},
		{Tick: 0, Dur: h, Key: 48, Vel: 90},
		{Tick: h, Dur: h, Key: 76, Vel: 90},
		{Tick: h, Dur: h, Key: 43, Vel: 90},
	}}
	sc := Notate(&Performance{Tracks: []*Track{tr}}, NotateOptions{Key: &KeySig{}})
	p := sc.Parts[0]
	if p.Staves != 2 {
		t.Fatalf("staves %d", p.Staves)
	}
	if got, want := dump(p, 0)+" / "+dump(p, 1), "c5 2 e5 2 / c3 2 g2 2"; got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
	if c := p.Measures[0].Clefs; len(c) != 2 || c[1].Clef.Sign != ClefF {
		t.Errorf("clefs %+v", c)
	}
}

func TestEstimateKey(t *testing.T) {
	// D major scale
	var hist [12]float64
	for _, k := range []int{62, 64, 66, 67, 69, 71, 73, 74, 62, 66, 69} {
		hist[k%12]++
	}
	if k := estimateKey(hist, 0); k.Fifths != 2 || k.Minor {
		t.Errorf("D major: %+v", k)
	}
}

func TestSpell(t *testing.T) {
	cases := []struct {
		key  int
		sig  KeySig
		hint int
		want string
	}{
		{61, KeySig{}, 0, "c#4"},
		{63, KeySig{}, 0, "eb4"},
		{66, KeySig{}, 0, "f#4"},
		{70, KeySig{}, 0, "bb4"},
		{63, KeySig{}, 1, "d#4"},
		{68, KeySig{Minor: true}, 0, "g#4"},
		{61, KeySig{Fifths: -1, Minor: true}, 0, "c#4"},
		{65, KeySig{Fifths: 6}, 0, "e#4"},
		{59, KeySig{Fifths: -6}, 0, "cb4"},
		{60, KeySig{Fifths: 7}, 0, "b#3"},
	}
	for _, c := range cases {
		if got := pitchName(spell(c.key, c.sig, c.hint)); got != c.want {
			t.Errorf("spell(%d, %+v, %d) = %s, want %s", c.key, c.sig, c.hint, got, c.want)
		}
	}
}

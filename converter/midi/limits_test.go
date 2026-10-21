package midi

import (
	"testing"

	"github.com/shibukawa/bdf/converter/internal/music"
)

func TestTimeSignatureAfterBarLine(t *testing.T) {
	// a time signature one tick after a bar line, under a note: the
	// measure of one tick is written
	a := trk(nil).meta(0, 0x58, "\x04\x02\x18\x08").ev(0, 0x90, 60, 100).meta(1, 0x58, "\x03\x02\x18\x08").ev(959, 0x80, 60, 64).
		note(0, 0, 62, 960).end(0)
	for _, division := range []int{960, 1000} {
		p := parse(t, smf(0, division, a))
		s := music.Notate(p.Perf, music.NotateOptions{})
		if len(s.Measures) != 2 || s.Measures[0].Length != 1 || s.Measures[1].Length != 3*music.PPQ {
			t.Fatalf("division %d: %d measures, the first of %d ticks", division, len(s.Measures), s.Measures[0].Length)
		}
	}
}

func TestTracksBudget(t *testing.T) {
	// 300 track chunks of 15 channels each are 4500 tracks (more than a
	// system has staves, which music.Notate refuses: see its tests)
	var tracks []trk
	for i := 0; i < 300; i++ {
		a := trk(nil)
		for ch := 0; ch < 15; ch++ {
			a = a.note(0, ch, 60, 96)
		}
		tracks = append(tracks, a.end(0))
	}
	p := parse(t, smf(1, 96, tracks...))
	if len(p.Perf.Tracks) != 4500 {
		t.Fatalf("%d tracks", len(p.Perf.Tracks))
	}
}

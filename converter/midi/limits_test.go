package midi

import (
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/music"
)

func TestTrackAfterTheLastMeasure(t *testing.T) {
	// the score is cut after 20000 measures, here of 1/64: a track that
	// starts after them has no part
	a := trk(nil).meta(0, 0x58, "\x01\x06\x18\x08").note(0, 0, 60, 1).end(0)
	b := trk(nil).note(1300, 1, 64, 1).end(0)
	p := parse(t, smf(1, 1, a, b))
	if len(p.Perf.Tracks) != 2 {
		t.Fatalf("%d tracks", len(p.Perf.Tracks))
	}
	s := music.Notate(p.Perf, music.NotateOptions{})
	if len(s.Measures) != 20000 || len(s.Parts) != 1 || len(s.Parts[0].Measures) != 20000 {
		t.Errorf("%d measures, %d parts", len(s.Measures), len(s.Parts))
	}
}

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
		s.SMF = p.SMF
		res, err := music.Build(s, music.Options{Source: "midi", FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true})
		if err != nil || res.Notes != 3 {
			t.Errorf("division %d: %v", division, err)
		}
	}
}

func TestTracksBudget(t *testing.T) {
	// 300 track chunks of 15 channels: more tracks than a system has staves
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
	s := music.Notate(p.Perf, music.NotateOptions{})
	_, err := music.Build(s, music.Options{Source: "midi", NoSystemFonts: true})
	if len(s.Parts) != 0 || err == nil || !strings.Contains(err.Error(), "more than 4096 staves") {
		t.Errorf("%d parts, error %v", len(s.Parts), err)
	}

	// 60 tracks of 20000 measures: more than a million measures of staves
	tracks = tracks[:4]
	tracks[0] = slices.Concat(tracks[0][:len(tracks[0])-4], trk(nil).note(96*4*20000-96*2, 0, 60, 96).end(0))
	p = parse(t, smf(1, 96, tracks...))
	s = music.Notate(p.Perf, music.NotateOptions{})
	_, err = music.Build(s, music.Options{Source: "midi", NoSystemFonts: true})
	if len(s.Parts) != 0 || err == nil || !strings.Contains(err.Error(), "measures on its staves (60 tracks of 20000 measures)") {
		t.Errorf("%d parts, error %v", len(s.Parts), err)
	}
}

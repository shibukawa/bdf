package mml

import (
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/music"
)

func TestTracksBudget(t *testing.T) {
	for _, src := range []string{
		strings.Repeat("c;", maxTracks+1),
		"MML@" + strings.Repeat("c,", maxTracks) + "c;",
	} {
		perf, warnings, err := Parse([]byte(src), nil)
		if err != nil {
			t.Fatalf("%.20q: %v", src, err)
		}
		if len(perf.Tracks) != maxTracks {
			t.Errorf("%.20q: %d tracks", src, len(perf.Tracks))
		}
		if !slices.ContainsFunc(warnings, func(w string) bool { return strings.Contains(w, "more than 1024 tracks") }) {
			t.Errorf("%.20q: warnings %q", src, warnings)
		}
	}
	perf, warnings, err := Parse([]byte(strings.Repeat("c;", maxTracks)), nil)
	if err != nil || len(perf.Tracks) != maxTracks || len(warnings) != 0 {
		t.Errorf("%d tracks: %v, warnings %q", maxTracks, err, warnings)
	}
}

func TestMeasuresOfTracks(t *testing.T) {
	// 101 tracks, one of 10000 measures: more than a million measures of
	// staves, refused before they are written
	perf, _, err := Parse([]byte("[c1]10000;"+strings.Repeat("c;", 100)), nil)
	if err != nil || len(perf.Tracks) != 101 {
		t.Fatal(err)
	}
	s := music.Notate(perf, music.NotateOptions{})
	_, err = music.Build(s, music.Options{Source: "mml", NoSystemFonts: true})
	if len(s.Parts) != 0 || err == nil || !strings.Contains(err.Error(), "measures on its staves (101 tracks of 10000 measures)") {
		t.Errorf("%d parts, error %v", len(s.Parts), err)
	}
}

func TestTrackAfterTheLastMeasure(t *testing.T) {
	// in 2/4 the second track starts where the 20000 measures written end
	perf, _, err := Parse([]byte("c;[r1]10000c;"), nil)
	if err != nil || len(perf.Tracks) != 2 {
		t.Fatal(err)
	}
	s := music.Notate(perf, music.NotateOptions{Time: &music.TimeSig{Beats: 2, BeatType: 4}})
	if len(s.Measures) != 20000 || len(s.Parts) != 1 {
		t.Errorf("%d measures, %d parts", len(s.Measures), len(s.Parts))
	}
}

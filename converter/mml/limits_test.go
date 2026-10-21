package mml

import (
	"slices"
	"strings"
	"testing"
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

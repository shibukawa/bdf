package midi

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/music"
)

// engrave writes a performance as a score and engraves it, as the
// converter does: it may fail, but must not panic. Music that takes long
// is left to the tests of the budgets.
func engrave(perf *music.Performance, source string) {
	notes := 0
	for _, t := range perf.Tracks {
		notes += len(t.Notes)
	}
	if notes > 20_000 || len(perf.Tracks)*(perf.End/music.PPQ+1) > 100_000 {
		return
	}
	score := music.Notate(perf, music.NotateOptions{})
	staves := 0
	for _, p := range score.Parts {
		staves += p.Staves
	}
	if staves*len(score.Measures) > 3000 {
		return
	}
	music.Build(score, music.Options{Source: source, NoSystemFonts: true, Warn: func(string) {}})
}

// FuzzParse feeds damaged files to the reader: it may fail, but must not
// panic, and what it reads must hold together and be engraved. The seeds
// are the test files and some damaged ones.
func FuzzParse(f *testing.F) {
	for _, pattern := range []string{"testdata/*.mid", "testdata/*.kar", "testdata/*.rmi"} {
		files, _ := filepath.Glob(pattern)
		for _, name := range files {
			if b, err := os.ReadFile(name); err == nil && len(b) < 64<<10 {
				f.Add(b)
			}
		}
	}
	f.Add(smf(0, 0xE728, trk(nil).note(1000, 0, 60, 500).end(0)))
	f.Add(smf(2, 100, trk(nil).meta(0, 0x20, "\x01").meta(0, 0x05, "a ").note(3, 1, 60, 7).ev(0, 0xB0, 123, 0), trk(nil).ev(0, 0xF0, 1, 0xF7)))
	f.Add([]byte("MThd\x00\x00\x00\x06\x00\x01\x00\x03\x01\xE0MTrk\x00\x00\x00\x10\x00\x90\x3C\x40\x60\x3C\x00junkMTrk\x00\x00"))
	// a time signature one tick after a bar line, and a track that starts
	// after the last measure written
	f.Add(smf(0, 960, trk(nil).ev(0, 0x90, 60, 100).meta(1, 0x58, "\x03\x02\x18\x08").ev(959, 0x80, 60, 64).end(0)))
	f.Add(smf(1, 1, trk(nil).meta(0, 0x58, "\x01\x06\x18\x08").note(0, 0, 60, 1).end(0), trk(nil).note(1300, 1, 64, 1).end(0)))
	f.Fuzz(func(t *testing.T, b []byte) {
		p, err := Parse(b, &Options{Warn: func(string) {}})
		if err != nil {
			return
		}
		count := 0
		for _, tr := range p.Perf.Tracks {
			if len(tr.Notes) == 0 || tr.Channel < 0 || tr.Channel > 15 || tr.Program < 0 || tr.Program > 127 {
				t.Fatalf("track %+v", tr)
			}
			if !slices.IsSortedFunc(tr.Notes, func(a, b music.PlayNote) int { return a.Tick - b.Tick }) {
				t.Fatal("notes out of order")
			}
			for _, n := range tr.Notes {
				if n.Dur < 1 || n.Tick < 0 || n.Key < 0 || n.Key > 127 || n.Vel < 1 || n.Vel > 127 || n.Tick+n.Dur > p.Perf.End {
					t.Fatalf("note %+v, end %d", n, p.Perf.End)
				}
			}
			count += len(tr.Notes)
		}
		engrave(p.Perf, "midi")
		if p.SMF == nil {
			return
		}
		// the file stored reads as the same music
		q, err := Parse(p.SMF, &Options{Warn: func(string) {}})
		if err != nil {
			t.Fatalf("the SMF stored does not read: %v", err)
		}
		recount := 0
		for _, tr := range q.Perf.Tracks {
			recount += len(tr.Notes)
		}
		if recount != count || q.Perf.End != p.Perf.End || len(q.Perf.Tempo) != len(p.Perf.Tempo) {
			t.Fatalf("the SMF stored reads %d notes to %d, not %d to %d", recount, q.Perf.End, count, p.Perf.End)
		}
	})
}

package mml

import (
	"os"
	"path/filepath"
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

// FuzzParse feeds damaged MML to the reader: it may fail, but must not
// panic, and what it reads must be a performance the engine can take
// (notes in order, in the MIDI range, of positive length) and engraves.
// The seeds are the test files and some commands.
func FuzzParse(f *testing.F) {
	files, _ := filepath.Glob("testdata/*.mml")
	for _, name := range files {
		if b, err := os.ReadFile(name); err == nil && len(b) < 64<<10 {
			f.Add(b)
		}
	}
	for _, s := range []string{
		"t120 l8 o4 [c d | e]3 /:2 f / g :/ {abc}4 c4&c8 d^16 n60,8 c%96 r4.. >c<c;",
		"$a=cde; $b=$a$a; @1 @W12 @p30 @v100 $b (2 ) ns-2 c; #OCTAVE REVERSE\n",
		"MML@t150l8o5cde>c<b,o3l2cg,v15o4e1;MML@efg;",
		"#TITLE x\n#OCTAVE-REV\n@v0 = { 15 14\n}\nABC t150 l8 ; comment\nA @2 K2 [c d|e]2 {cde}4 L w\nD @0 c\nE c",
		"{{{c}}}[[[c]]]/:/:c:/:/[c:d|e]{}[]0",
	} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		Detect(b[:min(len(b), 1024)], nil, int64(len(b)))
		perf, _, err := Parse(b, nil)
		if err != nil {
			return
		}
		if len(perf.Tempo) == 0 || perf.Tempo[0].Tick != 0 {
			t.Fatalf("tempo %v", perf.Tempo)
		}
		for i, tm := range perf.Tempo {
			if tm.BPM <= 0 || i > 0 && tm.Tick <= perf.Tempo[i-1].Tick {
				t.Fatalf("tempo %v", perf.Tempo)
			}
		}
		for _, tr := range perf.Tracks {
			if tr.Channel < 0 || tr.Channel > 15 || tr.Channel == 9 || tr.Program < 0 || tr.Program > 127 || len(tr.Notes) == 0 {
				t.Fatalf("track %q: channel %d program %d, %d notes", tr.Name, tr.Channel, tr.Program, len(tr.Notes))
			}
			for i, n := range tr.Notes {
				if n.Dur <= 0 || n.Key < 0 || n.Key > 127 || n.Vel < 1 || n.Vel > 127 || n.Tick < 0 ||
					i > 0 && n.Tick < tr.Notes[i-1].Tick+tr.Notes[i-1].Dur || n.Tick+n.Dur > perf.End {
					t.Fatalf("track %q note %d: %+v (end %d)", tr.Name, i, n, perf.End)
				}
			}
			for i, c := range tr.Controls {
				if c.Tick < 0 || i > 0 && c.Tick < tr.Controls[i-1].Tick {
					t.Fatalf("track %q controls %v", tr.Name, tr.Controls)
				}
			}
		}
		engrave(perf, "mml")
	})
}

package mml

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/music"
	xunicode "golang.org/x/text/encoding/unicode"
)

// The test files are self-written: traditional tunes (Twinkle, Twinkle,
// Frère Jacques) and Beethoven's Ode to Joy.

// n is a note as the tests check it.
type n struct{ tick, dur, key, spell int }

func notes(tr *music.Track) []n {
	var out []n
	for _, x := range tr.Notes {
		out = append(out, n{x.Tick, x.Dur, x.Key, x.Spell})
	}
	return out
}

func keys(tr *music.Track) []int {
	var out []int
	for _, x := range tr.Notes {
		out = append(out, x.Key)
	}
	return out
}

func parse(t *testing.T, src string, o *Options) (*music.Performance, []string) {
	t.Helper()
	perf, warnings, err := Parse([]byte(src), o)
	if err != nil {
		t.Fatalf("%q: %v", src, err)
	}
	return perf, warnings
}

func parseFile(t *testing.T, name string) (*music.Performance, []string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	perf, warnings, err := Parse(b, nil)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return perf, warnings
}

const q, e, h, w = music.PPQ, music.PPQ / 2, 2 * music.PPQ, 4 * music.PPQ

func TestNotes(t *testing.T) {
	normal := &Options{Octave: OctaveNormal}
	for _, tc := range []struct {
		name, src string
		o         *Options
		want      []n
	}{
		{"scale", "t120 l4 o4 cdefgab>c", nil, []n{{0, q, 60, 0}, {q, q, 62, 0}, {2 * q, q, 64, 0}, {3 * q, q, 65, 0},
			{4 * q, q, 67, 0}, {5 * q, q, 69, 0}, {6 * q, q, 71, 0}, {7 * q, q, 72, 0}}},
		{"lengths", "c8 d4. e2.. f16 r4 g l8 a b.", normal, []n{{0, e, 60, 0}, {e, q + e, 62, 0},
			{2 * q, h + q + e, 64, 0}, {5*q + e, q / 4, 65, 0}, {6*q + e + q/4, q, 67, 0},
			{7*q + e + q/4, e, 69, 0}, {8*q + q/4, e + e/2, 71, 0}}},
		{"triplet lengths", "c12c12c12 d3", normal, []n{{0, 320, 60, 0}, {320, 320, 60, 0}, {640, 320, 60, 0}, {q, 1280, 62, 0}}},
		{"accidentals", "c+ d# e- f++ b- c+-", normal, []n{{0, q, 61, 1}, {q, q, 63, 1}, {2 * q, q, 63, -1},
			{3 * q, q, 67, 1}, {4 * q, q, 70, -1}, {5 * q, q, 60, 0}}},
		{"tie", "c4&c8 d", normal, []n{{0, q + e, 60, 0}, {q + e, q, 62, 0}}},
		{"legato", "c4&d4", normal, []n{{0, q, 60, 0}, {q, q, 62, 0}}},
		{"tie over commands", "c4& o4 l8 c", normal, []n{{0, q + e, 60, 0}}},
		{"no tie after a rest", "c4& r c", normal, []n{{0, q, 60, 0}, {2 * q, q, 60, 0}}},
		{"caret", "c4^8^ d", normal, []n{{0, q + e + q, 60, 0}, {2*q + e, q, 62, 0}}},
		{"caret after a rest", "r4^8 c", normal, []n{{q + e, q, 60, 0}}},
		{"note numbers", "n48 n60,8 n61.", normal, []n{{0, q, 60, 0}, {q, e, 72, 0}, {q + e, q + e, 73, 0}}},
		{"percent generic", "c%96 d%48", normal, []n{{0, q, 60, 0}, {q, e, 62, 0}}},
		{"percent ppmck", "A c%48 d%24", &Options{Dialect: PPMCK}, []n{{0, q, 60, 0}, {q, e, 62, 0}}},
		{"octave", "o5c o3c o-1c", normal, []n{{0, q, 72, 0}, {q, q, 48, 0}, {2 * q, q, 0, 0}}},
		{"octave normal", "o4c>c<<c", normal, []n{{0, q, 60, 0}, {q, q, 72, 0}, {2 * q, q, 48, 0}}},
		{"octave reverse", "o4c>c<<c", &Options{Octave: OctaveReverse}, []n{{0, q, 60, 0}, {q, q, 48, 0}, {2 * q, q, 72, 0}}},
		{"out of range", "o9 g a b >c", normal, []n{{0, q, 127, 0}}},
		{"upper case", "T120 L8 O5 C D+ E-", nil, []n{{0, e, 72, 0}, {e, e, 75, 1}, {2 * e, e, 75, -1}}},
		{"loop", "[cd|e]3", normal, []n{{0, q, 60, 0}, {q, q, 62, 0}, {2 * q, q, 64, 0}, {3 * q, q, 60, 0},
			{4 * q, q, 62, 0}, {5 * q, q, 64, 0}, {6 * q, q, 60, 0}, {7 * q, q, 62, 0}}},
		{"loop default count and colon", "l8[c:d]", normal, []n{{0, e, 60, 0}, {e, e, 62, 0}, {2 * e, e, 60, 0}}},
		{"nested loops", "l8[c[d]2]2", normal, []n{{0, e, 60, 0}, {e, e, 62, 0}, {2 * e, e, 62, 0},
			{3 * e, e, 60, 0}, {4 * e, e, 62, 0}, {5 * e, e, 62, 0}}},
		{"loop keeps state", "o4 l8 [c>]3", normal, []n{{0, e, 60, 0}, {e, e, 72, 0}, {2 * e, e, 84, 0}}},
		{"flmml loop", "l8 /:3 c / d :/", normal, []n{{0, e, 60, 0}, {e, e, 62, 0}, {2 * e, e, 60, 0},
			{3 * e, e, 62, 0}, {4 * e, e, 60, 0}}},
		{"flmml nested loops", "l8 /:2 c /: d :/ :/", normal, []n{{0, e, 60, 0}, {e, e, 62, 0}, {2 * e, e, 62, 0},
			{3 * e, e, 60, 0}, {4 * e, e, 62, 0}, {5 * e, e, 62, 0}}},
		{"loop once", "[c]0 [d]1", normal, []n{{0, q, 60, 0}, {q, q, 62, 0}}},
		{"unclosed loop", "[cd", normal, []n{{0, q, 60, 0}, {q, q, 62, 0}}},
		{"tuplet", "{cde}4 f", normal, []n{{0, 320, 60, 0}, {320, 320, 62, 0}, {640, 320, 64, 0}, {q, q, 65, 0}}},
		{"tuplet of eighths", "{c8d8e8}2", normal, []n{{0, 640, 60, 0}, {640, 640, 62, 0}, {1280, 640, 64, 0}}},
		{"tuplet with a rest", "{crd}8", normal, []n{{0, 160, 60, 0}, {320, 160, 62, 0}}},
		{"tuplet of unequal notes", "{c4d8}4", normal, []n{{0, 640, 60, 0}, {640, 320, 62, 0}}},
		{"tuplet tie", "c4&{cd}4", normal, []n{{0, q + e, 60, 0}, {q + e, e, 62, 0}}},
		{"tuplet then tie", "{cd}4&d4", normal, []n{{0, e, 60, 0}, {e, e + q, 62, 0}}},
		{"silent v0", "v15c v0d v8e", normal, []n{{0, q, 60, 0}, {2 * q, q, 64, 0}}},
		{"comments", "c /* d */ e // f\n g", normal, []n{{0, q, 60, 0}, {q, q, 64, 0}, {2 * q, q, 67, 0}}},
		{"unknown commands", "c z12 d k3 @z e", normal, []n{{0, q, 60, 0}, {q, q, 62, 0}, {2 * q, q, 64, 0}}},
		{"ignored commands", "q6 @e1,0,10,64 @l10,20 @f1 x1 c", normal, []n{{0, q, 60, 0}}},
		{"flmml note shift", "ns2 c ns-3 c", normal, []n{{0, q, 62, 0}, {q, q, 59, 0}}},
		{"ppmck transpose", "A K2 c K-1 c K0 c", &Options{Dialect: PPMCK}, []n{{0, q, 62, 0}, {q, q, 59, 0}, {2 * q, q, 60, 0}}},
		{"ppmck wait and loop point", "A c w L c", &Options{Dialect: PPMCK}, []n{{0, q, 60, 0}, {2 * q, q, 60, 0}}},
		{"ppmck ignored", "A EN0 MP1 D2 @v0 @q3 s1,2 c ENOF", &Options{Dialect: PPMCK}, []n{{0, q, 60, 0}}},
		{"ppmck semicolon comment", "A c ; d\nA e", &Options{Dialect: PPMCK}, []n{{0, q, 60, 0}, {q, q, 64, 0}}},
		{"mabinogi octave", "MML@o4c>c<c;", nil, []n{{0, q, 60, 0}, {q, q, 72, 0}, {2 * q, q, 60, 0}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			perf, _ := parse(t, tc.src, tc.o)
			if got := notes(perf.Tracks[0]); !slices.Equal(got, tc.want) {
				t.Errorf("%q:\n got %v\nwant %v", tc.src, got, tc.want)
			}
		})
	}
}

func TestTupletSplit(t *testing.T) {
	perf, _ := parse(t, "{cdefgab}2 c", &Options{Octave: OctaveNormal})
	ns := perf.Tracks[0].Notes
	if len(ns) != 8 {
		t.Fatalf("%d notes", len(ns))
	}
	tick := 0
	for i, x := range ns[:7] {
		if x.Tick != tick || x.Dur != h/7 && x.Dur != h/7+1 {
			t.Errorf("note %d: tick %d dur %d", i, x.Tick, x.Dur)
		}
		tick += x.Dur
	}
	if tick != h || ns[7].Tick != h {
		t.Errorf("the tuplet ends at %d, the note after it starts at %d", tick, ns[7].Tick)
	}
}

func TestOctaveDirection(t *testing.T) {
	for _, tc := range []struct {
		src  string
		o    *Options
		want []int
	}{
		// the generic dialect: FlMML raises with <, MSX with >
		{"o4 l8 gab<cde", nil, []int{67, 69, 71, 72, 74, 76}},
		{"o4 l8 gab>cde", nil, []int{67, 69, 71, 72, 74, 76}},
		{"o5 l8 c<bag >c", nil, []int{72, 71, 69, 67, 72}},
		{"#OCTAVE REVERSE\no4 l8 gab>cde", nil, []int{67, 69, 71, 72, 74, 76}},
		// no evidence: > raises
		{"o4 c>c<c", nil, []int{60, 72, 60}},
		// set
		{"o4 l8 gab<cde", &Options{Octave: OctaveNormal}, []int{67, 69, 71, 48, 50, 52}},
		{"o4 l8 gab>cde", &Options{Octave: OctaveReverse}, []int{67, 69, 71, 48, 50, 52}},
		// PPMCK: > raises unless #OCTAVE-REV
		{"A o4 l8 gab<cde", &Options{Dialect: PPMCK}, []int{67, 69, 71, 48, 50, 52}},
		{"#OCTAVE-REV\nA o4 l8 gab<cde", nil, []int{67, 69, 71, 72, 74, 76}},
		{"#OCTAVE-REV\nA o4 l8 gab<cde", &Options{Octave: OctaveNormal}, []int{67, 69, 71, 48, 50, 52}},
		// Mabinogi: > raises
		{"MML@o4l8gab<cde;", nil, []int{67, 69, 71, 48, 50, 52}},
	} {
		perf, _ := parse(t, tc.src, tc.o)
		if got := keys(perf.Tracks[0]); !slices.Equal(got, tc.want) {
			t.Errorf("%q: got %v, want %v", tc.src, got, tc.want)
		}
	}
}

func TestVelocity(t *testing.T) {
	perf, _ := parse(t, "c v15c v8c v+2c v-10c v0c v1c @v100c (2c )c", &Options{Octave: OctaveNormal})
	var got []int
	for _, x := range perf.Tracks[0].Notes {
		got = append(got, x.Vel)
	}
	// v12 by default; the notes of v0 are silent; @v is a velocity; ( and
	// ) step from the last v
	want := []int{102, 127, 68, 85, 8, 100, 25, 17}
	if !slices.Equal(got, want) {
		t.Errorf("velocities %v, want %v", got, want)
	}
	perf, _ = parse(t, "MML@c,v15c;", nil)
	if perf.Tracks[0].Notes[0].Vel != 68 || perf.Tracks[1].Notes[0].Vel != 127 {
		t.Errorf("Mabinogi velocities %d %d", perf.Tracks[0].Notes[0].Vel, perf.Tracks[1].Notes[0].Vel)
	}
	perf, _ = parse(t, "A c", &Options{Dialect: PPMCK})
	if perf.Tracks[0].Notes[0].Vel != 127 {
		t.Errorf("PPMCK velocity %d", perf.Tracks[0].Notes[0].Vel)
	}
}

func TestTempo(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want []music.Tempo
	}{
		{"c", []music.Tempo{{Tick: 0, BPM: 120}}},
		{"t90 c", []music.Tempo{{Tick: 0, BPM: 90}}},
		{"c t150 c", []music.Tempo{{Tick: 0, BPM: 120}, {Tick: q, BPM: 150}}},
		{"t100 c d t150 e; t100 c d e", []music.Tempo{{Tick: 0, BPM: 100}, {Tick: h, BPM: 150}}},
		{"t100 c t100 c", []music.Tempo{{Tick: 0, BPM: 100}}},
		{"t100 c; t140 c", []music.Tempo{{Tick: 0, BPM: 140}}},
		{"t120.5 c", []music.Tempo{{Tick: 0, BPM: 120.5}}},
		{"t0 c", []music.Tempo{{Tick: 0, BPM: 120}}},
		// a track of tempos only
		{"t80 r1 t160; c1 c1", []music.Tempo{{Tick: 0, BPM: 80}, {Tick: w, BPM: 160}}},
		{"{c t150 d}4", []music.Tempo{{Tick: 0, BPM: 120}, {Tick: e, BPM: 150}}},
	} {
		perf, _ := parse(t, tc.src, &Options{Octave: OctaveNormal})
		if !slices.Equal(perf.Tempo, tc.want) {
			t.Errorf("%q: tempo %v, want %v", tc.src, perf.Tempo, tc.want)
		}
	}
}

func TestTracks(t *testing.T) {
	perf, _ := parse(t, "c d; ; t200; e1 f", &Options{Octave: OctaveNormal})
	if len(perf.Tracks) != 2 {
		t.Fatalf("%d tracks", len(perf.Tracks))
	}
	for i, tr := range perf.Tracks {
		if want := []string{"Track 1", "Track 2"}[i]; tr.Name != want || tr.Channel != i || tr.Program != 80 || tr.Volume != -1 || tr.Pan != -1 {
			t.Errorf("track %d: %q channel %d program %d volume %d pan %d", i, tr.Name, tr.Channel, tr.Program, tr.Volume, tr.Pan)
		}
	}
	if perf.End != w+q {
		t.Errorf("end %d", perf.End)
	}

	// channel 9 is left to percussion
	src := strings.Repeat("c;", 17)
	perf, _ = parse(t, src, nil)
	var chans []int
	for _, tr := range perf.Tracks {
		chans = append(chans, tr.Channel)
	}
	if want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 14, 15, 0, 1}; !slices.Equal(chans, want) {
		t.Errorf("channels %v", chans)
	}

	// the end of the longest track, with its rests
	perf, _ = parse(t, "c r1; d", nil)
	if perf.End != w+q {
		t.Errorf("end %d", perf.End)
	}
}

func TestPrograms(t *testing.T) {
	perf, _ := parse(t, "@1 @p30 @W12 c @W50 d @4 e @5 f; @0c; @2c; @7c; @13c; c", &Options{Octave: OctaveNormal})
	tr := perf.Tracks[0]
	if tr.Program != 81 || tr.Pan != 30 {
		t.Errorf("program %d pan %d", tr.Program, tr.Pan)
	}
	want := []music.Control{
		{Tick: 0, Kind: music.ControlChange, Num: 70, Value: 32},
		{Tick: q, Kind: music.ControlChange, Num: 70, Value: 0},
		{Tick: 2 * q, Kind: music.ControlProgram, Value: 122},
		{Tick: 3 * q, Kind: music.ControlProgram, Value: 80},
	}
	if !slices.Equal(tr.Controls, want) {
		t.Errorf("controls %v, want %v", tr.Controls, want)
	}
	var progs []int
	for _, tr := range perf.Tracks[1:] {
		progs = append(progs, tr.Program)
	}
	if want := []int{79, 82, 122, 80, 80}; !slices.Equal(progs, want) {
		t.Errorf("programs %v, want %v", progs, want)
	}

	// Program for the tracks without a tone
	piano := 0
	perf, _ = parse(t, "c; @1c; c @2 d", &Options{Program: &piano})
	progs = nil
	for _, tr := range perf.Tracks {
		progs = append(progs, tr.Program)
	}
	if want := []int{0, 81, 0}; !slices.Equal(progs, want) {
		t.Errorf("programs %v, want %v", progs, want)
	}
	if c := perf.Tracks[2].Controls; len(c) != 1 || c[0] != (music.Control{Tick: q, Kind: music.ControlProgram, Value: 82}) {
		t.Errorf("controls %v", c)
	}
	bad := 128
	if _, _, err := Parse([]byte("c"), &Options{Program: &bad}); err == nil {
		t.Error("program 128 accepted")
	}
}

func TestMacros(t *testing.T) {
	for _, tc := range []struct {
		src  string
		want []int
	}{
		{"$a=cd; $b=$a e; $b $a", []int{60, 62, 64, 60, 62}},
		{"$ab=e; $a=c; $ab $a", []int{64, 60}},
		{"$a=c; $a $x d", []int{60, 62}},
		{"$m{x}=c%x; $m{4} d", []int{62}},
		{"$a = c\nd ; $a$a", []int{60, 62, 60, 62}},
	} {
		perf, _ := parse(t, tc.src, &Options{Octave: OctaveNormal})
		if got := keys(perf.Tracks[0]); !slices.Equal(got, tc.want) {
			t.Errorf("%q: got %v, want %v", tc.src, got, tc.want)
		}
	}
	// a macro bomb stops at the text budget
	var b strings.Builder
	b.WriteString("$a=cccccccccc;")
	for c := 'b'; c <= 'i'; c++ {
		b.WriteString("$" + string(c) + "=" + strings.Repeat("$"+string(c-1), 10) + ";")
	}
	b.WriteString("d $i")
	start := time.Now()
	_, warnings, err := Parse([]byte(b.String()), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.ContainsFunc(warnings, func(s string) bool { return strings.Contains(s, "macros expand") }) {
		t.Errorf("warnings %q", warnings)
	}
	if d := time.Since(start); d > 20*time.Second {
		t.Errorf("the macro bomb took %v", d)
	}
}

func TestHostileLoops(t *testing.T) {
	for _, src := range []string{
		"[[[[[[c]255]255]255]255]255]255",
		"l64[[[[[[c]255]255]255]255]255]255",
		"c[[[[[[ ]255]255]255]255]255]255",
		"c" + strings.Repeat("[", 100) + "c" + strings.Repeat("]9", 100),
		"/:65535 /:65535 /:65535 c :/ :/ :/",
		"c" + strings.Repeat("{", 100) + "c" + strings.Repeat("}", 100),
	} {
		start := time.Now()
		perf, warnings, err := Parse([]byte(src), nil)
		if err != nil {
			t.Errorf("%.40q: %v", src, err)
			continue
		}
		if len(warnings) == 0 {
			t.Errorf("%.40q: no warning", src)
		}
		count := 0
		for _, tr := range perf.Tracks {
			count += len(tr.Notes)
		}
		if count > maxNotes || perf.End > maxTicks+w {
			t.Errorf("%.40q: %d notes, end %d", src, count, perf.End)
		}
		if d := time.Since(start); d > 20*time.Second {
			t.Errorf("%.40q took %v", src, d)
		}
	}
}

func TestWarnings(t *testing.T) {
	_, warnings := parse(t, "c z z1 z2 j d y", &Options{Octave: OctaveNormal})
	if len(warnings) != 3 {
		t.Errorf("warnings %q", warnings)
	}
	var got []string
	_, warnings, err := Parse([]byte("c z"), &Options{Warn: func(s string) { got = append(got, s) }})
	if err != nil || len(warnings) != 0 || len(got) != 1 {
		t.Errorf("warnings %q, to Warn %q, %v", warnings, got, err)
	}
	if _, _, err := Parse([]byte("r1 /* nothing */"), nil); err == nil {
		t.Error("no error for MML without notes")
	}
	if _, _, err := Parse([]byte("c"), &Options{Dialect: "mucom"}); err == nil {
		t.Error("no error for an unknown dialect")
	}
}

func TestGenericFile(t *testing.T) {
	perf, warnings := parseFile(t, "twinkle.mml")
	if len(warnings) > 0 {
		t.Errorf("warnings %q", warnings)
	}
	if perf.Title != "Twinkle, Twinkle, Little Star" || perf.Composer != "Traditional" || perf.Arranger != "bdf" {
		t.Errorf("credits %q %q %q", perf.Title, perf.Composer, perf.Arranger)
	}
	if !slices.Equal(perf.Tempo, []music.Tempo{{Tick: 0, BPM: 100}}) {
		t.Errorf("tempo %v", perf.Tempo)
	}
	if len(perf.Tracks) != 2 {
		t.Fatalf("%d tracks", len(perf.Tracks))
	}
	mel, bass := perf.Tracks[0], perf.Tracks[1]
	if mel.Program != 80 || bass.Program != 82 || len(mel.Notes) != 42 || len(bass.Notes) != 24 {
		t.Errorf("programs %d %d, notes %d %d", mel.Program, bass.Program, len(mel.Notes), len(bass.Notes))
	}
	if mel.Controls[0] != (music.Control{Kind: music.ControlChange, Num: 70, Value: 64}) {
		t.Errorf("melody controls %v", mel.Controls)
	}
	// < raises the octave here (FlMML)
	if got, want := keys(bass)[:8], []int{48, 60, 53, 60, 53, 60, 55, 48}; !slices.Equal(got, want) {
		t.Errorf("bass %v, want %v", got, want)
	}
	if got, want := keys(mel)[:7], []int{60, 60, 67, 67, 69, 69, 67}; !slices.Equal(got, want) {
		t.Errorf("melody %v, want %v", got, want)
	}
	if perf.End != 12*w || mel.Notes[len(mel.Notes)-1].Tick+mel.Notes[len(mel.Notes)-1].Dur != 12*w {
		t.Errorf("end %d", perf.End)
	}
}

func TestMabinogiFile(t *testing.T) {
	perf, warnings := parseFile(t, "ode.mml")
	if len(warnings) > 0 {
		t.Errorf("warnings %q", warnings)
	}
	var names []string
	for _, tr := range perf.Tracks {
		names = append(names, tr.Name)
		if tr.Program != 0 {
			t.Errorf("%s: program %d", tr.Name, tr.Program)
		}
	}
	if !slices.Equal(names, []string{"Melody", "Chord 1", "Chord 2"}) {
		t.Errorf("tracks %q", names)
	}
	mel := perf.Tracks[0]
	if got, want := notes(mel)[:3], []n{{0, q, 76, 0}, {q, q, 76, 0}, {2 * q, q, 77, 0}}; !slices.Equal(got, want) {
		t.Errorf("melody %v", got)
	}
	// e. d8 d2
	if got, want := notes(mel)[12:15], []n{{12 * q, q + e, 76, 0}, {13*q + e, e, 74, 0}, {14 * q, h, 74, 0}}; !slices.Equal(got, want) {
		t.Errorf("melody %v", got)
	}
	if perf.End != 8*w {
		t.Errorf("end %d", perf.End)
	}

	perf, _ = parse(t, "Song\nMML@t150cde,,e;\nMML@efg,c;", nil)
	names = nil
	for _, tr := range perf.Tracks {
		names = append(names, tr.Name)
	}
	if !slices.Equal(names, []string{"Melody (1)", "Chord 2 (1)", "Melody (2)", "Chord 1 (2)"}) {
		t.Errorf("tracks %q", names)
	}
	if !slices.Equal(perf.Tempo, []music.Tempo{{Tick: 0, BPM: 150}}) {
		t.Errorf("tempo %v", perf.Tempo)
	}
	// an e-mail address is not MML@
	if d := detectDialect("write to mml@example.com\nc d e"); d != Generic {
		t.Errorf("dialect %s", d)
	}
}

func TestPPMCKFile(t *testing.T) {
	perf, warnings := parseFile(t, "frere.mml")
	if len(warnings) != 0 {
		t.Errorf("warnings %q", warnings)
	}
	if perf.Title != "Frere Jacques" || perf.Composer != "Traditional" || perf.Copyright != "Public domain" || perf.Arranger != "bdf" {
		t.Errorf("credits %q %q %q %q", perf.Title, perf.Composer, perf.Copyright, perf.Arranger)
	}
	if !slices.Equal(perf.Tempo, []music.Tempo{{Tick: 0, BPM: 112}}) {
		t.Errorf("tempo %v", perf.Tempo)
	}
	type tk struct {
		name          string
		channel, prog int
		notes         int
	}
	var got []tk
	for _, tr := range perf.Tracks {
		got = append(got, tk{tr.Name, tr.Channel, tr.Program, len(tr.Notes)})
	}
	if want := []tk{{"A", 0, 80, 32}, {"B", 1, 80, 32}, {"C", 2, 82, 16}, {"D", 3, 122, 16}}; !slices.Equal(got, want) {
		t.Errorf("tracks %v, want %v", got, want)
	}
	a, b := perf.Tracks[0], perf.Tracks[1]
	if want := []music.Control{{Kind: music.ControlChange, Num: 70, Value: 0}}; !slices.Equal(a.Controls, want) {
		t.Errorf("A controls %v", a.Controls)
	}
	if want := []music.Control{{Kind: music.ControlChange, Num: 70, Value: 64}}; !slices.Equal(b.Controls, want) {
		t.Errorf("B controls %v", b.Controls)
	}
	if b.Notes[0].Tick != 2*w || a.Notes[0].Vel != 102 || b.Notes[0].Vel != 85 {
		t.Errorf("B starts at %d; velocities %d %d", b.Notes[0].Tick, a.Notes[0].Vel, b.Notes[0].Vel)
	}
	// c<g>c2: > raises
	if got, want := keys(a)[26:], []int{60, 55, 60, 60, 55, 60}; !slices.Equal(got, want) {
		t.Errorf("A ends %v, want %v", got, want)
	}
	if perf.End != 10*w {
		t.Errorf("end %d", perf.End)
	}

	// the DPCM track is not played, a line without a track name is skipped,
	// the default duty is 12.5%
	perf, warnings = parse(t, "#TITLE x\nA c\n  d\nE c\ncde\nZ e", &Options{Dialect: PPMCK})
	var names []string
	for _, tr := range perf.Tracks {
		names = append(names, tr.Name)
	}
	if !slices.Equal(names, []string{"A", "Z"}) || len(warnings) != 2 {
		t.Errorf("tracks %q, warnings %q", names, warnings)
	}
	if got := keys(perf.Tracks[0]); !slices.Equal(got, []int{60, 62}) {
		t.Errorf("A %v", got)
	}
	if c := perf.Tracks[0].Controls; len(c) != 1 || c[0].Num != 70 || c[0].Value != 32 {
		t.Errorf("A controls %v", c)
	}
	if tr := perf.Tracks[1]; tr.Program != 80 || len(tr.Controls) != 0 {
		t.Errorf("Z program %d controls %v", tr.Program, tr.Controls)
	}
	// the duty changes, and is left out on another program
	perf, _ = parse(t, "A @3 c @1 d", &Options{Dialect: PPMCK})
	if c := perf.Tracks[0].Controls; len(c) != 2 || c[0].Value != 96 || c[1] != (music.Control{Tick: q, Kind: music.ControlChange, Num: 70, Value: 64}) {
		t.Errorf("controls %v", c)
	}
	organ := 19
	perf, _ = parse(t, "A @3 c\nD c", &Options{Dialect: PPMCK, Program: &organ})
	if tr := perf.Tracks[0]; tr.Program != 19 || len(tr.Controls) != 0 || perf.Tracks[1].Program != 19 {
		t.Errorf("program %d controls %v", tr.Program, tr.Controls)
	}
	// definitions over several lines, and without braces
	perf, _ = parse(t, "#TITLE x\n@v0 =\n{ 1 2\n 3 }\n@N0 = 5\nA c\n@EN0 = { 1 }\nB d", nil)
	if len(perf.Tracks) != 2 || keys(perf.Tracks[1])[0] != 62 {
		t.Errorf("tracks %v", perf.Tracks)
	}
}

func TestDialect(t *testing.T) {
	for _, tc := range []struct{ src, want string }{
		{"MML@cde;", Mabinogi},
		{"t120 l8 cdefgab;", Generic},
		{"#TITLE x\n#ARTIST y\nt120 cde;", Generic},
		{"#TITLE x\nA t120 cde\nB c4 de", PPMCK},
		{"#OCTAVE-REV\nA cde", PPMCK},
		{"@v0 = { 15 }\nA cde", PPMCK},
		{"A t120 cde", Generic},
		{"#ARTIST y\nA t120 cde\nB c4 d4", Generic},
		{"CDEF GAB;", Generic},
	} {
		if got := detectDialect(tc.src); got != tc.want {
			t.Errorf("%q: %s, want %s", tc.src, got, tc.want)
		}
	}
}

func TestEncodings(t *testing.T) {
	perf, _ := parseFile(t, "sjis.mml")
	if perf.Title != "きらきら星" || perf.Composer != "作者不詳" || len(perf.Tracks[0].Notes) != 14 {
		t.Errorf("Shift_JIS: %q %q", perf.Title, perf.Composer)
	}
	perf, _ = parse(t, "\xEF\xBB\xBF#TITLE Song\r\nc\r\n", nil)
	if perf.Title != "Song" {
		t.Errorf("UTF-8 with BOM: %q", perf.Title)
	}
	u16, err := xunicode.UTF16(xunicode.LittleEndian, xunicode.UseBOM).NewEncoder().Bytes([]byte("#TITLE 歌\nc d"))
	if err != nil {
		t.Fatal(err)
	}
	perf, _, err = Parse(u16, nil)
	if err != nil || perf.Title != "歌" || len(perf.Tracks[0].Notes) != 2 {
		t.Errorf("UTF-16: %v", err)
	}
	// full-width spaces are spaces
	perf, warnings := parse(t, "c\u3000d", nil)
	if len(warnings) != 0 || len(perf.Tracks[0].Notes) != 2 {
		t.Errorf("warnings %q", warnings)
	}
}

func TestDetect(t *testing.T) {
	detect := func(s string) bool { return Detect([]byte(s), strings.NewReader(s), int64(len(s))) }
	files, _ := filepath.Glob("testdata/*.mml")
	for _, name := range files {
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if !Detect(b[:min(len(b), 1024)], nil, int64(len(b))) {
			t.Errorf("%s is not detected", name)
		}
	}
	for _, s := range []string{
		"MML@t120l8cdef,o3c1,o4e1;",
		"\xEF\xBB\xBFmml@t150l4o5eefg;",
		"#TITLE Song\n#COMPOSER Me\n\nA t150 o4 l8 cdef\n",
		"A t150 o4 l8 cdef\nB o3 l4 cegc\n",
		"ABC t150 l8\nA o4 cdefgab\n",
		"#TITLE Twinkle\n#ARTIST Traditional\nt100 l4 o4 ccggaag2;\n",
		"#TITLE x\n@v0 = { 15 14 13 }\n",
		"#OCTAVE REVERSE\n@1 l8 o5 c<c>c;",
	} {
		if !detect(s) {
			t.Errorf("%q is not detected", s)
		}
	}
	for _, s := range []string{
		"",
		"cdefgab>c",
		"t120 l8 o4 cdefgab>c;\n",
		"A bad face.\nI read a book about music.\nThe cafe was dead.\n",
		"A cab fed a dog.\nB bad deed.\n",
		"name,age\nbob,30\nann,25\n",
		"a1,b2,c3\nd4,e5,f6\n",
		"A 1,B 2\nC 3,D 4\n",
		"#include <stdio.h>\nint main(void) { return 0; }\n",
		"#!/bin/sh\n# TITLE\necho cde\n",
		"# Title\n\nSome text about MML: `cdefg`.\n\n```\n#TITLE Song\nt120 l8 cdefgab\n```\n",
		"```mml\nMML@t120l8cdef;\n```\n",
		"IN;SP1;PU0,0;PD100,100;\n",
		"G04 comment*\n%FSLAX24Y24*%\n",
		"M48\n;LAYER\nMETRIC\nT1C0.8\n%\nT1\nX100Y200\n",
		"#TITLE A song\nI love to sing it.\n",
		"<html><body>cde</body></html>",
		"\x00MML@cde;",
		"LDA #$00\nSTA $2000\n",
	} {
		if detect(s) {
			t.Errorf("%q is detected", s)
		}
	}
}

func TestParams(t *testing.T) {
	f := conv.Lookup("mml")
	if f == nil || len(f.Params) != 6 || !slices.Equal(f.Extensions, []string{".mml"}) {
		t.Fatalf("format %+v", f)
	}
	for _, tc := range []struct {
		params map[string]string
		ok     bool
	}{
		{nil, true},
		{map[string]string{"dialect": "ppmck", "octave": "reverse", "program": "0", "time": "3/4", "key": "F#m"}, true},
		{map[string]string{"dialect": "mucom"}, false},
		{map[string]string{"octave": "up"}, false},
		{map[string]string{"program": "128"}, false},
		{map[string]string{"program": "x"}, false},
		{map[string]string{"time": "3/5"}, false},
		{map[string]string{"key": "H"}, false},
	} {
		opts, no, err := options(&conv.Options{Params: tc.params})
		if (err == nil) != tc.ok {
			t.Errorf("%v: %v", tc.params, err)
			continue
		}
		if tc.params["dialect"] == "ppmck" && (opts.Dialect != PPMCK || opts.Octave != OctaveReverse || *opts.Program != 0 ||
			*no.Time != (music.TimeSig{Beats: 3, BeatType: 4}) || *no.Key != (music.KeySig{Fifths: 3, Minor: true})) {
			t.Errorf("%v: %+v %+v", tc.params, opts, no)
		}
	}
}

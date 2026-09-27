package music

import (
	"bytes"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestBuildPlay(t *testing.T) {
	p := samples()["twinkle"]
	res, err := Build(Notate(p, NotateOptions{}), Options{Source: "test", FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	doc := res.Doc
	v := doc.Views[0]
	if v.Play == nil || v.Play.Seq == "" || v.Play.Cues == "" {
		t.Fatalf("play %+v", v.Play)
	}
	seq := doc.Part(mustHash(t, v.Play.Seq))
	if seq == nil || seq.Type != bdf.PartSeq || !bytes.HasPrefix(seq.Data, []byte("MThd")) {
		t.Fatalf("seq part %+v", seq)
	}
	cp := doc.Part(mustHash(t, v.Play.Cues))
	cues, err := bdf.DecodeCues(cp.Data)
	if err != nil {
		t.Fatal(err)
	}
	if len(cues.Systems) != res.Systems {
		t.Errorf("%d cue systems, %d systems", len(cues.Systems), res.Systems)
	}
	var last uint32
	for i, c := range cues.Cues {
		if c.Tick < last {
			t.Fatalf("cue %d goes back: %d < %d", i, c.Tick, last)
		}
		last = c.Tick
		s := cues.Systems[c.System]
		if c.X < s.X-1 || c.X > s.X+s.W+1 {
			t.Errorf("cue %d at x %g outside its system %+v", i, c.X, s)
		}
	}
	// 28 notes in 8 measures of 4/4: the last cue is the end of the music
	if want := uint32(8 * 4 * PPQ); last != want {
		t.Errorf("last cue at %d, want %d", last, want)
	}
	if res.Seconds < 19 || res.Seconds > 20 {
		t.Errorf("%g seconds", res.Seconds)
	}
}

func TestBuildPlayOrder(t *testing.T) {
	// a repeat: measures 0 1 0 1 2
	s := Notate(melody(60, 4*q, 62, 4*q, 64, 4*q), NotateOptions{})
	s.PlayOrder = []PlayedMeasure{{0, 0}, {4 * q, 1}, {8 * q, 0}, {12 * q, 1}, {16 * q, 2}}
	s.Play.Division = 480
	res, err := Build(s, Options{Source: "test", FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	v := res.Doc.Views[0]
	cues, err := bdf.DecodeCues(res.Doc.Part(mustHash(t, v.Play.Cues)).Data)
	if err != nil {
		t.Fatal(err)
	}
	var ticks []uint32
	for _, c := range cues.Cues {
		ticks = append(ticks, c.Tick)
	}
	// each measure: its note and its end, in the file's 480 ticks a quarter
	want := []uint32{0, 1920, 1920, 3840, 3840, 5760, 5760, 7680, 7680, 9600}
	if len(ticks) != len(want) {
		t.Fatalf("ticks %v, want %v", ticks, want)
	}
	for i := range want {
		if ticks[i] != want[i] {
			t.Fatalf("ticks %v, want %v", ticks, want)
		}
	}
	// the third cue returns to the first measure's place
	if cues.Cues[4].X != cues.Cues[0].X {
		t.Errorf("repeat at x %g, first at %g", cues.Cues[4].X, cues.Cues[0].X)
	}
}

func mustHash(t *testing.T, s string) bdf.Hash {
	t.Helper()
	h, err := bdf.ParseHash(s)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

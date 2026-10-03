package music

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf/raster/imagebdf"
)

func TestTabCandidatesAndAuthoredPosition(t *testing.T) {
	note := &Note{Pitch: Pitch{Step: 2, Octave: 4}} // E4
	got := tabCandidates(tabInput{note: note, pitch: 64}, false)
	want := []TabPosition{{1, 0}, {2, 5}, {3, 9}, {4, 14}, {5, 19}, {6, 24}}
	if len(got) != len(want) {
		t.Fatalf("E4 candidates: %+v", got)
	}
	for i := range want {
		if got[i].pos != want[i] {
			t.Fatalf("candidate %d: %+v, want %+v", i, got[i].pos, want[i])
		}
	}
	// An authored fingering wins even when the written pitch is octave
	// transposed. Alternatives remain available if it collides in a chord.
	note.Tab = &TabPosition{String: 2, Fret: 5}
	chosen, omitted := chooseTab([]tabOnset{{0, []tabInput{{note: note, pitch: 64}}}}, tabPresets[0])
	if omitted != 0 || chosen[note] != *note.Tab {
		t.Fatalf("authored position: %+v, omitted %d", chosen[note], omitted)
	}
}

func TestTabChordConstraints(t *testing.T) {
	a := &Note{Tab: &TabPosition{String: 1, Fret: 0}}
	b := &Note{Tab: &TabPosition{String: 1, Fret: 0}}
	shapes := tabShapes([]tabInput{{note: a, pitch: 64}, {note: b, pitch: 64}}, tabPresets[0], true)
	if len(shapes) == 0 {
		t.Fatal("two E4 notes can use different strings")
	}
	for _, shape := range shapes {
		if shape.positions[0].String == shape.positions[1].String {
			t.Fatalf("two notes on the same string: %+v", shape)
		}
	}
	chosen, omitted := chooseTab([]tabOnset{{0, []tabInput{{note: a, pitch: 64}, {note: b, pitch: 64}}}}, tabPresets[0])
	if omitted != 0 || chosen[a].String == chosen[b].String {
		t.Fatalf("authored-string collision was not resolved: %+v %+v, omitted %d", chosen[a], chosen[b], omitted)
	}
	// Fixed authored frets twenty frets apart cannot be held as one chord.
	c := &Note{Tab: &TabPosition{String: 6, Fret: 1}}
	d := &Note{Tab: &TabPosition{String: 1, Fret: 20}}
	for _, shape := range tabShapes([]tabInput{{note: c, pitch: 41}, {note: d, pitch: 84}}, tabPresets[0], true) {
		if shape.positions[0].Fret == 1 && shape.positions[1].Fret == 20 {
			t.Fatalf("unplayable stretch: %+v", shape)
		}
	}
	seven := make([]tabInput, 7)
	for i := range seven {
		seven[i] = tabInput{note: &Note{}, pitch: 64}
	}
	if got := tabShapes(seven, tabPresets[0], false); len(got) != 0 {
		t.Fatalf("seven-note chord: %+v", got)
	}
}

func TestTabPhraseAndViews(t *testing.T) {
	held := &Note{Tab: &TabPosition{String: 1, Fret: 0}}
	overlap := &Note{}
	overlapped, omitted := chooseTab([]tabOnset{
		{0, []tabInput{{note: held, pitch: 64, end: 2 * PPQ}}},
		{PPQ, []tabInput{{note: overlap, pitch: 64, end: 3 * PPQ}}},
	}, tabPresets[0])
	if omitted != 0 || overlapped[held].String != 1 || overlapped[overlap].String == 1 {
		t.Fatalf("overlapping notes reuse a string: %+v %+v, omitted %d", overlapped[held], overlapped[overlap], omitted)
	}
	first, second := &Note{}, &Note{}
	position, _ := chooseTab([]tabOnset{{0, []tabInput{{note: first, pitch: 64}}},
		{PPQ, []tabInput{{note: second, pitch: 69}}}}, tabPresets[1])
	if position[first] != (TabPosition{String: 2, Fret: 5}) || position[second] != (TabPosition{String: 1, Fret: 5}) {
		t.Fatalf("phrase optimization did not keep the hand near fret 5: %+v %+v", position[first], position[second])
	}
	var onsets []tabOnset
	for i, pitch := range []int{64, 66, 67, 69, 67, 66, 64} {
		onsets = append(onsets, tabOnset{tick: i * PPQ, notes: []tabInput{{note: &Note{}, pitch: pitch}}})
	}
	chosen, omitted := chooseTab(onsets, tabPresets[1])
	if omitted != 0 || len(chosen) != len(onsets) {
		t.Fatalf("phrase: %d choices, %d omitted", len(chosen), omitted)
	}
	res, err := Build(Notate(samples()["twinkle"], NotateOptions{}), Options{
		Source: "test", GuitarTAB: true, FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Doc.Views) != 5 || res.Doc.Views[0].ID != "score" {
		t.Fatalf("views: %+v", res.Doc.Views)
	}
	scorePlay := res.Doc.Views[0].Play
	if scorePlay == nil {
		t.Fatal("score is missing playback data")
	}
	for _, view := range res.Doc.Views[1:] {
		if !strings.Contains(view.Title, "estimated") || len(view.Pages) == 0 || view.Play == nil || view.Play.Cues == "" {
			t.Fatalf("incomplete TAB view: %+v", view)
		}
		if view.Play.Seq != scorePlay.Seq {
			t.Fatalf("TAB view %q has different playback data from the score", view.ID)
		}
	}
	if h := res.Doc.Views[1].Pages[0].H; h > 350 {
		t.Fatalf("short TAB has an almost empty page: height %g", h)
	}
	r := imagebdf.New(res.Doc, &imagebdf.Options{FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
	if _, err := r.Page(res.Doc.Views[1], 0, 1); err != nil {
		t.Fatalf("render TAB page: %v", err)
	}
}

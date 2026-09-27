package music

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster"
)

// sb builds a score by hand, measure by measure.
type sb struct {
	s    *Score
	part *Part
}

func newSB(title string) *sb {
	return &sb{s: &Score{Title: title, Composer: "Test"}}
}

func (b *sb) addPart(name string, staves int, clefs ...Clef) {
	p := &Part{Name: name, Staves: staves}
	b.s.Parts = append(b.s.Parts, p)
	b.part = p
	for len(p.Measures) < len(b.s.Measures) {
		p.Measures = append(p.Measures, &PartMeasure{})
	}
	if len(p.Measures) == 0 {
		return
	}
	for i, c := range clefs {
		p.Measures[0].Clefs = append(p.Measures[0].Clefs, ClefChange{Staff: i, Clef: c})
	}
}

func (b *sb) measure(m *Measure) int {
	if m.Length == 0 {
		m.Length = 4 * PPQ
	}
	b.s.Measures = append(b.s.Measures, m)
	for _, p := range b.s.Parts {
		p.Measures = append(p.Measures, &PartMeasure{})
	}
	return len(b.s.Measures) - 1
}

// ev adds an event: notes like "c5", "f#4", "bb3"; "r" for a rest.
func (b *sb) ev(pi, mi, off int, typ NoteType, dots int, notes ...string) *Event {
	e := &Event{Offset: off, Type: typ, Dots: dots, Voice: 1, Duration: typ.Ticks(dots)}
	for _, n := range notes {
		if n == "r" {
			e.Rest = true
			continue
		}
		e.Notes = append(e.Notes, &Note{Pitch: parsePitch(n)})
	}
	pm := b.s.Parts[pi].Measures[mi]
	pm.Events = append(pm.Events, e)
	return e
}

func parsePitch(s string) Pitch {
	p := Pitch{Step: map[byte]int{'c': 0, 'd': 1, 'e': 2, 'f': 3, 'g': 4, 'a': 5, 'b': 6}[s[0]]}
	i := 1
	for i < len(s) && (s[i] == '#' || s[i] == 'b') {
		if s[i] == '#' {
			p.Alter++
		} else {
			p.Alter--
		}
		i++
	}
	p.Octave = int(s[i] - '0')
	return p
}

func featureScore() *Score {
	b := newSB("Features")
	b.s.Subtitle = "every mark the engraver draws"
	b.s.Lyricist = "Words: Anonymous"
	m0 := b.measure(&Measure{Time: &TimeSig{4, 4, ""}, Tempo: []TempoMark{{Text: "Allegro", Beat: Quarter, PerMinute: 132}},
		Left: BarRepeatStart, Rehearsal: "A"})
	m1 := b.measure(&Measure{})
	m2 := b.measure(&Measure{Ending: &Ending{Text: "1.", Start: true, Stop: true}, Right: BarRepeatEnd})
	m3 := b.measure(&Measure{Ending: &Ending{Text: "2.", Start: true, Stop: true, Open: true}, Right: BarDouble})
	m4 := b.measure(&Measure{Time: &TimeSig{3, 4, ""}, Length: 3 * PPQ})
	m5 := b.measure(&Measure{Time: &TimeSig{6, 8, ""}, Length: 3 * PPQ})
	m6 := b.measure(&Measure{Time: &TimeSig{4, 4, "common"}, Right: BarFinal})
	b.addPart("Voice", 1, Clef{Sign: ClefG})
	b.addPart("Piano", 2, Clef{Sign: ClefG}, Clef{Sign: ClefF})
	for _, p := range b.s.Parts {
		p.Measures[0].Key = &KeySig{Fifths: -3}
		p.Measures[m4].Key = &KeySig{Fifths: 2}
	}
	b.s.Parts[1].Measures[m0].Clefs = append(b.s.Parts[1].Measures[m0].Clefs, ClefChange{Offset: 2 * PPQ, Staff: 1, Clef: Clef{Sign: ClefG}})
	b.s.Parts[1].Measures[m1].Clefs = append(b.s.Parts[1].Measures[m1].Clefs, ClefChange{Offset: 0, Staff: 1, Clef: Clef{Sign: ClefF}})

	// the voice: slurs, articulations, lyrics in two verses, grace notes
	v := func(mi, off int, t NoteType, dots int, n string, words ...string) *Event {
		e := b.ev(0, mi, off, t, dots, n)
		for i, w := range words {
			syl := SyllableSingle
			if len(w) > 0 && w[len(w)-1] == '-' {
				w, syl = w[:len(w)-1], SyllableBegin
			}
			e.Lyrics = append(e.Lyrics, &LyricSyllable{Verse: i + 1, Text: w, Syllabic: syl})
		}
		return e
	}
	e := v(m0, 0, Quarter, 0, "eb5", "Hap-", "Sing-")
	e.Slurs = []SlurMark{{Number: 1, Start: true}}
	e.Artics = []string{"articAccent"}
	v(m0, PPQ, Eighth, 1, "d5", "py", "ing")
	v(m0, PPQ+3*PPQ/4, N16th, 0, "c5", "", "")
	e = v(m0, 2*PPQ, Half, 0, "bb4", "day", "loud")
	e.Slurs = []SlurMark{{Number: 1}}
	e.Artics = []string{"fermata"}
	e = v(m1, 0, Quarter, 0, "ab4", "to", "all")
	e.Artics = []string{"articStaccato", "articTenuto"}
	g := b.ev(0, m1, PPQ, Eighth, 0, "c5")
	g.Grace, g.Slash, g.Duration = true, true, 0
	tup := &Tuplet{Actual: 3, Normal: 2, ShowNumber: true}
	for i, n := range []string{"bb4", "c5", "db5"} {
		e := v(m1, PPQ+i*PPQ/3, Eighth, 0, n, "you", "the", "day")
		e.Tuplet, e.Duration = tup, PPQ/3
	}
	e = v(m1, 2*PPQ, Half, 0, "eb5", "dear-", "long-")
	e.Notes[0].TieStart = true
	e.Lyrics[0].Extend = true
	e = v(m2, 0, Whole, 0, "eb5")
	e.Notes[0].TieStop = true
	e.Artics = []string{"ornamentTrill"}
	v(m3, 0, Half, 0, "r")
	v(m3, 2*PPQ, Half, 0, "c#5", "end.")
	v(m4, 0, Quarter, 0, "d5")
	v(m4, PPQ, Quarter, 0, "f#5")
	v(m4, 2*PPQ, Quarter, 0, "a5")
	for i, n := range []string{"b5", "a5", "g5", "f#5", "e5", "d5"} {
		v(m5, i*PPQ/2, Eighth, 0, n)
	}
	e = v(m6, 0, Whole, 0, "d5")
	e.Artics = []string{"fermata"}

	// the piano: chords with seconds, two voices, dynamics, a hairpin,
	// pedal, an octave line and chord symbols
	pm := func(mi int) *PartMeasure { return b.s.Parts[1].Measures[mi] }
	b.ev(1, m0, 0, Half, 0, "eb4", "g4", "bb4", "c5")
	b.ev(1, m0, 2*PPQ, Quarter, 0, "f4", "g4", "ab4")
	b.ev(1, m0, 3*PPQ, Quarter, 0, "r")
	for i, n := range []string{"eb3", "bb3", "g3", "bb3"} {
		e := b.ev(1, m0, i*PPQ, Quarter, 0, n)
		e.Staff = 1
	}
	pm(m0).Directions = []*Direction{
		{Offset: 0, Kind: DirDynamic, Text: "mf"},
		{Offset: PPQ, Kind: DirCrescendo, Number: 1},
		{Offset: 3 * PPQ, Kind: DirCrescendo, Number: 1, Stop: true},
		{Offset: 0, Kind: DirChord, Text: "Cm7"},
		{Offset: 2 * PPQ, Kind: DirChord, Text: "Fm/Ab"},
		{Offset: 0, Staff: 1, Kind: DirPedal},
		{Offset: 3 * PPQ, Staff: 1, Kind: DirPedal, Stop: true},
	}
	// two voices in the right hand
	for i, n := range []string{"eb5", "f5", "g5", "ab5", "bb5", "ab5", "g5", "f5"} {
		e := b.ev(1, m1, i*PPQ/2, Eighth, 0, n)
		e.Voice = 1
	}
	e = b.ev(1, m1, 0, Whole, 0, "g4")
	e.Voice = 2
	e = b.ev(1, m1, 0, Whole, 0, "eb3", "bb3")
	e.Staff = 1
	pm(m1).Directions = []*Direction{{Offset: 0, Kind: DirOctave, Size: 8, Number: 1}, {Offset: 3 * PPQ, Kind: DirOctave, Stop: true, Number: 1},
		{Offset: 0, Kind: DirDynamic, Text: "p"}, {Offset: 2 * PPQ, Kind: DirWords, Text: "dolce", Italic: true}}
	for _, mi := range []int{m2, m3} {
		b.ev(1, mi, 0, Whole, 0, "g4", "bb4", "eb5")
		e := b.ev(1, mi, 0, Whole, 0, "eb3")
		e.Staff = 1
	}
	for i, n := range []string{"d4", "e4", "f#4", "g4", "a4", "b4"} {
		b.ev(1, m4, i*PPQ/2, Eighth, 0, n)
	}
	for i, n := range []string{"d3", "a3", "f#3"} {
		e := b.ev(1, m4, i*PPQ, Quarter, 0, n)
		e.Staff = 1
	}
	b.ev(1, m5, 0, Quarter, 1, "d4", "f#4", "a4")
	b.ev(1, m5, 3*PPQ/2, Quarter, 1, "c#4", "e4", "a4")
	e = b.ev(1, m5, 0, Half, 1, "d3")
	e.Staff = 1
	b.ev(1, m6, 0, Whole, 0, "d4", "f#4", "a4", "d5")
	e = b.ev(1, m6, 0, Whole, 0, "d2", "d3")
	e.Staff = 1
	pm(m6).Directions = []*Direction{{Offset: 0, Kind: DirDynamic, Text: "pp"}}
	return b.s
}

func TestFeatures(t *testing.T) {
	res, err := Build(featureScore(), Options{Source: "test", FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	dir := os.Getenv("MUSIC_SAMPLES")
	if dir == "" {
		return
	}
	writePages(t, res.Doc, filepath.Join(dir, "features"))
	f, err := os.Create(filepath.Join(dir, "features.bdf"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := res.Doc.WriteSingle(f); err != nil {
		t.Fatal(err)
	}
}

func writePages(t *testing.T, doc *bdf.Document, prefix string) {
	r := raster.New(doc, &raster.Options{FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
	v := doc.Views[0]
	for i := range v.Pages {
		img, err := r.Page(v, i, 2)
		if err != nil {
			t.Fatal(err)
		}
		f, err := os.Create(prefix + "-" + string(rune('1'+i)) + ".png")
		if err != nil {
			t.Fatal(err)
		}
		png.Encode(f, img)
		f.Close()
	}
}

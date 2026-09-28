package musicxml

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/music"
	"github.com/shibukawa/bdf/converter/midi"
)

func warned(ws []string, part string) bool {
	return slices.ContainsFunc(ws, func(w string) bool { return strings.Contains(w, part) })
}

// engrave builds a score without the fonts of the system.
func engrave(t *testing.T, s *music.Score) *music.Result {
	t.Helper()
	res, err := music.Build(s, music.Options{Source: "musicxml", FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// parts makes a partwise score of parts of 16 staves; the first has the
// measures.
func parts(n int, measures ...string) string {
	var b strings.Builder
	b.WriteString(`<score-partwise><part-list>`)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, `<score-part id="P%d"><part-name>Music</part-name></score-part>`, i)
	}
	b.WriteString(`</part-list>`)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, `<part id="P%d"><measure number="1"><attributes><divisions>1</divisions><staves>16</staves></attributes>`, i)
		if i == 0 {
			b.WriteString(strings.Join(measures, `</measure><measure>`))
		}
		b.WriteString(`</measure></part>`)
	}
	b.WriteString(`</score-partwise>`)
	return b.String()
}

func TestStaffMeasuresBudget(t *testing.T) {
	// 256 staves of 3907 measures are more than a million measures of staves
	doc := parts(16, make([]string, 3907)...)
	s, w, err := Parse(strings.NewReader(doc), int64(len(doc)), &Options{NoPlay: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Measures) != 3906 || len(s.Parts) != 16 || len(s.Parts[15].Measures) != 3906 {
		t.Errorf("%d measures, %d parts", len(s.Measures), len(s.Parts))
	}
	if !warned(w, "more than 1000000 measures on the staves (256 staves of 3907 measures); the measures after 3906 are left out") {
		t.Errorf("warnings %q", w)
	}

	// at the limit and past it: what is left is engraved and played
	defer func(n int) { maxStaffMeasures = n }(maxStaffMeasures)
	maxStaffMeasures = 96
	whole := n("C5", 4, "whole")
	s, w = parse(t, parts(2, whole, whole, whole))
	if res := engrave(t, s); len(w) != 0 || res.Measures != 3 || res.Staves != 32 {
		t.Errorf("3 measures: %s, warnings %q", res.Summary(), w)
	}
	s, w = parse(t, parts(2, whole, whole, whole, whole+`<barline><repeat direction="backward"/></barline>`))
	if res := engrave(t, s); !warned(w, "the measures after 3 are left out") || res.Measures != 3 || s.Play.End != 3*4*music.PPQ {
		t.Errorf("4 measures: %s, played to %d, warnings %q", res.Summary(), s.Play.End, w)
	}

	// the measures of the parts count as they are read: those of a part
	// that starts late in a timewise score too
	maxStaffMeasures = 12
	var b strings.Builder
	b.WriteString(`<score-timewise>`)
	for i := 0; i < 10; i++ {
		b.WriteString(`<measure><part id="A">` + whole + `</part>`)
		if i >= 8 {
			b.WriteString(`<part id="B">` + whole + `</part>`)
		}
		b.WriteString(`</measure>`)
	}
	b.WriteString(`</score-timewise>`)
	s, w = parse(t, b.String())
	if !warned(w, "more than 12 measures in the parts") || !warned(w, "the measures after 6 are left out") {
		t.Errorf("warnings %q", w)
	}
	if len(s.Measures) != 6 || len(s.Parts) != 2 || len(events(s, 0, 5)) != 1 || len(events(s, 1, 5)) != 0 {
		t.Errorf("%d measures, %d parts", len(s.Measures), len(s.Parts))
	}
}

func TestMeasureBudget(t *testing.T) {
	rest := `<note><rest/><duration>1</duration></note>`
	slurred := `<note><pitch><step>C</step><octave>5</octave></pitch><duration>1</duration>` +
		`<notations><slur type="start"/><slur type="stop" number="2"/></notations><lyric><text>la</text></lyric><lyric number="2"><text>lu</text></lyric></note>`
	words := `<direction><direction-type><words>dolce</words></direction-type></direction>`
	clef := `<attributes><clef><sign>F</sign></clef></attributes>`

	// the measure takes 10000 notes
	s, w := parse(t, score(attrs(1)+strings.Repeat(rest, maxInMeasure+1), rest))
	if got := len(events(s, 0, 0)); got != maxInMeasure || len(events(s, 0, 1)) != 1 {
		t.Errorf("%d notes", got)
	}
	if !warned(w, "more than 10000 notes, slurs, syllables and directions in a measure of a part") {
		t.Errorf("warnings %q", w)
	}
	// slurs, syllables, directions and clef changes count: a note with
	// its slurs and its first syllable fills the measure
	s, w = parse(t, score(attrs(1)+strings.Repeat(rest, maxInMeasure-4)+slurred+words+clef+rest))
	evs := events(s, 0, 0)
	pm := s.Parts[0].Measures[0]
	if last := evs[len(evs)-1]; len(evs) != maxInMeasure-3 || len(last.Slurs) != 2 || len(last.Lyrics) != 1 || len(pm.Directions) != 0 || len(pm.Clefs) != 1 {
		t.Errorf("%d notes, the last with %d slurs and %d syllables; %d directions, %d clefs", len(evs), len(last.Slurs), len(last.Lyrics),
			len(pm.Directions), len(pm.Clefs))
	}
	if len(w) != 1 {
		t.Errorf("warnings %q", w)
	}
	// with three rests less all is read
	s, w = parse(t, score(attrs(1)+strings.Repeat(rest, maxInMeasure-7)+slurred+words+clef))
	evs = events(s, 0, 0)
	pm = s.Parts[0].Measures[0]
	if last := evs[len(evs)-1]; len(w) != 0 || len(last.Slurs) != 2 || len(last.Lyrics) != 2 || len(pm.Directions) != 1 || len(pm.Clefs) != 2 {
		t.Errorf("the last note with %d slurs and %d syllables; %d directions, %d clefs; warnings %q", len(last.Slurs), len(last.Lyrics),
			len(pm.Directions), len(pm.Clefs), w)
	}
}

func TestTempoMarksBudget(t *testing.T) {
	var b strings.Builder
	b.WriteString(attrs(64))
	for i := 0; i <= maxTempoMarks; i++ {
		fmt.Fprintf(&b, `<direction><direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>%d</per-minute></metronome></direction-type></direction>`, 60+i)
		b.WriteString(`<forward><duration>1</duration></forward>`)
	}
	// a mark of the measure once more is no new one
	b.WriteString(`<backup><duration>65</duration></backup>` +
		`<direction><direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>60</per-minute></metronome></direction-type></direction>`)
	s, w := parse(t, score(b.String()+n("C5", 256, "whole")))
	if got := len(s.Measures[0].Tempo); got != maxTempoMarks {
		t.Errorf("%d tempo marks", got)
	}
	if len(w) != 1 || !warned(w, "more than 64 tempo marks in a measure") {
		t.Errorf("warnings %q", w)
	}
}

func TestVersesBudget(t *testing.T) {
	var b strings.Builder
	for i := 0; i < maxVerses; i++ {
		fmt.Fprintf(&b, `<lyric number="verse%d"><text>la</text></lyric>`, i)
	}
	// verses by their numbers, and one name more than there are verses
	b.WriteString(`<lyric number="7"><text>seven</text></lyric><lyric number="one more"><text>lu</text></lyric>`)
	s, w := parse(t, score(attrs(1)+n("C5", 4, "whole", b.String())+n("C5", 4, "whole", `<lyric number="verse98"><text>li</text></lyric>`)))
	ls := events(s, 0, 0)[0].Lyrics
	if len(ls) != maxVerses || ls[0].Verse != 1 || ls[maxVerses-1].Verse != maxVerses {
		t.Fatalf("%d syllables", len(ls))
	}
	if last := events(s, 0, 0)[1].Lyrics; len(last) != 1 || last[0].Verse != maxVerses || last[0].Text != "li" {
		t.Errorf("the second note: %+v", last)
	}
	if len(w) != 1 || !warned(w, "more than 99 verses of lyrics") {
		t.Errorf("warnings %q", w)
	}
}

func TestEndingNumbers(t *testing.T) {
	upTo := func(n int) []int {
		out := make([]int, n)
		for i := range out {
			out[i] = i + 1
		}
		return out
	}
	for _, c := range []struct {
		number string
		want   []int
		more   bool
	}{
		{"", nil, false},
		{"1", []int{1}, false},
		{"1, 2", []int{1, 2}, false},
		{"1 3;5.", []int{1, 3, 5}, false},
		{"2-4", []int{2, 3, 4}, false},
		{"0, x, 101, 4-2, 7", []int{7}, false},
		{"1-100", upTo(100), false},
		{"1-100,1", upTo(100), true},
		{"1-99, 100, 100", upTo(100), true},
		{"1-60 61-100 1-100", upTo(100), true},
	} {
		got, more := endingNumbers(c.number)
		if !slices.Equal(got, c.want) || more != c.more {
			t.Errorf("%q: %v %v, want %v %v", c.number, got, more, c.want, c.more)
		}
	}
	// a number of many numbers is read as far as the budget goes
	number := strings.Repeat("1-100,", 1000)
	m := attrs(1) + `<barline location="left"><ending number="` + number + `" type="start"/></barline>` + n("C5", 4, "whole") +
		`<barline location="right"><ending number="` + number + `" type="stop"/><repeat direction="backward"/></barline>`
	s, w := parse(t, score(m, n("C5", 4, "whole")))
	if !warned(w, "more than 100 numbers on a volta") {
		t.Errorf("warnings %q", w)
	}
	if e := s.Measures[0].Ending; e == nil || len(e.Text) > 500 || len(s.PlayOrder) != 101 {
		t.Errorf("volta %+v, %d measures played", e, len(s.PlayOrder))
	}
}

func TestPlayLength(t *testing.T) {
	// measures as long as they may be: 245 of them are longer than the
	// performance may be
	long := `<forward><duration>4092</duration></forward>` + n("C5", 4, "whole")
	ms := slices.Repeat([]string{long}, 245)
	ms[0] = attrs(1) + long
	s, w := parse(t, score(ms...))
	if !warned(w, "plays longer than 1000000 quarter notes") {
		t.Errorf("warnings %q", w)
	}
	if len(s.PlayOrder) != 244 || s.Play.End != 244*4096*music.PPQ || len(s.Play.Tracks[0].Notes) != 244 {
		t.Fatalf("%d measures played to %d", len(s.PlayOrder), s.Play.End)
	}
	// the waits between the notes are longer than a delta time of a
	// Standard MIDI File: the file written reads as the same music
	p, err := midi.Parse(s.Play.SMF(), nil)
	if err != nil || len(p.Warnings) != 0 {
		t.Fatalf("the file written: %v, warnings %q", err, p.Warnings)
	}
	notes := p.Perf.Tracks[0].Notes
	if len(notes) != 244 || notes[243].Tick != 243*4096*music.PPQ+4092*music.PPQ || p.Perf.End != s.Play.End {
		t.Errorf("%d notes to %d", len(notes), p.Perf.End)
	}
	if res := engrave(t, s); res.Measures != 245 {
		t.Errorf("%s", res.Summary())
	}
}

func TestPlayBudgets(t *testing.T) {
	defer func(n int) { maxPlayChanges = n }(maxPlayChanges)
	maxPlayChanges = 10
	// a measure of three tempo changes and four syllables, played 5 times
	var b strings.Builder
	b.WriteString(attrs(1))
	for i := 0; i < 4; i++ {
		fmt.Fprintf(&b, `<sound tempo="%d"/>`, 60+10*(i%2))
		b.WriteString(n("C5", 1, "quarter", `<lyric><text>la</text></lyric>`))
	}
	b.WriteString(`<barline location="right"><repeat direction="backward" times="5"/></barline>`)
	s, w := parse(t, score(b.String()))
	if len(s.PlayOrder) != 5 || len(s.Play.Tempo) != 10 || len(s.Play.Tracks[0].Lyrics) != 10 || len(s.Play.Tracks[0].Notes) != 20 {
		t.Errorf("%d measures played, %d tempo changes, %d syllables", len(s.PlayOrder), len(s.Play.Tempo), len(s.Play.Tracks[0].Lyrics))
	}
	if len(w) != 2 || !warned(w, "more than 10 tempo changes are played") || !warned(w, "more than 10 syllables are sung") {
		t.Errorf("warnings %q", w)
	}
}

func TestClefNone(t *testing.T) {
	// a clef that is not drawn, from the middle of the measure
	m := attrs(1) + n("C5", 2, "half") + `<attributes><clef><sign>none</sign></clef></attributes>` + n("D5", 2, "half")
	s, _ := parse(t, score(m))
	if cs := s.Parts[0].Measures[0].Clefs; len(cs) != 2 || cs[1].Clef.Sign != music.ClefNone || cs[1].Offset != 2*music.PPQ {
		t.Fatalf("clefs %+v", cs)
	}
	if res := engrave(t, s); res.Notes != 2 {
		t.Errorf("%s", res.Summary())
	}
}

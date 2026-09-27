package musicxml

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/music"
)

// parse reads a document, failing the test on an error.
func parse(t *testing.T, doc string) (*music.Score, []string) {
	t.Helper()
	s, w, err := Parse(strings.NewReader(doc), int64(len(doc)), nil)
	if err != nil {
		t.Fatal(err)
	}
	return s, w
}

// score makes a partwise score of one part from its measures.
func score(measures ...string) string {
	return scoreWith(`<score-part id="P1"><part-name>Music</part-name></score-part>`, measures...)
}

func scoreWith(partList string, measures ...string) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<score-partwise version="4.0"><part-list>` + partList + `</part-list><part id="P1">`)
	for i, m := range measures {
		fmt.Fprintf(&b, `<measure number="%d">%s</measure>`, i+1, m)
	}
	b.WriteString(`</part></score-partwise>`)
	return b.String()
}

// attrs are the attributes of a first measure: divisions and 4/4.
func attrs(div int) string {
	return fmt.Sprintf(`<attributes><divisions>%d</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>`, div)
}

// n is a note: "C4", "F#5", "Bb3", or "r" for a rest.
func n(pitch string, dur int, typ string, extra ...string) string {
	var b strings.Builder
	b.WriteString("<note>")
	rest := strings.Join(extra, "")
	if strings.Contains(rest, "<chord/>") {
		b.WriteString("<chord/>")
		rest = strings.Replace(rest, "<chord/>", "", 1)
	}
	if strings.Contains(rest, "<grace") {
		i := strings.Index(rest, "<grace")
		j := strings.Index(rest[i:], "/>") + i + 2
		b.WriteString(rest[i:j])
		rest = rest[:i] + rest[j:]
	}
	if pitch == "r" {
		b.WriteString("<rest/>")
	} else {
		alter := ""
		switch {
		case strings.Contains(pitch, "#"):
			alter = "<alter>1</alter>"
		case len(pitch) == 3 && pitch[1] == 'b':
			alter = "<alter>-1</alter>"
		}
		fmt.Fprintf(&b, "<pitch><step>%c</step>%s<octave>%c</octave></pitch>", pitch[0], alter, pitch[len(pitch)-1])
	}
	if dur >= 0 {
		fmt.Fprintf(&b, "<duration>%d</duration>", dur)
	}
	b.WriteString(rest)
	if typ != "" {
		fmt.Fprintf(&b, "<type>%s</type>", typ)
	}
	b.WriteString("</note>")
	return b.String()
}

func events(s *music.Score, part, measure int) []*music.Event {
	return s.Parts[part].Measures[measure].Events
}

func offsets(evs []*music.Event) []int {
	var out []int
	for _, e := range evs {
		out = append(out, e.Offset)
	}
	return out
}

func TestDivisions(t *testing.T) {
	// thirds of a quarter, then a change of divisions in the second measure
	var m1, m2 strings.Builder
	m1.WriteString(`<attributes><divisions>3</divisions><time><beats>2</beats><beat-type>4</beat-type></time></attributes>`)
	for i := 0; i < 6; i++ {
		m1.WriteString(n("C5", 1, "eighth", `<time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>`))
	}
	m2.WriteString(`<attributes><divisions>480</divisions></attributes>`)
	m2.WriteString(n("D5", 480, "quarter") + n("E5", 240, "eighth") + n("F5", 240, "eighth"))
	s, _ := parse(t, score(m1.String(), m2.String()))
	if got, want := offsets(events(s, 0, 0)), []int{0, 320, 640, 960, 1280, 1600}; !slices.Equal(got, want) {
		t.Errorf("offsets %v, want %v", got, want)
	}
	if got, want := offsets(events(s, 0, 1)), []int{0, 960, 1440}; !slices.Equal(got, want) {
		t.Errorf("offsets after the change %v, want %v", got, want)
	}
	if s.Measures[0].Length != 1920 || s.Measures[1].Length != 1920 {
		t.Errorf("lengths %d %d", s.Measures[0].Length, s.Measures[1].Length)
	}

	// sevenths do not divide the ticks: positions are rounded from the
	// start of the measure, so they do not drift
	var m strings.Builder
	m.WriteString(`<attributes><divisions>7</divisions></attributes>`)
	for i := 0; i < 28; i++ {
		m.WriteString(n("C5", 1, "16th"))
	}
	s, _ = parse(t, score(m.String()))
	evs := events(s, 0, 0)
	if last := evs[27]; last.Offset != 3703 || last.Offset+last.Duration != 3840 {
		t.Errorf("last seventh at %d for %d", last.Offset, last.Duration)
	}
	sum := 0
	for _, e := range evs {
		sum += e.Duration
	}
	if sum != 3840 {
		t.Errorf("durations add up to %d", sum)
	}
}

func TestChordsVoicesBackup(t *testing.T) {
	m := attrs(2) +
		n("C4", 4, "half", "<voice>1</voice>") + n("E4", 4, "half", "<chord/><voice>1</voice>") + n("G4", 4, "half", "<chord/><voice>1</voice>") +
		n("D4", 4, "half", "<voice>1</voice>") +
		`<backup><duration>8</duration></backup>` +
		n("C3", 2, "quarter", "<voice>a</voice>") + `<forward><duration>4</duration></forward>` + n("G3", 2, "quarter", "<voice>a</voice>")
	s, _ := parse(t, score(m))
	evs := events(s, 0, 0)
	if len(evs) != 4 {
		t.Fatalf("%d events", len(evs))
	}
	if len(evs[0].Notes) != 3 || evs[0].Offset != 0 || evs[0].Duration != 1920 {
		t.Errorf("chord %+v with %d notes", evs[0], len(evs[0].Notes))
	}
	if evs[1].Offset != 1920 {
		t.Errorf("after the chord at %d", evs[1].Offset)
	}
	if evs[2].Voice != 2 || evs[2].Offset != 0 || evs[3].Offset != 2880 || evs[3].Voice != 2 {
		t.Errorf("second voice %+v %+v", evs[2], evs[3])
	}
	if s.Measures[0].Length != 3840 {
		t.Errorf("length %d", s.Measures[0].Length)
	}
}

func TestStavesAndClefs(t *testing.T) {
	m1 := `<attributes><divisions>1</divisions><staves>2</staves><clef number="1"><sign>G</sign><line>2</line></clef></attributes>` +
		n("C5", 4, "whole", "<staff>1</staff>") + `<backup><duration>4</duration></backup>` +
		n("C3", 2, "half", "<staff>2</staff>") +
		`<attributes><clef number="2"><sign>G</sign><line>2</line><clef-octave-change>-1</clef-octave-change></clef></attributes>` +
		n("C4", 2, "half", "<staff>2</staff>")
	m2 := `<attributes><clef number="1"><sign>C</sign><line>4</line></clef><clef number="2"><sign>G</sign><line>2</line><clef-octave-change>-1</clef-octave-change></clef></attributes>` +
		n("C4", 4, "whole", "<staff>1</staff>")
	s, _ := parse(t, score(m1, m2))
	p := s.Parts[0]
	if p.Staves != 2 {
		t.Fatalf("%d staves", p.Staves)
	}
	want := []music.ClefChange{
		{Offset: 0, Staff: 1, Clef: music.Clef{Sign: music.ClefF}}, // the default of a lower staff
		{Offset: 0, Staff: 0, Clef: music.Clef{Sign: music.ClefG}},
		{Offset: 1920, Staff: 1, Clef: music.Clef{Sign: music.ClefG, Octave: -1}},
	}
	if !reflect.DeepEqual(p.Measures[0].Clefs, want) {
		t.Errorf("clefs %+v", p.Measures[0].Clefs)
	}
	// the unchanged clef of the second staff is no change
	if want := []music.ClefChange{{Staff: 0, Clef: music.Clef{Sign: music.ClefC, Line: 4}}}; !reflect.DeepEqual(p.Measures[1].Clefs, want) {
		t.Errorf("clefs of measure 2 %+v", p.Measures[1].Clefs)
	}
	if evs := events(s, 0, 0); evs[1].Staff != 1 || evs[0].Staff != 0 {
		t.Errorf("staves %d %d", evs[0].Staff, evs[1].Staff)
	}
}

func TestTuplets(t *testing.T) {
	tm := `<time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>`
	// with tuplet elements
	m1 := attrs(3) +
		n("C5", 1, "eighth", tm, `<notations><tuplet type="start" bracket="no" show-number="none"/></notations>`) +
		n("D5", 1, "eighth", tm) +
		n("E5", 1, "eighth", tm, `<notations><tuplet type="stop"/></notations>`) +
		n("F5", 1, "eighth", tm, `<notations><tuplet type="start" bracket="yes"/></notations>`) +
		n("G5", 1, "eighth", tm) +
		n("A5", 1, "eighth", tm, `<notations><tuplet type="stop"/></notations>`) +
		n("B5", 6, "half")
	// without: groups that fill 3 eighths
	m2 := n("C5", 1, "eighth", tm) + n("D5", 1, "eighth", tm) + n("E5", 1, "eighth", tm) +
		n("F5", 2, "quarter", tm) + n("G5", 1, "eighth", tm) +
		n("A5", 3, "quarter") + n("B5", 3, "quarter")
	s, _ := parse(t, score(m1, m2))
	evs := events(s, 0, 0)
	a, b := evs[0].Tuplet, evs[3].Tuplet
	if a == nil || b == nil || a == b || evs[1].Tuplet != a || evs[2].Tuplet != a || evs[4].Tuplet != b || evs[5].Tuplet != b || evs[6].Tuplet != nil {
		t.Fatalf("groups %p %p", a, b)
	}
	if a.Actual != 3 || a.Normal != 2 || a.ShowNumber || a.Bracket == nil || *a.Bracket || b.Bracket == nil || !*b.Bracket || !b.ShowNumber {
		t.Errorf("tuplets %+v %+v", *a, *b)
	}
	if evs[1].Duration != 320 || evs[6].Offset != 1920 {
		t.Errorf("tuplet timing %d %d", evs[1].Duration, evs[6].Offset)
	}
	evs = events(s, 0, 1)
	a, b = evs[0].Tuplet, evs[3].Tuplet
	if a == nil || b == nil || a == b || evs[2].Tuplet != a || evs[4].Tuplet != b || evs[5].Tuplet != nil {
		t.Errorf("groups without tuplet elements: %v", []*music.Tuplet{evs[0].Tuplet, evs[1].Tuplet, evs[2].Tuplet, evs[3].Tuplet, evs[4].Tuplet})
	}
	if evs[3].Duration != 640 || evs[3].Type != music.Quarter {
		t.Errorf("quarter of a triplet: %d, %v", evs[3].Duration, evs[3].Type)
	}
}

func TestInferredTypes(t *testing.T) {
	tm := `<time-modification><actual-notes>3</actual-notes><normal-notes>2</normal-notes></time-modification>`
	m := attrs(12) + n("C5", 18, "") + n("D5", 6, "") + n("E5", 4, "", tm) + n("F5", 4, "", tm) + n("G5", 4, "", tm) + n("r", 12, "")
	s, _ := parse(t, score(m))
	evs := events(s, 0, 0)
	for i, want := range []struct {
		typ  music.NoteType
		dots int
	}{{music.Quarter, 1}, {music.Eighth, 0}, {music.Eighth, 0}, {music.Eighth, 0}, {music.Eighth, 0}, {music.Quarter, 0}} {
		if evs[i].Type != want.typ || evs[i].Dots != want.dots {
			t.Errorf("event %d: %v with %d dots, want %v with %d", i, evs[i].Type, evs[i].Dots, want.typ, want.dots)
		}
	}
	// a rest of the whole measure without a type is a measure rest
	s, _ = parse(t, score(`<attributes><divisions>2</divisions><time><beats>3</beats><beat-type>4</beat-type></time></attributes>`+n("r", 6, ""),
		n("r", 6, "", `<rest measure="yes"/>`)))
	if e := events(s, 0, 0)[0]; !e.Rest || !e.MeasureRest {
		t.Errorf("rest %+v", e)
	}
}

func TestGraceNotes(t *testing.T) {
	m := attrs(2) +
		n("D5", -1, "eighth", `<grace slash="yes"/>`) +
		n("C5", 2, "quarter") +
		n("E5", -1, "16th", `<grace/>`) + n("F5", -1, "16th", `<grace/>`) +
		n("G5", 6, "half", "<dot/>")
	s, _ := parse(t, score(m))
	evs := events(s, 0, 0)
	if len(evs) != 5 {
		t.Fatalf("%d events", len(evs))
	}
	g := evs[0]
	if !g.Grace || !g.Slash || g.Duration != 0 || g.Offset != 0 || evs[1].Offset != 0 {
		t.Errorf("grace %+v", g)
	}
	if !evs[2].Grace || evs[2].Slash || evs[2].Offset != 960 || evs[3].Offset != 960 || evs[4].Offset != 960 {
		t.Errorf("graces before G: %v", offsets(evs))
	}
	// played before the principal notes, in time taken from them
	notes := s.Play.Tracks[0].Notes
	want := []music.PlayNote{
		{Tick: 0, Dur: 240, Key: 74}, {Tick: 240, Dur: 720, Key: 72},
		{Tick: 960, Dur: 240, Key: 76}, {Tick: 1200, Dur: 240, Key: 77}, {Tick: 1440, Dur: 2400, Key: 79},
	}
	if len(notes) != len(want) {
		t.Fatalf("played %+v", notes)
	}
	for i := range want {
		if notes[i].Tick != want[i].Tick || notes[i].Dur != want[i].Dur || notes[i].Key != want[i].Key {
			t.Errorf("note %d: %+v, want %+v", i, notes[i], want[i])
		}
	}
}

func TestBeams(t *testing.T) {
	beamed := attrs(2) +
		n("C5", 1, "eighth", `<beam number="1">begin</beam>`) +
		n("D5", 1, "eighth", `<beam number="1">continue</beam>`) +
		n("E5", 1, "eighth", `<beam number="1">end</beam>`) +
		n("F5", 1, "eighth") +
		n("G5", 1, "16th", `<beam number="1">begin</beam><beam number="2">begin</beam>`) +
		n("A5", 1, "16th", `<beam number="1">continue</beam><beam number="2">end</beam>`) +
		n("B5", 1, "16th", `<beam number="1">end</beam><beam number="2">backward hook</beam>`) +
		n("C6", 1, "eighth")
	s, _ := parse(t, score(beamed))
	evs := events(s, 0, 0)
	if !slices.Equal(evs[0].Beams, []music.Beam{music.BeamBegin}) || !slices.Equal(evs[2].Beams, []music.Beam{music.BeamEnd}) {
		t.Errorf("beams %v %v", evs[0].Beams, evs[2].Beams)
	}
	if evs[3].Beams == nil || len(evs[3].Beams) != 0 {
		t.Errorf("an unbeamed eighth in a file with beams: %#v", evs[3].Beams)
	}
	if !slices.Equal(evs[4].Beams, []music.Beam{music.BeamBegin, music.BeamBegin}) || !slices.Equal(evs[6].Beams, []music.Beam{music.BeamEnd, music.BeamBackwardHook}) {
		t.Errorf("sixteenths %v %v", evs[4].Beams, evs[6].Beams)
	}
	s, _ = parse(t, score(attrs(2)+n("C5", 1, "eighth")+n("D5", 1, "eighth")))
	if evs := events(s, 0, 0); evs[0].Beams != nil || evs[1].Beams != nil {
		t.Errorf("a file without beams: %#v", evs[0].Beams)
	}
}

func TestAccidentals(t *testing.T) {
	m := attrs(1) +
		n("F#4", 1, "quarter", `<accidental>sharp</accidental>`) +
		n("F#4", 1, "quarter") +
		n("Bb4", 1, "quarter", `<accidental cautionary="yes">flat</accidental>`) +
		n("C4", 1, "quarter", `<accidental parentheses="yes">natural</accidental>`)
	m2 := attrs(1) +
		n("C4", 1, "quarter", `<accidental>double-sharp</accidental>`) +
		n("C4", 1, "quarter", `<accidental>flat-flat</accidental>`) +
		n("C4", 1, "quarter", `<accidental>quarter-sharp</accidental>`) +
		n("C4", 1, "quarter", `<accidental>three-quarters-flat</accidental>`)
	s, _ := parse(t, score(m, m2))
	var got []music.Accidental
	var cautionary []bool
	for i := 0; i < 2; i++ {
		for _, e := range events(s, 0, i) {
			got = append(got, e.Notes[0].Accidental)
			cautionary = append(cautionary, e.Notes[0].Cautionary)
		}
	}
	want := []music.Accidental{music.AccSharp, music.AccNone, music.AccFlat, music.AccNatural,
		music.AccDoubleSharp, music.AccDoubleFlat, music.AccSharp, music.AccFlat}
	if !slices.Equal(got, want) {
		t.Errorf("accidentals %v, want %v", got, want)
	}
	if !slices.Equal(cautionary[:4], []bool{false, false, true, true}) {
		t.Errorf("cautionary %v", cautionary)
	}
	// altered notes and no accidental elements: the engine chooses
	s, _ = parse(t, score(attrs(1)+n("F#4", 1, "quarter")+n("G4", 1, "quarter")+n("F#4", 2, "half")))
	for _, e := range events(s, 0, 0) {
		if e.Notes[0].Accidental != music.AccAuto {
			t.Errorf("accidental %v without accidental elements", e.Notes[0].Accidental)
		}
	}
	// no alterations: none
	s, _ = parse(t, score(attrs(1)+n("F4", 4, "whole")))
	if a := events(s, 0, 0)[0].Notes[0].Accidental; a != music.AccNone {
		t.Errorf("accidental %v", a)
	}
}

func TestTies(t *testing.T) {
	m1 := attrs(1) + n("C5", 2, "half") + n("G4", 2, "half", `<tie type="start"/><notations><tied type="start"/></notations>`)
	m2 := n("G4", 1, "quarter", `<tie type="stop"/><tie type="start"/><notations><tied type="stop"/><tied type="start"/></notations>`) +
		n("G4", 1, "quarter", `<tie type="stop"/><notations><tied type="stop"/></notations>`) +
		n("A4", 2, "half", `<notations><tied type="start"/></notations>`)
	m3 := n("A4", 4, "whole", `<notations><tied type="stop"/></notations>`)
	s, _ := parse(t, score(m1, m2, m3))
	g1 := events(s, 0, 0)[1].Notes[0]
	g2 := events(s, 0, 1)[0].Notes[0]
	if !g1.TieStart || g1.TieStop || !g2.TieStart || !g2.TieStop {
		t.Errorf("ties %+v %+v", g1, g2)
	}
	notes := s.Play.Tracks[0].Notes
	want := []music.PlayNote{{Tick: 0, Dur: 1920, Key: 72}, {Tick: 1920, Dur: 3840, Key: 67}, {Tick: 5760, Dur: 5760, Key: 69}}
	if len(notes) != len(want) {
		t.Fatalf("played %+v", notes)
	}
	for i := range want {
		if notes[i].Tick != want[i].Tick || notes[i].Dur != want[i].Dur || notes[i].Key != want[i].Key {
			t.Errorf("note %d: %+v, want %+v", i, notes[i], want[i])
		}
	}
}

func TestLyrics(t *testing.T) {
	lyric := func(num, syl, text, extra string) string {
		return fmt.Sprintf(`<lyric number="%s"><syllabic>%s</syllabic><text>%s</text>%s</lyric>`, num, syl, text, extra)
	}
	m := attrs(1) +
		n("C4", 1, "quarter", lyric("1", "begin", "Hal", ""), lyric("chorus", "single", "Oh", "")) +
		n("D4", 1, "quarter", lyric("1", "end", "le", `<extend/>`)) +
		n("E4", 1, "quarter", `<lyric number="1"><extend type="stop"/></lyric>`) +
		n("F4", 1, "quarter", `<lyric number="1"><syllabic>single</syllabic><text>a</text><elision>‿</elision><syllabic>single</syllabic><text>te</text></lyric>`)
	s, _ := parse(t, score(m))
	evs := events(s, 0, 0)
	if l := evs[0].Lyrics; len(l) != 2 || *l[0] != (music.LyricSyllable{Verse: 1, Text: "Hal", Syllabic: music.SyllableBegin}) ||
		*l[1] != (music.LyricSyllable{Verse: 2, Text: "Oh"}) {
		t.Errorf("lyrics of the first note %+v", l)
	}
	if l := evs[1].Lyrics; len(l) != 1 || *l[0] != (music.LyricSyllable{Verse: 1, Text: "le", Syllabic: music.SyllableEnd, Extend: true}) {
		t.Errorf("lyrics of the second note %+v", l)
	}
	if len(evs[2].Lyrics) != 0 {
		t.Errorf("an extend alone %+v", evs[2].Lyrics[0])
	}
	if l := evs[3].Lyrics; len(l) != 1 || l[0].Text != "a‿te" || l[0].Syllabic != music.SyllableSingle {
		t.Errorf("elision %+v", l[0])
	}
	got := s.Play.Tracks[0].Lyrics
	want := []music.Lyric{{Tick: 0, Text: "Hal-"}, {Tick: 960, Text: "le"}, {Tick: 2880, Text: "a‿te"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("played lyrics %+v", got)
	}
}

func TestDirections(t *testing.T) {
	dir := func(placement, types string, extra ...string) string {
		return `<direction placement="` + placement + `"><direction-type>` + types + `</direction-type>` + strings.Join(extra, "") + `</direction>`
	}
	m1 := attrs(2) +
		dir("below", `<dynamics><p/></dynamics>`) +
		dir("above", `<words font-style="italic">dolce</words>`) +
		dir("below", `<wedge type="diminuendo" number="2"/>`, `<staff>1</staff>`) +
		n("C5", 2, "quarter") +
		dir("below", `<wedge type="stop" number="2"/>`) +
		dir("below", `<dynamics><sfz/></dynamics>`) +
		dir("below", `<pedal type="start" line="yes"/>`) +
		n("D5", 2, "quarter") +
		dir("below", `<pedal type="change" line="yes"/>`) +
		dir("above", `<octave-shift type="down" size="8" number="1"/>`) +
		n("C6", 2, "quarter") +
		n("D6", 2, "quarter") +
		dir("above", `<octave-shift type="stop" size="8" number="1"/>`) +
		dir("below", `<pedal type="stop" line="yes"/>`, `<offset>-1</offset>`)
	m2 := dir("above", `<rehearsal>A</rehearsal>`) + dir("above", `<segno/>`) +
		dir("above", `<words>Allegro</words>`) +
		dir("below", `<dynamics><other-dynamics>sfzp</other-dynamics></dynamics>`) +
		dir("above", `<words>cresc.</words>`) +
		n("C6", 8, "whole") + dir("above", `<coda/>`)
	s, _ := parse(t, score(m1, m2))
	type d struct {
		Offset int
		Kind   music.DirectionKind
		Text   string
		Stop   bool
		Number int
		Size   int
		Italic bool
		Above  string
	}
	conv := func(ds []*music.Direction) []d {
		var out []d
		for _, x := range ds {
			a := ""
			if x.Above != nil {
				a = fmt.Sprint(*x.Above)
			}
			out = append(out, d{x.Offset, x.Kind, x.Text, x.Stop, x.Number, x.Size, x.Italic, a})
		}
		return out
	}
	want := []d{
		{0, music.DirDynamic, "p", false, 0, 0, false, "false"},
		{0, music.DirWords, "dolce", false, 0, 0, true, "true"},
		{0, music.DirDiminuendo, "", false, 2, 0, false, "false"},
		{960, music.DirDiminuendo, "", true, 2, 0, false, "false"},
		{960, music.DirDynamic, "sfz", false, 0, 0, false, "false"},
		{960, music.DirPedal, "", false, 1, 0, false, "false"},
		{1920, music.DirPedal, "", true, 1, 0, false, "false"},
		{1920, music.DirPedal, "", false, 1, 0, false, "false"},
		{1920, music.DirOctave, "", false, 1, 8, false, "true"},
		{3840, music.DirOctave, "", true, 1, 0, false, "true"},
		{3360, music.DirPedal, "", true, 1, 0, false, "false"},
	}
	if got := conv(s.Parts[0].Measures[0].Directions); !reflect.DeepEqual(got, want) {
		t.Errorf("directions\n%+v\nwant\n%+v", got, want)
	}
	want = []d{
		{0, music.DirSegno, "", false, 0, 0, false, "true"},
		{0, music.DirDynamic, "sfzp", false, 0, 0, false, "false"},
		{0, music.DirWords, "cresc.", false, 0, 0, false, "true"},
		{3840, music.DirCoda, "", false, 0, 0, false, "true"},
	}
	if got := conv(s.Parts[0].Measures[1].Directions); !reflect.DeepEqual(got, want) {
		t.Errorf("directions of measure 2\n%+v\nwant\n%+v", got, want)
	}
	if m := s.Measures[1]; m.Rehearsal != "A" || len(m.Tempo) != 1 || m.Tempo[0].Text != "Allegro" || m.Tempo[0].PerMinute != 0 {
		t.Errorf("rehearsal %q, tempo %+v", m.Rehearsal, m.Tempo)
	}
	// under the 8va line the notes are written an octave lower and played
	// as they are
	evs := events(s, 0, 0)
	if evs[2].Notes[0].Pitch.Octave != 5 || evs[3].Notes[0].Pitch.Octave != 5 || evs[1].Notes[0].Pitch.Octave != 5 {
		t.Errorf("octaves under the line %d %d", evs[2].Notes[0].Pitch.Octave, evs[3].Notes[0].Pitch.Octave)
	}
	notes := s.Play.Tracks[0].Notes
	if notes[2].Key != 84 {
		t.Errorf("played %d under the line", notes[2].Key)
	}
	// dynamics: p, then an accent on sfz, then sfzp accents and plays p
	if notes[0].Vel != 49 || notes[1].Vel != 69 || notes[2].Vel != 49 || notes[4].Vel != 69 {
		t.Errorf("velocities %d %d %d %d", notes[0].Vel, notes[1].Vel, notes[2].Vel, notes[4].Vel)
	}
}

func TestOctaveShiftBelow(t *testing.T) {
	m := attrs(1) + `<attributes><clef><sign>F</sign><line>4</line></clef></attributes>` +
		`<direction placement="below"><direction-type><octave-shift type="up" size="15"/></direction-type></direction>` +
		n("C1", 4, "whole") +
		`<direction><direction-type><octave-shift type="stop" size="15"/></direction-type></direction>`
	s, _ := parse(t, score(m, n("C1", 4, "whole")))
	if o := events(s, 0, 0)[0].Notes[0].Pitch.Octave; o != 3 {
		t.Errorf("written octave %d under 15mb", o)
	}
	if o := events(s, 0, 1)[0].Notes[0].Pitch.Octave; o != 1 {
		t.Errorf("written octave %d after the line", o)
	}
	if d := s.Parts[0].Measures[0].Directions[0]; d.Size != -15 {
		t.Errorf("size %d", d.Size)
	}
	if k := s.Play.Tracks[0].Notes[0].Key; k != 24 {
		t.Errorf("played %d", k)
	}
}

func TestHarmony(t *testing.T) {
	for _, c := range []struct{ xml, want string }{
		{`<root><root-step>C</root-step></root><kind>major</kind>`, "C"},
		{`<root><root-step>B</root-step><root-alter>-1</root-alter></root><kind>minor-seventh</kind>`, "Bbm7"},
		{`<root><root-step>F</root-step><root-alter>1</root-alter></root><kind>half-diminished</kind>`, "F#m7b5"},
		{`<root><root-step>G</root-step></root><kind text="7">dominant</kind><degree><degree-value>9</degree-value><degree-alter>-1</degree-alter><degree-type>alter</degree-type></degree>`, "G7b9"},
		{`<root><root-step>C</root-step></root><kind>major</kind><degree><degree-value>9</degree-value><degree-alter>0</degree-alter><degree-type>add</degree-type></degree><bass><bass-step>E</bass-step></bass>`, "Cadd9/E"},
		{`<root><root-step>D</root-step></root><kind>suspended-fourth</kind><bass><bass-step>B</bass-step><bass-alter>-1</bass-alter></bass>`, "Dsus4/Bb"},
		{`<root><root-step>A</root-step></root><kind>dominant</kind><degree><degree-value>3</degree-value><degree-alter>0</degree-alter><degree-type>subtract</degree-type></degree>`, "A7no3"},
		{`<kind>none</kind>`, "N.C."},
		{`<root><root-step>E</root-step></root><kind text="min">minor</kind>`, "Emin"},
	} {
		s, _ := parse(t, score(attrs(1)+`<harmony placement="above">`+c.xml+`<offset>2</offset></harmony>`+n("C4", 4, "whole")))
		ds := s.Parts[0].Measures[0].Directions
		if len(ds) != 1 || ds[0].Kind != music.DirChord || ds[0].Text != c.want || ds[0].Offset != 1920 {
			t.Errorf("%s: %+v", c.want, ds)
		}
	}
}

func TestMeasures(t *testing.T) {
	bar := func(loc, inner string) string { return `<barline location="` + loc + `">` + inner + `</barline>` }
	whole := n("C5", 4, "whole")
	doc := score(
		`<attributes><divisions>1</divisions><time><beats>4</beats><beat-type>4</beat-type></time></attributes>`+n("G4", 1, "quarter"),
		bar("left", `<bar-style>heavy-light</bar-style><repeat direction="forward"/>`)+whole,
		whole+bar("right", `<bar-style>light-light</bar-style>`),
		`<print new-system="yes"/>`+bar("left", `<ending number="1, 2" type="start"/>`)+whole+bar("right", `<ending number="1, 2" type="stop"/><repeat direction="backward" times="3"/>`),
		bar("left", `<ending number="3" type="start">Last time</ending>`)+whole+bar("right", `<ending number="3" type="discontinue"/>`),
		`<print new-page="yes"/><attributes><time symbol="cut"><beats>2</beats><beat-type>2</beat-type></time></attributes>`+whole+bar("right", `<bar-style>dashed</bar-style>`),
		`<attributes><divisions>2</divisions><time><beats>3+2</beats><beat-type>8</beat-type></time></attributes>`+n("C5", 2, "quarter")+n("C5", 3, "quarter", "<dot/>")+bar("right", `<bar-style>light-heavy</bar-style>`),
	)
	doc = strings.Replace(doc, `<measure number="1">`, `<measure number="0" implicit="yes">`, 1)
	doc = strings.Replace(doc, `<measure number="3">`, `<measure number="3" text="3a">`, 1)
	s, _ := parse(t, doc)
	ms := s.Measures
	if len(ms) != 7 {
		t.Fatalf("%d measures", len(ms))
	}
	if ms[0].Number != "0" || !ms[0].Implicit || ms[0].Length != 960 || ms[0].Time == nil || *ms[0].Time != (music.TimeSig{Beats: 4, BeatType: 4}) {
		t.Errorf("pickup %+v", *ms[0])
	}
	if ms[1].Left != music.BarRepeatStart || ms[1].Length != 3840 || ms[1].Time != nil || ms[1].Number != "2" {
		t.Errorf("measure 2 %+v", *ms[1])
	}
	if ms[2].Right != music.BarDouble || ms[2].Number != "3a" {
		t.Errorf("measure 3 %+v", *ms[2])
	}
	if m := ms[3]; m.Right != music.BarRepeatEnd || m.Times != 3 || m.Break != music.BreakSystem || m.Ending == nil ||
		*m.Ending != (music.Ending{Text: "1, 2.", Start: true, Stop: true}) {
		t.Errorf("measure 4 %+v %+v", *m, m.Ending)
	}
	if m := ms[4]; m.Ending == nil || *m.Ending != (music.Ending{Text: "Last time", Start: true, Stop: true, Open: true}) {
		t.Errorf("measure 5 ending %+v", m.Ending)
	}
	if m := ms[5]; m.Break != music.BreakPage || m.Time == nil || *m.Time != (music.TimeSig{Beats: 2, BeatType: 2, Symbol: "cut"}) || m.Right != music.BarDashed {
		t.Errorf("measure 6 %+v", *m)
	}
	if m := ms[6]; m.Time == nil || *m.Time != (music.TimeSig{Beats: 5, BeatType: 8}) || m.Length != 2400 || m.Right != music.BarFinal {
		t.Errorf("measure 7 %+v %v", *m, m.Time)
	}
	// played: pickup, 2 3 4 three times, then the third ending
	var order []int
	for _, p := range s.PlayOrder {
		order = append(order, p.Measure)
	}
	if want := []int{0, 1, 2, 3, 1, 2, 3, 1, 2, 4, 5, 6}; !slices.Equal(order, want) {
		t.Errorf("play order %v, want %v", order, want)
	}
	if s.PlayOrder[1].Tick != 960 || s.PlayOrder[4].Tick != 960+3*3840 {
		t.Errorf("ticks %+v", s.PlayOrder[:5])
	}
}

// playOrder parses a score and returns its measures as played.
func playOrder(t *testing.T, measures ...string) []int {
	t.Helper()
	for i := range measures {
		measures[i] += n("C5", 4, "whole")
	}
	measures[0] = attrs(1) + measures[0]
	s, _ := parse(t, score(measures...))
	if s.PlayOrder == nil {
		out := make([]int, len(s.Measures))
		for i := range out {
			out[i] = i
		}
		return out
	}
	var out []int
	tick := 0
	for _, p := range s.PlayOrder {
		if p.Tick != tick {
			t.Errorf("measure %d played at %d, want %d", p.Measure, p.Tick, tick)
		}
		tick += s.Measures[p.Measure].Length
		out = append(out, p.Measure)
	}
	if s.Play.End != tick {
		t.Errorf("performance ends at %d, want %d", s.Play.End, tick)
	}
	return out
}

func TestPlayOrder(t *testing.T) {
	fwd := `<barline location="left"><repeat direction="forward"/></barline>`
	bwd := `<barline location="right"><repeat direction="backward"/></barline>`
	end := func(num, typ string) string {
		loc := "left"
		if typ != "start" {
			loc = "right"
		}
		return `<barline location="` + loc + `"><ending number="` + num + `" type="` + typ + `"/></barline>`
	}
	snd := func(attrs string) string { return `<sound ` + attrs + `/>` }
	words := func(w string) string {
		return `<direction><direction-type><words>` + w + `</words></direction-type></direction>`
	}
	for _, c := range []struct {
		name     string
		measures []string
		want     []int
	}{
		{"no repeats", []string{"", "", ""}, []int{0, 1, 2}},
		{"simple repeat", []string{"", fwd, bwd, ""}, []int{0, 1, 2, 1, 2, 3}},
		{"repeat from the start", []string{"", bwd, "", bwd}, []int{0, 1, 0, 1, 2, 3, 2, 3}},
		{"voltas", []string{fwd, "", end("1", "start") + end("1", "stop") + bwd, end("2", "start") + end("2", "discontinue"), ""},
			[]int{0, 1, 2, 0, 1, 3, 4}},
		{"two measure first volta", []string{fwd, end("1", "start"), end("1", "stop") + bwd, end("2", "start") + end("2", "stop"), fwd, bwd},
			[]int{0, 1, 2, 0, 3, 4, 5, 4, 5}},
		{"first volta only", []string{fwd, end("1", "start") + end("1", "stop") + bwd, "", bwd},
			[]int{0, 1, 0, 2, 3, 2, 3}},
		{"D.C. al Fine", []string{"", snd(`fine="yes"`), snd(`dacapo="yes"`)}, []int{0, 1, 2, 0, 1}},
		{"D.C. al Fine with a repeat", []string{fwd, bwd + snd(`fine="yes"`), "", snd(`dacapo="yes"`)}, []int{0, 1, 0, 1, 2, 3, 0, 1}},
		{"D.S. al Coda", []string{"", snd(`segno="s1"`), snd(`tocoda="c1"`), snd(`dalsegno="s1"`), snd(`coda="c1"`)},
			[]int{0, 1, 2, 3, 1, 2, 4}},
		{"D.C. al Fine in words", []string{"", words("Fine"), words("D.C. al Fine")}, []int{0, 1, 2, 0, 1}},
		{"D.S. al Coda in signs", []string{
			"", `<direction><direction-type><segno/></direction-type></direction>`,
			`<direction><direction-type><coda/></direction-type></direction>`, words("D.S. al Coda"),
			`<direction><direction-type><coda/></direction-type></direction>`,
		}, []int{0, 1, 2, 3, 1, 2, 4}},
		{"D.C. takes the last volta", []string{fwd, end("1", "start") + end("1", "stop") + bwd, end("2", "start") + end("2", "stop"), snd(`dacapo="yes"`)},
			[]int{0, 1, 0, 2, 3, 0, 2, 3}},
		{"voltas without numbers", []string{"", end("", "start") + end("", "stop") + bwd, end("", "start") + end("", "discontinue"), ""},
			[]int{0, 1, 0, 2, 3}},
		{"three passes, two voltas", []string{fwd, end("1, 2", "start") + end("1, 2", "stop") + bwd, end("3", "start") + end("3", "stop"), ""},
			[]int{0, 1, 0, 1, 0, 2, 3}},
		{"a repeat inside the first volta", []string{"", end("1", "start"), fwd + bwd, end("1", "stop") + bwd, end("2", "start") + end("2", "stop"), ""},
			[]int{0, 1, 2, 2, 3, 0, 4, 5}},
		{"the second volta starts a repeat", []string{fwd, end("1", "start") + end("1", "stop") + bwd, end("2", "start") + fwd + end("2", "stop"), bwd, ""},
			[]int{0, 1, 0, 2, 3, 2, 3, 4}},
		{"a backward repeat goes back to the last one", []string{"", bwd, "", bwd}, []int{0, 1, 0, 1, 2, 3, 2, 3}},
		{"D.C. al Coda passes Fine, then D.S. al Fine", []string{
			"", snd(`tocoda="c"`), snd(`segno="s"`), snd(`fine="yes"`) + words("Fine"),
			words("D.C. al Coda") + snd(`dacapo="yes"`), snd(`coda="c"`), words("D.S. al Fine") + snd(`dalsegno="s"`),
		}, []int{0, 1, 2, 3, 4, 0, 1, 5, 6, 2, 3}},
	} {
		if got := playOrder(t, c.measures...); !slices.Equal(got, c.want) {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}

func TestPlayOrderLimit(t *testing.T) {
	bwd := `<barline location="right"><repeat direction="backward" times="100"/></barline>`
	var ms []string
	for i := 0; i < 200; i++ {
		ms = append(ms, n("C5", 4, "whole"))
	}
	ms[0] = attrs(1) + ms[0]
	ms[199] += bwd
	s, w := parse(t, score(ms...))
	if len(s.PlayOrder) != maxPlays {
		t.Errorf("%d measures played", len(s.PlayOrder))
	}
	if !slices.ContainsFunc(w, func(s string) bool { return strings.Contains(s, "repeats play more") }) {
		t.Errorf("warnings %v", w)
	}
}

func TestTempo(t *testing.T) {
	fwd := `<barline location="left"><repeat direction="forward"/></barline>`
	bwd := `<barline location="right"><repeat direction="backward"/></barline>`
	metro := func(unit, dot, perMinute string) string {
		return `<direction placement="above"><direction-type><metronome><beat-unit>` + unit + `</beat-unit>` + dot +
			`<per-minute>` + perMinute + `</per-minute></metronome></direction-type></direction>`
	}
	w := n("C5", 4, "whole")
	s, _ := parse(t, score(
		attrs(1)+`<direction placement="above"><direction-type><words>Andante</words></direction-type><direction-type><metronome><beat-unit>quarter</beat-unit><per-minute>72</per-minute></metronome></direction-type><sound tempo="72"/></direction>`+w,
		fwd+w,
		metro("half", "<beat-unit-dot/>", "c. 60")+n("C5", 2, "half")+`<sound tempo="100"/>`+n("C5", 2, "half")+bwd,
		w,
	))
	if m := s.Measures[0]; len(m.Tempo) != 1 || m.Tempo[0] != (music.TempoMark{Text: "Andante", Beat: music.Quarter, PerMinute: 72}) {
		t.Errorf("tempo mark %+v", m.Tempo)
	}
	if m := s.Measures[2]; len(m.Tempo) != 1 || m.Tempo[0] != (music.TempoMark{Beat: music.Half, Dots: 1, PerMinute: 60}) {
		t.Errorf("metronome mark %+v", m.Tempo)
	}
	// measure 3 starts at a dotted half = 60 (180 quarters), changes to 100
	// in the middle; the repeat goes back to 72 at measure 2
	want := []music.Tempo{{Tick: 0, BPM: 72}, {Tick: 7680, BPM: 180}, {Tick: 9600, BPM: 100}, {Tick: 11520, BPM: 72},
		{Tick: 15360, BPM: 180}, {Tick: 17280, BPM: 100}}
	if !reflect.DeepEqual(s.Play.Tempo, want) {
		t.Errorf("tempo %+v", s.Play.Tempo)
	}
	// no tempo: 120
	s, _ = parse(t, score(attrs(1)+w))
	if !reflect.DeepEqual(s.Play.Tempo, []music.Tempo{{BPM: 120}}) {
		t.Errorf("default tempo %+v", s.Play.Tempo)
	}
}

func TestTransposingInstrument(t *testing.T) {
	part := `<score-part id="P1"><part-name>Clarinet in B♭</part-name><part-abbreviation>Cl.</part-abbreviation>` +
		`<midi-instrument id="P1-I1"><midi-channel>3</midi-channel><midi-program>72</midi-program><volume>50</volume><pan>-90</pan></midi-instrument></score-part>`
	m := `<attributes><divisions>1</divisions><key><fifths>2</fifths></key><transpose><diatonic>-1</diatonic><chromatic>-2</chromatic></transpose></attributes>` +
		n("D5", 4, "whole")
	m2 := `<attributes><transpose><diatonic>-1</diatonic><chromatic>-2</chromatic><octave-change>-1</octave-change></transpose></attributes>` + n("D5", 4, "whole")
	s, _ := parse(t, scoreWith(part, m, m2))
	p := s.Parts[0]
	if p.Name != "Clarinet in B♭" || p.Abbrev != "Cl." {
		t.Errorf("names %q %q", p.Name, p.Abbrev)
	}
	if e := events(s, 0, 0)[0]; e.Notes[0].Pitch != (music.Pitch{Step: 1, Octave: 5}) {
		t.Errorf("written %+v", e.Notes[0].Pitch)
	}
	tr := s.Play.Tracks[0]
	if tr.Notes[0].Key != 72 || tr.Notes[1].Key != 60 {
		t.Errorf("sounding %d %d", tr.Notes[0].Key, tr.Notes[1].Key)
	}
	if tr.Channel != 2 || tr.Program != 71 || tr.Volume != 64 || tr.Pan != 0 || tr.Name != "Clarinet in B♭" {
		t.Errorf("track %+v", *tr)
	}
}

func TestParts(t *testing.T) {
	doc := `<?xml version="1.0"?><score-partwise><part-list>
<score-part id="P1"><part-name print-object="no">Soprano</part-name></score-part>
<score-part id="P2"><part-name>Piano</part-name><part-name-display><display-text>Piano</display-text><display-text> in E</display-text><accidental-text>flat</accidental-text></part-name-display></score-part>
<score-part id="P3"><part-name>Drums</part-name><score-instrument id="P3-I36"><instrument-name>Bass Drum</instrument-name></score-instrument>
<score-instrument id="P3-I39"><instrument-name>Snare</instrument-name></score-instrument>
<midi-instrument id="P3-I36"><midi-channel>10</midi-channel><midi-unpitched>36</midi-unpitched></midi-instrument>
<midi-instrument id="P3-I39"><midi-channel>10</midi-channel><midi-unpitched>39</midi-unpitched></midi-instrument></score-part>
<score-part id="P4"><part-name>Triangle</part-name></score-part>
<score-part id="P5"><part-name>Empty</part-name></score-part>
</part-list>
<part id="P1"><measure number="1"><attributes><divisions>1</divisions></attributes>` + n("C5", 4, "whole") + `</measure></part>
<part id="P2"><measure number="1"><attributes><divisions>1</divisions><staves>2</staves></attributes>` + n("C5", 4, "whole") + `</measure></part>
<part id="P3"><measure number="1"><attributes><divisions>1</divisions><clef><sign>percussion</sign></clef></attributes>
<note><unpitched><display-step>F</display-step><display-octave>4</display-octave></unpitched><duration>2</duration><instrument id="P3-I36"/><type>half</type></note>
<note><unpitched><display-step>C</display-step><display-octave>5</display-octave></unpitched><duration>1</duration><instrument id="P3-I39"/><type>quarter</type><notehead>x</notehead></note>
<note><unpitched><display-step>G</display-step><display-octave>5</display-octave></unpitched><duration>1</duration><type>quarter</type></note>
</measure></part>
<part id="P4"><measure number="1"><attributes><divisions>1</divisions></attributes><note><unpitched><display-step>E</display-step><display-octave>4</display-octave></unpitched><duration>4</duration><type>whole</type></note></measure></part>
</score-partwise>`
	s, w := parse(t, doc)
	if len(s.Parts) != 4 {
		t.Fatalf("%d parts", len(s.Parts))
	}
	if !slices.ContainsFunc(w, func(s string) bool { return strings.Contains(s, `"P5" has no music`) }) {
		t.Errorf("warnings %v", w)
	}
	if s.Parts[0].Name != "" || s.Parts[1].Name != "Piano in Eb" || s.Parts[1].Staves != 2 {
		t.Errorf("names %q %q", s.Parts[0].Name, s.Parts[1].Name)
	}
	drums := s.Parts[2]
	if !drums.Percussion || drums.Measures[0].Clefs[0].Clef.Sign != music.ClefPercussion {
		t.Errorf("drums %+v", drums.Measures[0].Clefs)
	}
	if e := drums.Measures[0].Events[1]; e.Notes[0].Head != music.HeadX || e.Notes[0].Pitch != (music.Pitch{Step: 0, Octave: 5}) {
		t.Errorf("snare %+v", e.Notes[0])
	}
	var chans []int
	for _, tr := range s.Play.Tracks {
		chans = append(chans, tr.Channel)
	}
	if !slices.Equal(chans, []int{0, 1, 9, 9}) {
		t.Errorf("channels %v", chans)
	}
	var keys []int
	for _, n := range s.Play.Tracks[2].Notes {
		keys = append(keys, n.Key)
	}
	if !slices.Equal(keys, []int{35, 38, 42}) {
		t.Errorf("drum keys %v", keys)
	}
	if s.Play.Tracks[0].Name != "Soprano" {
		t.Errorf("track name %q", s.Play.Tracks[0].Name)
	}
}

func TestCredits(t *testing.T) {
	head := `<?xml version="1.0"?><score-partwise><work><work-title>Sonata</work-title></work><movement-title>Allegro</movement-title>
<identification><creator type="composer">W. A. Mozart</creator><creator type="poet">Nobody</creator><rights>Public domain</rights></identification>`
	body := `<part-list><score-part id="P1"><part-name>Piano</part-name></score-part></part-list><part id="P1"><measure number="1">` + n("C4", 1, "quarter") + `</measure></part></score-partwise>`
	s, _ := parse(t, head+body)
	if s.Title != "Sonata" || s.Subtitle != "Allegro" || s.Composer != "W. A. Mozart" || s.Lyricist != "Nobody" || s.Rights != "Public domain" {
		t.Errorf("credits %q %q %q %q %q", s.Title, s.Subtitle, s.Composer, s.Lyricist, s.Rights)
	}
	typed := head + `<credit page="1"><credit-type>title</credit-type><credit-words font-size="24">Sonata
 in C</credit-words></credit><credit page="1"><credit-type>composer</credit-type><credit-words>Mozart</credit-words></credit>` +
		`<credit page="1"><credit-type>arranger</credit-type><credit-words>arr. Someone</credit-words></credit>` + body
	s, _ = parse(t, typed)
	if s.Title != "Sonata in C" || s.Subtitle != "" || s.Composer != "Mozart" || s.Arranger != "arr. Someone" {
		t.Errorf("typed credits %q %q %q %q", s.Title, s.Subtitle, s.Composer, s.Arranger)
	}
	untyped := `<?xml version="1.0"?><score-partwise><credit page="1"><credit-words justify="center" font-size="12">A subtitle</credit-words></credit>
<credit page="1"><credit-words justify="center" font-size="24">The Title</credit-words></credit>
<credit page="1"><credit-words justify="right" font-size="10">The Composer</credit-words></credit>` + body
	s, _ = parse(t, untyped)
	if s.Title != "The Title" || s.Composer != "The Composer" {
		t.Errorf("untyped credits %q %q", s.Title, s.Composer)
	}
	s, _ = parse(t, `<score-partwise><movement-title>Only &quot;this&quot; &amp; that&nbsp;&unknown;</movement-title>`+body)
	if s.Title != "Only \"this\" & that &unknown;" {
		t.Errorf("title %q", s.Title)
	}
}

// ode reads a test score.
func readFile(t *testing.T, name string) (*music.Score, []string) {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	s, w, err := Parse(bytes.NewReader(b), int64(len(b)), nil)
	if err != nil {
		t.Fatal(err)
	}
	return s, w
}

func TestTimewise(t *testing.T) {
	s, w := readFile(t, "twinkle_timewise.musicxml")
	if len(w) > 0 {
		t.Errorf("warnings %v", w)
	}
	// the same score, partwise
	tw, _ := os.ReadFile("testdata/twinkle_timewise.musicxml")
	pw := partwiseOf(t, string(tw))
	s2, _ := parse(t, pw)
	if !reflect.DeepEqual(s, s2) {
		t.Errorf("timewise and partwise differ")
	}
	if len(s.Parts) != 2 || len(s.Measures) != 4 || s.Title != "Ah! vous dirai-je, maman" || s.Measures[3].Right != music.BarFinal {
		t.Errorf("%d parts, %d measures, %q", len(s.Parts), len(s.Measures), s.Title)
	}
	if got := len(s.Play.Tracks[0].Notes) + len(s.Play.Tracks[1].Notes); got != 14+8 {
		t.Errorf("%d notes played", got)
	}
}

// partwiseOf rewrites the timewise test score partwise.
func partwiseOf(t *testing.T, tw string) string {
	t.Helper()
	head, rest, _ := strings.Cut(tw, "<measure ")
	head = strings.Replace(head, "score-timewise", "score-partwise", -1)
	head = strings.Replace(head, "Timewise", "Partwise", -1)
	parts := map[string][]string{}
	for _, m := range strings.Split("<measure "+rest, "<measure ")[1:] {
		num := m[:strings.Index(m, ">")]
		for _, p := range strings.Split(m, `<part id="`)[1:] {
			id := p[:strings.Index(p, `"`)]
			body := p[strings.Index(p, ">")+1 : strings.Index(p, "</part>")]
			parts[id] = append(parts[id], "<measure "+num+">"+body+"</measure>")
		}
	}
	var b strings.Builder
	b.WriteString(head)
	for _, id := range []string{"P1", "P2"} {
		b.WriteString(`<part id="` + id + `">` + strings.Join(parts[id], "\n") + "</part>\n")
	}
	b.WriteString("</score-partwise>\n")
	return b.String()
}

func TestTestScores(t *testing.T) {
	ode, w := readFile(t, "ode_to_joy.musicxml")
	if len(w) > 0 {
		t.Errorf("warnings %v", w)
	}
	if ode.Title != "Ode to Joy" || ode.Subtitle != "Symphony No. 9, finale" || ode.Composer != "Ludwig van Beethoven" || ode.Lyricist != "Friedrich Schiller" {
		t.Errorf("credits %q %q %q %q", ode.Title, ode.Subtitle, ode.Composer, ode.Lyricist)
	}
	if len(ode.Measures) != 16 || len(ode.PlayOrder) != 24 {
		t.Errorf("%d measures, %d played", len(ode.Measures), len(ode.PlayOrder))
	}
	if k := ode.Parts[0].Measures[0].Key; k == nil || k.Fifths != 2 {
		t.Errorf("key %v", k)
	}
	// F sharps without accidental elements: the engine chooses
	if a := events(ode, 0, 0)[0].Notes[0].Accidental; a != music.AccAuto {
		t.Errorf("accidental %v", a)
	}
	tr := ode.Play.Tracks[0]
	if tr.Program != 52 || tr.Channel != 0 || tr.Notes[0].Key != 66 || tr.Notes[0].Vel != 80 {
		t.Errorf("track %+v, first note %+v", *tr, tr.Notes[0])
	}
	if len(tr.Lyrics) != 60+30 || tr.Lyrics[0].Text != "Freu-" || tr.Lyrics[1].Text != "de," {
		t.Errorf("%d lyrics: %v", len(tr.Lyrics), tr.Lyrics[:2])
	}
	if m := ode.Measures[0]; len(m.Tempo) != 1 || m.Tempo[0].Text != "Allegro assai" || m.Tempo[0].PerMinute != 120 {
		t.Errorf("tempo %+v", m.Tempo)
	}

	mxl, w := readFile(t, "ode_to_joy.mxl")
	if len(w) > 0 || !reflect.DeepEqual(ode, mxl) {
		t.Errorf("the compressed score differs (warnings %v)", w)
	}

	minuet, w := readFile(t, "minuet.musicxml")
	if len(w) > 0 {
		t.Errorf("warnings %v", w)
	}
	p := minuet.Parts[0]
	if p.Staves != 2 || len(p.Measures[0].Clefs) != 2 || p.Measures[0].Clefs[1].Clef.Sign != music.ClefF {
		t.Errorf("piano staves %d, clefs %+v", p.Staves, p.Measures[0].Clefs)
	}
	var order []int
	for _, pm := range minuet.PlayOrder {
		order = append(order, pm.Measure)
	}
	if want := []int{0, 1, 2, 3, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 6, 8}; !slices.Equal(order, want) {
		t.Errorf("minuet played %v", order)
	}
	if e := minuet.Measures[8].Ending; e == nil || e.Text != "2." || !e.Open {
		t.Errorf("second ending %+v", e)
	}
	if v := minuet.Play.Tracks[0].Notes[0].Vel; v != 49 {
		t.Errorf("velocity %d under p", v)
	}
}

// utf16Doc encodes a document as UTF-16 with a byte order mark.
func utf16Doc(s string, bigEndian bool) []byte {
	u := utf16.Encode([]rune("\ufeff" + s))
	b := make([]byte, 2*len(u))
	for i, c := range u {
		if bigEndian {
			binary.BigEndian.PutUint16(b[2*i:], c)
		} else {
			binary.LittleEndian.PutUint16(b[2*i:], c)
		}
	}
	return b
}

func TestEncodings(t *testing.T) {
	doc := score(attrs(1) + n("C5", 4, "whole", `<lyric><text>Grüß</text></lyric>`))
	want, _ := parse(t, doc)
	check := func(name string, b []byte) {
		t.Helper()
		if !Detect(b[:min(len(b), 1024)], bytes.NewReader(b), int64(len(b))) {
			t.Errorf("%s: not detected", name)
		}
		s, _, err := Parse(bytes.NewReader(b), int64(len(b)), nil)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got := events(s, 0, 0)[0].Lyrics[0].Text; got != "Grüß" {
			t.Errorf("%s: lyric %q", name, got)
		}
		if !reflect.DeepEqual(s, want) {
			t.Errorf("%s: the score differs", name)
		}
	}
	utf16Text := strings.Replace(doc, `encoding="UTF-8"`, `encoding="UTF-16"`, 1)
	check("UTF-16LE", utf16Doc(utf16Text, false))
	check("UTF-16BE", utf16Doc(utf16Text, true))
	check("UTF-8 BOM", append([]byte("\xef\xbb\xbf"), doc...))
	latin := strings.Replace(doc, `encoding="UTF-8"`, `encoding="ISO-8859-1"`, 1)
	latin = strings.Replace(latin, "Grüß", "Gr\xfc\xdf", 1)
	check("ISO-8859-1", []byte(latin))
	check("Windows-1252 said to be UTF-8", []byte(strings.Replace(doc, "Grüß", "Gr\xfc\xdf", 1)))
	withDTD := strings.Replace(doc, `<score-partwise`, `<!DOCTYPE score-partwise PUBLIC "-//Recordare//DTD MusicXML 4.0 Partwise//EN" "http://www.musicxml.org/dtds/partwise.dtd" [<!ENTITY x "y">]>`+"\n"+`<score-partwise`, 1)
	check("DOCTYPE", []byte(withDTD))
}

// mxl makes a compressed MusicXML file of entries (name, content).
func mxl(t *testing.T, entries ...string) []byte {
	t.Helper()
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	for i := 0; i+1 < len(entries); i += 2 {
		method := zip.Deflate
		if entries[i] == "mimetype" {
			method = zip.Store
		}
		w, err := z.CreateHeader(&zip.FileHeader{Name: entries[i], Method: method})
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(entries[i+1]))
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

const container = `<?xml version="1.0" encoding="UTF-8"?><container><rootfiles><rootfile full-path="%s" media-type="%s"/></rootfiles></container>`

func TestMXL(t *testing.T) {
	doc := score(attrs(1) + n("C5", 4, "whole"))
	for _, c := range []struct {
		name    string
		entries []string
	}{
		{"mimetype and container", []string{"mimetype", "application/vnd.recordare.musicxml", "META-INF/container.xml",
			fmt.Sprintf(container, "score.musicxml", "application/vnd.recordare.musicxml+xml"), "score.musicxml", doc}},
		{"container only", []string{"META-INF/container.xml", fmt.Sprintf(container, "music/score.xml", "application/vnd.recordare.musicxml+xml"),
			"music/other.xml", "<x/>", "music/score.xml", doc}},
		{"container without media type", []string{"META-INF/container.xml", fmt.Sprintf(container, "score.xml", ""), "score.xml", doc}},
		{"a PDF rootfile first", []string{"META-INF/container.xml", `<container><rootfiles><rootfile full-path="score.pdf" media-type="application/pdf"/>` +
			`<rootfile full-path="score.musicxml" media-type="application/vnd.recordare.musicxml+xml"/></rootfiles></container>`, "score.pdf", "%PDF-1.4", "score.musicxml", doc}},
	} {
		b := mxl(t, c.entries...)
		if !Detect(b[:min(len(b), 1024)], bytes.NewReader(b), int64(len(b))) {
			t.Errorf("%s: not detected", c.name)
		}
		s, _, err := Parse(bytes.NewReader(b), int64(len(b)), nil)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if len(s.Measures) != 1 || len(events(s, 0, 0)) != 1 {
			t.Errorf("%s: %d measures", c.name, len(s.Measures))
		}
	}
	// no container: the first score outside META-INF
	b := mxl(t, "score.musicxml", doc)
	if _, _, err := Parse(bytes.NewReader(b), int64(len(b)), nil); err != nil {
		t.Errorf("without a container: %v", err)
	}
	b = mxl(t, "readme.txt", "hello")
	if _, _, err := Parse(bytes.NewReader(b), int64(len(b)), nil); err == nil {
		t.Errorf("a ZIP without a score parses")
	}
}

func TestDetect(t *testing.T) {
	doc := score(attrs(1) + n("C5", 4, "whole"))
	detect := func(b []byte) bool {
		return Detect(b[:min(len(b), 1024)], bytes.NewReader(b), int64(len(b)))
	}
	longComment := `<?xml version="1.0"?><!DOCTYPE score-partwise PUBLIC "-//Recordare//DTD MusicXML 3.1 Partwise//EN" "http://www.musicxml.org/dtds/partwise.dtd"><!--` +
		strings.Repeat("x", 2000) + `--><score-partwise/>`
	for _, c := range []struct {
		name string
		b    []byte
	}{
		{"partwise", []byte(doc)},
		{"timewise", []byte(`<?xml version="1.0"?><score-timewise version="3.0"><part-list/></score-timewise>`)},
		{"no declaration", []byte(`<score-partwise version="3.0"><part-list/></score-partwise>`)},
		{"DOCTYPE with the root beyond the head", []byte(longComment)},
		{"UTF-16 without a byte order mark", utf16Doc(doc, false)[2:]},
		{"mimetype", mxl(t, "mimetype", "application/vnd.recordare.musicxml", "x.musicxml", doc)},
	} {
		if !detect(c.b) {
			t.Errorf("%s: not detected", c.name)
		}
	}
	for _, c := range []struct {
		name string
		b    []byte
	}{
		{"empty", nil},
		{"HTML", []byte(`<!DOCTYPE html><html><body><pre>&lt;score-partwise&gt;</pre></body></html>`)},
		{"SVG", []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"/>`)},
		{"text that mentions it", []byte(`Convert a <score-partwise> file`)},
		{"opus", []byte(`<?xml version="1.0"?><!DOCTYPE opus PUBLIC "-//Recordare//DTD MusicXML 4.0 Opus//EN" "http://www.musicxml.org/dtds/opus.dtd"><opus/>`)},
		{"container", []byte(`<?xml version="1.0"?><!DOCTYPE container PUBLIC "-//Recordare//DTD MusicXML 4.0 Container//EN" "http://www.musicxml.org/dtds/container.dtd"><container/>`)},
		{"Office document", mxl(t, "[Content_Types].xml", `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"/>`,
			"_rels/.rels", `<Relationships/>`, "word/document.xml", `<w:document/>`)},
		{"EPUB", mxl(t, "mimetype", "application/epub+zip", "META-INF/container.xml", fmt.Sprintf(container, "OEBPS/content.opf", "application/oebps-package+xml"),
			"OEBPS/content.opf", "<package/>")},
		{"ZIP of XML", mxl(t, "a.xml", "<a/>", "b.gbr", "G04*")},
		{"container naming other XML", mxl(t, "META-INF/container.xml", fmt.Sprintf(container, "a.xml", ""), "a.xml", "<a/>")},
	} {
		if detect(c.b) {
			t.Errorf("%s: detected", c.name)
		}
	}
}

func TestRegister(t *testing.T) {
	f := conv.Lookup("musicxml")
	if f == nil || !slices.Equal(f.Extensions, []string{".musicxml", ".mxl"}) || f.Description != "MusicXML score" {
		t.Fatalf("format %+v", f)
	}
	b, _ := os.ReadFile("testdata/minuet.musicxml")
	if got := conv.Detect(bytes.NewReader(b), int64(len(b))); got != f {
		t.Errorf("detected as %v", got)
	}
}

func TestErrors(t *testing.T) {
	for _, doc := range []string{
		``,
		`<?xml version="1.0"?>`,
		`<svg/>`,
		`<score-partwise></score-partwise>`,
		`<opus/>`,
	} {
		if _, _, err := Parse(strings.NewReader(doc), int64(len(doc)), nil); err == nil {
			t.Errorf("%q parses", doc)
		}
	}
	// a document cut short keeps what it has
	doc := score(attrs(1)+n("C5", 4, "whole"), n("D5", 4, "whole"))
	cut := doc[:strings.Index(doc, `<measure number="2">`)+30]
	s, w, err := Parse(strings.NewReader(cut), int64(len(cut)), nil)
	if err != nil || len(s.Measures) != 1 || len(w) == 0 {
		t.Errorf("cut short: %v, %v", err, w)
	}
}

func TestNoPlay(t *testing.T) {
	doc := score(attrs(1) + n("C5", 4, "whole"))
	s, _, err := Parse(strings.NewReader(doc), int64(len(doc)), &Options{NoPlay: true})
	if err != nil || s.Play != nil || s.PlayOrder != nil {
		t.Errorf("NoPlay: %v %v", err, s.Play)
	}
}

func TestWordsRuns(t *testing.T) {
	m := attrs(1) + `<direction placement="below"><direction-type><words font-style="italic">cre</words><words font-weight="bold">sc. </words><words>molto</words></direction-type></direction>` +
		n("C5", 4, "whole")
	s, _ := parse(t, score(m))
	ds := s.Parts[0].Measures[0].Directions
	if len(ds) != 1 || ds[0].Text != "cresc. molto" || !ds[0].Italic || ds[0].Kind != music.DirWords {
		t.Errorf("words %+v", ds)
	}
}

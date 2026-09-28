package music

import (
	"encoding/binary"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// engrave builds a score without the fonts of the system.
func engrave(t *testing.T, s *Score) (*Result, error) {
	t.Helper()
	return Build(s, Options{Source: "test", FontDirs: []string{"../../pptx/testdata/fonts"}, NoSystemFonts: true})
}

func hasWarning(ws []string, part string) bool {
	return slices.ContainsFunc(ws, func(w string) bool { return strings.Contains(w, part) })
}

// manyTracks makes a performance of n tracks of one quarter note each; the
// first has another one at tick last.
func manyTracks(n, last int) *Performance {
	p := &Performance{}
	for i := 0; i < n; i++ {
		tr := &Track{Volume: -1, Pan: -1, Notes: []PlayNote{{Tick: 0, Dur: q, Key: 60, Vel: 90}}}
		if i == 0 {
			tr.Notes = append(tr.Notes, PlayNote{Tick: last, Dur: q, Key: 62, Vel: 90})
		}
		p.Tracks = append(p.Tracks, tr)
	}
	return p
}

func TestNotateBudget(t *testing.T) {
	// 51 tracks of 20000 measures are more than a million measures of
	// staves: they are refused before they are written
	s := Notate(manyTracks(51, 20000*4*q-q), NotateOptions{})
	if len(s.Measures) != 20000 || len(s.Parts) != 0 {
		t.Errorf("%d measures, %d parts", len(s.Measures), len(s.Parts))
	}
	if _, err := engrave(t, s); err == nil || !strings.Contains(err.Error(), "measures on its staves (51 tracks of 20000 measures)") {
		t.Errorf("error %v", err)
	}
	// more tracks than a system has staves
	s = Notate(manyTracks(maxStaves+1, 0), NotateOptions{})
	if _, err := engrave(t, s); len(s.Parts) != 0 || err == nil || !strings.Contains(err.Error(), "staves") {
		t.Errorf("%d parts, error %v", len(s.Parts), err)
	}

	// at the limit and past it
	defer func(n int) { maxStaffMeasures = n }(maxStaffMeasures)
	maxStaffMeasures = 12
	s = Notate(manyTracks(3, 4*4*q-q), NotateOptions{})
	if res, err := engrave(t, s); err != nil || res.Parts != 3 || res.Measures != 4 {
		t.Errorf("3 tracks of 4 measures: %v", err)
	}
	s = Notate(manyTracks(3, 5*4*q-q), NotateOptions{})
	if _, err := engrave(t, s); len(s.Parts) != 0 || err == nil {
		t.Errorf("3 tracks of 5 measures: %d parts, error %v", len(s.Parts), err)
	}
}

func TestBuildBudget(t *testing.T) {
	score := func(parts, staves, measures int) *Score {
		s := &Score{}
		for i := 0; i < measures; i++ {
			s.Measures = append(s.Measures, &Measure{Length: 4 * PPQ})
		}
		for i := 0; i < parts; i++ {
			s.Parts = append(s.Parts, &Part{Staves: staves})
		}
		return s
	}
	for _, c := range []struct {
		name                    string
		parts, staves, measures int
		want                    string
	}{
		{"staves times measures", 1001, 1, 1000, "more than 1000000 measures on its staves (1001 staves of 1000 measures)"},
		{"the staves of grand staves count", 501, 2, 1000, "(1002 staves of 1000 measures)"},
		{"staves", maxStaves + 1, 1, 1, "more than 4096 staves"},
		{"grand staves", maxStaves/2 + 1, 2, 1, "more than 4096 staves"},
	} {
		if _, err := engrave(t, score(c.parts, c.staves, c.measures)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v", c.name, err)
		}
	}
	if _, err := engrave(t, score(3, 2, 4)); err != nil {
		t.Errorf("a small score: %v", err)
	}
}

func TestNotateLateTrack(t *testing.T) {
	// measures of 1/64 (60 ticks): the music is cut after 20000 of them,
	// and the second track starts after that
	end := maxMeasures * 60
	p := &Performance{Tracks: []*Track{
		{Name: "early", Volume: -1, Pan: -1, Notes: []PlayNote{{Tick: 0, Dur: 60, Key: 60, Vel: 90}}},
		{Name: "late", Volume: -1, Pan: -1, Notes: []PlayNote{{Tick: end + 120, Dur: 60, Key: 60, Vel: 90}}},
		{Name: "late drums", Channel: 9, Volume: -1, Pan: -1, Notes: []PlayNote{{Tick: end, Dur: 60, Key: 38, Vel: 90}}},
	}}
	s := Notate(p, NotateOptions{Time: &TimeSig{Beats: 1, BeatType: 64}})
	if len(s.Measures) != maxMeasures || len(s.Parts) != 1 || s.Parts[0].Name != "early" {
		t.Fatalf("%d measures, %d parts", len(s.Measures), len(s.Parts))
	}
	if !hasWarning(s.warnings, "longer than 20000 measures") {
		t.Errorf("warnings %q", s.warnings)
	}

	// Build reports what Notate left out
	s = Notate(melody(60, q), NotateOptions{})
	s.warnings = []string{"left out"}
	if res, err := engrave(t, s); err != nil || !hasWarning(res.Warnings, "left out") {
		t.Errorf("warnings %v, error %v", res, err)
	}
}

func TestNotateShortBeat(t *testing.T) {
	// a time signature that changes one tick after a bar line makes a
	// measure, and a beat, of one tick
	p := melody(60, q, 62, q, 64, h, 65, 3*q)
	p.Time = []TimeChange{{Tick: 0, Time: TimeSig{Beats: 4, BeatType: 4}}, {Tick: 1, Time: TimeSig{Beats: 3, BeatType: 4}}}
	s := Notate(p, NotateOptions{})
	if len(s.Measures) != 4 || s.Measures[0].Length != 1 || s.Measures[1].Length != 3*q {
		t.Fatalf("%d measures, the first of %d ticks", len(s.Measures), s.Measures[0].Length)
	}
	first := s.Parts[0].Measures[0].Events
	if len(first) != 1 || first[0].Tuplet != nil || first[0].Duration != 1 || len(first[0].Notes) != 1 || !first[0].Notes[0].TieStart {
		t.Errorf("the measure of one tick: %+v", first)
	}
	if res, err := engrave(t, s); err != nil || res.Notes == 0 {
		t.Errorf("build: %v", err)
	}
}

func TestClefNone(t *testing.T) {
	// a clef that is not drawn takes no room where it changes, on its
	// staff and beside the clef changes of the others
	b := newSB("Clefs")
	b.addPart("One", 1)
	b.addPart("Two", 1)
	m := b.measure(&Measure{})
	for pi, sign := range []ClefSign{ClefNone, ClefF} {
		b.ev(pi, m, 0, Half, 0, "c4")
		b.ev(pi, m, h, Half, 0, "e4")
		pm := b.s.Parts[pi].Measures[m]
		pm.Clefs = append(pm.Clefs, ClefChange{Staff: 0, Clef: Clef{Sign: ClefG}}, ClefChange{Offset: h, Staff: 0, Clef: Clef{Sign: sign}})
	}
	res, err := engrave(t, b.s)
	if err != nil || res.Notes != 4 {
		t.Fatalf("build: %v", err)
	}
}

func TestCuesBudget(t *testing.T) {
	defer func(n int) { maxCues = n }(maxCues)
	maxCues = 7
	// three measures played four times: two cues each
	s := Notate(melody(60, 4*q, 62, 4*q, 64, 4*q), NotateOptions{})
	for i := 0; i < 12; i++ {
		s.PlayOrder = append(s.PlayOrder, PlayedMeasure{Tick: i * 4 * q, Measure: i % 3})
	}
	res, err := engrave(t, s)
	if err != nil {
		t.Fatal(err)
	}
	if !hasWarning(res.Warnings, "more than 7 cues") {
		t.Errorf("warnings %q", res.Warnings)
	}
	cues, err := bdf.DecodeCues(res.Doc.Part(mustHash(t, res.Doc.Views[0].Play.Cues)).Data)
	if err != nil {
		t.Fatal(err)
	}
	if len(cues.Cues) != 7 || cues.Cues[6].Tick != uint32(3*4*q) {
		t.Errorf("%d cues, the last at %d", len(cues.Cues), cues.Cues[len(cues.Cues)-1].Tick)
	}
}

// smfEvent is an event of a track of a Standard MIDI File.
type smfEvent struct {
	tick int
	data []byte // status and data; of a meta event its type and data
}

// readSMF reads the tracks of a file written by Performance.SMF; a delta
// time of more than four bytes fails the test.
func readSMF(t *testing.T, b []byte) [][]smfEvent {
	t.Helper()
	if len(b) < 14 || string(b[:4]) != "MThd" {
		t.Fatalf("header % x", b[:min(len(b), 14)])
	}
	var out [][]smfEvent
	for b = b[14:]; len(b) >= 8; {
		n := int(binary.BigEndian.Uint32(b[4:]))
		body := b[8 : 8+n]
		b = b[8+n:]
		quantity := func() int {
			v := 0
			for i := 0; i < 4; i++ {
				c := body[0]
				body = body[1:]
				v = v<<7 | int(c&0x7F)
				if c < 0x80 {
					return v
				}
			}
			t.Fatal("a variable-length quantity of more than four bytes")
			return 0
		}
		var evs []smfEvent
		tick, running := 0, byte(0)
		for len(body) > 0 {
			tick += quantity()
			st := body[0]
			switch {
			case st == 0xFF:
				typ := body[1]
				body = body[2:]
				l := quantity()
				evs = append(evs, smfEvent{tick, append([]byte{0xFF, typ}, body[:l]...)})
				body = body[l:]
				running = 0
				continue
			case st >= 0x80:
				running = st
				body = body[1:]
			}
			l := 2
			if running&0xE0 == 0xC0 {
				l = 1
			}
			evs = append(evs, smfEvent{tick, append([]byte{running}, body[:l]...)})
			body = body[l:]
		}
		out = append(out, evs)
	}
	return out
}

func TestSMFLongWait(t *testing.T) {
	// a wait of three times what a delta time holds, and more
	late := 3*maxDelta + 5
	p := &Performance{Tracks: []*Track{{Volume: -1, Pan: -1, Program: 0, Notes: []PlayNote{
		{Tick: 0, Dur: q, Key: 60, Vel: 90}, {Tick: late, Dur: q, Key: 62, Vel: 90},
	}}}}
	tr := readSMF(t, p.SMF())
	if len(tr) != 2 {
		t.Fatalf("%d tracks", len(tr))
	}
	if last := tr[0][len(tr[0])-1]; last.tick != late+q || last.data[1] != 0x2F {
		t.Errorf("the conductor track ends with % x at %d", last.data, last.tick)
	}
	var on []int
	for _, e := range tr[1] {
		if e.data[0] == 0x90 {
			on = append(on, e.tick)
		}
	}
	if !slices.Equal(on, []int{0, late}) {
		t.Errorf("notes at %v", on)
	}

	// an end further than a quantity of five bytes holds
	p = &Performance{End: 1<<35 + 1, Tracks: []*Track{{Volume: -1, Pan: -1, Notes: []PlayNote{{Tick: 0, Dur: q, Key: 60, Vel: 90}}}}}
	tr = readSMF(t, p.SMF())
	if last := tr[1][len(tr[1])-1]; last.tick != 1<<35+1 || last.data[1] != 0x2F {
		t.Errorf("the track ends at %d", last.tick)
	}
}

func TestSMFWaits(t *testing.T) {
	// waits a delta time holds are written as they were
	w := &trackWriter{}
	w.event(maxDelta, []byte{0x90, 60, 90})
	if want := []byte{0xFF, 0xFF, 0xFF, 0x7F, 0x90, 60, 90}; !slices.Equal(w.b, want) {
		t.Errorf("% x, want % x", w.b, want)
	}
	w = &trackWriter{}
	w.event(maxDelta+1, []byte{0x90, 60, 90})
	if want := []byte{0xFF, 0xFF, 0xFF, 0x7F, 0xFF, 0x01, 0x00, 0x01, 0x90, 60, 90}; !slices.Equal(w.b, want) {
		t.Errorf("% x, want % x", w.b, want)
	}
}

package midi

import (
	"bytes"
	"encoding/binary"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/music"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

var update = flag.Bool("update", false, "rewrite the files of testdata")

// trk builds the events of a track chunk.
type trk []byte

// vlq appends a variable-length quantity.
func vlq(b []byte, v int) []byte {
	var tmp [5]byte
	i := len(tmp) - 1
	tmp[i] = byte(v & 0x7F)
	for v >>= 7; v > 0; v >>= 7 {
		i--
		tmp[i] = byte(v&0x7F) | 0x80
	}
	return append(b, tmp[i:]...)
}

// ev appends an event after delta ticks.
func (t trk) ev(delta int, b ...byte) trk { return append(trk(vlq(t, delta)), b...) }

// meta appends a meta event.
func (t trk) meta(delta int, typ byte, data string) trk {
	t = t.ev(delta, 0xFF, typ)
	return append(trk(vlq(t, len(data))), data...)
}

// end appends the end of track.
func (t trk) end(delta int) trk { return t.ev(delta, 0xFF, 0x2F, 0) }

// note appends a note-on and, dur ticks later, its note-off.
func (t trk) note(delta, ch, key, dur int) trk {
	return t.ev(delta, 0x90|byte(ch), byte(key), 100).ev(dur, 0x80|byte(ch), byte(key), 64)
}

// smf writes a Standard MIDI File.
func smf(format, division int, tracks ...trk) []byte {
	b := []byte("MThd\x00\x00\x00\x06")
	b = binary.BigEndian.AppendUint16(b, uint16(format))
	b = binary.BigEndian.AppendUint16(b, uint16(len(tracks)))
	b = binary.BigEndian.AppendUint16(b, uint16(division))
	for _, t := range tracks {
		b = append(b, "MTrk"...)
		b = binary.BigEndian.AppendUint32(b, uint32(len(t)))
		b = append(b, t...)
	}
	return b
}

func parse(t *testing.T, data []byte) *Parsed {
	t.Helper()
	p, err := Parse(data, nil)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func noWarnings(t *testing.T, p *Parsed) {
	t.Helper()
	if len(p.Warnings) > 0 {
		t.Errorf("warnings: %q", p.Warnings)
	}
}

func hasWarning(t *testing.T, p *Parsed, part string) {
	t.Helper()
	for _, w := range p.Warnings {
		if strings.Contains(w, part) {
			return
		}
	}
	t.Errorf("no warning with %q in %q", part, p.Warnings)
}

func notes(ns ...int) []music.PlayNote {
	var out []music.PlayNote
	for i := 0; i+2 < len(ns); i += 3 {
		out = append(out, music.PlayNote{Tick: ns[i], Dur: ns[i+1], Key: ns[i+2], Vel: 100})
	}
	return out
}

func checkNotes(t *testing.T, tr *music.Track, want []music.PlayNote) {
	t.Helper()
	got := slices.Clone(tr.Notes)
	for i := range got {
		got[i].Vel = 100
	}
	if !slices.Equal(got, want) {
		t.Errorf("track %q notes:\n got %v\nwant %v", tr.Name, got, want)
	}
}

func names(p *music.Performance) []string {
	var out []string
	for _, t := range p.Tracks {
		out = append(out, t.Name)
	}
	return out
}

func TestFormat0Channels(t *testing.T) {
	data := smf(0, 480, trk(nil).
		meta(0, 0x03, "Song").
		ev(0, 0xC0, 0).
		ev(0, 0xC1, 40).
		ev(0, 0xC9, 5).
		ev(0, 0x90, 60, 90).
		ev(0, 0x91, 67, 80).
		ev(0, 0x99, 36, 127).
		ev(240, 0x89, 36, 0).
		ev(240, 0x80, 60, 0).
		ev(480, 0x81, 67, 0).
		end(960))
	p := parse(t, data)
	noWarnings(t, p)
	perf := p.Perf
	if p.Format != 0 || perf.Title != "Song" || perf.Division != 480 || perf.End != 3840 {
		t.Errorf("format %d, title %q, division %d, end %d", p.Format, perf.Title, perf.Division, perf.End)
	}
	if !bytes.Equal(p.SMF, data) {
		t.Error("the SMF stored is not the input")
	}
	if got, want := names(perf), []string{"Acoustic Grand Piano", "Violin", "Drums"}; !slices.Equal(got, want) {
		t.Fatalf("tracks %q, want %q", got, want)
	}
	for i, ch := range []int{0, 1, 9} {
		if tr := perf.Tracks[i]; tr.Channel != ch || tr.Volume != -1 || tr.Pan != -1 {
			t.Errorf("track %d: channel %d, volume %d, pan %d", i, tr.Channel, tr.Volume, tr.Pan)
		}
	}
	if perf.Tracks[1].Program != 40 || perf.Tracks[2].Program != 0 {
		t.Errorf("programs %d %d", perf.Tracks[1].Program, perf.Tracks[2].Program)
	}
	checkNotes(t, perf.Tracks[0], notes(0, 960, 60))
	checkNotes(t, perf.Tracks[1], notes(0, 1920, 67))
	checkNotes(t, perf.Tracks[2], notes(0, 480, 36))
	if v := perf.Tracks[0].Notes[0].Vel; v != 90 {
		t.Errorf("velocity %d", v)
	}
}

func TestFormat1Conductor(t *testing.T) {
	conductor := trk(nil).
		meta(0, 0x03, "Twinkle").
		meta(0, 0x02, "Public domain").
		ev(0, 0xFF, 0x51, 3, 0x07, 0xA1, 0x20).   // 500000 µs: 120 BPM
		ev(0, 0xFF, 0x58, 4, 3, 2, 24, 8).        // 3/4
		ev(0, 0xFF, 0x59, 2, 0xFF, 0).            // F major
		ev(960, 0xFF, 0x51, 3, 0x06, 0x1A, 0x80). // 400000 µs: 150 BPM
		ev(1920, 0xFF, 0x58, 4, 6, 3, 24, 8).     // 6/8
		ev(0, 0xFF, 0x59, 2, 1, 1).               // E minor
		end(0)
	melody := trk(nil).
		meta(0, 0x03, "Melody").
		ev(0, 0xC0, 73).
		ev(0, 0xB0, 7, 100).
		ev(0, 0xB0, 10, 64).
		note(0, 0, 60, 960).
		ev(0, 0xB0, 7, 90).
		note(0, 0, 62, 960).
		end(0)
	bass := trk(nil).
		meta(0, 0x03, "Bass").
		note(0, 1, 36, 1920).
		end(0)
	p := parse(t, smf(1, 960, conductor, melody, bass))
	noWarnings(t, p)
	perf := p.Perf
	if perf.Title != "Twinkle" || perf.Copyright != "Public domain" {
		t.Errorf("title %q, copyright %q", perf.Title, perf.Copyright)
	}
	if got, want := names(perf), []string{"Melody", "Bass"}; !slices.Equal(got, want) {
		t.Fatalf("tracks %q, want %q", got, want)
	}
	m, b := perf.Tracks[0], perf.Tracks[1]
	if m.Program != 73 || m.Volume != 100 || m.Pan != 64 || b.Program != 0 || b.Volume != -1 || b.Pan != -1 {
		t.Errorf("melody %d/%d/%d, bass %d/%d/%d", m.Program, m.Volume, m.Pan, b.Program, b.Volume, b.Pan)
	}
	if len(m.Controls) != 4 || m.Controls[3] != (music.Control{Tick: 960, Kind: music.ControlChange, Num: 7, Value: 90}) {
		t.Errorf("controls %v", m.Controls)
	}
	checkNotes(t, m, notes(0, 960, 60, 960, 960, 62))
	checkNotes(t, b, notes(0, 1920, 36))
	if want := []music.Tempo{{Tick: 0, BPM: 120}, {Tick: 960, BPM: 150}}; !slices.Equal(perf.Tempo, want) {
		t.Errorf("tempo %v", perf.Tempo)
	}
	if want := []music.TimeChange{{Tick: 0, Time: music.TimeSig{Beats: 3, BeatType: 4}}, {Tick: 2880, Time: music.TimeSig{Beats: 6, BeatType: 8}}}; !slices.Equal(perf.Time, want) {
		t.Errorf("time %v", perf.Time)
	}
	if want := []music.KeyChange{{Tick: 0, Key: music.KeySig{Fifths: -1}}, {Tick: 2880, Key: music.KeySig{Fifths: 1, Minor: true}}}; !slices.Equal(perf.Key, want) {
		t.Errorf("key %v", perf.Key)
	}
	if perf.End != 2880 {
		t.Errorf("end %d", perf.End)
	}
}

func TestRunningStatus(t *testing.T) {
	data := smf(0, 960, trk(nil).
		ev(0, 0x90, 60, 100).
		ev(240, 60, 0). // running status, velocity 0: off
		ev(0, 62, 100).
		meta(0, 0x01, "text").
		ev(240, 62, 0).                   // running status after a meta event
		ev(0, 0xF0, 3, 0x7E, 0x7F, 0xF7). // system exclusive
		ev(0, 64, 100).                   // and after a system exclusive
		ev(480, 0x80, 64, 0).
		end(0))
	p := parse(t, data)
	noWarnings(t, p)
	checkNotes(t, p.Perf.Tracks[0], notes(0, 240, 60, 240, 240, 62, 480, 480, 64))
}

func TestOverlappingNotes(t *testing.T) {
	// two notes of one key: the first off ends the first note
	data := smf(0, 960, trk(nil).
		ev(0, 0x90, 60, 100).
		ev(240, 0x90, 60, 100).
		ev(240, 0x80, 60, 0).
		ev(240, 0x90, 60, 0).
		ev(0, 0x90, 62, 100).
		ev(0, 0x90, 64, 100).
		ev(240, 0xB0, 123, 0). // all notes off
		end(0))
	p := parse(t, data)
	noWarnings(t, p)
	checkNotes(t, p.Perf.Tracks[0], notes(0, 480, 60, 240, 480, 60, 720, 240, 62, 720, 240, 64))
}

func TestUnendedNotes(t *testing.T) {
	data := smf(0, 960, trk(nil).
		ev(0, 0x90, 60, 100).
		note(0, 0, 64, 480).
		ev(0, 0x90, 67, 100).
		ev(0, 0x80, 67, 0). // no length: left out
		end(1440))
	p := parse(t, data)
	hasWarning(t, p, "without a note-off")
	hasWarning(t, p, "no length")
	checkNotes(t, p.Perf.Tracks[0], notes(0, 1920, 60, 0, 480, 64))
}

func TestDivision(t *testing.T) {
	for _, tc := range []struct {
		division int
		ticks    []int // absolute ticks of note-on, note-off pairs
		want     []int // PPQ ticks
	}{
		{480, []int{0, 480, 720, 960}, []int{0, 960, 1440, 1920}},
		{96, []int{0, 24, 96, 97}, []int{0, 240, 960, 970}},
		// 9.6 PPQ ticks a tick, rounded from the absolute tick
		{100, []int{1, 3, 101, 203, 1001, 1100}, []int{10, 29, 970, 1949, 9610, 10560}},
		{1920, []int{1, 2, 3, 5}, []int{1, 1, 2, 3}},
	} {
		tr := trk(nil)
		at := 0
		for i := 0; i < len(tc.ticks); i += 2 {
			tr = tr.ev(tc.ticks[i]-at, 0x90, 60, 100).ev(tc.ticks[i+1]-tc.ticks[i], 0x80, 60, 0)
			at = tc.ticks[i+1]
		}
		p := parse(t, smf(0, tc.division, tr.end(0)))
		var got []int
		for _, n := range p.Perf.Tracks[0].Notes {
			got = append(got, n.Tick, n.Tick+max(n.Dur, 0))
		}
		want := slices.Clone(tc.want)
		if tc.division == 1920 {
			want[1] = 2 // at least a tick long
		}
		if !slices.Equal(got, want) {
			t.Errorf("division %d: %v, want %v", tc.division, got, want)
		}
		if p.Perf.Division != tc.division {
			t.Errorf("division %d: Division %d", tc.division, p.Perf.Division)
		}
	}
}

func TestMetaMerge(t *testing.T) {
	// changes at one tick: the later chunk wins; a change to the value in
	// effect is dropped
	a := trk(nil).
		ev(0, 0xFF, 0x51, 3, 0x07, 0xA1, 0x20). // 120
		ev(0, 0xFF, 0x58, 4, 4, 2, 24, 8).
		ev(1920, 0xFF, 0x58, 4, 4, 2, 24, 8).
		ev(0, 0xFF, 0x51, 3, 0x09, 0x27, 0xC0). // 100
		end(0)
	b := trk(nil).
		ev(0, 0xFF, 0x51, 3, 0x0A, 0x2C, 0x2B). // 90
		note(0, 0, 60, 3840).
		end(0)
	p := parse(t, smf(1, 480, a, b))
	noWarnings(t, p)
	if want := []music.Tempo{{Tick: 0, BPM: 90}, {Tick: 3840, BPM: 100}}; !slices.Equal(p.Perf.Tempo, want) {
		t.Errorf("tempo %v", p.Perf.Tempo)
	}
	if len(p.Perf.Time) != 1 {
		t.Errorf("time %v", p.Perf.Time)
	}
}

func TestLyricsConductor(t *testing.T) {
	// lyrics in a track without notes go to the melody: the track with
	// notes where they are sung
	words := trk(nil).
		meta(0, 0x05, "Twin").
		meta(450, 0x05, "kle"). // 30 ticks early
		meta(510, 0x05, "star").
		end(0)
	accomp := trk(nil).
		meta(0, 0x03, "Piano").
		note(0, 0, 48, 1920).
		end(0)
	melody := trk(nil).
		meta(0, 0x03, "Voice").
		note(0, 1, 60, 480).
		note(0, 1, 60, 480).
		note(0, 1, 67, 960).
		end(0)
	p := parse(t, smf(1, 480, words, accomp, melody))
	noWarnings(t, p)
	if l := p.Perf.Tracks[0].Lyrics; len(l) != 0 {
		t.Errorf("piano lyrics %v", l)
	}
	want := []music.Lyric{{Tick: 0, Text: "Twin"}, {Tick: 960, Text: "kle"}, {Tick: 1920, Text: "star"}}
	if got := p.Perf.Tracks[1].Lyrics; !slices.Equal(got, want) {
		t.Errorf("lyrics %v, want %v", got, want)
	}
}

func TestLyricsOwnChunk(t *testing.T) {
	// lyrics of a format 0 file go to the channel whose notes they start
	tr := trk(nil).
		meta(0, 0x05, "Ah").
		ev(0, 0x90, 60, 100).
		ev(0, 0x91, 48, 100).
		ev(0, 0x99, 36, 100).
		ev(480, 0x80, 60, 0).
		ev(0, 0x89, 36, 0).
		meta(0, 0x05, "ve").
		ev(0, 0x90, 62, 100).
		ev(0, 0x99, 36, 100).
		ev(480, 0x80, 62, 0).
		ev(0, 0x89, 36, 0).
		ev(0, 0x81, 48, 0).
		meta(0, 0x05, "Ma").
		note(0, 0, 64, 480).
		meta(0, 0x05, "ri").
		ev(0, 0x90, 65, 100).
		ev(0, 0x91, 50, 100).
		ev(480, 0x80, 65, 0).
		ev(0, 0x81, 50, 0).
		meta(0, 0x05, "a").
		note(0, 0, 67, 960).
		end(0)
	p := parse(t, smf(0, 480, tr))
	noWarnings(t, p)
	if got, want := names(p.Perf), []string{"Channel 1", "Channel 2", "Drums"}; !slices.Equal(got, want) {
		t.Fatalf("tracks %q, want %q", got, want)
	}
	var texts []string
	for _, l := range p.Perf.Tracks[0].Lyrics {
		texts = append(texts, l.Text)
	}
	if got, want := strings.Join(texts, " "), "Ah ve Ma ri a"; got != want {
		t.Errorf("lyrics %q, want %q", got, want)
	}
}

func TestLyricsSpacing(t *testing.T) {
	// words written with spaces (karaoke style) get hyphens between the
	// syllables of a word; kana do not
	tr := trk(nil)
	for _, s := range []string{"Twin", "kle ", "twin", "kle ", "lit", "tle ", "star", "\\きら", "きら"} {
		tr = tr.meta(0, 0x05, s).note(0, 0, 60, 480)
	}
	p := parse(t, smf(0, 480, tr.end(0)))
	var texts []string
	for _, l := range p.Perf.Tracks[0].Lyrics {
		texts = append(texts, l.Text)
	}
	if got, want := strings.Join(texts, "|"), "Twin-|kle|twin-|kle|lit-|tle|star|きら|きら"; got != want {
		t.Errorf("lyrics %q, want %q", got, want)
	}
}

func TestKaraoke(t *testing.T) {
	p := parse(t, twinkleKar())
	noWarnings(t, p)
	perf := p.Perf
	if perf.Title != "Twinkle, Twinkle, Little Star" {
		t.Errorf("title %q", perf.Title)
	}
	if got, want := names(perf), []string{"Melody", "Piano"}; !slices.Equal(got, want) {
		t.Fatalf("tracks %q, want %q", got, want)
	}
	var texts []string
	for _, l := range perf.Tracks[0].Lyrics {
		texts = append(texts, l.Text)
	}
	if got, want := strings.Join(texts, " "), "Twin- kle, twin- kle, lit- tle star. How I won- der what you are."; got != want {
		t.Errorf("lyrics\n got %q\nwant %q", got, want)
	}
	if len(perf.Tracks[1].Lyrics) != 0 {
		t.Error("lyrics on the piano")
	}
	if n := len(perf.Tracks[0].Lyrics); n != len(perf.Tracks[0].Notes) {
		t.Errorf("%d syllables, %d notes", n, len(perf.Tracks[0].Notes))
	}
	for i, l := range perf.Tracks[0].Lyrics {
		if l.Tick != perf.Tracks[0].Notes[i].Tick {
			t.Errorf("syllable %q at %d, note at %d", l.Text, l.Tick, perf.Tracks[0].Notes[i].Tick)
		}
	}
}

func TestEncodings(t *testing.T) {
	sjis := func(s string) string {
		b, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(s))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	latin := func(s string) string {
		b, err := charmap.Windows1252.NewEncoder().Bytes([]byte(s))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	for _, tc := range []struct {
		enc                            func(string) string
		title, track, credit, composer string
	}{
		{sjis, "荒城の月", "ﾋﾟｱﾉ", "作曲：滝廉太郎", "滝廉太郎"},
		{latin, "Für Elise", "Piano à queue", "Composer: Ludwig van Beethoven", "Ludwig van Beethoven"},
		{func(s string) string { return s }, "Ода к радости", "ピアノ", "作曲 : ベートーヴェン", "ベートーヴェン"},
	} {
		conductor := trk(nil).
			meta(0, 0x03, tc.enc(tc.title)+"\x00").
			meta(0, 0x01, tc.enc(tc.credit)).
			end(0)
		piano := trk(nil).
			meta(0, 0x03, tc.enc(tc.track)+"  ").
			note(0, 0, 60, 480).
			end(0)
		p := parse(t, smf(1, 480, conductor, piano))
		noWarnings(t, p)
		if p.Perf.Title != tc.title || p.Perf.Tracks[0].Name != tc.track || p.Perf.Composer != tc.composer {
			t.Errorf("title %q, track %q, composer %q; want %q, %q, %q", p.Perf.Title, p.Perf.Tracks[0].Name, p.Perf.Composer, tc.title, tc.track, tc.composer)
		}
	}
}

func TestCredits(t *testing.T) {
	for s, want := range map[string][2]string{
		"Composer: J. S. Bach":    {"composer", "J. S. Bach"},
		"composed by Beethoven":   {"composer", "Beethoven"},
		"Music by Stephen Foster": {"composer", "Stephen Foster"},
		"MUSIC : Mozart":          {"composer", "Mozart"},
		"Words: Jane Taylor":      {"lyricist", "Jane Taylor"},
		"作詞：土井晩翠":                 {"lyricist", "土井晩翠"},
		"編曲: 山田":                  {"arranger", "山田"},
		"Arranged by":             {},
		"Musical box":             {},
		"Composer:":               {},
		"Composedby X":            {},
		"The composer: X":         {},
	} {
		field, v := credit(s)
		if field != want[0] || v != want[1] {
			t.Errorf("credit(%q) = %q, %q", s, field, v)
		}
	}
}

func TestChannelPrefix(t *testing.T) {
	data := smf(0, 480, trk(nil).
		meta(0, 0x03, "Title").
		meta(0, 0x20, "\x01").
		meta(0, 0x03, "Strings").
		ev(0, 0xC1, 48).
		ev(0, 0xC2, 56).
		note(0, 1, 60, 480).
		note(0, 2, 72, 480).
		end(0))
	p := parse(t, data)
	if got, want := names(p.Perf), []string{"Strings", "Trumpet"}; !slices.Equal(got, want) {
		t.Errorf("tracks %q, want %q", got, want)
	}
	if p.Perf.Title != "Title" {
		t.Errorf("title %q", p.Perf.Title)
	}
}

func TestProgramFromOtherTrack(t *testing.T) {
	// a setup track sets the programs of the channels
	setup := trk(nil).
		meta(0, 0x03, "Setup").
		ev(0, 0xC0, 24).
		ev(0, 0xB0, 7, 80).
		end(0)
	guitar := trk(nil).
		note(10, 0, 52, 480).
		end(0)
	p := parse(t, smf(1, 480, setup, guitar))
	tr := p.Perf.Tracks[0]
	if tr.Program != 24 || tr.Volume != 80 || tr.Name != "Acoustic Guitar (nylon)" {
		t.Errorf("program %d, volume %d, name %q", tr.Program, tr.Volume, tr.Name)
	}
	if p.Perf.Title != "" {
		t.Errorf("title %q", p.Perf.Title)
	}
}

func TestRMID(t *testing.T) {
	inner := smf(0, 480, trk(nil).note(0, 0, 60, 480).end(0))
	data := rmid(inner, "Ode")
	if !detect(data) {
		t.Error("not detected")
	}
	p := parse(t, data)
	noWarnings(t, p)
	if !bytes.Equal(p.SMF, inner) {
		t.Error("the SMF is not the data chunk")
	}
	if p.Perf.Title != "Ode" {
		t.Errorf("title %q", p.Perf.Title)
	}
	if _, err := Parse([]byte("RIFF\x04\x00\x00\x00RMID"), nil); err == nil {
		t.Error("no error for a RIFF MIDI file without data")
	}
}

func TestSMPTE(t *testing.T) {
	// 25 frames of 40 ticks: 1000 ticks a second, 500 a quarter note at 120
	data := smf(0, 0xE728, trk(nil).
		ev(0, 0xFF, 0x51, 3, 0x0F, 0x42, 0x40). // ignored
		note(1000, 0, 60, 500).
		end(0))
	p := parse(t, data)
	hasWarning(t, p, "SMPTE")
	if p.SMF != nil || p.Perf.Division != 0 {
		t.Errorf("SMF %d bytes, division %d", len(p.SMF), p.Perf.Division)
	}
	if want := []music.Tempo{{Tick: 0, BPM: 120}}; !slices.Equal(p.Perf.Tempo, want) {
		t.Errorf("tempo %v", p.Perf.Tempo)
	}
	checkNotes(t, p.Perf.Tracks[0], notes(1920, 960, 60))
}

func TestFormat2(t *testing.T) {
	a := trk(nil).note(0, 0, 60, 480).end(0)
	b := trk(nil).note(0, 0, 64, 480).end(0)
	data := smf(2, 480, a, b)
	p := parse(t, data)
	hasWarning(t, p, "format 2")
	if len(p.Perf.Tracks) != 2 || p.Perf.Tracks[1].Notes[0].Tick != 0 {
		t.Fatalf("tracks %v", p.Perf.Tracks)
	}
	if p.SMF[9] != 1 || !bytes.Equal(p.SMF[10:], data[10:]) {
		t.Errorf("SMF %x", p.SMF[:14])
	}
}

func TestDamaged(t *testing.T) {
	good := smf(1, 480,
		trk(nil).meta(0, 0x03, "Song").end(0),
		trk(nil).note(0, 0, 60, 480).note(0, 0, 62, 480).note(0, 0, 64, 480).end(0))
	for _, tc := range []struct {
		name    string
		data    []byte
		warning string
		notes   []music.PlayNote
	}{
		// the note-off of the last note is cut: the note has no length
		{"truncated", good[:len(good)-6], "cut short", notes(0, 960, 60, 960, 960, 62)},
		{"count", append(slices.Clone(good[:10]), append([]byte{0, 3}, good[12:]...)...), "counts 3 tracks", nil},
		{"trailing", append(slices.Clone(good), "garbage"...), "bytes after the last track", nil},
		{"junk between", slices.Concat(good[:34], []byte("\x01\x02\x03"), good[34:]), "not a chunk", nil},
		{"format 0", append(slices.Clone(good[:9]), append([]byte{0}, good[10:]...)...), "format 0 file of 2 tracks", nil},
		{"no status", smf(0, 480, trk(nil).ev(0, 60, 64).end(0)), "without a status", []music.PlayNote{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := parse(t, tc.data)
			hasWarning(t, p, tc.warning)
			if bytes.Equal(p.SMF, tc.data) {
				t.Error("the damaged file is stored")
			}
			want := tc.notes
			if want == nil {
				want = notes(0, 480, 60, 480, 480, 62, 960, 480, 64)
				for i := range want {
					want[i].Tick *= 2
					want[i].Dur *= 2
				}
			}
			var got []music.PlayNote
			for _, tr := range p.Perf.Tracks {
				got = append(got, tr.Notes...)
			}
			for i := range got {
				got[i].Vel = 100
			}
			if len(want) == 0 {
				want = nil
			}
			if !slices.Equal(got, want) {
				t.Errorf("notes %v, want %v", got, want)
			}
			// the repaired file reads without damage
			q := parse(t, p.SMF)
			for _, w := range q.Warnings {
				if !strings.Contains(w, "note-off") && !strings.Contains(w, "no length") {
					t.Errorf("repaired file: %s", w)
				}
			}
		})
	}
}

func TestErrors(t *testing.T) {
	for name, data := range map[string][]byte{
		"empty":       nil,
		"not midi":    []byte("RIFFxxxxWAVEfmt "),
		"header":      []byte("MThd\x00\x00\x00\x06\x00\x01"),
		"short":       []byte("MThd\x00\x00\x00\x04\x00\x01\x00\x01"),
		"no tracks":   []byte("MThd\x00\x00\x00\x06\x00\x01\x00\x00\x01\xE0"),
		"division 0":  []byte("MThd\x00\x00\x00\x06\x00\x00\x00\x01\x00\x00MTrk\x00\x00\x00\x04\x00\xFF\x2F\x00"),
		"bad SMPTE":   []byte("MThd\x00\x00\x00\x06\x00\x00\x00\x01\x80\x10MTrk\x00\x00\x00\x04\x00\xFF\x2F\x00"),
		"many tracks": append([]byte("MThd\x00\x00\x00\x06\x00\x01\x04\x01\x01\xE0"), bytes.Repeat([]byte("MTrk\x00\x00\x00\x04\x00\xFF\x2F\x00"), maxChunks+1)...),
	} {
		if _, err := Parse(data, nil); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

func TestDetect(t *testing.T) {
	for _, tc := range []struct {
		head string
		want bool
	}{
		{"MThd\x00\x00\x00\x06\x00\x01", true},
		{"MThd\x00\x00\x00\x08", true},
		{"MThd\x00\x00\x00\x02", false},
		{"MThd", false},
		{"RIFF\x00\x10\x00\x00RMIDdata", true},
		{"RIFF\x00\x10\x00\x00WAVEfmt ", false},
		{"MTrk\x00\x00\x00\x06", false},
	} {
		if got := detect([]byte(tc.head)); got != tc.want {
			t.Errorf("detect(%q) = %v", tc.head, got)
		}
	}
	f := conv.Lookup("midi")
	if f == nil || !slices.Contains(f.Extensions, ".kar") || !f.Detect([]byte("MThd\x00\x00\x00\x06"), nil, 0) {
		t.Error("the format is not registered")
	}
}

// twinkle is "Twinkle, Twinkle, Little Star" (a traditional melody, words
// by Jane Taylor, 1806): the melody and its syllables.
var twinkle = []struct {
	key, beats int
	syllable   string
}{
	{60, 1, "Twin"}, {60, 1, "kle, "}, {67, 1, "twin"}, {67, 1, "kle, "},
	{69, 1, "lit"}, {69, 1, "tle "}, {67, 2, "star. "},
	{65, 1, "/How "}, {65, 1, "I "}, {64, 1, "won"}, {64, 1, "der "},
	{62, 1, "what "}, {62, 1, "you "}, {60, 2, "are."},
}

// twinklePerformance is Twinkle with a bass, as a performance.
func twinklePerformance() *music.Performance {
	melody := &music.Track{Name: "Melody", Channel: 0, Program: 73, Volume: 100, Pan: 64}
	bass := &music.Track{Name: "Bass", Channel: 1, Program: 32, Volume: 90, Pan: -1}
	tick := 0
	for _, n := range twinkle {
		melody.Notes = append(melody.Notes, music.PlayNote{Tick: tick, Dur: n.beats*music.PPQ - 60, Key: n.key, Vel: 96})
		s := strings.TrimSpace(strings.TrimPrefix(n.syllable, "/"))
		if !strings.HasSuffix(n.syllable, " ") && !strings.HasSuffix(s, ".") {
			s += "-"
		}
		melody.Lyrics = append(melody.Lyrics, music.Lyric{Tick: tick, Text: s})
		tick += n.beats * music.PPQ
	}
	for i, k := range []int{48, 52, 53, 48, 53, 48, 55, 48} {
		bass.Notes = append(bass.Notes, music.PlayNote{Tick: i * 2 * music.PPQ, Dur: 2 * music.PPQ, Key: k, Vel: 80})
	}
	bass.Controls = []music.Control{{Tick: 4 * music.PPQ, Kind: music.ControlChange, Num: 64, Value: 127}}
	return &music.Performance{
		Title:  "Twinkle, Twinkle, Little Star",
		Tempo:  []music.Tempo{{Tick: 0, BPM: 100}, {Tick: 12 * music.PPQ, BPM: 80}},
		Time:   []music.TimeChange{{Tick: 0, Time: music.TimeSig{Beats: 4, BeatType: 4}}},
		Key:    []music.KeyChange{{Tick: 0, Key: music.KeySig{}}},
		Tracks: []*music.Track{melody, bass},
		End:    16 * music.PPQ,
	}
}

// twinkleKar is Twinkle as a karaoke file: the syllables are text
// events of a track named Words.
func twinkleKar() []byte {
	head := trk(nil).
		meta(0, 0x01, "@KMIDI KARAOKE FILE").
		meta(0, 0x01, "@V0100").
		meta(0, 0x01, "@Itraditional").
		ev(0, 0xFF, 0x51, 3, 0x07, 0xA1, 0x20).
		end(0)
	words := trk(nil).
		meta(0, 0x03, "Words").
		meta(0, 0x01, "@LENGL").
		meta(0, 0x01, "@TTwinkle, Twinkle, Little Star").
		meta(0, 0x01, "@TTraditional")
	melody := trk(nil).meta(0, 0x03, "Melody").ev(0, 0xC0, 73)
	delta, rest := 0, 0
	for i, n := range twinkle {
		s := n.syllable
		if i == 0 {
			s = "\\" + s
		}
		words = words.meta(delta, 0x01, s)
		melody = melody.ev(rest, 0x90, byte(n.key), 100).ev(n.beats*480-40, 0x80, byte(n.key), 0)
		delta, rest = n.beats*480, 40
	}
	piano := trk(nil).meta(0, 0x03, "Piano")
	for i := 0; i < 4; i++ {
		piano = piano.ev(0, 0x91, 48, 70).ev(0, 64, 70).ev(0, 67, 70).
			ev(1920, 48, 0).ev(0, 64, 0).ev(0, 67, 0)
	}
	return smf(1, 480, head, words.end(0), melody.end(40), piano.end(0))
}

// rmid wraps a Standard MIDI File in a RIFF MIDI file with a title.
func rmid(smf []byte, title string) []byte {
	riff := func(id string, body []byte) []byte {
		b := append([]byte(id), binary.LittleEndian.AppendUint32(nil, uint32(len(body)))...)
		b = append(b, body...)
		if len(body)%2 == 1 {
			b = append(b, 0)
		}
		return b
	}
	info := append([]byte("INFO"), riff("INAM", []byte(title+"\x00"))...)
	body := append([]byte("RMID"), riff("LIST", info)...)
	return riff("RIFF", append(body, riff("data", smf)...))
}

// odeFormat0 is the beginning of "Ode to Joy" (Beethoven, 1824) as a
// format 0 file of 96 ticks a quarter note: flute, piano and drums on
// three channels of one track, written with running status.
func odeFormat0() []byte {
	tr := trk(nil).
		meta(0, 0x03, "Ode to Joy").
		meta(0, 0x02, "Public domain").
		ev(0, 0xFF, 0x51, 3, 0x07, 0xA1, 0x20).
		ev(0, 0xFF, 0x58, 4, 4, 2, 24, 8).
		ev(0, 0xFF, 0x59, 2, 2, 0). // D major
		ev(0, 0xC0, 73).
		ev(0, 0xC1, 0).
		ev(0, 0xB0, 7, 100)
	melody := []int{66, 66, 67, 69, 69, 67, 66, 64, 62, 62, 64, 66, 66, 64, 64}
	lengths := []int{96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 144, 48, 192}
	chords := [][]int{{50, 54, 57}, {55, 59, 62}, {50, 54, 57}, {45, 49, 57}}
	type ev struct {
		tick int
		b    []byte
	}
	var evs []ev
	at := 0
	for i, k := range melody {
		evs = append(evs, ev{at, []byte{0x90, byte(k), 100}}, ev{at + lengths[i] - 8, []byte{0x90, byte(k), 0}})
		at += lengths[i]
	}
	for m, c := range chords {
		for _, k := range c {
			evs = append(evs, ev{m * 384, []byte{0x91, byte(k), 70}}, ev{m*384 + 384, []byte{0x81, byte(k), 0}})
		}
		for b := 0; b < 4; b++ {
			key := byte(36)
			if b%2 == 1 {
				key = 38
			}
			evs = append(evs, ev{m*384 + b*96, []byte{0x99, key, 110}}, ev{m*384 + b*96 + 24, []byte{0x99, key, 0}})
		}
	}
	slices.SortStableFunc(evs, func(a, b ev) int { return a.tick - b.tick })
	var running byte
	last := 0
	for _, e := range evs {
		b := e.b
		if b[0] == running {
			b = b[1:]
		}
		running = e.b[0]
		tr = tr.ev(e.tick-last, b...)
		last = e.tick
	}
	return smf(0, 96, tr.end(0))
}

func TestRoundTrip(t *testing.T) {
	// a performance written as a file reads back the same
	want := twinklePerformance()
	p := parse(t, want.SMF())
	noWarnings(t, p)
	got := p.Perf
	if got.Title != want.Title || !slices.Equal(got.Tempo, want.Tempo) || !slices.Equal(got.Time, want.Time) ||
		!slices.Equal(got.Key, want.Key) || got.End != want.End || got.Division != music.PPQ {
		t.Errorf("got %q %v %v %v end %d division %d", got.Title, got.Tempo, got.Time, got.Key, got.End, got.Division)
	}
	if len(got.Tracks) != len(want.Tracks) {
		t.Fatalf("%d tracks", len(got.Tracks))
	}
	for i, w := range want.Tracks {
		g := got.Tracks[i]
		if g.Name != w.Name || g.Channel != w.Channel || g.Program != w.Program || g.Volume != w.Volume || g.Pan != w.Pan {
			t.Errorf("track %d: %q ch %d program %d volume %d pan %d", i, g.Name, g.Channel, g.Program, g.Volume, g.Pan)
		}
		if !slices.Equal(g.Notes, w.Notes) {
			t.Errorf("track %d notes:\n got %v\nwant %v", i, g.Notes, w.Notes)
		}
		if !slices.Equal(g.Lyrics, w.Lyrics) {
			t.Errorf("track %d lyrics:\n got %v\nwant %v", i, g.Lyrics, w.Lyrics)
		}
		for _, c := range w.Controls {
			if !slices.Contains(g.Controls, c) {
				t.Errorf("track %d: control %v is lost in %v", i, c, g.Controls)
			}
		}
	}
}

// The files of testdata are written by the functions above:
//
//	go test ./converter/midi -run TestTestdata -update
func TestTestdata(t *testing.T) {
	files := map[string][]byte{
		"twinkle.mid": twinklePerformance().SMF(),
		"twinkle.kar": twinkleKar(),
		"ode.mid":     odeFormat0(),
		"ode.rmi":     rmid(odeFormat0(), "Ode to Joy"),
	}
	if *update {
		for name, b := range files {
			if err := os.WriteFile(filepath.Join("testdata", name), b, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	read := func(name string) *Parsed {
		t.Helper()
		b, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		p := parse(t, b)
		noWarnings(t, p)
		return p
	}

	p := read("twinkle.mid")
	if p.Perf.Title != "Twinkle, Twinkle, Little Star" || len(p.Perf.Tracks) != 2 || len(p.Perf.Tracks[0].Lyrics) != 14 {
		t.Errorf("twinkle.mid: %q, %d tracks", p.Perf.Title, len(p.Perf.Tracks))
	}
	p = read("twinkle.kar")
	if p.Perf.Title != "Twinkle, Twinkle, Little Star" || len(p.Perf.Tracks[0].Lyrics) != 14 {
		t.Errorf("twinkle.kar: %q, %d syllables", p.Perf.Title, len(p.Perf.Tracks[0].Lyrics))
	}
	for _, name := range []string{"ode.mid", "ode.rmi"} {
		p = read(name)
		perf := p.Perf
		if got, want := names(perf), []string{"Flute", "Acoustic Grand Piano", "Drums"}; !slices.Equal(got, want) {
			t.Fatalf("%s: tracks %q, want %q", name, got, want)
		}
		if perf.Title != "Ode to Joy" || perf.Copyright != "Public domain" || perf.Division != 96 || p.Format != 0 ||
			len(perf.Tracks[0].Notes) != 15 || len(perf.Key) != 1 || perf.Key[0].Key != (music.KeySig{Fifths: 2}) || perf.End != 15360 {
			t.Errorf("%s: %q %q division %d format %d, %d notes, key %v, end %d", name, perf.Title, perf.Copyright,
				perf.Division, p.Format, len(perf.Tracks[0].Notes), perf.Key, perf.End)
		}
		if n := perf.Tracks[0].Notes[12]; n != (music.PlayNote{Tick: 11520, Dur: 1360, Key: 66, Vel: 100}) {
			t.Errorf("%s: the dotted note is %v", name, n)
		}
		if !bytes.Equal(p.SMF, files["ode.mid"]) {
			t.Errorf("%s: the SMF stored is not the file", name)
		}
	}
}

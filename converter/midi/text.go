package midi

import (
	"bytes"
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf/converter/internal/music"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

// decoder decodes the text of a file: UTF-8 as it is, and the other
// strings in Shift_JIS when all of them are Shift_JIS, else Windows-1252.
type decoder struct {
	legacy encoding.Encoding
}

func newDecoder(chunks []*chunk, info riffInfo) *decoder {
	d := &decoder{legacy: japanese.ShiftJIS}
	sjis := func(b []byte) bool {
		b = dropNUL(b)
		if utf8.Valid(b) {
			return true
		}
		_, ok := decodeSJIS(b)
		return ok
	}
	ok := sjis(info.title) && sjis(info.copyright)
	for _, c := range chunks {
		if !ok {
			break
		}
		ok = sjis(c.name) && sjis(c.instrument) && sjis(c.copyright)
		for _, n := range c.chanNames {
			ok = ok && sjis(n)
		}
		for _, t := range c.texts {
			ok = ok && sjis(t.data)
		}
	}
	if !ok {
		d.legacy = charmap.Windows1252
	}
	return d
}

// decodeSJIS decodes Shift_JIS, failing on bytes that are not.
func decodeSJIS(b []byte) (string, bool) {
	out, err := japanese.ShiftJIS.NewDecoder().Bytes(b)
	if err != nil {
		return "", false
	}
	for _, r := range string(out) {
		if r == utf8.RuneError || r >= 0x80 && r < 0xA0 {
			return "", false
		}
	}
	return string(out), true
}

func dropNUL(b []byte) []byte {
	if bytes.IndexByte(b, 0) < 0 {
		return b
	}
	return bytes.ReplaceAll(b, []byte{0}, nil)
}

// text decodes a string as it is (without NULs and a byte order mark).
func (d *decoder) text(b []byte) string {
	b = dropNUL(b)
	var s string
	switch {
	case utf8.Valid(b):
		s = string(b)
	case d.legacy == japanese.ShiftJIS:
		s, _ = decodeSJIS(b)
	default:
		out, _ := d.legacy.NewDecoder().Bytes(b)
		s = string(out)
	}
	return strings.TrimPrefix(s, "\ufeff")
}

// name decodes a name: trimmed, without control characters.
func (d *decoder) name(b []byte) string {
	s := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, d.text(b))
	return strings.TrimSpace(s)
}

// trackName names the track of a channel of a chunk of nchans channels:
// the channel's name, the chunk's track name (unless it is the title) or
// instrument name, else the General MIDI program when one is set.
func (d *decoder) trackName(c *chunk, ch, nchans int, hasProgram bool, program int, titled bool) string {
	if s := d.name(c.chanNames[ch]); s != "" {
		return s
	}
	var s string
	if !titled {
		s = d.name(c.name)
	}
	if s == "" {
		s = d.name(c.instrument)
	}
	switch {
	case s != "" && nchans > 1:
		return fmt.Sprintf("%s (ch %d)", s, ch+1)
	case s != "":
		return s
	case ch == 9:
		return "Drums"
	case hasProgram:
		return music.ProgramName(program)
	}
	return fmt.Sprintf("Channel %d", ch+1)
}

// placeholders are track names that do not name the music.
var placeholders = []string{"untitled", "tempo", "tempo track", "tempo map", "conductor", "conductor track",
	"setup", "system setup", "global", "sequence", "track 0", "track 1"}

// credits reads the title, the copyright and the credits.
func (p *parser) credits(perf *music.Performance, d *decoder, chunks []*chunk, format int, info riffInfo) {
	// the title: @T of a karaoke file, the name of the sequence (the first
	// track of format 0, a conductor track), the RIFF INAM
	for _, c := range chunks {
		for _, t := range c.texts {
			if t.typ == 1 && perf.Title == "" && bytes.HasPrefix(t.data, []byte("@T")) {
				perf.Title = d.name(t.data[2:])
			}
		}
	}
	c0 := chunks[0]
	if s := d.name(c0.name); perf.Title == "" && (format == 0 || !c0.hasNotes()) &&
		!slices.ContainsFunc(placeholders, func(p string) bool { return strings.EqualFold(p, s) }) {
		perf.Title = s
	}
	if perf.Title == "" {
		perf.Title = d.name(info.title)
	}

	// "Composer: …" and the like in text events
	for _, c := range chunks {
		for _, t := range c.texts {
			if t.typ != 1 {
				continue
			}
			field, v := credit(d.name(t.data))
			switch {
			case field == "composer" && perf.Composer == "":
				perf.Composer = v
			case field == "lyricist" && perf.Lyricist == "":
				perf.Lyricist = v
			case field == "arranger" && perf.Arranger == "":
				perf.Arranger = v
			}
		}
	}

	// else the first text of the first track
	if perf.Title == "" && !isKaraoke(c0, d) {
		for _, t := range c0.texts {
			if t.typ != 1 || t.tick != 0 {
				continue
			}
			if s := d.name(t.data); s != "" && s[0] != '@' {
				if field, _ := credit(s); field == "" {
					perf.Title = s
					break
				}
			}
		}
	}

	for _, c := range chunks {
		if s := d.name(c.copyright); s != "" {
			perf.Copyright = s
			break
		}
	}
	if perf.Copyright == "" {
		perf.Copyright = d.name(info.copyright)
	}
}

var creditLabels = []struct{ label, field string }{
	{"composer", "composer"}, {"composed by", "composer"}, {"music", "composer"}, {"music by", "composer"}, {"作曲", "composer"},
	{"lyricist", "lyricist"}, {"lyrics", "lyricist"}, {"lyrics by", "lyricist"}, {"words", "lyricist"}, {"words by", "lyricist"}, {"作詞", "lyricist"},
	{"arranger", "arranger"}, {"arranged by", "arranger"}, {"arrangement", "arranger"}, {"編曲", "arranger"},
}

// credit reads a text such as "Composer: J. S. Bach" or "作曲：滝廉太郎".
func credit(s string) (field, value string) {
	for _, l := range creditLabels {
		if len(s) <= len(l.label) || !strings.EqualFold(s[:len(l.label)], l.label) {
			continue
		}
		rest := strings.TrimLeft(s[len(l.label):], " \t")
		switch {
		case strings.HasPrefix(rest, ":"):
			rest = rest[1:]
		case strings.HasPrefix(rest, "："):
			rest = rest[len("："):]
		case strings.HasSuffix(l.label, " by") && len(rest) < len(s)-len(l.label):
		default:
			continue
		}
		if v := strings.TrimSpace(rest); v != "" {
			return l.field, v
		}
	}
	return "", ""
}

// isKaraoke tells whether a chunk holds the lyrics of a karaoke file as
// text events: it is named "Words" or has @L (language) or @T (title).
func isKaraoke(c *chunk, d *decoder) bool {
	if strings.EqualFold(d.name(c.name), "Words") {
		return true
	}
	for _, t := range c.texts {
		if t.typ == 1 && len(t.data) >= 2 && t.data[0] == '@' && (t.data[1] == 'L' || t.data[1] == 'T') {
			return true
		}
	}
	return false
}

// lyricTolerance is how far a syllable may be from the start of its note.
const lyricTolerance = music.PPQ / 16

// syllable is a syllable of lyrics in PPQ ticks. lead and trail are set
// when a word (or a line) starts before it or ends after it.
type syllable struct {
	tick        int
	text        string
	lead, trail bool
}

// lyrics gives the lyrics of each chunk to the track they are sung on: the
// chunk's own track, else the track with the most notes that start with a
// syllable.
func (p *parser) lyrics(perf *music.Performance, d *decoder, chunks []*chunk, owners []int) {
	if len(perf.Tracks) == 0 {
		return
	}
	budget := 50_000_000 // notes looked up to choose tracks
	for _, c := range chunks {
		syls := p.syllables(c, d)
		if len(syls) == 0 {
			continue
		}
		var cands []*music.Track
		for i, t := range perf.Tracks {
			if owners[i] == c.index {
				cands = append(cands, t)
			}
		}
		if len(cands) == 0 {
			cands = perf.Tracks
		}
		if tonal := slices.DeleteFunc(slices.Clone(cands), func(t *music.Track) bool { return t.Channel == 9 }); len(tonal) > 0 {
			cands = tonal
		}
		best, bestScore := cands[0], 0
		for _, t := range cands {
			if len(cands) == 1 || budget < len(syls) {
				break
			}
			budget -= len(syls)
			score := 0
			for _, s := range syls {
				if _, ok := noteAt(t, s.tick); ok {
					score++
				}
			}
			if score > bestScore {
				best, bestScore = t, score
			}
		}
		for _, s := range syls {
			tick := s.tick
			if at, ok := noteAt(best, tick); ok {
				tick = at
			}
			best.Lyrics = append(best.Lyrics, music.Lyric{Tick: tick, Text: s.text})
		}
	}
	// syllables at one tick sing on one note
	for _, t := range perf.Tracks {
		sort.SliceStable(t.Lyrics, func(i, j int) bool { return t.Lyrics[i].Tick < t.Lyrics[j].Tick })
		out := t.Lyrics[:0]
		for _, l := range t.Lyrics {
			if n := len(out); n > 0 && out[n-1].Tick == l.Tick {
				out[n-1].Text = join(out[n-1].Text, l.Text)
				continue
			}
			out = append(out, l)
		}
		t.Lyrics = slices.Clip(out)
	}
}

// noteAt returns the start of the note of a track nearest to a tick,
// within lyricTolerance.
func noteAt(t *music.Track, tick int) (int, bool) {
	k := sort.Search(len(t.Notes), func(i int) bool { return t.Notes[i].Tick >= tick })
	at, ok := 0, false
	if k < len(t.Notes) && t.Notes[k].Tick-tick <= lyricTolerance {
		at, ok = t.Notes[k].Tick, true
	}
	if k > 0 && tick-t.Notes[k-1].Tick <= lyricTolerance && (!ok || tick-t.Notes[k-1].Tick < at-tick) {
		at, ok = t.Notes[k-1].Tick, true
	}
	return at, ok
}

// join joins two syllables sung on one note.
func join(a, b string) string {
	if s, ok := strings.CutSuffix(a, "-"); ok {
		return s + b
	}
	if spaced(lastRune(a)) && spaced(firstRune(b)) {
		return a + " " + b
	}
	return a + b
}

// syllables returns the lyrics of a chunk: its lyric events, or the text
// events of a karaoke file. The marks of new lines and paragraphs (/ and \)
// are removed; when the words are separated by spaces, a syllable that goes
// on into the next one gets a hyphen, as notation writes it ("Twin-", "kle"),
// if both are letters of a script that separates words with spaces.
func (p *parser) syllables(c *chunk, d *decoder) []syllable {
	var events []rawText
	for _, t := range c.texts {
		if t.typ == 5 {
			events = append(events, t)
		}
	}
	karaoke := len(events) == 0 && isKaraoke(c, d)
	if karaoke {
		for _, t := range c.texts {
			if t.typ == 1 && !bytes.HasPrefix(t.data, []byte("@")) {
				events = append(events, t)
			}
		}
	}
	var out []syllable
	marks := 0 // syllables with spaces or marks around them
	newLine := false
	for _, t := range events {
		s := strings.Map(func(r rune) rune {
			switch {
			case r == '\r' || r == '\n' || r == '\t':
				return ' '
			case unicode.IsControl(r):
				return -1
			}
			return r
		}, d.text(t.data))
		lead := newLine
		for len(s) > 0 && (s[0] == '/' || s[0] == '\\') {
			s, lead = s[1:], true
		}
		trimmed := strings.TrimSpace(s)
		if trimmed != s || lead {
			marks++
		}
		lead = lead || strings.TrimLeftFunc(s, unicode.IsSpace) != s
		trail := strings.TrimRightFunc(s, unicode.IsSpace) != s
		if trimmed == "" {
			// a mark alone ends the syllable before and starts the next
			if n := len(out); n > 0 && (lead || trail) {
				out[n-1].trail = true
			}
			newLine = lead || trail
			continue
		}
		newLine = false
		out = append(out, syllable{tick: p.tm.ppq(t.tick), text: trimmed, lead: lead, trail: trail})
	}
	// karaoke text separates words with spaces; lyric events do when many
	// of them have spaces
	if karaoke || marks*5 >= len(events) {
		for i := 0; i+1 < len(out); i++ {
			a, b := &out[i], out[i+1]
			if !a.trail && !b.lead && !strings.HasSuffix(a.text, "-") && spaced(lastRune(a.text)) && spaced(firstRune(b.text)) {
				a.text += "-"
			}
		}
	}
	return out
}

// spaced tells whether a rune is a letter of a script that separates words
// with spaces.
func spaced(r rune) bool {
	return unicode.IsLetter(r) && (unicode.Is(unicode.Latin, r) || unicode.Is(unicode.Greek, r) || unicode.Is(unicode.Cyrillic, r))
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

func lastRune(s string) rune {
	r, _ := utf8.DecodeLastRuneInString(s)
	return r
}

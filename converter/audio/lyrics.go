package audio

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// lyricLine is a line of the lyrics, as the card sets it.
type lyricLine struct {
	// text is the line. The lines an LRC text leaves without a time stamp
	// follow the line before them, after a line feed.
	text string
	// at is when the line is sung, in milliseconds from the start; -1 for a
	// line of lyrics that are not synchronized.
	at int64
	// stanza reports a line that starts a stanza: a blank line came before.
	stanza bool
}

// maxLyricLines bounds the lines of lyrics the card sets.
const maxLyricLines = 2000

// The time stamps of LRC lyrics: "[mm:ss.xx]" before a line (several when
// the line is sung more than once), "<mm:ss.xx>" before its words (the
// enhanced format), and the tags that describe the file ("[ar:Artist]").
var (
	lrcStamps = regexp.MustCompile(`^(?:\s*\[\d{1,3}:\d{1,2}(?:[.:]\d{1,3})?\])+`)
	lrcStamp  = regexp.MustCompile(`\[(\d{1,3}):(\d{1,2})(?:[.:](\d{1,3}))?\]`)
	lrcWord   = regexp.MustCompile(`<\d{1,3}:\d{1,2}(?:[.:]\d{1,3})?>`)
	lrcTag    = regexp.MustCompile(`^\[([A-Za-z_#]+):(.*)\]$`)
)

// lrcTags are the tags of an LRC file that are not lyrics.
var lrcTags = map[string]bool{
	"ar": true, "ti": true, "al": true, "au": true, "lr": true, "by": true, "offset": true, "length": true,
	"re": true, "ve": true, "la": true, "tool": true, "id": true, "#": true,
}

// lrcTime reads the time of a stamp in milliseconds: minutes, seconds and
// tenths, hundredths or thousandths of a second.
func lrcTime(m []string) int64 {
	minutes, _ := strconv.ParseInt(m[1], 10, 64)
	seconds, _ := strconv.ParseInt(m[2], 10, 64)
	ms := minutes*60000 + seconds*1000
	if f := m[3]; f != "" {
		n, _ := strconv.ParseInt(f, 10, 64)
		for i := len(f); i < 3; i++ {
			n *= 10
		}
		ms += n
	}
	return ms
}

// parseLyrics reads the lyrics of a tag: LRC text, whose lines start with
// the times they are sung at, gives synchronized lines in the order of
// their times; any other text gives its lines as they are. It is LRC when
// lines with a time stamp are at least as many as the other lines.
func parseLyrics(text string) (lines []lyricLine, synced bool) {
	// what the text holds: lines, and the breaks between stanzas (a blank
	// line, a time alone), each with the time it is sorted by: its own, or
	// that of the line before
	type entry struct {
		lyricLine
		key   int64
		pause bool
	}
	var entries []entry
	var offset int64
	timed, plain, count := 0, 0, 0
	key := int64(-1)
	pause := func(at int64) {
		if len(entries) > 0 && !entries[len(entries)-1].pause {
			entries = append(entries, entry{key: at, pause: true})
		}
	}
	for _, raw := range strings.Split(text, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			pause(key)
			continue
		}
		if m := lrcTag.FindStringSubmatch(line); m != nil && lrcTags[strings.ToLower(m[1])] {
			if strings.EqualFold(m[1], "offset") {
				offset, _ = strconv.ParseInt(strings.TrimSpace(m[2]), 10, 64)
			}
			continue
		}
		stamps := lrcStamps.FindString(line)
		body := strings.Join(strings.Fields(lrcWord.ReplaceAllString(line[len(stamps):], " ")), " ")
		if stamps == "" {
			plain++
			if last := len(entries) - 1; last >= 0 && !entries[last].pause && entries[last].at >= 0 {
				entries[last].text += "\n" + body
			} else if count < maxLyricLines {
				entries = append(entries, entry{lyricLine: lyricLine{text: body, at: -1}, key: key})
				count++
			}
			continue
		}
		timed++
		for _, m := range lrcStamp.FindAllStringSubmatch(stamps, -1) {
			key = lrcTime(m)
			switch {
			case body == "": // a time alone ends the line before it: a break in the singing
				entries = append(entries, entry{key: key, pause: true})
			case count < maxLyricLines:
				entries = append(entries, entry{lyricLine: lyricLine{text: body, at: key}, key: key})
				count++
			}
		}
	}
	if timed == 0 || timed < plain {
		return plainLyrics(text), false
	}
	// in the order they are sung: a line with several times is repeated
	slices.SortStableFunc(entries, func(a, b entry) int {
		switch {
		case a.key < b.key:
			return -1
		case a.key > b.key:
			return 1
		}
		return 0
	})
	stanza := false
	for _, e := range entries {
		if e.pause {
			stanza = len(lines) > 0
			continue
		}
		l := e.lyricLine
		l.stanza = stanza
		stanza = false
		if l.at >= 0 {
			// a positive offset brings the lyrics forward
			l.at = max(0, l.at-offset)
		}
		lines = append(lines, l)
	}
	return lines, true
}

// plainLyrics gives the lines of lyrics that are not synchronized: blank
// lines part the stanzas.
func plainLyrics(text string) []lyricLine {
	var out []lyricLine
	blank := false
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			blank = true
			continue
		}
		if len(out) >= maxLyricLines {
			break
		}
		out = append(out, lyricLine{text: l, at: -1, stanza: blank && len(out) > 0})
		blank = false
	}
	return out
}

// lyricsText joins the lines of lyrics into a text, a blank line between
// the stanzas.
func lyricsText(lines []lyricLine) string {
	var b strings.Builder
	for i, l := range lines {
		if i > 0 {
			b.WriteByte('\n')
			if l.stanza {
				b.WriteByte('\n')
			}
		}
		b.WriteString(l.text)
	}
	return b.String()
}

// syltEntry is an entry of a SYLT frame: a piece of the text (a line, a
// word or a syllable) and its time.
type syltEntry struct {
	text string
	at   int64
}

// syltLines makes the lines of the lyrics of the entries of a SYLT frame.
// Entries that hold line feeds break the text into lines where those are,
// a line taking the time of its first piece; without any, entries that
// carry the spaces between words are the words of one line, and other
// entries are a line each.
func syltLines(entries []syltEntry, timed bool) []lyricLine {
	feeds, spaced := false, false
	for _, e := range entries {
		if strings.Contains(e.text, "\n") {
			feeds = true
		}
		if e.text != strings.TrimSpace(e.text) {
			spaced = true
		}
	}
	var out []lyricLine
	var cur strings.Builder
	at := int64(-1)
	blank := false
	flush := func() {
		s := strings.Join(strings.Fields(cur.String()), " ")
		cur.Reset()
		if s == "" {
			blank = true
			return
		}
		if len(out) < maxLyricLines {
			out = append(out, lyricLine{text: s, at: at, stanza: blank && len(out) > 0})
		}
		blank = false
		at = -1
	}
	for _, e := range entries {
		t := int64(-1)
		if timed {
			t = e.at
		}
		if !feeds && !spaced {
			cur.WriteString(e.text)
			at = t
			flush()
			continue
		}
		for i, piece := range strings.Split(e.text, "\n") {
			if i > 0 {
				flush()
			}
			if strings.TrimSpace(piece) != "" && at < 0 && strings.TrimSpace(cur.String()) == "" {
				at = t
			}
			cur.WriteString(piece)
		}
	}
	flush()
	return out
}

// setLyrics picks the lyrics the card shows: the synchronized ones when the
// file has them (a SYLT frame, or LRC text in a lyrics field), else the
// first lyrics field.
func (t *track) setLyrics() {
	if len(t.synced) > 0 && t.synced[0].at >= 0 {
		t.lyrics, t.lyricsSynced = t.synced, true
		return
	}
	texts := t.tags[keyLyrics]
	for _, s := range texts {
		if lines, ok := parseLyrics(s); ok && len(lines) > 0 {
			t.lyrics, t.lyricsSynced = lines, true
			return
		}
	}
	switch {
	case len(texts) > 0:
		t.lyrics = plainLyrics(texts[0])
	default:
		t.lyrics = t.synced // a SYLT frame without times
	}
}

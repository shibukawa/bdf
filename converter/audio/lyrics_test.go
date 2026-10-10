package audio

import (
	"slices"
	"testing"
)

func TestParseLyricsLRC(t *testing.T) {
	lines, synced := parseLyrics("[ti:Sunny Field]\n[ar:Taro Yamada]\n[offset:+250]\n" +
		"[00:01.50]Over the field the sun is high\n" +
		"[00:04.5]Not a cloud in the summer sky\n(softly)\n" +
		"[00:07.250]\n" +
		"[00:20.00][00:08.00]Walking home <00:08.40>along <00:08.90>the lane\n" +
		"[1:02]Hoping it will never rain\n")
	if !synced {
		t.Fatal("LRC text was not taken as synchronized lyrics")
	}
	want := []lyricLine{
		{"Over the field the sun is high", 1250, false},
		// a line without a time follows the line before it
		{"Not a cloud in the summer sky\n(softly)", 4250, false},
		// a time alone is a break before the line that is sung next; the words' own times are dropped
		{"Walking home along the lane", 7750, true},
		// a line with two times is sung twice, each in its place
		{"Walking home along the lane", 19750, false},
		{"Hoping it will never rain", 61750, false},
	}
	if !slices.Equal(lines, want) {
		t.Errorf("lines = %+v\nwant    %+v", lines, want)
	}
	// an offset never takes a line before the start
	if lines, _ := parseLyrics("[offset:5000]\n[00:01.00]a\n[00:09.00]b"); lines[0].at != 0 || lines[1].at != 4000 {
		t.Errorf("offset: %+v", lines)
	}
}

func TestParseLyricsPlain(t *testing.T) {
	// headings in brackets and a stray time are not LRC
	text := "[Chorus: All]\nOver the field\n[0:30] spoken\nNot a cloud\n\n\nWalking home"
	lines, synced := parseLyrics(text)
	want := []lyricLine{
		{"[Chorus: All]", -1, false}, {"Over the field", -1, false}, {"[0:30] spoken", -1, false}, {"Not a cloud", -1, false}, {"Walking home", -1, true},
	}
	if synced || !slices.Equal(lines, want) {
		t.Errorf("synced %v, lines = %+v", synced, lines)
	}
	if got := lyricsText(lines); got != "[Chorus: All]\nOver the field\n[0:30] spoken\nNot a cloud\n\nWalking home" {
		t.Errorf("lyricsText = %q", got)
	}
	if lines, synced := parseLyrics(""); synced || len(lines) != 0 {
		t.Errorf("empty: %v %+v", synced, lines)
	}
	// the lines kept are bounded
	long := ""
	for range maxLyricLines + 10 {
		long += "[00:01.00]la\n"
	}
	if lines, synced := parseLyrics(long); !synced || len(lines) != maxLyricLines {
		t.Errorf("%d lines, synced %v", len(lines), synced)
	}
}

func TestSyltLines(t *testing.T) {
	for _, c := range []struct {
		name    string
		entries []syltEntry
		timed   bool
		want    []lyricLine
	}{
		{"a line an entry", []syltEntry{{"Over the field", 1000}, {"Not a cloud", 4000}}, true,
			[]lyricLine{{"Over the field", 1000, false}, {"Not a cloud", 4000, false}}},
		{"line feeds before the lines, syllables", []syltEntry{{"O", 1000}, {"ver ", 1200}, {"the field", 1500}, {"\nNot ", 4000}, {"a cloud", 4400}, {"\n\nWalk", 8000}, {"ing", 8200}}, true,
			[]lyricLine{{"Over the field", 1000, false}, {"Not a cloud", 4000, false}, {"Walking", 8000, true}}},
		{"line feeds after the lines", []syltEntry{{"Over the field\n", 1000}, {"Not a cloud\n", 4000}}, true,
			[]lyricLine{{"Over the field", 1000, false}, {"Not a cloud", 4000, false}}},
		{"a line feed in an entry of its own", []syltEntry{{"Over", 1000}, {"\n", 3000}, {"Not", 4000}}, true,
			[]lyricLine{{"Over", 1000, false}, {"Not", 4000, false}}},
		{"words of one line", []syltEntry{{"Over ", 1000}, {"the ", 1300}, {"field", 1600}}, true,
			[]lyricLine{{"Over the field", 1000, false}}},
		{"stamps that count frames", []syltEntry{{"Over the field", 38}, {"Not a cloud", 153}}, false,
			[]lyricLine{{"Over the field", -1, false}, {"Not a cloud", -1, false}}},
		{"nothing", nil, true, nil},
	} {
		if got := syltLines(c.entries, c.timed); !slices.Equal(got, c.want) {
			t.Errorf("%s: %+v, want %+v", c.name, got, c.want)
		}
	}
}

// syltFrame makes a SYLT frame: its stamps in milliseconds (format 2) or
// frames (1), its content type, and entries of a text and a time.
func syltFrame(enc, format, content byte, entries ...syltEntry) []byte {
	data := slices.Concat([]byte{enc}, []byte("eng"), []byte{format, content}, encodeText(enc, "desc"))
	end := []byte{0}
	if enc == 1 || enc == 2 {
		end = []byte{0, 0}
	}
	data = append(data, end...)
	for _, e := range entries {
		data = slices.Concat(data, encodeText(enc, e.text), end, u32be(int(e.at)))
	}
	return frameOf(3, "SYLT", 0, data)
}

func TestSYLT(t *testing.T) {
	entries := []syltEntry{{"野原の上に太陽", 1500}, {"\n雲ひとつない夏の空", 4500}}
	want := []lyricLine{{"野原の上に太陽", 1500, false}, {"雲ひとつない夏の空", 4500, false}}
	for _, enc := range []byte{1, 2, 3} {
		tr := readTrack(t, slices.Concat(id3Tag(3, 0, syltFrame(enc, 2, 1, entries...)), mpegFrames(2, 0)))
		tr.setLyrics()
		if !slices.Equal(tr.lyrics, want) || !tr.lyricsSynced {
			t.Errorf("encoding %d: lyrics = %+v, synced %v", enc, tr.lyrics, tr.lyricsSynced)
		}
		// the text is a lyrics field too, for a file without USLT
		wantTag(t, tr, "野原の上に太陽\n雲ひとつない夏の空")
	}
	// synchronized lyrics come before the plain ones of the same file
	uslt := frameOf(3, "USLT", 0, slices.Concat([]byte{0}, []byte("eng"), []byte{0}, []byte("plain words")))
	tr := readTrack(t, slices.Concat(id3Tag(3, 0, uslt, syltFrame(3, 2, 1, entries...)), mpegFrames(2, 0)))
	tr.setLyrics()
	if !slices.Equal(tr.lyrics, want) || !tr.lyricsSynced {
		t.Errorf("with USLT: lyrics = %+v", tr.lyrics)
	}
	wantTag(t, tr, "plain words")
	// stamps in frames give the lines without their times, after USLT
	tr = readTrack(t, slices.Concat(id3Tag(3, 0, syltFrame(3, 1, 1, entries...)), mpegFrames(2, 0)))
	tr.setLyrics()
	if tr.lyricsSynced || len(tr.lyrics) != 2 || tr.lyrics[0].at != -1 || tr.lyrics[1].text != "雲ひとつない夏の空" {
		t.Errorf("frame stamps: %+v, synced %v", tr.lyrics, tr.lyricsSynced)
	}
	// a transcription, a chord chart or the like is not the lyrics
	tr = readTrack(t, slices.Concat(id3Tag(3, 0, syltFrame(3, 2, 2, entries...)), mpegFrames(2, 0)))
	tr.setLyrics()
	if len(tr.lyrics) != 0 {
		t.Errorf("content type 2: %+v", tr.lyrics)
	}
	// LRC text in the lyrics frame
	lrc := frameOf(3, "USLT", 0, slices.Concat([]byte{3}, []byte("eng"), []byte{0}, []byte("[00:01.50]Over the field\r\n[00:04.50]Not a cloud")))
	tr = readTrack(t, slices.Concat(id3Tag(3, 0, lrc), mpegFrames(2, 0)))
	tr.setLyrics()
	if !tr.lyricsSynced || !slices.Equal(tr.lyrics, []lyricLine{{"Over the field", 1500, false}, {"Not a cloud", 4500, false}}) {
		t.Errorf("LRC in USLT: %+v, synced %v", tr.lyrics, tr.lyricsSynced)
	}
}

func wantTag(t *testing.T, tr *track, text string) {
	t.Helper()
	want(t, tr.tags, keyLyrics, text)
}

func TestVorbisSyncedLyrics(t *testing.T) {
	// LYRICS holds plain text and SYNCEDLYRICS the same lines with times: the synchronized ones are shown
	comment := vorbisComment("TITLE=Sunny Field", "LYRICS=Over the field\nNot a cloud", "SYNCEDLYRICS=[00:01.00]Over the field\n[00:04.00]Not a cloud")
	data := slices.Concat([]byte("fLaC"), flacBlock(0, false, streamInfoBlock(44100, 2, 16, 441000)), flacBlock(4, true, comment))
	tr := readTrack(t, data)
	tr.setLyrics()
	if !tr.lyricsSynced || !slices.Equal(tr.lyrics, []lyricLine{{"Over the field", 1000, false}, {"Not a cloud", 4000, false}}) {
		t.Errorf("lyrics = %+v, synced %v", tr.lyrics, tr.lyricsSynced)
	}
}

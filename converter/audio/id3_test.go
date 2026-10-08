package audio

import (
	"bytes"
	"slices"
	"strings"
	"testing"
)

func TestID3v22(t *testing.T) {
	pic := slices.Concat([]byte{0}, []byte("PNG"), []byte{3}, []byte("cover\x00"), tinyPNG)
	tag := id3Tag(2, 0,
		textFrame(2, "TT2", 0, "Old Title"), textFrame(2, "TP1", 0, "Old Artist"), textFrame(2, "TCO", 0, "(17)"),
		textFrame(2, "TYE", 0, "1999"), textFrame(2, "TRK", 0, "7/9"),
		frameOf(2, "PIC", 0, pic), frameOf(2, "XYZ", 0, []byte("unknown")))
	tr := readTrack(t, slices.Concat(tag, mpegFrames(3, 0)))
	want(t, tr.tags, keyTitle, "Old Title")
	want(t, tr.tags, keyArtist, "Old Artist")
	want(t, tr.tags, keyGenre, "Rock")
	want(t, tr.tags, keyDate, "1999")
	want(t, tr.tags, keyTrack, "7")
	want(t, tr.tags, keyTrackTotal, "9")
	if len(tr.pictures) != 1 || tr.pictures[0].mime != "image/png" || tr.pictures[0].typ != 3 || !bytes.Equal(tr.pictures[0].data, tinyPNG) {
		t.Errorf("pictures = %+v", tr.pictures)
	}
}

func TestID3v23(t *testing.T) {
	// a frame with FF 00 in it, unsynchronised at the tag level
	title := textFrame(3, "TIT2", 1, "Title ÿ")
	// a compressed frame: the decompressed size, then zlib data
	comp := []byte("\x00A compressed album")
	album := frameOf(3, "TALB", 0x0080, slices.Concat(u32be(len(comp)), deflate(comp)))
	encrypted := frameOf(3, "TPUB", 0x0040, []byte("\x01secret"))
	grouped := frameOf(3, "TCOM", 0x0020, []byte("\x07\x00Composer"))
	comm := func(lang, desc, text string) []byte {
		return frameOf(3, "COMM", 0, slices.Concat([]byte{0}, []byte(lang), []byte(desc), []byte{0}, []byte(text)))
	}
	uslt := frameOf(3, "USLT", 0, slices.Concat([]byte{1}, []byte("eng"), encodeText(1, "words"), []byte{0, 0}, encodeText(1, "Line one\r\nLine two")))
	txxx := func(desc, value string) []byte {
		return frameOf(3, "TXXX", 0, slices.Concat([]byte{0}, []byte(desc), []byte{0}, []byte(value)))
	}
	body := slices.Concat(title, album, encrypted, grouped,
		textFrame(3, "TPE1", 0, "Ann; Bob"), textFrame(3, "TCON", 0, "(17)(RX)Rock"),
		textFrame(3, "TYER", 0, "2026"), textFrame(3, "TDAT", 0, "0109"), textFrame(3, "TIME", 0, "1020"),
		comm("eng", "iTunNORM", " 00000012"), comm("eng", "note", "A described comment"), comm("eng", "", "The comment"),
		uslt, txxx("ALBUM ARTIST", "The Band"), txxx("ISRC", "XX0000000000"), txxx("TITLE", "ignored: the frame comes first"))
	body = unsyncBytes(body)
	tag := slices.Concat([]byte("ID3\x03\x00\x80"), syncsafeBytes(len(body)), body)
	tr := readTrack(t, slices.Concat(tag, mpegFrames(2, 0)))
	want(t, tr.tags, keyTitle, "Title ÿ")
	want(t, tr.tags, keyAlbum, "A compressed album")
	want(t, tr.tags, keyPublisher)
	want(t, tr.tags, keyComposer, "Composer")
	want(t, tr.tags, keyArtist, "Ann", "Bob")
	want(t, tr.tags, keyGenre, "Rock", "Remix")
	want(t, tr.tags, keyDate, "2026-09-01T10:20")
	want(t, tr.tags, keyComment, "The comment")
	want(t, tr.tags, keyLyrics, "Line one\nLine two")
	want(t, tr.tags, keyAlbumArtist, "The Band")
	want(t, tr.tags, keyISRC, "XX0000000000")
	if d := tr.dublinCore(); d.Date.First() != "2026-09-01T10:20" || d.Description.First() != "The comment" {
		t.Errorf("dc = %+v", d)
	}
}

// unsyncBytes applies unsynchronisation: FF followed by 00 or E0-FF gets
// a 00.
func unsyncBytes(b []byte) []byte {
	var out []byte
	for i, c := range b {
		out = append(out, c)
		if c == 0xFF && i+1 < len(b) && (b[i+1] == 0 || b[i+1] >= 0xE0) {
			out = append(out, 0)
		}
	}
	return out
}

func TestID3v24(t *testing.T) {
	// a frame unsynchronised by itself (the byte order mark and ÿ hold
	// FF), with a data length indicator
	raw := slices.Concat([]byte{1}, encodeText(1, "Français ÿ"))
	title := frameOf(4, "TIT2", 0x0003, slices.Concat(syncsafeBytes(len(raw)), unsyncBytes(raw)))
	// a frame whose size is written plain, not syncsafe
	long := []byte("\x03" + strings.Repeat("x", 200))
	plain := slices.Concat([]byte("TPUB"), u32be(len(long)), []byte{0, 0}, long)
	tag := id3Tag(4, 0x10,
		title, textFrame(4, "TPE1", 3, "山田太郎", "佐藤花子"), textFrame(4, "TALB", 2, "Album"),
		textFrame(4, "TDRC", 3, "2026-09-01T10:20:30"), textFrame(4, "TDOR", 3, "1999"), plain,
		textFrame(4, "TRCK", 3, "3"), textFrame(4, "TLAN", 3, "jpn"))
	data := slices.Concat(tag, mpegFrames(2, 0))
	if n := id3v2Size(data); n != int64(len(tag)) {
		t.Errorf("id3v2Size = %d, want %d (with the footer)", n, len(tag))
	}
	tr := readTrack(t, data)
	want(t, tr.tags, keyTitle, "Français ÿ")
	want(t, tr.tags, keyArtist, "山田太郎", "佐藤花子")
	want(t, tr.tags, keyAlbum, "Album")
	want(t, tr.tags, keyDate, "2026-09-01T10:20:30")
	want(t, tr.tags, keyOriginalDate, "1999")
	want(t, tr.tags, keyPublisher, strings.Repeat("x", 200))
	want(t, tr.tags, keyTrack, "3")
	if dc := tr.dublinCore(); dc.Language.First() != "ja" || dc.Creator[1] != "佐藤花子" {
		t.Errorf("dc = %+v", dc)
	}
}

func TestLegacyEncodings(t *testing.T) {
	sjis := []byte("\x90\xb0\x82\xea\x82\xbd\x96\xec\x8c\xb4")                                                                        // 晴れた野原
	tag := id3Tag(3, 0, frameOf(3, "TIT2", 0, slices.Concat([]byte{0}, sjis)), frameOf(3, "TPE1", 0, []byte("\x00\x8e\x52\x93\x63"))) // 山田
	tr := readTrack(t, slices.Concat(tag, mpegFrames(2, 0)))
	want(t, tr.tags, keyTitle, "晴れた野原")
	want(t, tr.tags, keyArtist, "山田")
	// one frame that is not Shift_JIS makes the file Windows-1252
	tag = id3Tag(3, 0, frameOf(3, "TIT2", 0, slices.Concat([]byte{0}, sjis)), frameOf(3, "TPE1", 0, []byte("\x00Caf\xe9")))
	tr = readTrack(t, slices.Concat(tag, mpegFrames(2, 0)))
	want(t, tr.tags, keyArtist, "Café")
	// UTF-8 in a "Latin-1" frame is read as UTF-8
	tag = id3Tag(3, 0, frameOf(3, "TIT2", 0, []byte("\x00晴れた野原")))
	tr = readTrack(t, slices.Concat(tag, mpegFrames(2, 0)))
	want(t, tr.tags, keyTitle, "晴れた野原")
	// ID3v1 in Shift_JIS
	tr = readTrack(t, slices.Concat(mpegFrames(2, 0), id3v1(string(sjis), "", "", "2001", "", 0, 17)))
	want(t, tr.tags, keyTitle, "晴れた野原")
	want(t, tr.tags, keyGenre, "Rock")
}

func TestID3v1Fills(t *testing.T) {
	tag := id3Tag(3, 0, textFrame(3, "TIT2", 0, "v2 title"))
	data := slices.Concat(tag, mpegFrames(4, 0), id3v1("v1 title", "v1 artist", "v1 album", "1988", "a comment", 5, 13))
	tr := readTrack(t, data)
	want(t, tr.tags, keyTitle, "v2 title")
	want(t, tr.tags, keyArtist, "v1 artist")
	want(t, tr.tags, keyAlbum, "v1 album")
	want(t, tr.tags, keyDate, "1988")
	want(t, tr.tags, keyTrack, "5")
	want(t, tr.tags, keyComment, "a comment")
	want(t, tr.tags, keyGenre, "Pop")
	// the frames are scanned without the tag: 4 frames of 1152 samples
	if tr.frames() != 4 {
		t.Errorf("frames = %d", tr.frames())
	}
}

func TestChapters(t *testing.T) {
	chap := func(id string, start int, title string) []byte {
		return frameOf(4, "CHAP", 0, slices.Concat([]byte(id), []byte{0}, u32be(start), u32be(start+1000), u32be(-1), u32be(-1), textFrame(4, "TIT2", 3, title)))
	}
	tag := id3Tag(4, 0, chap("c2", 65000, "Second"), chap("c1", 0, "First"), chap("c3", 3600_000, ""))
	tr := readTrack(t, slices.Concat(tag, mpegFrames(2, 0)))
	if len(tr.chapters) != 3 || tr.chapters[0].title != "Second" || tr.chapters[0].start != 65 || tr.chapters[2].title != "c3" || tr.chapters[2].start != 3600 {
		t.Errorf("chapters = %+v", tr.chapters)
	}
}

func TestGenres(t *testing.T) {
	for in, out := range map[string][]string{
		"(17)": {"Rock"}, "17": {"Rock"}, "(17)Pop": {"Rock", "Pop"}, "((Parens": {"(Parens"}, "(RX)(CR)": {"Remix", "Cover"},
		"Folk": {"Folk"}, "(999)": {"999"}, "": nil,
	} {
		if got := genreNames(in); !slices.Equal(got, out) {
			t.Errorf("genreNames(%q) = %q, want %q", in, got, out)
		}
	}
}

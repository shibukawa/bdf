package audio

import (
	"math"
	"slices"
	"testing"
)

func TestFLAC(t *testing.T) {
	data := slices.Concat([]byte("fLaC"), flacBlock(0, false, streamInfoBlock(44100, 2, 16, 44100*10)),
		flacBlock(2, false, make([]byte, 40)), // an application block
		flacBlock(6, false, pictureBlock(4, "image/png", "back", tinyPNG)),
		flacBlock(6, false, pictureBlock(3, "image/jpeg", "front", []byte("\xff\xd8\xff"))),
		flacBlock(4, false, vorbisComment("title=Lower Case", "artist=Ann", "LYRICS=l1\nl2", "UNSYNCEDLYRICS=ignored: lyrics came first", "ORGANIZATION=Org")),
		flacBlock(1, true, make([]byte, 100)), // padding
		make([]byte, 1000))
	tr := readTrack(t, data)
	if tr.format != "flac" || tr.codec != "FLAC" || tr.sampleRate != 44100 || tr.channels != 2 || tr.bits != 16 || math.Abs(tr.duration-10) > 1e-9 || tr.bitrate != 800 {
		t.Errorf("track = %+v", tr)
	}
	want(t, tr.tags, keyTitle, "Lower Case")
	want(t, tr.tags, keyArtist, "Ann")
	want(t, tr.tags, keyLyrics, "l1\nl2", "ignored: lyrics came first")
	want(t, tr.tags, keyPublisher, "Org")
	if c := tr.cover(); c == nil || c.typ != 3 || c.mime != "image/jpeg" {
		t.Errorf("cover = %+v", c)
	}
	// an ID3v2 tag before the stream fills what the comment lacks
	tagged := slices.Concat(id3Tag(3, 0, textFrame(3, "TIT2", 0, "ID3 Title"), textFrame(3, "TALB", 0, "ID3 Album")), data)
	tr = readTrack(t, tagged)
	want(t, tr.tags, keyTitle, "Lower Case")
	want(t, tr.tags, keyAlbum, "ID3 Album")
}

func TestCover(t *testing.T) {
	tr := &track{pictures: []picture{{typ: 4}, {typ: 0}, {typ: 3}, {typ: 0}}}
	if c := tr.cover(); c.typ != 3 {
		t.Errorf("cover = %+v", c)
	}
	tr.pictures = tr.pictures[:2]
	if c := tr.cover(); c != &tr.pictures[1] {
		t.Errorf("cover = %+v, want the picture of no kind", c)
	}
	tr.pictures = tr.pictures[:1]
	if c := tr.cover(); c != &tr.pictures[0] {
		t.Errorf("cover = %+v, want the first", c)
	}
	if (&track{}).cover() != nil {
		t.Error("a cover without pictures")
	}
}

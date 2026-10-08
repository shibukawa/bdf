package audio

import (
	"bytes"
	"encoding/base64"
	"math"
	"slices"
	"strings"
	"testing"
)

func vorbisIdent(rate, channels, nominal int) []byte {
	return slices.Concat([]byte("\x01vorbis"), u32le(0), []byte{byte(channels)}, u32le(rate), u32le(0), u32le(nominal), u32le(0), []byte{0xB8, 1})
}

func TestOggVorbis(t *testing.T) {
	pic := base64.StdEncoding.EncodeToString(pictureBlock(3, "image/png", "front", tinyPNG))
	// a comment longer than a page, with a long field
	comment := slices.Concat([]byte("\x03vorbis"), vorbisComment("TITLE=Ogg Title", "ARTIST=Ann", "ARTIST=Bob", "ALBUM ARTIST=The Band", "DATE=2026",
		"TRACKNUMBER=2/9", "METADATA_BLOCK_PICTURE="+pic, "DESCRIPTION=desc", "COMMENT=comment", "CHAPTER001=00:00:00.000", "CHAPTER001NAME=One",
		"CHAPTER002=00:01:30.500", "CHAPTER010=01:00:00.000", "CHAPTER010NAME=Ten", "LICENSE=CC0", "junk="+strings.Repeat("j", 80000)), []byte{1})
	data := oggPages(7, 44100*3, vorbisIdent(44100, 2, 160000), comment, []byte("\x05vorbis setup"), []byte("audio"))
	tr := readTrack(t, data)
	if tr.format != "ogg" || tr.codec != "Vorbis" || tr.sampleRate != 44100 || tr.channels != 2 || tr.bitrate != 160000 || math.Abs(tr.duration-3) > 1e-9 {
		t.Errorf("track = %+v", tr)
	}
	want(t, tr.tags, keyTitle, "Ogg Title")
	want(t, tr.tags, keyArtist, "Ann", "Bob")
	want(t, tr.tags, keyAlbumArtist, "The Band")
	want(t, tr.tags, keyTrack, "2")
	want(t, tr.tags, keyTrackTotal, "9")
	want(t, tr.tags, keyComment, "desc", "comment")
	want(t, tr.tags, keyCopyright, "CC0")
	if len(tr.pictures) != 1 || tr.pictures[0].typ != 3 || !bytes.Equal(tr.pictures[0].data, tinyPNG) {
		t.Errorf("pictures = %+v", tr.pictures)
	}
	if len(tr.chapters) != 3 || tr.chapters[0].title != "One" || tr.chapters[1].title != "Chapter 2" || tr.chapters[1].start != 90.5 || tr.chapters[2].start != 3600 {
		t.Errorf("chapters = %+v", tr.chapters)
	}
}

func TestOggOpusAndFLAC(t *testing.T) {
	head := slices.Concat([]byte("OpusHead"), []byte{1, 2}, u16le(312), u32le(48000), u16le(0), []byte{0})
	tags := slices.Concat([]byte("OpusTags"), vorbisComment("TITLE=Opus Title"))
	tr := readTrack(t, oggPages(9, 48000*2+312, head, tags, []byte("audio")))
	if tr.codec != "Opus" || tr.channels != 2 || tr.sampleRate != 48000 || math.Abs(tr.duration-2) > 1e-9 {
		t.Errorf("opus: track = %+v", tr)
	}
	want(t, tr.tags, keyTitle, "Opus Title")

	ident := slices.Concat([]byte("\x7fFLAC\x01\x00"), u16be(2), []byte("fLaC"), flacBlock(0, false, streamInfoBlock(96000, 1, 24, 96000*4)))
	comment := flacBlock(4, false, vorbisComment("TITLE=FLAC in Ogg"))
	picture := flacBlock(6, true, pictureBlock(0, "image/png", "", tinyPNG))
	tr = readTrack(t, oggPages(3, 96000*4, ident, comment, picture, []byte("audio")))
	if tr.codec != "FLAC" || tr.sampleRate != 96000 || tr.bits != 24 || math.Abs(tr.duration-4) > 1e-9 || len(tr.pictures) != 1 {
		t.Errorf("flac: track = %+v", tr)
	}
	want(t, tr.tags, keyTitle, "FLAC in Ogg")
}

func TestOggStreams(t *testing.T) {
	// a Theora stream first: the audio stream is still read, with a warning
	theora := oggPages(1, 0, slices.Concat([]byte("\x80theora"), make([]byte, 40)))
	vorbis := oggPages(2, 44100, vorbisIdent(44100, 1, 0), slices.Concat([]byte("\x03vorbis"), vorbisComment("TITLE=With video"), []byte{1}), []byte("x"))
	tr := readTrack(t, slices.Concat(theora, vorbis))
	want(t, tr.tags, keyTitle, "With video")
	if len(tr.warnings) != 1 || tr.duration != 1 {
		t.Errorf("warnings %q duration %v", tr.warnings, tr.duration)
	}
	// no audio stream
	if _, err := read(bytes.NewReader(theora), int64(len(theora))); err == nil {
		t.Error("a Theora file was read")
	}
}

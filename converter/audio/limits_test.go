package audio

import (
	"bytes"
	"io"
	"slices"
	"strings"
	"testing"
)

// Files that claim more than they hold, or more than the converter keeps:
// what is read is bounded by the file and by the limits, without a panic.

// sparse is a file of a size whose bytes beyond a prefix are zero, for
// files too large to make.
type sparse struct {
	prefix []byte
	size   int64
}

func (s *sparse) ReadAt(b []byte, off int64) (int, error) {
	if off >= s.size {
		return 0, io.EOF
	}
	n := 0
	for ; n < len(b) && off+int64(n) < s.size; n++ {
		if i := off + int64(n); i < int64(len(s.prefix)) {
			b[n] = s.prefix[i]
		} else {
			b[n] = 0
		}
	}
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

func TestID3Claims(t *testing.T) {
	// a tag that claims the largest size, on a short file
	tag := slices.Concat([]byte("ID3\x04\x00\x00\x7f\x7f\x7f\x7f"), textFrame(4, "TIT2", 3, "Short"))
	tr, err := read(bytes.NewReader(tag), int64(len(tag)))
	if err != nil || tr.tags.first(keyTitle) != "Short" || len(tr.warnings) == 0 {
		t.Errorf("track %+v, err %v", tr, err)
	}
	// a frame that claims more than the tag: a truncated last frame
	frame := slices.Concat([]byte("TIT2"), u32be(1<<20), []byte{0, 0}, []byte("\x00Cut"))
	tag = id3Tag(3, 0, frame)
	tr = readTrack(t, slices.Concat(tag, mpegFrames(1, 0)))
	want(t, tr.tags, keyTitle, "Cut")
	// a compressed frame that inflates to too much
	big := frameOf(3, "TIT2", 0x0080, slices.Concat(u32be(maxInflate*2), deflate(bytes.Repeat([]byte{'a'}, maxInflate+100))))
	tr = readTrack(t, slices.Concat(id3Tag(3, 0, big), mpegFrames(1, 0)))
	want(t, tr.tags, keyTitle)
	// a hundred pictures: at most maxValues are kept
	var pics [][]byte
	for range 100 {
		pics = append(pics, frameOf(3, "APIC", 0, slices.Concat([]byte("\x00image/png\x00\x03\x00"), tinyPNG)))
	}
	tr = readTrack(t, slices.Concat(id3Tag(3, 0, pics...), mpegFrames(1, 0)))
	if len(tr.pictures) != maxValues {
		t.Errorf("%d pictures kept", len(tr.pictures))
	}
	// a hundred values of one field: at most maxValues
	var vals []string
	for i := range 100 {
		vals = append(vals, strings.Repeat("v", i+1))
	}
	tr = readTrack(t, slices.Concat(id3Tag(4, 0, textFrame(4, "TPE1", 3, vals...)), mpegFrames(1, 0)))
	if len(tr.tags[keyArtist]) != maxValues {
		t.Errorf("%d artists kept", len(tr.tags[keyArtist]))
	}
	// lyrics are cut to maxLyrics
	lyrics := frameOf(4, "USLT", 0, slices.Concat([]byte("\x03eng\x00"), bytes.Repeat([]byte("line\n"), maxLyrics/5+100)))
	tr = readTrack(t, slices.Concat(id3Tag(4, 0, lyrics), mpegFrames(1, 0)))
	if n := len(tr.tags.first(keyLyrics)); n > maxLyrics {
		t.Errorf("lyrics of %d bytes kept", n)
	}
}

func TestFLACClaims(t *testing.T) {
	// a block longer than the file
	data := slices.Concat([]byte("fLaC"), []byte{0x84, 0x10, 0x00, 0x00}, make([]byte, 100))
	if _, err := read(bytes.NewReader(data), int64(len(data))); err == nil {
		t.Error("a block longer than the file was read")
	}
	// a comment block larger than the limit, in a file that large: skipped
	hdr := slices.Concat([]byte("fLaC"), flacBlock(0, false, streamInfoBlock(44100, 1, 16, 44100)), []byte{0x84, 0xFF, 0xFF, 0xFF})
	f := &sparse{prefix: hdr, size: int64(len(hdr)) + 0xFFFFFF + 10}
	tr, err := read(f, f.size)
	if err != nil || len(tr.tags) != 0 || len(tr.warnings) != 1 {
		t.Errorf("track %+v, err %v", tr, err)
	}
}

func TestMP4Claims(t *testing.T) {
	// a movie box larger than the limit: not read
	hdr := slices.Concat(box("ftyp", []byte("M4A "), u32be(0)), u32be(maxMoov+100), []byte("moov"))
	f := &sparse{prefix: hdr, size: int64(len(hdr)) + maxMoov + 100}
	if _, err := read(f, f.size); err == nil {
		t.Error("a huge movie box was read")
	}
	if Detect(hdr, f, f.size) {
		t.Error("detected without a readable movie box")
	}
	// a box that claims more than the file: the walk stops
	data := slices.Concat(box("ftyp", []byte("M4A "), u32be(0)), u32be(1<<30), []byte("moov"), make([]byte, 100))
	if _, err := read(bytes.NewReader(data), int64(len(data))); err == nil {
		t.Error("a truncated movie box was read")
	}
}

func TestOggClaims(t *testing.T) {
	// thousands of pages of a stream that is not audio, then nothing
	var data []byte
	for range maxOggPages + 10 {
		data = append(data, oggPages(5, 0, []byte("\x80theora"))...)
	}
	if _, err := read(bytes.NewReader(data), int64(len(data))); err == nil {
		t.Error("read a file without audio")
	}
	// the comment header on a page cut short: the stream is read without
	// its comments
	data = slices.Concat(oggPages(1, 0, vorbisIdent(44100, 1, 0)), []byte("OggS\x00\x00"), make([]byte, 8), u32le(1), make([]byte, 8), []byte{255}, bytes.Repeat([]byte{255}, 100))
	tr := readTrack(t, data)
	if tr.codec != "Vorbis" || tr.tags.has(keyTitle) || tr.duration != 0 {
		t.Errorf("track = %+v", tr)
	}
}

func TestRIFFClaims(t *testing.T) {
	// chunks that claim more than the file
	data := slices.Concat([]byte("RIFF"), u32le(1<<30), []byte("WAVE"), riffChunk("fmt ", u16le(1), u16le(1), u32le(8000), u32le(16000), u16le(2), u16le(16)),
		[]byte("LIST"), u32le(1<<30), []byte("INFO"), []byte("INAM"), u32le(1<<30), []byte("Cut title"))
	tr := readTrack(t, data)
	if tr.sampleRate != 8000 {
		t.Errorf("track = %+v", tr)
	}
	data = slices.Concat([]byte("FORM"), u32be(1<<30), []byte("AIFF"), []byte("COMM"), u32be(1<<30), []byte{0, 1})
	if _, err := read(bytes.NewReader(data), int64(len(data))); err == nil {
		t.Error("a truncated common chunk was read")
	}
}

func TestMPEGJunk(t *testing.T) {
	// junk past the limit before the first frame: not MPEG audio
	data := slices.Concat(bytes.Repeat([]byte{'j'}, maxJunk+100), mpegFrames(2, 0))
	if _, err := read(bytes.NewReader(data), int64(len(data))); err == nil {
		t.Error("read audio after too much junk")
	}
	// junk between frames: skipped
	data = slices.Concat(mpegFrames(2, 0), bytes.Repeat([]byte{'j'}, 1000), mpegFrames(2, 0))
	tr := readTrack(t, data)
	if tr.frames() != 4 {
		t.Errorf("%d frames", tr.frames())
	}
}

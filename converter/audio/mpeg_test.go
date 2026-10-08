package audio

import (
	"bytes"
	"math"
	"slices"
	"testing"
)

// frames is the number of frames counted, for tests.
func (t *track) frames() int { return t.scanned }

func TestMPEGHeader(t *testing.T) {
	for _, c := range []struct {
		b      []byte
		ok     bool
		length int
		rate   int
	}{
		{[]byte{0xFF, 0xFB, 0x90, 0x00}, true, 417, 44100}, // MPEG-1 III 128k 44.1k
		{[]byte{0xFF, 0xFB, 0x92, 0x00}, true, 418, 44100}, // padded
		{[]byte{0xFF, 0xF3, 0x70, 0xC0}, true, 182, 22050}, // MPEG-2 III 56k 22.05k mono
		{[]byte{0xFF, 0xE3, 0x40, 0xC0}, true, 208, 11025}, // MPEG-2.5 III 32k 11.025k
		{[]byte{0xFF, 0xFE, 0x90, 0x00}, true, 312, 44100}, // layer I 288k
		{[]byte{0xFF, 0xFD, 0x90, 0x00}, true, 522, 44100}, // layer II 160k
		{[]byte{0xFF, 0xFB, 0xF0, 0x00}, false, 0, 0},      // bad bitrate index
		{[]byte{0xFF, 0xFB, 0x00, 0x00}, false, 0, 0},      // free format
		{[]byte{0xFF, 0xFB, 0x9C, 0x00}, false, 0, 0},      // reserved sample rate
		{[]byte{0xFF, 0xEB, 0x90, 0x00}, false, 0, 0},      // reserved version
		{[]byte{0xFF, 0xF9, 0x90, 0x00}, false, 0, 0},      // reserved layer
		{[]byte{0xFF, 0xFB, 0x90, 0x02}, false, 0, 0},      // reserved emphasis
		{[]byte{0x00, 0xFB, 0x90, 0x00}, false, 0, 0},
	} {
		f, ok := parseMPEGHeader(c.b)
		if ok != c.ok || ok && (f.length != c.length || f.sampleRate != c.rate) {
			t.Errorf("%x: %+v %v, want ok %v length %d rate %d", c.b, f, ok, c.ok, c.length, c.rate)
		}
	}
}

func TestScanMPEG(t *testing.T) {
	// junk, then frames, then an ID3v1 tag and an APE tag
	junk := bytes.Repeat([]byte{0xFF, 0xE0, 0x12}, 50)
	ape := slices.Concat(bytes.Repeat([]byte{1}, 40), []byte("APETAGEX"), u32le(2000), u32le(72), u32le(0), u32le(0), bytes.Repeat([]byte{0}, 8))
	data := slices.Concat(junk, mpegFrames(10, 0), ape, id3v1("t", "", "", "", "", 0, 255))
	tr := readTrack(t, data)
	if tr.format != "mp3" || tr.codec != "MP3" || tr.frames() != 10 || math.Abs(tr.duration-10*1152.0/44100) > 1e-9 || tr.bitrate != 128000 || tr.vbr || tr.channels != 2 {
		t.Errorf("track = %+v", tr)
	}
	// a Xing header counts the frames of the whole file; the bitrate is
	// the average over the bytes it counts
	tr = readTrack(t, mpegFrames(3, 1000))
	seconds := 1000 * 1152.0 / 44100
	if math.Abs(tr.duration-seconds) > 1e-9 || tr.bitrate != int(1000*417*8/seconds+0.5) || !tr.vbr {
		t.Errorf("with Xing: duration %v bitrate %d vbr %v", tr.duration, tr.bitrate, tr.vbr)
	}
	// an Info header: the same, but constant
	info := mpegFrames(3, 1000)
	copy(info[36:], "Info")
	if tr = readTrack(t, info); tr.vbr || tr.scanned != 1000 {
		t.Errorf("with Info: vbr %v frames %d", tr.vbr, tr.scanned)
	}
	// a VBRI header
	frame := mpegFrames(1, 0)
	copy(frame[36:], "VBRI")
	copy(frame[46:], u32be(417*500))
	copy(frame[50:], u32be(500))
	tr = readTrack(t, slices.Concat(frame, mpegFrames(2, 0)))
	if math.Abs(tr.duration-500*1152.0/44100) > 1e-9 {
		t.Errorf("with VBRI: duration %v", tr.duration)
	}
	// a single frame, cut short, is still MPEG audio
	if _, err := read(bytes.NewReader(mpegFrames(1, 0)[:300]), 300); err != nil {
		t.Errorf("a short frame: %v", err)
	}
	// text is not
	if _, err := read(bytes.NewReader([]byte("hello, world")), 12); err == nil {
		t.Error("text read as audio")
	}
}

func TestScanADTS(t *testing.T) {
	tr := readTrack(t, adtsFrames(5, 300))
	if tr.format != "aac" || tr.codec != "AAC LC" || tr.sampleRate != 22050 || tr.channels != 1 || math.Abs(tr.duration-5*1024.0/22050) > 1e-9 {
		t.Errorf("track = %+v", tr)
	}
	seconds := 5 * 1024.0 / 22050
	if want := int(5*300*8/seconds + 0.5); tr.bitrate != want {
		t.Errorf("bitrate = %d", tr.bitrate)
	}
	// with an ID3v2 tag in front
	tr = readTrack(t, slices.Concat(id3Tag(3, 0, textFrame(3, "TIT2", 0, "Raw AAC")), adtsFrames(2, 100)))
	if tr.format != "aac" || tr.tags.first(keyTitle) != "Raw AAC" {
		t.Errorf("track = %+v", tr)
	}
}

func TestDetect(t *testing.T) {
	for name, c := range map[string]struct {
		data []byte
		want bool
	}{
		"mp3":          {mpegFrames(2, 0), true},
		"id3":          {slices.Concat(id3Tag(3, 0), mpegFrames(1, 0)), true},
		"adts":         {adtsFrames(2, 50), true},
		"one sync":     {slices.Concat([]byte{0xFF, 0xFB, 0x90, 0x00}, bytes.Repeat([]byte{'x'}, 500)), false},
		"text":         {[]byte("ID2 is not a tag"), false},
		"flac":         {[]byte("fLaC\x80\x00\x00\x22"), true},
		"wav":          {[]byte("RIFF\x00\x00\x00\x00WAVEfmt "), true},
		"aiff":         {[]byte("FORM\x00\x00\x00\x00AIFFCOMM"), true},
		"ogg vorbis":   {oggPages(1, 0, slices.Concat([]byte("\x01vorbis"), make([]byte, 23))), true},
		"ogg theora":   {oggPages(1, 0, slices.Concat([]byte("\x80theora"), make([]byte, 40))), false},
		"avif":         {box("ftyp", []byte("avif"), u32be(0), []byte("avifmif1")), false},
		"mp4 no moov":  {box("ftyp", []byte("M4A "), u32be(0), []byte("M4A isom")), false},
		"ape tag only": {[]byte("APETAGEX\xd0\x07\x00\x00"), false},
	} {
		head := c.data[:min(len(c.data), 1024)]
		if got := Detect(head, bytes.NewReader(c.data), int64(len(c.data))); got != c.want {
			t.Errorf("%s: Detect = %v, want %v", name, got, c.want)
		}
	}
}

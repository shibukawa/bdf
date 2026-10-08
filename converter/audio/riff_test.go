package audio

import (
	"math"
	"slices"
	"testing"
)

func TestWAV(t *testing.T) {
	// RF64 with an extensible format chunk (PCM), the data size in ds64,
	// an INFO list, a broadcast extension and an ID3 chunk
	fmtChunk := riffChunk("fmt ", u16le(0xFFFE), u16le(2), u32le(48000), u32le(48000*4), u16le(4), u16le(16), u16le(22), u16le(16), u32le(3),
		u16le(1), []byte("\x00\x00\x00\x00\x10\x00\x80\x00\x00\xaa\x00\x38\x9b\x71"))
	ds64 := riffChunk("ds64", u64le(1000), u64le(48000*4*3), u64le(48000*3))
	info := riffChunk("LIST", []byte("INFO"), riffChunk("INAM", []byte("Wave Title\x00")), riffChunk("IART", []byte("Ann; Bob\x00")),
		riffChunk("ITRK", []byte("3/12\x00")), riffChunk("ICRD", []byte("2026-09-01\x00")), riffChunk("IMUS", []byte("Composer\x00")), riffChunk("ILNG", []byte("eng\x00")))
	bext := make([]byte, 602)
	copy(bext, "A description")
	copy(bext[256:], "Originator")
	copy(bext[320:], "2026:09:02")
	copy(bext[330:], "10:20:30")
	id3 := riffChunk("id3 ", id3Tag(3, 0, textFrame(3, "TALB", 0, "ID3 Album"), textFrame(3, "TIT2", 0, "ID3 Title")))
	data := slices.Concat(riffChunk("RF64", []byte("WAVE"), ds64, fmtChunk, info, riffChunk("bext", bext), id3, []byte("data\xff\xff\xff\xff"), make([]byte, 100)))
	tr := readTrack(t, data)
	if tr.format != "wav" || tr.codec != "PCM" || tr.sampleRate != 48000 || tr.channels != 2 || tr.bits != 16 || math.Abs(tr.duration-3) > 1e-9 || tr.bitrate != 48000*4*8 {
		t.Errorf("track = %+v", tr)
	}
	want(t, tr.tags, keyTitle, "ID3 Title") // the ID3 chunk comes first
	want(t, tr.tags, keyAlbum, "ID3 Album")
	want(t, tr.tags, keyArtist, "Ann", "Bob")
	want(t, tr.tags, keyTrack, "3")
	want(t, tr.tags, keyTrackTotal, "12")
	want(t, tr.tags, keyDate, "2026-09-01")
	want(t, tr.tags, keyComposer, "Composer")
	want(t, tr.tags, keyLanguage, "eng")
	want(t, tr.tags, keyComment, "A description")
	// without INFO, the broadcast extension's date
	data = slices.Concat(riffChunk("RIFF", []byte("WAVE"), riffChunk("fmt ", u16le(7), u16le(1), u32le(8000), u32le(8000), u16le(1), u16le(8)), riffChunk("bext", bext), riffChunk("data", make([]byte, 16000))))
	tr = readTrack(t, data)
	if tr.codec != "μ-law" || tr.bits != 0 || math.Abs(tr.duration-2) > 1e-9 {
		t.Errorf("track = %+v", tr)
	}
	want(t, tr.tags, keyDate, "2026-09-02T10:20:30")
}

func TestAIFF(t *testing.T) {
	rate := []byte{0x40, 0x0E, 0xAC, 0x44, 0, 0, 0, 0, 0, 0} // 44100
	comm := iffChunk("COMM", u16be(2), u32be(44100*2), u16be(24), rate, []byte("sowt"), []byte{0})
	data := slices.Concat(iffChunk("FORM", []byte("AIFC"), comm, iffChunk("NAME", []byte("AIFF Title")), iffChunk("AUTH", []byte("Ann")),
		iffChunk("(c) ", []byte("2026 Ann")), iffChunk("ANNO", []byte("a note")), iffChunk("SSND", make([]byte, 100)),
		iffChunk("ID3 ", id3Tag(4, 0, textFrame(4, "TIT2", 3, "ID3 Title")))))
	tr := readTrack(t, data)
	if tr.format != "aiff" || tr.codec != "PCM" || tr.sampleRate != 44100 || tr.channels != 2 || tr.bits != 24 || math.Abs(tr.duration-2) > 1e-9 || tr.bitrate != 44100*2*24 {
		t.Errorf("track = %+v", tr)
	}
	want(t, tr.tags, keyTitle, "ID3 Title")
	want(t, tr.tags, keyArtist, "Ann")
	want(t, tr.tags, keyCopyright, "2026 Ann")
	want(t, tr.tags, keyComment, "a note")
	// μ-law AIFF-C
	comm = iffChunk("COMM", u16be(1), u32be(8000), u16be(8), []byte{0x40, 0x0B, 0xFA, 0, 0, 0, 0, 0, 0, 0}, []byte("ulaw"), []byte{0})
	tr = readTrack(t, slices.Concat(iffChunk("FORM", []byte("AIFC"), comm)))
	if tr.codec != "μ-law" || tr.sampleRate != 8000 || tr.duration != 1 {
		t.Errorf("track = %+v", tr)
	}
}

func TestExtended(t *testing.T) {
	for b, want := range map[[10]byte]float64{
		{0x40, 0x0E, 0xAC, 0x44}: 44100, {0x40, 0x0E, 0xBB, 0x80}: 48000, {0x40, 0x0B, 0xFA}: 8000, {0x40, 0x0D, 0xAC, 0x44}: 22050, {}: 0,
	} {
		if got := extended(b[:]); math.Abs(got-want) > 1e-6 {
			t.Errorf("extended(%x) = %v, want %v", b, got, want)
		}
	}
}

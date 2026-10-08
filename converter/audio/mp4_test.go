package audio

import (
	"bytes"
	"math"
	"slices"
	"testing"
)

// m4a makes an M4A file: the movie box (before or after the media data)
// with a sound track of AAC at 44.1 kHz, stereo, 2.5 seconds, and the
// given user data.
func m4a(faststart bool, udta []byte, extraTraks ...[]byte) []byte {
	esds := fullBox("esds", 0,
		[]byte{0x03, 0x19, 0x00, 0x01, 0x00},                                 // ES_Descriptor
		[]byte{0x04, 0x11, 0x40, 0x15, 0, 0, 0}, u32be(128000), u32be(96000), // DecoderConfigDescriptor: AAC, average 96 kbit/s
		[]byte{0x05, 0x02, 0x12, 0x10}) // AudioSpecificConfig: LC, 44.1 kHz, stereo
	entry := slices.Concat(make([]byte, 6), u16be(1), u16be(0), u16be(0), u32be(0), u16be(2), u16be(16), u16be(0), u16be(0), u32be(44100<<16), esds)
	stsd := fullBox("stsd", 0, u32be(1), box("mp4a", entry))
	mdhd := fullBox("mdhd", 0, u32be(0), u32be(0), u32be(44100), u32be(44100*5/2), u16be(0x15C7), u16be(0)) // language "eng"
	trak := box("trak", box("mdia", mdhd, fullBox("hdlr", 0, u32be(0), []byte("soun"), make([]byte, 12), []byte("Sound\x00")),
		box("minf", box("stbl", stsd))))
	mvhd := fullBox("mvhd", 0, u32be(0), u32be(0), u32be(1000), u32be(2500), make([]byte, 80))
	moov := box("moov", slices.Concat(mvhd, trak, slices.Concat(extraTraks...), udta))
	mdat := slices.Concat(u32be(1), []byte("mdat"), be.AppendUint64(nil, 16+30000), make([]byte, 30000)) // a 64-bit size
	ftyp := box("ftyp", []byte("M4A "), u32be(0), []byte("M4A mp42isom"))
	if faststart {
		return slices.Concat(ftyp, moov, mdat)
	}
	return slices.Concat(ftyp, mdat, moov)
}

func TestMP4(t *testing.T) {
	ilst := box("ilst",
		item("\xa9nam", 1, []byte("Title")), item("\xa9ART", 1, []byte("Artist")), item("aART", 1, []byte("Album Artist")),
		item("\xa9alb", 2, encodeText(2, "Album")), item("\xa9day", 1, []byte("2026-09-01")), item("gnre", 0, u16be(18)),
		item("trkn", 0, []byte{0, 0, 0, 3, 0, 12, 0, 0}), item("disk", 0, []byte{0, 0, 0, 1, 0, 2}),
		item("\xa9lyr", 1, []byte("la la\nla")), item("tmpo", 21, u16be(120)), item("cprt", 1, []byte("© me")),
		box("----", box("mean", u32be(0), []byte("com.apple.iTunes")), box("name", u32be(0), []byte("ISRC")), box("data", u32be(1), u32be(0), []byte("XX0000000001"))),
		box("----", box("mean", u32be(0), []byte("com.apple.iTunes")), box("name", u32be(0), []byte("LABEL")), box("data", u32be(1), u32be(0), []byte("A Label"))),
		item("covr", 14, tinyPNG), item("covr", 13, []byte("\xff\xd8\xff")))
	meta := fullBox("meta", 0, fullBox("hdlr", 0, u32be(0), []byte("mdirappl"), make([]byte, 10)), ilst)
	chpl := fullBox("chpl", 0, []byte{2}, be.AppendUint64(nil, 0), []byte{5}, []byte("Intro"), be.AppendUint64(nil, 15e7), []byte{6}, []byte("Middle"))
	udta := box("udta", meta, chpl)
	for _, faststart := range []bool{true, false} {
		data := m4a(faststart, udta)
		if kind := mp4Kind(data[:64], bytes.NewReader(data), int64(len(data))); kind != "audio" {
			t.Fatalf("mp4Kind = %q", kind)
		}
		tr := readTrack(t, data)
		if tr.format != "m4a" || tr.codec != "AAC" || tr.sampleRate != 44100 || tr.channels != 2 || math.Abs(tr.duration-2.5) > 1e-9 || tr.bitrate != 96000 {
			t.Errorf("faststart %v: track = %+v", faststart, tr)
		}
		want(t, tr.tags, keyTitle, "Title")
		want(t, tr.tags, keyArtist, "Artist")
		want(t, tr.tags, keyAlbumArtist, "Album Artist")
		want(t, tr.tags, keyAlbum, "Album")
		want(t, tr.tags, keyDate, "2026-09-01")
		want(t, tr.tags, keyGenre, "Rock")
		want(t, tr.tags, keyTrack, "3")
		want(t, tr.tags, keyTrackTotal, "12")
		want(t, tr.tags, keyDisc, "1")
		want(t, tr.tags, keyDiscTotal, "2")
		want(t, tr.tags, keyLyrics, "la la\nla")
		want(t, tr.tags, keyBPM, "120")
		want(t, tr.tags, keyCopyright, "© me")
		want(t, tr.tags, keyISRC, "XX0000000001")
		want(t, tr.tags, keyPublisher, "A Label")
		want(t, tr.tags, keyLanguage, "eng")
		if len(tr.pictures) != 2 || tr.pictures[0].mime != "image/png" || !bytes.Equal(tr.cover().data, tinyPNG) {
			t.Errorf("pictures = %+v", tr.pictures)
		}
		if len(tr.chapters) != 2 || tr.chapters[1].title != "Middle" || tr.chapters[1].start != 15 {
			t.Errorf("chapters = %+v", tr.chapters)
		}
	}
}

func TestMP4Video(t *testing.T) {
	video := box("trak", box("mdia", fullBox("hdlr", 0, u32be(0), []byte("vide"), make([]byte, 12))))
	data := m4a(true, nil, video)
	if kind := mp4Kind(data[:64], bytes.NewReader(data), int64(len(data))); kind != "video" {
		t.Errorf("mp4Kind = %q", kind)
	}
	if _, err := read(bytes.NewReader(data), int64(len(data))); err == nil {
		t.Error("a video file was read as audio")
	}
	if Detect(data[:64], bytes.NewReader(data), int64(len(data))) {
		t.Error("a video file is detected as audio")
	}
}

func TestMP4QuickTimeKeys(t *testing.T) {
	// QuickTime's metadata: keys named in a keys box, a meta box without
	// version and flags, and text atoms under udta
	keys := fullBox("keys", 0, u32be(2), box("mdta", []byte("com.apple.quicktime.title")), box("mdta", []byte("com.apple.quicktime.artist")))
	ilst := box("ilst", item("\x00\x00\x00\x01", 1, []byte("QT Title")), item("\x00\x00\x00\x02", 1, []byte("QT Artist")))
	meta := box("meta", fullBox("hdlr", 0, u32be(0), []byte("mdta"), make([]byte, 10)), keys, ilst)
	qtText := box("\xa9alb", u16be(8), u16be(0), []byte("QT Album"))
	tr := readTrack(t, m4a(true, box("udta", meta, qtText)))
	want(t, tr.tags, keyTitle, "QT Title")
	want(t, tr.tags, keyArtist, "QT Artist")
	want(t, tr.tags, keyAlbum, "QT Album")
}

func TestMP4Codecs(t *testing.T) {
	for _, c := range []struct{ entry, codec string }{{"alac", "Apple Lossless"}, {"fLaC", "FLAC"}, {"Opus", "Opus"}, {"ac-3", "AC-3"}, {"drms", "AAC"}, {"twos", "PCM"}, {"xyzw", "xyzw"}} {
		entry := slices.Concat(make([]byte, 6), u16be(1), u16be(0), u16be(0), u32be(0), u16be(1), u16be(16), u16be(0), u16be(0), u32be(48000<<16))
		stsd := fullBox("stsd", 0, u32be(1), box(c.entry, entry))
		trak := box("trak", box("mdia", fullBox("mdhd", 0, u32be(0), u32be(0), u32be(48000), u32be(48000), u16be(0), u16be(0)),
			fullBox("hdlr", 0, u32be(0), []byte("soun"), make([]byte, 12)), box("minf", box("stbl", stsd))))
		data := slices.Concat(box("ftyp", []byte("isom"), u32be(0)), box("moov", trak), box("mdat", make([]byte, 100)))
		tr := readTrack(t, data)
		if tr.codec != c.codec || tr.sampleRate != 48000 || tr.channels != 1 || tr.duration != 1 || tr.protected != (c.entry == "drms") {
			t.Errorf("%s: track = %+v", c.entry, tr)
		}
	}
}

func TestDescriptor(t *testing.T) {
	// a length in two bytes (0x81 0x02 = 130)
	long := slices.Concat([]byte{0x04, 0x81, 0x02}, bytes.Repeat([]byte{7}, 130))
	b := slices.Concat([]byte{0x03, 0x01, 0x00}, long)
	if d := descriptor(b, 0x04); len(d) != 130 || d[0] != 7 {
		t.Errorf("descriptor = %v", d)
	}
	if d := descriptor([]byte{0x04, 0x7F, 0x01}, 0x04); len(d) != 1 { // claims more than it holds
		t.Errorf("short descriptor = %v", d)
	}
}

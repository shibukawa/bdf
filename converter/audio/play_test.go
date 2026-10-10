package audio

import (
	"bytes"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

// cuesOf decodes the cues of the card's view; nil when it has none.
func cuesOf(t *testing.T, doc *bdf.Document) *bdf.Cues {
	t.Helper()
	play := doc.Views[0].Play
	if play == nil || play.Cues == "" {
		return nil
	}
	h, err := bdf.ParseHash(play.Cues)
	if err != nil {
		t.Fatal(err)
	}
	part := doc.Part(h)
	if part == nil || part.Type != bdf.PartIndex {
		t.Fatalf("cues part %v", part)
	}
	cues, err := bdf.DecodeCues(part.Data)
	if err != nil {
		t.Fatal(err)
	}
	return cues
}

// audioOf returns the audio part of the card's view; nil when it has none.
func audioOf(t *testing.T, doc *bdf.Document) *bdf.Part {
	t.Helper()
	play := doc.Views[0].Play
	if play == nil {
		return nil
	}
	if play.Seq != "" || play.Audio == "" {
		t.Fatalf("play = %+v", play)
	}
	h, err := bdf.ParseHash(play.Audio)
	if err != nil {
		t.Fatal(err)
	}
	return doc.Part(h)
}

func TestPlay(t *testing.T) {
	// every format is stored as the file is, under its media type
	for file, mime := range map[string]string{
		"tagged.mp3": "audio/mpeg", "tagged.m4a": "audio/mp4", "plain.aac": "audio/aac", "tagged.flac": "audio/flac",
		"tagged.ogg": "audio/ogg", "tagged.opus": "audio/ogg", "tagged.wav": "audio/wav", "tagged.aiff": "audio/aiff",
	} {
		res := convert(t, file, nil)
		raw, err := os.ReadFile(filepath.Join("testdata", file))
		if err != nil {
			t.Fatal(err)
		}
		part := audioOf(t, res.Doc)
		if part == nil || part.Type != bdf.PartAudio || !bytes.Equal(part.Data, raw) || res.Audio != int64(len(raw)) {
			t.Errorf("%s: audio part %v, Audio %d of %d", file, part != nil, res.Audio, len(raw))
			continue
		}
		if got := res.Doc.Views[0].Play.Type; got != mime {
			t.Errorf("%s: type %q, want %q", file, got, mime)
		}
		// the audio is not compressed again, and the document reads back
		var buf bytes.Buffer
		if err := res.Doc.WriteSingle(&buf); err != nil {
			t.Fatal(err)
		}
		r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range r.Manifest.Parts {
			if e.T == bdf.PartAudio && (e.Enc != bdf.EncIdentity || e.Len != len(raw)) {
				t.Errorf("%s: audio stored as %s, %d bytes", file, e.Enc, e.Len)
			}
		}
		back, err := r.ToDocument()
		if err != nil {
			t.Fatal(err)
		}
		if p := audioOf(t, back); p == nil || !bytes.Equal(p.Data, raw) {
			t.Errorf("%s: the audio does not read back", file)
		}
	}
}

func TestPlayOptions(t *testing.T) {
	// NoPlay: the card alone
	res := convert(t, "tagged.mp3", &Options{NoPlay: true})
	if res.Doc.Views[0].Play != nil || res.Audio != 0 || len(res.Warnings) != 0 {
		t.Errorf("NoPlay: play %+v, Audio %d, warnings %q", res.Doc.Views[0].Play, res.Audio, res.Warnings)
	}
	for _, p := range res.Doc.Parts() {
		if p.Type == bdf.PartAudio {
			t.Error("NoPlay: an audio part")
		}
	}
	if s := res.Summary(); !strings.HasSuffix(s, "mono; 4 lines of lyrics, 2 chapters)") {
		t.Errorf("NoPlay: summary %q", s)
	}
	// a file larger than MaxAudio converts, without its audio, and says so
	res = convert(t, "tagged.mp3", &Options{MaxAudio: 1024})
	if res.Doc.Views[0].Play != nil || len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "larger than 1 KiB") {
		t.Errorf("MaxAudio: play %+v, warnings %q", res.Doc.Views[0].Play, res.Warnings)
	}
	if text := pageText(t, res.Doc); !strings.Contains(text, "Over the field the sun is high") {
		t.Errorf("MaxAudio: the card lost its text: %q", text)
	}
	if res = convert(t, "tagged.mp3", &Options{MaxAudio: 1 << 20}); res.Audio == 0 {
		t.Error("MaxAudio above the size: no audio")
	}
	// protected audio is not stored: no player plays it
	tr := newTrack("m4a")
	tr.protected = true
	var warned []string
	if b := storedAudio(bytes.NewReader([]byte("data")), 4, tr, &Options{}, func(m string) { warned = append(warned, m) }); b != nil || len(warned) != 0 {
		t.Errorf("protected: %q, warnings %q", b, warned)
	}
	// a file that ends before its size
	if b := storedAudio(bytes.NewReader([]byte("data")), 8, newTrack("mp3"), &Options{}, func(m string) { warned = append(warned, m) }); b != nil || len(warned) != 1 {
		t.Errorf("short read: %q, warnings %q", b, warned)
	}
	for n, want := range map[int64]string{256 << 20: "256 MiB", 1536 << 10: "1.5 MiB", 2048: "2 KiB", 100: "100 bytes"} {
		if got := byteSize(n); got != want {
			t.Errorf("byteSize(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestPlayParams(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "tagged.mp3"))
	if err != nil {
		t.Fatal(err)
	}
	conv := func(o *converter.Options) (*converter.Result, error) {
		return converter.Convert(bytes.NewReader(raw), int64(len(raw)), "audio", o)
	}
	plays := func(o *converter.Options) bool {
		t.Helper()
		res, err := conv(o)
		if err != nil {
			t.Fatal(err)
		}
		return res.Doc.Views[0].Play != nil
	}
	if !plays(&converter.Options{}) {
		t.Error("the default leaves the audio out")
	}
	// the page selected for a thumbnail: no audio, unless it is asked for
	if plays(&converter.Options{Pages: converter.PageList(1)}) {
		t.Error("a selection of pages stores the audio")
	}
	if !plays(&converter.Options{Pages: converter.PageList(1), Params: map[string]string{"play": "true"}}) {
		t.Error("play=true with a selection of pages leaves the audio out")
	}
	if plays(&converter.Options{Params: map[string]string{"play": "false"}}) {
		t.Error("play=false stores the audio")
	}
	if plays(&converter.Options{Params: map[string]string{"maxaudio": "0.001"}}) {
		t.Error("maxaudio below the size stores the audio")
	}
	if !plays(&converter.Options{Params: map[string]string{"maxaudio": "1"}}) {
		t.Error("maxaudio above the size leaves the audio out")
	}
	for _, p := range []map[string]string{{"play": "maybe"}, {"maxaudio": "0"}, {"maxaudio": "-1"}, {"maxaudio": "big"}, {"maxaudio": "NaN"}} {
		if _, err := conv(&converter.Options{Params: p}); err == nil {
			t.Errorf("params %v: no error", p)
		}
	}
}

func TestCues(t *testing.T) {
	// synchronized lyrics: a cue a line, on the rows the line is set on
	lrc := "[00:01.50]Over the field the sun is high\n[00:04.50]Not a cloud in the summer sky\n\n" +
		"[00:08.00]" + strings.Repeat("Walking home along the lane ", 6) + "\n[00:12.00]Hoping it will never rain"
	uslt := frameOf(3, "USLT", 0, slices.Concat([]byte{3}, []byte("eng"), []byte{0}, []byte(lrc)))
	chap := func(id string, start int, title string) []byte {
		return frameOf(3, "CHAP", 0, slices.Concat([]byte(id), []byte{0}, u32be(start), u32be(start+1000), u32be(-1), u32be(-1), textFrame(3, "TIT2", 0, title)))
	}
	chapters := [][]byte{chap("c1", 0, "Verse"), chap("c2", 8000, "Chorus")}
	data := slices.Concat(id3Tag(3, 0, slices.Concat(uslt, chapters[0], chapters[1])), mpegFrames(4, 0))
	res, err := Convert(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Synced || res.Lyrics != 4 || res.Chapters != 2 {
		t.Errorf("result: synced %v, %d lines, %d chapters", res.Synced, res.Lyrics, res.Chapters)
	}
	if s := res.Summary(); !strings.HasSuffix(s, "; plays, 4 lines of synchronized lyrics, 2 chapters)") {
		t.Errorf("summary %q", s)
	}
	// the time stamps are not on the card
	if text := pageText(t, res.Doc); strings.Contains(text, "[00:") || !strings.Contains(text, "Over the field the sun is high") {
		t.Errorf("card text %q", text)
	}
	cues := cuesOf(t, res.Doc)
	if cues == nil {
		t.Fatal("no cues")
	}
	page := res.Doc.Views[0].Pages[0]
	var ticks []uint32
	for i, c := range cues.Cues {
		ticks = append(ticks, c.Tick)
		if int(c.System) != i || c.X != margin {
			t.Errorf("cue %d: %+v", i, c)
		}
	}
	if !slices.Equal(ticks, []uint32{1500, 4500, 8000, 12000}) {
		t.Errorf("ticks = %v", ticks)
	}
	if len(cues.Systems) != 4 {
		t.Fatalf("%d systems", len(cues.Systems))
	}
	for i, s := range cues.Systems {
		if s.Page != 0 || s.X != margin || s.W != page.W-2*margin || s.Y <= 0 || s.Y+s.H > page.H {
			t.Errorf("system %d: %+v on a page of %v × %v", i, s, page.W, page.H)
		}
		if i > 0 && s.Y < cues.Systems[i-1].Y+cues.Systems[i-1].H {
			t.Errorf("system %d overlaps the one before: %+v, %+v", i, cues.Systems[i-1], s)
		}
	}
	// a line is 16 units high, a wrapped one as many times as it has rows; a stanza starts further down
	if cues.Systems[0].H != 16 || cues.Systems[1].Y != cues.Systems[0].Y+16 || cues.Systems[2].Y != cues.Systems[1].Y+16+8 || cues.Systems[2].H < 32 {
		t.Errorf("systems = %+v", cues.Systems)
	}

	// without synchronized lyrics the chapters have the cues
	plain := frameOf(3, "USLT", 0, slices.Concat([]byte{3}, []byte("eng"), []byte{0}, []byte("Over the field\nNot a cloud")))
	data = slices.Concat(id3Tag(3, 0, slices.Concat(plain, chapters[0], chapters[1])), mpegFrames(4, 0))
	if res, err = Convert(bytes.NewReader(data), int64(len(data)), nil); err != nil {
		t.Fatal(err)
	}
	cues = cuesOf(t, res.Doc)
	if cues == nil || len(cues.Cues) != 2 || cues.Cues[0].Tick != 0 || cues.Cues[1].Tick != 8000 || cues.Systems[0].H != 17 || cues.Systems[1].Y != cues.Systems[0].Y+17 {
		t.Errorf("chapter cues = %+v", cues)
	}

	// neither: the audio plays without cues, and cues are never stored without the audio
	data = slices.Concat(id3Tag(3, 0, plain), mpegFrames(4, 0))
	if res, err = Convert(bytes.NewReader(data), int64(len(data)), nil); err != nil {
		t.Fatal(err)
	}
	if play := res.Doc.Views[0].Play; play == nil || play.Audio == "" || play.Cues != "" {
		t.Errorf("plain lyrics: play %+v", play)
	}
	data = slices.Concat(id3Tag(3, 0, uslt), mpegFrames(4, 0))
	if res, err = Convert(bytes.NewReader(data), int64(len(data)), &Options{NoPlay: true}); err != nil {
		t.Fatal(err)
	}
	for _, p := range res.Doc.Parts() {
		if p.Type == bdf.PartIndex {
			if _, err := bdf.DecodeCues(p.Data); err == nil {
				t.Error("NoPlay: cues without the audio")
			}
		}
	}
}

// The files of test/audio/gen.sh that have the times of their lyrics.
func TestSyncedFiles(t *testing.T) {
	// an ID3 SYLT frame beside the plain lyrics of USLT: the song the viewer's sample plays
	res := convert(t, "synced.mp3", nil)
	if res.Format != "mp3" || res.Cover != "jpeg" || math.Abs(res.Duration-15) > 0.2 || !res.Synced || res.Lyrics != 4 || len(res.Warnings) != 0 {
		t.Errorf("synced.mp3: %s, cover %q, %v s, synced %v, %d lines, warnings %q", res.Format, res.Cover, res.Duration, res.Synced, res.Lyrics, res.Warnings)
	}
	cues := cuesOf(t, res.Doc)
	if cues == nil || len(cues.Systems) != 4 {
		t.Fatalf("synced.mp3: cues %+v", cues)
	}
	for i, at := range []uint32{1500, 4500, 7500, 10500} {
		if c := cues.Cues[i]; c.Tick != at || int(c.System) != i {
			t.Errorf("synced.mp3: cue %d = %+v, want tick %d", i, c, at)
		}
	}
	// the second stanza starts further down
	if s := cues.Systems; s[1].Y != s[0].Y+16 || s[2].Y != s[1].Y+16+8 {
		t.Errorf("synced.mp3: systems %+v", cues.Systems)
	}
	if text := pageText(t, res.Doc); !strings.Contains(text, "Lyrics\nOver the field the sun is high Not a cloud in the summer sky\nWalking home along the lane Hoping it will never rain") {
		t.Errorf("synced.mp3: text %q", text)
	}

	// LRC text in a Vorbis comment: its tags and stamps are not shown, and the line with two times is sung twice
	res = convert(t, "lrc.flac", nil)
	if res.Format != "flac" || math.Abs(res.Duration-4) > 0.1 || !res.Synced || res.Lyrics != 3 || len(res.Warnings) != 0 {
		t.Errorf("lrc.flac: %s, %v s, synced %v, %d lines, warnings %q", res.Format, res.Duration, res.Synced, res.Lyrics, res.Warnings)
	}
	cues = cuesOf(t, res.Doc)
	if cues == nil || len(cues.Cues) != 3 || cues.Cues[0].Tick != 500 || cues.Cues[1].Tick != 1500 || cues.Cues[2].Tick != 3000 {
		t.Fatalf("lrc.flac: cues %+v", cues)
	}
	if text := pageText(t, res.Doc); !strings.Contains(text, "Lyrics\n野原の上に太陽 雲ひとつない夏の空\n雲ひとつない夏の空") || strings.Contains(text, "[") {
		t.Errorf("lrc.flac: text %q", text)
	}
}

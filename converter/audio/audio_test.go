package audio

import (
	"encoding/json"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/thumbnail"
)

// convert converts a file of the test data.
func convert(t *testing.T, name string, opts *Options) *Result {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// dcJSON writes Dublin Core as the manifest does.
func dcJSON(t *testing.T, dc bdf.DublinCore) string {
	t.Helper()
	b, err := json.Marshal(dc)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// pageText returns the text of the first page, as a search index gets it.
func pageText(t *testing.T, doc *bdf.Document) string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Views) != 1 || len(st.Views[0].Pages) != 1 {
		t.Fatalf("text of %d views", len(st.Views))
	}
	return st.Views[0].Pages[0].Text
}

func TestConvert(t *testing.T) {
	for _, c := range []struct {
		file, format, codec string
		rate, channels      int
		cover               string
		dc                  string
		text                []string
	}{
		{"tagged.mp3", "mp3", "MP3", 22050, 1, "jpeg",
			`{"title":"Sunny Field","creator":"Taro Yamada","subject":"Folk","description":"A song about a field in the sun","publisher":"Example Records",` +
				`"contributor":["Various Artists","Hanako Sato"],"date":"2026-09-01","type":"Sound","format":"audio/mpeg","identifier":"JPA012600001",` +
				`"language":"en","relation":"Four Seasons","rights":"© 2026 Example Records"}`,
			[]string{"Sunny Field\nTaro Yamada\nFour Seasons\n", "Track\n3 of 12\n", "Hoping it will never rain", "Chapters\n0:00 Verse\n0:01 Chorus"}},
		// ID3v2.4 in UTF-8, two artists in one frame; the composer is one
		// of them, so not a contributor
		{"japanese.mp3", "mp3", "MP3", 22050, 1, "",
			`{"title":"晴れた野原","creator":["山田太郎","佐藤花子"],"subject":"フォーク","date":"2026","type":"Sound","format":"audio/mpeg","language":"ja","relation":"四季"}`,
			[]string{"晴れた野原\n山田太郎, 佐藤花子\n四季\n", "Lyrics\n野原の上に太陽 雲ひとつない夏の空"}},
		// the movie box after the media data; the language is undetermined
		{"tagged.m4a", "m4a", "AAC", 22050, 1, "png",
			`{"title":"Sunny Field","creator":"Taro Yamada","subject":"Folk","description":"A song about a field in the sun",` +
				`"contributor":["Various Artists","Hanako Sato"],"date":"2026-09-01","type":"Sound","format":"audio/mp4","relation":"Four Seasons","rights":"© 2026 Example Records"}`,
			[]string{"Grouping\nSummer\n", "Not a cloud in the summer sky"}},
		// Vorbis comments in lower case; the lyricist is the composer
		{"tagged.flac", "flac", "FLAC", 22050, 1, "png",
			`{"title":"Sunny Field","creator":"Taro Yamada","subject":"Folk","description":"A song about a field in the sun","contributor":"Hanako Sato",` +
				`"date":"2026-09-01","type":"Sound","format":"audio/flac","identifier":"JPA012600001","relation":"Four Seasons","rights":"© 2026 Example Records"}`,
			[]string{"Format\nFLAC, 16-bit, "}},
		{"tagged.ogg", "ogg", "Vorbis", 22050, 1, "png",
			`{"title":"Sunny Field","creator":"Taro Yamada","date":"2026","type":"Sound","format":"audio/ogg","relation":"Four Seasons"}`,
			[]string{"Track\n3\n", "Format\nVorbis, "}},
		{"tagged.opus", "ogg", "Opus", 48000, 1, "png",
			`{"title":"Sunny Field","creator":"Taro Yamada","date":"2026","type":"Sound","format":"audio/ogg","relation":"Four Seasons"}`,
			[]string{"Format\nOpus, "}},
		// RIFF INFO, and the broadcast extension's description
		{"tagged.wav", "wav", "PCM", 8000, 1, "",
			`{"title":"Sunny Field","creator":"Taro Yamada","subject":"Folk","description":"A song about a field in the sun","date":"2026-09-01",` +
				`"type":"Sound","format":"audio/wav","relation":"Four Seasons","rights":"© 2026 Example Records"}`,
			[]string{"Format\nPCM, 16-bit, 128 kbps, 8 kHz, mono"}},
		// an ID3 chunk and a NAME chunk
		{"tagged.aiff", "aiff", "PCM", 8000, 1, "jpeg",
			`{"title":"Sunny Field","creator":"Taro Yamada","date":"2026","type":"Sound","format":"audio/aiff","relation":"Four Seasons"}`,
			nil},
		{"plain.aac", "aac", "AAC LC", 22050, 1, "", `{"type":"Sound","format":"audio/aac"}`, []string{"plain\n", "Format\nAAC LC, "}},
	} {
		t.Run(c.file, func(t *testing.T) {
			res := convert(t, c.file, nil)
			if res.Format != c.format || res.Codec != c.codec || res.SampleRate != c.rate || res.Channels != c.channels || res.Cover != c.cover {
				t.Errorf("got %s %s %d Hz %d ch cover %q, want %s %s %d Hz %d ch cover %q", res.Format, res.Codec, res.SampleRate, res.Channels, res.Cover, c.format, c.codec, c.rate, c.channels, c.cover)
			}
			if math.Abs(res.Duration-1) > 0.15 || res.Bitrate <= 0 {
				t.Errorf("duration %v, bitrate %d", res.Duration, res.Bitrate)
			}
			if c.cover != "" && (res.CoverWidth != 320 || res.CoverHeight != 320) {
				t.Errorf("cover %d × %d", res.CoverWidth, res.CoverHeight)
			}
			if len(res.Warnings) != 0 {
				t.Errorf("warnings %q", res.Warnings)
			}
			if got := dcJSON(t, res.Doc.Meta.DC); got != c.dc {
				t.Errorf("dc\n got %s\nwant %s", got, c.dc)
			}
			if res.Doc.Meta.Source != c.format {
				t.Errorf("source %q", res.Doc.Meta.Source)
			}
			v := res.Doc.Views[0]
			if len(res.Doc.Views) != 1 || v.Kind != bdf.ViewFixed || len(v.Pages) != 1 || v.TextIndex == "" {
				t.Fatalf("views %+v", res.Doc.Views)
			}
			p := v.Pages[0]
			wantW := float32(placeholderSide)
			if c.cover != "" {
				wantW = coverMin // a 320 px cover (240 pt) is scaled up to the narrowest page
			}
			if p.W != wantW || p.H < p.W {
				t.Errorf("page %v × %v, want %v wide and at least as tall", p.W, p.H, wantW)
			}
			text := pageText(t, res.Doc)
			for _, s := range c.text {
				if !strings.Contains(text, s) {
					t.Errorf("text lacks %q:\n%s", s, text)
				}
			}
		})
	}
}

func TestOptions(t *testing.T) {
	// a title given, no text index, and the cover re-encoded
	res := convert(t, "tagged.mp3", &Options{Title: "Given", NoTextIndex: true, Images: imgconv.Options{Mode: imgconv.Convert}})
	if res.Doc.Meta.DC.Title.First() != "Given" || res.Doc.Views[0].Title != "Given" || res.Doc.Views[0].TextIndex != "" {
		t.Errorf("doc %+v", res.Doc.Views[0])
	}
	if imgconv.Available() && res.Cover != "webp" {
		t.Errorf("cover %q, want webp", res.Cover)
	}
	// a file without a title is named after the file
	res = convert(t, "plain.aac", nil)
	if res.Doc.Views[0].Title != "plain" || res.Doc.Meta.DC.Title.First() != "" {
		t.Errorf("title %q, dc %+v", res.Doc.Views[0].Title, res.Doc.Meta.DC)
	}
	if s := res.Summary(); !strings.HasPrefix(s, "1 page (no cover; ADTS AAC LC, 0:01, ") || !strings.HasSuffix(s, " kbps, 22.05 kHz, mono; plays)") {
		t.Errorf("summary %q", s)
	}
	res = convert(t, "tagged.mp3", nil)
	if s := res.Summary(); !strings.HasPrefix(s, "1 page (cover 320 × 320 px, JPEG; MP3, 0:01, ") || !strings.HasSuffix(s, "; plays, 4 lines of lyrics, 2 chapters)") {
		t.Errorf("summary %q", s)
	}
	// warnings go to Warn when it is set
	var warned []string
	data := slicesConcat(id3Tag(3, 0, frameOf(3, "APIC", 0, []byte("\x00image/tiff\x00\x03\x00II*\x00"))), mpegFrames(2, 0))
	r, err := Convert(bytesReader(data), int64(len(data)), &Options{Warn: func(msg string) { warned = append(warned, msg) }})
	if err != nil || len(warned) != 1 || len(r.Warnings) != 0 || r.Cover != "" || !strings.Contains(warned[0], "image/tiff") {
		t.Errorf("warnings %q %q, err %v", warned, r.Warnings, err)
	}
}

func TestThumbnail(t *testing.T) {
	// the thumbnail of a file with a cover is the cover: the top square of
	// the card (the fixture's cover is a blue sky over a green field, with
	// a yellow sun in the upper right)
	res := convert(t, "tagged.mp3", nil)
	th, err := thumbnail.Make(res.Doc, &thumbnail.Options{Size: 64, Raster: imagebdf.Options{NoSystemFonts: true}})
	if err != nil {
		t.Fatal(err)
	}
	if th.Mode != thumbnail.Crop || th.Image.Bounds().Dx() != 64 || th.Image.Bounds().Dy() != 64 {
		t.Fatalf("mode %v, size %v", th.Mode, th.Image.Bounds())
	}
	near := func(x, y int, r, g, b uint8) {
		t.Helper()
		c := th.Image.RGBAAt(x, y)
		if d := math.Abs(float64(c.R)-float64(r)) + math.Abs(float64(c.G)-float64(g)) + math.Abs(float64(c.B)-float64(b)); d > 60 {
			t.Errorf("pixel (%d, %d) = %v, want about #%02x%02x%02x", x, y, c, r, g, b)
		}
	}
	near(46, 18, 0xf6, 0xc3, 0x44)
	near(4, 4, 0x2f, 0x5f, 0x9e)
	near(32, 60, 0x3a, 0x9d, 0x3a)
	// without a cover, the placeholder
	res = convert(t, "plain.aac", nil)
	if th, err = thumbnail.Make(res.Doc, &thumbnail.Options{Size: 64, Raster: imagebdf.Options{NoSystemFonts: true}}); err != nil || th.Mode != thumbnail.Crop {
		t.Fatalf("mode %v, err %v", th.Mode, err)
	}
	near(2, 2, 0xe6, 0xea, 0xee)
}

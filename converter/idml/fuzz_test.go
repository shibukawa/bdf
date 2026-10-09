package idml

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"
	"time"
)

// FuzzConvert converts arbitrary packages: each converts or fails, soon,
// without an internal error. The seeds are the test documents and their
// parts with one part replaced by the fuzzed bytes.
func FuzzConvert(f *testing.F) {
	for _, name := range []string{"testdata/basic.idml", "testdata/vertical.idml"} {
		b, err := os.ReadFile(name)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		start := time.Now()
		opts := &Options{NoSystemFonts: true, NoTextIndex: true, Warn: func(m string) {
			if bytes.Contains([]byte(m), []byte("internal error")) {
				t.Error(m)
			}
		}}
		Convert(bytes.NewReader(b), int64(len(b)), opts)
		if d := time.Since(start); d > 2*time.Second {
			t.Errorf("%v to convert %d bytes", d, len(b))
		}
	})
}

// FuzzStory converts a document whose story is the fuzzed bytes.
func FuzzStory(f *testing.F) {
	f.Add(`<idPkg:Story><Story Self="u100"><ParagraphStyleRange><CharacterStyleRange><Content>text<?ACE 18?></Content><Br/></CharacterStyleRange></ParagraphStyleRange></Story></idPkg:Story>`)
	f.Add(`<idPkg:Story><Story Self="u100"><StoryPreference StoryOrientation="Vertical"/><ParagraphStyleRange Justification="FullyJustified"><CharacterStyleRange PointSize="1e9" Tracking="-1e9"><Content>縦組み&#x2028;x</Content></CharacterStyleRange></ParagraphStyleRange></Story></idPkg:Story>`)
	f.Fuzz(func(t *testing.T, story string) {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		for name, data := range map[string]string{
			"designmap.xml":          designmap(`<idPkg:Spread src="Spreads/Spread_ucc.xml"/><idPkg:Story src="Stories/Story_u100.xml"/>`),
			"Spreads/Spread_ucc.xml": spreadPart(`<TextFrame Self="u10" ItemLayer="ub3" ItemTransform="1 0 0 1 0 0" ParentStory="u100"><Properties><PathGeometry><GeometryPathType PathOpen="false"><PathPointArray><PathPointType Anchor="10 -90"/><PathPointType Anchor="10 -10"/><PathPointType Anchor="90 -10"/><PathPointType Anchor="90 -90"/></PathPointArray></GeometryPathType></PathGeometry></Properties><TextFramePreference TextColumnCount="2"/></TextFrame>`),
			"Stories/Story_u100.xml": story,
		} {
			w, _ := zw.Create(name)
			w.Write([]byte(data))
		}
		zw.Close()
		start := time.Now()
		opts := &Options{NoSystemFonts: true, NoTextIndex: true, Warn: func(m string) {
			if bytes.Contains([]byte(m), []byte("internal error")) {
				t.Error(m)
			}
		}}
		Convert(bytes.NewReader(buf.Bytes()), int64(buf.Len()), opts)
		if d := time.Since(start); d > 2*time.Second {
			t.Errorf("%v to convert", d)
		}
	})
}

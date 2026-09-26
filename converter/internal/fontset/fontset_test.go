package fontset

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
)

// testFonts are the M PLUS 1p subsets of the PowerPoint converter's tests
// (ASCII, Latin-1, kana and some kanji).
const testFonts = "../../pptx/testdata/fonts"

func TestEstimatedMetrics(t *testing.T) {
	var warnings []string
	s := New(fontdb.New(nil, nil, false), func(msg string) { warnings = append(warnings, msg) })
	fc := s.FaceFor("Calibri", "ＭＳ 明朝", true, false, 'a')
	if fc.Loaded != nil || fc.Use.Requested != "Calibri" || !fc.Use.Bold {
		t.Errorf("choice = %+v", fc)
	}
	if ea := s.FaceFor("Calibri", "ＭＳ 明朝", true, false, 'あ'); ea.Use.Requested != "ＭＳ 明朝" {
		t.Errorf("East Asian text uses %q", ea.Use.Requested)
	}
	if s.Choose("Calibri", true, false, false) != fc {
		t.Error("choices are not cached")
	}
	for r, want := range map[rune]float64{'a': 0.55, ' ': 0.25, 'あ': 1, '​': 0} {
		if got := s.Advance(fc, r); got != want {
			t.Errorf("Advance(%q) = %g, want %g", r, got, want)
		}
	}
	// without font files nothing is missing: the viewer's fonts draw the text
	s.ReportMissing()
	if len(warnings) != 0 {
		t.Errorf("warnings: %v", warnings)
	}
	f := s.Font(fc.Use)
	if f.Kind != bdf.FontSystem || f.Family != `"Calibri", sans-serif` || f.Weight != 700 {
		t.Errorf("font record = %+v", f)
	}
}

func TestEmbed(t *testing.T) {
	var warnings []string
	s := New(fontdb.New(nil, []string{testFonts}, false), func(msg string) { warnings = append(warnings, msg) })
	fc := s.FaceFor("M PLUS 1p", "", false, false, 'A')
	if fc.Loaded == nil {
		t.Fatal("test font not found")
	}
	for _, r := range "Aあ" {
		if w := s.Advance(s.FaceFor("M PLUS 1p", "", false, false, r), r); w <= 0 || w > 1.5 {
			t.Errorf("Advance(%q) = %g", r, w)
		}
	}
	s.Advance(fc, '𝄞') // no test font has it
	doc := bdf.NewDocument()
	if n := s.Embed(doc, EmbedOptions{NoWOFF2: true}); n != 1 {
		t.Fatalf("embedded %d fonts", n)
	}
	f := s.Font(fc.Use)
	if f.Kind != bdf.FontEmbedded || doc.Part(f.Hash) == nil {
		t.Errorf("font record = %+v", f)
	}
	if size := len(doc.Part(f.Hash).Data); size >= fc.Loaded.Size() {
		t.Errorf("the embedded font (%d bytes) is not a subset of %d bytes", size, fc.Loaded.Size())
	}
	s.ReportMissing()
	if len(warnings) != 1 || !strings.Contains(warnings[0], "U+1D11E") {
		t.Errorf("warnings: %v", warnings)
	}
	// a face that measured nothing is referred to by name
	if f := s.Font(Use{Requested: "Nowhere", Generic: "serif"}); f.Kind != bdf.FontSystem || f.Family != `"Nowhere", serif` {
		t.Errorf("font record = %+v", f)
	}
}

func TestFaceForFamilies(t *testing.T) {
	s := New(fontdb.New(nil, []string{testFonts}, false), nil)
	list := []string{"M PLUS 1p", "Nowhere"}
	// the first family draws what it has, CJK or not
	fc := s.FaceForFamilies(list, false, false, 'A')
	if fc.Loaded == nil || fc.Use.Requested != "M PLUS 1p" {
		t.Errorf("choice = %+v", fc.Use)
	}
	if ea := s.FaceForFamilies(list, false, false, 'あ'); ea.Loaded == nil || ea.Use.Requested != "M PLUS 1p" {
		t.Errorf("East Asian text: choice = %+v", ea.Use)
	}
	// a character no font has stays with the first family
	if o := s.FaceForFamilies(list, false, false, '𝄞'); o != fc {
		t.Errorf("choice = %+v", o.Use)
	}
	// no list: draw.io's default family
	if fc := s.FaceForFamilies(nil, true, false, 'A'); fc.Use.Requested != "Helvetica" || !fc.Use.Bold {
		t.Errorf("choice = %+v", fc.Use)
	}
}

func TestSymbol(t *testing.T) {
	if !IsSymbol("Wingdings") || !IsSymbol("SYMBOL") || IsSymbol("Arial") {
		t.Error("IsSymbol")
	}
	for _, c := range [][3]string{
		{"Wingdings", "Ø", "▪➢"},
		{"Symbol", "·a", "•α"},
		{"Arial", "", ""},
	} {
		if got := MapSymbol(c[0], c[1]); got != c[2] {
			t.Errorf("MapSymbol(%s, %q) = %q, want %q", c[0], c[1], got, c[2])
		}
	}
}

func TestScript(t *testing.T) {
	for lang, want := range map[string]string{"ja-JP": "Jpan", "zh-TW": "Hant", "zh-Hant-HK": "Hant", "zh-CN": "Hans", "ko": "Hang", "en-US": "", "": ""} {
		if got := Script(lang); got != want {
			t.Errorf("Script(%q) = %q, want %q", lang, got, want)
		}
	}
	for r, want := range map[rune]string{'あ': "ja", 'ア': "ja", '한': "ko", '漢': "", 'a': ""} {
		if got := RuneLang(r); got != want {
			t.Errorf("RuneLang(%q) = %q, want %q", r, got, want)
		}
	}
}

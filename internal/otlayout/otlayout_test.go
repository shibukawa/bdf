package otlayout

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/shibukawa/bdf/internal/sfnt"
)

// The test fonts are those of the font converter (test/font/gen.py).
// features.ttf is made from a feature file with a lookup of each kind; its
// glyphs are .notdef space a b c d e f i a.alt b.alt c.alt f_i f_f_i
// acutecomb, numbered from 0. What this package reads of the fonts was
// checked lookup by lookup against fontTools (types, substitutions, rules
// and the lookups they apply, pair counts, marks and bases).
const (
	gA, gB, gC, gD, gE, gF, gI = 2, 3, 4, 5, 6, 7, 8
	gAalt, gBalt, gCalt        = 9, 10, 11
	gFI, gFFI, gAcute          = 12, 13, 14
)

func loadFont(t testing.TB, name string) *sfnt.Font {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("../../converter/font/testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	f, err := sfnt.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func tables(t *testing.T, name string) (gsub, gpos *Table) {
	t.Helper()
	f := loadFont(t, name)
	gsub, err := ParseGSUB(f.Tables["GSUB"])
	if err != nil {
		t.Fatal(err)
	}
	gpos, err = ParseGPOS(f.Tables["GPOS"])
	if err != nil {
		t.Fatal(err)
	}
	return gsub, gpos
}

func allSubsts(t *Table, i int) []Subst {
	var out []Subst
	t.Substs(i, func(s Subst) bool { out = append(out, s); return true })
	return out
}

// feature returns the lookups of the first feature with a tag.
func feature(t *testing.T, tb *Table, tag string) *Feature {
	t.Helper()
	for i := range tb.Features {
		if tb.Features[i].Tag == tag {
			return &tb.Features[i]
		}
	}
	t.Fatalf("no feature %s", tag)
	return nil
}

func TestScriptsAndFeatures(t *testing.T) {
	gsub, _ := tables(t, "features.ttf")
	var tags []string
	for _, s := range gsub.Scripts {
		tags = append(tags, s.Tag)
	}
	if !slices.Equal(tags, []string{"DFLT", "latn"}) {
		t.Fatalf("scripts %v", tags)
	}
	latn := gsub.Scripts[1]
	if latn.Default == nil || len(latn.Langs) != 1 || latn.Langs[0].Tag != "TRK " {
		t.Fatalf("latn %+v", latn)
	}
	users := gsub.ScriptFeatures()
	for i, f := range gsub.Features {
		want := []string{"DFLT", "latn", "latn/TRK "}
		if f.Tag == "locl" {
			want = []string{"latn/TRK "}
		}
		if !slices.Equal(users[i], want) {
			t.Errorf("%s is used by %v, want %v", f.Tag, users[i], want)
		}
	}
	cv := feature(t, gsub, "cv01")
	if cv.UIName != 256 || string(cv.Chars) != "ab" {
		t.Errorf("cv01 params: name %d, characters %q", cv.UIName, string(cv.Chars))
	}
	if ss := feature(t, gsub, "ss01"); ss.UIName != 257 {
		t.Errorf("ss01 name %d", ss.UIName)
	}
}

func TestSubstitutions(t *testing.T) {
	gsub, _ := tables(t, "features.ttf")
	lookup := func(tag string) int { return feature(t, gsub, tag).Lookups[0] }
	cases := []struct {
		tag  string
		typ  int
		want []Subst
	}{
		{"smcp", SubstSingle, []Subst{{[]uint16{gA}, []uint16{gAalt}}, {[]uint16{gB}, []uint16{gBalt}}, {[]uint16{gC}, []uint16{gCalt}}}},
		{"locl", SubstSingle, []Subst{{[]uint16{gA}, []uint16{gE}}, {[]uint16{gB}, []uint16{gC}}}},
		{"ccmp", SubstMultiple, []Subst{{[]uint16{gD}, []uint16{gA, gB}}}},
		{"salt", SubstAlternate, []Subst{{[]uint16{gA}, []uint16{gAalt, gBalt, gCalt}}}},
		{"liga", SubstLigature, []Subst{{[]uint16{gF, gF, gI}, []uint16{gFFI}}, {[]uint16{gF, gI}, []uint16{gFI}}}},
		{"rclt", SubstReverseChain, []Subst{{[]uint16{gC}, []uint16{gE}}}},
	}
	for _, c := range cases {
		i := lookup(c.tag)
		if typ := gsub.Lookups[i].Type; typ != c.typ {
			t.Errorf("%s: lookup type %d, want %d", c.tag, typ, c.typ)
		}
		got := allSubsts(gsub, i)
		if !slices.EqualFunc(got, c.want, func(a, b Subst) bool { return slices.Equal(a.In, b.In) && slices.Equal(a.Out, b.Out) }) {
			t.Errorf("%s: %v, want %v", c.tag, got, c.want)
		}
	}
	// calt is a chained context in an extension lookup, applying a lookup
	// made for its rule
	i := lookup("calt")
	if gsub.Lookups[i].Type != SubstChainContext {
		t.Fatalf("calt: type %d", gsub.Lookups[i].Type)
	}
	c, ok := gsub.Context(i)
	if !ok || c.Rules != 1 || len(c.Nested) != 1 {
		t.Fatalf("calt: %+v", c)
	}
	if got := allSubsts(gsub, c.Nested[0]); len(got) != 1 || got[0].In[0] != gA || got[0].Out[0] != gAalt {
		t.Errorf("calt applies %v", got)
	}
	if cov := gsub.Covered(i); !slices.Equal(cov, []uint16{gA}) {
		t.Errorf("calt covers %v", cov)
	}
	if name := gsub.TypeName(SubstChainContext); name != "chained contextual" {
		t.Errorf("type name %q", name)
	}
}

func TestPositioning(t *testing.T) {
	_, gpos := tables(t, "features.ttf")
	kern := feature(t, gpos, "kern").Lookups[0]
	for _, c := range []struct {
		a, b uint16
		adv  int
		ok   bool
	}{{gA, gB, -50, true}, {gC, gE, -30, true}, {gD, gF, -30, true}, {gA, gC, 0, false}, {gB, gA, 0, false}} {
		v1, v2, ok := gpos.PairValue(kern, c.a, c.b)
		if ok && !(v1.Zero() && v2.Zero()) != c.ok || v1.XAdvance != c.adv {
			t.Errorf("pair %d %d: %+v %+v %v", c.a, c.b, v1, v2, ok)
		}
	}
	if n := gpos.PairCount(kern); n != 5 {
		t.Errorf("%d pairs, want 5", n)
	}
	var pairs []Pair
	gpos.Pairs(kern, func(p Pair) bool { pairs = append(pairs, p); return true })
	if len(pairs) != 5 {
		t.Errorf("pairs %v", pairs)
	}
	var singles []Value
	gpos.Singles(feature(t, gpos, "palt").Lookups[0], func(g uint16, v Value) bool {
		if g != gA {
			t.Errorf("single adjustment of glyph %d", g)
		}
		singles = append(singles, v)
		return true
	})
	if len(singles) != 1 || singles[0] != (Value{XPlacement: -10, XAdvance: -20}) {
		t.Errorf("singles %v", singles)
	}
	at := gpos.Attachments(feature(t, gpos, "mark").Lookups[0])
	if len(at) != 1 || at[0].Marks[gAcute] != (MarkAnchor{0, Anchor{100, 500}}) || len(at[0].Bases) != 2 || *at[0].Bases[gB][0] != (Anchor{250, 450}) {
		t.Errorf("mark to base %+v", at)
	}
	mk := gpos.Attachments(feature(t, gpos, "mkmk").Lookups[0])
	if len(mk) != 1 || *mk[0].Bases[gAcute][0] != (Anchor{100, 700}) {
		t.Errorf("mark to mark %+v", mk)
	}
	if m, b := gpos.MarkCounts(feature(t, gpos, "curs").Lookups[0]); m != 0 || b != 1 {
		t.Errorf("cursive: %d %d", m, b)
	}
}

func TestGDEF(t *testing.T) {
	f := loadFont(t, "features.ttf")
	g := ParseGDEF(f.Tables["GDEF"])
	if g == nil || g.Classes[gA] != GlyphBase || g.Classes[gFI] != GlyphLigature || g.Classes[gAcute] != GlyphMark || g.Classes[1] != 0 {
		t.Fatalf("GDEF %+v", g)
	}
}

// TestSTIX checks values of a real font against fontTools.
func TestSTIX(t *testing.T) {
	gsub, gpos := tables(t, "stix.ttf")
	if len(gsub.Lookups) != 34 || len(gpos.Lookups) != 6 {
		t.Fatalf("%d GSUB and %d GPOS lookups", len(gsub.Lookups), len(gpos.Lookups))
	}
	const a, aSmcp, A, V, T, circumflex, grave, both, acute = 60, 196, 3, 24, 22, 292, 288, 294, 290
	if s := allSubsts(gsub, 8); !slices.ContainsFunc(s, func(s Subst) bool { return s.In[0] == a && s.Out[0] == aSmcp }) {
		t.Errorf("smcp: a is not a.smcp")
	}
	if s := allSubsts(gsub, 1); len(s) == 0 || !slices.Equal(s[0].In, []uint16{circumflex, grave}) || s[0].Out[0] != both {
		t.Errorf("ccmp ligature %v", s)
	}
	kern := gpos.Features[0].Lookups[0]
	for _, c := range []struct {
		a, b uint16
		adv  int
	}{{A, V, -100}, {T, a, -60}} {
		if v1, _, _ := gpos.PairValue(kern, c.a, c.b); v1.XAdvance != c.adv {
			t.Errorf("kern %d %d: %d, want %d", c.a, c.b, v1.XAdvance, c.adv)
		}
	}
	at := gpos.Attachments(1)
	if len(at) == 0 || at[0].Marks[acute].Anchor != (Anchor{-230, 473}) || *at[0].Bases[a][at[0].Marks[acute].Class] != (Anchor{231, 473}) {
		t.Errorf("mark: acute on a %+v", at)
	}
}

func TestNames(t *testing.T) {
	for tag, want := range map[string]string{"liga": "Standard Ligatures", "ss07": "Stylistic Set 7", "cv42": "Character Variant 42", "zzzz": "", "ss21": ""} {
		if got := FeatureName(tag); got != want {
			t.Errorf("%s: %q, want %q", tag, got, want)
		}
	}
	if FeatureDefault("liga") != Default || FeatureDefault("init") != ScriptDefault || FeatureDefault("vert") != Vertical || FeatureDefault("smcp") != Optional {
		t.Error("feature defaults")
	}
	if ScriptName("latn") != "Latin" || LangName("TRK ") != "Turkish" || LangName("XYZ") != "" {
		t.Error("script and language names")
	}
}

// TestDamaged reads every truncation of the tables: nothing panics.
func TestDamaged(t *testing.T) {
	f := loadFont(t, "features.ttf")
	for _, tag := range []string{"GSUB", "GPOS", "GDEF"} {
		b := f.Tables[tag]
		for n := 0; n < len(b); n++ {
			exercise(b[:n])
		}
	}
}

// exercise reads a table as GSUB, GPOS and GDEF with every method.
func exercise(b []byte) {
	ParseGDEF(b)
	for _, gpos := range []bool{false, true} {
		t, err := parse(b, gpos)
		if err != nil {
			continue
		}
		t.ScriptFeatures()
		for i := range t.Lookups {
			t.Substs(i, func(Subst) bool { return true })
			t.Context(i)
			t.Covered(i)
			t.PairCount(i)
			n := 0
			t.Pairs(i, func(Pair) bool { n++; return n < 1000 })
			t.PairValue(i, 2, 3)
			t.Singles(i, func(uint16, Value) bool { return true })
			t.Attachments(i)
			t.MarkCounts(i)
		}
	}
}

func FuzzParse(f *testing.F) {
	for _, name := range []string{"features.ttf", "stix.otf"} {
		fn := loadFont(f, name)
		for _, tag := range []string{"GSUB", "GPOS", "GDEF"} {
			f.Add(fn.Tables[tag])
		}
	}
	f.Fuzz(func(t *testing.T, b []byte) { exercise(b) })
}

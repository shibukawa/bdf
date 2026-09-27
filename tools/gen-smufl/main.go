// Command gen-smufl writes the music symbols that the score engraver
// embeds (converter/internal/music/smufl): it reads a release of Bravura,
// the SMuFL reference font, and its metadata, and writes the outlines of
// the glyphs the engraver draws, their advance widths, bounding boxes and
// anchors (where stems and flags attach), and the font's engraving defaults
// (line thicknesses) as a Go table, together with the license.
//
// Usage (from the repository root):
//
//	go run ./tools/gen-smufl [-src URL-or-dir] [-out dir]
//
// -src holds Bravura.otf, Bravura.json and OFL.txt: by default the release
// on GitHub the table was made from, or a local directory.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/internal/sfnt"
)

// release is the Bravura version the table comes from.
const release = "bravura-1.482"

// glyphs are the SMuFL glyphs embedded, by name and code point.
var glyphs = []struct {
	name string
	cp   rune
}{
	{"noteheadDoubleWhole", 0xE0A0}, {"noteheadWhole", 0xE0A2}, {"noteheadHalf", 0xE0A3}, {"noteheadBlack", 0xE0A4},
	{"noteheadXWhole", 0xE0A7}, {"noteheadXHalf", 0xE0A8}, {"noteheadXBlack", 0xE0A9},

	{"gClef", 0xE050}, {"gClef8vb", 0xE052}, {"gClef8va", 0xE053}, {"cClef", 0xE05C}, {"fClef", 0xE062},
	{"fClef8vb", 0xE064}, {"unpitchedPercussionClef1", 0xE069}, {"6stringTabClef", 0xE06D},
	{"gClefChange", 0xE07A}, {"cClefChange", 0xE07B}, {"fClefChange", 0xE07C},

	{"accidentalFlat", 0xE260}, {"accidentalNatural", 0xE261}, {"accidentalSharp", 0xE262},
	{"accidentalDoubleSharp", 0xE263}, {"accidentalDoubleFlat", 0xE264},
	{"accidentalParensLeft", 0xE26A}, {"accidentalParensRight", 0xE26B},

	{"restDoubleWhole", 0xE4E2}, {"restWhole", 0xE4E3}, {"restHalf", 0xE4E4}, {"restQuarter", 0xE4E5},
	{"rest8th", 0xE4E6}, {"rest16th", 0xE4E7}, {"rest32nd", 0xE4E8}, {"rest64th", 0xE4E9}, {"rest128th", 0xE4EA},
	{"restWholeLegerLine", 0xE4F4}, {"restHalfLegerLine", 0xE4F5},

	{"flag8thUp", 0xE240}, {"flag8thDown", 0xE241}, {"flag16thUp", 0xE242}, {"flag16thDown", 0xE243},
	{"flag32ndUp", 0xE244}, {"flag32ndDown", 0xE245}, {"flag64thUp", 0xE246}, {"flag64thDown", 0xE247},
	{"flag128thUp", 0xE248}, {"flag128thDown", 0xE249},

	{"augmentationDot", 0xE1E7}, {"repeatDots", 0xE043}, {"repeatDot", 0xE044},

	{"timeSig0", 0xE080}, {"timeSig1", 0xE081}, {"timeSig2", 0xE082}, {"timeSig3", 0xE083}, {"timeSig4", 0xE084},
	{"timeSig5", 0xE085}, {"timeSig6", 0xE086}, {"timeSig7", 0xE087}, {"timeSig8", 0xE088}, {"timeSig9", 0xE089},
	{"timeSigCommon", 0xE08A}, {"timeSigCutCommon", 0xE08B},

	{"tuplet0", 0xE880}, {"tuplet1", 0xE881}, {"tuplet2", 0xE882}, {"tuplet3", 0xE883}, {"tuplet4", 0xE884},
	{"tuplet5", 0xE885}, {"tuplet6", 0xE886}, {"tuplet7", 0xE887}, {"tuplet8", 0xE888}, {"tuplet9", 0xE889},
	{"tupletColon", 0xE88A},

	{"dynamicPiano", 0xE520}, {"dynamicMezzo", 0xE521}, {"dynamicForte", 0xE522}, {"dynamicRinforzando", 0xE523},
	{"dynamicSforzando", 0xE524}, {"dynamicZ", 0xE525}, {"dynamicNiente", 0xE526},

	{"articAccentAbove", 0xE4A0}, {"articAccentBelow", 0xE4A1}, {"articStaccatoAbove", 0xE4A2},
	{"articStaccatoBelow", 0xE4A3}, {"articTenutoAbove", 0xE4A4}, {"articTenutoBelow", 0xE4A5},
	{"articStaccatissimoAbove", 0xE4A6}, {"articStaccatissimoBelow", 0xE4A7},
	{"articMarcatoAbove", 0xE4AC}, {"articMarcatoBelow", 0xE4AD},
	{"fermataAbove", 0xE4C0}, {"fermataBelow", 0xE4C1}, {"breathMarkComma", 0xE4CE}, {"caesura", 0xE4D1},

	{"ornamentTrill", 0xE566}, {"ornamentTurn", 0xE567}, {"ornamentShortTrill", 0xE56C}, {"ornamentMordent", 0xE56D},

	{"segno", 0xE047}, {"coda", 0xE048}, {"brace", 0xE000}, {"bracketTop", 0xE003}, {"bracketBottom", 0xE004},

	{"metNoteWhole", 0xECA2}, {"metNoteHalfUp", 0xECA3}, {"metNoteQuarterUp", 0xECA5}, {"metNote8thUp", 0xECA7},
	{"metNote16thUp", 0xECA9}, {"metAugmentationDot", 0xECB7},

	{"ottava", 0xE510}, {"ottavaAlta", 0xE511}, {"quindicesima", 0xE514}, {"quindicesimaAlta", 0xE515},
	{"ottavaBassaVb", 0xE51C},

	{"keyboardPedalPed", 0xE650}, {"keyboardPedalUp", 0xE655},

	{"graceNoteAcciaccaturaStemUp", 0xE560}, {"graceNoteAppoggiaturaStemUp", 0xE562},

	{"tremolo1", 0xE220}, {"tremolo2", 0xE221}, {"tremolo3", 0xE222},

	{"fingering0", 0xED10}, {"fingering1", 0xED11}, {"fingering2", 0xED12}, {"fingering3", 0xED13},
	{"fingering4", 0xED14}, {"fingering5", 0xED15},
}

// metadata is the part of Bravura.json the table uses.
type metadata struct {
	FontVersion       float64                         `json:"fontVersion"`
	EngravingDefaults map[string]any                  `json:"engravingDefaults"`
	GlyphBBoxes       map[string]map[string][]float64 `json:"glyphBBoxes"`
	GlyphsWithAnchors map[string]map[string][]float64 `json:"glyphsWithAnchors"`
}

func main() {
	src := flag.String("src", "https://github.com/steinbergmedia/bravura/releases/download/"+release, "Bravura release (URL or directory)")
	out := flag.String("out", "converter/internal/music/smufl", "output directory")
	flag.Parse()

	otf, err := read(*src, "Bravura.otf")
	if err != nil {
		fail(err)
	}
	js, err := read(*src, "Bravura.json")
	if err != nil {
		fail(err)
	}
	license, err := read(*src, "OFL.txt")
	if err != nil {
		fail(err)
	}
	var meta metadata
	if err := json.Unmarshal(js, &meta); err != nil {
		fail(fmt.Errorf("Bravura.json: %w", err))
	}
	font, err := sfnt.Parse(otf)
	if err != nil {
		fail(fmt.Errorf("Bravura.otf: %w", err))
	}
	if font.UnitsPerEm != 1000 {
		fail(fmt.Errorf("Bravura.otf: %d units per em, want 1000", font.UnitsPerEm))
	}
	outlines := sfnt.NewOutlines(font)

	var b bytes.Buffer
	fmt.Fprintf(&b, "// Code generated by tools/gen-smufl from Bravura %s; DO NOT EDIT.\n\n", strings.TrimPrefix(release, "bravura-"))
	b.WriteString("package smufl\n\n")

	// Engraving defaults that are numbers, in staff spaces.
	var keys []string
	for k, v := range meta.EngravingDefaults {
		if _, ok := v.(float64); ok {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	b.WriteString("// defaults are Bravura's engraving defaults, in staff spaces.\nvar defaults = Defaults{\n")
	for _, k := range keys {
		fmt.Fprintf(&b, "\t%s: %s,\n", strings.ToUpper(k[:1])+k[1:], num(meta.EngravingDefaults[k].(float64)))
	}
	b.WriteString("}\n\n")

	b.WriteString("// glyphs are the embedded glyphs in font units (1000 per em, 250 per staff\n// space, y up).\nvar glyphs = map[string]*Glyph{\n")
	for _, g := range glyphs {
		gid, ok := font.Cmap[uint32(g.cp)]
		if !ok {
			fail(fmt.Errorf("%s (U+%04X) is not in the font", g.name, g.cp))
		}
		var p pathWriter
		if !outlines.Outline(gid, &p) {
			fail(fmt.Errorf("%s: cannot read the outline", g.name))
		}
		bbox := meta.GlyphBBoxes[g.name]
		if bbox == nil {
			fail(fmt.Errorf("%s has no bounding box in Bravura.json", g.name))
		}
		fmt.Fprintf(&b, "\t%q: {Rune: 0x%04X, Advance: %d, BBox: [4]float64{%s, %s, %s, %s}", g.name, g.cp, font.Advance(gid),
			units(bbox["bBoxSW"][0]), units(bbox["bBoxSW"][1]), units(bbox["bBoxNE"][0]), units(bbox["bBoxNE"][1]))
		if a := meta.GlyphsWithAnchors[g.name]; len(a) > 0 {
			var names []string
			for n := range a {
				names = append(names, n)
			}
			sort.Strings(names)
			b.WriteString(", Anchors: map[string][2]float64{")
			for i, n := range names {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(&b, "%q: {%s, %s}", n, units(a[n][0]), units(a[n][1]))
			}
			b.WriteString("}")
		}
		fmt.Fprintf(&b, ",\n\t\tPath: %q},\n", p.String())
	}
	b.WriteString("}\n")

	src2, err := format.Source(b.Bytes())
	if err != nil {
		fail(fmt.Errorf("format: %w", err))
	}
	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "bravura.go"), src2, 0o644); err != nil {
		fail(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "OFL.txt"), license, 0o644); err != nil {
		fail(err)
	}
	fmt.Printf("%d glyphs, %d bytes\n", len(glyphs), len(src2))
}

// pathWriter writes an outline as SVG path data in font units ("M10 20L…Z").
type pathWriter struct{ strings.Builder }

func (p *pathWriter) pt(cmd byte, xy ...float64) {
	p.WriteByte(cmd)
	for i, v := range xy {
		if i > 0 && v >= 0 {
			p.WriteByte(' ')
		}
		p.WriteString(num(v))
	}
}

func (p *pathWriter) MoveTo(x, y float64)             { p.pt('M', x, y) }
func (p *pathWriter) LineTo(x, y float64)             { p.pt('L', x, y) }
func (p *pathWriter) QuadTo(cx, cy, x, y float64)     { p.pt('Q', cx, cy, x, y) }
func (p *pathWriter) CubeTo(a, b, c, d, x, y float64) { p.pt('C', a, b, c, d, x, y) }
func (p *pathWriter) Close()                          { p.WriteByte('Z') }

// num writes v with at most two decimals.
func num(v float64) string {
	v = math.Round(v*100) / 100
	if v == 0 {
		v = 0 // no -0
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// units converts staff spaces to font units.
func units(v float64) string { return num(v * 250) }

func read(src, name string) ([]byte, error) {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		resp, err := http.Get(src + "/" + name)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: %s", name, resp.Status)
		}
		return io.ReadAll(resp.Body)
	}
	return os.ReadFile(filepath.Join(src, name))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen-smufl:", err)
	os.Exit(1)
}

// Package font shows font files: TrueType and OpenType fonts (.ttf, .otf),
// collections of them (.ttc, .otc) and web fonts (.woff, .woff2). The
// format registered as "font" makes a document of scroll views for each
// font of the file:
//
//   - Overview: the name large, a specimen and the text of the scripts
//     the font covers at several sizes, then the name table, the license
//     and embedding permissions, the metrics, the variation axes and named
//     instances, the scripts and features of the layout tables, the color
//     palettes and the list of tables.
//   - Characters: the characters the font maps, by Unicode block, with the
//     block's coverage, as the code charts of the Unicode Standard.
//   - Glyphs: every glyph by glyph ID with its name and characters.
//   - Features: the OpenType features of GSUB and GPOS, with the scripts
//     and languages that use them and whether browsers apply them without
//     being asked; the substitutions of each (f + i → ﬁ, a → ᴀ, the
//     alternates of a glyph), and kerning, single adjustments and mark
//     attachments as pairs drawn before and after.
//
// BDF draws text with the Canvas 2D API, which applies no feature a
// document asks for, so a feature's effect is shown glyph by glyph rather
// than as shaped text. The glyphs are drawn with the font itself, embedded
// twice: as a "glyph font" whose cmap maps every glyph from a private use
// character and which has no layout tables, so that each glyph is drawn as
// it is (color glyphs included), and as a "sample font" with its layout
// tables but only the glyphs the sample text can reach, so that the
// browser shapes the sample text (joining Arabic letters, forming Indic
// conjuncts) as applications do. A font whose OS/2 fsType forbids
// embedding it (unless Options.IgnoreFSType) is not embedded: its glyphs
// are drawn as outlines and its sample text glyph by glyph. A variable font
// is drawn at its default instance.
package font

import (
	"encoding/binary"
	"fmt"
	"io"

	conv "github.com/shibukawa/bdf/converter"
)

// maxSize bounds the input read.
const maxSize = 256 << 20

func init() {
	conv.Register(&conv.Format{
		Name:        "font",
		Description: "Font file (TrueType, OpenType, WOFF)",
		Extensions:  []string{".ttf", ".otf", ".ttc", ".otc", ".woff", ".woff2"},
		Params: []conv.Param{
			{Name: "font", Usage: "fonts of a collection to show, 1-based and separated by commas, or all (default: all, or the first ones of a collection of more than 262,144 glyphs)"},
			{Name: "text", Usage: "sample text shown at several sizes instead of the pangram of the font's script"},
			{Name: "examples", Usage: "substitutions and pairs shown per feature, or all (default 200)"},
		},
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return detect(head) },
		Convert: func(r io.ReaderAt, size int64, o *conv.Options) (*conv.Result, error) {
			if size > maxSize {
				return nil, fmt.Errorf("font: the file is larger than %d bytes", maxSize)
			}
			data := make([]byte, size)
			if n, err := r.ReadAt(data, 0); n < len(data) {
				if err == nil {
					err = io.ErrUnexpectedEOF
				}
				return nil, fmt.Errorf("font: %w", err)
			}
			var warnings []string
			warn := o.Warn
			if warn == nil {
				warn = func(msg string) { warnings = append(warnings, msg) }
			}
			res, err := convert(data, o, warn)
			if err != nil {
				return nil, err
			}
			res.Warnings = append(warnings, res.Warnings...)
			return res, nil
		},
	})
}

// detect tells a font file by its signature and, for sfnt files, a table
// directory that holds together.
func detect(head []byte) bool {
	if len(head) < 12 {
		return false
	}
	switch string(head[:4]) {
	case "wOFF", "wOF2":
		return true
	case "ttcf":
		v := binary.BigEndian.Uint16(head[4:])
		n := binary.BigEndian.Uint32(head[8:])
		return (v == 1 || v == 2) && n >= 1 && n <= 1<<16
	case "\x00\x01\x00\x00", "OTTO", "true":
	default:
		return false
	}
	n := int(binary.BigEndian.Uint16(head[4:]))
	if n < 1 || n > 256 {
		return false
	}
	// the table tags are printable, and those read include the tables
	// every font has
	seen := 0
	shown := min(n, (len(head)-12)/16)
	for i := range shown {
		tag := string(head[12+i*16 : 16+i*16])
		for _, c := range []byte(tag) {
			if c < 0x20 || c > 0x7E {
				return false
			}
		}
		switch tag {
		case "cmap", "head", "maxp":
			seen++
		}
	}
	return shown > 0 && (seen == 3 || shown < n)
}

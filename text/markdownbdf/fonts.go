package markdownbdf

import (
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/font/woff2"
	"github.com/shibukawa/bdf/internal/sfnt"
)

// face is what the tables of an embedded font say of its style.
type face struct {
	bold, italic bool
}

// style tells what the Markdown makes of a font: strong text for a bold
// one, emphasis for an italic one, code for a monospaced one. An embedded
// font is recorded with the weight and style the browser is to synthesize
// (docs/design.md), so a bold or italic face says so in its own OS/2
// table, which is read once per font.
func (w *writer) style(f *bdf.Font) (bold, italic, mono bool) {
	if f == nil {
		return false, false, false
	}
	bold = f.Weight >= 600
	italic = f.Style != bdf.StyleNormal
	mono = isMono(f.Family)
	if f.Kind == bdf.FontEmbedded {
		if fc := w.face(f.Hash); fc != nil {
			bold = bold || fc.bold
			italic = italic || fc.italic
		}
	}
	return bold, italic, mono
}

// face reads the style of an embedded font from its tables, once.
func (w *writer) face(h bdf.Hash) *face {
	if fc, ok := w.faces[h]; ok {
		return fc
	}
	var fc *face
	if p := w.doc.Part(h); p != nil {
		data := p.Data
		if woff2.IsWOFF(data) {
			if d, err := woff2.Decode(data); err == nil {
				data = d
			} else {
				data = nil
			}
		}
		if data != nil {
			if sf, err := sfnt.Parse(data); err == nil {
				fc = faceOf(sf)
			}
		}
	}
	w.faces[h] = fc
	return fc
}

// faceOf reads the weight class and the selection flags of the OS/2 table,
// or the style of the head table of a font without one.
func faceOf(sf *sfnt.Font) *face {
	fc := &face{}
	if os2 := sf.Tables["OS/2"]; len(os2) >= 64 {
		sel := sfnt.BE16(os2, 62)
		fc.bold = sfnt.BE16(os2, 4) >= 600 || sel&0x20 != 0
		fc.italic = sel&0x01 != 0
	} else if head := sf.Tables["head"]; len(head) >= 46 {
		mac := sfnt.BE16(head, 44)
		fc.bold, fc.italic = mac&1 != 0, mac&2 != 0
	}
	return fc
}

// isMono reports whether a CSS family list names a monospaced font: the
// generic family, or a well-known one.
func isMono(family string) bool {
	fam := strings.ToLower(family)
	for _, k := range []string{"monospace", "courier", "mono", "consolas", "menlo", "source code", "fira code", "cascadia"} {
		if strings.Contains(fam, k) {
			return true
		}
	}
	return false
}

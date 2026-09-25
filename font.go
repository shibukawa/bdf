package bdf

// Font kinds.
const (
	FontEmbedded byte = 0
	FontSystem   byte = 1
)

// Font styles.
const (
	StyleNormal  byte = 0
	StyleItalic  byte = 1
	StyleOblique byte = 2
)

// Font is a font reference used by an object.
type Font struct {
	Kind   byte
	Hash   Hash   // embedded: font part
	Family string // CSS family list; fallback for embedded, primary for system
	Weight uint16
	Style  byte
}

// EmbeddedFont references a font part added with Document.AddFont.
func EmbeddedFont(h Hash, weight uint16, style byte) Font {
	return Font{Kind: FontEmbedded, Hash: h, Weight: weight, Style: style}
}

// SystemFont references a font by CSS family list.
func SystemFont(family string, weight uint16, style byte) Font {
	return Font{Kind: FontSystem, Family: family, Weight: weight, Style: style}
}

func (f *Font) encode(w *buf) {
	w.u8(f.Kind)
	if f.Kind == FontEmbedded {
		w.hash(f.Hash)
	}
	w.str(f.Family)
	w.u16(f.Weight)
	w.u8(f.Style)
}

func decodeFont(r *reader) Font {
	f := Font{Kind: r.u8()}
	if f.Kind == FontEmbedded {
		f.Hash = r.hash()
	}
	f.Family = r.str()
	f.Weight = r.u16()
	f.Style = r.u8()
	return f
}

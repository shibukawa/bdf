package fontset

import (
	"sort"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/woff2"
)

// EmbedOptions controls how Embed stores fonts.
type EmbedOptions struct {
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores fonts as TrueType/OpenType instead of WOFF2.
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting.
	IgnoreFSType bool
}

// maxWholeFont is the largest font file embedded whole when it cannot be
// subset (CFF outlines).
const maxWholeFont = 2 << 20

// Embed adds the faces that measured characters to doc, as subsets of those
// characters, and returns how many it embedded. Call it once all the text
// has been measured; faces that cannot be embedded are referred to by name
// (see Font).
func (s *Set) Embed(doc *bdf.Document, opts EmbedOptions) int {
	faces := make([]*fontdb.Face, 0, len(s.runes))
	for f := range s.runes {
		faces = append(faces, f)
	}
	sort.Slice(faces, func(i, j int) bool {
		if faces[i].Path != faces[j].Path {
			return faces[i].Path < faces[j].Path
		}
		return faces[i].Index < faces[j].Index
	})
	for _, f := range faces {
		l, err := f.Load()
		if err != nil {
			continue
		}
		if !l.CanSubset(opts.IgnoreFSType) && !opts.NoSubset && l.Size() > maxWholeFont {
			s.warnf("font %s is not embedded: its outlines cannot be subset and the file is %d KB", f.Family, l.Size()/1024)
			continue
		}
		runes := make([]rune, 0, len(s.runes[f]))
		for r := range s.runes[f] {
			runes = append(runes, r)
		}
		data, ok := l.Program(runes, opts.NoSubset, opts.IgnoreFSType)
		if !ok {
			s.warnf("font %s is not embedded: its license does not allow embedding", f.Family)
			continue
		}
		if !opts.NoWOFF2 {
			if w, err := woff2.Encode(data); err == nil {
				data = w
			} else if err != woff2.ErrNotAvailable {
				s.warnf("font %s: not stored as WOFF2: %v", f.Family, err)
			}
		}
		s.embedded[f] = doc.AddFont(data)
	}
	return len(s.embedded)
}

// Font returns the font record for u: the embedded face, or the requested
// family followed by the resolved one and the generic family, for the
// viewer's fonts.
func (s *Set) Font(u Use) bdf.Font {
	if h, ok := s.embedded[u.Face]; ok && u.Face != nil {
		// The face registered in the browser is the one that was measured;
		// only styles it lacks are synthesized.
		f := bdf.EmbeddedFont(h, 400, bdf.StyleNormal)
		if u.SynthBold {
			f.Weight = 700
		}
		if u.SynthItalic {
			f.Style = bdf.StyleItalic
		}
		f.Family = u.Generic
		return f
	}
	family := cssQuote(u.Requested)
	if u.Face != nil && u.Face.Family != u.Requested {
		family += ", " + cssQuote(u.Face.Family)
	}
	if family != "" {
		family += ", "
	}
	family += u.Generic
	var w uint16 = 400
	if u.Bold {
		w = 700
	}
	st := bdf.StyleNormal
	if u.Italic {
		st = bdf.StyleItalic
	}
	return bdf.SystemFont(family, w, st)
}

func cssQuote(s string) string {
	if s == "" {
		return ""
	}
	out := []rune{'"'}
	for _, r := range s {
		if r == '"' || r == '\\' {
			out = append(out, '\\')
		}
		out = append(out, r)
	}
	return string(append(out, '"'))
}

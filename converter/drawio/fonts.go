package drawio

import (
	"strings"

	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Fonts: labels are laid out here with the metrics of real fonts, which
// are then embedded as subsets of the characters used (fontset), as the
// Office converters do (docs/design.md §3.4).

// finalize embeds the fonts the labels were measured with and encodes every
// layer object, whose font references stay placeholders until then.
func (c *converter) finalize() {
	if !c.opts.SystemFonts {
		c.embeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{
			NoSubset: c.opts.NoSubset, NoWOFF2: c.opts.NoWOFF2, IgnoreFSType: c.opts.IgnoreFSType})
	}
	c.objs.Encode()
}

// parseFontFamilies splits a CSS font-family list ("Helvetica, 'Noto Sans
// JP', sans-serif") into family names; generic names map to the fonts
// fontdb resolves them with.
func parseFontFamilies(s string) []string {
	var out []string
	for _, f := range splitArgs(s) {
		f = strings.TrimSpace(f)
		f = strings.Trim(f, `"'`)
		if f == "" {
			continue
		}
		switch strings.ToLower(f) {
		case "sans-serif", "system-ui", "-apple-system", "blinkmacsystemfont", "ui-sans-serif":
			f = fontdb.Sans
		case "serif", "ui-serif":
			f = fontdb.Serif
		case "monospace", "ui-monospace":
			f = fontdb.Mono
		}
		out = append(out, f)
	}
	return out
}

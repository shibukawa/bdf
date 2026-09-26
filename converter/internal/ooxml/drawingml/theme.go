package drawingml

import "github.com/shibukawa/bdf/converter/internal/ooxml"

// theme is the part of a DrawingML theme the renderer needs.
type theme struct {
	colors map[string]rgba // dk1, lt1, dk2, lt2, accent1-6, hlink, folHlink
	major  fontScheme
	minor  fontScheme
	// format scheme: style matrices referenced by p:style (1-based indices)
	fills, lines, effects, bgFills []*ooxml.Node
	part                           string // relationships of blip fills
}

type fontScheme struct {
	latin, ea, cs string
	scripts       map[string]string // script tag (Jpan, Hans, Hant, Hang …) → typeface
}

var defaultTheme = &theme{
	colors: map[string]rgba{
		"dk1": black, "lt1": white,
		"dk2": {0x44 / 255.0, 0x54 / 255.0, 0x6a / 255.0, 1}, "lt2": {0xe7 / 255.0, 0xe6 / 255.0, 0xe6 / 255.0, 1},
		"accent1": {0x44 / 255.0, 0x72 / 255.0, 0xc4 / 255.0, 1}, "accent2": {0xed / 255.0, 0x7d / 255.0, 0x31 / 255.0, 1},
		"accent3": {0xa5 / 255.0, 0xa5 / 255.0, 0xa5 / 255.0, 1}, "accent4": {1, 0xc0 / 255.0, 0, 1},
		"accent5": {0x5b / 255.0, 0x9b / 255.0, 0xd5 / 255.0, 1}, "accent6": {0x70 / 255.0, 0xad / 255.0, 0x47 / 255.0, 1},
		"hlink": {0x05 / 255.0, 0x63 / 255.0, 0xc1 / 255.0, 1}, "folHlink": {0x95 / 255.0, 0x4f / 255.0, 0x72 / 255.0, 1},
	},
	major: fontScheme{latin: "Calibri Light"},
	minor: fontScheme{latin: "Calibri"},
}

func parseTheme(root *ooxml.Node) *theme {
	if root == nil {
		return defaultTheme
	}
	te := root.Child("themeElements")
	th := &theme{colors: map[string]rgba{}}
	for k, v := range defaultTheme.colors {
		th.colors[k] = v
	}
	if cs := te.Child("clrScheme"); cs != nil {
		plain := &colorCtx{}
		for _, k := range cs.Kids {
			if c, ok := plain.color(k); ok {
				th.colors[k.Name] = c
			}
		}
	}
	fs := te.Child("fontScheme")
	th.major = parseFontScheme(fs.Child("majorFont"), defaultTheme.major)
	th.minor = parseFontScheme(fs.Child("minorFont"), defaultTheme.minor)
	fm := te.Child("fmtScheme")
	th.fills = fm.Child("fillStyleLst").Elements()
	th.lines = fm.Child("lnStyleLst").Elements()
	th.effects = fm.Child("effectStyleLst").Elements()
	th.bgFills = fm.Child("bgFillStyleLst").Elements()
	return th
}

func parseFontScheme(n *ooxml.Node, def fontScheme) fontScheme {
	if n == nil {
		return def
	}
	fs := fontScheme{
		latin:   n.Child("latin").AttrStr("typeface", def.latin),
		ea:      n.Child("ea").AttrStr("typeface", ""),
		cs:      n.Child("cs").AttrStr("typeface", ""),
		scripts: map[string]string{},
	}
	for _, f := range n.Children("font") {
		fs.scripts[f.AttrStr("script", "")] = f.AttrStr("typeface", "")
	}
	return fs
}

// fontFor resolves a theme font reference (+mj-lt, +mn-ea …). script is the
// script tag to use when the theme's ea/cs font is empty.
func (th *theme) fontFor(typeface, script string) string {
	if len(typeface) < 4 || typeface[0] != '+' {
		return typeface
	}
	fs := th.minor
	if typeface[1:3] == "mj" {
		fs = th.major
	}
	switch typeface[3:] {
	case "-lt":
		return fs.latin
	case "-ea":
		if fs.ea != "" {
			return fs.ea
		}
		return fs.scripts[script]
	case "-cs":
		if fs.cs != "" {
			return fs.cs
		}
		return fs.scripts[script]
	}
	return typeface
}

// styleEntry returns entry idx (1-based) of a style matrix list.
func styleEntry(list []*ooxml.Node, idx int) *ooxml.Node {
	if idx < 1 || idx > len(list) {
		return nil
	}
	return list[idx-1]
}

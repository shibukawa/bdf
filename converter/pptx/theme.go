package pptx

// theme is the part of a DrawingML theme the renderer needs.
type theme struct {
	colors map[string]rgba // dk1, lt1, dk2, lt2, accent1-6, hlink, folHlink
	major  fontScheme
	minor  fontScheme
	// format scheme: style matrices referenced by p:style (1-based indices)
	fills, lines, effects, bgFills []*node
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

func parseTheme(root *node) *theme {
	if root == nil {
		return defaultTheme
	}
	te := root.child("themeElements")
	th := &theme{colors: map[string]rgba{}}
	for k, v := range defaultTheme.colors {
		th.colors[k] = v
	}
	if cs := te.child("clrScheme"); cs != nil {
		plain := &colorCtx{}
		for _, k := range cs.Kids {
			if c, ok := plain.color(k); ok {
				th.colors[k.Name] = c
			}
		}
	}
	fs := te.child("fontScheme")
	th.major = parseFontScheme(fs.child("majorFont"), defaultTheme.major)
	th.minor = parseFontScheme(fs.child("minorFont"), defaultTheme.minor)
	fm := te.child("fmtScheme")
	th.fills = fm.child("fillStyleLst").kids()
	th.lines = fm.child("lnStyleLst").kids()
	th.effects = fm.child("effectStyleLst").kids()
	th.bgFills = fm.child("bgFillStyleLst").kids()
	return th
}

func parseFontScheme(n *node, def fontScheme) fontScheme {
	if n == nil {
		return def
	}
	fs := fontScheme{
		latin:   n.child("latin").attrStr("typeface", def.latin),
		ea:      n.child("ea").attrStr("typeface", ""),
		cs:      n.child("cs").attrStr("typeface", ""),
		scripts: map[string]string{},
	}
	for _, f := range n.children("font") {
		fs.scripts[f.attrStr("script", "")] = f.attrStr("typeface", "")
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
func styleEntry(list []*node, idx int) *node {
	if idx < 1 || idx > len(list) {
		return nil
	}
	return list[idx-1]
}

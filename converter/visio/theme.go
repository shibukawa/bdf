package visio

import (
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// A dynamic theme ([MS-VSDX] §2.2.7.4) is a DrawingML theme with Visio
// extensions: variant color schemes and style schemes, a connector format
// scheme next to the shape one, and line, fill and font properties that
// DrawingML lacks. Shapes pick a slice of it with their QuickStyle cells:
// a color (dk1, lt1, accent1-6, the background, or a variant color) for
// lines, fills, shadows and text, and an index into the format scheme's
// style lists for lines, fills, effects and fonts.

type theme struct {
	colors    map[string]bdf.Color // clrScheme, with the Visio background as "bkgnd"
	varColors [][]bdf.Color        // variationClrScheme: varColor1-7
	varStyles [][]varStyle         // variationStyleScheme: four styles each
	major     fontScheme
	minor     fontScheme
	shapes    formatScheme // fmtScheme
	conns     formatScheme // fmtConnectorScheme
}

type varStyle struct{ fill, line, effect, font int }

type fontScheme struct {
	latin, ea string
	scripts   map[string]string // script tag (JPAN …) → typeface
}

// formatScheme is the style matrix of shapes or of connectors.
type formatScheme struct {
	fills, lines, effects []*ooxml.Node
	lineEx                []*ooxml.Node // lineStyles: rounding, arrows, pattern
	fillProps             []*ooxml.Node // fillStyles: pattern
	fontProps             []*ooxml.Node // fontStylesGroup: style, color
}

// descendant finds the first element with a local name under n (breadth
// first), for the extensions whose URIs do not matter here.
func descendant(n *ooxml.Node, name string) *ooxml.Node {
	queue := []*ooxml.Node{n}
	for len(queue) > 0 {
		k := queue[0]
		queue = queue[1:]
		if k == nil {
			continue
		}
		for _, c := range k.Kids {
			if c.Name == name {
				return c
			}
		}
		queue = append(queue, k.Kids...)
	}
	return nil
}

func parseTheme(root *ooxml.Node) *theme {
	te := root.Child("themeElements")
	th := &theme{colors: map[string]bdf.Color{}}
	cs := te.Child("clrScheme")
	for _, k := range cs.Elements() {
		if k.Name == "extLst" {
			continue
		}
		if c, ok := drawingml.ResolveColor(k, nil, 0); ok {
			th.colors[k.Name] = c
		}
	}
	if bg := descendant(cs.Child("extLst"), "bkgnd"); bg != nil {
		if c, ok := drawingml.ResolveColor(bg, nil, 0); ok {
			th.colors["bkgnd"] = c
		}
	}
	for _, v := range descendant(cs, "variationClrSchemeLst").Children("variationClrScheme") {
		var list []bdf.Color
		for i := 1; i <= 7; i++ {
			c, _ := drawingml.ResolveColor(v.Child("varColor"+string(rune('0'+i))), th.colors, 0)
			list = append(list, c)
		}
		th.varColors = append(th.varColors, list)
	}
	for _, v := range descendant(te, "variationStyleSchemeLst").Children("variationStyleScheme") {
		var list []varStyle
		for _, s := range v.Children("varStyle") {
			list = append(list, varStyle{attrInt(s, "fillIdx", 0), attrInt(s, "lineIdx", 0), attrInt(s, "effectIdx", 0), attrInt(s, "fontIdx", 0)})
		}
		th.varStyles = append(th.varStyles, list)
	}
	fs := te.Child("fontScheme")
	th.major = parseFontScheme(fs.Child("majorFont"))
	th.minor = parseFontScheme(fs.Child("minorFont"))
	fm := te.Child("fmtScheme")
	th.shapes = formatScheme{fills: fm.Child("fillStyleLst").Elements(), lines: fm.Child("lnStyleLst").Elements(),
		effects: fm.Child("effectStyleLst").Elements()}
	if cf := descendant(te.Child("extLst"), "fmtConnectorScheme"); cf != nil {
		th.conns = formatScheme{fills: cf.Child("fillStyleLst").Elements(), lines: cf.Child("lnStyleLst").Elements(),
			effects: cf.Child("effectStyleLst").Elements()}
	} else {
		th.conns = th.shapes
	}
	lineEx := func(n *ooxml.Node) []*ooxml.Node {
		var out []*ooxml.Node
		for _, s := range n.Children("lineStyle") {
			out = append(out, s.Child("lineEx"))
		}
		return out
	}
	if ls := descendant(te, "lineStyles"); ls != nil {
		th.shapes.lineEx = lineEx(ls.Child("fmtSchemeLineStyles"))
		th.conns.lineEx = lineEx(ls.Child("fmtConnectorSchemeLineStyles"))
	}
	if f := descendant(te, "fillStyles"); f != nil {
		th.shapes.fillProps = f.Children("fillProps")
		th.conns.fillProps = th.shapes.fillProps
	}
	if g := descendant(te, "fontStylesGroup"); g != nil {
		th.shapes.fontProps = g.Child("fontStyles").Children("fontProps")
		th.conns.fontProps = g.Child("connectorFontStyles").Children("fontProps")
	}
	return th
}

func parseFontScheme(n *ooxml.Node) fontScheme {
	fs := fontScheme{latin: n.Child("latin").AttrStr("typeface", "Calibri"), ea: n.Child("ea").AttrStr("typeface", ""),
		scripts: map[string]string{}}
	for _, f := range n.Children("font") {
		fs.scripts[f.AttrStr("script", "")] = f.AttrStr("typeface", "")
	}
	return fs
}

// script returns the typeface of a script ("Jpan"; themes write the tags
// in either case).
func (fs fontScheme) script(tag string) string {
	for _, k := range []string{tag, strings.ToUpper(tag), strings.ToLower(tag)} {
		if v, ok := fs.scripts[k]; ok {
			return v
		}
	}
	return ""
}

// entry returns element k (1-based) of a style list, or nil.
func entry(list []*ooxml.Node, k int) *ooxml.Node {
	if k < 1 || k > len(list) {
		return nil
	}
	return list[k-1]
}

// lineExFor returns the Visio line properties for a line matrix index. The
// list has an entry for index 0 before those of the style list when it is
// one longer.
func (f *formatScheme) lineExFor(k int) *ooxml.Node {
	if len(f.lineEx) > len(f.lines) {
		if k >= 0 && k < len(f.lineEx) {
			return f.lineEx[k]
		}
		return nil
	}
	return entry(f.lineEx, k)
}

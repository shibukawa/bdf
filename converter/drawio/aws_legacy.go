package drawio

import (
	"math"
	"strings"
)

// Older AWS icons: draw.io still opens diagrams drawn with the AWS icon
// sets that came before the current one (mxgraph.aws4): mxgraph.aws (the
// first, line icons), mxgraph.aws2 (2015), mxgraph.aws3 ("AWS17") and
// mxgraph.aws3d (isometric). Their libraries are not embedded; instead of
// drawing them as rectangles, the converter draws each with the current
// icon of the same service or resource, styled the way draw.io's current
// AWS palette styles it. The targets come from aws_legacy_table.go, made by
// tools/gen-drawio-awsmap from draw.io's palettes.
//
// The substitution rewrites the model before the view computes the states,
// so that edges are routed to, and labels placed at, the new icon:
//
//   - the style keeps the cell's own keys (label, font, label position,
//     spacing, rotation, flips, ...), loses the keys that draw the old icon
//     (shape, colors, outline) and gains the target's: its shape and icon,
//     colors and outline, and for groups their container keys and the
//     placement of the label next to the group's icon. A font color the
//     cell sets is kept (the old palettes set none);
//   - icons get the aspect of the new icon: the largest box of the target's
//     aspect centered in the cell's bounds (the old icons are often taller
//     than wide, the current service icons square), with the cell's
//     children kept where they were. Groups keep their bounds.
//
// Old names without a counterpart stay rectangles, with a warning.

//go:generate go run ../../tools/gen-drawio-awsmap -out aws_legacy_table.go -overrides ../../tools/gen-drawio-awsmap/overrides.txt

// awsLegacyTarget is the current style an old AWS shape is drawn with.
type awsLegacyTarget struct {
	// style holds the keys the target sets: its shape, icon and colors,
	// and for groups the container and label placement keys.
	style string
	// w, h are the size of the icon, whose aspect the cell's bounds are
	// fitted to; zero keeps the bounds (groups).
	w, h float64
	// keepColors keeps the cell's colors: a draw.io shape stands in for
	// a shape that is not an icon (the plain frame of old groups).
	keepColors bool
}

// awsLegacyIconKeys are the style keys that draw an old icon; they are
// dropped even when the target does not set them.
var awsLegacyIconKeys = []string{"shape", "resIcon", "prIcon", "grIcon", "grStroke", "fillColor",
	"gradientColor", "gradientDirection", "strokeColor", "strokeColor2", "strokeWidth", "dashed",
	"dashPattern", "fixDash", "aspect", "rounded", "arcSize", "absoluteArcSize", "glass"}

// isAWSLegacyShape reports whether a shape name belongs to an older AWS
// icon set.
func isAWSLegacyShape(name string) bool {
	for _, p := range []string{"mxgraph.aws.", "mxgraph.aws2.", "mxgraph.aws3.", "mxgraph.aws3d."} {
		if len(name) > len(p) && strings.EqualFold(name[:len(p)], p) {
			return true
		}
	}
	return false
}

// lookupAWSLegacy returns the target of an old AWS shape name; ok is false
// for names that are not in the table, and t is nil for names without a
// counterpart. Stencil names are matched without case, as draw.io's
// stencil registry does.
func lookupAWSLegacy(name string) (t *awsLegacyTarget, ok bool) {
	i, ok := awsLegacyShapes[name]
	if !ok {
		i, ok = awsLegacyShapes[strings.ToLower(name)]
	}
	if !ok || i < 0 {
		return nil, ok
	}
	return &awsLegacyTargets[i], true
}

// substituteAWSLegacy rewrites the vertices of m drawn with old AWS shapes
// (see above) and returns how many it rewrote. Names the converter draws
// itself (a registered shape or an embedded stencil) are left alone.
func (c *converter) substituteAWSLegacy(m *model) int {
	n := 0
	for _, cl := range m.ordered {
		if !cl.vertex || !strings.Contains(cl.styleStr, "aws") {
			continue
		}
		name := parseStyle(cl.styleStr, false)["shape"]
		if !isAWSLegacyShape(name) {
			continue
		}
		if _, ok := shapeRegistry[name]; ok || c.stencil(name) != nil {
			continue
		}
		t, ok := lookupAWSLegacy(name)
		if !ok {
			continue // not a known name: newShape reports it
		}
		if t == nil {
			// the same key as newShape's warning, which it replaces
			c.warnOnce("shape:"+name, "shape %q of an older AWS icon set has no current AWS counterpart; drawn as a rectangle", name)
			continue
		}
		cl.styleStr = rewriteAWSLegacyStyle(cl.styleStr, t)
		if t.w > 0 && t.h > 0 && cl.geo != nil {
			w, h := t.w, t.h
			if d := parseStyle(cl.styleStr, false)["direction"]; d == "north" || d == "south" {
				w, h = h, w
			}
			fitAWSLegacyGeometry(cl, w, h)
		}
		n++
	}
	return n
}

// reportAWSLegacy warns once per conversion about the shapes of older AWS
// icon sets drawn with their current counterparts.
func (c *converter) reportAWSLegacy() {
	switch c.awsLegacy {
	case 0:
	case 1:
		c.warnf("1 shape of an older AWS icon set is drawn with its current AWS counterpart")
	default:
		c.warnf("%d shapes of older AWS icon sets are drawn with their current AWS counterparts", c.awsLegacy)
	}
}

// rewriteAWSLegacyStyle returns a cell's style with the target's keys in
// place of those that drew the old shape.
func rewriteAWSLegacyStyle(style string, t *awsLegacyTarget) string {
	toks := strings.Split(style, ";")
	ownFont := false // a font color of the cell's own (the old palettes set none)
	for _, tok := range toks {
		if v, ok := strings.CutPrefix(tok, "fontColor="); ok && v != "" {
			ownFont = true
		}
	}
	drop := map[string]bool{"shape": true}
	if !t.keepColors {
		for _, k := range awsLegacyIconKeys {
			drop[k] = true
		}
	}
	var add []string
	for _, kv := range strings.Split(t.style, ";") {
		k, _, _ := strings.Cut(kv, "=")
		if k == "fontColor" && ownFont {
			continue
		}
		drop[k] = true
		add = append(add, kv)
	}
	var out []string
	for i, tok := range toks {
		k, _, isKV := strings.Cut(tok, "=")
		// empty entries go, but a leading one (no default style) stays
		if tok == "" && (i > 0 || style == "") || isKV && drop[k] {
			continue
		}
		out = append(out, tok)
	}
	return strings.Join(append(out, add...), ";")
}

// fitAWSLegacyGeometry shrinks a vertex's bounds to the largest box of
// aspect w:h centered in them, leaving its children where they are.
func fitAWSLegacyGeometry(cl *cell, w, h float64) {
	g := cl.geo
	if g.w <= 0 || g.h <= 0 {
		return
	}
	s := math.Min(g.w/w, g.h/h)
	nw, nh := w*s, h*s
	dx, dy := (g.w-nw)/2, (g.h-nh)/2
	if math.Abs(dx) < 1e-9 && math.Abs(dy) < 1e-9 {
		return
	}
	ow, oh := g.w, g.h
	if g.relative {
		// the origin is the relative position plus the offset
		off := point{dx, dy}
		if g.offset != nil {
			off.x += g.offset.x
			off.y += g.offset.y
		}
		g.offset = &off
	} else {
		g.x += dx
		g.y += dy
	}
	g.w, g.h = nw, nh
	// children are placed from the origin (and relative ones from the
	// size), which moved
	for _, k := range cl.children {
		kg := k.geo
		if kg == nil {
			continue
		}
		switch {
		case k.edge:
			for i := range kg.points {
				kg.points[i].x -= dx
				kg.points[i].y -= dy
			}
			for _, p := range []*point{kg.sourcePoint, kg.targetPoint} {
				if p != nil {
					p.x -= dx
					p.y -= dy
				}
			}
		case kg.relative:
			off := point{kg.x*(ow-nw) - dx, kg.y*(oh-nh) - dy}
			if kg.offset != nil {
				off.x += kg.offset.x
				off.y += kg.offset.y
			}
			kg.offset = &off
		default:
			kg.x -= dx
			kg.y -= dy
		}
	}
}

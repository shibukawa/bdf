package drawio

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
)

// The shapes Shapes.js adds for diagrams imported from mermaid: git graph
// tags and commits, mindmap bangs, the Ishikawa fish head, the "odd"
// node, block arrows and sankey links.

func init() {
	registerShape("gitTag", gitTagShape)
	registerShape("gitMergeCommit", gitMergeCommitShape)
	registerShape("gitCherryPick", gitCherryPickShape)
	registerShape("mindmapBang", mindmapBangShape)
	registerShape("ishikawaHead", ishikawaHeadShape)
	registerShape("mermaidOdd", mermaidOddShape)
	registerShape("mermaidBlockArrow", mermaidBlockArrowShape)
	registerShape("mermaidSankeyLink", &shapeDef{
		paintEdge:  paintSankeyLink,
		noRotation: true,
		noInvert:   true,
		augmentBounds: func(s *shape, r rect) rect {
			// the band ends at the terminals: it only grows vertically
			b := pointsBounds(s.points)
			g := sankeyWidth(s)/2 + s.strokewidth
			return r.union(rect{b.x, b.y - g, b.w, b.h + 2*g})
		},
	})
}

// gitTagShape ports GitTagShape: a tag with a pointed tab on the left
// (tabSize, tabInset) and a hole (holeSize, holeColor).
var gitTagShape = func() *shapeDef {
	path := func(s *shape, c *c2d, x, y, w, h float64) {
		tabSize := clamp(s.style.num("tabSize", 8), 0, w)
		tabInset := clamp(s.style.num("tabInset", 4), 0, h)
		tabY1 := (h - tabInset) / 2
		tabY2 := tabY1 + tabInset
		c.moveTo(0, tabY1)
		c.lineTo(0, tabY2)
		c.lineTo(tabSize, h)
		c.lineTo(w, h)
		c.lineTo(w, 0)
		c.lineTo(tabSize, 0)
		c.close()
	}
	actor := actorShape(path)
	return &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			actor.paintVertex(s, c, x, y, w, h)
			holeSize := math.Max(0, s.style.num("holeSize", 1))
			if holeSize <= 0 {
				return
			}
			tabSize := s.style.num("tabSize", 8)
			holeColor := s.style.get("holeColor", s.style.get("fontColor", "#333333"))
			c.setFillColor(holeColor)
			c.setStrokeColor(holeColor)
			c.begin()
			c.ellipse(tabSize/2-holeSize, h/2-holeSize, holeSize*2, holeSize*2)
			c.fillAndStroke()
		},
		labelBounds: func(s *shape, r rect) rect {
			tabSize := s.style.num("tabSize", 8)
			return rect{r.x + tabSize, r.y, r.w - tabSize, r.h}
		},
	}
}()

// gitMergeCommitShape ports GitMergeCommitShape: a circle with an inner
// circle (innerColor).
var gitMergeCommitShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	innerColor := s.style.get("innerColor", "#ECECFF")
	c.translate(x, y)
	c.ellipse(0, 0, w, h)
	c.fillAndStroke()
	innerD := math.Min(w, h) * 0.6
	c.setFillColor(innerColor)
	c.setStrokeColor(innerColor)
	c.ellipse((w-innerD)/2, (h-innerD)/2, innerD, innerD)
	c.fillAndStroke()
}}

// gitCherryPickShape ports GitCherryPickShape: a circle with a cherry
// (featureColor).
var gitCherryPickShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	featureColor := s.style.get("featureColor", "#fff")
	c.translate(x, y)
	cx, cy := w/2, h/2
	sc := math.Min(w, h) / 20
	c.ellipse(0, 0, w, h)
	c.fillAndStroke()
	c.setFillColor(featureColor)
	c.setStrokeColor(featureColor)
	c.setStrokeWidth(0)
	eyeR := 2.75 * sc
	c.ellipse(cx-3*sc-eyeR, cy+2*sc-eyeR, eyeR*2, eyeR*2)
	c.fillAndStroke()
	c.ellipse(cx+3*sc-eyeR, cy+2*sc-eyeR, eyeR*2, eyeR*2)
	c.fillAndStroke()
	c.setStrokeWidth(1 * sc)
	c.begin()
	c.moveTo(cx+3*sc, cy+1*sc)
	c.lineTo(cx, cy-5*sc)
	c.stroke()
	c.begin()
	c.moveTo(cx-3*sc, cy+1*sc)
	c.lineTo(cx, cy-5*sc)
	c.stroke()
}}

// mindmapBangShape ports MindmapBangShape: a starburst of arcs.
var mindmapBangShape = func() *shapeDef {
	d := actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		// the spikes extend beyond the design rectangle, 80% of the bounds
		W, H := w*0.8, h*0.8
		r := W * 0.15
		r80 := r * 0.8
		px, py := W*0.10, H*0.10
		c.moveTo(px, py)
		arc := func(r, dx, dy float64) {
			px += dx
			py += dy
			c.arcTo(r, r, 0, false, false, px, py)
		}
		// top: corners spike up
		arc(r, W*0.25, -H*0.10)
		arc(r, W*0.25, 0)
		arc(r, W*0.25, 0)
		arc(r, W*0.25, H*0.10)
		// right: the middle spike protrudes
		arc(r, W*0.15, H*0.33)
		arc(r80, 0, H*0.34)
		arc(r, -W*0.15, H*0.33)
		// bottom: corners spike down
		arc(r, -W*0.25, H*0.15)
		arc(r, -W*0.25, 0)
		arc(r, -W*0.25, 0)
		arc(r, -W*0.25, -H*0.15)
		// left
		arc(r, -W*0.10, -H*0.33)
		arc(r80, 0, -H*0.34)
		arc(r, W*0.10, -H*0.33)
		c.close()
	})
	d.labelBounds = func(s *shape, r rect) rect {
		ix, iy := r.w*0.10, r.h*0.10
		return rect{r.x + ix, r.y + iy, r.w - ix*2, r.h - iy*2}
	}
	return d
}()

// ishikawaHeadShape ports IshikawaHeadShape: a teardrop with a flat left side.
var ishikawaHeadShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	c.moveTo(0, 0)
	c.lineTo(0, h)
	c.quadTo(2*w, h/2, 0, 0)
	c.close()
})

// mermaidOddShape ports OddShape: a rectangle with a notch on the left.
var mermaidOddShape = func() *shapeDef {
	d := actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		notch := h / 4
		c.moveTo(0, 0)
		c.lineTo(notch, h/2)
		c.lineTo(0, h)
		c.lineTo(w, h)
		c.lineTo(w, 0)
		c.close()
	})
	d.labelBounds = func(s *shape, r rect) rect {
		notch := r.h / 4
		return rect{r.x + notch, r.y, r.w - notch, r.h}
	}
	return d
}()

// mermaidBlockArrowShape ports MermaidBlockArrowShape: the block arrows of
// mermaid's block diagrams, pointing to the directions listed in dirs
// (right, left, up, down, x, y).
var mermaidBlockArrowShape = func() *shapeDef {
	actor := actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		dirs := map[string]bool{}
		for _, d := range strings.FieldsFunc(s.style.get("dirs", "right"), func(r rune) bool { return r == ',' || r == '|' || r == ' ' }) {
			switch d = strings.ToLower(strings.TrimSpace(d)); d {
			case "x":
				dirs["right"], dirs["left"] = true, true
			case "y":
				dirs["up"], dirs["down"] = true, true
			default:
				dirs[d] = true
			}
		}
		pad := s.style.num("nodePadding", 8) / 2
		mid := h / 2
		right, left, up, down := dirs["right"], dirs["left"], dirs["up"], dirs["down"]
		var pts [][2]float64
		switch {
		case right && left && up && down:
			pts = [][2]float64{{0, 0}, {mid, 0}, {w / 2, 2 * pad}, {w - mid, 0}, {w, 0},
				{w, -h / 3}, {w + 2*pad, -h / 2}, {w, -2 * h / 3}, {w, -h},
				{w - mid, -h}, {w / 2, -h - 2*pad}, {mid, -h},
				{0, -h}, {0, -2 * h / 3}, {-2 * pad, -h / 2}, {0, -h / 3}}
		case right && left && up:
			pts = [][2]float64{{mid, 0}, {w - mid, 0}, {w, -h / 2}, {w - mid, -h}, {mid, -h}, {0, -h / 2}}
		case right && left && down:
			pts = [][2]float64{{0, 0}, {mid, -h}, {w - mid, -h}, {w, 0}}
		case right && up && down:
			pts = [][2]float64{{0, 0}, {w, -mid}, {w, -h + mid}, {0, -h}}
		case left && up && down:
			pts = [][2]float64{{w, 0}, {0, -mid}, {0, -h + mid}, {w, -h}}
		case right && left:
			pts = [][2]float64{{mid, 0}, {mid, -pad}, {w - mid, -pad}, {w - mid, 0}, {w, -h / 2},
				{w - mid, -h}, {w - mid, -h + pad}, {mid, -h + pad}, {mid, -h}, {0, -h / 2}}
		case up && down:
			pts = [][2]float64{{w / 2, 0}, {0, -pad}, {mid, -pad}, {mid, -h + pad}, {0, -h + pad},
				{w / 2, -h}, {w, -h + pad}, {w - mid, -h + pad}, {w - mid, -pad}, {w, -pad}}
		case right && up:
			pts = [][2]float64{{0, 0}, {w, -mid}, {0, -h}}
		case right && down:
			pts = [][2]float64{{0, 0}, {w, 0}, {0, -h}}
		case left && up:
			pts = [][2]float64{{w, 0}, {0, -mid}, {w, -h}}
		case left && down:
			pts = [][2]float64{{w, 0}, {0, 0}, {w, -h}}
		case right:
			pts = [][2]float64{{mid, -pad}, {w - mid, -pad}, {w - mid, 0}, {w, -h / 2},
				{w - mid, -h}, {w - mid, -h + pad}, {mid, -h + pad}}
		case left:
			pts = [][2]float64{{mid, 0}, {mid, -pad}, {w - mid, -pad}, {w - mid, -h + pad},
				{mid, -h + pad}, {mid, -h}, {0, -h / 2}}
		case up:
			pts = [][2]float64{{mid, -pad}, {mid, -h + pad}, {0, -h + pad}, {w / 2, -h},
				{w, -h + pad}, {w - mid, -h + pad}, {w - mid, -pad}}
		case down:
			pts = [][2]float64{{w / 2, 0}, {0, -pad}, {mid, -pad}, {mid, -h + pad},
				{w - mid, -h + pad}, {w - mid, -pad}, {w, -pad}}
		default:
			pts = [][2]float64{{0, 0}, {w, 0}, {w, -h}, {0, -h}}
		}
		c.moveTo(pts[0][0], h+pts[0][1])
		for _, p := range pts[1:] {
			c.lineTo(p[0], h+p[1])
		}
		c.close()
	})
	return &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.setLineJoin("round")
		actor.paintVertex(s, c, x, y, w, h)
	}}
}()

// sankeyWidth is MermaidSankeyLinkShape.getEdgeWidth.
func sankeyWidth(s *shape) float64 { return math.Max(1, s.style.num("width", 10)) }

// paintSankeyLink paints the band of a sankey link: d3-sankey's cubic
// from the first to the last point, offset by half the width on both
// sides and filled (MermaidSankeyLinkShape.paintEdgeShape). The band is
// composited with multiply, as mermaid's are, so that crossings darken.
func paintSankeyLink(s *shape, c *c2d, pts []point) {
	p0, pe := pts[0], pts[len(pts)-1]
	x0, y0, x1, y1 := p0.x, p0.y, pe.x, pe.y
	mx := (x0 + x1) / 2
	hw := sankeyWidth(s) / 2
	const n = 32
	var side1, side2 [n + 1]point
	for i := 0; i <= n; i++ {
		t := float64(i) / n
		u := 1 - t
		x := u*u*u*x0 + 3*u*u*t*mx + 3*u*t*t*mx + t*t*t*x1
		y := u*u*u*y0 + 3*u*u*t*y0 + 3*u*t*t*y1 + t*t*t*y1
		dx := 3*u*u*(mx-x0) + 3*t*t*(x1-mx)
		dy := 6 * u * t * (y1 - y0)
		l := math.Sqrt(dx*dx + dy*dy)
		if l == 0 {
			dx, dy, l = 1, 0, 1
		}
		nx, ny := -dy/l*hw, dx/l*hw
		side1[i] = point{x + nx, y + ny}
		side2[i] = point{x - nx, y - ny}
	}
	c.begin()
	c.moveTo(side1[0].x, side1[0].y)
	for i := 1; i <= n; i++ {
		c.lineTo(side1[i].x, side1[i].y)
	}
	for i := n; i >= 0; i-- {
		c.lineTo(side2[i].x, side2[i].y)
	}
	c.close()
	// the SVG canvas sets mix-blend-mode on the band's element
	b := c.pathBox.grow(s.strokewidth*c.st.scale + 1)
	c.save()
	c.obj.GroupBegin(1, bdf.BlendMultiply, f32(b.x), f32(b.y), f32(b.w), f32(b.h))
	c.fillAndStroke()
	c.obj.GroupEnd()
	c.restore()
}

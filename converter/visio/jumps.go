package visio

import (
	"math"
	"sort"

	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// Where dynamic connectors cross, Visio draws a line jump on one of them:
// an arc, a gap, a square or a polygon of 2 to 7 sides. The jumps are not
// saved; Visio places them when it draws the page, and so are they here,
// before the shapes are drawn. The page decides which of two crossing
// connectors jumps (LineJumpCode: the horizontal one by default, the
// vertical one, or the one above or below in the drawing order), how
// jumps look (LineJumpStyle) and how large they are (LineJumpFactorX
// times LineToLineX, and the same in Y for vertical connectors); a
// connector can always jump, never jump, leave the jump to the other
// connector or have no jumps at its crossings (ConLineJumpCode), and set
// its own style and direction.

// jump is a line jump to insert into a connector.
type jump struct {
	at    pt      // the crossing, page points
	width float64 // points, along the connector
	style int     // 1 arc, 2 gap, 3 square, 4-9 polygons of 2-7 sides
	up    bool    // for horizontal connectors: bulge upwards (else down)
	left  bool    // for vertical connectors: bulge to the left (else right)
}

// connPath is the stroked path of a connector, as straight pieces.
type connPath struct {
	s     *shape
	z     int // drawing order
	segs  []connSeg
	code  int // ConLineJumpCode
	style int // ConLineJumpStyle
	dirX  int // ConLineJumpDirX
	dirY  int // ConLineJumpDirY
}

type connSeg struct {
	a, b     pt
	straight bool // a line of the geometry (not a piece of a curve)
}

// routable reports whether a shape is a dynamic connector: a routable
// 1-D shape (ObjType).
func (p *pageCtx) routable(s *shape) bool {
	t := int(p.num(s, "ObjType", 0))
	return t&2 != 0 && t&4 == 0
}

// collectJumps finds the line jumps of the connectors of the page.
func (p *pageCtx) collectJumps() {
	code := p.pageInt("LineJumpCode", 1)
	var conns []*connPath
	z := 0
	var walk func(list []*shape, parent canvas.Matrix, depth int)
	walk = func(list []*shape, parent canvas.Matrix, depth int) {
		for _, s := range list {
			if s.del || depth > maxChain || s.typ == "Guide" {
				continue
			}
			m := parent.Mul(p.localMatrix(s))
			z++
			if len(s.kids) > 0 {
				if int(p.num(s, "DisplayMode", 2)) != 0 {
					walk(s.kids, m, depth+1)
				}
				continue
			}
			if p.layerHidden(s) || !p.routable(s) {
				continue
			}
			if c := p.connectorPath(s, m, z); len(c.segs) > 0 {
				conns = append(conns, c)
			}
		}
	}
	walk(p.pg.shapes, p.pageMatrix(), 0)
	for i, a := range conns {
		for _, b := range conns[i+1:] {
			p.crossings(a, b, code)
		}
	}
}

// connectorPath collects the stroked pieces of a connector.
func (p *pageCtx) connectorPath(s *shape, m canvas.Matrix, z int) *connPath {
	c := &connPath{s: s, z: z, code: int(p.num(s, "ConLineJumpCode", 0)), style: int(p.num(s, "ConLineJumpStyle", 0)),
		dirX: int(p.num(s, "ConLineJumpDirX", 0)), dirY: int(p.num(s, "ConLineJumpDirY", 0))}
	if p.lineStyle(s) == nil {
		return c
	}
	w, h := p.num(s, "Width", 0), p.num(s, "Height", 0)
	for _, g := range s.geometry() {
		if g.flag("NoShow") || g.flag("NoLine") {
			continue
		}
		for _, sp := range buildGeometry(g, m, w, h) {
			prev := sp.start
			for _, sg := range sp.segs {
				if !sg.cubic {
					c.segs = append(c.segs, connSeg{prev, sg.p, true})
				} else {
					// curves only take part as what others jump over
					const n = 8
					q := prev
					for i := 1; i <= n; i++ {
						r := cubicAt(prev, sg.c1, sg.c2, sg.p, float64(i)/n)
						c.segs = append(c.segs, connSeg{q, r, false})
						q = r
					}
				}
				prev = sg.p
			}
		}
	}
	return c
}

func cubicAt(p0, p1, p2, p3 pt, t float64) pt {
	u := 1 - t
	return p0.mul(u * u * u).add(p1.mul(3 * u * u * t)).add(p2.mul(3 * u * t * t)).add(p3.mul(t * t * t))
}

// pageInt reads an integer cell of the page sheet.
func (p *pageCtx) pageInt(name string, def int) int {
	c, ok := p.pg.sheet.cells[name]
	if !ok {
		return def
	}
	return atoi(c.v, def)
}

func (p *pageCtx) pageNum(name string, def float64) float64 {
	c, ok := p.pg.sheet.cells[name]
	if !ok || c.v == "" {
		return def
	}
	return num(c.v)
}

// crossings records the jumps where two connectors cross.
func (p *pageCtx) crossings(a, b *connPath, code int) {
	for _, sa := range a.segs {
		for _, sb := range b.segs {
			at, ok := segmentCrossing(sa.a, sa.b, sb.a, sb.b)
			if !ok {
				continue
			}
			j := jumper(a, b, sa, sb, code)
			if j == nil {
				continue
			}
			seg := sa
			if j == b {
				seg = sb
			}
			if !seg.straight {
				continue
			}
			jp := p.jumpAt(j, seg, at)
			// the jump must fit between the ends of the piece
			if at.dist(seg.a) < jp.width/2 || at.dist(seg.b) < jp.width/2 {
				continue
			}
			p.jumps[j.s] = append(p.jumps[j.s], jp)
		}
	}
}

// segmentCrossing returns where two segments cross (not where they only
// touch at their ends or overlap).
func segmentCrossing(a0, a1, b0, b1 pt) (pt, bool) {
	d1, d2 := a1.sub(a0), b1.sub(b0)
	den := d1.x*d2.y - d1.y*d2.x
	if math.Abs(den) < 1e-9*d1.len()*d2.len() || d1.len() == 0 || d2.len() == 0 {
		return pt{}, false
	}
	e := b0.sub(a0)
	t := (e.x*d2.y - e.y*d2.x) / den
	u := (e.x*d1.y - e.y*d1.x) / den
	const eps = 1e-6
	if t <= eps || t >= 1-eps || u <= eps || u >= 1-eps {
		return pt{}, false
	}
	return a0.add(d1.mul(t)), true
}

// horizontality is how close a piece is to horizontal (1) rather than
// vertical (0).
func horizontality(s connSeg) float64 {
	d := s.b.sub(s.a)
	l := d.len()
	if l == 0 {
		return 0
	}
	return math.Abs(d.x) / l
}

// jumper decides which of two crossing connectors jumps (nil: neither).
func jumper(a, b *connPath, sa, sb connSeg, code int) *connPath {
	if a.code == 4 || b.code == 4 {
		return nil
	}
	var j *connPath
	switch {
	case a.code == 2 && b.code != 2, b.code == 3 && a.code != 3:
		j = a
	case b.code == 2 && a.code != 2, a.code == 3 && b.code != 3:
		j = b
	default:
		// the page decides
		ha, hb := horizontality(sa), horizontality(sb)
		switch code {
		case 1: // horizontal lines
			if ha > hb && ha >= math.Sqrt2/2 {
				j = a
			} else if hb > ha && hb >= math.Sqrt2/2 {
				j = b
			}
		case 2: // vertical lines
			if ha < hb && ha <= math.Sqrt2/2 {
				j = a
			} else if hb < ha && hb <= math.Sqrt2/2 {
				j = b
			}
		case 3, 4: // the last routed (taken as the last drawn) or last drawn line
			j = a
			if b.z > a.z {
				j = b
			}
		case 5: // the first drawn line
			j = a
			if b.z < a.z {
				j = b
			}
		}
	}
	if j != nil && j.code == 1 {
		return nil // it never jumps
	}
	return j
}

// jumpAt describes the jump a connector makes at a crossing on one of its
// pieces.
func (p *pageCtx) jumpAt(c *connPath, seg connSeg, at pt) jump {
	horiz := horizontality(seg) >= math.Sqrt2/2
	factor, spacing := p.pageNum("LineJumpFactorX", 2.0/3), p.pageNum("LineToLineX", 0.125)
	if !horiz {
		factor, spacing = p.pageNum("LineJumpFactorY", 2.0/3), p.pageNum("LineToLineY", 0.125)
	}
	style := c.style
	if style == 0 {
		style = p.pageInt("LineJumpStyle", 0)
	}
	if style < 1 || style > 9 {
		style = 1
	}
	dirX := c.dirX
	if dirX == 0 {
		dirX = p.pageInt("PageLineJumpDirX", 0)
	}
	dirY := c.dirY
	if dirY == 0 {
		dirY = p.pageInt("PageLineJumpDirY", 0)
	}
	return jump{at: at, width: factor * spacing * 72 * p.scale, style: style, up: dirX != 2, left: dirY != 2}
}

// applyJumps inserts the jumps of a connector into its stroked figures;
// gaps split figures.
func applyJumps(subs []*subpath, jumps []jump) []*subpath {
	var out []*subpath
	for _, sp := range subs {
		cur := &subpath{start: sp.start}
		prev := sp.start
		for _, sg := range sp.segs {
			if sg.cubic {
				cur.segs = append(cur.segs, sg)
				prev = sg.p
				continue
			}
			for _, iv := range jumpsOn(prev, sg.p, jumps) {
				cur.segs = append(cur.segs, seg{p: iv.a})
				if iv.j.style == 2 {
					out = append(out, cur)
					cur = &subpath{start: iv.b}
					continue
				}
				bump(cur, iv)
			}
			cur.segs = append(cur.segs, sg)
			prev = sg.p
		}
		out = append(out, cur)
	}
	return out
}

// interval is the part of a line a jump (or merged jumps) replaces.
type interval struct {
	a, b pt
	j    jump
	u, n pt // along the line, and the side the jump bulges to
}

// jumpsOn returns the jumps on the line from a to b in order, overlapping
// ones merged.
func jumpsOn(a, b pt, jumps []jump) []interval {
	d := b.sub(a)
	l := d.len()
	if l == 0 {
		return nil
	}
	u := d.mul(1 / l)
	type hit struct {
		t0, t1 float64
		j      jump
	}
	var hits []hit
	for _, j := range jumps {
		e := j.at.sub(a)
		t := e.x*u.x + e.y*u.y
		if t <= 0 || t >= l || math.Abs(e.x*u.y-e.y*u.x) > 0.01 {
			continue
		}
		hits = append(hits, hit{t - j.width/2, t + j.width/2, j})
	}
	sort.Slice(hits, func(i, k int) bool { return hits[i].t0 < hits[k].t0 })
	// the side: up or down for horizontal lines, left or right for vertical ones
	n := pt{u.y, -u.x}
	var out []interval
	for i := 0; i < len(hits); {
		h := hits[i]
		for i++; i < len(hits) && hits[i].t0 <= h.t1; i++ {
			h.t1 = math.Max(h.t1, hits[i].t1)
		}
		nn := n
		if math.Abs(u.x) >= math.Sqrt2/2 {
			if (nn.y < 0) != h.j.up {
				nn = nn.mul(-1)
			}
		} else if (nn.x < 0) != h.j.left {
			nn = nn.mul(-1)
		}
		out = append(out, interval{a: a.add(u.mul(math.Max(h.t0, 0))), b: a.add(u.mul(math.Min(h.t1, l))), j: h.j, u: u, n: nn})
	}
	return out
}

// bump draws the jump over an interval (from its start, where the figure
// is, to its end).
func bump(sp *subpath, iv interval) {
	w := iv.a.dist(iv.b)
	h := iv.j.width / 2
	c := iv.a.add(iv.b).mul(0.5)
	switch style := iv.j.style; {
	case style == 3:
		// a square
		sp.segs = append(sp.segs, seg{p: iv.a.add(iv.n.mul(h))}, seg{p: iv.b.add(iv.n.mul(h))}, seg{p: iv.b})
	case style >= 4:
		// a polygon of 2-7 sides on the arc
		k := style - 2
		for i := 1; i <= k; i++ {
			a := math.Pi * float64(i) / float64(k)
			sp.segs = append(sp.segs, seg{p: c.sub(iv.u.mul(w / 2 * math.Cos(a))).add(iv.n.mul(h * math.Sin(a)))})
		}
	default:
		// an arc: half an ellipse as wide as the jumps it holds
		const k = 0.5522847498
		top := c.add(iv.n.mul(h))
		sp.segs = append(sp.segs,
			seg{cubic: true, c1: iv.a.add(iv.n.mul(h * k)), c2: top.sub(iv.u.mul(w / 2 * k)), p: top},
			seg{cubic: true, c1: top.add(iv.u.mul(w / 2 * k)), c2: iv.b.add(iv.n.mul(h * k)), p: iv.b})
	}
}

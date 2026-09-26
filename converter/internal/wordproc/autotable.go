package wordproc

import "math"

// autoColumns chooses the widths of a table's columns from their content,
// as CSS's automatic table layout does: a column is at least as wide as
// the widest content of its cells that cannot break (min) and needs no
// more than their widest line (max); what the table has between the two
// is shared in proportion to the difference.
type autoColumns struct {
	min, max []float64 // per grid column, the cells' padding included
	fill     bool      // the table takes the whole width
	width    float64   // the width the table asks for (0: none)
}

// widths returns the column widths in a column of width avail.
func (a *autoColumns) widths(avail float64) []float64 {
	n := len(a.min)
	out := make([]float64, n)
	if n == 0 {
		return out
	}
	var smin, smax float64
	for i := range n {
		smin += a.min[i]
		smax += a.max[i]
	}
	target := avail
	if a.width > 0 {
		target = math.Min(target, a.width)
	}
	switch {
	case smax <= target && !a.fill && a.width == 0:
		copy(out, a.max)
	case smax <= target:
		// wider than the content: the columns grow with their content
		for i := range out {
			share := 1 / float64(n)
			if smax > 0 {
				share = a.max[i] / smax
			}
			out[i] = a.max[i] + (target-smax)*share
		}
	case smin >= target:
		// narrower than the content can go: the minimum widths shrink,
		// and the lines break inside words
		for i := range out {
			share := 1 / float64(n)
			if smin > 0 {
				share = a.min[i] / smin
			}
			out[i] = target * share
		}
	default:
		f := (target - smin) / (smax - smin)
		for i := range out {
			out[i] = a.min[i] + (a.max[i]-a.min[i])*f
		}
	}
	return out
}

// autoLayout measures the content of a table's cells for autoColumns.
func autoLayout(t *table) *autoColumns {
	n := len(t.grid)
	a := &autoColumns{min: make([]float64, n), max: make([]float64, n)}
	type spanning struct {
		col, span int
		min, max  float64
	}
	var spans []spanning
	for _, r := range t.rows {
		for _, ce := range r.cells {
			if ce.vmerge == "continue" || ce.col >= n {
				continue
			}
			mn, mx := blockWidths(ce.blocks)
			pad := ce.mar[1] + ce.mar[3]
			mn, mx = mn+pad, math.Max(mx, mn)+pad
			span := min(ce.span, n-ce.col)
			if span <= 1 {
				a.min[ce.col] = math.Max(a.min[ce.col], mn)
				a.max[ce.col] = math.Max(a.max[ce.col], mx)
				continue
			}
			spans = append(spans, spanning{ce.col, span, mn, mx})
		}
	}
	// cells that span columns widen them evenly where they need more
	for _, s := range spans {
		var have, haveMax float64
		for i := s.col; i < s.col+s.span; i++ {
			have += a.min[i]
			haveMax += a.max[i]
		}
		for i := s.col; i < s.col+s.span; i++ {
			if s.min > have {
				a.min[i] += (s.min - have) / float64(s.span)
			}
			if s.max > haveMax {
				a.max[i] += (s.max - haveMax) / float64(s.span)
			}
		}
	}
	for i := range a.max {
		a.max[i] = math.Max(a.max[i], a.min[i])
	}
	return a
}

// blockWidths returns the width of the widest content of blocks that
// cannot break (min) and of their widest line (max).
func blockWidths(bs []block) (mn, mx float64) {
	for _, b := range bs {
		var a, z float64
		switch b := b.(type) {
		case *para:
			a, z = paraWidths(b)
		case *table:
			if b.auto != nil {
				for i := range b.auto.min {
					a += b.auto.min[i]
					z += b.auto.max[i]
				}
				if b.auto.width > 0 {
					a, z = math.Max(a, b.auto.width), math.Max(z, b.auto.width)
				}
			} else {
				for _, g := range b.grid {
					a += g
				}
				z = a
			}
			a, z = a+b.ind, z+b.ind
		}
		mn, mx = math.Max(mn, a), math.Max(mx, z)
	}
	return
}

// paraWidths measures a paragraph for blockWidths. Pictures that shrink
// to the line (inlineObj.fit) need no width.
func paraWidths(p *para) (mn, mx float64) {
	indL, indR, first := p.indents()
	start := indL + first
	if p.label != nil {
		start = indL
	}
	x, word := start, 0.0
	for _, it := range p.items {
		switch it.kind {
		case kBreak, kPage, kColumn:
			mx = math.Max(mx, x+indR)
			x, word = indL, 0
			continue
		}
		w := it.w + it.gap
		x += w
		switch {
		case it.kind == kChar && lbSpace(it.r):
		case it.kind == kObject && it.obj.fit:
		default:
			word += w
			mn = math.Max(mn, indL+word+indR)
		}
		if it.brk || it.kind == kChar && lbSpace(it.r) {
			word = 0
		}
	}
	mx = math.Max(mx, x+indR)
	return mn, math.Max(mx, mn)
}

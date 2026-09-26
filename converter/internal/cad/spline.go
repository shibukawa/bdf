package cad

import "math"

// maxDegree is the highest degree of B-spline evaluated.
const maxDegree = 25

// BSpline adds a B-spline curve (a NURBS curve when weights are given) as
// a new subpath. Non-rational curves of degree 3 or less with clamped knot
// vectors become their exact Bézier segments; other curves are sampled.
// Missing or inconsistent knots are replaced by a clamped uniform vector.
// Curves of a degree above maxDegree are drawn as their control polygon.
func (p *Path) BSpline(degree int, knots []float64, ctrl []Point, weights []float64) *Path {
	n := len(ctrl)
	if n == 0 {
		return p
	}
	if degree < 1 {
		degree = 1
	}
	if degree >= n {
		degree = n - 1
	}
	if degree == 0 {
		return p.MoveTo(ctrl[0].X, ctrl[0].Y)
	}
	if degree > maxDegree {
		// evaluating costs the square of the degree: no drawing needs so
		// high a degree, so draw the control polygon
		return p.Polyline(ctrl, false)
	}
	if len(knots) != n+degree+1 || !nonDecreasing(knots) || knots[degree] >= knots[n] {
		knots = clampedKnots(n, degree)
	}
	if len(weights) != n {
		weights = nil
	}
	rational := false
	for _, w := range weights {
		if w != weights[0] || w <= 0 {
			rational = true
			break
		}
	}
	p.open = false
	if !rational && degree <= 3 && clamped(knots, degree) {
		for i, seg := range bezierSegments(degree, knots, ctrl) {
			if i == 0 {
				p.MoveTo(seg[0].X, seg[0].Y)
			}
			switch degree {
			case 1:
				p.LineTo(seg[1].X, seg[1].Y)
			case 2:
				// elevate the quadratic segment to a cubic one
				c1 := seg[0].Add(seg[1].Sub(seg[0]).Mul(2.0 / 3))
				c2 := seg[2].Add(seg[1].Sub(seg[2]).Mul(2.0 / 3))
				p.CubicTo(c1, c2, seg[2])
			case 3:
				p.CubicTo(seg[1], seg[2], seg[3])
			}
		}
		return p
	}
	// sample every knot span
	t0, t1 := knots[degree], knots[n]
	const perSpan = 16
	spans := 0
	for i := degree; i < n; i++ {
		if knots[i+1] > knots[i] {
			spans++
		}
	}
	steps := max(spans*perSpan, 8)
	steps = min(steps, 4096)
	for i := 0; i <= steps; i++ {
		t := t0 + (t1-t0)*float64(i)/float64(steps)
		q := deBoor(degree, knots, ctrl, weights, t)
		if i == 0 {
			p.MoveTo(q.X, q.Y)
		} else {
			p.LineTo(q.X, q.Y)
		}
	}
	return p
}

func nonDecreasing(v []float64) bool {
	for i := 1; i < len(v); i++ {
		if v[i] < v[i-1] || math.IsNaN(v[i]) {
			return false
		}
	}
	return true
}

func clamped(knots []float64, degree int) bool {
	m := len(knots) - 1
	for i := 1; i <= degree; i++ {
		if knots[i] != knots[0] || knots[m-i] != knots[m] {
			return false
		}
	}
	return true
}

// clampedKnots returns a clamped uniform knot vector.
func clampedKnots(n, degree int) []float64 {
	k := make([]float64, n+degree+1)
	inner := n - degree
	for i := range k {
		switch {
		case i <= degree:
			k[i] = 0
		case i >= n:
			k[i] = 1
		default:
			k[i] = float64(i-degree) / float64(inner)
		}
	}
	return k
}

// bezierSegments decomposes a clamped non-rational B-spline into Bézier
// segments by knot insertion (The NURBS Book, algorithm A5.6).
func bezierSegments(p int, knots []float64, ctrl []Point) [][]Point {
	n := len(ctrl) - 1
	m := n + p + 1
	a, b := p, p+1
	alphas := make([]float64, p+1)
	cur := append([]Point(nil), ctrl[:p+1]...)
	var out [][]Point
	for b < m {
		next := make([]Point, p+1)
		i := b
		for b < m && knots[b+1] == knots[b] {
			b++
		}
		mult := b - i + 1
		if mult < p {
			numer := knots[b] - knots[a]
			for j := p; j > mult; j-- {
				alphas[j-mult-1] = numer / (knots[a+j] - knots[a])
			}
			r := p - mult
			for j := 1; j <= r; j++ {
				save := r - j
				s := mult + j
				for k := p; k >= s; k-- {
					alpha := alphas[k-s]
					cur[k] = cur[k].Mul(alpha).Add(cur[k-1].Mul(1 - alpha))
				}
				if b < m {
					next[save] = cur[p]
				}
			}
		}
		out = append(out, cur)
		if b < m {
			for i := max(p-mult, 0); i <= p; i++ {
				next[i] = ctrl[b-p+i]
			}
			a = b
			b++
			cur = next
		}
	}
	return out
}

// deBoor evaluates a (rational) B-spline at t.
func deBoor(p int, knots []float64, ctrl []Point, weights []float64, t float64) Point {
	n := len(ctrl)
	k := p
	for k < n-1 && knots[k+1] <= t {
		k++
	}
	type hp struct{ x, y, w float64 }
	d := make([]hp, p+1)
	for j := 0; j <= p; j++ {
		c := ctrl[j+k-p]
		w := 1.0
		if weights != nil {
			w = weights[j+k-p]
		}
		d[j] = hp{c.X * w, c.Y * w, w}
	}
	for r := 1; r <= p; r++ {
		for j := p; j >= r; j-- {
			den := knots[j+1+k-r] - knots[j+k-p]
			alpha := 0.0
			if den != 0 {
				alpha = (t - knots[j+k-p]) / den
			}
			d[j] = hp{(1-alpha)*d[j-1].x + alpha*d[j].x, (1-alpha)*d[j-1].y + alpha*d[j].y, (1-alpha)*d[j-1].w + alpha*d[j].w}
		}
	}
	if d[p].w == 0 {
		return Point{d[p].x, d[p].y}
	}
	return Point{d[p].x / d[p].w, d[p].y / d[p].w}
}

// Interpolate adds a C2 cubic spline through the fit points (parametrized
// by chord length) as a new subpath. start and end, when not nil, are the
// tangent directions at the ends; otherwise the ends are natural.
func (p *Path) Interpolate(fit []Point, start, end *Point) *Path {
	// drop repeated points
	pts := fit[:0:0]
	for _, q := range fit {
		if len(pts) == 0 || q != pts[len(pts)-1] {
			pts = append(pts, q)
		}
	}
	p.open = false
	switch len(pts) {
	case 0:
		return p
	case 1:
		return p.MoveTo(pts[0].X, pts[0].Y)
	case 2:
		if start == nil && end == nil {
			return p.MoveTo(pts[0].X, pts[0].Y).LineTo(pts[1].X, pts[1].Y)
		}
	}
	n := len(pts) - 1
	h := make([]float64, n)
	for i := range h {
		h[i] = pts[i+1].Sub(pts[i]).Len()
	}
	unit := func(v *Point) *Point {
		if v == nil {
			return nil
		}
		l := v.Len()
		if l == 0 {
			return nil
		}
		u := v.Mul(1 / l)
		return &u
	}
	start, end = unit(start), unit(end)
	solve := func(y []float64, d0, dn *float64) []float64 {
		// second derivatives M from the tridiagonal system
		a := make([]float64, n+1)
		b := make([]float64, n+1)
		c := make([]float64, n+1)
		r := make([]float64, n+1)
		if d0 != nil {
			b[0], c[0] = 2*h[0], h[0]
			r[0] = 6 * ((y[1]-y[0])/h[0] - *d0)
		} else {
			b[0] = 1
		}
		for i := 1; i < n; i++ {
			a[i], b[i], c[i] = h[i-1], 2*(h[i-1]+h[i]), h[i]
			r[i] = 6 * ((y[i+1]-y[i])/h[i] - (y[i]-y[i-1])/h[i-1])
		}
		if dn != nil {
			a[n], b[n] = h[n-1], 2*h[n-1]
			r[n] = 6 * (*dn - (y[n]-y[n-1])/h[n-1])
		} else {
			b[n] = 1
		}
		// Thomas algorithm
		for i := 1; i <= n; i++ {
			w := a[i] / b[i-1]
			b[i] -= w * c[i-1]
			r[i] -= w * r[i-1]
		}
		m := make([]float64, n+1)
		m[n] = r[n] / b[n]
		for i := n - 1; i >= 0; i-- {
			m[i] = (r[i] - c[i]*m[i+1]) / b[i]
		}
		return m
	}
	xs := make([]float64, n+1)
	ys := make([]float64, n+1)
	for i, q := range pts {
		xs[i], ys[i] = q.X, q.Y
	}
	var sx, sy, ex, ey *float64
	if start != nil {
		sx, sy = &start.X, &start.Y
	}
	if end != nil {
		ex, ey = &end.X, &end.Y
	}
	mx := solve(xs, sx, ex)
	my := solve(ys, sy, ey)
	p.MoveTo(pts[0].X, pts[0].Y)
	for i := 0; i < n; i++ {
		d0 := func(v []float64, m []float64) float64 {
			return (v[i+1]-v[i])/h[i] - h[i]*(2*m[i]+m[i+1])/6
		}
		d1 := func(v []float64, m []float64) float64 {
			return (v[i+1]-v[i])/h[i] + h[i]*(m[i]+2*m[i+1])/6
		}
		c1 := Point{xs[i] + h[i]/3*d0(xs, mx), ys[i] + h[i]/3*d0(ys, my)}
		c2 := Point{xs[i+1] - h[i]/3*d1(xs, mx), ys[i+1] - h[i]/3*d1(ys, my)}
		p.CubicTo(c1, c2, pts[i+1])
	}
	return p
}

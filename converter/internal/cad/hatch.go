package cad

import "math"

// PatternLine is a family of parallel lines of a hatch pattern, in drawing
// units: the line through Base at Angle (radians), repeated every Offset
// (the displacement from one line to the next), dashed by Dash in the
// convention of CAD line types (positive: a dash, negative: a gap, 0: a
// dot; nil: continuous). The pattern of every line starts at the point
// Base + k·Offset.
type PatternLine struct {
	Angle  float64
	Base   Point
	Offset Point
	Dash   []float64
}

// MaxHatchLines bounds the lines drawn for one hatch.
const MaxHatchLines = 20000

// Hatch draws the lines of a pattern clipped by boundary (even-odd rule).
// It reports false, drawing nothing, when the pattern would need more than
// MaxHatchLines lines: the caller then fills the boundary instead.
func (d *Drawing) Hatch(boundary *Path, lines []PatternLine, pen Pen) bool {
	b := boundary.Bounds()
	if !b.ok {
		return true
	}
	corners := []Point{b.Min, {b.Max.X, b.Min.Y}, b.Max, {b.Min.X, b.Max.Y}}
	type family struct {
		path *Path
		pen  Pen
	}
	var fams []family
	total := 0
	for _, pl := range lines {
		u := Point{math.Cos(pl.Angle), math.Sin(pl.Angle)}
		n := Point{-u.Y, u.X}
		dn := pl.Offset.X*n.X + pl.Offset.Y*n.Y
		if math.Abs(dn) < 1e-9*math.Max(b.W(), b.H()) || math.IsNaN(dn) {
			continue
		}
		du := pl.Offset.X*u.X + pl.Offset.Y*u.Y
		bn := pl.Base.X*n.X + pl.Base.Y*n.Y
		bu := pl.Base.X*u.X + pl.Base.Y*u.Y
		nmin, nmax := math.Inf(1), math.Inf(-1)
		umin, umax := math.Inf(1), math.Inf(-1)
		for _, c := range corners {
			cn, cu := c.X*n.X+c.Y*n.Y, c.X*u.X+c.Y*u.Y
			nmin, nmax = min(nmin, cn), max(nmax, cn)
			umin, umax = min(umin, cu), max(umax, cu)
		}
		k0, k1 := (nmin-bn)/dn, (nmax-bn)/dn
		if k0 > k1 {
			k0, k1 = k1, k0
		}
		ka, kb := math.Ceil(k0), math.Floor(k1)
		count := kb - ka + 1
		if count <= 0 {
			continue
		}
		total += int(count)
		if total > MaxHatchLines {
			return false
		}
		dash, _ := DashPattern(pl.Dash)
		plen := 0.0
		for _, v := range dash {
			plen += v
		}
		path := &Path{}
		for k := ka; k <= kb; k++ {
			// the start of this line's pattern, along the line
			su := bu + k*du
			start := umin
			if plen > 0 {
				start = su + math.Floor((umin-su)/plen)*plen
			}
			// the point of the line at the coordinates (start, n)
			ln := bn + k*dn
			a := Point{u.X*start + n.X*ln, u.Y*start + n.Y*ln}
			e := Point{u.X*umax + n.X*ln, u.Y*umax + n.Y*ln}
			path.MoveTo(a.X, a.Y).LineTo(e.X, e.Y)
		}
		p := pen
		if dash != nil {
			p.Dash = dash
			_, p.DashOffset = DashPattern(pl.Dash)
		}
		fams = append(fams, family{path, p})
	}
	d.Begin(boundary)
	for _, f := range fams {
		d.Stroke(f.path, f.pen)
	}
	d.End()
	return true
}

// DashPattern converts a CAD line type pattern (positive: a dash,
// negative: a gap, 0: a dot) into a dash array that starts with a dash and
// alternates dashes and gaps, and the offset into it at which the pattern
// begins. It returns nil for a continuous line.
func DashPattern(elems []float64) (dash []float64, offset float64) {
	var out []float64 // alternating, starting with a dash (possibly a gap first, handled below)
	lastDash := false
	started := false
	lead := 0.0 // leading gap
	for _, v := range elems {
		isDash := v >= 0
		l := math.Abs(v)
		if !started {
			if !isDash {
				lead += l
				continue
			}
			started = true
			out = append(out, l)
			lastDash = true
			continue
		}
		if isDash == lastDash {
			out[len(out)-1] += l
			continue
		}
		out = append(out, l)
		lastDash = isDash
	}
	if !started {
		return nil, 0
	}
	if lastDash {
		// the pattern ends with a dash: the leading gap follows it
		if lead > 0 {
			out = append(out, lead)
		} else {
			// dashes only: continuous, unless there are dots
			allSolid := true
			for _, v := range out {
				if v == 0 {
					allSolid = false
				}
			}
			if allSolid {
				return nil, 0
			}
			out = append(out, 0)
		}
	} else {
		out[len(out)-1] += lead
	}
	if len(out)%2 == 1 {
		out = append(out, 0)
	}
	total := 0.0
	for _, v := range out {
		total += v
	}
	if total <= 0 {
		return nil, 0
	}
	if lead > 0 {
		// the pattern starts in the gap that closes the array
		offset = total - lead
	}
	return out, offset
}

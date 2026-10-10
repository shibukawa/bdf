package imagebdf

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf/image/svg"
)

// parsePathData reads SVG path data into a path, up to the first error.
func parsePathData(d string) *path {
	p := &path{}
	var cmd byte
	var cx, cy, sx, sy float64 // current point, sub-path start
	var lcx, lcy float64       // last control point
	var lastCmd byte
	s := d
	nums := func(n int) ([]float64, bool) {
		out := make([]float64, n)
		for i := range out {
			v, rest, ok := svg.ParseNumber(s)
			if !ok {
				return nil, false
			}
			out[i] = v
			s = rest
		}
		return out, true
	}
	flag := func() (bool, bool) {
		s = strings.TrimLeft(s, " ,\t\n\r")
		if s == "" || (s[0] != '0' && s[0] != '1') {
			return false, false
		}
		f := s[0] == '1'
		s = s[1:]
		return f, true
	}
	for {
		s = strings.TrimLeft(s, " ,\t\n\r")
		if s == "" {
			return p
		}
		c := s[0]
		if strings.IndexByte("MmLlHhVvCcSsQqTtAaZz", c) >= 0 {
			cmd = c
			s = s[1:]
		} else if cmd == 0 {
			return p
		}
		rel := cmd >= 'a'
		ox, oy := 0.0, 0.0
		if rel {
			ox, oy = cx, cy
		}
		switch cmd {
		case 'M', 'm':
			a, ok := nums(2)
			if !ok {
				return p
			}
			cx, cy = a[0]+ox, a[1]+oy
			sx, sy = cx, cy
			p.moveTo(cx, cy)
			// further pairs are line-tos
			if cmd == 'M' {
				cmd = 'L'
			} else {
				cmd = 'l'
			}
			lastCmd = 'M'
			continue
		case 'L', 'l':
			a, ok := nums(2)
			if !ok {
				return p
			}
			cx, cy = a[0]+ox, a[1]+oy
			p.lineTo(cx, cy)
		case 'H', 'h':
			a, ok := nums(1)
			if !ok {
				return p
			}
			cx = a[0] + ox
			p.lineTo(cx, cy)
		case 'V', 'v':
			a, ok := nums(1)
			if !ok {
				return p
			}
			cy = a[0] + oy
			p.lineTo(cx, cy)
		case 'C', 'c':
			a, ok := nums(6)
			if !ok {
				return p
			}
			p.cubicTo(a[0]+ox, a[1]+oy, a[2]+ox, a[3]+oy, a[4]+ox, a[5]+oy)
			lcx, lcy = a[2]+ox, a[3]+oy
			cx, cy = a[4]+ox, a[5]+oy
		case 'S', 's':
			a, ok := nums(4)
			if !ok {
				return p
			}
			x1, y1 := cx, cy
			if lastCmd == 'C' || lastCmd == 'S' {
				x1, y1 = 2*cx-lcx, 2*cy-lcy
			}
			p.cubicTo(x1, y1, a[0]+ox, a[1]+oy, a[2]+ox, a[3]+oy)
			lcx, lcy = a[0]+ox, a[1]+oy
			cx, cy = a[2]+ox, a[3]+oy
		case 'Q', 'q':
			a, ok := nums(4)
			if !ok {
				return p
			}
			p.quadTo(a[0]+ox, a[1]+oy, a[2]+ox, a[3]+oy)
			lcx, lcy = a[0]+ox, a[1]+oy
			cx, cy = a[2]+ox, a[3]+oy
		case 'T', 't':
			a, ok := nums(2)
			if !ok {
				return p
			}
			x1, y1 := cx, cy
			if lastCmd == 'Q' || lastCmd == 'T' {
				x1, y1 = 2*cx-lcx, 2*cy-lcy
			}
			p.quadTo(x1, y1, a[0]+ox, a[1]+oy)
			lcx, lcy = x1, y1
			cx, cy = a[0]+ox, a[1]+oy
		case 'A', 'a':
			a, ok := nums(3)
			if !ok {
				return p
			}
			large, ok1 := flag()
			sweep, ok2 := flag()
			e, ok3 := nums(2)
			if !ok1 || !ok2 || !ok3 {
				return p
			}
			x, y := e[0]+ox, e[1]+oy
			arcTo(p, cx, cy, a[0], a[1], a[2], large, sweep, x, y)
			cx, cy = x, y
		case 'Z', 'z':
			p.close()
			cx, cy = sx, sy
		}
		lastCmd = cmd &^ 0x20 // upper case
	}
}

// arcTo adds an SVG elliptical arc from (x0, y0) to (x, y).
func arcTo(p *path, x0, y0, rx, ry, rotDeg float64, large, sweep bool, x, y float64) {
	if x0 == x && y0 == y {
		return
	}
	rx, ry = math.Abs(rx), math.Abs(ry)
	if rx == 0 || ry == 0 {
		p.lineTo(x, y)
		return
	}
	phi := rotDeg * math.Pi / 180
	sn, cs := math.Sincos(phi)
	dx, dy := (x0-x)/2, (y0-y)/2
	x1p, y1p := cs*dx+sn*dy, -sn*dx+cs*dy
	lambda := x1p*x1p/(rx*rx) + y1p*y1p/(ry*ry)
	if lambda > 1 {
		k := math.Sqrt(lambda)
		rx, ry = rx*k, ry*k
	}
	num := rx*rx*ry*ry - rx*rx*y1p*y1p - ry*ry*x1p*x1p
	den := rx*rx*y1p*y1p + ry*ry*x1p*x1p
	co := 0.0
	if den > 0 && num > 0 {
		co = math.Sqrt(num / den)
	}
	if large == sweep {
		co = -co
	}
	cxp, cyp := co*rx*y1p/ry, -co*ry*x1p/rx
	cx := cs*cxp - sn*cyp + (x0+x)/2
	cy := sn*cxp + cs*cyp + (y0+y)/2
	ang := func(ux, uy, vx, vy float64) float64 {
		a := math.Atan2(ux*vy-uy*vx, ux*vx+uy*vy)
		return a
	}
	t1 := ang(1, 0, (x1p-cxp)/rx, (y1p-cyp)/ry)
	dt := ang((x1p-cxp)/rx, (y1p-cyp)/ry, (-x1p-cxp)/rx, (-y1p-cyp)/ry)
	if !sweep && dt > 0 {
		dt -= 2 * math.Pi
	} else if sweep && dt < 0 {
		dt += 2 * math.Pi
	}
	p.ellipse(cx, cy, rx, ry, phi, t1, t1+dt, !sweep)
}

package drawio

// The arrow markers draw.io adds to mxMarker: those of Shapes.js (dash,
// box, cross, circle, circlePlus, halfCircle, async, openAsync and the
// mermaid markers) and those of the ER and SysML shape libraries (mxER.js,
// mxSysML.js), which draw.io always loads.

func init() {
	markers["dash"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		nx, ny := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x-nx/2-ny/2, pe.y-ny/2+nx/2)
			c.lineTo(pe.x+ny/2-3*nx/2, pe.y-3*ny/2-nx/2)
			c.stroke()
		}
	}
	markers["box"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		nx, ny := unitX*(size+sw+1), unitY*(size+sw+1)
		px, py := pe.x+nx/2, pe.y+ny/2
		pe.x -= nx
		pe.y -= ny
		return func() {
			c.begin()
			c.moveTo(px-nx/2-ny/2, py-ny/2+nx/2)
			c.lineTo(px-nx/2+ny/2, py-ny/2-nx/2)
			c.lineTo(px+ny/2-3*nx/2, py-3*ny/2-nx/2)
			c.lineTo(px-ny/2-3*nx/2, py-3*ny/2+nx/2)
			c.close()
			fillOrStroke(c, filled)
		}
	}
	markers["cross"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		nx, ny := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x-nx/2-ny/2, pe.y-ny/2+nx/2)
			c.lineTo(pe.x+ny/2-3*nx/2, pe.y-3*ny/2-nx/2)
			c.moveTo(pe.x-nx/2+ny/2, pe.y-ny/2-nx/2)
			c.lineTo(pe.x-ny/2-3*nx/2, pe.y-3*ny/2+nx/2)
			c.stroke()
		}
	}
	markers["circle"] = circleMarker
	markers["circlePlus"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		pt := *pe
		fn := circleMarker(c, s, typ, pe, unitX, unitY, size, source, sw, filled)
		nx, ny := unitX*(size+2*sw), unitY*(size+2*sw)
		return func() {
			fn()
			c.begin()
			c.moveTo(pt.x-unitX*sw, pt.y-unitY*sw)
			c.lineTo(pt.x-2*nx+unitX*sw, pt.y-2*ny+unitY*sw)
			c.moveTo(pt.x-nx-ny+unitY*sw, pt.y-ny+nx-unitX*sw)
			c.lineTo(pt.x+ny-nx-unitY*sw, pt.y-ny-nx+unitX*sw)
			c.stroke()
		}
	}
	markers["halfCircle"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		nx, ny := unitX*(size+sw+1), unitY*(size+sw+1)
		pt := *pe
		pe.x -= nx
		pe.y -= ny
		return func() {
			c.begin()
			c.moveTo(pt.x-ny, pt.y+nx)
			c.quadTo(pe.x-ny, pe.y+nx, pe.x, pe.y)
			c.quadTo(pe.x+ny, pe.y-nx, pt.x+ny, pt.y-nx)
			c.stroke()
		}
	}
	markers["async"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		// see createArrow for the 1.118
		endOffsetX, endOffsetY := unitX*sw*1.118, unitY*sw*1.118
		unitX *= size + sw
		unitY *= size + sw
		pt := point{pe.x - endOffsetX, pe.y - endOffsetY}
		pe.x += -unitX - endOffsetX
		pe.y += -unitY - endOffsetY
		return func() {
			c.begin()
			c.moveTo(pt.x, pt.y)
			if source {
				c.lineTo(pt.x-unitX-unitY/2, pt.y-unitY+unitX/2)
			} else {
				c.lineTo(pt.x+unitY/2-unitX, pt.y-unitY-unitX/2)
			}
			c.lineTo(pt.x-unitX, pt.y-unitY)
			c.close()
			fillOrStroke(c, filled)
		}
	}
	markers["openAsync"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		const widthFactor = 2
		unitX *= size + sw
		unitY *= size + sw
		pt := *pe
		return func() {
			c.begin()
			c.moveTo(pt.x, pt.y)
			if source {
				c.lineTo(pt.x-unitX-unitY/widthFactor, pt.y-unitY+unitX/widthFactor)
			} else {
				c.lineTo(pt.x+unitY/widthFactor-unitX, pt.y-unitY-unitX/widthFactor)
			}
			c.stroke()
		}
	}
	// Mermaid's class diagram markers: 17 long, 12 wide (17:6 from the tip
	// to the side), with the tip on the end point.
	markers["mermaidExtension"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		const wf = 17.0 / 6
		unitX *= size + sw
		unitY *= size + sw
		pt := *pe
		pe.x -= unitX
		pe.y -= unitY
		return func() {
			c.begin()
			c.moveTo(pt.x, pt.y)
			c.lineTo(pt.x-unitX-unitY/wf, pt.y-unitY+unitX/wf)
			c.lineTo(pt.x-unitX+unitY/wf, pt.y-unitY-unitX/wf)
			c.close()
			fillOrStroke(c, filled)
		}
	}
	markers["mermaidDiamond"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		const wf = 17.0 / 6
		unitX *= size + sw
		unitY *= size + sw
		pt := *pe
		pe.x -= unitX
		pe.y -= unitY
		return func() {
			c.begin()
			c.moveTo(pt.x, pt.y)
			c.lineTo(pt.x-unitX/2-unitY/wf, pt.y-unitY/2+unitX/wf)
			c.lineTo(pt.x-unitX, pt.y-unitY)
			c.lineTo(pt.x-unitX/2+unitY/wf, pt.y-unitY/2-unitX/wf)
			c.close()
			fillOrStroke(c, filled)
		}
	}

	// mxER.js: crow's feet
	markers["ERone"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x-m/2-n/2, pe.y-n/2+m/2)
			c.lineTo(pe.x-m/2+n/2, pe.y-n/2-m/2)
			c.stroke()
		}
	}
	markers["ERmandOne"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x-m/2-n/2, pe.y-n/2+m/2)
			c.lineTo(pe.x-m/2+n/2, pe.y-n/2-m/2)
			c.moveTo(pe.x-m-n/2, pe.y-n+m/2)
			c.lineTo(pe.x-m+n/2, pe.y-n-m/2)
			c.stroke()
		}
	}
	markers["ERmany"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x+n/2, pe.y-m/2)
			c.lineTo(pe.x-m, pe.y-n)
			c.lineTo(pe.x-n/2, pe.y+m/2)
			c.stroke()
		}
	}
	markers["ERoneToMany"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x-m-n/2, pe.y-n+m/2)
			c.lineTo(pe.x-m+n/2, pe.y-n-m/2)
			c.moveTo(pe.x+n/2, pe.y-m/2)
			c.lineTo(pe.x-m, pe.y-n)
			c.lineTo(pe.x-n/2, pe.y+m/2)
			c.stroke()
		}
	}
	markers["ERzeroToMany"] = erZeroMarker(true)
	markers["ERzeroToOne"] = erZeroMarker(false)

	// mxSysML.js
	markers["sysMLx"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n := unitX*(size+sw+1), unitY*(size+sw+1)
		return func() {
			c.begin()
			c.moveTo(pe.x-m/2-n/2, pe.y-n/2+m/2)
			c.lineTo(pe.x+m/2+n/2, pe.y+n/2-m/2)
			c.moveTo(pe.x+m/2-n/2, pe.y+n/2+m/2)
			c.lineTo(pe.x-m/2+n/2, pe.y-n/2-m/2)
			c.stroke()
		}
	}
	markers["sysMLLost"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n, p := unitX*(size+sw+1), unitY*(size+sw+1), size/2
		return func() {
			c.begin()
			c.moveTo(pe.x-1.5*m-n/2, pe.y-1.5*n+m/2)
			c.lineTo(pe.x-m/2, pe.y-n/2)
			c.lineTo(pe.x-1.5*m+n/2, pe.y-1.5*n-m/2)
			c.stroke()
			c.ellipse(pe.x-0.5*m-p, pe.y-0.5*n-p, 2*p, 2*p)
			c.setFillColor(s.style.get("strokeColor", "#000000"))
			c.fillAndStroke()
		}
	}
	markers["sysMLFound"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n, p := unitX*(size+sw+1), unitY*(size+sw+1), size/2
		return func() {
			c.ellipse(pe.x-0.5*m-p, pe.y-0.5*n-p, 2*p, 2*p)
			c.setFillColor(s.style.get("strokeColor", "#000000"))
			c.fillAndStroke()
		}
	}
	markers["sysMLPackCont"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n, p := unitX*(size+sw+1), unitY*(size+sw+1), size/2
		return func() {
			c.begin()
			c.moveTo(pe.x-m/2-n/2, pe.y-n/2+m/2)
			c.lineTo(pe.x-m/2+n/2, pe.y-n/2-m/2)
			c.stroke()
			c.ellipse(pe.x-0.5*m-p, pe.y-0.5*n-p, 2*p, 2*p)
			c.stroke()
		}
	}
	markers["sysMLReqInt"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n, p := unitX*(size+sw+1), unitY*(size+sw+1), size/2
		return func() {
			c.setFillColor(s.style.get("fillColor", "none"))
			c.ellipse(pe.x-0.5*m-p, pe.y-0.5*n-p, 2*p, 2*p)
			c.fillAndStroke()
		}
	}
	markers["sysMLProvInt"] = func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n, p := unitX*(size+sw+1), unitY*(size+sw+1), size/2
		return func() {
			c.setFillColor(s.style.get("fillColor", "none"))
			c.begin()
			c.moveTo(pe.x-n/2, pe.y+m/2)
			c.arcTo(p, p, 0, false, true, pe.x+n/2, pe.y-m/2)
			c.fillAndStroke()
		}
	}
}

// circleMarker is Shapes.js's circleMarker: a circle behind the end point.
func circleMarker(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
	size += sw
	pt := *pe
	pe.x -= unitX * (2*size + sw)
	pe.y -= unitY * (2*size + sw)
	unitX *= size + sw
	unitY *= size + sw
	return func() {
		c.ellipse(pt.x-unitX-size, pt.y-unitY-size, 2*size, 2*size)
		fillOrStroke(c, filled)
	}
}

// erZeroMarker is ERzeroToMany (many) and ERzeroToOne of mxER.js: a circle
// and a crow's foot or a bar. Filled, the circle is white over the line;
// else the line stops at the circle.
func erZeroMarker(many bool) markerFunc {
	return func(c *c2d, s *shape, typ string, pe *point, unitX, unitY, size float64, source bool, sw float64, filled bool) func() {
		m, n, p := unitX*(size+sw+1), unitY*(size+sw+1), size/2
		q, u := pe.x, pe.y
		if !filled {
			pe.x -= 2*m - unitX*sw/2
			pe.y -= 2*n - unitY*sw/2
		}
		return func() {
			c.begin()
			c.ellipse(q-1.5*m-p, u-1.5*n-p, 2*p, 2*p)
			if filled {
				t := s.style.get("strokeColor", "#666666")
				c.setFillColor("#ffffff")
				c.fillAndStroke()
				c.setFillColor(t)
			} else {
				c.stroke()
			}
			c.begin()
			if many {
				c.moveTo(q+n/2, u-m/2)
				c.lineTo(q-m, u-n)
				c.lineTo(q-n/2, u+m/2)
				if !filled {
					c.moveTo(q-m, u-n)
					c.lineTo(q, u)
				}
			} else {
				c.moveTo(q-m/2-n/2, u-n/2+m/2)
				c.lineTo(q-m/2+n/2, u-n/2-m/2)
				if !filled {
					c.moveTo(q-m-unitX*sw/2, u-n-unitY*sw/2)
					c.lineTo(q, u)
				}
			}
			c.stroke()
		}
	}
}

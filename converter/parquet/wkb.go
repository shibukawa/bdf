package parquet

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"
)

// Geometries (the GEOMETRY and GEOGRAPHY types, and the columns GeoParquet
// metadata marks) are stored as well-known binary (WKB), and shown as
// well-known text (WKT). ISO WKB and the extended WKB of PostGIS (with the
// Z, M and SRID flags) are read.

var wkbTypes = [...]string{1: "POINT", 2: "LINESTRING", 3: "POLYGON", 4: "MULTIPOINT", 5: "MULTILINESTRING", 6: "MULTIPOLYGON", 7: "GEOMETRYCOLLECTION"}

// maxWKTSize bounds the text made of a geometry.
const maxWKTSize = maxCellText + 1

type wkbReader struct {
	b   []byte
	pos int
	out strings.Builder
	bad bool
}

// wkt converts a WKB geometry into WKT.
func wkt(b []byte) (string, bool) {
	r := &wkbReader{b: b}
	r.geometry(0, true)
	if r.bad || r.pos != len(b) && r.out.Len() < maxWKTSize {
		return "", false
	}
	return r.out.String(), true
}

func (r *wkbReader) geometry(depth int, tagged bool) {
	if depth > 32 || len(r.b)-r.pos < 5 {
		r.bad = true
		return
	}
	var order binary.ByteOrder = binary.LittleEndian
	switch r.b[r.pos] {
	case 0:
		order = binary.BigEndian
	case 1:
	default:
		r.bad = true
		return
	}
	typ := order.Uint32(r.b[r.pos+1:])
	r.pos += 5
	z, m := typ&0x80000000 != 0, typ&0x40000000 != 0
	if typ&0x20000000 != 0 {
		r.pos += 4 // the SRID
	}
	typ &= 0x0fffffff
	switch typ / 1000 {
	case 1:
		z = true
	case 2:
		m = true
	case 3:
		z, m = true, true
	}
	typ %= 1000
	if typ == 0 || typ >= uint32(len(wkbTypes)) || r.pos > len(r.b) {
		r.bad = true
		return
	}
	dims := 2
	if z {
		dims++
	}
	if m {
		dims++
	}
	if tagged {
		r.out.WriteString(wkbTypes[typ])
		switch {
		case z && m:
			r.out.WriteString(" ZM")
		case z:
			r.out.WriteString(" Z")
		case m:
			r.out.WriteString(" M")
		}
		r.out.WriteString(" ")
	}
	switch typ {
	case 1:
		r.point(order, dims)
	case 2:
		r.points(order, dims)
	case 3:
		r.rings(order, dims)
	default:
		n := r.count(order, 5)
		if n == 0 {
			r.out.WriteString("EMPTY")
			return
		}
		r.out.WriteString("(")
		for i := range n {
			if r.bad || r.out.Len() >= maxWKTSize {
				return
			}
			if i > 0 {
				r.out.WriteString(", ")
			}
			// the parts of a multi geometry are shown without their type
			// (as POINT, LINESTRING or POLYGON), those of a collection with
			r.geometry(depth+1, typ == 7)
		}
		r.out.WriteString(")")
	}
}

// count reads a number of items that take size bytes each at least.
func (r *wkbReader) count(order binary.ByteOrder, size int) int {
	if len(r.b)-r.pos < 4 {
		r.bad = true
		return 0
	}
	n := order.Uint32(r.b[r.pos:])
	r.pos += 4
	if uint64(n)*uint64(size) > uint64(len(r.b)-r.pos) {
		r.bad = true
		return 0
	}
	return int(n)
}

func (r *wkbReader) coords(order binary.ByteOrder, dims int) []float64 {
	if len(r.b)-r.pos < 8*dims {
		r.bad = true
		return nil
	}
	c := make([]float64, dims)
	for i := range c {
		c[i] = math.Float64frombits(order.Uint64(r.b[r.pos:]))
		r.pos += 8
	}
	return c
}

func (r *wkbReader) writeCoords(c []float64) {
	for i, v := range c {
		if i > 0 {
			r.out.WriteByte(' ')
		}
		r.out.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	}
}

// point writes a point, which is empty when its coordinates are NaN.
func (r *wkbReader) point(order binary.ByteOrder, dims int) {
	c := r.coords(order, dims)
	if r.bad {
		return
	}
	if math.IsNaN(c[0]) && math.IsNaN(c[1]) {
		r.out.WriteString("EMPTY")
		return
	}
	r.out.WriteString("(")
	r.writeCoords(c)
	r.out.WriteString(")")
}

func (r *wkbReader) points(order binary.ByteOrder, dims int) {
	n := r.count(order, 8*dims)
	if r.bad {
		return
	}
	if n == 0 {
		r.out.WriteString("EMPTY")
		return
	}
	r.pointList(order, dims, n)
}

func (r *wkbReader) pointList(order binary.ByteOrder, dims, n int) {
	r.out.WriteString("(")
	for i := range n {
		if r.out.Len() >= maxWKTSize {
			r.pos += (n - i) * 8 * dims
			break
		}
		if i > 0 {
			r.out.WriteString(", ")
		}
		r.writeCoords(r.coords(order, dims))
	}
	r.out.WriteString(")")
}

func (r *wkbReader) rings(order binary.ByteOrder, dims int) {
	n := r.count(order, 4)
	if r.bad {
		return
	}
	if n == 0 {
		r.out.WriteString("EMPTY")
		return
	}
	r.out.WriteString("(")
	for i := range n {
		if i > 0 {
			r.out.WriteString(", ")
		}
		k := r.count(order, 8*dims)
		if r.bad {
			return
		}
		r.pointList(order, dims, k)
	}
	r.out.WriteString(")")
}

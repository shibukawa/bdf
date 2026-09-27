package parquet

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Values are shown as data tools show them: integers and decimals in full,
// floating-point numbers in their shortest form (with an exponent only
// when they are very large or small, as JavaScript writes them), dates,
// times and timestamps in ISO 8601 with a space (in UTC for the instants),
// UUIDs in their canonical form, geometries in WKT, and byte arrays as
// text when they are, else in hexadecimal. Numbers, dates and times are
// aligned to the right, as Excel aligns them.

type fmtKind uint8

const (
	fBool fmtKind = iota
	fInt
	fUint
	fDecimal
	fFloat
	fDouble
	fFloat16
	fDate
	fTime
	fTimestamp
	fInt96
	fString
	fJSON
	fUUID
	fInterval
	fGeometry
	fBinary // text when it is, else hexadecimal
	fBytes  // hexadecimal
	fNull
)

// format shows the values of a leaf.
type format struct {
	kind  fmtKind
	phys  int32
	unit  int16 // of times and timestamps
	scale int32 // of decimals
}

func formatOf(n *node) format {
	f := format{kind: fBinary, phys: n.phys}
	p, l := n.phys, n.logical
	bytes := p == typeBinary || p == typeFixed
	switch {
	case (l.kind == logString || l.kind == logEnum) && bytes:
		f.kind = fString
	case l.kind == logJSON && bytes:
		f.kind = fJSON
	case l.kind == logBSON && bytes:
		f.kind = fBytes
	case l.kind == logUUID && p == typeFixed && n.typeLen == 16:
		f.kind = fUUID
	case l.kind == logFloat16 && p == typeFixed && n.typeLen == 2:
		f.kind = fFloat16
	case l.kind == logDecimal && (p == typeInt32 || p == typeInt64 || bytes):
		f.kind, f.scale = fDecimal, l.scale
	case l.kind == logDate && p == typeInt32:
		f.kind = fDate
	case l.kind == logTime && (p == typeInt32 && l.unit == unitMillis || p == typeInt64 && (l.unit == unitMicros || l.unit == unitNanos)):
		f.kind, f.unit = fTime, l.unit
	case l.kind == logTimestamp && p == typeInt64 && l.unit >= unitMillis && l.unit <= unitNanos:
		f.kind, f.unit = fTimestamp, l.unit
	case l.kind == logInteger && (p == typeInt32 || p == typeInt64):
		f.kind = fInt
		if !l.signed {
			f.kind = fUint
		}
	case l.kind == logUnknown:
		f.kind = fNull
	case (l.kind == logGeometry || l.kind == logGeography) && bytes:
		f.kind = fGeometry
	default:
		f = formatOfConverted(n)
	}
	return f
}

func formatOfConverted(n *node) format {
	f := format{kind: fBinary, phys: n.phys}
	p := n.phys
	bytes := p == typeBinary || p == typeFixed
	switch c := n.conv; {
	case (c == convUTF8 || c == convEnum) && bytes:
		f.kind = fString
	case c == convJSON && bytes:
		f.kind = fJSON
	case c == convBSON && bytes:
		f.kind = fBytes
	case c == convDate && p == typeInt32:
		f.kind = fDate
	case c == convTimeMillis && p == typeInt32:
		f.kind, f.unit = fTime, unitMillis
	case c == convTimeMicros && p == typeInt64:
		f.kind, f.unit = fTime, unitMicros
	case c == convTSMillis && p == typeInt64:
		f.kind, f.unit = fTimestamp, unitMillis
	case c == convTSMicros && p == typeInt64:
		f.kind, f.unit = fTimestamp, unitMicros
	case c >= convUint8 && c <= convUint64 && (p == typeInt32 || p == typeInt64):
		f.kind = fUint
	case c == convInterval && p == typeFixed && n.typeLen == 12:
		f.kind = fInterval
	case n.geo && bytes:
		f.kind = fGeometry
	case p == typeBoolean:
		f.kind = fBool
	case p == typeInt32 || p == typeInt64:
		f.kind = fInt
	case p == typeInt96:
		f.kind = fInt96
	case p == typeFloat:
		f.kind = fFloat
	case p == typeDouble:
		f.kind = fDouble
	}
	return f
}

// number reports whether the values are aligned to the right.
func (f format) number() bool {
	switch f.kind {
	case fInt, fUint, fDecimal, fFloat, fDouble, fFloat16, fDate, fTime, fTimestamp, fInt96, fInterval:
		return true
	}
	return false
}

// quoted reports whether the values are strings within nested values.
func (f format) quoted() bool {
	switch f.kind {
	case fBool, fInt, fUint, fDecimal, fFloat, fDouble, fFloat16, fJSON, fNull:
		return false
	}
	return true
}

// timed reports whether the values have fractions of seconds, which a
// column shows with as many digits for all its values.
func (f format) timed() bool { return f.kind == fTime || f.kind == fTimestamp || f.kind == fInt96 }

// maxDigits is the number of digits of the fractions of seconds of the
// unit.
func (f format) maxDigits() int {
	switch {
	case f.kind == fInt96 || f.unit == unitNanos:
		return 9
	case f.unit == unitMicros:
		return 6
	}
	return 3
}

// nanos returns the seconds and nanoseconds of a time (of the day) or a
// timestamp (since 1970).
func (f format) nanos(v scalar) (sec, nsec int64) {
	x := int64(v.u)
	switch {
	case f.kind == fInt96:
		days := int64(v.x) - 2440588 // the Julian day of 1970-01-01
		sec, nsec = days*86400+x/1e9, x%1e9
	case f.unit == unitMillis:
		sec, nsec = x/1e3, x%1e3*1e6
	case f.unit == unitMicros:
		sec, nsec = x/1e6, x%1e6*1e3
	default:
		sec, nsec = x/1e9, x%1e9
	}
	if nsec < 0 {
		sec, nsec = sec-1, nsec+1e9
	}
	return sec, nsec
}

// digits returns the digits of the fraction of a second a value needs.
func (f format) digits(v scalar) int {
	_, ns := f.nanos(v)
	return min(fracDigits(ns), f.maxDigits())
}

func fracDigits(ns int64) int {
	switch {
	case ns == 0:
		return 0
	case ns%1e6 == 0:
		return 3
	case ns%1e3 == 0:
		return 6
	}
	return 9
}

// text shows a value; digits is the number of digits of the fractions of
// seconds of times and timestamps, or -1 for those the value needs.
func (f format) text(v scalar, digits int) string {
	switch f.kind {
	case fBool:
		if v.u != 0 {
			return "true"
		}
		return "false"
	case fInt:
		return strconv.FormatInt(int64(v.u), 10)
	case fUint:
		if f.phys == typeInt32 {
			return strconv.FormatUint(uint64(uint32(v.u)), 10)
		}
		return strconv.FormatUint(v.u, 10)
	case fDecimal:
		switch f.phys {
		case typeInt32, typeInt64:
			return decimalText(big.NewInt(int64(v.u)), f.scale)
		}
		return decimalText(bigEndianInt(v.b), f.scale)
	case fFloat:
		return floatText(float64(math.Float32frombits(uint32(v.u))), 32)
	case fDouble:
		return floatText(math.Float64frombits(v.u), 64)
	case fFloat16:
		return halfText(binary.LittleEndian.Uint16(v.b))
	case fDate:
		return dateText(int64(int32(v.u)))
	case fTime, fTimestamp, fInt96:
		if digits < 0 {
			digits = f.digits(v)
		}
		sec, ns := f.nanos(v)
		if f.kind == fTime {
			return clockText(sec, ns, digits)
		}
		return timestampText(sec, ns, digits)
	case fString, fJSON:
		return validText(v.b)
	case fUUID:
		return uuidText(v.b)
	case fInterval:
		return intervalText(v.b)
	case fGeometry:
		if s, ok := wkt(v.b); ok {
			return s
		}
		return hexText(v.b)
	case fBytes:
		return hexText(v.b)
	case fNull:
		return ""
	}
	return binaryText(v.b)
}

func validText(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	return strings.ToValidUTF8(string(b), "�")
}

// bigEndianInt reads a two's complement big-endian integer.
func bigEndianInt(b []byte) *big.Int {
	if len(b) <= 8 {
		var v int64
		for _, c := range b {
			v = v<<8 | int64(c)
		}
		if len(b) > 0 && len(b) < 8 && b[0]&0x80 != 0 {
			v -= 1 << (8 * len(b))
		}
		return big.NewInt(v)
	}
	x := new(big.Int).SetBytes(b)
	if b[0]&0x80 != 0 {
		x.Sub(x, new(big.Int).Lsh(big.NewInt(1), uint(8*len(b))))
	}
	return x
}

// decimalText shows an unscaled integer with a scale.
func decimalText(x *big.Int, scale int32) string {
	s := x.String()
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	switch {
	case scale > 0 && scale <= 1000:
		if len(s) <= int(scale) {
			s = strings.Repeat("0", int(scale)-len(s)+1) + s
		}
		s = s[:len(s)-int(scale)] + "." + s[len(s)-int(scale):]
	case scale < 0 && scale >= -1000:
		s += strings.Repeat("0", int(-scale))
	}
	if neg {
		return "-" + s
	}
	return s
}

// floatText shows a floating-point number in its shortest form.
func floatText(v float64, bitSize int) string {
	switch {
	case math.IsNaN(v):
		return "NaN"
	case math.IsInf(v, 1):
		return "Infinity"
	case math.IsInf(v, -1):
		return "-Infinity"
	}
	if a := math.Abs(v); a == 0 || a >= 1e-6 && a < 1e21 {
		return strconv.FormatFloat(v, 'f', -1, bitSize)
	}
	return strconv.FormatFloat(v, 'e', -1, bitSize)
}

// halfValue converts an IEEE 754 half-precision number.
func halfValue(h uint16) float64 {
	sign := 1.0
	if h&0x8000 != 0 {
		sign = -1
	}
	exp, frac := int(h>>10&0x1f), float64(h&0x3ff)
	switch exp {
	case 0:
		return sign * math.Ldexp(frac, -24)
	case 31:
		if frac != 0 {
			return math.NaN()
		}
		return math.Inf(int(sign))
	}
	return sign * math.Ldexp(1024+frac, exp-25)
}

// halfText shows a half-precision number with the fewest digits that
// read back as it.
func halfText(h uint16) string {
	v := halfValue(h)
	m := h & 0x7fff
	if m == 0 || m >= 0x7c00 {
		return floatText(v, 64)
	}
	a := math.Abs(v)
	lo := halfValue(m - 1)
	hi := halfValue(m + 1)
	if m+1 >= 0x7c00 {
		hi = a + (a - lo)
	}
	for p := range 6 {
		d, _ := strconv.ParseFloat(strconv.FormatFloat(a, 'e', p, 64), 64)
		if math.Abs(d-a) < math.Abs(d-lo) && math.Abs(d-a) < math.Abs(d-hi) {
			if v < 0 {
				d = -d
			}
			return floatText(d, 64)
		}
	}
	return floatText(v, 64)
}

func dateText(days int64) string {
	return time.Unix(days*86400, 0).UTC().Format(time.DateOnly)
}

// fracText shows nanoseconds as a fraction of a second in digits.
func fracText(ns int64, digits int) string {
	if digits <= 0 {
		return ""
	}
	s := strconv.FormatInt(ns+1e9, 10) // 1 and nine digits
	return "." + s[1:1+min(digits, 9)]
}

func clockText(sec, ns int64, digits int) string {
	sec %= 86400
	if sec < 0 {
		sec += 86400
	}
	b := make([]byte, 0, 18)
	b = appendTwo(b, sec/3600)
	b = append(b, ':')
	b = appendTwo(b, sec/60%60)
	b = append(b, ':')
	b = appendTwo(b, sec%60)
	return string(b) + fracText(ns, digits)
}

func appendTwo(b []byte, v int64) []byte { return append(b, byte('0'+v/10), byte('0'+v%10)) }

func timestampText(sec, ns int64, digits int) string {
	return time.Unix(sec, 0).UTC().Format(time.DateTime) + fracText(ns, digits)
}

func uuidText(b []byte) string {
	if len(b) != 16 {
		return hexText(b)
	}
	h := hex.EncodeToString(b)
	return h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

// intervalText shows an INTERVAL (months, days and milliseconds) as an ISO
// 8601 duration.
func intervalText(b []byte) string {
	if len(b) != 12 {
		return hexText(b)
	}
	months := binary.LittleEndian.Uint32(b)
	days := binary.LittleEndian.Uint32(b[4:])
	ms := binary.LittleEndian.Uint32(b[8:])
	var s strings.Builder
	s.WriteString("P")
	if y := months / 12; y > 0 {
		s.WriteString(strconv.FormatUint(uint64(y), 10) + "Y")
	}
	if m := months % 12; m > 0 {
		s.WriteString(strconv.FormatUint(uint64(m), 10) + "M")
	}
	if days > 0 {
		s.WriteString(strconv.FormatUint(uint64(days), 10) + "D")
	}
	if ms > 0 || s.Len() == 1 {
		s.WriteString("T")
		sec := ms / 1000
		if h := sec / 3600; h > 0 {
			s.WriteString(strconv.FormatUint(uint64(h), 10) + "H")
		}
		if m := sec / 60 % 60; m > 0 {
			s.WriteString(strconv.FormatUint(uint64(m), 10) + "M")
		}
		if sec%60 > 0 || ms%1000 > 0 || s.Len() == 2 {
			s.WriteString(strconv.FormatUint(uint64(sec%60), 10))
			if f := ms % 1000; f > 0 {
				s.WriteString(strings.TrimRight(fracText(int64(f)*1e6, 3), "0"))
			}
			s.WriteString("S")
		}
	}
	return s.String()
}

// maxHexBytes bounds the bytes shown of a byte array in hexadecimal.
const maxHexBytes = 32

func hexText(b []byte) string {
	if len(b) <= maxHexBytes {
		return "0x" + hex.EncodeToString(b)
	}
	return "0x" + hex.EncodeToString(b[:maxHexBytes]) + "… (" + strconv.Itoa(len(b)) + " bytes)"
}

// binaryText shows a byte array without a type: as text when it is text
// (strings written without their annotation), else in hexadecimal.
func binaryText(b []byte) string {
	if len(b) > 0 && utf8.Valid(b) {
		text := true
		for _, c := range b {
			if c < 0x20 && c != '\t' && c != '\n' && c != '\r' || c == 0x7f {
				text = false
				break
			}
		}
		if text {
			return string(b)
		}
	}
	return hexText(b)
}

//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"math"
)

func F___math_divzero(m *base.Module, l0 int32) float64 {
	var v4 float64
	_ = v4
	if l0 != 0 {
		v4 = math.Float64frombits(uint64(0xfff0000000000000))
	} else {
		v4 = math.Float64frombits(uint64(0x7ff0000000000000))
	}
	return v4
}
func F___math_divzerof(m *base.Module, l0 int32) float32 {
	var v4 float32
	_ = v4
	if l0 != 0 {
		v4 = math.Float32frombits(uint32(0xff800000))
	} else {
		v4 = math.Float32frombits(uint32(0x7f800000))
	}
	return v4
}
func F___math_invalid(m *base.Module, l0 float64) float64 {
	var v2 float64
	_ = v2
	v2 = base.F64_sub(l0, l0)
	return base.F64_div(v2, v2)
}
func F___math_invalidf(m *base.Module, l0 float32) float32 {
	var v2 float32
	_ = v2
	v2 = base.F32_sub(l0, l0)
	return base.F32_div(v2, v2)
}
func F___math_oflow(m *base.Module, l0 int32) float64 {
	var v2 float64
	_ = v2
	var v4 float64
	_ = v4
	v2 = float64(3.105036184601418e+231)
	if l0 != 0 {
		v4 = base.F64_neg(v2)
	} else {
		v4 = v2
	}
	return base.F64_mul(v4, v2)
}
func F___math_oflowf(m *base.Module, l0 int32) float32 {
	var v2 float32
	_ = v2
	var v4 float32
	_ = v4
	v2 = float32(1.5845633e+29)
	if l0 != 0 {
		v4 = base.F32_neg(v2)
	} else {
		v4 = v2
	}
	return base.F32_mul(v4, v2)
}
func F___math_uflow(m *base.Module, l0 int32) float64 {
	var v2 float64
	_ = v2
	var v4 float64
	_ = v4
	v2 = float64(1.2882297539194267e-231)
	if l0 != 0 {
		v4 = base.F64_neg(v2)
	} else {
		v4 = v2
	}
	return base.F64_mul(v4, v2)
}
func F___math_uflowf(m *base.Module, l0 int32) float32 {
	var v2 float32
	_ = v2
	var v4 float32
	_ = v4
	v2 = float32(2.524355e-29)
	if l0 != 0 {
		v4 = base.F32_neg(v2)
	} else {
		v4 = v2
	}
	return base.F32_mul(v4, v2)
}
func F___math_xflow(m *base.Module, l0 int32, l1 float64) float64 {
	var v4 float64
	_ = v4
	if l0 != 0 {
		v4 = base.F64_neg(l1)
	} else {
		v4 = l1
	}
	return base.F64_mul(v4, l1)
}
func F___math_xflowf(m *base.Module, l0 int32, l1 float32) float32 {
	var v4 float32
	_ = v4
	if l0 != 0 {
		v4 = base.F32_neg(l1)
	} else {
		v4 = l1
	}
	return base.F32_mul(v4, l1)
}

//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"math"
	"unsafe"
)

func F_End(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
	return
}
func F_Execute(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if v3 == int32(0) {
	} else {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		v8 = m.T0[v3].(func(*base.Module, int32, int32) int32)(m, v6, v7)
		mBase = m.M
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v9 | base.B2i32(v8 == int32(0))
	}
	return
}
func F_expf(m *base.Module, l0 float32) float32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v28 float32
	_ = v28
	var v36 float32
	_ = v36
	var v39 int32
	_ = v39
	var v40 float64
	_ = v40
	var v42 float64
	_ = v42
	var v44 float64
	_ = v44
	var v46 float64
	_ = v46
	var v47 float64
	_ = v47
	var v49 float64
	_ = v49
	var v52 float64
	_ = v52
	var v57 float64
	_ = v57
	var v62 int64
	_ = v62
	var v72 int64
	_ = v72
	var v77 float32
	_ = v77
	v12 = int32(base.Ui32(base.I32_reinterpret_f32(l0))>>(uint(int32(20))%32)) & int32(2047)
	if base.Ui32(v12) < base.Ui32(int32(1067)) {
		v39 = int32(0)
		v40 = *(*float64)(unsafe.Add(mBase, _c_F_expf[0]))
		v42 = *(*float64)(unsafe.Add(mBase, _c_F_expf[1]))
		v44 = base.F64_mul(v42, base.F64_promote_f32(l0))
		v46 = *(*float64)(unsafe.Add(mBase, _c_F_expf[2]))
		v47 = base.F64_add(v44, v46)
		v49 = base.F64_sub(v44, base.F64_sub(v47, v46))
		v52 = *(*float64)(unsafe.Add(mBase, _c_F_expf[3]))
		v57 = *(*float64)(unsafe.Add(mBase, _c_F_expf[4]))
		v62 = base.I64_reinterpret_f64(v47)
		v72 = *(*int64)(unsafe.Add(mBase, uint32(base.I32_wrap_i64(v62)&int32(31)<<(uint(int32(3))%32))+uint32(_c_F_expf[5])))
		v77 = base.F32_demote_f64(base.F64_mul(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v40, v49), v52), base.F64_mul(v49, v49)), base.F64_add(base.F64_mul(v57, v49), float64(1))), base.F64_reinterpret_i64(v62<<(uint(int64(47))%64)+v72)))
		return v77
	} else {
		if base.F32_eq(l0, math.Float32frombits(uint32(0xff800000))) != 0 {
			v77 = float32(0)
			return v77
		} else {
			if base.Ui32(v12) < base.Ui32(int32(2040)) {
				if base.F32_gt(l0, float32(88.72283)) == int32(0) {
					if base.F32_lt(l0, float32(-103.97208)) == int32(0) {
						v39 = int32(0)
						v40 = *(*float64)(unsafe.Add(mBase, _c_F_expf[0]))
						v42 = *(*float64)(unsafe.Add(mBase, _c_F_expf[1]))
						v44 = base.F64_mul(v42, base.F64_promote_f32(l0))
						v46 = *(*float64)(unsafe.Add(mBase, _c_F_expf[2]))
						v47 = base.F64_add(v44, v46)
						v49 = base.F64_sub(v44, base.F64_sub(v47, v46))
						v52 = *(*float64)(unsafe.Add(mBase, _c_F_expf[3]))
						v57 = *(*float64)(unsafe.Add(mBase, _c_F_expf[4]))
						v62 = base.I64_reinterpret_f64(v47)
						v72 = *(*int64)(unsafe.Add(mBase, uint32(base.I32_wrap_i64(v62)&int32(31)<<(uint(int32(3))%32))+uint32(_c_F_expf[5])))
						v77 = base.F32_demote_f64(base.F64_mul(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v40, v49), v52), base.F64_mul(v49, v49)), base.F64_add(base.F64_mul(v57, v49), float64(1))), base.F64_reinterpret_i64(v62<<(uint(int64(47))%64)+v72)))
						return v77
					} else {
						v36 = F___math_xflowf(m, int32(0), float32(2.524355e-29))
						mBase = m.M
						return v36
					}
				} else {
					v28 = F___math_xflowf(m, int32(0), float32(1.5845633e+29))
					mBase = m.M
					return v28
				}
			} else {
				return base.F32_add(l0, l0)
			}
		}
	}
}

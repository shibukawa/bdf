//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"math"
	"unsafe"
)

func F_InitSharpYuvSSE2(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	v2 = m.G2
	v3 = m.G127
	*(*int32)(unsafe.Add(mBase, uint32(v3))) = v2 + int32(337)
	v7 = m.G128
	*(*int32)(unsafe.Add(mBase, uint32(v7))) = v2 + int32(338)
	v11 = m.G129
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v2 + int32(339)
	return
}
func F_SharpYuvConvert(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32, l15 int32) int32 {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v27 int32
	_ = v27
	v18 = m.G0
	v19 = int32(16)
	v20 = v18 - v19
	m.G0 = v20
	*(*int32)(unsafe.Add(mBase, uint32(v20)+12)) = int32(13)
	*(*int32)(unsafe.Add(mBase, uint32(v20)+8)) = l15
	v27 = F_SharpYuvConvertWithOptions(m, l0, l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11, l12, l13, l14, v20+int32(8))
	mBase = m.M
	m.G0 = v20 + v19
	return v27
}
func F_SharpYuvGammaToLinear(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v45 int32
	_ = v45
	var v51 int32
	_ = v51
	var v56 float32
	_ = v56
	var v66 float32
	_ = v66
	var v77 float64
	_ = v77
	var v79 float32
	_ = v79
	var v82 float32
	_ = v82
	var v85 float64
	_ = v85
	var v87 float32
	_ = v87
	var v90 float32
	_ = v90
	var v93 float64
	_ = v93
	var v101 float32
	_ = v101
	var v112 float64
	_ = v112
	var v120 float32
	_ = v120
	var v123 float32
	_ = v123
	var v125 float32
	_ = v125
	var v128 float64
	_ = v128
	var v136 float32
	_ = v136
	var v139 float32
	_ = v139
	var v145 float64
	_ = v145
	var v159 float64
	_ = v159
	var v167 float32
	_ = v167
	var v178 float64
	_ = v178
	var v180 float32
	_ = v180
	var v187 float64
	_ = v187
	var v188 float32
	_ = v188
	var v190 float32
	_ = v190
	var v191 float32
	_ = v191
	var v194 float32
	_ = v194
	var v198 float32
	_ = v198
	var v199 float32
	_ = v199
	var v202 float32
	_ = v202
	var v206 float64
	_ = v206
	var v208 float32
	_ = v208
	var v211 float32
	_ = v211
	var v214 float64
	_ = v214
	var v227 float64
	_ = v227
	var v232 float32
	_ = v232
	var v243 int32
	_ = v243
	var v257 float32
	_ = v257
	var v263 float32
	_ = v263
	var v265 int32
	_ = v265
	var v266 float64
	_ = v266
	var v268 float64
	_ = v268
	var v270 float64
	_ = v270
	var v272 float64
	_ = v272
	var v273 float64
	_ = v273
	var v275 float64
	_ = v275
	var v278 float64
	_ = v278
	var v283 float64
	_ = v283
	var v288 int64
	_ = v288
	var v298 int64
	_ = v298
	var v303 float32
	_ = v303
	var v313 float32
	_ = v313
	var v320 float64
	_ = v320
	var v323 float32
	_ = v323
	var v325 float32
	_ = v325
	var v334 float32
	_ = v334
	var v342 int32
	_ = v342
	var v345 int32
	_ = v345
	if l2 != int32(13) {
		v51 = int32(-1)
		v56 = base.F32_div(base.F32_convert_i32_u(l0), base.F32_convert_i32_u(v51<<(uint(l1)%32)^v51))
		switch l2 + v51 {
		case 0, 5, 13, 14:
			if base.F32_lt(v56, float32(0.08124286)) == int32(0) {
				v66 = float32(1)
				if base.F32_lt(v56, v66) == int32(0) {
					v323 = v66
				} else {
					v77 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.09929682)), float32(1.0992968))), float64(2.222222328186035))
					mBase = m.M
					v323 = base.F32_demote_f64(v77)
				}
			} else {
				v323 = base.F32_div(v56, float32(4.5))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		default:
			v323 = float32(0)
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 3:
			v79 = float32(1)
			if base.F32_gt(v56, v79) != 0 {
				v82 = v79
			} else {
				v82 = v56
			}
			v85 = F_pow(m, base.F64_promote_f32(v82), float64(2.200000047683716))
			mBase = m.M
			v323 = base.F32_demote_f64(v85)
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 4:
			v87 = float32(1)
			if base.F32_gt(v56, v87) != 0 {
				v90 = v87
			} else {
				v90 = v56
			}
			v93 = F_pow(m, base.F64_promote_f32(v90), float64(2.799999952316284))
			mBase = m.M
			v323 = base.F32_demote_f64(v93)
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 6:
			if base.F32_lt(v56, float32(0.09128634)) == int32(0) {
				v101 = float32(1)
				if base.F32_lt(v56, v101) == int32(0) {
					v323 = v101
				} else {
					v112 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.1115722)), float32(1.1115721))), float64(2.222222328186035))
					mBase = m.M
					v323 = base.F32_demote_f64(v112)
				}
			} else {
				v323 = base.F32_mul(v56, float32(0.25))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 7:
			v345 = l0
			return v345
		case 8:
			if base.F32_le(v56, float32(0)) == int32(0) {
				v120 = float32(1)
				if base.F32_lt(v56, v120) != 0 {
					v123 = v56
				} else {
					v123 = v120
				}
				v125 = base.F32_add(v123, float32(-1))
				v128 = F_pow(m, float64(10), base.F64_promote_f32(base.F32_add(v125, v125)))
				mBase = m.M
				v323 = base.F32_demote_f64(v128)
			} else {
				v323 = float32(0.005)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 9:
			if base.F32_le(v56, float32(0)) == int32(0) {
				v136 = float32(1)
				if base.F32_lt(v56, v136) != 0 {
					v139 = v56
				} else {
					v139 = v136
				}
				v145 = F_pow(m, float64(10), base.F64_promote_f32(base.F32_mul(base.F32_add(v139, float32(-1)), float32(2.5))))
				mBase = m.M
				v323 = base.F32_demote_f64(v145)
			} else {
				v323 = float32(0.0015811388)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 10:
			if base.F32_lt(v56, float32(0.08124286)) == int32(0) {
				v159 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.09929682)), float32(1.0992968))), float64(2.222222328186035))
				mBase = m.M
				v323 = base.F32_demote_f64(v159)
			} else {
				v323 = base.F32_div(v56, float32(4.5))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 11:
			if base.F32_lt(v56, float32(0.08124286)) == int32(0) {
				v167 = float32(1)
				if base.F32_lt(v56, v167) == int32(0) {
					v323 = v167
				} else {
					v178 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.09929682)), float32(1.0992968))), float64(2.222222328186035))
					mBase = m.M
					v323 = base.F32_demote_f64(v178)
				}
			} else {
				v323 = base.F32_div(v56, float32(4.5))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 15:
			v180 = float32(0)
			if base.F32_gt(v56, v180) == int32(0) {
				v323 = v180
			} else {
				v187 = F_pow(m, base.F64_promote_f32(v56), float64(0.012683313339948654))
				mBase = m.M
				v188 = base.F32_demote_f64(v187)
				v190 = base.F32_add(v188, float32(-0.8359375))
				v191 = float32(0)
				if base.F32_gt(v190, v191) != 0 {
					v194 = v190
				} else {
					v194 = v191
				}
				v198 = base.F32_add(base.F32_mul(v188, float32(-18.6875)), float32(18.851562))
				v199 = float32(1.1754944e-38)
				if base.F32_gt(v198, v199) != 0 {
					v202 = v198
				} else {
					v202 = v199
				}
				v206 = F_pow(m, base.F64_promote_f32(base.F32_div(v194, v202)), float64(6.27258825302124))
				mBase = m.M
				v323 = base.F32_demote_f64(v206)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 16:
			v208 = float32(0)
			if base.F32_gt(v56, v208) != 0 {
				v211 = v56
			} else {
				v211 = v208
			}
			v214 = F_pow(m, base.F64_promote_f32(v211), float64(2.5999999046325684))
			mBase = m.M
			v323 = base.F32_div(base.F32_demote_f64(v214), float32(0.9165553))
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 17:
			if base.F32_le(v56, float32(0.5)) == int32(0) {
				v232 = base.F32_div(base.F32_add(v56, float32(-0.5599107)), float32(0.17883277))
				v243 = int32(base.Ui32(base.I32_reinterpret_f32(v232))>>(uint(int32(20))%32)) & int32(2047)
				if base.Ui32(v243) < base.Ui32(int32(1067)) {
					v265 = int32(0)
					v266 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[0]))
					v268 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[1]))
					v270 = base.F64_mul(v268, base.F64_promote_f32(v232))
					v272 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[2]))
					v273 = base.F64_add(v270, v272)
					v275 = base.F64_sub(v270, base.F64_sub(v273, v272))
					v278 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[3]))
					v283 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[4]))
					v288 = base.I64_reinterpret_f64(v273)
					v298 = *(*int64)(unsafe.Add(mBase, uint32(base.I32_wrap_i64(v288)&int32(31)<<(uint(int32(3))%32))+uint32(_c_F_SharpYuvGammaToLinear[5])))
					v303 = base.F32_demote_f64(base.F64_mul(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v266, v275), v278), base.F64_mul(v275, v275)), base.F64_add(base.F64_mul(v283, v275), float64(1))), base.F64_reinterpret_i64(v288<<(uint(int64(47))%64)+v298)))
					v313 = v303
				} else {
					if base.F32_eq(v232, math.Float32frombits(uint32(0xff800000))) != 0 {
						v303 = float32(0)
						v313 = v303
					} else {
						if base.Ui32(v243) < base.Ui32(int32(2040)) {
							if base.F32_gt(v232, float32(88.72283)) == int32(0) {
								if base.F32_lt(v232, float32(-103.97208)) == int32(0) {
									v265 = int32(0)
									v266 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[0]))
									v268 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[1]))
									v270 = base.F64_mul(v268, base.F64_promote_f32(v232))
									v272 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[2]))
									v273 = base.F64_add(v270, v272)
									v275 = base.F64_sub(v270, base.F64_sub(v273, v272))
									v278 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[3]))
									v283 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[4]))
									v288 = base.I64_reinterpret_f64(v273)
									v298 = *(*int64)(unsafe.Add(mBase, uint32(base.I32_wrap_i64(v288)&int32(31)<<(uint(int32(3))%32))+uint32(_c_F_SharpYuvGammaToLinear[5])))
									v303 = base.F32_demote_f64(base.F64_mul(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v266, v275), v278), base.F64_mul(v275, v275)), base.F64_add(base.F64_mul(v283, v275), float64(1))), base.F64_reinterpret_i64(v288<<(uint(int64(47))%64)+v298)))
									v313 = v303
								} else {
									v263 = F___math_uflowf(m, int32(0))
									mBase = m.M
									v313 = v263
								}
							} else {
								v257 = F___math_oflowf(m, int32(0))
								mBase = m.M
								v313 = v257
							}
						} else {
							v313 = base.F32_add(v232, v232)
						}
					}
				}
				v320 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v313, float32(0.28466892)), float32(12))), float64(1.2000000476837158))
				mBase = m.M
				v323 = base.F32_demote_f64(v320)
			} else {
				v227 = F_pow(m, base.F64_promote_f32(base.F32_mul(base.F32_mul(v56, v56), float32(0.33333334))), float64(1.2000000476837158))
				mBase = m.M
				v323 = base.F32_demote_f64(v227)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		}
	} else {
		if int32(9) < l1 {
			v23 = m.G1
			v27 = l1 + int32(-10)
			v28 = int32(base.Ui32(l0) >> (uint(v27) % 32))
			v31 = v23 + int32(_a_F_SharpYuvGammaToLinear_0) + v28<<(uint(int32(2))%32)
			v34 = *(*int32)(unsafe.Add(mBase, uint32(v31+int32(4))))
			v35 = *(*int32)(unsafe.Add(mBase, uint32(v31)))
			if v27 != 0 {
				v45 = int32(1) << (uint(l1+int32(-11)) % 32)
			} else {
				v45 = int32(0)
			}
			return int32(base.Ui32((v34-v35)*(l0-v28<<(uint(v27)%32))+v45)>>(uint(v27)%32)) + v35
		} else {
			v12 = m.G1
			v21 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(_a_F_SharpYuvGammaToLinear_0)+l0<<(uint(int32(10)-l1)%32)<<(uint(int32(2))%32))))
			return v21
		}
	}
}
func F_SharpYuvGetConversionMatrix(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	if base.Ui32(int32(5)) < base.Ui32(l0) {
		v13 = int32(0)
	} else {
		v6 = m.G1
		v12 = *(*int32)(unsafe.Add(mBase, uint32(v6+int32(_a_F_SharpYuvGetConversionMatrix_0)+l0<<(uint(int32(2))%32))))
		v13 = v12
	}
	return v13
}
func F_SharpYuvInit(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	v2 = m.G130
	if l0 != v2 {
		v6 = m.G130
		*(*int32)(unsafe.Add(mBase, uint32(v6))) = l0
		v8 = l0
	} else {
		v4 = m.G130
		v5 = *(*int32)(unsafe.Add(mBase, uint32(v4)))
		v8 = v5
	}
	v9 = m.G1
	v12 = *(*int32)(unsafe.Add(mBase, uint32(v9)+uint32(_c_F_SharpYuvInit[0])))
	if v12 == v8 {
	} else {
		v14 = m.G1
		F_SharpYuvInitDsp(m)
		mBase = m.M
		F_SharpYuvInitGammaTables(m)
		mBase = m.M
		v19 = m.G130
		v20 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_SharpYuvInit[0]))) = v20
	}
	return
}
func F_SharpYuvInitDsp(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	v3 = m.G2
	v4 = m.G1
	*(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_SharpYuvInitDsp[0]))) = v3 + int32(340)
	*(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_SharpYuvInitDsp[1]))) = v3 + int32(341)
	*(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_SharpYuvInitDsp[2]))) = v3 + int32(342)
	v20 = m.G130
	v21 = *(*int32)(unsafe.Add(mBase, uint32(v20)))
	if v21 == int32(0) {
	} else {
		v24 = int32(0)
		v25 = m.T0[v21].(func(*base.Module, int32) int32)(m, v24)
		mBase = m.M
		if v25 == v24 {
		} else {
			v29 = m.G2
			v30 = m.G127
			*(*int32)(unsafe.Add(mBase, uint32(v30))) = v29 + int32(337)
			v34 = m.G128
			*(*int32)(unsafe.Add(mBase, uint32(v34))) = v29 + int32(338)
			v38 = m.G129
			*(*int32)(unsafe.Add(mBase, uint32(v38))) = v29 + int32(339)
		}
	}
	return
}
func F_SharpYuvInitGammaTables(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v12 float64
	_ = v12
	var v13 int32
	_ = v13
	var v18 float64
	_ = v18
	var v30 float64
	_ = v30
	var v31 float64
	_ = v31
	var v32 int32
	_ = v32
	var v41 float64
	_ = v41
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v64 float64
	_ = v64
	var v65 int32
	_ = v65
	var v70 float64
	_ = v70
	var v78 float64
	_ = v78
	var v83 float64
	_ = v83
	var v84 int32
	_ = v84
	var v93 float64
	_ = v93
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v116 int32
	_ = v116
	v6 = m.G1
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v6)+uint32(_c_F_SharpYuvInitGammaTables[0])))
	if v9 != 0 {
	} else {
		v12 = float64(0)
		v13 = int32(-4100)
		for {
			v18 = base.F64_mul(v12, float64(0.0009765625))
			if base.F64_le(v18, float64(0.08124285829863151)) == int32(0) {
				v30 = F_pow(m, base.F64_mul(base.F64_add(v18, float64(0.09929682680944)), float64(0.909672415686275)), float64(2.2222222222222223))
				mBase = m.M
				v31 = v30
			} else {
				v31 = base.F64_div(v18, float64(4.5))
			}
			v32 = m.G1
			v41 = base.F64_add(base.F64_mul(v31, float64(65536)), float64(0.5))
			if base.F64_lt(v41, float64(4.294967296e+09))&base.F64_ge(v41, float64(0)) == int32(0) {
				v51 = int32(0)
			} else {
				v49 = base.I32_trunc_f64_u(v41)
				v51 = v49
			}
			*(*int32)(unsafe.Add(mBase, uint32(v32+int32(_a_F_SharpYuvInitGammaTables_0)+v13)+uint32(_c_F_SharpYuvInitGammaTables[1]))) = v51
			v56 = v13 + int32(4)
			if v56 != 0 {
				v12 = base.F64_add(v12, float64(1))
				v13 = v56
				continue
			} else {
				break
			}
			break
		}
		v57 = m.G1
		v60 = *(*int32)(unsafe.Add(mBase, uint32(v57)+uint32(_c_F_SharpYuvInitGammaTables[2])))
		*(*int32)(unsafe.Add(mBase, uint32(v57)+uint32(_c_F_SharpYuvInitGammaTables[3]))) = v60
		v64 = float64(0)
		v65 = int32(-2052)
		for {
			v70 = base.F64_mul(v64, float64(0.001953125))
			if base.F64_le(v70, float64(0.018053968510807)) == int32(0) {
				v78 = F_pow(m, v70, float64(0.44999999999999996))
				mBase = m.M
				v83 = base.F64_add(base.F64_mul(v78, float64(1.09929682680944)), float64(-0.09929682680944))
			} else {
				v83 = base.F64_mul(v70, float64(4.5))
			}
			v84 = m.G1
			v93 = base.F64_add(base.F64_mul(v83, float64(65536)), float64(0.5))
			if base.F64_lt(v93, float64(4.294967296e+09))&base.F64_ge(v93, float64(0)) == int32(0) {
				v103 = int32(0)
			} else {
				v101 = base.I32_trunc_f64_u(v93)
				v103 = v101
			}
			*(*int32)(unsafe.Add(mBase, uint32(v84+int32(_a_F_SharpYuvInitGammaTables_1)+v65+int32(2052)))) = v103
			v108 = v65 + int32(4)
			if v108 != 0 {
				v64 = base.F64_add(v64, float64(1))
				v65 = v108
				continue
			} else {
				break
			}
			break
		}
		v109 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v109)+uint32(_c_F_SharpYuvInitGammaTables[0]))) = int32(1)
		v116 = *(*int32)(unsafe.Add(mBase, uint32(v109)+uint32(_c_F_SharpYuvInitGammaTables[4])))
		*(*int32)(unsafe.Add(mBase, uint32(v109)+uint32(_c_F_SharpYuvInitGammaTables[5]))) = v116
	}
	return
}
func F_SharpYuvLinearToGamma(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v46 float32
	_ = v46
	var v56 float32
	_ = v56
	var v63 float64
	_ = v63
	var v69 float32
	_ = v69
	var v72 float32
	_ = v72
	var v75 float64
	_ = v75
	var v77 float32
	_ = v77
	var v80 float32
	_ = v80
	var v83 float64
	_ = v83
	var v91 float32
	_ = v91
	var v98 float64
	_ = v98
	var v107 float32
	_ = v107
	var v110 float32
	_ = v110
	var v111 float64
	_ = v111
	var v122 int64
	_ = v122
	var v137 int32
	_ = v137
	var v139 int64
	_ = v139
	var v148 int64
	_ = v148
	var v153 int64
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v158 int32
	_ = v158
	var v162 float64
	_ = v162
	var v164 float64
	_ = v164
	var v177 float64
	_ = v177
	var v180 float64
	_ = v180
	var v185 float64
	_ = v185
	var v186 float64
	_ = v186
	var v187 float64
	_ = v187
	var v188 float64
	_ = v188
	var v193 float64
	_ = v193
	var v194 float64
	_ = v194
	var v195 float64
	_ = v195
	var v220 float64
	_ = v220
	var v232 float64
	_ = v232
	var v254 float64
	_ = v254
	var v263 float32
	_ = v263
	var v266 float32
	_ = v266
	var v267 float64
	_ = v267
	var v278 int64
	_ = v278
	var v293 int32
	_ = v293
	var v295 int64
	_ = v295
	var v304 int64
	_ = v304
	var v309 int64
	_ = v309
	var v310 int32
	_ = v310
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
	var v318 float64
	_ = v318
	var v320 float64
	_ = v320
	var v333 float64
	_ = v333
	var v336 float64
	_ = v336
	var v341 float64
	_ = v341
	var v342 float64
	_ = v342
	var v343 float64
	_ = v343
	var v344 float64
	_ = v344
	var v349 float64
	_ = v349
	var v350 float64
	_ = v350
	var v351 float64
	_ = v351
	var v376 float64
	_ = v376
	var v388 float64
	_ = v388
	var v410 float64
	_ = v410
	var v424 float64
	_ = v424
	var v436 float32
	_ = v436
	var v443 float64
	_ = v443
	var v449 float32
	_ = v449
	var v456 float64
	_ = v456
	var v457 float32
	_ = v457
	var v469 float64
	_ = v469
	var v471 float32
	_ = v471
	var v474 float32
	_ = v474
	var v479 float64
	_ = v479
	var v481 float32
	_ = v481
	var v484 float64
	_ = v484
	var v485 float32
	_ = v485
	var v498 float32
	_ = v498
	var v503 int32
	_ = v503
	var v509 int32
	_ = v509
	var v511 float32
	_ = v511
	var v518 float32
	_ = v518
	var v524 int32
	_ = v524
	var v526 int32
	_ = v526
	var v527 float64
	_ = v527
	var v529 int32
	_ = v529
	var v538 int32
	_ = v538
	var v541 float64
	_ = v541
	var v544 float64
	_ = v544
	var v545 float64
	_ = v545
	var v548 float64
	_ = v548
	var v551 float64
	_ = v551
	var v559 float64
	_ = v559
	var v563 float64
	_ = v563
	var v568 float32
	_ = v568
	var v578 float32
	_ = v578
	var v584 float32
	_ = v584
	var v585 int32
	_ = v585
	var v590 float32
	_ = v590
	var v599 float32
	_ = v599
	var v607 int32
	_ = v607
	var v609 int32
	_ = v609
	if l2 != int32(13) {
		v46 = base.F32_div(base.F32_convert_i32_u(l0), float32(65535))
		switch l2 + int32(-1) {
		case 0, 5, 13, 14:
			if base.F32_lt(v46, float32(0.01805397)) == int32(0) {
				v56 = float32(1)
				if base.F32_lt(v46, v56) == int32(0) {
					v584 = v56
				} else {
					v63 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
					mBase = m.M
					v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v63), float32(1.0992968)), float32(-0.09929682))
				}
			} else {
				v584 = base.F32_mul(v46, float32(4.5))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		default:
			v584 = float32(0)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 3:
			v69 = float32(1)
			if base.F32_gt(v46, v69) != 0 {
				v72 = v69
			} else {
				v72 = v46
			}
			v75 = F_pow(m, base.F64_promote_f32(v72), float64(0.45454543828964233))
			mBase = m.M
			v584 = base.F32_demote_f64(v75)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 4:
			v77 = float32(1)
			if base.F32_gt(v46, v77) != 0 {
				v80 = v77
			} else {
				v80 = v46
			}
			v83 = F_pow(m, base.F64_promote_f32(v80), float64(0.3571428656578064))
			mBase = m.M
			v584 = base.F32_demote_f64(v83)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 6:
			if base.F32_lt(v46, float32(0.022821585)) == int32(0) {
				v91 = float32(1)
				if base.F32_lt(v46, v91) == int32(0) {
					v584 = v91
				} else {
					v98 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
					mBase = m.M
					v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v98), float32(1.1115721)), float32(-0.1115722))
				}
			} else {
				v584 = base.F32_mul(v46, float32(4))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 7:
			v609 = l0
		case 8:
			if base.F32_lt(v46, float32(0.01)) != 0 {
				v584 = float32(0)
			} else {
				v107 = float32(1)
				if base.F32_lt(v46, v107) != 0 {
					v110 = v46
				} else {
					v110 = v107
				}
				v111 = base.F64_promote_f32(v110)
				v122 = base.I64_reinterpret_f64(v111)
				if int64(4503599627370495) < v122 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v122) {
						v232 = v111
						v254 = v232
					} else {
						v137 = int32(-1023)
						v139 = int64(base.Ui64(v122) >> (uint(int64(32)) % 64))
						if v139 == int64(1072693248) {
							if base.I32_wrap_i64(v122) != 0 {
								v153 = v122
								v154 = v137
								v156 = int32(1072693248)
								v158 = v156 + int32(614242)
								v162 = base.F64_convert_i32_s(v154 + int32(base.Ui32(v158)>>(uint(int32(20))%32)))
								v164 = base.F64_mul(v162, float64(0.30102999566361177))
								v177 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v158&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v153&int64(4294967295)), float64(-1))
								v180 = base.F64_mul(v177, base.F64_mul(v177, float64(0.5)))
								v185 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v177, v180)) & int64(-4294967296))
								v186 = float64(0.4342944818781689)
								v187 = base.F64_mul(v185, v186)
								v188 = base.F64_add(v164, v187)
								v193 = base.F64_div(v177, base.F64_add(v177, float64(2)))
								v194 = base.F64_mul(v193, v193)
								v195 = base.F64_mul(v194, v194)
								v220 = base.F64_add(base.F64_mul(v193, base.F64_add(v180, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v194, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v177, v185), v180))
								v232 = base.F64_add(v188, base.F64_add(base.F64_add(v187, base.F64_sub(v164, v188)), base.F64_add(base.F64_mul(v220, v186), base.F64_add(base.F64_mul(v162, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v220, v185), float64(2.5082946711645275e-11))))))
								v254 = v232
							} else {
								v254 = float64(0)
							}
						} else {
							v153 = v122
							v154 = v137
							v156 = base.I32_wrap_i64(v139)
							v158 = v156 + int32(614242)
							v162 = base.F64_convert_i32_s(v154 + int32(base.Ui32(v158)>>(uint(int32(20))%32)))
							v164 = base.F64_mul(v162, float64(0.30102999566361177))
							v177 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v158&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v153&int64(4294967295)), float64(-1))
							v180 = base.F64_mul(v177, base.F64_mul(v177, float64(0.5)))
							v185 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v177, v180)) & int64(-4294967296))
							v186 = float64(0.4342944818781689)
							v187 = base.F64_mul(v185, v186)
							v188 = base.F64_add(v164, v187)
							v193 = base.F64_div(v177, base.F64_add(v177, float64(2)))
							v194 = base.F64_mul(v193, v193)
							v195 = base.F64_mul(v194, v194)
							v220 = base.F64_add(base.F64_mul(v193, base.F64_add(v180, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v194, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v177, v185), v180))
							v232 = base.F64_add(v188, base.F64_add(base.F64_add(v187, base.F64_sub(v164, v188)), base.F64_add(base.F64_mul(v220, v186), base.F64_add(base.F64_mul(v162, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v220, v185), float64(2.5082946711645275e-11))))))
							v254 = v232
						}
					}
				} else {
					if base.F64_ne(v111, float64(0)) != 0 {
						if int64(-1) < v122 {
							v148 = base.I64_reinterpret_f64(base.F64_mul(v111, float64(1.8014398509481984e+16)))
							v153 = v148
							v154 = int32(-1077)
							v156 = base.I32_wrap_i64(int64(base.Ui64(v148) >> (uint(int64(32)) % 64)))
							v158 = v156 + int32(614242)
							v162 = base.F64_convert_i32_s(v154 + int32(base.Ui32(v158)>>(uint(int32(20))%32)))
							v164 = base.F64_mul(v162, float64(0.30102999566361177))
							v177 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v158&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v153&int64(4294967295)), float64(-1))
							v180 = base.F64_mul(v177, base.F64_mul(v177, float64(0.5)))
							v185 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v177, v180)) & int64(-4294967296))
							v186 = float64(0.4342944818781689)
							v187 = base.F64_mul(v185, v186)
							v188 = base.F64_add(v164, v187)
							v193 = base.F64_div(v177, base.F64_add(v177, float64(2)))
							v194 = base.F64_mul(v193, v193)
							v195 = base.F64_mul(v194, v194)
							v220 = base.F64_add(base.F64_mul(v193, base.F64_add(v180, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v194, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v177, v185), v180))
							v232 = base.F64_add(v188, base.F64_add(base.F64_add(v187, base.F64_sub(v164, v188)), base.F64_add(base.F64_mul(v220, v186), base.F64_add(base.F64_mul(v162, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v220, v185), float64(2.5082946711645275e-11))))))
							v254 = v232
						} else {
							v254 = base.F64_div(base.F64_sub(v111, v111), float64(0))
						}
					} else {
						v254 = base.F64_div(float64(-1), base.F64_mul(v111, v111))
					}
				}
				v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v254), float32(0.5)), float32(1))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 9:
			if base.F32_lt(v46, float32(0.0031622776)) != 0 {
				v584 = float32(0)
			} else {
				v263 = float32(1)
				if base.F32_lt(v46, v263) != 0 {
					v266 = v46
				} else {
					v266 = v263
				}
				v267 = base.F64_promote_f32(v266)
				v278 = base.I64_reinterpret_f64(v267)
				if int64(4503599627370495) < v278 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v278) {
						v388 = v267
						v410 = v388
					} else {
						v293 = int32(-1023)
						v295 = int64(base.Ui64(v278) >> (uint(int64(32)) % 64))
						if v295 == int64(1072693248) {
							if base.I32_wrap_i64(v278) != 0 {
								v309 = v278
								v310 = v293
								v312 = int32(1072693248)
								v314 = v312 + int32(614242)
								v318 = base.F64_convert_i32_s(v310 + int32(base.Ui32(v314)>>(uint(int32(20))%32)))
								v320 = base.F64_mul(v318, float64(0.30102999566361177))
								v333 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v314&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v309&int64(4294967295)), float64(-1))
								v336 = base.F64_mul(v333, base.F64_mul(v333, float64(0.5)))
								v341 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v333, v336)) & int64(-4294967296))
								v342 = float64(0.4342944818781689)
								v343 = base.F64_mul(v341, v342)
								v344 = base.F64_add(v320, v343)
								v349 = base.F64_div(v333, base.F64_add(v333, float64(2)))
								v350 = base.F64_mul(v349, v349)
								v351 = base.F64_mul(v350, v350)
								v376 = base.F64_add(base.F64_mul(v349, base.F64_add(v336, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v350, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v333, v341), v336))
								v388 = base.F64_add(v344, base.F64_add(base.F64_add(v343, base.F64_sub(v320, v344)), base.F64_add(base.F64_mul(v376, v342), base.F64_add(base.F64_mul(v318, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v376, v341), float64(2.5082946711645275e-11))))))
								v410 = v388
							} else {
								v410 = float64(0)
							}
						} else {
							v309 = v278
							v310 = v293
							v312 = base.I32_wrap_i64(v295)
							v314 = v312 + int32(614242)
							v318 = base.F64_convert_i32_s(v310 + int32(base.Ui32(v314)>>(uint(int32(20))%32)))
							v320 = base.F64_mul(v318, float64(0.30102999566361177))
							v333 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v314&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v309&int64(4294967295)), float64(-1))
							v336 = base.F64_mul(v333, base.F64_mul(v333, float64(0.5)))
							v341 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v333, v336)) & int64(-4294967296))
							v342 = float64(0.4342944818781689)
							v343 = base.F64_mul(v341, v342)
							v344 = base.F64_add(v320, v343)
							v349 = base.F64_div(v333, base.F64_add(v333, float64(2)))
							v350 = base.F64_mul(v349, v349)
							v351 = base.F64_mul(v350, v350)
							v376 = base.F64_add(base.F64_mul(v349, base.F64_add(v336, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v350, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v333, v341), v336))
							v388 = base.F64_add(v344, base.F64_add(base.F64_add(v343, base.F64_sub(v320, v344)), base.F64_add(base.F64_mul(v376, v342), base.F64_add(base.F64_mul(v318, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v376, v341), float64(2.5082946711645275e-11))))))
							v410 = v388
						}
					}
				} else {
					if base.F64_ne(v267, float64(0)) != 0 {
						if int64(-1) < v278 {
							v304 = base.I64_reinterpret_f64(base.F64_mul(v267, float64(1.8014398509481984e+16)))
							v309 = v304
							v310 = int32(-1077)
							v312 = base.I32_wrap_i64(int64(base.Ui64(v304) >> (uint(int64(32)) % 64)))
							v314 = v312 + int32(614242)
							v318 = base.F64_convert_i32_s(v310 + int32(base.Ui32(v314)>>(uint(int32(20))%32)))
							v320 = base.F64_mul(v318, float64(0.30102999566361177))
							v333 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v314&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v309&int64(4294967295)), float64(-1))
							v336 = base.F64_mul(v333, base.F64_mul(v333, float64(0.5)))
							v341 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v333, v336)) & int64(-4294967296))
							v342 = float64(0.4342944818781689)
							v343 = base.F64_mul(v341, v342)
							v344 = base.F64_add(v320, v343)
							v349 = base.F64_div(v333, base.F64_add(v333, float64(2)))
							v350 = base.F64_mul(v349, v349)
							v351 = base.F64_mul(v350, v350)
							v376 = base.F64_add(base.F64_mul(v349, base.F64_add(v336, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v350, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v333, v341), v336))
							v388 = base.F64_add(v344, base.F64_add(base.F64_add(v343, base.F64_sub(v320, v344)), base.F64_add(base.F64_mul(v376, v342), base.F64_add(base.F64_mul(v318, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v376, v341), float64(2.5082946711645275e-11))))))
							v410 = v388
						} else {
							v410 = base.F64_div(base.F64_sub(v267, v267), float64(0))
						}
					} else {
						v410 = base.F64_div(float64(-1), base.F64_mul(v267, v267))
					}
				}
				v584 = base.F32_add(base.F32_div(base.F32_demote_f64(v410), float32(2.5)), float32(1))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 10:
			if base.F32_lt(v46, float32(0.01805397)) == int32(0) {
				v424 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
				mBase = m.M
				v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v424), float32(1.0992968)), float32(-0.09929682))
			} else {
				v584 = base.F32_mul(v46, float32(4.5))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 11:
			if base.F32_lt(v46, float32(0.01805397)) == int32(0) {
				v436 = float32(1)
				if base.F32_lt(v46, v436) == int32(0) {
					v584 = v436
				} else {
					v443 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
					mBase = m.M
					v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v443), float32(1.0992968)), float32(-0.09929682))
				}
			} else {
				v584 = base.F32_mul(v46, float32(4.5))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 15:
			v449 = float32(0)
			if base.F32_gt(v46, v449) == int32(0) {
				v584 = v449
			} else {
				v456 = F_pow(m, base.F64_promote_f32(v46), float64(0.159423828125))
				mBase = m.M
				v457 = base.F32_demote_f64(v456)
				v469 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(base.F32_mul(v457, float32(18.851562)), float32(0.8359375)), base.F32_add(base.F32_mul(v457, float32(18.6875)), float32(1)))), float64(78.84375))
				mBase = m.M
				v584 = base.F32_demote_f64(v469)
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 16:
			v471 = float32(0)
			if base.F32_gt(v46, v471) != 0 {
				v474 = v46
			} else {
				v474 = v471
			}
			v479 = F_pow(m, base.F64_promote_f32(base.F32_mul(v474, float32(0.9165553))), float64(0.38461539149284363))
			mBase = m.M
			v584 = base.F32_demote_f64(v479)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 17:
			v481 = float32(0)
			v484 = F_pow(m, base.F64_promote_f32(v46), float64(0.8333333134651184))
			mBase = m.M
			v485 = base.F32_demote_f64(v484)
			if base.F32_lt(v485, v481) != 0 {
				v584 = v481
			} else {
				if base.F32_le(v485, float32(0.083333336)) == int32(0) {
					v498 = base.F32_add(base.F32_mul(v485, float32(12)), float32(-0.28466892))
					v503 = base.I32_reinterpret_f32(v498)
					if base.Ui32(int32(-2130706433)) < base.Ui32(v503+int32(-2139095040)) {
						v524 = v503
						v526 = int32(0)
						v527 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[0]))
						v529 = v524 + int32(-1060306944)
						v538 = int32(base.Ui32(v529)>>(uint(int32(15))%32)) & int32(240)
						v541 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[1])))
						v544 = base.F64_add(base.F64_mul(base.F64_promote_f32(base.F32_reinterpret_i32(v524-v529&int32(-8388608))), v541), float64(-1))
						v545 = base.F64_mul(v544, v544)
						v548 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[2]))
						v551 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[3]))
						v559 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[4]))
						v563 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[5])))
						v568 = base.F32_demote_f64(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v527, v545), base.F64_add(base.F64_mul(v548, v544), v551)), v545), base.F64_add(base.F64_add(base.F64_mul(base.F64_convert_i32_s(v529>>(uint(int32(23))%32)), v559), v563), v544)))
						v578 = v568
					} else {
						v509 = v503 << (uint(int32(1)) % 32)
						if v509 != 0 {
							if v503 == int32(2139095040) {
								v568 = v498
								v578 = v568
							} else {
								if v503 < int32(0) {
									v518 = F___math_invalidf(m, v498)
									mBase = m.M
									v578 = v518
								} else {
									if base.Ui32(v509) < base.Ui32(int32(-16777216)) {
										v524 = base.I32_reinterpret_f32(base.F32_mul(v498, float32(8.388608e+06))) + int32(-192937984)
										v526 = int32(0)
										v527 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[0]))
										v529 = v524 + int32(-1060306944)
										v538 = int32(base.Ui32(v529)>>(uint(int32(15))%32)) & int32(240)
										v541 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[1])))
										v544 = base.F64_add(base.F64_mul(base.F64_promote_f32(base.F32_reinterpret_i32(v524-v529&int32(-8388608))), v541), float64(-1))
										v545 = base.F64_mul(v544, v544)
										v548 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[2]))
										v551 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[3]))
										v559 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[4]))
										v563 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[5])))
										v568 = base.F32_demote_f64(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v527, v545), base.F64_add(base.F64_mul(v548, v544), v551)), v545), base.F64_add(base.F64_add(base.F64_mul(base.F64_convert_i32_s(v529>>(uint(int32(23))%32)), v559), v563), v544)))
										v578 = v568
									} else {
										v518 = F___math_invalidf(m, v498)
										mBase = m.M
										v578 = v518
									}
								}
							}
						} else {
							v511 = F___math_divzerof(m, int32(1))
							mBase = m.M
							v578 = v511
						}
					}
					v584 = base.F32_add(base.F32_mul(v578, float32(0.17883277)), float32(0.5599107))
				} else {
					v584 = base.F32_sqrt(base.F32_mul(v485, float32(3)))
				}
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		}
	} else {
		v11 = m.G1
		v18 = v11 + int32(_a_F_SharpYuvLinearToGamma_0) + l0>>(uint(int32(7))%32)<<(uint(int32(2))%32)
		v21 = *(*int32)(unsafe.Add(mBase, uint32(v18+int32(4))))
		v22 = int32(16)
		v23 = v22 - l1
		v26 = l1 + int32(-16)
		v29 = base.B2i32(l1 < v22)
		if l1 < v22 {
			v30 = v21 >> (uint(v23) % 32)
		} else {
			v30 = v21 << (uint(v26) % 32)
		}
		v31 = *(*int32)(unsafe.Add(mBase, uint32(v18)))
		if l1 < v22 {
			v34 = v31 >> (uint(v23) % 32)
		} else {
			v34 = v31 << (uint(v26) % 32)
		}
		v609 = int32(base.Ui32((v30-v34)*(l0&int32(127))+int32(64))>>(uint(int32(7))%32)) + v34
	}
	return v609 & int32(_a_F_SharpYuvLinearToGamma_1)
}

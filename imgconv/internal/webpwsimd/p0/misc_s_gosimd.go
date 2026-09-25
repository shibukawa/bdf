//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_SSE16x16_C(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 int32
	_ = v3
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v18 int32
	_ = v18
	var v19 base.V128
	_ = v19
	var v24 base.V128
	_ = v24
	var v27 base.V128
	_ = v27
	var v29 base.V128
	_ = v29
	var v37 base.V128
	_ = v37
	var v40 base.V128
	_ = v40
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v58 base.V128
	_ = v58
	var v61 base.V128
	_ = v61
	var v64 base.V128
	_ = v64
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	v3 = int32(0)
	v12 = v3
	v13 = v3
	for {
		v18 = int32(0)
		v19 = base.Simd_g_v128_load(m, l0+v12, v18)
		v24 = base.Simd_g_v128_load(m, l1+v12, v18)
		v27 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v19)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v24)))
		v29 = base.Simd_g_const(&F_SSE16x16_C__k0)
		v37 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v19, base.Simd_g_const(&F_SSE16x16_C__k0), base.Simd_g_const(&F_SSE16x16_C__k1)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v24, v19, base.Simd_g_const(&F_SSE16x16_C__k0), base.Simd_g_const(&F_SSE16x16_C__k1)))))
		v40 = base.Simd_g_const(&F_SSE16x16_C__k2)
		v48 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v19, base.Simd_g_const(&F_SSE16x16_C__k2), base.Simd_g_const(&F_SSE16x16_C__k1)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v24, v19, base.Simd_g_const(&F_SSE16x16_C__k2), base.Simd_g_const(&F_SSE16x16_C__k1)))))
		v50 = base.Simd_g_const(&F_SSE16x16_C__k3)
		v58 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v19, base.Simd_g_const(&F_SSE16x16_C__k3), base.Simd_g_const(&F_SSE16x16_C__k1)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v24, v19, base.Simd_g_const(&F_SSE16x16_C__k3), base.Simd_g_const(&F_SSE16x16_C__k1)))))
		v61 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v27, v27), base.Simd_g_i32x4_mul(v37, v37)), base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v48, v48), base.Simd_g_i32x4_mul(v58, v58)))
		v64 = base.Simd_g_i32x4_add(v61, base.Simd_g_i8x16_shuffle2(v61, v61, base.Simd_g_const(&F_SSE16x16_C__k4), base.Simd_g_const(&F_SSE16x16_C__k1)))
		v70 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v64, base.Simd_g_i8x16_shuffle2(v64, v64, base.Simd_g_const(&F_SSE16x16_C__k5), base.Simd_g_const(&F_SSE16x16_C__k1)))) + v13
		v72 = v12 + int32(32)
		if v72 != int32(512) {
			v12 = v72
			v13 = v70
			continue
		} else {
			break
		}
		break
	}
	return v70
}

var F_SSE16x16_C__k0 = [2]uint64{0xb0a0908, 0x0}
var F_SSE16x16_C__k1 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SSE16x16_C__k2 = [2]uint64{0x7060504, 0x0}
var F_SSE16x16_C__k3 = [2]uint64{0xf0e0d0c, 0x0}
var F_SSE16x16_C__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_SSE16x16_C__k5 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_SSE16x16_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v13 int32
	_ = v13
	var v14 base.V128
	_ = v14
	var v16 base.V128
	_ = v16
	var v19 base.V128
	_ = v19
	var v21 base.V128
	_ = v21
	var v22 base.V128
	_ = v22
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v37 base.V128
	_ = v37
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v47 base.V128
	_ = v47
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v54 int32
	_ = v54
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v64 int32
	_ = v64
	var v65 base.V128
	_ = v65
	var v67 base.V128
	_ = v67
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v74 int32
	_ = v74
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v84 int32
	_ = v84
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v90 base.V128
	_ = v90
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v98 int32
	_ = v98
	var v99 base.V128
	_ = v99
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v106 base.V128
	_ = v106
	var v110 base.V128
	_ = v110
	var v115 base.V128
	_ = v115
	var v118 int32
	_ = v118
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v130 base.V128
	_ = v130
	var v135 base.V128
	_ = v135
	var v138 int32
	_ = v138
	var v139 base.V128
	_ = v139
	var v141 base.V128
	_ = v141
	var v144 base.V128
	_ = v144
	var v146 base.V128
	_ = v146
	var v150 base.V128
	_ = v150
	var v155 base.V128
	_ = v155
	var v158 int32
	_ = v158
	var v159 base.V128
	_ = v159
	var v161 base.V128
	_ = v161
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v170 base.V128
	_ = v170
	var v175 base.V128
	_ = v175
	var v178 int32
	_ = v178
	var v179 base.V128
	_ = v179
	var v181 base.V128
	_ = v181
	var v184 base.V128
	_ = v184
	var v186 base.V128
	_ = v186
	var v190 base.V128
	_ = v190
	var v195 base.V128
	_ = v195
	var v198 int32
	_ = v198
	var v199 base.V128
	_ = v199
	var v201 base.V128
	_ = v201
	var v204 base.V128
	_ = v204
	var v206 base.V128
	_ = v206
	var v210 base.V128
	_ = v210
	var v215 base.V128
	_ = v215
	var v218 int32
	_ = v218
	var v219 base.V128
	_ = v219
	var v221 base.V128
	_ = v221
	var v224 base.V128
	_ = v224
	var v226 base.V128
	_ = v226
	var v230 base.V128
	_ = v230
	var v235 base.V128
	_ = v235
	var v238 int32
	_ = v238
	var v239 base.V128
	_ = v239
	var v241 base.V128
	_ = v241
	var v244 base.V128
	_ = v244
	var v246 base.V128
	_ = v246
	var v250 base.V128
	_ = v250
	var v252 base.V128
	_ = v252
	v3 = base.Simd_g_const(&F_SSE16x16_SSE2__k0)
	v13 = int32(448)
	v14 = base.Simd_g_v128_load_rng(m, l1, v13, int32(0), int32(496))
	v16 = base.Simd_g_v128_load_rng(m, l0, v13, int32(0), int32(496))
	v19 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v14, v16), base.Simd_g_i8x16_sub_sat_u(v16, v14))
	v21 = base.Simd_g_const(&F_SSE16x16_SSE2__k1)
	v22 = base.Simd_g_i8x16_shuffle2(v19, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v24 = int32(384)
	v25 = base.Simd_g_v128_load_nc(m, l1, v24)
	v27 = base.Simd_g_v128_load_nc(m, l0, v24)
	v30 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v25, v27), base.Simd_g_i8x16_sub_sat_u(v27, v25))
	v32 = base.Simd_g_i8x16_shuffle2(v30, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v34 = int32(320)
	v35 = base.Simd_g_v128_load_nc(m, l1, v34)
	v37 = base.Simd_g_v128_load_nc(m, l0, v34)
	v40 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v35, v37), base.Simd_g_i8x16_sub_sat_u(v37, v35))
	v42 = base.Simd_g_i8x16_shuffle2(v40, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v44 = int32(256)
	v45 = base.Simd_g_v128_load_nc(m, l1, v44)
	v47 = base.Simd_g_v128_load_nc(m, l0, v44)
	v50 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v45, v47), base.Simd_g_i8x16_sub_sat_u(v47, v45))
	v52 = base.Simd_g_i8x16_shuffle2(v50, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v54 = int32(192)
	v55 = base.Simd_g_v128_load_nc(m, l1, v54)
	v57 = base.Simd_g_v128_load_nc(m, l0, v54)
	v60 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v55, v57), base.Simd_g_i8x16_sub_sat_u(v57, v55))
	v62 = base.Simd_g_i8x16_shuffle2(v60, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v64 = int32(128)
	v65 = base.Simd_g_v128_load_nc(m, l1, v64)
	v67 = base.Simd_g_v128_load_nc(m, l0, v64)
	v70 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v65, v67), base.Simd_g_i8x16_sub_sat_u(v67, v65))
	v72 = base.Simd_g_i8x16_shuffle2(v70, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v74 = int32(64)
	v75 = base.Simd_g_v128_load_nc(m, l1, v74)
	v77 = base.Simd_g_v128_load_nc(m, l0, v74)
	v80 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v75, v77), base.Simd_g_i8x16_sub_sat_u(v77, v75))
	v82 = base.Simd_g_i8x16_shuffle2(v80, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v84 = int32(0)
	v85 = base.Simd_g_v128_load_nc(m, l1, v84)
	v87 = base.Simd_g_v128_load_nc(m, l0, v84)
	v90 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v85, v87), base.Simd_g_i8x16_sub_sat_u(v87, v85))
	v92 = base.Simd_g_i8x16_shuffle2(v90, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v94 = base.Simd_g_const(&F_SSE16x16_SSE2__k4)
	v95 = base.Simd_g_i8x16_shuffle2(v90, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v98 = int32(32)
	v99 = base.Simd_g_v128_load_nc(m, l1, v98)
	v101 = base.Simd_g_v128_load_nc(m, l0, v98)
	v104 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v99, v101), base.Simd_g_i8x16_sub_sat_u(v101, v99))
	v106 = base.Simd_g_i8x16_shuffle2(v104, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v110 = base.Simd_g_i8x16_shuffle2(v104, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v115 = base.Simd_g_i8x16_shuffle2(v80, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v118 = int32(96)
	v119 = base.Simd_g_v128_load_nc(m, l1, v118)
	v121 = base.Simd_g_v128_load_nc(m, l0, v118)
	v124 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v119, v121), base.Simd_g_i8x16_sub_sat_u(v121, v119))
	v126 = base.Simd_g_i8x16_shuffle2(v124, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v130 = base.Simd_g_i8x16_shuffle2(v124, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v135 = base.Simd_g_i8x16_shuffle2(v70, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v138 = int32(160)
	v139 = base.Simd_g_v128_load_nc(m, l1, v138)
	v141 = base.Simd_g_v128_load_nc(m, l0, v138)
	v144 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v139, v141), base.Simd_g_i8x16_sub_sat_u(v141, v139))
	v146 = base.Simd_g_i8x16_shuffle2(v144, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v150 = base.Simd_g_i8x16_shuffle2(v144, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v155 = base.Simd_g_i8x16_shuffle2(v60, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v158 = int32(224)
	v159 = base.Simd_g_v128_load_nc(m, l1, v158)
	v161 = base.Simd_g_v128_load_nc(m, l0, v158)
	v164 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v159, v161), base.Simd_g_i8x16_sub_sat_u(v161, v159))
	v166 = base.Simd_g_i8x16_shuffle2(v164, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v170 = base.Simd_g_i8x16_shuffle2(v164, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v175 = base.Simd_g_i8x16_shuffle2(v50, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v178 = int32(288)
	v179 = base.Simd_g_v128_load_nc(m, l1, v178)
	v181 = base.Simd_g_v128_load_nc(m, l0, v178)
	v184 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v179, v181), base.Simd_g_i8x16_sub_sat_u(v181, v179))
	v186 = base.Simd_g_i8x16_shuffle2(v184, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v190 = base.Simd_g_i8x16_shuffle2(v184, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v195 = base.Simd_g_i8x16_shuffle2(v40, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v198 = int32(352)
	v199 = base.Simd_g_v128_load_nc(m, l1, v198)
	v201 = base.Simd_g_v128_load_nc(m, l0, v198)
	v204 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v199, v201), base.Simd_g_i8x16_sub_sat_u(v201, v199))
	v206 = base.Simd_g_i8x16_shuffle2(v204, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v210 = base.Simd_g_i8x16_shuffle2(v204, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v215 = base.Simd_g_i8x16_shuffle2(v30, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v218 = int32(416)
	v219 = base.Simd_g_v128_load_nc(m, l1, v218)
	v221 = base.Simd_g_v128_load_nc(m, l0, v218)
	v224 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v219, v221), base.Simd_g_i8x16_sub_sat_u(v221, v219))
	v226 = base.Simd_g_i8x16_shuffle2(v224, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v230 = base.Simd_g_i8x16_shuffle2(v224, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v235 = base.Simd_g_i8x16_shuffle2(v19, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v238 = int32(480)
	v239 = base.Simd_g_v128_load_nc(m, l1, v238)
	v241 = base.Simd_g_v128_load_nc(m, l0, v238)
	v244 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v239, v241), base.Simd_g_i8x16_sub_sat_u(v241, v239))
	v246 = base.Simd_g_i8x16_shuffle2(v244, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k2), base.Simd_g_const(&F_SSE16x16_SSE2__k3))
	v250 = base.Simd_g_i8x16_shuffle2(v244, v3, base.Simd_g_const(&F_SSE16x16_SSE2__k5), base.Simd_g_const(&F_SSE16x16_SSE2__k6))
	v252 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v22, v22), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v32, v32), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v42, v42), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v52, v52), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v62, v62), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v72, v72), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v82, v82), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v92, v92), base.Simd_g_i32x4_dot_i16x8_s(v95, v95)), base.Simd_g_i32x4_dot_i16x8_s(v106, v106)), base.Simd_g_i32x4_dot_i16x8_s(v110, v110))), base.Simd_g_i32x4_dot_i16x8_s(v115, v115)), base.Simd_g_i32x4_dot_i16x8_s(v126, v126)), base.Simd_g_i32x4_dot_i16x8_s(v130, v130))), base.Simd_g_i32x4_dot_i16x8_s(v135, v135)), base.Simd_g_i32x4_dot_i16x8_s(v146, v146)), base.Simd_g_i32x4_dot_i16x8_s(v150, v150))), base.Simd_g_i32x4_dot_i16x8_s(v155, v155)), base.Simd_g_i32x4_dot_i16x8_s(v166, v166)), base.Simd_g_i32x4_dot_i16x8_s(v170, v170))), base.Simd_g_i32x4_dot_i16x8_s(v175, v175)), base.Simd_g_i32x4_dot_i16x8_s(v186, v186)), base.Simd_g_i32x4_dot_i16x8_s(v190, v190))), base.Simd_g_i32x4_dot_i16x8_s(v195, v195)), base.Simd_g_i32x4_dot_i16x8_s(v206, v206)), base.Simd_g_i32x4_dot_i16x8_s(v210, v210))), base.Simd_g_i32x4_dot_i16x8_s(v215, v215)), base.Simd_g_i32x4_dot_i16x8_s(v226, v226)), base.Simd_g_i32x4_dot_i16x8_s(v230, v230))), base.Simd_g_i32x4_dot_i16x8_s(v235, v235)), base.Simd_g_i32x4_dot_i16x8_s(v246, v246)), base.Simd_g_i32x4_dot_i16x8_s(v250, v250))
	return base.Simd_g_i32x4_extract_lane_l3(v252) + base.Simd_g_i32x4_extract_lane_l2(v252) + base.Simd_g_i32x4_extract_lane_l1(v252) + base.Simd_g_i32x4_extract_lane_l0(v252)
}

var F_SSE16x16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SSE16x16_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_SSE16x16_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_SSE16x16_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_SSE16x16_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SSE16x16_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SSE16x16_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_SSE16x8_C(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 int32
	_ = v3
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v18 int32
	_ = v18
	var v19 base.V128
	_ = v19
	var v24 base.V128
	_ = v24
	var v27 base.V128
	_ = v27
	var v29 base.V128
	_ = v29
	var v37 base.V128
	_ = v37
	var v40 base.V128
	_ = v40
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v58 base.V128
	_ = v58
	var v61 base.V128
	_ = v61
	var v64 base.V128
	_ = v64
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	v3 = int32(0)
	v12 = v3
	v13 = v3
	for {
		v18 = int32(0)
		v19 = base.Simd_g_v128_load(m, l0+v12, v18)
		v24 = base.Simd_g_v128_load(m, l1+v12, v18)
		v27 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v19)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v24)))
		v29 = base.Simd_g_const(&F_SSE16x8_C__k0)
		v37 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v19, base.Simd_g_const(&F_SSE16x8_C__k0), base.Simd_g_const(&F_SSE16x8_C__k1)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v24, v19, base.Simd_g_const(&F_SSE16x8_C__k0), base.Simd_g_const(&F_SSE16x8_C__k1)))))
		v40 = base.Simd_g_const(&F_SSE16x8_C__k2)
		v48 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v19, base.Simd_g_const(&F_SSE16x8_C__k2), base.Simd_g_const(&F_SSE16x8_C__k1)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v24, v19, base.Simd_g_const(&F_SSE16x8_C__k2), base.Simd_g_const(&F_SSE16x8_C__k1)))))
		v50 = base.Simd_g_const(&F_SSE16x8_C__k3)
		v58 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v19, v19, base.Simd_g_const(&F_SSE16x8_C__k3), base.Simd_g_const(&F_SSE16x8_C__k1)))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v24, v19, base.Simd_g_const(&F_SSE16x8_C__k3), base.Simd_g_const(&F_SSE16x8_C__k1)))))
		v61 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v27, v27), base.Simd_g_i32x4_mul(v37, v37)), base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v48, v48), base.Simd_g_i32x4_mul(v58, v58)))
		v64 = base.Simd_g_i32x4_add(v61, base.Simd_g_i8x16_shuffle2(v61, v61, base.Simd_g_const(&F_SSE16x8_C__k4), base.Simd_g_const(&F_SSE16x8_C__k1)))
		v70 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v64, base.Simd_g_i8x16_shuffle2(v64, v64, base.Simd_g_const(&F_SSE16x8_C__k5), base.Simd_g_const(&F_SSE16x8_C__k1)))) + v13
		v72 = v12 + int32(32)
		if v72 != int32(256) {
			v12 = v72
			v13 = v70
			continue
		} else {
			break
		}
		break
	}
	return v70
}

var F_SSE16x8_C__k0 = [2]uint64{0xb0a0908, 0x0}
var F_SSE16x8_C__k1 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SSE16x8_C__k2 = [2]uint64{0x7060504, 0x0}
var F_SSE16x8_C__k3 = [2]uint64{0xf0e0d0c, 0x0}
var F_SSE16x8_C__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_SSE16x8_C__k5 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_SSE16x8_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v9 int32
	_ = v9
	var v10 base.V128
	_ = v10
	var v12 base.V128
	_ = v12
	var v15 base.V128
	_ = v15
	var v17 base.V128
	_ = v17
	var v18 base.V128
	_ = v18
	var v20 int32
	_ = v20
	var v21 base.V128
	_ = v21
	var v23 base.V128
	_ = v23
	var v26 base.V128
	_ = v26
	var v28 base.V128
	_ = v28
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v33 base.V128
	_ = v33
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v40 int32
	_ = v40
	var v41 base.V128
	_ = v41
	var v43 base.V128
	_ = v43
	var v46 base.V128
	_ = v46
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v54 int32
	_ = v54
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v66 base.V128
	_ = v66
	var v71 base.V128
	_ = v71
	var v74 int32
	_ = v74
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v86 base.V128
	_ = v86
	var v91 base.V128
	_ = v91
	var v94 int32
	_ = v94
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v106 base.V128
	_ = v106
	var v111 base.V128
	_ = v111
	var v114 int32
	_ = v114
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v120 base.V128
	_ = v120
	var v122 base.V128
	_ = v122
	var v126 base.V128
	_ = v126
	var v128 base.V128
	_ = v128
	v3 = base.Simd_g_const(&F_SSE16x8_SSE2__k0)
	v9 = int32(192)
	v10 = base.Simd_g_v128_load_rng(m, l1, v9, int32(0), int32(240))
	v12 = base.Simd_g_v128_load_rng(m, l0, v9, int32(0), int32(240))
	v15 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v10, v12), base.Simd_g_i8x16_sub_sat_u(v12, v10))
	v17 = base.Simd_g_const(&F_SSE16x8_SSE2__k1)
	v18 = base.Simd_g_i8x16_shuffle2(v15, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v20 = int32(128)
	v21 = base.Simd_g_v128_load_nc(m, l1, v20)
	v23 = base.Simd_g_v128_load_nc(m, l0, v20)
	v26 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v21, v23), base.Simd_g_i8x16_sub_sat_u(v23, v21))
	v28 = base.Simd_g_i8x16_shuffle2(v26, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v30 = int32(64)
	v31 = base.Simd_g_v128_load_nc(m, l1, v30)
	v33 = base.Simd_g_v128_load_nc(m, l0, v30)
	v36 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v31, v33), base.Simd_g_i8x16_sub_sat_u(v33, v31))
	v38 = base.Simd_g_i8x16_shuffle2(v36, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v40 = int32(0)
	v41 = base.Simd_g_v128_load_nc(m, l1, v40)
	v43 = base.Simd_g_v128_load_nc(m, l0, v40)
	v46 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v41, v43), base.Simd_g_i8x16_sub_sat_u(v43, v41))
	v48 = base.Simd_g_i8x16_shuffle2(v46, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v50 = base.Simd_g_const(&F_SSE16x8_SSE2__k4)
	v51 = base.Simd_g_i8x16_shuffle2(v46, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v54 = int32(32)
	v55 = base.Simd_g_v128_load_nc(m, l1, v54)
	v57 = base.Simd_g_v128_load_nc(m, l0, v54)
	v60 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v55, v57), base.Simd_g_i8x16_sub_sat_u(v57, v55))
	v62 = base.Simd_g_i8x16_shuffle2(v60, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v66 = base.Simd_g_i8x16_shuffle2(v60, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v71 = base.Simd_g_i8x16_shuffle2(v36, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v74 = int32(96)
	v75 = base.Simd_g_v128_load_nc(m, l1, v74)
	v77 = base.Simd_g_v128_load_nc(m, l0, v74)
	v80 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v75, v77), base.Simd_g_i8x16_sub_sat_u(v77, v75))
	v82 = base.Simd_g_i8x16_shuffle2(v80, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v86 = base.Simd_g_i8x16_shuffle2(v80, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v91 = base.Simd_g_i8x16_shuffle2(v26, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v94 = int32(160)
	v95 = base.Simd_g_v128_load_nc(m, l1, v94)
	v97 = base.Simd_g_v128_load_nc(m, l0, v94)
	v100 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v95, v97), base.Simd_g_i8x16_sub_sat_u(v97, v95))
	v102 = base.Simd_g_i8x16_shuffle2(v100, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v106 = base.Simd_g_i8x16_shuffle2(v100, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v111 = base.Simd_g_i8x16_shuffle2(v15, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v114 = int32(224)
	v115 = base.Simd_g_v128_load_nc(m, l1, v114)
	v117 = base.Simd_g_v128_load_nc(m, l0, v114)
	v120 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v115, v117), base.Simd_g_i8x16_sub_sat_u(v117, v115))
	v122 = base.Simd_g_i8x16_shuffle2(v120, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k2), base.Simd_g_const(&F_SSE16x8_SSE2__k3))
	v126 = base.Simd_g_i8x16_shuffle2(v120, v3, base.Simd_g_const(&F_SSE16x8_SSE2__k5), base.Simd_g_const(&F_SSE16x8_SSE2__k6))
	v128 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v18, v18), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v28, v28), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v38, v38), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v48, v48), base.Simd_g_i32x4_dot_i16x8_s(v51, v51)), base.Simd_g_i32x4_dot_i16x8_s(v62, v62)), base.Simd_g_i32x4_dot_i16x8_s(v66, v66))), base.Simd_g_i32x4_dot_i16x8_s(v71, v71)), base.Simd_g_i32x4_dot_i16x8_s(v82, v82)), base.Simd_g_i32x4_dot_i16x8_s(v86, v86))), base.Simd_g_i32x4_dot_i16x8_s(v91, v91)), base.Simd_g_i32x4_dot_i16x8_s(v102, v102)), base.Simd_g_i32x4_dot_i16x8_s(v106, v106))), base.Simd_g_i32x4_dot_i16x8_s(v111, v111)), base.Simd_g_i32x4_dot_i16x8_s(v122, v122)), base.Simd_g_i32x4_dot_i16x8_s(v126, v126))
	return base.Simd_g_i32x4_extract_lane_l3(v128) + base.Simd_g_i32x4_extract_lane_l2(v128) + base.Simd_g_i32x4_extract_lane_l1(v128) + base.Simd_g_i32x4_extract_lane_l0(v128)
}

var F_SSE16x8_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SSE16x8_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_SSE16x8_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_SSE16x8_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_SSE16x8_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SSE16x8_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SSE16x8_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_SSE4x4_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v9 base.V128
	_ = v9
	var v12 base.V128
	_ = v12
	var v15 base.V128
	_ = v15
	var v17 base.V128
	_ = v17
	var v22 base.V128
	_ = v22
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v33 base.V128
	_ = v33
	var v35 base.V128
	_ = v35
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	v3 = base.Simd_g_const(&F_SSE4x4_SSE2__k0)
	v5 = int32(64)
	v6 = base.Simd_g_v128_load64_zero(m, l0, v5)
	v7 = int32(96)
	v8 = base.Simd_g_v128_load64_zero(m, l0, v7)
	v9 = base.Simd_g_const(&F_SSE4x4_SSE2__k1)
	v12 = base.Simd_g_const(&F_SSE4x4_SSE2__k2)
	v15 = base.Simd_g_v128_load64_zero(m, l1, v5)
	v17 = base.Simd_g_v128_load64_zero(m, l1, v7)
	v22 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v6, v8, base.Simd_g_const(&F_SSE4x4_SSE2__k3), base.Simd_g_const(&F_SSE4x4_SSE2__k4)), v3, base.Simd_g_const(&F_SSE4x4_SSE2__k5), base.Simd_g_const(&F_SSE4x4_SSE2__k6)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v15, v17, base.Simd_g_const(&F_SSE4x4_SSE2__k3), base.Simd_g_const(&F_SSE4x4_SSE2__k4)), v3, base.Simd_g_const(&F_SSE4x4_SSE2__k5), base.Simd_g_const(&F_SSE4x4_SSE2__k6)))
	v24 = int32(0)
	v25 = base.Simd_g_v128_load64_zero(m, l0, v24)
	v26 = int32(32)
	v27 = base.Simd_g_v128_load64_zero(m, l0, v26)
	v33 = base.Simd_g_v128_load64_zero(m, l1, v24)
	v35 = base.Simd_g_v128_load64_zero(m, l1, v26)
	v40 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v25, v27, base.Simd_g_const(&F_SSE4x4_SSE2__k3), base.Simd_g_const(&F_SSE4x4_SSE2__k4)), v3, base.Simd_g_const(&F_SSE4x4_SSE2__k5), base.Simd_g_const(&F_SSE4x4_SSE2__k6)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v33, v35, base.Simd_g_const(&F_SSE4x4_SSE2__k3), base.Simd_g_const(&F_SSE4x4_SSE2__k4)), v3, base.Simd_g_const(&F_SSE4x4_SSE2__k5), base.Simd_g_const(&F_SSE4x4_SSE2__k6)))
	v42 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v22, v22), base.Simd_g_i32x4_dot_i16x8_s(v40, v40))
	return base.Simd_g_i32x4_extract_lane_l3(v42) + base.Simd_g_i32x4_extract_lane_l2(v42) + base.Simd_g_i32x4_extract_lane_l1(v42) + base.Simd_g_i32x4_extract_lane_l0(v42)
}

var F_SSE4x4_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SSE4x4_SSE2__k1 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_SSE4x4_SSE2__k2 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SSE4x4_SSE2__k3 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_SSE4x4_SSE2__k4 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_SSE4x4_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SSE4x4_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_SSE8x8_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v14 base.V128
	_ = v14
	var v16 int32
	_ = v16
	var v17 base.V128
	_ = v17
	var v21 base.V128
	_ = v21
	var v24 base.V128
	_ = v24
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v31 base.V128
	_ = v31
	var v34 base.V128
	_ = v34
	var v36 int32
	_ = v36
	var v37 base.V128
	_ = v37
	var v41 base.V128
	_ = v41
	var v44 base.V128
	_ = v44
	var v46 int32
	_ = v46
	var v47 base.V128
	_ = v47
	var v51 base.V128
	_ = v51
	var v54 base.V128
	_ = v54
	var v58 int32
	_ = v58
	var v59 base.V128
	_ = v59
	var v63 base.V128
	_ = v63
	var v66 base.V128
	_ = v66
	var v70 int32
	_ = v70
	var v71 base.V128
	_ = v71
	var v75 base.V128
	_ = v75
	var v78 base.V128
	_ = v78
	var v82 int32
	_ = v82
	var v83 base.V128
	_ = v83
	var v87 base.V128
	_ = v87
	var v90 base.V128
	_ = v90
	var v92 base.V128
	_ = v92
	v3 = base.Simd_g_const(&F_SSE8x8_SSE2__k0)
	v5 = int32(192)
	v6 = base.Simd_g_v128_load64_zero(m, l0, v5)
	v8 = base.Simd_g_const(&F_SSE8x8_SSE2__k1)
	v11 = base.Simd_g_v128_load64_zero(m, l1, v5)
	v14 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v6, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v11, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v16 = int32(128)
	v17 = base.Simd_g_v128_load64_zero(m, l0, v16)
	v21 = base.Simd_g_v128_load64_zero(m, l1, v16)
	v24 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v17, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v21, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v26 = int32(64)
	v27 = base.Simd_g_v128_load64_zero(m, l0, v26)
	v31 = base.Simd_g_v128_load64_zero(m, l1, v26)
	v34 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v27, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v31, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v36 = int32(0)
	v37 = base.Simd_g_v128_load64_zero(m, l0, v36)
	v41 = base.Simd_g_v128_load64_zero(m, l1, v36)
	v44 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v37, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v41, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v46 = int32(32)
	v47 = base.Simd_g_v128_load64_zero(m, l0, v46)
	v51 = base.Simd_g_v128_load64_zero(m, l1, v46)
	v54 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v47, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v51, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v58 = int32(96)
	v59 = base.Simd_g_v128_load64_zero(m, l0, v58)
	v63 = base.Simd_g_v128_load64_zero(m, l1, v58)
	v66 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v59, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v63, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v70 = int32(160)
	v71 = base.Simd_g_v128_load64_zero(m, l0, v70)
	v75 = base.Simd_g_v128_load64_zero(m, l1, v70)
	v78 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v71, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v75, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v82 = int32(224)
	v83 = base.Simd_g_v128_load64_zero(m, l0, v82)
	v87 = base.Simd_g_v128_load64_zero(m, l1, v82)
	v90 = base.Simd_g_i16x8_sub_sat_s(base.Simd_g_i8x16_shuffle2(v83, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v87, v3, base.Simd_g_const(&F_SSE8x8_SSE2__k2), base.Simd_g_const(&F_SSE8x8_SSE2__k3)))
	v92 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v14, v14), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v24, v24), base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v34, v34), base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v44, v44), base.Simd_g_i32x4_dot_i16x8_s(v54, v54))), base.Simd_g_i32x4_dot_i16x8_s(v66, v66))), base.Simd_g_i32x4_dot_i16x8_s(v78, v78))), base.Simd_g_i32x4_dot_i16x8_s(v90, v90))
	return base.Simd_g_i32x4_extract_lane_l3(v92) + base.Simd_g_i32x4_extract_lane_l2(v92) + base.Simd_g_i32x4_extract_lane_l1(v92) + base.Simd_g_i32x4_extract_lane_l0(v92)
}

var F_SSE8x8_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SSE8x8_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SSE8x8_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SSE8x8_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_SSIMGetClipped_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32) float64 {
	mBase := m.M
	_ = mBase
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v51 int32
	_ = v51
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int64
	_ = v64
	var v68 int64
	_ = v68
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v115 int32
	_ = v115
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v166 int32
	_ = v166
	var v172 int32
	_ = v172
	var v175 base.V128
	_ = v175
	var v176 int32
	_ = v176
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v200 int32
	_ = v200
	var v221 base.V128
	_ = v221
	var v222 base.V128
	_ = v222
	var v223 base.V128
	_ = v223
	var v224 base.V128
	_ = v224
	var v225 base.V128
	_ = v225
	var v226 base.V128
	_ = v226
	var v237 int32
	_ = v237
	var v238 base.V128
	_ = v238
	var v239 base.V128
	_ = v239
	var v242 base.V128
	_ = v242
	var v244 base.V128
	_ = v244
	var v245 base.V128
	_ = v245
	var v246 base.V128
	_ = v246
	var v249 base.V128
	_ = v249
	var v251 base.V128
	_ = v251
	var v252 base.V128
	_ = v252
	var v253 base.V128
	_ = v253
	var v255 base.V128
	_ = v255
	var v257 base.V128
	_ = v257
	var v259 base.V128
	_ = v259
	var v264 base.V128
	_ = v264
	var v266 int32
	_ = v266
	var v267 base.V128
	_ = v267
	var v269 base.V128
	_ = v269
	var v270 base.V128
	_ = v270
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v277 base.V128
	_ = v277
	var v282 int32
	_ = v282
	var v285 base.V128
	_ = v285
	var v290 int32
	_ = v290
	var v293 base.V128
	_ = v293
	var v298 int32
	_ = v298
	var v301 base.V128
	_ = v301
	var v306 int32
	_ = v306
	var v309 base.V128
	_ = v309
	var v314 int32
	_ = v314
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v325 int32
	_ = v325
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v345 int32
	_ = v345
	var v367 int32
	_ = v367
	var v371 int32
	_ = v371
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v393 int32
	_ = v393
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v415 int32
	_ = v415
	var v416 int32
	_ = v416
	var v418 int32
	_ = v418
	var v419 int32
	_ = v419
	var v420 int32
	_ = v420
	var v422 int32
	_ = v422
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v430 int32
	_ = v430
	var v433 int32
	_ = v433
	var v435 int32
	_ = v435
	var v437 int32
	_ = v437
	var v446 int32
	_ = v446
	var v463 int32
	_ = v463
	var v464 int32
	_ = v464
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v503 int64
	_ = v503
	var v504 int32
	_ = v504
	var v508 int64
	_ = v508
	var v509 int64
	_ = v509
	var v510 int64
	_ = v510
	var v511 int64
	_ = v511
	var v546 int64
	_ = v546
	var v547 int32
	_ = v547
	var v552 int64
	_ = v552
	var v554 int64
	_ = v554
	var v555 int64
	_ = v555
	var v556 int64
	_ = v556
	var v559 int64
	_ = v559
	var v560 int64
	_ = v560
	var v564 int64
	_ = v564
	var v566 int64
	_ = v566
	var v572 int64
	_ = v572
	var v588 float64
	_ = v588
	v9 = int64(0)
	v10 = int32(0)
	v51 = int32(3)
	if v51 < l5 {
		v54 = l5
	} else {
		v54 = v51
	}
	v56 = v54 + int32(-3)
	v58 = l5 + int32(3)
	v60 = l7 + int32(-1)
	if v58 < v60 {
		v62 = v58
	} else {
		v62 = v60
	}
	if v56 <= v62 {
		v68 = int64(0)
		v72 = int32(3)
		if v72 < l4 {
			v75 = l4
		} else {
			v75 = v72
		}
		v77 = v75 + int32(-3)
		v79 = l4 + int32(3)
		v81 = l6 + int32(-1)
		if v79 < v81 {
			v83 = v79
		} else {
			v83 = v81
		}
		if v83 < v77 {
			v503 = v9
			v504 = v10
			v508 = v68
			v509 = v68
			v510 = v68
			v511 = v68
		} else {
			v85 = m.G1
			v87 = v85 + int32(_a_F_SSIMGetClipped_C_0)
			v88 = int32(2)
			v89 = l4 << (uint(v88) % 32)
			v105 = v83 - v75 + int32(4)
			v107 = v105 & int32(-4)
			v109 = int32(0)
			v115 = v109
			v124 = v109
			v125 = v56
			v136 = l2 + v56*l3
			v137 = l0 + v56*l1
			v141 = v109
			v142 = v109
			v143 = v109
			v144 = v109
			for {
				v166 = m.G1
				v172 = *(*int32)(unsafe.Add(mBase, uint32((v125-l5)<<(uint(int32(2))%32)+(v166+int32(_a_F_SSIMGetClipped_C_0))+int32(12))))
				if base.Ui32(int32(4)) <= base.Ui32(v105) {
					v175 = base.Simd_g_const(&F_SSIMGetClipped_C__k0)
					v176 = int32(0)
					v193 = v77
					v196 = v87 + (v75<<(uint(v88)%32) - v89)
					v200 = v107
					v221 = base.Simd_g_i32x4_replace_lane_l0(v175, v124)
					v222 = base.Simd_g_i32x4_replace_lane_l0(v175, v144)
					v223 = base.Simd_g_i32x4_replace_lane_l0(v175, v143)
					v224 = base.Simd_g_i32x4_replace_lane_l0(v175, v142)
					v225 = base.Simd_g_i32x4_replace_lane_l0(v175, v141)
					v226 = base.Simd_g_i32x4_replace_lane_l0(v175, v115)
					for {
						v237 = int32(0)
						v238 = base.Simd_g_v128_load(m, v196, v237)
						v239 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_splat(v172), v238)
						v242 = base.Simd_g_v128_load32_zero(m, v136+v193, v237)
						v244 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v242))
						v245 = base.Simd_g_i32x4_mul(v239, v244)
						v246 = base.Simd_g_i32x4_add(v245, v223)
						v249 = base.Simd_g_v128_load32_zero(m, v137+v193, v237)
						v251 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v249))
						v252 = base.Simd_g_i32x4_mul(v239, v251)
						v253 = base.Simd_g_i32x4_add(v252, v222)
						v255 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v245, v244), v226)
						v257 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v252, v244), v225)
						v259 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v252, v251), v224)
						v264 = base.Simd_g_i32x4_add(v239, v221)
						v266 = v200 + int32(-4)
						if v266 != 0 {
							v193 = v193 + int32(4)
							v196 = v196 + int32(16)
							v200 = v266
							v221 = v264
							v222 = v253
							v223 = v246
							v224 = v259
							v225 = v257
							v226 = v255
							continue
						} else {
							break
						}
						break
					}
					v267 = base.Simd_g_const(&F_SSIMGetClipped_C__k1)
					v269 = base.Simd_g_i32x4_add(v264, base.Simd_g_i8x16_shuffle2(v264, v239, base.Simd_g_const(&F_SSIMGetClipped_C__k1), base.Simd_g_const(&F_SSIMGetClipped_C__k2)))
					v270 = base.Simd_g_const(&F_SSIMGetClipped_C__k3)
					v273 = int32(0)
					v274 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v269, base.Simd_g_i8x16_shuffle2(v269, v269, base.Simd_g_const(&F_SSIMGetClipped_C__k3), base.Simd_g_const(&F_SSIMGetClipped_C__k2))))
					v277 = base.Simd_g_i32x4_add(v253, base.Simd_g_i8x16_shuffle2(v253, v269, base.Simd_g_const(&F_SSIMGetClipped_C__k1), base.Simd_g_const(&F_SSIMGetClipped_C__k2)))
					v282 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v277, base.Simd_g_i8x16_shuffle2(v277, v277, base.Simd_g_const(&F_SSIMGetClipped_C__k3), base.Simd_g_const(&F_SSIMGetClipped_C__k2))))
					v285 = base.Simd_g_i32x4_add(v246, base.Simd_g_i8x16_shuffle2(v246, v277, base.Simd_g_const(&F_SSIMGetClipped_C__k1), base.Simd_g_const(&F_SSIMGetClipped_C__k2)))
					v290 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v285, base.Simd_g_i8x16_shuffle2(v285, v285, base.Simd_g_const(&F_SSIMGetClipped_C__k3), base.Simd_g_const(&F_SSIMGetClipped_C__k2))))
					v293 = base.Simd_g_i32x4_add(v259, base.Simd_g_i8x16_shuffle2(v259, v285, base.Simd_g_const(&F_SSIMGetClipped_C__k1), base.Simd_g_const(&F_SSIMGetClipped_C__k2)))
					v298 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v293, base.Simd_g_i8x16_shuffle2(v293, v293, base.Simd_g_const(&F_SSIMGetClipped_C__k3), base.Simd_g_const(&F_SSIMGetClipped_C__k2))))
					v301 = base.Simd_g_i32x4_add(v257, base.Simd_g_i8x16_shuffle2(v257, v293, base.Simd_g_const(&F_SSIMGetClipped_C__k1), base.Simd_g_const(&F_SSIMGetClipped_C__k2)))
					v306 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v301, base.Simd_g_i8x16_shuffle2(v301, v301, base.Simd_g_const(&F_SSIMGetClipped_C__k3), base.Simd_g_const(&F_SSIMGetClipped_C__k2))))
					v309 = base.Simd_g_i32x4_add(v255, base.Simd_g_i8x16_shuffle2(v255, v301, base.Simd_g_const(&F_SSIMGetClipped_C__k1), base.Simd_g_const(&F_SSIMGetClipped_C__k2)))
					v314 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v309, base.Simd_g_i8x16_shuffle2(v309, v309, base.Simd_g_const(&F_SSIMGetClipped_C__k3), base.Simd_g_const(&F_SSIMGetClipped_C__k2))))
					if v105 == v107 {
						v437 = v314
						v446 = v274
						v463 = v306
						v464 = v298
						v465 = v290
						v466 = v282
					} else {
						v316 = v314
						v320 = v77 + v107
						v325 = v274
						v342 = v306
						v343 = v298
						v344 = v290
						v345 = v282
						v367 = v316
						v371 = v320
						v374 = v87 - v89 + int32(12) + v320<<(uint(int32(2))%32)
						v376 = v325
						v393 = v342
						v394 = v343
						v395 = v344
						v396 = v345
						for {
							v415 = *(*int32)(unsafe.Add(mBase, uint32(v374)))
							v416 = v172 * v415
							v418 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v136+v371))))
							v419 = v416 * v418
							v420 = v419 + v395
							v422 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v137+v371))))
							v423 = v416 * v422
							v424 = v423 + v396
							v426 = v419*v418 + v367
							v428 = v423*v418 + v393
							v430 = v423*v422 + v394
							v433 = v416 + v376
							v435 = v371 + int32(1)
							if v83+int32(1) != v435 {
								v367 = v426
								v371 = v435
								v374 = v374 + int32(4)
								v376 = v433
								v393 = v428
								v394 = v430
								v395 = v420
								v396 = v424
								continue
							} else {
								break
							}
							break
						}
						v437 = v426
						v446 = v433
						v463 = v428
						v464 = v430
						v465 = v420
						v466 = v424
					}
				} else {
					v316 = v115
					v320 = v77
					v325 = v124
					v342 = v141
					v343 = v142
					v344 = v143
					v345 = v144
					v367 = v316
					v371 = v320
					v374 = v87 - v89 + int32(12) + v320<<(uint(int32(2))%32)
					v376 = v325
					v393 = v342
					v394 = v343
					v395 = v344
					v396 = v345
					for {
						v415 = *(*int32)(unsafe.Add(mBase, uint32(v374)))
						v416 = v172 * v415
						v418 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v136+v371))))
						v419 = v416 * v418
						v420 = v419 + v395
						v422 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v137+v371))))
						v423 = v416 * v422
						v424 = v423 + v396
						v426 = v419*v418 + v367
						v428 = v423*v418 + v393
						v430 = v423*v422 + v394
						v433 = v416 + v376
						v435 = v371 + int32(1)
						if v83+int32(1) != v435 {
							v367 = v426
							v371 = v435
							v374 = v374 + int32(4)
							v376 = v433
							v393 = v428
							v394 = v430
							v395 = v420
							v396 = v424
							continue
						} else {
							break
						}
						break
					}
					v437 = v426
					v446 = v433
					v463 = v428
					v464 = v430
					v465 = v420
					v466 = v424
				}
				if v125 < v62 {
					v115 = v437
					v124 = v446
					v125 = v125 + int32(1)
					v136 = v136 + l3
					v137 = v137 + l1
					v141 = v463
					v142 = v464
					v143 = v465
					v144 = v466
					continue
				} else {
					break
				}
				break
			}
			v503 = base.I64_extend_i32_u(v466)
			v504 = v446
			v508 = base.I64_extend_i32_u(v465)
			v509 = base.I64_extend_i32_u(v464)
			v510 = base.I64_extend_i32_u(v463)
			v511 = base.I64_extend_i32_u(v437)
		}
	} else {
		v64 = int64(0)
		v503 = v9
		v504 = v10
		v508 = v64
		v509 = v64
		v510 = v64
		v511 = v64
	}
	v546 = v508*v508 + v503*v503
	v547 = v504 * v504
	if base.Ui64(v546) < base.Ui64(base.I64_extend_i32_u(v547<<(uint(int32(6))%32))) {
		v588 = float64(1)
	} else {
		v552 = base.I64_extend_i32_u(v504)
		v554 = v508 * v503
		v555 = v510*v552 - v554
		v556 = int64(0)
		if v556 < v555 {
			v559 = v555
		} else {
			v559 = v556
		}
		v560 = int64(1)
		v564 = base.I64_extend_i32_u(v547 * int32(60))
		v566 = int64(8)
		v572 = base.I64_extend_i32_u(v547 * int32(20))
		v588 = base.F64_div(base.F64_convert_i64_u(int64(base.Ui64(v559<<(uint(v560)%64)+v564)>>(uint(v566)%64))*(v554<<(uint(v560)%64)+v572)), base.F64_convert_i64_u(int64(base.Ui64(v564-v546+(v511+v509)*v552)>>(uint(v566)%64))*(v546+v572)))
	}
	return v588
}

var F_SSIMGetClipped_C__k0 = [2]uint64{0x0, 0x0}
var F_SSIMGetClipped_C__k1 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_SSIMGetClipped_C__k2 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SSIMGetClipped_C__k3 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_SSIMGet_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) float64 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v7 base.V128
	_ = v7
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v34 base.V128
	_ = v34
	var v36 base.V128
	_ = v36
	var v37 base.V128
	_ = v37
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v49 int32
	_ = v49
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v55 base.V128
	_ = v55
	var v58 int32
	_ = v58
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v67 int32
	_ = v67
	var v69 base.V128
	_ = v69
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v75 int32
	_ = v75
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v90 base.V128
	_ = v90
	var v91 base.V128
	_ = v91
	var v93 base.V128
	_ = v93
	var v94 base.V128
	_ = v94
	var v100 int32
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v107 base.V128
	_ = v107
	var v109 base.V128
	_ = v109
	var v112 int32
	_ = v112
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v119 int32
	_ = v119
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v126 int32
	_ = v126
	var v128 base.V128
	_ = v128
	var v130 base.V128
	_ = v130
	var v133 int32
	_ = v133
	var v135 base.V128
	_ = v135
	var v137 base.V128
	_ = v137
	var v142 base.V128
	_ = v142
	var v144 base.V128
	_ = v144
	var v146 base.V128
	_ = v146
	var v149 base.V128
	_ = v149
	var v156 base.V128
	_ = v156
	var v158 base.V128
	_ = v158
	var v161 base.V128
	_ = v161
	var v164 base.V128
	_ = v164
	var v167 base.V128
	_ = v167
	var v170 base.V128
	_ = v170
	var v173 base.V128
	_ = v173
	var v175 base.V128
	_ = v175
	var v178 base.V128
	_ = v178
	var v190 base.V128
	_ = v190
	var v193 base.V128
	_ = v193
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v199 int32
	_ = v199
	var v211 base.V128
	_ = v211
	var v214 base.V128
	_ = v214
	var v228 int32
	_ = v228
	var v234 int64
	_ = v234
	var v236 int64
	_ = v236
	var v238 int64
	_ = v238
	var v241 int64
	_ = v241
	var v244 int64
	_ = v244
	var v245 int64
	_ = v245
	var v246 int64
	_ = v246
	var v249 int64
	_ = v249
	var v250 int64
	_ = v250
	var v252 int64
	_ = v252
	var v254 int64
	_ = v254
	var v258 int64
	_ = v258
	var v262 int64
	_ = v262
	var v263 int64
	_ = v263
	var v277 float64
	_ = v277
	v5 = int32(0)
	v7 = base.Simd_g_const(&F_SSIMGet_SSE2__k0)
	v28 = m.G0
	v30 = v28 - int32(32)
	m.G0 = v30
	v32 = l2 + l3
	v34 = base.Simd_g_v128_load64_zero(m, v32, v5)
	v36 = base.Simd_g_const(&F_SSIMGet_SSE2__k1)
	v37 = base.Simd_g_i8x16_shuffle2(v34, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v38 = base.Simd_g_const(&F_SSIMGet_SSE2__k4)
	v39 = base.Simd_g_i16x8_mul(v37, v38)
	v42 = base.Simd_g_v128_load64_zero(m, l2, v5)
	v44 = base.Simd_g_i8x16_shuffle2(v42, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v45 = base.Simd_g_const(&F_SSIMGet_SSE2__k5)
	v46 = base.Simd_g_i16x8_mul(v44, v45)
	v49 = v32 + l3
	v51 = base.Simd_g_v128_load64_zero(m, v49, v5)
	v53 = base.Simd_g_i8x16_shuffle2(v51, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v54 = base.Simd_g_const(&F_SSIMGet_SSE2__k6)
	v55 = base.Simd_g_i16x8_mul(v53, v54)
	v58 = v49 + l3
	v60 = base.Simd_g_v128_load64_zero(m, v58, v5)
	v62 = base.Simd_g_i8x16_shuffle2(v60, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v63 = base.Simd_g_const(&F_SSIMGet_SSE2__k7)
	v64 = base.Simd_g_i16x8_mul(v62, v63)
	v67 = v58 + l3
	v69 = base.Simd_g_v128_load64_zero(m, v67, v5)
	v71 = base.Simd_g_i8x16_shuffle2(v69, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v72 = base.Simd_g_i16x8_mul(v71, v54)
	v75 = v67 + l3
	v77 = base.Simd_g_v128_load64_zero(m, v75, v5)
	v79 = base.Simd_g_i8x16_shuffle2(v77, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v80 = base.Simd_g_i16x8_mul(v79, v38)
	v85 = base.Simd_g_v128_load64_zero(m, v75+l3, v5)
	v87 = base.Simd_g_i8x16_shuffle2(v85, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v88 = base.Simd_g_i16x8_mul(v87, v45)
	v90 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v37, v39), base.Simd_g_i32x4_dot_i16x8_s(v44, v46)), base.Simd_g_i32x4_dot_i16x8_s(v53, v55)), base.Simd_g_i32x4_dot_i16x8_s(v62, v64)), base.Simd_g_i32x4_dot_i16x8_s(v71, v72)), base.Simd_g_i32x4_dot_i16x8_s(v79, v80)), base.Simd_g_i32x4_dot_i16x8_s(v87, v88))
	v91 = base.Simd_g_const(&F_SSIMGet_SSE2__k8)
	v93 = base.Simd_g_i32x4_add(v90, base.Simd_g_i8x16_shuffle2(v90, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k9), base.Simd_g_const(&F_SSIMGet_SSE2__k10)))
	v94 = base.Simd_g_const(&F_SSIMGet_SSE2__k11)
	base.Simd_g_v128_store32_lane_l0(m, v30, int32(28), base.Simd_g_i32x4_add(v93, base.Simd_g_i8x16_shuffle2(v93, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k12), base.Simd_g_const(&F_SSIMGet_SSE2__k13))))
	v100 = l0 + l1
	v102 = base.Simd_g_v128_load64_zero(m, v100, v5)
	v104 = base.Simd_g_i8x16_shuffle2(v102, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v107 = base.Simd_g_v128_load64_zero(m, l0, v5)
	v109 = base.Simd_g_i8x16_shuffle2(v107, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v112 = v100 + l1
	v114 = base.Simd_g_v128_load64_zero(m, v112, v5)
	v116 = base.Simd_g_i8x16_shuffle2(v114, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v119 = v112 + l1
	v121 = base.Simd_g_v128_load64_zero(m, v119, v5)
	v123 = base.Simd_g_i8x16_shuffle2(v121, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v126 = v119 + l1
	v128 = base.Simd_g_v128_load64_zero(m, v126, v5)
	v130 = base.Simd_g_i8x16_shuffle2(v128, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v133 = v126 + l1
	v135 = base.Simd_g_v128_load64_zero(m, v133, v5)
	v137 = base.Simd_g_i8x16_shuffle2(v135, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v142 = base.Simd_g_v128_load64_zero(m, v133+l1, v5)
	v144 = base.Simd_g_i8x16_shuffle2(v142, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k2), base.Simd_g_const(&F_SSIMGet_SSE2__k3))
	v146 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v104, v39), base.Simd_g_i32x4_dot_i16x8_s(v109, v46)), base.Simd_g_i32x4_dot_i16x8_s(v116, v55)), base.Simd_g_i32x4_dot_i16x8_s(v123, v64)), base.Simd_g_i32x4_dot_i16x8_s(v130, v72)), base.Simd_g_i32x4_dot_i16x8_s(v137, v80)), base.Simd_g_i32x4_dot_i16x8_s(v144, v88))
	v149 = base.Simd_g_i32x4_add(v146, base.Simd_g_i8x16_shuffle2(v146, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k9), base.Simd_g_const(&F_SSIMGet_SSE2__k10)))
	base.Simd_g_v128_store32_lane_l0(m, v30, int32(24), base.Simd_g_i32x4_add(v149, base.Simd_g_i8x16_shuffle2(v149, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k12), base.Simd_g_const(&F_SSIMGet_SSE2__k13))))
	v156 = base.Simd_g_i16x8_mul(v104, v38)
	v158 = base.Simd_g_i16x8_mul(v109, v45)
	v161 = base.Simd_g_i16x8_mul(v116, v54)
	v164 = base.Simd_g_i16x8_mul(v123, v63)
	v167 = base.Simd_g_i16x8_mul(v130, v54)
	v170 = base.Simd_g_i16x8_mul(v137, v38)
	v173 = base.Simd_g_i16x8_mul(v144, v45)
	v175 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v104, v156), base.Simd_g_i32x4_dot_i16x8_s(v109, v158)), base.Simd_g_i32x4_dot_i16x8_s(v116, v161)), base.Simd_g_i32x4_dot_i16x8_s(v123, v164)), base.Simd_g_i32x4_dot_i16x8_s(v130, v167)), base.Simd_g_i32x4_dot_i16x8_s(v137, v170)), base.Simd_g_i32x4_dot_i16x8_s(v144, v173))
	v178 = base.Simd_g_i32x4_add(v175, base.Simd_g_i8x16_shuffle2(v175, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k9), base.Simd_g_const(&F_SSIMGet_SSE2__k10)))
	base.Simd_g_v128_store32_lane_l0(m, v30, int32(20), base.Simd_g_i32x4_add(v178, base.Simd_g_i8x16_shuffle2(v178, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k12), base.Simd_g_const(&F_SSIMGet_SSE2__k13))))
	v190 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v39, v46), v55), v64), v72), v80), v88)
	v193 = base.Simd_g_i16x8_add(v190, base.Simd_g_i8x16_shuffle2(v190, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k9), base.Simd_g_const(&F_SSIMGet_SSE2__k10)))
	v194 = int32(3)
	v196 = int32(2)
	v199 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v30)+16)) = base.Simd_g_i16x8_extract_lane_u_l3(v193) + base.Simd_g_i16x8_extract_lane_u_l2(v193) + base.Simd_g_i16x8_extract_lane_u_l1(v193) + base.Simd_g_i16x8_extract_lane_u_l0(v193)
	v211 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v156, v158), v161), v164), v167), v170), v173)
	v214 = base.Simd_g_i16x8_add(v211, base.Simd_g_i8x16_shuffle2(v211, v7, base.Simd_g_const(&F_SSIMGet_SSE2__k9), base.Simd_g_const(&F_SSIMGet_SSE2__k10)))
	*(*int32)(unsafe.Add(mBase, uint32(v30)+12)) = base.Simd_g_i16x8_extract_lane_u_l3(v214) + base.Simd_g_i16x8_extract_lane_u_l2(v214) + base.Simd_g_i16x8_extract_lane_u_l1(v214) + base.Simd_g_i16x8_extract_lane_u_l0(v214)
	v228 = v30 + int32(8)
	v234 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v228)+8)))
	v236 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v228)+4)))
	v238 = v234*v234 + v236*v236
	if base.Ui64(v238) < base.Ui64(int64(4194304)) {
		v277 = float64(1)
	} else {
		v241 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v228)+16)))
		v244 = v234 * v236
		v245 = v241<<(uint(int64(8))%64) - v244
		v246 = int64(0)
		if v246 < v245 {
			v249 = v245
		} else {
			v249 = v246
		}
		v250 = int64(1)
		v252 = int64(3932160)
		v254 = int64(8)
		v258 = int64(1310720)
		v262 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v228)+20)))
		v263 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v228)+12)))
		v277 = base.F64_div(base.F64_convert_i64_u(int64(base.Ui64(v249<<(uint(v250)%64)+v252)>>(uint(v254)%64))*(v244<<(uint(v250)%64)+v258)), base.F64_convert_i64_u(int64(base.Ui64((v262+v263)<<(uint(v254)%64)-v238+v252)>>(uint(v254)%64))*(v238+v258)))
	}
	m.G0 = v30 + int32(32)
	return v277
}

var F_SSIMGet_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SSIMGet_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SSIMGet_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SSIMGet_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_SSIMGet_SSE2__k4 = [2]uint64{0x8000600040002, 0x200040006}
var F_SSIMGet_SSE2__k5 = [2]uint64{0x4000300020001, 0x100020003}
var F_SSIMGet_SSE2__k6 = [2]uint64{0xc000900060003, 0x300060009}
var F_SSIMGet_SSE2__k7 = [2]uint64{0x10000c00080004, 0x40008000c}
var F_SSIMGet_SSE2__k8 = [2]uint64{0xf0e0d0c0b0a0908, 0x1716151413121110}
var F_SSIMGet_SSE2__k9 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_SSIMGet_SSE2__k10 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_SSIMGet_SSE2__k11 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_SSIMGet_SSE2__k12 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_SSIMGet_SSE2__k13 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_SetLoopParams(m *base.Module, l0 int32, l1 float32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v19 float32
	_ = v19
	var v22 float32
	_ = v22
	var v25 float32
	_ = v25
	var v27 base.V128
	_ = v27
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v69 int32
	_ = v69
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v91 int32
	_ = v91
	var v102 int32
	_ = v102
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v130 base.V128
	_ = v130
	var v133 int32
	_ = v133
	var v136 int32
	_ = v136
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v148 int32
	_ = v148
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v158 int32
	_ = v158
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v169 int32
	_ = v169
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v180 int32
	_ = v180
	var v184 int32
	_ = v184
	var v187 int32
	_ = v187
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v194 int32
	_ = v194
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v213 int32
	_ = v213
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v231 int32
	_ = v231
	var v237 int32
	_ = v237
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v249 int32
	_ = v249
	var v251 int32
	_ = v251
	var v254 int32
	_ = v254
	var v260 int32
	_ = v260
	var v268 int32
	_ = v268
	var v277 int32
	_ = v277
	var v285 int32
	_ = v285
	var v296 int32
	_ = v296
	var v304 int32
	_ = v304
	var v356 int32
	_ = v356
	var v371 int32
	_ = v371
	var v372 int32
	_ = v372
	var v373 int32
	_ = v373
	var v375 int32
	_ = v375
	var v378 int32
	_ = v378
	var v381 int32
	_ = v381
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v390 int32
	_ = v390
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v399 int32
	_ = v399
	var v402 int32
	_ = v402
	var v405 int32
	_ = v405
	var v408 int32
	_ = v408
	var v411 int32
	_ = v411
	var v414 int32
	_ = v414
	var v417 int32
	_ = v417
	var v422 int32
	_ = v422
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v445 int32
	_ = v445
	var v464 int32
	_ = v464
	var v496 int32
	_ = v496
	var v497 int32
	_ = v497
	var v498 int32
	_ = v498
	var v499 int32
	_ = v499
	var v512 int32
	_ = v512
	var v513 int32
	_ = v513
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v518 int32
	_ = v518
	var v522 int32
	_ = v522
	var v529 int32
	_ = v529
	var v562 int32
	_ = v562
	var v573 int32
	_ = v573
	var v578 int32
	_ = v578
	var v579 int32
	_ = v579
	var v581 int32
	_ = v581
	var v590 int32
	_ = v590
	var v617 int32
	_ = v617
	var v619 int32
	_ = v619
	var v620 int32
	_ = v620
	var v628 int32
	_ = v628
	var v629 int32
	_ = v629
	var v631 int32
	_ = v631
	var v640 int32
	_ = v640
	var v642 int32
	_ = v642
	var v644 int32
	_ = v644
	var v688 int32
	_ = v688
	var v693 int32
	_ = v693
	var v696 int32
	_ = v696
	var v699 int32
	_ = v699
	var v702 int32
	_ = v702
	var v703 int32
	_ = v703
	var v704 int32
	_ = v704
	var v708 int32
	_ = v708
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v715 int32
	_ = v715
	var v716 int32
	_ = v716
	var v723 int32
	_ = v723
	var v759 int32
	_ = v759
	var v770 int32
	_ = v770
	var v775 int32
	_ = v775
	var v776 int32
	_ = v776
	var v778 int32
	_ = v778
	var v787 int32
	_ = v787
	var v814 int32
	_ = v814
	var v816 int32
	_ = v816
	var v817 int32
	_ = v817
	var v825 int32
	_ = v825
	var v826 int32
	_ = v826
	var v828 int32
	_ = v828
	var v837 int32
	_ = v837
	var v839 int32
	_ = v839
	var v841 int32
	_ = v841
	var v885 int32
	_ = v885
	var v890 int32
	_ = v890
	var v893 int32
	_ = v893
	var v896 int32
	_ = v896
	var v899 int32
	_ = v899
	var v900 int32
	_ = v900
	var v901 int32
	_ = v901
	var v905 int32
	_ = v905
	var v906 int32
	_ = v906
	var v907 int32
	_ = v907
	var v912 int32
	_ = v912
	var v913 int32
	_ = v913
	var v920 int32
	_ = v920
	var v956 int32
	_ = v956
	var v967 int32
	_ = v967
	var v972 int32
	_ = v972
	var v973 int32
	_ = v973
	var v975 int32
	_ = v975
	var v984 int32
	_ = v984
	var v1011 int32
	_ = v1011
	var v1013 int32
	_ = v1013
	var v1014 int32
	_ = v1014
	var v1022 int32
	_ = v1022
	var v1023 int32
	_ = v1023
	var v1025 int32
	_ = v1025
	var v1034 int32
	_ = v1034
	var v1036 int32
	_ = v1036
	var v1038 int32
	_ = v1038
	var v1082 int32
	_ = v1082
	var v1087 int32
	_ = v1087
	var v1090 int32
	_ = v1090
	var v1093 int32
	_ = v1093
	var v1096 int32
	_ = v1096
	var v1103 int32
	_ = v1103
	var v1108 int32
	_ = v1108
	var v1109 int32
	_ = v1109
	var v1111 int32
	_ = v1111
	var v1113 int32
	_ = v1113
	var v1115 int32
	_ = v1115
	var v1117 int32
	_ = v1117
	var v1119 int32
	_ = v1119
	var v1121 int32
	_ = v1121
	var v1123 int32
	_ = v1123
	var v1125 int32
	_ = v1125
	var v1127 int32
	_ = v1127
	var v1129 int32
	_ = v1129
	var v1131 int32
	_ = v1131
	var v1133 int32
	_ = v1133
	var v1135 int32
	_ = v1135
	var v1137 int32
	_ = v1137
	var v1139 int32
	_ = v1139
	var v1141 int32
	_ = v1141
	var v1144 int32
	_ = v1144
	var v1237 int32
	_ = v1237
	var v1244 int32
	_ = v1244
	var v1290 int64
	_ = v1290
	var v1296 int32
	_ = v1296
	v14 = m.G0
	v16 = v14 - int32(16)
	m.G0 = v16
	v19 = float32(100)
	if base.F32_gt(l1, v19) != 0 {
		v22 = v19
	} else {
		v22 = l1
	}
	if base.F32_lt(l1, float32(0)) != 0 {
		v25 = float32(0)
	} else {
		v25 = v22
	}
	F_VP8SetSegmentParams(m, l0, v25)
	mBase = m.M
	v27 = base.Simd_g_const(&F_SetLoopParams__k0)
	base.Simd_g_v128_store(m, v16, int32(0), v27)
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v32 = v30 * v31
	if v32 < int32(1) {
	} else {
		v35 = int32(1)
		v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[0])))
		if v32 == v35 {
			v91 = int32(0)
		} else {
			v48 = v37
			v51 = int32(0)
			for {
				v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48))))
				v58 = int32(3)
				v60 = int32(12)
				v62 = v16 | int32(base.Ui32(v57)>>(uint(v58)%32))&v60
				v63 = *(*int32)(unsafe.Add(mBase, uint32(v62)))
				v64 = int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v62))) = v63 + v64
				v69 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v48+int32(4)))))
				v74 = v16 | int32(base.Ui32(v69)>>(uint(v58)%32))&v60
				v75 = *(*int32)(unsafe.Add(mBase, uint32(v74)))
				*(*int32)(unsafe.Add(mBase, uint32(v74))) = v75 + v64
				v82 = v51 + int32(2)
				if v32&int32(2147483646) != v82 {
					v48 = v48 + int32(8)
					v51 = v82
					continue
				} else {
					break
				}
				break
			}
			v91 = v82
		}
		if v32&v35 == int32(0) {
		} else {
			v102 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37+v91<<(uint(int32(2))%32)))))
			v107 = v16 | int32(base.Ui32(v102)>>(uint(int32(3))%32))&int32(12)
			v108 = *(*int32)(unsafe.Add(mBase, uint32(v107)))
			*(*int32)(unsafe.Add(mBase, uint32(v107))) = v108 + int32(1)
		}
	}
	v125 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v126 = *(*int32)(unsafe.Add(mBase, uint32(v125)+88))
	if v126 == int32(0) {
	} else {
		v130 = base.Simd_g_v128_load(m, v16, int32(0))
		base.Simd_g_v128_store(m, v126, int32(92), v130)
	}
	v133 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if v133 < int32(2) {
		v226 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v226
		v304 = v226
	} else {
		v136 = int32(255)
		v138 = *(*int32)(unsafe.Add(mBase, uint32(v16)+12))
		v139 = *(*int32)(unsafe.Add(mBase, uint32(v16)+8))
		v140 = v138 + v139
		v141 = *(*int32)(unsafe.Add(mBase, uint32(v16)+4))
		v142 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
		v143 = v141 + v142
		v144 = v140 + v143
		if v144 == int32(0) {
			v153 = v136
		} else {
			v148 = base.I32_div_s(v144, int32(2))
			v152 = base.I32_div_s(v148+v143*int32(255), v144)
			v153 = v152
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+3416)) = uint8(v153)
		if v143 == int32(0) {
			v163 = v136
		} else {
			v158 = base.I32_div_s(v143, int32(2))
			v162 = base.I32_div_s(v158+v142*int32(255), v143)
			v163 = v162
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)) = uint8(v163)
		if v140 == int32(0) {
			v174 = int32(255)
		} else {
			v169 = base.I32_div_s(v140, int32(2))
			v173 = base.I32_div_s(v169+v139*int32(255), v140)
			v174 = v173
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)) = uint8(v174)
		v176 = int32(255)
		if v153&v176 != v176 {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
			v237 = v153
			v241 = v174
			v242 = v163
		} else {
			v180 = int32(255)
			if v163&v180 != v180 {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
				v237 = v153
				v241 = v174
				v242 = v163
			} else {
				v184 = int32(255)
				v187 = base.B2i32(v174&v184 != v184)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v187
				if v174&v184 != v184 {
					v231 = int32(255)
					v237 = v231
					v241 = v174
					v242 = v231
				} else {
					v189 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
					v190 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
					if v189*v190 < int32(1) {
						v231 = int32(255)
						v237 = v231
						v241 = v174
						v242 = v231
					} else {
						v194 = int32(0)
						v200 = v194
						v203 = v194
						for {
							v209 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[0])))
							v210 = v209 + v200
							v211 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v210))))
							v213 = v211 & int32(159)
							*(*uint8)(unsafe.Add(mBase, uint32(v210))) = uint8(v213)
							v218 = v203 + int32(1)
							v219 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
							v220 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
							if v218 < v219*v220 {
								v200 = v200 + int32(4)
								v203 = v218
								continue
							} else {
								break
							}
							break
						}
						v223 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)))
						v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)))
						v225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3416)))
						v237 = v225
						v241 = v223
						v242 = v224
					}
				}
			}
		}
		v246 = m.G79
		v247 = int32(-1)
		v249 = int32(255)
		v251 = int32(1)
		v254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v246+(v242^v247)&v249<<(uint(v251)%32)))))
		v260 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v246+v237&v249<<(uint(v251)%32)))))
		v268 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v246+v242&v249<<(uint(v251)%32)))))
		v277 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v246+v241&v249<<(uint(v251)%32)))))
		v285 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v246+(v237^v247)&v249<<(uint(v251)%32)))))
		v296 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v246+(v241^v247)&v249<<(uint(v251)%32)))))
		v304 = (v254+v260)*v141 + (v268+v260)*v142 + (v277+v285)*v139 + (v296+v285)*v138
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v304
	v356 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[1])))
	if v356 == int32(0) {
	} else {
		v371 = m.G81
		v372 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+15)))
		v373 = int32(408)
		v375 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+14)))
		v378 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+13)))
		v381 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+12)))
		v384 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+11)))
		v387 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+10)))
		v390 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+9)))
		v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+8)))
		v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+7)))
		v399 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+6)))
		v402 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+5)))
		v405 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+4)))
		v408 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+3)))
		v411 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+2)))
		v414 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371)+1)))
		v417 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v371))))
		v422 = l0 + int32(3444)
		v423 = l0 + int32(3433)
		v424 = l0 + int32(3422)
		v445 = int32(0)
		for {
			v464 = l0 + int32(_a_F_SetLoopParams_0) + v445*int32(3264)
			v496 = v422
			v497 = v423
			v498 = v424
			v499 = int32(0)
			for {
				v512 = l0 + int32(3420) + v445*int32(264) + v499*int32(33)
				v513 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v512)+1)))
				v514 = int32(1)
				v517 = v464 + v499*int32(408)
				v518 = m.G79
				v522 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v518+v513<<(uint(v514)%32)))))
				*(*uint16)(unsafe.Add(mBase, uint32(v517))) = uint16(v522)
				v529 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v518+(v513^int32(255))<<(uint(v514)%32)))))
				v562 = v514
				for {
					v573 = m.G1
					v578 = v562<<(uint(int32(2))%32) + (v573 + int32(_a_F_SetLoopParams_1)) + int32(-4)
					v579 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v578))))
					if v579 != 0 {
						v581 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v578)+2)))
						v590 = v579
						v617 = v581
						v619 = int32(0)
						v620 = v498
						for {
							if v590&int32(1) == int32(0) {
								v642 = v619
							} else {
								v628 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v620))))
								v629 = m.G79
								v631 = int32(1)
								v640 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v629+(v628^(int32(0)-v617&v631)&int32(255))<<(uint(v631)%32)))))
								v642 = v619 + v640
							}
							v644 = int32(1)
							if base.Ui32(v644) < base.Ui32(v590) {
								v590 = int32(base.Ui32(v590) >> (uint(v644) % 32))
								v617 = int32(base.Ui32(v617) >> (uint(v644) % 32))
								v619 = v642
								v620 = v620 + v644
								continue
							} else {
								break
							}
							break
						}
						v688 = v642
					} else {
						v688 = int32(0)
					}
					v693 = int32(1)
					v696 = v529 + v688
					*(*uint16)(unsafe.Add(mBase, uint32(v517+v562<<(uint(v693)%32)))) = uint16(v696)
					v699 = v562 + v693
					if v699 != int32(68) {
						v562 = v699
						continue
					} else {
						break
					}
					break
				}
				v702 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v512)+12)))
				v703 = int32(1)
				v704 = m.G79
				v708 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v704+v702<<(uint(v703)%32)))))
				v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v512)+11)))
				v710 = int32(255)
				v715 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v704+(v709^v710)<<(uint(v703)%32)))))
				v716 = v708 + v715
				*(*uint16)(unsafe.Add(mBase, uint32(v517)+136)) = uint16(v716)
				v723 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v704+(v702^v710)<<(uint(v703)%32)))))
				v759 = v703
				for {
					v770 = m.G1
					v775 = v759<<(uint(int32(2))%32) + (v770 + int32(_a_F_SetLoopParams_1)) + int32(-4)
					v776 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v775))))
					if v776 != 0 {
						v778 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v775)+2)))
						v787 = v776
						v814 = v778
						v816 = int32(0)
						v817 = v497
						for {
							if v787&int32(1) == int32(0) {
								v839 = v816
							} else {
								v825 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v817))))
								v826 = m.G79
								v828 = int32(1)
								v837 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v826+(v825^(int32(0)-v814&v828)&int32(255))<<(uint(v828)%32)))))
								v839 = v816 + v837
							}
							v841 = int32(1)
							if base.Ui32(v841) < base.Ui32(v787) {
								v787 = int32(base.Ui32(v787) >> (uint(v841) % 32))
								v814 = int32(base.Ui32(v814) >> (uint(v841) % 32))
								v816 = v839
								v817 = v817 + v841
								continue
							} else {
								break
							}
							break
						}
						v885 = v839
					} else {
						v885 = int32(0)
					}
					v890 = int32(1)
					v893 = v715 + v723 + v885
					*(*uint16)(unsafe.Add(mBase, uint32(v517+int32(136)+v759<<(uint(v890)%32)))) = uint16(v893)
					v896 = v759 + v890
					if v896 != int32(68) {
						v759 = v896
						continue
					} else {
						break
					}
					break
				}
				v899 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v512)+23)))
				v900 = int32(1)
				v901 = m.G79
				v905 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v901+v899<<(uint(v900)%32)))))
				v906 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v512)+22)))
				v907 = int32(255)
				v912 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v901+(v906^v907)<<(uint(v900)%32)))))
				v913 = v905 + v912
				*(*uint16)(unsafe.Add(mBase, uint32(v517)+272)) = uint16(v913)
				v920 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v901+(v899^v907)<<(uint(v900)%32)))))
				v956 = v900
				for {
					v967 = m.G1
					v972 = v956<<(uint(int32(2))%32) + (v967 + int32(_a_F_SetLoopParams_1)) + int32(-4)
					v973 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v972))))
					if v973 != 0 {
						v975 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v972)+2)))
						v984 = v973
						v1011 = v975
						v1013 = int32(0)
						v1014 = v496
						for {
							if v984&int32(1) == int32(0) {
								v1036 = v1013
							} else {
								v1022 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1014))))
								v1023 = m.G79
								v1025 = int32(1)
								v1034 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1023+(v1022^(int32(0)-v1011&v1025)&int32(255))<<(uint(v1025)%32)))))
								v1036 = v1013 + v1034
							}
							v1038 = int32(1)
							if base.Ui32(v1038) < base.Ui32(v984) {
								v984 = int32(base.Ui32(v984) >> (uint(v1038) % 32))
								v1011 = int32(base.Ui32(v1011) >> (uint(v1038) % 32))
								v1013 = v1036
								v1014 = v1014 + v1038
								continue
							} else {
								break
							}
							break
						}
						v1082 = v1036
					} else {
						v1082 = int32(0)
					}
					v1087 = int32(1)
					v1090 = v912 + v920 + v1082
					*(*uint16)(unsafe.Add(mBase, uint32(v517+int32(272)+v956<<(uint(v1087)%32)))) = uint16(v1090)
					v1093 = v956 + v1087
					if v1093 != int32(68) {
						v956 = v1093
						continue
					} else {
						break
					}
					break
				}
				v1096 = int32(33)
				v1103 = v499 + int32(1)
				if v1103 != int32(8) {
					v496 = v496 + v1096
					v497 = v497 + v1096
					v498 = v498 + v1096
					v499 = v1103
					continue
				} else {
					break
				}
				break
			}
			v1108 = l0 + int32(_a_F_SetLoopParams_2) + v445*int32(192)
			v1109 = v464 + v372*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+180)) = v1109
			v1111 = v464 + v375*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+168)) = v1111
			v1113 = v464 + v378*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+156)) = v1113
			v1115 = v464 + v381*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+144)) = v1115
			v1117 = v464 + v384*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+132)) = v1117
			v1119 = v464 + v387*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+120)) = v1119
			v1121 = v464 + v390*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+108)) = v1121
			v1123 = v464 + v393*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+96)) = v1123
			v1125 = v464 + v396*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+84)) = v1125
			v1127 = v464 + v399*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+72)) = v1127
			v1129 = v464 + v402*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+60)) = v1129
			v1131 = v464 + v405*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+48)) = v1131
			v1133 = v464 + v408*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+36)) = v1133
			v1135 = v464 + v411*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+24)) = v1135
			v1137 = v464 + v414*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+12)) = v1137
			v1139 = v464 + v417*v373
			*(*int32)(unsafe.Add(mBase, uint32(v1108))) = v1139
			v1141 = int32(272)
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+188)) = v1109 + v1141
			v1144 = int32(136)
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+184)) = v1109 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+176)) = v1111 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+172)) = v1111 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+164)) = v1113 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+160)) = v1113 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+152)) = v1115 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+148)) = v1115 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+140)) = v1117 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+136)) = v1117 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+128)) = v1119 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+124)) = v1119 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+116)) = v1121 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+112)) = v1121 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+104)) = v1123 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+100)) = v1123 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+92)) = v1125 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+88)) = v1125 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+80)) = v1127 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+76)) = v1127 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+68)) = v1129 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+64)) = v1129 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+56)) = v1131 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+52)) = v1131 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+44)) = v1133 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+40)) = v1133 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+32)) = v1135 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+28)) = v1135 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+20)) = v1137 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+16)) = v1137 + v1144
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+8)) = v1139 + v1141
			*(*int32)(unsafe.Add(mBase, uint32(v1108)+4)) = v1139 + v1144
			v1237 = int32(264)
			v1244 = v445 + int32(1)
			if v1244 != int32(4) {
				v422 = v422 + v1237
				v423 = v423 + v1237
				v424 = v424 + v1237
				v445 = v1244
				continue
			} else {
				break
			}
			break
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[1]))) = int32(0)
	}
	v1290 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[2]))) = v1290
	*(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[3]))) = v1290
	v1296 = int32(0)
	base.Simd_g_v128_store(m, l0+int32(_a_F_SetLoopParams_3), v1296, v27)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_SetLoopParams[4]))) = v1296
	m.G0 = v16 + int32(16)
	return
}

var F_SetLoopParams__k0 = [2]uint64{0x0, 0x0}

func F_SetResidualCoeffs_SSE2(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v6 base.V128
	_ = v6
	var v8 base.V128
	_ = v8
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v20 int32
	_ = v20
	*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = l0
	v6 = base.Simd_g_v128_load_rng(m, l0, int32(0), int32(0), int32(32))
	v8 = base.Simd_g_v128_load_nc(m, l0, int32(16))
	v12 = base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_s(v6, v8), base.Simd_g_const(&F_SetResidualCoeffs_SSE2__k0)))
	v13 = int32(_a_F_SetResidualCoeffs_SSE2_0)
	if v12 == v13 {
		v20 = int32(-1)
	} else {
		v20 = base.I32_clz(v12^v13) ^ int32(31)
	}
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v20
	return
}

var F_SetResidualCoeffs_SSE2__k0 = [2]uint64{0x0, 0x0}

func F_SharpYuvFilterRow_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v72 base.V128
	_ = v72
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v83 int32
	_ = v83
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v113 int32
	_ = v113
	var v114 base.V128
	_ = v114
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v120 base.V128
	_ = v120
	var v122 base.V128
	_ = v122
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v126 int32
	_ = v126
	var v129 base.V128
	_ = v129
	var v130 base.V128
	_ = v130
	var v133 base.V128
	_ = v133
	var v134 base.V128
	_ = v134
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v141 base.V128
	_ = v141
	var v144 base.V128
	_ = v144
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v156 int32
	_ = v156
	var v161 int32
	_ = v161
	var v165 base.V128
	_ = v165
	var v168 base.V128
	_ = v168
	var v171 base.V128
	_ = v171
	var v174 base.V128
	_ = v174
	var v176 base.V128
	_ = v176
	var v178 base.V128
	_ = v178
	var v180 base.V128
	_ = v180
	var v182 base.V128
	_ = v182
	var v186 base.V128
	_ = v186
	var v190 base.V128
	_ = v190
	var v193 base.V128
	_ = v193
	var v201 base.V128
	_ = v201
	var v205 int32
	_ = v205
	var v210 int32
	_ = v210
	var v215 int32
	_ = v215
	var v220 int32
	_ = v220
	var v223 base.V128
	_ = v223
	var v226 base.V128
	_ = v226
	var v229 base.V128
	_ = v229
	var v232 base.V128
	_ = v232
	var v234 base.V128
	_ = v234
	var v241 base.V128
	_ = v241
	var v243 base.V128
	_ = v243
	var v265 int32
	_ = v265
	var v269 int32
	_ = v269
	var v273 int32
	_ = v273
	var v275 base.V128
	_ = v275
	var v276 base.V128
	_ = v276
	var v280 int32
	_ = v280
	var v293 int32
	_ = v293
	var v298 int32
	_ = v298
	var v303 int32
	_ = v303
	var v307 base.V128
	_ = v307
	var v310 base.V128
	_ = v310
	var v313 base.V128
	_ = v313
	var v316 base.V128
	_ = v316
	var v318 base.V128
	_ = v318
	var v329 base.V128
	_ = v329
	var v333 int32
	_ = v333
	var v338 int32
	_ = v338
	var v343 int32
	_ = v343
	var v348 int32
	_ = v348
	var v351 base.V128
	_ = v351
	var v354 base.V128
	_ = v354
	var v357 base.V128
	_ = v357
	var v360 base.V128
	_ = v360
	var v362 base.V128
	_ = v362
	var v370 base.V128
	_ = v370
	var v403 int32
	_ = v403
	var v409 int32
	_ = v409
	var v424 int32
	_ = v424
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v455 int32
	_ = v455
	var v459 int32
	_ = v459
	var v462 int32
	_ = v462
	var v463 int32
	_ = v463
	var v467 int32
	_ = v467
	var v498 int32
	_ = v498
	var v499 int32
	_ = v499
	var v502 int32
	_ = v502
	var v503 int32
	_ = v503
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v510 int32
	_ = v510
	var v519 int32
	_ = v519
	var v520 int32
	_ = v520
	var v522 int32
	_ = v522
	var v525 int32
	_ = v525
	var v527 int32
	_ = v527
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v546 int32
	_ = v546
	var v549 int32
	_ = v549
	var v551 int32
	_ = v551
	var v558 int32
	_ = v558
	if l2 < int32(1) {
	} else {
		v41 = int32(-1)
		v44 = v41<<(uint(l5)%32) ^ v41
		if base.Ui32(l2) < base.Ui32(int32(8)) {
			v424 = int32(0)
			v425 = l1
			v426 = l0
			v455 = v424 << (uint(int32(2)) % 32)
			v459 = int32(0)
			v462 = l3 + v455
			v463 = l4 + v455
			v467 = l2 - v424
			for {
				v498 = v426 + v459
				v499 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498))))
				v502 = v425 + v459
				v503 = int32(2)
				v505 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502+v503))))
				v507 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502))))
				v510 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498+v503))))
				v519 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462))))
				v520 = (v499*int32(9)+v505+(v507+v510)*int32(3)+int32(8))>>(uint(int32(4))%32) + v519
				if v520 < v44 {
					v522 = v520
				} else {
					v522 = v44
				}
				if v520 < int32(0) {
					v525 = int32(0)
				} else {
					v525 = v522
				}
				*(*uint16)(unsafe.Add(mBase, uint32(v463))) = uint16(v525)
				v527 = int32(2)
				v543 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462+v527))))
				v544 = (v507+v510*int32(9)+(v505+v499)*int32(3)+int32(8))>>(uint(int32(4))%32) + v543
				if v544 < v44 {
					v546 = v544
				} else {
					v546 = v44
				}
				if v544 < int32(0) {
					v549 = int32(0)
				} else {
					v549 = v546
				}
				*(*uint16)(unsafe.Add(mBase, uint32(v463+v527))) = uint16(v549)
				v551 = int32(4)
				v558 = v467 + int32(-1)
				if v558 != 0 {
					v459 = v459 + int32(2)
					v462 = v462 + v551
					v463 = v463 + v551
					v467 = v558
					continue
				} else {
					break
				}
				break
			}
		} else {
			v49 = int32(2)
			v50 = l2<<(uint(int32(1))%32) + v49
			v54 = l2 << (uint(v49) % 32)
			v55 = l4 + v54
			if base.B2i32(base.Ui32(l4) < base.Ui32(l0+v50))&base.B2i32(base.Ui32(l0) < base.Ui32(v55)) != 0 {
				v424 = int32(0)
				v425 = l1
				v426 = l0
				v455 = v424 << (uint(int32(2)) % 32)
				v459 = int32(0)
				v462 = l3 + v455
				v463 = l4 + v455
				v467 = l2 - v424
				for {
					v498 = v426 + v459
					v499 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498))))
					v502 = v425 + v459
					v503 = int32(2)
					v505 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502+v503))))
					v507 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502))))
					v510 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498+v503))))
					v519 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462))))
					v520 = (v499*int32(9)+v505+(v507+v510)*int32(3)+int32(8))>>(uint(int32(4))%32) + v519
					if v520 < v44 {
						v522 = v520
					} else {
						v522 = v44
					}
					if v520 < int32(0) {
						v525 = int32(0)
					} else {
						v525 = v522
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v463))) = uint16(v525)
					v527 = int32(2)
					v543 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462+v527))))
					v544 = (v507+v510*int32(9)+(v505+v499)*int32(3)+int32(8))>>(uint(int32(4))%32) + v543
					if v544 < v44 {
						v546 = v544
					} else {
						v546 = v44
					}
					if v544 < int32(0) {
						v549 = int32(0)
					} else {
						v549 = v546
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v463+v527))) = uint16(v549)
					v551 = int32(4)
					v558 = v467 + int32(-1)
					if v558 != 0 {
						v459 = v459 + int32(2)
						v462 = v462 + v551
						v463 = v463 + v551
						v467 = v558
						continue
					} else {
						break
					}
					break
				}
			} else {
				if base.B2i32(base.Ui32(l4) < base.Ui32(l1+v50))&base.B2i32(base.Ui32(l1) < base.Ui32(v55)) != 0 {
					v424 = int32(0)
					v425 = l1
					v426 = l0
					v455 = v424 << (uint(int32(2)) % 32)
					v459 = int32(0)
					v462 = l3 + v455
					v463 = l4 + v455
					v467 = l2 - v424
					for {
						v498 = v426 + v459
						v499 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498))))
						v502 = v425 + v459
						v503 = int32(2)
						v505 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502+v503))))
						v507 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502))))
						v510 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498+v503))))
						v519 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462))))
						v520 = (v499*int32(9)+v505+(v507+v510)*int32(3)+int32(8))>>(uint(int32(4))%32) + v519
						if v520 < v44 {
							v522 = v520
						} else {
							v522 = v44
						}
						if v520 < int32(0) {
							v525 = int32(0)
						} else {
							v525 = v522
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v463))) = uint16(v525)
						v527 = int32(2)
						v543 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462+v527))))
						v544 = (v507+v510*int32(9)+(v505+v499)*int32(3)+int32(8))>>(uint(int32(4))%32) + v543
						if v544 < v44 {
							v546 = v544
						} else {
							v546 = v44
						}
						if v544 < int32(0) {
							v549 = int32(0)
						} else {
							v549 = v546
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v463+v527))) = uint16(v549)
						v551 = int32(4)
						v558 = v467 + int32(-1)
						if v558 != 0 {
							v459 = v459 + int32(2)
							v462 = v462 + v551
							v463 = v463 + v551
							v467 = v558
							continue
						} else {
							break
						}
						break
					}
				} else {
					if base.B2i32(base.Ui32(l4) < base.Ui32(l3+v54))&base.B2i32(base.Ui32(l3) < base.Ui32(v55)) != 0 {
						v424 = int32(0)
						v425 = l1
						v426 = l0
						v455 = v424 << (uint(int32(2)) % 32)
						v459 = int32(0)
						v462 = l3 + v455
						v463 = l4 + v455
						v467 = l2 - v424
						for {
							v498 = v426 + v459
							v499 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498))))
							v502 = v425 + v459
							v503 = int32(2)
							v505 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502+v503))))
							v507 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502))))
							v510 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498+v503))))
							v519 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462))))
							v520 = (v499*int32(9)+v505+(v507+v510)*int32(3)+int32(8))>>(uint(int32(4))%32) + v519
							if v520 < v44 {
								v522 = v520
							} else {
								v522 = v44
							}
							if v520 < int32(0) {
								v525 = int32(0)
							} else {
								v525 = v522
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v463))) = uint16(v525)
							v527 = int32(2)
							v543 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462+v527))))
							v544 = (v507+v510*int32(9)+(v505+v499)*int32(3)+int32(8))>>(uint(int32(4))%32) + v543
							if v544 < v44 {
								v546 = v544
							} else {
								v546 = v44
							}
							if v544 < int32(0) {
								v549 = int32(0)
							} else {
								v549 = v546
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v463+v527))) = uint16(v549)
							v551 = int32(4)
							v558 = v467 + int32(-1)
							if v558 != 0 {
								v459 = v459 + int32(2)
								v462 = v462 + v551
								v463 = v463 + v551
								v467 = v558
								continue
							} else {
								break
							}
							break
						}
					} else {
						v67 = l2 & int32(2147483640)
						v69 = v67 << (uint(int32(1)) % 32)
						v72 = base.Simd_g_i32x4_splat(v44)
						v75 = l0
						v76 = l1
						v83 = v67
						v88 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k0)
						v89 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k1)
						for {
							v113 = int32(1)
							v114 = base.Simd_g_i32x4_shl(v89, v113)
							v115 = int32(0)
							v118 = base.Simd_g_i32x4_extract_lane_l0(v114) << (uint(v113) % 32)
							v120 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k2)
							v122 = base.Simd_g_v128_load_rng(m, v75, v115, int32(0), int32(18))
							v123 = base.Simd_g_i32x4_extend_low_i16x8_s(v122)
							v124 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k3)
							v126 = int32(2)
							v129 = base.Simd_g_v128_load_rng(m, v76+v126, v115, int32(-2), int32(18))
							v130 = base.Simd_g_i32x4_extend_low_i16x8_s(v129)
							v133 = base.Simd_g_v128_load_nc(m, v76, v115)
							v134 = base.Simd_g_i32x4_extend_low_i16x8_s(v133)
							v138 = base.Simd_g_v128_load_nc(m, v75+v126, v115)
							v139 = base.Simd_g_i32x4_extend_low_i16x8_s(v138)
							v141 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k4)
							v144 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k5)
							v146 = int32(4)
							v148 = int32(3)
							v151 = base.Simd_g_i32x4_extract_lane_l3(v114) << (uint(v113) % 32)
							v156 = base.Simd_g_i32x4_extract_lane_l2(v114) << (uint(v113) % 32)
							v161 = base.Simd_g_i32x4_extract_lane_l1(v114) << (uint(v113) % 32)
							v165 = base.Simd_g_v128_load16_splat(m, l3+v118, v115)
							v168 = base.Simd_g_v128_load16_lane_l1(m, l3+v161, v115, v165)
							v171 = base.Simd_g_v128_load16_lane_l2(m, l3+v156, v115, v168)
							v174 = base.Simd_g_v128_load16_lane_l3(m, l3+v151, v115, v171)
							v176 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v123, v124), v130), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v134, v139), v141)), v144), v146), base.Simd_g_i32x4_extend_low_i16x8_u(v174))
							v178 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k6)
							v180 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k7)
							v182 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v122, v122, base.Simd_g_const(&F_SharpYuvFilterRow_C__k7), base.Simd_g_const(&F_SharpYuvFilterRow_C__k8)))
							v186 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v129, v122, base.Simd_g_const(&F_SharpYuvFilterRow_C__k7), base.Simd_g_const(&F_SharpYuvFilterRow_C__k8)))
							v190 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v133, v122, base.Simd_g_const(&F_SharpYuvFilterRow_C__k7), base.Simd_g_const(&F_SharpYuvFilterRow_C__k8)))
							v193 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v138, v122, base.Simd_g_const(&F_SharpYuvFilterRow_C__k7), base.Simd_g_const(&F_SharpYuvFilterRow_C__k8)))
							v201 = base.Simd_g_i32x4_shl(v88, v113)
							v205 = base.Simd_g_i32x4_extract_lane_l3(v201) << (uint(v113) % 32)
							v210 = base.Simd_g_i32x4_extract_lane_l2(v201) << (uint(v113) % 32)
							v215 = base.Simd_g_i32x4_extract_lane_l1(v201) << (uint(v113) % 32)
							v220 = base.Simd_g_i32x4_extract_lane_l0(v201) << (uint(v113) % 32)
							v223 = base.Simd_g_v128_load16_splat(m, l3+v220, v115)
							v226 = base.Simd_g_v128_load16_lane_l1(m, l3+v215, v115, v223)
							v229 = base.Simd_g_v128_load16_lane_l2(m, l3+v210, v115, v226)
							v232 = base.Simd_g_v128_load16_lane_l3(m, l3+v205, v115, v229)
							v234 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v182, v124), v186), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v190, v193), v141)), v144), v146), base.Simd_g_i32x4_extend_low_i16x8_u(v232))
							v241 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k9)
							v243 = base.Simd_g_v128_bitselect(v120, base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v176, v72), v178), base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v234, v72), v178)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v176, v120), base.Simd_g_i32x4_lt_s(v234, v120), base.Simd_g_const(&F_SharpYuvFilterRow_C__k10), base.Simd_g_const(&F_SharpYuvFilterRow_C__k11)))
							base.Simd_g_v128_store16_lane_l0(m, l4+v118, v115, v243)
							base.Simd_g_v128_store16_lane_l1(m, l4+v161, v115, v243)
							base.Simd_g_v128_store16_lane_l2(m, l4+v156, v115, v243)
							base.Simd_g_v128_store16_lane_l3(m, l4+v151, v115, v243)
							base.Simd_g_v128_store16_lane_l4(m, l4+v220, v115, v243)
							v265 = int32(5)
							base.Simd_g_v128_store16_lane_l5(m, l4+v215, v115, v243)
							v269 = int32(6)
							base.Simd_g_v128_store16_lane_l6(m, l4+v210, v115, v243)
							v273 = int32(7)
							base.Simd_g_v128_store16_lane_l7(m, l4+v205, v115, v243)
							v275 = base.Simd_g_const(&F_SharpYuvFilterRow_C__k12)
							v276 = base.Simd_g_v128_or(v114, v275)
							v280 = base.Simd_g_i32x4_extract_lane_l0(v276) << (uint(v113) % 32)
							v293 = base.Simd_g_i32x4_extract_lane_l3(v276) << (uint(v113) % 32)
							v298 = base.Simd_g_i32x4_extract_lane_l2(v276) << (uint(v113) % 32)
							v303 = base.Simd_g_i32x4_extract_lane_l1(v276) << (uint(v113) % 32)
							v307 = base.Simd_g_v128_load16_splat(m, l3+v280, v115)
							v310 = base.Simd_g_v128_load16_lane_l1(m, l3+v303, v115, v307)
							v313 = base.Simd_g_v128_load16_lane_l2(m, l3+v298, v115, v310)
							v316 = base.Simd_g_v128_load16_lane_l3(m, l3+v293, v115, v313)
							v318 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v139, v124), v134), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v130, v123), v141)), v144), v146), base.Simd_g_i32x4_extend_low_i16x8_u(v316))
							v329 = base.Simd_g_v128_or(v201, v275)
							v333 = base.Simd_g_i32x4_extract_lane_l3(v329) << (uint(v113) % 32)
							v338 = base.Simd_g_i32x4_extract_lane_l2(v329) << (uint(v113) % 32)
							v343 = base.Simd_g_i32x4_extract_lane_l1(v329) << (uint(v113) % 32)
							v348 = base.Simd_g_i32x4_extract_lane_l0(v329) << (uint(v113) % 32)
							v351 = base.Simd_g_v128_load16_splat(m, l3+v348, v115)
							v354 = base.Simd_g_v128_load16_lane_l1(m, l3+v343, v115, v351)
							v357 = base.Simd_g_v128_load16_lane_l2(m, l3+v338, v115, v354)
							v360 = base.Simd_g_v128_load16_lane_l3(m, l3+v333, v115, v357)
							v362 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v193, v124), v190), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v186, v182), v141)), v144), v146), base.Simd_g_i32x4_extend_low_i16x8_u(v360))
							v370 = base.Simd_g_v128_bitselect(v120, base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v318, v72), v178), base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v362, v72), v178)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v318, v120), base.Simd_g_i32x4_lt_s(v362, v120), base.Simd_g_const(&F_SharpYuvFilterRow_C__k10), base.Simd_g_const(&F_SharpYuvFilterRow_C__k11)))
							base.Simd_g_v128_store16_lane_l0(m, l4+v280, v115, v370)
							base.Simd_g_v128_store16_lane_l1(m, l4+v303, v115, v370)
							base.Simd_g_v128_store16_lane_l2(m, l4+v298, v115, v370)
							base.Simd_g_v128_store16_lane_l3(m, l4+v293, v115, v370)
							base.Simd_g_v128_store16_lane_l4(m, l4+v348, v115, v370)
							base.Simd_g_v128_store16_lane_l5(m, l4+v343, v115, v370)
							base.Simd_g_v128_store16_lane_l6(m, l4+v338, v115, v370)
							base.Simd_g_v128_store16_lane_l7(m, l4+v333, v115, v370)
							v403 = int32(16)
							v409 = v83 + int32(-8)
							if v409 != 0 {
								v75 = v75 + v403
								v76 = v76 + v403
								v83 = v409
								v88 = base.Simd_g_i32x4_add(v88, v144)
								v89 = base.Simd_g_i32x4_add(v89, v144)
								continue
							} else {
								break
							}
							break
						}
						if v67 == l2 {
						} else {
							v424 = v67
							v425 = l1 + v69
							v426 = l0 + v69
							v455 = v424 << (uint(int32(2)) % 32)
							v459 = int32(0)
							v462 = l3 + v455
							v463 = l4 + v455
							v467 = l2 - v424
							for {
								v498 = v426 + v459
								v499 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498))))
								v502 = v425 + v459
								v503 = int32(2)
								v505 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502+v503))))
								v507 = int32(*(*int16)(unsafe.Add(mBase, uint32(v502))))
								v510 = int32(*(*int16)(unsafe.Add(mBase, uint32(v498+v503))))
								v519 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462))))
								v520 = (v499*int32(9)+v505+(v507+v510)*int32(3)+int32(8))>>(uint(int32(4))%32) + v519
								if v520 < v44 {
									v522 = v520
								} else {
									v522 = v44
								}
								if v520 < int32(0) {
									v525 = int32(0)
								} else {
									v525 = v522
								}
								*(*uint16)(unsafe.Add(mBase, uint32(v463))) = uint16(v525)
								v527 = int32(2)
								v543 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v462+v527))))
								v544 = (v507+v510*int32(9)+(v505+v499)*int32(3)+int32(8))>>(uint(int32(4))%32) + v543
								if v544 < v44 {
									v546 = v544
								} else {
									v546 = v44
								}
								if v544 < int32(0) {
									v549 = int32(0)
								} else {
									v549 = v546
								}
								*(*uint16)(unsafe.Add(mBase, uint32(v463+v527))) = uint16(v549)
								v551 = int32(4)
								v558 = v467 + int32(-1)
								if v558 != 0 {
									v459 = v459 + int32(2)
									v462 = v462 + v551
									v463 = v463 + v551
									v467 = v558
									continue
								} else {
									break
								}
								break
							}
						}
					}
				}
			}
		}
	}
	return
}

var F_SharpYuvFilterRow_C__k0 = [2]uint64{0x500000004, 0x700000006}
var F_SharpYuvFilterRow_C__k1 = [2]uint64{0x100000000, 0x300000002}
var F_SharpYuvFilterRow_C__k2 = [2]uint64{0x0, 0x0}
var F_SharpYuvFilterRow_C__k3 = [2]uint64{0x900000009, 0x900000009}
var F_SharpYuvFilterRow_C__k4 = [2]uint64{0x300000003, 0x300000003}
var F_SharpYuvFilterRow_C__k5 = [2]uint64{0x800000008, 0x800000008}
var F_SharpYuvFilterRow_C__k6 = [2]uint64{0xffff0000ffff, 0xffff0000ffff}
var F_SharpYuvFilterRow_C__k7 = [2]uint64{0xf0e0d0c0b0a0908, 0x100010001000100}
var F_SharpYuvFilterRow_C__k8 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SharpYuvFilterRow_C__k9 = [2]uint64{0xd0c090805040100, 0x1d1c191815141110}
var F_SharpYuvFilterRow_C__k10 = [2]uint64{0xd0c090805040100, 0x8080808080808080}
var F_SharpYuvFilterRow_C__k11 = [2]uint64{0x8080808080808080, 0xd0c090805040100}
var F_SharpYuvFilterRow_C__k12 = [2]uint64{0x100000001, 0x100000001}

func F_SharpYuvFilterRow_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 base.V128
	_ = v41
	var v47 int32
	_ = v47
	var v54 int32
	_ = v54
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 base.V128
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v97 int32
	_ = v97
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v105 base.V128
	_ = v105
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v111 base.V128
	_ = v111
	var v112 int32
	_ = v112
	var v115 int32
	_ = v115
	var v119 base.V128
	_ = v119
	var v127 base.V128
	_ = v127
	var v132 base.V128
	_ = v132
	var v133 base.V128
	_ = v133
	var v140 base.V128
	_ = v140
	var v152 int32
	_ = v152
	var v162 int32
	_ = v162
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v235 base.V128
	_ = v235
	var v241 int32
	_ = v241
	var v245 int32
	_ = v245
	var v250 base.V128
	_ = v250
	var v256 int32
	_ = v256
	var v260 int32
	_ = v260
	var v261 int32
	_ = v261
	var v266 base.V128
	_ = v266
	var v272 base.V128
	_ = v272
	var v287 int32
	_ = v287
	var v288 base.V128
	_ = v288
	var v289 int32
	_ = v289
	var v292 int32
	_ = v292
	var v294 base.V128
	_ = v294
	var v296 base.V128
	_ = v296
	var v302 int32
	_ = v302
	var v303 int32
	_ = v303
	var v306 base.V128
	_ = v306
	var v307 base.V128
	_ = v307
	var v308 base.V128
	_ = v308
	var v312 base.V128
	_ = v312
	var v313 base.V128
	_ = v313
	var v314 int32
	_ = v314
	var v319 base.V128
	_ = v319
	var v321 base.V128
	_ = v321
	var v323 base.V128
	_ = v323
	var v324 base.V128
	_ = v324
	var v326 int32
	_ = v326
	var v331 int32
	_ = v331
	var v336 int32
	_ = v336
	var v341 int32
	_ = v341
	var v345 base.V128
	_ = v345
	var v348 base.V128
	_ = v348
	var v351 base.V128
	_ = v351
	var v354 base.V128
	_ = v354
	var v356 base.V128
	_ = v356
	var v358 base.V128
	_ = v358
	var v360 base.V128
	_ = v360
	var v365 base.V128
	_ = v365
	var v366 base.V128
	_ = v366
	var v371 base.V128
	_ = v371
	var v378 base.V128
	_ = v378
	var v380 base.V128
	_ = v380
	var v385 base.V128
	_ = v385
	var v389 int32
	_ = v389
	var v394 int32
	_ = v394
	var v399 int32
	_ = v399
	var v404 int32
	_ = v404
	var v407 base.V128
	_ = v407
	var v410 base.V128
	_ = v410
	var v413 base.V128
	_ = v413
	var v416 base.V128
	_ = v416
	var v418 base.V128
	_ = v418
	var v425 base.V128
	_ = v425
	var v427 base.V128
	_ = v427
	var v449 int32
	_ = v449
	var v453 int32
	_ = v453
	var v457 int32
	_ = v457
	var v459 base.V128
	_ = v459
	var v460 base.V128
	_ = v460
	var v464 int32
	_ = v464
	var v477 int32
	_ = v477
	var v482 int32
	_ = v482
	var v487 int32
	_ = v487
	var v491 base.V128
	_ = v491
	var v494 base.V128
	_ = v494
	var v497 base.V128
	_ = v497
	var v500 base.V128
	_ = v500
	var v502 base.V128
	_ = v502
	var v513 base.V128
	_ = v513
	var v517 int32
	_ = v517
	var v522 int32
	_ = v522
	var v527 int32
	_ = v527
	var v532 int32
	_ = v532
	var v535 base.V128
	_ = v535
	var v538 base.V128
	_ = v538
	var v541 base.V128
	_ = v541
	var v544 base.V128
	_ = v544
	var v546 base.V128
	_ = v546
	var v554 base.V128
	_ = v554
	var v587 int32
	_ = v587
	var v593 int32
	_ = v593
	var v603 int32
	_ = v603
	var v607 int32
	_ = v607
	var v608 int32
	_ = v608
	var v609 int32
	_ = v609
	var v610 int32
	_ = v610
	var v634 int32
	_ = v634
	var v635 base.V128
	_ = v635
	var v636 base.V128
	_ = v636
	var v638 int32
	_ = v638
	var v640 int32
	_ = v640
	var v643 base.V128
	_ = v643
	var v647 base.V128
	_ = v647
	var v648 base.V128
	_ = v648
	var v652 base.V128
	_ = v652
	var v658 base.V128
	_ = v658
	var v662 base.V128
	_ = v662
	var v663 base.V128
	_ = v663
	var v666 base.V128
	_ = v666
	var v667 int32
	_ = v667
	var v670 int32
	_ = v670
	var v671 base.V128
	_ = v671
	var v676 base.V128
	_ = v676
	var v681 base.V128
	_ = v681
	var v686 base.V128
	_ = v686
	var v693 base.V128
	_ = v693
	var v697 base.V128
	_ = v697
	var v700 int32
	_ = v700
	var v711 int32
	_ = v711
	var v732 int32
	_ = v732
	var v750 int32
	_ = v750
	var v753 int32
	_ = v753
	var v757 int32
	_ = v757
	var v760 int32
	_ = v760
	var v767 int32
	_ = v767
	var v768 int32
	_ = v768
	var v769 int32
	_ = v769
	var v770 int32
	_ = v770
	var v773 int32
	_ = v773
	var v777 int32
	_ = v777
	var v780 int32
	_ = v780
	var v781 int32
	_ = v781
	var v794 base.V128
	_ = v794
	var v800 int32
	_ = v800
	var v804 int32
	_ = v804
	var v806 base.V128
	_ = v806
	var v812 int32
	_ = v812
	var v816 int32
	_ = v816
	var v817 int32
	_ = v817
	var v822 base.V128
	_ = v822
	var v828 base.V128
	_ = v828
	var v843 int32
	_ = v843
	var v844 base.V128
	_ = v844
	var v845 int32
	_ = v845
	var v848 int32
	_ = v848
	var v850 base.V128
	_ = v850
	var v852 base.V128
	_ = v852
	var v858 int32
	_ = v858
	var v859 int32
	_ = v859
	var v862 base.V128
	_ = v862
	var v863 base.V128
	_ = v863
	var v864 base.V128
	_ = v864
	var v868 base.V128
	_ = v868
	var v869 base.V128
	_ = v869
	var v870 int32
	_ = v870
	var v875 base.V128
	_ = v875
	var v877 base.V128
	_ = v877
	var v879 base.V128
	_ = v879
	var v880 base.V128
	_ = v880
	var v882 int32
	_ = v882
	var v887 int32
	_ = v887
	var v892 int32
	_ = v892
	var v897 int32
	_ = v897
	var v901 base.V128
	_ = v901
	var v904 base.V128
	_ = v904
	var v907 base.V128
	_ = v907
	var v910 base.V128
	_ = v910
	var v912 base.V128
	_ = v912
	var v914 base.V128
	_ = v914
	var v916 base.V128
	_ = v916
	var v921 base.V128
	_ = v921
	var v922 base.V128
	_ = v922
	var v927 base.V128
	_ = v927
	var v934 base.V128
	_ = v934
	var v936 base.V128
	_ = v936
	var v941 base.V128
	_ = v941
	var v945 int32
	_ = v945
	var v950 int32
	_ = v950
	var v955 int32
	_ = v955
	var v960 int32
	_ = v960
	var v963 base.V128
	_ = v963
	var v966 base.V128
	_ = v966
	var v969 base.V128
	_ = v969
	var v972 base.V128
	_ = v972
	var v974 base.V128
	_ = v974
	var v981 base.V128
	_ = v981
	var v983 base.V128
	_ = v983
	var v1005 int32
	_ = v1005
	var v1009 int32
	_ = v1009
	var v1013 int32
	_ = v1013
	var v1015 base.V128
	_ = v1015
	var v1016 base.V128
	_ = v1016
	var v1020 int32
	_ = v1020
	var v1033 int32
	_ = v1033
	var v1038 int32
	_ = v1038
	var v1043 int32
	_ = v1043
	var v1047 base.V128
	_ = v1047
	var v1050 base.V128
	_ = v1050
	var v1053 base.V128
	_ = v1053
	var v1056 base.V128
	_ = v1056
	var v1058 base.V128
	_ = v1058
	var v1069 base.V128
	_ = v1069
	var v1073 int32
	_ = v1073
	var v1078 int32
	_ = v1078
	var v1083 int32
	_ = v1083
	var v1088 int32
	_ = v1088
	var v1091 base.V128
	_ = v1091
	var v1094 base.V128
	_ = v1094
	var v1097 base.V128
	_ = v1097
	var v1100 base.V128
	_ = v1100
	var v1102 base.V128
	_ = v1102
	var v1110 base.V128
	_ = v1110
	var v1143 int32
	_ = v1143
	var v1149 int32
	_ = v1149
	var v1170 int32
	_ = v1170
	var v1189 int32
	_ = v1189
	var v1191 int32
	_ = v1191
	var v1197 int32
	_ = v1197
	var v1198 int32
	_ = v1198
	var v1199 int32
	_ = v1199
	var v1203 int32
	_ = v1203
	var v1204 int32
	_ = v1204
	var v1230 int32
	_ = v1230
	var v1232 int32
	_ = v1232
	var v1233 int32
	_ = v1233
	var v1234 int32
	_ = v1234
	var v1235 int32
	_ = v1235
	var v1236 int32
	_ = v1236
	var v1239 int32
	_ = v1239
	var v1244 int32
	_ = v1244
	var v1245 int32
	_ = v1245
	var v1246 int32
	_ = v1246
	var v1249 int32
	_ = v1249
	var v1253 int32
	_ = v1253
	var v1254 int32
	_ = v1254
	var v1255 int32
	_ = v1255
	var v1257 int32
	_ = v1257
	var v1260 int32
	_ = v1260
	var v1262 int32
	_ = v1262
	var v1275 int32
	_ = v1275
	var v1276 int32
	_ = v1276
	var v1278 int32
	_ = v1278
	var v1281 int32
	_ = v1281
	var v1283 int32
	_ = v1283
	var v1288 int32
	_ = v1288
	var v1297 int32
	_ = v1297
	var v1327 int32
	_ = v1327
	var v1329 int32
	_ = v1329
	var v1335 int32
	_ = v1335
	var v1336 int32
	_ = v1336
	var v1337 int32
	_ = v1337
	var v1341 int32
	_ = v1341
	var v1342 int32
	_ = v1342
	var v1368 int32
	_ = v1368
	var v1370 int32
	_ = v1370
	var v1371 int32
	_ = v1371
	var v1372 int32
	_ = v1372
	var v1373 int32
	_ = v1373
	var v1374 int32
	_ = v1374
	var v1377 int32
	_ = v1377
	var v1382 int32
	_ = v1382
	var v1383 int32
	_ = v1383
	var v1384 int32
	_ = v1384
	var v1387 int32
	_ = v1387
	var v1391 int32
	_ = v1391
	var v1392 int32
	_ = v1392
	var v1393 int32
	_ = v1393
	var v1395 int32
	_ = v1395
	var v1398 int32
	_ = v1398
	var v1400 int32
	_ = v1400
	var v1413 int32
	_ = v1413
	var v1414 int32
	_ = v1414
	var v1416 int32
	_ = v1416
	var v1419 int32
	_ = v1419
	var v1421 int32
	_ = v1421
	var v1426 int32
	_ = v1426
	v37 = int32(-1)
	v40 = v37<<(uint(l5)%32) ^ v37
	v41 = base.Simd_g_i16x8_splat(v40)
	if int32(10) < l5 {
		if l2 < int32(4) {
			v732 = int32(0)
		} else {
			v603 = l0
			v607 = l1
			v608 = int32(0)
			v609 = l3
			v610 = l4
			for {
				v634 = int32(0)
				v635 = base.Simd_g_v128_load64_zero(m, v607, v634)
				v636 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k0)
				v638 = int32(16)
				v640 = int32(2)
				v643 = base.Simd_g_v128_load64_zero(m, v603+v640, v634)
				v647 = base.Simd_g_i32x4_shr_s(base.Simd_g_i8x16_shuffle2(v643, v643, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k1), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k2)), v638)
				v648 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i8x16_shuffle2(v635, v635, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k1), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k2)), v638), v647)
				v652 = base.Simd_g_v128_load64_zero(m, v607+v640, v634)
				v658 = base.Simd_g_v128_load64_zero(m, v603, v634)
				v662 = base.Simd_g_i32x4_shr_s(base.Simd_g_i8x16_shuffle2(v658, v658, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k1), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k2)), v638)
				v663 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i8x16_shuffle2(v652, v652, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k1), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k2)), v638), v662)
				v666 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v648, v663), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k3))
				v667 = int32(1)
				v670 = int32(3)
				v671 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(v666, base.Simd_g_i32x4_shl(v648, v667)), v670)
				v676 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(v671, v662), base.Simd_g_i32x4_shr_s(base.Simd_g_v128_xor(v671, v662), v667))
				v681 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v663, v667), v666), v670)
				v686 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(v681, v647), base.Simd_g_i32x4_shr_s(base.Simd_g_v128_xor(v681, v647), v667))
				v693 = base.Simd_g_v128_load(m, v609, v634)
				v697 = base.Simd_g_i16x8_max_s(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i8x16_shuffle2(v676, v686, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k4), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v676, v686, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k6), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k7))), v693), v41), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k8))
				base.Simd_g_v128_store(m, v610, v634, v697)
				v700 = int32(8)
				v711 = v608 + int32(4)
				if v608+v700 <= l2 {
					v603 = v603 + v700
					v607 = v607 + v700
					v608 = v711
					v609 = v609 + v638
					v610 = v610 + v638
					continue
				} else {
					break
				}
				break
			}
			v732 = v711
		}
		if l2 <= v732 {
		} else {
			v750 = l2 - v732
			if base.Ui32(v750) <= base.Ui32(int32(15)) {
				v1170 = v732
				v1189 = v1170 << (uint(int32(2)) % 32)
				v1191 = v1170 << (uint(int32(1)) % 32)
				v1197 = l3
				v1198 = l4
				v1199 = l0 + v1191
				v1203 = l1 + v1191
				v1204 = l2 - v1170
				for {
					v1230 = v1198 + v1189
					v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
					v1233 = int32(2)
					v1234 = v1199 + v1233
					v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
					v1236 = v1232 + v1235
					v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
					v1244 = v1203 + v1233
					v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
					v1246 = v1245 + v1239
					v1249 = v1246 + v1236 + int32(8)
					v1253 = v1197 + v1189
					v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
					v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
					if v1255 < v40 {
						v1257 = v1255
					} else {
						v1257 = v40
					}
					if v1255 < int32(0) {
						v1260 = int32(0)
					} else {
						v1260 = v1257
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
					v1262 = int32(2)
					v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
					v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
					if v1276 < v40 {
						v1278 = v1276
					} else {
						v1278 = v40
					}
					if v1276 < int32(0) {
						v1281 = int32(0)
					} else {
						v1281 = v1278
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
					v1283 = int32(4)
					v1288 = v1204 + int32(-1)
					if v1288 != 0 {
						v1197 = v1197 + v1283
						v1198 = v1198 + v1283
						v1199 = v1234
						v1203 = v1244
						v1204 = v1288
						continue
					} else {
						break
					}
					break
				}
			} else {
				v753 = int32(2)
				v757 = l4 + v732<<(uint(v753)%32) + v753
				v760 = v732 ^ int32(-1) + l2
				if base.Ui32(v757+v760<<(uint(v753)%32)) < base.Ui32(v757) {
					v1170 = v732
					v1189 = v1170 << (uint(int32(2)) % 32)
					v1191 = v1170 << (uint(int32(1)) % 32)
					v1197 = l3
					v1198 = l4
					v1199 = l0 + v1191
					v1203 = l1 + v1191
					v1204 = l2 - v1170
					for {
						v1230 = v1198 + v1189
						v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
						v1233 = int32(2)
						v1234 = v1199 + v1233
						v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
						v1236 = v1232 + v1235
						v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
						v1244 = v1203 + v1233
						v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
						v1246 = v1245 + v1239
						v1249 = v1246 + v1236 + int32(8)
						v1253 = v1197 + v1189
						v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
						v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
						if v1255 < v40 {
							v1257 = v1255
						} else {
							v1257 = v40
						}
						if v1255 < int32(0) {
							v1260 = int32(0)
						} else {
							v1260 = v1257
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
						v1262 = int32(2)
						v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
						v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
						if v1276 < v40 {
							v1278 = v1276
						} else {
							v1278 = v40
						}
						if v1276 < int32(0) {
							v1281 = int32(0)
						} else {
							v1281 = v1278
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
						v1283 = int32(4)
						v1288 = v1204 + int32(-1)
						if v1288 != 0 {
							v1197 = v1197 + v1283
							v1198 = v1198 + v1283
							v1199 = v1234
							v1203 = v1244
							v1204 = v1288
							continue
						} else {
							break
						}
						break
					}
				} else {
					if base.Ui32(int32(1073741823)) < base.Ui32(v760) {
						v1170 = v732
						v1189 = v1170 << (uint(int32(2)) % 32)
						v1191 = v1170 << (uint(int32(1)) % 32)
						v1197 = l3
						v1198 = l4
						v1199 = l0 + v1191
						v1203 = l1 + v1191
						v1204 = l2 - v1170
						for {
							v1230 = v1198 + v1189
							v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
							v1233 = int32(2)
							v1234 = v1199 + v1233
							v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
							v1236 = v1232 + v1235
							v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
							v1244 = v1203 + v1233
							v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
							v1246 = v1245 + v1239
							v1249 = v1246 + v1236 + int32(8)
							v1253 = v1197 + v1189
							v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
							v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
							if v1255 < v40 {
								v1257 = v1255
							} else {
								v1257 = v40
							}
							if v1255 < int32(0) {
								v1260 = int32(0)
							} else {
								v1260 = v1257
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
							v1262 = int32(2)
							v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
							v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
							if v1276 < v40 {
								v1278 = v1276
							} else {
								v1278 = v40
							}
							if v1276 < int32(0) {
								v1281 = int32(0)
							} else {
								v1281 = v1278
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
							v1283 = int32(4)
							v1288 = v1204 + int32(-1)
							if v1288 != 0 {
								v1197 = v1197 + v1283
								v1198 = v1198 + v1283
								v1199 = v1234
								v1203 = v1244
								v1204 = v1288
								continue
							} else {
								break
							}
							break
						}
					} else {
						v767 = int32(2)
						v768 = v732 << (uint(v767) % 32)
						v769 = l4 + v768
						v770 = int32(1)
						v773 = l2<<(uint(v770)%32) + v767
						v777 = v732 << (uint(v770) % 32)
						v780 = l2 << (uint(v767) % 32)
						v781 = l4 + v780
						if base.B2i32(base.Ui32(v769) < base.Ui32(l0+v773))&base.B2i32(base.Ui32(l0+v777) < base.Ui32(v781)) != 0 {
							v1170 = v732
							v1189 = v1170 << (uint(int32(2)) % 32)
							v1191 = v1170 << (uint(int32(1)) % 32)
							v1197 = l3
							v1198 = l4
							v1199 = l0 + v1191
							v1203 = l1 + v1191
							v1204 = l2 - v1170
							for {
								v1230 = v1198 + v1189
								v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
								v1233 = int32(2)
								v1234 = v1199 + v1233
								v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
								v1236 = v1232 + v1235
								v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
								v1244 = v1203 + v1233
								v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
								v1246 = v1245 + v1239
								v1249 = v1246 + v1236 + int32(8)
								v1253 = v1197 + v1189
								v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
								v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
								if v1255 < v40 {
									v1257 = v1255
								} else {
									v1257 = v40
								}
								if v1255 < int32(0) {
									v1260 = int32(0)
								} else {
									v1260 = v1257
								}
								*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
								v1262 = int32(2)
								v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
								v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
								if v1276 < v40 {
									v1278 = v1276
								} else {
									v1278 = v40
								}
								if v1276 < int32(0) {
									v1281 = int32(0)
								} else {
									v1281 = v1278
								}
								*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
								v1283 = int32(4)
								v1288 = v1204 + int32(-1)
								if v1288 != 0 {
									v1197 = v1197 + v1283
									v1198 = v1198 + v1283
									v1199 = v1234
									v1203 = v1244
									v1204 = v1288
									continue
								} else {
									break
								}
								break
							}
						} else {
							if base.B2i32(base.Ui32(v769) < base.Ui32(l1+v773))&base.B2i32(base.Ui32(l1+v777) < base.Ui32(v781)) != 0 {
								v1170 = v732
								v1189 = v1170 << (uint(int32(2)) % 32)
								v1191 = v1170 << (uint(int32(1)) % 32)
								v1197 = l3
								v1198 = l4
								v1199 = l0 + v1191
								v1203 = l1 + v1191
								v1204 = l2 - v1170
								for {
									v1230 = v1198 + v1189
									v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
									v1233 = int32(2)
									v1234 = v1199 + v1233
									v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
									v1236 = v1232 + v1235
									v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
									v1244 = v1203 + v1233
									v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
									v1246 = v1245 + v1239
									v1249 = v1246 + v1236 + int32(8)
									v1253 = v1197 + v1189
									v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
									v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
									if v1255 < v40 {
										v1257 = v1255
									} else {
										v1257 = v40
									}
									if v1255 < int32(0) {
										v1260 = int32(0)
									} else {
										v1260 = v1257
									}
									*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
									v1262 = int32(2)
									v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
									v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
									if v1276 < v40 {
										v1278 = v1276
									} else {
										v1278 = v40
									}
									if v1276 < int32(0) {
										v1281 = int32(0)
									} else {
										v1281 = v1278
									}
									*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
									v1283 = int32(4)
									v1288 = v1204 + int32(-1)
									if v1288 != 0 {
										v1197 = v1197 + v1283
										v1198 = v1198 + v1283
										v1199 = v1234
										v1203 = v1244
										v1204 = v1288
										continue
									} else {
										break
									}
									break
								}
							} else {
								if base.B2i32(base.Ui32(v769) < base.Ui32(l3+v780))&base.B2i32(base.Ui32(l3+v768) < base.Ui32(v781)) != 0 {
									v1170 = v732
									v1189 = v1170 << (uint(int32(2)) % 32)
									v1191 = v1170 << (uint(int32(1)) % 32)
									v1197 = l3
									v1198 = l4
									v1199 = l0 + v1191
									v1203 = l1 + v1191
									v1204 = l2 - v1170
									for {
										v1230 = v1198 + v1189
										v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
										v1233 = int32(2)
										v1234 = v1199 + v1233
										v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
										v1236 = v1232 + v1235
										v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
										v1244 = v1203 + v1233
										v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
										v1246 = v1245 + v1239
										v1249 = v1246 + v1236 + int32(8)
										v1253 = v1197 + v1189
										v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
										v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
										if v1255 < v40 {
											v1257 = v1255
										} else {
											v1257 = v40
										}
										if v1255 < int32(0) {
											v1260 = int32(0)
										} else {
											v1260 = v1257
										}
										*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
										v1262 = int32(2)
										v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
										v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
										if v1276 < v40 {
											v1278 = v1276
										} else {
											v1278 = v40
										}
										if v1276 < int32(0) {
											v1281 = int32(0)
										} else {
											v1281 = v1278
										}
										*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
										v1283 = int32(4)
										v1288 = v1204 + int32(-1)
										if v1288 != 0 {
											v1197 = v1197 + v1283
											v1198 = v1198 + v1283
											v1199 = v1234
											v1203 = v1244
											v1204 = v1288
											continue
										} else {
											break
										}
										break
									}
								} else {
									v794 = base.Simd_g_i32x4_splat(v732)
									v800 = v732 << (uint(int32(1)) % 32)
									v804 = v750 & int32(-8)
									v806 = base.Simd_g_i32x4_splat(v40)
									v812 = l1 + v800
									v816 = l0 + v800
									v817 = v804
									v822 = base.Simd_g_i32x4_add(v794, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k9))
									v828 = base.Simd_g_i32x4_add(v794, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k10))
									for {
										v843 = int32(1)
										v844 = base.Simd_g_i32x4_shl(v822, v843)
										v845 = int32(0)
										v848 = base.Simd_g_i32x4_extract_lane_l0(v844) << (uint(v843) % 32)
										v850 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k8)
										v852 = base.Simd_g_v128_load(m, v812, v845)
										v858 = int32(2)
										v859 = base.Simd_g_i32x4_extract_lane_l0(v822)<<(uint(v843)%32) + v858
										v862 = base.Simd_g_v128_load(m, l0+v859, v845)
										v863 = base.Simd_g_i32x4_extend_low_i16x8_s(v862)
										v864 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v852), v863)
										v868 = base.Simd_g_v128_load(m, v816, v845)
										v869 = base.Simd_g_i32x4_extend_low_i16x8_s(v868)
										v870 = int32(3)
										v875 = base.Simd_g_v128_load(m, l1+v859, v845)
										v877 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v875), v869)
										v879 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k3)
										v880 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v877, v864), v879)
										v882 = int32(4)
										v887 = base.Simd_g_i32x4_extract_lane_l3(v844) << (uint(v843) % 32)
										v892 = base.Simd_g_i32x4_extract_lane_l2(v844) << (uint(v843) % 32)
										v897 = base.Simd_g_i32x4_extract_lane_l1(v844) << (uint(v843) % 32)
										v901 = base.Simd_g_v128_load16_splat(m, l3+v848, v845)
										v904 = base.Simd_g_v128_load16_lane_l1(m, l3+v897, v845, v901)
										v907 = base.Simd_g_v128_load16_lane_l2(m, l3+v892, v845, v904)
										v910 = base.Simd_g_v128_load16_lane_l3(m, l3+v887, v845, v907)
										v912 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v864, v843), base.Simd_g_i32x4_shl(v869, v870)), v880), v882), base.Simd_g_i32x4_extend_low_i16x8_u(v910))
										v914 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k11)
										v916 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12)
										v921 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v862, v852, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13)))
										v922 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v852, v852, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13))), v921)
										v927 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v868, v922, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13)))
										v934 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v875, v922, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13))), v927)
										v936 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v934, v922), v879)
										v941 = base.Simd_g_i32x4_shl(v828, v843)
										v945 = base.Simd_g_i32x4_extract_lane_l3(v941) << (uint(v843) % 32)
										v950 = base.Simd_g_i32x4_extract_lane_l2(v941) << (uint(v843) % 32)
										v955 = base.Simd_g_i32x4_extract_lane_l1(v941) << (uint(v843) % 32)
										v960 = base.Simd_g_i32x4_extract_lane_l0(v941) << (uint(v843) % 32)
										v963 = base.Simd_g_v128_load16_splat(m, l3+v960, v845)
										v966 = base.Simd_g_v128_load16_lane_l1(m, l3+v955, v845, v963)
										v969 = base.Simd_g_v128_load16_lane_l2(m, l3+v950, v845, v966)
										v972 = base.Simd_g_v128_load16_lane_l3(m, l3+v945, v845, v969)
										v974 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v922, v843), base.Simd_g_i32x4_shl(v927, v870)), v936), v882), base.Simd_g_i32x4_extend_low_i16x8_u(v972))
										v981 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k14)
										v983 = base.Simd_g_v128_bitselect(v850, base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v912, v806), v914), base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v974, v806), v914)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v912, v850), base.Simd_g_i32x4_lt_s(v974, v850), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k15), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k16)))
										base.Simd_g_v128_store16_lane_l0(m, l4+v848, v845, v983)
										base.Simd_g_v128_store16_lane_l1(m, l4+v897, v845, v983)
										base.Simd_g_v128_store16_lane_l2(m, l4+v892, v845, v983)
										base.Simd_g_v128_store16_lane_l3(m, l4+v887, v845, v983)
										base.Simd_g_v128_store16_lane_l4(m, l4+v960, v845, v983)
										v1005 = int32(5)
										base.Simd_g_v128_store16_lane_l5(m, l4+v955, v845, v983)
										v1009 = int32(6)
										base.Simd_g_v128_store16_lane_l6(m, l4+v950, v845, v983)
										v1013 = int32(7)
										base.Simd_g_v128_store16_lane_l7(m, l4+v945, v845, v983)
										v1015 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k17)
										v1016 = base.Simd_g_v128_or(v844, v1015)
										v1020 = base.Simd_g_i32x4_extract_lane_l0(v1016) << (uint(v843) % 32)
										v1033 = base.Simd_g_i32x4_extract_lane_l3(v1016) << (uint(v843) % 32)
										v1038 = base.Simd_g_i32x4_extract_lane_l2(v1016) << (uint(v843) % 32)
										v1043 = base.Simd_g_i32x4_extract_lane_l1(v1016) << (uint(v843) % 32)
										v1047 = base.Simd_g_v128_load16_splat(m, l3+v1020, v845)
										v1050 = base.Simd_g_v128_load16_lane_l1(m, l3+v1043, v845, v1047)
										v1053 = base.Simd_g_v128_load16_lane_l2(m, l3+v1038, v845, v1050)
										v1056 = base.Simd_g_v128_load16_lane_l3(m, l3+v1033, v845, v1053)
										v1058 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v863, v870), base.Simd_g_i32x4_shl(v877, v843)), v880), v882), base.Simd_g_i32x4_extend_low_i16x8_u(v1056))
										v1069 = base.Simd_g_v128_or(v941, v1015)
										v1073 = base.Simd_g_i32x4_extract_lane_l3(v1069) << (uint(v843) % 32)
										v1078 = base.Simd_g_i32x4_extract_lane_l2(v1069) << (uint(v843) % 32)
										v1083 = base.Simd_g_i32x4_extract_lane_l1(v1069) << (uint(v843) % 32)
										v1088 = base.Simd_g_i32x4_extract_lane_l0(v1069) << (uint(v843) % 32)
										v1091 = base.Simd_g_v128_load16_splat(m, l3+v1088, v845)
										v1094 = base.Simd_g_v128_load16_lane_l1(m, l3+v1083, v845, v1091)
										v1097 = base.Simd_g_v128_load16_lane_l2(m, l3+v1078, v845, v1094)
										v1100 = base.Simd_g_v128_load16_lane_l3(m, l3+v1073, v845, v1097)
										v1102 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v921, v870), base.Simd_g_i32x4_shl(v934, v843)), v936), v882), base.Simd_g_i32x4_extend_low_i16x8_u(v1100))
										v1110 = base.Simd_g_v128_bitselect(v850, base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v1058, v806), v914), base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v1102, v806), v914)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v1058, v850), base.Simd_g_i32x4_lt_s(v1102, v850), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k15), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k16)))
										base.Simd_g_v128_store16_lane_l0(m, l4+v1020, v845, v1110)
										base.Simd_g_v128_store16_lane_l1(m, l4+v1043, v845, v1110)
										base.Simd_g_v128_store16_lane_l2(m, l4+v1038, v845, v1110)
										base.Simd_g_v128_store16_lane_l3(m, l4+v1033, v845, v1110)
										base.Simd_g_v128_store16_lane_l4(m, l4+v1088, v845, v1110)
										base.Simd_g_v128_store16_lane_l5(m, l4+v1083, v845, v1110)
										base.Simd_g_v128_store16_lane_l6(m, l4+v1078, v845, v1110)
										base.Simd_g_v128_store16_lane_l7(m, l4+v1073, v845, v1110)
										v1143 = int32(16)
										v1149 = v817 + int32(-8)
										if v1149 != 0 {
											v812 = v812 + v1143
											v816 = v816 + v1143
											v817 = v1149
											v822 = base.Simd_g_i32x4_add(v822, v879)
											v828 = base.Simd_g_i32x4_add(v828, v879)
											continue
										} else {
											break
										}
										break
									}
									if v750 == v804 {
									} else {
										v1170 = v732 + v804
										v1189 = v1170 << (uint(int32(2)) % 32)
										v1191 = v1170 << (uint(int32(1)) % 32)
										v1197 = l3
										v1198 = l4
										v1199 = l0 + v1191
										v1203 = l1 + v1191
										v1204 = l2 - v1170
										for {
											v1230 = v1198 + v1189
											v1232 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1203))))
											v1233 = int32(2)
											v1234 = v1199 + v1233
											v1235 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1234))))
											v1236 = v1232 + v1235
											v1239 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1199))))
											v1244 = v1203 + v1233
											v1245 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1244))))
											v1246 = v1245 + v1239
											v1249 = v1246 + v1236 + int32(8)
											v1253 = v1197 + v1189
											v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253))))
											v1255 = (v1236<<(uint(int32(1))%32)+v1239<<(uint(int32(3))%32)+v1249)>>(uint(int32(4))%32) + v1254
											if v1255 < v40 {
												v1257 = v1255
											} else {
												v1257 = v40
											}
											if v1255 < int32(0) {
												v1260 = int32(0)
											} else {
												v1260 = v1257
											}
											*(*uint16)(unsafe.Add(mBase, uint32(v1230))) = uint16(v1260)
											v1262 = int32(2)
											v1275 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1253+v1262))))
											v1276 = (v1235<<(uint(int32(3))%32)+v1246<<(uint(int32(1))%32)+v1249)>>(uint(int32(4))%32) + v1275
											if v1276 < v40 {
												v1278 = v1276
											} else {
												v1278 = v40
											}
											if v1276 < int32(0) {
												v1281 = int32(0)
											} else {
												v1281 = v1278
											}
											*(*uint16)(unsafe.Add(mBase, uint32(v1230+v1262))) = uint16(v1281)
											v1283 = int32(4)
											v1288 = v1204 + int32(-1)
											if v1288 != 0 {
												v1197 = v1197 + v1283
												v1198 = v1198 + v1283
												v1199 = v1234
												v1203 = v1244
												v1204 = v1288
												continue
											} else {
												break
											}
											break
										}
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		if l2 < int32(8) {
			v162 = int32(0)
		} else {
			v47 = int32(0)
			v54 = v47
			v58 = l0
			v59 = l1
			v60 = v47
			for {
				v85 = l3 + v54
				v86 = int32(0)
				v87 = base.Simd_g_v128_load_rng(m, v85, v86, int32(0), int32(32))
				v88 = l4 + v54
				v89 = int32(16)
				v94 = base.Simd_g_v128_load_nc(m, v85+v89, v86)
				v96 = base.Simd_g_v128_load_rng(m, v59, v86, int32(0), int32(18))
				v97 = int32(2)
				v100 = base.Simd_g_v128_load_rng(m, v58+v97, v86, int32(-2), int32(18))
				v101 = base.Simd_g_i16x8_add(v96, v100)
				v105 = base.Simd_g_v128_load_nc(m, v59+v97, v86)
				v107 = base.Simd_g_v128_load_nc(m, v58, v86)
				v108 = base.Simd_g_i16x8_add(v105, v107)
				v111 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v101, v108), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k18))
				v112 = int32(1)
				v115 = int32(3)
				v119 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v111, base.Simd_g_i16x8_shl(v101, v112)), v115), v107), v112)
				v127 = base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_shl(v108, v112), v111), v115), v100), v112)
				v132 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k8)
				v133 = base.Simd_g_i16x8_max_s(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_add(v94, base.Simd_g_i8x16_shuffle2(v119, v127, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k19), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k20))), v41), v132)
				base.Simd_g_v128_store(m, v88+v89, v86, v133)
				v140 = base.Simd_g_i16x8_max_s(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_add(v87, base.Simd_g_i8x16_shuffle2(v119, v127, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k1), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k2))), v41), v132)
				base.Simd_g_v128_store(m, v88, v86, v140)
				v152 = v60 + int32(8)
				if v60+v89 <= l2 {
					v54 = v54 + int32(32)
					v58 = v58 + v89
					v59 = v59 + v89
					v60 = v152
					continue
				} else {
					break
				}
				break
			}
			v162 = v152
		}
		if l2 <= v162 {
		} else {
			v191 = l2 - v162
			if base.Ui32(v191) <= base.Ui32(int32(15)) {
				v1297 = v162
				v1327 = v1297 << (uint(int32(2)) % 32)
				v1329 = v1297 << (uint(int32(1)) % 32)
				v1335 = l3
				v1336 = l4
				v1337 = l0 + v1329
				v1341 = l1 + v1329
				v1342 = l2 - v1297
				for {
					v1368 = v1336 + v1327
					v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
					v1371 = int32(2)
					v1372 = v1337 + v1371
					v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
					v1374 = v1370 + v1373
					v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
					v1382 = v1341 + v1371
					v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
					v1384 = v1383 + v1377
					v1387 = v1384 + v1374 + int32(8)
					v1391 = v1335 + v1327
					v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
					v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
					if v1393 < v40 {
						v1395 = v1393
					} else {
						v1395 = v40
					}
					if v1393 < int32(0) {
						v1398 = int32(0)
					} else {
						v1398 = v1395
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
					v1400 = int32(2)
					v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
					v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
					if v1414 < v40 {
						v1416 = v1414
					} else {
						v1416 = v40
					}
					if v1414 < int32(0) {
						v1419 = int32(0)
					} else {
						v1419 = v1416
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
					v1421 = int32(4)
					v1426 = v1342 + int32(-1)
					if v1426 != 0 {
						v1335 = v1335 + v1421
						v1336 = v1336 + v1421
						v1337 = v1372
						v1341 = v1382
						v1342 = v1426
						continue
					} else {
						break
					}
					break
				}
			} else {
				v194 = int32(2)
				v198 = l4 + v162<<(uint(v194)%32) + v194
				v201 = v162 ^ int32(-1) + l2
				if base.Ui32(v198+v201<<(uint(v194)%32)) < base.Ui32(v198) {
					v1297 = v162
					v1327 = v1297 << (uint(int32(2)) % 32)
					v1329 = v1297 << (uint(int32(1)) % 32)
					v1335 = l3
					v1336 = l4
					v1337 = l0 + v1329
					v1341 = l1 + v1329
					v1342 = l2 - v1297
					for {
						v1368 = v1336 + v1327
						v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
						v1371 = int32(2)
						v1372 = v1337 + v1371
						v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
						v1374 = v1370 + v1373
						v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
						v1382 = v1341 + v1371
						v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
						v1384 = v1383 + v1377
						v1387 = v1384 + v1374 + int32(8)
						v1391 = v1335 + v1327
						v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
						v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
						if v1393 < v40 {
							v1395 = v1393
						} else {
							v1395 = v40
						}
						if v1393 < int32(0) {
							v1398 = int32(0)
						} else {
							v1398 = v1395
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
						v1400 = int32(2)
						v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
						v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
						if v1414 < v40 {
							v1416 = v1414
						} else {
							v1416 = v40
						}
						if v1414 < int32(0) {
							v1419 = int32(0)
						} else {
							v1419 = v1416
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
						v1421 = int32(4)
						v1426 = v1342 + int32(-1)
						if v1426 != 0 {
							v1335 = v1335 + v1421
							v1336 = v1336 + v1421
							v1337 = v1372
							v1341 = v1382
							v1342 = v1426
							continue
						} else {
							break
						}
						break
					}
				} else {
					if base.Ui32(int32(1073741823)) < base.Ui32(v201) {
						v1297 = v162
						v1327 = v1297 << (uint(int32(2)) % 32)
						v1329 = v1297 << (uint(int32(1)) % 32)
						v1335 = l3
						v1336 = l4
						v1337 = l0 + v1329
						v1341 = l1 + v1329
						v1342 = l2 - v1297
						for {
							v1368 = v1336 + v1327
							v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
							v1371 = int32(2)
							v1372 = v1337 + v1371
							v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
							v1374 = v1370 + v1373
							v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
							v1382 = v1341 + v1371
							v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
							v1384 = v1383 + v1377
							v1387 = v1384 + v1374 + int32(8)
							v1391 = v1335 + v1327
							v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
							v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
							if v1393 < v40 {
								v1395 = v1393
							} else {
								v1395 = v40
							}
							if v1393 < int32(0) {
								v1398 = int32(0)
							} else {
								v1398 = v1395
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
							v1400 = int32(2)
							v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
							v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
							if v1414 < v40 {
								v1416 = v1414
							} else {
								v1416 = v40
							}
							if v1414 < int32(0) {
								v1419 = int32(0)
							} else {
								v1419 = v1416
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
							v1421 = int32(4)
							v1426 = v1342 + int32(-1)
							if v1426 != 0 {
								v1335 = v1335 + v1421
								v1336 = v1336 + v1421
								v1337 = v1372
								v1341 = v1382
								v1342 = v1426
								continue
							} else {
								break
							}
							break
						}
					} else {
						v208 = int32(2)
						v209 = v162 << (uint(v208) % 32)
						v210 = l4 + v209
						v211 = int32(1)
						v214 = l2<<(uint(v211)%32) + v208
						v218 = v162 << (uint(v211) % 32)
						v221 = l2 << (uint(v208) % 32)
						v222 = l4 + v221
						if base.B2i32(base.Ui32(v210) < base.Ui32(l0+v214))&base.B2i32(base.Ui32(l0+v218) < base.Ui32(v222)) != 0 {
							v1297 = v162
							v1327 = v1297 << (uint(int32(2)) % 32)
							v1329 = v1297 << (uint(int32(1)) % 32)
							v1335 = l3
							v1336 = l4
							v1337 = l0 + v1329
							v1341 = l1 + v1329
							v1342 = l2 - v1297
							for {
								v1368 = v1336 + v1327
								v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
								v1371 = int32(2)
								v1372 = v1337 + v1371
								v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
								v1374 = v1370 + v1373
								v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
								v1382 = v1341 + v1371
								v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
								v1384 = v1383 + v1377
								v1387 = v1384 + v1374 + int32(8)
								v1391 = v1335 + v1327
								v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
								v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
								if v1393 < v40 {
									v1395 = v1393
								} else {
									v1395 = v40
								}
								if v1393 < int32(0) {
									v1398 = int32(0)
								} else {
									v1398 = v1395
								}
								*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
								v1400 = int32(2)
								v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
								v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
								if v1414 < v40 {
									v1416 = v1414
								} else {
									v1416 = v40
								}
								if v1414 < int32(0) {
									v1419 = int32(0)
								} else {
									v1419 = v1416
								}
								*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
								v1421 = int32(4)
								v1426 = v1342 + int32(-1)
								if v1426 != 0 {
									v1335 = v1335 + v1421
									v1336 = v1336 + v1421
									v1337 = v1372
									v1341 = v1382
									v1342 = v1426
									continue
								} else {
									break
								}
								break
							}
						} else {
							if base.B2i32(base.Ui32(v210) < base.Ui32(l1+v214))&base.B2i32(base.Ui32(l1+v218) < base.Ui32(v222)) != 0 {
								v1297 = v162
								v1327 = v1297 << (uint(int32(2)) % 32)
								v1329 = v1297 << (uint(int32(1)) % 32)
								v1335 = l3
								v1336 = l4
								v1337 = l0 + v1329
								v1341 = l1 + v1329
								v1342 = l2 - v1297
								for {
									v1368 = v1336 + v1327
									v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
									v1371 = int32(2)
									v1372 = v1337 + v1371
									v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
									v1374 = v1370 + v1373
									v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
									v1382 = v1341 + v1371
									v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
									v1384 = v1383 + v1377
									v1387 = v1384 + v1374 + int32(8)
									v1391 = v1335 + v1327
									v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
									v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
									if v1393 < v40 {
										v1395 = v1393
									} else {
										v1395 = v40
									}
									if v1393 < int32(0) {
										v1398 = int32(0)
									} else {
										v1398 = v1395
									}
									*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
									v1400 = int32(2)
									v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
									v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
									if v1414 < v40 {
										v1416 = v1414
									} else {
										v1416 = v40
									}
									if v1414 < int32(0) {
										v1419 = int32(0)
									} else {
										v1419 = v1416
									}
									*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
									v1421 = int32(4)
									v1426 = v1342 + int32(-1)
									if v1426 != 0 {
										v1335 = v1335 + v1421
										v1336 = v1336 + v1421
										v1337 = v1372
										v1341 = v1382
										v1342 = v1426
										continue
									} else {
										break
									}
									break
								}
							} else {
								if base.B2i32(base.Ui32(v210) < base.Ui32(l3+v221))&base.B2i32(base.Ui32(l3+v209) < base.Ui32(v222)) != 0 {
									v1297 = v162
									v1327 = v1297 << (uint(int32(2)) % 32)
									v1329 = v1297 << (uint(int32(1)) % 32)
									v1335 = l3
									v1336 = l4
									v1337 = l0 + v1329
									v1341 = l1 + v1329
									v1342 = l2 - v1297
									for {
										v1368 = v1336 + v1327
										v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
										v1371 = int32(2)
										v1372 = v1337 + v1371
										v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
										v1374 = v1370 + v1373
										v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
										v1382 = v1341 + v1371
										v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
										v1384 = v1383 + v1377
										v1387 = v1384 + v1374 + int32(8)
										v1391 = v1335 + v1327
										v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
										v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
										if v1393 < v40 {
											v1395 = v1393
										} else {
											v1395 = v40
										}
										if v1393 < int32(0) {
											v1398 = int32(0)
										} else {
											v1398 = v1395
										}
										*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
										v1400 = int32(2)
										v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
										v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
										if v1414 < v40 {
											v1416 = v1414
										} else {
											v1416 = v40
										}
										if v1414 < int32(0) {
											v1419 = int32(0)
										} else {
											v1419 = v1416
										}
										*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
										v1421 = int32(4)
										v1426 = v1342 + int32(-1)
										if v1426 != 0 {
											v1335 = v1335 + v1421
											v1336 = v1336 + v1421
											v1337 = v1372
											v1341 = v1382
											v1342 = v1426
											continue
										} else {
											break
										}
										break
									}
								} else {
									v235 = base.Simd_g_i32x4_splat(v162)
									v241 = v162 << (uint(int32(1)) % 32)
									v245 = l2 & int32(7)
									v250 = base.Simd_g_i32x4_splat(v40)
									v256 = l1 + v241
									v260 = l0 + v241
									v261 = v162 + v245 - l2
									v266 = base.Simd_g_v128_or(v235, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k9))
									v272 = base.Simd_g_v128_or(v235, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k10))
									for {
										v287 = int32(1)
										v288 = base.Simd_g_i32x4_shl(v266, v287)
										v289 = int32(0)
										v292 = base.Simd_g_i32x4_extract_lane_l0(v288) << (uint(v287) % 32)
										v294 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k8)
										v296 = base.Simd_g_v128_load(m, v256, v289)
										v302 = int32(2)
										v303 = base.Simd_g_i32x4_extract_lane_l0(v266)<<(uint(v287)%32) + v302
										v306 = base.Simd_g_v128_load(m, l0+v303, v289)
										v307 = base.Simd_g_i32x4_extend_low_i16x8_s(v306)
										v308 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v296), v307)
										v312 = base.Simd_g_v128_load(m, v260, v289)
										v313 = base.Simd_g_i32x4_extend_low_i16x8_s(v312)
										v314 = int32(3)
										v319 = base.Simd_g_v128_load(m, l1+v303, v289)
										v321 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v319), v313)
										v323 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k3)
										v324 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v321, v308), v323)
										v326 = int32(4)
										v331 = base.Simd_g_i32x4_extract_lane_l3(v288) << (uint(v287) % 32)
										v336 = base.Simd_g_i32x4_extract_lane_l2(v288) << (uint(v287) % 32)
										v341 = base.Simd_g_i32x4_extract_lane_l1(v288) << (uint(v287) % 32)
										v345 = base.Simd_g_v128_load16_splat(m, l3+v292, v289)
										v348 = base.Simd_g_v128_load16_lane_l1(m, l3+v341, v289, v345)
										v351 = base.Simd_g_v128_load16_lane_l2(m, l3+v336, v289, v348)
										v354 = base.Simd_g_v128_load16_lane_l3(m, l3+v331, v289, v351)
										v356 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v308, v287), base.Simd_g_i32x4_shl(v313, v314)), v324), v326), base.Simd_g_i32x4_extend_low_i16x8_u(v354))
										v358 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k11)
										v360 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12)
										v365 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v306, v296, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13)))
										v366 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v296, v296, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13))), v365)
										v371 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v312, v366, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13)))
										v378 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_shuffle2(v319, v366, base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k12), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k13))), v371)
										v380 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v378, v366), v323)
										v385 = base.Simd_g_i32x4_shl(v272, v287)
										v389 = base.Simd_g_i32x4_extract_lane_l3(v385) << (uint(v287) % 32)
										v394 = base.Simd_g_i32x4_extract_lane_l2(v385) << (uint(v287) % 32)
										v399 = base.Simd_g_i32x4_extract_lane_l1(v385) << (uint(v287) % 32)
										v404 = base.Simd_g_i32x4_extract_lane_l0(v385) << (uint(v287) % 32)
										v407 = base.Simd_g_v128_load16_splat(m, l3+v404, v289)
										v410 = base.Simd_g_v128_load16_lane_l1(m, l3+v399, v289, v407)
										v413 = base.Simd_g_v128_load16_lane_l2(m, l3+v394, v289, v410)
										v416 = base.Simd_g_v128_load16_lane_l3(m, l3+v389, v289, v413)
										v418 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v366, v287), base.Simd_g_i32x4_shl(v371, v314)), v380), v326), base.Simd_g_i32x4_extend_low_i16x8_u(v416))
										v425 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k14)
										v427 = base.Simd_g_v128_bitselect(v294, base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v356, v250), v358), base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v418, v250), v358)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v356, v294), base.Simd_g_i32x4_lt_s(v418, v294), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k15), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k16)))
										base.Simd_g_v128_store16_lane_l0(m, l4+v292, v289, v427)
										base.Simd_g_v128_store16_lane_l1(m, l4+v341, v289, v427)
										base.Simd_g_v128_store16_lane_l2(m, l4+v336, v289, v427)
										base.Simd_g_v128_store16_lane_l3(m, l4+v331, v289, v427)
										base.Simd_g_v128_store16_lane_l4(m, l4+v404, v289, v427)
										v449 = int32(5)
										base.Simd_g_v128_store16_lane_l5(m, l4+v399, v289, v427)
										v453 = int32(6)
										base.Simd_g_v128_store16_lane_l6(m, l4+v394, v289, v427)
										v457 = int32(7)
										base.Simd_g_v128_store16_lane_l7(m, l4+v389, v289, v427)
										v459 = base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k17)
										v460 = base.Simd_g_v128_or(v288, v459)
										v464 = base.Simd_g_i32x4_extract_lane_l0(v460) << (uint(v287) % 32)
										v477 = base.Simd_g_i32x4_extract_lane_l3(v460) << (uint(v287) % 32)
										v482 = base.Simd_g_i32x4_extract_lane_l2(v460) << (uint(v287) % 32)
										v487 = base.Simd_g_i32x4_extract_lane_l1(v460) << (uint(v287) % 32)
										v491 = base.Simd_g_v128_load16_splat(m, l3+v464, v289)
										v494 = base.Simd_g_v128_load16_lane_l1(m, l3+v487, v289, v491)
										v497 = base.Simd_g_v128_load16_lane_l2(m, l3+v482, v289, v494)
										v500 = base.Simd_g_v128_load16_lane_l3(m, l3+v477, v289, v497)
										v502 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v307, v314), base.Simd_g_i32x4_shl(v321, v287)), v324), v326), base.Simd_g_i32x4_extend_low_i16x8_u(v500))
										v513 = base.Simd_g_v128_or(v385, v459)
										v517 = base.Simd_g_i32x4_extract_lane_l3(v513) << (uint(v287) % 32)
										v522 = base.Simd_g_i32x4_extract_lane_l2(v513) << (uint(v287) % 32)
										v527 = base.Simd_g_i32x4_extract_lane_l1(v513) << (uint(v287) % 32)
										v532 = base.Simd_g_i32x4_extract_lane_l0(v513) << (uint(v287) % 32)
										v535 = base.Simd_g_v128_load16_splat(m, l3+v532, v289)
										v538 = base.Simd_g_v128_load16_lane_l1(m, l3+v527, v289, v535)
										v541 = base.Simd_g_v128_load16_lane_l2(m, l3+v522, v289, v538)
										v544 = base.Simd_g_v128_load16_lane_l3(m, l3+v517, v289, v541)
										v546 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_shl(v365, v314), base.Simd_g_i32x4_shl(v378, v287)), v380), v326), base.Simd_g_i32x4_extend_low_i16x8_u(v544))
										v554 = base.Simd_g_v128_bitselect(v294, base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v502, v250), v358), base.Simd_g_v128_and(base.Simd_g_i32x4_min_s(v546, v250), v358)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v502, v294), base.Simd_g_i32x4_lt_s(v546, v294), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k15), base.Simd_g_const(&F_SharpYuvFilterRow_SSE2__k16)))
										base.Simd_g_v128_store16_lane_l0(m, l4+v464, v289, v554)
										base.Simd_g_v128_store16_lane_l1(m, l4+v487, v289, v554)
										base.Simd_g_v128_store16_lane_l2(m, l4+v482, v289, v554)
										base.Simd_g_v128_store16_lane_l3(m, l4+v477, v289, v554)
										base.Simd_g_v128_store16_lane_l4(m, l4+v532, v289, v554)
										base.Simd_g_v128_store16_lane_l5(m, l4+v527, v289, v554)
										base.Simd_g_v128_store16_lane_l6(m, l4+v522, v289, v554)
										base.Simd_g_v128_store16_lane_l7(m, l4+v517, v289, v554)
										v587 = int32(16)
										v593 = v261 + int32(8)
										if v593 != 0 {
											v256 = v256 + v587
											v260 = v260 + v587
											v261 = v593
											v266 = base.Simd_g_i32x4_add(v266, v323)
											v272 = base.Simd_g_i32x4_add(v272, v323)
											continue
										} else {
											break
										}
										break
									}
									if v245 != 0 {
										v1297 = v162 + (v191 - v245)
										v1327 = v1297 << (uint(int32(2)) % 32)
										v1329 = v1297 << (uint(int32(1)) % 32)
										v1335 = l3
										v1336 = l4
										v1337 = l0 + v1329
										v1341 = l1 + v1329
										v1342 = l2 - v1297
										for {
											v1368 = v1336 + v1327
											v1370 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1341))))
											v1371 = int32(2)
											v1372 = v1337 + v1371
											v1373 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1372))))
											v1374 = v1370 + v1373
											v1377 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1337))))
											v1382 = v1341 + v1371
											v1383 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1382))))
											v1384 = v1383 + v1377
											v1387 = v1384 + v1374 + int32(8)
											v1391 = v1335 + v1327
											v1392 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391))))
											v1393 = (v1374<<(uint(int32(1))%32)+v1377<<(uint(int32(3))%32)+v1387)>>(uint(int32(4))%32) + v1392
											if v1393 < v40 {
												v1395 = v1393
											} else {
												v1395 = v40
											}
											if v1393 < int32(0) {
												v1398 = int32(0)
											} else {
												v1398 = v1395
											}
											*(*uint16)(unsafe.Add(mBase, uint32(v1368))) = uint16(v1398)
											v1400 = int32(2)
											v1413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1391+v1400))))
											v1414 = (v1373<<(uint(int32(3))%32)+v1384<<(uint(int32(1))%32)+v1387)>>(uint(int32(4))%32) + v1413
											if v1414 < v40 {
												v1416 = v1414
											} else {
												v1416 = v40
											}
											if v1414 < int32(0) {
												v1419 = int32(0)
											} else {
												v1419 = v1416
											}
											*(*uint16)(unsafe.Add(mBase, uint32(v1368+v1400))) = uint16(v1419)
											v1421 = int32(4)
											v1426 = v1342 + int32(-1)
											if v1426 != 0 {
												v1335 = v1335 + v1421
												v1336 = v1336 + v1421
												v1337 = v1372
												v1341 = v1382
												v1342 = v1426
												continue
											} else {
												break
											}
											break
										}
									} else {
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

var F_SharpYuvFilterRow_SSE2__k0 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_SharpYuvFilterRow_SSE2__k1 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_SharpYuvFilterRow_SSE2__k2 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_SharpYuvFilterRow_SSE2__k3 = [2]uint64{0x800000008, 0x800000008}
var F_SharpYuvFilterRow_SSE2__k4 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_SharpYuvFilterRow_SSE2__k5 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_SharpYuvFilterRow_SSE2__k6 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_SharpYuvFilterRow_SSE2__k7 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_SharpYuvFilterRow_SSE2__k8 = [2]uint64{0x0, 0x0}
var F_SharpYuvFilterRow_SSE2__k9 = [2]uint64{0x100000000, 0x300000002}
var F_SharpYuvFilterRow_SSE2__k10 = [2]uint64{0x500000004, 0x700000006}
var F_SharpYuvFilterRow_SSE2__k11 = [2]uint64{0xffff0000ffff, 0xffff0000ffff}
var F_SharpYuvFilterRow_SSE2__k12 = [2]uint64{0xf0e0d0c0b0a0908, 0x100010001000100}
var F_SharpYuvFilterRow_SSE2__k13 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SharpYuvFilterRow_SSE2__k14 = [2]uint64{0xd0c090805040100, 0x1d1c191815141110}
var F_SharpYuvFilterRow_SSE2__k15 = [2]uint64{0xd0c090805040100, 0x8080808080808080}
var F_SharpYuvFilterRow_SSE2__k16 = [2]uint64{0x8080808080808080, 0xd0c090805040100}
var F_SharpYuvFilterRow_SSE2__k17 = [2]uint64{0x100000001, 0x100000001}
var F_SharpYuvFilterRow_SSE2__k18 = [2]uint64{0x8000800080008, 0x8000800080008}
var F_SharpYuvFilterRow_SSE2__k19 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_SharpYuvFilterRow_SSE2__k20 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}

func F_SharpYuvUpdateRGB_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v27 int32
	_ = v27
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v43 base.V128
	_ = v43
	var v44 base.V128
	_ = v44
	var v47 int32
	_ = v47
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v127 int32
	_ = v127
	if l3 < int32(1) {
	} else {
		v12 = int32(0)
		if base.Ui32(l3) < base.Ui32(int32(16)) {
			v60 = v12
			if l3&int32(1) == int32(0) {
				v82 = v60
			} else {
				v69 = int32(1)
				v70 = v60 << (uint(v69) % 32)
				v71 = l2 + v70
				v73 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v70))))
				v75 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v70))))
				v77 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v71))))
				v78 = v73 - v75 + v77
				*(*uint16)(unsafe.Add(mBase, uint32(v71))) = uint16(v78)
				v82 = v60 | v69
			}
			if v60 == l3+int32(-1) {
			} else {
				v89 = v82 << (uint(int32(1)) % 32)
				v98 = l2 + v89
				v99 = l0 + v89
				v100 = l3 - v82
				v101 = l1 + v89
				for {
					v102 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99))))
					v103 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101))))
					v105 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v98))))
					v106 = v102 - v103 + v105
					*(*uint16)(unsafe.Add(mBase, uint32(v98))) = uint16(v106)
					v108 = int32(2)
					v109 = v98 + v108
					v112 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99+v108))))
					v115 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101+v108))))
					v117 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v109))))
					v118 = v112 - v115 + v117
					*(*uint16)(unsafe.Add(mBase, uint32(v109))) = uint16(v118)
					v120 = int32(4)
					v127 = v100 + int32(-2)
					if v127 != 0 {
						v98 = v98 + v120
						v99 = v99 + v120
						v100 = v127
						v101 = v101 + v120
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v16 = l3 << (uint(int32(1)) % 32)
			v19 = l2 + v16
			if base.B2i32(base.Ui32(l2) < base.Ui32(l0+v16))&base.B2i32(base.Ui32(l0) < base.Ui32(v19)) != 0 {
				v60 = v12
				if l3&int32(1) == int32(0) {
					v82 = v60
				} else {
					v69 = int32(1)
					v70 = v60 << (uint(v69) % 32)
					v71 = l2 + v70
					v73 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v70))))
					v75 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v70))))
					v77 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v71))))
					v78 = v73 - v75 + v77
					*(*uint16)(unsafe.Add(mBase, uint32(v71))) = uint16(v78)
					v82 = v60 | v69
				}
				if v60 == l3+int32(-1) {
				} else {
					v89 = v82 << (uint(int32(1)) % 32)
					v98 = l2 + v89
					v99 = l0 + v89
					v100 = l3 - v82
					v101 = l1 + v89
					for {
						v102 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99))))
						v103 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101))))
						v105 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v98))))
						v106 = v102 - v103 + v105
						*(*uint16)(unsafe.Add(mBase, uint32(v98))) = uint16(v106)
						v108 = int32(2)
						v109 = v98 + v108
						v112 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99+v108))))
						v115 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101+v108))))
						v117 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v109))))
						v118 = v112 - v115 + v117
						*(*uint16)(unsafe.Add(mBase, uint32(v109))) = uint16(v118)
						v120 = int32(4)
						v127 = v100 + int32(-2)
						if v127 != 0 {
							v98 = v98 + v120
							v99 = v99 + v120
							v100 = v127
							v101 = v101 + v120
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.B2i32(base.Ui32(l2) < base.Ui32(l1+v16))&base.B2i32(base.Ui32(l1) < base.Ui32(v19)) != 0 {
					v60 = v12
					if l3&int32(1) == int32(0) {
						v82 = v60
					} else {
						v69 = int32(1)
						v70 = v60 << (uint(v69) % 32)
						v71 = l2 + v70
						v73 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v70))))
						v75 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v70))))
						v77 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v71))))
						v78 = v73 - v75 + v77
						*(*uint16)(unsafe.Add(mBase, uint32(v71))) = uint16(v78)
						v82 = v60 | v69
					}
					if v60 == l3+int32(-1) {
					} else {
						v89 = v82 << (uint(int32(1)) % 32)
						v98 = l2 + v89
						v99 = l0 + v89
						v100 = l3 - v82
						v101 = l1 + v89
						for {
							v102 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99))))
							v103 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101))))
							v105 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v98))))
							v106 = v102 - v103 + v105
							*(*uint16)(unsafe.Add(mBase, uint32(v98))) = uint16(v106)
							v108 = int32(2)
							v109 = v98 + v108
							v112 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99+v108))))
							v115 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101+v108))))
							v117 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v109))))
							v118 = v112 - v115 + v117
							*(*uint16)(unsafe.Add(mBase, uint32(v109))) = uint16(v118)
							v120 = int32(4)
							v127 = v100 + int32(-2)
							if v127 != 0 {
								v98 = v98 + v120
								v99 = v99 + v120
								v100 = v127
								v101 = v101 + v120
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v27 = l3 & int32(2147483640)
					v33 = l2
					v34 = l1
					v35 = v27
					v36 = l0
					for {
						v37 = int32(0)
						v38 = base.Simd_g_v128_load(m, v36, v37)
						v40 = base.Simd_g_v128_load(m, v34, v37)
						v43 = base.Simd_g_v128_load(m, v33, v37)
						v44 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v38, v40), v43)
						base.Simd_g_v128_store(m, v33, v37, v44)
						v47 = int32(16)
						v54 = v35 + int32(-8)
						if v54 != 0 {
							v33 = v33 + v47
							v34 = v34 + v47
							v35 = v54
							v36 = v36 + v47
							continue
						} else {
							break
						}
						break
					}
					if v27 == l3 {
					} else {
						v60 = v27
						if l3&int32(1) == int32(0) {
							v82 = v60
						} else {
							v69 = int32(1)
							v70 = v60 << (uint(v69) % 32)
							v71 = l2 + v70
							v73 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v70))))
							v75 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v70))))
							v77 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v71))))
							v78 = v73 - v75 + v77
							*(*uint16)(unsafe.Add(mBase, uint32(v71))) = uint16(v78)
							v82 = v60 | v69
						}
						if v60 == l3+int32(-1) {
						} else {
							v89 = v82 << (uint(int32(1)) % 32)
							v98 = l2 + v89
							v99 = l0 + v89
							v100 = l3 - v82
							v101 = l1 + v89
							for {
								v102 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99))))
								v103 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101))))
								v105 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v98))))
								v106 = v102 - v103 + v105
								*(*uint16)(unsafe.Add(mBase, uint32(v98))) = uint16(v106)
								v108 = int32(2)
								v109 = v98 + v108
								v112 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v99+v108))))
								v115 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v101+v108))))
								v117 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v109))))
								v118 = v112 - v115 + v117
								*(*uint16)(unsafe.Add(mBase, uint32(v109))) = uint16(v118)
								v120 = int32(4)
								v127 = v100 + int32(-2)
								if v127 != 0 {
									v98 = v98 + v120
									v99 = v99 + v120
									v100 = v127
									v101 = v101 + v120
									continue
								} else {
									break
								}
								break
							}
						}
					}
				}
			}
		}
	}
	return
}
func F_SharpYuvUpdateRGB_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v54 base.V128
	_ = v54
	var v58 base.V128
	_ = v58
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v65 int32
	_ = v65
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v95 base.V128
	_ = v95
	var v98 base.V128
	_ = v98
	var v101 base.V128
	_ = v101
	var v102 base.V128
	_ = v102
	var v110 int32
	_ = v110
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v128 int32
	_ = v128
	var v137 int32
	_ = v137
	var v142 int32
	_ = v142
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v155 int32
	_ = v155
	var v156 int32
	_ = v156
	var v158 int32
	_ = v158
	var v159 base.V128
	_ = v159
	var v161 base.V128
	_ = v161
	var v164 base.V128
	_ = v164
	var v165 base.V128
	_ = v165
	var v168 int32
	_ = v168
	var v175 int32
	_ = v175
	var v183 int32
	_ = v183
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v208 int32
	_ = v208
	var v214 int32
	_ = v214
	var v222 int32
	_ = v222
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v247 int32
	_ = v247
	var v254 int32
	_ = v254
	v12 = int32(8)
	if l3 < v12 {
		v110 = int32(0)
	} else {
		v17 = l3 + int32(-8)
		if base.Ui32(int32(8)) <= base.Ui32(v17) {
			v32 = int32(8)
			v33 = (int32(base.Ui32(v17)>>(uint(int32(3))%32)) + int32(1)) & int32(1073741822)
			v35 = l0
			v36 = l1
			v37 = l2
			for {
				v39 = int32(0)
				v40 = base.Simd_g_v128_load(m, v35, v39)
				v42 = base.Simd_g_v128_load(m, v36, v39)
				v45 = base.Simd_g_v128_load(m, v37, v39)
				v46 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v40, v42), v45)
				base.Simd_g_v128_store(m, v37, v39, v46)
				v49 = int32(16)
				v50 = v37 + v49
				v54 = base.Simd_g_v128_load(m, v35+v49, v39)
				v58 = base.Simd_g_v128_load(m, v36+v49, v39)
				v61 = base.Simd_g_v128_load(m, v50, v39)
				v62 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v54, v58), v61)
				base.Simd_g_v128_store(m, v50, v39, v62)
				v65 = int32(32)
				v72 = v32 + v49
				v74 = v33 + int32(-2)
				if v74 != 0 {
					v32 = v72
					v33 = v74
					v35 = v35 + v65
					v36 = v36 + v65
					v37 = v37 + v65
					continue
				} else {
					break
				}
				break
			}
			v81 = v72
			v82 = v32 + int32(8)
		} else {
			v81 = v12
			v82 = int32(0)
		}
		if v17&int32(8) != 0 {
			v110 = v82
		} else {
			v91 = v82 << (uint(int32(1)) % 32)
			v92 = l2 + v91
			v94 = int32(0)
			v95 = base.Simd_g_v128_load(m, l0+v91, v94)
			v98 = base.Simd_g_v128_load(m, l1+v91, v94)
			v101 = base.Simd_g_v128_load(m, v92, v94)
			v102 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v95, v98), v101)
			base.Simd_g_v128_store(m, v92, v94, v102)
			v110 = v81
		}
	}
	if l3 <= v110 {
	} else {
		v117 = l3 - v110
		if base.Ui32(v117) <= base.Ui32(int32(15)) {
			v183 = v110
			if (l3-v183)&int32(1) == int32(0) {
				v208 = v183
			} else {
				v194 = int32(1)
				v195 = v183 << (uint(v194) % 32)
				v196 = l2 + v195
				v198 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v195))))
				v200 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v195))))
				v202 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v196))))
				v203 = v198 - v200 + v202
				*(*uint16)(unsafe.Add(mBase, uint32(v196))) = uint16(v203)
				v208 = v183 + v194
			}
			if v183 == l3+int32(-1) {
			} else {
				v214 = v208 << (uint(int32(1)) % 32)
				v222 = l3 - v208
				v225 = l0 + v214
				v226 = l1 + v214
				v227 = l2 + v214
				for {
					v229 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
					v230 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226))))
					v232 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v227))))
					v233 = v229 - v230 + v232
					*(*uint16)(unsafe.Add(mBase, uint32(v227))) = uint16(v233)
					v235 = int32(2)
					v236 = v227 + v235
					v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225+v235))))
					v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226+v235))))
					v244 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v236))))
					v245 = v239 - v242 + v244
					*(*uint16)(unsafe.Add(mBase, uint32(v236))) = uint16(v245)
					v247 = int32(4)
					v254 = v222 + int32(-2)
					if v254 != 0 {
						v222 = v254
						v225 = v225 + v247
						v226 = v226 + v247
						v227 = v227 + v247
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v120 = int32(1)
			v121 = v110 << (uint(v120) % 32)
			v122 = l2 + v121
			v124 = l3 << (uint(v120) % 32)
			v128 = l2 + v124
			if base.B2i32(base.Ui32(v122) < base.Ui32(l0+v124))&base.B2i32(base.Ui32(l0+v121) < base.Ui32(v128)) != 0 {
				v183 = v110
				if (l3-v183)&int32(1) == int32(0) {
					v208 = v183
				} else {
					v194 = int32(1)
					v195 = v183 << (uint(v194) % 32)
					v196 = l2 + v195
					v198 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v195))))
					v200 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v195))))
					v202 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v196))))
					v203 = v198 - v200 + v202
					*(*uint16)(unsafe.Add(mBase, uint32(v196))) = uint16(v203)
					v208 = v183 + v194
				}
				if v183 == l3+int32(-1) {
				} else {
					v214 = v208 << (uint(int32(1)) % 32)
					v222 = l3 - v208
					v225 = l0 + v214
					v226 = l1 + v214
					v227 = l2 + v214
					for {
						v229 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
						v230 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226))))
						v232 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v227))))
						v233 = v229 - v230 + v232
						*(*uint16)(unsafe.Add(mBase, uint32(v227))) = uint16(v233)
						v235 = int32(2)
						v236 = v227 + v235
						v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225+v235))))
						v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226+v235))))
						v244 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v236))))
						v245 = v239 - v242 + v244
						*(*uint16)(unsafe.Add(mBase, uint32(v236))) = uint16(v245)
						v247 = int32(4)
						v254 = v222 + int32(-2)
						if v254 != 0 {
							v222 = v254
							v225 = v225 + v247
							v226 = v226 + v247
							v227 = v227 + v247
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.B2i32(base.Ui32(v122) < base.Ui32(l1+v124))&base.B2i32(base.Ui32(l1+v121) < base.Ui32(v128)) != 0 {
					v183 = v110
					if (l3-v183)&int32(1) == int32(0) {
						v208 = v183
					} else {
						v194 = int32(1)
						v195 = v183 << (uint(v194) % 32)
						v196 = l2 + v195
						v198 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v195))))
						v200 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v195))))
						v202 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v196))))
						v203 = v198 - v200 + v202
						*(*uint16)(unsafe.Add(mBase, uint32(v196))) = uint16(v203)
						v208 = v183 + v194
					}
					if v183 == l3+int32(-1) {
					} else {
						v214 = v208 << (uint(int32(1)) % 32)
						v222 = l3 - v208
						v225 = l0 + v214
						v226 = l1 + v214
						v227 = l2 + v214
						for {
							v229 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
							v230 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226))))
							v232 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v227))))
							v233 = v229 - v230 + v232
							*(*uint16)(unsafe.Add(mBase, uint32(v227))) = uint16(v233)
							v235 = int32(2)
							v236 = v227 + v235
							v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225+v235))))
							v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226+v235))))
							v244 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v236))))
							v245 = v239 - v242 + v244
							*(*uint16)(unsafe.Add(mBase, uint32(v236))) = uint16(v245)
							v247 = int32(4)
							v254 = v222 + int32(-2)
							if v254 != 0 {
								v222 = v254
								v225 = v225 + v247
								v226 = v226 + v247
								v227 = v227 + v247
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v137 = v110 << (uint(int32(1)) % 32)
					v142 = l3 & int32(7)
					v151 = v110 + v142 - l3
					v154 = l1 + v137
					v155 = l0 + v137
					v156 = l2 + v137
					for {
						v158 = int32(0)
						v159 = base.Simd_g_v128_load(m, v155, v158)
						v161 = base.Simd_g_v128_load(m, v154, v158)
						v164 = base.Simd_g_v128_load(m, v156, v158)
						v165 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v159, v161), v164)
						base.Simd_g_v128_store(m, v156, v158, v165)
						v168 = int32(16)
						v175 = v151 + int32(8)
						if v175 != 0 {
							v151 = v175
							v154 = v154 + v168
							v155 = v155 + v168
							v156 = v156 + v168
							continue
						} else {
							break
						}
						break
					}
					if v142 == int32(0) {
					} else {
						v183 = v110 + (v117 - v142)
						if (l3-v183)&int32(1) == int32(0) {
							v208 = v183
						} else {
							v194 = int32(1)
							v195 = v183 << (uint(v194) % 32)
							v196 = l2 + v195
							v198 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0+v195))))
							v200 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l1+v195))))
							v202 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v196))))
							v203 = v198 - v200 + v202
							*(*uint16)(unsafe.Add(mBase, uint32(v196))) = uint16(v203)
							v208 = v183 + v194
						}
						if v183 == l3+int32(-1) {
						} else {
							v214 = v208 << (uint(int32(1)) % 32)
							v222 = l3 - v208
							v225 = l0 + v214
							v226 = l1 + v214
							v227 = l2 + v214
							for {
								v229 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
								v230 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226))))
								v232 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v227))))
								v233 = v229 - v230 + v232
								*(*uint16)(unsafe.Add(mBase, uint32(v227))) = uint16(v233)
								v235 = int32(2)
								v236 = v227 + v235
								v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225+v235))))
								v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v226+v235))))
								v244 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v236))))
								v245 = v239 - v242 + v244
								*(*uint16)(unsafe.Add(mBase, uint32(v236))) = uint16(v245)
								v247 = int32(4)
								v254 = v222 + int32(-2)
								if v254 != 0 {
									v222 = v254
									v225 = v225 + v247
									v226 = v226 + v247
									v227 = v227 + v247
									continue
								} else {
									break
								}
								break
							}
						}
					}
				}
			}
		}
	}
	return
}
func F_SharpYuvUpdateY_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int64 {
	mBase := m.M
	_ = mBase
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v25 int64
	_ = v25
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v47 int32
	_ = v47
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v58 base.V128
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v64 base.V128
	_ = v64
	var v65 int32
	_ = v65
	var v66 base.V128
	_ = v66
	var v69 base.V128
	_ = v69
	var v71 base.V128
	_ = v71
	var v73 base.V128
	_ = v73
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v87 int32
	_ = v87
	var v95 base.V128
	_ = v95
	var v97 int32
	_ = v97
	var v102 int64
	_ = v102
	var v110 int64
	_ = v110
	var v111 int32
	_ = v111
	var v122 int32
	_ = v122
	var v130 int32
	_ = v130
	var v132 int64
	_ = v132
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v141 int32
	_ = v141
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v159 int64
	_ = v159
	var v160 int32
	_ = v160
	var v167 int32
	_ = v167
	var v174 int64
	_ = v174
	if int32(1) <= l3 {
		v21 = int32(-1)
		v24 = v21<<(uint(l4)%32) ^ v21
		v25 = int64(0)
		if base.Ui32(int32(6)) <= base.Ui32(l3) {
			v30 = l3 << (uint(int32(1)) % 32)
			v33 = l2 + v30
			if base.B2i32(base.Ui32(l2) < base.Ui32(l0+v30))&base.B2i32(base.Ui32(l0) < base.Ui32(v33)) == int32(0) {
				if base.B2i32(base.Ui32(l2) < base.Ui32(l1+v30))&base.B2i32(base.Ui32(l1) < base.Ui32(v33)) != 0 {
					v110 = v25
					v111 = int32(0)
					v122 = v111 << (uint(int32(1)) % 32)
					v130 = l2 + v122
					v132 = v110
					v134 = l0 + v122
					v138 = l1 + v122
					v141 = l3 - v111
					for {
						v143 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v134))))
						v144 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v138))))
						v145 = v143 - v144
						v146 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v130))))
						v147 = v145 + v146
						if v147 < v24 {
							v149 = v147
						} else {
							v149 = v24
						}
						if v147 < int32(0) {
							v152 = int32(0)
						} else {
							v152 = v149
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v130))) = uint16(v152)
						v155 = v145 >> (uint(int32(31)) % 32)
						v159 = v132 + base.I64_extend_i32_u(v145^v155-v155)
						v160 = int32(2)
						v167 = v141 + int32(-1)
						if v167 != 0 {
							v130 = v130 + v160
							v132 = v159
							v134 = v134 + v160
							v138 = v138 + v160
							v141 = v167
							continue
						} else {
							break
						}
						break
					}
					v174 = v159
				} else {
					v47 = l3 & int32(2147483646)
					v52 = l2
					v56 = l1
					v58 = base.Simd_g_const(&F_SharpYuvUpdateY_C__k0)
					v59 = v47
					v60 = l0
					for {
						v64 = base.Simd_g_const(&F_SharpYuvUpdateY_C__k0)
						v65 = int32(0)
						v66 = base.Simd_g_v128_load32_zero(m, v60, v65)
						v69 = base.Simd_g_v128_load32_zero(m, v56, v65)
						v71 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(v66), base.Simd_g_i32x4_extend_low_i16x8_u(v69))
						v73 = base.Simd_g_v128_load32_zero(m, v52, v65)
						v75 = base.Simd_g_i32x4_add(v71, base.Simd_g_i32x4_extend_low_i16x8_u(v73))
						v77 = base.Simd_g_const(&F_SharpYuvUpdateY_C__k1)
						base.Simd_g_v128_store32_lane_l0(m, v52, v65, base.Simd_g_v128_bitselect(v64, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(v75, base.Simd_g_i32x4_splat(v24)), v58, base.Simd_g_const(&F_SharpYuvUpdateY_C__k1), base.Simd_g_const(&F_SharpYuvUpdateY_C__k2)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v75, v64), v58, base.Simd_g_const(&F_SharpYuvUpdateY_C__k1), base.Simd_g_const(&F_SharpYuvUpdateY_C__k2))))
						v87 = int32(4)
						v95 = base.Simd_g_i64x2_add(v58, base.Simd_g_i64x2_extend_low_i32x4_u(base.Simd_g_i32x4_abs(v71)))
						v97 = v59 + int32(-2)
						if v97 != 0 {
							v52 = v52 + v87
							v56 = v56 + v87
							v58 = v95
							v59 = v97
							v60 = v60 + v87
							continue
						} else {
							break
						}
						break
					}
					v102 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_i64x2_add(v95, base.Simd_g_i8x16_shuffle2(v95, v95, base.Simd_g_const(&F_SharpYuvUpdateY_C__k3), base.Simd_g_const(&F_SharpYuvUpdateY_C__k2))))
					if v47 == l3 {
						v174 = v102
					} else {
						v110 = v102
						v111 = v47
						v122 = v111 << (uint(int32(1)) % 32)
						v130 = l2 + v122
						v132 = v110
						v134 = l0 + v122
						v138 = l1 + v122
						v141 = l3 - v111
						for {
							v143 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v134))))
							v144 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v138))))
							v145 = v143 - v144
							v146 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v130))))
							v147 = v145 + v146
							if v147 < v24 {
								v149 = v147
							} else {
								v149 = v24
							}
							if v147 < int32(0) {
								v152 = int32(0)
							} else {
								v152 = v149
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v130))) = uint16(v152)
							v155 = v145 >> (uint(int32(31)) % 32)
							v159 = v132 + base.I64_extend_i32_u(v145^v155-v155)
							v160 = int32(2)
							v167 = v141 + int32(-1)
							if v167 != 0 {
								v130 = v130 + v160
								v132 = v159
								v134 = v134 + v160
								v138 = v138 + v160
								v141 = v167
								continue
							} else {
								break
							}
							break
						}
						v174 = v159
					}
				}
			} else {
				v110 = v25
				v111 = int32(0)
				v122 = v111 << (uint(int32(1)) % 32)
				v130 = l2 + v122
				v132 = v110
				v134 = l0 + v122
				v138 = l1 + v122
				v141 = l3 - v111
				for {
					v143 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v134))))
					v144 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v138))))
					v145 = v143 - v144
					v146 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v130))))
					v147 = v145 + v146
					if v147 < v24 {
						v149 = v147
					} else {
						v149 = v24
					}
					if v147 < int32(0) {
						v152 = int32(0)
					} else {
						v152 = v149
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v130))) = uint16(v152)
					v155 = v145 >> (uint(int32(31)) % 32)
					v159 = v132 + base.I64_extend_i32_u(v145^v155-v155)
					v160 = int32(2)
					v167 = v141 + int32(-1)
					if v167 != 0 {
						v130 = v130 + v160
						v132 = v159
						v134 = v134 + v160
						v138 = v138 + v160
						v141 = v167
						continue
					} else {
						break
					}
					break
				}
				v174 = v159
			}
		} else {
			v110 = v25
			v111 = int32(0)
			v122 = v111 << (uint(int32(1)) % 32)
			v130 = l2 + v122
			v132 = v110
			v134 = l0 + v122
			v138 = l1 + v122
			v141 = l3 - v111
			for {
				v143 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v134))))
				v144 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v138))))
				v145 = v143 - v144
				v146 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v130))))
				v147 = v145 + v146
				if v147 < v24 {
					v149 = v147
				} else {
					v149 = v24
				}
				if v147 < int32(0) {
					v152 = int32(0)
				} else {
					v152 = v149
				}
				*(*uint16)(unsafe.Add(mBase, uint32(v130))) = uint16(v152)
				v155 = v145 >> (uint(int32(31)) % 32)
				v159 = v132 + base.I64_extend_i32_u(v145^v155-v155)
				v160 = int32(2)
				v167 = v141 + int32(-1)
				if v167 != 0 {
					v130 = v130 + v160
					v132 = v159
					v134 = v134 + v160
					v138 = v138 + v160
					v141 = v167
					continue
				} else {
					break
				}
				break
			}
			v174 = v159
		}
		return v174
	} else {
		return int64(0)
	}
}

var F_SharpYuvUpdateY_C__k0 = [2]uint64{0x0, 0x0}
var F_SharpYuvUpdateY_C__k1 = [2]uint64{0x100010005040100, 0x100010001000100}
var F_SharpYuvUpdateY_C__k2 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SharpYuvUpdateY_C__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}

func F_SharpYuvUpdateY_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int64 {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v33 int32
	_ = v33
	var v36 base.V128
	_ = v36
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v46 int32
	_ = v46
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v56 base.V128
	_ = v56
	var v64 base.V128
	_ = v64
	var v65 int32
	_ = v65
	var v74 int32
	_ = v74
	var v82 int32
	_ = v82
	var v83 base.V128
	_ = v83
	var v104 int64
	_ = v104
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v134 int32
	_ = v134
	var v144 int32
	_ = v144
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v152 base.V128
	_ = v152
	var v157 base.V128
	_ = v157
	var v158 int32
	_ = v158
	var v159 base.V128
	_ = v159
	var v162 base.V128
	_ = v162
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v170 base.V128
	_ = v170
	var v180 int32
	_ = v180
	var v188 base.V128
	_ = v188
	var v190 int32
	_ = v190
	var v194 int32
	_ = v194
	var v195 int64
	_ = v195
	var v204 int32
	_ = v204
	var v212 int64
	_ = v212
	var v217 int32
	_ = v217
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v235 int64
	_ = v235
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v245 int32
	_ = v245
	var v248 int32
	_ = v248
	var v251 int32
	_ = v251
	var v255 int64
	_ = v255
	var v256 int32
	_ = v256
	var v263 int32
	_ = v263
	var v278 int64
	_ = v278
	v18 = int32(-1)
	v21 = v18<<(uint(l4)%32) ^ v18
	if int32(8) <= l3 {
		v33 = l2
		v36 = base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k0)
		v38 = int32(0)
		v39 = l1
		v40 = l0
		for {
			v46 = int32(0)
			v47 = base.Simd_g_v128_load(m, v33, v46)
			v49 = base.Simd_g_v128_load(m, v40, v46)
			v51 = base.Simd_g_v128_load(m, v39, v46)
			v52 = base.Simd_g_i16x8_sub(v49, v51)
			v56 = base.Simd_g_i16x8_max_s(base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_add(v47, v52), base.Simd_g_i16x8_splat(v21)), base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k0))
			base.Simd_g_v128_store(m, v33, v46, v56)
			v64 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v52, base.Simd_g_v128_or(base.Simd_g_i16x8_shr_s(v52, int32(15)), base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k1))), v36)
			v65 = int32(16)
			v74 = v38 + int32(8)
			if v38+v65 <= l3 {
				v33 = v33 + v65
				v36 = v64
				v38 = v74
				v39 = v39 + v65
				v40 = v40 + v65
				continue
			} else {
				break
			}
			break
		}
		v82 = v74
		v83 = v64
	} else {
		v82 = int32(0)
		v83 = base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k0)
	}
	v104 = base.I64_extend_i32_u(base.Simd_g_i32x4_extract_lane_l3(v83) + base.Simd_g_i32x4_extract_lane_l2(v83) + base.Simd_g_i32x4_extract_lane_l1(v83) + base.Simd_g_i32x4_extract_lane_l0(v83))
	if l3 <= v82 {
		v278 = v104
	} else {
		v106 = l3 - v82
		if base.Ui32(v106) < base.Ui32(int32(8)) {
			v204 = v82
			v212 = v104
			v217 = v204 << (uint(int32(1)) % 32)
			v224 = l3 - v204
			v225 = l2 + v217
			v230 = l0 + v217
			v231 = l1 + v217
			v235 = v212
			for {
				v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v230))))
				v240 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v231))))
				v241 = v239 - v240
				v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
				v243 = v241 + v242
				if v243 < v21 {
					v245 = v243
				} else {
					v245 = v21
				}
				if v243 < int32(0) {
					v248 = int32(0)
				} else {
					v248 = v245
				}
				*(*uint16)(unsafe.Add(mBase, uint32(v225))) = uint16(v248)
				v251 = v241 >> (uint(int32(31)) % 32)
				v255 = v235 + base.I64_extend_i32_u(v241^v251-v251)
				v256 = int32(2)
				v263 = v224 + int32(-1)
				if v263 != 0 {
					v224 = v263
					v225 = v225 + v256
					v230 = v230 + v256
					v231 = v231 + v256
					v235 = v255
					continue
				} else {
					break
				}
				break
			}
			v278 = v255
		} else {
			v109 = int32(1)
			v110 = v82 << (uint(v109) % 32)
			v111 = l2 + v110
			v113 = l3 << (uint(v109) % 32)
			v117 = l2 + v113
			if base.B2i32(base.Ui32(v111) < base.Ui32(l0+v113))&base.B2i32(base.Ui32(l0+v110) < base.Ui32(v117)) != 0 {
				v204 = v82
				v212 = v104
				v217 = v204 << (uint(int32(1)) % 32)
				v224 = l3 - v204
				v225 = l2 + v217
				v230 = l0 + v217
				v231 = l1 + v217
				v235 = v212
				for {
					v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v230))))
					v240 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v231))))
					v241 = v239 - v240
					v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
					v243 = v241 + v242
					if v243 < v21 {
						v245 = v243
					} else {
						v245 = v21
					}
					if v243 < int32(0) {
						v248 = int32(0)
					} else {
						v248 = v245
					}
					*(*uint16)(unsafe.Add(mBase, uint32(v225))) = uint16(v248)
					v251 = v241 >> (uint(int32(31)) % 32)
					v255 = v235 + base.I64_extend_i32_u(v241^v251-v251)
					v256 = int32(2)
					v263 = v224 + int32(-1)
					if v263 != 0 {
						v224 = v263
						v225 = v225 + v256
						v230 = v230 + v256
						v231 = v231 + v256
						v235 = v255
						continue
					} else {
						break
					}
					break
				}
				v278 = v255
			} else {
				if base.B2i32(base.Ui32(v111) < base.Ui32(l1+v113))&base.B2i32(base.Ui32(l1+v110) < base.Ui32(v117)) != 0 {
					v204 = v82
					v212 = v104
					v217 = v204 << (uint(int32(1)) % 32)
					v224 = l3 - v204
					v225 = l2 + v217
					v230 = l0 + v217
					v231 = l1 + v217
					v235 = v212
					for {
						v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v230))))
						v240 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v231))))
						v241 = v239 - v240
						v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
						v243 = v241 + v242
						if v243 < v21 {
							v245 = v243
						} else {
							v245 = v21
						}
						if v243 < int32(0) {
							v248 = int32(0)
						} else {
							v248 = v245
						}
						*(*uint16)(unsafe.Add(mBase, uint32(v225))) = uint16(v248)
						v251 = v241 >> (uint(int32(31)) % 32)
						v255 = v235 + base.I64_extend_i32_u(v241^v251-v251)
						v256 = int32(2)
						v263 = v224 + int32(-1)
						if v263 != 0 {
							v224 = v263
							v225 = v225 + v256
							v230 = v230 + v256
							v231 = v231 + v256
							v235 = v255
							continue
						} else {
							break
						}
						break
					}
					v278 = v255
				} else {
					v128 = int32(1)
					v129 = v82 << (uint(v128) % 32)
					v134 = l3 & v128
					v144 = l2 + v129
					v149 = l1 + v129
					v150 = l0 + v129
					v151 = v82 + v134 - l3
					v152 = base.Simd_g_i64x2_replace_lane_l0(base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k0), v104)
					for {
						v157 = base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k0)
						v158 = int32(0)
						v159 = base.Simd_g_v128_load32_zero(m, v150, v158)
						v162 = base.Simd_g_v128_load32_zero(m, v149, v158)
						v164 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(v159), base.Simd_g_i32x4_extend_low_i16x8_u(v162))
						v166 = base.Simd_g_v128_load32_zero(m, v144, v158)
						v168 = base.Simd_g_i32x4_add(v164, base.Simd_g_i32x4_extend_low_i16x8_u(v166))
						v170 = base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k2)
						base.Simd_g_v128_store32_lane_l0(m, v144, v158, base.Simd_g_v128_bitselect(v157, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(v168, base.Simd_g_i32x4_splat(v21)), v152, base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k2), base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k3)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_s(v168, v157), v152, base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k2), base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k3))))
						v180 = int32(4)
						v188 = base.Simd_g_i64x2_add(v152, base.Simd_g_i64x2_extend_low_i32x4_u(base.Simd_g_i32x4_abs(v164)))
						v190 = v151 + int32(2)
						if v190 != 0 {
							v144 = v144 + v180
							v149 = v149 + v180
							v150 = v150 + v180
							v151 = v190
							v152 = v188
							continue
						} else {
							break
						}
						break
					}
					v194 = int32(0)
					v195 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_i64x2_add(v188, base.Simd_g_i8x16_shuffle2(v188, v188, base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k4), base.Simd_g_const(&F_SharpYuvUpdateY_SSE2__k3))))
					if v134 == v194 {
						v278 = v195
					} else {
						v204 = v82 + (v106 - v134)
						v212 = v195
						v217 = v204 << (uint(int32(1)) % 32)
						v224 = l3 - v204
						v225 = l2 + v217
						v230 = l0 + v217
						v231 = l1 + v217
						v235 = v212
						for {
							v239 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v230))))
							v240 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v231))))
							v241 = v239 - v240
							v242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v225))))
							v243 = v241 + v242
							if v243 < v21 {
								v245 = v243
							} else {
								v245 = v21
							}
							if v243 < int32(0) {
								v248 = int32(0)
							} else {
								v248 = v245
							}
							*(*uint16)(unsafe.Add(mBase, uint32(v225))) = uint16(v248)
							v251 = v241 >> (uint(int32(31)) % 32)
							v255 = v235 + base.I64_extend_i32_u(v241^v251-v251)
							v256 = int32(2)
							v263 = v224 + int32(-1)
							if v263 != 0 {
								v224 = v263
								v225 = v225 + v256
								v230 = v230 + v256
								v231 = v231 + v256
								v235 = v255
								continue
							} else {
								break
							}
							break
						}
						v278 = v255
					}
				}
			}
		}
	}
	return v278
}

var F_SharpYuvUpdateY_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SharpYuvUpdateY_SSE2__k1 = [2]uint64{0x1000100010001, 0x1000100010001}
var F_SharpYuvUpdateY_SSE2__k2 = [2]uint64{0x100010005040100, 0x100010001000100}
var F_SharpYuvUpdateY_SSE2__k3 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SharpYuvUpdateY_SSE2__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}

func F_SimpleHFilter16_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v4 int32
	_ = v4
	var v8 base.V128
	_ = v8
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 base.V128
	_ = v34
	var v37 base.V128
	_ = v37
	var v40 base.V128
	_ = v40
	var v42 int32
	_ = v42
	var v43 base.V128
	_ = v43
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v55 base.V128
	_ = v55
	var v58 base.V128
	_ = v58
	var v61 base.V128
	_ = v61
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v82 base.V128
	_ = v82
	var v85 base.V128
	_ = v85
	var v88 base.V128
	_ = v88
	var v91 base.V128
	_ = v91
	var v95 int32
	_ = v95
	var v97 base.V128
	_ = v97
	var v100 base.V128
	_ = v100
	var v103 base.V128
	_ = v103
	var v106 base.V128
	_ = v106
	var v108 base.V128
	_ = v108
	var v110 base.V128
	_ = v110
	var v112 base.V128
	_ = v112
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v117 base.V128
	_ = v117
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v120 base.V128
	_ = v120
	var v121 base.V128
	_ = v121
	var v122 base.V128
	_ = v122
	var v124 base.V128
	_ = v124
	var v125 base.V128
	_ = v125
	var v127 base.V128
	_ = v127
	var v129 base.V128
	_ = v129
	var v130 base.V128
	_ = v130
	var v131 base.V128
	_ = v131
	var v134 base.V128
	_ = v134
	var v142 base.V128
	_ = v142
	var v155 base.V128
	_ = v155
	var v157 base.V128
	_ = v157
	var v160 int32
	_ = v160
	var v168 base.V128
	_ = v168
	var v170 base.V128
	_ = v170
	var v172 base.V128
	_ = v172
	var v183 base.V128
	_ = v183
	var v185 base.V128
	_ = v185
	var v187 base.V128
	_ = v187
	var v191 base.V128
	_ = v191
	var v192 base.V128
	_ = v192
	var v196 int32
	_ = v196
	var v198 base.V128
	_ = v198
	var v209 base.V128
	_ = v209
	var v213 int32
	_ = v213
	var v215 base.V128
	_ = v215
	var v219 int32
	_ = v219
	var v221 base.V128
	_ = v221
	var v232 base.V128
	_ = v232
	var v234 base.V128
	_ = v234
	var v236 base.V128
	_ = v236
	var v241 base.V128
	_ = v241
	var v245 int32
	_ = v245
	var v247 base.V128
	_ = v247
	var v258 base.V128
	_ = v258
	var v262 int32
	_ = v262
	var v264 base.V128
	_ = v264
	var v268 int32
	_ = v268
	var v270 base.V128
	_ = v270
	v4 = int32(0)
	v8 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k0)
	v23 = l0 + int32(-2)
	v25 = l1 * int32(6)
	v27 = int32(1)
	v28 = l1 << (uint(v27) % 32)
	v30 = int32(2)
	v31 = l1 << (uint(v30) % 32)
	v32 = v23 + v31
	v34 = base.Simd_g_v128_load32_splat(m, v23, v4)
	v37 = base.Simd_g_v128_load32_lane_l1(m, v32, v4, v34)
	v40 = base.Simd_g_v128_load32_lane_l2(m, v23+v28, v4, v37)
	v42 = int32(3)
	v43 = base.Simd_g_v128_load32_lane_l3(m, v23+v25, v4, v40)
	v45 = l1 * int32(7)
	v48 = l1 * v42
	v51 = l1 * int32(5)
	v53 = v23 + l1
	v55 = base.Simd_g_v128_load32_splat(m, v53, v4)
	v58 = base.Simd_g_v128_load32_lane_l1(m, v23+v51, v4, v55)
	v61 = base.Simd_g_v128_load32_lane_l2(m, v23+v48, v4, v58)
	v64 = base.Simd_g_v128_load32_lane_l3(m, v23+v45, v4, v61)
	v65 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k1)
	v66 = base.Simd_g_i8x16_shuffle2(v43, v64, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k3))
	v67 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k4)
	v68 = base.Simd_g_i8x16_shuffle2(v43, v64, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k6))
	v69 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k7)
	v70 = base.Simd_g_i8x16_shuffle2(v66, v68, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k9))
	v71 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k10)
	v72 = base.Simd_g_i8x16_shuffle2(v66, v68, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k12))
	v73 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k13)
	v74 = base.Simd_g_i8x16_shuffle2(v70, v72, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k15))
	v77 = v23 + l1<<(uint(v42)%32)
	v80 = v77 + v31
	v82 = base.Simd_g_v128_load32_splat(m, v77, v4)
	v85 = base.Simd_g_v128_load32_lane_l1(m, v80, v4, v82)
	v88 = base.Simd_g_v128_load32_lane_l2(m, v77+v28, v4, v85)
	v91 = base.Simd_g_v128_load32_lane_l3(m, v77+v25, v4, v88)
	v95 = v77 + l1
	v97 = base.Simd_g_v128_load32_splat(m, v95, v4)
	v100 = base.Simd_g_v128_load32_lane_l1(m, v77+v51, v4, v97)
	v103 = base.Simd_g_v128_load32_lane_l2(m, v77+v48, v4, v100)
	v106 = base.Simd_g_v128_load32_lane_l3(m, v77+v45, v4, v103)
	v108 = base.Simd_g_i8x16_shuffle2(v91, v106, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k3))
	v110 = base.Simd_g_i8x16_shuffle2(v91, v106, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k6))
	v112 = base.Simd_g_i8x16_shuffle2(v108, v110, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k9))
	v114 = base.Simd_g_i8x16_shuffle2(v108, v110, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k12))
	v116 = base.Simd_g_i8x16_shuffle2(v112, v114, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k15))
	v117 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k16)
	v118 = base.Simd_g_i8x16_shuffle2(v74, v116, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k18))
	v119 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k19)
	v120 = base.Simd_g_i8x16_shuffle2(v74, v116, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k21))
	v121 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k22)
	v122 = base.Simd_g_v128_xor(v120, v121)
	v124 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k23)
	v125 = base.Simd_g_i8x16_shuffle2(v70, v72, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k25))
	v127 = base.Simd_g_i8x16_shuffle2(v112, v114, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k25))
	v129 = base.Simd_g_i8x16_shuffle2(v125, v127, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k18))
	v130 = base.Simd_g_v128_xor(v129, v121)
	v131 = base.Simd_g_i8x16_sub_sat_s(v130, v122)
	v134 = base.Simd_g_i8x16_shuffle2(v125, v127, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k21))
	v142 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v120, v129), base.Simd_g_i8x16_sub_sat_u(v129, v120))
	v155 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v131, base.Simd_g_i8x16_add_sat_s(v131, base.Simd_g_i8x16_add_sat_s(v131, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v118, v121), base.Simd_g_v128_xor(v134, v121))))), v8, base.Simd_g_i8x16_eq(v8, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v142, v142), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v118, v134), base.Simd_g_i8x16_sub_sat_u(v134, v118)), v27), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k26))), base.Simd_g_i8x16_splat(l2))))
	v157 = base.Simd_g_i8x16_add_sat_s(v155, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k27))
	v160 = int32(11)
	v168 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v122, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v8, v157, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k3)), v160), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v8, v157, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k6)), v160))), v121)
	v170 = base.Simd_g_i8x16_shuffle2(v118, v168, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k3))
	v172 = base.Simd_g_i8x16_add_sat_s(v155, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k28))
	v183 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v130, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v8, v172, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k3)), v160), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v8, v172, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k6)), v160))), v121)
	v185 = base.Simd_g_i8x16_shuffle2(v183, v134, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k3))
	v187 = base.Simd_g_i8x16_shuffle2(v170, v185, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v23, v4, v187)
	v191 = base.Simd_g_const(&F_SimpleHFilter16_SSE2__k29)
	v192 = base.Simd_g_i8x16_shuffle2(v187, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v53, v4, v192)
	v196 = v53 + l1
	v198 = base.Simd_g_i8x16_shuffle2(v192, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v196, v4, v198)
	base.Simd_g_v128_store32_lane_l0(m, v196+l1, v4, base.Simd_g_i8x16_shuffle2(v198, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31)))
	v209 = base.Simd_g_i8x16_shuffle2(v170, v185, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v32, v4, v209)
	v213 = v32 + l1
	v215 = base.Simd_g_i8x16_shuffle2(v209, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v213, v4, v215)
	v219 = v213 + l1
	v221 = base.Simd_g_i8x16_shuffle2(v215, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v219, v4, v221)
	base.Simd_g_v128_store32_lane_l0(m, v219+l1, v4, base.Simd_g_i8x16_shuffle2(v221, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31)))
	v232 = base.Simd_g_i8x16_shuffle2(v118, v168, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k6))
	v234 = base.Simd_g_i8x16_shuffle2(v183, v134, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k6))
	v236 = base.Simd_g_i8x16_shuffle2(v232, v234, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v77, v4, v236)
	v241 = base.Simd_g_i8x16_shuffle2(v236, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v95, v4, v241)
	v245 = v95 + l1
	v247 = base.Simd_g_i8x16_shuffle2(v241, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v245, v4, v247)
	base.Simd_g_v128_store32_lane_l0(m, v245+l1, v4, base.Simd_g_i8x16_shuffle2(v247, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31)))
	v258 = base.Simd_g_i8x16_shuffle2(v232, v234, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v80, v4, v258)
	v262 = v80 + l1
	v264 = base.Simd_g_i8x16_shuffle2(v258, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v262, v4, v264)
	v268 = v262 + l1
	v270 = base.Simd_g_i8x16_shuffle2(v264, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v268, v4, v270)
	base.Simd_g_v128_store32_lane_l0(m, v268+l1, v4, base.Simd_g_i8x16_shuffle2(v270, v8, base.Simd_g_const(&F_SimpleHFilter16_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16_SSE2__k31)))
	return
}

var F_SimpleHFilter16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SimpleHFilter16_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SimpleHFilter16_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SimpleHFilter16_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_SimpleHFilter16_SSE2__k4 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_SimpleHFilter16_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_SimpleHFilter16_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_SimpleHFilter16_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_SimpleHFilter16_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_SimpleHFilter16_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_SimpleHFilter16_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_SimpleHFilter16_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_SimpleHFilter16_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_SimpleHFilter16_SSE2__k13 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_SimpleHFilter16_SSE2__k14 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_SimpleHFilter16_SSE2__k15 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_SimpleHFilter16_SSE2__k16 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_SimpleHFilter16_SSE2__k17 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_SimpleHFilter16_SSE2__k18 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_SimpleHFilter16_SSE2__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_SimpleHFilter16_SSE2__k20 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_SimpleHFilter16_SSE2__k21 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_SimpleHFilter16_SSE2__k22 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SimpleHFilter16_SSE2__k23 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_SimpleHFilter16_SSE2__k24 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_SimpleHFilter16_SSE2__k25 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_SimpleHFilter16_SSE2__k26 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_SimpleHFilter16_SSE2__k27 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_SimpleHFilter16_SSE2__k28 = [2]uint64{0x404040404040404, 0x404040404040404}
var F_SimpleHFilter16_SSE2__k29 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_SimpleHFilter16_SSE2__k30 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_SimpleHFilter16_SSE2__k31 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_SimpleHFilter16i_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v6 int32
	_ = v6
	var v10 base.V128
	_ = v10
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 base.V128
	_ = v36
	var v39 base.V128
	_ = v39
	var v42 base.V128
	_ = v42
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v57 base.V128
	_ = v57
	var v60 base.V128
	_ = v60
	var v63 base.V128
	_ = v63
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v84 base.V128
	_ = v84
	var v87 base.V128
	_ = v87
	var v90 base.V128
	_ = v90
	var v93 base.V128
	_ = v93
	var v97 int32
	_ = v97
	var v99 base.V128
	_ = v99
	var v102 base.V128
	_ = v102
	var v105 base.V128
	_ = v105
	var v108 base.V128
	_ = v108
	var v110 base.V128
	_ = v110
	var v112 base.V128
	_ = v112
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v120 base.V128
	_ = v120
	var v121 base.V128
	_ = v121
	var v122 base.V128
	_ = v122
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v127 base.V128
	_ = v127
	var v129 base.V128
	_ = v129
	var v131 base.V128
	_ = v131
	var v132 base.V128
	_ = v132
	var v133 base.V128
	_ = v133
	var v136 base.V128
	_ = v136
	var v144 base.V128
	_ = v144
	var v157 base.V128
	_ = v157
	var v159 base.V128
	_ = v159
	var v162 int32
	_ = v162
	var v170 base.V128
	_ = v170
	var v172 base.V128
	_ = v172
	var v174 base.V128
	_ = v174
	var v185 base.V128
	_ = v185
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v193 base.V128
	_ = v193
	var v194 base.V128
	_ = v194
	var v198 int32
	_ = v198
	var v200 base.V128
	_ = v200
	var v211 base.V128
	_ = v211
	var v215 int32
	_ = v215
	var v217 base.V128
	_ = v217
	var v221 int32
	_ = v221
	var v223 base.V128
	_ = v223
	var v234 base.V128
	_ = v234
	var v236 base.V128
	_ = v236
	var v238 base.V128
	_ = v238
	var v243 base.V128
	_ = v243
	var v247 int32
	_ = v247
	var v249 base.V128
	_ = v249
	var v260 base.V128
	_ = v260
	var v264 int32
	_ = v264
	var v266 base.V128
	_ = v266
	var v270 int32
	_ = v270
	var v272 base.V128
	_ = v272
	var v284 int32
	_ = v284
	var v288 base.V128
	_ = v288
	var v303 int32
	_ = v303
	var v305 int32
	_ = v305
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v311 int32
	_ = v311
	var v312 int32
	_ = v312
	var v314 base.V128
	_ = v314
	var v317 base.V128
	_ = v317
	var v320 base.V128
	_ = v320
	var v322 int32
	_ = v322
	var v323 base.V128
	_ = v323
	var v325 int32
	_ = v325
	var v328 int32
	_ = v328
	var v331 int32
	_ = v331
	var v333 int32
	_ = v333
	var v335 base.V128
	_ = v335
	var v338 base.V128
	_ = v338
	var v341 base.V128
	_ = v341
	var v344 base.V128
	_ = v344
	var v345 base.V128
	_ = v345
	var v346 base.V128
	_ = v346
	var v347 base.V128
	_ = v347
	var v348 base.V128
	_ = v348
	var v349 base.V128
	_ = v349
	var v350 base.V128
	_ = v350
	var v351 base.V128
	_ = v351
	var v352 base.V128
	_ = v352
	var v353 base.V128
	_ = v353
	var v354 base.V128
	_ = v354
	var v357 int32
	_ = v357
	var v360 int32
	_ = v360
	var v362 base.V128
	_ = v362
	var v365 base.V128
	_ = v365
	var v368 base.V128
	_ = v368
	var v371 base.V128
	_ = v371
	var v375 int32
	_ = v375
	var v377 base.V128
	_ = v377
	var v380 base.V128
	_ = v380
	var v383 base.V128
	_ = v383
	var v386 base.V128
	_ = v386
	var v388 base.V128
	_ = v388
	var v390 base.V128
	_ = v390
	var v392 base.V128
	_ = v392
	var v394 base.V128
	_ = v394
	var v396 base.V128
	_ = v396
	var v397 base.V128
	_ = v397
	var v398 base.V128
	_ = v398
	var v399 base.V128
	_ = v399
	var v400 base.V128
	_ = v400
	var v401 base.V128
	_ = v401
	var v402 base.V128
	_ = v402
	var v404 base.V128
	_ = v404
	var v405 base.V128
	_ = v405
	var v407 base.V128
	_ = v407
	var v409 base.V128
	_ = v409
	var v410 base.V128
	_ = v410
	var v411 base.V128
	_ = v411
	var v414 base.V128
	_ = v414
	var v422 base.V128
	_ = v422
	var v435 base.V128
	_ = v435
	var v437 base.V128
	_ = v437
	var v440 int32
	_ = v440
	var v448 base.V128
	_ = v448
	var v450 base.V128
	_ = v450
	var v452 base.V128
	_ = v452
	var v463 base.V128
	_ = v463
	var v465 base.V128
	_ = v465
	var v467 base.V128
	_ = v467
	var v471 base.V128
	_ = v471
	var v472 base.V128
	_ = v472
	var v476 int32
	_ = v476
	var v478 base.V128
	_ = v478
	var v489 base.V128
	_ = v489
	var v493 int32
	_ = v493
	var v495 base.V128
	_ = v495
	var v499 int32
	_ = v499
	var v501 base.V128
	_ = v501
	var v512 base.V128
	_ = v512
	var v514 base.V128
	_ = v514
	var v516 base.V128
	_ = v516
	var v521 base.V128
	_ = v521
	var v525 int32
	_ = v525
	var v527 base.V128
	_ = v527
	var v538 base.V128
	_ = v538
	var v542 int32
	_ = v542
	var v544 base.V128
	_ = v544
	var v548 int32
	_ = v548
	var v550 base.V128
	_ = v550
	var v562 int32
	_ = v562
	var v566 base.V128
	_ = v566
	var v581 int32
	_ = v581
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v588 int32
	_ = v588
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v592 base.V128
	_ = v592
	var v595 base.V128
	_ = v595
	var v598 base.V128
	_ = v598
	var v600 int32
	_ = v600
	var v601 base.V128
	_ = v601
	var v603 int32
	_ = v603
	var v606 int32
	_ = v606
	var v609 int32
	_ = v609
	var v611 int32
	_ = v611
	var v613 base.V128
	_ = v613
	var v616 base.V128
	_ = v616
	var v619 base.V128
	_ = v619
	var v622 base.V128
	_ = v622
	var v623 base.V128
	_ = v623
	var v624 base.V128
	_ = v624
	var v625 base.V128
	_ = v625
	var v626 base.V128
	_ = v626
	var v627 base.V128
	_ = v627
	var v628 base.V128
	_ = v628
	var v629 base.V128
	_ = v629
	var v630 base.V128
	_ = v630
	var v631 base.V128
	_ = v631
	var v632 base.V128
	_ = v632
	var v635 int32
	_ = v635
	var v638 int32
	_ = v638
	var v640 base.V128
	_ = v640
	var v643 base.V128
	_ = v643
	var v646 base.V128
	_ = v646
	var v649 base.V128
	_ = v649
	var v653 int32
	_ = v653
	var v655 base.V128
	_ = v655
	var v658 base.V128
	_ = v658
	var v661 base.V128
	_ = v661
	var v664 base.V128
	_ = v664
	var v666 base.V128
	_ = v666
	var v668 base.V128
	_ = v668
	var v670 base.V128
	_ = v670
	var v672 base.V128
	_ = v672
	var v674 base.V128
	_ = v674
	var v675 base.V128
	_ = v675
	var v676 base.V128
	_ = v676
	var v677 base.V128
	_ = v677
	var v678 base.V128
	_ = v678
	var v679 base.V128
	_ = v679
	var v680 base.V128
	_ = v680
	var v682 base.V128
	_ = v682
	var v683 base.V128
	_ = v683
	var v685 base.V128
	_ = v685
	var v687 base.V128
	_ = v687
	var v688 base.V128
	_ = v688
	var v689 base.V128
	_ = v689
	var v692 base.V128
	_ = v692
	var v700 base.V128
	_ = v700
	var v713 base.V128
	_ = v713
	var v715 base.V128
	_ = v715
	var v718 int32
	_ = v718
	var v726 base.V128
	_ = v726
	var v728 base.V128
	_ = v728
	var v730 base.V128
	_ = v730
	var v741 base.V128
	_ = v741
	var v743 base.V128
	_ = v743
	var v745 base.V128
	_ = v745
	var v749 base.V128
	_ = v749
	var v750 base.V128
	_ = v750
	var v754 int32
	_ = v754
	var v756 base.V128
	_ = v756
	var v767 base.V128
	_ = v767
	var v771 int32
	_ = v771
	var v773 base.V128
	_ = v773
	var v777 int32
	_ = v777
	var v779 base.V128
	_ = v779
	var v790 base.V128
	_ = v790
	var v792 base.V128
	_ = v792
	var v794 base.V128
	_ = v794
	var v799 base.V128
	_ = v799
	var v803 int32
	_ = v803
	var v805 base.V128
	_ = v805
	var v816 base.V128
	_ = v816
	var v820 int32
	_ = v820
	var v822 base.V128
	_ = v822
	var v826 int32
	_ = v826
	var v828 base.V128
	_ = v828
	var v838 int32
	_ = v838
	var v839 int32
	_ = v839
	v6 = int32(0)
	v10 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k0)
	v838 = int32(2)
	v25 = l0 + v838
	v27 = l1 * int32(6)
	v29 = int32(1)
	v30 = l1 << (uint(v29) % 32)
	v33 = l1 << (uint(v838) % 32)
	v34 = v25 + v33
	v36 = base.Simd_g_v128_load32_splat(m, v25, v6)
	v39 = base.Simd_g_v128_load32_lane_l1(m, v34, v6, v36)
	v42 = base.Simd_g_v128_load32_lane_l2(m, v25+v30, v6, v39)
	v44 = int32(3)
	v45 = base.Simd_g_v128_load32_lane_l3(m, v25+v27, v6, v42)
	v47 = l1 * int32(7)
	v50 = l1 * v44
	v53 = l1 * int32(5)
	v55 = v25 + l1
	v57 = base.Simd_g_v128_load32_splat(m, v55, v6)
	v60 = base.Simd_g_v128_load32_lane_l1(m, v25+v53, v6, v57)
	v63 = base.Simd_g_v128_load32_lane_l2(m, v25+v50, v6, v60)
	v66 = base.Simd_g_v128_load32_lane_l3(m, v25+v47, v6, v63)
	v67 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k1)
	v68 = base.Simd_g_i8x16_shuffle2(v45, v66, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v69 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k4)
	v70 = base.Simd_g_i8x16_shuffle2(v45, v66, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v71 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k7)
	v72 = base.Simd_g_i8x16_shuffle2(v68, v70, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	v73 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k10)
	v74 = base.Simd_g_i8x16_shuffle2(v68, v70, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	v75 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k13)
	v76 = base.Simd_g_i8x16_shuffle2(v72, v74, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k15))
	v79 = v25 + l1<<(uint(v44)%32)
	v82 = v79 + v33
	v84 = base.Simd_g_v128_load32_splat(m, v79, v6)
	v87 = base.Simd_g_v128_load32_lane_l1(m, v82, v6, v84)
	v90 = base.Simd_g_v128_load32_lane_l2(m, v79+v30, v6, v87)
	v93 = base.Simd_g_v128_load32_lane_l3(m, v79+v27, v6, v90)
	v97 = v79 + l1
	v99 = base.Simd_g_v128_load32_splat(m, v97, v6)
	v102 = base.Simd_g_v128_load32_lane_l1(m, v79+v53, v6, v99)
	v105 = base.Simd_g_v128_load32_lane_l2(m, v79+v50, v6, v102)
	v108 = base.Simd_g_v128_load32_lane_l3(m, v79+v47, v6, v105)
	v110 = base.Simd_g_i8x16_shuffle2(v93, v108, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v112 = base.Simd_g_i8x16_shuffle2(v93, v108, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v114 = base.Simd_g_i8x16_shuffle2(v110, v112, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	v116 = base.Simd_g_i8x16_shuffle2(v110, v112, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	v118 = base.Simd_g_i8x16_shuffle2(v114, v116, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k15))
	v119 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k16)
	v120 = base.Simd_g_i8x16_shuffle2(v76, v118, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k18))
	v121 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k19)
	v122 = base.Simd_g_i8x16_shuffle2(v76, v118, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k21))
	v123 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k22)
	v124 = base.Simd_g_v128_xor(v122, v123)
	v126 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k23)
	v127 = base.Simd_g_i8x16_shuffle2(v72, v74, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k25))
	v129 = base.Simd_g_i8x16_shuffle2(v114, v116, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k25))
	v131 = base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k18))
	v132 = base.Simd_g_v128_xor(v131, v123)
	v133 = base.Simd_g_i8x16_sub_sat_s(v132, v124)
	v136 = base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k21))
	v144 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v122, v131), base.Simd_g_i8x16_sub_sat_u(v131, v122))
	v157 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v133, base.Simd_g_i8x16_add_sat_s(v133, base.Simd_g_i8x16_add_sat_s(v133, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v120, v123), base.Simd_g_v128_xor(v136, v123))))), v10, base.Simd_g_i8x16_eq(v10, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v144, v144), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v120, v136), base.Simd_g_i8x16_sub_sat_u(v136, v120)), v29), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k26))), base.Simd_g_i8x16_splat(l2))))
	v159 = base.Simd_g_i8x16_add_sat_s(v157, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k27))
	v162 = int32(11)
	v170 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v124, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v159, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3)), v162), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v159, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6)), v162))), v123)
	v172 = base.Simd_g_i8x16_shuffle2(v120, v170, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v174 = base.Simd_g_i8x16_add_sat_s(v157, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k28))
	v185 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v132, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v174, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3)), v162), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v174, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6)), v162))), v123)
	v187 = base.Simd_g_i8x16_shuffle2(v185, v136, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v189 = base.Simd_g_i8x16_shuffle2(v172, v187, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v25, v6, v189)
	v193 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k29)
	v194 = base.Simd_g_i8x16_shuffle2(v189, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v55, v6, v194)
	v198 = v55 + l1
	v200 = base.Simd_g_i8x16_shuffle2(v194, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v198, v6, v200)
	base.Simd_g_v128_store32_lane_l0(m, v198+l1, v6, base.Simd_g_i8x16_shuffle2(v200, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v211 = base.Simd_g_i8x16_shuffle2(v172, v187, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v34, v6, v211)
	v215 = v34 + l1
	v217 = base.Simd_g_i8x16_shuffle2(v211, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v215, v6, v217)
	v221 = v215 + l1
	v223 = base.Simd_g_i8x16_shuffle2(v217, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v221, v6, v223)
	base.Simd_g_v128_store32_lane_l0(m, v221+l1, v6, base.Simd_g_i8x16_shuffle2(v223, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v234 = base.Simd_g_i8x16_shuffle2(v120, v170, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v236 = base.Simd_g_i8x16_shuffle2(v185, v136, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v238 = base.Simd_g_i8x16_shuffle2(v234, v236, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v79, v6, v238)
	v243 = base.Simd_g_i8x16_shuffle2(v238, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v97, v6, v243)
	v247 = v97 + l1
	v249 = base.Simd_g_i8x16_shuffle2(v243, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v247, v6, v249)
	base.Simd_g_v128_store32_lane_l0(m, v247+l1, v6, base.Simd_g_i8x16_shuffle2(v249, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v260 = base.Simd_g_i8x16_shuffle2(v234, v236, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v82, v6, v260)
	v264 = v82 + l1
	v266 = base.Simd_g_i8x16_shuffle2(v260, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v264, v6, v266)
	v270 = v264 + l1
	v272 = base.Simd_g_i8x16_shuffle2(v266, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v270, v6, v272)
	base.Simd_g_v128_store32_lane_l0(m, v270+l1, v6, base.Simd_g_i8x16_shuffle2(v272, v10, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v284 = int32(0)
	v288 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k0)
	v839 = int32(6)
	v303 = l0 + v839
	v305 = l1 * v839
	v307 = int32(1)
	v308 = l1 << (uint(v307) % 32)
	v310 = int32(2)
	v311 = l1 << (uint(v310) % 32)
	v312 = v303 + v311
	v314 = base.Simd_g_v128_load32_splat(m, v303, v284)
	v317 = base.Simd_g_v128_load32_lane_l1(m, v312, v284, v314)
	v320 = base.Simd_g_v128_load32_lane_l2(m, v303+v308, v284, v317)
	v322 = int32(3)
	v323 = base.Simd_g_v128_load32_lane_l3(m, v303+v305, v284, v320)
	v325 = l1 * int32(7)
	v328 = l1 * v322
	v331 = l1 * int32(5)
	v333 = v303 + l1
	v335 = base.Simd_g_v128_load32_splat(m, v333, v284)
	v338 = base.Simd_g_v128_load32_lane_l1(m, v303+v331, v284, v335)
	v341 = base.Simd_g_v128_load32_lane_l2(m, v303+v328, v284, v338)
	v344 = base.Simd_g_v128_load32_lane_l3(m, v303+v325, v284, v341)
	v345 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k1)
	v346 = base.Simd_g_i8x16_shuffle2(v323, v344, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v347 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k4)
	v348 = base.Simd_g_i8x16_shuffle2(v323, v344, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v349 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k7)
	v350 = base.Simd_g_i8x16_shuffle2(v346, v348, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	v351 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k10)
	v352 = base.Simd_g_i8x16_shuffle2(v346, v348, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	v353 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k13)
	v354 = base.Simd_g_i8x16_shuffle2(v350, v352, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k15))
	v357 = v303 + l1<<(uint(v322)%32)
	v360 = v357 + v311
	v362 = base.Simd_g_v128_load32_splat(m, v357, v284)
	v365 = base.Simd_g_v128_load32_lane_l1(m, v360, v284, v362)
	v368 = base.Simd_g_v128_load32_lane_l2(m, v357+v308, v284, v365)
	v371 = base.Simd_g_v128_load32_lane_l3(m, v357+v305, v284, v368)
	v375 = v357 + l1
	v377 = base.Simd_g_v128_load32_splat(m, v375, v284)
	v380 = base.Simd_g_v128_load32_lane_l1(m, v357+v331, v284, v377)
	v383 = base.Simd_g_v128_load32_lane_l2(m, v357+v328, v284, v380)
	v386 = base.Simd_g_v128_load32_lane_l3(m, v357+v325, v284, v383)
	v388 = base.Simd_g_i8x16_shuffle2(v371, v386, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v390 = base.Simd_g_i8x16_shuffle2(v371, v386, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v392 = base.Simd_g_i8x16_shuffle2(v388, v390, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	v394 = base.Simd_g_i8x16_shuffle2(v388, v390, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	v396 = base.Simd_g_i8x16_shuffle2(v392, v394, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k15))
	v397 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k16)
	v398 = base.Simd_g_i8x16_shuffle2(v354, v396, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k18))
	v399 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k19)
	v400 = base.Simd_g_i8x16_shuffle2(v354, v396, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k21))
	v401 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k22)
	v402 = base.Simd_g_v128_xor(v400, v401)
	v404 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k23)
	v405 = base.Simd_g_i8x16_shuffle2(v350, v352, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k25))
	v407 = base.Simd_g_i8x16_shuffle2(v392, v394, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k25))
	v409 = base.Simd_g_i8x16_shuffle2(v405, v407, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k18))
	v410 = base.Simd_g_v128_xor(v409, v401)
	v411 = base.Simd_g_i8x16_sub_sat_s(v410, v402)
	v414 = base.Simd_g_i8x16_shuffle2(v405, v407, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k21))
	v422 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v400, v409), base.Simd_g_i8x16_sub_sat_u(v409, v400))
	v435 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v411, base.Simd_g_i8x16_add_sat_s(v411, base.Simd_g_i8x16_add_sat_s(v411, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v398, v401), base.Simd_g_v128_xor(v414, v401))))), v288, base.Simd_g_i8x16_eq(v288, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v422, v422), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v398, v414), base.Simd_g_i8x16_sub_sat_u(v414, v398)), v307), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k26))), base.Simd_g_i8x16_splat(l2))))
	v437 = base.Simd_g_i8x16_add_sat_s(v435, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k27))
	v440 = int32(11)
	v448 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v402, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v288, v437, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3)), v440), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v288, v437, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6)), v440))), v401)
	v450 = base.Simd_g_i8x16_shuffle2(v398, v448, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v452 = base.Simd_g_i8x16_add_sat_s(v435, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k28))
	v463 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v410, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v288, v452, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3)), v440), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v288, v452, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6)), v440))), v401)
	v465 = base.Simd_g_i8x16_shuffle2(v463, v414, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v467 = base.Simd_g_i8x16_shuffle2(v450, v465, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v303, v284, v467)
	v471 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k29)
	v472 = base.Simd_g_i8x16_shuffle2(v467, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v333, v284, v472)
	v476 = v333 + l1
	v478 = base.Simd_g_i8x16_shuffle2(v472, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v476, v284, v478)
	base.Simd_g_v128_store32_lane_l0(m, v476+l1, v284, base.Simd_g_i8x16_shuffle2(v478, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v489 = base.Simd_g_i8x16_shuffle2(v450, v465, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v312, v284, v489)
	v493 = v312 + l1
	v495 = base.Simd_g_i8x16_shuffle2(v489, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v493, v284, v495)
	v499 = v493 + l1
	v501 = base.Simd_g_i8x16_shuffle2(v495, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v499, v284, v501)
	base.Simd_g_v128_store32_lane_l0(m, v499+l1, v284, base.Simd_g_i8x16_shuffle2(v501, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v512 = base.Simd_g_i8x16_shuffle2(v398, v448, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v514 = base.Simd_g_i8x16_shuffle2(v463, v414, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v516 = base.Simd_g_i8x16_shuffle2(v512, v514, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v357, v284, v516)
	v521 = base.Simd_g_i8x16_shuffle2(v516, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v375, v284, v521)
	v525 = v375 + l1
	v527 = base.Simd_g_i8x16_shuffle2(v521, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v525, v284, v527)
	base.Simd_g_v128_store32_lane_l0(m, v525+l1, v284, base.Simd_g_i8x16_shuffle2(v527, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v538 = base.Simd_g_i8x16_shuffle2(v512, v514, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v360, v284, v538)
	v542 = v360 + l1
	v544 = base.Simd_g_i8x16_shuffle2(v538, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v542, v284, v544)
	v548 = v542 + l1
	v550 = base.Simd_g_i8x16_shuffle2(v544, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v548, v284, v550)
	base.Simd_g_v128_store32_lane_l0(m, v548+l1, v284, base.Simd_g_i8x16_shuffle2(v550, v288, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v562 = int32(0)
	v566 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k0)
	v581 = l0 + int32(10)
	v583 = l1 * int32(6)
	v585 = int32(1)
	v586 = l1 << (uint(v585) % 32)
	v588 = int32(2)
	v589 = l1 << (uint(v588) % 32)
	v590 = v581 + v589
	v592 = base.Simd_g_v128_load32_splat(m, v581, v562)
	v595 = base.Simd_g_v128_load32_lane_l1(m, v590, v562, v592)
	v598 = base.Simd_g_v128_load32_lane_l2(m, v581+v586, v562, v595)
	v600 = int32(3)
	v601 = base.Simd_g_v128_load32_lane_l3(m, v581+v583, v562, v598)
	v603 = l1 * int32(7)
	v606 = l1 * v600
	v609 = l1 * int32(5)
	v611 = v581 + l1
	v613 = base.Simd_g_v128_load32_splat(m, v611, v562)
	v616 = base.Simd_g_v128_load32_lane_l1(m, v581+v609, v562, v613)
	v619 = base.Simd_g_v128_load32_lane_l2(m, v581+v606, v562, v616)
	v622 = base.Simd_g_v128_load32_lane_l3(m, v581+v603, v562, v619)
	v623 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k1)
	v624 = base.Simd_g_i8x16_shuffle2(v601, v622, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v625 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k4)
	v626 = base.Simd_g_i8x16_shuffle2(v601, v622, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v627 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k7)
	v628 = base.Simd_g_i8x16_shuffle2(v624, v626, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	v629 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k10)
	v630 = base.Simd_g_i8x16_shuffle2(v624, v626, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	v631 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k13)
	v632 = base.Simd_g_i8x16_shuffle2(v628, v630, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k15))
	v635 = v581 + l1<<(uint(v600)%32)
	v638 = v635 + v589
	v640 = base.Simd_g_v128_load32_splat(m, v635, v562)
	v643 = base.Simd_g_v128_load32_lane_l1(m, v638, v562, v640)
	v646 = base.Simd_g_v128_load32_lane_l2(m, v635+v586, v562, v643)
	v649 = base.Simd_g_v128_load32_lane_l3(m, v635+v583, v562, v646)
	v653 = v635 + l1
	v655 = base.Simd_g_v128_load32_splat(m, v653, v562)
	v658 = base.Simd_g_v128_load32_lane_l1(m, v635+v609, v562, v655)
	v661 = base.Simd_g_v128_load32_lane_l2(m, v635+v606, v562, v658)
	v664 = base.Simd_g_v128_load32_lane_l3(m, v635+v603, v562, v661)
	v666 = base.Simd_g_i8x16_shuffle2(v649, v664, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v668 = base.Simd_g_i8x16_shuffle2(v649, v664, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v670 = base.Simd_g_i8x16_shuffle2(v666, v668, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	v672 = base.Simd_g_i8x16_shuffle2(v666, v668, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	v674 = base.Simd_g_i8x16_shuffle2(v670, v672, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k14), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k15))
	v675 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k16)
	v676 = base.Simd_g_i8x16_shuffle2(v632, v674, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k18))
	v677 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k19)
	v678 = base.Simd_g_i8x16_shuffle2(v632, v674, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k21))
	v679 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k22)
	v680 = base.Simd_g_v128_xor(v678, v679)
	v682 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k23)
	v683 = base.Simd_g_i8x16_shuffle2(v628, v630, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k25))
	v685 = base.Simd_g_i8x16_shuffle2(v670, v672, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k24), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k25))
	v687 = base.Simd_g_i8x16_shuffle2(v683, v685, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k17), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k18))
	v688 = base.Simd_g_v128_xor(v687, v679)
	v689 = base.Simd_g_i8x16_sub_sat_s(v688, v680)
	v692 = base.Simd_g_i8x16_shuffle2(v683, v685, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k20), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k21))
	v700 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v678, v687), base.Simd_g_i8x16_sub_sat_u(v687, v678))
	v713 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v689, base.Simd_g_i8x16_add_sat_s(v689, base.Simd_g_i8x16_add_sat_s(v689, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v676, v679), base.Simd_g_v128_xor(v692, v679))))), v566, base.Simd_g_i8x16_eq(v566, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v700, v700), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v676, v692), base.Simd_g_i8x16_sub_sat_u(v692, v676)), v585), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k26))), base.Simd_g_i8x16_splat(l2))))
	v715 = base.Simd_g_i8x16_add_sat_s(v713, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k27))
	v718 = int32(11)
	v726 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v680, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v566, v715, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3)), v718), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v566, v715, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6)), v718))), v679)
	v728 = base.Simd_g_i8x16_shuffle2(v676, v726, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v730 = base.Simd_g_i8x16_add_sat_s(v713, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k28))
	v741 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v688, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v566, v730, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3)), v718), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v566, v730, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6)), v718))), v679)
	v743 = base.Simd_g_i8x16_shuffle2(v741, v692, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k2), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k3))
	v745 = base.Simd_g_i8x16_shuffle2(v728, v743, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v581, v562, v745)
	v749 = base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k29)
	v750 = base.Simd_g_i8x16_shuffle2(v745, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v611, v562, v750)
	v754 = v611 + l1
	v756 = base.Simd_g_i8x16_shuffle2(v750, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v754, v562, v756)
	base.Simd_g_v128_store32_lane_l0(m, v754+l1, v562, base.Simd_g_i8x16_shuffle2(v756, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v767 = base.Simd_g_i8x16_shuffle2(v728, v743, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v590, v562, v767)
	v771 = v590 + l1
	v773 = base.Simd_g_i8x16_shuffle2(v767, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v771, v562, v773)
	v777 = v771 + l1
	v779 = base.Simd_g_i8x16_shuffle2(v773, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v777, v562, v779)
	base.Simd_g_v128_store32_lane_l0(m, v777+l1, v562, base.Simd_g_i8x16_shuffle2(v779, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v790 = base.Simd_g_i8x16_shuffle2(v676, v726, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v792 = base.Simd_g_i8x16_shuffle2(v741, v692, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k5), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k6))
	v794 = base.Simd_g_i8x16_shuffle2(v790, v792, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v635, v562, v794)
	v799 = base.Simd_g_i8x16_shuffle2(v794, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v653, v562, v799)
	v803 = v653 + l1
	v805 = base.Simd_g_i8x16_shuffle2(v799, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v803, v562, v805)
	base.Simd_g_v128_store32_lane_l0(m, v803+l1, v562, base.Simd_g_i8x16_shuffle2(v805, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	v816 = base.Simd_g_i8x16_shuffle2(v790, v792, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k11), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v638, v562, v816)
	v820 = v638 + l1
	v822 = base.Simd_g_i8x16_shuffle2(v816, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v820, v562, v822)
	v826 = v820 + l1
	v828 = base.Simd_g_i8x16_shuffle2(v822, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31))
	base.Simd_g_v128_store32_lane_l0(m, v826, v562, v828)
	base.Simd_g_v128_store32_lane_l0(m, v826+l1, v562, base.Simd_g_i8x16_shuffle2(v828, v566, base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k30), base.Simd_g_const(&F_SimpleHFilter16i_SSE2__k31)))
	return
}

var F_SimpleHFilter16i_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SimpleHFilter16i_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SimpleHFilter16i_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SimpleHFilter16i_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_SimpleHFilter16i_SSE2__k4 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_SimpleHFilter16i_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_SimpleHFilter16i_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_SimpleHFilter16i_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_SimpleHFilter16i_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_SimpleHFilter16i_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_SimpleHFilter16i_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_SimpleHFilter16i_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_SimpleHFilter16i_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_SimpleHFilter16i_SSE2__k13 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_SimpleHFilter16i_SSE2__k14 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_SimpleHFilter16i_SSE2__k15 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_SimpleHFilter16i_SSE2__k16 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_SimpleHFilter16i_SSE2__k17 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_SimpleHFilter16i_SSE2__k18 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_SimpleHFilter16i_SSE2__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_SimpleHFilter16i_SSE2__k20 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_SimpleHFilter16i_SSE2__k21 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_SimpleHFilter16i_SSE2__k22 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SimpleHFilter16i_SSE2__k23 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_SimpleHFilter16i_SSE2__k24 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_SimpleHFilter16i_SSE2__k25 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_SimpleHFilter16i_SSE2__k26 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_SimpleHFilter16i_SSE2__k27 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_SimpleHFilter16i_SSE2__k28 = [2]uint64{0x404040404040404, 0x404040404040404}
var F_SimpleHFilter16i_SSE2__k29 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_SimpleHFilter16i_SSE2__k30 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_SimpleHFilter16i_SSE2__k31 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_SimpleVFilter16_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v4 int32
	_ = v4
	var v5 base.V128
	_ = v5
	var v13 int32
	_ = v13
	var v15 base.V128
	_ = v15
	var v16 base.V128
	_ = v16
	var v17 base.V128
	_ = v17
	var v20 base.V128
	_ = v20
	var v21 base.V128
	_ = v21
	var v22 base.V128
	_ = v22
	var v23 int32
	_ = v23
	var v27 base.V128
	_ = v27
	var v31 base.V128
	_ = v31
	var v39 base.V128
	_ = v39
	var v52 base.V128
	_ = v52
	var v54 base.V128
	_ = v54
	var v55 base.V128
	_ = v55
	var v57 int32
	_ = v57
	var v59 base.V128
	_ = v59
	var v66 base.V128
	_ = v66
	var v70 base.V128
	_ = v70
	var v81 base.V128
	_ = v81
	v4 = int32(0)
	v5 = base.Simd_g_const(&F_SimpleVFilter16_SSE2__k0)
	v13 = l0 - l1
	v15 = base.Simd_g_v128_load(m, v13, v4)
	v16 = base.Simd_g_const(&F_SimpleVFilter16_SSE2__k1)
	v17 = base.Simd_g_v128_xor(v15, v16)
	v20 = base.Simd_g_v128_load(m, l0, v4)
	v21 = base.Simd_g_v128_xor(v20, v16)
	v22 = base.Simd_g_i8x16_sub_sat_s(v21, v17)
	v23 = int32(1)
	v27 = base.Simd_g_v128_load(m, l0-l1<<(uint(v23)%32), v4)
	v31 = base.Simd_g_v128_load(m, l0+l1, v4)
	v39 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v15, v20), base.Simd_g_i8x16_sub_sat_u(v20, v15))
	v52 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v22, base.Simd_g_i8x16_add_sat_s(v22, base.Simd_g_i8x16_add_sat_s(v22, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v27, v16), base.Simd_g_v128_xor(v31, v16))))), v5, base.Simd_g_i8x16_eq(v5, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v39, v39), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v27, v31), base.Simd_g_i8x16_sub_sat_u(v31, v27)), v23), base.Simd_g_const(&F_SimpleVFilter16_SSE2__k2))), base.Simd_g_i8x16_splat(l2))))
	v54 = base.Simd_g_i8x16_add_sat_s(v52, base.Simd_g_const(&F_SimpleVFilter16_SSE2__k3))
	v55 = base.Simd_g_const(&F_SimpleVFilter16_SSE2__k4)
	v57 = int32(11)
	v59 = base.Simd_g_const(&F_SimpleVFilter16_SSE2__k5)
	v66 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v17, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v5, v54, base.Simd_g_const(&F_SimpleVFilter16_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16_SSE2__k7)), v57), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v5, v54, base.Simd_g_const(&F_SimpleVFilter16_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16_SSE2__k9)), v57))), v16)
	base.Simd_g_v128_store(m, v13, v4, v66)
	v70 = base.Simd_g_i8x16_add_sat_s(v52, base.Simd_g_const(&F_SimpleVFilter16_SSE2__k10))
	v81 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v21, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v5, v70, base.Simd_g_const(&F_SimpleVFilter16_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16_SSE2__k7)), v57), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v5, v70, base.Simd_g_const(&F_SimpleVFilter16_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16_SSE2__k9)), v57))), v16)
	base.Simd_g_v128_store(m, l0, v4, v81)
	return
}

var F_SimpleVFilter16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SimpleVFilter16_SSE2__k1 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SimpleVFilter16_SSE2__k2 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_SimpleVFilter16_SSE2__k3 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_SimpleVFilter16_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SimpleVFilter16_SSE2__k5 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_SimpleVFilter16_SSE2__k6 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SimpleVFilter16_SSE2__k7 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_SimpleVFilter16_SSE2__k8 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_SimpleVFilter16_SSE2__k9 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_SimpleVFilter16_SSE2__k10 = [2]uint64{0x404040404040404, 0x404040404040404}

func F_SimpleVFilter16i_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v4 int32
	_ = v4
	var v6 base.V128
	_ = v6
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v36 base.V128
	_ = v36
	var v40 base.V128
	_ = v40
	var v48 base.V128
	_ = v48
	var v55 base.V128
	_ = v55
	var v58 base.V128
	_ = v58
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 int32
	_ = v66
	var v68 base.V128
	_ = v68
	var v75 base.V128
	_ = v75
	var v78 base.V128
	_ = v78
	var v79 base.V128
	_ = v79
	var v90 base.V128
	_ = v90
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v108 base.V128
	_ = v108
	var v116 base.V128
	_ = v116
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v139 base.V128
	_ = v139
	var v142 base.V128
	_ = v142
	var v153 base.V128
	_ = v153
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v159 base.V128
	_ = v159
	var v160 base.V128
	_ = v160
	var v162 base.V128
	_ = v162
	var v163 base.V128
	_ = v163
	var v164 base.V128
	_ = v164
	var v167 base.V128
	_ = v167
	var v171 base.V128
	_ = v171
	var v179 base.V128
	_ = v179
	var v190 base.V128
	_ = v190
	var v191 base.V128
	_ = v191
	var v202 base.V128
	_ = v202
	var v205 base.V128
	_ = v205
	var v216 base.V128
	_ = v216
	v4 = int32(0)
	v6 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k0)
	v20 = l1 << (uint(int32(2)) % 32)
	v21 = l0 + v20
	v22 = v21 - l1
	v24 = base.Simd_g_v128_load(m, v22, v4)
	v25 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k1)
	v26 = base.Simd_g_v128_xor(v24, v25)
	v29 = base.Simd_g_v128_load(m, v21, v4)
	v30 = base.Simd_g_v128_xor(v29, v25)
	v31 = base.Simd_g_i8x16_sub_sat_s(v30, v26)
	v32 = int32(1)
	v33 = l1 << (uint(v32) % 32)
	v36 = base.Simd_g_v128_load(m, v21-v33, v4)
	v40 = base.Simd_g_v128_load(m, v21+l1, v4)
	v48 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v24, v29), base.Simd_g_i8x16_sub_sat_u(v29, v24))
	v55 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k2)
	v58 = base.Simd_g_i8x16_splat(l2)
	v61 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v31, base.Simd_g_i8x16_add_sat_s(v31, base.Simd_g_i8x16_add_sat_s(v31, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v36, v25), base.Simd_g_v128_xor(v40, v25))))), v6, base.Simd_g_i8x16_eq(v6, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v48, v48), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v36, v40), base.Simd_g_i8x16_sub_sat_u(v40, v36)), v32), v55)), v58)))
	v62 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k3)
	v63 = base.Simd_g_i8x16_add_sat_s(v61, v62)
	v64 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k4)
	v66 = int32(11)
	v68 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k5)
	v75 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v26, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v63, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k7)), v66), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v63, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k9)), v66))), v25)
	base.Simd_g_v128_store(m, v22, v4, v75)
	v78 = base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k10)
	v79 = base.Simd_g_i8x16_add_sat_s(v61, v78)
	v90 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v30, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v79, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k7)), v66), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v79, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k9)), v66))), v25)
	base.Simd_g_v128_store(m, v21, v4, v90)
	v93 = v21 + v20
	v94 = v93 - l1
	v96 = base.Simd_g_v128_load(m, v94, v4)
	v97 = base.Simd_g_v128_xor(v96, v25)
	v99 = base.Simd_g_v128_load(m, v93, v4)
	v100 = base.Simd_g_v128_xor(v99, v25)
	v101 = base.Simd_g_i8x16_sub_sat_s(v100, v97)
	v104 = base.Simd_g_v128_load(m, v93-v33, v4)
	v108 = base.Simd_g_v128_load(m, v93+l1, v4)
	v116 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v96, v99), base.Simd_g_i8x16_sub_sat_u(v99, v96))
	v127 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v101, base.Simd_g_i8x16_add_sat_s(v101, base.Simd_g_i8x16_add_sat_s(v101, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v104, v25), base.Simd_g_v128_xor(v108, v25))))), v6, base.Simd_g_i8x16_eq(v6, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v116, v116), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v104, v108), base.Simd_g_i8x16_sub_sat_u(v108, v104)), v32), v55)), v58)))
	v128 = base.Simd_g_i8x16_add_sat_s(v127, v62)
	v139 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v97, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v128, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k7)), v66), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v128, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k9)), v66))), v25)
	base.Simd_g_v128_store(m, v94, v4, v139)
	v142 = base.Simd_g_i8x16_add_sat_s(v127, v78)
	v153 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v100, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v142, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k7)), v66), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v142, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k9)), v66))), v25)
	base.Simd_g_v128_store(m, v93, v4, v153)
	v156 = v93 + v20
	v157 = v156 - l1
	v159 = base.Simd_g_v128_load(m, v157, v4)
	v160 = base.Simd_g_v128_xor(v159, v25)
	v162 = base.Simd_g_v128_load(m, v156, v4)
	v163 = base.Simd_g_v128_xor(v162, v25)
	v164 = base.Simd_g_i8x16_sub_sat_s(v163, v160)
	v167 = base.Simd_g_v128_load(m, v156-v33, v4)
	v171 = base.Simd_g_v128_load(m, v156+l1, v4)
	v179 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v159, v162), base.Simd_g_i8x16_sub_sat_u(v162, v159))
	v190 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(v164, base.Simd_g_i8x16_add_sat_s(v164, base.Simd_g_i8x16_add_sat_s(v164, base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v167, v25), base.Simd_g_v128_xor(v171, v25))))), v6, base.Simd_g_i8x16_eq(v6, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v179, v179), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v167, v171), base.Simd_g_i8x16_sub_sat_u(v171, v167)), v32), v55)), v58)))
	v191 = base.Simd_g_i8x16_add_sat_s(v190, v62)
	v202 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v160, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v191, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k7)), v66), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v191, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k9)), v66))), v25)
	base.Simd_g_v128_store(m, v157, v4, v202)
	v205 = base.Simd_g_i8x16_add_sat_s(v190, v78)
	v216 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v163, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v205, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k6), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k7)), v66), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v6, v205, base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k8), base.Simd_g_const(&F_SimpleVFilter16i_SSE2__k9)), v66))), v25)
	base.Simd_g_v128_store(m, v156, v4, v216)
	return
}

var F_SimpleVFilter16i_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SimpleVFilter16i_SSE2__k1 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_SimpleVFilter16i_SSE2__k2 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_SimpleVFilter16i_SSE2__k3 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_SimpleVFilter16i_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_SimpleVFilter16i_SSE2__k5 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_SimpleVFilter16i_SSE2__k6 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_SimpleVFilter16i_SSE2__k7 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_SimpleVFilter16i_SSE2__k8 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_SimpleVFilter16i_SSE2__k9 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_SimpleVFilter16i_SSE2__k10 = [2]uint64{0x404040404040404, 0x404040404040404}

func F_StoreHuffmanCode(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v38 int32
	_ = v38
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v122 int64
	_ = v122
	var v125 int32
	_ = v125
	var v127 int64
	_ = v127
	var v130 int32
	_ = v130
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v155 int32
	_ = v155
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v159 int32
	_ = v159
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v229 int32
	_ = v229
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v247 int32
	_ = v247
	var v248 int32
	_ = v248
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v269 int32
	_ = v269
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v282 int64
	_ = v282
	var v285 int32
	_ = v285
	var v287 int64
	_ = v287
	var v290 int32
	_ = v290
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v305 int32
	_ = v305
	var v306 int32
	_ = v306
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v319 int32
	_ = v319
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v332 int32
	_ = v332
	var v335 int32
	_ = v335
	var v337 int32
	_ = v337
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v382 int32
	_ = v382
	var v383 int32
	_ = v383
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v409 int32
	_ = v409
	var v410 int32
	_ = v410
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v414 int32
	_ = v414
	var v415 int32
	_ = v415
	var v419 int32
	_ = v419
	var v420 int32
	_ = v420
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v431 int32
	_ = v431
	var v432 int32
	_ = v432
	var v435 int64
	_ = v435
	var v438 int32
	_ = v438
	var v440 int64
	_ = v440
	var v443 int32
	_ = v443
	var v450 int32
	_ = v450
	var v452 int32
	_ = v452
	var v456 int32
	_ = v456
	var v457 int32
	_ = v457
	var v458 int32
	_ = v458
	var v459 int32
	_ = v459
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v472 int32
	_ = v472
	var v477 int32
	_ = v477
	var v478 int32
	_ = v478
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v490 int32
	_ = v490
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v496 int32
	_ = v496
	var v497 int32
	_ = v497
	var v534 int32
	_ = v534
	var v535 int32
	_ = v535
	var v548 int32
	_ = v548
	var v549 int32
	_ = v549
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v562 int32
	_ = v562
	var v563 int32
	_ = v563
	var v564 int32
	_ = v564
	var v565 int32
	_ = v565
	var v567 int32
	_ = v567
	var v568 int32
	_ = v568
	var v572 int32
	_ = v572
	var v573 int32
	_ = v573
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v584 int32
	_ = v584
	var v585 int32
	_ = v585
	var v588 int64
	_ = v588
	var v591 int32
	_ = v591
	var v593 int64
	_ = v593
	var v596 int32
	_ = v596
	var v603 int32
	_ = v603
	var v605 int32
	_ = v605
	var v609 int32
	_ = v609
	var v610 int32
	_ = v610
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v621 int32
	_ = v621
	var v622 int32
	_ = v622
	var v623 int32
	_ = v623
	var v625 int32
	_ = v625
	var v630 int32
	_ = v630
	var v631 int32
	_ = v631
	var v638 int32
	_ = v638
	var v641 int32
	_ = v641
	var v643 int32
	_ = v643
	var v647 int32
	_ = v647
	var v648 int32
	_ = v648
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v689 int32
	_ = v689
	var v701 int32
	_ = v701
	var v702 int32
	_ = v702
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v715 int32
	_ = v715
	var v716 int32
	_ = v716
	var v717 int32
	_ = v717
	var v718 int32
	_ = v718
	var v720 int32
	_ = v720
	var v721 int32
	_ = v721
	var v725 int32
	_ = v725
	var v726 int32
	_ = v726
	var v727 int32
	_ = v727
	var v728 int32
	_ = v728
	var v737 int32
	_ = v737
	var v738 int32
	_ = v738
	var v741 int64
	_ = v741
	var v744 int32
	_ = v744
	var v746 int64
	_ = v746
	var v749 int32
	_ = v749
	var v756 int32
	_ = v756
	var v758 int32
	_ = v758
	var v762 int32
	_ = v762
	var v763 int32
	_ = v763
	var v764 int32
	_ = v764
	var v765 int32
	_ = v765
	var v774 int32
	_ = v774
	var v775 int32
	_ = v775
	var v776 int32
	_ = v776
	var v778 int32
	_ = v778
	var v783 int32
	_ = v783
	var v784 int32
	_ = v784
	var v791 int32
	_ = v791
	var v794 int32
	_ = v794
	var v796 int32
	_ = v796
	var v800 int32
	_ = v800
	var v801 int32
	_ = v801
	var v802 int32
	_ = v802
	var v803 int32
	_ = v803
	var v842 int32
	_ = v842
	var v854 int32
	_ = v854
	var v855 int32
	_ = v855
	var v859 int32
	_ = v859
	var v860 int32
	_ = v860
	var v868 int32
	_ = v868
	var v869 int32
	_ = v869
	var v870 int32
	_ = v870
	var v871 int32
	_ = v871
	var v873 int32
	_ = v873
	var v874 int32
	_ = v874
	var v878 int32
	_ = v878
	var v879 int32
	_ = v879
	var v880 int32
	_ = v880
	var v881 int32
	_ = v881
	var v890 int32
	_ = v890
	var v891 int32
	_ = v891
	var v894 int64
	_ = v894
	var v897 int32
	_ = v897
	var v899 int64
	_ = v899
	var v902 int32
	_ = v902
	var v909 int32
	_ = v909
	var v911 int32
	_ = v911
	var v915 int32
	_ = v915
	var v916 int32
	_ = v916
	var v917 int32
	_ = v917
	var v918 int32
	_ = v918
	var v927 int32
	_ = v927
	var v928 int32
	_ = v928
	var v929 int32
	_ = v929
	var v931 int32
	_ = v931
	var v936 int32
	_ = v936
	var v937 int32
	_ = v937
	var v944 int32
	_ = v944
	var v947 int32
	_ = v947
	var v949 int32
	_ = v949
	var v953 int32
	_ = v953
	var v954 int32
	_ = v954
	var v955 int32
	_ = v955
	var v956 int32
	_ = v956
	var v995 int32
	_ = v995
	var v998 base.V128
	_ = v998
	var v999 int32
	_ = v999
	var v1009 int32
	_ = v1009
	var v1011 int32
	_ = v1011
	var v1021 int32
	_ = v1021
	var v1033 int32
	_ = v1033
	var v1034 int32
	_ = v1034
	var v1038 int32
	_ = v1038
	var v1039 int32
	_ = v1039
	var v1047 int32
	_ = v1047
	var v1048 int32
	_ = v1048
	var v1049 int32
	_ = v1049
	var v1050 int32
	_ = v1050
	var v1052 int32
	_ = v1052
	var v1053 int32
	_ = v1053
	var v1057 int32
	_ = v1057
	var v1058 int32
	_ = v1058
	var v1059 int32
	_ = v1059
	var v1060 int32
	_ = v1060
	var v1069 int32
	_ = v1069
	var v1070 int32
	_ = v1070
	var v1073 int64
	_ = v1073
	var v1076 int32
	_ = v1076
	var v1078 int64
	_ = v1078
	var v1081 int32
	_ = v1081
	var v1088 int32
	_ = v1088
	var v1090 int32
	_ = v1090
	var v1094 int32
	_ = v1094
	var v1095 int32
	_ = v1095
	var v1096 int32
	_ = v1096
	var v1097 int32
	_ = v1097
	var v1106 int32
	_ = v1106
	var v1107 int32
	_ = v1107
	var v1108 int32
	_ = v1108
	var v1110 int32
	_ = v1110
	var v1115 int32
	_ = v1115
	var v1116 int32
	_ = v1116
	var v1123 int32
	_ = v1123
	var v1126 int32
	_ = v1126
	var v1128 int32
	_ = v1128
	var v1132 int32
	_ = v1132
	var v1133 int32
	_ = v1133
	var v1134 int32
	_ = v1134
	var v1135 int32
	_ = v1135
	var v1172 int32
	_ = v1172
	var v1304 int32
	_ = v1304
	var v1317 int32
	_ = v1317
	var v1319 int32
	_ = v1319
	var v1323 int32
	_ = v1323
	var v1324 int32
	_ = v1324
	var v1325 int32
	_ = v1325
	var v1327 int32
	_ = v1327
	var v1328 int32
	_ = v1328
	var v1329 int32
	_ = v1329
	var v1336 int32
	_ = v1336
	var v1339 int32
	_ = v1339
	var v1340 int32
	_ = v1340
	var v1346 int32
	_ = v1346
	var v1348 int32
	_ = v1348
	var v1351 int32
	_ = v1351
	var v1352 int32
	_ = v1352
	var v1360 int32
	_ = v1360
	var v1363 int32
	_ = v1363
	var v1364 int32
	_ = v1364
	var v1371 int32
	_ = v1371
	var v1382 int32
	_ = v1382
	var v1395 int32
	_ = v1395
	var v1397 int32
	_ = v1397
	var v1404 int32
	_ = v1404
	var v1405 int32
	_ = v1405
	var v1407 int32
	_ = v1407
	var v1408 int32
	_ = v1408
	var v1415 int32
	_ = v1415
	var v1436 int32
	_ = v1436
	var v1437 int32
	_ = v1437
	var v1441 int32
	_ = v1441
	var v1445 int32
	_ = v1445
	var v1449 int32
	_ = v1449
	var v1453 int32
	_ = v1453
	var v1457 int32
	_ = v1457
	var v1461 int32
	_ = v1461
	var v1465 int32
	_ = v1465
	var v1469 int32
	_ = v1469
	var v1473 int32
	_ = v1473
	var v1477 int32
	_ = v1477
	var v1481 int32
	_ = v1481
	var v1485 int32
	_ = v1485
	var v1491 int32
	_ = v1491
	var v1492 int32
	_ = v1492
	var v1493 int32
	_ = v1493
	var v1495 int32
	_ = v1495
	var v1496 int32
	_ = v1496
	var v1508 int32
	_ = v1508
	var v1509 int32
	_ = v1509
	var v1513 int32
	_ = v1513
	var v1514 int32
	_ = v1514
	var v1522 int32
	_ = v1522
	var v1523 int32
	_ = v1523
	var v1524 int32
	_ = v1524
	var v1525 int32
	_ = v1525
	var v1527 int32
	_ = v1527
	var v1528 int32
	_ = v1528
	var v1532 int32
	_ = v1532
	var v1533 int32
	_ = v1533
	var v1534 int32
	_ = v1534
	var v1535 int32
	_ = v1535
	var v1544 int32
	_ = v1544
	var v1545 int32
	_ = v1545
	var v1548 int64
	_ = v1548
	var v1551 int32
	_ = v1551
	var v1553 int64
	_ = v1553
	var v1556 int32
	_ = v1556
	var v1563 int32
	_ = v1563
	var v1565 int32
	_ = v1565
	var v1569 int32
	_ = v1569
	var v1570 int32
	_ = v1570
	var v1571 int32
	_ = v1571
	var v1572 int32
	_ = v1572
	var v1581 int32
	_ = v1581
	var v1582 int32
	_ = v1582
	var v1583 int32
	_ = v1583
	var v1585 int32
	_ = v1585
	var v1590 int32
	_ = v1590
	var v1591 int32
	_ = v1591
	var v1598 int32
	_ = v1598
	var v1601 int32
	_ = v1601
	var v1603 int32
	_ = v1603
	var v1607 int32
	_ = v1607
	var v1608 int32
	_ = v1608
	var v1609 int32
	_ = v1609
	var v1610 int32
	_ = v1610
	var v1655 int32
	_ = v1655
	var v1662 int32
	_ = v1662
	var v1666 int32
	_ = v1666
	var v1668 int32
	_ = v1668
	var v1669 int32
	_ = v1669
	var v1681 int32
	_ = v1681
	var v1682 int32
	_ = v1682
	var v1686 int32
	_ = v1686
	var v1687 int32
	_ = v1687
	var v1695 int32
	_ = v1695
	var v1696 int32
	_ = v1696
	var v1697 int32
	_ = v1697
	var v1698 int32
	_ = v1698
	var v1700 int32
	_ = v1700
	var v1701 int32
	_ = v1701
	var v1705 int32
	_ = v1705
	var v1706 int32
	_ = v1706
	var v1707 int32
	_ = v1707
	var v1708 int32
	_ = v1708
	var v1717 int32
	_ = v1717
	var v1718 int32
	_ = v1718
	var v1721 int64
	_ = v1721
	var v1724 int32
	_ = v1724
	var v1726 int64
	_ = v1726
	var v1729 int32
	_ = v1729
	var v1736 int32
	_ = v1736
	var v1738 int32
	_ = v1738
	var v1742 int32
	_ = v1742
	var v1743 int32
	_ = v1743
	var v1744 int32
	_ = v1744
	var v1745 int32
	_ = v1745
	var v1754 int32
	_ = v1754
	var v1755 int32
	_ = v1755
	var v1756 int32
	_ = v1756
	var v1758 int32
	_ = v1758
	var v1763 int32
	_ = v1763
	var v1764 int32
	_ = v1764
	var v1771 int32
	_ = v1771
	var v1774 int32
	_ = v1774
	var v1776 int32
	_ = v1776
	var v1780 int32
	_ = v1780
	var v1781 int32
	_ = v1781
	var v1782 int32
	_ = v1782
	var v1783 int32
	_ = v1783
	var v1821 int32
	_ = v1821
	var v1823 int32
	_ = v1823
	var v1826 int32
	_ = v1826
	var v1833 int32
	_ = v1833
	var v1834 int32
	_ = v1834
	var v1835 int32
	_ = v1835
	var v1840 int32
	_ = v1840
	var v1846 int32
	_ = v1846
	var v1851 int32
	_ = v1851
	var v1852 int32
	_ = v1852
	var v1859 int32
	_ = v1859
	var v1861 int32
	_ = v1861
	var v1866 int32
	_ = v1866
	var v1868 int32
	_ = v1868
	var v1870 int32
	_ = v1870
	var v1877 int32
	_ = v1877
	var v1878 int32
	_ = v1878
	var v1906 int32
	_ = v1906
	var v1907 int32
	_ = v1907
	var v1910 int32
	_ = v1910
	var v1913 int32
	_ = v1913
	var v1925 int32
	_ = v1925
	var v1926 int32
	_ = v1926
	var v1933 int32
	_ = v1933
	var v1943 int32
	_ = v1943
	var v1945 int32
	_ = v1945
	var v1955 int32
	_ = v1955
	var v1957 int32
	_ = v1957
	var v1960 int32
	_ = v1960
	var v1962 int32
	_ = v1962
	var v1975 int32
	_ = v1975
	var v1976 int32
	_ = v1976
	var v1980 int32
	_ = v1980
	var v1981 int32
	_ = v1981
	var v1989 int32
	_ = v1989
	var v1990 int32
	_ = v1990
	var v1991 int32
	_ = v1991
	var v1992 int32
	_ = v1992
	var v1994 int32
	_ = v1994
	var v1995 int32
	_ = v1995
	var v1999 int32
	_ = v1999
	var v2000 int32
	_ = v2000
	var v2001 int32
	_ = v2001
	var v2002 int32
	_ = v2002
	var v2011 int32
	_ = v2011
	var v2012 int32
	_ = v2012
	var v2015 int64
	_ = v2015
	var v2018 int32
	_ = v2018
	var v2020 int64
	_ = v2020
	var v2023 int32
	_ = v2023
	var v2030 int32
	_ = v2030
	var v2032 int32
	_ = v2032
	var v2036 int32
	_ = v2036
	var v2037 int32
	_ = v2037
	var v2038 int32
	_ = v2038
	var v2039 int32
	_ = v2039
	var v2048 int32
	_ = v2048
	var v2049 int32
	_ = v2049
	var v2050 int32
	_ = v2050
	var v2052 int32
	_ = v2052
	var v2057 int32
	_ = v2057
	var v2058 int32
	_ = v2058
	var v2065 int32
	_ = v2065
	var v2068 int32
	_ = v2068
	var v2070 int32
	_ = v2070
	var v2074 int32
	_ = v2074
	var v2075 int32
	_ = v2075
	var v2076 int32
	_ = v2076
	var v2077 int32
	_ = v2077
	var v2116 int32
	_ = v2116
	var v2119 int32
	_ = v2119
	var v2120 int32
	_ = v2120
	var v2132 int32
	_ = v2132
	var v2133 int32
	_ = v2133
	var v2137 int32
	_ = v2137
	var v2138 int32
	_ = v2138
	var v2146 int32
	_ = v2146
	var v2147 int32
	_ = v2147
	var v2148 int32
	_ = v2148
	var v2149 int32
	_ = v2149
	var v2151 int32
	_ = v2151
	var v2152 int32
	_ = v2152
	var v2156 int32
	_ = v2156
	var v2157 int32
	_ = v2157
	var v2158 int32
	_ = v2158
	var v2159 int32
	_ = v2159
	var v2168 int32
	_ = v2168
	var v2169 int32
	_ = v2169
	var v2172 int64
	_ = v2172
	var v2175 int32
	_ = v2175
	var v2177 int64
	_ = v2177
	var v2180 int32
	_ = v2180
	var v2187 int32
	_ = v2187
	var v2189 int32
	_ = v2189
	var v2193 int32
	_ = v2193
	var v2194 int32
	_ = v2194
	var v2195 int32
	_ = v2195
	var v2196 int32
	_ = v2196
	var v2205 int32
	_ = v2205
	var v2206 int32
	_ = v2206
	var v2207 int32
	_ = v2207
	var v2209 int32
	_ = v2209
	var v2214 int32
	_ = v2214
	var v2215 int32
	_ = v2215
	var v2222 int32
	_ = v2222
	var v2225 int32
	_ = v2225
	var v2227 int32
	_ = v2227
	var v2231 int32
	_ = v2231
	var v2232 int32
	_ = v2232
	var v2233 int32
	_ = v2233
	var v2234 int32
	_ = v2234
	var v2272 int32
	_ = v2272
	var v2275 int32
	_ = v2275
	var v2277 int32
	_ = v2277
	var v2278 int32
	_ = v2278
	var v2290 int32
	_ = v2290
	var v2291 int32
	_ = v2291
	var v2295 int32
	_ = v2295
	var v2296 int32
	_ = v2296
	var v2304 int32
	_ = v2304
	var v2305 int32
	_ = v2305
	var v2306 int32
	_ = v2306
	var v2307 int32
	_ = v2307
	var v2309 int32
	_ = v2309
	var v2310 int32
	_ = v2310
	var v2314 int32
	_ = v2314
	var v2315 int32
	_ = v2315
	var v2316 int32
	_ = v2316
	var v2317 int32
	_ = v2317
	var v2326 int32
	_ = v2326
	var v2327 int32
	_ = v2327
	var v2330 int64
	_ = v2330
	var v2333 int32
	_ = v2333
	var v2335 int64
	_ = v2335
	var v2338 int32
	_ = v2338
	var v2345 int32
	_ = v2345
	var v2347 int32
	_ = v2347
	var v2351 int32
	_ = v2351
	var v2352 int32
	_ = v2352
	var v2353 int32
	_ = v2353
	var v2354 int32
	_ = v2354
	var v2363 int32
	_ = v2363
	var v2364 int32
	_ = v2364
	var v2365 int32
	_ = v2365
	var v2367 int32
	_ = v2367
	var v2372 int32
	_ = v2372
	var v2373 int32
	_ = v2373
	var v2380 int32
	_ = v2380
	var v2383 int32
	_ = v2383
	var v2385 int32
	_ = v2385
	var v2389 int32
	_ = v2389
	var v2390 int32
	_ = v2390
	var v2391 int32
	_ = v2391
	var v2392 int32
	_ = v2392
	var v2432 int32
	_ = v2432
	var v2444 int32
	_ = v2444
	var v2445 int32
	_ = v2445
	var v2449 int32
	_ = v2449
	var v2450 int32
	_ = v2450
	var v2458 int32
	_ = v2458
	var v2459 int32
	_ = v2459
	var v2460 int32
	_ = v2460
	var v2461 int32
	_ = v2461
	var v2463 int32
	_ = v2463
	var v2464 int32
	_ = v2464
	var v2468 int32
	_ = v2468
	var v2469 int32
	_ = v2469
	var v2470 int32
	_ = v2470
	var v2471 int32
	_ = v2471
	var v2480 int32
	_ = v2480
	var v2481 int32
	_ = v2481
	var v2484 int64
	_ = v2484
	var v2487 int32
	_ = v2487
	var v2489 int64
	_ = v2489
	var v2492 int32
	_ = v2492
	var v2499 int32
	_ = v2499
	var v2501 int32
	_ = v2501
	var v2505 int32
	_ = v2505
	var v2506 int32
	_ = v2506
	var v2507 int32
	_ = v2507
	var v2508 int32
	_ = v2508
	var v2517 int32
	_ = v2517
	var v2518 int32
	_ = v2518
	var v2519 int32
	_ = v2519
	var v2521 int32
	_ = v2521
	var v2526 int32
	_ = v2526
	var v2527 int32
	_ = v2527
	var v2534 int32
	_ = v2534
	var v2537 int32
	_ = v2537
	var v2539 int32
	_ = v2539
	var v2543 int32
	_ = v2543
	var v2544 int32
	_ = v2544
	var v2545 int32
	_ = v2545
	var v2546 int32
	_ = v2546
	var v2586 int32
	_ = v2586
	var v2589 int32
	_ = v2589
	var v2595 int32
	_ = v2595
	var v2599 int32
	_ = v2599
	var v2601 int32
	_ = v2601
	var v2602 int32
	_ = v2602
	var v2603 int32
	_ = v2603
	var v2607 int32
	_ = v2607
	var v2608 int32
	_ = v2608
	var v2610 int32
	_ = v2610
	var v2622 int32
	_ = v2622
	var v2623 int32
	_ = v2623
	var v2627 int32
	_ = v2627
	var v2628 int32
	_ = v2628
	var v2636 int32
	_ = v2636
	var v2637 int32
	_ = v2637
	var v2638 int32
	_ = v2638
	var v2639 int32
	_ = v2639
	var v2641 int32
	_ = v2641
	var v2642 int32
	_ = v2642
	var v2646 int32
	_ = v2646
	var v2647 int32
	_ = v2647
	var v2648 int32
	_ = v2648
	var v2649 int32
	_ = v2649
	var v2658 int32
	_ = v2658
	var v2659 int32
	_ = v2659
	var v2662 int64
	_ = v2662
	var v2665 int32
	_ = v2665
	var v2667 int64
	_ = v2667
	var v2670 int32
	_ = v2670
	var v2677 int32
	_ = v2677
	var v2679 int32
	_ = v2679
	var v2683 int32
	_ = v2683
	var v2684 int32
	_ = v2684
	var v2685 int32
	_ = v2685
	var v2686 int32
	_ = v2686
	var v2695 int32
	_ = v2695
	var v2696 int32
	_ = v2696
	var v2697 int32
	_ = v2697
	var v2699 int32
	_ = v2699
	var v2704 int32
	_ = v2704
	var v2705 int32
	_ = v2705
	var v2712 int32
	_ = v2712
	var v2715 int32
	_ = v2715
	var v2717 int32
	_ = v2717
	var v2721 int32
	_ = v2721
	var v2722 int32
	_ = v2722
	var v2723 int32
	_ = v2723
	var v2724 int32
	_ = v2724
	var v2764 int32
	_ = v2764
	var v2767 int32
	_ = v2767
	var v2773 int32
	_ = v2773
	var v2785 int32
	_ = v2785
	var v2786 int32
	_ = v2786
	var v2790 int32
	_ = v2790
	var v2791 int32
	_ = v2791
	var v2799 int32
	_ = v2799
	var v2800 int32
	_ = v2800
	var v2801 int32
	_ = v2801
	var v2802 int32
	_ = v2802
	var v2804 int32
	_ = v2804
	var v2805 int32
	_ = v2805
	var v2809 int32
	_ = v2809
	var v2810 int32
	_ = v2810
	var v2811 int32
	_ = v2811
	var v2812 int32
	_ = v2812
	var v2821 int32
	_ = v2821
	var v2822 int32
	_ = v2822
	var v2825 int64
	_ = v2825
	var v2828 int32
	_ = v2828
	var v2830 int64
	_ = v2830
	var v2833 int32
	_ = v2833
	var v2840 int32
	_ = v2840
	var v2842 int32
	_ = v2842
	var v2846 int32
	_ = v2846
	var v2847 int32
	_ = v2847
	var v2848 int32
	_ = v2848
	var v2849 int32
	_ = v2849
	var v2858 int32
	_ = v2858
	var v2859 int32
	_ = v2859
	var v2860 int32
	_ = v2860
	var v2862 int32
	_ = v2862
	var v2867 int32
	_ = v2867
	var v2868 int32
	_ = v2868
	var v2875 int32
	_ = v2875
	var v2878 int32
	_ = v2878
	var v2880 int32
	_ = v2880
	var v2884 int32
	_ = v2884
	var v2885 int32
	_ = v2885
	var v2886 int32
	_ = v2886
	var v2887 int32
	_ = v2887
	var v2927 int32
	_ = v2927
	v13 = m.G0
	v15 = v13 - int32(224)
	m.G0 = v15
	*(*int64)(unsafe.Add(mBase, uint32(v15)+8)) = int64(0)
	v19 = *(*int32)(unsafe.Add(mBase, uint32(l3)))
	if v19 < int32(1) {
		goto L3
	} else {
		goto L4
	}
L1:
	;
	m.G0 = v15 + int32(224)
	return
L2:
	;
	if int32(2) < v51 {
		goto L41
	} else {
		goto L42
	}
L3:
	;
	v69 = int32(1)
	v70 = int32(4)
	goto L17
L4:
	;
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l3)+4))
	v23 = int32(0)
	v32 = v23
	v33 = v23
	goto L5
L5:
	;
	v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v22+v32))))
	if v38 == int32(0) {
		v51 = v33
		goto L7
	} else {
		goto L8
	}
L6:
	;
	if v51 != 0 {
		goto L2
	} else {
		goto L14
	}
L7:
	;
	v53 = v32 + int32(1)
	if v19 <= v53 {
		goto L11
	} else {
		goto L12
	}
L8:
	;
	if int32(1) < v33 {
		goto L9
	} else {
		goto L10
	}
L9:
	;
	v51 = v33 + int32(1)
	goto L7
L10:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v15+int32(8)+v33<<(uint(int32(2))%32)))) = v32
	goto L9
L11:
	;
	goto L6
L12:
	;
	if v51 < int32(3) {
		v32 = v53
		v33 = v51
		goto L5
	} else {
		goto L13
	}
L13:
	;
	goto L11
L14:
	;
	goto L3
L15:
	;
	goto L1
L16:
	;
	goto L15
L17:
	;
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v83+v70 < int32(32) {
		goto L20
	} else {
		goto L21
	}
L18:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v184 + v182
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v181<<(uint(v184)%32) | v183
	goto L16
L19:
	;
	v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v102 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v106 = v98
	v107 = v99
	v108 = v102
	v109 = v101
	goto L23
L20:
	;
	if v83 < int32(16) {
		v181 = v69
		v182 = v70
		v183 = v82
		v184 = v83
		goto L18
	} else {
		goto L22
	}
L21:
	;
	v87 = int32(32)
	v88 = v87 - v83
	v96 = int32(base.Ui32(v69) >> (uint(v88) % 32))
	v97 = v70 - v88
	v98 = v69<<(uint(v83)%32) | v82
	v99 = v87
	goto L19
L22:
	;
	v96 = v69
	v97 = v70
	v98 = v82
	v99 = v83
	goto L19
L23:
	;
	if base.Ui32(v108+int32(2)) <= base.Ui32(v109) {
		v164 = v108
		v165 = v109
		goto L25
	} else {
		goto L26
	}
L24:
	;
	v181 = v96
	v182 = v97
	v183 = v177
	v184 = v175
	goto L18
L25:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v164))) = uint16(v106)
	v172 = v164 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v172
	v175 = v107 + int32(-16)
	v177 = int32(base.Ui32(v106) >> (uint(int32(16)) % 32))
	if int32(31) < v107 {
		v106 = v177
		v107 = v175
		v108 = v172
		v109 = v165
		goto L23
	} else {
		goto L40
	}
L26:
	;
	v118 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v119 = v109 - v118
	v122 = base.I64_extend_i32_s(v119) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v122) {
		v146 = v118
		goto L28
	} else {
		goto L29
	}
L27:
	;
	if v108 == v118 {
		goto L38
	} else {
		goto L39
	}
L28:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v146
	goto L15
L29:
	;
	v125 = v108 - v118
	v127 = v122 + base.I64_extend_i32_u(v125)
	if base.Ui64(int64(4294967295)) < base.Ui64(v127) {
		v146 = v118
		goto L28
	} else {
		goto L30
	}
L30:
	;
	v130 = base.I32_wrap_i64(v127)
	if v109 == v118 {
		goto L31
	} else {
		goto L32
	}
L31:
	;
	v137 = int32(base.Ui32(v119*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v130) < base.Ui32(v137) {
		goto L34
	} else {
		goto L35
	}
L32:
	;
	if base.Ui32(v130) <= base.Ui32(v119) {
		v164 = v108
		v165 = v109
		goto L25
	} else {
		goto L33
	}
L33:
	;
	goto L31
L34:
	;
	v139 = v137
	goto L36
L35:
	;
	v139 = v130
	goto L36
L36:
	;
	v143 = v139&int32(-1024) + int32(1024)
	v144 = F_WebPSafeMalloc(m, int64(1), v143)
	mBase = m.M
	if v144 != 0 {
		goto L27
	} else {
		goto L37
	}
L37:
	;
	v145 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v146 = v145
	goto L28
L38:
	;
	v157 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v157)
	mBase = m.M
	v159 = v144 + v143
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v159
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v144
	v164 = v144 + v125
	v165 = v159
	goto L25
L39:
	;
	v155 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v156 = F_memcpy(m, v144, v155, v125)
	mBase = m.M
	goto L38
L40:
	;
	goto L24
L41:
	;
	v995 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v15)+207)) = v995
	v998 = base.Simd_g_const(&F_StoreHuffmanCode__k0)
	v999 = int32(192)
	base.Simd_g_v128_store(m, v15, v999, v998)
	*(*int64)(unsafe.Add(mBase, uint32(v15+int32(174)))) = int64(0)
	base.Simd_g_v128_store(m, v15+int32(160), v995, v998)
	v1009 = int32(144)
	base.Simd_g_v128_store(m, v15, v1009, v998)
	v1011 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(v15)+132)) = v1011
	*(*int32)(unsafe.Add(mBase, uint32(v15)+140)) = v15 + v1009
	*(*int32)(unsafe.Add(mBase, uint32(v15)+136)) = v15 + v999
	v1021 = int32(1)
	goto L181
L42:
	;
	v223 = *(*int32)(unsafe.Add(mBase, uint32(v15)+8))
	if int32(255) < v223 {
		goto L41
	} else {
		goto L43
	}
L43:
	;
	v226 = *(*int32)(unsafe.Add(mBase, uint32(v15)+12))
	if int32(255) < v226 {
		goto L41
	} else {
		goto L44
	}
L44:
	;
	v229 = int32(1)
	goto L47
L45:
	;
	v382 = v51 + int32(-1)
	v383 = int32(1)
	goto L73
L46:
	;
	goto L45
L47:
	;
	v242 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v243 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v243+v229 < int32(32) {
		goto L50
	} else {
		goto L51
	}
L48:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v344 + v342
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v341<<(uint(v344)%32) | v343
	goto L46
L49:
	;
	v261 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v262 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v266 = v258
	v267 = v259
	v268 = v262
	v269 = v261
	goto L53
L50:
	;
	if v243 < int32(16) {
		v341 = v229
		v342 = v229
		v343 = v242
		v344 = v243
		goto L48
	} else {
		goto L52
	}
L51:
	;
	v247 = int32(32)
	v248 = v247 - v243
	v256 = int32(base.Ui32(v229) >> (uint(v248) % 32))
	v257 = v229 - v248
	v258 = v229<<(uint(v243)%32) | v242
	v259 = v247
	goto L49
L52:
	;
	v256 = v229
	v257 = v229
	v258 = v242
	v259 = v243
	goto L49
L53:
	;
	if base.Ui32(v268+int32(2)) <= base.Ui32(v269) {
		v324 = v268
		v325 = v269
		goto L55
	} else {
		goto L56
	}
L54:
	;
	v341 = v256
	v342 = v257
	v343 = v337
	v344 = v335
	goto L48
L55:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v324))) = uint16(v266)
	v332 = v324 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v332
	v335 = v267 + int32(-16)
	v337 = int32(base.Ui32(v266) >> (uint(int32(16)) % 32))
	if int32(31) < v267 {
		v266 = v337
		v267 = v335
		v268 = v332
		v269 = v325
		goto L53
	} else {
		goto L70
	}
L56:
	;
	v278 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v279 = v269 - v278
	v282 = base.I64_extend_i32_s(v279) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v282) {
		v306 = v278
		goto L58
	} else {
		goto L59
	}
L57:
	;
	if v268 == v278 {
		goto L68
	} else {
		goto L69
	}
L58:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v306
	goto L45
L59:
	;
	v285 = v268 - v278
	v287 = v282 + base.I64_extend_i32_u(v285)
	if base.Ui64(int64(4294967295)) < base.Ui64(v287) {
		v306 = v278
		goto L58
	} else {
		goto L60
	}
L60:
	;
	v290 = base.I32_wrap_i64(v287)
	if v269 == v278 {
		goto L61
	} else {
		goto L62
	}
L61:
	;
	v297 = int32(base.Ui32(v279*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v290) < base.Ui32(v297) {
		goto L64
	} else {
		goto L65
	}
L62:
	;
	if base.Ui32(v290) <= base.Ui32(v279) {
		v324 = v268
		v325 = v269
		goto L55
	} else {
		goto L63
	}
L63:
	;
	goto L61
L64:
	;
	v299 = v297
	goto L66
L65:
	;
	v299 = v290
	goto L66
L66:
	;
	v303 = v299&int32(-1024) + int32(1024)
	v304 = F_WebPSafeMalloc(m, int64(1), v303)
	mBase = m.M
	if v304 != 0 {
		goto L57
	} else {
		goto L67
	}
L67:
	;
	v305 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v306 = v305
	goto L58
L68:
	;
	v317 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v317)
	mBase = m.M
	v319 = v304 + v303
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v319
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v304
	v324 = v304 + v285
	v325 = v319
	goto L55
L69:
	;
	v315 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v316 = F_memcpy(m, v304, v315, v285)
	mBase = m.M
	goto L68
L70:
	;
	goto L54
L71:
	;
	v534 = int32(1)
	v535 = base.B2i32(v534 < v223)
	goto L99
L72:
	;
	goto L71
L73:
	;
	v395 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v396 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v396+v383 < int32(32) {
		goto L76
	} else {
		goto L77
	}
L74:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v497 + v495
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v494<<(uint(v497)%32) | v496
	goto L72
L75:
	;
	v414 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v415 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v419 = v411
	v420 = v412
	v421 = v415
	v422 = v414
	goto L79
L76:
	;
	if v396 < int32(16) {
		v494 = v382
		v495 = v383
		v496 = v395
		v497 = v396
		goto L74
	} else {
		goto L78
	}
L77:
	;
	v400 = int32(32)
	v401 = v400 - v396
	v409 = int32(base.Ui32(v382) >> (uint(v401) % 32))
	v410 = v383 - v401
	v411 = v382<<(uint(v396)%32) | v395
	v412 = v400
	goto L75
L78:
	;
	v409 = v382
	v410 = v383
	v411 = v395
	v412 = v396
	goto L75
L79:
	;
	if base.Ui32(v421+int32(2)) <= base.Ui32(v422) {
		v477 = v421
		v478 = v422
		goto L81
	} else {
		goto L82
	}
L80:
	;
	v494 = v409
	v495 = v410
	v496 = v490
	v497 = v488
	goto L74
L81:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v477))) = uint16(v419)
	v485 = v477 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v485
	v488 = v420 + int32(-16)
	v490 = int32(base.Ui32(v419) >> (uint(int32(16)) % 32))
	if int32(31) < v420 {
		v419 = v490
		v420 = v488
		v421 = v485
		v422 = v478
		goto L79
	} else {
		goto L96
	}
L82:
	;
	v431 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v432 = v422 - v431
	v435 = base.I64_extend_i32_s(v432) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v435) {
		v459 = v431
		goto L84
	} else {
		goto L85
	}
L83:
	;
	if v421 == v431 {
		goto L94
	} else {
		goto L95
	}
L84:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v459
	goto L71
L85:
	;
	v438 = v421 - v431
	v440 = v435 + base.I64_extend_i32_u(v438)
	if base.Ui64(int64(4294967295)) < base.Ui64(v440) {
		v459 = v431
		goto L84
	} else {
		goto L86
	}
L86:
	;
	v443 = base.I32_wrap_i64(v440)
	if v422 == v431 {
		goto L87
	} else {
		goto L88
	}
L87:
	;
	v450 = int32(base.Ui32(v432*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v443) < base.Ui32(v450) {
		goto L90
	} else {
		goto L91
	}
L88:
	;
	if base.Ui32(v443) <= base.Ui32(v432) {
		v477 = v421
		v478 = v422
		goto L81
	} else {
		goto L89
	}
L89:
	;
	goto L87
L90:
	;
	v452 = v450
	goto L92
L91:
	;
	v452 = v443
	goto L92
L92:
	;
	v456 = v452&int32(-1024) + int32(1024)
	v457 = F_WebPSafeMalloc(m, int64(1), v456)
	mBase = m.M
	if v457 != 0 {
		goto L83
	} else {
		goto L93
	}
L93:
	;
	v458 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v459 = v458
	goto L84
L94:
	;
	v470 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v470)
	mBase = m.M
	v472 = v457 + v456
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v472
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v457
	v477 = v457 + v438
	v478 = v472
	goto L81
L95:
	;
	v468 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v469 = F_memcpy(m, v457, v468, v438)
	mBase = m.M
	goto L94
L96:
	;
	goto L80
L97:
	;
	if v534 < v223 {
		goto L123
	} else {
		goto L124
	}
L98:
	;
	goto L97
L99:
	;
	v548 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v549 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v549+v534 < int32(32) {
		goto L102
	} else {
		goto L103
	}
L100:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v650 + v648
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v647<<(uint(v650)%32) | v649
	goto L98
L101:
	;
	v567 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v568 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v572 = v564
	v573 = v565
	v574 = v568
	v575 = v567
	goto L105
L102:
	;
	if v549 < int32(16) {
		v647 = v535
		v648 = v534
		v649 = v548
		v650 = v549
		goto L100
	} else {
		goto L104
	}
L103:
	;
	v553 = int32(32)
	v554 = v553 - v549
	v562 = int32(base.Ui32(v535) >> (uint(v554) % 32))
	v563 = v534 - v554
	v564 = v535<<(uint(v549)%32) | v548
	v565 = v553
	goto L101
L104:
	;
	v562 = v535
	v563 = v534
	v564 = v548
	v565 = v549
	goto L101
L105:
	;
	if base.Ui32(v574+int32(2)) <= base.Ui32(v575) {
		v630 = v574
		v631 = v575
		goto L107
	} else {
		goto L108
	}
L106:
	;
	v647 = v562
	v648 = v563
	v649 = v643
	v650 = v641
	goto L100
L107:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v630))) = uint16(v572)
	v638 = v630 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v638
	v641 = v573 + int32(-16)
	v643 = int32(base.Ui32(v572) >> (uint(int32(16)) % 32))
	if int32(31) < v573 {
		v572 = v643
		v573 = v641
		v574 = v638
		v575 = v631
		goto L105
	} else {
		goto L122
	}
L108:
	;
	v584 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v585 = v575 - v584
	v588 = base.I64_extend_i32_s(v585) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v588) {
		v612 = v584
		goto L110
	} else {
		goto L111
	}
L109:
	;
	if v574 == v584 {
		goto L120
	} else {
		goto L121
	}
L110:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v612
	goto L97
L111:
	;
	v591 = v574 - v584
	v593 = v588 + base.I64_extend_i32_u(v591)
	if base.Ui64(int64(4294967295)) < base.Ui64(v593) {
		v612 = v584
		goto L110
	} else {
		goto L112
	}
L112:
	;
	v596 = base.I32_wrap_i64(v593)
	if v575 == v584 {
		goto L113
	} else {
		goto L114
	}
L113:
	;
	v603 = int32(base.Ui32(v585*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v596) < base.Ui32(v603) {
		goto L116
	} else {
		goto L117
	}
L114:
	;
	if base.Ui32(v596) <= base.Ui32(v585) {
		v630 = v574
		v631 = v575
		goto L107
	} else {
		goto L115
	}
L115:
	;
	goto L113
L116:
	;
	v605 = v603
	goto L118
L117:
	;
	v605 = v596
	goto L118
L118:
	;
	v609 = v605&int32(-1024) + int32(1024)
	v610 = F_WebPSafeMalloc(m, int64(1), v609)
	mBase = m.M
	if v610 != 0 {
		goto L109
	} else {
		goto L119
	}
L119:
	;
	v611 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v612 = v611
	goto L110
L120:
	;
	v623 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v623)
	mBase = m.M
	v625 = v610 + v609
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v625
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v610
	v630 = v610 + v591
	v631 = v625
	goto L107
L121:
	;
	v621 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v622 = F_memcpy(m, v610, v621, v591)
	mBase = m.M
	goto L120
L122:
	;
	goto L106
L123:
	;
	v689 = int32(8)
	goto L125
L124:
	;
	v689 = int32(1)
	goto L125
L125:
	;
	if v689 < int32(1) {
		goto L127
	} else {
		goto L128
	}
L126:
	;
	if v51 != int32(2) {
		goto L1
	} else {
		goto L152
	}
L127:
	;
	goto L126
L128:
	;
	v701 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v702 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v702+v689 < int32(32) {
		goto L131
	} else {
		goto L132
	}
L129:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v803 + v801
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v800<<(uint(v803)%32) | v802
	goto L127
L130:
	;
	v720 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v721 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v725 = v717
	v726 = v718
	v727 = v721
	v728 = v720
	goto L134
L131:
	;
	if v702 < int32(16) {
		v800 = v223
		v801 = v689
		v802 = v701
		v803 = v702
		goto L129
	} else {
		goto L133
	}
L132:
	;
	v706 = int32(32)
	v707 = v706 - v702
	v715 = int32(base.Ui32(v223) >> (uint(v707) % 32))
	v716 = v689 - v707
	v717 = v223<<(uint(v702)%32) | v701
	v718 = v706
	goto L130
L133:
	;
	v715 = v223
	v716 = v689
	v717 = v701
	v718 = v702
	goto L130
L134:
	;
	if base.Ui32(v727+int32(2)) <= base.Ui32(v728) {
		v783 = v727
		v784 = v728
		goto L136
	} else {
		goto L137
	}
L135:
	;
	v800 = v715
	v801 = v716
	v802 = v796
	v803 = v794
	goto L129
L136:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v783))) = uint16(v725)
	v791 = v783 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v791
	v794 = v726 + int32(-16)
	v796 = int32(base.Ui32(v725) >> (uint(int32(16)) % 32))
	if int32(31) < v726 {
		v725 = v796
		v726 = v794
		v727 = v791
		v728 = v784
		goto L134
	} else {
		goto L151
	}
L137:
	;
	v737 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v738 = v728 - v737
	v741 = base.I64_extend_i32_s(v738) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v741) {
		v765 = v737
		goto L139
	} else {
		goto L140
	}
L138:
	;
	if v727 == v737 {
		goto L149
	} else {
		goto L150
	}
L139:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v765
	goto L126
L140:
	;
	v744 = v727 - v737
	v746 = v741 + base.I64_extend_i32_u(v744)
	if base.Ui64(int64(4294967295)) < base.Ui64(v746) {
		v765 = v737
		goto L139
	} else {
		goto L141
	}
L141:
	;
	v749 = base.I32_wrap_i64(v746)
	if v728 == v737 {
		goto L142
	} else {
		goto L143
	}
L142:
	;
	v756 = int32(base.Ui32(v738*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v749) < base.Ui32(v756) {
		goto L145
	} else {
		goto L146
	}
L143:
	;
	if base.Ui32(v749) <= base.Ui32(v738) {
		v783 = v727
		v784 = v728
		goto L136
	} else {
		goto L144
	}
L144:
	;
	goto L142
L145:
	;
	v758 = v756
	goto L147
L146:
	;
	v758 = v749
	goto L147
L147:
	;
	v762 = v758&int32(-1024) + int32(1024)
	v763 = F_WebPSafeMalloc(m, int64(1), v762)
	mBase = m.M
	if v763 != 0 {
		goto L138
	} else {
		goto L148
	}
L148:
	;
	v764 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v765 = v764
	goto L139
L149:
	;
	v776 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v776)
	mBase = m.M
	v778 = v763 + v762
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v778
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v763
	v783 = v763 + v744
	v784 = v778
	goto L136
L150:
	;
	v774 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v775 = F_memcpy(m, v763, v774, v744)
	mBase = m.M
	goto L149
L151:
	;
	goto L135
L152:
	;
	v842 = int32(8)
	goto L155
L153:
	;
	goto L1
L154:
	;
	goto L153
L155:
	;
	v854 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v855 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v855+v842 < int32(32) {
		goto L158
	} else {
		goto L159
	}
L156:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v956 + v954
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v953<<(uint(v956)%32) | v955
	goto L154
L157:
	;
	v873 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v874 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v878 = v870
	v879 = v871
	v880 = v874
	v881 = v873
	goto L161
L158:
	;
	if v855 < int32(16) {
		v953 = v226
		v954 = v842
		v955 = v854
		v956 = v855
		goto L156
	} else {
		goto L160
	}
L159:
	;
	v859 = int32(32)
	v860 = v859 - v855
	v868 = int32(base.Ui32(v226) >> (uint(v860) % 32))
	v869 = v842 - v860
	v870 = v226<<(uint(v855)%32) | v854
	v871 = v859
	goto L157
L160:
	;
	v868 = v226
	v869 = v842
	v870 = v854
	v871 = v855
	goto L157
L161:
	;
	if base.Ui32(v880+int32(2)) <= base.Ui32(v881) {
		v936 = v880
		v937 = v881
		goto L163
	} else {
		goto L164
	}
L162:
	;
	v953 = v868
	v954 = v869
	v955 = v949
	v956 = v947
	goto L156
L163:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v936))) = uint16(v878)
	v944 = v936 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v944
	v947 = v879 + int32(-16)
	v949 = int32(base.Ui32(v878) >> (uint(int32(16)) % 32))
	if int32(31) < v879 {
		v878 = v949
		v879 = v947
		v880 = v944
		v881 = v937
		goto L161
	} else {
		goto L178
	}
L164:
	;
	v890 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v891 = v881 - v890
	v894 = base.I64_extend_i32_s(v891) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v894) {
		v918 = v890
		goto L166
	} else {
		goto L167
	}
L165:
	;
	if v880 == v890 {
		goto L176
	} else {
		goto L177
	}
L166:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v918
	goto L153
L167:
	;
	v897 = v880 - v890
	v899 = v894 + base.I64_extend_i32_u(v897)
	if base.Ui64(int64(4294967295)) < base.Ui64(v899) {
		v918 = v890
		goto L166
	} else {
		goto L168
	}
L168:
	;
	v902 = base.I32_wrap_i64(v899)
	if v881 == v890 {
		goto L169
	} else {
		goto L170
	}
L169:
	;
	v909 = int32(base.Ui32(v891*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v902) < base.Ui32(v909) {
		goto L172
	} else {
		goto L173
	}
L170:
	;
	if base.Ui32(v902) <= base.Ui32(v891) {
		v936 = v880
		v937 = v881
		goto L163
	} else {
		goto L171
	}
L171:
	;
	goto L169
L172:
	;
	v911 = v909
	goto L174
L173:
	;
	v911 = v902
	goto L174
L174:
	;
	v915 = v911&int32(-1024) + int32(1024)
	v916 = F_WebPSafeMalloc(m, int64(1), v915)
	mBase = m.M
	if v916 != 0 {
		goto L165
	} else {
		goto L175
	}
L175:
	;
	v917 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v918 = v917
	goto L166
L176:
	;
	v929 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v929)
	mBase = m.M
	v931 = v916 + v915
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v931
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v916
	v936 = v916 + v897
	v937 = v931
	goto L163
L177:
	;
	v927 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v928 = F_memcpy(m, v916, v927, v897)
	mBase = m.M
	goto L176
L178:
	;
	goto L162
L179:
	;
	v1172 = F_VP8LCreateCompressedHuffmanTree(m, l3, l2, v19)
	mBase = m.M
	goto L207
L180:
	;
	goto L179
L181:
	;
	v1033 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1034 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1034+v1021 < int32(32) {
		goto L184
	} else {
		goto L185
	}
L182:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1135 + v1133
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1132<<(uint(v1135)%32) | v1134
	goto L180
L183:
	;
	v1052 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1053 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1057 = v1049
	v1058 = v1050
	v1059 = v1053
	v1060 = v1052
	goto L187
L184:
	;
	if v1034 < int32(16) {
		v1132 = v995
		v1133 = v1021
		v1134 = v1033
		v1135 = v1034
		goto L182
	} else {
		goto L186
	}
L185:
	;
	v1038 = int32(32)
	v1039 = v1038 - v1034
	v1047 = int32(base.Ui32(v995) >> (uint(v1039) % 32))
	v1048 = v1021 - v1039
	v1049 = v995<<(uint(v1034)%32) | v1033
	v1050 = v1038
	goto L183
L186:
	;
	v1047 = v995
	v1048 = v1021
	v1049 = v1033
	v1050 = v1034
	goto L183
L187:
	;
	if base.Ui32(v1059+int32(2)) <= base.Ui32(v1060) {
		v1115 = v1059
		v1116 = v1060
		goto L189
	} else {
		goto L190
	}
L188:
	;
	v1132 = v1047
	v1133 = v1048
	v1134 = v1128
	v1135 = v1126
	goto L182
L189:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1115))) = uint16(v1057)
	v1123 = v1115 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1123
	v1126 = v1058 + int32(-16)
	v1128 = int32(base.Ui32(v1057) >> (uint(int32(16)) % 32))
	if int32(31) < v1058 {
		v1057 = v1128
		v1058 = v1126
		v1059 = v1123
		v1060 = v1116
		goto L187
	} else {
		goto L204
	}
L190:
	;
	v1069 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1070 = v1060 - v1069
	v1073 = base.I64_extend_i32_s(v1070) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1073) {
		v1097 = v1069
		goto L192
	} else {
		goto L193
	}
L191:
	;
	if v1059 == v1069 {
		goto L202
	} else {
		goto L203
	}
L192:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1097
	goto L179
L193:
	;
	v1076 = v1059 - v1069
	v1078 = v1073 + base.I64_extend_i32_u(v1076)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1078) {
		v1097 = v1069
		goto L192
	} else {
		goto L194
	}
L194:
	;
	v1081 = base.I32_wrap_i64(v1078)
	if v1060 == v1069 {
		goto L195
	} else {
		goto L196
	}
L195:
	;
	v1088 = int32(base.Ui32(v1070*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1081) < base.Ui32(v1088) {
		goto L198
	} else {
		goto L199
	}
L196:
	;
	if base.Ui32(v1081) <= base.Ui32(v1070) {
		v1115 = v1059
		v1116 = v1060
		goto L189
	} else {
		goto L197
	}
L197:
	;
	goto L195
L198:
	;
	v1090 = v1088
	goto L200
L199:
	;
	v1090 = v1081
	goto L200
L200:
	;
	v1094 = v1090&int32(-1024) + int32(1024)
	v1095 = F_WebPSafeMalloc(m, int64(1), v1094)
	mBase = m.M
	if v1095 != 0 {
		goto L191
	} else {
		goto L201
	}
L201:
	;
	v1096 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1097 = v1096
	goto L192
L202:
	;
	v1108 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1108)
	mBase = m.M
	v1110 = v1095 + v1094
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1110
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1095
	v1115 = v1095 + v1076
	v1116 = v1110
	goto L189
L203:
	;
	v1106 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1107 = F_memcpy(m, v1095, v1106, v1076)
	mBase = m.M
	goto L202
L204:
	;
	goto L188
L205:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v15)+31)) = int32(0)
	base.Simd_g_v128_store(m, v15, int32(16), v998)
	if v1172 < int32(1) {
		goto L219
	} else {
		goto L220
	}
L207:
	;
	base.MemoryFill(m, v15+int32(48), int32(0), int32(76))
	goto L205
L219:
	;
	F_VP8LCreateHuffmanTree(m, v15+int32(48), int32(7), v15+int32(16), l1, v15+int32(132))
	mBase = m.M
	v1436 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+207)))
	if v1436 != 0 {
		v1493 = v1011
		goto L230
	} else {
		goto L231
	}
L220:
	;
	v1304 = v1172 & int32(3)
	if base.Ui32(v1172) < base.Ui32(int32(4)) {
		v1382 = v995
		goto L221
	} else {
		goto L222
	}
L221:
	;
	if v1304 == int32(0) {
		goto L219
	} else {
		goto L226
	}
L222:
	;
	v1317 = l2
	v1319 = int32(0)
	goto L223
L223:
	;
	v1323 = v15 + int32(48)
	v1324 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1317))))
	v1325 = int32(2)
	v1327 = v1323 + v1324<<(uint(v1325)%32)
	v1328 = *(*int32)(unsafe.Add(mBase, uint32(v1327)))
	v1329 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v1327))) = v1328 + v1329
	v1336 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1317+v1325))))
	v1339 = v1323 + v1336<<(uint(v1325)%32)
	v1340 = *(*int32)(unsafe.Add(mBase, uint32(v1339)))
	*(*int32)(unsafe.Add(mBase, uint32(v1339))) = v1340 + v1329
	v1346 = int32(4)
	v1348 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1317+v1346))))
	v1351 = v1323 + v1348<<(uint(v1325)%32)
	v1352 = *(*int32)(unsafe.Add(mBase, uint32(v1351)))
	*(*int32)(unsafe.Add(mBase, uint32(v1351))) = v1352 + v1329
	v1360 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1317+int32(6)))))
	v1363 = v1323 + v1360<<(uint(v1325)%32)
	v1364 = *(*int32)(unsafe.Add(mBase, uint32(v1363)))
	*(*int32)(unsafe.Add(mBase, uint32(v1363))) = v1364 + v1329
	v1371 = v1319 + v1346
	if v1172&int32(2147483644) != v1371 {
		v1317 = v1317 + int32(8)
		v1319 = v1371
		goto L223
	} else {
		goto L225
	}
L224:
	;
	v1382 = v1371
	goto L221
L225:
	;
	goto L224
L226:
	;
	v1395 = v1304
	v1397 = l2 + v1382<<(uint(int32(1))%32)
	goto L227
L227:
	;
	v1404 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1397))))
	v1405 = int32(2)
	v1407 = v15 + int32(48) + v1404<<(uint(v1405)%32)
	v1408 = *(*int32)(unsafe.Add(mBase, uint32(v1407)))
	*(*int32)(unsafe.Add(mBase, uint32(v1407))) = v1408 + int32(1)
	v1415 = v1395 + int32(-1)
	if v1415 != 0 {
		v1395 = v1415
		v1397 = v1397 + v1405
		goto L227
	} else {
		goto L229
	}
L228:
	;
	goto L219
L229:
	;
	goto L228
L230:
	;
	v1495 = v1493 + int32(-4)
	v1496 = int32(4)
	goto L263
L231:
	;
	v1437 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+206)))
	if v1437 == int32(0) {
		goto L232
	} else {
		goto L233
	}
L232:
	;
	v1441 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+205)))
	if v1441 == int32(0) {
		goto L234
	} else {
		goto L235
	}
L233:
	;
	v1493 = int32(18)
	goto L230
L234:
	;
	v1445 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+204)))
	if v1445 == int32(0) {
		goto L236
	} else {
		goto L237
	}
L235:
	;
	v1493 = int32(17)
	goto L230
L236:
	;
	v1449 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+203)))
	if v1449 == int32(0) {
		goto L238
	} else {
		goto L239
	}
L237:
	;
	v1493 = int32(16)
	goto L230
L238:
	;
	v1453 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+202)))
	if v1453 == int32(0) {
		goto L240
	} else {
		goto L241
	}
L239:
	;
	v1493 = int32(15)
	goto L230
L240:
	;
	v1457 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+201)))
	if v1457 == int32(0) {
		goto L242
	} else {
		goto L243
	}
L241:
	;
	v1493 = int32(14)
	goto L230
L242:
	;
	v1461 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+200)))
	if v1461 == int32(0) {
		goto L244
	} else {
		goto L245
	}
L243:
	;
	v1493 = int32(13)
	goto L230
L244:
	;
	v1465 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+199)))
	if v1465 == int32(0) {
		goto L246
	} else {
		goto L247
	}
L245:
	;
	v1493 = int32(12)
	goto L230
L246:
	;
	v1469 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+198)))
	if v1469 == int32(0) {
		goto L248
	} else {
		goto L249
	}
L247:
	;
	v1493 = int32(11)
	goto L230
L248:
	;
	v1473 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+208)))
	if v1473 == int32(0) {
		goto L250
	} else {
		goto L251
	}
L249:
	;
	v1493 = int32(10)
	goto L230
L250:
	;
	v1477 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+197)))
	if v1477 == int32(0) {
		goto L252
	} else {
		goto L253
	}
L251:
	;
	v1493 = int32(9)
	goto L230
L252:
	;
	v1481 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+196)))
	if v1481 == int32(0) {
		goto L254
	} else {
		goto L255
	}
L253:
	;
	v1493 = int32(8)
	goto L230
L254:
	;
	v1485 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+195)))
	if v1485 == int32(0) {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	v1493 = int32(7)
	goto L230
L256:
	;
	v1491 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15)+194)))
	if v1491 != 0 {
		goto L258
	} else {
		goto L259
	}
L257:
	;
	v1493 = int32(6)
	goto L230
L258:
	;
	v1492 = int32(5)
	goto L260
L259:
	;
	v1492 = int32(4)
	goto L260
L260:
	;
	v1493 = v1492
	goto L230
L261:
	;
	v1655 = int32(0)
	goto L287
L262:
	;
	goto L261
L263:
	;
	v1508 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1509 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1509+v1496 < int32(32) {
		goto L266
	} else {
		goto L267
	}
L264:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1610 + v1608
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1607<<(uint(v1610)%32) | v1609
	goto L262
L265:
	;
	v1527 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1528 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1532 = v1524
	v1533 = v1525
	v1534 = v1528
	v1535 = v1527
	goto L269
L266:
	;
	if v1509 < int32(16) {
		v1607 = v1495
		v1608 = v1496
		v1609 = v1508
		v1610 = v1509
		goto L264
	} else {
		goto L268
	}
L267:
	;
	v1513 = int32(32)
	v1514 = v1513 - v1509
	v1522 = int32(base.Ui32(v1495) >> (uint(v1514) % 32))
	v1523 = v1496 - v1514
	v1524 = v1495<<(uint(v1509)%32) | v1508
	v1525 = v1513
	goto L265
L268:
	;
	v1522 = v1495
	v1523 = v1496
	v1524 = v1508
	v1525 = v1509
	goto L265
L269:
	;
	if base.Ui32(v1534+int32(2)) <= base.Ui32(v1535) {
		v1590 = v1534
		v1591 = v1535
		goto L271
	} else {
		goto L272
	}
L270:
	;
	v1607 = v1522
	v1608 = v1523
	v1609 = v1603
	v1610 = v1601
	goto L264
L271:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1590))) = uint16(v1532)
	v1598 = v1590 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1598
	v1601 = v1533 + int32(-16)
	v1603 = int32(base.Ui32(v1532) >> (uint(int32(16)) % 32))
	if int32(31) < v1533 {
		v1532 = v1603
		v1533 = v1601
		v1534 = v1598
		v1535 = v1591
		goto L269
	} else {
		goto L286
	}
L272:
	;
	v1544 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1545 = v1535 - v1544
	v1548 = base.I64_extend_i32_s(v1545) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1548) {
		v1572 = v1544
		goto L274
	} else {
		goto L275
	}
L273:
	;
	if v1534 == v1544 {
		goto L284
	} else {
		goto L285
	}
L274:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1572
	goto L261
L275:
	;
	v1551 = v1534 - v1544
	v1553 = v1548 + base.I64_extend_i32_u(v1551)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1553) {
		v1572 = v1544
		goto L274
	} else {
		goto L276
	}
L276:
	;
	v1556 = base.I32_wrap_i64(v1553)
	if v1535 == v1544 {
		goto L277
	} else {
		goto L278
	}
L277:
	;
	v1563 = int32(base.Ui32(v1545*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1556) < base.Ui32(v1563) {
		goto L280
	} else {
		goto L281
	}
L278:
	;
	if base.Ui32(v1556) <= base.Ui32(v1545) {
		v1590 = v1534
		v1591 = v1535
		goto L271
	} else {
		goto L279
	}
L279:
	;
	goto L277
L280:
	;
	v1565 = v1563
	goto L282
L281:
	;
	v1565 = v1556
	goto L282
L282:
	;
	v1569 = v1565&int32(-1024) + int32(1024)
	v1570 = F_WebPSafeMalloc(m, int64(1), v1569)
	mBase = m.M
	if v1570 != 0 {
		goto L273
	} else {
		goto L283
	}
L283:
	;
	v1571 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1572 = v1571
	goto L274
L284:
	;
	v1583 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1583)
	mBase = m.M
	v1585 = v1570 + v1569
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1585
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1570
	v1590 = v1570 + v1551
	v1591 = v1585
	goto L271
L285:
	;
	v1581 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1582 = F_memcpy(m, v1570, v1581, v1551)
	mBase = m.M
	goto L284
L286:
	;
	goto L270
L287:
	;
	v1662 = m.G1
	v1666 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1662+int32(_a_F_StoreHuffmanCode_0)+v1655))))
	v1668 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15+int32(192)+v1666))))
	v1669 = int32(3)
	goto L291
L288:
	;
	v1823 = *(*int32)(unsafe.Add(mBase, uint32(v15)+132))
	if v1823 < int32(1) {
		goto L316
	} else {
		goto L317
	}
L289:
	;
	v1821 = v1655 + int32(1)
	if v1493 != v1821 {
		v1655 = v1821
		goto L287
	} else {
		goto L315
	}
L290:
	;
	goto L289
L291:
	;
	v1681 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1682 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1682+v1669 < int32(32) {
		goto L294
	} else {
		goto L295
	}
L292:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1783 + v1781
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1780<<(uint(v1783)%32) | v1782
	goto L290
L293:
	;
	v1700 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1701 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1705 = v1697
	v1706 = v1698
	v1707 = v1701
	v1708 = v1700
	goto L297
L294:
	;
	if v1682 < int32(16) {
		v1780 = v1668
		v1781 = v1669
		v1782 = v1681
		v1783 = v1682
		goto L292
	} else {
		goto L296
	}
L295:
	;
	v1686 = int32(32)
	v1687 = v1686 - v1682
	v1695 = int32(base.Ui32(v1668) >> (uint(v1687) % 32))
	v1696 = v1669 - v1687
	v1697 = v1668<<(uint(v1682)%32) | v1681
	v1698 = v1686
	goto L293
L296:
	;
	v1695 = v1668
	v1696 = v1669
	v1697 = v1681
	v1698 = v1682
	goto L293
L297:
	;
	if base.Ui32(v1707+int32(2)) <= base.Ui32(v1708) {
		v1763 = v1707
		v1764 = v1708
		goto L299
	} else {
		goto L300
	}
L298:
	;
	v1780 = v1695
	v1781 = v1696
	v1782 = v1776
	v1783 = v1774
	goto L292
L299:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1763))) = uint16(v1705)
	v1771 = v1763 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1771
	v1774 = v1706 + int32(-16)
	v1776 = int32(base.Ui32(v1705) >> (uint(int32(16)) % 32))
	if int32(31) < v1706 {
		v1705 = v1776
		v1706 = v1774
		v1707 = v1771
		v1708 = v1764
		goto L297
	} else {
		goto L314
	}
L300:
	;
	v1717 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1718 = v1708 - v1717
	v1721 = base.I64_extend_i32_s(v1718) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1721) {
		v1745 = v1717
		goto L302
	} else {
		goto L303
	}
L301:
	;
	if v1707 == v1717 {
		goto L312
	} else {
		goto L313
	}
L302:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1745
	goto L289
L303:
	;
	v1724 = v1707 - v1717
	v1726 = v1721 + base.I64_extend_i32_u(v1724)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1726) {
		v1745 = v1717
		goto L302
	} else {
		goto L304
	}
L304:
	;
	v1729 = base.I32_wrap_i64(v1726)
	if v1708 == v1717 {
		goto L305
	} else {
		goto L306
	}
L305:
	;
	v1736 = int32(base.Ui32(v1718*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1729) < base.Ui32(v1736) {
		goto L308
	} else {
		goto L309
	}
L306:
	;
	if base.Ui32(v1729) <= base.Ui32(v1718) {
		v1763 = v1707
		v1764 = v1708
		goto L299
	} else {
		goto L307
	}
L307:
	;
	goto L305
L308:
	;
	v1738 = v1736
	goto L310
L309:
	;
	v1738 = v1729
	goto L310
L310:
	;
	v1742 = v1738&int32(-1024) + int32(1024)
	v1743 = F_WebPSafeMalloc(m, int64(1), v1742)
	mBase = m.M
	if v1743 != 0 {
		goto L301
	} else {
		goto L311
	}
L311:
	;
	v1744 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1745 = v1744
	goto L302
L312:
	;
	v1756 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v1756)
	mBase = m.M
	v1758 = v1743 + v1742
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v1758
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v1743
	v1763 = v1743 + v1724
	v1764 = v1758
	goto L299
L313:
	;
	v1754 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1755 = F_memcpy(m, v1743, v1754, v1724)
	mBase = m.M
	goto L312
L314:
	;
	goto L298
L315:
	;
	goto L288
L316:
	;
	if int32(1) <= v1172 {
		goto L328
	} else {
		goto L329
	}
L317:
	;
	v1826 = *(*int32)(unsafe.Add(mBase, uint32(v15)+136))
	v1833 = v1823
	v1834 = int32(0)
	v1835 = v1826
	goto L318
L318:
	;
	v1840 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1835))))
	if v1840 == int32(0) {
		v1846 = v1834
		goto L320
	} else {
		goto L321
	}
L319:
	;
	v1852 = int32(0)
	v1859 = v1852
	v1861 = v1852
	goto L324
L320:
	;
	v1851 = v1833 + int32(-1)
	if v1851 != 0 {
		v1833 = v1851
		v1834 = v1846
		v1835 = v1835 + int32(1)
		goto L318
	} else {
		goto L323
	}
L321:
	;
	if int32(0) < v1834 {
		goto L316
	} else {
		goto L322
	}
L322:
	;
	v1846 = int32(1)
	goto L320
L323:
	;
	goto L319
L324:
	;
	v1866 = *(*int32)(unsafe.Add(mBase, uint32(v15)+136))
	v1868 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v1866+v1859))) = uint8(v1868)
	v1870 = *(*int32)(unsafe.Add(mBase, uint32(v15)+140))
	*(*uint16)(unsafe.Add(mBase, uint32(v1870+v1861))) = uint16(v1868)
	v1877 = v1859 + int32(1)
	v1878 = *(*int32)(unsafe.Add(mBase, uint32(v15)+132))
	if v1877 < v1878 {
		v1859 = v1877
		v1861 = v1861 + int32(2)
		goto L324
	} else {
		goto L326
	}
L325:
	;
	goto L316
L326:
	;
	goto L325
L327:
	;
	v1960 = int32(1)
	v1962 = base.B2i32(v1960 < v1957) & v1955
	goto L341
L328:
	;
	v1906 = int32(0)
	v1907 = v1172<<(uint(int32(1))%32) + l2 + int32(-2)
	v1910 = v1172
	goto L331
L329:
	;
	v1955 = int32(0)
	v1957 = v1172
	goto L327
L330:
	;
	v1955 = base.B2i32(int32(12) < v1943)
	v1957 = v1945
	goto L327
L331:
	;
	v1913 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1907))))
	if base.Ui32(int32(18)) < base.Ui32(v1913) {
		v1943 = v1906
		v1945 = v1910
		goto L330
	} else {
		goto L333
	}
L332:
	;
	v1943 = v1933
	v1945 = int32(0)
	goto L330
L333:
	;
	if int32(1)<<(uint(v1913)%32)&int32(393217) == int32(0) {
		v1943 = v1906
		v1945 = v1910
		goto L330
	} else {
		goto L334
	}
L334:
	;
	v1925 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v15+int32(192)+v1913))))
	v1926 = v1906 + v1925
	switch v1913 + int32(-17) {
	case 0:
		goto L337
	case 1:
		goto L336
	default:
		v1933 = v1926
		goto L335
	}
L335:
	;
	if base.Ui32(int32(1)) < base.Ui32(v1910+int32(0)) {
		v1906 = v1933
		v1907 = v1907 + int32(-2)
		v1910 = v1910 + int32(-1)
		goto L331
	} else {
		goto L338
	}
L336:
	;
	v1933 = v1926 + int32(7)
	goto L335
L337:
	;
	v1933 = v1926 + int32(3)
	goto L335
L338:
	;
	goto L332
L339:
	;
	if v1962 != int32(1) {
		goto L366
	} else {
		goto L367
	}
L340:
	;
	goto L339
L341:
	;
	v1975 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1976 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v1976+v1960 < int32(32) {
		goto L344
	} else {
		goto L345
	}
L342:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2077 + v2075
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2074<<(uint(v2077)%32) | v2076
	goto L340
L343:
	;
	v1994 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1995 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1999 = v1991
	v2000 = v1992
	v2001 = v1995
	v2002 = v1994
	goto L347
L344:
	;
	if v1976 < int32(16) {
		v2074 = v1962
		v2075 = v1960
		v2076 = v1975
		v2077 = v1976
		goto L342
	} else {
		goto L346
	}
L345:
	;
	v1980 = int32(32)
	v1981 = v1980 - v1976
	v1989 = int32(base.Ui32(v1962) >> (uint(v1981) % 32))
	v1990 = v1960 - v1981
	v1991 = v1962<<(uint(v1976)%32) | v1975
	v1992 = v1980
	goto L343
L346:
	;
	v1989 = v1962
	v1990 = v1960
	v1991 = v1975
	v1992 = v1976
	goto L343
L347:
	;
	if base.Ui32(v2001+int32(2)) <= base.Ui32(v2002) {
		v2057 = v2001
		v2058 = v2002
		goto L349
	} else {
		goto L350
	}
L348:
	;
	v2074 = v1989
	v2075 = v1990
	v2076 = v2070
	v2077 = v2068
	goto L342
L349:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2057))) = uint16(v1999)
	v2065 = v2057 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2065
	v2068 = v2000 + int32(-16)
	v2070 = int32(base.Ui32(v1999) >> (uint(int32(16)) % 32))
	if int32(31) < v2000 {
		v1999 = v2070
		v2000 = v2068
		v2001 = v2065
		v2002 = v2058
		goto L347
	} else {
		goto L364
	}
L350:
	;
	v2011 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2012 = v2002 - v2011
	v2015 = base.I64_extend_i32_s(v2012) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2015) {
		v2039 = v2011
		goto L352
	} else {
		goto L353
	}
L351:
	;
	if v2001 == v2011 {
		goto L362
	} else {
		goto L363
	}
L352:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2039
	goto L339
L353:
	;
	v2018 = v2001 - v2011
	v2020 = v2015 + base.I64_extend_i32_u(v2018)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2020) {
		v2039 = v2011
		goto L352
	} else {
		goto L354
	}
L354:
	;
	v2023 = base.I32_wrap_i64(v2020)
	if v2002 == v2011 {
		goto L355
	} else {
		goto L356
	}
L355:
	;
	v2030 = int32(base.Ui32(v2012*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2023) < base.Ui32(v2030) {
		goto L358
	} else {
		goto L359
	}
L356:
	;
	if base.Ui32(v2023) <= base.Ui32(v2012) {
		v2057 = v2001
		v2058 = v2002
		goto L349
	} else {
		goto L357
	}
L357:
	;
	goto L355
L358:
	;
	v2032 = v2030
	goto L360
L359:
	;
	v2032 = v2023
	goto L360
L360:
	;
	v2036 = v2032&int32(-1024) + int32(1024)
	v2037 = F_WebPSafeMalloc(m, int64(1), v2036)
	mBase = m.M
	if v2037 != 0 {
		goto L351
	} else {
		goto L361
	}
L361:
	;
	v2038 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2039 = v2038
	goto L352
L362:
	;
	v2050 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2050)
	mBase = m.M
	v2052 = v2037 + v2036
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2052
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2037
	v2057 = v2037 + v2018
	v2058 = v2052
	goto L349
L363:
	;
	v2048 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2049 = F_memcpy(m, v2037, v2048, v2018)
	mBase = m.M
	goto L362
L364:
	;
	goto L348
L365:
	;
	v2589 = l2
	v2595 = v2586
	goto L449
L366:
	;
	if v1172 < int32(1) {
		goto L1
	} else {
		goto L448
	}
L367:
	;
	v2116 = int32(2)
	if v1957 != v2116 {
		goto L368
	} else {
		goto L369
	}
L368:
	;
	v2272 = v1957 + int32(-2)
	v2275 = base.I32_clz(v2272) ^ int32(31)
	v2277 = int32(base.Ui32(v2275) >> (uint(int32(1)) % 32))
	v2278 = int32(3)
	goto L398
L369:
	;
	v2119 = int32(0)
	v2120 = int32(5)
	goto L372
L370:
	;
	v2586 = v2116
	goto L365
L371:
	;
	goto L370
L372:
	;
	v2132 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2133 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2133+v2120 < int32(32) {
		goto L375
	} else {
		goto L376
	}
L373:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2234 + v2232
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2231<<(uint(v2234)%32) | v2233
	goto L371
L374:
	;
	v2151 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2152 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2156 = v2148
	v2157 = v2149
	v2158 = v2152
	v2159 = v2151
	goto L378
L375:
	;
	if v2133 < int32(16) {
		v2231 = v2119
		v2232 = v2120
		v2233 = v2132
		v2234 = v2133
		goto L373
	} else {
		goto L377
	}
L376:
	;
	v2137 = int32(32)
	v2138 = v2137 - v2133
	v2146 = int32(base.Ui32(v2119) >> (uint(v2138) % 32))
	v2147 = v2120 - v2138
	v2148 = v2119<<(uint(v2133)%32) | v2132
	v2149 = v2137
	goto L374
L377:
	;
	v2146 = v2119
	v2147 = v2120
	v2148 = v2132
	v2149 = v2133
	goto L374
L378:
	;
	if base.Ui32(v2158+int32(2)) <= base.Ui32(v2159) {
		v2214 = v2158
		v2215 = v2159
		goto L380
	} else {
		goto L381
	}
L379:
	;
	v2231 = v2146
	v2232 = v2147
	v2233 = v2227
	v2234 = v2225
	goto L373
L380:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2214))) = uint16(v2156)
	v2222 = v2214 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2222
	v2225 = v2157 + int32(-16)
	v2227 = int32(base.Ui32(v2156) >> (uint(int32(16)) % 32))
	if int32(31) < v2157 {
		v2156 = v2227
		v2157 = v2225
		v2158 = v2222
		v2159 = v2215
		goto L378
	} else {
		goto L395
	}
L381:
	;
	v2168 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2169 = v2159 - v2168
	v2172 = base.I64_extend_i32_s(v2169) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2172) {
		v2196 = v2168
		goto L383
	} else {
		goto L384
	}
L382:
	;
	if v2158 == v2168 {
		goto L393
	} else {
		goto L394
	}
L383:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2196
	goto L370
L384:
	;
	v2175 = v2158 - v2168
	v2177 = v2172 + base.I64_extend_i32_u(v2175)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2177) {
		v2196 = v2168
		goto L383
	} else {
		goto L385
	}
L385:
	;
	v2180 = base.I32_wrap_i64(v2177)
	if v2159 == v2168 {
		goto L386
	} else {
		goto L387
	}
L386:
	;
	v2187 = int32(base.Ui32(v2169*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2180) < base.Ui32(v2187) {
		goto L389
	} else {
		goto L390
	}
L387:
	;
	if base.Ui32(v2180) <= base.Ui32(v2169) {
		v2214 = v2158
		v2215 = v2159
		goto L380
	} else {
		goto L388
	}
L388:
	;
	goto L386
L389:
	;
	v2189 = v2187
	goto L391
L390:
	;
	v2189 = v2180
	goto L391
L391:
	;
	v2193 = v2189&int32(-1024) + int32(1024)
	v2194 = F_WebPSafeMalloc(m, int64(1), v2193)
	mBase = m.M
	if v2194 != 0 {
		goto L382
	} else {
		goto L392
	}
L392:
	;
	v2195 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2196 = v2195
	goto L383
L393:
	;
	v2207 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2207)
	mBase = m.M
	v2209 = v2194 + v2193
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2209
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2194
	v2214 = v2194 + v2175
	v2215 = v2209
	goto L380
L394:
	;
	v2205 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2206 = F_memcpy(m, v2194, v2205, v2175)
	mBase = m.M
	goto L393
L395:
	;
	goto L379
L396:
	;
	v2432 = v2275&int32(62) + int32(2)
	if v2432 < int32(1) {
		goto L423
	} else {
		goto L424
	}
L397:
	;
	goto L396
L398:
	;
	v2290 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2291 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2291+v2278 < int32(32) {
		goto L401
	} else {
		goto L402
	}
L399:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2392 + v2390
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2389<<(uint(v2392)%32) | v2391
	goto L397
L400:
	;
	v2309 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2310 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2314 = v2306
	v2315 = v2307
	v2316 = v2310
	v2317 = v2309
	goto L404
L401:
	;
	if v2291 < int32(16) {
		v2389 = v2277
		v2390 = v2278
		v2391 = v2290
		v2392 = v2291
		goto L399
	} else {
		goto L403
	}
L402:
	;
	v2295 = int32(32)
	v2296 = v2295 - v2291
	v2304 = int32(base.Ui32(v2277) >> (uint(v2296) % 32))
	v2305 = v2278 - v2296
	v2306 = v2277<<(uint(v2291)%32) | v2290
	v2307 = v2295
	goto L400
L403:
	;
	v2304 = v2277
	v2305 = v2278
	v2306 = v2290
	v2307 = v2291
	goto L400
L404:
	;
	if base.Ui32(v2316+int32(2)) <= base.Ui32(v2317) {
		v2372 = v2316
		v2373 = v2317
		goto L406
	} else {
		goto L407
	}
L405:
	;
	v2389 = v2304
	v2390 = v2305
	v2391 = v2385
	v2392 = v2383
	goto L399
L406:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2372))) = uint16(v2314)
	v2380 = v2372 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2380
	v2383 = v2315 + int32(-16)
	v2385 = int32(base.Ui32(v2314) >> (uint(int32(16)) % 32))
	if int32(31) < v2315 {
		v2314 = v2385
		v2315 = v2383
		v2316 = v2380
		v2317 = v2373
		goto L404
	} else {
		goto L421
	}
L407:
	;
	v2326 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2327 = v2317 - v2326
	v2330 = base.I64_extend_i32_s(v2327) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2330) {
		v2354 = v2326
		goto L409
	} else {
		goto L410
	}
L408:
	;
	if v2316 == v2326 {
		goto L419
	} else {
		goto L420
	}
L409:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2354
	goto L396
L410:
	;
	v2333 = v2316 - v2326
	v2335 = v2330 + base.I64_extend_i32_u(v2333)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2335) {
		v2354 = v2326
		goto L409
	} else {
		goto L411
	}
L411:
	;
	v2338 = base.I32_wrap_i64(v2335)
	if v2317 == v2326 {
		goto L412
	} else {
		goto L413
	}
L412:
	;
	v2345 = int32(base.Ui32(v2327*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2338) < base.Ui32(v2345) {
		goto L415
	} else {
		goto L416
	}
L413:
	;
	if base.Ui32(v2338) <= base.Ui32(v2327) {
		v2372 = v2316
		v2373 = v2317
		goto L406
	} else {
		goto L414
	}
L414:
	;
	goto L412
L415:
	;
	v2347 = v2345
	goto L417
L416:
	;
	v2347 = v2338
	goto L417
L417:
	;
	v2351 = v2347&int32(-1024) + int32(1024)
	v2352 = F_WebPSafeMalloc(m, int64(1), v2351)
	mBase = m.M
	if v2352 != 0 {
		goto L408
	} else {
		goto L418
	}
L418:
	;
	v2353 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2354 = v2353
	goto L409
L419:
	;
	v2365 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2365)
	mBase = m.M
	v2367 = v2352 + v2351
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2367
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2352
	v2372 = v2352 + v2333
	v2373 = v2367
	goto L406
L420:
	;
	v2363 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2364 = F_memcpy(m, v2352, v2363, v2333)
	mBase = m.M
	goto L419
L421:
	;
	goto L405
L422:
	;
	v2586 = v1957
	goto L365
L423:
	;
	goto L422
L424:
	;
	v2444 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2445 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2445+v2432 < int32(32) {
		goto L427
	} else {
		goto L428
	}
L425:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2546 + v2544
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2543<<(uint(v2546)%32) | v2545
	goto L423
L426:
	;
	v2463 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2464 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2468 = v2460
	v2469 = v2461
	v2470 = v2464
	v2471 = v2463
	goto L430
L427:
	;
	if v2445 < int32(16) {
		v2543 = v2272
		v2544 = v2432
		v2545 = v2444
		v2546 = v2445
		goto L425
	} else {
		goto L429
	}
L428:
	;
	v2449 = int32(32)
	v2450 = v2449 - v2445
	v2458 = int32(base.Ui32(v2272) >> (uint(v2450) % 32))
	v2459 = v2432 - v2450
	v2460 = v2272<<(uint(v2445)%32) | v2444
	v2461 = v2449
	goto L426
L429:
	;
	v2458 = v2272
	v2459 = v2432
	v2460 = v2444
	v2461 = v2445
	goto L426
L430:
	;
	if base.Ui32(v2470+int32(2)) <= base.Ui32(v2471) {
		v2526 = v2470
		v2527 = v2471
		goto L432
	} else {
		goto L433
	}
L431:
	;
	v2543 = v2458
	v2544 = v2459
	v2545 = v2539
	v2546 = v2537
	goto L425
L432:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2526))) = uint16(v2468)
	v2534 = v2526 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2534
	v2537 = v2469 + int32(-16)
	v2539 = int32(base.Ui32(v2468) >> (uint(int32(16)) % 32))
	if int32(31) < v2469 {
		v2468 = v2539
		v2469 = v2537
		v2470 = v2534
		v2471 = v2527
		goto L430
	} else {
		goto L447
	}
L433:
	;
	v2480 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2481 = v2471 - v2480
	v2484 = base.I64_extend_i32_s(v2481) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2484) {
		v2508 = v2480
		goto L435
	} else {
		goto L436
	}
L434:
	;
	if v2470 == v2480 {
		goto L445
	} else {
		goto L446
	}
L435:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2508
	goto L422
L436:
	;
	v2487 = v2470 - v2480
	v2489 = v2484 + base.I64_extend_i32_u(v2487)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2489) {
		v2508 = v2480
		goto L435
	} else {
		goto L437
	}
L437:
	;
	v2492 = base.I32_wrap_i64(v2489)
	if v2471 == v2480 {
		goto L438
	} else {
		goto L439
	}
L438:
	;
	v2499 = int32(base.Ui32(v2481*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2492) < base.Ui32(v2499) {
		goto L441
	} else {
		goto L442
	}
L439:
	;
	if base.Ui32(v2492) <= base.Ui32(v2481) {
		v2526 = v2470
		v2527 = v2471
		goto L432
	} else {
		goto L440
	}
L440:
	;
	goto L438
L441:
	;
	v2501 = v2499
	goto L443
L442:
	;
	v2501 = v2492
	goto L443
L443:
	;
	v2505 = v2501&int32(-1024) + int32(1024)
	v2506 = F_WebPSafeMalloc(m, int64(1), v2505)
	mBase = m.M
	if v2506 != 0 {
		goto L434
	} else {
		goto L444
	}
L444:
	;
	v2507 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2508 = v2507
	goto L435
L445:
	;
	v2519 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2519)
	mBase = m.M
	v2521 = v2506 + v2505
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2521
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2506
	v2526 = v2506 + v2487
	v2527 = v2521
	goto L432
L446:
	;
	v2517 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2518 = F_memcpy(m, v2506, v2517, v2487)
	mBase = m.M
	goto L445
L447:
	;
	goto L431
L448:
	;
	v2586 = v1172
	goto L365
L449:
	;
	v2599 = int32(1)
	v2601 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2589+v2599))))
	v2602 = *(*int32)(unsafe.Add(mBase, uint32(v15)+140))
	v2603 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2589))))
	v2607 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2602+v2603<<(uint(v2599)%32)))))
	v2608 = *(*int32)(unsafe.Add(mBase, uint32(v15)+136))
	v2610 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2608+v2603))))
	if v2610 < v2599 {
		goto L452
	} else {
		goto L453
	}
L450:
	;
	goto L1
L451:
	;
	v2764 = (v2603 + int32(-16)) & int32(255)
	if base.Ui32(int32(2)) < base.Ui32(v2764) {
		goto L477
	} else {
		goto L478
	}
L452:
	;
	goto L451
L453:
	;
	v2622 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2623 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2623+v2610 < int32(32) {
		goto L456
	} else {
		goto L457
	}
L454:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2724 + v2722
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2721<<(uint(v2724)%32) | v2723
	goto L452
L455:
	;
	v2641 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2642 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2646 = v2638
	v2647 = v2639
	v2648 = v2642
	v2649 = v2641
	goto L459
L456:
	;
	if v2623 < int32(16) {
		v2721 = v2607
		v2722 = v2610
		v2723 = v2622
		v2724 = v2623
		goto L454
	} else {
		goto L458
	}
L457:
	;
	v2627 = int32(32)
	v2628 = v2627 - v2623
	v2636 = int32(base.Ui32(v2607) >> (uint(v2628) % 32))
	v2637 = v2610 - v2628
	v2638 = v2607<<(uint(v2623)%32) | v2622
	v2639 = v2627
	goto L455
L458:
	;
	v2636 = v2607
	v2637 = v2610
	v2638 = v2622
	v2639 = v2623
	goto L455
L459:
	;
	if base.Ui32(v2648+int32(2)) <= base.Ui32(v2649) {
		v2704 = v2648
		v2705 = v2649
		goto L461
	} else {
		goto L462
	}
L460:
	;
	v2721 = v2636
	v2722 = v2637
	v2723 = v2717
	v2724 = v2715
	goto L454
L461:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2704))) = uint16(v2646)
	v2712 = v2704 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2712
	v2715 = v2647 + int32(-16)
	v2717 = int32(base.Ui32(v2646) >> (uint(int32(16)) % 32))
	if int32(31) < v2647 {
		v2646 = v2717
		v2647 = v2715
		v2648 = v2712
		v2649 = v2705
		goto L459
	} else {
		goto L476
	}
L462:
	;
	v2658 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2659 = v2649 - v2658
	v2662 = base.I64_extend_i32_s(v2659) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2662) {
		v2686 = v2658
		goto L464
	} else {
		goto L465
	}
L463:
	;
	if v2648 == v2658 {
		goto L474
	} else {
		goto L475
	}
L464:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2686
	goto L451
L465:
	;
	v2665 = v2648 - v2658
	v2667 = v2662 + base.I64_extend_i32_u(v2665)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2667) {
		v2686 = v2658
		goto L464
	} else {
		goto L466
	}
L466:
	;
	v2670 = base.I32_wrap_i64(v2667)
	if v2649 == v2658 {
		goto L467
	} else {
		goto L468
	}
L467:
	;
	v2677 = int32(base.Ui32(v2659*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2670) < base.Ui32(v2677) {
		goto L470
	} else {
		goto L471
	}
L468:
	;
	if base.Ui32(v2670) <= base.Ui32(v2659) {
		v2704 = v2648
		v2705 = v2649
		goto L461
	} else {
		goto L469
	}
L469:
	;
	goto L467
L470:
	;
	v2679 = v2677
	goto L472
L471:
	;
	v2679 = v2670
	goto L472
L472:
	;
	v2683 = v2679&int32(-1024) + int32(1024)
	v2684 = F_WebPSafeMalloc(m, int64(1), v2683)
	mBase = m.M
	if v2684 != 0 {
		goto L463
	} else {
		goto L473
	}
L473:
	;
	v2685 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2686 = v2685
	goto L464
L474:
	;
	v2697 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2697)
	mBase = m.M
	v2699 = v2684 + v2683
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2699
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2684
	v2704 = v2684 + v2665
	v2705 = v2699
	goto L461
L475:
	;
	v2695 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2696 = F_memcpy(m, v2684, v2695, v2665)
	mBase = m.M
	goto L474
L476:
	;
	goto L460
L477:
	;
	v2927 = v2595 + int32(-1)
	if v2927 != 0 {
		v2589 = v2589 + int32(2)
		v2595 = v2927
		goto L449
	} else {
		goto L505
	}
L478:
	;
	v2767 = m.G1
	v2773 = *(*int32)(unsafe.Add(mBase, uint32(v2767+int32(_a_F_StoreHuffmanCode_1)+v2764<<(uint(int32(2))%32))))
	if v2773 < int32(1) {
		goto L480
	} else {
		goto L481
	}
L479:
	;
	goto L477
L480:
	;
	goto L479
L481:
	;
	v2785 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v2786 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v2786+v2773 < int32(32) {
		goto L484
	} else {
		goto L485
	}
L482:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v2887 + v2885
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v2884<<(uint(v2887)%32) | v2886
	goto L480
L483:
	;
	v2804 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v2805 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v2809 = v2801
	v2810 = v2802
	v2811 = v2805
	v2812 = v2804
	goto L487
L484:
	;
	if v2786 < int32(16) {
		v2884 = v2601
		v2885 = v2773
		v2886 = v2785
		v2887 = v2786
		goto L482
	} else {
		goto L486
	}
L485:
	;
	v2790 = int32(32)
	v2791 = v2790 - v2786
	v2799 = int32(base.Ui32(v2601) >> (uint(v2791) % 32))
	v2800 = v2773 - v2791
	v2801 = v2601<<(uint(v2786)%32) | v2785
	v2802 = v2790
	goto L483
L486:
	;
	v2799 = v2601
	v2800 = v2773
	v2801 = v2785
	v2802 = v2786
	goto L483
L487:
	;
	if base.Ui32(v2811+int32(2)) <= base.Ui32(v2812) {
		v2867 = v2811
		v2868 = v2812
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v2884 = v2799
	v2885 = v2800
	v2886 = v2880
	v2887 = v2878
	goto L482
L489:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2867))) = uint16(v2809)
	v2875 = v2867 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2875
	v2878 = v2810 + int32(-16)
	v2880 = int32(base.Ui32(v2809) >> (uint(int32(16)) % 32))
	if int32(31) < v2810 {
		v2809 = v2880
		v2810 = v2878
		v2811 = v2875
		v2812 = v2868
		goto L487
	} else {
		goto L504
	}
L490:
	;
	v2821 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2822 = v2812 - v2821
	v2825 = base.I64_extend_i32_s(v2822) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2825) {
		v2849 = v2821
		goto L492
	} else {
		goto L493
	}
L491:
	;
	if v2811 == v2821 {
		goto L502
	} else {
		goto L503
	}
L492:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v2849
	goto L479
L493:
	;
	v2828 = v2811 - v2821
	v2830 = v2825 + base.I64_extend_i32_u(v2828)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2830) {
		v2849 = v2821
		goto L492
	} else {
		goto L494
	}
L494:
	;
	v2833 = base.I32_wrap_i64(v2830)
	if v2812 == v2821 {
		goto L495
	} else {
		goto L496
	}
L495:
	;
	v2840 = int32(base.Ui32(v2822*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v2833) < base.Ui32(v2840) {
		goto L498
	} else {
		goto L499
	}
L496:
	;
	if base.Ui32(v2833) <= base.Ui32(v2822) {
		v2867 = v2811
		v2868 = v2812
		goto L489
	} else {
		goto L497
	}
L497:
	;
	goto L495
L498:
	;
	v2842 = v2840
	goto L500
L499:
	;
	v2842 = v2833
	goto L500
L500:
	;
	v2846 = v2842&int32(-1024) + int32(1024)
	v2847 = F_WebPSafeMalloc(m, int64(1), v2846)
	mBase = m.M
	if v2847 != 0 {
		goto L491
	} else {
		goto L501
	}
L501:
	;
	v2848 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2849 = v2848
	goto L492
L502:
	;
	v2860 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v2860)
	mBase = m.M
	v2862 = v2847 + v2846
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v2862
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v2847
	v2867 = v2847 + v2828
	v2868 = v2862
	goto L489
L503:
	;
	v2858 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v2859 = F_memcpy(m, v2847, v2858, v2828)
	mBase = m.M
	goto L502
L504:
	;
	goto L488
L505:
	;
	goto L450
}

var F_StoreHuffmanCode__k0 = [2]uint64{0x0, 0x0}

func F_SubtractGreenFromBlueAndRed_SSE2(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v15 int32
	_ = v15
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v37 int32
	_ = v37
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v48 int32
	_ = v48
	var v50 base.V128
	_ = v50
	var v57 base.V128
	_ = v57
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 base.V128
	_ = v83
	var v86 base.V128
	_ = v86
	var v91 base.V128
	_ = v91
	var v97 int32
	_ = v97
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v128 base.V128
	_ = v128
	var v130 base.V128
	_ = v130
	var v137 int32
	_ = v137
	var v142 base.V128
	_ = v142
	var v148 int32
	_ = v148
	var v152 int32
	_ = v152
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v168 int32
	_ = v168
	var v170 int32
	_ = v170
	var v187 int32
	_ = v187
	v10 = int32(4)
	if l1 < v10 {
		v97 = int32(0)
	} else {
		v15 = l1 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v15) {
			v28 = int32(4)
			v29 = l0
			v31 = (int32(base.Ui32(v15)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v35 = int32(0)
				v36 = base.Simd_g_v128_load(m, v29, v35)
				v37 = int32(8)
				v39 = base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k0)
				v40 = base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k1)
				v42 = base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k2)
				v44 = base.Simd_g_i8x16_sub(v36, base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_shr_u(v36, v37), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k1)), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k2)))
				base.Simd_g_v128_store(m, v29, v35, v44)
				v48 = v29 + int32(16)
				v50 = base.Simd_g_v128_load(m, v48, v35)
				v57 = base.Simd_g_i8x16_sub(v50, base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_shr_u(v50, v37), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k1)), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k2)))
				base.Simd_g_v128_store(m, v48, v35, v57)
				v63 = v28 + v37
				v65 = v31 + int32(-2)
				if v65 != 0 {
					v28 = v63
					v29 = v29 + int32(32)
					v31 = v65
					continue
				} else {
					break
				}
				break
			}
			v70 = v63
			v71 = v28 + int32(4)
		} else {
			v70 = v10
			v71 = int32(0)
		}
		if v15&int32(4) != 0 {
			v97 = v71
		} else {
			v81 = l0 + v71<<(uint(int32(2))%32)
			v82 = int32(0)
			v83 = base.Simd_g_v128_load(m, v81, v82)
			v86 = base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k0)
			v91 = base.Simd_g_i8x16_sub(v83, base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_shr_u(v83, int32(8)), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k1)), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k2)))
			base.Simd_g_v128_store(m, v81, v82, v91)
			v97 = v70
		}
	}
	if l1 == v97 {
	} else {
		v106 = l0 + v97<<(uint(int32(2))%32)
		v107 = l1 - v97
		if v107 < int32(1) {
		} else {
			if base.Ui32(v107) < base.Ui32(int32(4)) {
				v152 = int32(0)
				v162 = v107 - v152
				v165 = v106 + v152<<(uint(int32(2))%32)
				for {
					v168 = *(*int32)(unsafe.Add(mBase, uint32(v165)))
					v170 = int32(base.Ui32(v168) >> (uint(int32(8)) % 32))
					*(*int32)(unsafe.Add(mBase, uint32(v165))) = (v168-v170)&int32(255) | v168&int32(-16711936) | (v168-v170<<(uint(int32(16))%32))&int32(16711680)
					v187 = v162 + int32(-1)
					if v187 != 0 {
						v162 = v187
						v165 = v165 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			} else {
				v119 = v107 & int32(2147483644)
				v123 = v119
				v124 = v106
				for {
					v127 = int32(0)
					v128 = base.Simd_g_v128_load(m, v124, v127)
					v130 = base.Simd_g_i32x4_shr_u(v128, int32(8))
					v137 = int32(16)
					v142 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v128, v130), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k3)), base.Simd_g_v128_and(v128, base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k4))), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v128, base.Simd_g_i32x4_shl(v130, v137)), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE2__k5)))
					base.Simd_g_v128_store(m, v124, v127, v142)
					v148 = v123 + int32(-4)
					if v148 != 0 {
						v123 = v148
						v124 = v124 + v137
						continue
					} else {
						break
					}
					break
				}
				if v119 == v107 {
				} else {
					v152 = v119
					v162 = v107 - v152
					v165 = v106 + v152<<(uint(int32(2))%32)
					for {
						v168 = *(*int32)(unsafe.Add(mBase, uint32(v165)))
						v170 = int32(base.Ui32(v168) >> (uint(int32(8)) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(v165))) = (v168-v170)&int32(255) | v168&int32(-16711936) | (v168-v170<<(uint(int32(16))%32))&int32(16711680)
						v187 = v162 + int32(-1)
						if v187 != 0 {
							v162 = v187
							v165 = v165 + int32(4)
							continue
						} else {
							break
						}
						break
					}
				}
			}
		}
	}
	return
}

var F_SubtractGreenFromBlueAndRed_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_SubtractGreenFromBlueAndRed_SSE2__k1 = [2]uint64{0x504050401000100, 0xf0e0d0c0b0a0908}
var F_SubtractGreenFromBlueAndRed_SSE2__k2 = [2]uint64{0x706050403020100, 0xd0c0d0c09080908}
var F_SubtractGreenFromBlueAndRed_SSE2__k3 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_SubtractGreenFromBlueAndRed_SSE2__k4 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_SubtractGreenFromBlueAndRed_SSE2__k5 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}

func F_SubtractGreenFromBlueAndRed_SSE41(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v41 base.V128
	_ = v41
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v53 int32
	_ = v53
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v61 int32
	_ = v61
	var v63 base.V128
	_ = v63
	var v65 base.V128
	_ = v65
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v89 int32
	_ = v89
	var __phi89 int32
	_ = __phi89
	var v90 int32
	_ = v90
	var __phi90 int32
	_ = __phi90
	var v92 int32
	_ = v92
	var __phi92 int32
	_ = __phi92
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v100 base.V128
	_ = v100
	var v103 base.V128
	_ = v103
	var v109 int32
	_ = v109
	var v115 int32
	_ = v115
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v137 int32
	_ = v137
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v146 base.V128
	_ = v146
	var v148 base.V128
	_ = v148
	var v155 int32
	_ = v155
	var v160 base.V128
	_ = v160
	var v166 int32
	_ = v166
	var v170 int32
	_ = v170
	var v180 int32
	_ = v180
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v205 int32
	_ = v205
	v10 = int32(4)
	if l1 < v10 {
		v115 = int32(0)
	} else {
		v15 = l1 + int32(-4)
		v19 = int32(base.Ui32(v15)>>(uint(int32(2))%32)) + int32(1)
		v21 = v19 & int32(3)
		if base.Ui32(int32(12)) <= base.Ui32(v15) {
			v30 = int32(4)
			v31 = l0
			v32 = v19 & int32(2147483644)
			for {
				v37 = int32(0)
				v38 = base.Simd_g_v128_load(m, v31, v37)
				v39 = base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE41__k0)
				v41 = base.Simd_g_i8x16_sub(v38, base.Simd_g_i8x16_swizzle(v38, v39))
				base.Simd_g_v128_store(m, v31, v37, v41)
				v44 = int32(16)
				v45 = v31 + v44
				v47 = base.Simd_g_v128_load(m, v45, v37)
				v49 = base.Simd_g_i8x16_sub(v47, base.Simd_g_i8x16_swizzle(v47, v39))
				base.Simd_g_v128_store(m, v45, v37, v49)
				v53 = v31 + int32(32)
				v55 = base.Simd_g_v128_load(m, v53, v37)
				v57 = base.Simd_g_i8x16_sub(v55, base.Simd_g_i8x16_swizzle(v55, v39))
				base.Simd_g_v128_store(m, v53, v37, v57)
				v61 = v31 + int32(48)
				v63 = base.Simd_g_v128_load(m, v61, v37)
				v65 = base.Simd_g_i8x16_sub(v63, base.Simd_g_i8x16_swizzle(v63, v39))
				base.Simd_g_v128_store(m, v61, v37, v65)
				v71 = v30 + v44
				v73 = v32 + int32(-4)
				if v73 != 0 {
					v30 = v71
					v31 = v31 + int32(64)
					v32 = v73
					continue
				} else {
					break
				}
				break
			}
			v78 = v71
			v79 = v30 + int32(12)
		} else {
			v78 = v10
			v79 = int32(0)
		}
		if v21 == int32(0) {
			v115 = v79
		} else {
			__phi89 = v78
			__phi90 = v79
			__phi92 = v21
			v89 = __phi89
			v90 = __phi90
			v92 = __phi92
			for {
				v98 = l0 + v90<<(uint(int32(2))%32)
				v99 = int32(0)
				v100 = base.Simd_g_v128_load(m, v98, v99)
				v103 = base.Simd_g_i8x16_sub(v100, base.Simd_g_i8x16_swizzle(v100, base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE41__k0)))
				base.Simd_g_v128_store(m, v98, v99, v103)
				v109 = v92 + int32(-1)
				if v109 != 0 {
					__phi89 = v89 + int32(4)
					__phi90 = v89
					__phi92 = v109
					v89 = __phi89
					v90 = __phi90
					v92 = __phi92
					continue
				} else {
					break
				}
				break
			}
			v115 = v89 + int32(0)
		}
	}
	if l1 == v115 {
	} else {
		v124 = l0 + v115<<(uint(int32(2))%32)
		v125 = l1 - v115
		if v125 < int32(1) {
		} else {
			if base.Ui32(v125) < base.Ui32(int32(4)) {
				v170 = int32(0)
				v180 = v125 - v170
				v183 = v124 + v170<<(uint(int32(2))%32)
				for {
					v186 = *(*int32)(unsafe.Add(mBase, uint32(v183)))
					v188 = int32(base.Ui32(v186) >> (uint(int32(8)) % 32))
					*(*int32)(unsafe.Add(mBase, uint32(v183))) = (v186-v188)&int32(255) | v186&int32(-16711936) | (v186-v188<<(uint(int32(16))%32))&int32(16711680)
					v205 = v180 + int32(-1)
					if v205 != 0 {
						v180 = v205
						v183 = v183 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			} else {
				v137 = v125 & int32(2147483644)
				v141 = v137
				v142 = v124
				for {
					v145 = int32(0)
					v146 = base.Simd_g_v128_load(m, v142, v145)
					v148 = base.Simd_g_i32x4_shr_u(v146, int32(8))
					v155 = int32(16)
					v160 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v146, v148), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE41__k1)), base.Simd_g_v128_and(v146, base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE41__k2))), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(v146, base.Simd_g_i32x4_shl(v148, v155)), base.Simd_g_const(&F_SubtractGreenFromBlueAndRed_SSE41__k3)))
					base.Simd_g_v128_store(m, v142, v145, v160)
					v166 = v141 + int32(-4)
					if v166 != 0 {
						v141 = v166
						v142 = v142 + v155
						continue
					} else {
						break
					}
					break
				}
				if v137 == v125 {
				} else {
					v170 = v137
					v180 = v125 - v170
					v183 = v124 + v170<<(uint(int32(2))%32)
					for {
						v186 = *(*int32)(unsafe.Add(mBase, uint32(v183)))
						v188 = int32(base.Ui32(v186) >> (uint(int32(8)) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(v183))) = (v186-v188)&int32(255) | v186&int32(-16711936) | (v186-v188<<(uint(int32(16))%32))&int32(16711680)
						v205 = v180 + int32(-1)
						if v205 != 0 {
							v180 = v205
							v183 = v183 + int32(4)
							continue
						} else {
							break
						}
						break
					}
				}
			}
		}
	}
	return
}

var F_SubtractGreenFromBlueAndRed_SSE41__k0 = [2]uint64{0x8f058f058f018f01, 0x8f0d8f0d8f098f09}
var F_SubtractGreenFromBlueAndRed_SSE41__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_SubtractGreenFromBlueAndRed_SSE41__k2 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_SubtractGreenFromBlueAndRed_SSE41__k3 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}

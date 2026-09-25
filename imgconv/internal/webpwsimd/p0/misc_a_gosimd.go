//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_AccumulateSSE_C(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v26 base.V128
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v35 base.V128
	_ = v35
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v41 int32
	_ = v41
	var v46 int32
	_ = v46
	var v49 base.V128
	_ = v49
	var v54 int32
	_ = v54
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v87 int32
	_ = v87
	var v92 int32
	_ = v92
	if int32(1) <= l2 {
		v14 = int32(0)
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v59 = v14
			v60 = v14
			v70 = l2 - v59
			v72 = v60
			v74 = l0 + v59
			v75 = l1 + v59
			for {
				v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v74))))
				v78 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v75))))
				v79 = v77 - v78
				v81 = v79*v79 + v72
				v82 = int32(1)
				v87 = v70 + int32(-1)
				if v87 != 0 {
					v70 = v87
					v72 = v81
					v74 = v74 + v82
					v75 = v75 + v82
					continue
				} else {
					break
				}
				break
			}
			v92 = v81
		} else {
			v20 = l2 & int32(2147483644)
			v25 = v20
			v26 = base.Simd_g_const(&F_AccumulateSSE_C__k0)
			v27 = l1
			v28 = l0
			for {
				v30 = int32(0)
				v31 = base.Simd_g_v128_load32_zero(m, v28, v30)
				v35 = base.Simd_g_v128_load32_zero(m, v27, v30)
				v38 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v31)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v35)))
				v40 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v38, v38), v26)
				v41 = int32(4)
				v46 = v25 + int32(-4)
				if v46 != 0 {
					v25 = v46
					v26 = v40
					v27 = v27 + v41
					v28 = v28 + v41
					continue
				} else {
					break
				}
				break
			}
			v49 = base.Simd_g_i32x4_add(v40, base.Simd_g_i8x16_shuffle2(v40, v40, base.Simd_g_const(&F_AccumulateSSE_C__k1), base.Simd_g_const(&F_AccumulateSSE_C__k2)))
			v54 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v49, base.Simd_g_i8x16_shuffle2(v49, v49, base.Simd_g_const(&F_AccumulateSSE_C__k3), base.Simd_g_const(&F_AccumulateSSE_C__k2))))
			if v20 == l2 {
				v92 = v54
			} else {
				v59 = v20
				v60 = v54
				v70 = l2 - v59
				v72 = v60
				v74 = l0 + v59
				v75 = l1 + v59
				for {
					v77 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v74))))
					v78 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v75))))
					v79 = v77 - v78
					v81 = v79*v79 + v72
					v82 = int32(1)
					v87 = v70 + int32(-1)
					if v87 != 0 {
						v70 = v87
						v72 = v81
						v74 = v74 + v82
						v75 = v75 + v82
						continue
					} else {
						break
					}
					break
				}
				v92 = v81
			}
		}
		return v92
	} else {
		return int32(0)
	}
}

var F_AccumulateSSE_C__k0 = [2]uint64{0x0, 0x0}
var F_AccumulateSSE_C__k1 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_AccumulateSSE_C__k2 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_AccumulateSSE_C__k3 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_AccumulateSSE_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v18 int32
	_ = v18
	var v19 base.V128
	_ = v19
	var v21 base.V128
	_ = v21
	var v26 int32
	_ = v26
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v44 int32
	_ = v44
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 int32
	_ = v64
	var v68 base.V128
	_ = v68
	var v71 base.V128
	_ = v71
	var v73 base.V128
	_ = v73
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v82 int32
	_ = v82
	var v85 base.V128
	_ = v85
	var v89 base.V128
	_ = v89
	var v99 int32
	_ = v99
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v102 base.V128
	_ = v102
	var v110 base.V128
	_ = v110
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v117 base.V128
	_ = v117
	var v119 base.V128
	_ = v119
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v145 int32
	_ = v145
	var v154 int32
	_ = v154
	var v159 int32
	_ = v159
	var v161 base.V128
	_ = v161
	var v166 int32
	_ = v166
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v170 base.V128
	_ = v170
	var v174 base.V128
	_ = v174
	var v177 base.V128
	_ = v177
	var v179 base.V128
	_ = v179
	var v180 int32
	_ = v180
	var v185 int32
	_ = v185
	var v188 base.V128
	_ = v188
	var v193 int32
	_ = v193
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v214 int32
	_ = v214
	var v215 int32
	_ = v215
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v234 int32
	_ = v234
	var v238 int32
	_ = v238
	v4 = int32(0)
	if l2 < int32(16) {
		v134 = v4
		v135 = v4
	} else {
		v18 = int32(0)
		v19 = base.Simd_g_v128_load(m, l0, v18)
		v21 = base.Simd_g_v128_load(m, l1, v18)
		if base.Ui32(int32(48)) <= base.Ui32(l2) {
			v26 = int32(-32)
			v39 = v19
			v40 = v21
			v41 = base.Simd_g_const(&F_AccumulateSSE_SSE2__k0)
			v44 = int32(0)
			for {
				v49 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v40, v39), base.Simd_g_i8x16_sub_sat_u(v39, v40))
				v50 = base.Simd_g_const(&F_AccumulateSSE_SSE2__k0)
				v51 = base.Simd_g_const(&F_AccumulateSSE_SSE2__k1)
				v52 = base.Simd_g_i8x16_shuffle2(v49, v50, base.Simd_g_const(&F_AccumulateSSE_SSE2__k2), base.Simd_g_const(&F_AccumulateSSE_SSE2__k3))
				v55 = base.Simd_g_const(&F_AccumulateSSE_SSE2__k4)
				v56 = base.Simd_g_i8x16_shuffle2(v49, v50, base.Simd_g_const(&F_AccumulateSSE_SSE2__k5), base.Simd_g_const(&F_AccumulateSSE_SSE2__k6))
				v59 = l1 + v44
				v60 = int32(16)
				v62 = int32(0)
				v63 = base.Simd_g_v128_load_rng(m, v59+v60, v62, int32(0), int32(32))
				v64 = l0 + v44
				v68 = base.Simd_g_v128_load_rng(m, v64+v60, v62, int32(0), int32(32))
				v71 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v63, v68), base.Simd_g_i8x16_sub_sat_u(v68, v63))
				v73 = base.Simd_g_i8x16_shuffle2(v71, v50, base.Simd_g_const(&F_AccumulateSSE_SSE2__k2), base.Simd_g_const(&F_AccumulateSSE_SSE2__k3))
				v77 = base.Simd_g_i8x16_shuffle2(v71, v50, base.Simd_g_const(&F_AccumulateSSE_SSE2__k5), base.Simd_g_const(&F_AccumulateSSE_SSE2__k6))
				v79 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v52, v52), v41), base.Simd_g_i32x4_dot_i16x8_s(v56, v56)), base.Simd_g_i32x4_dot_i16x8_s(v73, v73)), base.Simd_g_i32x4_dot_i16x8_s(v77, v77))
				v82 = int32(32)
				v85 = base.Simd_g_v128_load_nc(m, v64+v82, v62)
				v89 = base.Simd_g_v128_load_nc(m, v59+v82, v62)
				if v44+int32(48) <= l2+v26 {
					v39 = v85
					v40 = v89
					v41 = v79
					v44 = v44 + v82
					continue
				} else {
					break
				}
				break
			}
			v99 = (l2+int32(-48))&v26 + int32(48)
			v100 = v85
			v101 = v89
			v102 = v79
		} else {
			v99 = int32(16)
			v100 = v19
			v101 = v21
			v102 = base.Simd_g_const(&F_AccumulateSSE_SSE2__k0)
		}
		v110 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v101, v100), base.Simd_g_i8x16_sub_sat_u(v100, v101))
		v111 = base.Simd_g_const(&F_AccumulateSSE_SSE2__k0)
		v113 = base.Simd_g_i8x16_shuffle2(v110, v111, base.Simd_g_const(&F_AccumulateSSE_SSE2__k2), base.Simd_g_const(&F_AccumulateSSE_SSE2__k3))
		v117 = base.Simd_g_i8x16_shuffle2(v110, v111, base.Simd_g_const(&F_AccumulateSSE_SSE2__k5), base.Simd_g_const(&F_AccumulateSSE_SSE2__k6))
		v119 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v113, v113), v102), base.Simd_g_i32x4_dot_i16x8_s(v117, v117))
		v134 = base.Simd_g_i32x4_extract_lane_l3(v119) + base.Simd_g_i32x4_extract_lane_l2(v119) + base.Simd_g_i32x4_extract_lane_l1(v119) + base.Simd_g_i32x4_extract_lane_l0(v119)
		v135 = v99
	}
	if l2 <= v135 {
		v238 = v134
	} else {
		v145 = l2 - v135
		if base.Ui32(v145) < base.Ui32(int32(4)) {
			v198 = v134
			v199 = v135
			v214 = v198
			v215 = l2 - v199
			v221 = l0 + v199
			v223 = l1 + v199
			for {
				v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v221))))
				v225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v223))))
				v226 = v224 - v225
				v228 = v226*v226 + v214
				v229 = int32(1)
				v234 = v215 + int32(-1)
				if v234 != 0 {
					v214 = v228
					v215 = v234
					v221 = v221 + v229
					v223 = v223 + v229
					continue
				} else {
					break
				}
				break
			}
			v238 = v228
		} else {
			v154 = v145 & int32(-4)
			v159 = v154
			v161 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_AccumulateSSE_SSE2__k0), v134)
			v166 = l1 + v135
			v168 = l0 + v135
			for {
				v169 = int32(0)
				v170 = base.Simd_g_v128_load32_zero(m, v168, v169)
				v174 = base.Simd_g_v128_load32_zero(m, v166, v169)
				v177 = base.Simd_g_i32x4_sub(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v170)), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v174)))
				v179 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v177, v177), v161)
				v180 = int32(4)
				v185 = v159 + int32(-4)
				if v185 != 0 {
					v159 = v185
					v161 = v179
					v166 = v166 + v180
					v168 = v168 + v180
					continue
				} else {
					break
				}
				break
			}
			v188 = base.Simd_g_i32x4_add(v179, base.Simd_g_i8x16_shuffle2(v179, v179, base.Simd_g_const(&F_AccumulateSSE_SSE2__k7), base.Simd_g_const(&F_AccumulateSSE_SSE2__k8)))
			v193 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v188, base.Simd_g_i8x16_shuffle2(v188, v188, base.Simd_g_const(&F_AccumulateSSE_SSE2__k9), base.Simd_g_const(&F_AccumulateSSE_SSE2__k8))))
			if v145 == v154 {
				v238 = v193
			} else {
				v198 = v193
				v199 = v135 + v154
				v214 = v198
				v215 = l2 - v199
				v221 = l0 + v199
				v223 = l1 + v199
				for {
					v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v221))))
					v225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v223))))
					v226 = v224 - v225
					v228 = v226*v226 + v214
					v229 = int32(1)
					v234 = v215 + int32(-1)
					if v234 != 0 {
						v214 = v228
						v215 = v234
						v221 = v221 + v229
						v223 = v223 + v229
						continue
					} else {
						break
					}
					break
				}
				v238 = v228
			}
		}
	}
	return v238
}

var F_AccumulateSSE_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_AccumulateSSE_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_AccumulateSSE_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_AccumulateSSE_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_AccumulateSSE_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_AccumulateSSE_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_AccumulateSSE_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_AccumulateSSE_SSE2__k7 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_AccumulateSSE_SSE2__k8 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_AccumulateSSE_SSE2__k9 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_AddGreenToBlueAndRed_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v16 int32
	_ = v16
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v39 int32
	_ = v39
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v46 base.V128
	_ = v46
	var v49 int32
	_ = v49
	var v54 base.V128
	_ = v54
	var v61 base.V128
	_ = v61
	var v64 int32
	_ = v64
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v91 base.V128
	_ = v91
	var v94 base.V128
	_ = v94
	var v99 base.V128
	_ = v99
	var v106 int32
	_ = v106
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v127 int32
	_ = v127
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v145 int32
	_ = v145
	var v146 base.V128
	_ = v146
	var v148 base.V128
	_ = v148
	var v151 base.V128
	_ = v151
	var v154 int32
	_ = v154
	var v160 base.V128
	_ = v160
	var v168 int32
	_ = v168
	var v173 int32
	_ = v173
	var v185 int32
	_ = v185
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v207 int32
	_ = v207
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v235 int32
	_ = v235
	var v240 int32
	_ = v240
	var v244 int32
	_ = v244
	var v248 int32
	_ = v248
	var v250 int32
	_ = v250
	var v270 int32
	_ = v270
	v11 = int32(4)
	if l1 < v11 {
		v106 = int32(0)
	} else {
		v16 = l1 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v16) {
			v30 = int32(4)
			v31 = l0
			v33 = (int32(base.Ui32(v16)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			v34 = l2
			for {
				v37 = int32(0)
				v38 = base.Simd_g_v128_load(m, v31, v37)
				v39 = int32(8)
				v41 = base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k0)
				v42 = base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k1)
				v44 = base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k2)
				v46 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_shr_u(v38, v39), base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k1)), base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k2)), v38)
				base.Simd_g_v128_store(m, v34, v37, v46)
				v49 = int32(16)
				v54 = base.Simd_g_v128_load(m, v31+v49, v37)
				v61 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_shr_u(v54, v39), base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k1)), base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k2)), v54)
				base.Simd_g_v128_store(m, v34+v49, v37, v61)
				v64 = int32(32)
				v69 = v30 + v39
				v71 = v33 + int32(-2)
				if v71 != 0 {
					v30 = v69
					v31 = v31 + v64
					v33 = v71
					v34 = v34 + v64
					continue
				} else {
					break
				}
				break
			}
			v77 = v69
			v78 = v30 + int32(4)
		} else {
			v77 = v11
			v78 = int32(0)
		}
		if v16&int32(4) != 0 {
			v106 = v78
		} else {
			v87 = v78 << (uint(int32(2)) % 32)
			v90 = int32(0)
			v91 = base.Simd_g_v128_load(m, l0+v87, v90)
			v94 = base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k0)
			v99 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_shr_u(v91, int32(8)), base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k1)), base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k2)), v91)
			base.Simd_g_v128_store(m, l2+v87, v90, v99)
			v106 = v77
		}
	}
	if l1 == v106 {
	} else {
		v114 = v106 << (uint(int32(2)) % 32)
		v115 = l0 + v114
		v116 = l1 - v106
		v117 = l2 + v114
		if v116 < int32(1) {
		} else {
			v127 = int32(0)
			if base.Ui32(v116) < base.Ui32(int32(4)) {
				v173 = v127
				if v116&int32(1) == int32(0) {
					v207 = v173
				} else {
					v185 = v173 << (uint(int32(2)) % 32)
					v188 = *(*int32)(unsafe.Add(mBase, uint32(v115+v185)))
					v190 = int32(base.Ui32(v188) >> (uint(int32(8)) % 32))
					v193 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v117+v185))) = (v190&int32(255)+v188&v193+v190<<(uint(int32(16))%32))&v193 | v188&int32(-16711936)
					v207 = v173 | int32(1)
				}
				if v173 == v116+int32(-1) {
				} else {
					v214 = v207 << (uint(int32(2)) % 32)
					v218 = v116 - v207
					v222 = v115 + v214
					v223 = v117 + v214
					for {
						v227 = *(*int32)(unsafe.Add(mBase, uint32(v222)))
						v228 = int32(8)
						v229 = int32(base.Ui32(v227) >> (uint(v228) % 32))
						v230 = int32(255)
						v232 = int32(16711935)
						v235 = int32(16)
						v240 = int32(-16711936)
						*(*int32)(unsafe.Add(mBase, uint32(v223))) = (v229&v230+v227&v232+v229<<(uint(v235)%32))&v232 | v227&v240
						v244 = int32(4)
						v248 = *(*int32)(unsafe.Add(mBase, uint32(v222+v244)))
						v250 = int32(base.Ui32(v248) >> (uint(v228) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(v223+v244))) = (v250&v230+v248&v232+v250<<(uint(v235)%32))&v232 | v248&v240
						v270 = v218 + int32(-2)
						if v270 != 0 {
							v218 = v270
							v222 = v222 + v228
							v223 = v223 + v228
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if base.Ui32(v117-v115) < base.Ui32(int32(16)) {
					v173 = v127
					if v116&int32(1) == int32(0) {
						v207 = v173
					} else {
						v185 = v173 << (uint(int32(2)) % 32)
						v188 = *(*int32)(unsafe.Add(mBase, uint32(v115+v185)))
						v190 = int32(base.Ui32(v188) >> (uint(int32(8)) % 32))
						v193 = int32(16711935)
						*(*int32)(unsafe.Add(mBase, uint32(v117+v185))) = (v190&int32(255)+v188&v193+v190<<(uint(int32(16))%32))&v193 | v188&int32(-16711936)
						v207 = v173 | int32(1)
					}
					if v173 == v116+int32(-1) {
					} else {
						v214 = v207 << (uint(int32(2)) % 32)
						v218 = v116 - v207
						v222 = v115 + v214
						v223 = v117 + v214
						for {
							v227 = *(*int32)(unsafe.Add(mBase, uint32(v222)))
							v228 = int32(8)
							v229 = int32(base.Ui32(v227) >> (uint(v228) % 32))
							v230 = int32(255)
							v232 = int32(16711935)
							v235 = int32(16)
							v240 = int32(-16711936)
							*(*int32)(unsafe.Add(mBase, uint32(v223))) = (v229&v230+v227&v232+v229<<(uint(v235)%32))&v232 | v227&v240
							v244 = int32(4)
							v248 = *(*int32)(unsafe.Add(mBase, uint32(v222+v244)))
							v250 = int32(base.Ui32(v248) >> (uint(v228) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(v223+v244))) = (v250&v230+v248&v232+v250<<(uint(v235)%32))&v232 | v248&v240
							v270 = v218 + int32(-2)
							if v270 != 0 {
								v218 = v270
								v222 = v222 + v228
								v223 = v223 + v228
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v134 = v116 & int32(2147483644)
					v139 = v134
					v140 = v117
					v141 = v115
					for {
						v145 = int32(0)
						v146 = base.Simd_g_v128_load(m, v141, v145)
						v148 = base.Simd_g_i32x4_shr_u(v146, int32(8))
						v151 = base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k3)
						v154 = int32(16)
						v160 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v148, base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k4)), base.Simd_g_v128_and(v146, v151)), base.Simd_g_i32x4_shl(v148, v154)), v151), base.Simd_g_v128_and(v146, base.Simd_g_const(&F_AddGreenToBlueAndRed_SSE2__k5)))
						base.Simd_g_v128_store(m, v140, v145, v160)
						v168 = v139 + int32(-4)
						if v168 != 0 {
							v139 = v168
							v140 = v140 + v154
							v141 = v141 + v154
							continue
						} else {
							break
						}
						break
					}
					if v134 == v116 {
					} else {
						v173 = v134
						if v116&int32(1) == int32(0) {
							v207 = v173
						} else {
							v185 = v173 << (uint(int32(2)) % 32)
							v188 = *(*int32)(unsafe.Add(mBase, uint32(v115+v185)))
							v190 = int32(base.Ui32(v188) >> (uint(int32(8)) % 32))
							v193 = int32(16711935)
							*(*int32)(unsafe.Add(mBase, uint32(v117+v185))) = (v190&int32(255)+v188&v193+v190<<(uint(int32(16))%32))&v193 | v188&int32(-16711936)
							v207 = v173 | int32(1)
						}
						if v173 == v116+int32(-1) {
						} else {
							v214 = v207 << (uint(int32(2)) % 32)
							v218 = v116 - v207
							v222 = v115 + v214
							v223 = v117 + v214
							for {
								v227 = *(*int32)(unsafe.Add(mBase, uint32(v222)))
								v228 = int32(8)
								v229 = int32(base.Ui32(v227) >> (uint(v228) % 32))
								v230 = int32(255)
								v232 = int32(16711935)
								v235 = int32(16)
								v240 = int32(-16711936)
								*(*int32)(unsafe.Add(mBase, uint32(v223))) = (v229&v230+v227&v232+v229<<(uint(v235)%32))&v232 | v227&v240
								v244 = int32(4)
								v248 = *(*int32)(unsafe.Add(mBase, uint32(v222+v244)))
								v250 = int32(base.Ui32(v248) >> (uint(v228) % 32))
								*(*int32)(unsafe.Add(mBase, uint32(v223+v244))) = (v250&v230+v248&v232+v250<<(uint(v235)%32))&v232 | v248&v240
								v270 = v218 + int32(-2)
								if v270 != 0 {
									v218 = v270
									v222 = v222 + v228
									v223 = v223 + v228
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

var F_AddGreenToBlueAndRed_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_AddGreenToBlueAndRed_SSE2__k1 = [2]uint64{0x504050401000100, 0xf0e0d0c0b0a0908}
var F_AddGreenToBlueAndRed_SSE2__k2 = [2]uint64{0x706050403020100, 0xd0c0d0c09080908}
var F_AddGreenToBlueAndRed_SSE2__k3 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_AddGreenToBlueAndRed_SSE2__k4 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_AddGreenToBlueAndRed_SSE2__k5 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}

func F_AddVectorEq_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 base.V128
	_ = v23
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v29 int32
	_ = v29
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v45 int32
	_ = v45
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v64 int32
	_ = v64
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v39 = int32(0)
			v45 = v39 << (uint(int32(2)) % 32)
			v52 = l2 - v39
			v53 = l1 + v45
			v54 = l0 + v45
			for {
				v55 = *(*int32)(unsafe.Add(mBase, uint32(v53)))
				v56 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
				*(*int32)(unsafe.Add(mBase, uint32(v53))) = v55 + v56
				v59 = int32(4)
				v64 = v52 + int32(-1)
				if v64 != 0 {
					v52 = v64
					v53 = v53 + v59
					v54 = v54 + v59
					continue
				} else {
					break
				}
				break
			}
		} else {
			v14 = l2 & int32(2147483644)
			v19 = v14
			v20 = l1
			v21 = l0
			for {
				v22 = int32(0)
				v23 = base.Simd_g_v128_load(m, v20, v22)
				v25 = base.Simd_g_v128_load(m, v21, v22)
				v26 = base.Simd_g_i32x4_add(v23, v25)
				base.Simd_g_v128_store(m, v20, v22, v26)
				v29 = int32(16)
				v34 = v19 + int32(-4)
				if v34 != 0 {
					v19 = v34
					v20 = v20 + v29
					v21 = v21 + v29
					continue
				} else {
					break
				}
				break
			}
			if v14 == l2 {
			} else {
				v39 = v14
				v45 = v39 << (uint(int32(2)) % 32)
				v52 = l2 - v39
				v53 = l1 + v45
				v54 = l0 + v45
				for {
					v55 = *(*int32)(unsafe.Add(mBase, uint32(v53)))
					v56 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
					*(*int32)(unsafe.Add(mBase, uint32(v53))) = v55 + v56
					v59 = int32(4)
					v64 = v52 + int32(-1)
					if v64 != 0 {
						v52 = v64
						v53 = v53 + v59
						v54 = v54 + v59
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}
func F_AddVectorEq_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	var v4 int32
	_ = v4
	var v12 int32
	_ = v12
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v28 int32
	_ = v28
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v37 base.V128
	_ = v37
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v48 base.V128
	_ = v48
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 base.V128
	_ = v59
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v80 base.V128
	_ = v80
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v85 base.V128
	_ = v85
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v102 int32
	_ = v102
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 base.V128
	_ = v115
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 base.V128
	_ = v126
	var v129 base.V128
	_ = v129
	v4 = int32(0)
	v12 = l2 & int32(-16)
	v19 = v4
	v20 = v4
	for {
		v25 = l1 + v19
		v26 = int32(0)
		v27 = base.Simd_g_v128_load(m, v25, v26)
		v28 = l0 + v19
		v30 = base.Simd_g_v128_load(m, v28, v26)
		v31 = base.Simd_g_i32x4_add(v27, v30)
		base.Simd_g_v128_store(m, v25, v26, v31)
		v34 = int32(48)
		v35 = v25 + v34
		v37 = base.Simd_g_v128_load(m, v35, v26)
		v41 = base.Simd_g_v128_load(m, v28+v34, v26)
		v42 = base.Simd_g_i32x4_add(v37, v41)
		base.Simd_g_v128_store(m, v35, v26, v42)
		v45 = int32(32)
		v46 = v25 + v45
		v48 = base.Simd_g_v128_load(m, v46, v26)
		v52 = base.Simd_g_v128_load(m, v28+v45, v26)
		v53 = base.Simd_g_i32x4_add(v48, v52)
		base.Simd_g_v128_store(m, v46, v26, v53)
		v56 = int32(16)
		v57 = v25 + v56
		v59 = base.Simd_g_v128_load(m, v57, v26)
		v63 = base.Simd_g_v128_load(m, v28+v56, v26)
		v64 = base.Simd_g_i32x4_add(v59, v63)
		base.Simd_g_v128_store(m, v57, v26, v64)
		v68 = v19 + int32(64)
		v70 = v20 + v56
		if v12 != v70 {
			v19 = v68
			v20 = v70
			continue
		} else {
			break
		}
		break
	}
	if l2&int32(8) == int32(0) {
		v102 = v12
	} else {
		v77 = int32(16)
		v78 = l1 + v68 + v77
		v79 = int32(0)
		v80 = base.Simd_g_v128_load(m, v78, v79)
		v82 = v12 << (uint(int32(2)) % 32)
		v83 = l1 + v82
		v85 = base.Simd_g_v128_load(m, v83, v79)
		v88 = base.Simd_g_v128_load(m, l0+v82, v79)
		v89 = base.Simd_g_i32x4_add(v85, v88)
		base.Simd_g_v128_store(m, v83, v79, v89)
		v96 = base.Simd_g_v128_load(m, l0+v68+v77, v79)
		v97 = base.Simd_g_i32x4_add(v80, v96)
		base.Simd_g_v128_store(m, v78, v79, v97)
		v102 = v20 + int32(24)
	}
	switch l2&int32(7) + int32(-2) {
	case 0:
		v123 = v102 << (uint(int32(2)) % 32)
		v124 = l1 + v123
		v125 = int32(0)
		v126 = base.Simd_g_v128_load64_splat(m, v124, v125)
		v129 = base.Simd_g_v128_load64_splat(m, l0+v123, v125)
		base.Simd_g_v128_store64_lane_l0(m, v124, v125, base.Simd_g_i32x4_add(v126, v129))
		return
	default:
		return
	case 2:
		v112 = v102 << (uint(int32(2)) % 32)
		v113 = l1 + v112
		v114 = int32(0)
		v115 = base.Simd_g_v128_load(m, v113, v114)
		v118 = base.Simd_g_v128_load(m, l0+v112, v114)
		v119 = base.Simd_g_i32x4_add(v115, v118)
		base.Simd_g_v128_store(m, v113, v114, v119)
		return
	}
}
func F_AddVector_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v33 int32
	_ = v33
	var v40 int32
	_ = v40
	var v46 int32
	_ = v46
	var v53 int32
	_ = v53
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	if l3 < int32(1) {
	} else {
		if base.Ui32(l3) < base.Ui32(int32(4)) {
			v46 = int32(0)
			v53 = v46 << (uint(int32(2)) % 32)
			v62 = l3 - v46
			v63 = l0 + v53
			v64 = l1 + v53
			v65 = l2 + v53
			for {
				v66 = *(*int32)(unsafe.Add(mBase, uint32(v64)))
				v67 = *(*int32)(unsafe.Add(mBase, uint32(v63)))
				*(*int32)(unsafe.Add(mBase, uint32(v65))) = v66 + v67
				v70 = int32(4)
				v77 = v62 + int32(-1)
				if v77 != 0 {
					v62 = v77
					v63 = v63 + v70
					v64 = v64 + v70
					v65 = v65 + v70
					continue
				} else {
					break
				}
				break
			}
		} else {
			v16 = l3 & int32(2147483644)
			v22 = v16
			v23 = l2
			v24 = l1
			v25 = l0
			for {
				v26 = int32(0)
				v27 = base.Simd_g_v128_load(m, v24, v26)
				v29 = base.Simd_g_v128_load(m, v25, v26)
				v30 = base.Simd_g_i32x4_add(v27, v29)
				base.Simd_g_v128_store(m, v23, v26, v30)
				v33 = int32(16)
				v40 = v22 + int32(-4)
				if v40 != 0 {
					v22 = v40
					v23 = v23 + v33
					v24 = v24 + v33
					v25 = v25 + v33
					continue
				} else {
					break
				}
				break
			}
			if v16 == l3 {
			} else {
				v46 = v16
				v53 = v46 << (uint(int32(2)) % 32)
				v62 = l3 - v46
				v63 = l0 + v53
				v64 = l1 + v53
				v65 = l2 + v53
				for {
					v66 = *(*int32)(unsafe.Add(mBase, uint32(v64)))
					v67 = *(*int32)(unsafe.Add(mBase, uint32(v63)))
					*(*int32)(unsafe.Add(mBase, uint32(v65))) = v66 + v67
					v70 = int32(4)
					v77 = v62 + int32(-1)
					if v77 != 0 {
						v62 = v77
						v63 = v63 + v70
						v64 = v64 + v70
						v65 = v65 + v70
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}
func F_AddVector_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	var v5 int32
	_ = v5
	var v12 int32
	_ = v12
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 base.V128
	_ = v28
	var v29 int32
	_ = v29
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v35 int32
	_ = v35
	var v40 base.V128
	_ = v40
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v48 int32
	_ = v48
	var v53 base.V128
	_ = v53
	var v57 base.V128
	_ = v57
	var v58 base.V128
	_ = v58
	var v61 int32
	_ = v61
	var v66 base.V128
	_ = v66
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v88 base.V128
	_ = v88
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v96 int32
	_ = v96
	var v102 base.V128
	_ = v102
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v113 int32
	_ = v113
	var v120 int32
	_ = v120
	var v123 int32
	_ = v123
	var v124 base.V128
	_ = v124
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v136 base.V128
	_ = v136
	var v139 base.V128
	_ = v139
	v5 = int32(0)
	v12 = l3 & int32(-16)
	v20 = v5
	v21 = v5
	for {
		v25 = l2 + v20
		v26 = l1 + v20
		v27 = int32(0)
		v28 = base.Simd_g_v128_load(m, v26, v27)
		v29 = l0 + v20
		v31 = base.Simd_g_v128_load(m, v29, v27)
		v32 = base.Simd_g_i32x4_add(v28, v31)
		base.Simd_g_v128_store(m, v25, v27, v32)
		v35 = int32(48)
		v40 = base.Simd_g_v128_load(m, v26+v35, v27)
		v44 = base.Simd_g_v128_load(m, v29+v35, v27)
		v45 = base.Simd_g_i32x4_add(v40, v44)
		base.Simd_g_v128_store(m, v25+v35, v27, v45)
		v48 = int32(32)
		v53 = base.Simd_g_v128_load(m, v26+v48, v27)
		v57 = base.Simd_g_v128_load(m, v29+v48, v27)
		v58 = base.Simd_g_i32x4_add(v53, v57)
		base.Simd_g_v128_store(m, v25+v48, v27, v58)
		v61 = int32(16)
		v66 = base.Simd_g_v128_load(m, v26+v61, v27)
		v70 = base.Simd_g_v128_load(m, v29+v61, v27)
		v71 = base.Simd_g_i32x4_add(v66, v70)
		base.Simd_g_v128_store(m, v25+v61, v27, v71)
		v75 = v20 + int32(64)
		v77 = v21 + v61
		if v12 != v77 {
			v20 = v75
			v21 = v77
			continue
		} else {
			break
		}
		break
	}
	if l3&int32(8) == int32(0) {
		v113 = v12
	} else {
		v84 = v12 << (uint(int32(2)) % 32)
		v87 = int32(0)
		v88 = base.Simd_g_v128_load(m, l1+v84, v87)
		v91 = base.Simd_g_v128_load(m, l0+v84, v87)
		v92 = base.Simd_g_i32x4_add(v88, v91)
		base.Simd_g_v128_store(m, l2+v84, v87, v92)
		v96 = int32(16)
		v102 = base.Simd_g_v128_load(m, l1+v75+v96, v87)
		v107 = base.Simd_g_v128_load(m, l0+v75+v96, v87)
		v108 = base.Simd_g_i32x4_add(v102, v107)
		base.Simd_g_v128_store(m, l2+v75+v96, v87, v108)
		v113 = v21 + int32(24)
	}
	switch l3&int32(7) + int32(-2) {
	case 0:
		v132 = v113 << (uint(int32(2)) % 32)
		v135 = int32(0)
		v136 = base.Simd_g_v128_load64_splat(m, l1+v132, v135)
		v139 = base.Simd_g_v128_load64_splat(m, l0+v132, v135)
		base.Simd_g_v128_store64_lane_l0(m, l2+v132, v135, base.Simd_g_i32x4_add(v136, v139))
		return
	default:
		return
	case 2:
		v120 = v113 << (uint(int32(2)) % 32)
		v123 = int32(0)
		v124 = base.Simd_g_v128_load(m, l1+v120, v123)
		v127 = base.Simd_g_v128_load(m, l0+v120, v123)
		v128 = base.Simd_g_i32x4_add(v124, v127)
		base.Simd_g_v128_store(m, l2+v120, v123, v128)
		return
	}
}
func F_AlphaReplace_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v23 base.V128
	_ = v23
	var v25 base.V128
	_ = v25
	var v33 int32
	_ = v33
	var v63 int32
	_ = v63
	var v68 int32
	_ = v68
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v90 int32
	_ = v90
	if l1 < int32(1) {
	} else {
		if base.Ui32(l1) < base.Ui32(int32(4)) {
			v68 = int32(0)
			v80 = l1 - v68
			v81 = l0 + v68<<(uint(int32(2))%32)
			for {
				v83 = *(*int32)(unsafe.Add(mBase, uint32(v81)))
				if base.Ui32(int32(16777215)) < base.Ui32(v83) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v81))) = l2
				}
				v90 = v80 + int32(-1)
				if v90 != 0 {
					v80 = v90
					v81 = v81 + int32(4)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v14 = l1 & int32(2147483644)
			v19 = v14
			v20 = l0
			for {
				v22 = int32(0)
				v23 = base.Simd_g_v128_load(m, v20, v22)
				v25 = base.Simd_g_i32x4_lt_u(v23, base.Simd_g_const(&F_AlphaReplace_C__k0))
				if base.Simd_g_i32x4_extract_lane_l0(v25)&int32(1) == v22 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v20))) = l2
				}
				v33 = int32(1)
				if base.Simd_g_i32x4_extract_lane_l1(v25)&v33 == int32(0) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v20+int32(4)))) = l2
				}
				if base.Simd_g_i32x4_extract_lane_l2(v25)&int32(1) == int32(0) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v20+int32(8)))) = l2
				}
				if base.Simd_g_i32x4_extract_lane_l3(v25)&int32(1) == int32(0) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v20+int32(12)))) = l2
				}
				v63 = v19 + int32(-4)
				if v63 != 0 {
					v19 = v63
					v20 = v20 + int32(16)
					continue
				} else {
					break
				}
				break
			}
			if v14 == l1 {
			} else {
				v68 = v14
				v80 = l1 - v68
				v81 = l0 + v68<<(uint(int32(2))%32)
				for {
					v83 = *(*int32)(unsafe.Add(mBase, uint32(v81)))
					if base.Ui32(int32(16777215)) < base.Ui32(v83) {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v81))) = l2
					}
					v90 = v80 + int32(-1)
					if v90 != 0 {
						v80 = v90
						v81 = v81 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}

var F_AlphaReplace_C__k0 = [2]uint64{0x100000001000000, 0x100000001000000}

func F_AlphaReplace_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v16 base.V128
	_ = v16
	var v18 int32
	_ = v18
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v56 int32
	_ = v56
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v64 int32
	_ = v64
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v102 int32
	_ = v102
	var v103 base.V128
	_ = v103
	var v105 base.V128
	_ = v105
	var v112 int32
	_ = v112
	var v120 int32
	_ = v120
	var v127 int32
	_ = v127
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v144 base.V128
	_ = v144
	var v146 base.V128
	_ = v146
	var v154 int32
	_ = v154
	var v184 int32
	_ = v184
	var v191 int32
	_ = v191
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v213 int32
	_ = v213
	var v220 int32
	_ = v220
	v12 = int32(8)
	if l1 < v12 {
		v112 = int32(0)
	} else {
		v16 = base.Simd_g_i32x4_splat(l2)
		v18 = l1 + int32(-8)
		if base.Ui32(int32(8)) <= base.Ui32(v18) {
			v32 = int32(8)
			v33 = (int32(base.Ui32(v18)>>(uint(int32(3))%32)) + int32(1)) & int32(1073741822)
			v36 = l0
			for {
				v40 = int32(0)
				v41 = base.Simd_g_v128_load(m, v36, v40)
				v42 = base.Simd_g_const(&F_AlphaReplace_SSE2__k0)
				v44 = base.Simd_g_v128_bitselect(v41, v16, base.Simd_g_i32x4_gt_u(v41, v42))
				base.Simd_g_v128_store(m, v36, v40, v44)
				v47 = int32(16)
				v48 = v36 + v47
				v50 = base.Simd_g_v128_load(m, v48, v40)
				v52 = base.Simd_g_v128_bitselect(v50, v16, base.Simd_g_i32x4_gt_u(v50, v42))
				base.Simd_g_v128_store(m, v48, v40, v52)
				v56 = v36 + int32(32)
				v58 = base.Simd_g_v128_load(m, v56, v40)
				v60 = base.Simd_g_v128_bitselect(v58, v16, base.Simd_g_i32x4_gt_u(v58, v42))
				base.Simd_g_v128_store(m, v56, v40, v60)
				v64 = v36 + int32(48)
				v66 = base.Simd_g_v128_load(m, v64, v40)
				v68 = base.Simd_g_v128_bitselect(v66, v16, base.Simd_g_i32x4_gt_u(v66, v42))
				base.Simd_g_v128_store(m, v64, v40, v68)
				v74 = v32 + v47
				v76 = v33 + int32(-2)
				if v76 != 0 {
					v32 = v74
					v33 = v76
					v36 = v36 + int32(64)
					continue
				} else {
					break
				}
				break
			}
			v82 = v74
			v83 = v32 + int32(8)
		} else {
			v82 = v12
			v83 = int32(0)
		}
		if v18&int32(8) != 0 {
			v112 = v83
		} else {
			v94 = l0 + v83<<(uint(int32(2))%32)
			v95 = int32(16)
			v96 = base.Simd_g_v128_load(m, v94, v95)
			v97 = base.Simd_g_const(&F_AlphaReplace_SSE2__k0)
			v99 = base.Simd_g_v128_bitselect(v96, v16, base.Simd_g_i32x4_gt_u(v96, v97))
			base.Simd_g_v128_store(m, v94, v95, v99)
			v102 = int32(0)
			v103 = base.Simd_g_v128_load(m, v94, v102)
			v105 = base.Simd_g_v128_bitselect(v103, v16, base.Simd_g_i32x4_gt_u(v103, v97))
			base.Simd_g_v128_store(m, v94, v102, v105)
			v112 = v82
		}
	}
	if l1 <= v112 {
	} else {
		v120 = l1 - v112
		if base.Ui32(v120) <= base.Ui32(int32(3)) {
			v191 = v112
			v205 = l1 - v191
			v209 = l0 + v191<<(uint(int32(2))%32)
			for {
				v213 = *(*int32)(unsafe.Add(mBase, uint32(v209)))
				if base.Ui32(int32(16777215)) < base.Ui32(v213) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v209))) = l2
				}
				v220 = v205 + int32(-1)
				if v220 != 0 {
					v205 = v220
					v209 = v209 + int32(4)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v127 = l1 & int32(3)
			v135 = v112 + v127 - l1
			v139 = l0 + v112<<(uint(int32(2))%32)
			for {
				v143 = int32(0)
				v144 = base.Simd_g_v128_load(m, v139, v143)
				v146 = base.Simd_g_i32x4_lt_u(v144, base.Simd_g_const(&F_AlphaReplace_SSE2__k1))
				if base.Simd_g_i32x4_extract_lane_l0(v146)&int32(1) == v143 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v139))) = l2
				}
				v154 = int32(1)
				if base.Simd_g_i32x4_extract_lane_l1(v146)&v154 == int32(0) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v139+int32(4)))) = l2
				}
				if base.Simd_g_i32x4_extract_lane_l2(v146)&int32(1) == int32(0) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v139+int32(8)))) = l2
				}
				if base.Simd_g_i32x4_extract_lane_l3(v146)&int32(1) == int32(0) {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v139+int32(12)))) = l2
				}
				v184 = v135 + int32(4)
				if v184 != 0 {
					v135 = v184
					v139 = v139 + int32(16)
					continue
				} else {
					break
				}
				break
			}
			if v127 == int32(0) {
			} else {
				v191 = v112 + (v120 - v127)
				v205 = l1 - v191
				v209 = l0 + v191<<(uint(int32(2))%32)
				for {
					v213 = *(*int32)(unsafe.Add(mBase, uint32(v209)))
					if base.Ui32(int32(16777215)) < base.Ui32(v213) {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v209))) = l2
					}
					v220 = v205 + int32(-1)
					if v220 != 0 {
						v205 = v220
						v209 = v209 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	return
}

var F_AlphaReplace_SSE2__k0 = [2]uint64{0xffffff00ffffff, 0xffffff00ffffff}
var F_AlphaReplace_SSE2__k1 = [2]uint64{0x100000001000000, 0x100000001000000}

func F_ApplyAlphaMultiply_16b_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v126 int32
	_ = v126
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v129 base.V128
	_ = v129
	var v130 base.V128
	_ = v130
	var v173 int32
	_ = v173
	var v174 base.V128
	_ = v174
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v179 base.V128
	_ = v179
	var v180 base.V128
	_ = v180
	var v181 base.V128
	_ = v181
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v195 base.V128
	_ = v195
	var v196 base.V128
	_ = v196
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v208 int32
	_ = v208
	var v210 base.V128
	_ = v210
	var v211 base.V128
	_ = v211
	var v214 int32
	_ = v214
	var v217 int32
	_ = v217
	var v220 int32
	_ = v220
	var v223 int32
	_ = v223
	var v224 base.V128
	_ = v224
	var v227 int32
	_ = v227
	var v230 int32
	_ = v230
	var v233 int32
	_ = v233
	var v236 int32
	_ = v236
	var v238 base.V128
	_ = v238
	var v241 base.V128
	_ = v241
	var v244 base.V128
	_ = v244
	var v247 base.V128
	_ = v247
	var v249 int32
	_ = v249
	var v250 base.V128
	_ = v250
	var v252 int32
	_ = v252
	var v253 base.V128
	_ = v253
	var v255 int32
	_ = v255
	var v256 base.V128
	_ = v256
	var v258 int32
	_ = v258
	var v259 base.V128
	_ = v259
	var v261 int32
	_ = v261
	var v262 base.V128
	_ = v262
	var v264 int32
	_ = v264
	var v265 base.V128
	_ = v265
	var v267 int32
	_ = v267
	var v268 base.V128
	_ = v268
	var v270 int32
	_ = v270
	var v271 base.V128
	_ = v271
	var v273 int32
	_ = v273
	var v274 base.V128
	_ = v274
	var v276 int32
	_ = v276
	var v277 base.V128
	_ = v277
	var v279 int32
	_ = v279
	var v280 base.V128
	_ = v280
	var v282 int32
	_ = v282
	var v283 base.V128
	_ = v283
	var v284 base.V128
	_ = v284
	var v285 base.V128
	_ = v285
	var v288 base.V128
	_ = v288
	var v289 base.V128
	_ = v289
	var v292 int32
	_ = v292
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v304 int32
	_ = v304
	var v307 int32
	_ = v307
	var v310 int32
	_ = v310
	var v313 int32
	_ = v313
	var v316 int32
	_ = v316
	var v319 int32
	_ = v319
	var v322 int32
	_ = v322
	var v325 int32
	_ = v325
	var v328 int32
	_ = v328
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v336 base.V128
	_ = v336
	var v339 base.V128
	_ = v339
	var v342 base.V128
	_ = v342
	var v345 base.V128
	_ = v345
	var v348 base.V128
	_ = v348
	var v351 base.V128
	_ = v351
	var v354 base.V128
	_ = v354
	var v357 base.V128
	_ = v357
	var v360 base.V128
	_ = v360
	var v363 base.V128
	_ = v363
	var v366 base.V128
	_ = v366
	var v369 base.V128
	_ = v369
	var v372 base.V128
	_ = v372
	var v375 base.V128
	_ = v375
	var v378 base.V128
	_ = v378
	var v381 base.V128
	_ = v381
	var v382 base.V128
	_ = v382
	var v386 base.V128
	_ = v386
	var v390 int32
	_ = v390
	var v392 base.V128
	_ = v392
	var v396 base.V128
	_ = v396
	var v405 base.V128
	_ = v405
	var v409 base.V128
	_ = v409
	var v417 base.V128
	_ = v417
	var v421 base.V128
	_ = v421
	var v435 base.V128
	_ = v435
	var v439 int32
	_ = v439
	var v441 base.V128
	_ = v441
	var v470 base.V128
	_ = v470
	var v477 base.V128
	_ = v477
	var v508 base.V128
	_ = v508
	var v602 base.V128
	_ = v602
	var v608 int32
	_ = v608
	var v616 int32
	_ = v616
	var v678 int32
	_ = v678
	var v684 int32
	_ = v684
	var v719 int32
	_ = v719
	var v721 int32
	_ = v721
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v724 int32
	_ = v724
	var v726 int32
	_ = v726
	var v727 int32
	_ = v727
	var v729 int32
	_ = v729
	var v733 int32
	_ = v733
	var v737 int32
	_ = v737
	var v739 int32
	_ = v739
	var v760 int32
	_ = v760
	var v765 int32
	_ = v765
	if l1 < int32(1) {
	} else {
		if l2 < int32(1) {
		} else {
			v59 = l1 & int32(2147483632)
			v62 = l0
			v64 = l2
			for {
				if base.Ui32(l1) < base.Ui32(int32(16)) {
					v616 = int32(0)
					v678 = v616 << (uint(int32(1)) % 32)
					v684 = l1 - v616
					for {
						v719 = v62 + v678
						v721 = v719 + int32(1)
						v722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v721))))
						v723 = int32(15)
						v724 = v722 & v723
						v726 = v724 * int32(_a_F_ApplyAlphaMultiply_16b_C_0)
						v727 = int32(240)
						v729 = int32(4)
						v733 = int32(16)
						v737 = int32(base.Ui32(v726*(v722&v727|int32(base.Ui32(v722)>>(uint(v729)%32))))>>(uint(v733)%32))&v727 | v724
						*(*uint8)(unsafe.Add(mBase, uint32(v721))) = uint8(v737)
						v739 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v719))))
						v760 = int32(base.Ui32(v726*(v739&v727|int32(base.Ui32(v739)>>(uint(v729)%32))))>>(uint(v733)%32))&v727 | int32(base.Ui32(v726*((v739&v723|v739<<(uint(v729)%32))&int32(255)))>>(uint(int32(20))%32))
						*(*uint8)(unsafe.Add(mBase, uint32(v719))) = uint8(v760)
						v765 = v684 + int32(-1)
						if v765 != 0 {
							v678 = v678 + int32(2)
							v684 = v765
							continue
						} else {
							break
						}
						break
					}
				} else {
					v126 = v59
					v127 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k0)
					v128 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k1)
					v129 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k2)
					v130 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k3)
					for {
						v173 = int32(1)
						v174 = base.Simd_g_i32x4_shl(v130, v173)
						v175 = int32(0)
						v177 = v62 + base.Simd_g_i32x4_extract_lane_l0(v174)
						v179 = base.Simd_g_i32x4_shl(v127, v173)
						v180 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k4)
						v181 = base.Simd_g_v128_or(v179, v180)
						v182 = int32(3)
						v184 = v62 + base.Simd_g_i32x4_extract_lane_l3(v181)
						v185 = int32(2)
						v187 = v62 + base.Simd_g_i32x4_extract_lane_l2(v181)
						v190 = v62 + base.Simd_g_i32x4_extract_lane_l1(v181)
						v193 = v62 + base.Simd_g_i32x4_extract_lane_l0(v181)
						v195 = base.Simd_g_i32x4_shl(v128, v173)
						v196 = base.Simd_g_v128_or(v195, v180)
						v199 = v62 + base.Simd_g_i32x4_extract_lane_l3(v196)
						v202 = v62 + base.Simd_g_i32x4_extract_lane_l2(v196)
						v205 = v62 + base.Simd_g_i32x4_extract_lane_l1(v196)
						v208 = v62 + base.Simd_g_i32x4_extract_lane_l0(v196)
						v210 = base.Simd_g_i32x4_shl(v129, v173)
						v211 = base.Simd_g_v128_or(v210, v180)
						v214 = v62 + base.Simd_g_i32x4_extract_lane_l3(v211)
						v217 = v62 + base.Simd_g_i32x4_extract_lane_l2(v211)
						v220 = v62 + base.Simd_g_i32x4_extract_lane_l1(v211)
						v223 = v62 + base.Simd_g_i32x4_extract_lane_l0(v211)
						v224 = base.Simd_g_v128_or(v174, v180)
						v227 = v62 + base.Simd_g_i32x4_extract_lane_l3(v224)
						v230 = v62 + base.Simd_g_i32x4_extract_lane_l2(v224)
						v233 = v62 + base.Simd_g_i32x4_extract_lane_l1(v224)
						v236 = v62 + base.Simd_g_i32x4_extract_lane_l0(v224)
						v238 = base.Simd_g_v128_load8_splat(m, v236, v175)
						v241 = base.Simd_g_v128_load8_lane_l1(m, v233, v175, v238)
						v244 = base.Simd_g_v128_load8_lane_l2(m, v230, v175, v241)
						v247 = base.Simd_g_v128_load8_lane_l3(m, v227, v175, v244)
						v249 = int32(4)
						v250 = base.Simd_g_v128_load8_lane_l4(m, v223, v175, v247)
						v252 = int32(5)
						v253 = base.Simd_g_v128_load8_lane_l5(m, v220, v175, v250)
						v255 = int32(6)
						v256 = base.Simd_g_v128_load8_lane_l6(m, v217, v175, v253)
						v258 = int32(7)
						v259 = base.Simd_g_v128_load8_lane_l7(m, v214, v175, v256)
						v261 = int32(8)
						v262 = base.Simd_g_v128_load8_lane_l8(m, v208, v175, v259)
						v264 = int32(9)
						v265 = base.Simd_g_v128_load8_lane_l9(m, v205, v175, v262)
						v267 = int32(10)
						v268 = base.Simd_g_v128_load8_lane_l10(m, v202, v175, v265)
						v270 = int32(11)
						v271 = base.Simd_g_v128_load8_lane_l11(m, v199, v175, v268)
						v273 = int32(12)
						v274 = base.Simd_g_v128_load8_lane_l12(m, v193, v175, v271)
						v276 = int32(13)
						v277 = base.Simd_g_v128_load8_lane_l13(m, v190, v175, v274)
						v279 = int32(14)
						v280 = base.Simd_g_v128_load8_lane_l14(m, v187, v175, v277)
						v282 = int32(15)
						v283 = base.Simd_g_v128_load8_lane_l15(m, v184, v175, v280)
						v284 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k5)
						v285 = base.Simd_g_v128_and(v283, v284)
						v288 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k6)
						v289 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v285)), v288)
						v292 = v62 + base.Simd_g_i32x4_extract_lane_l3(v179)
						v295 = v62 + base.Simd_g_i32x4_extract_lane_l2(v179)
						v298 = v62 + base.Simd_g_i32x4_extract_lane_l1(v179)
						v301 = v62 + base.Simd_g_i32x4_extract_lane_l0(v179)
						v304 = v62 + base.Simd_g_i32x4_extract_lane_l3(v195)
						v307 = v62 + base.Simd_g_i32x4_extract_lane_l2(v195)
						v310 = v62 + base.Simd_g_i32x4_extract_lane_l1(v195)
						v313 = v62 + base.Simd_g_i32x4_extract_lane_l0(v195)
						v316 = v62 + base.Simd_g_i32x4_extract_lane_l3(v210)
						v319 = v62 + base.Simd_g_i32x4_extract_lane_l2(v210)
						v322 = v62 + base.Simd_g_i32x4_extract_lane_l1(v210)
						v325 = v62 + base.Simd_g_i32x4_extract_lane_l0(v210)
						v328 = v62 + base.Simd_g_i32x4_extract_lane_l3(v174)
						v331 = v62 + base.Simd_g_i32x4_extract_lane_l2(v174)
						v334 = v62 + base.Simd_g_i32x4_extract_lane_l1(v174)
						v336 = base.Simd_g_v128_load8_splat(m, v177, v175)
						v339 = base.Simd_g_v128_load8_lane_l1(m, v334, v175, v336)
						v342 = base.Simd_g_v128_load8_lane_l2(m, v331, v175, v339)
						v345 = base.Simd_g_v128_load8_lane_l3(m, v328, v175, v342)
						v348 = base.Simd_g_v128_load8_lane_l4(m, v325, v175, v345)
						v351 = base.Simd_g_v128_load8_lane_l5(m, v322, v175, v348)
						v354 = base.Simd_g_v128_load8_lane_l6(m, v319, v175, v351)
						v357 = base.Simd_g_v128_load8_lane_l7(m, v316, v175, v354)
						v360 = base.Simd_g_v128_load8_lane_l8(m, v313, v175, v357)
						v363 = base.Simd_g_v128_load8_lane_l9(m, v310, v175, v360)
						v366 = base.Simd_g_v128_load8_lane_l10(m, v307, v175, v363)
						v369 = base.Simd_g_v128_load8_lane_l11(m, v304, v175, v366)
						v372 = base.Simd_g_v128_load8_lane_l12(m, v301, v175, v369)
						v375 = base.Simd_g_v128_load8_lane_l13(m, v298, v175, v372)
						v378 = base.Simd_g_v128_load8_lane_l14(m, v295, v175, v375)
						v381 = base.Simd_g_v128_load8_lane_l15(m, v292, v175, v378)
						v382 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k7)
						v386 = base.Simd_g_v128_or(base.Simd_g_v128_and(v381, v382), base.Simd_g_i8x16_shr_u(v381, v249))
						v390 = int32(16)
						v392 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k8)
						v396 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v285, v381, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k8), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9)))), v288)
						v405 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k10)
						v409 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v285, v381, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9)))), v288)
						v417 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k11)
						v421 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v285, v381, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k11), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9)))), v288)
						v435 = base.Simd_g_v128_or(base.Simd_g_v128_and(v381, v284), base.Simd_g_i8x16_shl(v381, v249))
						v439 = int32(20)
						v441 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k12)
						v470 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v289, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v386))), v390), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v396, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v386, v381, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k8), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v390)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v409, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v386, v381, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v390), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v421, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v386, v381, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k11), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v390))), v382), base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v289, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v435))), v439), v441), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v396, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v435, v435, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k8), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v439), v441)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v409, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v435, v435, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v439), v441), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v421, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v435, v435, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k11), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v439), v441))))
						base.Simd_g_v128_store8_lane_l0(m, v177, v175, v470)
						v477 = base.Simd_g_v128_or(base.Simd_g_v128_and(v283, v382), base.Simd_g_i8x16_shr_u(v283, v249))
						v508 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v289, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v477))), v390), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v396, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v477, v470, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k8), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v390)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v409, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v477, v470, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v390), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_mul(v421, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_shuffle2(v477, v470, base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k11), base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k9))))), v390))), v382), v285)
						base.Simd_g_v128_store8_lane_l0(m, v236, v175, v508)
						base.Simd_g_v128_store8_lane_l1(m, v334, v175, v470)
						base.Simd_g_v128_store8_lane_l1(m, v233, v175, v508)
						base.Simd_g_v128_store8_lane_l2(m, v331, v175, v470)
						base.Simd_g_v128_store8_lane_l2(m, v230, v175, v508)
						base.Simd_g_v128_store8_lane_l3(m, v328, v175, v470)
						base.Simd_g_v128_store8_lane_l3(m, v227, v175, v508)
						base.Simd_g_v128_store8_lane_l4(m, v325, v175, v470)
						base.Simd_g_v128_store8_lane_l4(m, v223, v175, v508)
						base.Simd_g_v128_store8_lane_l5(m, v322, v175, v470)
						base.Simd_g_v128_store8_lane_l5(m, v220, v175, v508)
						base.Simd_g_v128_store8_lane_l6(m, v319, v175, v470)
						base.Simd_g_v128_store8_lane_l6(m, v217, v175, v508)
						base.Simd_g_v128_store8_lane_l7(m, v316, v175, v470)
						base.Simd_g_v128_store8_lane_l7(m, v214, v175, v508)
						base.Simd_g_v128_store8_lane_l8(m, v313, v175, v470)
						base.Simd_g_v128_store8_lane_l8(m, v208, v175, v508)
						base.Simd_g_v128_store8_lane_l9(m, v310, v175, v470)
						base.Simd_g_v128_store8_lane_l9(m, v205, v175, v508)
						base.Simd_g_v128_store8_lane_l10(m, v307, v175, v470)
						base.Simd_g_v128_store8_lane_l10(m, v202, v175, v508)
						base.Simd_g_v128_store8_lane_l11(m, v304, v175, v470)
						base.Simd_g_v128_store8_lane_l11(m, v199, v175, v508)
						base.Simd_g_v128_store8_lane_l12(m, v301, v175, v470)
						base.Simd_g_v128_store8_lane_l12(m, v193, v175, v508)
						base.Simd_g_v128_store8_lane_l13(m, v298, v175, v470)
						base.Simd_g_v128_store8_lane_l13(m, v190, v175, v508)
						base.Simd_g_v128_store8_lane_l14(m, v295, v175, v470)
						base.Simd_g_v128_store8_lane_l14(m, v187, v175, v508)
						base.Simd_g_v128_store8_lane_l15(m, v292, v175, v470)
						base.Simd_g_v128_store8_lane_l15(m, v184, v175, v508)
						v602 = base.Simd_g_const(&F_ApplyAlphaMultiply_16b_C__k13)
						v608 = v126 + int32(-16)
						if v608 != 0 {
							v126 = v608
							v127 = base.Simd_g_i32x4_add(v127, v602)
							v128 = base.Simd_g_i32x4_add(v128, v602)
							v129 = base.Simd_g_i32x4_add(v129, v602)
							v130 = base.Simd_g_i32x4_add(v130, v602)
							continue
						} else {
							break
						}
						break
					}
					if v59 == l1 {
					} else {
						v616 = v59
						v678 = v616 << (uint(int32(1)) % 32)
						v684 = l1 - v616
						for {
							v719 = v62 + v678
							v721 = v719 + int32(1)
							v722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v721))))
							v723 = int32(15)
							v724 = v722 & v723
							v726 = v724 * int32(_a_F_ApplyAlphaMultiply_16b_C_0)
							v727 = int32(240)
							v729 = int32(4)
							v733 = int32(16)
							v737 = int32(base.Ui32(v726*(v722&v727|int32(base.Ui32(v722)>>(uint(v729)%32))))>>(uint(v733)%32))&v727 | v724
							*(*uint8)(unsafe.Add(mBase, uint32(v721))) = uint8(v737)
							v739 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v719))))
							v760 = int32(base.Ui32(v726*(v739&v727|int32(base.Ui32(v739)>>(uint(v729)%32))))>>(uint(v733)%32))&v727 | int32(base.Ui32(v726*((v739&v723|v739<<(uint(v729)%32))&int32(255)))>>(uint(int32(20))%32))
							*(*uint8)(unsafe.Add(mBase, uint32(v719))) = uint8(v760)
							v765 = v684 + int32(-1)
							if v765 != 0 {
								v678 = v678 + int32(2)
								v684 = v765
								continue
							} else {
								break
							}
							break
						}
					}
				}
				if int32(1) < v64 {
					v62 = v62 + l3
					v64 = v64 + int32(-1)
					continue
				} else {
					break
				}
				break
			}
		}
	}
	return
}

var F_ApplyAlphaMultiply_16b_C__k0 = [2]uint64{0xd0000000c, 0xf0000000e}
var F_ApplyAlphaMultiply_16b_C__k1 = [2]uint64{0x900000008, 0xb0000000a}
var F_ApplyAlphaMultiply_16b_C__k2 = [2]uint64{0x500000004, 0x700000006}
var F_ApplyAlphaMultiply_16b_C__k3 = [2]uint64{0x100000000, 0x300000002}
var F_ApplyAlphaMultiply_16b_C__k4 = [2]uint64{0x100000001, 0x100000001}
var F_ApplyAlphaMultiply_16b_C__k5 = [2]uint64{0xf0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f}
var F_ApplyAlphaMultiply_16b_C__k6 = [2]uint64{0x111100001111, 0x111100001111}
var F_ApplyAlphaMultiply_16b_C__k7 = [2]uint64{0xf0f0f0f0f0f0f0f0, 0xf0f0f0f0f0f0f0f0}
var F_ApplyAlphaMultiply_16b_C__k8 = [2]uint64{0x7060504, 0x0}
var F_ApplyAlphaMultiply_16b_C__k9 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ApplyAlphaMultiply_16b_C__k10 = [2]uint64{0xb0a0908, 0x0}
var F_ApplyAlphaMultiply_16b_C__k11 = [2]uint64{0xf0e0d0c, 0x0}
var F_ApplyAlphaMultiply_16b_C__k12 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ApplyAlphaMultiply_16b_C__k13 = [2]uint64{0x1000000010, 0x1000000010}

func F_ApplyAlphaMultiply_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v75 int32
	_ = v75
	var v76 base.V128
	_ = v76
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v84 base.V128
	_ = v84
	var v86 base.V128
	_ = v86
	var v88 base.V128
	_ = v88
	var v92 base.V128
	_ = v92
	var v94 int32
	_ = v94
	var v97 base.V128
	_ = v97
	var v103 base.V128
	_ = v103
	var v112 base.V128
	_ = v112
	var v120 int32
	_ = v120
	var v128 int32
	_ = v128
	var v144 int32
	_ = v144
	var v148 int32
	_ = v148
	var v150 int32
	_ = v150
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v188 int32
	_ = v188
	var v192 int32
	_ = v192
	var v197 int32
	_ = v197
	var v227 int32
	_ = v227
	var v232 int32
	_ = v232
	var v237 int32
	_ = v237
	var v242 int32
	_ = v242
	var v249 int32
	_ = v249
	var v253 int32
	_ = v253
	var v257 int32
	_ = v257
	var v259 int32
	_ = v259
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v272 base.V128
	_ = v272
	var v273 base.V128
	_ = v273
	var v275 base.V128
	_ = v275
	var v276 base.V128
	_ = v276
	var v278 base.V128
	_ = v278
	var v280 base.V128
	_ = v280
	var v282 base.V128
	_ = v282
	var v284 base.V128
	_ = v284
	var v288 base.V128
	_ = v288
	var v290 int32
	_ = v290
	var v293 base.V128
	_ = v293
	var v299 base.V128
	_ = v299
	var v308 base.V128
	_ = v308
	var v311 int32
	_ = v311
	var v320 int32
	_ = v320
	var v326 int32
	_ = v326
	var v329 int32
	_ = v329
	var v331 int32
	_ = v331
	var v342 int32
	_ = v342
	var v346 int32
	_ = v346
	var v347 int32
	_ = v347
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	var v353 int32
	_ = v353
	var v354 int32
	_ = v354
	var v357 int32
	_ = v357
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v373 int32
	_ = v373
	var v407 int32
	_ = v407
	var v410 int32
	_ = v410
	var v426 int32
	_ = v426
	var v427 int32
	_ = v427
	var v428 int32
	_ = v428
	var v432 int32
	_ = v432
	var v433 int32
	_ = v433
	var v435 int32
	_ = v435
	var v436 int32
	_ = v436
	var v439 int32
	_ = v439
	var v440 int32
	_ = v440
	var v443 int32
	_ = v443
	var v446 int32
	_ = v446
	var v447 int32
	_ = v447
	var v450 int32
	_ = v450
	var v456 int32
	_ = v456
	var v460 int32
	_ = v460
	var v462 int32
	_ = v462
	var v463 int32
	_ = v463
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v473 int32
	_ = v473
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v480 int32
	_ = v480
	var v486 int32
	_ = v486
	var v490 int32
	_ = v490
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v495 int32
	_ = v495
	var v496 int32
	_ = v496
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v503 int32
	_ = v503
	var v506 int32
	_ = v506
	var v507 int32
	_ = v507
	var v510 int32
	_ = v510
	var v518 int32
	_ = v518
	if l3 < int32(1) {
	} else {
		v23 = l3 + int32(-1)
		if l1 != 0 {
			v26 = int32(0)
		} else {
			v26 = int32(3)
		}
		v27 = int32(0)
		v28 = base.B2i32(l1 != v27)
		if l1 == v27 {
			if l2 < int32(4) {
				if l2 < int32(1) {
				} else {
					v407 = l0
					v410 = l3
					for {
						v426 = v407 + v28
						v427 = v407 + v26
						v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v427))))
						if v428 == int32(255) {
						} else {
							v432 = v428 * int32(_a_F_ApplyAlphaMultiply_SSE2_0)
							v433 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v426))))
							v435 = int32(23)
							v436 = int32(base.Ui32(v432*v433) >> (uint(v435) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v426))) = uint8(v436)
							v439 = v426 + int32(1)
							v440 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v439))))
							v443 = int32(base.Ui32(v432*v440) >> (uint(v435) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v439))) = uint8(v443)
							v446 = v426 + int32(2)
							v447 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v446))))
							v450 = int32(base.Ui32(v432*v447) >> (uint(v435) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v446))) = uint8(v450)
						}
						if l2 == int32(1) {
						} else {
							v456 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v427+int32(4)))))
							if v456 == int32(255) {
							} else {
								v460 = v426 + int32(4)
								v462 = v456 * int32(_a_F_ApplyAlphaMultiply_SSE2_0)
								v463 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v460))))
								v465 = int32(23)
								v466 = int32(base.Ui32(v462*v463) >> (uint(v465) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v460))) = uint8(v466)
								v469 = v426 + int32(5)
								v470 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v469))))
								v473 = int32(base.Ui32(v462*v470) >> (uint(v465) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v469))) = uint8(v473)
								v476 = v426 + int32(6)
								v477 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v476))))
								v480 = int32(base.Ui32(v462*v477) >> (uint(v465) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v476))) = uint8(v480)
							}
							if l2 == int32(2) {
							} else {
								v486 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v427+int32(8)))))
								if v486 == int32(255) {
								} else {
									v490 = v426 + int32(8)
									v492 = v486 * int32(_a_F_ApplyAlphaMultiply_SSE2_0)
									v493 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v490))))
									v495 = int32(23)
									v496 = int32(base.Ui32(v492*v493) >> (uint(v495) % 32))
									*(*uint8)(unsafe.Add(mBase, uint32(v490))) = uint8(v496)
									v499 = v426 + int32(9)
									v500 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v499))))
									v503 = int32(base.Ui32(v492*v500) >> (uint(v495) % 32))
									*(*uint8)(unsafe.Add(mBase, uint32(v499))) = uint8(v503)
									v506 = v426 + int32(10)
									v507 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v506))))
									v510 = int32(base.Ui32(v492*v507) >> (uint(v495) % 32))
									*(*uint8)(unsafe.Add(mBase, uint32(v506))) = uint8(v510)
								}
							}
						}
						v518 = v410 + int32(-1)
						if int32(0) < v518 {
							v407 = v407 + l4
							v410 = v518
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				v227 = int32(16)
				v232 = l0
				v237 = v23
				v242 = v26 + l0 + v227
				v249 = v28 + l0 + v227
				for {
					v253 = v232
					v257 = v249
					v259 = v242
					v268 = int32(0)
					for {
						v271 = int32(0)
						v272 = base.Simd_g_v128_load(m, v253, v271)
						v273 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k0)
						v275 = base.Simd_g_i8x16_shuffle2(v272, v273, base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k1), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k2))
						v276 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k3)
						v278 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k4)
						v280 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k5)
						v282 = base.Simd_g_i16x8_mul(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_or(v275, v276), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k4)), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k5)), v275)
						v284 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k6)
						v288 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k7)
						v290 = int32(7)
						v293 = base.Simd_g_i8x16_shuffle2(v272, v273, base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k8), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k9))
						v299 = base.Simd_g_i16x8_mul(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_or(v293, v276), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k4)), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k5)), v293)
						v308 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v282), v284), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v282), v284), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k11)), v290), base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v299), v284), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v299), v284), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k11)), v290))
						base.Simd_g_v128_store(m, v253, v271, v308)
						v311 = int32(16)
						v320 = v268 + int32(4)
						if v268+int32(8) <= l2 {
							v253 = v253 + v311
							v257 = v257 + v311
							v259 = v259 + v311
							v268 = v320
							continue
						} else {
							break
						}
						break
					}
					if l2 <= v320 {
					} else {
						v326 = v257
						v329 = v259
						v331 = v320
						for {
							v342 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v329))))
							if v342 == int32(255) {
							} else {
								v346 = v342 * int32(_a_F_ApplyAlphaMultiply_SSE2_0)
								v347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v326))))
								v349 = int32(23)
								v350 = int32(base.Ui32(v346*v347) >> (uint(v349) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v326))) = uint8(v350)
								v353 = v326 + int32(1)
								v354 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v353))))
								v357 = int32(base.Ui32(v346*v354) >> (uint(v349) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v353))) = uint8(v357)
								v360 = v326 + int32(2)
								v361 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v360))))
								v364 = int32(base.Ui32(v346*v361) >> (uint(v349) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v360))) = uint8(v364)
							}
							v368 = int32(4)
							v373 = v331 + int32(1)
							if v373 < l2 {
								v326 = v326 + v368
								v329 = v329 + v368
								v331 = v373
								continue
							} else {
								break
							}
							break
						}
					}
					if int32(0) < v237 {
						v232 = v232 + l4
						v237 = v237 + int32(-1)
						v242 = v242 + l4
						v249 = v249 + l4
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v35 = l0
			v40 = v23
			v43 = l0 + v28
			v44 = l0 + v26
			for {
				if l2 < int32(4) {
					v128 = int32(0)
				} else {
					v57 = v35
					v59 = int32(0)
					for {
						v75 = int32(0)
						v76 = base.Simd_g_v128_load(m, v57, v75)
						v77 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k0)
						v79 = base.Simd_g_i8x16_shuffle2(v76, v77, base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k1), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k2))
						v80 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k3)
						v82 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k12)
						v84 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k13)
						v86 = base.Simd_g_i16x8_mul(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_or(v79, v80), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k12)), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k13)), v79)
						v88 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k6)
						v92 = base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k7)
						v94 = int32(7)
						v97 = base.Simd_g_i8x16_shuffle2(v76, v77, base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k8), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k9))
						v103 = base.Simd_g_i16x8_mul(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_or(v97, v80), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k12)), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k13)), v97)
						v112 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v86), v88), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v86), v88), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k11)), v94), base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v103), v88), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v103), v88), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k10), base.Simd_g_const(&F_ApplyAlphaMultiply_SSE2__k11)), v94))
						base.Simd_g_v128_store(m, v57, v75, v112)
						v120 = v59 + int32(4)
						if v59+int32(8) <= l2 {
							v57 = v57 + int32(16)
							v59 = v120
							continue
						} else {
							break
						}
						break
					}
					v128 = v120
				}
				if l2 <= v128 {
				} else {
					v144 = v128 << (uint(int32(2)) % 32)
					v148 = v43 + v144
					v150 = v44 + v144
					v163 = l2 - v128
					for {
						v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v150))))
						if v166 == int32(255) {
						} else {
							v170 = v166 * int32(_a_F_ApplyAlphaMultiply_SSE2_0)
							v171 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v148))))
							v173 = int32(23)
							v174 = int32(base.Ui32(v170*v171) >> (uint(v173) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v148))) = uint8(v174)
							v177 = v148 + int32(1)
							v178 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v177))))
							v181 = int32(base.Ui32(v170*v178) >> (uint(v173) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v177))) = uint8(v181)
							v184 = v148 + int32(2)
							v185 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v184))))
							v188 = int32(base.Ui32(v170*v185) >> (uint(v173) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v184))) = uint8(v188)
						}
						v192 = int32(4)
						v197 = v163 + int32(-1)
						if v197 != 0 {
							v148 = v148 + v192
							v150 = v150 + v192
							v163 = v197
							continue
						} else {
							break
						}
						break
					}
				}
				if int32(0) < v40 {
					v35 = v35 + l4
					v40 = v40 + int32(-1)
					v43 = v43 + l4
					v44 = v44 + l4
					continue
				} else {
					break
				}
				break
			}
		}
	}
	return
}

var F_ApplyAlphaMultiply_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_ApplyAlphaMultiply_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ApplyAlphaMultiply_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ApplyAlphaMultiply_SSE2__k3 = [2]uint64{0xff00ff0000, 0xff00ff0000}
var F_ApplyAlphaMultiply_SSE2__k4 = [2]uint64{0x504070607060706, 0xf0e0d0c0b0a0908}
var F_ApplyAlphaMultiply_SSE2__k5 = [2]uint64{0x706050403020100, 0xd0c0f0e0f0e0f0e}
var F_ApplyAlphaMultiply_SSE2__k6 = [2]uint64{0x808100008081, 0x808100008081}
var F_ApplyAlphaMultiply_SSE2__k7 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_ApplyAlphaMultiply_SSE2__k8 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ApplyAlphaMultiply_SSE2__k9 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ApplyAlphaMultiply_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_ApplyAlphaMultiply_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_ApplyAlphaMultiply_SSE2__k12 = [2]uint64{0x100010001000302, 0xf0e0d0c0b0a0908}
var F_ApplyAlphaMultiply_SSE2__k13 = [2]uint64{0x706050403020100, 0x908090809080b0a}

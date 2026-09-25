//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_CollectColorBlueTransforms_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v42 base.V128
	_ = v42
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v87 int32
	_ = v87
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v95 base.V128
	_ = v95
	var v97 int32
	_ = v97
	var v98 base.V128
	_ = v98
	var v105 int32
	_ = v105
	var v109 base.V128
	_ = v109
	var v114 base.V128
	_ = v114
	var v115 base.V128
	_ = v115
	var v123 base.V128
	_ = v123
	var v135 base.V128
	_ = v135
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v211 int32
	_ = v211
	var v215 int32
	_ = v215
	var v238 int32
	_ = v238
	var v241 int32
	_ = v241
	var v259 int32
	_ = v259
	var v262 int32
	_ = v262
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v283 int32
	_ = v283
	var v286 int32
	_ = v286
	var v289 int32
	_ = v289
	var v304 int32
	_ = v304
	var v305 int32
	_ = v305
	var v312 int32
	_ = v312
	if l3 < int32(1) {
	} else {
		if l2 < int32(8) {
		} else {
			v25 = int32(8)
			v30 = base.Simd_g_i32x4_splat(base.I32_extend16_s(l5<<(uint(v25)%32)) << (uint(int32(11)) % 32))
			v31 = base.Simd_g_i32x4_extend_high_i16x8_s(v30)
			v32 = base.Simd_g_i32x4_extend_low_i16x8_s(v30)
			v40 = base.Simd_g_i32x4_splat(int32(base.Ui32(base.I32_extend16_s(l4<<(uint(v25)%32)))>>(uint(int32(5))%32)) & int32(_a_F_CollectColorBlueTransforms_SSE2_0))
			v41 = base.Simd_g_i32x4_extend_high_i16x8_s(v40)
			v42 = base.Simd_g_i32x4_extend_low_i16x8_s(v40)
			v59 = int32(0)
			v60 = l0
			for {
				v82 = int32(8)
				v83 = v60
				for {
					v87 = int32(0)
					v88 = base.Simd_g_v128_load_rng(m, v83, v87, int32(0), int32(32))
					v89 = base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k0)
					v90 = base.Simd_g_v128_and(v88, v89)
					v95 = base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k1)
					v97 = int32(8)
					v98 = base.Simd_g_i16x8_shl(v88, v97)
					v105 = int32(16)
					v109 = base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k2)
					v114 = base.Simd_g_v128_load_nc(m, v83+v105, v87)
					v115 = base.Simd_g_v128_and(v114, v89)
					v123 = base.Simd_g_i16x8_shl(v114, v97)
					v135 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(base.Simd_g_i8x16_sub(v88, base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v90), v42), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v90), v41), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k3), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k4)), base.Simd_g_i32x4_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v98), v32), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v98), v31), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k3), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k4)), v105))), v109), base.Simd_g_v128_and(base.Simd_g_i8x16_sub(v114, base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v115), v42), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v115), v41), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k3), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k4)), base.Simd_g_i32x4_shr_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v123), v32), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v123), v31), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k3), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE2__k4)), v105))), v109))
					v138 = int32(2)
					v140 = l6 + base.Simd_g_i16x8_extract_lane_u_l0(v135)<<(uint(v138)%32)
					v141 = *(*int32)(unsafe.Add(mBase, uint32(v140)))
					v142 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v140))) = v141 + v142
					v149 = l6 + base.Simd_g_i16x8_extract_lane_u_l1(v135)<<(uint(v138)%32)
					v150 = *(*int32)(unsafe.Add(mBase, uint32(v149)))
					*(*int32)(unsafe.Add(mBase, uint32(v149))) = v150 + v142
					v158 = l6 + base.Simd_g_i16x8_extract_lane_u_l2(v135)<<(uint(v138)%32)
					v159 = *(*int32)(unsafe.Add(mBase, uint32(v158)))
					*(*int32)(unsafe.Add(mBase, uint32(v158))) = v159 + v142
					v167 = l6 + base.Simd_g_i16x8_extract_lane_u_l3(v135)<<(uint(v138)%32)
					v168 = *(*int32)(unsafe.Add(mBase, uint32(v167)))
					*(*int32)(unsafe.Add(mBase, uint32(v167))) = v168 + v142
					v176 = l6 + base.Simd_g_i16x8_extract_lane_u_l4(v135)<<(uint(v138)%32)
					v177 = *(*int32)(unsafe.Add(mBase, uint32(v176)))
					*(*int32)(unsafe.Add(mBase, uint32(v176))) = v177 + v142
					v185 = l6 + base.Simd_g_i16x8_extract_lane_u_l5(v135)<<(uint(v138)%32)
					v186 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
					*(*int32)(unsafe.Add(mBase, uint32(v185))) = v186 + v142
					v194 = l6 + base.Simd_g_i16x8_extract_lane_u_l6(v135)<<(uint(v138)%32)
					v195 = *(*int32)(unsafe.Add(mBase, uint32(v194)))
					*(*int32)(unsafe.Add(mBase, uint32(v194))) = v195 + v142
					v203 = l6 + base.Simd_g_i16x8_extract_lane_u_l7(v135)<<(uint(v138)%32)
					v204 = *(*int32)(unsafe.Add(mBase, uint32(v203)))
					*(*int32)(unsafe.Add(mBase, uint32(v203))) = v204 + v142
					v211 = v82 + v97
					if v211 <= l2 {
						v82 = v211
						v83 = v83 + int32(32)
						continue
					} else {
						break
					}
					break
				}
				v215 = v59 + int32(1)
				if v215 != l3 {
					v59 = v215
					v60 = v60 + l1<<(uint(int32(2))%32)
					continue
				} else {
					break
				}
				break
			}
		}
	}
	v238 = l2 & int32(7)
	if v238 == int32(0) {
	} else {
		v241 = int32(2)
		if l3 < int32(1) {
		} else {
			if v238 < int32(1) {
			} else {
				v259 = l0 + l2<<(uint(v241)%32) - v238<<(uint(v241)%32)
				v262 = l3
				for {
					v273 = v238
					v277 = v259
					for {
						v283 = *(*int32)(unsafe.Add(mBase, uint32(v277)))
						v286 = int32(24)
						v289 = int32(5)
						v304 = l6 + (v283-(int32(base.Ui32(v283<<(uint(int32(16))%32)>>(uint(v286)%32)*base.I32_extend8_s(l4))>>(uint(v289)%32))+int32(base.Ui32(v283<<(uint(int32(8))%32)>>(uint(v286)%32)*base.I32_extend8_s(l5))>>(uint(v289)%32))))&int32(255)<<(uint(int32(2))%32)
						v305 = *(*int32)(unsafe.Add(mBase, uint32(v304)))
						*(*int32)(unsafe.Add(mBase, uint32(v304))) = v305 + int32(1)
						v312 = v273 + int32(-1)
						if v312 != 0 {
							v273 = v312
							v277 = v277 + int32(4)
							continue
						} else {
							break
						}
						break
					}
					if int32(1) < v262 {
						v259 = v259 + l1<<(uint(int32(2))%32)
						v262 = v262 + int32(-1)
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

var F_CollectColorBlueTransforms_SSE2__k0 = [2]uint64{0xff000000ff00, 0xff000000ff00}
var F_CollectColorBlueTransforms_SSE2__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_CollectColorBlueTransforms_SSE2__k2 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_CollectColorBlueTransforms_SSE2__k3 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_CollectColorBlueTransforms_SSE2__k4 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}

func F_CollectColorBlueTransforms_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v37 base.V128
	_ = v37
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v76 base.V128
	_ = v76
	var v79 base.V128
	_ = v79
	var v90 base.V128
	_ = v90
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v137 base.V128
	_ = v137
	var v138 base.V128
	_ = v138
	var v145 base.V128
	_ = v145
	var v146 int32
	_ = v146
	var v148 base.V128
	_ = v148
	var v152 int32
	_ = v152
	var v161 base.V128
	_ = v161
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v209 int32
	_ = v209
	var v229 int32
	_ = v229
	var v232 int32
	_ = v232
	var v250 int32
	_ = v250
	var v253 int32
	_ = v253
	var v264 int32
	_ = v264
	var v268 int32
	_ = v268
	var v274 int32
	_ = v274
	var v277 int32
	_ = v277
	var v280 int32
	_ = v280
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v303 int32
	_ = v303
	if l2 < int32(4) {
	} else {
		if l3 < int32(1) {
		} else {
			v22 = int32(8)
			v37 = base.Simd_g_i32x4_splat(base.I32_extend16_s(l5<<(uint(v22)%32))<<(uint(int32(11))%32) | int32(base.Ui32(base.I32_extend16_s(l4<<(uint(v22)%32)))>>(uint(int32(5))%32))&int32(_a_F_CollectColorBlueTransforms_SSE41_0) + int32(16777216))
			v38 = base.Simd_g_i32x4_extend_high_i16x8_s(v37)
			v39 = base.Simd_g_i32x4_extend_low_i16x8_s(v37)
			v56 = l0 + int32(16)
			v57 = int32(0)
			for {
				v67 = base.Simd_g_v128_load(m, l0+v57*l1<<(uint(int32(2))%32), int32(0))
				v68 = base.Simd_g_const(&F_CollectColorBlueTransforms_SSE41__k0)
				v69 = base.Simd_g_i8x16_swizzle(v67, v68)
				v76 = base.Simd_g_i16x8_sub(v67, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v69), v39), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v69), v38), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE41__k1), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE41__k2)))
				v79 = base.Simd_g_i16x8_add(v76, base.Simd_g_i32x4_shr_u(v76, int32(16)))
				if base.Ui32(l2) < base.Ui32(int32(8)) {
					v161 = v79
				} else {
					v90 = v79
					v97 = int32(8)
					v98 = v56
					for {
						v100 = int32(0)
						v102 = int32(2)
						v104 = l6 + base.Simd_g_i8x16_extract_lane_u_l0(v90)<<(uint(v102)%32)
						v105 = *(*int32)(unsafe.Add(mBase, uint32(v104)))
						v106 = int32(1)
						*(*int32)(unsafe.Add(mBase, uint32(v104))) = v105 + v106
						v109 = int32(4)
						v113 = l6 + base.Simd_g_i8x16_extract_lane_u_l4(v90)<<(uint(v102)%32)
						v114 = *(*int32)(unsafe.Add(mBase, uint32(v113)))
						*(*int32)(unsafe.Add(mBase, uint32(v113))) = v114 + v106
						v122 = l6 + base.Simd_g_i8x16_extract_lane_u_l8(v90)<<(uint(v102)%32)
						v123 = *(*int32)(unsafe.Add(mBase, uint32(v122)))
						*(*int32)(unsafe.Add(mBase, uint32(v122))) = v123 + v106
						v131 = l6 + base.Simd_g_i8x16_extract_lane_u_l12(v90)<<(uint(v102)%32)
						v132 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
						*(*int32)(unsafe.Add(mBase, uint32(v131))) = v132 + v106
						v137 = base.Simd_g_v128_load(m, v98, v100)
						v138 = base.Simd_g_i8x16_swizzle(v137, v68)
						v145 = base.Simd_g_i16x8_sub(v137, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v138), v39), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v138), v38), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE41__k1), base.Simd_g_const(&F_CollectColorBlueTransforms_SSE41__k2)))
						v146 = int32(16)
						v148 = base.Simd_g_i16x8_add(v145, base.Simd_g_i32x4_shr_u(v145, v146))
						v152 = v97 + v109
						if v152 <= l2 {
							v90 = v148
							v97 = v152
							v98 = v98 + v146
							continue
						} else {
							break
						}
						break
					}
					v161 = v148
				}
				v173 = int32(2)
				v175 = l6 + base.Simd_g_i8x16_extract_lane_u_l0(v161)<<(uint(v173)%32)
				v176 = *(*int32)(unsafe.Add(mBase, uint32(v175)))
				v177 = int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v175))) = v176 + v177
				v184 = l6 + base.Simd_g_i8x16_extract_lane_u_l4(v161)<<(uint(v173)%32)
				v185 = *(*int32)(unsafe.Add(mBase, uint32(v184)))
				*(*int32)(unsafe.Add(mBase, uint32(v184))) = v185 + v177
				v193 = l6 + base.Simd_g_i8x16_extract_lane_u_l8(v161)<<(uint(v173)%32)
				v194 = *(*int32)(unsafe.Add(mBase, uint32(v193)))
				*(*int32)(unsafe.Add(mBase, uint32(v193))) = v194 + v177
				v202 = l6 + base.Simd_g_i8x16_extract_lane_u_l12(v161)<<(uint(v173)%32)
				v203 = *(*int32)(unsafe.Add(mBase, uint32(v202)))
				*(*int32)(unsafe.Add(mBase, uint32(v202))) = v203 + v177
				v209 = v57 + v177
				if v209 != l3 {
					v56 = v56 + l1<<(uint(int32(2))%32)
					v57 = v209
					continue
				} else {
					break
				}
				break
			}
		}
	}
	v229 = l2 & int32(3)
	if v229 == int32(0) {
	} else {
		v232 = int32(2)
		if l3 < int32(1) {
		} else {
			if v229 < int32(1) {
			} else {
				v250 = l0 + l2<<(uint(v232)%32) - v229<<(uint(v232)%32)
				v253 = l3
				for {
					v264 = v229
					v268 = v250
					for {
						v274 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
						v277 = int32(24)
						v280 = int32(5)
						v295 = l6 + (v274-(int32(base.Ui32(v274<<(uint(int32(16))%32)>>(uint(v277)%32)*base.I32_extend8_s(l4))>>(uint(v280)%32))+int32(base.Ui32(v274<<(uint(int32(8))%32)>>(uint(v277)%32)*base.I32_extend8_s(l5))>>(uint(v280)%32))))&int32(255)<<(uint(int32(2))%32)
						v296 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
						*(*int32)(unsafe.Add(mBase, uint32(v295))) = v296 + int32(1)
						v303 = v264 + int32(-1)
						if v303 != 0 {
							v264 = v303
							v268 = v268 + int32(4)
							continue
						} else {
							break
						}
						break
					}
					if int32(1) < v253 {
						v250 = v250 + l1<<(uint(int32(2))%32)
						v253 = v253 + int32(-1)
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

var F_CollectColorBlueTransforms_SSE41__k0 = [2]uint64{0x68f058f028f018f, 0xe8f0d8f0a8f098f}
var F_CollectColorBlueTransforms_SSE41__k1 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_CollectColorBlueTransforms_SSE41__k2 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}

func F_CollectColorRedTransforms_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	var v71 base.V128
	_ = v71
	var v72 int32
	_ = v72
	var v74 base.V128
	_ = v74
	var v75 base.V128
	_ = v75
	var v80 base.V128
	_ = v80
	var v83 base.V128
	_ = v83
	var v88 base.V128
	_ = v88
	var v91 base.V128
	_ = v91
	var v100 base.V128
	_ = v100
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v176 int32
	_ = v176
	var v180 int32
	_ = v180
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v218 int32
	_ = v218
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v241 int32
	_ = v241
	var v244 int32
	_ = v244
	var v251 int32
	_ = v251
	var v252 int32
	_ = v252
	var v256 int32
	_ = v256
	var v259 int32
	_ = v259
	var v262 int32
	_ = v262
	var v264 int32
	_ = v264
	var v266 int32
	_ = v266
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v273 int32
	_ = v273
	var v288 int32
	_ = v288
	var v289 int32
	_ = v289
	var v296 int32
	_ = v296
	var v299 int32
	_ = v299
	var v311 int32
	_ = v311
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v329 int32
	_ = v329
	var v330 int32
	_ = v330
	if l3 < int32(1) {
	} else {
		if l2 < int32(8) {
		} else {
			v29 = base.Simd_g_i32x4_splat(int32(base.Ui32(base.I32_extend16_s(l4<<(uint(int32(8))%32)))>>(uint(int32(5))%32)) & int32(_a_F_CollectColorRedTransforms_SSE2_0))
			v30 = base.Simd_g_i32x4_extend_high_i16x8_s(v29)
			v31 = base.Simd_g_i32x4_extend_low_i16x8_s(v29)
			v45 = int32(0)
			v46 = l0
			for {
				v65 = int32(8)
				v66 = v46
				for {
					v70 = int32(0)
					v71 = base.Simd_g_v128_load_rng(m, v66, v70, int32(0), int32(32))
					v72 = int32(16)
					v74 = base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k0)
					v75 = base.Simd_g_v128_and(v71, v74)
					v80 = base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k1)
					v83 = base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k2)
					v88 = base.Simd_g_v128_load_nc(m, v66+v72, v70)
					v91 = base.Simd_g_v128_and(v88, v74)
					v100 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(base.Simd_g_i8x16_sub(base.Simd_g_i32x4_shr_u(v71, v72), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v75), v31), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v75), v30), base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k3), base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k4))), v83), base.Simd_g_v128_and(base.Simd_g_i8x16_sub(base.Simd_g_i32x4_shr_u(v88, v72), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v91), v31), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v91), v30), base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k3), base.Simd_g_const(&F_CollectColorRedTransforms_SSE2__k4))), v83))
					v103 = int32(2)
					v105 = l5 + base.Simd_g_i16x8_extract_lane_u_l0(v100)<<(uint(v103)%32)
					v106 = *(*int32)(unsafe.Add(mBase, uint32(v105)))
					v107 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v105))) = v106 + v107
					v114 = l5 + base.Simd_g_i16x8_extract_lane_u_l1(v100)<<(uint(v103)%32)
					v115 = *(*int32)(unsafe.Add(mBase, uint32(v114)))
					*(*int32)(unsafe.Add(mBase, uint32(v114))) = v115 + v107
					v123 = l5 + base.Simd_g_i16x8_extract_lane_u_l2(v100)<<(uint(v103)%32)
					v124 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
					*(*int32)(unsafe.Add(mBase, uint32(v123))) = v124 + v107
					v132 = l5 + base.Simd_g_i16x8_extract_lane_u_l3(v100)<<(uint(v103)%32)
					v133 = *(*int32)(unsafe.Add(mBase, uint32(v132)))
					*(*int32)(unsafe.Add(mBase, uint32(v132))) = v133 + v107
					v141 = l5 + base.Simd_g_i16x8_extract_lane_u_l4(v100)<<(uint(v103)%32)
					v142 = *(*int32)(unsafe.Add(mBase, uint32(v141)))
					*(*int32)(unsafe.Add(mBase, uint32(v141))) = v142 + v107
					v150 = l5 + base.Simd_g_i16x8_extract_lane_u_l5(v100)<<(uint(v103)%32)
					v151 = *(*int32)(unsafe.Add(mBase, uint32(v150)))
					*(*int32)(unsafe.Add(mBase, uint32(v150))) = v151 + v107
					v159 = l5 + base.Simd_g_i16x8_extract_lane_u_l6(v100)<<(uint(v103)%32)
					v160 = *(*int32)(unsafe.Add(mBase, uint32(v159)))
					*(*int32)(unsafe.Add(mBase, uint32(v159))) = v160 + v107
					v168 = l5 + base.Simd_g_i16x8_extract_lane_u_l7(v100)<<(uint(v103)%32)
					v169 = *(*int32)(unsafe.Add(mBase, uint32(v168)))
					*(*int32)(unsafe.Add(mBase, uint32(v168))) = v169 + v107
					v176 = v65 + int32(8)
					if v176 <= l2 {
						v65 = v176
						v66 = v66 + int32(32)
						continue
					} else {
						break
					}
					break
				}
				v180 = v45 + int32(1)
				if v180 != l3 {
					v45 = v180
					v46 = v46 + l1<<(uint(int32(2))%32)
					continue
				} else {
					break
				}
				break
			}
		}
	}
	v200 = l2 & int32(7)
	if v200 == int32(0) {
	} else {
		v203 = int32(2)
		if l3 < int32(1) {
		} else {
			if v200 < int32(1) {
			} else {
				v218 = base.I32_extend8_s(l4)
				v225 = l0 + l2<<(uint(v203)%32) - v200<<(uint(v203)%32)
				v228 = l3
				for {
					if v200 == int32(1) {
						v299 = int32(0)
					} else {
						v241 = int32(0)
						v244 = v225
						for {
							v251 = *(*int32)(unsafe.Add(mBase, uint32(v244)))
							v252 = int32(16)
							v256 = int32(24)
							v259 = int32(5)
							v262 = int32(255)
							v264 = int32(2)
							v266 = l5 + (int32(base.Ui32(v251)>>(uint(v252)%32))-int32(base.Ui32(v251<<(uint(v252)%32)>>(uint(v256)%32)*v218)>>(uint(v259)%32)))&v262<<(uint(v264)%32)
							v267 = *(*int32)(unsafe.Add(mBase, uint32(v266)))
							v268 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(v266))) = v267 + v268
							v273 = *(*int32)(unsafe.Add(mBase, uint32(v244+int32(4))))
							v288 = l5 + (int32(base.Ui32(v273)>>(uint(v252)%32))-int32(base.Ui32(v273<<(uint(v252)%32)>>(uint(v256)%32)*v218)>>(uint(v259)%32)))&v262<<(uint(v264)%32)
							v289 = *(*int32)(unsafe.Add(mBase, uint32(v288)))
							*(*int32)(unsafe.Add(mBase, uint32(v288))) = v289 + v268
							v296 = v241 + v264
							if v200&int32(2147483646) != v296 {
								v241 = v296
								v244 = v244 + int32(8)
								continue
							} else {
								break
							}
							break
						}
						v299 = v296
					}
					if v200&int32(1) == int32(0) {
					} else {
						v311 = int32(2)
						v314 = *(*int32)(unsafe.Add(mBase, uint32(v225+v299<<(uint(v311)%32))))
						v315 = int32(16)
						v329 = l5 + (int32(base.Ui32(v314)>>(uint(v315)%32))-int32(base.Ui32(v314<<(uint(v315)%32)>>(uint(int32(24))%32)*v218)>>(uint(int32(5))%32)))&int32(255)<<(uint(v311)%32)
						v330 = *(*int32)(unsafe.Add(mBase, uint32(v329)))
						*(*int32)(unsafe.Add(mBase, uint32(v329))) = v330 + int32(1)
					}
					if int32(1) < v228 {
						v225 = v225 + l1<<(uint(int32(2))%32)
						v228 = v228 + int32(-1)
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

var F_CollectColorRedTransforms_SSE2__k0 = [2]uint64{0xff000000ff00, 0xff000000ff00}
var F_CollectColorRedTransforms_SSE2__k1 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_CollectColorRedTransforms_SSE2__k2 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_CollectColorRedTransforms_SSE2__k3 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_CollectColorRedTransforms_SSE2__k4 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}

func F_CollectColorRedTransforms_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v27 base.V128
	_ = v27
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v71 base.V128
	_ = v71
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v113 base.V128
	_ = v113
	var v116 base.V128
	_ = v116
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v128 int32
	_ = v128
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v160 int32
	_ = v160
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 base.V128
	_ = v180
	var v184 base.V128
	_ = v184
	var v185 int32
	_ = v185
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v191 int32
	_ = v191
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v223 int32
	_ = v223
	var v240 int32
	_ = v240
	var v243 int32
	_ = v243
	var v258 int32
	_ = v258
	var v265 int32
	_ = v265
	var v268 int32
	_ = v268
	var v281 int32
	_ = v281
	var v284 int32
	_ = v284
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v296 int32
	_ = v296
	var v299 int32
	_ = v299
	var v302 int32
	_ = v302
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v313 int32
	_ = v313
	var v328 int32
	_ = v328
	var v329 int32
	_ = v329
	var v336 int32
	_ = v336
	var v339 int32
	_ = v339
	var v351 int32
	_ = v351
	var v354 int32
	_ = v354
	var v355 int32
	_ = v355
	var v369 int32
	_ = v369
	var v370 int32
	_ = v370
	if l2 < int32(4) {
	} else {
		if l3 < int32(1) {
		} else {
			v20 = int32(8)
			v27 = base.Simd_g_i32x4_splat(int32(base.Ui32(base.I32_extend16_s(l4<<(uint(v20)%32)))>>(uint(int32(5))%32)) & int32(_a_F_CollectColorRedTransforms_SSE41_0))
			if base.Ui32(l2) < base.Ui32(v20) {
				v176 = l3
				v177 = l0
				for {
					v180 = base.Simd_g_v128_load(m, v177, int32(0))
					v184 = base.Simd_g_i16x8_sub(v180, base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_v128_and(v180, base.Simd_g_const(&F_CollectColorRedTransforms_SSE41__k0)), v27))
					v185 = int32(2)
					v189 = l5 + base.Simd_g_i8x16_extract_lane_u_l2(v184)<<(uint(v185)%32)
					v190 = *(*int32)(unsafe.Add(mBase, uint32(v189)))
					v191 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v189))) = v190 + v191
					v198 = l5 + base.Simd_g_i8x16_extract_lane_u_l6(v184)<<(uint(v185)%32)
					v199 = *(*int32)(unsafe.Add(mBase, uint32(v198)))
					*(*int32)(unsafe.Add(mBase, uint32(v198))) = v199 + v191
					v207 = l5 + base.Simd_g_i8x16_extract_lane_u_l10(v184)<<(uint(v185)%32)
					v208 = *(*int32)(unsafe.Add(mBase, uint32(v207)))
					*(*int32)(unsafe.Add(mBase, uint32(v207))) = v208 + v191
					v216 = l5 + base.Simd_g_i8x16_extract_lane_u_l14(v184)<<(uint(v185)%32)
					v217 = *(*int32)(unsafe.Add(mBase, uint32(v216)))
					*(*int32)(unsafe.Add(mBase, uint32(v216))) = v217 + v191
					v223 = v176 + int32(-1)
					if v223 != 0 {
						v176 = v223
						v177 = v177 + l1<<(uint(int32(2))%32)
						continue
					} else {
						break
					}
					break
				}
			} else {
				v43 = l0 + int32(16)
				v44 = int32(0)
				for {
					v55 = base.Simd_g_v128_load(m, l0+v44*l1<<(uint(int32(2))%32), int32(0))
					v56 = base.Simd_g_const(&F_CollectColorRedTransforms_SSE41__k0)
					v71 = base.Simd_g_i16x8_sub(v55, base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_v128_and(v55, v56), v27))
					v73 = int32(8)
					v74 = v43
					for {
						v76 = int32(2)
						v80 = l5 + base.Simd_g_i8x16_extract_lane_u_l2(v71)<<(uint(v76)%32)
						v81 = *(*int32)(unsafe.Add(mBase, uint32(v80)))
						v82 = int32(1)
						*(*int32)(unsafe.Add(mBase, uint32(v80))) = v81 + v82
						v89 = l5 + base.Simd_g_i8x16_extract_lane_u_l6(v71)<<(uint(v76)%32)
						v90 = *(*int32)(unsafe.Add(mBase, uint32(v89)))
						*(*int32)(unsafe.Add(mBase, uint32(v89))) = v90 + v82
						v98 = l5 + base.Simd_g_i8x16_extract_lane_u_l10(v71)<<(uint(v76)%32)
						v99 = *(*int32)(unsafe.Add(mBase, uint32(v98)))
						*(*int32)(unsafe.Add(mBase, uint32(v98))) = v99 + v82
						v107 = l5 + base.Simd_g_i8x16_extract_lane_u_l14(v71)<<(uint(v76)%32)
						v108 = *(*int32)(unsafe.Add(mBase, uint32(v107)))
						*(*int32)(unsafe.Add(mBase, uint32(v107))) = v108 + v82
						v113 = base.Simd_g_v128_load(m, v74, int32(0))
						v116 = base.Simd_g_i16x8_sub(v113, base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_v128_and(v113, v56), v27))
						v120 = v73 + int32(4)
						if v120 <= l2 {
							v71 = v116
							v73 = v120
							v74 = v74 + int32(16)
							continue
						} else {
							break
						}
						break
					}
					v122 = int32(2)
					v126 = l5 + base.Simd_g_i8x16_extract_lane_u_l2(v116)<<(uint(v122)%32)
					v127 = *(*int32)(unsafe.Add(mBase, uint32(v126)))
					v128 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v126))) = v127 + v128
					v135 = l5 + base.Simd_g_i8x16_extract_lane_u_l6(v116)<<(uint(v122)%32)
					v136 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
					*(*int32)(unsafe.Add(mBase, uint32(v135))) = v136 + v128
					v144 = l5 + base.Simd_g_i8x16_extract_lane_u_l10(v116)<<(uint(v122)%32)
					v145 = *(*int32)(unsafe.Add(mBase, uint32(v144)))
					*(*int32)(unsafe.Add(mBase, uint32(v144))) = v145 + v128
					v153 = l5 + base.Simd_g_i8x16_extract_lane_u_l14(v116)<<(uint(v122)%32)
					v154 = *(*int32)(unsafe.Add(mBase, uint32(v153)))
					*(*int32)(unsafe.Add(mBase, uint32(v153))) = v154 + v128
					v160 = v44 + v128
					if v160 != l3 {
						v43 = v43 + l1<<(uint(int32(2))%32)
						v44 = v160
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	v240 = l2 & int32(3)
	if v240 == int32(0) {
	} else {
		v243 = int32(2)
		if l3 < int32(1) {
		} else {
			if v240 < int32(1) {
			} else {
				v258 = base.I32_extend8_s(l4)
				v265 = l0 + l2<<(uint(v243)%32) - v240<<(uint(v243)%32)
				v268 = l3
				for {
					if v240 == int32(1) {
						v339 = int32(0)
					} else {
						v281 = int32(0)
						v284 = v265
						for {
							v291 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
							v292 = int32(16)
							v296 = int32(24)
							v299 = int32(5)
							v302 = int32(255)
							v304 = int32(2)
							v306 = l5 + (int32(base.Ui32(v291)>>(uint(v292)%32))-int32(base.Ui32(v291<<(uint(v292)%32)>>(uint(v296)%32)*v258)>>(uint(v299)%32)))&v302<<(uint(v304)%32)
							v307 = *(*int32)(unsafe.Add(mBase, uint32(v306)))
							v308 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(v306))) = v307 + v308
							v313 = *(*int32)(unsafe.Add(mBase, uint32(v284+int32(4))))
							v328 = l5 + (int32(base.Ui32(v313)>>(uint(v292)%32))-int32(base.Ui32(v313<<(uint(v292)%32)>>(uint(v296)%32)*v258)>>(uint(v299)%32)))&v302<<(uint(v304)%32)
							v329 = *(*int32)(unsafe.Add(mBase, uint32(v328)))
							*(*int32)(unsafe.Add(mBase, uint32(v328))) = v329 + v308
							v336 = v281 + v304
							if v240&int32(2147483646) != v336 {
								v281 = v336
								v284 = v284 + int32(8)
								continue
							} else {
								break
							}
							break
						}
						v339 = v336
					}
					if v240&int32(1) == int32(0) {
					} else {
						v351 = int32(2)
						v354 = *(*int32)(unsafe.Add(mBase, uint32(v265+v339<<(uint(v351)%32))))
						v355 = int32(16)
						v369 = l5 + (int32(base.Ui32(v354)>>(uint(v355)%32))-int32(base.Ui32(v354<<(uint(v355)%32)>>(uint(int32(24))%32)*v258)>>(uint(int32(5))%32)))&int32(255)<<(uint(v351)%32)
						v370 = *(*int32)(unsafe.Add(mBase, uint32(v369)))
						*(*int32)(unsafe.Add(mBase, uint32(v369))) = v370 + int32(1)
					}
					if int32(1) < v268 {
						v265 = v265 + l1<<(uint(int32(2))%32)
						v268 = v268 + int32(-1)
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

var F_CollectColorRedTransforms_SSE41__k0 = [2]uint64{0xff000000ff00, 0xff000000ff00}

func F_CollectHistogram_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v143 int32
	_ = v143
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 base.V128
	_ = v164
	var v165 int32
	_ = v165
	var v166 base.V128
	_ = v166
	var v167 base.V128
	_ = v167
	var v169 base.V128
	_ = v169
	var v170 base.V128
	_ = v170
	var v172 int32
	_ = v172
	var v174 base.V128
	_ = v174
	var v176 base.V128
	_ = v176
	var v182 base.V128
	_ = v182
	var v183 base.V128
	_ = v183
	var v184 int32
	_ = v184
	var v185 base.V128
	_ = v185
	var v186 int32
	_ = v186
	var v187 base.V128
	_ = v187
	var v193 base.V128
	_ = v193
	var v195 base.V128
	_ = v195
	var v202 base.V128
	_ = v202
	var v203 base.V128
	_ = v203
	var v204 base.V128
	_ = v204
	var v205 base.V128
	_ = v205
	var v206 base.V128
	_ = v206
	var v207 base.V128
	_ = v207
	var v212 base.V128
	_ = v212
	var v213 base.V128
	_ = v213
	var v218 int32
	_ = v218
	var v226 base.V128
	_ = v226
	var v228 base.V128
	_ = v228
	var v230 base.V128
	_ = v230
	var v234 base.V128
	_ = v234
	var v236 base.V128
	_ = v236
	var v237 base.V128
	_ = v237
	var v239 base.V128
	_ = v239
	var v241 base.V128
	_ = v241
	var v243 int32
	_ = v243
	var v245 base.V128
	_ = v245
	var v249 base.V128
	_ = v249
	var v254 int32
	_ = v254
	var v255 base.V128
	_ = v255
	var v262 int32
	_ = v262
	var v264 base.V128
	_ = v264
	var v265 base.V128
	_ = v265
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v288 int32
	_ = v288
	var v289 int32
	_ = v289
	var v297 int32
	_ = v297
	var v298 int32
	_ = v298
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v311 int32
	_ = v311
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v329 int32
	_ = v329
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v346 base.V128
	_ = v346
	var v353 base.V128
	_ = v353
	var v358 int32
	_ = v358
	var v359 int32
	_ = v359
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v385 int32
	_ = v385
	var v386 int32
	_ = v386
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v403 int32
	_ = v403
	var v404 int32
	_ = v404
	var v412 int32
	_ = v412
	var v413 int32
	_ = v413
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v429 int32
	_ = v429
	var v444 int32
	_ = v444
	var v536 int32
	_ = v536
	var v538 int32
	_ = v538
	var v539 int32
	_ = v539
	var v543 int32
	_ = v543
	var v545 int32
	_ = v545
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v549 int32
	_ = v549
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v555 int32
	_ = v555
	var v557 int32
	_ = v557
	var v558 int32
	_ = v558
	var v559 int32
	_ = v559
	var v561 int32
	_ = v561
	var v562 int32
	_ = v562
	var v563 int32
	_ = v563
	var v565 int32
	_ = v565
	var v566 int32
	_ = v566
	var v567 int32
	_ = v567
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v573 int32
	_ = v573
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v577 int32
	_ = v577
	var v578 int32
	_ = v578
	var v579 int32
	_ = v579
	var v581 int32
	_ = v581
	var v582 int32
	_ = v582
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v593 int32
	_ = v593
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v597 int32
	_ = v597
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v602 int32
	_ = v602
	var v603 int32
	_ = v603
	var v605 int32
	_ = v605
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v609 int32
	_ = v609
	var v610 int32
	_ = v610
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v617 int32
	_ = v617
	var v618 int32
	_ = v618
	var v619 int32
	_ = v619
	var v621 int32
	_ = v621
	var v622 int32
	_ = v622
	var v623 int32
	_ = v623
	var v625 int32
	_ = v625
	var v626 int32
	_ = v626
	var v627 int32
	_ = v627
	var v629 int32
	_ = v629
	var v630 int32
	_ = v630
	var v631 int32
	_ = v631
	var v633 int32
	_ = v633
	var v634 int32
	_ = v634
	var v635 int32
	_ = v635
	var v637 int32
	_ = v637
	var v638 int32
	_ = v638
	var v639 int32
	_ = v639
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v643 int32
	_ = v643
	var v645 int32
	_ = v645
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v655 int32
	_ = v655
	var v657 int32
	_ = v657
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v662 int32
	_ = v662
	var v664 int32
	_ = v664
	var v667 int32
	_ = v667
	var v669 int32
	_ = v669
	var v670 int32
	_ = v670
	var v672 int32
	_ = v672
	var v673 int32
	_ = v673
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v681 int32
	_ = v681
	var v682 int32
	_ = v682
	var v684 int32
	_ = v684
	var v685 int32
	_ = v685
	var v687 int32
	_ = v687
	var v688 int32
	_ = v688
	var v690 int32
	_ = v690
	var v691 int32
	_ = v691
	var v693 int32
	_ = v693
	var v694 int32
	_ = v694
	var v696 int32
	_ = v696
	var v697 int32
	_ = v697
	var v699 int32
	_ = v699
	var v700 int32
	_ = v700
	var v702 int32
	_ = v702
	var v703 int32
	_ = v703
	var v705 int32
	_ = v705
	var v706 int32
	_ = v706
	var v708 int32
	_ = v708
	var v709 int32
	_ = v709
	var v711 int32
	_ = v711
	var v712 int32
	_ = v712
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v717 int32
	_ = v717
	var v718 int32
	_ = v718
	var v720 int32
	_ = v720
	var v721 int32
	_ = v721
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
	var v730 int32
	_ = v730
	var v732 int32
	_ = v732
	var v733 int32
	_ = v733
	var v735 int32
	_ = v735
	var v736 int32
	_ = v736
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v741 int32
	_ = v741
	var v742 int32
	_ = v742
	var v744 int32
	_ = v744
	var v745 int32
	_ = v745
	var v747 int32
	_ = v747
	var v748 int32
	_ = v748
	var v750 int32
	_ = v750
	var v751 int32
	_ = v751
	var v753 int32
	_ = v753
	var v754 int32
	_ = v754
	var v756 int32
	_ = v756
	var v757 int32
	_ = v757
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	v15 = m.G0
	v16 = int32(128)
	v17 = v15 - v16
	m.G0 = v17
	base.MemoryFill(m, v17, int32(0), v16)
	if l3 <= l2 {
	} else {
		v143 = m.G76
		v153 = l3 - l2
		v154 = v143 + l2<<(uint(int32(2))%32)
		for {
			v161 = *(*int32)(unsafe.Add(mBase, uint32(v154)))
			v162 = l0 + v161
			v163 = int32(0)
			v164 = base.Simd_g_v128_load64_zero(m, v162, v163)
			v165 = int32(32)
			v166 = base.Simd_g_v128_load64_zero(m, v162, v165)
			v167 = base.Simd_g_const(&F_CollectHistogram_SSE2__k0)
			v169 = base.Simd_g_const(&F_CollectHistogram_SSE2__k1)
			v170 = base.Simd_g_const(&F_CollectHistogram_SSE2__k2)
			v172 = l1 + v161
			v174 = base.Simd_g_v128_load64_zero(m, v172, v163)
			v176 = base.Simd_g_v128_load64_zero(m, v172, v165)
			v182 = base.Simd_g_const(&F_CollectHistogram_SSE2__k3)
			v183 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v164, v166, base.Simd_g_const(&F_CollectHistogram_SSE2__k4), base.Simd_g_const(&F_CollectHistogram_SSE2__k5)), v169, base.Simd_g_const(&F_CollectHistogram_SSE2__k6), base.Simd_g_const(&F_CollectHistogram_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v174, v176, base.Simd_g_const(&F_CollectHistogram_SSE2__k4), base.Simd_g_const(&F_CollectHistogram_SSE2__k5)), v169, base.Simd_g_const(&F_CollectHistogram_SSE2__k6), base.Simd_g_const(&F_CollectHistogram_SSE2__k7))), base.Simd_g_const(&F_CollectHistogram_SSE2__k3))
			v184 = int32(64)
			v185 = base.Simd_g_v128_load64_zero(m, v162, v184)
			v186 = int32(96)
			v187 = base.Simd_g_v128_load64_zero(m, v162, v186)
			v193 = base.Simd_g_v128_load64_zero(m, v172, v184)
			v195 = base.Simd_g_v128_load64_zero(m, v172, v186)
			v202 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v185, v187, base.Simd_g_const(&F_CollectHistogram_SSE2__k4), base.Simd_g_const(&F_CollectHistogram_SSE2__k5)), v169, base.Simd_g_const(&F_CollectHistogram_SSE2__k6), base.Simd_g_const(&F_CollectHistogram_SSE2__k7)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v193, v195, base.Simd_g_const(&F_CollectHistogram_SSE2__k4), base.Simd_g_const(&F_CollectHistogram_SSE2__k5)), v169, base.Simd_g_const(&F_CollectHistogram_SSE2__k6), base.Simd_g_const(&F_CollectHistogram_SSE2__k7))), base.Simd_g_const(&F_CollectHistogram_SSE2__k3))
			v203 = base.Simd_g_const(&F_CollectHistogram_SSE2__k8)
			v204 = base.Simd_g_i8x16_shuffle2(v183, v202, base.Simd_g_const(&F_CollectHistogram_SSE2__k9), base.Simd_g_const(&F_CollectHistogram_SSE2__k10))
			v205 = base.Simd_g_const(&F_CollectHistogram_SSE2__k11)
			v206 = base.Simd_g_i8x16_shuffle2(v183, v202, base.Simd_g_const(&F_CollectHistogram_SSE2__k12), base.Simd_g_const(&F_CollectHistogram_SSE2__k13))
			v207 = base.Simd_g_i16x8_add(v204, v206)
			v212 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(v207, base.Simd_g_const(&F_CollectHistogram_SSE2__k14)), base.Simd_g_i32x4_dot_i16x8_s(v207, base.Simd_g_const(&F_CollectHistogram_SSE2__k15)))
			v213 = base.Simd_g_i16x8_sub(v204, v206)
			v218 = int32(9)
			v226 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v213, base.Simd_g_const(&F_CollectHistogram_SSE2__k16)), base.Simd_g_const(&F_CollectHistogram_SSE2__k17)), v218), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v213, base.Simd_g_const(&F_CollectHistogram_SSE2__k18)), base.Simd_g_const(&F_CollectHistogram_SSE2__k19)), v218))
			v228 = base.Simd_g_i8x16_shuffle2(v212, v226, base.Simd_g_const(&F_CollectHistogram_SSE2__k4), base.Simd_g_const(&F_CollectHistogram_SSE2__k5))
			v230 = base.Simd_g_i8x16_shuffle2(v212, v226, base.Simd_g_const(&F_CollectHistogram_SSE2__k20), base.Simd_g_const(&F_CollectHistogram_SSE2__k21))
			v234 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_shuffle2(v228, v230, base.Simd_g_const(&F_CollectHistogram_SSE2__k22), base.Simd_g_const(&F_CollectHistogram_SSE2__k23)), base.Simd_g_const(&F_CollectHistogram_SSE2__k24))
			v236 = base.Simd_g_i8x16_shuffle2(v228, v230, base.Simd_g_const(&F_CollectHistogram_SSE2__k25), base.Simd_g_const(&F_CollectHistogram_SSE2__k26))
			v237 = base.Simd_g_i16x8_add(v234, v236)
			v239 = base.Simd_g_i16x8_add(v237, base.Simd_g_const(&F_CollectHistogram_SSE2__k27))
			v241 = base.Simd_g_i8x16_shuffle2(v237, v237, base.Simd_g_const(&F_CollectHistogram_SSE2__k12), base.Simd_g_const(&F_CollectHistogram_SSE2__k13))
			v243 = int32(4)
			v245 = base.Simd_g_i16x8_sub(v236, v234)
			v249 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v245, v245, base.Simd_g_const(&F_CollectHistogram_SSE2__k12), base.Simd_g_const(&F_CollectHistogram_SSE2__k13)), v245, base.Simd_g_const(&F_CollectHistogram_SSE2__k4), base.Simd_g_const(&F_CollectHistogram_SSE2__k5))
			v254 = int32(16)
			v255 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v249, base.Simd_g_const(&F_CollectHistogram_SSE2__k28)), base.Simd_g_const(&F_CollectHistogram_SSE2__k29)), v254)
			v262 = int32(3)
			v264 = base.Simd_g_const(&F_CollectHistogram_SSE2__k30)
			v265 = base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v239, v241), v243), base.Simd_g_i16x8_add(base.Simd_g_i16x8_narrow_i32x4_s(v255, v255), base.Simd_g_i16x8_eq(v236, v234)), base.Simd_g_const(&F_CollectHistogram_SSE2__k9), base.Simd_g_const(&F_CollectHistogram_SSE2__k10))), v262), v264)
			v268 = int32(2)
			v270 = v17 + base.Simd_g_i16x8_extract_lane_s_l0(v265)<<(uint(v268)%32)
			v271 = *(*int32)(unsafe.Add(mBase, uint32(v270)))
			v272 = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v270))) = v271 + v272
			v279 = v17 + base.Simd_g_i16x8_extract_lane_s_l1(v265)<<(uint(v268)%32)
			v280 = *(*int32)(unsafe.Add(mBase, uint32(v279)))
			*(*int32)(unsafe.Add(mBase, uint32(v279))) = v280 + v272
			v288 = v17 + base.Simd_g_i16x8_extract_lane_s_l2(v265)<<(uint(v268)%32)
			v289 = *(*int32)(unsafe.Add(mBase, uint32(v288)))
			*(*int32)(unsafe.Add(mBase, uint32(v288))) = v289 + v272
			v297 = v17 + base.Simd_g_i16x8_extract_lane_s_l3(v265)<<(uint(v268)%32)
			v298 = *(*int32)(unsafe.Add(mBase, uint32(v297)))
			*(*int32)(unsafe.Add(mBase, uint32(v297))) = v298 + v272
			v306 = v17 + base.Simd_g_i16x8_extract_lane_s_l4(v265)<<(uint(v268)%32)
			v307 = *(*int32)(unsafe.Add(mBase, uint32(v306)))
			*(*int32)(unsafe.Add(mBase, uint32(v306))) = v307 + v272
			v311 = int32(5)
			v315 = v17 + base.Simd_g_i16x8_extract_lane_s_l5(v265)<<(uint(v268)%32)
			v316 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
			*(*int32)(unsafe.Add(mBase, uint32(v315))) = v316 + v272
			v320 = int32(6)
			v324 = v17 + base.Simd_g_i16x8_extract_lane_s_l6(v265)<<(uint(v268)%32)
			v325 = *(*int32)(unsafe.Add(mBase, uint32(v324)))
			*(*int32)(unsafe.Add(mBase, uint32(v324))) = v325 + v272
			v329 = int32(7)
			v333 = v17 + base.Simd_g_i16x8_extract_lane_s_l7(v265)<<(uint(v268)%32)
			v334 = *(*int32)(unsafe.Add(mBase, uint32(v333)))
			*(*int32)(unsafe.Add(mBase, uint32(v333))) = v334 + v272
			v346 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v249, base.Simd_g_const(&F_CollectHistogram_SSE2__k31)), base.Simd_g_const(&F_CollectHistogram_SSE2__k32)), v254)
			v353 = base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_abs(base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(v239, v241), v243), base.Simd_g_i16x8_narrow_i32x4_s(v346, v346), base.Simd_g_const(&F_CollectHistogram_SSE2__k9), base.Simd_g_const(&F_CollectHistogram_SSE2__k10))), v262), v264)
			v358 = v17 + base.Simd_g_i16x8_extract_lane_s_l0(v353)<<(uint(v268)%32)
			v359 = *(*int32)(unsafe.Add(mBase, uint32(v358)))
			*(*int32)(unsafe.Add(mBase, uint32(v358))) = v359 + v272
			v367 = v17 + base.Simd_g_i16x8_extract_lane_s_l1(v353)<<(uint(v268)%32)
			v368 = *(*int32)(unsafe.Add(mBase, uint32(v367)))
			*(*int32)(unsafe.Add(mBase, uint32(v367))) = v368 + v272
			v376 = v17 + base.Simd_g_i16x8_extract_lane_s_l2(v353)<<(uint(v268)%32)
			v377 = *(*int32)(unsafe.Add(mBase, uint32(v376)))
			*(*int32)(unsafe.Add(mBase, uint32(v376))) = v377 + v272
			v385 = v17 + base.Simd_g_i16x8_extract_lane_s_l3(v353)<<(uint(v268)%32)
			v386 = *(*int32)(unsafe.Add(mBase, uint32(v385)))
			*(*int32)(unsafe.Add(mBase, uint32(v385))) = v386 + v272
			v394 = v17 + base.Simd_g_i16x8_extract_lane_s_l4(v353)<<(uint(v268)%32)
			v395 = *(*int32)(unsafe.Add(mBase, uint32(v394)))
			*(*int32)(unsafe.Add(mBase, uint32(v394))) = v395 + v272
			v403 = v17 + base.Simd_g_i16x8_extract_lane_s_l5(v353)<<(uint(v268)%32)
			v404 = *(*int32)(unsafe.Add(mBase, uint32(v403)))
			*(*int32)(unsafe.Add(mBase, uint32(v403))) = v404 + v272
			v412 = v17 + base.Simd_g_i16x8_extract_lane_s_l6(v353)<<(uint(v268)%32)
			v413 = *(*int32)(unsafe.Add(mBase, uint32(v412)))
			*(*int32)(unsafe.Add(mBase, uint32(v412))) = v413 + v272
			v421 = v17 + base.Simd_g_i16x8_extract_lane_s_l7(v353)<<(uint(v268)%32)
			v422 = *(*int32)(unsafe.Add(mBase, uint32(v421)))
			*(*int32)(unsafe.Add(mBase, uint32(v421))) = v422 + v272
			v429 = v153 + int32(-1)
			if v429 != 0 {
				v153 = v429
				v154 = v154 + v243
				continue
			} else {
				break
			}
			break
		}
	}
	v444 = int32(0)
	v536 = *(*int32)(unsafe.Add(mBase, uint32(v17)+4))
	v538 = base.B2i32(v444 < v536)
	v539 = *(*int32)(unsafe.Add(mBase, uint32(v17)))
	v543 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v545 = base.B2i32(v444 < v543)
	if v444 < v543 {
		v546 = int32(2)
	} else {
		v546 = v538 | base.B2i32(v539 < int32(1))
	}
	v547 = *(*int32)(unsafe.Add(mBase, uint32(v17)+12))
	v549 = base.B2i32(int32(0) < v547)
	if int32(0) < v547 {
		v550 = int32(3)
	} else {
		v550 = v546
	}
	v551 = *(*int32)(unsafe.Add(mBase, uint32(v17)+16))
	v553 = base.B2i32(int32(0) < v551)
	if int32(0) < v551 {
		v554 = int32(4)
	} else {
		v554 = v550
	}
	v555 = *(*int32)(unsafe.Add(mBase, uint32(v17)+20))
	v557 = base.B2i32(int32(0) < v555)
	if int32(0) < v555 {
		v558 = int32(5)
	} else {
		v558 = v554
	}
	v559 = *(*int32)(unsafe.Add(mBase, uint32(v17)+24))
	v561 = base.B2i32(int32(0) < v559)
	if int32(0) < v559 {
		v562 = int32(6)
	} else {
		v562 = v558
	}
	v563 = *(*int32)(unsafe.Add(mBase, uint32(v17)+28))
	v565 = base.B2i32(int32(0) < v563)
	if int32(0) < v563 {
		v566 = int32(7)
	} else {
		v566 = v562
	}
	v567 = *(*int32)(unsafe.Add(mBase, uint32(v17)+32))
	v569 = base.B2i32(int32(0) < v567)
	if int32(0) < v567 {
		v570 = int32(8)
	} else {
		v570 = v566
	}
	v571 = *(*int32)(unsafe.Add(mBase, uint32(v17)+36))
	v573 = base.B2i32(int32(0) < v571)
	if int32(0) < v571 {
		v574 = int32(9)
	} else {
		v574 = v570
	}
	v575 = *(*int32)(unsafe.Add(mBase, uint32(v17)+40))
	v577 = base.B2i32(int32(0) < v575)
	if int32(0) < v575 {
		v578 = int32(10)
	} else {
		v578 = v574
	}
	v579 = *(*int32)(unsafe.Add(mBase, uint32(v17)+44))
	v581 = base.B2i32(int32(0) < v579)
	if int32(0) < v579 {
		v582 = int32(11)
	} else {
		v582 = v578
	}
	v583 = *(*int32)(unsafe.Add(mBase, uint32(v17)+48))
	v585 = base.B2i32(int32(0) < v583)
	if int32(0) < v583 {
		v586 = int32(12)
	} else {
		v586 = v582
	}
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v17)+52))
	v589 = base.B2i32(int32(0) < v587)
	if int32(0) < v587 {
		v590 = int32(13)
	} else {
		v590 = v586
	}
	v591 = *(*int32)(unsafe.Add(mBase, uint32(v17)+56))
	v593 = base.B2i32(int32(0) < v591)
	if int32(0) < v591 {
		v594 = int32(14)
	} else {
		v594 = v590
	}
	v595 = *(*int32)(unsafe.Add(mBase, uint32(v17)+60))
	v597 = base.B2i32(int32(0) < v595)
	if int32(0) < v595 {
		v598 = int32(15)
	} else {
		v598 = v594
	}
	v599 = *(*int32)(unsafe.Add(mBase, uint32(v17)+64))
	v601 = base.B2i32(int32(0) < v599)
	if int32(0) < v599 {
		v602 = int32(16)
	} else {
		v602 = v598
	}
	v603 = *(*int32)(unsafe.Add(mBase, uint32(v17)+68))
	v605 = base.B2i32(int32(0) < v603)
	if int32(0) < v603 {
		v606 = int32(17)
	} else {
		v606 = v602
	}
	v607 = *(*int32)(unsafe.Add(mBase, uint32(v17)+72))
	v609 = base.B2i32(int32(0) < v607)
	if int32(0) < v607 {
		v610 = int32(18)
	} else {
		v610 = v606
	}
	v611 = *(*int32)(unsafe.Add(mBase, uint32(v17)+76))
	v613 = base.B2i32(int32(0) < v611)
	if int32(0) < v611 {
		v614 = int32(19)
	} else {
		v614 = v610
	}
	v615 = *(*int32)(unsafe.Add(mBase, uint32(v17)+80))
	v617 = base.B2i32(int32(0) < v615)
	if int32(0) < v615 {
		v618 = int32(20)
	} else {
		v618 = v614
	}
	v619 = *(*int32)(unsafe.Add(mBase, uint32(v17)+84))
	v621 = base.B2i32(int32(0) < v619)
	if int32(0) < v619 {
		v622 = int32(21)
	} else {
		v622 = v618
	}
	v623 = *(*int32)(unsafe.Add(mBase, uint32(v17)+88))
	v625 = base.B2i32(int32(0) < v623)
	if int32(0) < v623 {
		v626 = int32(22)
	} else {
		v626 = v622
	}
	v627 = *(*int32)(unsafe.Add(mBase, uint32(v17)+92))
	v629 = base.B2i32(int32(0) < v627)
	if int32(0) < v627 {
		v630 = int32(23)
	} else {
		v630 = v626
	}
	v631 = *(*int32)(unsafe.Add(mBase, uint32(v17)+96))
	v633 = base.B2i32(int32(0) < v631)
	if int32(0) < v631 {
		v634 = int32(24)
	} else {
		v634 = v630
	}
	v635 = *(*int32)(unsafe.Add(mBase, uint32(v17)+100))
	v637 = base.B2i32(int32(0) < v635)
	if int32(0) < v635 {
		v638 = int32(25)
	} else {
		v638 = v634
	}
	v639 = *(*int32)(unsafe.Add(mBase, uint32(v17)+104))
	v641 = base.B2i32(int32(0) < v639)
	if int32(0) < v639 {
		v642 = int32(26)
	} else {
		v642 = v638
	}
	v643 = *(*int32)(unsafe.Add(mBase, uint32(v17)+108))
	v645 = base.B2i32(int32(0) < v643)
	if int32(0) < v643 {
		v646 = int32(27)
	} else {
		v646 = v642
	}
	v647 = *(*int32)(unsafe.Add(mBase, uint32(v17)+112))
	v649 = base.B2i32(int32(0) < v647)
	if int32(0) < v647 {
		v650 = int32(28)
	} else {
		v650 = v646
	}
	v651 = *(*int32)(unsafe.Add(mBase, uint32(v17)+116))
	v653 = base.B2i32(int32(0) < v651)
	if int32(0) < v651 {
		v654 = int32(29)
	} else {
		v654 = v650
	}
	v655 = *(*int32)(unsafe.Add(mBase, uint32(v17)+120))
	v657 = base.B2i32(int32(0) < v655)
	if int32(0) < v655 {
		v658 = int32(30)
	} else {
		v658 = v654
	}
	v659 = *(*int32)(unsafe.Add(mBase, uint32(v17)+124))
	v661 = base.B2i32(int32(0) < v659)
	if int32(0) < v659 {
		v662 = int32(31)
	} else {
		v662 = v658
	}
	*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v662
	v664 = int32(0)
	if v664 < v539 {
		v667 = v539
	} else {
		v667 = v664
	}
	if v667 < v536 {
		v669 = v536
	} else {
		v669 = v667
	}
	if v444 < v536 {
		v670 = v669
	} else {
		v670 = v667
	}
	if v670 < v543 {
		v672 = v543
	} else {
		v672 = v670
	}
	if v444 < v543 {
		v673 = v672
	} else {
		v673 = v670
	}
	if v673 < v547 {
		v675 = v547
	} else {
		v675 = v673
	}
	if int32(0) < v547 {
		v676 = v675
	} else {
		v676 = v673
	}
	if v676 < v551 {
		v678 = v551
	} else {
		v678 = v676
	}
	if int32(0) < v551 {
		v679 = v678
	} else {
		v679 = v676
	}
	if v679 < v555 {
		v681 = v555
	} else {
		v681 = v679
	}
	if int32(0) < v555 {
		v682 = v681
	} else {
		v682 = v679
	}
	if v682 < v559 {
		v684 = v559
	} else {
		v684 = v682
	}
	if int32(0) < v559 {
		v685 = v684
	} else {
		v685 = v682
	}
	if v685 < v563 {
		v687 = v563
	} else {
		v687 = v685
	}
	if int32(0) < v563 {
		v688 = v687
	} else {
		v688 = v685
	}
	if v688 < v567 {
		v690 = v567
	} else {
		v690 = v688
	}
	if int32(0) < v567 {
		v691 = v690
	} else {
		v691 = v688
	}
	if v691 < v571 {
		v693 = v571
	} else {
		v693 = v691
	}
	if int32(0) < v571 {
		v694 = v693
	} else {
		v694 = v691
	}
	if v694 < v575 {
		v696 = v575
	} else {
		v696 = v694
	}
	if int32(0) < v575 {
		v697 = v696
	} else {
		v697 = v694
	}
	if v697 < v579 {
		v699 = v579
	} else {
		v699 = v697
	}
	if int32(0) < v579 {
		v700 = v699
	} else {
		v700 = v697
	}
	if v700 < v583 {
		v702 = v583
	} else {
		v702 = v700
	}
	if int32(0) < v583 {
		v703 = v702
	} else {
		v703 = v700
	}
	if v703 < v587 {
		v705 = v587
	} else {
		v705 = v703
	}
	if int32(0) < v587 {
		v706 = v705
	} else {
		v706 = v703
	}
	if v706 < v591 {
		v708 = v591
	} else {
		v708 = v706
	}
	if int32(0) < v591 {
		v709 = v708
	} else {
		v709 = v706
	}
	if v709 < v595 {
		v711 = v595
	} else {
		v711 = v709
	}
	if int32(0) < v595 {
		v712 = v711
	} else {
		v712 = v709
	}
	if v712 < v599 {
		v714 = v599
	} else {
		v714 = v712
	}
	if int32(0) < v599 {
		v715 = v714
	} else {
		v715 = v712
	}
	if v715 < v603 {
		v717 = v603
	} else {
		v717 = v715
	}
	if int32(0) < v603 {
		v718 = v717
	} else {
		v718 = v715
	}
	if v718 < v607 {
		v720 = v607
	} else {
		v720 = v718
	}
	if int32(0) < v607 {
		v721 = v720
	} else {
		v721 = v718
	}
	if v721 < v611 {
		v723 = v611
	} else {
		v723 = v721
	}
	if int32(0) < v611 {
		v724 = v723
	} else {
		v724 = v721
	}
	if v724 < v615 {
		v726 = v615
	} else {
		v726 = v724
	}
	if int32(0) < v615 {
		v727 = v726
	} else {
		v727 = v724
	}
	if v727 < v619 {
		v729 = v619
	} else {
		v729 = v727
	}
	if int32(0) < v619 {
		v730 = v729
	} else {
		v730 = v727
	}
	if v730 < v623 {
		v732 = v623
	} else {
		v732 = v730
	}
	if int32(0) < v623 {
		v733 = v732
	} else {
		v733 = v730
	}
	if v733 < v627 {
		v735 = v627
	} else {
		v735 = v733
	}
	if int32(0) < v627 {
		v736 = v735
	} else {
		v736 = v733
	}
	if v736 < v631 {
		v738 = v631
	} else {
		v738 = v736
	}
	if int32(0) < v631 {
		v739 = v738
	} else {
		v739 = v736
	}
	if v739 < v635 {
		v741 = v635
	} else {
		v741 = v739
	}
	if int32(0) < v635 {
		v742 = v741
	} else {
		v742 = v739
	}
	if v742 < v639 {
		v744 = v639
	} else {
		v744 = v742
	}
	if int32(0) < v639 {
		v745 = v744
	} else {
		v745 = v742
	}
	if v745 < v643 {
		v747 = v643
	} else {
		v747 = v745
	}
	if int32(0) < v643 {
		v748 = v747
	} else {
		v748 = v745
	}
	if v748 < v647 {
		v750 = v647
	} else {
		v750 = v748
	}
	if int32(0) < v647 {
		v751 = v750
	} else {
		v751 = v748
	}
	if v751 < v651 {
		v753 = v651
	} else {
		v753 = v751
	}
	if int32(0) < v651 {
		v754 = v753
	} else {
		v754 = v751
	}
	if v754 < v655 {
		v756 = v655
	} else {
		v756 = v754
	}
	if int32(0) < v655 {
		v757 = v756
	} else {
		v757 = v754
	}
	if v757 < v659 {
		v759 = v659
	} else {
		v759 = v757
	}
	if int32(0) < v659 {
		v760 = v759
	} else {
		v760 = v757
	}
	*(*int32)(unsafe.Add(mBase, uint32(l4))) = v760
	m.G0 = v17 + int32(128)
	return
}

var F_CollectHistogram_SSE2__k0 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_CollectHistogram_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_CollectHistogram_SSE2__k2 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_CollectHistogram_SSE2__k3 = [2]uint64{0x706050403020100, 0xd0c0f0e09080b0a}
var F_CollectHistogram_SSE2__k4 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_CollectHistogram_SSE2__k5 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_CollectHistogram_SSE2__k6 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_CollectHistogram_SSE2__k7 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_CollectHistogram_SSE2__k8 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_CollectHistogram_SSE2__k9 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_CollectHistogram_SSE2__k10 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_CollectHistogram_SSE2__k11 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_CollectHistogram_SSE2__k12 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_CollectHistogram_SSE2__k13 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_CollectHistogram_SSE2__k14 = [2]uint64{0x8000800080008, 0x8000800080008}
var F_CollectHistogram_SSE2__k15 = [2]uint64{0xfff80008fff80008, 0xfff80008fff80008}
var F_CollectHistogram_SSE2__k16 = [2]uint64{0x8a914e808a914e8, 0x8a914e808a914e8}
var F_CollectHistogram_SSE2__k17 = [2]uint64{0x71400000714, 0x71400000714}
var F_CollectHistogram_SSE2__k18 = [2]uint64{0xeb1808a9eb1808a9, 0xeb1808a9eb1808a9}
var F_CollectHistogram_SSE2__k19 = [2]uint64{0x3a9000003a9, 0x3a9000003a9}
var F_CollectHistogram_SSE2__k20 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_CollectHistogram_SSE2__k21 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_CollectHistogram_SSE2__k22 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_CollectHistogram_SSE2__k23 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_CollectHistogram_SSE2__k24 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}
var F_CollectHistogram_SSE2__k25 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_CollectHistogram_SSE2__k26 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_CollectHistogram_SSE2__k27 = [2]uint64{0x7000700070007, 0x7000700070007}
var F_CollectHistogram_SSE2__k28 = [2]uint64{0x14e808a914e808a9, 0x14e808a914e808a9}
var F_CollectHistogram_SSE2__k29 = [2]uint64{0x12ee000012ee0, 0x12ee000012ee0}
var F_CollectHistogram_SSE2__k30 = [2]uint64{0x1f001f001f001f, 0x1f001f001f001f}
var F_CollectHistogram_SSE2__k31 = [2]uint64{0x8a9eb1808a9eb18, 0x8a9eb1808a9eb18}
var F_CollectHistogram_SSE2__k32 = [2]uint64{0xc7380000c738, 0xc7380000c738}

func F_CollectHistogram_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v140 int32
	_ = v140
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v159 int32
	_ = v159
	var v160 base.V128
	_ = v160
	var v162 int32
	_ = v162
	var v164 base.V128
	_ = v164
	var v165 base.V128
	_ = v165
	var v169 int32
	_ = v169
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v186 int32
	_ = v186
	var v187 int32
	_ = v187
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v215 int32
	_ = v215
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v257 int32
	_ = v257
	var v258 base.V128
	_ = v258
	var v262 base.V128
	_ = v262
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v305 int32
	_ = v305
	var v306 int32
	_ = v306
	var v312 int32
	_ = v312
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v322 int32
	_ = v322
	var v325 int32
	_ = v325
	var v326 int32
	_ = v326
	var v332 int32
	_ = v332
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v342 int32
	_ = v342
	var v345 int32
	_ = v345
	var v346 int32
	_ = v346
	var v353 int32
	_ = v353
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v457 int32
	_ = v457
	var v459 int32
	_ = v459
	var v460 int32
	_ = v460
	var v464 int32
	_ = v464
	var v466 int32
	_ = v466
	var v467 int32
	_ = v467
	var v468 int32
	_ = v468
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v472 int32
	_ = v472
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v476 int32
	_ = v476
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v480 int32
	_ = v480
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v484 int32
	_ = v484
	var v486 int32
	_ = v486
	var v487 int32
	_ = v487
	var v488 int32
	_ = v488
	var v490 int32
	_ = v490
	var v491 int32
	_ = v491
	var v492 int32
	_ = v492
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v496 int32
	_ = v496
	var v498 int32
	_ = v498
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v502 int32
	_ = v502
	var v503 int32
	_ = v503
	var v504 int32
	_ = v504
	var v506 int32
	_ = v506
	var v507 int32
	_ = v507
	var v508 int32
	_ = v508
	var v510 int32
	_ = v510
	var v511 int32
	_ = v511
	var v512 int32
	_ = v512
	var v514 int32
	_ = v514
	var v515 int32
	_ = v515
	var v516 int32
	_ = v516
	var v518 int32
	_ = v518
	var v519 int32
	_ = v519
	var v520 int32
	_ = v520
	var v522 int32
	_ = v522
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v534 int32
	_ = v534
	var v535 int32
	_ = v535
	var v536 int32
	_ = v536
	var v538 int32
	_ = v538
	var v539 int32
	_ = v539
	var v540 int32
	_ = v540
	var v542 int32
	_ = v542
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v552 int32
	_ = v552
	var v554 int32
	_ = v554
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v558 int32
	_ = v558
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v562 int32
	_ = v562
	var v563 int32
	_ = v563
	var v564 int32
	_ = v564
	var v566 int32
	_ = v566
	var v567 int32
	_ = v567
	var v568 int32
	_ = v568
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v572 int32
	_ = v572
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v576 int32
	_ = v576
	var v578 int32
	_ = v578
	var v579 int32
	_ = v579
	var v580 int32
	_ = v580
	var v582 int32
	_ = v582
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v588 int32
	_ = v588
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v593 int32
	_ = v593
	var v594 int32
	_ = v594
	var v596 int32
	_ = v596
	var v597 int32
	_ = v597
	var v599 int32
	_ = v599
	var v600 int32
	_ = v600
	var v602 int32
	_ = v602
	var v603 int32
	_ = v603
	var v605 int32
	_ = v605
	var v606 int32
	_ = v606
	var v608 int32
	_ = v608
	var v609 int32
	_ = v609
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v617 int32
	_ = v617
	var v618 int32
	_ = v618
	var v620 int32
	_ = v620
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v624 int32
	_ = v624
	var v626 int32
	_ = v626
	var v627 int32
	_ = v627
	var v629 int32
	_ = v629
	var v630 int32
	_ = v630
	var v632 int32
	_ = v632
	var v633 int32
	_ = v633
	var v635 int32
	_ = v635
	var v636 int32
	_ = v636
	var v638 int32
	_ = v638
	var v639 int32
	_ = v639
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v647 int32
	_ = v647
	var v648 int32
	_ = v648
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v656 int32
	_ = v656
	var v657 int32
	_ = v657
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v662 int32
	_ = v662
	var v663 int32
	_ = v663
	var v665 int32
	_ = v665
	var v666 int32
	_ = v666
	var v668 int32
	_ = v668
	var v669 int32
	_ = v669
	var v671 int32
	_ = v671
	var v672 int32
	_ = v672
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v680 int32
	_ = v680
	var v681 int32
	_ = v681
	v10 = m.G0
	v12 = v10 - int32(160)
	m.G0 = v12
	base.MemoryFill(m, v12+int32(32), int32(0), int32(128))
	if l3 <= l2 {
	} else {
		v140 = m.G76
		v146 = v140 + l2<<(uint(int32(2))%32)
		v147 = l3 - l2
		for {
			v153 = *(*int32)(unsafe.Add(mBase, uint32(v146)))
			v156 = m.G66
			v157 = *(*int32)(unsafe.Add(mBase, uint32(v156)))
			m.T0[v157].(func(*base.Module, int32, int32, int32))(m, l0+v153, l1+v153, v12)
			mBase = m.M
			v159 = int32(0)
			v160 = base.Simd_g_v128_load(m, v12, v159)
			v162 = int32(3)
			v164 = base.Simd_g_const(&F_CollectHistogram_SSE41__k0)
			v165 = base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_abs(v160), v162), v164)
			base.Simd_g_v128_store(m, v12, v159, v165)
			v169 = v12 + int32(32)
			v173 = int32(2)
			v175 = v169 + base.I32_extend16_s(base.Simd_g_i32x4_extract_lane_l0(v165))<<(uint(v173)%32)
			v176 = *(*int32)(unsafe.Add(mBase, uint32(v175)))
			v177 = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v175))) = v176 + v177
			v186 = v169 + base.Simd_g_i16x8_extract_lane_s_l1(v165)<<(uint(v173)%32)
			v187 = *(*int32)(unsafe.Add(mBase, uint32(v186)))
			*(*int32)(unsafe.Add(mBase, uint32(v186))) = v187 + v177
			v197 = v169 + base.Simd_g_i16x8_extract_lane_s_l2(v165)<<(uint(v173)%32)
			v198 = *(*int32)(unsafe.Add(mBase, uint32(v197)))
			*(*int32)(unsafe.Add(mBase, uint32(v197))) = v198 + v177
			v208 = v169 + base.Simd_g_i16x8_extract_lane_s_l3(v165)<<(uint(v173)%32)
			v209 = *(*int32)(unsafe.Add(mBase, uint32(v208)))
			*(*int32)(unsafe.Add(mBase, uint32(v208))) = v209 + v177
			v215 = int32(4)
			v219 = v169 + base.Simd_g_i16x8_extract_lane_s_l4(v165)<<(uint(v173)%32)
			v220 = *(*int32)(unsafe.Add(mBase, uint32(v219)))
			*(*int32)(unsafe.Add(mBase, uint32(v219))) = v220 + v177
			v230 = v169 + base.Simd_g_i16x8_extract_lane_s_l5(v165)<<(uint(v173)%32)
			v231 = *(*int32)(unsafe.Add(mBase, uint32(v230)))
			*(*int32)(unsafe.Add(mBase, uint32(v230))) = v231 + v177
			v241 = v169 + base.Simd_g_i16x8_extract_lane_s_l6(v165)<<(uint(v173)%32)
			v242 = *(*int32)(unsafe.Add(mBase, uint32(v241)))
			*(*int32)(unsafe.Add(mBase, uint32(v241))) = v242 + v177
			v252 = v169 + base.Simd_g_i16x8_extract_lane_s_l7(v165)<<(uint(v173)%32)
			v253 = *(*int32)(unsafe.Add(mBase, uint32(v252)))
			*(*int32)(unsafe.Add(mBase, uint32(v252))) = v253 + v177
			v257 = int32(16)
			v258 = base.Simd_g_v128_load(m, v12, v257)
			v262 = base.Simd_g_i16x8_min_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_abs(v258), v162), v164)
			base.Simd_g_v128_store(m, v12, v257, v262)
			v272 = v169 + base.I32_extend16_s(base.Simd_g_i32x4_extract_lane_l0(v262))<<(uint(v173)%32)
			v273 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
			*(*int32)(unsafe.Add(mBase, uint32(v272))) = v273 + v177
			v283 = v169 + base.Simd_g_i16x8_extract_lane_s_l1(v262)<<(uint(v173)%32)
			v284 = *(*int32)(unsafe.Add(mBase, uint32(v283)))
			*(*int32)(unsafe.Add(mBase, uint32(v283))) = v284 + v177
			v294 = v169 + base.Simd_g_i16x8_extract_lane_s_l2(v262)<<(uint(v173)%32)
			v295 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
			*(*int32)(unsafe.Add(mBase, uint32(v294))) = v295 + v177
			v305 = v169 + base.Simd_g_i16x8_extract_lane_s_l3(v262)<<(uint(v173)%32)
			v306 = *(*int32)(unsafe.Add(mBase, uint32(v305)))
			*(*int32)(unsafe.Add(mBase, uint32(v305))) = v306 + v177
			v312 = int32(*(*int16)(unsafe.Add(mBase, uint32(v12)+24)))
			v315 = v169 + v312<<(uint(v173)%32)
			v316 = *(*int32)(unsafe.Add(mBase, uint32(v315)))
			*(*int32)(unsafe.Add(mBase, uint32(v315))) = v316 + v177
			v322 = int32(*(*int16)(unsafe.Add(mBase, uint32(v12)+26)))
			v325 = v169 + v322<<(uint(v173)%32)
			v326 = *(*int32)(unsafe.Add(mBase, uint32(v325)))
			*(*int32)(unsafe.Add(mBase, uint32(v325))) = v326 + v177
			v332 = int32(*(*int16)(unsafe.Add(mBase, uint32(v12)+28)))
			v335 = v169 + v332<<(uint(v173)%32)
			v336 = *(*int32)(unsafe.Add(mBase, uint32(v335)))
			*(*int32)(unsafe.Add(mBase, uint32(v335))) = v336 + v177
			v342 = int32(*(*int16)(unsafe.Add(mBase, uint32(v12)+30)))
			v345 = v169 + v342<<(uint(v173)%32)
			v346 = *(*int32)(unsafe.Add(mBase, uint32(v345)))
			*(*int32)(unsafe.Add(mBase, uint32(v345))) = v346 + v177
			v353 = v147 + int32(-1)
			if v353 != 0 {
				v146 = v146 + v215
				v147 = v353
				continue
			} else {
				break
			}
			break
		}
	}
	v364 = v12 + int32(32)
	v365 = int32(0)
	v457 = *(*int32)(unsafe.Add(mBase, uint32(v364)+4))
	v459 = base.B2i32(v365 < v457)
	v460 = *(*int32)(unsafe.Add(mBase, uint32(v364)))
	v464 = *(*int32)(unsafe.Add(mBase, uint32(v364)+8))
	v466 = base.B2i32(v365 < v464)
	if v365 < v464 {
		v467 = int32(2)
	} else {
		v467 = v459 | base.B2i32(v460 < int32(1))
	}
	v468 = *(*int32)(unsafe.Add(mBase, uint32(v364)+12))
	v470 = base.B2i32(int32(0) < v468)
	if int32(0) < v468 {
		v471 = int32(3)
	} else {
		v471 = v467
	}
	v472 = *(*int32)(unsafe.Add(mBase, uint32(v364)+16))
	v474 = base.B2i32(int32(0) < v472)
	if int32(0) < v472 {
		v475 = int32(4)
	} else {
		v475 = v471
	}
	v476 = *(*int32)(unsafe.Add(mBase, uint32(v364)+20))
	v478 = base.B2i32(int32(0) < v476)
	if int32(0) < v476 {
		v479 = int32(5)
	} else {
		v479 = v475
	}
	v480 = *(*int32)(unsafe.Add(mBase, uint32(v364)+24))
	v482 = base.B2i32(int32(0) < v480)
	if int32(0) < v480 {
		v483 = int32(6)
	} else {
		v483 = v479
	}
	v484 = *(*int32)(unsafe.Add(mBase, uint32(v364)+28))
	v486 = base.B2i32(int32(0) < v484)
	if int32(0) < v484 {
		v487 = int32(7)
	} else {
		v487 = v483
	}
	v488 = *(*int32)(unsafe.Add(mBase, uint32(v364)+32))
	v490 = base.B2i32(int32(0) < v488)
	if int32(0) < v488 {
		v491 = int32(8)
	} else {
		v491 = v487
	}
	v492 = *(*int32)(unsafe.Add(mBase, uint32(v364)+36))
	v494 = base.B2i32(int32(0) < v492)
	if int32(0) < v492 {
		v495 = int32(9)
	} else {
		v495 = v491
	}
	v496 = *(*int32)(unsafe.Add(mBase, uint32(v364)+40))
	v498 = base.B2i32(int32(0) < v496)
	if int32(0) < v496 {
		v499 = int32(10)
	} else {
		v499 = v495
	}
	v500 = *(*int32)(unsafe.Add(mBase, uint32(v364)+44))
	v502 = base.B2i32(int32(0) < v500)
	if int32(0) < v500 {
		v503 = int32(11)
	} else {
		v503 = v499
	}
	v504 = *(*int32)(unsafe.Add(mBase, uint32(v364)+48))
	v506 = base.B2i32(int32(0) < v504)
	if int32(0) < v504 {
		v507 = int32(12)
	} else {
		v507 = v503
	}
	v508 = *(*int32)(unsafe.Add(mBase, uint32(v364)+52))
	v510 = base.B2i32(int32(0) < v508)
	if int32(0) < v508 {
		v511 = int32(13)
	} else {
		v511 = v507
	}
	v512 = *(*int32)(unsafe.Add(mBase, uint32(v364)+56))
	v514 = base.B2i32(int32(0) < v512)
	if int32(0) < v512 {
		v515 = int32(14)
	} else {
		v515 = v511
	}
	v516 = *(*int32)(unsafe.Add(mBase, uint32(v364)+60))
	v518 = base.B2i32(int32(0) < v516)
	if int32(0) < v516 {
		v519 = int32(15)
	} else {
		v519 = v515
	}
	v520 = *(*int32)(unsafe.Add(mBase, uint32(v364)+64))
	v522 = base.B2i32(int32(0) < v520)
	if int32(0) < v520 {
		v523 = int32(16)
	} else {
		v523 = v519
	}
	v524 = *(*int32)(unsafe.Add(mBase, uint32(v364)+68))
	v526 = base.B2i32(int32(0) < v524)
	if int32(0) < v524 {
		v527 = int32(17)
	} else {
		v527 = v523
	}
	v528 = *(*int32)(unsafe.Add(mBase, uint32(v364)+72))
	v530 = base.B2i32(int32(0) < v528)
	if int32(0) < v528 {
		v531 = int32(18)
	} else {
		v531 = v527
	}
	v532 = *(*int32)(unsafe.Add(mBase, uint32(v364)+76))
	v534 = base.B2i32(int32(0) < v532)
	if int32(0) < v532 {
		v535 = int32(19)
	} else {
		v535 = v531
	}
	v536 = *(*int32)(unsafe.Add(mBase, uint32(v364)+80))
	v538 = base.B2i32(int32(0) < v536)
	if int32(0) < v536 {
		v539 = int32(20)
	} else {
		v539 = v535
	}
	v540 = *(*int32)(unsafe.Add(mBase, uint32(v364)+84))
	v542 = base.B2i32(int32(0) < v540)
	if int32(0) < v540 {
		v543 = int32(21)
	} else {
		v543 = v539
	}
	v544 = *(*int32)(unsafe.Add(mBase, uint32(v364)+88))
	v546 = base.B2i32(int32(0) < v544)
	if int32(0) < v544 {
		v547 = int32(22)
	} else {
		v547 = v543
	}
	v548 = *(*int32)(unsafe.Add(mBase, uint32(v364)+92))
	v550 = base.B2i32(int32(0) < v548)
	if int32(0) < v548 {
		v551 = int32(23)
	} else {
		v551 = v547
	}
	v552 = *(*int32)(unsafe.Add(mBase, uint32(v364)+96))
	v554 = base.B2i32(int32(0) < v552)
	if int32(0) < v552 {
		v555 = int32(24)
	} else {
		v555 = v551
	}
	v556 = *(*int32)(unsafe.Add(mBase, uint32(v364)+100))
	v558 = base.B2i32(int32(0) < v556)
	if int32(0) < v556 {
		v559 = int32(25)
	} else {
		v559 = v555
	}
	v560 = *(*int32)(unsafe.Add(mBase, uint32(v364)+104))
	v562 = base.B2i32(int32(0) < v560)
	if int32(0) < v560 {
		v563 = int32(26)
	} else {
		v563 = v559
	}
	v564 = *(*int32)(unsafe.Add(mBase, uint32(v364)+108))
	v566 = base.B2i32(int32(0) < v564)
	if int32(0) < v564 {
		v567 = int32(27)
	} else {
		v567 = v563
	}
	v568 = *(*int32)(unsafe.Add(mBase, uint32(v364)+112))
	v570 = base.B2i32(int32(0) < v568)
	if int32(0) < v568 {
		v571 = int32(28)
	} else {
		v571 = v567
	}
	v572 = *(*int32)(unsafe.Add(mBase, uint32(v364)+116))
	v574 = base.B2i32(int32(0) < v572)
	if int32(0) < v572 {
		v575 = int32(29)
	} else {
		v575 = v571
	}
	v576 = *(*int32)(unsafe.Add(mBase, uint32(v364)+120))
	v578 = base.B2i32(int32(0) < v576)
	if int32(0) < v576 {
		v579 = int32(30)
	} else {
		v579 = v575
	}
	v580 = *(*int32)(unsafe.Add(mBase, uint32(v364)+124))
	v582 = base.B2i32(int32(0) < v580)
	if int32(0) < v580 {
		v583 = int32(31)
	} else {
		v583 = v579
	}
	*(*int32)(unsafe.Add(mBase, uint32(l4)+4)) = v583
	v585 = int32(0)
	if v585 < v460 {
		v588 = v460
	} else {
		v588 = v585
	}
	if v588 < v457 {
		v590 = v457
	} else {
		v590 = v588
	}
	if v365 < v457 {
		v591 = v590
	} else {
		v591 = v588
	}
	if v591 < v464 {
		v593 = v464
	} else {
		v593 = v591
	}
	if v365 < v464 {
		v594 = v593
	} else {
		v594 = v591
	}
	if v594 < v468 {
		v596 = v468
	} else {
		v596 = v594
	}
	if int32(0) < v468 {
		v597 = v596
	} else {
		v597 = v594
	}
	if v597 < v472 {
		v599 = v472
	} else {
		v599 = v597
	}
	if int32(0) < v472 {
		v600 = v599
	} else {
		v600 = v597
	}
	if v600 < v476 {
		v602 = v476
	} else {
		v602 = v600
	}
	if int32(0) < v476 {
		v603 = v602
	} else {
		v603 = v600
	}
	if v603 < v480 {
		v605 = v480
	} else {
		v605 = v603
	}
	if int32(0) < v480 {
		v606 = v605
	} else {
		v606 = v603
	}
	if v606 < v484 {
		v608 = v484
	} else {
		v608 = v606
	}
	if int32(0) < v484 {
		v609 = v608
	} else {
		v609 = v606
	}
	if v609 < v488 {
		v611 = v488
	} else {
		v611 = v609
	}
	if int32(0) < v488 {
		v612 = v611
	} else {
		v612 = v609
	}
	if v612 < v492 {
		v614 = v492
	} else {
		v614 = v612
	}
	if int32(0) < v492 {
		v615 = v614
	} else {
		v615 = v612
	}
	if v615 < v496 {
		v617 = v496
	} else {
		v617 = v615
	}
	if int32(0) < v496 {
		v618 = v617
	} else {
		v618 = v615
	}
	if v618 < v500 {
		v620 = v500
	} else {
		v620 = v618
	}
	if int32(0) < v500 {
		v621 = v620
	} else {
		v621 = v618
	}
	if v621 < v504 {
		v623 = v504
	} else {
		v623 = v621
	}
	if int32(0) < v504 {
		v624 = v623
	} else {
		v624 = v621
	}
	if v624 < v508 {
		v626 = v508
	} else {
		v626 = v624
	}
	if int32(0) < v508 {
		v627 = v626
	} else {
		v627 = v624
	}
	if v627 < v512 {
		v629 = v512
	} else {
		v629 = v627
	}
	if int32(0) < v512 {
		v630 = v629
	} else {
		v630 = v627
	}
	if v630 < v516 {
		v632 = v516
	} else {
		v632 = v630
	}
	if int32(0) < v516 {
		v633 = v632
	} else {
		v633 = v630
	}
	if v633 < v520 {
		v635 = v520
	} else {
		v635 = v633
	}
	if int32(0) < v520 {
		v636 = v635
	} else {
		v636 = v633
	}
	if v636 < v524 {
		v638 = v524
	} else {
		v638 = v636
	}
	if int32(0) < v524 {
		v639 = v638
	} else {
		v639 = v636
	}
	if v639 < v528 {
		v641 = v528
	} else {
		v641 = v639
	}
	if int32(0) < v528 {
		v642 = v641
	} else {
		v642 = v639
	}
	if v642 < v532 {
		v644 = v532
	} else {
		v644 = v642
	}
	if int32(0) < v532 {
		v645 = v644
	} else {
		v645 = v642
	}
	if v645 < v536 {
		v647 = v536
	} else {
		v647 = v645
	}
	if int32(0) < v536 {
		v648 = v647
	} else {
		v648 = v645
	}
	if v648 < v540 {
		v650 = v540
	} else {
		v650 = v648
	}
	if int32(0) < v540 {
		v651 = v650
	} else {
		v651 = v648
	}
	if v651 < v544 {
		v653 = v544
	} else {
		v653 = v651
	}
	if int32(0) < v544 {
		v654 = v653
	} else {
		v654 = v651
	}
	if v654 < v548 {
		v656 = v548
	} else {
		v656 = v654
	}
	if int32(0) < v548 {
		v657 = v656
	} else {
		v657 = v654
	}
	if v657 < v552 {
		v659 = v552
	} else {
		v659 = v657
	}
	if int32(0) < v552 {
		v660 = v659
	} else {
		v660 = v657
	}
	if v660 < v556 {
		v662 = v556
	} else {
		v662 = v660
	}
	if int32(0) < v556 {
		v663 = v662
	} else {
		v663 = v660
	}
	if v663 < v560 {
		v665 = v560
	} else {
		v665 = v663
	}
	if int32(0) < v560 {
		v666 = v665
	} else {
		v666 = v663
	}
	if v666 < v564 {
		v668 = v564
	} else {
		v668 = v666
	}
	if int32(0) < v564 {
		v669 = v668
	} else {
		v669 = v666
	}
	if v669 < v568 {
		v671 = v568
	} else {
		v671 = v669
	}
	if int32(0) < v568 {
		v672 = v671
	} else {
		v672 = v669
	}
	if v672 < v572 {
		v674 = v572
	} else {
		v674 = v672
	}
	if int32(0) < v572 {
		v675 = v674
	} else {
		v675 = v672
	}
	if v675 < v576 {
		v677 = v576
	} else {
		v677 = v675
	}
	if int32(0) < v576 {
		v678 = v677
	} else {
		v678 = v675
	}
	if v678 < v580 {
		v680 = v580
	} else {
		v680 = v678
	}
	if int32(0) < v580 {
		v681 = v680
	} else {
		v681 = v678
	}
	*(*int32)(unsafe.Add(mBase, uint32(l4))) = v681
	m.G0 = v12 + int32(160)
	return
}

var F_CollectHistogram_SSE41__k0 = [2]uint64{0x1f001f001f001f, 0x1f001f001f001f}

func F_CombinedShannonEntropy_SSE2(m *base.Module, l0 int32, l1 int32) int64 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v21 int64
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 base.V128
	_ = v37
	var v38 int32
	_ = v38
	var v39 base.V128
	_ = v39
	var v41 int32
	_ = v41
	var v42 base.V128
	_ = v42
	var v43 int32
	_ = v43
	var v44 base.V128
	_ = v44
	var v47 base.V128
	_ = v47
	var v50 int32
	_ = v50
	var v52 base.V128
	_ = v52
	var v54 base.V128
	_ = v54
	var v57 base.V128
	_ = v57
	var v59 base.V128
	_ = v59
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v69 int64
	_ = v69
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v81 int32
	_ = v81
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v97 int64
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v100 int64
	_ = v100
	var v101 int64
	_ = v101
	var v104 int64
	_ = v104
	var v105 int32
	_ = v105
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v122 int64
	_ = v122
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v125 int64
	_ = v125
	var v126 int64
	_ = v126
	var v127 int32
	_ = v127
	var v128 int64
	_ = v128
	var v131 int32
	_ = v131
	var v134 int64
	_ = v134
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v152 int32
	_ = v152
	var v156 int64
	_ = v156
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v159 int64
	_ = v159
	var v160 int64
	_ = v160
	var v163 int32
	_ = v163
	var v167 int64
	_ = v167
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v170 int64
	_ = v170
	var v171 int64
	_ = v171
	v4 = int32(0)
	v21 = int64(0)
	v22 = v4
	v23 = v4
	v24 = v4
	for {
		v34 = v22 << (uint(int32(2)) % 32)
		v35 = l1 + v34
		v36 = int32(0)
		v37 = base.Simd_g_v128_load_rng(m, v35, v36, int32(0), int32(64))
		v38 = int32(16)
		v39 = base.Simd_g_v128_load_nc(m, v35, v38)
		v41 = int32(32)
		v42 = base.Simd_g_v128_load_nc(m, v35, v41)
		v43 = int32(48)
		v44 = base.Simd_g_v128_load_nc(m, v35, v43)
		v47 = base.Simd_g_const(&F_CombinedShannonEntropy_SSE2__k0)
		v50 = l0 + v34
		v52 = base.Simd_g_v128_load_rng(m, v50, v36, int32(0), int32(64))
		v54 = base.Simd_g_v128_load_nc(m, v50, v38)
		v57 = base.Simd_g_v128_load_nc(m, v50, v41)
		v59 = base.Simd_g_v128_load_nc(m, v50, v43)
		v63 = base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_gt_s(base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_narrow_i32x4_s(v52, v54), base.Simd_g_i16x8_narrow_i32x4_s(v57, v59)), v47))
		v64 = base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_gt_s(base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_narrow_i32x4_s(v37, v39), base.Simd_g_i16x8_narrow_i32x4_s(v42, v44)), v47)) | v63
		if v64 == v36 {
			v134 = v21
			v136 = v23
			v137 = v24
		} else {
			v69 = v21
			v71 = v23
			v72 = v24
			v74 = v64
			for {
				v81 = base.I32_ctz(v74)
				if int32(base.Ui32(v63)>>(uint(v81)%32))&int32(1) == int32(0) {
					v104 = v69
					v105 = v72
				} else {
					v90 = *(*int32)(unsafe.Add(mBase, uint32(v50+v81<<(uint(int32(2))%32))))
					if base.Ui32(int32(255)) < base.Ui32(v90) {
						v98 = m.G110
						v99 = *(*int32)(unsafe.Add(mBase, uint32(v98)))
						v100 = m.T0[v99].(func(*base.Module, int32) int64)(m, v90)
						mBase = m.M
						v101 = v100
					} else {
						v93 = m.G109
						v97 = *(*int64)(unsafe.Add(mBase, uint32(v93+v90<<(uint(int32(3))%32))))
						v101 = v97
					}
					v104 = v101 + v69
					v105 = v90 + v72
				}
				v110 = (v81 + v22) << (uint(int32(2)) % 32)
				v112 = *(*int32)(unsafe.Add(mBase, uint32(l1+v110)))
				v114 = *(*int32)(unsafe.Add(mBase, uint32(l0+v110)))
				v115 = v112 + v114
				if base.Ui32(int32(255)) < base.Ui32(v115) {
					v123 = m.G110
					v124 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
					v125 = m.T0[v124].(func(*base.Module, int32) int64)(m, v115)
					mBase = m.M
					v126 = v125
				} else {
					v118 = m.G109
					v122 = *(*int64)(unsafe.Add(mBase, uint32(v118+v115<<(uint(int32(3))%32))))
					v126 = v122
				}
				v127 = v115 + v71
				v128 = v126 + v104
				v131 = (v74 + int32(-1)) & v74
				if v131 != 0 {
					v69 = v128
					v71 = v127
					v72 = v105
					v74 = v131
					continue
				} else {
					break
				}
				break
			}
			v134 = v128
			v136 = v127
			v137 = v105
		}
		if base.Ui32(v22) < base.Ui32(int32(240)) {
			v21 = v134
			v22 = v22 + int32(16)
			v23 = v136
			v24 = v137
			continue
		} else {
			break
		}
		break
	}
	if base.Ui32(int32(255)) < base.Ui32(v136) {
		v157 = m.G110
		v158 = *(*int32)(unsafe.Add(mBase, uint32(v157)))
		v159 = m.T0[v158].(func(*base.Module, int32) int64)(m, v136)
		mBase = m.M
		v160 = v159
	} else {
		v152 = m.G109
		v156 = *(*int64)(unsafe.Add(mBase, uint32(v152+v136<<(uint(int32(3))%32))))
		v160 = v156
	}
	if base.Ui32(int32(255)) < base.Ui32(v137) {
		v168 = m.G110
		v169 = *(*int32)(unsafe.Add(mBase, uint32(v168)))
		v170 = m.T0[v169].(func(*base.Module, int32) int64)(m, v137)
		mBase = m.M
		v171 = v170
	} else {
		v163 = m.G109
		v167 = *(*int64)(unsafe.Add(mBase, uint32(v163+v137<<(uint(int32(3))%32))))
		v171 = v167
	}
	return v160 - v134 + v171
}

var F_CombinedShannonEntropy_SSE2__k0 = [2]uint64{0x0, 0x0}

func F_CompressAlphaJob(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v41 int32
	_ = v41
	var v49 int32
	_ = v49
	var v54 int64
	_ = v54
	var v55 int32
	_ = v55
	var v62 int64
	_ = v62
	var v63 int32
	_ = v63
	var v64 int64
	_ = v64
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v99 int32
	_ = v99
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
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v157 int32
	_ = v157
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v315 int32
	_ = v315
	var v329 int32
	_ = v329
	var v336 int32
	_ = v336
	var v345 int32
	_ = v345
	var v351 int32
	_ = v351
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v367 base.V128
	_ = v367
	var v369 base.V128
	_ = v369
	var v373 int32
	_ = v373
	var v384 int32
	_ = v384
	var v391 int32
	_ = v391
	var v401 base.V128
	_ = v401
	var v433 base.V128
	_ = v433
	var v465 base.V128
	_ = v465
	var v495 int32
	_ = v495
	var v500 int32
	_ = v500
	var v522 int32
	_ = v522
	var v544 int32
	_ = v544
	var v546 int32
	_ = v546
	var v549 int32
	_ = v549
	var v575 int32
	_ = v575
	var v615 int32
	_ = v615
	var v621 base.V128
	_ = v621
	var v624 int32
	_ = v624
	var v625 int32
	_ = v625
	var v626 base.V128
	_ = v626
	var v627 base.V128
	_ = v627
	var v629 base.V128
	_ = v629
	var v631 base.V128
	_ = v631
	var v637 base.V128
	_ = v637
	var v646 base.V128
	_ = v646
	var v655 base.V128
	_ = v655
	var v660 base.V128
	_ = v660
	var v662 int32
	_ = v662
	var v665 int32
	_ = v665
	var v668 base.V128
	_ = v668
	var v673 int32
	_ = v673
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v681 int32
	_ = v681
	var v685 base.V128
	_ = v685
	var v688 base.V128
	_ = v688
	var v721 int32
	_ = v721
	var v728 int32
	_ = v728
	var v734 int32
	_ = v734
	var v738 base.V128
	_ = v738
	var v741 base.V128
	_ = v741
	var v791 int32
	_ = v791
	var v796 int32
	_ = v796
	var v800 int32
	_ = v800
	var v804 int32
	_ = v804
	var v808 base.V128
	_ = v808
	var v811 base.V128
	_ = v811
	var v848 int32
	_ = v848
	var v862 int32
	_ = v862
	var v867 int64
	_ = v867
	var v874 int64
	_ = v874
	var v875 int32
	_ = v875
	var v876 int64
	_ = v876
	var v886 int32
	_ = v886
	var v888 int32
	_ = v888
	var v893 int32
	_ = v893
	var v898 int32
	_ = v898
	var v905 int32
	_ = v905
	var v917 int32
	_ = v917
	var v924 int32
	_ = v924
	var v927 int32
	_ = v927
	var v928 int32
	_ = v928
	var v935 int32
	_ = v935
	var v936 int32
	_ = v936
	var v938 base.V128
	_ = v938
	var v939 int32
	_ = v939
	var v950 int32
	_ = v950
	var v956 int32
	_ = v956
	var v957 int32
	_ = v957
	var v959 base.V128
	_ = v959
	var v960 int32
	_ = v960
	var v966 int32
	_ = v966
	var v969 int32
	_ = v969
	var v985 int32
	_ = v985
	var v993 int32
	_ = v993
	var v1002 int32
	_ = v1002
	var v1005 base.V128
	_ = v1005
	var v1009 base.V128
	_ = v1009
	var v1014 int32
	_ = v1014
	var v1017 int32
	_ = v1017
	var v1024 int32
	_ = v1024
	var v1025 int32
	_ = v1025
	var v1027 base.V128
	_ = v1027
	var v1028 int32
	_ = v1028
	var v1057 int32
	_ = v1057
	var v1060 int32
	_ = v1060
	var v1063 int32
	_ = v1063
	var v1065 int32
	_ = v1065
	var v1066 int32
	_ = v1066
	var v1072 int32
	_ = v1072
	var v1073 int32
	_ = v1073
	var v1075 int32
	_ = v1075
	var v1085 int32
	_ = v1085
	var v1088 int64
	_ = v1088
	var v1090 int32
	_ = v1090
	var v1093 int32
	_ = v1093
	var v1095 int32
	_ = v1095
	var v1096 int32
	_ = v1096
	var v1102 int32
	_ = v1102
	var v1103 int32
	_ = v1103
	var v1105 int32
	_ = v1105
	var v1116 int32
	_ = v1116
	var v1118 int32
	_ = v1118
	var v1119 int32
	_ = v1119
	var v1125 int32
	_ = v1125
	var v1126 int32
	_ = v1126
	var v1142 int32
	_ = v1142
	v23 = m.G0
	v25 = v23 - int32(496)
	m.G0 = v25
	v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v28 = *(*int32)(unsafe.Add(mBase, uint32(v27)+48))
	v29 = *(*int32)(unsafe.Add(mBase, uint32(v27)+8))
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v31 = *(*int32)(unsafe.Add(mBase, uint32(v30)+12))
	v32 = *(*int32)(unsafe.Add(mBase, uint32(v30)+8))
	v33 = *(*int32)(unsafe.Add(mBase, uint32(v27)+56))
	v34 = *(*int32)(unsafe.Add(mBase, uint32(v27)+52))
	*(*int64)(unsafe.Add(mBase, uint32(v25)+8)) = int64(0)
	if base.Ui32(v33) < base.Ui32(int32(101)) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	if v1119 != 0 {
		goto L163
	} else {
		goto L164
	}
L2:
	;
	if base.Ui32(v28) < base.Ui32(int32(2)) {
		goto L7
	} else {
		goto L8
	}
L3:
	;
	v41 = *(*int32)(unsafe.Add(mBase, uint32(v30)+92))
	if v41 != 0 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	v1116 = l0
	v1118 = v25
	v1119 = int32(0)
	v1125 = int32(0)
	v1126 = int32(0)
	goto L1
L5:
	;
	goto L4
L6:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v30)+92)) = int32(4)
	goto L5
L7:
	;
	v54 = int64(1)
	v55 = v31 * v32
	goto L16
L8:
	;
	v49 = *(*int32)(unsafe.Add(mBase, uint32(v30)+92))
	if v49 != 0 {
		goto L10
	} else {
		goto L11
	}
L9:
	;
	v1116 = l0
	v1118 = v25
	v1119 = int32(0)
	v1125 = int32(0)
	v1126 = int32(0)
	goto L1
L10:
	;
	goto L9
L11:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v30)+92)) = int32(4)
	goto L10
L12:
	;
	v83 = *(*int32)(unsafe.Add(mBase, uint32(v30)+36))
	v84 = *(*int32)(unsafe.Add(mBase, uint32(v30)+40))
	if v31 < int32(1) {
		goto L24
	} else {
		goto L25
	}
L13:
	;
	if v76 != 0 {
		goto L12
	} else {
		goto L19
	}
L14:
	;
	goto L13
L15:
	;
	v74 = F_malloc(m, base.I32_wrap_i64(v54)*v55)
	mBase = m.M
	v76 = v74
	goto L14
L16:
	;
	v62 = base.I64_div_u_s(int64(2147418112), v54)
	v63 = int32(0)
	v64 = base.I64_extend_i32_u(v55)
	if base.Ui64(int64(4294967295)) < base.Ui64(v64*v54) {
		v76 = v63
		goto L14
	} else {
		goto L17
	}
L17:
	;
	if base.Ui64(v62) < base.Ui64(v64) {
		v76 = v63
		goto L14
	} else {
		goto L18
	}
L18:
	;
	goto L15
L19:
	;
	v79 = *(*int32)(unsafe.Add(mBase, uint32(v30)+92))
	if v79 != 0 {
		goto L21
	} else {
		goto L22
	}
L20:
	;
	v1116 = l0
	v1118 = v25
	v1119 = int32(0)
	v1125 = int32(0)
	v1126 = int32(0)
	goto L1
L21:
	;
	goto L20
L22:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v30)+92)) = int32(1)
	goto L21
L23:
	;
	if int32(99) < v33 {
		goto L37
	} else {
		goto L38
	}
L24:
	;
	goto L23
L25:
	;
	v90 = v31 & int32(3)
	if v90 != 0 {
		goto L27
	} else {
		goto L28
	}
L26:
	;
	if base.Ui32(v31) < base.Ui32(int32(4)) {
		goto L24
	} else {
		goto L32
	}
L27:
	;
	v93 = v83
	v95 = v76
	v99 = v90
	goto L29
L28:
	;
	v106 = v83
	v108 = v76
	v113 = v31
	goto L26
L29:
	;
	v101 = F_memcpy(m, v95, v93, v32)
	mBase = m.M
	v102 = v101 + v32
	v103 = v93 + v84
	v105 = v99 + int32(-1)
	if v105 != 0 {
		v93 = v103
		v95 = v102
		v99 = v105
		goto L29
	} else {
		goto L31
	}
L30:
	;
	v106 = v103
	v108 = v102
	v113 = v31 & int32(-4)
	goto L26
L31:
	;
	goto L30
L32:
	;
	v118 = v106
	v120 = v108
	v124 = v113 + int32(-1)
	goto L33
L33:
	;
	v126 = F_memcpy(m, v120, v118, v32)
	mBase = m.M
	v128 = v118 + v84
	v129 = F_memcpy(m, v126+v32, v128, v32)
	mBase = m.M
	v131 = v128 + v84
	v132 = F_memcpy(m, v129+v32, v131, v32)
	mBase = m.M
	v134 = v131 + v84
	v135 = F_memcpy(m, v132+v32, v134, v32)
	mBase = m.M
	v139 = v124 + int32(-4)
	if base.Ui32(v139) < base.Ui32(int32(-2)) {
		v118 = v134 + v84
		v120 = v135 + v32
		v124 = v139
		goto L33
	} else {
		goto L35
	}
L34:
	;
	goto L24
L35:
	;
	goto L34
L36:
	;
	F_free(m, v1105)
	mBase = m.M
	goto L161
L37:
	;
	v173 = base.B2i32(v33 < int32(100))
	F_VP8FiltersInit(m)
	mBase = m.M
	v175 = *(*int32)(unsafe.Add(mBase, uint32(v30)+88))
	if v34 == int32(1) {
		goto L51
	} else {
		goto L52
	}
L38:
	;
	if base.Ui32(int32(70)) < base.Ui32(v33) {
		goto L40
	} else {
		goto L41
	}
L39:
	;
	v167 = F_QuantizeLevels(m, v76, v32, v31, v164, v25+int32(8))
	mBase = m.M
	if v167 != 0 {
		goto L37
	} else {
		goto L42
	}
L40:
	;
	v164 = v33<<(uint(int32(3))%32) + int32(-544)
	goto L39
L41:
	;
	v157 = base.I32_div_u_s(v33&int32(255), int32(5))
	v164 = v157 + int32(2)
	goto L39
L42:
	;
	v168 = int32(0)
	v1093 = l0
	v1095 = v25
	v1096 = v168
	v1102 = v168
	v1103 = v168
	v1105 = v76
	goto L36
L43:
	;
	v1085 = *(*int32)(unsafe.Add(mBase, uint32(v30)+88))
	if v1085 == int32(0) {
		v1093 = v1063
		v1095 = v1065
		v1096 = v1066
		v1102 = v1072
		v1103 = v1073
		v1105 = v1075
		goto L36
	} else {
		goto L160
	}
L44:
	;
	v1057 = *(*int32)(unsafe.Add(mBase, uint32(v30)+92))
	if v1057 != 0 {
		goto L158
	} else {
		goto L159
	}
L45:
	;
	if v985 == int32(0) {
		goto L150
	} else {
		goto L151
	}
L46:
	;
	v867 = int64(1)
	goto L129
L47:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v25)+16)) = int32(-1)
	v804 = v25 + int32(20)
	v808 = base.Simd_g_const(&F_CompressAlphaJob__k0)
	base.Simd_g_v128_store(m, v804, int32(0), v808)
	v811 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	base.Simd_g_v128_store(m, v804, int32(16), v811)
	goto L117
L48:
	;
	v796 = int32(0)
	v800 = F_EncodeAlphaInternal(m, v76, v32, v31, v28, v796, v173, v29, v796, v25+int32(16))
	mBase = m.M
	v985 = v800
	v993 = v791
	goto L45
L49:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v25)+16)) = int32(-1)
	v734 = v25 + int32(20)
	v738 = base.Simd_g_const(&F_CompressAlphaJob__k0)
	base.Simd_g_v128_store(m, v734, int32(0), v738)
	v741 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	base.Simd_g_v128_store(m, v734, int32(16), v741)
	goto L107
L50:
	;
	v185 = int32(0)
	goto L62
L51:
	;
	v180 = int32(6)
	goto L53
L52:
	;
	v180 = int32(5)
	goto L53
L53:
	;
	if v34 != 0 {
		goto L54
	} else {
		goto L55
	}
L54:
	;
	v182 = v180
	goto L56
L55:
	;
	v182 = int32(0)
	goto L56
L56:
	;
	if v28 != 0 {
		goto L57
	} else {
		goto L58
	}
L57:
	;
	v184 = v182
	goto L59
L58:
	;
	v184 = int32(0)
	goto L59
L59:
	;
	switch v184 {
	case 0:
		goto L49
	default:
		goto L47
	case 6:
		goto L50
	}
L60:
	;
	if v32 < int32(1) {
		goto L74
	} else {
		goto L75
	}
L62:
	;
	base.MemoryFill(m, v25+int32(240), v185, int32(256))
	goto L60
L74:
	;
	v615 = v185
	v621 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	goto L90
L75:
	;
	if v31 < int32(1) {
		goto L74
	} else {
		goto L76
	}
L76:
	;
	v315 = v32 & int32(2147483632)
	v329 = v76
	v336 = int32(0)
	goto L77
L77:
	;
	if base.Ui32(v32) < base.Ui32(int32(16)) {
		v500 = int32(0)
		goto L80
	} else {
		goto L81
	}
L78:
	;
	goto L74
L79:
	;
	v575 = v336 + int32(1)
	if v575 != v31 {
		v329 = v329 + v32
		v336 = v575
		goto L77
	} else {
		goto L89
	}
L80:
	;
	v522 = v500
	goto L86
L81:
	;
	v345 = v329
	v351 = v315
	goto L82
L82:
	;
	v365 = v25 + int32(240)
	v366 = int32(0)
	v367 = base.Simd_g_v128_load(m, v345, v366)
	v369 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v367))
	v373 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l0(v369)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l1(v369)))) = uint8(v373)
	v384 = int32(2)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l2(v369)))) = uint8(v373)
	v391 = int32(3)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l3(v369)))) = uint8(v373)
	v401 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_swizzle_c(v367, base.Simd_g_const(&F_CompressAlphaJob__k2))))
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l0(v401)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l1(v401)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l2(v401)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l3(v401)))) = uint8(v373)
	v433 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_swizzle_c(v367, base.Simd_g_const(&F_CompressAlphaJob__k3))))
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l0(v433)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l1(v433)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l2(v433)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l3(v433)))) = uint8(v373)
	v465 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(base.Simd_g_i8x16_swizzle_c(v367, base.Simd_g_const(&F_CompressAlphaJob__k4))))
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l0(v465)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l1(v465)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l2(v465)))) = uint8(v373)
	*(*uint8)(unsafe.Add(mBase, uint32(v365+base.Simd_g_i32x4_extract_lane_l3(v465)))) = uint8(v373)
	v495 = v351 + int32(-16)
	if v495 != 0 {
		v345 = v345 + int32(16)
		v351 = v495
		goto L82
	} else {
		goto L84
	}
L83:
	;
	if v32 == v315 {
		goto L79
	} else {
		goto L85
	}
L84:
	;
	goto L83
L85:
	;
	v500 = v315
	goto L80
L86:
	;
	v544 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v329+v522))))
	v546 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v25+int32(240)+v544))) = uint8(v546)
	v549 = v522 + v546
	if v32 != v549 {
		v522 = v549
		goto L86
	} else {
		goto L88
	}
L87:
	;
	goto L79
L88:
	;
	goto L87
L89:
	;
	goto L78
L90:
	;
	v624 = v25 + int32(240) + v615
	v625 = int32(0)
	v626 = base.Simd_g_v128_load(m, v624, v625)
	v627 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	v629 = base.Simd_g_const(&F_CompressAlphaJob__k5)
	v631 = base.Simd_g_const(&F_CompressAlphaJob__k6)
	v637 = base.Simd_g_v128_load32_zero(m, v624+int32(4), v625)
	v646 = base.Simd_g_v128_load64_zero(m, v624+int32(8), v625)
	v655 = base.Simd_g_v128_load32_zero(m, v624+int32(12), v625)
	v660 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v621, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_ne(v626, v627), base.Simd_g_const(&F_CompressAlphaJob__k5)), v631)), base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_ne(v637, v627), base.Simd_g_const(&F_CompressAlphaJob__k5)), v631)), base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_ne(v646, v627), base.Simd_g_const(&F_CompressAlphaJob__k5)), v631)), base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_ne(v655, v627), base.Simd_g_const(&F_CompressAlphaJob__k5)), v631))
	v662 = v615 + int32(16)
	if v662 != int32(256) {
		v615 = v662
		v621 = v660
		goto L90
	} else {
		goto L92
	}
L91:
	;
	v665 = int32(0)
	v668 = base.Simd_g_i32x4_add(v660, base.Simd_g_i8x16_swizzle_c(v660, base.Simd_g_const(&F_CompressAlphaJob__k7)))
	v673 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v668, base.Simd_g_i8x16_swizzle_c(v668, base.Simd_g_const(&F_CompressAlphaJob__k8))))
	if base.Ui32(v673) < base.Ui32(int32(17)) {
		v677 = v665
		goto L93
	} else {
		goto L94
	}
L92:
	;
	goto L91
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v25)+16)) = int32(-1)
	v681 = v25 + int32(20)
	v685 = base.Simd_g_const(&F_CompressAlphaJob__k0)
	base.Simd_g_v128_store(m, v681, int32(0), v685)
	v688 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	base.Simd_g_v128_store(m, v681, int32(16), v688)
	goto L96
L94:
	;
	v676 = F_WebPEstimateBestFilter(m, v76, v32, v31, v32)
	mBase = m.M
	v677 = v676
	goto L93
L95:
	;
	v721 = int32(1)
	v728 = v721<<(uint(v677)%32) | (base.B2i32(int32(3) < v29) | base.B2i32(base.Ui32(int32(192)) < base.Ui32(v673)))
	if v728 == v721 {
		v791 = v681
		goto L48
	} else {
		goto L105
	}
L96:
	;
	goto L95
L105:
	;
	v848 = v728
	v862 = v681
	goto L46
L106:
	;
	v791 = v734
	goto L48
L107:
	;
	goto L106
L116:
	;
	v848 = int32(15)
	v862 = v804
	goto L46
L117:
	;
	goto L116
L126:
	;
	if v888 == int32(0) {
		goto L44
	} else {
		goto L132
	}
L127:
	;
	goto L126
L128:
	;
	v886 = F_malloc(m, base.I32_wrap_i64(v867)*v55)
	mBase = m.M
	v888 = v886
	goto L127
L129:
	;
	v874 = base.I64_div_u_s(int64(2147418112), v867)
	v875 = int32(0)
	v876 = base.I64_extend_i32_u(v55)
	if base.Ui64(int64(4294967295)) < base.Ui64(v876*v867) {
		v888 = v875
		goto L127
	} else {
		goto L130
	}
L130:
	;
	if base.Ui64(v874) < base.Ui64(v876) {
		v888 = v875
		goto L127
	} else {
		goto L131
	}
L131:
	;
	goto L128
L132:
	;
	v893 = v25 + int32(244)
	v898 = v848
	v905 = int32(0)
	goto L133
L133:
	;
	v917 = int32(1)
	if v898&v917 == int32(0) {
		v966 = v917
		goto L135
	} else {
		goto L136
	}
L134:
	;
	F_free(m, v888)
	mBase = m.M
	goto L149
L135:
	;
	if base.Ui32(v898) < base.Ui32(int32(2)) {
		goto L146
	} else {
		goto L147
	}
L136:
	;
	v924 = F_EncodeAlphaInternal(m, v76, v32, v31, v28, v905, v173, v29, v888, v25+int32(240))
	mBase = m.M
	if v924 == int32(0) {
		goto L137
	} else {
		goto L138
	}
L137:
	;
	if v893 == int32(0) {
		goto L144
	} else {
		goto L145
	}
L138:
	;
	v927 = *(*int32)(unsafe.Add(mBase, uint32(v25)+240))
	v928 = *(*int32)(unsafe.Add(mBase, uint32(v25)+16))
	if base.Ui32(v928) <= base.Ui32(v927) {
		goto L137
	} else {
		goto L139
	}
L139:
	;
	if v862 == int32(0) {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v950 = F_memcpy(m, v25+int32(16), v25+int32(240), int32(224))
	mBase = m.M
	v966 = v924
	goto L135
L141:
	;
	goto L140
L142:
	;
	v935 = v862 + int32(16)
	v936 = *(*int32)(unsafe.Add(mBase, uint32(v935)))
	F_WebPSafeFree(m, v936)
	mBase = m.M
	v938 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	v939 = int32(0)
	base.Simd_g_v128_store(m, v935, v939, v938)
	base.Simd_g_v128_store(m, v862, v939, v938)
	goto L141
L143:
	;
	v966 = v924
	goto L135
L144:
	;
	goto L143
L145:
	;
	v956 = v25 + int32(260)
	v957 = *(*int32)(unsafe.Add(mBase, uint32(v956)))
	F_WebPSafeFree(m, v957)
	mBase = m.M
	v959 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	v960 = int32(0)
	base.Simd_g_v128_store(m, v956, v960, v959)
	base.Simd_g_v128_store(m, v893, v960, v959)
	goto L144
L146:
	;
	goto L134
L147:
	;
	v969 = int32(1)
	if v966 != 0 {
		v898 = int32(base.Ui32(v898) >> (uint(v969) % 32))
		v905 = v905 + v969
		goto L133
	} else {
		goto L148
	}
L148:
	;
	goto L146
L149:
	;
	v985 = v966
	v993 = v862
	goto L45
L150:
	;
	if v993 == int32(0) {
		goto L155
	} else {
		goto L156
	}
L151:
	;
	if v175 == int32(0) {
		goto L152
	} else {
		goto L153
	}
L152:
	;
	v1014 = *(*int32)(unsafe.Add(mBase, uint32(v993+int32(16))))
	v1017 = *(*int32)(unsafe.Add(mBase, uint32(v993+int32(20))))
	v1063 = l0
	v1065 = v25
	v1066 = int32(1)
	v1072 = v1017
	v1073 = v1014
	v1075 = v76
	goto L43
L153:
	;
	v1002 = *(*int32)(unsafe.Add(mBase, uint32(v25)+232))
	*(*int32)(unsafe.Add(mBase, uint32(v175)+180)) = v1002
	v1005 = base.Simd_g_v128_load(m, v25, int32(200))
	base.Simd_g_v128_store(m, v175, int32(148), v1005)
	v1009 = base.Simd_g_v128_load(m, v25, int32(216))
	base.Simd_g_v128_store(m, v175, int32(164), v1009)
	goto L152
L154:
	;
	goto L44
L155:
	;
	goto L154
L156:
	;
	v1024 = v993 + int32(16)
	v1025 = *(*int32)(unsafe.Add(mBase, uint32(v1024)))
	F_WebPSafeFree(m, v1025)
	mBase = m.M
	v1027 = base.Simd_g_const(&F_CompressAlphaJob__k1)
	v1028 = int32(0)
	base.Simd_g_v128_store(m, v1024, v1028, v1027)
	base.Simd_g_v128_store(m, v993, v1028, v1027)
	goto L155
L157:
	;
	v1060 = int32(0)
	v1063 = l0
	v1065 = v25
	v1066 = v1060
	v1072 = v1060
	v1073 = v1060
	v1075 = v76
	goto L43
L158:
	;
	goto L157
L159:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v30)+92)) = int32(1)
	goto L158
L160:
	;
	v1088 = *(*int64)(unsafe.Add(mBase, uint32(v1065)+8))
	*(*int64)(unsafe.Add(mBase, uint32(v1063)+uint32(_c_F_CompressAlphaJob[0]))) = v1088
	v1090 = *(*int32)(unsafe.Add(mBase, uint32(v1085)))
	*(*int32)(unsafe.Add(mBase, uint32(v1085))) = v1090 + v1072
	v1093 = v1063
	v1095 = v1065
	v1096 = v1066
	v1102 = v1072
	v1103 = v1073
	v1105 = v1075
	goto L36
L161:
	;
	v1116 = v1093
	v1118 = v1095
	v1119 = v1096
	v1125 = v1102
	v1126 = v1103
	goto L1
L162:
	;
	m.G0 = v1118 + int32(496)
	return v1142
L163:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1116)+376)) = v1126
	*(*int32)(unsafe.Add(mBase, uint32(v1116)+380)) = v1125
	v1142 = int32(1)
	goto L162
L164:
	;
	v1142 = int32(0)
	goto L162
}

var F_CompressAlphaJob__k0 = [2]uint64{0xfe, 0xfffffff800000000}
var F_CompressAlphaJob__k1 = [2]uint64{0x0, 0x0}
var F_CompressAlphaJob__k2 = [2]uint64{0x7060504, 0x0}
var F_CompressAlphaJob__k3 = [2]uint64{0xb0a0908, 0x0}
var F_CompressAlphaJob__k4 = [2]uint64{0xf0e0d0c, 0x0}
var F_CompressAlphaJob__k5 = [2]uint64{0x100000000, 0x300000002}
var F_CompressAlphaJob__k6 = [2]uint64{0x100000001, 0x100000001}
var F_CompressAlphaJob__k7 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_CompressAlphaJob__k8 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_ConvertARGBToUV_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v50 int32
	_ = v50
	var v51 base.V128
	_ = v51
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v57 base.V128
	_ = v57
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v61 base.V128
	_ = v61
	var v63 base.V128
	_ = v63
	var v69 base.V128
	_ = v69
	var v73 base.V128
	_ = v73
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v81 base.V128
	_ = v81
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v86 base.V128
	_ = v86
	var v89 base.V128
	_ = v89
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v99 base.V128
	_ = v99
	var v100 base.V128
	_ = v100
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v109 base.V128
	_ = v109
	var v110 base.V128
	_ = v110
	var v113 base.V128
	_ = v113
	var v120 base.V128
	_ = v120
	var v122 base.V128
	_ = v122
	var v123 base.V128
	_ = v123
	var v126 base.V128
	_ = v126
	var v128 int32
	_ = v128
	var v130 base.V128
	_ = v130
	var v131 base.V128
	_ = v131
	var v134 base.V128
	_ = v134
	var v144 base.V128
	_ = v144
	var v148 base.V128
	_ = v148
	var v150 base.V128
	_ = v150
	var v152 base.V128
	_ = v152
	var v154 base.V128
	_ = v154
	var v156 base.V128
	_ = v156
	var v162 base.V128
	_ = v162
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v170 base.V128
	_ = v170
	var v172 base.V128
	_ = v172
	var v174 base.V128
	_ = v174
	var v178 base.V128
	_ = v178
	var v185 base.V128
	_ = v185
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v191 base.V128
	_ = v191
	var v198 base.V128
	_ = v198
	var v200 base.V128
	_ = v200
	var v203 base.V128
	_ = v203
	var v210 base.V128
	_ = v210
	var v212 base.V128
	_ = v212
	var v219 base.V128
	_ = v219
	var v222 base.V128
	_ = v222
	var v229 base.V128
	_ = v229
	var v230 base.V128
	_ = v230
	var v232 base.V128
	_ = v232
	var v258 base.V128
	_ = v258
	var v259 int32
	_ = v259
	var v260 base.V128
	_ = v260
	var v263 base.V128
	_ = v263
	var v265 base.V128
	_ = v265
	var v266 base.V128
	_ = v266
	var v267 int32
	_ = v267
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v278 int32
	_ = v278
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v285 int32
	_ = v285
	var v305 int32
	_ = v305
	var v306 int32
	_ = v306
	var v323 int32
	_ = v323
	var v324 int32
	_ = v324
	var v335 int32
	_ = v335
	var v343 base.V128
	_ = v343
	var v344 int32
	_ = v344
	var v345 int32
	_ = v345
	var v346 int32
	_ = v346
	var v357 int32
	_ = v357
	var v358 base.V128
	_ = v358
	var v360 base.V128
	_ = v360
	var v361 int32
	_ = v361
	var v363 int32
	_ = v363
	var v376 int32
	_ = v376
	var v382 base.V128
	_ = v382
	var v385 base.V128
	_ = v385
	var v388 base.V128
	_ = v388
	var v391 base.V128
	_ = v391
	var v392 int32
	_ = v392
	var v394 base.V128
	_ = v394
	var v417 base.V128
	_ = v417
	var v420 base.V128
	_ = v420
	var v423 base.V128
	_ = v423
	var v426 base.V128
	_ = v426
	var v430 base.V128
	_ = v430
	var v433 int32
	_ = v433
	var v439 base.V128
	_ = v439
	var v449 base.V128
	_ = v449
	var v450 base.V128
	_ = v450
	var v453 base.V128
	_ = v453
	var v455 int32
	_ = v455
	var v457 base.V128
	_ = v457
	var v477 int32
	_ = v477
	var v484 int32
	_ = v484
	var v491 int32
	_ = v491
	var v499 base.V128
	_ = v499
	var v500 int32
	_ = v500
	var v501 int32
	_ = v501
	var v502 int32
	_ = v502
	var v513 int32
	_ = v513
	var v514 base.V128
	_ = v514
	var v516 base.V128
	_ = v516
	var v517 int32
	_ = v517
	var v519 int32
	_ = v519
	var v532 int32
	_ = v532
	var v538 base.V128
	_ = v538
	var v541 base.V128
	_ = v541
	var v544 base.V128
	_ = v544
	var v547 base.V128
	_ = v547
	var v548 int32
	_ = v548
	var v550 base.V128
	_ = v550
	var v573 base.V128
	_ = v573
	var v576 base.V128
	_ = v576
	var v579 base.V128
	_ = v579
	var v582 base.V128
	_ = v582
	var v586 base.V128
	_ = v586
	var v589 int32
	_ = v589
	var v595 base.V128
	_ = v595
	var v605 base.V128
	_ = v605
	var v606 base.V128
	_ = v606
	var v609 base.V128
	_ = v609
	var v611 int32
	_ = v611
	var v613 base.V128
	_ = v613
	var v615 int32
	_ = v615
	var v617 int32
	_ = v617
	var v619 int32
	_ = v619
	var v622 base.V128
	_ = v622
	var v624 base.V128
	_ = v624
	var v625 base.V128
	_ = v625
	var v631 base.V128
	_ = v631
	var v655 base.V128
	_ = v655
	var v657 base.V128
	_ = v657
	var v658 base.V128
	_ = v658
	var v676 int32
	_ = v676
	var v684 int32
	_ = v684
	var v708 int32
	_ = v708
	var v711 int32
	_ = v711
	var v723 int32
	_ = v723
	var v726 int32
	_ = v726
	var v727 int32
	_ = v727
	var v729 int32
	_ = v729
	var v731 int32
	_ = v731
	var v736 int32
	_ = v736
	var v739 int32
	_ = v739
	var v747 int32
	_ = v747
	var v751 int32
	_ = v751
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v767 int32
	_ = v767
	var v772 int32
	_ = v772
	var v774 int32
	_ = v774
	var v787 int32
	_ = v787
	var v792 int32
	_ = v792
	var v797 int32
	_ = v797
	var v805 int32
	_ = v805
	var v829 int32
	_ = v829
	var v832 int32
	_ = v832
	var v847 int32
	_ = v847
	var v848 int32
	_ = v848
	var v850 int32
	_ = v850
	var v852 int32
	_ = v852
	var v857 int32
	_ = v857
	var v860 int32
	_ = v860
	var v868 int32
	_ = v868
	var v872 int32
	_ = v872
	var v880 int32
	_ = v880
	var v881 int32
	_ = v881
	var v884 int32
	_ = v884
	var v886 int32
	_ = v886
	var v887 int32
	_ = v887
	var v901 int32
	_ = v901
	var v906 int32
	_ = v906
	var v913 int32
	_ = v913
	var v936 int32
	_ = v936
	var v939 int32
	_ = v939
	var v940 int32
	_ = v940
	var v941 int32
	_ = v941
	var v946 int32
	_ = v946
	var v953 int32
	_ = v953
	var v957 int32
	_ = v957
	var v959 int32
	_ = v959
	var v960 int32
	_ = v960
	var v972 int32
	_ = v972
	var v979 int32
	_ = v979
	var v980 int32
	_ = v980
	var v982 int32
	_ = v982
	var v985 int32
	_ = v985
	var v987 int32
	_ = v987
	var v988 int32
	_ = v988
	var v993 int32
	_ = v993
	v6 = int32(0)
	v25 = l3 & int32(-32)
	if v25 < int32(1) {
		v281 = l1
		v282 = l2
		v285 = v6
	} else {
		v29 = l1
		v30 = l2
		v33 = v6
		v35 = l0
		for {
			v50 = int32(0)
			v51 = base.Simd_g_v128_load_rng(m, v35, v50, int32(0), int32(128))
			v55 = base.Simd_g_v128_load_nc(m, v35+int32(16), v50)
			v56 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k0)
			v57 = base.Simd_g_i8x16_shuffle2(v51, v55, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v58 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k3)
			v59 = base.Simd_g_i8x16_shuffle2(v51, v55, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v61 = base.Simd_g_i8x16_shuffle2(v57, v59, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v63 = base.Simd_g_i8x16_shuffle2(v57, v59, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v69 = base.Simd_g_v128_load_nc(m, v35+int32(32), v50)
			v73 = base.Simd_g_v128_load_nc(m, v35+int32(48), v50)
			v75 = base.Simd_g_i8x16_shuffle2(v69, v73, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v77 = base.Simd_g_i8x16_shuffle2(v69, v73, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v79 = base.Simd_g_i8x16_shuffle2(v75, v77, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v81 = base.Simd_g_i8x16_shuffle2(v75, v77, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v84 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k6)
			v85 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v61, v63, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v79, v81, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k7), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k8))
			v86 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k9)
			v89 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k10)
			v94 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v85, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2)), v89), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v85, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), v89))
			v96 = base.Simd_g_i8x16_shuffle2(v61, v63, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v98 = base.Simd_g_i8x16_shuffle2(v79, v81, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v99 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k11)
			v100 = base.Simd_g_i8x16_shuffle2(v96, v98, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k12), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k13))
			v107 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v100, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2)), v89), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v100, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), v89))
			v108 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k14)
			v109 = base.Simd_g_i8x16_shuffle2(v94, v107, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k15), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k16))
			v110 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k17)
			v113 = base.Simd_g_i8x16_shuffle2(v96, v98, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k7), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k8))
			v120 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v113, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2)), v89), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v113, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), v89))
			v122 = base.Simd_g_i8x16_shuffle2(v107, v120, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k15), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k16))
			v123 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k18)
			v126 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k19)
			v128 = int32(18)
			v130 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k20)
			v131 = base.Simd_g_i8x16_shuffle2(v94, v107, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k21), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k22))
			v134 = base.Simd_g_i8x16_shuffle2(v107, v120, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k21), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k22))
			v144 = base.Simd_g_v128_load_nc(m, v35+int32(64), v50)
			v148 = base.Simd_g_v128_load_nc(m, v35+int32(80), v50)
			v150 = base.Simd_g_i8x16_shuffle2(v144, v148, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v152 = base.Simd_g_i8x16_shuffle2(v144, v148, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v154 = base.Simd_g_i8x16_shuffle2(v150, v152, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v156 = base.Simd_g_i8x16_shuffle2(v150, v152, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v162 = base.Simd_g_v128_load_nc(m, v35+int32(96), v50)
			v166 = base.Simd_g_v128_load_nc(m, v35+int32(112), v50)
			v168 = base.Simd_g_i8x16_shuffle2(v162, v166, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v170 = base.Simd_g_i8x16_shuffle2(v162, v166, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v172 = base.Simd_g_i8x16_shuffle2(v168, v170, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v174 = base.Simd_g_i8x16_shuffle2(v168, v170, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5))
			v178 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v154, v156, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v172, v174, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k7), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k8))
			v185 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v178, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2)), v89), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v178, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), v89))
			v187 = base.Simd_g_i8x16_shuffle2(v154, v156, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v189 = base.Simd_g_i8x16_shuffle2(v172, v174, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2))
			v191 = base.Simd_g_i8x16_shuffle2(v187, v189, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k12), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k13))
			v198 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v191, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2)), v89), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v191, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), v89))
			v200 = base.Simd_g_i8x16_shuffle2(v185, v198, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k15), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k16))
			v203 = base.Simd_g_i8x16_shuffle2(v187, v189, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k7), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k8))
			v210 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v203, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k2)), v89), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v203, v86, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k5)), v89))
			v212 = base.Simd_g_i8x16_shuffle2(v198, v210, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k15), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k16))
			v219 = base.Simd_g_i8x16_shuffle2(v185, v198, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k21), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k22))
			v222 = base.Simd_g_i8x16_shuffle2(v198, v210, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k21), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k22))
			v229 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v109, v110), base.Simd_g_i32x4_dot_i16x8_s(v122, v123)), v126), v128), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v131, v110), base.Simd_g_i32x4_dot_i16x8_s(v134, v123)), v126), v128)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v200, v110), base.Simd_g_i32x4_dot_i16x8_s(v212, v123)), v126), v128), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v219, v110), base.Simd_g_i32x4_dot_i16x8_s(v222, v123)), v126), v128)))
			v230 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k23)
			v232 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k24)
			v258 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v109, v230), base.Simd_g_i32x4_dot_i16x8_s(v122, v232)), v126), v128), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v131, v230), base.Simd_g_i32x4_dot_i16x8_s(v134, v232)), v126), v128)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v200, v230), base.Simd_g_i32x4_dot_i16x8_s(v212, v232)), v126), v128), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v219, v230), base.Simd_g_i32x4_dot_i16x8_s(v222, v232)), v126), v128)))
			if l4 != 0 {
				v265 = v258
				v266 = v229
			} else {
				v259 = int32(0)
				v260 = base.Simd_g_v128_load(m, v30, v259)
				v263 = base.Simd_g_v128_load(m, v29, v259)
				v265 = base.Simd_g_i8x16_avgr_u(v258, v263)
				v266 = base.Simd_g_i8x16_avgr_u(v229, v260)
			}
			v267 = int32(0)
			base.Simd_g_v128_store(m, v30, v267, v266)
			base.Simd_g_v128_store(m, v29, v267, v265)
			v273 = int32(16)
			v274 = v30 + v273
			v276 = v29 + v273
			v278 = v33 + int32(32)
			if v278 < v25 {
				v29 = v276
				v30 = v274
				v33 = v278
				v35 = v35 + int32(128)
				continue
			} else {
				break
			}
			break
		}
		v281 = v276
		v282 = v274
		v285 = v278
	}
	if l3 <= v285 {
	} else {
		v305 = l0 + v285<<(uint(int32(2))%32)
		v306 = l3 - v285
		v323 = int32(1)
		v324 = v306 >> (uint(v323) % 32)
		if v323 <= v324 {
			if l4 == int32(0) {
				if base.Ui32(v324) < base.Ui32(int32(4)) {
					v684 = int32(0)
					v708 = v684
					v711 = v305 + v684<<(uint(int32(3))%32)
					for {
						v723 = v281 + v708
						v726 = *(*int32)(unsafe.Add(mBase, uint32(v711+int32(4))))
						v727 = int32(15)
						v729 = int32(510)
						v731 = *(*int32)(unsafe.Add(mBase, uint32(v711)))
						v736 = int32(base.Ui32(v726)>>(uint(v727)%32))&v729 + int32(base.Ui32(v731)>>(uint(v727)%32))&v729
						v739 = int32(7)
						v747 = int32(base.Ui32(v726)>>(uint(v739)%32))&v729 + int32(base.Ui32(v731)>>(uint(v739)%32))&v729
						v751 = int32(1)
						v759 = v726<<(uint(v751)%32)&v729 + v731<<(uint(v751)%32)&v729
						v760 = int32(_a_F_ConvertARGBToUV_SSE2_0)
						v763 = int32(33685504)
						v765 = int32(18)
						v767 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v723))))
						v772 = int32(base.Ui32(int32(base.Ui32(v736*int32(-9719)+v747*int32(-19081)+v759*v760+v763)>>(uint(v765)%32))+v767+v751) >> (uint(v751) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v723))) = uint8(v772)
						v774 = v282 + v708
						v787 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v774))))
						v792 = int32(base.Ui32(int32(base.Ui32(v736*v760+v747*int32(-24116)+v759*int32(-4684)+v763)>>(uint(v765)%32))+v787+v751) >> (uint(v751) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v774))) = uint8(v792)
						v797 = v708 + v751
						if v324 != v797 {
							v708 = v797
							v711 = v711 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v913 = v324
				} else {
					v491 = v324 & int32(2147483644)
					v499 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k25)
					v500 = v491
					v501 = v282
					v502 = v281
					for {
						v513 = int32(1)
						v514 = base.Simd_g_i32x4_shl(v499, v513)
						v516 = base.Simd_g_v128_or(v514, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k26))
						v517 = int32(3)
						v519 = int32(2)
						v532 = int32(0)
						v538 = base.Simd_g_v128_load32_splat(m, v305+base.Simd_g_i32x4_extract_lane_l0(v516)<<(uint(v519)%32), v532)
						v541 = base.Simd_g_v128_load32_lane_l1(m, v305+base.Simd_g_i32x4_extract_lane_l1(v516)<<(uint(v519)%32), v532, v538)
						v544 = base.Simd_g_v128_load32_lane_l2(m, v305+base.Simd_g_i32x4_extract_lane_l2(v516)<<(uint(v519)%32), v532, v541)
						v547 = base.Simd_g_v128_load32_lane_l3(m, v305+base.Simd_g_i32x4_extract_lane_l3(v516)<<(uint(v519)%32), v532, v544)
						v548 = int32(15)
						v550 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k27)
						v573 = base.Simd_g_v128_load32_splat(m, v305+base.Simd_g_i32x4_extract_lane_l0(v514)<<(uint(v519)%32), v532)
						v576 = base.Simd_g_v128_load32_lane_l1(m, v305+base.Simd_g_i32x4_extract_lane_l1(v514)<<(uint(v519)%32), v532, v573)
						v579 = base.Simd_g_v128_load32_lane_l2(m, v305+base.Simd_g_i32x4_extract_lane_l2(v514)<<(uint(v519)%32), v532, v576)
						v582 = base.Simd_g_v128_load32_lane_l3(m, v305+base.Simd_g_i32x4_extract_lane_l3(v514)<<(uint(v519)%32), v532, v579)
						v586 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v547, v548), v550), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v582, v548), v550))
						v589 = int32(7)
						v595 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v547, v589), v550), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v582, v589), v550))
						v605 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v547, v513), v550), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v582, v513), v550))
						v606 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k17)
						v609 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k19)
						v611 = int32(18)
						v613 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k28)
						v615 = int32(4)
						v617 = int32(5)
						v619 = int32(6)
						v622 = base.Simd_g_i16x8_replace_lane_l7(base.Simd_g_i16x8_replace_lane_l6(base.Simd_g_i16x8_replace_lane_l5(base.Simd_g_i16x8_replace_lane_l4(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v586, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k29)), base.Simd_g_i32x4_mul(v595, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k30))), base.Simd_g_i32x4_mul(v605, v606)), v609), v611), v605, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k28), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k31)), v305), v305), v305), v305)
						v624 = base.Simd_g_v128_load32_zero(m, v502, v532)
						v625 = base.Simd_g_i16x8_extend_low_i8x16_u(v624)
						v631 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k32)
						base.Simd_g_v128_store32_lane_l0(m, v502, v532, base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_sub(base.Simd_g_v128_or(v622, v625), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v622, v625), v513)), v605, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k32), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k31)))
						v655 = base.Simd_g_i16x8_replace_lane_l7(base.Simd_g_i16x8_replace_lane_l6(base.Simd_g_i16x8_replace_lane_l5(base.Simd_g_i16x8_replace_lane_l4(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v586, v606), base.Simd_g_i32x4_mul(v595, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k33))), base.Simd_g_i32x4_mul(v605, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k34))), v609), v611), v605, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k28), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k31)), v305), v305), v305), v305)
						v657 = base.Simd_g_v128_load32_zero(m, v501, v532)
						v658 = base.Simd_g_i16x8_extend_low_i8x16_u(v657)
						base.Simd_g_v128_store32_lane_l0(m, v501, v532, base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_sub(base.Simd_g_v128_or(v655, v658), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v655, v658), v513)), v655, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k32), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k31)))
						v676 = v500 + int32(-4)
						if v676 != 0 {
							v499 = base.Simd_g_i32x4_add(v499, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k35))
							v500 = v676
							v501 = v501 + v615
							v502 = v502 + v615
							continue
						} else {
							break
						}
						break
					}
					if v324 == v491 {
						v913 = v324
					} else {
						v684 = v491
						v708 = v684
						v711 = v305 + v684<<(uint(int32(3))%32)
						for {
							v723 = v281 + v708
							v726 = *(*int32)(unsafe.Add(mBase, uint32(v711+int32(4))))
							v727 = int32(15)
							v729 = int32(510)
							v731 = *(*int32)(unsafe.Add(mBase, uint32(v711)))
							v736 = int32(base.Ui32(v726)>>(uint(v727)%32))&v729 + int32(base.Ui32(v731)>>(uint(v727)%32))&v729
							v739 = int32(7)
							v747 = int32(base.Ui32(v726)>>(uint(v739)%32))&v729 + int32(base.Ui32(v731)>>(uint(v739)%32))&v729
							v751 = int32(1)
							v759 = v726<<(uint(v751)%32)&v729 + v731<<(uint(v751)%32)&v729
							v760 = int32(_a_F_ConvertARGBToUV_SSE2_0)
							v763 = int32(33685504)
							v765 = int32(18)
							v767 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v723))))
							v772 = int32(base.Ui32(int32(base.Ui32(v736*int32(-9719)+v747*int32(-19081)+v759*v760+v763)>>(uint(v765)%32))+v767+v751) >> (uint(v751) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v723))) = uint8(v772)
							v774 = v282 + v708
							v787 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v774))))
							v792 = int32(base.Ui32(int32(base.Ui32(v736*v760+v747*int32(-24116)+v759*int32(-4684)+v763)>>(uint(v765)%32))+v787+v751) >> (uint(v751) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v774))) = uint8(v792)
							v797 = v708 + v751
							if v324 != v797 {
								v708 = v797
								v711 = v711 + int32(8)
								continue
							} else {
								break
							}
							break
						}
						v913 = v324
					}
				}
			} else {
				if base.Ui32(v324) <= base.Ui32(int32(3)) {
					v805 = int32(0)
					v829 = v805
					v832 = v305 + v805<<(uint(int32(3))%32)
					for {
						v847 = *(*int32)(unsafe.Add(mBase, uint32(v832+int32(4))))
						v848 = int32(15)
						v850 = int32(510)
						v852 = *(*int32)(unsafe.Add(mBase, uint32(v832)))
						v857 = int32(base.Ui32(v847)>>(uint(v848)%32))&v850 + int32(base.Ui32(v852)>>(uint(v848)%32))&v850
						v860 = int32(7)
						v868 = int32(base.Ui32(v847)>>(uint(v860)%32))&v850 + int32(base.Ui32(v852)>>(uint(v860)%32))&v850
						v872 = int32(1)
						v880 = v847<<(uint(v872)%32)&v850 + v852<<(uint(v872)%32)&v850
						v881 = int32(_a_F_ConvertARGBToUV_SSE2_0)
						v884 = int32(33685504)
						v886 = int32(18)
						v887 = int32(base.Ui32(v857*int32(67099145)+v868*int32(67089783)+v880*v881+v884) >> (uint(v886) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v281+v829))) = uint8(v887)
						v901 = int32(base.Ui32(v857*v881+v868*int32(67084748)+v880*int32(67104180)+v884) >> (uint(v886) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v282+v829))) = uint8(v901)
						v906 = v829 + v872
						if v324 != v906 {
							v829 = v906
							v832 = v832 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v913 = v324
				} else {
					v335 = v324 & int32(2147483644)
					v343 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k25)
					v344 = v335
					v345 = v282
					v346 = v281
					for {
						v357 = int32(1)
						v358 = base.Simd_g_i32x4_shl(v343, v357)
						v360 = base.Simd_g_v128_or(v358, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k26))
						v361 = int32(3)
						v363 = int32(2)
						v376 = int32(0)
						v382 = base.Simd_g_v128_load32_splat(m, v305+base.Simd_g_i32x4_extract_lane_l0(v360)<<(uint(v363)%32), v376)
						v385 = base.Simd_g_v128_load32_lane_l1(m, v305+base.Simd_g_i32x4_extract_lane_l1(v360)<<(uint(v363)%32), v376, v382)
						v388 = base.Simd_g_v128_load32_lane_l2(m, v305+base.Simd_g_i32x4_extract_lane_l2(v360)<<(uint(v363)%32), v376, v385)
						v391 = base.Simd_g_v128_load32_lane_l3(m, v305+base.Simd_g_i32x4_extract_lane_l3(v360)<<(uint(v363)%32), v376, v388)
						v392 = int32(15)
						v394 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k27)
						v417 = base.Simd_g_v128_load32_splat(m, v305+base.Simd_g_i32x4_extract_lane_l0(v358)<<(uint(v363)%32), v376)
						v420 = base.Simd_g_v128_load32_lane_l1(m, v305+base.Simd_g_i32x4_extract_lane_l1(v358)<<(uint(v363)%32), v376, v417)
						v423 = base.Simd_g_v128_load32_lane_l2(m, v305+base.Simd_g_i32x4_extract_lane_l2(v358)<<(uint(v363)%32), v376, v420)
						v426 = base.Simd_g_v128_load32_lane_l3(m, v305+base.Simd_g_i32x4_extract_lane_l3(v358)<<(uint(v363)%32), v376, v423)
						v430 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v391, v392), v394), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v426, v392), v394))
						v433 = int32(7)
						v439 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v391, v433), v394), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v426, v433), v394))
						v449 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v391, v357), v394), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v426, v357), v394))
						v450 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k17)
						v453 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k19)
						v455 = int32(18)
						v457 = base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k36)
						base.Simd_g_v128_store32_lane_l0(m, v346, v376, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v430, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k37)), base.Simd_g_i32x4_mul(v439, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k38))), base.Simd_g_i32x4_mul(v449, v450)), v453), v455), v449, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k36), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k31)))
						base.Simd_g_v128_store32_lane_l0(m, v345, v376, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v430, v450), base.Simd_g_i32x4_mul(v439, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k39))), base.Simd_g_i32x4_mul(v449, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k40))), v453), v455), v449, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k36), base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k31)))
						v477 = int32(4)
						v484 = v344 + int32(-4)
						if v484 != 0 {
							v343 = base.Simd_g_i32x4_add(v343, base.Simd_g_const(&F_ConvertARGBToUV_SSE2__k35))
							v344 = v484
							v345 = v345 + v477
							v346 = v346 + v477
							continue
						} else {
							break
						}
						break
					}
					if v324 != v335 {
						v805 = v335
						v829 = v805
						v832 = v305 + v805<<(uint(int32(3))%32)
						for {
							v847 = *(*int32)(unsafe.Add(mBase, uint32(v832+int32(4))))
							v848 = int32(15)
							v850 = int32(510)
							v852 = *(*int32)(unsafe.Add(mBase, uint32(v832)))
							v857 = int32(base.Ui32(v847)>>(uint(v848)%32))&v850 + int32(base.Ui32(v852)>>(uint(v848)%32))&v850
							v860 = int32(7)
							v868 = int32(base.Ui32(v847)>>(uint(v860)%32))&v850 + int32(base.Ui32(v852)>>(uint(v860)%32))&v850
							v872 = int32(1)
							v880 = v847<<(uint(v872)%32)&v850 + v852<<(uint(v872)%32)&v850
							v881 = int32(_a_F_ConvertARGBToUV_SSE2_0)
							v884 = int32(33685504)
							v886 = int32(18)
							v887 = int32(base.Ui32(v857*int32(67099145)+v868*int32(67089783)+v880*v881+v884) >> (uint(v886) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v281+v829))) = uint8(v887)
							v901 = int32(base.Ui32(v857*v881+v868*int32(67084748)+v880*int32(67104180)+v884) >> (uint(v886) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v282+v829))) = uint8(v901)
							v906 = v829 + v872
							if v324 != v906 {
								v829 = v906
								v832 = v832 + int32(8)
								continue
							} else {
								break
							}
							break
						}
						v913 = v324
					} else {
						v913 = v324
					}
				}
			}
		} else {
			v913 = int32(0)
		}
		if v306&int32(1) == int32(0) {
		} else {
			v936 = *(*int32)(unsafe.Add(mBase, uint32(v305+v913<<(uint(int32(3))%32))))
			v939 = int32(1020)
			v940 = int32(base.Ui32(v936)>>(uint(int32(14))%32)) & v939
			v941 = int32(_a_F_ConvertARGBToUV_SSE2_0)
			v946 = int32(base.Ui32(v936)>>(uint(int32(6))%32)) & v939
			v953 = v936 << (uint(int32(2)) % 32) & v939
			v957 = int32(33685504)
			v959 = int32(18)
			v960 = int32(base.Ui32(v940*v941+v946*int32(-24116)+v953*int32(-4684)+v957) >> (uint(v959) % 32))
			v972 = int32(base.Ui32(v940*int32(-9719)+v946*int32(-19081)+v953*v941+v957) >> (uint(v959) % 32))
			if l4 == int32(0) {
				v979 = v281 + v913
				v980 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v979))))
				v982 = int32(1)
				v985 = int32(base.Ui32(v972+v980+v982) >> (uint(v982) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v979))) = uint8(v985)
				v987 = v282 + v913
				v988 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v987))))
				v993 = int32(base.Ui32(v960+v988+v982) >> (uint(v982) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v987))) = uint8(v993)
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v282+v913))) = uint8(v960)
				*(*uint8)(unsafe.Add(mBase, uint32(v281+v913))) = uint8(v972)
			}
		}
	}
	return
}

var F_ConvertARGBToUV_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertARGBToUV_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertARGBToUV_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertARGBToUV_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertARGBToUV_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertARGBToUV_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertARGBToUV_SSE2__k6 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertARGBToUV_SSE2__k7 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertARGBToUV_SSE2__k8 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertARGBToUV_SSE2__k9 = [2]uint64{0x0, 0x0}
var F_ConvertARGBToUV_SSE2__k10 = [2]uint64{0x2000200020002, 0x2000200020002}
var F_ConvertARGBToUV_SSE2__k11 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_ConvertARGBToUV_SSE2__k12 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertARGBToUV_SSE2__k13 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertARGBToUV_SSE2__k14 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertARGBToUV_SSE2__k15 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertARGBToUV_SSE2__k16 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertARGBToUV_SSE2__k17 = [2]uint64{0x708000007080, 0x708000007080}
var F_ConvertARGBToUV_SSE2__k18 = [2]uint64{0xedb4a1ccedb4a1cc, 0xedb4a1ccedb4a1cc}
var F_ConvertARGBToUV_SSE2__k19 = [2]uint64{0x202000002020000, 0x202000002020000}
var F_ConvertARGBToUV_SSE2__k20 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertARGBToUV_SSE2__k21 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertARGBToUV_SSE2__k22 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertARGBToUV_SSE2__k23 = [2]uint64{0xb577da09b577da09, 0xb577da09b577da09}
var F_ConvertARGBToUV_SSE2__k24 = [2]uint64{0x7080000070800000, 0x7080000070800000}
var F_ConvertARGBToUV_SSE2__k25 = [2]uint64{0x100000000, 0x300000002}
var F_ConvertARGBToUV_SSE2__k26 = [2]uint64{0x100000001, 0x100000001}
var F_ConvertARGBToUV_SSE2__k27 = [2]uint64{0x1fe000001fe, 0x1fe000001fe}
var F_ConvertARGBToUV_SSE2__k28 = [2]uint64{0xd0c090805040100, 0x100010001000100}
var F_ConvertARGBToUV_SSE2__k29 = [2]uint64{0xffffda09ffffda09, 0xffffda09ffffda09}
var F_ConvertARGBToUV_SSE2__k30 = [2]uint64{0xffffb577ffffb577, 0xffffb577ffffb577}
var F_ConvertARGBToUV_SSE2__k31 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ConvertARGBToUV_SSE2__k32 = [2]uint64{0x6040200, 0x0}
var F_ConvertARGBToUV_SSE2__k33 = [2]uint64{0xffffa1ccffffa1cc, 0xffffa1ccffffa1cc}
var F_ConvertARGBToUV_SSE2__k34 = [2]uint64{0xffffedb4ffffedb4, 0xffffedb4ffffedb4}
var F_ConvertARGBToUV_SSE2__k35 = [2]uint64{0x400000004, 0x400000004}
var F_ConvertARGBToUV_SSE2__k36 = [2]uint64{0xc080400, 0x0}
var F_ConvertARGBToUV_SSE2__k37 = [2]uint64{0x3ffda0903ffda09, 0x3ffda0903ffda09}
var F_ConvertARGBToUV_SSE2__k38 = [2]uint64{0x3ffb57703ffb577, 0x3ffb57703ffb577}
var F_ConvertARGBToUV_SSE2__k39 = [2]uint64{0x3ffa1cc03ffa1cc, 0x3ffa1cc03ffa1cc}
var F_ConvertARGBToUV_SSE2__k40 = [2]uint64{0x3ffedb403ffedb4, 0x3ffedb403ffedb4}

func F_ConvertARGBToUV_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v50 int32
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v57 base.V128
	_ = v57
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v78 base.V128
	_ = v78
	var v80 base.V128
	_ = v80
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v98 base.V128
	_ = v98
	var v99 base.V128
	_ = v99
	var v102 base.V128
	_ = v102
	var v109 base.V128
	_ = v109
	var v111 base.V128
	_ = v111
	var v112 base.V128
	_ = v112
	var v115 base.V128
	_ = v115
	var v117 int32
	_ = v117
	var v119 base.V128
	_ = v119
	var v120 base.V128
	_ = v120
	var v123 base.V128
	_ = v123
	var v133 base.V128
	_ = v133
	var v134 base.V128
	_ = v134
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v145 base.V128
	_ = v145
	var v146 base.V128
	_ = v146
	var v150 base.V128
	_ = v150
	var v151 base.V128
	_ = v151
	var v155 base.V128
	_ = v155
	var v162 base.V128
	_ = v162
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v175 base.V128
	_ = v175
	var v177 base.V128
	_ = v177
	var v180 base.V128
	_ = v180
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v196 base.V128
	_ = v196
	var v199 base.V128
	_ = v199
	var v206 base.V128
	_ = v206
	var v207 base.V128
	_ = v207
	var v209 base.V128
	_ = v209
	var v235 base.V128
	_ = v235
	var v236 int32
	_ = v236
	var v237 base.V128
	_ = v237
	var v240 base.V128
	_ = v240
	var v242 base.V128
	_ = v242
	var v243 base.V128
	_ = v243
	var v244 int32
	_ = v244
	var v250 int32
	_ = v250
	var v251 int32
	_ = v251
	var v253 int32
	_ = v253
	var v255 int32
	_ = v255
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v262 int32
	_ = v262
	var v282 int32
	_ = v282
	var v283 int32
	_ = v283
	var v300 int32
	_ = v300
	var v301 int32
	_ = v301
	var v312 int32
	_ = v312
	var v320 base.V128
	_ = v320
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v334 int32
	_ = v334
	var v335 base.V128
	_ = v335
	var v337 base.V128
	_ = v337
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v353 int32
	_ = v353
	var v359 base.V128
	_ = v359
	var v362 base.V128
	_ = v362
	var v365 base.V128
	_ = v365
	var v368 base.V128
	_ = v368
	var v369 int32
	_ = v369
	var v371 base.V128
	_ = v371
	var v394 base.V128
	_ = v394
	var v397 base.V128
	_ = v397
	var v400 base.V128
	_ = v400
	var v403 base.V128
	_ = v403
	var v407 base.V128
	_ = v407
	var v410 int32
	_ = v410
	var v416 base.V128
	_ = v416
	var v426 base.V128
	_ = v426
	var v427 base.V128
	_ = v427
	var v430 base.V128
	_ = v430
	var v432 int32
	_ = v432
	var v434 base.V128
	_ = v434
	var v454 int32
	_ = v454
	var v461 int32
	_ = v461
	var v468 int32
	_ = v468
	var v476 base.V128
	_ = v476
	var v477 int32
	_ = v477
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v490 int32
	_ = v490
	var v491 base.V128
	_ = v491
	var v493 base.V128
	_ = v493
	var v494 int32
	_ = v494
	var v496 int32
	_ = v496
	var v509 int32
	_ = v509
	var v515 base.V128
	_ = v515
	var v518 base.V128
	_ = v518
	var v521 base.V128
	_ = v521
	var v524 base.V128
	_ = v524
	var v525 int32
	_ = v525
	var v527 base.V128
	_ = v527
	var v550 base.V128
	_ = v550
	var v553 base.V128
	_ = v553
	var v556 base.V128
	_ = v556
	var v559 base.V128
	_ = v559
	var v563 base.V128
	_ = v563
	var v566 int32
	_ = v566
	var v572 base.V128
	_ = v572
	var v582 base.V128
	_ = v582
	var v583 base.V128
	_ = v583
	var v586 base.V128
	_ = v586
	var v588 int32
	_ = v588
	var v590 base.V128
	_ = v590
	var v592 int32
	_ = v592
	var v594 int32
	_ = v594
	var v596 int32
	_ = v596
	var v599 base.V128
	_ = v599
	var v601 base.V128
	_ = v601
	var v602 base.V128
	_ = v602
	var v608 base.V128
	_ = v608
	var v632 base.V128
	_ = v632
	var v634 base.V128
	_ = v634
	var v635 base.V128
	_ = v635
	var v653 int32
	_ = v653
	var v661 int32
	_ = v661
	var v685 int32
	_ = v685
	var v688 int32
	_ = v688
	var v700 int32
	_ = v700
	var v703 int32
	_ = v703
	var v704 int32
	_ = v704
	var v706 int32
	_ = v706
	var v708 int32
	_ = v708
	var v713 int32
	_ = v713
	var v716 int32
	_ = v716
	var v724 int32
	_ = v724
	var v728 int32
	_ = v728
	var v736 int32
	_ = v736
	var v737 int32
	_ = v737
	var v740 int32
	_ = v740
	var v742 int32
	_ = v742
	var v744 int32
	_ = v744
	var v749 int32
	_ = v749
	var v751 int32
	_ = v751
	var v764 int32
	_ = v764
	var v769 int32
	_ = v769
	var v774 int32
	_ = v774
	var v782 int32
	_ = v782
	var v806 int32
	_ = v806
	var v809 int32
	_ = v809
	var v824 int32
	_ = v824
	var v825 int32
	_ = v825
	var v827 int32
	_ = v827
	var v829 int32
	_ = v829
	var v834 int32
	_ = v834
	var v837 int32
	_ = v837
	var v845 int32
	_ = v845
	var v849 int32
	_ = v849
	var v857 int32
	_ = v857
	var v858 int32
	_ = v858
	var v861 int32
	_ = v861
	var v863 int32
	_ = v863
	var v864 int32
	_ = v864
	var v878 int32
	_ = v878
	var v883 int32
	_ = v883
	var v890 int32
	_ = v890
	var v913 int32
	_ = v913
	var v916 int32
	_ = v916
	var v917 int32
	_ = v917
	var v918 int32
	_ = v918
	var v923 int32
	_ = v923
	var v930 int32
	_ = v930
	var v934 int32
	_ = v934
	var v936 int32
	_ = v936
	var v937 int32
	_ = v937
	var v949 int32
	_ = v949
	var v956 int32
	_ = v956
	var v957 int32
	_ = v957
	var v959 int32
	_ = v959
	var v962 int32
	_ = v962
	var v964 int32
	_ = v964
	var v965 int32
	_ = v965
	var v970 int32
	_ = v970
	v6 = int32(0)
	v25 = l3 & int32(-32)
	if v25 < int32(1) {
		v258 = l1
		v259 = l2
		v262 = v6
	} else {
		v29 = l1
		v30 = l2
		v33 = v6
		v35 = l0
		for {
			v50 = int32(0)
			v51 = base.Simd_g_v128_load_rng(m, v35, v50, int32(0), int32(128))
			v52 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k0)
			v53 = base.Simd_g_i8x16_swizzle(v51, v52)
			v57 = base.Simd_g_v128_load_nc(m, v35+int32(16), v50)
			v58 = base.Simd_g_i8x16_swizzle(v57, v52)
			v59 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k1)
			v64 = base.Simd_g_v128_load_nc(m, v35+int32(32), v50)
			v65 = base.Simd_g_i8x16_swizzle(v64, v52)
			v69 = base.Simd_g_v128_load_nc(m, v35+int32(48), v50)
			v70 = base.Simd_g_i8x16_swizzle(v69, v52)
			v73 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k2)
			v74 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v53, v58, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k3), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k4)), base.Simd_g_i8x16_shuffle2(v65, v70, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k3), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k4)), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k5), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k6))
			v75 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k7)
			v76 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k8)
			v78 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k9)
			v80 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k10)
			v83 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v74, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k11), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k12)), v78), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v74, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k13), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k14)), v78))
			v84 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k15)
			v85 = base.Simd_g_i8x16_shuffle2(v53, v58, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k16), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k17))
			v87 = base.Simd_g_i8x16_shuffle2(v65, v70, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k16), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k17))
			v88 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k18)
			v89 = base.Simd_g_i8x16_shuffle2(v85, v87, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k19), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k20))
			v96 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v89, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k11), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k12)), v78), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v89, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k13), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k14)), v78))
			v97 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k21)
			v98 = base.Simd_g_i8x16_shuffle2(v83, v96, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k23))
			v99 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k24)
			v102 = base.Simd_g_i8x16_shuffle2(v85, v87, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k5), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k6))
			v109 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v102, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k11), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k12)), v78), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v102, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k13), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k14)), v78))
			v111 = base.Simd_g_i8x16_shuffle2(v96, v109, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k23))
			v112 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k25)
			v115 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k26)
			v117 = int32(18)
			v119 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k27)
			v120 = base.Simd_g_i8x16_shuffle2(v83, v96, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k28), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k29))
			v123 = base.Simd_g_i8x16_shuffle2(v96, v109, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k28), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k29))
			v133 = base.Simd_g_v128_load_nc(m, v35+int32(64), v50)
			v134 = base.Simd_g_i8x16_swizzle(v133, v52)
			v138 = base.Simd_g_v128_load_nc(m, v35+int32(80), v50)
			v139 = base.Simd_g_i8x16_swizzle(v138, v52)
			v145 = base.Simd_g_v128_load_nc(m, v35+int32(96), v50)
			v146 = base.Simd_g_i8x16_swizzle(v145, v52)
			v150 = base.Simd_g_v128_load_nc(m, v35+int32(112), v50)
			v151 = base.Simd_g_i8x16_swizzle(v150, v52)
			v155 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v134, v139, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k3), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k4)), base.Simd_g_i8x16_shuffle2(v146, v151, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k3), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k4)), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k5), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k6))
			v162 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v155, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k11), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k12)), v78), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v155, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k13), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k14)), v78))
			v164 = base.Simd_g_i8x16_shuffle2(v134, v139, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k16), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k17))
			v166 = base.Simd_g_i8x16_shuffle2(v146, v151, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k16), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k17))
			v168 = base.Simd_g_i8x16_shuffle2(v164, v166, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k19), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k20))
			v175 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v168, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k11), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k12)), v78), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v168, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k13), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k14)), v78))
			v177 = base.Simd_g_i8x16_shuffle2(v162, v175, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k23))
			v180 = base.Simd_g_i8x16_shuffle2(v164, v166, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k5), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k6))
			v187 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v180, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k11), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k12)), v78), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v180, v75, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k13), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k14)), v78))
			v189 = base.Simd_g_i8x16_shuffle2(v175, v187, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k23))
			v196 = base.Simd_g_i8x16_shuffle2(v162, v175, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k28), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k29))
			v199 = base.Simd_g_i8x16_shuffle2(v175, v187, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k28), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k29))
			v206 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v98, v99), base.Simd_g_i32x4_dot_i16x8_s(v111, v112)), v115), v117), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v120, v99), base.Simd_g_i32x4_dot_i16x8_s(v123, v112)), v115), v117)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v177, v99), base.Simd_g_i32x4_dot_i16x8_s(v189, v112)), v115), v117), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v196, v99), base.Simd_g_i32x4_dot_i16x8_s(v199, v112)), v115), v117)))
			v207 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k30)
			v209 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k31)
			v235 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v98, v207), base.Simd_g_i32x4_dot_i16x8_s(v111, v209)), v115), v117), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v120, v207), base.Simd_g_i32x4_dot_i16x8_s(v123, v209)), v115), v117)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v177, v207), base.Simd_g_i32x4_dot_i16x8_s(v189, v209)), v115), v117), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v196, v207), base.Simd_g_i32x4_dot_i16x8_s(v199, v209)), v115), v117)))
			if l4 != 0 {
				v242 = v206
				v243 = v235
			} else {
				v236 = int32(0)
				v237 = base.Simd_g_v128_load(m, v30, v236)
				v240 = base.Simd_g_v128_load(m, v29, v236)
				v242 = base.Simd_g_i8x16_avgr_u(v206, v237)
				v243 = base.Simd_g_i8x16_avgr_u(v235, v240)
			}
			v244 = int32(0)
			base.Simd_g_v128_store(m, v30, v244, v242)
			base.Simd_g_v128_store(m, v29, v244, v243)
			v250 = int32(16)
			v251 = v30 + v250
			v253 = v29 + v250
			v255 = v33 + int32(32)
			if v255 < v25 {
				v29 = v253
				v30 = v251
				v33 = v255
				v35 = v35 + int32(128)
				continue
			} else {
				break
			}
			break
		}
		v258 = v253
		v259 = v251
		v262 = v255
	}
	if l3 <= v262 {
	} else {
		v282 = l0 + v262<<(uint(int32(2))%32)
		v283 = l3 - v262
		v300 = int32(1)
		v301 = v283 >> (uint(v300) % 32)
		if v300 <= v301 {
			if l4 == int32(0) {
				if base.Ui32(v301) < base.Ui32(int32(4)) {
					v661 = int32(0)
					v685 = v661
					v688 = v282 + v661<<(uint(int32(3))%32)
					for {
						v700 = v258 + v685
						v703 = *(*int32)(unsafe.Add(mBase, uint32(v688+int32(4))))
						v704 = int32(15)
						v706 = int32(510)
						v708 = *(*int32)(unsafe.Add(mBase, uint32(v688)))
						v713 = int32(base.Ui32(v703)>>(uint(v704)%32))&v706 + int32(base.Ui32(v708)>>(uint(v704)%32))&v706
						v716 = int32(7)
						v724 = int32(base.Ui32(v703)>>(uint(v716)%32))&v706 + int32(base.Ui32(v708)>>(uint(v716)%32))&v706
						v728 = int32(1)
						v736 = v703<<(uint(v728)%32)&v706 + v708<<(uint(v728)%32)&v706
						v737 = int32(_a_F_ConvertARGBToUV_SSE41_0)
						v740 = int32(33685504)
						v742 = int32(18)
						v744 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v700))))
						v749 = int32(base.Ui32(int32(base.Ui32(v713*int32(-9719)+v724*int32(-19081)+v736*v737+v740)>>(uint(v742)%32))+v744+v728) >> (uint(v728) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v700))) = uint8(v749)
						v751 = v259 + v685
						v764 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v751))))
						v769 = int32(base.Ui32(int32(base.Ui32(v713*v737+v724*int32(-24116)+v736*int32(-4684)+v740)>>(uint(v742)%32))+v764+v728) >> (uint(v728) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v751))) = uint8(v769)
						v774 = v685 + v728
						if v301 != v774 {
							v685 = v774
							v688 = v688 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v890 = v301
				} else {
					v468 = v301 & int32(2147483644)
					v476 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k32)
					v477 = v468
					v478 = v259
					v479 = v258
					for {
						v490 = int32(1)
						v491 = base.Simd_g_i32x4_shl(v476, v490)
						v493 = base.Simd_g_v128_or(v491, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k33))
						v494 = int32(3)
						v496 = int32(2)
						v509 = int32(0)
						v515 = base.Simd_g_v128_load32_splat(m, v282+base.Simd_g_i32x4_extract_lane_l0(v493)<<(uint(v496)%32), v509)
						v518 = base.Simd_g_v128_load32_lane_l1(m, v282+base.Simd_g_i32x4_extract_lane_l1(v493)<<(uint(v496)%32), v509, v515)
						v521 = base.Simd_g_v128_load32_lane_l2(m, v282+base.Simd_g_i32x4_extract_lane_l2(v493)<<(uint(v496)%32), v509, v518)
						v524 = base.Simd_g_v128_load32_lane_l3(m, v282+base.Simd_g_i32x4_extract_lane_l3(v493)<<(uint(v496)%32), v509, v521)
						v525 = int32(15)
						v527 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k34)
						v550 = base.Simd_g_v128_load32_splat(m, v282+base.Simd_g_i32x4_extract_lane_l0(v491)<<(uint(v496)%32), v509)
						v553 = base.Simd_g_v128_load32_lane_l1(m, v282+base.Simd_g_i32x4_extract_lane_l1(v491)<<(uint(v496)%32), v509, v550)
						v556 = base.Simd_g_v128_load32_lane_l2(m, v282+base.Simd_g_i32x4_extract_lane_l2(v491)<<(uint(v496)%32), v509, v553)
						v559 = base.Simd_g_v128_load32_lane_l3(m, v282+base.Simd_g_i32x4_extract_lane_l3(v491)<<(uint(v496)%32), v509, v556)
						v563 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v524, v525), v527), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v559, v525), v527))
						v566 = int32(7)
						v572 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v524, v566), v527), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v559, v566), v527))
						v582 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v524, v490), v527), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v559, v490), v527))
						v583 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k24)
						v586 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k26)
						v588 = int32(18)
						v590 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k35)
						v592 = int32(4)
						v594 = int32(5)
						v596 = int32(6)
						v599 = base.Simd_g_i16x8_replace_lane_l7(base.Simd_g_i16x8_replace_lane_l6(base.Simd_g_i16x8_replace_lane_l5(base.Simd_g_i16x8_replace_lane_l4(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v563, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k36)), base.Simd_g_i32x4_mul(v572, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k37))), base.Simd_g_i32x4_mul(v582, v583)), v586), v588), v582, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k35), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k38)), v282), v282), v282), v282)
						v601 = base.Simd_g_v128_load32_zero(m, v479, v509)
						v602 = base.Simd_g_i16x8_extend_low_i8x16_u(v601)
						v608 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k39)
						base.Simd_g_v128_store32_lane_l0(m, v479, v509, base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_sub(base.Simd_g_v128_or(v599, v602), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v599, v602), v490)), v582, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k39), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k38)))
						v632 = base.Simd_g_i16x8_replace_lane_l7(base.Simd_g_i16x8_replace_lane_l6(base.Simd_g_i16x8_replace_lane_l5(base.Simd_g_i16x8_replace_lane_l4(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v563, v583), base.Simd_g_i32x4_mul(v572, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k40))), base.Simd_g_i32x4_mul(v582, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k41))), v586), v588), v582, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k35), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k38)), v282), v282), v282), v282)
						v634 = base.Simd_g_v128_load32_zero(m, v478, v509)
						v635 = base.Simd_g_i16x8_extend_low_i8x16_u(v634)
						base.Simd_g_v128_store32_lane_l0(m, v478, v509, base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_sub(base.Simd_g_v128_or(v632, v635), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v632, v635), v490)), v632, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k39), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k38)))
						v653 = v477 + int32(-4)
						if v653 != 0 {
							v476 = base.Simd_g_i32x4_add(v476, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k42))
							v477 = v653
							v478 = v478 + v592
							v479 = v479 + v592
							continue
						} else {
							break
						}
						break
					}
					if v301 == v468 {
						v890 = v301
					} else {
						v661 = v468
						v685 = v661
						v688 = v282 + v661<<(uint(int32(3))%32)
						for {
							v700 = v258 + v685
							v703 = *(*int32)(unsafe.Add(mBase, uint32(v688+int32(4))))
							v704 = int32(15)
							v706 = int32(510)
							v708 = *(*int32)(unsafe.Add(mBase, uint32(v688)))
							v713 = int32(base.Ui32(v703)>>(uint(v704)%32))&v706 + int32(base.Ui32(v708)>>(uint(v704)%32))&v706
							v716 = int32(7)
							v724 = int32(base.Ui32(v703)>>(uint(v716)%32))&v706 + int32(base.Ui32(v708)>>(uint(v716)%32))&v706
							v728 = int32(1)
							v736 = v703<<(uint(v728)%32)&v706 + v708<<(uint(v728)%32)&v706
							v737 = int32(_a_F_ConvertARGBToUV_SSE41_0)
							v740 = int32(33685504)
							v742 = int32(18)
							v744 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v700))))
							v749 = int32(base.Ui32(int32(base.Ui32(v713*int32(-9719)+v724*int32(-19081)+v736*v737+v740)>>(uint(v742)%32))+v744+v728) >> (uint(v728) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v700))) = uint8(v749)
							v751 = v259 + v685
							v764 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v751))))
							v769 = int32(base.Ui32(int32(base.Ui32(v713*v737+v724*int32(-24116)+v736*int32(-4684)+v740)>>(uint(v742)%32))+v764+v728) >> (uint(v728) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v751))) = uint8(v769)
							v774 = v685 + v728
							if v301 != v774 {
								v685 = v774
								v688 = v688 + int32(8)
								continue
							} else {
								break
							}
							break
						}
						v890 = v301
					}
				}
			} else {
				if base.Ui32(v301) <= base.Ui32(int32(3)) {
					v782 = int32(0)
					v806 = v782
					v809 = v282 + v782<<(uint(int32(3))%32)
					for {
						v824 = *(*int32)(unsafe.Add(mBase, uint32(v809+int32(4))))
						v825 = int32(15)
						v827 = int32(510)
						v829 = *(*int32)(unsafe.Add(mBase, uint32(v809)))
						v834 = int32(base.Ui32(v824)>>(uint(v825)%32))&v827 + int32(base.Ui32(v829)>>(uint(v825)%32))&v827
						v837 = int32(7)
						v845 = int32(base.Ui32(v824)>>(uint(v837)%32))&v827 + int32(base.Ui32(v829)>>(uint(v837)%32))&v827
						v849 = int32(1)
						v857 = v824<<(uint(v849)%32)&v827 + v829<<(uint(v849)%32)&v827
						v858 = int32(_a_F_ConvertARGBToUV_SSE41_0)
						v861 = int32(33685504)
						v863 = int32(18)
						v864 = int32(base.Ui32(v834*int32(67099145)+v845*int32(67089783)+v857*v858+v861) >> (uint(v863) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v258+v806))) = uint8(v864)
						v878 = int32(base.Ui32(v834*v858+v845*int32(67084748)+v857*int32(67104180)+v861) >> (uint(v863) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v259+v806))) = uint8(v878)
						v883 = v806 + v849
						if v301 != v883 {
							v806 = v883
							v809 = v809 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v890 = v301
				} else {
					v312 = v301 & int32(2147483644)
					v320 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k32)
					v321 = v312
					v322 = v259
					v323 = v258
					for {
						v334 = int32(1)
						v335 = base.Simd_g_i32x4_shl(v320, v334)
						v337 = base.Simd_g_v128_or(v335, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k33))
						v338 = int32(3)
						v340 = int32(2)
						v353 = int32(0)
						v359 = base.Simd_g_v128_load32_splat(m, v282+base.Simd_g_i32x4_extract_lane_l0(v337)<<(uint(v340)%32), v353)
						v362 = base.Simd_g_v128_load32_lane_l1(m, v282+base.Simd_g_i32x4_extract_lane_l1(v337)<<(uint(v340)%32), v353, v359)
						v365 = base.Simd_g_v128_load32_lane_l2(m, v282+base.Simd_g_i32x4_extract_lane_l2(v337)<<(uint(v340)%32), v353, v362)
						v368 = base.Simd_g_v128_load32_lane_l3(m, v282+base.Simd_g_i32x4_extract_lane_l3(v337)<<(uint(v340)%32), v353, v365)
						v369 = int32(15)
						v371 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k34)
						v394 = base.Simd_g_v128_load32_splat(m, v282+base.Simd_g_i32x4_extract_lane_l0(v335)<<(uint(v340)%32), v353)
						v397 = base.Simd_g_v128_load32_lane_l1(m, v282+base.Simd_g_i32x4_extract_lane_l1(v335)<<(uint(v340)%32), v353, v394)
						v400 = base.Simd_g_v128_load32_lane_l2(m, v282+base.Simd_g_i32x4_extract_lane_l2(v335)<<(uint(v340)%32), v353, v397)
						v403 = base.Simd_g_v128_load32_lane_l3(m, v282+base.Simd_g_i32x4_extract_lane_l3(v335)<<(uint(v340)%32), v353, v400)
						v407 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v368, v369), v371), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v403, v369), v371))
						v410 = int32(7)
						v416 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v368, v410), v371), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v403, v410), v371))
						v426 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v368, v334), v371), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v403, v334), v371))
						v427 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k24)
						v430 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k26)
						v432 = int32(18)
						v434 = base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k43)
						base.Simd_g_v128_store32_lane_l0(m, v323, v353, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v407, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k44)), base.Simd_g_i32x4_mul(v416, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k45))), base.Simd_g_i32x4_mul(v426, v427)), v430), v432), v426, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k43), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k38)))
						base.Simd_g_v128_store32_lane_l0(m, v322, v353, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v407, v427), base.Simd_g_i32x4_mul(v416, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k46))), base.Simd_g_i32x4_mul(v426, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k47))), v430), v432), v426, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k43), base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k38)))
						v454 = int32(4)
						v461 = v321 + int32(-4)
						if v461 != 0 {
							v320 = base.Simd_g_i32x4_add(v320, base.Simd_g_const(&F_ConvertARGBToUV_SSE41__k42))
							v321 = v461
							v322 = v322 + v454
							v323 = v323 + v454
							continue
						} else {
							break
						}
						break
					}
					if v301 != v312 {
						v782 = v312
						v806 = v782
						v809 = v282 + v782<<(uint(int32(3))%32)
						for {
							v824 = *(*int32)(unsafe.Add(mBase, uint32(v809+int32(4))))
							v825 = int32(15)
							v827 = int32(510)
							v829 = *(*int32)(unsafe.Add(mBase, uint32(v809)))
							v834 = int32(base.Ui32(v824)>>(uint(v825)%32))&v827 + int32(base.Ui32(v829)>>(uint(v825)%32))&v827
							v837 = int32(7)
							v845 = int32(base.Ui32(v824)>>(uint(v837)%32))&v827 + int32(base.Ui32(v829)>>(uint(v837)%32))&v827
							v849 = int32(1)
							v857 = v824<<(uint(v849)%32)&v827 + v829<<(uint(v849)%32)&v827
							v858 = int32(_a_F_ConvertARGBToUV_SSE41_0)
							v861 = int32(33685504)
							v863 = int32(18)
							v864 = int32(base.Ui32(v834*int32(67099145)+v845*int32(67089783)+v857*v858+v861) >> (uint(v863) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v258+v806))) = uint8(v864)
							v878 = int32(base.Ui32(v834*v858+v845*int32(67084748)+v857*int32(67104180)+v861) >> (uint(v863) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v259+v806))) = uint8(v878)
							v883 = v806 + v849
							if v301 != v883 {
								v806 = v883
								v809 = v809 + int32(8)
								continue
							} else {
								break
							}
							break
						}
						v890 = v301
					} else {
						v890 = v301
					}
				}
			}
		} else {
			v890 = int32(0)
		}
		if v283&int32(1) == int32(0) {
		} else {
			v913 = *(*int32)(unsafe.Add(mBase, uint32(v282+v890<<(uint(int32(3))%32))))
			v916 = int32(1020)
			v917 = int32(base.Ui32(v913)>>(uint(int32(14))%32)) & v916
			v918 = int32(_a_F_ConvertARGBToUV_SSE41_0)
			v923 = int32(base.Ui32(v913)>>(uint(int32(6))%32)) & v916
			v930 = v913 << (uint(int32(2)) % 32) & v916
			v934 = int32(33685504)
			v936 = int32(18)
			v937 = int32(base.Ui32(v917*v918+v923*int32(-24116)+v930*int32(-4684)+v934) >> (uint(v936) % 32))
			v949 = int32(base.Ui32(v917*int32(-9719)+v923*int32(-19081)+v930*v918+v934) >> (uint(v936) % 32))
			if l4 == int32(0) {
				v956 = v258 + v890
				v957 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v956))))
				v959 = int32(1)
				v962 = int32(base.Ui32(v949+v957+v959) >> (uint(v959) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v956))) = uint8(v962)
				v964 = v259 + v890
				v965 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v964))))
				v970 = int32(base.Ui32(v937+v965+v959) >> (uint(v959) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v964))) = uint8(v970)
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v259+v890))) = uint8(v937)
				*(*uint8)(unsafe.Add(mBase, uint32(v258+v890))) = uint8(v949)
			}
		}
	}
	return
}

var F_ConvertARGBToUV_SSE41__k0 = [2]uint64{0xd0905010c080400, 0xf0b07030e0a0602}
var F_ConvertARGBToUV_SSE41__k1 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_ConvertARGBToUV_SSE41__k2 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertARGBToUV_SSE41__k3 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_ConvertARGBToUV_SSE41__k4 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_ConvertARGBToUV_SSE41__k5 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertARGBToUV_SSE41__k6 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertARGBToUV_SSE41__k7 = [2]uint64{0x0, 0x0}
var F_ConvertARGBToUV_SSE41__k8 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertARGBToUV_SSE41__k9 = [2]uint64{0x2000200020002, 0x2000200020002}
var F_ConvertARGBToUV_SSE41__k10 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertARGBToUV_SSE41__k11 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertARGBToUV_SSE41__k12 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertARGBToUV_SSE41__k13 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertARGBToUV_SSE41__k14 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertARGBToUV_SSE41__k15 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_ConvertARGBToUV_SSE41__k16 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_ConvertARGBToUV_SSE41__k17 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_ConvertARGBToUV_SSE41__k18 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_ConvertARGBToUV_SSE41__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertARGBToUV_SSE41__k20 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertARGBToUV_SSE41__k21 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertARGBToUV_SSE41__k22 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertARGBToUV_SSE41__k23 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertARGBToUV_SSE41__k24 = [2]uint64{0x708000007080, 0x708000007080}
var F_ConvertARGBToUV_SSE41__k25 = [2]uint64{0xedb4a1ccedb4a1cc, 0xedb4a1ccedb4a1cc}
var F_ConvertARGBToUV_SSE41__k26 = [2]uint64{0x202000002020000, 0x202000002020000}
var F_ConvertARGBToUV_SSE41__k27 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertARGBToUV_SSE41__k28 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertARGBToUV_SSE41__k29 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertARGBToUV_SSE41__k30 = [2]uint64{0xb577da09b577da09, 0xb577da09b577da09}
var F_ConvertARGBToUV_SSE41__k31 = [2]uint64{0x7080000070800000, 0x7080000070800000}
var F_ConvertARGBToUV_SSE41__k32 = [2]uint64{0x100000000, 0x300000002}
var F_ConvertARGBToUV_SSE41__k33 = [2]uint64{0x100000001, 0x100000001}
var F_ConvertARGBToUV_SSE41__k34 = [2]uint64{0x1fe000001fe, 0x1fe000001fe}
var F_ConvertARGBToUV_SSE41__k35 = [2]uint64{0xd0c090805040100, 0x100010001000100}
var F_ConvertARGBToUV_SSE41__k36 = [2]uint64{0xffffda09ffffda09, 0xffffda09ffffda09}
var F_ConvertARGBToUV_SSE41__k37 = [2]uint64{0xffffb577ffffb577, 0xffffb577ffffb577}
var F_ConvertARGBToUV_SSE41__k38 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ConvertARGBToUV_SSE41__k39 = [2]uint64{0x6040200, 0x0}
var F_ConvertARGBToUV_SSE41__k40 = [2]uint64{0xffffa1ccffffa1cc, 0xffffa1ccffffa1cc}
var F_ConvertARGBToUV_SSE41__k41 = [2]uint64{0xffffedb4ffffedb4, 0xffffedb4ffffedb4}
var F_ConvertARGBToUV_SSE41__k42 = [2]uint64{0x400000004, 0x400000004}
var F_ConvertARGBToUV_SSE41__k43 = [2]uint64{0xc080400, 0x0}
var F_ConvertARGBToUV_SSE41__k44 = [2]uint64{0x3ffda0903ffda09, 0x3ffda0903ffda09}
var F_ConvertARGBToUV_SSE41__k45 = [2]uint64{0x3ffb57703ffb577, 0x3ffb57703ffb577}
var F_ConvertARGBToUV_SSE41__k46 = [2]uint64{0x3ffa1cc03ffa1cc, 0x3ffa1cc03ffa1cc}
var F_ConvertARGBToUV_SSE41__k47 = [2]uint64{0x3ffedb403ffedb4, 0x3ffedb403ffedb4}

func F_ConvertARGBToY_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v28 base.V128
	_ = v28
	var v32 int32
	_ = v32
	var v58 int32
	_ = v58
	var v63 int32
	_ = v63
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v88 int32
	_ = v88
	var v105 int32
	_ = v105
	var v112 int32
	_ = v112
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v63 = int32(0)
			v76 = l2 - v63
			v79 = l1 + v63
			v80 = l0 + v63<<(uint(int32(2))%32)
			for {
				v83 = *(*int32)(unsafe.Add(mBase, uint32(v80)))
				v84 = int32(255)
				v88 = int32(16)
				v105 = int32(base.Ui32(v83&v84*int32(_a_F_ConvertARGBToY_C_0)+int32(base.Ui32(v83)>>(uint(v88)%32))&v84*int32(_a_F_ConvertARGBToY_C_1)+int32(base.Ui32(v83)>>(uint(int32(8))%32))&v84*int32(_a_F_ConvertARGBToY_C_2)+int32(1081344)) >> (uint(v88) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v79))) = uint8(v105)
				v112 = v76 + int32(-1)
				if v112 != 0 {
					v76 = v112
					v79 = v79 + int32(1)
					v80 = v80 + int32(4)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v16 = l2 & int32(2147483644)
			v21 = v16
			v22 = l1
			v23 = l0
			for {
				v26 = int32(0)
				v27 = base.Simd_g_v128_load(m, v23, v26)
				v28 = base.Simd_g_const(&F_ConvertARGBToY_C__k0)
				v32 = int32(16)
				base.Simd_g_v128_store32_lane_l0(m, v22, v26, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_v128_and(v27, v28), base.Simd_g_const(&F_ConvertARGBToY_C__k1)), base.Simd_g_i32x4_mul(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v27, v32), v28), base.Simd_g_const(&F_ConvertARGBToY_C__k2))), base.Simd_g_i32x4_mul(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v27, int32(8)), v28), base.Simd_g_const(&F_ConvertARGBToY_C__k3))), base.Simd_g_const(&F_ConvertARGBToY_C__k4)), v32), v27, base.Simd_g_const(&F_ConvertARGBToY_C__k5), base.Simd_g_const(&F_ConvertARGBToY_C__k6)))
				v58 = v21 + int32(-4)
				if v58 != 0 {
					v21 = v58
					v22 = v22 + int32(4)
					v23 = v23 + v32
					continue
				} else {
					break
				}
				break
			}
			if v16 == l2 {
			} else {
				v63 = v16
				v76 = l2 - v63
				v79 = l1 + v63
				v80 = l0 + v63<<(uint(int32(2))%32)
				for {
					v83 = *(*int32)(unsafe.Add(mBase, uint32(v80)))
					v84 = int32(255)
					v88 = int32(16)
					v105 = int32(base.Ui32(v83&v84*int32(_a_F_ConvertARGBToY_C_0)+int32(base.Ui32(v83)>>(uint(v88)%32))&v84*int32(_a_F_ConvertARGBToY_C_1)+int32(base.Ui32(v83)>>(uint(int32(8))%32))&v84*int32(_a_F_ConvertARGBToY_C_2)+int32(1081344)) >> (uint(v88) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v79))) = uint8(v105)
					v112 = v76 + int32(-1)
					if v112 != 0 {
						v76 = v112
						v79 = v79 + int32(1)
						v80 = v80 + int32(4)
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

var F_ConvertARGBToY_C__k0 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ConvertARGBToY_C__k1 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertARGBToY_C__k2 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertARGBToY_C__k3 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertARGBToY_C__k4 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertARGBToY_C__k5 = [2]uint64{0xc080400, 0x0}
var F_ConvertARGBToY_C__k6 = [2]uint64{0x8080808080808080, 0x8080808080808080}

func F_ConvertARGBToY_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 int32
	_ = v47
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v64 base.V128
	_ = v64
	var v68 base.V128
	_ = v68
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v76 base.V128
	_ = v76
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v85 base.V128
	_ = v85
	var v87 base.V128
	_ = v87
	var v89 base.V128
	_ = v89
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v102 base.V128
	_ = v102
	var v105 base.V128
	_ = v105
	var v109 base.V128
	_ = v109
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v128 base.V128
	_ = v128
	var v147 base.V128
	_ = v147
	var v153 int32
	_ = v153
	var v158 int32
	_ = v158
	var v175 int32
	_ = v175
	var v183 int32
	_ = v183
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 base.V128
	_ = v205
	var v206 base.V128
	_ = v206
	var v210 int32
	_ = v210
	var v236 int32
	_ = v236
	var v241 int32
	_ = v241
	var v263 int32
	_ = v263
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v286 int32
	_ = v286
	var v303 int32
	_ = v303
	var v310 int32
	_ = v310
	v4 = int32(0)
	v22 = l2 & int32(-16)
	if v22 < int32(1) {
		v158 = v4
	} else {
		v28 = v4
		v30 = l0
		for {
			v45 = int32(0)
			v46 = base.Simd_g_v128_load_rng(m, v30, v45, int32(0), int32(64))
			v47 = int32(16)
			v50 = base.Simd_g_v128_load_nc(m, v30+v47, v45)
			v51 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k0)
			v52 = base.Simd_g_i8x16_shuffle2(v46, v50, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v53 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k3)
			v54 = base.Simd_g_i8x16_shuffle2(v46, v50, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v56 = base.Simd_g_i8x16_shuffle2(v52, v54, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v58 = base.Simd_g_i8x16_shuffle2(v52, v54, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v64 = base.Simd_g_v128_load_nc(m, v30+int32(32), v45)
			v68 = base.Simd_g_v128_load_nc(m, v30+int32(48), v45)
			v70 = base.Simd_g_i8x16_shuffle2(v64, v68, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v72 = base.Simd_g_i8x16_shuffle2(v64, v68, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v74 = base.Simd_g_i8x16_shuffle2(v70, v72, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v76 = base.Simd_g_i8x16_shuffle2(v70, v72, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v79 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k6)
			v80 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v56, v58, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v74, v76, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5)), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k7), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k8))
			v81 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k9)
			v83 = base.Simd_g_i8x16_shuffle2(v80, v81, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v85 = base.Simd_g_i8x16_shuffle2(v56, v58, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v87 = base.Simd_g_i8x16_shuffle2(v74, v76, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v89 = base.Simd_g_i8x16_shuffle2(v85, v87, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k10), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k11))
			v91 = base.Simd_g_i8x16_shuffle2(v89, v81, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v92 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k12)
			v94 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k13)
			v97 = base.Simd_g_i8x16_shuffle2(v85, v87, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k7), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k8))
			v99 = base.Simd_g_i8x16_shuffle2(v97, v81, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k1), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k2))
			v102 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k14)
			v105 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k15)
			v109 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k16)
			v121 = base.Simd_g_i8x16_shuffle2(v80, v81, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v123 = base.Simd_g_i8x16_shuffle2(v89, v81, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v128 = base.Simd_g_i8x16_shuffle2(v97, v81, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k4), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k5))
			v147 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v83, v91, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k17), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k18)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v91, v99, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k17), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k18)), v102)), v105), v47), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v83, v91, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k19), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k20)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v91, v99, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k19), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k20)), v102)), v105), v47)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v121, v123, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k17), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k18)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v123, v128, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k17), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k18)), v102)), v105), v47), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v121, v123, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k19), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k20)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v123, v128, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k19), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k20)), v102)), v105), v47)))
			base.Simd_g_v128_store(m, l1+v28, v45, v147)
			v153 = v28 + v47
			if v153 < v22 {
				v28 = v153
				v30 = v30 + int32(64)
				continue
			} else {
				break
			}
			break
		}
		v158 = v153
	}
	if l2 <= v158 {
	} else {
		v175 = l2 - v158
		if base.Ui32(v175) <= base.Ui32(int32(3)) {
			v241 = v158
			v263 = l2 - v241
			v265 = l0 + v241<<(uint(int32(2))%32)
			v267 = l1 + v241
			for {
				v281 = *(*int32)(unsafe.Add(mBase, uint32(v265)))
				v282 = int32(255)
				v286 = int32(16)
				v303 = int32(base.Ui32(v281&v282*int32(_a_F_ConvertARGBToY_SSE2_0)+int32(base.Ui32(v281)>>(uint(v286)%32))&v282*int32(_a_F_ConvertARGBToY_SSE2_1)+int32(base.Ui32(v281)>>(uint(int32(8))%32))&v282*int32(_a_F_ConvertARGBToY_SSE2_2)+int32(1081344)) >> (uint(v286) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v267))) = uint8(v303)
				v310 = v263 + int32(-1)
				if v310 != 0 {
					v263 = v310
					v265 = v265 + int32(4)
					v267 = v267 + int32(1)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v183 = v175 & int32(-4)
			v189 = l0 + v158<<(uint(int32(2))%32)
			v190 = l1 + v158
			v203 = v183
			for {
				v204 = int32(0)
				v205 = base.Simd_g_v128_load(m, v189, v204)
				v206 = base.Simd_g_const(&F_ConvertARGBToY_SSE2__k21)
				v210 = int32(16)
				base.Simd_g_v128_store32_lane_l0(m, v190, v204, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_v128_and(v205, v206), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k22)), base.Simd_g_i32x4_mul(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v205, v210), v206), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k23))), base.Simd_g_i32x4_mul(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v205, int32(8)), v206), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k24))), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k15)), v210), v205, base.Simd_g_const(&F_ConvertARGBToY_SSE2__k25), base.Simd_g_const(&F_ConvertARGBToY_SSE2__k26)))
				v236 = v203 + int32(-4)
				if v236 != 0 {
					v189 = v189 + v210
					v190 = v190 + int32(4)
					v203 = v236
					continue
				} else {
					break
				}
				break
			}
			if v175 == v183 {
			} else {
				v241 = v158 + v183
				v263 = l2 - v241
				v265 = l0 + v241<<(uint(int32(2))%32)
				v267 = l1 + v241
				for {
					v281 = *(*int32)(unsafe.Add(mBase, uint32(v265)))
					v282 = int32(255)
					v286 = int32(16)
					v303 = int32(base.Ui32(v281&v282*int32(_a_F_ConvertARGBToY_SSE2_0)+int32(base.Ui32(v281)>>(uint(v286)%32))&v282*int32(_a_F_ConvertARGBToY_SSE2_1)+int32(base.Ui32(v281)>>(uint(int32(8))%32))&v282*int32(_a_F_ConvertARGBToY_SSE2_2)+int32(1081344)) >> (uint(v286) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v267))) = uint8(v303)
					v310 = v263 + int32(-1)
					if v310 != 0 {
						v263 = v310
						v265 = v265 + int32(4)
						v267 = v267 + int32(1)
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

var F_ConvertARGBToY_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertARGBToY_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertARGBToY_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertARGBToY_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertARGBToY_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertARGBToY_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertARGBToY_SSE2__k6 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertARGBToY_SSE2__k7 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertARGBToY_SSE2__k8 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertARGBToY_SSE2__k9 = [2]uint64{0x0, 0x0}
var F_ConvertARGBToY_SSE2__k10 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertARGBToY_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertARGBToY_SSE2__k12 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertARGBToY_SSE2__k13 = [2]uint64{0x412341c7412341c7, 0x412341c7412341c7}
var F_ConvertARGBToY_SSE2__k14 = [2]uint64{0x1914400019144000, 0x1914400019144000}
var F_ConvertARGBToY_SSE2__k15 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertARGBToY_SSE2__k16 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertARGBToY_SSE2__k17 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertARGBToY_SSE2__k18 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertARGBToY_SSE2__k19 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertARGBToY_SSE2__k20 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertARGBToY_SSE2__k21 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ConvertARGBToY_SSE2__k22 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertARGBToY_SSE2__k23 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertARGBToY_SSE2__k24 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertARGBToY_SSE2__k25 = [2]uint64{0xc080400, 0x0}
var F_ConvertARGBToY_SSE2__k26 = [2]uint64{0x8080808080808080, 0x8080808080808080}

func F_ConvertARGBToY_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v49 int32
	_ = v49
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v59 base.V128
	_ = v59
	var v60 base.V128
	_ = v60
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
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
	var v76 base.V128
	_ = v76
	var v78 base.V128
	_ = v78
	var v80 base.V128
	_ = v80
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v86 base.V128
	_ = v86
	var v88 base.V128
	_ = v88
	var v91 base.V128
	_ = v91
	var v94 base.V128
	_ = v94
	var v98 base.V128
	_ = v98
	var v109 base.V128
	_ = v109
	var v110 base.V128
	_ = v110
	var v112 base.V128
	_ = v112
	var v117 base.V128
	_ = v117
	var v136 base.V128
	_ = v136
	var v142 int32
	_ = v142
	var v147 int32
	_ = v147
	var v164 int32
	_ = v164
	var v172 int32
	_ = v172
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v194 base.V128
	_ = v194
	var v195 base.V128
	_ = v195
	var v199 int32
	_ = v199
	var v225 int32
	_ = v225
	var v230 int32
	_ = v230
	var v252 int32
	_ = v252
	var v254 int32
	_ = v254
	var v256 int32
	_ = v256
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v275 int32
	_ = v275
	var v292 int32
	_ = v292
	var v299 int32
	_ = v299
	v4 = int32(0)
	v22 = l2 & int32(-16)
	if v22 < int32(1) {
		v147 = v4
	} else {
		v28 = v4
		v30 = l0
		for {
			v45 = int32(0)
			v46 = base.Simd_g_v128_load_rng(m, v30, v45, int32(0), int32(64))
			v47 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k0)
			v48 = base.Simd_g_i8x16_swizzle(v46, v47)
			v49 = int32(16)
			v52 = base.Simd_g_v128_load_nc(m, v30+v49, v45)
			v53 = base.Simd_g_i8x16_swizzle(v52, v47)
			v54 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k1)
			v59 = base.Simd_g_v128_load_nc(m, v30+int32(32), v45)
			v60 = base.Simd_g_i8x16_swizzle(v59, v47)
			v64 = base.Simd_g_v128_load_nc(m, v30+int32(48), v45)
			v65 = base.Simd_g_i8x16_swizzle(v64, v47)
			v68 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k2)
			v69 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v48, v53, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k3), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k4)), base.Simd_g_i8x16_shuffle2(v60, v65, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k3), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k4)), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k5), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k6))
			v70 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k7)
			v71 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k8)
			v72 = base.Simd_g_i8x16_shuffle2(v69, v70, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k9), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k10))
			v73 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k11)
			v74 = base.Simd_g_i8x16_shuffle2(v48, v53, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k12), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k13))
			v76 = base.Simd_g_i8x16_shuffle2(v60, v65, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k12), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k13))
			v78 = base.Simd_g_i8x16_shuffle2(v74, v76, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k14), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k15))
			v80 = base.Simd_g_i8x16_shuffle2(v78, v70, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k9), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k10))
			v81 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k16)
			v83 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k17)
			v86 = base.Simd_g_i8x16_shuffle2(v74, v76, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k5), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k6))
			v88 = base.Simd_g_i8x16_shuffle2(v86, v70, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k9), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k10))
			v91 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k18)
			v94 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k19)
			v98 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k20)
			v109 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k21)
			v110 = base.Simd_g_i8x16_shuffle2(v69, v70, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k23))
			v112 = base.Simd_g_i8x16_shuffle2(v78, v70, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k23))
			v117 = base.Simd_g_i8x16_shuffle2(v86, v70, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k22), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k23))
			v136 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v72, v80, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k24), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k25)), v83), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v80, v88, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k24), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k25)), v91)), v94), v49), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v72, v80, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k26), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k27)), v83), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v80, v88, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k26), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k27)), v91)), v94), v49)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v110, v112, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k24), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k25)), v83), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v112, v117, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k24), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k25)), v91)), v94), v49), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v110, v112, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k26), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k27)), v83), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v112, v117, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k26), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k27)), v91)), v94), v49)))
			base.Simd_g_v128_store(m, l1+v28, v45, v136)
			v142 = v28 + v49
			if v142 < v22 {
				v28 = v142
				v30 = v30 + int32(64)
				continue
			} else {
				break
			}
			break
		}
		v147 = v142
	}
	if l2 <= v147 {
	} else {
		v164 = l2 - v147
		if base.Ui32(v164) <= base.Ui32(int32(3)) {
			v230 = v147
			v252 = l2 - v230
			v254 = l0 + v230<<(uint(int32(2))%32)
			v256 = l1 + v230
			for {
				v270 = *(*int32)(unsafe.Add(mBase, uint32(v254)))
				v271 = int32(255)
				v275 = int32(16)
				v292 = int32(base.Ui32(v270&v271*int32(_a_F_ConvertARGBToY_SSE41_0)+int32(base.Ui32(v270)>>(uint(v275)%32))&v271*int32(_a_F_ConvertARGBToY_SSE41_1)+int32(base.Ui32(v270)>>(uint(int32(8))%32))&v271*int32(_a_F_ConvertARGBToY_SSE41_2)+int32(1081344)) >> (uint(v275) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v256))) = uint8(v292)
				v299 = v252 + int32(-1)
				if v299 != 0 {
					v252 = v299
					v254 = v254 + int32(4)
					v256 = v256 + int32(1)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v172 = v164 & int32(-4)
			v178 = l0 + v147<<(uint(int32(2))%32)
			v179 = l1 + v147
			v192 = v172
			for {
				v193 = int32(0)
				v194 = base.Simd_g_v128_load(m, v178, v193)
				v195 = base.Simd_g_const(&F_ConvertARGBToY_SSE41__k28)
				v199 = int32(16)
				base.Simd_g_v128_store32_lane_l0(m, v179, v193, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_v128_and(v194, v195), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k29)), base.Simd_g_i32x4_mul(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v194, v199), v195), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k30))), base.Simd_g_i32x4_mul(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v194, int32(8)), v195), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k31))), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k19)), v199), v194, base.Simd_g_const(&F_ConvertARGBToY_SSE41__k32), base.Simd_g_const(&F_ConvertARGBToY_SSE41__k33)))
				v225 = v192 + int32(-4)
				if v225 != 0 {
					v178 = v178 + v199
					v179 = v179 + int32(4)
					v192 = v225
					continue
				} else {
					break
				}
				break
			}
			if v164 == v172 {
			} else {
				v230 = v147 + v172
				v252 = l2 - v230
				v254 = l0 + v230<<(uint(int32(2))%32)
				v256 = l1 + v230
				for {
					v270 = *(*int32)(unsafe.Add(mBase, uint32(v254)))
					v271 = int32(255)
					v275 = int32(16)
					v292 = int32(base.Ui32(v270&v271*int32(_a_F_ConvertARGBToY_SSE41_0)+int32(base.Ui32(v270)>>(uint(v275)%32))&v271*int32(_a_F_ConvertARGBToY_SSE41_1)+int32(base.Ui32(v270)>>(uint(int32(8))%32))&v271*int32(_a_F_ConvertARGBToY_SSE41_2)+int32(1081344)) >> (uint(v275) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v256))) = uint8(v292)
					v299 = v252 + int32(-1)
					if v299 != 0 {
						v252 = v299
						v254 = v254 + int32(4)
						v256 = v256 + int32(1)
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

var F_ConvertARGBToY_SSE41__k0 = [2]uint64{0xd0905010c080400, 0xf0b07030e0a0602}
var F_ConvertARGBToY_SSE41__k1 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_ConvertARGBToY_SSE41__k2 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertARGBToY_SSE41__k3 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_ConvertARGBToY_SSE41__k4 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_ConvertARGBToY_SSE41__k5 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertARGBToY_SSE41__k6 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertARGBToY_SSE41__k7 = [2]uint64{0x0, 0x0}
var F_ConvertARGBToY_SSE41__k8 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertARGBToY_SSE41__k9 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertARGBToY_SSE41__k10 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertARGBToY_SSE41__k11 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_ConvertARGBToY_SSE41__k12 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_ConvertARGBToY_SSE41__k13 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_ConvertARGBToY_SSE41__k14 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertARGBToY_SSE41__k15 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertARGBToY_SSE41__k16 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertARGBToY_SSE41__k17 = [2]uint64{0x412341c7412341c7, 0x412341c7412341c7}
var F_ConvertARGBToY_SSE41__k18 = [2]uint64{0x1914400019144000, 0x1914400019144000}
var F_ConvertARGBToY_SSE41__k19 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertARGBToY_SSE41__k20 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertARGBToY_SSE41__k21 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertARGBToY_SSE41__k22 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertARGBToY_SSE41__k23 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertARGBToY_SSE41__k24 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertARGBToY_SSE41__k25 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertARGBToY_SSE41__k26 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertARGBToY_SSE41__k27 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertARGBToY_SSE41__k28 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ConvertARGBToY_SSE41__k29 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertARGBToY_SSE41__k30 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertARGBToY_SSE41__k31 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertARGBToY_SSE41__k32 = [2]uint64{0xc080400, 0x0}
var F_ConvertARGBToY_SSE41__k33 = [2]uint64{0x8080808080808080, 0x8080808080808080}

func F_ConvertBGR24ToY_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v44 base.V128
	_ = v44
	var v46 int32
	_ = v46
	var v47 base.V128
	_ = v47
	var v50 base.V128
	_ = v50
	var v52 int32
	_ = v52
	var v53 base.V128
	_ = v53
	var v56 base.V128
	_ = v56
	var v67 base.V128
	_ = v67
	var v70 base.V128
	_ = v70
	var v73 base.V128
	_ = v73
	var v76 base.V128
	_ = v76
	var v79 base.V128
	_ = v79
	var v89 base.V128
	_ = v89
	var v92 base.V128
	_ = v92
	var v95 base.V128
	_ = v95
	var v98 base.V128
	_ = v98
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v106 int32
	_ = v106
	var v108 base.V128
	_ = v108
	var v119 base.V128
	_ = v119
	var v122 base.V128
	_ = v122
	var v125 base.V128
	_ = v125
	var v128 base.V128
	_ = v128
	var v141 base.V128
	_ = v141
	var v144 base.V128
	_ = v144
	var v147 base.V128
	_ = v147
	var v150 base.V128
	_ = v150
	var v164 base.V128
	_ = v164
	var v167 base.V128
	_ = v167
	var v170 base.V128
	_ = v170
	var v173 base.V128
	_ = v173
	var v192 base.V128
	_ = v192
	var v195 base.V128
	_ = v195
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v214 base.V128
	_ = v214
	var v217 base.V128
	_ = v217
	var v220 base.V128
	_ = v220
	var v223 base.V128
	_ = v223
	var v237 base.V128
	_ = v237
	var v240 base.V128
	_ = v240
	var v243 base.V128
	_ = v243
	var v246 base.V128
	_ = v246
	var v264 base.V128
	_ = v264
	var v267 base.V128
	_ = v267
	var v270 base.V128
	_ = v270
	var v273 base.V128
	_ = v273
	var v286 base.V128
	_ = v286
	var v289 base.V128
	_ = v289
	var v292 base.V128
	_ = v292
	var v295 base.V128
	_ = v295
	var v309 base.V128
	_ = v309
	var v312 base.V128
	_ = v312
	var v315 base.V128
	_ = v315
	var v318 base.V128
	_ = v318
	var v328 base.V128
	_ = v328
	var v336 int32
	_ = v336
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v352 int32
	_ = v352
	var v356 int32
	_ = v356
	var v358 int32
	_ = v358
	var v364 int32
	_ = v364
	var v367 int32
	_ = v367
	var v371 int32
	_ = v371
	var v378 int32
	_ = v378
	var v385 int32
	_ = v385
	if l2 < int32(1) {
	} else {
		if base.Ui32(int32(16)) <= base.Ui32(l2) {
			v19 = l2 & int32(2147483632)
			v23 = l0
			v28 = v19
			v29 = l1
			for {
				v41 = int32(2)
				v43 = int32(0)
				v44 = base.Simd_g_v128_load8_splat(m, v23+v41, v43)
				v46 = int32(1)
				v47 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(5), v43, v44)
				v50 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(8), v43, v47)
				v52 = int32(3)
				v53 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(11), v43, v50)
				v56 = base.Simd_g_const(&F_ConvertBGR24ToY_C__k0)
				v67 = base.Simd_g_v128_load8_splat(m, v23+v46, v43)
				v70 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(4), v43, v67)
				v73 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(7), v43, v70)
				v76 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(10), v43, v73)
				v79 = base.Simd_g_const(&F_ConvertBGR24ToY_C__k1)
				v89 = base.Simd_g_v128_load8_splat(m, v23, v43)
				v92 = base.Simd_g_v128_load8_lane_l1(m, v23+v52, v43, v89)
				v95 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(6), v43, v92)
				v98 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(9), v43, v95)
				v101 = base.Simd_g_const(&F_ConvertBGR24ToY_C__k2)
				v104 = base.Simd_g_const(&F_ConvertBGR24ToY_C__k3)
				v106 = int32(16)
				v108 = base.Simd_g_const(&F_ConvertBGR24ToY_C__k4)
				v119 = base.Simd_g_v128_load8_splat(m, v23+int32(14), v43)
				v122 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(17), v43, v119)
				v125 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(20), v43, v122)
				v128 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(23), v43, v125)
				v141 = base.Simd_g_v128_load8_splat(m, v23+int32(13), v43)
				v144 = base.Simd_g_v128_load8_lane_l1(m, v23+v106, v43, v141)
				v147 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(19), v43, v144)
				v150 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(22), v43, v147)
				v164 = base.Simd_g_v128_load8_splat(m, v23+int32(12), v43)
				v167 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(15), v43, v164)
				v170 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(18), v43, v167)
				v173 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(21), v43, v170)
				v192 = base.Simd_g_v128_load8_splat(m, v23+int32(26), v43)
				v195 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(29), v43, v192)
				v198 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(32), v43, v195)
				v201 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(35), v43, v198)
				v214 = base.Simd_g_v128_load8_splat(m, v23+int32(25), v43)
				v217 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(28), v43, v214)
				v220 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(31), v43, v217)
				v223 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(34), v43, v220)
				v237 = base.Simd_g_v128_load8_splat(m, v23+int32(24), v43)
				v240 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(27), v43, v237)
				v243 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(30), v43, v240)
				v246 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(33), v43, v243)
				v264 = base.Simd_g_v128_load8_splat(m, v23+int32(38), v43)
				v267 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(41), v43, v264)
				v270 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(44), v43, v267)
				v273 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(47), v43, v270)
				v286 = base.Simd_g_v128_load8_splat(m, v23+int32(37), v43)
				v289 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(40), v43, v286)
				v292 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(43), v43, v289)
				v295 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(46), v43, v292)
				v309 = base.Simd_g_v128_load8_splat(m, v23+int32(36), v43)
				v312 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(39), v43, v309)
				v315 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(42), v43, v312)
				v318 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(45), v43, v315)
				v328 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v53)), v56), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v76)), v79)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v98)), v101)), v104), v106), v108), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v128)), v56), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v150)), v79)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v173)), v101)), v104), v106), v108)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v201)), v56), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v223)), v79)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v246)), v101)), v104), v106), v108), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v273)), v56), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v295)), v79)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v318)), v101)), v104), v106), v108)))
				base.Simd_g_v128_store(m, v29, v43, v328)
				v336 = v28 + int32(-16)
				if v336 != 0 {
					v23 = v23 + int32(48)
					v28 = v336
					v29 = v29 + v106
					continue
				} else {
					break
				}
				break
			}
			if v19 == l2 {
			} else {
				v341 = v19
				v342 = l0 + v19*int32(3)
				v352 = l1 + v341
				v356 = v342
				v358 = l2 - v341
				for {
					v364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+2)))
					v367 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+1)))
					v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356))))
					v378 = int32(base.Ui32(v364*int32(_a_F_ConvertBGR24ToY_C_0)+v367*int32(_a_F_ConvertBGR24ToY_C_1)+v371*int32(_a_F_ConvertBGR24ToY_C_2)+int32(1081344)) >> (uint(int32(16)) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v352))) = uint8(v378)
					v385 = v358 + int32(-1)
					if v385 != 0 {
						v352 = v352 + int32(1)
						v356 = v356 + int32(3)
						v358 = v385
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v341 = int32(0)
			v342 = l0
			v352 = l1 + v341
			v356 = v342
			v358 = l2 - v341
			for {
				v364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+2)))
				v367 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+1)))
				v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356))))
				v378 = int32(base.Ui32(v364*int32(_a_F_ConvertBGR24ToY_C_0)+v367*int32(_a_F_ConvertBGR24ToY_C_1)+v371*int32(_a_F_ConvertBGR24ToY_C_2)+int32(1081344)) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v352))) = uint8(v378)
				v385 = v358 + int32(-1)
				if v385 != 0 {
					v352 = v352 + int32(1)
					v356 = v356 + int32(3)
					v358 = v385
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

var F_ConvertBGR24ToY_C__k0 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertBGR24ToY_C__k1 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertBGR24ToY_C__k2 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertBGR24ToY_C__k3 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertBGR24ToY_C__k4 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_ConvertBGR24ToY_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v75 int32
	_ = v75
	var v76 base.V128
	_ = v76
	var v78 base.V128
	_ = v78
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
	var v90 base.V128
	_ = v90
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v106 base.V128
	_ = v106
	var v108 base.V128
	_ = v108
	var v109 base.V128
	_ = v109
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v129 base.V128
	_ = v129
	var v131 base.V128
	_ = v131
	var v133 base.V128
	_ = v133
	var v135 base.V128
	_ = v135
	var v138 base.V128
	_ = v138
	var v141 base.V128
	_ = v141
	var v145 base.V128
	_ = v145
	var v157 base.V128
	_ = v157
	var v159 base.V128
	_ = v159
	var v164 base.V128
	_ = v164
	var v183 base.V128
	_ = v183
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v191 base.V128
	_ = v191
	var v193 base.V128
	_ = v193
	var v198 base.V128
	_ = v198
	var v200 base.V128
	_ = v200
	var v220 base.V128
	_ = v220
	var v222 base.V128
	_ = v222
	var v227 base.V128
	_ = v227
	var v246 base.V128
	_ = v246
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v294 int32
	_ = v294
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v323 int32
	_ = v323
	var v330 int32
	_ = v330
	var v332 int32
	_ = v332
	var v333 base.V128
	_ = v333
	var v335 int32
	_ = v335
	var v336 base.V128
	_ = v336
	var v339 base.V128
	_ = v339
	var v341 int32
	_ = v341
	var v342 base.V128
	_ = v342
	var v345 base.V128
	_ = v345
	var v356 base.V128
	_ = v356
	var v359 base.V128
	_ = v359
	var v362 base.V128
	_ = v362
	var v365 base.V128
	_ = v365
	var v368 base.V128
	_ = v368
	var v378 base.V128
	_ = v378
	var v381 base.V128
	_ = v381
	var v384 base.V128
	_ = v384
	var v387 base.V128
	_ = v387
	var v390 base.V128
	_ = v390
	var v393 base.V128
	_ = v393
	var v395 int32
	_ = v395
	var v397 base.V128
	_ = v397
	var v408 base.V128
	_ = v408
	var v411 base.V128
	_ = v411
	var v414 base.V128
	_ = v414
	var v417 base.V128
	_ = v417
	var v430 base.V128
	_ = v430
	var v433 base.V128
	_ = v433
	var v436 base.V128
	_ = v436
	var v439 base.V128
	_ = v439
	var v453 base.V128
	_ = v453
	var v456 base.V128
	_ = v456
	var v459 base.V128
	_ = v459
	var v462 base.V128
	_ = v462
	var v481 base.V128
	_ = v481
	var v484 base.V128
	_ = v484
	var v487 base.V128
	_ = v487
	var v490 base.V128
	_ = v490
	var v503 base.V128
	_ = v503
	var v506 base.V128
	_ = v506
	var v509 base.V128
	_ = v509
	var v512 base.V128
	_ = v512
	var v526 base.V128
	_ = v526
	var v529 base.V128
	_ = v529
	var v532 base.V128
	_ = v532
	var v535 base.V128
	_ = v535
	var v553 base.V128
	_ = v553
	var v556 base.V128
	_ = v556
	var v559 base.V128
	_ = v559
	var v562 base.V128
	_ = v562
	var v575 base.V128
	_ = v575
	var v578 base.V128
	_ = v578
	var v581 base.V128
	_ = v581
	var v584 base.V128
	_ = v584
	var v598 base.V128
	_ = v598
	var v601 base.V128
	_ = v601
	var v604 base.V128
	_ = v604
	var v607 base.V128
	_ = v607
	var v617 base.V128
	_ = v617
	var v625 int32
	_ = v625
	var v631 int32
	_ = v631
	var v633 int32
	_ = v633
	var v657 int32
	_ = v657
	var v660 int32
	_ = v660
	var v664 int32
	_ = v664
	var v671 int32
	_ = v671
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v684 int32
	_ = v684
	var v687 int32
	_ = v687
	var v689 int32
	_ = v689
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v712 int32
	_ = v712
	var v713 int32
	_ = v713
	var v716 int32
	_ = v716
	var v717 int32
	_ = v717
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v727 int32
	_ = v727
	var v730 int32
	_ = v730
	var v734 int32
	_ = v734
	var v741 int32
	_ = v741
	var v748 int32
	_ = v748
	v4 = int32(0)
	v28 = l2 & int32(-32)
	if v28 < int32(1) {
		v254 = l0
		v257 = v4
	} else {
		v31 = l0
		v34 = v4
		for {
			v56 = l1 + v34
			v57 = int32(16)
			v59 = int32(0)
			v60 = base.Simd_g_v128_load_rng(m, v31, v59, int32(0), int32(96))
			v62 = base.Simd_g_v128_load_nc(m, v31, int32(48))
			v63 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k0)
			v64 = base.Simd_g_i8x16_shuffle2(v60, v62, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v66 = base.Simd_g_v128_load_nc(m, v31, v57)
			v68 = base.Simd_g_v128_load_nc(m, v31, int32(64))
			v69 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k3)
			v70 = base.Simd_g_i8x16_shuffle2(v66, v68, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v72 = base.Simd_g_i8x16_shuffle2(v64, v70, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v74 = base.Simd_g_i8x16_shuffle2(v60, v62, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v75 = int32(32)
			v76 = base.Simd_g_v128_load_nc(m, v31, v75)
			v78 = base.Simd_g_v128_load_nc(m, v31, int32(80))
			v80 = base.Simd_g_i8x16_shuffle2(v76, v78, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v82 = base.Simd_g_i8x16_shuffle2(v74, v80, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v84 = base.Simd_g_i8x16_shuffle2(v72, v82, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v86 = base.Simd_g_i8x16_shuffle2(v64, v70, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v88 = base.Simd_g_i8x16_shuffle2(v66, v68, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v90 = base.Simd_g_i8x16_shuffle2(v76, v78, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v92 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v94 = base.Simd_g_i8x16_shuffle2(v86, v92, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v96 = base.Simd_g_i8x16_shuffle2(v84, v94, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v98 = base.Simd_g_i8x16_shuffle2(v86, v92, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v100 = base.Simd_g_i8x16_shuffle2(v74, v80, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v102 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v104 = base.Simd_g_i8x16_shuffle2(v100, v102, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v106 = base.Simd_g_i8x16_shuffle2(v98, v104, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v108 = base.Simd_g_i8x16_shuffle2(v96, v106, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v109 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k6)
			v111 = base.Simd_g_i8x16_shuffle2(v108, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v113 = base.Simd_g_i8x16_shuffle2(v84, v94, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v115 = base.Simd_g_i8x16_shuffle2(v72, v82, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v117 = base.Simd_g_i8x16_shuffle2(v100, v102, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v119 = base.Simd_g_i8x16_shuffle2(v115, v117, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v121 = base.Simd_g_i8x16_shuffle2(v113, v119, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v123 = base.Simd_g_i8x16_shuffle2(v121, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v124 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k7)
			v126 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k8)
			v129 = base.Simd_g_i8x16_shuffle2(v115, v117, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v131 = base.Simd_g_i8x16_shuffle2(v98, v104, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v133 = base.Simd_g_i8x16_shuffle2(v129, v131, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v135 = base.Simd_g_i8x16_shuffle2(v133, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v138 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k9)
			v141 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k10)
			v145 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k11)
			v157 = base.Simd_g_i8x16_shuffle2(v108, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v159 = base.Simd_g_i8x16_shuffle2(v121, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v164 = base.Simd_g_i8x16_shuffle2(v133, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v183 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v111, v123, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v135, v111, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v111, v123, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v135, v111, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v138)), v141), v57)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v157, v159, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v164, v157, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v157, v159, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v164, v157, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v138)), v141), v57)))
			base.Simd_g_v128_store(m, v56+v57, v59, v183)
			v187 = base.Simd_g_i8x16_shuffle2(v96, v106, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v189 = base.Simd_g_i8x16_shuffle2(v187, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v191 = base.Simd_g_i8x16_shuffle2(v113, v119, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v193 = base.Simd_g_i8x16_shuffle2(v191, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v198 = base.Simd_g_i8x16_shuffle2(v129, v131, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v200 = base.Simd_g_i8x16_shuffle2(v198, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k2))
			v220 = base.Simd_g_i8x16_shuffle2(v187, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v222 = base.Simd_g_i8x16_shuffle2(v191, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v227 = base.Simd_g_i8x16_shuffle2(v198, v109, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k5))
			v246 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v189, v193, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v200, v189, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v189, v193, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v200, v189, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v138)), v141), v57)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v220, v222, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v227, v220, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v220, v222, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v227, v220, base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k15)), v138)), v141), v57)))
			base.Simd_g_v128_store(m, v56, v59, v246)
			v250 = v31 + int32(96)
			v252 = v34 + v75
			if v252 < v28 {
				v31 = v250
				v34 = v252
				continue
			} else {
				break
			}
			break
		}
		v254 = v250
		v257 = v252
	}
	if l2 <= v257 {
	} else {
		v280 = l2 - v257
		if base.Ui32(v280) <= base.Ui32(int32(15)) {
			v631 = v257
			v633 = v254
			if (l2-v631)&int32(1) != 0 {
				v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
				v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
				v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
				v671 = int32(base.Ui32(v657*int32(_a_F_ConvertBGR24ToY_SSE2_0)+v660*int32(_a_F_ConvertBGR24ToY_SSE2_1)+v664*int32(_a_F_ConvertBGR24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
				v677 = v631 + int32(1)
				v678 = v633 + int32(3)
			} else {
				v677 = v631
				v678 = v633
			}
			if v631 == l2+int32(-1) {
			} else {
				v684 = l1 + v677
				v687 = l2 - v677
				v689 = v678
				for {
					v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
					v710 = int32(_a_F_ConvertBGR24ToY_SSE2_0)
					v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
					v713 = int32(_a_F_ConvertBGR24ToY_SSE2_1)
					v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
					v717 = int32(_a_F_ConvertBGR24ToY_SSE2_2)
					v720 = int32(1081344)
					v722 = int32(16)
					v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
					v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
					v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
					v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
					v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
					v748 = v687 + int32(-2)
					if v748 != 0 {
						v684 = v684 + int32(2)
						v687 = v748
						v689 = v689 + int32(6)
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v283 = l1 + v257
			if base.Ui32(v254+(v257*int32(-3)+l2*int32(3))) <= base.Ui32(v283) {
				v294 = v280 & int32(-16)
				v299 = v254
				v303 = v283
				v323 = v294
				for {
					v330 = int32(2)
					v332 = int32(0)
					v333 = base.Simd_g_v128_load8_splat(m, v299+v330, v332)
					v335 = int32(1)
					v336 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(5), v332, v333)
					v339 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(8), v332, v336)
					v341 = int32(3)
					v342 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(11), v332, v339)
					v345 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k16)
					v356 = base.Simd_g_v128_load8_splat(m, v299+v335, v332)
					v359 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(4), v332, v356)
					v362 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(7), v332, v359)
					v365 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(10), v332, v362)
					v368 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k17)
					v378 = base.Simd_g_v128_load8_splat(m, v299, v332)
					v381 = base.Simd_g_v128_load8_lane_l1(m, v299+v341, v332, v378)
					v384 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(6), v332, v381)
					v387 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(9), v332, v384)
					v390 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k18)
					v393 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k10)
					v395 = int32(16)
					v397 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k19)
					v408 = base.Simd_g_v128_load8_splat(m, v299+int32(14), v332)
					v411 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(17), v332, v408)
					v414 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(20), v332, v411)
					v417 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(23), v332, v414)
					v430 = base.Simd_g_v128_load8_splat(m, v299+int32(13), v332)
					v433 = base.Simd_g_v128_load8_lane_l1(m, v299+v395, v332, v430)
					v436 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(19), v332, v433)
					v439 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(22), v332, v436)
					v453 = base.Simd_g_v128_load8_splat(m, v299+int32(12), v332)
					v456 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(15), v332, v453)
					v459 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(18), v332, v456)
					v462 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(21), v332, v459)
					v481 = base.Simd_g_v128_load8_splat(m, v299+int32(26), v332)
					v484 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(29), v332, v481)
					v487 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(32), v332, v484)
					v490 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(35), v332, v487)
					v503 = base.Simd_g_v128_load8_splat(m, v299+int32(25), v332)
					v506 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(28), v332, v503)
					v509 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(31), v332, v506)
					v512 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(34), v332, v509)
					v526 = base.Simd_g_v128_load8_splat(m, v299+int32(24), v332)
					v529 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(27), v332, v526)
					v532 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(30), v332, v529)
					v535 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(33), v332, v532)
					v553 = base.Simd_g_v128_load8_splat(m, v299+int32(38), v332)
					v556 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(41), v332, v553)
					v559 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(44), v332, v556)
					v562 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(47), v332, v559)
					v575 = base.Simd_g_v128_load8_splat(m, v299+int32(37), v332)
					v578 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(40), v332, v575)
					v581 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(43), v332, v578)
					v584 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(46), v332, v581)
					v598 = base.Simd_g_v128_load8_splat(m, v299+int32(36), v332)
					v601 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(39), v332, v598)
					v604 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(42), v332, v601)
					v607 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(45), v332, v604)
					v617 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v342)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v365)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v387)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v417)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v439)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v462)), v390)), v393), v395), v397)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v490)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v512)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v535)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v562)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v584)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v607)), v390)), v393), v395), v397)))
					base.Simd_g_v128_store(m, v303, v332, v617)
					v625 = v323 + int32(-16)
					if v625 != 0 {
						v299 = v299 + int32(48)
						v303 = v303 + v395
						v323 = v625
						continue
					} else {
						break
					}
					break
				}
				if v280 != v294 {
					v631 = v257 + v294
					v633 = v254 + v294*int32(3)
					if (l2-v631)&int32(1) != 0 {
						v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
						v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
						v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
						v671 = int32(base.Ui32(v657*int32(_a_F_ConvertBGR24ToY_SSE2_0)+v660*int32(_a_F_ConvertBGR24ToY_SSE2_1)+v664*int32(_a_F_ConvertBGR24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
						v677 = v631 + int32(1)
						v678 = v633 + int32(3)
					} else {
						v677 = v631
						v678 = v633
					}
					if v631 == l2+int32(-1) {
					} else {
						v684 = l1 + v677
						v687 = l2 - v677
						v689 = v678
						for {
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
							v710 = int32(_a_F_ConvertBGR24ToY_SSE2_0)
							v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
							v713 = int32(_a_F_ConvertBGR24ToY_SSE2_1)
							v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
							v717 = int32(_a_F_ConvertBGR24ToY_SSE2_2)
							v720 = int32(1081344)
							v722 = int32(16)
							v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
							v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
							v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
							v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
							v748 = v687 + int32(-2)
							if v748 != 0 {
								v684 = v684 + int32(2)
								v687 = v748
								v689 = v689 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
				}
			} else {
				if base.Ui32(v254) < base.Ui32(l1+l2) {
					v631 = v257
					v633 = v254
					if (l2-v631)&int32(1) != 0 {
						v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
						v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
						v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
						v671 = int32(base.Ui32(v657*int32(_a_F_ConvertBGR24ToY_SSE2_0)+v660*int32(_a_F_ConvertBGR24ToY_SSE2_1)+v664*int32(_a_F_ConvertBGR24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
						v677 = v631 + int32(1)
						v678 = v633 + int32(3)
					} else {
						v677 = v631
						v678 = v633
					}
					if v631 == l2+int32(-1) {
					} else {
						v684 = l1 + v677
						v687 = l2 - v677
						v689 = v678
						for {
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
							v710 = int32(_a_F_ConvertBGR24ToY_SSE2_0)
							v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
							v713 = int32(_a_F_ConvertBGR24ToY_SSE2_1)
							v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
							v717 = int32(_a_F_ConvertBGR24ToY_SSE2_2)
							v720 = int32(1081344)
							v722 = int32(16)
							v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
							v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
							v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
							v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
							v748 = v687 + int32(-2)
							if v748 != 0 {
								v684 = v684 + int32(2)
								v687 = v748
								v689 = v689 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v294 = v280 & int32(-16)
					v299 = v254
					v303 = v283
					v323 = v294
					for {
						v330 = int32(2)
						v332 = int32(0)
						v333 = base.Simd_g_v128_load8_splat(m, v299+v330, v332)
						v335 = int32(1)
						v336 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(5), v332, v333)
						v339 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(8), v332, v336)
						v341 = int32(3)
						v342 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(11), v332, v339)
						v345 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k16)
						v356 = base.Simd_g_v128_load8_splat(m, v299+v335, v332)
						v359 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(4), v332, v356)
						v362 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(7), v332, v359)
						v365 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(10), v332, v362)
						v368 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k17)
						v378 = base.Simd_g_v128_load8_splat(m, v299, v332)
						v381 = base.Simd_g_v128_load8_lane_l1(m, v299+v341, v332, v378)
						v384 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(6), v332, v381)
						v387 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(9), v332, v384)
						v390 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k18)
						v393 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k10)
						v395 = int32(16)
						v397 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE2__k19)
						v408 = base.Simd_g_v128_load8_splat(m, v299+int32(14), v332)
						v411 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(17), v332, v408)
						v414 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(20), v332, v411)
						v417 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(23), v332, v414)
						v430 = base.Simd_g_v128_load8_splat(m, v299+int32(13), v332)
						v433 = base.Simd_g_v128_load8_lane_l1(m, v299+v395, v332, v430)
						v436 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(19), v332, v433)
						v439 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(22), v332, v436)
						v453 = base.Simd_g_v128_load8_splat(m, v299+int32(12), v332)
						v456 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(15), v332, v453)
						v459 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(18), v332, v456)
						v462 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(21), v332, v459)
						v481 = base.Simd_g_v128_load8_splat(m, v299+int32(26), v332)
						v484 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(29), v332, v481)
						v487 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(32), v332, v484)
						v490 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(35), v332, v487)
						v503 = base.Simd_g_v128_load8_splat(m, v299+int32(25), v332)
						v506 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(28), v332, v503)
						v509 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(31), v332, v506)
						v512 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(34), v332, v509)
						v526 = base.Simd_g_v128_load8_splat(m, v299+int32(24), v332)
						v529 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(27), v332, v526)
						v532 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(30), v332, v529)
						v535 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(33), v332, v532)
						v553 = base.Simd_g_v128_load8_splat(m, v299+int32(38), v332)
						v556 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(41), v332, v553)
						v559 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(44), v332, v556)
						v562 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(47), v332, v559)
						v575 = base.Simd_g_v128_load8_splat(m, v299+int32(37), v332)
						v578 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(40), v332, v575)
						v581 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(43), v332, v578)
						v584 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(46), v332, v581)
						v598 = base.Simd_g_v128_load8_splat(m, v299+int32(36), v332)
						v601 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(39), v332, v598)
						v604 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(42), v332, v601)
						v607 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(45), v332, v604)
						v617 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v342)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v365)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v387)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v417)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v439)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v462)), v390)), v393), v395), v397)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v490)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v512)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v535)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v562)), v345), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v584)), v368)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v607)), v390)), v393), v395), v397)))
						base.Simd_g_v128_store(m, v303, v332, v617)
						v625 = v323 + int32(-16)
						if v625 != 0 {
							v299 = v299 + int32(48)
							v303 = v303 + v395
							v323 = v625
							continue
						} else {
							break
						}
						break
					}
					if v280 != v294 {
						v631 = v257 + v294
						v633 = v254 + v294*int32(3)
						if (l2-v631)&int32(1) != 0 {
							v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
							v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
							v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
							v671 = int32(base.Ui32(v657*int32(_a_F_ConvertBGR24ToY_SSE2_0)+v660*int32(_a_F_ConvertBGR24ToY_SSE2_1)+v664*int32(_a_F_ConvertBGR24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
							v677 = v631 + int32(1)
							v678 = v633 + int32(3)
						} else {
							v677 = v631
							v678 = v633
						}
						if v631 == l2+int32(-1) {
						} else {
							v684 = l1 + v677
							v687 = l2 - v677
							v689 = v678
							for {
								v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
								v710 = int32(_a_F_ConvertBGR24ToY_SSE2_0)
								v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
								v713 = int32(_a_F_ConvertBGR24ToY_SSE2_1)
								v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
								v717 = int32(_a_F_ConvertBGR24ToY_SSE2_2)
								v720 = int32(1081344)
								v722 = int32(16)
								v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
								v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
								v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
								v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
								v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
								v748 = v687 + int32(-2)
								if v748 != 0 {
									v684 = v684 + int32(2)
									v687 = v748
									v689 = v689 + int32(6)
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
					}
				}
			}
		}
	}
	return
}

var F_ConvertBGR24ToY_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertBGR24ToY_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertBGR24ToY_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertBGR24ToY_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertBGR24ToY_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertBGR24ToY_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertBGR24ToY_SSE2__k6 = [2]uint64{0x0, 0x0}
var F_ConvertBGR24ToY_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertBGR24ToY_SSE2__k8 = [2]uint64{0x1914400019144000, 0x1914400019144000}
var F_ConvertBGR24ToY_SSE2__k9 = [2]uint64{0x412341c7412341c7, 0x412341c7412341c7}
var F_ConvertBGR24ToY_SSE2__k10 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertBGR24ToY_SSE2__k11 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertBGR24ToY_SSE2__k12 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertBGR24ToY_SSE2__k13 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertBGR24ToY_SSE2__k14 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertBGR24ToY_SSE2__k15 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertBGR24ToY_SSE2__k16 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertBGR24ToY_SSE2__k17 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertBGR24ToY_SSE2__k18 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertBGR24ToY_SSE2__k19 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_ConvertBGR24ToY_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v78 base.V128
	_ = v78
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v81 base.V128
	_ = v81
	var v82 base.V128
	_ = v82
	var v84 base.V128
	_ = v84
	var v87 base.V128
	_ = v87
	var v89 base.V128
	_ = v89
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v101 base.V128
	_ = v101
	var v103 base.V128
	_ = v103
	var v105 base.V128
	_ = v105
	var v108 base.V128
	_ = v108
	var v111 base.V128
	_ = v111
	var v115 base.V128
	_ = v115
	var v126 base.V128
	_ = v126
	var v127 base.V128
	_ = v127
	var v129 base.V128
	_ = v129
	var v134 base.V128
	_ = v134
	var v153 base.V128
	_ = v153
	var v154 int32
	_ = v154
	var v157 base.V128
	_ = v157
	var v160 base.V128
	_ = v160
	var v163 int32
	_ = v163
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v173 base.V128
	_ = v173
	var v175 base.V128
	_ = v175
	var v183 base.V128
	_ = v183
	var v185 base.V128
	_ = v185
	var v205 base.V128
	_ = v205
	var v207 base.V128
	_ = v207
	var v212 base.V128
	_ = v212
	var v231 base.V128
	_ = v231
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v282 int32
	_ = v282
	var v287 int32
	_ = v287
	var v291 int32
	_ = v291
	var v314 int32
	_ = v314
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v324 base.V128
	_ = v324
	var v326 int32
	_ = v326
	var v327 base.V128
	_ = v327
	var v330 base.V128
	_ = v330
	var v332 int32
	_ = v332
	var v333 base.V128
	_ = v333
	var v336 base.V128
	_ = v336
	var v347 base.V128
	_ = v347
	var v350 base.V128
	_ = v350
	var v353 base.V128
	_ = v353
	var v356 base.V128
	_ = v356
	var v359 base.V128
	_ = v359
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
	var v384 base.V128
	_ = v384
	var v386 int32
	_ = v386
	var v388 base.V128
	_ = v388
	var v399 base.V128
	_ = v399
	var v402 base.V128
	_ = v402
	var v405 base.V128
	_ = v405
	var v408 base.V128
	_ = v408
	var v421 base.V128
	_ = v421
	var v424 base.V128
	_ = v424
	var v427 base.V128
	_ = v427
	var v430 base.V128
	_ = v430
	var v444 base.V128
	_ = v444
	var v447 base.V128
	_ = v447
	var v450 base.V128
	_ = v450
	var v453 base.V128
	_ = v453
	var v472 base.V128
	_ = v472
	var v475 base.V128
	_ = v475
	var v478 base.V128
	_ = v478
	var v481 base.V128
	_ = v481
	var v494 base.V128
	_ = v494
	var v497 base.V128
	_ = v497
	var v500 base.V128
	_ = v500
	var v503 base.V128
	_ = v503
	var v517 base.V128
	_ = v517
	var v520 base.V128
	_ = v520
	var v523 base.V128
	_ = v523
	var v526 base.V128
	_ = v526
	var v544 base.V128
	_ = v544
	var v547 base.V128
	_ = v547
	var v550 base.V128
	_ = v550
	var v553 base.V128
	_ = v553
	var v566 base.V128
	_ = v566
	var v569 base.V128
	_ = v569
	var v572 base.V128
	_ = v572
	var v575 base.V128
	_ = v575
	var v589 base.V128
	_ = v589
	var v592 base.V128
	_ = v592
	var v595 base.V128
	_ = v595
	var v598 base.V128
	_ = v598
	var v608 base.V128
	_ = v608
	var v616 int32
	_ = v616
	var v622 int32
	_ = v622
	var v624 int32
	_ = v624
	var v651 int32
	_ = v651
	var v654 int32
	_ = v654
	var v658 int32
	_ = v658
	var v665 int32
	_ = v665
	var v671 int32
	_ = v671
	var v672 int32
	_ = v672
	var v678 int32
	_ = v678
	var v681 int32
	_ = v681
	var v683 int32
	_ = v683
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v713 int32
	_ = v713
	var v714 int32
	_ = v714
	var v717 int32
	_ = v717
	var v719 int32
	_ = v719
	var v720 int32
	_ = v720
	var v724 int32
	_ = v724
	var v727 int32
	_ = v727
	var v731 int32
	_ = v731
	var v738 int32
	_ = v738
	var v745 int32
	_ = v745
	v4 = int32(0)
	v31 = l2 & int32(-32)
	if v31 < int32(1) {
		v239 = l0
		v242 = v4
	} else {
		v34 = l0
		v37 = v4
		for {
			v62 = l1 + v37
			v63 = int32(16)
			v66 = base.Simd_g_v128_load_rng(m, v34, int32(64), int32(48), int32(48))
			v67 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k0)
			v70 = base.Simd_g_v128_load_nc(m, v34, int32(48))
			v71 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k1)
			v75 = base.Simd_g_v128_load_nc(m, v34, int32(80))
			v76 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k2)
			v78 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v66, v67), base.Simd_g_i8x16_swizzle(v70, v71)), base.Simd_g_i8x16_swizzle(v75, v76))
			v79 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k3)
			v80 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k4)
			v81 = base.Simd_g_i8x16_shuffle2(v78, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k6))
			v82 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k7)
			v84 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k8)
			v87 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k9)
			v89 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v66, v82), base.Simd_g_i8x16_swizzle(v70, v84)), base.Simd_g_i8x16_swizzle(v75, v87))
			v91 = base.Simd_g_i8x16_shuffle2(v89, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k6))
			v92 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k10)
			v94 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k11)
			v96 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k12)
			v98 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k13)
			v101 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k14)
			v103 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v66, v96), base.Simd_g_i8x16_swizzle(v70, v98)), base.Simd_g_i8x16_swizzle(v75, v101))
			v105 = base.Simd_g_i8x16_shuffle2(v103, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k6))
			v108 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k15)
			v111 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k16)
			v115 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k17)
			v126 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k18)
			v127 = base.Simd_g_i8x16_shuffle2(v78, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k20))
			v129 = base.Simd_g_i8x16_shuffle2(v89, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k20))
			v134 = base.Simd_g_i8x16_shuffle2(v103, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k20))
			v153 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v81, v91, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v105, v81, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v81, v91, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v105, v81, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v108)), v111), v63)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v134, v127, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v134, v127, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v108)), v111), v63)))
			v154 = int32(0)
			base.Simd_g_v128_store(m, v62+v63, v154, v153)
			v157 = base.Simd_g_v128_load_rng(m, v34, v63, int32(0), int32(48))
			v160 = base.Simd_g_v128_load_nc(m, v34, v154)
			v163 = int32(32)
			v164 = base.Simd_g_v128_load_nc(m, v34, v163)
			v166 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v157, v67), base.Simd_g_i8x16_swizzle(v160, v71)), base.Simd_g_i8x16_swizzle(v164, v76))
			v168 = base.Simd_g_i8x16_shuffle2(v166, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k6))
			v173 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v157, v82), base.Simd_g_i8x16_swizzle(v160, v84)), base.Simd_g_i8x16_swizzle(v164, v87))
			v175 = base.Simd_g_i8x16_shuffle2(v173, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k6))
			v183 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v157, v96), base.Simd_g_i8x16_swizzle(v160, v98)), base.Simd_g_i8x16_swizzle(v164, v101))
			v185 = base.Simd_g_i8x16_shuffle2(v183, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k6))
			v205 = base.Simd_g_i8x16_shuffle2(v166, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k20))
			v207 = base.Simd_g_i8x16_shuffle2(v173, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k20))
			v212 = base.Simd_g_i8x16_shuffle2(v183, v79, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k20))
			v231 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v168, v175, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v185, v168, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v168, v175, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v185, v168, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v108)), v111), v63)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v205, v207, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v212, v205, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v205, v207, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v212, v205, base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k24)), v108)), v111), v63)))
			base.Simd_g_v128_store(m, v62, v154, v231)
			v235 = v34 + int32(96)
			v237 = v37 + v163
			if v237 < v31 {
				v34 = v235
				v37 = v237
				continue
			} else {
				break
			}
			break
		}
		v239 = v235
		v242 = v237
	}
	if l2 <= v242 {
	} else {
		v268 = l2 - v242
		if base.Ui32(v268) <= base.Ui32(int32(15)) {
			v622 = v242
			v624 = v239
			if (l2-v622)&int32(1) != 0 {
				v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
				v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
				v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
				v665 = int32(base.Ui32(v651*int32(_a_F_ConvertBGR24ToY_SSE41_0)+v654*int32(_a_F_ConvertBGR24ToY_SSE41_1)+v658*int32(_a_F_ConvertBGR24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
				v671 = v622 + int32(1)
				v672 = v624 + int32(3)
			} else {
				v671 = v622
				v672 = v624
			}
			if v622 == l2+int32(-1) {
			} else {
				v678 = l1 + v671
				v681 = l2 - v671
				v683 = v672
				for {
					v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
					v707 = int32(_a_F_ConvertBGR24ToY_SSE41_0)
					v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
					v710 = int32(_a_F_ConvertBGR24ToY_SSE41_1)
					v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
					v714 = int32(_a_F_ConvertBGR24ToY_SSE41_2)
					v717 = int32(1081344)
					v719 = int32(16)
					v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
					v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
					v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
					v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
					v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
					v745 = v681 + int32(-2)
					if v745 != 0 {
						v678 = v678 + int32(2)
						v681 = v745
						v683 = v683 + int32(6)
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v271 = l1 + v242
			if base.Ui32(v239+(v242*int32(-3)+l2*int32(3))) <= base.Ui32(v271) {
				v282 = v268 & int32(-16)
				v287 = v239
				v291 = v271
				v314 = v282
				for {
					v321 = int32(2)
					v323 = int32(0)
					v324 = base.Simd_g_v128_load8_splat(m, v287+v321, v323)
					v326 = int32(1)
					v327 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(5), v323, v324)
					v330 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(8), v323, v327)
					v332 = int32(3)
					v333 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(11), v323, v330)
					v336 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k25)
					v347 = base.Simd_g_v128_load8_splat(m, v287+v326, v323)
					v350 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(4), v323, v347)
					v353 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(7), v323, v350)
					v356 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(10), v323, v353)
					v359 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k26)
					v369 = base.Simd_g_v128_load8_splat(m, v287, v323)
					v372 = base.Simd_g_v128_load8_lane_l1(m, v287+v332, v323, v369)
					v375 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(6), v323, v372)
					v378 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(9), v323, v375)
					v381 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k27)
					v384 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k16)
					v386 = int32(16)
					v388 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k28)
					v399 = base.Simd_g_v128_load8_splat(m, v287+int32(14), v323)
					v402 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(17), v323, v399)
					v405 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(20), v323, v402)
					v408 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(23), v323, v405)
					v421 = base.Simd_g_v128_load8_splat(m, v287+int32(13), v323)
					v424 = base.Simd_g_v128_load8_lane_l1(m, v287+v386, v323, v421)
					v427 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(19), v323, v424)
					v430 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(22), v323, v427)
					v444 = base.Simd_g_v128_load8_splat(m, v287+int32(12), v323)
					v447 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(15), v323, v444)
					v450 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(18), v323, v447)
					v453 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(21), v323, v450)
					v472 = base.Simd_g_v128_load8_splat(m, v287+int32(26), v323)
					v475 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(29), v323, v472)
					v478 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(32), v323, v475)
					v481 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(35), v323, v478)
					v494 = base.Simd_g_v128_load8_splat(m, v287+int32(25), v323)
					v497 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(28), v323, v494)
					v500 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(31), v323, v497)
					v503 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(34), v323, v500)
					v517 = base.Simd_g_v128_load8_splat(m, v287+int32(24), v323)
					v520 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(27), v323, v517)
					v523 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(30), v323, v520)
					v526 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(33), v323, v523)
					v544 = base.Simd_g_v128_load8_splat(m, v287+int32(38), v323)
					v547 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(41), v323, v544)
					v550 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(44), v323, v547)
					v553 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(47), v323, v550)
					v566 = base.Simd_g_v128_load8_splat(m, v287+int32(37), v323)
					v569 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(40), v323, v566)
					v572 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(43), v323, v569)
					v575 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(46), v323, v572)
					v589 = base.Simd_g_v128_load8_splat(m, v287+int32(36), v323)
					v592 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(39), v323, v589)
					v595 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(42), v323, v592)
					v598 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(45), v323, v595)
					v608 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v333)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v356)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v378)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v408)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v430)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v453)), v381)), v384), v386), v388)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v481)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v503)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v526)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v553)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v575)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v598)), v381)), v384), v386), v388)))
					base.Simd_g_v128_store(m, v291, v323, v608)
					v616 = v314 + int32(-16)
					if v616 != 0 {
						v287 = v287 + int32(48)
						v291 = v291 + v386
						v314 = v616
						continue
					} else {
						break
					}
					break
				}
				if v268 != v282 {
					v622 = v242 + v282
					v624 = v239 + v282*int32(3)
					if (l2-v622)&int32(1) != 0 {
						v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
						v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
						v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
						v665 = int32(base.Ui32(v651*int32(_a_F_ConvertBGR24ToY_SSE41_0)+v654*int32(_a_F_ConvertBGR24ToY_SSE41_1)+v658*int32(_a_F_ConvertBGR24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
						v671 = v622 + int32(1)
						v672 = v624 + int32(3)
					} else {
						v671 = v622
						v672 = v624
					}
					if v622 == l2+int32(-1) {
					} else {
						v678 = l1 + v671
						v681 = l2 - v671
						v683 = v672
						for {
							v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
							v707 = int32(_a_F_ConvertBGR24ToY_SSE41_0)
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
							v710 = int32(_a_F_ConvertBGR24ToY_SSE41_1)
							v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
							v714 = int32(_a_F_ConvertBGR24ToY_SSE41_2)
							v717 = int32(1081344)
							v719 = int32(16)
							v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
							v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
							v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
							v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
							v745 = v681 + int32(-2)
							if v745 != 0 {
								v678 = v678 + int32(2)
								v681 = v745
								v683 = v683 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
				}
			} else {
				if base.Ui32(v239) < base.Ui32(l1+l2) {
					v622 = v242
					v624 = v239
					if (l2-v622)&int32(1) != 0 {
						v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
						v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
						v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
						v665 = int32(base.Ui32(v651*int32(_a_F_ConvertBGR24ToY_SSE41_0)+v654*int32(_a_F_ConvertBGR24ToY_SSE41_1)+v658*int32(_a_F_ConvertBGR24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
						v671 = v622 + int32(1)
						v672 = v624 + int32(3)
					} else {
						v671 = v622
						v672 = v624
					}
					if v622 == l2+int32(-1) {
					} else {
						v678 = l1 + v671
						v681 = l2 - v671
						v683 = v672
						for {
							v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
							v707 = int32(_a_F_ConvertBGR24ToY_SSE41_0)
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
							v710 = int32(_a_F_ConvertBGR24ToY_SSE41_1)
							v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
							v714 = int32(_a_F_ConvertBGR24ToY_SSE41_2)
							v717 = int32(1081344)
							v719 = int32(16)
							v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
							v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
							v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
							v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
							v745 = v681 + int32(-2)
							if v745 != 0 {
								v678 = v678 + int32(2)
								v681 = v745
								v683 = v683 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v282 = v268 & int32(-16)
					v287 = v239
					v291 = v271
					v314 = v282
					for {
						v321 = int32(2)
						v323 = int32(0)
						v324 = base.Simd_g_v128_load8_splat(m, v287+v321, v323)
						v326 = int32(1)
						v327 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(5), v323, v324)
						v330 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(8), v323, v327)
						v332 = int32(3)
						v333 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(11), v323, v330)
						v336 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k25)
						v347 = base.Simd_g_v128_load8_splat(m, v287+v326, v323)
						v350 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(4), v323, v347)
						v353 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(7), v323, v350)
						v356 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(10), v323, v353)
						v359 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k26)
						v369 = base.Simd_g_v128_load8_splat(m, v287, v323)
						v372 = base.Simd_g_v128_load8_lane_l1(m, v287+v332, v323, v369)
						v375 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(6), v323, v372)
						v378 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(9), v323, v375)
						v381 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k27)
						v384 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k16)
						v386 = int32(16)
						v388 = base.Simd_g_const(&F_ConvertBGR24ToY_SSE41__k28)
						v399 = base.Simd_g_v128_load8_splat(m, v287+int32(14), v323)
						v402 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(17), v323, v399)
						v405 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(20), v323, v402)
						v408 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(23), v323, v405)
						v421 = base.Simd_g_v128_load8_splat(m, v287+int32(13), v323)
						v424 = base.Simd_g_v128_load8_lane_l1(m, v287+v386, v323, v421)
						v427 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(19), v323, v424)
						v430 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(22), v323, v427)
						v444 = base.Simd_g_v128_load8_splat(m, v287+int32(12), v323)
						v447 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(15), v323, v444)
						v450 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(18), v323, v447)
						v453 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(21), v323, v450)
						v472 = base.Simd_g_v128_load8_splat(m, v287+int32(26), v323)
						v475 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(29), v323, v472)
						v478 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(32), v323, v475)
						v481 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(35), v323, v478)
						v494 = base.Simd_g_v128_load8_splat(m, v287+int32(25), v323)
						v497 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(28), v323, v494)
						v500 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(31), v323, v497)
						v503 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(34), v323, v500)
						v517 = base.Simd_g_v128_load8_splat(m, v287+int32(24), v323)
						v520 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(27), v323, v517)
						v523 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(30), v323, v520)
						v526 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(33), v323, v523)
						v544 = base.Simd_g_v128_load8_splat(m, v287+int32(38), v323)
						v547 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(41), v323, v544)
						v550 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(44), v323, v547)
						v553 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(47), v323, v550)
						v566 = base.Simd_g_v128_load8_splat(m, v287+int32(37), v323)
						v569 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(40), v323, v566)
						v572 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(43), v323, v569)
						v575 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(46), v323, v572)
						v589 = base.Simd_g_v128_load8_splat(m, v287+int32(36), v323)
						v592 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(39), v323, v589)
						v595 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(42), v323, v592)
						v598 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(45), v323, v595)
						v608 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v333)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v356)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v378)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v408)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v430)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v453)), v381)), v384), v386), v388)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v481)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v503)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v526)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v553)), v336), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v575)), v359)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v598)), v381)), v384), v386), v388)))
						base.Simd_g_v128_store(m, v291, v323, v608)
						v616 = v314 + int32(-16)
						if v616 != 0 {
							v287 = v287 + int32(48)
							v291 = v291 + v386
							v314 = v616
							continue
						} else {
							break
						}
						break
					}
					if v268 != v282 {
						v622 = v242 + v282
						v624 = v239 + v282*int32(3)
						if (l2-v622)&int32(1) != 0 {
							v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
							v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
							v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
							v665 = int32(base.Ui32(v651*int32(_a_F_ConvertBGR24ToY_SSE41_0)+v654*int32(_a_F_ConvertBGR24ToY_SSE41_1)+v658*int32(_a_F_ConvertBGR24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
							v671 = v622 + int32(1)
							v672 = v624 + int32(3)
						} else {
							v671 = v622
							v672 = v624
						}
						if v622 == l2+int32(-1) {
						} else {
							v678 = l1 + v671
							v681 = l2 - v671
							v683 = v672
							for {
								v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
								v707 = int32(_a_F_ConvertBGR24ToY_SSE41_0)
								v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
								v710 = int32(_a_F_ConvertBGR24ToY_SSE41_1)
								v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
								v714 = int32(_a_F_ConvertBGR24ToY_SSE41_2)
								v717 = int32(1081344)
								v719 = int32(16)
								v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
								v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
								v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
								v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
								v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
								v745 = v681 + int32(-2)
								if v745 != 0 {
									v678 = v678 + int32(2)
									v681 = v745
									v683 = v683 + int32(6)
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
					}
				}
			}
		}
	}
	return
}

var F_ConvertBGR24ToY_SSE41__k0 = [2]uint64{0x603008f8f8f8f8f, 0x8f8f8f8f8f0f0c09}
var F_ConvertBGR24ToY_SSE41__k1 = [2]uint64{0x8f8f8f0d0a070401, 0x8f8f8f8f8f8f8f8f}
var F_ConvertBGR24ToY_SSE41__k2 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xe0b0805028f8f8f}
var F_ConvertBGR24ToY_SSE41__k3 = [2]uint64{0x0, 0x0}
var F_ConvertBGR24ToY_SSE41__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertBGR24ToY_SSE41__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertBGR24ToY_SSE41__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertBGR24ToY_SSE41__k7 = [2]uint64{0x5028f8f8f8f8f8f, 0x8f8f8f8f8f0e0b08}
var F_ConvertBGR24ToY_SSE41__k8 = [2]uint64{0x8f8f0f0c09060300, 0x8f8f8f8f8f8f8f8f}
var F_ConvertBGR24ToY_SSE41__k9 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xd0a0704018f8f8f}
var F_ConvertBGR24ToY_SSE41__k10 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertBGR24ToY_SSE41__k11 = [2]uint64{0x1914400019144000, 0x1914400019144000}
var F_ConvertBGR24ToY_SSE41__k12 = [2]uint64{0x704018f8f8f8f8f, 0x8f8f8f8f8f8f0d0a}
var F_ConvertBGR24ToY_SSE41__k13 = [2]uint64{0x8f8f8f0e0b080502, 0x8f8f8f8f8f8f8f8f}
var F_ConvertBGR24ToY_SSE41__k14 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xf0c090603008f8f}
var F_ConvertBGR24ToY_SSE41__k15 = [2]uint64{0x412341c7412341c7, 0x412341c7412341c7}
var F_ConvertBGR24ToY_SSE41__k16 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertBGR24ToY_SSE41__k17 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertBGR24ToY_SSE41__k18 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertBGR24ToY_SSE41__k19 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertBGR24ToY_SSE41__k20 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertBGR24ToY_SSE41__k21 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertBGR24ToY_SSE41__k22 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertBGR24ToY_SSE41__k23 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertBGR24ToY_SSE41__k24 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertBGR24ToY_SSE41__k25 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertBGR24ToY_SSE41__k26 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertBGR24ToY_SSE41__k27 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertBGR24ToY_SSE41__k28 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_ConvertBGRAToBGR_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v28 int32
	_ = v28
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v34 base.V128
	_ = v34
	var v38 base.V128
	_ = v38
	var v39 base.V128
	_ = v39
	var v45 base.V128
	_ = v45
	var v50 base.V128
	_ = v50
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	if int32(9) <= l1 {
		v16 = l0
		v17 = l1
		v18 = l2
		for {
			v26 = int32(0)
			v27 = base.Simd_g_v128_load(m, v16, v26)
			v28 = int32(8)
			v30 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k0)
			v32 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k1)
			v34 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i64x2_shr_u(v27, v28), v30), base.Simd_g_v128_and(v27, v32))
			base.Simd_g_v128_store64_lane_l0(m, v18, v26, v34)
			v38 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k2)
			v39 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k3)
			base.Simd_g_v128_store64_lane_l0(m, v18, int32(6), base.Simd_g_i8x16_shuffle2(v34, v38, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k5)))
			v45 = base.Simd_g_v128_load(m, v16, int32(16))
			v50 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i64x2_shr_u(v45, v28), v30), base.Simd_g_v128_and(v45, v32))
			base.Simd_g_v128_store64_lane_l0(m, v18, int32(12), v50)
			base.Simd_g_v128_store64_lane_l0(m, v18, int32(18), base.Simd_g_i8x16_shuffle2(v50, v38, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToBGR_SSE2__k5)))
			v60 = v16 + int32(32)
			v62 = v17 + int32(-8)
			v66 = v18 + int32(24)
			if base.Ui32(v18+int32(50)) <= base.Ui32(l2+l1*int32(3)) {
				v16 = v60
				v17 = v62
				v18 = v66
				continue
			} else {
				break
			}
			break
		}
		v68 = v60
		v69 = v62
		v71 = v66
	} else {
		v68 = l0
		v69 = l1
		v71 = l2
	}
	if v69 < int32(1) {
	} else {
		if v69 < int32(1) {
		} else {
			v86 = v68
			v88 = v71
			for {
				v90 = *(*int32)(unsafe.Add(mBase, uint32(v86)))
				*(*uint8)(unsafe.Add(mBase, uint32(v88))) = uint8(v90)
				v93 = int32(base.Ui32(v90) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+2)) = uint8(v93)
				v96 = int32(base.Ui32(v90) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+1)) = uint8(v96)
				v101 = v86 + int32(4)
				if base.Ui32(v101) < base.Ui32(v68+v69<<(uint(int32(2))%32)) {
					v86 = v101
					v88 = v88 + int32(3)
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

var F_ConvertBGRAToBGR_SSE2__k0 = [2]uint64{0xffffff000000, 0xffffff000000}
var F_ConvertBGRAToBGR_SSE2__k1 = [2]uint64{0xffffff, 0xffffff}
var F_ConvertBGRAToBGR_SSE2__k2 = [2]uint64{0x0, 0x0}
var F_ConvertBGRAToBGR_SSE2__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x1716151413121110}
var F_ConvertBGRAToBGR_SSE2__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertBGRAToBGR_SSE2__k5 = [2]uint64{0x8080808080808080, 0x706050403020100}

func F_ConvertBGRAToBGR_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v13 base.V128
	_ = v13
	var v14 base.V128
	_ = v14
	var v17 base.V128
	_ = v17
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v37 base.V128
	_ = v37
	var v38 int32
	_ = v38
	var v39 base.V128
	_ = v39
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v52 int32
	_ = v52
	var v53 base.V128
	_ = v53
	var v57 base.V128
	_ = v57
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	if int32(16) <= l1 {
		v13 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k0)
		v14 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k1)
		v17 = base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k2)
		v25 = l0
		v26 = l1
		v27 = l2
		for {
			v35 = int32(32)
			v36 = base.Simd_g_v128_load_rng(m, v25, v35, int32(32), int32(32))
			v37 = base.Simd_g_i8x16_swizzle(v36, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(v13, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k3)), v17))
			v38 = int32(48)
			v39 = base.Simd_g_v128_load_nc(m, v25, v38)
			v42 = base.Simd_g_i8x16_shuffle2(v37, base.Simd_g_i8x16_swizzle(v39, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(v13, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k4)), v17)), base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k5), base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k6))
			base.Simd_g_v128_store(m, v27, v35, v42)
			v45 = int32(16)
			v46 = base.Simd_g_v128_load(m, v25, v45)
			v47 = base.Simd_g_i8x16_swizzle(v46, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(v13, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k7)), v17))
			v49 = base.Simd_g_i8x16_shuffle2(v47, v37, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k8), base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k9))
			base.Simd_g_v128_store(m, v27, v45, v49)
			v52 = int32(0)
			v53 = base.Simd_g_v128_load(m, v25, v52)
			v57 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_swizzle(v53, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k10)), v47, base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k11), base.Simd_g_const(&F_ConvertBGRAToBGR_SSE41__k12))
			base.Simd_g_v128_store(m, v27, v52, v57)
			v61 = v27 + v38
			v63 = v25 + int32(64)
			v67 = v26 + int32(-16)
			if base.Ui32(int32(31)) < base.Ui32(v26) {
				v25 = v63
				v26 = v67
				v27 = v61
				continue
			} else {
				break
			}
			break
		}
		v68 = v63
		v70 = v61
		v71 = v67
	} else {
		v68 = l0
		v70 = l2
		v71 = l1
	}
	if v71 < int32(1) {
	} else {
		if v71 < int32(1) {
		} else {
			v86 = v68
			v88 = v70
			for {
				v90 = *(*int32)(unsafe.Add(mBase, uint32(v86)))
				*(*uint8)(unsafe.Add(mBase, uint32(v88))) = uint8(v90)
				v93 = int32(base.Ui32(v90) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+2)) = uint8(v93)
				v96 = int32(base.Ui32(v90) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+1)) = uint8(v96)
				v101 = v86 + int32(4)
				if base.Ui32(v101) < base.Ui32(v68+v71<<(uint(int32(2))%32)) {
					v86 = v101
					v88 = v88 + int32(3)
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

var F_ConvertBGRAToBGR_SSE41__k0 = [2]uint64{0x908060504020100, 0xffffffff0e0d0c0a}
var F_ConvertBGRAToBGR_SSE41__k1 = [2]uint64{0x0, 0x0}
var F_ConvertBGRAToBGR_SSE41__k2 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0x8f8f8f8f8f8f8f8f}
var F_ConvertBGRAToBGR_SSE41__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}
var F_ConvertBGRAToBGR_SSE41__k4 = [2]uint64{0x30201000f0e0d0c, 0xb0a090807060504}
var F_ConvertBGRAToBGR_SSE41__k5 = [2]uint64{0x8080808003020100, 0x8080808080808080}
var F_ConvertBGRAToBGR_SSE41__k6 = [2]uint64{0x706050480808080, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToBGR_SSE41__k7 = [2]uint64{0xb0a090807060504, 0x30201000f0e0d0c}
var F_ConvertBGRAToBGR_SSE41__k8 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertBGRAToBGR_SSE41__k9 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToBGR_SSE41__k10 = [2]uint64{0x908060504020100, 0x8f8f8f8f0e0d0c0a}
var F_ConvertBGRAToBGR_SSE41__k11 = [2]uint64{0x706050403020100, 0x808080800b0a0908}
var F_ConvertBGRAToBGR_SSE41__k12 = [2]uint64{0x8080808080808080, 0xf0e0d0c80808080}

func F_ConvertBGRAToRGB565_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v20 base.V128
	_ = v20
	var v21 int32
	_ = v21
	var v22 base.V128
	_ = v22
	var v23 base.V128
	_ = v23
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v28 base.V128
	_ = v28
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v34 base.V128
	_ = v34
	var v36 base.V128
	_ = v36
	var v44 base.V128
	_ = v44
	var v46 int32
	_ = v46
	var v57 base.V128
	_ = v57
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v141 int32
	_ = v141
	var v142 base.V128
	_ = v142
	var v143 int32
	_ = v143
	var v151 base.V128
	_ = v151
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v170 int32
	_ = v170
	var v178 base.V128
	_ = v178
	var v180 int32
	_ = v180
	var v185 int32
	_ = v185
	var v190 int32
	_ = v190
	var v195 int32
	_ = v195
	var v202 int32
	_ = v202
	var v206 int32
	_ = v206
	var v209 int32
	_ = v209
	var v216 int32
	_ = v216
	var v219 int32
	_ = v219
	var v225 int32
	_ = v225
	var v234 int32
	_ = v234
	var v244 int32
	_ = v244
	var v249 int32
	_ = v249
	if int32(8) <= l1 {
		v11 = l0
		v12 = l1
		v13 = l2
		for {
			v19 = int32(0)
			v20 = base.Simd_g_v128_load_rng(m, v11, v19, int32(0), int32(32))
			v21 = int32(16)
			v22 = base.Simd_g_v128_load_nc(m, v11, v21)
			v23 = base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k0)
			v24 = base.Simd_g_i8x16_shuffle2(v20, v22, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k2))
			v25 = base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k3)
			v26 = base.Simd_g_i8x16_shuffle2(v20, v22, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k5))
			v28 = base.Simd_g_i8x16_shuffle2(v24, v26, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k2))
			v30 = base.Simd_g_i8x16_shuffle2(v24, v26, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k5))
			v32 = base.Simd_g_i8x16_shuffle2(v28, v30, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k2))
			v34 = base.Simd_g_i8x16_shuffle2(v28, v30, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k5))
			v36 = base.Simd_g_i8x16_shuffle2(v32, v34, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k6), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k7))
			v44 = base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v34, v32, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k8), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k9)), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k10))
			v46 = int32(3)
			v57 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(v36, int32(5)), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k11)), v44), base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shl(v36, v46), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k12)), base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v44, base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k13), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k6), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k9)), v46)), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k2))
			base.Simd_g_v128_store(m, v13, v19, v57)
			v61 = v11 + int32(32)
			v63 = v13 + v21
			v67 = v12 + int32(-8)
			if base.Ui32(int32(15)) < base.Ui32(v12) {
				v11 = v61
				v12 = v67
				v13 = v63
				continue
			} else {
				break
			}
			break
		}
		v68 = v61
		v70 = v63
		v71 = v67
	} else {
		v68 = l0
		v70 = l2
		v71 = l1
	}
	if v71 < int32(1) {
	} else {
		if v71 < int32(1) {
		} else {
			v89 = v68 + v71<<(uint(int32(2))%32)
			v91 = v68 + int32(4)
			if base.Ui32(v91) < base.Ui32(v89) {
				v93 = v89
			} else {
				v93 = v91
			}
			v95 = v68 ^ int32(-1)
			v96 = v93 + v95
			if base.Ui32(v96) < base.Ui32(int32(44)) {
				v206 = v70
				v209 = v68
				v216 = v206
				v219 = v209
				for {
					v225 = *(*int32)(unsafe.Add(mBase, uint32(v219)))
					v234 = int32(base.Ui32(v225)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v225)>>(uint(int32(3))%32))&int32(31)
					*(*uint8)(unsafe.Add(mBase, uint32(v216)+1)) = uint8(v234)
					v244 = int32(base.Ui32(v225)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v225)>>(uint(int32(13))%32))&int32(7)
					*(*uint8)(unsafe.Add(mBase, uint32(v216))) = uint8(v244)
					v249 = v219 + int32(4)
					if base.Ui32(v249) < base.Ui32(v89) {
						v216 = v216 + int32(2)
						v219 = v249
						continue
					} else {
						break
					}
					break
				}
			} else {
				v101 = v71<<(uint(int32(2))%32) + v68
				if base.Ui32(v91) < base.Ui32(v101) {
					v103 = v101
				} else {
					v103 = v91
				}
				v104 = v103 + v95
				if base.Ui32(v91+v104&int32(-4)) <= base.Ui32(v70) {
					v117 = int32(2)
					v119 = int32(1)
					v120 = int32(base.Ui32(v96)>>(uint(v117)%32)) + v119
					v122 = v120 & int32(2147483644)
					v129 = v68
					v131 = v70
					v134 = v122
					for {
						v141 = int32(0)
						v142 = base.Simd_g_v128_load(m, v129, v141)
						v143 = int32(16)
						v151 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, v143), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k14)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, int32(13)), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k15)))
						v152 = int32(3)
						v153 = base.Simd_g_i32x4_extract_lane_l3(v151)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+int32(6)))) = uint8(v153)
						v157 = int32(2)
						v158 = base.Simd_g_i32x4_extract_lane_l2(v151)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+int32(4)))) = uint8(v158)
						v162 = int32(1)
						v163 = base.Simd_g_i32x4_extract_lane_l1(v151)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+v157))) = uint8(v163)
						v166 = base.Simd_g_i32x4_extract_lane_l0(v151)
						*(*uint8)(unsafe.Add(mBase, uint32(v131))) = uint8(v166)
						v170 = int32(5)
						v178 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, v170), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k16)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, v152), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k17)))
						v180 = base.Simd_g_i32x4_extract_lane_l3(v178)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+int32(7)))) = uint8(v180)
						v185 = base.Simd_g_i32x4_extract_lane_l2(v178)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+v170))) = uint8(v185)
						v190 = base.Simd_g_i32x4_extract_lane_l1(v178)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+v152))) = uint8(v190)
						v195 = base.Simd_g_i32x4_extract_lane_l0(v178)
						*(*uint8)(unsafe.Add(mBase, uint32(v131+v162))) = uint8(v195)
						v202 = v134 + int32(-4)
						if v202 != 0 {
							v129 = v129 + v143
							v131 = v131 + int32(8)
							v134 = v202
							continue
						} else {
							break
						}
						break
					}
					if v120 == v122 {
					} else {
						v206 = v70 + v122<<(uint(v119)%32)
						v209 = v68 + v122<<(uint(v117)%32)
						v216 = v206
						v219 = v209
						for {
							v225 = *(*int32)(unsafe.Add(mBase, uint32(v219)))
							v234 = int32(base.Ui32(v225)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v225)>>(uint(int32(3))%32))&int32(31)
							*(*uint8)(unsafe.Add(mBase, uint32(v216)+1)) = uint8(v234)
							v244 = int32(base.Ui32(v225)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v225)>>(uint(int32(13))%32))&int32(7)
							*(*uint8)(unsafe.Add(mBase, uint32(v216))) = uint8(v244)
							v249 = v219 + int32(4)
							if base.Ui32(v249) < base.Ui32(v89) {
								v216 = v216 + int32(2)
								v219 = v249
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					if base.Ui32(v68) < base.Ui32(v70+int32(base.Ui32(v104)>>(uint(int32(1))%32))&int32(2147483646)+int32(2)) {
						v206 = v70
						v209 = v68
						v216 = v206
						v219 = v209
						for {
							v225 = *(*int32)(unsafe.Add(mBase, uint32(v219)))
							v234 = int32(base.Ui32(v225)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v225)>>(uint(int32(3))%32))&int32(31)
							*(*uint8)(unsafe.Add(mBase, uint32(v216)+1)) = uint8(v234)
							v244 = int32(base.Ui32(v225)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v225)>>(uint(int32(13))%32))&int32(7)
							*(*uint8)(unsafe.Add(mBase, uint32(v216))) = uint8(v244)
							v249 = v219 + int32(4)
							if base.Ui32(v249) < base.Ui32(v89) {
								v216 = v216 + int32(2)
								v219 = v249
								continue
							} else {
								break
							}
							break
						}
					} else {
						v117 = int32(2)
						v119 = int32(1)
						v120 = int32(base.Ui32(v96)>>(uint(v117)%32)) + v119
						v122 = v120 & int32(2147483644)
						v129 = v68
						v131 = v70
						v134 = v122
						for {
							v141 = int32(0)
							v142 = base.Simd_g_v128_load(m, v129, v141)
							v143 = int32(16)
							v151 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, v143), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k14)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, int32(13)), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k15)))
							v152 = int32(3)
							v153 = base.Simd_g_i32x4_extract_lane_l3(v151)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+int32(6)))) = uint8(v153)
							v157 = int32(2)
							v158 = base.Simd_g_i32x4_extract_lane_l2(v151)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+int32(4)))) = uint8(v158)
							v162 = int32(1)
							v163 = base.Simd_g_i32x4_extract_lane_l1(v151)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+v157))) = uint8(v163)
							v166 = base.Simd_g_i32x4_extract_lane_l0(v151)
							*(*uint8)(unsafe.Add(mBase, uint32(v131))) = uint8(v166)
							v170 = int32(5)
							v178 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, v170), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k16)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v142, v152), base.Simd_g_const(&F_ConvertBGRAToRGB565_SSE2__k17)))
							v180 = base.Simd_g_i32x4_extract_lane_l3(v178)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+int32(7)))) = uint8(v180)
							v185 = base.Simd_g_i32x4_extract_lane_l2(v178)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+v170))) = uint8(v185)
							v190 = base.Simd_g_i32x4_extract_lane_l1(v178)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+v152))) = uint8(v190)
							v195 = base.Simd_g_i32x4_extract_lane_l0(v178)
							*(*uint8)(unsafe.Add(mBase, uint32(v131+v162))) = uint8(v195)
							v202 = v134 + int32(-4)
							if v202 != 0 {
								v129 = v129 + v143
								v131 = v131 + int32(8)
								v134 = v202
								continue
							} else {
								break
							}
							break
						}
						if v120 == v122 {
						} else {
							v206 = v70 + v122<<(uint(v119)%32)
							v209 = v68 + v122<<(uint(v117)%32)
							v216 = v206
							v219 = v209
							for {
								v225 = *(*int32)(unsafe.Add(mBase, uint32(v219)))
								v234 = int32(base.Ui32(v225)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v225)>>(uint(int32(3))%32))&int32(31)
								*(*uint8)(unsafe.Add(mBase, uint32(v216)+1)) = uint8(v234)
								v244 = int32(base.Ui32(v225)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v225)>>(uint(int32(13))%32))&int32(7)
								*(*uint8)(unsafe.Add(mBase, uint32(v216))) = uint8(v244)
								v249 = v219 + int32(4)
								if base.Ui32(v249) < base.Ui32(v89) {
									v216 = v216 + int32(2)
									v219 = v249
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

var F_ConvertBGRAToRGB565_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertBGRAToRGB565_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertBGRAToRGB565_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertBGRAToRGB565_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertBGRAToRGB565_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertBGRAToRGB565_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertBGRAToRGB565_SSE2__k6 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertBGRAToRGB565_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToRGB565_SSE2__k8 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertBGRAToRGB565_SSE2__k9 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertBGRAToRGB565_SSE2__k10 = [2]uint64{0xf8f8f8f8f8f8f8f8, 0xf8f8f8f8f8f8f8f8}
var F_ConvertBGRAToRGB565_SSE2__k11 = [2]uint64{0x707070707070707, 0x707070707070707}
var F_ConvertBGRAToRGB565_SSE2__k12 = [2]uint64{0xe0e0e0e0e0e0e0e0, 0xe0e0e0e0e0e0e0e0}
var F_ConvertBGRAToRGB565_SSE2__k13 = [2]uint64{0x0, 0x0}
var F_ConvertBGRAToRGB565_SSE2__k14 = [2]uint64{0xf8000000f8, 0xf8000000f8}
var F_ConvertBGRAToRGB565_SSE2__k15 = [2]uint64{0x700000007, 0x700000007}
var F_ConvertBGRAToRGB565_SSE2__k16 = [2]uint64{0xe0000000e0, 0xe0000000e0}
var F_ConvertBGRAToRGB565_SSE2__k17 = [2]uint64{0x1f0000001f, 0x1f0000001f}

func F_ConvertBGRAToRGBA4444_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v20 base.V128
	_ = v20
	var v21 int32
	_ = v21
	var v22 base.V128
	_ = v22
	var v23 base.V128
	_ = v23
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v28 base.V128
	_ = v28
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v34 base.V128
	_ = v34
	var v45 base.V128
	_ = v45
	var v50 base.V128
	_ = v50
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v136 int32
	_ = v136
	var v137 base.V128
	_ = v137
	var v138 base.V128
	_ = v138
	var v142 base.V128
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v159 int32
	_ = v159
	var v163 int32
	_ = v163
	var v170 base.V128
	_ = v170
	var v172 int32
	_ = v172
	var v177 int32
	_ = v177
	var v182 int32
	_ = v182
	var v185 int32
	_ = v185
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v199 int32
	_ = v199
	var v207 int32
	_ = v207
	var v210 int32
	_ = v210
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v222 int32
	_ = v222
	var v232 int32
	_ = v232
	var v237 int32
	_ = v237
	if int32(8) <= l1 {
		v11 = l0
		v12 = l1
		v13 = l2
		for {
			v19 = int32(0)
			v20 = base.Simd_g_v128_load_rng(m, v11, v19, int32(0), int32(32))
			v21 = int32(16)
			v22 = base.Simd_g_v128_load_nc(m, v11, v21)
			v23 = base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k0)
			v24 = base.Simd_g_i8x16_shuffle2(v20, v22, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k2))
			v25 = base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k3)
			v26 = base.Simd_g_i8x16_shuffle2(v20, v22, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k5))
			v28 = base.Simd_g_i8x16_shuffle2(v24, v26, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k2))
			v30 = base.Simd_g_i8x16_shuffle2(v24, v26, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k5))
			v32 = base.Simd_g_i8x16_shuffle2(v28, v30, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k2))
			v34 = base.Simd_g_i8x16_shuffle2(v28, v30, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k5))
			v45 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_i8x16_shuffle2(v32, v34, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k6), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k7)), int32(4)), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k8)), base.Simd_g_v128_and(base.Simd_g_i8x16_shuffle2(v34, v32, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k9), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k10)), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k11)))
			v50 = base.Simd_g_i8x16_shuffle2(v45, base.Simd_g_i8x16_shuffle2(v45, base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k12), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k6), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k10)), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k2))
			base.Simd_g_v128_store(m, v13, v19, v50)
			v54 = v11 + int32(32)
			v56 = v13 + v21
			v60 = v12 + int32(-8)
			if base.Ui32(int32(15)) < base.Ui32(v12) {
				v11 = v54
				v12 = v60
				v13 = v56
				continue
			} else {
				break
			}
			break
		}
		v61 = v54
		v63 = v56
		v64 = v60
	} else {
		v61 = l0
		v63 = l2
		v64 = l1
	}
	if v64 < int32(1) {
	} else {
		if v64 < int32(1) {
		} else {
			v83 = v61 + v64<<(uint(int32(2))%32)
			v85 = v61 + int32(4)
			if base.Ui32(v85) < base.Ui32(v83) {
				v87 = v83
			} else {
				v87 = v85
			}
			v89 = v61 ^ int32(-1)
			v90 = v87 + v89
			if base.Ui32(v90) < base.Ui32(int32(44)) {
				v196 = v63
				v199 = v61
				v207 = v196
				v210 = v199
				for {
					v217 = *(*int32)(unsafe.Add(mBase, uint32(v210)))
					v218 = int32(240)
					v222 = v217&v218 | int32(base.Ui32(v217)>>(uint(int32(28))%32))
					*(*uint8)(unsafe.Add(mBase, uint32(v207)+1)) = uint8(v222)
					v232 = int32(base.Ui32(v217)>>(uint(int32(16))%32))&v218 | int32(base.Ui32(v217)>>(uint(int32(12))%32))&int32(15)
					*(*uint8)(unsafe.Add(mBase, uint32(v207))) = uint8(v232)
					v237 = v210 + int32(4)
					if base.Ui32(v237) < base.Ui32(v83) {
						v207 = v207 + int32(2)
						v210 = v237
						continue
					} else {
						break
					}
					break
				}
			} else {
				v95 = v64<<(uint(int32(2))%32) + v61
				if base.Ui32(v85) < base.Ui32(v95) {
					v97 = v95
				} else {
					v97 = v85
				}
				v98 = v97 + v89
				if base.Ui32(v85+v98&int32(-4)) <= base.Ui32(v63) {
					v111 = int32(2)
					v113 = int32(1)
					v114 = int32(base.Ui32(v90)>>(uint(v111)%32)) + v113
					v116 = v114 & int32(2147483644)
					v123 = v61
					v125 = v63
					v128 = v116
					for {
						v136 = int32(0)
						v137 = base.Simd_g_v128_load(m, v123, v136)
						v138 = base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k13)
						v142 = base.Simd_g_v128_or(base.Simd_g_v128_and(v137, v138), base.Simd_g_i32x4_shr_u(v137, int32(28)))
						v143 = int32(3)
						v144 = base.Simd_g_i32x4_extract_lane_l3(v142)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(7)))) = uint8(v144)
						v148 = int32(2)
						v149 = base.Simd_g_i32x4_extract_lane_l2(v142)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(5)))) = uint8(v149)
						v153 = int32(1)
						v154 = base.Simd_g_i32x4_extract_lane_l1(v142)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+v143))) = uint8(v154)
						v159 = base.Simd_g_i32x4_extract_lane_l0(v142)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+v153))) = uint8(v159)
						v163 = int32(16)
						v170 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v137, v163), v138), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v137, int32(12)), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k14)))
						v172 = base.Simd_g_i32x4_extract_lane_l3(v170)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(6)))) = uint8(v172)
						v177 = base.Simd_g_i32x4_extract_lane_l2(v170)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(4)))) = uint8(v177)
						v182 = base.Simd_g_i32x4_extract_lane_l1(v170)
						*(*uint8)(unsafe.Add(mBase, uint32(v125+v148))) = uint8(v182)
						v185 = base.Simd_g_i32x4_extract_lane_l0(v170)
						*(*uint8)(unsafe.Add(mBase, uint32(v125))) = uint8(v185)
						v192 = v128 + int32(-4)
						if v192 != 0 {
							v123 = v123 + v163
							v125 = v125 + int32(8)
							v128 = v192
							continue
						} else {
							break
						}
						break
					}
					if v114 == v116 {
					} else {
						v196 = v63 + v116<<(uint(v113)%32)
						v199 = v61 + v116<<(uint(v111)%32)
						v207 = v196
						v210 = v199
						for {
							v217 = *(*int32)(unsafe.Add(mBase, uint32(v210)))
							v218 = int32(240)
							v222 = v217&v218 | int32(base.Ui32(v217)>>(uint(int32(28))%32))
							*(*uint8)(unsafe.Add(mBase, uint32(v207)+1)) = uint8(v222)
							v232 = int32(base.Ui32(v217)>>(uint(int32(16))%32))&v218 | int32(base.Ui32(v217)>>(uint(int32(12))%32))&int32(15)
							*(*uint8)(unsafe.Add(mBase, uint32(v207))) = uint8(v232)
							v237 = v210 + int32(4)
							if base.Ui32(v237) < base.Ui32(v83) {
								v207 = v207 + int32(2)
								v210 = v237
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					if base.Ui32(v61) < base.Ui32(v63+int32(base.Ui32(v98)>>(uint(int32(1))%32))&int32(2147483646)+int32(2)) {
						v196 = v63
						v199 = v61
						v207 = v196
						v210 = v199
						for {
							v217 = *(*int32)(unsafe.Add(mBase, uint32(v210)))
							v218 = int32(240)
							v222 = v217&v218 | int32(base.Ui32(v217)>>(uint(int32(28))%32))
							*(*uint8)(unsafe.Add(mBase, uint32(v207)+1)) = uint8(v222)
							v232 = int32(base.Ui32(v217)>>(uint(int32(16))%32))&v218 | int32(base.Ui32(v217)>>(uint(int32(12))%32))&int32(15)
							*(*uint8)(unsafe.Add(mBase, uint32(v207))) = uint8(v232)
							v237 = v210 + int32(4)
							if base.Ui32(v237) < base.Ui32(v83) {
								v207 = v207 + int32(2)
								v210 = v237
								continue
							} else {
								break
							}
							break
						}
					} else {
						v111 = int32(2)
						v113 = int32(1)
						v114 = int32(base.Ui32(v90)>>(uint(v111)%32)) + v113
						v116 = v114 & int32(2147483644)
						v123 = v61
						v125 = v63
						v128 = v116
						for {
							v136 = int32(0)
							v137 = base.Simd_g_v128_load(m, v123, v136)
							v138 = base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k13)
							v142 = base.Simd_g_v128_or(base.Simd_g_v128_and(v137, v138), base.Simd_g_i32x4_shr_u(v137, int32(28)))
							v143 = int32(3)
							v144 = base.Simd_g_i32x4_extract_lane_l3(v142)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(7)))) = uint8(v144)
							v148 = int32(2)
							v149 = base.Simd_g_i32x4_extract_lane_l2(v142)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(5)))) = uint8(v149)
							v153 = int32(1)
							v154 = base.Simd_g_i32x4_extract_lane_l1(v142)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+v143))) = uint8(v154)
							v159 = base.Simd_g_i32x4_extract_lane_l0(v142)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+v153))) = uint8(v159)
							v163 = int32(16)
							v170 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v137, v163), v138), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v137, int32(12)), base.Simd_g_const(&F_ConvertBGRAToRGBA4444_SSE2__k14)))
							v172 = base.Simd_g_i32x4_extract_lane_l3(v170)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(6)))) = uint8(v172)
							v177 = base.Simd_g_i32x4_extract_lane_l2(v170)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+int32(4)))) = uint8(v177)
							v182 = base.Simd_g_i32x4_extract_lane_l1(v170)
							*(*uint8)(unsafe.Add(mBase, uint32(v125+v148))) = uint8(v182)
							v185 = base.Simd_g_i32x4_extract_lane_l0(v170)
							*(*uint8)(unsafe.Add(mBase, uint32(v125))) = uint8(v185)
							v192 = v128 + int32(-4)
							if v192 != 0 {
								v123 = v123 + v163
								v125 = v125 + int32(8)
								v128 = v192
								continue
							} else {
								break
							}
							break
						}
						if v114 == v116 {
						} else {
							v196 = v63 + v116<<(uint(v113)%32)
							v199 = v61 + v116<<(uint(v111)%32)
							v207 = v196
							v210 = v199
							for {
								v217 = *(*int32)(unsafe.Add(mBase, uint32(v210)))
								v218 = int32(240)
								v222 = v217&v218 | int32(base.Ui32(v217)>>(uint(int32(28))%32))
								*(*uint8)(unsafe.Add(mBase, uint32(v207)+1)) = uint8(v222)
								v232 = int32(base.Ui32(v217)>>(uint(int32(16))%32))&v218 | int32(base.Ui32(v217)>>(uint(int32(12))%32))&int32(15)
								*(*uint8)(unsafe.Add(mBase, uint32(v207))) = uint8(v232)
								v237 = v210 + int32(4)
								if base.Ui32(v237) < base.Ui32(v83) {
									v207 = v207 + int32(2)
									v210 = v237
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

var F_ConvertBGRAToRGBA4444_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertBGRAToRGBA4444_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertBGRAToRGBA4444_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertBGRAToRGBA4444_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertBGRAToRGBA4444_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertBGRAToRGBA4444_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertBGRAToRGBA4444_SSE2__k6 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertBGRAToRGBA4444_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToRGBA4444_SSE2__k8 = [2]uint64{0xf0f0f0f0f0f0f0f, 0xf0f0f0f0f0f0f0f}
var F_ConvertBGRAToRGBA4444_SSE2__k9 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertBGRAToRGBA4444_SSE2__k10 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertBGRAToRGBA4444_SSE2__k11 = [2]uint64{0xf0f0f0f0f0f0f0f0, 0xf0f0f0f0f0f0f0f0}
var F_ConvertBGRAToRGBA4444_SSE2__k12 = [2]uint64{0x0, 0x0}
var F_ConvertBGRAToRGBA4444_SSE2__k13 = [2]uint64{0xf0000000f0, 0xf0000000f0}
var F_ConvertBGRAToRGBA4444_SSE2__k14 = [2]uint64{0xf0000000f, 0xf0000000f}

func F_ConvertBGRAToRGBA_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v21 int32
	_ = v21
	var v22 base.V128
	_ = v22
	var v23 base.V128
	_ = v23
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v28 base.V128
	_ = v28
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v43 base.V128
	_ = v43
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	if int32(8) <= l1 {
		v12 = l0
		v13 = l1
		v14 = l2
		for {
			v21 = int32(16)
			v22 = base.Simd_g_v128_load(m, v12, v21)
			v23 = base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k0)
			v25 = base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k1)
			v27 = base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k2)
			v28 = base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k3)
			v30 = base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k4)
			v32 = base.Simd_g_v128_or(base.Simd_g_v128_and(v22, v23), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_and(v22, v25), base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k3)), base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k4)))
			base.Simd_g_v128_store(m, v14, v21, v32)
			v35 = int32(0)
			v36 = base.Simd_g_v128_load(m, v12, v35)
			v43 = base.Simd_g_v128_or(base.Simd_g_v128_and(v36, v23), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_and(v36, v25), base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k3)), base.Simd_g_const(&F_ConvertBGRAToRGBA_SSE2__k4)))
			base.Simd_g_v128_store(m, v14, v35, v43)
			v46 = int32(32)
			v47 = v12 + v46
			v49 = v14 + v46
			v53 = v13 + int32(-8)
			if base.Ui32(int32(15)) < base.Ui32(v13) {
				v12 = v47
				v13 = v53
				v14 = v49
				continue
			} else {
				break
			}
			break
		}
		v54 = v47
		v56 = v49
		v57 = v53
	} else {
		v54 = l0
		v56 = l2
		v57 = l1
	}
	if v57 < int32(1) {
	} else {
		if v57 < int32(1) {
		} else {
			v71 = v54
			v73 = v56
			for {
				v75 = *(*int32)(unsafe.Add(mBase, uint32(v71)))
				*(*uint8)(unsafe.Add(mBase, uint32(v73)+2)) = uint8(v75)
				v78 = int32(base.Ui32(v75) >> (uint(int32(24)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v73)+3)) = uint8(v78)
				v81 = int32(base.Ui32(v75) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v73)+1)) = uint8(v81)
				v84 = int32(base.Ui32(v75) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v73))) = uint8(v84)
				v86 = int32(4)
				v89 = v71 + v86
				if base.Ui32(v89) < base.Ui32(v54+v57<<(uint(int32(2))%32)) {
					v71 = v89
					v73 = v73 + v86
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

var F_ConvertBGRAToRGBA_SSE2__k0 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_ConvertBGRAToRGBA_SSE2__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_ConvertBGRAToRGBA_SSE2__k2 = [2]uint64{0x0, 0x0}
var F_ConvertBGRAToRGBA_SSE2__k3 = [2]uint64{0x504070601000302, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToRGBA_SSE2__k4 = [2]uint64{0x706050403020100, 0xd0c0f0e09080b0a}

func F_ConvertBGRAToRGB_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 int32
	_ = v47
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v58 base.V128
	_ = v58
	var v59 int32
	_ = v59
	var v60 base.V128
	_ = v60
	var v61 int32
	_ = v61
	var v62 base.V128
	_ = v62
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v75 base.V128
	_ = v75
	var v77 int32
	_ = v77
	var v78 base.V128
	_ = v78
	var v79 int32
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
	var v90 base.V128
	_ = v90
	var v91 int32
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v106 base.V128
	_ = v106
	var v108 base.V128
	_ = v108
	var v109 int32
	_ = v109
	var v116 base.V128
	_ = v116
	var v124 base.V128
	_ = v124
	var v127 base.V128
	_ = v127
	var v130 base.V128
	_ = v130
	var v132 base.V128
	_ = v132
	var v133 base.V128
	_ = v133
	var v137 base.V128
	_ = v137
	var v140 base.V128
	_ = v140
	var v147 base.V128
	_ = v147
	var v150 base.V128
	_ = v150
	var v152 base.V128
	_ = v152
	var v157 base.V128
	_ = v157
	var v161 base.V128
	_ = v161
	var v163 base.V128
	_ = v163
	var v168 base.V128
	_ = v168
	var v171 base.V128
	_ = v171
	var v174 base.V128
	_ = v174
	var v179 base.V128
	_ = v179
	var v186 base.V128
	_ = v186
	var v189 base.V128
	_ = v189
	var v196 base.V128
	_ = v196
	var v199 base.V128
	_ = v199
	var v202 base.V128
	_ = v202
	var v207 base.V128
	_ = v207
	var v212 base.V128
	_ = v212
	var v217 base.V128
	_ = v217
	var v220 base.V128
	_ = v220
	var v223 base.V128
	_ = v223
	var v228 base.V128
	_ = v228
	var v233 base.V128
	_ = v233
	var v236 base.V128
	_ = v236
	var v241 base.V128
	_ = v241
	var v246 base.V128
	_ = v246
	var v251 base.V128
	_ = v251
	var v255 int32
	_ = v255
	var v257 int32
	_ = v257
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v291 int32
	_ = v291
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v306 int32
	_ = v306
	if int32(32) <= l1 {
		v24 = l0
		v25 = l1
		v26 = l2
		for {
			v45 = int32(0)
			v46 = base.Simd_g_v128_load_rng(m, v24, v45, int32(0), int32(128))
			v47 = int32(16)
			v48 = base.Simd_g_v128_load_nc(m, v24, v47)
			v49 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k0)
			v50 = base.Simd_g_i8x16_shuffle2(v46, v48, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v51 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k3)
			v52 = base.Simd_g_i8x16_shuffle2(v46, v48, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v54 = base.Simd_g_i8x16_shuffle2(v50, v52, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v56 = base.Simd_g_i8x16_shuffle2(v50, v52, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v58 = base.Simd_g_i8x16_shuffle2(v54, v56, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v59 = int32(32)
			v60 = base.Simd_g_v128_load_nc(m, v24, v59)
			v61 = int32(48)
			v62 = base.Simd_g_v128_load_nc(m, v24, v61)
			v64 = base.Simd_g_i8x16_shuffle2(v60, v62, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v66 = base.Simd_g_i8x16_shuffle2(v60, v62, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v68 = base.Simd_g_i8x16_shuffle2(v64, v66, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v70 = base.Simd_g_i8x16_shuffle2(v64, v66, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v72 = base.Simd_g_i8x16_shuffle2(v68, v70, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v73 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k6)
			v74 = base.Simd_g_i8x16_shuffle2(v58, v72, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k7), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k8))
			v75 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k9)
			v77 = int32(64)
			v78 = base.Simd_g_v128_load_nc(m, v24, v77)
			v79 = int32(80)
			v80 = base.Simd_g_v128_load_nc(m, v24, v79)
			v82 = base.Simd_g_i8x16_shuffle2(v78, v80, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v84 = base.Simd_g_i8x16_shuffle2(v78, v80, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v86 = base.Simd_g_i8x16_shuffle2(v82, v84, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v88 = base.Simd_g_i8x16_shuffle2(v82, v84, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v90 = base.Simd_g_i8x16_shuffle2(v86, v88, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v91 = int32(96)
			v92 = base.Simd_g_v128_load_nc(m, v24, v91)
			v94 = base.Simd_g_v128_load_nc(m, v24, int32(112))
			v96 = base.Simd_g_i8x16_shuffle2(v92, v94, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v98 = base.Simd_g_i8x16_shuffle2(v92, v94, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v100 = base.Simd_g_i8x16_shuffle2(v96, v98, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v102 = base.Simd_g_i8x16_shuffle2(v96, v98, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5))
			v104 = base.Simd_g_i8x16_shuffle2(v100, v102, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k1), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k2))
			v106 = base.Simd_g_i8x16_shuffle2(v90, v104, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k7), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k8))
			v108 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v74, v75), base.Simd_g_v128_and(v106, v75))
			v109 = int32(8)
			v116 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v54, v56, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v68, v70, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5)), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k7), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k8))
			v124 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v86, v88, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v100, v102, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k4), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k5)), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k7), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k8))
			v127 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v116, v109), base.Simd_g_i16x8_shr_u(v124, v109))
			v130 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v108, v109), base.Simd_g_i16x8_shr_u(v127, v109))
			v132 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k10)
			v133 = base.Simd_g_i8x16_shuffle2(v58, v72, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k11), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k12))
			v137 = base.Simd_g_i8x16_shuffle2(v90, v104, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k11), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE2__k12))
			v140 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v133, v109), base.Simd_g_i16x8_shr_u(v137, v109))
			v147 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v74, v109), base.Simd_g_i16x8_shr_u(v106, v109))
			v150 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v140, v109), base.Simd_g_i16x8_shr_u(v147, v109))
			v152 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v130, v75), base.Simd_g_v128_and(v150, v75))
			v157 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v116, v75), base.Simd_g_v128_and(v124, v75))
			v161 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v133, v75), base.Simd_g_v128_and(v137, v75))
			v163 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v157, v75), base.Simd_g_v128_and(v161, v75))
			v168 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v108, v75), base.Simd_g_v128_and(v127, v75))
			v171 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v163, v109), base.Simd_g_i16x8_shr_u(v168, v109))
			v174 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v152, v109), base.Simd_g_i16x8_shr_u(v171, v109))
			v179 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v140, v75), base.Simd_g_v128_and(v147, v75))
			v186 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v157, v109), base.Simd_g_i16x8_shr_u(v161, v109))
			v189 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v179, v109), base.Simd_g_i16x8_shr_u(v186, v109))
			v196 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v130, v109), base.Simd_g_i16x8_shr_u(v150, v109))
			v199 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v189, v109), base.Simd_g_i16x8_shr_u(v196, v109))
			v202 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v174, v109), base.Simd_g_i16x8_shr_u(v199, v109))
			base.Simd_g_v128_store(m, v26, v79, v202)
			v207 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v189, v75), base.Simd_g_v128_and(v196, v75))
			v212 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v163, v75), base.Simd_g_v128_and(v168, v75))
			v217 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v179, v75), base.Simd_g_v128_and(v186, v75))
			v220 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v212, v109), base.Simd_g_i16x8_shr_u(v217, v109))
			v223 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v207, v109), base.Simd_g_i16x8_shr_u(v220, v109))
			base.Simd_g_v128_store(m, v26, v77, v223)
			v228 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v212, v75), base.Simd_g_v128_and(v217, v75))
			v233 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v152, v75), base.Simd_g_v128_and(v171, v75))
			v236 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v228, v109), base.Simd_g_i16x8_shr_u(v233, v109))
			base.Simd_g_v128_store(m, v26, v61, v236)
			v241 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v174, v75), base.Simd_g_v128_and(v199, v75))
			base.Simd_g_v128_store(m, v26, v59, v241)
			v246 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v207, v75), base.Simd_g_v128_and(v220, v75))
			base.Simd_g_v128_store(m, v26, v47, v246)
			v251 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v228, v75), base.Simd_g_v128_and(v233, v75))
			base.Simd_g_v128_store(m, v26, v45, v251)
			v255 = v26 + v91
			v257 = v24 + int32(128)
			v261 = v25 + int32(-32)
			if base.Ui32(int32(63)) < base.Ui32(v25) {
				v24 = v257
				v25 = v261
				v26 = v255
				continue
			} else {
				break
			}
			break
		}
		v262 = v257
		v264 = v255
		v265 = v261
	} else {
		v262 = l0
		v264 = l2
		v265 = l1
	}
	if v265 < int32(1) {
	} else {
		if v265 < int32(1) {
		} else {
			v291 = v262
			v293 = v264
			for {
				v295 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
				*(*uint8)(unsafe.Add(mBase, uint32(v293)+2)) = uint8(v295)
				v298 = int32(base.Ui32(v295) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v293)+1)) = uint8(v298)
				v301 = int32(base.Ui32(v295) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v293))) = uint8(v301)
				v306 = v291 + int32(4)
				if base.Ui32(v306) < base.Ui32(v262+v265<<(uint(int32(2))%32)) {
					v291 = v306
					v293 = v293 + int32(3)
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

var F_ConvertBGRAToRGB_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertBGRAToRGB_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertBGRAToRGB_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertBGRAToRGB_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertBGRAToRGB_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertBGRAToRGB_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertBGRAToRGB_SSE2__k6 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertBGRAToRGB_SSE2__k7 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertBGRAToRGB_SSE2__k8 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertBGRAToRGB_SSE2__k9 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_ConvertBGRAToRGB_SSE2__k10 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_ConvertBGRAToRGB_SSE2__k11 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertBGRAToRGB_SSE2__k12 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}

func F_ConvertBGRAToRGB_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v13 base.V128
	_ = v13
	var v14 base.V128
	_ = v14
	var v17 base.V128
	_ = v17
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v37 base.V128
	_ = v37
	var v38 int32
	_ = v38
	var v39 base.V128
	_ = v39
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v52 int32
	_ = v52
	var v53 base.V128
	_ = v53
	var v57 base.V128
	_ = v57
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	if int32(16) <= l1 {
		v13 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k0)
		v14 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k1)
		v17 = base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k2)
		v25 = l0
		v26 = l1
		v27 = l2
		for {
			v35 = int32(32)
			v36 = base.Simd_g_v128_load_rng(m, v25, v35, int32(32), int32(32))
			v37 = base.Simd_g_i8x16_swizzle(v36, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(v13, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k3)), v17))
			v38 = int32(48)
			v39 = base.Simd_g_v128_load_nc(m, v25, v38)
			v42 = base.Simd_g_i8x16_shuffle2(v37, base.Simd_g_i8x16_swizzle(v39, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(v13, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k4)), v17)), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k5), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k6))
			base.Simd_g_v128_store(m, v27, v35, v42)
			v45 = int32(16)
			v46 = base.Simd_g_v128_load(m, v25, v45)
			v47 = base.Simd_g_i8x16_swizzle(v46, base.Simd_g_v128_and(base.Simd_g_i8x16_swizzle_c(v13, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k7)), v17))
			v49 = base.Simd_g_i8x16_shuffle2(v47, v37, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k8), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k9))
			base.Simd_g_v128_store(m, v27, v45, v49)
			v52 = int32(0)
			v53 = base.Simd_g_v128_load(m, v25, v52)
			v57 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_swizzle(v53, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k10)), v47, base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k11), base.Simd_g_const(&F_ConvertBGRAToRGB_SSE41__k12))
			base.Simd_g_v128_store(m, v27, v52, v57)
			v61 = v27 + v38
			v63 = v25 + int32(64)
			v67 = v26 + int32(-16)
			if base.Ui32(int32(31)) < base.Ui32(v26) {
				v25 = v63
				v26 = v67
				v27 = v61
				continue
			} else {
				break
			}
			break
		}
		v68 = v63
		v70 = v61
		v71 = v67
	} else {
		v68 = l0
		v70 = l2
		v71 = l1
	}
	if v71 < int32(1) {
	} else {
		if v71 < int32(1) {
		} else {
			v86 = v68
			v88 = v70
			for {
				v90 = *(*int32)(unsafe.Add(mBase, uint32(v86)))
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+2)) = uint8(v90)
				v93 = int32(base.Ui32(v90) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v88)+1)) = uint8(v93)
				v96 = int32(base.Ui32(v90) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v88))) = uint8(v96)
				v101 = v86 + int32(4)
				if base.Ui32(v101) < base.Ui32(v68+v71<<(uint(int32(2))%32)) {
					v86 = v101
					v88 = v88 + int32(3)
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

var F_ConvertBGRAToRGB_SSE41__k0 = [2]uint64{0x90a040506000102, 0xffffffff0c0d0e08}
var F_ConvertBGRAToRGB_SSE41__k1 = [2]uint64{0x0, 0x0}
var F_ConvertBGRAToRGB_SSE41__k2 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0x8f8f8f8f8f8f8f8f}
var F_ConvertBGRAToRGB_SSE41__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}
var F_ConvertBGRAToRGB_SSE41__k4 = [2]uint64{0x30201000f0e0d0c, 0xb0a090807060504}
var F_ConvertBGRAToRGB_SSE41__k5 = [2]uint64{0x8080808003020100, 0x8080808080808080}
var F_ConvertBGRAToRGB_SSE41__k6 = [2]uint64{0x706050480808080, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToRGB_SSE41__k7 = [2]uint64{0xb0a090807060504, 0x30201000f0e0d0c}
var F_ConvertBGRAToRGB_SSE41__k8 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertBGRAToRGB_SSE41__k9 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertBGRAToRGB_SSE41__k10 = [2]uint64{0x90a040506000102, 0x8f8f8f8f0c0d0e08}
var F_ConvertBGRAToRGB_SSE41__k11 = [2]uint64{0x706050403020100, 0x808080800b0a0908}
var F_ConvertBGRAToRGB_SSE41__k12 = [2]uint64{0x8080808080808080, 0xf0e0d0c80808080}

func F_ConvertRGB24ToY_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v47 int32
	_ = v47
	var v48 base.V128
	_ = v48
	var v51 base.V128
	_ = v51
	var v54 base.V128
	_ = v54
	var v65 base.V128
	_ = v65
	var v68 base.V128
	_ = v68
	var v71 base.V128
	_ = v71
	var v74 base.V128
	_ = v74
	var v77 base.V128
	_ = v77
	var v89 base.V128
	_ = v89
	var v92 base.V128
	_ = v92
	var v95 base.V128
	_ = v95
	var v98 base.V128
	_ = v98
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v106 int32
	_ = v106
	var v108 base.V128
	_ = v108
	var v119 base.V128
	_ = v119
	var v122 base.V128
	_ = v122
	var v125 base.V128
	_ = v125
	var v128 base.V128
	_ = v128
	var v141 base.V128
	_ = v141
	var v144 base.V128
	_ = v144
	var v147 base.V128
	_ = v147
	var v150 base.V128
	_ = v150
	var v164 base.V128
	_ = v164
	var v167 base.V128
	_ = v167
	var v170 base.V128
	_ = v170
	var v173 base.V128
	_ = v173
	var v192 base.V128
	_ = v192
	var v195 base.V128
	_ = v195
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v214 base.V128
	_ = v214
	var v217 base.V128
	_ = v217
	var v220 base.V128
	_ = v220
	var v223 base.V128
	_ = v223
	var v237 base.V128
	_ = v237
	var v240 base.V128
	_ = v240
	var v243 base.V128
	_ = v243
	var v246 base.V128
	_ = v246
	var v264 base.V128
	_ = v264
	var v267 base.V128
	_ = v267
	var v270 base.V128
	_ = v270
	var v273 base.V128
	_ = v273
	var v286 base.V128
	_ = v286
	var v289 base.V128
	_ = v289
	var v292 base.V128
	_ = v292
	var v295 base.V128
	_ = v295
	var v309 base.V128
	_ = v309
	var v312 base.V128
	_ = v312
	var v315 base.V128
	_ = v315
	var v318 base.V128
	_ = v318
	var v328 base.V128
	_ = v328
	var v336 int32
	_ = v336
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v352 int32
	_ = v352
	var v356 int32
	_ = v356
	var v358 int32
	_ = v358
	var v364 int32
	_ = v364
	var v367 int32
	_ = v367
	var v371 int32
	_ = v371
	var v378 int32
	_ = v378
	var v385 int32
	_ = v385
	if l2 < int32(1) {
	} else {
		if base.Ui32(int32(16)) <= base.Ui32(l2) {
			v19 = l2 & int32(2147483632)
			v23 = l0
			v28 = v19
			v29 = l1
			for {
				v39 = int32(3)
				v41 = int32(0)
				v42 = base.Simd_g_v128_load8_splat(m, v23, v41)
				v44 = int32(1)
				v45 = base.Simd_g_v128_load8_lane_l1(m, v23+v39, v41, v42)
				v47 = int32(2)
				v48 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(6), v41, v45)
				v51 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(9), v41, v48)
				v54 = base.Simd_g_const(&F_ConvertRGB24ToY_C__k0)
				v65 = base.Simd_g_v128_load8_splat(m, v23+v44, v41)
				v68 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(4), v41, v65)
				v71 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(7), v41, v68)
				v74 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(10), v41, v71)
				v77 = base.Simd_g_const(&F_ConvertRGB24ToY_C__k1)
				v89 = base.Simd_g_v128_load8_splat(m, v23+v47, v41)
				v92 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(5), v41, v89)
				v95 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(8), v41, v92)
				v98 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(11), v41, v95)
				v101 = base.Simd_g_const(&F_ConvertRGB24ToY_C__k2)
				v104 = base.Simd_g_const(&F_ConvertRGB24ToY_C__k3)
				v106 = int32(16)
				v108 = base.Simd_g_const(&F_ConvertRGB24ToY_C__k4)
				v119 = base.Simd_g_v128_load8_splat(m, v23+int32(12), v41)
				v122 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(15), v41, v119)
				v125 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(18), v41, v122)
				v128 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(21), v41, v125)
				v141 = base.Simd_g_v128_load8_splat(m, v23+int32(13), v41)
				v144 = base.Simd_g_v128_load8_lane_l1(m, v23+v106, v41, v141)
				v147 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(19), v41, v144)
				v150 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(22), v41, v147)
				v164 = base.Simd_g_v128_load8_splat(m, v23+int32(14), v41)
				v167 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(17), v41, v164)
				v170 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(20), v41, v167)
				v173 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(23), v41, v170)
				v192 = base.Simd_g_v128_load8_splat(m, v23+int32(24), v41)
				v195 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(27), v41, v192)
				v198 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(30), v41, v195)
				v201 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(33), v41, v198)
				v214 = base.Simd_g_v128_load8_splat(m, v23+int32(25), v41)
				v217 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(28), v41, v214)
				v220 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(31), v41, v217)
				v223 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(34), v41, v220)
				v237 = base.Simd_g_v128_load8_splat(m, v23+int32(26), v41)
				v240 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(29), v41, v237)
				v243 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(32), v41, v240)
				v246 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(35), v41, v243)
				v264 = base.Simd_g_v128_load8_splat(m, v23+int32(36), v41)
				v267 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(39), v41, v264)
				v270 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(42), v41, v267)
				v273 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(45), v41, v270)
				v286 = base.Simd_g_v128_load8_splat(m, v23+int32(37), v41)
				v289 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(40), v41, v286)
				v292 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(43), v41, v289)
				v295 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(46), v41, v292)
				v309 = base.Simd_g_v128_load8_splat(m, v23+int32(38), v41)
				v312 = base.Simd_g_v128_load8_lane_l1(m, v23+int32(41), v41, v309)
				v315 = base.Simd_g_v128_load8_lane_l2(m, v23+int32(44), v41, v312)
				v318 = base.Simd_g_v128_load8_lane_l3(m, v23+int32(47), v41, v315)
				v328 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v51)), v54), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v74)), v77)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v98)), v101)), v104), v106), v108), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v128)), v54), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v150)), v77)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v173)), v101)), v104), v106), v108)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v201)), v54), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v223)), v77)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v246)), v101)), v104), v106), v108), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v273)), v54), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v295)), v77)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v318)), v101)), v104), v106), v108)))
				base.Simd_g_v128_store(m, v29, v41, v328)
				v336 = v28 + int32(-16)
				if v336 != 0 {
					v23 = v23 + int32(48)
					v28 = v336
					v29 = v29 + v106
					continue
				} else {
					break
				}
				break
			}
			if v19 == l2 {
			} else {
				v341 = v19
				v342 = l0 + v19*int32(3)
				v352 = l1 + v341
				v356 = v342
				v358 = l2 - v341
				for {
					v364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356))))
					v367 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+1)))
					v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+2)))
					v378 = int32(base.Ui32(v364*int32(_a_F_ConvertRGB24ToY_C_0)+v367*int32(_a_F_ConvertRGB24ToY_C_1)+v371*int32(_a_F_ConvertRGB24ToY_C_2)+int32(1081344)) >> (uint(int32(16)) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v352))) = uint8(v378)
					v385 = v358 + int32(-1)
					if v385 != 0 {
						v352 = v352 + int32(1)
						v356 = v356 + int32(3)
						v358 = v385
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v341 = int32(0)
			v342 = l0
			v352 = l1 + v341
			v356 = v342
			v358 = l2 - v341
			for {
				v364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356))))
				v367 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+1)))
				v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v356)+2)))
				v378 = int32(base.Ui32(v364*int32(_a_F_ConvertRGB24ToY_C_0)+v367*int32(_a_F_ConvertRGB24ToY_C_1)+v371*int32(_a_F_ConvertRGB24ToY_C_2)+int32(1081344)) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v352))) = uint8(v378)
				v385 = v358 + int32(-1)
				if v385 != 0 {
					v352 = v352 + int32(1)
					v356 = v356 + int32(3)
					v358 = v385
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

var F_ConvertRGB24ToY_C__k0 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertRGB24ToY_C__k1 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertRGB24ToY_C__k2 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertRGB24ToY_C__k3 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertRGB24ToY_C__k4 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_ConvertRGB24ToY_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v75 int32
	_ = v75
	var v76 base.V128
	_ = v76
	var v78 base.V128
	_ = v78
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
	var v90 base.V128
	_ = v90
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v106 base.V128
	_ = v106
	var v108 base.V128
	_ = v108
	var v109 base.V128
	_ = v109
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v129 base.V128
	_ = v129
	var v131 base.V128
	_ = v131
	var v133 base.V128
	_ = v133
	var v135 base.V128
	_ = v135
	var v138 base.V128
	_ = v138
	var v141 base.V128
	_ = v141
	var v145 base.V128
	_ = v145
	var v157 base.V128
	_ = v157
	var v159 base.V128
	_ = v159
	var v164 base.V128
	_ = v164
	var v183 base.V128
	_ = v183
	var v187 base.V128
	_ = v187
	var v189 base.V128
	_ = v189
	var v191 base.V128
	_ = v191
	var v193 base.V128
	_ = v193
	var v198 base.V128
	_ = v198
	var v200 base.V128
	_ = v200
	var v220 base.V128
	_ = v220
	var v222 base.V128
	_ = v222
	var v227 base.V128
	_ = v227
	var v246 base.V128
	_ = v246
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v294 int32
	_ = v294
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v331 base.V128
	_ = v331
	var v333 int32
	_ = v333
	var v334 base.V128
	_ = v334
	var v336 int32
	_ = v336
	var v337 base.V128
	_ = v337
	var v340 base.V128
	_ = v340
	var v343 base.V128
	_ = v343
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
	var v378 base.V128
	_ = v378
	var v381 base.V128
	_ = v381
	var v384 base.V128
	_ = v384
	var v387 base.V128
	_ = v387
	var v390 base.V128
	_ = v390
	var v393 base.V128
	_ = v393
	var v395 int32
	_ = v395
	var v397 base.V128
	_ = v397
	var v408 base.V128
	_ = v408
	var v411 base.V128
	_ = v411
	var v414 base.V128
	_ = v414
	var v417 base.V128
	_ = v417
	var v430 base.V128
	_ = v430
	var v433 base.V128
	_ = v433
	var v436 base.V128
	_ = v436
	var v439 base.V128
	_ = v439
	var v453 base.V128
	_ = v453
	var v456 base.V128
	_ = v456
	var v459 base.V128
	_ = v459
	var v462 base.V128
	_ = v462
	var v481 base.V128
	_ = v481
	var v484 base.V128
	_ = v484
	var v487 base.V128
	_ = v487
	var v490 base.V128
	_ = v490
	var v503 base.V128
	_ = v503
	var v506 base.V128
	_ = v506
	var v509 base.V128
	_ = v509
	var v512 base.V128
	_ = v512
	var v526 base.V128
	_ = v526
	var v529 base.V128
	_ = v529
	var v532 base.V128
	_ = v532
	var v535 base.V128
	_ = v535
	var v553 base.V128
	_ = v553
	var v556 base.V128
	_ = v556
	var v559 base.V128
	_ = v559
	var v562 base.V128
	_ = v562
	var v575 base.V128
	_ = v575
	var v578 base.V128
	_ = v578
	var v581 base.V128
	_ = v581
	var v584 base.V128
	_ = v584
	var v598 base.V128
	_ = v598
	var v601 base.V128
	_ = v601
	var v604 base.V128
	_ = v604
	var v607 base.V128
	_ = v607
	var v617 base.V128
	_ = v617
	var v625 int32
	_ = v625
	var v631 int32
	_ = v631
	var v633 int32
	_ = v633
	var v657 int32
	_ = v657
	var v660 int32
	_ = v660
	var v664 int32
	_ = v664
	var v671 int32
	_ = v671
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v684 int32
	_ = v684
	var v687 int32
	_ = v687
	var v689 int32
	_ = v689
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v712 int32
	_ = v712
	var v713 int32
	_ = v713
	var v716 int32
	_ = v716
	var v717 int32
	_ = v717
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v727 int32
	_ = v727
	var v730 int32
	_ = v730
	var v734 int32
	_ = v734
	var v741 int32
	_ = v741
	var v748 int32
	_ = v748
	v4 = int32(0)
	v28 = l2 & int32(-32)
	if v28 < int32(1) {
		v254 = l0
		v257 = v4
	} else {
		v31 = l0
		v34 = v4
		for {
			v56 = l1 + v34
			v57 = int32(16)
			v59 = int32(0)
			v60 = base.Simd_g_v128_load_rng(m, v31, v59, int32(0), int32(96))
			v62 = base.Simd_g_v128_load_nc(m, v31, int32(48))
			v63 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k0)
			v64 = base.Simd_g_i8x16_shuffle2(v60, v62, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v66 = base.Simd_g_v128_load_nc(m, v31, v57)
			v68 = base.Simd_g_v128_load_nc(m, v31, int32(64))
			v69 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k3)
			v70 = base.Simd_g_i8x16_shuffle2(v66, v68, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v72 = base.Simd_g_i8x16_shuffle2(v64, v70, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v74 = base.Simd_g_i8x16_shuffle2(v60, v62, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v75 = int32(32)
			v76 = base.Simd_g_v128_load_nc(m, v31, v75)
			v78 = base.Simd_g_v128_load_nc(m, v31, int32(80))
			v80 = base.Simd_g_i8x16_shuffle2(v76, v78, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v82 = base.Simd_g_i8x16_shuffle2(v74, v80, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v84 = base.Simd_g_i8x16_shuffle2(v72, v82, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v86 = base.Simd_g_i8x16_shuffle2(v64, v70, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v88 = base.Simd_g_i8x16_shuffle2(v66, v68, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v90 = base.Simd_g_i8x16_shuffle2(v76, v78, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v92 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v94 = base.Simd_g_i8x16_shuffle2(v86, v92, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v96 = base.Simd_g_i8x16_shuffle2(v84, v94, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v98 = base.Simd_g_i8x16_shuffle2(v72, v82, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v100 = base.Simd_g_i8x16_shuffle2(v74, v80, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v102 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v104 = base.Simd_g_i8x16_shuffle2(v100, v102, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v106 = base.Simd_g_i8x16_shuffle2(v98, v104, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v108 = base.Simd_g_i8x16_shuffle2(v96, v106, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v109 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k6)
			v111 = base.Simd_g_i8x16_shuffle2(v108, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v113 = base.Simd_g_i8x16_shuffle2(v84, v94, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v115 = base.Simd_g_i8x16_shuffle2(v86, v92, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v117 = base.Simd_g_i8x16_shuffle2(v100, v102, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v119 = base.Simd_g_i8x16_shuffle2(v115, v117, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v121 = base.Simd_g_i8x16_shuffle2(v113, v119, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v123 = base.Simd_g_i8x16_shuffle2(v121, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v124 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k7)
			v126 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k8)
			v129 = base.Simd_g_i8x16_shuffle2(v98, v104, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v131 = base.Simd_g_i8x16_shuffle2(v115, v117, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v133 = base.Simd_g_i8x16_shuffle2(v129, v131, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v135 = base.Simd_g_i8x16_shuffle2(v133, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v138 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k9)
			v141 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k10)
			v145 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k11)
			v157 = base.Simd_g_i8x16_shuffle2(v108, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v159 = base.Simd_g_i8x16_shuffle2(v121, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v164 = base.Simd_g_i8x16_shuffle2(v133, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v183 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v111, v123, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v123, v135, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v111, v123, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v123, v135, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v138)), v141), v57)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v157, v159, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v159, v164, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v157, v159, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v159, v164, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v138)), v141), v57)))
			base.Simd_g_v128_store(m, v56+v57, v59, v183)
			v187 = base.Simd_g_i8x16_shuffle2(v96, v106, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v189 = base.Simd_g_i8x16_shuffle2(v187, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v191 = base.Simd_g_i8x16_shuffle2(v113, v119, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v193 = base.Simd_g_i8x16_shuffle2(v191, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v198 = base.Simd_g_i8x16_shuffle2(v129, v131, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v200 = base.Simd_g_i8x16_shuffle2(v198, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k1), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k2))
			v220 = base.Simd_g_i8x16_shuffle2(v187, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v222 = base.Simd_g_i8x16_shuffle2(v191, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v227 = base.Simd_g_i8x16_shuffle2(v198, v109, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k4), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k5))
			v246 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v189, v193, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v193, v200, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v189, v193, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v193, v200, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v138)), v141), v57)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v220, v222, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v222, v227, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k12), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k13)), v138)), v141), v57), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v220, v222, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v126), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v222, v227, base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k14), base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k15)), v138)), v141), v57)))
			base.Simd_g_v128_store(m, v56, v59, v246)
			v250 = v31 + int32(96)
			v252 = v34 + v75
			if v252 < v28 {
				v31 = v250
				v34 = v252
				continue
			} else {
				break
			}
			break
		}
		v254 = v250
		v257 = v252
	}
	if l2 <= v257 {
	} else {
		v280 = l2 - v257
		if base.Ui32(v280) <= base.Ui32(int32(15)) {
			v631 = v257
			v633 = v254
			if (l2-v631)&int32(1) != 0 {
				v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
				v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
				v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
				v671 = int32(base.Ui32(v657*int32(_a_F_ConvertRGB24ToY_SSE2_0)+v660*int32(_a_F_ConvertRGB24ToY_SSE2_1)+v664*int32(_a_F_ConvertRGB24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
				v677 = v631 + int32(1)
				v678 = v633 + int32(3)
			} else {
				v677 = v631
				v678 = v633
			}
			if v631 == l2+int32(-1) {
			} else {
				v684 = l1 + v677
				v687 = l2 - v677
				v689 = v678
				for {
					v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
					v710 = int32(_a_F_ConvertRGB24ToY_SSE2_0)
					v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
					v713 = int32(_a_F_ConvertRGB24ToY_SSE2_1)
					v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
					v717 = int32(_a_F_ConvertRGB24ToY_SSE2_2)
					v720 = int32(1081344)
					v722 = int32(16)
					v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
					v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
					v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
					v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
					v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
					v748 = v687 + int32(-2)
					if v748 != 0 {
						v684 = v684 + int32(2)
						v687 = v748
						v689 = v689 + int32(6)
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v283 = l1 + v257
			if base.Ui32(v254+(v257*int32(-3)+l2*int32(3))) <= base.Ui32(v283) {
				v294 = v280 & int32(-16)
				v299 = v254
				v303 = v283
				v323 = v294
				for {
					v328 = int32(3)
					v330 = int32(0)
					v331 = base.Simd_g_v128_load8_splat(m, v299, v330)
					v333 = int32(1)
					v334 = base.Simd_g_v128_load8_lane_l1(m, v299+v328, v330, v331)
					v336 = int32(2)
					v337 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(6), v330, v334)
					v340 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(9), v330, v337)
					v343 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k16)
					v354 = base.Simd_g_v128_load8_splat(m, v299+v333, v330)
					v357 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(4), v330, v354)
					v360 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(7), v330, v357)
					v363 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(10), v330, v360)
					v366 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k17)
					v378 = base.Simd_g_v128_load8_splat(m, v299+v336, v330)
					v381 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(5), v330, v378)
					v384 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(8), v330, v381)
					v387 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(11), v330, v384)
					v390 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k18)
					v393 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k10)
					v395 = int32(16)
					v397 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k19)
					v408 = base.Simd_g_v128_load8_splat(m, v299+int32(12), v330)
					v411 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(15), v330, v408)
					v414 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(18), v330, v411)
					v417 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(21), v330, v414)
					v430 = base.Simd_g_v128_load8_splat(m, v299+int32(13), v330)
					v433 = base.Simd_g_v128_load8_lane_l1(m, v299+v395, v330, v430)
					v436 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(19), v330, v433)
					v439 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(22), v330, v436)
					v453 = base.Simd_g_v128_load8_splat(m, v299+int32(14), v330)
					v456 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(17), v330, v453)
					v459 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(20), v330, v456)
					v462 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(23), v330, v459)
					v481 = base.Simd_g_v128_load8_splat(m, v299+int32(24), v330)
					v484 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(27), v330, v481)
					v487 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(30), v330, v484)
					v490 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(33), v330, v487)
					v503 = base.Simd_g_v128_load8_splat(m, v299+int32(25), v330)
					v506 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(28), v330, v503)
					v509 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(31), v330, v506)
					v512 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(34), v330, v509)
					v526 = base.Simd_g_v128_load8_splat(m, v299+int32(26), v330)
					v529 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(29), v330, v526)
					v532 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(32), v330, v529)
					v535 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(35), v330, v532)
					v553 = base.Simd_g_v128_load8_splat(m, v299+int32(36), v330)
					v556 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(39), v330, v553)
					v559 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(42), v330, v556)
					v562 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(45), v330, v559)
					v575 = base.Simd_g_v128_load8_splat(m, v299+int32(37), v330)
					v578 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(40), v330, v575)
					v581 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(43), v330, v578)
					v584 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(46), v330, v581)
					v598 = base.Simd_g_v128_load8_splat(m, v299+int32(38), v330)
					v601 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(41), v330, v598)
					v604 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(44), v330, v601)
					v607 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(47), v330, v604)
					v617 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v340)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v363)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v387)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v417)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v439)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v462)), v390)), v393), v395), v397)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v490)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v512)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v535)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v562)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v584)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v607)), v390)), v393), v395), v397)))
					base.Simd_g_v128_store(m, v303, v330, v617)
					v625 = v323 + int32(-16)
					if v625 != 0 {
						v299 = v299 + int32(48)
						v303 = v303 + v395
						v323 = v625
						continue
					} else {
						break
					}
					break
				}
				if v280 != v294 {
					v631 = v257 + v294
					v633 = v254 + v294*int32(3)
					if (l2-v631)&int32(1) != 0 {
						v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
						v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
						v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
						v671 = int32(base.Ui32(v657*int32(_a_F_ConvertRGB24ToY_SSE2_0)+v660*int32(_a_F_ConvertRGB24ToY_SSE2_1)+v664*int32(_a_F_ConvertRGB24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
						v677 = v631 + int32(1)
						v678 = v633 + int32(3)
					} else {
						v677 = v631
						v678 = v633
					}
					if v631 == l2+int32(-1) {
					} else {
						v684 = l1 + v677
						v687 = l2 - v677
						v689 = v678
						for {
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
							v710 = int32(_a_F_ConvertRGB24ToY_SSE2_0)
							v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
							v713 = int32(_a_F_ConvertRGB24ToY_SSE2_1)
							v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
							v717 = int32(_a_F_ConvertRGB24ToY_SSE2_2)
							v720 = int32(1081344)
							v722 = int32(16)
							v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
							v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
							v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
							v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
							v748 = v687 + int32(-2)
							if v748 != 0 {
								v684 = v684 + int32(2)
								v687 = v748
								v689 = v689 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
				}
			} else {
				if base.Ui32(v254) < base.Ui32(l1+l2) {
					v631 = v257
					v633 = v254
					if (l2-v631)&int32(1) != 0 {
						v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
						v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
						v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
						v671 = int32(base.Ui32(v657*int32(_a_F_ConvertRGB24ToY_SSE2_0)+v660*int32(_a_F_ConvertRGB24ToY_SSE2_1)+v664*int32(_a_F_ConvertRGB24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
						v677 = v631 + int32(1)
						v678 = v633 + int32(3)
					} else {
						v677 = v631
						v678 = v633
					}
					if v631 == l2+int32(-1) {
					} else {
						v684 = l1 + v677
						v687 = l2 - v677
						v689 = v678
						for {
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
							v710 = int32(_a_F_ConvertRGB24ToY_SSE2_0)
							v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
							v713 = int32(_a_F_ConvertRGB24ToY_SSE2_1)
							v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
							v717 = int32(_a_F_ConvertRGB24ToY_SSE2_2)
							v720 = int32(1081344)
							v722 = int32(16)
							v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
							v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
							v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
							v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
							v748 = v687 + int32(-2)
							if v748 != 0 {
								v684 = v684 + int32(2)
								v687 = v748
								v689 = v689 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v294 = v280 & int32(-16)
					v299 = v254
					v303 = v283
					v323 = v294
					for {
						v328 = int32(3)
						v330 = int32(0)
						v331 = base.Simd_g_v128_load8_splat(m, v299, v330)
						v333 = int32(1)
						v334 = base.Simd_g_v128_load8_lane_l1(m, v299+v328, v330, v331)
						v336 = int32(2)
						v337 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(6), v330, v334)
						v340 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(9), v330, v337)
						v343 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k16)
						v354 = base.Simd_g_v128_load8_splat(m, v299+v333, v330)
						v357 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(4), v330, v354)
						v360 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(7), v330, v357)
						v363 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(10), v330, v360)
						v366 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k17)
						v378 = base.Simd_g_v128_load8_splat(m, v299+v336, v330)
						v381 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(5), v330, v378)
						v384 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(8), v330, v381)
						v387 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(11), v330, v384)
						v390 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k18)
						v393 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k10)
						v395 = int32(16)
						v397 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE2__k19)
						v408 = base.Simd_g_v128_load8_splat(m, v299+int32(12), v330)
						v411 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(15), v330, v408)
						v414 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(18), v330, v411)
						v417 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(21), v330, v414)
						v430 = base.Simd_g_v128_load8_splat(m, v299+int32(13), v330)
						v433 = base.Simd_g_v128_load8_lane_l1(m, v299+v395, v330, v430)
						v436 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(19), v330, v433)
						v439 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(22), v330, v436)
						v453 = base.Simd_g_v128_load8_splat(m, v299+int32(14), v330)
						v456 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(17), v330, v453)
						v459 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(20), v330, v456)
						v462 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(23), v330, v459)
						v481 = base.Simd_g_v128_load8_splat(m, v299+int32(24), v330)
						v484 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(27), v330, v481)
						v487 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(30), v330, v484)
						v490 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(33), v330, v487)
						v503 = base.Simd_g_v128_load8_splat(m, v299+int32(25), v330)
						v506 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(28), v330, v503)
						v509 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(31), v330, v506)
						v512 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(34), v330, v509)
						v526 = base.Simd_g_v128_load8_splat(m, v299+int32(26), v330)
						v529 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(29), v330, v526)
						v532 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(32), v330, v529)
						v535 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(35), v330, v532)
						v553 = base.Simd_g_v128_load8_splat(m, v299+int32(36), v330)
						v556 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(39), v330, v553)
						v559 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(42), v330, v556)
						v562 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(45), v330, v559)
						v575 = base.Simd_g_v128_load8_splat(m, v299+int32(37), v330)
						v578 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(40), v330, v575)
						v581 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(43), v330, v578)
						v584 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(46), v330, v581)
						v598 = base.Simd_g_v128_load8_splat(m, v299+int32(38), v330)
						v601 = base.Simd_g_v128_load8_lane_l1(m, v299+int32(41), v330, v598)
						v604 = base.Simd_g_v128_load8_lane_l2(m, v299+int32(44), v330, v601)
						v607 = base.Simd_g_v128_load8_lane_l3(m, v299+int32(47), v330, v604)
						v617 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v340)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v363)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v387)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v417)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v439)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v462)), v390)), v393), v395), v397)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v490)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v512)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v535)), v390)), v393), v395), v397), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v562)), v343), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v584)), v366)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v607)), v390)), v393), v395), v397)))
						base.Simd_g_v128_store(m, v303, v330, v617)
						v625 = v323 + int32(-16)
						if v625 != 0 {
							v299 = v299 + int32(48)
							v303 = v303 + v395
							v323 = v625
							continue
						} else {
							break
						}
						break
					}
					if v280 != v294 {
						v631 = v257 + v294
						v633 = v254 + v294*int32(3)
						if (l2-v631)&int32(1) != 0 {
							v657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633))))
							v660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+1)))
							v664 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633)+2)))
							v671 = int32(base.Ui32(v657*int32(_a_F_ConvertRGB24ToY_SSE2_0)+v660*int32(_a_F_ConvertRGB24ToY_SSE2_1)+v664*int32(_a_F_ConvertRGB24ToY_SSE2_2)+int32(1081344)) >> (uint(int32(16)) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(l1+v631))) = uint8(v671)
							v677 = v631 + int32(1)
							v678 = v633 + int32(3)
						} else {
							v677 = v631
							v678 = v633
						}
						if v631 == l2+int32(-1) {
						} else {
							v684 = l1 + v677
							v687 = l2 - v677
							v689 = v678
							for {
								v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689))))
								v710 = int32(_a_F_ConvertRGB24ToY_SSE2_0)
								v712 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+1)))
								v713 = int32(_a_F_ConvertRGB24ToY_SSE2_1)
								v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+2)))
								v717 = int32(_a_F_ConvertRGB24ToY_SSE2_2)
								v720 = int32(1081344)
								v722 = int32(16)
								v723 = int32(base.Ui32(v709*v710+v712*v713+v716*v717+v720) >> (uint(v722) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v684))) = uint8(v723)
								v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+3)))
								v730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+4)))
								v734 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v689)+5)))
								v741 = int32(base.Ui32(v727*v710+v730*v713+v734*v717+v720) >> (uint(v722) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v684+int32(1)))) = uint8(v741)
								v748 = v687 + int32(-2)
								if v748 != 0 {
									v684 = v684 + int32(2)
									v687 = v748
									v689 = v689 + int32(6)
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
					}
				}
			}
		}
	}
	return
}

var F_ConvertRGB24ToY_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertRGB24ToY_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertRGB24ToY_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertRGB24ToY_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertRGB24ToY_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertRGB24ToY_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertRGB24ToY_SSE2__k6 = [2]uint64{0x0, 0x0}
var F_ConvertRGB24ToY_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertRGB24ToY_SSE2__k8 = [2]uint64{0x412341c7412341c7, 0x412341c7412341c7}
var F_ConvertRGB24ToY_SSE2__k9 = [2]uint64{0x1914400019144000, 0x1914400019144000}
var F_ConvertRGB24ToY_SSE2__k10 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertRGB24ToY_SSE2__k11 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertRGB24ToY_SSE2__k12 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertRGB24ToY_SSE2__k13 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertRGB24ToY_SSE2__k14 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertRGB24ToY_SSE2__k15 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertRGB24ToY_SSE2__k16 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertRGB24ToY_SSE2__k17 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertRGB24ToY_SSE2__k18 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertRGB24ToY_SSE2__k19 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_ConvertRGB24ToY_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v66 base.V128
	_ = v66
	var v67 base.V128
	_ = v67
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v78 base.V128
	_ = v78
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v81 base.V128
	_ = v81
	var v82 base.V128
	_ = v82
	var v84 base.V128
	_ = v84
	var v87 base.V128
	_ = v87
	var v89 base.V128
	_ = v89
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v98 base.V128
	_ = v98
	var v101 base.V128
	_ = v101
	var v103 base.V128
	_ = v103
	var v105 base.V128
	_ = v105
	var v108 base.V128
	_ = v108
	var v111 base.V128
	_ = v111
	var v115 base.V128
	_ = v115
	var v126 base.V128
	_ = v126
	var v127 base.V128
	_ = v127
	var v129 base.V128
	_ = v129
	var v134 base.V128
	_ = v134
	var v153 base.V128
	_ = v153
	var v154 int32
	_ = v154
	var v157 base.V128
	_ = v157
	var v160 base.V128
	_ = v160
	var v163 int32
	_ = v163
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v173 base.V128
	_ = v173
	var v175 base.V128
	_ = v175
	var v183 base.V128
	_ = v183
	var v185 base.V128
	_ = v185
	var v205 base.V128
	_ = v205
	var v207 base.V128
	_ = v207
	var v212 base.V128
	_ = v212
	var v231 base.V128
	_ = v231
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v282 int32
	_ = v282
	var v287 int32
	_ = v287
	var v291 int32
	_ = v291
	var v314 int32
	_ = v314
	var v319 int32
	_ = v319
	var v321 int32
	_ = v321
	var v322 base.V128
	_ = v322
	var v324 int32
	_ = v324
	var v325 base.V128
	_ = v325
	var v327 int32
	_ = v327
	var v328 base.V128
	_ = v328
	var v331 base.V128
	_ = v331
	var v334 base.V128
	_ = v334
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
	var v384 base.V128
	_ = v384
	var v386 int32
	_ = v386
	var v388 base.V128
	_ = v388
	var v399 base.V128
	_ = v399
	var v402 base.V128
	_ = v402
	var v405 base.V128
	_ = v405
	var v408 base.V128
	_ = v408
	var v421 base.V128
	_ = v421
	var v424 base.V128
	_ = v424
	var v427 base.V128
	_ = v427
	var v430 base.V128
	_ = v430
	var v444 base.V128
	_ = v444
	var v447 base.V128
	_ = v447
	var v450 base.V128
	_ = v450
	var v453 base.V128
	_ = v453
	var v472 base.V128
	_ = v472
	var v475 base.V128
	_ = v475
	var v478 base.V128
	_ = v478
	var v481 base.V128
	_ = v481
	var v494 base.V128
	_ = v494
	var v497 base.V128
	_ = v497
	var v500 base.V128
	_ = v500
	var v503 base.V128
	_ = v503
	var v517 base.V128
	_ = v517
	var v520 base.V128
	_ = v520
	var v523 base.V128
	_ = v523
	var v526 base.V128
	_ = v526
	var v544 base.V128
	_ = v544
	var v547 base.V128
	_ = v547
	var v550 base.V128
	_ = v550
	var v553 base.V128
	_ = v553
	var v566 base.V128
	_ = v566
	var v569 base.V128
	_ = v569
	var v572 base.V128
	_ = v572
	var v575 base.V128
	_ = v575
	var v589 base.V128
	_ = v589
	var v592 base.V128
	_ = v592
	var v595 base.V128
	_ = v595
	var v598 base.V128
	_ = v598
	var v608 base.V128
	_ = v608
	var v616 int32
	_ = v616
	var v622 int32
	_ = v622
	var v624 int32
	_ = v624
	var v651 int32
	_ = v651
	var v654 int32
	_ = v654
	var v658 int32
	_ = v658
	var v665 int32
	_ = v665
	var v671 int32
	_ = v671
	var v672 int32
	_ = v672
	var v678 int32
	_ = v678
	var v681 int32
	_ = v681
	var v683 int32
	_ = v683
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v713 int32
	_ = v713
	var v714 int32
	_ = v714
	var v717 int32
	_ = v717
	var v719 int32
	_ = v719
	var v720 int32
	_ = v720
	var v724 int32
	_ = v724
	var v727 int32
	_ = v727
	var v731 int32
	_ = v731
	var v738 int32
	_ = v738
	var v745 int32
	_ = v745
	v4 = int32(0)
	v31 = l2 & int32(-32)
	if v31 < int32(1) {
		v239 = l0
		v242 = v4
	} else {
		v34 = l0
		v37 = v4
		for {
			v62 = l1 + v37
			v63 = int32(16)
			v66 = base.Simd_g_v128_load_rng(m, v34, int32(64), int32(48), int32(48))
			v67 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k0)
			v70 = base.Simd_g_v128_load_nc(m, v34, int32(48))
			v71 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k1)
			v75 = base.Simd_g_v128_load_nc(m, v34, int32(80))
			v76 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k2)
			v78 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v66, v67), base.Simd_g_i8x16_swizzle(v70, v71)), base.Simd_g_i8x16_swizzle(v75, v76))
			v79 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k3)
			v80 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k4)
			v81 = base.Simd_g_i8x16_shuffle2(v78, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k6))
			v82 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k7)
			v84 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k8)
			v87 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k9)
			v89 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v66, v82), base.Simd_g_i8x16_swizzle(v70, v84)), base.Simd_g_i8x16_swizzle(v75, v87))
			v91 = base.Simd_g_i8x16_shuffle2(v89, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k6))
			v92 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k10)
			v94 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k11)
			v96 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k12)
			v98 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k13)
			v101 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k14)
			v103 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v66, v96), base.Simd_g_i8x16_swizzle(v70, v98)), base.Simd_g_i8x16_swizzle(v75, v101))
			v105 = base.Simd_g_i8x16_shuffle2(v103, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k6))
			v108 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k15)
			v111 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k16)
			v115 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k17)
			v126 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k18)
			v127 = base.Simd_g_i8x16_shuffle2(v78, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k20))
			v129 = base.Simd_g_i8x16_shuffle2(v89, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k20))
			v134 = base.Simd_g_i8x16_shuffle2(v103, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k20))
			v153 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v81, v91, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v91, v105, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v81, v91, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v91, v105, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v108)), v111), v63)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v129, v134, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v127, v129, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v129, v134, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v108)), v111), v63)))
			v154 = int32(0)
			base.Simd_g_v128_store(m, v62+v63, v154, v153)
			v157 = base.Simd_g_v128_load_rng(m, v34, v63, int32(0), int32(48))
			v160 = base.Simd_g_v128_load_nc(m, v34, v154)
			v163 = int32(32)
			v164 = base.Simd_g_v128_load_nc(m, v34, v163)
			v166 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v157, v67), base.Simd_g_i8x16_swizzle(v160, v71)), base.Simd_g_i8x16_swizzle(v164, v76))
			v168 = base.Simd_g_i8x16_shuffle2(v166, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k6))
			v173 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v157, v82), base.Simd_g_i8x16_swizzle(v160, v84)), base.Simd_g_i8x16_swizzle(v164, v87))
			v175 = base.Simd_g_i8x16_shuffle2(v173, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k6))
			v183 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v157, v96), base.Simd_g_i8x16_swizzle(v160, v98)), base.Simd_g_i8x16_swizzle(v164, v101))
			v185 = base.Simd_g_i8x16_shuffle2(v183, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k5), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k6))
			v205 = base.Simd_g_i8x16_shuffle2(v166, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k20))
			v207 = base.Simd_g_i8x16_shuffle2(v173, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k20))
			v212 = base.Simd_g_i8x16_shuffle2(v183, v79, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k19), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k20))
			v231 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v168, v175, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v175, v185, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v168, v175, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v175, v185, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v108)), v111), v63)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v205, v207, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v207, v212, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k21), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k22)), v108)), v111), v63), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v205, v207, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v94), base.Simd_g_i32x4_dot_i16x8_s(base.Simd_g_i8x16_shuffle2(v207, v212, base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k23), base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k24)), v108)), v111), v63)))
			base.Simd_g_v128_store(m, v62, v154, v231)
			v235 = v34 + int32(96)
			v237 = v37 + v163
			if v237 < v31 {
				v34 = v235
				v37 = v237
				continue
			} else {
				break
			}
			break
		}
		v239 = v235
		v242 = v237
	}
	if l2 <= v242 {
	} else {
		v268 = l2 - v242
		if base.Ui32(v268) <= base.Ui32(int32(15)) {
			v622 = v242
			v624 = v239
			if (l2-v622)&int32(1) != 0 {
				v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
				v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
				v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
				v665 = int32(base.Ui32(v651*int32(_a_F_ConvertRGB24ToY_SSE41_0)+v654*int32(_a_F_ConvertRGB24ToY_SSE41_1)+v658*int32(_a_F_ConvertRGB24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
				v671 = v622 + int32(1)
				v672 = v624 + int32(3)
			} else {
				v671 = v622
				v672 = v624
			}
			if v622 == l2+int32(-1) {
			} else {
				v678 = l1 + v671
				v681 = l2 - v671
				v683 = v672
				for {
					v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
					v707 = int32(_a_F_ConvertRGB24ToY_SSE41_0)
					v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
					v710 = int32(_a_F_ConvertRGB24ToY_SSE41_1)
					v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
					v714 = int32(_a_F_ConvertRGB24ToY_SSE41_2)
					v717 = int32(1081344)
					v719 = int32(16)
					v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
					v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
					v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
					v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
					v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
					v745 = v681 + int32(-2)
					if v745 != 0 {
						v678 = v678 + int32(2)
						v681 = v745
						v683 = v683 + int32(6)
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v271 = l1 + v242
			if base.Ui32(v239+(v242*int32(-3)+l2*int32(3))) <= base.Ui32(v271) {
				v282 = v268 & int32(-16)
				v287 = v239
				v291 = v271
				v314 = v282
				for {
					v319 = int32(3)
					v321 = int32(0)
					v322 = base.Simd_g_v128_load8_splat(m, v287, v321)
					v324 = int32(1)
					v325 = base.Simd_g_v128_load8_lane_l1(m, v287+v319, v321, v322)
					v327 = int32(2)
					v328 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(6), v321, v325)
					v331 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(9), v321, v328)
					v334 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k25)
					v345 = base.Simd_g_v128_load8_splat(m, v287+v324, v321)
					v348 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(4), v321, v345)
					v351 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(7), v321, v348)
					v354 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(10), v321, v351)
					v357 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k26)
					v369 = base.Simd_g_v128_load8_splat(m, v287+v327, v321)
					v372 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(5), v321, v369)
					v375 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(8), v321, v372)
					v378 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(11), v321, v375)
					v381 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k27)
					v384 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k16)
					v386 = int32(16)
					v388 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k28)
					v399 = base.Simd_g_v128_load8_splat(m, v287+int32(12), v321)
					v402 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(15), v321, v399)
					v405 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(18), v321, v402)
					v408 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(21), v321, v405)
					v421 = base.Simd_g_v128_load8_splat(m, v287+int32(13), v321)
					v424 = base.Simd_g_v128_load8_lane_l1(m, v287+v386, v321, v421)
					v427 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(19), v321, v424)
					v430 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(22), v321, v427)
					v444 = base.Simd_g_v128_load8_splat(m, v287+int32(14), v321)
					v447 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(17), v321, v444)
					v450 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(20), v321, v447)
					v453 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(23), v321, v450)
					v472 = base.Simd_g_v128_load8_splat(m, v287+int32(24), v321)
					v475 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(27), v321, v472)
					v478 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(30), v321, v475)
					v481 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(33), v321, v478)
					v494 = base.Simd_g_v128_load8_splat(m, v287+int32(25), v321)
					v497 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(28), v321, v494)
					v500 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(31), v321, v497)
					v503 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(34), v321, v500)
					v517 = base.Simd_g_v128_load8_splat(m, v287+int32(26), v321)
					v520 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(29), v321, v517)
					v523 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(32), v321, v520)
					v526 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(35), v321, v523)
					v544 = base.Simd_g_v128_load8_splat(m, v287+int32(36), v321)
					v547 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(39), v321, v544)
					v550 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(42), v321, v547)
					v553 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(45), v321, v550)
					v566 = base.Simd_g_v128_load8_splat(m, v287+int32(37), v321)
					v569 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(40), v321, v566)
					v572 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(43), v321, v569)
					v575 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(46), v321, v572)
					v589 = base.Simd_g_v128_load8_splat(m, v287+int32(38), v321)
					v592 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(41), v321, v589)
					v595 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(44), v321, v592)
					v598 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(47), v321, v595)
					v608 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v331)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v354)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v378)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v408)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v430)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v453)), v381)), v384), v386), v388)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v481)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v503)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v526)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v553)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v575)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v598)), v381)), v384), v386), v388)))
					base.Simd_g_v128_store(m, v291, v321, v608)
					v616 = v314 + int32(-16)
					if v616 != 0 {
						v287 = v287 + int32(48)
						v291 = v291 + v386
						v314 = v616
						continue
					} else {
						break
					}
					break
				}
				if v268 != v282 {
					v622 = v242 + v282
					v624 = v239 + v282*int32(3)
					if (l2-v622)&int32(1) != 0 {
						v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
						v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
						v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
						v665 = int32(base.Ui32(v651*int32(_a_F_ConvertRGB24ToY_SSE41_0)+v654*int32(_a_F_ConvertRGB24ToY_SSE41_1)+v658*int32(_a_F_ConvertRGB24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
						v671 = v622 + int32(1)
						v672 = v624 + int32(3)
					} else {
						v671 = v622
						v672 = v624
					}
					if v622 == l2+int32(-1) {
					} else {
						v678 = l1 + v671
						v681 = l2 - v671
						v683 = v672
						for {
							v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
							v707 = int32(_a_F_ConvertRGB24ToY_SSE41_0)
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
							v710 = int32(_a_F_ConvertRGB24ToY_SSE41_1)
							v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
							v714 = int32(_a_F_ConvertRGB24ToY_SSE41_2)
							v717 = int32(1081344)
							v719 = int32(16)
							v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
							v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
							v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
							v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
							v745 = v681 + int32(-2)
							if v745 != 0 {
								v678 = v678 + int32(2)
								v681 = v745
								v683 = v683 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
				}
			} else {
				if base.Ui32(v239) < base.Ui32(l1+l2) {
					v622 = v242
					v624 = v239
					if (l2-v622)&int32(1) != 0 {
						v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
						v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
						v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
						v665 = int32(base.Ui32(v651*int32(_a_F_ConvertRGB24ToY_SSE41_0)+v654*int32(_a_F_ConvertRGB24ToY_SSE41_1)+v658*int32(_a_F_ConvertRGB24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
						v671 = v622 + int32(1)
						v672 = v624 + int32(3)
					} else {
						v671 = v622
						v672 = v624
					}
					if v622 == l2+int32(-1) {
					} else {
						v678 = l1 + v671
						v681 = l2 - v671
						v683 = v672
						for {
							v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
							v707 = int32(_a_F_ConvertRGB24ToY_SSE41_0)
							v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
							v710 = int32(_a_F_ConvertRGB24ToY_SSE41_1)
							v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
							v714 = int32(_a_F_ConvertRGB24ToY_SSE41_2)
							v717 = int32(1081344)
							v719 = int32(16)
							v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
							v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
							v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
							v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
							v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
							v745 = v681 + int32(-2)
							if v745 != 0 {
								v678 = v678 + int32(2)
								v681 = v745
								v683 = v683 + int32(6)
								continue
							} else {
								break
							}
							break
						}
					}
				} else {
					v282 = v268 & int32(-16)
					v287 = v239
					v291 = v271
					v314 = v282
					for {
						v319 = int32(3)
						v321 = int32(0)
						v322 = base.Simd_g_v128_load8_splat(m, v287, v321)
						v324 = int32(1)
						v325 = base.Simd_g_v128_load8_lane_l1(m, v287+v319, v321, v322)
						v327 = int32(2)
						v328 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(6), v321, v325)
						v331 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(9), v321, v328)
						v334 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k25)
						v345 = base.Simd_g_v128_load8_splat(m, v287+v324, v321)
						v348 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(4), v321, v345)
						v351 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(7), v321, v348)
						v354 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(10), v321, v351)
						v357 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k26)
						v369 = base.Simd_g_v128_load8_splat(m, v287+v327, v321)
						v372 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(5), v321, v369)
						v375 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(8), v321, v372)
						v378 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(11), v321, v375)
						v381 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k27)
						v384 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k16)
						v386 = int32(16)
						v388 = base.Simd_g_const(&F_ConvertRGB24ToY_SSE41__k28)
						v399 = base.Simd_g_v128_load8_splat(m, v287+int32(12), v321)
						v402 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(15), v321, v399)
						v405 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(18), v321, v402)
						v408 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(21), v321, v405)
						v421 = base.Simd_g_v128_load8_splat(m, v287+int32(13), v321)
						v424 = base.Simd_g_v128_load8_lane_l1(m, v287+v386, v321, v421)
						v427 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(19), v321, v424)
						v430 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(22), v321, v427)
						v444 = base.Simd_g_v128_load8_splat(m, v287+int32(14), v321)
						v447 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(17), v321, v444)
						v450 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(20), v321, v447)
						v453 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(23), v321, v450)
						v472 = base.Simd_g_v128_load8_splat(m, v287+int32(24), v321)
						v475 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(27), v321, v472)
						v478 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(30), v321, v475)
						v481 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(33), v321, v478)
						v494 = base.Simd_g_v128_load8_splat(m, v287+int32(25), v321)
						v497 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(28), v321, v494)
						v500 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(31), v321, v497)
						v503 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(34), v321, v500)
						v517 = base.Simd_g_v128_load8_splat(m, v287+int32(26), v321)
						v520 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(29), v321, v517)
						v523 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(32), v321, v520)
						v526 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(35), v321, v523)
						v544 = base.Simd_g_v128_load8_splat(m, v287+int32(36), v321)
						v547 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(39), v321, v544)
						v550 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(42), v321, v547)
						v553 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(45), v321, v550)
						v566 = base.Simd_g_v128_load8_splat(m, v287+int32(37), v321)
						v569 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(40), v321, v566)
						v572 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(43), v321, v569)
						v575 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(46), v321, v572)
						v589 = base.Simd_g_v128_load8_splat(m, v287+int32(38), v321)
						v592 = base.Simd_g_v128_load8_lane_l1(m, v287+int32(41), v321, v589)
						v595 = base.Simd_g_v128_load8_lane_l2(m, v287+int32(44), v321, v592)
						v598 = base.Simd_g_v128_load8_lane_l3(m, v287+int32(47), v321, v595)
						v608 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v331)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v354)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v378)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v408)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v430)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v453)), v381)), v384), v386), v388)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v481)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v503)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v526)), v381)), v384), v386), v388), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v553)), v334), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v575)), v357)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v598)), v381)), v384), v386), v388)))
						base.Simd_g_v128_store(m, v291, v321, v608)
						v616 = v314 + int32(-16)
						if v616 != 0 {
							v287 = v287 + int32(48)
							v291 = v291 + v386
							v314 = v616
							continue
						} else {
							break
						}
						break
					}
					if v268 != v282 {
						v622 = v242 + v282
						v624 = v239 + v282*int32(3)
						if (l2-v622)&int32(1) != 0 {
							v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624))))
							v654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+1)))
							v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v624)+2)))
							v665 = int32(base.Ui32(v651*int32(_a_F_ConvertRGB24ToY_SSE41_0)+v654*int32(_a_F_ConvertRGB24ToY_SSE41_1)+v658*int32(_a_F_ConvertRGB24ToY_SSE41_2)+int32(1081344)) >> (uint(int32(16)) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(l1+v622))) = uint8(v665)
							v671 = v622 + int32(1)
							v672 = v624 + int32(3)
						} else {
							v671 = v622
							v672 = v624
						}
						if v622 == l2+int32(-1) {
						} else {
							v678 = l1 + v671
							v681 = l2 - v671
							v683 = v672
							for {
								v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683))))
								v707 = int32(_a_F_ConvertRGB24ToY_SSE41_0)
								v709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+1)))
								v710 = int32(_a_F_ConvertRGB24ToY_SSE41_1)
								v713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+2)))
								v714 = int32(_a_F_ConvertRGB24ToY_SSE41_2)
								v717 = int32(1081344)
								v719 = int32(16)
								v720 = int32(base.Ui32(v706*v707+v709*v710+v713*v714+v717) >> (uint(v719) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v720)
								v724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+3)))
								v727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+4)))
								v731 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683)+5)))
								v738 = int32(base.Ui32(v724*v707+v727*v710+v731*v714+v717) >> (uint(v719) % 32))
								*(*uint8)(unsafe.Add(mBase, uint32(v678+int32(1)))) = uint8(v738)
								v745 = v681 + int32(-2)
								if v745 != 0 {
									v678 = v678 + int32(2)
									v681 = v745
									v683 = v683 + int32(6)
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
					}
				}
			}
		}
	}
	return
}

var F_ConvertRGB24ToY_SSE41__k0 = [2]uint64{0x5028f8f8f8f8f8f, 0x8f8f8f8f8f0e0b08}
var F_ConvertRGB24ToY_SSE41__k1 = [2]uint64{0x8f8f0f0c09060300, 0x8f8f8f8f8f8f8f8f}
var F_ConvertRGB24ToY_SSE41__k2 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xd0a0704018f8f8f}
var F_ConvertRGB24ToY_SSE41__k3 = [2]uint64{0x0, 0x0}
var F_ConvertRGB24ToY_SSE41__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_ConvertRGB24ToY_SSE41__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_ConvertRGB24ToY_SSE41__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_ConvertRGB24ToY_SSE41__k7 = [2]uint64{0x603008f8f8f8f8f, 0x8f8f8f8f8f0f0c09}
var F_ConvertRGB24ToY_SSE41__k8 = [2]uint64{0x8f8f8f0d0a070401, 0x8f8f8f8f8f8f8f8f}
var F_ConvertRGB24ToY_SSE41__k9 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xe0b0805028f8f8f}
var F_ConvertRGB24ToY_SSE41__k10 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertRGB24ToY_SSE41__k11 = [2]uint64{0x412341c7412341c7, 0x412341c7412341c7}
var F_ConvertRGB24ToY_SSE41__k12 = [2]uint64{0x704018f8f8f8f8f, 0x8f8f8f8f8f8f0d0a}
var F_ConvertRGB24ToY_SSE41__k13 = [2]uint64{0x8f8f8f0e0b080502, 0x8f8f8f8f8f8f8f8f}
var F_ConvertRGB24ToY_SSE41__k14 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xf0c090603008f8f}
var F_ConvertRGB24ToY_SSE41__k15 = [2]uint64{0x1914400019144000, 0x1914400019144000}
var F_ConvertRGB24ToY_SSE41__k16 = [2]uint64{0x10800000108000, 0x10800000108000}
var F_ConvertRGB24ToY_SSE41__k17 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertRGB24ToY_SSE41__k18 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_ConvertRGB24ToY_SSE41__k19 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_ConvertRGB24ToY_SSE41__k20 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_ConvertRGB24ToY_SSE41__k21 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertRGB24ToY_SSE41__k22 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertRGB24ToY_SSE41__k23 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertRGB24ToY_SSE41__k24 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertRGB24ToY_SSE41__k25 = [2]uint64{0x41c7000041c7, 0x41c7000041c7}
var F_ConvertRGB24ToY_SSE41__k26 = [2]uint64{0x812300008123, 0x812300008123}
var F_ConvertRGB24ToY_SSE41__k27 = [2]uint64{0x191400001914, 0x191400001914}
var F_ConvertRGB24ToY_SSE41__k28 = [2]uint64{0xff000000ff, 0xff000000ff}

func F_ConvertRGBA32ToUV_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 int32
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v59 base.V128
	_ = v59
	var v61 base.V128
	_ = v61
	var v63 base.V128
	_ = v63
	var v65 base.V128
	_ = v65
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
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v81 base.V128
	_ = v81
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v87 base.V128
	_ = v87
	var v89 int32
	_ = v89
	var v92 base.V128
	_ = v92
	var v95 base.V128
	_ = v95
	var v103 base.V128
	_ = v103
	var v105 base.V128
	_ = v105
	var v107 base.V128
	_ = v107
	var v109 base.V128
	_ = v109
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v125 base.V128
	_ = v125
	var v127 base.V128
	_ = v127
	var v134 base.V128
	_ = v134
	var v136 base.V128
	_ = v136
	var v143 base.V128
	_ = v143
	var v146 base.V128
	_ = v146
	var v153 base.V128
	_ = v153
	var v156 base.V128
	_ = v156
	var v158 base.V128
	_ = v158
	var v184 base.V128
	_ = v184
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v196 int32
	_ = v196
	var v216 int32
	_ = v216
	var v239 int32
	_ = v239
	var v243 int32
	_ = v243
	var v249 int32
	_ = v249
	var v250 int32
	_ = v250
	var v251 int32
	_ = v251
	var v267 int32
	_ = v267
	var v269 int32
	_ = v269
	var v270 base.V128
	_ = v270
	var v272 int32
	_ = v272
	var v273 base.V128
	_ = v273
	var v275 int32
	_ = v275
	var v276 base.V128
	_ = v276
	var v278 int32
	_ = v278
	var v279 base.V128
	_ = v279
	var v280 base.V128
	_ = v280
	var v281 base.V128
	_ = v281
	var v285 int32
	_ = v285
	var v292 base.V128
	_ = v292
	var v295 base.V128
	_ = v295
	var v298 base.V128
	_ = v298
	var v301 base.V128
	_ = v301
	var v302 base.V128
	_ = v302
	var v303 base.V128
	_ = v303
	var v315 base.V128
	_ = v315
	var v318 base.V128
	_ = v318
	var v321 base.V128
	_ = v321
	var v324 base.V128
	_ = v324
	var v325 base.V128
	_ = v325
	var v326 base.V128
	_ = v326
	var v329 base.V128
	_ = v329
	var v333 base.V128
	_ = v333
	var v335 base.V128
	_ = v335
	var v346 base.V128
	_ = v346
	var v349 base.V128
	_ = v349
	var v352 base.V128
	_ = v352
	var v355 base.V128
	_ = v355
	var v356 base.V128
	_ = v356
	var v367 base.V128
	_ = v367
	var v370 base.V128
	_ = v370
	var v373 base.V128
	_ = v373
	var v376 base.V128
	_ = v376
	var v377 base.V128
	_ = v377
	var v389 base.V128
	_ = v389
	var v392 base.V128
	_ = v392
	var v395 base.V128
	_ = v395
	var v398 base.V128
	_ = v398
	var v399 base.V128
	_ = v399
	var v407 base.V128
	_ = v407
	var v412 base.V128
	_ = v412
	var v414 base.V128
	_ = v414
	var v446 int32
	_ = v446
	var v452 int32
	_ = v452
	var v453 int32
	_ = v453
	var v471 int32
	_ = v471
	var v475 int32
	_ = v475
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v491 int32
	_ = v491
	var v494 int32
	_ = v494
	var v498 int32
	_ = v498
	var v505 int32
	_ = v505
	var v506 int32
	_ = v506
	var v509 int32
	_ = v509
	var v510 int32
	_ = v510
	var v513 int32
	_ = v513
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v534 int32
	_ = v534
	var v536 int32
	_ = v536
	var v543 int32
	_ = v543
	v22 = l3 & int32(-16)
	if v22 < int32(1) {
		v194 = l0
		v195 = l1
		v196 = l2
	} else {
		v28 = l0
		v29 = l1
		v30 = l2
		for {
			v48 = int32(0)
			v49 = base.Simd_g_v128_load_rng(m, v28, v48, int32(0), int32(128))
			v50 = int32(16)
			v51 = base.Simd_g_v128_load_nc(m, v28, v50)
			v52 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k0)
			v53 = base.Simd_g_i8x16_shuffle2(v49, v51, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v54 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k3)
			v55 = base.Simd_g_i8x16_shuffle2(v49, v51, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v57 = base.Simd_g_i8x16_shuffle2(v53, v55, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v59 = base.Simd_g_v128_load_nc(m, v28, int32(32))
			v61 = base.Simd_g_v128_load_nc(m, v28, int32(48))
			v63 = base.Simd_g_i8x16_shuffle2(v59, v61, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v65 = base.Simd_g_i8x16_shuffle2(v59, v61, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v67 = base.Simd_g_i8x16_shuffle2(v63, v65, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v68 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k6)
			v69 = base.Simd_g_i8x16_shuffle2(v57, v67, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k7), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k8))
			v70 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k9)
			v71 = base.Simd_g_i8x16_shuffle2(v57, v67, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k10), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k11))
			v73 = base.Simd_g_i8x16_shuffle2(v69, v71, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v74 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k12)
			v81 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v53, v55, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v63, v65, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5)), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k7), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k8))
			v83 = base.Simd_g_i8x16_shuffle2(v71, v81, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v84 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k13)
			v87 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k14)
			v89 = int32(18)
			v92 = base.Simd_g_i8x16_shuffle2(v69, v71, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v95 = base.Simd_g_i8x16_shuffle2(v71, v81, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v103 = base.Simd_g_v128_load_nc(m, v28, int32(64))
			v105 = base.Simd_g_v128_load_nc(m, v28, int32(80))
			v107 = base.Simd_g_i8x16_shuffle2(v103, v105, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v109 = base.Simd_g_i8x16_shuffle2(v103, v105, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v111 = base.Simd_g_i8x16_shuffle2(v107, v109, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v113 = base.Simd_g_v128_load_nc(m, v28, int32(96))
			v115 = base.Simd_g_v128_load_nc(m, v28, int32(112))
			v117 = base.Simd_g_i8x16_shuffle2(v113, v115, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v119 = base.Simd_g_i8x16_shuffle2(v113, v115, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v121 = base.Simd_g_i8x16_shuffle2(v117, v119, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v123 = base.Simd_g_i8x16_shuffle2(v111, v121, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k7), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k8))
			v125 = base.Simd_g_i8x16_shuffle2(v111, v121, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k10), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k11))
			v127 = base.Simd_g_i8x16_shuffle2(v123, v125, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v134 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v107, v109, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5)), base.Simd_g_i8x16_shuffle2(v117, v119, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5)), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k7), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k8))
			v136 = base.Simd_g_i8x16_shuffle2(v125, v134, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k1), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k2))
			v143 = base.Simd_g_i8x16_shuffle2(v123, v125, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v146 = base.Simd_g_i8x16_shuffle2(v125, v134, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k4), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k5))
			v153 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v73, v74), base.Simd_g_i32x4_dot_i16x8_s(v83, v84)), v87), v89), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v92, v74), base.Simd_g_i32x4_dot_i16x8_s(v95, v84)), v87), v89)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v127, v74), base.Simd_g_i32x4_dot_i16x8_s(v136, v84)), v87), v89), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v143, v74), base.Simd_g_i32x4_dot_i16x8_s(v146, v84)), v87), v89)))
			base.Simd_g_v128_store(m, v29, v48, v153)
			v156 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k15)
			v158 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k16)
			v184 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v73, v156), base.Simd_g_i32x4_dot_i16x8_s(v83, v158)), v87), v89), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v92, v156), base.Simd_g_i32x4_dot_i16x8_s(v95, v158)), v87), v89)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v127, v156), base.Simd_g_i32x4_dot_i16x8_s(v136, v158)), v87), v89), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v143, v156), base.Simd_g_i32x4_dot_i16x8_s(v146, v158)), v87), v89)))
			base.Simd_g_v128_store(m, v30, v48, v184)
			v188 = v30 + v50
			v190 = v29 + v50
			v192 = v28 + int32(128)
			if base.Ui32(v192) < base.Ui32(l0+v22<<(uint(int32(3))%32)) {
				v28 = v192
				v29 = v190
				v30 = v188
				continue
			} else {
				break
			}
			break
		}
		v194 = v192
		v195 = v190
		v196 = v188
	}
	if v22 == l3 {
	} else {
		v216 = l3 & int32(15)
		if v216 < int32(1) {
		} else {
			if base.Ui32(int32(8)) <= base.Ui32(v216) {
				v239 = v216 & int32(2147483640)
				v243 = v194
				v249 = v239
				v250 = v196
				v251 = v195
				for {
					v267 = int32(8)
					v269 = int32(0)
					v270 = base.Simd_g_v128_load16_splat(m, v243, v269)
					v272 = int32(1)
					v273 = base.Simd_g_v128_load16_lane_l1(m, v243+v267, v269, v270)
					v275 = int32(2)
					v276 = base.Simd_g_v128_load16_lane_l2(m, v243+int32(16), v269, v273)
					v278 = int32(3)
					v279 = base.Simd_g_v128_load16_lane_l3(m, v243+int32(24), v269, v276)
					v280 = base.Simd_g_i32x4_extend_low_i16x8_u(v279)
					v281 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k15)
					v285 = int32(18)
					v292 = base.Simd_g_v128_load16_splat(m, v243+v275, v269)
					v295 = base.Simd_g_v128_load16_lane_l1(m, v243+int32(10), v269, v292)
					v298 = base.Simd_g_v128_load16_lane_l2(m, v243+v285, v269, v295)
					v301 = base.Simd_g_v128_load16_lane_l3(m, v243+int32(26), v269, v298)
					v302 = base.Simd_g_i32x4_extend_low_i16x8_u(v301)
					v303 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k17)
					v315 = base.Simd_g_v128_load16_splat(m, v243+int32(4), v269)
					v318 = base.Simd_g_v128_load16_lane_l1(m, v243+int32(12), v269, v315)
					v321 = base.Simd_g_v128_load16_lane_l2(m, v243+int32(20), v269, v318)
					v324 = base.Simd_g_v128_load16_lane_l3(m, v243+int32(28), v269, v321)
					v325 = base.Simd_g_i32x4_extend_low_i16x8_u(v324)
					v326 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k18)
					v329 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k14)
					v333 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k19)
					v335 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k20)
					v346 = base.Simd_g_v128_load16_splat(m, v243+int32(32), v269)
					v349 = base.Simd_g_v128_load16_lane_l1(m, v243+int32(40), v269, v346)
					v352 = base.Simd_g_v128_load16_lane_l2(m, v243+int32(48), v269, v349)
					v355 = base.Simd_g_v128_load16_lane_l3(m, v243+int32(56), v269, v352)
					v356 = base.Simd_g_i32x4_extend_low_i16x8_u(v355)
					v367 = base.Simd_g_v128_load16_splat(m, v243+int32(34), v269)
					v370 = base.Simd_g_v128_load16_lane_l1(m, v243+int32(42), v269, v367)
					v373 = base.Simd_g_v128_load16_lane_l2(m, v243+int32(50), v269, v370)
					v376 = base.Simd_g_v128_load16_lane_l3(m, v243+int32(58), v269, v373)
					v377 = base.Simd_g_i32x4_extend_low_i16x8_u(v376)
					v389 = base.Simd_g_v128_load16_splat(m, v243+int32(36), v269)
					v392 = base.Simd_g_v128_load16_lane_l1(m, v243+int32(44), v269, v389)
					v395 = base.Simd_g_v128_load16_lane_l2(m, v243+int32(52), v269, v392)
					v398 = base.Simd_g_v128_load16_lane_l3(m, v243+int32(60), v269, v395)
					v399 = base.Simd_g_i32x4_extend_low_i16x8_u(v398)
					v407 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k21)
					base.Simd_g_v128_store64_lane_l0(m, v250, v269, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v280, v281), base.Simd_g_i32x4_mul(v302, v303)), base.Simd_g_i32x4_mul(v325, v326)), v329), v285), v333), v335), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v356, v281), base.Simd_g_i32x4_mul(v377, v303)), base.Simd_g_i32x4_mul(v399, v326)), v329), v285), v333), v335), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k22), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k23)))
					v412 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k24)
					v414 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k25)
					base.Simd_g_v128_store64_lane_l0(m, v251, v269, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v280, v412), base.Simd_g_i32x4_mul(v302, v414)), base.Simd_g_i32x4_mul(v325, v281)), v329), v285), v333), v335), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v356, v412), base.Simd_g_i32x4_mul(v377, v414)), base.Simd_g_i32x4_mul(v399, v281)), v329), v285), v333), v335), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k22), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE2__k23)))
					v446 = v249 + int32(-8)
					if v446 != 0 {
						v243 = v243 + int32(64)
						v249 = v446
						v250 = v250 + v267
						v251 = v251 + v267
						continue
					} else {
						break
					}
					break
				}
				if v239 == v216 {
				} else {
					v452 = v194 + v239<<(uint(int32(3))%32)
					v453 = v239
					v471 = v195 + v453
					v475 = v452
					v478 = v196 + v453
					v479 = v216 - v453
					for {
						v491 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v475))))
						v494 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v475)+2)))
						v498 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v475)+4)))
						v505 = (v491*int32(-9719) + v494*int32(-19081) + v498*int32(_a_F_ConvertRGBA32ToUV_SSE2_0) + int32(33685504)) >> (uint(int32(18)) % 32)
						v506 = int32(0)
						if v506 < v505 {
							v509 = v505
						} else {
							v509 = v506
						}
						v510 = int32(255)
						if v509 < v510 {
							v513 = v509
						} else {
							v513 = v510
						}
						*(*uint8)(unsafe.Add(mBase, uint32(v471))) = uint8(v513)
						v526 = (v491*int32(_a_F_ConvertRGBA32ToUV_SSE2_0) + v494*int32(-24116) + v498*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
						v527 = int32(0)
						if v527 < v526 {
							v530 = v526
						} else {
							v530 = v527
						}
						v531 = int32(255)
						if v530 < v531 {
							v534 = v530
						} else {
							v534 = v531
						}
						*(*uint8)(unsafe.Add(mBase, uint32(v478))) = uint8(v534)
						v536 = int32(1)
						v543 = v479 + int32(-1)
						if v543 != 0 {
							v471 = v471 + v536
							v475 = v475 + int32(8)
							v478 = v478 + v536
							v479 = v543
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				v452 = v194
				v453 = int32(0)
				v471 = v195 + v453
				v475 = v452
				v478 = v196 + v453
				v479 = v216 - v453
				for {
					v491 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v475))))
					v494 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v475)+2)))
					v498 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v475)+4)))
					v505 = (v491*int32(-9719) + v494*int32(-19081) + v498*int32(_a_F_ConvertRGBA32ToUV_SSE2_0) + int32(33685504)) >> (uint(int32(18)) % 32)
					v506 = int32(0)
					if v506 < v505 {
						v509 = v505
					} else {
						v509 = v506
					}
					v510 = int32(255)
					if v509 < v510 {
						v513 = v509
					} else {
						v513 = v510
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v471))) = uint8(v513)
					v526 = (v491*int32(_a_F_ConvertRGBA32ToUV_SSE2_0) + v494*int32(-24116) + v498*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
					v527 = int32(0)
					if v527 < v526 {
						v530 = v526
					} else {
						v530 = v527
					}
					v531 = int32(255)
					if v530 < v531 {
						v534 = v530
					} else {
						v534 = v531
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v478))) = uint8(v534)
					v536 = int32(1)
					v543 = v479 + int32(-1)
					if v543 != 0 {
						v471 = v471 + v536
						v475 = v475 + int32(8)
						v478 = v478 + v536
						v479 = v543
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

var F_ConvertRGBA32ToUV_SSE2__k0 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertRGBA32ToUV_SSE2__k1 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertRGBA32ToUV_SSE2__k2 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertRGBA32ToUV_SSE2__k3 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertRGBA32ToUV_SSE2__k4 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertRGBA32ToUV_SSE2__k5 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertRGBA32ToUV_SSE2__k6 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertRGBA32ToUV_SSE2__k7 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertRGBA32ToUV_SSE2__k8 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertRGBA32ToUV_SSE2__k9 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_ConvertRGBA32ToUV_SSE2__k10 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertRGBA32ToUV_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertRGBA32ToUV_SSE2__k12 = [2]uint64{0xb577da09b577da09, 0xb577da09b577da09}
var F_ConvertRGBA32ToUV_SSE2__k13 = [2]uint64{0x7080000070800000, 0x7080000070800000}
var F_ConvertRGBA32ToUV_SSE2__k14 = [2]uint64{0x202000002020000, 0x202000002020000}
var F_ConvertRGBA32ToUV_SSE2__k15 = [2]uint64{0x708000007080, 0x708000007080}
var F_ConvertRGBA32ToUV_SSE2__k16 = [2]uint64{0xedb4a1ccedb4a1cc, 0xedb4a1ccedb4a1cc}
var F_ConvertRGBA32ToUV_SSE2__k17 = [2]uint64{0xffffa1ccffffa1cc, 0xffffa1ccffffa1cc}
var F_ConvertRGBA32ToUV_SSE2__k18 = [2]uint64{0xffffedb4ffffedb4, 0xffffedb4ffffedb4}
var F_ConvertRGBA32ToUV_SSE2__k19 = [2]uint64{0x0, 0x0}
var F_ConvertRGBA32ToUV_SSE2__k20 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ConvertRGBA32ToUV_SSE2__k21 = [2]uint64{0x1c1814100c080400, 0x0}
var F_ConvertRGBA32ToUV_SSE2__k22 = [2]uint64{0x808080800c080400, 0x0}
var F_ConvertRGBA32ToUV_SSE2__k23 = [2]uint64{0xc08040080808080, 0x8080808080808080}
var F_ConvertRGBA32ToUV_SSE2__k24 = [2]uint64{0xffffda09ffffda09, 0xffffda09ffffda09}
var F_ConvertRGBA32ToUV_SSE2__k25 = [2]uint64{0xffffb577ffffb577, 0xffffb577ffffb577}

func F_ConvertRGBA32ToUV_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 int32
	_ = v52
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v57 base.V128
	_ = v57
	var v59 base.V128
	_ = v59
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
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
	var v77 base.V128
	_ = v77
	var v79 base.V128
	_ = v79
	var v80 base.V128
	_ = v80
	var v83 base.V128
	_ = v83
	var v85 int32
	_ = v85
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v91 base.V128
	_ = v91
	var v99 base.V128
	_ = v99
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v105 base.V128
	_ = v105
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v110 base.V128
	_ = v110
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v115 base.V128
	_ = v115
	var v117 base.V128
	_ = v117
	var v119 base.V128
	_ = v119
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v133 base.V128
	_ = v133
	var v136 base.V128
	_ = v136
	var v143 base.V128
	_ = v143
	var v146 base.V128
	_ = v146
	var v148 base.V128
	_ = v148
	var v174 base.V128
	_ = v174
	var v178 int32
	_ = v178
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v206 int32
	_ = v206
	var v229 int32
	_ = v229
	var v233 int32
	_ = v233
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v257 int32
	_ = v257
	var v259 int32
	_ = v259
	var v260 base.V128
	_ = v260
	var v262 int32
	_ = v262
	var v263 base.V128
	_ = v263
	var v265 int32
	_ = v265
	var v266 base.V128
	_ = v266
	var v268 int32
	_ = v268
	var v269 base.V128
	_ = v269
	var v270 base.V128
	_ = v270
	var v271 base.V128
	_ = v271
	var v275 int32
	_ = v275
	var v282 base.V128
	_ = v282
	var v285 base.V128
	_ = v285
	var v288 base.V128
	_ = v288
	var v291 base.V128
	_ = v291
	var v292 base.V128
	_ = v292
	var v293 base.V128
	_ = v293
	var v305 base.V128
	_ = v305
	var v308 base.V128
	_ = v308
	var v311 base.V128
	_ = v311
	var v314 base.V128
	_ = v314
	var v315 base.V128
	_ = v315
	var v316 base.V128
	_ = v316
	var v319 base.V128
	_ = v319
	var v323 base.V128
	_ = v323
	var v325 base.V128
	_ = v325
	var v336 base.V128
	_ = v336
	var v339 base.V128
	_ = v339
	var v342 base.V128
	_ = v342
	var v345 base.V128
	_ = v345
	var v346 base.V128
	_ = v346
	var v357 base.V128
	_ = v357
	var v360 base.V128
	_ = v360
	var v363 base.V128
	_ = v363
	var v366 base.V128
	_ = v366
	var v367 base.V128
	_ = v367
	var v379 base.V128
	_ = v379
	var v382 base.V128
	_ = v382
	var v385 base.V128
	_ = v385
	var v388 base.V128
	_ = v388
	var v389 base.V128
	_ = v389
	var v397 base.V128
	_ = v397
	var v402 base.V128
	_ = v402
	var v404 base.V128
	_ = v404
	var v436 int32
	_ = v436
	var v442 int32
	_ = v442
	var v443 int32
	_ = v443
	var v461 int32
	_ = v461
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v481 int32
	_ = v481
	var v484 int32
	_ = v484
	var v488 int32
	_ = v488
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
	var v516 int32
	_ = v516
	var v517 int32
	_ = v517
	var v520 int32
	_ = v520
	var v521 int32
	_ = v521
	var v524 int32
	_ = v524
	var v526 int32
	_ = v526
	var v533 int32
	_ = v533
	v22 = l3 & int32(-16)
	if v22 < int32(1) {
		v184 = l0
		v185 = l1
		v186 = l2
	} else {
		v28 = l0
		v29 = l1
		v30 = l2
		for {
			v48 = int32(0)
			v49 = base.Simd_g_v128_load_rng(m, v28, v48, int32(0), int32(128))
			v50 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k0)
			v51 = base.Simd_g_i8x16_swizzle(v49, v50)
			v52 = int32(16)
			v53 = base.Simd_g_v128_load_nc(m, v28, v52)
			v54 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k1)
			v55 = base.Simd_g_i8x16_swizzle(v53, v54)
			v56 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k2)
			v57 = base.Simd_g_i8x16_shuffle2(v51, v55, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k3), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k4))
			v59 = base.Simd_g_v128_load_nc(m, v28, int32(32))
			v60 = base.Simd_g_i8x16_swizzle(v59, v50)
			v62 = base.Simd_g_v128_load_nc(m, v28, int32(48))
			v63 = base.Simd_g_i8x16_swizzle(v62, v54)
			v65 = base.Simd_g_i8x16_shuffle2(v60, v63, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k3), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k4))
			v66 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k5)
			v67 = base.Simd_g_i8x16_shuffle2(v57, v65, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k6), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k7))
			v68 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k8)
			v69 = base.Simd_g_i8x16_shuffle2(v57, v65, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k9), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k10))
			v70 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k11)
			v71 = base.Simd_g_i8x16_shuffle2(v67, v69, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k12), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k13))
			v72 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k14)
			v77 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(v55, v51), base.Simd_g_v128_or(v63, v60), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k9), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k10))
			v79 = base.Simd_g_i8x16_shuffle2(v69, v77, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k12), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k13))
			v80 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k15)
			v83 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k16)
			v85 = int32(18)
			v87 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k17)
			v88 = base.Simd_g_i8x16_shuffle2(v67, v69, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k18), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k19))
			v91 = base.Simd_g_i8x16_shuffle2(v69, v77, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k18), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k19))
			v99 = base.Simd_g_v128_load_nc(m, v28, int32(64))
			v100 = base.Simd_g_i8x16_swizzle(v99, v50)
			v102 = base.Simd_g_v128_load_nc(m, v28, int32(80))
			v103 = base.Simd_g_i8x16_swizzle(v102, v54)
			v105 = base.Simd_g_i8x16_shuffle2(v100, v103, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k3), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k4))
			v107 = base.Simd_g_v128_load_nc(m, v28, int32(96))
			v108 = base.Simd_g_i8x16_swizzle(v107, v50)
			v110 = base.Simd_g_v128_load_nc(m, v28, int32(112))
			v111 = base.Simd_g_i8x16_swizzle(v110, v54)
			v113 = base.Simd_g_i8x16_shuffle2(v108, v111, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k3), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k4))
			v115 = base.Simd_g_i8x16_shuffle2(v105, v113, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k6), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k7))
			v117 = base.Simd_g_i8x16_shuffle2(v105, v113, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k9), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k10))
			v119 = base.Simd_g_i8x16_shuffle2(v115, v117, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k12), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k13))
			v124 = base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(v103, v100), base.Simd_g_v128_or(v111, v108), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k9), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k10))
			v126 = base.Simd_g_i8x16_shuffle2(v117, v124, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k12), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k13))
			v133 = base.Simd_g_i8x16_shuffle2(v115, v117, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k18), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k19))
			v136 = base.Simd_g_i8x16_shuffle2(v117, v124, base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k18), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k19))
			v143 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v71, v72), base.Simd_g_i32x4_dot_i16x8_s(v79, v80)), v83), v85), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v88, v72), base.Simd_g_i32x4_dot_i16x8_s(v91, v80)), v83), v85)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v119, v72), base.Simd_g_i32x4_dot_i16x8_s(v126, v80)), v83), v85), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v133, v72), base.Simd_g_i32x4_dot_i16x8_s(v136, v80)), v83), v85)))
			base.Simd_g_v128_store(m, v29, v48, v143)
			v146 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k20)
			v148 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k21)
			v174 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v71, v146), base.Simd_g_i32x4_dot_i16x8_s(v79, v148)), v83), v85), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v88, v146), base.Simd_g_i32x4_dot_i16x8_s(v91, v148)), v83), v85)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v119, v146), base.Simd_g_i32x4_dot_i16x8_s(v126, v148)), v83), v85), base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_dot_i16x8_s(v133, v146), base.Simd_g_i32x4_dot_i16x8_s(v136, v148)), v83), v85)))
			base.Simd_g_v128_store(m, v30, v48, v174)
			v178 = v30 + v52
			v180 = v29 + v52
			v182 = v28 + int32(128)
			if base.Ui32(v182) < base.Ui32(l0+v22<<(uint(int32(3))%32)) {
				v28 = v182
				v29 = v180
				v30 = v178
				continue
			} else {
				break
			}
			break
		}
		v184 = v182
		v185 = v180
		v186 = v178
	}
	if v22 == l3 {
	} else {
		v206 = l3 & int32(15)
		if v206 < int32(1) {
		} else {
			if base.Ui32(int32(8)) <= base.Ui32(v206) {
				v229 = v206 & int32(2147483640)
				v233 = v184
				v239 = v229
				v240 = v186
				v241 = v185
				for {
					v257 = int32(8)
					v259 = int32(0)
					v260 = base.Simd_g_v128_load16_splat(m, v233, v259)
					v262 = int32(1)
					v263 = base.Simd_g_v128_load16_lane_l1(m, v233+v257, v259, v260)
					v265 = int32(2)
					v266 = base.Simd_g_v128_load16_lane_l2(m, v233+int32(16), v259, v263)
					v268 = int32(3)
					v269 = base.Simd_g_v128_load16_lane_l3(m, v233+int32(24), v259, v266)
					v270 = base.Simd_g_i32x4_extend_low_i16x8_u(v269)
					v271 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k20)
					v275 = int32(18)
					v282 = base.Simd_g_v128_load16_splat(m, v233+v265, v259)
					v285 = base.Simd_g_v128_load16_lane_l1(m, v233+int32(10), v259, v282)
					v288 = base.Simd_g_v128_load16_lane_l2(m, v233+v275, v259, v285)
					v291 = base.Simd_g_v128_load16_lane_l3(m, v233+int32(26), v259, v288)
					v292 = base.Simd_g_i32x4_extend_low_i16x8_u(v291)
					v293 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k22)
					v305 = base.Simd_g_v128_load16_splat(m, v233+int32(4), v259)
					v308 = base.Simd_g_v128_load16_lane_l1(m, v233+int32(12), v259, v305)
					v311 = base.Simd_g_v128_load16_lane_l2(m, v233+int32(20), v259, v308)
					v314 = base.Simd_g_v128_load16_lane_l3(m, v233+int32(28), v259, v311)
					v315 = base.Simd_g_i32x4_extend_low_i16x8_u(v314)
					v316 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k23)
					v319 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k16)
					v323 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k24)
					v325 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k25)
					v336 = base.Simd_g_v128_load16_splat(m, v233+int32(32), v259)
					v339 = base.Simd_g_v128_load16_lane_l1(m, v233+int32(40), v259, v336)
					v342 = base.Simd_g_v128_load16_lane_l2(m, v233+int32(48), v259, v339)
					v345 = base.Simd_g_v128_load16_lane_l3(m, v233+int32(56), v259, v342)
					v346 = base.Simd_g_i32x4_extend_low_i16x8_u(v345)
					v357 = base.Simd_g_v128_load16_splat(m, v233+int32(34), v259)
					v360 = base.Simd_g_v128_load16_lane_l1(m, v233+int32(42), v259, v357)
					v363 = base.Simd_g_v128_load16_lane_l2(m, v233+int32(50), v259, v360)
					v366 = base.Simd_g_v128_load16_lane_l3(m, v233+int32(58), v259, v363)
					v367 = base.Simd_g_i32x4_extend_low_i16x8_u(v366)
					v379 = base.Simd_g_v128_load16_splat(m, v233+int32(36), v259)
					v382 = base.Simd_g_v128_load16_lane_l1(m, v233+int32(44), v259, v379)
					v385 = base.Simd_g_v128_load16_lane_l2(m, v233+int32(52), v259, v382)
					v388 = base.Simd_g_v128_load16_lane_l3(m, v233+int32(60), v259, v385)
					v389 = base.Simd_g_i32x4_extend_low_i16x8_u(v388)
					v397 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k26)
					base.Simd_g_v128_store64_lane_l0(m, v240, v259, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v270, v271), base.Simd_g_i32x4_mul(v292, v293)), base.Simd_g_i32x4_mul(v315, v316)), v319), v275), v323), v325), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v346, v271), base.Simd_g_i32x4_mul(v367, v293)), base.Simd_g_i32x4_mul(v389, v316)), v319), v275), v323), v325), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k27), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k28)))
					v402 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k29)
					v404 = base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k30)
					base.Simd_g_v128_store64_lane_l0(m, v241, v259, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v270, v402), base.Simd_g_i32x4_mul(v292, v404)), base.Simd_g_i32x4_mul(v315, v271)), v319), v275), v323), v325), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v346, v402), base.Simd_g_i32x4_mul(v367, v404)), base.Simd_g_i32x4_mul(v389, v271)), v319), v275), v323), v325), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k27), base.Simd_g_const(&F_ConvertRGBA32ToUV_SSE41__k28)))
					v436 = v239 + int32(-8)
					if v436 != 0 {
						v233 = v233 + int32(64)
						v239 = v436
						v240 = v240 + v257
						v241 = v241 + v257
						continue
					} else {
						break
					}
					break
				}
				if v229 == v206 {
				} else {
					v442 = v184 + v229<<(uint(int32(3))%32)
					v443 = v229
					v461 = v185 + v443
					v465 = v442
					v468 = v186 + v443
					v469 = v206 - v443
					for {
						v481 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v465))))
						v484 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v465)+2)))
						v488 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v465)+4)))
						v495 = (v481*int32(-9719) + v484*int32(-19081) + v488*int32(_a_F_ConvertRGBA32ToUV_SSE41_0) + int32(33685504)) >> (uint(int32(18)) % 32)
						v496 = int32(0)
						if v496 < v495 {
							v499 = v495
						} else {
							v499 = v496
						}
						v500 = int32(255)
						if v499 < v500 {
							v503 = v499
						} else {
							v503 = v500
						}
						*(*uint8)(unsafe.Add(mBase, uint32(v461))) = uint8(v503)
						v516 = (v481*int32(_a_F_ConvertRGBA32ToUV_SSE41_0) + v484*int32(-24116) + v488*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
						v517 = int32(0)
						if v517 < v516 {
							v520 = v516
						} else {
							v520 = v517
						}
						v521 = int32(255)
						if v520 < v521 {
							v524 = v520
						} else {
							v524 = v521
						}
						*(*uint8)(unsafe.Add(mBase, uint32(v468))) = uint8(v524)
						v526 = int32(1)
						v533 = v469 + int32(-1)
						if v533 != 0 {
							v461 = v461 + v526
							v465 = v465 + int32(8)
							v468 = v468 + v526
							v469 = v533
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				v442 = v184
				v443 = int32(0)
				v461 = v185 + v443
				v465 = v442
				v468 = v186 + v443
				v469 = v206 - v443
				for {
					v481 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v465))))
					v484 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v465)+2)))
					v488 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v465)+4)))
					v495 = (v481*int32(-9719) + v484*int32(-19081) + v488*int32(_a_F_ConvertRGBA32ToUV_SSE41_0) + int32(33685504)) >> (uint(int32(18)) % 32)
					v496 = int32(0)
					if v496 < v495 {
						v499 = v495
					} else {
						v499 = v496
					}
					v500 = int32(255)
					if v499 < v500 {
						v503 = v499
					} else {
						v503 = v500
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v461))) = uint8(v503)
					v516 = (v481*int32(_a_F_ConvertRGBA32ToUV_SSE41_0) + v484*int32(-24116) + v488*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
					v517 = int32(0)
					if v517 < v516 {
						v520 = v516
					} else {
						v520 = v517
					}
					v521 = int32(255)
					if v520 < v521 {
						v524 = v520
					} else {
						v524 = v521
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v468))) = uint8(v524)
					v526 = int32(1)
					v533 = v469 + int32(-1)
					if v533 != 0 {
						v461 = v461 + v526
						v465 = v465 + int32(8)
						v468 = v468 + v526
						v469 = v533
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

var F_ConvertRGBA32ToUV_SSE41__k0 = [2]uint64{0xb0a030209080100, 0x8f8f8f8f0d0c0504}
var F_ConvertRGBA32ToUV_SSE41__k1 = [2]uint64{0xb0a030209080100, 0xd0c05048f8f8f8f}
var F_ConvertRGBA32ToUV_SSE41__k2 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_ConvertRGBA32ToUV_SSE41__k3 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_ConvertRGBA32ToUV_SSE41__k4 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_ConvertRGBA32ToUV_SSE41__k5 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_ConvertRGBA32ToUV_SSE41__k6 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ConvertRGBA32ToUV_SSE41__k7 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ConvertRGBA32ToUV_SSE41__k8 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_ConvertRGBA32ToUV_SSE41__k9 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ConvertRGBA32ToUV_SSE41__k10 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_ConvertRGBA32ToUV_SSE41__k11 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_ConvertRGBA32ToUV_SSE41__k12 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_ConvertRGBA32ToUV_SSE41__k13 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_ConvertRGBA32ToUV_SSE41__k14 = [2]uint64{0xb577da09b577da09, 0xb577da09b577da09}
var F_ConvertRGBA32ToUV_SSE41__k15 = [2]uint64{0x7080000070800000, 0x7080000070800000}
var F_ConvertRGBA32ToUV_SSE41__k16 = [2]uint64{0x202000002020000, 0x202000002020000}
var F_ConvertRGBA32ToUV_SSE41__k17 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_ConvertRGBA32ToUV_SSE41__k18 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_ConvertRGBA32ToUV_SSE41__k19 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_ConvertRGBA32ToUV_SSE41__k20 = [2]uint64{0x708000007080, 0x708000007080}
var F_ConvertRGBA32ToUV_SSE41__k21 = [2]uint64{0xedb4a1ccedb4a1cc, 0xedb4a1ccedb4a1cc}
var F_ConvertRGBA32ToUV_SSE41__k22 = [2]uint64{0xffffa1ccffffa1cc, 0xffffa1ccffffa1cc}
var F_ConvertRGBA32ToUV_SSE41__k23 = [2]uint64{0xffffedb4ffffedb4, 0xffffedb4ffffedb4}
var F_ConvertRGBA32ToUV_SSE41__k24 = [2]uint64{0x0, 0x0}
var F_ConvertRGBA32ToUV_SSE41__k25 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ConvertRGBA32ToUV_SSE41__k26 = [2]uint64{0x1c1814100c080400, 0x0}
var F_ConvertRGBA32ToUV_SSE41__k27 = [2]uint64{0x808080800c080400, 0x0}
var F_ConvertRGBA32ToUV_SSE41__k28 = [2]uint64{0xc08040080808080, 0x8080808080808080}
var F_ConvertRGBA32ToUV_SSE41__k29 = [2]uint64{0xffffda09ffffda09, 0xffffda09ffffda09}
var F_ConvertRGBA32ToUV_SSE41__k30 = [2]uint64{0xffffb577ffffb577, 0xffffb577ffffb577}

func F_Copy16x8_C(m *base.Module, l0 int32, l1 int32) {
	var v3 int32
	_ = v3
	var v4 base.V128
	_ = v4
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v11 int32
	_ = v11
	var v12 base.V128
	_ = v12
	var v15 int32
	_ = v15
	var v16 base.V128
	_ = v16
	var v19 int32
	_ = v19
	var v20 base.V128
	_ = v20
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v27 int32
	_ = v27
	var v28 base.V128
	_ = v28
	var v31 int32
	_ = v31
	var v32 base.V128
	_ = v32
	v3 = int32(0)
	v4 = base.Simd_g_v128_load(m, l0, v3)
	base.Simd_g_v128_store(m, l1, v3, v4)
	v7 = int32(32)
	v8 = base.Simd_g_v128_load(m, l0, v7)
	base.Simd_g_v128_store(m, l1, v7, v8)
	v11 = int32(64)
	v12 = base.Simd_g_v128_load(m, l0, v11)
	base.Simd_g_v128_store(m, l1, v11, v12)
	v15 = int32(96)
	v16 = base.Simd_g_v128_load(m, l0, v15)
	base.Simd_g_v128_store(m, l1, v15, v16)
	v19 = int32(128)
	v20 = base.Simd_g_v128_load(m, l0, v19)
	base.Simd_g_v128_store(m, l1, v19, v20)
	v23 = int32(160)
	v24 = base.Simd_g_v128_load(m, l0, v23)
	base.Simd_g_v128_store(m, l1, v23, v24)
	v27 = int32(192)
	v28 = base.Simd_g_v128_load(m, l0, v27)
	base.Simd_g_v128_store(m, l1, v27, v28)
	v31 = int32(224)
	v32 = base.Simd_g_v128_load(m, l0, v31)
	base.Simd_g_v128_store(m, l1, v31, v32)
	return
}

//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_MaxDiffsForRow(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v28 base.V128
	_ = v28
	var v31 base.V128
	_ = v31
	var v41 base.V128
	_ = v41
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v58 int32
	_ = v58
	var __phi58 int32
	_ = __phi58
	var v59 int32
	_ = v59
	var __phi59 int32
	_ = __phi59
	var v66 int32
	_ = v66
	var __phi66 int32
	_ = __phi66
	var v67 int32
	_ = v67
	var __phi67 int32
	_ = __phi67
	var v68 int32
	_ = v68
	var __phi68 int32
	_ = __phi68
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v103 int32
	_ = v103
	var v118 int32
	_ = v118
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v153 int32
	_ = v153
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v170 int32
	_ = v170
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v197 int32
	_ = v197
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v203 int32
	_ = v203
	var v206 int32
	_ = v206
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v217 int32
	_ = v217
	var v219 int32
	_ = v219
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v241 int32
	_ = v241
	var v243 int32
	_ = v243
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v255 int32
	_ = v255
	var v257 int32
	_ = v257
	var v259 int32
	_ = v259
	var v261 int32
	_ = v261
	var v263 int32
	_ = v263
	var v266 int32
	_ = v266
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v275 int32
	_ = v275
	var v277 int32
	_ = v277
	var v279 int32
	_ = v279
	var v281 int32
	_ = v281
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v301 int32
	_ = v301
	var v303 int32
	_ = v303
	var v305 int32
	_ = v305
	var v310 int32
	_ = v310
	if l0 < int32(3) {
	} else {
		v23 = int32(0)
		v24 = base.Simd_g_v128_load64_zero(m, l2, v23)
		if l4 == v23 {
			v41 = v24
		} else {
			v28 = base.Simd_g_i32x4_shr_u(v24, int32(8))
			v31 = base.Simd_g_const(&F_MaxDiffsForRow__k0)
			v41 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v28, base.Simd_g_const(&F_MaxDiffsForRow__k1)), base.Simd_g_v128_and(v24, v31)), base.Simd_g_i32x4_shl(v28, int32(16))), v31), base.Simd_g_v128_and(v24, base.Simd_g_const(&F_MaxDiffsForRow__k2)))
		}
		v46 = int32(1)
		v50 = int32(0)
		v52 = l1 << (uint(int32(2)) % 32)
		__phi58 = base.Simd_g_i32x4_extract_lane_l1(v41)
		__phi59 = base.Simd_g_i32x4_extract_lane_l0(v41)
		__phi66 = l2 + int32(4)
		__phi67 = l3 + v46
		__phi68 = l0 + int32(-2)
		v58 = __phi58
		v59 = __phi59
		v66 = __phi66
		v67 = __phi67
		v68 = __phi68
		for {
			v79 = v66 + int32(4)
			v80 = *(*int32)(unsafe.Add(mBase, uint32(v79)))
			v82 = *(*int32)(unsafe.Add(mBase, uint32(v66+v52)))
			v84 = *(*int32)(unsafe.Add(mBase, uint32(v66+(v50-v52))))
			if l4 == int32(0) {
				v132 = v80
				v133 = v82
				v134 = v84
			} else {
				v87 = int32(8)
				v88 = int32(base.Ui32(v80) >> (uint(v87) % 32))
				v89 = int32(255)
				v91 = int32(16711935)
				v94 = int32(16)
				v99 = int32(-16711936)
				v103 = int32(base.Ui32(v82) >> (uint(v87) % 32))
				v118 = int32(base.Ui32(v84) >> (uint(v87) % 32))
				v132 = (v88&v89+v80&v91+v88<<(uint(v94)%32))&v91 | v80&v99
				v133 = (v103&v89+v82&v91+v103<<(uint(v94)%32))&v91 | v82&v99
				v134 = (v118&v89+v84&v91+v118<<(uint(v94)%32))&v91 | v84&v99
			}
			v136 = int32(24)
			v137 = int32(base.Ui32(v58) >> (uint(v136) % 32))
			v140 = v137 - int32(base.Ui32(v134)>>(uint(v136)%32))
			v141 = int32(31)
			v142 = v140 >> (uint(v141) % 32)
			v144 = v140 ^ v142 - v142
			v145 = int32(16)
			v147 = int32(255)
			v148 = int32(base.Ui32(v58)>>(uint(v145)%32)) & v147
			v153 = v148 - int32(base.Ui32(v134)>>(uint(v145)%32))&v147
			v155 = v153 >> (uint(v141) % 32)
			v157 = v153 ^ v155 - v155
			if base.Ui32(v157) < base.Ui32(v144) {
				v159 = v144
			} else {
				v159 = v157
			}
			v160 = int32(8)
			v162 = int32(255)
			v163 = int32(base.Ui32(v58)>>(uint(v160)%32)) & v162
			v168 = v163 - int32(base.Ui32(v134)>>(uint(v160)%32))&v162
			v169 = int32(31)
			v170 = v168 >> (uint(v169) % 32)
			v172 = v168 ^ v170 - v170
			v174 = v58 & v162
			v177 = v174 - v134&v162
			v179 = v177 >> (uint(v169) % 32)
			v181 = v177 ^ v179 - v179
			if base.Ui32(v181) < base.Ui32(v172) {
				v183 = v172
			} else {
				v183 = v181
			}
			if base.Ui32(v183) < base.Ui32(v159) {
				v185 = v159
			} else {
				v185 = v183
			}
			v188 = v137 - int32(base.Ui32(v133)>>(uint(int32(24))%32))
			v189 = int32(31)
			v190 = v188 >> (uint(v189) % 32)
			v192 = v188 ^ v190 - v190
			v197 = v148 - int32(base.Ui32(v133)>>(uint(int32(16))%32))&int32(255)
			v199 = v197 >> (uint(v189) % 32)
			v201 = v197 ^ v199 - v199
			if base.Ui32(v201) < base.Ui32(v192) {
				v203 = v192
			} else {
				v203 = v201
			}
			v206 = int32(255)
			v208 = v163 - int32(base.Ui32(v133)>>(uint(int32(8))%32))&v206
			v209 = int32(31)
			v210 = v208 >> (uint(v209) % 32)
			v212 = v208 ^ v210 - v210
			v215 = v174 - v133&v206
			v217 = v215 >> (uint(v209) % 32)
			v219 = v215 ^ v217 - v217
			if base.Ui32(v219) < base.Ui32(v212) {
				v221 = v212
			} else {
				v221 = v219
			}
			if base.Ui32(v221) < base.Ui32(v203) {
				v223 = v203
			} else {
				v223 = v221
			}
			if base.Ui32(v223) < base.Ui32(v185) {
				v225 = v185
			} else {
				v225 = v223
			}
			v228 = v137 - int32(base.Ui32(v59)>>(uint(int32(24))%32))
			v229 = int32(31)
			v230 = v228 >> (uint(v229) % 32)
			v232 = v228 ^ v230 - v230
			v237 = v148 - int32(base.Ui32(v59)>>(uint(int32(16))%32))&int32(255)
			v239 = v237 >> (uint(v229) % 32)
			v241 = v237 ^ v239 - v239
			if base.Ui32(v241) < base.Ui32(v232) {
				v243 = v232
			} else {
				v243 = v241
			}
			v246 = int32(255)
			v248 = v163 - int32(base.Ui32(v59)>>(uint(int32(8))%32))&v246
			v249 = int32(31)
			v250 = v248 >> (uint(v249) % 32)
			v252 = v248 ^ v250 - v250
			v255 = v174 - v59&v246
			v257 = v255 >> (uint(v249) % 32)
			v259 = v255 ^ v257 - v257
			if base.Ui32(v259) < base.Ui32(v252) {
				v261 = v252
			} else {
				v261 = v259
			}
			if base.Ui32(v261) < base.Ui32(v243) {
				v263 = v243
			} else {
				v263 = v261
			}
			v266 = v137 - int32(base.Ui32(v132)>>(uint(int32(24))%32))
			v267 = int32(31)
			v268 = v266 >> (uint(v267) % 32)
			v270 = v266 ^ v268 - v268
			v275 = v148 - int32(base.Ui32(v132)>>(uint(int32(16))%32))&int32(255)
			v277 = v275 >> (uint(v267) % 32)
			v279 = v275 ^ v277 - v277
			if base.Ui32(v279) < base.Ui32(v270) {
				v281 = v270
			} else {
				v281 = v279
			}
			v284 = int32(255)
			v286 = v163 - int32(base.Ui32(v132)>>(uint(int32(8))%32))&v284
			v287 = int32(31)
			v288 = v286 >> (uint(v287) % 32)
			v290 = v286 ^ v288 - v288
			v293 = v174 - v132&v284
			v295 = v293 >> (uint(v287) % 32)
			v297 = v293 ^ v295 - v295
			if base.Ui32(v297) < base.Ui32(v290) {
				v299 = v290
			} else {
				v299 = v297
			}
			if base.Ui32(v299) < base.Ui32(v281) {
				v301 = v281
			} else {
				v301 = v299
			}
			if base.Ui32(v301) < base.Ui32(v263) {
				v303 = v263
			} else {
				v303 = v301
			}
			if base.Ui32(v303) < base.Ui32(v225) {
				v305 = v225
			} else {
				v305 = v303
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v67))) = uint8(v305)
			v310 = v68 + int32(-1)
			if v310 != 0 {
				__phi58 = v132
				__phi59 = v58
				__phi66 = v79
				__phi67 = v67 + int32(1)
				__phi68 = v310
				v58 = __phi58
				v59 = __phi59
				v66 = __phi66
				v67 = __phi67
				v68 = __phi68
				continue
			} else {
				break
			}
			break
		}
	}
	return
}

var F_MaxDiffsForRow__k0 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_MaxDiffsForRow__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_MaxDiffsForRow__k2 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}

func F_Mean16x4_SSE2(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v7 base.V128
	_ = v7
	var v8 base.V128
	_ = v8
	var v10 int32
	_ = v10
	var v11 base.V128
	_ = v11
	var v14 int32
	_ = v14
	var v18 base.V128
	_ = v18
	var v25 base.V128
	_ = v25
	var v33 base.V128
	_ = v33
	v7 = base.Simd_g_v128_load_rng(m, l0, int32(32), int32(0), int32(112))
	v8 = base.Simd_g_const(&F_Mean16x4_SSE2__k0)
	v10 = int32(0)
	v11 = base.Simd_g_v128_load_nc(m, l0, v10)
	v14 = int32(8)
	v18 = base.Simd_g_v128_load_nc(m, l0, int32(64))
	v25 = base.Simd_g_v128_load_nc(m, l0, int32(96))
	v33 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v7, v8), base.Simd_g_v128_and(v11, v8)), base.Simd_g_i16x8_shr_u(v11, v14)), base.Simd_g_v128_and(v18, v8)), base.Simd_g_i16x8_shr_u(v7, v14)), base.Simd_g_v128_and(v25, v8)), base.Simd_g_i16x8_shr_u(v18, v14)), base.Simd_g_i16x8_shr_u(v25, v14))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = base.Simd_g_i16x8_extract_lane_u_l6(v33) + base.Simd_g_i16x8_extract_lane_u_l7(v33)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = base.Simd_g_i16x8_extract_lane_u_l4(v33) + base.Simd_g_i16x8_extract_lane_u_l5(v33)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = base.Simd_g_i16x8_extract_lane_u_l2(v33) + base.Simd_g_i16x8_extract_lane_u_l3(v33)
	*(*int32)(unsafe.Add(mBase, uint32(l1))) = base.Simd_g_i16x8_extract_lane_u_l0(v33) + base.Simd_g_i16x8_extract_lane_u_l1(v33)
	return
}

var F_Mean16x4_SSE2__k0 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_MultARGBRow_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v54 int32
	_ = v54
	var v59 int32
	_ = v59
	var v65 int32
	_ = v65
	var v70 int32
	_ = v70
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v95 base.V128
	_ = v95
	var v97 int32
	_ = v97
	var v98 base.V128
	_ = v98
	var v99 base.V128
	_ = v99
	var v102 base.V128
	_ = v102
	var v114 int32
	_ = v114
	var v116 base.V128
	_ = v116
	var v117 int32
	_ = v117
	var v124 int32
	_ = v124
	var v128 int32
	_ = v128
	var v131 base.V128
	_ = v131
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v146 base.V128
	_ = v146
	var v154 int32
	_ = v154
	var v158 int32
	_ = v158
	var v161 base.V128
	_ = v161
	var v163 base.V128
	_ = v163
	var v166 base.V128
	_ = v166
	var v172 int32
	_ = v172
	var v177 int32
	_ = v177
	var v194 base.V128
	_ = v194
	var v195 base.V128
	_ = v195
	var v196 int32
	_ = v196
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v214 int32
	_ = v214
	var v225 int32
	_ = v225
	var v236 int32
	_ = v236
	var v241 int32
	_ = v241
	var v247 int32
	_ = v247
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v259 int32
	_ = v259
	var v260 base.V128
	_ = v260
	var v261 int32
	_ = v261
	var v264 base.V128
	_ = v264
	var v265 base.V128
	_ = v265
	var v268 base.V128
	_ = v268
	var v272 base.V128
	_ = v272
	var v275 int32
	_ = v275
	var v280 int32
	_ = v280
	var v297 base.V128
	_ = v297
	var v298 base.V128
	_ = v298
	var v305 int32
	_ = v305
	var v308 int32
	_ = v308
	var v317 int32
	_ = v317
	var v328 int32
	_ = v328
	var v339 int32
	_ = v339
	var v344 int32
	_ = v344
	var v349 int32
	_ = v349
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v370 int32
	_ = v370
	var v376 int32
	_ = v376
	var v379 int32
	_ = v379
	var v380 int32
	_ = v380
	var v383 int32
	_ = v383
	var v390 int32
	_ = v390
	var v397 int32
	_ = v397
	var v414 int32
	_ = v414
	var v420 int32
	_ = v420
	var v424 int32
	_ = v424
	var v436 int32
	_ = v436
	var v437 int32
	_ = v437
	var v445 int32
	_ = v445
	var v451 int32
	_ = v451
	var v452 int32
	_ = v452
	var v454 int32
	_ = v454
	var v455 int32
	_ = v455
	var v458 int32
	_ = v458
	var v465 int32
	_ = v465
	var v472 int32
	_ = v472
	var v489 int32
	_ = v489
	var v495 int32
	_ = v495
	v4 = int32(0)
	if l2 != 0 {
		v59 = v4
	} else {
		if l1 < int32(2) {
			v59 = v4
		} else {
			v18 = int32(0)
			v19 = l0
			for {
				v23 = int32(0)
				v24 = base.Simd_g_v128_load64_zero(m, v19, v23)
				v25 = base.Simd_g_const(&F_MultARGBRow_SSE2__k0)
				v27 = base.Simd_g_i8x16_shuffle2(v24, v25, base.Simd_g_const(&F_MultARGBRow_SSE2__k1), base.Simd_g_const(&F_MultARGBRow_SSE2__k2))
				v36 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_mul(base.Simd_g_i8x16_swizzle_c(base.Simd_g_i8x16_swizzle_c(base.Simd_g_v128_or(v27, base.Simd_g_const(&F_MultARGBRow_SSE2__k3)), base.Simd_g_const(&F_MultARGBRow_SSE2__k4)), base.Simd_g_const(&F_MultARGBRow_SSE2__k5)), v27), base.Simd_g_const(&F_MultARGBRow_SSE2__k6))
				v38 = base.Simd_g_const(&F_MultARGBRow_SSE2__k7)
				base.Simd_g_v128_store64_lane_l0(m, v19, v23, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v36), v38), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v36), v38), base.Simd_g_const(&F_MultARGBRow_SSE2__k8), base.Simd_g_const(&F_MultARGBRow_SSE2__k9)), v25))
				v54 = v18 + int32(2)
				if v18+int32(4) <= l1 {
					v18 = v54
					v19 = v19 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v59 = v54
		}
	}
	v65 = l1 - v59
	if v65 < int32(1) {
	} else {
		v70 = l0 + v59<<(uint(int32(2))%32)
		if v65 < int32(1) {
		} else {
			if l2 == int32(0) {
				if base.Ui32(v65) < base.Ui32(int32(4)) {
					v349 = int32(0)
					v361 = v65 - v349
					v362 = v70 + v349<<(uint(int32(2))%32)
					for {
						v370 = *(*int32)(unsafe.Add(mBase, uint32(v362)))
						if base.Ui32(int32(-16777217)) < base.Ui32(v370) {
						} else {
							if base.Ui32(v370) < base.Ui32(int32(16777216)) {
								v414 = int32(0)
							} else {
								v376 = int32(24)
								v379 = int32(base.Ui32(v370)>>(uint(v376)%32)) * int32(_a_F_MultARGBRow_SSE2_0)
								v380 = int32(255)
								v383 = int32(8388608)
								v390 = int32(8)
								v397 = int32(16)
								v414 = int32(base.Ui32(v379*(v370&v380)+v383)>>(uint(v376)%32)) | v370&int32(-16777216) | int32(base.Ui32(v379*(int32(base.Ui32(v370)>>(uint(v390)%32))&v380)+v383)>>(uint(v397)%32))&int32(_a_F_MultARGBRow_SSE2_1) | int32(base.Ui32(v379*(int32(base.Ui32(v370)>>(uint(v397)%32))&v380)+v383)>>(uint(v390)%32))&int32(16711680)
							}
							*(*int32)(unsafe.Add(mBase, uint32(v362))) = v414
						}
						v420 = v361 + int32(-1)
						if v420 != 0 {
							v361 = v420
							v362 = v362 + int32(4)
							continue
						} else {
							break
						}
						break
					}
				} else {
					v247 = v65 & int32(2147483644)
					v250 = v70
					v252 = v247
					for {
						v259 = int32(0)
						v260 = base.Simd_g_v128_load(m, v250, v259)
						v261 = int32(24)
						v264 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_u(v260, v261), base.Simd_g_const(&F_MultARGBRow_SSE2__k10))
						v265 = base.Simd_g_const(&F_MultARGBRow_SSE2__k11)
						v268 = base.Simd_g_const(&F_MultARGBRow_SSE2__k12)
						v272 = base.Simd_g_const(&F_MultARGBRow_SSE2__k13)
						v275 = int32(8)
						v280 = int32(16)
						v297 = base.Simd_g_v128_bitselect(base.Simd_g_const(&F_MultARGBRow_SSE2__k0), base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v264, base.Simd_g_v128_and(v260, v265)), v268), v261), base.Simd_g_v128_and(v260, v272)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v264, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v260, v275), v265)), v268), v280), base.Simd_g_const(&F_MultARGBRow_SSE2__k14))), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v264, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v260, v280), v265)), v268), v275), base.Simd_g_const(&F_MultARGBRow_SSE2__k15))), base.Simd_g_i32x4_lt_u(v260, base.Simd_g_const(&F_MultARGBRow_SSE2__k16)))
						v298 = base.Simd_g_i32x4_lt_u(v260, v272)
						if base.Simd_g_i32x4_extract_lane_l0(v298)&int32(1) == v259 {
						} else {
							v305 = int32(0)
							base.Simd_g_v128_store32_lane_l0(m, v250, v305, v297)
						}
						v308 = int32(1)
						if base.Simd_g_i32x4_extract_lane_l1(v298)&v308 == int32(0) {
						} else {
							v317 = int32(1)
							base.Simd_g_v128_store32_lane_l1(m, v250+int32(4), int32(0), v297)
						}
						if base.Simd_g_i32x4_extract_lane_l2(v298)&int32(1) == int32(0) {
						} else {
							v328 = int32(2)
							base.Simd_g_v128_store32_lane_l2(m, v250+int32(8), int32(0), v297)
						}
						if base.Simd_g_i32x4_extract_lane_l3(v298)&int32(1) == int32(0) {
						} else {
							v339 = int32(3)
							base.Simd_g_v128_store32_lane_l3(m, v250+int32(12), int32(0), v297)
						}
						v344 = v252 + int32(-4)
						if v344 != 0 {
							v250 = v250 + int32(16)
							v252 = v344
							continue
						} else {
							break
						}
						break
					}
					if v247 == v65 {
					} else {
						v349 = v247
						v361 = v65 - v349
						v362 = v70 + v349<<(uint(int32(2))%32)
						for {
							v370 = *(*int32)(unsafe.Add(mBase, uint32(v362)))
							if base.Ui32(int32(-16777217)) < base.Ui32(v370) {
							} else {
								if base.Ui32(v370) < base.Ui32(int32(16777216)) {
									v414 = int32(0)
								} else {
									v376 = int32(24)
									v379 = int32(base.Ui32(v370)>>(uint(v376)%32)) * int32(_a_F_MultARGBRow_SSE2_0)
									v380 = int32(255)
									v383 = int32(8388608)
									v390 = int32(8)
									v397 = int32(16)
									v414 = int32(base.Ui32(v379*(v370&v380)+v383)>>(uint(v376)%32)) | v370&int32(-16777216) | int32(base.Ui32(v379*(int32(base.Ui32(v370)>>(uint(v390)%32))&v380)+v383)>>(uint(v397)%32))&int32(_a_F_MultARGBRow_SSE2_1) | int32(base.Ui32(v379*(int32(base.Ui32(v370)>>(uint(v397)%32))&v380)+v383)>>(uint(v390)%32))&int32(16711680)
								}
								*(*int32)(unsafe.Add(mBase, uint32(v362))) = v414
							}
							v420 = v361 + int32(-1)
							if v420 != 0 {
								v361 = v420
								v362 = v362 + int32(4)
								continue
							} else {
								break
							}
							break
						}
					}
				}
			} else {
				if base.Ui32(v65) <= base.Ui32(int32(3)) {
					v424 = int32(0)
					v436 = v65 - v424
					v437 = v70 + v424<<(uint(int32(2))%32)
					for {
						v445 = *(*int32)(unsafe.Add(mBase, uint32(v437)))
						if base.Ui32(int32(-16777217)) < base.Ui32(v445) {
						} else {
							if base.Ui32(v445) < base.Ui32(int32(16777216)) {
								v489 = int32(0)
							} else {
								v451 = int32(-16777216)
								v452 = int32(24)
								v454 = base.I32_div_u_s(v451, int32(base.Ui32(v445)>>(uint(v452)%32)))
								v455 = int32(255)
								v458 = int32(8388608)
								v465 = int32(8)
								v472 = int32(16)
								v489 = int32(base.Ui32(v454*(v445&v455)+v458)>>(uint(v452)%32)) | v445&v451 | int32(base.Ui32(v454*(int32(base.Ui32(v445)>>(uint(v465)%32))&v455)+v458)>>(uint(v472)%32))&int32(_a_F_MultARGBRow_SSE2_1) | int32(base.Ui32(v454*(int32(base.Ui32(v445)>>(uint(v472)%32))&v455)+v458)>>(uint(v465)%32))&int32(16711680)
							}
							*(*int32)(unsafe.Add(mBase, uint32(v437))) = v489
						}
						v495 = v436 + int32(-1)
						if v495 != 0 {
							v436 = v495
							v437 = v437 + int32(4)
							continue
						} else {
							break
						}
						break
					}
				} else {
					v86 = v65 & int32(2147483644)
					v89 = v70
					v91 = v86
					v95 = base.Simd_g_const(&F_MultARGBRow_SSE2__k0)
					for {
						v97 = int32(0)
						v98 = base.Simd_g_v128_load(m, v89, v97)
						v99 = base.Simd_g_const(&F_MultARGBRow_SSE2__k13)
						v102 = base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_add(v98, v99), base.Simd_g_const(&F_MultARGBRow_SSE2__k17))
						if base.Simd_g_i32x4_extract_lane_l0(v102)&int32(1) == v97 {
							v116 = v95
						} else {
							v114 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(v98))>>(uint(int32(24))%32)))
							v116 = base.Simd_g_i32x4_splat(v114)
						}
						v117 = int32(1)
						if base.Simd_g_i32x4_extract_lane_l1(v102)&v117 == int32(0) {
							v131 = v116
						} else {
							v124 = int32(1)
							v128 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l1(v98))>>(uint(int32(24))%32)))
							v131 = base.Simd_g_i32x4_replace_lane_l1(v116, v128)
						}
						if base.Simd_g_i32x4_extract_lane_l2(v102)&int32(1) == int32(0) {
							v146 = v131
						} else {
							v139 = int32(2)
							v143 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l2(v98))>>(uint(int32(24))%32)))
							v146 = base.Simd_g_i32x4_replace_lane_l2(v131, v143)
						}
						if base.Simd_g_i32x4_extract_lane_l3(v102)&int32(1) == int32(0) {
							v161 = v146
						} else {
							v154 = int32(3)
							v158 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l3(v98))>>(uint(int32(24))%32)))
							v161 = base.Simd_g_i32x4_replace_lane_l3(v146, v158)
						}
						v163 = base.Simd_g_const(&F_MultARGBRow_SSE2__k11)
						v166 = base.Simd_g_const(&F_MultARGBRow_SSE2__k12)
						v172 = int32(8)
						v177 = int32(16)
						v194 = base.Simd_g_v128_bitselect(base.Simd_g_const(&F_MultARGBRow_SSE2__k0), base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v161, base.Simd_g_v128_and(v98, v163)), v166), int32(24)), base.Simd_g_v128_and(v98, v99)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v161, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v98, v172), v163)), v166), v177), base.Simd_g_const(&F_MultARGBRow_SSE2__k14))), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v161, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v98, v177), v163)), v166), v172), base.Simd_g_const(&F_MultARGBRow_SSE2__k15))), base.Simd_g_i32x4_lt_u(v98, base.Simd_g_const(&F_MultARGBRow_SSE2__k16)))
						v195 = base.Simd_g_i32x4_lt_u(v98, v99)
						v196 = int32(0)
						if base.Simd_g_i32x4_extract_lane_l0(v195)&int32(1) == v196 {
						} else {
							v202 = int32(0)
							base.Simd_g_v128_store32_lane_l0(m, v89, v202, v194)
						}
						v205 = int32(1)
						if base.Simd_g_i32x4_extract_lane_l1(v195)&v205 == int32(0) {
						} else {
							v214 = int32(1)
							base.Simd_g_v128_store32_lane_l1(m, v89+int32(4), int32(0), v194)
						}
						if base.Simd_g_i32x4_extract_lane_l2(v195)&int32(1) == int32(0) {
						} else {
							v225 = int32(2)
							base.Simd_g_v128_store32_lane_l2(m, v89+int32(8), int32(0), v194)
						}
						if base.Simd_g_i32x4_extract_lane_l3(v195)&int32(1) == int32(0) {
						} else {
							v236 = int32(3)
							base.Simd_g_v128_store32_lane_l3(m, v89+int32(12), int32(0), v194)
						}
						v241 = v91 + int32(-4)
						if v241 != 0 {
							v89 = v89 + int32(16)
							v91 = v241
							v95 = v194
							continue
						} else {
							break
						}
						break
					}
					if v86 != v65 {
						v424 = v86
						v436 = v65 - v424
						v437 = v70 + v424<<(uint(int32(2))%32)
						for {
							v445 = *(*int32)(unsafe.Add(mBase, uint32(v437)))
							if base.Ui32(int32(-16777217)) < base.Ui32(v445) {
							} else {
								if base.Ui32(v445) < base.Ui32(int32(16777216)) {
									v489 = int32(0)
								} else {
									v451 = int32(-16777216)
									v452 = int32(24)
									v454 = base.I32_div_u_s(v451, int32(base.Ui32(v445)>>(uint(v452)%32)))
									v455 = int32(255)
									v458 = int32(8388608)
									v465 = int32(8)
									v472 = int32(16)
									v489 = int32(base.Ui32(v454*(v445&v455)+v458)>>(uint(v452)%32)) | v445&v451 | int32(base.Ui32(v454*(int32(base.Ui32(v445)>>(uint(v465)%32))&v455)+v458)>>(uint(v472)%32))&int32(_a_F_MultARGBRow_SSE2_1) | int32(base.Ui32(v454*(int32(base.Ui32(v445)>>(uint(v472)%32))&v455)+v458)>>(uint(v465)%32))&int32(16711680)
								}
								*(*int32)(unsafe.Add(mBase, uint32(v437))) = v489
							}
							v495 = v436 + int32(-1)
							if v495 != 0 {
								v436 = v495
								v437 = v437 + int32(4)
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
	return
}

var F_MultARGBRow_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_MultARGBRow_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_MultARGBRow_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_MultARGBRow_SSE2__k3 = [2]uint64{0xff00000000, 0xff00000000}
var F_MultARGBRow_SSE2__k4 = [2]uint64{0x504070607060706, 0xf0e0d0c0b0a0908}
var F_MultARGBRow_SSE2__k5 = [2]uint64{0x706050403020100, 0xd0c0f0e0f0e0f0e}
var F_MultARGBRow_SSE2__k6 = [2]uint64{0x80008000800080, 0x80008000800080}
var F_MultARGBRow_SSE2__k7 = [2]uint64{0x10100000101, 0x10100000101}
var F_MultARGBRow_SSE2__k8 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_MultARGBRow_SSE2__k9 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_MultARGBRow_SSE2__k10 = [2]uint64{0x1010100010101, 0x1010100010101}
var F_MultARGBRow_SSE2__k11 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_MultARGBRow_SSE2__k12 = [2]uint64{0x80000000800000, 0x80000000800000}
var F_MultARGBRow_SSE2__k13 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}
var F_MultARGBRow_SSE2__k14 = [2]uint64{0xff000000ff00, 0xff000000ff00}
var F_MultARGBRow_SSE2__k15 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_MultARGBRow_SSE2__k16 = [2]uint64{0x100000001000000, 0x100000001000000}
var F_MultARGBRow_SSE2__k17 = [2]uint64{0xfe000000fe000000, 0xfe000000fe000000}

func F_MultRow_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 base.V128
	_ = v26
	var v27 base.V128
	_ = v27
	var v28 base.V128
	_ = v28
	var v31 base.V128
	_ = v31
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v52 int32
	_ = v52
	var v58 int32
	_ = v58
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v119 int32
	_ = v119
	v5 = int32(0)
	if l3 != 0 {
		v58 = v5
	} else {
		if l2 < int32(8) {
			v58 = v5
		} else {
			v19 = int32(0)
			for {
				v23 = l0 + v19
				v25 = int32(0)
				v26 = base.Simd_g_v128_load64_zero(m, l1+v19, v25)
				v27 = base.Simd_g_const(&F_MultRow_SSE2__k0)
				v28 = base.Simd_g_const(&F_MultRow_SSE2__k1)
				v31 = base.Simd_g_v128_load64_zero(m, v23, v25)
				v36 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_mul(base.Simd_g_i8x16_shuffle2(v26, v27, base.Simd_g_const(&F_MultRow_SSE2__k2), base.Simd_g_const(&F_MultRow_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v31, v27, base.Simd_g_const(&F_MultRow_SSE2__k2), base.Simd_g_const(&F_MultRow_SSE2__k3))), base.Simd_g_const(&F_MultRow_SSE2__k4))
				v38 = base.Simd_g_const(&F_MultRow_SSE2__k5)
				base.Simd_g_v128_store64_lane_l0(m, v23, v25, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v36), v38), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v36), v38), base.Simd_g_const(&F_MultRow_SSE2__k6), base.Simd_g_const(&F_MultRow_SSE2__k7)), v27))
				v52 = v19 + int32(8)
				if v19+int32(16) <= l2 {
					v19 = v52
					continue
				} else {
					break
				}
				break
			}
			v58 = v52
		}
	}
	v63 = l2 - v58
	if v63 < int32(1) {
	} else {
		v66 = l0 + v58
		v67 = l1 + v58
		if v63 < int32(1) {
		} else {
			if l3 == int32(0) {
				v96 = v66
				v97 = v67
				v98 = v63
				for {
					v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v97))))
					if v100 == int32(255) {
					} else {
						if v100 != 0 {
							v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v96))))
							v112 = int32(base.Ui32(v100*v105*int32(_a_F_MultRow_SSE2_0)+int32(8388608)) >> (uint(int32(24)) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v112)
						} else {
							v103 = int32(0)
							*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v103)
						}
					}
					v114 = int32(1)
					v119 = v98 + int32(-1)
					if v119 != 0 {
						v96 = v96 + v114
						v97 = v97 + v114
						v98 = v119
						continue
					} else {
						break
					}
					break
				}
			} else {
				v72 = v66
				v73 = v67
				v74 = v63
				for {
					v76 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v73))))
					if v76 == int32(255) {
					} else {
						if v76 != 0 {
							v82 = base.I32_div_u_s(int32(-16777216), v76)
							v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v72))))
							v88 = int32(base.Ui32(v82*v83+int32(8388608)) >> (uint(int32(24)) % 32))
							*(*uint8)(unsafe.Add(mBase, uint32(v72))) = uint8(v88)
						} else {
							v79 = int32(0)
							*(*uint8)(unsafe.Add(mBase, uint32(v72))) = uint8(v79)
						}
					}
					v90 = int32(1)
					v95 = v74 + int32(-1)
					if v95 != 0 {
						v72 = v72 + v90
						v73 = v73 + v90
						v74 = v95
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

var F_MultRow_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_MultRow_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_MultRow_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_MultRow_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_MultRow_SSE2__k4 = [2]uint64{0x80008000800080, 0x80008000800080}
var F_MultRow_SSE2__k5 = [2]uint64{0x10100000101, 0x10100000101}
var F_MultRow_SSE2__k6 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_MultRow_SSE2__k7 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}

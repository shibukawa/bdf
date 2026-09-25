//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_DeleteVP8Encoder(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v5 int32
	_ = v5
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v21 int32
	_ = v21
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	v3 = F_VP8EncDeleteAlpha(m, l0)
	mBase = m.M
	v5 = l0 + int32(344)
	if v5 == int32(0) {
	} else {
		v10 = *(*int32)(unsafe.Add(mBase, uint32(v5)))
		if v10 == int32(0) {
		} else {
			v14 = v10
			for {
				v16 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
				F_WebPSafeFree(m, v14)
				mBase = m.M
				if v16 != 0 {
					v14 = v16
					continue
				} else {
					break
				}
				break
			}
		}
		v21 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v5)+20)) = v21
		*(*int64)(unsafe.Add(mBase, uint32(v5)+8)) = int64(0)
		*(*int32)(unsafe.Add(mBase, uint32(v5))) = v21
		v27 = *(*int32)(unsafe.Add(mBase, uint32(v5)+16))
		v28 = int32(_a_F_DeleteVP8Encoder_0)
		if v28 < v27 {
			v31 = v27
		} else {
			v31 = v28
		}
		*(*int32)(unsafe.Add(mBase, uint32(v5)+16)) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v5)+4)) = v5
	}
	F_free(m, l0)
	mBase = m.M
	return v3
}
func F_InitVP8Encoder(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 float32
	_ = v51
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v73 int64
	_ = v73
	var v74 int32
	_ = v74
	var v81 int64
	_ = v81
	var v82 int32
	_ = v82
	var v83 int64
	_ = v83
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v236 int32
	_ = v236
	var v243 int32
	_ = v243
	var v249 int32
	_ = v249
	var v253 int32
	_ = v253
	var v257 int32
	_ = v257
	var v263 int32
	_ = v263
	var v266 int32
	_ = v266
	var v273 int32
	_ = v273
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v282 int32
	_ = v282
	var v283 int32
	_ = v283
	var v292 int32
	_ = v292
	var v295 int32
	_ = v295
	var v297 int32
	_ = v297
	var v303 int32
	_ = v303
	var v306 int32
	_ = v306
	var v309 float32
	_ = v309
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v323 int32
	_ = v323
	var v333 int32
	_ = v333
	var v337 int32
	_ = v337
	var v340 int32
	_ = v340
	var v346 int32
	_ = v346
	var v347 int32
	_ = v347
	var v352 int32
	_ = v352
	var v353 int32
	_ = v353
	var v356 int32
	_ = v356
	var v360 int32
	_ = v360
	var v378 int32
	_ = v378
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v404 int32
	_ = v404
	var v411 int32
	_ = v411
	var v428 int32
	_ = v428
	var v431 int32
	_ = v431
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v457 int32
	_ = v457
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v491 int32
	_ = v491
	var v492 float32
	_ = v492
	var v500 float32
	_ = v500
	var v506 int32
	_ = v506
	var v508 int32
	_ = v508
	var v509 int32
	_ = v509
	var v515 int32
	_ = v515
	var v518 int32
	_ = v518
	v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	if v21 != 0 {
		v22 = int32(2079)
	} else {
		v22 = int32(0)
	}
	v23 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v24 = int32(15)
	v26 = int32(4)
	v27 = (v23 + v24) >> (uint(v26) % 32)
	v28 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v30 = v28 + v24
	v32 = v30 >> (uint(v26) % 32)
	v33 = v27 * v32
	v34 = int32(2)
	v35 = v33 << (uint(v34) % 32)
	v38 = int32(1)
	v41 = v32 << (uint(v34) % 32)
	v43 = v41 | v38
	v44 = (v27<<(uint(v34)%32) | v38) * v43
	v46 = v30 & int32(-16)
	v48 = v46 << (uint(v38) % 32)
	v50 = v41 + int32(35)
	v51 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
	if base.F32_le(v51, float32(98)) != 0 {
		v59 = v41
	} else {
		v55 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
		if v55 < int32(2) {
			v59 = int32(0)
		} else {
			v59 = v41
		}
	}
	v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
	v73 = base.I64_extend_i32_u(v48) + base.I64_extend_i32_u(v50) + base.I64_extend_i32_u(v22) + base.I64_extend_i32_u(v35) + base.I64_extend_i32_u(v44) + base.I64_extend_i32_u(v59) + int64(22758)
	v74 = int32(1)
	if v73 == int64(0) {
		v93 = F_malloc(m, base.I32_wrap_i64(v73)*v74)
		mBase = m.M
		v95 = v93
	} else {
		v81 = base.I64_div_u_s(int64(2147418112), v73)
		v82 = int32(0)
		v83 = base.I64_extend_i32_u(v74)
		if base.Ui64(int64(4294967295)) < base.Ui64(v83*v73) {
			v95 = v82
		} else {
			if base.Ui64(v81) < base.Ui64(v83) {
				v95 = v82
			} else {
				v93 = F_malloc(m, base.I32_wrap_i64(v73)*v74)
				mBase = m.M
				v95 = v93
			}
		}
	}
	if v95 != 0 {
		base.MemoryFill(m, v95, int32(0), int32(_a_F_InitVP8Encoder_0))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+48)) = v43
		*(*int32)(unsafe.Add(mBase, uint32(v95)+44)) = v27
		*(*int32)(unsafe.Add(mBase, uint32(v95)+40)) = v32
		*(*int32)(unsafe.Add(mBase, uint32(v95))) = l0
		v229 = int32(-32)
		v230 = (v95 + int32(_a_F_InitVP8Encoder_1)) & v229
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[0]))) = v230
		v232 = int32(1)
		v233 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+52)) = v232 << (uint(v233) % 32)
		v236 = v230 + v35
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[1]))) = v236 + v43 + v232
		v243 = v236 + v44 + int32(31)
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[2]))) = v243&v229 | int32(4)
		v249 = v243 + v50
		if v21 != 0 {
			v253 = v249 & v229
		} else {
			v253 = int32(0)
		}
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[3]))) = v253
		v257 = (v249 + v22) & int32(-32)
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[4]))) = v257
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[5]))) = v257 + v46
		if v59 != 0 {
			v263 = v257 + v48
		} else {
			v263 = int32(0)
		}
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[6]))) = v263
		v266 = int32(0)
		if base.B2i32(v266 < v60)|base.B2i32(v266 < v21) == v266 {
			v276 = int32(2)
		} else {
			v273 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
			v276 = base.B2i32(v273 != int32(1))
		}
		*(*int32)(unsafe.Add(mBase, uint32(v95)+4)) = l1
		*(*int32)(unsafe.Add(mBase, uint32(v95)+36)) = v276
		v279 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[7]))) = v279
		v282 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
		v283 = int32(100) - v282
		if v279 <= int32(5) {
			if v279 != int32(5) {
				v292 = base.B2i32(int32(2) < v279)
			} else {
				v292 = int32(2)
			}
		} else {
			v292 = int32(3)
		}
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[8]))) = v292
		v295 = base.I32_div_s(int32(1069547520), v33)
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[9]))) = v295
		v297 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[10]))) = v297
		v303 = base.I32_div_u_s(v283*v283<<(uint(int32(16))%32), int32(_a_F_InitVP8Encoder_2))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[11]))) = v303
		v306 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		if int32(0) < v306 {
			v312 = int32(1)
		} else {
			v309 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
			v312 = base.F32_gt(v309, float32(0))
		}
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[12]))) = v312
		v314 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
		if v314 != 0 {
		} else {
			v315 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[13]))) = base.B2i32(v292 != v315)
			if v292 == v315 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v95)+52)) = int32(1)
			}
		}
		F_VP8EncDspInit(m)
		mBase = m.M
		v323 = int32(-1)
		*(*uint16)(unsafe.Add(mBase, uint32(v95)+3416)) = uint16(v323)
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[14]))) = int32(0)
		*(*uint8)(unsafe.Add(mBase, uint32(v95+int32(3418)))) = uint8(v323)
		v333 = m.G1
		v337 = F_memcpy(m, v95+int32(3420), v333+int32(_a_F_InitVP8Encoder_3), int32(1056))
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[15]))) = int32(1)
		v340 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v95)+32)) = v340
		*(*int64)(unsafe.Add(mBase, uint32(v95)+16)) = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(v95)+8)) = int64(1)
		v346 = *(*int32)(unsafe.Add(mBase, uint32(v95)))
		v347 = *(*int32)(unsafe.Add(mBase, uint32(v346)+24))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+24)) = v347
		*(*int32)(unsafe.Add(mBase, uint32(v95)+28)) = base.B2i32(int32(1) < v347)
		v352 = *(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[1])))
		v353 = *(*int32)(unsafe.Add(mBase, uint32(v95)+40))
		if v353 < v340 {
		} else {
			v356 = *(*int32)(unsafe.Add(mBase, uint32(v95)+48))
			v360 = int32(-1)
			for {
				v378 = int32(0)
				*(*uint8)(unsafe.Add(mBase, uint32(v352-v356+v360))) = uint8(v378)
				v381 = v360 + int32(1)
				v382 = *(*int32)(unsafe.Add(mBase, uint32(v95)+40))
				if v381 < v382<<(uint(int32(2))%32) {
					v360 = v381
					continue
				} else {
					break
				}
				break
			}
		}
		v404 = *(*int32)(unsafe.Add(mBase, uint32(v95)+44))
		if v404 < int32(1) {
		} else {
			v411 = int32(0)
			for {
				v428 = *(*int32)(unsafe.Add(mBase, uint32(v95)+48))
				v431 = int32(0)
				*(*uint8)(unsafe.Add(mBase, uint32(v352+int32(-1)+v428*v411))) = uint8(v431)
				v434 = v411 + int32(1)
				v435 = *(*int32)(unsafe.Add(mBase, uint32(v95)+44))
				if v434 < v435<<(uint(int32(2))%32) {
					v411 = v434
					continue
				} else {
					break
				}
				break
			}
		}
		v457 = *(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[2])))
		*(*int32)(unsafe.Add(mBase, uint32(v457+int32(-4)))) = int32(0)
		v465 = m.G1
		v468 = *(*int32)(unsafe.Add(mBase, uint32(v465)+uint32(_c_F_InitVP8Encoder[16])))
		v469 = m.G3
		v470 = *(*int32)(unsafe.Add(mBase, uint32(v469)))
		if v468 == v470 {
		} else {
			v472 = m.G2
			v473 = m.G1
			*(*int32)(unsafe.Add(mBase, uint32(v473)+uint32(_c_F_InitVP8Encoder[17]))) = v472 + int32(128)
			*(*int32)(unsafe.Add(mBase, uint32(v473)+uint32(_c_F_InitVP8Encoder[18]))) = v472 + int32(129)
			*(*int32)(unsafe.Add(mBase, uint32(v473)+uint32(_c_F_InitVP8Encoder[16]))) = v470
		}
		F_VP8EncInitAlpha(m, v95)
		mBase = m.M
		v491 = v95 + int32(344)
		v492 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
		v500 = base.F32_mul(base.F32_add(base.F32_div(base.F32_mul(v492, float32(5)), float32(100)), float32(1)), base.F32_convert_i32_s(v35))
		if base.F32_lt(base.F32_abs(v500), float32(2.1474836e+09)) == int32(0) {
			v508 = int32(-2147483648)
		} else {
			v506 = base.I32_trunc_f32_s(v500)
			v508 = v506
		}
		v509 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v491)+20)) = v509
		*(*int32)(unsafe.Add(mBase, uint32(v491))) = v509
		*(*int64)(unsafe.Add(mBase, uint32(v491)+8)) = int64(0)
		v515 = int32(_a_F_InitVP8Encoder_4)
		if v515 < v508 {
			v518 = v508
		} else {
			v518 = v515
		}
		*(*int32)(unsafe.Add(mBase, uint32(v491)+16)) = v518
		*(*int32)(unsafe.Add(mBase, uint32(v491)+4)) = v491
		return v95
	} else {
		v97 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
		if v97 != 0 {
			return v95
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
			return v95
		}
	}
}
func F_VP8AdjustFilterStrength(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v19 float64
	_ = v19
	var v24 int32
	_ = v24
	var v27 float64
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v36 float64
	_ = v36
	var v39 float64
	_ = v39
	var v42 float64
	_ = v42
	var v43 int32
	_ = v43
	var v44 float64
	_ = v44
	var v45 int32
	_ = v45
	var v46 float64
	_ = v46
	var v47 int32
	_ = v47
	var v48 float64
	_ = v48
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v65 float64
	_ = v65
	var v70 int32
	_ = v70
	var v73 float64
	_ = v73
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v82 float64
	_ = v82
	var v85 float64
	_ = v85
	var v88 float64
	_ = v88
	var v89 int32
	_ = v89
	var v90 float64
	_ = v90
	var v91 int32
	_ = v91
	var v92 float64
	_ = v92
	var v93 int32
	_ = v93
	var v94 float64
	_ = v94
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v105 int32
	_ = v105
	var v111 float64
	_ = v111
	var v116 int32
	_ = v116
	var v119 float64
	_ = v119
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v128 float64
	_ = v128
	var v131 float64
	_ = v131
	var v134 float64
	_ = v134
	var v135 int32
	_ = v135
	var v136 float64
	_ = v136
	var v137 int32
	_ = v137
	var v138 float64
	_ = v138
	var v139 int32
	_ = v139
	var v140 float64
	_ = v140
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v151 int32
	_ = v151
	var v157 float64
	_ = v157
	var v162 int32
	_ = v162
	var v165 float64
	_ = v165
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v174 float64
	_ = v174
	var v177 float64
	_ = v177
	var v180 float64
	_ = v180
	var v181 int32
	_ = v181
	var v182 float64
	_ = v182
	var v183 int32
	_ = v183
	var v184 float64
	_ = v184
	var v185 int32
	_ = v185
	var v186 float64
	_ = v186
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v197 int32
	_ = v197
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v211 int32
	_ = v211
	var v213 int32
	_ = v213
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v227 int32
	_ = v227
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v240 int32
	_ = v240
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v248 int32
	_ = v248
	var v250 int32
	_ = v250
	var v251 int32
	_ = v251
	var v255 int32
	_ = v255
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v263 int32
	_ = v263
	var v265 int32
	_ = v265
	var v268 int32
	_ = v268
	var v269 int32
	_ = v269
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v280 int32
	_ = v280
	var v281 int32
	_ = v281
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v289 int32
	_ = v289
	var v291 int32
	_ = v291
	var v293 int32
	_ = v293
	var v297 int32
	_ = v297
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
	if v14 == int32(0) {
		v201 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
		v202 = *(*int32)(unsafe.Add(mBase, uint32(v201)+32))
		if v202 < int32(1) {
		} else {
			v205 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v13)+634)))
			v206 = *(*int32)(unsafe.Add(mBase, uint32(v13)+1096))
			v207 = m.G1
			v208 = *(*int32)(unsafe.Add(mBase, uint32(v13)+1092))
			v211 = *(*int32)(unsafe.Add(mBase, uint32(v13)+16))
			v213 = v211 << (uint(int32(6)) % 32)
			v217 = v206 * v205 >> (uint(int32(3)) % 32)
			v218 = int32(63)
			if v217 < v218 {
				v221 = v217
			} else {
				v221 = v218
			}
			v223 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v207+int32(_a_F_VP8AdjustFilterStrength_0)+v213+v221))))
			if v223 <= v208 {
				v226 = v208
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v13)+1092)) = v223
				v226 = v223
			}
			v227 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v13)+1378)))
			v228 = *(*int32)(unsafe.Add(mBase, uint32(v13)+1840))
			v229 = m.G1
			v230 = *(*int32)(unsafe.Add(mBase, uint32(v13)+1836))
			v236 = v228 * v227 >> (uint(int32(3)) % 32)
			v237 = int32(63)
			if v236 < v237 {
				v240 = v236
			} else {
				v240 = v237
			}
			v242 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v229+int32(_a_F_VP8AdjustFilterStrength_0)+v213+v240))))
			if v242 <= v230 {
				v245 = v230
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v13)+1836)) = v242
				v245 = v242
			}
			v246 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v13)+2122)))
			v247 = *(*int32)(unsafe.Add(mBase, uint32(v13)+2584))
			v248 = m.G1
			if base.Ui32(v245) < base.Ui32(v226) {
				v250 = v226
			} else {
				v250 = v245
			}
			v251 = *(*int32)(unsafe.Add(mBase, uint32(v13)+2580))
			v255 = v211 << (uint(int32(6)) % 32)
			v259 = v247 * v246 >> (uint(int32(3)) % 32)
			v260 = int32(63)
			if v259 < v260 {
				v263 = v259
			} else {
				v263 = v260
			}
			v265 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v248+int32(_a_F_VP8AdjustFilterStrength_0)+v255+v263))))
			if v265 <= v251 {
				v268 = v251
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v13)+2580)) = v265
				v268 = v265
			}
			v269 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v13)+2866)))
			v270 = *(*int32)(unsafe.Add(mBase, uint32(v13)+3328))
			v271 = m.G1
			if base.Ui32(v268) < base.Ui32(v250) {
				v273 = v250
			} else {
				v273 = v268
			}
			v274 = *(*int32)(unsafe.Add(mBase, uint32(v13)+3324))
			v280 = v270 * v269 >> (uint(int32(3)) % 32)
			v281 = int32(63)
			if v280 < v281 {
				v284 = v280
			} else {
				v284 = v281
			}
			v286 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v271+int32(_a_F_VP8AdjustFilterStrength_0)+v255+v284))))
			if v286 <= v274 {
				v289 = v274
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v13)+3324)) = v286
				v289 = v286
			}
			if base.Ui32(v289) < base.Ui32(v273) {
				v291 = v273
			} else {
				v291 = v289
			}
			v293 = int32(12)
			v297 = v291
			*(*int32)(unsafe.Add(mBase, uint32(v13+v293))) = v297
		}
	} else {
		v19 = *(*float64)(unsafe.Add(mBase, uint32(v14)))
		v24 = v14 + int32(24)
		v27 = base.F64_mul(v19, float64(1.00001))
		v28 = int32(0)
		v29 = int32(1)
		for {
			v36 = *(*float64)(unsafe.Add(mBase, uint32(v24)))
			v39 = *(*float64)(unsafe.Add(mBase, uint32(v24+int32(-8))))
			v42 = *(*float64)(unsafe.Add(mBase, uint32(v24+int32(-16))))
			v43 = base.F64_gt(v42, v27)
			if v43 != 0 {
				v44 = v42
			} else {
				v44 = v27
			}
			v45 = base.F64_gt(v39, v44)
			if v45 != 0 {
				v46 = v39
			} else {
				v46 = v44
			}
			v47 = base.F64_gt(v36, v46)
			if v47 != 0 {
				v48 = v36
			} else {
				v48 = v46
			}
			if v43 != 0 {
				v53 = v29
			} else {
				v53 = v28
			}
			if v45 != 0 {
				v54 = v29 + int32(1)
			} else {
				v54 = v53
			}
			if v47 != 0 {
				v55 = v29 + int32(2)
			} else {
				v55 = v54
			}
			v59 = v29 + int32(3)
			if v59 != int32(64) {
				v24 = v24 + int32(24)
				v27 = v48
				v28 = v55
				v29 = v59
				continue
			} else {
				break
			}
			break
		}
		*(*int32)(unsafe.Add(mBase, uint32(v13)+1092)) = v55
		v65 = *(*float64)(unsafe.Add(mBase, uint32(v14)+512))
		v70 = v14 + int32(536)
		v73 = base.F64_mul(v65, float64(1.00001))
		v74 = int32(0)
		v75 = int32(1)
		for {
			v82 = *(*float64)(unsafe.Add(mBase, uint32(v70)))
			v85 = *(*float64)(unsafe.Add(mBase, uint32(v70+int32(-8))))
			v88 = *(*float64)(unsafe.Add(mBase, uint32(v70+int32(-16))))
			v89 = base.F64_gt(v88, v73)
			if v89 != 0 {
				v90 = v88
			} else {
				v90 = v73
			}
			v91 = base.F64_gt(v85, v90)
			if v91 != 0 {
				v92 = v85
			} else {
				v92 = v90
			}
			v93 = base.F64_gt(v82, v92)
			if v93 != 0 {
				v94 = v82
			} else {
				v94 = v92
			}
			if v89 != 0 {
				v99 = v75
			} else {
				v99 = v74
			}
			if v91 != 0 {
				v100 = v75 + int32(1)
			} else {
				v100 = v99
			}
			if v93 != 0 {
				v101 = v75 + int32(2)
			} else {
				v101 = v100
			}
			v105 = v75 + int32(3)
			if v105 != int32(64) {
				v70 = v70 + int32(24)
				v73 = v94
				v74 = v101
				v75 = v105
				continue
			} else {
				break
			}
			break
		}
		*(*int32)(unsafe.Add(mBase, uint32(v13)+1836)) = v101
		v111 = *(*float64)(unsafe.Add(mBase, uint32(v14)+1024))
		v116 = v14 + int32(1048)
		v119 = base.F64_mul(v111, float64(1.00001))
		v120 = int32(0)
		v121 = int32(1)
		for {
			v128 = *(*float64)(unsafe.Add(mBase, uint32(v116)))
			v131 = *(*float64)(unsafe.Add(mBase, uint32(v116+int32(-8))))
			v134 = *(*float64)(unsafe.Add(mBase, uint32(v116+int32(-16))))
			v135 = base.F64_gt(v134, v119)
			if v135 != 0 {
				v136 = v134
			} else {
				v136 = v119
			}
			v137 = base.F64_gt(v131, v136)
			if v137 != 0 {
				v138 = v131
			} else {
				v138 = v136
			}
			v139 = base.F64_gt(v128, v138)
			if v139 != 0 {
				v140 = v128
			} else {
				v140 = v138
			}
			if v135 != 0 {
				v145 = v121
			} else {
				v145 = v120
			}
			if v137 != 0 {
				v146 = v121 + int32(1)
			} else {
				v146 = v145
			}
			if v139 != 0 {
				v147 = v121 + int32(2)
			} else {
				v147 = v146
			}
			v151 = v121 + int32(3)
			if v151 != int32(64) {
				v116 = v116 + int32(24)
				v119 = v140
				v120 = v147
				v121 = v151
				continue
			} else {
				break
			}
			break
		}
		*(*int32)(unsafe.Add(mBase, uint32(v13)+2580)) = v147
		v157 = *(*float64)(unsafe.Add(mBase, uint32(v14)+1536))
		v162 = v14 + int32(1560)
		v165 = base.F64_mul(v157, float64(1.00001))
		v166 = int32(0)
		v167 = int32(1)
		for {
			v174 = *(*float64)(unsafe.Add(mBase, uint32(v162)))
			v177 = *(*float64)(unsafe.Add(mBase, uint32(v162+int32(-8))))
			v180 = *(*float64)(unsafe.Add(mBase, uint32(v162+int32(-16))))
			v181 = base.F64_gt(v180, v165)
			if v181 != 0 {
				v182 = v180
			} else {
				v182 = v165
			}
			v183 = base.F64_gt(v177, v182)
			if v183 != 0 {
				v184 = v177
			} else {
				v184 = v182
			}
			v185 = base.F64_gt(v174, v184)
			if v185 != 0 {
				v186 = v174
			} else {
				v186 = v184
			}
			if v181 != 0 {
				v191 = v167
			} else {
				v191 = v166
			}
			if v183 != 0 {
				v192 = v167 + int32(1)
			} else {
				v192 = v191
			}
			if v185 != 0 {
				v193 = v167 + int32(2)
			} else {
				v193 = v192
			}
			v197 = v167 + int32(3)
			if v197 != int32(64) {
				v162 = v162 + int32(24)
				v165 = v186
				v166 = v193
				v167 = v197
				continue
			} else {
				break
			}
			break
		}
		v293 = int32(3324)
		v297 = v193
		*(*int32)(unsafe.Add(mBase, uint32(v13+v293))) = v297
	}
	return
}
func F_VP8ApplyNearLossless(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v18 int64
	_ = v18
	var v19 int32
	_ = v19
	var v26 int64
	_ = v26
	var v27 int32
	_ = v27
	var v28 int64
	_ = v28
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v139 int32
	_ = v139
	v12 = base.I32_div_s(l1, int32(-20))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v18 = base.I64_extend_i32_s(v15 * int32(3))
	v19 = int32(4)
	if v18 == int64(0) {
		v38 = F_malloc(m, base.I32_wrap_i64(v18)*v19)
		mBase = m.M
		v40 = v38
	} else {
		v26 = base.I64_div_u_s(int64(2147418112), v18)
		v27 = int32(0)
		v28 = base.I64_extend_i32_u(v19)
		if base.Ui64(int64(4294967295)) < base.Ui64(v28*v18) {
			v40 = v27
		} else {
			if base.Ui64(v26) < base.Ui64(v28) {
				v40 = v27
			} else {
				v38 = F_malloc(m, base.I32_wrap_i64(v18)*v19)
				mBase = m.M
				v40 = v38
			}
		}
	}
	if v40 != 0 {
		v44 = int32(64)
		if base.B2i32(v15 < v44)&base.B2i32(v14 < v44) != 0 {
			if v14 < int32(1) {
			} else {
				v53 = int32(1)
				v56 = v15 << (uint(int32(2)) % 32)
				if v14 == v53 {
					v99 = int32(0)
				} else {
					v68 = int32(0)
					v70 = l2
					for {
						v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
						v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v78 = int32(2)
						v81 = F_memcpy(m, v70, v75+v76*v68<<(uint(v78)%32), v56)
						mBase = m.M
						v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
						v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v91 = F_memcpy(m, v81+v56, v83+v84*(v68+int32(1))<<(uint(v78)%32), v56)
						mBase = m.M
						v94 = v68 + v78
						if v14&int32(2147483646) != v94 {
							v68 = v94
							v70 = v81 + v15<<(uint(int32(3))%32)
							continue
						} else {
							break
						}
						break
					}
					v99 = v94
				}
				if v14&v53 == int32(0) {
				} else {
					v109 = int32(2)
					v112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
					v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
					v118 = F_memcpy(m, l2+v99*v15<<(uint(v109)%32), v112+v113*v99<<(uint(v109)%32), v56)
					mBase = m.M
				}
			}
		} else {
			if int32(2) < v14 {
				v119 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
				F_NearLossless(m, v15, v14, v119, v13, v12+int32(5), v40, l2)
				mBase = m.M
				v124 = v12 + int32(4)
				if v124 == int32(0) {
				} else {
					v127 = v124
					for {
						F_NearLossless(m, v15, v14, l2, v15, v127, v40, l2)
						mBase = m.M
						v139 = v127 + int32(-1)
						if v139 != 0 {
							v127 = v139
							continue
						} else {
							break
						}
						break
					}
				}
			} else {
				if v14 < int32(1) {
				} else {
					v53 = int32(1)
					v56 = v15 << (uint(int32(2)) % 32)
					if v14 == v53 {
						v99 = int32(0)
					} else {
						v68 = int32(0)
						v70 = l2
						for {
							v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
							v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
							v78 = int32(2)
							v81 = F_memcpy(m, v70, v75+v76*v68<<(uint(v78)%32), v56)
							mBase = m.M
							v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
							v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
							v91 = F_memcpy(m, v81+v56, v83+v84*(v68+int32(1))<<(uint(v78)%32), v56)
							mBase = m.M
							v94 = v68 + v78
							if v14&int32(2147483646) != v94 {
								v68 = v94
								v70 = v81 + v15<<(uint(int32(3))%32)
								continue
							} else {
								break
							}
							break
						}
						v99 = v94
					}
					if v14&v53 == int32(0) {
					} else {
						v109 = int32(2)
						v112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
						v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v118 = F_memcpy(m, l2+v99*v15<<(uint(v109)%32), v112+v113*v99<<(uint(v109)%32), v56)
						mBase = m.M
					}
				}
			}
		}
		F_free(m, v40)
		mBase = m.M
		return int32(1)
	} else {
		return int32(0)
	}
}
func F_VP8BitWriterAppend(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v15 int64
	_ = v15
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int64
	_ = v26
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v41 int64
	_ = v41
	var v42 int32
	_ = v42
	var v43 int64
	_ = v43
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v80 int32
	_ = v80
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v9 != int32(-8) {
		v80 = int32(0)
		return v80
	} else {
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
		v15 = base.I64_extend_i32_u(v12) + base.I64_extend_i32_u(l2)
		if base.Ui64(v15) < base.Ui64(int64(4294967296)) {
			v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
			v23 = base.I32_wrap_i64(v15)
			if base.Ui32(v22) < base.Ui32(v23) {
				v26 = int64(1)
				v28 = v22 << (uint(int32(1)) % 32)
				if base.Ui32(v23) < base.Ui32(v28) {
					v30 = v28
				} else {
					v30 = v23
				}
				v31 = int32(1024)
				if base.Ui32(v31) < base.Ui32(v30) {
					v34 = v30
				} else {
					v34 = v31
				}
				v41 = base.I64_div_u_s(int64(2147418112), v26)
				v42 = int32(0)
				v43 = base.I64_extend_i32_u(v34)
				if base.Ui64(int64(4294967295)) < base.Ui64(v43*v26) {
					v55 = v42
				} else {
					if base.Ui64(v41) < base.Ui64(v43) {
						v55 = v42
					} else {
						v53 = F_malloc(m, base.I32_wrap_i64(v26)*v34)
						mBase = m.M
						v55 = v53
					}
				}
				if v55 != 0 {
					v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					if v61 == int32(0) {
					} else {
						v64 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v65 = F_memcpy(m, v55, v64, v61)
						mBase = m.M
					}
					v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
					F_free(m, v66)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v34
					*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v55
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					v71 = v70
					v73 = v55
					v75 = F_memcpy(m, v73+v71, l1, l2)
					mBase = m.M
					v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v76 + l2
					v80 = int32(1)
					return v80
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
					return int32(0)
				}
			} else {
				v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				v71 = v12
				v73 = v25
				v75 = F_memcpy(m, v73+v71, l1, l2)
				mBase = m.M
				v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v76 + l2
				v80 = int32(1)
				return v80
			}
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
			return int32(0)
		}
	}
}
func F_VP8BitWriterFinish(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v63 int64
	_ = v63
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v135 int32
	_ = v135
	var v141 int32
	_ = v141
	var v147 int32
	_ = v147
	var v159 int32
	_ = v159
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v178 int32
	_ = v178
	var v182 int32
	_ = v182
	var v186 int32
	_ = v186
	var v190 int32
	_ = v190
	var v192 int64
	_ = v192
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v225 int32
	_ = v225
	var v230 int32
	_ = v230
	var v233 int32
	_ = v233
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v239 int32
	_ = v239
	var v244 int32
	_ = v244
	var v246 int32
	_ = v246
	var v252 int32
	_ = v252
	var v255 int32
	_ = v255
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v264 int32
	_ = v264
	var v270 int32
	_ = v270
	var v276 int32
	_ = v276
	var v287 int32
	_ = v287
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v11 = int32(1) << (uint(int32(8)-v7) % 32)
	for {
		v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v16 = v14 >> (uint(int32(1)) % 32)
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v16
		if int32(126) < v16 {
		} else {
			v20 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v21 = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v20 << (uint(v21) % 32)
			v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v24 + v21
			v28 = m.G1
			v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(_a_F_VP8BitWriterFinish_0)+v16))))
			*(*int32)(unsafe.Add(mBase, uint32(l0))) = v32
			if v24 < int32(0) {
			} else {
				v42 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v42 + int32(-8)
				v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v48 = v42 + int32(8)
				v49 = v46 >> (uint(v48) % 32)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v46 - v49<<(uint(v48)%32)
				v53 = int32(255)
				if v49&v53 == v53 {
					v147 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v147 + int32(1)
				} else {
					v57 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					v63 = base.I64_extend_i32_u(v57+int32(1)) + base.I64_extend_i32_u(v61)
					if base.Ui64(v63) < base.Ui64(int64(4294967296)) {
						v68 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
						v69 = base.I32_wrap_i64(v63)
						if base.Ui32(v69) <= base.Ui32(v68) {
							v96 = v49 & int32(256)
							if v96 == int32(0) {
							} else {
								if v61 == int32(0) {
								} else {
									v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v104 = v101 + v61 + int32(-1)
									v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104))))
									v107 = v105 + int32(1)
									*(*uint8)(unsafe.Add(mBase, uint32(v104))) = uint8(v107)
								}
							}
							v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v110 < int32(1) {
								v135 = v61
							} else {
								if v96 != 0 {
									v115 = int32(0)
								} else {
									v115 = int32(-1)
								}
								v117 = v61
								for {
									v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v123+v117))) = uint8(v115)
									v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v126 + int32(-1)
									v130 = int32(1)
									v131 = v117 + v130
									if v130 < v126 {
										v117 = v131
										continue
									} else {
										break
									}
									break
								}
								v135 = v131
							}
							v141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v141+v135))) = uint8(v49)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v135 + int32(1)
						} else {
							v73 = v68 << (uint(int32(1)) % 32)
							if base.Ui32(v69) < base.Ui32(v73) {
								v75 = v73
							} else {
								v75 = v69
							}
							v76 = int32(1024)
							if base.Ui32(v76) < base.Ui32(v75) {
								v79 = v75
							} else {
								v79 = v76
							}
							v80 = F_WebPSafeMalloc(m, int64(1), v79)
							mBase = m.M
							if v80 != 0 {
								v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
								if v83 == int32(0) {
								} else {
									v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v87 = F_memcpy(m, v80, v86, v83)
									mBase = m.M
								}
								v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								F_WebPSafeFree(m, v88)
								mBase = m.M
								*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v79
								*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v80
								v96 = v49 & int32(256)
								if v96 == int32(0) {
								} else {
									if v61 == int32(0) {
									} else {
										v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										v104 = v101 + v61 + int32(-1)
										v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104))))
										v107 = v105 + int32(1)
										*(*uint8)(unsafe.Add(mBase, uint32(v104))) = uint8(v107)
									}
								}
								v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								if v110 < int32(1) {
									v135 = v61
								} else {
									if v96 != 0 {
										v115 = int32(0)
									} else {
										v115 = int32(-1)
									}
									v117 = v61
									for {
										v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										*(*uint8)(unsafe.Add(mBase, uint32(v123+v117))) = uint8(v115)
										v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
										*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v126 + int32(-1)
										v130 = int32(1)
										v131 = v117 + v130
										if v130 < v126 {
											v117 = v131
											continue
										} else {
											break
										}
										break
									}
									v135 = v131
								}
								v141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								*(*uint8)(unsafe.Add(mBase, uint32(v141+v135))) = uint8(v49)
								*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v135 + int32(1)
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
							}
						}
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
					}
				}
			}
		}
		v159 = int32(1)
		if base.Ui32(v159) < base.Ui32(v11) {
			v11 = int32(base.Ui32(v11) >> (uint(v159) % 32))
			continue
		} else {
			break
		}
		break
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = int32(-8)
	v175 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v176 = int32(8)
	v178 = v175 >> (uint(v176) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v175 - v178<<(uint(v176)%32)
	v182 = int32(255)
	if v178&v182 == v182 {
		v276 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v276 + int32(1)
	} else {
		v186 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		v190 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
		v192 = base.I64_extend_i32_u(v186+int32(1)) + base.I64_extend_i32_u(v190)
		if base.Ui64(v192) < base.Ui64(int64(4294967296)) {
			v197 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
			v198 = base.I32_wrap_i64(v192)
			if base.Ui32(v198) <= base.Ui32(v197) {
				v225 = v178 & int32(256)
				if v225 == int32(0) {
				} else {
					if v190 == int32(0) {
					} else {
						v230 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v233 = v230 + v190 + int32(-1)
						v234 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v233))))
						v236 = v234 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v233))) = uint8(v236)
					}
				}
				v239 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				if v239 < int32(1) {
					v264 = v190
				} else {
					if v225 != 0 {
						v244 = int32(0)
					} else {
						v244 = int32(-1)
					}
					v246 = v190
					for {
						v252 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						*(*uint8)(unsafe.Add(mBase, uint32(v252+v246))) = uint8(v244)
						v255 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v255 + int32(-1)
						v259 = int32(1)
						v260 = v246 + v259
						if v259 < v255 {
							v246 = v260
							continue
						} else {
							break
						}
						break
					}
					v264 = v260
				}
				v270 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				*(*uint8)(unsafe.Add(mBase, uint32(v270+v264))) = uint8(v178)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v264 + int32(1)
			} else {
				v202 = v197 << (uint(int32(1)) % 32)
				if base.Ui32(v198) < base.Ui32(v202) {
					v204 = v202
				} else {
					v204 = v198
				}
				v205 = int32(1024)
				if base.Ui32(v205) < base.Ui32(v204) {
					v208 = v204
				} else {
					v208 = v205
				}
				v209 = F_WebPSafeMalloc(m, int64(1), v208)
				mBase = m.M
				if v209 != 0 {
					v212 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					if v212 == int32(0) {
					} else {
						v215 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v216 = F_memcpy(m, v209, v215, v212)
						mBase = m.M
					}
					v217 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
					F_WebPSafeFree(m, v217)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v208
					*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v209
					v225 = v178 & int32(256)
					if v225 == int32(0) {
					} else {
						if v190 == int32(0) {
						} else {
							v230 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							v233 = v230 + v190 + int32(-1)
							v234 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v233))))
							v236 = v234 + int32(1)
							*(*uint8)(unsafe.Add(mBase, uint32(v233))) = uint8(v236)
						}
					}
					v239 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					if v239 < int32(1) {
						v264 = v190
					} else {
						if v225 != 0 {
							v244 = int32(0)
						} else {
							v244 = int32(-1)
						}
						v246 = v190
						for {
							v252 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v252+v246))) = uint8(v244)
							v255 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v255 + int32(-1)
							v259 = int32(1)
							v260 = v246 + v259
							if v259 < v255 {
								v246 = v260
								continue
							} else {
								break
							}
							break
						}
						v264 = v260
					}
					v270 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
					*(*uint8)(unsafe.Add(mBase, uint32(v270+v264))) = uint8(v178)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v264 + int32(1)
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
				}
			}
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
		}
	}
	v287 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	return v287
}
func F_VP8BitWriterInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int64
	_ = v5
	var v17 int64
	_ = v17
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v28 int64
	_ = v28
	var v29 int32
	_ = v29
	var v30 int64
	_ = v30
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v48 int32
	_ = v48
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	v5 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+8)) = int64(-34359738368)
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(254)
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v5
	if l1 == int32(0) {
		return int32(1)
	} else {
		v17 = int64(1)
		v18 = int32(1024)
		if base.Ui32(v18) < base.Ui32(l1) {
			v21 = l1
		} else {
			v21 = v18
		}
		v28 = base.I64_div_u_s(int64(2147418112), v17)
		v29 = int32(0)
		v30 = base.I64_extend_i32_u(v21)
		if base.Ui64(int64(4294967295)) < base.Ui64(v30*v17) {
			v42 = v29
		} else {
			if base.Ui64(v28) < base.Ui64(v30) {
				v42 = v29
			} else {
				v40 = F_malloc(m, base.I32_wrap_i64(v17)*v21)
				mBase = m.M
				v42 = v40
			}
		}
		if v42 != 0 {
			v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
			if v48 == int32(0) {
			} else {
				v53 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(16))))
				v54 = F_memcpy(m, v42, v53, v48)
				mBase = m.M
			}
			v55 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
			F_free(m, v55)
			mBase = m.M
			*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v21
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v42
			return int32(1)
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
			return int32(0)
		}
	}
}
func F_VP8BitWriterWipeOut(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v11 int64
	_ = v11
	if l0 == int32(0) {
	} else {
		v6 = l0 + int32(16)
		v7 = *(*int32)(unsafe.Add(mBase, uint32(v6)))
		F_free(m, v7)
		mBase = m.M
		v11 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v11
		*(*int64)(unsafe.Add(mBase, uint32(v6))) = v11
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(8)))) = v11
		*(*int64)(unsafe.Add(mBase, uint32(l0))) = v11
	}
	return
}
func F_VP8CalculateLevelCosts(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v42 int32
	_ = v42
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v103 int32
	_ = v103
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v131 int32
	_ = v131
	var v150 int32
	_ = v150
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v208 int32
	_ = v208
	var v215 int32
	_ = v215
	var v248 int32
	_ = v248
	var v259 int32
	_ = v259
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v276 int32
	_ = v276
	var v303 int32
	_ = v303
	var v305 int32
	_ = v305
	var v306 int32
	_ = v306
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v317 int32
	_ = v317
	var v326 int32
	_ = v326
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v374 int32
	_ = v374
	var v379 int32
	_ = v379
	var v382 int32
	_ = v382
	var v385 int32
	_ = v385
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v401 int32
	_ = v401
	var v402 int32
	_ = v402
	var v409 int32
	_ = v409
	var v445 int32
	_ = v445
	var v456 int32
	_ = v456
	var v461 int32
	_ = v461
	var v462 int32
	_ = v462
	var v464 int32
	_ = v464
	var v473 int32
	_ = v473
	var v500 int32
	_ = v500
	var v502 int32
	_ = v502
	var v503 int32
	_ = v503
	var v511 int32
	_ = v511
	var v512 int32
	_ = v512
	var v514 int32
	_ = v514
	var v523 int32
	_ = v523
	var v525 int32
	_ = v525
	var v527 int32
	_ = v527
	var v571 int32
	_ = v571
	var v576 int32
	_ = v576
	var v579 int32
	_ = v579
	var v582 int32
	_ = v582
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v591 int32
	_ = v591
	var v592 int32
	_ = v592
	var v593 int32
	_ = v593
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v606 int32
	_ = v606
	var v642 int32
	_ = v642
	var v653 int32
	_ = v653
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v670 int32
	_ = v670
	var v697 int32
	_ = v697
	var v699 int32
	_ = v699
	var v700 int32
	_ = v700
	var v708 int32
	_ = v708
	var v709 int32
	_ = v709
	var v711 int32
	_ = v711
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v724 int32
	_ = v724
	var v768 int32
	_ = v768
	var v773 int32
	_ = v773
	var v776 int32
	_ = v776
	var v779 int32
	_ = v779
	var v782 int32
	_ = v782
	var v789 int32
	_ = v789
	var v794 int32
	_ = v794
	var v795 int32
	_ = v795
	var v797 int32
	_ = v797
	var v799 int32
	_ = v799
	var v801 int32
	_ = v801
	var v803 int32
	_ = v803
	var v805 int32
	_ = v805
	var v807 int32
	_ = v807
	var v809 int32
	_ = v809
	var v811 int32
	_ = v811
	var v813 int32
	_ = v813
	var v815 int32
	_ = v815
	var v817 int32
	_ = v817
	var v819 int32
	_ = v819
	var v821 int32
	_ = v821
	var v823 int32
	_ = v823
	var v825 int32
	_ = v825
	var v827 int32
	_ = v827
	var v830 int32
	_ = v830
	var v923 int32
	_ = v923
	var v930 int32
	_ = v930
	v42 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CalculateLevelCosts[0])))
	if v42 == int32(0) {
	} else {
		v57 = m.G23
		v58 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+15)))
		v59 = int32(408)
		v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+14)))
		v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+13)))
		v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+12)))
		v70 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+11)))
		v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+10)))
		v76 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+9)))
		v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+8)))
		v82 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+7)))
		v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+6)))
		v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+5)))
		v91 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+4)))
		v94 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+3)))
		v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+2)))
		v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57)+1)))
		v103 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57))))
		v108 = l0 + int32(28)
		v109 = l0 + int32(17)
		v110 = l0 + int32(6)
		v131 = int32(0)
		for {
			v150 = l0 + int32(_a_F_VP8CalculateLevelCosts_0) + v131*int32(3264)
			v182 = v108
			v183 = v109
			v184 = v110
			v185 = int32(0)
			for {
				v198 = l0 + int32(4) + v131*int32(264) + v185*int32(33)
				v199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198)+1)))
				v200 = int32(1)
				v203 = v150 + v185*int32(408)
				v204 = m.G24
				v208 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v204+v199<<(uint(v200)%32)))))
				*(*uint16)(unsafe.Add(mBase, uint32(v203))) = uint16(v208)
				v215 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v204+(v199^int32(255))<<(uint(v200)%32)))))
				v248 = v200
				for {
					v259 = m.G1
					v264 = v248<<(uint(int32(2))%32) + (v259 + int32(_a_F_VP8CalculateLevelCosts_1)) + int32(-4)
					v265 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v264))))
					if v265 != 0 {
						v267 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v264)+2)))
						v276 = v265
						v303 = v267
						v305 = int32(0)
						v306 = v184
						for {
							if v276&int32(1) == int32(0) {
								v328 = v305
							} else {
								v314 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v306))))
								v315 = m.G24
								v317 = int32(1)
								v326 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v315+(v314^(int32(0)-v303&v317)&int32(255))<<(uint(v317)%32)))))
								v328 = v305 + v326
							}
							v330 = int32(1)
							if base.Ui32(v330) < base.Ui32(v276) {
								v276 = int32(base.Ui32(v276) >> (uint(v330) % 32))
								v303 = int32(base.Ui32(v303) >> (uint(v330) % 32))
								v305 = v328
								v306 = v306 + v330
								continue
							} else {
								break
							}
							break
						}
						v374 = v328
					} else {
						v374 = int32(0)
					}
					v379 = int32(1)
					v382 = v215 + v374
					*(*uint16)(unsafe.Add(mBase, uint32(v203+v248<<(uint(v379)%32)))) = uint16(v382)
					v385 = v248 + v379
					if v385 != int32(68) {
						v248 = v385
						continue
					} else {
						break
					}
					break
				}
				v388 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198)+12)))
				v389 = int32(1)
				v390 = m.G24
				v394 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v390+v388<<(uint(v389)%32)))))
				v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198)+11)))
				v396 = int32(255)
				v401 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v390+(v395^v396)<<(uint(v389)%32)))))
				v402 = v394 + v401
				*(*uint16)(unsafe.Add(mBase, uint32(v203)+136)) = uint16(v402)
				v409 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v390+(v388^v396)<<(uint(v389)%32)))))
				v445 = v389
				for {
					v456 = m.G1
					v461 = v445<<(uint(int32(2))%32) + (v456 + int32(_a_F_VP8CalculateLevelCosts_1)) + int32(-4)
					v462 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v461))))
					if v462 != 0 {
						v464 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v461)+2)))
						v473 = v462
						v500 = v464
						v502 = int32(0)
						v503 = v183
						for {
							if v473&int32(1) == int32(0) {
								v525 = v502
							} else {
								v511 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v503))))
								v512 = m.G24
								v514 = int32(1)
								v523 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v512+(v511^(int32(0)-v500&v514)&int32(255))<<(uint(v514)%32)))))
								v525 = v502 + v523
							}
							v527 = int32(1)
							if base.Ui32(v527) < base.Ui32(v473) {
								v473 = int32(base.Ui32(v473) >> (uint(v527) % 32))
								v500 = int32(base.Ui32(v500) >> (uint(v527) % 32))
								v502 = v525
								v503 = v503 + v527
								continue
							} else {
								break
							}
							break
						}
						v571 = v525
					} else {
						v571 = int32(0)
					}
					v576 = int32(1)
					v579 = v401 + v409 + v571
					*(*uint16)(unsafe.Add(mBase, uint32(v203+int32(136)+v445<<(uint(v576)%32)))) = uint16(v579)
					v582 = v445 + v576
					if v582 != int32(68) {
						v445 = v582
						continue
					} else {
						break
					}
					break
				}
				v585 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198)+23)))
				v586 = int32(1)
				v587 = m.G24
				v591 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v587+v585<<(uint(v586)%32)))))
				v592 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198)+22)))
				v593 = int32(255)
				v598 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v587+(v592^v593)<<(uint(v586)%32)))))
				v599 = v591 + v598
				*(*uint16)(unsafe.Add(mBase, uint32(v203)+272)) = uint16(v599)
				v606 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v587+(v585^v593)<<(uint(v586)%32)))))
				v642 = v586
				for {
					v653 = m.G1
					v658 = v642<<(uint(int32(2))%32) + (v653 + int32(_a_F_VP8CalculateLevelCosts_1)) + int32(-4)
					v659 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v658))))
					if v659 != 0 {
						v661 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v658)+2)))
						v670 = v659
						v697 = v661
						v699 = int32(0)
						v700 = v182
						for {
							if v670&int32(1) == int32(0) {
								v722 = v699
							} else {
								v708 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v700))))
								v709 = m.G24
								v711 = int32(1)
								v720 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v709+(v708^(int32(0)-v697&v711)&int32(255))<<(uint(v711)%32)))))
								v722 = v699 + v720
							}
							v724 = int32(1)
							if base.Ui32(v724) < base.Ui32(v670) {
								v670 = int32(base.Ui32(v670) >> (uint(v724) % 32))
								v697 = int32(base.Ui32(v697) >> (uint(v724) % 32))
								v699 = v722
								v700 = v700 + v724
								continue
							} else {
								break
							}
							break
						}
						v768 = v722
					} else {
						v768 = int32(0)
					}
					v773 = int32(1)
					v776 = v598 + v606 + v768
					*(*uint16)(unsafe.Add(mBase, uint32(v203+int32(272)+v642<<(uint(v773)%32)))) = uint16(v776)
					v779 = v642 + v773
					if v779 != int32(68) {
						v642 = v779
						continue
					} else {
						break
					}
					break
				}
				v782 = int32(33)
				v789 = v185 + int32(1)
				if v789 != int32(8) {
					v182 = v182 + v782
					v183 = v183 + v782
					v184 = v184 + v782
					v185 = v789
					continue
				} else {
					break
				}
				break
			}
			v794 = l0 + int32(_a_F_VP8CalculateLevelCosts_2) + v131*int32(192)
			v795 = v150 + v58*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+180)) = v795
			v797 = v150 + v61*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+168)) = v797
			v799 = v150 + v64*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+156)) = v799
			v801 = v150 + v67*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+144)) = v801
			v803 = v150 + v70*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+132)) = v803
			v805 = v150 + v73*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+120)) = v805
			v807 = v150 + v76*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+108)) = v807
			v809 = v150 + v79*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+96)) = v809
			v811 = v150 + v82*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+84)) = v811
			v813 = v150 + v85*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+72)) = v813
			v815 = v150 + v88*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+60)) = v815
			v817 = v150 + v91*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+48)) = v817
			v819 = v150 + v94*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+36)) = v819
			v821 = v150 + v97*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+24)) = v821
			v823 = v150 + v100*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794)+12)) = v823
			v825 = v150 + v103*v59
			*(*int32)(unsafe.Add(mBase, uint32(v794))) = v825
			v827 = int32(272)
			*(*int32)(unsafe.Add(mBase, uint32(v794)+188)) = v795 + v827
			v830 = int32(136)
			*(*int32)(unsafe.Add(mBase, uint32(v794)+184)) = v795 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+176)) = v797 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+172)) = v797 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+164)) = v799 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+160)) = v799 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+152)) = v801 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+148)) = v801 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+140)) = v803 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+136)) = v803 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+128)) = v805 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+124)) = v805 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+116)) = v807 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+112)) = v807 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+104)) = v809 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+100)) = v809 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+92)) = v811 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+88)) = v811 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+80)) = v813 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+76)) = v813 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+68)) = v815 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+64)) = v815 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+56)) = v817 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+52)) = v817 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+44)) = v819 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+40)) = v819 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+32)) = v821 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+28)) = v821 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+20)) = v823 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+16)) = v823 + v830
			*(*int32)(unsafe.Add(mBase, uint32(v794)+8)) = v825 + v827
			*(*int32)(unsafe.Add(mBase, uint32(v794)+4)) = v825 + v830
			v923 = int32(264)
			v930 = v131 + int32(1)
			if v930 != int32(4) {
				v108 = v108 + v923
				v109 = v109 + v923
				v110 = v110 + v923
				v131 = v930
				continue
			} else {
				break
			}
			break
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CalculateLevelCosts[0]))) = int32(0)
	}
	return
}
func F_VP8CodeIntraModes(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v44 int32
	_ = v44
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v65 int64
	_ = v65
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v74 int32
	_ = v74
	var v79 int32
	_ = v79
	var v80 int64
	_ = v80
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v129 int32
	_ = v129
	var v135 int32
	_ = v135
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v159 int32
	_ = v159
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v187 int32
	_ = v187
	var v189 int32
	_ = v189
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v204 int32
	_ = v204
	var v207 int32
	_ = v207
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v217 int32
	_ = v217
	var v221 int32
	_ = v221
	var v225 int32
	_ = v225
	var v229 int32
	_ = v229
	var v231 int32
	_ = v231
	var v235 int32
	_ = v235
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v247 int32
	_ = v247
	var v250 int32
	_ = v250
	var v255 int32
	_ = v255
	var v257 int32
	_ = v257
	var v260 int32
	_ = v260
	var v263 int32
	_ = v263
	var v265 int32
	_ = v265
	var v269 int32
	_ = v269
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v279 int32
	_ = v279
	var v283 int32
	_ = v283
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v294 int32
	_ = v294
	var v297 int32
	_ = v297
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v304 int32
	_ = v304
	var v307 int32
	_ = v307
	var v309 int32
	_ = v309
	var v313 int32
	_ = v313
	var v317 int32
	_ = v317
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v327 int32
	_ = v327
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v340 int32
	_ = v340
	var v341 int32
	_ = v341
	var v346 int32
	_ = v346
	var v349 int32
	_ = v349
	var v352 int32
	_ = v352
	var v355 int32
	_ = v355
	var v357 int32
	_ = v357
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v375 int32
	_ = v375
	var v378 int32
	_ = v378
	var v379 int32
	_ = v379
	var v386 int32
	_ = v386
	var v390 int32
	_ = v390
	var v392 int32
	_ = v392
	var v395 int32
	_ = v395
	var v398 int32
	_ = v398
	var v400 int32
	_ = v400
	var v404 int32
	_ = v404
	var v408 int32
	_ = v408
	var v412 int32
	_ = v412
	var v414 int32
	_ = v414
	var v418 int32
	_ = v418
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v429 int32
	_ = v429
	var v431 int32
	_ = v431
	var v437 int32
	_ = v437
	var v441 int32
	_ = v441
	var v451 int32
	_ = v451
	var v459 int32
	_ = v459
	var v460 int32
	_ = v460
	var v468 int32
	_ = v468
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v478 int32
	_ = v478
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v491 int32
	_ = v491
	var v493 int32
	_ = v493
	var v497 int32
	_ = v497
	var v501 int32
	_ = v501
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v511 int32
	_ = v511
	var v514 int32
	_ = v514
	var v515 int32
	_ = v515
	var v525 int32
	_ = v525
	var v526 int32
	_ = v526
	var v528 int32
	_ = v528
	var v531 int32
	_ = v531
	var v534 int32
	_ = v534
	var v536 int32
	_ = v536
	var v540 int32
	_ = v540
	var v544 int32
	_ = v544
	var v548 int32
	_ = v548
	var v550 int32
	_ = v550
	var v554 int32
	_ = v554
	var v557 int32
	_ = v557
	var v558 int32
	_ = v558
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v571 int32
	_ = v571
	var v574 int32
	_ = v574
	var v577 int32
	_ = v577
	var v579 int32
	_ = v579
	var v583 int32
	_ = v583
	var v587 int32
	_ = v587
	var v591 int32
	_ = v591
	var v593 int32
	_ = v593
	var v597 int32
	_ = v597
	var v600 int32
	_ = v600
	var v601 int32
	_ = v601
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v614 int32
	_ = v614
	var v617 int32
	_ = v617
	var v620 int32
	_ = v620
	var v622 int32
	_ = v622
	var v626 int32
	_ = v626
	var v630 int32
	_ = v630
	var v634 int32
	_ = v634
	var v636 int32
	_ = v636
	var v640 int32
	_ = v640
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v652 int32
	_ = v652
	var v653 int32
	_ = v653
	var v655 int32
	_ = v655
	var v658 int32
	_ = v658
	var v661 int32
	_ = v661
	var v663 int32
	_ = v663
	var v667 int32
	_ = v667
	var v671 int32
	_ = v671
	var v675 int32
	_ = v675
	var v677 int32
	_ = v677
	var v681 int32
	_ = v681
	var v684 int32
	_ = v684
	var v685 int32
	_ = v685
	var v697 int32
	_ = v697
	var v698 int32
	_ = v698
	var v700 int32
	_ = v700
	var v703 int32
	_ = v703
	var v706 int32
	_ = v706
	var v708 int32
	_ = v708
	var v712 int32
	_ = v712
	var v716 int32
	_ = v716
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v726 int32
	_ = v726
	var v729 int32
	_ = v729
	var v730 int32
	_ = v730
	var v740 int32
	_ = v740
	var v741 int32
	_ = v741
	var v743 int32
	_ = v743
	var v746 int32
	_ = v746
	var v749 int32
	_ = v749
	var v751 int32
	_ = v751
	var v755 int32
	_ = v755
	var v759 int32
	_ = v759
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v769 int32
	_ = v769
	var v772 int32
	_ = v772
	var v773 int32
	_ = v773
	var v782 int32
	_ = v782
	var v784 int32
	_ = v784
	var v785 int32
	_ = v785
	var v788 int32
	_ = v788
	var v790 int32
	_ = v790
	var v793 int32
	_ = v793
	var v796 int32
	_ = v796
	var v798 int32
	_ = v798
	var v802 int32
	_ = v802
	var v806 int32
	_ = v806
	var v810 int32
	_ = v810
	var v812 int32
	_ = v812
	var v816 int32
	_ = v816
	var v819 int32
	_ = v819
	var v820 int32
	_ = v820
	var v830 int32
	_ = v830
	var v835 int32
	_ = v835
	var v853 int32
	_ = v853
	var v857 int32
	_ = v857
	var v858 int32
	_ = v858
	var v859 int32
	_ = v859
	var v862 int32
	_ = v862
	var v865 int32
	_ = v865
	var v868 int32
	_ = v868
	var v870 int32
	_ = v870
	var v874 int32
	_ = v874
	var v878 int32
	_ = v878
	var v882 int32
	_ = v882
	var v884 int32
	_ = v884
	var v888 int32
	_ = v888
	var v891 int32
	_ = v891
	var v892 int32
	_ = v892
	var v902 int32
	_ = v902
	var v905 int32
	_ = v905
	var v908 int32
	_ = v908
	var v911 int32
	_ = v911
	var v913 int32
	_ = v913
	var v917 int32
	_ = v917
	var v921 int32
	_ = v921
	var v925 int32
	_ = v925
	var v927 int32
	_ = v927
	var v931 int32
	_ = v931
	var v934 int32
	_ = v934
	var v935 int32
	_ = v935
	var v948 int32
	_ = v948
	var v951 int32
	_ = v951
	var v954 int32
	_ = v954
	var v956 int32
	_ = v956
	var v960 int32
	_ = v960
	var v964 int32
	_ = v964
	var v968 int32
	_ = v968
	var v970 int32
	_ = v970
	var v974 int32
	_ = v974
	var v977 int32
	_ = v977
	var v978 int32
	_ = v978
	var v989 int32
	_ = v989
	var v991 int32
	_ = v991
	var v993 int32
	_ = v993
	var v994 int32
	_ = v994
	var v998 int32
	_ = v998
	var v1000 int64
	_ = v1000
	var v1002 int32
	_ = v1002
	var v1004 int32
	_ = v1004
	var v1006 int32
	_ = v1006
	var v1007 int32
	_ = v1007
	var v1009 int32
	_ = v1009
	var v1013 int32
	_ = v1013
	var v1019 int32
	_ = v1019
	var v1020 int32
	_ = v1020
	var v1029 int32
	_ = v1029
	var v1036 int32
	_ = v1036
	var v1038 int32
	_ = v1038
	var v1039 int32
	_ = v1039
	var v1042 int32
	_ = v1042
	var v1046 int32
	_ = v1046
	var v1047 int64
	_ = v1047
	var v1053 int32
	_ = v1053
	var v1056 int32
	_ = v1056
	var v1059 int32
	_ = v1059
	var v1061 int32
	_ = v1061
	var v1066 int32
	_ = v1066
	var v1067 int32
	_ = v1067
	var v1070 int32
	_ = v1070
	var v1074 int32
	_ = v1074
	var v1078 int32
	_ = v1078
	var v1079 int32
	_ = v1079
	var v1082 int32
	_ = v1082
	var v1089 int32
	_ = v1089
	v2 = int32(0)
	v16 = m.G0
	v18 = v16 - int32(3776)
	m.G0 = v18
	*(*int64)(unsafe.Add(mBase, uint32(v18))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v18)+24)) = l0
	v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[0])))
	*(*int32)(unsafe.Add(mBase, uint32(v18)+280)) = v25
	v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	*(*int32)(unsafe.Add(mBase, uint32(v18)+296)) = v27
	v29 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[1])))
	*(*int32)(unsafe.Add(mBase, uint32(v18)+304)) = v29
	*(*int32)(unsafe.Add(mBase, uint32(v18)+32)) = l0 + int32(88)
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v18)+36)) = v34
	v38 = int32(-32)
	v39 = (v18 + int32(447)) & v38
	*(*int32)(unsafe.Add(mBase, uint32(v18)+8)) = v39
	v44 = (v18 + int32(360)) & v38
	*(*int32)(unsafe.Add(mBase, uint32(v18)+308)) = v44
	*(*int32)(unsafe.Add(mBase, uint32(v18)+20)) = v39 + int32(1536)
	*(*int32)(unsafe.Add(mBase, uint32(v18)+16)) = v39 + int32(1024)
	*(*int32)(unsafe.Add(mBase, uint32(v18)+12)) = v39 + int32(512)
	*(*int32)(unsafe.Add(mBase, uint32(v18)+316)) = v44 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(v18)+312)) = v44 + int32(32)
	v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v18)+40)) = v61
	v63 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[4])))
	*(*int32)(unsafe.Add(mBase, uint32(v18)+28)) = v63
	v65 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[5])))
	*(*int64)(unsafe.Add(mBase, uint32(v18)+320)) = v65
	v67 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v44)+47)) = uint8(v67)
	v69 = *(*int32)(unsafe.Add(mBase, uint32(v18)+312))
	v70 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v69+v70))) = uint8(v67)
	v74 = *(*int32)(unsafe.Add(mBase, uint32(v18)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v74+v70))) = uint8(v67)
	v79 = *(*int32)(unsafe.Add(mBase, uint32(v18)+308))
	v80 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v79))) = v80
	*(*int64)(unsafe.Add(mBase, uint32(v79+int32(8)))) = v80
	v86 = *(*int32)(unsafe.Add(mBase, uint32(v18)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v86))) = v80
	v89 = *(*int32)(unsafe.Add(mBase, uint32(v18)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v89))) = v80
	*(*int32)(unsafe.Add(mBase, uint32(v18)+160)) = v2
	v94 = *(*int32)(unsafe.Add(mBase, uint32(v18)+304))
	if v94 == v2 {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(v18)+300)) = int32(0)
	}
	v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v100 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v101 = v99 * v100
	*(*int32)(unsafe.Add(mBase, uint32(v18)+292)) = v101
	*(*int32)(unsafe.Add(mBase, uint32(v18)+288)) = v101
	v104 = *(*int32)(unsafe.Add(mBase, uint32(v18)+24))
	v105 = *(*int32)(unsafe.Add(mBase, uint32(v104)+uint32(_c_F_VP8CodeIntraModes[5])))
	v107 = *(*int32)(unsafe.Add(mBase, uint32(v104)+40))
	v110 = F_memset(m, v105, int32(127), v107<<(uint(int32(5))%32))
	mBase = m.M
	v111 = *(*int32)(unsafe.Add(mBase, uint32(v104)+uint32(_c_F_VP8CodeIntraModes[3])))
	v112 = int32(0)
	v113 = *(*int32)(unsafe.Add(mBase, uint32(v104)+40))
	v116 = F_memset(m, v111, v112, v113<<(uint(int32(2))%32))
	mBase = m.M
	v117 = *(*int32)(unsafe.Add(mBase, uint32(v104)+uint32(_c_F_VP8CodeIntraModes[1])))
	if v117 == v112 {
	} else {
		v121 = *(*int32)(unsafe.Add(mBase, uint32(v104)+40))
		v124 = F_memset(m, v117, int32(0), v121<<(uint(int32(2))%32))
		mBase = m.M
	}
	v127 = int32(0)
	v129 = F_memset(m, v18+int32(168), v127, int32(96))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v18)+284)) = v127
	v135 = l0 + int32(56)
	for {
		v151 = *(*int32)(unsafe.Add(mBase, uint32(v18)+36))
		v152 = *(*int32)(unsafe.Add(mBase, uint32(v18)+28))
		v153 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
		if v153 == int32(0) {
		} else {
			v156 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v152))))
			v159 = int32(1)
			v164 = int32(base.Ui32(v156)>>(uint(int32(6))%32)) & v159
			v165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3416)))
			v167 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
			v170 = v167 * v165 >> (uint(int32(8)) % 32)
			if v164 == int32(0) {
				v179 = v170
			} else {
				v173 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v175 = v170 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v173 + v175
				v179 = v167 - v175
			}
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v179
			if int32(126) < v179 {
			} else {
				v183 = m.G1
				v187 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v183+int32(_a_F_VP8CodeIntraModes_0)+v179))))
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v187
				v189 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v193 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v183+int32(_a_F_VP8CodeIntraModes_1)+v179))))
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v189 << (uint(v193) % 32)
				v196 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
				v197 = v193 + v196
				*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v197
				if v197 < int32(1) {
				} else {
					F_Flush(m, v135)
					mBase = m.M
				}
			}
			v204 = int32(0)
			v207 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3417)+base.B2i32(v164 != v204)))))
			v209 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
			v212 = v209 * v207 >> (uint(int32(8)) % 32)
			if int32(base.Ui32(v156)>>(uint(int32(5))%32))&v159 == v204 {
				v221 = v212
			} else {
				v215 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v217 = v212 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v215 + v217
				v221 = v209 - v217
			}
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v221
			if int32(126) < v221 {
			} else {
				v225 = m.G1
				v229 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v225+int32(_a_F_VP8CodeIntraModes_0)+v221))))
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v229
				v231 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v235 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v225+int32(_a_F_VP8CodeIntraModes_1)+v221))))
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v231 << (uint(v235) % 32)
				v238 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
				v239 = v235 + v238
				*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v239
				if v239 < int32(1) {
				} else {
					F_Flush(m, v135)
					mBase = m.M
				}
			}
		}
		v247 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8CodeIntraModes[6])))
		if v247 == int32(0) {
		} else {
			v250 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v152))))
			v255 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3419)))
			v257 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
			v260 = v257 * v255 >> (uint(int32(8)) % 32)
			if int32(base.Ui32(v250)>>(uint(int32(4))%32))&int32(1) == int32(0) {
				v269 = v260
			} else {
				v263 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v265 = v260 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v263 + v265
				v269 = v257 - v265
			}
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v269
			if int32(126) < v269 {
			} else {
				v273 = m.G1
				v277 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v273+int32(_a_F_VP8CodeIntraModes_0)+v269))))
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v277
				v279 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v283 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v273+int32(_a_F_VP8CodeIntraModes_1)+v269))))
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v279 << (uint(v283) % 32)
				v286 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
				v287 = v283 + v286
				*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v287
				if v287 < int32(1) {
				} else {
					F_Flush(m, v135)
					mBase = m.M
				}
			}
		}
		v294 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v152))))
		v297 = int32(0)
		v298 = base.B2i32(v294&int32(3) != v297)
		v301 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
		v304 = v301 * int32(145) >> (uint(int32(8)) % 32)
		if v298 == v297 {
			v313 = v304
		} else {
			v307 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
			v309 = v304 + int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v307 + v309
			v313 = v301 - v309
		}
		*(*int32)(unsafe.Add(mBase, uint32(v135))) = v313
		if int32(126) < v313 {
		} else {
			v317 = m.G1
			v321 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v317+int32(_a_F_VP8CodeIntraModes_0)+v313))))
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v321
			v323 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
			v327 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v317+int32(_a_F_VP8CodeIntraModes_1)+v313))))
			*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v323 << (uint(v327) % 32)
			v330 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
			v331 = v327 + v330
			*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v331
			if v331 < int32(1) {
			} else {
				F_Flush(m, v135)
				mBase = m.M
			}
		}
		if v298 == int32(0) {
			v429 = int32(0)
			v431 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
			v437 = v151
			v441 = v429
			for {
				v451 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v437+int32(-1)))))
				v459 = int32(0)
				v460 = v451
				for {
					v468 = m.G1
					v472 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v437+v459))))
					v473 = int32(0)
					v474 = base.B2i32(v472 != v473)
					v478 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v437+(v429-v431)+v459))))
					v482 = v460*int32(9) + (v468 + int32(_a_F_VP8CodeIntraModes_2) + v478*int32(90))
					v483 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482))))
					v485 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
					v488 = v485 * v483 >> (uint(int32(8)) % 32)
					if v474 == v473 {
						v497 = v488
					} else {
						v491 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
						v493 = v488 + int32(1)
						*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v491 + v493
						v497 = v485 - v493
					}
					*(*int32)(unsafe.Add(mBase, uint32(v135))) = v497
					if int32(126) < v497 {
					} else {
						v501 = m.G1
						v505 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v501+int32(_a_F_VP8CodeIntraModes_0)+v497))))
						*(*int32)(unsafe.Add(mBase, uint32(v135))) = v505
						v507 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
						v511 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v501+int32(_a_F_VP8CodeIntraModes_1)+v497))))
						*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v507 << (uint(v511) % 32)
						v514 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
						v515 = v511 + v514
						*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v515
						if v515 < int32(1) {
						} else {
							F_Flush(m, v135)
							mBase = m.M
						}
					}
					if v474 == int32(0) {
					} else {
						v525 = base.B2i32(v472 != int32(1))
						v526 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482)+1)))
						v528 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
						v531 = v528 * v526 >> (uint(int32(8)) % 32)
						if v525 == int32(0) {
							v540 = v531
						} else {
							v534 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
							v536 = v531 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v534 + v536
							v540 = v528 - v536
						}
						*(*int32)(unsafe.Add(mBase, uint32(v135))) = v540
						if int32(126) < v540 {
						} else {
							v544 = m.G1
							v548 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v544+int32(_a_F_VP8CodeIntraModes_0)+v540))))
							*(*int32)(unsafe.Add(mBase, uint32(v135))) = v548
							v550 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
							v554 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v544+int32(_a_F_VP8CodeIntraModes_1)+v540))))
							*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v550 << (uint(v554) % 32)
							v557 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
							v558 = v554 + v557
							*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v558
							if v558 < int32(1) {
							} else {
								F_Flush(m, v135)
								mBase = m.M
							}
						}
						if v525 == int32(0) {
						} else {
							v568 = base.B2i32(v472 != int32(2))
							v569 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482)+2)))
							v571 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
							v574 = v571 * v569 >> (uint(int32(8)) % 32)
							if v568 == int32(0) {
								v583 = v574
							} else {
								v577 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
								v579 = v574 + int32(1)
								*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v577 + v579
								v583 = v571 - v579
							}
							*(*int32)(unsafe.Add(mBase, uint32(v135))) = v583
							if int32(126) < v583 {
							} else {
								v587 = m.G1
								v591 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v587+int32(_a_F_VP8CodeIntraModes_0)+v583))))
								*(*int32)(unsafe.Add(mBase, uint32(v135))) = v591
								v593 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
								v597 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v587+int32(_a_F_VP8CodeIntraModes_1)+v583))))
								*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v593 << (uint(v597) % 32)
								v600 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
								v601 = v597 + v600
								*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v601
								if v601 < int32(1) {
								} else {
									F_Flush(m, v135)
									mBase = m.M
								}
							}
							if v568 == int32(0) {
							} else {
								v611 = base.B2i32(base.Ui32(int32(5)) < base.Ui32(v472))
								v612 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482)+3)))
								v614 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
								v617 = v614 * v612 >> (uint(int32(8)) % 32)
								if v611 == int32(0) {
									v626 = v617
								} else {
									v620 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
									v622 = v617 + int32(1)
									*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v620 + v622
									v626 = v614 - v622
								}
								*(*int32)(unsafe.Add(mBase, uint32(v135))) = v626
								if int32(126) < v626 {
								} else {
									v630 = m.G1
									v634 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v630+int32(_a_F_VP8CodeIntraModes_0)+v626))))
									*(*int32)(unsafe.Add(mBase, uint32(v135))) = v634
									v636 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
									v640 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v630+int32(_a_F_VP8CodeIntraModes_1)+v626))))
									*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v636 << (uint(v640) % 32)
									v643 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
									v644 = v640 + v643
									*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v644
									if v644 < int32(1) {
									} else {
										F_Flush(m, v135)
										mBase = m.M
									}
								}
								if base.Ui32(int32(5)) < base.Ui32(v472) {
									v697 = base.B2i32(v472 != int32(6))
									v698 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482)+6)))
									v700 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
									v703 = v700 * v698 >> (uint(int32(8)) % 32)
									if v697 == int32(0) {
										v712 = v703
									} else {
										v706 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
										v708 = v703 + int32(1)
										*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v706 + v708
										v712 = v700 - v708
									}
									*(*int32)(unsafe.Add(mBase, uint32(v135))) = v712
									if int32(126) < v712 {
									} else {
										v716 = m.G1
										v720 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v716+int32(_a_F_VP8CodeIntraModes_0)+v712))))
										*(*int32)(unsafe.Add(mBase, uint32(v135))) = v720
										v722 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
										v726 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v716+int32(_a_F_VP8CodeIntraModes_1)+v712))))
										*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v722 << (uint(v726) % 32)
										v729 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
										v730 = v726 + v729
										*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v730
										if v730 < int32(1) {
										} else {
											F_Flush(m, v135)
											mBase = m.M
										}
									}
									if v697 == int32(0) {
									} else {
										v740 = base.B2i32(v472 != int32(7))
										v741 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482)+7)))
										v743 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
										v746 = v743 * v741 >> (uint(int32(8)) % 32)
										if v740 == int32(0) {
											v755 = v746
										} else {
											v749 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
											v751 = v746 + int32(1)
											*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v749 + v751
											v755 = v743 - v751
										}
										*(*int32)(unsafe.Add(mBase, uint32(v135))) = v755
										if int32(126) < v755 {
										} else {
											v759 = m.G1
											v763 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v759+int32(_a_F_VP8CodeIntraModes_0)+v755))))
											*(*int32)(unsafe.Add(mBase, uint32(v135))) = v763
											v765 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
											v769 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v759+int32(_a_F_VP8CodeIntraModes_1)+v755))))
											*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v765 << (uint(v769) % 32)
											v772 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
											v773 = v769 + v772
											*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v773
											if v773 < int32(1) {
											} else {
												F_Flush(m, v135)
												mBase = m.M
											}
										}
										if v740 == int32(0) {
										} else {
											v782 = int32(8)
											v784 = v782
											v785 = v782
											v788 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482+v784))))
											v790 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
											v793 = v790 * v788 >> (uint(int32(8)) % 32)
											if base.B2i32(v785 != v472) == int32(0) {
												v802 = v793
											} else {
												v796 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
												v798 = v793 + int32(1)
												*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v796 + v798
												v802 = v790 - v798
											}
											*(*int32)(unsafe.Add(mBase, uint32(v135))) = v802
											if int32(126) < v802 {
											} else {
												v806 = m.G1
												v810 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v806+int32(_a_F_VP8CodeIntraModes_0)+v802))))
												*(*int32)(unsafe.Add(mBase, uint32(v135))) = v810
												v812 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
												v816 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v806+int32(_a_F_VP8CodeIntraModes_1)+v802))))
												*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v812 << (uint(v816) % 32)
												v819 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
												v820 = v816 + v819
												*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v820
												if v820 < int32(1) {
												} else {
													F_Flush(m, v135)
													mBase = m.M
												}
											}
										}
									}
								} else {
									v652 = base.B2i32(v472 != int32(3))
									v653 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482)+4)))
									v655 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
									v658 = v655 * v653 >> (uint(int32(8)) % 32)
									if v652 == int32(0) {
										v667 = v658
									} else {
										v661 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
										v663 = v658 + int32(1)
										*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v661 + v663
										v667 = v655 - v663
									}
									*(*int32)(unsafe.Add(mBase, uint32(v135))) = v667
									if int32(126) < v667 {
									} else {
										v671 = m.G1
										v675 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v671+int32(_a_F_VP8CodeIntraModes_0)+v667))))
										*(*int32)(unsafe.Add(mBase, uint32(v135))) = v675
										v677 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
										v681 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v671+int32(_a_F_VP8CodeIntraModes_1)+v667))))
										*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v677 << (uint(v681) % 32)
										v684 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
										v685 = v681 + v684
										*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v685
										if v685 < int32(1) {
										} else {
											F_Flush(m, v135)
											mBase = m.M
										}
									}
									if v652 == int32(0) {
									} else {
										v784 = int32(5)
										v785 = int32(4)
										v788 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v482+v784))))
										v790 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
										v793 = v790 * v788 >> (uint(int32(8)) % 32)
										if base.B2i32(v785 != v472) == int32(0) {
											v802 = v793
										} else {
											v796 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
											v798 = v793 + int32(1)
											*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v796 + v798
											v802 = v790 - v798
										}
										*(*int32)(unsafe.Add(mBase, uint32(v135))) = v802
										if int32(126) < v802 {
										} else {
											v806 = m.G1
											v810 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v806+int32(_a_F_VP8CodeIntraModes_0)+v802))))
											*(*int32)(unsafe.Add(mBase, uint32(v135))) = v810
											v812 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
											v816 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v806+int32(_a_F_VP8CodeIntraModes_1)+v802))))
											*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v812 << (uint(v816) % 32)
											v819 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
											v820 = v816 + v819
											*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v820
											if v820 < int32(1) {
											} else {
												F_Flush(m, v135)
												mBase = m.M
											}
										}
									}
								}
							}
						}
					}
					v830 = v459 + int32(1)
					if v830 != int32(4) {
						v459 = v830
						v460 = v472
						continue
					} else {
						break
					}
					break
				}
				v835 = v441 + int32(1)
				if v835 != int32(4) {
					v437 = v437 + v431
					v441 = v835
					continue
				} else {
					break
				}
				break
			}
		} else {
			v340 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v151))))
			v341 = int32(1)
			v346 = base.B2i32(v340&int32(253) == v341)
			v349 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
			v352 = v349 * int32(156) >> (uint(int32(8)) % 32)
			if v346 == int32(0) {
				v361 = v352
			} else {
				v355 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v357 = v352 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v355 + v357
				v361 = v349 - v357
			}
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v361
			if int32(126) < v361 {
			} else {
				v365 = m.G1
				v369 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v365+int32(_a_F_VP8CodeIntraModes_0)+v361))))
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v369
				v371 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v375 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v365+int32(_a_F_VP8CodeIntraModes_1)+v361))))
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v371 << (uint(v375) % 32)
				v378 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
				v379 = v375 + v378
				*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v379
				if v379 < int32(1) {
				} else {
					F_Flush(m, v135)
					mBase = m.M
				}
			}
			if v340&int32(253) == v341 {
				v386 = v341
			} else {
				v386 = int32(2)
			}
			if v340&int32(253) == v341 {
				v390 = int32(128)
			} else {
				v390 = int32(163)
			}
			v392 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
			v395 = v392 * v390 >> (uint(int32(8)) % 32)
			if base.B2i32(v340 == v386) == int32(0) {
				v404 = v395
			} else {
				v398 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v400 = v395 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v398 + v400
				v404 = v392 - v400
			}
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v404
			if int32(126) < v404 {
			} else {
				v408 = m.G1
				v412 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v408+int32(_a_F_VP8CodeIntraModes_0)+v404))))
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v412
				v414 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v418 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v408+int32(_a_F_VP8CodeIntraModes_1)+v404))))
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v414 << (uint(v418) % 32)
				v421 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
				v422 = v418 + v421
				*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v422
				if v422 < int32(1) {
				} else {
					F_Flush(m, v135)
					mBase = m.M
				}
			}
		}
		v853 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v152))))
		v857 = int32(base.Ui32(v853)>>(uint(int32(2))%32)) & int32(3)
		v858 = int32(0)
		v859 = base.B2i32(v857 != v858)
		v862 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
		v865 = v862 * int32(142) >> (uint(int32(8)) % 32)
		if v859 == v858 {
			v874 = v865
		} else {
			v868 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
			v870 = v865 + int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v868 + v870
			v874 = v862 - v870
		}
		*(*int32)(unsafe.Add(mBase, uint32(v135))) = v874
		if int32(126) < v874 {
		} else {
			v878 = m.G1
			v882 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v878+int32(_a_F_VP8CodeIntraModes_0)+v874))))
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v882
			v884 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
			v888 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v878+int32(_a_F_VP8CodeIntraModes_1)+v874))))
			*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v884 << (uint(v888) % 32)
			v891 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
			v892 = v888 + v891
			*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v892
			if v892 < int32(1) {
			} else {
				F_Flush(m, v135)
				mBase = m.M
			}
		}
		if v859 == int32(0) {
		} else {
			v902 = base.B2i32(v857 != int32(2))
			v905 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
			v908 = v905 * int32(114) >> (uint(int32(8)) % 32)
			if v902 == int32(0) {
				v917 = v908
			} else {
				v911 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v913 = v908 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v911 + v913
				v917 = v905 - v913
			}
			*(*int32)(unsafe.Add(mBase, uint32(v135))) = v917
			if int32(126) < v917 {
			} else {
				v921 = m.G1
				v925 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v921+int32(_a_F_VP8CodeIntraModes_0)+v917))))
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v925
				v927 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
				v931 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v921+int32(_a_F_VP8CodeIntraModes_1)+v917))))
				*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v927 << (uint(v931) % 32)
				v934 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
				v935 = v931 + v934
				*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v935
				if v935 < int32(1) {
				} else {
					F_Flush(m, v135)
					mBase = m.M
				}
			}
			if v902 == int32(0) {
			} else {
				v948 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
				v951 = v948 * int32(183) >> (uint(int32(8)) % 32)
				if base.B2i32(v857 != int32(3)) == int32(0) {
					v960 = v951
				} else {
					v954 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
					v956 = v951 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v954 + v956
					v960 = v948 - v956
				}
				*(*int32)(unsafe.Add(mBase, uint32(v135))) = v960
				if int32(126) < v960 {
				} else {
					v964 = m.G1
					v968 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v964+int32(_a_F_VP8CodeIntraModes_0)+v960))))
					*(*int32)(unsafe.Add(mBase, uint32(v135))) = v968
					v970 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
					v974 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v964+int32(_a_F_VP8CodeIntraModes_1)+v960))))
					*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v970 << (uint(v974) % 32)
					v977 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
					v978 = v974 + v977
					*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v978
					if v978 < int32(1) {
					} else {
						F_Flush(m, v135)
						mBase = m.M
					}
				}
			}
		}
		v989 = *(*int32)(unsafe.Add(mBase, uint32(v18)))
		v991 = v989 + int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(v18))) = v991
		v993 = *(*int32)(unsafe.Add(mBase, uint32(v18)+24))
		v994 = *(*int32)(unsafe.Add(mBase, uint32(v993)+40))
		if v991 != v994 {
			v1066 = *(*int32)(unsafe.Add(mBase, uint32(v18)+36))
			v1067 = int32(4)
			*(*int32)(unsafe.Add(mBase, uint32(v18)+36)) = v1066 + v1067
			v1070 = *(*int32)(unsafe.Add(mBase, uint32(v18)+28))
			*(*int32)(unsafe.Add(mBase, uint32(v18)+28)) = v1070 + v1067
			v1074 = *(*int32)(unsafe.Add(mBase, uint32(v18)+40))
			*(*int32)(unsafe.Add(mBase, uint32(v18)+40)) = v1074 + v1067
			v1078 = *(*int32)(unsafe.Add(mBase, uint32(v18)+320))
			v1079 = int32(16)
			*(*int32)(unsafe.Add(mBase, uint32(v18)+320)) = v1078 + v1079
			v1082 = *(*int32)(unsafe.Add(mBase, uint32(v18)+324))
			*(*int32)(unsafe.Add(mBase, uint32(v18)+324)) = v1082 + v1079
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(v18))) = int32(0)
			v998 = *(*int32)(unsafe.Add(mBase, uint32(v993)+uint32(_c_F_VP8CodeIntraModes[3])))
			*(*int32)(unsafe.Add(mBase, uint32(v18)+40)) = v998
			v1000 = *(*int64)(unsafe.Add(mBase, uint32(v993)+uint32(_c_F_VP8CodeIntraModes[5])))
			*(*int64)(unsafe.Add(mBase, uint32(v18)+320)) = v1000
			v1002 = *(*int32)(unsafe.Add(mBase, uint32(v18)+4))
			v1004 = v1002 + int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v18)+4)) = v1004
			v1006 = *(*int32)(unsafe.Add(mBase, uint32(v993)+uint32(_c_F_VP8CodeIntraModes[2])))
			v1007 = *(*int32)(unsafe.Add(mBase, uint32(v993)+48))
			v1009 = int32(2)
			*(*int32)(unsafe.Add(mBase, uint32(v18)+36)) = v1006 + v1004*v1007<<(uint(v1009)%32)
			v1013 = *(*int32)(unsafe.Add(mBase, uint32(v993)+uint32(_c_F_VP8CodeIntraModes[4])))
			*(*int32)(unsafe.Add(mBase, uint32(v18)+28)) = v1013 + v1004*v991<<(uint(v1009)%32)
			v1019 = *(*int32)(unsafe.Add(mBase, uint32(v993)+52))
			v1020 = int32(-1)
			*(*int32)(unsafe.Add(mBase, uint32(v18)+32)) = v993 + (v1019+v1020)&v1004<<(uint(int32(5))%32) + int32(88)
			v1029 = *(*int32)(unsafe.Add(mBase, uint32(v18)+316))
			if v1020 < v1002 {
				v1036 = int32(-127)
			} else {
				v1036 = int32(127)
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v1029+v1020))) = uint8(v1036)
			v1038 = *(*int32)(unsafe.Add(mBase, uint32(v18)+312))
			v1039 = int32(-1)
			*(*uint8)(unsafe.Add(mBase, uint32(v1038+v1039))) = uint8(v1036)
			v1042 = *(*int32)(unsafe.Add(mBase, uint32(v18)+308))
			*(*uint8)(unsafe.Add(mBase, uint32(v1042+v1039))) = uint8(v1036)
			v1046 = *(*int32)(unsafe.Add(mBase, uint32(v18)+308))
			v1047 = int64(-9114861777597660799)
			*(*int64)(unsafe.Add(mBase, uint32(v1046))) = v1047
			*(*int64)(unsafe.Add(mBase, uint32(v1046+int32(8)))) = v1047
			v1053 = *(*int32)(unsafe.Add(mBase, uint32(v18)+312))
			*(*int64)(unsafe.Add(mBase, uint32(v1053))) = v1047
			v1056 = *(*int32)(unsafe.Add(mBase, uint32(v18)+316))
			*(*int64)(unsafe.Add(mBase, uint32(v1056))) = v1047
			v1059 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(v18)+160)) = v1059
			v1061 = *(*int32)(unsafe.Add(mBase, uint32(v18)+304))
			if v1061 == v1059 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v18)+300)) = int32(0)
			}
		}
		v1089 = *(*int32)(unsafe.Add(mBase, uint32(v18)+288))
		*(*int32)(unsafe.Add(mBase, uint32(v18)+288)) = v1089 + int32(-1)
		if int32(1) < v1089 {
			continue
		} else {
			break
		}
		break
	}
	m.G0 = v18 + int32(3776)
	return
}
func F_VP8DefaultProbas(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	v2 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+3416)) = uint16(v2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8DefaultProbas[0]))) = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3418)))) = uint8(v2)
	v12 = m.G1
	v16 = F_memcpy(m, l0+int32(3420), v12+int32(_a_F_VP8DefaultProbas_0), int32(1056))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8DefaultProbas[1]))) = int32(1)
	return
}
func F_VP8DspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_VP8DspInit[0])))
	v8 = m.G3
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[1]))) = v11 + int32(66)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[2]))) = v11 + int32(67)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[3]))) = v11 + int32(68)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[4]))) = v11 + int32(69)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[5]))) = v11 + int32(70)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[6]))) = v11 + int32(71)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[7]))) = v11 + int32(72)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[8]))) = v11 + int32(73)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[9]))) = v11 + int32(74)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[10]))) = v11 + int32(75)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[11]))) = v11 + int32(76)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[12]))) = v11 + int32(77)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[13]))) = v11 + int32(78)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[14]))) = v11 + int32(79)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[15]))) = v11 + int32(80)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[16]))) = v11 + int32(81)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[17]))) = v11 + int32(82)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[18]))) = v11 + int32(83)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[19]))) = v11 + int32(84)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[20]))) = v11 + int32(85)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[21]))) = v11 + int32(86)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[22]))) = v11 + int32(87)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[23]))) = v11 + int32(88)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[24]))) = v11 + int32(89)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[25]))) = v11 + int32(90)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[26]))) = v11 + int32(91)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[27]))) = v11 + int32(92)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[28]))) = v11 + int32(93)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[29]))) = v11 + int32(94)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[30]))) = v11 + int32(95)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[31]))) = v11 + int32(96)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[32]))) = v11 + int32(97)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[33]))) = v11 + int32(98)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[34]))) = v11 + int32(99)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[35]))) = v11 + int32(100)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[36]))) = v11 + int32(101)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[37]))) = v11 + int32(102)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[38]))) = v11 + int32(103)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[39]))) = v11 + int32(104)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[40]))) = v11 + int32(105)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[41]))) = v11 + int32(106)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[42]))) = v11 + int32(107)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[43]))) = v11 + int32(108)
		v188 = m.G3
		v189 = *(*int32)(unsafe.Add(mBase, uint32(v188)))
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8DspInit[0]))) = v189
	}
	return
}
func F_VP8EmitTokens(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v21 int32
	_ = v21
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v105 int32
	_ = v105
	var v123 int32
	_ = v123
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v187 int32
	_ = v187
	var v189 int32
	_ = v189
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v207 int32
	_ = v207
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v12 == int32(0) {
	} else {
		if l3 == int32(0) {
			v123 = v12
			for {
				v131 = *(*int32)(unsafe.Add(mBase, uint32(v123)))
				if v131 != 0 {
					v133 = int32(0)
				} else {
					v132 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v133 = v132
				}
				v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				if v134 <= v133 {
				} else {
					v148 = v134
					v149 = v123 + v134<<(uint(int32(1))%32) + int32(2)
					for {
						v152 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v149))))
						if v152&int32(_a_F_VP8EmitTokens_0) == int32(0) {
							v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v152&int32(_a_F_VP8EmitTokens_1)))))
							v165 = v164
						} else {
							v165 = v152 & int32(255)
						}
						v167 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v170 = v167 * v165 >> (uint(int32(8)) % 32)
						if int32(base.Ui32(v152)>>(uint(int32(15))%32)) == int32(0) {
							v179 = v170
						} else {
							v173 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
							v175 = v170 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v173 + v175
							v179 = v167 - v175
						}
						*(*int32)(unsafe.Add(mBase, uint32(l1))) = v179
						if int32(126) < v179 {
						} else {
							v183 = m.G1
							v187 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v183+int32(_a_F_VP8EmitTokens_2)+v179))))
							*(*int32)(unsafe.Add(mBase, uint32(l1))) = v187
							v189 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
							v193 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v183+int32(_a_F_VP8EmitTokens_3)+v179))))
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v189 << (uint(v193) % 32)
							v196 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
							v197 = v193 + v196
							*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v197
							if v197 < int32(1) {
							} else {
								F_Flush(m, l1)
								mBase = m.M
							}
						}
						v207 = v148 + int32(-1)
						if v133 < v207 {
							v148 = v207
							v149 = v149 + int32(-2)
							continue
						} else {
							break
						}
						break
					}
				}
				if v131 != 0 {
					v123 = v131
					continue
				} else {
					break
				}
				break
			}
		} else {
			v21 = v12
			for {
				v29 = *(*int32)(unsafe.Add(mBase, uint32(v21)))
				if v29 != 0 {
					v31 = int32(0)
				} else {
					v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v31 = v30
				}
				v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				if v32 <= v31 {
				} else {
					v46 = v32
					v47 = v21 + v32<<(uint(int32(1))%32) + int32(2)
					for {
						v50 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47))))
						if v50&int32(_a_F_VP8EmitTokens_0) == int32(0) {
							v62 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v50&int32(_a_F_VP8EmitTokens_1)))))
							v63 = v62
						} else {
							v63 = v50 & int32(255)
						}
						v65 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v68 = v65 * v63 >> (uint(int32(8)) % 32)
						if int32(base.Ui32(v50)>>(uint(int32(15))%32)) == int32(0) {
							v77 = v68
						} else {
							v71 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
							v73 = v68 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v71 + v73
							v77 = v65 - v73
						}
						*(*int32)(unsafe.Add(mBase, uint32(l1))) = v77
						if int32(126) < v77 {
						} else {
							v81 = m.G1
							v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v81+int32(_a_F_VP8EmitTokens_2)+v77))))
							*(*int32)(unsafe.Add(mBase, uint32(l1))) = v85
							v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
							v91 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v81+int32(_a_F_VP8EmitTokens_3)+v77))))
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v87 << (uint(v91) % 32)
							v94 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
							v95 = v91 + v94
							*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v95
							if v95 < int32(1) {
							} else {
								F_Flush(m, l1)
								mBase = m.M
							}
						}
						v105 = v46 + int32(-1)
						if v31 < v105 {
							v46 = v105
							v47 = v47 + int32(-2)
							continue
						} else {
							break
						}
						break
					}
				}
				F_free(m, v21)
				mBase = m.M
				if v29 != 0 {
					v21 = v29
					continue
				} else {
					break
				}
				break
			}
		}
	}
	if l3 == int32(0) {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(0)
	}
	return int32(1)
}
func F_VP8EstimateTokenSize(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v92 int32
	_ = v92
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v11 != 0 {
		v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		v21 = v11
		v22 = int32(0)
		for {
			v30 = *(*int32)(unsafe.Add(mBase, uint32(v21)))
			if v30 != 0 {
				v32 = int32(0)
			} else {
				v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				v32 = v31
			}
			if v13 <= v32 {
				v82 = v22
			} else {
				v38 = v22
				v42 = v21 + (v13<<(uint(int32(1))%32) + int32(2))
				v43 = v13
				for {
					v45 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v42))))
					if v45&int32(_a_F_VP8EstimateTokenSize_0) == int32(0) {
						v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v45&int32(_a_F_VP8EstimateTokenSize_1)))))
						v67 = v60 ^ int32(base.Ui32(base.I32_extend16_s(v45))>>(uint(int32(15))%32))&int32(255)
					} else {
						v67 = (v45<<(uint(int32(16))%32)>>(uint(int32(31))%32) ^ v45) & int32(255)
					}
					v70 = m.G24
					v74 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v70+v67<<(uint(int32(1))%32)))))
					v75 = v38 + v74
					v77 = v43 + int32(-1)
					if v32 < v77 {
						v38 = v75
						v42 = v42 + int32(-2)
						v43 = v77
						continue
					} else {
						break
					}
					break
				}
				v82 = v75
			}
			if v30 != 0 {
				v21 = v30
				v22 = v82
				continue
			} else {
				break
			}
			break
		}
		v92 = v82
	} else {
		v92 = int32(0)
	}
	return v92
}
func F_VP8FilterStrengthFromDelta(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	v3 = m.G1
	v9 = int32(63)
	if l1 < v9 {
		v12 = l1
	} else {
		v12 = v9
	}
	v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3+int32(_a_F_VP8FilterStrengthFromDelta_0)+l0<<(uint(int32(6))%32)+v12))))
	return v14
}
func F_VP8FiltersInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	v5 = m.G1
	v8 = *(*int32)(unsafe.Add(mBase, uint32(v5)+uint32(_c_F_VP8FiltersInit[0])))
	v9 = m.G3
	v10 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	if v8 == v10 {
	} else {
		v12 = m.G2
		v13 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[1]))) = v12 + int32(131)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[2]))) = v12 + int32(132)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[3]))) = v12 + int32(133)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[4]))) = v12 + int32(134)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[5]))) = v12 + int32(135)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[6]))) = v12 + int32(136)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[7]))) = v12 + int32(137)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[8]))) = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v13)+uint32(_c_F_VP8FiltersInit[0]))) = v10
	}
	return
}
func F_VP8GetCostLuma16(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v42 int32
	_ = v42
	var v52 int32
	_ = v52
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v180 int32
	_ = v180
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v187 int32
	_ = v187
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v195 int32
	_ = v195
	var v200 int32
	_ = v200
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v207 int32
	_ = v207
	var v211 int32
	_ = v211
	var v214 int32
	_ = v214
	var v219 int32
	_ = v219
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v230 int32
	_ = v230
	var v233 int32
	_ = v233
	var v236 int32
	_ = v236
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v247 int32
	_ = v247
	var v253 int32
	_ = v253
	var v257 int32
	_ = v257
	v15 = m.G0
	v17 = v15 - int32(32)
	m.G0 = v17
	v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v25 = *(*int32)(unsafe.Add(mBase, uint32(v22+int32(-4))))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
	v29 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+124)) = int32(base.Ui32(v26)>>(uint(int32(24))%32)) & v29
	v32 = int32(23)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = int32(base.Ui32(v26)>>(uint(v32)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = int32(base.Ui32(v26)>>(uint(int32(22))%32)) & v29
	v42 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = int32(base.Ui32(v26)>>(uint(v42)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = int32(base.Ui32(v26)>>(uint(int32(18))%32)) & v29
	v52 = int32(15)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = int32(base.Ui32(v26)>>(uint(v52)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = int32(base.Ui32(v26)>>(uint(int32(14))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = int32(base.Ui32(v26)>>(uint(int32(13))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(base.Ui32(v26)>>(uint(int32(12))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = int32(base.Ui32(v25)>>(uint(v32)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = int32(base.Ui32(v25)>>(uint(int32(21))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = int32(base.Ui32(v25)>>(uint(v42)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = int32(base.Ui32(v25)>>(uint(int32(17))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+140)) = int32(base.Ui32(v25)>>(uint(v52)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+136)) = int32(base.Ui32(v25)>>(uint(int32(11))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+132)) = int32(base.Ui32(v25)>>(uint(int32(7))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+128)) = int32(base.Ui32(v25)>>(uint(int32(3))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(v17)+28)) = v19 + int32(_a_F_VP8GetCostLuma16_0)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+24)) = v19 + int32(_a_F_VP8GetCostLuma16_1)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+20)) = v19 + int32(3684)
	v121 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+16)) = v121
	v123 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+4)) = v123
	v129 = v17 + int32(4)
	v130 = m.G25
	v131 = *(*int32)(unsafe.Add(mBase, uint32(v130)))
	m.T0[v131].(func(*base.Module, int32, int32))(m, l1+int32(40), v129)
	mBase = m.M
	v133 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
	v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+124))
	v138 = m.G26
	v139 = *(*int32)(unsafe.Add(mBase, uint32(v138)))
	v140 = m.T0[v139].(func(*base.Module, int32, int32) int32)(m, v133+v134, v129)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v17)+28)) = v19 + int32(_a_F_VP8GetCostLuma16_2)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+24)) = v19 + int32(_a_F_VP8GetCostLuma16_3)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+20)) = v19 + int32(3420)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+16)) = v123
	*(*int32)(unsafe.Add(mBase, uint32(v17)+4)) = v121
	v161 = l1 + int32(168)
	v162 = v123
	v163 = v140
	for {
		v172 = l0 + int32(128) + v162
		v173 = *(*int32)(unsafe.Add(mBase, uint32(v172)))
		v174 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		v177 = int32(4)
		v178 = v17 + v177
		v179 = m.G25
		v180 = *(*int32)(unsafe.Add(mBase, uint32(v179)))
		m.T0[v180].(func(*base.Module, int32, int32))(m, v161+int32(-96), v178)
		mBase = m.M
		v185 = m.G26
		v186 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
		v187 = m.T0[v186].(func(*base.Module, int32, int32) int32)(m, v173+v174, v178)
		mBase = m.M
		v188 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
		v189 = int32(-1)
		v191 = int32(31)
		v192 = int32(base.Ui32(v188^v189) >> (uint(v191) % 32))
		*(*int32)(unsafe.Add(mBase, uint32(v172))) = v192
		*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v192
		v195 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
		v200 = *(*int32)(unsafe.Add(mBase, uint32(v179)))
		m.T0[v200].(func(*base.Module, int32, int32))(m, v161+int32(-64), v178)
		mBase = m.M
		v205 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
		v206 = m.T0[v205].(func(*base.Module, int32, int32) int32)(m, v195+v192, v178)
		mBase = m.M
		v207 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
		v211 = int32(base.Ui32(v207^v189) >> (uint(v191) % 32))
		*(*int32)(unsafe.Add(mBase, uint32(v172))) = v211
		*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = v211
		v214 = *(*int32)(unsafe.Add(mBase, uint32(l0)+100))
		v219 = *(*int32)(unsafe.Add(mBase, uint32(v179)))
		m.T0[v219].(func(*base.Module, int32, int32))(m, v161+int32(-32), v178)
		mBase = m.M
		v224 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
		v225 = m.T0[v224].(func(*base.Module, int32, int32) int32)(m, v214+v211, v178)
		mBase = m.M
		v226 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
		v230 = int32(base.Ui32(v226^v189) >> (uint(v191) % 32))
		*(*int32)(unsafe.Add(mBase, uint32(v172))) = v230
		*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = v230
		v233 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
		v236 = *(*int32)(unsafe.Add(mBase, uint32(v179)))
		m.T0[v236].(func(*base.Module, int32, int32))(m, v161, v178)
		mBase = m.M
		v241 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
		v242 = m.T0[v241].(func(*base.Module, int32, int32) int32)(m, v233+v230, v178)
		mBase = m.M
		v243 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
		v247 = int32(base.Ui32(v243^v189) >> (uint(v191) % 32))
		*(*int32)(unsafe.Add(mBase, uint32(v172))) = v247
		*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = v247
		v253 = v242 + (v225 + (v206 + (v187 + v163)))
		v257 = v162 + v177
		if v257 != int32(16) {
			v161 = v161 + int32(128)
			v162 = v257
			v163 = v253
			continue
		} else {
			break
		}
		break
	}
	m.G0 = v17 + int32(32)
	return v253
}
func F_VP8GetCostLuma4(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	v6 = m.G0
	v7 = int32(32)
	v8 = v6 - v7
	m.G0 = v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = int32(0)
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	*(*int32)(unsafe.Add(mBase, uint32(v8)+28)) = v13 + int32(_a_F_VP8GetCostLuma4_0)
	*(*int32)(unsafe.Add(mBase, uint32(v8)+24)) = v13 + int32(_a_F_VP8GetCostLuma4_1)
	*(*int32)(unsafe.Add(mBase, uint32(v8)+20)) = v13 + int32(_a_F_VP8GetCostLuma4_2)
	v28 = *(*int32)(unsafe.Add(mBase, uint32(l0+v10&int32(-4)+int32(128))))
	v29 = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v8)+16)) = v29
	v38 = *(*int32)(unsafe.Add(mBase, uint32(l0+v10&v29<<(uint(int32(2))%32)+int32(92))))
	v40 = v8 + int32(4)
	v41 = m.G25
	v42 = *(*int32)(unsafe.Add(mBase, uint32(v41)))
	m.T0[v42].(func(*base.Module, int32, int32))(m, l1, v40)
	mBase = m.M
	v47 = m.G26
	v48 = *(*int32)(unsafe.Add(mBase, uint32(v47)))
	v49 = m.T0[v48].(func(*base.Module, int32, int32) int32)(m, v28+v38, v40)
	mBase = m.M
	m.G0 = v8 + v7
	return v49
}
func F_VP8GetCostUV(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v42 int32
	_ = v42
	var v52 int32
	_ = v52
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v152 int32
	_ = v152
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v172 int32
	_ = v172
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v191 int32
	_ = v191
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v211 int32
	_ = v211
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v222 int32
	_ = v222
	var v225 int32
	_ = v225
	var v230 int32
	_ = v230
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v241 int32
	_ = v241
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v250 int32
	_ = v250
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v261 int32
	_ = v261
	var v264 int32
	_ = v264
	var v269 int32
	_ = v269
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v280 int32
	_ = v280
	v15 = m.G0
	v17 = v15 - int32(32)
	m.G0 = v17
	v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v25 = *(*int32)(unsafe.Add(mBase, uint32(v22+int32(-4))))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
	v29 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+124)) = int32(base.Ui32(v26)>>(uint(int32(24))%32)) & v29
	v32 = int32(23)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = int32(base.Ui32(v26)>>(uint(v32)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = int32(base.Ui32(v26)>>(uint(int32(22))%32)) & v29
	v42 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = int32(base.Ui32(v26)>>(uint(v42)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = int32(base.Ui32(v26)>>(uint(int32(18))%32)) & v29
	v52 = int32(15)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = int32(base.Ui32(v26)>>(uint(v52)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = int32(base.Ui32(v26)>>(uint(int32(14))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = int32(base.Ui32(v26)>>(uint(int32(13))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(base.Ui32(v26)>>(uint(int32(12))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = int32(base.Ui32(v25)>>(uint(v32)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = int32(base.Ui32(v25)>>(uint(int32(21))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = int32(base.Ui32(v25)>>(uint(v42)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = int32(base.Ui32(v25)>>(uint(int32(17))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+140)) = int32(base.Ui32(v25)>>(uint(v52)%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+136)) = int32(base.Ui32(v25)>>(uint(int32(11))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+132)) = int32(base.Ui32(v25)>>(uint(int32(7))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(l0)+128)) = int32(base.Ui32(v25)>>(uint(int32(3))%32)) & v29
	*(*int32)(unsafe.Add(mBase, uint32(v17)+28)) = v19 + int32(_a_F_VP8GetCostUV_0)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+24)) = v19 + int32(_a_F_VP8GetCostUV_1)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+20)) = v19 + int32(3948)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+16)) = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v17)+4)) = int32(0)
	v125 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
	v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+144))
	v130 = v17 + int32(4)
	v131 = m.G25
	v132 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v132].(func(*base.Module, int32, int32))(m, l1+int32(584), v130)
	mBase = m.M
	v137 = m.G26
	v138 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v139 = m.T0[v138].(func(*base.Module, int32, int32) int32)(m, v126+v125, v130)
	mBase = m.M
	v140 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v141 = int32(-1)
	v143 = int32(31)
	v144 = int32(base.Ui32(v140^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = v144
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = v144
	v147 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
	v152 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v152].(func(*base.Module, int32, int32))(m, l1+int32(616), v130)
	mBase = m.M
	v157 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v158 = m.T0[v157].(func(*base.Module, int32, int32) int32)(m, v147+v144, v130)
	mBase = m.M
	v159 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v163 = int32(base.Ui32(v159^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = v163
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = v163
	v166 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
	v167 = *(*int32)(unsafe.Add(mBase, uint32(l0)+148))
	v172 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v172].(func(*base.Module, int32, int32))(m, l1+int32(648), v130)
	mBase = m.M
	v177 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v178 = m.T0[v177].(func(*base.Module, int32, int32) int32)(m, v167+v166, v130)
	mBase = m.M
	v179 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v183 = int32(base.Ui32(v179^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = v183
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = v183
	v186 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
	v191 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v191].(func(*base.Module, int32, int32))(m, l1+int32(680), v130)
	mBase = m.M
	v196 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v197 = m.T0[v196].(func(*base.Module, int32, int32) int32)(m, v186+v183, v130)
	mBase = m.M
	v198 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v202 = int32(base.Ui32(v198^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = v202
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = v202
	v205 = *(*int32)(unsafe.Add(mBase, uint32(l0)+116))
	v206 = *(*int32)(unsafe.Add(mBase, uint32(l0)+152))
	v211 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v211].(func(*base.Module, int32, int32))(m, l1+int32(712), v130)
	mBase = m.M
	v216 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v217 = m.T0[v216].(func(*base.Module, int32, int32) int32)(m, v206+v205, v130)
	mBase = m.M
	v218 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v222 = int32(base.Ui32(v218^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = v222
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = v222
	v225 = *(*int32)(unsafe.Add(mBase, uint32(l0)+120))
	v230 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v230].(func(*base.Module, int32, int32))(m, l1+int32(744), v130)
	mBase = m.M
	v235 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v236 = m.T0[v235].(func(*base.Module, int32, int32) int32)(m, v225+v222, v130)
	mBase = m.M
	v237 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v241 = int32(base.Ui32(v237^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = v241
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = v241
	v244 = *(*int32)(unsafe.Add(mBase, uint32(l0)+116))
	v245 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
	v250 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v250].(func(*base.Module, int32, int32))(m, l1+int32(776), v130)
	mBase = m.M
	v255 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v256 = m.T0[v255].(func(*base.Module, int32, int32) int32)(m, v245+v244, v130)
	mBase = m.M
	v257 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v261 = int32(base.Ui32(v257^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v261
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = v261
	v264 = *(*int32)(unsafe.Add(mBase, uint32(l0)+120))
	v269 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
	m.T0[v269].(func(*base.Module, int32, int32))(m, l1+int32(808), v130)
	mBase = m.M
	v274 = *(*int32)(unsafe.Add(mBase, uint32(v137)))
	v275 = m.T0[v274].(func(*base.Module, int32, int32) int32)(m, v264+v261, v130)
	mBase = m.M
	v276 = *(*int32)(unsafe.Add(mBase, uint32(v17)+8))
	v280 = int32(base.Ui32(v276^v141) >> (uint(v143) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v280
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = v280
	m.G0 = v17 + int32(32)
	return v275 + (v256 + (v236 + (v217 + (v197 + (v178 + (v158 + v139))))))
}
func F_VP8InitClipTables(m *base.Module) {
	return
}
func F_VP8InitFilter(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
	if v2 == int32(0) {
	} else {
		base.MemoryFill(m, v2, int32(0), int32(2048))
		v130 = m.G1
		v133 = *(*int32)(unsafe.Add(mBase, uint32(v130)+uint32(_c_F_VP8InitFilter[0])))
		v134 = m.G3
		v135 = *(*int32)(unsafe.Add(mBase, uint32(v134)))
		if v133 == v135 {
		} else {
			v137 = m.G2
			v138 = m.G1
			*(*int32)(unsafe.Add(mBase, uint32(v138)+uint32(_c_F_VP8InitFilter[1]))) = v137 + int32(63)
			*(*int32)(unsafe.Add(mBase, uint32(v138)+uint32(_c_F_VP8InitFilter[2]))) = v137 + int32(64)
			*(*int32)(unsafe.Add(mBase, uint32(v138)+uint32(_c_F_VP8InitFilter[3]))) = v137 + int32(65)
			*(*int32)(unsafe.Add(mBase, uint32(v138)+uint32(_c_F_VP8InitFilter[0]))) = v135
		}
	}
	return
}
func F_VP8InitRandom(m *base.Module, l0 int32, l1 float32) {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v12 int32
	_ = v12
	var v18 float32
	_ = v18
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	v8 = m.G1
	v12 = F_memcpy(m, l0+int32(8), v8+int32(_a_F_VP8InitRandom_0), int32(220))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(133143986176)
	v18 = base.F32_mul(l1, float32(256))
	if base.F32_lt(v18, float32(4.2949673e+09))&base.F32_ge(v18, float32(0)) == int32(0) {
		v28 = int32(0)
	} else {
		v26 = base.I32_trunc_f32_u(v18)
		v28 = v26
	}
	if base.F32_gt(l1, float32(1)) != 0 {
		v31 = int32(256)
	} else {
		v31 = v28
	}
	if base.F32_lt(l1, float32(0)) != 0 {
		v34 = int32(0)
	} else {
		v34 = v31
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+228)) = v34
	return
}
func F_VP8InitResidual(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	*(*int32)(unsafe.Add(mBase, uint32(l3)+12)) = l1
	*(*int32)(unsafe.Add(mBase, uint32(l3))) = l0
	*(*int32)(unsafe.Add(mBase, uint32(l3)+24)) = l2 + l1*int32(192) + int32(_a_F_VP8InitResidual_0)
	*(*int32)(unsafe.Add(mBase, uint32(l3)+20)) = l2 + l1*int32(1056) + int32(_a_F_VP8InitResidual_1)
	*(*int32)(unsafe.Add(mBase, uint32(l3)+16)) = l2 + l1*int32(264) + int32(3420)
	return
}
func F_VP8MakeChroma8Preds(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	v2 = int32(0)
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v6 == v2 {
		v10 = v2
	} else {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+312))
		v10 = v9
	}
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v11 == int32(0) {
		v15 = v2
	} else {
		v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+324))
		v15 = v14
	}
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v17 = m.G28
	v18 = *(*int32)(unsafe.Add(mBase, uint32(v17)))
	m.T0[v18].(func(*base.Module, int32, int32, int32))(m, v16, v10, v15)
	mBase = m.M
	return
}
func F_VP8MakeLuma16Preds(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	v2 = int32(0)
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v6 == v2 {
		v10 = v2
	} else {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+308))
		v10 = v9
	}
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v11 == int32(0) {
		v15 = v2
	} else {
		v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+320))
		v15 = v14
	}
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v17 = m.G27
	v18 = *(*int32)(unsafe.Add(mBase, uint32(v17)))
	m.T0[v18].(func(*base.Module, int32, int32, int32))(m, v16, v10, v15)
	mBase = m.M
	return
}
func F_VP8PutBit(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v60 int32
	_ = v60
	var v64 int32
	_ = v64
	var v66 int64
	_ = v66
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v99 int32
	_ = v99
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v144 int32
	_ = v144
	var v150 int32
	_ = v150
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v8 = v5 * l2 >> (uint(int32(8)) % 32)
	if l1 == int32(0) {
		v17 = v8
	} else {
		v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v13 = v8 + int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v11 + v13
		v17 = v5 - v13
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v17
	if int32(126) < v17 {
	} else {
		v21 = m.G1
		v25 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+int32(_a_F_VP8PutBit_0)+v17))))
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v25
		v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v31 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v21+int32(_a_F_VP8PutBit_1)+v17))))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v27 << (uint(v31) % 32)
		v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		v35 = v31 + v34
		*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v35
		if v35 < int32(1) {
		} else {
			v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v45 + int32(-8)
			v49 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v51 = v45 + int32(8)
			v52 = v49 >> (uint(v51) % 32)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v49 - v52<<(uint(v51)%32)
			v56 = int32(255)
			if v52&v56 == v56 {
				v150 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v150 + int32(1)
			} else {
				v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v64 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				v66 = base.I64_extend_i32_u(v60+int32(1)) + base.I64_extend_i32_u(v64)
				if base.Ui64(v66) < base.Ui64(int64(4294967296)) {
					v71 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
					v72 = base.I32_wrap_i64(v66)
					if base.Ui32(v72) <= base.Ui32(v71) {
						v99 = v52 & int32(256)
						if v99 == int32(0) {
						} else {
							if v64 == int32(0) {
							} else {
								v104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v107 = v104 + v64 + int32(-1)
								v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v107))))
								v110 = v108 + int32(1)
								*(*uint8)(unsafe.Add(mBase, uint32(v107))) = uint8(v110)
							}
						}
						v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if v113 < int32(1) {
							v138 = v64
						} else {
							if v99 != 0 {
								v118 = int32(0)
							} else {
								v118 = int32(-1)
							}
							v120 = v64
							for {
								v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								*(*uint8)(unsafe.Add(mBase, uint32(v126+v120))) = uint8(v118)
								v129 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v129 + int32(-1)
								v133 = int32(1)
								v134 = v120 + v133
								if v133 < v129 {
									v120 = v134
									continue
								} else {
									break
								}
								break
							}
							v138 = v134
						}
						v144 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						*(*uint8)(unsafe.Add(mBase, uint32(v144+v138))) = uint8(v52)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v138 + int32(1)
					} else {
						v76 = v71 << (uint(int32(1)) % 32)
						if base.Ui32(v72) < base.Ui32(v76) {
							v78 = v76
						} else {
							v78 = v72
						}
						v79 = int32(1024)
						if base.Ui32(v79) < base.Ui32(v78) {
							v82 = v78
						} else {
							v82 = v79
						}
						v83 = F_WebPSafeMalloc(m, int64(1), v82)
						mBase = m.M
						if v83 != 0 {
							v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
							if v86 == int32(0) {
							} else {
								v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v90 = F_memcpy(m, v83, v89, v86)
								mBase = m.M
							}
							v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							F_WebPSafeFree(m, v91)
							mBase = m.M
							*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v82
							*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v83
							v99 = v52 & int32(256)
							if v99 == int32(0) {
							} else {
								if v64 == int32(0) {
								} else {
									v104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v107 = v104 + v64 + int32(-1)
									v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v107))))
									v110 = v108 + int32(1)
									*(*uint8)(unsafe.Add(mBase, uint32(v107))) = uint8(v110)
								}
							}
							v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v113 < int32(1) {
								v138 = v64
							} else {
								if v99 != 0 {
									v118 = int32(0)
								} else {
									v118 = int32(-1)
								}
								v120 = v64
								for {
									v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v126+v120))) = uint8(v118)
									v129 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v129 + int32(-1)
									v133 = int32(1)
									v134 = v120 + v133
									if v133 < v129 {
										v120 = v134
										continue
									} else {
										break
									}
									break
								}
								v138 = v134
							}
							v144 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v144+v138))) = uint8(v52)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v138 + int32(1)
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
						}
					}
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
				}
			}
		}
	}
	return l1
}
func F_VP8PutBitUniform(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v63 int64
	_ = v63
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v115 int32
	_ = v115
	var v117 int32
	_ = v117
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v135 int32
	_ = v135
	var v141 int32
	_ = v141
	var v147 int32
	_ = v147
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v7 = v5 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
		v16 = v7
	} else {
		v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v12 = v7 + int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v10 + v12
		v16 = v5 - v12
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v16
	if int32(126) < v16 {
	} else {
		v20 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v21 = int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v20 << (uint(v21) % 32)
		v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v24 + v21
		v28 = m.G1
		v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(_a_F_VP8PutBitUniform_0)+v16))))
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v32
		if v24 < int32(0) {
		} else {
			v42 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v42 + int32(-8)
			v46 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v48 = v42 + int32(8)
			v49 = v46 >> (uint(v48) % 32)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v46 - v49<<(uint(v48)%32)
			v53 = int32(255)
			if v49&v53 == v53 {
				v147 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v147 + int32(1)
			} else {
				v57 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				v63 = base.I64_extend_i32_u(v57+int32(1)) + base.I64_extend_i32_u(v61)
				if base.Ui64(v63) < base.Ui64(int64(4294967296)) {
					v68 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
					v69 = base.I32_wrap_i64(v63)
					if base.Ui32(v69) <= base.Ui32(v68) {
						v96 = v49 & int32(256)
						if v96 == int32(0) {
						} else {
							if v61 == int32(0) {
							} else {
								v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v104 = v101 + v61 + int32(-1)
								v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104))))
								v107 = v105 + int32(1)
								*(*uint8)(unsafe.Add(mBase, uint32(v104))) = uint8(v107)
							}
						}
						v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if v110 < int32(1) {
							v135 = v61
						} else {
							if v96 != 0 {
								v115 = int32(0)
							} else {
								v115 = int32(-1)
							}
							v117 = v61
							for {
								v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								*(*uint8)(unsafe.Add(mBase, uint32(v123+v117))) = uint8(v115)
								v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v126 + int32(-1)
								v130 = int32(1)
								v131 = v117 + v130
								if v130 < v126 {
									v117 = v131
									continue
								} else {
									break
								}
								break
							}
							v135 = v131
						}
						v141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						*(*uint8)(unsafe.Add(mBase, uint32(v141+v135))) = uint8(v49)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v135 + int32(1)
					} else {
						v73 = v68 << (uint(int32(1)) % 32)
						if base.Ui32(v69) < base.Ui32(v73) {
							v75 = v73
						} else {
							v75 = v69
						}
						v76 = int32(1024)
						if base.Ui32(v76) < base.Ui32(v75) {
							v79 = v75
						} else {
							v79 = v76
						}
						v80 = F_WebPSafeMalloc(m, int64(1), v79)
						mBase = m.M
						if v80 != 0 {
							v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
							if v83 == int32(0) {
							} else {
								v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v87 = F_memcpy(m, v80, v86, v83)
								mBase = m.M
							}
							v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							F_WebPSafeFree(m, v88)
							mBase = m.M
							*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v79
							*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v80
							v96 = v49 & int32(256)
							if v96 == int32(0) {
							} else {
								if v61 == int32(0) {
								} else {
									v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v104 = v101 + v61 + int32(-1)
									v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104))))
									v107 = v105 + int32(1)
									*(*uint8)(unsafe.Add(mBase, uint32(v104))) = uint8(v107)
								}
							}
							v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v110 < int32(1) {
								v135 = v61
							} else {
								if v96 != 0 {
									v115 = int32(0)
								} else {
									v115 = int32(-1)
								}
								v117 = v61
								for {
									v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v123+v117))) = uint8(v115)
									v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v126 + int32(-1)
									v130 = int32(1)
									v131 = v117 + v130
									if v130 < v126 {
										v117 = v131
										continue
									} else {
										break
									}
									break
								}
								v135 = v131
							}
							v141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v141+v135))) = uint8(v49)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v135 + int32(1)
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
						}
					}
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
				}
			}
		}
	}
	return l1
}
func F_VP8PutBits(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	var v74 int32
	_ = v74
	var v76 int64
	_ = v76
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v109 int32
	_ = v109
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v123 int32
	_ = v123
	var v128 int32
	_ = v128
	var v130 int32
	_ = v130
	var v136 int32
	_ = v136
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v148 int32
	_ = v148
	var v154 int32
	_ = v154
	var v160 int32
	_ = v160
	var v172 int32
	_ = v172
	v14 = int32(1) << (uint(l2+int32(-1)) % 32)
	for {
		v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v19 = v17 >> (uint(int32(1)) % 32)
		if v14&l1 == int32(0) {
			v29 = v19
		} else {
			v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v25 = v19 + int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v23 + v25
			v29 = v17 - v25
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v29
		if int32(126) < v29 {
		} else {
			v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v34 = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v33 << (uint(v34) % 32)
			v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v37 + v34
			v41 = m.G1
			v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v41+int32(_a_F_VP8PutBits_0)+v29))))
			*(*int32)(unsafe.Add(mBase, uint32(l0))) = v45
			if v37 < int32(0) {
			} else {
				v55 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v55 + int32(-8)
				v59 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v61 = v55 + int32(8)
				v62 = v59 >> (uint(v61) % 32)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v59 - v62<<(uint(v61)%32)
				v66 = int32(255)
				if v62&v66 == v66 {
					v160 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v160 + int32(1)
				} else {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v74 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					v76 = base.I64_extend_i32_u(v70+int32(1)) + base.I64_extend_i32_u(v74)
					if base.Ui64(v76) < base.Ui64(int64(4294967296)) {
						v81 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
						v82 = base.I32_wrap_i64(v76)
						if base.Ui32(v82) <= base.Ui32(v81) {
							v109 = v62 & int32(256)
							if v109 == int32(0) {
							} else {
								if v74 == int32(0) {
								} else {
									v114 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v117 = v114 + v74 + int32(-1)
									v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v117))))
									v120 = v118 + int32(1)
									*(*uint8)(unsafe.Add(mBase, uint32(v117))) = uint8(v120)
								}
							}
							v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v123 < int32(1) {
								v148 = v74
							} else {
								if v109 != 0 {
									v128 = int32(0)
								} else {
									v128 = int32(-1)
								}
								v130 = v74
								for {
									v136 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v136+v130))) = uint8(v128)
									v139 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v139 + int32(-1)
									v143 = int32(1)
									v144 = v130 + v143
									if v143 < v139 {
										v130 = v144
										continue
									} else {
										break
									}
									break
								}
								v148 = v144
							}
							v154 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v154+v148))) = uint8(v62)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v148 + int32(1)
						} else {
							v86 = v81 << (uint(int32(1)) % 32)
							if base.Ui32(v82) < base.Ui32(v86) {
								v88 = v86
							} else {
								v88 = v82
							}
							v89 = int32(1024)
							if base.Ui32(v89) < base.Ui32(v88) {
								v92 = v88
							} else {
								v92 = v89
							}
							v93 = F_WebPSafeMalloc(m, int64(1), v92)
							mBase = m.M
							if v93 != 0 {
								v96 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
								if v96 == int32(0) {
								} else {
									v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v100 = F_memcpy(m, v93, v99, v96)
									mBase = m.M
								}
								v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								F_WebPSafeFree(m, v101)
								mBase = m.M
								*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v92
								*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v93
								v109 = v62 & int32(256)
								if v109 == int32(0) {
								} else {
									if v74 == int32(0) {
									} else {
										v114 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										v117 = v114 + v74 + int32(-1)
										v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v117))))
										v120 = v118 + int32(1)
										*(*uint8)(unsafe.Add(mBase, uint32(v117))) = uint8(v120)
									}
								}
								v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								if v123 < int32(1) {
									v148 = v74
								} else {
									if v109 != 0 {
										v128 = int32(0)
									} else {
										v128 = int32(-1)
									}
									v130 = v74
									for {
										v136 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										*(*uint8)(unsafe.Add(mBase, uint32(v136+v130))) = uint8(v128)
										v139 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
										*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v139 + int32(-1)
										v143 = int32(1)
										v144 = v130 + v143
										if v143 < v139 {
											v130 = v144
											continue
										} else {
											break
										}
										break
									}
									v148 = v144
								}
								v154 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								*(*uint8)(unsafe.Add(mBase, uint32(v154+v148))) = uint8(v62)
								*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v148 + int32(1)
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
							}
						}
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
					}
				}
			}
		}
		v172 = int32(1)
		if base.Ui32(v172) < base.Ui32(v14) {
			v14 = int32(base.Ui32(v14) >> (uint(v172) % 32))
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_VP8PutSignedBits(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v44 int32
	_ = v44
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v65 int64
	_ = v65
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v98 int32
	_ = v98
	var v103 int32
	_ = v103
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v137 int32
	_ = v137
	var v143 int32
	_ = v143
	var v149 int32
	_ = v149
	var v165 int32
	_ = v165
	var v173 int32
	_ = v173
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v189 int32
	_ = v189
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v197 int32
	_ = v197
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v215 int32
	_ = v215
	var v219 int32
	_ = v219
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v226 int32
	_ = v226
	var v230 int32
	_ = v230
	var v234 int32
	_ = v234
	var v236 int64
	_ = v236
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v261 int32
	_ = v261
	var v269 int32
	_ = v269
	var v274 int32
	_ = v274
	var v277 int32
	_ = v277
	var v278 int32
	_ = v278
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v296 int32
	_ = v296
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v308 int32
	_ = v308
	var v314 int32
	_ = v314
	var v320 int32
	_ = v320
	var v336 int32
	_ = v336
	var v342 int32
	_ = v342
	var v346 int32
	_ = v346
	var v348 int32
	_ = v348
	var v352 int32
	_ = v352
	var v354 int32
	_ = v354
	var v358 int32
	_ = v358
	var v362 int32
	_ = v362
	var v363 int32
	_ = v363
	var v366 int32
	_ = v366
	var v370 int32
	_ = v370
	var v374 int32
	_ = v374
	var v384 int32
	_ = v384
	var v388 int32
	_ = v388
	var v390 int32
	_ = v390
	var v391 int32
	_ = v391
	var v395 int32
	_ = v395
	var v399 int32
	_ = v399
	var v403 int32
	_ = v403
	var v405 int64
	_ = v405
	var v410 int32
	_ = v410
	var v411 int32
	_ = v411
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v418 int32
	_ = v418
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v425 int32
	_ = v425
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v430 int32
	_ = v430
	var v438 int32
	_ = v438
	var v443 int32
	_ = v443
	var v446 int32
	_ = v446
	var v447 int32
	_ = v447
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v457 int32
	_ = v457
	var v459 int32
	_ = v459
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v477 int32
	_ = v477
	var v483 int32
	_ = v483
	var v489 int32
	_ = v489
	var v501 int32
	_ = v501
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v9 = v7 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
		v18 = v9
	} else {
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v14 = v9 + int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v12 + v14
		v18 = v7 - v14
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v18
	if int32(126) < v18 {
	} else {
		v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v23 = int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v22 << (uint(v23) % 32)
		v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v26 + v23
		v30 = m.G1
		v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v30+int32(_a_F_VP8PutSignedBits_0)+v18))))
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v34
		if v26 < int32(0) {
		} else {
			v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v44 + int32(-8)
			v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v50 = v44 + int32(8)
			v51 = v48 >> (uint(v50) % 32)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v48 - v51<<(uint(v50)%32)
			v55 = int32(255)
			if v51&v55 == v55 {
				v149 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v149 + int32(1)
			} else {
				v59 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v63 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				v65 = base.I64_extend_i32_u(v59+int32(1)) + base.I64_extend_i32_u(v63)
				if base.Ui64(v65) < base.Ui64(int64(4294967296)) {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
					v71 = base.I32_wrap_i64(v65)
					if base.Ui32(v71) <= base.Ui32(v70) {
						v98 = v51 & int32(256)
						if v98 == int32(0) {
						} else {
							if v63 == int32(0) {
							} else {
								v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v106 = v103 + v63 + int32(-1)
								v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v106))))
								v109 = v107 + int32(1)
								*(*uint8)(unsafe.Add(mBase, uint32(v106))) = uint8(v109)
							}
						}
						v112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if v112 < int32(1) {
							v137 = v63
						} else {
							if v98 != 0 {
								v117 = int32(0)
							} else {
								v117 = int32(-1)
							}
							v119 = v63
							for {
								v125 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								*(*uint8)(unsafe.Add(mBase, uint32(v125+v119))) = uint8(v117)
								v128 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v128 + int32(-1)
								v132 = int32(1)
								v133 = v119 + v132
								if v132 < v128 {
									v119 = v133
									continue
								} else {
									break
								}
								break
							}
							v137 = v133
						}
						v143 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						*(*uint8)(unsafe.Add(mBase, uint32(v143+v137))) = uint8(v51)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v137 + int32(1)
					} else {
						v75 = v70 << (uint(int32(1)) % 32)
						if base.Ui32(v71) < base.Ui32(v75) {
							v77 = v75
						} else {
							v77 = v71
						}
						v78 = int32(1024)
						if base.Ui32(v78) < base.Ui32(v77) {
							v81 = v77
						} else {
							v81 = v78
						}
						v82 = F_WebPSafeMalloc(m, int64(1), v81)
						mBase = m.M
						if v82 != 0 {
							v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
							if v85 == int32(0) {
							} else {
								v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v89 = F_memcpy(m, v82, v88, v85)
								mBase = m.M
							}
							v90 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							F_WebPSafeFree(m, v90)
							mBase = m.M
							*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v81
							*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v82
							v98 = v51 & int32(256)
							if v98 == int32(0) {
							} else {
								if v63 == int32(0) {
								} else {
									v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v106 = v103 + v63 + int32(-1)
									v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v106))))
									v109 = v107 + int32(1)
									*(*uint8)(unsafe.Add(mBase, uint32(v106))) = uint8(v109)
								}
							}
							v112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v112 < int32(1) {
								v137 = v63
							} else {
								if v98 != 0 {
									v117 = int32(0)
								} else {
									v117 = int32(-1)
								}
								v119 = v63
								for {
									v125 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v125+v119))) = uint8(v117)
									v128 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v128 + int32(-1)
									v132 = int32(1)
									v133 = v119 + v132
									if v132 < v128 {
										v119 = v133
										continue
									} else {
										break
									}
									break
								}
								v137 = v133
							}
							v143 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
							*(*uint8)(unsafe.Add(mBase, uint32(v143+v137))) = uint8(v51)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v137 + int32(1)
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
						}
					}
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
				}
			}
		}
	}
	if l1 == int32(0) {
	} else {
		if int32(-1) < l1 {
			v336 = int32(1)
			v342 = v336 << (uint(l2) % 32)
			for {
				v346 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v348 = v346 >> (uint(int32(1)) % 32)
				if v342&(l1<<(uint(v336)%32)) == int32(0) {
					v358 = v348
				} else {
					v352 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v354 = v348 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v352 + v354
					v358 = v346 - v354
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v358
				if int32(126) < v358 {
				} else {
					v362 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v363 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v362 << (uint(v363) % 32)
					v366 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v366 + v363
					v370 = m.G1
					v374 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v370+int32(_a_F_VP8PutSignedBits_0)+v358))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v374
					if v366 < int32(0) {
					} else {
						v384 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v384 + int32(-8)
						v388 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
						v390 = v384 + int32(8)
						v391 = v388 >> (uint(v390) % 32)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v388 - v391<<(uint(v390)%32)
						v395 = int32(255)
						if v391&v395 == v395 {
							v489 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v489 + int32(1)
						} else {
							v399 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							v403 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
							v405 = base.I64_extend_i32_u(v399+int32(1)) + base.I64_extend_i32_u(v403)
							if base.Ui64(v405) < base.Ui64(int64(4294967296)) {
								v410 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
								v411 = base.I32_wrap_i64(v405)
								if base.Ui32(v411) <= base.Ui32(v410) {
									v438 = v391 & int32(256)
									if v438 == int32(0) {
									} else {
										if v403 == int32(0) {
										} else {
											v443 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
											v446 = v443 + v403 + int32(-1)
											v447 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v446))))
											v449 = v447 + int32(1)
											*(*uint8)(unsafe.Add(mBase, uint32(v446))) = uint8(v449)
										}
									}
									v452 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									if v452 < int32(1) {
										v477 = v403
									} else {
										if v438 != 0 {
											v457 = int32(0)
										} else {
											v457 = int32(-1)
										}
										v459 = v403
										for {
											v465 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
											*(*uint8)(unsafe.Add(mBase, uint32(v465+v459))) = uint8(v457)
											v468 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
											*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v468 + int32(-1)
											v472 = int32(1)
											v473 = v459 + v472
											if v472 < v468 {
												v459 = v473
												continue
											} else {
												break
											}
											break
										}
										v477 = v473
									}
									v483 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v483+v477))) = uint8(v391)
									*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v477 + int32(1)
								} else {
									v415 = v410 << (uint(int32(1)) % 32)
									if base.Ui32(v411) < base.Ui32(v415) {
										v417 = v415
									} else {
										v417 = v411
									}
									v418 = int32(1024)
									if base.Ui32(v418) < base.Ui32(v417) {
										v421 = v417
									} else {
										v421 = v418
									}
									v422 = F_WebPSafeMalloc(m, int64(1), v421)
									mBase = m.M
									if v422 != 0 {
										v425 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
										if v425 == int32(0) {
										} else {
											v428 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
											v429 = F_memcpy(m, v422, v428, v425)
											mBase = m.M
										}
										v430 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										F_WebPSafeFree(m, v430)
										mBase = m.M
										*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v421
										*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v422
										v438 = v391 & int32(256)
										if v438 == int32(0) {
										} else {
											if v403 == int32(0) {
											} else {
												v443 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
												v446 = v443 + v403 + int32(-1)
												v447 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v446))))
												v449 = v447 + int32(1)
												*(*uint8)(unsafe.Add(mBase, uint32(v446))) = uint8(v449)
											}
										}
										v452 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
										if v452 < int32(1) {
											v477 = v403
										} else {
											if v438 != 0 {
												v457 = int32(0)
											} else {
												v457 = int32(-1)
											}
											v459 = v403
											for {
												v465 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
												*(*uint8)(unsafe.Add(mBase, uint32(v465+v459))) = uint8(v457)
												v468 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
												*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v468 + int32(-1)
												v472 = int32(1)
												v473 = v459 + v472
												if v472 < v468 {
													v459 = v473
													continue
												} else {
													break
												}
												break
											}
											v477 = v473
										}
										v483 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										*(*uint8)(unsafe.Add(mBase, uint32(v483+v477))) = uint8(v391)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v477 + int32(1)
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
									}
								}
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
							}
						}
					}
				}
				v501 = int32(1)
				if base.Ui32(v501) < base.Ui32(v342) {
					v342 = int32(base.Ui32(v342) >> (uint(v501) % 32))
					continue
				} else {
					break
				}
				break
			}
		} else {
			v165 = int32(1)
			v173 = v165 << (uint(l2) % 32)
			for {
				v177 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v179 = v177 >> (uint(int32(1)) % 32)
				if v173&(v165-l1<<(uint(v165)%32)) == int32(0) {
					v189 = v179
				} else {
					v183 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v185 = v179 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v183 + v185
					v189 = v177 - v185
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v189
				if int32(126) < v189 {
				} else {
					v193 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v194 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v193 << (uint(v194) % 32)
					v197 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v197 + v194
					v201 = m.G1
					v205 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v201+int32(_a_F_VP8PutSignedBits_0)+v189))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v205
					if v197 < int32(0) {
					} else {
						v215 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v215 + int32(-8)
						v219 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
						v221 = v215 + int32(8)
						v222 = v219 >> (uint(v221) % 32)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v219 - v222<<(uint(v221)%32)
						v226 = int32(255)
						if v222&v226 == v226 {
							v320 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v320 + int32(1)
						} else {
							v230 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							v234 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
							v236 = base.I64_extend_i32_u(v230+int32(1)) + base.I64_extend_i32_u(v234)
							if base.Ui64(v236) < base.Ui64(int64(4294967296)) {
								v241 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
								v242 = base.I32_wrap_i64(v236)
								if base.Ui32(v242) <= base.Ui32(v241) {
									v269 = v222 & int32(256)
									if v269 == int32(0) {
									} else {
										if v234 == int32(0) {
										} else {
											v274 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
											v277 = v274 + v234 + int32(-1)
											v278 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v277))))
											v280 = v278 + int32(1)
											*(*uint8)(unsafe.Add(mBase, uint32(v277))) = uint8(v280)
										}
									}
									v283 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									if v283 < int32(1) {
										v308 = v234
									} else {
										if v269 != 0 {
											v288 = int32(0)
										} else {
											v288 = int32(-1)
										}
										v290 = v234
										for {
											v296 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
											*(*uint8)(unsafe.Add(mBase, uint32(v296+v290))) = uint8(v288)
											v299 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
											*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v299 + int32(-1)
											v303 = int32(1)
											v304 = v290 + v303
											if v303 < v299 {
												v290 = v304
												continue
											} else {
												break
											}
											break
										}
										v308 = v304
									}
									v314 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*uint8)(unsafe.Add(mBase, uint32(v314+v308))) = uint8(v222)
									*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v308 + int32(1)
								} else {
									v246 = v241 << (uint(int32(1)) % 32)
									if base.Ui32(v242) < base.Ui32(v246) {
										v248 = v246
									} else {
										v248 = v242
									}
									v249 = int32(1024)
									if base.Ui32(v249) < base.Ui32(v248) {
										v252 = v248
									} else {
										v252 = v249
									}
									v253 = F_WebPSafeMalloc(m, int64(1), v252)
									mBase = m.M
									if v253 != 0 {
										v256 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
										if v256 == int32(0) {
										} else {
											v259 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
											v260 = F_memcpy(m, v253, v259, v256)
											mBase = m.M
										}
										v261 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										F_WebPSafeFree(m, v261)
										mBase = m.M
										*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v252
										*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v253
										v269 = v222 & int32(256)
										if v269 == int32(0) {
										} else {
											if v234 == int32(0) {
											} else {
												v274 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
												v277 = v274 + v234 + int32(-1)
												v278 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v277))))
												v280 = v278 + int32(1)
												*(*uint8)(unsafe.Add(mBase, uint32(v277))) = uint8(v280)
											}
										}
										v283 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
										if v283 < int32(1) {
											v308 = v234
										} else {
											if v269 != 0 {
												v288 = int32(0)
											} else {
												v288 = int32(-1)
											}
											v290 = v234
											for {
												v296 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
												*(*uint8)(unsafe.Add(mBase, uint32(v296+v290))) = uint8(v288)
												v299 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
												*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v299 + int32(-1)
												v303 = int32(1)
												v304 = v290 + v303
												if v303 < v299 {
													v290 = v304
													continue
												} else {
													break
												}
												break
											}
											v308 = v304
										}
										v314 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
										*(*uint8)(unsafe.Add(mBase, uint32(v314+v308))) = uint8(v222)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v308 + int32(1)
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
									}
								}
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
							}
						}
					}
				}
				if base.Ui32(int32(2)) <= base.Ui32(v173) {
					v173 = int32(base.Ui32(v173) >> (uint(int32(1)) % 32))
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
func F_VP8RecordCoeffTokens(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int64
	_ = v40
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
	var v52 int64
	_ = v52
	var v53 int32
	_ = v53
	var v54 int64
	_ = v54
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v73 int32
	_ = v73
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v108 int32
	_ = v108
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v119 int32
	_ = v119
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v143 int64
	_ = v143
	var v144 int32
	_ = v144
	var v148 int32
	_ = v148
	var v155 int64
	_ = v155
	var v156 int32
	_ = v156
	var v157 int64
	_ = v157
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v176 int32
	_ = v176
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v199 int32
	_ = v199
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v225 int32
	_ = v225
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v244 int64
	_ = v244
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v256 int64
	_ = v256
	var v257 int32
	_ = v257
	var v258 int64
	_ = v258
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	var v298 int32
	_ = v298
	var v302 int32
	_ = v302
	var v303 int32
	_ = v303
	var v312 int32
	_ = v312
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v323 int32
	_ = v323
	var v326 int32
	_ = v326
	var v327 int32
	_ = v327
	var v328 int64
	_ = v328
	var v329 int32
	_ = v329
	var v333 int32
	_ = v333
	var v340 int64
	_ = v340
	var v341 int32
	_ = v341
	var v342 int64
	_ = v342
	var v352 int32
	_ = v352
	var v354 int32
	_ = v354
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v382 int32
	_ = v382
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v395 int32
	_ = v395
	var v400 int32
	_ = v400
	var v403 int32
	_ = v403
	var v408 int32
	_ = v408
	var v409 int32
	_ = v409
	var v410 int64
	_ = v410
	var v411 int32
	_ = v411
	var v415 int32
	_ = v415
	var v422 int64
	_ = v422
	var v423 int32
	_ = v423
	var v424 int64
	_ = v424
	var v434 int32
	_ = v434
	var v436 int32
	_ = v436
	var v443 int32
	_ = v443
	var v447 int32
	_ = v447
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v453 int32
	_ = v453
	var v464 int32
	_ = v464
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v478 int32
	_ = v478
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v486 int32
	_ = v486
	var v487 int32
	_ = v487
	var v490 int32
	_ = v490
	var v491 int32
	_ = v491
	var v492 int64
	_ = v492
	var v493 int32
	_ = v493
	var v497 int32
	_ = v497
	var v504 int64
	_ = v504
	var v505 int32
	_ = v505
	var v506 int64
	_ = v506
	var v516 int32
	_ = v516
	var v518 int32
	_ = v518
	var v525 int32
	_ = v525
	var v529 int32
	_ = v529
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v533 int32
	_ = v533
	var v535 int32
	_ = v535
	var v546 int32
	_ = v546
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v559 int32
	_ = v559
	var v564 int32
	_ = v564
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int64
	_ = v571
	var v572 int32
	_ = v572
	var v576 int32
	_ = v576
	var v583 int64
	_ = v583
	var v584 int32
	_ = v584
	var v585 int64
	_ = v585
	var v595 int32
	_ = v595
	var v597 int32
	_ = v597
	var v604 int32
	_ = v604
	var v608 int32
	_ = v608
	var v610 int32
	_ = v610
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v614 int32
	_ = v614
	var v625 int32
	_ = v625
	var v629 int32
	_ = v629
	var v630 int32
	_ = v630
	var v638 int32
	_ = v638
	var v643 int32
	_ = v643
	var v648 int32
	_ = v648
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v653 int64
	_ = v653
	var v654 int32
	_ = v654
	var v658 int32
	_ = v658
	var v665 int64
	_ = v665
	var v666 int32
	_ = v666
	var v667 int64
	_ = v667
	var v677 int32
	_ = v677
	var v679 int32
	_ = v679
	var v686 int32
	_ = v686
	var v690 int32
	_ = v690
	var v692 int32
	_ = v692
	var v693 int32
	_ = v693
	var v694 int32
	_ = v694
	var v696 int32
	_ = v696
	var v707 int32
	_ = v707
	var v711 int32
	_ = v711
	var v712 int32
	_ = v712
	var v720 int32
	_ = v720
	var v725 int32
	_ = v725
	var v728 int32
	_ = v728
	var v733 int32
	_ = v733
	var v734 int32
	_ = v734
	var v735 int64
	_ = v735
	var v736 int32
	_ = v736
	var v740 int32
	_ = v740
	var v747 int64
	_ = v747
	var v748 int32
	_ = v748
	var v749 int64
	_ = v749
	var v759 int32
	_ = v759
	var v761 int32
	_ = v761
	var v767 int32
	_ = v767
	var v771 int32
	_ = v771
	var v773 int32
	_ = v773
	var v774 int32
	_ = v774
	var v775 int32
	_ = v775
	var v777 int32
	_ = v777
	var v786 int32
	_ = v786
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v792 int64
	_ = v792
	var v793 int32
	_ = v793
	var v797 int32
	_ = v797
	var v804 int64
	_ = v804
	var v805 int32
	_ = v805
	var v806 int64
	_ = v806
	var v816 int32
	_ = v816
	var v818 int32
	_ = v818
	var v820 int32
	_ = v820
	var v821 int32
	_ = v821
	var v826 int32
	_ = v826
	var v830 int32
	_ = v830
	var v832 int32
	_ = v832
	var v833 int32
	_ = v833
	var v834 int32
	_ = v834
	var v836 int32
	_ = v836
	var v845 int32
	_ = v845
	var v847 int32
	_ = v847
	var v851 int32
	_ = v851
	var v852 int32
	_ = v852
	var v853 int64
	_ = v853
	var v854 int32
	_ = v854
	var v858 int32
	_ = v858
	var v865 int64
	_ = v865
	var v866 int32
	_ = v866
	var v867 int64
	_ = v867
	var v877 int32
	_ = v877
	var v879 int32
	_ = v879
	var v885 int32
	_ = v885
	var v889 int32
	_ = v889
	var v891 int32
	_ = v891
	var v892 int32
	_ = v892
	var v893 int32
	_ = v893
	var v895 int32
	_ = v895
	var v903 int32
	_ = v903
	var v906 int32
	_ = v906
	var v909 int32
	_ = v909
	var v912 int32
	_ = v912
	var v913 int32
	_ = v913
	var v914 int64
	_ = v914
	var v915 int32
	_ = v915
	var v919 int32
	_ = v919
	var v926 int64
	_ = v926
	var v927 int32
	_ = v927
	var v928 int64
	_ = v928
	var v938 int32
	_ = v938
	var v940 int32
	_ = v940
	var v947 int32
	_ = v947
	var v951 int32
	_ = v951
	var v953 int32
	_ = v953
	var v954 int32
	_ = v954
	var v955 int32
	_ = v955
	var v957 int32
	_ = v957
	var v963 int32
	_ = v963
	var v967 int32
	_ = v967
	var v968 int32
	_ = v968
	var v976 int32
	_ = v976
	var v980 int32
	_ = v980
	var v983 int32
	_ = v983
	var v984 int32
	_ = v984
	var v985 int64
	_ = v985
	var v986 int32
	_ = v986
	var v990 int32
	_ = v990
	var v997 int64
	_ = v997
	var v998 int32
	_ = v998
	var v999 int64
	_ = v999
	var v1009 int32
	_ = v1009
	var v1011 int32
	_ = v1011
	var v1018 int32
	_ = v1018
	var v1022 int32
	_ = v1022
	var v1024 int32
	_ = v1024
	var v1025 int32
	_ = v1025
	var v1026 int32
	_ = v1026
	var v1028 int32
	_ = v1028
	var v1034 int32
	_ = v1034
	var v1038 int32
	_ = v1038
	var v1041 int32
	_ = v1041
	var v1042 int32
	_ = v1042
	var v1043 int64
	_ = v1043
	var v1044 int32
	_ = v1044
	var v1048 int32
	_ = v1048
	var v1055 int64
	_ = v1055
	var v1056 int32
	_ = v1056
	var v1057 int64
	_ = v1057
	var v1067 int32
	_ = v1067
	var v1069 int32
	_ = v1069
	var v1076 int32
	_ = v1076
	var v1080 int32
	_ = v1080
	var v1082 int32
	_ = v1082
	var v1083 int32
	_ = v1083
	var v1084 int32
	_ = v1084
	var v1086 int32
	_ = v1086
	var v1092 int32
	_ = v1092
	var v1096 int32
	_ = v1096
	var v1097 int32
	_ = v1097
	var v1105 int32
	_ = v1105
	var v1109 int32
	_ = v1109
	var v1112 int32
	_ = v1112
	var v1113 int32
	_ = v1113
	var v1114 int64
	_ = v1114
	var v1115 int32
	_ = v1115
	var v1119 int32
	_ = v1119
	var v1126 int64
	_ = v1126
	var v1127 int32
	_ = v1127
	var v1128 int64
	_ = v1128
	var v1138 int32
	_ = v1138
	var v1140 int32
	_ = v1140
	var v1147 int32
	_ = v1147
	var v1151 int32
	_ = v1151
	var v1153 int32
	_ = v1153
	var v1154 int32
	_ = v1154
	var v1155 int32
	_ = v1155
	var v1157 int32
	_ = v1157
	var v1165 int32
	_ = v1165
	var v1167 int32
	_ = v1167
	var v1172 int32
	_ = v1172
	var v1173 int32
	_ = v1173
	var v1174 int64
	_ = v1174
	var v1175 int32
	_ = v1175
	var v1179 int32
	_ = v1179
	var v1186 int64
	_ = v1186
	var v1187 int32
	_ = v1187
	var v1188 int64
	_ = v1188
	var v1198 int32
	_ = v1198
	var v1200 int32
	_ = v1200
	var v1207 int32
	_ = v1207
	var v1211 int32
	_ = v1211
	var v1213 int32
	_ = v1213
	var v1214 int32
	_ = v1214
	var v1215 int32
	_ = v1215
	var v1217 int32
	_ = v1217
	var v1225 int32
	_ = v1225
	var v1229 int32
	_ = v1229
	var v1230 int32
	_ = v1230
	var v1238 int32
	_ = v1238
	var v1242 int32
	_ = v1242
	var v1245 int32
	_ = v1245
	var v1246 int32
	_ = v1246
	var v1247 int64
	_ = v1247
	var v1248 int32
	_ = v1248
	var v1252 int32
	_ = v1252
	var v1259 int64
	_ = v1259
	var v1260 int32
	_ = v1260
	var v1261 int64
	_ = v1261
	var v1271 int32
	_ = v1271
	var v1273 int32
	_ = v1273
	var v1280 int32
	_ = v1280
	var v1284 int32
	_ = v1284
	var v1286 int32
	_ = v1286
	var v1287 int32
	_ = v1287
	var v1288 int32
	_ = v1288
	var v1290 int32
	_ = v1290
	var v1296 int32
	_ = v1296
	var v1300 int32
	_ = v1300
	var v1301 int32
	_ = v1301
	var v1302 int64
	_ = v1302
	var v1303 int32
	_ = v1303
	var v1307 int32
	_ = v1307
	var v1314 int64
	_ = v1314
	var v1315 int32
	_ = v1315
	var v1316 int64
	_ = v1316
	var v1326 int32
	_ = v1326
	var v1328 int32
	_ = v1328
	var v1335 int32
	_ = v1335
	var v1339 int32
	_ = v1339
	var v1341 int32
	_ = v1341
	var v1342 int32
	_ = v1342
	var v1343 int32
	_ = v1343
	var v1345 int32
	_ = v1345
	var v1353 int32
	_ = v1353
	var v1357 int32
	_ = v1357
	var v1358 int32
	_ = v1358
	var v1366 int32
	_ = v1366
	var v1370 int32
	_ = v1370
	var v1373 int32
	_ = v1373
	var v1374 int32
	_ = v1374
	var v1375 int64
	_ = v1375
	var v1376 int32
	_ = v1376
	var v1380 int32
	_ = v1380
	var v1387 int64
	_ = v1387
	var v1388 int32
	_ = v1388
	var v1389 int64
	_ = v1389
	var v1399 int32
	_ = v1399
	var v1401 int32
	_ = v1401
	var v1408 int32
	_ = v1408
	var v1412 int32
	_ = v1412
	var v1414 int32
	_ = v1414
	var v1415 int32
	_ = v1415
	var v1416 int32
	_ = v1416
	var v1418 int32
	_ = v1418
	var v1426 int32
	_ = v1426
	var v1430 int32
	_ = v1430
	var v1435 int32
	_ = v1435
	var v1441 int32
	_ = v1441
	var v1444 int32
	_ = v1444
	var v1445 int32
	_ = v1445
	var v1446 int32
	_ = v1446
	var v1451 int32
	_ = v1451
	var v1455 int32
	_ = v1455
	var v1456 int32
	_ = v1456
	var v1457 int32
	_ = v1457
	var v1459 int32
	_ = v1459
	var v1460 int32
	_ = v1460
	var v1461 int32
	_ = v1461
	var v1469 int32
	_ = v1469
	var v1473 int32
	_ = v1473
	var v1481 int32
	_ = v1481
	var v1488 int32
	_ = v1488
	var v1489 int32
	_ = v1489
	var v1492 int32
	_ = v1492
	var v1493 int32
	_ = v1493
	var v1494 int64
	_ = v1494
	var v1495 int32
	_ = v1495
	var v1499 int32
	_ = v1499
	var v1506 int64
	_ = v1506
	var v1507 int32
	_ = v1507
	var v1508 int64
	_ = v1508
	var v1518 int32
	_ = v1518
	var v1520 int32
	_ = v1520
	var v1527 int32
	_ = v1527
	var v1531 int32
	_ = v1531
	var v1533 int32
	_ = v1533
	var v1534 int32
	_ = v1534
	var v1535 int32
	_ = v1535
	var v1537 int32
	_ = v1537
	var v1545 int32
	_ = v1545
	var v1546 int32
	_ = v1546
	var v1550 int32
	_ = v1550
	var v1585 int32
	_ = v1585
	var v1590 int32
	_ = v1590
	var v1592 int32
	_ = v1592
	var v1594 int32
	_ = v1594
	var v1595 int32
	_ = v1595
	var v1596 int32
	_ = v1596
	var v1599 int32
	_ = v1599
	var v1600 int32
	_ = v1600
	var v1601 int64
	_ = v1601
	var v1602 int32
	_ = v1602
	var v1606 int32
	_ = v1606
	var v1613 int64
	_ = v1613
	var v1614 int32
	_ = v1614
	var v1615 int64
	_ = v1615
	var v1625 int32
	_ = v1625
	var v1627 int32
	_ = v1627
	var v1634 int32
	_ = v1634
	var v1638 int32
	_ = v1638
	var v1640 int32
	_ = v1640
	var v1641 int32
	_ = v1641
	var v1642 int32
	_ = v1642
	var v1644 int32
	_ = v1644
	var v1652 int32
	_ = v1652
	var v1663 int32
	_ = v1663
	var v1667 int32
	_ = v1667
	var v1668 int32
	_ = v1668
	var v1671 int32
	_ = v1671
	var v1672 int32
	_ = v1672
	var v1673 int64
	_ = v1673
	var v1674 int32
	_ = v1674
	var v1678 int32
	_ = v1678
	var v1685 int64
	_ = v1685
	var v1686 int32
	_ = v1686
	var v1687 int64
	_ = v1687
	var v1697 int32
	_ = v1697
	var v1699 int32
	_ = v1699
	var v1706 int32
	_ = v1706
	var v1710 int32
	_ = v1710
	var v1712 int32
	_ = v1712
	var v1713 int32
	_ = v1713
	var v1714 int32
	_ = v1714
	var v1716 int32
	_ = v1716
	var v1724 int32
	_ = v1724
	var v1728 int32
	_ = v1728
	var v1729 int32
	_ = v1729
	var v1738 int32
	_ = v1738
	var v1742 int32
	_ = v1742
	var v1746 int32
	_ = v1746
	var v1750 int32
	_ = v1750
	var v1789 int32
	_ = v1789
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v23 = v16 + v17*int32(132) + l0*int32(44)
	v24 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v25 = int32(3)
	v26 = v24 << (uint(v25) % 32)
	v32 = ((v17+v26)*v25 + l0) * int32(11)
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v35 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v35 < int32(1) {
		goto L3
	} else {
		goto L4
	}
L1:
	;
	v99 = *(*int32)(unsafe.Add(mBase, uint32(v23)))
	v100 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v99) {
		goto L15
	} else {
		goto L16
	}
L2:
	;
	v82 = int32(-1)
	v83 = v80 + v82
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v83
	v94 = v32 | int32(base.Ui32(v33^v82)>>(uint(int32(16))%32))&int32(_a_F_VP8RecordCoeffTokens_0)
	*(*uint16)(unsafe.Add(mBase, uint32(v81+v83<<(uint(int32(1))%32)))) = uint16(v94)
	goto L1
L3:
	;
	v39 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v39 != 0 {
		goto L6
	} else {
		goto L7
	}
L4:
	;
	v38 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v80 = v35
	v81 = v38
	goto L2
L5:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v66))) = int32(0)
	v73 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v73))) = v66
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v66
	v77 = v66 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v77
	v79 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v80 = v79
	v81 = v77
	goto L2
L6:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L1
L7:
	;
	v40 = int64(1)
	v41 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v45 = v41<<(uint(int32(1))%32) + int32(4)
	goto L11
L8:
	;
	if v66 != 0 {
		goto L5
	} else {
		goto L14
	}
L9:
	;
	goto L8
L10:
	;
	v64 = F_malloc(m, base.I32_wrap_i64(v40)*v45)
	mBase = m.M
	v66 = v64
	goto L9
L11:
	;
	v52 = base.I64_div_u_s(int64(2147418112), v40)
	v53 = int32(0)
	v54 = base.I64_extend_i32_u(v45)
	if base.Ui64(int64(4294967295)) < base.Ui64(v54*v40) {
		v66 = v53
		goto L9
	} else {
		goto L12
	}
L12:
	;
	if base.Ui64(v52) < base.Ui64(v54) {
		v66 = v53
		goto L9
	} else {
		goto L13
	}
L13:
	;
	goto L10
L14:
	;
	goto L6
L15:
	;
	v108 = int32(base.Ui32(v99+v100)>>(uint(v100)%32)) & int32(2147450879)
	goto L17
L16:
	;
	v108 = v99
	goto L17
L17:
	;
	v112 = base.B2i32(v33 < int32(0))
	if v33 < int32(0) {
		goto L18
	} else {
		goto L19
	}
L18:
	;
	v113 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L20
L19:
	;
	v113 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L20
L20:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v23))) = v108 + v113
	if v33 < int32(0) {
		v1789 = int32(0)
		goto L21
	} else {
		goto L22
	}
L21:
	;
	return v1789
L22:
	;
	if int32(15) < v17 {
		v1789 = int32(1)
		goto L21
	} else {
		goto L23
	}
L23:
	;
	v119 = v32
	v122 = v17
	v123 = v23
	goto L24
L24:
	;
	v134 = int32(1)
	v137 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v34+v122<<(uint(v134)%32)))))
	v138 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v138 < v134 {
		goto L28
	} else {
		goto L29
	}
L25:
	;
	v1789 = int32(1)
	goto L21
L26:
	;
	v203 = int32(1)
	v204 = v122 + v203
	v205 = *(*int32)(unsafe.Add(mBase, uint32(v123)+4))
	if base.Ui32(int32(-131073)) < base.Ui32(v205) {
		goto L40
	} else {
		goto L41
	}
L27:
	;
	v186 = v183 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v186
	v188 = int32(1)
	v199 = base.B2i32(v137&int32(_a_F_VP8RecordCoeffTokens_3) != int32(0))<<(uint(int32(15))%32) | (v119 + v188)
	*(*uint16)(unsafe.Add(mBase, uint32(v184+v186<<(uint(v188)%32)))) = uint16(v199)
	goto L26
L28:
	;
	v142 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v142 != 0 {
		goto L31
	} else {
		goto L32
	}
L29:
	;
	v141 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v183 = v138
	v184 = v141
	goto L27
L30:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v169))) = int32(0)
	v176 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v176))) = v169
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v169
	v180 = v169 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v180
	v182 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v183 = v182
	v184 = v180
	goto L27
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L26
L32:
	;
	v143 = int64(1)
	v144 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v148 = v144<<(uint(int32(1))%32) + int32(4)
	goto L36
L33:
	;
	if v169 != 0 {
		goto L30
	} else {
		goto L39
	}
L34:
	;
	goto L33
L35:
	;
	v167 = F_malloc(m, base.I32_wrap_i64(v143)*v148)
	mBase = m.M
	v169 = v167
	goto L34
L36:
	;
	v155 = base.I64_div_u_s(int64(2147418112), v143)
	v156 = int32(0)
	v157 = base.I64_extend_i32_u(v148)
	if base.Ui64(int64(4294967295)) < base.Ui64(v157*v143) {
		v169 = v156
		goto L34
	} else {
		goto L37
	}
L37:
	;
	if base.Ui64(v155) < base.Ui64(v157) {
		v169 = v156
		goto L34
	} else {
		goto L38
	}
L38:
	;
	goto L35
L39:
	;
	goto L31
L40:
	;
	v214 = int32(base.Ui32(v205+v203)>>(uint(v203)%32)) & int32(2147450879)
	goto L42
L41:
	;
	v214 = v205
	goto L42
L42:
	;
	v218 = v137 & int32(_a_F_VP8RecordCoeffTokens_3)
	if v218 != 0 {
		goto L43
	} else {
		goto L44
	}
L43:
	;
	v219 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L45
L44:
	;
	v219 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L45
L45:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+4)) = v214 + v219
	if v218 != 0 {
		goto L48
	} else {
		goto L49
	}
L46:
	;
	goto L25
L47:
	;
	if v204 != int32(16) {
		v119 = v1746
		v122 = v204
		v123 = v1750
		goto L24
	} else {
		goto L415
	}
L48:
	;
	v234 = base.I32_extend16_s(v137) >> (uint(int32(15)) % 32)
	v236 = v137 ^ v234 - v234
	v238 = v236 & int32(_a_F_VP8RecordCoeffTokens_3)
	v239 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v239 < int32(1) {
		goto L52
	} else {
		goto L53
	}
L49:
	;
	v222 = m.G23
	v223 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v222+v204))))
	v1746 = (v26 + v225) * int32(33)
	v1750 = v223 + v225*int32(132)
	goto L47
L50:
	;
	v302 = int32(1)
	v303 = *(*int32)(unsafe.Add(mBase, uint32(v123)+8))
	if base.Ui32(int32(-131073)) < base.Ui32(v303) {
		goto L64
	} else {
		goto L65
	}
L51:
	;
	v287 = v284 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v287
	v289 = int32(1)
	v298 = base.B2i32(base.Ui32(v289) < base.Ui32(v238))<<(uint(int32(15))%32) | (v119 + int32(2))
	*(*uint16)(unsafe.Add(mBase, uint32(v285+v287<<(uint(v289)%32)))) = uint16(v298)
	goto L50
L52:
	;
	v243 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v243 != 0 {
		goto L55
	} else {
		goto L56
	}
L53:
	;
	v242 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v284 = v239
	v285 = v242
	goto L51
L54:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v270))) = int32(0)
	v277 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v277))) = v270
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v270
	v281 = v270 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v281
	v283 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v284 = v283
	v285 = v281
	goto L51
L55:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L50
L56:
	;
	v244 = int64(1)
	v245 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v249 = v245<<(uint(int32(1))%32) + int32(4)
	goto L60
L57:
	;
	if v270 != 0 {
		goto L54
	} else {
		goto L63
	}
L58:
	;
	goto L57
L59:
	;
	v268 = F_malloc(m, base.I32_wrap_i64(v244)*v249)
	mBase = m.M
	v270 = v268
	goto L58
L60:
	;
	v256 = base.I64_div_u_s(int64(2147418112), v244)
	v257 = int32(0)
	v258 = base.I64_extend_i32_u(v249)
	if base.Ui64(int64(4294967295)) < base.Ui64(v258*v244) {
		v270 = v257
		goto L58
	} else {
		goto L61
	}
L61:
	;
	if base.Ui64(v256) < base.Ui64(v258) {
		v270 = v257
		goto L58
	} else {
		goto L62
	}
L62:
	;
	goto L59
L63:
	;
	goto L55
L64:
	;
	v312 = int32(base.Ui32(v303+v302)>>(uint(v302)%32)) & int32(2147450879)
	goto L66
L65:
	;
	v312 = v303
	goto L66
L66:
	;
	v316 = base.B2i32(base.Ui32(v238) < base.Ui32(int32(2)))
	if base.Ui32(v238) < base.Ui32(int32(2)) {
		goto L67
	} else {
		goto L68
	}
L67:
	;
	v317 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L69
L68:
	;
	v317 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L69
L69:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+8)) = v312 + v317
	if v316 == int32(0) {
		goto L71
	} else {
		goto L72
	}
L70:
	;
	v1592 = m.G23
	v1594 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1592+v204))))
	v1595 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v1596 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1596 < int32(1) {
		goto L381
	} else {
		goto L382
	}
L71:
	;
	v323 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v323 < int32(1) {
		goto L75
	} else {
		goto L76
	}
L72:
	;
	v1585 = int32(11)
	v1590 = v302
	goto L70
L73:
	;
	v386 = *(*int32)(unsafe.Add(mBase, uint32(v123)+12))
	v387 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v386) {
		goto L87
	} else {
		goto L88
	}
L74:
	;
	v371 = v368 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v371
	v382 = base.B2i32(base.Ui32(int32(4)) < base.Ui32(v238))<<(uint(int32(15))%32) | (v119 + int32(3))
	*(*uint16)(unsafe.Add(mBase, uint32(v369+v371<<(uint(int32(1))%32)))) = uint16(v382)
	goto L73
L75:
	;
	v327 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v327 != 0 {
		goto L78
	} else {
		goto L79
	}
L76:
	;
	v326 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v368 = v323
	v369 = v326
	goto L74
L77:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v354))) = int32(0)
	v361 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v361))) = v354
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v354
	v365 = v354 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v365
	v367 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v368 = v367
	v369 = v365
	goto L74
L78:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L73
L79:
	;
	v328 = int64(1)
	v329 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v333 = v329<<(uint(int32(1))%32) + int32(4)
	goto L83
L80:
	;
	if v354 != 0 {
		goto L77
	} else {
		goto L86
	}
L81:
	;
	goto L80
L82:
	;
	v352 = F_malloc(m, base.I32_wrap_i64(v328)*v333)
	mBase = m.M
	v354 = v352
	goto L81
L83:
	;
	v340 = base.I64_div_u_s(int64(2147418112), v328)
	v341 = int32(0)
	v342 = base.I64_extend_i32_u(v333)
	if base.Ui64(int64(4294967295)) < base.Ui64(v342*v328) {
		v354 = v341
		goto L81
	} else {
		goto L84
	}
L84:
	;
	if base.Ui64(v340) < base.Ui64(v342) {
		v354 = v341
		goto L81
	} else {
		goto L85
	}
L85:
	;
	goto L82
L86:
	;
	goto L78
L87:
	;
	v395 = int32(base.Ui32(v386+v387)>>(uint(v387)%32)) & int32(2147450879)
	goto L89
L88:
	;
	v395 = v386
	goto L89
L89:
	;
	if base.Ui32(v238) < base.Ui32(int32(5)) {
		goto L90
	} else {
		goto L91
	}
L90:
	;
	v400 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L92
L91:
	;
	v400 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L92
L92:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+12)) = v395 + v400
	v403 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if base.Ui32(int32(4)) < base.Ui32(v238) {
		goto L93
	} else {
		goto L94
	}
L93:
	;
	if v403 < int32(1) {
		goto L138
	} else {
		goto L139
	}
L94:
	;
	if v403 < int32(1) {
		goto L97
	} else {
		goto L98
	}
L95:
	;
	v468 = int32(2)
	v469 = *(*int32)(unsafe.Add(mBase, uint32(v123)+16))
	v470 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v469) {
		goto L109
	} else {
		goto L110
	}
L96:
	;
	v453 = v450 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v453
	v464 = base.B2i32(v238 != int32(2))<<(uint(int32(15))%32) | (v119 + int32(4))
	*(*uint16)(unsafe.Add(mBase, uint32(v451+v453<<(uint(int32(1))%32)))) = uint16(v464)
	goto L95
L97:
	;
	v409 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v409 != 0 {
		goto L100
	} else {
		goto L101
	}
L98:
	;
	v408 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v450 = v403
	v451 = v408
	goto L96
L99:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v436))) = int32(0)
	v443 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v443))) = v436
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v436
	v447 = v436 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v447
	v449 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v450 = v449
	v451 = v447
	goto L96
L100:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L95
L101:
	;
	v410 = int64(1)
	v411 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v415 = v411<<(uint(int32(1))%32) + int32(4)
	goto L105
L102:
	;
	if v436 != 0 {
		goto L99
	} else {
		goto L108
	}
L103:
	;
	goto L102
L104:
	;
	v434 = F_malloc(m, base.I32_wrap_i64(v410)*v415)
	mBase = m.M
	v436 = v434
	goto L103
L105:
	;
	v422 = base.I64_div_u_s(int64(2147418112), v410)
	v423 = int32(0)
	v424 = base.I64_extend_i32_u(v415)
	if base.Ui64(int64(4294967295)) < base.Ui64(v424*v410) {
		v436 = v423
		goto L103
	} else {
		goto L106
	}
L106:
	;
	if base.Ui64(v422) < base.Ui64(v424) {
		v436 = v423
		goto L103
	} else {
		goto L107
	}
L107:
	;
	goto L104
L108:
	;
	goto L100
L109:
	;
	v478 = int32(base.Ui32(v469+v470)>>(uint(v470)%32)) & int32(2147450879)
	goto L111
L110:
	;
	v478 = v469
	goto L111
L111:
	;
	v482 = base.B2i32(v238 == int32(2))
	if v238 == int32(2) {
		goto L112
	} else {
		goto L113
	}
L112:
	;
	v483 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L114
L113:
	;
	v483 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L114
L114:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+16)) = v478 + v483
	v486 = int32(22)
	if v238 == int32(2) {
		v1585 = v486
		v1590 = v468
		goto L70
	} else {
		goto L115
	}
L115:
	;
	v487 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v487 < int32(1) {
		goto L118
	} else {
		goto L119
	}
L116:
	;
	v550 = *(*int32)(unsafe.Add(mBase, uint32(v123)+20))
	v551 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v550) {
		goto L130
	} else {
		goto L131
	}
L117:
	;
	v535 = v532 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v535
	v546 = base.B2i32(v238 == int32(4))<<(uint(int32(15))%32) | (v119 + int32(5))
	*(*uint16)(unsafe.Add(mBase, uint32(v533+v535<<(uint(int32(1))%32)))) = uint16(v546)
	goto L116
L118:
	;
	v491 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v491 != 0 {
		goto L121
	} else {
		goto L122
	}
L119:
	;
	v490 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v532 = v487
	v533 = v490
	goto L117
L120:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v518))) = int32(0)
	v525 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v525))) = v518
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v518
	v529 = v518 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v529
	v531 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v532 = v531
	v533 = v529
	goto L117
L121:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L116
L122:
	;
	v492 = int64(1)
	v493 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v497 = v493<<(uint(int32(1))%32) + int32(4)
	goto L126
L123:
	;
	if v518 != 0 {
		goto L120
	} else {
		goto L129
	}
L124:
	;
	goto L123
L125:
	;
	v516 = F_malloc(m, base.I32_wrap_i64(v492)*v497)
	mBase = m.M
	v518 = v516
	goto L124
L126:
	;
	v504 = base.I64_div_u_s(int64(2147418112), v492)
	v505 = int32(0)
	v506 = base.I64_extend_i32_u(v497)
	if base.Ui64(int64(4294967295)) < base.Ui64(v506*v492) {
		v518 = v505
		goto L124
	} else {
		goto L127
	}
L127:
	;
	if base.Ui64(v504) < base.Ui64(v506) {
		v518 = v505
		goto L124
	} else {
		goto L128
	}
L128:
	;
	goto L125
L129:
	;
	goto L121
L130:
	;
	v559 = int32(base.Ui32(v550+v551)>>(uint(v551)%32)) & int32(2147450879)
	goto L132
L131:
	;
	v559 = v550
	goto L132
L132:
	;
	if v238 == int32(4) {
		goto L133
	} else {
		goto L134
	}
L133:
	;
	v564 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L135
L134:
	;
	v564 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L135
L135:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+20)) = v559 + v564
	v1585 = v486
	v1590 = v468
	goto L70
L136:
	;
	v629 = *(*int32)(unsafe.Add(mBase, uint32(v123)+24))
	v630 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v629) {
		goto L150
	} else {
		goto L151
	}
L137:
	;
	v614 = v611 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v614
	v625 = base.B2i32(base.Ui32(int32(10)) < base.Ui32(v238))<<(uint(int32(15))%32) | (v119 + int32(6))
	*(*uint16)(unsafe.Add(mBase, uint32(v612+v614<<(uint(int32(1))%32)))) = uint16(v625)
	goto L136
L138:
	;
	v570 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v570 != 0 {
		goto L141
	} else {
		goto L142
	}
L139:
	;
	v569 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v611 = v403
	v612 = v569
	goto L137
L140:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v597))) = int32(0)
	v604 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v604))) = v597
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v597
	v608 = v597 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v608
	v610 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v611 = v610
	v612 = v608
	goto L137
L141:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L136
L142:
	;
	v571 = int64(1)
	v572 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v576 = v572<<(uint(int32(1))%32) + int32(4)
	goto L146
L143:
	;
	if v597 != 0 {
		goto L140
	} else {
		goto L149
	}
L144:
	;
	goto L143
L145:
	;
	v595 = F_malloc(m, base.I32_wrap_i64(v571)*v576)
	mBase = m.M
	v597 = v595
	goto L144
L146:
	;
	v583 = base.I64_div_u_s(int64(2147418112), v571)
	v584 = int32(0)
	v585 = base.I64_extend_i32_u(v576)
	if base.Ui64(int64(4294967295)) < base.Ui64(v585*v571) {
		v597 = v584
		goto L144
	} else {
		goto L147
	}
L147:
	;
	if base.Ui64(v583) < base.Ui64(v585) {
		v597 = v584
		goto L144
	} else {
		goto L148
	}
L148:
	;
	goto L145
L149:
	;
	goto L141
L150:
	;
	v638 = int32(base.Ui32(v629+v630)>>(uint(v630)%32)) & int32(2147450879)
	goto L152
L151:
	;
	v638 = v629
	goto L152
L152:
	;
	if base.Ui32(v238) < base.Ui32(int32(11)) {
		goto L153
	} else {
		goto L154
	}
L153:
	;
	v643 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L155
L154:
	;
	v643 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L155
L155:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+24)) = v638 + v643
	if base.Ui32(int32(10)) < base.Ui32(v238) {
		goto L158
	} else {
		goto L159
	}
L156:
	;
	v1585 = int32(22)
	v1590 = int32(2)
	goto L70
L157:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L156
L158:
	;
	v906 = v238 + int32(-3)
	if base.Ui32(int32(15)) < base.Ui32(v906) {
		goto L226
	} else {
		goto L227
	}
L159:
	;
	v648 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v648 < int32(1) {
		goto L162
	} else {
		goto L163
	}
L160:
	;
	v711 = *(*int32)(unsafe.Add(mBase, uint32(v123)+28))
	v712 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v711) {
		goto L174
	} else {
		goto L175
	}
L161:
	;
	v696 = v693 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v696
	v707 = base.B2i32(base.Ui32(int32(6)) < base.Ui32(v238))<<(uint(int32(15))%32) | (v119 + int32(7))
	*(*uint16)(unsafe.Add(mBase, uint32(v694+v696<<(uint(int32(1))%32)))) = uint16(v707)
	goto L160
L162:
	;
	v652 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v652 != 0 {
		goto L165
	} else {
		goto L166
	}
L163:
	;
	v651 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v693 = v648
	v694 = v651
	goto L161
L164:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v679))) = int32(0)
	v686 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v686))) = v679
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v679
	v690 = v679 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v690
	v692 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v693 = v692
	v694 = v690
	goto L161
L165:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L160
L166:
	;
	v653 = int64(1)
	v654 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v658 = v654<<(uint(int32(1))%32) + int32(4)
	goto L170
L167:
	;
	if v679 != 0 {
		goto L164
	} else {
		goto L173
	}
L168:
	;
	goto L167
L169:
	;
	v677 = F_malloc(m, base.I32_wrap_i64(v653)*v658)
	mBase = m.M
	v679 = v677
	goto L168
L170:
	;
	v665 = base.I64_div_u_s(int64(2147418112), v653)
	v666 = int32(0)
	v667 = base.I64_extend_i32_u(v658)
	if base.Ui64(int64(4294967295)) < base.Ui64(v667*v653) {
		v679 = v666
		goto L168
	} else {
		goto L171
	}
L171:
	;
	if base.Ui64(v665) < base.Ui64(v667) {
		v679 = v666
		goto L168
	} else {
		goto L172
	}
L172:
	;
	goto L169
L173:
	;
	goto L165
L174:
	;
	v720 = int32(base.Ui32(v711+v712)>>(uint(v712)%32)) & int32(2147450879)
	goto L176
L175:
	;
	v720 = v711
	goto L176
L176:
	;
	if base.Ui32(v238) < base.Ui32(int32(7)) {
		goto L177
	} else {
		goto L178
	}
L177:
	;
	v725 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L179
L178:
	;
	v725 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L179
L179:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+28)) = v720 + v725
	v728 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if base.Ui32(int32(6)) < base.Ui32(v238) {
		goto L180
	} else {
		goto L181
	}
L180:
	;
	if v728 < int32(1) {
		goto L198
	} else {
		goto L199
	}
L181:
	;
	if v728 < int32(1) {
		goto L183
	} else {
		goto L184
	}
L182:
	;
	v777 = v774 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v777
	if v238 == int32(6) {
		goto L193
	} else {
		goto L194
	}
L183:
	;
	v734 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v734 != 0 {
		goto L157
	} else {
		goto L185
	}
L184:
	;
	v733 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v774 = v728
	v775 = v733
	goto L182
L185:
	;
	v735 = int64(1)
	v736 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v740 = v736<<(uint(int32(1))%32) + int32(4)
	goto L189
L186:
	;
	if v761 == int32(0) {
		goto L157
	} else {
		goto L192
	}
L187:
	;
	goto L186
L188:
	;
	v759 = F_malloc(m, base.I32_wrap_i64(v735)*v740)
	mBase = m.M
	v761 = v759
	goto L187
L189:
	;
	v747 = base.I64_div_u_s(int64(2147418112), v735)
	v748 = int32(0)
	v749 = base.I64_extend_i32_u(v740)
	if base.Ui64(int64(4294967295)) < base.Ui64(v749*v735) {
		v761 = v748
		goto L187
	} else {
		goto L190
	}
L190:
	;
	if base.Ui64(v747) < base.Ui64(v749) {
		v761 = v748
		goto L187
	} else {
		goto L191
	}
L191:
	;
	goto L188
L192:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v761))) = int32(0)
	v767 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v767))) = v761
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v761
	v771 = v761 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v771
	v773 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v774 = v773
	v775 = v771
	goto L182
L193:
	;
	v786 = int32(-16225)
	goto L195
L194:
	;
	v786 = int32(_a_F_VP8RecordCoeffTokens_4)
	goto L195
L195:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v775+v777<<(uint(int32(1))%32)))) = uint16(v786)
	goto L156
L196:
	;
	if v847 < int32(1) {
		goto L214
	} else {
		goto L215
	}
L197:
	;
	v836 = v833 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v836
	if base.Ui32(int32(8)) < base.Ui32(v238) {
		goto L210
	} else {
		goto L211
	}
L198:
	;
	v791 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v791 != 0 {
		v821 = v728
		goto L201
	} else {
		goto L202
	}
L199:
	;
	v790 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v833 = v728
	v834 = v790
	goto L197
L200:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v818))) = int32(0)
	v826 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v826))) = v818
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v818
	v830 = v818 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v830
	v832 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v833 = v832
	v834 = v830
	goto L197
L201:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	v847 = v821
	goto L196
L202:
	;
	v792 = int64(1)
	v793 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v797 = v793<<(uint(int32(1))%32) + int32(4)
	goto L206
L203:
	;
	if v818 != 0 {
		goto L200
	} else {
		goto L209
	}
L204:
	;
	goto L203
L205:
	;
	v816 = F_malloc(m, base.I32_wrap_i64(v792)*v797)
	mBase = m.M
	v818 = v816
	goto L204
L206:
	;
	v804 = base.I64_div_u_s(int64(2147418112), v792)
	v805 = int32(0)
	v806 = base.I64_extend_i32_u(v797)
	if base.Ui64(int64(4294967295)) < base.Ui64(v806*v792) {
		v818 = v805
		goto L204
	} else {
		goto L207
	}
L207:
	;
	if base.Ui64(v804) < base.Ui64(v806) {
		v818 = v805
		goto L204
	} else {
		goto L208
	}
L208:
	;
	goto L205
L209:
	;
	v820 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	v821 = v820
	goto L201
L210:
	;
	v845 = int32(-16219)
	goto L212
L211:
	;
	v845 = int32(_a_F_VP8RecordCoeffTokens_5)
	goto L212
L212:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v834+v836<<(uint(int32(1))%32)))) = uint16(v845)
	v847 = v836
	goto L196
L213:
	;
	v895 = v892 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v895
	v903 = v236<<(uint(int32(15))%32) ^ int32(_a_F_VP8RecordCoeffTokens_6)
	*(*uint16)(unsafe.Add(mBase, uint32(v893+v895<<(uint(int32(1))%32)))) = uint16(v903)
	goto L156
L214:
	;
	v852 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v852 != 0 {
		goto L157
	} else {
		goto L216
	}
L215:
	;
	v851 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v892 = v847
	v893 = v851
	goto L213
L216:
	;
	v853 = int64(1)
	v854 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v858 = v854<<(uint(int32(1))%32) + int32(4)
	goto L220
L217:
	;
	if v879 == int32(0) {
		goto L157
	} else {
		goto L223
	}
L218:
	;
	goto L217
L219:
	;
	v877 = F_malloc(m, base.I32_wrap_i64(v853)*v858)
	mBase = m.M
	v879 = v877
	goto L218
L220:
	;
	v865 = base.I64_div_u_s(int64(2147418112), v853)
	v866 = int32(0)
	v867 = base.I64_extend_i32_u(v858)
	if base.Ui64(int64(4294967295)) < base.Ui64(v867*v853) {
		v879 = v866
		goto L218
	} else {
		goto L221
	}
L221:
	;
	if base.Ui64(v865) < base.Ui64(v867) {
		v879 = v866
		goto L218
	} else {
		goto L222
	}
L222:
	;
	goto L219
L223:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v879))) = int32(0)
	v885 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v885))) = v879
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v879
	v889 = v879 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v889
	v891 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v892 = v891
	v893 = v889
	goto L213
L224:
	;
	v1460 = *(*int32)(unsafe.Add(mBase, uint32(v123)+36))
	v1461 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v1460) {
		goto L356
	} else {
		goto L357
	}
L225:
	;
	v1451 = m.G79
	v1455 = int32(4)
	v1456 = v1451
	v1457 = int32(-11)
	v1459 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L224
L226:
	;
	if base.Ui32(int32(31)) < base.Ui32(v906) {
		goto L260
	} else {
		goto L261
	}
L227:
	;
	v909 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v909 < int32(1) {
		goto L230
	} else {
		goto L231
	}
L228:
	;
	v967 = *(*int32)(unsafe.Add(mBase, uint32(v123)+32))
	v968 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v967) {
		goto L242
	} else {
		goto L243
	}
L229:
	;
	v957 = v954 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v957
	v963 = v119 + int32(8)
	*(*uint16)(unsafe.Add(mBase, uint32(v955+v957<<(uint(int32(1))%32)))) = uint16(v963)
	goto L228
L230:
	;
	v913 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v913 != 0 {
		goto L233
	} else {
		goto L234
	}
L231:
	;
	v912 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v954 = v909
	v955 = v912
	goto L229
L232:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v940))) = int32(0)
	v947 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v947))) = v940
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v940
	v951 = v940 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v951
	v953 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v954 = v953
	v955 = v951
	goto L229
L233:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L228
L234:
	;
	v914 = int64(1)
	v915 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v919 = v915<<(uint(int32(1))%32) + int32(4)
	goto L238
L235:
	;
	if v940 != 0 {
		goto L232
	} else {
		goto L241
	}
L236:
	;
	goto L235
L237:
	;
	v938 = F_malloc(m, base.I32_wrap_i64(v914)*v919)
	mBase = m.M
	v940 = v938
	goto L236
L238:
	;
	v926 = base.I64_div_u_s(int64(2147418112), v914)
	v927 = int32(0)
	v928 = base.I64_extend_i32_u(v919)
	if base.Ui64(int64(4294967295)) < base.Ui64(v928*v914) {
		v940 = v927
		goto L236
	} else {
		goto L239
	}
L239:
	;
	if base.Ui64(v926) < base.Ui64(v928) {
		v940 = v927
		goto L236
	} else {
		goto L240
	}
L240:
	;
	goto L237
L241:
	;
	goto L233
L242:
	;
	v976 = int32(base.Ui32(v967+v968)>>(uint(v968)%32)) & int32(2147450879)
	goto L244
L243:
	;
	v976 = v967
	goto L244
L244:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+32)) = v976 + int32(_a_F_VP8RecordCoeffTokens_1)
	v980 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v980 < int32(1) {
		goto L246
	} else {
		goto L247
	}
L245:
	;
	v1028 = v1025 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1028
	v1034 = v119 + int32(9)
	*(*uint16)(unsafe.Add(mBase, uint32(v1026+v1028<<(uint(int32(1))%32)))) = uint16(v1034)
	goto L225
L246:
	;
	v984 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v984 != 0 {
		goto L249
	} else {
		goto L250
	}
L247:
	;
	v983 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1025 = v980
	v1026 = v983
	goto L245
L248:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1011))) = int32(0)
	v1018 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1018))) = v1011
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1011
	v1022 = v1011 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1022
	v1024 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1025 = v1024
	v1026 = v1022
	goto L245
L249:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L225
L250:
	;
	v985 = int64(1)
	v986 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v990 = v986<<(uint(int32(1))%32) + int32(4)
	goto L254
L251:
	;
	if v1011 != 0 {
		goto L248
	} else {
		goto L257
	}
L252:
	;
	goto L251
L253:
	;
	v1009 = F_malloc(m, base.I32_wrap_i64(v985)*v990)
	mBase = m.M
	v1011 = v1009
	goto L252
L254:
	;
	v997 = base.I64_div_u_s(int64(2147418112), v985)
	v998 = int32(0)
	v999 = base.I64_extend_i32_u(v990)
	if base.Ui64(int64(4294967295)) < base.Ui64(v999*v985) {
		v1011 = v998
		goto L252
	} else {
		goto L255
	}
L255:
	;
	if base.Ui64(v997) < base.Ui64(v999) {
		v1011 = v998
		goto L252
	} else {
		goto L256
	}
L256:
	;
	goto L253
L257:
	;
	goto L249
L258:
	;
	v1455 = v1444
	v1456 = v1445
	v1457 = v1446
	v1459 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L224
L259:
	;
	v1441 = m.G78
	v1444 = int32(8)
	v1445 = v1441
	v1446 = int32(-19)
	goto L258
L260:
	;
	v1167 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if base.Ui32(int32(63)) < base.Ui32(v906) {
		goto L293
	} else {
		goto L294
	}
L261:
	;
	v1038 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1038 < int32(1) {
		goto L264
	} else {
		goto L265
	}
L262:
	;
	v1096 = *(*int32)(unsafe.Add(mBase, uint32(v123)+32))
	v1097 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v1096) {
		goto L276
	} else {
		goto L277
	}
L263:
	;
	v1086 = v1083 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1086
	v1092 = v119 + int32(8)
	*(*uint16)(unsafe.Add(mBase, uint32(v1084+v1086<<(uint(int32(1))%32)))) = uint16(v1092)
	goto L262
L264:
	;
	v1042 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1042 != 0 {
		goto L267
	} else {
		goto L268
	}
L265:
	;
	v1041 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1083 = v1038
	v1084 = v1041
	goto L263
L266:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1069))) = int32(0)
	v1076 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1076))) = v1069
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1069
	v1080 = v1069 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1080
	v1082 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1083 = v1082
	v1084 = v1080
	goto L263
L267:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L262
L268:
	;
	v1043 = int64(1)
	v1044 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1048 = v1044<<(uint(int32(1))%32) + int32(4)
	goto L272
L269:
	;
	if v1069 != 0 {
		goto L266
	} else {
		goto L275
	}
L270:
	;
	goto L269
L271:
	;
	v1067 = F_malloc(m, base.I32_wrap_i64(v1043)*v1048)
	mBase = m.M
	v1069 = v1067
	goto L270
L272:
	;
	v1055 = base.I64_div_u_s(int64(2147418112), v1043)
	v1056 = int32(0)
	v1057 = base.I64_extend_i32_u(v1048)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1057*v1043) {
		v1069 = v1056
		goto L270
	} else {
		goto L273
	}
L273:
	;
	if base.Ui64(v1055) < base.Ui64(v1057) {
		v1069 = v1056
		goto L270
	} else {
		goto L274
	}
L274:
	;
	goto L271
L275:
	;
	goto L267
L276:
	;
	v1105 = int32(base.Ui32(v1096+v1097)>>(uint(v1097)%32)) & int32(2147450879)
	goto L278
L277:
	;
	v1105 = v1096
	goto L278
L278:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+32)) = v1105 + int32(_a_F_VP8RecordCoeffTokens_1)
	v1109 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1109 < int32(1) {
		goto L280
	} else {
		goto L281
	}
L279:
	;
	v1157 = v1154 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1157
	v1165 = v119 + int32(9) | int32(_a_F_VP8RecordCoeffTokens_0)
	*(*uint16)(unsafe.Add(mBase, uint32(v1155+v1157<<(uint(int32(1))%32)))) = uint16(v1165)
	goto L259
L280:
	;
	v1113 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1113 != 0 {
		goto L283
	} else {
		goto L284
	}
L281:
	;
	v1112 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1154 = v1109
	v1155 = v1112
	goto L279
L282:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1140))) = int32(0)
	v1147 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1147))) = v1140
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1140
	v1151 = v1140 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1151
	v1153 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1154 = v1153
	v1155 = v1151
	goto L279
L283:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L259
L284:
	;
	v1114 = int64(1)
	v1115 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1119 = v1115<<(uint(int32(1))%32) + int32(4)
	goto L288
L285:
	;
	if v1140 != 0 {
		goto L282
	} else {
		goto L291
	}
L286:
	;
	goto L285
L287:
	;
	v1138 = F_malloc(m, base.I32_wrap_i64(v1114)*v1119)
	mBase = m.M
	v1140 = v1138
	goto L286
L288:
	;
	v1126 = base.I64_div_u_s(int64(2147418112), v1114)
	v1127 = int32(0)
	v1128 = base.I64_extend_i32_u(v1119)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1128*v1114) {
		v1140 = v1127
		goto L286
	} else {
		goto L289
	}
L289:
	;
	if base.Ui64(v1126) < base.Ui64(v1128) {
		v1140 = v1127
		goto L286
	} else {
		goto L290
	}
L290:
	;
	goto L287
L291:
	;
	goto L283
L292:
	;
	v1435 = m.G77
	v1455 = int32(16)
	v1456 = v1435
	v1457 = int32(-35)
	v1459 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L224
L293:
	;
	if v1167 < int32(1) {
		goto L327
	} else {
		goto L328
	}
L294:
	;
	if v1167 < int32(1) {
		goto L297
	} else {
		goto L298
	}
L295:
	;
	v1229 = *(*int32)(unsafe.Add(mBase, uint32(v123)+32))
	v1230 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v1229) {
		goto L309
	} else {
		goto L310
	}
L296:
	;
	v1217 = v1215 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1217
	v1225 = v119 + int32(8) | int32(_a_F_VP8RecordCoeffTokens_0)
	*(*uint16)(unsafe.Add(mBase, uint32(v1214+v1217<<(uint(int32(1))%32)))) = uint16(v1225)
	goto L295
L297:
	;
	v1173 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1173 != 0 {
		goto L300
	} else {
		goto L301
	}
L298:
	;
	v1172 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1214 = v1172
	v1215 = v1167
	goto L296
L299:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1200))) = int32(0)
	v1207 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1207))) = v1200
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1200
	v1211 = v1200 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1211
	v1213 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1214 = v1211
	v1215 = v1213
	goto L296
L300:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L295
L301:
	;
	v1174 = int64(1)
	v1175 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1179 = v1175<<(uint(int32(1))%32) + int32(4)
	goto L305
L302:
	;
	if v1200 != 0 {
		goto L299
	} else {
		goto L308
	}
L303:
	;
	goto L302
L304:
	;
	v1198 = F_malloc(m, base.I32_wrap_i64(v1174)*v1179)
	mBase = m.M
	v1200 = v1198
	goto L303
L305:
	;
	v1186 = base.I64_div_u_s(int64(2147418112), v1174)
	v1187 = int32(0)
	v1188 = base.I64_extend_i32_u(v1179)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1188*v1174) {
		v1200 = v1187
		goto L303
	} else {
		goto L306
	}
L306:
	;
	if base.Ui64(v1186) < base.Ui64(v1188) {
		v1200 = v1187
		goto L303
	} else {
		goto L307
	}
L307:
	;
	goto L304
L308:
	;
	goto L300
L309:
	;
	v1238 = int32(base.Ui32(v1229+v1230)>>(uint(v1230)%32)) & int32(2147450879)
	goto L311
L310:
	;
	v1238 = v1229
	goto L311
L311:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+32)) = v1238 + int32(_a_F_VP8RecordCoeffTokens_2)
	v1242 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1242 < int32(1) {
		goto L313
	} else {
		goto L314
	}
L312:
	;
	v1290 = v1287 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1290
	v1296 = v119 + int32(10)
	*(*uint16)(unsafe.Add(mBase, uint32(v1288+v1290<<(uint(int32(1))%32)))) = uint16(v1296)
	goto L292
L313:
	;
	v1246 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1246 != 0 {
		goto L316
	} else {
		goto L317
	}
L314:
	;
	v1245 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1287 = v1242
	v1288 = v1245
	goto L312
L315:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1273))) = int32(0)
	v1280 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1280))) = v1273
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1273
	v1284 = v1273 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1284
	v1286 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1287 = v1286
	v1288 = v1284
	goto L312
L316:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L292
L317:
	;
	v1247 = int64(1)
	v1248 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1252 = v1248<<(uint(int32(1))%32) + int32(4)
	goto L321
L318:
	;
	if v1273 != 0 {
		goto L315
	} else {
		goto L324
	}
L319:
	;
	goto L318
L320:
	;
	v1271 = F_malloc(m, base.I32_wrap_i64(v1247)*v1252)
	mBase = m.M
	v1273 = v1271
	goto L319
L321:
	;
	v1259 = base.I64_div_u_s(int64(2147418112), v1247)
	v1260 = int32(0)
	v1261 = base.I64_extend_i32_u(v1252)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1261*v1247) {
		v1273 = v1260
		goto L319
	} else {
		goto L322
	}
L322:
	;
	if base.Ui64(v1259) < base.Ui64(v1261) {
		v1273 = v1260
		goto L319
	} else {
		goto L323
	}
L323:
	;
	goto L320
L324:
	;
	goto L316
L325:
	;
	v1357 = *(*int32)(unsafe.Add(mBase, uint32(v123)+32))
	v1358 = int32(1)
	if base.Ui32(int32(-131073)) < base.Ui32(v1357) {
		goto L339
	} else {
		goto L340
	}
L326:
	;
	v1345 = v1343 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1345
	v1353 = v119 + int32(8) | int32(_a_F_VP8RecordCoeffTokens_0)
	*(*uint16)(unsafe.Add(mBase, uint32(v1342+v1345<<(uint(int32(1))%32)))) = uint16(v1353)
	goto L325
L327:
	;
	v1301 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1301 != 0 {
		goto L330
	} else {
		goto L331
	}
L328:
	;
	v1300 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1342 = v1300
	v1343 = v1167
	goto L326
L329:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1328))) = int32(0)
	v1335 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1335))) = v1328
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1328
	v1339 = v1328 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1339
	v1341 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1342 = v1339
	v1343 = v1341
	goto L326
L330:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L325
L331:
	;
	v1302 = int64(1)
	v1303 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1307 = v1303<<(uint(int32(1))%32) + int32(4)
	goto L335
L332:
	;
	if v1328 != 0 {
		goto L329
	} else {
		goto L338
	}
L333:
	;
	goto L332
L334:
	;
	v1326 = F_malloc(m, base.I32_wrap_i64(v1302)*v1307)
	mBase = m.M
	v1328 = v1326
	goto L333
L335:
	;
	v1314 = base.I64_div_u_s(int64(2147418112), v1302)
	v1315 = int32(0)
	v1316 = base.I64_extend_i32_u(v1307)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1316*v1302) {
		v1328 = v1315
		goto L333
	} else {
		goto L336
	}
L336:
	;
	if base.Ui64(v1314) < base.Ui64(v1316) {
		v1328 = v1315
		goto L333
	} else {
		goto L337
	}
L337:
	;
	goto L334
L338:
	;
	goto L330
L339:
	;
	v1366 = int32(base.Ui32(v1357+v1358)>>(uint(v1358)%32)) & int32(2147450879)
	goto L341
L340:
	;
	v1366 = v1357
	goto L341
L341:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+32)) = v1366 + int32(_a_F_VP8RecordCoeffTokens_2)
	v1370 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1370 < int32(1) {
		goto L344
	} else {
		goto L345
	}
L342:
	;
	v1430 = m.G76
	v1444 = int32(1024)
	v1445 = v1430
	v1446 = int32(-67)
	goto L258
L343:
	;
	v1418 = v1415 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1418
	v1426 = v119 + int32(10) | int32(_a_F_VP8RecordCoeffTokens_0)
	*(*uint16)(unsafe.Add(mBase, uint32(v1416+v1418<<(uint(int32(1))%32)))) = uint16(v1426)
	goto L342
L344:
	;
	v1374 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1374 != 0 {
		goto L347
	} else {
		goto L348
	}
L345:
	;
	v1373 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1415 = v1370
	v1416 = v1373
	goto L343
L346:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1401))) = int32(0)
	v1408 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1408))) = v1401
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1401
	v1412 = v1401 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1412
	v1414 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1415 = v1414
	v1416 = v1412
	goto L343
L347:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L342
L348:
	;
	v1375 = int64(1)
	v1376 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1380 = v1376<<(uint(int32(1))%32) + int32(4)
	goto L352
L349:
	;
	if v1401 != 0 {
		goto L346
	} else {
		goto L355
	}
L350:
	;
	goto L349
L351:
	;
	v1399 = F_malloc(m, base.I32_wrap_i64(v1375)*v1380)
	mBase = m.M
	v1401 = v1399
	goto L350
L352:
	;
	v1387 = base.I64_div_u_s(int64(2147418112), v1375)
	v1388 = int32(0)
	v1389 = base.I64_extend_i32_u(v1380)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1389*v1375) {
		v1401 = v1388
		goto L350
	} else {
		goto L353
	}
L353:
	;
	if base.Ui64(v1387) < base.Ui64(v1389) {
		v1401 = v1388
		goto L350
	} else {
		goto L354
	}
L354:
	;
	goto L351
L355:
	;
	goto L347
L356:
	;
	v1469 = int32(base.Ui32(v1460+v1461)>>(uint(v1461)%32)) & int32(2147450879)
	goto L358
L357:
	;
	v1469 = v1460
	goto L358
L358:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v123)+36)) = v1469 + v1459
	v1473 = v1455
	v1481 = v1456
	goto L359
L359:
	;
	v1488 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1481))))
	v1489 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1489 < int32(1) {
		goto L363
	} else {
		goto L364
	}
L361:
	;
	v1550 = int32(1)
	if base.Ui32(v1550) < base.Ui32(v1473) {
		v1473 = int32(base.Ui32(v1473) >> (uint(v1550) % 32))
		v1481 = v1481 + v1550
		goto L359
	} else {
		goto L378
	}
L362:
	;
	v1537 = v1534 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1537
	if v1473&(v1457+v238) != 0 {
		goto L375
	} else {
		goto L376
	}
L363:
	;
	v1493 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1493 != 0 {
		goto L366
	} else {
		goto L367
	}
L364:
	;
	v1492 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1534 = v1489
	v1535 = v1492
	goto L362
L365:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1520))) = int32(0)
	v1527 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1527))) = v1520
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1520
	v1531 = v1520 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1531
	v1533 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1534 = v1533
	v1535 = v1531
	goto L362
L366:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L361
L367:
	;
	v1494 = int64(1)
	v1495 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1499 = v1495<<(uint(int32(1))%32) + int32(4)
	goto L371
L368:
	;
	if v1520 != 0 {
		goto L365
	} else {
		goto L374
	}
L369:
	;
	goto L368
L370:
	;
	v1518 = F_malloc(m, base.I32_wrap_i64(v1494)*v1499)
	mBase = m.M
	v1520 = v1518
	goto L369
L371:
	;
	v1506 = base.I64_div_u_s(int64(2147418112), v1494)
	v1507 = int32(0)
	v1508 = base.I64_extend_i32_u(v1499)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1508*v1494) {
		v1520 = v1507
		goto L369
	} else {
		goto L372
	}
L372:
	;
	if base.Ui64(v1506) < base.Ui64(v1508) {
		v1520 = v1507
		goto L369
	} else {
		goto L373
	}
L373:
	;
	goto L370
L374:
	;
	goto L366
L375:
	;
	v1545 = int32(-16384)
	goto L377
L376:
	;
	v1545 = int32(_a_F_VP8RecordCoeffTokens_7)
	goto L377
L377:
	;
	v1546 = v1545 | v1488
	*(*uint16)(unsafe.Add(mBase, uint32(v1535+v1537<<(uint(int32(1))%32)))) = uint16(v1546)
	goto L361
L378:
	;
	goto L156
L379:
	;
	if v204 == int32(16) {
		goto L46
	} else {
		goto L393
	}
L380:
	;
	v1644 = v1642 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1644
	v1652 = v137&int32(-32768) | int32(_a_F_VP8RecordCoeffTokens_8)
	*(*uint16)(unsafe.Add(mBase, uint32(v1641+v1644<<(uint(int32(1))%32)))) = uint16(v1652)
	goto L379
L381:
	;
	v1600 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1600 != 0 {
		goto L384
	} else {
		goto L385
	}
L382:
	;
	v1599 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1641 = v1599
	v1642 = v1596
	goto L380
L383:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1627))) = int32(0)
	v1634 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1634))) = v1627
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1627
	v1638 = v1627 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1638
	v1640 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1641 = v1638
	v1642 = v1640
	goto L380
L384:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L379
L385:
	;
	v1601 = int64(1)
	v1602 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1606 = v1602<<(uint(int32(1))%32) + int32(4)
	goto L389
L386:
	;
	if v1627 != 0 {
		goto L383
	} else {
		goto L392
	}
L387:
	;
	goto L386
L388:
	;
	v1625 = F_malloc(m, base.I32_wrap_i64(v1601)*v1606)
	mBase = m.M
	v1627 = v1625
	goto L387
L389:
	;
	v1613 = base.I64_div_u_s(int64(2147418112), v1601)
	v1614 = int32(0)
	v1615 = base.I64_extend_i32_u(v1606)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1615*v1601) {
		v1627 = v1614
		goto L387
	} else {
		goto L390
	}
L390:
	;
	if base.Ui64(v1613) < base.Ui64(v1615) {
		v1627 = v1614
		goto L387
	} else {
		goto L391
	}
L391:
	;
	goto L388
L392:
	;
	goto L384
L393:
	;
	v1663 = v1595 + v1594*int32(132) + v1590*int32(44)
	v1667 = (v26+v1594)*int32(33) + v1585
	v1668 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	if v1668 < int32(1) {
		goto L396
	} else {
		goto L397
	}
L394:
	;
	v1728 = int32(1)
	v1729 = *(*int32)(unsafe.Add(mBase, uint32(v1663)))
	if base.Ui32(int32(-131073)) < base.Ui32(v1729) {
		goto L408
	} else {
		goto L409
	}
L395:
	;
	v1716 = v1713 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v1716
	v1724 = v1667 | base.B2i32(v122 < v33)<<(uint(int32(15))%32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1714+v1716<<(uint(int32(1))%32)))) = uint16(v1724)
	goto L394
L396:
	;
	v1672 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	if v1672 != 0 {
		goto L399
	} else {
		goto L400
	}
L397:
	;
	v1671 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v1713 = v1668
	v1714 = v1671
	goto L395
L398:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1699))) = int32(0)
	v1706 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1706))) = v1699
	*(*int32)(unsafe.Add(mBase, uint32(l2)+4)) = v1699
	v1710 = v1699 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v1710
	v1712 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1713 = v1712
	v1714 = v1710
	goto L395
L399:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = int32(1)
	goto L394
L400:
	;
	v1673 = int64(1)
	v1674 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v1678 = v1674<<(uint(int32(1))%32) + int32(4)
	goto L404
L401:
	;
	if v1699 != 0 {
		goto L398
	} else {
		goto L407
	}
L402:
	;
	goto L401
L403:
	;
	v1697 = F_malloc(m, base.I32_wrap_i64(v1673)*v1678)
	mBase = m.M
	v1699 = v1697
	goto L402
L404:
	;
	v1685 = base.I64_div_u_s(int64(2147418112), v1673)
	v1686 = int32(0)
	v1687 = base.I64_extend_i32_u(v1678)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1687*v1673) {
		v1699 = v1686
		goto L402
	} else {
		goto L405
	}
L405:
	;
	if base.Ui64(v1685) < base.Ui64(v1687) {
		v1699 = v1686
		goto L402
	} else {
		goto L406
	}
L406:
	;
	goto L403
L407:
	;
	goto L399
L408:
	;
	v1738 = int32(base.Ui32(v1729+v1728)>>(uint(v1728)%32)) & int32(2147450879)
	goto L410
L409:
	;
	v1738 = v1729
	goto L410
L410:
	;
	if v122 < v33 {
		goto L411
	} else {
		goto L412
	}
L411:
	;
	v1742 = int32(_a_F_VP8RecordCoeffTokens_2)
	goto L413
L412:
	;
	v1742 = int32(_a_F_VP8RecordCoeffTokens_1)
	goto L413
L413:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1663))) = v1738 + v1742
	if v33 <= v122 {
		v1789 = v1728
		goto L21
	} else {
		goto L414
	}
L414:
	;
	v1746 = v1667
	v1750 = v1663
	goto L47
L415:
	;
	goto L46
}
func F_VP8RecordCoeffs(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v25 int32
	_ = v25
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v154 int32
	_ = v154
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v163 int32
	_ = v163
	var v167 int32
	_ = v167
	var v171 int32
	_ = v171
	var v172 int32
	_ = v172
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v191 int32
	_ = v191
	var v197 int32
	_ = v197
	var v203 int32
	_ = v203
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v231 int32
	_ = v231
	var v234 int32
	_ = v234
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v265 int32
	_ = v265
	var v273 int32
	_ = v273
	v3 = int32(0)
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v19 = v12 + v13*int32(132) + l0*int32(44)
	v21 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	if v21 < v3 {
		v245 = v19
		v249 = v3
		v256 = *(*int32)(unsafe.Add(mBase, uint32(v245)))
		v257 = int32(1)
		if base.Ui32(int32(-131073)) < base.Ui32(v256) {
			v265 = int32(base.Ui32(v256+v257)>>(uint(v257)%32)) & int32(2147450879)
		} else {
			v265 = v256
		}
		*(*int32)(unsafe.Add(mBase, uint32(v245))) = v265 + int32(_a_F_VP8RecordCoeffs_0)
		v273 = v249
	} else {
		if v21 < v13 {
			v231 = v19
			v234 = v13
		} else {
			v25 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
			v28 = v19
			v31 = v13
			for {
				v39 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
				v40 = int32(1)
				if base.Ui32(int32(-131073)) < base.Ui32(v39) {
					v48 = int32(base.Ui32(v39+v40)>>(uint(v40)%32)) & int32(2147450879)
				} else {
					v48 = v39
				}
				*(*int32)(unsafe.Add(mBase, uint32(v28))) = v48 + int32(_a_F_VP8RecordCoeffs_1)
				v52 = int32(1)
				v53 = v31 + v52
				v55 = v31 << (uint(v52) % 32)
				v57 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v25+v55))))
				if v57 == int32(0) {
					v61 = v28
					v65 = v53
					v66 = v25 + int32(2) + v55
					for {
						v72 = *(*int32)(unsafe.Add(mBase, uint32(v61)+4))
						v73 = int32(1)
						if base.Ui32(int32(-131073)) < base.Ui32(v72) {
							v81 = int32(base.Ui32(v72+v73)>>(uint(v73)%32)) & int32(2147450879)
						} else {
							v81 = v72
						}
						*(*int32)(unsafe.Add(mBase, uint32(v61)+4)) = v81 + int32(_a_F_VP8RecordCoeffs_0)
						v85 = m.G23
						v87 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v85+v65))))
						v90 = v12 + v87*int32(132)
						v91 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v66))))
						v95 = v65 + int32(1)
						if v91 == int32(0) {
							v61 = v90
							v65 = v95
							v66 = v66 + int32(2)
							continue
						} else {
							break
						}
						break
					}
					v98 = v90
					v101 = v95
					v106 = v91
				} else {
					v98 = v28
					v101 = v53
					v106 = v57
				}
				v109 = int32(1)
				v110 = *(*int32)(unsafe.Add(mBase, uint32(v98)+4))
				if base.Ui32(int32(-131073)) < base.Ui32(v110) {
					v119 = int32(base.Ui32(v110+v109)>>(uint(v109)%32)) & int32(2147450879)
				} else {
					v119 = v110
				}
				*(*int32)(unsafe.Add(mBase, uint32(v98)+4)) = v119 + int32(_a_F_VP8RecordCoeffs_1)
				v123 = *(*int32)(unsafe.Add(mBase, uint32(v98)+8))
				v124 = int32(1)
				if base.Ui32(int32(-131073)) < base.Ui32(v123) {
					v132 = int32(base.Ui32(v123+v124)>>(uint(v124)%32)) & int32(2147450879)
				} else {
					v132 = v123
				}
				v135 = base.I32_extend16_s(v106)
				v139 = base.B2i32(base.Ui32(v135+int32(1)) < base.Ui32(int32(3)))
				if base.Ui32(v135+int32(1)) < base.Ui32(int32(3)) {
					v140 = int32(_a_F_VP8RecordCoeffs_0)
				} else {
					v140 = int32(_a_F_VP8RecordCoeffs_1)
				}
				*(*int32)(unsafe.Add(mBase, uint32(v98)+8)) = v132 + v140
				if base.Ui32(v135+int32(1)) < base.Ui32(int32(3)) {
					v218 = v109
				} else {
					v143 = int32(2)
					v145 = v135 >> (uint(int32(31)) % 32)
					v147 = v135 ^ v145 - v145
					v148 = int32(67)
					if base.Ui32(v147) < base.Ui32(v148) {
						v151 = v147
					} else {
						v151 = v148
					}
					v152 = int32(2)
					v154 = m.G1
					v159 = v151<<(uint(v152)%32) + (v154 + int32(_a_F_VP8RecordCoeffs_2)) + int32(-4)
					v160 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v159))))
					if base.Ui32(v160) < base.Ui32(v152) {
						v218 = v143
					} else {
						v163 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v159)+2)))
						v167 = v98 + int32(12)
						v171 = v160
						v172 = int32(0)
						for {
							if v171&int32(2) == int32(0) {
							} else {
								v182 = *(*int32)(unsafe.Add(mBase, uint32(v167)))
								v183 = int32(1)
								if base.Ui32(int32(-131073)) < base.Ui32(v182) {
									v191 = int32(base.Ui32(v182+v183)>>(uint(v183)%32)) & int32(2147450879)
								} else {
									v191 = v182
								}
								if int32(base.Ui32(v163)>>(uint(v172)%32))&int32(2) != 0 {
									v197 = int32(_a_F_VP8RecordCoeffs_1)
								} else {
									v197 = int32(_a_F_VP8RecordCoeffs_0)
								}
								*(*int32)(unsafe.Add(mBase, uint32(v167))) = v191 + v197
							}
							v203 = int32(1)
							if base.Ui32(int32(3)) < base.Ui32(v171) {
								v167 = v167 + int32(4)
								v171 = int32(base.Ui32(v171) >> (uint(v203) % 32))
								v172 = v172 + v203
								continue
							} else {
								break
							}
							break
						}
						v218 = v143
					}
				}
				v220 = m.G23
				v222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v220+v101))))
				v228 = v12 + v222*int32(132) + v218*int32(44)
				v229 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
				if v101 <= v229 {
					v28 = v228
					v31 = v101
					continue
				} else {
					break
				}
				break
			}
			v231 = v228
			v234 = v101
		}
		v242 = int32(1)
		if int32(15) < v234 {
			v273 = v242
		} else {
			v245 = v231
			v249 = v242
			v256 = *(*int32)(unsafe.Add(mBase, uint32(v245)))
			v257 = int32(1)
			if base.Ui32(int32(-131073)) < base.Ui32(v256) {
				v265 = int32(base.Ui32(v256+v257)>>(uint(v257)%32)) & int32(2147450879)
			} else {
				v265 = v256
			}
			*(*int32)(unsafe.Add(mBase, uint32(v245))) = v265 + int32(_a_F_VP8RecordCoeffs_0)
			v273 = v249
		}
	}
	return v273
}
func F_VP8SSIMDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_VP8SSIMDspInit[0])))
	v8 = m.G3
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8SSIMDspInit[1]))) = v11 + int32(63)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8SSIMDspInit[2]))) = v11 + int32(64)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8SSIMDspInit[3]))) = v11 + int32(65)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8SSIMDspInit[0]))) = v9
	}
	return
}
func F_VP8SetHistogramData(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
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
	var v138 int32
	_ = v138
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
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v150 int32
	_ = v150
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v170 int32
	_ = v170
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
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
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v231 int32
	_ = v231
	var v232 int32
	_ = v232
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v238 int32
	_ = v238
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v243 int32
	_ = v243
	var v244 int32
	_ = v244
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v249 int32
	_ = v249
	var v250 int32
	_ = v250
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v277 int32
	_ = v277
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v282 int32
	_ = v282
	var v283 int32
	_ = v283
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v289 int32
	_ = v289
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v297 int32
	_ = v297
	var v298 int32
	_ = v298
	var v300 int32
	_ = v300
	var v301 int32
	_ = v301
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v309 int32
	_ = v309
	var v310 int32
	_ = v310
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v318 int32
	_ = v318
	var v319 int32
	_ = v319
	v3 = int32(0)
	v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v97 = base.B2i32(v3 < v95)
	v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v102 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v104 = base.B2i32(v3 < v102)
	if v3 < v102 {
		v105 = int32(2)
	} else {
		v105 = v97 | base.B2i32(v98 < int32(1))
	}
	v106 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v108 = base.B2i32(int32(0) < v106)
	if int32(0) < v106 {
		v109 = int32(3)
	} else {
		v109 = v105
	}
	v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v112 = base.B2i32(int32(0) < v110)
	if int32(0) < v110 {
		v113 = int32(4)
	} else {
		v113 = v109
	}
	v114 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v116 = base.B2i32(int32(0) < v114)
	if int32(0) < v114 {
		v117 = int32(5)
	} else {
		v117 = v113
	}
	v118 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v120 = base.B2i32(int32(0) < v118)
	if int32(0) < v118 {
		v121 = int32(6)
	} else {
		v121 = v117
	}
	v122 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v124 = base.B2i32(int32(0) < v122)
	if int32(0) < v122 {
		v125 = int32(7)
	} else {
		v125 = v121
	}
	v126 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
	v128 = base.B2i32(int32(0) < v126)
	if int32(0) < v126 {
		v129 = int32(8)
	} else {
		v129 = v125
	}
	v130 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	v132 = base.B2i32(int32(0) < v130)
	if int32(0) < v130 {
		v133 = int32(9)
	} else {
		v133 = v129
	}
	v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v136 = base.B2i32(int32(0) < v134)
	if int32(0) < v134 {
		v137 = int32(10)
	} else {
		v137 = v133
	}
	v138 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v140 = base.B2i32(int32(0) < v138)
	if int32(0) < v138 {
		v141 = int32(11)
	} else {
		v141 = v137
	}
	v142 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
	v144 = base.B2i32(int32(0) < v142)
	if int32(0) < v142 {
		v145 = int32(12)
	} else {
		v145 = v141
	}
	v146 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v148 = base.B2i32(int32(0) < v146)
	if int32(0) < v146 {
		v149 = int32(13)
	} else {
		v149 = v145
	}
	v150 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	v152 = base.B2i32(int32(0) < v150)
	if int32(0) < v150 {
		v153 = int32(14)
	} else {
		v153 = v149
	}
	v154 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
	v156 = base.B2i32(int32(0) < v154)
	if int32(0) < v154 {
		v157 = int32(15)
	} else {
		v157 = v153
	}
	v158 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v160 = base.B2i32(int32(0) < v158)
	if int32(0) < v158 {
		v161 = int32(16)
	} else {
		v161 = v157
	}
	v162 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	v164 = base.B2i32(int32(0) < v162)
	if int32(0) < v162 {
		v165 = int32(17)
	} else {
		v165 = v161
	}
	v166 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
	v168 = base.B2i32(int32(0) < v166)
	if int32(0) < v166 {
		v169 = int32(18)
	} else {
		v169 = v165
	}
	v170 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
	v172 = base.B2i32(int32(0) < v170)
	if int32(0) < v170 {
		v173 = int32(19)
	} else {
		v173 = v169
	}
	v174 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
	v176 = base.B2i32(int32(0) < v174)
	if int32(0) < v174 {
		v177 = int32(20)
	} else {
		v177 = v173
	}
	v178 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	v180 = base.B2i32(int32(0) < v178)
	if int32(0) < v178 {
		v181 = int32(21)
	} else {
		v181 = v177
	}
	v182 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	v184 = base.B2i32(int32(0) < v182)
	if int32(0) < v182 {
		v185 = int32(22)
	} else {
		v185 = v181
	}
	v186 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
	v188 = base.B2i32(int32(0) < v186)
	if int32(0) < v186 {
		v189 = int32(23)
	} else {
		v189 = v185
	}
	v190 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	v192 = base.B2i32(int32(0) < v190)
	if int32(0) < v190 {
		v193 = int32(24)
	} else {
		v193 = v189
	}
	v194 = *(*int32)(unsafe.Add(mBase, uint32(l0)+100))
	v196 = base.B2i32(int32(0) < v194)
	if int32(0) < v194 {
		v197 = int32(25)
	} else {
		v197 = v193
	}
	v198 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
	v200 = base.B2i32(int32(0) < v198)
	if int32(0) < v198 {
		v201 = int32(26)
	} else {
		v201 = v197
	}
	v202 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
	v204 = base.B2i32(int32(0) < v202)
	if int32(0) < v202 {
		v205 = int32(27)
	} else {
		v205 = v201
	}
	v206 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
	v208 = base.B2i32(int32(0) < v206)
	if int32(0) < v206 {
		v209 = int32(28)
	} else {
		v209 = v205
	}
	v210 = *(*int32)(unsafe.Add(mBase, uint32(l0)+116))
	v212 = base.B2i32(int32(0) < v210)
	if int32(0) < v210 {
		v213 = int32(29)
	} else {
		v213 = v209
	}
	v214 = *(*int32)(unsafe.Add(mBase, uint32(l0)+120))
	v216 = base.B2i32(int32(0) < v214)
	if int32(0) < v214 {
		v217 = int32(30)
	} else {
		v217 = v213
	}
	v218 = *(*int32)(unsafe.Add(mBase, uint32(l0)+124))
	v220 = base.B2i32(int32(0) < v218)
	if int32(0) < v218 {
		v221 = int32(31)
	} else {
		v221 = v217
	}
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v221
	v223 = int32(0)
	if v223 < v98 {
		v226 = v98
	} else {
		v226 = v223
	}
	if v226 < v95 {
		v228 = v95
	} else {
		v228 = v226
	}
	if v3 < v95 {
		v229 = v228
	} else {
		v229 = v226
	}
	if v229 < v102 {
		v231 = v102
	} else {
		v231 = v229
	}
	if v3 < v102 {
		v232 = v231
	} else {
		v232 = v229
	}
	if v232 < v106 {
		v234 = v106
	} else {
		v234 = v232
	}
	if int32(0) < v106 {
		v235 = v234
	} else {
		v235 = v232
	}
	if v235 < v110 {
		v237 = v110
	} else {
		v237 = v235
	}
	if int32(0) < v110 {
		v238 = v237
	} else {
		v238 = v235
	}
	if v238 < v114 {
		v240 = v114
	} else {
		v240 = v238
	}
	if int32(0) < v114 {
		v241 = v240
	} else {
		v241 = v238
	}
	if v241 < v118 {
		v243 = v118
	} else {
		v243 = v241
	}
	if int32(0) < v118 {
		v244 = v243
	} else {
		v244 = v241
	}
	if v244 < v122 {
		v246 = v122
	} else {
		v246 = v244
	}
	if int32(0) < v122 {
		v247 = v246
	} else {
		v247 = v244
	}
	if v247 < v126 {
		v249 = v126
	} else {
		v249 = v247
	}
	if int32(0) < v126 {
		v250 = v249
	} else {
		v250 = v247
	}
	if v250 < v130 {
		v252 = v130
	} else {
		v252 = v250
	}
	if int32(0) < v130 {
		v253 = v252
	} else {
		v253 = v250
	}
	if v253 < v134 {
		v255 = v134
	} else {
		v255 = v253
	}
	if int32(0) < v134 {
		v256 = v255
	} else {
		v256 = v253
	}
	if v256 < v138 {
		v258 = v138
	} else {
		v258 = v256
	}
	if int32(0) < v138 {
		v259 = v258
	} else {
		v259 = v256
	}
	if v259 < v142 {
		v261 = v142
	} else {
		v261 = v259
	}
	if int32(0) < v142 {
		v262 = v261
	} else {
		v262 = v259
	}
	if v262 < v146 {
		v264 = v146
	} else {
		v264 = v262
	}
	if int32(0) < v146 {
		v265 = v264
	} else {
		v265 = v262
	}
	if v265 < v150 {
		v267 = v150
	} else {
		v267 = v265
	}
	if int32(0) < v150 {
		v268 = v267
	} else {
		v268 = v265
	}
	if v268 < v154 {
		v270 = v154
	} else {
		v270 = v268
	}
	if int32(0) < v154 {
		v271 = v270
	} else {
		v271 = v268
	}
	if v271 < v158 {
		v273 = v158
	} else {
		v273 = v271
	}
	if int32(0) < v158 {
		v274 = v273
	} else {
		v274 = v271
	}
	if v274 < v162 {
		v276 = v162
	} else {
		v276 = v274
	}
	if int32(0) < v162 {
		v277 = v276
	} else {
		v277 = v274
	}
	if v277 < v166 {
		v279 = v166
	} else {
		v279 = v277
	}
	if int32(0) < v166 {
		v280 = v279
	} else {
		v280 = v277
	}
	if v280 < v170 {
		v282 = v170
	} else {
		v282 = v280
	}
	if int32(0) < v170 {
		v283 = v282
	} else {
		v283 = v280
	}
	if v283 < v174 {
		v285 = v174
	} else {
		v285 = v283
	}
	if int32(0) < v174 {
		v286 = v285
	} else {
		v286 = v283
	}
	if v286 < v178 {
		v288 = v178
	} else {
		v288 = v286
	}
	if int32(0) < v178 {
		v289 = v288
	} else {
		v289 = v286
	}
	if v289 < v182 {
		v291 = v182
	} else {
		v291 = v289
	}
	if int32(0) < v182 {
		v292 = v291
	} else {
		v292 = v289
	}
	if v292 < v186 {
		v294 = v186
	} else {
		v294 = v292
	}
	if int32(0) < v186 {
		v295 = v294
	} else {
		v295 = v292
	}
	if v295 < v190 {
		v297 = v190
	} else {
		v297 = v295
	}
	if int32(0) < v190 {
		v298 = v297
	} else {
		v298 = v295
	}
	if v298 < v194 {
		v300 = v194
	} else {
		v300 = v298
	}
	if int32(0) < v194 {
		v301 = v300
	} else {
		v301 = v298
	}
	if v301 < v198 {
		v303 = v198
	} else {
		v303 = v301
	}
	if int32(0) < v198 {
		v304 = v303
	} else {
		v304 = v301
	}
	if v304 < v202 {
		v306 = v202
	} else {
		v306 = v304
	}
	if int32(0) < v202 {
		v307 = v306
	} else {
		v307 = v304
	}
	if v307 < v206 {
		v309 = v206
	} else {
		v309 = v307
	}
	if int32(0) < v206 {
		v310 = v309
	} else {
		v310 = v307
	}
	if v310 < v210 {
		v312 = v210
	} else {
		v312 = v310
	}
	if int32(0) < v210 {
		v313 = v312
	} else {
		v313 = v310
	}
	if v313 < v214 {
		v315 = v214
	} else {
		v315 = v313
	}
	if int32(0) < v214 {
		v316 = v315
	} else {
		v316 = v313
	}
	if v316 < v218 {
		v318 = v218
	} else {
		v318 = v316
	}
	if int32(0) < v218 {
		v319 = v318
	} else {
		v319 = v316
	}
	*(*int32)(unsafe.Add(mBase, uint32(l1))) = v319
	return
}
func F_VP8SetIntra16Mode(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	v8 = l1 & int32(255) * int32(16843009)
	*(*int32)(unsafe.Add(mBase, uint32(v4))) = v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v11 = *(*int32)(unsafe.Add(mBase, uint32(v10)+48))
	v12 = v4 + v11
	*(*int32)(unsafe.Add(mBase, uint32(v12))) = v8
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(v14)+48))
	v16 = v12 + v15
	*(*int32)(unsafe.Add(mBase, uint32(v16))) = v8
	v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v19 = *(*int32)(unsafe.Add(mBase, uint32(v18)+48))
	*(*int32)(unsafe.Add(mBase, uint32(v16+v19))) = v8
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v22))))
	v27 = v23&int32(252) | int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v22))) = uint8(v27)
	return
}
func F_VP8SetIntra4Mode(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	*(*int32)(unsafe.Add(mBase, uint32(v4))) = v5
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(v7)+48))
	v9 = v4 + v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v10
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v12)+48))
	v14 = v9 + v13
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v14))) = v15
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v18 = *(*int32)(unsafe.Add(mBase, uint32(v17)+48))
	v20 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v14+v18))) = v20
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v22))))
	v25 = v23 & int32(252)
	*(*uint8)(unsafe.Add(mBase, uint32(v22))) = uint8(v25)
	return
}
func F_VP8SetIntraUVMode(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	var v11 int32
	_ = v11
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v4 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3))))
	v11 = v4&int32(243) | l1<<(uint(int32(2))%32)&int32(12)
	*(*uint8)(unsafe.Add(mBase, uint32(v3))) = uint8(v11)
	return
}
func F_VP8SetSegment(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	var v11 int32
	_ = v11
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v4 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3))))
	v11 = v4&int32(159) | l1<<(uint(int32(5))%32)&int32(96)
	*(*uint8)(unsafe.Add(mBase, uint32(v3))) = uint8(v11)
	return
}
func F_VP8SetSegmentParams(m *base.Module, l0 int32, l1 float32) {
	mBase := m.M
	_ = mBase
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v32 float64
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v43 float64
	_ = v43
	var v46 float64
	_ = v46
	var v55 float64
	_ = v55
	var v56 float64
	_ = v56
	var v57 float64
	_ = v57
	var v65 float64
	_ = v65
	var v67 float64
	_ = v67
	var v68 float64
	_ = v68
	var v74 int32
	_ = v74
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v112 float64
	_ = v112
	var v115 int32
	_ = v115
	var v120 float64
	_ = v120
	var v123 float64
	_ = v123
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v139 int32
	_ = v139
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v148 int32
	_ = v148
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v181 int32
	_ = v181
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v219 int32
	_ = v219
	var v233 int32
	_ = v233
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v306 int32
	_ = v306
	var v332 int32
	_ = v332
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v358 int32
	_ = v358
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v380 int32
	_ = v380
	var v383 int32
	_ = v383
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v391 int32
	_ = v391
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v399 int32
	_ = v399
	var v400 int32
	_ = v400
	var v403 int32
	_ = v403
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v411 int32
	_ = v411
	var v413 int32
	_ = v413
	var v414 int32
	_ = v414
	var v420 int32
	_ = v420
	var v423 int32
	_ = v423
	var v425 int32
	_ = v425
	var v427 int32
	_ = v427
	var v429 int32
	_ = v429
	var v432 int32
	_ = v432
	var v433 int32
	_ = v433
	var v436 int32
	_ = v436
	var v439 int32
	_ = v439
	var v442 int32
	_ = v442
	var v443 int32
	_ = v443
	var v444 int32
	_ = v444
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v451 int32
	_ = v451
	var v455 int32
	_ = v455
	var v457 int32
	_ = v457
	var v458 int32
	_ = v458
	var v464 int32
	_ = v464
	var v467 int32
	_ = v467
	var v469 int32
	_ = v469
	var v471 int32
	_ = v471
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v478 int32
	_ = v478
	var v481 int32
	_ = v481
	var v484 int32
	_ = v484
	var v485 int32
	_ = v485
	var v486 int32
	_ = v486
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v497 int32
	_ = v497
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v506 int32
	_ = v506
	var v509 int32
	_ = v509
	var v511 int32
	_ = v511
	var v513 int32
	_ = v513
	var v516 int32
	_ = v516
	var v517 int32
	_ = v517
	var v520 int32
	_ = v520
	var v523 int32
	_ = v523
	var v525 int32
	_ = v525
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
	var v538 int32
	_ = v538
	var v540 int32
	_ = v540
	var v541 int32
	_ = v541
	var v547 int32
	_ = v547
	var v550 int32
	_ = v550
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v558 int32
	_ = v558
	var v559 int32
	_ = v559
	var v564 int32
	_ = v564
	var v567 int32
	_ = v567
	var v568 int32
	_ = v568
	var v571 int32
	_ = v571
	var v574 int32
	_ = v574
	var v578 int32
	_ = v578
	var v581 int64
	_ = v581
	var v585 int64
	_ = v585
	var v587 int32
	_ = v587
	var v588 int32
	_ = v588
	var v591 int32
	_ = v591
	var v596 int32
	_ = v596
	var v608 int32
	_ = v608
	var v610 int32
	_ = v610
	var v625 int32
	_ = v625
	var v627 int32
	_ = v627
	var v628 int32
	_ = v628
	var v632 int32
	_ = v632
	var v633 int32
	_ = v633
	var v636 int32
	_ = v636
	var v643 int32
	_ = v643
	var v646 int32
	_ = v646
	var v653 int32
	_ = v653
	var v665 int32
	_ = v665
	var v667 int32
	_ = v667
	var v682 int32
	_ = v682
	var v684 int32
	_ = v684
	var v685 int32
	_ = v685
	var v690 int32
	_ = v690
	var v693 int32
	_ = v693
	var v697 int32
	_ = v697
	var v701 int32
	_ = v701
	var v710 int32
	_ = v710
	var v722 int32
	_ = v722
	var v724 int32
	_ = v724
	var v739 int32
	_ = v739
	var v741 int32
	_ = v741
	var v742 int32
	_ = v742
	var v747 int32
	_ = v747
	var v752 int32
	_ = v752
	var v756 int32
	_ = v756
	var v765 int32
	_ = v765
	var v786 int32
	_ = v786
	var v787 int32
	_ = v787
	var v788 int32
	_ = v788
	var v795 int32
	_ = v795
	var v797 int32
	_ = v797
	var v800 int32
	_ = v800
	var v801 int32
	_ = v801
	var v807 int32
	_ = v807
	var v814 int32
	_ = v814
	var v817 int32
	_ = v817
	var v836 int32
	_ = v836
	var v838 int32
	_ = v838
	var v851 int32
	_ = v851
	var v852 int32
	_ = v852
	var v853 int32
	_ = v853
	var v854 int32
	_ = v854
	var v856 int32
	_ = v856
	var v859 int32
	_ = v859
	var v860 int32
	_ = v860
	var v862 int32
	_ = v862
	var v864 int32
	_ = v864
	var v866 int32
	_ = v866
	var v868 int32
	_ = v868
	var v871 int32
	_ = v871
	var v872 int32
	_ = v872
	var v878 int32
	_ = v878
	var v885 int32
	_ = v885
	var v890 int32
	_ = v890
	var v921 int32
	_ = v921
	var v923 int32
	_ = v923
	var v926 int32
	_ = v926
	var v940 int32
	_ = v940
	var v944 int32
	_ = v944
	var v955 int32
	_ = v955
	var v956 int32
	_ = v956
	var v960 int32
	_ = v960
	var v967 int32
	_ = v967
	var v1003 int32
	_ = v1003
	var v1007 int32
	_ = v1007
	var v1018 int32
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1023 int32
	_ = v1023
	var v1027 int32
	_ = v1027
	var v1031 int32
	_ = v1031
	var v1035 int32
	_ = v1035
	var v1061 int32
	_ = v1061
	var v1064 int32
	_ = v1064
	var v1065 int32
	_ = v1065
	var v1066 int32
	_ = v1066
	var v1067 int32
	_ = v1067
	var v1076 int32
	_ = v1076
	var v1079 int32
	_ = v1079
	var v1097 int32
	_ = v1097
	var v1100 int32
	_ = v1100
	var v1102 int32
	_ = v1102
	var v1103 int32
	_ = v1103
	var v1106 int32
	_ = v1106
	var v1107 int32
	_ = v1107
	var v1110 int32
	_ = v1110
	var v1114 int32
	_ = v1114
	var v1117 int32
	_ = v1117
	var v1118 int32
	_ = v1118
	var v1119 int32
	_ = v1119
	var v1120 int32
	_ = v1120
	var v1123 int32
	_ = v1123
	var v1124 int32
	_ = v1124
	var v1127 int32
	_ = v1127
	var v1129 int32
	_ = v1129
	var v1132 int32
	_ = v1132
	var v1133 int32
	_ = v1133
	var v1134 int32
	_ = v1134
	var v1135 int32
	_ = v1135
	var v1138 int32
	_ = v1138
	var v1139 int32
	_ = v1139
	var v1142 int32
	_ = v1142
	var v1144 int32
	_ = v1144
	var v1147 int32
	_ = v1147
	var v1148 int32
	_ = v1148
	var v1149 int32
	_ = v1149
	var v1150 int32
	_ = v1150
	var v1153 int32
	_ = v1153
	var v1154 int32
	_ = v1154
	var v1157 int32
	_ = v1157
	var v1159 int32
	_ = v1159
	var v1161 int32
	_ = v1161
	var v1167 int32
	_ = v1167
	var v1168 int32
	_ = v1168
	var v1169 int32
	_ = v1169
	var v1172 int32
	_ = v1172
	var v1173 int32
	_ = v1173
	var v1176 int32
	_ = v1176
	var v1180 int32
	_ = v1180
	var v1184 int32
	_ = v1184
	var v1185 int32
	_ = v1185
	var v1186 int32
	_ = v1186
	var v1189 int32
	_ = v1189
	var v1190 int32
	_ = v1190
	var v1193 int32
	_ = v1193
	var v1197 int32
	_ = v1197
	var v1199 int32
	_ = v1199
	var v1205 int32
	_ = v1205
	var v1209 int32
	_ = v1209
	var v1210 int32
	_ = v1210
	var v1211 int32
	_ = v1211
	var v1214 int32
	_ = v1214
	var v1219 int32
	_ = v1219
	var v1225 int32
	_ = v1225
	var v1226 int32
	_ = v1226
	var v1227 int32
	_ = v1227
	var v1229 int32
	_ = v1229
	var v1231 int32
	_ = v1231
	var v1236 int32
	_ = v1236
	var v1238 int32
	_ = v1238
	var v1240 int32
	_ = v1240
	var v1246 int32
	_ = v1246
	var v1285 int32
	_ = v1285
	var v1299 int32
	_ = v1299
	var v1303 int32
	_ = v1303
	var v1343 int32
	_ = v1343
	var v1345 int32
	_ = v1345
	var v1346 int32
	_ = v1346
	var v1350 int32
	_ = v1350
	var v1353 int32
	_ = v1353
	var v1379 int32
	_ = v1379
	var v1382 int32
	_ = v1382
	var v1395 int32
	_ = v1395
	var v1399 int32
	_ = v1399
	var v1403 int32
	_ = v1403
	var v1407 int32
	_ = v1407
	var v1411 int32
	_ = v1411
	var v1415 int32
	_ = v1415
	var v1417 int32
	_ = v1417
	var v1418 int32
	_ = v1418
	var v1439 int32
	_ = v1439
	var v1441 int32
	_ = v1441
	var v1446 int32
	_ = v1446
	var v1450 int32
	_ = v1450
	var v1451 int32
	_ = v1451
	var v1452 int32
	_ = v1452
	var v1455 int32
	_ = v1455
	var v1460 int32
	_ = v1460
	var v1466 int32
	_ = v1466
	var v1467 int32
	_ = v1467
	var v1468 int32
	_ = v1468
	var v1470 int32
	_ = v1470
	var v1472 int32
	_ = v1472
	var v1477 int32
	_ = v1477
	var v1479 int32
	_ = v1479
	var v1481 int32
	_ = v1481
	var v1487 int32
	_ = v1487
	var v1526 int32
	_ = v1526
	var v1540 int32
	_ = v1540
	var v1544 int32
	_ = v1544
	var v1549 int64
	_ = v1549
	var v1563 int32
	_ = v1563
	var v1564 int32
	_ = v1564
	var v1565 int32
	_ = v1565
	var v1567 int32
	_ = v1567
	var v1569 int32
	_ = v1569
	var v1680 int32
	_ = v1680
	var v1682 int32
	_ = v1682
	var v1687 int32
	_ = v1687
	var v1691 int32
	_ = v1691
	var v1692 int32
	_ = v1692
	var v1693 int32
	_ = v1693
	var v1696 int32
	_ = v1696
	var v1701 int32
	_ = v1701
	var v1707 int32
	_ = v1707
	var v1708 int32
	_ = v1708
	var v1709 int32
	_ = v1709
	var v1711 int32
	_ = v1711
	var v1713 int32
	_ = v1713
	var v1718 int32
	_ = v1718
	var v1720 int32
	_ = v1720
	var v1722 int32
	_ = v1722
	var v1728 int32
	_ = v1728
	var v1767 int32
	_ = v1767
	var v1781 int32
	_ = v1781
	var v1785 int32
	_ = v1785
	var v1790 int64
	_ = v1790
	var v1804 int32
	_ = v1804
	var v1805 int32
	_ = v1805
	var v1806 int32
	_ = v1806
	var v1808 int32
	_ = v1808
	var v1810 int32
	_ = v1810
	var v1921 int32
	_ = v1921
	var v1923 int32
	_ = v1923
	var v1924 int32
	_ = v1924
	var v1925 int32
	_ = v1925
	var v1929 int32
	_ = v1929
	var v1932 int32
	_ = v1932
	var v1935 int32
	_ = v1935
	var v1938 int32
	_ = v1938
	var v1939 int32
	_ = v1939
	var v1942 int32
	_ = v1942
	var v1943 int32
	_ = v1943
	var v1945 int32
	_ = v1945
	var v1948 int32
	_ = v1948
	var v1952 int32
	_ = v1952
	var v1955 int32
	_ = v1955
	var v1960 int32
	_ = v1960
	var v1961 int32
	_ = v1961
	var v1963 int32
	_ = v1963
	var v1966 int32
	_ = v1966
	var v1970 int32
	_ = v1970
	var v2010 int32
	_ = v2010
	var v2023 int32
	_ = v2023
	v26 = m.G0
	v28 = v26 - int32(16)
	m.G0 = v28
	v32 = base.F64_div(base.F64_promote_f32(l1), float64(100))
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v34 = *(*int32)(unsafe.Add(mBase, uint32(v33)+28))
	v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v36 = *(*int32)(unsafe.Add(mBase, uint32(v33)+80))
	if v36 == int32(0) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	if int32(0) < v35 {
		goto L14
	} else {
		goto L15
	}
L2:
	;
	if base.F64_lt(v32, float64(0.75)) != 0 {
		goto L9
	} else {
		goto L10
	}
L3:
	;
	v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3388))
	v43 = base.F64_div(base.F64_convert_i32_s(v40), float64(255))
	if base.F64_gt(v43, float64(0.85)) != 0 {
		v56 = float64(0.4)
		goto L4
	} else {
		goto L5
	}
L4:
	;
	v57 = F_pow(m, v32, v56)
	mBase = m.M
	v68 = v57
	goto L1
L5:
	;
	v46 = float64(0.9)
	if base.F64_lt(v43, float64(0.3)) != 0 {
		goto L6
	} else {
		goto L7
	}
L6:
	;
	v55 = v46
	goto L8
L7:
	;
	v55 = base.F64_add(base.F64_mul(base.F64_add(v43, float64(-0.3)), float64(-0.9090909090909091)), v46)
	goto L8
L8:
	;
	v56 = v55
	goto L4
L9:
	;
	v65 = base.F64_mul(v32, float64(0.6666666666666666))
	goto L11
L10:
	;
	v65 = base.F64_add(base.F64_add(v32, v32), float64(-1))
	goto L11
L11:
	;
	v67 = F_pow(m, v65, float64(0.3333333333333333))
	mBase = m.M
	v68 = v67
	goto L1
L12:
	;
	v358 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3404)) = v358
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3396)) = int64(0)
	v364 = base.I32_div_s(v34, int32(-25))
	v365 = int32(15)
	if v364 < v365 {
		goto L39
	} else {
		goto L40
	}
L13:
	;
	v181 = (int32(4) - v35) & int32(7)
	if v181 == int32(0) {
		v233 = v35
		goto L29
	} else {
		goto L30
	}
L14:
	;
	v96 = l0 + int32(1088)
	v97 = v35
	goto L16
L15:
	;
	v74 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3384)) = v74
	v162 = l0 + int32(408)
	v163 = v74
	v164 = l0 + int32(1088)
	goto L13
L16:
	;
	v112 = float64(1)
	v115 = *(*int32)(unsafe.Add(mBase, uint32(v96+int32(-8))))
	v120 = F_pow(m, v68, base.F64_add(base.F64_mul(base.F64_mul(base.F64_div(base.F64_mul(base.F64_convert_i32_s(v34), float64(0.9)), float64(100)), float64(-0.0078125)), base.F64_convert_i32_s(v115)), v112))
	mBase = m.M
	v123 = base.F64_mul(base.F64_sub(v112, v120), float64(127))
	if base.F64_lt(base.F64_abs(v123), float64(2.147483648e+09)) == int32(0) {
		goto L19
	} else {
		goto L20
	}
L17:
	;
	v145 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3384)) = v145
	v148 = l0 + int32(408)
	if int32(3) < v35 {
		v342 = v148
		v343 = v145
		goto L12
	} else {
		goto L28
	}
L18:
	;
	v132 = int32(127)
	if v131 < v132 {
		goto L21
	} else {
		goto L22
	}
L19:
	;
	v131 = int32(-2147483648)
	goto L18
L20:
	;
	v129 = base.I32_trunc_f64_s(v123)
	v131 = v129
	goto L18
L21:
	;
	v135 = v131
	goto L23
L22:
	;
	v135 = v132
	goto L23
L23:
	;
	v136 = int32(0)
	if v136 < v135 {
		goto L24
	} else {
		goto L25
	}
L24:
	;
	v139 = v135
	goto L26
L25:
	;
	v139 = v136
	goto L26
L26:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v96))) = v139
	v144 = v97 + int32(-1)
	if v144 != 0 {
		v96 = v96 + int32(744)
		v97 = v144
		goto L16
	} else {
		goto L27
	}
L27:
	;
	goto L17
L28:
	;
	v162 = v148
	v163 = v145
	v164 = l0 + int32(1088)
	goto L13
L29:
	;
	if base.Ui32(v35+int32(3)) < base.Ui32(int32(7)) {
		goto L34
	} else {
		goto L35
	}
L30:
	;
	v202 = v35*int32(744) + l0 + int32(1088)
	v204 = v181
	goto L31
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v202))) = v163
	v219 = v204 + int32(-1)
	if v219 != 0 {
		v202 = v202 + int32(744)
		v204 = v219
		goto L31
	} else {
		goto L33
	}
L32:
	;
	v233 = v35 + v181
	goto L29
L33:
	;
	goto L32
L34:
	;
	v332 = *(*int32)(unsafe.Add(mBase, uint32(v164)))
	v342 = v162
	v343 = v332
	goto L12
L35:
	;
	v266 = l0 + v233*int32(744)
	v268 = v233 + int32(-4)
	goto L36
L36:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v266)+uint32(_c_F_VP8SetSegmentParams[0]))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266)+uint32(_c_F_VP8SetSegmentParams[1]))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266)+uint32(_c_F_VP8SetSegmentParams[2]))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266+int32(4064)))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266+int32(3320)))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266+int32(2576)))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266+int32(1832)))) = v163
	*(*int32)(unsafe.Add(mBase, uint32(v266+int32(1088)))) = v163
	v306 = v268 + int32(8)
	if v306 != 0 {
		v266 = v266 + int32(_a_F_VP8SetSegmentParams_0)
		v268 = v306
		goto L36
	} else {
		goto L38
	}
L37:
	;
	goto L34
L38:
	;
	goto L37
L39:
	;
	v368 = v364
	goto L41
L40:
	;
	v368 = v365
	goto L41
L41:
	;
	v369 = int32(-15)
	if v369 < v368 {
		goto L42
	} else {
		goto L43
	}
L42:
	;
	v372 = v368
	goto L44
L43:
	;
	v372 = v369
	goto L44
L44:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3408)) = v372
	v374 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3392))
	v380 = base.I32_div_s(v374*int32(10)+int32(-640), int32(70))
	v383 = base.I32_div_s(v380*v34, int32(100))
	v384 = int32(6)
	if v383 < v384 {
		goto L45
	} else {
		goto L46
	}
L45:
	;
	v387 = v383
	goto L47
L46:
	;
	v387 = v384
	goto L47
L47:
	;
	v388 = int32(-4)
	if v388 < v387 {
		goto L48
	} else {
		goto L49
	}
L48:
	;
	v391 = v387
	goto L50
L49:
	;
	v391 = v388
	goto L50
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3412)) = v391
	v394 = m.G1
	v395 = *(*int32)(unsafe.Add(mBase, uint32(v33)+32))
	v397 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v399 = v394 + int32(_a_F_VP8SetSegmentParams_1)
	v400 = int32(127)
	if v343 < v400 {
		goto L51
	} else {
		goto L52
	}
L51:
	;
	v403 = v343
	goto L53
L52:
	;
	v403 = v400
	goto L53
L53:
	;
	v404 = int32(0)
	if v404 < v403 {
		goto L54
	} else {
		goto L55
	}
L54:
	;
	v407 = v403
	goto L56
L55:
	;
	v407 = v404
	goto L56
L56:
	;
	v411 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v399+v407<<(uint(int32(1))%32)))))
	v413 = int32(base.Ui32(v411) >> (uint(int32(2)) % 32))
	v414 = m.G1
	v420 = int32(63)
	if v413 < v420 {
		goto L58
	} else {
		goto L59
	}
L57:
	;
	v427 = v395 * int32(5)
	v429 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1084))
	v432 = base.I32_div_s(v425*v427, v429+int32(256))
	v433 = int32(63)
	if v432 < v433 {
		goto L61
	} else {
		goto L62
	}
L58:
	;
	v423 = v413
	goto L60
L59:
	;
	v423 = v420
	goto L60
L60:
	;
	v425 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v414+int32(_a_F_VP8SetSegmentParams_2)+v397<<(uint(int32(6))%32)+v423))))
	goto L57
L61:
	;
	v436 = v432
	goto L63
L62:
	;
	v436 = v433
	goto L63
L63:
	;
	if v432 < int32(2) {
		goto L64
	} else {
		goto L65
	}
L64:
	;
	v439 = int32(0)
	goto L66
L65:
	;
	v439 = v436
	goto L66
L66:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+1092)) = v439
	v442 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v443 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1832))
	v444 = int32(127)
	if v443 < v444 {
		goto L67
	} else {
		goto L68
	}
L67:
	;
	v447 = v443
	goto L69
L68:
	;
	v447 = v444
	goto L69
L69:
	;
	v448 = int32(0)
	if v448 < v447 {
		goto L70
	} else {
		goto L71
	}
L70:
	;
	v451 = v447
	goto L72
L71:
	;
	v451 = v448
	goto L72
L72:
	;
	v455 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v399+v451<<(uint(int32(1))%32)))))
	v457 = int32(base.Ui32(v455) >> (uint(int32(2)) % 32))
	v458 = m.G1
	v464 = int32(63)
	if v457 < v464 {
		goto L74
	} else {
		goto L75
	}
L73:
	;
	v471 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1828))
	v474 = base.I32_div_s(v469*v427, v471+int32(256))
	v475 = int32(63)
	if v474 < v475 {
		goto L77
	} else {
		goto L78
	}
L74:
	;
	v467 = v457
	goto L76
L75:
	;
	v467 = v464
	goto L76
L76:
	;
	v469 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v458+int32(_a_F_VP8SetSegmentParams_2)+v442<<(uint(int32(6))%32)+v467))))
	goto L73
L77:
	;
	v478 = v474
	goto L79
L78:
	;
	v478 = v475
	goto L79
L79:
	;
	if v474 < int32(2) {
		goto L80
	} else {
		goto L81
	}
L80:
	;
	v481 = int32(0)
	goto L82
L81:
	;
	v481 = v478
	goto L82
L82:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+1836)) = v481
	v484 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v485 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2576))
	v486 = int32(127)
	if v485 < v486 {
		goto L83
	} else {
		goto L84
	}
L83:
	;
	v489 = v485
	goto L85
L84:
	;
	v489 = v486
	goto L85
L85:
	;
	v490 = int32(0)
	if v490 < v489 {
		goto L86
	} else {
		goto L87
	}
L86:
	;
	v493 = v489
	goto L88
L87:
	;
	v493 = v490
	goto L88
L88:
	;
	v497 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v399+v493<<(uint(int32(1))%32)))))
	v499 = int32(base.Ui32(v497) >> (uint(int32(2)) % 32))
	v500 = m.G1
	v506 = int32(63)
	if v499 < v506 {
		goto L90
	} else {
		goto L91
	}
L89:
	;
	v513 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2572))
	v516 = base.I32_div_s(v511*v427, v513+int32(256))
	v517 = int32(63)
	if v516 < v517 {
		goto L93
	} else {
		goto L94
	}
L90:
	;
	v509 = v499
	goto L92
L91:
	;
	v509 = v506
	goto L92
L92:
	;
	v511 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v500+int32(_a_F_VP8SetSegmentParams_2)+v484<<(uint(int32(6))%32)+v509))))
	goto L89
L93:
	;
	v520 = v516
	goto L95
L94:
	;
	v520 = v517
	goto L95
L95:
	;
	if v516 < int32(2) {
		goto L96
	} else {
		goto L97
	}
L96:
	;
	v523 = int32(0)
	goto L98
L97:
	;
	v523 = v520
	goto L98
L98:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+2580)) = v523
	v525 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v526 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3320))
	v527 = int32(127)
	if v526 < v527 {
		goto L99
	} else {
		goto L100
	}
L99:
	;
	v530 = v526
	goto L101
L100:
	;
	v530 = v527
	goto L101
L101:
	;
	v531 = int32(0)
	if v531 < v530 {
		goto L102
	} else {
		goto L103
	}
L102:
	;
	v534 = v530
	goto L104
L103:
	;
	v534 = v531
	goto L104
L104:
	;
	v538 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v399+v534<<(uint(int32(1))%32)))))
	v540 = int32(base.Ui32(v538) >> (uint(int32(2)) % 32))
	v541 = m.G1
	v547 = int32(63)
	if v540 < v547 {
		goto L106
	} else {
		goto L107
	}
L105:
	;
	v553 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1092))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v553
	v555 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v556 = *(*int32)(unsafe.Add(mBase, uint32(v555)+36))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v556
	v558 = *(*int32)(unsafe.Add(mBase, uint32(v555)+40))
	v559 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = base.B2i32(v558 == v559)
	v564 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3316))
	v567 = base.I32_div_s(v552*v427, v564+int32(256))
	v568 = int32(63)
	if v567 < v568 {
		goto L109
	} else {
		goto L110
	}
L106:
	;
	v550 = v540
	goto L108
L107:
	;
	v550 = v547
	goto L108
L108:
	;
	v552 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v541+int32(_a_F_VP8SetSegmentParams_2)+v525<<(uint(int32(6))%32)+v550))))
	goto L105
L109:
	;
	v571 = v567
	goto L111
L110:
	;
	v571 = v568
	goto L111
L111:
	;
	if v567 < int32(2) {
		goto L112
	} else {
		goto L113
	}
L112:
	;
	v574 = v559
	goto L114
L113:
	;
	v574 = v571
	goto L114
L114:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3324)) = v574
	if v35 < int32(2) {
		goto L115
	} else {
		goto L116
	}
L115:
	;
	v1061 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[3])))
	if v1061 < int32(4) {
		v1066 = v358
		goto L169
	} else {
		goto L170
	}
L116:
	;
	v578 = m.G1
	v581 = *(*int64)(unsafe.Add(mBase, uint32(v578)+uint32(_c_F_VP8SetSegmentParams[4])))
	*(*int64)(unsafe.Add(mBase, uint32(v28))) = v581
	v585 = *(*int64)(unsafe.Add(mBase, uint32(v578)+uint32(_c_F_VP8SetSegmentParams[5])))
	*(*int64)(unsafe.Add(mBase, uint32(v28)+8)) = v585
	v587 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v588 = int32(4)
	if v587 < v588 {
		goto L117
	} else {
		goto L118
	}
L117:
	;
	v591 = v587
	goto L119
L118:
	;
	v591 = v588
	goto L119
L119:
	;
	if v587 < int32(2) {
		v765 = int32(1)
		goto L120
	} else {
		goto L121
	}
L120:
	;
	if v591 <= v765 {
		goto L115
	} else {
		goto L151
	}
L121:
	;
	v596 = *(*int32)(unsafe.Add(mBase, uint32(v342)+1424))
	v608 = int32(0)
	v610 = v342 + int32(684)
	goto L124
L122:
	;
	if v587 == int32(2) {
		v765 = v643
		goto L120
	} else {
		goto L130
	}
L123:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v28)+4)) = v608
	v643 = int32(1)
	v646 = int32(0)
	goto L122
L124:
	;
	v625 = *(*int32)(unsafe.Add(mBase, uint32(v610+int32(-4))))
	if v596 != v625 {
		goto L126
	} else {
		goto L127
	}
L125:
	;
	v636 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+4)) = v636
	v643 = int32(2)
	v646 = v636
	goto L122
L126:
	;
	v632 = int32(1)
	v633 = v608 + v632
	if v633 != v632 {
		v608 = v633
		v610 = v610 + int32(744)
		goto L124
	} else {
		goto L129
	}
L127:
	;
	v627 = *(*int32)(unsafe.Add(mBase, uint32(v342)+1428))
	v628 = *(*int32)(unsafe.Add(mBase, uint32(v610)))
	if v627 == v628 {
		goto L123
	} else {
		goto L128
	}
L128:
	;
	goto L126
L129:
	;
	goto L125
L130:
	;
	v653 = *(*int32)(unsafe.Add(mBase, uint32(v342)+2168))
	v665 = int32(0)
	v667 = v342 + int32(684)
	goto L133
L131:
	;
	if v587 == int32(3) {
		v765 = v701
		goto L120
	} else {
		goto L141
	}
L132:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v28)+8)) = v665
	v701 = v643
	goto L131
L133:
	;
	v682 = *(*int32)(unsafe.Add(mBase, uint32(v667+int32(-4))))
	if v653 != v682 {
		goto L135
	} else {
		goto L136
	}
L134:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v28)+8)) = v643
	if v646 != 0 {
		goto L139
	} else {
		goto L140
	}
L135:
	;
	v690 = v665 + int32(1)
	if v643 != v690 {
		v665 = v690
		v667 = v667 + int32(744)
		goto L133
	} else {
		goto L138
	}
L136:
	;
	v684 = *(*int32)(unsafe.Add(mBase, uint32(v342)+2172))
	v685 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	if v684 == v685 {
		goto L132
	} else {
		goto L137
	}
L137:
	;
	goto L135
L138:
	;
	goto L134
L139:
	;
	v701 = v643 + int32(1)
	goto L131
L140:
	;
	v693 = int32(744)
	v697 = F_memcpy(m, v342+v643*v693, v342+int32(1488), v693)
	mBase = m.M
	goto L139
L141:
	;
	v710 = *(*int32)(unsafe.Add(mBase, uint32(v342)+2912))
	v722 = int32(0)
	v724 = v342 + int32(684)
	goto L143
L142:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v28)+12)) = v722
	v765 = v701
	goto L120
L143:
	;
	v739 = *(*int32)(unsafe.Add(mBase, uint32(v724+int32(-4))))
	if v710 != v739 {
		goto L145
	} else {
		goto L146
	}
L144:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v28)+12)) = v701
	if v701 == int32(3) {
		goto L149
	} else {
		goto L150
	}
L145:
	;
	v747 = v722 + int32(1)
	if v701 != v747 {
		v722 = v747
		v724 = v724 + int32(744)
		goto L143
	} else {
		goto L148
	}
L146:
	;
	v741 = *(*int32)(unsafe.Add(mBase, uint32(v342)+2916))
	v742 = *(*int32)(unsafe.Add(mBase, uint32(v724)))
	if v741 == v742 {
		goto L142
	} else {
		goto L147
	}
L147:
	;
	goto L145
L148:
	;
	goto L144
L149:
	;
	v765 = v701 + int32(1)
	goto L120
L150:
	;
	v752 = int32(744)
	v756 = F_memcpy(m, v342+v701*v752, v342+int32(2232), v752)
	mBase = m.M
	goto L149
L151:
	;
	v786 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v787 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v788 = v786 * v787
	if v788 < int32(1) {
		goto L152
	} else {
		goto L153
	}
L152:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v765
	v921 = v765*int32(744) + v342
	v923 = v921 + int32(-744)
	v926 = (v591 - v765) & int32(3)
	if v926 == int32(0) {
		v967 = v765
		goto L160
	} else {
		goto L161
	}
L153:
	;
	if v788&int32(1) == int32(0) {
		v817 = v788
		goto L154
	} else {
		goto L155
	}
L154:
	;
	if v788 == int32(1) {
		goto L152
	} else {
		goto L156
	}
L155:
	;
	v795 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[6])))
	v797 = v788 + int32(-1)
	v800 = v795 + v797<<(uint(int32(2))%32)
	v801 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v800))))
	v807 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28|int32(base.Ui32(v801)>>(uint(int32(3))%32))&int32(12)))))
	v814 = v807<<(uint(int32(5))%32)&int32(96) | v801&int32(159)
	*(*uint8)(unsafe.Add(mBase, uint32(v800))) = uint8(v814)
	v817 = v797
	goto L154
L156:
	;
	v836 = v817 + int32(1)
	v838 = v817<<(uint(int32(2))%32) + int32(-4)
	goto L157
L157:
	;
	v851 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[6])))
	v852 = v851 + v838
	v853 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v852))))
	v854 = int32(3)
	v856 = int32(12)
	v859 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28|int32(base.Ui32(v853)>>(uint(v854)%32))&v856))))
	v860 = int32(5)
	v862 = int32(96)
	v864 = int32(159)
	v866 = v859<<(uint(v860)%32)&v862 | v853&v864
	*(*uint8)(unsafe.Add(mBase, uint32(v852))) = uint8(v866)
	v868 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[6])))
	v871 = v868 + v838 + int32(-4)
	v872 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v871))))
	v878 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28|int32(base.Ui32(v872)>>(uint(v854)%32))&v856))))
	v885 = v878<<(uint(v860)%32)&v862 | v872&v864
	*(*uint8)(unsafe.Add(mBase, uint32(v871))) = uint8(v885)
	v890 = v836 + int32(-2)
	if base.Ui32(int32(1)) < base.Ui32(v890) {
		v836 = v890
		v838 = v838 + int32(-8)
		goto L157
	} else {
		goto L159
	}
L158:
	;
	goto L152
L159:
	;
	goto L158
L160:
	;
	if base.Ui32(int32(-4)) < base.Ui32(v765-v591) {
		goto L115
	} else {
		goto L165
	}
L161:
	;
	v940 = v926
	v944 = v921
	goto L162
L162:
	;
	v955 = int32(744)
	v956 = F_memcpy(m, v944, v923, v955)
	mBase = m.M
	v960 = v940 + int32(-1)
	if v960 != 0 {
		v940 = v960
		v944 = v956 + v955
		goto L162
	} else {
		goto L164
	}
L163:
	;
	v967 = v765 + v926
	goto L160
L164:
	;
	goto L163
L165:
	;
	v1003 = v342 + v967*int32(744)
	v1007 = v591 - v967
	goto L166
L166:
	;
	v1018 = int32(744)
	v1019 = F_memcpy(m, v1003, v923, v1018)
	mBase = m.M
	v1023 = F_memcpy(m, v1019+v1018, v923, v1018)
	mBase = m.M
	v1027 = F_memcpy(m, v1019+int32(1488), v923, v1018)
	mBase = m.M
	v1031 = F_memcpy(m, v1019+int32(2232), v923, v1018)
	mBase = m.M
	v1035 = v1007 + int32(-4)
	if v1035 != 0 {
		v1003 = v1019 + int32(2976)
		v1007 = v1035
		goto L166
	} else {
		goto L168
	}
L167:
	;
	goto L115
L168:
	;
	goto L167
L169:
	;
	v1067 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if v1067 < int32(1) {
		goto L171
	} else {
		goto L172
	}
L170:
	;
	v1064 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1065 = *(*int32)(unsafe.Add(mBase, uint32(v1064)+28))
	v1066 = v1065
	goto L169
L171:
	;
	m.G0 = v28 + int32(16)
	return
L172:
	;
	v1076 = v1067
	v1079 = v342
	goto L173
L173:
	;
	v1097 = *(*int32)(unsafe.Add(mBase, uint32(v1079+int32(680))))
	v1100 = m.G1
	v1102 = v1100 + int32(_a_F_VP8SetSegmentParams_1)
	v1103 = int32(127)
	if v1097 < v1103 {
		goto L175
	} else {
		goto L176
	}
L174:
	;
	goto L171
L175:
	;
	v1106 = v1097
	goto L177
L176:
	;
	v1106 = v1103
	goto L177
L177:
	;
	v1107 = int32(0)
	if v1107 < v1106 {
		goto L178
	} else {
		goto L179
	}
L178:
	;
	v1110 = v1106
	goto L180
L179:
	;
	v1110 = v1107
	goto L180
L180:
	;
	v1114 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1102+v1110<<(uint(int32(1))%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(2)))) = uint16(v1114)
	v1117 = v1100 + int32(_a_F_VP8SetSegmentParams_3)
	v1118 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3396))
	v1119 = v1097 + v1118
	v1120 = int32(127)
	if v1119 < v1120 {
		goto L181
	} else {
		goto L182
	}
L181:
	;
	v1123 = v1119
	goto L183
L182:
	;
	v1123 = v1120
	goto L183
L183:
	;
	v1124 = int32(0)
	if v1124 < v1123 {
		goto L184
	} else {
		goto L185
	}
L184:
	;
	v1127 = v1123
	goto L186
L185:
	;
	v1127 = v1124
	goto L186
L186:
	;
	v1129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1117+v1127))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079))) = uint16(v1129)
	v1132 = v1079 + int32(448)
	v1133 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3408))
	v1134 = v1097 + v1133
	v1135 = int32(117)
	if v1134 < v1135 {
		goto L187
	} else {
		goto L188
	}
L187:
	;
	v1138 = v1134
	goto L189
L188:
	;
	v1138 = v1135
	goto L189
L189:
	;
	v1139 = int32(0)
	if v1139 < v1138 {
		goto L190
	} else {
		goto L191
	}
L190:
	;
	v1142 = v1138
	goto L192
L191:
	;
	v1142 = v1139
	goto L192
L192:
	;
	v1144 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1117+v1142))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1132))) = uint16(v1144)
	v1147 = v1079 + int32(224)
	v1148 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3400))
	v1149 = v1097 + v1148
	v1150 = int32(127)
	if v1149 < v1150 {
		goto L193
	} else {
		goto L194
	}
L193:
	;
	v1153 = v1149
	goto L195
L194:
	;
	v1153 = v1150
	goto L195
L195:
	;
	v1154 = int32(0)
	if v1154 < v1153 {
		goto L196
	} else {
		goto L197
	}
L196:
	;
	v1157 = v1153
	goto L198
L197:
	;
	v1157 = v1154
	goto L198
L198:
	;
	v1159 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1117+v1157))))
	v1161 = v1159 << (uint(int32(1)) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147))) = uint16(v1161)
	v1167 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3404))
	v1168 = v1097 + v1167
	v1169 = int32(127)
	if v1168 < v1169 {
		goto L199
	} else {
		goto L200
	}
L199:
	;
	v1172 = v1168
	goto L201
L200:
	;
	v1172 = v1169
	goto L201
L201:
	;
	v1173 = int32(0)
	if v1173 < v1172 {
		goto L202
	} else {
		goto L203
	}
L202:
	;
	v1176 = v1172
	goto L204
L203:
	;
	v1176 = v1173
	goto L204
L204:
	;
	v1180 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1100+int32(_a_F_VP8SetSegmentParams_4)+v1176<<(uint(int32(1))%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(226)))) = uint16(v1180)
	v1184 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3412))
	v1185 = v1097 + v1184
	v1186 = int32(127)
	if v1185 < v1186 {
		goto L205
	} else {
		goto L206
	}
L205:
	;
	v1189 = v1185
	goto L207
L206:
	;
	v1189 = v1186
	goto L207
L207:
	;
	v1190 = int32(0)
	if v1190 < v1189 {
		goto L208
	} else {
		goto L209
	}
L208:
	;
	v1193 = v1189
	goto L210
L209:
	;
	v1193 = v1190
	goto L210
L210:
	;
	v1197 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1102+v1193<<(uint(int32(1))%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(450)))) = uint16(v1197)
	v1199 = int32(0)
	v1205 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+4)) = uint16(v1205)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+6)) = uint16(v1205)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+8)) = uint16(v1205)
	v1209 = int32(131072)
	v1210 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079))))
	v1211 = base.I32_div_u_s(v1209, v1210)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+32)) = uint16(v1211)
	v1214 = base.I32_div_u_s(v1209, v1205)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+34)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+36)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+38)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+40)) = uint16(v1214)
	v1219 = m.G1
	v1225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1219)+uint32(_c_F_VP8SetSegmentParams[7]))))
	v1226 = int32(9)
	v1227 = v1225 << (uint(v1226) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+64)) = v1227
	v1229 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1219)+uint32(_c_F_VP8SetSegmentParams[8]))))
	v1231 = v1229 << (uint(v1226) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+68)) = v1231
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+72)) = v1231
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+76)) = v1231
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+80)) = v1231
	v1236 = int32(131071)
	v1238 = int32(_a_F_VP8SetSegmentParams_5)
	v1240 = base.I32_div_u_s(v1227^v1236, v1211&v1238)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+128)) = v1240
	v1246 = base.I32_div_u_s(v1231^v1236, v1214&v1238)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+132)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+136)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+140)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+144)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+148)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+84)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+42)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+10)) = uint16(v1205)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+12)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+152)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+88)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+44)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+14)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+156)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+92)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+46)) = uint16(v1214)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+160)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+96)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+48)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+16)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+164)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+100)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+50)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+18)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+168)) = v1246
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+104)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+52)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+20)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+108)) = v1231
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+172)) = v1246
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+54)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+22)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+112)) = v1231
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+176)) = v1246
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+56)) = uint16(v1214)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+24)) = uint16(v1205)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+26)) = uint16(v1205)
	v1285 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+58)) = uint16(v1285)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+116)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+60)) = uint16(v1285)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+28)) = uint16(v1205)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+120)) = v1231
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+30)) = uint16(v1205)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079)+62)) = uint16(v1285)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(192)))) = uint16(v1199)
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+124)) = v1231
	v1299 = *(*int32)(unsafe.Add(mBase, uint32(v1079)+132))
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+180)) = v1299
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+184)) = v1299
	*(*int32)(unsafe.Add(mBase, uint32(v1079)+188)) = v1299
	v1303 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079))))
	goto L213
L211:
	;
	v1441 = int32(0)
	v1446 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+4)) = uint16(v1446)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+6)) = uint16(v1446)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+8)) = uint16(v1446)
	v1450 = int32(131072)
	v1451 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147))))
	v1452 = base.I32_div_u_s(v1450, v1451)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+32)) = uint16(v1452)
	v1455 = base.I32_div_u_s(v1450, v1446)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+34)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+36)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+38)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+40)) = uint16(v1455)
	v1460 = m.G1
	v1466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1460)+uint32(_c_F_VP8SetSegmentParams[9]))))
	v1467 = int32(9)
	v1468 = v1466 << (uint(v1467) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+64)) = v1468
	v1470 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1460)+uint32(_c_F_VP8SetSegmentParams[10]))))
	v1472 = v1470 << (uint(v1467) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+68)) = v1472
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+72)) = v1472
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+76)) = v1472
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+80)) = v1472
	v1477 = int32(131071)
	v1479 = int32(_a_F_VP8SetSegmentParams_5)
	v1481 = base.I32_div_u_s(v1468^v1477, v1452&v1479)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+128)) = v1481
	v1487 = base.I32_div_u_s(v1472^v1477, v1455&v1479)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+132)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+136)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+140)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+144)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+148)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+84)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+42)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+10)) = uint16(v1446)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+12)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+152)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+88)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+44)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+14)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+156)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+92)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+46)) = uint16(v1455)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+160)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+96)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+48)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+16)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+164)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+100)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+50)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+18)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+168)) = v1487
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+104)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+52)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+20)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+108)) = v1472
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+172)) = v1487
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+54)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+22)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+112)) = v1472
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+176)) = v1487
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+56)) = uint16(v1455)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+24)) = uint16(v1446)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+26)) = uint16(v1446)
	v1526 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+58)) = uint16(v1526)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+116)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+60)) = uint16(v1526)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+28)) = uint16(v1446)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+120)) = v1472
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+30)) = uint16(v1446)
	*(*uint16)(unsafe.Add(mBase, uint32(v1147)+62)) = uint16(v1526)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(416)))) = uint16(v1441)
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+124)) = v1472
	v1540 = *(*int32)(unsafe.Add(mBase, uint32(v1147)+132))
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+180)) = v1540
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+184)) = v1540
	*(*int32)(unsafe.Add(mBase, uint32(v1147)+188)) = v1540
	v1544 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147))))
	goto L218
L212:
	;
	v1439 = int32(base.Ui32(v1411+(v1403+(v1395+v1418))+v1418+v1205<<(uint(int32(2))%32)+((v1205+v1205)<<(uint(v1417)%32)+v1303)+int32(8)) >> (uint(int32(4)) % 32))
	goto L211
L213:
	;
	v1343 = int32(30)
	v1345 = int32(11)
	v1346 = int32(base.Ui32(v1205*v1343) >> (uint(v1345) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(194)))) = uint16(v1346)
	v1350 = int32(90)
	v1353 = int32(base.Ui32(v1205*v1350) >> (uint(v1345) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(222)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(220)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(218)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(216)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(214)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(212)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(210)))) = uint16(v1353)
	v1379 = int32(60)
	v1382 = int32(base.Ui32(v1205*v1379) >> (uint(v1345) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(208)))) = uint16(v1382)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(206)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(204)))) = uint16(v1353)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(202)))) = uint16(v1382)
	v1395 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079)+4)))
	v1399 = int32(base.Ui32(v1395*v1379) >> (uint(v1345) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(196)))) = uint16(v1399)
	v1403 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079)+6)))
	v1407 = int32(base.Ui32(v1403*v1350) >> (uint(v1345) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(198)))) = uint16(v1407)
	v1411 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079)+8)))
	v1415 = int32(base.Ui32(v1411*v1343) >> (uint(v1345) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(200)))) = uint16(v1415)
	v1417 = int32(1)
	v1418 = v1205 << (uint(v1417) % 32)
	goto L212
L215:
	;
	v1682 = int32(0)
	v1687 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+4)) = uint16(v1687)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+6)) = uint16(v1687)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+8)) = uint16(v1687)
	v1691 = int32(131072)
	v1692 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132))))
	v1693 = base.I32_div_u_s(v1691, v1692)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+32)) = uint16(v1693)
	v1696 = base.I32_div_u_s(v1691, v1687)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+34)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+36)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+38)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+40)) = uint16(v1696)
	v1701 = m.G1
	v1707 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1701)+uint32(_c_F_VP8SetSegmentParams[11]))))
	v1708 = int32(9)
	v1709 = v1707 << (uint(v1708) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+64)) = v1709
	v1711 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1701)+uint32(_c_F_VP8SetSegmentParams[12]))))
	v1713 = v1711 << (uint(v1708) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+68)) = v1713
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+72)) = v1713
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+76)) = v1713
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+80)) = v1713
	v1718 = int32(131071)
	v1720 = int32(_a_F_VP8SetSegmentParams_5)
	v1722 = base.I32_div_u_s(v1709^v1718, v1693&v1720)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+128)) = v1722
	v1728 = base.I32_div_u_s(v1713^v1718, v1696&v1720)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+132)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+136)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+140)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+144)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+148)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+84)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+42)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+10)) = uint16(v1687)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+12)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+152)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+88)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+44)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+14)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+156)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+92)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+46)) = uint16(v1696)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+160)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+96)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+48)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+16)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+164)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+100)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+50)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+18)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+168)) = v1728
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+104)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+52)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+20)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+108)) = v1713
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+172)) = v1728
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+54)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+22)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+112)) = v1713
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+176)) = v1728
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+56)) = uint16(v1696)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+24)) = uint16(v1687)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+26)) = uint16(v1687)
	v1767 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+58)) = uint16(v1767)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+116)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+60)) = uint16(v1767)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+28)) = uint16(v1687)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+120)) = v1713
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+30)) = uint16(v1687)
	*(*uint16)(unsafe.Add(mBase, uint32(v1132)+62)) = uint16(v1767)
	*(*uint16)(unsafe.Add(mBase, uint32(v1079+int32(640)))) = uint16(v1682)
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+124)) = v1713
	v1781 = *(*int32)(unsafe.Add(mBase, uint32(v1132)+132))
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+180)) = v1781
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+184)) = v1781
	*(*int32)(unsafe.Add(mBase, uint32(v1132)+188)) = v1781
	v1785 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132))))
	goto L222
L216:
	;
	v1680 = int32(base.Ui32(v1564+v1565+v1567+v1569+v1564+v1446<<(uint(int32(2))%32)+((v1446+v1446)<<(uint(v1563)%32)+v1544)+int32(8)) >> (uint(int32(4)) % 32))
	goto L215
L218:
	;
	v1549 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(418)))) = v1549
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(426)))) = v1549
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(434)))) = v1549
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(440)))) = v1549
	v1563 = int32(1)
	v1564 = v1446 << (uint(v1563) % 32)
	v1565 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147)+4)))
	v1567 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147)+6)))
	v1569 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1147)+8)))
	goto L216
L219:
	;
	v1923 = v1079 + int32(708)
	v1924 = v1439 * v1439
	v1925 = int32(7)
	*(*int32)(unsafe.Add(mBase, uint32(v1923))) = int32(base.Ui32(v1924) >> (uint(v1925) % 32))
	v1929 = v1079 + int32(716)
	v1932 = v1439 * v1066 >> (uint(int32(5)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1929))) = v1932
	v1935 = v1079 + int32(724)
	v1938 = int32(3)
	v1939 = v1924 * v1925 >> (uint(v1938) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1935))) = v1939
	v1942 = v1079 + int32(696)
	v1943 = v1680 * v1680
	v1945 = v1943 * v1938
	*(*int32)(unsafe.Add(mBase, uint32(v1942))) = v1945
	v1948 = v1079 + int32(700)
	v1952 = v1924 * v1938 >> (uint(v1925) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1948))) = v1952
	v1955 = v1079 + int32(720)
	*(*int32)(unsafe.Add(mBase, uint32(v1955))) = int32(base.Ui32(v1943) >> (uint(int32(2)) % 32))
	v1960 = v1079 + int32(728)
	v1961 = v1921 * v1921
	v1963 = v1961 << (uint(int32(1)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1960))) = v1963
	v1966 = v1079 + int32(704)
	v1970 = v1961 * v1938 >> (uint(int32(6)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1966))) = v1970
	if int32(0) < v1952 {
		goto L223
	} else {
		goto L224
	}
L220:
	;
	v1921 = int32(base.Ui32(v1805+v1806+v1808+v1810+v1805+v1687<<(uint(int32(2))%32)+((v1687+v1687)<<(uint(v1804)%32)+v1785)+int32(8)) >> (uint(int32(4)) % 32))
	goto L219
L222:
	;
	v1790 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(642)))) = v1790
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(650)))) = v1790
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(658)))) = v1790
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(664)))) = v1790
	v1804 = int32(1)
	v1805 = v1687 << (uint(v1804) % 32)
	v1806 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132)+4)))
	v1808 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132)+6)))
	v1810 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1132)+8)))
	goto L220
L223:
	;
	if int32(0) < v1945 {
		goto L225
	} else {
		goto L226
	}
L224:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1948))) = int32(1)
	goto L223
L225:
	;
	if int32(0) < v1970 {
		goto L227
	} else {
		goto L228
	}
L226:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1942))) = int32(1)
	goto L225
L227:
	;
	if base.Ui32(int32(127)) < base.Ui32(v1924) {
		goto L229
	} else {
		goto L230
	}
L228:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1966))) = int32(1)
	goto L227
L229:
	;
	if int32(0) < v1939 {
		goto L231
	} else {
		goto L232
	}
L230:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1923))) = int32(1)
	goto L229
L231:
	;
	if base.Ui32(int32(3)) < base.Ui32(v1943) {
		goto L233
	} else {
		goto L234
	}
L232:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1935))) = int32(1)
	goto L231
L233:
	;
	if int32(0) < v1963 {
		goto L235
	} else {
		goto L236
	}
L234:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1955))) = int32(1)
	goto L233
L235:
	;
	if int32(0) < v1932 {
		goto L237
	} else {
		goto L238
	}
L236:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1960))) = int32(1)
	goto L235
L237:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1079+int32(688)))) = int32(0)
	v2010 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1079))))
	*(*int32)(unsafe.Add(mBase, uint32(v1079+int32(692)))) = v2010 * int32(20)
	*(*int64)(unsafe.Add(mBase, uint32(v1079+int32(736)))) = base.I64_extend_i32_s(v1924 * int32(1000))
	v2023 = v1076 + int32(-1)
	if v2023 != 0 {
		v1076 = v2023
		v1079 = v1079 + int32(744)
		goto L173
	} else {
		goto L239
	}
L238:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1929))) = int32(1)
	goto L237
L239:
	;
	goto L174
}
func F_VP8SetSkip(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	var v11 int32
	_ = v11
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v4 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3))))
	v11 = v4&int32(239) | l1<<(uint(int32(4))%32)&int32(16)
	*(*uint8)(unsafe.Add(mBase, uint32(v3))) = uint8(v11)
	return
}
func F_VP8StoreFilterStats(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 float64
	_ = v42
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v47 float64
	_ = v47
	var v54 int32
	_ = v54
	var v59 int32
	_ = v59
	var v69 int32
	_ = v69
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v80 int32
	_ = v80
	var v90 int32
	_ = v90
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v156 int32
	_ = v156
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v163 float64
	_ = v163
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v166 float64
	_ = v166
	var v178 int32
	_ = v178
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v18 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v17))))
	if v18&int32(19) == int32(17) {
	} else {
		v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
		if v23 == int32(0) {
		} else {
			v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
			v30 = int32(base.Ui32(v18)>>(uint(int32(5))%32)) & int32(3)
			v33 = v26 + v30*int32(744)
			v36 = *(*int32)(unsafe.Add(mBase, uint32(v33+int32(1088))))
			v39 = *(*int32)(unsafe.Add(mBase, uint32(v33+int32(1092))))
			v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v42 = F_GetMBSSIM(m, v40, v41)
			mBase = m.M
			v43 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
			v45 = v30 << (uint(int32(9)) % 32)
			v46 = v43 + v45
			v47 = *(*float64)(unsafe.Add(mBase, uint32(v46)))
			*(*float64)(unsafe.Add(mBase, uint32(v46))) = base.F64_add(v42, v47)
			if v36 < int32(0) {
			} else {
				v54 = int32(1)
				v59 = int32(3)
				if v54 < v36 {
					v69 = int32(4)
				} else {
					v69 = v54
				}
				v75 = int32(0) - v36
				v76 = v39<<(uint(v54)%32) - v36<<(uint(v54)%32)
				v80 = v45 + v39<<(uint(v59)%32) - v36<<(uint(v59)%32)
				for {
					v90 = v39 + v75
					if base.Ui32(v90+int32(-64)) < base.Ui32(int32(-63)) {
					} else {
						v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
						v96 = *(*int32)(unsafe.Add(mBase, uint32(v95)))
						v97 = *(*int32)(unsafe.Add(mBase, uint32(v96)+36))
						if v97 < int32(1) {
							v110 = v90
						} else {
							if base.Ui32(int32(4)) < base.Ui32(v97) {
								v104 = int32(2)
							} else {
								v104 = int32(1)
							}
							v105 = int32(base.Ui32(v90) >> (uint(v104) % 32))
							v107 = int32(9) - v97
							if v105 < v107 {
								v109 = v105
							} else {
								v109 = v107
							}
							v110 = v109
						}
						v112 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						v115 = F_memcpy(m, v112, v113, int32(512))
						mBase = m.M
						v116 = int32(1)
						if v116 < v110 {
							v119 = v110
						} else {
							v119 = v116
						}
						v120 = v119 + v76
						v121 = *(*int32)(unsafe.Add(mBase, uint32(v95)+8))
						if v121 != int32(1) {
							if base.Ui32(int32(39)) < base.Ui32(v90) {
								v138 = int32(2)
							} else {
								v138 = base.B2i32(base.Ui32(int32(14)) < base.Ui32(v90))
							}
							v139 = m.G18
							v140 = *(*int32)(unsafe.Add(mBase, uint32(v139)))
							m.T0[v140].(func(*base.Module, int32, int32, int32, int32, int32))(m, v115, int32(32), v120, v119, v138)
							mBase = m.M
							v143 = v115 + int32(16)
							v145 = v115 + int32(24)
							v146 = int32(32)
							v147 = m.G19
							v148 = *(*int32)(unsafe.Add(mBase, uint32(v147)))
							m.T0[v148].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v143, v145, v146, v120, v119, v138)
							mBase = m.M
							v151 = m.G20
							v152 = *(*int32)(unsafe.Add(mBase, uint32(v151)))
							m.T0[v152].(func(*base.Module, int32, int32, int32, int32, int32))(m, v115, v146, v120, v119, v138)
							mBase = m.M
							v155 = m.G21
							v156 = *(*int32)(unsafe.Add(mBase, uint32(v155)))
							m.T0[v156].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v143, v145, v146, v120, v119, v138)
							mBase = m.M
						} else {
							v124 = int32(32)
							v125 = m.G16
							v126 = *(*int32)(unsafe.Add(mBase, uint32(v125)))
							m.T0[v126].(func(*base.Module, int32, int32, int32))(m, v115, v124, v120)
							mBase = m.M
							v129 = m.G17
							v130 = *(*int32)(unsafe.Add(mBase, uint32(v129)))
							m.T0[v130].(func(*base.Module, int32, int32, int32))(m, v115, v124, v120)
							mBase = m.M
						}
						v161 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v162 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						v163 = F_GetMBSSIM(m, v161, v162)
						mBase = m.M
						v164 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
						v165 = v164 + v80
						v166 = *(*float64)(unsafe.Add(mBase, uint32(v165)))
						*(*float64)(unsafe.Add(mBase, uint32(v165))) = base.F64_add(v163, v166)
					}
					v178 = v75 + v69
					if v178 <= v36 {
						v75 = v178
						v76 = v76 + v69<<(uint(int32(1))%32)
						v80 = v80 + v69<<(uint(int32(3))%32)
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
func F_VP8TBufferClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v27 int32
	_ = v27
	if l0 == int32(0) {
	} else {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v6 == int32(0) {
		} else {
			v10 = v6
			for {
				v12 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
				F_free(m, v10)
				mBase = m.M
				if v12 != 0 {
					v10 = v12
					continue
				} else {
					break
				}
				break
			}
		}
		v17 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v17
		*(*int64)(unsafe.Add(mBase, uint32(l0)+8)) = int64(0)
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v17
		v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		v24 = int32(_a_F_VP8TBufferClear_0)
		if v24 < v23 {
			v27 = v23
		} else {
			v27 = v24
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v27
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = l0
	}
	return
}
func F_VP8TBufferInit(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	v3 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v3
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v3
	*(*int64)(unsafe.Add(mBase, uint32(l0)+8)) = int64(0)
	v9 = int32(_a_F_VP8TBufferInit_0)
	if v9 < l1 {
		v12 = l1
	} else {
		v12 = v9
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v12
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = l0
	return
}
func F_VP8WriteProbas(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v52 int32
	_ = v52
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v89 int32
	_ = v89
	var v93 int32
	_ = v93
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v103 int32
	_ = v103
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v154 int32
	_ = v154
	var v158 int32
	_ = v158
	var v164 int32
	_ = v164
	var v168 int32
	_ = v168
	var v171 int32
	_ = v171
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v194 int32
	_ = v194
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v206 int32
	_ = v206
	var v208 int32
	_ = v208
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v236 int32
	_ = v236
	var v239 int32
	_ = v239
	var v241 int32
	_ = v241
	var v245 int32
	_ = v245
	var v247 int32
	_ = v247
	var v251 int32
	_ = v251
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v267 int32
	_ = v267
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v280 int32
	_ = v280
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v303 int32
	_ = v303
	var v307 int32
	_ = v307
	var v311 int32
	_ = v311
	var v315 int32
	_ = v315
	var v317 int32
	_ = v317
	var v321 int32
	_ = v321
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v345 int32
	_ = v345
	var v348 int32
	_ = v348
	var v350 int32
	_ = v350
	var v354 int32
	_ = v354
	var v356 int32
	_ = v356
	var v360 int32
	_ = v360
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v368 int32
	_ = v368
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v382 int32
	_ = v382
	var v386 int32
	_ = v386
	var v389 int32
	_ = v389
	var v393 int32
	_ = v393
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v402 int32
	_ = v402
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v410 int32
	_ = v410
	var v412 int32
	_ = v412
	var v416 int32
	_ = v416
	var v420 int32
	_ = v420
	var v424 int32
	_ = v424
	var v426 int32
	_ = v426
	var v430 int32
	_ = v430
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v454 int32
	_ = v454
	var v457 int32
	_ = v457
	var v459 int32
	_ = v459
	var v463 int32
	_ = v463
	var v465 int32
	_ = v465
	var v469 int32
	_ = v469
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v477 int32
	_ = v477
	var v481 int32
	_ = v481
	var v485 int32
	_ = v485
	var v491 int32
	_ = v491
	var v495 int32
	_ = v495
	var v498 int32
	_ = v498
	var v502 int32
	_ = v502
	var v504 int32
	_ = v504
	var v505 int32
	_ = v505
	var v511 int32
	_ = v511
	var v513 int32
	_ = v513
	var v516 int32
	_ = v516
	var v519 int32
	_ = v519
	var v521 int32
	_ = v521
	var v525 int32
	_ = v525
	var v529 int32
	_ = v529
	var v533 int32
	_ = v533
	var v535 int32
	_ = v535
	var v539 int32
	_ = v539
	var v542 int32
	_ = v542
	var v543 int32
	_ = v543
	var v563 int32
	_ = v563
	var v566 int32
	_ = v566
	var v568 int32
	_ = v568
	var v572 int32
	_ = v572
	var v574 int32
	_ = v574
	var v578 int32
	_ = v578
	var v582 int32
	_ = v582
	var v583 int32
	_ = v583
	var v586 int32
	_ = v586
	var v590 int32
	_ = v590
	var v594 int32
	_ = v594
	var v600 int32
	_ = v600
	var v604 int32
	_ = v604
	var v607 int32
	_ = v607
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v614 int32
	_ = v614
	var v620 int32
	_ = v620
	var v622 int32
	_ = v622
	var v625 int32
	_ = v625
	var v628 int32
	_ = v628
	var v630 int32
	_ = v630
	var v634 int32
	_ = v634
	var v638 int32
	_ = v638
	var v642 int32
	_ = v642
	var v644 int32
	_ = v644
	var v648 int32
	_ = v648
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v672 int32
	_ = v672
	var v675 int32
	_ = v675
	var v677 int32
	_ = v677
	var v681 int32
	_ = v681
	var v683 int32
	_ = v683
	var v687 int32
	_ = v687
	var v691 int32
	_ = v691
	var v692 int32
	_ = v692
	var v695 int32
	_ = v695
	var v699 int32
	_ = v699
	var v703 int32
	_ = v703
	var v709 int32
	_ = v709
	var v713 int32
	_ = v713
	var v716 int32
	_ = v716
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v729 int32
	_ = v729
	var v731 int32
	_ = v731
	var v734 int32
	_ = v734
	var v737 int32
	_ = v737
	var v739 int32
	_ = v739
	var v743 int32
	_ = v743
	var v747 int32
	_ = v747
	var v751 int32
	_ = v751
	var v753 int32
	_ = v753
	var v757 int32
	_ = v757
	var v760 int32
	_ = v760
	var v761 int32
	_ = v761
	var v781 int32
	_ = v781
	var v784 int32
	_ = v784
	var v786 int32
	_ = v786
	var v790 int32
	_ = v790
	var v792 int32
	_ = v792
	var v796 int32
	_ = v796
	var v800 int32
	_ = v800
	var v801 int32
	_ = v801
	var v804 int32
	_ = v804
	var v808 int32
	_ = v808
	var v812 int32
	_ = v812
	var v818 int32
	_ = v818
	var v822 int32
	_ = v822
	var v825 int32
	_ = v825
	var v829 int32
	_ = v829
	var v831 int32
	_ = v831
	var v832 int32
	_ = v832
	var v838 int32
	_ = v838
	var v840 int32
	_ = v840
	var v843 int32
	_ = v843
	var v846 int32
	_ = v846
	var v848 int32
	_ = v848
	var v852 int32
	_ = v852
	var v856 int32
	_ = v856
	var v860 int32
	_ = v860
	var v862 int32
	_ = v862
	var v866 int32
	_ = v866
	var v869 int32
	_ = v869
	var v870 int32
	_ = v870
	var v890 int32
	_ = v890
	var v893 int32
	_ = v893
	var v895 int32
	_ = v895
	var v899 int32
	_ = v899
	var v901 int32
	_ = v901
	var v905 int32
	_ = v905
	var v909 int32
	_ = v909
	var v910 int32
	_ = v910
	var v913 int32
	_ = v913
	var v917 int32
	_ = v917
	var v921 int32
	_ = v921
	var v927 int32
	_ = v927
	var v931 int32
	_ = v931
	var v934 int32
	_ = v934
	var v938 int32
	_ = v938
	var v940 int32
	_ = v940
	var v941 int32
	_ = v941
	var v947 int32
	_ = v947
	var v949 int32
	_ = v949
	var v952 int32
	_ = v952
	var v955 int32
	_ = v955
	var v957 int32
	_ = v957
	var v961 int32
	_ = v961
	var v965 int32
	_ = v965
	var v969 int32
	_ = v969
	var v971 int32
	_ = v971
	var v975 int32
	_ = v975
	var v978 int32
	_ = v978
	var v979 int32
	_ = v979
	var v999 int32
	_ = v999
	var v1002 int32
	_ = v1002
	var v1004 int32
	_ = v1004
	var v1008 int32
	_ = v1008
	var v1010 int32
	_ = v1010
	var v1014 int32
	_ = v1014
	var v1018 int32
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1022 int32
	_ = v1022
	var v1026 int32
	_ = v1026
	var v1030 int32
	_ = v1030
	var v1036 int32
	_ = v1036
	var v1040 int32
	_ = v1040
	var v1043 int32
	_ = v1043
	var v1047 int32
	_ = v1047
	var v1049 int32
	_ = v1049
	var v1050 int32
	_ = v1050
	var v1056 int32
	_ = v1056
	var v1058 int32
	_ = v1058
	var v1061 int32
	_ = v1061
	var v1064 int32
	_ = v1064
	var v1066 int32
	_ = v1066
	var v1070 int32
	_ = v1070
	var v1074 int32
	_ = v1074
	var v1078 int32
	_ = v1078
	var v1080 int32
	_ = v1080
	var v1084 int32
	_ = v1084
	var v1087 int32
	_ = v1087
	var v1088 int32
	_ = v1088
	var v1108 int32
	_ = v1108
	var v1111 int32
	_ = v1111
	var v1113 int32
	_ = v1113
	var v1117 int32
	_ = v1117
	var v1119 int32
	_ = v1119
	var v1123 int32
	_ = v1123
	var v1127 int32
	_ = v1127
	var v1128 int32
	_ = v1128
	var v1131 int32
	_ = v1131
	var v1135 int32
	_ = v1135
	var v1139 int32
	_ = v1139
	var v1145 int32
	_ = v1145
	var v1149 int32
	_ = v1149
	var v1152 int32
	_ = v1152
	var v1156 int32
	_ = v1156
	var v1158 int32
	_ = v1158
	var v1159 int32
	_ = v1159
	var v1165 int32
	_ = v1165
	var v1167 int32
	_ = v1167
	var v1170 int32
	_ = v1170
	var v1173 int32
	_ = v1173
	var v1175 int32
	_ = v1175
	var v1179 int32
	_ = v1179
	var v1183 int32
	_ = v1183
	var v1187 int32
	_ = v1187
	var v1189 int32
	_ = v1189
	var v1193 int32
	_ = v1193
	var v1196 int32
	_ = v1196
	var v1197 int32
	_ = v1197
	var v1217 int32
	_ = v1217
	var v1220 int32
	_ = v1220
	var v1222 int32
	_ = v1222
	var v1226 int32
	_ = v1226
	var v1228 int32
	_ = v1228
	var v1232 int32
	_ = v1232
	var v1236 int32
	_ = v1236
	var v1237 int32
	_ = v1237
	var v1240 int32
	_ = v1240
	var v1244 int32
	_ = v1244
	var v1248 int32
	_ = v1248
	var v1254 int32
	_ = v1254
	var v1259 int32
	_ = v1259
	var v1260 int32
	_ = v1260
	var v1265 int32
	_ = v1265
	var v1268 int32
	_ = v1268
	var v1273 int32
	_ = v1273
	var v1276 int32
	_ = v1276
	var v1279 int32
	_ = v1279
	var v1281 int32
	_ = v1281
	var v1284 int32
	_ = v1284
	var v1286 int32
	_ = v1286
	var v1290 int32
	_ = v1290
	var v1294 int32
	_ = v1294
	var v1295 int32
	_ = v1295
	var v1298 int32
	_ = v1298
	var v1302 int32
	_ = v1302
	var v1306 int32
	_ = v1306
	var v1314 int32
	_ = v1314
	var v1326 int32
	_ = v1326
	var v1329 int32
	_ = v1329
	var v1331 int32
	_ = v1331
	var v1335 int32
	_ = v1335
	var v1337 int32
	_ = v1337
	var v1341 int32
	_ = v1341
	var v1345 int32
	_ = v1345
	var v1346 int32
	_ = v1346
	var v1349 int32
	_ = v1349
	var v1353 int32
	_ = v1353
	var v1357 int32
	_ = v1357
	var v1363 int32
	_ = v1363
	v3 = int32(0)
	v18 = v3
	v19 = l1
	v20 = v3
	for {
		v35 = v19
		v36 = v18
		v37 = int32(0)
		for {
			v52 = int32(-33)
			for {
				v57 = m.G1
				v58 = v35 + v52
				v61 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(37)))))
				v64 = v36 + v52
				v66 = int32(33)
				v68 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(_a_F_VP8WriteProbas_0)+v64+v66))))
				v69 = base.B2i32(v61 != v68)
				v75 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v57+int32(_a_F_VP8WriteProbas_1)+v64+v66))))
				v77 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v80 = v77 * v75 >> (uint(int32(8)) % 32)
				if v69 == int32(0) {
					v89 = v80
				} else {
					v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v85 = v80 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v83 + v85
					v89 = v77 - v85
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v89
				if int32(126) < v89 {
				} else {
					v93 = m.G1
					v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v93+int32(_a_F_VP8WriteProbas_2)+v89))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v97
					v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v103 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v93+int32(_a_F_VP8WriteProbas_3)+v89))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v99 << (uint(v103) % 32)
					v106 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v107 = v103 + v106
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v107
					if v107 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v69 == int32(0) {
				} else {
					v127 = int32(128)
					for {
						v130 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v132 = v130 >> (uint(int32(1)) % 32)
						if v127&v61 == int32(0) {
							v142 = v132
						} else {
							v136 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v138 = v132 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v136 + v138
							v142 = v130 - v138
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v142
						if int32(126) < v142 {
						} else {
							v146 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v147 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v146 << (uint(v147) % 32)
							v150 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v150 + v147
							v154 = m.G1
							v158 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v154+int32(_a_F_VP8WriteProbas_2)+v142))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v158
							if v150 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v164 = int32(1)
						if base.Ui32(v164) < base.Ui32(v127) {
							v127 = int32(base.Ui32(v127) >> (uint(v164) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v168 = m.G1
				v171 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(38)))))
				v175 = int32(34)
				v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v168+int32(_a_F_VP8WriteProbas_0)+v64+v175))))
				v178 = base.B2i32(v171 != v177)
				v184 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v168+int32(_a_F_VP8WriteProbas_1)+v64+v175))))
				v186 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v189 = v186 * v184 >> (uint(int32(8)) % 32)
				if v178 == int32(0) {
					v198 = v189
				} else {
					v192 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v194 = v189 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v192 + v194
					v198 = v186 - v194
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v198
				if int32(126) < v198 {
				} else {
					v202 = m.G1
					v206 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v202+int32(_a_F_VP8WriteProbas_2)+v198))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v206
					v208 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v212 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v202+int32(_a_F_VP8WriteProbas_3)+v198))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v208 << (uint(v212) % 32)
					v215 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v216 = v212 + v215
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v216
					if v216 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v178 == int32(0) {
				} else {
					v236 = int32(128)
					for {
						v239 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v241 = v239 >> (uint(int32(1)) % 32)
						if v236&v171 == int32(0) {
							v251 = v241
						} else {
							v245 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v247 = v241 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v245 + v247
							v251 = v239 - v247
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v251
						if int32(126) < v251 {
						} else {
							v255 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v256 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v255 << (uint(v256) % 32)
							v259 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v259 + v256
							v263 = m.G1
							v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v263+int32(_a_F_VP8WriteProbas_2)+v251))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v267
							if v259 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v273 = int32(1)
						if base.Ui32(v273) < base.Ui32(v236) {
							v236 = int32(base.Ui32(v236) >> (uint(v273) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v277 = m.G1
				v280 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(39)))))
				v284 = int32(35)
				v286 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v277+int32(_a_F_VP8WriteProbas_0)+v64+v284))))
				v287 = base.B2i32(v280 != v286)
				v293 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v277+int32(_a_F_VP8WriteProbas_1)+v64+v284))))
				v295 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v298 = v295 * v293 >> (uint(int32(8)) % 32)
				if v287 == int32(0) {
					v307 = v298
				} else {
					v301 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v303 = v298 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v301 + v303
					v307 = v295 - v303
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v307
				if int32(126) < v307 {
				} else {
					v311 = m.G1
					v315 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v311+int32(_a_F_VP8WriteProbas_2)+v307))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v315
					v317 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v321 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v311+int32(_a_F_VP8WriteProbas_3)+v307))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v317 << (uint(v321) % 32)
					v324 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v325 = v321 + v324
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v325
					if v325 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v287 == int32(0) {
				} else {
					v345 = int32(128)
					for {
						v348 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v350 = v348 >> (uint(int32(1)) % 32)
						if v345&v280 == int32(0) {
							v360 = v350
						} else {
							v354 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v356 = v350 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v354 + v356
							v360 = v348 - v356
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v360
						if int32(126) < v360 {
						} else {
							v364 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v365 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v364 << (uint(v365) % 32)
							v368 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v368 + v365
							v372 = m.G1
							v376 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v372+int32(_a_F_VP8WriteProbas_2)+v360))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v376
							if v368 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v382 = int32(1)
						if base.Ui32(v382) < base.Ui32(v345) {
							v345 = int32(base.Ui32(v345) >> (uint(v382) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v386 = m.G1
				v389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(40)))))
				v393 = int32(36)
				v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v386+int32(_a_F_VP8WriteProbas_0)+v64+v393))))
				v396 = base.B2i32(v389 != v395)
				v402 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v386+int32(_a_F_VP8WriteProbas_1)+v64+v393))))
				v404 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v407 = v404 * v402 >> (uint(int32(8)) % 32)
				if v396 == int32(0) {
					v416 = v407
				} else {
					v410 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v412 = v407 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v410 + v412
					v416 = v404 - v412
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v416
				if int32(126) < v416 {
				} else {
					v420 = m.G1
					v424 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v420+int32(_a_F_VP8WriteProbas_2)+v416))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v424
					v426 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v430 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v420+int32(_a_F_VP8WriteProbas_3)+v416))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v426 << (uint(v430) % 32)
					v433 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v434 = v430 + v433
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v434
					if v434 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v396 == int32(0) {
				} else {
					v454 = int32(128)
					for {
						v457 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v459 = v457 >> (uint(int32(1)) % 32)
						if v454&v389 == int32(0) {
							v469 = v459
						} else {
							v463 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v465 = v459 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v463 + v465
							v469 = v457 - v465
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v469
						if int32(126) < v469 {
						} else {
							v473 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v474 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v473 << (uint(v474) % 32)
							v477 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v477 + v474
							v481 = m.G1
							v485 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v481+int32(_a_F_VP8WriteProbas_2)+v469))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v485
							if v477 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v491 = int32(1)
						if base.Ui32(v491) < base.Ui32(v454) {
							v454 = int32(base.Ui32(v454) >> (uint(v491) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v495 = m.G1
				v498 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(41)))))
				v502 = int32(37)
				v504 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v495+int32(_a_F_VP8WriteProbas_0)+v64+v502))))
				v505 = base.B2i32(v498 != v504)
				v511 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v495+int32(_a_F_VP8WriteProbas_1)+v64+v502))))
				v513 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v516 = v513 * v511 >> (uint(int32(8)) % 32)
				if v505 == int32(0) {
					v525 = v516
				} else {
					v519 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v521 = v516 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v519 + v521
					v525 = v513 - v521
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v525
				if int32(126) < v525 {
				} else {
					v529 = m.G1
					v533 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v529+int32(_a_F_VP8WriteProbas_2)+v525))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v533
					v535 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v539 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v529+int32(_a_F_VP8WriteProbas_3)+v525))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v535 << (uint(v539) % 32)
					v542 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v543 = v539 + v542
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v543
					if v543 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v505 == int32(0) {
				} else {
					v563 = int32(128)
					for {
						v566 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v568 = v566 >> (uint(int32(1)) % 32)
						if v563&v498 == int32(0) {
							v578 = v568
						} else {
							v572 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v574 = v568 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v572 + v574
							v578 = v566 - v574
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v578
						if int32(126) < v578 {
						} else {
							v582 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v583 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v582 << (uint(v583) % 32)
							v586 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v586 + v583
							v590 = m.G1
							v594 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v590+int32(_a_F_VP8WriteProbas_2)+v578))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v594
							if v586 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v600 = int32(1)
						if base.Ui32(v600) < base.Ui32(v563) {
							v563 = int32(base.Ui32(v563) >> (uint(v600) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v604 = m.G1
				v607 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(42)))))
				v611 = int32(38)
				v613 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v604+int32(_a_F_VP8WriteProbas_0)+v64+v611))))
				v614 = base.B2i32(v607 != v613)
				v620 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v604+int32(_a_F_VP8WriteProbas_1)+v64+v611))))
				v622 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v625 = v622 * v620 >> (uint(int32(8)) % 32)
				if v614 == int32(0) {
					v634 = v625
				} else {
					v628 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v630 = v625 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v628 + v630
					v634 = v622 - v630
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v634
				if int32(126) < v634 {
				} else {
					v638 = m.G1
					v642 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v638+int32(_a_F_VP8WriteProbas_2)+v634))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v642
					v644 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v648 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v638+int32(_a_F_VP8WriteProbas_3)+v634))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v644 << (uint(v648) % 32)
					v651 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v652 = v648 + v651
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v652
					if v652 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v614 == int32(0) {
				} else {
					v672 = int32(128)
					for {
						v675 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v677 = v675 >> (uint(int32(1)) % 32)
						if v672&v607 == int32(0) {
							v687 = v677
						} else {
							v681 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v683 = v677 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v681 + v683
							v687 = v675 - v683
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v687
						if int32(126) < v687 {
						} else {
							v691 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v692 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v691 << (uint(v692) % 32)
							v695 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v695 + v692
							v699 = m.G1
							v703 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v699+int32(_a_F_VP8WriteProbas_2)+v687))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v703
							if v695 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v709 = int32(1)
						if base.Ui32(v709) < base.Ui32(v672) {
							v672 = int32(base.Ui32(v672) >> (uint(v709) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v713 = m.G1
				v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(43)))))
				v720 = int32(39)
				v722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v713+int32(_a_F_VP8WriteProbas_0)+v64+v720))))
				v723 = base.B2i32(v716 != v722)
				v729 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v713+int32(_a_F_VP8WriteProbas_1)+v64+v720))))
				v731 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v734 = v731 * v729 >> (uint(int32(8)) % 32)
				if v723 == int32(0) {
					v743 = v734
				} else {
					v737 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v739 = v734 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v737 + v739
					v743 = v731 - v739
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v743
				if int32(126) < v743 {
				} else {
					v747 = m.G1
					v751 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v747+int32(_a_F_VP8WriteProbas_2)+v743))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v751
					v753 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v757 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v747+int32(_a_F_VP8WriteProbas_3)+v743))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v753 << (uint(v757) % 32)
					v760 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v761 = v757 + v760
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v761
					if v761 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v723 == int32(0) {
				} else {
					v781 = int32(128)
					for {
						v784 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v786 = v784 >> (uint(int32(1)) % 32)
						if v781&v716 == int32(0) {
							v796 = v786
						} else {
							v790 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v792 = v786 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v790 + v792
							v796 = v784 - v792
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v796
						if int32(126) < v796 {
						} else {
							v800 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v801 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v800 << (uint(v801) % 32)
							v804 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v804 + v801
							v808 = m.G1
							v812 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v808+int32(_a_F_VP8WriteProbas_2)+v796))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v812
							if v804 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v818 = int32(1)
						if base.Ui32(v818) < base.Ui32(v781) {
							v781 = int32(base.Ui32(v781) >> (uint(v818) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v822 = m.G1
				v825 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(44)))))
				v829 = int32(40)
				v831 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v822+int32(_a_F_VP8WriteProbas_0)+v64+v829))))
				v832 = base.B2i32(v825 != v831)
				v838 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v822+int32(_a_F_VP8WriteProbas_1)+v64+v829))))
				v840 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v843 = v840 * v838 >> (uint(int32(8)) % 32)
				if v832 == int32(0) {
					v852 = v843
				} else {
					v846 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v848 = v843 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v846 + v848
					v852 = v840 - v848
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v852
				if int32(126) < v852 {
				} else {
					v856 = m.G1
					v860 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v856+int32(_a_F_VP8WriteProbas_2)+v852))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v860
					v862 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v866 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v856+int32(_a_F_VP8WriteProbas_3)+v852))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v862 << (uint(v866) % 32)
					v869 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v870 = v866 + v869
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v870
					if v870 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v832 == int32(0) {
				} else {
					v890 = int32(128)
					for {
						v893 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v895 = v893 >> (uint(int32(1)) % 32)
						if v890&v825 == int32(0) {
							v905 = v895
						} else {
							v899 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v901 = v895 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v899 + v901
							v905 = v893 - v901
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v905
						if int32(126) < v905 {
						} else {
							v909 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v910 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v909 << (uint(v910) % 32)
							v913 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v913 + v910
							v917 = m.G1
							v921 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v917+int32(_a_F_VP8WriteProbas_2)+v905))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v921
							if v913 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v927 = int32(1)
						if base.Ui32(v927) < base.Ui32(v890) {
							v890 = int32(base.Ui32(v890) >> (uint(v927) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v931 = m.G1
				v934 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(45)))))
				v938 = int32(41)
				v940 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v931+int32(_a_F_VP8WriteProbas_0)+v64+v938))))
				v941 = base.B2i32(v934 != v940)
				v947 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v931+int32(_a_F_VP8WriteProbas_1)+v64+v938))))
				v949 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v952 = v949 * v947 >> (uint(int32(8)) % 32)
				if v941 == int32(0) {
					v961 = v952
				} else {
					v955 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v957 = v952 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v955 + v957
					v961 = v949 - v957
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v961
				if int32(126) < v961 {
				} else {
					v965 = m.G1
					v969 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v965+int32(_a_F_VP8WriteProbas_2)+v961))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v969
					v971 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v975 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v965+int32(_a_F_VP8WriteProbas_3)+v961))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v971 << (uint(v975) % 32)
					v978 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v979 = v975 + v978
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v979
					if v979 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v941 == int32(0) {
				} else {
					v999 = int32(128)
					for {
						v1002 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v1004 = v1002 >> (uint(int32(1)) % 32)
						if v999&v934 == int32(0) {
							v1014 = v1004
						} else {
							v1008 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v1010 = v1004 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1008 + v1010
							v1014 = v1002 - v1010
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1014
						if int32(126) < v1014 {
						} else {
							v1018 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v1019 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1018 << (uint(v1019) % 32)
							v1022 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1022 + v1019
							v1026 = m.G1
							v1030 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1026+int32(_a_F_VP8WriteProbas_2)+v1014))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1030
							if v1022 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v1036 = int32(1)
						if base.Ui32(v1036) < base.Ui32(v999) {
							v999 = int32(base.Ui32(v999) >> (uint(v1036) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v1040 = m.G1
				v1043 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(46)))))
				v1047 = int32(42)
				v1049 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1040+int32(_a_F_VP8WriteProbas_0)+v64+v1047))))
				v1050 = base.B2i32(v1043 != v1049)
				v1056 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1040+int32(_a_F_VP8WriteProbas_1)+v64+v1047))))
				v1058 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v1061 = v1058 * v1056 >> (uint(int32(8)) % 32)
				if v1050 == int32(0) {
					v1070 = v1061
				} else {
					v1064 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v1066 = v1061 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1064 + v1066
					v1070 = v1058 - v1066
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1070
				if int32(126) < v1070 {
				} else {
					v1074 = m.G1
					v1078 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1074+int32(_a_F_VP8WriteProbas_2)+v1070))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1078
					v1080 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v1084 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1074+int32(_a_F_VP8WriteProbas_3)+v1070))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1080 << (uint(v1084) % 32)
					v1087 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v1088 = v1084 + v1087
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1088
					if v1088 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v1050 == int32(0) {
				} else {
					v1108 = int32(128)
					for {
						v1111 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v1113 = v1111 >> (uint(int32(1)) % 32)
						if v1108&v1043 == int32(0) {
							v1123 = v1113
						} else {
							v1117 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v1119 = v1113 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1117 + v1119
							v1123 = v1111 - v1119
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1123
						if int32(126) < v1123 {
						} else {
							v1127 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v1128 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1127 << (uint(v1128) % 32)
							v1131 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1131 + v1128
							v1135 = m.G1
							v1139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1135+int32(_a_F_VP8WriteProbas_2)+v1123))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1139
							if v1131 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v1145 = int32(1)
						if base.Ui32(v1145) < base.Ui32(v1108) {
							v1108 = int32(base.Ui32(v1108) >> (uint(v1145) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v1149 = m.G1
				v1152 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+int32(47)))))
				v1156 = int32(43)
				v1158 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1149+int32(_a_F_VP8WriteProbas_0)+v64+v1156))))
				v1159 = base.B2i32(v1152 != v1158)
				v1165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1149+int32(_a_F_VP8WriteProbas_1)+v64+v1156))))
				v1167 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v1170 = v1167 * v1165 >> (uint(int32(8)) % 32)
				if v1159 == int32(0) {
					v1179 = v1170
				} else {
					v1173 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v1175 = v1170 + int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1173 + v1175
					v1179 = v1167 - v1175
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1179
				if int32(126) < v1179 {
				} else {
					v1183 = m.G1
					v1187 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1183+int32(_a_F_VP8WriteProbas_2)+v1179))))
					*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1187
					v1189 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v1193 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1183+int32(_a_F_VP8WriteProbas_3)+v1179))))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1189 << (uint(v1193) % 32)
					v1196 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					v1197 = v1193 + v1196
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1197
					if v1197 < int32(1) {
					} else {
						F_Flush(m, l0)
						mBase = m.M
					}
				}
				if v1159 == int32(0) {
				} else {
					v1217 = int32(128)
					for {
						v1220 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v1222 = v1220 >> (uint(int32(1)) % 32)
						if v1217&v1152 == int32(0) {
							v1232 = v1222
						} else {
							v1226 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v1228 = v1222 + int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1226 + v1228
							v1232 = v1220 - v1228
						}
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1232
						if int32(126) < v1232 {
						} else {
							v1236 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							v1237 = int32(1)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1236 << (uint(v1237) % 32)
							v1240 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1240 + v1237
							v1244 = m.G1
							v1248 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1244+int32(_a_F_VP8WriteProbas_2)+v1232))))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1248
							if v1240 < int32(0) {
							} else {
								F_Flush(m, l0)
								mBase = m.M
							}
						}
						v1254 = int32(1)
						if base.Ui32(v1254) < base.Ui32(v1217) {
							v1217 = int32(base.Ui32(v1217) >> (uint(v1254) % 32))
							continue
						} else {
							break
						}
						break
					}
				}
				v1259 = v52 + int32(11)
				if v1259 != 0 {
					v52 = v1259
					continue
				} else {
					break
				}
				break
			}
			v1260 = int32(33)
			v1265 = v37 + int32(1)
			if v1265 != int32(8) {
				v35 = v35 + v1260
				v36 = v36 + v1260
				v37 = v1265
				continue
			} else {
				break
			}
			break
		}
		v1268 = int32(264)
		v1273 = v20 + int32(1)
		if v1273 != int32(4) {
			v18 = v18 + v1268
			v19 = v19 + v1268
			v20 = v1273
			continue
		} else {
			break
		}
		break
	}
	v1276 = *(*int32)(unsafe.Add(mBase, uint32(l1)+uint32(_c_F_VP8WriteProbas[0])))
	v1279 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1281 = v1279 >> (uint(int32(1)) % 32)
	if v1276 == int32(0) {
		v1290 = v1281
	} else {
		v1284 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v1286 = v1281 + int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1284 + v1286
		v1290 = v1279 - v1286
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1290
	if int32(126) < v1290 {
	} else {
		v1294 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v1295 = int32(1)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1294 << (uint(v1295) % 32)
		v1298 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1298 + v1295
		v1302 = m.G1
		v1306 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1302+int32(_a_F_VP8WriteProbas_2)+v1290))))
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1306
		if v1298 < int32(0) {
		} else {
			F_Flush(m, l0)
			mBase = m.M
		}
	}
	if v1276 == int32(0) {
	} else {
		v1314 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1)+3)))
		v1326 = int32(128)
		for {
			v1329 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
			v1331 = v1329 >> (uint(int32(1)) % 32)
			if v1326&v1314 == int32(0) {
				v1341 = v1331
			} else {
				v1335 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v1337 = v1331 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1335 + v1337
				v1341 = v1329 - v1337
			}
			*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1341
			if int32(126) < v1341 {
			} else {
				v1345 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v1346 = int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v1345 << (uint(v1346) % 32)
				v1349 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v1349 + v1346
				v1353 = m.G1
				v1357 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1353+int32(_a_F_VP8WriteProbas_2)+v1341))))
				*(*int32)(unsafe.Add(mBase, uint32(l0))) = v1357
				if v1349 < int32(0) {
				} else {
					F_Flush(m, l0)
					mBase = m.M
				}
			}
			v1363 = int32(1)
			if base.Ui32(v1363) < base.Ui32(v1326) {
				v1326 = int32(base.Ui32(v1326) >> (uint(v1363) % 32))
				continue
			} else {
				break
			}
			break
		}
	}
	return
}

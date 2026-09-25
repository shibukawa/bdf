//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

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
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v345 base.V128
	_ = v345
	var v351 int32
	_ = v351
	var v352 int32
	_ = v352
	var v355 int32
	_ = v355
	var v359 int32
	_ = v359
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v403 int32
	_ = v403
	var v410 int32
	_ = v410
	var v427 int32
	_ = v427
	var v430 int32
	_ = v430
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v456 int32
	_ = v456
	var v459 int32
	_ = v459
	var v464 int32
	_ = v464
	var v465 float32
	_ = v465
	var v473 float32
	_ = v473
	var v479 int32
	_ = v479
	var v481 int32
	_ = v481
	var v482 int32
	_ = v482
	var v488 int32
	_ = v488
	var v491 int32
	_ = v491
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
		v342 = *(*int32)(unsafe.Add(mBase, uint32(v95)))
		v343 = *(*int32)(unsafe.Add(mBase, uint32(v342)+24))
		*(*int32)(unsafe.Add(mBase, uint32(v95)+24)) = v343
		v345 = base.Simd_g_const(&F_InitVP8Encoder__k0)
		base.Simd_g_v128_store(m, v95, int32(8), v345)
		*(*int32)(unsafe.Add(mBase, uint32(v95)+28)) = base.B2i32(int32(1) < v343)
		v351 = *(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[1])))
		v352 = *(*int32)(unsafe.Add(mBase, uint32(v95)+40))
		if v352 < v340 {
		} else {
			v355 = *(*int32)(unsafe.Add(mBase, uint32(v95)+48))
			v359 = int32(-1)
			for {
				v377 = int32(0)
				*(*uint8)(unsafe.Add(mBase, uint32(v351-v355+v359))) = uint8(v377)
				v380 = v359 + int32(1)
				v381 = *(*int32)(unsafe.Add(mBase, uint32(v95)+40))
				if v380 < v381<<(uint(int32(2))%32) {
					v359 = v380
					continue
				} else {
					break
				}
				break
			}
		}
		v403 = *(*int32)(unsafe.Add(mBase, uint32(v95)+44))
		if v403 < int32(1) {
		} else {
			v410 = int32(0)
			for {
				v427 = *(*int32)(unsafe.Add(mBase, uint32(v95)+48))
				v430 = int32(0)
				*(*uint8)(unsafe.Add(mBase, uint32(v351+int32(-1)+v427*v410))) = uint8(v430)
				v433 = v410 + int32(1)
				v434 = *(*int32)(unsafe.Add(mBase, uint32(v95)+44))
				if v433 < v434<<(uint(int32(2))%32) {
					v410 = v433
					continue
				} else {
					break
				}
				break
			}
		}
		v456 = *(*int32)(unsafe.Add(mBase, uint32(v95)+uint32(_c_F_InitVP8Encoder[2])))
		v459 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v456+int32(-4)))) = v459
		F_VP8EncDspCostInit(m)
		mBase = m.M
		F_VP8EncInitAlpha(m, v95)
		mBase = m.M
		v464 = v95 + int32(344)
		v465 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
		v473 = base.F32_mul(base.F32_add(base.F32_div(base.F32_mul(v465, float32(5)), float32(100)), float32(1)), base.F32_convert_i32_s(v35))
		if base.F32_lt(base.F32_abs(v473), float32(2.1474836e+09)) == v459 {
			v481 = int32(-2147483648)
		} else {
			v479 = base.I32_trunc_f32_s(v473)
			v481 = v479
		}
		v482 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v464)+20)) = v482
		*(*int32)(unsafe.Add(mBase, uint32(v464))) = v482
		*(*int64)(unsafe.Add(mBase, uint32(v464)+8)) = int64(0)
		v488 = int32(_a_F_InitVP8Encoder_4)
		if v488 < v481 {
			v491 = v481
		} else {
			v491 = v488
		}
		*(*int32)(unsafe.Add(mBase, uint32(v464)+16)) = v491
		*(*int32)(unsafe.Add(mBase, uint32(v464)+4)) = v464
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

var F_InitVP8Encoder__k0 = [2]uint64{0x1, 0x0}

func F_VP8BitWriterInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v5 base.V128
	_ = v5
	var v8 base.V128
	_ = v8
	var v13 int64
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v24 int64
	_ = v24
	var v25 int32
	_ = v25
	var v26 int64
	_ = v26
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v44 int32
	_ = v44
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	v3 = int32(0)
	v5 = base.Simd_g_const(&F_VP8BitWriterInit__k0)
	base.Simd_g_v128_store(m, l0, v3, v5)
	v8 = base.Simd_g_const(&F_VP8BitWriterInit__k1)
	base.Simd_g_v128_store(m, l0, int32(16), v8)
	if l1 == v3 {
		return int32(1)
	} else {
		v13 = int64(1)
		v14 = int32(1024)
		if base.Ui32(v14) < base.Ui32(l1) {
			v17 = l1
		} else {
			v17 = v14
		}
		v24 = base.I64_div_u_s(int64(2147418112), v13)
		v25 = int32(0)
		v26 = base.I64_extend_i32_u(v17)
		if base.Ui64(int64(4294967295)) < base.Ui64(v26*v13) {
			v38 = v25
		} else {
			if base.Ui64(v24) < base.Ui64(v26) {
				v38 = v25
			} else {
				v36 = F_malloc(m, base.I32_wrap_i64(v13)*v17)
				mBase = m.M
				v38 = v36
			}
		}
		if v38 != 0 {
			v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
			if v44 == int32(0) {
			} else {
				v49 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(16))))
				v50 = F_memcpy(m, v38, v49, v44)
				mBase = m.M
			}
			v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
			F_free(m, v51)
			mBase = m.M
			*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v17
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v38
			return int32(1)
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = int32(1)
			return int32(0)
		}
	}
}

var F_VP8BitWriterInit__k0 = [2]uint64{0xfe, 0xfffffff800000000}
var F_VP8BitWriterInit__k1 = [2]uint64{0x0, 0x0}

func F_VP8BitWriterWipeOut(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v10 base.V128
	_ = v10
	var v11 int32
	_ = v11
	if l0 == int32(0) {
	} else {
		v7 = l0 + int32(16)
		v8 = *(*int32)(unsafe.Add(mBase, uint32(v7)))
		F_free(m, v8)
		mBase = m.M
		v10 = base.Simd_g_const(&F_VP8BitWriterWipeOut__k0)
		v11 = int32(0)
		base.Simd_g_v128_store(m, v7, v11, v10)
		base.Simd_g_v128_store(m, l0, v11, v10)
	}
	return
}

var F_VP8BitWriterWipeOut__k0 = [2]uint64{0x0, 0x0}

func F_VP8SetSegmentParams(m *base.Module, l0 int32, l1 float32) {
	mBase := m.M
	_ = mBase
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v37 float64
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v45 int32
	_ = v45
	var v48 float64
	_ = v48
	var v51 float64
	_ = v51
	var v60 float64
	_ = v60
	var v61 float64
	_ = v61
	var v62 float64
	_ = v62
	var v70 float64
	_ = v70
	var v72 float64
	_ = v72
	var v73 float64
	_ = v73
	var v76 int32
	_ = v76
	var v81 int32
	_ = v81
	var v91 float64
	_ = v91
	var v94 base.V128
	_ = v94
	var v96 int32
	_ = v96
	var v107 int32
	_ = v107
	var v111 int32
	_ = v111
	var v127 base.V128
	_ = v127
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	var v137 int32
	_ = v137
	var v140 base.V128
	_ = v140
	var v141 int32
	_ = v141
	var v143 float64
	_ = v143
	var v147 float64
	_ = v147
	var v151 base.V128
	_ = v151
	var v152 base.V128
	_ = v152
	var v154 float64
	_ = v154
	var v160 int32
	_ = v160
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 float64
	_ = v164
	var v170 int32
	_ = v170
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v178 int32
	_ = v178
	var v183 int32
	_ = v183
	var v188 base.V128
	_ = v188
	var v189 int32
	_ = v189
	var v191 float64
	_ = v191
	var v195 float64
	_ = v195
	var v199 base.V128
	_ = v199
	var v201 float64
	_ = v201
	var v207 int32
	_ = v207
	var v209 int32
	_ = v209
	var v213 float64
	_ = v213
	var v219 int32
	_ = v219
	var v221 int32
	_ = v221
	var v224 int32
	_ = v224
	var v229 base.V128
	_ = v229
	var v230 int32
	_ = v230
	var v236 int32
	_ = v236
	var v241 int32
	_ = v241
	var v251 int32
	_ = v251
	var v262 int32
	_ = v262
	var v299 int32
	_ = v299
	var v300 int32
	_ = v300
	var v319 float64
	_ = v319
	var v322 int32
	_ = v322
	var v327 float64
	_ = v327
	var v330 float64
	_ = v330
	var v336 int32
	_ = v336
	var v338 int32
	_ = v338
	var v339 int32
	_ = v339
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v346 int32
	_ = v346
	var v351 int32
	_ = v351
	var v382 int32
	_ = v382
	var v385 int32
	_ = v385
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v402 int32
	_ = v402
	var v423 int32
	_ = v423
	var v446 int32
	_ = v446
	var v450 int32
	_ = v450
	var v466 int32
	_ = v466
	var v476 int32
	_ = v476
	var v520 int32
	_ = v520
	var v524 int32
	_ = v524
	var v563 int32
	_ = v563
	var v594 int32
	_ = v594
	var v605 int32
	_ = v605
	var v606 int32
	_ = v606
	var v625 int32
	_ = v625
	var v631 int32
	_ = v631
	var v632 int32
	_ = v632
	var v635 int32
	_ = v635
	var v636 int32
	_ = v636
	var v639 int32
	_ = v639
	var v641 int32
	_ = v641
	var v647 int32
	_ = v647
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v654 int32
	_ = v654
	var v655 int32
	_ = v655
	var v658 int32
	_ = v658
	var v661 int32
	_ = v661
	var v662 int32
	_ = v662
	var v664 int32
	_ = v664
	var v666 int32
	_ = v666
	var v667 int32
	_ = v667
	var v670 int32
	_ = v670
	var v671 int32
	_ = v671
	var v674 int32
	_ = v674
	var v678 int32
	_ = v678
	var v680 int32
	_ = v680
	var v681 int32
	_ = v681
	var v687 int32
	_ = v687
	var v690 int32
	_ = v690
	var v692 int32
	_ = v692
	var v694 int32
	_ = v694
	var v696 int32
	_ = v696
	var v699 int32
	_ = v699
	var v700 int32
	_ = v700
	var v703 int32
	_ = v703
	var v706 int32
	_ = v706
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v718 int32
	_ = v718
	var v722 int32
	_ = v722
	var v724 int32
	_ = v724
	var v725 int32
	_ = v725
	var v731 int32
	_ = v731
	var v734 int32
	_ = v734
	var v736 int32
	_ = v736
	var v738 int32
	_ = v738
	var v741 int32
	_ = v741
	var v742 int32
	_ = v742
	var v745 int32
	_ = v745
	var v748 int32
	_ = v748
	var v751 int32
	_ = v751
	var v752 int32
	_ = v752
	var v753 int32
	_ = v753
	var v756 int32
	_ = v756
	var v757 int32
	_ = v757
	var v760 int32
	_ = v760
	var v764 int32
	_ = v764
	var v766 int32
	_ = v766
	var v767 int32
	_ = v767
	var v773 int32
	_ = v773
	var v776 int32
	_ = v776
	var v778 int32
	_ = v778
	var v780 int32
	_ = v780
	var v783 int32
	_ = v783
	var v784 int32
	_ = v784
	var v787 int32
	_ = v787
	var v790 int32
	_ = v790
	var v792 int32
	_ = v792
	var v793 int32
	_ = v793
	var v794 int32
	_ = v794
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v801 int32
	_ = v801
	var v805 int32
	_ = v805
	var v807 int32
	_ = v807
	var v808 int32
	_ = v808
	var v814 int32
	_ = v814
	var v817 int32
	_ = v817
	var v819 int32
	_ = v819
	var v820 int32
	_ = v820
	var v822 int32
	_ = v822
	var v823 int32
	_ = v823
	var v825 int32
	_ = v825
	var v826 int32
	_ = v826
	var v831 int32
	_ = v831
	var v834 int32
	_ = v834
	var v835 int32
	_ = v835
	var v838 int32
	_ = v838
	var v841 int32
	_ = v841
	var v845 int32
	_ = v845
	var v848 int32
	_ = v848
	var v849 base.V128
	_ = v849
	var v852 int32
	_ = v852
	var v853 int32
	_ = v853
	var v856 int32
	_ = v856
	var v861 int32
	_ = v861
	var v874 int32
	_ = v874
	var v877 int32
	_ = v877
	var v895 int32
	_ = v895
	var v897 int32
	_ = v897
	var v898 int32
	_ = v898
	var v902 int32
	_ = v902
	var v903 int32
	_ = v903
	var v906 int32
	_ = v906
	var v913 int32
	_ = v913
	var v914 int32
	_ = v914
	var v923 int32
	_ = v923
	var v936 int32
	_ = v936
	var v939 int32
	_ = v939
	var v957 int32
	_ = v957
	var v959 int32
	_ = v959
	var v960 int32
	_ = v960
	var v965 int32
	_ = v965
	var v968 int32
	_ = v968
	var v972 int32
	_ = v972
	var v976 int32
	_ = v976
	var v985 int32
	_ = v985
	var v998 int32
	_ = v998
	var v1001 int32
	_ = v1001
	var v1019 int32
	_ = v1019
	var v1021 int32
	_ = v1021
	var v1022 int32
	_ = v1022
	var v1027 int32
	_ = v1027
	var v1032 int32
	_ = v1032
	var v1036 int32
	_ = v1036
	var v1045 int32
	_ = v1045
	var v1071 int32
	_ = v1071
	var v1072 int32
	_ = v1072
	var v1073 int32
	_ = v1073
	var v1080 int32
	_ = v1080
	var v1082 int32
	_ = v1082
	var v1085 int32
	_ = v1085
	var v1086 int32
	_ = v1086
	var v1092 int32
	_ = v1092
	var v1099 int32
	_ = v1099
	var v1102 int32
	_ = v1102
	var v1122 int32
	_ = v1122
	var v1125 int32
	_ = v1125
	var v1141 int32
	_ = v1141
	var v1142 int32
	_ = v1142
	var v1143 int32
	_ = v1143
	var v1144 int32
	_ = v1144
	var v1146 int32
	_ = v1146
	var v1149 int32
	_ = v1149
	var v1150 int32
	_ = v1150
	var v1152 int32
	_ = v1152
	var v1154 int32
	_ = v1154
	var v1156 int32
	_ = v1156
	var v1158 int32
	_ = v1158
	var v1161 int32
	_ = v1161
	var v1162 int32
	_ = v1162
	var v1168 int32
	_ = v1168
	var v1175 int32
	_ = v1175
	var v1180 int32
	_ = v1180
	var v1216 int32
	_ = v1216
	var v1218 int32
	_ = v1218
	var v1221 int32
	_ = v1221
	var v1236 int32
	_ = v1236
	var v1243 int32
	_ = v1243
	var v1255 int32
	_ = v1255
	var v1256 int32
	_ = v1256
	var v1260 int32
	_ = v1260
	var v1267 int32
	_ = v1267
	var v1309 int32
	_ = v1309
	var v1316 int32
	_ = v1316
	var v1328 int32
	_ = v1328
	var v1329 int32
	_ = v1329
	var v1333 int32
	_ = v1333
	var v1337 int32
	_ = v1337
	var v1341 int32
	_ = v1341
	var v1345 int32
	_ = v1345
	var v1376 int32
	_ = v1376
	var v1379 int32
	_ = v1379
	var v1380 int32
	_ = v1380
	var v1381 int32
	_ = v1381
	var v1382 int32
	_ = v1382
	var v1390 int32
	_ = v1390
	var v1395 int32
	_ = v1395
	var v1417 int32
	_ = v1417
	var v1420 int32
	_ = v1420
	var v1422 int32
	_ = v1422
	var v1423 int32
	_ = v1423
	var v1426 int32
	_ = v1426
	var v1427 int32
	_ = v1427
	var v1430 int32
	_ = v1430
	var v1434 int32
	_ = v1434
	var v1437 int32
	_ = v1437
	var v1438 int32
	_ = v1438
	var v1439 int32
	_ = v1439
	var v1440 int32
	_ = v1440
	var v1443 int32
	_ = v1443
	var v1444 int32
	_ = v1444
	var v1447 int32
	_ = v1447
	var v1449 int32
	_ = v1449
	var v1452 int32
	_ = v1452
	var v1453 int32
	_ = v1453
	var v1454 int32
	_ = v1454
	var v1455 int32
	_ = v1455
	var v1458 int32
	_ = v1458
	var v1459 int32
	_ = v1459
	var v1462 int32
	_ = v1462
	var v1464 int32
	_ = v1464
	var v1467 int32
	_ = v1467
	var v1468 int32
	_ = v1468
	var v1469 int32
	_ = v1469
	var v1470 int32
	_ = v1470
	var v1473 int32
	_ = v1473
	var v1474 int32
	_ = v1474
	var v1477 int32
	_ = v1477
	var v1479 int32
	_ = v1479
	var v1481 int32
	_ = v1481
	var v1487 int32
	_ = v1487
	var v1488 int32
	_ = v1488
	var v1489 int32
	_ = v1489
	var v1492 int32
	_ = v1492
	var v1493 int32
	_ = v1493
	var v1496 int32
	_ = v1496
	var v1500 int32
	_ = v1500
	var v1504 int32
	_ = v1504
	var v1505 int32
	_ = v1505
	var v1506 int32
	_ = v1506
	var v1509 int32
	_ = v1509
	var v1510 int32
	_ = v1510
	var v1513 int32
	_ = v1513
	var v1514 int32
	_ = v1514
	var v1517 int32
	_ = v1517
	var v1519 int32
	_ = v1519
	var v1526 int32
	_ = v1526
	var v1530 int32
	_ = v1530
	var v1531 int32
	_ = v1531
	var v1532 int32
	_ = v1532
	var v1535 int32
	_ = v1535
	var v1540 int32
	_ = v1540
	var v1546 int32
	_ = v1546
	var v1547 int32
	_ = v1547
	var v1548 int32
	_ = v1548
	var v1550 int32
	_ = v1550
	var v1552 int32
	_ = v1552
	var v1557 int32
	_ = v1557
	var v1559 int32
	_ = v1559
	var v1561 int32
	_ = v1561
	var v1567 int32
	_ = v1567
	var v1597 int32
	_ = v1597
	var v1612 base.V128
	_ = v1612
	var v1623 int32
	_ = v1623
	var v1660 int32
	_ = v1660
	var v1662 int32
	_ = v1662
	var v1663 int32
	_ = v1663
	var v1667 int32
	_ = v1667
	var v1670 int32
	_ = v1670
	var v1696 int32
	_ = v1696
	var v1699 int32
	_ = v1699
	var v1712 int32
	_ = v1712
	var v1716 int32
	_ = v1716
	var v1720 int32
	_ = v1720
	var v1724 int32
	_ = v1724
	var v1728 int32
	_ = v1728
	var v1732 int32
	_ = v1732
	var v1734 int32
	_ = v1734
	var v1735 int32
	_ = v1735
	var v1756 int32
	_ = v1756
	var v1758 int32
	_ = v1758
	var v1764 int32
	_ = v1764
	var v1768 int32
	_ = v1768
	var v1769 int32
	_ = v1769
	var v1770 int32
	_ = v1770
	var v1773 int32
	_ = v1773
	var v1778 int32
	_ = v1778
	var v1784 int32
	_ = v1784
	var v1785 int32
	_ = v1785
	var v1786 int32
	_ = v1786
	var v1788 int32
	_ = v1788
	var v1790 int32
	_ = v1790
	var v1795 int32
	_ = v1795
	var v1797 int32
	_ = v1797
	var v1799 int32
	_ = v1799
	var v1805 int32
	_ = v1805
	var v1835 int32
	_ = v1835
	var v1850 base.V128
	_ = v1850
	var v1861 int32
	_ = v1861
	var v1866 int64
	_ = v1866
	var v1870 base.V128
	_ = v1870
	var v1877 int32
	_ = v1877
	var v1878 int32
	_ = v1878
	var v1879 int32
	_ = v1879
	var v1881 int32
	_ = v1881
	var v1883 int32
	_ = v1883
	var v1994 int32
	_ = v1994
	var v1996 int32
	_ = v1996
	var v2002 int32
	_ = v2002
	var v2006 int32
	_ = v2006
	var v2007 int32
	_ = v2007
	var v2008 int32
	_ = v2008
	var v2011 int32
	_ = v2011
	var v2016 int32
	_ = v2016
	var v2022 int32
	_ = v2022
	var v2023 int32
	_ = v2023
	var v2024 int32
	_ = v2024
	var v2026 int32
	_ = v2026
	var v2028 int32
	_ = v2028
	var v2033 int32
	_ = v2033
	var v2035 int32
	_ = v2035
	var v2037 int32
	_ = v2037
	var v2043 int32
	_ = v2043
	var v2073 int32
	_ = v2073
	var v2088 base.V128
	_ = v2088
	var v2099 int32
	_ = v2099
	var v2104 int64
	_ = v2104
	var v2108 base.V128
	_ = v2108
	var v2115 int32
	_ = v2115
	var v2116 int32
	_ = v2116
	var v2117 int32
	_ = v2117
	var v2119 int32
	_ = v2119
	var v2121 int32
	_ = v2121
	var v2232 int32
	_ = v2232
	var v2234 int32
	_ = v2234
	var v2235 int32
	_ = v2235
	var v2236 int32
	_ = v2236
	var v2240 int32
	_ = v2240
	var v2243 int32
	_ = v2243
	var v2246 int32
	_ = v2246
	var v2249 int32
	_ = v2249
	var v2250 int32
	_ = v2250
	var v2253 int32
	_ = v2253
	var v2254 int32
	_ = v2254
	var v2256 int32
	_ = v2256
	var v2259 int32
	_ = v2259
	var v2263 int32
	_ = v2263
	var v2266 int32
	_ = v2266
	var v2271 int32
	_ = v2271
	var v2272 int32
	_ = v2272
	var v2274 int32
	_ = v2274
	var v2277 int32
	_ = v2277
	var v2281 int32
	_ = v2281
	var v2321 int32
	_ = v2321
	var v2334 int32
	_ = v2334
	v31 = m.G0
	v33 = v31 - int32(16)
	m.G0 = v33
	v37 = base.F64_div(base.F64_promote_f32(l1), float64(100))
	v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v39 = *(*int32)(unsafe.Add(mBase, uint32(v38)+28))
	v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v41 = *(*int32)(unsafe.Add(mBase, uint32(v38)+80))
	if v41 == int32(0) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	v76 = int32(0)
	if v76 < v40 {
		goto L14
	} else {
		goto L15
	}
L2:
	;
	if base.F64_lt(v37, float64(0.75)) != 0 {
		goto L9
	} else {
		goto L10
	}
L3:
	;
	v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3388))
	v48 = base.F64_div(base.F64_convert_i32_s(v45), float64(255))
	if base.F64_gt(v48, float64(0.85)) != 0 {
		v61 = float64(0.4)
		goto L4
	} else {
		goto L5
	}
L4:
	;
	v62 = F_pow(m, v37, v61)
	mBase = m.M
	v73 = v62
	goto L1
L5:
	;
	v51 = float64(0.9)
	if base.F64_lt(v48, float64(0.3)) != 0 {
		goto L6
	} else {
		goto L7
	}
L6:
	;
	v60 = v51
	goto L8
L7:
	;
	v60 = base.F64_add(base.F64_mul(base.F64_add(v48, float64(-0.3)), float64(-0.9090909090909091)), v51)
	goto L8
L8:
	;
	v61 = v60
	goto L4
L9:
	;
	v70 = base.F64_mul(v37, float64(0.6666666666666666))
	goto L11
L10:
	;
	v70 = base.F64_add(base.F64_add(v37, v37), float64(-1))
	goto L11
L11:
	;
	v72 = F_pow(m, v70, float64(0.3333333333333333))
	mBase = m.M
	v73 = v72
	goto L1
L12:
	;
	v625 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3404)) = v625
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3396)) = int64(0)
	v631 = base.I32_div_s(v39, int32(-25))
	v632 = int32(15)
	if v631 < v632 {
		goto L58
	} else {
		goto L59
	}
L13:
	;
	v423 = (int32(4) - v40) & int32(7)
	if v423 == int32(0) {
		v476 = v40
		goto L48
	} else {
		goto L49
	}
L14:
	;
	v91 = base.F64_mul(base.F64_div(base.F64_mul(base.F64_convert_i32_s(v39), float64(0.9)), float64(100)), float64(-0.0078125))
	if base.Ui32(v40) < base.Ui32(int32(4)) {
		v262 = v76
		goto L17
	} else {
		goto L18
	}
L15:
	;
	v81 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3384)) = v81
	v400 = l0 + int32(408)
	v401 = v81
	v402 = l0 + int32(1088)
	goto L13
L16:
	;
	v382 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3384)) = v382
	v385 = l0 + int32(408)
	if int32(3) < v40 {
		v605 = v385
		v606 = v382
		goto L12
	} else {
		goto L47
	}
L17:
	;
	v299 = v262*int32(744) + l0 + int32(1088)
	v300 = v40 - v262
	goto L35
L18:
	;
	v94 = base.Simd_g_f64x2_splat(v91)
	v96 = v40 & int32(2147483644)
	v107 = l0
	v111 = v96
	goto L19
L19:
	;
	v127 = base.Simd_g_const(&F_VP8SetSegmentParams__k0)
	v130 = *(*int32)(unsafe.Add(mBase, uint32(v107+int32(1080))))
	v135 = *(*int32)(unsafe.Add(mBase, uint32(v107+int32(1824))))
	v137 = int32(1)
	v140 = base.Simd_g_f64x2_add(v127, base.Simd_g_f64x2_mul(v94, base.Simd_g_f64x2_replace_lane_l1(base.Simd_g_f64x2_splat(base.F64_convert_i32_s(v130)), base.F64_convert_i32_s(v135))))
	v141 = int32(0)
	v143 = F_pow(m, v73, base.Simd_g_f64x2_extract_lane_l0(v140))
	mBase = m.M
	v147 = F_pow(m, v73, base.Simd_g_f64x2_extract_lane_l1(v140))
	mBase = m.M
	v151 = base.Simd_g_const(&F_VP8SetSegmentParams__k1)
	v152 = base.Simd_g_f64x2_mul(base.Simd_g_f64x2_sub(v127, base.Simd_g_f64x2_replace_lane_l1(base.Simd_g_f64x2_splat(v143), v147)), v151)
	v154 = base.Simd_g_f64x2_extract_lane_l1(v152)
	if base.F64_lt(base.F64_abs(v154), float64(2.147483648e+09)) == v141 {
		goto L22
	} else {
		goto L23
	}
L20:
	;
	if v40 == v96 {
		goto L16
	} else {
		goto L34
	}
L21:
	;
	v163 = int32(0)
	v164 = base.Simd_g_f64x2_extract_lane_l0(v152)
	if base.F64_lt(base.F64_abs(v164), float64(2.147483648e+09)) == v163 {
		goto L25
	} else {
		goto L26
	}
L22:
	;
	v162 = int32(-2147483648)
	goto L21
L23:
	;
	v160 = base.I32_trunc_f64_s(v154)
	v162 = v160
	goto L21
L24:
	;
	v174 = int32(1)
	v178 = *(*int32)(unsafe.Add(mBase, uint32(v107+int32(2568))))
	v183 = *(*int32)(unsafe.Add(mBase, uint32(v107+int32(3312))))
	v188 = base.Simd_g_f64x2_add(base.Simd_g_f64x2_mul(v94, base.Simd_g_f64x2_replace_lane_l1(base.Simd_g_f64x2_splat(base.F64_convert_i32_s(v178)), base.F64_convert_i32_s(v183))), v127)
	v189 = int32(0)
	v191 = F_pow(m, v73, base.Simd_g_f64x2_extract_lane_l0(v188))
	mBase = m.M
	v195 = F_pow(m, v73, base.Simd_g_f64x2_extract_lane_l1(v188))
	mBase = m.M
	v199 = base.Simd_g_f64x2_mul(base.Simd_g_f64x2_sub(v127, base.Simd_g_f64x2_replace_lane_l1(base.Simd_g_f64x2_splat(v191), v195)), v151)
	v201 = base.Simd_g_f64x2_extract_lane_l0(v199)
	if base.F64_lt(base.F64_abs(v201), float64(2.147483648e+09)) == v189 {
		goto L28
	} else {
		goto L29
	}
L25:
	;
	v172 = int32(-2147483648)
	goto L24
L26:
	;
	v170 = base.I32_trunc_f64_s(v164)
	v172 = v170
	goto L24
L27:
	;
	v213 = base.Simd_g_f64x2_extract_lane_l1(v199)
	if base.F64_lt(base.F64_abs(v213), float64(2.147483648e+09)) == int32(0) {
		goto L31
	} else {
		goto L32
	}
L28:
	;
	v209 = int32(-2147483648)
	goto L27
L29:
	;
	v207 = base.I32_trunc_f64_s(v201)
	v209 = v207
	goto L27
L30:
	;
	v224 = int32(3)
	v229 = base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v172), v162), v209), v221), base.Simd_g_const(&F_VP8SetSegmentParams__k2)), base.Simd_g_const(&F_VP8SetSegmentParams__k3))
	v230 = int32(0)
	base.Simd_g_v128_store32_lane_l3(m, v107+int32(3320), v230, v229)
	v236 = int32(2)
	base.Simd_g_v128_store32_lane_l2(m, v107+int32(2576), v230, v229)
	v241 = int32(1)
	base.Simd_g_v128_store32_lane_l1(m, v107+int32(1832), v230, v229)
	base.Simd_g_v128_store32_lane_l0(m, v107+int32(1088), v230, v229)
	v251 = v111 + int32(-4)
	if v251 != 0 {
		v107 = v107 + int32(2976)
		v111 = v251
		goto L19
	} else {
		goto L33
	}
L31:
	;
	v221 = int32(-2147483648)
	goto L30
L32:
	;
	v219 = base.I32_trunc_f64_s(v213)
	v221 = v219
	goto L30
L33:
	;
	goto L20
L34:
	;
	v262 = v96
	goto L17
L35:
	;
	v319 = float64(1)
	v322 = *(*int32)(unsafe.Add(mBase, uint32(v299+int32(-8))))
	v327 = F_pow(m, v73, base.F64_add(base.F64_mul(v91, base.F64_convert_i32_s(v322)), v319))
	mBase = m.M
	v330 = base.F64_mul(base.F64_sub(v319, v327), float64(127))
	if base.F64_lt(base.F64_abs(v330), float64(2.147483648e+09)) == int32(0) {
		goto L38
	} else {
		goto L39
	}
L36:
	;
	goto L16
L37:
	;
	v339 = int32(127)
	if v338 < v339 {
		goto L40
	} else {
		goto L41
	}
L38:
	;
	v338 = int32(-2147483648)
	goto L37
L39:
	;
	v336 = base.I32_trunc_f64_s(v330)
	v338 = v336
	goto L37
L40:
	;
	v342 = v338
	goto L42
L41:
	;
	v342 = v339
	goto L42
L42:
	;
	v343 = int32(0)
	if v343 < v342 {
		goto L43
	} else {
		goto L44
	}
L43:
	;
	v346 = v342
	goto L45
L44:
	;
	v346 = v343
	goto L45
L45:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v299))) = v346
	v351 = v300 + int32(-1)
	if v351 != 0 {
		v299 = v299 + int32(744)
		v300 = v351
		goto L35
	} else {
		goto L46
	}
L46:
	;
	goto L36
L47:
	;
	v400 = v385
	v401 = v382
	v402 = l0 + int32(1088)
	goto L13
L48:
	;
	if base.Ui32(v40+int32(3)) < base.Ui32(int32(7)) {
		goto L53
	} else {
		goto L54
	}
L49:
	;
	v446 = v40*int32(744) + l0 + int32(1088)
	v450 = v423
	goto L50
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v446))) = v401
	v466 = v450 + int32(-1)
	if v466 != 0 {
		v446 = v446 + int32(744)
		v450 = v466
		goto L50
	} else {
		goto L52
	}
L51:
	;
	v476 = v40 + v423
	goto L48
L52:
	;
	goto L51
L53:
	;
	v594 = *(*int32)(unsafe.Add(mBase, uint32(v402)))
	v605 = v400
	v606 = v594
	goto L12
L54:
	;
	v520 = l0 + v476*int32(744)
	v524 = v476 + int32(-4)
	goto L55
L55:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v520)+uint32(_c_F_VP8SetSegmentParams[0]))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520)+uint32(_c_F_VP8SetSegmentParams[1]))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520)+uint32(_c_F_VP8SetSegmentParams[2]))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520+int32(4064)))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520+int32(3320)))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520+int32(2576)))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520+int32(1832)))) = v401
	*(*int32)(unsafe.Add(mBase, uint32(v520+int32(1088)))) = v401
	v563 = v524 + int32(8)
	if v563 != 0 {
		v520 = v520 + int32(_a_F_VP8SetSegmentParams_0)
		v524 = v563
		goto L55
	} else {
		goto L57
	}
L56:
	;
	goto L53
L57:
	;
	goto L56
L58:
	;
	v635 = v631
	goto L60
L59:
	;
	v635 = v632
	goto L60
L60:
	;
	v636 = int32(-15)
	if v636 < v635 {
		goto L61
	} else {
		goto L62
	}
L61:
	;
	v639 = v635
	goto L63
L62:
	;
	v639 = v636
	goto L63
L63:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3408)) = v639
	v641 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3392))
	v647 = base.I32_div_s(v641*int32(10)+int32(-640), int32(70))
	v650 = base.I32_div_s(v647*v39, int32(100))
	v651 = int32(6)
	if v650 < v651 {
		goto L64
	} else {
		goto L65
	}
L64:
	;
	v654 = v650
	goto L66
L65:
	;
	v654 = v651
	goto L66
L66:
	;
	v655 = int32(-4)
	if v655 < v654 {
		goto L67
	} else {
		goto L68
	}
L67:
	;
	v658 = v654
	goto L69
L68:
	;
	v658 = v655
	goto L69
L69:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3412)) = v658
	v661 = m.G1
	v662 = *(*int32)(unsafe.Add(mBase, uint32(v38)+32))
	v664 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v666 = v661 + int32(_a_F_VP8SetSegmentParams_1)
	v667 = int32(127)
	if v606 < v667 {
		goto L70
	} else {
		goto L71
	}
L70:
	;
	v670 = v606
	goto L72
L71:
	;
	v670 = v667
	goto L72
L72:
	;
	v671 = int32(0)
	if v671 < v670 {
		goto L73
	} else {
		goto L74
	}
L73:
	;
	v674 = v670
	goto L75
L74:
	;
	v674 = v671
	goto L75
L75:
	;
	v678 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v666+v674<<(uint(int32(1))%32)))))
	v680 = int32(base.Ui32(v678) >> (uint(int32(2)) % 32))
	v681 = m.G1
	v687 = int32(63)
	if v680 < v687 {
		goto L77
	} else {
		goto L78
	}
L76:
	;
	v694 = v662 * int32(5)
	v696 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1084))
	v699 = base.I32_div_s(v692*v694, v696+int32(256))
	v700 = int32(63)
	if v699 < v700 {
		goto L80
	} else {
		goto L81
	}
L77:
	;
	v690 = v680
	goto L79
L78:
	;
	v690 = v687
	goto L79
L79:
	;
	v692 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v681+int32(_a_F_VP8SetSegmentParams_2)+v664<<(uint(int32(6))%32)+v690))))
	goto L76
L80:
	;
	v703 = v699
	goto L82
L81:
	;
	v703 = v700
	goto L82
L82:
	;
	if v699 < int32(2) {
		goto L83
	} else {
		goto L84
	}
L83:
	;
	v706 = int32(0)
	goto L85
L84:
	;
	v706 = v703
	goto L85
L85:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+1092)) = v706
	v709 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v710 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1832))
	v711 = int32(127)
	if v710 < v711 {
		goto L86
	} else {
		goto L87
	}
L86:
	;
	v714 = v710
	goto L88
L87:
	;
	v714 = v711
	goto L88
L88:
	;
	v715 = int32(0)
	if v715 < v714 {
		goto L89
	} else {
		goto L90
	}
L89:
	;
	v718 = v714
	goto L91
L90:
	;
	v718 = v715
	goto L91
L91:
	;
	v722 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v666+v718<<(uint(int32(1))%32)))))
	v724 = int32(base.Ui32(v722) >> (uint(int32(2)) % 32))
	v725 = m.G1
	v731 = int32(63)
	if v724 < v731 {
		goto L93
	} else {
		goto L94
	}
L92:
	;
	v738 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1828))
	v741 = base.I32_div_s(v736*v694, v738+int32(256))
	v742 = int32(63)
	if v741 < v742 {
		goto L96
	} else {
		goto L97
	}
L93:
	;
	v734 = v724
	goto L95
L94:
	;
	v734 = v731
	goto L95
L95:
	;
	v736 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v725+int32(_a_F_VP8SetSegmentParams_2)+v709<<(uint(int32(6))%32)+v734))))
	goto L92
L96:
	;
	v745 = v741
	goto L98
L97:
	;
	v745 = v742
	goto L98
L98:
	;
	if v741 < int32(2) {
		goto L99
	} else {
		goto L100
	}
L99:
	;
	v748 = int32(0)
	goto L101
L100:
	;
	v748 = v745
	goto L101
L101:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+1836)) = v748
	v751 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v752 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2576))
	v753 = int32(127)
	if v752 < v753 {
		goto L102
	} else {
		goto L103
	}
L102:
	;
	v756 = v752
	goto L104
L103:
	;
	v756 = v753
	goto L104
L104:
	;
	v757 = int32(0)
	if v757 < v756 {
		goto L105
	} else {
		goto L106
	}
L105:
	;
	v760 = v756
	goto L107
L106:
	;
	v760 = v757
	goto L107
L107:
	;
	v764 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v666+v760<<(uint(int32(1))%32)))))
	v766 = int32(base.Ui32(v764) >> (uint(int32(2)) % 32))
	v767 = m.G1
	v773 = int32(63)
	if v766 < v773 {
		goto L109
	} else {
		goto L110
	}
L108:
	;
	v780 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2572))
	v783 = base.I32_div_s(v778*v694, v780+int32(256))
	v784 = int32(63)
	if v783 < v784 {
		goto L112
	} else {
		goto L113
	}
L109:
	;
	v776 = v766
	goto L111
L110:
	;
	v776 = v773
	goto L111
L111:
	;
	v778 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v767+int32(_a_F_VP8SetSegmentParams_2)+v751<<(uint(int32(6))%32)+v776))))
	goto L108
L112:
	;
	v787 = v783
	goto L114
L113:
	;
	v787 = v784
	goto L114
L114:
	;
	if v783 < int32(2) {
		goto L115
	} else {
		goto L116
	}
L115:
	;
	v790 = int32(0)
	goto L117
L116:
	;
	v790 = v787
	goto L117
L117:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+2580)) = v790
	v792 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v793 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3320))
	v794 = int32(127)
	if v793 < v794 {
		goto L118
	} else {
		goto L119
	}
L118:
	;
	v797 = v793
	goto L120
L119:
	;
	v797 = v794
	goto L120
L120:
	;
	v798 = int32(0)
	if v798 < v797 {
		goto L121
	} else {
		goto L122
	}
L121:
	;
	v801 = v797
	goto L123
L122:
	;
	v801 = v798
	goto L123
L123:
	;
	v805 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v666+v801<<(uint(int32(1))%32)))))
	v807 = int32(base.Ui32(v805) >> (uint(int32(2)) % 32))
	v808 = m.G1
	v814 = int32(63)
	if v807 < v814 {
		goto L125
	} else {
		goto L126
	}
L124:
	;
	v820 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1092))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v820
	v822 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v823 = *(*int32)(unsafe.Add(mBase, uint32(v822)+36))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v823
	v825 = *(*int32)(unsafe.Add(mBase, uint32(v822)+40))
	v826 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = base.B2i32(v825 == v826)
	v831 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3316))
	v834 = base.I32_div_s(v819*v694, v831+int32(256))
	v835 = int32(63)
	if v834 < v835 {
		goto L128
	} else {
		goto L129
	}
L125:
	;
	v817 = v807
	goto L127
L126:
	;
	v817 = v814
	goto L127
L127:
	;
	v819 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v808+int32(_a_F_VP8SetSegmentParams_2)+v792<<(uint(int32(6))%32)+v817))))
	goto L124
L128:
	;
	v838 = v834
	goto L130
L129:
	;
	v838 = v835
	goto L130
L130:
	;
	if v834 < int32(2) {
		goto L131
	} else {
		goto L132
	}
L131:
	;
	v841 = v826
	goto L133
L132:
	;
	v841 = v838
	goto L133
L133:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3324)) = v841
	if v40 < int32(2) {
		goto L134
	} else {
		goto L135
	}
L134:
	;
	v1376 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[3])))
	if v1376 < int32(4) {
		v1381 = v625
		goto L188
	} else {
		goto L189
	}
L135:
	;
	v845 = m.G1
	v848 = int32(0)
	v849 = base.Simd_g_v128_load(m, v845+int32(_a_F_VP8SetSegmentParams_3), v848)
	base.Simd_g_v128_store(m, v33, v848, v849)
	v852 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v853 = int32(4)
	if v852 < v853 {
		goto L136
	} else {
		goto L137
	}
L136:
	;
	v856 = v852
	goto L138
L137:
	;
	v856 = v853
	goto L138
L138:
	;
	if v852 < int32(2) {
		v1045 = int32(1)
		goto L139
	} else {
		goto L140
	}
L139:
	;
	if v856 <= v1045 {
		goto L134
	} else {
		goto L170
	}
L140:
	;
	v861 = *(*int32)(unsafe.Add(mBase, uint32(v605)+1424))
	v874 = int32(0)
	v877 = v605 + int32(684)
	goto L143
L141:
	;
	if v852 == int32(2) {
		v1045 = v913
		goto L139
	} else {
		goto L149
	}
L142:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+4)) = v874
	v913 = int32(1)
	v914 = int32(0)
	goto L141
L143:
	;
	v895 = *(*int32)(unsafe.Add(mBase, uint32(v877+int32(-4))))
	if v861 != v895 {
		goto L145
	} else {
		goto L146
	}
L144:
	;
	v906 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+4)) = v906
	v913 = int32(2)
	v914 = v906
	goto L141
L145:
	;
	v902 = int32(1)
	v903 = v874 + v902
	if v903 != v902 {
		v874 = v903
		v877 = v877 + int32(744)
		goto L143
	} else {
		goto L148
	}
L146:
	;
	v897 = *(*int32)(unsafe.Add(mBase, uint32(v605)+1428))
	v898 = *(*int32)(unsafe.Add(mBase, uint32(v877)))
	if v897 == v898 {
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
	v923 = *(*int32)(unsafe.Add(mBase, uint32(v605)+2168))
	v936 = int32(0)
	v939 = v605 + int32(684)
	goto L152
L150:
	;
	if v852 == int32(3) {
		v1045 = v976
		goto L139
	} else {
		goto L160
	}
L151:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+8)) = v936
	v976 = v913
	goto L150
L152:
	;
	v957 = *(*int32)(unsafe.Add(mBase, uint32(v939+int32(-4))))
	if v923 != v957 {
		goto L154
	} else {
		goto L155
	}
L153:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+8)) = v913
	if v914 != 0 {
		goto L158
	} else {
		goto L159
	}
L154:
	;
	v965 = v936 + int32(1)
	if v913 != v965 {
		v936 = v965
		v939 = v939 + int32(744)
		goto L152
	} else {
		goto L157
	}
L155:
	;
	v959 = *(*int32)(unsafe.Add(mBase, uint32(v605)+2172))
	v960 = *(*int32)(unsafe.Add(mBase, uint32(v939)))
	if v959 == v960 {
		goto L151
	} else {
		goto L156
	}
L156:
	;
	goto L154
L157:
	;
	goto L153
L158:
	;
	v976 = v913 + int32(1)
	goto L150
L159:
	;
	v968 = int32(744)
	v972 = F_memcpy(m, v605+v913*v968, v605+int32(1488), v968)
	mBase = m.M
	goto L158
L160:
	;
	v985 = *(*int32)(unsafe.Add(mBase, uint32(v605)+2912))
	v998 = int32(0)
	v1001 = v605 + int32(684)
	goto L162
L161:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+12)) = v998
	v1045 = v976
	goto L139
L162:
	;
	v1019 = *(*int32)(unsafe.Add(mBase, uint32(v1001+int32(-4))))
	if v985 != v1019 {
		goto L164
	} else {
		goto L165
	}
L163:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+12)) = v976
	if v976 == int32(3) {
		goto L168
	} else {
		goto L169
	}
L164:
	;
	v1027 = v998 + int32(1)
	if v976 != v1027 {
		v998 = v1027
		v1001 = v1001 + int32(744)
		goto L162
	} else {
		goto L167
	}
L165:
	;
	v1021 = *(*int32)(unsafe.Add(mBase, uint32(v605)+2916))
	v1022 = *(*int32)(unsafe.Add(mBase, uint32(v1001)))
	if v1021 == v1022 {
		goto L161
	} else {
		goto L166
	}
L166:
	;
	goto L164
L167:
	;
	goto L163
L168:
	;
	v1045 = v976 + int32(1)
	goto L139
L169:
	;
	v1032 = int32(744)
	v1036 = F_memcpy(m, v605+v976*v1032, v605+int32(2232), v1032)
	mBase = m.M
	goto L168
L170:
	;
	v1071 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v1072 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v1073 = v1071 * v1072
	if v1073 < int32(1) {
		goto L171
	} else {
		goto L172
	}
L171:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v1045
	v1216 = v1045*int32(744) + v605
	v1218 = v1216 + int32(-744)
	v1221 = (v856 - v1045) & int32(3)
	if v1221 == int32(0) {
		v1267 = v1045
		goto L179
	} else {
		goto L180
	}
L172:
	;
	if v1073&int32(1) == int32(0) {
		v1102 = v1073
		goto L173
	} else {
		goto L174
	}
L173:
	;
	if v1073 == int32(1) {
		goto L171
	} else {
		goto L175
	}
L174:
	;
	v1080 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[4])))
	v1082 = v1073 + int32(-1)
	v1085 = v1080 + v1082<<(uint(int32(2))%32)
	v1086 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1085))))
	v1092 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33|int32(base.Ui32(v1086)>>(uint(int32(3))%32))&int32(12)))))
	v1099 = v1092<<(uint(int32(5))%32)&int32(96) | v1086&int32(159)
	*(*uint8)(unsafe.Add(mBase, uint32(v1085))) = uint8(v1099)
	v1102 = v1082
	goto L173
L175:
	;
	v1122 = v1102 + int32(1)
	v1125 = v1102<<(uint(int32(2))%32) + int32(-4)
	goto L176
L176:
	;
	v1141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[4])))
	v1142 = v1141 + v1125
	v1143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1142))))
	v1144 = int32(3)
	v1146 = int32(12)
	v1149 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33|int32(base.Ui32(v1143)>>(uint(v1144)%32))&v1146))))
	v1150 = int32(5)
	v1152 = int32(96)
	v1154 = int32(159)
	v1156 = v1149<<(uint(v1150)%32)&v1152 | v1143&v1154
	*(*uint8)(unsafe.Add(mBase, uint32(v1142))) = uint8(v1156)
	v1158 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8SetSegmentParams[4])))
	v1161 = v1158 + v1125 + int32(-4)
	v1162 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1161))))
	v1168 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v33|int32(base.Ui32(v1162)>>(uint(v1144)%32))&v1146))))
	v1175 = v1168<<(uint(v1150)%32)&v1152 | v1162&v1154
	*(*uint8)(unsafe.Add(mBase, uint32(v1161))) = uint8(v1175)
	v1180 = v1122 + int32(-2)
	if base.Ui32(int32(1)) < base.Ui32(v1180) {
		v1122 = v1180
		v1125 = v1125 + int32(-8)
		goto L176
	} else {
		goto L178
	}
L177:
	;
	goto L171
L178:
	;
	goto L177
L179:
	;
	if base.Ui32(int32(-4)) < base.Ui32(v1045-v856) {
		goto L134
	} else {
		goto L184
	}
L180:
	;
	v1236 = v1221
	v1243 = v1216
	goto L181
L181:
	;
	v1255 = int32(744)
	v1256 = F_memcpy(m, v1243, v1218, v1255)
	mBase = m.M
	v1260 = v1236 + int32(-1)
	if v1260 != 0 {
		v1236 = v1260
		v1243 = v1256 + v1255
		goto L181
	} else {
		goto L183
	}
L182:
	;
	v1267 = v1045 + v1221
	goto L179
L183:
	;
	goto L182
L184:
	;
	v1309 = v605 + v1267*int32(744)
	v1316 = v856 - v1267
	goto L185
L185:
	;
	v1328 = int32(744)
	v1329 = F_memcpy(m, v1309, v1218, v1328)
	mBase = m.M
	v1333 = F_memcpy(m, v1329+v1328, v1218, v1328)
	mBase = m.M
	v1337 = F_memcpy(m, v1329+int32(1488), v1218, v1328)
	mBase = m.M
	v1341 = F_memcpy(m, v1329+int32(2232), v1218, v1328)
	mBase = m.M
	v1345 = v1316 + int32(-4)
	if v1345 != 0 {
		v1309 = v1329 + int32(2976)
		v1316 = v1345
		goto L185
	} else {
		goto L187
	}
L186:
	;
	goto L134
L187:
	;
	goto L186
L188:
	;
	v1382 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if v1382 < int32(1) {
		goto L190
	} else {
		goto L191
	}
L189:
	;
	v1379 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1380 = *(*int32)(unsafe.Add(mBase, uint32(v1379)+28))
	v1381 = v1380
	goto L188
L190:
	;
	m.G0 = v33 + int32(16)
	return
L191:
	;
	v1390 = v1382
	v1395 = v605
	goto L192
L192:
	;
	v1417 = *(*int32)(unsafe.Add(mBase, uint32(v1395+int32(680))))
	v1420 = m.G1
	v1422 = v1420 + int32(_a_F_VP8SetSegmentParams_1)
	v1423 = int32(127)
	if v1417 < v1423 {
		goto L194
	} else {
		goto L195
	}
L193:
	;
	goto L190
L194:
	;
	v1426 = v1417
	goto L196
L195:
	;
	v1426 = v1423
	goto L196
L196:
	;
	v1427 = int32(0)
	if v1427 < v1426 {
		goto L197
	} else {
		goto L198
	}
L197:
	;
	v1430 = v1426
	goto L199
L198:
	;
	v1430 = v1427
	goto L199
L199:
	;
	v1434 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1422+v1430<<(uint(int32(1))%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(2)))) = uint16(v1434)
	v1437 = v1420 + int32(_a_F_VP8SetSegmentParams_4)
	v1438 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3396))
	v1439 = v1417 + v1438
	v1440 = int32(127)
	if v1439 < v1440 {
		goto L200
	} else {
		goto L201
	}
L200:
	;
	v1443 = v1439
	goto L202
L201:
	;
	v1443 = v1440
	goto L202
L202:
	;
	v1444 = int32(0)
	if v1444 < v1443 {
		goto L203
	} else {
		goto L204
	}
L203:
	;
	v1447 = v1443
	goto L205
L204:
	;
	v1447 = v1444
	goto L205
L205:
	;
	v1449 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1437+v1447))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395))) = uint16(v1449)
	v1452 = v1395 + int32(448)
	v1453 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3408))
	v1454 = v1417 + v1453
	v1455 = int32(117)
	if v1454 < v1455 {
		goto L206
	} else {
		goto L207
	}
L206:
	;
	v1458 = v1454
	goto L208
L207:
	;
	v1458 = v1455
	goto L208
L208:
	;
	v1459 = int32(0)
	if v1459 < v1458 {
		goto L209
	} else {
		goto L210
	}
L209:
	;
	v1462 = v1458
	goto L211
L210:
	;
	v1462 = v1459
	goto L211
L211:
	;
	v1464 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1437+v1462))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1452))) = uint16(v1464)
	v1467 = v1395 + int32(224)
	v1468 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3400))
	v1469 = v1417 + v1468
	v1470 = int32(127)
	if v1469 < v1470 {
		goto L212
	} else {
		goto L213
	}
L212:
	;
	v1473 = v1469
	goto L214
L213:
	;
	v1473 = v1470
	goto L214
L214:
	;
	v1474 = int32(0)
	if v1474 < v1473 {
		goto L215
	} else {
		goto L216
	}
L215:
	;
	v1477 = v1473
	goto L217
L216:
	;
	v1477 = v1474
	goto L217
L217:
	;
	v1479 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1437+v1477))))
	v1481 = v1479 << (uint(int32(1)) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467))) = uint16(v1481)
	v1487 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3404))
	v1488 = v1417 + v1487
	v1489 = int32(127)
	if v1488 < v1489 {
		goto L218
	} else {
		goto L219
	}
L218:
	;
	v1492 = v1488
	goto L220
L219:
	;
	v1492 = v1489
	goto L220
L220:
	;
	v1493 = int32(0)
	if v1493 < v1492 {
		goto L221
	} else {
		goto L222
	}
L221:
	;
	v1496 = v1492
	goto L223
L222:
	;
	v1496 = v1493
	goto L223
L223:
	;
	v1500 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1420+int32(_a_F_VP8SetSegmentParams_5)+v1496<<(uint(int32(1))%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(226)))) = uint16(v1500)
	v1504 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3412))
	v1505 = v1417 + v1504
	v1506 = int32(127)
	if v1505 < v1506 {
		goto L224
	} else {
		goto L225
	}
L224:
	;
	v1509 = v1505
	goto L226
L225:
	;
	v1509 = v1506
	goto L226
L226:
	;
	v1510 = int32(0)
	if v1510 < v1509 {
		goto L227
	} else {
		goto L228
	}
L227:
	;
	v1513 = v1509
	goto L229
L228:
	;
	v1513 = v1510
	goto L229
L229:
	;
	v1514 = int32(1)
	v1517 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1422+v1513<<(uint(v1514)%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(450)))) = uint16(v1517)
	v1519 = int32(0)
	v1526 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+4)) = uint16(v1526)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+6)) = uint16(v1526)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+8)) = uint16(v1526)
	v1530 = int32(131072)
	v1531 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395))))
	v1532 = base.I32_div_u_s(v1530, v1531)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+32)) = uint16(v1532)
	v1535 = base.I32_div_u_s(v1530, v1526)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+34)) = uint16(v1535)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+36)) = uint16(v1535)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+38)) = uint16(v1535)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+40)) = uint16(v1535)
	v1540 = m.G1
	v1546 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1540)+uint32(_c_F_VP8SetSegmentParams[5]))))
	v1547 = int32(9)
	v1548 = v1546 << (uint(v1547) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+64)) = v1548
	v1550 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1540)+uint32(_c_F_VP8SetSegmentParams[6]))))
	v1552 = v1550 << (uint(v1547) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+68)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+72)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+76)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+80)) = v1552
	v1557 = int32(131071)
	v1559 = int32(_a_F_VP8SetSegmentParams_6)
	v1561 = base.I32_div_u_s(v1548^v1557, v1532&v1559)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+128)) = v1561
	v1567 = base.I32_div_u_s(v1552^v1557, v1535&v1559)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+132)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+136)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+140)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+144)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+148)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+84)) = v1552
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+42)) = uint16(v1535)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+10)) = uint16(v1526)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+88)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+152)) = v1567
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+44)) = uint16(v1535)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+92)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+156)) = v1567
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+46)) = uint16(v1535)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+96)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+160)) = v1567
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+48)) = uint16(v1535)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+100)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+164)) = v1567
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+50)) = uint16(v1535)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+52)) = uint16(v1535)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+104)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+168)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+172)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+108)) = v1552
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+54)) = uint16(v1535)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+56)) = uint16(v1535)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+112)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+176)) = v1567
	v1597 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+58)) = uint16(v1597)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+116)) = v1552
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+180)) = v1567
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+60)) = uint16(v1597)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+28)) = uint16(v1526)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+184)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+120)) = v1552
	v1612 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_replace_lane_l1(base.Simd_g_i16x8_splat(v1526), v1526), base.Simd_g_const(&F_VP8SetSegmentParams__k4))
	base.Simd_g_v128_store(m, v1395, int32(12), v1612)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(192)))) = uint16(v1519)
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+188)) = v1567
	*(*int32)(unsafe.Add(mBase, uint32(v1395)+124)) = v1552
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+62)) = uint16(v1597)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395)+30)) = uint16(v1526)
	v1623 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395))))
	goto L232
L230:
	;
	v1758 = int32(0)
	v1764 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+4)) = uint16(v1764)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+6)) = uint16(v1764)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+8)) = uint16(v1764)
	v1768 = int32(131072)
	v1769 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467))))
	v1770 = base.I32_div_u_s(v1768, v1769)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+32)) = uint16(v1770)
	v1773 = base.I32_div_u_s(v1768, v1764)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+34)) = uint16(v1773)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+36)) = uint16(v1773)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+38)) = uint16(v1773)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+40)) = uint16(v1773)
	v1778 = m.G1
	v1784 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1778)+uint32(_c_F_VP8SetSegmentParams[7]))))
	v1785 = int32(9)
	v1786 = v1784 << (uint(v1785) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+64)) = v1786
	v1788 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1778)+uint32(_c_F_VP8SetSegmentParams[8]))))
	v1790 = v1788 << (uint(v1785) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+68)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+72)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+76)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+80)) = v1790
	v1795 = int32(131071)
	v1797 = int32(_a_F_VP8SetSegmentParams_6)
	v1799 = base.I32_div_u_s(v1786^v1795, v1770&v1797)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+128)) = v1799
	v1805 = base.I32_div_u_s(v1790^v1795, v1773&v1797)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+132)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+136)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+140)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+144)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+148)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+84)) = v1790
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+42)) = uint16(v1773)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+10)) = uint16(v1764)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+88)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+152)) = v1805
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+44)) = uint16(v1773)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+92)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+156)) = v1805
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+46)) = uint16(v1773)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+96)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+160)) = v1805
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+48)) = uint16(v1773)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+100)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+164)) = v1805
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+50)) = uint16(v1773)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+52)) = uint16(v1773)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+104)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+168)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+172)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+108)) = v1790
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+54)) = uint16(v1773)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+56)) = uint16(v1773)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+112)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+176)) = v1805
	v1835 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+58)) = uint16(v1835)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+116)) = v1790
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+180)) = v1805
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+60)) = uint16(v1835)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+28)) = uint16(v1764)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+184)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+120)) = v1790
	v1850 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_replace_lane_l1(base.Simd_g_i16x8_splat(v1764), v1764), base.Simd_g_const(&F_VP8SetSegmentParams__k4))
	base.Simd_g_v128_store(m, v1467, int32(12), v1850)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(416)))) = uint16(v1758)
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+188)) = v1805
	*(*int32)(unsafe.Add(mBase, uint32(v1467)+124)) = v1790
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+62)) = uint16(v1835)
	*(*uint16)(unsafe.Add(mBase, uint32(v1467)+30)) = uint16(v1764)
	v1861 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467))))
	goto L237
L231:
	;
	v1756 = int32(base.Ui32(v1728+(v1720+(v1712+v1735))+v1735+v1526<<(uint(int32(2))%32)+((v1526+v1526)<<(uint(v1734)%32)+v1623)+int32(8)) >> (uint(int32(4)) % 32))
	goto L230
L232:
	;
	v1660 = int32(30)
	v1662 = int32(11)
	v1663 = int32(base.Ui32(v1526*v1660) >> (uint(v1662) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(194)))) = uint16(v1663)
	v1667 = int32(90)
	v1670 = int32(base.Ui32(v1526*v1667) >> (uint(v1662) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(222)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(220)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(218)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(216)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(214)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(212)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(210)))) = uint16(v1670)
	v1696 = int32(60)
	v1699 = int32(base.Ui32(v1526*v1696) >> (uint(v1662) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(208)))) = uint16(v1699)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(206)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(204)))) = uint16(v1670)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(202)))) = uint16(v1699)
	v1712 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395)+4)))
	v1716 = int32(base.Ui32(v1712*v1696) >> (uint(v1662) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(196)))) = uint16(v1716)
	v1720 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395)+6)))
	v1724 = int32(base.Ui32(v1720*v1667) >> (uint(v1662) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(198)))) = uint16(v1724)
	v1728 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395)+8)))
	v1732 = int32(base.Ui32(v1728*v1660) >> (uint(v1662) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(200)))) = uint16(v1732)
	v1734 = int32(1)
	v1735 = v1526 << (uint(v1734) % 32)
	goto L231
L234:
	;
	v1996 = int32(0)
	v2002 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+4)) = uint16(v2002)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+6)) = uint16(v2002)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+8)) = uint16(v2002)
	v2006 = int32(131072)
	v2007 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452))))
	v2008 = base.I32_div_u_s(v2006, v2007)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+32)) = uint16(v2008)
	v2011 = base.I32_div_u_s(v2006, v2002)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+34)) = uint16(v2011)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+36)) = uint16(v2011)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+38)) = uint16(v2011)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+40)) = uint16(v2011)
	v2016 = m.G1
	v2022 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2016)+uint32(_c_F_VP8SetSegmentParams[9]))))
	v2023 = int32(9)
	v2024 = v2022 << (uint(v2023) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+64)) = v2024
	v2026 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2016)+uint32(_c_F_VP8SetSegmentParams[10]))))
	v2028 = v2026 << (uint(v2023) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+68)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+72)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+76)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+80)) = v2028
	v2033 = int32(131071)
	v2035 = int32(_a_F_VP8SetSegmentParams_6)
	v2037 = base.I32_div_u_s(v2024^v2033, v2008&v2035)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+128)) = v2037
	v2043 = base.I32_div_u_s(v2028^v2033, v2011&v2035)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+132)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+136)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+140)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+144)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+148)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+84)) = v2028
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+42)) = uint16(v2011)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+10)) = uint16(v2002)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+88)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+152)) = v2043
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+44)) = uint16(v2011)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+92)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+156)) = v2043
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+46)) = uint16(v2011)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+96)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+160)) = v2043
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+48)) = uint16(v2011)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+100)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+164)) = v2043
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+50)) = uint16(v2011)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+52)) = uint16(v2011)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+104)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+168)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+172)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+108)) = v2028
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+54)) = uint16(v2011)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+56)) = uint16(v2011)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+112)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+176)) = v2043
	v2073 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+58)) = uint16(v2073)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+116)) = v2028
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+180)) = v2043
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+60)) = uint16(v2073)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+28)) = uint16(v2002)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+184)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+120)) = v2028
	v2088 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_replace_lane_l1(base.Simd_g_i16x8_splat(v2002), v2002), base.Simd_g_const(&F_VP8SetSegmentParams__k4))
	base.Simd_g_v128_store(m, v1452, int32(12), v2088)
	*(*uint16)(unsafe.Add(mBase, uint32(v1395+int32(640)))) = uint16(v1996)
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+188)) = v2043
	*(*int32)(unsafe.Add(mBase, uint32(v1452)+124)) = v2028
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+62)) = uint16(v2073)
	*(*uint16)(unsafe.Add(mBase, uint32(v1452)+30)) = uint16(v2002)
	v2099 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452))))
	goto L241
L235:
	;
	v1994 = int32(base.Ui32(v1878+v1879+v1881+v1883+v1878+v1764<<(uint(int32(2))%32)+((v1764+v1764)<<(uint(v1877)%32)+v1861)+int32(8)) >> (uint(int32(4)) % 32))
	goto L234
L237:
	;
	v1866 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v1395+int32(418)))) = v1866
	v1870 = base.Simd_g_const(&F_VP8SetSegmentParams__k3)
	base.Simd_g_v128_store(m, v1395+int32(426), int32(0), v1870)
	*(*int64)(unsafe.Add(mBase, uint32(v1395+int32(440)))) = v1866
	v1877 = int32(1)
	v1878 = v1764 << (uint(v1877) % 32)
	v1879 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467)+4)))
	v1881 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467)+6)))
	v1883 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1467)+8)))
	goto L235
L238:
	;
	v2234 = v1395 + int32(708)
	v2235 = v1756 * v1756
	v2236 = int32(7)
	*(*int32)(unsafe.Add(mBase, uint32(v2234))) = int32(base.Ui32(v2235) >> (uint(v2236) % 32))
	v2240 = v1395 + int32(716)
	v2243 = v1756 * v1381 >> (uint(int32(5)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v2240))) = v2243
	v2246 = v1395 + int32(724)
	v2249 = int32(3)
	v2250 = v2235 * v2236 >> (uint(v2249) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v2246))) = v2250
	v2253 = v1395 + int32(696)
	v2254 = v1994 * v1994
	v2256 = v2254 * v2249
	*(*int32)(unsafe.Add(mBase, uint32(v2253))) = v2256
	v2259 = v1395 + int32(700)
	v2263 = v2235 * v2249 >> (uint(v2236) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v2259))) = v2263
	v2266 = v1395 + int32(720)
	*(*int32)(unsafe.Add(mBase, uint32(v2266))) = int32(base.Ui32(v2254) >> (uint(int32(2)) % 32))
	v2271 = v1395 + int32(728)
	v2272 = v2232 * v2232
	v2274 = v2272 << (uint(int32(1)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v2271))) = v2274
	v2277 = v1395 + int32(704)
	v2281 = v2272 * v2249 >> (uint(int32(6)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v2277))) = v2281
	if int32(0) < v2263 {
		goto L242
	} else {
		goto L243
	}
L239:
	;
	v2232 = int32(base.Ui32(v2116+v2117+v2119+v2121+v2116+v2002<<(uint(int32(2))%32)+((v2002+v2002)<<(uint(v2115)%32)+v2099)+int32(8)) >> (uint(int32(4)) % 32))
	goto L238
L241:
	;
	v2104 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v1395+int32(642)))) = v2104
	v2108 = base.Simd_g_const(&F_VP8SetSegmentParams__k3)
	base.Simd_g_v128_store(m, v1395+int32(650), int32(0), v2108)
	*(*int64)(unsafe.Add(mBase, uint32(v1395+int32(664)))) = v2104
	v2115 = int32(1)
	v2116 = v2002 << (uint(v2115) % 32)
	v2117 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452)+4)))
	v2119 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452)+6)))
	v2121 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1452)+8)))
	goto L239
L242:
	;
	if int32(0) < v2256 {
		goto L244
	} else {
		goto L245
	}
L243:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2259))) = int32(1)
	goto L242
L244:
	;
	if int32(0) < v2281 {
		goto L246
	} else {
		goto L247
	}
L245:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2253))) = int32(1)
	goto L244
L246:
	;
	if base.Ui32(int32(127)) < base.Ui32(v2235) {
		goto L248
	} else {
		goto L249
	}
L247:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2277))) = int32(1)
	goto L246
L248:
	;
	if int32(0) < v2250 {
		goto L250
	} else {
		goto L251
	}
L249:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2234))) = int32(1)
	goto L248
L250:
	;
	if base.Ui32(int32(3)) < base.Ui32(v2254) {
		goto L252
	} else {
		goto L253
	}
L251:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2246))) = int32(1)
	goto L250
L252:
	;
	if int32(0) < v2274 {
		goto L254
	} else {
		goto L255
	}
L253:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2266))) = int32(1)
	goto L252
L254:
	;
	if int32(0) < v2243 {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2271))) = int32(1)
	goto L254
L256:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1395+int32(688)))) = int32(0)
	v2321 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1395))))
	*(*int32)(unsafe.Add(mBase, uint32(v1395+int32(692)))) = v2321 * int32(20)
	*(*int64)(unsafe.Add(mBase, uint32(v1395+int32(736)))) = base.I64_extend_i32_s(v2235 * int32(1000))
	v2334 = v1390 + int32(-1)
	if v2334 != 0 {
		v1390 = v2334
		v1395 = v1395 + int32(744)
		goto L192
	} else {
		goto L258
	}
L257:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2240))) = int32(1)
	goto L256
L258:
	;
	goto L193
}

var F_VP8SetSegmentParams__k0 = [2]uint64{0x3ff0000000000000, 0x3ff0000000000000}
var F_VP8SetSegmentParams__k1 = [2]uint64{0x405fc00000000000, 0x405fc00000000000}
var F_VP8SetSegmentParams__k2 = [2]uint64{0x7f0000007f, 0x7f0000007f}
var F_VP8SetSegmentParams__k3 = [2]uint64{0x0, 0x0}
var F_VP8SetSegmentParams__k4 = [2]uint64{0x100010001000100, 0x302010001000100}

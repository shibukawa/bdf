//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_EncodeAlphaInternal(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) int32 {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v34 base.V128
	_ = v34
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v80 int32
	_ = v80
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v99 int64
	_ = v99
	var v107 base.V128
	_ = v107
	var v108 int32
	_ = v108
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v140 int64
	_ = v140
	var v145 base.V128
	_ = v145
	var v148 base.V128
	_ = v148
	var v151 base.V128
	_ = v151
	var v154 int32
	_ = v154
	var v162 base.V128
	_ = v162
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v200 float32
	_ = v200
	var v204 float32
	_ = v204
	var v205 float32
	_ = v205
	var v210 int32
	_ = v210
	var v213 int32
	_ = v213
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v222 int64
	_ = v222
	var v230 base.V128
	_ = v230
	var v231 int32
	_ = v231
	var v239 int32
	_ = v239
	var v243 int32
	_ = v243
	var v248 int32
	_ = v248
	var v254 base.V128
	_ = v254
	var v261 int32
	_ = v261
	var v266 int32
	_ = v266
	var v272 base.V128
	_ = v272
	var v275 base.V128
	_ = v275
	var v280 int32
	_ = v280
	var v284 int32
	_ = v284
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v297 int32
	_ = v297
	var v303 int64
	_ = v303
	var v308 int32
	_ = v308
	var v309 int32
	_ = v309
	var v310 int32
	_ = v310
	var v311 int32
	_ = v311
	var v318 int32
	_ = v318
	var v320 int32
	_ = v320
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v332 int32
	_ = v332
	var v339 int32
	_ = v339
	var v340 int32
	_ = v340
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v354 int32
	_ = v354
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v363 int32
	_ = v363
	var v380 int32
	_ = v380
	var v390 int32
	_ = v390
	var v391 int32
	_ = v391
	var v395 int32
	_ = v395
	var v400 int32
	_ = v400
	var v406 base.V128
	_ = v406
	var v409 base.V128
	_ = v409
	var v414 int32
	_ = v414
	var v417 int32
	_ = v417
	var v422 int32
	_ = v422
	var v423 int32
	_ = v423
	var v425 int32
	_ = v425
	var v430 int32
	_ = v430
	var v435 int32
	_ = v435
	var v441 base.V128
	_ = v441
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v452 int32
	_ = v452
	var v453 int32
	_ = v453
	var v457 int32
	_ = v457
	var v460 int32
	_ = v460
	var v463 int32
	_ = v463
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v468 base.V128
	_ = v468
	var v471 base.V128
	_ = v471
	var v477 int32
	_ = v477
	var v480 int32
	_ = v480
	var v481 int32
	_ = v481
	var v485 int32
	_ = v485
	var v490 int32
	_ = v490
	var v491 int32
	_ = v491
	var v492 int32
	_ = v492
	var v503 int32
	_ = v503
	var v508 int32
	_ = v508
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v520 int64
	_ = v520
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v529 int32
	_ = v529
	var v532 int32
	_ = v532
	var v534 int32
	_ = v534
	var v535 int32
	_ = v535
	var v538 int32
	_ = v538
	var v539 int32
	_ = v539
	var v543 int32
	_ = v543
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v555 int32
	_ = v555
	var v557 int32
	_ = v557
	var v558 int32
	_ = v558
	var v562 int32
	_ = v562
	var v570 int32
	_ = v570
	var v578 int32
	_ = v578
	var v581 int32
	_ = v581
	var v584 int64
	_ = v584
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v593 int32
	_ = v593
	var v596 int32
	_ = v596
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v602 int32
	_ = v602
	var v603 int32
	_ = v603
	var v607 int32
	_ = v607
	var v610 int32
	_ = v610
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v616 int32
	_ = v616
	var v617 int32
	_ = v617
	var v619 int32
	_ = v619
	var v621 int32
	_ = v621
	var v622 int32
	_ = v622
	var v626 int32
	_ = v626
	var v634 int32
	_ = v634
	var v637 int32
	_ = v637
	var v639 int32
	_ = v639
	var v644 int32
	_ = v644
	var v650 base.V128
	_ = v650
	var v653 int32
	_ = v653
	var v656 int32
	_ = v656
	var v659 int32
	_ = v659
	var v662 int32
	_ = v662
	var v664 int32
	_ = v664
	v14 = m.G0
	v16 = v14 - int32(320)
	m.G0 = v16
	v18 = m.G92
	v22 = *(*int32)(unsafe.Add(mBase, uint32(v18+l4<<(uint(int32(2))%32))))
	if v22 != 0 {
		m.T0[v22].(func(*base.Module, int32, int32, int32, int32, int32))(m, l0, l1, l2, l1, l7)
		mBase = m.M
		v24 = l7
	} else {
		v24 = l0
	}
	v25 = l2 * l1
	v26 = int32(0)
	if l3 != 0 {
		v30 = v16 + int32(8)
		v34 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
		base.Simd_g_v128_store(m, v30, int32(0), v34)
		*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
		v45 = int32(base.Ui32(v25)>>(uint(int32(3))%32))&int32(-1024) + int32(1024)
		v46 = F_WebPSafeMalloc(m, int64(1), v45)
		mBase = m.M
		if v46 != 0 {
			v50 = *(*int32)(unsafe.Add(mBase, uint32(v30)+8))
			F_WebPSafeFree(m, v50)
			mBase = m.M
			*(*int32)(unsafe.Add(mBase, uint32(v30)+16)) = v46 + v45
			*(*int32)(unsafe.Add(mBase, uint32(v30)+12)) = v46
			*(*int32)(unsafe.Add(mBase, uint32(v30)+8)) = v46
			v57 = int32(1)
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(v30)+20)) = int32(1)
			v57 = int32(0)
		}
		if v57 == int32(0) {
			v261 = v16 + int32(8)
			if v261 == int32(0) {
			} else {
				v266 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
				F_WebPSafeFree(m, v266)
				mBase = m.M
				*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
				v272 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
				base.Simd_g_v128_store(m, v261, int32(0), v272)
			}
			v275 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
			base.Simd_g_v128_store(m, l8, int32(4), v275)
			v280 = int32(0)
			base.Simd_g_v128_store(m, l8+int32(20), v280, v275)
			v664 = v280
		} else {
			v61 = v16 + int32(32)
			if v61 == int32(0) {
			} else {
				v72 = int32(0)
				v74 = F_memset(m, v61, v72, int32(172))
				mBase = m.M
				v75 = m.G2
				*(*int32)(unsafe.Add(mBase, uint32(v74)+72)) = v75 + int32(1)
				v80 = F_WebPEncodingSetError(m, v74, v72)
				mBase = m.M
			}
			*(*int32)(unsafe.Add(mBase, uint32(v16)+44)) = l2
			*(*int32)(unsafe.Add(mBase, uint32(v16)+40)) = l1
			*(*int32)(unsafe.Add(mBase, uint32(v16)+32)) = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v16)+120)) = l8 + int32(36)
			v93 = v16 + int32(32)
			if v93 != 0 {
				v95 = *(*int32)(unsafe.Add(mBase, uint32(v93)+156))
				F_WebPSafeFree(m, v95)
				mBase = m.M
				v97 = *(*int32)(unsafe.Add(mBase, uint32(v93)+160))
				F_WebPSafeFree(m, v97)
				mBase = m.M
				v99 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(v93)+52)) = v99
				*(*int64)(unsafe.Add(mBase, uint32(v93)+156)) = v99
				*(*int64)(unsafe.Add(mBase, uint32(v93)+16)) = v99
				v107 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
				v108 = int32(0)
				base.Simd_g_v128_store(m, v16+int32(56), v108, v107)
				*(*int32)(unsafe.Add(mBase, uint32(v16+int32(72)))) = v108
				v114 = *(*int32)(unsafe.Add(mBase, uint32(v93)))
				if v114 != 0 {
					v116 = F_WebPPictureAllocARGB(m, v93)
					mBase = m.M
					v117 = v116
				} else {
					v115 = F_WebPPictureAllocYUVA(m, v93)
					mBase = m.M
					v117 = v115
				}
			} else {
				v117 = int32(1)
			}
			if v117 == int32(0) {
				v261 = v16 + int32(8)
				if v261 == int32(0) {
				} else {
					v266 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
					F_WebPSafeFree(m, v266)
					mBase = m.M
					*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
					v272 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
					base.Simd_g_v128_store(m, v261, int32(0), v272)
				}
				v275 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
				base.Simd_g_v128_store(m, l8, int32(4), v275)
				v280 = int32(0)
				base.Simd_g_v128_store(m, l8+int32(20), v280, v275)
				v664 = v280
			} else {
				v120 = m.G7
				v121 = *(*int32)(unsafe.Add(mBase, uint32(v16)+40))
				v122 = *(*int32)(unsafe.Add(mBase, uint32(v16)+44))
				v123 = *(*int32)(unsafe.Add(mBase, uint32(v16)+84))
				v124 = *(*int32)(unsafe.Add(mBase, uint32(v16)+88))
				v125 = *(*int32)(unsafe.Add(mBase, uint32(v120)))
				m.T0[v125].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v24, l1, v121, v122, v123, v124)
				mBase = m.M
				v128 = v16 + int32(204)
				v129 = int32(0)
				if v128 == v129 {
					v190 = v129
				} else {
					v140 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v128)+16)) = v140
					*(*float32)(unsafe.Add(mBase, uint32(v128)+4)) = float32(75)
					*(*int32)(unsafe.Add(mBase, uint32(v128)+112)) = int32(100)
					v145 = base.Simd_g_const(&F_EncodeAlphaInternal__k1)
					base.Simd_g_v128_store(m, v128, int32(24), v145)
					v148 = base.Simd_g_const(&F_EncodeAlphaInternal__k2)
					base.Simd_g_v128_store(m, v128, int32(56), v148)
					v151 = base.Simd_g_const(&F_EncodeAlphaInternal__k3)
					base.Simd_g_v128_store(m, v128, int32(40), v151)
					v154 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(v128)+96)) = v154
					*(*int32)(unsafe.Add(mBase, uint32(v128))) = v154
					*(*int64)(unsafe.Add(mBase, uint32(v128)+104)) = v140
					*(*int64)(unsafe.Add(mBase, uint32(v128)+88)) = int64(429496729600)
					v162 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
					base.Simd_g_v128_store(m, v128, int32(72), v162)
					*(*int64)(unsafe.Add(mBase, uint32(v128)+8)) = int64(4)
					switch int32(-1) {
					case 0:
						*(*int32)(unsafe.Add(mBase, uint32(v128)+36)) = int32(4)
						*(*int64)(unsafe.Add(mBase, uint32(v128)+28)) = int64(150323855440)
					case 1:
						*(*int32)(unsafe.Add(mBase, uint32(v128)+36)) = int32(3)
						*(*int64)(unsafe.Add(mBase, uint32(v128)+28)) = int64(128849018960)
						*(*int32)(unsafe.Add(mBase, uint32(v128)+68)) = int32(2)
					case 2:
						*(*int32)(unsafe.Add(mBase, uint32(v128)+36)) = int32(6)
						*(*int64)(unsafe.Add(mBase, uint32(v128)+28)) = int64(42949672985)
					case 3:
						*(*int64)(unsafe.Add(mBase, uint32(v128)+28)) = int64(0)
					case 4:
						*(*int32)(unsafe.Add(mBase, uint32(v128)+32)) = int32(0)
						*(*int64)(unsafe.Add(mBase, uint32(v128)+24)) = int64(2)
					default:
					}
					v189 = F_WebPValidateConfig(m, v128)
					mBase = m.M
					v190 = v189
				}
				if v190 == int32(0) {
					v261 = v16 + int32(8)
					if v261 == int32(0) {
					} else {
						v266 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
						F_WebPSafeFree(m, v266)
						mBase = m.M
						*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
						v272 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
						base.Simd_g_v128_store(m, v261, int32(0), v272)
					}
					v275 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
					base.Simd_g_v128_store(m, l8, int32(4), v275)
					v280 = int32(0)
					base.Simd_g_v128_store(m, l8+int32(20), v280, v275)
					v664 = v280
				} else {
					v193 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v16)+300)) = v193
					*(*int32)(unsafe.Add(mBase, uint32(v16)+204)) = v193
					*(*int32)(unsafe.Add(mBase, uint32(v16)+212)) = l6
					v200 = base.F32_mul(base.F32_convert_i32_s(l6), float32(8))
					if l6 == int32(6) {
						v204 = float32(100)
					} else {
						v204 = v200
					}
					if l5 != 0 {
						v205 = v200
					} else {
						v205 = v204
					}
					*(*float32)(unsafe.Add(mBase, uint32(v16)+208)) = v205
					v210 = v16 + int32(32)
					v213 = F_VP8LEncodeStream(m, v16+int32(204), v210, v16+int32(8))
					mBase = m.M
					if v210 == int32(0) {
					} else {
						v218 = *(*int32)(unsafe.Add(mBase, uint32(v210)+156))
						F_WebPSafeFree(m, v218)
						mBase = m.M
						v220 = *(*int32)(unsafe.Add(mBase, uint32(v210)+160))
						F_WebPSafeFree(m, v220)
						mBase = m.M
						v222 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(v210)+52)) = v222
						*(*int64)(unsafe.Add(mBase, uint32(v210)+156)) = v222
						*(*int64)(unsafe.Add(mBase, uint32(v210)+16)) = v222
						v230 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
						v231 = int32(0)
						base.Simd_g_v128_store(m, v16+int32(56), v231, v230)
						*(*int32)(unsafe.Add(mBase, uint32(v16+int32(72)))) = v231
					}
					if v213 == int32(0) {
						v243 = v16 + int32(8)
						if v243 == int32(0) {
						} else {
							v248 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
							F_WebPSafeFree(m, v248)
							mBase = m.M
							*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
							v254 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
							base.Simd_g_v128_store(m, v243, int32(0), v254)
						}
						v261 = v16 + int32(8)
						if v261 == int32(0) {
						} else {
							v266 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
							F_WebPSafeFree(m, v266)
							mBase = m.M
							*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
							v272 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
							base.Simd_g_v128_store(m, v261, int32(0), v272)
						}
						v275 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
						base.Simd_g_v128_store(m, l8, int32(4), v275)
						v280 = int32(0)
						base.Simd_g_v128_store(m, l8+int32(20), v280, v275)
						v664 = v280
					} else {
						v239 = *(*int32)(unsafe.Add(mBase, uint32(v16)+28))
						if v239 == int32(0) {
							v284 = v16 + int32(8)
							v293 = *(*int32)(unsafe.Add(mBase, uint32(v284)+12))
							v294 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
							v295 = v293 - v294
							v297 = *(*int32)(unsafe.Add(mBase, uint32(v284)+4))
							v303 = base.I64_extend_i32_u(v295) + base.I64_extend_i32_u((v297+int32(7))>>(uint(int32(3))%32))
							if base.Ui64(v303) < base.Ui64(int64(4294967296)) {
								v309 = base.I32_wrap_i64(v303)
								v310 = *(*int32)(unsafe.Add(mBase, uint32(v284)+16))
								v311 = v310 - v294
								if v310 == v294 {
									v318 = int32(base.Ui32(v311*int32(3)) >> (uint(int32(1)) % 32))
									if base.Ui32(v309) < base.Ui32(v318) {
										v320 = v318
									} else {
										v320 = v309
									}
									v324 = v320&int32(-1024) + int32(1024)
									v325 = F_WebPSafeMalloc(m, int64(1), v324)
									mBase = m.M
									if v325 != 0 {
										if v293 == v294 {
										} else {
											v330 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
											v331 = F_memcpy(m, v325, v330, v295)
											mBase = m.M
										}
										v332 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
										F_WebPSafeFree(m, v332)
										mBase = m.M
										*(*int32)(unsafe.Add(mBase, uint32(v284)+16)) = v325 + v324
										*(*int32)(unsafe.Add(mBase, uint32(v284)+12)) = v325 + v295
										*(*int32)(unsafe.Add(mBase, uint32(v284)+8)) = v325
										v339 = *(*int32)(unsafe.Add(mBase, uint32(v284)+4))
										v340 = v339
										if v340 < int32(1) {
										} else {
											v344 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
											v347 = v344
											for {
												v354 = *(*int32)(unsafe.Add(mBase, uint32(v284)+12))
												*(*int32)(unsafe.Add(mBase, uint32(v284)+12)) = v354 + int32(1)
												*(*uint8)(unsafe.Add(mBase, uint32(v354))) = uint8(v347)
												v359 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
												v360 = int32(8)
												v361 = int32(base.Ui32(v359) >> (uint(v360) % 32))
												*(*int32)(unsafe.Add(mBase, uint32(v284))) = v361
												v363 = *(*int32)(unsafe.Add(mBase, uint32(v284)+4))
												*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = v363 + int32(-8)
												if v360 < v363 {
													v347 = v361
													continue
												} else {
													break
												}
												break
											}
										}
										*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = int32(0)
										v380 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
										v390 = v380
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(v284)+20)) = int32(1)
										v328 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
										v390 = v328
									}
								} else {
									if base.Ui32(v309) <= base.Ui32(v311) {
										v340 = v297
										if v340 < int32(1) {
										} else {
											v344 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
											v347 = v344
											for {
												v354 = *(*int32)(unsafe.Add(mBase, uint32(v284)+12))
												*(*int32)(unsafe.Add(mBase, uint32(v284)+12)) = v354 + int32(1)
												*(*uint8)(unsafe.Add(mBase, uint32(v354))) = uint8(v347)
												v359 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
												v360 = int32(8)
												v361 = int32(base.Ui32(v359) >> (uint(v360) % 32))
												*(*int32)(unsafe.Add(mBase, uint32(v284))) = v361
												v363 = *(*int32)(unsafe.Add(mBase, uint32(v284)+4))
												*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = v363 + int32(-8)
												if v360 < v363 {
													v347 = v361
													continue
												} else {
													break
												}
												break
											}
										}
										*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = int32(0)
										v380 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
										v390 = v380
									} else {
										v318 = int32(base.Ui32(v311*int32(3)) >> (uint(int32(1)) % 32))
										if base.Ui32(v309) < base.Ui32(v318) {
											v320 = v318
										} else {
											v320 = v309
										}
										v324 = v320&int32(-1024) + int32(1024)
										v325 = F_WebPSafeMalloc(m, int64(1), v324)
										mBase = m.M
										if v325 != 0 {
											if v293 == v294 {
											} else {
												v330 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
												v331 = F_memcpy(m, v325, v330, v295)
												mBase = m.M
											}
											v332 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
											F_WebPSafeFree(m, v332)
											mBase = m.M
											*(*int32)(unsafe.Add(mBase, uint32(v284)+16)) = v325 + v324
											*(*int32)(unsafe.Add(mBase, uint32(v284)+12)) = v325 + v295
											*(*int32)(unsafe.Add(mBase, uint32(v284)+8)) = v325
											v339 = *(*int32)(unsafe.Add(mBase, uint32(v284)+4))
											v340 = v339
											if v340 < int32(1) {
											} else {
												v344 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
												v347 = v344
												for {
													v354 = *(*int32)(unsafe.Add(mBase, uint32(v284)+12))
													*(*int32)(unsafe.Add(mBase, uint32(v284)+12)) = v354 + int32(1)
													*(*uint8)(unsafe.Add(mBase, uint32(v354))) = uint8(v347)
													v359 = *(*int32)(unsafe.Add(mBase, uint32(v284)))
													v360 = int32(8)
													v361 = int32(base.Ui32(v359) >> (uint(v360) % 32))
													*(*int32)(unsafe.Add(mBase, uint32(v284))) = v361
													v363 = *(*int32)(unsafe.Add(mBase, uint32(v284)+4))
													*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = v363 + int32(-8)
													if v360 < v363 {
														v347 = v361
														continue
													} else {
														break
													}
													break
												}
											}
											*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = int32(0)
											v380 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
											v390 = v380
										} else {
											*(*int32)(unsafe.Add(mBase, uint32(v284)+20)) = int32(1)
											v328 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
											v390 = v328
										}
									}
								}
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v284)+20)) = int32(1)
								v308 = *(*int32)(unsafe.Add(mBase, uint32(v284)+8))
								v390 = v308
							}
							v391 = *(*int32)(unsafe.Add(mBase, uint32(v16)+28))
							if v391 == int32(0) {
								v417 = *(*int32)(unsafe.Add(mBase, uint32(v16)+12))
								v422 = *(*int32)(unsafe.Add(mBase, uint32(v16)+20))
								v423 = *(*int32)(unsafe.Add(mBase, uint32(v16)+16))
								v425 = (v417+int32(7))>>(uint(int32(3))%32) + (v422 - v423)
								if base.Ui32(v25) < base.Ui32(v425) {
									v430 = v16 + int32(8)
									if v430 == int32(0) {
									} else {
										v435 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
										F_WebPSafeFree(m, v435)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
										v441 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
										base.Simd_g_v128_store(m, v430, int32(0), v441)
									}
									v447 = int32(0)
									v448 = int32(1)
									v450 = v25
									v451 = v447
									v452 = v448
									v453 = v24
								} else {
									v450 = v425
									v451 = int32(1)
									v452 = int32(0)
									v453 = v390
								}
								v457 = v451 | l4<<(uint(int32(2))%32)
								if l5 != 0 {
									v460 = v457 | int32(16)
								} else {
									v460 = v457
								}
								*(*uint8)(unsafe.Add(mBase, uint32(v16)+32)) = uint8(v460)
								v463 = l8 + int32(4)
								v465 = v450 + int32(1)
								v466 = int32(0)
								v468 = base.Simd_g_const(&F_EncodeAlphaInternal__k4)
								base.Simd_g_v128_store(m, v463, v466, v468)
								v471 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
								base.Simd_g_v128_store(m, v463, int32(16), v471)
								if v465 == v466 {
									v503 = int32(1)
								} else {
									v477 = int32(1024)
									if base.Ui32(v477) < base.Ui32(v465) {
										v480 = v465
									} else {
										v480 = v477
									}
									v481 = F_WebPSafeMalloc(m, int64(1), v480)
									mBase = m.M
									if v481 != 0 {
										v485 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
										if v485 == int32(0) {
										} else {
											v490 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(20))))
											v491 = F_memcpy(m, v481, v490, v485)
											mBase = m.M
										}
										v492 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
										F_WebPSafeFree(m, v492)
										mBase = m.M
										*(*int32)(unsafe.Add(mBase, uint32(v463)+24)) = v480
										*(*int32)(unsafe.Add(mBase, uint32(v463)+16)) = v481
										v503 = int32(1)
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
										v503 = int32(0)
									}
								}
								if v503 == int32(0) {
									v637 = v26
								} else {
									v508 = int32(1)
									v514 = *(*int32)(unsafe.Add(mBase, uint32(v463)+12))
									if v514 != int32(-8) {
										v562 = int32(0)
										v570 = v562
									} else {
										v517 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
										v520 = base.I64_extend_i32_u(v517) + base.I64_extend_i32_u(v508)
										if base.Ui64(v520) < base.Ui64(int64(4294967296)) {
											v526 = *(*int32)(unsafe.Add(mBase, uint32(v463)+24))
											v527 = base.I32_wrap_i64(v520)
											if base.Ui32(v526) < base.Ui32(v527) {
												v532 = v526 << (uint(int32(1)) % 32)
												if base.Ui32(v527) < base.Ui32(v532) {
													v534 = v532
												} else {
													v534 = v527
												}
												v535 = int32(1024)
												if base.Ui32(v535) < base.Ui32(v534) {
													v538 = v534
												} else {
													v538 = v535
												}
												v539 = F_WebPSafeMalloc(m, int64(1), v538)
												mBase = m.M
												if v539 != 0 {
													v543 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
													if v543 == int32(0) {
													} else {
														v546 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
														v547 = F_memcpy(m, v539, v546, v543)
														mBase = m.M
													}
													v548 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
													F_WebPSafeFree(m, v548)
													mBase = m.M
													*(*int32)(unsafe.Add(mBase, uint32(v463)+24)) = v538
													*(*int32)(unsafe.Add(mBase, uint32(v463)+16)) = v539
													v552 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
													v553 = v552
													v555 = v539
													v557 = F_memcpy(m, v555+v553, v16+int32(32), v508)
													mBase = m.M
													v558 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
													*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v558 + v508
													v562 = int32(1)
													v570 = v562
												} else {
													*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
													v570 = int32(0)
												}
											} else {
												v529 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
												v553 = v517
												v555 = v529
												v557 = F_memcpy(m, v555+v553, v16+int32(32), v508)
												mBase = m.M
												v558 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
												*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v558 + v508
												v562 = int32(1)
												v570 = v562
											}
										} else {
											*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
											v570 = int32(0)
										}
									}
									if v570 == int32(0) {
										v637 = v26
									} else {
										v578 = *(*int32)(unsafe.Add(mBase, uint32(v463)+12))
										if v578 != int32(-8) {
											v626 = int32(0)
											v634 = v626
										} else {
											v581 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
											v584 = base.I64_extend_i32_u(v581) + base.I64_extend_i32_u(v450)
											if base.Ui64(v584) < base.Ui64(int64(4294967296)) {
												v590 = *(*int32)(unsafe.Add(mBase, uint32(v463)+24))
												v591 = base.I32_wrap_i64(v584)
												if base.Ui32(v590) < base.Ui32(v591) {
													v596 = v590 << (uint(int32(1)) % 32)
													if base.Ui32(v591) < base.Ui32(v596) {
														v598 = v596
													} else {
														v598 = v591
													}
													v599 = int32(1024)
													if base.Ui32(v599) < base.Ui32(v598) {
														v602 = v598
													} else {
														v602 = v599
													}
													v603 = F_WebPSafeMalloc(m, int64(1), v602)
													mBase = m.M
													if v603 != 0 {
														v607 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
														if v607 == int32(0) {
														} else {
															v610 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
															v611 = F_memcpy(m, v603, v610, v607)
															mBase = m.M
														}
														v612 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
														F_WebPSafeFree(m, v612)
														mBase = m.M
														*(*int32)(unsafe.Add(mBase, uint32(v463)+24)) = v602
														*(*int32)(unsafe.Add(mBase, uint32(v463)+16)) = v603
														v616 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
														v617 = v616
														v619 = v603
														v621 = F_memcpy(m, v619+v617, v453, v450)
														mBase = m.M
														v622 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
														*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v622 + v450
														v626 = int32(1)
														v634 = v626
													} else {
														*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
														v634 = int32(0)
													}
												} else {
													v593 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
													v617 = v581
													v619 = v593
													v621 = F_memcpy(m, v619+v617, v453, v450)
													mBase = m.M
													v622 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
													*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v622 + v450
													v626 = int32(1)
													v634 = v626
												}
											} else {
												*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
												v634 = int32(0)
											}
										}
										v637 = base.B2i32(v634 != int32(0))
									}
								}
								if v452 != 0 {
								} else {
									v639 = v16 + int32(8)
									if v639 == int32(0) {
									} else {
										v644 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
										F_WebPSafeFree(m, v644)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
										v650 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
										base.Simd_g_v128_store(m, v639, int32(0), v650)
									}
								}
								v653 = int32(0)
								if v637 == v653 {
									v659 = v653
								} else {
									v656 = *(*int32)(unsafe.Add(mBase, uint32(l8)+32))
									v659 = base.B2i32(v656 == int32(0))
								}
								v662 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(24))))
								*(*int32)(unsafe.Add(mBase, uint32(l8))) = v662
								v664 = v659
							} else {
								v395 = v16 + int32(8)
								if v395 == int32(0) {
								} else {
									v400 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
									F_WebPSafeFree(m, v400)
									mBase = m.M
									*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
									v406 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
									base.Simd_g_v128_store(m, v395, int32(0), v406)
								}
								v409 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
								base.Simd_g_v128_store(m, l8, int32(4), v409)
								v414 = int32(0)
								base.Simd_g_v128_store(m, l8+int32(20), v414, v409)
								v664 = v414
							}
						} else {
							v243 = v16 + int32(8)
							if v243 == int32(0) {
							} else {
								v248 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
								F_WebPSafeFree(m, v248)
								mBase = m.M
								*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
								v254 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
								base.Simd_g_v128_store(m, v243, int32(0), v254)
							}
							v261 = v16 + int32(8)
							if v261 == int32(0) {
							} else {
								v266 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
								F_WebPSafeFree(m, v266)
								mBase = m.M
								*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
								v272 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
								base.Simd_g_v128_store(m, v261, int32(0), v272)
							}
							v275 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
							base.Simd_g_v128_store(m, l8, int32(4), v275)
							v280 = int32(0)
							base.Simd_g_v128_store(m, l8+int32(20), v280, v275)
							v664 = v280
						}
					}
				}
			}
		}
	} else {
		v447 = int32(0)
		v448 = int32(1)
		v450 = v25
		v451 = v447
		v452 = v448
		v453 = v24
		v457 = v451 | l4<<(uint(int32(2))%32)
		if l5 != 0 {
			v460 = v457 | int32(16)
		} else {
			v460 = v457
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v16)+32)) = uint8(v460)
		v463 = l8 + int32(4)
		v465 = v450 + int32(1)
		v466 = int32(0)
		v468 = base.Simd_g_const(&F_EncodeAlphaInternal__k4)
		base.Simd_g_v128_store(m, v463, v466, v468)
		v471 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
		base.Simd_g_v128_store(m, v463, int32(16), v471)
		if v465 == v466 {
			v503 = int32(1)
		} else {
			v477 = int32(1024)
			if base.Ui32(v477) < base.Ui32(v465) {
				v480 = v465
			} else {
				v480 = v477
			}
			v481 = F_WebPSafeMalloc(m, int64(1), v480)
			mBase = m.M
			if v481 != 0 {
				v485 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
				if v485 == int32(0) {
				} else {
					v490 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(20))))
					v491 = F_memcpy(m, v481, v490, v485)
					mBase = m.M
				}
				v492 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
				F_WebPSafeFree(m, v492)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(v463)+24)) = v480
				*(*int32)(unsafe.Add(mBase, uint32(v463)+16)) = v481
				v503 = int32(1)
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
				v503 = int32(0)
			}
		}
		if v503 == int32(0) {
			v637 = v26
		} else {
			v508 = int32(1)
			v514 = *(*int32)(unsafe.Add(mBase, uint32(v463)+12))
			if v514 != int32(-8) {
				v562 = int32(0)
				v570 = v562
			} else {
				v517 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
				v520 = base.I64_extend_i32_u(v517) + base.I64_extend_i32_u(v508)
				if base.Ui64(v520) < base.Ui64(int64(4294967296)) {
					v526 = *(*int32)(unsafe.Add(mBase, uint32(v463)+24))
					v527 = base.I32_wrap_i64(v520)
					if base.Ui32(v526) < base.Ui32(v527) {
						v532 = v526 << (uint(int32(1)) % 32)
						if base.Ui32(v527) < base.Ui32(v532) {
							v534 = v532
						} else {
							v534 = v527
						}
						v535 = int32(1024)
						if base.Ui32(v535) < base.Ui32(v534) {
							v538 = v534
						} else {
							v538 = v535
						}
						v539 = F_WebPSafeMalloc(m, int64(1), v538)
						mBase = m.M
						if v539 != 0 {
							v543 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
							if v543 == int32(0) {
							} else {
								v546 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
								v547 = F_memcpy(m, v539, v546, v543)
								mBase = m.M
							}
							v548 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
							F_WebPSafeFree(m, v548)
							mBase = m.M
							*(*int32)(unsafe.Add(mBase, uint32(v463)+24)) = v538
							*(*int32)(unsafe.Add(mBase, uint32(v463)+16)) = v539
							v552 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
							v553 = v552
							v555 = v539
							v557 = F_memcpy(m, v555+v553, v16+int32(32), v508)
							mBase = m.M
							v558 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
							*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v558 + v508
							v562 = int32(1)
							v570 = v562
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
							v570 = int32(0)
						}
					} else {
						v529 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
						v553 = v517
						v555 = v529
						v557 = F_memcpy(m, v555+v553, v16+int32(32), v508)
						mBase = m.M
						v558 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
						*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v558 + v508
						v562 = int32(1)
						v570 = v562
					}
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
					v570 = int32(0)
				}
			}
			if v570 == int32(0) {
				v637 = v26
			} else {
				v578 = *(*int32)(unsafe.Add(mBase, uint32(v463)+12))
				if v578 != int32(-8) {
					v626 = int32(0)
					v634 = v626
				} else {
					v581 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
					v584 = base.I64_extend_i32_u(v581) + base.I64_extend_i32_u(v450)
					if base.Ui64(v584) < base.Ui64(int64(4294967296)) {
						v590 = *(*int32)(unsafe.Add(mBase, uint32(v463)+24))
						v591 = base.I32_wrap_i64(v584)
						if base.Ui32(v590) < base.Ui32(v591) {
							v596 = v590 << (uint(int32(1)) % 32)
							if base.Ui32(v591) < base.Ui32(v596) {
								v598 = v596
							} else {
								v598 = v591
							}
							v599 = int32(1024)
							if base.Ui32(v599) < base.Ui32(v598) {
								v602 = v598
							} else {
								v602 = v599
							}
							v603 = F_WebPSafeMalloc(m, int64(1), v602)
							mBase = m.M
							if v603 != 0 {
								v607 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
								if v607 == int32(0) {
								} else {
									v610 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
									v611 = F_memcpy(m, v603, v610, v607)
									mBase = m.M
								}
								v612 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
								F_WebPSafeFree(m, v612)
								mBase = m.M
								*(*int32)(unsafe.Add(mBase, uint32(v463)+24)) = v602
								*(*int32)(unsafe.Add(mBase, uint32(v463)+16)) = v603
								v616 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
								v617 = v616
								v619 = v603
								v621 = F_memcpy(m, v619+v617, v453, v450)
								mBase = m.M
								v622 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
								*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v622 + v450
								v626 = int32(1)
								v634 = v626
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
								v634 = int32(0)
							}
						} else {
							v593 = *(*int32)(unsafe.Add(mBase, uint32(v463)+16))
							v617 = v581
							v619 = v593
							v621 = F_memcpy(m, v619+v617, v453, v450)
							mBase = m.M
							v622 = *(*int32)(unsafe.Add(mBase, uint32(v463)+20))
							*(*int32)(unsafe.Add(mBase, uint32(v463)+20)) = v622 + v450
							v626 = int32(1)
							v634 = v626
						}
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v463)+28)) = int32(1)
						v634 = int32(0)
					}
				}
				v637 = base.B2i32(v634 != int32(0))
			}
		}
		if v452 != 0 {
		} else {
			v639 = v16 + int32(8)
			if v639 == int32(0) {
			} else {
				v644 = *(*int32)(unsafe.Add(mBase, uint32(v16+int32(16))))
				F_WebPSafeFree(m, v644)
				mBase = m.M
				*(*int64)(unsafe.Add(mBase, uint32(v16+int32(24)))) = int64(0)
				v650 = base.Simd_g_const(&F_EncodeAlphaInternal__k0)
				base.Simd_g_v128_store(m, v639, int32(0), v650)
			}
		}
		v653 = int32(0)
		if v637 == v653 {
			v659 = v653
		} else {
			v656 = *(*int32)(unsafe.Add(mBase, uint32(l8)+32))
			v659 = base.B2i32(v656 == int32(0))
		}
		v662 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(24))))
		*(*int32)(unsafe.Add(mBase, uint32(l8))) = v662
		v664 = v659
	}
	m.G0 = v16 + int32(320)
	return v664
}

var F_EncodeAlphaInternal__k0 = [2]uint64{0x0, 0x0}
var F_EncodeAlphaInternal__k1 = [2]uint64{0x3200000004, 0x3c}
var F_EncodeAlphaInternal__k2 = [2]uint64{0x100000064, 0x0}
var F_EncodeAlphaInternal__k3 = [2]uint64{0x1, 0x100000001}
var F_EncodeAlphaInternal__k4 = [2]uint64{0xfe, 0xfffffff800000000}

func F_EncodeImageNoHuffman(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v14 base.V128
	_ = v14
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v36 int32
	_ = v36
	var v48 int64
	_ = v48
	var v56 int64
	_ = v56
	var v57 int32
	_ = v57
	var v58 int64
	_ = v58
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v101 int32
	_ = v101
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v157 int32
	_ = v157
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v163 int32
	_ = v163
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v195 int32
	_ = v195
	var v201 int32
	_ = v201
	var v203 int32
	_ = v203
	var v207 int32
	_ = v207
	var v212 int32
	_ = v212
	var v214 int32
	_ = v214
	var v216 int32
	_ = v216
	var v221 int32
	_ = v221
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v248 int32
	_ = v248
	var v249 int32
	_ = v249
	var v257 int32
	_ = v257
	var v263 int32
	_ = v263
	var v273 int32
	_ = v273
	var v282 int32
	_ = v282
	var v283 int32
	_ = v283
	var v288 int32
	_ = v288
	var v299 int32
	_ = v299
	var v302 int32
	_ = v302
	var v303 int32
	_ = v303
	var v305 int32
	_ = v305
	var v307 int32
	_ = v307
	var v312 int32
	_ = v312
	var v315 int32
	_ = v315
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v325 int32
	_ = v325
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v342 int32
	_ = v342
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v348 int32
	_ = v348
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v379 int32
	_ = v379
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v396 int32
	_ = v396
	var v397 int32
	_ = v397
	var v400 int64
	_ = v400
	var v403 int32
	_ = v403
	var v405 int64
	_ = v405
	var v408 int32
	_ = v408
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v421 int32
	_ = v421
	var v422 int32
	_ = v422
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v437 int32
	_ = v437
	var v442 int32
	_ = v442
	var v443 int32
	_ = v443
	var v450 int32
	_ = v450
	var v453 int32
	_ = v453
	var v455 int32
	_ = v455
	var v459 int32
	_ = v459
	var v460 int32
	_ = v460
	var v461 int32
	_ = v461
	var v462 int32
	_ = v462
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v502 int32
	_ = v502
	var v503 int32
	_ = v503
	var v505 int32
	_ = v505
	var v506 int32
	_ = v506
	var v508 int32
	_ = v508
	var v509 int32
	_ = v509
	var v511 int32
	_ = v511
	var v512 int32
	_ = v512
	var v515 int32
	_ = v515
	var v516 int64
	_ = v516
	var v517 int32
	_ = v517
	var v524 int64
	_ = v524
	var v525 int32
	_ = v525
	var v526 int64
	_ = v526
	var v536 int32
	_ = v536
	var v538 int32
	_ = v538
	var v543 int32
	_ = v543
	var v549 int32
	_ = v549
	var v552 int32
	_ = v552
	var v560 int32
	_ = v560
	var v563 int32
	_ = v563
	var v564 int32
	_ = v564
	var v571 int32
	_ = v571
	var v578 int32
	_ = v578
	var v582 int32
	_ = v582
	var v583 int32
	_ = v583
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v602 int32
	_ = v602
	var v604 int32
	_ = v604
	var v606 int32
	_ = v606
	var v613 int32
	_ = v613
	var v614 int32
	_ = v614
	var v636 int32
	_ = v636
	var v639 int32
	_ = v639
	var v647 int32
	_ = v647
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v658 int32
	_ = v658
	var v665 int32
	_ = v665
	var v669 int32
	_ = v669
	var v670 int32
	_ = v670
	var v681 int32
	_ = v681
	var v682 int32
	_ = v682
	var v689 int32
	_ = v689
	var v691 int32
	_ = v691
	var v693 int32
	_ = v693
	var v700 int32
	_ = v700
	var v701 int32
	_ = v701
	var v723 int32
	_ = v723
	var v726 int32
	_ = v726
	var v734 int32
	_ = v734
	var v737 int32
	_ = v737
	var v738 int32
	_ = v738
	var v745 int32
	_ = v745
	var v752 int32
	_ = v752
	var v756 int32
	_ = v756
	var v757 int32
	_ = v757
	var v768 int32
	_ = v768
	var v769 int32
	_ = v769
	var v776 int32
	_ = v776
	var v778 int32
	_ = v778
	var v780 int32
	_ = v780
	var v787 int32
	_ = v787
	var v788 int32
	_ = v788
	var v810 int32
	_ = v810
	var v813 int32
	_ = v813
	var v821 int32
	_ = v821
	var v824 int32
	_ = v824
	var v825 int32
	_ = v825
	var v832 int32
	_ = v832
	var v839 int32
	_ = v839
	var v843 int32
	_ = v843
	var v844 int32
	_ = v844
	var v855 int32
	_ = v855
	var v856 int32
	_ = v856
	var v863 int32
	_ = v863
	var v865 int32
	_ = v865
	var v867 int32
	_ = v867
	var v874 int32
	_ = v874
	var v875 int32
	_ = v875
	var v895 int32
	_ = v895
	var v898 int32
	_ = v898
	var v906 int32
	_ = v906
	var v909 int32
	_ = v909
	var v910 int32
	_ = v910
	var v917 int32
	_ = v917
	var v924 int32
	_ = v924
	var v928 int32
	_ = v928
	var v929 int32
	_ = v929
	var v940 int32
	_ = v940
	var v941 int32
	_ = v941
	var v948 int32
	_ = v948
	var v950 int32
	_ = v950
	var v952 int32
	_ = v952
	var v959 int32
	_ = v959
	var v960 int32
	_ = v960
	var v984 int32
	_ = v984
	var v986 int32
	_ = v986
	var v1002 int32
	_ = v1002
	var v1005 int32
	_ = v1005
	var v1010 int32
	_ = v1010
	var v1012 int32
	_ = v1012
	v12 = int32(0)
	v14 = base.Simd_g_const(&F_EncodeImageNoHuffman__k0)
	v18 = m.G0
	v20 = v18 - int32(80)
	m.G0 = v20
	*(*int32)(unsafe.Add(mBase, uint32(v20+int32(72)))) = v12
	*(*int64)(unsafe.Add(mBase, uint32(v20+int32(64)))) = int64(0)
	base.Simd_g_v128_store(m, v20+int32(48), v12, v14)
	v36 = int32(16)
	base.Simd_g_v128_store(m, v20+int32(32), v12, v14)
	base.Simd_g_v128_store(m, v20, v36, v14)
	*(*int32)(unsafe.Add(mBase, uint32(v20)+12)) = v12
	*(*int32)(unsafe.Add(mBase, uint32(v20)+8)) = v12
	v48 = int64(57)
	goto L6
L1:
	;
	F_free(m, v1005)
	mBase = m.M
	goto L167
L2:
	;
	v77 = int32(0)
	v79 = base.I32_div_s(l9, int32(2))
	v80 = F_VP8LHashChainFill(m, l2, l6, l1, l4, l5, l7, l8, v79, l10)
	mBase = m.M
	if v80 == v77 {
		v1002 = v12
		v1005 = v77
		goto L1
	} else {
		goto L13
	}
L3:
	;
	if v70 != 0 {
		goto L2
	} else {
		goto L9
	}
L4:
	;
	goto L3
L5:
	;
	v68 = F_malloc(m, base.I32_wrap_i64(v48)*v36)
	mBase = m.M
	v70 = v68
	goto L4
L6:
	;
	v56 = base.I64_div_u_s(int64(2147418112), v48)
	v57 = int32(0)
	v58 = base.I64_extend_i32_u(v36)
	if base.Ui64(int64(4294967295)) < base.Ui64(v58*v48) {
		v70 = v57
		goto L4
	} else {
		goto L7
	}
L7:
	;
	if base.Ui64(v56) < base.Ui64(v58) {
		v70 = v57
		goto L4
	} else {
		goto L8
	}
L8:
	;
	goto L5
L9:
	;
	v73 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v73 != 0 {
		goto L11
	} else {
		goto L12
	}
L10:
	;
	v1002 = v12
	v1005 = int32(0)
	goto L1
L11:
	;
	goto L10
L12:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l8)+92)) = int32(1)
	goto L11
L13:
	;
	v83 = int32(0)
	v91 = F_VP8LGetBackwardReferences(m, l4, l5, l1, l6, v83, int32(3), v83, v83, l2, l3, v20+int32(8), l8, l9-v79, l10)
	mBase = m.M
	if v91 != 0 {
		goto L14
	} else {
		goto L15
	}
L14:
	;
	v94 = *(*int32)(unsafe.Add(mBase, uint32(v20)+8))
	v95 = F_VP8LAllocateHistogramSet(m, int32(1), v94)
	mBase = m.M
	if v95 != 0 {
		goto L16
	} else {
		goto L17
	}
L15:
	;
	v1002 = int32(0)
	v1005 = v83
	goto L1
L16:
	;
	v101 = int32(0)
	v113 = *(*int32)(unsafe.Add(mBase, uint32(v95)+8))
	v114 = *(*int32)(unsafe.Add(mBase, uint32(v113)))
	v115 = *(*int32)(unsafe.Add(mBase, uint32(v114)+3236))
	v116 = int32(4) << (uint(v115) % 32)
	v117 = int32(_a_F_EncodeImageNoHuffman_0)
	v121 = base.B2i32(v101 < v115)
	if v101 < v115 {
		goto L22
	} else {
		goto L23
	}
L17:
	;
	v97 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v97 != 0 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	v1002 = int32(0)
	v1005 = v83
	goto L1
L19:
	;
	goto L18
L20:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l8)+92)) = int32(1)
	goto L19
L21:
	;
	v299 = int32(0)
	v302 = *(*int32)(unsafe.Add(mBase, uint32(v95)+8))
	v303 = *(*int32)(unsafe.Add(mBase, uint32(v302)))
	v305 = m.G0
	v307 = v305 - int32(16)
	m.G0 = v307
	F_VP8LRefsCursorInit(m, v307+int32(4), l3)
	mBase = m.M
	v312 = *(*int32)(unsafe.Add(mBase, uint32(v307)+4))
	if v312 == v299 {
		goto L48
	} else {
		goto L49
	}
L22:
	;
	v122 = v116 + v117
	goto L24
L23:
	;
	v122 = v117
	goto L24
L24:
	;
	v123 = *(*int32)(unsafe.Add(mBase, uint32(v95)+4))
	v125 = int32(12)
	v127 = F_memset(m, v95, v101, v122*v123+v125)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v127)+4)) = v123
	v130 = v127 + v125
	*(*int32)(unsafe.Add(mBase, uint32(v127)+8)) = v130
	*(*int32)(unsafe.Add(mBase, uint32(v127))) = v123
	if v123 < int32(1) {
		goto L25
	} else {
		goto L26
	}
L25:
	;
	goto L21
L26:
	;
	v135 = int32(1)
	v139 = v130 + v123<<(uint(int32(2))%32)
	if v123 == v135 {
		v195 = v101
		v201 = v130
		v203 = v139
		goto L27
	} else {
		goto L28
	}
L27:
	;
	if v123&v135 == int32(0) {
		v221 = v201
		goto L35
	} else {
		goto L36
	}
L28:
	;
	v142 = int32(_a_F_EncodeImageNoHuffman_1)
	if v101 < v115 {
		goto L29
	} else {
		goto L30
	}
L29:
	;
	v145 = v116 + v142
	goto L31
L30:
	;
	v145 = v142
	goto L31
L31:
	;
	v150 = int32(4)
	v151 = int32(0)
	v157 = v130
	v159 = v139
	goto L32
L32:
	;
	v161 = int32(-4)
	v163 = int32(31)
	v165 = int32(-32)
	v166 = (v159 + v163) & v165
	*(*int32)(unsafe.Add(mBase, uint32(v157+v150+v161))) = v166
	v168 = *(*int32)(unsafe.Add(mBase, uint32(v127)+8))
	v169 = v168 + v150
	v172 = *(*int32)(unsafe.Add(mBase, uint32(v169+v161)))
	v173 = int32(3312)
	*(*int32)(unsafe.Add(mBase, uint32(v172))) = v166 + v173
	v180 = (v166 + v145 + v163) & v165
	*(*int32)(unsafe.Add(mBase, uint32(v169))) = v180
	v182 = *(*int32)(unsafe.Add(mBase, uint32(v127)+8))
	v184 = *(*int32)(unsafe.Add(mBase, uint32(v182+v150)))
	*(*int32)(unsafe.Add(mBase, uint32(v184))) = v180 + v173
	v190 = v180 + v145
	v192 = v151 + int32(2)
	if v123&int32(2147483646) != v192 {
		v150 = v150 + int32(8)
		v151 = v192
		v157 = v182
		v159 = v190
		goto L32
	} else {
		goto L34
	}
L33:
	;
	v195 = v192
	v201 = v182
	v203 = v190
	goto L27
L34:
	;
	goto L33
L35:
	;
	if v123 < int32(1) {
		goto L25
	} else {
		goto L37
	}
L36:
	;
	v207 = v195 << (uint(int32(2)) % 32)
	v212 = (v203 + int32(31)) & int32(-32)
	*(*int32)(unsafe.Add(mBase, uint32(v201+v207))) = v212
	v214 = *(*int32)(unsafe.Add(mBase, uint32(v127)+8))
	v216 = *(*int32)(unsafe.Add(mBase, uint32(v214+v207)))
	*(*int32)(unsafe.Add(mBase, uint32(v216))) = v212 + int32(3312)
	v221 = v214
	goto L35
L37:
	;
	if base.Ui32(v123) < base.Ui32(int32(4)) {
		v263 = int32(0)
		goto L38
	} else {
		goto L39
	}
L38:
	;
	v273 = v221 + v263<<(uint(int32(2))%32)
	v282 = v123 - v263
	goto L44
L39:
	;
	v229 = v123 & int32(2147483644)
	v230 = v221
	v239 = v229
	goto L40
L40:
	;
	v242 = *(*int32)(unsafe.Add(mBase, uint32(v230+int32(12))))
	v245 = *(*int32)(unsafe.Add(mBase, uint32(v230+int32(8))))
	v248 = *(*int32)(unsafe.Add(mBase, uint32(v230+int32(4))))
	v249 = *(*int32)(unsafe.Add(mBase, uint32(v230)))
	*(*int32)(unsafe.Add(mBase, uint32(v249)+3236)) = v115
	*(*int32)(unsafe.Add(mBase, uint32(v248)+3236)) = v115
	*(*int32)(unsafe.Add(mBase, uint32(v245)+3236)) = v115
	*(*int32)(unsafe.Add(mBase, uint32(v242)+3236)) = v115
	v257 = v239 + int32(-4)
	if v257 != 0 {
		v230 = v230 + int32(16)
		v239 = v257
		goto L40
	} else {
		goto L42
	}
L41:
	;
	if v123 == v229 {
		goto L25
	} else {
		goto L43
	}
L42:
	;
	goto L41
L43:
	;
	v263 = v229
	goto L38
L44:
	;
	v283 = *(*int32)(unsafe.Add(mBase, uint32(v273)))
	*(*int32)(unsafe.Add(mBase, uint32(v283)+3236)) = v115
	v288 = v282 + int32(-1)
	if v288 != 0 {
		v273 = v273 + int32(4)
		v282 = v288
		goto L44
	} else {
		goto L46
	}
L45:
	;
	goto L25
L46:
	;
	goto L45
L47:
	;
	v342 = F_GetHuffBitLengthsAndCodes(m, v95, v20+int32(16))
	mBase = m.M
	if v342 != 0 {
		goto L55
	} else {
		goto L56
	}
L48:
	;
	m.G0 = v307 + int32(16)
	goto L47
L49:
	;
	v315 = v312
	goto L50
L50:
	;
	F_HistogramAddSinglePixOrCopy(m, v303, v315, v299, v299)
	mBase = m.M
	v321 = *(*int32)(unsafe.Add(mBase, uint32(v307)+4))
	v323 = v321 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v307)+4)) = v323
	v325 = *(*int32)(unsafe.Add(mBase, uint32(v307)+12))
	if v323 != v325 {
		v331 = v323
		goto L52
	} else {
		goto L53
	}
L51:
	;
	goto L48
L52:
	;
	if v331 != 0 {
		v315 = v331
		goto L50
	} else {
		goto L54
	}
L53:
	;
	F_VP8LRefsCursorNextBlock(m, v307+int32(4))
	mBase = m.M
	v330 = *(*int32)(unsafe.Add(mBase, uint32(v307)+4))
	v331 = v330
	goto L52
L54:
	;
	goto L51
L55:
	;
	v347 = int32(0)
	v348 = int32(1)
	goto L62
L56:
	;
	v344 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v344 != 0 {
		goto L58
	} else {
		goto L59
	}
L57:
	;
	v1002 = v95
	v1005 = v299
	goto L1
L58:
	;
	goto L57
L59:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l8)+92)) = int32(1)
	goto L58
L60:
	;
	v499 = *(*int32)(unsafe.Add(mBase, uint32(v20)+16))
	v500 = *(*int32)(unsafe.Add(mBase, uint32(v20)+28))
	if v500 < v499 {
		goto L87
	} else {
		goto L88
	}
L61:
	;
	goto L60
L62:
	;
	v360 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v361 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v361+v348 < int32(32) {
		goto L65
	} else {
		goto L66
	}
L63:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v462 + v460
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v459<<(uint(v462)%32) | v461
	goto L61
L64:
	;
	v379 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v380 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v384 = v376
	v385 = v377
	v386 = v380
	v387 = v379
	goto L68
L65:
	;
	if v361 < int32(16) {
		v459 = v347
		v460 = v348
		v461 = v360
		v462 = v361
		goto L63
	} else {
		goto L67
	}
L66:
	;
	v365 = int32(32)
	v366 = v365 - v361
	v374 = int32(base.Ui32(v347) >> (uint(v366) % 32))
	v375 = v348 - v366
	v376 = v347<<(uint(v361)%32) | v360
	v377 = v365
	goto L64
L67:
	;
	v374 = v347
	v375 = v348
	v376 = v360
	v377 = v361
	goto L64
L68:
	;
	if base.Ui32(v386+int32(2)) <= base.Ui32(v387) {
		v442 = v386
		v443 = v387
		goto L70
	} else {
		goto L71
	}
L69:
	;
	v459 = v374
	v460 = v375
	v461 = v455
	v462 = v453
	goto L63
L70:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v442))) = uint16(v384)
	v450 = v442 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v450
	v453 = v385 + int32(-16)
	v455 = int32(base.Ui32(v384) >> (uint(int32(16)) % 32))
	if int32(31) < v385 {
		v384 = v455
		v385 = v453
		v386 = v450
		v387 = v443
		goto L68
	} else {
		goto L85
	}
L71:
	;
	v396 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v397 = v387 - v396
	v400 = base.I64_extend_i32_s(v397) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v400) {
		v424 = v396
		goto L73
	} else {
		goto L74
	}
L72:
	;
	if v386 == v396 {
		goto L83
	} else {
		goto L84
	}
L73:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v424
	goto L60
L74:
	;
	v403 = v386 - v396
	v405 = v400 + base.I64_extend_i32_u(v403)
	if base.Ui64(int64(4294967295)) < base.Ui64(v405) {
		v424 = v396
		goto L73
	} else {
		goto L75
	}
L75:
	;
	v408 = base.I32_wrap_i64(v405)
	if v387 == v396 {
		goto L76
	} else {
		goto L77
	}
L76:
	;
	v415 = int32(base.Ui32(v397*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v408) < base.Ui32(v415) {
		goto L79
	} else {
		goto L80
	}
L77:
	;
	if base.Ui32(v408) <= base.Ui32(v397) {
		v442 = v386
		v443 = v387
		goto L70
	} else {
		goto L78
	}
L78:
	;
	goto L76
L79:
	;
	v417 = v415
	goto L81
L80:
	;
	v417 = v408
	goto L81
L81:
	;
	v421 = v417&int32(-1024) + int32(1024)
	v422 = F_WebPSafeMalloc(m, int64(1), v421)
	mBase = m.M
	if v422 != 0 {
		goto L72
	} else {
		goto L82
	}
L82:
	;
	v423 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v424 = v423
	goto L73
L83:
	;
	v435 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v435)
	mBase = m.M
	v437 = v422 + v421
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v437
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v422
	v442 = v422 + v403
	v443 = v437
	goto L70
L84:
	;
	v433 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v434 = F_memcpy(m, v422, v433, v403)
	mBase = m.M
	goto L83
L85:
	;
	goto L69
L86:
	;
	v986 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v986 != 0 {
		goto L165
	} else {
		goto L166
	}
L87:
	;
	v502 = v499
	goto L89
L88:
	;
	v502 = v500
	goto L89
L89:
	;
	v503 = *(*int32)(unsafe.Add(mBase, uint32(v20)+40))
	if v503 < v502 {
		goto L90
	} else {
		goto L91
	}
L90:
	;
	v505 = v502
	goto L92
L91:
	;
	v505 = v503
	goto L92
L92:
	;
	v506 = *(*int32)(unsafe.Add(mBase, uint32(v20)+52))
	if v506 < v505 {
		goto L93
	} else {
		goto L94
	}
L93:
	;
	v508 = v505
	goto L95
L94:
	;
	v508 = v506
	goto L95
L95:
	;
	v509 = *(*int32)(unsafe.Add(mBase, uint32(v20)+64))
	if v509 < v508 {
		goto L96
	} else {
		goto L97
	}
L96:
	;
	v511 = v508
	goto L98
L97:
	;
	v511 = v509
	goto L98
L98:
	;
	v512 = int32(0)
	if v512 < v511 {
		goto L99
	} else {
		goto L100
	}
L99:
	;
	v515 = v511
	goto L101
L100:
	;
	v515 = v512
	goto L101
L101:
	;
	v516 = base.I64_extend_i32_u(v515)
	v517 = int32(2)
	if v516 == int64(0) {
		goto L104
	} else {
		goto L105
	}
L102:
	;
	if v538 == int32(0) {
		goto L86
	} else {
		goto L108
	}
L103:
	;
	goto L102
L104:
	;
	v536 = F_malloc(m, base.I32_wrap_i64(v516)*v517)
	mBase = m.M
	v538 = v536
	goto L103
L105:
	;
	v524 = base.I64_div_u_s(int64(2147418112), v516)
	v525 = int32(0)
	v526 = base.I64_extend_i32_u(v517)
	if base.Ui64(int64(4294967295)) < base.Ui64(v526*v516) {
		v538 = v525
		goto L103
	} else {
		goto L106
	}
L106:
	;
	if base.Ui64(v524) < base.Ui64(v526) {
		v538 = v525
		goto L103
	} else {
		goto L107
	}
L107:
	;
	goto L104
L108:
	;
	v543 = v20 + int32(16)
	F_StoreHuffmanCode(m, l0, v70, v538, v543)
	mBase = m.M
	v549 = *(*int32)(unsafe.Add(mBase, uint32(v20)+16))
	if v549 <= int32(0) {
		goto L109
	} else {
		goto L110
	}
L109:
	;
	F_StoreHuffmanCode(m, l0, v70, v538, v543|int32(12))
	mBase = m.M
	v636 = *(*int32)(unsafe.Add(mBase, uint32(v20)+28))
	if v636 < int32(1) {
		goto L120
	} else {
		goto L121
	}
L110:
	;
	v552 = *(*int32)(unsafe.Add(mBase, uint32(v20)+20))
	v560 = int32(0)
	v563 = v549
	v564 = v552
	goto L111
L111:
	;
	v571 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v564))))
	if v571 == int32(0) {
		v578 = v560
		goto L113
	} else {
		goto L114
	}
L112:
	;
	v583 = int32(0)
	v594 = v583
	v595 = v583
	goto L117
L113:
	;
	v582 = v563 + int32(-1)
	if v582 != 0 {
		v560 = v578
		v563 = v582
		v564 = v564 + int32(1)
		goto L111
	} else {
		goto L116
	}
L114:
	;
	if int32(0) < v560 {
		goto L109
	} else {
		goto L115
	}
L115:
	;
	v578 = int32(1)
	goto L113
L116:
	;
	goto L112
L117:
	;
	v602 = *(*int32)(unsafe.Add(mBase, uint32(v20)+20))
	v604 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v602+v594))) = uint8(v604)
	v606 = *(*int32)(unsafe.Add(mBase, uint32(v20)+24))
	*(*uint16)(unsafe.Add(mBase, uint32(v606+v595))) = uint16(v604)
	v613 = v594 + int32(1)
	v614 = *(*int32)(unsafe.Add(mBase, uint32(v20)+16))
	if v613 < v614 {
		v594 = v613
		v595 = v595 + int32(2)
		goto L117
	} else {
		goto L119
	}
L118:
	;
	goto L109
L119:
	;
	goto L118
L120:
	;
	F_StoreHuffmanCode(m, l0, v70, v538, v20+int32(40))
	mBase = m.M
	v723 = *(*int32)(unsafe.Add(mBase, uint32(v20)+40))
	if v723 < int32(1) {
		goto L131
	} else {
		goto L132
	}
L121:
	;
	v639 = *(*int32)(unsafe.Add(mBase, uint32(v20)+32))
	v647 = int32(0)
	v650 = v636
	v651 = v639
	goto L122
L122:
	;
	v658 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v651))))
	if v658 == int32(0) {
		v665 = v647
		goto L124
	} else {
		goto L125
	}
L123:
	;
	v670 = int32(0)
	v681 = v670
	v682 = v670
	goto L128
L124:
	;
	v669 = v650 + int32(-1)
	if v669 != 0 {
		v647 = v665
		v650 = v669
		v651 = v651 + int32(1)
		goto L122
	} else {
		goto L127
	}
L125:
	;
	if int32(0) < v647 {
		goto L120
	} else {
		goto L126
	}
L126:
	;
	v665 = int32(1)
	goto L124
L127:
	;
	goto L123
L128:
	;
	v689 = *(*int32)(unsafe.Add(mBase, uint32(v20)+32))
	v691 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v689+v681))) = uint8(v691)
	v693 = *(*int32)(unsafe.Add(mBase, uint32(v20)+36))
	*(*uint16)(unsafe.Add(mBase, uint32(v693+v682))) = uint16(v691)
	v700 = v681 + int32(1)
	v701 = *(*int32)(unsafe.Add(mBase, uint32(v20)+28))
	if v700 < v701 {
		v681 = v700
		v682 = v682 + int32(2)
		goto L128
	} else {
		goto L130
	}
L129:
	;
	goto L120
L130:
	;
	goto L129
L131:
	;
	F_StoreHuffmanCode(m, l0, v70, v538, v20+int32(52))
	mBase = m.M
	v810 = *(*int32)(unsafe.Add(mBase, uint32(v20)+52))
	if v810 < int32(1) {
		goto L142
	} else {
		goto L143
	}
L132:
	;
	v726 = *(*int32)(unsafe.Add(mBase, uint32(v20)+44))
	v734 = int32(0)
	v737 = v723
	v738 = v726
	goto L133
L133:
	;
	v745 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v738))))
	if v745 == int32(0) {
		v752 = v734
		goto L135
	} else {
		goto L136
	}
L134:
	;
	v757 = int32(0)
	v768 = v757
	v769 = v757
	goto L139
L135:
	;
	v756 = v737 + int32(-1)
	if v756 != 0 {
		v734 = v752
		v737 = v756
		v738 = v738 + int32(1)
		goto L133
	} else {
		goto L138
	}
L136:
	;
	if int32(0) < v734 {
		goto L131
	} else {
		goto L137
	}
L137:
	;
	v752 = int32(1)
	goto L135
L138:
	;
	goto L134
L139:
	;
	v776 = *(*int32)(unsafe.Add(mBase, uint32(v20)+44))
	v778 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v776+v768))) = uint8(v778)
	v780 = *(*int32)(unsafe.Add(mBase, uint32(v20)+48))
	*(*uint16)(unsafe.Add(mBase, uint32(v780+v769))) = uint16(v778)
	v787 = v768 + int32(1)
	v788 = *(*int32)(unsafe.Add(mBase, uint32(v20)+40))
	if v787 < v788 {
		v768 = v787
		v769 = v769 + int32(2)
		goto L139
	} else {
		goto L141
	}
L140:
	;
	goto L131
L141:
	;
	goto L140
L142:
	;
	F_StoreHuffmanCode(m, l0, v70, v538, v20+int32(64))
	mBase = m.M
	v895 = *(*int32)(unsafe.Add(mBase, uint32(v20)+64))
	if v895 < int32(1) {
		goto L153
	} else {
		goto L154
	}
L143:
	;
	v813 = *(*int32)(unsafe.Add(mBase, uint32(v20)+56))
	v821 = int32(0)
	v824 = v810
	v825 = v813
	goto L144
L144:
	;
	v832 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v825))))
	if v832 == int32(0) {
		v839 = v821
		goto L146
	} else {
		goto L147
	}
L145:
	;
	v844 = int32(0)
	v855 = v844
	v856 = v844
	goto L150
L146:
	;
	v843 = v824 + int32(-1)
	if v843 != 0 {
		v821 = v839
		v824 = v843
		v825 = v825 + int32(1)
		goto L144
	} else {
		goto L149
	}
L147:
	;
	if int32(0) < v821 {
		goto L142
	} else {
		goto L148
	}
L148:
	;
	v839 = int32(1)
	goto L146
L149:
	;
	goto L145
L150:
	;
	v863 = *(*int32)(unsafe.Add(mBase, uint32(v20)+56))
	v865 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v863+v855))) = uint8(v865)
	v867 = *(*int32)(unsafe.Add(mBase, uint32(v20)+60))
	*(*uint16)(unsafe.Add(mBase, uint32(v867+v856))) = uint16(v865)
	v874 = v855 + int32(1)
	v875 = *(*int32)(unsafe.Add(mBase, uint32(v20)+52))
	if v874 < v875 {
		v855 = v874
		v856 = v856 + int32(2)
		goto L150
	} else {
		goto L152
	}
L151:
	;
	goto L142
L152:
	;
	goto L151
L153:
	;
	v984 = F_StoreImageToBitMask(m, l0, l4, int32(0), l3, v20+int32(12), v20+int32(16), l8)
	mBase = m.M
	v1002 = v95
	v1005 = v538
	goto L1
L154:
	;
	v898 = *(*int32)(unsafe.Add(mBase, uint32(v20)+68))
	v906 = int32(0)
	v909 = v895
	v910 = v898
	goto L155
L155:
	;
	v917 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v910))))
	if v917 == int32(0) {
		v924 = v906
		goto L157
	} else {
		goto L158
	}
L156:
	;
	v929 = int32(0)
	v940 = v929
	v941 = v929
	goto L161
L157:
	;
	v928 = v909 + int32(-1)
	if v928 != 0 {
		v906 = v924
		v909 = v928
		v910 = v910 + int32(1)
		goto L155
	} else {
		goto L160
	}
L158:
	;
	if int32(0) < v906 {
		goto L153
	} else {
		goto L159
	}
L159:
	;
	v924 = int32(1)
	goto L157
L160:
	;
	goto L156
L161:
	;
	v948 = *(*int32)(unsafe.Add(mBase, uint32(v20)+68))
	v950 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v948+v940))) = uint8(v950)
	v952 = *(*int32)(unsafe.Add(mBase, uint32(v20)+72))
	*(*uint16)(unsafe.Add(mBase, uint32(v952+v941))) = uint16(v950)
	v959 = v940 + int32(1)
	v960 = *(*int32)(unsafe.Add(mBase, uint32(v20)+64))
	if v959 < v960 {
		v940 = v959
		v941 = v941 + int32(2)
		goto L161
	} else {
		goto L163
	}
L162:
	;
	goto L153
L163:
	;
	goto L162
L164:
	;
	v1002 = v95
	v1005 = int32(0)
	goto L1
L165:
	;
	goto L164
L166:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l8)+92)) = int32(1)
	goto L165
L167:
	;
	F_free(m, v70)
	mBase = m.M
	goto L168
L168:
	;
	F_WebPSafeFree(m, v1002)
	mBase = m.M
	goto L169
L169:
	;
	v1010 = *(*int32)(unsafe.Add(mBase, uint32(v20)+24))
	F_free(m, v1010)
	mBase = m.M
	goto L170
L170:
	;
	v1012 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	m.G0 = v20 + int32(80)
	return base.B2i32(v1012 == int32(0))
}

var F_EncodeImageNoHuffman__k0 = [2]uint64{0x0, 0x0}

func F_ExpandMatrix(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v18 int32
	_ = v18
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v80 int32
	_ = v80
	var v95 base.V128
	_ = v95
	var v106 int32
	_ = v106
	var v111 int64
	_ = v111
	var v115 base.V128
	_ = v115
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var v195 int32
	_ = v195
	var v199 int32
	_ = v199
	var v203 int32
	_ = v203
	var v207 int32
	_ = v207
	var v211 int32
	_ = v211
	var v215 int32
	_ = v215
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v233 int32
	_ = v233
	v3 = int32(0)
	v9 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)) = uint16(v9)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)) = uint16(v9)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)) = uint16(v9)
	v13 = int32(131072)
	v14 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0))))
	v15 = base.I32_div_u_s(v13, v14)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+32)) = uint16(v15)
	v18 = base.I32_div_u_s(v13, v9)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+34)) = uint16(v18)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+36)) = uint16(v18)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+38)) = uint16(v18)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+40)) = uint16(v18)
	v23 = m.G1
	v26 = int32(1)
	v28 = v23 + int32(_a_F_ExpandMatrix_0) + l1<<(uint(v26)%32)
	v29 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28))))
	v30 = int32(9)
	v31 = v29 << (uint(v30) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+64)) = v31
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v28)+1)))
	v35 = v33 << (uint(v30) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+68)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+72)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+76)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+80)) = v35
	v40 = int32(131071)
	v42 = int32(_a_F_ExpandMatrix_1)
	v44 = base.I32_div_u_s(v31^v40, v15&v42)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+128)) = v44
	v50 = base.I32_div_u_s(v35^v40, v18&v42)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+132)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+136)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+140)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+84)) = v35
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+42)) = uint16(v18)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+10)) = uint16(v9)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+88)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = v50
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+44)) = uint16(v18)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v50
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+46)) = uint16(v18)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v50
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+48)) = uint16(v18)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+164)) = v50
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+50)) = uint16(v18)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+52)) = uint16(v18)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+168)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+172)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = v35
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+54)) = uint16(v18)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+56)) = uint16(v18)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+176)) = v50
	v80 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+58)) = uint16(v80)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = v35
	*(*int32)(unsafe.Add(mBase, uint32(l0)+180)) = v50
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+60)) = uint16(v80)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+28)) = uint16(v9)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+184)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = v35
	v95 = base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_replace_lane_l1(base.Simd_g_i16x8_splat(v9), v9), base.Simd_g_const(&F_ExpandMatrix__k0))
	base.Simd_g_v128_store(m, l0, int32(12), v95)
	*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(192)))) = uint16(v3)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+188)) = v50
	*(*int32)(unsafe.Add(mBase, uint32(l0)+124)) = v35
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+62)) = uint16(v80)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+30)) = uint16(v9)
	v106 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0))))
	if l1 == v3 {
		v143 = int32(30)
		v145 = int32(11)
		v146 = int32(base.Ui32(v9*v143) >> (uint(v145) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(194)))) = uint16(v146)
		v150 = int32(90)
		v153 = int32(base.Ui32(v9*v150) >> (uint(v145) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(222)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(220)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(218)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(216)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(214)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(212)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(210)))) = uint16(v153)
		v179 = int32(60)
		v182 = int32(base.Ui32(v9*v179) >> (uint(v145) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(208)))) = uint16(v182)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(206)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(204)))) = uint16(v153)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(202)))) = uint16(v182)
		v195 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)))
		v199 = int32(base.Ui32(v195*v179) >> (uint(v145) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(196)))) = uint16(v199)
		v203 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)))
		v207 = int32(base.Ui32(v203*v150) >> (uint(v145) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(198)))) = uint16(v207)
		v211 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)))
		v215 = int32(base.Ui32(v211*v143) >> (uint(v145) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(200)))) = uint16(v215)
		v217 = int32(1)
		v218 = v9 << (uint(v217) % 32)
		v233 = v211 + (v203 + (v195 + v218)) + v218 + v9<<(uint(int32(2))%32) + ((v9+v9)<<(uint(v217)%32) + v106)
	} else {
		v111 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(194)))) = v111
		v115 = base.Simd_g_const(&F_ExpandMatrix__k1)
		base.Simd_g_v128_store(m, l0+int32(202), int32(0), v115)
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(216)))) = v111
		v122 = int32(1)
		v123 = v9 << (uint(v122) % 32)
		v124 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)))
		v126 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)))
		v128 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)))
		v233 = v123 + v124 + v126 + v128 + v123 + v9<<(uint(int32(2))%32) + ((v9+v9)<<(uint(v122)%32) + v106)
	}
	return int32(base.Ui32(v233+int32(8)) >> (uint(int32(4)) % 32))
}

var F_ExpandMatrix__k0 = [2]uint64{0x100010001000100, 0x302010001000100}
var F_ExpandMatrix__k1 = [2]uint64{0x0, 0x0}

func F_ExtraCost_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v30 int32
	_ = v30
	var v36 int32
	_ = v36
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v47 int32
	_ = v47
	var v48 base.V128
	_ = v48
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v68 base.V128
	_ = v68
	var v71 base.V128
	_ = v71
	var v74 base.V128
	_ = v74
	var v77 base.V128
	_ = v77
	var v78 int32
	_ = v78
	var v85 base.V128
	_ = v85
	var v88 base.V128
	_ = v88
	var v91 base.V128
	_ = v91
	var v94 base.V128
	_ = v94
	var v97 base.V128
	_ = v97
	var v99 base.V128
	_ = v99
	var v101 int32
	_ = v101
	var v104 base.V128
	_ = v104
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
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
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v16 = v14 + v15
	if l1 < int32(8) {
		v158 = v16
	} else {
		v21 = int32(base.Ui32(l1) >> (uint(int32(1)) % 32))
		v23 = v21 + int32(-3)
		if base.Ui32(v23) < base.Ui32(int32(4)) {
			v113 = v16
			v114 = int32(2)
			v131 = v114<<(uint(int32(3))%32) + l0 + int32(8)
			v133 = v113
			v134 = v114
			for {
				v146 = *(*int32)(unsafe.Add(mBase, uint32(v131+int32(4))))
				v147 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
				v150 = (v146+v147)*v134 + v133
				v154 = v134 + int32(1)
				if v21+int32(-1) != v154 {
					v131 = v131 + int32(8)
					v133 = v150
					v134 = v154
					continue
				} else {
					break
				}
				break
			}
			v158 = v150
		} else {
			v30 = v23 & int32(-4)
			v36 = v30
			v40 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_ExtraCost_C__k0), v16)
			v42 = base.Simd_g_const(&F_ExtraCost_C__k1)
			for {
				v47 = int32(3)
				v48 = base.Simd_g_i32x4_shl(v42, v47)
				v51 = l0 + base.Simd_g_i32x4_extract_lane_l3(v48)
				v52 = int32(12)
				v54 = int32(2)
				v56 = l0 + base.Simd_g_i32x4_extract_lane_l2(v48)
				v59 = int32(1)
				v61 = l0 + base.Simd_g_i32x4_extract_lane_l1(v48)
				v64 = int32(0)
				v66 = l0 + base.Simd_g_i32x4_extract_lane_l0(v48)
				v68 = base.Simd_g_v128_load32_splat(m, v66, v52)
				v71 = base.Simd_g_v128_load32_lane_l1(m, v61+v52, v64, v68)
				v74 = base.Simd_g_v128_load32_lane_l2(m, v56+v52, v64, v71)
				v77 = base.Simd_g_v128_load32_lane_l3(m, v51+v52, v64, v74)
				v78 = int32(8)
				v85 = base.Simd_g_v128_load32_splat(m, v66, v78)
				v88 = base.Simd_g_v128_load32_lane_l1(m, v61+v78, v64, v85)
				v91 = base.Simd_g_v128_load32_lane_l2(m, v56+v78, v64, v88)
				v94 = base.Simd_g_v128_load32_lane_l3(m, v51+v78, v64, v91)
				v97 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(v77, v94), v42), v40)
				v99 = base.Simd_g_i32x4_add(v42, base.Simd_g_const(&F_ExtraCost_C__k2))
				v101 = v36 + int32(-4)
				if v101 != 0 {
					v36 = v101
					v40 = v97
					v42 = v99
					continue
				} else {
					break
				}
				break
			}
			v104 = base.Simd_g_i32x4_add(v97, base.Simd_g_i8x16_shuffle2(v97, v99, base.Simd_g_const(&F_ExtraCost_C__k3), base.Simd_g_const(&F_ExtraCost_C__k4)))
			v109 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v104, base.Simd_g_i8x16_shuffle2(v104, v104, base.Simd_g_const(&F_ExtraCost_C__k5), base.Simd_g_const(&F_ExtraCost_C__k4))))
			if v23 == v30 {
				v158 = v109
			} else {
				v113 = v109
				v114 = v30 | int32(2)
				v131 = v114<<(uint(int32(3))%32) + l0 + int32(8)
				v133 = v113
				v134 = v114
				for {
					v146 = *(*int32)(unsafe.Add(mBase, uint32(v131+int32(4))))
					v147 = *(*int32)(unsafe.Add(mBase, uint32(v131)))
					v150 = (v146+v147)*v134 + v133
					v154 = v134 + int32(1)
					if v21+int32(-1) != v154 {
						v131 = v131 + int32(8)
						v133 = v150
						v134 = v154
						continue
					} else {
						break
					}
					break
				}
				v158 = v150
			}
		}
	}
	return v158
}

var F_ExtraCost_C__k0 = [2]uint64{0x0, 0x0}
var F_ExtraCost_C__k1 = [2]uint64{0x300000002, 0x500000004}
var F_ExtraCost_C__k2 = [2]uint64{0x400000004, 0x400000004}
var F_ExtraCost_C__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_ExtraCost_C__k4 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ExtraCost_C__k5 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_ExtraCost_SSE41(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v9 base.V128
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v19 base.V128
	_ = v19
	var v25 int32
	_ = v25
	var v27 base.V128
	_ = v27
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v33 base.V128
	_ = v33
	var v37 base.V128
	_ = v37
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v61 base.V128
	_ = v61
	var v65 int32
	_ = v65
	var v69 base.V128
	_ = v69
	var v74 base.V128
	_ = v74
	var v77 base.V128
	_ = v77
	v8 = int32(16)
	v9 = base.Simd_g_v128_load64_zero(m, l0, v8)
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v11 = int32(1)
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v19 = base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(v9, v10<<(uint(v11)%32)), v15<<(uint(v11)%32))
	if l1 < v8 {
		v69 = v19
	} else {
		v25 = l0 + int32(32)
		v27 = v19
		v28 = int32(16)
		for {
			v32 = int32(0)
			v33 = base.Simd_g_v128_load_rng(m, v25, v32, int32(0), int32(32))
			v37 = base.Simd_g_v128_load_nc(m, v25+int32(16), v32)
			v45 = int32(1)
			v46 = (v28 + int32(-10)) >> (uint(v45) % 32)
			v52 = int32(2)
			v56 = int32(3)
			v61 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i8x16_shuffle2(v33, v37, base.Simd_g_const(&F_ExtraCost_SSE41__k0), base.Simd_g_const(&F_ExtraCost_SSE41__k1)), base.Simd_g_i8x16_shuffle2(v33, v37, base.Simd_g_const(&F_ExtraCost_SSE41__k2), base.Simd_g_const(&F_ExtraCost_SSE41__k3))), base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v46), v46+v45), v46+v52), v46+v56)), v27)
			v65 = v28 + int32(8)
			if v65 <= l1 {
				v25 = v25 + int32(32)
				v27 = v61
				v28 = v65
				continue
			} else {
				break
			}
			break
		}
		v69 = v61
	}
	v74 = base.Simd_g_const(&F_ExtraCost_SSE41__k4)
	v77 = base.Simd_g_i32x4_add(v69, base.Simd_g_i8x16_shuffle2(v69, v74, base.Simd_g_const(&F_ExtraCost_SSE41__k5), base.Simd_g_const(&F_ExtraCost_SSE41__k6)))
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v77, base.Simd_g_i8x16_shuffle2(v77, v74, base.Simd_g_const(&F_ExtraCost_SSE41__k7), base.Simd_g_const(&F_ExtraCost_SSE41__k8))))
}

var F_ExtraCost_SSE41__k0 = [2]uint64{0xb0a090803020100, 0x8080808080808080}
var F_ExtraCost_SSE41__k1 = [2]uint64{0x8080808080808080, 0xb0a090803020100}
var F_ExtraCost_SSE41__k2 = [2]uint64{0xf0e0d0c07060504, 0x8080808080808080}
var F_ExtraCost_SSE41__k3 = [2]uint64{0x8080808080808080, 0xf0e0d0c07060504}
var F_ExtraCost_SSE41__k4 = [2]uint64{0x0, 0x0}
var F_ExtraCost_SSE41__k5 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_ExtraCost_SSE41__k6 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_ExtraCost_SSE41__k7 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_ExtraCost_SSE41__k8 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_ExtractAlpha_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v26 int32
	_ = v26
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v54 int32
	_ = v54
	var v67 int32
	_ = v67
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
	var v81 int32
	_ = v81
	var v82 base.V128
	_ = v82
	var v83 int32
	_ = v83
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v96 base.V128
	_ = v96
	var v110 base.V128
	_ = v110
	var v124 base.V128
	_ = v124
	var v138 base.V128
	_ = v138
	var v141 base.V128
	_ = v141
	var v144 base.V128
	_ = v144
	var v147 base.V128
	_ = v147
	var v150 base.V128
	_ = v150
	var v153 base.V128
	_ = v153
	var v156 base.V128
	_ = v156
	var v159 base.V128
	_ = v159
	var v162 base.V128
	_ = v162
	var v165 base.V128
	_ = v165
	var v168 base.V128
	_ = v168
	var v171 base.V128
	_ = v171
	var v174 base.V128
	_ = v174
	var v177 base.V128
	_ = v177
	var v180 base.V128
	_ = v180
	var v183 base.V128
	_ = v183
	var v186 base.V128
	_ = v186
	var v187 base.V128
	_ = v187
	var v193 int32
	_ = v193
	var v197 base.V128
	_ = v197
	var v200 base.V128
	_ = v200
	var v203 base.V128
	_ = v203
	var v208 int32
	_ = v208
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v238 int32
	_ = v238
	var v240 int32
	_ = v240
	var v249 int32
	_ = v249
	var v252 int32
	_ = v252
	var v256 int32
	_ = v256
	var v258 int32
	_ = v258
	var v268 int32
	_ = v268
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v294 int32
	_ = v294
	v20 = int32(1)
	if l3 < v20 {
		v294 = v20
	} else {
		if l2 < int32(1) {
			v294 = v20
		} else {
			v26 = l2 & int32(2147483632)
			v31 = l0
			v35 = l4
			v39 = int32(255)
			v40 = int32(0)
			for {
				if base.B2i32(base.Ui32(l2) < base.Ui32(int32(16))) == int32(0) {
					v54 = int32(0)
					v67 = v54
					v72 = base.Simd_g_i8x16_replace_lane_l0(base.Simd_g_const(&F_ExtractAlpha_C__k0), v39)
					v73 = base.Simd_g_const(&F_ExtractAlpha_C__k1)
					v74 = base.Simd_g_const(&F_ExtractAlpha_C__k2)
					v75 = base.Simd_g_const(&F_ExtractAlpha_C__k3)
					v76 = base.Simd_g_const(&F_ExtractAlpha_C__k4)
					for {
						v81 = int32(2)
						v82 = base.Simd_g_i32x4_shl(v73, v81)
						v83 = int32(3)
						v89 = int32(1)
						v92 = int32(0)
						v96 = base.Simd_g_i32x4_shl(v74, v81)
						v110 = base.Simd_g_i32x4_shl(v75, v81)
						v124 = base.Simd_g_i32x4_shl(v76, v81)
						v138 = base.Simd_g_v128_load8_splat(m, v31+base.Simd_g_i32x4_extract_lane_l0(v124), v92)
						v141 = base.Simd_g_v128_load8_lane_l1(m, v31+base.Simd_g_i32x4_extract_lane_l1(v124), v92, v138)
						v144 = base.Simd_g_v128_load8_lane_l2(m, v31+base.Simd_g_i32x4_extract_lane_l2(v124), v92, v141)
						v147 = base.Simd_g_v128_load8_lane_l3(m, v31+base.Simd_g_i32x4_extract_lane_l3(v124), v92, v144)
						v150 = base.Simd_g_v128_load8_lane_l4(m, v31+base.Simd_g_i32x4_extract_lane_l0(v110), v92, v147)
						v153 = base.Simd_g_v128_load8_lane_l5(m, v31+base.Simd_g_i32x4_extract_lane_l1(v110), v92, v150)
						v156 = base.Simd_g_v128_load8_lane_l6(m, v31+base.Simd_g_i32x4_extract_lane_l2(v110), v92, v153)
						v159 = base.Simd_g_v128_load8_lane_l7(m, v31+base.Simd_g_i32x4_extract_lane_l3(v110), v92, v156)
						v162 = base.Simd_g_v128_load8_lane_l8(m, v31+base.Simd_g_i32x4_extract_lane_l0(v96), v92, v159)
						v165 = base.Simd_g_v128_load8_lane_l9(m, v31+base.Simd_g_i32x4_extract_lane_l1(v96), v92, v162)
						v168 = base.Simd_g_v128_load8_lane_l10(m, v31+base.Simd_g_i32x4_extract_lane_l2(v96), v92, v165)
						v171 = base.Simd_g_v128_load8_lane_l11(m, v31+base.Simd_g_i32x4_extract_lane_l3(v96), v92, v168)
						v174 = base.Simd_g_v128_load8_lane_l12(m, v31+base.Simd_g_i32x4_extract_lane_l0(v82), v92, v171)
						v177 = base.Simd_g_v128_load8_lane_l13(m, v31+base.Simd_g_i32x4_extract_lane_l1(v82), v92, v174)
						v180 = base.Simd_g_v128_load8_lane_l14(m, v31+base.Simd_g_i32x4_extract_lane_l2(v82), v92, v177)
						v183 = base.Simd_g_v128_load8_lane_l15(m, v31+base.Simd_g_i32x4_extract_lane_l3(v82), v92, v180)
						base.Simd_g_v128_store(m, v35+v67, v92, v183)
						v186 = base.Simd_g_v128_and(v183, v72)
						v187 = base.Simd_g_const(&F_ExtractAlpha_C__k5)
						v193 = v67 + int32(16)
						if v26 != v193 {
							v67 = v193
							v72 = v186
							v73 = base.Simd_g_i32x4_add(v73, v187)
							v74 = base.Simd_g_i32x4_add(v74, v187)
							v75 = base.Simd_g_i32x4_add(v75, v187)
							v76 = base.Simd_g_i32x4_add(v76, v187)
							continue
						} else {
							break
						}
						break
					}
					v197 = base.Simd_g_v128_and(v186, base.Simd_g_i8x16_swizzle_c(v186, base.Simd_g_const(&F_ExtractAlpha_C__k6)))
					v200 = base.Simd_g_v128_and(v197, base.Simd_g_i8x16_shuffle2(v197, v197, base.Simd_g_const(&F_ExtractAlpha_C__k7), base.Simd_g_const(&F_ExtractAlpha_C__k8)))
					v203 = base.Simd_g_v128_and(v200, base.Simd_g_i8x16_shuffle2(v200, v200, base.Simd_g_const(&F_ExtractAlpha_C__k9), base.Simd_g_const(&F_ExtractAlpha_C__k8)))
					v208 = base.Simd_g_i8x16_extract_lane_u_l0(base.Simd_g_v128_and(v203, base.Simd_g_i8x16_shuffle2(v203, v203, base.Simd_g_const(&F_ExtractAlpha_C__k10), base.Simd_g_const(&F_ExtractAlpha_C__k8))))
					if v26 == l2 {
						v268 = v208
					} else {
						v216 = v26
						v218 = v208
						v238 = v216
						v240 = v218
						v249 = v31 + v216<<(uint(int32(2))%32)
						for {
							v252 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v249))))
							*(*uint8)(unsafe.Add(mBase, uint32(v35+v238))) = uint8(v252)
							v256 = v252 & v240
							v258 = v238 + int32(1)
							if l2 != v258 {
								v238 = v258
								v240 = v256
								v249 = v249 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						v268 = v256
					}
				} else {
					v216 = int32(0)
					v218 = v39
					v238 = v216
					v240 = v218
					v249 = v31 + v216<<(uint(int32(2))%32)
					for {
						v252 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v249))))
						*(*uint8)(unsafe.Add(mBase, uint32(v35+v238))) = uint8(v252)
						v256 = v252 & v240
						v258 = v238 + int32(1)
						if l2 != v258 {
							v238 = v258
							v240 = v256
							v249 = v249 + int32(4)
							continue
						} else {
							break
						}
						break
					}
					v268 = v256
				}
				v282 = v40 + int32(1)
				if v282 != l3 {
					v31 = v31 + l1
					v35 = v35 + l5
					v39 = v268
					v40 = v282
					continue
				} else {
					break
				}
				break
			}
			v284 = int32(255)
			v294 = base.B2i32(v268&v284 == v284)
		}
	}
	return v294
}

var F_ExtractAlpha_C__k0 = [2]uint64{0xffffffffffffff00, 0xffffffffffffffff}
var F_ExtractAlpha_C__k1 = [2]uint64{0xd0000000c, 0xf0000000e}
var F_ExtractAlpha_C__k2 = [2]uint64{0x900000008, 0xb0000000a}
var F_ExtractAlpha_C__k3 = [2]uint64{0x500000004, 0x700000006}
var F_ExtractAlpha_C__k4 = [2]uint64{0x100000000, 0x300000002}
var F_ExtractAlpha_C__k5 = [2]uint64{0x1000000010, 0x1000000010}
var F_ExtractAlpha_C__k6 = [2]uint64{0xf0e0d0c0b0a0908, 0x0}
var F_ExtractAlpha_C__k7 = [2]uint64{0x7060504, 0x0}
var F_ExtractAlpha_C__k8 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ExtractAlpha_C__k9 = [2]uint64{0x302, 0x0}
var F_ExtractAlpha_C__k10 = [2]uint64{0x1, 0x0}

func F_ExtractAlpha_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v22 base.V128
	_ = v22
	var v29 int32
	_ = v29
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v74 base.V128
	_ = v74
	var v75 int32
	_ = v75
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v99 base.V128
	_ = v99
	var v102 int32
	_ = v102
	var v106 int32
	_ = v106
	var v107 base.V128
	_ = v107
	var v108 base.V128
	_ = v108
	var v110 int32
	_ = v110
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v122 base.V128
	_ = v122
	var v124 base.V128
	_ = v124
	var v128 base.V128
	_ = v128
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v134 int32
	_ = v134
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v150 base.V128
	_ = v150
	var v157 int32
	_ = v157
	var v158 base.V128
	_ = v158
	var v159 base.V128
	_ = v159
	var v162 base.V128
	_ = v162
	var v164 base.V128
	_ = v164
	var v165 base.V128
	_ = v165
	var v172 int32
	_ = v172
	var v173 base.V128
	_ = v173
	var v192 base.V128
	_ = v192
	var v195 int32
	_ = v195
	var v204 int32
	_ = v204
	var v206 base.V128
	_ = v206
	var v207 int32
	_ = v207
	var v208 base.V128
	_ = v208
	var v209 int32
	_ = v209
	var v215 int32
	_ = v215
	var v218 int32
	_ = v218
	var v222 base.V128
	_ = v222
	var v225 base.V128
	_ = v225
	var v228 base.V128
	_ = v228
	var v231 base.V128
	_ = v231
	var v237 base.V128
	_ = v237
	var v243 int32
	_ = v243
	var v246 base.V128
	_ = v246
	var v250 int32
	_ = v250
	var v251 int32
	_ = v251
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v300 int32
	_ = v300
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v315 int32
	_ = v315
	var v332 int32
	_ = v332
	var v340 int32
	_ = v340
	var v345 int32
	_ = v345
	var v349 int32
	_ = v349
	var v352 int32
	_ = v352
	var v358 int32
	_ = v358
	var v370 int32
	_ = v370
	var v382 int32
	_ = v382
	var v389 base.V128
	_ = v389
	var v394 base.V128
	_ = v394
	var v396 int32
	_ = v396
	var v397 base.V128
	_ = v397
	var v398 int32
	_ = v398
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v411 base.V128
	_ = v411
	var v414 base.V128
	_ = v414
	var v417 base.V128
	_ = v417
	var v420 base.V128
	_ = v420
	var v426 base.V128
	_ = v426
	var v430 int32
	_ = v430
	var v434 base.V128
	_ = v434
	var v439 int32
	_ = v439
	var v448 int32
	_ = v448
	var v449 int32
	_ = v449
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v487 int32
	_ = v487
	var v491 int32
	_ = v491
	var v493 int32
	_ = v493
	var v502 int32
	_ = v502
	var v519 int32
	_ = v519
	var v527 base.V128
	_ = v527
	var v528 int32
	_ = v528
	v22 = base.Simd_g_const(&F_ExtractAlpha_SSE2__k0)
	if int32(1) <= l3 {
		v29 = (l2 + int32(-1)) & int32(-8)
		if v29 < int32(1) {
			if int32(1) <= l2 {
				v340 = l2 & int32(2147483644)
				v345 = l0
				v349 = l4
				v352 = int32(255)
				v358 = int32(0)
				for {
					if base.B2i32(base.Ui32(l2) < base.Ui32(int32(4))) == int32(0) {
						v370 = int32(0)
						v382 = v370
						v389 = base.Simd_g_const(&F_ExtractAlpha_SSE2__k1)
						v394 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_ExtractAlpha_SSE2__k2), v352)
						for {
							v396 = int32(2)
							v397 = base.Simd_g_i32x4_shl(v389, v396)
							v398 = int32(3)
							v404 = int32(1)
							v407 = int32(0)
							v411 = base.Simd_g_v128_load8_splat(m, v345+base.Simd_g_i32x4_extract_lane_l0(v397), v407)
							v414 = base.Simd_g_v128_load8_lane_l1(m, v345+base.Simd_g_i32x4_extract_lane_l1(v397), v407, v411)
							v417 = base.Simd_g_v128_load8_lane_l2(m, v345+base.Simd_g_i32x4_extract_lane_l2(v397), v407, v414)
							v420 = base.Simd_g_v128_load8_lane_l3(m, v345+base.Simd_g_i32x4_extract_lane_l3(v397), v407, v417)
							base.Simd_g_v128_store32_lane_l0(m, v349+v382, v407, v420)
							v426 = base.Simd_g_v128_and(v394, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v420)))
							v430 = v382 + int32(4)
							if v340 != v430 {
								v382 = v430
								v389 = base.Simd_g_i32x4_add(v389, base.Simd_g_const(&F_ExtractAlpha_SSE2__k3))
								v394 = v426
								continue
							} else {
								break
							}
							break
						}
						v434 = base.Simd_g_v128_and(v426, base.Simd_g_i8x16_shuffle2(v426, v420, base.Simd_g_const(&F_ExtractAlpha_SSE2__k4), base.Simd_g_const(&F_ExtractAlpha_SSE2__k5)))
						v439 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_v128_and(v434, base.Simd_g_i8x16_shuffle2(v434, v434, base.Simd_g_const(&F_ExtractAlpha_SSE2__k6), base.Simd_g_const(&F_ExtractAlpha_SSE2__k5))))
						if v340 == l2 {
							v502 = v439
						} else {
							v448 = v439
							v449 = v340
							v472 = v448
							v473 = v449
							v474 = v345 + v449<<(uint(int32(2))%32)
							for {
								v487 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v474))))
								*(*uint8)(unsafe.Add(mBase, uint32(v349+v473))) = uint8(v487)
								v491 = v472 & v487
								v493 = v473 + int32(1)
								if l2 != v493 {
									v472 = v491
									v473 = v493
									v474 = v474 + int32(4)
									continue
								} else {
									break
								}
								break
							}
							v502 = v491
						}
					} else {
						v448 = v352
						v449 = int32(0)
						v472 = v448
						v473 = v449
						v474 = v345 + v449<<(uint(int32(2))%32)
						for {
							v487 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v474))))
							*(*uint8)(unsafe.Add(mBase, uint32(v349+v473))) = uint8(v487)
							v491 = v472 & v487
							v493 = v473 + int32(1)
							if l2 != v493 {
								v472 = v491
								v473 = v493
								v474 = v474 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						v502 = v491
					}
					v519 = v358 + int32(1)
					if v519 != l3 {
						v345 = v345 + l1
						v349 = v349 + l5
						v352 = v502
						v358 = v519
						continue
					} else {
						break
					}
					break
				}
				v527 = v22
				v528 = v502
			} else {
				v527 = v22
				v528 = int32(255)
			}
		} else {
			v33 = v29 + int32(-1)
			v35 = v33 & int32(-8)
			v37 = v35 + int32(9)
			if v37 < l2 {
				v39 = l2
			} else {
				v39 = v37
			}
			v42 = v39 - v35 + int32(-8)
			v43 = int32(3)
			v44 = v39 & v43
			v45 = v42 - v44
			v59 = l0
			v63 = l4
			v66 = int32(255)
			v74 = base.Simd_g_const(&F_ExtractAlpha_SSE2__k7)
			v75 = int32(0)
			for {
				if base.B2i32(v33 == int32(7)) == int32(0) {
					v92 = v59
					v93 = int32(0)
					v99 = v74
					v102 = (int32(base.Ui32(v33)>>(uint(v43)%32)) + int32(1)) & int32(1073741822)
					for {
						v106 = int32(0)
						v107 = base.Simd_g_v128_load_rng(m, v92, v106, int32(0), int32(64))
						v108 = base.Simd_g_const(&F_ExtractAlpha_SSE2__k8)
						v110 = int32(16)
						v111 = base.Simd_g_v128_load_nc(m, v92, v110)
						v113 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(v107, v108), base.Simd_g_v128_and(v111, v108))
						v114 = base.Simd_g_i8x16_narrow_i16x8_u(v113, v113)
						v116 = base.Simd_g_v128_load_nc(m, v92, int32(32))
						v119 = base.Simd_g_v128_load_nc(m, v92, int32(48))
						v121 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(v116, v108), base.Simd_g_v128_and(v119, v108))
						v122 = base.Simd_g_i8x16_narrow_i16x8_u(v121, v121)
						v124 = base.Simd_g_i8x16_shuffle2(v114, v122, base.Simd_g_const(&F_ExtractAlpha_SSE2__k9), base.Simd_g_const(&F_ExtractAlpha_SSE2__k10))
						base.Simd_g_v128_store(m, v63+v93, v106, v124)
						v128 = base.Simd_g_v128_and(base.Simd_g_v128_and(v99, v114), v122)
						v130 = v93 + v110
						v132 = v92 + int32(64)
						v134 = v102 + int32(-2)
						if v134 != 0 {
							v92 = v132
							v93 = v130
							v99 = v128
							v102 = v134
							continue
						} else {
							break
						}
						break
					}
					v143 = v132
					v144 = v130
					v150 = v128
				} else {
					v143 = v59
					v144 = int32(0)
					v150 = v74
				}
				if v33&int32(8) != 0 {
					v172 = v144
					v173 = v150
				} else {
					v157 = int32(0)
					v158 = base.Simd_g_v128_load_rng(m, v143, v157, int32(0), int32(32))
					v159 = base.Simd_g_const(&F_ExtractAlpha_SSE2__k8)
					v162 = base.Simd_g_v128_load_nc(m, v143, int32(16))
					v164 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(v158, v159), base.Simd_g_v128_and(v162, v159))
					v165 = base.Simd_g_i8x16_narrow_i16x8_u(v164, v164)
					base.Simd_g_v128_store64_lane_l0(m, v63+v144, v157, v165)
					v172 = v144 + int32(8)
					v173 = base.Simd_g_v128_and(v150, v165)
				}
				if l2 <= v172 {
					v315 = v66
				} else {
					if base.Ui32(int32(3)) < base.Ui32(v42) {
						v192 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_ExtractAlpha_SSE2__k2), v66)
						v195 = v45
						v204 = v63 + v172
						v206 = base.Simd_g_v128_or(base.Simd_g_i32x4_splat(v172), base.Simd_g_const(&F_ExtractAlpha_SSE2__k1))
						for {
							v207 = int32(2)
							v208 = base.Simd_g_i32x4_shl(v206, v207)
							v209 = int32(3)
							v215 = int32(1)
							v218 = int32(0)
							v222 = base.Simd_g_v128_load8_splat(m, v59+base.Simd_g_i32x4_extract_lane_l0(v208), v218)
							v225 = base.Simd_g_v128_load8_lane_l1(m, v59+base.Simd_g_i32x4_extract_lane_l1(v208), v218, v222)
							v228 = base.Simd_g_v128_load8_lane_l2(m, v59+base.Simd_g_i32x4_extract_lane_l2(v208), v218, v225)
							v231 = base.Simd_g_v128_load8_lane_l3(m, v59+base.Simd_g_i32x4_extract_lane_l3(v208), v218, v228)
							base.Simd_g_v128_store32_lane_l0(m, v204, v218, v231)
							v237 = base.Simd_g_v128_and(v192, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v231)))
							v243 = v195 + int32(-4)
							if v243 != 0 {
								v192 = v237
								v195 = v243
								v204 = v204 + int32(4)
								v206 = base.Simd_g_i32x4_add(v206, base.Simd_g_const(&F_ExtractAlpha_SSE2__k3))
								continue
							} else {
								break
							}
							break
						}
						v246 = base.Simd_g_v128_and(v237, base.Simd_g_i8x16_shuffle2(v237, v231, base.Simd_g_const(&F_ExtractAlpha_SSE2__k4), base.Simd_g_const(&F_ExtractAlpha_SSE2__k5)))
						v250 = int32(0)
						v251 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_v128_and(v246, base.Simd_g_i8x16_shuffle2(v246, v246, base.Simd_g_const(&F_ExtractAlpha_SSE2__k6), base.Simd_g_const(&F_ExtractAlpha_SSE2__k5))))
						if v44 == v250 {
							v315 = v251
						} else {
							v261 = v251
							v262 = v172 + v45
							v285 = v261
							v286 = v262
							v287 = v59 + v262<<(uint(int32(2))%32)
							for {
								v300 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v287))))
								*(*uint8)(unsafe.Add(mBase, uint32(v63+v286))) = uint8(v300)
								v304 = v285 & v300
								v306 = v286 + int32(1)
								if v306 < l2 {
									v285 = v304
									v286 = v306
									v287 = v287 + int32(4)
									continue
								} else {
									break
								}
								break
							}
							v315 = v304
						}
					} else {
						v261 = v66
						v262 = v172
						v285 = v261
						v286 = v262
						v287 = v59 + v262<<(uint(int32(2))%32)
						for {
							v300 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v287))))
							*(*uint8)(unsafe.Add(mBase, uint32(v63+v286))) = uint8(v300)
							v304 = v285 & v300
							v306 = v286 + int32(1)
							if v306 < l2 {
								v285 = v304
								v286 = v306
								v287 = v287 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						v315 = v304
					}
				}
				v332 = v75 + int32(1)
				if v332 != l3 {
					v59 = v59 + l1
					v63 = v63 + l5
					v66 = v315
					v74 = v173
					v75 = v332
					continue
				} else {
					break
				}
				break
			}
			v527 = base.Simd_g_i8x16_eq(v173, base.Simd_g_const(&F_ExtractAlpha_SSE2__k7))
			v528 = v315
		}
	} else {
		v527 = v22
		v528 = int32(255)
	}
	return base.B2i32(base.Simd_g_i8x16_bitmask(v527)&v528 == int32(255))
}

var F_ExtractAlpha_SSE2__k0 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}
var F_ExtractAlpha_SSE2__k1 = [2]uint64{0x100000000, 0x300000002}
var F_ExtractAlpha_SSE2__k2 = [2]uint64{0xffffffff00000000, 0xffffffffffffffff}
var F_ExtractAlpha_SSE2__k3 = [2]uint64{0x400000004, 0x400000004}
var F_ExtractAlpha_SSE2__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_ExtractAlpha_SSE2__k5 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ExtractAlpha_SSE2__k6 = [2]uint64{0x302010007060504, 0x302010003020100}
var F_ExtractAlpha_SSE2__k7 = [2]uint64{0xffffffffffffffff, 0x0}
var F_ExtractAlpha_SSE2__k8 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ExtractAlpha_SSE2__k9 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_ExtractAlpha_SSE2__k10 = [2]uint64{0x8080808080808080, 0x706050403020100}

func F_ExtractAlpha_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v30 int32
	_ = v30
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v72 int32
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 int32
	_ = v74
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v96 base.V128
	_ = v96
	var v98 int32
	_ = v98
	var v103 int32
	_ = v103
	var v104 base.V128
	_ = v104
	var v107 int32
	_ = v107
	var v108 base.V128
	_ = v108
	var v113 base.V128
	_ = v113
	var v118 base.V128
	_ = v118
	var v121 base.V128
	_ = v121
	var v124 base.V128
	_ = v124
	var v130 int32
	_ = v130
	var v149 base.V128
	_ = v149
	var v152 int32
	_ = v152
	var v162 int32
	_ = v162
	var v163 base.V128
	_ = v163
	var v164 int32
	_ = v164
	var v165 base.V128
	_ = v165
	var v166 int32
	_ = v166
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v179 base.V128
	_ = v179
	var v182 base.V128
	_ = v182
	var v185 base.V128
	_ = v185
	var v188 base.V128
	_ = v188
	var v194 base.V128
	_ = v194
	var v200 int32
	_ = v200
	var v203 base.V128
	_ = v203
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v217 int32
	_ = v217
	var v220 int32
	_ = v220
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v246 int32
	_ = v246
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v265 int32
	_ = v265
	var v273 int32
	_ = v273
	var v293 int32
	_ = v293
	var v300 int32
	_ = v300
	var v305 int32
	_ = v305
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v313 int32
	_ = v313
	var v331 int32
	_ = v331
	var v344 int32
	_ = v344
	var v354 base.V128
	_ = v354
	var v356 base.V128
	_ = v356
	var v358 int32
	_ = v358
	var v359 base.V128
	_ = v359
	var v360 int32
	_ = v360
	var v366 int32
	_ = v366
	var v369 int32
	_ = v369
	var v373 base.V128
	_ = v373
	var v376 base.V128
	_ = v376
	var v379 base.V128
	_ = v379
	var v382 base.V128
	_ = v382
	var v388 base.V128
	_ = v388
	var v392 int32
	_ = v392
	var v396 base.V128
	_ = v396
	var v401 int32
	_ = v401
	var v409 int32
	_ = v409
	var v412 int32
	_ = v412
	var v434 int32
	_ = v434
	var v437 int32
	_ = v437
	var v438 int32
	_ = v438
	var v451 int32
	_ = v451
	var v455 int32
	_ = v455
	var v457 int32
	_ = v457
	var v465 int32
	_ = v465
	var v484 int32
	_ = v484
	var v492 int32
	_ = v492
	var v493 base.V128
	_ = v493
	v23 = int32(255)
	v24 = base.Simd_g_const(&F_ExtractAlpha_SSE41__k0)
	if l3 < int32(1) {
		v492 = v23
		v493 = v24
	} else {
		v30 = (l2 + int32(-1)) & int32(-16)
		if v30 < int32(1) {
			if l2 < int32(1) {
				v492 = v23
				v493 = v24
			} else {
				v300 = l2 & int32(2147483644)
				v305 = l0
				v309 = l4
				v311 = int32(255)
				v313 = int32(0)
				for {
					if base.B2i32(base.Ui32(l2) < base.Ui32(int32(4))) == int32(0) {
						v331 = int32(0)
						v344 = v331
						v354 = base.Simd_g_const(&F_ExtractAlpha_SSE41__k1)
						v356 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_ExtractAlpha_SSE41__k2), v311)
						for {
							v358 = int32(2)
							v359 = base.Simd_g_i32x4_shl(v354, v358)
							v360 = int32(3)
							v366 = int32(1)
							v369 = int32(0)
							v373 = base.Simd_g_v128_load8_splat(m, v305+base.Simd_g_i32x4_extract_lane_l0(v359), v369)
							v376 = base.Simd_g_v128_load8_lane_l1(m, v305+base.Simd_g_i32x4_extract_lane_l1(v359), v369, v373)
							v379 = base.Simd_g_v128_load8_lane_l2(m, v305+base.Simd_g_i32x4_extract_lane_l2(v359), v369, v376)
							v382 = base.Simd_g_v128_load8_lane_l3(m, v305+base.Simd_g_i32x4_extract_lane_l3(v359), v369, v379)
							base.Simd_g_v128_store32_lane_l0(m, v309+v344, v369, v382)
							v388 = base.Simd_g_v128_and(v356, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v382)))
							v392 = v344 + int32(4)
							if v300 != v392 {
								v344 = v392
								v354 = base.Simd_g_i32x4_add(v354, base.Simd_g_const(&F_ExtractAlpha_SSE41__k3))
								v356 = v388
								continue
							} else {
								break
							}
							break
						}
						v396 = base.Simd_g_v128_and(v388, base.Simd_g_i8x16_shuffle2(v388, v382, base.Simd_g_const(&F_ExtractAlpha_SSE41__k4), base.Simd_g_const(&F_ExtractAlpha_SSE41__k5)))
						v401 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_v128_and(v396, base.Simd_g_i8x16_shuffle2(v396, v396, base.Simd_g_const(&F_ExtractAlpha_SSE41__k6), base.Simd_g_const(&F_ExtractAlpha_SSE41__k5))))
						if v300 == l2 {
							v465 = v401
						} else {
							v409 = v401
							v412 = v300
							v434 = v409
							v437 = v412
							v438 = v305 + v412<<(uint(int32(2))%32)
							for {
								v451 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v438))))
								*(*uint8)(unsafe.Add(mBase, uint32(v309+v437))) = uint8(v451)
								v455 = v434 & v451
								v457 = v437 + int32(1)
								if l2 != v457 {
									v434 = v455
									v437 = v457
									v438 = v438 + int32(4)
									continue
								} else {
									break
								}
								break
							}
							v465 = v455
						}
					} else {
						v409 = v311
						v412 = int32(0)
						v434 = v409
						v437 = v412
						v438 = v305 + v412<<(uint(int32(2))%32)
						for {
							v451 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v438))))
							*(*uint8)(unsafe.Add(mBase, uint32(v309+v437))) = uint8(v451)
							v455 = v434 & v451
							v457 = v437 + int32(1)
							if l2 != v457 {
								v434 = v455
								v437 = v457
								v438 = v438 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						v465 = v455
					}
					v484 = v313 + int32(1)
					if v484 != l3 {
						v305 = v305 + l1
						v309 = v309 + l5
						v311 = v465
						v313 = v484
						continue
					} else {
						break
					}
					break
				}
				v492 = v465
				v493 = v24
			}
		} else {
			v36 = (v30 + int32(-1)) & int32(-16)
			v38 = v36 + int32(17)
			if v38 < l2 {
				v40 = l2
			} else {
				v40 = v38
			}
			v42 = int32(-16)
			v43 = v40 - v36 + v42
			v45 = v40 & int32(3)
			v57 = l0
			v61 = l4
			v63 = int32(255)
			v72 = l4 + int32(16)
			v73 = base.Simd_g_const(&F_ExtractAlpha_SSE41__k0)
			v74 = int32(0)
			for {
				v89 = v57
				v90 = int32(0)
				v96 = v73
				v98 = v72
				for {
					v103 = int32(16)
					v104 = base.Simd_g_v128_load_rng(m, v89, v103, int32(0), int32(64))
					v107 = int32(0)
					v108 = base.Simd_g_v128_load_nc(m, v89, v107)
					v113 = base.Simd_g_v128_load_nc(m, v89, int32(32))
					v118 = base.Simd_g_v128_load_nc(m, v89, int32(48))
					v121 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v104, base.Simd_g_const(&F_ExtractAlpha_SSE41__k7)), base.Simd_g_i8x16_swizzle(v108, base.Simd_g_const(&F_ExtractAlpha_SSE41__k8))), base.Simd_g_i8x16_swizzle(v113, base.Simd_g_const(&F_ExtractAlpha_SSE41__k9))), base.Simd_g_i8x16_swizzle(v118, base.Simd_g_const(&F_ExtractAlpha_SSE41__k10)))
					base.Simd_g_v128_store(m, v61+v90, v107, v121)
					v124 = base.Simd_g_v128_and(v96, v121)
					v130 = v90 + v103
					if v130 < v30 {
						v89 = v89 + int32(64)
						v90 = v130
						v96 = v124
						v98 = v98 + v103
						continue
					} else {
						break
					}
					break
				}
				if l2 <= v130 {
					v273 = v63
				} else {
					if base.Ui32(int32(3)) < base.Ui32(v43) {
						v149 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_ExtractAlpha_SSE41__k2), v63)
						v152 = v43 - v45
						v162 = v98
						v163 = base.Simd_g_v128_or(base.Simd_g_i32x4_splat(v130), base.Simd_g_const(&F_ExtractAlpha_SSE41__k1))
						for {
							v164 = int32(2)
							v165 = base.Simd_g_i32x4_shl(v163, v164)
							v166 = int32(3)
							v172 = int32(1)
							v175 = int32(0)
							v179 = base.Simd_g_v128_load8_splat(m, v57+base.Simd_g_i32x4_extract_lane_l0(v165), v175)
							v182 = base.Simd_g_v128_load8_lane_l1(m, v57+base.Simd_g_i32x4_extract_lane_l1(v165), v175, v179)
							v185 = base.Simd_g_v128_load8_lane_l2(m, v57+base.Simd_g_i32x4_extract_lane_l2(v165), v175, v182)
							v188 = base.Simd_g_v128_load8_lane_l3(m, v57+base.Simd_g_i32x4_extract_lane_l3(v165), v175, v185)
							base.Simd_g_v128_store32_lane_l0(m, v162, v175, v188)
							v194 = base.Simd_g_v128_and(v149, base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v188)))
							v200 = v152 + int32(-4)
							if v200 != 0 {
								v149 = v194
								v152 = v200
								v162 = v162 + int32(4)
								v163 = base.Simd_g_i32x4_add(v163, base.Simd_g_const(&F_ExtractAlpha_SSE41__k3))
								continue
							} else {
								break
							}
							break
						}
						v203 = base.Simd_g_v128_and(v194, base.Simd_g_i8x16_shuffle2(v194, v124, base.Simd_g_const(&F_ExtractAlpha_SSE41__k4), base.Simd_g_const(&F_ExtractAlpha_SSE41__k5)))
						v207 = int32(0)
						v208 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_v128_and(v203, base.Simd_g_i8x16_shuffle2(v203, v124, base.Simd_g_const(&F_ExtractAlpha_SSE41__k6), base.Simd_g_const(&F_ExtractAlpha_SSE41__k5))))
						if v45 == v207 {
							v273 = v208
						} else {
							v217 = v208
							v220 = v40&int32(-4) - v36 + v42 + v130
							v242 = v217
							v245 = v220
							v246 = v57 + v220<<(uint(int32(2))%32)
							for {
								v259 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v246))))
								*(*uint8)(unsafe.Add(mBase, uint32(v61+v245))) = uint8(v259)
								v263 = v242 & v259
								v265 = v245 + int32(1)
								if v265 < l2 {
									v242 = v263
									v245 = v265
									v246 = v246 + int32(4)
									continue
								} else {
									break
								}
								break
							}
							v273 = v263
						}
					} else {
						v217 = v63
						v220 = v130
						v242 = v217
						v245 = v220
						v246 = v57 + v220<<(uint(int32(2))%32)
						for {
							v259 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v246))))
							*(*uint8)(unsafe.Add(mBase, uint32(v61+v245))) = uint8(v259)
							v263 = v242 & v259
							v265 = v245 + int32(1)
							if v265 < l2 {
								v242 = v263
								v245 = v265
								v246 = v246 + int32(4)
								continue
							} else {
								break
							}
							break
						}
						v273 = v263
					}
				}
				v293 = v74 + int32(1)
				if v293 != l3 {
					v57 = v57 + l1
					v61 = v61 + l5
					v63 = v273
					v72 = v72 + l5
					v73 = v124
					v74 = v293
					continue
				} else {
					break
				}
				break
			}
			v492 = v273
			v493 = base.Simd_g_i8x16_eq(v124, base.Simd_g_const(&F_ExtractAlpha_SSE41__k0))
		}
	}
	return base.B2i32(base.Simd_g_i8x16_bitmask(v493)&(v492|int32(_a_F_ExtractAlpha_SSE41_0)) == int32(_a_F_ExtractAlpha_SSE41_1))
}

var F_ExtractAlpha_SSE41__k0 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}
var F_ExtractAlpha_SSE41__k1 = [2]uint64{0x100000000, 0x300000002}
var F_ExtractAlpha_SSE41__k2 = [2]uint64{0xffffffff00000000, 0xffffffffffffffff}
var F_ExtractAlpha_SSE41__k3 = [2]uint64{0x400000004, 0x400000004}
var F_ExtractAlpha_SSE41__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_ExtractAlpha_SSE41__k5 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_ExtractAlpha_SSE41__k6 = [2]uint64{0x302010007060504, 0x302010003020100}
var F_ExtractAlpha_SSE41__k7 = [2]uint64{0xc0804008f8f8f8f, 0x8f8f8f8f8f8f8f8f}
var F_ExtractAlpha_SSE41__k8 = [2]uint64{0x8f8f8f8f0c080400, 0x8f8f8f8f8f8f8f8f}
var F_ExtractAlpha_SSE41__k9 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0x8f8f8f8f0c080400}
var F_ExtractAlpha_SSE41__k10 = [2]uint64{0x8f8f8f8f8f8f8f8f, 0xc0804008f8f8f8f}

func F_ExtractGreen_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v38 int32
	_ = v38
	var v43 int32
	_ = v43
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v70 int32
	_ = v70
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v43 = int32(0)
			v57 = l2 - v43
			v58 = l1 + v43
			v59 = l0 + v43<<(uint(int32(2))%32)
			for {
				v61 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
				v63 = int32(base.Ui32(v61) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v58))) = uint8(v63)
				v70 = v57 + int32(-1)
				if v70 != 0 {
					v57 = v70
					v58 = v58 + int32(1)
					v59 = v59 + int32(4)
					continue
				} else {
					break
				}
				break
			}
		} else {
			v15 = l2 & int32(2147483644)
			v20 = v15
			v21 = l1
			v22 = l0
			for {
				v24 = int32(0)
				v25 = base.Simd_g_v128_load(m, v22, v24)
				base.Simd_g_v128_store32_lane_l0(m, v21, v24, base.Simd_g_i8x16_swizzle_c(base.Simd_g_i32x4_shr_u(v25, int32(8)), base.Simd_g_const(&F_ExtractGreen_C__k0)))
				v38 = v20 + int32(-4)
				if v38 != 0 {
					v20 = v38
					v21 = v21 + int32(4)
					v22 = v22 + int32(16)
					continue
				} else {
					break
				}
				break
			}
			if v15 == l2 {
			} else {
				v43 = v15
				v57 = l2 - v43
				v58 = l1 + v43
				v59 = l0 + v43<<(uint(int32(2))%32)
				for {
					v61 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
					v63 = int32(base.Ui32(v61) >> (uint(int32(8)) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v58))) = uint8(v63)
					v70 = v57 + int32(-1)
					if v70 != 0 {
						v57 = v70
						v58 = v58 + int32(1)
						v59 = v59 + int32(4)
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

var F_ExtractGreen_C__k0 = [2]uint64{0xc080400, 0x0}

func F_ExtractGreen_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v28 int32
	_ = v28
	var v30 base.V128
	_ = v30
	var v32 int32
	_ = v32
	var v33 base.V128
	_ = v33
	var v38 int32
	_ = v38
	var v39 base.V128
	_ = v39
	var v44 base.V128
	_ = v44
	var v49 base.V128
	_ = v49
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v65 base.V128
	_ = v65
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v74 base.V128
	_ = v74
	var v75 int32
	_ = v75
	var v77 base.V128
	_ = v77
	var v80 base.V128
	_ = v80
	var v84 base.V128
	_ = v84
	var v89 base.V128
	_ = v89
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v100 int32
	_ = v100
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v113 base.V128
	_ = v113
	var v126 int32
	_ = v126
	var v136 int32
	_ = v136
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v155 int32
	_ = v155
	var v162 int32
	_ = v162
	if l2 < int32(16) {
		v62 = int32(0)
		v63 = l0
		v65 = base.Simd_g_const(&F_ExtractGreen_SSE2__k0)
	} else {
		v19 = l0
		v20 = int32(0)
		for {
			v26 = int32(0)
			v27 = base.Simd_g_v128_load_rng(m, v19, v26, int32(0), int32(64))
			v28 = int32(8)
			v30 = base.Simd_g_const(&F_ExtractGreen_SSE2__k1)
			v32 = int32(16)
			v33 = base.Simd_g_v128_load_nc(m, v19, v32)
			v38 = int32(32)
			v39 = base.Simd_g_v128_load_nc(m, v19, v38)
			v44 = base.Simd_g_v128_load_nc(m, v19, int32(48))
			v49 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v27, v28), v30), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v33, v28), v30)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v39, v28), v30), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v44, v28), v30)))
			base.Simd_g_v128_store(m, l1+v20, v26, v49)
			v53 = v19 + int32(64)
			v57 = v20 + v32
			if v20+v38 <= l2 {
				v19 = v53
				v20 = v57
				continue
			} else {
				break
			}
			break
		}
		v62 = v57
		v63 = v53
		v65 = v30
	}
	v70 = v62 | int32(8)
	if v70 <= l2 {
		v73 = int32(0)
		v74 = base.Simd_g_v128_load_rng(m, v63, v73, int32(0), int32(32))
		v75 = int32(8)
		v77 = base.Simd_g_const(&F_ExtractGreen_SSE2__k1)
		v80 = base.Simd_g_v128_load_nc(m, v63, int32(16))
		v84 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v74, v75), v77), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v80, v75), v77))
		base.Simd_g_v128_store64_lane_l0(m, l1+v62, v73, base.Simd_g_i8x16_narrow_i16x8_u(v84, v84))
		v89 = v84
		v90 = v70
	} else {
		v89 = v65
		v90 = v62
	}
	if l2 <= v90 {
	} else {
		v92 = l2 - v90
		if base.Ui32(v92) < base.Ui32(int32(4)) {
			v136 = v90
			v147 = l1 + v136
			v148 = l0 + v136<<(uint(int32(2))%32)
			v150 = l2 - v136
			for {
				v153 = *(*int32)(unsafe.Add(mBase, uint32(v148)))
				v155 = int32(base.Ui32(v153) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v147))) = uint8(v155)
				v162 = v150 + int32(-1)
				if v162 != 0 {
					v147 = v147 + int32(1)
					v148 = v148 + int32(4)
					v150 = v162
					continue
				} else {
					break
				}
				break
			}
		} else {
			v100 = v92 & int32(-4)
			v106 = l1 + v90
			v107 = l0 + v90<<(uint(int32(2))%32)
			v109 = v100
			for {
				v112 = int32(0)
				v113 = base.Simd_g_v128_load(m, v107, v112)
				base.Simd_g_v128_store32_lane_l0(m, v106, v112, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(v113, int32(8)), v89, base.Simd_g_const(&F_ExtractGreen_SSE2__k2), base.Simd_g_const(&F_ExtractGreen_SSE2__k3)))
				v126 = v109 + int32(-4)
				if v126 != 0 {
					v106 = v106 + int32(4)
					v107 = v107 + int32(16)
					v109 = v126
					continue
				} else {
					break
				}
				break
			}
			if v92 == v100 {
			} else {
				v136 = v90 + v100
				v147 = l1 + v136
				v148 = l0 + v136<<(uint(int32(2))%32)
				v150 = l2 - v136
				for {
					v153 = *(*int32)(unsafe.Add(mBase, uint32(v148)))
					v155 = int32(base.Ui32(v153) >> (uint(int32(8)) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v147))) = uint8(v155)
					v162 = v150 + int32(-1)
					if v162 != 0 {
						v147 = v147 + int32(1)
						v148 = v148 + int32(4)
						v150 = v162
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

var F_ExtractGreen_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_ExtractGreen_SSE2__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_ExtractGreen_SSE2__k2 = [2]uint64{0xc080400, 0x0}
var F_ExtractGreen_SSE2__k3 = [2]uint64{0x8080808080808080, 0x8080808080808080}

func F_encode(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v31 int64
	_ = v31
	var v36 base.V128
	_ = v36
	var v39 base.V128
	_ = v39
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v53 base.V128
	_ = v53
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v98 float32
	_ = v98
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v107 float32
	_ = v107
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v201 int32
	_ = v201
	var v212 int32
	_ = v212
	var v214 int32
	_ = v214
	var v215 int32
	_ = v215
	var v220 int32
	_ = v220
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v238 int32
	_ = v238
	var v240 int32
	_ = v240
	var v242 int64
	_ = v242
	var v250 base.V128
	_ = v250
	var v251 int32
	_ = v251
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v273 int32
	_ = v273
	var v275 int32
	_ = v275
	var v277 int64
	_ = v277
	var v285 base.V128
	_ = v285
	var v286 int32
	_ = v286
	var v294 int32
	_ = v294
	var v301 int32
	_ = v301
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	v9 = int32(0)
	v11 = m.G0
	v13 = v11 - int32(304)
	m.G0 = v13
	*(*int32)(unsafe.Add(mBase, uint32(l3))) = v9
	v19 = v13 + int32(188)
	if v19 == v9 {
		v81 = v9
	} else {
		v31 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+16)) = v31
		*(*float32)(unsafe.Add(mBase, uint32(v19)+4)) = float32(75)
		*(*int32)(unsafe.Add(mBase, uint32(v19)+112)) = int32(100)
		v36 = base.Simd_g_const(&F_encode__k0)
		base.Simd_g_v128_store(m, v19, int32(24), v36)
		v39 = base.Simd_g_const(&F_encode__k1)
		base.Simd_g_v128_store(m, v19, int32(56), v39)
		v42 = base.Simd_g_const(&F_encode__k2)
		base.Simd_g_v128_store(m, v19, int32(40), v42)
		v45 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v19)+96)) = v45
		*(*int32)(unsafe.Add(mBase, uint32(v19))) = v45
		*(*int64)(unsafe.Add(mBase, uint32(v19)+104)) = v31
		*(*int64)(unsafe.Add(mBase, uint32(v19)+88)) = int64(429496729600)
		v53 = base.Simd_g_const(&F_encode__k3)
		base.Simd_g_v128_store(m, v19, int32(72), v53)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+8)) = int64(4)
		switch int32(-1) {
		case 0:
			*(*int32)(unsafe.Add(mBase, uint32(v19)+36)) = int32(4)
			*(*int64)(unsafe.Add(mBase, uint32(v19)+28)) = int64(150323855440)
		case 1:
			*(*int32)(unsafe.Add(mBase, uint32(v19)+36)) = int32(3)
			*(*int64)(unsafe.Add(mBase, uint32(v19)+28)) = int64(128849018960)
			*(*int32)(unsafe.Add(mBase, uint32(v19)+68)) = int32(2)
		case 2:
			*(*int32)(unsafe.Add(mBase, uint32(v19)+36)) = int32(6)
			*(*int64)(unsafe.Add(mBase, uint32(v19)+28)) = int64(42949672985)
		case 3:
			*(*int64)(unsafe.Add(mBase, uint32(v19)+28)) = int64(0)
		case 4:
			*(*int32)(unsafe.Add(mBase, uint32(v19)+32)) = int32(0)
			*(*int64)(unsafe.Add(mBase, uint32(v19)+24)) = int64(2)
		default:
		}
		v80 = F_WebPValidateConfig(m, v19)
		mBase = m.M
		v81 = v80
	}
	if v81 == int32(0) {
		v304 = v9
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(v13)+196)) = l5
		*(*int32)(unsafe.Add(mBase, uint32(v13)+284)) = l7
		*(*int32)(unsafe.Add(mBase, uint32(v13)+188)) = l6
		*(*float32)(unsafe.Add(mBase, uint32(v13)+192)) = base.F32_convert_i32_s(l4)
		v90 = v13 + int32(188)
		v91 = int32(0)
		if v90 == v91 {
			v193 = v91
		} else {
			v98 = *(*float32)(unsafe.Add(mBase, uint32(v90)+4))
			if base.F32_lt(v98, float32(0)) != 0 {
				v193 = v91
			} else {
				if base.F32_gt(v98, float32(100)) != 0 {
					v193 = v91
				} else {
					v103 = int32(0)
					v104 = *(*int32)(unsafe.Add(mBase, uint32(v90)+16))
					if v104 < v103 {
						v193 = v103
					} else {
						v107 = *(*float32)(unsafe.Add(mBase, uint32(v90)+20))
						if base.F32_lt(v107, float32(0)) != 0 {
							v193 = v103
						} else {
							v110 = *(*int32)(unsafe.Add(mBase, uint32(v90)+8))
							if base.Ui32(int32(6)) < base.Ui32(v110) {
								v193 = v103
							} else {
								v113 = *(*int32)(unsafe.Add(mBase, uint32(v90)+24))
								if base.Ui32(v113+int32(-5)) < base.Ui32(int32(-4)) {
									v193 = v103
								} else {
									v118 = *(*int32)(unsafe.Add(mBase, uint32(v90)+28))
									if base.Ui32(int32(100)) < base.Ui32(v118) {
										v193 = v103
									} else {
										v121 = *(*int32)(unsafe.Add(mBase, uint32(v90)+32))
										if base.Ui32(int32(100)) < base.Ui32(v121) {
											v193 = v103
										} else {
											v124 = *(*int32)(unsafe.Add(mBase, uint32(v90)+36))
											if base.Ui32(int32(7)) < base.Ui32(v124) {
												v193 = v103
											} else {
												v127 = *(*int32)(unsafe.Add(mBase, uint32(v90)+40))
												if base.Ui32(int32(1)) < base.Ui32(v127) {
													v193 = v103
												} else {
													v130 = *(*int32)(unsafe.Add(mBase, uint32(v90)+44))
													if base.Ui32(int32(1)) < base.Ui32(v130) {
														v193 = v103
													} else {
														v133 = *(*int32)(unsafe.Add(mBase, uint32(v90)+60))
														if base.Ui32(v133+int32(-11)) < base.Ui32(int32(-10)) {
															v193 = v103
														} else {
															v138 = int32(0)
															v139 = *(*int32)(unsafe.Add(mBase, uint32(v90)+108))
															if v139 < v138 {
																v193 = v138
															} else {
																v142 = *(*int32)(unsafe.Add(mBase, uint32(v90)+112))
																if int32(100) < v142 {
																	v193 = v138
																} else {
																	if v142 < v139 {
																		v193 = v138
																	} else {
																		v146 = *(*int32)(unsafe.Add(mBase, uint32(v90)+64))
																		if base.Ui32(int32(1)) < base.Ui32(v146) {
																			v193 = v138
																		} else {
																			v149 = *(*int32)(unsafe.Add(mBase, uint32(v90)+68))
																			if base.Ui32(int32(7)) < base.Ui32(v149) {
																				v193 = v138
																			} else {
																				v152 = *(*int32)(unsafe.Add(mBase, uint32(v90)+72))
																				if base.Ui32(int32(3)) < base.Ui32(v152) {
																					v193 = v138
																				} else {
																					v155 = *(*int32)(unsafe.Add(mBase, uint32(v90)+76))
																					if base.Ui32(int32(100)) < base.Ui32(v155) {
																						v193 = v138
																					} else {
																						v158 = int32(0)
																						v159 = *(*int32)(unsafe.Add(mBase, uint32(v90)+48))
																						if v159 < v158 {
																							v193 = v158
																						} else {
																							v162 = int32(0)
																							v163 = *(*int32)(unsafe.Add(mBase, uint32(v90)+52))
																							if v163 < v162 {
																								v193 = v162
																							} else {
																								v166 = *(*int32)(unsafe.Add(mBase, uint32(v90)+56))
																								if base.Ui32(int32(100)) < base.Ui32(v166) {
																									v193 = v162
																								} else {
																									v169 = *(*int32)(unsafe.Add(mBase, uint32(v90)))
																									if base.Ui32(int32(1)) < base.Ui32(v169) {
																										v193 = v162
																									} else {
																										v172 = *(*int32)(unsafe.Add(mBase, uint32(v90)+92))
																										if base.Ui32(int32(100)) < base.Ui32(v172) {
																											v193 = v162
																										} else {
																											v175 = *(*int32)(unsafe.Add(mBase, uint32(v90)+12))
																											if base.Ui32(int32(3)) < base.Ui32(v175) {
																												v193 = v162
																											} else {
																												v178 = *(*int32)(unsafe.Add(mBase, uint32(v90)+80))
																												if base.Ui32(int32(1)) < base.Ui32(v178) {
																													v193 = v162
																												} else {
																													v181 = *(*int32)(unsafe.Add(mBase, uint32(v90)+84))
																													if base.Ui32(int32(1)) < base.Ui32(v181) {
																														v193 = v162
																													} else {
																														v184 = *(*int32)(unsafe.Add(mBase, uint32(v90)+88))
																														if base.Ui32(int32(1)) < base.Ui32(v184) {
																															v193 = v162
																														} else {
																															v187 = *(*int32)(unsafe.Add(mBase, uint32(v90)+96))
																															if base.Ui32(int32(1)) < base.Ui32(v187) {
																																v193 = v162
																															} else {
																																v190 = *(*int32)(unsafe.Add(mBase, uint32(v90)+104))
																																v193 = base.B2i32(base.Ui32(v190) < base.Ui32(int32(2)))
																															}
																														}
																													}
																												}
																											}
																										}
																									}
																								}
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if v193 == int32(0) {
			v304 = v9
		} else {
			v201 = v13 + int32(16)
			if v201 == int32(0) {
			} else {
				v212 = int32(0)
				v214 = F_memset(m, v201, v212, int32(172))
				mBase = m.M
				v215 = m.G2
				*(*int32)(unsafe.Add(mBase, uint32(v214)+72)) = v215 + int32(1)
				v220 = F_WebPEncodingSetError(m, v214, v212)
				mBase = m.M
			}
			*(*int32)(unsafe.Add(mBase, uint32(v13)+28)) = l2
			*(*int32)(unsafe.Add(mBase, uint32(v13)+16)) = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v13)+24)) = l1
			v233 = F_WebPPictureImportRGBA(m, v13+int32(16), l0, l1<<(uint(int32(2))%32))
			mBase = m.M
			if v233 != 0 {
				*(*int32)(unsafe.Add(mBase, uint32(v13)+8)) = int32(0)
				*(*int64)(unsafe.Add(mBase, uint32(v13))) = int64(0)
				*(*int32)(unsafe.Add(mBase, uint32(v13)+88)) = int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v13)+92)) = v13
				v267 = v13 + int32(16)
				v268 = F_WebPEncode(m, v13+int32(188), v267)
				mBase = m.M
				if v267 == int32(0) {
				} else {
					v273 = *(*int32)(unsafe.Add(mBase, uint32(v267)+156))
					F_WebPSafeFree(m, v273)
					mBase = m.M
					v275 = *(*int32)(unsafe.Add(mBase, uint32(v267)+160))
					F_WebPSafeFree(m, v275)
					mBase = m.M
					v277 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v267)+52)) = v277
					*(*int64)(unsafe.Add(mBase, uint32(v267)+156)) = v277
					*(*int64)(unsafe.Add(mBase, uint32(v267)+16)) = v277
					v285 = base.Simd_g_const(&F_encode__k3)
					v286 = int32(0)
					base.Simd_g_v128_store(m, v13+int32(40), v286, v285)
					*(*int32)(unsafe.Add(mBase, uint32(v13+int32(56)))) = v286
				}
				if v268 != 0 {
					v301 = *(*int32)(unsafe.Add(mBase, uint32(v13)+4))
					*(*int32)(unsafe.Add(mBase, uint32(l3))) = v301
					v303 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
					v304 = v303
				} else {
					if v13 == int32(0) {
					} else {
						v294 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
						F_WebPSafeFree(m, v294)
						mBase = m.M
						*(*int32)(unsafe.Add(mBase, uint32(v13)+8)) = int32(0)
						*(*int64)(unsafe.Add(mBase, uint32(v13))) = int64(0)
					}
					v304 = int32(0)
				}
			} else {
				v235 = v13 + int32(16)
				if v235 == int32(0) {
				} else {
					v238 = *(*int32)(unsafe.Add(mBase, uint32(v235)+156))
					F_WebPSafeFree(m, v238)
					mBase = m.M
					v240 = *(*int32)(unsafe.Add(mBase, uint32(v235)+160))
					F_WebPSafeFree(m, v240)
					mBase = m.M
					v242 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v235)+52)) = v242
					*(*int64)(unsafe.Add(mBase, uint32(v235)+156)) = v242
					*(*int64)(unsafe.Add(mBase, uint32(v235)+16)) = v242
					v250 = base.Simd_g_const(&F_encode__k3)
					v251 = int32(0)
					base.Simd_g_v128_store(m, v13+int32(40), v251, v250)
					*(*int32)(unsafe.Add(mBase, uint32(v13+int32(56)))) = v251
				}
				v304 = int32(0)
			}
		}
	}
	m.G0 = v13 + int32(304)
	return v304
}

var F_encode__k0 = [2]uint64{0x3200000004, 0x3c}
var F_encode__k1 = [2]uint64{0x100000064, 0x0}
var F_encode__k2 = [2]uint64{0x1, 0x100000001}
var F_encode__k3 = [2]uint64{0x0, 0x0}

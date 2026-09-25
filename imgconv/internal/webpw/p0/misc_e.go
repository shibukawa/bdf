//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"math"
	"unsafe"
)

func F_EncodeAlphaInternal(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) int32 {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v33 int64
	_ = v33
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v59 int32
	_ = v59
	var v63 int32
	_ = v63
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v101 int64
	_ = v101
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v145 int64
	_ = v145
	var v152 int64
	_ = v152
	var v166 int32
	_ = v166
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v209 float32
	_ = v209
	var v213 float32
	_ = v213
	var v214 float32
	_ = v214
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v227 int32
	_ = v227
	var v229 int32
	_ = v229
	var v231 int64
	_ = v231
	var v251 int32
	_ = v251
	var v255 int32
	_ = v255
	var v260 int32
	_ = v260
	var v261 int32
	_ = v261
	var v265 int64
	_ = v265
	var v276 int32
	_ = v276
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v286 int64
	_ = v286
	var v293 int64
	_ = v293
	var v309 int32
	_ = v309
	var v318 int32
	_ = v318
	var v319 int32
	_ = v319
	var v320 int32
	_ = v320
	var v322 int32
	_ = v322
	var v328 int64
	_ = v328
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v343 int32
	_ = v343
	var v345 int32
	_ = v345
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	var v353 int32
	_ = v353
	var v355 int32
	_ = v355
	var v356 int32
	_ = v356
	var v357 int32
	_ = v357
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v369 int32
	_ = v369
	var v372 int32
	_ = v372
	var v379 int32
	_ = v379
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v386 int32
	_ = v386
	var v388 int32
	_ = v388
	var v405 int32
	_ = v405
	var v415 int32
	_ = v415
	var v416 int32
	_ = v416
	var v420 int32
	_ = v420
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v430 int64
	_ = v430
	var v437 int64
	_ = v437
	var v452 int32
	_ = v452
	var v457 int32
	_ = v457
	var v458 int32
	_ = v458
	var v460 int32
	_ = v460
	var v465 int32
	_ = v465
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v475 int64
	_ = v475
	var v485 int32
	_ = v485
	var v486 int32
	_ = v486
	var v488 int32
	_ = v488
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v491 int32
	_ = v491
	var v495 int32
	_ = v495
	var v498 int32
	_ = v498
	var v501 int32
	_ = v501
	var v503 int32
	_ = v503
	var v506 int64
	_ = v506
	var v519 int32
	_ = v519
	var v522 int32
	_ = v522
	var v523 int32
	_ = v523
	var v527 int32
	_ = v527
	var v532 int32
	_ = v532
	var v533 int32
	_ = v533
	var v534 int32
	_ = v534
	var v545 int32
	_ = v545
	var v550 int32
	_ = v550
	var v556 int32
	_ = v556
	var v559 int32
	_ = v559
	var v562 int64
	_ = v562
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v571 int32
	_ = v571
	var v574 int32
	_ = v574
	var v576 int32
	_ = v576
	var v577 int32
	_ = v577
	var v580 int32
	_ = v580
	var v581 int32
	_ = v581
	var v585 int32
	_ = v585
	var v588 int32
	_ = v588
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v597 int32
	_ = v597
	var v599 int32
	_ = v599
	var v600 int32
	_ = v600
	var v604 int32
	_ = v604
	var v612 int32
	_ = v612
	var v620 int32
	_ = v620
	var v623 int32
	_ = v623
	var v626 int64
	_ = v626
	var v632 int32
	_ = v632
	var v633 int32
	_ = v633
	var v635 int32
	_ = v635
	var v638 int32
	_ = v638
	var v640 int32
	_ = v640
	var v641 int32
	_ = v641
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v649 int32
	_ = v649
	var v652 int32
	_ = v652
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v663 int32
	_ = v663
	var v664 int32
	_ = v664
	var v668 int32
	_ = v668
	var v676 int32
	_ = v676
	var v679 int32
	_ = v679
	var v681 int32
	_ = v681
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v691 int64
	_ = v691
	var v698 int32
	_ = v698
	var v701 int32
	_ = v701
	var v704 int32
	_ = v704
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	v13 = m.G0
	v15 = v13 - int32(320)
	m.G0 = v15
	v17 = m.G74
	v21 = *(*int32)(unsafe.Add(mBase, uint32(v17+l4<<(uint(int32(2))%32))))
	if v21 != 0 {
		m.T0[v21].(func(*base.Module, int32, int32, int32, int32, int32))(m, l0, l1, l2, l1, l7)
		mBase = m.M
		v23 = l7
	} else {
		v23 = l0
	}
	v24 = l2 * l1
	v25 = int32(0)
	if l3 != 0 {
		v29 = v15 + int32(8)
		v33 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(v29))) = v33
		*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v33
		*(*int64)(unsafe.Add(mBase, uint32(v15+int32(16)))) = v33
		v47 = int32(base.Ui32(v24)>>(uint(int32(3))%32))&int32(-1024) + int32(1024)
		v48 = F_WebPSafeMalloc(m, int64(1), v47)
		mBase = m.M
		if v48 != 0 {
			v52 = *(*int32)(unsafe.Add(mBase, uint32(v29)+8))
			F_WebPSafeFree(m, v52)
			mBase = m.M
			*(*int32)(unsafe.Add(mBase, uint32(v29)+16)) = v48 + v47
			*(*int32)(unsafe.Add(mBase, uint32(v29)+12)) = v48
			*(*int32)(unsafe.Add(mBase, uint32(v29)+8)) = v48
			v59 = int32(1)
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(v29)+20)) = int32(1)
			v59 = int32(0)
		}
		if v59 == int32(0) {
			v276 = v15 + int32(8)
			if v276 == int32(0) {
			} else {
				v281 = v15 + int32(16)
				v282 = *(*int32)(unsafe.Add(mBase, uint32(v281)))
				F_WebPSafeFree(m, v282)
				mBase = m.M
				v286 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v286
				*(*int64)(unsafe.Add(mBase, uint32(v281))) = v286
				*(*int64)(unsafe.Add(mBase, uint32(v276))) = v286
			}
			v293 = int64(0)
			*(*int64)(unsafe.Add(mBase, uint32(l8)+4)) = v293
			*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v293
			*(*int64)(unsafe.Add(mBase, uint32(l8+int32(20)))) = v293
			*(*int64)(unsafe.Add(mBase, uint32(l8+int32(12)))) = v293
			v709 = int32(0)
		} else {
			v63 = v15 + int32(32)
			if v63 == int32(0) {
			} else {
				v74 = int32(0)
				v76 = F_memset(m, v63, v74, int32(172))
				mBase = m.M
				v77 = m.G2
				*(*int32)(unsafe.Add(mBase, uint32(v76)+72)) = v77 + int32(1)
				v82 = F_WebPEncodingSetError(m, v76, v74)
				mBase = m.M
			}
			*(*int32)(unsafe.Add(mBase, uint32(v15)+44)) = l2
			*(*int32)(unsafe.Add(mBase, uint32(v15)+40)) = l1
			*(*int32)(unsafe.Add(mBase, uint32(v15)+32)) = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v15)+120)) = l8 + int32(36)
			v95 = v15 + int32(32)
			if v95 != 0 {
				v97 = *(*int32)(unsafe.Add(mBase, uint32(v95)+156))
				F_WebPSafeFree(m, v97)
				mBase = m.M
				v99 = *(*int32)(unsafe.Add(mBase, uint32(v95)+160))
				F_WebPSafeFree(m, v99)
				mBase = m.M
				v101 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(v95)+52)) = v101
				*(*int64)(unsafe.Add(mBase, uint32(v95)+156)) = v101
				*(*int64)(unsafe.Add(mBase, uint32(v95)+16)) = v101
				*(*int64)(unsafe.Add(mBase, uint32(v15+int32(56)))) = v101
				*(*int64)(unsafe.Add(mBase, uint32(v15+int32(64)))) = v101
				*(*int32)(unsafe.Add(mBase, uint32(v15+int32(72)))) = int32(0)
				v119 = *(*int32)(unsafe.Add(mBase, uint32(v95)))
				if v119 != 0 {
					v121 = F_WebPPictureAllocARGB(m, v95)
					mBase = m.M
					v122 = v121
				} else {
					v120 = F_WebPPictureAllocYUVA(m, v95)
					mBase = m.M
					v122 = v120
				}
			} else {
				v122 = int32(1)
			}
			if v122 == int32(0) {
				v276 = v15 + int32(8)
				if v276 == int32(0) {
				} else {
					v281 = v15 + int32(16)
					v282 = *(*int32)(unsafe.Add(mBase, uint32(v281)))
					F_WebPSafeFree(m, v282)
					mBase = m.M
					v286 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v286
					*(*int64)(unsafe.Add(mBase, uint32(v281))) = v286
					*(*int64)(unsafe.Add(mBase, uint32(v276))) = v286
				}
				v293 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(l8)+4)) = v293
				*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v293
				*(*int64)(unsafe.Add(mBase, uint32(l8+int32(20)))) = v293
				*(*int64)(unsafe.Add(mBase, uint32(l8+int32(12)))) = v293
				v709 = int32(0)
			} else {
				v125 = m.G75
				v126 = *(*int32)(unsafe.Add(mBase, uint32(v15)+40))
				v127 = *(*int32)(unsafe.Add(mBase, uint32(v15)+44))
				v128 = *(*int32)(unsafe.Add(mBase, uint32(v15)+84))
				v129 = *(*int32)(unsafe.Add(mBase, uint32(v15)+88))
				v130 = *(*int32)(unsafe.Add(mBase, uint32(v125)))
				m.T0[v130].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v23, l1, v126, v127, v128, v129)
				mBase = m.M
				v133 = v15 + int32(204)
				v134 = int32(0)
				if v133 == v134 {
					v199 = v134
				} else {
					v145 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+16)) = v145
					*(*float32)(unsafe.Add(mBase, uint32(v133)+4)) = float32(75)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+32)) = int64(60)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+8)) = int64(4)
					v152 = int64(429496729600)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+108)) = v152
					*(*int64)(unsafe.Add(mBase, uint32(v133)+24)) = int64(214748364804)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+64)) = v145
					*(*int64)(unsafe.Add(mBase, uint32(v133)+72)) = v145
					*(*int64)(unsafe.Add(mBase, uint32(v133)+40)) = int64(1)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+56)) = int64(4294967396)
					*(*int64)(unsafe.Add(mBase, uint32(v133)+48)) = int64(4294967297)
					v166 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(v133))) = v166
					*(*int32)(unsafe.Add(mBase, uint32(v133)+96)) = v166
					*(*int64)(unsafe.Add(mBase, uint32(v133)+88)) = v152
					*(*int64)(unsafe.Add(mBase, uint32(v133)+80)) = v145
					*(*int32)(unsafe.Add(mBase, uint32(v133)+104)) = v166
					switch int32(-1) {
					case 0:
						*(*int32)(unsafe.Add(mBase, uint32(v133)+36)) = int32(4)
						*(*int64)(unsafe.Add(mBase, uint32(v133)+28)) = int64(150323855440)
					case 1:
						*(*int32)(unsafe.Add(mBase, uint32(v133)+36)) = int32(3)
						*(*int64)(unsafe.Add(mBase, uint32(v133)+28)) = int64(128849018960)
						*(*int32)(unsafe.Add(mBase, uint32(v133)+68)) = int32(2)
					case 2:
						*(*int32)(unsafe.Add(mBase, uint32(v133)+36)) = int32(6)
						*(*int64)(unsafe.Add(mBase, uint32(v133)+28)) = int64(42949672985)
					case 3:
						*(*int64)(unsafe.Add(mBase, uint32(v133)+28)) = int64(0)
					case 4:
						*(*int32)(unsafe.Add(mBase, uint32(v133)+32)) = int32(0)
						*(*int64)(unsafe.Add(mBase, uint32(v133)+24)) = int64(2)
					default:
					}
					v198 = F_WebPValidateConfig(m, v133)
					mBase = m.M
					v199 = v198
				}
				if v199 == int32(0) {
					v276 = v15 + int32(8)
					if v276 == int32(0) {
					} else {
						v281 = v15 + int32(16)
						v282 = *(*int32)(unsafe.Add(mBase, uint32(v281)))
						F_WebPSafeFree(m, v282)
						mBase = m.M
						v286 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v286
						*(*int64)(unsafe.Add(mBase, uint32(v281))) = v286
						*(*int64)(unsafe.Add(mBase, uint32(v276))) = v286
					}
					v293 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(l8)+4)) = v293
					*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v293
					*(*int64)(unsafe.Add(mBase, uint32(l8+int32(20)))) = v293
					*(*int64)(unsafe.Add(mBase, uint32(l8+int32(12)))) = v293
					v709 = int32(0)
				} else {
					v202 = int32(1)
					*(*int32)(unsafe.Add(mBase, uint32(v15)+300)) = v202
					*(*int32)(unsafe.Add(mBase, uint32(v15)+204)) = v202
					*(*int32)(unsafe.Add(mBase, uint32(v15)+212)) = l6
					v209 = base.F32_mul(base.F32_convert_i32_s(l6), float32(8))
					if l6 == int32(6) {
						v213 = float32(100)
					} else {
						v213 = v209
					}
					if l5 != 0 {
						v214 = v209
					} else {
						v214 = v213
					}
					*(*float32)(unsafe.Add(mBase, uint32(v15)+208)) = v214
					v219 = v15 + int32(32)
					v222 = F_VP8LEncodeStream(m, v15+int32(204), v219, v15+int32(8))
					mBase = m.M
					if v219 == int32(0) {
					} else {
						v227 = *(*int32)(unsafe.Add(mBase, uint32(v219)+156))
						F_WebPSafeFree(m, v227)
						mBase = m.M
						v229 = *(*int32)(unsafe.Add(mBase, uint32(v219)+160))
						F_WebPSafeFree(m, v229)
						mBase = m.M
						v231 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(v219)+52)) = v231
						*(*int64)(unsafe.Add(mBase, uint32(v219)+156)) = v231
						*(*int64)(unsafe.Add(mBase, uint32(v219)+16)) = v231
						*(*int64)(unsafe.Add(mBase, uint32(v15+int32(56)))) = v231
						*(*int64)(unsafe.Add(mBase, uint32(v15+int32(64)))) = v231
						*(*int32)(unsafe.Add(mBase, uint32(v15+int32(72)))) = int32(0)
					}
					if v222 == int32(0) {
						v255 = v15 + int32(8)
						if v255 == int32(0) {
						} else {
							v260 = v15 + int32(16)
							v261 = *(*int32)(unsafe.Add(mBase, uint32(v260)))
							F_WebPSafeFree(m, v261)
							mBase = m.M
							v265 = int64(0)
							*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v265
							*(*int64)(unsafe.Add(mBase, uint32(v260))) = v265
							*(*int64)(unsafe.Add(mBase, uint32(v255))) = v265
						}
						v276 = v15 + int32(8)
						if v276 == int32(0) {
						} else {
							v281 = v15 + int32(16)
							v282 = *(*int32)(unsafe.Add(mBase, uint32(v281)))
							F_WebPSafeFree(m, v282)
							mBase = m.M
							v286 = int64(0)
							*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v286
							*(*int64)(unsafe.Add(mBase, uint32(v281))) = v286
							*(*int64)(unsafe.Add(mBase, uint32(v276))) = v286
						}
						v293 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(l8)+4)) = v293
						*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v293
						*(*int64)(unsafe.Add(mBase, uint32(l8+int32(20)))) = v293
						*(*int64)(unsafe.Add(mBase, uint32(l8+int32(12)))) = v293
						v709 = int32(0)
					} else {
						v251 = *(*int32)(unsafe.Add(mBase, uint32(v15)+28))
						if v251 == int32(0) {
							v309 = v15 + int32(8)
							v318 = *(*int32)(unsafe.Add(mBase, uint32(v309)+12))
							v319 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
							v320 = v318 - v319
							v322 = *(*int32)(unsafe.Add(mBase, uint32(v309)+4))
							v328 = base.I64_extend_i32_u(v320) + base.I64_extend_i32_u((v322+int32(7))>>(uint(int32(3))%32))
							if base.Ui64(v328) < base.Ui64(int64(4294967296)) {
								v334 = base.I32_wrap_i64(v328)
								v335 = *(*int32)(unsafe.Add(mBase, uint32(v309)+16))
								v336 = v335 - v319
								if v335 == v319 {
									v343 = int32(base.Ui32(v336*int32(3)) >> (uint(int32(1)) % 32))
									if base.Ui32(v334) < base.Ui32(v343) {
										v345 = v343
									} else {
										v345 = v334
									}
									v349 = v345&int32(-1024) + int32(1024)
									v350 = F_WebPSafeMalloc(m, int64(1), v349)
									mBase = m.M
									if v350 != 0 {
										if v318 == v319 {
										} else {
											v355 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
											v356 = F_memcpy(m, v350, v355, v320)
											mBase = m.M
										}
										v357 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
										F_WebPSafeFree(m, v357)
										mBase = m.M
										*(*int32)(unsafe.Add(mBase, uint32(v309)+16)) = v350 + v349
										*(*int32)(unsafe.Add(mBase, uint32(v309)+12)) = v350 + v320
										*(*int32)(unsafe.Add(mBase, uint32(v309)+8)) = v350
										v364 = *(*int32)(unsafe.Add(mBase, uint32(v309)+4))
										v365 = v364
										if v365 < int32(1) {
										} else {
											v369 = *(*int32)(unsafe.Add(mBase, uint32(v309)))
											v372 = v369
											for {
												v379 = *(*int32)(unsafe.Add(mBase, uint32(v309)+12))
												*(*int32)(unsafe.Add(mBase, uint32(v309)+12)) = v379 + int32(1)
												*(*uint8)(unsafe.Add(mBase, uint32(v379))) = uint8(v372)
												v384 = *(*int32)(unsafe.Add(mBase, uint32(v309)))
												v385 = int32(8)
												v386 = int32(base.Ui32(v384) >> (uint(v385) % 32))
												*(*int32)(unsafe.Add(mBase, uint32(v309))) = v386
												v388 = *(*int32)(unsafe.Add(mBase, uint32(v309)+4))
												*(*int32)(unsafe.Add(mBase, uint32(v309)+4)) = v388 + int32(-8)
												if v385 < v388 {
													v372 = v386
													continue
												} else {
													break
												}
												break
											}
										}
										*(*int32)(unsafe.Add(mBase, uint32(v309)+4)) = int32(0)
										v405 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
										v415 = v405
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(v309)+20)) = int32(1)
										v353 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
										v415 = v353
									}
								} else {
									if base.Ui32(v334) <= base.Ui32(v336) {
										v365 = v322
										if v365 < int32(1) {
										} else {
											v369 = *(*int32)(unsafe.Add(mBase, uint32(v309)))
											v372 = v369
											for {
												v379 = *(*int32)(unsafe.Add(mBase, uint32(v309)+12))
												*(*int32)(unsafe.Add(mBase, uint32(v309)+12)) = v379 + int32(1)
												*(*uint8)(unsafe.Add(mBase, uint32(v379))) = uint8(v372)
												v384 = *(*int32)(unsafe.Add(mBase, uint32(v309)))
												v385 = int32(8)
												v386 = int32(base.Ui32(v384) >> (uint(v385) % 32))
												*(*int32)(unsafe.Add(mBase, uint32(v309))) = v386
												v388 = *(*int32)(unsafe.Add(mBase, uint32(v309)+4))
												*(*int32)(unsafe.Add(mBase, uint32(v309)+4)) = v388 + int32(-8)
												if v385 < v388 {
													v372 = v386
													continue
												} else {
													break
												}
												break
											}
										}
										*(*int32)(unsafe.Add(mBase, uint32(v309)+4)) = int32(0)
										v405 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
										v415 = v405
									} else {
										v343 = int32(base.Ui32(v336*int32(3)) >> (uint(int32(1)) % 32))
										if base.Ui32(v334) < base.Ui32(v343) {
											v345 = v343
										} else {
											v345 = v334
										}
										v349 = v345&int32(-1024) + int32(1024)
										v350 = F_WebPSafeMalloc(m, int64(1), v349)
										mBase = m.M
										if v350 != 0 {
											if v318 == v319 {
											} else {
												v355 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
												v356 = F_memcpy(m, v350, v355, v320)
												mBase = m.M
											}
											v357 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
											F_WebPSafeFree(m, v357)
											mBase = m.M
											*(*int32)(unsafe.Add(mBase, uint32(v309)+16)) = v350 + v349
											*(*int32)(unsafe.Add(mBase, uint32(v309)+12)) = v350 + v320
											*(*int32)(unsafe.Add(mBase, uint32(v309)+8)) = v350
											v364 = *(*int32)(unsafe.Add(mBase, uint32(v309)+4))
											v365 = v364
											if v365 < int32(1) {
											} else {
												v369 = *(*int32)(unsafe.Add(mBase, uint32(v309)))
												v372 = v369
												for {
													v379 = *(*int32)(unsafe.Add(mBase, uint32(v309)+12))
													*(*int32)(unsafe.Add(mBase, uint32(v309)+12)) = v379 + int32(1)
													*(*uint8)(unsafe.Add(mBase, uint32(v379))) = uint8(v372)
													v384 = *(*int32)(unsafe.Add(mBase, uint32(v309)))
													v385 = int32(8)
													v386 = int32(base.Ui32(v384) >> (uint(v385) % 32))
													*(*int32)(unsafe.Add(mBase, uint32(v309))) = v386
													v388 = *(*int32)(unsafe.Add(mBase, uint32(v309)+4))
													*(*int32)(unsafe.Add(mBase, uint32(v309)+4)) = v388 + int32(-8)
													if v385 < v388 {
														v372 = v386
														continue
													} else {
														break
													}
													break
												}
											}
											*(*int32)(unsafe.Add(mBase, uint32(v309)+4)) = int32(0)
											v405 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
											v415 = v405
										} else {
											*(*int32)(unsafe.Add(mBase, uint32(v309)+20)) = int32(1)
											v353 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
											v415 = v353
										}
									}
								}
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v309)+20)) = int32(1)
								v333 = *(*int32)(unsafe.Add(mBase, uint32(v309)+8))
								v415 = v333
							}
							v416 = *(*int32)(unsafe.Add(mBase, uint32(v15)+28))
							if v416 == int32(0) {
								v452 = *(*int32)(unsafe.Add(mBase, uint32(v15)+12))
								v457 = *(*int32)(unsafe.Add(mBase, uint32(v15)+20))
								v458 = *(*int32)(unsafe.Add(mBase, uint32(v15)+16))
								v460 = (v452+int32(7))>>(uint(int32(3))%32) + (v457 - v458)
								if base.Ui32(v24) < base.Ui32(v460) {
									v465 = v15 + int32(8)
									if v465 == int32(0) {
									} else {
										v470 = v15 + int32(16)
										v471 = *(*int32)(unsafe.Add(mBase, uint32(v470)))
										F_WebPSafeFree(m, v471)
										mBase = m.M
										v475 = int64(0)
										*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v475
										*(*int64)(unsafe.Add(mBase, uint32(v470))) = v475
										*(*int64)(unsafe.Add(mBase, uint32(v465))) = v475
									}
									v485 = int32(0)
									v486 = int32(1)
									v488 = v24
									v489 = v485
									v490 = v486
									v491 = v23
								} else {
									v488 = v460
									v489 = int32(1)
									v490 = int32(0)
									v491 = v415
								}
								v495 = v489 | l4<<(uint(int32(2))%32)
								if l5 != 0 {
									v498 = v495 | int32(16)
								} else {
									v498 = v495
								}
								*(*uint8)(unsafe.Add(mBase, uint32(v15)+32)) = uint8(v498)
								v501 = l8 + int32(4)
								v503 = v488 + int32(1)
								v506 = int64(0)
								*(*int64)(unsafe.Add(mBase, uint32(v501)+16)) = v506
								*(*int64)(unsafe.Add(mBase, uint32(v501)+8)) = int64(-34359738368)
								*(*int64)(unsafe.Add(mBase, uint32(v501))) = int64(254)
								*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v506
								if v503 == int32(0) {
									v545 = int32(1)
								} else {
									v519 = int32(1024)
									if base.Ui32(v519) < base.Ui32(v503) {
										v522 = v503
									} else {
										v522 = v519
									}
									v523 = F_WebPSafeMalloc(m, int64(1), v522)
									mBase = m.M
									if v523 != 0 {
										v527 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
										if v527 == int32(0) {
										} else {
											v532 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(20))))
											v533 = F_memcpy(m, v523, v532, v527)
											mBase = m.M
										}
										v534 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
										F_WebPSafeFree(m, v534)
										mBase = m.M
										*(*int32)(unsafe.Add(mBase, uint32(v501)+24)) = v522
										*(*int32)(unsafe.Add(mBase, uint32(v501)+16)) = v523
										v545 = int32(1)
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
										v545 = int32(0)
									}
								}
								if v545 == int32(0) {
									v679 = v25
								} else {
									v550 = int32(1)
									v556 = *(*int32)(unsafe.Add(mBase, uint32(v501)+12))
									if v556 != int32(-8) {
										v604 = int32(0)
										v612 = v604
									} else {
										v559 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
										v562 = base.I64_extend_i32_u(v559) + base.I64_extend_i32_u(v550)
										if base.Ui64(v562) < base.Ui64(int64(4294967296)) {
											v568 = *(*int32)(unsafe.Add(mBase, uint32(v501)+24))
											v569 = base.I32_wrap_i64(v562)
											if base.Ui32(v568) < base.Ui32(v569) {
												v574 = v568 << (uint(int32(1)) % 32)
												if base.Ui32(v569) < base.Ui32(v574) {
													v576 = v574
												} else {
													v576 = v569
												}
												v577 = int32(1024)
												if base.Ui32(v577) < base.Ui32(v576) {
													v580 = v576
												} else {
													v580 = v577
												}
												v581 = F_WebPSafeMalloc(m, int64(1), v580)
												mBase = m.M
												if v581 != 0 {
													v585 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
													if v585 == int32(0) {
													} else {
														v588 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
														v589 = F_memcpy(m, v581, v588, v585)
														mBase = m.M
													}
													v590 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
													F_WebPSafeFree(m, v590)
													mBase = m.M
													*(*int32)(unsafe.Add(mBase, uint32(v501)+24)) = v580
													*(*int32)(unsafe.Add(mBase, uint32(v501)+16)) = v581
													v594 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
													v595 = v594
													v597 = v581
													v599 = F_memcpy(m, v597+v595, v15+int32(32), v550)
													mBase = m.M
													v600 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
													*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v600 + v550
													v604 = int32(1)
													v612 = v604
												} else {
													*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
													v612 = int32(0)
												}
											} else {
												v571 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
												v595 = v559
												v597 = v571
												v599 = F_memcpy(m, v597+v595, v15+int32(32), v550)
												mBase = m.M
												v600 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
												*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v600 + v550
												v604 = int32(1)
												v612 = v604
											}
										} else {
											*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
											v612 = int32(0)
										}
									}
									if v612 == int32(0) {
										v679 = v25
									} else {
										v620 = *(*int32)(unsafe.Add(mBase, uint32(v501)+12))
										if v620 != int32(-8) {
											v668 = int32(0)
											v676 = v668
										} else {
											v623 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
											v626 = base.I64_extend_i32_u(v623) + base.I64_extend_i32_u(v488)
											if base.Ui64(v626) < base.Ui64(int64(4294967296)) {
												v632 = *(*int32)(unsafe.Add(mBase, uint32(v501)+24))
												v633 = base.I32_wrap_i64(v626)
												if base.Ui32(v632) < base.Ui32(v633) {
													v638 = v632 << (uint(int32(1)) % 32)
													if base.Ui32(v633) < base.Ui32(v638) {
														v640 = v638
													} else {
														v640 = v633
													}
													v641 = int32(1024)
													if base.Ui32(v641) < base.Ui32(v640) {
														v644 = v640
													} else {
														v644 = v641
													}
													v645 = F_WebPSafeMalloc(m, int64(1), v644)
													mBase = m.M
													if v645 != 0 {
														v649 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
														if v649 == int32(0) {
														} else {
															v652 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
															v653 = F_memcpy(m, v645, v652, v649)
															mBase = m.M
														}
														v654 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
														F_WebPSafeFree(m, v654)
														mBase = m.M
														*(*int32)(unsafe.Add(mBase, uint32(v501)+24)) = v644
														*(*int32)(unsafe.Add(mBase, uint32(v501)+16)) = v645
														v658 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
														v659 = v658
														v661 = v645
														v663 = F_memcpy(m, v661+v659, v491, v488)
														mBase = m.M
														v664 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
														*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v664 + v488
														v668 = int32(1)
														v676 = v668
													} else {
														*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
														v676 = int32(0)
													}
												} else {
													v635 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
													v659 = v623
													v661 = v635
													v663 = F_memcpy(m, v661+v659, v491, v488)
													mBase = m.M
													v664 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
													*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v664 + v488
													v668 = int32(1)
													v676 = v668
												}
											} else {
												*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
												v676 = int32(0)
											}
										}
										v679 = base.B2i32(v676 != int32(0))
									}
								}
								if v490 != 0 {
								} else {
									v681 = v15 + int32(8)
									if v681 == int32(0) {
									} else {
										v686 = v15 + int32(16)
										v687 = *(*int32)(unsafe.Add(mBase, uint32(v686)))
										F_WebPSafeFree(m, v687)
										mBase = m.M
										v691 = int64(0)
										*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v691
										*(*int64)(unsafe.Add(mBase, uint32(v686))) = v691
										*(*int64)(unsafe.Add(mBase, uint32(v681))) = v691
									}
								}
								v698 = int32(0)
								if v679 == v698 {
									v704 = v698
								} else {
									v701 = *(*int32)(unsafe.Add(mBase, uint32(l8)+32))
									v704 = base.B2i32(v701 == int32(0))
								}
								v707 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(24))))
								*(*int32)(unsafe.Add(mBase, uint32(l8))) = v707
								v709 = v704
							} else {
								v420 = v15 + int32(8)
								if v420 == int32(0) {
								} else {
									v425 = v15 + int32(16)
									v426 = *(*int32)(unsafe.Add(mBase, uint32(v425)))
									F_WebPSafeFree(m, v426)
									mBase = m.M
									v430 = int64(0)
									*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v430
									*(*int64)(unsafe.Add(mBase, uint32(v425))) = v430
									*(*int64)(unsafe.Add(mBase, uint32(v420))) = v430
								}
								v437 = int64(0)
								*(*int64)(unsafe.Add(mBase, uint32(l8)+4)) = v437
								*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v437
								*(*int64)(unsafe.Add(mBase, uint32(l8+int32(20)))) = v437
								*(*int64)(unsafe.Add(mBase, uint32(l8+int32(12)))) = v437
								v709 = int32(0)
							}
						} else {
							v255 = v15 + int32(8)
							if v255 == int32(0) {
							} else {
								v260 = v15 + int32(16)
								v261 = *(*int32)(unsafe.Add(mBase, uint32(v260)))
								F_WebPSafeFree(m, v261)
								mBase = m.M
								v265 = int64(0)
								*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v265
								*(*int64)(unsafe.Add(mBase, uint32(v260))) = v265
								*(*int64)(unsafe.Add(mBase, uint32(v255))) = v265
							}
							v276 = v15 + int32(8)
							if v276 == int32(0) {
							} else {
								v281 = v15 + int32(16)
								v282 = *(*int32)(unsafe.Add(mBase, uint32(v281)))
								F_WebPSafeFree(m, v282)
								mBase = m.M
								v286 = int64(0)
								*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v286
								*(*int64)(unsafe.Add(mBase, uint32(v281))) = v286
								*(*int64)(unsafe.Add(mBase, uint32(v276))) = v286
							}
							v293 = int64(0)
							*(*int64)(unsafe.Add(mBase, uint32(l8)+4)) = v293
							*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v293
							*(*int64)(unsafe.Add(mBase, uint32(l8+int32(20)))) = v293
							*(*int64)(unsafe.Add(mBase, uint32(l8+int32(12)))) = v293
							v709 = int32(0)
						}
					}
				}
			}
		}
	} else {
		v485 = int32(0)
		v486 = int32(1)
		v488 = v24
		v489 = v485
		v490 = v486
		v491 = v23
		v495 = v489 | l4<<(uint(int32(2))%32)
		if l5 != 0 {
			v498 = v495 | int32(16)
		} else {
			v498 = v495
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v15)+32)) = uint8(v498)
		v501 = l8 + int32(4)
		v503 = v488 + int32(1)
		v506 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(v501)+16)) = v506
		*(*int64)(unsafe.Add(mBase, uint32(v501)+8)) = int64(-34359738368)
		*(*int64)(unsafe.Add(mBase, uint32(v501))) = int64(254)
		*(*int64)(unsafe.Add(mBase, uint32(l8+int32(28)))) = v506
		if v503 == int32(0) {
			v545 = int32(1)
		} else {
			v519 = int32(1024)
			if base.Ui32(v519) < base.Ui32(v503) {
				v522 = v503
			} else {
				v522 = v519
			}
			v523 = F_WebPSafeMalloc(m, int64(1), v522)
			mBase = m.M
			if v523 != 0 {
				v527 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
				if v527 == int32(0) {
				} else {
					v532 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(20))))
					v533 = F_memcpy(m, v523, v532, v527)
					mBase = m.M
				}
				v534 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
				F_WebPSafeFree(m, v534)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(v501)+24)) = v522
				*(*int32)(unsafe.Add(mBase, uint32(v501)+16)) = v523
				v545 = int32(1)
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
				v545 = int32(0)
			}
		}
		if v545 == int32(0) {
			v679 = v25
		} else {
			v550 = int32(1)
			v556 = *(*int32)(unsafe.Add(mBase, uint32(v501)+12))
			if v556 != int32(-8) {
				v604 = int32(0)
				v612 = v604
			} else {
				v559 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
				v562 = base.I64_extend_i32_u(v559) + base.I64_extend_i32_u(v550)
				if base.Ui64(v562) < base.Ui64(int64(4294967296)) {
					v568 = *(*int32)(unsafe.Add(mBase, uint32(v501)+24))
					v569 = base.I32_wrap_i64(v562)
					if base.Ui32(v568) < base.Ui32(v569) {
						v574 = v568 << (uint(int32(1)) % 32)
						if base.Ui32(v569) < base.Ui32(v574) {
							v576 = v574
						} else {
							v576 = v569
						}
						v577 = int32(1024)
						if base.Ui32(v577) < base.Ui32(v576) {
							v580 = v576
						} else {
							v580 = v577
						}
						v581 = F_WebPSafeMalloc(m, int64(1), v580)
						mBase = m.M
						if v581 != 0 {
							v585 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
							if v585 == int32(0) {
							} else {
								v588 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
								v589 = F_memcpy(m, v581, v588, v585)
								mBase = m.M
							}
							v590 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
							F_WebPSafeFree(m, v590)
							mBase = m.M
							*(*int32)(unsafe.Add(mBase, uint32(v501)+24)) = v580
							*(*int32)(unsafe.Add(mBase, uint32(v501)+16)) = v581
							v594 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
							v595 = v594
							v597 = v581
							v599 = F_memcpy(m, v597+v595, v15+int32(32), v550)
							mBase = m.M
							v600 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
							*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v600 + v550
							v604 = int32(1)
							v612 = v604
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
							v612 = int32(0)
						}
					} else {
						v571 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
						v595 = v559
						v597 = v571
						v599 = F_memcpy(m, v597+v595, v15+int32(32), v550)
						mBase = m.M
						v600 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
						*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v600 + v550
						v604 = int32(1)
						v612 = v604
					}
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
					v612 = int32(0)
				}
			}
			if v612 == int32(0) {
				v679 = v25
			} else {
				v620 = *(*int32)(unsafe.Add(mBase, uint32(v501)+12))
				if v620 != int32(-8) {
					v668 = int32(0)
					v676 = v668
				} else {
					v623 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
					v626 = base.I64_extend_i32_u(v623) + base.I64_extend_i32_u(v488)
					if base.Ui64(v626) < base.Ui64(int64(4294967296)) {
						v632 = *(*int32)(unsafe.Add(mBase, uint32(v501)+24))
						v633 = base.I32_wrap_i64(v626)
						if base.Ui32(v632) < base.Ui32(v633) {
							v638 = v632 << (uint(int32(1)) % 32)
							if base.Ui32(v633) < base.Ui32(v638) {
								v640 = v638
							} else {
								v640 = v633
							}
							v641 = int32(1024)
							if base.Ui32(v641) < base.Ui32(v640) {
								v644 = v640
							} else {
								v644 = v641
							}
							v645 = F_WebPSafeMalloc(m, int64(1), v644)
							mBase = m.M
							if v645 != 0 {
								v649 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
								if v649 == int32(0) {
								} else {
									v652 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
									v653 = F_memcpy(m, v645, v652, v649)
									mBase = m.M
								}
								v654 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
								F_WebPSafeFree(m, v654)
								mBase = m.M
								*(*int32)(unsafe.Add(mBase, uint32(v501)+24)) = v644
								*(*int32)(unsafe.Add(mBase, uint32(v501)+16)) = v645
								v658 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
								v659 = v658
								v661 = v645
								v663 = F_memcpy(m, v661+v659, v491, v488)
								mBase = m.M
								v664 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
								*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v664 + v488
								v668 = int32(1)
								v676 = v668
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
								v676 = int32(0)
							}
						} else {
							v635 = *(*int32)(unsafe.Add(mBase, uint32(v501)+16))
							v659 = v623
							v661 = v635
							v663 = F_memcpy(m, v661+v659, v491, v488)
							mBase = m.M
							v664 = *(*int32)(unsafe.Add(mBase, uint32(v501)+20))
							*(*int32)(unsafe.Add(mBase, uint32(v501)+20)) = v664 + v488
							v668 = int32(1)
							v676 = v668
						}
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v501)+28)) = int32(1)
						v676 = int32(0)
					}
				}
				v679 = base.B2i32(v676 != int32(0))
			}
		}
		if v490 != 0 {
		} else {
			v681 = v15 + int32(8)
			if v681 == int32(0) {
			} else {
				v686 = v15 + int32(16)
				v687 = *(*int32)(unsafe.Add(mBase, uint32(v686)))
				F_WebPSafeFree(m, v687)
				mBase = m.M
				v691 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(v15+int32(24)))) = v691
				*(*int64)(unsafe.Add(mBase, uint32(v686))) = v691
				*(*int64)(unsafe.Add(mBase, uint32(v681))) = v691
			}
		}
		v698 = int32(0)
		if v679 == v698 {
			v704 = v698
		} else {
			v701 = *(*int32)(unsafe.Add(mBase, uint32(l8)+32))
			v704 = base.B2i32(v701 == int32(0))
		}
		v707 = *(*int32)(unsafe.Add(mBase, uint32(l8+int32(24))))
		*(*int32)(unsafe.Add(mBase, uint32(l8))) = v707
		v709 = v704
	}
	m.G0 = v15 + int32(320)
	return v709
}
func F_EncodeImageNoHuffman(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v19 int32
	_ = v19
	var v28 int64
	_ = v28
	var v42 int32
	_ = v42
	var v56 int64
	_ = v56
	var v64 int64
	_ = v64
	var v65 int32
	_ = v65
	var v66 int64
	_ = v66
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v91 int32
	_ = v91
	var v99 int32
	_ = v99
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v109 int32
	_ = v109
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
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v133 int32
	_ = v133
	var v135 int32
	_ = v135
	var v138 int32
	_ = v138
	var v143 int32
	_ = v143
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v171 int32
	_ = v171
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v209 int32
	_ = v209
	var v211 int32
	_ = v211
	var v215 int32
	_ = v215
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v224 int32
	_ = v224
	var v229 int32
	_ = v229
	var v234 int32
	_ = v234
	var v241 int32
	_ = v241
	var v247 int32
	_ = v247
	var v251 int32
	_ = v251
	var v253 int32
	_ = v253
	var v255 int32
	_ = v255
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v268 int32
	_ = v268
	var v276 int32
	_ = v276
	var v285 int32
	_ = v285
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v300 int32
	_ = v300
	var v311 int32
	_ = v311
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v317 int32
	_ = v317
	var v319 int32
	_ = v319
	var v324 int32
	_ = v324
	var v327 int32
	_ = v327
	var v333 int32
	_ = v333
	var v335 int32
	_ = v335
	var v337 int32
	_ = v337
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v354 int32
	_ = v354
	var v356 int32
	_ = v356
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v372 int32
	_ = v372
	var v373 int32
	_ = v373
	var v377 int32
	_ = v377
	var v378 int32
	_ = v378
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v396 int32
	_ = v396
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v399 int32
	_ = v399
	var v408 int32
	_ = v408
	var v409 int32
	_ = v409
	var v412 int64
	_ = v412
	var v415 int32
	_ = v415
	var v417 int64
	_ = v417
	var v420 int32
	_ = v420
	var v427 int32
	_ = v427
	var v429 int32
	_ = v429
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v436 int32
	_ = v436
	var v445 int32
	_ = v445
	var v446 int32
	_ = v446
	var v447 int32
	_ = v447
	var v449 int32
	_ = v449
	var v454 int32
	_ = v454
	var v455 int32
	_ = v455
	var v462 int32
	_ = v462
	var v465 int32
	_ = v465
	var v467 int32
	_ = v467
	var v471 int32
	_ = v471
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v511 int32
	_ = v511
	var v512 int32
	_ = v512
	var v514 int32
	_ = v514
	var v515 int32
	_ = v515
	var v517 int32
	_ = v517
	var v518 int32
	_ = v518
	var v520 int32
	_ = v520
	var v521 int32
	_ = v521
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v527 int32
	_ = v527
	var v528 int64
	_ = v528
	var v529 int32
	_ = v529
	var v536 int64
	_ = v536
	var v537 int32
	_ = v537
	var v538 int64
	_ = v538
	var v548 int32
	_ = v548
	var v550 int32
	_ = v550
	var v555 int32
	_ = v555
	var v561 int32
	_ = v561
	var v564 int32
	_ = v564
	var v572 int32
	_ = v572
	var v575 int32
	_ = v575
	var v576 int32
	_ = v576
	var v582 int32
	_ = v582
	var v589 int32
	_ = v589
	var v593 int32
	_ = v593
	var v594 int32
	_ = v594
	var v605 int32
	_ = v605
	var v606 int32
	_ = v606
	var v612 int32
	_ = v612
	var v614 int32
	_ = v614
	var v616 int32
	_ = v616
	var v623 int32
	_ = v623
	var v624 int32
	_ = v624
	var v645 int32
	_ = v645
	var v648 int32
	_ = v648
	var v656 int32
	_ = v656
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v666 int32
	_ = v666
	var v673 int32
	_ = v673
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v689 int32
	_ = v689
	var v690 int32
	_ = v690
	var v696 int32
	_ = v696
	var v698 int32
	_ = v698
	var v700 int32
	_ = v700
	var v707 int32
	_ = v707
	var v708 int32
	_ = v708
	var v729 int32
	_ = v729
	var v732 int32
	_ = v732
	var v740 int32
	_ = v740
	var v743 int32
	_ = v743
	var v744 int32
	_ = v744
	var v750 int32
	_ = v750
	var v757 int32
	_ = v757
	var v761 int32
	_ = v761
	var v762 int32
	_ = v762
	var v773 int32
	_ = v773
	var v774 int32
	_ = v774
	var v780 int32
	_ = v780
	var v782 int32
	_ = v782
	var v784 int32
	_ = v784
	var v791 int32
	_ = v791
	var v792 int32
	_ = v792
	var v813 int32
	_ = v813
	var v816 int32
	_ = v816
	var v824 int32
	_ = v824
	var v827 int32
	_ = v827
	var v828 int32
	_ = v828
	var v834 int32
	_ = v834
	var v841 int32
	_ = v841
	var v845 int32
	_ = v845
	var v846 int32
	_ = v846
	var v857 int32
	_ = v857
	var v858 int32
	_ = v858
	var v864 int32
	_ = v864
	var v866 int32
	_ = v866
	var v868 int32
	_ = v868
	var v875 int32
	_ = v875
	var v876 int32
	_ = v876
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
	var v916 int32
	_ = v916
	var v923 int32
	_ = v923
	var v927 int32
	_ = v927
	var v928 int32
	_ = v928
	var v939 int32
	_ = v939
	var v940 int32
	_ = v940
	var v946 int32
	_ = v946
	var v948 int32
	_ = v948
	var v950 int32
	_ = v950
	var v957 int32
	_ = v957
	var v958 int32
	_ = v958
	var v981 int32
	_ = v981
	var v983 int32
	_ = v983
	var v999 int32
	_ = v999
	var v1001 int32
	_ = v1001
	var v1006 int32
	_ = v1006
	var v1008 int32
	_ = v1008
	v12 = int32(0)
	v17 = m.G0
	v19 = v17 - int32(80)
	m.G0 = v19
	*(*int32)(unsafe.Add(mBase, uint32(v19+int32(72)))) = v12
	v28 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v19+int32(64)))) = v28
	*(*int64)(unsafe.Add(mBase, uint32(v19+int32(56)))) = v28
	*(*int64)(unsafe.Add(mBase, uint32(v19+int32(48)))) = v28
	*(*int64)(unsafe.Add(mBase, uint32(v19+int32(40)))) = v28
	v42 = int32(16)
	*(*int64)(unsafe.Add(mBase, uint32(v19+int32(32)))) = v28
	*(*int64)(unsafe.Add(mBase, uint32(v19)+24)) = v28
	*(*int64)(unsafe.Add(mBase, uint32(v19)+16)) = v28
	*(*int32)(unsafe.Add(mBase, uint32(v19)+12)) = v12
	*(*int32)(unsafe.Add(mBase, uint32(v19)+8)) = v12
	v56 = int64(57)
	goto L6
L1:
	;
	F_free(m, v1001)
	mBase = m.M
	goto L167
L2:
	;
	v85 = int32(0)
	v87 = base.I32_div_s(l9, int32(2))
	v88 = F_VP8LHashChainFill(m, l2, l6, l1, l4, l5, l7, l8, v87, l10)
	mBase = m.M
	if v88 == v85 {
		v999 = v12
		v1001 = v85
		goto L1
	} else {
		goto L13
	}
L3:
	;
	if v78 != 0 {
		goto L2
	} else {
		goto L9
	}
L4:
	;
	goto L3
L5:
	;
	v76 = F_malloc(m, base.I32_wrap_i64(v56)*v42)
	mBase = m.M
	v78 = v76
	goto L4
L6:
	;
	v64 = base.I64_div_u_s(int64(2147418112), v56)
	v65 = int32(0)
	v66 = base.I64_extend_i32_u(v42)
	if base.Ui64(int64(4294967295)) < base.Ui64(v66*v56) {
		v78 = v65
		goto L4
	} else {
		goto L7
	}
L7:
	;
	if base.Ui64(v64) < base.Ui64(v66) {
		v78 = v65
		goto L4
	} else {
		goto L8
	}
L8:
	;
	goto L5
L9:
	;
	v81 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v81 != 0 {
		goto L11
	} else {
		goto L12
	}
L10:
	;
	v999 = v12
	v1001 = int32(0)
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
	v91 = int32(0)
	v99 = F_VP8LGetBackwardReferences(m, l4, l5, l1, l6, v91, int32(3), v91, v91, l2, l3, v19+int32(8), l8, l9-v87, l10)
	mBase = m.M
	if v99 != 0 {
		goto L14
	} else {
		goto L15
	}
L14:
	;
	v102 = *(*int32)(unsafe.Add(mBase, uint32(v19)+8))
	v103 = F_VP8LAllocateHistogramSet(m, int32(1), v102)
	mBase = m.M
	if v103 != 0 {
		goto L16
	} else {
		goto L17
	}
L15:
	;
	v999 = int32(0)
	v1001 = v91
	goto L1
L16:
	;
	v109 = int32(0)
	v121 = *(*int32)(unsafe.Add(mBase, uint32(v103)+8))
	v122 = *(*int32)(unsafe.Add(mBase, uint32(v121)))
	v123 = *(*int32)(unsafe.Add(mBase, uint32(v122)+3236))
	v124 = int32(4) << (uint(v123) % 32)
	v125 = int32(_a_F_EncodeImageNoHuffman_0)
	v129 = base.B2i32(v109 < v123)
	if v109 < v123 {
		goto L22
	} else {
		goto L23
	}
L17:
	;
	v105 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v105 != 0 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	v999 = int32(0)
	v1001 = v91
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
	v311 = int32(0)
	v314 = *(*int32)(unsafe.Add(mBase, uint32(v103)+8))
	v315 = *(*int32)(unsafe.Add(mBase, uint32(v314)))
	v317 = m.G0
	v319 = v317 - int32(16)
	m.G0 = v319
	F_VP8LRefsCursorInit(m, v319+int32(4), l3)
	mBase = m.M
	v324 = *(*int32)(unsafe.Add(mBase, uint32(v319)+4))
	if v324 == v311 {
		goto L48
	} else {
		goto L49
	}
L22:
	;
	v130 = v124 + v125
	goto L24
L23:
	;
	v130 = v125
	goto L24
L24:
	;
	v131 = *(*int32)(unsafe.Add(mBase, uint32(v103)+4))
	v133 = int32(12)
	v135 = F_memset(m, v103, v109, v130*v131+v133)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v135)+4)) = v131
	v138 = v135 + v133
	*(*int32)(unsafe.Add(mBase, uint32(v135)+8)) = v138
	*(*int32)(unsafe.Add(mBase, uint32(v135))) = v131
	if v131 < int32(1) {
		goto L25
	} else {
		goto L26
	}
L25:
	;
	goto L21
L26:
	;
	v143 = int32(1)
	v147 = v138 + v131<<(uint(int32(2))%32)
	if v131 == v143 {
		v203 = v109
		v209 = v138
		v211 = v147
		goto L27
	} else {
		goto L28
	}
L27:
	;
	if v131&v143 == int32(0) {
		v229 = v209
		goto L35
	} else {
		goto L36
	}
L28:
	;
	v150 = int32(_a_F_EncodeImageNoHuffman_1)
	if v109 < v123 {
		goto L29
	} else {
		goto L30
	}
L29:
	;
	v153 = v124 + v150
	goto L31
L30:
	;
	v153 = v150
	goto L31
L31:
	;
	v158 = int32(4)
	v159 = int32(0)
	v165 = v138
	v167 = v147
	goto L32
L32:
	;
	v169 = int32(-4)
	v171 = int32(31)
	v173 = int32(-32)
	v174 = (v167 + v171) & v173
	*(*int32)(unsafe.Add(mBase, uint32(v165+v158+v169))) = v174
	v176 = *(*int32)(unsafe.Add(mBase, uint32(v135)+8))
	v177 = v176 + v158
	v180 = *(*int32)(unsafe.Add(mBase, uint32(v177+v169)))
	v181 = int32(3312)
	*(*int32)(unsafe.Add(mBase, uint32(v180))) = v174 + v181
	v188 = (v174 + v153 + v171) & v173
	*(*int32)(unsafe.Add(mBase, uint32(v177))) = v188
	v190 = *(*int32)(unsafe.Add(mBase, uint32(v135)+8))
	v192 = *(*int32)(unsafe.Add(mBase, uint32(v190+v158)))
	*(*int32)(unsafe.Add(mBase, uint32(v192))) = v188 + v181
	v198 = v188 + v153
	v200 = v159 + int32(2)
	if v131&int32(2147483646) != v200 {
		v158 = v158 + int32(8)
		v159 = v200
		v165 = v190
		v167 = v198
		goto L32
	} else {
		goto L34
	}
L33:
	;
	v203 = v200
	v209 = v190
	v211 = v198
	goto L27
L34:
	;
	goto L33
L35:
	;
	if v131 < int32(1) {
		goto L25
	} else {
		goto L37
	}
L36:
	;
	v215 = v203 << (uint(int32(2)) % 32)
	v220 = (v211 + int32(31)) & int32(-32)
	*(*int32)(unsafe.Add(mBase, uint32(v209+v215))) = v220
	v222 = *(*int32)(unsafe.Add(mBase, uint32(v135)+8))
	v224 = *(*int32)(unsafe.Add(mBase, uint32(v222+v215)))
	*(*int32)(unsafe.Add(mBase, uint32(v224))) = v220 + int32(3312)
	v229 = v222
	goto L35
L37:
	;
	v234 = v131 & int32(3)
	if base.Ui32(v131) < base.Ui32(int32(4)) {
		v276 = int32(0)
		goto L38
	} else {
		goto L39
	}
L38:
	;
	if v234 == int32(0) {
		goto L25
	} else {
		goto L43
	}
L39:
	;
	v241 = v229
	v247 = int32(0)
	goto L40
L40:
	;
	v251 = *(*int32)(unsafe.Add(mBase, uint32(v241)))
	*(*int32)(unsafe.Add(mBase, uint32(v251)+3236)) = v123
	v253 = int32(4)
	v255 = *(*int32)(unsafe.Add(mBase, uint32(v241+v253)))
	*(*int32)(unsafe.Add(mBase, uint32(v255)+3236)) = v123
	v259 = *(*int32)(unsafe.Add(mBase, uint32(v241+int32(8))))
	*(*int32)(unsafe.Add(mBase, uint32(v259)+3236)) = v123
	v263 = *(*int32)(unsafe.Add(mBase, uint32(v241+int32(12))))
	*(*int32)(unsafe.Add(mBase, uint32(v263)+3236)) = v123
	v268 = v247 + v253
	if v131&int32(2147483644) != v268 {
		v241 = v241 + int32(16)
		v247 = v268
		goto L40
	} else {
		goto L42
	}
L41:
	;
	v276 = v268
	goto L38
L42:
	;
	goto L41
L43:
	;
	v285 = v229 + v276<<(uint(int32(2))%32)
	v294 = v234
	goto L44
L44:
	;
	v295 = *(*int32)(unsafe.Add(mBase, uint32(v285)))
	*(*int32)(unsafe.Add(mBase, uint32(v295)+3236)) = v123
	v300 = v294 + int32(-1)
	if v300 != 0 {
		v285 = v285 + int32(4)
		v294 = v300
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
	v354 = F_GetHuffBitLengthsAndCodes(m, v103, v19+int32(16))
	mBase = m.M
	if v354 != 0 {
		goto L55
	} else {
		goto L56
	}
L48:
	;
	m.G0 = v319 + int32(16)
	goto L47
L49:
	;
	v327 = v324
	goto L50
L50:
	;
	F_HistogramAddSinglePixOrCopy(m, v315, v327, v311, v311)
	mBase = m.M
	v333 = *(*int32)(unsafe.Add(mBase, uint32(v319)+4))
	v335 = v333 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v319)+4)) = v335
	v337 = *(*int32)(unsafe.Add(mBase, uint32(v319)+12))
	if v335 != v337 {
		v343 = v335
		goto L52
	} else {
		goto L53
	}
L51:
	;
	goto L48
L52:
	;
	if v343 != 0 {
		v327 = v343
		goto L50
	} else {
		goto L54
	}
L53:
	;
	F_VP8LRefsCursorNextBlock(m, v319+int32(4))
	mBase = m.M
	v342 = *(*int32)(unsafe.Add(mBase, uint32(v319)+4))
	v343 = v342
	goto L52
L54:
	;
	goto L51
L55:
	;
	v359 = int32(0)
	v360 = int32(1)
	goto L62
L56:
	;
	v356 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v356 != 0 {
		goto L58
	} else {
		goto L59
	}
L57:
	;
	v999 = v103
	v1001 = v311
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
	v511 = *(*int32)(unsafe.Add(mBase, uint32(v19)+16))
	v512 = *(*int32)(unsafe.Add(mBase, uint32(v19)+28))
	if v512 < v511 {
		goto L87
	} else {
		goto L88
	}
L61:
	;
	goto L60
L62:
	;
	v372 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v373 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v373+v360 < int32(32) {
		goto L65
	} else {
		goto L66
	}
L63:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v474 + v472
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v471<<(uint(v474)%32) | v473
	goto L61
L64:
	;
	v391 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v392 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v396 = v388
	v397 = v389
	v398 = v392
	v399 = v391
	goto L68
L65:
	;
	if v373 < int32(16) {
		v471 = v359
		v472 = v360
		v473 = v372
		v474 = v373
		goto L63
	} else {
		goto L67
	}
L66:
	;
	v377 = int32(32)
	v378 = v377 - v373
	v386 = int32(base.Ui32(v359) >> (uint(v378) % 32))
	v387 = v360 - v378
	v388 = v359<<(uint(v373)%32) | v372
	v389 = v377
	goto L64
L67:
	;
	v386 = v359
	v387 = v360
	v388 = v372
	v389 = v373
	goto L64
L68:
	;
	if base.Ui32(v398+int32(2)) <= base.Ui32(v399) {
		v454 = v398
		v455 = v399
		goto L70
	} else {
		goto L71
	}
L69:
	;
	v471 = v386
	v472 = v387
	v473 = v467
	v474 = v465
	goto L63
L70:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v454))) = uint16(v396)
	v462 = v454 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v462
	v465 = v397 + int32(-16)
	v467 = int32(base.Ui32(v396) >> (uint(int32(16)) % 32))
	if int32(31) < v397 {
		v396 = v467
		v397 = v465
		v398 = v462
		v399 = v455
		goto L68
	} else {
		goto L85
	}
L71:
	;
	v408 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v409 = v399 - v408
	v412 = base.I64_extend_i32_s(v409) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v412) {
		v436 = v408
		goto L73
	} else {
		goto L74
	}
L72:
	;
	if v398 == v408 {
		goto L83
	} else {
		goto L84
	}
L73:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v436
	goto L60
L74:
	;
	v415 = v398 - v408
	v417 = v412 + base.I64_extend_i32_u(v415)
	if base.Ui64(int64(4294967295)) < base.Ui64(v417) {
		v436 = v408
		goto L73
	} else {
		goto L75
	}
L75:
	;
	v420 = base.I32_wrap_i64(v417)
	if v399 == v408 {
		goto L76
	} else {
		goto L77
	}
L76:
	;
	v427 = int32(base.Ui32(v409*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v420) < base.Ui32(v427) {
		goto L79
	} else {
		goto L80
	}
L77:
	;
	if base.Ui32(v420) <= base.Ui32(v409) {
		v454 = v398
		v455 = v399
		goto L70
	} else {
		goto L78
	}
L78:
	;
	goto L76
L79:
	;
	v429 = v427
	goto L81
L80:
	;
	v429 = v420
	goto L81
L81:
	;
	v433 = v429&int32(-1024) + int32(1024)
	v434 = F_WebPSafeMalloc(m, int64(1), v433)
	mBase = m.M
	if v434 != 0 {
		goto L72
	} else {
		goto L82
	}
L82:
	;
	v435 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v436 = v435
	goto L73
L83:
	;
	v447 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_WebPSafeFree(m, v447)
	mBase = m.M
	v449 = v434 + v433
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v449
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v434
	v454 = v434 + v415
	v455 = v449
	goto L70
L84:
	;
	v445 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v446 = F_memcpy(m, v434, v445, v415)
	mBase = m.M
	goto L83
L85:
	;
	goto L69
L86:
	;
	v983 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	if v983 != 0 {
		goto L165
	} else {
		goto L166
	}
L87:
	;
	v514 = v511
	goto L89
L88:
	;
	v514 = v512
	goto L89
L89:
	;
	v515 = *(*int32)(unsafe.Add(mBase, uint32(v19)+40))
	if v515 < v514 {
		goto L90
	} else {
		goto L91
	}
L90:
	;
	v517 = v514
	goto L92
L91:
	;
	v517 = v515
	goto L92
L92:
	;
	v518 = *(*int32)(unsafe.Add(mBase, uint32(v19)+52))
	if v518 < v517 {
		goto L93
	} else {
		goto L94
	}
L93:
	;
	v520 = v517
	goto L95
L94:
	;
	v520 = v518
	goto L95
L95:
	;
	v521 = *(*int32)(unsafe.Add(mBase, uint32(v19)+64))
	if v521 < v520 {
		goto L96
	} else {
		goto L97
	}
L96:
	;
	v523 = v520
	goto L98
L97:
	;
	v523 = v521
	goto L98
L98:
	;
	v524 = int32(0)
	if v524 < v523 {
		goto L99
	} else {
		goto L100
	}
L99:
	;
	v527 = v523
	goto L101
L100:
	;
	v527 = v524
	goto L101
L101:
	;
	v528 = base.I64_extend_i32_u(v527)
	v529 = int32(2)
	if v528 == int64(0) {
		goto L104
	} else {
		goto L105
	}
L102:
	;
	if v550 == int32(0) {
		goto L86
	} else {
		goto L108
	}
L103:
	;
	goto L102
L104:
	;
	v548 = F_malloc(m, base.I32_wrap_i64(v528)*v529)
	mBase = m.M
	v550 = v548
	goto L103
L105:
	;
	v536 = base.I64_div_u_s(int64(2147418112), v528)
	v537 = int32(0)
	v538 = base.I64_extend_i32_u(v529)
	if base.Ui64(int64(4294967295)) < base.Ui64(v538*v528) {
		v550 = v537
		goto L103
	} else {
		goto L106
	}
L106:
	;
	if base.Ui64(v536) < base.Ui64(v538) {
		v550 = v537
		goto L103
	} else {
		goto L107
	}
L107:
	;
	goto L104
L108:
	;
	v555 = v19 + int32(16)
	F_StoreHuffmanCode(m, l0, v78, v550, v555)
	mBase = m.M
	v561 = *(*int32)(unsafe.Add(mBase, uint32(v19)+16))
	if v561 <= int32(0) {
		goto L109
	} else {
		goto L110
	}
L109:
	;
	F_StoreHuffmanCode(m, l0, v78, v550, v555|int32(12))
	mBase = m.M
	v645 = *(*int32)(unsafe.Add(mBase, uint32(v19)+28))
	if v645 < int32(1) {
		goto L120
	} else {
		goto L121
	}
L110:
	;
	v564 = *(*int32)(unsafe.Add(mBase, uint32(v19)+20))
	v572 = int32(0)
	v575 = v561
	v576 = v564
	goto L111
L111:
	;
	v582 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v576))))
	if v582 == int32(0) {
		v589 = v572
		goto L113
	} else {
		goto L114
	}
L112:
	;
	v594 = int32(0)
	v605 = v594
	v606 = v594
	goto L117
L113:
	;
	v593 = v575 + int32(-1)
	if v593 != 0 {
		v572 = v589
		v575 = v593
		v576 = v576 + int32(1)
		goto L111
	} else {
		goto L116
	}
L114:
	;
	if int32(0) < v572 {
		goto L109
	} else {
		goto L115
	}
L115:
	;
	v589 = int32(1)
	goto L113
L116:
	;
	goto L112
L117:
	;
	v612 = *(*int32)(unsafe.Add(mBase, uint32(v19)+20))
	v614 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v612+v605))) = uint8(v614)
	v616 = *(*int32)(unsafe.Add(mBase, uint32(v19)+24))
	*(*uint16)(unsafe.Add(mBase, uint32(v616+v606))) = uint16(v614)
	v623 = v605 + int32(1)
	v624 = *(*int32)(unsafe.Add(mBase, uint32(v19)+16))
	if v623 < v624 {
		v605 = v623
		v606 = v606 + int32(2)
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
	F_StoreHuffmanCode(m, l0, v78, v550, v19+int32(40))
	mBase = m.M
	v729 = *(*int32)(unsafe.Add(mBase, uint32(v19)+40))
	if v729 < int32(1) {
		goto L131
	} else {
		goto L132
	}
L121:
	;
	v648 = *(*int32)(unsafe.Add(mBase, uint32(v19)+32))
	v656 = int32(0)
	v659 = v645
	v660 = v648
	goto L122
L122:
	;
	v666 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v660))))
	if v666 == int32(0) {
		v673 = v656
		goto L124
	} else {
		goto L125
	}
L123:
	;
	v678 = int32(0)
	v689 = v678
	v690 = v678
	goto L128
L124:
	;
	v677 = v659 + int32(-1)
	if v677 != 0 {
		v656 = v673
		v659 = v677
		v660 = v660 + int32(1)
		goto L122
	} else {
		goto L127
	}
L125:
	;
	if int32(0) < v656 {
		goto L120
	} else {
		goto L126
	}
L126:
	;
	v673 = int32(1)
	goto L124
L127:
	;
	goto L123
L128:
	;
	v696 = *(*int32)(unsafe.Add(mBase, uint32(v19)+32))
	v698 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v696+v689))) = uint8(v698)
	v700 = *(*int32)(unsafe.Add(mBase, uint32(v19)+36))
	*(*uint16)(unsafe.Add(mBase, uint32(v700+v690))) = uint16(v698)
	v707 = v689 + int32(1)
	v708 = *(*int32)(unsafe.Add(mBase, uint32(v19)+28))
	if v707 < v708 {
		v689 = v707
		v690 = v690 + int32(2)
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
	F_StoreHuffmanCode(m, l0, v78, v550, v19+int32(52))
	mBase = m.M
	v813 = *(*int32)(unsafe.Add(mBase, uint32(v19)+52))
	if v813 < int32(1) {
		goto L142
	} else {
		goto L143
	}
L132:
	;
	v732 = *(*int32)(unsafe.Add(mBase, uint32(v19)+44))
	v740 = int32(0)
	v743 = v729
	v744 = v732
	goto L133
L133:
	;
	v750 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v744))))
	if v750 == int32(0) {
		v757 = v740
		goto L135
	} else {
		goto L136
	}
L134:
	;
	v762 = int32(0)
	v773 = v762
	v774 = v762
	goto L139
L135:
	;
	v761 = v743 + int32(-1)
	if v761 != 0 {
		v740 = v757
		v743 = v761
		v744 = v744 + int32(1)
		goto L133
	} else {
		goto L138
	}
L136:
	;
	if int32(0) < v740 {
		goto L131
	} else {
		goto L137
	}
L137:
	;
	v757 = int32(1)
	goto L135
L138:
	;
	goto L134
L139:
	;
	v780 = *(*int32)(unsafe.Add(mBase, uint32(v19)+44))
	v782 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v780+v773))) = uint8(v782)
	v784 = *(*int32)(unsafe.Add(mBase, uint32(v19)+48))
	*(*uint16)(unsafe.Add(mBase, uint32(v784+v774))) = uint16(v782)
	v791 = v773 + int32(1)
	v792 = *(*int32)(unsafe.Add(mBase, uint32(v19)+40))
	if v791 < v792 {
		v773 = v791
		v774 = v774 + int32(2)
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
	F_StoreHuffmanCode(m, l0, v78, v550, v19+int32(64))
	mBase = m.M
	v895 = *(*int32)(unsafe.Add(mBase, uint32(v19)+64))
	if v895 < int32(1) {
		goto L153
	} else {
		goto L154
	}
L143:
	;
	v816 = *(*int32)(unsafe.Add(mBase, uint32(v19)+56))
	v824 = int32(0)
	v827 = v813
	v828 = v816
	goto L144
L144:
	;
	v834 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v828))))
	if v834 == int32(0) {
		v841 = v824
		goto L146
	} else {
		goto L147
	}
L145:
	;
	v846 = int32(0)
	v857 = v846
	v858 = v846
	goto L150
L146:
	;
	v845 = v827 + int32(-1)
	if v845 != 0 {
		v824 = v841
		v827 = v845
		v828 = v828 + int32(1)
		goto L144
	} else {
		goto L149
	}
L147:
	;
	if int32(0) < v824 {
		goto L142
	} else {
		goto L148
	}
L148:
	;
	v841 = int32(1)
	goto L146
L149:
	;
	goto L145
L150:
	;
	v864 = *(*int32)(unsafe.Add(mBase, uint32(v19)+56))
	v866 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v864+v857))) = uint8(v866)
	v868 = *(*int32)(unsafe.Add(mBase, uint32(v19)+60))
	*(*uint16)(unsafe.Add(mBase, uint32(v868+v858))) = uint16(v866)
	v875 = v857 + int32(1)
	v876 = *(*int32)(unsafe.Add(mBase, uint32(v19)+52))
	if v875 < v876 {
		v857 = v875
		v858 = v858 + int32(2)
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
	v981 = F_StoreImageToBitMask(m, l0, l4, int32(0), l3, v19+int32(12), v19+int32(16), l8)
	mBase = m.M
	v999 = v103
	v1001 = v550
	goto L1
L154:
	;
	v898 = *(*int32)(unsafe.Add(mBase, uint32(v19)+68))
	v906 = int32(0)
	v909 = v895
	v910 = v898
	goto L155
L155:
	;
	v916 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v910))))
	if v916 == int32(0) {
		v923 = v906
		goto L157
	} else {
		goto L158
	}
L156:
	;
	v928 = int32(0)
	v939 = v928
	v940 = v928
	goto L161
L157:
	;
	v927 = v909 + int32(-1)
	if v927 != 0 {
		v906 = v923
		v909 = v927
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
	v923 = int32(1)
	goto L157
L160:
	;
	goto L156
L161:
	;
	v946 = *(*int32)(unsafe.Add(mBase, uint32(v19)+68))
	v948 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v946+v939))) = uint8(v948)
	v950 = *(*int32)(unsafe.Add(mBase, uint32(v19)+72))
	*(*uint16)(unsafe.Add(mBase, uint32(v950+v940))) = uint16(v948)
	v957 = v939 + int32(1)
	v958 = *(*int32)(unsafe.Add(mBase, uint32(v19)+64))
	if v957 < v958 {
		v939 = v957
		v940 = v940 + int32(2)
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
	v999 = v103
	v1001 = int32(0)
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
	F_free(m, v78)
	mBase = m.M
	goto L168
L168:
	;
	F_WebPSafeFree(m, v999)
	mBase = m.M
	goto L169
L169:
	;
	v1006 = *(*int32)(unsafe.Add(mBase, uint32(v19)+24))
	F_free(m, v1006)
	mBase = m.M
	goto L170
L170:
	;
	v1008 = *(*int32)(unsafe.Add(mBase, uint32(l8)+92))
	m.G0 = v19 + int32(80)
	return base.B2i32(v1008 == int32(0))
}
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
func F_ExpandMatrix(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v8 int32
	_ = v8
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v22 int32
	_ = v22
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v49 int32
	_ = v49
	var v88 int32
	_ = v88
	var v102 int32
	_ = v102
	var v106 int32
	_ = v106
	var v111 int64
	_ = v111
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v182 int32
	_ = v182
	var v185 int32
	_ = v185
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v206 int32
	_ = v206
	var v210 int32
	_ = v210
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v236 int32
	_ = v236
	v3 = int32(0)
	v8 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+2)))
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)) = uint16(v8)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)) = uint16(v8)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)) = uint16(v8)
	v12 = int32(131072)
	v13 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0))))
	v14 = base.I32_div_u_s(v12, v13)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+32)) = uint16(v14)
	v17 = base.I32_div_u_s(v12, v8)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+34)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+36)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+38)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+40)) = uint16(v17)
	v22 = m.G1
	v27 = v22 + int32(_a_F_ExpandMatrix_0) + l1<<(uint(int32(1))%32)
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v27))))
	v29 = int32(9)
	v30 = v28 << (uint(v29) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+64)) = v30
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v27)+1)))
	v34 = v32 << (uint(v29) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+68)) = v34
	*(*int32)(unsafe.Add(mBase, uint32(l0)+72)) = v34
	*(*int32)(unsafe.Add(mBase, uint32(l0)+76)) = v34
	*(*int32)(unsafe.Add(mBase, uint32(l0)+80)) = v34
	v39 = int32(131071)
	v41 = int32(_a_F_ExpandMatrix_1)
	v43 = base.I32_div_u_s(v30^v39, v14&v41)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+128)) = v43
	v49 = base.I32_div_u_s(v34^v39, v17&v41)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+132)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+136)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+140)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+144)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+148)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+84)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+42)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+10)) = uint16(v8)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+12)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+152)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+88)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+44)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+14)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+46)) = uint16(v17)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+48)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+16)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+164)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+100)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+50)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+18)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+168)) = v49
	*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+52)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+20)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+108)) = v34
	*(*int32)(unsafe.Add(mBase, uint32(l0)+172)) = v49
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+54)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+22)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+112)) = v34
	*(*int32)(unsafe.Add(mBase, uint32(l0)+176)) = v49
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+56)) = uint16(v17)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+24)) = uint16(v8)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+26)) = uint16(v8)
	v88 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+34)))
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+58)) = uint16(v88)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+116)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+60)) = uint16(v88)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+28)) = uint16(v8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+120)) = v34
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+30)) = uint16(v8)
	*(*uint16)(unsafe.Add(mBase, uint32(l0)+62)) = uint16(v88)
	*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(192)))) = uint16(v3)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+124)) = v34
	v102 = *(*int32)(unsafe.Add(mBase, uint32(l0)+132))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+180)) = v102
	*(*int32)(unsafe.Add(mBase, uint32(l0)+184)) = v102
	*(*int32)(unsafe.Add(mBase, uint32(l0)+188)) = v102
	v106 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0))))
	if l1 == v3 {
		v146 = int32(30)
		v148 = int32(11)
		v149 = int32(base.Ui32(v8*v146) >> (uint(v148) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(194)))) = uint16(v149)
		v153 = int32(90)
		v156 = int32(base.Ui32(v8*v153) >> (uint(v148) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(222)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(220)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(218)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(216)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(214)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(212)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(210)))) = uint16(v156)
		v182 = int32(60)
		v185 = int32(base.Ui32(v8*v182) >> (uint(v148) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(208)))) = uint16(v185)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(206)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(204)))) = uint16(v156)
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(202)))) = uint16(v185)
		v198 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)))
		v202 = int32(base.Ui32(v198*v182) >> (uint(v148) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(196)))) = uint16(v202)
		v206 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)))
		v210 = int32(base.Ui32(v206*v153) >> (uint(v148) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(198)))) = uint16(v210)
		v214 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)))
		v218 = int32(base.Ui32(v214*v146) >> (uint(v148) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(200)))) = uint16(v218)
		v220 = int32(1)
		v221 = v8 << (uint(v220) % 32)
		v236 = v214 + (v206 + (v198 + v221)) + v221 + v8<<(uint(int32(2))%32) + ((v8+v8)<<(uint(v220)%32) + v106)
	} else {
		v111 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(194)))) = v111
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(202)))) = v111
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(210)))) = v111
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(216)))) = v111
		v125 = int32(1)
		v126 = v8 << (uint(v125) % 32)
		v127 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+4)))
		v129 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+6)))
		v131 = int32(*(*uint16)(unsafe.Add(mBase, uint32(l0)+8)))
		v236 = v126 + v127 + v129 + v131 + v126 + v8<<(uint(int32(2))%32) + ((v8+v8)<<(uint(v125)%32) + v106)
	}
	return int32(base.Ui32(v236+int32(8)) >> (uint(int32(4)) % 32))
}
func F_ExtraCost_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v75 int32
	_ = v75
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v9 = v7 + v8
	if l1 < int32(8) {
		v75 = v9
	} else {
		v12 = int32(1)
		v13 = int32(base.Ui32(l1) >> (uint(v12) % 32))
		v15 = v13 + int32(-3)
		if v13 != int32(4) {
			v27 = l0 + int32(36)
			v28 = v9
			v29 = int32(0)
			for {
				v32 = *(*int32)(unsafe.Add(mBase, uint32(v27)))
				v35 = *(*int32)(unsafe.Add(mBase, uint32(v27+int32(-4))))
				v42 = *(*int32)(unsafe.Add(mBase, uint32(v27+int32(-8))))
				v45 = *(*int32)(unsafe.Add(mBase, uint32(v27+int32(-12))))
				v48 = v29 + int32(2)
				v51 = (v32+v35)*(v29+int32(3)) + ((v42+v45)*v48 + v28)
				if v15&int32(-2) != v48 {
					v27 = v27 + int32(16)
					v28 = v51
					v29 = v48
					continue
				} else {
					break
				}
				break
			}
			v58 = v29 + int32(4)
			v59 = v51
		} else {
			v58 = int32(2)
			v59 = v9
		}
		if v15&v12 == int32(0) {
			v75 = v59
		} else {
			v67 = l0 + v58<<(uint(int32(3))%32)
			v68 = *(*int32)(unsafe.Add(mBase, uint32(v67)+12))
			v69 = *(*int32)(unsafe.Add(mBase, uint32(v67)+8))
			v75 = (v68+v69)*v58 + v59
		}
	}
	return v75
}
func F_ExtractAlpha_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v25 int32
	_ = v25
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v52 int32
	_ = v52
	var v56 int32
	_ = v56
	var v61 int32
	_ = v61
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v79 int32
	_ = v79
	var v85 int32
	_ = v85
	var v90 int32
	_ = v90
	var v94 int32
	_ = v94
	var v102 int32
	_ = v102
	var v107 int32
	_ = v107
	var v120 int32
	_ = v120
	var v124 int32
	_ = v124
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v149 int32
	_ = v149
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v174 int32
	_ = v174
	v17 = int32(1)
	if l3 < v17 {
		v174 = v17
	} else {
		if l2 < int32(1) {
			v174 = v17
		} else {
			v25 = l2 & int32(3)
			v30 = l0
			v34 = l4
			v36 = int32(255)
			v39 = int32(0)
			for {
				if base.B2i32(base.Ui32(l2) < base.Ui32(int32(4))) == int32(0) {
					v52 = v30
					v56 = v36
					v61 = int32(0)
					for {
						v66 = v34 + v61
						v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v52))))
						*(*uint8)(unsafe.Add(mBase, uint32(v66))) = uint8(v67)
						v71 = int32(4)
						v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v52+v71))))
						*(*uint8)(unsafe.Add(mBase, uint32(v66+int32(1)))) = uint8(v73)
						v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v52+int32(8)))))
						*(*uint8)(unsafe.Add(mBase, uint32(v66+int32(2)))) = uint8(v79)
						v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v52+int32(12)))))
						*(*uint8)(unsafe.Add(mBase, uint32(v66+int32(3)))) = uint8(v85)
						v90 = v85 & (v79 & (v73 & (v67 & v56)))
						v94 = v61 + v71
						if l2&int32(2147483644) != v94 {
							v52 = v52 + int32(16)
							v56 = v90
							v61 = v94
							continue
						} else {
							break
						}
						break
					}
					v102 = v90
					v107 = v94
				} else {
					v102 = v36
					v107 = int32(0)
				}
				if v25 == int32(0) {
					v149 = v102
				} else {
					v120 = v34 + v107
					v124 = v102
					v129 = v25
					v130 = v30 + v107<<(uint(int32(2))%32)
					for {
						v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130))))
						*(*uint8)(unsafe.Add(mBase, uint32(v120))) = uint8(v134)
						v140 = v134 & v124
						v142 = v129 + int32(-1)
						if v142 != 0 {
							v120 = v120 + int32(1)
							v124 = v140
							v129 = v142
							v130 = v130 + int32(4)
							continue
						} else {
							break
						}
						break
					}
					v149 = v140
				}
				v162 = v39 + int32(1)
				if v162 != l3 {
					v30 = v30 + l1
					v34 = v34 + l5
					v36 = v149
					v39 = v162
					continue
				} else {
					break
				}
				break
			}
			v164 = int32(255)
			v174 = base.B2i32(v149&v164 == v164)
		}
	}
	return v174
}
func F_ExtractGreen_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v42 int32
	_ = v42
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v57 int32
	_ = v57
	var v63 int32
	_ = v63
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v88 int32
	_ = v88
	if l2 < int32(1) {
	} else {
		v11 = l2 & int32(3)
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v63 = int32(0)
		} else {
			v20 = l0
			v22 = int32(0)
			for {
				v25 = l1 + v22
				v26 = *(*int32)(unsafe.Add(mBase, uint32(v20)))
				v27 = int32(8)
				v28 = int32(base.Ui32(v26) >> (uint(v27) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v25))) = uint8(v28)
				v32 = int32(4)
				v34 = *(*int32)(unsafe.Add(mBase, uint32(v20+v32)))
				v36 = int32(base.Ui32(v34) >> (uint(v27) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v25+int32(1)))) = uint8(v36)
				v42 = *(*int32)(unsafe.Add(mBase, uint32(v20+v27)))
				v44 = int32(base.Ui32(v42) >> (uint(v27) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v25+int32(2)))) = uint8(v44)
				v50 = *(*int32)(unsafe.Add(mBase, uint32(v20+int32(12))))
				v52 = int32(base.Ui32(v50) >> (uint(v27) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v25+int32(3)))) = uint8(v52)
				v57 = v22 + v32
				if l2&int32(2147483644) != v57 {
					v20 = v20 + int32(16)
					v22 = v57
					continue
				} else {
					break
				}
				break
			}
			v63 = v57
		}
		if v11 == int32(0) {
		} else {
			v74 = l1 + v63
			v75 = v11
			v78 = l0 + v63<<(uint(int32(2))%32)
			for {
				v79 = *(*int32)(unsafe.Add(mBase, uint32(v78)))
				v81 = int32(base.Ui32(v79) >> (uint(int32(8)) % 32))
				*(*uint8)(unsafe.Add(mBase, uint32(v74))) = uint8(v81)
				v88 = v75 + int32(-1)
				if v88 != 0 {
					v74 = v74 + int32(1)
					v75 = v88
					v78 = v78 + int32(4)
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
	var v38 int64
	_ = v38
	var v52 int32
	_ = v52
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v102 float32
	_ = v102
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v111 float32
	_ = v111
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v159 int32
	_ = v159
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var v185 int32
	_ = v185
	var v188 int32
	_ = v188
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v197 int32
	_ = v197
	var v205 int32
	_ = v205
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v224 int32
	_ = v224
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v242 int32
	_ = v242
	var v244 int32
	_ = v244
	var v246 int64
	_ = v246
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v280 int32
	_ = v280
	var v282 int32
	_ = v282
	var v284 int64
	_ = v284
	var v304 int32
	_ = v304
	var v311 int32
	_ = v311
	var v313 int32
	_ = v313
	var v314 int32
	_ = v314
	v9 = int32(0)
	v11 = m.G0
	v13 = v11 - int32(304)
	m.G0 = v13
	*(*int32)(unsafe.Add(mBase, uint32(l3))) = v9
	v19 = v13 + int32(188)
	if v19 == v9 {
		v85 = v9
	} else {
		v31 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+16)) = v31
		*(*float32)(unsafe.Add(mBase, uint32(v19)+4)) = float32(75)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+32)) = int64(60)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+8)) = int64(4)
		v38 = int64(429496729600)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+108)) = v38
		*(*int64)(unsafe.Add(mBase, uint32(v19)+24)) = int64(214748364804)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+64)) = v31
		*(*int64)(unsafe.Add(mBase, uint32(v19)+72)) = v31
		*(*int64)(unsafe.Add(mBase, uint32(v19)+40)) = int64(1)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+56)) = int64(4294967396)
		*(*int64)(unsafe.Add(mBase, uint32(v19)+48)) = int64(4294967297)
		v52 = int32(0)
		*(*int32)(unsafe.Add(mBase, uint32(v19))) = v52
		*(*int32)(unsafe.Add(mBase, uint32(v19)+96)) = v52
		*(*int64)(unsafe.Add(mBase, uint32(v19)+88)) = v38
		*(*int64)(unsafe.Add(mBase, uint32(v19)+80)) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v19)+104)) = v52
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
		v84 = F_WebPValidateConfig(m, v19)
		mBase = m.M
		v85 = v84
	}
	if v85 == int32(0) {
		v314 = v9
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(v13)+196)) = l5
		*(*int32)(unsafe.Add(mBase, uint32(v13)+284)) = l7
		*(*int32)(unsafe.Add(mBase, uint32(v13)+188)) = l6
		*(*float32)(unsafe.Add(mBase, uint32(v13)+192)) = base.F32_convert_i32_s(l4)
		v94 = v13 + int32(188)
		v95 = int32(0)
		if v94 == v95 {
			v197 = v95
		} else {
			v102 = *(*float32)(unsafe.Add(mBase, uint32(v94)+4))
			if base.F32_lt(v102, float32(0)) != 0 {
				v197 = v95
			} else {
				if base.F32_gt(v102, float32(100)) != 0 {
					v197 = v95
				} else {
					v107 = int32(0)
					v108 = *(*int32)(unsafe.Add(mBase, uint32(v94)+16))
					if v108 < v107 {
						v197 = v107
					} else {
						v111 = *(*float32)(unsafe.Add(mBase, uint32(v94)+20))
						if base.F32_lt(v111, float32(0)) != 0 {
							v197 = v107
						} else {
							v114 = *(*int32)(unsafe.Add(mBase, uint32(v94)+8))
							if base.Ui32(int32(6)) < base.Ui32(v114) {
								v197 = v107
							} else {
								v117 = *(*int32)(unsafe.Add(mBase, uint32(v94)+24))
								if base.Ui32(v117+int32(-5)) < base.Ui32(int32(-4)) {
									v197 = v107
								} else {
									v122 = *(*int32)(unsafe.Add(mBase, uint32(v94)+28))
									if base.Ui32(int32(100)) < base.Ui32(v122) {
										v197 = v107
									} else {
										v125 = *(*int32)(unsafe.Add(mBase, uint32(v94)+32))
										if base.Ui32(int32(100)) < base.Ui32(v125) {
											v197 = v107
										} else {
											v128 = *(*int32)(unsafe.Add(mBase, uint32(v94)+36))
											if base.Ui32(int32(7)) < base.Ui32(v128) {
												v197 = v107
											} else {
												v131 = *(*int32)(unsafe.Add(mBase, uint32(v94)+40))
												if base.Ui32(int32(1)) < base.Ui32(v131) {
													v197 = v107
												} else {
													v134 = *(*int32)(unsafe.Add(mBase, uint32(v94)+44))
													if base.Ui32(int32(1)) < base.Ui32(v134) {
														v197 = v107
													} else {
														v137 = *(*int32)(unsafe.Add(mBase, uint32(v94)+60))
														if base.Ui32(v137+int32(-11)) < base.Ui32(int32(-10)) {
															v197 = v107
														} else {
															v142 = int32(0)
															v143 = *(*int32)(unsafe.Add(mBase, uint32(v94)+108))
															if v143 < v142 {
																v197 = v142
															} else {
																v146 = *(*int32)(unsafe.Add(mBase, uint32(v94)+112))
																if int32(100) < v146 {
																	v197 = v142
																} else {
																	if v146 < v143 {
																		v197 = v142
																	} else {
																		v150 = *(*int32)(unsafe.Add(mBase, uint32(v94)+64))
																		if base.Ui32(int32(1)) < base.Ui32(v150) {
																			v197 = v142
																		} else {
																			v153 = *(*int32)(unsafe.Add(mBase, uint32(v94)+68))
																			if base.Ui32(int32(7)) < base.Ui32(v153) {
																				v197 = v142
																			} else {
																				v156 = *(*int32)(unsafe.Add(mBase, uint32(v94)+72))
																				if base.Ui32(int32(3)) < base.Ui32(v156) {
																					v197 = v142
																				} else {
																					v159 = *(*int32)(unsafe.Add(mBase, uint32(v94)+76))
																					if base.Ui32(int32(100)) < base.Ui32(v159) {
																						v197 = v142
																					} else {
																						v162 = int32(0)
																						v163 = *(*int32)(unsafe.Add(mBase, uint32(v94)+48))
																						if v163 < v162 {
																							v197 = v162
																						} else {
																							v166 = int32(0)
																							v167 = *(*int32)(unsafe.Add(mBase, uint32(v94)+52))
																							if v167 < v166 {
																								v197 = v166
																							} else {
																								v170 = *(*int32)(unsafe.Add(mBase, uint32(v94)+56))
																								if base.Ui32(int32(100)) < base.Ui32(v170) {
																									v197 = v166
																								} else {
																									v173 = *(*int32)(unsafe.Add(mBase, uint32(v94)))
																									if base.Ui32(int32(1)) < base.Ui32(v173) {
																										v197 = v166
																									} else {
																										v176 = *(*int32)(unsafe.Add(mBase, uint32(v94)+92))
																										if base.Ui32(int32(100)) < base.Ui32(v176) {
																											v197 = v166
																										} else {
																											v179 = *(*int32)(unsafe.Add(mBase, uint32(v94)+12))
																											if base.Ui32(int32(3)) < base.Ui32(v179) {
																												v197 = v166
																											} else {
																												v182 = *(*int32)(unsafe.Add(mBase, uint32(v94)+80))
																												if base.Ui32(int32(1)) < base.Ui32(v182) {
																													v197 = v166
																												} else {
																													v185 = *(*int32)(unsafe.Add(mBase, uint32(v94)+84))
																													if base.Ui32(int32(1)) < base.Ui32(v185) {
																														v197 = v166
																													} else {
																														v188 = *(*int32)(unsafe.Add(mBase, uint32(v94)+88))
																														if base.Ui32(int32(1)) < base.Ui32(v188) {
																															v197 = v166
																														} else {
																															v191 = *(*int32)(unsafe.Add(mBase, uint32(v94)+96))
																															if base.Ui32(int32(1)) < base.Ui32(v191) {
																																v197 = v166
																															} else {
																																v194 = *(*int32)(unsafe.Add(mBase, uint32(v94)+104))
																																v197 = base.B2i32(base.Ui32(v194) < base.Ui32(int32(2)))
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
		if v197 == int32(0) {
			v314 = v9
		} else {
			v205 = v13 + int32(16)
			if v205 == int32(0) {
			} else {
				v216 = int32(0)
				v218 = F_memset(m, v205, v216, int32(172))
				mBase = m.M
				v219 = m.G2
				*(*int32)(unsafe.Add(mBase, uint32(v218)+72)) = v219 + int32(1)
				v224 = F_WebPEncodingSetError(m, v218, v216)
				mBase = m.M
			}
			*(*int32)(unsafe.Add(mBase, uint32(v13)+28)) = l2
			*(*int32)(unsafe.Add(mBase, uint32(v13)+16)) = int32(1)
			*(*int32)(unsafe.Add(mBase, uint32(v13)+24)) = l1
			v237 = F_WebPPictureImportRGBA(m, v13+int32(16), l0, l1<<(uint(int32(2))%32))
			mBase = m.M
			if v237 != 0 {
				*(*int32)(unsafe.Add(mBase, uint32(v13)+8)) = int32(0)
				*(*int64)(unsafe.Add(mBase, uint32(v13))) = int64(0)
				*(*int32)(unsafe.Add(mBase, uint32(v13)+88)) = int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(v13)+92)) = v13
				v274 = v13 + int32(16)
				v275 = F_WebPEncode(m, v13+int32(188), v274)
				mBase = m.M
				if v274 == int32(0) {
				} else {
					v280 = *(*int32)(unsafe.Add(mBase, uint32(v274)+156))
					F_WebPSafeFree(m, v280)
					mBase = m.M
					v282 = *(*int32)(unsafe.Add(mBase, uint32(v274)+160))
					F_WebPSafeFree(m, v282)
					mBase = m.M
					v284 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v274)+52)) = v284
					*(*int64)(unsafe.Add(mBase, uint32(v274)+156)) = v284
					*(*int64)(unsafe.Add(mBase, uint32(v274)+16)) = v284
					*(*int64)(unsafe.Add(mBase, uint32(v13+int32(40)))) = v284
					*(*int64)(unsafe.Add(mBase, uint32(v13+int32(48)))) = v284
					*(*int32)(unsafe.Add(mBase, uint32(v13+int32(56)))) = int32(0)
				}
				if v275 != 0 {
					v311 = *(*int32)(unsafe.Add(mBase, uint32(v13)+4))
					*(*int32)(unsafe.Add(mBase, uint32(l3))) = v311
					v313 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
					v314 = v313
				} else {
					if v13 == int32(0) {
					} else {
						v304 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
						F_WebPSafeFree(m, v304)
						mBase = m.M
						*(*int32)(unsafe.Add(mBase, uint32(v13)+8)) = int32(0)
						*(*int64)(unsafe.Add(mBase, uint32(v13))) = int64(0)
					}
					v314 = int32(0)
				}
			} else {
				v239 = v13 + int32(16)
				if v239 == int32(0) {
				} else {
					v242 = *(*int32)(unsafe.Add(mBase, uint32(v239)+156))
					F_WebPSafeFree(m, v242)
					mBase = m.M
					v244 = *(*int32)(unsafe.Add(mBase, uint32(v239)+160))
					F_WebPSafeFree(m, v244)
					mBase = m.M
					v246 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v239)+52)) = v246
					*(*int64)(unsafe.Add(mBase, uint32(v239)+156)) = v246
					*(*int64)(unsafe.Add(mBase, uint32(v239)+16)) = v246
					*(*int64)(unsafe.Add(mBase, uint32(v13+int32(40)))) = v246
					*(*int64)(unsafe.Add(mBase, uint32(v13+int32(48)))) = v246
					*(*int32)(unsafe.Add(mBase, uint32(v13+int32(56)))) = int32(0)
				}
				v314 = int32(0)
			}
		}
	}
	m.G0 = v13 + int32(304)
	return v314
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

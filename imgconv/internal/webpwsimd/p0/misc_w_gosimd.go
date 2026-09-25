//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_WebPConvertARGBToUV_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v34 int32
	_ = v34
	var v42 base.V128
	_ = v42
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v56 int32
	_ = v56
	var v57 base.V128
	_ = v57
	var v59 base.V128
	_ = v59
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v75 int32
	_ = v75
	var v81 base.V128
	_ = v81
	var v84 base.V128
	_ = v84
	var v87 base.V128
	_ = v87
	var v90 base.V128
	_ = v90
	var v91 int32
	_ = v91
	var v93 base.V128
	_ = v93
	var v116 base.V128
	_ = v116
	var v119 base.V128
	_ = v119
	var v122 base.V128
	_ = v122
	var v125 base.V128
	_ = v125
	var v129 base.V128
	_ = v129
	var v132 int32
	_ = v132
	var v138 base.V128
	_ = v138
	var v148 base.V128
	_ = v148
	var v149 base.V128
	_ = v149
	var v152 base.V128
	_ = v152
	var v154 int32
	_ = v154
	var v156 base.V128
	_ = v156
	var v176 int32
	_ = v176
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	var v198 base.V128
	_ = v198
	var v199 int32
	_ = v199
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v212 int32
	_ = v212
	var v213 base.V128
	_ = v213
	var v215 base.V128
	_ = v215
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v231 int32
	_ = v231
	var v237 base.V128
	_ = v237
	var v240 base.V128
	_ = v240
	var v243 base.V128
	_ = v243
	var v246 base.V128
	_ = v246
	var v247 int32
	_ = v247
	var v249 base.V128
	_ = v249
	var v272 base.V128
	_ = v272
	var v275 base.V128
	_ = v275
	var v278 base.V128
	_ = v278
	var v281 base.V128
	_ = v281
	var v285 base.V128
	_ = v285
	var v288 int32
	_ = v288
	var v294 base.V128
	_ = v294
	var v304 base.V128
	_ = v304
	var v305 base.V128
	_ = v305
	var v308 base.V128
	_ = v308
	var v310 int32
	_ = v310
	var v312 base.V128
	_ = v312
	var v314 int32
	_ = v314
	var v316 int32
	_ = v316
	var v318 int32
	_ = v318
	var v321 base.V128
	_ = v321
	var v323 base.V128
	_ = v323
	var v324 base.V128
	_ = v324
	var v330 base.V128
	_ = v330
	var v354 base.V128
	_ = v354
	var v356 base.V128
	_ = v356
	var v357 base.V128
	_ = v357
	var v375 int32
	_ = v375
	var v383 int32
	_ = v383
	var v407 int32
	_ = v407
	var v410 int32
	_ = v410
	var v422 int32
	_ = v422
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v430 int32
	_ = v430
	var v435 int32
	_ = v435
	var v438 int32
	_ = v438
	var v446 int32
	_ = v446
	var v450 int32
	_ = v450
	var v458 int32
	_ = v458
	var v459 int32
	_ = v459
	var v462 int32
	_ = v462
	var v464 int32
	_ = v464
	var v466 int32
	_ = v466
	var v471 int32
	_ = v471
	var v473 int32
	_ = v473
	var v486 int32
	_ = v486
	var v491 int32
	_ = v491
	var v496 int32
	_ = v496
	var v504 int32
	_ = v504
	var v528 int32
	_ = v528
	var v531 int32
	_ = v531
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v549 int32
	_ = v549
	var v551 int32
	_ = v551
	var v556 int32
	_ = v556
	var v559 int32
	_ = v559
	var v567 int32
	_ = v567
	var v571 int32
	_ = v571
	var v579 int32
	_ = v579
	var v580 int32
	_ = v580
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v600 int32
	_ = v600
	var v605 int32
	_ = v605
	var v612 int32
	_ = v612
	var v635 int32
	_ = v635
	var v638 int32
	_ = v638
	var v639 int32
	_ = v639
	var v640 int32
	_ = v640
	var v645 int32
	_ = v645
	var v652 int32
	_ = v652
	var v656 int32
	_ = v656
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v671 int32
	_ = v671
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v681 int32
	_ = v681
	var v684 int32
	_ = v684
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v692 int32
	_ = v692
	v22 = int32(1)
	v23 = l3 >> (uint(v22) % 32)
	if v22 <= v23 {
		if l4 == int32(0) {
			if base.Ui32(v23) < base.Ui32(int32(4)) {
				v383 = int32(0)
				v407 = v383
				v410 = l0 + v383<<(uint(int32(3))%32)
				for {
					v422 = l1 + v407
					v425 = *(*int32)(unsafe.Add(mBase, uint32(v410+int32(4))))
					v426 = int32(15)
					v428 = int32(510)
					v430 = *(*int32)(unsafe.Add(mBase, uint32(v410)))
					v435 = int32(base.Ui32(v425)>>(uint(v426)%32))&v428 + int32(base.Ui32(v430)>>(uint(v426)%32))&v428
					v438 = int32(7)
					v446 = int32(base.Ui32(v425)>>(uint(v438)%32))&v428 + int32(base.Ui32(v430)>>(uint(v438)%32))&v428
					v450 = int32(1)
					v458 = v425<<(uint(v450)%32)&v428 + v430<<(uint(v450)%32)&v428
					v459 = int32(_a_F_WebPConvertARGBToUV_C_0)
					v462 = int32(33685504)
					v464 = int32(18)
					v466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v422))))
					v471 = int32(base.Ui32(int32(base.Ui32(v435*int32(-9719)+v446*int32(-19081)+v458*v459+v462)>>(uint(v464)%32))+v466+v450) >> (uint(v450) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v422))) = uint8(v471)
					v473 = l2 + v407
					v486 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v473))))
					v491 = int32(base.Ui32(int32(base.Ui32(v435*v459+v446*int32(-24116)+v458*int32(-4684)+v462)>>(uint(v464)%32))+v486+v450) >> (uint(v450) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(v473))) = uint8(v491)
					v496 = v407 + v450
					if v23 != v496 {
						v407 = v496
						v410 = v410 + int32(8)
						continue
					} else {
						break
					}
					break
				}
				v612 = v23
			} else {
				v190 = v23 & int32(2147483644)
				v198 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k0)
				v199 = v190
				v200 = l2
				v201 = l1
				for {
					v212 = int32(1)
					v213 = base.Simd_g_i32x4_shl(v198, v212)
					v215 = base.Simd_g_v128_or(v213, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k1))
					v216 = int32(3)
					v218 = int32(2)
					v231 = int32(0)
					v237 = base.Simd_g_v128_load32_splat(m, l0+base.Simd_g_i32x4_extract_lane_l0(v215)<<(uint(v218)%32), v231)
					v240 = base.Simd_g_v128_load32_lane_l1(m, l0+base.Simd_g_i32x4_extract_lane_l1(v215)<<(uint(v218)%32), v231, v237)
					v243 = base.Simd_g_v128_load32_lane_l2(m, l0+base.Simd_g_i32x4_extract_lane_l2(v215)<<(uint(v218)%32), v231, v240)
					v246 = base.Simd_g_v128_load32_lane_l3(m, l0+base.Simd_g_i32x4_extract_lane_l3(v215)<<(uint(v218)%32), v231, v243)
					v247 = int32(15)
					v249 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k2)
					v272 = base.Simd_g_v128_load32_splat(m, l0+base.Simd_g_i32x4_extract_lane_l0(v213)<<(uint(v218)%32), v231)
					v275 = base.Simd_g_v128_load32_lane_l1(m, l0+base.Simd_g_i32x4_extract_lane_l1(v213)<<(uint(v218)%32), v231, v272)
					v278 = base.Simd_g_v128_load32_lane_l2(m, l0+base.Simd_g_i32x4_extract_lane_l2(v213)<<(uint(v218)%32), v231, v275)
					v281 = base.Simd_g_v128_load32_lane_l3(m, l0+base.Simd_g_i32x4_extract_lane_l3(v213)<<(uint(v218)%32), v231, v278)
					v285 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v246, v247), v249), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v281, v247), v249))
					v288 = int32(7)
					v294 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v246, v288), v249), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v281, v288), v249))
					v304 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v246, v212), v249), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v281, v212), v249))
					v305 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k3)
					v308 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k4)
					v310 = int32(18)
					v312 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k5)
					v314 = int32(4)
					v316 = int32(5)
					v318 = int32(6)
					v321 = base.Simd_g_i16x8_replace_lane_l7(base.Simd_g_i16x8_replace_lane_l6(base.Simd_g_i16x8_replace_lane_l5(base.Simd_g_i16x8_replace_lane_l4(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v285, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k6)), base.Simd_g_i32x4_mul(v294, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k7))), base.Simd_g_i32x4_mul(v304, v305)), v308), v310), v304, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k5), base.Simd_g_const(&F_WebPConvertARGBToUV_C__k8)), l0), l0), l0), l0)
					v323 = base.Simd_g_v128_load32_zero(m, v201, v231)
					v324 = base.Simd_g_i16x8_extend_low_i8x16_u(v323)
					v330 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k9)
					base.Simd_g_v128_store32_lane_l0(m, v201, v231, base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_sub(base.Simd_g_v128_or(v321, v324), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v321, v324), v212)), v304, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k9), base.Simd_g_const(&F_WebPConvertARGBToUV_C__k8)))
					v354 = base.Simd_g_i16x8_replace_lane_l7(base.Simd_g_i16x8_replace_lane_l6(base.Simd_g_i16x8_replace_lane_l5(base.Simd_g_i16x8_replace_lane_l4(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v285, v305), base.Simd_g_i32x4_mul(v294, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k10))), base.Simd_g_i32x4_mul(v304, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k11))), v308), v310), v304, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k5), base.Simd_g_const(&F_WebPConvertARGBToUV_C__k8)), l0), l0), l0), l0)
					v356 = base.Simd_g_v128_load32_zero(m, v200, v231)
					v357 = base.Simd_g_i16x8_extend_low_i8x16_u(v356)
					base.Simd_g_v128_store32_lane_l0(m, v200, v231, base.Simd_g_i8x16_shuffle2(base.Simd_g_i16x8_sub(base.Simd_g_v128_or(v354, v357), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v354, v357), v212)), v354, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k9), base.Simd_g_const(&F_WebPConvertARGBToUV_C__k8)))
					v375 = v199 + int32(-4)
					if v375 != 0 {
						v198 = base.Simd_g_i32x4_add(v198, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k12))
						v199 = v375
						v200 = v200 + v314
						v201 = v201 + v314
						continue
					} else {
						break
					}
					break
				}
				if v23 == v190 {
					v612 = v23
				} else {
					v383 = v190
					v407 = v383
					v410 = l0 + v383<<(uint(int32(3))%32)
					for {
						v422 = l1 + v407
						v425 = *(*int32)(unsafe.Add(mBase, uint32(v410+int32(4))))
						v426 = int32(15)
						v428 = int32(510)
						v430 = *(*int32)(unsafe.Add(mBase, uint32(v410)))
						v435 = int32(base.Ui32(v425)>>(uint(v426)%32))&v428 + int32(base.Ui32(v430)>>(uint(v426)%32))&v428
						v438 = int32(7)
						v446 = int32(base.Ui32(v425)>>(uint(v438)%32))&v428 + int32(base.Ui32(v430)>>(uint(v438)%32))&v428
						v450 = int32(1)
						v458 = v425<<(uint(v450)%32)&v428 + v430<<(uint(v450)%32)&v428
						v459 = int32(_a_F_WebPConvertARGBToUV_C_0)
						v462 = int32(33685504)
						v464 = int32(18)
						v466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v422))))
						v471 = int32(base.Ui32(int32(base.Ui32(v435*int32(-9719)+v446*int32(-19081)+v458*v459+v462)>>(uint(v464)%32))+v466+v450) >> (uint(v450) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v422))) = uint8(v471)
						v473 = l2 + v407
						v486 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v473))))
						v491 = int32(base.Ui32(int32(base.Ui32(v435*v459+v446*int32(-24116)+v458*int32(-4684)+v462)>>(uint(v464)%32))+v486+v450) >> (uint(v450) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v473))) = uint8(v491)
						v496 = v407 + v450
						if v23 != v496 {
							v407 = v496
							v410 = v410 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v612 = v23
				}
			}
		} else {
			if base.Ui32(v23) <= base.Ui32(int32(3)) {
				v504 = int32(0)
				v528 = v504
				v531 = l0 + v504<<(uint(int32(3))%32)
				for {
					v546 = *(*int32)(unsafe.Add(mBase, uint32(v531+int32(4))))
					v547 = int32(15)
					v549 = int32(510)
					v551 = *(*int32)(unsafe.Add(mBase, uint32(v531)))
					v556 = int32(base.Ui32(v546)>>(uint(v547)%32))&v549 + int32(base.Ui32(v551)>>(uint(v547)%32))&v549
					v559 = int32(7)
					v567 = int32(base.Ui32(v546)>>(uint(v559)%32))&v549 + int32(base.Ui32(v551)>>(uint(v559)%32))&v549
					v571 = int32(1)
					v579 = v546<<(uint(v571)%32)&v549 + v551<<(uint(v571)%32)&v549
					v580 = int32(_a_F_WebPConvertARGBToUV_C_0)
					v583 = int32(33685504)
					v585 = int32(18)
					v586 = int32(base.Ui32(v556*int32(67099145)+v567*int32(67089783)+v579*v580+v583) >> (uint(v585) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(l1+v528))) = uint8(v586)
					v600 = int32(base.Ui32(v556*v580+v567*int32(67084748)+v579*int32(67104180)+v583) >> (uint(v585) % 32))
					*(*uint8)(unsafe.Add(mBase, uint32(l2+v528))) = uint8(v600)
					v605 = v528 + v571
					if v23 != v605 {
						v528 = v605
						v531 = v531 + int32(8)
						continue
					} else {
						break
					}
					break
				}
				v612 = v23
			} else {
				v34 = v23 & int32(2147483644)
				v42 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k0)
				v43 = v34
				v44 = l2
				v45 = l1
				for {
					v56 = int32(1)
					v57 = base.Simd_g_i32x4_shl(v42, v56)
					v59 = base.Simd_g_v128_or(v57, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k1))
					v60 = int32(3)
					v62 = int32(2)
					v75 = int32(0)
					v81 = base.Simd_g_v128_load32_splat(m, l0+base.Simd_g_i32x4_extract_lane_l0(v59)<<(uint(v62)%32), v75)
					v84 = base.Simd_g_v128_load32_lane_l1(m, l0+base.Simd_g_i32x4_extract_lane_l1(v59)<<(uint(v62)%32), v75, v81)
					v87 = base.Simd_g_v128_load32_lane_l2(m, l0+base.Simd_g_i32x4_extract_lane_l2(v59)<<(uint(v62)%32), v75, v84)
					v90 = base.Simd_g_v128_load32_lane_l3(m, l0+base.Simd_g_i32x4_extract_lane_l3(v59)<<(uint(v62)%32), v75, v87)
					v91 = int32(15)
					v93 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k2)
					v116 = base.Simd_g_v128_load32_splat(m, l0+base.Simd_g_i32x4_extract_lane_l0(v57)<<(uint(v62)%32), v75)
					v119 = base.Simd_g_v128_load32_lane_l1(m, l0+base.Simd_g_i32x4_extract_lane_l1(v57)<<(uint(v62)%32), v75, v116)
					v122 = base.Simd_g_v128_load32_lane_l2(m, l0+base.Simd_g_i32x4_extract_lane_l2(v57)<<(uint(v62)%32), v75, v119)
					v125 = base.Simd_g_v128_load32_lane_l3(m, l0+base.Simd_g_i32x4_extract_lane_l3(v57)<<(uint(v62)%32), v75, v122)
					v129 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v90, v91), v93), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v125, v91), v93))
					v132 = int32(7)
					v138 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v90, v132), v93), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v125, v132), v93))
					v148 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v90, v56), v93), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(v125, v56), v93))
					v149 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k3)
					v152 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k4)
					v154 = int32(18)
					v156 = base.Simd_g_const(&F_WebPConvertARGBToUV_C__k13)
					base.Simd_g_v128_store32_lane_l0(m, v45, v75, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v129, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k14)), base.Simd_g_i32x4_mul(v138, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k15))), base.Simd_g_i32x4_mul(v148, v149)), v152), v154), v148, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k13), base.Simd_g_const(&F_WebPConvertARGBToUV_C__k8)))
					base.Simd_g_v128_store32_lane_l0(m, v44, v75, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v129, v149), base.Simd_g_i32x4_mul(v138, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k16))), base.Simd_g_i32x4_mul(v148, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k17))), v152), v154), v148, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k13), base.Simd_g_const(&F_WebPConvertARGBToUV_C__k8)))
					v176 = int32(4)
					v183 = v43 + int32(-4)
					if v183 != 0 {
						v42 = base.Simd_g_i32x4_add(v42, base.Simd_g_const(&F_WebPConvertARGBToUV_C__k12))
						v43 = v183
						v44 = v44 + v176
						v45 = v45 + v176
						continue
					} else {
						break
					}
					break
				}
				if v23 != v34 {
					v504 = v34
					v528 = v504
					v531 = l0 + v504<<(uint(int32(3))%32)
					for {
						v546 = *(*int32)(unsafe.Add(mBase, uint32(v531+int32(4))))
						v547 = int32(15)
						v549 = int32(510)
						v551 = *(*int32)(unsafe.Add(mBase, uint32(v531)))
						v556 = int32(base.Ui32(v546)>>(uint(v547)%32))&v549 + int32(base.Ui32(v551)>>(uint(v547)%32))&v549
						v559 = int32(7)
						v567 = int32(base.Ui32(v546)>>(uint(v559)%32))&v549 + int32(base.Ui32(v551)>>(uint(v559)%32))&v549
						v571 = int32(1)
						v579 = v546<<(uint(v571)%32)&v549 + v551<<(uint(v571)%32)&v549
						v580 = int32(_a_F_WebPConvertARGBToUV_C_0)
						v583 = int32(33685504)
						v585 = int32(18)
						v586 = int32(base.Ui32(v556*int32(67099145)+v567*int32(67089783)+v579*v580+v583) >> (uint(v585) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l1+v528))) = uint8(v586)
						v600 = int32(base.Ui32(v556*v580+v567*int32(67084748)+v579*int32(67104180)+v583) >> (uint(v585) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(l2+v528))) = uint8(v600)
						v605 = v528 + v571
						if v23 != v605 {
							v528 = v605
							v531 = v531 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v612 = v23
				} else {
					v612 = v23
				}
			}
		}
	} else {
		v612 = int32(0)
	}
	if l3&int32(1) == int32(0) {
		return
	} else {
		v635 = *(*int32)(unsafe.Add(mBase, uint32(l0+v612<<(uint(int32(3))%32))))
		v638 = int32(1020)
		v639 = int32(base.Ui32(v635)>>(uint(int32(14))%32)) & v638
		v640 = int32(_a_F_WebPConvertARGBToUV_C_0)
		v645 = int32(base.Ui32(v635)>>(uint(int32(6))%32)) & v638
		v652 = v635 << (uint(int32(2)) % 32) & v638
		v656 = int32(33685504)
		v658 = int32(18)
		v659 = int32(base.Ui32(v639*v640+v645*int32(-24116)+v652*int32(-4684)+v656) >> (uint(v658) % 32))
		v671 = int32(base.Ui32(v639*int32(-9719)+v645*int32(-19081)+v652*v640+v656) >> (uint(v658) % 32))
		if l4 == int32(0) {
			v678 = l1 + v612
			v679 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v678))))
			v681 = int32(1)
			v684 = int32(base.Ui32(v671+v679+v681) >> (uint(v681) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v678))) = uint8(v684)
			v686 = l2 + v612
			v687 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v686))))
			v692 = int32(base.Ui32(v659+v687+v681) >> (uint(v681) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v686))) = uint8(v692)
			return
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l2+v612))) = uint8(v659)
			*(*uint8)(unsafe.Add(mBase, uint32(l1+v612))) = uint8(v671)
			return
		}
	}
}

var F_WebPConvertARGBToUV_C__k0 = [2]uint64{0x100000000, 0x300000002}
var F_WebPConvertARGBToUV_C__k1 = [2]uint64{0x100000001, 0x100000001}
var F_WebPConvertARGBToUV_C__k2 = [2]uint64{0x1fe000001fe, 0x1fe000001fe}
var F_WebPConvertARGBToUV_C__k3 = [2]uint64{0x708000007080, 0x708000007080}
var F_WebPConvertARGBToUV_C__k4 = [2]uint64{0x202000002020000, 0x202000002020000}
var F_WebPConvertARGBToUV_C__k5 = [2]uint64{0xd0c090805040100, 0x100010001000100}
var F_WebPConvertARGBToUV_C__k6 = [2]uint64{0xffffda09ffffda09, 0xffffda09ffffda09}
var F_WebPConvertARGBToUV_C__k7 = [2]uint64{0xffffb577ffffb577, 0xffffb577ffffb577}
var F_WebPConvertARGBToUV_C__k8 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_WebPConvertARGBToUV_C__k9 = [2]uint64{0x6040200, 0x0}
var F_WebPConvertARGBToUV_C__k10 = [2]uint64{0xffffa1ccffffa1cc, 0xffffa1ccffffa1cc}
var F_WebPConvertARGBToUV_C__k11 = [2]uint64{0xffffedb4ffffedb4, 0xffffedb4ffffedb4}
var F_WebPConvertARGBToUV_C__k12 = [2]uint64{0x400000004, 0x400000004}
var F_WebPConvertARGBToUV_C__k13 = [2]uint64{0xc080400, 0x0}
var F_WebPConvertARGBToUV_C__k14 = [2]uint64{0x3ffda0903ffda09, 0x3ffda0903ffda09}
var F_WebPConvertARGBToUV_C__k15 = [2]uint64{0x3ffb57703ffb577, 0x3ffb57703ffb577}
var F_WebPConvertARGBToUV_C__k16 = [2]uint64{0x3ffa1cc03ffa1cc, 0x3ffa1cc03ffa1cc}
var F_WebPConvertARGBToUV_C__k17 = [2]uint64{0x3ffedb403ffedb4, 0x3ffedb403ffedb4}

func F_WebPConvertRGBA32ToUV_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v58 base.V128
	_ = v58
	var v60 int32
	_ = v60
	var v61 base.V128
	_ = v61
	var v63 int32
	_ = v63
	var v64 base.V128
	_ = v64
	var v66 int32
	_ = v66
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v73 int32
	_ = v73
	var v80 base.V128
	_ = v80
	var v83 base.V128
	_ = v83
	var v86 base.V128
	_ = v86
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v91 base.V128
	_ = v91
	var v103 base.V128
	_ = v103
	var v106 base.V128
	_ = v106
	var v109 base.V128
	_ = v109
	var v112 base.V128
	_ = v112
	var v113 base.V128
	_ = v113
	var v114 base.V128
	_ = v114
	var v117 base.V128
	_ = v117
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v134 base.V128
	_ = v134
	var v137 base.V128
	_ = v137
	var v140 base.V128
	_ = v140
	var v143 base.V128
	_ = v143
	var v144 base.V128
	_ = v144
	var v155 base.V128
	_ = v155
	var v158 base.V128
	_ = v158
	var v161 base.V128
	_ = v161
	var v164 base.V128
	_ = v164
	var v165 base.V128
	_ = v165
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
	var v195 base.V128
	_ = v195
	var v200 base.V128
	_ = v200
	var v202 base.V128
	_ = v202
	var v234 int32
	_ = v234
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v266 int32
	_ = v266
	var v267 int32
	_ = v267
	var v279 int32
	_ = v279
	var v282 int32
	_ = v282
	var v286 int32
	_ = v286
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v297 int32
	_ = v297
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v318 int32
	_ = v318
	var v319 int32
	_ = v319
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v331 int32
	_ = v331
	if l3 < int32(1) {
	} else {
		if base.Ui32(int32(8)) <= base.Ui32(l3) {
			v27 = l3 & int32(2147483640)
			v31 = l0
			v37 = v27
			v38 = l2
			v39 = l1
			for {
				v55 = int32(8)
				v57 = int32(0)
				v58 = base.Simd_g_v128_load16_splat(m, v31, v57)
				v60 = int32(1)
				v61 = base.Simd_g_v128_load16_lane_l1(m, v31+v55, v57, v58)
				v63 = int32(2)
				v64 = base.Simd_g_v128_load16_lane_l2(m, v31+int32(16), v57, v61)
				v66 = int32(3)
				v67 = base.Simd_g_v128_load16_lane_l3(m, v31+int32(24), v57, v64)
				v68 = base.Simd_g_i32x4_extend_low_i16x8_u(v67)
				v69 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k0)
				v73 = int32(18)
				v80 = base.Simd_g_v128_load16_splat(m, v31+v63, v57)
				v83 = base.Simd_g_v128_load16_lane_l1(m, v31+int32(10), v57, v80)
				v86 = base.Simd_g_v128_load16_lane_l2(m, v31+v73, v57, v83)
				v89 = base.Simd_g_v128_load16_lane_l3(m, v31+int32(26), v57, v86)
				v90 = base.Simd_g_i32x4_extend_low_i16x8_u(v89)
				v91 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k1)
				v103 = base.Simd_g_v128_load16_splat(m, v31+int32(4), v57)
				v106 = base.Simd_g_v128_load16_lane_l1(m, v31+int32(12), v57, v103)
				v109 = base.Simd_g_v128_load16_lane_l2(m, v31+int32(20), v57, v106)
				v112 = base.Simd_g_v128_load16_lane_l3(m, v31+int32(28), v57, v109)
				v113 = base.Simd_g_i32x4_extend_low_i16x8_u(v112)
				v114 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k2)
				v117 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k3)
				v121 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k4)
				v123 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k5)
				v134 = base.Simd_g_v128_load16_splat(m, v31+int32(32), v57)
				v137 = base.Simd_g_v128_load16_lane_l1(m, v31+int32(40), v57, v134)
				v140 = base.Simd_g_v128_load16_lane_l2(m, v31+int32(48), v57, v137)
				v143 = base.Simd_g_v128_load16_lane_l3(m, v31+int32(56), v57, v140)
				v144 = base.Simd_g_i32x4_extend_low_i16x8_u(v143)
				v155 = base.Simd_g_v128_load16_splat(m, v31+int32(34), v57)
				v158 = base.Simd_g_v128_load16_lane_l1(m, v31+int32(42), v57, v155)
				v161 = base.Simd_g_v128_load16_lane_l2(m, v31+int32(50), v57, v158)
				v164 = base.Simd_g_v128_load16_lane_l3(m, v31+int32(58), v57, v161)
				v165 = base.Simd_g_i32x4_extend_low_i16x8_u(v164)
				v177 = base.Simd_g_v128_load16_splat(m, v31+int32(36), v57)
				v180 = base.Simd_g_v128_load16_lane_l1(m, v31+int32(44), v57, v177)
				v183 = base.Simd_g_v128_load16_lane_l2(m, v31+int32(52), v57, v180)
				v186 = base.Simd_g_v128_load16_lane_l3(m, v31+int32(60), v57, v183)
				v187 = base.Simd_g_i32x4_extend_low_i16x8_u(v186)
				v195 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k6)
				base.Simd_g_v128_store64_lane_l0(m, v38, v57, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v68, v69), base.Simd_g_i32x4_mul(v90, v91)), base.Simd_g_i32x4_mul(v113, v114)), v117), v73), v121), v123), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v144, v69), base.Simd_g_i32x4_mul(v165, v91)), base.Simd_g_i32x4_mul(v187, v114)), v117), v73), v121), v123), base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k7), base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k8)))
				v200 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k9)
				v202 = base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k10)
				base.Simd_g_v128_store64_lane_l0(m, v39, v57, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v68, v200), base.Simd_g_i32x4_mul(v90, v202)), base.Simd_g_i32x4_mul(v113, v69)), v117), v73), v121), v123), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v144, v200), base.Simd_g_i32x4_mul(v165, v202)), base.Simd_g_i32x4_mul(v187, v69)), v117), v73), v121), v123), base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k7), base.Simd_g_const(&F_WebPConvertRGBA32ToUV_C__k8)))
				v234 = v37 + int32(-8)
				if v234 != 0 {
					v31 = v31 + int32(64)
					v37 = v234
					v38 = v38 + v55
					v39 = v39 + v55
					continue
				} else {
					break
				}
				break
			}
			if v27 == l3 {
			} else {
				v240 = l0 + v27<<(uint(int32(3))%32)
				v241 = v27
				v259 = l1 + v241
				v263 = v240
				v266 = l2 + v241
				v267 = l3 - v241
				for {
					v279 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v263))))
					v282 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v263)+2)))
					v286 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v263)+4)))
					v293 = (v279*int32(-9719) + v282*int32(-19081) + v286*int32(_a_F_WebPConvertRGBA32ToUV_C_0) + int32(33685504)) >> (uint(int32(18)) % 32)
					v294 = int32(0)
					if v294 < v293 {
						v297 = v293
					} else {
						v297 = v294
					}
					v298 = int32(255)
					if v297 < v298 {
						v301 = v297
					} else {
						v301 = v298
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v259))) = uint8(v301)
					v314 = (v279*int32(_a_F_WebPConvertRGBA32ToUV_C_0) + v282*int32(-24116) + v286*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
					v315 = int32(0)
					if v315 < v314 {
						v318 = v314
					} else {
						v318 = v315
					}
					v319 = int32(255)
					if v318 < v319 {
						v322 = v318
					} else {
						v322 = v319
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v266))) = uint8(v322)
					v324 = int32(1)
					v331 = v267 + int32(-1)
					if v331 != 0 {
						v259 = v259 + v324
						v263 = v263 + int32(8)
						v266 = v266 + v324
						v267 = v331
						continue
					} else {
						break
					}
					break
				}
			}
		} else {
			v240 = l0
			v241 = int32(0)
			v259 = l1 + v241
			v263 = v240
			v266 = l2 + v241
			v267 = l3 - v241
			for {
				v279 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v263))))
				v282 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v263)+2)))
				v286 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v263)+4)))
				v293 = (v279*int32(-9719) + v282*int32(-19081) + v286*int32(_a_F_WebPConvertRGBA32ToUV_C_0) + int32(33685504)) >> (uint(int32(18)) % 32)
				v294 = int32(0)
				if v294 < v293 {
					v297 = v293
				} else {
					v297 = v294
				}
				v298 = int32(255)
				if v297 < v298 {
					v301 = v297
				} else {
					v301 = v298
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v259))) = uint8(v301)
				v314 = (v279*int32(_a_F_WebPConvertRGBA32ToUV_C_0) + v282*int32(-24116) + v286*int32(-4684) + int32(33685504)) >> (uint(int32(18)) % 32)
				v315 = int32(0)
				if v315 < v314 {
					v318 = v314
				} else {
					v318 = v315
				}
				v319 = int32(255)
				if v318 < v319 {
					v322 = v318
				} else {
					v322 = v319
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v266))) = uint8(v322)
				v324 = int32(1)
				v331 = v267 + int32(-1)
				if v331 != 0 {
					v259 = v259 + v324
					v263 = v263 + int32(8)
					v266 = v266 + v324
					v267 = v331
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

var F_WebPConvertRGBA32ToUV_C__k0 = [2]uint64{0x708000007080, 0x708000007080}
var F_WebPConvertRGBA32ToUV_C__k1 = [2]uint64{0xffffa1ccffffa1cc, 0xffffa1ccffffa1cc}
var F_WebPConvertRGBA32ToUV_C__k2 = [2]uint64{0xffffedb4ffffedb4, 0xffffedb4ffffedb4}
var F_WebPConvertRGBA32ToUV_C__k3 = [2]uint64{0x202000002020000, 0x202000002020000}
var F_WebPConvertRGBA32ToUV_C__k4 = [2]uint64{0x0, 0x0}
var F_WebPConvertRGBA32ToUV_C__k5 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_WebPConvertRGBA32ToUV_C__k6 = [2]uint64{0x1c1814100c080400, 0x0}
var F_WebPConvertRGBA32ToUV_C__k7 = [2]uint64{0x808080800c080400, 0x0}
var F_WebPConvertRGBA32ToUV_C__k8 = [2]uint64{0xc08040080808080, 0x8080808080808080}
var F_WebPConvertRGBA32ToUV_C__k9 = [2]uint64{0xffffda09ffffda09, 0xffffda09ffffda09}
var F_WebPConvertRGBA32ToUV_C__k10 = [2]uint64{0xffffb577ffffb577, 0xffffb577ffffb577}

func F_WebPMultARGBRow_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v28 base.V128
	_ = v28
	var v30 int32
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v35 base.V128
	_ = v35
	var v47 int32
	_ = v47
	var v49 base.V128
	_ = v49
	var v50 int32
	_ = v50
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v64 base.V128
	_ = v64
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v79 base.V128
	_ = v79
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v94 base.V128
	_ = v94
	var v96 base.V128
	_ = v96
	var v99 base.V128
	_ = v99
	var v105 int32
	_ = v105
	var v110 int32
	_ = v110
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v129 int32
	_ = v129
	var v135 int32
	_ = v135
	var v138 int32
	_ = v138
	var v147 int32
	_ = v147
	var v158 int32
	_ = v158
	var v169 int32
	_ = v169
	var v174 int32
	_ = v174
	var v180 int32
	_ = v180
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v192 int32
	_ = v192
	var v193 base.V128
	_ = v193
	var v194 int32
	_ = v194
	var v197 base.V128
	_ = v197
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v205 base.V128
	_ = v205
	var v208 int32
	_ = v208
	var v213 int32
	_ = v213
	var v230 base.V128
	_ = v230
	var v231 base.V128
	_ = v231
	var v238 int32
	_ = v238
	var v241 int32
	_ = v241
	var v250 int32
	_ = v250
	var v261 int32
	_ = v261
	var v272 int32
	_ = v272
	var v277 int32
	_ = v277
	var v282 int32
	_ = v282
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v303 int32
	_ = v303
	var v309 int32
	_ = v309
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v316 int32
	_ = v316
	var v323 int32
	_ = v323
	var v330 int32
	_ = v330
	var v347 int32
	_ = v347
	var v353 int32
	_ = v353
	var v357 int32
	_ = v357
	var v369 int32
	_ = v369
	var v370 int32
	_ = v370
	var v378 int32
	_ = v378
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v391 int32
	_ = v391
	var v398 int32
	_ = v398
	var v405 int32
	_ = v405
	var v422 int32
	_ = v422
	var v428 int32
	_ = v428
	if l1 < int32(1) {
	} else {
		if l2 == int32(0) {
			if base.Ui32(l1) < base.Ui32(int32(4)) {
				v282 = int32(0)
				v294 = l1 - v282
				v295 = l0 + v282<<(uint(int32(2))%32)
				for {
					v303 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
					if base.Ui32(int32(-16777217)) < base.Ui32(v303) {
					} else {
						if base.Ui32(v303) < base.Ui32(int32(16777216)) {
							v347 = int32(0)
						} else {
							v309 = int32(24)
							v312 = int32(base.Ui32(v303)>>(uint(v309)%32)) * int32(_a_F_WebPMultARGBRow_C_0)
							v313 = int32(255)
							v316 = int32(8388608)
							v323 = int32(8)
							v330 = int32(16)
							v347 = int32(base.Ui32(v312*(v303&v313)+v316)>>(uint(v309)%32)) | v303&int32(-16777216) | int32(base.Ui32(v312*(int32(base.Ui32(v303)>>(uint(v323)%32))&v313)+v316)>>(uint(v330)%32))&int32(_a_F_WebPMultARGBRow_C_1) | int32(base.Ui32(v312*(int32(base.Ui32(v303)>>(uint(v330)%32))&v313)+v316)>>(uint(v323)%32))&int32(16711680)
						}
						*(*int32)(unsafe.Add(mBase, uint32(v295))) = v347
					}
					v353 = v294 + int32(-1)
					if v353 != 0 {
						v294 = v353
						v295 = v295 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			} else {
				v180 = l1 & int32(2147483644)
				v183 = l0
				v185 = v180
				for {
					v192 = int32(0)
					v193 = base.Simd_g_v128_load(m, v183, v192)
					v194 = int32(24)
					v197 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_shr_u(v193, v194), base.Simd_g_const(&F_WebPMultARGBRow_C__k0))
					v198 = base.Simd_g_const(&F_WebPMultARGBRow_C__k1)
					v201 = base.Simd_g_const(&F_WebPMultARGBRow_C__k2)
					v205 = base.Simd_g_const(&F_WebPMultARGBRow_C__k3)
					v208 = int32(8)
					v213 = int32(16)
					v230 = base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPMultARGBRow_C__k4), base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v197, base.Simd_g_v128_and(v193, v198)), v201), v194), base.Simd_g_v128_and(v193, v205)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v197, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v193, v208), v198)), v201), v213), base.Simd_g_const(&F_WebPMultARGBRow_C__k5))), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v197, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v193, v213), v198)), v201), v208), base.Simd_g_const(&F_WebPMultARGBRow_C__k6))), base.Simd_g_i32x4_lt_u(v193, base.Simd_g_const(&F_WebPMultARGBRow_C__k7)))
					v231 = base.Simd_g_i32x4_lt_u(v193, v205)
					if base.Simd_g_i32x4_extract_lane_l0(v231)&int32(1) == v192 {
					} else {
						v238 = int32(0)
						base.Simd_g_v128_store32_lane_l0(m, v183, v238, v230)
					}
					v241 = int32(1)
					if base.Simd_g_i32x4_extract_lane_l1(v231)&v241 == int32(0) {
					} else {
						v250 = int32(1)
						base.Simd_g_v128_store32_lane_l1(m, v183+int32(4), int32(0), v230)
					}
					if base.Simd_g_i32x4_extract_lane_l2(v231)&int32(1) == int32(0) {
					} else {
						v261 = int32(2)
						base.Simd_g_v128_store32_lane_l2(m, v183+int32(8), int32(0), v230)
					}
					if base.Simd_g_i32x4_extract_lane_l3(v231)&int32(1) == int32(0) {
					} else {
						v272 = int32(3)
						base.Simd_g_v128_store32_lane_l3(m, v183+int32(12), int32(0), v230)
					}
					v277 = v185 + int32(-4)
					if v277 != 0 {
						v183 = v183 + int32(16)
						v185 = v277
						continue
					} else {
						break
					}
					break
				}
				if v180 == l1 {
				} else {
					v282 = v180
					v294 = l1 - v282
					v295 = l0 + v282<<(uint(int32(2))%32)
					for {
						v303 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
						if base.Ui32(int32(-16777217)) < base.Ui32(v303) {
						} else {
							if base.Ui32(v303) < base.Ui32(int32(16777216)) {
								v347 = int32(0)
							} else {
								v309 = int32(24)
								v312 = int32(base.Ui32(v303)>>(uint(v309)%32)) * int32(_a_F_WebPMultARGBRow_C_0)
								v313 = int32(255)
								v316 = int32(8388608)
								v323 = int32(8)
								v330 = int32(16)
								v347 = int32(base.Ui32(v312*(v303&v313)+v316)>>(uint(v309)%32)) | v303&int32(-16777216) | int32(base.Ui32(v312*(int32(base.Ui32(v303)>>(uint(v323)%32))&v313)+v316)>>(uint(v330)%32))&int32(_a_F_WebPMultARGBRow_C_1) | int32(base.Ui32(v312*(int32(base.Ui32(v303)>>(uint(v330)%32))&v313)+v316)>>(uint(v323)%32))&int32(16711680)
							}
							*(*int32)(unsafe.Add(mBase, uint32(v295))) = v347
						}
						v353 = v294 + int32(-1)
						if v353 != 0 {
							v294 = v353
							v295 = v295 + int32(4)
							continue
						} else {
							break
						}
						break
					}
				}
			}
		} else {
			if base.Ui32(l1) <= base.Ui32(int32(3)) {
				v357 = int32(0)
				v369 = l1 - v357
				v370 = l0 + v357<<(uint(int32(2))%32)
				for {
					v378 = *(*int32)(unsafe.Add(mBase, uint32(v370)))
					if base.Ui32(int32(-16777217)) < base.Ui32(v378) {
					} else {
						if base.Ui32(v378) < base.Ui32(int32(16777216)) {
							v422 = int32(0)
						} else {
							v384 = int32(-16777216)
							v385 = int32(24)
							v387 = base.I32_div_u_s(v384, int32(base.Ui32(v378)>>(uint(v385)%32)))
							v388 = int32(255)
							v391 = int32(8388608)
							v398 = int32(8)
							v405 = int32(16)
							v422 = int32(base.Ui32(v387*(v378&v388)+v391)>>(uint(v385)%32)) | v378&v384 | int32(base.Ui32(v387*(int32(base.Ui32(v378)>>(uint(v398)%32))&v388)+v391)>>(uint(v405)%32))&int32(_a_F_WebPMultARGBRow_C_1) | int32(base.Ui32(v387*(int32(base.Ui32(v378)>>(uint(v405)%32))&v388)+v391)>>(uint(v398)%32))&int32(16711680)
						}
						*(*int32)(unsafe.Add(mBase, uint32(v370))) = v422
					}
					v428 = v369 + int32(-1)
					if v428 != 0 {
						v369 = v428
						v370 = v370 + int32(4)
						continue
					} else {
						break
					}
					break
				}
			} else {
				v19 = l1 & int32(2147483644)
				v22 = l0
				v24 = v19
				v28 = base.Simd_g_const(&F_WebPMultARGBRow_C__k4)
				for {
					v30 = int32(0)
					v31 = base.Simd_g_v128_load(m, v22, v30)
					v32 = base.Simd_g_const(&F_WebPMultARGBRow_C__k3)
					v35 = base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_add(v31, v32), base.Simd_g_const(&F_WebPMultARGBRow_C__k8))
					if base.Simd_g_i32x4_extract_lane_l0(v35)&int32(1) == v30 {
						v49 = v28
					} else {
						v47 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l0(v31))>>(uint(int32(24))%32)))
						v49 = base.Simd_g_i32x4_splat(v47)
					}
					v50 = int32(1)
					if base.Simd_g_i32x4_extract_lane_l1(v35)&v50 == int32(0) {
						v64 = v49
					} else {
						v57 = int32(1)
						v61 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l1(v31))>>(uint(int32(24))%32)))
						v64 = base.Simd_g_i32x4_replace_lane_l1(v49, v61)
					}
					if base.Simd_g_i32x4_extract_lane_l2(v35)&int32(1) == int32(0) {
						v79 = v64
					} else {
						v72 = int32(2)
						v76 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l2(v31))>>(uint(int32(24))%32)))
						v79 = base.Simd_g_i32x4_replace_lane_l2(v64, v76)
					}
					if base.Simd_g_i32x4_extract_lane_l3(v35)&int32(1) == int32(0) {
						v94 = v79
					} else {
						v87 = int32(3)
						v91 = base.I32_div_u_s(int32(-16777216), int32(base.Ui32(base.Simd_g_i32x4_extract_lane_l3(v31))>>(uint(int32(24))%32)))
						v94 = base.Simd_g_i32x4_replace_lane_l3(v79, v91)
					}
					v96 = base.Simd_g_const(&F_WebPMultARGBRow_C__k1)
					v99 = base.Simd_g_const(&F_WebPMultARGBRow_C__k2)
					v105 = int32(8)
					v110 = int32(16)
					v127 = base.Simd_g_v128_bitselect(base.Simd_g_const(&F_WebPMultARGBRow_C__k4), base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v94, base.Simd_g_v128_and(v31, v96)), v99), int32(24)), base.Simd_g_v128_and(v31, v32)), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v94, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v31, v105), v96)), v99), v110), base.Simd_g_const(&F_WebPMultARGBRow_C__k5))), base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v94, base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v31, v110), v96)), v99), v105), base.Simd_g_const(&F_WebPMultARGBRow_C__k6))), base.Simd_g_i32x4_lt_u(v31, base.Simd_g_const(&F_WebPMultARGBRow_C__k7)))
					v128 = base.Simd_g_i32x4_lt_u(v31, v32)
					v129 = int32(0)
					if base.Simd_g_i32x4_extract_lane_l0(v128)&int32(1) == v129 {
					} else {
						v135 = int32(0)
						base.Simd_g_v128_store32_lane_l0(m, v22, v135, v127)
					}
					v138 = int32(1)
					if base.Simd_g_i32x4_extract_lane_l1(v128)&v138 == int32(0) {
					} else {
						v147 = int32(1)
						base.Simd_g_v128_store32_lane_l1(m, v22+int32(4), int32(0), v127)
					}
					if base.Simd_g_i32x4_extract_lane_l2(v128)&int32(1) == int32(0) {
					} else {
						v158 = int32(2)
						base.Simd_g_v128_store32_lane_l2(m, v22+int32(8), int32(0), v127)
					}
					if base.Simd_g_i32x4_extract_lane_l3(v128)&int32(1) == int32(0) {
					} else {
						v169 = int32(3)
						base.Simd_g_v128_store32_lane_l3(m, v22+int32(12), int32(0), v127)
					}
					v174 = v24 + int32(-4)
					if v174 != 0 {
						v22 = v22 + int32(16)
						v24 = v174
						v28 = v127
						continue
					} else {
						break
					}
					break
				}
				if v19 != l1 {
					v357 = v19
					v369 = l1 - v357
					v370 = l0 + v357<<(uint(int32(2))%32)
					for {
						v378 = *(*int32)(unsafe.Add(mBase, uint32(v370)))
						if base.Ui32(int32(-16777217)) < base.Ui32(v378) {
						} else {
							if base.Ui32(v378) < base.Ui32(int32(16777216)) {
								v422 = int32(0)
							} else {
								v384 = int32(-16777216)
								v385 = int32(24)
								v387 = base.I32_div_u_s(v384, int32(base.Ui32(v378)>>(uint(v385)%32)))
								v388 = int32(255)
								v391 = int32(8388608)
								v398 = int32(8)
								v405 = int32(16)
								v422 = int32(base.Ui32(v387*(v378&v388)+v391)>>(uint(v385)%32)) | v378&v384 | int32(base.Ui32(v387*(int32(base.Ui32(v378)>>(uint(v398)%32))&v388)+v391)>>(uint(v405)%32))&int32(_a_F_WebPMultARGBRow_C_1) | int32(base.Ui32(v387*(int32(base.Ui32(v378)>>(uint(v405)%32))&v388)+v391)>>(uint(v398)%32))&int32(16711680)
							}
							*(*int32)(unsafe.Add(mBase, uint32(v370))) = v422
						}
						v428 = v369 + int32(-1)
						if v428 != 0 {
							v369 = v428
							v370 = v370 + int32(4)
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
	return
}

var F_WebPMultARGBRow_C__k0 = [2]uint64{0x1010100010101, 0x1010100010101}
var F_WebPMultARGBRow_C__k1 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_WebPMultARGBRow_C__k2 = [2]uint64{0x80000000800000, 0x80000000800000}
var F_WebPMultARGBRow_C__k3 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}
var F_WebPMultARGBRow_C__k4 = [2]uint64{0x0, 0x0}
var F_WebPMultARGBRow_C__k5 = [2]uint64{0xff000000ff00, 0xff000000ff00}
var F_WebPMultARGBRow_C__k6 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_WebPMultARGBRow_C__k7 = [2]uint64{0x100000001000000, 0x100000001000000}
var F_WebPMultARGBRow_C__k8 = [2]uint64{0xfe000000fe000000, 0xfe000000fe000000}

func F_WriteImage(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v31 int64
	_ = v31
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v108 int32
	_ = v108
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v135 base.V128
	_ = v135
	var v140 int64
	_ = v140
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v149 int32
	_ = v149
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
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
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v188 int32
	_ = v188
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v194 int32
	_ = v194
	var v201 int32
	_ = v201
	v9 = m.G0
	v11 = v9 - int32(32)
	m.G0 = v11
	v21 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v23 = v21 - v22
	v25 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v31 = base.I64_extend_i32_u(v23) + base.I64_extend_i32_u((v25+int32(7))>>(uint(int32(3))%32))
	if base.Ui64(v31) < base.Ui64(int64(4294967296)) {
		v37 = base.I32_wrap_i64(v31)
		v38 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
		v39 = v38 - v22
		if v38 == v22 {
			v46 = int32(base.Ui32(v39*int32(3)) >> (uint(int32(1)) % 32))
			if base.Ui32(v37) < base.Ui32(v46) {
				v48 = v46
			} else {
				v48 = v37
			}
			v52 = v48&int32(-1024) + int32(1024)
			v53 = F_WebPSafeMalloc(m, int64(1), v52)
			mBase = m.M
			if v53 != 0 {
				if v21 == v22 {
				} else {
					v58 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v59 = F_memcpy(m, v53, v58, v23)
					mBase = m.M
				}
				v60 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				F_WebPSafeFree(m, v60)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v53 + v52
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v53 + v23
				*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v53
				v67 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
				v68 = v67
				if v68 < int32(1) {
				} else {
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
					v75 = v72
					for {
						v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v82 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v75)
						v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v88 = int32(8)
						v89 = int32(base.Ui32(v87) >> (uint(v88) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l1))) = v89
						v91 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v91 + int32(-8)
						if v88 < v91 {
							v75 = v89
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
				v108 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v118 = v108
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
				v56 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v118 = v56
			}
		} else {
			if base.Ui32(v37) <= base.Ui32(v39) {
				v68 = v25
				if v68 < int32(1) {
				} else {
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
					v75 = v72
					for {
						v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v82 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v75)
						v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v88 = int32(8)
						v89 = int32(base.Ui32(v87) >> (uint(v88) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l1))) = v89
						v91 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v91 + int32(-8)
						if v88 < v91 {
							v75 = v89
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
				v108 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v118 = v108
			} else {
				v46 = int32(base.Ui32(v39*int32(3)) >> (uint(int32(1)) % 32))
				if base.Ui32(v37) < base.Ui32(v46) {
					v48 = v46
				} else {
					v48 = v37
				}
				v52 = v48&int32(-1024) + int32(1024)
				v53 = F_WebPSafeMalloc(m, int64(1), v52)
				mBase = m.M
				if v53 != 0 {
					if v21 == v22 {
					} else {
						v58 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
						v59 = F_memcpy(m, v53, v58, v23)
						mBase = m.M
					}
					v60 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					F_WebPSafeFree(m, v60)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v53 + v52
					*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v53 + v23
					*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v53
					v67 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
					v68 = v67
					if v68 < int32(1) {
					} else {
						v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
						v75 = v72
						for {
							v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v82 + int32(1)
							*(*uint8)(unsafe.Add(mBase, uint32(v82))) = uint8(v75)
							v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
							v88 = int32(8)
							v89 = int32(base.Ui32(v87) >> (uint(v88) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(l1))) = v89
							v91 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
							*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v91 + int32(-8)
							if v88 < v91 {
								v75 = v89
								continue
							} else {
								break
							}
							break
						}
					}
					*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
					v108 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v118 = v108
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
					v56 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v118 = v56
				}
			}
		}
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
		v36 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
		v118 = v36
	}
	v119 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l2))) = v119
	v121 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	if v121 == v119 {
		v128 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
		v129 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
		v130 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
		v131 = m.G1
		v134 = int32(0)
		v135 = base.Simd_g_v128_load(m, v131+int32(_a_F_WriteImage_0), v134)
		base.Simd_g_v128_store(m, v11, v134, v135)
		v140 = *(*int64)(unsafe.Add(mBase, uint32(v131)+uint32(_c_F_WriteImage[0])))
		*(*int64)(unsafe.Add(mBase, uint32(v11)+13)) = v140
		v147 = (v130+int32(7))>>(uint(int32(3))%32) + (v129 - v128)
		v148 = int32(1)
		v149 = v147 + v148
		v151 = v149 & v148
		v152 = v147 + v151
		v154 = v152 + int32(13)
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+4)) = uint8(v154)
		v156 = int32(24)
		v157 = int32(base.Ui32(v154) >> (uint(v156) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+7)) = uint8(v157)
		v159 = int32(16)
		v160 = int32(base.Ui32(v154) >> (uint(v159) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+6)) = uint8(v160)
		v162 = int32(8)
		v163 = int32(base.Ui32(v154) >> (uint(v162) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+5)) = uint8(v163)
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+16)) = uint8(v149)
		v167 = int32(base.Ui32(v149) >> (uint(v156) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+19)) = uint8(v167)
		v170 = int32(base.Ui32(v149) >> (uint(v159) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+18)) = uint8(v170)
		v173 = int32(base.Ui32(v149) >> (uint(v162) % 32))
		*(*uint8)(unsafe.Add(mBase, uint32(v11)+17)) = uint8(v173)
		v176 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
		v177 = m.T0[v176].(func(*base.Module, int32, int32, int32) int32)(m, v11, int32(21), l0)
		mBase = m.M
		if v177 == v134 {
			v183 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v183 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(8)
			}
			v201 = int32(0)
		} else {
			v180 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
			v181 = m.T0[v180].(func(*base.Module, int32, int32, int32) int32)(m, v118, v147, l0)
			mBase = m.M
			if v181 != 0 {
				if v151 == int32(0) {
					*(*int32)(unsafe.Add(mBase, uint32(l2))) = v152 + int32(21)
					v201 = int32(1)
				} else {
					v188 = int32(0)
					*(*uint8)(unsafe.Add(mBase, uint32(v11))) = uint8(v188)
					v191 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
					v192 = m.T0[v191].(func(*base.Module, int32, int32, int32) int32)(m, v11, int32(1), l0)
					mBase = m.M
					if v192 != 0 {
						*(*int32)(unsafe.Add(mBase, uint32(l2))) = v152 + int32(21)
						v201 = int32(1)
					} else {
						v194 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
						if v194 != 0 {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(8)
						}
						v201 = int32(0)
					}
				}
			} else {
				v183 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v183 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(8)
				}
				v201 = int32(0)
			}
		}
	} else {
		v125 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		if v125 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(1)
		}
		v201 = int32(0)
	}
	m.G0 = v11 + int32(32)
	return v201
}

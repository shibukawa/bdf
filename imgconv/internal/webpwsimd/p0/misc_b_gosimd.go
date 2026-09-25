//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_BundleColorMap_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v37 int32
	_ = v37
	var v41 int32
	_ = v41
	var v46 int32
	_ = v46
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v70 base.V128
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v77 base.V128
	_ = v77
	var v78 base.V128
	_ = v78
	var v81 base.V128
	_ = v81
	var v82 base.V128
	_ = v82
	var v85 base.V128
	_ = v85
	var v86 base.V128
	_ = v86
	var v88 base.V128
	_ = v88
	var v94 base.V128
	_ = v94
	var v100 base.V128
	_ = v100
	var v102 base.V128
	_ = v102
	var v104 base.V128
	_ = v104
	var v108 base.V128
	_ = v108
	var v109 int32
	_ = v109
	var v112 base.V128
	_ = v112
	var v114 base.V128
	_ = v114
	var v118 base.V128
	_ = v118
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v144 base.V128
	_ = v144
	var v146 int32
	_ = v146
	var v147 base.V128
	_ = v147
	var v149 base.V128
	_ = v149
	var v150 base.V128
	_ = v150
	var v151 base.V128
	_ = v151
	var v152 base.V128
	_ = v152
	var v155 base.V128
	_ = v155
	var v156 base.V128
	_ = v156
	var v160 base.V128
	_ = v160
	var v162 base.V128
	_ = v162
	var v166 base.V128
	_ = v166
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v193 base.V128
	_ = v193
	var v194 base.V128
	_ = v194
	var v196 base.V128
	_ = v196
	var v197 base.V128
	_ = v197
	var v198 base.V128
	_ = v198
	var v199 base.V128
	_ = v199
	var v202 int32
	_ = v202
	var v204 base.V128
	_ = v204
	var v205 base.V128
	_ = v205
	var v211 base.V128
	_ = v211
	var v213 base.V128
	_ = v213
	var v215 base.V128
	_ = v215
	var v219 base.V128
	_ = v219
	var v223 int32
	_ = v223
	var v225 int32
	_ = v225
	var v227 int32
	_ = v227
	var v233 int32
	_ = v233
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v246 int32
	_ = v246
	var v247 base.V128
	_ = v247
	var v250 base.V128
	_ = v250
	var v251 base.V128
	_ = v251
	var v253 base.V128
	_ = v253
	var v257 base.V128
	_ = v257
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v282 int32
	_ = v282
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v285 base.V128
	_ = v285
	var v286 base.V128
	_ = v286
	var v288 base.V128
	_ = v288
	var v289 base.V128
	_ = v289
	var v290 int32
	_ = v290
	var v293 base.V128
	_ = v293
	var v294 base.V128
	_ = v294
	var v297 int32
	_ = v297
	var v302 base.V128
	_ = v302
	var v304 base.V128
	_ = v304
	var v308 base.V128
	_ = v308
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
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
	var v339 base.V128
	_ = v339
	var v344 base.V128
	_ = v344
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v369 int32
	_ = v369
	var v370 int32
	_ = v370
	var v371 base.V128
	_ = v371
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v377 int32
	_ = v377
	var v383 base.V128
	_ = v383
	var v386 int32
	_ = v386
	var v392 int32
	_ = v392
	var v409 int32
	_ = v409
	var v411 int32
	_ = v411
	var v413 int32
	_ = v413
	var v419 int32
	_ = v419
	var v420 int32
	_ = v420
	var v421 int32
	_ = v421
	var v433 base.V128
	_ = v433
	var v436 int32
	_ = v436
	var v437 int32
	_ = v437
	var v439 int32
	_ = v439
	var v442 int32
	_ = v442
	var v454 int32
	_ = v454
	var v456 int32
	_ = v456
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v482 int32
	_ = v482
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v489 int32
	_ = v489
	var v494 int32
	_ = v494
	var v495 base.V128
	_ = v495
	var v501 base.V128
	_ = v501
	var v509 int32
	_ = v509
	var v514 int32
	_ = v514
	var v515 int32
	_ = v515
	var v518 int32
	_ = v518
	var v519 int32
	_ = v519
	var v530 int32
	_ = v530
	var v535 int32
	_ = v535
	var v544 int32
	_ = v544
	var v545 int32
	_ = v545
	var v546 int32
	_ = v546
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v563 int32
	_ = v563
	var v564 int32
	_ = v564
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v574 int32
	_ = v574
	var v577 int32
	_ = v577
	var v582 int32
	_ = v582
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v601 int32
	_ = v601
	var v608 int32
	_ = v608
	var v622 int32
	_ = v622
	var v625 int32
	_ = v625
	var v626 int32
	_ = v626
	var v631 int32
	_ = v631
	var v642 int32
	_ = v642
	switch l2 {
	case 0:
		v32 = int32(16)
		if v32 <= l1 {
			v37 = l1 + int32(-16)
			if base.Ui32(int32(16)) <= base.Ui32(v37) {
				v60 = l3
				v61 = int32(16)
				v62 = (int32(base.Ui32(v37)>>(uint(int32(4))%32)) + int32(1)) & int32(536870910)
				for {
					v70 = base.Simd_g_const(&F_BundleColorMap_SSE2__k0)
					v71 = l0 + v61
					v72 = int32(0)
					v73 = base.Simd_g_v128_load(m, v71, v72)
					v74 = base.Simd_g_const(&F_BundleColorMap_SSE2__k1)
					v75 = base.Simd_g_i8x16_shuffle2(v70, v73, base.Simd_g_const(&F_BundleColorMap_SSE2__k2), base.Simd_g_const(&F_BundleColorMap_SSE2__k3))
					v76 = base.Simd_g_const(&F_BundleColorMap_SSE2__k4)
					v77 = base.Simd_g_const(&F_BundleColorMap_SSE2__k5)
					v78 = base.Simd_g_i8x16_shuffle2(v75, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
					base.Simd_g_v128_store(m, v60, int32(112), v78)
					v81 = base.Simd_g_const(&F_BundleColorMap_SSE2__k8)
					v82 = base.Simd_g_i8x16_shuffle2(v75, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
					base.Simd_g_v128_store(m, v60, int32(96), v82)
					v85 = base.Simd_g_const(&F_BundleColorMap_SSE2__k11)
					v86 = base.Simd_g_i8x16_shuffle2(v70, v73, base.Simd_g_const(&F_BundleColorMap_SSE2__k12), base.Simd_g_const(&F_BundleColorMap_SSE2__k13))
					v88 = base.Simd_g_i8x16_shuffle2(v86, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
					base.Simd_g_v128_store(m, v60, int32(80), v88)
					v94 = base.Simd_g_i8x16_shuffle2(v86, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
					base.Simd_g_v128_store(m, v60+int32(64), v72, v94)
					v100 = base.Simd_g_v128_load(m, v71+int32(-16), v72)
					v102 = base.Simd_g_i8x16_shuffle2(v70, v100, base.Simd_g_const(&F_BundleColorMap_SSE2__k2), base.Simd_g_const(&F_BundleColorMap_SSE2__k3))
					v104 = base.Simd_g_i8x16_shuffle2(v102, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
					base.Simd_g_v128_store(m, v60, int32(48), v104)
					v108 = base.Simd_g_i8x16_shuffle2(v102, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
					v109 = int32(32)
					base.Simd_g_v128_store(m, v60, v109, v108)
					v112 = base.Simd_g_i8x16_shuffle2(v70, v100, base.Simd_g_const(&F_BundleColorMap_SSE2__k12), base.Simd_g_const(&F_BundleColorMap_SSE2__k13))
					v114 = base.Simd_g_i8x16_shuffle2(v112, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
					base.Simd_g_v128_store(m, v60, int32(16), v114)
					v118 = base.Simd_g_i8x16_shuffle2(v112, v76, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
					base.Simd_g_v128_store(m, v60, v72, v118)
					v122 = v61 + v109
					v124 = v60 + int32(128)
					v126 = v62 + int32(-2)
					if v126 != 0 {
						v60 = v124
						v61 = v122
						v62 = v126
						continue
					} else {
						break
					}
					break
				}
				v132 = v124
				v133 = v122
				v134 = v61 + int32(16)
			} else {
				v132 = l3
				v133 = v32
				v134 = int32(0)
			}
			if v37&int32(16) != 0 {
				v454 = v132
				v456 = v134
			} else {
				v144 = base.Simd_g_const(&F_BundleColorMap_SSE2__k0)
				v146 = int32(0)
				v147 = base.Simd_g_v128_load(m, l0+v134, v146)
				v149 = base.Simd_g_i8x16_shuffle2(v144, v147, base.Simd_g_const(&F_BundleColorMap_SSE2__k2), base.Simd_g_const(&F_BundleColorMap_SSE2__k3))
				v150 = base.Simd_g_const(&F_BundleColorMap_SSE2__k4)
				v151 = base.Simd_g_const(&F_BundleColorMap_SSE2__k5)
				v152 = base.Simd_g_i8x16_shuffle2(v149, v150, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
				base.Simd_g_v128_store(m, v132, int32(48), v152)
				v155 = base.Simd_g_const(&F_BundleColorMap_SSE2__k8)
				v156 = base.Simd_g_i8x16_shuffle2(v149, v150, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
				base.Simd_g_v128_store(m, v132, int32(32), v156)
				v160 = base.Simd_g_i8x16_shuffle2(v144, v147, base.Simd_g_const(&F_BundleColorMap_SSE2__k12), base.Simd_g_const(&F_BundleColorMap_SSE2__k13))
				v162 = base.Simd_g_i8x16_shuffle2(v160, v150, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
				base.Simd_g_v128_store(m, v132, int32(16), v162)
				v166 = base.Simd_g_i8x16_shuffle2(v160, v150, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
				base.Simd_g_v128_store(m, v132, v146, v166)
				v454 = v132 + int32(64)
				v456 = v133
			}
		} else {
			v454 = l3
			v456 = int32(0)
		}
	case 1:
		v23 = int32(16)
		if v23 <= l1 {
			v28 = l1 + int32(-16)
			if base.Ui32(int32(16)) <= base.Ui32(v28) {
				v181 = l3
				v182 = int32(16)
				v183 = (int32(base.Ui32(v28)>>(uint(int32(4))%32)) + int32(1)) & int32(536870910)
				for {
					v191 = l0 + v182
					v192 = int32(0)
					v193 = base.Simd_g_v128_load(m, v191, v192)
					v194 = base.Simd_g_const(&F_BundleColorMap_SSE2__k14)
					v196 = base.Simd_g_const(&F_BundleColorMap_SSE2__k4)
					v197 = base.Simd_g_v128_and(base.Simd_g_i16x8_mul(v193, v194), v196)
					v198 = base.Simd_g_const(&F_BundleColorMap_SSE2__k5)
					v199 = base.Simd_g_i8x16_shuffle2(v197, v196, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
					base.Simd_g_v128_store(m, v181, int32(48), v199)
					v202 = int32(32)
					v204 = base.Simd_g_const(&F_BundleColorMap_SSE2__k8)
					v205 = base.Simd_g_i8x16_shuffle2(v197, v196, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
					base.Simd_g_v128_store(m, v181+v202, v192, v205)
					v211 = base.Simd_g_v128_load(m, v191+int32(-16), v192)
					v213 = base.Simd_g_v128_and(base.Simd_g_i16x8_mul(v211, v194), v196)
					v215 = base.Simd_g_i8x16_shuffle2(v213, v196, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
					base.Simd_g_v128_store(m, v181, int32(16), v215)
					v219 = base.Simd_g_i8x16_shuffle2(v213, v196, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
					base.Simd_g_v128_store(m, v181, v192, v219)
					v223 = v182 + v202
					v225 = v181 + int32(64)
					v227 = v183 + int32(-2)
					if v227 != 0 {
						v181 = v225
						v182 = v223
						v183 = v227
						continue
					} else {
						break
					}
					break
				}
				v233 = v225
				v234 = v223
				v235 = v182 + int32(16)
			} else {
				v233 = l3
				v234 = v23
				v235 = int32(0)
			}
			if v28&int32(16) != 0 {
				v454 = v233
				v456 = v235
			} else {
				v246 = int32(0)
				v247 = base.Simd_g_v128_load(m, l0+v235, v246)
				v250 = base.Simd_g_const(&F_BundleColorMap_SSE2__k4)
				v251 = base.Simd_g_v128_and(base.Simd_g_i16x8_mul(v247, base.Simd_g_const(&F_BundleColorMap_SSE2__k14)), v250)
				v253 = base.Simd_g_i8x16_shuffle2(v251, v250, base.Simd_g_const(&F_BundleColorMap_SSE2__k6), base.Simd_g_const(&F_BundleColorMap_SSE2__k7))
				base.Simd_g_v128_store(m, v233, int32(16), v253)
				v257 = base.Simd_g_i8x16_shuffle2(v251, v250, base.Simd_g_const(&F_BundleColorMap_SSE2__k9), base.Simd_g_const(&F_BundleColorMap_SSE2__k10))
				base.Simd_g_v128_store(m, v233, v246, v257)
				v454 = v233 + int32(32)
				v456 = v234
			}
		} else {
			v454 = l3
			v456 = int32(0)
		}
	case 2:
		v14 = int32(16)
		if v14 <= l1 {
			v19 = l1 + int32(-16)
			if base.Ui32(int32(16)) <= base.Ui32(v19) {
				v273 = (int32(base.Ui32(v19)>>(uint(int32(4))%32)) + int32(1)) & int32(536870910)
				v274 = int32(0)
				for {
					v282 = l3 + v274
					v283 = l0 + v274
					v284 = int32(0)
					v285 = base.Simd_g_v128_load(m, v283, v284)
					v286 = base.Simd_g_const(&F_BundleColorMap_SSE2__k15)
					v288 = base.Simd_g_const(&F_BundleColorMap_SSE2__k16)
					v289 = base.Simd_g_v128_and(base.Simd_g_i16x8_mul(v285, v286), v288)
					v290 = int32(12)
					v293 = base.Simd_g_const(&F_BundleColorMap_SSE2__k17)
					v294 = base.Simd_g_v128_or(base.Simd_g_v128_or(v289, base.Simd_g_i32x4_shr_u(v289, v290)), v293)
					base.Simd_g_v128_store(m, v282, v284, v294)
					v297 = int32(16)
					v302 = base.Simd_g_v128_load(m, v283+v297, v284)
					v304 = base.Simd_g_v128_and(base.Simd_g_i16x8_mul(v302, v286), v288)
					v308 = base.Simd_g_v128_or(base.Simd_g_v128_or(v304, base.Simd_g_i32x4_shr_u(v304, v290)), v293)
					base.Simd_g_v128_store(m, v282+v297, v284, v308)
					v312 = v274 + int32(32)
					v314 = v273 + int32(-2)
					if v314 != 0 {
						v273 = v314
						v274 = v312
						continue
					} else {
						break
					}
					break
				}
				v321 = l3 + v312
				v322 = v274 + int32(48)
				v323 = v312
			} else {
				v321 = l3
				v322 = v14
				v323 = int32(0)
			}
			if v19&int32(16) != 0 {
				v454 = v321
				v456 = v323
			} else {
				v334 = int32(0)
				v335 = base.Simd_g_v128_load(m, l0+v323, v334)
				v339 = base.Simd_g_v128_and(base.Simd_g_i16x8_mul(v335, base.Simd_g_const(&F_BundleColorMap_SSE2__k15)), base.Simd_g_const(&F_BundleColorMap_SSE2__k16))
				v344 = base.Simd_g_v128_or(base.Simd_g_v128_or(v339, base.Simd_g_i32x4_shr_u(v339, int32(12))), base.Simd_g_const(&F_BundleColorMap_SSE2__k17))
				base.Simd_g_v128_store(m, v321, v334, v344)
				v454 = v321 + int32(16)
				v456 = v322
			}
		} else {
			v454 = l3
			v456 = int32(0)
		}
	default:
		v41 = int32(16)
		if v41 <= l1 {
			v46 = l1 + int32(-16)
			if base.Ui32(int32(16)) <= base.Ui32(v46) {
				v359 = l3
				v360 = int32(16)
				v361 = (int32(base.Ui32(v46)>>(uint(int32(4))%32)) + int32(1)) & int32(536870910)
				for {
					v369 = l0 + v360
					v370 = int32(0)
					v371 = base.Simd_g_v128_load(m, v369, v370)
					v372 = int32(7)
					v374 = base.Simd_g_i8x16_bitmask(base.Simd_g_i64x2_shl(v371, v372))
					v375 = int32(_a_F_BundleColorMap_SSE2_0)
					v377 = int32(-16777216)
					*(*int32)(unsafe.Add(mBase, uint32(v359)+12)) = v374&v375 | v377
					v383 = base.Simd_g_v128_load(m, v369+int32(-16), v370)
					v386 = base.Simd_g_i8x16_bitmask(base.Simd_g_i64x2_shl(v383, v372))
					*(*int32)(unsafe.Add(mBase, uint32(v359)+4)) = v386&v375 | v377
					v392 = int32(8)
					*(*int32)(unsafe.Add(mBase, uint32(v359+v392))) = v374<<(uint(v392)%32)&v375 | v377
					*(*int32)(unsafe.Add(mBase, uint32(v359))) = v386<<(uint(v392)%32)&v375 | v377
					v409 = v360 + int32(32)
					v411 = v359 + int32(16)
					v413 = v361 + int32(-2)
					if v413 != 0 {
						v359 = v411
						v360 = v409
						v361 = v413
						continue
					} else {
						break
					}
					break
				}
				v419 = v411
				v420 = v409
				v421 = v360 + int32(16)
			} else {
				v419 = l3
				v420 = v41
				v421 = int32(0)
			}
			if v46&int32(16) != 0 {
				v454 = v419
				v456 = v421
			} else {
				v433 = base.Simd_g_v128_load(m, l0+v421, int32(0))
				v436 = base.Simd_g_i8x16_bitmask(base.Simd_g_i64x2_shl(v433, int32(7)))
				v437 = int32(_a_F_BundleColorMap_SSE2_0)
				v439 = int32(-16777216)
				*(*int32)(unsafe.Add(mBase, uint32(v419)+4)) = v436&v437 | v439
				v442 = int32(8)
				*(*int32)(unsafe.Add(mBase, uint32(v419))) = v436<<(uint(v442)%32)&v437 | v439
				v454 = v419 + v442
				v456 = v420
			}
		} else {
			v454 = l3
			v456 = int32(0)
		}
	}
	if l1 == v456 {
	} else {
		v465 = l0 + v456
		v466 = l1 - v456
		if int32(0) < l2 {
			if v466 < int32(1) {
			} else {
				v514 = int32(3) - l2
				v515 = int32(-1)
				v518 = v515<<(uint(l2)%32) ^ v515
				v519 = int32(1)
				if v466 != v519 {
					v530 = int32(-16777216)
					v535 = int32(0)
					for {
						v544 = v465 + v535
						v545 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v544))))
						v546 = v535 & v518
						if v546 != 0 {
							v552 = v530
						} else {
							v552 = int32(-16777216)
						}
						v553 = v545<<(uint(v546<<(uint(v514)%32)+int32(8))%32) | v552
						*(*int32)(unsafe.Add(mBase, uint32(v454+int32(base.Ui32(v535)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v553
						v555 = int32(1)
						v556 = v535 + v555
						v563 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v544+v555))))
						v564 = v556 & v518
						if v564 != 0 {
							v570 = v553
						} else {
							v570 = int32(-16777216)
						}
						v571 = v563<<(uint(v564<<(uint(v514)%32)+int32(8))%32) | v570
						*(*int32)(unsafe.Add(mBase, uint32(v454+int32(base.Ui32(v556)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v571
						v574 = v535 + int32(2)
						if v574 != v466&int32(2147483646) {
							v530 = v571
							v535 = v574
							continue
						} else {
							break
						}
						break
					}
					v577 = v571
					v582 = v574
				} else {
					v577 = int32(-16777216)
					v582 = int32(0)
				}
				if v466&v519 == int32(0) {
				} else {
					v594 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v465+v582))))
					v595 = v582 & v518
					if v595 != 0 {
						v601 = v577
					} else {
						v601 = int32(-16777216)
					}
					*(*int32)(unsafe.Add(mBase, uint32(v454+int32(base.Ui32(v582)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v594<<(uint(v595<<(uint(v514)%32)+int32(8))%32) | v601
				}
			}
		} else {
			if v466 < int32(1) {
			} else {
				if base.Ui32(v466) <= base.Ui32(int32(3)) {
					v608 = int32(0)
					v622 = v454 + v608<<(uint(int32(2))%32)
					v625 = v466 - v608
					v626 = v465 + v608
					for {
						v631 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v626))))
						*(*int32)(unsafe.Add(mBase, uint32(v622))) = v631<<(uint(int32(8))%32) | int32(-16777216)
						v642 = v625 + int32(-1)
						if v642 != 0 {
							v622 = v622 + int32(4)
							v625 = v642
							v626 = v626 + int32(1)
							continue
						} else {
							break
						}
						break
					}
				} else {
					v482 = v466 & int32(2147483644)
					v485 = v465
					v488 = v482
					v489 = v454
					for {
						v494 = int32(0)
						v495 = base.Simd_g_v128_load32_zero(m, v485, v494)
						v501 = base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v495)), int32(8)), base.Simd_g_const(&F_BundleColorMap_SSE2__k17))
						base.Simd_g_v128_store(m, v489, v494, v501)
						v509 = v488 + int32(-4)
						if v509 != 0 {
							v485 = v485 + int32(4)
							v488 = v509
							v489 = v489 + int32(16)
							continue
						} else {
							break
						}
						break
					}
					if v482 != v466 {
						v608 = v482
						v622 = v454 + v608<<(uint(int32(2))%32)
						v625 = v466 - v608
						v626 = v465 + v608
						for {
							v631 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v626))))
							*(*int32)(unsafe.Add(mBase, uint32(v622))) = v631<<(uint(int32(8))%32) | int32(-16777216)
							v642 = v625 + int32(-1)
							if v642 != 0 {
								v622 = v622 + int32(4)
								v625 = v642
								v626 = v626 + int32(1)
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

var F_BundleColorMap_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_BundleColorMap_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_BundleColorMap_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_BundleColorMap_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_BundleColorMap_SSE2__k4 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_BundleColorMap_SSE2__k5 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_BundleColorMap_SSE2__k6 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_BundleColorMap_SSE2__k7 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_BundleColorMap_SSE2__k8 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_BundleColorMap_SSE2__k9 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_BundleColorMap_SSE2__k10 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_BundleColorMap_SSE2__k11 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_BundleColorMap_SSE2__k12 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_BundleColorMap_SSE2__k13 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_BundleColorMap_SSE2__k14 = [2]uint64{0x110011001100110, 0x110011001100110}
var F_BundleColorMap_SSE2__k15 = [2]uint64{0x104010401040104, 0x104010401040104}
var F_BundleColorMap_SSE2__k16 = [2]uint64{0xf000f000f000f00, 0xf000f000f000f00}
var F_BundleColorMap_SSE2__k17 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}

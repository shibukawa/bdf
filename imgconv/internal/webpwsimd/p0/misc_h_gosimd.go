//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_HE16_SSE2(m *base.Module, l0 int32) {
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v11 base.V128
	_ = v11
	var v17 int32
	_ = v17
	v6 = l0 + int32(-1)
	v7 = int32(17)
	for {
		v8 = int32(1)
		v10 = int32(0)
		v11 = base.Simd_g_v128_load8_splat(m, v6, v10)
		base.Simd_g_v128_store(m, v6+v8, v10, v11)
		v17 = v7 + int32(-1)
		if base.Ui32(v8) < base.Ui32(v17) {
			v6 = v6 + int32(32)
			v7 = v17
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_HE16_SSE41(m *base.Module, l0 int32) {
	var v4 base.V128
	_ = v4
	var v5 base.V128
	_ = v5
	var v6 base.V128
	_ = v6
	var v10 base.V128
	_ = v10
	var v11 base.V128
	_ = v11
	var v15 base.V128
	_ = v15
	var v16 base.V128
	_ = v16
	var v20 base.V128
	_ = v20
	var v21 base.V128
	_ = v21
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v35 base.V128
	_ = v35
	var v36 base.V128
	_ = v36
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v46 int32
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v52 base.V128
	_ = v52
	var v53 base.V128
	_ = v53
	var v57 base.V128
	_ = v57
	var v58 base.V128
	_ = v58
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v72 base.V128
	_ = v72
	var v73 base.V128
	_ = v73
	var v77 base.V128
	_ = v77
	var v78 base.V128
	_ = v78
	var v82 base.V128
	_ = v82
	var v83 base.V128
	_ = v83
	v4 = base.Simd_g_v128_load32_zero(m, l0, int32(28))
	v5 = base.Simd_g_const(&F_HE16_SSE41__k0)
	v6 = base.Simd_g_i8x16_swizzle(v4, v5)
	base.Simd_g_v128_store(m, l0, int32(32), v6)
	v10 = base.Simd_g_v128_load32_zero(m, l0, int32(60))
	v11 = base.Simd_g_i8x16_swizzle(v10, v5)
	base.Simd_g_v128_store(m, l0, int32(64), v11)
	v15 = base.Simd_g_v128_load32_zero(m, l0, int32(92))
	v16 = base.Simd_g_i8x16_swizzle(v15, v5)
	base.Simd_g_v128_store(m, l0, int32(96), v16)
	v20 = base.Simd_g_v128_load32_zero(m, l0, int32(124))
	v21 = base.Simd_g_i8x16_swizzle(v20, v5)
	base.Simd_g_v128_store(m, l0, int32(128), v21)
	v25 = base.Simd_g_v128_load32_zero(m, l0, int32(156))
	v26 = base.Simd_g_i8x16_swizzle(v25, v5)
	base.Simd_g_v128_store(m, l0, int32(160), v26)
	v30 = base.Simd_g_v128_load32_zero(m, l0, int32(188))
	v31 = base.Simd_g_i8x16_swizzle(v30, v5)
	base.Simd_g_v128_store(m, l0, int32(192), v31)
	v35 = base.Simd_g_v128_load32_zero(m, l0, int32(220))
	v36 = base.Simd_g_i8x16_swizzle(v35, v5)
	base.Simd_g_v128_store(m, l0, int32(224), v36)
	v40 = base.Simd_g_v128_load32_zero(m, l0, int32(252))
	v41 = base.Simd_g_i8x16_swizzle(v40, v5)
	base.Simd_g_v128_store(m, l0, int32(256), v41)
	v46 = int32(0)
	v47 = base.Simd_g_v128_load32_zero(m, l0+int32(-4), v46)
	v48 = base.Simd_g_i8x16_swizzle(v47, v5)
	base.Simd_g_v128_store(m, l0, v46, v48)
	v52 = base.Simd_g_v128_load32_zero(m, l0, int32(284))
	v53 = base.Simd_g_i8x16_swizzle(v52, v5)
	base.Simd_g_v128_store(m, l0, int32(288), v53)
	v57 = base.Simd_g_v128_load32_zero(m, l0, int32(316))
	v58 = base.Simd_g_i8x16_swizzle(v57, v5)
	base.Simd_g_v128_store(m, l0, int32(320), v58)
	v62 = base.Simd_g_v128_load32_zero(m, l0, int32(348))
	v63 = base.Simd_g_i8x16_swizzle(v62, v5)
	base.Simd_g_v128_store(m, l0, int32(352), v63)
	v67 = base.Simd_g_v128_load32_zero(m, l0, int32(380))
	v68 = base.Simd_g_i8x16_swizzle(v67, v5)
	base.Simd_g_v128_store(m, l0, int32(384), v68)
	v72 = base.Simd_g_v128_load32_zero(m, l0, int32(412))
	v73 = base.Simd_g_i8x16_swizzle(v72, v5)
	base.Simd_g_v128_store(m, l0, int32(416), v73)
	v77 = base.Simd_g_v128_load32_zero(m, l0, int32(444))
	v78 = base.Simd_g_i8x16_swizzle(v77, v5)
	base.Simd_g_v128_store(m, l0, int32(448), v78)
	v82 = base.Simd_g_v128_load32_zero(m, l0, int32(476))
	v83 = base.Simd_g_i8x16_swizzle(v82, v5)
	base.Simd_g_v128_store(m, l0, int32(480), v83)
	return
}

var F_HE16_SSE41__k0 = [2]uint64{0x303030303030303, 0x303030303030303}

func F_HFilter16_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	var v6 int32
	_ = v6
	var v11 base.V128
	_ = v11
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v52 base.V128
	_ = v52
	var v55 base.V128
	_ = v55
	var v58 base.V128
	_ = v58
	var v60 int32
	_ = v60
	var v61 base.V128
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v73 base.V128
	_ = v73
	var v76 base.V128
	_ = v76
	var v79 base.V128
	_ = v79
	var v82 base.V128
	_ = v82
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v86 base.V128
	_ = v86
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v100 base.V128
	_ = v100
	var v103 base.V128
	_ = v103
	var v106 base.V128
	_ = v106
	var v109 base.V128
	_ = v109
	var v113 int32
	_ = v113
	var v115 base.V128
	_ = v115
	var v118 base.V128
	_ = v118
	var v121 base.V128
	_ = v121
	var v124 base.V128
	_ = v124
	var v126 base.V128
	_ = v126
	var v128 base.V128
	_ = v128
	var v130 base.V128
	_ = v130
	var v132 base.V128
	_ = v132
	var v134 base.V128
	_ = v134
	var v135 base.V128
	_ = v135
	var v136 base.V128
	_ = v136
	var v137 base.V128
	_ = v137
	var v138 base.V128
	_ = v138
	var v139 base.V128
	_ = v139
	var v142 base.V128
	_ = v142
	var v143 base.V128
	_ = v143
	var v145 base.V128
	_ = v145
	var v147 base.V128
	_ = v147
	var v150 int32
	_ = v150
	var v152 base.V128
	_ = v152
	var v155 base.V128
	_ = v155
	var v158 base.V128
	_ = v158
	var v161 base.V128
	_ = v161
	var v165 int32
	_ = v165
	var v167 base.V128
	_ = v167
	var v170 base.V128
	_ = v170
	var v173 base.V128
	_ = v173
	var v176 base.V128
	_ = v176
	var v178 base.V128
	_ = v178
	var v180 base.V128
	_ = v180
	var v182 base.V128
	_ = v182
	var v184 base.V128
	_ = v184
	var v186 base.V128
	_ = v186
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v192 base.V128
	_ = v192
	var v195 base.V128
	_ = v195
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v205 int32
	_ = v205
	var v207 base.V128
	_ = v207
	var v210 base.V128
	_ = v210
	var v213 base.V128
	_ = v213
	var v216 base.V128
	_ = v216
	var v218 base.V128
	_ = v218
	var v220 base.V128
	_ = v220
	var v222 base.V128
	_ = v222
	var v224 base.V128
	_ = v224
	var v226 base.V128
	_ = v226
	var v228 base.V128
	_ = v228
	var v231 base.V128
	_ = v231
	var v234 base.V128
	_ = v234
	var v236 base.V128
	_ = v236
	var v249 base.V128
	_ = v249
	var v260 base.V128
	_ = v260
	var v263 base.V128
	_ = v263
	var v265 base.V128
	_ = v265
	var v267 base.V128
	_ = v267
	var v269 base.V128
	_ = v269
	var v281 base.V128
	_ = v281
	var v282 base.V128
	_ = v282
	var v283 base.V128
	_ = v283
	var v284 base.V128
	_ = v284
	var v285 base.V128
	_ = v285
	var v286 base.V128
	_ = v286
	var v290 base.V128
	_ = v290
	var v295 base.V128
	_ = v295
	var v296 base.V128
	_ = v296
	var v298 base.V128
	_ = v298
	var v300 base.V128
	_ = v300
	var v304 base.V128
	_ = v304
	var v305 base.V128
	_ = v305
	var v306 base.V128
	_ = v306
	var v307 base.V128
	_ = v307
	var v311 base.V128
	_ = v311
	var v317 base.V128
	_ = v317
	var v318 base.V128
	_ = v318
	var v321 base.V128
	_ = v321
	var v323 base.V128
	_ = v323
	var v325 base.V128
	_ = v325
	var v326 base.V128
	_ = v326
	var v329 base.V128
	_ = v329
	var v332 base.V128
	_ = v332
	var v334 base.V128
	_ = v334
	var v336 base.V128
	_ = v336
	var v338 base.V128
	_ = v338
	var v341 int32
	_ = v341
	var v355 base.V128
	_ = v355
	var v357 base.V128
	_ = v357
	var v359 base.V128
	_ = v359
	var v361 base.V128
	_ = v361
	var v365 base.V128
	_ = v365
	var v366 base.V128
	_ = v366
	var v370 int32
	_ = v370
	var v372 base.V128
	_ = v372
	var v383 base.V128
	_ = v383
	var v387 int32
	_ = v387
	var v389 base.V128
	_ = v389
	var v393 int32
	_ = v393
	var v395 base.V128
	_ = v395
	var v406 base.V128
	_ = v406
	var v408 base.V128
	_ = v408
	var v410 base.V128
	_ = v410
	var v415 base.V128
	_ = v415
	var v419 int32
	_ = v419
	var v421 base.V128
	_ = v421
	var v432 base.V128
	_ = v432
	var v436 int32
	_ = v436
	var v438 base.V128
	_ = v438
	var v442 int32
	_ = v442
	var v444 base.V128
	_ = v444
	var v455 base.V128
	_ = v455
	var v467 base.V128
	_ = v467
	var v469 base.V128
	_ = v469
	var v471 base.V128
	_ = v471
	var v474 base.V128
	_ = v474
	var v476 base.V128
	_ = v476
	var v478 base.V128
	_ = v478
	var v483 base.V128
	_ = v483
	var v487 int32
	_ = v487
	var v489 base.V128
	_ = v489
	var v500 base.V128
	_ = v500
	var v504 int32
	_ = v504
	var v506 base.V128
	_ = v506
	var v510 int32
	_ = v510
	var v512 base.V128
	_ = v512
	var v523 base.V128
	_ = v523
	var v525 base.V128
	_ = v525
	var v527 base.V128
	_ = v527
	var v532 base.V128
	_ = v532
	var v536 int32
	_ = v536
	var v538 base.V128
	_ = v538
	var v549 base.V128
	_ = v549
	var v553 int32
	_ = v553
	var v555 base.V128
	_ = v555
	var v559 int32
	_ = v559
	var v561 base.V128
	_ = v561
	v6 = int32(0)
	v11 = base.Simd_g_const(&F_HFilter16_SSE2__k0)
	v41 = l0 + int32(-4)
	v43 = l1 * int32(6)
	v45 = int32(1)
	v46 = l1 << (uint(v45) % 32)
	v48 = int32(2)
	v49 = l1 << (uint(v48) % 32)
	v50 = v41 + v49
	v52 = base.Simd_g_v128_load32_splat(m, v41, v6)
	v55 = base.Simd_g_v128_load32_lane_l1(m, v50, v6, v52)
	v58 = base.Simd_g_v128_load32_lane_l2(m, v41+v46, v6, v55)
	v60 = int32(3)
	v61 = base.Simd_g_v128_load32_lane_l3(m, v41+v43, v6, v58)
	v62 = int32(7)
	v63 = l1 * v62
	v66 = l1 * v60
	v69 = l1 * int32(5)
	v71 = v41 + l1
	v73 = base.Simd_g_v128_load32_splat(m, v71, v6)
	v76 = base.Simd_g_v128_load32_lane_l1(m, v41+v69, v6, v73)
	v79 = base.Simd_g_v128_load32_lane_l2(m, v41+v66, v6, v76)
	v82 = base.Simd_g_v128_load32_lane_l3(m, v41+v63, v6, v79)
	v83 = base.Simd_g_const(&F_HFilter16_SSE2__k1)
	v84 = base.Simd_g_i8x16_shuffle2(v61, v82, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v85 = base.Simd_g_const(&F_HFilter16_SSE2__k4)
	v86 = base.Simd_g_i8x16_shuffle2(v61, v82, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v87 = base.Simd_g_const(&F_HFilter16_SSE2__k7)
	v88 = base.Simd_g_i8x16_shuffle2(v84, v86, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	v89 = base.Simd_g_const(&F_HFilter16_SSE2__k10)
	v90 = base.Simd_g_i8x16_shuffle2(v84, v86, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	v91 = base.Simd_g_const(&F_HFilter16_SSE2__k13)
	v92 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_HFilter16_SSE2__k14), base.Simd_g_const(&F_HFilter16_SSE2__k15))
	v94 = l1 << (uint(v60) % 32)
	v95 = v41 + v94
	v98 = v95 + v49
	v100 = base.Simd_g_v128_load32_splat(m, v95, v6)
	v103 = base.Simd_g_v128_load32_lane_l1(m, v98, v6, v100)
	v106 = base.Simd_g_v128_load32_lane_l2(m, v95+v46, v6, v103)
	v109 = base.Simd_g_v128_load32_lane_l3(m, v95+v43, v6, v106)
	v113 = v95 + l1
	v115 = base.Simd_g_v128_load32_splat(m, v113, v6)
	v118 = base.Simd_g_v128_load32_lane_l1(m, v95+v69, v6, v115)
	v121 = base.Simd_g_v128_load32_lane_l2(m, v95+v66, v6, v118)
	v124 = base.Simd_g_v128_load32_lane_l3(m, v95+v63, v6, v121)
	v126 = base.Simd_g_i8x16_shuffle2(v109, v124, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v128 = base.Simd_g_i8x16_shuffle2(v109, v124, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v130 = base.Simd_g_i8x16_shuffle2(v126, v128, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	v132 = base.Simd_g_i8x16_shuffle2(v126, v128, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	v134 = base.Simd_g_i8x16_shuffle2(v130, v132, base.Simd_g_const(&F_HFilter16_SSE2__k14), base.Simd_g_const(&F_HFilter16_SSE2__k15))
	v135 = base.Simd_g_const(&F_HFilter16_SSE2__k16)
	v136 = base.Simd_g_i8x16_shuffle2(v92, v134, base.Simd_g_const(&F_HFilter16_SSE2__k17), base.Simd_g_const(&F_HFilter16_SSE2__k18))
	v137 = base.Simd_g_const(&F_HFilter16_SSE2__k19)
	v138 = base.Simd_g_i8x16_shuffle2(v92, v134, base.Simd_g_const(&F_HFilter16_SSE2__k20), base.Simd_g_const(&F_HFilter16_SSE2__k21))
	v139 = base.Simd_g_const(&F_HFilter16_SSE2__k22)
	v142 = base.Simd_g_const(&F_HFilter16_SSE2__k23)
	v143 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_HFilter16_SSE2__k24), base.Simd_g_const(&F_HFilter16_SSE2__k25))
	v145 = base.Simd_g_i8x16_shuffle2(v130, v132, base.Simd_g_const(&F_HFilter16_SSE2__k24), base.Simd_g_const(&F_HFilter16_SSE2__k25))
	v147 = base.Simd_g_i8x16_shuffle2(v143, v145, base.Simd_g_const(&F_HFilter16_SSE2__k20), base.Simd_g_const(&F_HFilter16_SSE2__k21))
	v150 = l0 + v49
	v152 = base.Simd_g_v128_load32_splat(m, l0, v6)
	v155 = base.Simd_g_v128_load32_lane_l1(m, v150, v6, v152)
	v158 = base.Simd_g_v128_load32_lane_l2(m, l0+v46, v6, v155)
	v161 = base.Simd_g_v128_load32_lane_l3(m, l0+v43, v6, v158)
	v165 = l0 + l1
	v167 = base.Simd_g_v128_load32_splat(m, v165, v6)
	v170 = base.Simd_g_v128_load32_lane_l1(m, l0+v69, v6, v167)
	v173 = base.Simd_g_v128_load32_lane_l2(m, l0+v66, v6, v170)
	v176 = base.Simd_g_v128_load32_lane_l3(m, l0+v63, v6, v173)
	v178 = base.Simd_g_i8x16_shuffle2(v161, v176, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v180 = base.Simd_g_i8x16_shuffle2(v161, v176, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v182 = base.Simd_g_i8x16_shuffle2(v178, v180, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	v184 = base.Simd_g_i8x16_shuffle2(v178, v180, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	v186 = base.Simd_g_i8x16_shuffle2(v182, v184, base.Simd_g_const(&F_HFilter16_SSE2__k14), base.Simd_g_const(&F_HFilter16_SSE2__k15))
	v187 = l0 + v94
	v190 = v187 + v49
	v192 = base.Simd_g_v128_load32_splat(m, v187, v6)
	v195 = base.Simd_g_v128_load32_lane_l1(m, v190, v6, v192)
	v198 = base.Simd_g_v128_load32_lane_l2(m, v187+v46, v6, v195)
	v201 = base.Simd_g_v128_load32_lane_l3(m, v187+v43, v6, v198)
	v205 = v187 + l1
	v207 = base.Simd_g_v128_load32_splat(m, v205, v6)
	v210 = base.Simd_g_v128_load32_lane_l1(m, v187+v69, v6, v207)
	v213 = base.Simd_g_v128_load32_lane_l2(m, v187+v66, v6, v210)
	v216 = base.Simd_g_v128_load32_lane_l3(m, v187+v63, v6, v213)
	v218 = base.Simd_g_i8x16_shuffle2(v201, v216, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v220 = base.Simd_g_i8x16_shuffle2(v201, v216, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v222 = base.Simd_g_i8x16_shuffle2(v218, v220, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	v224 = base.Simd_g_i8x16_shuffle2(v218, v220, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	v226 = base.Simd_g_i8x16_shuffle2(v222, v224, base.Simd_g_const(&F_HFilter16_SSE2__k14), base.Simd_g_const(&F_HFilter16_SSE2__k15))
	v228 = base.Simd_g_i8x16_shuffle2(v186, v226, base.Simd_g_const(&F_HFilter16_SSE2__k17), base.Simd_g_const(&F_HFilter16_SSE2__k18))
	v231 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v147, v228), base.Simd_g_i8x16_sub_sat_u(v228, v147))
	v234 = base.Simd_g_i8x16_shuffle2(v143, v145, base.Simd_g_const(&F_HFilter16_SSE2__k17), base.Simd_g_const(&F_HFilter16_SSE2__k18))
	v236 = base.Simd_g_i8x16_shuffle2(v186, v226, base.Simd_g_const(&F_HFilter16_SSE2__k20), base.Simd_g_const(&F_HFilter16_SSE2__k21))
	v249 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v234, v147), base.Simd_g_i8x16_sub_sat_u(v147, v234))
	v260 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v236, v228), base.Simd_g_i8x16_sub_sat_u(v228, v236))
	v263 = base.Simd_g_i8x16_shuffle2(v182, v184, base.Simd_g_const(&F_HFilter16_SSE2__k24), base.Simd_g_const(&F_HFilter16_SSE2__k25))
	v265 = base.Simd_g_i8x16_shuffle2(v222, v224, base.Simd_g_const(&F_HFilter16_SSE2__k24), base.Simd_g_const(&F_HFilter16_SSE2__k25))
	v267 = base.Simd_g_i8x16_shuffle2(v263, v265, base.Simd_g_const(&F_HFilter16_SSE2__k20), base.Simd_g_const(&F_HFilter16_SSE2__k21))
	v269 = base.Simd_g_i8x16_shuffle2(v263, v265, base.Simd_g_const(&F_HFilter16_SSE2__k17), base.Simd_g_const(&F_HFilter16_SSE2__k18))
	v281 = base.Simd_g_i8x16_eq(v11, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v231, v231), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v234, v236), base.Simd_g_i8x16_sub_sat_u(v236, v234)), v45), base.Simd_g_const(&F_HFilter16_SSE2__k26))), base.Simd_g_i8x16_splat(l2)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v249, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v136, v138), base.Simd_g_i8x16_sub_sat_u(v138, v136))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v138, v234), base.Simd_g_i8x16_sub_sat_u(v234, v138))), v260), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v267, v269), base.Simd_g_i8x16_sub_sat_u(v269, v267))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v269, v236), base.Simd_g_i8x16_sub_sat_u(v236, v269))), base.Simd_g_i8x16_splat(l3))))
	v282 = base.Simd_g_v128_xor(v228, v139)
	v283 = base.Simd_g_v128_xor(v147, v139)
	v284 = base.Simd_g_i8x16_sub_sat_s(v282, v283)
	v285 = base.Simd_g_v128_xor(v234, v139)
	v286 = base.Simd_g_v128_xor(v236, v139)
	v290 = base.Simd_g_i8x16_add_sat_s(v284, base.Simd_g_i8x16_add_sat_s(v284, base.Simd_g_i8x16_add_sat_s(v284, base.Simd_g_i8x16_sub_sat_s(v285, v286))))
	v295 = base.Simd_g_i8x16_eq(v11, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v249, v260), base.Simd_g_i8x16_splat(l4)))
	v296 = base.Simd_g_v128_and(base.Simd_g_v128_and(v281, v290), v295)
	v298 = base.Simd_g_i8x16_shuffle2(v11, v296, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v300 = base.Simd_g_const(&F_HFilter16_SSE2__k27)
	v304 = base.Simd_g_const(&F_HFilter16_SSE2__k28)
	v305 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v298), v300), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v298), v300), base.Simd_g_const(&F_HFilter16_SSE2__k29), base.Simd_g_const(&F_HFilter16_SSE2__k30))
	v306 = base.Simd_g_const(&F_HFilter16_SSE2__k31)
	v307 = base.Simd_g_i16x8_add(v305, v306)
	v311 = base.Simd_g_i8x16_shuffle2(v11, v296, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v317 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v311), v300), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v311), v300), base.Simd_g_const(&F_HFilter16_SSE2__k29), base.Simd_g_const(&F_HFilter16_SSE2__k30))
	v318 = base.Simd_g_i16x8_add(v317, v306)
	v321 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v307, v62), base.Simd_g_i16x8_shr_s(v318, v62))
	v323 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_xor(v138, v139), v321), v139)
	v325 = base.Simd_g_i8x16_shuffle2(v136, v323, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v326 = base.Simd_g_i16x8_add(v307, v305)
	v329 = base.Simd_g_i16x8_add(v318, v317)
	v332 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v326, v62), base.Simd_g_i16x8_shr_s(v329, v62))
	v334 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v285, v332), v139)
	v336 = base.Simd_g_v128_and(base.Simd_g_v128_andnot(v290, v295), v281)
	v338 = base.Simd_g_i8x16_add_sat_s(v336, base.Simd_g_const(&F_HFilter16_SSE2__k32))
	v341 = int32(11)
	v355 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v326, v305), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v329, v317), v62))
	v357 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(v283, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v338, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3)), v341), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v338, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6)), v341))), v355), v139)
	v359 = base.Simd_g_i8x16_shuffle2(v334, v357, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v361 = base.Simd_g_i8x16_shuffle2(v325, v359, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v41, v6, v361)
	v365 = base.Simd_g_const(&F_HFilter16_SSE2__k33)
	v366 = base.Simd_g_i8x16_shuffle2(v361, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v71, v6, v366)
	v370 = v71 + l1
	v372 = base.Simd_g_i8x16_shuffle2(v366, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v370, v6, v372)
	base.Simd_g_v128_store32_lane_l0(m, v370+l1, v6, base.Simd_g_i8x16_shuffle2(v372, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v383 = base.Simd_g_i8x16_shuffle2(v325, v359, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v50, v6, v383)
	v387 = v50 + l1
	v389 = base.Simd_g_i8x16_shuffle2(v383, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v387, v6, v389)
	v393 = v387 + l1
	v395 = base.Simd_g_i8x16_shuffle2(v389, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v393, v6, v395)
	base.Simd_g_v128_store32_lane_l0(m, v393+l1, v6, base.Simd_g_i8x16_shuffle2(v395, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v406 = base.Simd_g_i8x16_shuffle2(v136, v323, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v408 = base.Simd_g_i8x16_shuffle2(v334, v357, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v410 = base.Simd_g_i8x16_shuffle2(v406, v408, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v95, v6, v410)
	v415 = base.Simd_g_i8x16_shuffle2(v410, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v113, v6, v415)
	v419 = v113 + l1
	v421 = base.Simd_g_i8x16_shuffle2(v415, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v419, v6, v421)
	base.Simd_g_v128_store32_lane_l0(m, v419+l1, v6, base.Simd_g_i8x16_shuffle2(v421, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v432 = base.Simd_g_i8x16_shuffle2(v406, v408, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v98, v6, v432)
	v436 = v98 + l1
	v438 = base.Simd_g_i8x16_shuffle2(v432, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v436, v6, v438)
	v442 = v436 + l1
	v444 = base.Simd_g_i8x16_shuffle2(v438, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v442, v6, v444)
	base.Simd_g_v128_store32_lane_l0(m, v442+l1, v6, base.Simd_g_i8x16_shuffle2(v444, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v455 = base.Simd_g_i8x16_add_sat_s(v336, base.Simd_g_const(&F_HFilter16_SSE2__k36))
	v467 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_i8x16_sub_sat_s(v282, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v455, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3)), v341), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v11, v455, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6)), v341))), v355), v139)
	v469 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v286, v332), v139)
	v471 = base.Simd_g_i8x16_shuffle2(v467, v469, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v474 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v269, v139), v321), v139)
	v476 = base.Simd_g_i8x16_shuffle2(v474, v267, base.Simd_g_const(&F_HFilter16_SSE2__k2), base.Simd_g_const(&F_HFilter16_SSE2__k3))
	v478 = base.Simd_g_i8x16_shuffle2(v471, v476, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, l0, v6, v478)
	v483 = base.Simd_g_i8x16_shuffle2(v478, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v165, v6, v483)
	v487 = v165 + l1
	v489 = base.Simd_g_i8x16_shuffle2(v483, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v487, v6, v489)
	base.Simd_g_v128_store32_lane_l0(m, v487+l1, v6, base.Simd_g_i8x16_shuffle2(v489, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v500 = base.Simd_g_i8x16_shuffle2(v471, v476, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v150, v6, v500)
	v504 = v150 + l1
	v506 = base.Simd_g_i8x16_shuffle2(v500, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v504, v6, v506)
	v510 = v504 + l1
	v512 = base.Simd_g_i8x16_shuffle2(v506, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v510, v6, v512)
	base.Simd_g_v128_store32_lane_l0(m, v510+l1, v6, base.Simd_g_i8x16_shuffle2(v512, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v523 = base.Simd_g_i8x16_shuffle2(v467, v469, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v525 = base.Simd_g_i8x16_shuffle2(v474, v267, base.Simd_g_const(&F_HFilter16_SSE2__k5), base.Simd_g_const(&F_HFilter16_SSE2__k6))
	v527 = base.Simd_g_i8x16_shuffle2(v523, v525, base.Simd_g_const(&F_HFilter16_SSE2__k8), base.Simd_g_const(&F_HFilter16_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v187, v6, v527)
	v532 = base.Simd_g_i8x16_shuffle2(v527, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v205, v6, v532)
	v536 = v205 + l1
	v538 = base.Simd_g_i8x16_shuffle2(v532, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v536, v6, v538)
	base.Simd_g_v128_store32_lane_l0(m, v536+l1, v6, base.Simd_g_i8x16_shuffle2(v538, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	v549 = base.Simd_g_i8x16_shuffle2(v523, v525, base.Simd_g_const(&F_HFilter16_SSE2__k11), base.Simd_g_const(&F_HFilter16_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v190, v6, v549)
	v553 = v190 + l1
	v555 = base.Simd_g_i8x16_shuffle2(v549, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v553, v6, v555)
	v559 = v553 + l1
	v561 = base.Simd_g_i8x16_shuffle2(v555, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v559, v6, v561)
	base.Simd_g_v128_store32_lane_l0(m, v559+l1, v6, base.Simd_g_i8x16_shuffle2(v561, v11, base.Simd_g_const(&F_HFilter16_SSE2__k34), base.Simd_g_const(&F_HFilter16_SSE2__k35)))
	return
}

var F_HFilter16_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_HFilter16_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_HFilter16_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_HFilter16_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_HFilter16_SSE2__k4 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_HFilter16_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_HFilter16_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_HFilter16_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_HFilter16_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_HFilter16_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_HFilter16_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_HFilter16_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_HFilter16_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_HFilter16_SSE2__k13 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_HFilter16_SSE2__k14 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_HFilter16_SSE2__k15 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_HFilter16_SSE2__k16 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_HFilter16_SSE2__k17 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_HFilter16_SSE2__k18 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_HFilter16_SSE2__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_HFilter16_SSE2__k20 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_HFilter16_SSE2__k21 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_HFilter16_SSE2__k22 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_HFilter16_SSE2__k23 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_HFilter16_SSE2__k24 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_HFilter16_SSE2__k25 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_HFilter16_SSE2__k26 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_HFilter16_SSE2__k27 = [2]uint64{0x90000000900, 0x90000000900}
var F_HFilter16_SSE2__k28 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_HFilter16_SSE2__k29 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_HFilter16_SSE2__k30 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_HFilter16_SSE2__k31 = [2]uint64{0x3f003f003f003f, 0x3f003f003f003f}
var F_HFilter16_SSE2__k32 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_HFilter16_SSE2__k33 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_HFilter16_SSE2__k34 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_HFilter16_SSE2__k35 = [2]uint64{0x8080808080808080, 0x302010080808080}
var F_HFilter16_SSE2__k36 = [2]uint64{0x404040404040404, 0x404040404040404}

func F_HFilter16i_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	var v6 int32
	_ = v6
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v65 base.V128
	_ = v65
	var v68 base.V128
	_ = v68
	var v71 base.V128
	_ = v71
	var v73 int32
	_ = v73
	var v74 base.V128
	_ = v74
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v86 base.V128
	_ = v86
	var v89 base.V128
	_ = v89
	var v92 base.V128
	_ = v92
	var v95 base.V128
	_ = v95
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v98 base.V128
	_ = v98
	var v99 base.V128
	_ = v99
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v102 base.V128
	_ = v102
	var v103 base.V128
	_ = v103
	var v104 base.V128
	_ = v104
	var v105 base.V128
	_ = v105
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v113 base.V128
	_ = v113
	var v116 base.V128
	_ = v116
	var v119 base.V128
	_ = v119
	var v122 base.V128
	_ = v122
	var v128 base.V128
	_ = v128
	var v131 base.V128
	_ = v131
	var v134 base.V128
	_ = v134
	var v137 base.V128
	_ = v137
	var v139 base.V128
	_ = v139
	var v141 base.V128
	_ = v141
	var v143 base.V128
	_ = v143
	var v145 base.V128
	_ = v145
	var v147 base.V128
	_ = v147
	var v148 base.V128
	_ = v148
	var v150 base.V128
	_ = v150
	var v152 base.V128
	_ = v152
	var v153 base.V128
	_ = v153
	var v155 base.V128
	_ = v155
	var v178 int32
	_ = v178
	var v190 base.V128
	_ = v190
	var v194 int32
	_ = v194
	var v195 base.V128
	_ = v195
	var v196 base.V128
	_ = v196
	var v198 base.V128
	_ = v198
	var v232 base.V128
	_ = v232
	var v233 base.V128
	_ = v233
	var v234 base.V128
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v238 int32
	_ = v238
	var v241 int32
	_ = v241
	var v245 base.V128
	_ = v245
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v248 base.V128
	_ = v248
	var v250 int32
	_ = v250
	var v251 base.V128
	_ = v251
	var v253 int32
	_ = v253
	var v254 base.V128
	_ = v254
	var v255 int32
	_ = v255
	var v258 int32
	_ = v258
	var v261 int32
	_ = v261
	var v264 int32
	_ = v264
	var v268 base.V128
	_ = v268
	var v271 base.V128
	_ = v271
	var v274 base.V128
	_ = v274
	var v277 base.V128
	_ = v277
	var v278 base.V128
	_ = v278
	var v279 base.V128
	_ = v279
	var v280 base.V128
	_ = v280
	var v281 base.V128
	_ = v281
	var v282 base.V128
	_ = v282
	var v283 base.V128
	_ = v283
	var v284 base.V128
	_ = v284
	var v285 base.V128
	_ = v285
	var v286 base.V128
	_ = v286
	var v287 base.V128
	_ = v287
	var v288 int32
	_ = v288
	var v291 int32
	_ = v291
	var v294 int32
	_ = v294
	var v297 int32
	_ = v297
	var v301 base.V128
	_ = v301
	var v304 base.V128
	_ = v304
	var v307 base.V128
	_ = v307
	var v310 base.V128
	_ = v310
	var v311 int32
	_ = v311
	var v314 int32
	_ = v314
	var v317 int32
	_ = v317
	var v320 int32
	_ = v320
	var v324 base.V128
	_ = v324
	var v327 base.V128
	_ = v327
	var v330 base.V128
	_ = v330
	var v333 base.V128
	_ = v333
	var v335 base.V128
	_ = v335
	var v337 base.V128
	_ = v337
	var v339 base.V128
	_ = v339
	var v341 base.V128
	_ = v341
	var v343 base.V128
	_ = v343
	var v344 base.V128
	_ = v344
	var v345 base.V128
	_ = v345
	var v346 base.V128
	_ = v346
	var v350 base.V128
	_ = v350
	var v351 base.V128
	_ = v351
	var v352 base.V128
	_ = v352
	var v355 base.V128
	_ = v355
	var v358 base.V128
	_ = v358
	var v360 base.V128
	_ = v360
	var v361 base.V128
	_ = v361
	var v362 base.V128
	_ = v362
	var v368 base.V128
	_ = v368
	var v387 base.V128
	_ = v387
	var v388 base.V128
	_ = v388
	var v390 base.V128
	_ = v390
	var v392 base.V128
	_ = v392
	var v394 base.V128
	_ = v394
	var v407 base.V128
	_ = v407
	var v409 base.V128
	_ = v409
	var v412 int32
	_ = v412
	var v418 base.V128
	_ = v418
	var v423 base.V128
	_ = v423
	var v425 base.V128
	_ = v425
	var v427 base.V128
	_ = v427
	var v438 base.V128
	_ = v438
	var v440 base.V128
	_ = v440
	var v443 base.V128
	_ = v443
	var v445 base.V128
	_ = v445
	var v447 base.V128
	_ = v447
	var v449 base.V128
	_ = v449
	var v455 base.V128
	_ = v455
	var v456 base.V128
	_ = v456
	var v463 base.V128
	_ = v463
	var v477 base.V128
	_ = v477
	var v484 base.V128
	_ = v484
	var v491 base.V128
	_ = v491
	var v505 base.V128
	_ = v505
	var v507 base.V128
	_ = v507
	var v509 base.V128
	_ = v509
	var v516 base.V128
	_ = v516
	var v523 base.V128
	_ = v523
	var v537 base.V128
	_ = v537
	var v544 base.V128
	_ = v544
	var v551 base.V128
	_ = v551
	var v565 int32
	_ = v565
	v6 = int32(0)
	v56 = l1 * int32(6)
	v58 = int32(1)
	v59 = l1 << (uint(v58) % 32)
	v61 = int32(2)
	v62 = l1 << (uint(v61) % 32)
	v65 = base.Simd_g_v128_load32_splat(m, l0, v6)
	v68 = base.Simd_g_v128_load32_lane_l1(m, l0+v62, v6, v65)
	v71 = base.Simd_g_v128_load32_lane_l2(m, l0+v59, v6, v68)
	v73 = int32(3)
	v74 = base.Simd_g_v128_load32_lane_l3(m, l0+v56, v6, v71)
	v76 = l1 * int32(7)
	v79 = l1 * v73
	v82 = l1 * int32(5)
	v86 = base.Simd_g_v128_load32_splat(m, l0+l1, v6)
	v89 = base.Simd_g_v128_load32_lane_l1(m, l0+v82, v6, v86)
	v92 = base.Simd_g_v128_load32_lane_l2(m, l0+v79, v6, v89)
	v95 = base.Simd_g_v128_load32_lane_l3(m, l0+v76, v6, v92)
	v96 = base.Simd_g_const(&F_HFilter16i_SSE2__k0)
	v97 = base.Simd_g_i8x16_shuffle2(v74, v95, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2))
	v98 = base.Simd_g_const(&F_HFilter16i_SSE2__k3)
	v99 = base.Simd_g_i8x16_shuffle2(v74, v95, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5))
	v100 = base.Simd_g_const(&F_HFilter16i_SSE2__k6)
	v101 = base.Simd_g_i8x16_shuffle2(v97, v99, base.Simd_g_const(&F_HFilter16i_SSE2__k7), base.Simd_g_const(&F_HFilter16i_SSE2__k8))
	v102 = base.Simd_g_const(&F_HFilter16i_SSE2__k9)
	v103 = base.Simd_g_i8x16_shuffle2(v97, v99, base.Simd_g_const(&F_HFilter16i_SSE2__k10), base.Simd_g_const(&F_HFilter16i_SSE2__k11))
	v104 = base.Simd_g_const(&F_HFilter16i_SSE2__k12)
	v105 = base.Simd_g_i8x16_shuffle2(v101, v103, base.Simd_g_const(&F_HFilter16i_SSE2__k13), base.Simd_g_const(&F_HFilter16i_SSE2__k14))
	v107 = l1 << (uint(v73) % 32)
	v108 = l0 + v107
	v113 = base.Simd_g_v128_load32_splat(m, v108, v6)
	v116 = base.Simd_g_v128_load32_lane_l1(m, v108+v62, v6, v113)
	v119 = base.Simd_g_v128_load32_lane_l2(m, v108+v59, v6, v116)
	v122 = base.Simd_g_v128_load32_lane_l3(m, v108+v56, v6, v119)
	v128 = base.Simd_g_v128_load32_splat(m, v108+l1, v6)
	v131 = base.Simd_g_v128_load32_lane_l1(m, v108+v82, v6, v128)
	v134 = base.Simd_g_v128_load32_lane_l2(m, v108+v79, v6, v131)
	v137 = base.Simd_g_v128_load32_lane_l3(m, v108+v76, v6, v134)
	v139 = base.Simd_g_i8x16_shuffle2(v122, v137, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2))
	v141 = base.Simd_g_i8x16_shuffle2(v122, v137, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5))
	v143 = base.Simd_g_i8x16_shuffle2(v139, v141, base.Simd_g_const(&F_HFilter16i_SSE2__k7), base.Simd_g_const(&F_HFilter16i_SSE2__k8))
	v145 = base.Simd_g_i8x16_shuffle2(v139, v141, base.Simd_g_const(&F_HFilter16i_SSE2__k10), base.Simd_g_const(&F_HFilter16i_SSE2__k11))
	v147 = base.Simd_g_i8x16_shuffle2(v143, v145, base.Simd_g_const(&F_HFilter16i_SSE2__k13), base.Simd_g_const(&F_HFilter16i_SSE2__k14))
	v148 = base.Simd_g_const(&F_HFilter16i_SSE2__k15)
	v150 = base.Simd_g_const(&F_HFilter16i_SSE2__k16)
	v152 = base.Simd_g_const(&F_HFilter16i_SSE2__k17)
	v153 = base.Simd_g_i8x16_shuffle2(v101, v103, base.Simd_g_const(&F_HFilter16i_SSE2__k18), base.Simd_g_const(&F_HFilter16i_SSE2__k19))
	v155 = base.Simd_g_i8x16_shuffle2(v143, v145, base.Simd_g_const(&F_HFilter16i_SSE2__k18), base.Simd_g_const(&F_HFilter16i_SSE2__k19))
	v178 = l0
	v190 = base.Simd_g_i8x16_shuffle2(v105, v147, base.Simd_g_const(&F_HFilter16i_SSE2__k20), base.Simd_g_const(&F_HFilter16i_SSE2__k21))
	v194 = int32(4)
	v195 = base.Simd_g_i8x16_shuffle2(v105, v147, base.Simd_g_const(&F_HFilter16i_SSE2__k22), base.Simd_g_const(&F_HFilter16i_SSE2__k23))
	v196 = base.Simd_g_i8x16_shuffle2(v153, v155, base.Simd_g_const(&F_HFilter16i_SSE2__k22), base.Simd_g_const(&F_HFilter16i_SSE2__k23))
	v198 = base.Simd_g_i8x16_shuffle2(v153, v155, base.Simd_g_const(&F_HFilter16i_SSE2__k20), base.Simd_g_const(&F_HFilter16i_SSE2__k21))
	for {
		v232 = base.Simd_g_const(&F_HFilter16i_SSE2__k24)
		v233 = base.Simd_g_v128_xor(v190, v232)
		v234 = base.Simd_g_const(&F_HFilter16i_SSE2__k25)
		v235 = v178 + v56
		v236 = int32(4)
		v238 = v178 + v59
		v241 = v178 + v62
		v245 = base.Simd_g_v128_load32_splat(m, v178, v236)
		v246 = int32(0)
		v247 = int32(1)
		v248 = base.Simd_g_v128_load32_lane_l1(m, v241+v236, v246, v245)
		v250 = int32(2)
		v251 = base.Simd_g_v128_load32_lane_l2(m, v238+v236, v246, v248)
		v253 = int32(3)
		v254 = base.Simd_g_v128_load32_lane_l3(m, v235+v236, v246, v251)
		v255 = v178 + v76
		v258 = v178 + v79
		v261 = v178 + v82
		v264 = v178 + l1
		v268 = base.Simd_g_v128_load32_splat(m, v264+v236, v246)
		v271 = base.Simd_g_v128_load32_lane_l1(m, v261+v236, v246, v268)
		v274 = base.Simd_g_v128_load32_lane_l2(m, v258+v236, v246, v271)
		v277 = base.Simd_g_v128_load32_lane_l3(m, v255+v236, v246, v274)
		v278 = base.Simd_g_const(&F_HFilter16i_SSE2__k0)
		v279 = base.Simd_g_i8x16_shuffle2(v254, v277, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2))
		v280 = base.Simd_g_const(&F_HFilter16i_SSE2__k3)
		v281 = base.Simd_g_i8x16_shuffle2(v254, v277, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5))
		v282 = base.Simd_g_const(&F_HFilter16i_SSE2__k6)
		v283 = base.Simd_g_i8x16_shuffle2(v279, v281, base.Simd_g_const(&F_HFilter16i_SSE2__k7), base.Simd_g_const(&F_HFilter16i_SSE2__k8))
		v284 = base.Simd_g_const(&F_HFilter16i_SSE2__k9)
		v285 = base.Simd_g_i8x16_shuffle2(v279, v281, base.Simd_g_const(&F_HFilter16i_SSE2__k10), base.Simd_g_const(&F_HFilter16i_SSE2__k11))
		v286 = base.Simd_g_const(&F_HFilter16i_SSE2__k17)
		v287 = base.Simd_g_i8x16_shuffle2(v283, v285, base.Simd_g_const(&F_HFilter16i_SSE2__k18), base.Simd_g_const(&F_HFilter16i_SSE2__k19))
		v288 = v178 + l1*int32(14)
		v291 = v178 + l1*int32(10)
		v294 = v178 + l1*int32(12)
		v297 = v178 + v107
		v301 = base.Simd_g_v128_load32_splat(m, v297+v236, v246)
		v304 = base.Simd_g_v128_load32_lane_l1(m, v294+v236, v246, v301)
		v307 = base.Simd_g_v128_load32_lane_l2(m, v291+v236, v246, v304)
		v310 = base.Simd_g_v128_load32_lane_l3(m, v288+v236, v246, v307)
		v311 = v178 + l1*int32(15)
		v314 = v178 + l1*int32(11)
		v317 = v178 + l1*int32(13)
		v320 = v178 + l1*int32(9)
		v324 = base.Simd_g_v128_load32_splat(m, v320+v236, v246)
		v327 = base.Simd_g_v128_load32_lane_l1(m, v317+v236, v246, v324)
		v330 = base.Simd_g_v128_load32_lane_l2(m, v314+v236, v246, v327)
		v333 = base.Simd_g_v128_load32_lane_l3(m, v311+v236, v246, v330)
		v335 = base.Simd_g_i8x16_shuffle2(v310, v333, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2))
		v337 = base.Simd_g_i8x16_shuffle2(v310, v333, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5))
		v339 = base.Simd_g_i8x16_shuffle2(v335, v337, base.Simd_g_const(&F_HFilter16i_SSE2__k7), base.Simd_g_const(&F_HFilter16i_SSE2__k8))
		v341 = base.Simd_g_i8x16_shuffle2(v335, v337, base.Simd_g_const(&F_HFilter16i_SSE2__k10), base.Simd_g_const(&F_HFilter16i_SSE2__k11))
		v343 = base.Simd_g_i8x16_shuffle2(v339, v341, base.Simd_g_const(&F_HFilter16i_SSE2__k18), base.Simd_g_const(&F_HFilter16i_SSE2__k19))
		v344 = base.Simd_g_const(&F_HFilter16i_SSE2__k15)
		v345 = base.Simd_g_i8x16_shuffle2(v287, v343, base.Simd_g_const(&F_HFilter16i_SSE2__k22), base.Simd_g_const(&F_HFilter16i_SSE2__k23))
		v346 = base.Simd_g_v128_xor(v345, v232)
		v350 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v190, v195), base.Simd_g_i8x16_sub_sat_u(v195, v190))
		v351 = base.Simd_g_const(&F_HFilter16i_SSE2__k16)
		v352 = base.Simd_g_i8x16_shuffle2(v287, v343, base.Simd_g_const(&F_HFilter16i_SSE2__k20), base.Simd_g_const(&F_HFilter16i_SSE2__k21))
		v355 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v352, v345), base.Simd_g_i8x16_sub_sat_u(v345, v352))
		v358 = base.Simd_g_i8x16_eq(v234, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v350, v355), base.Simd_g_i8x16_splat(l4)))
		v360 = base.Simd_g_v128_xor(v352, v232)
		v361 = base.Simd_g_v128_xor(v195, v232)
		v362 = base.Simd_g_i8x16_sub_sat_s(v360, v361)
		v368 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v195, v352), base.Simd_g_i8x16_sub_sat_u(v352, v195))
		v387 = base.Simd_g_const(&F_HFilter16i_SSE2__k12)
		v388 = base.Simd_g_i8x16_shuffle2(v283, v285, base.Simd_g_const(&F_HFilter16i_SSE2__k13), base.Simd_g_const(&F_HFilter16i_SSE2__k14))
		v390 = base.Simd_g_i8x16_shuffle2(v339, v341, base.Simd_g_const(&F_HFilter16i_SSE2__k13), base.Simd_g_const(&F_HFilter16i_SSE2__k14))
		v392 = base.Simd_g_i8x16_shuffle2(v388, v390, base.Simd_g_const(&F_HFilter16i_SSE2__k20), base.Simd_g_const(&F_HFilter16i_SSE2__k21))
		v394 = base.Simd_g_i8x16_shuffle2(v388, v390, base.Simd_g_const(&F_HFilter16i_SSE2__k22), base.Simd_g_const(&F_HFilter16i_SSE2__k23))
		v407 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_bitselect(v234, base.Simd_g_i8x16_sub_sat_s(v233, v346), v358), v362), v362), v362), v234, base.Simd_g_i8x16_eq(v234, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v368, v368), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v190, v345), base.Simd_g_i8x16_sub_sat_u(v345, v190)), v247), base.Simd_g_const(&F_HFilter16i_SSE2__k26))), base.Simd_g_i8x16_splat(l2)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v350, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v198, v196), base.Simd_g_i8x16_sub_sat_u(v196, v198))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v196, v190), base.Simd_g_i8x16_sub_sat_u(v190, v196))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v392, v394), base.Simd_g_i8x16_sub_sat_u(v394, v392))), v355), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v345, v392), base.Simd_g_i8x16_sub_sat_u(v392, v345))), base.Simd_g_i8x16_splat(l3)))))
		v409 = base.Simd_g_i8x16_add_sat_s(v407, base.Simd_g_const(&F_HFilter16i_SSE2__k27))
		v412 = int32(11)
		v418 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v234, v409, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2)), v412), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v234, v409, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5)), v412))
		v423 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(base.Simd_g_v128_xor(v418, v232), v234), base.Simd_g_const(&F_HFilter16i_SSE2__k28)), v234, v358)
		v425 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v233, v423), v232)
		v427 = base.Simd_g_i8x16_add_sat_s(v407, base.Simd_g_const(&F_HFilter16i_SSE2__k29))
		v438 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v361, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v234, v427, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2)), v412), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v234, v427, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5)), v412))), v232)
		v440 = base.Simd_g_i8x16_shuffle2(v425, v438, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2))
		v443 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v360, v418), v232)
		v445 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v346, v423), v232)
		v447 = base.Simd_g_i8x16_shuffle2(v443, v445, base.Simd_g_const(&F_HFilter16i_SSE2__k1), base.Simd_g_const(&F_HFilter16i_SSE2__k2))
		v449 = base.Simd_g_i8x16_shuffle2(v440, v447, base.Simd_g_const(&F_HFilter16i_SSE2__k7), base.Simd_g_const(&F_HFilter16i_SSE2__k8))
		base.Simd_g_v128_store32_lane_l0(m, v178, v250, v449)
		v455 = base.Simd_g_const(&F_HFilter16i_SSE2__k30)
		v456 = base.Simd_g_i8x16_shuffle2(v449, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v264+v250, v246, v456)
		v463 = base.Simd_g_i8x16_shuffle2(v456, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v238+v250, v246, v463)
		base.Simd_g_v128_store32_lane_l0(m, v258+v250, v246, base.Simd_g_i8x16_shuffle2(v463, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32)))
		v477 = base.Simd_g_i8x16_shuffle2(v440, v447, base.Simd_g_const(&F_HFilter16i_SSE2__k10), base.Simd_g_const(&F_HFilter16i_SSE2__k11))
		base.Simd_g_v128_store32_lane_l0(m, v241+v250, v246, v477)
		v484 = base.Simd_g_i8x16_shuffle2(v477, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v261+v250, v246, v484)
		v491 = base.Simd_g_i8x16_shuffle2(v484, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v235+v250, v246, v491)
		base.Simd_g_v128_store32_lane_l0(m, v255+v250, v246, base.Simd_g_i8x16_shuffle2(v491, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32)))
		v505 = base.Simd_g_i8x16_shuffle2(v425, v438, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5))
		v507 = base.Simd_g_i8x16_shuffle2(v443, v445, base.Simd_g_const(&F_HFilter16i_SSE2__k4), base.Simd_g_const(&F_HFilter16i_SSE2__k5))
		v509 = base.Simd_g_i8x16_shuffle2(v505, v507, base.Simd_g_const(&F_HFilter16i_SSE2__k7), base.Simd_g_const(&F_HFilter16i_SSE2__k8))
		base.Simd_g_v128_store32_lane_l0(m, v297+v250, v246, v509)
		v516 = base.Simd_g_i8x16_shuffle2(v509, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v320+v250, v246, v516)
		v523 = base.Simd_g_i8x16_shuffle2(v516, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v291+v250, v246, v523)
		base.Simd_g_v128_store32_lane_l0(m, v314+v250, v246, base.Simd_g_i8x16_shuffle2(v523, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32)))
		v537 = base.Simd_g_i8x16_shuffle2(v505, v507, base.Simd_g_const(&F_HFilter16i_SSE2__k10), base.Simd_g_const(&F_HFilter16i_SSE2__k11))
		base.Simd_g_v128_store32_lane_l0(m, v294+v250, v246, v537)
		v544 = base.Simd_g_i8x16_shuffle2(v537, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v317+v250, v246, v544)
		v551 = base.Simd_g_i8x16_shuffle2(v544, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32))
		base.Simd_g_v128_store32_lane_l0(m, v288+v250, v246, v551)
		base.Simd_g_v128_store32_lane_l0(m, v311+v250, v246, base.Simd_g_i8x16_shuffle2(v551, v234, base.Simd_g_const(&F_HFilter16i_SSE2__k31), base.Simd_g_const(&F_HFilter16i_SSE2__k32)))
		v565 = v194 + int32(-1)
		if base.Ui32(v247) < base.Ui32(v565) {
			v178 = v178 + v236
			v190 = v392
			v194 = v565
			v195 = v394
			v196 = v445
			v198 = v443
			continue
		} else {
			break
		}
		break
	}
	return
}

var F_HFilter16i_SSE2__k0 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_HFilter16i_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_HFilter16i_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_HFilter16i_SSE2__k3 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_HFilter16i_SSE2__k4 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_HFilter16i_SSE2__k5 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_HFilter16i_SSE2__k6 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_HFilter16i_SSE2__k7 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_HFilter16i_SSE2__k8 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_HFilter16i_SSE2__k9 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_HFilter16i_SSE2__k10 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_HFilter16i_SSE2__k11 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_HFilter16i_SSE2__k12 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_HFilter16i_SSE2__k13 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_HFilter16i_SSE2__k14 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_HFilter16i_SSE2__k15 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_HFilter16i_SSE2__k16 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_HFilter16i_SSE2__k17 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_HFilter16i_SSE2__k18 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_HFilter16i_SSE2__k19 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_HFilter16i_SSE2__k20 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_HFilter16i_SSE2__k21 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_HFilter16i_SSE2__k22 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_HFilter16i_SSE2__k23 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_HFilter16i_SSE2__k24 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_HFilter16i_SSE2__k25 = [2]uint64{0x0, 0x0}
var F_HFilter16i_SSE2__k26 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_HFilter16i_SSE2__k27 = [2]uint64{0x404040404040404, 0x404040404040404}
var F_HFilter16i_SSE2__k28 = [2]uint64{0xc0c0c0c0c0c0c0c0, 0xc0c0c0c0c0c0c0c0}
var F_HFilter16i_SSE2__k29 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_HFilter16i_SSE2__k30 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_HFilter16i_SSE2__k31 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_HFilter16i_SSE2__k32 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_HFilter8_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	var v7 int32
	_ = v7
	var v12 base.V128
	_ = v12
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v52 base.V128
	_ = v52
	var v55 base.V128
	_ = v55
	var v58 base.V128
	_ = v58
	var v60 int32
	_ = v60
	var v61 base.V128
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v73 base.V128
	_ = v73
	var v76 base.V128
	_ = v76
	var v79 base.V128
	_ = v79
	var v82 base.V128
	_ = v82
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v86 base.V128
	_ = v86
	var v87 base.V128
	_ = v87
	var v88 base.V128
	_ = v88
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v94 int32
	_ = v94
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
	var v112 int32
	_ = v112
	var v114 base.V128
	_ = v114
	var v117 base.V128
	_ = v117
	var v120 base.V128
	_ = v120
	var v123 base.V128
	_ = v123
	var v125 base.V128
	_ = v125
	var v127 base.V128
	_ = v127
	var v129 base.V128
	_ = v129
	var v131 base.V128
	_ = v131
	var v133 base.V128
	_ = v133
	var v134 base.V128
	_ = v134
	var v135 base.V128
	_ = v135
	var v136 base.V128
	_ = v136
	var v137 base.V128
	_ = v137
	var v138 base.V128
	_ = v138
	var v141 base.V128
	_ = v141
	var v142 base.V128
	_ = v142
	var v144 base.V128
	_ = v144
	var v146 base.V128
	_ = v146
	var v149 int32
	_ = v149
	var v151 base.V128
	_ = v151
	var v154 base.V128
	_ = v154
	var v157 base.V128
	_ = v157
	var v160 base.V128
	_ = v160
	var v164 int32
	_ = v164
	var v166 base.V128
	_ = v166
	var v169 base.V128
	_ = v169
	var v172 base.V128
	_ = v172
	var v175 base.V128
	_ = v175
	var v177 base.V128
	_ = v177
	var v179 base.V128
	_ = v179
	var v181 base.V128
	_ = v181
	var v183 base.V128
	_ = v183
	var v185 base.V128
	_ = v185
	var v188 int32
	_ = v188
	var v190 base.V128
	_ = v190
	var v193 base.V128
	_ = v193
	var v196 base.V128
	_ = v196
	var v199 base.V128
	_ = v199
	var v203 int32
	_ = v203
	var v205 base.V128
	_ = v205
	var v208 base.V128
	_ = v208
	var v211 base.V128
	_ = v211
	var v214 base.V128
	_ = v214
	var v216 base.V128
	_ = v216
	var v218 base.V128
	_ = v218
	var v220 base.V128
	_ = v220
	var v222 base.V128
	_ = v222
	var v224 base.V128
	_ = v224
	var v226 base.V128
	_ = v226
	var v229 base.V128
	_ = v229
	var v232 base.V128
	_ = v232
	var v234 base.V128
	_ = v234
	var v247 base.V128
	_ = v247
	var v258 base.V128
	_ = v258
	var v261 base.V128
	_ = v261
	var v263 base.V128
	_ = v263
	var v265 base.V128
	_ = v265
	var v267 base.V128
	_ = v267
	var v279 base.V128
	_ = v279
	var v280 base.V128
	_ = v280
	var v281 base.V128
	_ = v281
	var v282 base.V128
	_ = v282
	var v283 base.V128
	_ = v283
	var v284 base.V128
	_ = v284
	var v288 base.V128
	_ = v288
	var v293 base.V128
	_ = v293
	var v294 base.V128
	_ = v294
	var v296 base.V128
	_ = v296
	var v298 base.V128
	_ = v298
	var v302 base.V128
	_ = v302
	var v303 base.V128
	_ = v303
	var v304 base.V128
	_ = v304
	var v305 base.V128
	_ = v305
	var v309 base.V128
	_ = v309
	var v315 base.V128
	_ = v315
	var v316 base.V128
	_ = v316
	var v319 base.V128
	_ = v319
	var v321 base.V128
	_ = v321
	var v323 base.V128
	_ = v323
	var v324 base.V128
	_ = v324
	var v327 base.V128
	_ = v327
	var v330 base.V128
	_ = v330
	var v332 base.V128
	_ = v332
	var v334 base.V128
	_ = v334
	var v336 base.V128
	_ = v336
	var v339 int32
	_ = v339
	var v353 base.V128
	_ = v353
	var v355 base.V128
	_ = v355
	var v357 base.V128
	_ = v357
	var v359 base.V128
	_ = v359
	var v364 base.V128
	_ = v364
	var v366 base.V128
	_ = v366
	var v368 base.V128
	_ = v368
	var v372 base.V128
	_ = v372
	var v373 base.V128
	_ = v373
	var v378 base.V128
	_ = v378
	var v382 int32
	_ = v382
	var v384 base.V128
	_ = v384
	var v388 int32
	_ = v388
	var v390 base.V128
	_ = v390
	var v407 base.V128
	_ = v407
	var v412 base.V128
	_ = v412
	var v416 int32
	_ = v416
	var v418 base.V128
	_ = v418
	var v422 int32
	_ = v422
	var v424 base.V128
	_ = v424
	var v428 int32
	_ = v428
	var v430 base.V128
	_ = v430
	var v434 int32
	_ = v434
	var v436 base.V128
	_ = v436
	var v453 base.V128
	_ = v453
	var v465 base.V128
	_ = v465
	var v467 base.V128
	_ = v467
	var v469 base.V128
	_ = v469
	var v472 base.V128
	_ = v472
	var v474 base.V128
	_ = v474
	var v476 base.V128
	_ = v476
	var v481 base.V128
	_ = v481
	var v483 base.V128
	_ = v483
	var v485 base.V128
	_ = v485
	var v490 base.V128
	_ = v490
	var v495 base.V128
	_ = v495
	var v499 int32
	_ = v499
	var v501 base.V128
	_ = v501
	var v505 int32
	_ = v505
	var v507 base.V128
	_ = v507
	var v524 base.V128
	_ = v524
	var v529 base.V128
	_ = v529
	var v533 int32
	_ = v533
	var v535 base.V128
	_ = v535
	var v539 int32
	_ = v539
	var v541 base.V128
	_ = v541
	var v545 int32
	_ = v545
	var v547 base.V128
	_ = v547
	var v551 int32
	_ = v551
	var v553 base.V128
	_ = v553
	v7 = int32(0)
	v12 = base.Simd_g_const(&F_HFilter8_SSE2__k0)
	v40 = int32(-4)
	v41 = l0 + v40
	v43 = l2 * int32(6)
	v45 = int32(1)
	v46 = l2 << (uint(v45) % 32)
	v48 = int32(2)
	v49 = l2 << (uint(v48) % 32)
	v50 = v41 + v49
	v52 = base.Simd_g_v128_load32_splat(m, v41, v7)
	v55 = base.Simd_g_v128_load32_lane_l1(m, v50, v7, v52)
	v58 = base.Simd_g_v128_load32_lane_l2(m, v41+v46, v7, v55)
	v60 = int32(3)
	v61 = base.Simd_g_v128_load32_lane_l3(m, v41+v43, v7, v58)
	v62 = int32(7)
	v63 = l2 * v62
	v66 = l2 * v60
	v69 = l2 * int32(5)
	v71 = v41 + l2
	v73 = base.Simd_g_v128_load32_splat(m, v71, v7)
	v76 = base.Simd_g_v128_load32_lane_l1(m, v41+v69, v7, v73)
	v79 = base.Simd_g_v128_load32_lane_l2(m, v41+v66, v7, v76)
	v82 = base.Simd_g_v128_load32_lane_l3(m, v41+v63, v7, v79)
	v83 = base.Simd_g_const(&F_HFilter8_SSE2__k1)
	v84 = base.Simd_g_i8x16_shuffle2(v61, v82, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v85 = base.Simd_g_const(&F_HFilter8_SSE2__k4)
	v86 = base.Simd_g_i8x16_shuffle2(v61, v82, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v87 = base.Simd_g_const(&F_HFilter8_SSE2__k7)
	v88 = base.Simd_g_i8x16_shuffle2(v84, v86, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	v89 = base.Simd_g_const(&F_HFilter8_SSE2__k10)
	v90 = base.Simd_g_i8x16_shuffle2(v84, v86, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	v91 = base.Simd_g_const(&F_HFilter8_SSE2__k13)
	v92 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_HFilter8_SSE2__k14), base.Simd_g_const(&F_HFilter8_SSE2__k15))
	v94 = l1 + v40
	v97 = v94 + v49
	v99 = base.Simd_g_v128_load32_splat(m, v94, v7)
	v102 = base.Simd_g_v128_load32_lane_l1(m, v97, v7, v99)
	v105 = base.Simd_g_v128_load32_lane_l2(m, v94+v46, v7, v102)
	v108 = base.Simd_g_v128_load32_lane_l3(m, v94+v43, v7, v105)
	v112 = v94 + l2
	v114 = base.Simd_g_v128_load32_splat(m, v112, v7)
	v117 = base.Simd_g_v128_load32_lane_l1(m, v94+v69, v7, v114)
	v120 = base.Simd_g_v128_load32_lane_l2(m, v94+v66, v7, v117)
	v123 = base.Simd_g_v128_load32_lane_l3(m, v94+v63, v7, v120)
	v125 = base.Simd_g_i8x16_shuffle2(v108, v123, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v127 = base.Simd_g_i8x16_shuffle2(v108, v123, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v129 = base.Simd_g_i8x16_shuffle2(v125, v127, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	v131 = base.Simd_g_i8x16_shuffle2(v125, v127, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	v133 = base.Simd_g_i8x16_shuffle2(v129, v131, base.Simd_g_const(&F_HFilter8_SSE2__k14), base.Simd_g_const(&F_HFilter8_SSE2__k15))
	v134 = base.Simd_g_const(&F_HFilter8_SSE2__k16)
	v135 = base.Simd_g_i8x16_shuffle2(v92, v133, base.Simd_g_const(&F_HFilter8_SSE2__k17), base.Simd_g_const(&F_HFilter8_SSE2__k18))
	v136 = base.Simd_g_const(&F_HFilter8_SSE2__k19)
	v137 = base.Simd_g_i8x16_shuffle2(v92, v133, base.Simd_g_const(&F_HFilter8_SSE2__k20), base.Simd_g_const(&F_HFilter8_SSE2__k21))
	v138 = base.Simd_g_const(&F_HFilter8_SSE2__k22)
	v141 = base.Simd_g_const(&F_HFilter8_SSE2__k23)
	v142 = base.Simd_g_i8x16_shuffle2(v88, v90, base.Simd_g_const(&F_HFilter8_SSE2__k24), base.Simd_g_const(&F_HFilter8_SSE2__k25))
	v144 = base.Simd_g_i8x16_shuffle2(v129, v131, base.Simd_g_const(&F_HFilter8_SSE2__k24), base.Simd_g_const(&F_HFilter8_SSE2__k25))
	v146 = base.Simd_g_i8x16_shuffle2(v142, v144, base.Simd_g_const(&F_HFilter8_SSE2__k20), base.Simd_g_const(&F_HFilter8_SSE2__k21))
	v149 = l0 + v49
	v151 = base.Simd_g_v128_load32_splat(m, l0, v7)
	v154 = base.Simd_g_v128_load32_lane_l1(m, v149, v7, v151)
	v157 = base.Simd_g_v128_load32_lane_l2(m, l0+v46, v7, v154)
	v160 = base.Simd_g_v128_load32_lane_l3(m, l0+v43, v7, v157)
	v164 = l0 + l2
	v166 = base.Simd_g_v128_load32_splat(m, v164, v7)
	v169 = base.Simd_g_v128_load32_lane_l1(m, l0+v69, v7, v166)
	v172 = base.Simd_g_v128_load32_lane_l2(m, l0+v66, v7, v169)
	v175 = base.Simd_g_v128_load32_lane_l3(m, l0+v63, v7, v172)
	v177 = base.Simd_g_i8x16_shuffle2(v160, v175, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v179 = base.Simd_g_i8x16_shuffle2(v160, v175, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v181 = base.Simd_g_i8x16_shuffle2(v177, v179, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	v183 = base.Simd_g_i8x16_shuffle2(v177, v179, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	v185 = base.Simd_g_i8x16_shuffle2(v181, v183, base.Simd_g_const(&F_HFilter8_SSE2__k14), base.Simd_g_const(&F_HFilter8_SSE2__k15))
	v188 = l1 + v49
	v190 = base.Simd_g_v128_load32_splat(m, l1, v7)
	v193 = base.Simd_g_v128_load32_lane_l1(m, v188, v7, v190)
	v196 = base.Simd_g_v128_load32_lane_l2(m, l1+v46, v7, v193)
	v199 = base.Simd_g_v128_load32_lane_l3(m, l1+v43, v7, v196)
	v203 = l1 + l2
	v205 = base.Simd_g_v128_load32_splat(m, v203, v7)
	v208 = base.Simd_g_v128_load32_lane_l1(m, l1+v69, v7, v205)
	v211 = base.Simd_g_v128_load32_lane_l2(m, l1+v66, v7, v208)
	v214 = base.Simd_g_v128_load32_lane_l3(m, l1+v63, v7, v211)
	v216 = base.Simd_g_i8x16_shuffle2(v199, v214, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v218 = base.Simd_g_i8x16_shuffle2(v199, v214, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v220 = base.Simd_g_i8x16_shuffle2(v216, v218, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	v222 = base.Simd_g_i8x16_shuffle2(v216, v218, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	v224 = base.Simd_g_i8x16_shuffle2(v220, v222, base.Simd_g_const(&F_HFilter8_SSE2__k14), base.Simd_g_const(&F_HFilter8_SSE2__k15))
	v226 = base.Simd_g_i8x16_shuffle2(v185, v224, base.Simd_g_const(&F_HFilter8_SSE2__k17), base.Simd_g_const(&F_HFilter8_SSE2__k18))
	v229 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v146, v226), base.Simd_g_i8x16_sub_sat_u(v226, v146))
	v232 = base.Simd_g_i8x16_shuffle2(v142, v144, base.Simd_g_const(&F_HFilter8_SSE2__k17), base.Simd_g_const(&F_HFilter8_SSE2__k18))
	v234 = base.Simd_g_i8x16_shuffle2(v185, v224, base.Simd_g_const(&F_HFilter8_SSE2__k20), base.Simd_g_const(&F_HFilter8_SSE2__k21))
	v247 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v232, v146), base.Simd_g_i8x16_sub_sat_u(v146, v232))
	v258 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v234, v226), base.Simd_g_i8x16_sub_sat_u(v226, v234))
	v261 = base.Simd_g_i8x16_shuffle2(v181, v183, base.Simd_g_const(&F_HFilter8_SSE2__k24), base.Simd_g_const(&F_HFilter8_SSE2__k25))
	v263 = base.Simd_g_i8x16_shuffle2(v220, v222, base.Simd_g_const(&F_HFilter8_SSE2__k24), base.Simd_g_const(&F_HFilter8_SSE2__k25))
	v265 = base.Simd_g_i8x16_shuffle2(v261, v263, base.Simd_g_const(&F_HFilter8_SSE2__k20), base.Simd_g_const(&F_HFilter8_SSE2__k21))
	v267 = base.Simd_g_i8x16_shuffle2(v261, v263, base.Simd_g_const(&F_HFilter8_SSE2__k17), base.Simd_g_const(&F_HFilter8_SSE2__k18))
	v279 = base.Simd_g_i8x16_eq(v12, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v229, v229), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v232, v234), base.Simd_g_i8x16_sub_sat_u(v234, v232)), v45), base.Simd_g_const(&F_HFilter8_SSE2__k26))), base.Simd_g_i8x16_splat(l3)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v247, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v135, v137), base.Simd_g_i8x16_sub_sat_u(v137, v135))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v137, v232), base.Simd_g_i8x16_sub_sat_u(v232, v137))), v258), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v265, v267), base.Simd_g_i8x16_sub_sat_u(v267, v265))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v267, v234), base.Simd_g_i8x16_sub_sat_u(v234, v267))), base.Simd_g_i8x16_splat(l4))))
	v280 = base.Simd_g_v128_xor(v226, v138)
	v281 = base.Simd_g_v128_xor(v146, v138)
	v282 = base.Simd_g_i8x16_sub_sat_s(v280, v281)
	v283 = base.Simd_g_v128_xor(v232, v138)
	v284 = base.Simd_g_v128_xor(v234, v138)
	v288 = base.Simd_g_i8x16_add_sat_s(v282, base.Simd_g_i8x16_add_sat_s(v282, base.Simd_g_i8x16_add_sat_s(v282, base.Simd_g_i8x16_sub_sat_s(v283, v284))))
	v293 = base.Simd_g_i8x16_eq(v12, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v247, v258), base.Simd_g_i8x16_splat(l5)))
	v294 = base.Simd_g_v128_and(base.Simd_g_v128_and(v279, v288), v293)
	v296 = base.Simd_g_i8x16_shuffle2(v12, v294, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v298 = base.Simd_g_const(&F_HFilter8_SSE2__k27)
	v302 = base.Simd_g_const(&F_HFilter8_SSE2__k28)
	v303 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v296), v298), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v296), v298), base.Simd_g_const(&F_HFilter8_SSE2__k29), base.Simd_g_const(&F_HFilter8_SSE2__k30))
	v304 = base.Simd_g_const(&F_HFilter8_SSE2__k31)
	v305 = base.Simd_g_i16x8_add(v303, v304)
	v309 = base.Simd_g_i8x16_shuffle2(v12, v294, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v315 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_s(v309), v298), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_s(v309), v298), base.Simd_g_const(&F_HFilter8_SSE2__k29), base.Simd_g_const(&F_HFilter8_SSE2__k30))
	v316 = base.Simd_g_i16x8_add(v315, v304)
	v319 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v305, v62), base.Simd_g_i16x8_shr_s(v316, v62))
	v321 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_xor(v137, v138), v319), v138)
	v323 = base.Simd_g_i8x16_shuffle2(v135, v321, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v324 = base.Simd_g_i16x8_add(v305, v303)
	v327 = base.Simd_g_i16x8_add(v316, v315)
	v330 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(v324, v62), base.Simd_g_i16x8_shr_s(v327, v62))
	v332 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v283, v330), v138)
	v334 = base.Simd_g_v128_and(base.Simd_g_v128_andnot(v288, v293), v279)
	v336 = base.Simd_g_i8x16_add_sat_s(v334, base.Simd_g_const(&F_HFilter8_SSE2__k32))
	v339 = int32(11)
	v353 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v324, v303), v62), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(v327, v315), v62))
	v355 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(v281, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v12, v336, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3)), v339), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v12, v336, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6)), v339))), v353), v138)
	v357 = base.Simd_g_i8x16_shuffle2(v332, v355, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v359 = base.Simd_g_i8x16_shuffle2(v323, v357, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v41, v7, v359)
	v364 = base.Simd_g_i8x16_shuffle2(v135, v321, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v366 = base.Simd_g_i8x16_shuffle2(v332, v355, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v368 = base.Simd_g_i8x16_shuffle2(v364, v366, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, v94, v7, v368)
	v372 = base.Simd_g_const(&F_HFilter8_SSE2__k33)
	v373 = base.Simd_g_i8x16_shuffle2(v359, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v71, v7, v373)
	v378 = base.Simd_g_i8x16_shuffle2(v368, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v112, v7, v378)
	v382 = v71 + l2
	v384 = base.Simd_g_i8x16_shuffle2(v373, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v382, v7, v384)
	v388 = v112 + l2
	v390 = base.Simd_g_i8x16_shuffle2(v378, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v388, v7, v390)
	base.Simd_g_v128_store32_lane_l0(m, v382+l2, v7, base.Simd_g_i8x16_shuffle2(v384, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	base.Simd_g_v128_store32_lane_l0(m, v388+l2, v7, base.Simd_g_i8x16_shuffle2(v390, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	v407 = base.Simd_g_i8x16_shuffle2(v323, v357, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v50, v7, v407)
	v412 = base.Simd_g_i8x16_shuffle2(v364, v366, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v97, v7, v412)
	v416 = v50 + l2
	v418 = base.Simd_g_i8x16_shuffle2(v407, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v416, v7, v418)
	v422 = v97 + l2
	v424 = base.Simd_g_i8x16_shuffle2(v412, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v422, v7, v424)
	v428 = v416 + l2
	v430 = base.Simd_g_i8x16_shuffle2(v418, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v428, v7, v430)
	v434 = v422 + l2
	v436 = base.Simd_g_i8x16_shuffle2(v424, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v434, v7, v436)
	base.Simd_g_v128_store32_lane_l0(m, v428+l2, v7, base.Simd_g_i8x16_shuffle2(v430, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	base.Simd_g_v128_store32_lane_l0(m, v434+l2, v7, base.Simd_g_i8x16_shuffle2(v436, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	v453 = base.Simd_g_i8x16_add_sat_s(v334, base.Simd_g_const(&F_HFilter8_SSE2__k36))
	v465 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_i8x16_sub_sat_s(v280, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v12, v453, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3)), v339), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v12, v453, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6)), v339))), v353), v138)
	v467 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v284, v330), v138)
	v469 = base.Simd_g_i8x16_shuffle2(v465, v467, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v472 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(base.Simd_g_v128_xor(v267, v138), v319), v138)
	v474 = base.Simd_g_i8x16_shuffle2(v472, v265, base.Simd_g_const(&F_HFilter8_SSE2__k2), base.Simd_g_const(&F_HFilter8_SSE2__k3))
	v476 = base.Simd_g_i8x16_shuffle2(v469, v474, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, l0, v7, v476)
	v481 = base.Simd_g_i8x16_shuffle2(v465, v467, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v483 = base.Simd_g_i8x16_shuffle2(v472, v265, base.Simd_g_const(&F_HFilter8_SSE2__k5), base.Simd_g_const(&F_HFilter8_SSE2__k6))
	v485 = base.Simd_g_i8x16_shuffle2(v481, v483, base.Simd_g_const(&F_HFilter8_SSE2__k8), base.Simd_g_const(&F_HFilter8_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, l1, v7, v485)
	v490 = base.Simd_g_i8x16_shuffle2(v476, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v164, v7, v490)
	v495 = base.Simd_g_i8x16_shuffle2(v485, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v203, v7, v495)
	v499 = v164 + l2
	v501 = base.Simd_g_i8x16_shuffle2(v490, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v499, v7, v501)
	v505 = v203 + l2
	v507 = base.Simd_g_i8x16_shuffle2(v495, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v505, v7, v507)
	base.Simd_g_v128_store32_lane_l0(m, v499+l2, v7, base.Simd_g_i8x16_shuffle2(v501, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	base.Simd_g_v128_store32_lane_l0(m, v505+l2, v7, base.Simd_g_i8x16_shuffle2(v507, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	v524 = base.Simd_g_i8x16_shuffle2(v469, v474, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v149, v7, v524)
	v529 = base.Simd_g_i8x16_shuffle2(v481, v483, base.Simd_g_const(&F_HFilter8_SSE2__k11), base.Simd_g_const(&F_HFilter8_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v188, v7, v529)
	v533 = v149 + l2
	v535 = base.Simd_g_i8x16_shuffle2(v524, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v533, v7, v535)
	v539 = v188 + l2
	v541 = base.Simd_g_i8x16_shuffle2(v529, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v539, v7, v541)
	v545 = v533 + l2
	v547 = base.Simd_g_i8x16_shuffle2(v535, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v545, v7, v547)
	v551 = v539 + l2
	v553 = base.Simd_g_i8x16_shuffle2(v541, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35))
	base.Simd_g_v128_store32_lane_l0(m, v551, v7, v553)
	base.Simd_g_v128_store32_lane_l0(m, v545+l2, v7, base.Simd_g_i8x16_shuffle2(v547, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	base.Simd_g_v128_store32_lane_l0(m, v551+l2, v7, base.Simd_g_i8x16_shuffle2(v553, v12, base.Simd_g_const(&F_HFilter8_SSE2__k34), base.Simd_g_const(&F_HFilter8_SSE2__k35)))
	return
}

var F_HFilter8_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_HFilter8_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_HFilter8_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_HFilter8_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_HFilter8_SSE2__k4 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_HFilter8_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_HFilter8_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_HFilter8_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_HFilter8_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_HFilter8_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_HFilter8_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_HFilter8_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_HFilter8_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_HFilter8_SSE2__k13 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_HFilter8_SSE2__k14 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_HFilter8_SSE2__k15 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_HFilter8_SSE2__k16 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_HFilter8_SSE2__k17 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_HFilter8_SSE2__k18 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_HFilter8_SSE2__k19 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_HFilter8_SSE2__k20 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_HFilter8_SSE2__k21 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_HFilter8_SSE2__k22 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_HFilter8_SSE2__k23 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_HFilter8_SSE2__k24 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_HFilter8_SSE2__k25 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_HFilter8_SSE2__k26 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_HFilter8_SSE2__k27 = [2]uint64{0x90000000900, 0x90000000900}
var F_HFilter8_SSE2__k28 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_HFilter8_SSE2__k29 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_HFilter8_SSE2__k30 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_HFilter8_SSE2__k31 = [2]uint64{0x3f003f003f003f, 0x3f003f003f003f}
var F_HFilter8_SSE2__k32 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_HFilter8_SSE2__k33 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_HFilter8_SSE2__k34 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_HFilter8_SSE2__k35 = [2]uint64{0x8080808080808080, 0x302010080808080}
var F_HFilter8_SSE2__k36 = [2]uint64{0x404040404040404, 0x404040404040404}

func F_HFilter8i_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	var v7 int32
	_ = v7
	var v10 base.V128
	_ = v10
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v45 base.V128
	_ = v45
	var v48 base.V128
	_ = v48
	var v51 base.V128
	_ = v51
	var v53 int32
	_ = v53
	var v54 base.V128
	_ = v54
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v66 base.V128
	_ = v66
	var v69 base.V128
	_ = v69
	var v72 base.V128
	_ = v72
	var v75 base.V128
	_ = v75
	var v76 base.V128
	_ = v76
	var v77 base.V128
	_ = v77
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
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v85 base.V128
	_ = v85
	var v90 base.V128
	_ = v90
	var v93 base.V128
	_ = v93
	var v96 base.V128
	_ = v96
	var v99 base.V128
	_ = v99
	var v105 base.V128
	_ = v105
	var v108 base.V128
	_ = v108
	var v111 base.V128
	_ = v111
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v118 base.V128
	_ = v118
	var v120 base.V128
	_ = v120
	var v122 base.V128
	_ = v122
	var v124 base.V128
	_ = v124
	var v125 base.V128
	_ = v125
	var v126 base.V128
	_ = v126
	var v127 base.V128
	_ = v127
	var v128 base.V128
	_ = v128
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v136 base.V128
	_ = v136
	var v139 base.V128
	_ = v139
	var v142 base.V128
	_ = v142
	var v145 base.V128
	_ = v145
	var v151 base.V128
	_ = v151
	var v154 base.V128
	_ = v154
	var v157 base.V128
	_ = v157
	var v160 base.V128
	_ = v160
	var v162 base.V128
	_ = v162
	var v164 base.V128
	_ = v164
	var v166 base.V128
	_ = v166
	var v168 base.V128
	_ = v168
	var v169 base.V128
	_ = v169
	var v170 base.V128
	_ = v170
	var v172 int32
	_ = v172
	var v177 base.V128
	_ = v177
	var v180 base.V128
	_ = v180
	var v183 base.V128
	_ = v183
	var v186 base.V128
	_ = v186
	var v192 base.V128
	_ = v192
	var v195 base.V128
	_ = v195
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v203 base.V128
	_ = v203
	var v205 base.V128
	_ = v205
	var v207 base.V128
	_ = v207
	var v209 base.V128
	_ = v209
	var v211 base.V128
	_ = v211
	var v212 base.V128
	_ = v212
	var v213 base.V128
	_ = v213
	var v214 base.V128
	_ = v214
	var v217 base.V128
	_ = v217
	var v220 base.V128
	_ = v220
	var v222 base.V128
	_ = v222
	var v225 base.V128
	_ = v225
	var v229 base.V128
	_ = v229
	var v231 base.V128
	_ = v231
	var v232 base.V128
	_ = v232
	var v233 base.V128
	_ = v233
	var v239 base.V128
	_ = v239
	var v252 base.V128
	_ = v252
	var v254 base.V128
	_ = v254
	var v256 base.V128
	_ = v256
	var v258 base.V128
	_ = v258
	var v269 base.V128
	_ = v269
	var v271 base.V128
	_ = v271
	var v273 base.V128
	_ = v273
	var v275 base.V128
	_ = v275
	var v288 base.V128
	_ = v288
	var v290 base.V128
	_ = v290
	var v293 int32
	_ = v293
	var v299 base.V128
	_ = v299
	var v304 base.V128
	_ = v304
	var v306 base.V128
	_ = v306
	var v308 base.V128
	_ = v308
	var v319 base.V128
	_ = v319
	var v321 base.V128
	_ = v321
	var v323 base.V128
	_ = v323
	var v325 base.V128
	_ = v325
	var v327 base.V128
	_ = v327
	var v329 base.V128
	_ = v329
	var v334 base.V128
	_ = v334
	var v336 base.V128
	_ = v336
	var v338 base.V128
	_ = v338
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v345 base.V128
	_ = v345
	var v346 base.V128
	_ = v346
	var v351 int32
	_ = v351
	var v352 int32
	_ = v352
	var v354 base.V128
	_ = v354
	var v358 int32
	_ = v358
	var v360 base.V128
	_ = v360
	var v364 int32
	_ = v364
	var v366 base.V128
	_ = v366
	var v382 int32
	_ = v382
	var v384 base.V128
	_ = v384
	var v388 int32
	_ = v388
	var v390 base.V128
	_ = v390
	var v394 int32
	_ = v394
	var v396 base.V128
	_ = v396
	var v400 int32
	_ = v400
	var v402 base.V128
	_ = v402
	var v406 int32
	_ = v406
	var v408 base.V128
	_ = v408
	var v412 int32
	_ = v412
	var v414 base.V128
	_ = v414
	v7 = int32(0)
	v10 = base.Simd_g_const(&F_HFilter8i_SSE2__k0)
	v36 = l2 * int32(6)
	v38 = int32(1)
	v39 = l2 << (uint(v38) % 32)
	v41 = int32(2)
	v42 = l2 << (uint(v41) % 32)
	v45 = base.Simd_g_v128_load32_splat(m, l0, v7)
	v48 = base.Simd_g_v128_load32_lane_l1(m, l0+v42, v7, v45)
	v51 = base.Simd_g_v128_load32_lane_l2(m, l0+v39, v7, v48)
	v53 = int32(3)
	v54 = base.Simd_g_v128_load32_lane_l3(m, l0+v36, v7, v51)
	v56 = l2 * int32(7)
	v59 = l2 * v53
	v62 = l2 * int32(5)
	v66 = base.Simd_g_v128_load32_splat(m, l0+l2, v7)
	v69 = base.Simd_g_v128_load32_lane_l1(m, l0+v62, v7, v66)
	v72 = base.Simd_g_v128_load32_lane_l2(m, l0+v59, v7, v69)
	v75 = base.Simd_g_v128_load32_lane_l3(m, l0+v56, v7, v72)
	v76 = base.Simd_g_const(&F_HFilter8i_SSE2__k1)
	v77 = base.Simd_g_i8x16_shuffle2(v54, v75, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3))
	v78 = base.Simd_g_const(&F_HFilter8i_SSE2__k4)
	v79 = base.Simd_g_i8x16_shuffle2(v54, v75, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6))
	v80 = base.Simd_g_const(&F_HFilter8i_SSE2__k7)
	v81 = base.Simd_g_i8x16_shuffle2(v77, v79, base.Simd_g_const(&F_HFilter8i_SSE2__k8), base.Simd_g_const(&F_HFilter8i_SSE2__k9))
	v82 = base.Simd_g_const(&F_HFilter8i_SSE2__k10)
	v83 = base.Simd_g_i8x16_shuffle2(v77, v79, base.Simd_g_const(&F_HFilter8i_SSE2__k11), base.Simd_g_const(&F_HFilter8i_SSE2__k12))
	v84 = base.Simd_g_const(&F_HFilter8i_SSE2__k13)
	v85 = base.Simd_g_i8x16_shuffle2(v81, v83, base.Simd_g_const(&F_HFilter8i_SSE2__k14), base.Simd_g_const(&F_HFilter8i_SSE2__k15))
	v90 = base.Simd_g_v128_load32_splat(m, l1, v7)
	v93 = base.Simd_g_v128_load32_lane_l1(m, l1+v42, v7, v90)
	v96 = base.Simd_g_v128_load32_lane_l2(m, l1+v39, v7, v93)
	v99 = base.Simd_g_v128_load32_lane_l3(m, l1+v36, v7, v96)
	v105 = base.Simd_g_v128_load32_splat(m, l1+l2, v7)
	v108 = base.Simd_g_v128_load32_lane_l1(m, l1+v62, v7, v105)
	v111 = base.Simd_g_v128_load32_lane_l2(m, l1+v59, v7, v108)
	v114 = base.Simd_g_v128_load32_lane_l3(m, l1+v56, v7, v111)
	v116 = base.Simd_g_i8x16_shuffle2(v99, v114, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3))
	v118 = base.Simd_g_i8x16_shuffle2(v99, v114, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6))
	v120 = base.Simd_g_i8x16_shuffle2(v116, v118, base.Simd_g_const(&F_HFilter8i_SSE2__k8), base.Simd_g_const(&F_HFilter8i_SSE2__k9))
	v122 = base.Simd_g_i8x16_shuffle2(v116, v118, base.Simd_g_const(&F_HFilter8i_SSE2__k11), base.Simd_g_const(&F_HFilter8i_SSE2__k12))
	v124 = base.Simd_g_i8x16_shuffle2(v120, v122, base.Simd_g_const(&F_HFilter8i_SSE2__k14), base.Simd_g_const(&F_HFilter8i_SSE2__k15))
	v125 = base.Simd_g_const(&F_HFilter8i_SSE2__k16)
	v126 = base.Simd_g_i8x16_shuffle2(v85, v124, base.Simd_g_const(&F_HFilter8i_SSE2__k17), base.Simd_g_const(&F_HFilter8i_SSE2__k18))
	v127 = base.Simd_g_const(&F_HFilter8i_SSE2__k19)
	v128 = base.Simd_g_v128_xor(v126, v127)
	v130 = int32(4)
	v131 = l0 + v130
	v136 = base.Simd_g_v128_load32_splat(m, l0, v130)
	v139 = base.Simd_g_v128_load32_lane_l1(m, v131+v42, v7, v136)
	v142 = base.Simd_g_v128_load32_lane_l2(m, v131+v39, v7, v139)
	v145 = base.Simd_g_v128_load32_lane_l3(m, v131+v36, v7, v142)
	v151 = base.Simd_g_v128_load32_splat(m, v131+l2, v7)
	v154 = base.Simd_g_v128_load32_lane_l1(m, v131+v62, v7, v151)
	v157 = base.Simd_g_v128_load32_lane_l2(m, v131+v59, v7, v154)
	v160 = base.Simd_g_v128_load32_lane_l3(m, v131+v56, v7, v157)
	v162 = base.Simd_g_i8x16_shuffle2(v145, v160, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3))
	v164 = base.Simd_g_i8x16_shuffle2(v145, v160, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6))
	v166 = base.Simd_g_i8x16_shuffle2(v162, v164, base.Simd_g_const(&F_HFilter8i_SSE2__k8), base.Simd_g_const(&F_HFilter8i_SSE2__k9))
	v168 = base.Simd_g_i8x16_shuffle2(v162, v164, base.Simd_g_const(&F_HFilter8i_SSE2__k11), base.Simd_g_const(&F_HFilter8i_SSE2__k12))
	v169 = base.Simd_g_const(&F_HFilter8i_SSE2__k20)
	v170 = base.Simd_g_i8x16_shuffle2(v166, v168, base.Simd_g_const(&F_HFilter8i_SSE2__k21), base.Simd_g_const(&F_HFilter8i_SSE2__k22))
	v172 = l1 + v130
	v177 = base.Simd_g_v128_load32_splat(m, l1, v130)
	v180 = base.Simd_g_v128_load32_lane_l1(m, v172+v42, v7, v177)
	v183 = base.Simd_g_v128_load32_lane_l2(m, v172+v39, v7, v180)
	v186 = base.Simd_g_v128_load32_lane_l3(m, v172+v36, v7, v183)
	v192 = base.Simd_g_v128_load32_splat(m, v172+l2, v7)
	v195 = base.Simd_g_v128_load32_lane_l1(m, v172+v62, v7, v192)
	v198 = base.Simd_g_v128_load32_lane_l2(m, v172+v59, v7, v195)
	v201 = base.Simd_g_v128_load32_lane_l3(m, v172+v56, v7, v198)
	v203 = base.Simd_g_i8x16_shuffle2(v186, v201, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3))
	v205 = base.Simd_g_i8x16_shuffle2(v186, v201, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6))
	v207 = base.Simd_g_i8x16_shuffle2(v203, v205, base.Simd_g_const(&F_HFilter8i_SSE2__k8), base.Simd_g_const(&F_HFilter8i_SSE2__k9))
	v209 = base.Simd_g_i8x16_shuffle2(v203, v205, base.Simd_g_const(&F_HFilter8i_SSE2__k11), base.Simd_g_const(&F_HFilter8i_SSE2__k12))
	v211 = base.Simd_g_i8x16_shuffle2(v207, v209, base.Simd_g_const(&F_HFilter8i_SSE2__k21), base.Simd_g_const(&F_HFilter8i_SSE2__k22))
	v212 = base.Simd_g_const(&F_HFilter8i_SSE2__k23)
	v213 = base.Simd_g_i8x16_shuffle2(v170, v211, base.Simd_g_const(&F_HFilter8i_SSE2__k24), base.Simd_g_const(&F_HFilter8i_SSE2__k25))
	v214 = base.Simd_g_v128_xor(v213, v127)
	v217 = base.Simd_g_i8x16_shuffle2(v85, v124, base.Simd_g_const(&F_HFilter8i_SSE2__k24), base.Simd_g_const(&F_HFilter8i_SSE2__k25))
	v220 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v126, v217), base.Simd_g_i8x16_sub_sat_u(v217, v126))
	v222 = base.Simd_g_i8x16_shuffle2(v170, v211, base.Simd_g_const(&F_HFilter8i_SSE2__k17), base.Simd_g_const(&F_HFilter8i_SSE2__k18))
	v225 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v213, v222), base.Simd_g_i8x16_sub_sat_u(v222, v213))
	v229 = base.Simd_g_i8x16_eq(v10, base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(v220, v225), base.Simd_g_i8x16_splat(l5)))
	v231 = base.Simd_g_v128_xor(v222, v127)
	v232 = base.Simd_g_v128_xor(v217, v127)
	v233 = base.Simd_g_i8x16_sub_sat_s(v231, v232)
	v239 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v217, v222), base.Simd_g_i8x16_sub_sat_u(v222, v217))
	v252 = base.Simd_g_i8x16_shuffle2(v81, v83, base.Simd_g_const(&F_HFilter8i_SSE2__k21), base.Simd_g_const(&F_HFilter8i_SSE2__k22))
	v254 = base.Simd_g_i8x16_shuffle2(v120, v122, base.Simd_g_const(&F_HFilter8i_SSE2__k21), base.Simd_g_const(&F_HFilter8i_SSE2__k22))
	v256 = base.Simd_g_i8x16_shuffle2(v252, v254, base.Simd_g_const(&F_HFilter8i_SSE2__k17), base.Simd_g_const(&F_HFilter8i_SSE2__k18))
	v258 = base.Simd_g_i8x16_shuffle2(v252, v254, base.Simd_g_const(&F_HFilter8i_SSE2__k24), base.Simd_g_const(&F_HFilter8i_SSE2__k25))
	v269 = base.Simd_g_i8x16_shuffle2(v166, v168, base.Simd_g_const(&F_HFilter8i_SSE2__k14), base.Simd_g_const(&F_HFilter8i_SSE2__k15))
	v271 = base.Simd_g_i8x16_shuffle2(v207, v209, base.Simd_g_const(&F_HFilter8i_SSE2__k14), base.Simd_g_const(&F_HFilter8i_SSE2__k15))
	v273 = base.Simd_g_i8x16_shuffle2(v269, v271, base.Simd_g_const(&F_HFilter8i_SSE2__k24), base.Simd_g_const(&F_HFilter8i_SSE2__k25))
	v275 = base.Simd_g_i8x16_shuffle2(v269, v271, base.Simd_g_const(&F_HFilter8i_SSE2__k17), base.Simd_g_const(&F_HFilter8i_SSE2__k18))
	v288 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_i8x16_add_sat_s(base.Simd_g_v128_bitselect(v10, base.Simd_g_i8x16_sub_sat_s(v128, v214), v229), v233), v233), v233), v10, base.Simd_g_i8x16_eq(v10, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_add_sat_u(base.Simd_g_i8x16_add_sat_u(v239, v239), base.Simd_g_v128_and(base.Simd_g_i16x8_shr_u(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v126, v213), base.Simd_g_i8x16_sub_sat_u(v213, v126)), v38), base.Simd_g_const(&F_HFilter8i_SSE2__k26))), base.Simd_g_i8x16_splat(l3)), base.Simd_g_i8x16_sub_sat_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(base.Simd_g_i8x16_max_u(v220, base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v256, v258), base.Simd_g_i8x16_sub_sat_u(v258, v256))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v258, v126), base.Simd_g_i8x16_sub_sat_u(v126, v258))), v225), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v273, v275), base.Simd_g_i8x16_sub_sat_u(v275, v273))), base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v275, v213), base.Simd_g_i8x16_sub_sat_u(v213, v275))), base.Simd_g_i8x16_splat(l4)))))
	v290 = base.Simd_g_i8x16_add_sat_s(v288, base.Simd_g_const(&F_HFilter8i_SSE2__k27))
	v293 = int32(11)
	v299 = base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v290, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3)), v293), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v290, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6)), v293))
	v304 = base.Simd_g_v128_bitselect(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(base.Simd_g_v128_xor(v299, v127), v10), base.Simd_g_const(&F_HFilter8i_SSE2__k28)), v10, v229)
	v306 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v128, v304), v127)
	v308 = base.Simd_g_i8x16_add_sat_s(v288, base.Simd_g_const(&F_HFilter8i_SSE2__k29))
	v319 = base.Simd_g_v128_xor(base.Simd_g_i8x16_add_sat_s(v232, base.Simd_g_i8x16_narrow_i16x8_s(base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v308, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3)), v293), base.Simd_g_i16x8_shr_s(base.Simd_g_i8x16_shuffle2(v10, v308, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6)), v293))), v127)
	v321 = base.Simd_g_i8x16_shuffle2(v306, v319, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3))
	v323 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v231, v299), v127)
	v325 = base.Simd_g_v128_xor(base.Simd_g_i8x16_sub_sat_s(v214, v304), v127)
	v327 = base.Simd_g_i8x16_shuffle2(v323, v325, base.Simd_g_const(&F_HFilter8i_SSE2__k2), base.Simd_g_const(&F_HFilter8i_SSE2__k3))
	v329 = base.Simd_g_i8x16_shuffle2(v321, v327, base.Simd_g_const(&F_HFilter8i_SSE2__k8), base.Simd_g_const(&F_HFilter8i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, l0, v41, v329)
	v334 = base.Simd_g_i8x16_shuffle2(v306, v319, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6))
	v336 = base.Simd_g_i8x16_shuffle2(v323, v325, base.Simd_g_const(&F_HFilter8i_SSE2__k5), base.Simd_g_const(&F_HFilter8i_SSE2__k6))
	v338 = base.Simd_g_i8x16_shuffle2(v334, v336, base.Simd_g_const(&F_HFilter8i_SSE2__k8), base.Simd_g_const(&F_HFilter8i_SSE2__k9))
	base.Simd_g_v128_store32_lane_l0(m, l1, v41, v338)
	v343 = l0 + v41
	v344 = v343 + l2
	v345 = base.Simd_g_const(&F_HFilter8i_SSE2__k30)
	v346 = base.Simd_g_i8x16_shuffle2(v329, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v344, v7, v346)
	v351 = l1 + v41
	v352 = v351 + l2
	v354 = base.Simd_g_i8x16_shuffle2(v338, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v352, v7, v354)
	v358 = v344 + l2
	v360 = base.Simd_g_i8x16_shuffle2(v346, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v358, v7, v360)
	v364 = v352 + l2
	v366 = base.Simd_g_i8x16_shuffle2(v354, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v364, v7, v366)
	base.Simd_g_v128_store32_lane_l0(m, v358+l2, v7, base.Simd_g_i8x16_shuffle2(v360, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32)))
	base.Simd_g_v128_store32_lane_l0(m, v364+l2, v7, base.Simd_g_i8x16_shuffle2(v366, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32)))
	v382 = v343 + v42
	v384 = base.Simd_g_i8x16_shuffle2(v321, v327, base.Simd_g_const(&F_HFilter8i_SSE2__k11), base.Simd_g_const(&F_HFilter8i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v382, v7, v384)
	v388 = v351 + v42
	v390 = base.Simd_g_i8x16_shuffle2(v334, v336, base.Simd_g_const(&F_HFilter8i_SSE2__k11), base.Simd_g_const(&F_HFilter8i_SSE2__k12))
	base.Simd_g_v128_store32_lane_l0(m, v388, v7, v390)
	v394 = v382 + l2
	v396 = base.Simd_g_i8x16_shuffle2(v384, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v394, v7, v396)
	v400 = v388 + l2
	v402 = base.Simd_g_i8x16_shuffle2(v390, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v400, v7, v402)
	v406 = v394 + l2
	v408 = base.Simd_g_i8x16_shuffle2(v396, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v406, v7, v408)
	v412 = v400 + l2
	v414 = base.Simd_g_i8x16_shuffle2(v402, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32))
	base.Simd_g_v128_store32_lane_l0(m, v412, v7, v414)
	base.Simd_g_v128_store32_lane_l0(m, v406+l2, v7, base.Simd_g_i8x16_shuffle2(v408, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32)))
	base.Simd_g_v128_store32_lane_l0(m, v412+l2, v7, base.Simd_g_i8x16_shuffle2(v414, v10, base.Simd_g_const(&F_HFilter8i_SSE2__k31), base.Simd_g_const(&F_HFilter8i_SSE2__k32)))
	return
}

var F_HFilter8i_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_HFilter8i_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_HFilter8i_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_HFilter8i_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_HFilter8i_SSE2__k4 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_HFilter8i_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_HFilter8i_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_HFilter8i_SSE2__k7 = [2]uint64{0x1312030211100100, 0x1716070615140504}
var F_HFilter8i_SSE2__k8 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_HFilter8i_SSE2__k9 = [2]uint64{0x302808001008080, 0x706808005048080}
var F_HFilter8i_SSE2__k10 = [2]uint64{0x1b1a0b0a19180908, 0x1f1e0f0e1d1c0d0c}
var F_HFilter8i_SSE2__k11 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_HFilter8i_SSE2__k12 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_HFilter8i_SSE2__k13 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_HFilter8i_SSE2__k14 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_HFilter8i_SSE2__k15 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_HFilter8i_SSE2__k16 = [2]uint64{0x706050403020100, 0x1716151413121110}
var F_HFilter8i_SSE2__k17 = [2]uint64{0x706050403020100, 0x8080808080808080}
var F_HFilter8i_SSE2__k18 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_HFilter8i_SSE2__k19 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_HFilter8i_SSE2__k20 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_HFilter8i_SSE2__k21 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_HFilter8i_SSE2__k22 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_HFilter8i_SSE2__k23 = [2]uint64{0xf0e0d0c0b0a0908, 0x1f1e1d1c1b1a1918}
var F_HFilter8i_SSE2__k24 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_HFilter8i_SSE2__k25 = [2]uint64{0x8080808080808080, 0xf0e0d0c0b0a0908}
var F_HFilter8i_SSE2__k26 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_HFilter8i_SSE2__k27 = [2]uint64{0x404040404040404, 0x404040404040404}
var F_HFilter8i_SSE2__k28 = [2]uint64{0xc0c0c0c0c0c0c0c0, 0xc0c0c0c0c0c0c0c0}
var F_HFilter8i_SSE2__k29 = [2]uint64{0x303030303030303, 0x303030303030303}
var F_HFilter8i_SSE2__k30 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_HFilter8i_SSE2__k31 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_HFilter8i_SSE2__k32 = [2]uint64{0x8080808080808080, 0x302010080808080}

func F_HasAlpha32b_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v34 base.V128
	_ = v34
	var v35 base.V128
	_ = v35
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v48 int32
	_ = v48
	var v54 int32
	_ = v54
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 base.V128
	_ = v62
	var v63 base.V128
	_ = v63
	var v66 base.V128
	_ = v66
	var v70 base.V128
	_ = v70
	var v75 base.V128
	_ = v75
	var v95 int32
	_ = v95
	v10 = l1<<(uint(int32(2))%32) + int32(-3)
	v13 = int32(0)
	goto L3
L1:
	;
	return v95
L2:
	;
	v95 = int32(1)
	goto L1
L3:
	;
	v19 = v13 + int32(64)
	if v19 <= v10 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	goto L2
L5:
	;
	v60 = l0 + v13
	v61 = int32(0)
	v62 = base.Simd_g_v128_load_rng(m, v60, v61, int32(0), int32(64))
	v63 = base.Simd_g_const(&F_HasAlpha32b_SSE2__k0)
	v66 = base.Simd_g_v128_load_nc(m, v60, int32(16))
	v70 = base.Simd_g_v128_load_nc(m, v60, int32(32))
	v75 = base.Simd_g_v128_load_nc(m, v60+int32(48), v61)
	if base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(v62, v63), base.Simd_g_v128_and(v66, v63)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(v70, v63), base.Simd_g_v128_and(v75, v63))), base.Simd_g_const(&F_HasAlpha32b_SSE2__k1))) == int32(_a_F_HasAlpha32b_SSE2_0) {
		v13 = v19
		goto L3
	} else {
		goto L18
	}
L6:
	;
	v22 = v13
	goto L8
L7:
	;
	v48 = v22
	goto L14
L8:
	;
	v28 = v22 + int32(32)
	if v28 <= v10 {
		goto L10
	} else {
		goto L11
	}
L10:
	;
	v32 = l0 + v22
	v34 = base.Simd_g_v128_load_rng(m, v32, int32(0), int32(0), int32(32))
	v35 = base.Simd_g_const(&F_HasAlpha32b_SSE2__k0)
	v38 = base.Simd_g_v128_load_nc(m, v32, int32(16))
	v40 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_v128_and(v34, v35), base.Simd_g_v128_and(v38, v35))
	if base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(base.Simd_g_i8x16_narrow_i16x8_u(v40, v40), base.Simd_g_const(&F_HasAlpha32b_SSE2__k1))) == int32(_a_F_HasAlpha32b_SSE2_0) {
		v22 = v28
		goto L8
	} else {
		goto L13
	}
L11:
	;
	v30 = int32(0)
	if v10 < v22 {
		v95 = v30
		goto L1
	} else {
		goto L12
	}
L12:
	;
	goto L7
L13:
	;
	goto L2
L14:
	;
	v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v48))))
	if v54 != int32(255) {
		goto L2
	} else {
		goto L16
	}
L16:
	;
	v58 = v48 + int32(4)
	if v10 < v58 {
		v95 = v30
		goto L1
	} else {
		goto L17
	}
L17:
	;
	v48 = v58
	goto L14
L18:
	;
	goto L4
}

var F_HasAlpha32b_SSE2__k0 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_HasAlpha32b_SSE2__k1 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}

func F_HasAlpha8b_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v21 base.V128
	_ = v21
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v38 int32
	_ = v38
	var v51 int32
	_ = v51
	v9 = int32(0)
	goto L4
L1:
	;
	return v51
L2:
	;
	v51 = int32(1)
	goto L1
L3:
	;
	v29 = l0 + v9
	v30 = l1 - v9
	goto L10
L4:
	;
	v13 = v9 + int32(16)
	if v13 <= l1 {
		goto L6
	} else {
		goto L7
	}
L6:
	;
	v21 = base.Simd_g_v128_load(m, l0+v9, int32(0))
	if base.Simd_g_i8x16_bitmask(base.Simd_g_i8x16_eq(v21, base.Simd_g_const(&F_HasAlpha8b_SSE2__k0))) == int32(_a_F_HasAlpha8b_SSE2_0) {
		v9 = v13
		goto L4
	} else {
		goto L9
	}
L7:
	;
	v15 = int32(0)
	if l1 <= v9 {
		v51 = v15
		goto L1
	} else {
		goto L8
	}
L8:
	;
	goto L3
L9:
	;
	goto L2
L10:
	;
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v29))))
	if v32 != int32(255) {
		goto L2
	} else {
		goto L12
	}
L12:
	;
	v38 = v30 + int32(-1)
	if v38 == int32(0) {
		v51 = v15
		goto L1
	} else {
		goto L13
	}
L13:
	;
	v29 = v29 + int32(1)
	v30 = v38
	goto L10
}

var F_HasAlpha8b_SSE2__k0 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}

func F_HorizontalFilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v23 int32
	_ = v23
	var v30 int32
	_ = v30
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v50 int32
	_ = v50
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v57 int32
	_ = v57
	var v62 int32
	_ = v62
	var v71 int32
	_ = v71
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v114 int32
	_ = v114
	var v137 int32
	_ = v137
	var v141 int32
	_ = v141
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v150 int32
	_ = v150
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v187 int32
	_ = v187
	var v188 base.V128
	_ = v188
	var v190 base.V128
	_ = v190
	var v191 base.V128
	_ = v191
	var v197 int32
	_ = v197
	var v207 int32
	_ = v207
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v236 int32
	_ = v236
	var v238 int32
	_ = v238
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v248 int32
	_ = v248
	var v267 int32
	_ = v267
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v284 int32
	_ = v284
	var v287 int32
	_ = v287
	var v290 int32
	_ = v290
	var v301 int32
	_ = v301
	var v305 int32
	_ = v305
	var v306 int32
	_ = v306
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v321 int32
	_ = v321
	var v325 int32
	_ = v325
	var v326 int32
	_ = v326
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v333 int32
	_ = v333
	var v335 int32
	_ = v335
	var v338 int32
	_ = v338
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v362 int32
	_ = v362
	var v365 int32
	_ = v365
	var v367 int32
	_ = v367
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v381 int32
	_ = v381
	v18 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v18)
	if l1 < int32(2) {
		if l2 < int32(2) {
		} else {
			v272 = l2 + int32(-1)
			v273 = int32(3)
			v274 = v272 & v273
			if base.Ui32(l2+int32(-2)) < base.Ui32(v273) {
				v338 = l0
				v342 = l4
				v343 = v18
			} else {
				v284 = l3 << (uint(int32(1)) % 32)
				v287 = l3 * int32(3)
				v290 = l3 << (uint(int32(2)) % 32)
				v301 = v18
				v305 = v272 & int32(-4)
				v306 = int32(0)
				for {
					v315 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l3+v306))))
					v316 = v315 - v301
					*(*uint8)(unsafe.Add(mBase, uint32(l4+l3+v306))) = uint8(v316)
					v320 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v284+v306))))
					v321 = v320 - v315
					*(*uint8)(unsafe.Add(mBase, uint32(l4+v284+v306))) = uint8(v321)
					v325 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v287+v306))))
					v326 = v325 - v320
					*(*uint8)(unsafe.Add(mBase, uint32(l4+v287+v306))) = uint8(v326)
					v330 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v290+v306))))
					v331 = v330 - v325
					*(*uint8)(unsafe.Add(mBase, uint32(l4+v290+v306))) = uint8(v331)
					v333 = v306 + v290
					v335 = v305 + int32(-4)
					if v335 != 0 {
						v301 = v330
						v305 = v335
						v306 = v333
						continue
					} else {
						break
					}
					break
				}
				v338 = l0 + v333
				v342 = l4 + v333
				v343 = v330
			}
			if v274 == int32(0) {
			} else {
				v362 = v343
				v365 = v274
				v367 = l3
				for {
					v376 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v338+v367))))
					v377 = v376 - v362
					*(*uint8)(unsafe.Add(mBase, uint32(v342+v367))) = uint8(v377)
					v381 = v365 + int32(-1)
					if v381 != 0 {
						v362 = v376
						v365 = v381
						v367 = v367 + l3
						continue
					} else {
						break
					}
					break
				}
			}
		}
	} else {
		v23 = l1 + int32(-1)
		if base.Ui32(l1) < base.Ui32(int32(17)) {
			v71 = int32(0)
			v96 = v71 + l4 + int32(1)
			v97 = v71 ^ int32(-1) + l1
			v98 = l0 + v71
			for {
				v105 = int32(1)
				v106 = v98 + v105
				v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v106))))
				v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v98))))
				v109 = v107 - v108
				*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v109)
				v114 = v97 + int32(-1)
				if v114 != 0 {
					v96 = v96 + v105
					v97 = v114
					v98 = v106
					continue
				} else {
					break
				}
				break
			}
		} else {
			v30 = v23 & int32(-16)
			v39 = l4 + int32(1)
			v40 = v30
			v41 = l0
			for {
				v50 = int32(0)
				v51 = base.Simd_g_v128_load_rng(m, v41+int32(1), v50, int32(-1), int32(17))
				v53 = base.Simd_g_v128_load_nc(m, v41, v50)
				v54 = base.Simd_g_i8x16_sub(v51, v53)
				base.Simd_g_v128_store(m, v39, v50, v54)
				v57 = int32(16)
				v62 = v40 + int32(-16)
				if v62 != 0 {
					v39 = v39 + v57
					v40 = v62
					v41 = v41 + v57
					continue
				} else {
					break
				}
				break
			}
			if v23 == v30 {
			} else {
				v71 = v30
				v96 = v71 + l4 + int32(1)
				v97 = v71 ^ int32(-1) + l1
				v98 = l0 + v71
				for {
					v105 = int32(1)
					v106 = v98 + v105
					v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v106))))
					v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v98))))
					v109 = v107 - v108
					*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v109)
					v114 = v97 + int32(-1)
					if v114 != 0 {
						v96 = v96 + v105
						v97 = v114
						v98 = v106
						continue
					} else {
						break
					}
					break
				}
			}
		}
		if l2 < int32(2) {
		} else {
			v137 = v23 & int32(-16)
			v141 = l0
			v145 = l4
			v146 = v18
			v150 = int32(1)
			for {
				v158 = v145 + l3
				v159 = v141 + l3
				v160 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v159))))
				v161 = v160 - v146
				*(*uint8)(unsafe.Add(mBase, uint32(v158))) = uint8(v161)
				if base.Ui32(l1) < base.Ui32(int32(17)) {
					v207 = int32(0)
					v226 = l1 + int32(-1) - v207
					v228 = l3 + v207
					for {
						v236 = int32(1)
						v238 = v141 + v228
						v241 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v238+v236))))
						v242 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v238))))
						v243 = v241 - v242
						*(*uint8)(unsafe.Add(mBase, uint32(v145+v228+v236))) = uint8(v243)
						v248 = v226 + int32(-1)
						if v248 != 0 {
							v226 = v248
							v228 = v228 + v236
							continue
						} else {
							break
						}
						break
					}
				} else {
					v172 = v137
					v174 = l3
					for {
						v182 = int32(1)
						v184 = v141 + v174
						v187 = int32(0)
						v188 = base.Simd_g_v128_load_rng(m, v184+v182, v187, int32(-1), int32(17))
						v190 = base.Simd_g_v128_load_nc(m, v184, v187)
						v191 = base.Simd_g_i8x16_sub(v188, v190)
						base.Simd_g_v128_store(m, v145+v174+v182, v187, v191)
						v197 = v172 + int32(-16)
						if v197 != 0 {
							v172 = v197
							v174 = v174 + int32(16)
							continue
						} else {
							break
						}
						break
					}
					if v23 == v137 {
					} else {
						v207 = v137
						v226 = l1 + int32(-1) - v207
						v228 = l3 + v207
						for {
							v236 = int32(1)
							v238 = v141 + v228
							v241 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v238+v236))))
							v242 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v238))))
							v243 = v241 - v242
							*(*uint8)(unsafe.Add(mBase, uint32(v145+v228+v236))) = uint8(v243)
							v248 = v226 + int32(-1)
							if v248 != 0 {
								v226 = v248
								v228 = v228 + v236
								continue
							} else {
								break
							}
							break
						}
					}
				}
				v267 = v150 + int32(1)
				if v267 != l2 {
					v141 = v159
					v145 = v158
					v146 = v160
					v150 = v267
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
func F_HorizontalFilter_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v30 int32
	_ = v30
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v67 base.V128
	_ = v67
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v75 int32
	_ = v75
	var v80 base.V128
	_ = v80
	var v82 base.V128
	_ = v82
	var v83 base.V128
	_ = v83
	var v86 int32
	_ = v86
	var v91 base.V128
	_ = v91
	var v95 base.V128
	_ = v95
	var v96 base.V128
	_ = v96
	var v99 int32
	_ = v99
	var v104 base.V128
	_ = v104
	var v108 base.V128
	_ = v108
	var v109 base.V128
	_ = v109
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v124 int32
	_ = v124
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 base.V128
	_ = v144
	var v146 base.V128
	_ = v146
	var v147 base.V128
	_ = v147
	var v150 int32
	_ = v150
	var v151 base.V128
	_ = v151
	var v155 base.V128
	_ = v155
	var v156 base.V128
	_ = v156
	var v169 int32
	_ = v169
	var v183 int32
	_ = v183
	var v191 int32
	_ = v191
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v214 int32
	_ = v214
	var v215 base.V128
	_ = v215
	var v217 base.V128
	_ = v217
	var v218 base.V128
	_ = v218
	var v221 int32
	_ = v221
	var v226 int32
	_ = v226
	var v236 int32
	_ = v236
	var v262 int32
	_ = v262
	var v265 int32
	_ = v265
	var v266 int32
	_ = v266
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v277 int32
	_ = v277
	var v282 int32
	_ = v282
	var v307 int32
	_ = v307
	var v309 int32
	_ = v309
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v322 int32
	_ = v322
	var v329 int32
	_ = v329
	var v330 int32
	_ = v330
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v338 int32
	_ = v338
	var v339 int32
	_ = v339
	var v348 int32
	_ = v348
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	var v351 int32
	_ = v351
	var v362 int32
	_ = v362
	var v363 int32
	_ = v363
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v381 base.V128
	_ = v381
	var v385 base.V128
	_ = v385
	var v386 base.V128
	_ = v386
	var v389 int32
	_ = v389
	var v394 base.V128
	_ = v394
	var v396 base.V128
	_ = v396
	var v397 base.V128
	_ = v397
	var v400 int32
	_ = v400
	var v405 base.V128
	_ = v405
	var v409 base.V128
	_ = v409
	var v410 base.V128
	_ = v410
	var v413 int32
	_ = v413
	var v418 base.V128
	_ = v418
	var v422 base.V128
	_ = v422
	var v423 base.V128
	_ = v423
	var v427 int32
	_ = v427
	var v429 int32
	_ = v429
	var v438 int32
	_ = v438
	var v449 int32
	_ = v449
	var v451 int32
	_ = v451
	var v454 int32
	_ = v454
	var v455 int32
	_ = v455
	var v456 base.V128
	_ = v456
	var v458 base.V128
	_ = v458
	var v459 base.V128
	_ = v459
	var v462 int32
	_ = v462
	var v463 base.V128
	_ = v463
	var v467 base.V128
	_ = v467
	var v468 base.V128
	_ = v468
	var v473 int32
	_ = v473
	var v488 int32
	_ = v488
	var v491 int32
	_ = v491
	var v500 int32
	_ = v500
	var v502 int32
	_ = v502
	var v505 int32
	_ = v505
	var v506 base.V128
	_ = v506
	var v508 base.V128
	_ = v508
	var v509 base.V128
	_ = v509
	var v515 int32
	_ = v515
	var v530 int32
	_ = v530
	var v549 int32
	_ = v549
	var v557 int32
	_ = v557
	var v559 int32
	_ = v559
	var v562 int32
	_ = v562
	var v563 int32
	_ = v563
	var v564 int32
	_ = v564
	var v567 int32
	_ = v567
	var v591 int32
	_ = v591
	var v596 int32
	_ = v596
	var v597 int32
	_ = v597
	var v598 int32
	_ = v598
	var v608 int32
	_ = v608
	var v611 int32
	_ = v611
	var v614 int32
	_ = v614
	var v625 int32
	_ = v625
	var v631 int32
	_ = v631
	var v632 int32
	_ = v632
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v656 int32
	_ = v656
	var v657 int32
	_ = v657
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v664 int32
	_ = v664
	var v668 int32
	_ = v668
	var v669 int32
	_ = v669
	var v690 int32
	_ = v690
	var v693 int32
	_ = v693
	var v697 int32
	_ = v697
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v711 int32
	_ = v711
	var v715 int32
	_ = v715
	var v719 int32
	_ = v719
	var v723 int32
	_ = v723
	var v724 int32
	_ = v724
	var v729 int32
	_ = v729
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v740 int32
	_ = v740
	var v741 int32
	_ = v741
	var v752 int32
	_ = v752
	var v756 int32
	_ = v756
	var v764 int32
	_ = v764
	var v766 int32
	_ = v766
	var v769 int32
	_ = v769
	var v770 base.V128
	_ = v770
	var v772 base.V128
	_ = v772
	var v773 base.V128
	_ = v773
	var v779 int32
	_ = v779
	var v789 int32
	_ = v789
	var v810 int32
	_ = v810
	var v814 int32
	_ = v814
	var v822 int32
	_ = v822
	var v824 int32
	_ = v824
	var v827 int32
	_ = v827
	var v828 int32
	_ = v828
	var v829 int32
	_ = v829
	var v834 int32
	_ = v834
	var v855 int32
	_ = v855
	v20 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v20)
	v23 = l1 + int32(-1)
	v25 = v23 & int32(-32)
	if int32(1) <= v25 {
		v30 = v25 + int32(-1)
		if v30 != int32(31) {
			v49 = int32(0)
			v51 = (int32(base.Ui32(v30)>>(uint(int32(5))%32)) + int32(1)) & int32(268435454)
			for {
				v60 = l4 + v49
				v61 = int32(17)
				v63 = l0 + v49
				v66 = int32(0)
				v67 = base.Simd_g_v128_load_rng(m, v63+v61, v66, int32(-1), int32(17))
				v71 = base.Simd_g_v128_load_nc(m, v63+int32(16), v66)
				v72 = base.Simd_g_i8x16_sub(v67, v71)
				base.Simd_g_v128_store(m, v60+v61, v66, v72)
				v75 = int32(1)
				v80 = base.Simd_g_v128_load_rng(m, v63+v75, v66, int32(-1), int32(17))
				v82 = base.Simd_g_v128_load_nc(m, v63, v66)
				v83 = base.Simd_g_i8x16_sub(v80, v82)
				base.Simd_g_v128_store(m, v60+v75, v66, v83)
				v86 = int32(49)
				v91 = base.Simd_g_v128_load_rng(m, v63+v86, v66, int32(-1), int32(17))
				v95 = base.Simd_g_v128_load_nc(m, v63+int32(48), v66)
				v96 = base.Simd_g_i8x16_sub(v91, v95)
				base.Simd_g_v128_store(m, v60+v86, v66, v96)
				v99 = int32(33)
				v104 = base.Simd_g_v128_load_rng(m, v63+v99, v66, int32(-1), int32(17))
				v108 = base.Simd_g_v128_load_nc(m, v63+int32(32), v66)
				v109 = base.Simd_g_i8x16_sub(v104, v108)
				base.Simd_g_v128_store(m, v60+v99, v66, v109)
				v113 = v49 + int32(64)
				v115 = v51 + int32(-2)
				if v115 != 0 {
					v49 = v113
					v51 = v115
					continue
				} else {
					break
				}
				break
			}
			v124 = v113
		} else {
			v124 = int32(0)
		}
		if v30&int32(32) != 0 {
			v169 = v124
		} else {
			v137 = int32(1)
			v139 = l4 + v137 + v124
			v142 = l0 + v137 + v124
			v143 = int32(16)
			v144 = base.Simd_g_v128_load_rng(m, v142, v143, int32(15), int32(17))
			v146 = base.Simd_g_v128_load_nc(m, v142, int32(15))
			v147 = base.Simd_g_i8x16_sub(v144, v146)
			base.Simd_g_v128_store(m, v139, v143, v147)
			v150 = int32(0)
			v151 = base.Simd_g_v128_load(m, v142, v150)
			v155 = base.Simd_g_v128_load(m, v142+int32(-1), v150)
			v156 = base.Simd_g_i8x16_sub(v151, v155)
			base.Simd_g_v128_store(m, v139, v150, v156)
			v169 = v124 + int32(32)
		}
	} else {
		v169 = int32(0)
	}
	if v23 <= v169 {
	} else {
		v183 = v169 ^ int32(-1) + l1
		if base.Ui32(v183) <= base.Ui32(int32(15)) {
			v236 = v169
			v262 = v236 ^ int32(-1) + l1
			v265 = v236 + l4 + int32(1)
			v266 = l0 + v236
			for {
				v273 = int32(1)
				v274 = v266 + v273
				v275 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v274))))
				v276 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v266))))
				v277 = v275 - v276
				*(*uint8)(unsafe.Add(mBase, uint32(v265))) = uint8(v277)
				v282 = v262 + int32(-1)
				if v282 != 0 {
					v262 = v282
					v265 = v265 + v273
					v266 = v274
					continue
				} else {
					break
				}
				break
			}
		} else {
			v191 = v183 & int32(-16)
			v203 = v191
			v204 = v169 + l4 + int32(1)
			v205 = l0 + v169
			for {
				v214 = int32(0)
				v215 = base.Simd_g_v128_load_rng(m, v205+int32(1), v214, int32(-1), int32(17))
				v217 = base.Simd_g_v128_load_nc(m, v205, v214)
				v218 = base.Simd_g_i8x16_sub(v215, v217)
				base.Simd_g_v128_store(m, v204, v214, v218)
				v221 = int32(16)
				v226 = v203 + int32(-16)
				if v226 != 0 {
					v203 = v226
					v204 = v204 + v221
					v205 = v205 + v221
					continue
				} else {
					break
				}
				break
			}
			if v183 == v191 {
			} else {
				v236 = v169 + v191
				v262 = v236 ^ int32(-1) + l1
				v265 = v236 + l4 + int32(1)
				v266 = l0 + v236
				for {
					v273 = int32(1)
					v274 = v266 + v273
					v275 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v274))))
					v276 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v266))))
					v277 = v275 - v276
					*(*uint8)(unsafe.Add(mBase, uint32(v265))) = uint8(v277)
					v282 = v262 + int32(-1)
					if v282 != 0 {
						v262 = v282
						v265 = v265 + v273
						v266 = v274
						continue
					} else {
						break
					}
					break
				}
			}
		}
	}
	if l2 < int32(2) {
	} else {
		if v25 < int32(1) {
			if int32(1) < l1 {
				v715 = v23 & int32(-16)
				v719 = l0
				v723 = l4
				v724 = v20
				v729 = int32(1)
				for {
					v738 = v723 + l3
					v739 = v719 + l3
					v740 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v739))))
					v741 = v740 - v724
					*(*uint8)(unsafe.Add(mBase, uint32(v738))) = uint8(v741)
					if base.Ui32(l1) < base.Ui32(int32(17)) {
						v789 = int32(0)
						v810 = l1 + int32(-1) - v789
						v814 = l3 + v789
						for {
							v822 = int32(1)
							v824 = v719 + v814
							v827 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v824+v822))))
							v828 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v824))))
							v829 = v827 - v828
							*(*uint8)(unsafe.Add(mBase, uint32(v723+v814+v822))) = uint8(v829)
							v834 = v810 + int32(-1)
							if v834 != 0 {
								v810 = v834
								v814 = v814 + v822
								continue
							} else {
								break
							}
							break
						}
					} else {
						v752 = v715
						v756 = l3
						for {
							v764 = int32(1)
							v766 = v719 + v756
							v769 = int32(0)
							v770 = base.Simd_g_v128_load_rng(m, v766+v764, v769, int32(-1), int32(17))
							v772 = base.Simd_g_v128_load_nc(m, v766, v769)
							v773 = base.Simd_g_i8x16_sub(v770, v772)
							base.Simd_g_v128_store(m, v723+v756+v764, v769, v773)
							v779 = v752 + int32(-16)
							if v779 != 0 {
								v752 = v779
								v756 = v756 + int32(16)
								continue
							} else {
								break
							}
							break
						}
						if v23 == v715 {
						} else {
							v789 = v715
							v810 = l1 + int32(-1) - v789
							v814 = l3 + v789
							for {
								v822 = int32(1)
								v824 = v719 + v814
								v827 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v824+v822))))
								v828 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v824))))
								v829 = v827 - v828
								*(*uint8)(unsafe.Add(mBase, uint32(v723+v814+v822))) = uint8(v829)
								v834 = v810 + int32(-1)
								if v834 != 0 {
									v810 = v834
									v814 = v814 + v822
									continue
								} else {
									break
								}
								break
							}
						}
					}
					v855 = v729 + int32(1)
					if v855 != l2 {
						v719 = v739
						v723 = v738
						v724 = v740
						v729 = v855
						continue
					} else {
						break
					}
					break
				}
			} else {
				v596 = l2 + int32(-1)
				v597 = int32(3)
				v598 = v596 & v597
				if base.Ui32(l2+int32(-2)) < base.Ui32(v597) {
					v664 = l0
					v668 = l4
					v669 = v20
				} else {
					v608 = l3 << (uint(int32(1)) % 32)
					v611 = l3 * int32(3)
					v614 = l3 << (uint(int32(2)) % 32)
					v625 = v20
					v631 = v596 & int32(-4)
					v632 = int32(0)
					for {
						v641 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l3+v632))))
						v642 = v641 - v625
						*(*uint8)(unsafe.Add(mBase, uint32(l4+l3+v632))) = uint8(v642)
						v646 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v608+v632))))
						v647 = v646 - v641
						*(*uint8)(unsafe.Add(mBase, uint32(l4+v608+v632))) = uint8(v647)
						v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v611+v632))))
						v652 = v651 - v646
						*(*uint8)(unsafe.Add(mBase, uint32(l4+v611+v632))) = uint8(v652)
						v656 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v614+v632))))
						v657 = v656 - v651
						*(*uint8)(unsafe.Add(mBase, uint32(l4+v614+v632))) = uint8(v657)
						v659 = v632 + v614
						v661 = v631 + int32(-4)
						if v661 != 0 {
							v625 = v656
							v631 = v661
							v632 = v659
							continue
						} else {
							break
						}
						break
					}
					v664 = l0 + v659
					v668 = l4 + v659
					v669 = v656
				}
				if v598 == int32(0) {
				} else {
					v690 = v669
					v693 = v598
					v697 = l3
					for {
						v706 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v664+v697))))
						v707 = v706 - v690
						*(*uint8)(unsafe.Add(mBase, uint32(v668+v697))) = uint8(v707)
						v711 = v693 + int32(-1)
						if v711 != 0 {
							v690 = v706
							v693 = v711
							v697 = v697 + l3
							continue
						} else {
							break
						}
						break
					}
				}
			}
		} else {
			v307 = l1 + int32(-33)
			v309 = v25 + int32(-1)
			v312 = v307 - v309&int32(-32)
			v314 = v307 & int32(15)
			v315 = v312 - v314
			v322 = int32(1)
			v329 = l0
			v330 = l0 + l3
			v333 = l4
			v334 = v20
			v338 = v322
			v339 = l4 + l3
			for {
				v348 = v333 + l3
				v349 = v329 + l3
				v350 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v349))))
				v351 = v350 - v334
				*(*uint8)(unsafe.Add(mBase, uint32(v348))) = uint8(v351)
				if v309 == int32(31) {
					v438 = int32(0)
				} else {
					v362 = (int32(base.Ui32(v309)>>(uint(int32(5))%32)) + v322) & int32(268435454)
					v363 = int32(0)
					for {
						v374 = v339 + v363
						v375 = int32(17)
						v377 = v330 + v363
						v380 = int32(0)
						v381 = base.Simd_g_v128_load_rng(m, v377+v375, v380, int32(-1), int32(17))
						v385 = base.Simd_g_v128_load_nc(m, v377+int32(16), v380)
						v386 = base.Simd_g_i8x16_sub(v381, v385)
						base.Simd_g_v128_store(m, v374+v375, v380, v386)
						v389 = int32(1)
						v394 = base.Simd_g_v128_load_rng(m, v377+v389, v380, int32(-1), int32(17))
						v396 = base.Simd_g_v128_load_nc(m, v377, v380)
						v397 = base.Simd_g_i8x16_sub(v394, v396)
						base.Simd_g_v128_store(m, v374+v389, v380, v397)
						v400 = int32(49)
						v405 = base.Simd_g_v128_load_rng(m, v377+v400, v380, int32(-1), int32(17))
						v409 = base.Simd_g_v128_load_nc(m, v377+int32(48), v380)
						v410 = base.Simd_g_i8x16_sub(v405, v409)
						base.Simd_g_v128_store(m, v374+v400, v380, v410)
						v413 = int32(33)
						v418 = base.Simd_g_v128_load_rng(m, v377+v413, v380, int32(-1), int32(17))
						v422 = base.Simd_g_v128_load_nc(m, v377+int32(32), v380)
						v423 = base.Simd_g_i8x16_sub(v418, v422)
						base.Simd_g_v128_store(m, v374+v413, v380, v423)
						v427 = v363 + int32(64)
						v429 = v362 + int32(-2)
						if v429 != 0 {
							v362 = v429
							v363 = v427
							continue
						} else {
							break
						}
						break
					}
					v438 = v427
				}
				if v309&int32(32) != 0 {
					v473 = v438
				} else {
					v449 = int32(1)
					v451 = v348 + v449 + v438
					v454 = v349 + v449 + v438
					v455 = int32(16)
					v456 = base.Simd_g_v128_load_rng(m, v454, v455, int32(15), int32(17))
					v458 = base.Simd_g_v128_load_nc(m, v454, int32(15))
					v459 = base.Simd_g_i8x16_sub(v456, v458)
					base.Simd_g_v128_store(m, v451, v455, v459)
					v462 = int32(0)
					v463 = base.Simd_g_v128_load(m, v454, v462)
					v467 = base.Simd_g_v128_load(m, v454+int32(-1), v462)
					v468 = base.Simd_g_i8x16_sub(v463, v467)
					base.Simd_g_v128_store(m, v451, v462, v468)
					v473 = v438 + int32(32)
				}
				if v23 <= v473 {
				} else {
					if base.Ui32(int32(15)) < base.Ui32(v312) {
						v488 = v473
						v491 = v315
						for {
							v500 = int32(1)
							v502 = v330 + v488
							v505 = int32(0)
							v506 = base.Simd_g_v128_load_rng(m, v502+v500, v505, int32(-1), int32(17))
							v508 = base.Simd_g_v128_load_nc(m, v502, v505)
							v509 = base.Simd_g_i8x16_sub(v506, v508)
							base.Simd_g_v128_store(m, v339+v488+v500, v505, v509)
							v515 = v491 + int32(-16)
							if v515 != 0 {
								v488 = v488 + int32(16)
								v491 = v515
								continue
							} else {
								break
							}
							break
						}
						if v314 == int32(0) {
						} else {
							v530 = v473 + v315
							v549 = v530
							for {
								v557 = int32(1)
								v559 = v330 + v549
								v562 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v559+v557))))
								v563 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v559))))
								v564 = v562 - v563
								*(*uint8)(unsafe.Add(mBase, uint32(v339+v549+v557))) = uint8(v564)
								v567 = v549 + v557
								if v23 != v567 {
									v549 = v567
									continue
								} else {
									break
								}
								break
							}
						}
					} else {
						v530 = v473
						v549 = v530
						for {
							v557 = int32(1)
							v559 = v330 + v549
							v562 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v559+v557))))
							v563 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v559))))
							v564 = v562 - v563
							*(*uint8)(unsafe.Add(mBase, uint32(v339+v549+v557))) = uint8(v564)
							v567 = v549 + v557
							if v23 != v567 {
								v549 = v567
								continue
							} else {
								break
							}
							break
						}
					}
				}
				v591 = v338 + int32(1)
				if v591 != l2 {
					v329 = v349
					v330 = v330 + l3
					v333 = v348
					v334 = v350
					v338 = v591
					v339 = v339 + l3
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
func F_HorizontalUnfilter_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v22 int32
	_ = v22
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v39 base.V128
	_ = v39
	var v43 base.V128
	_ = v43
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v50 base.V128
	_ = v50
	var v53 base.V128
	_ = v53
	var v56 base.V128
	_ = v56
	var v74 int32
	_ = v74
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v110 int32
	_ = v110
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v139 int32
	_ = v139
	var v141 int32
	_ = v141
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v155 int32
	_ = v155
	var v160 int32
	_ = v160
	v11 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
	if l0 != 0 {
		v13 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
		v14 = v13
	} else {
		v14 = int32(0)
	}
	v15 = v14 + v11
	*(*uint8)(unsafe.Add(mBase, uint32(l2))) = uint8(v15)
	if l3 < int32(2) {
	} else {
		if base.Ui32(l3) < base.Ui32(int32(9)) {
			v74 = int32(1)
		} else {
			v22 = int32(1)
			v29 = int32(0)
			v32 = v29
			v39 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k0), v15&int32(255))
			for {
				v43 = base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k0)
				v45 = int32(0)
				v46 = base.Simd_g_v128_load64_zero(m, l1+v22+v32, v45)
				v47 = base.Simd_g_i8x16_add(v46, v39)
				v50 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v43, v47, base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k1), base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k2)), v47)
				v53 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v43, v50, base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k3), base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k4)), v50)
				v56 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v43, v53, base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k5), base.Simd_g_const(&F_HorizontalUnfilter_SSE2__k6)), v53)
				base.Simd_g_v128_store64_lane_l0(m, l2+v22+v32, v45, v56)
				if v32+int32(17) <= l3 {
					v32 = v32 + int32(8)
					v39 = base.Simd_g_i64x2_shr_u(v56, int32(56))
					continue
				} else {
					break
				}
				break
			}
			v74 = v32 + int32(9)
		}
		if l3 <= v74 {
		} else {
			v83 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v74+int32(-1)))))
			v86 = (l3 - v74) & int32(3)
			if v86 != 0 {
				v87 = v74
				v91 = v83
				v93 = v86
				for {
					v99 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v87))))
					v100 = v91 + v99
					*(*uint8)(unsafe.Add(mBase, uint32(l2+v87))) = uint8(v100)
					v103 = v87 + int32(1)
					v105 = v93 + int32(-1)
					if v105 != 0 {
						v87 = v103
						v91 = v100
						v93 = v105
						continue
					} else {
						break
					}
					break
				}
				v106 = v103
				v110 = v100
			} else {
				v106 = v74
				v110 = v83
			}
			if base.Ui32(int32(-4)) < base.Ui32(v74-l3) {
			} else {
				v120 = l1
				v121 = l2
				v122 = l3
				v123 = v110
				for {
					v129 = v121 + v106
					v130 = v120 + v106
					v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130))))
					v132 = v123 + v131
					*(*uint8)(unsafe.Add(mBase, uint32(v129))) = uint8(v132)
					v134 = int32(1)
					v138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130+v134))))
					v139 = v132 + v138
					*(*uint8)(unsafe.Add(mBase, uint32(v129+v134))) = uint8(v139)
					v141 = int32(2)
					v145 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130+v141))))
					v146 = v139 + v145
					*(*uint8)(unsafe.Add(mBase, uint32(v129+v141))) = uint8(v146)
					v148 = int32(3)
					v152 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130+v148))))
					v153 = v146 + v152
					*(*uint8)(unsafe.Add(mBase, uint32(v129+v148))) = uint8(v153)
					v155 = int32(4)
					v160 = v122 + int32(-4)
					if v106 != v160 {
						v120 = v120 + v155
						v121 = v121 + v155
						v122 = v160
						v123 = v153
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

var F_HorizontalUnfilter_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_HorizontalUnfilter_SSE2__k1 = [2]uint64{0x808080808080800f, 0x8080808080808080}
var F_HorizontalUnfilter_SSE2__k2 = [2]uint64{0x605040302010080, 0xe0d0c0b0a090807}
var F_HorizontalUnfilter_SSE2__k3 = [2]uint64{0x8080808080800f0e, 0x8080808080808080}
var F_HorizontalUnfilter_SSE2__k4 = [2]uint64{0x504030201008080, 0xd0c0b0a09080706}
var F_HorizontalUnfilter_SSE2__k5 = [2]uint64{0x808080800f0e0d0c, 0x8080808080808080}
var F_HorizontalUnfilter_SSE2__k6 = [2]uint64{0x302010080808080, 0xb0a090807060504}

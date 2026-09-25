//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_PaletteCompareColorsForQsort(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	if base.Ui32(v5) < base.Ui32(v6) {
		v8 = int32(-1)
	} else {
		v8 = int32(1)
	}
	return v8
}
func F_PopulationCost(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) int64 {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v36 int32
	_ = v36
	var v48 int64
	_ = v48
	var v51 int64
	_ = v51
	var v52 int64
	_ = v52
	var v55 int64
	_ = v55
	var v58 int64
	_ = v58
	var v60 int64
	_ = v60
	var v61 int64
	_ = v61
	var v64 int32
	_ = v64
	var v67 int32
	_ = v67
	var v76 int64
	_ = v76
	var v79 int64
	_ = v79
	var v82 int64
	_ = v82
	var v84 int64
	_ = v84
	var v85 int64
	_ = v85
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	v10 = m.G0
	v12 = v10 - int32(48)
	m.G0 = v12
	v16 = m.G114
	v17 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
	m.T0[v17].(func(*base.Module, int32, int32, int32, int32))(m, l0, l1, v12+int32(24), v12)
	mBase = m.M
	if l2 == int32(0) {
	} else {
		v21 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v12)+44)))
		v23 = *(*int32)(unsafe.Add(mBase, uint32(v12)+36))
		if v23 == int32(1) {
			v26 = v21
		} else {
			v26 = int32(-1)
		}
		*(*uint16)(unsafe.Add(mBase, uint32(l2))) = uint16(v26)
	}
	if l3 == int32(0) {
	} else {
		v30 = *(*int32)(unsafe.Add(mBase, uint32(v12)+16))
		v31 = *(*int32)(unsafe.Add(mBase, uint32(v12)+20))
		*(*uint8)(unsafe.Add(mBase, uint32(l3))) = uint8(base.B2i32(v30|v31 != int32(0)))
	}
	v36 = *(*int32)(unsafe.Add(mBase, uint32(v12)+36))
	if v36 <= int32(4) {
		if v36 < int32(2) {
			v85 = int64(0)
		} else {
			switch v36 + int32(-2) {
			case 0:
				v48 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v12)+32)))
				v51 = *(*int64)(unsafe.Add(mBase, uint32(v12)+24))
				v52 = v48*int64(830472192) + v51
				if v52 < int64(0) {
					v55 = int64(-50)
				} else {
					v55 = int64(50)
				}
				v58 = base.I64_div_s(v55+v52, int64(100))
				v85 = v58
			case 1:
				v60 = int64(950)
				v61 = *(*int64)(unsafe.Add(mBase, uint32(v12)+24))
				v64 = *(*int32)(unsafe.Add(mBase, uint32(v12)+32))
				v67 = *(*int32)(unsafe.Add(mBase, uint32(v12)+40))
				v76 = v60*base.I64_extend_i32_u(v64<<(uint(int32(1))%32)-v67)<<(uint(int64(23))%64) + v61*(int64(1000)-v60)
				if v76 < int64(0) {
					v79 = int64(-500)
				} else {
					v79 = int64(500)
				}
				v82 = base.I64_div_s(v79+v76, int64(1000))
				if base.Ui64(v82) < base.Ui64(v61) {
					v84 = v61
				} else {
					v84 = v82
				}
				v85 = v84
			default:
				v60 = int64(700)
				v61 = *(*int64)(unsafe.Add(mBase, uint32(v12)+24))
				v64 = *(*int32)(unsafe.Add(mBase, uint32(v12)+32))
				v67 = *(*int32)(unsafe.Add(mBase, uint32(v12)+40))
				v76 = v60*base.I64_extend_i32_u(v64<<(uint(int32(1))%32)-v67)<<(uint(int64(23))%64) + v61*(int64(1000)-v60)
				if v76 < int64(0) {
					v79 = int64(-500)
				} else {
					v79 = int64(500)
				}
				v82 = base.I64_div_s(v79+v76, int64(1000))
				if base.Ui64(v82) < base.Ui64(v61) {
					v84 = v61
				} else {
					v84 = v82
				}
				v85 = v84
			}
		}
	} else {
		v60 = int64(627)
		v61 = *(*int64)(unsafe.Add(mBase, uint32(v12)+24))
		v64 = *(*int32)(unsafe.Add(mBase, uint32(v12)+32))
		v67 = *(*int32)(unsafe.Add(mBase, uint32(v12)+40))
		v76 = v60*base.I64_extend_i32_u(v64<<(uint(int32(1))%32)-v67)<<(uint(int64(23))%64) + v61*(int64(1000)-v60)
		if v76 < int64(0) {
			v79 = int64(-500)
		} else {
			v79 = int64(500)
		}
		v82 = base.I64_div_s(v79+v76, int64(1000))
		if base.Ui64(v82) < base.Ui64(v61) {
			v84 = v61
		} else {
			v84 = v82
		}
		v85 = v84
	}
	v87 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
	v88 = *(*int32)(unsafe.Add(mBase, uint32(v12)+12))
	v89 = *(*int32)(unsafe.Add(mBase, uint32(v12)+4))
	v90 = *(*int32)(unsafe.Add(mBase, uint32(v12)+20))
	v91 = *(*int32)(unsafe.Add(mBase, uint32(v12)+8))
	v92 = *(*int32)(unsafe.Add(mBase, uint32(v12)+16))
	m.G0 = v12 + int32(48)
	return v85 + base.I64_extend_i32_u(v88*int32(240)+v87*int32(1600)+v89*int32(2640)+v90*int32(720)+v91*int32(1840)+v92*int32(3360))<<(uint(int64(13))%64) + int64(401814323)
}
func F_PostLoopFinalize(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v29 int32
	_ = v29
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v91 int64
	_ = v91
	var v92 int64
	_ = v92
	var v94 int64
	_ = v94
	var v95 int64
	_ = v95
	var v97 int64
	_ = v97
	var v101 int64
	_ = v101
	var v103 int64
	_ = v103
	var v107 int64
	_ = v107
	var v109 int64
	_ = v109
	var v113 int64
	_ = v113
	var v115 int64
	_ = v115
	var v119 int64
	_ = v119
	var v121 int64
	_ = v121
	var v125 int64
	_ = v125
	var v127 int64
	_ = v127
	var v131 int64
	_ = v131
	var v133 int64
	_ = v133
	var v137 int64
	_ = v137
	var v139 int64
	_ = v139
	var v143 int64
	_ = v143
	var v145 int64
	_ = v145
	var v149 int64
	_ = v149
	var v151 int64
	_ = v151
	var v155 int64
	_ = v155
	var v157 int64
	_ = v157
	var v161 int64
	_ = v161
	var v174 int32
	_ = v174
	var v175 int32
	_ = v175
	var v180 float64
	_ = v180
	var v185 int32
	_ = v185
	var v188 float64
	_ = v188
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v197 float64
	_ = v197
	var v200 float64
	_ = v200
	var v203 float64
	_ = v203
	var v204 int32
	_ = v204
	var v205 float64
	_ = v205
	var v206 int32
	_ = v206
	var v207 float64
	_ = v207
	var v208 int32
	_ = v208
	var v209 float64
	_ = v209
	var v214 int32
	_ = v214
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v220 int32
	_ = v220
	var v226 float64
	_ = v226
	var v231 int32
	_ = v231
	var v234 float64
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v243 float64
	_ = v243
	var v246 float64
	_ = v246
	var v249 float64
	_ = v249
	var v250 int32
	_ = v250
	var v251 float64
	_ = v251
	var v252 int32
	_ = v252
	var v253 float64
	_ = v253
	var v254 int32
	_ = v254
	var v255 float64
	_ = v255
	var v260 int32
	_ = v260
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v272 float64
	_ = v272
	var v277 int32
	_ = v277
	var v280 float64
	_ = v280
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v289 float64
	_ = v289
	var v292 float64
	_ = v292
	var v295 float64
	_ = v295
	var v296 int32
	_ = v296
	var v297 float64
	_ = v297
	var v298 int32
	_ = v298
	var v299 float64
	_ = v299
	var v300 int32
	_ = v300
	var v301 float64
	_ = v301
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v312 int32
	_ = v312
	var v318 float64
	_ = v318
	var v323 int32
	_ = v323
	var v326 float64
	_ = v326
	var v327 int32
	_ = v327
	var v328 int32
	_ = v328
	var v335 float64
	_ = v335
	var v338 float64
	_ = v338
	var v341 float64
	_ = v341
	var v342 int32
	_ = v342
	var v343 float64
	_ = v343
	var v344 int32
	_ = v344
	var v345 float64
	_ = v345
	var v346 int32
	_ = v346
	var v347 float64
	_ = v347
	var v352 int32
	_ = v352
	var v353 int32
	_ = v353
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
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v378 int32
	_ = v378
	var v379 int32
	_ = v379
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v391 int32
	_ = v391
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v401 int32
	_ = v401
	var v403 int32
	_ = v403
	var v406 int32
	_ = v406
	var v407 int32
	_ = v407
	var v408 int32
	_ = v408
	var v409 int32
	_ = v409
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v416 int32
	_ = v416
	var v420 int32
	_ = v420
	var v421 int32
	_ = v421
	var v424 int32
	_ = v424
	var v426 int32
	_ = v426
	var v429 int32
	_ = v429
	var v430 int32
	_ = v430
	var v431 int32
	_ = v431
	var v432 int32
	_ = v432
	var v434 int32
	_ = v434
	var v435 int32
	_ = v435
	var v441 int32
	_ = v441
	var v442 int32
	_ = v442
	var v445 int32
	_ = v445
	var v447 int32
	_ = v447
	var v450 int32
	_ = v450
	var v452 int32
	_ = v452
	var v454 int32
	_ = v454
	var v458 int32
	_ = v458
	var v492 int32
	_ = v492
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v505 int32
	_ = v505
	var v506 int32
	_ = v506
	var v511 int32
	_ = v511
	var v513 int32
	_ = v513
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if l1 == int32(0) {
		F_VP8BitWriterWipeOut(m, v7+int32(56))
		mBase = m.M
		v492 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
		if v492 < int32(1) {
		} else {
			v499 = v7 + int32(88)
			v500 = int32(0)
			for {
				F_VP8BitWriterWipeOut(m, v499)
				mBase = m.M
				v505 = v500 + int32(1)
				v506 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
				if v505 < v506 {
					v499 = v499 + int32(32)
					v500 = v505
					continue
				} else {
					break
				}
				break
			}
		}
		v511 = *(*int32)(unsafe.Add(mBase, uint32(v7)+4))
		v513 = *(*int32)(unsafe.Add(mBase, uint32(v511)+92))
		if v513 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(v511)+92)) = int32(1)
		}
		return int32(0)
	} else {
		v10 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
		if v10 < int32(1) {
			v82 = l1
			v87 = *(*int32)(unsafe.Add(mBase, uint32(v7)+4))
			v88 = *(*int32)(unsafe.Add(mBase, uint32(v87)+88))
			if v88 == int32(0) {
			} else {
				v91 = *(*int64)(unsafe.Add(mBase, uint32(l0)+168))
				v92 = int64(7)
				v94 = int64(3)
				v95 = int64(base.Ui64(v91+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[0]))) = uint32(v95)
				v97 = *(*int64)(unsafe.Add(mBase, uint32(l0)+192))
				v101 = int64(base.Ui64(v97+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[1]))) = uint32(v101)
				v103 = *(*int64)(unsafe.Add(mBase, uint32(l0)+216))
				v107 = int64(base.Ui64(v103+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[2]))) = uint32(v107)
				v109 = *(*int64)(unsafe.Add(mBase, uint32(l0)+240))
				v113 = int64(base.Ui64(v109+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[3]))) = uint32(v113)
				v115 = *(*int64)(unsafe.Add(mBase, uint32(l0)+176))
				v119 = int64(base.Ui64(v115+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[4]))) = uint32(v119)
				v121 = *(*int64)(unsafe.Add(mBase, uint32(l0)+200))
				v125 = int64(base.Ui64(v121+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[5]))) = uint32(v125)
				v127 = *(*int64)(unsafe.Add(mBase, uint32(l0)+224))
				v131 = int64(base.Ui64(v127+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[6]))) = uint32(v131)
				v133 = *(*int64)(unsafe.Add(mBase, uint32(l0)+248))
				v137 = int64(base.Ui64(v133+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[7]))) = uint32(v137)
				v139 = *(*int64)(unsafe.Add(mBase, uint32(l0)+184))
				v143 = int64(base.Ui64(v139+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[8]))) = uint32(v143)
				v145 = *(*int64)(unsafe.Add(mBase, uint32(l0)+208))
				v149 = int64(base.Ui64(v145+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[9]))) = uint32(v149)
				v151 = *(*int64)(unsafe.Add(mBase, uint32(l0)+232))
				v155 = int64(base.Ui64(v151+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[10]))) = uint32(v155)
				v157 = *(*int64)(unsafe.Add(mBase, uint32(l0)+256))
				v161 = int64(base.Ui64(v157+v92) >> (uint(v94) % 64))
				*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[11]))) = uint32(v161)
			}
			v174 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
			v175 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
			if v175 == int32(0) {
				v362 = *(*int32)(unsafe.Add(mBase, uint32(v174)))
				v363 = *(*int32)(unsafe.Add(mBase, uint32(v362)+32))
				if v363 < int32(1) {
				} else {
					v366 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+634)))
					v367 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1096))
					v368 = m.G1
					v369 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1092))
					v372 = *(*int32)(unsafe.Add(mBase, uint32(v174)+16))
					v374 = v372 << (uint(int32(6)) % 32)
					v378 = v367 * v366 >> (uint(int32(3)) % 32)
					v379 = int32(63)
					if v378 < v379 {
						v382 = v378
					} else {
						v382 = v379
					}
					v384 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v368+int32(_a_F_PostLoopFinalize_0)+v374+v382))))
					if v384 <= v369 {
						v387 = v369
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v174)+1092)) = v384
						v387 = v384
					}
					v388 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+1378)))
					v389 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1840))
					v390 = m.G1
					v391 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1836))
					v397 = v389 * v388 >> (uint(int32(3)) % 32)
					v398 = int32(63)
					if v397 < v398 {
						v401 = v397
					} else {
						v401 = v398
					}
					v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v390+int32(_a_F_PostLoopFinalize_0)+v374+v401))))
					if v403 <= v391 {
						v406 = v391
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v174)+1836)) = v403
						v406 = v403
					}
					v407 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+2122)))
					v408 = *(*int32)(unsafe.Add(mBase, uint32(v174)+2584))
					v409 = m.G1
					if base.Ui32(v406) < base.Ui32(v387) {
						v411 = v387
					} else {
						v411 = v406
					}
					v412 = *(*int32)(unsafe.Add(mBase, uint32(v174)+2580))
					v416 = v372 << (uint(int32(6)) % 32)
					v420 = v408 * v407 >> (uint(int32(3)) % 32)
					v421 = int32(63)
					if v420 < v421 {
						v424 = v420
					} else {
						v424 = v421
					}
					v426 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v409+int32(_a_F_PostLoopFinalize_0)+v416+v424))))
					if v426 <= v412 {
						v429 = v412
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v174)+2580)) = v426
						v429 = v426
					}
					v430 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+2866)))
					v431 = *(*int32)(unsafe.Add(mBase, uint32(v174)+3328))
					v432 = m.G1
					if base.Ui32(v429) < base.Ui32(v411) {
						v434 = v411
					} else {
						v434 = v429
					}
					v435 = *(*int32)(unsafe.Add(mBase, uint32(v174)+3324))
					v441 = v431 * v430 >> (uint(int32(3)) % 32)
					v442 = int32(63)
					if v441 < v442 {
						v445 = v441
					} else {
						v445 = v442
					}
					v447 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v432+int32(_a_F_PostLoopFinalize_0)+v416+v445))))
					if v447 <= v435 {
						v450 = v435
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(v174)+3324)) = v447
						v450 = v447
					}
					if base.Ui32(v450) < base.Ui32(v434) {
						v452 = v434
					} else {
						v452 = v450
					}
					v454 = int32(12)
					v458 = v452
					*(*int32)(unsafe.Add(mBase, uint32(v174+v454))) = v458
				}
			} else {
				v180 = *(*float64)(unsafe.Add(mBase, uint32(v175)))
				v185 = v175 + int32(24)
				v188 = base.F64_mul(v180, float64(1.00001))
				v189 = int32(0)
				v190 = int32(1)
				for {
					v197 = *(*float64)(unsafe.Add(mBase, uint32(v185)))
					v200 = *(*float64)(unsafe.Add(mBase, uint32(v185+int32(-8))))
					v203 = *(*float64)(unsafe.Add(mBase, uint32(v185+int32(-16))))
					v204 = base.F64_gt(v203, v188)
					if v204 != 0 {
						v205 = v203
					} else {
						v205 = v188
					}
					v206 = base.F64_gt(v200, v205)
					if v206 != 0 {
						v207 = v200
					} else {
						v207 = v205
					}
					v208 = base.F64_gt(v197, v207)
					if v208 != 0 {
						v209 = v197
					} else {
						v209 = v207
					}
					if v204 != 0 {
						v214 = v190
					} else {
						v214 = v189
					}
					if v206 != 0 {
						v215 = v190 + int32(1)
					} else {
						v215 = v214
					}
					if v208 != 0 {
						v216 = v190 + int32(2)
					} else {
						v216 = v215
					}
					v220 = v190 + int32(3)
					if v220 != int32(64) {
						v185 = v185 + int32(24)
						v188 = v209
						v189 = v216
						v190 = v220
						continue
					} else {
						break
					}
					break
				}
				*(*int32)(unsafe.Add(mBase, uint32(v174)+1092)) = v216
				v226 = *(*float64)(unsafe.Add(mBase, uint32(v175)+512))
				v231 = v175 + int32(536)
				v234 = base.F64_mul(v226, float64(1.00001))
				v235 = int32(0)
				v236 = int32(1)
				for {
					v243 = *(*float64)(unsafe.Add(mBase, uint32(v231)))
					v246 = *(*float64)(unsafe.Add(mBase, uint32(v231+int32(-8))))
					v249 = *(*float64)(unsafe.Add(mBase, uint32(v231+int32(-16))))
					v250 = base.F64_gt(v249, v234)
					if v250 != 0 {
						v251 = v249
					} else {
						v251 = v234
					}
					v252 = base.F64_gt(v246, v251)
					if v252 != 0 {
						v253 = v246
					} else {
						v253 = v251
					}
					v254 = base.F64_gt(v243, v253)
					if v254 != 0 {
						v255 = v243
					} else {
						v255 = v253
					}
					if v250 != 0 {
						v260 = v236
					} else {
						v260 = v235
					}
					if v252 != 0 {
						v261 = v236 + int32(1)
					} else {
						v261 = v260
					}
					if v254 != 0 {
						v262 = v236 + int32(2)
					} else {
						v262 = v261
					}
					v266 = v236 + int32(3)
					if v266 != int32(64) {
						v231 = v231 + int32(24)
						v234 = v255
						v235 = v262
						v236 = v266
						continue
					} else {
						break
					}
					break
				}
				*(*int32)(unsafe.Add(mBase, uint32(v174)+1836)) = v262
				v272 = *(*float64)(unsafe.Add(mBase, uint32(v175)+1024))
				v277 = v175 + int32(1048)
				v280 = base.F64_mul(v272, float64(1.00001))
				v281 = int32(0)
				v282 = int32(1)
				for {
					v289 = *(*float64)(unsafe.Add(mBase, uint32(v277)))
					v292 = *(*float64)(unsafe.Add(mBase, uint32(v277+int32(-8))))
					v295 = *(*float64)(unsafe.Add(mBase, uint32(v277+int32(-16))))
					v296 = base.F64_gt(v295, v280)
					if v296 != 0 {
						v297 = v295
					} else {
						v297 = v280
					}
					v298 = base.F64_gt(v292, v297)
					if v298 != 0 {
						v299 = v292
					} else {
						v299 = v297
					}
					v300 = base.F64_gt(v289, v299)
					if v300 != 0 {
						v301 = v289
					} else {
						v301 = v299
					}
					if v296 != 0 {
						v306 = v282
					} else {
						v306 = v281
					}
					if v298 != 0 {
						v307 = v282 + int32(1)
					} else {
						v307 = v306
					}
					if v300 != 0 {
						v308 = v282 + int32(2)
					} else {
						v308 = v307
					}
					v312 = v282 + int32(3)
					if v312 != int32(64) {
						v277 = v277 + int32(24)
						v280 = v301
						v281 = v308
						v282 = v312
						continue
					} else {
						break
					}
					break
				}
				*(*int32)(unsafe.Add(mBase, uint32(v174)+2580)) = v308
				v318 = *(*float64)(unsafe.Add(mBase, uint32(v175)+1536))
				v323 = v175 + int32(1560)
				v326 = base.F64_mul(v318, float64(1.00001))
				v327 = int32(0)
				v328 = int32(1)
				for {
					v335 = *(*float64)(unsafe.Add(mBase, uint32(v323)))
					v338 = *(*float64)(unsafe.Add(mBase, uint32(v323+int32(-8))))
					v341 = *(*float64)(unsafe.Add(mBase, uint32(v323+int32(-16))))
					v342 = base.F64_gt(v341, v326)
					if v342 != 0 {
						v343 = v341
					} else {
						v343 = v326
					}
					v344 = base.F64_gt(v338, v343)
					if v344 != 0 {
						v345 = v338
					} else {
						v345 = v343
					}
					v346 = base.F64_gt(v335, v345)
					if v346 != 0 {
						v347 = v335
					} else {
						v347 = v345
					}
					if v342 != 0 {
						v352 = v328
					} else {
						v352 = v327
					}
					if v344 != 0 {
						v353 = v328 + int32(1)
					} else {
						v353 = v352
					}
					if v346 != 0 {
						v354 = v328 + int32(2)
					} else {
						v354 = v353
					}
					v358 = v328 + int32(3)
					if v358 != int32(64) {
						v323 = v323 + int32(24)
						v326 = v347
						v327 = v354
						v328 = v358
						continue
					} else {
						break
					}
					break
				}
				v454 = int32(3324)
				v458 = v354
				*(*int32)(unsafe.Add(mBase, uint32(v174+v454))) = v458
			}
			return v82
		} else {
			v21 = l1 & int32(1)
			v22 = v7 + int32(88)
			v23 = int32(0)
			for {
				v29 = *(*int32)(unsafe.Add(mBase, uint32(v22)+12))
				v33 = int32(1) << (uint(int32(8)-v29) % 32)
				for {
					v36 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
					v38 = v36 >> (uint(int32(1)) % 32)
					*(*int32)(unsafe.Add(mBase, uint32(v22))) = v38
					if int32(126) < v38 {
					} else {
						v42 = *(*int32)(unsafe.Add(mBase, uint32(v22)+4))
						v43 = int32(1)
						*(*int32)(unsafe.Add(mBase, uint32(v22)+4)) = v42 << (uint(v43) % 32)
						v46 = *(*int32)(unsafe.Add(mBase, uint32(v22)+12))
						*(*int32)(unsafe.Add(mBase, uint32(v22)+12)) = v46 + v43
						v50 = m.G1
						v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v50+int32(_a_F_PostLoopFinalize_1)+v38))))
						*(*int32)(unsafe.Add(mBase, uint32(v22))) = v54
						if v46 < int32(0) {
						} else {
							F_Flush(m, v22)
							mBase = m.M
						}
					}
					v60 = int32(1)
					if base.Ui32(v60) < base.Ui32(v33) {
						v33 = int32(base.Ui32(v33) >> (uint(v60) % 32))
						continue
					} else {
						break
					}
					break
				}
				*(*int32)(unsafe.Add(mBase, uint32(v22)+12)) = int32(0)
				F_Flush(m, v22)
				mBase = m.M
				v71 = *(*int32)(unsafe.Add(mBase, uint32(v22+int32(28))))
				if v71 != 0 {
					v72 = int32(0)
				} else {
					v72 = v21
				}
				v76 = v23 + int32(1)
				v77 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
				if v76 < v77 {
					v21 = v72
					v22 = v22 + int32(32)
					v23 = v76
					continue
				} else {
					break
				}
				break
			}
			if v72 == int32(0) {
				F_VP8BitWriterWipeOut(m, v7+int32(56))
				mBase = m.M
				v492 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
				if v492 < int32(1) {
				} else {
					v499 = v7 + int32(88)
					v500 = int32(0)
					for {
						F_VP8BitWriterWipeOut(m, v499)
						mBase = m.M
						v505 = v500 + int32(1)
						v506 = *(*int32)(unsafe.Add(mBase, uint32(v7)+52))
						if v505 < v506 {
							v499 = v499 + int32(32)
							v500 = v505
							continue
						} else {
							break
						}
						break
					}
				}
				v511 = *(*int32)(unsafe.Add(mBase, uint32(v7)+4))
				v513 = *(*int32)(unsafe.Add(mBase, uint32(v511)+92))
				if v513 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v511)+92)) = int32(1)
				}
				return int32(0)
			} else {
				v82 = v21
				v87 = *(*int32)(unsafe.Add(mBase, uint32(v7)+4))
				v88 = *(*int32)(unsafe.Add(mBase, uint32(v87)+88))
				if v88 == int32(0) {
				} else {
					v91 = *(*int64)(unsafe.Add(mBase, uint32(l0)+168))
					v92 = int64(7)
					v94 = int64(3)
					v95 = int64(base.Ui64(v91+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[0]))) = uint32(v95)
					v97 = *(*int64)(unsafe.Add(mBase, uint32(l0)+192))
					v101 = int64(base.Ui64(v97+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[1]))) = uint32(v101)
					v103 = *(*int64)(unsafe.Add(mBase, uint32(l0)+216))
					v107 = int64(base.Ui64(v103+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[2]))) = uint32(v107)
					v109 = *(*int64)(unsafe.Add(mBase, uint32(l0)+240))
					v113 = int64(base.Ui64(v109+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[3]))) = uint32(v113)
					v115 = *(*int64)(unsafe.Add(mBase, uint32(l0)+176))
					v119 = int64(base.Ui64(v115+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[4]))) = uint32(v119)
					v121 = *(*int64)(unsafe.Add(mBase, uint32(l0)+200))
					v125 = int64(base.Ui64(v121+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[5]))) = uint32(v125)
					v127 = *(*int64)(unsafe.Add(mBase, uint32(l0)+224))
					v131 = int64(base.Ui64(v127+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[6]))) = uint32(v131)
					v133 = *(*int64)(unsafe.Add(mBase, uint32(l0)+248))
					v137 = int64(base.Ui64(v133+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[7]))) = uint32(v137)
					v139 = *(*int64)(unsafe.Add(mBase, uint32(l0)+184))
					v143 = int64(base.Ui64(v139+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[8]))) = uint32(v143)
					v145 = *(*int64)(unsafe.Add(mBase, uint32(l0)+208))
					v149 = int64(base.Ui64(v145+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[9]))) = uint32(v149)
					v151 = *(*int64)(unsafe.Add(mBase, uint32(l0)+232))
					v155 = int64(base.Ui64(v151+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[10]))) = uint32(v155)
					v157 = *(*int64)(unsafe.Add(mBase, uint32(l0)+256))
					v161 = int64(base.Ui64(v157+v92) >> (uint(v94) % 64))
					*(*uint32)(unsafe.Add(mBase, uint32(v7)+uint32(_c_F_PostLoopFinalize[11]))) = uint32(v161)
				}
				v174 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
				v175 = *(*int32)(unsafe.Add(mBase, uint32(l0)+280))
				if v175 == int32(0) {
					v362 = *(*int32)(unsafe.Add(mBase, uint32(v174)))
					v363 = *(*int32)(unsafe.Add(mBase, uint32(v362)+32))
					if v363 < int32(1) {
					} else {
						v366 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+634)))
						v367 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1096))
						v368 = m.G1
						v369 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1092))
						v372 = *(*int32)(unsafe.Add(mBase, uint32(v174)+16))
						v374 = v372 << (uint(int32(6)) % 32)
						v378 = v367 * v366 >> (uint(int32(3)) % 32)
						v379 = int32(63)
						if v378 < v379 {
							v382 = v378
						} else {
							v382 = v379
						}
						v384 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v368+int32(_a_F_PostLoopFinalize_0)+v374+v382))))
						if v384 <= v369 {
							v387 = v369
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v174)+1092)) = v384
							v387 = v384
						}
						v388 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+1378)))
						v389 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1840))
						v390 = m.G1
						v391 = *(*int32)(unsafe.Add(mBase, uint32(v174)+1836))
						v397 = v389 * v388 >> (uint(int32(3)) % 32)
						v398 = int32(63)
						if v397 < v398 {
							v401 = v397
						} else {
							v401 = v398
						}
						v403 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v390+int32(_a_F_PostLoopFinalize_0)+v374+v401))))
						if v403 <= v391 {
							v406 = v391
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v174)+1836)) = v403
							v406 = v403
						}
						v407 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+2122)))
						v408 = *(*int32)(unsafe.Add(mBase, uint32(v174)+2584))
						v409 = m.G1
						if base.Ui32(v406) < base.Ui32(v387) {
							v411 = v387
						} else {
							v411 = v406
						}
						v412 = *(*int32)(unsafe.Add(mBase, uint32(v174)+2580))
						v416 = v372 << (uint(int32(6)) % 32)
						v420 = v408 * v407 >> (uint(int32(3)) % 32)
						v421 = int32(63)
						if v420 < v421 {
							v424 = v420
						} else {
							v424 = v421
						}
						v426 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v409+int32(_a_F_PostLoopFinalize_0)+v416+v424))))
						if v426 <= v412 {
							v429 = v412
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v174)+2580)) = v426
							v429 = v426
						}
						v430 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v174)+2866)))
						v431 = *(*int32)(unsafe.Add(mBase, uint32(v174)+3328))
						v432 = m.G1
						if base.Ui32(v429) < base.Ui32(v411) {
							v434 = v411
						} else {
							v434 = v429
						}
						v435 = *(*int32)(unsafe.Add(mBase, uint32(v174)+3324))
						v441 = v431 * v430 >> (uint(int32(3)) % 32)
						v442 = int32(63)
						if v441 < v442 {
							v445 = v441
						} else {
							v445 = v442
						}
						v447 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v432+int32(_a_F_PostLoopFinalize_0)+v416+v445))))
						if v447 <= v435 {
							v450 = v435
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v174)+3324)) = v447
							v450 = v447
						}
						if base.Ui32(v450) < base.Ui32(v434) {
							v452 = v434
						} else {
							v452 = v450
						}
						v454 = int32(12)
						v458 = v452
						*(*int32)(unsafe.Add(mBase, uint32(v174+v454))) = v458
					}
				} else {
					v180 = *(*float64)(unsafe.Add(mBase, uint32(v175)))
					v185 = v175 + int32(24)
					v188 = base.F64_mul(v180, float64(1.00001))
					v189 = int32(0)
					v190 = int32(1)
					for {
						v197 = *(*float64)(unsafe.Add(mBase, uint32(v185)))
						v200 = *(*float64)(unsafe.Add(mBase, uint32(v185+int32(-8))))
						v203 = *(*float64)(unsafe.Add(mBase, uint32(v185+int32(-16))))
						v204 = base.F64_gt(v203, v188)
						if v204 != 0 {
							v205 = v203
						} else {
							v205 = v188
						}
						v206 = base.F64_gt(v200, v205)
						if v206 != 0 {
							v207 = v200
						} else {
							v207 = v205
						}
						v208 = base.F64_gt(v197, v207)
						if v208 != 0 {
							v209 = v197
						} else {
							v209 = v207
						}
						if v204 != 0 {
							v214 = v190
						} else {
							v214 = v189
						}
						if v206 != 0 {
							v215 = v190 + int32(1)
						} else {
							v215 = v214
						}
						if v208 != 0 {
							v216 = v190 + int32(2)
						} else {
							v216 = v215
						}
						v220 = v190 + int32(3)
						if v220 != int32(64) {
							v185 = v185 + int32(24)
							v188 = v209
							v189 = v216
							v190 = v220
							continue
						} else {
							break
						}
						break
					}
					*(*int32)(unsafe.Add(mBase, uint32(v174)+1092)) = v216
					v226 = *(*float64)(unsafe.Add(mBase, uint32(v175)+512))
					v231 = v175 + int32(536)
					v234 = base.F64_mul(v226, float64(1.00001))
					v235 = int32(0)
					v236 = int32(1)
					for {
						v243 = *(*float64)(unsafe.Add(mBase, uint32(v231)))
						v246 = *(*float64)(unsafe.Add(mBase, uint32(v231+int32(-8))))
						v249 = *(*float64)(unsafe.Add(mBase, uint32(v231+int32(-16))))
						v250 = base.F64_gt(v249, v234)
						if v250 != 0 {
							v251 = v249
						} else {
							v251 = v234
						}
						v252 = base.F64_gt(v246, v251)
						if v252 != 0 {
							v253 = v246
						} else {
							v253 = v251
						}
						v254 = base.F64_gt(v243, v253)
						if v254 != 0 {
							v255 = v243
						} else {
							v255 = v253
						}
						if v250 != 0 {
							v260 = v236
						} else {
							v260 = v235
						}
						if v252 != 0 {
							v261 = v236 + int32(1)
						} else {
							v261 = v260
						}
						if v254 != 0 {
							v262 = v236 + int32(2)
						} else {
							v262 = v261
						}
						v266 = v236 + int32(3)
						if v266 != int32(64) {
							v231 = v231 + int32(24)
							v234 = v255
							v235 = v262
							v236 = v266
							continue
						} else {
							break
						}
						break
					}
					*(*int32)(unsafe.Add(mBase, uint32(v174)+1836)) = v262
					v272 = *(*float64)(unsafe.Add(mBase, uint32(v175)+1024))
					v277 = v175 + int32(1048)
					v280 = base.F64_mul(v272, float64(1.00001))
					v281 = int32(0)
					v282 = int32(1)
					for {
						v289 = *(*float64)(unsafe.Add(mBase, uint32(v277)))
						v292 = *(*float64)(unsafe.Add(mBase, uint32(v277+int32(-8))))
						v295 = *(*float64)(unsafe.Add(mBase, uint32(v277+int32(-16))))
						v296 = base.F64_gt(v295, v280)
						if v296 != 0 {
							v297 = v295
						} else {
							v297 = v280
						}
						v298 = base.F64_gt(v292, v297)
						if v298 != 0 {
							v299 = v292
						} else {
							v299 = v297
						}
						v300 = base.F64_gt(v289, v299)
						if v300 != 0 {
							v301 = v289
						} else {
							v301 = v299
						}
						if v296 != 0 {
							v306 = v282
						} else {
							v306 = v281
						}
						if v298 != 0 {
							v307 = v282 + int32(1)
						} else {
							v307 = v306
						}
						if v300 != 0 {
							v308 = v282 + int32(2)
						} else {
							v308 = v307
						}
						v312 = v282 + int32(3)
						if v312 != int32(64) {
							v277 = v277 + int32(24)
							v280 = v301
							v281 = v308
							v282 = v312
							continue
						} else {
							break
						}
						break
					}
					*(*int32)(unsafe.Add(mBase, uint32(v174)+2580)) = v308
					v318 = *(*float64)(unsafe.Add(mBase, uint32(v175)+1536))
					v323 = v175 + int32(1560)
					v326 = base.F64_mul(v318, float64(1.00001))
					v327 = int32(0)
					v328 = int32(1)
					for {
						v335 = *(*float64)(unsafe.Add(mBase, uint32(v323)))
						v338 = *(*float64)(unsafe.Add(mBase, uint32(v323+int32(-8))))
						v341 = *(*float64)(unsafe.Add(mBase, uint32(v323+int32(-16))))
						v342 = base.F64_gt(v341, v326)
						if v342 != 0 {
							v343 = v341
						} else {
							v343 = v326
						}
						v344 = base.F64_gt(v338, v343)
						if v344 != 0 {
							v345 = v338
						} else {
							v345 = v343
						}
						v346 = base.F64_gt(v335, v345)
						if v346 != 0 {
							v347 = v335
						} else {
							v347 = v345
						}
						if v342 != 0 {
							v352 = v328
						} else {
							v352 = v327
						}
						if v344 != 0 {
							v353 = v328 + int32(1)
						} else {
							v353 = v352
						}
						if v346 != 0 {
							v354 = v328 + int32(2)
						} else {
							v354 = v353
						}
						v358 = v328 + int32(3)
						if v358 != int32(64) {
							v323 = v323 + int32(24)
							v326 = v347
							v327 = v354
							v328 = v358
							continue
						} else {
							break
						}
						break
					}
					v454 = int32(3324)
					v458 = v354
					*(*int32)(unsafe.Add(mBase, uint32(v174+v454))) = v458
				}
				return v82
			}
		}
	}
}
func F_PredictorAdd10_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v44 int32
	_ = v44
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v74 int32
	_ = v74
	if l2 < int32(1) {
	} else {
		v11 = int32(-4)
		v15 = *(*int32)(unsafe.Add(mBase, uint32(l3+v11)))
		v16 = l0
		v17 = l1 + v11
		v18 = l2
		v19 = l3
		v20 = v15
		for {
			v26 = *(*int32)(unsafe.Add(mBase, uint32(v17+int32(8))))
			v27 = int32(4)
			v28 = v17 + v27
			v29 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
			v31 = int32(1)
			v33 = int32(2139062143)
			v36 = int32(base.Ui32(v26^v29)>>(uint(v31)%32))&v33 + v26&v29
			v37 = *(*int32)(unsafe.Add(mBase, uint32(v17)))
			v44 = int32(base.Ui32(v37^v20)>>(uint(v31)%32))&v33 + v37&v20
			v51 = int32(base.Ui32(v36^v44)>>(uint(v31)%32))&v33 + v36&v44
			v52 = int32(-16711936)
			v54 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
			v60 = int32(16711935)
			v67 = (v51&v52+v54&v52)&v52 | (v51&v60+v54&v60)&v60
			*(*int32)(unsafe.Add(mBase, uint32(v19))) = v67
			v74 = v18 + int32(-1)
			if v74 != 0 {
				v16 = v16 + v27
				v17 = v28
				v18 = v74
				v19 = v19 + v27
				v20 = v67
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorAdd11_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v92 int32
	_ = v92
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v128 int32
	_ = v128
	var v135 int32
	_ = v135
	var v137 int32
	_ = v137
	var v142 int32
	_ = v142
	if l2 < int32(1) {
	} else {
		v14 = int32(-4)
		v18 = *(*int32)(unsafe.Add(mBase, uint32(l3+v14)))
		v19 = l0
		v20 = v18
		v21 = l2
		v22 = l3
		v23 = l1 + v14
		for {
			v30 = *(*int32)(unsafe.Add(mBase, uint32(v23)))
			v32 = v23 + int32(4)
			v33 = *(*int32)(unsafe.Add(mBase, uint32(v32)))
			v34 = int32(255)
			v37 = v30 & v34
			v38 = v20&v34 - v37
			v39 = int32(31)
			v40 = v38 >> (uint(v39) % 32)
			v43 = int32(24)
			v46 = int32(base.Ui32(v30) >> (uint(v43) % 32))
			v47 = int32(base.Ui32(v20)>>(uint(v43)%32)) - v46
			v49 = v47 >> (uint(v39) % 32)
			v53 = int32(8)
			v60 = int32(base.Ui32(v30)>>(uint(v53)%32)) & v34
			v61 = int32(base.Ui32(v20)>>(uint(v53)%32))&v34 - v60
			v63 = v61 >> (uint(v39) % 32)
			v69 = v33&v34 - v37
			v71 = v69 >> (uint(v39) % 32)
			v76 = int32(base.Ui32(v33)>>(uint(v43)%32)) - v46
			v78 = v76 >> (uint(v39) % 32)
			v86 = int32(base.Ui32(v33)>>(uint(v53)%32))&v34 - v60
			v88 = v86 >> (uint(v39) % 32)
			v92 = int32(16)
			v99 = int32(base.Ui32(v30)>>(uint(v92)%32)) & v34
			v100 = int32(base.Ui32(v33)>>(uint(v92)%32))&v34 - v99
			v102 = v100 >> (uint(v39) % 32)
			v111 = int32(base.Ui32(v20)>>(uint(v92)%32))&v34 - v99
			v113 = v111 >> (uint(v39) % 32)
			if v38^v40-v40+(v47^v49-v49)+(v61^v63-v63)-(v69^v71-v71+(v76^v78-v78)+(v86^v88-v88)+(v100^v102-v102))+(v111^v113-v113) < int32(1) {
				v119 = v33
			} else {
				v119 = v20
			}
			v120 = int32(-16711936)
			v122 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
			v128 = int32(16711935)
			v135 = (v119&v120+v122&v120)&v120 | (v119&v128+v122&v128)&v128
			*(*int32)(unsafe.Add(mBase, uint32(v22))) = v135
			v137 = int32(4)
			v142 = v21 + int32(-1)
			if v142 != 0 {
				v19 = v19 + v137
				v20 = v135
				v21 = v142
				v22 = v22 + v137
				v23 = v32
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorAdd12_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v52 int32
	_ = v52
	var v59 int32
	_ = v59
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v74 int32
	_ = v74
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v117 int32
	_ = v117
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v131 int32
	_ = v131
	if l2 < int32(1) {
	} else {
		v11 = int32(-4)
		v15 = *(*int32)(unsafe.Add(mBase, uint32(l3+v11)))
		v16 = l0
		v17 = v15
		v18 = l2
		v19 = l3
		v20 = l1 + v11
		for {
			v25 = v20 + int32(4)
			v26 = *(*int32)(unsafe.Add(mBase, uint32(v25)))
			v27 = int32(24)
			v32 = *(*int32)(unsafe.Add(mBase, uint32(v20)))
			v35 = int32(base.Ui32(v26)>>(uint(v27)%32)) + int32(base.Ui32(v17)>>(uint(v27)%32)) - int32(base.Ui32(v32)>>(uint(v27)%32))
			if base.Ui32(v35) < base.Ui32(int32(256)) {
				v42 = v35
			} else {
				v42 = int32(base.Ui32(v35^int32(-1)) >> (uint(v27) % 32))
			}
			v43 = int32(24)
			v45 = int32(255)
			v52 = v26&v45 + v17&v45 - v32&v45
			if base.Ui32(v52) < base.Ui32(int32(256)) {
				v59 = v52
			} else {
				v59 = int32(base.Ui32(v52^int32(-1)) >> (uint(v43) % 32))
			}
			v61 = int32(16)
			v63 = int32(255)
			v74 = int32(base.Ui32(v26)>>(uint(v61)%32))&v63 + int32(base.Ui32(v17)>>(uint(v61)%32))&v63 - int32(base.Ui32(v32)>>(uint(v61)%32))&v63
			if base.Ui32(v74) < base.Ui32(int32(256)) {
				v81 = v74
			} else {
				v81 = int32(base.Ui32(v74^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v85 = int32(8)
			v87 = int32(255)
			v98 = int32(base.Ui32(v26)>>(uint(v85)%32))&v87 + int32(base.Ui32(v17)>>(uint(v85)%32))&v87 - int32(base.Ui32(v32)>>(uint(v85)%32))&v87
			if base.Ui32(v98) < base.Ui32(int32(256)) {
				v105 = v98
			} else {
				v105 = int32(base.Ui32(v98^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v108 = v42<<(uint(v43)%32) + v59 + v81<<(uint(int32(16))%32) + v105<<(uint(int32(8))%32)
			v109 = int32(-16711936)
			v111 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
			v117 = int32(16711935)
			v124 = (v108&v109+v111&v109)&v109 | (v108&v117+v111&v117)&v117
			*(*int32)(unsafe.Add(mBase, uint32(v19))) = v124
			v126 = int32(4)
			v131 = v18 + int32(-1)
			if v131 != 0 {
				v16 = v16 + v126
				v17 = v124
				v18 = v131
				v19 = v19 + v126
				v20 = v25
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorAdd13_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
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
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v122 int32
	_ = v122
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	if l2 < int32(1) {
	} else {
		v10 = int32(-4)
		v14 = *(*int32)(unsafe.Add(mBase, uint32(l3+v10)))
		v15 = l0
		v16 = v14
		v17 = l2
		v18 = l3
		v19 = l1 + v10
		for {
			v23 = v19 + int32(4)
			v24 = *(*int32)(unsafe.Add(mBase, uint32(v23)))
			v31 = int32(base.Ui32(v24^v16)>>(uint(int32(1))%32))&int32(2139062143) + v24&v16
			v32 = int32(24)
			v33 = int32(base.Ui32(v31) >> (uint(v32) % 32))
			v34 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
			v39 = base.I32_div_s(v33-int32(base.Ui32(v34)>>(uint(v32)%32)), int32(2))
			v41 = v33 + base.I32_extend16_s(v39)
			if base.Ui32(v41) < base.Ui32(int32(256)) {
				v48 = v41
			} else {
				v48 = int32(base.Ui32(v41^int32(-1)) >> (uint(v32) % 32))
			}
			v49 = int32(24)
			v51 = int32(255)
			v52 = v31 & v51
			v57 = base.I32_div_s(v52-v34&v51, int32(2))
			v59 = v52 + base.I32_extend16_s(v57)
			if base.Ui32(v59) < base.Ui32(int32(256)) {
				v66 = v59
			} else {
				v66 = int32(base.Ui32(v59^int32(-1)) >> (uint(v49) % 32))
			}
			v68 = int32(16)
			v70 = int32(255)
			v71 = int32(base.Ui32(v31)>>(uint(v68)%32)) & v70
			v78 = base.I32_div_s(v71-int32(base.Ui32(v34)>>(uint(v68)%32))&v70, int32(2))
			v80 = v71 + base.I32_extend16_s(v78)
			if base.Ui32(v80) < base.Ui32(int32(256)) {
				v87 = v80
			} else {
				v87 = int32(base.Ui32(v80^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v91 = int32(8)
			v93 = int32(255)
			v94 = int32(base.Ui32(v31)>>(uint(v91)%32)) & v93
			v101 = base.I32_div_s(v94-int32(base.Ui32(v34)>>(uint(v91)%32))&v93, int32(2))
			v103 = v94 + base.I32_extend16_s(v101)
			if base.Ui32(v103) < base.Ui32(int32(256)) {
				v110 = v103
			} else {
				v110 = int32(base.Ui32(v103^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v113 = v48<<(uint(v49)%32) + v66 + v87<<(uint(int32(16))%32) + v110<<(uint(int32(8))%32)
			v114 = int32(-16711936)
			v116 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
			v122 = int32(16711935)
			v129 = (v113&v114+v116&v114)&v114 | (v113&v122+v116&v122)&v122
			*(*int32)(unsafe.Add(mBase, uint32(v18))) = v129
			v131 = int32(4)
			v136 = v17 + int32(-1)
			if v136 != 0 {
				v15 = v15 + v131
				v16 = v129
				v17 = v136
				v18 = v18 + v131
				v19 = v23
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorAdd1_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v102 int32
	_ = v102
	if l2 < int32(1) {
	} else {
		v14 = int32(1)
		v18 = *(*int32)(unsafe.Add(mBase, uint32(l3+int32(-4))))
		if l2 == v14 {
			v82 = v18
			v83 = int32(0)
		} else {
			v27 = l3
			v30 = v18
			v31 = int32(0)
			v33 = l0
			for {
				v36 = *(*int32)(unsafe.Add(mBase, uint32(v33)))
				v37 = int32(-16711936)
				v43 = (v36&v37 + v30&v37) & v37
				v44 = int32(16711935)
				v50 = (v36&v44 + v30&v44) & v44
				*(*int32)(unsafe.Add(mBase, uint32(v27))) = v43 | v50
				v53 = int32(4)
				v57 = *(*int32)(unsafe.Add(mBase, uint32(v33+v53)))
				v68 = (v57&v37+v43)&v37 | (v57&v44+v50)&v44
				*(*int32)(unsafe.Add(mBase, uint32(v27+v53))) = v68
				v70 = int32(8)
				v75 = v31 + int32(2)
				if l2&int32(2147483646) != v75 {
					v27 = v27 + v70
					v30 = v68
					v31 = v75
					v33 = v33 + v70
					continue
				} else {
					break
				}
				break
			}
			v82 = v68
			v83 = v75
		}
		if l2&v14 == int32(0) {
		} else {
			v91 = v83 << (uint(int32(2)) % 32)
			v94 = *(*int32)(unsafe.Add(mBase, uint32(l0+v91)))
			v95 = int32(-16711936)
			v102 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(l3+v91))) = (v94&v95+v82&v95)&v95 | (v94&v102+v82&v102)&v102
		}
	}
	return
}
func F_PredictorAdd5_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
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
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v40 int32
	_ = v40
	var v46 int32
	_ = v46
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	if l2 < int32(1) {
	} else {
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l3+int32(-4))))
		v13 = l0
		v14 = l1
		v15 = l2
		v16 = l3
		v17 = v12
		for {
			v20 = int32(4)
			v21 = v14 + v20
			v22 = *(*int32)(unsafe.Add(mBase, uint32(v21)))
			v24 = int32(1)
			v26 = int32(2139062143)
			v29 = int32(base.Ui32(v22^v17)>>(uint(v24)%32))&v26 + v22&v17
			v30 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
			v37 = int32(base.Ui32(v29^v30)>>(uint(v24)%32))&v26 + v29&v30
			v38 = int32(-16711936)
			v40 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
			v46 = int32(16711935)
			v53 = (v37&v38+v40&v38)&v38 | (v37&v46+v40&v46)&v46
			*(*int32)(unsafe.Add(mBase, uint32(v16))) = v53
			v60 = v15 + int32(-1)
			if v60 != 0 {
				v13 = v13 + v20
				v14 = v21
				v15 = v60
				v16 = v16 + v20
				v17 = v53
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorAdd6_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
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
	var v20 int32
	_ = v20
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v36 int32
	_ = v36
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v52 int32
	_ = v52
	if l2 < int32(1) {
	} else {
		v9 = int32(-4)
		v13 = *(*int32)(unsafe.Add(mBase, uint32(l3+v9)))
		v14 = l0
		v15 = l1 + v9
		v16 = l2
		v17 = l3
		v18 = v13
		for {
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
			v27 = int32(base.Ui32(v20^v18)>>(uint(int32(1))%32))&int32(2139062143) + v20&v18
			v28 = int32(-16711936)
			v30 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
			v36 = int32(16711935)
			v43 = (v27&v28+v30&v28)&v28 | (v27&v36+v30&v36)&v36
			*(*int32)(unsafe.Add(mBase, uint32(v17))) = v43
			v45 = int32(4)
			v52 = v16 + int32(-1)
			if v52 != 0 {
				v14 = v14 + v45
				v15 = v15 + v45
				v16 = v52
				v17 = v17 + v45
				v18 = v43
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorAdd7_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v34 int32
	_ = v34
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v50 int32
	_ = v50
	if l2 < int32(1) {
	} else {
		v11 = *(*int32)(unsafe.Add(mBase, uint32(l3+int32(-4))))
		v12 = l0
		v13 = l1
		v14 = l2
		v15 = l3
		v16 = v11
		for {
			v18 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
			v25 = int32(base.Ui32(v18^v16)>>(uint(int32(1))%32))&int32(2139062143) + v18&v16
			v26 = int32(-16711936)
			v28 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v34 = int32(16711935)
			v41 = (v25&v26+v28&v26)&v26 | (v25&v34+v28&v34)&v34
			*(*int32)(unsafe.Add(mBase, uint32(v15))) = v41
			v43 = int32(4)
			v50 = v14 + int32(-1)
			if v50 != 0 {
				v12 = v12 + v43
				v13 = v13 + v43
				v14 = v50
				v15 = v15 + v43
				v16 = v41
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub10_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v38 int32
	_ = v38
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	var v70 int32
	_ = v70
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v12)+4))
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v22 = int32(1)
			v24 = int32(2139062143)
			v27 = int32(base.Ui32(v19^v20)>>(uint(v22)%32))&v24 + v19&v20
			v30 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v31 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v38 = int32(base.Ui32(v30^v31)>>(uint(v22)%32))&v24 + v30&v31
			v45 = int32(base.Ui32(v27^v38)>>(uint(v22)%32))&v24 + v27&v38
			v46 = int32(4)
			v47 = v11 + v46
			v48 = *(*int32)(unsafe.Add(mBase, uint32(v47)))
			v51 = int32(-16711936)
			v58 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v48|int32(16711680)-v45&v51)&v51 | (v48|int32(_a_F_PredictorSub10_C_0)-v45&v58)&v58
			v70 = v13 + int32(-1)
			if v70 != 0 {
				v11 = v47
				v12 = v12 + v46
				v13 = v70
				v14 = v14 + v46
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub11_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v125 int32
	_ = v125
	var v137 int32
	_ = v137
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v22 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v23 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v24 = int32(255)
			v28 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v30 = v28 & v24
			v31 = v23&v24 - v30
			v32 = int32(31)
			v33 = v31 >> (uint(v32) % 32)
			v36 = int32(24)
			v39 = int32(base.Ui32(v28) >> (uint(v36) % 32))
			v40 = int32(base.Ui32(v23)>>(uint(v36)%32)) - v39
			v42 = v40 >> (uint(v32) % 32)
			v46 = int32(8)
			v53 = int32(base.Ui32(v28)>>(uint(v46)%32)) & v24
			v54 = int32(base.Ui32(v23)>>(uint(v46)%32))&v24 - v53
			v56 = v54 >> (uint(v32) % 32)
			v62 = int32(base.Ui32(v22)>>(uint(v36)%32)) - v39
			v64 = v62 >> (uint(v32) % 32)
			v69 = v22&v24 - v30
			v71 = v69 >> (uint(v32) % 32)
			v79 = int32(base.Ui32(v22)>>(uint(v46)%32))&v24 - v53
			v81 = v79 >> (uint(v32) % 32)
			v85 = int32(16)
			v92 = int32(base.Ui32(v28)>>(uint(v85)%32)) & v24
			v93 = int32(base.Ui32(v22)>>(uint(v85)%32))&v24 - v92
			v95 = v93 >> (uint(v32) % 32)
			v104 = int32(base.Ui32(v23)>>(uint(v85)%32))&v24 - v92
			v106 = v104 >> (uint(v32) % 32)
			if v31^v33-v33+(v40^v42-v42)+(v54^v56-v56)-(v62^v64-v64+(v69^v71-v71)+(v79^v81-v81)+(v93^v95-v95))+(v104^v106-v106) < int32(1) {
				v112 = v22
			} else {
				v112 = v23
			}
			v113 = int32(4)
			v114 = v11 + v113
			v115 = *(*int32)(unsafe.Add(mBase, uint32(v114)))
			v118 = int32(-16711936)
			v125 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v115|int32(16711680)-v112&v118)&v118 | (v115|int32(_a_F_PredictorSub11_C_0)-v112&v125)&v125
			v137 = v13 + int32(-1)
			if v137 != 0 {
				v11 = v114
				v12 = v12 + v113
				v13 = v137
				v14 = v14 + v113
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub12_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v94 int32
	_ = v94
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v117 int32
	_ = v117
	var v129 int32
	_ = v129
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v20 = int32(24)
			v22 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v28 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v31 = int32(base.Ui32(v19)>>(uint(v20)%32)) + int32(base.Ui32(v22)>>(uint(v20)%32)) - int32(base.Ui32(v28)>>(uint(v20)%32))
			if base.Ui32(v31) < base.Ui32(int32(256)) {
				v38 = v31
			} else {
				v38 = int32(base.Ui32(v31^int32(-1)) >> (uint(v20) % 32))
			}
			v39 = int32(24)
			v41 = int32(255)
			v48 = v19&v41 + v22&v41 - v28&v41
			if base.Ui32(v48) < base.Ui32(int32(256)) {
				v55 = v48
			} else {
				v55 = int32(base.Ui32(v48^int32(-1)) >> (uint(v39) % 32))
			}
			v57 = int32(16)
			v59 = int32(255)
			v70 = int32(base.Ui32(v19)>>(uint(v57)%32))&v59 + int32(base.Ui32(v22)>>(uint(v57)%32))&v59 - int32(base.Ui32(v28)>>(uint(v57)%32))&v59
			if base.Ui32(v70) < base.Ui32(int32(256)) {
				v77 = v70
			} else {
				v77 = int32(base.Ui32(v70^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v81 = int32(8)
			v83 = int32(255)
			v94 = int32(base.Ui32(v19)>>(uint(v81)%32))&v83 + int32(base.Ui32(v22)>>(uint(v81)%32))&v83 - int32(base.Ui32(v28)>>(uint(v81)%32))&v83
			if base.Ui32(v94) < base.Ui32(int32(256)) {
				v101 = v94
			} else {
				v101 = int32(base.Ui32(v94^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v104 = v38<<(uint(v39)%32) + v55 + v77<<(uint(int32(16))%32) + v101<<(uint(int32(8))%32)
			v105 = int32(4)
			v106 = v11 + v105
			v107 = *(*int32)(unsafe.Add(mBase, uint32(v106)))
			v110 = int32(-16711936)
			v117 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v107|int32(16711680)-v104&v110)&v110 | (v107|int32(_a_F_PredictorSub12_C_0)-v104&v117)&v117
			v129 = v13 + int32(-1)
			if v129 != 0 {
				v11 = v106
				v12 = v12 + v105
				v13 = v129
				v14 = v14 + v105
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub13_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v84 int32
	_ = v84
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v98 int32
	_ = v98
	var v100 int32
	_ = v100
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
	var v123 int32
	_ = v123
	var v135 int32
	_ = v135
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v18 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v26 = int32(base.Ui32(v18^v19)>>(uint(int32(1))%32))&int32(2139062143) + v18&v19
			v27 = int32(24)
			v28 = int32(base.Ui32(v26) >> (uint(v27) % 32))
			v31 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v36 = base.I32_div_s(v28-int32(base.Ui32(v31)>>(uint(v27)%32)), int32(2))
			v38 = v28 + base.I32_extend16_s(v36)
			if base.Ui32(v38) < base.Ui32(int32(256)) {
				v45 = v38
			} else {
				v45 = int32(base.Ui32(v38^int32(-1)) >> (uint(v27) % 32))
			}
			v46 = int32(24)
			v48 = int32(255)
			v49 = v26 & v48
			v54 = base.I32_div_s(v49-v31&v48, int32(2))
			v56 = v49 + base.I32_extend16_s(v54)
			if base.Ui32(v56) < base.Ui32(int32(256)) {
				v63 = v56
			} else {
				v63 = int32(base.Ui32(v56^int32(-1)) >> (uint(v46) % 32))
			}
			v65 = int32(16)
			v67 = int32(255)
			v68 = int32(base.Ui32(v26)>>(uint(v65)%32)) & v67
			v75 = base.I32_div_s(v68-int32(base.Ui32(v31)>>(uint(v65)%32))&v67, int32(2))
			v77 = v68 + base.I32_extend16_s(v75)
			if base.Ui32(v77) < base.Ui32(int32(256)) {
				v84 = v77
			} else {
				v84 = int32(base.Ui32(v77^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v88 = int32(8)
			v90 = int32(255)
			v91 = int32(base.Ui32(v26)>>(uint(v88)%32)) & v90
			v98 = base.I32_div_s(v91-int32(base.Ui32(v31)>>(uint(v88)%32))&v90, int32(2))
			v100 = v91 + base.I32_extend16_s(v98)
			if base.Ui32(v100) < base.Ui32(int32(256)) {
				v107 = v100
			} else {
				v107 = int32(base.Ui32(v100^int32(-1)) >> (uint(int32(24)) % 32))
			}
			v110 = v45<<(uint(v46)%32) + v63 + v84<<(uint(int32(16))%32) + v107<<(uint(int32(8))%32)
			v111 = int32(4)
			v112 = v11 + v111
			v113 = *(*int32)(unsafe.Add(mBase, uint32(v112)))
			v116 = int32(-16711936)
			v123 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v113|int32(16711680)-v110&v116)&v116 | (v113|int32(_a_F_PredictorSub13_C_0)-v110&v123)&v123
			v135 = v13 + int32(-1)
			if v135 != 0 {
				v11 = v112
				v12 = v12 + v111
				v13 = v135
				v14 = v14 + v111
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub2_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v30 int32
	_ = v30
	var v42 int32
	_ = v42
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v17 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v18 = int32(4)
			v19 = v11 + v18
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
			v23 = int32(-16711936)
			v30 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v20|int32(16711680)-v17&v23)&v23 | (v20|int32(_a_F_PredictorSub2_C_0)-v17&v30)&v30
			v42 = v13 + int32(-1)
			if v42 != 0 {
				v11 = v19
				v12 = v12 + v18
				v13 = v42
				v14 = v14 + v18
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub3_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v30 int32
	_ = v30
	var v42 int32
	_ = v42
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v17 = *(*int32)(unsafe.Add(mBase, uint32(v12)+4))
			v18 = int32(4)
			v19 = v11 + v18
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
			v23 = int32(-16711936)
			v30 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v20|int32(16711680)-v17&v23)&v23 | (v20|int32(_a_F_PredictorSub3_C_0)-v17&v30)&v30
			v42 = v13 + int32(-1)
			if v42 != 0 {
				v11 = v19
				v12 = v12 + v18
				v13 = v42
				v14 = v14 + v18
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub4_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
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
	var v25 int32
	_ = v25
	var v32 int32
	_ = v32
	var v44 int32
	_ = v44
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v20 = int32(4)
			v21 = v11 + v20
			v22 = *(*int32)(unsafe.Add(mBase, uint32(v21)))
			v25 = int32(-16711936)
			v32 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v22|int32(16711680)-v19&v25)&v25 | (v22|int32(_a_F_PredictorSub4_C_0)-v19&v32)&v32
			v44 = v13 + int32(-1)
			if v44 != 0 {
				v11 = v21
				v12 = v12 + v20
				v13 = v44
				v14 = v14 + v20
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub5_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v47 int32
	_ = v47
	var v59 int32
	_ = v59
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v18 = *(*int32)(unsafe.Add(mBase, uint32(v12)+4))
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v21 = int32(1)
			v23 = int32(2139062143)
			v26 = int32(base.Ui32(v18^v19)>>(uint(v21)%32))&v23 + v18&v19
			v27 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v34 = int32(base.Ui32(v26^v27)>>(uint(v21)%32))&v23 + v26&v27
			v35 = int32(4)
			v36 = v11 + v35
			v37 = *(*int32)(unsafe.Add(mBase, uint32(v36)))
			v40 = int32(-16711936)
			v47 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v37|int32(16711680)-v34&v40)&v40 | (v37|int32(_a_F_PredictorSub5_C_0)-v34&v47)&v47
			v59 = v13 + int32(-1)
			if v59 != 0 {
				v11 = v36
				v12 = v12 + v35
				v13 = v59
				v14 = v14 + v35
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub6_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v40 int32
	_ = v40
	var v52 int32
	_ = v52
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v27 = int32(base.Ui32(v19^v20)>>(uint(int32(1))%32))&int32(2139062143) + v19&v20
			v28 = int32(4)
			v29 = v11 + v28
			v30 = *(*int32)(unsafe.Add(mBase, uint32(v29)))
			v33 = int32(-16711936)
			v40 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v30|int32(16711680)-v27&v33)&v33 | (v30|int32(_a_F_PredictorSub6_C_0)-v27&v40)&v40
			v52 = v13 + int32(-1)
			if v52 != 0 {
				v11 = v29
				v12 = v12 + v28
				v13 = v52
				v14 = v14 + v28
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub7_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v38 int32
	_ = v38
	var v50 int32
	_ = v50
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v17 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v18 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
			v25 = int32(base.Ui32(v17^v18)>>(uint(int32(1))%32))&int32(2139062143) + v17&v18
			v26 = int32(4)
			v27 = v11 + v26
			v28 = *(*int32)(unsafe.Add(mBase, uint32(v27)))
			v31 = int32(-16711936)
			v38 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v28|int32(16711680)-v25&v31)&v31 | (v28|int32(_a_F_PredictorSub7_C_0)-v25&v38)&v38
			v50 = v13 + int32(-1)
			if v50 != 0 {
				v11 = v27
				v12 = v12 + v26
				v13 = v50
				v14 = v14 + v26
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub8_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v41 int32
	_ = v41
	var v53 int32
	_ = v53
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v18 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v21 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(-4))))
			v28 = int32(base.Ui32(v18^v21)>>(uint(int32(1))%32))&int32(2139062143) + v18&v21
			v29 = int32(4)
			v30 = v11 + v29
			v31 = *(*int32)(unsafe.Add(mBase, uint32(v30)))
			v34 = int32(-16711936)
			v41 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v31|int32(16711680)-v28&v34)&v34 | (v31|int32(_a_F_PredictorSub8_C_0)-v28&v41)&v41
			v53 = v13 + int32(-1)
			if v53 != 0 {
				v11 = v30
				v12 = v12 + v29
				v13 = v53
				v14 = v14 + v29
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PredictorSub9_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v39 int32
	_ = v39
	var v51 int32
	_ = v51
	if l2 < int32(1) {
	} else {
		v11 = l0 + int32(-4)
		v12 = l1
		v13 = l2
		v14 = l3
		for {
			v18 = *(*int32)(unsafe.Add(mBase, uint32(v12)+4))
			v19 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			v26 = int32(base.Ui32(v18^v19)>>(uint(int32(1))%32))&int32(2139062143) + v18&v19
			v27 = int32(4)
			v28 = v11 + v27
			v29 = *(*int32)(unsafe.Add(mBase, uint32(v28)))
			v32 = int32(-16711936)
			v39 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(v14))) = (v29|int32(16711680)-v26&v32)&v32 | (v29|int32(_a_F_PredictorSub9_C_0)-v26&v39)&v39
			v51 = v13 + int32(-1)
			if v51 != 0 {
				v11 = v28
				v12 = v12 + v27
				v13 = v51
				v14 = v14 + v27
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PrepareMapToPalette(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v29 int32
	_ = v29
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v56 int32
	_ = v56
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v81 int32
	_ = v81
	v14 = F_memcpy(m, l2, l0, l1<<(uint(int32(2))%32))
	mBase = m.M
	v16 = m.G2
	F___qsort_r(m, v14, l1, int32(4), int32(344), v16+int32(328))
	mBase = m.M
	if l1 == int32(0) {
	} else {
		v29 = int32(0)
		for {
			v36 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
			v40 = *(*int32)(unsafe.Add(mBase, uint32(l0+v29<<(uint(int32(2))%32))))
			if v36 == v40 {
				v67 = int32(0)
			} else {
				v50 = int32(0)
				v51 = l1
				for {
					v56 = (v51 + v50) >> (uint(int32(1)) % 32)
					v60 = *(*int32)(unsafe.Add(mBase, uint32(v14+v56<<(uint(int32(2))%32))))
					v61 = base.B2i32(base.Ui32(v60) < base.Ui32(v40))
					if base.Ui32(v60) < base.Ui32(v40) {
						v62 = v51
					} else {
						v62 = v56
					}
					if base.Ui32(v60) < base.Ui32(v40) {
						v63 = v56
					} else {
						v63 = v50
					}
					if v60 != v40 {
						v50 = v63
						v51 = v62
						continue
					} else {
						break
					}
					break
				}
				v67 = v56
			}
			*(*int32)(unsafe.Add(mBase, uint32(l3+v67<<(uint(int32(2))%32)))) = v29
			v81 = v29 + int32(1)
			if v81 != l1 {
				v29 = v81
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_PushInterval(m *base.Module, l0 int32, l1 int64, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v26 int32
	_ = v26
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v47 int64
	_ = v47
	var v48 int64
	_ = v48
	var v49 int64
	_ = v49
	var v52 int32
	_ = v52
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v64 int32
	_ = v64
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v84 int32
	_ = v84
	var v90 int32
	_ = v90
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v99 int64
	_ = v99
	var v100 int64
	_ = v100
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v128 int64
	_ = v128
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v144 int32
	_ = v144
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v150 int32
	_ = v150
	var v157 int32
	_ = v157
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v172 int32
	_ = v172
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	if int32(9) < l3 {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	return
L2:
	;
	v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v67 == int32(0) {
		goto L1
	} else {
		goto L10
	}
L3:
	;
	if l3 < int32(1) {
		goto L1
	} else {
		goto L4
	}
L4:
	;
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_PushInterval[0])))
	v36 = l2 << (uint(int32(1)) % 32)
	v37 = l0 + int32(16)
	v38 = v26 + l2<<(uint(int32(3))%32)
	v39 = int32(0)
	goto L5
L5:
	;
	v47 = *(*int64)(unsafe.Add(mBase, uint32(v38)))
	v48 = *(*int64)(unsafe.Add(mBase, uint32(v37)))
	v49 = v48 + l1
	if v47 <= v49 {
		goto L7
	} else {
		goto L8
	}
L7:
	;
	v57 = int32(8)
	v64 = v39 + int32(1)
	if l2+v64 < l3+l2 {
		v36 = v36 + int32(2)
		v37 = v37 + v57
		v38 = v38 + v57
		v39 = v64
		goto L5
	} else {
		goto L9
	}
L8:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v38))) = v49
	v52 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_PushInterval[1])))
	v55 = v39 + int32(1)
	*(*uint16)(unsafe.Add(mBase, uint32(v52+v36))) = uint16(v55)
	goto L7
L9:
	;
	goto L1
L10:
	;
	v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v84 = v75
	v90 = int32(0)
	goto L11
L11:
	;
	v95 = v70 + v90<<(uint(int32(4))%32)
	v96 = *(*int32)(unsafe.Add(mBase, uint32(v95)+8))
	if l3 <= v96 {
		goto L1
	} else {
		goto L13
	}
L12:
	;
	goto L1
L13:
	;
	v98 = v96 + l2
	v99 = *(*int64)(unsafe.Add(mBase, uint32(v95)))
	v100 = v99 + l1
	v101 = *(*int32)(unsafe.Add(mBase, uint32(v95)+12))
	if v101 < l3 {
		goto L14
	} else {
		goto L15
	}
L14:
	;
	v103 = v101
	goto L16
L15:
	;
	v103 = l3
	goto L16
L16:
	;
	v104 = v103 + l2
	if v84 == int32(0) {
		v172 = v98
		goto L18
	} else {
		goto L19
	}
L17:
	;
	F_InsertInterval(m, l0, v191, v100, l2, v192, v104)
	mBase = m.M
	v202 = v90 + int32(1)
	v203 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if base.Ui32(v202) < base.Ui32(v203) {
		v84 = v191
		v90 = v202
		goto L11
	} else {
		goto L44
	}
L18:
	;
	v191 = int32(0)
	v192 = v172
	goto L17
L19:
	;
	v112 = v98
	v113 = v84
	goto L20
L20:
	;
	v123 = *(*int32)(unsafe.Add(mBase, uint32(v113)+8))
	if v104 <= v123 {
		v161 = v113
		goto L23
	} else {
		goto L24
	}
L21:
	;
	v172 = v164
	goto L18
L22:
	;
	if v125 != 0 {
		v112 = v164
		v113 = v125
		goto L20
	} else {
		goto L43
	}
L23:
	;
	v191 = v161
	v192 = v112
	goto L17
L24:
	;
	v125 = *(*int32)(unsafe.Add(mBase, uint32(v113)+24))
	v126 = *(*int32)(unsafe.Add(mBase, uint32(v113)+12))
	if v126 <= v112 {
		v164 = v112
		goto L22
	} else {
		goto L25
	}
L25:
	;
	v128 = *(*int64)(unsafe.Add(mBase, uint32(v113)))
	if v100 < v128 {
		goto L26
	} else {
		goto L27
	}
L26:
	;
	if v123 < v112 {
		goto L29
	} else {
		goto L30
	}
L27:
	;
	F_InsertInterval(m, l0, v113, v100, l2, v112, v123)
	mBase = m.M
	if v104 <= v126 {
		v191 = v113
		v192 = v126
		goto L17
	} else {
		goto L28
	}
L28:
	;
	v164 = v126
	goto L22
L29:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v113)+12)) = v112
	if v126 <= v104 {
		v164 = v112
		goto L22
	} else {
		goto L42
	}
L30:
	;
	if v104 < v126 {
		goto L31
	} else {
		goto L32
	}
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v113)+8)) = v104
	v161 = v113
	goto L23
L32:
	;
	v134 = *(*int32)(unsafe.Add(mBase, uint32(v113)+20))
	if v134 != 0 {
		goto L33
	} else {
		goto L34
	}
L33:
	;
	v137 = v134 + int32(24)
	goto L35
L34:
	;
	v137 = l0
	goto L35
L35:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v137))) = v125
	if v125 == int32(0) {
		goto L36
	} else {
		goto L37
	}
L36:
	;
	if base.Ui32(v113) < base.Ui32(l0+int32(_a_F_PushInterval_0)) {
		goto L39
	} else {
		goto L40
	}
L37:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v125)+20)) = v134
	goto L36
L38:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v113)+24)) = v148
	v150 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v150 + int32(-1)
	v164 = v112
	goto L22
L39:
	;
	v146 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_PushInterval[2])))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_PushInterval[2]))) = v113
	v148 = v146
	goto L38
L40:
	;
	if base.Ui32(l0+int32(_a_F_PushInterval_1)) < base.Ui32(v113) {
		goto L39
	} else {
		goto L41
	}
L41:
	;
	v144 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_PushInterval[3])))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_PushInterval[3]))) = v113
	v148 = v144
	goto L38
L42:
	;
	v157 = *(*int32)(unsafe.Add(mBase, uint32(v113)+16))
	F_InsertInterval(m, l0, v113, v128, v157, v104, v126)
	mBase = m.M
	v159 = *(*int32)(unsafe.Add(mBase, uint32(v113)+24))
	v161 = v159
	goto L23
L43:
	;
	goto L21
L44:
	;
	goto L12
}
func F_PutCoeffs(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
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
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v100 int32
	_ = v100
	var v104 int32
	_ = v104
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v132 int32
	_ = v132
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v155 int32
	_ = v155
	var v159 int32
	_ = v159
	var v163 int32
	_ = v163
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v187 int32
	_ = v187
	var v188 int32
	_ = v188
	var v190 int32
	_ = v190
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v206 int32
	_ = v206
	var v210 int32
	_ = v210
	var v212 int32
	_ = v212
	var v216 int32
	_ = v216
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v227 int32
	_ = v227
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v235 int32
	_ = v235
	var v238 int32
	_ = v238
	var v240 int32
	_ = v240
	var v244 int32
	_ = v244
	var v248 int32
	_ = v248
	var v252 int32
	_ = v252
	var v254 int32
	_ = v254
	var v258 int32
	_ = v258
	var v261 int32
	_ = v261
	var v262 int32
	_ = v262
	var v273 int32
	_ = v273
	var v275 int32
	_ = v275
	var v278 int32
	_ = v278
	var v281 int32
	_ = v281
	var v283 int32
	_ = v283
	var v287 int32
	_ = v287
	var v291 int32
	_ = v291
	var v295 int32
	_ = v295
	var v297 int32
	_ = v297
	var v301 int32
	_ = v301
	var v304 int32
	_ = v304
	var v305 int32
	_ = v305
	var v313 int32
	_ = v313
	var v314 int32
	_ = v314
	var v316 int32
	_ = v316
	var v319 int32
	_ = v319
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v328 int32
	_ = v328
	var v332 int32
	_ = v332
	var v336 int32
	_ = v336
	var v338 int32
	_ = v338
	var v342 int32
	_ = v342
	var v345 int32
	_ = v345
	var v346 int32
	_ = v346
	var v354 int32
	_ = v354
	var v355 int32
	_ = v355
	var v357 int32
	_ = v357
	var v360 int32
	_ = v360
	var v363 int32
	_ = v363
	var v365 int32
	_ = v365
	var v369 int32
	_ = v369
	var v373 int32
	_ = v373
	var v377 int32
	_ = v377
	var v379 int32
	_ = v379
	var v383 int32
	_ = v383
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v398 int32
	_ = v398
	var v401 int32
	_ = v401
	var v404 int32
	_ = v404
	var v406 int32
	_ = v406
	var v410 int32
	_ = v410
	var v414 int32
	_ = v414
	var v418 int32
	_ = v418
	var v420 int32
	_ = v420
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v428 int32
	_ = v428
	var v435 int32
	_ = v435
	var v439 int32
	_ = v439
	var v442 int32
	_ = v442
	var v445 int32
	_ = v445
	var v447 int32
	_ = v447
	var v451 int32
	_ = v451
	var v455 int32
	_ = v455
	var v459 int32
	_ = v459
	var v461 int32
	_ = v461
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v482 int32
	_ = v482
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v490 int32
	_ = v490
	var v494 int32
	_ = v494
	var v498 int32
	_ = v498
	var v502 int32
	_ = v502
	var v504 int32
	_ = v504
	var v508 int32
	_ = v508
	var v511 int32
	_ = v511
	var v512 int32
	_ = v512
	var v522 int32
	_ = v522
	var v524 int32
	_ = v524
	var v527 int32
	_ = v527
	var v540 int32
	_ = v540
	var v544 int32
	_ = v544
	var v546 int32
	_ = v546
	var v550 int32
	_ = v550
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v562 int32
	_ = v562
	var v564 int32
	_ = v564
	var v567 int32
	_ = v567
	var v580 int32
	_ = v580
	var v584 int32
	_ = v584
	var v586 int32
	_ = v586
	var v590 int32
	_ = v590
	var v593 int32
	_ = v593
	var v594 int32
	_ = v594
	var v601 int32
	_ = v601
	var v609 int32
	_ = v609
	var v611 int32
	_ = v611
	var v614 int32
	_ = v614
	var v627 int32
	_ = v627
	var v631 int32
	_ = v631
	var v633 int32
	_ = v633
	var v637 int32
	_ = v637
	var v640 int32
	_ = v640
	var v641 int32
	_ = v641
	var v649 int32
	_ = v649
	var v651 int32
	_ = v651
	var v657 int32
	_ = v657
	var v659 int32
	_ = v659
	var v662 int32
	_ = v662
	var v667 int32
	_ = v667
	var v671 int32
	_ = v671
	var v673 int32
	_ = v673
	var v677 int32
	_ = v677
	var v680 int32
	_ = v680
	var v681 int32
	_ = v681
	var v688 int32
	_ = v688
	var v694 int32
	_ = v694
	var v696 int32
	_ = v696
	var v702 int32
	_ = v702
	var v704 int32
	_ = v704
	var v707 int32
	_ = v707
	var v712 int32
	_ = v712
	var v716 int32
	_ = v716
	var v718 int32
	_ = v718
	var v722 int32
	_ = v722
	var v725 int32
	_ = v725
	var v726 int32
	_ = v726
	var v733 int32
	_ = v733
	var v738 int32
	_ = v738
	var v741 int32
	_ = v741
	var v754 int32
	_ = v754
	var v758 int32
	_ = v758
	var v760 int32
	_ = v760
	var v764 int32
	_ = v764
	var v767 int32
	_ = v767
	var v768 int32
	_ = v768
	var v775 int32
	_ = v775
	var v782 int32
	_ = v782
	var v788 int32
	_ = v788
	var v790 int32
	_ = v790
	var v793 int32
	_ = v793
	var v798 int32
	_ = v798
	var v802 int32
	_ = v802
	var v804 int32
	_ = v804
	var v808 int32
	_ = v808
	var v811 int32
	_ = v811
	var v812 int32
	_ = v812
	var v819 int32
	_ = v819
	var v824 int32
	_ = v824
	var v825 int32
	_ = v825
	var v826 int32
	_ = v826
	var v829 int32
	_ = v829
	var v833 int32
	_ = v833
	var v839 int32
	_ = v839
	var v841 int32
	_ = v841
	var v843 int32
	_ = v843
	var v846 int32
	_ = v846
	var v849 int32
	_ = v849
	var v851 int32
	_ = v851
	var v855 int32
	_ = v855
	var v859 int32
	_ = v859
	var v863 int32
	_ = v863
	var v865 int32
	_ = v865
	var v869 int32
	_ = v869
	var v872 int32
	_ = v872
	var v873 int32
	_ = v873
	var v880 int32
	_ = v880
	var v905 int32
	_ = v905
	var v907 int32
	_ = v907
	var v910 int32
	_ = v910
	var v912 int32
	_ = v912
	var v915 int32
	_ = v915
	var v917 int32
	_ = v917
	var v921 int32
	_ = v921
	var v925 int32
	_ = v925
	var v926 int32
	_ = v926
	var v929 int32
	_ = v929
	var v933 int32
	_ = v933
	var v937 int32
	_ = v937
	var v943 int32
	_ = v943
	var v946 int32
	_ = v946
	var v947 int32
	_ = v947
	var v949 int32
	_ = v949
	var v955 int32
	_ = v955
	var v956 int32
	_ = v956
	var v958 int32
	_ = v958
	var v961 int32
	_ = v961
	var v964 int32
	_ = v964
	var v966 int32
	_ = v966
	var v970 int32
	_ = v970
	var v974 int32
	_ = v974
	var v978 int32
	_ = v978
	var v980 int32
	_ = v980
	var v984 int32
	_ = v984
	var v987 int32
	_ = v987
	var v988 int32
	_ = v988
	var v998 int32
	_ = v998
	var v1027 int32
	_ = v1027
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	v15 = int32(base.Ui32(v11^int32(-1)) >> (uint(int32(31)) % 32))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
	v23 = v16 + v17*int32(33) + l1*int32(11)
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v23))))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v29 = v26 * v24 >> (uint(int32(8)) % 32)
	if v15 == int32(0) {
		v38 = v29
		goto L3
	} else {
		goto L4
	}
L1:
	;
	if int32(15) < v17 {
		goto L10
	} else {
		goto L11
	}
L2:
	;
	if v15 != 0 {
		goto L1
	} else {
		goto L8
	}
L3:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v38
	if int32(126) < v38 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v34 = v29 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v32 + v34
	v38 = v26 - v34
	goto L3
L5:
	;
	goto L2
L6:
	;
	v42 = m.G1
	v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v42+int32(_a_F_PutCoeffs_0)+v38))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v46
	v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v42+int32(_a_F_PutCoeffs_1)+v38))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v48 << (uint(v52) % 32)
	v55 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v56 = v52 + v55
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v56
	if v56 < int32(1) {
		goto L5
	} else {
		goto L7
	}
L7:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L5
L8:
	;
	return int32(0)
L9:
	;
	return v1027
L10:
	;
	v1027 = int32(1)
	goto L9
L11:
	;
	v68 = v23
	v70 = v17
	goto L12
L12:
	;
	v77 = int32(1)
	v78 = v70 + v77
	v79 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v83 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v79+v70<<(uint(v77)%32)))))
	v84 = int32(0)
	v85 = base.B2i32(v83 != v84)
	v86 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+1)))
	v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v91 = v88 * v86 >> (uint(int32(8)) % 32)
	if v85 == v84 {
		v100 = v91
		goto L17
	} else {
		goto L18
	}
L13:
	;
	goto L10
L14:
	;
	if v78 != int32(16) {
		v68 = v998
		v70 = v78
		goto L12
	} else {
		goto L159
	}
L15:
	;
	v132 = int32(15)
	v134 = int32(1)
	v138 = base.I32_extend16_s(v83) >> (uint(v132) % 32)
	v142 = (v83 ^ v138 - v138) & int32(_a_F_PutCoeffs_2)
	v144 = base.B2i32(base.Ui32(v134) < base.Ui32(v142))
	v145 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+2)))
	v147 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v150 = v147 * v145 >> (uint(int32(8)) % 32)
	if v144 == int32(0) {
		v159 = v150
		goto L25
	} else {
		goto L26
	}
L16:
	;
	if v83 != v84 {
		goto L15
	} else {
		goto L22
	}
L17:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v100
	if int32(126) < v100 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v96 = v91 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v94 + v96
	v100 = v88 - v96
	goto L17
L19:
	;
	goto L16
L20:
	;
	v104 = m.G1
	v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104+int32(_a_F_PutCoeffs_0)+v100))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v108
	v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104+int32(_a_F_PutCoeffs_1)+v100))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v110 << (uint(v114) % 32)
	v117 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v118 = v114 + v117
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v118
	if v118 < int32(1) {
		goto L19
	} else {
		goto L21
	}
L21:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L19
L22:
	;
	v125 = m.G81
	v126 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v125+v78))))
	v998 = v126 + v128*int32(33)
	goto L14
L23:
	;
	v907 = *(*int32)(unsafe.Add(mBase, uint32(l2)+16))
	v910 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v912 = v910 >> (uint(int32(1)) % 32)
	if int32(base.Ui32(v83)>>(uint(v132)%32)) == int32(0) {
		v921 = v912
		goto L146
	} else {
		goto L147
	}
L24:
	;
	if v144 == int32(0) {
		v905 = v134
		goto L23
	} else {
		goto L30
	}
L25:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v159
	if int32(126) < v159 {
		goto L27
	} else {
		goto L28
	}
L26:
	;
	v153 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v155 = v150 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v153 + v155
	v159 = v147 - v155
	goto L25
L27:
	;
	goto L24
L28:
	;
	v163 = m.G1
	v167 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v163+int32(_a_F_PutCoeffs_0)+v159))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v167
	v169 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v173 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v163+int32(_a_F_PutCoeffs_1)+v159))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v169 << (uint(v173) % 32)
	v176 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v177 = v173 + v176
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v177
	if v177 < int32(1) {
		goto L27
	} else {
		goto L29
	}
L29:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L27
L30:
	;
	v187 = base.B2i32(base.Ui32(int32(4)) < base.Ui32(v142))
	v188 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+3)))
	v190 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v193 = v190 * v188 >> (uint(int32(8)) % 32)
	if v187 == int32(0) {
		v202 = v193
		goto L33
	} else {
		goto L34
	}
L31:
	;
	v313 = base.B2i32(base.Ui32(int32(10)) < base.Ui32(v142))
	v314 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+6)))
	v316 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v319 = v316 * v314 >> (uint(int32(8)) % 32)
	if v313 == int32(0) {
		v328 = v319
		goto L55
	} else {
		goto L56
	}
L32:
	;
	if base.Ui32(int32(4)) < base.Ui32(v142) {
		goto L31
	} else {
		goto L38
	}
L33:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v202
	if int32(126) < v202 {
		goto L35
	} else {
		goto L36
	}
L34:
	;
	v196 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v198 = v193 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v196 + v198
	v202 = v190 - v198
	goto L33
L35:
	;
	goto L32
L36:
	;
	v206 = m.G1
	v210 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v206+int32(_a_F_PutCoeffs_0)+v202))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v210
	v212 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v216 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v206+int32(_a_F_PutCoeffs_1)+v202))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v212 << (uint(v216) % 32)
	v219 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v220 = v216 + v219
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v220
	if v220 < int32(1) {
		goto L35
	} else {
		goto L37
	}
L37:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L35
L38:
	;
	v227 = int32(2)
	v229 = base.B2i32(v142 != v227)
	v230 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+4)))
	v232 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v235 = v232 * v230 >> (uint(int32(8)) % 32)
	if v229 == int32(0) {
		v244 = v235
		goto L40
	} else {
		goto L41
	}
L39:
	;
	if v229 == int32(0) {
		v905 = v227
		goto L23
	} else {
		goto L45
	}
L40:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v244
	if int32(126) < v244 {
		goto L42
	} else {
		goto L43
	}
L41:
	;
	v238 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v240 = v235 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v238 + v240
	v244 = v232 - v240
	goto L40
L42:
	;
	goto L39
L43:
	;
	v248 = m.G1
	v252 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v248+int32(_a_F_PutCoeffs_0)+v244))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v252
	v254 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v258 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v248+int32(_a_F_PutCoeffs_1)+v244))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v254 << (uint(v258) % 32)
	v261 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v262 = v258 + v261
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v262
	if v262 < int32(1) {
		goto L42
	} else {
		goto L44
	}
L44:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L42
L45:
	;
	v273 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+5)))
	v275 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v278 = v275 * v273 >> (uint(int32(8)) % 32)
	if base.B2i32(v142 == int32(4)) == int32(0) {
		v287 = v278
		goto L47
	} else {
		goto L48
	}
L46:
	;
	v905 = v227
	goto L23
L47:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v287
	if int32(126) < v287 {
		goto L49
	} else {
		goto L50
	}
L48:
	;
	v281 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v283 = v278 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v281 + v283
	v287 = v275 - v283
	goto L47
L49:
	;
	goto L46
L50:
	;
	v291 = m.G1
	v295 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+int32(_a_F_PutCoeffs_0)+v287))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v295
	v297 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v301 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v291+int32(_a_F_PutCoeffs_1)+v287))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v297 << (uint(v301) % 32)
	v304 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v305 = v301 + v304
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v305
	if v305 < int32(1) {
		goto L49
	} else {
		goto L51
	}
L51:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L49
L52:
	;
	v905 = int32(2)
	goto L23
L53:
	;
	if base.Ui32(int32(18)) < base.Ui32(v142) {
		goto L88
	} else {
		goto L89
	}
L54:
	;
	if base.Ui32(int32(10)) < base.Ui32(v142) {
		goto L53
	} else {
		goto L60
	}
L55:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v328
	if int32(126) < v328 {
		goto L57
	} else {
		goto L58
	}
L56:
	;
	v322 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v324 = v319 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v322 + v324
	v328 = v316 - v324
	goto L55
L57:
	;
	goto L54
L58:
	;
	v332 = m.G1
	v336 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v332+int32(_a_F_PutCoeffs_0)+v328))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v336
	v338 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v342 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v332+int32(_a_F_PutCoeffs_1)+v328))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v338 << (uint(v342) % 32)
	v345 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v346 = v342 + v345
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v346
	if v346 < int32(1) {
		goto L57
	} else {
		goto L59
	}
L59:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L57
L60:
	;
	v354 = base.B2i32(base.Ui32(int32(6)) < base.Ui32(v142))
	v355 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+7)))
	v357 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v360 = v357 * v355 >> (uint(int32(8)) % 32)
	if v354 == int32(0) {
		v369 = v360
		goto L63
	} else {
		goto L64
	}
L61:
	;
	v435 = int32(8)
	v439 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v442 = v439 * int32(165) >> (uint(v435) % 32)
	if base.B2i32(base.Ui32(v435) < base.Ui32(v142)) == int32(0) {
		v451 = v442
		goto L76
	} else {
		goto L77
	}
L62:
	;
	if base.Ui32(int32(6)) < base.Ui32(v142) {
		goto L61
	} else {
		goto L68
	}
L63:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v369
	if int32(126) < v369 {
		goto L65
	} else {
		goto L66
	}
L64:
	;
	v363 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v365 = v360 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v363 + v365
	v369 = v357 - v365
	goto L63
L65:
	;
	goto L62
L66:
	;
	v373 = m.G1
	v377 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v373+int32(_a_F_PutCoeffs_0)+v369))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v377
	v379 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v383 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v373+int32(_a_F_PutCoeffs_1)+v369))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v379 << (uint(v383) % 32)
	v386 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v387 = v383 + v386
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v387
	if v387 < int32(1) {
		goto L65
	} else {
		goto L67
	}
L67:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L65
L68:
	;
	v398 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v401 = v398 * int32(159) >> (uint(int32(8)) % 32)
	if base.B2i32(v142 == int32(6)) == int32(0) {
		v410 = v401
		goto L70
	} else {
		goto L71
	}
L69:
	;
	goto L52
L70:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v410
	if int32(126) < v410 {
		goto L72
	} else {
		goto L73
	}
L71:
	;
	v404 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v406 = v401 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v404 + v406
	v410 = v398 - v406
	goto L70
L72:
	;
	goto L69
L73:
	;
	v414 = m.G1
	v418 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v414+int32(_a_F_PutCoeffs_0)+v410))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v418
	v420 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v424 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v414+int32(_a_F_PutCoeffs_1)+v410))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v420 << (uint(v424) % 32)
	v427 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v428 = v424 + v427
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v428
	if v428 < int32(1) {
		goto L72
	} else {
		goto L74
	}
L74:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L72
L75:
	;
	v482 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v485 = v482 * int32(145) >> (uint(int32(8)) % 32)
	if (v142^int32(-1))&int32(1) == int32(0) {
		v494 = v485
		goto L82
	} else {
		goto L83
	}
L76:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v451
	if int32(126) < v451 {
		goto L78
	} else {
		goto L79
	}
L77:
	;
	v445 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v447 = v442 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v445 + v447
	v451 = v439 - v447
	goto L76
L78:
	;
	goto L75
L79:
	;
	v455 = m.G1
	v459 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v455+int32(_a_F_PutCoeffs_0)+v451))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v459
	v461 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v465 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v455+int32(_a_F_PutCoeffs_1)+v451))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v461 << (uint(v465) % 32)
	v468 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v469 = v465 + v468
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v469
	if v469 < int32(1) {
		goto L78
	} else {
		goto L80
	}
L80:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L78
L81:
	;
	goto L52
L82:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v494
	if int32(126) < v494 {
		goto L84
	} else {
		goto L85
	}
L83:
	;
	v488 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v490 = v485 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v488 + v490
	v494 = v482 - v490
	goto L82
L84:
	;
	goto L81
L85:
	;
	v498 = m.G1
	v502 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v498+int32(_a_F_PutCoeffs_0)+v494))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v502
	v504 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v508 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v498+int32(_a_F_PutCoeffs_1)+v494))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v504 << (uint(v508) % 32)
	v511 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v512 = v508 + v511
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v512
	if v512 < int32(1) {
		goto L84
	} else {
		goto L86
	}
L86:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L84
L87:
	;
	v829 = v824
	v833 = v825
	goto L136
L88:
	;
	if base.Ui32(int32(34)) < base.Ui32(v142) {
		goto L102
	} else {
		goto L103
	}
L89:
	;
	v522 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+8)))
	v524 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v527 = v524 * v522 >> (uint(int32(8)) % 32)
	goto L91
L90:
	;
	v562 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+9)))
	v564 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v567 = v564 * v562 >> (uint(int32(8)) % 32)
	goto L97
L91:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v527
	if int32(126) < v527 {
		goto L93
	} else {
		goto L94
	}
L93:
	;
	goto L90
L94:
	;
	v540 = m.G1
	v544 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v540+int32(_a_F_PutCoeffs_0)+v527))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v544
	v546 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v550 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v540+int32(_a_F_PutCoeffs_1)+v527))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v546 << (uint(v550) % 32)
	v553 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v554 = v550 + v553
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v554
	if v554 < int32(1) {
		goto L93
	} else {
		goto L95
	}
L95:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L93
L96:
	;
	v601 = m.G1
	v824 = int32(4)
	v825 = v601 + int32(_a_F_PutCoeffs_3)
	v826 = int32(-11)
	goto L87
L97:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v567
	if int32(126) < v567 {
		goto L99
	} else {
		goto L100
	}
L99:
	;
	goto L96
L100:
	;
	v580 = m.G1
	v584 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v580+int32(_a_F_PutCoeffs_0)+v567))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v584
	v586 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v590 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v580+int32(_a_F_PutCoeffs_1)+v567))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v586 << (uint(v590) % 32)
	v593 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v594 = v590 + v593
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v594
	if v594 < int32(1) {
		goto L99
	} else {
		goto L101
	}
L101:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L99
L102:
	;
	v694 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+8)))
	v696 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	goto L118
L103:
	;
	v609 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+8)))
	v611 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v614 = v611 * v609 >> (uint(int32(8)) % 32)
	goto L105
L104:
	;
	v649 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+9)))
	v651 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	goto L112
L105:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v614
	if int32(126) < v614 {
		goto L107
	} else {
		goto L108
	}
L107:
	;
	goto L104
L108:
	;
	v627 = m.G1
	v631 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v627+int32(_a_F_PutCoeffs_0)+v614))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v631
	v633 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v637 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v627+int32(_a_F_PutCoeffs_1)+v614))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v633 << (uint(v637) % 32)
	v640 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v641 = v637 + v640
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v641
	if v641 < int32(1) {
		goto L107
	} else {
		goto L109
	}
L109:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L107
L110:
	;
	v688 = m.G1
	v824 = int32(8)
	v825 = v688 + int32(_a_F_PutCoeffs_4)
	v826 = int32(-19)
	goto L87
L111:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v662
	if int32(126) < v662 {
		goto L113
	} else {
		goto L114
	}
L112:
	;
	v657 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v659 = v651*v649>>(uint(int32(8))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v657 + v659
	v662 = v651 - v659
	goto L111
L113:
	;
	goto L110
L114:
	;
	v667 = m.G1
	v671 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v667+int32(_a_F_PutCoeffs_0)+v662))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v671
	v673 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v677 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v667+int32(_a_F_PutCoeffs_1)+v662))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v673 << (uint(v677) % 32)
	v680 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v681 = v677 + v680
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v681
	if v681 < int32(1) {
		goto L113
	} else {
		goto L115
	}
L115:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L113
L116:
	;
	v733 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v68)+10)))
	if base.Ui32(int32(66)) < base.Ui32(v142) {
		goto L122
	} else {
		goto L123
	}
L117:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v707
	if int32(126) < v707 {
		goto L119
	} else {
		goto L120
	}
L118:
	;
	v702 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v704 = v696*v694>>(uint(int32(8))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v702 + v704
	v707 = v696 - v704
	goto L117
L119:
	;
	goto L116
L120:
	;
	v712 = m.G1
	v716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v712+int32(_a_F_PutCoeffs_0)+v707))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v716
	v718 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v712+int32(_a_F_PutCoeffs_1)+v707))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v718 << (uint(v722) % 32)
	v725 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v726 = v722 + v725
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v726
	if v726 < int32(1) {
		goto L119
	} else {
		goto L121
	}
L121:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L119
L122:
	;
	v782 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	goto L132
L123:
	;
	v738 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v741 = v738 * v733 >> (uint(int32(8)) % 32)
	goto L125
L124:
	;
	v775 = m.G1
	v824 = int32(16)
	v825 = v775 + int32(_a_F_PutCoeffs_5)
	v826 = int32(-35)
	goto L87
L125:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v741
	if int32(126) < v741 {
		goto L127
	} else {
		goto L128
	}
L127:
	;
	goto L124
L128:
	;
	v754 = m.G1
	v758 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v754+int32(_a_F_PutCoeffs_0)+v741))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v758
	v760 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v764 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v754+int32(_a_F_PutCoeffs_1)+v741))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v760 << (uint(v764) % 32)
	v767 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v768 = v764 + v767
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v768
	if v768 < int32(1) {
		goto L127
	} else {
		goto L129
	}
L129:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L127
L130:
	;
	v819 = m.G1
	v824 = int32(1024)
	v825 = v819 + int32(_a_F_PutCoeffs_6)
	v826 = int32(-67)
	goto L87
L131:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v793
	if int32(126) < v793 {
		goto L133
	} else {
		goto L134
	}
L132:
	;
	v788 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v790 = v782*v733>>(uint(int32(8))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v788 + v790
	v793 = v782 - v790
	goto L131
L133:
	;
	goto L130
L134:
	;
	v798 = m.G1
	v802 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v798+int32(_a_F_PutCoeffs_0)+v793))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v802
	v804 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v808 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v798+int32(_a_F_PutCoeffs_1)+v793))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v804 << (uint(v808) % 32)
	v811 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v812 = v808 + v811
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v812
	if v812 < int32(1) {
		goto L133
	} else {
		goto L135
	}
L135:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L133
L136:
	;
	v839 = int32(0)
	v841 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v833))))
	v843 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v846 = v843 * v841 >> (uint(int32(8)) % 32)
	if base.B2i32(v829&(v826+v142) != v839) == v839 {
		v855 = v846
		goto L139
	} else {
		goto L140
	}
L137:
	;
	goto L52
L138:
	;
	v880 = int32(1)
	if base.Ui32(v880) < base.Ui32(v829) {
		v829 = int32(base.Ui32(v829) >> (uint(v880) % 32))
		v833 = v833 + v880
		goto L136
	} else {
		goto L144
	}
L139:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v855
	if int32(126) < v855 {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v849 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v851 = v846 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v849 + v851
	v855 = v843 - v851
	goto L139
L141:
	;
	goto L138
L142:
	;
	v859 = m.G1
	v863 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v859+int32(_a_F_PutCoeffs_0)+v855))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v863
	v865 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v869 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v859+int32(_a_F_PutCoeffs_1)+v855))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v865 << (uint(v869) % 32)
	v872 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v873 = v869 + v872
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v873
	if v873 < int32(1) {
		goto L141
	} else {
		goto L143
	}
L143:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L141
L144:
	;
	goto L137
L145:
	;
	v943 = m.G81
	if v78 == int32(16) {
		v1027 = v134
		goto L9
	} else {
		goto L151
	}
L146:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v921
	if int32(126) < v921 {
		goto L148
	} else {
		goto L149
	}
L147:
	;
	v915 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v917 = v912 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v915 + v917
	v921 = v910 - v917
	goto L146
L148:
	;
	goto L145
L149:
	;
	v925 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v926 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v925 << (uint(v926) % 32)
	v929 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v929 + v926
	v933 = m.G1
	v937 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v933+int32(_a_F_PutCoeffs_0)+v921))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v937
	if v929 < int32(0) {
		goto L148
	} else {
		goto L150
	}
L150:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L148
L151:
	;
	v946 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	v947 = base.B2i32(v70 < v946)
	v949 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v943+v78))))
	v955 = v907 + v949*int32(33) + v905*int32(11)
	v956 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v955))))
	v958 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v961 = v958 * v956 >> (uint(int32(8)) % 32)
	if v947 == int32(0) {
		v970 = v961
		goto L153
	} else {
		goto L154
	}
L152:
	;
	if v947 == int32(0) {
		v1027 = v134
		goto L9
	} else {
		goto L158
	}
L153:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v970
	if int32(126) < v970 {
		goto L155
	} else {
		goto L156
	}
L154:
	;
	v964 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v966 = v961 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v964 + v966
	v970 = v958 - v966
	goto L153
L155:
	;
	goto L152
L156:
	;
	v974 = m.G1
	v978 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v974+int32(_a_F_PutCoeffs_0)+v970))))
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v978
	v980 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v984 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v974+int32(_a_F_PutCoeffs_1)+v970))))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v980 << (uint(v984) % 32)
	v987 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v988 = v984 + v987
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v988
	if v988 < int32(1) {
		goto L155
	} else {
		goto L157
	}
L157:
	;
	F_Flush(m, l0)
	mBase = m.M
	goto L155
L158:
	;
	v998 = v955
	goto L14
L159:
	;
	goto L13
}
func F_pow(m *base.Module, l0 float64, l1 float64) float64 {
	mBase := m.M
	_ = mBase
	var v19 int64
	_ = v19
	var v20 int64
	_ = v20
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v27 int64
	_ = v27
	var v30 int32
	_ = v30
	var v40 int64
	_ = v40
	var v45 float64
	_ = v45
	var v51 int64
	_ = v51
	var v67 float64
	_ = v67
	var v75 float64
	_ = v75
	var v87 int32
	_ = v87
	var v98 int64
	_ = v98
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v114 float64
	_ = v114
	var v115 float64
	_ = v115
	var v132 int32
	_ = v132
	var v143 int64
	_ = v143
	var v154 int32
	_ = v154
	var v155 int32
	_ = v155
	var v157 float64
	_ = v157
	var v169 int64
	_ = v169
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v174 float64
	_ = v174
	var v186 float64
	_ = v186
	var v190 float64
	_ = v190
	var v199 int64
	_ = v199
	var v201 int32
	_ = v201
	var v203 int64
	_ = v203
	var v205 float64
	_ = v205
	var v207 int64
	_ = v207
	var v210 int64
	_ = v210
	var v215 float64
	_ = v215
	var v222 int32
	_ = v222
	var v225 float64
	_ = v225
	var v228 float64
	_ = v228
	var v229 int32
	_ = v229
	var v230 float64
	_ = v230
	var v231 float64
	_ = v231
	var v232 float64
	_ = v232
	var v233 int64
	_ = v233
	var v236 float64
	_ = v236
	var v238 float64
	_ = v238
	var v242 float64
	_ = v242
	var v243 float64
	_ = v243
	var v246 float64
	_ = v246
	var v247 float64
	_ = v247
	var v248 float64
	_ = v248
	var v249 float64
	_ = v249
	var v252 float64
	_ = v252
	var v256 float64
	_ = v256
	var v260 float64
	_ = v260
	var v267 float64
	_ = v267
	var v270 float64
	_ = v270
	var v273 float64
	_ = v273
	var v277 float64
	_ = v277
	var v280 float64
	_ = v280
	var v285 float64
	_ = v285
	var v288 float64
	_ = v288
	var v292 float64
	_ = v292
	var v293 float64
	_ = v293
	var v297 float64
	_ = v297
	var v298 float64
	_ = v298
	var v299 int64
	_ = v299
	var v304 int32
	_ = v304
	var v313 float64
	_ = v313
	var v321 float64
	_ = v321
	var v324 float64
	_ = v324
	var v326 int32
	_ = v326
	var v336 int32
	_ = v336
	var v337 float64
	_ = v337
	var v340 float64
	_ = v340
	var v341 float64
	_ = v341
	var v342 float64
	_ = v342
	var v344 float64
	_ = v344
	var v347 float64
	_ = v347
	var v351 float64
	_ = v351
	var v352 float64
	_ = v352
	var v355 float64
	_ = v355
	var v358 float64
	_ = v358
	var v362 float64
	_ = v362
	var v365 float64
	_ = v365
	var v368 int64
	_ = v368
	var v373 int32
	_ = v373
	var v376 float64
	_ = v376
	var v379 float64
	_ = v379
	var v382 int64
	_ = v382
	var v387 int64
	_ = v387
	var v398 float64
	_ = v398
	var v404 int64
	_ = v404
	var v405 float64
	_ = v405
	var v406 float64
	_ = v406
	var v407 float64
	_ = v407
	var v420 float64
	_ = v420
	var v421 float64
	_ = v421
	var v428 float64
	_ = v428
	var v431 float64
	_ = v431
	var v432 float64
	_ = v432
	var v443 float64
	_ = v443
	var v445 float64
	_ = v445
	var v456 float64
	_ = v456
	v19 = base.I64_reinterpret_f64(l1)
	v20 = int64(52)
	v24 = base.I32_wrap_i64(int64(base.Ui64(v19)>>(uint(v20)%64))) & int32(2047)
	v26 = v24 + int32(-1086)
	v27 = base.I64_reinterpret_f64(l0)
	v30 = base.I32_wrap_i64(int64(base.Ui64(v27) >> (uint(v20) % 64)))
	if base.Ui32(v30+int32(-2047)) < base.Ui32(int32(-2046)) {
		v40 = v19 << (uint(int64(1)) % 64)
		if base.Ui64(v40+int64(-1)) < base.Ui64(int64(-9007199254740993)) {
			if base.Ui64(v27<<(uint(int64(1))%64)+int64(-1)) < base.Ui64(int64(-9007199254740993)) {
				if int64(-1) < v27 {
					v169 = v27
					v170 = v30
					v171 = int32(0)
					if base.Ui32(int32(-129)) < base.Ui32(v26) {
						if v170 != 0 {
							v199 = v169
							v201 = v171
						} else {
							v199 = base.I64_reinterpret_f64(base.F64_mul(l0, float64(4.503599627370496e+15)))&int64(9223372036854775807) + int64(-234187180623265792)
							v201 = v171
						}
						v203 = int64(-134217728)
						v205 = base.F64_reinterpret_i64(v19 & v203)
						v207 = v199 + int64(-4604531861337669632)
						v210 = v199 - v207&int64(-4503599627370496)
						v215 = base.F64_reinterpret_i64((v210 + int64(2147483648)) & int64(-4294967296))
						v222 = base.I32_wrap_i64(int64(base.Ui64(v207)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(5)) % 32)
						v225 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[0])))
						v228 = base.F64_add(base.F64_mul(v215, v225), float64(-1))
						v229 = int32(0)
						v230 = *(*float64)(unsafe.Add(mBase, _c_F_pow[1]))
						v231 = base.F64_mul(v228, v230)
						v232 = base.F64_mul(v228, v231)
						v233 = int64(52)
						v236 = base.F64_convert_i32_s(base.I32_wrap_i64(v207 >> (uint(v233) % 64)))
						v238 = *(*float64)(unsafe.Add(mBase, _c_F_pow[2]))
						v242 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[3])))
						v243 = base.F64_add(base.F64_mul(v236, v238), v242)
						v246 = base.F64_mul(v225, base.F64_sub(base.F64_reinterpret_i64(v210), v215))
						v247 = base.F64_add(v228, v246)
						v248 = base.F64_add(v243, v247)
						v249 = base.F64_add(v232, v248)
						v252 = base.F64_mul(v230, v247)
						v256 = *(*float64)(unsafe.Add(mBase, _c_F_pow[4]))
						v260 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[5])))
						v267 = base.F64_mul(v247, v252)
						v270 = *(*float64)(unsafe.Add(mBase, _c_F_pow[6]))
						v273 = *(*float64)(unsafe.Add(mBase, _c_F_pow[7]))
						v277 = *(*float64)(unsafe.Add(mBase, _c_F_pow[8]))
						v280 = *(*float64)(unsafe.Add(mBase, _c_F_pow[9]))
						v285 = *(*float64)(unsafe.Add(mBase, _c_F_pow[10]))
						v288 = *(*float64)(unsafe.Add(mBase, _c_F_pow[11]))
						v292 = base.F64_add(base.F64_add(base.F64_add(v232, base.F64_sub(v248, v249)), base.F64_add(base.F64_mul(v246, base.F64_add(v231, v252)), base.F64_add(base.F64_add(base.F64_mul(v236, v256), v260), base.F64_add(v247, base.F64_sub(v243, v248))))), base.F64_mul(base.F64_mul(v247, v267), base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v247, v270), v273)), base.F64_add(base.F64_mul(v247, v277), v280))), base.F64_add(base.F64_mul(v247, v285), v288))))
						v293 = base.F64_add(v249, v292)
						v297 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(v293) & v203)
						v298 = base.F64_mul(v205, v297)
						v299 = base.I64_reinterpret_f64(v298)
						v304 = base.I32_wrap_i64(int64(base.Ui64(v299)>>(uint(v233)%64))) & int32(2047)
						if base.Ui32(v304+int32(-969)) < base.Ui32(int32(63)) {
							v326 = v304
							v336 = int32(0)
							v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
							v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
							v341 = base.F64_add(base.F64_mul(v298, v337), v340)
							v342 = base.F64_sub(v341, v340)
							v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
							v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
							v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
							v352 = base.F64_mul(v351, v351)
							v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
							v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
							v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
							v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
							v368 = base.I64_reinterpret_f64(v341)
							v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
							v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
							v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
							v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
							v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
							if v326 != 0 {
								v445 = base.F64_reinterpret_i64(v387)
								v456 = base.F64_add(base.F64_mul(v445, v379), v445)
								return v456
							} else {
								if v368&int64(2147483648) != int64(0) {
									v404 = v387 + int64(4602678819172646912)
									v405 = base.F64_reinterpret_i64(v404)
									v406 = base.F64_mul(v405, v379)
									v407 = base.F64_add(v406, v405)
									if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
										v432 = v407
									} else {
										if base.F64_lt(v407, float64(0)) != 0 {
											v420 = float64(-1)
										} else {
											v420 = float64(1)
										}
										v421 = base.F64_add(v407, v420)
										v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
										if base.F64_eq(v428, float64(0)) != 0 {
											v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
										} else {
											v431 = v428
										}
										v432 = v431
									}
									v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
								} else {
									v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
									v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
								}
								return v443
							}
						} else {
							if base.Ui32(int32(968)) < base.Ui32(v304) {
								if base.Ui32(v304) < base.Ui32(int32(1033)) {
									v326 = int32(0)
									v336 = int32(0)
									v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
									v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
									v341 = base.F64_add(base.F64_mul(v298, v337), v340)
									v342 = base.F64_sub(v341, v340)
									v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
									v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
									v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
									v352 = base.F64_mul(v351, v351)
									v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
									v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
									v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
									v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
									v368 = base.I64_reinterpret_f64(v341)
									v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
									v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
									v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
									v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
									v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
									if v326 != 0 {
										v445 = base.F64_reinterpret_i64(v387)
										v456 = base.F64_add(base.F64_mul(v445, v379), v445)
										return v456
									} else {
										if v368&int64(2147483648) != int64(0) {
											v404 = v387 + int64(4602678819172646912)
											v405 = base.F64_reinterpret_i64(v404)
											v406 = base.F64_mul(v405, v379)
											v407 = base.F64_add(v406, v405)
											if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
												v432 = v407
											} else {
												if base.F64_lt(v407, float64(0)) != 0 {
													v420 = float64(-1)
												} else {
													v420 = float64(1)
												}
												v421 = base.F64_add(v407, v420)
												v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
												if base.F64_eq(v428, float64(0)) != 0 {
													v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
												} else {
													v431 = v428
												}
												v432 = v431
											}
											v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
										} else {
											v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
											v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
										}
										return v443
									}
								} else {
									if int64(-1) < v299 {
										v324 = F___math_xflow(m, v201, float64(3.105036184601418e+231))
										mBase = m.M
										return v324
									} else {
										v321 = F___math_xflow(m, v201, float64(1.2882297539194267e-231))
										mBase = m.M
										return v321
									}
								}
							} else {
								if v201 != 0 {
									v313 = float64(-1)
								} else {
									v313 = float64(1)
								}
								return v313
							}
						}
					} else {
						v174 = float64(1)
						if v169 == int64(4607182418800017408) {
							v456 = v174
							return v456
						} else {
							if base.Ui32(v24) < base.Ui32(int32(958)) {
								v456 = v174
								return v456
							} else {
								if base.B2i32(v19 < int64(0)) == base.B2i32(base.Ui64(int64(4607182418800017408)) < base.Ui64(v169)) {
									v190 = F___math_xflow(m, int32(0), float64(1.2882297539194267e-231))
									mBase = m.M
									return v190
								} else {
									v186 = F___math_xflow(m, int32(0), float64(3.105036184601418e+231))
									mBase = m.M
									return v186
								}
							}
						}
					}
				} else {
					v132 = base.I32_wrap_i64(int64(base.Ui64(v19)>>(uint(int64(52))%64))) & int32(2047)
					if base.Ui32(v132) < base.Ui32(int32(1023)) {
						v155 = int32(0)
					} else {
						if base.Ui32(int32(1075)) < base.Ui32(v132) {
							v155 = int32(2)
						} else {
							v143 = int64(1) << (uint(base.I64_extend_i32_u(int32(1075)-v132)) % 64)
							if (v143+int64(-1))&v19 != int64(0) {
								v155 = int32(0)
							} else {
								if v143&v19 == int64(0) {
									v154 = int32(2)
								} else {
									v154 = int32(1)
								}
								v155 = v154
							}
						}
					}
					if v155 != 0 {
						v169 = base.I64_reinterpret_f64(l0) & int64(9223372036854775807)
						v170 = v30 & int32(2047)
						v171 = base.B2i32(v155 == int32(1)) << (uint(int32(18)) % 32)
						if base.Ui32(int32(-129)) < base.Ui32(v26) {
							if v170 != 0 {
								v199 = v169
								v201 = v171
							} else {
								v199 = base.I64_reinterpret_f64(base.F64_mul(l0, float64(4.503599627370496e+15)))&int64(9223372036854775807) + int64(-234187180623265792)
								v201 = v171
							}
							v203 = int64(-134217728)
							v205 = base.F64_reinterpret_i64(v19 & v203)
							v207 = v199 + int64(-4604531861337669632)
							v210 = v199 - v207&int64(-4503599627370496)
							v215 = base.F64_reinterpret_i64((v210 + int64(2147483648)) & int64(-4294967296))
							v222 = base.I32_wrap_i64(int64(base.Ui64(v207)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(5)) % 32)
							v225 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[0])))
							v228 = base.F64_add(base.F64_mul(v215, v225), float64(-1))
							v229 = int32(0)
							v230 = *(*float64)(unsafe.Add(mBase, _c_F_pow[1]))
							v231 = base.F64_mul(v228, v230)
							v232 = base.F64_mul(v228, v231)
							v233 = int64(52)
							v236 = base.F64_convert_i32_s(base.I32_wrap_i64(v207 >> (uint(v233) % 64)))
							v238 = *(*float64)(unsafe.Add(mBase, _c_F_pow[2]))
							v242 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[3])))
							v243 = base.F64_add(base.F64_mul(v236, v238), v242)
							v246 = base.F64_mul(v225, base.F64_sub(base.F64_reinterpret_i64(v210), v215))
							v247 = base.F64_add(v228, v246)
							v248 = base.F64_add(v243, v247)
							v249 = base.F64_add(v232, v248)
							v252 = base.F64_mul(v230, v247)
							v256 = *(*float64)(unsafe.Add(mBase, _c_F_pow[4]))
							v260 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[5])))
							v267 = base.F64_mul(v247, v252)
							v270 = *(*float64)(unsafe.Add(mBase, _c_F_pow[6]))
							v273 = *(*float64)(unsafe.Add(mBase, _c_F_pow[7]))
							v277 = *(*float64)(unsafe.Add(mBase, _c_F_pow[8]))
							v280 = *(*float64)(unsafe.Add(mBase, _c_F_pow[9]))
							v285 = *(*float64)(unsafe.Add(mBase, _c_F_pow[10]))
							v288 = *(*float64)(unsafe.Add(mBase, _c_F_pow[11]))
							v292 = base.F64_add(base.F64_add(base.F64_add(v232, base.F64_sub(v248, v249)), base.F64_add(base.F64_mul(v246, base.F64_add(v231, v252)), base.F64_add(base.F64_add(base.F64_mul(v236, v256), v260), base.F64_add(v247, base.F64_sub(v243, v248))))), base.F64_mul(base.F64_mul(v247, v267), base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v247, v270), v273)), base.F64_add(base.F64_mul(v247, v277), v280))), base.F64_add(base.F64_mul(v247, v285), v288))))
							v293 = base.F64_add(v249, v292)
							v297 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(v293) & v203)
							v298 = base.F64_mul(v205, v297)
							v299 = base.I64_reinterpret_f64(v298)
							v304 = base.I32_wrap_i64(int64(base.Ui64(v299)>>(uint(v233)%64))) & int32(2047)
							if base.Ui32(v304+int32(-969)) < base.Ui32(int32(63)) {
								v326 = v304
								v336 = int32(0)
								v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
								v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
								v341 = base.F64_add(base.F64_mul(v298, v337), v340)
								v342 = base.F64_sub(v341, v340)
								v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
								v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
								v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
								v352 = base.F64_mul(v351, v351)
								v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
								v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
								v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
								v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
								v368 = base.I64_reinterpret_f64(v341)
								v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
								v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
								v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
								v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
								v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
								if v326 != 0 {
									v445 = base.F64_reinterpret_i64(v387)
									v456 = base.F64_add(base.F64_mul(v445, v379), v445)
									return v456
								} else {
									if v368&int64(2147483648) != int64(0) {
										v404 = v387 + int64(4602678819172646912)
										v405 = base.F64_reinterpret_i64(v404)
										v406 = base.F64_mul(v405, v379)
										v407 = base.F64_add(v406, v405)
										if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
											v432 = v407
										} else {
											if base.F64_lt(v407, float64(0)) != 0 {
												v420 = float64(-1)
											} else {
												v420 = float64(1)
											}
											v421 = base.F64_add(v407, v420)
											v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
											if base.F64_eq(v428, float64(0)) != 0 {
												v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
											} else {
												v431 = v428
											}
											v432 = v431
										}
										v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
									} else {
										v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
										v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
									}
									return v443
								}
							} else {
								if base.Ui32(int32(968)) < base.Ui32(v304) {
									if base.Ui32(v304) < base.Ui32(int32(1033)) {
										v326 = int32(0)
										v336 = int32(0)
										v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
										v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
										v341 = base.F64_add(base.F64_mul(v298, v337), v340)
										v342 = base.F64_sub(v341, v340)
										v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
										v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
										v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
										v352 = base.F64_mul(v351, v351)
										v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
										v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
										v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
										v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
										v368 = base.I64_reinterpret_f64(v341)
										v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
										v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
										v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
										v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
										v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
										if v326 != 0 {
											v445 = base.F64_reinterpret_i64(v387)
											v456 = base.F64_add(base.F64_mul(v445, v379), v445)
											return v456
										} else {
											if v368&int64(2147483648) != int64(0) {
												v404 = v387 + int64(4602678819172646912)
												v405 = base.F64_reinterpret_i64(v404)
												v406 = base.F64_mul(v405, v379)
												v407 = base.F64_add(v406, v405)
												if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
													v432 = v407
												} else {
													if base.F64_lt(v407, float64(0)) != 0 {
														v420 = float64(-1)
													} else {
														v420 = float64(1)
													}
													v421 = base.F64_add(v407, v420)
													v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
													if base.F64_eq(v428, float64(0)) != 0 {
														v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
													} else {
														v431 = v428
													}
													v432 = v431
												}
												v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
											} else {
												v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
												v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
											}
											return v443
										}
									} else {
										if int64(-1) < v299 {
											v324 = F___math_xflow(m, v201, float64(3.105036184601418e+231))
											mBase = m.M
											return v324
										} else {
											v321 = F___math_xflow(m, v201, float64(1.2882297539194267e-231))
											mBase = m.M
											return v321
										}
									}
								} else {
									if v201 != 0 {
										v313 = float64(-1)
									} else {
										v313 = float64(1)
									}
									return v313
								}
							}
						} else {
							v174 = float64(1)
							if v169 == int64(4607182418800017408) {
								v456 = v174
								return v456
							} else {
								if base.Ui32(v24) < base.Ui32(int32(958)) {
									v456 = v174
									return v456
								} else {
									if base.B2i32(v19 < int64(0)) == base.B2i32(base.Ui64(int64(4607182418800017408)) < base.Ui64(v169)) {
										v190 = F___math_xflow(m, int32(0), float64(1.2882297539194267e-231))
										mBase = m.M
										return v190
									} else {
										v186 = F___math_xflow(m, int32(0), float64(3.105036184601418e+231))
										mBase = m.M
										return v186
									}
								}
							}
						}
					} else {
						v157 = base.F64_sub(l0, l0)
						return base.F64_div(v157, v157)
					}
				}
			} else {
				v75 = base.F64_mul(l0, l0)
				if int64(-1) < v27 {
					v115 = v75
				} else {
					v87 = base.I32_wrap_i64(int64(base.Ui64(v19)>>(uint(int64(52))%64))) & int32(2047)
					if base.Ui32(v87) < base.Ui32(int32(1023)) {
						v110 = int32(0)
					} else {
						if base.Ui32(int32(1075)) < base.Ui32(v87) {
							v110 = int32(2)
						} else {
							v98 = int64(1) << (uint(base.I64_extend_i32_u(int32(1075)-v87)) % 64)
							if (v98+int64(-1))&v19 != int64(0) {
								v110 = int32(0)
							} else {
								if v98&v19 == int64(0) {
									v109 = int32(2)
								} else {
									v109 = int32(1)
								}
								v110 = v109
							}
						}
					}
					if v110 == int32(1) {
						v114 = base.F64_neg(v75)
					} else {
						v114 = v75
					}
					v115 = v114
				}
				if int64(-1) < v19 {
					v456 = v115
					return v456
				} else {
					return base.F64_div(float64(1), v115)
				}
			}
		} else {
			v45 = float64(1)
			if v27 == int64(4607182418800017408) {
				v456 = v45
				return v456
			} else {
				if v40 == int64(0) {
					v456 = v45
					return v456
				} else {
					v51 = v27 << (uint(int64(1)) % 64)
					if base.Ui64(int64(-9007199254740992)) < base.Ui64(v51) {
						return base.F64_add(l0, l1)
					} else {
						if base.Ui64(v40) < base.Ui64(int64(-9007199254740991)) {
							if v51 == int64(9214364837600034816) {
								v456 = v45
								return v456
							} else {
								if base.B2i32(base.Ui64(v51) < base.Ui64(int64(9214364837600034816)))^base.B2i32(v19 < int64(0)) != 0 {
									v67 = float64(0)
								} else {
									v67 = base.F64_mul(l1, l1)
								}
								return v67
							}
						} else {
							return base.F64_add(l0, l1)
						}
					}
				}
			}
		}
	} else {
		if base.Ui32(int32(-129)) < base.Ui32(v26) {
			v199 = v27
			v201 = int32(0)
			v203 = int64(-134217728)
			v205 = base.F64_reinterpret_i64(v19 & v203)
			v207 = v199 + int64(-4604531861337669632)
			v210 = v199 - v207&int64(-4503599627370496)
			v215 = base.F64_reinterpret_i64((v210 + int64(2147483648)) & int64(-4294967296))
			v222 = base.I32_wrap_i64(int64(base.Ui64(v207)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(5)) % 32)
			v225 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[0])))
			v228 = base.F64_add(base.F64_mul(v215, v225), float64(-1))
			v229 = int32(0)
			v230 = *(*float64)(unsafe.Add(mBase, _c_F_pow[1]))
			v231 = base.F64_mul(v228, v230)
			v232 = base.F64_mul(v228, v231)
			v233 = int64(52)
			v236 = base.F64_convert_i32_s(base.I32_wrap_i64(v207 >> (uint(v233) % 64)))
			v238 = *(*float64)(unsafe.Add(mBase, _c_F_pow[2]))
			v242 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[3])))
			v243 = base.F64_add(base.F64_mul(v236, v238), v242)
			v246 = base.F64_mul(v225, base.F64_sub(base.F64_reinterpret_i64(v210), v215))
			v247 = base.F64_add(v228, v246)
			v248 = base.F64_add(v243, v247)
			v249 = base.F64_add(v232, v248)
			v252 = base.F64_mul(v230, v247)
			v256 = *(*float64)(unsafe.Add(mBase, _c_F_pow[4]))
			v260 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[5])))
			v267 = base.F64_mul(v247, v252)
			v270 = *(*float64)(unsafe.Add(mBase, _c_F_pow[6]))
			v273 = *(*float64)(unsafe.Add(mBase, _c_F_pow[7]))
			v277 = *(*float64)(unsafe.Add(mBase, _c_F_pow[8]))
			v280 = *(*float64)(unsafe.Add(mBase, _c_F_pow[9]))
			v285 = *(*float64)(unsafe.Add(mBase, _c_F_pow[10]))
			v288 = *(*float64)(unsafe.Add(mBase, _c_F_pow[11]))
			v292 = base.F64_add(base.F64_add(base.F64_add(v232, base.F64_sub(v248, v249)), base.F64_add(base.F64_mul(v246, base.F64_add(v231, v252)), base.F64_add(base.F64_add(base.F64_mul(v236, v256), v260), base.F64_add(v247, base.F64_sub(v243, v248))))), base.F64_mul(base.F64_mul(v247, v267), base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v247, v270), v273)), base.F64_add(base.F64_mul(v247, v277), v280))), base.F64_add(base.F64_mul(v247, v285), v288))))
			v293 = base.F64_add(v249, v292)
			v297 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(v293) & v203)
			v298 = base.F64_mul(v205, v297)
			v299 = base.I64_reinterpret_f64(v298)
			v304 = base.I32_wrap_i64(int64(base.Ui64(v299)>>(uint(v233)%64))) & int32(2047)
			if base.Ui32(v304+int32(-969)) < base.Ui32(int32(63)) {
				v326 = v304
				v336 = int32(0)
				v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
				v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
				v341 = base.F64_add(base.F64_mul(v298, v337), v340)
				v342 = base.F64_sub(v341, v340)
				v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
				v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
				v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
				v352 = base.F64_mul(v351, v351)
				v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
				v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
				v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
				v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
				v368 = base.I64_reinterpret_f64(v341)
				v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
				v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
				v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
				v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
				v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
				if v326 != 0 {
					v445 = base.F64_reinterpret_i64(v387)
					v456 = base.F64_add(base.F64_mul(v445, v379), v445)
					return v456
				} else {
					if v368&int64(2147483648) != int64(0) {
						v404 = v387 + int64(4602678819172646912)
						v405 = base.F64_reinterpret_i64(v404)
						v406 = base.F64_mul(v405, v379)
						v407 = base.F64_add(v406, v405)
						if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
							v432 = v407
						} else {
							if base.F64_lt(v407, float64(0)) != 0 {
								v420 = float64(-1)
							} else {
								v420 = float64(1)
							}
							v421 = base.F64_add(v407, v420)
							v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
							if base.F64_eq(v428, float64(0)) != 0 {
								v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
							} else {
								v431 = v428
							}
							v432 = v431
						}
						v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
					} else {
						v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
						v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
					}
					return v443
				}
			} else {
				if base.Ui32(int32(968)) < base.Ui32(v304) {
					if base.Ui32(v304) < base.Ui32(int32(1033)) {
						v326 = int32(0)
						v336 = int32(0)
						v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
						v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
						v341 = base.F64_add(base.F64_mul(v298, v337), v340)
						v342 = base.F64_sub(v341, v340)
						v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
						v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
						v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
						v352 = base.F64_mul(v351, v351)
						v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
						v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
						v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
						v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
						v368 = base.I64_reinterpret_f64(v341)
						v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
						v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
						v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
						v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
						v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
						if v326 != 0 {
							v445 = base.F64_reinterpret_i64(v387)
							v456 = base.F64_add(base.F64_mul(v445, v379), v445)
							return v456
						} else {
							if v368&int64(2147483648) != int64(0) {
								v404 = v387 + int64(4602678819172646912)
								v405 = base.F64_reinterpret_i64(v404)
								v406 = base.F64_mul(v405, v379)
								v407 = base.F64_add(v406, v405)
								if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
									v432 = v407
								} else {
									if base.F64_lt(v407, float64(0)) != 0 {
										v420 = float64(-1)
									} else {
										v420 = float64(1)
									}
									v421 = base.F64_add(v407, v420)
									v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
									if base.F64_eq(v428, float64(0)) != 0 {
										v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
									} else {
										v431 = v428
									}
									v432 = v431
								}
								v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
							} else {
								v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
								v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
							}
							return v443
						}
					} else {
						if int64(-1) < v299 {
							v324 = F___math_xflow(m, v201, float64(3.105036184601418e+231))
							mBase = m.M
							return v324
						} else {
							v321 = F___math_xflow(m, v201, float64(1.2882297539194267e-231))
							mBase = m.M
							return v321
						}
					}
				} else {
					if v201 != 0 {
						v313 = float64(-1)
					} else {
						v313 = float64(1)
					}
					return v313
				}
			}
		} else {
			v40 = v19 << (uint(int64(1)) % 64)
			if base.Ui64(v40+int64(-1)) < base.Ui64(int64(-9007199254740993)) {
				if base.Ui64(v27<<(uint(int64(1))%64)+int64(-1)) < base.Ui64(int64(-9007199254740993)) {
					if int64(-1) < v27 {
						v169 = v27
						v170 = v30
						v171 = int32(0)
						if base.Ui32(int32(-129)) < base.Ui32(v26) {
							if v170 != 0 {
								v199 = v169
								v201 = v171
							} else {
								v199 = base.I64_reinterpret_f64(base.F64_mul(l0, float64(4.503599627370496e+15)))&int64(9223372036854775807) + int64(-234187180623265792)
								v201 = v171
							}
							v203 = int64(-134217728)
							v205 = base.F64_reinterpret_i64(v19 & v203)
							v207 = v199 + int64(-4604531861337669632)
							v210 = v199 - v207&int64(-4503599627370496)
							v215 = base.F64_reinterpret_i64((v210 + int64(2147483648)) & int64(-4294967296))
							v222 = base.I32_wrap_i64(int64(base.Ui64(v207)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(5)) % 32)
							v225 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[0])))
							v228 = base.F64_add(base.F64_mul(v215, v225), float64(-1))
							v229 = int32(0)
							v230 = *(*float64)(unsafe.Add(mBase, _c_F_pow[1]))
							v231 = base.F64_mul(v228, v230)
							v232 = base.F64_mul(v228, v231)
							v233 = int64(52)
							v236 = base.F64_convert_i32_s(base.I32_wrap_i64(v207 >> (uint(v233) % 64)))
							v238 = *(*float64)(unsafe.Add(mBase, _c_F_pow[2]))
							v242 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[3])))
							v243 = base.F64_add(base.F64_mul(v236, v238), v242)
							v246 = base.F64_mul(v225, base.F64_sub(base.F64_reinterpret_i64(v210), v215))
							v247 = base.F64_add(v228, v246)
							v248 = base.F64_add(v243, v247)
							v249 = base.F64_add(v232, v248)
							v252 = base.F64_mul(v230, v247)
							v256 = *(*float64)(unsafe.Add(mBase, _c_F_pow[4]))
							v260 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[5])))
							v267 = base.F64_mul(v247, v252)
							v270 = *(*float64)(unsafe.Add(mBase, _c_F_pow[6]))
							v273 = *(*float64)(unsafe.Add(mBase, _c_F_pow[7]))
							v277 = *(*float64)(unsafe.Add(mBase, _c_F_pow[8]))
							v280 = *(*float64)(unsafe.Add(mBase, _c_F_pow[9]))
							v285 = *(*float64)(unsafe.Add(mBase, _c_F_pow[10]))
							v288 = *(*float64)(unsafe.Add(mBase, _c_F_pow[11]))
							v292 = base.F64_add(base.F64_add(base.F64_add(v232, base.F64_sub(v248, v249)), base.F64_add(base.F64_mul(v246, base.F64_add(v231, v252)), base.F64_add(base.F64_add(base.F64_mul(v236, v256), v260), base.F64_add(v247, base.F64_sub(v243, v248))))), base.F64_mul(base.F64_mul(v247, v267), base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v247, v270), v273)), base.F64_add(base.F64_mul(v247, v277), v280))), base.F64_add(base.F64_mul(v247, v285), v288))))
							v293 = base.F64_add(v249, v292)
							v297 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(v293) & v203)
							v298 = base.F64_mul(v205, v297)
							v299 = base.I64_reinterpret_f64(v298)
							v304 = base.I32_wrap_i64(int64(base.Ui64(v299)>>(uint(v233)%64))) & int32(2047)
							if base.Ui32(v304+int32(-969)) < base.Ui32(int32(63)) {
								v326 = v304
								v336 = int32(0)
								v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
								v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
								v341 = base.F64_add(base.F64_mul(v298, v337), v340)
								v342 = base.F64_sub(v341, v340)
								v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
								v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
								v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
								v352 = base.F64_mul(v351, v351)
								v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
								v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
								v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
								v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
								v368 = base.I64_reinterpret_f64(v341)
								v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
								v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
								v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
								v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
								v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
								if v326 != 0 {
									v445 = base.F64_reinterpret_i64(v387)
									v456 = base.F64_add(base.F64_mul(v445, v379), v445)
									return v456
								} else {
									if v368&int64(2147483648) != int64(0) {
										v404 = v387 + int64(4602678819172646912)
										v405 = base.F64_reinterpret_i64(v404)
										v406 = base.F64_mul(v405, v379)
										v407 = base.F64_add(v406, v405)
										if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
											v432 = v407
										} else {
											if base.F64_lt(v407, float64(0)) != 0 {
												v420 = float64(-1)
											} else {
												v420 = float64(1)
											}
											v421 = base.F64_add(v407, v420)
											v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
											if base.F64_eq(v428, float64(0)) != 0 {
												v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
											} else {
												v431 = v428
											}
											v432 = v431
										}
										v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
									} else {
										v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
										v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
									}
									return v443
								}
							} else {
								if base.Ui32(int32(968)) < base.Ui32(v304) {
									if base.Ui32(v304) < base.Ui32(int32(1033)) {
										v326 = int32(0)
										v336 = int32(0)
										v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
										v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
										v341 = base.F64_add(base.F64_mul(v298, v337), v340)
										v342 = base.F64_sub(v341, v340)
										v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
										v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
										v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
										v352 = base.F64_mul(v351, v351)
										v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
										v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
										v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
										v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
										v368 = base.I64_reinterpret_f64(v341)
										v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
										v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
										v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
										v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
										v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
										if v326 != 0 {
											v445 = base.F64_reinterpret_i64(v387)
											v456 = base.F64_add(base.F64_mul(v445, v379), v445)
											return v456
										} else {
											if v368&int64(2147483648) != int64(0) {
												v404 = v387 + int64(4602678819172646912)
												v405 = base.F64_reinterpret_i64(v404)
												v406 = base.F64_mul(v405, v379)
												v407 = base.F64_add(v406, v405)
												if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
													v432 = v407
												} else {
													if base.F64_lt(v407, float64(0)) != 0 {
														v420 = float64(-1)
													} else {
														v420 = float64(1)
													}
													v421 = base.F64_add(v407, v420)
													v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
													if base.F64_eq(v428, float64(0)) != 0 {
														v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
													} else {
														v431 = v428
													}
													v432 = v431
												}
												v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
											} else {
												v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
												v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
											}
											return v443
										}
									} else {
										if int64(-1) < v299 {
											v324 = F___math_xflow(m, v201, float64(3.105036184601418e+231))
											mBase = m.M
											return v324
										} else {
											v321 = F___math_xflow(m, v201, float64(1.2882297539194267e-231))
											mBase = m.M
											return v321
										}
									}
								} else {
									if v201 != 0 {
										v313 = float64(-1)
									} else {
										v313 = float64(1)
									}
									return v313
								}
							}
						} else {
							v174 = float64(1)
							if v169 == int64(4607182418800017408) {
								v456 = v174
								return v456
							} else {
								if base.Ui32(v24) < base.Ui32(int32(958)) {
									v456 = v174
									return v456
								} else {
									if base.B2i32(v19 < int64(0)) == base.B2i32(base.Ui64(int64(4607182418800017408)) < base.Ui64(v169)) {
										v190 = F___math_xflow(m, int32(0), float64(1.2882297539194267e-231))
										mBase = m.M
										return v190
									} else {
										v186 = F___math_xflow(m, int32(0), float64(3.105036184601418e+231))
										mBase = m.M
										return v186
									}
								}
							}
						}
					} else {
						v132 = base.I32_wrap_i64(int64(base.Ui64(v19)>>(uint(int64(52))%64))) & int32(2047)
						if base.Ui32(v132) < base.Ui32(int32(1023)) {
							v155 = int32(0)
						} else {
							if base.Ui32(int32(1075)) < base.Ui32(v132) {
								v155 = int32(2)
							} else {
								v143 = int64(1) << (uint(base.I64_extend_i32_u(int32(1075)-v132)) % 64)
								if (v143+int64(-1))&v19 != int64(0) {
									v155 = int32(0)
								} else {
									if v143&v19 == int64(0) {
										v154 = int32(2)
									} else {
										v154 = int32(1)
									}
									v155 = v154
								}
							}
						}
						if v155 != 0 {
							v169 = base.I64_reinterpret_f64(l0) & int64(9223372036854775807)
							v170 = v30 & int32(2047)
							v171 = base.B2i32(v155 == int32(1)) << (uint(int32(18)) % 32)
							if base.Ui32(int32(-129)) < base.Ui32(v26) {
								if v170 != 0 {
									v199 = v169
									v201 = v171
								} else {
									v199 = base.I64_reinterpret_f64(base.F64_mul(l0, float64(4.503599627370496e+15)))&int64(9223372036854775807) + int64(-234187180623265792)
									v201 = v171
								}
								v203 = int64(-134217728)
								v205 = base.F64_reinterpret_i64(v19 & v203)
								v207 = v199 + int64(-4604531861337669632)
								v210 = v199 - v207&int64(-4503599627370496)
								v215 = base.F64_reinterpret_i64((v210 + int64(2147483648)) & int64(-4294967296))
								v222 = base.I32_wrap_i64(int64(base.Ui64(v207)>>(uint(int64(45))%64))) & int32(127) << (uint(int32(5)) % 32)
								v225 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[0])))
								v228 = base.F64_add(base.F64_mul(v215, v225), float64(-1))
								v229 = int32(0)
								v230 = *(*float64)(unsafe.Add(mBase, _c_F_pow[1]))
								v231 = base.F64_mul(v228, v230)
								v232 = base.F64_mul(v228, v231)
								v233 = int64(52)
								v236 = base.F64_convert_i32_s(base.I32_wrap_i64(v207 >> (uint(v233) % 64)))
								v238 = *(*float64)(unsafe.Add(mBase, _c_F_pow[2]))
								v242 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[3])))
								v243 = base.F64_add(base.F64_mul(v236, v238), v242)
								v246 = base.F64_mul(v225, base.F64_sub(base.F64_reinterpret_i64(v210), v215))
								v247 = base.F64_add(v228, v246)
								v248 = base.F64_add(v243, v247)
								v249 = base.F64_add(v232, v248)
								v252 = base.F64_mul(v230, v247)
								v256 = *(*float64)(unsafe.Add(mBase, _c_F_pow[4]))
								v260 = *(*float64)(unsafe.Add(mBase, uint32(v222)+uint32(_c_F_pow[5])))
								v267 = base.F64_mul(v247, v252)
								v270 = *(*float64)(unsafe.Add(mBase, _c_F_pow[6]))
								v273 = *(*float64)(unsafe.Add(mBase, _c_F_pow[7]))
								v277 = *(*float64)(unsafe.Add(mBase, _c_F_pow[8]))
								v280 = *(*float64)(unsafe.Add(mBase, _c_F_pow[9]))
								v285 = *(*float64)(unsafe.Add(mBase, _c_F_pow[10]))
								v288 = *(*float64)(unsafe.Add(mBase, _c_F_pow[11]))
								v292 = base.F64_add(base.F64_add(base.F64_add(v232, base.F64_sub(v248, v249)), base.F64_add(base.F64_mul(v246, base.F64_add(v231, v252)), base.F64_add(base.F64_add(base.F64_mul(v236, v256), v260), base.F64_add(v247, base.F64_sub(v243, v248))))), base.F64_mul(base.F64_mul(v247, v267), base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v267, base.F64_add(base.F64_mul(v247, v270), v273)), base.F64_add(base.F64_mul(v247, v277), v280))), base.F64_add(base.F64_mul(v247, v285), v288))))
								v293 = base.F64_add(v249, v292)
								v297 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(v293) & v203)
								v298 = base.F64_mul(v205, v297)
								v299 = base.I64_reinterpret_f64(v298)
								v304 = base.I32_wrap_i64(int64(base.Ui64(v299)>>(uint(v233)%64))) & int32(2047)
								if base.Ui32(v304+int32(-969)) < base.Ui32(int32(63)) {
									v326 = v304
									v336 = int32(0)
									v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
									v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
									v341 = base.F64_add(base.F64_mul(v298, v337), v340)
									v342 = base.F64_sub(v341, v340)
									v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
									v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
									v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
									v352 = base.F64_mul(v351, v351)
									v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
									v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
									v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
									v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
									v368 = base.I64_reinterpret_f64(v341)
									v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
									v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
									v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
									v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
									v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
									if v326 != 0 {
										v445 = base.F64_reinterpret_i64(v387)
										v456 = base.F64_add(base.F64_mul(v445, v379), v445)
										return v456
									} else {
										if v368&int64(2147483648) != int64(0) {
											v404 = v387 + int64(4602678819172646912)
											v405 = base.F64_reinterpret_i64(v404)
											v406 = base.F64_mul(v405, v379)
											v407 = base.F64_add(v406, v405)
											if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
												v432 = v407
											} else {
												if base.F64_lt(v407, float64(0)) != 0 {
													v420 = float64(-1)
												} else {
													v420 = float64(1)
												}
												v421 = base.F64_add(v407, v420)
												v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
												if base.F64_eq(v428, float64(0)) != 0 {
													v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
												} else {
													v431 = v428
												}
												v432 = v431
											}
											v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
										} else {
											v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
											v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
										}
										return v443
									}
								} else {
									if base.Ui32(int32(968)) < base.Ui32(v304) {
										if base.Ui32(v304) < base.Ui32(int32(1033)) {
											v326 = int32(0)
											v336 = int32(0)
											v337 = *(*float64)(unsafe.Add(mBase, _c_F_pow[12]))
											v340 = *(*float64)(unsafe.Add(mBase, _c_F_pow[13]))
											v341 = base.F64_add(base.F64_mul(v298, v337), v340)
											v342 = base.F64_sub(v341, v340)
											v344 = *(*float64)(unsafe.Add(mBase, _c_F_pow[14]))
											v347 = *(*float64)(unsafe.Add(mBase, _c_F_pow[15]))
											v351 = base.F64_add(base.F64_add(base.F64_mul(base.F64_sub(l1, v205), v297), base.F64_mul(base.F64_add(base.F64_add(v292, base.F64_sub(v249, v293)), base.F64_sub(v293, v297)), l1)), base.F64_add(base.F64_mul(v342, v344), base.F64_add(base.F64_mul(v342, v347), v298)))
											v352 = base.F64_mul(v351, v351)
											v355 = *(*float64)(unsafe.Add(mBase, _c_F_pow[16]))
											v358 = *(*float64)(unsafe.Add(mBase, _c_F_pow[17]))
											v362 = *(*float64)(unsafe.Add(mBase, _c_F_pow[18]))
											v365 = *(*float64)(unsafe.Add(mBase, _c_F_pow[19]))
											v368 = base.I64_reinterpret_f64(v341)
											v373 = base.I32_wrap_i64(v368) << (uint(int32(4)) % 32) & int32(2032)
											v376 = *(*float64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[20])))
											v379 = base.F64_add(base.F64_mul(base.F64_mul(v352, v352), base.F64_add(base.F64_mul(v351, v355), v358)), base.F64_add(base.F64_mul(v352, base.F64_add(base.F64_mul(v351, v362), v365)), base.F64_add(v376, v351)))
											v382 = *(*int64)(unsafe.Add(mBase, uint32(v373)+uint32(_c_F_pow[21])))
											v387 = v382 + (v368+base.I64_extend_i32_u(v201))<<(uint(int64(45))%64)
											if v326 != 0 {
												v445 = base.F64_reinterpret_i64(v387)
												v456 = base.F64_add(base.F64_mul(v445, v379), v445)
												return v456
											} else {
												if v368&int64(2147483648) != int64(0) {
													v404 = v387 + int64(4602678819172646912)
													v405 = base.F64_reinterpret_i64(v404)
													v406 = base.F64_mul(v405, v379)
													v407 = base.F64_add(v406, v405)
													if base.F64_lt(base.F64_abs(v407), float64(1)) == int32(0) {
														v432 = v407
													} else {
														if base.F64_lt(v407, float64(0)) != 0 {
															v420 = float64(-1)
														} else {
															v420 = float64(1)
														}
														v421 = base.F64_add(v407, v420)
														v428 = base.F64_sub(base.F64_add(v421, base.F64_add(base.F64_add(v406, base.F64_sub(v405, v407)), base.F64_add(v407, base.F64_sub(v420, v421)))), v420)
														if base.F64_eq(v428, float64(0)) != 0 {
															v431 = base.F64_reinterpret_i64(v404 & int64(-9223372036854775807-1))
														} else {
															v431 = v428
														}
														v432 = v431
													}
													v443 = base.F64_mul(v432, float64(2.2250738585072014e-308))
												} else {
													v398 = base.F64_reinterpret_i64(v387 + int64(-4544132024016830464))
													v443 = base.F64_mul(base.F64_add(base.F64_mul(v398, v379), v398), float64(5.486124068793689e+303))
												}
												return v443
											}
										} else {
											if int64(-1) < v299 {
												v324 = F___math_xflow(m, v201, float64(3.105036184601418e+231))
												mBase = m.M
												return v324
											} else {
												v321 = F___math_xflow(m, v201, float64(1.2882297539194267e-231))
												mBase = m.M
												return v321
											}
										}
									} else {
										if v201 != 0 {
											v313 = float64(-1)
										} else {
											v313 = float64(1)
										}
										return v313
									}
								}
							} else {
								v174 = float64(1)
								if v169 == int64(4607182418800017408) {
									v456 = v174
									return v456
								} else {
									if base.Ui32(v24) < base.Ui32(int32(958)) {
										v456 = v174
										return v456
									} else {
										if base.B2i32(v19 < int64(0)) == base.B2i32(base.Ui64(int64(4607182418800017408)) < base.Ui64(v169)) {
											v190 = F___math_xflow(m, int32(0), float64(1.2882297539194267e-231))
											mBase = m.M
											return v190
										} else {
											v186 = F___math_xflow(m, int32(0), float64(3.105036184601418e+231))
											mBase = m.M
											return v186
										}
									}
								}
							}
						} else {
							v157 = base.F64_sub(l0, l0)
							return base.F64_div(v157, v157)
						}
					}
				} else {
					v75 = base.F64_mul(l0, l0)
					if int64(-1) < v27 {
						v115 = v75
					} else {
						v87 = base.I32_wrap_i64(int64(base.Ui64(v19)>>(uint(int64(52))%64))) & int32(2047)
						if base.Ui32(v87) < base.Ui32(int32(1023)) {
							v110 = int32(0)
						} else {
							if base.Ui32(int32(1075)) < base.Ui32(v87) {
								v110 = int32(2)
							} else {
								v98 = int64(1) << (uint(base.I64_extend_i32_u(int32(1075)-v87)) % 64)
								if (v98+int64(-1))&v19 != int64(0) {
									v110 = int32(0)
								} else {
									if v98&v19 == int64(0) {
										v109 = int32(2)
									} else {
										v109 = int32(1)
									}
									v110 = v109
								}
							}
						}
						if v110 == int32(1) {
							v114 = base.F64_neg(v75)
						} else {
							v114 = v75
						}
						v115 = v114
					}
					if int64(-1) < v19 {
						v456 = v115
						return v456
					} else {
						return base.F64_div(float64(1), v115)
					}
				}
			} else {
				v45 = float64(1)
				if v27 == int64(4607182418800017408) {
					v456 = v45
					return v456
				} else {
					if v40 == int64(0) {
						v456 = v45
						return v456
					} else {
						v51 = v27 << (uint(int64(1)) % 64)
						if base.Ui64(int64(-9007199254740992)) < base.Ui64(v51) {
							return base.F64_add(l0, l1)
						} else {
							if base.Ui64(v40) < base.Ui64(int64(-9007199254740991)) {
								if v51 == int64(9214364837600034816) {
									v456 = v45
									return v456
								} else {
									if base.B2i32(base.Ui64(v51) < base.Ui64(int64(9214364837600034816)))^base.B2i32(v19 < int64(0)) != 0 {
										v67 = float64(0)
									} else {
										v67 = base.F64_mul(l1, l1)
									}
									return v67
								}
							} else {
								return base.F64_add(l0, l1)
							}
						}
					}
				}
			}
		}
	}
}
func F_prepend_alloc(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v87 int32
	_ = v87
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var __phi95 int32
	_ = __phi95
	var v101 int32
	_ = v101
	var __phi101 int32
	_ = __phi101
	var v106 int32
	_ = v106
	var v109 int32
	_ = v109
	var v115 int32
	_ = v115
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v135 int32
	_ = v135
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v150 int32
	_ = v150
	var v155 int32
	_ = v155
	var v171 int32
	_ = v171
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v229 int32
	_ = v229
	var v234 int32
	_ = v234
	var v238 int32
	_ = v238
	var v240 int32
	_ = v240
	var v256 int32
	_ = v256
	var v258 int32
	_ = v258
	var v261 int32
	_ = v261
	var v266 int32
	_ = v266
	var v269 int32
	_ = v269
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v287 int32
	_ = v287
	v11 = int32(-8)
	v13 = int32(15)
	v15 = l0 + (v11-l0)&v13
	*(*int32)(unsafe.Add(mBase, uint32(v15)+4)) = l2 | int32(3)
	v23 = l1 + (v11-l1)&v13
	v24 = v15 + l2
	v25 = v23 - v24
	v27 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[0]))
	if v23 != v27 {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	return v15 + int32(8)
L2:
	;
	v40 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[1]))
	if v23 != v40 {
		goto L4
	} else {
		goto L5
	}
L3:
	;
	v29 = int32(0)
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[0])) = v24
	v33 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[2]))
	v34 = v33 + v25
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[2])) = v34
	*(*int32)(unsafe.Add(mBase, uint32(v24)+4)) = v34 | int32(1)
	goto L1
L4:
	;
	v54 = *(*int32)(unsafe.Add(mBase, uint32(v23)+4))
	if v54&int32(3) != int32(1) {
		v173 = v25
		v174 = v54
		v177 = v23
		goto L6
	} else {
		goto L7
	}
L5:
	;
	v42 = int32(0)
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[1])) = v24
	v46 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[3]))
	v47 = v46 + v25
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[3])) = v47
	*(*int32)(unsafe.Add(mBase, uint32(v24)+4)) = v47 | int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v24+v47))) = v47
	goto L1
L6:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v177)+4)) = v174 & int32(-2)
	*(*int32)(unsafe.Add(mBase, uint32(v24+v173))) = v173
	*(*int32)(unsafe.Add(mBase, uint32(v24)+4)) = v173 | int32(1)
	if base.Ui32(int32(255)) < base.Ui32(v173) {
		goto L37
	} else {
		goto L38
	}
L7:
	;
	v60 = v54 & int32(-8)
	v61 = *(*int32)(unsafe.Add(mBase, uint32(v23)+12))
	if base.Ui32(int32(255)) < base.Ui32(v54) {
		goto L9
	} else {
		goto L10
	}
L8:
	;
	v171 = v23 + v60
	v172 = *(*int32)(unsafe.Add(mBase, uint32(v171)+4))
	v173 = v60 + v25
	v174 = v172
	v177 = v171
	goto L6
L9:
	;
	v77 = *(*int32)(unsafe.Add(mBase, uint32(v23)+24))
	if v61 == v23 {
		goto L14
	} else {
		goto L15
	}
L10:
	;
	v64 = *(*int32)(unsafe.Add(mBase, uint32(v23)+8))
	if v61 != v64 {
		goto L11
	} else {
		goto L12
	}
L11:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v61)+8)) = v64
	*(*int32)(unsafe.Add(mBase, uint32(v64)+12)) = v61
	goto L8
L12:
	;
	v66 = int32(0)
	v68 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[4]))
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[4])) = v68 & base.I32_rotl(int32(-2), int32(base.Ui32(v54)>>(uint(int32(3))%32)))
	goto L8
L13:
	;
	if v77 == int32(0) {
		goto L8
	} else {
		goto L25
	}
L14:
	;
	v82 = *(*int32)(unsafe.Add(mBase, uint32(v23)+20))
	if v82 == int32(0) {
		goto L18
	} else {
		goto L19
	}
L15:
	;
	v79 = *(*int32)(unsafe.Add(mBase, uint32(v23)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v79)+12)) = v61
	*(*int32)(unsafe.Add(mBase, uint32(v61)+8)) = v79
	v115 = v61
	goto L13
L16:
	;
	v115 = int32(0)
	goto L13
L17:
	;
	__phi95 = v92
	__phi101 = v93
	v95 = __phi95
	v101 = __phi101
	goto L21
L18:
	;
	v87 = *(*int32)(unsafe.Add(mBase, uint32(v23)+16))
	if v87 == int32(0) {
		goto L16
	} else {
		goto L20
	}
L19:
	;
	v92 = v82
	v93 = v23 + int32(20)
	goto L17
L20:
	;
	v92 = v87
	v93 = v23 + int32(16)
	goto L17
L21:
	;
	v106 = *(*int32)(unsafe.Add(mBase, uint32(v95)+20))
	if v106 != 0 {
		__phi95 = v106
		__phi101 = v95 + int32(20)
		v95 = __phi95
		v101 = __phi101
		goto L21
	} else {
		goto L23
	}
L22:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v101))) = int32(0)
	v115 = v95
	goto L13
L23:
	;
	v109 = *(*int32)(unsafe.Add(mBase, uint32(v95)+16))
	if v109 != 0 {
		__phi95 = v109
		__phi101 = v95 + int32(16)
		v95 = __phi95
		v101 = __phi101
		goto L21
	} else {
		goto L24
	}
L24:
	;
	goto L22
L25:
	;
	v125 = *(*int32)(unsafe.Add(mBase, uint32(v23)+28))
	v127 = v125 << (uint(int32(2)) % 32)
	v130 = *(*int32)(unsafe.Add(mBase, uint32(v127)+uint32(_c_F_prepend_alloc[5])))
	if v23 != v130 {
		goto L27
	} else {
		goto L28
	}
L26:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v115)+24)) = v77
	v150 = *(*int32)(unsafe.Add(mBase, uint32(v23)+16))
	if v150 == int32(0) {
		goto L34
	} else {
		goto L35
	}
L27:
	;
	v142 = *(*int32)(unsafe.Add(mBase, uint32(v77)+16))
	if v142 == v23 {
		goto L30
	} else {
		goto L31
	}
L28:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v127)+uint32(_c_F_prepend_alloc[5]))) = v115
	if v115 != 0 {
		goto L26
	} else {
		goto L29
	}
L29:
	;
	v133 = int32(0)
	v135 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[6]))
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[6])) = v135 & base.I32_rotl(int32(-2), v125)
	goto L8
L30:
	;
	v144 = int32(16)
	goto L32
L31:
	;
	v144 = int32(20)
	goto L32
L32:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v77+v144))) = v115
	if v115 == int32(0) {
		goto L8
	} else {
		goto L33
	}
L33:
	;
	goto L26
L34:
	;
	v155 = *(*int32)(unsafe.Add(mBase, uint32(v23)+20))
	if v155 == int32(0) {
		goto L8
	} else {
		goto L36
	}
L35:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v115)+16)) = v150
	*(*int32)(unsafe.Add(mBase, uint32(v150)+24)) = v115
	goto L34
L36:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v115)+20)) = v155
	*(*int32)(unsafe.Add(mBase, uint32(v155)+24)) = v115
	goto L8
L37:
	;
	if base.Ui32(int32(16777215)) < base.Ui32(v173) {
		v229 = int32(31)
		goto L42
	} else {
		goto L43
	}
L38:
	;
	v194 = v173 & int32(-8)
	v196 = v194 + int32(_a_F_prepend_alloc_0)
	v198 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[4]))
	v202 = int32(1) << (uint(int32(base.Ui32(v173)>>(uint(int32(3))%32))) % 32)
	if v198&v202 != 0 {
		goto L40
	} else {
		goto L41
	}
L39:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v208)+12)) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v194)+uint32(_c_F_prepend_alloc[7]))) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v24)+12)) = v196
	*(*int32)(unsafe.Add(mBase, uint32(v24)+8)) = v208
	goto L1
L40:
	;
	v207 = *(*int32)(unsafe.Add(mBase, uint32(v194)+uint32(_c_F_prepend_alloc[7])))
	v208 = v207
	goto L39
L41:
	;
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[4])) = v198 | v202
	v208 = v196
	goto L39
L42:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v24)+28)) = v229
	*(*int64)(unsafe.Add(mBase, uint32(v24)+16)) = int64(0)
	v234 = v229 << (uint(int32(2)) % 32)
	v238 = *(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[6]))
	v240 = int32(1) << (uint(v229) % 32)
	if v238&v240 != 0 {
		goto L44
	} else {
		goto L45
	}
L43:
	;
	v219 = base.I32_clz(int32(base.Ui32(v173) >> (uint(int32(8)) % 32)))
	v222 = int32(1)
	v229 = int32(base.Ui32(v173)>>(uint(int32(38)-v219)%32))&v222 - v219<<(uint(v222)%32) + int32(62)
	goto L42
L44:
	;
	if v229 == int32(31) {
		goto L46
	} else {
		goto L47
	}
L45:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v234)+uint32(_c_F_prepend_alloc[5]))) = v24
	*(*int32)(unsafe.Add(mBase, _c_F_prepend_alloc[6])) = v238 | v240
	*(*int32)(unsafe.Add(mBase, uint32(v24)+24)) = v234 + int32(_a_F_prepend_alloc_1)
	*(*int32)(unsafe.Add(mBase, uint32(v24)+8)) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v24)+12)) = v24
	goto L1
L46:
	;
	v256 = int32(0)
	goto L48
L47:
	;
	v256 = int32(25) - int32(base.Ui32(v229)>>(uint(int32(1))%32))
	goto L48
L48:
	;
	v258 = *(*int32)(unsafe.Add(mBase, uint32(v234)+uint32(_c_F_prepend_alloc[5])))
	v261 = v173 << (uint(v256) % 32)
	v266 = v258
	goto L50
L49:
	;
	v287 = *(*int32)(unsafe.Add(mBase, uint32(v266)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v287)+12)) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v266)+8)) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v24)+24)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v24)+12)) = v266
	*(*int32)(unsafe.Add(mBase, uint32(v24)+8)) = v287
	goto L1
L50:
	;
	v269 = *(*int32)(unsafe.Add(mBase, uint32(v266)+4))
	if v269&int32(-8) == v173 {
		goto L49
	} else {
		goto L52
	}
L51:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v281))) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v24)+24)) = v266
	*(*int32)(unsafe.Add(mBase, uint32(v24)+12)) = v24
	*(*int32)(unsafe.Add(mBase, uint32(v24)+8)) = v24
	goto L1
L52:
	;
	v281 = v266 + int32(base.Ui32(v261)>>(uint(int32(29))%32))&int32(4) + int32(16)
	v282 = *(*int32)(unsafe.Add(mBase, uint32(v281)))
	if v282 != 0 {
		v261 = v261 << (uint(int32(1)) % 32)
		v266 = v282
		goto L50
	} else {
		goto L53
	}
L53:
	;
	goto L51
}

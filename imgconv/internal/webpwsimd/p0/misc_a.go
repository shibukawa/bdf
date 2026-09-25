//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_AccumulateRGB(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v36 int32
	_ = v36
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v108 int32
	_ = v108
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v156 int32
	_ = v156
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v169 int32
	_ = v169
	var v171 int32
	_ = v171
	var v175 int32
	_ = v175
	var v178 int32
	_ = v178
	var v182 int32
	_ = v182
	var v185 int32
	_ = v185
	var v189 int32
	_ = v189
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v197 int32
	_ = v197
	var v200 int32
	_ = v200
	var v204 int32
	_ = v204
	var v210 int32
	_ = v210
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v216 int32
	_ = v216
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v243 int32
	_ = v243
	var v245 int32
	_ = v245
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
	var v257 int32
	_ = v257
	var v258 int32
	_ = v258
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
	var v269 int32
	_ = v269
	var v270 int32
	_ = v270
	var v272 int32
	_ = v272
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v286 int32
	_ = v286
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v300 int32
	_ = v300
	var v305 int32
	_ = v305
	var v308 int32
	_ = v308
	var v312 int32
	_ = v312
	var v318 int32
	_ = v318
	var v320 int32
	_ = v320
	var v322 int32
	_ = v322
	var v326 int32
	_ = v326
	var v327 int32
	_ = v327
	var v331 int32
	_ = v331
	var v332 int32
	_ = v332
	var v336 int32
	_ = v336
	var v341 int32
	_ = v341
	var v344 int32
	_ = v344
	var v348 int32
	_ = v348
	var v354 int32
	_ = v354
	v23 = int32(1)
	v24 = l6 >> (uint(v23) % 32)
	if v23 <= v24 {
		v36 = l4 + l3
		v44 = int32(0)
		v46 = l5
		v48 = v24
		for {
			v63 = m.G1
			v65 = v63 + int32(_a_F_AccumulateRGB_0)
			v67 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l3+v44))))
			v68 = int32(1)
			v71 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v67<<(uint(v68)%32)))))
			v73 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v44))))
			v77 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v73<<(uint(v68)%32)))))
			v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+l4+v44))))
			v84 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v80<<(uint(v68)%32)))))
			v87 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v36+v44))))
			v91 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v87<<(uint(v68)%32)))))
			v92 = v71 + v77 + v84 + v91
			v93 = int32(511)
			v94 = v92 & v93
			v96 = v63 + int32(_a_F_AccumulateRGB_1)
			v97 = int32(7)
			v99 = int32(4092)
			v101 = v96 + int32(base.Ui32(v92)>>(uint(v97)%32))&v99
			v102 = int32(4)
			v104 = *(*int32)(unsafe.Add(mBase, uint32(v101+v102)))
			v106 = int32(512)
			v108 = *(*int32)(unsafe.Add(mBase, uint32(v101)))
			v111 = int32(64)
			v114 = int32(base.Ui32(v94*v104+(v106-v94)*v108+v111) >> (uint(v97) % 32))
			*(*uint16)(unsafe.Add(mBase, uint32(v46))) = uint16(v114)
			v117 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+l3+v44))))
			v121 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v117<<(uint(v68)%32)))))
			v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v44))))
			v127 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v123<<(uint(v68)%32)))))
			v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+l4+v44))))
			v134 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v130<<(uint(v68)%32)))))
			v137 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v36+v44))))
			v141 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v137<<(uint(v68)%32)))))
			v142 = v121 + v127 + v134 + v141
			v144 = v142 & v93
			v149 = v96 + int32(base.Ui32(v142)>>(uint(v97)%32))&v99
			v152 = *(*int32)(unsafe.Add(mBase, uint32(v149+v102)))
			v156 = *(*int32)(unsafe.Add(mBase, uint32(v149)))
			v162 = int32(base.Ui32(v144*v152+(v106-v144)*v156+v111) >> (uint(v97) % 32))
			*(*uint16)(unsafe.Add(mBase, uint32(v46)+2)) = uint16(v162)
			v165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+l3+v44))))
			v169 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v165<<(uint(v68)%32)))))
			v171 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v44))))
			v175 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v171<<(uint(v68)%32)))))
			v178 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+l4+v44))))
			v182 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v178<<(uint(v68)%32)))))
			v185 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v36+v44))))
			v189 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v65+v185<<(uint(v68)%32)))))
			v190 = v169 + v175 + v182 + v189
			v192 = v190 & v93
			v197 = v96 + int32(base.Ui32(v190)>>(uint(v97)%32))&v99
			v200 = *(*int32)(unsafe.Add(mBase, uint32(v197+v102)))
			v204 = *(*int32)(unsafe.Add(mBase, uint32(v197)))
			v210 = int32(base.Ui32(v192*v200+(v106-v192)*v204+v111) >> (uint(v97) % 32))
			*(*uint16)(unsafe.Add(mBase, uint32(v46)+4)) = uint16(v210)
			v213 = v46 + int32(8)
			v214 = v44 + l3<<(uint(int32(1))%32)
			v216 = v48 + int32(-1)
			if v216 != 0 {
				v44 = v214
				v46 = v213
				v48 = v216
				continue
			} else {
				break
			}
			break
		}
		v220 = v214
		v222 = v213
	} else {
		v220 = int32(0)
		v222 = l5
	}
	if l6&int32(1) == int32(0) {
	} else {
		v243 = m.G1
		v245 = v243 + int32(_a_F_AccumulateRGB_0)
		v246 = l0 + v220
		v248 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v246+l4))))
		v249 = int32(1)
		v252 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v245+v248<<(uint(v249)%32)))))
		v253 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v246))))
		v257 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v245+v253<<(uint(v249)%32)))))
		v258 = v252 + v257
		v261 = int32(510)
		v262 = v258 << (uint(v249) % 32) & v261
		v264 = v243 + int32(_a_F_AccumulateRGB_1)
		v265 = int32(6)
		v267 = int32(2044)
		v269 = v264 + int32(base.Ui32(v258)>>(uint(v265)%32))&v267
		v270 = int32(4)
		v272 = *(*int32)(unsafe.Add(mBase, uint32(v269+v270)))
		v274 = int32(512)
		v276 = *(*int32)(unsafe.Add(mBase, uint32(v269)))
		v279 = int32(64)
		v281 = int32(7)
		v282 = int32(base.Ui32(v262*v272+(v274-v262)*v276+v279) >> (uint(v281) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(v222))) = uint16(v282)
		v284 = l1 + v220
		v286 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v284+l4))))
		v290 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v245+v286<<(uint(v249)%32)))))
		v291 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v284))))
		v295 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v245+v291<<(uint(v249)%32)))))
		v296 = v290 + v295
		v300 = v296 << (uint(v249) % 32) & v261
		v305 = v264 + int32(base.Ui32(v296)>>(uint(v265)%32))&v267
		v308 = *(*int32)(unsafe.Add(mBase, uint32(v305+v270)))
		v312 = *(*int32)(unsafe.Add(mBase, uint32(v305)))
		v318 = int32(base.Ui32(v300*v308+(v274-v300)*v312+v279) >> (uint(v281) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(v222)+2)) = uint16(v318)
		v320 = l2 + v220
		v322 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+l4))))
		v326 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v245+v322<<(uint(v249)%32)))))
		v327 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320))))
		v331 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v245+v327<<(uint(v249)%32)))))
		v332 = v326 + v331
		v336 = v332 << (uint(v249) % 32) & v261
		v341 = v264 + int32(base.Ui32(v332)>>(uint(v265)%32))&v267
		v344 = *(*int32)(unsafe.Add(mBase, uint32(v341+v270)))
		v348 = *(*int32)(unsafe.Add(mBase, uint32(v341)))
		v354 = int32(base.Ui32(v336*v344+(v274-v336)*v348+v279) >> (uint(v281) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(v222)+4)) = uint16(v354)
	}
	return
}
func F_AccumulateRGBA(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v96 int32
	_ = v96
	var v100 int32
	_ = v100
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v119 int32
	_ = v119
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
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
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v145 int32
	_ = v145
	var v149 int32
	_ = v149
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v168 int32
	_ = v168
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v185 int32
	_ = v185
	var v187 int32
	_ = v187
	var v188 int32
	_ = v188
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v204 int32
	_ = v204
	var v206 int32
	_ = v206
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v214 int32
	_ = v214
	var v216 int32
	_ = v216
	var v220 int32
	_ = v220
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v233 int32
	_ = v233
	var v237 int32
	_ = v237
	var v245 int32
	_ = v245
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
	var v257 int32
	_ = v257
	var v260 int32
	_ = v260
	var v262 int32
	_ = v262
	var v264 int32
	_ = v264
	var v267 int32
	_ = v267
	var v269 int32
	_ = v269
	var v272 int32
	_ = v272
	var v276 int32
	_ = v276
	var v278 int32
	_ = v278
	var v282 int32
	_ = v282
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v290 int32
	_ = v290
	var v295 int32
	_ = v295
	var v299 int32
	_ = v299
	var v302 int32
	_ = v302
	var v306 int32
	_ = v306
	var v311 int32
	_ = v311
	var v314 int32
	_ = v314
	var v318 int32
	_ = v318
	var v323 int32
	_ = v323
	var v326 int32
	_ = v326
	var v330 int32
	_ = v330
	var v332 int32
	_ = v332
	var v336 int32
	_ = v336
	var v339 int32
	_ = v339
	var v340 int32
	_ = v340
	var v344 int32
	_ = v344
	var v349 int32
	_ = v349
	var v353 int32
	_ = v353
	var v356 int32
	_ = v356
	var v362 int32
	_ = v362
	var v363 int32
	_ = v363
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v370 int32
	_ = v370
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v382 int32
	_ = v382
	var v385 int32
	_ = v385
	var v386 int32
	_ = v386
	var v391 int32
	_ = v391
	var v394 int32
	_ = v394
	var v398 int32
	_ = v398
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v409 int32
	_ = v409
	var v416 int32
	_ = v416
	var v419 int32
	_ = v419
	var v439 int32
	_ = v439
	var v441 int32
	_ = v441
	var v442 int32
	_ = v442
	var v443 int32
	_ = v443
	var v445 int32
	_ = v445
	var v448 int32
	_ = v448
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v453 int32
	_ = v453
	var v454 int32
	_ = v454
	var v457 int32
	_ = v457
	var v458 int32
	_ = v458
	var v462 int32
	_ = v462
	var v463 int32
	_ = v463
	var v466 int32
	_ = v466
	var v467 int32
	_ = v467
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v472 int32
	_ = v472
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v477 int32
	_ = v477
	var v479 int32
	_ = v479
	var v481 int32
	_ = v481
	var v484 int32
	_ = v484
	var v486 int32
	_ = v486
	var v488 int32
	_ = v488
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v497 int32
	_ = v497
	var v498 int32
	_ = v498
	var v502 int32
	_ = v502
	var v507 int32
	_ = v507
	var v510 int32
	_ = v510
	var v514 int32
	_ = v514
	var v519 int32
	_ = v519
	var v521 int32
	_ = v521
	var v525 int32
	_ = v525
	var v526 int32
	_ = v526
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v538 int32
	_ = v538
	var v540 int32
	_ = v540
	var v541 int32
	_ = v541
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v547 int32
	_ = v547
	var v549 int32
	_ = v549
	var v553 int32
	_ = v553
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
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v571 int32
	_ = v571
	var v573 int32
	_ = v573
	var v574 int32
	_ = v574
	var v576 int32
	_ = v576
	var v578 int32
	_ = v578
	var v580 int32
	_ = v580
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
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
	var v604 int32
	_ = v604
	var v609 int32
	_ = v609
	var v612 int32
	_ = v612
	var v616 int32
	_ = v616
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v627 int32
	_ = v627
	var v629 int32
	_ = v629
	var v633 int32
	_ = v633
	var v636 int32
	_ = v636
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v655 int32
	_ = v655
	var v656 int32
	_ = v656
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v666 int32
	_ = v666
	var v669 int32
	_ = v669
	var v673 int32
	_ = v673
	var v679 int32
	_ = v679
	v25 = int32(1)
	v26 = l6 >> (uint(v25) % 32)
	if v25 <= v26 {
		v30 = l2 + l4
		v31 = l1 + l4
		v32 = l0 + l4
		v42 = v26
		v43 = int32(0)
		for {
			v59 = l3 + l4 + v43
			v60 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59))))
			v61 = l3 + v43
			v62 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61))))
			v64 = int32(4)
			v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v61+v64))))
			v70 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v64))))
			v71 = v60 + v62 + v66 + v70
			if v71 == int32(1020) {
				v74 = m.G1
				v76 = v74 + int32(_a_F_AccumulateRGBA_0)
				v77 = l1 + v43
				v78 = int32(4)
				v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v77+v78))))
				v81 = int32(1)
				v84 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v80<<(uint(v81)%32)))))
				v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v77))))
				v89 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v85<<(uint(v81)%32)))))
				v91 = v31 + v43
				v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
				v96 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v92<<(uint(v81)%32)))))
				v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91+v78))))
				v104 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v100<<(uint(v81)%32)))))
				v105 = v84 + v89 + v96 + v104
				v106 = int32(511)
				v107 = v105 & v106
				v109 = v74 + int32(_a_F_AccumulateRGBA_1)
				v110 = int32(7)
				v112 = int32(4092)
				v114 = v109 + int32(base.Ui32(v105)>>(uint(v110)%32))&v112
				v117 = *(*int32)(unsafe.Add(mBase, uint32(v114+v78)))
				v119 = int32(512)
				v121 = *(*int32)(unsafe.Add(mBase, uint32(v114)))
				v124 = int32(64)
				v126 = l0 + v43
				v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v126+v78))))
				v133 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v129<<(uint(v81)%32)))))
				v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v126))))
				v138 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v134<<(uint(v81)%32)))))
				v140 = v32 + v43
				v141 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v140))))
				v145 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v141<<(uint(v81)%32)))))
				v149 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v140+v78))))
				v153 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v149<<(uint(v81)%32)))))
				v154 = v133 + v138 + v145 + v153
				v156 = v154 & v106
				v161 = v109 + int32(base.Ui32(v154)>>(uint(v110)%32))&v112
				v164 = *(*int32)(unsafe.Add(mBase, uint32(v161+v78)))
				v168 = *(*int32)(unsafe.Add(mBase, uint32(v161)))
				v173 = l2 + v43
				v176 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v173+v78))))
				v180 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v176<<(uint(v81)%32)))))
				v181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v173))))
				v185 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v181<<(uint(v81)%32)))))
				v187 = v30 + v43
				v188 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v187))))
				v192 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v188<<(uint(v81)%32)))))
				v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v187+v78))))
				v200 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v196<<(uint(v81)%32)))))
				v201 = v180 + v185 + v192 + v200
				v362 = v201
				v363 = int32(base.Ui32(v201) >> (uint(int32(9)) % 32))
				v365 = v107*v117 + (v119-v107)*v121 + v124
				v366 = v156*v164 + (v119-v156)*v168 + v124
			} else {
				if v71 != 0 {
					v204 = m.G1
					v206 = v204 + int32(_a_F_AccumulateRGBA_0)
					v207 = l1 + v43
					v208 = int32(4)
					v210 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v207+v208))))
					v211 = int32(1)
					v214 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v210<<(uint(v211)%32)))))
					v216 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v207))))
					v220 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v216<<(uint(v211)%32)))))
					v223 = v31 + v43
					v224 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v223))))
					v228 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v224<<(uint(v211)%32)))))
					v233 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v223+v208))))
					v237 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v233<<(uint(v211)%32)))))
					v245 = *(*int32)(unsafe.Add(mBase, uint32(v204+int32(1024)+v71<<(uint(int32(2))%32))))
					v246 = (v214*v66 + v220*v62 + v228*v60 + v237*v70) * v245
					v247 = int32(17)
					v249 = int32(511)
					v250 = int32(base.Ui32(v246)>>(uint(v247)%32)) & v249
					v252 = v204 + int32(_a_F_AccumulateRGBA_1)
					v253 = int32(24)
					v255 = int32(252)
					v257 = v252 + int32(base.Ui32(v246)>>(uint(v253)%32))&v255
					v260 = *(*int32)(unsafe.Add(mBase, uint32(v257+v208)))
					v262 = int32(512)
					v264 = *(*int32)(unsafe.Add(mBase, uint32(v257)))
					v267 = int32(64)
					v269 = l0 + v43
					v272 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v269+v208))))
					v276 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v272<<(uint(v211)%32)))))
					v278 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v269))))
					v282 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v278<<(uint(v211)%32)))))
					v285 = v32 + v43
					v286 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v285))))
					v290 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v286<<(uint(v211)%32)))))
					v295 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v285+v208))))
					v299 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v295<<(uint(v211)%32)))))
					v302 = (v276*v66 + v282*v62 + v290*v60 + v299*v70) * v245
					v306 = int32(base.Ui32(v302)>>(uint(v247)%32)) & v249
					v311 = v252 + int32(base.Ui32(v302)>>(uint(v253)%32))&v255
					v314 = *(*int32)(unsafe.Add(mBase, uint32(v311+v208)))
					v318 = *(*int32)(unsafe.Add(mBase, uint32(v311)))
					v323 = l2 + v43
					v326 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323+v208))))
					v330 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v326<<(uint(v211)%32)))))
					v332 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v323))))
					v336 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v332<<(uint(v211)%32)))))
					v339 = v30 + v43
					v340 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v339))))
					v344 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v340<<(uint(v211)%32)))))
					v349 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v339+v208))))
					v353 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v206+v349<<(uint(v211)%32)))))
					v356 = (v330*v66 + v336*v62 + v344*v60 + v353*v70) * v245
					v362 = int32(base.Ui32(v356) >> (uint(v247) % 32))
					v363 = int32(base.Ui32(v356) >> (uint(int32(26)) % 32))
					v365 = v250*v260 + (v262-v250)*v264 + v267
					v366 = v306*v314 + (v262-v306)*v318 + v267
				} else {
					v74 = m.G1
					v76 = v74 + int32(_a_F_AccumulateRGBA_0)
					v77 = l1 + v43
					v78 = int32(4)
					v80 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v77+v78))))
					v81 = int32(1)
					v84 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v80<<(uint(v81)%32)))))
					v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v77))))
					v89 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v85<<(uint(v81)%32)))))
					v91 = v31 + v43
					v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
					v96 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v92<<(uint(v81)%32)))))
					v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91+v78))))
					v104 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v100<<(uint(v81)%32)))))
					v105 = v84 + v89 + v96 + v104
					v106 = int32(511)
					v107 = v105 & v106
					v109 = v74 + int32(_a_F_AccumulateRGBA_1)
					v110 = int32(7)
					v112 = int32(4092)
					v114 = v109 + int32(base.Ui32(v105)>>(uint(v110)%32))&v112
					v117 = *(*int32)(unsafe.Add(mBase, uint32(v114+v78)))
					v119 = int32(512)
					v121 = *(*int32)(unsafe.Add(mBase, uint32(v114)))
					v124 = int32(64)
					v126 = l0 + v43
					v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v126+v78))))
					v133 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v129<<(uint(v81)%32)))))
					v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v126))))
					v138 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v134<<(uint(v81)%32)))))
					v140 = v32 + v43
					v141 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v140))))
					v145 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v141<<(uint(v81)%32)))))
					v149 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v140+v78))))
					v153 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v149<<(uint(v81)%32)))))
					v154 = v133 + v138 + v145 + v153
					v156 = v154 & v106
					v161 = v109 + int32(base.Ui32(v154)>>(uint(v110)%32))&v112
					v164 = *(*int32)(unsafe.Add(mBase, uint32(v161+v78)))
					v168 = *(*int32)(unsafe.Add(mBase, uint32(v161)))
					v173 = l2 + v43
					v176 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v173+v78))))
					v180 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v176<<(uint(v81)%32)))))
					v181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v173))))
					v185 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v181<<(uint(v81)%32)))))
					v187 = v30 + v43
					v188 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v187))))
					v192 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v188<<(uint(v81)%32)))))
					v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v187+v78))))
					v200 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v76+v196<<(uint(v81)%32)))))
					v201 = v180 + v185 + v192 + v200
					v362 = v201
					v363 = int32(base.Ui32(v201) >> (uint(int32(9)) % 32))
					v365 = v107*v117 + (v119-v107)*v121 + v124
					v366 = v156*v164 + (v119-v156)*v168 + v124
				}
			}
			v370 = l5 + v43
			*(*uint16)(unsafe.Add(mBase, uint32(v370+int32(6)))) = uint16(v71)
			v374 = int32(2)
			v376 = int32(7)
			v377 = int32(base.Ui32(v365) >> (uint(v376) % 32))
			*(*uint16)(unsafe.Add(mBase, uint32(v370+v374))) = uint16(v377)
			v380 = int32(base.Ui32(v366) >> (uint(v376) % 32))
			*(*uint16)(unsafe.Add(mBase, uint32(v370))) = uint16(v380)
			v382 = int32(4)
			v385 = v362 & int32(511)
			v386 = m.G1
			v391 = v386 + int32(_a_F_AccumulateRGBA_1) + v363<<(uint(v374)%32)
			v394 = *(*int32)(unsafe.Add(mBase, uint32(v391+v382)))
			v398 = *(*int32)(unsafe.Add(mBase, uint32(v391)))
			v404 = int32(base.Ui32(v385*v394+(int32(512)-v385)*v398+int32(64)) >> (uint(v376) % 32))
			*(*uint16)(unsafe.Add(mBase, uint32(v370+v382))) = uint16(v404)
			v407 = v43 + int32(8)
			v409 = v42 + int32(-1)
			if v409 != 0 {
				v42 = v409
				v43 = v407
				continue
			} else {
				break
			}
			break
		}
		v416 = l5 + v407
		v419 = v407
	} else {
		v416 = l5
		v419 = int32(0)
	}
	if l6&int32(1) == int32(0) {
	} else {
		v439 = l3 + v419
		v441 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v439+l4))))
		v442 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v439))))
		v443 = v441 + v442
		v445 = v443 << (uint(int32(1)) % 32)
		if v443 == int32(510) {
			v448 = m.G1
			v450 = v448 + int32(_a_F_AccumulateRGBA_0)
			v451 = l1 + v419
			v453 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v451+l4))))
			v454 = int32(1)
			v457 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v453<<(uint(v454)%32)))))
			v458 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v451))))
			v462 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v458<<(uint(v454)%32)))))
			v463 = v457 + v462
			v466 = int32(510)
			v467 = v463 << (uint(v454) % 32) & v466
			v469 = v448 + int32(_a_F_AccumulateRGBA_1)
			v470 = int32(6)
			v472 = int32(2044)
			v474 = v469 + int32(base.Ui32(v463)>>(uint(v470)%32))&v472
			v475 = int32(4)
			v477 = *(*int32)(unsafe.Add(mBase, uint32(v474+v475)))
			v479 = int32(512)
			v481 = *(*int32)(unsafe.Add(mBase, uint32(v474)))
			v484 = int32(64)
			v486 = l0 + v419
			v488 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v486+l4))))
			v492 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v488<<(uint(v454)%32)))))
			v493 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v486))))
			v497 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v493<<(uint(v454)%32)))))
			v498 = v492 + v497
			v502 = v498 << (uint(v454) % 32) & v466
			v507 = v469 + int32(base.Ui32(v498)>>(uint(v470)%32))&v472
			v510 = *(*int32)(unsafe.Add(mBase, uint32(v507+v475)))
			v514 = *(*int32)(unsafe.Add(mBase, uint32(v507)))
			v519 = l2 + v419
			v521 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v519+l4))))
			v525 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v521<<(uint(v454)%32)))))
			v526 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v519))))
			v530 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v526<<(uint(v454)%32)))))
			v531 = v525 + v530
			v646 = v531 << (uint(v454) % 32) & v466
			v647 = int32(base.Ui32(v531) >> (uint(int32(8)) % 32))
			v649 = v467*v477 + (v479-v467)*v481 + v484
			v650 = v502*v510 + (v479-v502)*v514 + v484
		} else {
			if v443 != 0 {
				v538 = m.G1
				v540 = v538 + int32(_a_F_AccumulateRGBA_0)
				v541 = l1 + v419
				v543 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v541+l4))))
				v544 = int32(1)
				v547 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v540+v543<<(uint(v544)%32)))))
				v549 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v541))))
				v553 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v540+v549<<(uint(v544)%32)))))
				v561 = *(*int32)(unsafe.Add(mBase, uint32(v538+int32(1024)+v445<<(uint(int32(2))%32))))
				v562 = (v547*v441 + v553*v442) * v561
				v563 = int32(16)
				v565 = int32(511)
				v566 = int32(base.Ui32(v562)>>(uint(v563)%32)) & v565
				v568 = v538 + int32(_a_F_AccumulateRGBA_1)
				v569 = int32(23)
				v571 = int32(252)
				v573 = v568 + int32(base.Ui32(v562)>>(uint(v569)%32))&v571
				v574 = int32(4)
				v576 = *(*int32)(unsafe.Add(mBase, uint32(v573+v574)))
				v578 = int32(512)
				v580 = *(*int32)(unsafe.Add(mBase, uint32(v573)))
				v583 = int32(64)
				v585 = l0 + v419
				v587 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v585+l4))))
				v591 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v540+v587<<(uint(v544)%32)))))
				v593 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v585))))
				v597 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v540+v593<<(uint(v544)%32)))))
				v600 = (v591*v441 + v597*v442) * v561
				v604 = int32(base.Ui32(v600)>>(uint(v563)%32)) & v565
				v609 = v568 + int32(base.Ui32(v600)>>(uint(v569)%32))&v571
				v612 = *(*int32)(unsafe.Add(mBase, uint32(v609+v574)))
				v616 = *(*int32)(unsafe.Add(mBase, uint32(v609)))
				v621 = l2 + v419
				v623 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v621+l4))))
				v627 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v540+v623<<(uint(v544)%32)))))
				v629 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v621))))
				v633 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v540+v629<<(uint(v544)%32)))))
				v636 = (v627*v441 + v633*v442) * v561
				v646 = int32(base.Ui32(v636)>>(uint(v563)%32)) & v565
				v647 = int32(base.Ui32(v636)>>(uint(int32(25))%32)) & int32(63)
				v649 = v566*v576 + (v578-v566)*v580 + v583
				v650 = v604*v612 + (v578-v604)*v616 + v583
			} else {
				v448 = m.G1
				v450 = v448 + int32(_a_F_AccumulateRGBA_0)
				v451 = l1 + v419
				v453 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v451+l4))))
				v454 = int32(1)
				v457 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v453<<(uint(v454)%32)))))
				v458 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v451))))
				v462 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v458<<(uint(v454)%32)))))
				v463 = v457 + v462
				v466 = int32(510)
				v467 = v463 << (uint(v454) % 32) & v466
				v469 = v448 + int32(_a_F_AccumulateRGBA_1)
				v470 = int32(6)
				v472 = int32(2044)
				v474 = v469 + int32(base.Ui32(v463)>>(uint(v470)%32))&v472
				v475 = int32(4)
				v477 = *(*int32)(unsafe.Add(mBase, uint32(v474+v475)))
				v479 = int32(512)
				v481 = *(*int32)(unsafe.Add(mBase, uint32(v474)))
				v484 = int32(64)
				v486 = l0 + v419
				v488 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v486+l4))))
				v492 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v488<<(uint(v454)%32)))))
				v493 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v486))))
				v497 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v493<<(uint(v454)%32)))))
				v498 = v492 + v497
				v502 = v498 << (uint(v454) % 32) & v466
				v507 = v469 + int32(base.Ui32(v498)>>(uint(v470)%32))&v472
				v510 = *(*int32)(unsafe.Add(mBase, uint32(v507+v475)))
				v514 = *(*int32)(unsafe.Add(mBase, uint32(v507)))
				v519 = l2 + v419
				v521 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v519+l4))))
				v525 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v521<<(uint(v454)%32)))))
				v526 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v519))))
				v530 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v450+v526<<(uint(v454)%32)))))
				v531 = v525 + v530
				v646 = v531 << (uint(v454) % 32) & v466
				v647 = int32(base.Ui32(v531) >> (uint(int32(8)) % 32))
				v649 = v467*v477 + (v479-v467)*v481 + v484
				v650 = v502*v510 + (v479-v502)*v514 + v484
			}
		}
		*(*uint16)(unsafe.Add(mBase, uint32(v416)+6)) = uint16(v445)
		v655 = int32(7)
		v656 = int32(base.Ui32(v649) >> (uint(v655) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(v416)+2)) = uint16(v656)
		v659 = int32(base.Ui32(v650) >> (uint(v655) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(v416))) = uint16(v659)
		v661 = m.G1
		v666 = v661 + int32(_a_F_AccumulateRGBA_1) + v647<<(uint(int32(2))%32)
		v669 = *(*int32)(unsafe.Add(mBase, uint32(v666+int32(4))))
		v673 = *(*int32)(unsafe.Add(mBase, uint32(v666)))
		v679 = int32(base.Ui32(v646*v669+(int32(512)-v646)*v673+int32(64)) >> (uint(v655) % 32))
		*(*uint16)(unsafe.Add(mBase, uint32(v416)+4)) = uint16(v679)
	}
	return
}
func F_AllocateTransformBuffer(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v10 int64
	_ = v10
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v25 int64
	_ = v25
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v31 int64
	_ = v31
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v44 int32
	_ = v44
	var v45 int64
	_ = v45
	var v46 int64
	_ = v46
	var v50 int64
	_ = v50
	var v51 int32
	_ = v51
	var v54 int64
	_ = v54
	var v57 int64
	_ = v57
	var v59 int32
	_ = v59
	var v66 int64
	_ = v66
	var v67 int32
	_ = v67
	var v68 int64
	_ = v68
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v92 int32
	_ = v92
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	v10 = base.I64_extend_i32_s(l2) * base.I64_extend_i32_s(l1)
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
	if v11 == int32(0) {
		v24 = int32(0)
		v25 = int64(0)
		v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
		if v27 == v24 {
			v44 = v24
			v45 = v25
			v46 = v25
		} else {
			v30 = v24
			v31 = v25
			v33 = int32(3)
			v35 = int32(2)
			v44 = v30
			v45 = v31
			v46 = base.I64_extend_i32_u(int32(base.Ui32(l2+v33)>>(uint(v35)%32))) * base.I64_extend_i32_u(int32(base.Ui32(l1+v33)>>(uint(v35)%32)))
		}
	} else {
		v15 = l1 << (uint(int32(1)) % 32)
		v18 = int32(2)
		v22 = v15 + int32(base.Ui32(v15+int32(3))>>(uint(v18)%32)) + v18
		v30 = v22
		v31 = base.I64_extend_i32_u(v22)
		v33 = int32(3)
		v35 = int32(2)
		v44 = v30
		v45 = v31
		v46 = base.I64_extend_i32_u(int32(base.Ui32(l2+v33)>>(uint(v35)%32))) * base.I64_extend_i32_u(int32(base.Ui32(l1+v33)>>(uint(v35)%32)))
	}
	v50 = v10 + v45 + v46 + int64(16)
	v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if v51 == int32(0) {
		F_free(m, v51)
		mBase = m.M
		v57 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0)+24)) = v57
		v59 = int32(4)
		if v50 == v57 {
			v78 = F_malloc(m, base.I32_wrap_i64(v50)*v59)
			mBase = m.M
			v80 = v78
		} else {
			v66 = base.I64_div_u_s(int64(2147418112), v50)
			v67 = int32(0)
			v68 = base.I64_extend_i32_u(v59)
			if base.Ui64(int64(4294967295)) < base.Ui64(v68*v50) {
				v80 = v67
			} else {
				if base.Ui64(v66) < base.Ui64(v68) {
					v80 = v67
				} else {
					v78 = F_malloc(m, base.I32_wrap_i64(v50)*v59)
					mBase = m.M
					v80 = v78
				}
			}
		}
		if v80 != 0 {
			*(*uint32)(unsafe.Add(mBase, uint32(l0)+28)) = uint32(v50)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v80
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = int32(0)
			v92 = v80
			*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = l1
			*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v92
			v96 = int32(2)
			v99 = int32(31)
			v101 = int32(-32)
			v102 = (v92 + base.I32_wrap_i64(v10)<<(uint(v96)%32) + v99) & v101
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v102
			*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = (v102 + v44<<(uint(v96)%32) + v99) & v101
			return int32(1)
		} else {
			v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v84 = *(*int32)(unsafe.Add(mBase, uint32(v82)+92))
			if v84 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(v82)+92)) = int32(1)
			}
			return int32(0)
		}
	} else {
		v54 = int64(*(*uint32)(unsafe.Add(mBase, uint32(l0)+28)))
		if base.Ui64(v50) <= base.Ui64(v54) {
			v92 = v51
			*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = l1
			*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v92
			v96 = int32(2)
			v99 = int32(31)
			v101 = int32(-32)
			v102 = (v92 + base.I32_wrap_i64(v10)<<(uint(v96)%32) + v99) & v101
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v102
			*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = (v102 + v44<<(uint(v96)%32) + v99) & v101
			return int32(1)
		} else {
			F_free(m, v51)
			mBase = m.M
			v57 = int64(0)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+24)) = v57
			v59 = int32(4)
			if v50 == v57 {
				v78 = F_malloc(m, base.I32_wrap_i64(v50)*v59)
				mBase = m.M
				v80 = v78
			} else {
				v66 = base.I64_div_u_s(int64(2147418112), v50)
				v67 = int32(0)
				v68 = base.I64_extend_i32_u(v59)
				if base.Ui64(int64(4294967295)) < base.Ui64(v68*v50) {
					v80 = v67
				} else {
					if base.Ui64(v66) < base.Ui64(v68) {
						v80 = v67
					} else {
						v78 = F_malloc(m, base.I32_wrap_i64(v50)*v59)
						mBase = m.M
						v80 = v78
					}
				}
			}
			if v80 != 0 {
				*(*uint32)(unsafe.Add(mBase, uint32(l0)+28)) = uint32(v50)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v80
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = int32(0)
				v92 = v80
				*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = l1
				*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v92
				v96 = int32(2)
				v99 = int32(31)
				v101 = int32(-32)
				v102 = (v92 + base.I32_wrap_i64(v10)<<(uint(v96)%32) + v99) & v101
				*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v102
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = (v102 + v44<<(uint(v96)%32) + v99) & v101
				return int32(1)
			} else {
				v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v84 = *(*int32)(unsafe.Add(mBase, uint32(v82)+92))
				if v84 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v82)+92)) = int32(1)
				}
				return int32(0)
			}
		}
	}
}
func F_ApplyAlphaMultiply_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v21 int32
	_ = v21
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v78 int32
	_ = v78
	if l3 < int32(1) {
	} else {
		if l2 < int32(1) {
		} else {
			v16 = int32(0)
			if l1 != 0 {
				v21 = v16
			} else {
				v21 = int32(3)
			}
			v26 = l3
			v28 = l0 + base.B2i32(l1 != v16)
			v29 = l0 + v21
			for {
				v36 = v29
				v37 = v28
				v44 = l2
				for {
					v47 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v36))))
					if v47 == int32(255) {
					} else {
						v51 = v47 * int32(_a_F_ApplyAlphaMultiply_C_0)
						v52 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v37))))
						v54 = int32(23)
						v55 = int32(base.Ui32(v51*v52) >> (uint(v54) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v37))) = uint8(v55)
						v58 = v37 + int32(1)
						v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58))))
						v62 = int32(base.Ui32(v51*v59) >> (uint(v54) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v58))) = uint8(v62)
						v65 = v37 + int32(2)
						v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v65))))
						v69 = int32(base.Ui32(v51*v66) >> (uint(v54) % 32))
						*(*uint8)(unsafe.Add(mBase, uint32(v65))) = uint8(v69)
					}
					v73 = int32(4)
					v78 = v44 + int32(-1)
					if v78 != 0 {
						v36 = v36 + v73
						v37 = v37 + v73
						v44 = v78
						continue
					} else {
						break
					}
					break
				}
				if int32(1) < v26 {
					v26 = v26 + int32(-1)
					v28 = v28 + l4
					v29 = v29 + l4
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
func F_abort(m *base.Module) {
	base.Wasm_trap_unreachable()
	for {
	}
}

//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_UpdateChroma(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
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
	var v98 int32
	_ = v98
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
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
	var v133 int32
	_ = v133
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v156 int32
	_ = v156
	if l4 < int32(13) {
		v25 = int32(2)
	} else {
		v25 = int32(14) - l4
	}
	v26 = v25 + l4
	v27 = int32(1)
	if v27 < l3 {
		v30 = l3
	} else {
		v30 = v27
	}
	v34 = l3 << (uint(int32(2)) % 32)
	v37 = l3 << (uint(int32(3)) % 32)
	v44 = l2
	v45 = int32(0)
	v48 = v30
	for {
		v61 = l1 + v45
		v62 = int32(2)
		v64 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v61+v62))))
		v65 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v61))))
		v66 = l0 + v45
		v69 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v66+v62))))
		v70 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v66))))
		v71 = F_SharpYuvGammaToLinear(m, v70, v26, l5)
		mBase = m.M
		v72 = F_SharpYuvGammaToLinear(m, v69, v26, l5)
		mBase = m.M
		v74 = F_SharpYuvGammaToLinear(m, v65, v26, l5)
		mBase = m.M
		v76 = F_SharpYuvGammaToLinear(m, v64, v26, l5)
		mBase = m.M
		v82 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v71+v72+v74+v76+v62)>>(uint(v62)%32)), v26, l5)
		mBase = m.M
		v83 = l1 + v34 + v45
		v86 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v83+v62))))
		v87 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v83))))
		v88 = l0 + v34 + v45
		v91 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v88+v62))))
		v92 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v88))))
		v93 = F_SharpYuvGammaToLinear(m, v92, v26, l5)
		mBase = m.M
		v94 = F_SharpYuvGammaToLinear(m, v91, v26, l5)
		mBase = m.M
		v96 = F_SharpYuvGammaToLinear(m, v87, v26, l5)
		mBase = m.M
		v98 = F_SharpYuvGammaToLinear(m, v86, v26, l5)
		mBase = m.M
		v104 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v93+v94+v96+v98+v62)>>(uint(v62)%32)), v26, l5)
		mBase = m.M
		v105 = l1 + v37 + v45
		v108 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v105+v62))))
		v109 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v105))))
		v110 = l0 + v37 + v45
		v113 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v110+v62))))
		v121 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v110))))
		v122 = F_SharpYuvGammaToLinear(m, v121, v26, l5)
		mBase = m.M
		v123 = F_SharpYuvGammaToLinear(m, v113, v26, l5)
		mBase = m.M
		v125 = F_SharpYuvGammaToLinear(m, v109, v26, l5)
		mBase = m.M
		v127 = F_SharpYuvGammaToLinear(m, v108, v26, l5)
		mBase = m.M
		v133 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v122+v123+v125+v127+v62)>>(uint(v62)%32)), v26, l5)
		mBase = m.M
		v142 = base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v82)*int64(13933)+base.I64_extend_i32_u(v104)*int64(46871)+base.I64_extend_i32_u(v133)*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64)))
		v143 = v82 - v142
		*(*uint16)(unsafe.Add(mBase, uint32(v44))) = uint16(v143)
		v146 = v104 - v142
		*(*uint16)(unsafe.Add(mBase, uint32(v44+l3<<(uint(int32(1))%32)))) = uint16(v146)
		v149 = v133 - v142
		*(*uint16)(unsafe.Add(mBase, uint32(v44+v34))) = uint16(v149)
		v156 = v48 + int32(-1)
		if v156 != 0 {
			v44 = v44 + v62
			v45 = v45 + int32(4)
			v48 = v156
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_UpsampleArgbLinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v140 int32
	_ = v140
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v211 int32
	_ = v211
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v230 int32
	_ = v230
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v259 int32
	_ = v259
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v275 int32
	_ = v275
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v285 int32
	_ = v285
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v294 int32
	_ = v294
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v313 int32
	_ = v313
	var v316 int32
	_ = v316
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v337 int32
	_ = v337
	var v340 int32
	_ = v340
	var v345 int32
	_ = v345
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	var v353 int32
	_ = v353
	var v354 int32
	_ = v354
	var v357 int32
	_ = v357
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v368 int32
	_ = v368
	var v375 int32
	_ = v375
	var v378 int32
	_ = v378
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v390 int32
	_ = v390
	var v392 int32
	_ = v392
	var v399 int32
	_ = v399
	var v402 int32
	_ = v402
	var v404 int32
	_ = v404
	var v408 int32
	_ = v408
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v433 int32
	_ = v433
	var v440 int32
	_ = v440
	var v442 int32
	_ = v442
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v456 int32
	_ = v456
	var v459 int32
	_ = v459
	var v461 int32
	_ = v461
	var v467 int32
	_ = v467
	var v471 int32
	_ = v471
	var v472 int32
	_ = v472
	var v476 int32
	_ = v476
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v521 int32
	_ = v521
	var v523 int32
	_ = v523
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v537 int32
	_ = v537
	var v539 int32
	_ = v539
	var v546 int32
	_ = v546
	var v549 int32
	_ = v549
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v570 int32
	_ = v570
	var v573 int32
	_ = v573
	var v575 int32
	_ = v575
	var v579 int32
	_ = v579
	var v586 int32
	_ = v586
	var v588 int32
	_ = v588
	var v595 int32
	_ = v595
	var v598 int32
	_ = v598
	var v604 int32
	_ = v604
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v620 int32
	_ = v620
	var v623 int32
	_ = v623
	var v634 int32
	_ = v634
	var v641 int32
	_ = v641
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v680 int32
	_ = v680
	var v683 int32
	_ = v683
	var v684 int32
	_ = v684
	var v689 int32
	_ = v689
	var v693 int32
	_ = v693
	var v698 int32
	_ = v698
	var v700 int32
	_ = v700
	var v707 int32
	_ = v707
	var v710 int32
	_ = v710
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v729 int32
	_ = v729
	var v732 int32
	_ = v732
	var v736 int32
	_ = v736
	var v743 int32
	_ = v743
	var v745 int32
	_ = v745
	var v752 int32
	_ = v752
	var v755 int32
	_ = v755
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v763 int32
	_ = v763
	var v766 int32
	_ = v766
	var v767 int32
	_ = v767
	var v772 int32
	_ = v772
	var v776 int32
	_ = v776
	var v781 int32
	_ = v781
	var v783 int32
	_ = v783
	var v790 int32
	_ = v790
	var v793 int32
	_ = v793
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v803 int32
	_ = v803
	var v805 int32
	_ = v805
	var v812 int32
	_ = v812
	var v815 int32
	_ = v815
	var v819 int32
	_ = v819
	var v826 int32
	_ = v826
	var v828 int32
	_ = v828
	var v835 int32
	_ = v835
	var v838 int32
	_ = v838
	v30 = int32(255)
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v30)
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v33 = int32(16)
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v36 = v32<<(uint(v33)%32) | v35
	v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v43 = v39<<(uint(v33)%32) | v42
	v46 = v36*int32(3) + v43 + int32(131074)
	v50 = int32(base.Ui32(v46)>>(uint(int32(2))%32)) & v30
	v53 = int32(8)
	v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v59 = int32(base.Ui32(v55*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v53) % 32))
	v60 = int32(base.Ui32(v50*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v53)%32)) + v59
	v62 = v60 + int32(-17685)
	if base.Ui32(v60) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
		v69 = int32(0)
	} else {
		v69 = v30
	}
	if base.Ui32(v62) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
		v72 = int32(base.Ui32(v62) >> (uint(int32(6)) % 32))
	} else {
		v72 = v69
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+3)) = uint8(v72)
	v76 = int32(255)
	v77 = int32(base.Ui32(v46)>>(uint(int32(18))%32)) & v76
	v82 = int32(base.Ui32(v77*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32)) + v59
	v84 = v82 + int32(-14234)
	if base.Ui32(v82) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
		v91 = int32(0)
	} else {
		v91 = v76
	}
	if base.Ui32(v84) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
		v94 = int32(base.Ui32(v84) >> (uint(int32(6)) % 32))
	} else {
		v94 = v91
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v94)
	v98 = int32(8)
	v105 = v59 - (int32(base.Ui32(v77*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v98)%32)) + int32(base.Ui32(v50*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v98)%32)))
	v107 = v105 + int32(_a_F_UpsampleArgbLinePair_C_8)
	if v105 < int32(-8708) {
		v114 = int32(0)
	} else {
		v114 = int32(255)
	}
	if base.Ui32(v107) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
		v117 = int32(base.Ui32(v107) >> (uint(int32(6)) % 32))
	} else {
		v117 = v114
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v117)
	v120 = l8 + int32(-1)
	v122 = v120 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v125 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v125)
		v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v130 = int32(8)
		v131 = int32(base.Ui32(v127*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v130) % 32))
		v136 = v36 + v43*int32(3) + int32(131074)
		v140 = int32(base.Ui32(v136)>>(uint(int32(2))%32)) & v125
		v145 = v131 + int32(base.Ui32(v140*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v130)%32))
		v147 = v145 + int32(-17685)
		if base.Ui32(v145) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
			v154 = int32(0)
		} else {
			v154 = v125
		}
		if base.Ui32(v147) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
			v157 = int32(base.Ui32(v147) >> (uint(int32(6)) % 32))
		} else {
			v157 = v154
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+3)) = uint8(v157)
		v161 = int32(255)
		v162 = int32(base.Ui32(v136)>>(uint(int32(18))%32)) & v161
		v167 = v131 + int32(base.Ui32(v162*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32))
		v169 = v167 + int32(-14234)
		if base.Ui32(v167) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
			v176 = int32(0)
		} else {
			v176 = v161
		}
		if base.Ui32(v169) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
			v179 = int32(base.Ui32(v169) >> (uint(int32(6)) % 32))
		} else {
			v179 = v176
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v179)
		v183 = int32(8)
		v190 = v131 - (int32(base.Ui32(v140*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v183)%32)) + int32(base.Ui32(v162*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v183)%32)))
		v192 = v190 + int32(_a_F_UpsampleArgbLinePair_C_8)
		if v190 < int32(-8708) {
			v199 = int32(0)
		} else {
			v199 = int32(255)
		}
		if base.Ui32(v192) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
			v202 = int32(base.Ui32(v192) >> (uint(int32(6)) % 32))
		} else {
			v202 = v199
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v202)
	}
	if int32(1) <= v122 {
		v211 = int32(1)
		v223 = int32(0)
		v227 = v223
		v230 = v223
		v234 = v36
		v235 = v43
		v236 = l0 + v211
		v237 = l1 + v211
		for {
			v254 = l6 + v227
			v257 = int32(255)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(4)))) = uint8(v257)
			v259 = int32(8)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v259))) = uint8(v257)
			v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v211+v230))))
			v268 = int32(16)
			v271 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v211+v230))))
			v272 = v267<<(uint(v268)%32) | v271
			v275 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v211+v230))))
			v279 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v211+v230))))
			v280 = v275<<(uint(v268)%32) | v279
			v283 = v235 + v234 + v272 + v280 + int32(524296)
			v285 = int32(1)
			v289 = int32(base.Ui32(v283+(v272+v235)<<(uint(v285)%32)) >> (uint(int32(3)) % 32))
			v290 = v289 + v234
			v294 = int32(base.Ui32(v290)>>(uint(v285)%32)) & v257
			v299 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236))))
			v303 = int32(base.Ui32(v299*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v259) % 32))
			v304 = int32(base.Ui32(v294*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v259)%32)) + v303
			v306 = v304 + int32(-17685)
			if base.Ui32(v304) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
				v313 = int32(0)
			} else {
				v313 = v257
			}
			if base.Ui32(v306) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v316 = int32(base.Ui32(v306) >> (uint(int32(6)) % 32))
			} else {
				v316 = v313
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(7)))) = uint8(v316)
			v322 = int32(255)
			v323 = int32(base.Ui32(v290)>>(uint(int32(17))%32)) & v322
			v328 = int32(base.Ui32(v323*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32)) + v303
			v330 = v328 + int32(-14234)
			if base.Ui32(v328) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
				v337 = int32(0)
			} else {
				v337 = v322
			}
			if base.Ui32(v330) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v340 = int32(base.Ui32(v330) >> (uint(int32(6)) % 32))
			} else {
				v340 = v337
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(5)))) = uint8(v340)
			v345 = int32(1)
			v349 = int32(base.Ui32(v283+(v280+v234)<<(uint(v345)%32)) >> (uint(int32(3)) % 32))
			v350 = v349 + v272
			v353 = int32(255)
			v354 = int32(base.Ui32(v350)>>(uint(v345)%32)) & v353
			v357 = int32(8)
			v361 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236+v345))))
			v365 = int32(base.Ui32(v361*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v357) % 32))
			v366 = int32(base.Ui32(v354*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v357)%32)) + v365
			v368 = v366 + int32(-17685)
			if base.Ui32(v366) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
				v375 = int32(0)
			} else {
				v375 = v353
			}
			if base.Ui32(v368) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v378 = int32(base.Ui32(v368) >> (uint(int32(6)) % 32))
			} else {
				v378 = v375
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(11)))) = uint8(v378)
			v384 = int32(255)
			v385 = int32(base.Ui32(v350)>>(uint(int32(17))%32)) & v384
			v390 = int32(base.Ui32(v385*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32)) + v365
			v392 = v390 + int32(-14234)
			if base.Ui32(v390) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
				v399 = int32(0)
			} else {
				v399 = v384
			}
			if base.Ui32(v392) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v402 = int32(base.Ui32(v392) >> (uint(int32(6)) % 32))
			} else {
				v402 = v399
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(9)))) = uint8(v402)
			v404 = int32(6)
			v408 = int32(8)
			v415 = v303 - (int32(base.Ui32(v323*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v408)%32)) + int32(base.Ui32(v294*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v408)%32)))
			v417 = v415 + int32(_a_F_UpsampleArgbLinePair_C_8)
			if v415 < int32(-8708) {
				v424 = int32(0)
			} else {
				v424 = int32(255)
			}
			if base.Ui32(v417) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v427 = int32(base.Ui32(v417) >> (uint(v404) % 32))
			} else {
				v427 = v424
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v404))) = uint8(v427)
			v433 = int32(8)
			v440 = v365 - (int32(base.Ui32(v385*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v433)%32)) + int32(base.Ui32(v354*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v433)%32)))
			v442 = v440 + int32(_a_F_UpsampleArgbLinePair_C_8)
			if v440 < int32(-8708) {
				v449 = int32(0)
			} else {
				v449 = int32(255)
			}
			if base.Ui32(v442) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v452 = int32(base.Ui32(v442) >> (uint(int32(6)) % 32))
			} else {
				v452 = v449
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(10)))) = uint8(v452)
			if l1 == int32(0) {
			} else {
				v456 = l7 + v227
				v459 = int32(255)
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(4)))) = uint8(v459)
				v461 = int32(8)
				*(*uint8)(unsafe.Add(mBase, uint32(v456+v461))) = uint8(v459)
				v467 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237))))
				v471 = int32(base.Ui32(v467*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v461) % 32))
				v472 = v349 + v235
				v476 = int32(base.Ui32(v472)>>(uint(int32(1))%32)) & v459
				v481 = v471 + int32(base.Ui32(v476*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v461)%32))
				v483 = v481 + int32(-17685)
				if base.Ui32(v481) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
					v490 = int32(0)
				} else {
					v490 = v459
				}
				if base.Ui32(v483) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
					v493 = int32(base.Ui32(v483) >> (uint(int32(6)) % 32))
				} else {
					v493 = v490
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(7)))) = uint8(v493)
				v499 = int32(255)
				v500 = int32(base.Ui32(v472)>>(uint(int32(17))%32)) & v499
				v505 = v471 + int32(base.Ui32(v500*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32))
				v507 = v505 + int32(-14234)
				if base.Ui32(v505) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
					v514 = int32(0)
				} else {
					v514 = v499
				}
				if base.Ui32(v507) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
					v517 = int32(base.Ui32(v507) >> (uint(int32(6)) % 32))
				} else {
					v517 = v514
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(5)))) = uint8(v517)
				v521 = int32(1)
				v523 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237+v521))))
				v526 = int32(8)
				v527 = int32(base.Ui32(v523*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v526) % 32))
				v528 = v289 + v280
				v531 = int32(255)
				v532 = int32(base.Ui32(v528)>>(uint(v521)%32)) & v531
				v537 = v527 + int32(base.Ui32(v532*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v526)%32))
				v539 = v537 + int32(-17685)
				if base.Ui32(v537) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
					v546 = int32(0)
				} else {
					v546 = v531
				}
				if base.Ui32(v539) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
					v549 = int32(base.Ui32(v539) >> (uint(int32(6)) % 32))
				} else {
					v549 = v546
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(11)))) = uint8(v549)
				v555 = int32(255)
				v556 = int32(base.Ui32(v528)>>(uint(int32(17))%32)) & v555
				v561 = v527 + int32(base.Ui32(v556*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32))
				v563 = v561 + int32(-14234)
				if base.Ui32(v561) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
					v570 = int32(0)
				} else {
					v570 = v555
				}
				if base.Ui32(v563) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
					v573 = int32(base.Ui32(v563) >> (uint(int32(6)) % 32))
				} else {
					v573 = v570
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(9)))) = uint8(v573)
				v575 = int32(6)
				v579 = int32(8)
				v586 = v471 - (int32(base.Ui32(v476*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v579)%32)) + int32(base.Ui32(v500*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v579)%32)))
				v588 = v586 + int32(_a_F_UpsampleArgbLinePair_C_8)
				if v586 < int32(-8708) {
					v595 = int32(0)
				} else {
					v595 = int32(255)
				}
				if base.Ui32(v588) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
					v598 = int32(base.Ui32(v588) >> (uint(v575) % 32))
				} else {
					v598 = v595
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+v575))) = uint8(v598)
				v604 = int32(8)
				v611 = v527 - (int32(base.Ui32(v532*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v604)%32)) + int32(base.Ui32(v556*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v604)%32)))
				v613 = v611 + int32(_a_F_UpsampleArgbLinePair_C_8)
				if v611 < int32(-8708) {
					v620 = int32(0)
				} else {
					v620 = int32(255)
				}
				if base.Ui32(v613) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
					v623 = int32(base.Ui32(v613) >> (uint(int32(6)) % 32))
				} else {
					v623 = v620
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(10)))) = uint8(v623)
			}
			v634 = int32(2)
			v641 = v230 + int32(1)
			if v122 != v641 {
				v227 = v227 + int32(8)
				v230 = v641
				v234 = v272
				v235 = v280
				v236 = v236 + v634
				v237 = v237 + v634
				continue
			} else {
				break
			}
			break
		}
		v646 = v280
		v647 = v272
	} else {
		v646 = v43
		v647 = v36
	}
	if l8&int32(1) != 0 {
	} else {
		v674 = int32(2)
		v675 = v120 << (uint(v674) % 32)
		v676 = l6 + v675
		v677 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(v676))) = uint8(v677)
		v680 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v120))))
		v683 = int32(8)
		v684 = int32(base.Ui32(v680*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v683) % 32))
		v689 = v646 + v647*int32(3) + int32(131074)
		v693 = int32(base.Ui32(v689)>>(uint(v674)%32)) & v677
		v698 = v684 + int32(base.Ui32(v693*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v683)%32))
		v700 = v698 + int32(-17685)
		if base.Ui32(v698) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
			v707 = int32(0)
		} else {
			v707 = v677
		}
		if base.Ui32(v700) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
			v710 = int32(base.Ui32(v700) >> (uint(int32(6)) % 32))
		} else {
			v710 = v707
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+3)) = uint8(v710)
		v714 = int32(255)
		v715 = int32(base.Ui32(v689)>>(uint(int32(18))%32)) & v714
		v720 = v684 + int32(base.Ui32(v715*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32))
		v722 = v720 + int32(-14234)
		if base.Ui32(v720) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
			v729 = int32(0)
		} else {
			v729 = v714
		}
		if base.Ui32(v722) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
			v732 = int32(base.Ui32(v722) >> (uint(int32(6)) % 32))
		} else {
			v732 = v729
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+1)) = uint8(v732)
		v736 = int32(8)
		v743 = v684 - (int32(base.Ui32(v693*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v736)%32)) + int32(base.Ui32(v715*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v736)%32)))
		v745 = v743 + int32(_a_F_UpsampleArgbLinePair_C_8)
		if v743 < int32(-8708) {
			v752 = int32(0)
		} else {
			v752 = int32(255)
		}
		if base.Ui32(v745) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
			v755 = int32(base.Ui32(v745) >> (uint(int32(6)) % 32))
		} else {
			v755 = v752
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+2)) = uint8(v755)
		if l1 == int32(0) {
		} else {
			v759 = l7 + v675
			v760 = int32(255)
			*(*uint8)(unsafe.Add(mBase, uint32(v759))) = uint8(v760)
			v763 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v120))))
			v766 = int32(8)
			v767 = int32(base.Ui32(v763*int32(_a_F_UpsampleArgbLinePair_C_0)) >> (uint(v766) % 32))
			v772 = v647 + v646*int32(3) + int32(131074)
			v776 = int32(base.Ui32(v772)>>(uint(int32(2))%32)) & v760
			v781 = v767 + int32(base.Ui32(v776*int32(_a_F_UpsampleArgbLinePair_C_1))>>(uint(v766)%32))
			v783 = v781 + int32(-17685)
			if base.Ui32(v781) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_2)) {
				v790 = int32(0)
			} else {
				v790 = v760
			}
			if base.Ui32(v783) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v793 = int32(base.Ui32(v783) >> (uint(int32(6)) % 32))
			} else {
				v793 = v790
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+3)) = uint8(v793)
			v797 = int32(255)
			v798 = int32(base.Ui32(v772)>>(uint(int32(18))%32)) & v797
			v803 = v767 + int32(base.Ui32(v798*int32(_a_F_UpsampleArgbLinePair_C_4))>>(uint(int32(8))%32))
			v805 = v803 + int32(-14234)
			if base.Ui32(v803) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_5)) {
				v812 = int32(0)
			} else {
				v812 = v797
			}
			if base.Ui32(v805) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v815 = int32(base.Ui32(v805) >> (uint(int32(6)) % 32))
			} else {
				v815 = v812
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+1)) = uint8(v815)
			v819 = int32(8)
			v826 = v767 - (int32(base.Ui32(v776*int32(_a_F_UpsampleArgbLinePair_C_7))>>(uint(v819)%32)) + int32(base.Ui32(v798*int32(_a_F_UpsampleArgbLinePair_C_6))>>(uint(v819)%32)))
			v828 = v826 + int32(_a_F_UpsampleArgbLinePair_C_8)
			if v826 < int32(-8708) {
				v835 = int32(0)
			} else {
				v835 = int32(255)
			}
			if base.Ui32(v828) < base.Ui32(int32(_a_F_UpsampleArgbLinePair_C_3)) {
				v838 = int32(base.Ui32(v828) >> (uint(int32(6)) % 32))
			} else {
				v838 = v835
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+2)) = uint8(v838)
		}
	}
	return
}
func F_UpsampleBgrLinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
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
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v168 int32
	_ = v168
	var v171 int32
	_ = v171
	var v175 int32
	_ = v175
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v203 int32
	_ = v203
	var v215 int32
	_ = v215
	var v219 int32
	_ = v219
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
	var v246 int32
	_ = v246
	var v251 int32
	_ = v251
	var v252 int32
	_ = v252
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v264 int32
	_ = v264
	var v267 int32
	_ = v267
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v281 int32
	_ = v281
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v300 int32
	_ = v300
	var v303 int32
	_ = v303
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v317 int32
	_ = v317
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v340 int32
	_ = v340
	var v341 int32
	_ = v341
	var v346 int32
	_ = v346
	var v348 int32
	_ = v348
	var v355 int32
	_ = v355
	var v358 int32
	_ = v358
	var v360 int32
	_ = v360
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v370 int32
	_ = v370
	var v372 int32
	_ = v372
	var v379 int32
	_ = v379
	var v382 int32
	_ = v382
	var v388 int32
	_ = v388
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v413 int32
	_ = v413
	var v420 int32
	_ = v420
	var v422 int32
	_ = v422
	var v429 int32
	_ = v429
	var v432 int32
	_ = v432
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
	var v446 int32
	_ = v446
	var v451 int32
	_ = v451
	var v453 int32
	_ = v453
	var v460 int32
	_ = v460
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
	var v476 int32
	_ = v476
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v519 int32
	_ = v519
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v529 int32
	_ = v529
	var v531 int32
	_ = v531
	var v538 int32
	_ = v538
	var v541 int32
	_ = v541
	var v547 int32
	_ = v547
	var v554 int32
	_ = v554
	var v556 int32
	_ = v556
	var v563 int32
	_ = v563
	var v566 int32
	_ = v566
	var v572 int32
	_ = v572
	var v579 int32
	_ = v579
	var v581 int32
	_ = v581
	var v588 int32
	_ = v588
	var v591 int32
	_ = v591
	var v602 int32
	_ = v602
	var v609 int32
	_ = v609
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v642 int32
	_ = v642
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v646 int32
	_ = v646
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v655 int32
	_ = v655
	var v657 int32
	_ = v657
	var v662 int32
	_ = v662
	var v664 int32
	_ = v664
	var v671 int32
	_ = v671
	var v674 int32
	_ = v674
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v684 int32
	_ = v684
	var v686 int32
	_ = v686
	var v693 int32
	_ = v693
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v716 int32
	_ = v716
	var v719 int32
	_ = v719
	var v723 int32
	_ = v723
	var v725 int32
	_ = v725
	var v728 int32
	_ = v728
	var v729 int32
	_ = v729
	var v734 int32
	_ = v734
	var v736 int32
	_ = v736
	var v741 int32
	_ = v741
	var v743 int32
	_ = v743
	var v750 int32
	_ = v750
	var v753 int32
	_ = v753
	var v757 int32
	_ = v757
	var v758 int32
	_ = v758
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v772 int32
	_ = v772
	var v775 int32
	_ = v775
	var v779 int32
	_ = v779
	var v786 int32
	_ = v786
	var v788 int32
	_ = v788
	var v795 int32
	_ = v795
	var v798 int32
	_ = v798
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v31 = int32(16)
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v34 = v30<<(uint(v31)%32) | v33
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v41 = v37<<(uint(v31)%32) | v40
	v44 = v34*int32(3) + v41 + int32(131074)
	v46 = int32(base.Ui32(v44) >> (uint(int32(18)) % 32))
	v49 = int32(8)
	v51 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v55 = int32(base.Ui32(v51*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v49) % 32))
	v56 = int32(base.Ui32(v46*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v49)%32)) + v55
	v58 = v56 + int32(-14234)
	if base.Ui32(v56) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
		v65 = int32(0)
	} else {
		v65 = int32(255)
	}
	if base.Ui32(v58) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
		v68 = int32(base.Ui32(v58) >> (uint(int32(6)) % 32))
	} else {
		v68 = v65
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v68)
	v72 = int32(255)
	v73 = int32(base.Ui32(v44)>>(uint(int32(2))%32)) & v72
	v78 = int32(base.Ui32(v73*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32)) + v55
	v80 = v78 + int32(-17685)
	if base.Ui32(v78) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
		v87 = int32(0)
	} else {
		v87 = v72
	}
	if base.Ui32(v80) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
		v90 = int32(base.Ui32(v80) >> (uint(int32(6)) % 32))
	} else {
		v90 = v87
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v90)
	v94 = int32(8)
	v101 = v55 - (int32(base.Ui32(v46*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v94)%32)) + int32(base.Ui32(v73*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v94)%32)))
	v103 = v101 + int32(_a_F_UpsampleBgrLinePair_C_8)
	if v101 < int32(-8708) {
		v110 = int32(0)
	} else {
		v110 = int32(255)
	}
	if base.Ui32(v103) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
		v113 = int32(base.Ui32(v103) >> (uint(int32(6)) % 32))
	} else {
		v113 = v110
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v113)
	v116 = l8 + int32(-1)
	v118 = v116 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v124 = int32(8)
		v125 = int32(base.Ui32(v121*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v124) % 32))
		v130 = v34 + v41*int32(3) + int32(131074)
		v132 = int32(base.Ui32(v130) >> (uint(int32(18)) % 32))
		v137 = v125 + int32(base.Ui32(v132*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v124)%32))
		v139 = v137 + int32(-14234)
		if base.Ui32(v137) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
			v146 = int32(0)
		} else {
			v146 = int32(255)
		}
		if base.Ui32(v139) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
			v149 = int32(base.Ui32(v139) >> (uint(int32(6)) % 32))
		} else {
			v149 = v146
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v149)
		v153 = int32(255)
		v154 = int32(base.Ui32(v130)>>(uint(int32(2))%32)) & v153
		v159 = v125 + int32(base.Ui32(v154*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32))
		v161 = v159 + int32(-17685)
		if base.Ui32(v159) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
			v168 = int32(0)
		} else {
			v168 = v153
		}
		if base.Ui32(v161) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
			v171 = int32(base.Ui32(v161) >> (uint(int32(6)) % 32))
		} else {
			v171 = v168
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v171)
		v175 = int32(8)
		v182 = v125 - (int32(base.Ui32(v154*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v175)%32)) + int32(base.Ui32(v132*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v175)%32)))
		v184 = v182 + int32(_a_F_UpsampleBgrLinePair_C_8)
		if v182 < int32(-8708) {
			v191 = int32(0)
		} else {
			v191 = int32(255)
		}
		if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
			v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
		} else {
			v194 = v191
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v194)
	}
	if int32(1) <= v118 {
		v203 = int32(1)
		v215 = int32(0)
		v219 = v215
		v226 = v34
		v227 = v41
		v228 = l0 + v203
		v229 = l1 + v203
		v230 = v215
		for {
			v246 = l6 + v219
			v251 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v203+v230))))
			v252 = int32(16)
			v255 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v203+v230))))
			v256 = v251<<(uint(v252)%32) | v255
			v259 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v203+v230))))
			v263 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v203+v230))))
			v264 = v259<<(uint(v252)%32) | v263
			v267 = v227 + v226 + v256 + v264 + int32(524296)
			v273 = int32(base.Ui32(v267+(v256+v227)<<(uint(int32(1))%32)) >> (uint(int32(3)) % 32))
			v274 = v273 + v226
			v276 = int32(base.Ui32(v274) >> (uint(int32(17)) % 32))
			v279 = int32(8)
			v281 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v228))))
			v285 = int32(base.Ui32(v281*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v279) % 32))
			v286 = int32(base.Ui32(v276*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v279)%32)) + v285
			v288 = v286 + int32(-14234)
			if base.Ui32(v286) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
				v295 = int32(0)
			} else {
				v295 = int32(255)
			}
			if base.Ui32(v288) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v298 = int32(base.Ui32(v288) >> (uint(int32(6)) % 32))
			} else {
				v298 = v295
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(5)))) = uint8(v298)
			v300 = int32(8)
			v303 = int32(1)
			v307 = int32(base.Ui32(v267+(v264+v226)<<(uint(v303)%32)) >> (uint(int32(3)) % 32))
			v308 = v307 + v256
			v310 = int32(base.Ui32(v308) >> (uint(int32(17)) % 32))
			v317 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v228+v303))))
			v321 = int32(base.Ui32(v317*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v300) % 32))
			v322 = int32(base.Ui32(v310*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v300)%32)) + v321
			v324 = v322 + int32(-14234)
			if base.Ui32(v322) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
				v331 = int32(0)
			} else {
				v331 = int32(255)
			}
			if base.Ui32(v324) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v334 = int32(base.Ui32(v324) >> (uint(int32(6)) % 32))
			} else {
				v334 = v331
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+v300))) = uint8(v334)
			v340 = int32(255)
			v341 = int32(base.Ui32(v274)>>(uint(int32(1))%32)) & v340
			v346 = int32(base.Ui32(v341*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32)) + v285
			v348 = v346 + int32(-17685)
			if base.Ui32(v346) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
				v355 = int32(0)
			} else {
				v355 = v340
			}
			if base.Ui32(v348) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v358 = int32(base.Ui32(v348) >> (uint(int32(6)) % 32))
			} else {
				v358 = v355
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(3)))) = uint8(v358)
			v360 = int32(6)
			v364 = int32(255)
			v365 = int32(base.Ui32(v308)>>(uint(int32(1))%32)) & v364
			v370 = int32(base.Ui32(v365*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32)) + v321
			v372 = v370 + int32(-17685)
			if base.Ui32(v370) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
				v379 = int32(0)
			} else {
				v379 = v364
			}
			if base.Ui32(v372) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v382 = int32(base.Ui32(v372) >> (uint(v360) % 32))
			} else {
				v382 = v379
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+v360))) = uint8(v382)
			v388 = int32(8)
			v395 = v285 - (int32(base.Ui32(v276*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v388)%32)) + int32(base.Ui32(v341*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v388)%32)))
			v397 = v395 + int32(_a_F_UpsampleBgrLinePair_C_8)
			if v395 < int32(-8708) {
				v404 = int32(0)
			} else {
				v404 = int32(255)
			}
			if base.Ui32(v397) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v407 = int32(base.Ui32(v397) >> (uint(int32(6)) % 32))
			} else {
				v407 = v404
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(4)))) = uint8(v407)
			v413 = int32(8)
			v420 = v321 - (int32(base.Ui32(v310*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v413)%32)) + int32(base.Ui32(v365*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v413)%32)))
			v422 = v420 + int32(_a_F_UpsampleBgrLinePair_C_8)
			if v420 < int32(-8708) {
				v429 = int32(0)
			} else {
				v429 = int32(255)
			}
			if base.Ui32(v422) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v432 = int32(base.Ui32(v422) >> (uint(int32(6)) % 32))
			} else {
				v432 = v429
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(7)))) = uint8(v432)
			if l1 == int32(0) {
			} else {
				v436 = l7 + v219
				v439 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v229))))
				v442 = int32(8)
				v443 = int32(base.Ui32(v439*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v442) % 32))
				v444 = v307 + v227
				v446 = int32(base.Ui32(v444) >> (uint(int32(17)) % 32))
				v451 = v443 + int32(base.Ui32(v446*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v442)%32))
				v453 = v451 + int32(-14234)
				if base.Ui32(v451) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
					v460 = int32(0)
				} else {
					v460 = int32(255)
				}
				if base.Ui32(v453) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
					v463 = int32(base.Ui32(v453) >> (uint(int32(6)) % 32))
				} else {
					v463 = v460
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(5)))) = uint8(v463)
				v465 = int32(8)
				v469 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v229+int32(1)))))
				v473 = int32(base.Ui32(v469*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v465) % 32))
				v474 = v273 + v264
				v476 = int32(base.Ui32(v474) >> (uint(int32(17)) % 32))
				v481 = v473 + int32(base.Ui32(v476*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v465)%32))
				v483 = v481 + int32(-14234)
				if base.Ui32(v481) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
					v490 = int32(0)
				} else {
					v490 = int32(255)
				}
				if base.Ui32(v483) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
					v493 = int32(base.Ui32(v483) >> (uint(int32(6)) % 32))
				} else {
					v493 = v490
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+v465))) = uint8(v493)
				v499 = int32(255)
				v500 = int32(base.Ui32(v444)>>(uint(int32(1))%32)) & v499
				v505 = v443 + int32(base.Ui32(v500*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32))
				v507 = v505 + int32(-17685)
				if base.Ui32(v505) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
					v514 = int32(0)
				} else {
					v514 = v499
				}
				if base.Ui32(v507) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
					v517 = int32(base.Ui32(v507) >> (uint(int32(6)) % 32))
				} else {
					v517 = v514
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(3)))) = uint8(v517)
				v519 = int32(6)
				v523 = int32(255)
				v524 = int32(base.Ui32(v474)>>(uint(int32(1))%32)) & v523
				v529 = v473 + int32(base.Ui32(v524*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32))
				v531 = v529 + int32(-17685)
				if base.Ui32(v529) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
					v538 = int32(0)
				} else {
					v538 = v523
				}
				if base.Ui32(v531) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
					v541 = int32(base.Ui32(v531) >> (uint(v519) % 32))
				} else {
					v541 = v538
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+v519))) = uint8(v541)
				v547 = int32(8)
				v554 = v443 - (int32(base.Ui32(v500*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v547)%32)) + int32(base.Ui32(v446*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v547)%32)))
				v556 = v554 + int32(_a_F_UpsampleBgrLinePair_C_8)
				if v554 < int32(-8708) {
					v563 = int32(0)
				} else {
					v563 = int32(255)
				}
				if base.Ui32(v556) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
					v566 = int32(base.Ui32(v556) >> (uint(int32(6)) % 32))
				} else {
					v566 = v563
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(4)))) = uint8(v566)
				v572 = int32(8)
				v579 = v473 - (int32(base.Ui32(v524*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v572)%32)) + int32(base.Ui32(v476*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v572)%32)))
				v581 = v579 + int32(_a_F_UpsampleBgrLinePair_C_8)
				if v579 < int32(-8708) {
					v588 = int32(0)
				} else {
					v588 = int32(255)
				}
				if base.Ui32(v581) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
					v591 = int32(base.Ui32(v581) >> (uint(int32(6)) % 32))
				} else {
					v591 = v588
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(7)))) = uint8(v591)
			}
			v602 = int32(2)
			v609 = v230 + int32(1)
			if v118 != v609 {
				v219 = v219 + int32(6)
				v226 = v256
				v227 = v264
				v228 = v228 + v602
				v229 = v229 + v602
				v230 = v609
				continue
			} else {
				break
			}
			break
		}
		v614 = v264
		v615 = v256
	} else {
		v614 = v41
		v615 = v34
	}
	if l8&int32(1) != 0 {
	} else {
		v642 = int32(3)
		v643 = v116 * v642
		v644 = l6 + v643
		v646 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v116))))
		v649 = int32(8)
		v650 = int32(base.Ui32(v646*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v649) % 32))
		v655 = v614 + v615*v642 + int32(131074)
		v657 = int32(base.Ui32(v655) >> (uint(int32(18)) % 32))
		v662 = v650 + int32(base.Ui32(v657*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v649)%32))
		v664 = v662 + int32(-14234)
		if base.Ui32(v662) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
			v671 = int32(0)
		} else {
			v671 = int32(255)
		}
		if base.Ui32(v664) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
			v674 = int32(base.Ui32(v664) >> (uint(int32(6)) % 32))
		} else {
			v674 = v671
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v644)+2)) = uint8(v674)
		v678 = int32(255)
		v679 = int32(base.Ui32(v655)>>(uint(int32(2))%32)) & v678
		v684 = v650 + int32(base.Ui32(v679*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32))
		v686 = v684 + int32(-17685)
		if base.Ui32(v684) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
			v693 = int32(0)
		} else {
			v693 = v678
		}
		if base.Ui32(v686) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
			v696 = int32(base.Ui32(v686) >> (uint(int32(6)) % 32))
		} else {
			v696 = v693
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v644))) = uint8(v696)
		v700 = int32(8)
		v707 = v650 - (int32(base.Ui32(v679*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v700)%32)) + int32(base.Ui32(v657*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v700)%32)))
		v709 = v707 + int32(_a_F_UpsampleBgrLinePair_C_8)
		if v707 < int32(-8708) {
			v716 = int32(0)
		} else {
			v716 = int32(255)
		}
		if base.Ui32(v709) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
			v719 = int32(base.Ui32(v709) >> (uint(int32(6)) % 32))
		} else {
			v719 = v716
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v644)+1)) = uint8(v719)
		if l1 == int32(0) {
		} else {
			v723 = l7 + v643
			v725 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v116))))
			v728 = int32(8)
			v729 = int32(base.Ui32(v725*int32(_a_F_UpsampleBgrLinePair_C_0)) >> (uint(v728) % 32))
			v734 = v615 + v614*int32(3) + int32(131074)
			v736 = int32(base.Ui32(v734) >> (uint(int32(18)) % 32))
			v741 = v729 + int32(base.Ui32(v736*int32(_a_F_UpsampleBgrLinePair_C_1))>>(uint(v728)%32))
			v743 = v741 + int32(-14234)
			if base.Ui32(v741) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_2)) {
				v750 = int32(0)
			} else {
				v750 = int32(255)
			}
			if base.Ui32(v743) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v753 = int32(base.Ui32(v743) >> (uint(int32(6)) % 32))
			} else {
				v753 = v750
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v723)+2)) = uint8(v753)
			v757 = int32(255)
			v758 = int32(base.Ui32(v734)>>(uint(int32(2))%32)) & v757
			v763 = v729 + int32(base.Ui32(v758*int32(_a_F_UpsampleBgrLinePair_C_4))>>(uint(int32(8))%32))
			v765 = v763 + int32(-17685)
			if base.Ui32(v763) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_5)) {
				v772 = int32(0)
			} else {
				v772 = v757
			}
			if base.Ui32(v765) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v775 = int32(base.Ui32(v765) >> (uint(int32(6)) % 32))
			} else {
				v775 = v772
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v723))) = uint8(v775)
			v779 = int32(8)
			v786 = v729 - (int32(base.Ui32(v758*int32(_a_F_UpsampleBgrLinePair_C_7))>>(uint(v779)%32)) + int32(base.Ui32(v736*int32(_a_F_UpsampleBgrLinePair_C_6))>>(uint(v779)%32)))
			v788 = v786 + int32(_a_F_UpsampleBgrLinePair_C_8)
			if v786 < int32(-8708) {
				v795 = int32(0)
			} else {
				v795 = int32(255)
			}
			if base.Ui32(v788) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_C_3)) {
				v798 = int32(base.Ui32(v788) >> (uint(int32(6)) % 32))
			} else {
				v798 = v795
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v723)+1)) = uint8(v798)
		}
	}
	return
}
func F_UpsampleBgraLinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v140 int32
	_ = v140
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v211 int32
	_ = v211
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v230 int32
	_ = v230
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v263 int32
	_ = v263
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v275 int32
	_ = v275
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v294 int32
	_ = v294
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v313 int32
	_ = v313
	var v316 int32
	_ = v316
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v337 int32
	_ = v337
	var v340 int32
	_ = v340
	var v345 int32
	_ = v345
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	var v353 int32
	_ = v353
	var v354 int32
	_ = v354
	var v357 int32
	_ = v357
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v368 int32
	_ = v368
	var v375 int32
	_ = v375
	var v378 int32
	_ = v378
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v390 int32
	_ = v390
	var v392 int32
	_ = v392
	var v399 int32
	_ = v399
	var v402 int32
	_ = v402
	var v408 int32
	_ = v408
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v433 int32
	_ = v433
	var v440 int32
	_ = v440
	var v442 int32
	_ = v442
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v456 int32
	_ = v456
	var v459 int32
	_ = v459
	var v465 int32
	_ = v465
	var v467 int32
	_ = v467
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v472 int32
	_ = v472
	var v476 int32
	_ = v476
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v523 int32
	_ = v523
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v537 int32
	_ = v537
	var v539 int32
	_ = v539
	var v546 int32
	_ = v546
	var v549 int32
	_ = v549
	var v551 int32
	_ = v551
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v570 int32
	_ = v570
	var v573 int32
	_ = v573
	var v579 int32
	_ = v579
	var v586 int32
	_ = v586
	var v588 int32
	_ = v588
	var v595 int32
	_ = v595
	var v598 int32
	_ = v598
	var v604 int32
	_ = v604
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v620 int32
	_ = v620
	var v623 int32
	_ = v623
	var v634 int32
	_ = v634
	var v641 int32
	_ = v641
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v680 int32
	_ = v680
	var v683 int32
	_ = v683
	var v684 int32
	_ = v684
	var v689 int32
	_ = v689
	var v693 int32
	_ = v693
	var v698 int32
	_ = v698
	var v700 int32
	_ = v700
	var v707 int32
	_ = v707
	var v710 int32
	_ = v710
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v729 int32
	_ = v729
	var v732 int32
	_ = v732
	var v736 int32
	_ = v736
	var v743 int32
	_ = v743
	var v745 int32
	_ = v745
	var v752 int32
	_ = v752
	var v755 int32
	_ = v755
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v763 int32
	_ = v763
	var v766 int32
	_ = v766
	var v767 int32
	_ = v767
	var v772 int32
	_ = v772
	var v776 int32
	_ = v776
	var v781 int32
	_ = v781
	var v783 int32
	_ = v783
	var v790 int32
	_ = v790
	var v793 int32
	_ = v793
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v803 int32
	_ = v803
	var v805 int32
	_ = v805
	var v812 int32
	_ = v812
	var v815 int32
	_ = v815
	var v819 int32
	_ = v819
	var v826 int32
	_ = v826
	var v828 int32
	_ = v828
	var v835 int32
	_ = v835
	var v838 int32
	_ = v838
	v30 = int32(255)
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+3)) = uint8(v30)
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v33 = int32(16)
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v36 = v32<<(uint(v33)%32) | v35
	v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v43 = v39<<(uint(v33)%32) | v42
	v46 = v36*int32(3) + v43 + int32(131074)
	v50 = int32(base.Ui32(v46)>>(uint(int32(18))%32)) & v30
	v53 = int32(8)
	v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v59 = int32(base.Ui32(v55*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v53) % 32))
	v60 = int32(base.Ui32(v50*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v53)%32)) + v59
	v62 = v60 + int32(-14234)
	if base.Ui32(v60) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
		v69 = int32(0)
	} else {
		v69 = v30
	}
	if base.Ui32(v62) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
		v72 = int32(base.Ui32(v62) >> (uint(int32(6)) % 32))
	} else {
		v72 = v69
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v72)
	v76 = int32(255)
	v77 = int32(base.Ui32(v46)>>(uint(int32(2))%32)) & v76
	v82 = int32(base.Ui32(v77*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(int32(8))%32)) + v59
	v84 = v82 + int32(-17685)
	if base.Ui32(v82) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
		v91 = int32(0)
	} else {
		v91 = v76
	}
	if base.Ui32(v84) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
		v94 = int32(base.Ui32(v84) >> (uint(int32(6)) % 32))
	} else {
		v94 = v91
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v94)
	v98 = int32(8)
	v105 = v59 - (int32(base.Ui32(v50*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v98)%32)) + int32(base.Ui32(v77*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v98)%32)))
	v107 = v105 + int32(_a_F_UpsampleBgraLinePair_C_8)
	if v105 < int32(-8708) {
		v114 = int32(0)
	} else {
		v114 = int32(255)
	}
	if base.Ui32(v107) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
		v117 = int32(base.Ui32(v107) >> (uint(int32(6)) % 32))
	} else {
		v117 = v114
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v117)
	v120 = l8 + int32(-1)
	v122 = v120 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v125 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+3)) = uint8(v125)
		v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v130 = int32(8)
		v131 = int32(base.Ui32(v127*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v130) % 32))
		v136 = v36 + v43*int32(3) + int32(131074)
		v140 = int32(base.Ui32(v136)>>(uint(int32(18))%32)) & v125
		v145 = v131 + int32(base.Ui32(v140*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v130)%32))
		v147 = v145 + int32(-14234)
		if base.Ui32(v145) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
			v154 = int32(0)
		} else {
			v154 = v125
		}
		if base.Ui32(v147) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
			v157 = int32(base.Ui32(v147) >> (uint(int32(6)) % 32))
		} else {
			v157 = v154
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v157)
		v161 = int32(255)
		v162 = int32(base.Ui32(v136)>>(uint(int32(2))%32)) & v161
		v167 = v131 + int32(base.Ui32(v162*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(int32(8))%32))
		v169 = v167 + int32(-17685)
		if base.Ui32(v167) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
			v176 = int32(0)
		} else {
			v176 = v161
		}
		if base.Ui32(v169) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
			v179 = int32(base.Ui32(v169) >> (uint(int32(6)) % 32))
		} else {
			v179 = v176
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v179)
		v183 = int32(8)
		v190 = v131 - (int32(base.Ui32(v162*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v183)%32)) + int32(base.Ui32(v140*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v183)%32)))
		v192 = v190 + int32(_a_F_UpsampleBgraLinePair_C_8)
		if v190 < int32(-8708) {
			v199 = int32(0)
		} else {
			v199 = int32(255)
		}
		if base.Ui32(v192) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
			v202 = int32(base.Ui32(v192) >> (uint(int32(6)) % 32))
		} else {
			v202 = v199
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v202)
	}
	if int32(1) <= v122 {
		v211 = int32(1)
		v223 = int32(0)
		v227 = v223
		v230 = v223
		v234 = v36
		v235 = v43
		v236 = l0 + v211
		v237 = l1 + v211
		for {
			v254 = l6 + v227
			v257 = int32(255)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(7)))) = uint8(v257)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(11)))) = uint8(v257)
			v263 = int32(6)
			v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v211+v230))))
			v268 = int32(16)
			v271 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v211+v230))))
			v272 = v267<<(uint(v268)%32) | v271
			v275 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v211+v230))))
			v279 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v211+v230))))
			v280 = v275<<(uint(v268)%32) | v279
			v283 = v235 + v234 + v272 + v280 + int32(524296)
			v289 = int32(base.Ui32(v283+(v272+v235)<<(uint(int32(1))%32)) >> (uint(int32(3)) % 32))
			v290 = v289 + v234
			v294 = int32(base.Ui32(v290)>>(uint(int32(17))%32)) & v257
			v297 = int32(8)
			v299 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236))))
			v303 = int32(base.Ui32(v299*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v297) % 32))
			v304 = int32(base.Ui32(v294*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v297)%32)) + v303
			v306 = v304 + int32(-14234)
			if base.Ui32(v304) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
				v313 = int32(0)
			} else {
				v313 = v257
			}
			if base.Ui32(v306) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v316 = int32(base.Ui32(v306) >> (uint(v263) % 32))
			} else {
				v316 = v313
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v263))) = uint8(v316)
			v322 = int32(255)
			v323 = int32(base.Ui32(v290)>>(uint(int32(1))%32)) & v322
			v328 = int32(base.Ui32(v323*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(int32(8))%32)) + v303
			v330 = v328 + int32(-17685)
			if base.Ui32(v328) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
				v337 = int32(0)
			} else {
				v337 = v322
			}
			if base.Ui32(v330) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v340 = int32(base.Ui32(v330) >> (uint(int32(6)) % 32))
			} else {
				v340 = v337
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(4)))) = uint8(v340)
			v345 = int32(1)
			v349 = int32(base.Ui32(v283+(v280+v234)<<(uint(v345)%32)) >> (uint(int32(3)) % 32))
			v350 = v349 + v272
			v353 = int32(255)
			v354 = int32(base.Ui32(v350)>>(uint(int32(17))%32)) & v353
			v357 = int32(8)
			v361 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236+v345))))
			v365 = int32(base.Ui32(v361*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v357) % 32))
			v366 = int32(base.Ui32(v354*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v357)%32)) + v365
			v368 = v366 + int32(-14234)
			if base.Ui32(v366) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
				v375 = int32(0)
			} else {
				v375 = v353
			}
			if base.Ui32(v368) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v378 = int32(base.Ui32(v368) >> (uint(int32(6)) % 32))
			} else {
				v378 = v375
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(10)))) = uint8(v378)
			v380 = int32(8)
			v384 = int32(255)
			v385 = int32(base.Ui32(v350)>>(uint(int32(1))%32)) & v384
			v390 = int32(base.Ui32(v385*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(v380)%32)) + v365
			v392 = v390 + int32(-17685)
			if base.Ui32(v390) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
				v399 = int32(0)
			} else {
				v399 = v384
			}
			if base.Ui32(v392) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v402 = int32(base.Ui32(v392) >> (uint(int32(6)) % 32))
			} else {
				v402 = v399
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v380))) = uint8(v402)
			v408 = int32(8)
			v415 = v303 - (int32(base.Ui32(v294*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v408)%32)) + int32(base.Ui32(v323*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v408)%32)))
			v417 = v415 + int32(_a_F_UpsampleBgraLinePair_C_8)
			if v415 < int32(-8708) {
				v424 = int32(0)
			} else {
				v424 = int32(255)
			}
			if base.Ui32(v417) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v427 = int32(base.Ui32(v417) >> (uint(int32(6)) % 32))
			} else {
				v427 = v424
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(5)))) = uint8(v427)
			v433 = int32(8)
			v440 = v365 - (int32(base.Ui32(v354*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v433)%32)) + int32(base.Ui32(v385*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v433)%32)))
			v442 = v440 + int32(_a_F_UpsampleBgraLinePair_C_8)
			if v440 < int32(-8708) {
				v449 = int32(0)
			} else {
				v449 = int32(255)
			}
			if base.Ui32(v442) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v452 = int32(base.Ui32(v442) >> (uint(int32(6)) % 32))
			} else {
				v452 = v449
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(9)))) = uint8(v452)
			if l1 == int32(0) {
			} else {
				v456 = l7 + v227
				v459 = int32(255)
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(7)))) = uint8(v459)
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(11)))) = uint8(v459)
				v465 = int32(6)
				v467 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237))))
				v470 = int32(8)
				v471 = int32(base.Ui32(v467*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v470) % 32))
				v472 = v349 + v235
				v476 = int32(base.Ui32(v472)>>(uint(int32(17))%32)) & v459
				v481 = v471 + int32(base.Ui32(v476*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v470)%32))
				v483 = v481 + int32(-14234)
				if base.Ui32(v481) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
					v490 = int32(0)
				} else {
					v490 = v459
				}
				if base.Ui32(v483) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
					v493 = int32(base.Ui32(v483) >> (uint(v465) % 32))
				} else {
					v493 = v490
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+v465))) = uint8(v493)
				v499 = int32(255)
				v500 = int32(base.Ui32(v472)>>(uint(int32(1))%32)) & v499
				v505 = v471 + int32(base.Ui32(v500*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(int32(8))%32))
				v507 = v505 + int32(-17685)
				if base.Ui32(v505) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
					v514 = int32(0)
				} else {
					v514 = v499
				}
				if base.Ui32(v507) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
					v517 = int32(base.Ui32(v507) >> (uint(int32(6)) % 32))
				} else {
					v517 = v514
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(4)))) = uint8(v517)
				v523 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237+int32(1)))))
				v526 = int32(8)
				v527 = int32(base.Ui32(v523*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v526) % 32))
				v528 = v289 + v280
				v531 = int32(255)
				v532 = int32(base.Ui32(v528)>>(uint(int32(17))%32)) & v531
				v537 = v527 + int32(base.Ui32(v532*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v526)%32))
				v539 = v537 + int32(-14234)
				if base.Ui32(v537) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
					v546 = int32(0)
				} else {
					v546 = v531
				}
				if base.Ui32(v539) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
					v549 = int32(base.Ui32(v539) >> (uint(int32(6)) % 32))
				} else {
					v549 = v546
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(10)))) = uint8(v549)
				v551 = int32(8)
				v555 = int32(255)
				v556 = int32(base.Ui32(v528)>>(uint(int32(1))%32)) & v555
				v561 = v527 + int32(base.Ui32(v556*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(v551)%32))
				v563 = v561 + int32(-17685)
				if base.Ui32(v561) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
					v570 = int32(0)
				} else {
					v570 = v555
				}
				if base.Ui32(v563) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
					v573 = int32(base.Ui32(v563) >> (uint(int32(6)) % 32))
				} else {
					v573 = v570
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+v551))) = uint8(v573)
				v579 = int32(8)
				v586 = v471 - (int32(base.Ui32(v500*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v579)%32)) + int32(base.Ui32(v476*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v579)%32)))
				v588 = v586 + int32(_a_F_UpsampleBgraLinePair_C_8)
				if v586 < int32(-8708) {
					v595 = int32(0)
				} else {
					v595 = int32(255)
				}
				if base.Ui32(v588) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
					v598 = int32(base.Ui32(v588) >> (uint(int32(6)) % 32))
				} else {
					v598 = v595
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(5)))) = uint8(v598)
				v604 = int32(8)
				v611 = v527 - (int32(base.Ui32(v556*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v604)%32)) + int32(base.Ui32(v532*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v604)%32)))
				v613 = v611 + int32(_a_F_UpsampleBgraLinePair_C_8)
				if v611 < int32(-8708) {
					v620 = int32(0)
				} else {
					v620 = int32(255)
				}
				if base.Ui32(v613) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
					v623 = int32(base.Ui32(v613) >> (uint(int32(6)) % 32))
				} else {
					v623 = v620
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(9)))) = uint8(v623)
			}
			v634 = int32(2)
			v641 = v230 + int32(1)
			if v122 != v641 {
				v227 = v227 + int32(8)
				v230 = v641
				v234 = v272
				v235 = v280
				v236 = v236 + v634
				v237 = v237 + v634
				continue
			} else {
				break
			}
			break
		}
		v646 = v280
		v647 = v272
	} else {
		v646 = v43
		v647 = v36
	}
	if l8&int32(1) != 0 {
	} else {
		v675 = v120 << (uint(int32(2)) % 32)
		v676 = l6 + v675
		v677 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+3)) = uint8(v677)
		v680 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v120))))
		v683 = int32(8)
		v684 = int32(base.Ui32(v680*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v683) % 32))
		v689 = v646 + v647*int32(3) + int32(131074)
		v693 = int32(base.Ui32(v689)>>(uint(int32(18))%32)) & v677
		v698 = v684 + int32(base.Ui32(v693*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v683)%32))
		v700 = v698 + int32(-14234)
		if base.Ui32(v698) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
			v707 = int32(0)
		} else {
			v707 = v677
		}
		if base.Ui32(v700) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
			v710 = int32(base.Ui32(v700) >> (uint(int32(6)) % 32))
		} else {
			v710 = v707
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+2)) = uint8(v710)
		v714 = int32(255)
		v715 = int32(base.Ui32(v689)>>(uint(int32(2))%32)) & v714
		v720 = v684 + int32(base.Ui32(v715*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(int32(8))%32))
		v722 = v720 + int32(-17685)
		if base.Ui32(v720) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
			v729 = int32(0)
		} else {
			v729 = v714
		}
		if base.Ui32(v722) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
			v732 = int32(base.Ui32(v722) >> (uint(int32(6)) % 32))
		} else {
			v732 = v729
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676))) = uint8(v732)
		v736 = int32(8)
		v743 = v684 - (int32(base.Ui32(v715*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v736)%32)) + int32(base.Ui32(v693*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v736)%32)))
		v745 = v743 + int32(_a_F_UpsampleBgraLinePair_C_8)
		if v743 < int32(-8708) {
			v752 = int32(0)
		} else {
			v752 = int32(255)
		}
		if base.Ui32(v745) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
			v755 = int32(base.Ui32(v745) >> (uint(int32(6)) % 32))
		} else {
			v755 = v752
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+1)) = uint8(v755)
		if l1 == int32(0) {
		} else {
			v759 = l7 + v675
			v760 = int32(255)
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+3)) = uint8(v760)
			v763 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v120))))
			v766 = int32(8)
			v767 = int32(base.Ui32(v763*int32(_a_F_UpsampleBgraLinePair_C_0)) >> (uint(v766) % 32))
			v772 = v647 + v646*int32(3) + int32(131074)
			v776 = int32(base.Ui32(v772)>>(uint(int32(18))%32)) & v760
			v781 = v767 + int32(base.Ui32(v776*int32(_a_F_UpsampleBgraLinePair_C_1))>>(uint(v766)%32))
			v783 = v781 + int32(-14234)
			if base.Ui32(v781) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_2)) {
				v790 = int32(0)
			} else {
				v790 = v760
			}
			if base.Ui32(v783) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v793 = int32(base.Ui32(v783) >> (uint(int32(6)) % 32))
			} else {
				v793 = v790
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+2)) = uint8(v793)
			v797 = int32(255)
			v798 = int32(base.Ui32(v772)>>(uint(int32(2))%32)) & v797
			v803 = v767 + int32(base.Ui32(v798*int32(_a_F_UpsampleBgraLinePair_C_4))>>(uint(int32(8))%32))
			v805 = v803 + int32(-17685)
			if base.Ui32(v803) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_5)) {
				v812 = int32(0)
			} else {
				v812 = v797
			}
			if base.Ui32(v805) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v815 = int32(base.Ui32(v805) >> (uint(int32(6)) % 32))
			} else {
				v815 = v812
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759))) = uint8(v815)
			v819 = int32(8)
			v826 = v767 - (int32(base.Ui32(v798*int32(_a_F_UpsampleBgraLinePair_C_7))>>(uint(v819)%32)) + int32(base.Ui32(v776*int32(_a_F_UpsampleBgraLinePair_C_6))>>(uint(v819)%32)))
			v828 = v826 + int32(_a_F_UpsampleBgraLinePair_C_8)
			if v826 < int32(-8708) {
				v835 = int32(0)
			} else {
				v835 = int32(255)
			}
			if base.Ui32(v828) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_C_3)) {
				v838 = int32(base.Ui32(v828) >> (uint(int32(6)) % 32))
			} else {
				v838 = v835
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+1)) = uint8(v838)
		}
	}
	return
}
func F_UpsampleRgb565LinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v74 int32
	_ = v74
	var v77 int32
	_ = v77
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v141 int32
	_ = v141
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v155 int32
	_ = v155
	var v158 int32
	_ = v158
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v189 int32
	_ = v189
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v208 int32
	_ = v208
	var v211 int32
	_ = v211
	var v212 int32
	_ = v212
	var v221 int32
	_ = v221
	var v233 int32
	_ = v233
	var v237 int32
	_ = v237
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v248 int32
	_ = v248
	var v264 int32
	_ = v264
	var v267 int32
	_ = v267
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v282 int32
	_ = v282
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v290 int32
	_ = v290
	var v292 int32
	_ = v292
	var v296 int32
	_ = v296
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v313 int32
	_ = v313
	var v315 int32
	_ = v315
	var v322 int32
	_ = v322
	var v325 int32
	_ = v325
	var v332 int32
	_ = v332
	var v334 int32
	_ = v334
	var v341 int32
	_ = v341
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v351 int32
	_ = v351
	var v353 int32
	_ = v353
	var v356 int32
	_ = v356
	var v357 int32
	_ = v357
	var v363 int32
	_ = v363
	var v364 int32
	_ = v364
	var v366 int32
	_ = v366
	var v373 int32
	_ = v373
	var v374 int32
	_ = v374
	var v380 int32
	_ = v380
	var v382 int32
	_ = v382
	var v389 int32
	_ = v389
	var v392 int32
	_ = v392
	var v399 int32
	_ = v399
	var v401 int32
	_ = v401
	var v408 int32
	_ = v408
	var v411 int32
	_ = v411
	var v414 int32
	_ = v414
	var v416 int32
	_ = v416
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v435 int32
	_ = v435
	var v438 int32
	_ = v438
	var v439 int32
	_ = v439
	var v451 int32
	_ = v451
	var v453 int32
	_ = v453
	var v460 int32
	_ = v460
	var v463 int32
	_ = v463
	var v464 int32
	_ = v464
	var v468 int32
	_ = v468
	var v471 int32
	_ = v471
	var v474 int32
	_ = v474
	var v475 int32
	_ = v475
	var v476 int32
	_ = v476
	var v478 int32
	_ = v478
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v492 int32
	_ = v492
	var v495 int32
	_ = v495
	var v500 int32
	_ = v500
	var v501 int32
	_ = v501
	var v504 int32
	_ = v504
	var v511 int32
	_ = v511
	var v513 int32
	_ = v513
	var v520 int32
	_ = v520
	var v523 int32
	_ = v523
	var v526 int32
	_ = v526
	var v532 int32
	_ = v532
	var v535 int32
	_ = v535
	var v536 int32
	_ = v536
	var v537 int32
	_ = v537
	var v539 int32
	_ = v539
	var v544 int32
	_ = v544
	var v546 int32
	_ = v546
	var v553 int32
	_ = v553
	var v556 int32
	_ = v556
	var v561 int32
	_ = v561
	var v562 int32
	_ = v562
	var v565 int32
	_ = v565
	var v572 int32
	_ = v572
	var v574 int32
	_ = v574
	var v581 int32
	_ = v581
	var v584 int32
	_ = v584
	var v587 int32
	_ = v587
	var v589 int32
	_ = v589
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v608 int32
	_ = v608
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v624 int32
	_ = v624
	var v626 int32
	_ = v626
	var v633 int32
	_ = v633
	var v636 int32
	_ = v636
	var v637 int32
	_ = v637
	var v648 int32
	_ = v648
	var v655 int32
	_ = v655
	var v660 int32
	_ = v660
	var v662 int32
	_ = v662
	var v689 int32
	_ = v689
	var v690 int32
	_ = v690
	var v692 int32
	_ = v692
	var v695 int32
	_ = v695
	var v696 int32
	_ = v696
	var v701 int32
	_ = v701
	var v703 int32
	_ = v703
	var v708 int32
	_ = v708
	var v710 int32
	_ = v710
	var v717 int32
	_ = v717
	var v720 int32
	_ = v720
	var v725 int32
	_ = v725
	var v726 int32
	_ = v726
	var v729 int32
	_ = v729
	var v736 int32
	_ = v736
	var v738 int32
	_ = v738
	var v745 int32
	_ = v745
	var v748 int32
	_ = v748
	var v751 int32
	_ = v751
	var v761 int32
	_ = v761
	var v763 int32
	_ = v763
	var v770 int32
	_ = v770
	var v773 int32
	_ = v773
	var v774 int32
	_ = v774
	var v778 int32
	_ = v778
	var v780 int32
	_ = v780
	var v783 int32
	_ = v783
	var v784 int32
	_ = v784
	var v789 int32
	_ = v789
	var v791 int32
	_ = v791
	var v796 int32
	_ = v796
	var v798 int32
	_ = v798
	var v805 int32
	_ = v805
	var v808 int32
	_ = v808
	var v813 int32
	_ = v813
	var v814 int32
	_ = v814
	var v817 int32
	_ = v817
	var v824 int32
	_ = v824
	var v826 int32
	_ = v826
	var v833 int32
	_ = v833
	var v836 int32
	_ = v836
	var v839 int32
	_ = v839
	var v849 int32
	_ = v849
	var v851 int32
	_ = v851
	var v858 int32
	_ = v858
	var v861 int32
	_ = v861
	var v862 int32
	_ = v862
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v33 = int32(8)
	v34 = int32(base.Ui32(v30*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v33) % 32))
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v36 = int32(16)
	v38 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v39 = v35<<(uint(v36)%32) | v38
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v46 = v42<<(uint(v36)%32) | v45
	v49 = v39*int32(3) + v46 + int32(131074)
	v51 = int32(base.Ui32(v49) >> (uint(int32(18)) % 32))
	v58 = int32(255)
	v59 = int32(base.Ui32(v49)>>(uint(int32(2))%32)) & v58
	v65 = v34 - (int32(base.Ui32(v51*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v33)%32)) + int32(base.Ui32(v59*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v33)%32)))
	v67 = v65 + int32(_a_F_UpsampleRgb565LinePair_C_3)
	if v65 < int32(-8708) {
		v74 = int32(0)
	} else {
		v74 = v58
	}
	if base.Ui32(v67) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
		v77 = int32(base.Ui32(v67) >> (uint(int32(6)) % 32))
	} else {
		v77 = v74
	}
	v84 = int32(base.Ui32(v51*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(int32(8))%32)) + v34
	v86 = v84 + int32(-14234)
	if base.Ui32(v84) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
		v93 = int32(0)
	} else {
		v93 = int32(248)
	}
	if base.Ui32(v86) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
		v96 = int32(base.Ui32(v86) >> (uint(int32(6)) % 32))
	} else {
		v96 = v93
	}
	v99 = int32(base.Ui32(v77)>>(uint(int32(5))%32)) | v96&int32(248)
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v99)
	v109 = int32(base.Ui32(v59*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32)) + v34
	v111 = v109 + int32(-17685)
	if base.Ui32(v109) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
		v118 = int32(0)
	} else {
		v118 = int32(31)
	}
	if base.Ui32(v111) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
		v121 = int32(base.Ui32(v111) >> (uint(int32(9)) % 32))
	} else {
		v121 = v118
	}
	v122 = v77<<(uint(int32(3))%32)&int32(224) | v121
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v122)
	v125 = l8 + int32(-1)
	v127 = v125 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v133 = int32(8)
		v134 = int32(base.Ui32(v130*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v133) % 32))
		v139 = v39 + v46*int32(3) + int32(131074)
		v141 = int32(base.Ui32(v139) >> (uint(int32(18)) % 32))
		v146 = v134 + int32(base.Ui32(v141*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(v133)%32))
		v148 = v146 + int32(-14234)
		if base.Ui32(v146) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
			v155 = int32(0)
		} else {
			v155 = int32(248)
		}
		if base.Ui32(v148) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
			v158 = int32(base.Ui32(v148) >> (uint(int32(6)) % 32))
		} else {
			v158 = v155
		}
		v163 = int32(255)
		v164 = int32(base.Ui32(v139)>>(uint(int32(2))%32)) & v163
		v167 = int32(8)
		v174 = v134 - (int32(base.Ui32(v164*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v167)%32)) + int32(base.Ui32(v141*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v167)%32)))
		v176 = v174 + int32(_a_F_UpsampleRgb565LinePair_C_3)
		if v174 < int32(-8708) {
			v183 = int32(0)
		} else {
			v183 = v163
		}
		if base.Ui32(v176) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
			v186 = int32(base.Ui32(v176) >> (uint(int32(6)) % 32))
		} else {
			v186 = v183
		}
		v189 = v158&int32(248) | int32(base.Ui32(v186)>>(uint(int32(5))%32))
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v189)
		v199 = v134 + int32(base.Ui32(v164*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32))
		v201 = v199 + int32(-17685)
		if base.Ui32(v199) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
			v208 = int32(0)
		} else {
			v208 = int32(31)
		}
		if base.Ui32(v201) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
			v211 = int32(base.Ui32(v201) >> (uint(int32(9)) % 32))
		} else {
			v211 = v208
		}
		v212 = v186<<(uint(int32(3))%32)&int32(224) | v211
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v212)
	}
	if int32(1) <= v127 {
		v221 = int32(1)
		v233 = int32(0)
		v237 = v233
		v244 = v233
		v245 = v39
		v246 = v46
		v247 = l0 + v221
		v248 = l1 + v221
		for {
			v264 = l6 + v237
			v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v247))))
			v270 = int32(8)
			v271 = int32(base.Ui32(v267*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v270) % 32))
			v274 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v221+v244))))
			v275 = int32(16)
			v278 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v221+v244))))
			v279 = v274<<(uint(v275)%32) | v278
			v282 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v221+v244))))
			v286 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v221+v244))))
			v287 = v282<<(uint(v275)%32) | v286
			v290 = v246 + v245 + v279 + v287 + int32(524296)
			v292 = int32(1)
			v296 = int32(base.Ui32(v290+(v279+v246)<<(uint(v292)%32)) >> (uint(int32(3)) % 32))
			v297 = v296 + v245
			v299 = int32(base.Ui32(v297) >> (uint(int32(17)) % 32))
			v306 = int32(255)
			v307 = int32(base.Ui32(v297)>>(uint(v292)%32)) & v306
			v313 = v271 - (int32(base.Ui32(v299*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v270)%32)) + int32(base.Ui32(v307*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v270)%32)))
			v315 = v313 + int32(_a_F_UpsampleRgb565LinePair_C_3)
			if v313 < int32(-8708) {
				v322 = int32(0)
			} else {
				v322 = v306
			}
			if base.Ui32(v315) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v325 = int32(base.Ui32(v315) >> (uint(int32(6)) % 32))
			} else {
				v325 = v322
			}
			v332 = int32(base.Ui32(v299*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(int32(8))%32)) + v271
			v334 = v332 + int32(-14234)
			if base.Ui32(v332) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
				v341 = int32(0)
			} else {
				v341 = int32(248)
			}
			if base.Ui32(v334) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v344 = int32(base.Ui32(v334) >> (uint(int32(6)) % 32))
			} else {
				v344 = v341
			}
			v347 = int32(base.Ui32(v325)>>(uint(int32(5))%32)) | v344&int32(248)
			*(*uint8)(unsafe.Add(mBase, uint32(v264+int32(2)))) = uint8(v347)
			v351 = int32(1)
			v353 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v247+v351))))
			v356 = int32(8)
			v357 = int32(base.Ui32(v353*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v356) % 32))
			v363 = int32(base.Ui32(v290+(v287+v245)<<(uint(v351)%32)) >> (uint(int32(3)) % 32))
			v364 = v363 + v279
			v366 = int32(base.Ui32(v364) >> (uint(int32(17)) % 32))
			v373 = int32(255)
			v374 = int32(base.Ui32(v364)>>(uint(v351)%32)) & v373
			v380 = v357 - (int32(base.Ui32(v366*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v356)%32)) + int32(base.Ui32(v374*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v356)%32)))
			v382 = v380 + int32(_a_F_UpsampleRgb565LinePair_C_3)
			if v380 < int32(-8708) {
				v389 = int32(0)
			} else {
				v389 = v373
			}
			if base.Ui32(v382) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v392 = int32(base.Ui32(v382) >> (uint(int32(6)) % 32))
			} else {
				v392 = v389
			}
			v399 = int32(base.Ui32(v366*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(int32(8))%32)) + v357
			v401 = v399 + int32(-14234)
			if base.Ui32(v399) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
				v408 = int32(0)
			} else {
				v408 = int32(248)
			}
			if base.Ui32(v401) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v411 = int32(base.Ui32(v401) >> (uint(int32(6)) % 32))
			} else {
				v411 = v408
			}
			v414 = int32(base.Ui32(v392)>>(uint(int32(5))%32)) | v411&int32(248)
			*(*uint8)(unsafe.Add(mBase, uint32(v264+int32(4)))) = uint8(v414)
			v416 = int32(3)
			v426 = int32(base.Ui32(v307*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32)) + v271
			v428 = v426 + int32(-17685)
			if base.Ui32(v426) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
				v435 = int32(0)
			} else {
				v435 = int32(31)
			}
			if base.Ui32(v428) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v438 = int32(base.Ui32(v428) >> (uint(int32(9)) % 32))
			} else {
				v438 = v435
			}
			v439 = v325<<(uint(v416)%32)&int32(224) | v438
			*(*uint8)(unsafe.Add(mBase, uint32(v264+v416))) = uint8(v439)
			v451 = int32(base.Ui32(v374*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32)) + v357
			v453 = v451 + int32(-17685)
			if base.Ui32(v451) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
				v460 = int32(0)
			} else {
				v460 = int32(31)
			}
			if base.Ui32(v453) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v463 = int32(base.Ui32(v453) >> (uint(int32(9)) % 32))
			} else {
				v463 = v460
			}
			v464 = v392<<(uint(int32(3))%32)&int32(224) | v463
			*(*uint8)(unsafe.Add(mBase, uint32(v264+int32(5)))) = uint8(v464)
			if l1 == int32(0) {
			} else {
				v468 = l7 + v237
				v471 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v248))))
				v474 = int32(8)
				v475 = int32(base.Ui32(v471*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v474) % 32))
				v476 = v363 + v246
				v478 = int32(base.Ui32(v476) >> (uint(int32(17)) % 32))
				v483 = v475 + int32(base.Ui32(v478*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(v474)%32))
				v485 = v483 + int32(-14234)
				if base.Ui32(v483) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
					v492 = int32(0)
				} else {
					v492 = int32(248)
				}
				if base.Ui32(v485) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
					v495 = int32(base.Ui32(v485) >> (uint(int32(6)) % 32))
				} else {
					v495 = v492
				}
				v500 = int32(255)
				v501 = int32(base.Ui32(v476)>>(uint(int32(1))%32)) & v500
				v504 = int32(8)
				v511 = v475 - (int32(base.Ui32(v501*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v504)%32)) + int32(base.Ui32(v478*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v504)%32)))
				v513 = v511 + int32(_a_F_UpsampleRgb565LinePair_C_3)
				if v511 < int32(-8708) {
					v520 = int32(0)
				} else {
					v520 = v500
				}
				if base.Ui32(v513) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
					v523 = int32(base.Ui32(v513) >> (uint(int32(6)) % 32))
				} else {
					v523 = v520
				}
				v526 = v495&int32(248) | int32(base.Ui32(v523)>>(uint(int32(5))%32))
				*(*uint8)(unsafe.Add(mBase, uint32(v468+int32(2)))) = uint8(v526)
				v532 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v248+int32(1)))))
				v535 = int32(8)
				v536 = int32(base.Ui32(v532*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v535) % 32))
				v537 = v296 + v287
				v539 = int32(base.Ui32(v537) >> (uint(int32(17)) % 32))
				v544 = v536 + int32(base.Ui32(v539*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(v535)%32))
				v546 = v544 + int32(-14234)
				if base.Ui32(v544) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
					v553 = int32(0)
				} else {
					v553 = int32(248)
				}
				if base.Ui32(v546) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
					v556 = int32(base.Ui32(v546) >> (uint(int32(6)) % 32))
				} else {
					v556 = v553
				}
				v561 = int32(255)
				v562 = int32(base.Ui32(v537)>>(uint(int32(1))%32)) & v561
				v565 = int32(8)
				v572 = v536 - (int32(base.Ui32(v562*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v565)%32)) + int32(base.Ui32(v539*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v565)%32)))
				v574 = v572 + int32(_a_F_UpsampleRgb565LinePair_C_3)
				if v572 < int32(-8708) {
					v581 = int32(0)
				} else {
					v581 = v561
				}
				if base.Ui32(v574) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
					v584 = int32(base.Ui32(v574) >> (uint(int32(6)) % 32))
				} else {
					v584 = v581
				}
				v587 = v556&int32(248) | int32(base.Ui32(v584)>>(uint(int32(5))%32))
				*(*uint8)(unsafe.Add(mBase, uint32(v468+int32(4)))) = uint8(v587)
				v589 = int32(3)
				v599 = v475 + int32(base.Ui32(v501*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32))
				v601 = v599 + int32(-17685)
				if base.Ui32(v599) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
					v608 = int32(0)
				} else {
					v608 = int32(31)
				}
				if base.Ui32(v601) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
					v611 = int32(base.Ui32(v601) >> (uint(int32(9)) % 32))
				} else {
					v611 = v608
				}
				v612 = v523<<(uint(v589)%32)&int32(224) | v611
				*(*uint8)(unsafe.Add(mBase, uint32(v468+v589))) = uint8(v612)
				v624 = v536 + int32(base.Ui32(v562*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32))
				v626 = v624 + int32(-17685)
				if base.Ui32(v624) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
					v633 = int32(0)
				} else {
					v633 = int32(31)
				}
				if base.Ui32(v626) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
					v636 = int32(base.Ui32(v626) >> (uint(int32(9)) % 32))
				} else {
					v636 = v633
				}
				v637 = v584<<(uint(int32(3))%32)&int32(224) | v636
				*(*uint8)(unsafe.Add(mBase, uint32(v468+int32(5)))) = uint8(v637)
			}
			v648 = int32(2)
			v655 = v244 + int32(1)
			if v127 != v655 {
				v237 = v237 + int32(4)
				v244 = v655
				v245 = v279
				v246 = v287
				v247 = v247 + v648
				v248 = v248 + v648
				continue
			} else {
				break
			}
			break
		}
		v660 = v287
		v662 = v279
	} else {
		v660 = v46
		v662 = v39
	}
	if l8&int32(1) != 0 {
	} else {
		v689 = v125 << (uint(int32(1)) % 32)
		v690 = l6 + v689
		v692 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v125))))
		v695 = int32(8)
		v696 = int32(base.Ui32(v692*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v695) % 32))
		v701 = v660 + v662*int32(3) + int32(131074)
		v703 = int32(base.Ui32(v701) >> (uint(int32(18)) % 32))
		v708 = v696 + int32(base.Ui32(v703*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(v695)%32))
		v710 = v708 + int32(-14234)
		if base.Ui32(v708) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
			v717 = int32(0)
		} else {
			v717 = int32(248)
		}
		if base.Ui32(v710) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
			v720 = int32(base.Ui32(v710) >> (uint(int32(6)) % 32))
		} else {
			v720 = v717
		}
		v725 = int32(255)
		v726 = int32(base.Ui32(v701)>>(uint(int32(2))%32)) & v725
		v729 = int32(8)
		v736 = v696 - (int32(base.Ui32(v726*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v729)%32)) + int32(base.Ui32(v703*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v729)%32)))
		v738 = v736 + int32(_a_F_UpsampleRgb565LinePair_C_3)
		if v736 < int32(-8708) {
			v745 = int32(0)
		} else {
			v745 = v725
		}
		if base.Ui32(v738) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
			v748 = int32(base.Ui32(v738) >> (uint(int32(6)) % 32))
		} else {
			v748 = v745
		}
		v751 = v720&int32(248) | int32(base.Ui32(v748)>>(uint(int32(5))%32))
		*(*uint8)(unsafe.Add(mBase, uint32(v690))) = uint8(v751)
		v761 = v696 + int32(base.Ui32(v726*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32))
		v763 = v761 + int32(-17685)
		if base.Ui32(v761) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
			v770 = int32(0)
		} else {
			v770 = int32(31)
		}
		if base.Ui32(v763) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
			v773 = int32(base.Ui32(v763) >> (uint(int32(9)) % 32))
		} else {
			v773 = v770
		}
		v774 = v748<<(uint(int32(3))%32)&int32(224) | v773
		*(*uint8)(unsafe.Add(mBase, uint32(v690)+1)) = uint8(v774)
		if l1 == int32(0) {
		} else {
			v778 = l7 + v689
			v780 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v125))))
			v783 = int32(8)
			v784 = int32(base.Ui32(v780*int32(_a_F_UpsampleRgb565LinePair_C_0)) >> (uint(v783) % 32))
			v789 = v662 + v660*int32(3) + int32(131074)
			v791 = int32(base.Ui32(v789) >> (uint(int32(18)) % 32))
			v796 = v784 + int32(base.Ui32(v791*int32(_a_F_UpsampleRgb565LinePair_C_5))>>(uint(v783)%32))
			v798 = v796 + int32(-14234)
			if base.Ui32(v796) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_6)) {
				v805 = int32(0)
			} else {
				v805 = int32(248)
			}
			if base.Ui32(v798) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v808 = int32(base.Ui32(v798) >> (uint(int32(6)) % 32))
			} else {
				v808 = v805
			}
			v813 = int32(255)
			v814 = int32(base.Ui32(v789)>>(uint(int32(2))%32)) & v813
			v817 = int32(8)
			v824 = v784 - (int32(base.Ui32(v814*int32(_a_F_UpsampleRgb565LinePair_C_2))>>(uint(v817)%32)) + int32(base.Ui32(v791*int32(_a_F_UpsampleRgb565LinePair_C_1))>>(uint(v817)%32)))
			v826 = v824 + int32(_a_F_UpsampleRgb565LinePair_C_3)
			if v824 < int32(-8708) {
				v833 = int32(0)
			} else {
				v833 = v813
			}
			if base.Ui32(v826) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v836 = int32(base.Ui32(v826) >> (uint(int32(6)) % 32))
			} else {
				v836 = v833
			}
			v839 = v808&int32(248) | int32(base.Ui32(v836)>>(uint(int32(5))%32))
			*(*uint8)(unsafe.Add(mBase, uint32(v778))) = uint8(v839)
			v849 = v784 + int32(base.Ui32(v814*int32(_a_F_UpsampleRgb565LinePair_C_7))>>(uint(int32(8))%32))
			v851 = v849 + int32(-17685)
			if base.Ui32(v849) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_8)) {
				v858 = int32(0)
			} else {
				v858 = int32(31)
			}
			if base.Ui32(v851) < base.Ui32(int32(_a_F_UpsampleRgb565LinePair_C_4)) {
				v861 = int32(base.Ui32(v851) >> (uint(int32(9)) % 32))
			} else {
				v861 = v858
			}
			v862 = v836<<(uint(int32(3))%32)&int32(224) | v861
			*(*uint8)(unsafe.Add(mBase, uint32(v778)+1)) = uint8(v862)
		}
	}
	return
}
func F_UpsampleRgbLinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v58 int32
	_ = v58
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v78 int32
	_ = v78
	var v80 int32
	_ = v80
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
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
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v168 int32
	_ = v168
	var v171 int32
	_ = v171
	var v175 int32
	_ = v175
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v203 int32
	_ = v203
	var v215 int32
	_ = v215
	var v219 int32
	_ = v219
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
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v251 int32
	_ = v251
	var v252 int32
	_ = v252
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v264 int32
	_ = v264
	var v267 int32
	_ = v267
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v281 int32
	_ = v281
	var v285 int32
	_ = v285
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v300 int32
	_ = v300
	var v303 int32
	_ = v303
	var v307 int32
	_ = v307
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v313 int32
	_ = v313
	var v317 int32
	_ = v317
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v340 int32
	_ = v340
	var v341 int32
	_ = v341
	var v346 int32
	_ = v346
	var v348 int32
	_ = v348
	var v355 int32
	_ = v355
	var v358 int32
	_ = v358
	var v360 int32
	_ = v360
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v370 int32
	_ = v370
	var v372 int32
	_ = v372
	var v379 int32
	_ = v379
	var v382 int32
	_ = v382
	var v388 int32
	_ = v388
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v413 int32
	_ = v413
	var v420 int32
	_ = v420
	var v422 int32
	_ = v422
	var v429 int32
	_ = v429
	var v432 int32
	_ = v432
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
	var v446 int32
	_ = v446
	var v451 int32
	_ = v451
	var v453 int32
	_ = v453
	var v460 int32
	_ = v460
	var v463 int32
	_ = v463
	var v465 int32
	_ = v465
	var v469 int32
	_ = v469
	var v472 int32
	_ = v472
	var v473 int32
	_ = v473
	var v474 int32
	_ = v474
	var v476 int32
	_ = v476
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v519 int32
	_ = v519
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v529 int32
	_ = v529
	var v531 int32
	_ = v531
	var v538 int32
	_ = v538
	var v541 int32
	_ = v541
	var v547 int32
	_ = v547
	var v554 int32
	_ = v554
	var v556 int32
	_ = v556
	var v563 int32
	_ = v563
	var v566 int32
	_ = v566
	var v572 int32
	_ = v572
	var v579 int32
	_ = v579
	var v581 int32
	_ = v581
	var v588 int32
	_ = v588
	var v591 int32
	_ = v591
	var v602 int32
	_ = v602
	var v609 int32
	_ = v609
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v642 int32
	_ = v642
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v646 int32
	_ = v646
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v655 int32
	_ = v655
	var v657 int32
	_ = v657
	var v662 int32
	_ = v662
	var v664 int32
	_ = v664
	var v671 int32
	_ = v671
	var v674 int32
	_ = v674
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v684 int32
	_ = v684
	var v686 int32
	_ = v686
	var v693 int32
	_ = v693
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v716 int32
	_ = v716
	var v719 int32
	_ = v719
	var v723 int32
	_ = v723
	var v725 int32
	_ = v725
	var v728 int32
	_ = v728
	var v729 int32
	_ = v729
	var v734 int32
	_ = v734
	var v736 int32
	_ = v736
	var v741 int32
	_ = v741
	var v743 int32
	_ = v743
	var v750 int32
	_ = v750
	var v753 int32
	_ = v753
	var v757 int32
	_ = v757
	var v758 int32
	_ = v758
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v772 int32
	_ = v772
	var v775 int32
	_ = v775
	var v779 int32
	_ = v779
	var v786 int32
	_ = v786
	var v788 int32
	_ = v788
	var v795 int32
	_ = v795
	var v798 int32
	_ = v798
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v31 = int32(16)
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v34 = v30<<(uint(v31)%32) | v33
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v41 = v37<<(uint(v31)%32) | v40
	v44 = v34*int32(3) + v41 + int32(131074)
	v46 = int32(base.Ui32(v44) >> (uint(int32(18)) % 32))
	v49 = int32(8)
	v51 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v55 = int32(base.Ui32(v51*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v49) % 32))
	v56 = int32(base.Ui32(v46*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v49)%32)) + v55
	v58 = v56 + int32(-14234)
	if base.Ui32(v56) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
		v65 = int32(0)
	} else {
		v65 = int32(255)
	}
	if base.Ui32(v58) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
		v68 = int32(base.Ui32(v58) >> (uint(int32(6)) % 32))
	} else {
		v68 = v65
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v68)
	v72 = int32(255)
	v73 = int32(base.Ui32(v44)>>(uint(int32(2))%32)) & v72
	v78 = int32(base.Ui32(v73*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(int32(8))%32)) + v55
	v80 = v78 + int32(-17685)
	if base.Ui32(v78) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
		v87 = int32(0)
	} else {
		v87 = v72
	}
	if base.Ui32(v80) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
		v90 = int32(base.Ui32(v80) >> (uint(int32(6)) % 32))
	} else {
		v90 = v87
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v90)
	v94 = int32(8)
	v101 = v55 - (int32(base.Ui32(v46*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v94)%32)) + int32(base.Ui32(v73*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v94)%32)))
	v103 = v101 + int32(_a_F_UpsampleRgbLinePair_C_8)
	if v101 < int32(-8708) {
		v110 = int32(0)
	} else {
		v110 = int32(255)
	}
	if base.Ui32(v103) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
		v113 = int32(base.Ui32(v103) >> (uint(int32(6)) % 32))
	} else {
		v113 = v110
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v113)
	v116 = l8 + int32(-1)
	v118 = v116 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v124 = int32(8)
		v125 = int32(base.Ui32(v121*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v124) % 32))
		v130 = v34 + v41*int32(3) + int32(131074)
		v132 = int32(base.Ui32(v130) >> (uint(int32(18)) % 32))
		v137 = v125 + int32(base.Ui32(v132*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v124)%32))
		v139 = v137 + int32(-14234)
		if base.Ui32(v137) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
			v146 = int32(0)
		} else {
			v146 = int32(255)
		}
		if base.Ui32(v139) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
			v149 = int32(base.Ui32(v139) >> (uint(int32(6)) % 32))
		} else {
			v149 = v146
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v149)
		v153 = int32(255)
		v154 = int32(base.Ui32(v130)>>(uint(int32(2))%32)) & v153
		v159 = v125 + int32(base.Ui32(v154*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(int32(8))%32))
		v161 = v159 + int32(-17685)
		if base.Ui32(v159) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
			v168 = int32(0)
		} else {
			v168 = v153
		}
		if base.Ui32(v161) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
			v171 = int32(base.Ui32(v161) >> (uint(int32(6)) % 32))
		} else {
			v171 = v168
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v171)
		v175 = int32(8)
		v182 = v125 - (int32(base.Ui32(v154*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v175)%32)) + int32(base.Ui32(v132*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v175)%32)))
		v184 = v182 + int32(_a_F_UpsampleRgbLinePair_C_8)
		if v182 < int32(-8708) {
			v191 = int32(0)
		} else {
			v191 = int32(255)
		}
		if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
			v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
		} else {
			v194 = v191
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v194)
	}
	if int32(1) <= v118 {
		v203 = int32(1)
		v215 = int32(0)
		v219 = v215
		v226 = v34
		v227 = v41
		v228 = l0 + v203
		v229 = l1 + v203
		v230 = v215
		for {
			v246 = l6 + v219
			v247 = int32(3)
			v251 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v203+v230))))
			v252 = int32(16)
			v255 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v203+v230))))
			v256 = v251<<(uint(v252)%32) | v255
			v259 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v203+v230))))
			v263 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v203+v230))))
			v264 = v259<<(uint(v252)%32) | v263
			v267 = v227 + v226 + v256 + v264 + int32(524296)
			v273 = int32(base.Ui32(v267+(v256+v227)<<(uint(int32(1))%32)) >> (uint(v247) % 32))
			v274 = v273 + v226
			v276 = int32(base.Ui32(v274) >> (uint(int32(17)) % 32))
			v279 = int32(8)
			v281 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v228))))
			v285 = int32(base.Ui32(v281*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v279) % 32))
			v286 = int32(base.Ui32(v276*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v279)%32)) + v285
			v288 = v286 + int32(-14234)
			if base.Ui32(v286) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
				v295 = int32(0)
			} else {
				v295 = int32(255)
			}
			if base.Ui32(v288) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v298 = int32(base.Ui32(v288) >> (uint(int32(6)) % 32))
			} else {
				v298 = v295
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+v247))) = uint8(v298)
			v300 = int32(6)
			v303 = int32(1)
			v307 = int32(base.Ui32(v267+(v264+v226)<<(uint(v303)%32)) >> (uint(int32(3)) % 32))
			v308 = v307 + v256
			v310 = int32(base.Ui32(v308) >> (uint(int32(17)) % 32))
			v313 = int32(8)
			v317 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v228+v303))))
			v321 = int32(base.Ui32(v317*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v313) % 32))
			v322 = int32(base.Ui32(v310*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v313)%32)) + v321
			v324 = v322 + int32(-14234)
			if base.Ui32(v322) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
				v331 = int32(0)
			} else {
				v331 = int32(255)
			}
			if base.Ui32(v324) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v334 = int32(base.Ui32(v324) >> (uint(v300) % 32))
			} else {
				v334 = v331
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+v300))) = uint8(v334)
			v340 = int32(255)
			v341 = int32(base.Ui32(v274)>>(uint(int32(1))%32)) & v340
			v346 = int32(base.Ui32(v341*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(int32(8))%32)) + v285
			v348 = v346 + int32(-17685)
			if base.Ui32(v346) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
				v355 = int32(0)
			} else {
				v355 = v340
			}
			if base.Ui32(v348) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v358 = int32(base.Ui32(v348) >> (uint(int32(6)) % 32))
			} else {
				v358 = v355
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(5)))) = uint8(v358)
			v360 = int32(8)
			v364 = int32(255)
			v365 = int32(base.Ui32(v308)>>(uint(int32(1))%32)) & v364
			v370 = int32(base.Ui32(v365*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(v360)%32)) + v321
			v372 = v370 + int32(-17685)
			if base.Ui32(v370) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
				v379 = int32(0)
			} else {
				v379 = v364
			}
			if base.Ui32(v372) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v382 = int32(base.Ui32(v372) >> (uint(int32(6)) % 32))
			} else {
				v382 = v379
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+v360))) = uint8(v382)
			v388 = int32(8)
			v395 = v285 - (int32(base.Ui32(v276*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v388)%32)) + int32(base.Ui32(v341*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v388)%32)))
			v397 = v395 + int32(_a_F_UpsampleRgbLinePair_C_8)
			if v395 < int32(-8708) {
				v404 = int32(0)
			} else {
				v404 = int32(255)
			}
			if base.Ui32(v397) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v407 = int32(base.Ui32(v397) >> (uint(int32(6)) % 32))
			} else {
				v407 = v404
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(4)))) = uint8(v407)
			v413 = int32(8)
			v420 = v321 - (int32(base.Ui32(v310*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v413)%32)) + int32(base.Ui32(v365*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v413)%32)))
			v422 = v420 + int32(_a_F_UpsampleRgbLinePair_C_8)
			if v420 < int32(-8708) {
				v429 = int32(0)
			} else {
				v429 = int32(255)
			}
			if base.Ui32(v422) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v432 = int32(base.Ui32(v422) >> (uint(int32(6)) % 32))
			} else {
				v432 = v429
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v246+int32(7)))) = uint8(v432)
			if l1 == int32(0) {
			} else {
				v436 = l7 + v219
				v439 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v229))))
				v442 = int32(8)
				v443 = int32(base.Ui32(v439*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v442) % 32))
				v444 = v307 + v227
				v446 = int32(base.Ui32(v444) >> (uint(int32(17)) % 32))
				v451 = v443 + int32(base.Ui32(v446*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v442)%32))
				v453 = v451 + int32(-14234)
				if base.Ui32(v451) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
					v460 = int32(0)
				} else {
					v460 = int32(255)
				}
				if base.Ui32(v453) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
					v463 = int32(base.Ui32(v453) >> (uint(int32(6)) % 32))
				} else {
					v463 = v460
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(3)))) = uint8(v463)
				v465 = int32(6)
				v469 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v229+int32(1)))))
				v472 = int32(8)
				v473 = int32(base.Ui32(v469*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v472) % 32))
				v474 = v273 + v264
				v476 = int32(base.Ui32(v474) >> (uint(int32(17)) % 32))
				v481 = v473 + int32(base.Ui32(v476*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v472)%32))
				v483 = v481 + int32(-14234)
				if base.Ui32(v481) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
					v490 = int32(0)
				} else {
					v490 = int32(255)
				}
				if base.Ui32(v483) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
					v493 = int32(base.Ui32(v483) >> (uint(v465) % 32))
				} else {
					v493 = v490
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+v465))) = uint8(v493)
				v499 = int32(255)
				v500 = int32(base.Ui32(v444)>>(uint(int32(1))%32)) & v499
				v505 = v443 + int32(base.Ui32(v500*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(int32(8))%32))
				v507 = v505 + int32(-17685)
				if base.Ui32(v505) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
					v514 = int32(0)
				} else {
					v514 = v499
				}
				if base.Ui32(v507) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
					v517 = int32(base.Ui32(v507) >> (uint(int32(6)) % 32))
				} else {
					v517 = v514
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(5)))) = uint8(v517)
				v519 = int32(8)
				v523 = int32(255)
				v524 = int32(base.Ui32(v474)>>(uint(int32(1))%32)) & v523
				v529 = v473 + int32(base.Ui32(v524*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(v519)%32))
				v531 = v529 + int32(-17685)
				if base.Ui32(v529) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
					v538 = int32(0)
				} else {
					v538 = v523
				}
				if base.Ui32(v531) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
					v541 = int32(base.Ui32(v531) >> (uint(int32(6)) % 32))
				} else {
					v541 = v538
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+v519))) = uint8(v541)
				v547 = int32(8)
				v554 = v443 - (int32(base.Ui32(v500*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v547)%32)) + int32(base.Ui32(v446*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v547)%32)))
				v556 = v554 + int32(_a_F_UpsampleRgbLinePair_C_8)
				if v554 < int32(-8708) {
					v563 = int32(0)
				} else {
					v563 = int32(255)
				}
				if base.Ui32(v556) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
					v566 = int32(base.Ui32(v556) >> (uint(int32(6)) % 32))
				} else {
					v566 = v563
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(4)))) = uint8(v566)
				v572 = int32(8)
				v579 = v473 - (int32(base.Ui32(v524*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v572)%32)) + int32(base.Ui32(v476*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v572)%32)))
				v581 = v579 + int32(_a_F_UpsampleRgbLinePair_C_8)
				if v579 < int32(-8708) {
					v588 = int32(0)
				} else {
					v588 = int32(255)
				}
				if base.Ui32(v581) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
					v591 = int32(base.Ui32(v581) >> (uint(int32(6)) % 32))
				} else {
					v591 = v588
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(7)))) = uint8(v591)
			}
			v602 = int32(2)
			v609 = v230 + int32(1)
			if v118 != v609 {
				v219 = v219 + int32(6)
				v226 = v256
				v227 = v264
				v228 = v228 + v602
				v229 = v229 + v602
				v230 = v609
				continue
			} else {
				break
			}
			break
		}
		v614 = v264
		v615 = v256
	} else {
		v614 = v41
		v615 = v34
	}
	if l8&int32(1) != 0 {
	} else {
		v642 = int32(3)
		v643 = v116 * v642
		v644 = l6 + v643
		v646 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v116))))
		v649 = int32(8)
		v650 = int32(base.Ui32(v646*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v649) % 32))
		v655 = v614 + v615*v642 + int32(131074)
		v657 = int32(base.Ui32(v655) >> (uint(int32(18)) % 32))
		v662 = v650 + int32(base.Ui32(v657*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v649)%32))
		v664 = v662 + int32(-14234)
		if base.Ui32(v662) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
			v671 = int32(0)
		} else {
			v671 = int32(255)
		}
		if base.Ui32(v664) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
			v674 = int32(base.Ui32(v664) >> (uint(int32(6)) % 32))
		} else {
			v674 = v671
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v644))) = uint8(v674)
		v678 = int32(255)
		v679 = int32(base.Ui32(v655)>>(uint(int32(2))%32)) & v678
		v684 = v650 + int32(base.Ui32(v679*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(int32(8))%32))
		v686 = v684 + int32(-17685)
		if base.Ui32(v684) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
			v693 = int32(0)
		} else {
			v693 = v678
		}
		if base.Ui32(v686) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
			v696 = int32(base.Ui32(v686) >> (uint(int32(6)) % 32))
		} else {
			v696 = v693
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v644)+2)) = uint8(v696)
		v700 = int32(8)
		v707 = v650 - (int32(base.Ui32(v679*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v700)%32)) + int32(base.Ui32(v657*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v700)%32)))
		v709 = v707 + int32(_a_F_UpsampleRgbLinePair_C_8)
		if v707 < int32(-8708) {
			v716 = int32(0)
		} else {
			v716 = int32(255)
		}
		if base.Ui32(v709) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
			v719 = int32(base.Ui32(v709) >> (uint(int32(6)) % 32))
		} else {
			v719 = v716
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v644)+1)) = uint8(v719)
		if l1 == int32(0) {
		} else {
			v723 = l7 + v643
			v725 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v116))))
			v728 = int32(8)
			v729 = int32(base.Ui32(v725*int32(_a_F_UpsampleRgbLinePair_C_0)) >> (uint(v728) % 32))
			v734 = v615 + v614*int32(3) + int32(131074)
			v736 = int32(base.Ui32(v734) >> (uint(int32(18)) % 32))
			v741 = v729 + int32(base.Ui32(v736*int32(_a_F_UpsampleRgbLinePair_C_1))>>(uint(v728)%32))
			v743 = v741 + int32(-14234)
			if base.Ui32(v741) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_2)) {
				v750 = int32(0)
			} else {
				v750 = int32(255)
			}
			if base.Ui32(v743) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v753 = int32(base.Ui32(v743) >> (uint(int32(6)) % 32))
			} else {
				v753 = v750
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v723))) = uint8(v753)
			v757 = int32(255)
			v758 = int32(base.Ui32(v734)>>(uint(int32(2))%32)) & v757
			v763 = v729 + int32(base.Ui32(v758*int32(_a_F_UpsampleRgbLinePair_C_4))>>(uint(int32(8))%32))
			v765 = v763 + int32(-17685)
			if base.Ui32(v763) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_5)) {
				v772 = int32(0)
			} else {
				v772 = v757
			}
			if base.Ui32(v765) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v775 = int32(base.Ui32(v765) >> (uint(int32(6)) % 32))
			} else {
				v775 = v772
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v723)+2)) = uint8(v775)
			v779 = int32(8)
			v786 = v729 - (int32(base.Ui32(v758*int32(_a_F_UpsampleRgbLinePair_C_7))>>(uint(v779)%32)) + int32(base.Ui32(v736*int32(_a_F_UpsampleRgbLinePair_C_6))>>(uint(v779)%32)))
			v788 = v786 + int32(_a_F_UpsampleRgbLinePair_C_8)
			if v786 < int32(-8708) {
				v795 = int32(0)
			} else {
				v795 = int32(255)
			}
			if base.Ui32(v788) < base.Ui32(int32(_a_F_UpsampleRgbLinePair_C_3)) {
				v798 = int32(base.Ui32(v788) >> (uint(int32(6)) % 32))
			} else {
				v798 = v795
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v723)+1)) = uint8(v798)
		}
	}
	return
}
func F_UpsampleRgba4444LinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v48 int32
	_ = v48
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v97 int32
	_ = v97
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v128 int32
	_ = v128
	var v129 int32
	_ = v129
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v157 int32
	_ = v157
	var v160 int32
	_ = v160
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v182 int32
	_ = v182
	var v189 int32
	_ = v189
	var v191 int32
	_ = v191
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v211 int32
	_ = v211
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v238 int32
	_ = v238
	var v254 int32
	_ = v254
	var v255 int32
	_ = v255
	var v259 int32
	_ = v259
	var v260 int32
	_ = v260
	var v263 int32
	_ = v263
	var v264 int32
	_ = v264
	var v267 int32
	_ = v267
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v275 int32
	_ = v275
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v286 int32
	_ = v286
	var v289 int32
	_ = v289
	var v291 int32
	_ = v291
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v298 int32
	_ = v298
	var v305 int32
	_ = v305
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v315 int32
	_ = v315
	var v319 int32
	_ = v319
	var v320 int32
	_ = v320
	var v324 int32
	_ = v324
	var v327 int32
	_ = v327
	var v331 int32
	_ = v331
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v338 int32
	_ = v338
	var v345 int32
	_ = v345
	var v348 int32
	_ = v348
	var v350 int32
	_ = v350
	var v355 int32
	_ = v355
	var v360 int32
	_ = v360
	var v362 int32
	_ = v362
	var v369 int32
	_ = v369
	var v372 int32
	_ = v372
	var v377 int32
	_ = v377
	var v384 int32
	_ = v384
	var v386 int32
	_ = v386
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v397 int32
	_ = v397
	var v402 int32
	_ = v402
	var v407 int32
	_ = v407
	var v409 int32
	_ = v409
	var v416 int32
	_ = v416
	var v419 int32
	_ = v419
	var v424 int32
	_ = v424
	var v431 int32
	_ = v431
	var v433 int32
	_ = v433
	var v440 int32
	_ = v440
	var v443 int32
	_ = v443
	var v444 int32
	_ = v444
	var v448 int32
	_ = v448
	var v451 int32
	_ = v451
	var v454 int32
	_ = v454
	var v455 int32
	_ = v455
	var v456 int32
	_ = v456
	var v460 int32
	_ = v460
	var v465 int32
	_ = v465
	var v467 int32
	_ = v467
	var v474 int32
	_ = v474
	var v477 int32
	_ = v477
	var v479 int32
	_ = v479
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v494 int32
	_ = v494
	var v499 int32
	_ = v499
	var v501 int32
	_ = v501
	var v508 int32
	_ = v508
	var v511 int32
	_ = v511
	var v513 int32
	_ = v513
	var v518 int32
	_ = v518
	var v523 int32
	_ = v523
	var v525 int32
	_ = v525
	var v532 int32
	_ = v532
	var v535 int32
	_ = v535
	var v540 int32
	_ = v540
	var v547 int32
	_ = v547
	var v549 int32
	_ = v549
	var v556 int32
	_ = v556
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v565 int32
	_ = v565
	var v570 int32
	_ = v570
	var v572 int32
	_ = v572
	var v579 int32
	_ = v579
	var v582 int32
	_ = v582
	var v587 int32
	_ = v587
	var v594 int32
	_ = v594
	var v596 int32
	_ = v596
	var v603 int32
	_ = v603
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v618 int32
	_ = v618
	var v625 int32
	_ = v625
	var v630 int32
	_ = v630
	var v632 int32
	_ = v632
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v662 int32
	_ = v662
	var v665 int32
	_ = v665
	var v666 int32
	_ = v666
	var v671 int32
	_ = v671
	var v675 int32
	_ = v675
	var v680 int32
	_ = v680
	var v682 int32
	_ = v682
	var v689 int32
	_ = v689
	var v692 int32
	_ = v692
	var v694 int32
	_ = v694
	var v697 int32
	_ = v697
	var v702 int32
	_ = v702
	var v704 int32
	_ = v704
	var v711 int32
	_ = v711
	var v714 int32
	_ = v714
	var v719 int32
	_ = v719
	var v726 int32
	_ = v726
	var v728 int32
	_ = v728
	var v735 int32
	_ = v735
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v743 int32
	_ = v743
	var v745 int32
	_ = v745
	var v748 int32
	_ = v748
	var v749 int32
	_ = v749
	var v754 int32
	_ = v754
	var v758 int32
	_ = v758
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v772 int32
	_ = v772
	var v775 int32
	_ = v775
	var v777 int32
	_ = v777
	var v780 int32
	_ = v780
	var v785 int32
	_ = v785
	var v787 int32
	_ = v787
	var v794 int32
	_ = v794
	var v797 int32
	_ = v797
	var v802 int32
	_ = v802
	var v809 int32
	_ = v809
	var v811 int32
	_ = v811
	var v818 int32
	_ = v818
	var v821 int32
	_ = v821
	var v822 int32
	_ = v822
	v30 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v31 = int32(16)
	v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v34 = v30<<(uint(v31)%32) | v33
	v37 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v41 = v37<<(uint(v31)%32) | v40
	v44 = v34*int32(3) + v41 + int32(131074)
	v48 = int32(base.Ui32(v44)>>(uint(int32(2))%32)) & int32(255)
	v51 = int32(8)
	v53 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v57 = int32(base.Ui32(v53*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v51) % 32))
	v58 = int32(base.Ui32(v48*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v51)%32)) + v57
	v60 = v58 + int32(-17685)
	if base.Ui32(v58) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
		v67 = int32(0)
	} else {
		v67 = int32(240)
	}
	if base.Ui32(v60) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
		v70 = int32(base.Ui32(v60) >> (uint(int32(6)) % 32))
	} else {
		v70 = v67
	}
	v72 = v70 | int32(15)
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v72)
	v75 = int32(base.Ui32(v44) >> (uint(int32(18)) % 32))
	v80 = int32(base.Ui32(v75*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32)) + v57
	v82 = v80 + int32(-14234)
	if base.Ui32(v80) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
		v89 = int32(0)
	} else {
		v89 = int32(240)
	}
	if base.Ui32(v82) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
		v92 = int32(base.Ui32(v82) >> (uint(int32(6)) % 32))
	} else {
		v92 = v89
	}
	v97 = int32(8)
	v104 = v57 - (int32(base.Ui32(v75*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v97)%32)) + int32(base.Ui32(v48*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v97)%32)))
	v106 = v104 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
	if v104 < int32(-8708) {
		v113 = int32(0)
	} else {
		v113 = int32(15)
	}
	if base.Ui32(v106) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
		v116 = int32(base.Ui32(v106) >> (uint(int32(10)) % 32))
	} else {
		v116 = v113
	}
	v117 = v92&int32(240) | v116
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v117)
	v120 = l8 + int32(-1)
	v122 = v120 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v128 = int32(8)
		v129 = int32(base.Ui32(v125*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v128) % 32))
		v134 = v34 + v41*int32(3) + int32(131074)
		v138 = int32(base.Ui32(v134)>>(uint(int32(2))%32)) & int32(255)
		v143 = v129 + int32(base.Ui32(v138*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v128)%32))
		v145 = v143 + int32(-17685)
		if base.Ui32(v143) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
			v152 = int32(0)
		} else {
			v152 = int32(240)
		}
		if base.Ui32(v145) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
			v155 = int32(base.Ui32(v145) >> (uint(int32(6)) % 32))
		} else {
			v155 = v152
		}
		v157 = v155 | int32(15)
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v157)
		v160 = int32(base.Ui32(v134) >> (uint(int32(18)) % 32))
		v165 = v129 + int32(base.Ui32(v160*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32))
		v167 = v165 + int32(-14234)
		if base.Ui32(v165) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
			v174 = int32(0)
		} else {
			v174 = int32(240)
		}
		if base.Ui32(v167) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
			v177 = int32(base.Ui32(v167) >> (uint(int32(6)) % 32))
		} else {
			v177 = v174
		}
		v182 = int32(8)
		v189 = v129 - (int32(base.Ui32(v138*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v182)%32)) + int32(base.Ui32(v160*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v182)%32)))
		v191 = v189 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
		if v189 < int32(-8708) {
			v198 = int32(0)
		} else {
			v198 = int32(15)
		}
		if base.Ui32(v191) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
			v201 = int32(base.Ui32(v191) >> (uint(int32(10)) % 32))
		} else {
			v201 = v198
		}
		v202 = v177&int32(240) | v201
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v202)
	}
	if int32(1) <= v122 {
		v211 = int32(1)
		v223 = int32(0)
		v227 = v223
		v234 = v34
		v235 = v41
		v236 = l0 + v211
		v237 = l1 + v211
		v238 = v223
		for {
			v254 = l6 + v227
			v255 = int32(3)
			v259 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v211+v238))))
			v260 = int32(16)
			v263 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v211+v238))))
			v264 = v259<<(uint(v260)%32) | v263
			v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v211+v238))))
			v271 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v211+v238))))
			v272 = v267<<(uint(v260)%32) | v271
			v275 = v235 + v234 + v264 + v272 + int32(524296)
			v277 = int32(1)
			v281 = int32(base.Ui32(v275+(v264+v235)<<(uint(v277)%32)) >> (uint(v255) % 32))
			v282 = v281 + v234
			v286 = int32(base.Ui32(v282)>>(uint(v277)%32)) & int32(255)
			v289 = int32(8)
			v291 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236))))
			v295 = int32(base.Ui32(v291*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v289) % 32))
			v296 = int32(base.Ui32(v286*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v289)%32)) + v295
			v298 = v296 + int32(-17685)
			if base.Ui32(v296) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
				v305 = int32(0)
			} else {
				v305 = int32(240)
			}
			if base.Ui32(v298) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v308 = int32(base.Ui32(v298) >> (uint(int32(6)) % 32))
			} else {
				v308 = v305
			}
			v310 = v308 | int32(15)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v255))) = uint8(v310)
			v315 = int32(1)
			v319 = int32(base.Ui32(v275+(v272+v234)<<(uint(v315)%32)) >> (uint(int32(3)) % 32))
			v320 = v319 + v264
			v324 = int32(base.Ui32(v320)>>(uint(v315)%32)) & int32(255)
			v327 = int32(8)
			v331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236+v315))))
			v335 = int32(base.Ui32(v331*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v327) % 32))
			v336 = int32(base.Ui32(v324*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v327)%32)) + v335
			v338 = v336 + int32(-17685)
			if base.Ui32(v336) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
				v345 = int32(0)
			} else {
				v345 = int32(240)
			}
			if base.Ui32(v338) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v348 = int32(base.Ui32(v338) >> (uint(int32(6)) % 32))
			} else {
				v348 = v345
			}
			v350 = v348 | int32(15)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(5)))) = uint8(v350)
			v355 = int32(base.Ui32(v282) >> (uint(int32(17)) % 32))
			v360 = int32(base.Ui32(v355*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32)) + v295
			v362 = v360 + int32(-14234)
			if base.Ui32(v360) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
				v369 = int32(0)
			} else {
				v369 = int32(240)
			}
			if base.Ui32(v362) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v372 = int32(base.Ui32(v362) >> (uint(int32(6)) % 32))
			} else {
				v372 = v369
			}
			v377 = int32(8)
			v384 = v295 - (int32(base.Ui32(v355*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v377)%32)) + int32(base.Ui32(v286*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v377)%32)))
			v386 = v384 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
			if v384 < int32(-8708) {
				v393 = int32(0)
			} else {
				v393 = int32(15)
			}
			if base.Ui32(v386) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v396 = int32(base.Ui32(v386) >> (uint(int32(10)) % 32))
			} else {
				v396 = v393
			}
			v397 = v372&int32(240) | v396
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(2)))) = uint8(v397)
			v402 = int32(base.Ui32(v320) >> (uint(int32(17)) % 32))
			v407 = int32(base.Ui32(v402*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32)) + v335
			v409 = v407 + int32(-14234)
			if base.Ui32(v407) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
				v416 = int32(0)
			} else {
				v416 = int32(240)
			}
			if base.Ui32(v409) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v419 = int32(base.Ui32(v409) >> (uint(int32(6)) % 32))
			} else {
				v419 = v416
			}
			v424 = int32(8)
			v431 = v335 - (int32(base.Ui32(v402*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v424)%32)) + int32(base.Ui32(v324*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v424)%32)))
			v433 = v431 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
			if v431 < int32(-8708) {
				v440 = int32(0)
			} else {
				v440 = int32(15)
			}
			if base.Ui32(v433) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v443 = int32(base.Ui32(v433) >> (uint(int32(10)) % 32))
			} else {
				v443 = v440
			}
			v444 = v419&int32(240) | v443
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(4)))) = uint8(v444)
			if l1 == int32(0) {
			} else {
				v448 = l7 + v227
				v451 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237))))
				v454 = int32(8)
				v455 = int32(base.Ui32(v451*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v454) % 32))
				v456 = v319 + v235
				v460 = int32(base.Ui32(v456)>>(uint(int32(1))%32)) & int32(255)
				v465 = v455 + int32(base.Ui32(v460*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v454)%32))
				v467 = v465 + int32(-17685)
				if base.Ui32(v465) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
					v474 = int32(0)
				} else {
					v474 = int32(240)
				}
				if base.Ui32(v467) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
					v477 = int32(base.Ui32(v467) >> (uint(int32(6)) % 32))
				} else {
					v477 = v474
				}
				v479 = v477 | int32(15)
				*(*uint8)(unsafe.Add(mBase, uint32(v448+int32(3)))) = uint8(v479)
				v483 = int32(1)
				v485 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237+v483))))
				v488 = int32(8)
				v489 = int32(base.Ui32(v485*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v488) % 32))
				v490 = v281 + v272
				v494 = int32(base.Ui32(v490)>>(uint(v483)%32)) & int32(255)
				v499 = v489 + int32(base.Ui32(v494*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v488)%32))
				v501 = v499 + int32(-17685)
				if base.Ui32(v499) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
					v508 = int32(0)
				} else {
					v508 = int32(240)
				}
				if base.Ui32(v501) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
					v511 = int32(base.Ui32(v501) >> (uint(int32(6)) % 32))
				} else {
					v511 = v508
				}
				v513 = v511 | int32(15)
				*(*uint8)(unsafe.Add(mBase, uint32(v448+int32(5)))) = uint8(v513)
				v518 = int32(base.Ui32(v456) >> (uint(int32(17)) % 32))
				v523 = v455 + int32(base.Ui32(v518*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32))
				v525 = v523 + int32(-14234)
				if base.Ui32(v523) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
					v532 = int32(0)
				} else {
					v532 = int32(240)
				}
				if base.Ui32(v525) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
					v535 = int32(base.Ui32(v525) >> (uint(int32(6)) % 32))
				} else {
					v535 = v532
				}
				v540 = int32(8)
				v547 = v455 - (int32(base.Ui32(v460*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v540)%32)) + int32(base.Ui32(v518*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v540)%32)))
				v549 = v547 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
				if v547 < int32(-8708) {
					v556 = int32(0)
				} else {
					v556 = int32(15)
				}
				if base.Ui32(v549) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
					v559 = int32(base.Ui32(v549) >> (uint(int32(10)) % 32))
				} else {
					v559 = v556
				}
				v560 = v535&int32(240) | v559
				*(*uint8)(unsafe.Add(mBase, uint32(v448+int32(2)))) = uint8(v560)
				v565 = int32(base.Ui32(v490) >> (uint(int32(17)) % 32))
				v570 = v489 + int32(base.Ui32(v565*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32))
				v572 = v570 + int32(-14234)
				if base.Ui32(v570) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
					v579 = int32(0)
				} else {
					v579 = int32(240)
				}
				if base.Ui32(v572) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
					v582 = int32(base.Ui32(v572) >> (uint(int32(6)) % 32))
				} else {
					v582 = v579
				}
				v587 = int32(8)
				v594 = v489 - (int32(base.Ui32(v494*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v587)%32)) + int32(base.Ui32(v565*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v587)%32)))
				v596 = v594 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
				if v594 < int32(-8708) {
					v603 = int32(0)
				} else {
					v603 = int32(15)
				}
				if base.Ui32(v596) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
					v606 = int32(base.Ui32(v596) >> (uint(int32(10)) % 32))
				} else {
					v606 = v603
				}
				v607 = v582&int32(240) | v606
				*(*uint8)(unsafe.Add(mBase, uint32(v448+int32(4)))) = uint8(v607)
			}
			v618 = int32(2)
			v625 = v238 + int32(1)
			if v122 != v625 {
				v227 = v227 + int32(4)
				v234 = v264
				v235 = v272
				v236 = v236 + v618
				v237 = v237 + v618
				v238 = v625
				continue
			} else {
				break
			}
			break
		}
		v630 = v272
		v632 = v264
	} else {
		v630 = v41
		v632 = v34
	}
	if l8&int32(1) != 0 {
	} else {
		v659 = v120 << (uint(int32(1)) % 32)
		v660 = l6 + v659
		v662 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v120))))
		v665 = int32(8)
		v666 = int32(base.Ui32(v662*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v665) % 32))
		v671 = v630 + v632*int32(3) + int32(131074)
		v675 = int32(base.Ui32(v671)>>(uint(int32(2))%32)) & int32(255)
		v680 = v666 + int32(base.Ui32(v675*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v665)%32))
		v682 = v680 + int32(-17685)
		if base.Ui32(v680) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
			v689 = int32(0)
		} else {
			v689 = int32(240)
		}
		if base.Ui32(v682) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
			v692 = int32(base.Ui32(v682) >> (uint(int32(6)) % 32))
		} else {
			v692 = v689
		}
		v694 = v692 | int32(15)
		*(*uint8)(unsafe.Add(mBase, uint32(v660)+1)) = uint8(v694)
		v697 = int32(base.Ui32(v671) >> (uint(int32(18)) % 32))
		v702 = v666 + int32(base.Ui32(v697*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32))
		v704 = v702 + int32(-14234)
		if base.Ui32(v702) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
			v711 = int32(0)
		} else {
			v711 = int32(240)
		}
		if base.Ui32(v704) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
			v714 = int32(base.Ui32(v704) >> (uint(int32(6)) % 32))
		} else {
			v714 = v711
		}
		v719 = int32(8)
		v726 = v666 - (int32(base.Ui32(v675*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v719)%32)) + int32(base.Ui32(v697*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v719)%32)))
		v728 = v726 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
		if v726 < int32(-8708) {
			v735 = int32(0)
		} else {
			v735 = int32(15)
		}
		if base.Ui32(v728) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
			v738 = int32(base.Ui32(v728) >> (uint(int32(10)) % 32))
		} else {
			v738 = v735
		}
		v739 = v714&int32(240) | v738
		*(*uint8)(unsafe.Add(mBase, uint32(v660))) = uint8(v739)
		if l1 == int32(0) {
		} else {
			v743 = l7 + v659
			v745 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v120))))
			v748 = int32(8)
			v749 = int32(base.Ui32(v745*int32(_a_F_UpsampleRgba4444LinePair_C_0)) >> (uint(v748) % 32))
			v754 = v632 + v630*int32(3) + int32(131074)
			v758 = int32(base.Ui32(v754)>>(uint(int32(2))%32)) & int32(255)
			v763 = v749 + int32(base.Ui32(v758*int32(_a_F_UpsampleRgba4444LinePair_C_1))>>(uint(v748)%32))
			v765 = v763 + int32(-17685)
			if base.Ui32(v763) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_2)) {
				v772 = int32(0)
			} else {
				v772 = int32(240)
			}
			if base.Ui32(v765) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v775 = int32(base.Ui32(v765) >> (uint(int32(6)) % 32))
			} else {
				v775 = v772
			}
			v777 = v775 | int32(15)
			*(*uint8)(unsafe.Add(mBase, uint32(v743)+1)) = uint8(v777)
			v780 = int32(base.Ui32(v754) >> (uint(int32(18)) % 32))
			v785 = v749 + int32(base.Ui32(v780*int32(_a_F_UpsampleRgba4444LinePair_C_4))>>(uint(int32(8))%32))
			v787 = v785 + int32(-14234)
			if base.Ui32(v785) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_5)) {
				v794 = int32(0)
			} else {
				v794 = int32(240)
			}
			if base.Ui32(v787) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v797 = int32(base.Ui32(v787) >> (uint(int32(6)) % 32))
			} else {
				v797 = v794
			}
			v802 = int32(8)
			v809 = v749 - (int32(base.Ui32(v758*int32(_a_F_UpsampleRgba4444LinePair_C_7))>>(uint(v802)%32)) + int32(base.Ui32(v780*int32(_a_F_UpsampleRgba4444LinePair_C_6))>>(uint(v802)%32)))
			v811 = v809 + int32(_a_F_UpsampleRgba4444LinePair_C_8)
			if v809 < int32(-8708) {
				v818 = int32(0)
			} else {
				v818 = int32(15)
			}
			if base.Ui32(v811) < base.Ui32(int32(_a_F_UpsampleRgba4444LinePair_C_3)) {
				v821 = int32(base.Ui32(v811) >> (uint(int32(10)) % 32))
			} else {
				v821 = v818
			}
			v822 = v797&int32(240) | v821
			*(*uint8)(unsafe.Add(mBase, uint32(v743))) = uint8(v822)
		}
	}
	return
}
func F_UpsampleRgbaLinePair_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v30 int32
	_ = v30
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v82 int32
	_ = v82
	var v84 int32
	_ = v84
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v140 int32
	_ = v140
	var v145 int32
	_ = v145
	var v147 int32
	_ = v147
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v211 int32
	_ = v211
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v230 int32
	_ = v230
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v263 int32
	_ = v263
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v275 int32
	_ = v275
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v285 int32
	_ = v285
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v294 int32
	_ = v294
	var v297 int32
	_ = v297
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v313 int32
	_ = v313
	var v316 int32
	_ = v316
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v330 int32
	_ = v330
	var v337 int32
	_ = v337
	var v340 int32
	_ = v340
	var v345 int32
	_ = v345
	var v349 int32
	_ = v349
	var v350 int32
	_ = v350
	var v353 int32
	_ = v353
	var v354 int32
	_ = v354
	var v357 int32
	_ = v357
	var v361 int32
	_ = v361
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v368 int32
	_ = v368
	var v375 int32
	_ = v375
	var v378 int32
	_ = v378
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v390 int32
	_ = v390
	var v392 int32
	_ = v392
	var v399 int32
	_ = v399
	var v402 int32
	_ = v402
	var v408 int32
	_ = v408
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v433 int32
	_ = v433
	var v440 int32
	_ = v440
	var v442 int32
	_ = v442
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v456 int32
	_ = v456
	var v459 int32
	_ = v459
	var v465 int32
	_ = v465
	var v467 int32
	_ = v467
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v472 int32
	_ = v472
	var v476 int32
	_ = v476
	var v481 int32
	_ = v481
	var v483 int32
	_ = v483
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v514 int32
	_ = v514
	var v517 int32
	_ = v517
	var v521 int32
	_ = v521
	var v523 int32
	_ = v523
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v531 int32
	_ = v531
	var v532 int32
	_ = v532
	var v537 int32
	_ = v537
	var v539 int32
	_ = v539
	var v546 int32
	_ = v546
	var v549 int32
	_ = v549
	var v551 int32
	_ = v551
	var v555 int32
	_ = v555
	var v556 int32
	_ = v556
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v570 int32
	_ = v570
	var v573 int32
	_ = v573
	var v579 int32
	_ = v579
	var v586 int32
	_ = v586
	var v588 int32
	_ = v588
	var v595 int32
	_ = v595
	var v598 int32
	_ = v598
	var v604 int32
	_ = v604
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v620 int32
	_ = v620
	var v623 int32
	_ = v623
	var v634 int32
	_ = v634
	var v641 int32
	_ = v641
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v680 int32
	_ = v680
	var v683 int32
	_ = v683
	var v684 int32
	_ = v684
	var v689 int32
	_ = v689
	var v693 int32
	_ = v693
	var v698 int32
	_ = v698
	var v700 int32
	_ = v700
	var v707 int32
	_ = v707
	var v710 int32
	_ = v710
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v720 int32
	_ = v720
	var v722 int32
	_ = v722
	var v729 int32
	_ = v729
	var v732 int32
	_ = v732
	var v736 int32
	_ = v736
	var v743 int32
	_ = v743
	var v745 int32
	_ = v745
	var v752 int32
	_ = v752
	var v755 int32
	_ = v755
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v763 int32
	_ = v763
	var v766 int32
	_ = v766
	var v767 int32
	_ = v767
	var v772 int32
	_ = v772
	var v776 int32
	_ = v776
	var v781 int32
	_ = v781
	var v783 int32
	_ = v783
	var v790 int32
	_ = v790
	var v793 int32
	_ = v793
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v803 int32
	_ = v803
	var v805 int32
	_ = v805
	var v812 int32
	_ = v812
	var v815 int32
	_ = v815
	var v819 int32
	_ = v819
	var v826 int32
	_ = v826
	var v828 int32
	_ = v828
	var v835 int32
	_ = v835
	var v838 int32
	_ = v838
	v30 = int32(255)
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+3)) = uint8(v30)
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v33 = int32(16)
	v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v36 = v32<<(uint(v33)%32) | v35
	v39 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v43 = v39<<(uint(v33)%32) | v42
	v46 = v36*int32(3) + v43 + int32(131074)
	v50 = int32(base.Ui32(v46)>>(uint(int32(2))%32)) & v30
	v53 = int32(8)
	v55 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v59 = int32(base.Ui32(v55*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v53) % 32))
	v60 = int32(base.Ui32(v50*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v53)%32)) + v59
	v62 = v60 + int32(-17685)
	if base.Ui32(v60) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
		v69 = int32(0)
	} else {
		v69 = v30
	}
	if base.Ui32(v62) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
		v72 = int32(base.Ui32(v62) >> (uint(int32(6)) % 32))
	} else {
		v72 = v69
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v72)
	v76 = int32(255)
	v77 = int32(base.Ui32(v46)>>(uint(int32(18))%32)) & v76
	v82 = int32(base.Ui32(v77*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(int32(8))%32)) + v59
	v84 = v82 + int32(-14234)
	if base.Ui32(v82) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
		v91 = int32(0)
	} else {
		v91 = v76
	}
	if base.Ui32(v84) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
		v94 = int32(base.Ui32(v84) >> (uint(int32(6)) % 32))
	} else {
		v94 = v91
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v94)
	v98 = int32(8)
	v105 = v59 - (int32(base.Ui32(v77*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v98)%32)) + int32(base.Ui32(v50*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v98)%32)))
	v107 = v105 + int32(_a_F_UpsampleRgbaLinePair_C_8)
	if v105 < int32(-8708) {
		v114 = int32(0)
	} else {
		v114 = int32(255)
	}
	if base.Ui32(v107) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
		v117 = int32(base.Ui32(v107) >> (uint(int32(6)) % 32))
	} else {
		v117 = v114
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v117)
	v120 = l8 + int32(-1)
	v122 = v120 >> (uint(int32(1)) % 32)
	if l1 == int32(0) {
	} else {
		v125 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+3)) = uint8(v125)
		v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v130 = int32(8)
		v131 = int32(base.Ui32(v127*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v130) % 32))
		v136 = v36 + v43*int32(3) + int32(131074)
		v140 = int32(base.Ui32(v136)>>(uint(int32(2))%32)) & v125
		v145 = v131 + int32(base.Ui32(v140*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v130)%32))
		v147 = v145 + int32(-17685)
		if base.Ui32(v145) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
			v154 = int32(0)
		} else {
			v154 = v125
		}
		if base.Ui32(v147) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
			v157 = int32(base.Ui32(v147) >> (uint(int32(6)) % 32))
		} else {
			v157 = v154
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v157)
		v161 = int32(255)
		v162 = int32(base.Ui32(v136)>>(uint(int32(18))%32)) & v161
		v167 = v131 + int32(base.Ui32(v162*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(int32(8))%32))
		v169 = v167 + int32(-14234)
		if base.Ui32(v167) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
			v176 = int32(0)
		} else {
			v176 = v161
		}
		if base.Ui32(v169) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
			v179 = int32(base.Ui32(v169) >> (uint(int32(6)) % 32))
		} else {
			v179 = v176
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v179)
		v183 = int32(8)
		v190 = v131 - (int32(base.Ui32(v140*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v183)%32)) + int32(base.Ui32(v162*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v183)%32)))
		v192 = v190 + int32(_a_F_UpsampleRgbaLinePair_C_8)
		if v190 < int32(-8708) {
			v199 = int32(0)
		} else {
			v199 = int32(255)
		}
		if base.Ui32(v192) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
			v202 = int32(base.Ui32(v192) >> (uint(int32(6)) % 32))
		} else {
			v202 = v199
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v202)
	}
	if int32(1) <= v122 {
		v211 = int32(1)
		v223 = int32(0)
		v227 = v223
		v230 = v223
		v234 = v36
		v235 = v43
		v236 = l0 + v211
		v237 = l1 + v211
		for {
			v254 = l6 + v227
			v257 = int32(255)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(7)))) = uint8(v257)
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(11)))) = uint8(v257)
			v263 = int32(6)
			v267 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3+v211+v230))))
			v268 = int32(16)
			v271 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v211+v230))))
			v272 = v267<<(uint(v268)%32) | v271
			v275 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5+v211+v230))))
			v279 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4+v211+v230))))
			v280 = v275<<(uint(v268)%32) | v279
			v283 = v235 + v234 + v272 + v280 + int32(524296)
			v285 = int32(1)
			v289 = int32(base.Ui32(v283+(v272+v235)<<(uint(v285)%32)) >> (uint(int32(3)) % 32))
			v290 = v289 + v234
			v294 = int32(base.Ui32(v290)>>(uint(v285)%32)) & v257
			v297 = int32(8)
			v299 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236))))
			v303 = int32(base.Ui32(v299*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v297) % 32))
			v304 = int32(base.Ui32(v294*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v297)%32)) + v303
			v306 = v304 + int32(-17685)
			if base.Ui32(v304) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
				v313 = int32(0)
			} else {
				v313 = v257
			}
			if base.Ui32(v306) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v316 = int32(base.Ui32(v306) >> (uint(v263) % 32))
			} else {
				v316 = v313
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v263))) = uint8(v316)
			v322 = int32(255)
			v323 = int32(base.Ui32(v290)>>(uint(int32(17))%32)) & v322
			v328 = int32(base.Ui32(v323*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(int32(8))%32)) + v303
			v330 = v328 + int32(-14234)
			if base.Ui32(v328) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
				v337 = int32(0)
			} else {
				v337 = v322
			}
			if base.Ui32(v330) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v340 = int32(base.Ui32(v330) >> (uint(int32(6)) % 32))
			} else {
				v340 = v337
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(4)))) = uint8(v340)
			v345 = int32(1)
			v349 = int32(base.Ui32(v283+(v280+v234)<<(uint(v345)%32)) >> (uint(int32(3)) % 32))
			v350 = v349 + v272
			v353 = int32(255)
			v354 = int32(base.Ui32(v350)>>(uint(v345)%32)) & v353
			v357 = int32(8)
			v361 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v236+v345))))
			v365 = int32(base.Ui32(v361*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v357) % 32))
			v366 = int32(base.Ui32(v354*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v357)%32)) + v365
			v368 = v366 + int32(-17685)
			if base.Ui32(v366) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
				v375 = int32(0)
			} else {
				v375 = v353
			}
			if base.Ui32(v368) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v378 = int32(base.Ui32(v368) >> (uint(int32(6)) % 32))
			} else {
				v378 = v375
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(10)))) = uint8(v378)
			v380 = int32(8)
			v384 = int32(255)
			v385 = int32(base.Ui32(v350)>>(uint(int32(17))%32)) & v384
			v390 = int32(base.Ui32(v385*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(v380)%32)) + v365
			v392 = v390 + int32(-14234)
			if base.Ui32(v390) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
				v399 = int32(0)
			} else {
				v399 = v384
			}
			if base.Ui32(v392) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v402 = int32(base.Ui32(v392) >> (uint(int32(6)) % 32))
			} else {
				v402 = v399
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+v380))) = uint8(v402)
			v408 = int32(8)
			v415 = v303 - (int32(base.Ui32(v323*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v408)%32)) + int32(base.Ui32(v294*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v408)%32)))
			v417 = v415 + int32(_a_F_UpsampleRgbaLinePair_C_8)
			if v415 < int32(-8708) {
				v424 = int32(0)
			} else {
				v424 = int32(255)
			}
			if base.Ui32(v417) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v427 = int32(base.Ui32(v417) >> (uint(int32(6)) % 32))
			} else {
				v427 = v424
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(5)))) = uint8(v427)
			v433 = int32(8)
			v440 = v365 - (int32(base.Ui32(v385*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v433)%32)) + int32(base.Ui32(v354*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v433)%32)))
			v442 = v440 + int32(_a_F_UpsampleRgbaLinePair_C_8)
			if v440 < int32(-8708) {
				v449 = int32(0)
			} else {
				v449 = int32(255)
			}
			if base.Ui32(v442) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v452 = int32(base.Ui32(v442) >> (uint(int32(6)) % 32))
			} else {
				v452 = v449
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v254+int32(9)))) = uint8(v452)
			if l1 == int32(0) {
			} else {
				v456 = l7 + v227
				v459 = int32(255)
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(7)))) = uint8(v459)
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(11)))) = uint8(v459)
				v465 = int32(6)
				v467 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237))))
				v470 = int32(8)
				v471 = int32(base.Ui32(v467*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v470) % 32))
				v472 = v349 + v235
				v476 = int32(base.Ui32(v472)>>(uint(int32(1))%32)) & v459
				v481 = v471 + int32(base.Ui32(v476*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v470)%32))
				v483 = v481 + int32(-17685)
				if base.Ui32(v481) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
					v490 = int32(0)
				} else {
					v490 = v459
				}
				if base.Ui32(v483) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
					v493 = int32(base.Ui32(v483) >> (uint(v465) % 32))
				} else {
					v493 = v490
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+v465))) = uint8(v493)
				v499 = int32(255)
				v500 = int32(base.Ui32(v472)>>(uint(int32(17))%32)) & v499
				v505 = v471 + int32(base.Ui32(v500*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(int32(8))%32))
				v507 = v505 + int32(-14234)
				if base.Ui32(v505) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
					v514 = int32(0)
				} else {
					v514 = v499
				}
				if base.Ui32(v507) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
					v517 = int32(base.Ui32(v507) >> (uint(int32(6)) % 32))
				} else {
					v517 = v514
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(4)))) = uint8(v517)
				v521 = int32(1)
				v523 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v237+v521))))
				v526 = int32(8)
				v527 = int32(base.Ui32(v523*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v526) % 32))
				v528 = v289 + v280
				v531 = int32(255)
				v532 = int32(base.Ui32(v528)>>(uint(v521)%32)) & v531
				v537 = v527 + int32(base.Ui32(v532*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v526)%32))
				v539 = v537 + int32(-17685)
				if base.Ui32(v537) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
					v546 = int32(0)
				} else {
					v546 = v531
				}
				if base.Ui32(v539) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
					v549 = int32(base.Ui32(v539) >> (uint(int32(6)) % 32))
				} else {
					v549 = v546
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(10)))) = uint8(v549)
				v551 = int32(8)
				v555 = int32(255)
				v556 = int32(base.Ui32(v528)>>(uint(int32(17))%32)) & v555
				v561 = v527 + int32(base.Ui32(v556*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(v551)%32))
				v563 = v561 + int32(-14234)
				if base.Ui32(v561) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
					v570 = int32(0)
				} else {
					v570 = v555
				}
				if base.Ui32(v563) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
					v573 = int32(base.Ui32(v563) >> (uint(int32(6)) % 32))
				} else {
					v573 = v570
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+v551))) = uint8(v573)
				v579 = int32(8)
				v586 = v471 - (int32(base.Ui32(v476*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v579)%32)) + int32(base.Ui32(v500*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v579)%32)))
				v588 = v586 + int32(_a_F_UpsampleRgbaLinePair_C_8)
				if v586 < int32(-8708) {
					v595 = int32(0)
				} else {
					v595 = int32(255)
				}
				if base.Ui32(v588) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
					v598 = int32(base.Ui32(v588) >> (uint(int32(6)) % 32))
				} else {
					v598 = v595
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(5)))) = uint8(v598)
				v604 = int32(8)
				v611 = v527 - (int32(base.Ui32(v532*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v604)%32)) + int32(base.Ui32(v556*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v604)%32)))
				v613 = v611 + int32(_a_F_UpsampleRgbaLinePair_C_8)
				if v611 < int32(-8708) {
					v620 = int32(0)
				} else {
					v620 = int32(255)
				}
				if base.Ui32(v613) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
					v623 = int32(base.Ui32(v613) >> (uint(int32(6)) % 32))
				} else {
					v623 = v620
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v456+int32(9)))) = uint8(v623)
			}
			v634 = int32(2)
			v641 = v230 + int32(1)
			if v122 != v641 {
				v227 = v227 + int32(8)
				v230 = v641
				v234 = v272
				v235 = v280
				v236 = v236 + v634
				v237 = v237 + v634
				continue
			} else {
				break
			}
			break
		}
		v646 = v280
		v647 = v272
	} else {
		v646 = v43
		v647 = v36
	}
	if l8&int32(1) != 0 {
	} else {
		v674 = int32(2)
		v675 = v120 << (uint(v674) % 32)
		v676 = l6 + v675
		v677 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+3)) = uint8(v677)
		v680 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v120))))
		v683 = int32(8)
		v684 = int32(base.Ui32(v680*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v683) % 32))
		v689 = v646 + v647*int32(3) + int32(131074)
		v693 = int32(base.Ui32(v689)>>(uint(v674)%32)) & v677
		v698 = v684 + int32(base.Ui32(v693*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v683)%32))
		v700 = v698 + int32(-17685)
		if base.Ui32(v698) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
			v707 = int32(0)
		} else {
			v707 = v677
		}
		if base.Ui32(v700) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
			v710 = int32(base.Ui32(v700) >> (uint(int32(6)) % 32))
		} else {
			v710 = v707
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+2)) = uint8(v710)
		v714 = int32(255)
		v715 = int32(base.Ui32(v689)>>(uint(int32(18))%32)) & v714
		v720 = v684 + int32(base.Ui32(v715*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(int32(8))%32))
		v722 = v720 + int32(-14234)
		if base.Ui32(v720) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
			v729 = int32(0)
		} else {
			v729 = v714
		}
		if base.Ui32(v722) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
			v732 = int32(base.Ui32(v722) >> (uint(int32(6)) % 32))
		} else {
			v732 = v729
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676))) = uint8(v732)
		v736 = int32(8)
		v743 = v684 - (int32(base.Ui32(v693*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v736)%32)) + int32(base.Ui32(v715*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v736)%32)))
		v745 = v743 + int32(_a_F_UpsampleRgbaLinePair_C_8)
		if v743 < int32(-8708) {
			v752 = int32(0)
		} else {
			v752 = int32(255)
		}
		if base.Ui32(v745) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
			v755 = int32(base.Ui32(v745) >> (uint(int32(6)) % 32))
		} else {
			v755 = v752
		}
		*(*uint8)(unsafe.Add(mBase, uint32(v676)+1)) = uint8(v755)
		if l1 == int32(0) {
		} else {
			v759 = l7 + v675
			v760 = int32(255)
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+3)) = uint8(v760)
			v763 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v120))))
			v766 = int32(8)
			v767 = int32(base.Ui32(v763*int32(_a_F_UpsampleRgbaLinePair_C_0)) >> (uint(v766) % 32))
			v772 = v647 + v646*int32(3) + int32(131074)
			v776 = int32(base.Ui32(v772)>>(uint(int32(2))%32)) & v760
			v781 = v767 + int32(base.Ui32(v776*int32(_a_F_UpsampleRgbaLinePair_C_1))>>(uint(v766)%32))
			v783 = v781 + int32(-17685)
			if base.Ui32(v781) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_2)) {
				v790 = int32(0)
			} else {
				v790 = v760
			}
			if base.Ui32(v783) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v793 = int32(base.Ui32(v783) >> (uint(int32(6)) % 32))
			} else {
				v793 = v790
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+2)) = uint8(v793)
			v797 = int32(255)
			v798 = int32(base.Ui32(v772)>>(uint(int32(18))%32)) & v797
			v803 = v767 + int32(base.Ui32(v798*int32(_a_F_UpsampleRgbaLinePair_C_4))>>(uint(int32(8))%32))
			v805 = v803 + int32(-14234)
			if base.Ui32(v803) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_5)) {
				v812 = int32(0)
			} else {
				v812 = v797
			}
			if base.Ui32(v805) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v815 = int32(base.Ui32(v805) >> (uint(int32(6)) % 32))
			} else {
				v815 = v812
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759))) = uint8(v815)
			v819 = int32(8)
			v826 = v767 - (int32(base.Ui32(v776*int32(_a_F_UpsampleRgbaLinePair_C_7))>>(uint(v819)%32)) + int32(base.Ui32(v798*int32(_a_F_UpsampleRgbaLinePair_C_6))>>(uint(v819)%32)))
			v828 = v826 + int32(_a_F_UpsampleRgbaLinePair_C_8)
			if v826 < int32(-8708) {
				v835 = int32(0)
			} else {
				v835 = int32(255)
			}
			if base.Ui32(v828) < base.Ui32(int32(_a_F_UpsampleRgbaLinePair_C_3)) {
				v838 = int32(base.Ui32(v828) >> (uint(int32(6)) % 32))
			} else {
				v838 = v835
			}
			*(*uint8)(unsafe.Add(mBase, uint32(v759)+1)) = uint8(v838)
		}
	}
	return
}

//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

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

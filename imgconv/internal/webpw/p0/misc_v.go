//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_VE16_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v8 int64
	_ = v8
	var v20 int32
	_ = v20
	var v21 int64
	_ = v21
	var v46 int64
	_ = v46
	var v48 int64
	_ = v48
	v7 = l0 + int32(-32)
	v8 = *(*int64)(unsafe.Add(mBase, uint32(v7)))
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+96)) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+128)) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+160)) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+192)) = v8
	*(*int64)(unsafe.Add(mBase, uint32(l0)+224)) = v8
	v20 = l0 + int32(-24)
	v21 = *(*int64)(unsafe.Add(mBase, uint32(v20)))
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(8)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(72)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(104)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(136)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(168)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(200)))) = v21
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(232)))) = v21
	v46 = *(*int64)(unsafe.Add(mBase, uint32(v20)))
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(264)))) = v46
	v48 = *(*int64)(unsafe.Add(mBase, uint32(v7)))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+256)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(296)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0)+288)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(328)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0)+320)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(360)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0)+352)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0)+384)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(392)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(424)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0)+416)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0)+448)) = v48
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(456)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(488)))) = v46
	*(*int64)(unsafe.Add(mBase, uint32(l0)+480)) = v48
	return
}
func F_VE4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v56 int32
	_ = v56
	v10 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-30)))))
	v11 = int32(2)
	v12 = v10 + v11
	v15 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-29)))))
	v16 = int32(1)
	v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-28)))))
	v24 = int32(base.Ui32(v12+v15<<(uint(v16)%32)+v21) >> (uint(v11) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v24)
	v28 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-31)))))
	v30 = v28 + v11
	v36 = int32(base.Ui32(v15+(v30+v10<<(uint(v16)%32))) >> (uint(v11) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v36)
	v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-32)))))
	v46 = int32(base.Ui32(v12+v40+v28<<(uint(v16)%32)) >> (uint(v11) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v46)
	v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	v56 = int32(base.Ui32(v30+v50+v40<<(uint(v16)%32)) >> (uint(v11) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v56)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v24)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v36)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v46)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v56)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v24)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v36)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v46)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v56)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v24)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v36)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v46)
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v56)
	return
}
func F_VE8uv_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int64
	_ = v5
	v5 = *(*int64)(unsafe.Add(mBase, uint32(l0+int32(-32))))
	*(*int64)(unsafe.Add(mBase, uint32(l0)+224)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+192)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+160)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+128)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+96)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = v5
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v5
	return
}
func F_VFilter16_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v35 int32
	_ = v35
	var v43 int32
	_ = v43
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
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
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v103 int32
	_ = v103
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
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v164 int32
	_ = v164
	var v166 int32
	_ = v166
	var v171 int32
	_ = v171
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v197 int32
	_ = v197
	var v200 int32
	_ = v200
	var v204 int32
	_ = v204
	var v208 int32
	_ = v208
	var v211 int32
	_ = v211
	var v216 int32
	_ = v216
	var v218 int32
	_ = v218
	var v232 int32
	_ = v232
	v35 = int32(1)
	v43 = int32(0)
	v53 = l1 << (uint(v35) % 32)
	v55 = m.G13
	v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)))
	v57 = m.G12
	v58 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
	v59 = m.G14
	v60 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
	v61 = m.G15
	v62 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
	v63 = l0
	v67 = int32(17)
	for {
		v95 = v63 + (v43 - l1)
		v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v95))))
		v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63))))
		v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v96-v97)))))
		v103 = v63 + (v43 - v53)
		v104 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103))))
		v105 = v63 + l1
		v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v105))))
		v107 = v104 - v106
		v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+v107))))
		if l2<<(uint(v35)%32)|v35 < v100<<(uint(int32(2))%32)+v109 {
		} else {
			v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v43-l1<<(uint(int32(2))%32))))))
			v114 = v63 + l1*int32(-3)
			v115 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v114))))
			v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v113-v115)))))
			if l3 < v118 {
			} else {
				v122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v115-v104)))))
				if base.Ui32(l3) < base.Ui32(v122) {
				} else {
					v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v104-v96)))))
					if base.Ui32(l3) < base.Ui32(v126) {
					} else {
						v129 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+l1*int32(3)))))
						v130 = v63 + v53
						v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v130))))
						v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v129-v131)))))
						if base.Ui32(l3) < base.Ui32(v134) {
						} else {
							v138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v131-v106)))))
							if base.Ui32(l3) < base.Ui32(v138) {
							} else {
								v142 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v106-v97)))))
								if base.Ui32(l3) < base.Ui32(v142) {
								} else {
									v146 = (v97 - v96) * int32(3)
									v148 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+v107))))
									if l4 < v126 {
										v151 = v146 + v148
										v154 = int32(3)
										v157 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+int32(4))>>(uint(v154)%32)))))
										v164 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+v154)>>(uint(v154)%32)))))
										v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v96+v164))))
										*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v166)
										v211 = v97 - v157
										v216 = v63
									} else {
										if base.Ui32(v142) <= base.Ui32(l4) {
											v171 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+v146+v148))))
											v174 = int32(63)
											v176 = int32(7)
											v177 = (v171*int32(9) + v174) >> (uint(v176) % 32)
											v180 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v177+v115))))
											*(*uint8)(unsafe.Add(mBase, uint32(v114))) = uint8(v180)
											v187 = (v171*int32(18) + v174) >> (uint(v176) % 32)
											v190 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v187+v104))))
											*(*uint8)(unsafe.Add(mBase, uint32(v103))) = uint8(v190)
											v197 = (v171*int32(27) + v174) >> (uint(v176) % 32)
											v200 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v197+v96))))
											*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v200)
											v204 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v97-v197)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v63))) = uint8(v204)
											v208 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v106-v187)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v105))) = uint8(v208)
											v211 = v131 - v177
											v216 = v130
										} else {
											v151 = v146 + v148
											v154 = int32(3)
											v157 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+int32(4))>>(uint(v154)%32)))))
											v164 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+(v151+v154)>>(uint(v154)%32)))))
											v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v96+v164))))
											*(*uint8)(unsafe.Add(mBase, uint32(v95))) = uint8(v166)
											v211 = v97 - v157
											v216 = v63
										}
									}
									v218 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v211))))
									*(*uint8)(unsafe.Add(mBase, uint32(v216))) = uint8(v218)
								}
							}
						}
					}
				}
			}
		}
		v232 = v67 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v232) {
			v63 = v63 + int32(1)
			v67 = v232
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_VFilter16i_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v35 int32
	_ = v35
	var v43 int32
	_ = v43
	var v53 int32
	_ = v53
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
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
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v118 int32
	_ = v118
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
	var v138 int32
	_ = v138
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v160 int32
	_ = v160
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v183 int32
	_ = v183
	var v187 int32
	_ = v187
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v201 int32
	_ = v201
	var v212 int32
	_ = v212
	var v243 int32
	_ = v243
	var v269 int32
	_ = v269
	var v277 int32
	_ = v277
	var v287 int32
	_ = v287
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v297 int32
	_ = v297
	var v301 int32
	_ = v301
	var v325 int32
	_ = v325
	var v326 int32
	_ = v326
	var v327 int32
	_ = v327
	var v330 int32
	_ = v330
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v337 int32
	_ = v337
	var v339 int32
	_ = v339
	var v343 int32
	_ = v343
	var v345 int32
	_ = v345
	var v348 int32
	_ = v348
	var v352 int32
	_ = v352
	var v356 int32
	_ = v356
	var v359 int32
	_ = v359
	var v361 int32
	_ = v361
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v394 int32
	_ = v394
	var v396 int32
	_ = v396
	var v399 int32
	_ = v399
	var v404 int32
	_ = v404
	var v410 int32
	_ = v410
	var v411 int32
	_ = v411
	var v414 int32
	_ = v414
	var v417 int32
	_ = v417
	var v421 int32
	_ = v421
	var v425 int32
	_ = v425
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v435 int32
	_ = v435
	var v446 int32
	_ = v446
	var v503 int32
	_ = v503
	var v511 int32
	_ = v511
	var v521 int32
	_ = v521
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v525 int32
	_ = v525
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v529 int32
	_ = v529
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v535 int32
	_ = v535
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v564 int32
	_ = v564
	var v567 int32
	_ = v567
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v573 int32
	_ = v573
	var v577 int32
	_ = v577
	var v579 int32
	_ = v579
	var v582 int32
	_ = v582
	var v586 int32
	_ = v586
	var v590 int32
	_ = v590
	var v593 int32
	_ = v593
	var v595 int32
	_ = v595
	var v598 int32
	_ = v598
	var v602 int32
	_ = v602
	var v606 int32
	_ = v606
	var v610 int32
	_ = v610
	var v614 int32
	_ = v614
	var v615 int32
	_ = v615
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v628 int32
	_ = v628
	var v630 int32
	_ = v630
	var v633 int32
	_ = v633
	var v638 int32
	_ = v638
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v648 int32
	_ = v648
	var v651 int32
	_ = v651
	var v655 int32
	_ = v655
	var v659 int32
	_ = v659
	var v662 int32
	_ = v662
	var v663 int32
	_ = v663
	var v669 int32
	_ = v669
	var v680 int32
	_ = v680
	v8 = l1 << (uint(int32(2)) % 32)
	v9 = l0 + v8
	v35 = int32(1)
	v43 = int32(0)
	v53 = l1 << (uint(v35) % 32)
	v55 = m.G14
	v56 = *(*int32)(unsafe.Add(mBase, uint32(v55)))
	v57 = m.G12
	v58 = *(*int32)(unsafe.Add(mBase, uint32(v57)))
	v59 = m.G13
	v60 = *(*int32)(unsafe.Add(mBase, uint32(v59)))
	v61 = m.G15
	v62 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
	v63 = v9
	v67 = int32(17)
	for {
		v91 = v63 + (v43 - l1)
		v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
		v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63))))
		v96 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v92-v93)))))
		v99 = v63 + (v43 - v53)
		v100 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v99))))
		v101 = v63 + l1
		v102 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v101))))
		v103 = v100 - v102
		v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+v103))))
		if l2<<(uint(v35)%32)|v35 < v96<<(uint(int32(2))%32)+v105 {
		} else {
			v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v43-l1<<(uint(int32(2))%32))))))
			v111 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+l1*int32(-3)))))
			v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v109-v111)))))
			if l3 < v114 {
			} else {
				v118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v111-v100)))))
				if base.Ui32(l3) < base.Ui32(v118) {
				} else {
					v122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v100-v92)))))
					if base.Ui32(l3) < base.Ui32(v122) {
					} else {
						v125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+l1*int32(3)))))
						v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+v53))))
						v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v125-v127)))))
						if base.Ui32(l3) < base.Ui32(v130) {
						} else {
							v134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v127-v102)))))
							if base.Ui32(l3) < base.Ui32(v134) {
							} else {
								v138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v62+(v102-v93)))))
								if base.Ui32(l3) < base.Ui32(v138) {
								} else {
									v142 = (v93 - v92) * int32(3)
									if l4 < v122 {
										v146 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+v103))))
										v147 = v142 + v146
										v150 = int32(3)
										v153 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+(v147+int32(4))>>(uint(v150)%32)))))
										v160 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+(v147+v150)>>(uint(v150)%32)))))
										v162 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v92+v160))))
										*(*uint8)(unsafe.Add(mBase, uint32(v91))) = uint8(v162)
										v194 = v93 - v153
										v195 = v63
									} else {
										if base.Ui32(v138) <= base.Ui32(l4) {
											v165 = int32(3)
											v170 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+(v142+v165)>>(uint(v165)%32)))))
											v176 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+(v142+int32(4))>>(uint(v165)%32)))))
											v177 = int32(1)
											v180 = (v176 + v177) >> (uint(v177) % 32)
											v183 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v180+v100))))
											*(*uint8)(unsafe.Add(mBase, uint32(v99))) = uint8(v183)
											v187 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v170+(v58+v92)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v91))) = uint8(v187)
											v191 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+(v93-v176)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v63))) = uint8(v191)
											v194 = v102 - v180
											v195 = v101
										} else {
											v146 = int32(*(*int8)(unsafe.Add(mBase, uint32(v56+v103))))
											v147 = v142 + v146
											v150 = int32(3)
											v153 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+(v147+int32(4))>>(uint(v150)%32)))))
											v160 = int32(*(*int8)(unsafe.Add(mBase, uint32(v60+(v147+v150)>>(uint(v150)%32)))))
											v162 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v92+v160))))
											*(*uint8)(unsafe.Add(mBase, uint32(v91))) = uint8(v162)
											v194 = v93 - v153
											v195 = v63
										}
									}
									v201 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v58+v194))))
									*(*uint8)(unsafe.Add(mBase, uint32(v195))) = uint8(v201)
								}
							}
						}
					}
				}
			}
		}
		v212 = v67 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v212) {
			v63 = v63 + int32(1)
			v67 = v212
			continue
		} else {
			break
		}
		break
	}
	v243 = v9 + v8
	v269 = int32(1)
	v277 = int32(0)
	v287 = l1 << (uint(v269) % 32)
	v289 = m.G14
	v290 = *(*int32)(unsafe.Add(mBase, uint32(v289)))
	v291 = m.G12
	v292 = *(*int32)(unsafe.Add(mBase, uint32(v291)))
	v293 = m.G13
	v294 = *(*int32)(unsafe.Add(mBase, uint32(v293)))
	v295 = m.G15
	v296 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v297 = v243
	v301 = int32(17)
	for {
		v325 = v297 + (v277 - l1)
		v326 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325))))
		v327 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297))))
		v330 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v326-v327)))))
		v333 = v297 + (v277 - v287)
		v334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v333))))
		v335 = v297 + l1
		v336 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v335))))
		v337 = v334 - v336
		v339 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+v337))))
		if l2<<(uint(v269)%32)|v269 < v330<<(uint(int32(2))%32)+v339 {
		} else {
			v343 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v277-l1<<(uint(int32(2))%32))))))
			v345 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+l1*int32(-3)))))
			v348 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v343-v345)))))
			if l3 < v348 {
			} else {
				v352 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v345-v334)))))
				if base.Ui32(l3) < base.Ui32(v352) {
				} else {
					v356 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v334-v326)))))
					if base.Ui32(l3) < base.Ui32(v356) {
					} else {
						v359 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+l1*int32(3)))))
						v361 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+v287))))
						v364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v359-v361)))))
						if base.Ui32(l3) < base.Ui32(v364) {
						} else {
							v368 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v361-v336)))))
							if base.Ui32(l3) < base.Ui32(v368) {
							} else {
								v372 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v296+(v336-v327)))))
								if base.Ui32(l3) < base.Ui32(v372) {
								} else {
									v376 = (v327 - v326) * int32(3)
									if l4 < v356 {
										v380 = int32(*(*int8)(unsafe.Add(mBase, uint32(v290+v337))))
										v381 = v376 + v380
										v384 = int32(3)
										v387 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+int32(4))>>(uint(v384)%32)))))
										v394 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+v384)>>(uint(v384)%32)))))
										v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v326+v394))))
										*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v396)
										v428 = v327 - v387
										v429 = v297
									} else {
										if base.Ui32(v372) <= base.Ui32(l4) {
											v399 = int32(3)
											v404 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v376+v399)>>(uint(v399)%32)))))
											v410 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v376+int32(4))>>(uint(v399)%32)))))
											v411 = int32(1)
											v414 = (v410 + v411) >> (uint(v411) % 32)
											v417 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v414+v334))))
											*(*uint8)(unsafe.Add(mBase, uint32(v333))) = uint8(v417)
											v421 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v404+(v292+v326)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v421)
											v425 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+(v327-v410)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v297))) = uint8(v425)
											v428 = v336 - v414
											v429 = v335
										} else {
											v380 = int32(*(*int8)(unsafe.Add(mBase, uint32(v290+v337))))
											v381 = v376 + v380
											v384 = int32(3)
											v387 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+int32(4))>>(uint(v384)%32)))))
											v394 = int32(*(*int8)(unsafe.Add(mBase, uint32(v294+(v381+v384)>>(uint(v384)%32)))))
											v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v326+v394))))
											*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v396)
											v428 = v327 - v387
											v429 = v297
										}
									}
									v435 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v292+v428))))
									*(*uint8)(unsafe.Add(mBase, uint32(v429))) = uint8(v435)
								}
							}
						}
					}
				}
			}
		}
		v446 = v301 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v446) {
			v297 = v297 + int32(1)
			v301 = v446
			continue
		} else {
			break
		}
		break
	}
	v503 = int32(1)
	v511 = int32(0)
	v521 = l1 << (uint(v503) % 32)
	v523 = m.G14
	v524 = *(*int32)(unsafe.Add(mBase, uint32(v523)))
	v525 = m.G12
	v526 = *(*int32)(unsafe.Add(mBase, uint32(v525)))
	v527 = m.G13
	v528 = *(*int32)(unsafe.Add(mBase, uint32(v527)))
	v529 = m.G15
	v530 = *(*int32)(unsafe.Add(mBase, uint32(v529)))
	v531 = v243 + v8
	v535 = int32(17)
	for {
		v559 = v531 + (v511 - l1)
		v560 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v559))))
		v561 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531))))
		v564 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v560-v561)))))
		v567 = v531 + (v511 - v521)
		v568 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v567))))
		v569 = v531 + l1
		v570 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v569))))
		v571 = v568 - v570
		v573 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+v571))))
		if l2<<(uint(v503)%32)|v503 < v564<<(uint(int32(2))%32)+v573 {
		} else {
			v577 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+(v511-l1<<(uint(int32(2))%32))))))
			v579 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+l1*int32(-3)))))
			v582 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v577-v579)))))
			if l3 < v582 {
			} else {
				v586 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v579-v568)))))
				if base.Ui32(l3) < base.Ui32(v586) {
				} else {
					v590 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v568-v560)))))
					if base.Ui32(l3) < base.Ui32(v590) {
					} else {
						v593 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+l1*int32(3)))))
						v595 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v531+v521))))
						v598 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v593-v595)))))
						if base.Ui32(l3) < base.Ui32(v598) {
						} else {
							v602 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v595-v570)))))
							if base.Ui32(l3) < base.Ui32(v602) {
							} else {
								v606 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v530+(v570-v561)))))
								if base.Ui32(l3) < base.Ui32(v606) {
								} else {
									v610 = (v561 - v560) * int32(3)
									if l4 < v590 {
										v614 = int32(*(*int8)(unsafe.Add(mBase, uint32(v524+v571))))
										v615 = v610 + v614
										v618 = int32(3)
										v621 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+int32(4))>>(uint(v618)%32)))))
										v628 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+v618)>>(uint(v618)%32)))))
										v630 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v560+v628))))
										*(*uint8)(unsafe.Add(mBase, uint32(v559))) = uint8(v630)
										v662 = v561 - v621
										v663 = v531
									} else {
										if base.Ui32(v606) <= base.Ui32(l4) {
											v633 = int32(3)
											v638 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v610+v633)>>(uint(v633)%32)))))
											v644 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v610+int32(4))>>(uint(v633)%32)))))
											v645 = int32(1)
											v648 = (v644 + v645) >> (uint(v645) % 32)
											v651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v648+v568))))
											*(*uint8)(unsafe.Add(mBase, uint32(v567))) = uint8(v651)
											v655 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v638+(v526+v560)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v559))) = uint8(v655)
											v659 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+(v561-v644)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v531))) = uint8(v659)
											v662 = v570 - v648
											v663 = v569
										} else {
											v614 = int32(*(*int8)(unsafe.Add(mBase, uint32(v524+v571))))
											v615 = v610 + v614
											v618 = int32(3)
											v621 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+int32(4))>>(uint(v618)%32)))))
											v628 = int32(*(*int8)(unsafe.Add(mBase, uint32(v528+(v615+v618)>>(uint(v618)%32)))))
											v630 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v560+v628))))
											*(*uint8)(unsafe.Add(mBase, uint32(v559))) = uint8(v630)
											v662 = v561 - v621
											v663 = v531
										}
									}
									v669 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v526+v662))))
									*(*uint8)(unsafe.Add(mBase, uint32(v663))) = uint8(v669)
								}
							}
						}
					}
				}
			}
		}
		v680 = v535 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v680) {
			v531 = v531 + int32(1)
			v535 = v680
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_VFilter8_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v36 int32
	_ = v36
	var v44 int32
	_ = v44
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v68 int32
	_ = v68
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
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
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v147 int32
	_ = v147
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v158 int32
	_ = v158
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v188 int32
	_ = v188
	var v191 int32
	_ = v191
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v217 int32
	_ = v217
	var v219 int32
	_ = v219
	var v233 int32
	_ = v233
	var v297 int32
	_ = v297
	var v305 int32
	_ = v305
	var v315 int32
	_ = v315
	var v317 int32
	_ = v317
	var v318 int32
	_ = v318
	var v319 int32
	_ = v319
	var v320 int32
	_ = v320
	var v321 int32
	_ = v321
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v329 int32
	_ = v329
	var v357 int32
	_ = v357
	var v358 int32
	_ = v358
	var v359 int32
	_ = v359
	var v362 int32
	_ = v362
	var v365 int32
	_ = v365
	var v366 int32
	_ = v366
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v375 int32
	_ = v375
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v388 int32
	_ = v388
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v400 int32
	_ = v400
	var v404 int32
	_ = v404
	var v408 int32
	_ = v408
	var v410 int32
	_ = v410
	var v413 int32
	_ = v413
	var v416 int32
	_ = v416
	var v419 int32
	_ = v419
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v433 int32
	_ = v433
	var v436 int32
	_ = v436
	var v438 int32
	_ = v438
	var v439 int32
	_ = v439
	var v442 int32
	_ = v442
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v459 int32
	_ = v459
	var v462 int32
	_ = v462
	var v466 int32
	_ = v466
	var v470 int32
	_ = v470
	var v473 int32
	_ = v473
	var v478 int32
	_ = v478
	var v480 int32
	_ = v480
	var v494 int32
	_ = v494
	v36 = int32(1)
	v44 = int32(0)
	v54 = l2 << (uint(v36) % 32)
	v56 = m.G13
	v57 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
	v58 = m.G12
	v59 = *(*int32)(unsafe.Add(mBase, uint32(v58)))
	v60 = m.G14
	v61 = *(*int32)(unsafe.Add(mBase, uint32(v60)))
	v62 = m.G15
	v63 = *(*int32)(unsafe.Add(mBase, uint32(v62)))
	v64 = l0
	v68 = int32(9)
	for {
		v96 = v64 + (v44 - l2)
		v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v96))))
		v98 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64))))
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v97-v98)))))
		v104 = v64 + (v44 - v54)
		v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v104))))
		v106 = v64 + l2
		v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v106))))
		v108 = v105 - v107
		v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+v108))))
		if l3<<(uint(v36)%32)|v36 < v101<<(uint(int32(2))%32)+v110 {
		} else {
			v114 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+(v44-l2<<(uint(int32(2))%32))))))
			v115 = v64 + l2*int32(-3)
			v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v115))))
			v119 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v114-v116)))))
			if l4 < v119 {
			} else {
				v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v116-v105)))))
				if base.Ui32(l4) < base.Ui32(v123) {
				} else {
					v127 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v105-v97)))))
					if base.Ui32(l4) < base.Ui32(v127) {
					} else {
						v130 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+l2*int32(3)))))
						v131 = v64 + v54
						v132 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v131))))
						v135 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v130-v132)))))
						if base.Ui32(l4) < base.Ui32(v135) {
						} else {
							v139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v132-v107)))))
							if base.Ui32(l4) < base.Ui32(v139) {
							} else {
								v143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v107-v98)))))
								if base.Ui32(l4) < base.Ui32(v143) {
								} else {
									v147 = (v98 - v97) * int32(3)
									v149 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+v108))))
									if l5 < v127 {
										v152 = v147 + v149
										v155 = int32(3)
										v158 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+int32(4))>>(uint(v155)%32)))))
										v165 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+v155)>>(uint(v155)%32)))))
										v167 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v97+v165))))
										*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v167)
										v212 = v98 - v158
										v217 = v64
									} else {
										if base.Ui32(v143) <= base.Ui32(l5) {
											v172 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+v147+v149))))
											v175 = int32(63)
											v177 = int32(7)
											v178 = (v172*int32(9) + v175) >> (uint(v177) % 32)
											v181 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v178+v116))))
											*(*uint8)(unsafe.Add(mBase, uint32(v115))) = uint8(v181)
											v188 = (v172*int32(18) + v175) >> (uint(v177) % 32)
											v191 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v188+v105))))
											*(*uint8)(unsafe.Add(mBase, uint32(v104))) = uint8(v191)
											v198 = (v172*int32(27) + v175) >> (uint(v177) % 32)
											v201 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v198+v97))))
											*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v201)
											v205 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+(v98-v198)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v64))) = uint8(v205)
											v209 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+(v107-v188)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v106))) = uint8(v209)
											v212 = v132 - v178
											v217 = v131
										} else {
											v152 = v147 + v149
											v155 = int32(3)
											v158 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+int32(4))>>(uint(v155)%32)))))
											v165 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+(v152+v155)>>(uint(v155)%32)))))
											v167 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v97+v165))))
											*(*uint8)(unsafe.Add(mBase, uint32(v96))) = uint8(v167)
											v212 = v98 - v158
											v217 = v64
										}
									}
									v219 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v212))))
									*(*uint8)(unsafe.Add(mBase, uint32(v217))) = uint8(v219)
								}
							}
						}
					}
				}
			}
		}
		v233 = v68 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v233) {
			v64 = v64 + int32(1)
			v68 = v233
			continue
		} else {
			break
		}
		break
	}
	v297 = int32(1)
	v305 = int32(0)
	v315 = l2 << (uint(v297) % 32)
	v317 = m.G13
	v318 = *(*int32)(unsafe.Add(mBase, uint32(v317)))
	v319 = m.G12
	v320 = *(*int32)(unsafe.Add(mBase, uint32(v319)))
	v321 = m.G14
	v322 = *(*int32)(unsafe.Add(mBase, uint32(v321)))
	v323 = m.G15
	v324 = *(*int32)(unsafe.Add(mBase, uint32(v323)))
	v325 = l1
	v329 = int32(9)
	for {
		v357 = v325 + (v305 - l2)
		v358 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v357))))
		v359 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325))))
		v362 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v358-v359)))))
		v365 = v325 + (v305 - v315)
		v366 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v365))))
		v367 = v325 + l2
		v368 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v367))))
		v369 = v366 - v368
		v371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+v369))))
		if l3<<(uint(v297)%32)|v297 < v362<<(uint(int32(2))%32)+v371 {
		} else {
			v375 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325+(v305-l2<<(uint(int32(2))%32))))))
			v376 = v325 + l2*int32(-3)
			v377 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v376))))
			v380 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v375-v377)))))
			if l4 < v380 {
			} else {
				v384 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v377-v366)))))
				if base.Ui32(l4) < base.Ui32(v384) {
				} else {
					v388 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v366-v358)))))
					if base.Ui32(l4) < base.Ui32(v388) {
					} else {
						v391 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v325+l2*int32(3)))))
						v392 = v325 + v315
						v393 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v392))))
						v396 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v391-v393)))))
						if base.Ui32(l4) < base.Ui32(v396) {
						} else {
							v400 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v393-v368)))))
							if base.Ui32(l4) < base.Ui32(v400) {
							} else {
								v404 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v324+(v368-v359)))))
								if base.Ui32(l4) < base.Ui32(v404) {
								} else {
									v408 = (v359 - v358) * int32(3)
									v410 = int32(*(*int8)(unsafe.Add(mBase, uint32(v322+v369))))
									if l5 < v388 {
										v413 = v408 + v410
										v416 = int32(3)
										v419 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+int32(4))>>(uint(v416)%32)))))
										v426 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+v416)>>(uint(v416)%32)))))
										v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v358+v426))))
										*(*uint8)(unsafe.Add(mBase, uint32(v357))) = uint8(v428)
										v473 = v359 - v419
										v478 = v325
									} else {
										if base.Ui32(v404) <= base.Ui32(l5) {
											v433 = int32(*(*int8)(unsafe.Add(mBase, uint32(v322+v408+v410))))
											v436 = int32(63)
											v438 = int32(7)
											v439 = (v433*int32(9) + v436) >> (uint(v438) % 32)
											v442 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v439+v377))))
											*(*uint8)(unsafe.Add(mBase, uint32(v376))) = uint8(v442)
											v449 = (v433*int32(18) + v436) >> (uint(v438) % 32)
											v452 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v449+v366))))
											*(*uint8)(unsafe.Add(mBase, uint32(v365))) = uint8(v452)
											v459 = (v433*int32(27) + v436) >> (uint(v438) % 32)
											v462 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v459+v358))))
											*(*uint8)(unsafe.Add(mBase, uint32(v357))) = uint8(v462)
											v466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+(v359-v459)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v466)
											v470 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+(v368-v449)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v367))) = uint8(v470)
											v473 = v393 - v439
											v478 = v392
										} else {
											v413 = v408 + v410
											v416 = int32(3)
											v419 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+int32(4))>>(uint(v416)%32)))))
											v426 = int32(*(*int8)(unsafe.Add(mBase, uint32(v318+(v413+v416)>>(uint(v416)%32)))))
											v428 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v358+v426))))
											*(*uint8)(unsafe.Add(mBase, uint32(v357))) = uint8(v428)
											v473 = v359 - v419
											v478 = v325
										}
									}
									v480 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v320+v473))))
									*(*uint8)(unsafe.Add(mBase, uint32(v478))) = uint8(v480)
								}
							}
						}
					}
				}
			}
		}
		v494 = v329 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v494) {
			v325 = v325 + int32(1)
			v329 = v494
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_VFilter8i_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v36 int32
	_ = v36
	var v44 int32
	_ = v44
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v68 int32
	_ = v68
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v106 int32
	_ = v106
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v115 int32
	_ = v115
	var v119 int32
	_ = v119
	var v123 int32
	_ = v123
	var v126 int32
	_ = v126
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v135 int32
	_ = v135
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v154 int32
	_ = v154
	var v161 int32
	_ = v161
	var v163 int32
	_ = v163
	var v166 int32
	_ = v166
	var v171 int32
	_ = v171
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v192 int32
	_ = v192
	var v195 int32
	_ = v195
	var v196 int32
	_ = v196
	var v202 int32
	_ = v202
	var v213 int32
	_ = v213
	var v270 int32
	_ = v270
	var v278 int32
	_ = v278
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v297 int32
	_ = v297
	var v298 int32
	_ = v298
	var v302 int32
	_ = v302
	var v326 int32
	_ = v326
	var v327 int32
	_ = v327
	var v328 int32
	_ = v328
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v344 int32
	_ = v344
	var v346 int32
	_ = v346
	var v349 int32
	_ = v349
	var v353 int32
	_ = v353
	var v357 int32
	_ = v357
	var v360 int32
	_ = v360
	var v362 int32
	_ = v362
	var v365 int32
	_ = v365
	var v369 int32
	_ = v369
	var v373 int32
	_ = v373
	var v377 int32
	_ = v377
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v385 int32
	_ = v385
	var v388 int32
	_ = v388
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v400 int32
	_ = v400
	var v405 int32
	_ = v405
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v415 int32
	_ = v415
	var v418 int32
	_ = v418
	var v422 int32
	_ = v422
	var v426 int32
	_ = v426
	var v429 int32
	_ = v429
	var v430 int32
	_ = v430
	var v436 int32
	_ = v436
	var v447 int32
	_ = v447
	v9 = l2 << (uint(int32(2)) % 32)
	v36 = int32(1)
	v44 = int32(0)
	v54 = l2 << (uint(v36) % 32)
	v56 = m.G14
	v57 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
	v58 = m.G12
	v59 = *(*int32)(unsafe.Add(mBase, uint32(v58)))
	v60 = m.G13
	v61 = *(*int32)(unsafe.Add(mBase, uint32(v60)))
	v62 = m.G15
	v63 = *(*int32)(unsafe.Add(mBase, uint32(v62)))
	v64 = l0 + v9
	v68 = int32(9)
	for {
		v92 = v64 + (v44 - l2)
		v93 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v92))))
		v94 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64))))
		v97 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v93-v94)))))
		v100 = v64 + (v44 - v54)
		v101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v100))))
		v102 = v64 + l2
		v103 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v102))))
		v104 = v101 - v103
		v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+v104))))
		if l3<<(uint(v36)%32)|v36 < v97<<(uint(int32(2))%32)+v106 {
		} else {
			v110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+(v44-l2<<(uint(int32(2))%32))))))
			v112 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+l2*int32(-3)))))
			v115 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v110-v112)))))
			if l4 < v115 {
			} else {
				v119 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v112-v101)))))
				if base.Ui32(l4) < base.Ui32(v119) {
				} else {
					v123 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v101-v93)))))
					if base.Ui32(l4) < base.Ui32(v123) {
					} else {
						v126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+l2*int32(3)))))
						v128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v64+v54))))
						v131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v126-v128)))))
						if base.Ui32(l4) < base.Ui32(v131) {
						} else {
							v135 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v128-v103)))))
							if base.Ui32(l4) < base.Ui32(v135) {
							} else {
								v139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v63+(v103-v94)))))
								if base.Ui32(l4) < base.Ui32(v139) {
								} else {
									v143 = (v94 - v93) * int32(3)
									if l5 < v123 {
										v147 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+v104))))
										v148 = v143 + v147
										v151 = int32(3)
										v154 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+(v148+int32(4))>>(uint(v151)%32)))))
										v161 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+(v148+v151)>>(uint(v151)%32)))))
										v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v93+v161))))
										*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v163)
										v195 = v94 - v154
										v196 = v64
									} else {
										if base.Ui32(v139) <= base.Ui32(l5) {
											v166 = int32(3)
											v171 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+(v143+v166)>>(uint(v166)%32)))))
											v177 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+(v143+int32(4))>>(uint(v166)%32)))))
											v178 = int32(1)
											v181 = (v177 + v178) >> (uint(v178) % 32)
											v184 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v181+v101))))
											*(*uint8)(unsafe.Add(mBase, uint32(v100))) = uint8(v184)
											v188 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v171+(v59+v93)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v188)
											v192 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+(v94-v177)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v64))) = uint8(v192)
											v195 = v103 - v181
											v196 = v102
										} else {
											v147 = int32(*(*int8)(unsafe.Add(mBase, uint32(v57+v104))))
											v148 = v143 + v147
											v151 = int32(3)
											v154 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+(v148+int32(4))>>(uint(v151)%32)))))
											v161 = int32(*(*int8)(unsafe.Add(mBase, uint32(v61+(v148+v151)>>(uint(v151)%32)))))
											v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v93+v161))))
											*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v163)
											v195 = v94 - v154
											v196 = v64
										}
									}
									v202 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v59+v195))))
									*(*uint8)(unsafe.Add(mBase, uint32(v196))) = uint8(v202)
								}
							}
						}
					}
				}
			}
		}
		v213 = v68 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v213) {
			v64 = v64 + int32(1)
			v68 = v213
			continue
		} else {
			break
		}
		break
	}
	v270 = int32(1)
	v278 = int32(0)
	v288 = l2 << (uint(v270) % 32)
	v290 = m.G14
	v291 = *(*int32)(unsafe.Add(mBase, uint32(v290)))
	v292 = m.G12
	v293 = *(*int32)(unsafe.Add(mBase, uint32(v292)))
	v294 = m.G13
	v295 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v296 = m.G15
	v297 = *(*int32)(unsafe.Add(mBase, uint32(v296)))
	v298 = l1 + v9
	v302 = int32(9)
	for {
		v326 = v298 + (v278 - l2)
		v327 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v326))))
		v328 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v298))))
		v331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v327-v328)))))
		v334 = v298 + (v278 - v288)
		v335 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v334))))
		v336 = v298 + l2
		v337 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v336))))
		v338 = v335 - v337
		v340 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+v338))))
		if l3<<(uint(v270)%32)|v270 < v331<<(uint(int32(2))%32)+v340 {
		} else {
			v344 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v298+(v278-l2<<(uint(int32(2))%32))))))
			v346 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v298+l2*int32(-3)))))
			v349 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v344-v346)))))
			if l4 < v349 {
			} else {
				v353 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v346-v335)))))
				if base.Ui32(l4) < base.Ui32(v353) {
				} else {
					v357 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v335-v327)))))
					if base.Ui32(l4) < base.Ui32(v357) {
					} else {
						v360 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v298+l2*int32(3)))))
						v362 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v298+v288))))
						v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v360-v362)))))
						if base.Ui32(l4) < base.Ui32(v365) {
						} else {
							v369 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v362-v337)))))
							if base.Ui32(l4) < base.Ui32(v369) {
							} else {
								v373 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v297+(v337-v328)))))
								if base.Ui32(l4) < base.Ui32(v373) {
								} else {
									v377 = (v328 - v327) * int32(3)
									if l5 < v357 {
										v381 = int32(*(*int8)(unsafe.Add(mBase, uint32(v291+v338))))
										v382 = v377 + v381
										v385 = int32(3)
										v388 = int32(*(*int8)(unsafe.Add(mBase, uint32(v295+(v382+int32(4))>>(uint(v385)%32)))))
										v395 = int32(*(*int8)(unsafe.Add(mBase, uint32(v295+(v382+v385)>>(uint(v385)%32)))))
										v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v293+v327+v395))))
										*(*uint8)(unsafe.Add(mBase, uint32(v326))) = uint8(v397)
										v429 = v328 - v388
										v430 = v298
									} else {
										if base.Ui32(v373) <= base.Ui32(l5) {
											v400 = int32(3)
											v405 = int32(*(*int8)(unsafe.Add(mBase, uint32(v295+(v377+v400)>>(uint(v400)%32)))))
											v411 = int32(*(*int8)(unsafe.Add(mBase, uint32(v295+(v377+int32(4))>>(uint(v400)%32)))))
											v412 = int32(1)
											v415 = (v411 + v412) >> (uint(v412) % 32)
											v418 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v293+v415+v335))))
											*(*uint8)(unsafe.Add(mBase, uint32(v334))) = uint8(v418)
											v422 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v405+(v293+v327)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v326))) = uint8(v422)
											v426 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v293+(v328-v411)))))
											*(*uint8)(unsafe.Add(mBase, uint32(v298))) = uint8(v426)
											v429 = v337 - v415
											v430 = v336
										} else {
											v381 = int32(*(*int8)(unsafe.Add(mBase, uint32(v291+v338))))
											v382 = v377 + v381
											v385 = int32(3)
											v388 = int32(*(*int8)(unsafe.Add(mBase, uint32(v295+(v382+int32(4))>>(uint(v385)%32)))))
											v395 = int32(*(*int8)(unsafe.Add(mBase, uint32(v295+(v382+v385)>>(uint(v385)%32)))))
											v397 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v293+v327+v395))))
											*(*uint8)(unsafe.Add(mBase, uint32(v326))) = uint8(v397)
											v429 = v328 - v388
											v430 = v298
										}
									}
									v436 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v293+v429))))
									*(*uint8)(unsafe.Add(mBase, uint32(v430))) = uint8(v436)
								}
							}
						}
					}
				}
			}
		}
		v447 = v302 + int32(-1)
		if base.Ui32(int32(1)) < base.Ui32(v447) {
			v298 = v298 + int32(1)
			v302 = v447
			continue
		} else {
			break
		}
		break
	}
	return
}
func F_VL4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
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
	var v42 int32
	_ = v42
	var v47 int32
	_ = v47
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v67 int32
	_ = v67
	var v74 int32
	_ = v74
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v88 int32
	_ = v88
	var v91 int32
	_ = v91
	var v101 int32
	_ = v101
	var v110 int32
	_ = v110
	v11 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-31)))))
	v14 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-30)))))
	v16 = int32(1)
	v19 = int32(base.Ui32(v11+v14+v16) >> (uint(v16) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v19)
	v23 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-32)))))
	v28 = int32(base.Ui32(v11+v23+v16) >> (uint(v16) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v28)
	v32 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-29)))))
	v37 = int32(base.Ui32(v14+v32+v16) >> (uint(v16) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v37)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v19)
	v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-28)))))
	v47 = int32(base.Ui32(v32+v42+v16) >> (uint(v16) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v47)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v37)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v47)
	v51 = int32(2)
	v52 = v32 + v51
	v58 = int32(base.Ui32(v11+v52+v14<<(uint(v16)%32)) >> (uint(v51) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v58)
	v61 = v14 + v51
	v67 = int32(base.Ui32(v23+v61+v11<<(uint(v16)%32)) >> (uint(v51) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v67)
	v74 = int32(base.Ui32(v42+(v61+v32<<(uint(v16)%32))) >> (uint(v51) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v74)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v58)
	v79 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-25)))))
	v82 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-26)))))
	v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-27)))))
	v91 = int32(base.Ui32(v52+v42<<(uint(v16)%32)+v88) >> (uint(v51) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v91)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v74)
	v101 = int32(base.Ui32(v79+(v88+v82<<(uint(v16)%32))+v51) >> (uint(v51) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v101)
	v110 = int32(base.Ui32(v82+(v42+v88<<(uint(v16)%32))+v51) >> (uint(v51) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v110)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v91)
	return
}
func F_VP8LAddGreenToBlueAndRed_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v72 int32
	_ = v72
	var v78 int32
	_ = v78
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v94 int32
	_ = v94
	if l1 < int32(1) {
	} else {
		v12 = int32(1)
		if l1 == v12 {
			v78 = int32(0)
		} else {
			v21 = l2
			v24 = int32(0)
			v26 = l0
			for {
				v29 = *(*int32)(unsafe.Add(mBase, uint32(v26)))
				v30 = int32(8)
				v31 = int32(base.Ui32(v29) >> (uint(v30) % 32))
				v32 = int32(255)
				v34 = int32(16711935)
				v37 = int32(16)
				v42 = int32(-16711936)
				*(*int32)(unsafe.Add(mBase, uint32(v21))) = (v31&v32+v29&v34+v31<<(uint(v37)%32))&v34 | v29&v42
				v46 = int32(4)
				v50 = *(*int32)(unsafe.Add(mBase, uint32(v26+v46)))
				v52 = int32(base.Ui32(v50) >> (uint(v30) % 32))
				*(*int32)(unsafe.Add(mBase, uint32(v21+v46))) = (v52&v32+v50&v34+v52<<(uint(v37)%32))&v34 | v50&v42
				v72 = v24 + int32(2)
				if l1&int32(2147483646) != v72 {
					v21 = v21 + v30
					v24 = v72
					v26 = v26 + v30
					continue
				} else {
					break
				}
				break
			}
			v78 = v72
		}
		if l1&v12 == int32(0) {
		} else {
			v86 = v78 << (uint(int32(2)) % 32)
			v89 = *(*int32)(unsafe.Add(mBase, uint32(l0+v86)))
			v91 = int32(base.Ui32(v89) >> (uint(int32(8)) % 32))
			v94 = int32(16711935)
			*(*int32)(unsafe.Add(mBase, uint32(l2+v86))) = (v91&int32(255)+v89&v94+v91<<(uint(int32(16))%32))&v94 | v89&int32(-16711936)
		}
	}
	return
}
func F_VP8LBundleColorMap_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v56 int32
	_ = v56
	var v66 int32
	_ = v66
	var v75 int32
	_ = v75
	var v82 int32
	_ = v82
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v116 int32
	_ = v116
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v125 int32
	_ = v125
	var v136 int32
	_ = v136
	var v141 int32
	_ = v141
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v169 int32
	_ = v169
	var v170 int32
	_ = v170
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v180 int32
	_ = v180
	var v183 int32
	_ = v183
	var v188 int32
	_ = v188
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v207 int32
	_ = v207
	if int32(0) < l2 {
		if l1 < int32(1) {
		} else {
			v120 = int32(3) - l2
			v121 = int32(-1)
			v124 = v121<<(uint(l2)%32) ^ v121
			v125 = int32(1)
			if l1 != v125 {
				v136 = int32(0)
				v141 = int32(-16777216)
				for {
					v150 = l0 + v136
					v151 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v150))))
					v152 = v136 & v124
					if v152 != 0 {
						v158 = v141
					} else {
						v158 = int32(-16777216)
					}
					v159 = v151<<(uint(v152<<(uint(v120)%32)+int32(8))%32) | v158
					*(*int32)(unsafe.Add(mBase, uint32(l3+int32(base.Ui32(v136)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v159
					v161 = int32(1)
					v162 = v136 + v161
					v169 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v150+v161))))
					v170 = v162 & v124
					if v170 != 0 {
						v176 = v159
					} else {
						v176 = int32(-16777216)
					}
					v177 = v169<<(uint(v170<<(uint(v120)%32)+int32(8))%32) | v176
					*(*int32)(unsafe.Add(mBase, uint32(l3+int32(base.Ui32(v162)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v177
					v180 = v136 + int32(2)
					if v180 != l1&int32(2147483646) {
						v136 = v180
						v141 = v177
						continue
					} else {
						break
					}
					break
				}
				v183 = v180
				v188 = v177
			} else {
				v183 = int32(0)
				v188 = int32(-16777216)
			}
			if l1&v125 == int32(0) {
			} else {
				v200 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v183))))
				v201 = v183 & v124
				if v201 != 0 {
					v207 = v188
				} else {
					v207 = int32(-16777216)
				}
				*(*int32)(unsafe.Add(mBase, uint32(l3+int32(base.Ui32(v183)>>(uint(l2)%32))<<(uint(int32(2))%32)))) = v200<<(uint(v201<<(uint(v120)%32)+int32(8))%32) | v207
			}
		}
	} else {
		if l1 < int32(1) {
		} else {
			v17 = l1 & int32(3)
			if base.Ui32(l1) < base.Ui32(int32(4)) {
				v82 = int32(0)
			} else {
				v25 = l3
				v29 = int32(0)
				for {
					v35 = l0 + v29
					v36 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v35))))
					v37 = int32(8)
					v39 = int32(-16777216)
					*(*int32)(unsafe.Add(mBase, uint32(v25))) = v36<<(uint(v37)%32) | v39
					v42 = int32(4)
					v46 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(1)))))
					*(*int32)(unsafe.Add(mBase, uint32(v25+v42))) = v46<<(uint(v37)%32) | v39
					v56 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(2)))))
					*(*int32)(unsafe.Add(mBase, uint32(v25+v37))) = v56<<(uint(v37)%32) | v39
					v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v35+int32(3)))))
					*(*int32)(unsafe.Add(mBase, uint32(v25+int32(12)))) = v66<<(uint(v37)%32) | v39
					v75 = v29 + v42
					if l1&int32(2147483644) != v75 {
						v25 = v25 + int32(16)
						v29 = v75
						continue
					} else {
						break
					}
					break
				}
				v82 = v75
			}
			if v17 == int32(0) {
			} else {
				v95 = l0 + v82
				v96 = l3 + v82<<(uint(int32(2))%32)
				v98 = v17
				for {
					v105 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v95))))
					*(*int32)(unsafe.Add(mBase, uint32(v96))) = v105<<(uint(int32(8))%32) | int32(-16777216)
					v116 = v98 + int32(-1)
					if v116 != 0 {
						v95 = v95 + int32(1)
						v96 = v96 + int32(4)
						v98 = v116
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
func F_VP8LCollectColorBlueTransforms_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v34 int32
	_ = v34
	var v38 int32
	_ = v38
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v73 int32
	_ = v73
	if l3 < int32(1) {
	} else {
		if l2 < int32(1) {
		} else {
			v20 = l0
			v23 = l3
			for {
				v34 = l2
				v38 = v20
				for {
					v44 = *(*int32)(unsafe.Add(mBase, uint32(v38)))
					v47 = int32(24)
					v50 = int32(5)
					v65 = l6 + (v44-(int32(base.Ui32(v44<<(uint(int32(16))%32)>>(uint(v47)%32)*base.I32_extend8_s(l4))>>(uint(v50)%32))+int32(base.Ui32(v44<<(uint(int32(8))%32)>>(uint(v47)%32)*base.I32_extend8_s(l5))>>(uint(v50)%32))))&int32(255)<<(uint(int32(2))%32)
					v66 = *(*int32)(unsafe.Add(mBase, uint32(v65)))
					*(*int32)(unsafe.Add(mBase, uint32(v65))) = v66 + int32(1)
					v73 = v34 + int32(-1)
					if v73 != 0 {
						v34 = v73
						v38 = v38 + int32(4)
						continue
					} else {
						break
					}
					break
				}
				if int32(1) < v23 {
					v20 = v20 + l1<<(uint(int32(2))%32)
					v23 = v23 + int32(-1)
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
func F_VP8LCollectColorRedTransforms_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v71 int32
	_ = v71
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	var v127 int32
	_ = v127
	var v128 int32
	_ = v128
	if l3 < int32(1) {
	} else {
		if l2 < int32(1) {
		} else {
			v16 = base.I32_extend8_s(l4)
			v23 = l0
			v26 = l3
			for {
				if l2 == int32(1) {
					v97 = int32(0)
				} else {
					v39 = int32(0)
					v42 = v23
					for {
						v49 = *(*int32)(unsafe.Add(mBase, uint32(v42)))
						v50 = int32(16)
						v54 = int32(24)
						v57 = int32(5)
						v60 = int32(255)
						v62 = int32(2)
						v64 = l5 + (int32(base.Ui32(v49)>>(uint(v50)%32))-int32(base.Ui32(v49<<(uint(v50)%32)>>(uint(v54)%32)*v16)>>(uint(v57)%32)))&v60<<(uint(v62)%32)
						v65 = *(*int32)(unsafe.Add(mBase, uint32(v64)))
						v66 = int32(1)
						*(*int32)(unsafe.Add(mBase, uint32(v64))) = v65 + v66
						v71 = *(*int32)(unsafe.Add(mBase, uint32(v42+int32(4))))
						v86 = l5 + (int32(base.Ui32(v71)>>(uint(v50)%32))-int32(base.Ui32(v71<<(uint(v50)%32)>>(uint(v54)%32)*v16)>>(uint(v57)%32)))&v60<<(uint(v62)%32)
						v87 = *(*int32)(unsafe.Add(mBase, uint32(v86)))
						*(*int32)(unsafe.Add(mBase, uint32(v86))) = v87 + v66
						v94 = v39 + v62
						if l2&int32(2147483646) != v94 {
							v39 = v94
							v42 = v42 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v97 = v94
				}
				if l2&int32(1) == int32(0) {
				} else {
					v109 = int32(2)
					v112 = *(*int32)(unsafe.Add(mBase, uint32(v23+v97<<(uint(v109)%32))))
					v113 = int32(16)
					v127 = l5 + (int32(base.Ui32(v112)>>(uint(v113)%32))-int32(base.Ui32(v112<<(uint(v113)%32)>>(uint(int32(24))%32)*v16)>>(uint(int32(5))%32)))&int32(255)<<(uint(v109)%32)
					v128 = *(*int32)(unsafe.Add(mBase, uint32(v127)))
					*(*int32)(unsafe.Add(mBase, uint32(v127))) = v128 + int32(1)
				}
				if int32(1) < v26 {
					v23 = v23 + l1<<(uint(int32(2))%32)
					v26 = v26 + int32(-1)
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
func F_VP8LConvertBGRAToBGR_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	if l1 < int32(1) {
	} else {
		v10 = l0
		v12 = l2
		for {
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
			*(*uint8)(unsafe.Add(mBase, uint32(v12))) = uint8(v14)
			v17 = int32(base.Ui32(v14) >> (uint(int32(16)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+2)) = uint8(v17)
			v20 = int32(base.Ui32(v14) >> (uint(int32(8)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+1)) = uint8(v20)
			v25 = v10 + int32(4)
			if base.Ui32(v25) < base.Ui32(l0+l1<<(uint(int32(2))%32)) {
				v10 = v25
				v12 = v12 + int32(3)
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LConvertBGRAToRGB565_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v23 int32
	_ = v23
	var v33 int32
	_ = v33
	var v38 int32
	_ = v38
	if l1 < int32(1) {
	} else {
		v10 = l0
		v12 = l2
		for {
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
			v23 = int32(base.Ui32(v14)>>(uint(int32(5))%32))&int32(224) | int32(base.Ui32(v14)>>(uint(int32(3))%32))&int32(31)
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+1)) = uint8(v23)
			v33 = int32(base.Ui32(v14)>>(uint(int32(16))%32))&int32(248) | int32(base.Ui32(v14)>>(uint(int32(13))%32))&int32(7)
			*(*uint8)(unsafe.Add(mBase, uint32(v12))) = uint8(v33)
			v38 = v10 + int32(4)
			if base.Ui32(v38) < base.Ui32(l0+l1<<(uint(int32(2))%32)) {
				v10 = v38
				v12 = v12 + int32(2)
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LConvertBGRAToRGBA4444_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v19 int32
	_ = v19
	var v29 int32
	_ = v29
	var v34 int32
	_ = v34
	if l1 < int32(1) {
	} else {
		v10 = l0
		v12 = l2
		for {
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
			v15 = int32(240)
			v19 = v14&v15 | int32(base.Ui32(v14)>>(uint(int32(28))%32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+1)) = uint8(v19)
			v29 = int32(base.Ui32(v14)>>(uint(int32(16))%32))&v15 | int32(base.Ui32(v14)>>(uint(int32(12))%32))&int32(15)
			*(*uint8)(unsafe.Add(mBase, uint32(v12))) = uint8(v29)
			v34 = v10 + int32(4)
			if base.Ui32(v34) < base.Ui32(l0+l1<<(uint(int32(2))%32)) {
				v10 = v34
				v12 = v12 + int32(2)
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LConvertBGRAToRGBA_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v28 int32
	_ = v28
	if l1 < int32(1) {
	} else {
		v10 = l0
		v12 = l2
		for {
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+2)) = uint8(v14)
			v17 = int32(base.Ui32(v14) >> (uint(int32(24)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+3)) = uint8(v17)
			v20 = int32(base.Ui32(v14) >> (uint(int32(8)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+1)) = uint8(v20)
			v23 = int32(base.Ui32(v14) >> (uint(int32(16)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12))) = uint8(v23)
			v25 = int32(4)
			v28 = v10 + v25
			if base.Ui32(v28) < base.Ui32(l0+l1<<(uint(int32(2))%32)) {
				v10 = v28
				v12 = v12 + v25
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LConvertBGRAToRGB_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	if l1 < int32(1) {
	} else {
		v10 = l0
		v12 = l2
		for {
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+2)) = uint8(v14)
			v17 = int32(base.Ui32(v14) >> (uint(int32(8)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12)+1)) = uint8(v17)
			v20 = int32(base.Ui32(v14) >> (uint(int32(16)) % 32))
			*(*uint8)(unsafe.Add(mBase, uint32(v12))) = uint8(v20)
			v25 = v10 + int32(4)
			if base.Ui32(v25) < base.Ui32(l0+l1<<(uint(int32(2))%32)) {
				v10 = v25
				v12 = v12 + int32(3)
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LPredictor0_C(m *base.Module, l0 int32, l1 int32) int32 {
	return int32(-16777216)
}
func F_VP8LPredictor10_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v24 int32
	_ = v24
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v8 = int32(1)
	v10 = int32(2139062143)
	v13 = int32(base.Ui32(v5^v6)>>(uint(v8)%32))&v10 + v5&v6
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v24 = int32(base.Ui32(v16^v17)>>(uint(v8)%32))&v10 + v16&v17
	return int32(base.Ui32(v13^v24)>>(uint(v8)%32))&v10 + v13&v24
}
func F_VP8LPredictor11_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v28 int32
	_ = v28
	var v32 int32
	_ = v32
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v55 int32
	_ = v55
	var v57 int32
	_ = v57
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v71 int32
	_ = v71
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v81 int32
	_ = v81
	var v90 int32
	_ = v90
	var v92 int32
	_ = v92
	var v98 int32
	_ = v98
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v10 = int32(255)
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	v16 = v14 & v10
	v17 = v9&v10 - v16
	v18 = int32(31)
	v19 = v17 >> (uint(v18) % 32)
	v22 = int32(24)
	v25 = int32(base.Ui32(v14) >> (uint(v22) % 32))
	v26 = int32(base.Ui32(v9)>>(uint(v22)%32)) - v25
	v28 = v26 >> (uint(v18) % 32)
	v32 = int32(8)
	v39 = int32(base.Ui32(v14)>>(uint(v32)%32)) & v10
	v40 = int32(base.Ui32(v9)>>(uint(v32)%32))&v10 - v39
	v42 = v40 >> (uint(v18) % 32)
	v48 = int32(base.Ui32(v8)>>(uint(v22)%32)) - v25
	v50 = v48 >> (uint(v18) % 32)
	v55 = v8&v10 - v16
	v57 = v55 >> (uint(v18) % 32)
	v65 = int32(base.Ui32(v8)>>(uint(v32)%32))&v10 - v39
	v67 = v65 >> (uint(v18) % 32)
	v71 = int32(16)
	v78 = int32(base.Ui32(v14)>>(uint(v71)%32)) & v10
	v79 = int32(base.Ui32(v8)>>(uint(v71)%32))&v10 - v78
	v81 = v79 >> (uint(v18) % 32)
	v90 = int32(base.Ui32(v9)>>(uint(v71)%32))&v10 - v78
	v92 = v90 >> (uint(v18) % 32)
	if v17^v19-v19+(v26^v28-v28)+(v40^v42-v42)-(v48^v50-v50+(v55^v57-v57)+(v65^v67-v67)+(v79^v81-v81))+(v90^v92-v92) < int32(1) {
		v98 = v8
	} else {
		v98 = v9
	}
	return v98
}
func F_VP8LPredictor12_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v34 int32
	_ = v34
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v56 int32
	_ = v56
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v80 int32
	_ = v80
	var v87 int32
	_ = v87
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v6 = int32(24)
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	v17 = int32(base.Ui32(v5)>>(uint(v6)%32)) + int32(base.Ui32(v8)>>(uint(v6)%32)) - int32(base.Ui32(v14)>>(uint(v6)%32))
	if base.Ui32(v17) < base.Ui32(int32(256)) {
		v24 = v17
	} else {
		v24 = int32(base.Ui32(v17^int32(-1)) >> (uint(v6) % 32))
	}
	v25 = int32(24)
	v27 = int32(255)
	v34 = v5&v27 + v8&v27 - v14&v27
	if base.Ui32(v34) < base.Ui32(int32(256)) {
		v41 = v34
	} else {
		v41 = int32(base.Ui32(v34^int32(-1)) >> (uint(v25) % 32))
	}
	v43 = int32(16)
	v45 = int32(255)
	v56 = int32(base.Ui32(v5)>>(uint(v43)%32))&v45 + int32(base.Ui32(v8)>>(uint(v43)%32))&v45 - int32(base.Ui32(v14)>>(uint(v43)%32))&v45
	if base.Ui32(v56) < base.Ui32(int32(256)) {
		v63 = v56
	} else {
		v63 = int32(base.Ui32(v56^int32(-1)) >> (uint(int32(24)) % 32))
	}
	v67 = int32(8)
	v69 = int32(255)
	v80 = int32(base.Ui32(v5)>>(uint(v67)%32))&v69 + int32(base.Ui32(v8)>>(uint(v67)%32))&v69 - int32(base.Ui32(v14)>>(uint(v67)%32))&v69
	if base.Ui32(v80) < base.Ui32(int32(256)) {
		v87 = v80
	} else {
		v87 = int32(base.Ui32(v80^int32(-1)) >> (uint(int32(24)) % 32))
	}
	return v24<<(uint(v25)%32) + v41 + v63<<(uint(int32(16))%32) + v87<<(uint(int32(8))%32)
}
func F_VP8LPredictor13_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
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
	var v24 int32
	_ = v24
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v70 int32
	_ = v70
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v93 int32
	_ = v93
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v12 = int32(base.Ui32(v4^v5)>>(uint(int32(1))%32))&int32(2139062143) + v4&v5
	v13 = int32(24)
	v14 = int32(base.Ui32(v12) >> (uint(v13) % 32))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	v22 = base.I32_div_s(v14-int32(base.Ui32(v17)>>(uint(v13)%32)), int32(2))
	v24 = v14 + base.I32_extend16_s(v22)
	if base.Ui32(v24) < base.Ui32(int32(256)) {
		v31 = v24
	} else {
		v31 = int32(base.Ui32(v24^int32(-1)) >> (uint(v13) % 32))
	}
	v32 = int32(24)
	v34 = int32(255)
	v35 = v12 & v34
	v40 = base.I32_div_s(v35-v17&v34, int32(2))
	v42 = v35 + base.I32_extend16_s(v40)
	if base.Ui32(v42) < base.Ui32(int32(256)) {
		v49 = v42
	} else {
		v49 = int32(base.Ui32(v42^int32(-1)) >> (uint(v32) % 32))
	}
	v51 = int32(16)
	v53 = int32(255)
	v54 = int32(base.Ui32(v12)>>(uint(v51)%32)) & v53
	v61 = base.I32_div_s(v54-int32(base.Ui32(v17)>>(uint(v51)%32))&v53, int32(2))
	v63 = v54 + base.I32_extend16_s(v61)
	if base.Ui32(v63) < base.Ui32(int32(256)) {
		v70 = v63
	} else {
		v70 = int32(base.Ui32(v63^int32(-1)) >> (uint(int32(24)) % 32))
	}
	v74 = int32(8)
	v76 = int32(255)
	v77 = int32(base.Ui32(v12)>>(uint(v74)%32)) & v76
	v84 = base.I32_div_s(v77-int32(base.Ui32(v17)>>(uint(v74)%32))&v76, int32(2))
	v86 = v77 + base.I32_extend16_s(v84)
	if base.Ui32(v86) < base.Ui32(int32(256)) {
		v93 = v86
	} else {
		v93 = int32(base.Ui32(v86^int32(-1)) >> (uint(int32(24)) % 32))
	}
	return v31<<(uint(v32)%32) + v49 + v70<<(uint(int32(16))%32) + v93<<(uint(int32(8))%32)
}
func F_VP8LPredictor1_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	return v3
}
func F_VP8LPredictor2_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	return v3
}
func F_VP8LPredictor3_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	return v3
}
func F_VP8LPredictor4_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	return v5
}
func F_VP8LPredictor5_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v7 = int32(1)
	v9 = int32(2139062143)
	v12 = int32(base.Ui32(v4^v5)>>(uint(v7)%32))&v9 + v4&v5
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	return int32(base.Ui32(v12^v13)>>(uint(v7)%32))&v9 + v12&v13
}
func F_VP8LPredictor6_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	return int32(base.Ui32(v5^v6)>>(uint(int32(1))%32))&int32(2139062143) + v5&v6
}
func F_VP8LPredictor7_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	return int32(base.Ui32(v3^v4)>>(uint(int32(1))%32))&int32(2139062143) + v3&v4
}
func F_VP8LPredictor8_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l1+int32(-4))))
	return int32(base.Ui32(v4^v7)>>(uint(int32(1))%32))&int32(2139062143) + v4&v7
}
func F_VP8LPredictor9_C(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	return int32(base.Ui32(v4^v5)>>(uint(int32(1))%32))&int32(2139062143) + v4&v5
}
func F_VP8LSubtractGreenFromBlueAndRed_C(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v37 int32
	_ = v37
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v62 int32
	_ = v62
	var v67 int32
	_ = v67
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v78 int32
	_ = v78
	if l1 < int32(1) {
	} else {
		v10 = int32(1)
		if l1 == v10 {
			v67 = int32(0)
		} else {
			v19 = l0
			v21 = int32(0)
			for {
				v25 = *(*int32)(unsafe.Add(mBase, uint32(v19)))
				v26 = int32(8)
				v27 = int32(base.Ui32(v25) >> (uint(v26) % 32))
				v29 = int32(255)
				v31 = int32(-16711936)
				v34 = int32(16)
				v37 = int32(16711680)
				*(*int32)(unsafe.Add(mBase, uint32(v19))) = (v25-v27)&v29 | v25&v31 | (v25-v27<<(uint(v34)%32))&v37
				v42 = v19 + int32(4)
				v43 = *(*int32)(unsafe.Add(mBase, uint32(v42)))
				v45 = int32(base.Ui32(v43) >> (uint(v26) % 32))
				*(*int32)(unsafe.Add(mBase, uint32(v42))) = (v43-v45)&v29 | v43&v31 | (v43-v45<<(uint(v34)%32))&v37
				v62 = v21 + int32(2)
				if l1&int32(2147483646) != v62 {
					v19 = v19 + v26
					v21 = v62
					continue
				} else {
					break
				}
				break
			}
			v67 = v62
		}
		if l1&v10 == int32(0) {
		} else {
			v75 = l0 + v67<<(uint(int32(2))%32)
			v76 = *(*int32)(unsafe.Add(mBase, uint32(v75)))
			v78 = int32(base.Ui32(v76) >> (uint(int32(8)) % 32))
			*(*int32)(unsafe.Add(mBase, uint32(v75))) = (v76-v78)&int32(255) | v76&int32(-16711936) | (v76-v78<<(uint(int32(16))%32))&int32(16711680)
		}
	}
	return
}
func F_VP8LTransformColorInverse_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v16 int32
	_ = v16
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v55 int32
	_ = v55
	var v60 int32
	_ = v60
	if l2 < int32(1) {
	} else {
		v12 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
		v13 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
		v14 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
		v16 = l1
		v17 = l2
		v18 = l3
		for {
			v24 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
			v25 = int32(16)
			v28 = v24 << (uint(v25) % 32) >> (uint(int32(24)) % 32)
			v30 = int32(5)
			v34 = v28*v14>>(uint(v30)%32) + int32(base.Ui32(v24)>>(uint(v25)%32))
			*(*int32)(unsafe.Add(mBase, uint32(v18))) = v34<<(uint(v25)%32)&int32(16711680) | v24&int32(-16711936) | (int32(base.Ui32(v28*v13)>>(uint(v30)%32))+v24+int32(base.Ui32(base.I32_extend8_s(v34)*v12)>>(uint(v30)%32)))&int32(255)
			v55 = int32(4)
			v60 = v17 + int32(-1)
			if v60 != 0 {
				v16 = v16 + v55
				v17 = v60
				v18 = v18 + v55
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LTransformColor_C(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
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
	var v56 int32
	_ = v56
	if l2 < int32(1) {
	} else {
		v11 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+2)))
		v12 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0)+1)))
		v13 = int32(*(*int8)(unsafe.Add(mBase, uint32(l0))))
		v15 = l1
		v16 = l2
		for {
			v22 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
			v23 = int32(16)
			v24 = int32(base.Ui32(v22) >> (uint(v23) % 32))
			v28 = v22 << (uint(v23) % 32) >> (uint(int32(24)) % 32)
			v30 = int32(5)
			*(*int32)(unsafe.Add(mBase, uint32(v15))) = (v24-int32(base.Ui32(v28*v13)>>(uint(v30)%32)))<<(uint(v23)%32)&int32(16711680) | v22&int32(-16711936) | (v22-(int32(base.Ui32(v28*v12)>>(uint(v30)%32))+int32(base.Ui32(base.I32_extend8_s(v24)*v11)>>(uint(v30)%32))))&int32(255)
			v56 = v16 + int32(-1)
			if v56 != 0 {
				v15 = v15 + int32(4)
				v16 = v56
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VR4_C(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v24 int32
	_ = v24
	var v29 int32
	_ = v29
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v44 int32
	_ = v44
	var v49 int32
	_ = v49
	var v54 int32
	_ = v54
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	var v64 int32
	_ = v64
	var v71 int32
	_ = v71
	var v74 int32
	_ = v74
	var v80 int32
	_ = v80
	var v87 int32
	_ = v87
	var v97 int32
	_ = v97
	var v107 int32
	_ = v107
	v12 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-32)))))
	v15 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-33)))))
	v17 = int32(1)
	v20 = int32(base.Ui32(v12+v15+v17) >> (uint(v17) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+65)) = uint8(v20)
	v24 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-31)))))
	v29 = int32(base.Ui32(v12+v24+v17) >> (uint(v17) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+66)) = uint8(v29)
	*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v20)
	v34 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-30)))))
	v39 = int32(base.Ui32(v24+v34+v17) >> (uint(v17) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+67)) = uint8(v39)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v29)
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-29)))))
	v49 = int32(base.Ui32(v34+v44+v17) >> (uint(v17) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v49)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v39)
	v54 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(-1)))))
	v55 = int32(2)
	v56 = v54 + v55
	v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+63)))
	v59 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+31)))
	v64 = int32(base.Ui32(v56+v57+v59<<(uint(v17)%32)) >> (uint(v55) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+96)) = uint8(v64)
	v71 = int32(base.Ui32(v12+(v56+v15<<(uint(v17)%32))) >> (uint(v55) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+97)) = uint8(v71)
	v74 = v15 + v55
	v80 = int32(base.Ui32(v59+v74+v54<<(uint(v17)%32)) >> (uint(v55) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+64)) = uint8(v80)
	v87 = int32(base.Ui32(v24+(v74+v12<<(uint(v17)%32))) >> (uint(v55) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+98)) = uint8(v87)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+32)) = uint8(v71)
	v97 = int32(base.Ui32(v34+(v12+v24<<(uint(v17)%32))+v55) >> (uint(v55) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+99)) = uint8(v97)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+33)) = uint8(v87)
	v107 = int32(base.Ui32(v44+(v24+v34<<(uint(v17)%32))+v55) >> (uint(v55) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+35)) = uint8(v107)
	*(*uint8)(unsafe.Add(mBase, uint32(l0)+34)) = uint8(v97)
	return
}
func F_VectorMismatch_C(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	v4 = int32(0)
	if l2 < int32(1) {
		v25 = v4
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v25
L2:
	;
	v8 = l0
	v9 = l1
	v11 = v4
	goto L3
L3:
	;
	v12 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	if v12 != v13 {
		v25 = v11
		goto L1
	} else {
		goto L5
	}
L4:
	;
	v25 = l2
	goto L1
L5:
	;
	v15 = int32(4)
	v20 = v11 + int32(1)
	if l2 != v20 {
		v8 = v8 + v15
		v9 = v9 + v15
		v11 = v20
		goto L3
	} else {
		goto L6
	}
L6:
	;
	goto L4
}
func F_VerticalFilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v36 int32
	_ = v36
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v70 int32
	_ = v70
	var v79 int32
	_ = v79
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v123 int32
	_ = v123
	var v133 int32
	_ = v133
	var v157 int32
	_ = v157
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v175 int32
	_ = v175
	var v179 int32
	_ = v179
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v198 int32
	_ = v198
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v209 int32
	_ = v209
	var v212 int32
	_ = v212
	var v213 int32
	_ = v213
	var v215 int32
	_ = v215
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v228 int32
	_ = v228
	var v239 int32
	_ = v239
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v252 int32
	_ = v252
	var v254 int32
	_ = v254
	var v255 int32
	_ = v255
	var v259 int32
	_ = v259
	v15 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	*(*uint8)(unsafe.Add(mBase, uint32(l4))) = uint8(v15)
	if l1 < int32(2) {
	} else {
		v20 = l1 + int32(-1)
		if l1 == int32(2) {
			v70 = int32(0)
		} else {
			v36 = int32(0)
			for {
				v43 = l4 + v36
				v44 = int32(1)
				v46 = l0 + v36
				v49 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v46+v44))))
				v50 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v46))))
				v51 = v49 - v50
				*(*uint8)(unsafe.Add(mBase, uint32(v43+v44))) = uint8(v51)
				v53 = int32(2)
				v57 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v46+v53))))
				v58 = v57 - v49
				*(*uint8)(unsafe.Add(mBase, uint32(v43+v53))) = uint8(v58)
				v61 = v36 + v53
				if v20&int32(-2) != v61 {
					v36 = v61
					continue
				} else {
					break
				}
				break
			}
			v70 = v61
		}
		if v20&int32(1) == int32(0) {
		} else {
			v79 = int32(1)
			v85 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v79+v70))))
			v87 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v70))))
			v88 = v85 - v87
			*(*uint8)(unsafe.Add(mBase, uint32(l4+v79+v70))) = uint8(v88)
		}
	}
	if l2 < int32(2) {
	} else {
		if int32(0) < l1 {
			v166 = int32(1)
			v175 = l0
			v179 = l4
			v186 = l0 + v166
			v188 = v166
			for {
				if l1 == int32(1) {
					v239 = int32(0)
				} else {
					v198 = v186
					v200 = l3
					v202 = int32(0)
					for {
						v207 = v179 + v200
						v208 = v175 + v200
						v209 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v208))))
						v212 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198+int32(-1)))))
						v213 = v209 - v212
						*(*uint8)(unsafe.Add(mBase, uint32(v207))) = uint8(v213)
						v215 = int32(1)
						v219 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v208+v215))))
						v220 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v198))))
						v221 = v219 - v220
						*(*uint8)(unsafe.Add(mBase, uint32(v207+v215))) = uint8(v221)
						v223 = int32(2)
						v228 = v202 + int32(-2)
						if int32(0)-l1&int32(2147483646) != v228 {
							v198 = v198 + v223
							v200 = v200 + v223
							v202 = v228
							continue
						} else {
							break
						}
						break
					}
					v239 = int32(0) - v228
				}
				v246 = v175 + l3
				v247 = v179 + l3
				if l1&v166 == int32(0) {
				} else {
					v252 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v246+v239))))
					v254 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v175+v239))))
					v255 = v252 - v254
					*(*uint8)(unsafe.Add(mBase, uint32(v247+v239))) = uint8(v255)
				}
				v259 = v188 + int32(1)
				if v259 != l2 {
					v175 = v246
					v179 = v247
					v186 = v186 + l3
					v188 = v259
					continue
				} else {
					break
				}
				break
			}
		} else {
			v109 = l2 + int32(-1)
			v110 = int32(7)
			v111 = v109 & v110
			if base.Ui32(l2+int32(-2)) < base.Ui32(v110) {
			} else {
				v123 = v109 & int32(-8)
				for {
					v133 = v123 + int32(-8)
					if v133 != 0 {
						v123 = v133
						continue
					} else {
						break
					}
					break
				}
			}
			if v111 == int32(0) {
			} else {
				v157 = v111
				for {
					v165 = v157 + int32(-1)
					if v165 != 0 {
						v157 = v165
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
func F_VerticalUnfilter_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v14 int32
	_ = v14
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
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
	var v42 int32
	_ = v42
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v56 int32
	_ = v56
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v79 int32
	_ = v79
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v102 int32
	_ = v102
	var v106 int32
	_ = v106
	var v107 int32
	_ = v107
	var v109 int32
	_ = v109
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v132 int32
	_ = v132
	var v135 int32
	_ = v135
	var v137 int32
	_ = v137
	var v139 int32
	_ = v139
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v149 int32
	_ = v149
	if l0 == int32(0) {
		if l3 < int32(1) {
		} else {
			v72 = l3 & int32(3)
			v73 = int32(0)
			if base.Ui32(l3) < base.Ui32(int32(4)) {
				v124 = v73
				v126 = v73
			} else {
				v79 = int32(0)
				v86 = v79
				v88 = v79
				for {
					v90 = l2 + v88
					v91 = l1 + v88
					v92 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91))))
					v93 = v92 + v86
					*(*uint8)(unsafe.Add(mBase, uint32(v90))) = uint8(v93)
					v95 = int32(1)
					v99 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91+v95))))
					v100 = v99 + v93
					*(*uint8)(unsafe.Add(mBase, uint32(v90+v95))) = uint8(v100)
					v102 = int32(2)
					v106 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91+v102))))
					v107 = v106 + v100
					*(*uint8)(unsafe.Add(mBase, uint32(v90+v102))) = uint8(v107)
					v109 = int32(3)
					v113 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v91+v109))))
					v114 = v113 + v107
					*(*uint8)(unsafe.Add(mBase, uint32(v90+v109))) = uint8(v114)
					v117 = v88 + int32(4)
					if l3&int32(2147483644) != v117 {
						v86 = v114
						v88 = v117
						continue
					} else {
						break
					}
					break
				}
				v124 = v114
				v126 = v117
			}
			if v72 == int32(0) {
			} else {
				v132 = v72
				v135 = l1 + v126
				v137 = v124
				v139 = l2 + v126
				for {
					v141 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v135))))
					v142 = v141 + v137
					*(*uint8)(unsafe.Add(mBase, uint32(v139))) = uint8(v142)
					v144 = int32(1)
					v149 = v132 + int32(-1)
					if v149 != 0 {
						v132 = v149
						v135 = v135 + v144
						v137 = v142
						v139 = v139 + v144
						continue
					} else {
						break
					}
					break
				}
			}
		}
		return
	} else {
		if l3 < int32(1) {
			return
		} else {
			v14 = int32(1)
			if l3 == v14 {
				v56 = int32(0)
			} else {
				v27 = int32(0)
				for {
					v31 = l2 + v27
					v32 = l1 + v27
					v33 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v32))))
					v34 = l0 + v27
					v35 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v34))))
					v36 = v33 + v35
					*(*uint8)(unsafe.Add(mBase, uint32(v31))) = uint8(v36)
					v38 = int32(1)
					v42 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v32+v38))))
					v45 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v34+v38))))
					v46 = v42 + v45
					*(*uint8)(unsafe.Add(mBase, uint32(v31+v38))) = uint8(v46)
					v49 = v27 + int32(2)
					if l3&int32(2147483646) != v49 {
						v27 = v49
						continue
					} else {
						break
					}
					break
				}
				v56 = v49
			}
			if l3&v14 == int32(0) {
				return
			} else {
				v64 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v56))))
				v66 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v56))))
				v67 = v64 + v66
				*(*uint8)(unsafe.Add(mBase, uint32(l2+v56))) = uint8(v67)
				return
			}
		}
	}
}

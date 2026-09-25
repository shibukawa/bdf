//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_VP8LAllocateHistogram(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v10 int32
	_ = v10
	var v11 int64
	_ = v11
	var v12 int32
	_ = v12
	var v19 int64
	_ = v19
	var v20 int32
	_ = v20
	var v21 int64
	_ = v21
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v42 int64
	_ = v42
	var v49 int32
	_ = v49
	var v53 int32
	_ = v53
	v5 = int32(_a_F_VP8LAllocateHistogram_0)
	if int32(0) < l0 {
		v10 = int32(4)<<(uint(l0)%32) + v5
	} else {
		v10 = v5
	}
	v11 = base.I64_extend_i32_s(v10)
	v12 = int32(1)
	if v11 == int64(0) {
		v31 = F_malloc(m, base.I32_wrap_i64(v11)*v12)
		mBase = m.M
		v33 = v31
	} else {
		v19 = base.I64_div_u_s(int64(2147418112), v11)
		v20 = int32(0)
		v21 = base.I64_extend_i32_u(v12)
		if base.Ui64(int64(4294967295)) < base.Ui64(v21*v11) {
			v33 = v20
		} else {
			if base.Ui64(v19) < base.Ui64(v21) {
				v33 = v20
			} else {
				v31 = F_malloc(m, base.I32_wrap_i64(v11)*v12)
				mBase = m.M
				v33 = v31
			}
		}
	}
	if v33 == int32(0) {
	} else {
		*(*int64)(unsafe.Add(mBase, uint32(v33)+3240)) = int64(-1)
		*(*int32)(unsafe.Add(mBase, uint32(v33)+3236)) = l0
		*(*int32)(unsafe.Add(mBase, uint32(v33)+3304)) = int32(16843009)
		v42 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(v33)+3256)) = v42
		*(*int32)(unsafe.Add(mBase, uint32(v33))) = v33 + int32(3312)
		v49 = int32(-1)
		*(*uint16)(unsafe.Add(mBase, uint32(v33+int32(3248)))) = uint16(v49)
		v53 = int32(1)
		*(*uint8)(unsafe.Add(mBase, uint32(v33+int32(3308)))) = uint8(v53)
		*(*int64)(unsafe.Add(mBase, uint32(v33+int32(3264)))) = v42
		*(*int64)(unsafe.Add(mBase, uint32(v33+int32(3272)))) = v42
		*(*int64)(unsafe.Add(mBase, uint32(v33+int32(3280)))) = v42
		*(*int64)(unsafe.Add(mBase, uint32(v33+int32(3288)))) = v42
		*(*int64)(unsafe.Add(mBase, uint32(v33+int32(3296)))) = v42
	}
	return v33
}
func F_VP8LAllocateHistogramSet(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v23 int64
	_ = v23
	var v24 int32
	_ = v24
	var v31 int64
	_ = v31
	var v32 int32
	_ = v32
	var v33 int64
	_ = v33
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v52 int32
	_ = v52
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v67 int32
	_ = v67
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v96 int32
	_ = v96
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v138 int32
	_ = v138
	var v140 int32
	_ = v140
	var v160 int32
	_ = v160
	var v162 int32
	_ = v162
	var v166 int32
	_ = v166
	var v170 int32
	_ = v170
	var v171 int64
	_ = v171
	var v173 int32
	_ = v173
	var v175 int64
	_ = v175
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v189 int32
	_ = v189
	var v193 int32
	_ = v193
	var v197 int32
	_ = v197
	var v201 int32
	_ = v201
	var v206 int32
	_ = v206
	var v208 int32
	_ = v208
	var v247 int32
	_ = v247
	var v254 int32
	_ = v254
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v272 int64
	_ = v272
	var v276 int32
	_ = v276
	var v280 int32
	_ = v280
	v13 = int32(4) << (uint(l1) % 32)
	v14 = int32(_a_F_VP8LAllocateHistogramSet_0)
	v18 = base.B2i32(int32(0) < l1)
	if int32(0) < l1 {
		v19 = v13 + v14
	} else {
		v19 = v14
	}
	v23 = base.I64_extend_i32_u(v19*l0 + int32(12))
	v24 = int32(1)
	if v23 == int64(0) {
		v43 = F_malloc(m, base.I32_wrap_i64(v23)*v24)
		mBase = m.M
		v45 = v43
	} else {
		v31 = base.I64_div_u_s(int64(2147418112), v23)
		v32 = int32(0)
		v33 = base.I64_extend_i32_u(v24)
		if base.Ui64(int64(4294967295)) < base.Ui64(v33*v23) {
			v45 = v32
		} else {
			if base.Ui64(v31) < base.Ui64(v33) {
				v45 = v32
			} else {
				v43 = F_malloc(m, base.I32_wrap_i64(v23)*v24)
				mBase = m.M
				v45 = v43
			}
		}
	}
	if v45 == int32(0) {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(v45)+4)) = l0
		*(*int32)(unsafe.Add(mBase, uint32(v45))) = l0
		v52 = v45 + int32(12)
		*(*int32)(unsafe.Add(mBase, uint32(v45)+8)) = v52
		if l0 < int32(1) {
		} else {
			v60 = v52 + l0<<(uint(int32(2))%32)
			v62 = l0 + int32(-1)
			if v62 != 0 {
				v64 = int32(_a_F_VP8LAllocateHistogramSet_1)
				if int32(0) < l1 {
					v67 = v13 + v64
				} else {
					v67 = v64
				}
				v74 = int32(4)
				v75 = int32(0)
				v77 = v52
				v79 = v60
				for {
					v84 = int32(-4)
					v86 = int32(31)
					v88 = int32(-32)
					v89 = (v79 + v86) & v88
					*(*int32)(unsafe.Add(mBase, uint32(v77+v74+v84))) = v89
					v91 = *(*int32)(unsafe.Add(mBase, uint32(v45)+8))
					v92 = v91 + v74
					v95 = *(*int32)(unsafe.Add(mBase, uint32(v92+v84)))
					v96 = int32(3312)
					*(*int32)(unsafe.Add(mBase, uint32(v95))) = v89 + v96
					v103 = (v89 + v67 + v86) & v88
					*(*int32)(unsafe.Add(mBase, uint32(v92))) = v103
					v105 = *(*int32)(unsafe.Add(mBase, uint32(v45)+8))
					v107 = *(*int32)(unsafe.Add(mBase, uint32(v105+v74)))
					*(*int32)(unsafe.Add(mBase, uint32(v107))) = v103 + v96
					v113 = v103 + v67
					v115 = v75 + int32(2)
					if l0&int32(2147483646) != v115 {
						v74 = v74 + int32(8)
						v75 = v115
						v77 = v105
						v79 = v113
						continue
					} else {
						break
					}
					break
				}
				v120 = v115
				v122 = v105
				v124 = v113
			} else {
				v120 = int32(0)
				v122 = v52
				v124 = v60
			}
			if l0&int32(1) == int32(0) {
			} else {
				v131 = v120 << (uint(int32(2)) % 32)
				v136 = (v124 + int32(31)) & int32(-32)
				*(*int32)(unsafe.Add(mBase, uint32(v122+v131))) = v136
				v138 = *(*int32)(unsafe.Add(mBase, uint32(v45)+8))
				v140 = *(*int32)(unsafe.Add(mBase, uint32(v138+v131)))
				*(*int32)(unsafe.Add(mBase, uint32(v140))) = v136 + int32(3312)
			}
			if l0 < int32(1) {
			} else {
				if v62 != 0 {
					v160 = int32(0)
					v162 = int32(4)
					for {
						v166 = *(*int32)(unsafe.Add(mBase, uint32(v45)+8))
						v170 = *(*int32)(unsafe.Add(mBase, uint32(v166+v162+int32(-4))))
						v171 = int64(-1)
						*(*int64)(unsafe.Add(mBase, uint32(v170)+3240)) = v171
						v173 = int32(16843009)
						*(*int32)(unsafe.Add(mBase, uint32(v170)+3304)) = v173
						v175 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(v170)+3256)) = v175
						v177 = int32(3248)
						v179 = int32(-1)
						*(*uint16)(unsafe.Add(mBase, uint32(v170+v177))) = uint16(v179)
						v181 = int32(3308)
						v183 = int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v170+v181))) = uint8(v183)
						v185 = int32(3264)
						*(*int64)(unsafe.Add(mBase, uint32(v170+v185))) = v175
						v189 = int32(3272)
						*(*int64)(unsafe.Add(mBase, uint32(v170+v189))) = v175
						v193 = int32(3280)
						*(*int64)(unsafe.Add(mBase, uint32(v170+v193))) = v175
						v197 = int32(3288)
						*(*int64)(unsafe.Add(mBase, uint32(v170+v197))) = v175
						v201 = int32(3296)
						*(*int64)(unsafe.Add(mBase, uint32(v170+v201))) = v175
						*(*int32)(unsafe.Add(mBase, uint32(v170)+3236)) = l1
						v206 = *(*int32)(unsafe.Add(mBase, uint32(v45)+8))
						v208 = *(*int32)(unsafe.Add(mBase, uint32(v206+v162)))
						*(*int64)(unsafe.Add(mBase, uint32(v208)+3240)) = v171
						*(*int32)(unsafe.Add(mBase, uint32(v208)+3304)) = v173
						*(*int64)(unsafe.Add(mBase, uint32(v208)+3256)) = v175
						*(*int32)(unsafe.Add(mBase, uint32(v208)+3236)) = l1
						*(*uint16)(unsafe.Add(mBase, uint32(v208+v177))) = uint16(v179)
						*(*uint8)(unsafe.Add(mBase, uint32(v208+v181))) = uint8(v183)
						*(*int64)(unsafe.Add(mBase, uint32(v208+v185))) = v175
						*(*int64)(unsafe.Add(mBase, uint32(v208+v189))) = v175
						*(*int64)(unsafe.Add(mBase, uint32(v208+v193))) = v175
						*(*int64)(unsafe.Add(mBase, uint32(v208+v197))) = v175
						*(*int64)(unsafe.Add(mBase, uint32(v208+v201))) = v175
						v247 = v160 + int32(2)
						if l0&int32(2147483646) != v247 {
							v160 = v247
							v162 = v162 + int32(8)
							continue
						} else {
							break
						}
						break
					}
					v254 = v247
				} else {
					v254 = int32(0)
				}
				if l0&int32(1) == int32(0) {
				} else {
					v262 = *(*int32)(unsafe.Add(mBase, uint32(v45)+8))
					v266 = *(*int32)(unsafe.Add(mBase, uint32(v262+v254<<(uint(int32(2))%32))))
					*(*int64)(unsafe.Add(mBase, uint32(v266)+3240)) = int64(-1)
					*(*int32)(unsafe.Add(mBase, uint32(v266)+3236)) = l1
					*(*int32)(unsafe.Add(mBase, uint32(v266)+3304)) = int32(16843009)
					v272 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(v266)+3256)) = v272
					v276 = int32(-1)
					*(*uint16)(unsafe.Add(mBase, uint32(v266+int32(3248)))) = uint16(v276)
					v280 = int32(1)
					*(*uint8)(unsafe.Add(mBase, uint32(v266+int32(3308)))) = uint8(v280)
					*(*int64)(unsafe.Add(mBase, uint32(v266+int32(3264)))) = v272
					*(*int64)(unsafe.Add(mBase, uint32(v266+int32(3272)))) = v272
					*(*int64)(unsafe.Add(mBase, uint32(v266+int32(3280)))) = v272
					*(*int64)(unsafe.Add(mBase, uint32(v266+int32(3288)))) = v272
					*(*int64)(unsafe.Add(mBase, uint32(v266+int32(3296)))) = v272
				}
			}
		}
	}
	return v45
}
func F_VP8LBackwardReferencesTraceBackwards(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32 {
	mBase := m.M
	_ = mBase
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v65 int64
	_ = v65
	var v66 int32
	_ = v66
	var v73 int64
	_ = v73
	var v74 int32
	_ = v74
	var v75 int64
	_ = v75
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	var v90 int64
	_ = v90
	var v91 int32
	_ = v91
	var v98 int32
	_ = v98
	var v105 int64
	_ = v105
	var v106 int32
	_ = v106
	var v107 int64
	_ = v107
	var v116 int32
	_ = v116
	var v118 int32
	_ = v118
	var v121 int64
	_ = v121
	var v122 int32
	_ = v122
	var v129 int64
	_ = v129
	var v130 int32
	_ = v130
	var v131 int64
	_ = v131
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v157 int32
	_ = v157
	var v163 int32
	_ = v163
	var v171 int32
	_ = v171
	var v178 int32
	_ = v178
	var v183 int32
	_ = v183
	var v186 int32
	_ = v186
	var v194 int64
	_ = v194
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v238 int32
	_ = v238
	var v243 int32
	_ = v243
	var v244 int32
	_ = v244
	var v253 int64
	_ = v253
	var v257 int32
	_ = v257
	var v261 int32
	_ = v261
	var v283 int32
	_ = v283
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v292 int32
	_ = v292
	var v295 int32
	_ = v295
	var v301 int32
	_ = v301
	var v303 int32
	_ = v303
	var v305 int32
	_ = v305
	var v310 int32
	_ = v310
	var v311 int32
	_ = v311
	var v320 int32
	_ = v320
	var v322 int32
	_ = v322
	var v324 int32
	_ = v324
	var v329 int32
	_ = v329
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v337 int32
	_ = v337
	var v340 int32
	_ = v340
	var v356 int32
	_ = v356
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v400 int32
	_ = v400
	var v403 int32
	_ = v403
	var v406 int32
	_ = v406
	var v407 int32
	_ = v407
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v423 int32
	_ = v423
	var v425 int32
	_ = v425
	var v427 int32
	_ = v427
	var v431 int32
	_ = v431
	var v448 int32
	_ = v448
	var v450 int32
	_ = v450
	var v451 int32
	_ = v451
	var v489 int32
	_ = v489
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v497 int32
	_ = v497
	var v518 int32
	_ = v518
	var v519 int32
	_ = v519
	var v618 int32
	_ = v618
	var v620 int32
	_ = v620
	var v631 int32
	_ = v631
	var v653 int32
	_ = v653
	var v656 int32
	_ = v656
	var v657 int32
	_ = v657
	var v665 int32
	_ = v665
	var v666 int32
	_ = v666
	var v701 int32
	_ = v701
	var v702 int32
	_ = v702
	var v707 int64
	_ = v707
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v723 int32
	_ = v723
	var v743 int32
	_ = v743
	var v747 int32
	_ = v747
	var v748 int32
	_ = v748
	var v749 int32
	_ = v749
	var v750 int32
	_ = v750
	var v751 int32
	_ = v751
	var v767 int32
	_ = v767
	var v768 int32
	_ = v768
	var v771 int32
	_ = v771
	var v811 int32
	_ = v811
	var v814 int32
	_ = v814
	var v818 int32
	_ = v818
	var v819 int32
	_ = v819
	var v820 int32
	_ = v820
	var v821 int32
	_ = v821
	var v822 int32
	_ = v822
	var v825 int32
	_ = v825
	var v830 int32
	_ = v830
	var v892 int32
	_ = v892
	var v912 int32
	_ = v912
	var v915 int32
	_ = v915
	var v916 int32
	_ = v916
	var v954 int32
	_ = v954
	var v955 int32
	_ = v955
	var v957 int32
	_ = v957
	var v960 int32
	_ = v960
	var v963 int32
	_ = v963
	var v966 int32
	_ = v966
	var v970 int32
	_ = v970
	var v971 int32
	_ = v971
	var v982 int32
	_ = v982
	var v984 int32
	_ = v984
	var v988 int32
	_ = v988
	var v1115 int32
	_ = v1115
	var v1119 int32
	_ = v1119
	var v1120 int32
	_ = v1120
	var v1121 int32
	_ = v1121
	var v1122 int32
	_ = v1122
	var v1123 int32
	_ = v1123
	var v1126 int32
	_ = v1126
	var v1185 int32
	_ = v1185
	var v1188 int32
	_ = v1188
	var v1192 int32
	_ = v1192
	var v1193 int32
	_ = v1193
	var v1194 int32
	_ = v1194
	var v1195 int32
	_ = v1195
	var v1196 int32
	_ = v1196
	var v1201 int32
	_ = v1201
	var v1265 int32
	_ = v1265
	var v1285 int32
	_ = v1285
	var v1288 int32
	_ = v1288
	var v1289 int32
	_ = v1289
	var v1327 int32
	_ = v1327
	var v1330 int32
	_ = v1330
	var v1333 int32
	_ = v1333
	var v1336 int32
	_ = v1336
	var v1339 int32
	_ = v1339
	var v1343 int32
	_ = v1343
	var v1344 int32
	_ = v1344
	var v1355 int32
	_ = v1355
	var v1357 int32
	_ = v1357
	var v1361 int32
	_ = v1361
	var v1488 int32
	_ = v1488
	var v1492 int32
	_ = v1492
	var v1493 int32
	_ = v1493
	var v1494 int32
	_ = v1494
	var v1495 int32
	_ = v1495
	var v1496 int32
	_ = v1496
	var v1499 int32
	_ = v1499
	var v1558 int32
	_ = v1558
	var v1561 int32
	_ = v1561
	var v1565 int32
	_ = v1565
	var v1566 int32
	_ = v1566
	var v1567 int32
	_ = v1567
	var v1568 int32
	_ = v1568
	var v1569 int32
	_ = v1569
	var v1574 int32
	_ = v1574
	var v1638 int32
	_ = v1638
	var v1658 int32
	_ = v1658
	var v1661 int32
	_ = v1661
	var v1662 int32
	_ = v1662
	var v1700 int32
	_ = v1700
	var v1703 int32
	_ = v1703
	var v1706 int32
	_ = v1706
	var v1709 int32
	_ = v1709
	var v1712 int32
	_ = v1712
	var v1716 int32
	_ = v1716
	var v1717 int32
	_ = v1717
	var v1728 int32
	_ = v1728
	var v1730 int32
	_ = v1730
	var v1859 int32
	_ = v1859
	var v1863 int32
	_ = v1863
	var v1864 int32
	_ = v1864
	var v1865 int32
	_ = v1865
	var v1866 int32
	_ = v1866
	var v1867 int32
	_ = v1867
	var v1870 int32
	_ = v1870
	var v1929 int32
	_ = v1929
	var v1932 int32
	_ = v1932
	var v1936 int32
	_ = v1936
	var v1937 int32
	_ = v1937
	var v1938 int32
	_ = v1938
	var v1939 int32
	_ = v1939
	var v1940 int32
	_ = v1940
	var v1945 int32
	_ = v1945
	var v2008 int32
	_ = v2008
	var v2009 int32
	_ = v2009
	var v2010 int32
	_ = v2010
	var v2012 int32
	_ = v2012
	var v2016 int32
	_ = v2016
	var v2020 int32
	_ = v2020
	var v2024 int32
	_ = v2024
	var v2028 int32
	_ = v2028
	var v2032 int32
	_ = v2032
	var v2036 int32
	_ = v2036
	var v2040 int32
	_ = v2040
	var v2044 int32
	_ = v2044
	var v2048 int32
	_ = v2048
	var v2052 int32
	_ = v2052
	var v2056 int32
	_ = v2056
	var v2060 int32
	_ = v2060
	var v2064 int32
	_ = v2064
	var v2068 int32
	_ = v2068
	var v2072 int32
	_ = v2072
	var v2076 int32
	_ = v2076
	var v2080 int32
	_ = v2080
	var v2084 int32
	_ = v2084
	var v2088 int32
	_ = v2088
	var v2092 int32
	_ = v2092
	var v2096 int32
	_ = v2096
	var v2100 int32
	_ = v2100
	var v2104 int32
	_ = v2104
	var v2108 int32
	_ = v2108
	var v2112 int32
	_ = v2112
	var v2116 int32
	_ = v2116
	var v2120 int32
	_ = v2120
	var v2124 int32
	_ = v2124
	var v2128 int32
	_ = v2128
	var v2132 int32
	_ = v2132
	var v2136 int32
	_ = v2136
	var v2140 int32
	_ = v2140
	var v2144 int32
	_ = v2144
	var v2148 int32
	_ = v2148
	var v2152 int32
	_ = v2152
	var v2156 int32
	_ = v2156
	var v2160 int32
	_ = v2160
	var v2164 int32
	_ = v2164
	var v2330 int32
	_ = v2330
	var v2333 int32
	_ = v2333
	var v2337 int32
	_ = v2337
	var v2338 int32
	_ = v2338
	var v2339 int32
	_ = v2339
	var v2340 int32
	_ = v2340
	var v2341 int32
	_ = v2341
	var v2346 int32
	_ = v2346
	var v2405 int32
	_ = v2405
	var v2408 int32
	_ = v2408
	var v2412 int32
	_ = v2412
	var v2413 int32
	_ = v2413
	var v2414 int32
	_ = v2414
	var v2415 int32
	_ = v2415
	var v2416 int32
	_ = v2416
	var v2421 int32
	_ = v2421
	var v2484 int32
	_ = v2484
	var v2492 int32
	_ = v2492
	var v2519 int32
	_ = v2519
	var v2526 int32
	_ = v2526
	var v2529 int32
	_ = v2529
	var v2535 int32
	_ = v2535
	var v2539 int32
	_ = v2539
	var v2541 int32
	_ = v2541
	var v2543 int32
	_ = v2543
	var v2556 int32
	_ = v2556
	var v2562 int32
	_ = v2562
	var v2605 int32
	_ = v2605
	var v2606 int32
	_ = v2606
	var v2607 int32
	_ = v2607
	var v2608 int32
	_ = v2608
	var v2611 int32
	_ = v2611
	var v2613 int32
	_ = v2613
	var v2615 int32
	_ = v2615
	var v2620 int32
	_ = v2620
	var v2621 int32
	_ = v2621
	var v2625 int32
	_ = v2625
	var v2628 int64
	_ = v2628
	var v2636 int32
	_ = v2636
	var v2638 int32
	_ = v2638
	var v2646 int32
	_ = v2646
	var v2649 int64
	_ = v2649
	var v2660 int32
	_ = v2660
	var v2673 int32
	_ = v2673
	var v2679 int32
	_ = v2679
	var v2716 int64
	_ = v2716
	var v2720 int64
	_ = v2720
	var v2723 int32
	_ = v2723
	var v2725 int32
	_ = v2725
	var v2726 int64
	_ = v2726
	var v2729 int32
	_ = v2729
	var v2731 int32
	_ = v2731
	var v2735 int32
	_ = v2735
	var v2740 int32
	_ = v2740
	var v2753 int32
	_ = v2753
	var v2796 int64
	_ = v2796
	var v2803 int64
	_ = v2803
	var v2806 int32
	_ = v2806
	var v2808 int32
	_ = v2808
	var v2825 int32
	_ = v2825
	var v2868 int64
	_ = v2868
	var v2870 int32
	_ = v2870
	var v2877 int64
	_ = v2877
	var v2878 int32
	_ = v2878
	var v2879 int64
	_ = v2879
	var v2889 int32
	_ = v2889
	var v2891 int32
	_ = v2891
	var v2898 int64
	_ = v2898
	var v2902 int32
	_ = v2902
	var v2904 int32
	_ = v2904
	var v2915 int32
	_ = v2915
	var v2928 int32
	_ = v2928
	var v2934 int32
	_ = v2934
	var v2971 int64
	_ = v2971
	var v2975 int64
	_ = v2975
	var v2983 int32
	_ = v2983
	var v2985 int32
	_ = v2985
	var v2987 int64
	_ = v2987
	var v2993 int32
	_ = v2993
	var v3003 int32
	_ = v3003
	var v3016 int32
	_ = v3016
	var v3059 int64
	_ = v3059
	var v3068 int64
	_ = v3068
	var v3074 int32
	_ = v3074
	var v3137 int32
	_ = v3137
	var v3144 int64
	_ = v3144
	var v3145 int32
	_ = v3145
	var v3146 int64
	_ = v3146
	var v3156 int32
	_ = v3156
	var v3158 int32
	_ = v3158
	var v3166 int32
	_ = v3166
	var v3174 int32
	_ = v3174
	var v3193 int32
	_ = v3193
	var v3232 int64
	_ = v3232
	var v3258 int32
	_ = v3258
	var v3265 int32
	_ = v3265
	var v3287 int32
	_ = v3287
	var v3332 int32
	_ = v3332
	var v3345 int32
	_ = v3345
	var v3395 int32
	_ = v3395
	var v3455 int32
	_ = v3455
	var v3457 int32
	_ = v3457
	var v3460 int32
	_ = v3460
	var v3463 int32
	_ = v3463
	var v3464 int32
	_ = v3464
	var v3467 int32
	_ = v3467
	var v3470 int32
	_ = v3470
	var v3472 int32
	_ = v3472
	var v3478 int64
	_ = v3478
	var v3486 int32
	_ = v3486
	var v3489 int64
	_ = v3489
	var v3495 int64
	_ = v3495
	var v3497 int32
	_ = v3497
	var v3503 int64
	_ = v3503
	var v3510 int64
	_ = v3510
	var v3516 int64
	_ = v3516
	var v3517 int64
	_ = v3517
	var v3521 int64
	_ = v3521
	var v3523 int32
	_ = v3523
	var v3529 int32
	_ = v3529
	var v3533 int32
	_ = v3533
	var v3545 int64
	_ = v3545
	var v3552 int32
	_ = v3552
	var v3556 int32
	_ = v3556
	var v3561 int32
	_ = v3561
	var v3563 int32
	_ = v3563
	var v3564 int32
	_ = v3564
	var v3566 int32
	_ = v3566
	var v3568 int32
	_ = v3568
	var v3595 int32
	_ = v3595
	var v3597 int32
	_ = v3597
	var v3599 int32
	_ = v3599
	var v3601 int32
	_ = v3601
	var v3602 int32
	_ = v3602
	var v3604 int32
	_ = v3604
	var v3608 int64
	_ = v3608
	var v3611 int32
	_ = v3611
	var v3614 int32
	_ = v3614
	var v3615 int32
	_ = v3615
	var v3617 int32
	_ = v3617
	var v3618 int32
	_ = v3618
	var v3621 int32
	_ = v3621
	var v3623 int32
	_ = v3623
	var v3627 int64
	_ = v3627
	var v3636 int32
	_ = v3636
	var v3639 int64
	_ = v3639
	var v3645 int64
	_ = v3645
	var v3647 int32
	_ = v3647
	var v3653 int64
	_ = v3653
	var v3660 int64
	_ = v3660
	var v3667 int64
	_ = v3667
	var v3669 int32
	_ = v3669
	var v3671 int32
	_ = v3671
	var v3672 int32
	_ = v3672
	var v3673 int64
	_ = v3673
	var v3677 int64
	_ = v3677
	var v3678 int64
	_ = v3678
	var v3681 int32
	_ = v3681
	var v3687 int32
	_ = v3687
	var v3693 int32
	_ = v3693
	var v3695 int32
	_ = v3695
	var v3700 int32
	_ = v3700
	var v3709 int32
	_ = v3709
	var v3721 int32
	_ = v3721
	var v3727 int32
	_ = v3727
	var v3732 int32
	_ = v3732
	var v3735 int32
	_ = v3735
	var v3738 int32
	_ = v3738
	var v3739 int32
	_ = v3739
	var v3740 int32
	_ = v3740
	var v3741 int32
	_ = v3741
	var v3742 int32
	_ = v3742
	var v3745 int32
	_ = v3745
	var v3747 int32
	_ = v3747
	var v3749 int32
	_ = v3749
	var v3754 int32
	_ = v3754
	var v3755 int32
	_ = v3755
	var v3762 int64
	_ = v3762
	var v3763 int64
	_ = v3763
	var v3770 int32
	_ = v3770
	var v3795 int32
	_ = v3795
	var v3805 int32
	_ = v3805
	var v3841 int32
	_ = v3841
	var v3842 int32
	_ = v3842
	var v3847 int32
	_ = v3847
	var v3850 int32
	_ = v3850
	var v3856 int32
	_ = v3856
	var v3887 int32
	_ = v3887
	var v3915 int32
	_ = v3915
	var v3916 int32
	_ = v3916
	var v3921 int32
	_ = v3921
	var v3946 int32
	_ = v3946
	var v3981 int32
	_ = v3981
	var v3983 int32
	_ = v3983
	var v3984 int32
	_ = v3984
	var v3986 int64
	_ = v3986
	var v3987 int64
	_ = v3987
	var v3990 int32
	_ = v3990
	var v3991 int32
	_ = v3991
	var v3994 int32
	_ = v3994
	var v3997 int32
	_ = v3997
	var v4004 int32
	_ = v4004
	var v4026 int32
	_ = v4026
	var v4064 int32
	_ = v4064
	var v4066 int32
	_ = v4066
	var v4067 int32
	_ = v4067
	var v4069 int64
	_ = v4069
	var v4070 int64
	_ = v4070
	var v4073 int32
	_ = v4073
	var v4074 int32
	_ = v4074
	var v4077 int32
	_ = v4077
	var v4080 int32
	_ = v4080
	var v4145 int64
	_ = v4145
	var v4159 int64
	_ = v4159
	var v4180 int32
	_ = v4180
	var v4182 int32
	_ = v4182
	var v4209 int32
	_ = v4209
	var v4226 int32
	_ = v4226
	var v4271 int32
	_ = v4271
	var v4273 int32
	_ = v4273
	var v4274 int32
	_ = v4274
	var v4276 int32
	_ = v4276
	var v4279 int32
	_ = v4279
	var v4286 int32
	_ = v4286
	var v4288 int32
	_ = v4288
	var v4290 int32
	_ = v4290
	var v4292 int32
	_ = v4292
	var v4296 int32
	_ = v4296
	var v4297 int32
	_ = v4297
	var v4298 int64
	_ = v4298
	var v4299 int64
	_ = v4299
	var v4302 int32
	_ = v4302
	var v4303 int32
	_ = v4303
	var v4306 int32
	_ = v4306
	var v4309 int32
	_ = v4309
	var v4372 int32
	_ = v4372
	var v4377 int32
	_ = v4377
	var v4438 int32
	_ = v4438
	var v4441 int32
	_ = v4441
	var v4509 int32
	_ = v4509
	var v4511 int32
	_ = v4511
	var v4513 int32
	_ = v4513
	var v4521 int32
	_ = v4521
	var v4525 int32
	_ = v4525
	var v4534 int32
	_ = v4534
	var v4536 int32
	_ = v4536
	var v4544 int32
	_ = v4544
	var v4548 int32
	_ = v4548
	var v4559 int32
	_ = v4559
	var v4597 int32
	_ = v4597
	var v4658 int32
	_ = v4658
	var v4661 int32
	_ = v4661
	var v4666 int32
	_ = v4666
	var v4730 int32
	_ = v4730
	var v4732 int32
	_ = v4732
	var v4734 int32
	_ = v4734
	var v4742 int32
	_ = v4742
	var v4746 int32
	_ = v4746
	var v4755 int32
	_ = v4755
	var v4757 int32
	_ = v4757
	var v4765 int32
	_ = v4765
	var v4769 int32
	_ = v4769
	var v4780 int32
	_ = v4780
	var v4822 int32
	_ = v4822
	var v4824 int32
	_ = v4824
	var v4827 int32
	_ = v4827
	var v4847 int32
	_ = v4847
	var v4886 int32
	_ = v4886
	var v4887 int32
	_ = v4887
	var v4891 int32
	_ = v4891
	var v4914 int32
	_ = v4914
	var v4953 int32
	_ = v4953
	var v4954 int32
	_ = v4954
	var v4958 int32
	_ = v4958
	var v4964 int32
	_ = v4964
	var v4975 int32
	_ = v4975
	var v4978 int32
	_ = v4978
	var v4980 int32
	_ = v4980
	var v4982 int32
	_ = v4982
	var v4989 int32
	_ = v4989
	var v4992 int32
	_ = v4992
	var v4996 int32
	_ = v4996
	var v4999 int32
	_ = v4999
	var v5001 int32
	_ = v5001
	var v5003 int32
	_ = v5003
	var v5010 int32
	_ = v5010
	var v5013 int32
	_ = v5013
	var v5017 int32
	_ = v5017
	var v5037 int32
	_ = v5037
	var v5038 int32
	_ = v5038
	var v5075 int32
	_ = v5075
	var v5079 int32
	_ = v5079
	var v5081 int32
	_ = v5081
	var v5085 int32
	_ = v5085
	var v5089 int64
	_ = v5089
	var v5095 int32
	_ = v5095
	var v5098 int32
	_ = v5098
	var v5099 int32
	_ = v5099
	var v5102 int32
	_ = v5102
	var v5104 int32
	_ = v5104
	var v5109 int32
	_ = v5109
	var v5110 int32
	_ = v5110
	var v5117 int32
	_ = v5117
	var v5119 int32
	_ = v5119
	var v5120 int32
	_ = v5120
	var v5122 int32
	_ = v5122
	var v5129 int32
	_ = v5129
	var v5130 int32
	_ = v5130
	var v5134 int32
	_ = v5134
	var v5138 int64
	_ = v5138
	var v5145 int64
	_ = v5145
	var v5149 int64
	_ = v5149
	var v5154 int32
	_ = v5154
	var v5157 int32
	_ = v5157
	var v5158 int32
	_ = v5158
	var v5161 int32
	_ = v5161
	var v5163 int32
	_ = v5163
	var v5168 int32
	_ = v5168
	var v5169 int32
	_ = v5169
	var v5176 int32
	_ = v5176
	var v5178 int32
	_ = v5178
	var v5179 int32
	_ = v5179
	var v5181 int32
	_ = v5181
	var v5188 int32
	_ = v5188
	var v5189 int32
	_ = v5189
	var v5193 int32
	_ = v5193
	var v5197 int64
	_ = v5197
	var v5206 int32
	_ = v5206
	var v5231 int32
	_ = v5231
	var v5232 int32
	_ = v5232
	var v5267 int32
	_ = v5267
	var v5270 int32
	_ = v5270
	var v5274 int32
	_ = v5274
	var v5276 int32
	_ = v5276
	var v5278 int32
	_ = v5278
	var v5280 int32
	_ = v5280
	var v5284 int64
	_ = v5284
	var v5290 int32
	_ = v5290
	var v5293 int32
	_ = v5293
	var v5294 int32
	_ = v5294
	var v5297 int32
	_ = v5297
	var v5299 int32
	_ = v5299
	var v5304 int32
	_ = v5304
	var v5305 int32
	_ = v5305
	var v5312 int32
	_ = v5312
	var v5314 int32
	_ = v5314
	var v5315 int32
	_ = v5315
	var v5317 int32
	_ = v5317
	var v5324 int32
	_ = v5324
	var v5325 int32
	_ = v5325
	var v5329 int32
	_ = v5329
	var v5333 int64
	_ = v5333
	var v5344 int32
	_ = v5344
	var v5346 int32
	_ = v5346
	var v5358 int32
	_ = v5358
	var v5404 int32
	_ = v5404
	var v5405 int32
	_ = v5405
	var v5407 int32
	_ = v5407
	var v5409 int32
	_ = v5409
	var v5415 int32
	_ = v5415
	var v5418 int32
	_ = v5418
	var v5425 int32
	_ = v5425
	var v5427 int32
	_ = v5427
	var v5430 int32
	_ = v5430
	var v5433 int32
	_ = v5433
	var v5439 int32
	_ = v5439
	var v5440 int32
	_ = v5440
	var v5442 int32
	_ = v5442
	var v5443 int32
	_ = v5443
	var v5446 int32
	_ = v5446
	var v5447 int32
	_ = v5447
	var v5450 int32
	_ = v5450
	var v5451 int32
	_ = v5451
	var v5458 int32
	_ = v5458
	var v5460 int32
	_ = v5460
	var v5461 int64
	_ = v5461
	var v5465 int64
	_ = v5465
	var v5472 int32
	_ = v5472
	var v5475 int32
	_ = v5475
	var v5476 int32
	_ = v5476
	var v5479 int32
	_ = v5479
	var v5481 int32
	_ = v5481
	var v5486 int32
	_ = v5486
	var v5487 int32
	_ = v5487
	var v5494 int32
	_ = v5494
	var v5496 int32
	_ = v5496
	var v5497 int32
	_ = v5497
	var v5499 int32
	_ = v5499
	var v5506 int32
	_ = v5506
	var v5507 int32
	_ = v5507
	var v5511 int32
	_ = v5511
	var v5515 int64
	_ = v5515
	var v5580 int32
	_ = v5580
	var v5641 int32
	_ = v5641
	var v5643 int32
	_ = v5643
	var v5657 int32
	_ = v5657
	var v5706 int32
	_ = v5706
	var v5709 int32
	_ = v5709
	var v5724 int32
	_ = v5724
	v60 = m.G0
	v62 = v60 - int32(64)
	m.G0 = v62
	v64 = l1 * l0
	v65 = base.I64_extend_i32_s(v64)
	v66 = int32(2)
	if v65 == int64(0) {
		goto L5
	} else {
		goto L6
	}
L1:
	;
	F_free(m, v87)
	mBase = m.M
	goto L499
L2:
	;
	v90 = int64(1)
	v91 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_0)
	if l3 < int32(1) {
		goto L10
	} else {
		goto L11
	}
L3:
	;
	if v87 != 0 {
		goto L2
	} else {
		goto L9
	}
L4:
	;
	goto L3
L5:
	;
	v85 = F_malloc(m, base.I32_wrap_i64(v65)*v66)
	mBase = m.M
	v87 = v85
	goto L4
L6:
	;
	v73 = base.I64_div_u_s(int64(2147418112), v65)
	v74 = int32(0)
	v75 = base.I64_extend_i32_u(v66)
	if base.Ui64(int64(4294967295)) < base.Ui64(v75*v65) {
		v87 = v74
		goto L4
	} else {
		goto L7
	}
L7:
	;
	if base.Ui64(v73) < base.Ui64(v75) {
		v87 = v74
		goto L4
	} else {
		goto L8
	}
L8:
	;
	goto L5
L9:
	;
	v5724 = int32(0)
	goto L1
L10:
	;
	v98 = v91
	goto L12
L11:
	;
	v98 = int32(4)<<(uint(l3)%32) + v91
	goto L12
L12:
	;
	goto L16
L13:
	;
	v121 = int64(1)
	v122 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_1)
	goto L22
L14:
	;
	goto L13
L15:
	;
	v116 = F_calloc(m, base.I32_wrap_i64(v90), v98)
	mBase = m.M
	v118 = v116
	goto L14
L16:
	;
	v105 = base.I64_div_u_s(int64(2147418112), v90)
	v106 = int32(0)
	v107 = base.I64_extend_i32_u(v98)
	if base.Ui64(int64(4294967295)) < base.Ui64(v107*v90) {
		v118 = v106
		goto L14
	} else {
		goto L17
	}
L17:
	;
	if base.Ui64(v105) < base.Ui64(v107) {
		v118 = v106
		goto L14
	} else {
		goto L18
	}
L18:
	;
	goto L15
L19:
	;
	v144 = int32(0)
	if v118 == v144 {
		v4666 = v144
		goto L25
	} else {
		goto L26
	}
L20:
	;
	goto L19
L21:
	;
	v140 = F_calloc(m, base.I32_wrap_i64(v121), v122)
	mBase = m.M
	v142 = v140
	goto L20
L22:
	;
	v129 = base.I64_div_u_s(int64(2147418112), v121)
	v130 = int32(0)
	v131 = base.I64_extend_i32_u(v122)
	if base.Ui64(int64(4294967295)) < base.Ui64(v131*v121) {
		v142 = v130
		goto L20
	} else {
		goto L23
	}
L23:
	;
	if base.Ui64(v129) < base.Ui64(v131) {
		v142 = v130
		goto L20
	} else {
		goto L24
	}
L24:
	;
	goto L21
L25:
	;
	if v142 == int32(0) {
		goto L387
	} else {
		goto L388
	}
L26:
	;
	v147 = int32(0)
	if v142 == v147 {
		v4666 = v147
		goto L25
	} else {
		goto L27
	}
L27:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v118)+3232)) = v118 + int32(3236)
	if l3 < int32(1) {
		goto L28
	} else {
		goto L29
	}
L28:
	;
	v178 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_2)
	if int32(0) < l3 {
		goto L39
	} else {
		goto L40
	}
L29:
	;
	v157 = v62 + int32(52)
	v163 = F_WebPSafeCalloc(m, base.I64_extend_i32_s(int32(1)<<(uint(l3)%32)), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v157))) = v163
	if v163 != 0 {
		goto L31
	} else {
		goto L32
	}
L30:
	;
	if v171 == int32(0) {
		v4666 = int32(0)
		goto L25
	} else {
		goto L33
	}
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v157)+8)) = l3
	*(*int32)(unsafe.Add(mBase, uint32(v157)+4)) = int32(32) - l3
	v171 = int32(1)
	goto L30
L32:
	;
	v171 = int32(0)
	goto L30
L33:
	;
	goto L28
L34:
	;
	if l3 < int32(1) {
		v4666 = v4597
		goto L25
	} else {
		goto L382
	}
L35:
	;
	if v142 == int32(0) {
		goto L362
	} else {
		goto L363
	}
L36:
	;
	v4441 = int32(0)
	F_WebPSafeFree(m, v4441)
	mBase = m.M
	goto L360
L37:
	;
	if v186 == int32(0) {
		goto L36
	} else {
		goto L43
	}
L38:
	;
	goto L37
L39:
	;
	v183 = int32(4)<<(uint(l3)%32) + v178
	goto L41
L40:
	;
	v183 = v178
	goto L41
L41:
	;
	v186 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v183), int32(1))
	mBase = m.M
	if v186 == int32(0) {
		goto L38
	} else {
		goto L42
	}
L42:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v186)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v186)+3236)) = l3
	*(*int32)(unsafe.Add(mBase, uint32(v186)+3304)) = int32(16843009)
	v194 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v186)+3256)) = v194
	*(*int32)(unsafe.Add(mBase, uint32(v186))) = v186 + int32(3312)
	v201 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v186+int32(3248)))) = uint16(v201)
	v205 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v186+int32(3308)))) = uint8(v205)
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3264)))) = v194
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3272)))) = v194
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3280)))) = v194
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3288)))) = v194
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3296)))) = v194
	goto L38
L43:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v186)+3236)) = l3
	goto L46
L44:
	;
	v283 = m.G60
	v285 = m.G0
	v287 = v285 - int32(16)
	m.G0 = v287
	F_VP8LRefsCursorInit(m, v287+int32(4), l5)
	mBase = m.M
	v292 = *(*int32)(unsafe.Add(mBase, uint32(v287)+4))
	if v292 == int32(0) {
		goto L51
	} else {
		goto L52
	}
L45:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v186)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v186)+3304)) = int32(16843009)
	v253 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v186)+3256)) = v253
	v257 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v186+int32(3248)))) = uint16(v257)
	v261 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v186+int32(3308)))) = uint8(v261)
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3264)))) = v253
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3272)))) = v253
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3280)))) = v253
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3288)))) = v253
	*(*int64)(unsafe.Add(mBase, uint32(v186+int32(3296)))) = v253
	goto L44
L46:
	;
	v234 = *(*int32)(unsafe.Add(mBase, uint32(v186)))
	v235 = int32(0)
	v238 = int32(_a_F_VP8LBackwardReferencesTraceBackwards_2)
	if v235 < l3 {
		goto L47
	} else {
		goto L48
	}
L47:
	;
	v243 = int32(4)<<(uint(l3)%32) + v238
	goto L49
L48:
	;
	v243 = v238
	goto L49
L49:
	;
	v244 = F_memset(m, v186, v235, v243)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v244))) = v234
	*(*int32)(unsafe.Add(mBase, uint32(v244)+3236)) = l3
	goto L45
L50:
	;
	v320 = *(*int32)(unsafe.Add(mBase, uint32(v118)+3232))
	v322 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3236))
	v324 = int32(280)
	if int32(0) < v322 {
		goto L61
	} else {
		goto L62
	}
L51:
	;
	m.G0 = v287 + int32(16)
	goto L50
L52:
	;
	v295 = v292
	goto L53
L53:
	;
	F_HistogramAddSinglePixOrCopy(m, v186, v295, v283, l0)
	mBase = m.M
	v301 = *(*int32)(unsafe.Add(mBase, uint32(v287)+4))
	v303 = v301 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v287)+4)) = v303
	v305 = *(*int32)(unsafe.Add(mBase, uint32(v287)+12))
	if v303 != v305 {
		v311 = v303
		goto L55
	} else {
		goto L56
	}
L54:
	;
	goto L51
L55:
	;
	if v311 != 0 {
		v295 = v311
		goto L53
	} else {
		goto L57
	}
L56:
	;
	F_VP8LRefsCursorNextBlock(m, v287+int32(4))
	mBase = m.M
	v310 = *(*int32)(unsafe.Add(mBase, uint32(v287)+4))
	v311 = v310
	goto L55
L57:
	;
	goto L54
L58:
	;
	v892 = int32(0)
	v912 = v892
	v915 = v892
	v916 = v892
	goto L97
L59:
	;
	if base.Ui32(int32(255)) < base.Ui32(v518) {
		goto L89
	} else {
		goto L90
	}
L60:
	;
	v618 = int32(0)
	v620 = v329 << (uint(int32(2)) % 32)
	if base.Ui32(v620) < base.Ui32(int32(33)) {
		goto L75
	} else {
		goto L76
	}
L61:
	;
	v329 = int32(1)<<(uint(v322)%32) + v324
	goto L63
L62:
	;
	v329 = v324
	goto L63
L63:
	;
	if v329 < int32(1) {
		goto L60
	} else {
		goto L64
	}
L64:
	;
	v335 = v329 & int32(3)
	v336 = *(*int32)(unsafe.Add(mBase, uint32(v186)))
	v337 = int32(0)
	v340 = v336
	v356 = v329 & int32(2147483644)
	v359 = v337
	v360 = v337
	goto L65
L65:
	;
	v400 = *(*int32)(unsafe.Add(mBase, uint32(v340+int32(12))))
	v403 = *(*int32)(unsafe.Add(mBase, uint32(v340+int32(8))))
	v406 = *(*int32)(unsafe.Add(mBase, uint32(v340+int32(4))))
	v407 = *(*int32)(unsafe.Add(mBase, uint32(v340)))
	v411 = v400 + (v403 + (v406 + (v407 + v359)))
	v412 = int32(0)
	v423 = v360 + base.B2i32(v407 != v412) + base.B2i32(v406 != v412) + base.B2i32(v403 != v412) + base.B2i32(v400 != v412)
	v425 = v340 + int32(16)
	v427 = v356 + int32(-4)
	if v427 != 0 {
		v340 = v425
		v356 = v427
		v359 = v411
		v360 = v423
		goto L65
	} else {
		goto L67
	}
L66:
	;
	if v335 == int32(0) {
		v518 = v411
		v519 = v423
		goto L68
	} else {
		goto L69
	}
L67:
	;
	goto L66
L68:
	;
	if base.Ui32(int32(1)) < base.Ui32(v519) {
		goto L59
	} else {
		goto L73
	}
L69:
	;
	v431 = v425
	v448 = v335
	v450 = v411
	v451 = v423
	goto L70
L70:
	;
	v489 = *(*int32)(unsafe.Add(mBase, uint32(v431)))
	v490 = v489 + v450
	v493 = v451 + base.B2i32(v489 != int32(0))
	v497 = v448 + int32(-1)
	if v497 != 0 {
		v431 = v431 + int32(4)
		v448 = v497
		v450 = v490
		v451 = v493
		goto L70
	} else {
		goto L72
	}
L71:
	;
	v518 = v490
	v519 = v493
	goto L68
L72:
	;
	goto L71
L73:
	;
	goto L60
L74:
	;
	goto L58
L75:
	;
	if v620 == int32(0) {
		goto L77
	} else {
		goto L78
	}
L76:
	;
	base.MemoryFill(m, v320, v618, v620)
	goto L74
L77:
	;
	goto L74
L78:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v320))) = uint8(v618)
	v631 = v320 + v620
	*(*uint8)(unsafe.Add(mBase, uint32(v631+int32(-1)))) = uint8(v618)
	if base.Ui32(v620) < base.Ui32(int32(3)) {
		goto L77
	} else {
		goto L79
	}
L79:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v320)+2)) = uint8(v618)
	*(*uint8)(unsafe.Add(mBase, uint32(v320)+1)) = uint8(v618)
	*(*uint8)(unsafe.Add(mBase, uint32(v631+int32(-3)))) = uint8(v618)
	*(*uint8)(unsafe.Add(mBase, uint32(v631+int32(-2)))) = uint8(v618)
	if base.Ui32(v620) < base.Ui32(int32(7)) {
		goto L77
	} else {
		goto L80
	}
L80:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v320)+3)) = uint8(v618)
	*(*uint8)(unsafe.Add(mBase, uint32(v631+int32(-4)))) = uint8(v618)
	if base.Ui32(v620) < base.Ui32(int32(9)) {
		goto L77
	} else {
		goto L81
	}
L81:
	;
	v653 = int32(0)
	v656 = (v653 - v320) & int32(3)
	v657 = v320 + v656
	*(*int32)(unsafe.Add(mBase, uint32(v657))) = v653
	v665 = (v620 - v656) & int32(60)
	v666 = v657 + v665
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-4)))) = v653
	if base.Ui32(v665) < base.Ui32(int32(9)) {
		goto L77
	} else {
		goto L82
	}
L82:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v657)+8)) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v657)+4)) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-8)))) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-12)))) = v653
	if base.Ui32(v665) < base.Ui32(int32(25)) {
		goto L77
	} else {
		goto L83
	}
L83:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v657)+24)) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v657)+20)) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v657)+16)) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v657)+12)) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-16)))) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-20)))) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-24)))) = v653
	*(*int32)(unsafe.Add(mBase, uint32(v666+int32(-28)))) = v653
	v701 = v657&int32(4) | int32(24)
	v702 = v665 - v701
	if base.Ui32(v702) < base.Ui32(int32(32)) {
		goto L77
	} else {
		goto L84
	}
L84:
	;
	v707 = base.I64_extend_i32_u(v653) * int64(4294967297)
	v710 = v702
	v711 = v657 + v701
	goto L85
L85:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v711)+24)) = v707
	*(*int64)(unsafe.Add(mBase, uint32(v711)+16)) = v707
	*(*int64)(unsafe.Add(mBase, uint32(v711)+8)) = v707
	*(*int64)(unsafe.Add(mBase, uint32(v711))) = v707
	v723 = v710 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v723) {
		v710 = v723
		v711 = v711 + int32(32)
		goto L85
	} else {
		goto L87
	}
L86:
	;
	goto L77
L87:
	;
	goto L86
L88:
	;
	v767 = v320
	v768 = v329
	v771 = v336
	goto L91
L89:
	;
	v748 = m.G62
	v749 = *(*int32)(unsafe.Add(mBase, uint32(v748)))
	v750 = m.T0[v749].(func(*base.Module, int32) int32)(m, v518)
	mBase = m.M
	v751 = v750
	goto L88
L90:
	;
	v743 = m.G61
	v747 = *(*int32)(unsafe.Add(mBase, uint32(v743+v518<<(uint(int32(2))%32))))
	v751 = v747
	goto L88
L91:
	;
	v811 = *(*int32)(unsafe.Add(mBase, uint32(v771)))
	if base.Ui32(int32(255)) < base.Ui32(v811) {
		goto L94
	} else {
		goto L95
	}
L92:
	;
	goto L58
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v767))) = v751 - v822
	v825 = int32(4)
	v830 = v768 + int32(-1)
	if v830 != 0 {
		v767 = v767 + v825
		v768 = v830
		v771 = v771 + v825
		goto L91
	} else {
		goto L96
	}
L94:
	;
	v819 = m.G62
	v820 = *(*int32)(unsafe.Add(mBase, uint32(v819)))
	v821 = m.T0[v820].(func(*base.Module, int32) int32)(m, v811)
	mBase = m.M
	v822 = v821
	goto L93
L95:
	;
	v814 = m.G61
	v818 = *(*int32)(unsafe.Add(mBase, uint32(v814+v811<<(uint(int32(2))%32))))
	v822 = v818
	goto L93
L96:
	;
	goto L92
L97:
	;
	v954 = v186 + v915
	v955 = int32(16)
	v957 = *(*int32)(unsafe.Add(mBase, uint32(v954+v955)))
	v960 = *(*int32)(unsafe.Add(mBase, uint32(v954+int32(12))))
	v963 = *(*int32)(unsafe.Add(mBase, uint32(v954+int32(8))))
	v966 = *(*int32)(unsafe.Add(mBase, uint32(v954+int32(4))))
	v970 = v957 + (v960 + (v963 + (v966 + v916)))
	v971 = int32(0)
	v982 = v912 + base.B2i32(v966 != v971) + base.B2i32(v963 != v971) + base.B2i32(v960 != v971) + base.B2i32(v957 != v971)
	v984 = v915 + v955
	if v984 != int32(1024) {
		v912 = v982
		v915 = v984
		v916 = v970
		goto L97
	} else {
		goto L99
	}
L98:
	;
	v988 = v118 + int32(1024)
	if base.Ui32(int32(1)) < base.Ui32(v982) {
		goto L101
	} else {
		goto L102
	}
L99:
	;
	goto L98
L100:
	;
	v1265 = int32(0)
	v1285 = v1265
	v1288 = v1265
	v1289 = v1265
	goto L126
L101:
	;
	if base.Ui32(int32(255)) < base.Ui32(v970) {
		goto L118
	} else {
		goto L119
	}
L102:
	;
	goto L105
L103:
	;
	goto L100
L105:
	;
	base.MemoryFill(m, v988, int32(0), int32(1024))
	goto L103
L117:
	;
	v1126 = int32(0)
	goto L120
L118:
	;
	v1120 = m.G62
	v1121 = *(*int32)(unsafe.Add(mBase, uint32(v1120)))
	v1122 = m.T0[v1121].(func(*base.Module, int32) int32)(m, v970)
	mBase = m.M
	v1123 = v1122
	goto L117
L119:
	;
	v1115 = m.G61
	v1119 = *(*int32)(unsafe.Add(mBase, uint32(v1115+v970<<(uint(int32(2))%32))))
	v1123 = v1119
	goto L117
L120:
	;
	v1185 = *(*int32)(unsafe.Add(mBase, uint32(v186+int32(4)+v1126)))
	if base.Ui32(int32(255)) < base.Ui32(v1185) {
		goto L123
	} else {
		goto L124
	}
L121:
	;
	goto L100
L122:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v988+v1126))) = v1123 - v1196
	v1201 = v1126 + int32(4)
	if v1201 != int32(1024) {
		v1126 = v1201
		goto L120
	} else {
		goto L125
	}
L123:
	;
	v1193 = m.G62
	v1194 = *(*int32)(unsafe.Add(mBase, uint32(v1193)))
	v1195 = m.T0[v1194].(func(*base.Module, int32) int32)(m, v1185)
	mBase = m.M
	v1196 = v1195
	goto L122
L124:
	;
	v1188 = m.G61
	v1192 = *(*int32)(unsafe.Add(mBase, uint32(v1188+v1185<<(uint(int32(2))%32))))
	v1196 = v1192
	goto L122
L125:
	;
	goto L121
L126:
	;
	v1327 = v186 + v1288
	v1330 = *(*int32)(unsafe.Add(mBase, uint32(v1327+int32(1040))))
	v1333 = *(*int32)(unsafe.Add(mBase, uint32(v1327+int32(1036))))
	v1336 = *(*int32)(unsafe.Add(mBase, uint32(v1327+int32(1032))))
	v1339 = *(*int32)(unsafe.Add(mBase, uint32(v1327+int32(1028))))
	v1343 = v1330 + (v1333 + (v1336 + (v1339 + v1289)))
	v1344 = int32(0)
	v1355 = v1285 + base.B2i32(v1339 != v1344) + base.B2i32(v1336 != v1344) + base.B2i32(v1333 != v1344) + base.B2i32(v1330 != v1344)
	v1357 = v1288 + int32(16)
	if v1357 != int32(1024) {
		v1285 = v1355
		v1288 = v1357
		v1289 = v1343
		goto L126
	} else {
		goto L128
	}
L127:
	;
	v1361 = v118 + int32(2048)
	if base.Ui32(int32(1)) < base.Ui32(v1355) {
		goto L130
	} else {
		goto L131
	}
L128:
	;
	goto L127
L129:
	;
	v1638 = int32(0)
	v1658 = v1638
	v1661 = v1638
	v1662 = v1638
	goto L155
L130:
	;
	if base.Ui32(int32(255)) < base.Ui32(v1343) {
		goto L147
	} else {
		goto L148
	}
L131:
	;
	goto L134
L132:
	;
	goto L129
L134:
	;
	base.MemoryFill(m, v1361, int32(0), int32(1024))
	goto L132
L146:
	;
	v1499 = int32(0)
	goto L149
L147:
	;
	v1493 = m.G62
	v1494 = *(*int32)(unsafe.Add(mBase, uint32(v1493)))
	v1495 = m.T0[v1494].(func(*base.Module, int32) int32)(m, v1343)
	mBase = m.M
	v1496 = v1495
	goto L146
L148:
	;
	v1488 = m.G61
	v1492 = *(*int32)(unsafe.Add(mBase, uint32(v1488+v1343<<(uint(int32(2))%32))))
	v1496 = v1492
	goto L146
L149:
	;
	v1558 = *(*int32)(unsafe.Add(mBase, uint32(v186+int32(1028)+v1499)))
	if base.Ui32(int32(255)) < base.Ui32(v1558) {
		goto L152
	} else {
		goto L153
	}
L150:
	;
	goto L129
L151:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1361+v1499))) = v1496 - v1569
	v1574 = v1499 + int32(4)
	if v1574 != int32(1024) {
		v1499 = v1574
		goto L149
	} else {
		goto L154
	}
L152:
	;
	v1566 = m.G62
	v1567 = *(*int32)(unsafe.Add(mBase, uint32(v1566)))
	v1568 = m.T0[v1567].(func(*base.Module, int32) int32)(m, v1558)
	mBase = m.M
	v1569 = v1568
	goto L151
L153:
	;
	v1561 = m.G61
	v1565 = *(*int32)(unsafe.Add(mBase, uint32(v1561+v1558<<(uint(int32(2))%32))))
	v1569 = v1565
	goto L151
L154:
	;
	goto L150
L155:
	;
	v1700 = v186 + v1661
	v1703 = *(*int32)(unsafe.Add(mBase, uint32(v1700+int32(2064))))
	v1706 = *(*int32)(unsafe.Add(mBase, uint32(v1700+int32(2060))))
	v1709 = *(*int32)(unsafe.Add(mBase, uint32(v1700+int32(2056))))
	v1712 = *(*int32)(unsafe.Add(mBase, uint32(v1700+int32(2052))))
	v1716 = v1703 + (v1706 + (v1709 + (v1712 + v1662)))
	v1717 = int32(0)
	v1728 = v1658 + base.B2i32(v1712 != v1717) + base.B2i32(v1709 != v1717) + base.B2i32(v1706 != v1717) + base.B2i32(v1703 != v1717)
	v1730 = v1661 + int32(16)
	if v1730 != int32(1024) {
		v1658 = v1728
		v1661 = v1730
		v1662 = v1716
		goto L155
	} else {
		goto L157
	}
L156:
	;
	if base.Ui32(int32(1)) < base.Ui32(v1728) {
		goto L159
	} else {
		goto L160
	}
L157:
	;
	goto L156
L158:
	;
	v2008 = v118 + int32(3072)
	v2009 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3080))
	v2010 = int32(0)
	v2012 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3076))
	v2016 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3084))
	v2020 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3088))
	v2024 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3092))
	v2028 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3096))
	v2032 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3100))
	v2036 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3104))
	v2040 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3108))
	v2044 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3112))
	v2048 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3116))
	v2052 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3120))
	v2056 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3124))
	v2060 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3128))
	v2064 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3132))
	v2068 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3136))
	v2072 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3140))
	v2076 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3144))
	v2080 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3148))
	v2084 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3152))
	v2088 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3156))
	v2092 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3160))
	v2096 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3164))
	v2100 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3168))
	v2104 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3172))
	v2108 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3176))
	v2112 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3180))
	v2116 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3184))
	v2120 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3188))
	v2124 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3192))
	v2128 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3196))
	v2132 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3200))
	v2136 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3204))
	v2140 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3208))
	v2144 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3212))
	v2148 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3216))
	v2152 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3220))
	v2156 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3224))
	v2160 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3228))
	v2164 = *(*int32)(unsafe.Add(mBase, uint32(v186)+3232))
	if base.Ui32(int32(1)) < base.Ui32(base.B2i32(v2009 != v2010)+base.B2i32(v2012 != v2010)+base.B2i32(v2016 != v2010)+base.B2i32(v2020 != v2010)+base.B2i32(v2024 != v2010)+base.B2i32(v2028 != v2010)+base.B2i32(v2032 != v2010)+base.B2i32(v2036 != v2010)+base.B2i32(v2040 != v2010)+base.B2i32(v2044 != v2010)+base.B2i32(v2048 != v2010)+base.B2i32(v2052 != v2010)+base.B2i32(v2056 != v2010)+base.B2i32(v2060 != v2010)+base.B2i32(v2064 != v2010)+base.B2i32(v2068 != v2010)+base.B2i32(v2072 != v2010)+base.B2i32(v2076 != v2010)+base.B2i32(v2080 != v2010)+base.B2i32(v2084 != v2010)+base.B2i32(v2088 != v2010)+base.B2i32(v2092 != v2010)+base.B2i32(v2096 != v2010)+base.B2i32(v2100 != v2010)+base.B2i32(v2104 != v2010)+base.B2i32(v2108 != v2010)+base.B2i32(v2112 != v2010)+base.B2i32(v2116 != v2010)+base.B2i32(v2120 != v2010)+base.B2i32(v2124 != v2010)+base.B2i32(v2128 != v2010)+base.B2i32(v2132 != v2010)+base.B2i32(v2136 != v2010)+base.B2i32(v2140 != v2010)+base.B2i32(v2144 != v2010)+base.B2i32(v2148 != v2010)+base.B2i32(v2152 != v2010)+base.B2i32(v2156 != v2010)+base.B2i32(v2160 != v2010)+base.B2i32(v2164 != v2010)) {
		goto L185
	} else {
		goto L186
	}
L159:
	;
	if base.Ui32(int32(255)) < base.Ui32(v1716) {
		goto L176
	} else {
		goto L177
	}
L160:
	;
	goto L163
L161:
	;
	goto L158
L163:
	;
	base.MemoryFill(m, v118, int32(0), int32(1024))
	goto L161
L175:
	;
	v1870 = int32(0)
	goto L178
L176:
	;
	v1864 = m.G62
	v1865 = *(*int32)(unsafe.Add(mBase, uint32(v1864)))
	v1866 = m.T0[v1865].(func(*base.Module, int32) int32)(m, v1716)
	mBase = m.M
	v1867 = v1866
	goto L175
L177:
	;
	v1859 = m.G61
	v1863 = *(*int32)(unsafe.Add(mBase, uint32(v1859+v1716<<(uint(int32(2))%32))))
	v1867 = v1863
	goto L175
L178:
	;
	v1929 = *(*int32)(unsafe.Add(mBase, uint32(v186+int32(2052)+v1870)))
	if base.Ui32(int32(255)) < base.Ui32(v1929) {
		goto L181
	} else {
		goto L182
	}
L179:
	;
	goto L158
L180:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v118+v1870))) = v1867 - v1940
	v1945 = v1870 + int32(4)
	if v1945 != int32(1024) {
		v1870 = v1945
		goto L178
	} else {
		goto L183
	}
L181:
	;
	v1937 = m.G62
	v1938 = *(*int32)(unsafe.Add(mBase, uint32(v1937)))
	v1939 = m.T0[v1938].(func(*base.Module, int32) int32)(m, v1929)
	mBase = m.M
	v1940 = v1939
	goto L180
L182:
	;
	v1932 = m.G61
	v1936 = *(*int32)(unsafe.Add(mBase, uint32(v1932+v1929<<(uint(int32(2))%32))))
	v1940 = v1936
	goto L180
L183:
	;
	goto L179
L184:
	;
	F_WebPSafeFree(m, v186)
	mBase = m.M
	goto L210
L185:
	;
	v2330 = v2009 + v2012 + v2016 + v2020 + v2024 + v2028 + v2032 + v2036 + v2040 + v2044 + v2048 + v2052 + v2056 + v2060 + v2064 + v2068 + v2072 + v2076 + v2080 + v2084 + v2088 + v2092 + v2096 + v2100 + v2104 + v2108 + v2112 + v2116 + v2120 + v2124 + v2128 + v2132 + v2136 + v2140 + v2144 + v2148 + v2152 + v2156 + v2160 + v2164
	if base.Ui32(int32(255)) < base.Ui32(v2330) {
		goto L202
	} else {
		goto L203
	}
L186:
	;
	goto L189
L187:
	;
	goto L184
L189:
	;
	base.MemoryFill(m, v2008, int32(0), int32(160))
	goto L187
L201:
	;
	v2346 = int32(0)
	goto L204
L202:
	;
	v2338 = m.G62
	v2339 = *(*int32)(unsafe.Add(mBase, uint32(v2338)))
	v2340 = m.T0[v2339].(func(*base.Module, int32) int32)(m, v2330)
	mBase = m.M
	v2341 = v2340
	goto L201
L203:
	;
	v2333 = m.G61
	v2337 = *(*int32)(unsafe.Add(mBase, uint32(v2333+v2330<<(uint(int32(2))%32))))
	v2341 = v2337
	goto L201
L204:
	;
	v2405 = *(*int32)(unsafe.Add(mBase, uint32(v186+int32(3076)+v2346)))
	if base.Ui32(int32(255)) < base.Ui32(v2405) {
		goto L207
	} else {
		goto L208
	}
L205:
	;
	goto L184
L206:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2008+v2346))) = v2341 - v2416
	v2421 = v2346 + int32(4)
	if v2421 != int32(160) {
		v2346 = v2421
		goto L204
	} else {
		goto L209
	}
L207:
	;
	v2413 = m.G62
	v2414 = *(*int32)(unsafe.Add(mBase, uint32(v2413)))
	v2415 = m.T0[v2414].(func(*base.Module, int32) int32)(m, v2405)
	mBase = m.M
	v2416 = v2415
	goto L206
L208:
	;
	v2408 = m.G61
	v2412 = *(*int32)(unsafe.Add(mBase, uint32(v2408+v2405<<(uint(int32(2))%32))))
	v2416 = v2412
	goto L206
L209:
	;
	goto L205
L210:
	;
	v2484 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+8)) = v2484
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0]))) = v2484
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1]))) = v2484
	v2492 = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v2492
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[3]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[4]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_5)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[5]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_6)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[6]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_7)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[7]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_8)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[8]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_9)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[9]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_10)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[10]))) = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_11)
	v2519 = v142 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[11]))) = v2519
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[12]))) = v2484
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13]))) = v87
	*(*int64)(unsafe.Add(mBase, uint32(v142))) = int64(0)
	v2526 = int32(4095)
	if v64 < v2526 {
		goto L211
	} else {
		goto L212
	}
L211:
	;
	v2529 = v64
	goto L213
L212:
	;
	v2529 = v2526
	goto L213
L213:
	;
	if int32(0) < v64 {
		goto L215
	} else {
		goto L216
	}
L214:
	;
	v2870 = int32(16)
	if v2868 == int64(0) {
		goto L240
	} else {
		goto L241
	}
L215:
	;
	v2535 = *(*int32)(unsafe.Add(mBase, uint32(v118)+3232))
	v2539 = m.G54
	v2541 = v142 + int32(16)
	v2543 = int32(0)
	v2556 = v2539
	v2562 = v2541
	goto L217
L216:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v142)+12)) = int32(1)
	v2825 = v2484
	v2868 = int64(1)
	goto L214
L217:
	;
	if base.Ui32(int32(511)) < base.Ui32(v2543) {
		goto L220
	} else {
		goto L221
	}
L218:
	;
	v2638 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+12)) = v2638
	if v64 != v2638 {
		goto L223
	} else {
		goto L224
	}
L219:
	;
	v2625 = int32(2)
	v2628 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2535+int32(1024)+v2620<<(uint(v2625)%32)))))
	*(*int64)(unsafe.Add(mBase, uint32(v2562))) = base.I64_extend_i32_s(v2621)<<(uint(int64(23))%64) + v2628
	v2636 = v2543 + int32(1)
	if v2529 != v2636 {
		v2543 = v2636
		v2556 = v2556 + v2625
		v2562 = v2562 + int32(8)
		goto L217
	} else {
		goto L222
	}
L220:
	;
	v2607 = int32(-1)
	v2608 = v2543 + v2607
	v2611 = base.I32_clz(v2608) ^ int32(31)
	v2613 = v2611 + v2607
	v2615 = int32(1)
	v2620 = int32(base.Ui32(v2608)>>(uint(v2613)%32))&v2615 | v2611<<(uint(v2615)%32)
	v2621 = v2613
	goto L219
L221:
	;
	v2605 = int32(*(*int8)(unsafe.Add(mBase, uint32(v2556+int32(1)))))
	v2606 = int32(*(*int8)(unsafe.Add(mBase, uint32(v2556))))
	v2620 = v2606
	v2621 = v2605
	goto L219
L222:
	;
	goto L218
L223:
	;
	v2646 = v2529 + int32(-1)
	v2649 = *(*int64)(unsafe.Add(mBase, uint32(v2541)))
	if v64 != int32(2) {
		goto L226
	} else {
		goto L227
	}
L224:
	;
	v2825 = int32(0)
	v2868 = int64(1)
	goto L214
L225:
	;
	if v2646&int32(1) == int32(0) {
		v2808 = v2753
		goto L235
	} else {
		goto L236
	}
L226:
	;
	v2660 = v142 + int32(32)
	v2673 = int32(1)
	v2679 = int32(0)
	v2716 = v2649
	goto L228
L227:
	;
	v2740 = int32(1)
	v2753 = v2638
	v2796 = v2649
	goto L225
L228:
	;
	v2720 = *(*int64)(unsafe.Add(mBase, uint32(v2660+int32(-8))))
	if v2720 == v2716 {
		v2725 = v2673
		goto L230
	} else {
		goto L231
	}
L229:
	;
	v2740 = v2679 + int32(3)
	v2753 = v2731
	v2796 = v2726
	goto L225
L230:
	;
	v2726 = *(*int64)(unsafe.Add(mBase, uint32(v2660)))
	if v2726 == v2720 {
		v2731 = v2725
		goto L232
	} else {
		goto L233
	}
L231:
	;
	v2723 = v2673 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+12)) = v2723
	v2725 = v2723
	goto L230
L232:
	;
	v2735 = v2679 + int32(2)
	if v2646&int32(-2) != v2735 {
		v2660 = v2660 + int32(16)
		v2673 = v2731
		v2679 = v2735
		v2716 = v2726
		goto L228
	} else {
		goto L234
	}
L233:
	;
	v2729 = v2725 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+12)) = v2729
	v2731 = v2729
	goto L232
L234:
	;
	goto L229
L235:
	;
	v2825 = int32(1)
	v2868 = base.I64_extend_i32_u(v2808)
	goto L214
L236:
	;
	v2803 = *(*int64)(unsafe.Add(mBase, uint32(v2541+v2740<<(uint(int32(3))%32))))
	if v2803 == v2796 {
		v2808 = v2753
		goto L235
	} else {
		goto L237
	}
L237:
	;
	v2806 = v2753 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v142)+12)) = v2806
	v2808 = v2806
	goto L235
L238:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v142)+8)) = v2891
	if v2891 == int32(0) {
		goto L35
	} else {
		goto L244
	}
L239:
	;
	goto L238
L240:
	;
	v2889 = F_malloc(m, base.I32_wrap_i64(v2868)*v2870)
	mBase = m.M
	v2891 = v2889
	goto L239
L241:
	;
	v2877 = base.I64_div_u_s(int64(2147418112), v2868)
	v2878 = int32(0)
	v2879 = base.I64_extend_i32_u(v2870)
	if base.Ui64(int64(4294967295)) < base.Ui64(v2879*v2868) {
		v2891 = v2878
		goto L239
	} else {
		goto L242
	}
L242:
	;
	if base.Ui64(v2877) < base.Ui64(v2879) {
		v2891 = v2878
		goto L239
	} else {
		goto L243
	}
L243:
	;
	goto L240
L244:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2891)+8)) = int64(4294967296)
	v2898 = *(*int64)(unsafe.Add(mBase, uint32(v142)+16))
	*(*int64)(unsafe.Add(mBase, uint32(v2891))) = v2898
	if v2825 == int32(0) {
		goto L245
	} else {
		goto L246
	}
L245:
	;
	v3137 = int32(8)
	if v65 == int64(0) {
		goto L261
	} else {
		goto L262
	}
L246:
	;
	v2902 = int32(1)
	v2904 = v2529 + int32(-1)
	if v64 == int32(2) {
		v3003 = v2891
		v3016 = v2902
		v3059 = v2898
		goto L247
	} else {
		goto L248
	}
L247:
	;
	if v2904&v2902 == int32(0) {
		goto L245
	} else {
		goto L256
	}
L248:
	;
	v2915 = v2891
	v2928 = v142 + int32(32)
	v2934 = int32(0)
	v2971 = v2898
	goto L249
L249:
	;
	v2975 = *(*int64)(unsafe.Add(mBase, uint32(v2928+int32(-8))))
	if v2975 == v2971 {
		v2983 = v2915
		goto L251
	} else {
		goto L252
	}
L250:
	;
	v3003 = v2993
	v3016 = v2934 + int32(3)
	v3059 = v2987
	goto L247
L251:
	;
	v2985 = v2934 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v2983)+12)) = v2985
	v2987 = *(*int64)(unsafe.Add(mBase, uint32(v2928)))
	if v2987 == v2975 {
		v2993 = v2983
		goto L253
	} else {
		goto L254
	}
L252:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2915)+16)) = v2975
	*(*int32)(unsafe.Add(mBase, uint32(v2915)+24)) = v2934 + int32(1)
	v2983 = v2915 + int32(16)
	goto L251
L253:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2993)+12)) = v2934 + int32(3)
	if v2904&int32(-2) != v2985 {
		v2915 = v2993
		v2928 = v2928 + int32(16)
		v2934 = v2985
		v2971 = v2987
		goto L249
	} else {
		goto L255
	}
L254:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v2983)+16)) = v2987
	*(*int32)(unsafe.Add(mBase, uint32(v2983)+24)) = v2985
	v2993 = v2983 + int32(16)
	goto L253
L255:
	;
	goto L250
L256:
	;
	v3068 = *(*int64)(unsafe.Add(mBase, uint32(v142+int32(16)+v3016<<(uint(int32(3))%32))))
	if v3068 == v3059 {
		v3074 = v3003
		goto L257
	} else {
		goto L258
	}
L257:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3074)+12)) = v3016 + int32(1)
	goto L245
L258:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3003)+16)) = v3068
	*(*int32)(unsafe.Add(mBase, uint32(v3003)+24)) = v3016
	v3074 = v3003 + int32(16)
	goto L257
L259:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0]))) = v3158
	if v3158 == int32(0) {
		goto L35
	} else {
		goto L265
	}
L260:
	;
	goto L259
L261:
	;
	v3156 = F_malloc(m, base.I32_wrap_i64(v65)*v3137)
	mBase = m.M
	v3158 = v3156
	goto L260
L262:
	;
	v3144 = base.I64_div_u_s(int64(2147418112), v65)
	v3145 = int32(0)
	v3146 = base.I64_extend_i32_u(v3137)
	if base.Ui64(int64(4294967295)) < base.Ui64(v3146*v65) {
		v3158 = v3145
		goto L260
	} else {
		goto L263
	}
L263:
	;
	if base.Ui64(v3144) < base.Ui64(v3146) {
		v3158 = v3145
		goto L260
	} else {
		goto L264
	}
L264:
	;
	goto L261
L265:
	;
	if v64 < int32(1) {
		goto L266
	} else {
		goto L267
	}
L266:
	;
	v3455 = int32(0)
	*(*uint16)(unsafe.Add(mBase, uint32(v87))) = uint16(v3455)
	v3457 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
	if l3 < int32(1) {
		goto L278
	} else {
		goto L279
	}
L267:
	;
	v3166 = v64 & int32(7)
	if base.Ui32(v64) < base.Ui32(int32(8)) {
		v3287 = int32(0)
		goto L268
	} else {
		goto L269
	}
L268:
	;
	if v3166 == int32(0) {
		goto L266
	} else {
		goto L273
	}
L269:
	;
	v3174 = v3158
	v3193 = int32(0)
	goto L270
L270:
	;
	v3232 = int64(9223372036854775807)
	*(*int64)(unsafe.Add(mBase, uint32(v3174))) = v3232
	*(*int64)(unsafe.Add(mBase, uint32(v3174+int32(56)))) = v3232
	*(*int64)(unsafe.Add(mBase, uint32(v3174+int32(48)))) = v3232
	*(*int64)(unsafe.Add(mBase, uint32(v3174+int32(40)))) = v3232
	*(*int64)(unsafe.Add(mBase, uint32(v3174+int32(32)))) = v3232
	*(*int64)(unsafe.Add(mBase, uint32(v3174+int32(24)))) = v3232
	*(*int64)(unsafe.Add(mBase, uint32(v3174+int32(16)))) = v3232
	v3258 = int32(8)
	*(*int64)(unsafe.Add(mBase, uint32(v3174+v3258))) = v3232
	v3265 = v3193 + v3258
	if v64&int32(2147483640) != v3265 {
		v3174 = v3174 + int32(64)
		v3193 = v3265
		goto L270
	} else {
		goto L272
	}
L271:
	;
	v3287 = v3265
	goto L268
L272:
	;
	goto L271
L273:
	;
	v3332 = v3158 + v3287<<(uint(int32(3))%32)
	v3345 = v3166
	goto L274
L274:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3332))) = int64(9223372036854775807)
	v3395 = v3345 + int32(-1)
	if v3395 != 0 {
		v3332 = v3332 + int32(8)
		v3345 = v3395
		goto L274
	} else {
		goto L276
	}
L275:
	;
	goto L266
L276:
	;
	goto L275
L277:
	;
	v3517 = *(*int64)(unsafe.Add(mBase, uint32(v3158)))
	v3521 = base.I64_div_u_s(v3516+int64(50), int64(100))
	if v3517 <= v3521 {
		goto L283
	} else {
		goto L284
	}
L278:
	;
	v3486 = int32(1020)
	v3489 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v988+int32(base.Ui32(v3457)>>(uint(int32(14))%32))&v3486))))
	v3495 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v118+int32(base.Ui32(v3457)>>(uint(int32(22))%32))&v3486))))
	v3497 = *(*int32)(unsafe.Add(mBase, uint32(v118)+3232))
	v3503 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3497+int32(base.Ui32(v3457)>>(uint(int32(6))%32))&v3486))))
	v3510 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1361+v3457&int32(255)<<(uint(int32(2))%32)))))
	v3516 = (v3489 + v3495 + v3503 + v3510) * int64(82)
	goto L277
L279:
	;
	v3460 = *(*int32)(unsafe.Add(mBase, uint32(v62)+52))
	v3463 = *(*int32)(unsafe.Add(mBase, uint32(v62)+56))
	v3464 = int32(base.Ui32(v3457*int32(506832829)) >> (uint(v3463) % 32))
	v3467 = v3460 + v3464<<(uint(int32(2))%32)
	if v3464 < int32(0) {
		goto L280
	} else {
		goto L281
	}
L280:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3467))) = v3457
	goto L278
L281:
	;
	v3470 = *(*int32)(unsafe.Add(mBase, uint32(v3467)))
	if v3470 != v3457 {
		goto L280
	} else {
		goto L282
	}
L282:
	;
	v3472 = *(*int32)(unsafe.Add(mBase, uint32(v118)+3232))
	v3478 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3472+v3464<<(uint(int32(2))%32)+int32(1120)))))
	v3516 = v3478 * int64(68)
	goto L277
L283:
	;
	if v64 < int32(2) {
		goto L285
	} else {
		goto L286
	}
L284:
	;
	v3523 = int32(1)
	*(*uint16)(unsafe.Add(mBase, uint32(v87))) = uint16(v3523)
	*(*int64)(unsafe.Add(mBase, uint32(v3158))) = v3521
	goto L283
L285:
	;
	v4438 = *(*int32)(unsafe.Add(mBase, uint32(l5)+4))
	v4597 = base.B2i32(v4438 == int32(0))
	goto L34
L286:
	;
	v3529 = int32(0)
	v3533 = int32(-1)
	v3545 = int64(-1)
	v3552 = v3533
	v3556 = int32(1)
	v3561 = v3533
	v3563 = v3529
	v3564 = int32(4)
	v3566 = v3529
	v3568 = v3533
	goto L287
L287:
	;
	v3595 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v3597 = v3556 << (uint(int32(2)) % 32)
	v3599 = *(*int32)(unsafe.Add(mBase, uint32(v3595+v3597)))
	v3601 = *(*int32)(unsafe.Add(mBase, uint32(l2+v3597)))
	v3602 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	v3604 = v3556 + int32(-1)
	v3608 = *(*int64)(unsafe.Add(mBase, uint32(v3602+v3604<<(uint(int32(3))%32))))
	if l3 < int32(1) {
		goto L290
	} else {
		goto L291
	}
L288:
	;
	goto L285
L289:
	;
	v3669 = v3599 & int32(4095)
	v3671 = v3556 << (uint(int32(3)) % 32)
	v3672 = v3602 + v3671
	v3673 = *(*int64)(unsafe.Add(mBase, uint32(v3672)))
	v3677 = base.I64_div_u_s(v3667+int64(50), int64(100))
	v3678 = v3677 + v3608
	if v3673 <= v3678 {
		goto L295
	} else {
		goto L296
	}
L290:
	;
	v3636 = int32(1020)
	v3639 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v988+int32(base.Ui32(v3601)>>(uint(int32(14))%32))&v3636))))
	v3645 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v118+int32(base.Ui32(v3601)>>(uint(int32(22))%32))&v3636))))
	v3647 = *(*int32)(unsafe.Add(mBase, uint32(v118)+3232))
	v3653 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3647+int32(base.Ui32(v3601)>>(uint(int32(6))%32))&v3636))))
	v3660 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1361+v3601&int32(255)<<(uint(int32(2))%32)))))
	v3667 = (v3639 + v3645 + v3653 + v3660) * int64(82)
	goto L289
L291:
	;
	v3611 = *(*int32)(unsafe.Add(mBase, uint32(v62)+52))
	v3614 = *(*int32)(unsafe.Add(mBase, uint32(v62)+56))
	v3615 = int32(base.Ui32(v3601*int32(506832829)) >> (uint(v3614) % 32))
	v3617 = v3615 << (uint(int32(2)) % 32)
	v3618 = v3611 + v3617
	if v3615 < int32(0) {
		goto L292
	} else {
		goto L293
	}
L292:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3618))) = v3601
	goto L290
L293:
	;
	v3621 = *(*int32)(unsafe.Add(mBase, uint32(v3618)))
	if v3621 != v3601 {
		goto L292
	} else {
		goto L294
	}
L294:
	;
	v3623 = *(*int32)(unsafe.Add(mBase, uint32(v118)+3232))
	v3627 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3623+v3617+int32(1120)))))
	v3667 = v3627 * int64(68)
	goto L289
L295:
	;
	v3687 = int32(base.Ui32(v3599) >> (uint(int32(12)) % 32))
	if base.Ui32(v3669) < base.Ui32(int32(2)) {
		v4159 = v3545
		v4180 = v3566
		v4182 = v3568
		goto L297
	} else {
		goto L298
	}
L296:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3672))) = v3678
	v3681 = int32(1)
	*(*uint16)(unsafe.Add(mBase, uint32(v87+v3556<<(uint(v3681)%32)))) = uint16(v3681)
	goto L295
L297:
	;
	v4209 = *(*int32)(unsafe.Add(mBase, uint32(v142)))
	if v4209 == int32(0) {
		goto L340
	} else {
		goto L341
	}
L298:
	;
	if v3687 == v3552 {
		goto L299
	} else {
		goto L300
	}
L299:
	;
	if v3568 != 0 {
		goto L311
	} else {
		goto L312
	}
L300:
	;
	v3693 = base.I32_div_s(v3687, l0)
	v3695 = v3687 - v3693*l0
	if int32(7) < v3693 {
		goto L304
	} else {
		goto L305
	}
L301:
	;
	v3762 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2008+v3755<<(uint(int32(2))%32)))))
	v3763 = base.I64_extend_i32_s(v3754)<<(uint(int64(23))%64) + v3762
	F_PushInterval(m, v142, v3763+v3608, v3556, v3669)
	mBase = m.M
	v4159 = v3763
	v4180 = v3566
	v4182 = int32(1)
	goto L297
L302:
	;
	v3741 = int32(-1)
	v3742 = v3732 + v3741
	v3745 = base.I32_clz(v3742) ^ int32(31)
	v3747 = v3745 + v3741
	v3749 = int32(1)
	v3754 = v3747
	v3755 = int32(base.Ui32(v3742)>>(uint(v3747)%32))&v3749 | v3745<<(uint(v3749)%32)
	goto L301
L303:
	;
	if int32(511) < v3732 {
		goto L302
	} else {
		goto L310
	}
L304:
	;
	if int32(6) < v3693 {
		goto L307
	} else {
		goto L308
	}
L305:
	;
	if int32(8) < v3695 {
		goto L304
	} else {
		goto L306
	}
L306:
	;
	v3700 = m.G1
	v3709 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3700+int32(_a_F_VP8LBackwardReferencesTraceBackwards_13)+(v3693<<(uint(int32(4))%32)|int32(8)-v3695)))))
	v3732 = v3709 + int32(1)
	goto L303
L307:
	;
	v3732 = v3687 + int32(120)
	goto L303
L308:
	;
	if v3695 <= l0+int32(-8) {
		goto L307
	} else {
		goto L309
	}
L309:
	;
	v3721 = m.G1
	v3727 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v3693<<(uint(int32(4))%32)-v3695+(v3721+int32(_a_F_VP8LBackwardReferencesTraceBackwards_13))+int32(24)))))
	v3732 = v3727 + int32(1)
	goto L303
L310:
	;
	v3735 = m.G54
	v3738 = v3735 + v3732<<(uint(int32(1))%32)
	v3739 = int32(*(*int8)(unsafe.Add(mBase, uint32(v3738)+1)))
	v3740 = int32(*(*int8)(unsafe.Add(mBase, uint32(v3738))))
	v3754 = v3739
	v3755 = v3740
	goto L301
L311:
	;
	v3770 = v3556 + v3561 + int32(-2)
	goto L313
L312:
	;
	v3770 = v3566
	goto L313
L313:
	;
	if v3669+v3604 <= v3770 {
		v4159 = v3545
		v4180 = v3770
		v4182 = int32(0)
		goto L297
	} else {
		goto L314
	}
L314:
	;
	if v3770 < v3556 {
		v3856 = v3556
		v3887 = int32(0)
		goto L315
	} else {
		goto L316
	}
L315:
	;
	v3915 = v3856 + int32(-1)
	v3916 = *(*int32)(unsafe.Add(mBase, uint32(v142)))
	if v3916 == int32(0) {
		goto L323
	} else {
		goto L324
	}
L316:
	;
	v3795 = v3595 + v3564
	v3805 = int32(0)
	goto L319
L317:
	;
	v3856 = v3770 + int32(1)
	v3887 = v3842 & int32(4095)
	goto L315
L318:
	;
	v3850 = *(*int32)(unsafe.Add(mBase, uint32(v3795)))
	v3856 = v3556 - v3805
	v3887 = v3850 & int32(4095)
	goto L315
L319:
	;
	v3841 = v3795 + int32(4)
	v3842 = *(*int32)(unsafe.Add(mBase, uint32(v3841)))
	if int32(base.Ui32(v3842)>>(uint(int32(12))%32)) != v3552 {
		goto L318
	} else {
		goto L321
	}
L321:
	;
	v3847 = v3805 + int32(-1)
	if v3563-v3770 == v3847 {
		goto L317
	} else {
		goto L322
	}
L322:
	;
	v3795 = v3841
	v3805 = v3847
	goto L319
L323:
	;
	v4145 = *(*int64)(unsafe.Add(mBase, uint32(v3602+v3915<<(uint(int32(3))%32))))
	F_PushInterval(m, v142, v4145+v3545, v3856, v3887)
	mBase = m.M
	v4159 = v3545
	v4180 = v3915 + v3887
	v4182 = int32(0)
	goto L297
L324:
	;
	v3921 = v3602 + v3915<<(uint(int32(3))%32)
	v3946 = v3916
	goto L326
L325:
	;
	v4004 = v3602 + v3856<<(uint(int32(3))%32)
	v4026 = v3916
	goto L333
L326:
	;
	v3981 = *(*int32)(unsafe.Add(mBase, uint32(v3946)+8))
	if v3856 <= v3981 {
		goto L325
	} else {
		goto L328
	}
L327:
	;
	goto L325
L328:
	;
	v3983 = *(*int32)(unsafe.Add(mBase, uint32(v3946)+24))
	v3984 = *(*int32)(unsafe.Add(mBase, uint32(v3946)+12))
	if v3984 < v3856 {
		goto L329
	} else {
		goto L330
	}
L329:
	;
	if v3983 != 0 {
		v3946 = v3983
		goto L326
	} else {
		goto L332
	}
L330:
	;
	v3986 = *(*int64)(unsafe.Add(mBase, uint32(v3921)))
	v3987 = *(*int64)(unsafe.Add(mBase, uint32(v3946)))
	if v3986 <= v3987 {
		goto L329
	} else {
		goto L331
	}
L331:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3921))) = v3987
	v3990 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13])))
	v3991 = int32(1)
	v3994 = *(*int32)(unsafe.Add(mBase, uint32(v3946)+16))
	v3997 = v3915 - v3994 + v3991
	*(*uint16)(unsafe.Add(mBase, uint32(v3990+v3915<<(uint(v3991)%32)))) = uint16(v3997)
	goto L329
L332:
	;
	goto L327
L333:
	;
	v4064 = *(*int32)(unsafe.Add(mBase, uint32(v4026)+8))
	if v3856 < v4064 {
		goto L323
	} else {
		goto L335
	}
L334:
	;
	goto L323
L335:
	;
	v4066 = *(*int32)(unsafe.Add(mBase, uint32(v4026)+24))
	v4067 = *(*int32)(unsafe.Add(mBase, uint32(v4026)+12))
	if v4067 <= v3856 {
		goto L336
	} else {
		goto L337
	}
L336:
	;
	if v4066 != 0 {
		v4026 = v4066
		goto L333
	} else {
		goto L339
	}
L337:
	;
	v4069 = *(*int64)(unsafe.Add(mBase, uint32(v4004)))
	v4070 = *(*int64)(unsafe.Add(mBase, uint32(v4026)))
	if v4069 <= v4070 {
		goto L336
	} else {
		goto L338
	}
L338:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4004))) = v4070
	v4073 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13])))
	v4074 = int32(1)
	v4077 = *(*int32)(unsafe.Add(mBase, uint32(v4026)+16))
	v4080 = v3856 - v4077 + v4074
	*(*uint16)(unsafe.Add(mBase, uint32(v4073+v3856<<(uint(v4074)%32)))) = uint16(v4080)
	goto L336
L339:
	;
	goto L334
L340:
	;
	v4372 = int32(1)
	v4377 = v3556 + v4372
	if v4377 != v64 {
		v3545 = v4159
		v3552 = v3687
		v3556 = v4377
		v3561 = v3669
		v3563 = v3563 + v4372
		v3564 = v3564 + int32(4)
		v3566 = v4180
		v3568 = v4182
		goto L287
	} else {
		goto L359
	}
L341:
	;
	v4226 = v4209
	goto L342
L342:
	;
	v4271 = *(*int32)(unsafe.Add(mBase, uint32(v4226)+8))
	if v3556 < v4271 {
		goto L340
	} else {
		goto L344
	}
L343:
	;
	goto L340
L344:
	;
	v4273 = *(*int32)(unsafe.Add(mBase, uint32(v4226)+24))
	v4274 = *(*int32)(unsafe.Add(mBase, uint32(v4226)+12))
	if v3556 < v4274 {
		goto L346
	} else {
		goto L347
	}
L345:
	;
	if v4273 != 0 {
		v4226 = v4273
		goto L342
	} else {
		goto L358
	}
L346:
	;
	v4296 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	v4297 = v4296 + v3671
	v4298 = *(*int64)(unsafe.Add(mBase, uint32(v4297)))
	v4299 = *(*int64)(unsafe.Add(mBase, uint32(v4226)))
	if v4298 <= v4299 {
		goto L345
	} else {
		goto L357
	}
L347:
	;
	v4276 = *(*int32)(unsafe.Add(mBase, uint32(v4226)+20))
	if v4276 != 0 {
		goto L348
	} else {
		goto L349
	}
L348:
	;
	v4279 = v4276 + int32(24)
	goto L350
L349:
	;
	v4279 = v142
	goto L350
L350:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4279))) = v4273
	if v4273 == int32(0) {
		goto L351
	} else {
		goto L352
	}
L351:
	;
	if base.Ui32(v4226) < base.Ui32(v2519) {
		goto L354
	} else {
		goto L355
	}
L352:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4273)+20)) = v4276
	goto L351
L353:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4226)+24)) = v4290
	v4292 = *(*int32)(unsafe.Add(mBase, uint32(v142)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v142)+4)) = v4292 + int32(-1)
	goto L345
L354:
	;
	v4288 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1])))
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1]))) = v4226
	v4290 = v4288
	goto L353
L355:
	;
	if base.Ui32(v2492) < base.Ui32(v4226) {
		goto L354
	} else {
		goto L356
	}
L356:
	;
	v4286 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v4226
	v4290 = v4286
	goto L353
L357:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4297))) = v4299
	v4302 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[13])))
	v4303 = int32(1)
	v4306 = *(*int32)(unsafe.Add(mBase, uint32(v4226)+16))
	v4309 = v3556 - v4306 + v4303
	*(*uint16)(unsafe.Add(mBase, uint32(v4302+v3556<<(uint(v4303)%32)))) = uint16(v4309)
	goto L345
L358:
	;
	goto L343
L359:
	;
	goto L288
L360:
	;
	v4597 = v4441
	goto L34
L361:
	;
	v4597 = int32(0)
	goto L34
L362:
	;
	goto L361
L363:
	;
	v4509 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	F_WebPSafeFree(m, v4509)
	mBase = m.M
	v4511 = *(*int32)(unsafe.Add(mBase, uint32(v142)+8))
	F_WebPSafeFree(m, v4511)
	mBase = m.M
	v4513 = *(*int32)(unsafe.Add(mBase, uint32(v142)))
	if v4513 == int32(0) {
		goto L364
	} else {
		goto L365
	}
L364:
	;
	v4534 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v142))) = v4534
	v4536 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1])))
	if v4536 == v4534 {
		goto L373
	} else {
		goto L374
	}
L365:
	;
	v4521 = v4513
	goto L366
L366:
	;
	v4525 = *(*int32)(unsafe.Add(mBase, uint32(v4521)+24))
	if base.Ui32(v4521) < base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)) {
		goto L369
	} else {
		goto L370
	}
L367:
	;
	goto L364
L368:
	;
	if v4525 != 0 {
		v4521 = v4525
		goto L366
	} else {
		goto L372
	}
L369:
	;
	F_WebPSafeFree(m, v4521)
	mBase = m.M
	goto L368
L370:
	;
	if base.Ui32(v4521) <= base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)) {
		goto L368
	} else {
		goto L371
	}
L371:
	;
	goto L369
L372:
	;
	goto L367
L373:
	;
	v4559 = F_memset(m, v142, int32(0), int32(_a_F_VP8LBackwardReferencesTraceBackwards_1))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[3]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[4]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_5)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[5]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_6)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[6]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_7)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[7]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_8)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[8]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_9)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[9]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_10)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[10]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_11)
	*(*int32)(unsafe.Add(mBase, uint32(v4559)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[11]))) = v4559 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)
	goto L362
L374:
	;
	v4544 = v4536
	goto L375
L375:
	;
	v4548 = *(*int32)(unsafe.Add(mBase, uint32(v4544)+24))
	if base.Ui32(v4544) < base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)) {
		goto L378
	} else {
		goto L379
	}
L376:
	;
	goto L373
L377:
	;
	if v4548 != 0 {
		v4544 = v4548
		goto L375
	} else {
		goto L381
	}
L378:
	;
	F_WebPSafeFree(m, v4544)
	mBase = m.M
	goto L377
L379:
	;
	if base.Ui32(v4544) <= base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)) {
		goto L377
	} else {
		goto L380
	}
L380:
	;
	goto L378
L381:
	;
	goto L376
L382:
	;
	v4658 = v62 + int32(52)
	if v4658 == int32(0) {
		goto L384
	} else {
		goto L385
	}
L383:
	;
	v4666 = v4597
	goto L25
L384:
	;
	goto L383
L385:
	;
	v4661 = *(*int32)(unsafe.Add(mBase, uint32(v4658)))
	F_WebPSafeFree(m, v4661)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4658))) = int32(0)
	goto L384
L386:
	;
	F_free(m, v118)
	mBase = m.M
	goto L407
L387:
	;
	goto L386
L388:
	;
	v4730 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[0])))
	F_WebPSafeFree(m, v4730)
	mBase = m.M
	v4732 = *(*int32)(unsafe.Add(mBase, uint32(v142)+8))
	F_WebPSafeFree(m, v4732)
	mBase = m.M
	v4734 = *(*int32)(unsafe.Add(mBase, uint32(v142)))
	if v4734 == int32(0) {
		goto L389
	} else {
		goto L390
	}
L389:
	;
	v4755 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v142))) = v4755
	v4757 = *(*int32)(unsafe.Add(mBase, uint32(v142)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[1])))
	if v4757 == v4755 {
		goto L398
	} else {
		goto L399
	}
L390:
	;
	v4742 = v4734
	goto L391
L391:
	;
	v4746 = *(*int32)(unsafe.Add(mBase, uint32(v4742)+24))
	if base.Ui32(v4742) < base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)) {
		goto L394
	} else {
		goto L395
	}
L392:
	;
	goto L389
L393:
	;
	if v4746 != 0 {
		v4742 = v4746
		goto L391
	} else {
		goto L397
	}
L394:
	;
	F_WebPSafeFree(m, v4742)
	mBase = m.M
	goto L393
L395:
	;
	if base.Ui32(v4742) <= base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)) {
		goto L393
	} else {
		goto L396
	}
L396:
	;
	goto L394
L397:
	;
	goto L392
L398:
	;
	v4780 = F_memset(m, v142, int32(0), int32(_a_F_VP8LBackwardReferencesTraceBackwards_1))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[2]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[3]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_4)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[4]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_5)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[5]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_6)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[6]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_7)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[7]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_8)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[8]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_9)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[9]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_10)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[10]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_11)
	*(*int32)(unsafe.Add(mBase, uint32(v4780)+uint32(_c_F_VP8LBackwardReferencesTraceBackwards[11]))) = v4780 + int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)
	goto L387
L399:
	;
	v4765 = v4757
	goto L400
L400:
	;
	v4769 = *(*int32)(unsafe.Add(mBase, uint32(v4765)+24))
	if base.Ui32(v4765) < base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_12)) {
		goto L403
	} else {
		goto L404
	}
L401:
	;
	goto L398
L402:
	;
	if v4769 != 0 {
		v4765 = v4769
		goto L400
	} else {
		goto L406
	}
L403:
	;
	F_WebPSafeFree(m, v4765)
	mBase = m.M
	goto L402
L404:
	;
	if base.Ui32(v4765) <= base.Ui32(v142+int32(_a_F_VP8LBackwardReferencesTraceBackwards_3)) {
		goto L402
	} else {
		goto L405
	}
L405:
	;
	goto L403
L406:
	;
	goto L401
L407:
	;
	F_free(m, v142)
	mBase = m.M
	goto L408
L408:
	;
	if v4666 == int32(0) {
		v5724 = int32(0)
		goto L1
	} else {
		goto L409
	}
L409:
	;
	v4822 = v87 + v64<<(uint(int32(1))%32)
	v4824 = v4822 + int32(-2)
	if base.Ui32(v4824) < base.Ui32(v87) {
		v4914 = v4822
		goto L410
	} else {
		goto L411
	}
L410:
	;
	v4953 = int32(1)
	v4954 = (v4822 - v4914) >> (uint(v4953) % 32)
	if l3 < v4953 {
		goto L420
	} else {
		goto L421
	}
L411:
	;
	v4827 = v4824
	v4847 = v4822
	goto L412
L412:
	;
	v4886 = v4847 + int32(-2)
	v4887 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4827))))
	*(*uint16)(unsafe.Add(mBase, uint32(v4886))) = uint16(v4887)
	v4891 = v4827 - v4887<<(uint(int32(1))%32)
	if base.Ui32(v87) <= base.Ui32(v4891) {
		v4827 = v4891
		v4847 = v4886
		goto L412
	} else {
		goto L414
	}
L413:
	;
	v4914 = v4886
	goto L410
L414:
	;
	goto L413
L415:
	;
	v5706 = v62 + int32(52)
	if v5706 == int32(0) {
		goto L497
	} else {
		goto L498
	}
L416:
	;
	v5641 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	v5643 = base.B2i32(v5641 == int32(0))
	if l3 < int32(1) {
		v5724 = v5643
		goto L1
	} else {
		goto L495
	}
L417:
	;
	v5231 = v5010
	v5232 = int32(0)
	goto L460
L418:
	;
	v5017 = v4989
	v5037 = v4914
	v5038 = v4954
	goto L434
L419:
	;
	v4996 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	if v4996 == int32(0) {
		goto L431
	} else {
		goto L432
	}
L420:
	;
	v4975 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	if v4975 == int32(0) {
		goto L427
	} else {
		goto L428
	}
L421:
	;
	v4958 = v62 + int32(52)
	v4964 = F_WebPSafeCalloc(m, base.I64_extend_i32_s(int32(1)<<(uint(l3)%32)), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v4958))) = v4964
	if v4964 != 0 {
		goto L423
	} else {
		goto L424
	}
L422:
	;
	if v4964 != 0 {
		goto L419
	} else {
		goto L425
	}
L423:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4958)+8)) = l3
	*(*int32)(unsafe.Add(mBase, uint32(v4958)+4)) = int32(32) - l3
	goto L422
L424:
	;
	goto L422
L425:
	;
	v5724 = int32(0)
	goto L1
L426:
	;
	v4989 = int32(0)
	if v4989 < v4954 {
		goto L418
	} else {
		goto L429
	}
L427:
	;
	v4980 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v4980
	v4982 = *(*int32)(unsafe.Add(mBase, uint32(l6)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v4982
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = l6 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+8)) = v4980
	goto L426
L428:
	;
	v4978 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v4975))) = v4978
	goto L427
L429:
	;
	v4992 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	v5724 = base.B2i32(v4992 == int32(0))
	goto L1
L430:
	;
	v5010 = int32(0)
	if v5010 < v4954 {
		goto L417
	} else {
		goto L433
	}
L431:
	;
	v5001 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v5001
	v5003 = *(*int32)(unsafe.Add(mBase, uint32(l6)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v5003
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = l6 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l6)+8)) = v5001
	goto L430
L432:
	;
	v4999 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v4996))) = v4999
	goto L431
L433:
	;
	v5013 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	v5657 = base.B2i32(v5013 == int32(0))
	goto L415
L434:
	;
	v5075 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v5037))))
	if v5075 == int32(1) {
		goto L437
	} else {
		goto L438
	}
L436:
	;
	v5206 = v5038 + int32(-1)
	if v5206 != 0 {
		v5017 = v5017 + v5075
		v5037 = v5037 + int32(2)
		v5038 = v5206
		goto L434
	} else {
		goto L459
	}
L437:
	;
	v5145 = int64(*(*uint32)(unsafe.Add(mBase, uint32(l2+v5017<<(uint(int32(2))%32)))))
	v5149 = v5145<<(uint(int64(32))%64) | int64(65536)
	*(*int64)(unsafe.Add(mBase, uint32(v62)+32)) = v5149
	*(*int64)(unsafe.Add(mBase, uint32(v62))) = v5149
	v5154 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v5154 == int32(0) {
		goto L451
	} else {
		goto L452
	}
L438:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v62)+46)) = uint16(v5075)
	v5079 = int32(2)
	*(*uint8)(unsafe.Add(mBase, uint32(v62)+44)) = uint8(v5079)
	v5081 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v5085 = *(*int32)(unsafe.Add(mBase, uint32(v5081+v5017<<(uint(v5079)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v62)+48)) = int32(base.Ui32(v5085) >> (uint(int32(12)) % 32))
	v5089 = *(*int64)(unsafe.Add(mBase, uint32(v62)+44))
	*(*int64)(unsafe.Add(mBase, uint32(v62)+8)) = v5089
	v5095 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v5095 == int32(0) {
		goto L441
	} else {
		goto L442
	}
L439:
	;
	goto L436
L440:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5129)+8)) = v5130 + int32(1)
	v5134 = *(*int32)(unsafe.Add(mBase, uint32(v5129)+4))
	v5138 = *(*int64)(unsafe.Add(mBase, uint32(v62+int32(8))))
	*(*int64)(unsafe.Add(mBase, uint32(v5134+v5130<<(uint(int32(3))%32)))) = v5138
	goto L439
L441:
	;
	v5102 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v5102 != 0 {
		goto L445
	} else {
		goto L446
	}
L442:
	;
	v5098 = *(*int32)(unsafe.Add(mBase, uint32(v5095)+8))
	v5099 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v5098 != v5099 {
		v5129 = v5095
		v5130 = v5098
		goto L440
	} else {
		goto L443
	}
L443:
	;
	goto L441
L444:
	;
	v5120 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v5120))) = v5119
	v5122 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v5119)+8)) = v5122
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v5119
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v5119
	*(*int32)(unsafe.Add(mBase, uint32(v5119))) = v5122
	v5129 = v5119
	v5130 = v5122
	goto L440
L445:
	;
	v5117 = *(*int32)(unsafe.Add(mBase, uint32(v5102)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v5117
	v5119 = v5102
	goto L444
L446:
	;
	v5104 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v5109 = F_WebPSafeMalloc(m, int64(1), v5104<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v5109 != 0 {
		goto L447
	} else {
		goto L448
	}
L447:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5109)+4)) = v5109 + int32(12)
	v5119 = v5109
	goto L444
L448:
	;
	v5110 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v5110 | int32(1)
	goto L439
L449:
	;
	goto L436
L450:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5188)+8)) = v5189 + int32(1)
	v5193 = *(*int32)(unsafe.Add(mBase, uint32(v5188)+4))
	v5197 = *(*int64)(unsafe.Add(mBase, uint32(v62)))
	*(*int64)(unsafe.Add(mBase, uint32(v5193+v5189<<(uint(int32(3))%32)))) = v5197
	goto L449
L451:
	;
	v5161 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v5161 != 0 {
		goto L455
	} else {
		goto L456
	}
L452:
	;
	v5157 = *(*int32)(unsafe.Add(mBase, uint32(v5154)+8))
	v5158 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v5157 != v5158 {
		v5188 = v5154
		v5189 = v5157
		goto L450
	} else {
		goto L453
	}
L453:
	;
	goto L451
L454:
	;
	v5179 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v5179))) = v5178
	v5181 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v5178)+8)) = v5181
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v5178
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v5178
	*(*int32)(unsafe.Add(mBase, uint32(v5178))) = v5181
	v5188 = v5178
	v5189 = v5181
	goto L450
L455:
	;
	v5176 = *(*int32)(unsafe.Add(mBase, uint32(v5161)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v5176
	v5178 = v5161
	goto L454
L456:
	;
	v5163 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v5168 = F_WebPSafeMalloc(m, int64(1), v5163<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v5168 != 0 {
		goto L457
	} else {
		goto L458
	}
L457:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5168)+4)) = v5168 + int32(12)
	v5178 = v5168
	goto L454
L458:
	;
	v5169 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v5169 | int32(1)
	goto L449
L459:
	;
	goto L416
L460:
	;
	v5267 = int32(1)
	v5270 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4914+v5231<<(uint(v5267)%32)))))
	if v5270 == v5267 {
		goto L463
	} else {
		goto L464
	}
L461:
	;
	goto L416
L462:
	;
	v5580 = v5231 + int32(1)
	if v5580 != v4954 {
		v5231 = v5580
		v5232 = v5232 + v5270
		goto L460
	} else {
		goto L494
	}
L463:
	;
	v5439 = *(*int32)(unsafe.Add(mBase, uint32(v62)+52))
	v5440 = int32(2)
	v5442 = l2 + v5232<<(uint(v5440)%32)
	v5443 = *(*int32)(unsafe.Add(mBase, uint32(v5442)))
	v5446 = *(*int32)(unsafe.Add(mBase, uint32(v62)+56))
	v5447 = int32(base.Ui32(v5443*int32(506832829)) >> (uint(v5446) % 32))
	v5450 = v5439 + v5447<<(uint(v5440)%32)
	v5451 = *(*int32)(unsafe.Add(mBase, uint32(v5450)))
	if v5451 != v5443 {
		goto L481
	} else {
		goto L482
	}
L464:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v62)+46)) = uint16(v5270)
	v5274 = int32(2)
	*(*uint8)(unsafe.Add(mBase, uint32(v62)+44)) = uint8(v5274)
	v5276 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v5278 = v5232 << (uint(v5274) % 32)
	v5280 = *(*int32)(unsafe.Add(mBase, uint32(v5276+v5278)))
	*(*int32)(unsafe.Add(mBase, uint32(v62)+48)) = int32(base.Ui32(v5280) >> (uint(int32(12)) % 32))
	v5284 = *(*int64)(unsafe.Add(mBase, uint32(v62)+44))
	*(*int64)(unsafe.Add(mBase, uint32(v62)+24)) = v5284
	v5290 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v5290 == int32(0) {
		goto L467
	} else {
		goto L468
	}
L465:
	;
	if v5270 == int32(0) {
		goto L462
	} else {
		goto L475
	}
L466:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5324)+8)) = v5325 + int32(1)
	v5329 = *(*int32)(unsafe.Add(mBase, uint32(v5324)+4))
	v5333 = *(*int64)(unsafe.Add(mBase, uint32(v62+int32(24))))
	*(*int64)(unsafe.Add(mBase, uint32(v5329+v5325<<(uint(int32(3))%32)))) = v5333
	goto L465
L467:
	;
	v5297 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v5297 != 0 {
		goto L471
	} else {
		goto L472
	}
L468:
	;
	v5293 = *(*int32)(unsafe.Add(mBase, uint32(v5290)+8))
	v5294 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v5293 != v5294 {
		v5324 = v5290
		v5325 = v5293
		goto L466
	} else {
		goto L469
	}
L469:
	;
	goto L467
L470:
	;
	v5315 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v5315))) = v5314
	v5317 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v5314)+8)) = v5317
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v5314
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v5314
	*(*int32)(unsafe.Add(mBase, uint32(v5314))) = v5317
	v5324 = v5314
	v5325 = v5317
	goto L466
L471:
	;
	v5312 = *(*int32)(unsafe.Add(mBase, uint32(v5297)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v5312
	v5314 = v5297
	goto L470
L472:
	;
	v5299 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v5304 = F_WebPSafeMalloc(m, int64(1), v5299<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v5304 != 0 {
		goto L473
	} else {
		goto L474
	}
L473:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5304)+4)) = v5304 + int32(12)
	v5314 = v5304
	goto L470
L474:
	;
	v5305 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v5305 | int32(1)
	goto L465
L475:
	;
	v5344 = *(*int32)(unsafe.Add(mBase, uint32(v62)+52))
	v5346 = l2 + v5278
	v5358 = v5270 & int32(_a_F_VP8LBackwardReferencesTraceBackwards_14)
	goto L476
L476:
	;
	v5404 = *(*int32)(unsafe.Add(mBase, uint32(v5346)))
	v5405 = int32(506832829)
	v5407 = *(*int32)(unsafe.Add(mBase, uint32(v62)+56))
	v5409 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v5344+int32(base.Ui32(v5404*v5405)>>(uint(v5407)%32))<<(uint(v5409)%32)))) = v5404
	v5415 = *(*int32)(unsafe.Add(mBase, uint32(v5346+int32(4))))
	v5418 = *(*int32)(unsafe.Add(mBase, uint32(v62)+56))
	*(*int32)(unsafe.Add(mBase, uint32(v5344+int32(base.Ui32(v5415*v5405)>>(uint(v5418)%32))<<(uint(v5409)%32)))) = v5415
	v5425 = v5346 + int32(8)
	v5427 = v5358 + int32(-2)
	if v5427 != 0 {
		v5346 = v5425
		v5358 = v5427
		goto L476
	} else {
		goto L478
	}
L477:
	;
	if v5270&int32(1) == int32(0) {
		goto L462
	} else {
		goto L479
	}
L478:
	;
	goto L477
L479:
	;
	v5430 = *(*int32)(unsafe.Add(mBase, uint32(v5425)))
	v5433 = *(*int32)(unsafe.Add(mBase, uint32(v62)+56))
	*(*int32)(unsafe.Add(mBase, uint32(v5344+int32(base.Ui32(v5430*int32(506832829))>>(uint(v5433)%32))<<(uint(int32(2))%32)))) = v5430
	goto L462
L480:
	;
	v5465 = base.I64_extend_i32_u(v5460)<<(uint(int64(32))%64) | v5461
	*(*int64)(unsafe.Add(mBase, uint32(v62)+16)) = v5465
	*(*int64)(unsafe.Add(mBase, uint32(v62)+32)) = v5465
	v5472 = *(*int32)(unsafe.Add(mBase, uint32(l6)+20))
	if v5472 == int32(0) {
		goto L486
	} else {
		goto L487
	}
L481:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5450))) = v5443
	v5458 = *(*int32)(unsafe.Add(mBase, uint32(v5442)))
	v5460 = v5458
	v5461 = int64(65536)
	goto L480
L482:
	;
	if int32(-1) < v5447 {
		v5460 = v5447
		v5461 = int64(65537)
		goto L480
	} else {
		goto L483
	}
L483:
	;
	goto L481
L484:
	;
	goto L462
L485:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5506)+8)) = v5507 + int32(1)
	v5511 = *(*int32)(unsafe.Add(mBase, uint32(v5506)+4))
	v5515 = *(*int64)(unsafe.Add(mBase, uint32(v62+int32(16))))
	*(*int64)(unsafe.Add(mBase, uint32(v5511+v5507<<(uint(int32(3))%32)))) = v5515
	goto L484
L486:
	;
	v5479 = *(*int32)(unsafe.Add(mBase, uint32(l6)+16))
	if v5479 != 0 {
		goto L490
	} else {
		goto L491
	}
L487:
	;
	v5475 = *(*int32)(unsafe.Add(mBase, uint32(v5472)+8))
	v5476 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	if v5475 != v5476 {
		v5506 = v5472
		v5507 = v5475
		goto L485
	} else {
		goto L488
	}
L488:
	;
	goto L486
L489:
	;
	v5497 = *(*int32)(unsafe.Add(mBase, uint32(l6)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v5497))) = v5496
	v5499 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v5496)+8)) = v5499
	*(*int32)(unsafe.Add(mBase, uint32(l6)+20)) = v5496
	*(*int32)(unsafe.Add(mBase, uint32(l6)+12)) = v5496
	*(*int32)(unsafe.Add(mBase, uint32(v5496))) = v5499
	v5506 = v5496
	v5507 = v5499
	goto L485
L490:
	;
	v5494 = *(*int32)(unsafe.Add(mBase, uint32(v5479)))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+16)) = v5494
	v5496 = v5479
	goto L489
L491:
	;
	v5481 = *(*int32)(unsafe.Add(mBase, uint32(l6)))
	v5486 = F_WebPSafeMalloc(m, int64(1), v5481<<(uint(int32(3))%32)+int32(12))
	mBase = m.M
	if v5486 != 0 {
		goto L492
	} else {
		goto L493
	}
L492:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5486)+4)) = v5486 + int32(12)
	v5496 = v5486
	goto L489
L493:
	;
	v5487 = *(*int32)(unsafe.Add(mBase, uint32(l6)+4))
	*(*int32)(unsafe.Add(mBase, uint32(l6)+4)) = v5487 | int32(1)
	goto L484
L494:
	;
	goto L461
L495:
	;
	v5657 = v5643
	goto L415
L496:
	;
	v5724 = v5657
	goto L1
L497:
	;
	goto L496
L498:
	;
	v5709 = *(*int32)(unsafe.Add(mBase, uint32(v5706)))
	F_WebPSafeFree(m, v5709)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v5706))) = int32(0)
	goto L497
L499:
	;
	m.G0 = v62 + int32(64)
	return v5724
}
func F_VP8LBackwardRefsClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v4 == int32(0) {
	} else {
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		*(*int32)(unsafe.Add(mBase, uint32(v4))) = v7
	}
	v9 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v9
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v11
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = l0 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v9
	if v11 == v9 {
	} else {
		v22 = v11
		for {
			v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
			F_free(m, v22)
			mBase = m.M
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v23
			if v23 != 0 {
				v22 = v23
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8LBackwardRefsCursorAdd(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int64
	_ = v13
	var v14 int32
	_ = v14
	var v18 int32
	_ = v18
	var v25 int64
	_ = v25
	var v26 int32
	_ = v26
	var v27 int64
	_ = v27
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v69 int64
	_ = v69
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	if v5 == int32(0) {
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		if v12 != 0 {
			v48 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
			*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v48
			v50 = v12
			v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
			v53 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
			*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
			*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
			*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
			v60 = v50
			v61 = v53
			*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
			v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
			v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
			*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
			return
		} else {
			v13 = int64(1)
			v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
			v18 = v14<<(uint(int32(3))%32) + int32(12)
			v25 = base.I64_div_u_s(int64(2147418112), v13)
			v26 = int32(0)
			v27 = base.I64_extend_i32_u(v18)
			if base.Ui64(int64(4294967295)) < base.Ui64(v27*v13) {
				v39 = v26
			} else {
				if base.Ui64(v25) < base.Ui64(v27) {
					v39 = v26
				} else {
					v37 = F_malloc(m, base.I32_wrap_i64(v13)*v18)
					mBase = m.M
					v39 = v37
				}
			}
			if v39 != 0 {
				*(*int32)(unsafe.Add(mBase, uint32(v39)+4)) = v39 + int32(12)
				v50 = v39
				v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
				v53 = int32(0)
				*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
				v60 = v50
				v61 = v53
				*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
				v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
				v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
				*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
				return
			} else {
				v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v41 | int32(1)
				return
			}
		}
	} else {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(v5)+8))
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v8 != v9 {
			v60 = v5
			v61 = v8
			*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
			v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
			v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
			*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
			return
		} else {
			v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
			if v12 != 0 {
				v48 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
				*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v48
				v50 = v12
				v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
				*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
				v53 = int32(0)
				*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
				*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
				v60 = v50
				v61 = v53
				*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
				v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
				v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
				*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
				return
			} else {
				v13 = int64(1)
				v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				v18 = v14<<(uint(int32(3))%32) + int32(12)
				v25 = base.I64_div_u_s(int64(2147418112), v13)
				v26 = int32(0)
				v27 = base.I64_extend_i32_u(v18)
				if base.Ui64(int64(4294967295)) < base.Ui64(v27*v13) {
					v39 = v26
				} else {
					if base.Ui64(v25) < base.Ui64(v27) {
						v39 = v26
					} else {
						v37 = F_malloc(m, base.I32_wrap_i64(v13)*v18)
						mBase = m.M
						v39 = v37
					}
				}
				if v39 != 0 {
					*(*int32)(unsafe.Add(mBase, uint32(v39)+4)) = v39 + int32(12)
					v50 = v39
					v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					*(*int32)(unsafe.Add(mBase, uint32(v51))) = v50
					v53 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(v50)+8)) = v53
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v50
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v50
					*(*int32)(unsafe.Add(mBase, uint32(v50))) = v53
					v60 = v50
					v61 = v53
					*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v61 + int32(1)
					v65 = *(*int32)(unsafe.Add(mBase, uint32(v60)+4))
					v69 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
					*(*int64)(unsafe.Add(mBase, uint32(v65+v61<<(uint(int32(3))%32)))) = v69
					return
				} else {
					v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v41 | int32(1)
					return
				}
			}
		}
	}
}
func F_VP8LBackwardRefsInit(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int64
	_ = v6
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	v5 = l0 + int32(12)
	v6 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v5))) = v6
	*(*int64)(unsafe.Add(mBase, uint32(l0)+4)) = v6
	*(*int32)(unsafe.Add(mBase, uint32(l0+int32(20)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v5))) = l0 + int32(8)
	v17 = int32(256)
	if v17 < l1 {
		v20 = l1
	} else {
		v20 = v17
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v20
	return
}
func F_VP8LBitWriterClone(m *base.Module, l0 int32, l1 int32) int32 {
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
	var v20 int64
	_ = v20
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v32 int64
	_ = v32
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v49 int64
	_ = v49
	var v50 int32
	_ = v50
	var v51 int64
	_ = v51
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int64
	_ = v83
	var v85 int32
	_ = v85
	var v87 int32
	_ = v87
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v14 = v12 - v13
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v18 = v16 - v17
	v20 = base.I64_extend_i32_u(v14) + base.I64_extend_i32_u(v18)
	if base.Ui64(v20) < base.Ui64(int64(4294967296)) {
		v27 = base.I32_wrap_i64(v20)
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
		v29 = v28 - v13
		if v28 == v13 {
			v32 = int64(1)
			v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
			if base.Ui32(v27) < base.Ui32(v36) {
				v38 = v36
			} else {
				v38 = v27
			}
			v42 = v38&int32(-1024) + int32(1024)
			v49 = base.I64_div_u_s(int64(2147418112), v32)
			v50 = int32(0)
			v51 = base.I64_extend_i32_u(v42)
			if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
				v63 = v50
			} else {
				if base.Ui64(v49) < base.Ui64(v51) {
					v63 = v50
				} else {
					v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
					mBase = m.M
					v63 = v61
				}
			}
			if v63 != 0 {
				if v12 == v13 {
				} else {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v71 = F_memcpy(m, v63, v70, v14)
					mBase = m.M
				}
				v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				F_free(m, v72)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v63
				*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v63 + v42
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v63 + v14
				v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v80 = v79
				v81 = v63
				v82 = F_memcpy(m, v81, v80, v18)
				mBase = m.M
				v83 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
				*(*int64)(unsafe.Add(mBase, uint32(l1))) = v83
				v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v85
				v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v87 + v18
				return int32(1)
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
				return int32(0)
			}
		} else {
			if base.Ui32(v29) < base.Ui32(v27) {
				v32 = int64(1)
				v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
				if base.Ui32(v27) < base.Ui32(v36) {
					v38 = v36
				} else {
					v38 = v27
				}
				v42 = v38&int32(-1024) + int32(1024)
				v49 = base.I64_div_u_s(int64(2147418112), v32)
				v50 = int32(0)
				v51 = base.I64_extend_i32_u(v42)
				if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
					v63 = v50
				} else {
					if base.Ui64(v49) < base.Ui64(v51) {
						v63 = v50
					} else {
						v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
						mBase = m.M
						v63 = v61
					}
				}
				if v63 != 0 {
					if v12 == v13 {
					} else {
						v70 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
						v71 = F_memcpy(m, v63, v70, v14)
						mBase = m.M
					}
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					F_free(m, v72)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l1)+8)) = v63
					*(*int32)(unsafe.Add(mBase, uint32(l1)+16)) = v63 + v42
					*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v63 + v14
					v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v80 = v79
					v81 = v63
					v82 = F_memcpy(m, v81, v80, v18)
					mBase = m.M
					v83 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
					*(*int64)(unsafe.Add(mBase, uint32(l1))) = v83
					v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v85
					v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v87 + v18
					return int32(1)
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
					return int32(0)
				}
			} else {
				v80 = v17
				v81 = v13
				v82 = F_memcpy(m, v81, v80, v18)
				mBase = m.M
				v83 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
				*(*int64)(unsafe.Add(mBase, uint32(l1))) = v83
				v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v85
				v87 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v87 + v18
				return int32(1)
			}
		}
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = int32(1)
		return int32(0)
	}
}
func F_VP8LBitWriterFinish(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v20 int64
	_ = v20
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v32 int64
	_ = v32
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v42 int32
	_ = v42
	var v49 int64
	_ = v49
	var v50 int32
	_ = v50
	var v51 int64
	_ = v51
	var v61 int32
	_ = v61
	var v63 int32
	_ = v63
	var v67 int32
	_ = v67
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v94 int32
	_ = v94
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v120 int32
	_ = v120
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v12 = v10 - v11
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v20 = base.I64_extend_i32_u(v12) + base.I64_extend_i32_u((v14+int32(7))>>(uint(int32(3))%32))
	if base.Ui64(v20) < base.Ui64(int64(4294967296)) {
		v27 = base.I32_wrap_i64(v20)
		v28 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		v29 = v28 - v11
		if v28 == v11 {
			v32 = int64(1)
			v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
			if base.Ui32(v27) < base.Ui32(v36) {
				v38 = v36
			} else {
				v38 = v27
			}
			v42 = v38&int32(-1024) + int32(1024)
			v49 = base.I64_div_u_s(int64(2147418112), v32)
			v50 = int32(0)
			v51 = base.I64_extend_i32_u(v42)
			if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
				v63 = v50
			} else {
				if base.Ui64(v49) < base.Ui64(v51) {
					v63 = v50
				} else {
					v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
					mBase = m.M
					v63 = v61
				}
			}
			if v63 != 0 {
				if v10 == v11 {
				} else {
					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v71 = F_memcpy(m, v63, v70, v12)
					mBase = m.M
				}
				v72 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				F_free(m, v72)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v63 + v42
				*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v63 + v12
				*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v63
				v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v80 = v79
				if v80 < int32(1) {
				} else {
					v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
					v87 = v84
					for {
						v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v94 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v94))) = uint8(v87)
						v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v100 = int32(8)
						v101 = int32(base.Ui32(v99) >> (uint(v100) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v101
						v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v103 + int32(-8)
						if v100 < v103 {
							v87 = v101
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
				v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				return v120
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
				v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				return v67
			}
		} else {
			if base.Ui32(v27) <= base.Ui32(v29) {
				v80 = v14
				if v80 < int32(1) {
				} else {
					v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
					v87 = v84
					for {
						v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v94 + int32(1)
						*(*uint8)(unsafe.Add(mBase, uint32(v94))) = uint8(v87)
						v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v100 = int32(8)
						v101 = int32(base.Ui32(v99) >> (uint(v100) % 32))
						*(*int32)(unsafe.Add(mBase, uint32(l0))) = v101
						v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
						*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v103 + int32(-8)
						if v100 < v103 {
							v87 = v101
							continue
						} else {
							break
						}
						break
					}
				}
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
				v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				return v120
			} else {
				v32 = int64(1)
				v36 = int32(base.Ui32(v29*int32(3)) >> (uint(int32(1)) % 32))
				if base.Ui32(v27) < base.Ui32(v36) {
					v38 = v36
				} else {
					v38 = v27
				}
				v42 = v38&int32(-1024) + int32(1024)
				v49 = base.I64_div_u_s(int64(2147418112), v32)
				v50 = int32(0)
				v51 = base.I64_extend_i32_u(v42)
				if base.Ui64(int64(4294967295)) < base.Ui64(v51*v32) {
					v63 = v50
				} else {
					if base.Ui64(v49) < base.Ui64(v51) {
						v63 = v50
					} else {
						v61 = F_malloc(m, base.I32_wrap_i64(v32)*v42)
						mBase = m.M
						v63 = v61
					}
				}
				if v63 != 0 {
					if v10 == v11 {
					} else {
						v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v71 = F_memcpy(m, v63, v70, v12)
						mBase = m.M
					}
					v72 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					F_free(m, v72)
					mBase = m.M
					*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v63 + v42
					*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v63 + v12
					*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v63
					v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					v80 = v79
					if v80 < int32(1) {
					} else {
						v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						v87 = v84
						for {
							v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v94 + int32(1)
							*(*uint8)(unsafe.Add(mBase, uint32(v94))) = uint8(v87)
							v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
							v100 = int32(8)
							v101 = int32(base.Ui32(v99) >> (uint(v100) % 32))
							*(*int32)(unsafe.Add(mBase, uint32(l0))) = v101
							v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
							*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v103 + int32(-8)
							if v100 < v103 {
								v87 = v101
								continue
							} else {
								break
							}
							break
						}
					}
					*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
					v120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					return v120
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
					v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					return v67
				}
			}
		}
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
		v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		return v25
	}
}
func F_VP8LBitWriterInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int64
	_ = v4
	var v14 int64
	_ = v14
	var v18 int32
	_ = v18
	var v25 int64
	_ = v25
	var v26 int32
	_ = v26
	var v27 int64
	_ = v27
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v45 int32
	_ = v45
	v4 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(16)))) = v4
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(8)))) = v4
	v14 = int64(1)
	v18 = l1&int32(-1024) + int32(1024)
	v25 = base.I64_div_u_s(int64(2147418112), v14)
	v26 = int32(0)
	v27 = base.I64_extend_i32_u(v18)
	if base.Ui64(int64(4294967295)) < base.Ui64(v27*v14) {
		v39 = v26
	} else {
		if base.Ui64(v25) < base.Ui64(v27) {
			v39 = v26
		} else {
			v37 = F_malloc(m, base.I32_wrap_i64(v14)*v18)
			mBase = m.M
			v39 = v37
		}
	}
	if v39 != 0 {
		v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		F_free(m, v45)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v39 + v18
		*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v39
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v39
		return int32(1)
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
		return int32(0)
	}
}
func F_VP8LBitWriterReset(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v3 int64
	_ = v3
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	v3 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
	*(*int64)(unsafe.Add(mBase, uint32(l1))) = v3
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+20)) = v5
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l1)+12)) = v7 + (v8 - v9)
	return
}
func F_VP8LBitWriterSwap(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v11 int64
	_ = v11
	var v12 int64
	_ = v12
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int64
	_ = v16
	var v18 int32
	_ = v18
	var v19 int64
	_ = v19
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int64
	_ = v23
	var v25 int32
	_ = v25
	var v26 int64
	_ = v26
	v11 = *(*int64)(unsafe.Add(mBase, uint32(l0)))
	v12 = *(*int64)(unsafe.Add(mBase, uint32(l1)))
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = v12
	v14 = int32(16)
	v15 = l0 + v14
	v16 = *(*int64)(unsafe.Add(mBase, uint32(v15)))
	v18 = l1 + v14
	v19 = *(*int64)(unsafe.Add(mBase, uint32(v18)))
	*(*int64)(unsafe.Add(mBase, uint32(v15))) = v19
	v21 = int32(8)
	v22 = l0 + v21
	v23 = *(*int64)(unsafe.Add(mBase, uint32(v22)))
	v25 = l1 + v21
	v26 = *(*int64)(unsafe.Add(mBase, uint32(v25)))
	*(*int64)(unsafe.Add(mBase, uint32(v22))) = v26
	*(*int64)(unsafe.Add(mBase, uint32(l1))) = v11
	*(*int64)(unsafe.Add(mBase, uint32(v18))) = v16
	*(*int64)(unsafe.Add(mBase, uint32(v25))) = v23
	return
}
func F_VP8LBitWriterWipeOut(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v11 int64
	_ = v11
	if l0 == int32(0) {
	} else {
		v6 = l0 + int32(8)
		v7 = *(*int32)(unsafe.Add(mBase, uint32(v6)))
		F_free(m, v7)
		mBase = m.M
		v11 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(16)))) = v11
		*(*int64)(unsafe.Add(mBase, uint32(v6))) = v11
		*(*int64)(unsafe.Add(mBase, uint32(l0))) = v11
	}
	return
}
func F_VP8LBitsEntropy(m *base.Module, l0 int32, l1 int32) int64 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v13 int32
	_ = v13
	var v25 int64
	_ = v25
	var v28 int64
	_ = v28
	var v29 int64
	_ = v29
	var v32 int64
	_ = v32
	var v35 int64
	_ = v35
	var v37 int64
	_ = v37
	var v38 int64
	_ = v38
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v53 int64
	_ = v53
	var v56 int64
	_ = v56
	var v59 int64
	_ = v59
	var v61 int64
	_ = v61
	var v62 int64
	_ = v62
	v6 = m.G0
	v8 = v6 - int32(32)
	m.G0 = v8
	F_VP8LBitsEntropyUnrefined(m, l0, l1, v8+int32(8))
	mBase = m.M
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v8)+20))
	if v13 <= int32(4) {
		if v13 < int32(2) {
			v62 = int64(0)
		} else {
			switch v13 + int32(-2) {
			case 0:
				v25 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v8)+16)))
				v28 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
				v29 = v25*int64(830472192) + v28
				if v29 < int64(0) {
					v32 = int64(-50)
				} else {
					v32 = int64(50)
				}
				v35 = base.I64_div_s(v32+v29, int64(100))
				v62 = v35
			case 1:
				v37 = int64(950)
				v38 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
				v41 = *(*int32)(unsafe.Add(mBase, uint32(v8)+16))
				v44 = *(*int32)(unsafe.Add(mBase, uint32(v8)+24))
				v53 = v37*base.I64_extend_i32_u(v41<<(uint(int32(1))%32)-v44)<<(uint(int64(23))%64) + v38*(int64(1000)-v37)
				if v53 < int64(0) {
					v56 = int64(-500)
				} else {
					v56 = int64(500)
				}
				v59 = base.I64_div_s(v56+v53, int64(1000))
				if base.Ui64(v59) < base.Ui64(v38) {
					v61 = v38
				} else {
					v61 = v59
				}
				v62 = v61
			default:
				v37 = int64(700)
				v38 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
				v41 = *(*int32)(unsafe.Add(mBase, uint32(v8)+16))
				v44 = *(*int32)(unsafe.Add(mBase, uint32(v8)+24))
				v53 = v37*base.I64_extend_i32_u(v41<<(uint(int32(1))%32)-v44)<<(uint(int64(23))%64) + v38*(int64(1000)-v37)
				if v53 < int64(0) {
					v56 = int64(-500)
				} else {
					v56 = int64(500)
				}
				v59 = base.I64_div_s(v56+v53, int64(1000))
				if base.Ui64(v59) < base.Ui64(v38) {
					v61 = v38
				} else {
					v61 = v59
				}
				v62 = v61
			}
		}
	} else {
		v37 = int64(627)
		v38 = *(*int64)(unsafe.Add(mBase, uint32(v8)+8))
		v41 = *(*int32)(unsafe.Add(mBase, uint32(v8)+16))
		v44 = *(*int32)(unsafe.Add(mBase, uint32(v8)+24))
		v53 = v37*base.I64_extend_i32_u(v41<<(uint(int32(1))%32)-v44)<<(uint(int64(23))%64) + v38*(int64(1000)-v37)
		if v53 < int64(0) {
			v56 = int64(-500)
		} else {
			v56 = int64(500)
		}
		v59 = base.I64_div_s(v56+v53, int64(1000))
		if base.Ui64(v59) < base.Ui64(v38) {
			v61 = v38
		} else {
			v61 = v59
		}
		v62 = v61
	}
	m.G0 = v8 + int32(32)
	return v62
}
func F_VP8LBitsEntropyUnrefined(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v4 int64
	_ = v4
	var v22 int32
	_ = v22
	var v26 int32
	_ = v26
	var v29 int64
	_ = v29
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
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v47 int32
	_ = v47
	var v53 int64
	_ = v53
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v58 int64
	_ = v58
	var v59 int64
	_ = v59
	var v60 int64
	_ = v60
	var v63 int64
	_ = v63
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v66 int32
	_ = v66
	var v71 int32
	_ = v71
	var v78 int64
	_ = v78
	var v79 int32
	_ = v79
	var v85 int32
	_ = v85
	var v91 int64
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v96 int64
	_ = v96
	var v100 int64
	_ = v100
	var v106 int64
	_ = v106
	v4 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l2)+8)) = v4
	*(*int64)(unsafe.Add(mBase, uint32(l2+int32(16)))) = int64(281470681743360)
	if int32(1) <= l1 {
		v22 = int32(0)
		v26 = l0
		v29 = int64(0)
		v30 = v22
		v31 = v22
		v32 = v22
		v33 = v22
		for {
			v36 = *(*int32)(unsafe.Add(mBase, uint32(v26)))
			if v36 == int32(0) {
				v63 = v29
				v64 = v30
				v65 = v31
				v66 = v32
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l2)+20)) = v33
				v40 = v30 + v36
				*(*int32)(unsafe.Add(mBase, uint32(l2)+8)) = v40
				v43 = v32 + int32(1)
				*(*int32)(unsafe.Add(mBase, uint32(l2)+12)) = v43
				if base.Ui32(int32(255)) < base.Ui32(v36) {
					v54 = m.G1
					v57 = *(*int32)(unsafe.Add(mBase, uint32(v54)+uint32(_c_F_VP8LBitsEntropyUnrefined[0])))
					v58 = m.T0[v57].(func(*base.Module, int32) int64)(m, v36)
					mBase = m.M
					v59 = v58
				} else {
					v47 = m.G1
					v53 = *(*int64)(unsafe.Add(mBase, uint32(v47+int32(_a_F_VP8LBitsEntropyUnrefined_0)+v36<<(uint(int32(3))%32))))
					v59 = v53
				}
				v60 = v29 + v59
				if base.Ui32(v36) <= base.Ui32(v31) {
					v63 = v60
					v64 = v40
					v65 = v31
					v66 = v43
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l2)+16)) = v36
					v63 = v60
					v64 = v40
					v65 = v36
					v66 = v43
				}
			}
			v71 = v33 + int32(1)
			if l1 != v71 {
				v26 = v26 + int32(4)
				v29 = v63
				v30 = v64
				v31 = v65
				v32 = v66
				v33 = v71
				continue
			} else {
				break
			}
			break
		}
		if base.Ui32(int32(255)) < base.Ui32(v64) {
			v92 = m.G1
			v95 = *(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_VP8LBitsEntropyUnrefined[0])))
			v96 = m.T0[v95].(func(*base.Module, int32) int64)(m, v64)
			mBase = m.M
			v100 = v63
			v106 = v96
		} else {
			v78 = v63
			v79 = v64
			v85 = m.G1
			v91 = *(*int64)(unsafe.Add(mBase, uint32(v85+int32(_a_F_VP8LBitsEntropyUnrefined_0)+v79<<(uint(int32(3))%32))))
			v100 = v78
			v106 = v91
		}
	} else {
		v78 = v4
		v79 = int32(0)
		v85 = m.G1
		v91 = *(*int64)(unsafe.Add(mBase, uint32(v85+int32(_a_F_VP8LBitsEntropyUnrefined_0)+v79<<(uint(int32(3))%32))))
		v100 = v78
		v106 = v91
	}
	*(*int64)(unsafe.Add(mBase, uint32(l2))) = v106 - v100
	return
}
func F_VP8LClearBackwardRefs(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if v3 == int32(0) {
	} else {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		*(*int32)(unsafe.Add(mBase, uint32(v3))) = v6
	}
	v8 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v8
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v10
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = l0 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v8
	return
}
func F_VP8LColorCacheClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	if l0 == int32(0) {
	} else {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		F_free(m, v4)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(0)
	}
	return
}
func F_VP8LColorCacheInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int64
	_ = v6
	var v7 int32
	_ = v7
	var v14 int64
	_ = v14
	var v15 int32
	_ = v15
	var v16 int64
	_ = v16
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	v6 = base.I64_extend_i32_s(int32(1) << (uint(l1) % 32))
	v7 = int32(4)
	if v6 == int64(0) {
		v25 = F_calloc(m, base.I32_wrap_i64(v6), v7)
		mBase = m.M
		v27 = v25
	} else {
		v14 = base.I64_div_u_s(int64(2147418112), v6)
		v15 = int32(0)
		v16 = base.I64_extend_i32_u(v7)
		if base.Ui64(int64(4294967295)) < base.Ui64(v16*v6) {
			v27 = v15
		} else {
			if base.Ui64(v14) < base.Ui64(v16) {
				v27 = v15
			} else {
				v25 = F_calloc(m, base.I32_wrap_i64(v6), v7)
				mBase = m.M
				v27 = v25
			}
		}
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v27
	if v27 != 0 {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = l1
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(32) - l1
		return int32(1)
	} else {
		return int32(0)
	}
}
func F_VP8LColorSpaceTransform(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32) int32 {
	mBase := m.M
	_ = mBase
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v49 int32
	_ = v49
	var v298 int32
	_ = v298
	var v306 int32
	_ = v306
	var v308 int32
	_ = v308
	var v310 int32
	_ = v310
	var v312 int32
	_ = v312
	var v316 int32
	_ = v316
	var v320 int32
	_ = v320
	var v323 int32
	_ = v323
	var v325 int32
	_ = v325
	var v348 int32
	_ = v348
	var v371 int32
	_ = v371
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v377 int32
	_ = v377
	var v381 int32
	_ = v381
	var v387 int32
	_ = v387
	var v429 int32
	_ = v429
	var v440 int32
	_ = v440
	var v447 int32
	_ = v447
	var v450 int32
	_ = v450
	var v453 int32
	_ = v453
	var v456 int32
	_ = v456
	var v457 int32
	_ = v457
	var v460 int32
	_ = v460
	var v465 int32
	_ = v465
	var v466 int32
	_ = v466
	var v476 int32
	_ = v476
	var v478 int32
	_ = v478
	var v480 int32
	_ = v480
	var v490 int32
	_ = v490
	var v497 int32
	_ = v497
	var v499 int32
	_ = v499
	var v506 int32
	_ = v506
	var v510 int32
	_ = v510
	var v517 int32
	_ = v517
	var v519 int32
	_ = v519
	var v522 int32
	_ = v522
	var v537 int32
	_ = v537
	var v558 int32
	_ = v558
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v572 int32
	_ = v572
	var v574 int32
	_ = v574
	var v589 int32
	_ = v589
	var v610 int32
	_ = v610
	var v635 int32
	_ = v635
	var v657 int32
	_ = v657
	var v666 int32
	_ = v666
	var v667 int32
	_ = v667
	var v668 int32
	_ = v668
	var v672 int32
	_ = v672
	var v673 int32
	_ = v673
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v682 int32
	_ = v682
	var v683 int32
	_ = v683
	var v687 int32
	_ = v687
	var v735 int32
	_ = v735
	var v737 int32
	_ = v737
	var v739 int32
	_ = v739
	var v785 int32
	_ = v785
	var v786 int32
	_ = v786
	var v787 int32
	_ = v787
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v808 int32
	_ = v808
	var v811 int32
	_ = v811
	var v830 int32
	_ = v830
	var v836 int32
	_ = v836
	var v840 int32
	_ = v840
	var v843 int32
	_ = v843
	var v847 int32
	_ = v847
	var v850 int32
	_ = v850
	var v851 int32
	_ = v851
	var v853 int32
	_ = v853
	var v855 int32
	_ = v855
	var v858 int32
	_ = v858
	var v862 int32
	_ = v862
	var v870 int32
	_ = v870
	var v871 int32
	_ = v871
	var v872 int32
	_ = v872
	var v881 int32
	_ = v881
	var v882 int32
	_ = v882
	var v888 int32
	_ = v888
	var v895 int32
	_ = v895
	var v943 int32
	_ = v943
	var v990 int32
	_ = v990
	var v993 int32
	_ = v993
	var v995 int32
	_ = v995
	var v998 int32
	_ = v998
	var v1014 int32
	_ = v1014
	var v1017 int32
	_ = v1017
	var v1044 int32
	_ = v1044
	var v1046 int32
	_ = v1046
	var v1051 int32
	_ = v1051
	var v1113 int32
	_ = v1113
	var v1149 int32
	_ = v1149
	v45 = m.G0
	v47 = v45 - int32(2080)
	m.G0 = v47
	v49 = *(*int32)(unsafe.Add(mBase, uint32(l8)))
	goto L3
L1:
	;
	goto L17
L3:
	;
	base.MemoryFill(m, v47+int32(1040), int32(0), int32(1024))
	goto L1
L15:
	;
	v298 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+14)) = uint8(v298)
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+12)) = uint16(v298)
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+11)) = uint8(v298)
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+9)) = uint16(v298)
	v306 = int32(1)
	v308 = v306 << (uint(l2) % 32)
	v310 = v308 + int32(-1)
	v312 = int32(base.Ui32(v310+l1) >> (uint(l2) % 32))
	if v312 < v306 {
		goto L30
	} else {
		goto L31
	}
L17:
	;
	base.MemoryFill(m, v47+int32(16), int32(0), int32(1024))
	goto L15
L29:
	;
	m.G0 = v47 + int32(2080)
	return v1113
L30:
	;
	F_VP8LOptimizeSampling(m, l5, l0, l1, l2, int32(9), l9)
	mBase = m.M
	v1113 = v306
	goto L29
L31:
	;
	v316 = int32(base.Ui32(v310+l0) >> (uint(l2) % 32))
	if int32(0) < v316 {
		goto L33
	} else {
		goto L34
	}
L32:
	;
	v1014 = int32(0)
	v1017 = v312
	goto L91
L33:
	;
	v320 = int32(2)
	v323 = l0 << (uint(v320) % 32)
	v325 = l4 + int32(-8)
	v348 = int32(0)
	goto L35
L34:
	;
	goto L32
L35:
	;
	v371 = v348 << (uint(l2) % 32)
	v372 = v371 + v308
	if v372 < l1 {
		goto L37
	} else {
		goto L38
	}
L37:
	;
	v374 = v372
	goto L39
L38:
	;
	v374 = l1
	goto L39
L39:
	;
	v375 = l1 - v371
	if v308 < v375 {
		goto L40
	} else {
		goto L41
	}
L40:
	;
	v377 = v308
	goto L42
L41:
	;
	v377 = v375
	goto L42
L42:
	;
	v381 = v377 & int32(3)
	v387 = v371 * l0
	v429 = int32(0)
	goto L43
L43:
	;
	v440 = v429 + v348*v316
	if v348 == int32(0) {
		goto L46
	} else {
		goto L47
	}
L44:
	;
	v993 = base.I32_div_s(v348*l7, v312)
	v995 = F_WebPReportProgress(m, l6, v993+v49, l8)
	mBase = m.M
	if v995 != 0 {
		goto L87
	} else {
		goto L88
	}
L45:
	;
	v457 = v429 << (uint(l2) % 32)
	v460 = int32(2)
	v1149 = int32(8)
	v465 = v47 + int32(14)
	v466 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v465))))
	*(*uint8)(unsafe.Add(mBase, uint32(v47+v1149))) = uint8(v466)
	v476 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(11)))))
	*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(4)))) = uint8(v476)
	v478 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47)+12)))
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+6)) = uint16(v478)
	v480 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47)+9)))
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+2)) = uint16(v480)
	v490 = int32(16)
	F_GetBestColorTransformForTile(m, v47+int32(2072), v429, v456, l2, v47+int32(6), v47+v460, l3, l0, l1, v47+int32(1040), v47+v490, l4)
	mBase = m.M
	v497 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(2074)))))
	*(*uint8)(unsafe.Add(mBase, uint32(v465))) = uint8(v497)
	v499 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v47)+2072)))
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+12)) = uint16(v499)
	v506 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47)+13)))
	v510 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v47)+12)))
	*(*int32)(unsafe.Add(mBase, uint32(l5+v440<<(uint(v460)%32)))) = v497<<(uint(v490)%32) | v506<<(uint(v1149)%32) | v510 | int32(-16777216)
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+2074)) = uint8(v497)
	*(*uint16)(unsafe.Add(mBase, uint32(v47)+2072)) = uint16(v499)
	if v377 < int32(1) {
		goto L48
	} else {
		goto L49
	}
L46:
	;
	v456 = int32(0)
	goto L45
L47:
	;
	v447 = *(*int32)(unsafe.Add(mBase, uint32(l5+(v440-v316)<<(uint(int32(2))%32))))
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+9)) = uint8(v447)
	v450 = int32(base.Ui32(v447) >> (uint(int32(16)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+11)) = uint8(v450)
	v453 = int32(base.Ui32(v447) >> (uint(int32(8)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v47)+10)) = uint8(v453)
	v456 = v348
	goto L45
L48:
	;
	if v374 <= v371 {
		goto L63
	} else {
		goto L64
	}
L49:
	;
	v517 = l0 - v457
	if v308 < v517 {
		goto L50
	} else {
		goto L51
	}
L50:
	;
	v519 = v308
	goto L52
L51:
	;
	v519 = v517
	goto L52
L52:
	;
	v522 = l4 + v387<<(uint(int32(2))%32) + v457<<(uint(int32(2))%32)
	if v381 != 0 {
		goto L54
	} else {
		goto L55
	}
L53:
	;
	if base.Ui32(v377) < base.Ui32(int32(4)) {
		goto L48
	} else {
		goto L59
	}
L54:
	;
	v537 = v522
	v558 = v381
	goto L56
L55:
	;
	v589 = v522
	v610 = v377
	goto L53
L56:
	;
	v569 = m.G68
	v570 = *(*int32)(unsafe.Add(mBase, uint32(v569)))
	m.T0[v570].(func(*base.Module, int32, int32, int32))(m, v47+int32(2072), v537, v519)
	mBase = m.M
	v572 = v537 + v323
	v574 = v558 + int32(-1)
	if v574 != 0 {
		v537 = v572
		v558 = v574
		goto L56
	} else {
		goto L58
	}
L57:
	;
	v589 = v572
	v610 = v377 & int32(-4)
	goto L53
L58:
	;
	goto L57
L59:
	;
	v635 = v589
	v657 = v610 + int32(-1)
	goto L60
L60:
	;
	v666 = v47 + int32(2072)
	v667 = m.G68
	v668 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v668].(func(*base.Module, int32, int32, int32))(m, v666, v635, v519)
	mBase = m.M
	v672 = v635 + v323
	v673 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v673].(func(*base.Module, int32, int32, int32))(m, v666, v672, v519)
	mBase = m.M
	v677 = v672 + v323
	v678 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v678].(func(*base.Module, int32, int32, int32))(m, v666, v677, v519)
	mBase = m.M
	v682 = v677 + v323
	v683 = *(*int32)(unsafe.Add(mBase, uint32(v667)))
	m.T0[v683].(func(*base.Module, int32, int32, int32))(m, v666, v682, v519)
	mBase = m.M
	v687 = v657 + int32(-4)
	if base.Ui32(v687) < base.Ui32(int32(-2)) {
		v635 = v682 + v323
		v657 = v687
		goto L60
	} else {
		goto L62
	}
L61:
	;
	goto L48
L62:
	;
	goto L61
L63:
	;
	v990 = v429 + int32(1)
	if v990 != v316 {
		v429 = v990
		goto L43
	} else {
		goto L86
	}
L64:
	;
	v735 = v457 + v308
	if v735 < l0 {
		goto L65
	} else {
		goto L66
	}
L65:
	;
	v737 = v735
	goto L67
L66:
	;
	v737 = l0
	goto L67
L67:
	;
	v739 = int32(2)
	v785 = v325 + v323*(v371+int32(-1)) + v457<<(uint(v739)%32)
	v786 = v325 + (v387+v457)<<(uint(v739)%32)
	v787 = v371
	goto L68
L68:
	;
	v790 = v787 * l0
	v791 = v790 + v457
	if v790+v737 <= v791 {
		goto L70
	} else {
		goto L71
	}
L69:
	;
	goto L63
L70:
	;
	v943 = v787 + int32(1)
	if v943 != v374 {
		v785 = v785 + v323
		v786 = v786 + v323
		v787 = v943
		goto L68
	} else {
		goto L85
	}
L71:
	;
	v808 = v786
	v811 = v791
	v830 = v785
	v836 = v737 - v457
	goto L72
L72:
	;
	v840 = *(*int32)(unsafe.Add(mBase, uint32(v808+int32(8))))
	if v811 < int32(2) {
		goto L75
	} else {
		goto L76
	}
L73:
	;
	goto L70
L74:
	;
	v888 = int32(4)
	v895 = v836 + int32(-1)
	if v895 != 0 {
		v808 = v808 + v888
		v811 = v811 + int32(1)
		v830 = v830 + v888
		v836 = v895
		goto L72
	} else {
		goto L84
	}
L75:
	;
	if v811 < l0+v320 {
		goto L79
	} else {
		goto L80
	}
L76:
	;
	v843 = *(*int32)(unsafe.Add(mBase, uint32(v808)))
	if v840 != v843 {
		goto L75
	} else {
		goto L77
	}
L77:
	;
	v847 = *(*int32)(unsafe.Add(mBase, uint32(v808+int32(4))))
	if v840 == v847 {
		goto L74
	} else {
		goto L78
	}
L78:
	;
	goto L75
L79:
	;
	v870 = v47 + int32(1040) + int32(base.Ui32(v840)>>(uint(int32(14))%32))&int32(1020)
	v871 = *(*int32)(unsafe.Add(mBase, uint32(v870)))
	v872 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v870))) = v871 + v872
	v881 = v47 + int32(16) + v840&int32(255)<<(uint(int32(2))%32)
	v882 = *(*int32)(unsafe.Add(mBase, uint32(v881)))
	*(*int32)(unsafe.Add(mBase, uint32(v881))) = v882 + v872
	goto L74
L80:
	;
	v850 = *(*int32)(unsafe.Add(mBase, uint32(v808)))
	v851 = *(*int32)(unsafe.Add(mBase, uint32(v830)))
	if v850 != v851 {
		goto L79
	} else {
		goto L81
	}
L81:
	;
	v853 = int32(4)
	v855 = *(*int32)(unsafe.Add(mBase, uint32(v808+v853)))
	v858 = *(*int32)(unsafe.Add(mBase, uint32(v830+v853)))
	if v855 != v858 {
		goto L79
	} else {
		goto L82
	}
L82:
	;
	v862 = *(*int32)(unsafe.Add(mBase, uint32(v830+int32(8))))
	if v840 == v862 {
		goto L74
	} else {
		goto L83
	}
L83:
	;
	goto L79
L84:
	;
	goto L73
L85:
	;
	goto L69
L86:
	;
	goto L44
L87:
	;
	v998 = v348 + int32(1)
	if v998 != v312 {
		v348 = v998
		goto L35
	} else {
		goto L89
	}
L88:
	;
	v1113 = int32(0)
	goto L29
L89:
	;
	goto L30
L90:
	;
	v1113 = int32(0)
	goto L29
L91:
	;
	v1044 = base.I32_div_s(v1014, v312)
	v1046 = F_WebPReportProgress(m, l6, v1044+v49, l8)
	mBase = m.M
	if v1046 == int32(0) {
		goto L90
	} else {
		goto L93
	}
L93:
	;
	v1051 = v1017 + int32(-1)
	if v1051 == int32(0) {
		goto L30
	} else {
		goto L94
	}
L94:
	;
	v1014 = v1014 + l7
	v1017 = v1051
	goto L91
}
func F_VP8LCreateCompressedHuffmanTree(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v44 int32
	_ = v44
	var v53 int32
	_ = v53
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v74 int32
	_ = v74
	var v77 int32
	_ = v77
	var v86 int32
	_ = v86
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v106 int32
	_ = v106
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v176 int32
	_ = v176
	var v177 int32
	_ = v177
	var v182 int64
	_ = v182
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v198 int32
	_ = v198
	var v219 int32
	_ = v219
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v235 int32
	_ = v235
	var v241 int32
	_ = v241
	var v244 int32
	_ = v244
	var v255 int32
	_ = v255
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v269 int32
	_ = v269
	var v275 int32
	_ = v275
	var v280 int32
	_ = v280
	var v286 int32
	_ = v286
	var v290 int32
	_ = v290
	var v293 int32
	_ = v293
	var v295 int32
	_ = v295
	var v299 int32
	_ = v299
	var v306 int32
	_ = v306
	var v311 int32
	_ = v311
	var v325 int32
	_ = v325
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v361 int32
	_ = v361
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v381 int32
	_ = v381
	var v386 int32
	_ = v386
	var v392 int32
	_ = v392
	var v396 int32
	_ = v396
	var v398 int32
	_ = v398
	var v402 int32
	_ = v402
	var v404 int32
	_ = v404
	var v417 int32
	_ = v417
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v15 < int32(1) {
		v417 = l1
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return (v417 - l1) >> (uint(int32(1)) % 32)
L2:
	;
	v23 = l1
	v25 = int32(8)
	v26 = int32(0)
	goto L3
L3:
	;
	v35 = v26 + int32(1)
	if v35 < v15 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	v417 = v402
	goto L1
L5:
	;
	v37 = v15
	goto L7
L6:
	;
	v37 = v35
	goto L7
L7:
	;
	v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v40+v26))))
	v53 = v26
	goto L9
L8:
	;
	v68 = v65 - v26
	if v44 != 0 {
		goto L15
	} else {
		goto L16
	}
L9:
	;
	if v37+int32(-1) != v53 {
		goto L11
	} else {
		goto L12
	}
L10:
	;
	v65 = v62
	goto L8
L11:
	;
	v62 = v53 + int32(1)
	v63 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v40+int32(1)+v53))))
	if v63 == v44 {
		v53 = v62
		goto L9
	} else {
		goto L13
	}
L12:
	;
	v65 = v37
	goto L8
L13:
	;
	goto L10
L14:
	;
	if v65 < v15 {
		v23 = v402
		v25 = v404
		v26 = v65
		goto L3
	} else {
		goto L62
	}
L15:
	;
	if v25 == v44 {
		v262 = v23
		v263 = v68
		goto L41
	} else {
		goto L42
	}
L16:
	;
	if v68 < int32(1) {
		v402 = v23
		v404 = v25
		goto L14
	} else {
		goto L17
	}
L17:
	;
	v74 = base.I32_rem_u_s(v68+int32(-1), int32(138))
	v77 = v68 << (uint(int32(1)) % 32)
	v86 = v68
	v88 = v77
	v89 = v77
	v91 = int32(0)
	goto L18
L18:
	;
	v92 = v23 + v91
	if base.Ui32(int32(2)) < base.Ui32(v86) {
		goto L20
	} else {
		goto L21
	}
L19:
	;
	v402 = v23 + v244
	v404 = v25
	goto L14
L20:
	;
	if base.Ui32(int32(10)) < base.Ui32(v86) {
		goto L36
	} else {
		goto L37
	}
L21:
	;
	v95 = int32(0)
	if base.Ui32(v88) < base.Ui32(int32(33)) {
		goto L23
	} else {
		goto L24
	}
L22:
	;
	v402 = v23 + v89
	v404 = v25
	goto L14
L23:
	;
	if v88 == int32(0) {
		goto L25
	} else {
		goto L26
	}
L24:
	;
	base.MemoryFill(m, v92, v95, v88)
	goto L22
L25:
	;
	goto L22
L26:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v95)
	v106 = v92 + v88
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-1)))) = uint8(v95)
	if base.Ui32(v88) < base.Ui32(int32(3)) {
		goto L25
	} else {
		goto L27
	}
L27:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v92)+2)) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v92)+1)) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-3)))) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-2)))) = uint8(v95)
	if base.Ui32(v88) < base.Ui32(int32(7)) {
		goto L25
	} else {
		goto L28
	}
L28:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v92)+3)) = uint8(v95)
	*(*uint8)(unsafe.Add(mBase, uint32(v106+int32(-4)))) = uint8(v95)
	if base.Ui32(v88) < base.Ui32(int32(9)) {
		goto L25
	} else {
		goto L29
	}
L29:
	;
	v128 = int32(0)
	v131 = (v128 - v92) & int32(3)
	v132 = v92 + v131
	*(*int32)(unsafe.Add(mBase, uint32(v132))) = v128
	v140 = (v88 - v131) & int32(60)
	v141 = v132 + v140
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-4)))) = v128
	if base.Ui32(v140) < base.Ui32(int32(9)) {
		goto L25
	} else {
		goto L30
	}
L30:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v132)+8)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+4)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-8)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-12)))) = v128
	if base.Ui32(v140) < base.Ui32(int32(25)) {
		goto L25
	} else {
		goto L31
	}
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v132)+24)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+20)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+16)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v132)+12)) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-16)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-20)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-24)))) = v128
	*(*int32)(unsafe.Add(mBase, uint32(v141+int32(-28)))) = v128
	v176 = v132&int32(4) | int32(24)
	v177 = v140 - v176
	if base.Ui32(v177) < base.Ui32(int32(32)) {
		goto L25
	} else {
		goto L32
	}
L32:
	;
	v182 = base.I64_extend_i32_u(v128) * int64(4294967297)
	v185 = v177
	v186 = v132 + v176
	goto L33
L33:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v186)+24)) = v182
	*(*int64)(unsafe.Add(mBase, uint32(v186)+16)) = v182
	*(*int64)(unsafe.Add(mBase, uint32(v186)+8)) = v182
	*(*int64)(unsafe.Add(mBase, uint32(v186))) = v182
	v198 = v185 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v198) {
		v185 = v198
		v186 = v186 + int32(32)
		goto L33
	} else {
		goto L35
	}
L34:
	;
	goto L25
L35:
	;
	goto L34
L36:
	;
	v228 = int32(18)
	*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v228)
	if base.Ui32(int32(138)) < base.Ui32(v86) {
		goto L38
	} else {
		goto L39
	}
L37:
	;
	v219 = int32(17)
	*(*uint8)(unsafe.Add(mBase, uint32(v92))) = uint8(v219)
	v224 = v86 + int32(-3)
	*(*uint8)(unsafe.Add(mBase, uint32(v92+int32(1)))) = uint8(v224)
	v402 = v92 + int32(2)
	v404 = v25
	goto L14
L38:
	;
	v241 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v92+int32(1)))) = uint8(v241)
	v244 = v91 + int32(2)
	if int32(138) < v86 {
		v86 = v86 + int32(-138)
		v88 = v88 + int32(-276)
		v89 = v89 + int32(-274)
		v91 = v244
		goto L18
	} else {
		goto L40
	}
L39:
	;
	v235 = v74 + int32(-10)
	*(*uint8)(unsafe.Add(mBase, uint32(v92+int32(1)))) = uint8(v235)
	v402 = v92 + int32(2)
	v404 = v25
	goto L14
L40:
	;
	goto L19
L41:
	;
	if int32(1) <= v263 {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v255 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+1)) = uint8(v255)
	*(*uint8)(unsafe.Add(mBase, uint32(v23))) = uint8(v44)
	v262 = v23 + int32(2)
	v263 = v68 + int32(-1)
	goto L41
L43:
	;
	v269 = base.I32_rem_u_s(v263+int32(-1), int32(6))
	if base.Ui32(v263) < base.Ui32(int32(3)) {
		v306 = v262
		v311 = v263
		goto L45
	} else {
		goto L46
	}
L44:
	;
	v402 = v262
	v404 = v44
	goto L14
L45:
	;
	if base.Ui32(v311+int32(-1)) < base.Ui32(int32(7)) {
		v367 = v306
		goto L52
	} else {
		goto L53
	}
L46:
	;
	v275 = v262
	v280 = v263
	goto L48
L47:
	;
	v299 = v269 + int32(-2)
	*(*uint8)(unsafe.Add(mBase, uint32(v275)+1)) = uint8(v299)
	v402 = v275 + int32(2)
	v404 = v44
	goto L14
L48:
	;
	v286 = int32(16)
	*(*uint8)(unsafe.Add(mBase, uint32(v275))) = uint8(v286)
	if base.Ui32(v280) <= base.Ui32(int32(6)) {
		goto L47
	} else {
		goto L50
	}
L50:
	;
	v290 = int32(3)
	*(*uint8)(unsafe.Add(mBase, uint32(v275)+1)) = uint8(v290)
	v293 = v275 + int32(2)
	v295 = v280 + int32(-6)
	if base.Ui32(v295) < base.Ui32(v290) {
		v306 = v293
		v311 = v295
		goto L45
	} else {
		goto L51
	}
L51:
	;
	v275 = v293
	v280 = v295
	goto L48
L52:
	;
	if v311 != 0 {
		goto L57
	} else {
		goto L58
	}
L53:
	;
	v325 = v306
	v335 = int32(0)
	goto L54
L54:
	;
	v336 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+15)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+14)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+13)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+12)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+11)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+10)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+9)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+8)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+7)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+6)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+5)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+4)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+3)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+2)) = uint8(v44)
	*(*uint8)(unsafe.Add(mBase, uint32(v325)+1)) = uint8(v336)
	*(*uint8)(unsafe.Add(mBase, uint32(v325))) = uint8(v44)
	v361 = v325 + int32(16)
	v363 = v335 + int32(-8)
	if v363 != 0 {
		v325 = v361
		v335 = v363
		goto L54
	} else {
		goto L56
	}
L55:
	;
	v367 = v361
	goto L52
L56:
	;
	goto L55
L57:
	;
	v381 = v367
	v386 = v311
	goto L59
L58:
	;
	v402 = v367
	v404 = v44
	goto L14
L59:
	;
	v392 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v381)+1)) = uint8(v392)
	*(*uint8)(unsafe.Add(mBase, uint32(v381))) = uint8(v44)
	v396 = v381 + int32(2)
	v398 = v386 + int32(-1)
	if v398 != 0 {
		v381 = v396
		v386 = v398
		goto L59
	} else {
		goto L61
	}
L60:
	;
	v402 = v396
	v404 = v44
	goto L14
L61:
	;
	goto L60
L62:
	;
	goto L4
}
func F_VP8LCreateHuffmanTree(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v31 int32
	_ = v31
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v47 int32
	_ = v47
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v123 int64
	_ = v123
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v139 int32
	_ = v139
	var v161 int32
	_ = v161
	var v170 int32
	_ = v170
	var v177 int32
	_ = v177
	var v198 int32
	_ = v198
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v208 int32
	_ = v208
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v246 int32
	_ = v246
	var v247 int32
	_ = v247
	var v258 int32
	_ = v258
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v287 int32
	_ = v287
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v328 int32
	_ = v328
	var v329 int32
	_ = v329
	var v334 int64
	_ = v334
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v350 int32
	_ = v350
	var v368 int32
	_ = v368
	var v372 int32
	_ = v372
	var v373 int32
	_ = v373
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v386 int32
	_ = v386
	var v392 int32
	_ = v392
	var v395 int32
	_ = v395
	var v399 int32
	_ = v399
	var v400 int32
	_ = v400
	var v414 int32
	_ = v414
	var v415 int32
	_ = v415
	var v416 int32
	_ = v416
	var v421 int32
	_ = v421
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v428 int32
	_ = v428
	var v440 int32
	_ = v440
	var v443 int32
	_ = v443
	var v444 int32
	_ = v444
	var v446 int32
	_ = v446
	var v448 int32
	_ = v448
	var v453 int32
	_ = v453
	var v467 int32
	_ = v467
	var v468 int32
	_ = v468
	var v490 int32
	_ = v490
	var v509 int32
	_ = v509
	var v522 int32
	_ = v522
	var v542 int32
	_ = v542
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v582 int32
	_ = v582
	var v614 int32
	_ = v614
	var v616 int32
	_ = v616
	var v617 int32
	_ = v617
	var v618 int32
	_ = v618
	var v620 int32
	_ = v620
	var v622 int32
	_ = v622
	var v633 int32
	_ = v633
	var v635 int32
	_ = v635
	var v639 int32
	_ = v639
	var v641 int32
	_ = v641
	var v643 int32
	_ = v643
	var v649 int32
	_ = v649
	var v652 int32
	_ = v652
	var v661 int32
	_ = v661
	var v665 int32
	_ = v665
	var v719 int32
	_ = v719
	var v723 int32
	_ = v723
	var v726 int32
	_ = v726
	var v730 int32
	_ = v730
	var v734 int32
	_ = v734
	var v744 int32
	_ = v744
	var v751 int32
	_ = v751
	var v762 int32
	_ = v762
	var v763 int32
	_ = v763
	var v766 int32
	_ = v766
	var v768 int32
	_ = v768
	var v774 int32
	_ = v774
	var v780 int32
	_ = v780
	var v783 int32
	_ = v783
	var v787 int32
	_ = v787
	var v801 int32
	_ = v801
	var v808 int32
	_ = v808
	var v826 int32
	_ = v826
	var v835 int32
	_ = v835
	var v843 int32
	_ = v843
	var v854 int32
	_ = v854
	var v857 int32
	_ = v857
	var v861 int32
	_ = v861
	var v881 int32
	_ = v881
	var v894 int32
	_ = v894
	var v896 int32
	_ = v896
	var v901 int32
	_ = v901
	var v902 int32
	_ = v902
	var v905 int32
	_ = v905
	var v922 int32
	_ = v922
	var v942 int32
	_ = v942
	var v944 int32
	_ = v944
	var v948 int32
	_ = v948
	var v957 int32
	_ = v957
	var v958 int32
	_ = v958
	var v976 int32
	_ = v976
	var v981 int32
	_ = v981
	var v986 int32
	_ = v986
	var v990 int32
	_ = v990
	var v994 int32
	_ = v994
	var v999 int32
	_ = v999
	var v1006 int32
	_ = v1006
	var v1010 int32
	_ = v1010
	var v1015 int32
	_ = v1015
	var v1028 int32
	_ = v1028
	var v1029 int32
	_ = v1029
	var v1052 int32
	_ = v1052
	var v1057 int32
	_ = v1057
	var v1062 int32
	_ = v1062
	var v1067 int32
	_ = v1067
	var v1074 int32
	_ = v1074
	var v1076 int32
	_ = v1076
	var v1078 int32
	_ = v1078
	var v1085 int32
	_ = v1085
	var v1095 int32
	_ = v1095
	var v1113 int32
	_ = v1113
	var v1116 int32
	_ = v1116
	var v1118 int32
	_ = v1118
	var v1121 int32
	_ = v1121
	var v1123 int32
	_ = v1123
	var v1126 int32
	_ = v1126
	var v1128 int32
	_ = v1128
	var v1131 int32
	_ = v1131
	var v1133 int32
	_ = v1133
	var v1135 int32
	_ = v1135
	var v1141 int32
	_ = v1141
	var v1150 int32
	_ = v1150
	var v1174 int32
	_ = v1174
	var v1183 int32
	_ = v1183
	var v1184 int32
	_ = v1184
	var v1202 int32
	_ = v1202
	var v1204 int32
	_ = v1204
	var v1208 int32
	_ = v1208
	var v1211 int32
	_ = v1211
	var v1246 int32
	_ = v1246
	var v1249 int32
	_ = v1249
	var v1250 int32
	_ = v1250
	var v1253 int32
	_ = v1253
	var v1255 int32
	_ = v1255
	var v1269 int32
	_ = v1269
	var v1292 int32
	_ = v1292
	var v1296 int32
	_ = v1296
	var v1300 int32
	_ = v1300
	var v1309 int32
	_ = v1309
	var v1310 int32
	_ = v1310
	var v1328 int32
	_ = v1328
	var v1333 int32
	_ = v1333
	var v1338 int32
	_ = v1338
	var v1342 int32
	_ = v1342
	var v1346 int32
	_ = v1346
	var v1351 int32
	_ = v1351
	var v1358 int32
	_ = v1358
	var v1362 int32
	_ = v1362
	var v1367 int32
	_ = v1367
	var v1380 int32
	_ = v1380
	var v1381 int32
	_ = v1381
	var v1404 int32
	_ = v1404
	var v1409 int32
	_ = v1409
	var v1414 int32
	_ = v1414
	var v1419 int32
	_ = v1419
	var v1427 int32
	_ = v1427
	var v1434 int32
	_ = v1434
	var v1437 int32
	_ = v1437
	var v1455 int32
	_ = v1455
	var v1457 int32
	_ = v1457
	var v1460 int32
	_ = v1460
	var v1461 int64
	_ = v1461
	var v1463 int32
	_ = v1463
	var v1467 int64
	_ = v1467
	var v1470 int32
	_ = v1470
	var v1473 int32
	_ = v1473
	var v1475 int32
	_ = v1475
	var v1478 int32
	_ = v1478
	var v1479 int64
	_ = v1479
	var v1485 int64
	_ = v1485
	var v1489 int32
	_ = v1489
	var v1496 int32
	_ = v1496
	var v1505 int32
	_ = v1505
	var v1524 int32
	_ = v1524
	var v1529 int32
	_ = v1529
	var v1542 int32
	_ = v1542
	var v1563 int32
	_ = v1563
	var v1565 int32
	_ = v1565
	var v1567 int32
	_ = v1567
	var v1570 int32
	_ = v1570
	var v1578 int32
	_ = v1578
	var v1589 int32
	_ = v1589
	var v1597 int32
	_ = v1597
	var v1600 int32
	_ = v1600
	var v1602 int32
	_ = v1602
	var v1609 int32
	_ = v1609
	var v1612 int32
	_ = v1612
	var v1614 int32
	_ = v1614
	var v1621 int32
	_ = v1621
	var v1624 int32
	_ = v1624
	var v1626 int32
	_ = v1626
	var v1633 int32
	_ = v1633
	var v1648 int32
	_ = v1648
	var v1649 int32
	_ = v1649
	var v1651 int32
	_ = v1651
	var v1658 int32
	_ = v1658
	var v1659 int32
	_ = v1659
	var v1661 int32
	_ = v1661
	var v1668 int32
	_ = v1668
	var v1669 int32
	_ = v1669
	var v1671 int32
	_ = v1671
	var v1678 int32
	_ = v1678
	var v1681 int32
	_ = v1681
	var v1683 int32
	_ = v1683
	var v1689 int32
	_ = v1689
	var v1695 int32
	_ = v1695
	var v1704 int32
	_ = v1704
	var v1705 int32
	_ = v1705
	var v1711 int32
	_ = v1711
	var v1714 int32
	_ = v1714
	var v1716 int32
	_ = v1716
	var v1719 int32
	_ = v1719
	var v1732 int32
	_ = v1732
	var v1737 int32
	_ = v1737
	var v1738 int32
	_ = v1738
	var v1740 int32
	_ = v1740
	var v1743 int32
	_ = v1743
	var v1745 int32
	_ = v1745
	var v1749 int32
	_ = v1749
	var v1751 int32
	_ = v1751
	var v1755 int32
	_ = v1755
	var v1757 int32
	_ = v1757
	var v1760 int32
	_ = v1760
	var v1765 int32
	_ = v1765
	var v1773 int32
	_ = v1773
	var v1783 int32
	_ = v1783
	var v1784 int32
	_ = v1784
	var v1789 int32
	_ = v1789
	var v1791 int32
	_ = v1791
	var v1792 int32
	_ = v1792
	var v1794 int32
	_ = v1794
	var v1798 int32
	_ = v1798
	var v1811 int32
	_ = v1811
	var v1815 int32
	_ = v1815
	var v1816 int32
	_ = v1816
	var v1818 int32
	_ = v1818
	var v1821 int32
	_ = v1821
	var v1823 int32
	_ = v1823
	var v1827 int32
	_ = v1827
	var v1829 int32
	_ = v1829
	var v1833 int32
	_ = v1833
	var v1835 int32
	_ = v1835
	var v1838 int32
	_ = v1838
	var v1839 int32
	_ = v1839
	var v1840 int32
	_ = v1840
	var v1841 int32
	_ = v1841
	var v1845 int32
	_ = v1845
	var v1851 int32
	_ = v1851
	var v1858 int32
	_ = v1858
	var v1859 int32
	_ = v1859
	var v1860 int32
	_ = v1860
	var v1864 int32
	_ = v1864
	var v1866 int32
	_ = v1866
	var v1867 int32
	_ = v1867
	var v1869 int32
	_ = v1869
	var v1871 int32
	_ = v1871
	var v1873 int32
	_ = v1873
	var v1875 int32
	_ = v1875
	var v1877 int32
	_ = v1877
	var v1882 int32
	_ = v1882
	var v1884 int32
	_ = v1884
	var v1886 int32
	_ = v1886
	var v1888 int32
	_ = v1888
	var v1890 int32
	_ = v1890
	var v1892 int32
	_ = v1892
	var v1894 int32
	_ = v1894
	var v1896 int32
	_ = v1896
	var v1898 int32
	_ = v1898
	var v1900 int32
	_ = v1900
	var v1902 int32
	_ = v1902
	var v1904 int32
	_ = v1904
	var v1905 int32
	_ = v1905
	var v1907 int32
	_ = v1907
	var v1909 int32
	_ = v1909
	var v1913 int32
	_ = v1913
	var v1915 int32
	_ = v1915
	var v1917 int32
	_ = v1917
	var v1922 int32
	_ = v1922
	var v1926 int32
	_ = v1926
	var v1927 int32
	_ = v1927
	var v1928 int32
	_ = v1928
	var v1932 int32
	_ = v1932
	var v1934 int32
	_ = v1934
	var v1935 int32
	_ = v1935
	var v1937 int32
	_ = v1937
	var v1939 int32
	_ = v1939
	var v1941 int32
	_ = v1941
	var v1943 int32
	_ = v1943
	var v1944 int32
	_ = v1944
	var v1950 int32
	_ = v1950
	var v1952 int32
	_ = v1952
	var v1953 int32
	_ = v1953
	var v1956 int32
	_ = v1956
	var v1958 int32
	_ = v1958
	var v1960 int32
	_ = v1960
	var v1962 int32
	_ = v1962
	var v1964 int32
	_ = v1964
	var v1966 int32
	_ = v1966
	var v1968 int32
	_ = v1968
	var v1970 int32
	_ = v1970
	var v1972 int32
	_ = v1972
	var v1977 int32
	_ = v1977
	var v1987 int32
	_ = v1987
	var v1996 int32
	_ = v1996
	var v1998 int32
	_ = v1998
	var v2001 int32
	_ = v2001
	var v2004 int32
	_ = v2004
	var v2005 int32
	_ = v2005
	var v2006 int32
	_ = v2006
	var v2010 int32
	_ = v2010
	var v2012 int32
	_ = v2012
	var v2015 int32
	_ = v2015
	var v2016 int32
	_ = v2016
	var v2019 int32
	_ = v2019
	var v2022 int32
	_ = v2022
	var v2024 int32
	_ = v2024
	var v2027 int32
	_ = v2027
	var v2034 int32
	_ = v2034
	var v2044 int32
	_ = v2044
	var v2062 int32
	_ = v2062
	var v2065 int32
	_ = v2065
	var v2067 int32
	_ = v2067
	var v2070 int32
	_ = v2070
	var v2072 int32
	_ = v2072
	var v2075 int32
	_ = v2075
	var v2077 int32
	_ = v2077
	var v2080 int32
	_ = v2080
	var v2082 int32
	_ = v2082
	var v2084 int32
	_ = v2084
	var v2090 int32
	_ = v2090
	var v2099 int32
	_ = v2099
	var v2123 int32
	_ = v2123
	var v2132 int32
	_ = v2132
	var v2133 int32
	_ = v2133
	var v2151 int32
	_ = v2151
	var v2153 int32
	_ = v2153
	var v2157 int32
	_ = v2157
	var v2160 int32
	_ = v2160
	var v2209 int32
	_ = v2209
	var v2223 int64
	_ = v2223
	var v2249 int32
	_ = v2249
	var v2251 int32
	_ = v2251
	var v2254 int32
	_ = v2254
	var v2269 int32
	_ = v2269
	var v2270 int32
	_ = v2270
	var v2288 int32
	_ = v2288
	var v2307 int32
	_ = v2307
	var v2308 int32
	_ = v2308
	var v2309 int32
	_ = v2309
	var v2311 int32
	_ = v2311
	var v2312 int32
	_ = v2312
	var v2313 int32
	_ = v2313
	var v2318 int32
	_ = v2318
	var v2321 int32
	_ = v2321
	var v2322 int32
	_ = v2322
	var v2328 int32
	_ = v2328
	var v2331 int32
	_ = v2331
	var v2332 int32
	_ = v2332
	var v2338 int32
	_ = v2338
	var v2341 int32
	_ = v2341
	var v2342 int32
	_ = v2342
	var v2347 int32
	_ = v2347
	var v2360 int32
	_ = v2360
	var v2384 int32
	_ = v2384
	var v2394 int32
	_ = v2394
	var v2412 int32
	_ = v2412
	var v2415 int32
	_ = v2415
	var v2416 int32
	_ = v2416
	var v2417 int32
	_ = v2417
	var v2423 int32
	_ = v2423
	var v2454 int32
	_ = v2454
	var v2457 int32
	_ = v2457
	var v2458 int32
	_ = v2458
	var v2459 int32
	_ = v2459
	var v2460 int32
	_ = v2460
	var v2461 int32
	_ = v2461
	var v2462 int32
	_ = v2462
	var v2463 int32
	_ = v2463
	var v2464 int32
	_ = v2464
	var v2465 int32
	_ = v2465
	var v2466 int32
	_ = v2466
	var v2467 int32
	_ = v2467
	var v2468 int32
	_ = v2468
	var v2469 int32
	_ = v2469
	var v2472 int32
	_ = v2472
	var v2473 int32
	_ = v2473
	var v2477 int32
	_ = v2477
	var v2478 int32
	_ = v2478
	var v2480 int32
	_ = v2480
	var v2481 int32
	_ = v2481
	var v2482 int32
	_ = v2482
	var v2483 int32
	_ = v2483
	var v2484 int32
	_ = v2484
	var v2485 int32
	_ = v2485
	var v2487 int32
	_ = v2487
	var v2489 int32
	_ = v2489
	var v2490 int32
	_ = v2490
	var v2494 int32
	_ = v2494
	var v2504 int32
	_ = v2504
	var v2505 int32
	_ = v2505
	var v2509 int32
	_ = v2509
	var v2513 int32
	_ = v2513
	var v2517 int32
	_ = v2517
	var v2521 int32
	_ = v2521
	var v2525 int32
	_ = v2525
	var v2529 int32
	_ = v2529
	var v2533 int32
	_ = v2533
	var v2537 int32
	_ = v2537
	var v2541 int32
	_ = v2541
	var v2545 int32
	_ = v2545
	var v2549 int32
	_ = v2549
	var v2555 int32
	_ = v2555
	var v2556 int32
	_ = v2556
	var v2574 int32
	_ = v2574
	var v2591 int32
	_ = v2591
	var v2594 int32
	_ = v2594
	var v2595 int32
	_ = v2595
	var v2599 int32
	_ = v2599
	var v2603 int32
	_ = v2603
	var v2614 int32
	_ = v2614
	var v2618 int32
	_ = v2618
	var v2626 int32
	_ = v2626
	var v2627 int32
	_ = v2627
	var v2628 int32
	_ = v2628
	var v2646 int32
	_ = v2646
	var v2648 int32
	_ = v2648
	var v2651 int32
	_ = v2651
	var v2654 int32
	_ = v2654
	var v2655 int32
	_ = v2655
	var v2661 int32
	_ = v2661
	var v2666 int32
	_ = v2666
	var v2668 int32
	_ = v2668
	var v2670 int32
	_ = v2670
	var v2672 int32
	_ = v2672
	var v2677 int32
	_ = v2677
	var v2686 int32
	_ = v2686
	var v2687 int32
	_ = v2687
	var v2707 int32
	_ = v2707
	var v2713 int32
	_ = v2713
	var v2730 int32
	_ = v2730
	var v2748 int32
	_ = v2748
	var v2753 int32
	_ = v2753
	var v2756 int32
	_ = v2756
	v6 = int32(0)
	v31 = m.G0
	v33 = v31 - int32(128)
	m.G0 = v33
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	if base.Ui32(v36) < base.Ui32(int32(33)) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	if v36 < int32(0) {
		goto L17
	} else {
		goto L18
	}
L2:
	;
	if v36 == int32(0) {
		goto L4
	} else {
		goto L5
	}
L3:
	;
	base.MemoryFill(m, l2, v6, v36)
	goto L1
L4:
	;
	goto L1
L5:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l2))) = uint8(v6)
	v47 = l2 + v36
	*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(-1)))) = uint8(v6)
	if base.Ui32(v36) < base.Ui32(int32(3)) {
		goto L4
	} else {
		goto L6
	}
L6:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+2)) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+1)) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(-3)))) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(-2)))) = uint8(v6)
	if base.Ui32(v36) < base.Ui32(int32(7)) {
		goto L4
	} else {
		goto L7
	}
L7:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l2)+3)) = uint8(v6)
	*(*uint8)(unsafe.Add(mBase, uint32(v47+int32(-4)))) = uint8(v6)
	if base.Ui32(v36) < base.Ui32(int32(9)) {
		goto L4
	} else {
		goto L8
	}
L8:
	;
	v69 = int32(0)
	v72 = (v69 - l2) & int32(3)
	v73 = l2 + v72
	*(*int32)(unsafe.Add(mBase, uint32(v73))) = v69
	v81 = (v36 - v72) & int32(60)
	v82 = v73 + v81
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-4)))) = v69
	if base.Ui32(v81) < base.Ui32(int32(9)) {
		goto L4
	} else {
		goto L9
	}
L9:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v73)+8)) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v73)+4)) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-8)))) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-12)))) = v69
	if base.Ui32(v81) < base.Ui32(int32(25)) {
		goto L4
	} else {
		goto L10
	}
L10:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v73)+24)) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v73)+20)) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v73)+16)) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v73)+12)) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-16)))) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-20)))) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-24)))) = v69
	*(*int32)(unsafe.Add(mBase, uint32(v82+int32(-28)))) = v69
	v117 = v73&int32(4) | int32(24)
	v118 = v81 - v117
	if base.Ui32(v118) < base.Ui32(int32(32)) {
		goto L4
	} else {
		goto L11
	}
L11:
	;
	v123 = base.I64_extend_i32_u(v69) * int64(4294967297)
	v126 = v118
	v127 = v73 + v117
	goto L12
L12:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v127)+24)) = v123
	*(*int64)(unsafe.Add(mBase, uint32(v127)+16)) = v123
	*(*int64)(unsafe.Add(mBase, uint32(v127)+8)) = v123
	*(*int64)(unsafe.Add(mBase, uint32(v127))) = v123
	v139 = v126 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v139) {
		v126 = v139
		v127 = v127 + int32(32)
		goto L12
	} else {
		goto L14
	}
L13:
	;
	goto L4
L14:
	;
	goto L13
L15:
	;
	v2223 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v33+int32(56)))) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33+int32(48)))) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33+int32(40)))) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33+int32(32)))) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33+int32(24)))) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33+int32(16)))) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33)+8)) = v2223
	*(*int64)(unsafe.Add(mBase, uint32(v33))) = v2223
	v2249 = *(*int32)(unsafe.Add(mBase, uint32(l4)))
	v2251 = base.B2i32(v2249 < int32(1))
	if v2251 == int32(0) {
		goto L309
	} else {
		goto L310
	}
L16:
	;
	v719 = l4 + int32(4)
	if v36 == int32(0) {
		v2209 = v719
		goto L15
	} else {
		goto L93
	}
L17:
	;
	v2209 = l4 + int32(4)
	goto L15
L18:
	;
	if v36 == int32(0) {
		goto L19
	} else {
		goto L20
	}
L19:
	;
	v2209 = l4 + int32(4)
	goto L15
L20:
	;
	v161 = int32(-4)
	v170 = v36<<(uint(int32(2))%32) + l0 + v161
	v177 = v36
	goto L22
L21:
	;
	v203 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v204 = int32(0)
	v208 = v204
	v216 = v203
	v217 = v204
	v218 = l0
	goto L27
L22:
	;
	v198 = *(*int32)(unsafe.Add(mBase, uint32(v170)))
	if v198 != 0 {
		goto L21
	} else {
		goto L24
	}
L24:
	;
	v202 = v177 + int32(-1)
	if v202 != 0 {
		v170 = v170 + int32(-4)
		v177 = v202
		goto L22
	} else {
		goto L25
	}
L25:
	;
	goto L16
L26:
	;
	v380 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v381 = int32(0)
	v386 = v381
	v392 = l0 + v161
	v395 = v380
	v399 = v381
	v400 = v381
	goto L56
L27:
	;
	v236 = base.B2i32(v177 == v208)
	if v177 == v208 {
		goto L31
	} else {
		goto L32
	}
L29:
	;
	v208 = v208 + int32(1)
	v216 = v372
	v217 = v373
	v218 = v218 + int32(4)
	goto L27
L30:
	;
	v372 = v216
	v373 = v217 + int32(1)
	goto L29
L31:
	;
	if v216 != 0 {
		goto L36
	} else {
		goto L37
	}
L32:
	;
	v237 = *(*int32)(unsafe.Add(mBase, uint32(v218)))
	if v237 == v216 {
		goto L30
	} else {
		goto L33
	}
L33:
	;
	goto L31
L34:
	;
	if v177 == v208 {
		goto L26
	} else {
		goto L55
	}
L35:
	;
	v246 = l2 + (v208 - v217)
	v247 = int32(1)
	if base.Ui32(v217) < base.Ui32(int32(33)) {
		goto L42
	} else {
		goto L43
	}
L36:
	;
	if v216 == int32(0) {
		goto L34
	} else {
		goto L39
	}
L37:
	;
	if int32(4) < v217 {
		goto L35
	} else {
		goto L38
	}
L38:
	;
	goto L36
L39:
	;
	if v217 < int32(7) {
		goto L34
	} else {
		goto L40
	}
L40:
	;
	goto L35
L41:
	;
	goto L34
L42:
	;
	if v217 == int32(0) {
		goto L44
	} else {
		goto L45
	}
L43:
	;
	base.MemoryFill(m, v246, v247, v217)
	goto L41
L44:
	;
	goto L41
L45:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v246))) = uint8(v247)
	v258 = v246 + v217
	*(*uint8)(unsafe.Add(mBase, uint32(v258+int32(-1)))) = uint8(v247)
	if base.Ui32(v217) < base.Ui32(int32(3)) {
		goto L44
	} else {
		goto L46
	}
L46:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v246)+2)) = uint8(v247)
	*(*uint8)(unsafe.Add(mBase, uint32(v246)+1)) = uint8(v247)
	*(*uint8)(unsafe.Add(mBase, uint32(v258+int32(-3)))) = uint8(v247)
	*(*uint8)(unsafe.Add(mBase, uint32(v258+int32(-2)))) = uint8(v247)
	if base.Ui32(v217) < base.Ui32(int32(7)) {
		goto L44
	} else {
		goto L47
	}
L47:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v246)+3)) = uint8(v247)
	*(*uint8)(unsafe.Add(mBase, uint32(v258+int32(-4)))) = uint8(v247)
	if base.Ui32(v217) < base.Ui32(int32(9)) {
		goto L44
	} else {
		goto L48
	}
L48:
	;
	v283 = (int32(0) - v246) & int32(3)
	v284 = v246 + v283
	v287 = int32(16843009)
	*(*int32)(unsafe.Add(mBase, uint32(v284))) = v287
	v292 = (v217 - v283) & int32(60)
	v293 = v284 + v292
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-4)))) = v287
	if base.Ui32(v292) < base.Ui32(int32(9)) {
		goto L44
	} else {
		goto L49
	}
L49:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v284)+8)) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v284)+4)) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-8)))) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-12)))) = v287
	if base.Ui32(v292) < base.Ui32(int32(25)) {
		goto L44
	} else {
		goto L50
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v284)+24)) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v284)+20)) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v284)+16)) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v284)+12)) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-16)))) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-20)))) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-24)))) = v287
	*(*int32)(unsafe.Add(mBase, uint32(v293+int32(-28)))) = v287
	v328 = v284&int32(4) | int32(24)
	v329 = v292 - v328
	if base.Ui32(v329) < base.Ui32(int32(32)) {
		goto L44
	} else {
		goto L51
	}
L51:
	;
	v334 = base.I64_extend_i32_u(v287) * int64(4294967297)
	v337 = v329
	v338 = v284 + v328
	goto L52
L52:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v338)+24)) = v334
	*(*int64)(unsafe.Add(mBase, uint32(v338)+16)) = v334
	*(*int64)(unsafe.Add(mBase, uint32(v338)+8)) = v334
	*(*int64)(unsafe.Add(mBase, uint32(v338))) = v334
	v350 = v337 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v350) {
		v337 = v350
		v338 = v338 + int32(32)
		goto L52
	} else {
		goto L54
	}
L53:
	;
	goto L44
L54:
	;
	goto L53
L55:
	;
	v368 = *(*int32)(unsafe.Add(mBase, uint32(v218)))
	v372 = v368
	v373 = int32(1)
	goto L29
L56:
	;
	v414 = base.B2i32(v177 == v400)
	if v177 == v400 {
		goto L60
	} else {
		goto L61
	}
L58:
	;
	v386 = v652
	v392 = v392 + int32(4)
	v395 = v661
	v399 = v665
	v400 = v400 + int32(1)
	goto L56
L59:
	;
	v641 = v425 + v386
	v643 = v399 + int32(1)
	if base.Ui32(v643) < base.Ui32(int32(4)) {
		v652 = v641
		v661 = v395
		v665 = v643
		goto L58
	} else {
		goto L92
	}
L60:
	;
	if base.Ui32(int32(3)) < base.Ui32(v399) {
		goto L68
	} else {
		goto L69
	}
L61:
	;
	v415 = l2 + v400
	v416 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v415))))
	if v416 != 0 {
		goto L60
	} else {
		goto L62
	}
L62:
	;
	if v400 == int32(0) {
		goto L63
	} else {
		goto L64
	}
L63:
	;
	v425 = *(*int32)(unsafe.Add(mBase, uint32(l0+v400<<(uint(int32(2))%32))))
	v426 = v425 - v395
	v428 = v426 >> (uint(int32(31)) % 32)
	if base.Ui32(v426^v428-v428) < base.Ui32(int32(4)) {
		goto L59
	} else {
		goto L66
	}
L64:
	;
	v421 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v415+int32(-1)))))
	if v421 != 0 {
		goto L60
	} else {
		goto L65
	}
L65:
	;
	goto L63
L66:
	;
	goto L60
L67:
	;
	if v177+int32(-3) <= v400 {
		goto L88
	} else {
		goto L89
	}
L68:
	;
	v440 = int32(1)
	v443 = int32(base.Ui32(v399)>>(uint(v440)%32)) + v386
	v444 = base.I32_div_u_s(v443, v399)
	if base.Ui32(v443) < base.Ui32(v399) {
		goto L72
	} else {
		goto L73
	}
L69:
	;
	if v399 != int32(3) {
		goto L67
	} else {
		goto L70
	}
L70:
	;
	if v386 != 0 {
		goto L67
	} else {
		goto L71
	}
L71:
	;
	goto L68
L72:
	;
	v446 = v440
	goto L74
L73:
	;
	v446 = v444
	goto L74
L74:
	;
	if v386 != 0 {
		goto L75
	} else {
		goto L76
	}
L75:
	;
	v448 = v446
	goto L77
L76:
	;
	v448 = int32(0)
	goto L77
L77:
	;
	if base.Ui32(v399) < base.Ui32(int32(4)) {
		v522 = int32(0)
		goto L78
	} else {
		goto L79
	}
L78:
	;
	v542 = v399 & int32(3)
	if v542 == int32(0) {
		goto L67
	} else {
		goto L83
	}
L79:
	;
	v453 = l0 + v400<<(uint(int32(2))%32)
	v467 = v392
	v468 = int32(0)
	goto L80
L80:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v467))) = v448
	v490 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v453+(v468^int32(1073741821))<<(uint(v490)%32)))) = v448
	*(*int32)(unsafe.Add(mBase, uint32(v453+(v468^int32(1073741822))<<(uint(v490)%32)))) = v448
	*(*int32)(unsafe.Add(mBase, uint32(v453+(v468^int32(1073741820))<<(uint(v490)%32)))) = v448
	v509 = v468 + int32(4)
	if v399&int32(-4) != v509 {
		v467 = v467 + int32(-16)
		v468 = v509
		goto L80
	} else {
		goto L82
	}
L81:
	;
	v522 = v509
	goto L78
L82:
	;
	goto L81
L83:
	;
	v559 = v392 - v522<<(uint(int32(2))%32)
	v560 = v542
	goto L84
L84:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v559))) = v448
	v582 = v560 + int32(-1)
	if v582 != 0 {
		v559 = v559 + int32(-4)
		v560 = v582
		goto L84
	} else {
		goto L86
	}
L85:
	;
	goto L67
L86:
	;
	goto L85
L87:
	;
	if v177 == v400 {
		goto L16
	} else {
		goto L91
	}
L88:
	;
	if v177 <= v400 {
		v635 = int32(0)
		goto L87
	} else {
		goto L90
	}
L89:
	;
	v614 = int32(2)
	v616 = l0 + v400<<(uint(v614)%32)
	v617 = *(*int32)(unsafe.Add(mBase, uint32(v616)))
	v618 = *(*int32)(unsafe.Add(mBase, uint32(v616)+4))
	v620 = *(*int32)(unsafe.Add(mBase, uint32(v616)+8))
	v622 = *(*int32)(unsafe.Add(mBase, uint32(v616)+12))
	v635 = int32(base.Ui32(v617+v618+v620+v622+v614) >> (uint(v614) % 32))
	goto L87
L90:
	;
	v633 = *(*int32)(unsafe.Add(mBase, uint32(l0+v400<<(uint(int32(2))%32))))
	v635 = v633
	goto L87
L91:
	;
	v639 = *(*int32)(unsafe.Add(mBase, uint32(l0+v400<<(uint(int32(2))%32))))
	v652 = v639
	v661 = v635
	v665 = int32(1)
	goto L58
L92:
	;
	v649 = base.I32_div_u_s(v641+int32(base.Ui32(v643)>>(uint(int32(1))%32)), v643)
	v652 = v641
	v661 = v649
	v665 = v643
	goto L58
L93:
	;
	v723 = v36 & int32(3)
	if base.Ui32(int32(4)) <= base.Ui32(v36) {
		goto L95
	} else {
		goto L96
	}
L94:
	;
	if v723 == int32(0) {
		v881 = v808
		goto L100
	} else {
		goto L101
	}
L95:
	;
	v730 = int32(0)
	v734 = l0
	v744 = v730
	v751 = v730
	goto L97
L96:
	;
	v726 = int32(0)
	v801 = v726
	v808 = v726
	goto L94
L97:
	;
	v762 = *(*int32)(unsafe.Add(mBase, uint32(v734)))
	v763 = int32(0)
	v766 = int32(4)
	v768 = *(*int32)(unsafe.Add(mBase, uint32(v734+v766)))
	v774 = *(*int32)(unsafe.Add(mBase, uint32(v734+int32(8))))
	v780 = *(*int32)(unsafe.Add(mBase, uint32(v734+int32(12))))
	v783 = v751 + base.B2i32(v762 != v763) + base.B2i32(v768 != v763) + base.B2i32(v774 != v763) + base.B2i32(v780 != v763)
	v787 = v744 + v766
	if v36&int32(2147483644) != v787 {
		v734 = v734 + int32(16)
		v744 = v787
		v751 = v783
		goto L97
	} else {
		goto L99
	}
L98:
	;
	v801 = v787
	v808 = v783
	goto L94
L99:
	;
	goto L98
L100:
	;
	if v881 == int32(0) {
		v2209 = v719
		goto L15
	} else {
		goto L105
	}
L101:
	;
	v826 = l0 + v801<<(uint(int32(2))%32)
	v835 = v723
	v843 = v808
	goto L102
L102:
	;
	v854 = *(*int32)(unsafe.Add(mBase, uint32(v826)))
	v857 = v843 + base.B2i32(v854 != int32(0))
	v861 = v835 + int32(-1)
	if v861 != 0 {
		v826 = v826 + int32(4)
		v835 = v861
		v843 = v857
		goto L102
	} else {
		goto L104
	}
L103:
	;
	v881 = v857
	goto L100
L104:
	;
	goto L103
L105:
	;
	v894 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	v896 = v36 + int32(-1)
	if int32(1) < v881 {
		goto L106
	} else {
		goto L107
	}
L106:
	;
	v1246 = l3 + v881<<(uint(int32(4))%32)
	v1249 = int32(3)
	v1250 = v896 & v1249
	v1253 = int32(1)
	v1255 = int32(-2)
	v1269 = v1253
	goto L162
L107:
	;
	v901 = int32(3)
	v902 = v896 & v901
	v905 = int32(1)
	v922 = v905
	goto L108
L108:
	;
	if v896 != 0 {
		goto L111
	} else {
		goto L112
	}
L110:
	;
	if v36&v905 == int32(0) {
		goto L126
	} else {
		goto L127
	}
L111:
	;
	v944 = int32(0)
	v948 = l0
	v957 = v944
	v958 = v944
	goto L113
L112:
	;
	v942 = int32(0)
	v1028 = v942
	v1029 = v942
	goto L110
L113:
	;
	v976 = *(*int32)(unsafe.Add(mBase, uint32(v948)))
	if v976 == int32(0) {
		v990 = v957
		goto L115
	} else {
		goto L116
	}
L114:
	;
	v1028 = v1010
	v1029 = v1015
	goto L110
L115:
	;
	v994 = *(*int32)(unsafe.Add(mBase, uint32(v948+int32(4))))
	if v994 == int32(0) {
		v1010 = v990
		goto L120
	} else {
		goto L121
	}
L116:
	;
	v981 = l3 + v957<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v981)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v981)+4)) = v958
	if base.Ui32(v922) < base.Ui32(v976) {
		goto L117
	} else {
		goto L118
	}
L117:
	;
	v986 = v976
	goto L119
L118:
	;
	v986 = v922
	goto L119
L119:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v981))) = v986
	v990 = v957 + int32(1)
	goto L115
L120:
	;
	v1015 = v958 + int32(2)
	if v36&int32(2147483646) != v1015 {
		v948 = v948 + int32(8)
		v957 = v1010
		v958 = v1015
		goto L113
	} else {
		goto L125
	}
L121:
	;
	v999 = l3 + v990<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v999)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v999)+4)) = v958 + int32(1)
	if base.Ui32(v922) < base.Ui32(v994) {
		goto L122
	} else {
		goto L123
	}
L122:
	;
	v1006 = v994
	goto L124
L123:
	;
	v1006 = v922
	goto L124
L124:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v999))) = v1006
	v1010 = v990 + int32(1)
	goto L120
L125:
	;
	goto L114
L126:
	;
	v1067 = m.G2
	F___qsort_r(m, l3, v881, int32(16), int32(182), v1067+int32(177))
	mBase = m.M
	goto L132
L127:
	;
	v1052 = *(*int32)(unsafe.Add(mBase, uint32(l0+v1029<<(uint(int32(2))%32))))
	if v1052 == int32(0) {
		goto L126
	} else {
		goto L128
	}
L128:
	;
	v1057 = l3 + v1028<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1057)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1057)+4)) = v1029
	if base.Ui32(v922) < base.Ui32(v1052) {
		goto L129
	} else {
		goto L130
	}
L129:
	;
	v1062 = v1052
	goto L131
L130:
	;
	v1062 = v922
	goto L131
L131:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1057))) = v1062
	goto L126
L132:
	;
	if v881 != int32(1) {
		goto L133
	} else {
		goto L134
	}
L133:
	;
	v1078 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894))))
	if v36 == int32(1) {
		v1211 = v1078
		goto L135
	} else {
		goto L136
	}
L134:
	;
	v1074 = *(*int32)(unsafe.Add(mBase, uint32(l3)+4))
	v1076 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v894+v1074))) = uint8(v1076)
	goto L133
L135:
	;
	if l1 < v1211 {
		v922 = v922 << (uint(int32(1)) % 32)
		goto L108
	} else {
		goto L161
	}
L136:
	;
	if base.Ui32(v36+int32(-2)) < base.Ui32(v901) {
		v1141 = v1078
		v1150 = int32(1)
		goto L137
	} else {
		goto L138
	}
L137:
	;
	if v902 == int32(0) {
		v1211 = v1141
		goto L135
	} else {
		goto L154
	}
L138:
	;
	v1085 = v1078
	v1095 = int32(0)
	goto L139
L139:
	;
	v1113 = v894 + v1095
	v1116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1113+int32(1)))))
	if v1116 < v1085 {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v1141 = v1133
	v1150 = v1095 + int32(5)
	goto L137
L141:
	;
	v1118 = v1085
	goto L143
L142:
	;
	v1118 = v1116
	goto L143
L143:
	;
	v1121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1113+int32(2)))))
	if v1121 < v1118 {
		goto L144
	} else {
		goto L145
	}
L144:
	;
	v1123 = v1118
	goto L146
L145:
	;
	v1123 = v1121
	goto L146
L146:
	;
	v1126 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1113+int32(3)))))
	if v1126 < v1123 {
		goto L147
	} else {
		goto L148
	}
L147:
	;
	v1128 = v1123
	goto L149
L148:
	;
	v1128 = v1126
	goto L149
L149:
	;
	v1131 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1113+int32(4)))))
	if v1131 < v1128 {
		goto L150
	} else {
		goto L151
	}
L150:
	;
	v1133 = v1128
	goto L152
L151:
	;
	v1133 = v1131
	goto L152
L152:
	;
	v1135 = v1095 + int32(4)
	if v896&int32(-4) != v1135 {
		v1085 = v1133
		v1095 = v1135
		goto L139
	} else {
		goto L153
	}
L153:
	;
	goto L140
L154:
	;
	v1174 = v1141
	v1183 = v894 + v1150
	v1184 = v902
	goto L155
L155:
	;
	v1202 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1183))))
	if v1202 < v1174 {
		goto L157
	} else {
		goto L158
	}
L156:
	;
	v1211 = v1204
	goto L135
L157:
	;
	v1204 = v1174
	goto L159
L158:
	;
	v1204 = v1202
	goto L159
L159:
	;
	v1208 = v1184 + int32(-1)
	if v1208 != 0 {
		v1174 = v1204
		v1183 = v1183 + int32(1)
		v1184 = v1208
		goto L155
	} else {
		goto L160
	}
L160:
	;
	goto L156
L161:
	;
	v2209 = v719
	goto L15
L162:
	;
	v1292 = int32(0)
	if v896 == v1292 {
		v1380 = v1292
		v1381 = v1292
		goto L164
	} else {
		goto L165
	}
L163:
	;
	v2209 = v719
	goto L15
L164:
	;
	if v36&v1253 == int32(0) {
		goto L179
	} else {
		goto L180
	}
L165:
	;
	v1296 = int32(0)
	v1300 = l0
	v1309 = v1296
	v1310 = v1296
	goto L166
L166:
	;
	v1328 = *(*int32)(unsafe.Add(mBase, uint32(v1300)))
	if v1328 == int32(0) {
		v1342 = v1309
		goto L168
	} else {
		goto L169
	}
L167:
	;
	v1380 = v1362
	v1381 = v1367
	goto L164
L168:
	;
	v1346 = *(*int32)(unsafe.Add(mBase, uint32(v1300+int32(4))))
	if v1346 == int32(0) {
		v1362 = v1342
		goto L173
	} else {
		goto L174
	}
L169:
	;
	v1333 = l3 + v1309<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1333)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1333)+4)) = v1310
	if base.Ui32(v1269) < base.Ui32(v1328) {
		goto L170
	} else {
		goto L171
	}
L170:
	;
	v1338 = v1328
	goto L172
L171:
	;
	v1338 = v1269
	goto L172
L172:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1333))) = v1338
	v1342 = v1309 + int32(1)
	goto L168
L173:
	;
	v1367 = v1310 + int32(2)
	if v36&int32(2147483646) != v1367 {
		v1300 = v1300 + int32(8)
		v1309 = v1362
		v1310 = v1367
		goto L166
	} else {
		goto L178
	}
L174:
	;
	v1351 = l3 + v1342<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1351)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1351)+4)) = v1310 + int32(1)
	if base.Ui32(v1269) < base.Ui32(v1346) {
		goto L175
	} else {
		goto L176
	}
L175:
	;
	v1358 = v1346
	goto L177
L176:
	;
	v1358 = v1269
	goto L177
L177:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1351))) = v1358
	v1362 = v1342 + int32(1)
	goto L173
L178:
	;
	goto L167
L179:
	;
	v1419 = m.G2
	F___qsort_r(m, l3, v881, int32(16), int32(182), v1419+int32(177))
	mBase = m.M
	goto L185
L180:
	;
	v1404 = *(*int32)(unsafe.Add(mBase, uint32(l0+v1381<<(uint(int32(2))%32))))
	if v1404 == int32(0) {
		goto L179
	} else {
		goto L181
	}
L181:
	;
	v1409 = l3 + v1380<<(uint(int32(4))%32)
	*(*int64)(unsafe.Add(mBase, uint32(v1409)+8)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1409)+4)) = v1381
	if base.Ui32(v1269) < base.Ui32(v1404) {
		goto L182
	} else {
		goto L183
	}
L182:
	;
	v1414 = v1404
	goto L184
L183:
	;
	v1414 = v1269
	goto L184
L184:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1409))) = v1414
	goto L179
L185:
	;
	v1427 = v881
	v1434 = int32(0)
	v1437 = v881 + v1255
	goto L186
L186:
	;
	v1455 = int32(4)
	v1457 = v1246 + v1434<<(uint(v1455)%32)
	v1460 = l3 + int32(-16) + v1427<<(uint(v1455)%32)
	v1461 = *(*int64)(unsafe.Add(mBase, uint32(v1460)))
	*(*int64)(unsafe.Add(mBase, uint32(v1457))) = v1461
	v1463 = int32(8)
	v1467 = *(*int64)(unsafe.Add(mBase, uint32(v1460+v1463)))
	*(*int64)(unsafe.Add(mBase, uint32(v1457+v1463))) = v1467
	v1470 = v1434 | int32(1)
	v1473 = v1246 + v1470<<(uint(v1455)%32)
	v1475 = v1427 + int32(-2)
	v1478 = l3 + v1475<<(uint(v1455)%32)
	v1479 = *(*int64)(unsafe.Add(mBase, uint32(v1478)))
	*(*int64)(unsafe.Add(mBase, uint32(v1473))) = v1479
	v1485 = *(*int64)(unsafe.Add(mBase, uint32(v1478+v1463)))
	*(*int64)(unsafe.Add(mBase, uint32(v1473+v1463))) = v1485
	v1489 = base.I32_wrap_i64(v1461) + base.I32_wrap_i64(v1479)
	if v1427 < int32(3) {
		v1542 = int32(0)
		goto L188
	} else {
		goto L189
	}
L187:
	;
	v1996 = int32(0)
	v1998 = *(*int32)(unsafe.Add(mBase, uint32(l3)+8))
	if v1998 < v1996 {
		v2019 = l3
		v2022 = v1996
		goto L276
	} else {
		goto L277
	}
L188:
	;
	v1563 = int32(4)
	v1565 = l3 + v1542<<(uint(v1563)%32)
	v1567 = v1565 + int32(16)
	v1570 = (v1475 - v1542) << (uint(v1563) % 32)
	if base.Ui32(int32(33)) <= base.Ui32(v1570) {
		goto L197
	} else {
		goto L198
	}
L189:
	;
	v1496 = l3
	v1505 = int32(0)
	goto L190
L190:
	;
	v1524 = *(*int32)(unsafe.Add(mBase, uint32(v1496)))
	if base.Ui32(v1524) <= base.Ui32(v1489) {
		v1542 = v1505
		goto L188
	} else {
		goto L192
	}
L191:
	;
	v1542 = v1437
	goto L188
L192:
	;
	v1529 = v1505 + int32(1)
	if v1437 != v1529 {
		v1496 = v1496 + int32(16)
		v1505 = v1529
		goto L190
	} else {
		goto L193
	}
L193:
	;
	goto L191
L194:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1565)+12)) = v1434
	*(*int32)(unsafe.Add(mBase, uint32(v1565)+8)) = v1470
	v1987 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v1565)+4)) = v1987
	*(*int32)(unsafe.Add(mBase, uint32(v1565))) = v1489
	if int32(2) < v1427 {
		v1427 = v1427 + v1987
		v1434 = v1434 + int32(2)
		v1437 = v1437 + v1987
		goto L186
	} else {
		goto L274
	}
L195:
	;
	goto L194
L196:
	;
	v1589 = (v1565 ^ v1567) & int32(3)
	if base.Ui32(v1565) <= base.Ui32(v1567) {
		goto L203
	} else {
		goto L204
	}
L197:
	;
	base.MemoryCopy(m, v1567, v1565, v1570)
	goto L195
L198:
	;
	if v1567 == v1565 {
		goto L195
	} else {
		goto L199
	}
L199:
	;
	v1578 = v1567 + v1570
	if base.Ui32(int32(0)-v1570<<(uint(int32(1))%32)) < base.Ui32(v1565-v1578) {
		goto L196
	} else {
		goto L200
	}
L200:
	;
	goto L197
L201:
	;
	if v1917 == int32(0) {
		goto L195
	} else {
		goto L263
	}
L202:
	;
	if base.Ui32(v1841) < base.Ui32(int32(4)) {
		v1913 = v1839
		v1915 = v1840
		v1917 = v1841
		goto L201
	} else {
		goto L253
	}
L203:
	;
	if v1589 != 0 {
		v1765 = v1570
		goto L219
	} else {
		goto L220
	}
L204:
	;
	if v1589 == int32(0) {
		goto L205
	} else {
		goto L206
	}
L205:
	;
	if v1567&int32(3) != 0 {
		goto L207
	} else {
		goto L208
	}
L206:
	;
	v1913 = v1565
	v1915 = v1567
	v1917 = v1570
	goto L201
L207:
	;
	if v1570 == int32(0) {
		goto L195
	} else {
		goto L209
	}
L208:
	;
	v1839 = v1565
	v1840 = v1567
	v1841 = v1570
	goto L202
L209:
	;
	v1597 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1567))) = uint8(v1597)
	v1600 = v1570 + int32(-1)
	v1602 = v1565 + int32(17)
	if v1602&int32(3) != 0 {
		goto L210
	} else {
		goto L211
	}
L210:
	;
	if v1600 == int32(0) {
		goto L195
	} else {
		goto L212
	}
L211:
	;
	v1839 = v1565 + int32(1)
	v1840 = v1602
	v1841 = v1600
	goto L202
L212:
	;
	v1609 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565)+1)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1567)+1)) = uint8(v1609)
	v1612 = v1570 + int32(-2)
	v1614 = v1565 + int32(18)
	if v1614&int32(3) != 0 {
		goto L213
	} else {
		goto L214
	}
L213:
	;
	if v1612 == int32(0) {
		goto L195
	} else {
		goto L215
	}
L214:
	;
	v1839 = v1565 + int32(2)
	v1840 = v1614
	v1841 = v1612
	goto L202
L215:
	;
	v1621 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565)+2)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1567)+2)) = uint8(v1621)
	v1624 = v1570 + int32(-3)
	v1626 = v1565 + int32(19)
	if v1626&int32(3) != 0 {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	if v1624 == int32(0) {
		goto L195
	} else {
		goto L218
	}
L217:
	;
	v1839 = v1565 + int32(3)
	v1840 = v1626
	v1841 = v1624
	goto L202
L218:
	;
	v1633 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565)+3)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1567)+3)) = uint8(v1633)
	v1839 = v1565 + int32(4)
	v1840 = v1565 + int32(20)
	v1841 = v1570 + int32(-4)
	goto L202
L219:
	;
	if v1765 == int32(0) {
		goto L195
	} else {
		goto L243
	}
L220:
	;
	if v1578&int32(3) == int32(0) {
		v1683 = v1570
		goto L221
	} else {
		goto L222
	}
L221:
	;
	if base.Ui32(v1683) < base.Ui32(int32(4)) {
		v1765 = v1683
		goto L219
	} else {
		goto L233
	}
L222:
	;
	if v1570 == int32(0) {
		goto L195
	} else {
		goto L223
	}
L223:
	;
	v1648 = v1570 + int32(-1)
	v1649 = v1567 + v1648
	v1651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565+v1648))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1649))) = uint8(v1651)
	if v1649&int32(3) != 0 {
		goto L224
	} else {
		goto L225
	}
L224:
	;
	if v1648 == int32(0) {
		goto L195
	} else {
		goto L226
	}
L225:
	;
	v1683 = v1648
	goto L221
L226:
	;
	v1658 = v1570 + int32(-2)
	v1659 = v1567 + v1658
	v1661 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565+v1658))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1659))) = uint8(v1661)
	if v1659&int32(3) != 0 {
		goto L227
	} else {
		goto L228
	}
L227:
	;
	if v1658 == int32(0) {
		goto L195
	} else {
		goto L229
	}
L228:
	;
	v1683 = v1658
	goto L221
L229:
	;
	v1668 = v1570 + int32(-3)
	v1669 = v1567 + v1668
	v1671 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565+v1668))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1669))) = uint8(v1671)
	if v1669&int32(3) != 0 {
		goto L230
	} else {
		goto L231
	}
L230:
	;
	if v1668 == int32(0) {
		goto L195
	} else {
		goto L232
	}
L231:
	;
	v1683 = v1668
	goto L221
L232:
	;
	v1678 = v1570 + int32(-4)
	v1681 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565+v1678))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1567+v1678))) = uint8(v1681)
	v1683 = v1678
	goto L221
L233:
	;
	v1689 = v1683 + int32(-4)
	v1695 = (int32(base.Ui32(v1689)>>(uint(int32(2))%32)) + int32(1)) & int32(3)
	if v1695 == int32(0) {
		v1719 = v1683
		goto L234
	} else {
		goto L235
	}
L234:
	;
	if base.Ui32(v1689) < base.Ui32(int32(12)) {
		v1765 = v1719
		goto L219
	} else {
		goto L239
	}
L235:
	;
	v1704 = v1683
	v1705 = v1695
	goto L236
L236:
	;
	v1711 = *(*int32)(unsafe.Add(mBase, uint32(v1565+int32(-4)+v1704)))
	*(*int32)(unsafe.Add(mBase, uint32(v1565+int32(12)+v1704))) = v1711
	v1714 = v1704 + int32(-4)
	v1716 = v1705 + int32(-1)
	if v1716 != 0 {
		v1704 = v1714
		v1705 = v1716
		goto L236
	} else {
		goto L238
	}
L237:
	;
	v1719 = v1714
	goto L234
L238:
	;
	goto L237
L239:
	;
	v1732 = v1719
	goto L240
L240:
	;
	v1737 = v1565 + int32(0) + v1732
	v1738 = int32(12)
	v1740 = v1565 + int32(-16) + v1732
	v1743 = *(*int32)(unsafe.Add(mBase, uint32(v1740+v1738)))
	*(*int32)(unsafe.Add(mBase, uint32(v1737+v1738))) = v1743
	v1745 = int32(8)
	v1749 = *(*int32)(unsafe.Add(mBase, uint32(v1740+v1745)))
	*(*int32)(unsafe.Add(mBase, uint32(v1737+v1745))) = v1749
	v1751 = int32(4)
	v1755 = *(*int32)(unsafe.Add(mBase, uint32(v1740+v1751)))
	*(*int32)(unsafe.Add(mBase, uint32(v1737+v1751))) = v1755
	v1757 = *(*int32)(unsafe.Add(mBase, uint32(v1740)))
	*(*int32)(unsafe.Add(mBase, uint32(v1737))) = v1757
	v1760 = v1732 + int32(-16)
	if base.Ui32(int32(3)) < base.Ui32(v1760) {
		v1732 = v1760
		goto L240
	} else {
		goto L242
	}
L241:
	;
	v1765 = v1760
	goto L219
L242:
	;
	goto L241
L243:
	;
	v1773 = v1765 & int32(3)
	if v1773 == int32(0) {
		v1798 = v1765
		goto L244
	} else {
		goto L245
	}
L244:
	;
	if base.Ui32(v1765) < base.Ui32(int32(4)) {
		goto L195
	} else {
		goto L249
	}
L245:
	;
	v1783 = v1765
	v1784 = v1773
	goto L246
L246:
	;
	v1789 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1565+int32(-1)+v1783))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1565+int32(15)+v1783))) = uint8(v1789)
	v1791 = int32(-1)
	v1792 = v1783 + v1791
	v1794 = v1784 + v1791
	if v1794 != 0 {
		v1783 = v1792
		v1784 = v1794
		goto L246
	} else {
		goto L248
	}
L247:
	;
	v1798 = v1792
	goto L244
L248:
	;
	goto L247
L249:
	;
	v1811 = v1798
	goto L250
L250:
	;
	v1815 = v1565 + int32(12) + v1811
	v1816 = int32(3)
	v1818 = v1565 + int32(-4) + v1811
	v1821 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1818+v1816))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1815+v1816))) = uint8(v1821)
	v1823 = int32(2)
	v1827 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1818+v1823))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1815+v1823))) = uint8(v1827)
	v1829 = int32(1)
	v1833 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1818+v1829))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1815+v1829))) = uint8(v1833)
	v1835 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1818))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1815))) = uint8(v1835)
	v1838 = v1811 + int32(-4)
	if v1838 != 0 {
		v1811 = v1838
		goto L250
	} else {
		goto L252
	}
L252:
	;
	goto L195
L253:
	;
	v1845 = v1841 + int32(-4)
	v1851 = (int32(base.Ui32(v1845)>>(uint(int32(2))%32)) + int32(1)) & int32(7)
	if v1851 == int32(0) {
		v1873 = v1839
		v1875 = v1840
		v1877 = v1841
		goto L254
	} else {
		goto L255
	}
L254:
	;
	if base.Ui32(v1845) < base.Ui32(int32(28)) {
		v1913 = v1873
		v1915 = v1875
		v1917 = v1877
		goto L201
	} else {
		goto L259
	}
L255:
	;
	v1858 = v1839
	v1859 = v1851
	v1860 = v1840
	goto L256
L256:
	;
	v1864 = *(*int32)(unsafe.Add(mBase, uint32(v1858)))
	*(*int32)(unsafe.Add(mBase, uint32(v1860))) = v1864
	v1866 = int32(4)
	v1867 = v1858 + v1866
	v1869 = v1860 + v1866
	v1871 = v1859 + int32(-1)
	if v1871 != 0 {
		v1858 = v1867
		v1859 = v1871
		v1860 = v1869
		goto L256
	} else {
		goto L258
	}
L257:
	;
	v1873 = v1867
	v1875 = v1869
	v1877 = v1841 - v1851<<(uint(int32(2))%32)
	goto L254
L258:
	;
	goto L257
L259:
	;
	v1882 = v1873
	v1884 = v1875
	v1886 = v1877
	goto L260
L260:
	;
	v1888 = *(*int32)(unsafe.Add(mBase, uint32(v1882)))
	*(*int32)(unsafe.Add(mBase, uint32(v1884))) = v1888
	v1890 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+4)) = v1890
	v1892 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+8)) = v1892
	v1894 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+12)) = v1894
	v1896 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+16)) = v1896
	v1898 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+20)) = v1898
	v1900 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+24))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+24)) = v1900
	v1902 = *(*int32)(unsafe.Add(mBase, uint32(v1882)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v1884)+28)) = v1902
	v1904 = int32(32)
	v1905 = v1882 + v1904
	v1907 = v1884 + v1904
	v1909 = v1886 + int32(-32)
	if base.Ui32(int32(3)) < base.Ui32(v1909) {
		v1882 = v1905
		v1884 = v1907
		v1886 = v1909
		goto L260
	} else {
		goto L262
	}
L261:
	;
	v1913 = v1905
	v1915 = v1907
	v1917 = v1909
	goto L201
L262:
	;
	goto L261
L263:
	;
	v1922 = v1917 & int32(7)
	if v1922 != 0 {
		goto L265
	} else {
		goto L266
	}
L264:
	;
	if base.Ui32(v1917) < base.Ui32(int32(8)) {
		goto L195
	} else {
		goto L270
	}
L265:
	;
	v1926 = v1913
	v1927 = v1922
	v1928 = v1915
	goto L267
L266:
	;
	v1941 = v1913
	v1943 = v1915
	v1944 = v1917
	goto L264
L267:
	;
	v1932 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1926))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1928))) = uint8(v1932)
	v1934 = int32(1)
	v1935 = v1928 + v1934
	v1937 = v1926 + v1934
	v1939 = v1927 + int32(-1)
	if v1939 != 0 {
		v1926 = v1937
		v1927 = v1939
		v1928 = v1935
		goto L267
	} else {
		goto L269
	}
L268:
	;
	v1941 = v1937
	v1943 = v1935
	v1944 = v1917 & int32(-8)
	goto L264
L269:
	;
	goto L268
L270:
	;
	v1950 = v1941
	v1952 = v1943
	v1953 = v1944
	goto L271
L271:
	;
	v1956 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952))) = uint8(v1956)
	v1958 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+1)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+1)) = uint8(v1958)
	v1960 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+2)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+2)) = uint8(v1960)
	v1962 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+3)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+3)) = uint8(v1962)
	v1964 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+4)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+4)) = uint8(v1964)
	v1966 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+5)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+5)) = uint8(v1966)
	v1968 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+6)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+6)) = uint8(v1968)
	v1970 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1950)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(v1952)+7)) = uint8(v1970)
	v1972 = int32(8)
	v1977 = v1953 + int32(-8)
	if v1977 != 0 {
		v1950 = v1950 + v1972
		v1952 = v1952 + v1972
		v1953 = v1977
		goto L271
	} else {
		goto L273
	}
L272:
	;
	goto L195
L273:
	;
	goto L272
L274:
	;
	goto L187
L275:
	;
	v2027 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894))))
	if v36 == int32(1) {
		v2160 = v2027
		goto L281
	} else {
		goto L282
	}
L276:
	;
	v2024 = *(*int32)(unsafe.Add(mBase, uint32(v2019)+4))
	*(*uint8)(unsafe.Add(mBase, uint32(v894+v2024))) = uint8(v2022)
	goto L275
L277:
	;
	v2001 = l3
	v2004 = v1996
	v2005 = v1998
	goto L278
L278:
	;
	v2006 = int32(4)
	v2010 = v2004 + int32(1)
	F_SetBitDepths(m, v1246+v2005<<(uint(v2006)%32), v1246, v894, v2010)
	mBase = m.M
	v2012 = *(*int32)(unsafe.Add(mBase, uint32(v2001)+12))
	v2015 = v1246 + v2012<<(uint(v2006)%32)
	v2016 = *(*int32)(unsafe.Add(mBase, uint32(v2015)+8))
	if int32(-1) < v2016 {
		v2001 = v2015
		v2004 = v2010
		v2005 = v2016
		goto L278
	} else {
		goto L280
	}
L279:
	;
	v2019 = v2015
	v2022 = v2010
	goto L276
L280:
	;
	goto L279
L281:
	;
	if l1 < v2160 {
		v1269 = v1269 << (uint(int32(1)) % 32)
		goto L162
	} else {
		goto L307
	}
L282:
	;
	if base.Ui32(v36+v1255) < base.Ui32(v1249) {
		v2090 = v2027
		v2099 = int32(1)
		goto L283
	} else {
		goto L284
	}
L283:
	;
	if v1250 == int32(0) {
		v2160 = v2090
		goto L281
	} else {
		goto L300
	}
L284:
	;
	v2034 = v2027
	v2044 = int32(0)
	goto L285
L285:
	;
	v2062 = v894 + v2044
	v2065 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2062+int32(1)))))
	if v2065 < v2034 {
		goto L287
	} else {
		goto L288
	}
L286:
	;
	v2090 = v2082
	v2099 = v2044 + int32(5)
	goto L283
L287:
	;
	v2067 = v2034
	goto L289
L288:
	;
	v2067 = v2065
	goto L289
L289:
	;
	v2070 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2062+int32(2)))))
	if v2070 < v2067 {
		goto L290
	} else {
		goto L291
	}
L290:
	;
	v2072 = v2067
	goto L292
L291:
	;
	v2072 = v2070
	goto L292
L292:
	;
	v2075 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2062+int32(3)))))
	if v2075 < v2072 {
		goto L293
	} else {
		goto L294
	}
L293:
	;
	v2077 = v2072
	goto L295
L294:
	;
	v2077 = v2075
	goto L295
L295:
	;
	v2080 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2062+int32(4)))))
	if v2080 < v2077 {
		goto L296
	} else {
		goto L297
	}
L296:
	;
	v2082 = v2077
	goto L298
L297:
	;
	v2082 = v2080
	goto L298
L298:
	;
	v2084 = v2044 + int32(4)
	if v896&int32(-4) != v2084 {
		v2034 = v2082
		v2044 = v2084
		goto L285
	} else {
		goto L299
	}
L299:
	;
	goto L286
L300:
	;
	v2123 = v2090
	v2132 = v894 + v2099
	v2133 = v1250
	goto L301
L301:
	;
	v2151 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2132))))
	if v2151 < v2123 {
		goto L303
	} else {
		goto L304
	}
L302:
	;
	v2160 = v2153
	goto L281
L303:
	;
	v2153 = v2123
	goto L305
L304:
	;
	v2153 = v2151
	goto L305
L305:
	;
	v2157 = v2133 + int32(-1)
	if v2157 != 0 {
		v2123 = v2153
		v2132 = v2132 + int32(1)
		v2133 = v2157
		goto L301
	} else {
		goto L306
	}
L306:
	;
	goto L302
L307:
	;
	goto L163
L308:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v33)+64)) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+72)) = v2472
	v2504 = int32(1)
	v2505 = (v2472 + v2489) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+76)) = v2505
	v2509 = (v2505 + v2494) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+80)) = v2509
	v2513 = (v2509 + v2490) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+84)) = v2513
	v2517 = (v2513 + v2484) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+88)) = v2517
	v2521 = (v2517 + v2477) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+92)) = v2521
	v2525 = (v2521 + v2487) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+96)) = v2525
	v2529 = (v2525 + v2478) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+100)) = v2529
	v2533 = (v2529 + v2485) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+104)) = v2533
	v2537 = (v2533 + v2473) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+108)) = v2537
	v2541 = (v2537 + v2483) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+112)) = v2541
	v2545 = (v2541 + v2480) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+116)) = v2545
	v2549 = (v2545 + v2482) << (uint(v2504) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v33)+120)) = v2549
	*(*int32)(unsafe.Add(mBase, uint32(v33)+124)) = (v2549 + v2481) << (uint(v2504) % 32)
	if v2249 < int32(1) {
		goto L321
	} else {
		goto L322
	}
L309:
	;
	v2269 = v2249 & int32(3)
	v2270 = *(*int32)(unsafe.Add(mBase, uint32(v2209)))
	if base.Ui32(v2249) < base.Ui32(int32(4)) {
		v2360 = int32(0)
		goto L311
	} else {
		goto L312
	}
L310:
	;
	v2254 = int32(0)
	v2472 = v2254
	v2473 = v2254
	v2477 = v2254
	v2478 = v2254
	v2480 = v2254
	v2481 = v2254
	v2482 = v2254
	v2483 = v2254
	v2484 = v2254
	v2485 = v2254
	v2487 = v2254
	v2489 = v2254
	v2490 = v2254
	v2494 = v2254
	goto L308
L311:
	;
	if v2269 == int32(0) {
		goto L316
	} else {
		goto L317
	}
L312:
	;
	v2288 = int32(0)
	goto L313
L313:
	;
	v2307 = v2270 + v2288
	v2308 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2307))))
	v2309 = int32(2)
	v2311 = v33 + v2308<<(uint(v2309)%32)
	v2312 = *(*int32)(unsafe.Add(mBase, uint32(v2311)))
	v2313 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2311))) = v2312 + v2313
	v2318 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2307+v2313))))
	v2321 = v33 + v2318<<(uint(v2309)%32)
	v2322 = *(*int32)(unsafe.Add(mBase, uint32(v2321)))
	*(*int32)(unsafe.Add(mBase, uint32(v2321))) = v2322 + v2313
	v2328 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2307+v2309))))
	v2331 = v33 + v2328<<(uint(v2309)%32)
	v2332 = *(*int32)(unsafe.Add(mBase, uint32(v2331)))
	*(*int32)(unsafe.Add(mBase, uint32(v2331))) = v2332 + v2313
	v2338 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2307+int32(3)))))
	v2341 = v33 + v2338<<(uint(v2309)%32)
	v2342 = *(*int32)(unsafe.Add(mBase, uint32(v2341)))
	*(*int32)(unsafe.Add(mBase, uint32(v2341))) = v2342 + v2313
	v2347 = v2288 + int32(4)
	if v2249&int32(2147483644) != v2347 {
		v2288 = v2347
		goto L313
	} else {
		goto L315
	}
L314:
	;
	v2360 = v2347
	goto L311
L315:
	;
	goto L314
L316:
	;
	v2454 = *(*int32)(unsafe.Add(mBase, uint32(v33)+4))
	v2457 = *(*int32)(unsafe.Add(mBase, uint32(v33)+56))
	v2458 = *(*int32)(unsafe.Add(mBase, uint32(v33)+52))
	v2459 = *(*int32)(unsafe.Add(mBase, uint32(v33)+48))
	v2460 = *(*int32)(unsafe.Add(mBase, uint32(v33)+44))
	v2461 = *(*int32)(unsafe.Add(mBase, uint32(v33)+40))
	v2462 = *(*int32)(unsafe.Add(mBase, uint32(v33)+36))
	v2463 = *(*int32)(unsafe.Add(mBase, uint32(v33)+32))
	v2464 = *(*int32)(unsafe.Add(mBase, uint32(v33)+28))
	v2465 = *(*int32)(unsafe.Add(mBase, uint32(v33)+24))
	v2466 = *(*int32)(unsafe.Add(mBase, uint32(v33)+20))
	v2467 = *(*int32)(unsafe.Add(mBase, uint32(v33)+16))
	v2468 = *(*int32)(unsafe.Add(mBase, uint32(v33)+12))
	v2469 = *(*int32)(unsafe.Add(mBase, uint32(v33)+8))
	v2472 = v2454 << (uint(int32(1)) % 32)
	v2473 = v2461
	v2477 = v2465
	v2478 = v2463
	v2480 = v2459
	v2481 = v2457
	v2482 = v2458
	v2483 = v2460
	v2484 = v2466
	v2485 = v2462
	v2487 = v2464
	v2489 = v2469
	v2490 = v2467
	v2494 = v2468
	goto L308
L317:
	;
	v2384 = v2270 + v2360
	v2394 = v2269
	goto L318
L318:
	;
	v2412 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2384))))
	v2415 = v33 + v2412<<(uint(int32(2))%32)
	v2416 = *(*int32)(unsafe.Add(mBase, uint32(v2415)))
	v2417 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2415))) = v2416 + v2417
	v2423 = v2394 + int32(-1)
	if v2423 != 0 {
		v2384 = v2384 + v2417
		v2394 = v2423
		goto L318
	} else {
		goto L320
	}
L319:
	;
	goto L316
L320:
	;
	goto L319
L321:
	;
	m.G0 = v33 + int32(128)
	return
L322:
	;
	v2555 = *(*int32)(unsafe.Add(mBase, uint32(l4)+8))
	v2556 = *(*int32)(unsafe.Add(mBase, uint32(l4)+4))
	v2574 = int32(0)
	goto L323
L323:
	;
	v2591 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2556+v2574))))
	v2594 = v33 + int32(64) + v2591<<(uint(int32(2))%32)
	v2595 = *(*int32)(unsafe.Add(mBase, uint32(v2594)))
	*(*int32)(unsafe.Add(mBase, uint32(v2594))) = v2595 + int32(1)
	v2599 = int32(0)
	if v2591 == v2599 {
		v2730 = v2599
		goto L325
	} else {
		goto L326
	}
L324:
	;
	goto L321
L325:
	;
	v2748 = int32(1)
	v2753 = int32(base.Ui32(v2730) >> (uint(int32(16)-v2591) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v2555+v2574<<(uint(v2748)%32)))) = uint16(v2753)
	v2756 = v2574 + v2748
	if v2756 != v2249 {
		v2574 = v2756
		goto L323
	} else {
		goto L334
	}
L326:
	;
	v2603 = v2591 + int32(-1)
	if base.Ui32(int32(5)) <= base.Ui32(v2591) {
		goto L328
	} else {
		goto L329
	}
L327:
	;
	if v2603&int32(4) != 0 {
		v2730 = v2687
		goto L325
	} else {
		goto L333
	}
L328:
	;
	v2614 = int32(0)
	v2618 = v2595
	v2626 = (int32(base.Ui32(v2603)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
	v2627 = v2614
	v2628 = v2614
	goto L330
L329:
	;
	v2677 = v2595
	v2686 = int32(0)
	v2687 = int32(0)
	goto L327
L330:
	;
	v2646 = m.G1
	v2648 = v2646 + int32(_a_F_VP8LCreateHuffmanTree_0)
	v2651 = int32(15)
	v2654 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2648+int32(base.Ui32(v2618)>>(uint(int32(4))%32))&v2651))))
	v2655 = int32(8)
	v2661 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2648+v2618&v2651))))
	v2666 = v2654<<(uint(v2627+v2655)%32) | (v2661<<(uint(v2627+int32(12))%32) | v2628)
	v2668 = v2627 + int32(-8)
	v2670 = int32(base.Ui32(v2618) >> (uint(v2655) % 32))
	v2672 = v2626 + int32(-2)
	if v2672 != 0 {
		v2618 = v2670
		v2626 = v2672
		v2627 = v2668
		v2628 = v2666
		goto L330
	} else {
		goto L332
	}
L331:
	;
	v2677 = v2670
	v2686 = int32(0) - v2668
	v2687 = v2666
	goto L327
L332:
	;
	goto L331
L333:
	;
	v2707 = m.G1
	v2713 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2707+int32(_a_F_VP8LCreateHuffmanTree_0)+v2677&int32(15)))))
	v2730 = v2713<<(uint(int32(12)-v2686)%32) | v2687
	goto L325
L334:
	;
	goto L324
}
func F_VP8LDistanceToPlaneCode(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v12 int32
	_ = v12
	var v21 int32
	_ = v21
	var v34 int32
	_ = v34
	var v40 int32
	_ = v40
	v5 = base.I32_div_s(l1, l0)
	v7 = l1 - v5*l0
	if int32(7) < v5 {
		if int32(6) < v5 {
			return l1 + int32(120)
		} else {
			if v7 <= l0+int32(-8) {
				return l1 + int32(120)
			} else {
				v34 = m.G1
				v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v5<<(uint(int32(4))%32)-v7+(v34+int32(_a_F_VP8LDistanceToPlaneCode_0))+int32(24)))))
				return v40 + int32(1)
			}
		}
	} else {
		if int32(8) < v7 {
			if int32(6) < v5 {
				return l1 + int32(120)
			} else {
				if v7 <= l0+int32(-8) {
					return l1 + int32(120)
				} else {
					v34 = m.G1
					v40 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v5<<(uint(int32(4))%32)-v7+(v34+int32(_a_F_VP8LDistanceToPlaneCode_0))+int32(24)))))
					return v40 + int32(1)
				}
			}
		} else {
			v12 = m.G1
			v21 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v12+int32(_a_F_VP8LDistanceToPlaneCode_0)+(v5<<(uint(int32(4))%32)|int32(8)-v7)))))
			return v21 + int32(1)
		}
	}
}
func F_VP8LDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v31 int32
	_ = v31
	var v77 int32
	_ = v77
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v96 int32
	_ = v96
	var v99 int32
	_ = v99
	var v102 int32
	_ = v102
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	v19 = m.G1
	v22 = *(*int32)(unsafe.Add(mBase, uint32(v19)+uint32(_c_F_VP8LDspInit[0])))
	v23 = m.G3
	v24 = *(*int32)(unsafe.Add(mBase, uint32(v23)))
	if v22 == v24 {
	} else {
		v26 = m.G2
		v27 = m.G1
		v31 = v26 + int32(21)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[1]))) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[2]))) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[3]))) = v26 + int32(22)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[4]))) = v26 + int32(23)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[5]))) = v26 + int32(24)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[6]))) = v26 + int32(25)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[7]))) = v26 + int32(26)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[8]))) = v26 + int32(27)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[9]))) = v26 + int32(28)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[10]))) = v26 + int32(29)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[11]))) = v26 + int32(30)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[12]))) = v26 + int32(31)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[13]))) = v26 + int32(32)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[14]))) = v26 + int32(33)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[15]))) = v26 + int32(34)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[16]))) = v31
		v77 = v26 + int32(35)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[17]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[18]))) = v77
		v81 = v26 + int32(36)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[19]))) = v81
		v84 = v26 + int32(37)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[20]))) = v84
		v87 = v26 + int32(38)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[21]))) = v87
		v90 = v26 + int32(39)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[22]))) = v90
		v93 = v26 + int32(40)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[23]))) = v93
		v96 = v26 + int32(41)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[24]))) = v96
		v99 = v26 + int32(42)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[25]))) = v99
		v102 = v26 + int32(43)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[26]))) = v102
		v105 = v26 + int32(44)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[27]))) = v105
		v108 = v26 + int32(45)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[28]))) = v108
		v111 = v26 + int32(46)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[29]))) = v111
		v114 = v26 + int32(47)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[30]))) = v114
		v117 = v26 + int32(48)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[31]))) = v117
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[32]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[33]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[34]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[35]))) = v81
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[36]))) = v84
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[37]))) = v87
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[38]))) = v90
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[39]))) = v93
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[40]))) = v96
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[41]))) = v99
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[42]))) = v102
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[43]))) = v105
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[44]))) = v108
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[45]))) = v111
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[46]))) = v114
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[47]))) = v117
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[48]))) = v77
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[49]))) = v26 + int32(49)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[50]))) = v26 + int32(50)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[51]))) = v26 + int32(51)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[52]))) = v26 + int32(52)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[53]))) = v26 + int32(53)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[54]))) = v26 + int32(54)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[55]))) = v26 + int32(55)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[56]))) = v26 + int32(56)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[57]))) = v26 + int32(57)
		*(*int32)(unsafe.Add(mBase, uint32(v27)+uint32(_c_F_VP8LDspInit[0]))) = v24
	}
	return
}
func F_VP8LEncDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v45 int32
	_ = v45
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v103 int32
	_ = v103
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
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
	var v140 int32
	_ = v140
	var v143 int32
	_ = v143
	var v307 int32
	_ = v307
	var v311 int32
	_ = v311
	var v314 int32
	_ = v314
	var v317 int32
	_ = v317
	var v320 int32
	_ = v320
	var v323 int32
	_ = v323
	var v326 int32
	_ = v326
	var v329 int32
	_ = v329
	var v332 int32
	_ = v332
	var v335 int32
	_ = v335
	var v338 int32
	_ = v338
	var v341 int32
	_ = v341
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v370 int32
	_ = v370
	var v371 int32
	_ = v371
	v18 = m.G1
	v21 = *(*int32)(unsafe.Add(mBase, uint32(v18)+uint32(_c_F_VP8LEncDspInit[0])))
	v22 = m.G3
	v23 = *(*int32)(unsafe.Add(mBase, uint32(v22)))
	if v21 == v23 {
	} else {
		v25 = m.G2
		v26 = m.G1
		v45 = m.G1
		v48 = *(*int32)(unsafe.Add(mBase, uint32(v45)+uint32(_c_F_VP8LEncDspInit[1])))
		v49 = m.G3
		v50 = *(*int32)(unsafe.Add(mBase, uint32(v49)))
		if v48 == v50 {
		} else {
			v52 = m.G2
			v53 = m.G1
			v57 = v52 + int32(21)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[2]))) = v57
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[3]))) = v57
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[4]))) = v52 + int32(22)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[5]))) = v52 + int32(23)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[6]))) = v52 + int32(24)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[7]))) = v52 + int32(25)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[8]))) = v52 + int32(26)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[9]))) = v52 + int32(27)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[10]))) = v52 + int32(28)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[11]))) = v52 + int32(29)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[12]))) = v52 + int32(30)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[13]))) = v52 + int32(31)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[14]))) = v52 + int32(32)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[15]))) = v52 + int32(33)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[16]))) = v52 + int32(34)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[17]))) = v57
			v103 = v52 + int32(35)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[18]))) = v103
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[19]))) = v103
			v107 = v52 + int32(36)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[20]))) = v107
			v110 = v52 + int32(37)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[21]))) = v110
			v113 = v52 + int32(38)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[22]))) = v113
			v116 = v52 + int32(39)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[23]))) = v116
			v119 = v52 + int32(40)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[24]))) = v119
			v122 = v52 + int32(41)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[25]))) = v122
			v125 = v52 + int32(42)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[26]))) = v125
			v128 = v52 + int32(43)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[27]))) = v128
			v131 = v52 + int32(44)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[28]))) = v131
			v134 = v52 + int32(45)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[29]))) = v134
			v137 = v52 + int32(46)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[30]))) = v137
			v140 = v52 + int32(47)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[31]))) = v140
			v143 = v52 + int32(48)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[32]))) = v143
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[33]))) = v103
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[34]))) = v103
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[35]))) = v103
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[36]))) = v107
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[37]))) = v110
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[38]))) = v113
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[39]))) = v116
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[40]))) = v119
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[41]))) = v122
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[42]))) = v125
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[43]))) = v128
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[44]))) = v131
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[45]))) = v134
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[46]))) = v137
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[47]))) = v140
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[48]))) = v143
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[49]))) = v103
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[50]))) = v52 + int32(49)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[51]))) = v52 + int32(50)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[52]))) = v52 + int32(51)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[53]))) = v52 + int32(52)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[54]))) = v52 + int32(53)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[55]))) = v52 + int32(54)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[56]))) = v52 + int32(55)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[57]))) = v52 + int32(56)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[58]))) = v52 + int32(57)
			*(*int32)(unsafe.Add(mBase, uint32(v53)+uint32(_c_F_VP8LEncDspInit[1]))) = v50
		}
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[59]))) = v25 + int32(138)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[60]))) = v25 + int32(139)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[61]))) = v25 + int32(140)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[62]))) = v25 + int32(141)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[63]))) = v25 + int32(142)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[64]))) = v25 + int32(143)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[65]))) = v25 + int32(144)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[66]))) = v25 + int32(145)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[67]))) = v25 + int32(146)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[68]))) = v25 + int32(147)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[69]))) = v25 + int32(148)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[70]))) = v25 + int32(149)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[71]))) = v25 + int32(150)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[72]))) = v25 + int32(151)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[73]))) = v25 + int32(152)
		v307 = v25 + int32(153)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[74]))) = v307
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[75]))) = v307
		v311 = v25 + int32(154)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[76]))) = v311
		v314 = v25 + int32(155)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[77]))) = v314
		v317 = v25 + int32(156)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[78]))) = v317
		v320 = v25 + int32(157)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[79]))) = v320
		v323 = v25 + int32(158)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[80]))) = v323
		v326 = v25 + int32(159)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[81]))) = v326
		v329 = v25 + int32(160)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[82]))) = v329
		v332 = v25 + int32(161)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[83]))) = v332
		v335 = v25 + int32(162)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[84]))) = v335
		v338 = v25 + int32(163)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[85]))) = v338
		v341 = v25 + int32(164)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[86]))) = v341
		v344 = v25 + int32(165)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[87]))) = v344
		v347 = v25 + int32(166)
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[88]))) = v347
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[89]))) = v307
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[90]))) = v307
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[91]))) = v307
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[92]))) = v311
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[93]))) = v314
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[94]))) = v317
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[95]))) = v320
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[96]))) = v323
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[97]))) = v326
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[98]))) = v329
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[99]))) = v332
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[100]))) = v335
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[101]))) = v338
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[102]))) = v341
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[103]))) = v344
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[104]))) = v347
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[105]))) = v307
		v370 = m.G3
		v371 = *(*int32)(unsafe.Add(mBase, uint32(v370)))
		*(*int32)(unsafe.Add(mBase, uint32(v26)+uint32(_c_F_VP8LEncDspInit[0]))) = v371
	}
	return
}
func F_VP8LEncodeImage(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v31 int64
	_ = v31
	var v45 int32
	_ = v45
	var v46 int32
	_ = v46
	var v50 int32
	_ = v50
	var v59 int32
	_ = v59
	var v65 int32
	_ = v65
	var v68 int32
	_ = v68
	var v195 int64
	_ = v195
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v256 int64
	_ = v256
	var v259 int32
	_ = v259
	var v261 int64
	_ = v261
	var v264 int32
	_ = v264
	var v271 int32
	_ = v271
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v293 int32
	_ = v293
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v306 int32
	_ = v306
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v318 int32
	_ = v318
	var v356 int32
	_ = v356
	var v357 int32
	_ = v357
	var v369 int32
	_ = v369
	var v370 int32
	_ = v370
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v383 int32
	_ = v383
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v386 int32
	_ = v386
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v393 int32
	_ = v393
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v405 int32
	_ = v405
	var v406 int32
	_ = v406
	var v409 int64
	_ = v409
	var v412 int32
	_ = v412
	var v414 int64
	_ = v414
	var v417 int32
	_ = v417
	var v424 int32
	_ = v424
	var v426 int32
	_ = v426
	var v430 int32
	_ = v430
	var v431 int32
	_ = v431
	var v432 int32
	_ = v432
	var v433 int32
	_ = v433
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
	var v452 int32
	_ = v452
	var v459 int32
	_ = v459
	var v462 int32
	_ = v462
	var v464 int32
	_ = v464
	var v468 int32
	_ = v468
	var v469 int32
	_ = v469
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v508 int32
	_ = v508
	var v512 int32
	_ = v512
	var v515 int32
	_ = v515
	var v516 int32
	_ = v516
	var v528 int32
	_ = v528
	var v529 int32
	_ = v529
	var v533 int32
	_ = v533
	var v534 int32
	_ = v534
	var v542 int32
	_ = v542
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v545 int32
	_ = v545
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v555 int32
	_ = v555
	var v564 int32
	_ = v564
	var v565 int32
	_ = v565
	var v568 int64
	_ = v568
	var v571 int32
	_ = v571
	var v573 int64
	_ = v573
	var v576 int32
	_ = v576
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v592 int32
	_ = v592
	var v601 int32
	_ = v601
	var v602 int32
	_ = v602
	var v603 int32
	_ = v603
	var v605 int32
	_ = v605
	var v610 int32
	_ = v610
	var v611 int32
	_ = v611
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v627 int32
	_ = v627
	var v628 int32
	_ = v628
	var v629 int32
	_ = v629
	var v630 int32
	_ = v630
	var v667 int32
	_ = v667
	var v668 int32
	_ = v668
	var v680 int32
	_ = v680
	var v681 int32
	_ = v681
	var v685 int32
	_ = v685
	var v686 int32
	_ = v686
	var v694 int32
	_ = v694
	var v695 int32
	_ = v695
	var v696 int32
	_ = v696
	var v697 int32
	_ = v697
	var v699 int32
	_ = v699
	var v700 int32
	_ = v700
	var v704 int32
	_ = v704
	var v705 int32
	_ = v705
	var v706 int32
	_ = v706
	var v707 int32
	_ = v707
	var v716 int32
	_ = v716
	var v717 int32
	_ = v717
	var v720 int64
	_ = v720
	var v723 int32
	_ = v723
	var v725 int64
	_ = v725
	var v728 int32
	_ = v728
	var v735 int32
	_ = v735
	var v737 int32
	_ = v737
	var v741 int32
	_ = v741
	var v742 int32
	_ = v742
	var v743 int32
	_ = v743
	var v744 int32
	_ = v744
	var v753 int32
	_ = v753
	var v754 int32
	_ = v754
	var v755 int32
	_ = v755
	var v757 int32
	_ = v757
	var v762 int32
	_ = v762
	var v763 int32
	_ = v763
	var v770 int32
	_ = v770
	var v773 int32
	_ = v773
	var v775 int32
	_ = v775
	var v779 int32
	_ = v779
	var v780 int32
	_ = v780
	var v781 int32
	_ = v781
	var v782 int32
	_ = v782
	var v819 int32
	_ = v819
	var v823 int32
	_ = v823
	var v829 int32
	_ = v829
	var v832 int32
	_ = v832
	var v838 int32
	_ = v838
	var v843 int32
	_ = v843
	var v849 int32
	_ = v849
	var v852 int32
	_ = v852
	var v855 int32
	_ = v855
	var v857 int32
	_ = v857
	var v861 int32
	_ = v861
	var v864 int32
	_ = v864
	var v865 int32
	_ = v865
	var v867 int32
	_ = v867
	var v873 int32
	_ = v873
	var v884 int32
	_ = v884
	var v906 int32
	_ = v906
	var v909 int32
	_ = v909
	var v910 int32
	_ = v910
	var v918 int32
	_ = v918
	var v919 int32
	_ = v919
	var v954 int32
	_ = v954
	var v955 int32
	_ = v955
	var v960 int64
	_ = v960
	var v963 int32
	_ = v963
	var v964 int32
	_ = v964
	var v976 int32
	_ = v976
	var v996 int32
	_ = v996
	var v1001 int32
	_ = v1001
	var v1005 int32
	_ = v1005
	var v1012 int32
	_ = v1012
	var v1013 int32
	_ = v1013
	var v1017 int64
	_ = v1017
	var v1024 int32
	_ = v1024
	var v1028 int32
	_ = v1028
	v3 = int32(0)
	v7 = m.G0
	v9 = v7 - int32(32)
	m.G0 = v9
	*(*int32)(unsafe.Add(mBase, uint32(v9)+24)) = v3
	if l1 == v3 {
		v1028 = v3
		goto L1
	} else {
		goto L2
	}
L1:
	;
	m.G0 = v9 + int32(32)
	return v1028
L2:
	;
	if l0 == int32(0) {
		goto L4
	} else {
		goto L5
	}
L3:
	;
	v23 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v24 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v31 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v9))) = v31
	*(*int64)(unsafe.Add(mBase, uint32(v9+int32(16)))) = v31
	*(*int64)(unsafe.Add(mBase, uint32(v9+int32(8)))) = v31
	v45 = v23*v24<<(uint(base.B2i32(v26 != int32(3)))%32)&int32(-1024) + int32(1024)
	v46 = F_WebPSafeMalloc(m, int64(1), v45)
	mBase = m.M
	if v46 != 0 {
		goto L13
	} else {
		goto L14
	}
L4:
	;
	v20 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v20 != 0 {
		goto L8
	} else {
		goto L9
	}
L5:
	;
	v18 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v18 != 0 {
		goto L3
	} else {
		goto L6
	}
L6:
	;
	goto L4
L7:
	;
	v1028 = int32(0)
	goto L1
L8:
	;
	goto L7
L9:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(3)
	goto L8
L10:
	;
	v1001 = *(*int32)(unsafe.Add(mBase, uint32(v9)+20))
	if v1001 == int32(0) {
		goto L176
	} else {
		goto L177
	}
L11:
	;
	v65 = F_WebPReportProgress(m, l1, int32(1), v9+int32(24))
	mBase = m.M
	if v65 == int32(0) {
		goto L19
	} else {
		goto L20
	}
L12:
	;
	if v46 != 0 {
		goto L11
	} else {
		goto L15
	}
L13:
	;
	v50 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v50)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v46 + v45
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v46
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v46
	goto L12
L14:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	goto L12
L15:
	;
	v59 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v59 != 0 {
		goto L17
	} else {
		goto L18
	}
L16:
	;
	goto L10
L17:
	;
	goto L16
L18:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L17
L19:
	;
	v996 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v996 != 0 {
		goto L174
	} else {
		goto L175
	}
L20:
	;
	v68 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v68 == int32(0) {
		goto L21
	} else {
		goto L22
	}
L21:
	;
	v200 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v201 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v203 = v201 + int32(-1)
	v204 = int32(14)
	goto L39
L22:
	;
	goto L25
L23:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v68)+20)) = int32(1120272384)
	v195 = int64(4811533253012226048)
	*(*int64)(unsafe.Add(mBase, uint32(v68)+12)) = v195
	*(*int64)(unsafe.Add(mBase, uint32(v68)+4)) = v195
	goto L21
L25:
	;
	base.MemoryFill(m, v68, int32(0), int32(188))
	goto L23
L37:
	;
	v356 = v200 + int32(-1)
	v357 = int32(14)
	goto L65
L38:
	;
	goto L37
L39:
	;
	v216 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v217 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v217+v204 < int32(32) {
		goto L42
	} else {
		goto L43
	}
L40:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v318 + v316
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v315<<(uint(v318)%32) | v317
	goto L38
L41:
	;
	v235 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v236 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v240 = v232
	v241 = v233
	v242 = v236
	v243 = v235
	goto L45
L42:
	;
	if v217 < int32(16) {
		v315 = v203
		v316 = v204
		v317 = v216
		v318 = v217
		goto L40
	} else {
		goto L44
	}
L43:
	;
	v221 = int32(32)
	v222 = v221 - v217
	v230 = int32(base.Ui32(v203) >> (uint(v222) % 32))
	v231 = v204 - v222
	v232 = v203<<(uint(v217)%32) | v216
	v233 = v221
	goto L41
L44:
	;
	v230 = v203
	v231 = v204
	v232 = v216
	v233 = v217
	goto L41
L45:
	;
	if base.Ui32(v242+int32(2)) <= base.Ui32(v243) {
		v298 = v242
		v299 = v243
		goto L47
	} else {
		goto L48
	}
L46:
	;
	v315 = v230
	v316 = v231
	v317 = v311
	v318 = v309
	goto L40
L47:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v298))) = uint16(v240)
	v306 = v298 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v306
	v309 = v241 + int32(-16)
	v311 = int32(base.Ui32(v240) >> (uint(int32(16)) % 32))
	if int32(31) < v241 {
		v240 = v311
		v241 = v309
		v242 = v306
		v243 = v299
		goto L45
	} else {
		goto L62
	}
L48:
	;
	v252 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v253 = v243 - v252
	v256 = base.I64_extend_i32_s(v253) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v256) {
		v280 = v252
		goto L50
	} else {
		goto L51
	}
L49:
	;
	if v242 == v252 {
		goto L60
	} else {
		goto L61
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v280
	goto L37
L51:
	;
	v259 = v242 - v252
	v261 = v256 + base.I64_extend_i32_u(v259)
	if base.Ui64(int64(4294967295)) < base.Ui64(v261) {
		v280 = v252
		goto L50
	} else {
		goto L52
	}
L52:
	;
	v264 = base.I32_wrap_i64(v261)
	if v243 == v252 {
		goto L53
	} else {
		goto L54
	}
L53:
	;
	v271 = int32(base.Ui32(v253*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v264) < base.Ui32(v271) {
		goto L56
	} else {
		goto L57
	}
L54:
	;
	if base.Ui32(v264) <= base.Ui32(v253) {
		v298 = v242
		v299 = v243
		goto L47
	} else {
		goto L55
	}
L55:
	;
	goto L53
L56:
	;
	v273 = v271
	goto L58
L57:
	;
	v273 = v264
	goto L58
L58:
	;
	v277 = v273&int32(-1024) + int32(1024)
	v278 = F_WebPSafeMalloc(m, int64(1), v277)
	mBase = m.M
	if v278 != 0 {
		goto L49
	} else {
		goto L59
	}
L59:
	;
	v279 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v280 = v279
	goto L50
L60:
	;
	v291 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v291)
	mBase = m.M
	v293 = v278 + v277
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v293
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v278
	v298 = v278 + v259
	v299 = v293
	goto L47
L61:
	;
	v289 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v290 = F_memcpy(m, v278, v289, v259)
	mBase = m.M
	goto L60
L62:
	;
	goto L46
L63:
	;
	v508 = *(*int32)(unsafe.Add(mBase, uint32(v9)+20))
	if v508 == int32(0) {
		goto L89
	} else {
		goto L90
	}
L64:
	;
	goto L63
L65:
	;
	v369 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v370 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v370+v357 < int32(32) {
		goto L68
	} else {
		goto L69
	}
L66:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v471 + v469
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v468<<(uint(v471)%32) | v470
	goto L64
L67:
	;
	v388 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v389 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v393 = v385
	v394 = v386
	v395 = v389
	v396 = v388
	goto L71
L68:
	;
	if v370 < int32(16) {
		v468 = v356
		v469 = v357
		v470 = v369
		v471 = v370
		goto L66
	} else {
		goto L70
	}
L69:
	;
	v374 = int32(32)
	v375 = v374 - v370
	v383 = int32(base.Ui32(v356) >> (uint(v375) % 32))
	v384 = v357 - v375
	v385 = v356<<(uint(v370)%32) | v369
	v386 = v374
	goto L67
L70:
	;
	v383 = v356
	v384 = v357
	v385 = v369
	v386 = v370
	goto L67
L71:
	;
	if base.Ui32(v395+int32(2)) <= base.Ui32(v396) {
		v451 = v395
		v452 = v396
		goto L73
	} else {
		goto L74
	}
L72:
	;
	v468 = v383
	v469 = v384
	v470 = v464
	v471 = v462
	goto L66
L73:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v451))) = uint16(v393)
	v459 = v451 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v459
	v462 = v394 + int32(-16)
	v464 = int32(base.Ui32(v393) >> (uint(int32(16)) % 32))
	if int32(31) < v394 {
		v393 = v464
		v394 = v462
		v395 = v459
		v396 = v452
		goto L71
	} else {
		goto L88
	}
L74:
	;
	v405 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v406 = v396 - v405
	v409 = base.I64_extend_i32_s(v406) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v409) {
		v433 = v405
		goto L76
	} else {
		goto L77
	}
L75:
	;
	if v395 == v405 {
		goto L86
	} else {
		goto L87
	}
L76:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v433
	goto L63
L77:
	;
	v412 = v395 - v405
	v414 = v409 + base.I64_extend_i32_u(v412)
	if base.Ui64(int64(4294967295)) < base.Ui64(v414) {
		v433 = v405
		goto L76
	} else {
		goto L78
	}
L78:
	;
	v417 = base.I32_wrap_i64(v414)
	if v396 == v405 {
		goto L79
	} else {
		goto L80
	}
L79:
	;
	v424 = int32(base.Ui32(v406*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v417) < base.Ui32(v424) {
		goto L82
	} else {
		goto L83
	}
L80:
	;
	if base.Ui32(v417) <= base.Ui32(v406) {
		v451 = v395
		v452 = v396
		goto L73
	} else {
		goto L81
	}
L81:
	;
	goto L79
L82:
	;
	v426 = v424
	goto L84
L83:
	;
	v426 = v417
	goto L84
L84:
	;
	v430 = v426&int32(-1024) + int32(1024)
	v431 = F_WebPSafeMalloc(m, int64(1), v430)
	mBase = m.M
	if v431 != 0 {
		goto L75
	} else {
		goto L85
	}
L85:
	;
	v432 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v433 = v432
	goto L76
L86:
	;
	v444 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v444)
	mBase = m.M
	v446 = v431 + v430
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v446
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v431
	v451 = v431 + v412
	v452 = v446
	goto L73
L87:
	;
	v442 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v443 = F_memcpy(m, v431, v442, v412)
	mBase = m.M
	goto L86
L88:
	;
	goto L72
L89:
	;
	v515 = F_WebPPictureHasTransparency(m, l1)
	mBase = m.M
	v516 = int32(1)
	goto L96
L90:
	;
	v512 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v512 != 0 {
		goto L92
	} else {
		goto L93
	}
L91:
	;
	goto L10
L92:
	;
	goto L91
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L92
L94:
	;
	v667 = int32(0)
	v668 = int32(3)
	goto L122
L95:
	;
	goto L94
L96:
	;
	v528 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v529 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v529+v516 < int32(32) {
		goto L99
	} else {
		goto L100
	}
L97:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v630 + v628
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v627<<(uint(v630)%32) | v629
	goto L95
L98:
	;
	v547 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v548 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v552 = v544
	v553 = v545
	v554 = v548
	v555 = v547
	goto L102
L99:
	;
	if v529 < int32(16) {
		v627 = v515
		v628 = v516
		v629 = v528
		v630 = v529
		goto L97
	} else {
		goto L101
	}
L100:
	;
	v533 = int32(32)
	v534 = v533 - v529
	v542 = int32(base.Ui32(v515) >> (uint(v534) % 32))
	v543 = v516 - v534
	v544 = v515<<(uint(v529)%32) | v528
	v545 = v533
	goto L98
L101:
	;
	v542 = v515
	v543 = v516
	v544 = v528
	v545 = v529
	goto L98
L102:
	;
	if base.Ui32(v554+int32(2)) <= base.Ui32(v555) {
		v610 = v554
		v611 = v555
		goto L104
	} else {
		goto L105
	}
L103:
	;
	v627 = v542
	v628 = v543
	v629 = v623
	v630 = v621
	goto L97
L104:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v610))) = uint16(v552)
	v618 = v610 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v618
	v621 = v553 + int32(-16)
	v623 = int32(base.Ui32(v552) >> (uint(int32(16)) % 32))
	if int32(31) < v553 {
		v552 = v623
		v553 = v621
		v554 = v618
		v555 = v611
		goto L102
	} else {
		goto L119
	}
L105:
	;
	v564 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v565 = v555 - v564
	v568 = base.I64_extend_i32_s(v565) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v568) {
		v592 = v564
		goto L107
	} else {
		goto L108
	}
L106:
	;
	if v554 == v564 {
		goto L117
	} else {
		goto L118
	}
L107:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v592
	goto L94
L108:
	;
	v571 = v554 - v564
	v573 = v568 + base.I64_extend_i32_u(v571)
	if base.Ui64(int64(4294967295)) < base.Ui64(v573) {
		v592 = v564
		goto L107
	} else {
		goto L109
	}
L109:
	;
	v576 = base.I32_wrap_i64(v573)
	if v555 == v564 {
		goto L110
	} else {
		goto L111
	}
L110:
	;
	v583 = int32(base.Ui32(v565*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v576) < base.Ui32(v583) {
		goto L113
	} else {
		goto L114
	}
L111:
	;
	if base.Ui32(v576) <= base.Ui32(v565) {
		v610 = v554
		v611 = v555
		goto L104
	} else {
		goto L112
	}
L112:
	;
	goto L110
L113:
	;
	v585 = v583
	goto L115
L114:
	;
	v585 = v576
	goto L115
L115:
	;
	v589 = v585&int32(-1024) + int32(1024)
	v590 = F_WebPSafeMalloc(m, int64(1), v589)
	mBase = m.M
	if v590 != 0 {
		goto L106
	} else {
		goto L116
	}
L116:
	;
	v591 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v592 = v591
	goto L107
L117:
	;
	v603 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v603)
	mBase = m.M
	v605 = v590 + v589
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v605
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v590
	v610 = v590 + v571
	v611 = v605
	goto L104
L118:
	;
	v601 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v602 = F_memcpy(m, v590, v601, v571)
	mBase = m.M
	goto L117
L119:
	;
	goto L103
L120:
	;
	v819 = *(*int32)(unsafe.Add(mBase, uint32(v9)+20))
	if v819 == int32(0) {
		goto L146
	} else {
		goto L147
	}
L121:
	;
	goto L120
L122:
	;
	v680 = *(*int32)(unsafe.Add(mBase, uint32(v9)))
	v681 = *(*int32)(unsafe.Add(mBase, uint32(v9)+4))
	if v681+v668 < int32(32) {
		goto L125
	} else {
		goto L126
	}
L123:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+4)) = v782 + v780
	*(*int32)(unsafe.Add(mBase, uint32(v9))) = v779<<(uint(v782)%32) | v781
	goto L121
L124:
	;
	v699 = *(*int32)(unsafe.Add(mBase, uint32(v9)+16))
	v700 = *(*int32)(unsafe.Add(mBase, uint32(v9)+12))
	v704 = v696
	v705 = v697
	v706 = v700
	v707 = v699
	goto L128
L125:
	;
	if v681 < int32(16) {
		v779 = v667
		v780 = v668
		v781 = v680
		v782 = v681
		goto L123
	} else {
		goto L127
	}
L126:
	;
	v685 = int32(32)
	v686 = v685 - v681
	v694 = int32(base.Ui32(v667) >> (uint(v686) % 32))
	v695 = v668 - v686
	v696 = v667<<(uint(v681)%32) | v680
	v697 = v685
	goto L124
L127:
	;
	v694 = v667
	v695 = v668
	v696 = v680
	v697 = v681
	goto L124
L128:
	;
	if base.Ui32(v706+int32(2)) <= base.Ui32(v707) {
		v762 = v706
		v763 = v707
		goto L130
	} else {
		goto L131
	}
L129:
	;
	v779 = v694
	v780 = v695
	v781 = v775
	v782 = v773
	goto L123
L130:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v762))) = uint16(v704)
	v770 = v762 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v770
	v773 = v705 + int32(-16)
	v775 = int32(base.Ui32(v704) >> (uint(int32(16)) % 32))
	if int32(31) < v705 {
		v704 = v775
		v705 = v773
		v706 = v770
		v707 = v763
		goto L128
	} else {
		goto L145
	}
L131:
	;
	v716 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v717 = v707 - v716
	v720 = base.I64_extend_i32_s(v717) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v720) {
		v744 = v716
		goto L133
	} else {
		goto L134
	}
L132:
	;
	if v706 == v716 {
		goto L143
	} else {
		goto L144
	}
L133:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v9)+12)) = v744
	goto L120
L134:
	;
	v723 = v706 - v716
	v725 = v720 + base.I64_extend_i32_u(v723)
	if base.Ui64(int64(4294967295)) < base.Ui64(v725) {
		v744 = v716
		goto L133
	} else {
		goto L135
	}
L135:
	;
	v728 = base.I32_wrap_i64(v725)
	if v707 == v716 {
		goto L136
	} else {
		goto L137
	}
L136:
	;
	v735 = int32(base.Ui32(v717*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v728) < base.Ui32(v735) {
		goto L139
	} else {
		goto L140
	}
L137:
	;
	if base.Ui32(v728) <= base.Ui32(v717) {
		v762 = v706
		v763 = v707
		goto L130
	} else {
		goto L138
	}
L138:
	;
	goto L136
L139:
	;
	v737 = v735
	goto L141
L140:
	;
	v737 = v728
	goto L141
L141:
	;
	v741 = v737&int32(-1024) + int32(1024)
	v742 = F_WebPSafeMalloc(m, int64(1), v741)
	mBase = m.M
	if v742 != 0 {
		goto L132
	} else {
		goto L142
	}
L142:
	;
	v743 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v744 = v743
	goto L133
L143:
	;
	v755 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	F_WebPSafeFree(m, v755)
	mBase = m.M
	v757 = v742 + v741
	*(*int32)(unsafe.Add(mBase, uint32(v9)+16)) = v757
	*(*int32)(unsafe.Add(mBase, uint32(v9)+8)) = v742
	v762 = v742 + v723
	v763 = v757
	goto L130
L144:
	;
	v753 = *(*int32)(unsafe.Add(mBase, uint32(v9)+8))
	v754 = F_memcpy(m, v742, v753, v723)
	mBase = m.M
	goto L143
L145:
	;
	goto L129
L146:
	;
	v829 = F_WebPReportProgress(m, l1, int32(2), v9+int32(24))
	mBase = m.M
	if v829 == int32(0) {
		goto L19
	} else {
		goto L151
	}
L147:
	;
	v823 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v823 != 0 {
		goto L149
	} else {
		goto L150
	}
L148:
	;
	goto L10
L149:
	;
	goto L148
L150:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L149
L151:
	;
	v832 = F_VP8LEncodeStream(m, l0, l1, v9)
	mBase = m.M
	if v832 == int32(0) {
		goto L10
	} else {
		goto L152
	}
L152:
	;
	v838 = F_WebPReportProgress(m, l1, int32(99), v9+int32(24))
	mBase = m.M
	if v838 == int32(0) {
		goto L19
	} else {
		goto L153
	}
L153:
	;
	v843 = F_WriteImage(m, l1, v9, v9+int32(28))
	mBase = m.M
	if v843 == int32(0) {
		goto L10
	} else {
		goto L154
	}
L154:
	;
	v849 = F_WebPReportProgress(m, l1, int32(100), v9+int32(24))
	mBase = m.M
	if v849 == int32(0) {
		goto L19
	} else {
		goto L155
	}
L155:
	;
	v852 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v852 == int32(0) {
		goto L156
	} else {
		goto L157
	}
L156:
	;
	v861 = *(*int32)(unsafe.Add(mBase, uint32(l1)+84))
	if v861 == int32(0) {
		goto L10
	} else {
		goto L158
	}
L157:
	;
	v855 = *(*int32)(unsafe.Add(mBase, uint32(v9)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v852)+168)) = v855
	v857 = *(*int32)(unsafe.Add(mBase, uint32(v852)))
	*(*int32)(unsafe.Add(mBase, uint32(v852))) = v855 + v857
	goto L156
L158:
	;
	v864 = int32(0)
	v865 = int32(15)
	v867 = int32(4)
	v873 = (v23 + v865) >> (uint(v867) % 32) * ((v24 + v865) >> (uint(v867) % 32))
	if base.Ui32(v873) < base.Ui32(int32(33)) {
		goto L160
	} else {
		goto L161
	}
L159:
	;
	goto L10
L160:
	;
	if v873 == int32(0) {
		goto L162
	} else {
		goto L163
	}
L161:
	;
	base.MemoryFill(m, v861, v864, v873)
	goto L159
L162:
	;
	goto L159
L163:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v861))) = uint8(v864)
	v884 = v861 + v873
	*(*uint8)(unsafe.Add(mBase, uint32(v884+int32(-1)))) = uint8(v864)
	if base.Ui32(v873) < base.Ui32(int32(3)) {
		goto L162
	} else {
		goto L164
	}
L164:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v861)+2)) = uint8(v864)
	*(*uint8)(unsafe.Add(mBase, uint32(v861)+1)) = uint8(v864)
	*(*uint8)(unsafe.Add(mBase, uint32(v884+int32(-3)))) = uint8(v864)
	*(*uint8)(unsafe.Add(mBase, uint32(v884+int32(-2)))) = uint8(v864)
	if base.Ui32(v873) < base.Ui32(int32(7)) {
		goto L162
	} else {
		goto L165
	}
L165:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v861)+3)) = uint8(v864)
	*(*uint8)(unsafe.Add(mBase, uint32(v884+int32(-4)))) = uint8(v864)
	if base.Ui32(v873) < base.Ui32(int32(9)) {
		goto L162
	} else {
		goto L166
	}
L166:
	;
	v906 = int32(0)
	v909 = (v906 - v861) & int32(3)
	v910 = v861 + v909
	*(*int32)(unsafe.Add(mBase, uint32(v910))) = v906
	v918 = (v873 - v909) & int32(60)
	v919 = v910 + v918
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-4)))) = v906
	if base.Ui32(v918) < base.Ui32(int32(9)) {
		goto L162
	} else {
		goto L167
	}
L167:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v910)+8)) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v910)+4)) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-8)))) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-12)))) = v906
	if base.Ui32(v918) < base.Ui32(int32(25)) {
		goto L162
	} else {
		goto L168
	}
L168:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v910)+24)) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v910)+20)) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v910)+16)) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v910)+12)) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-16)))) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-20)))) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-24)))) = v906
	*(*int32)(unsafe.Add(mBase, uint32(v919+int32(-28)))) = v906
	v954 = v910&int32(4) | int32(24)
	v955 = v918 - v954
	if base.Ui32(v955) < base.Ui32(int32(32)) {
		goto L162
	} else {
		goto L169
	}
L169:
	;
	v960 = base.I64_extend_i32_u(v906) * int64(4294967297)
	v963 = v955
	v964 = v910 + v954
	goto L170
L170:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v964)+24)) = v960
	*(*int64)(unsafe.Add(mBase, uint32(v964)+16)) = v960
	*(*int64)(unsafe.Add(mBase, uint32(v964)+8)) = v960
	*(*int64)(unsafe.Add(mBase, uint32(v964))) = v960
	v976 = v963 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v976) {
		v963 = v976
		v964 = v964 + int32(32)
		goto L170
	} else {
		goto L172
	}
L171:
	;
	goto L162
L172:
	;
	goto L171
L173:
	;
	goto L10
L174:
	;
	goto L173
L175:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(10)
	goto L174
L176:
	;
	if v9 == int32(0) {
		goto L182
	} else {
		goto L183
	}
L177:
	;
	v1005 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v1005 != 0 {
		goto L179
	} else {
		goto L180
	}
L178:
	;
	goto L176
L179:
	;
	goto L178
L180:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L179
L181:
	;
	v1024 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	v1028 = base.B2i32(v1024 == int32(0))
	goto L1
L182:
	;
	goto L181
L183:
	;
	v1012 = v9 + int32(8)
	v1013 = *(*int32)(unsafe.Add(mBase, uint32(v1012)))
	F_WebPSafeFree(m, v1013)
	mBase = m.M
	v1017 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v9+int32(16)))) = v1017
	*(*int64)(unsafe.Add(mBase, uint32(v1012))) = v1017
	*(*int64)(unsafe.Add(mBase, uint32(v9))) = v1017
	goto L182
}
func F_VP8LEncodeStream(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v58 int64
	_ = v58
	var v59 int32
	_ = v59
	var v66 int64
	_ = v66
	var v67 int32
	_ = v67
	var v68 int64
	_ = v68
	var v77 int32
	_ = v77
	var v79 int32
	_ = v79
	var v82 int32
	_ = v82
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v101 int64
	_ = v101
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v120 int32
	_ = v120
	var v129 int32
	_ = v129
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	var v138 int32
	_ = v138
	var v141 int32
	_ = v141
	var v143 int32
	_ = v143
	var v145 int32
	_ = v145
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v172 int32
	_ = v172
	var v174 int32
	_ = v174
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v203 int32
	_ = v203
	var v214 int32
	_ = v214
	var v215 int32
	_ = v215
	var v222 int32
	_ = v222
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v230 int32
	_ = v230
	var v232 int32
	_ = v232
	var v243 int32
	_ = v243
	var v244 int32
	_ = v244
	var v250 int32
	_ = v250
	var v258 int32
	_ = v258
	var v262 int32
	_ = v262
	var v273 int32
	_ = v273
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v281 int32
	_ = v281
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v288 int32
	_ = v288
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v303 int32
	_ = v303
	var v305 int32
	_ = v305
	var v308 int32
	_ = v308
	var v312 int32
	_ = v312
	var v315 int32
	_ = v315
	var v316 int32
	_ = v316
	var v319 int32
	_ = v319
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v326 int32
	_ = v326
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v335 int32
	_ = v335
	var v336 int32
	_ = v336
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v350 int32
	_ = v350
	var v356 int32
	_ = v356
	var v363 int32
	_ = v363
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v370 int32
	_ = v370
	var v383 int32
	_ = v383
	var v390 int32
	_ = v390
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v395 int32
	_ = v395
	var v404 int32
	_ = v404
	var v405 int32
	_ = v405
	var v418 int32
	_ = v418
	var v429 int32
	_ = v429
	var v430 int32
	_ = v430
	var v440 int32
	_ = v440
	var v442 int32
	_ = v442
	var v447 int32
	_ = v447
	var v452 int32
	_ = v452
	var v463 int32
	_ = v463
	var v467 int32
	_ = v467
	var v468 int32
	_ = v468
	var v475 int32
	_ = v475
	var v479 int32
	_ = v479
	var v480 int32
	_ = v480
	var v486 int32
	_ = v486
	var v490 int32
	_ = v490
	var v493 int32
	_ = v493
	var v501 int32
	_ = v501
	var v505 int32
	_ = v505
	var v509 int32
	_ = v509
	var v513 int32
	_ = v513
	var v521 int32
	_ = v521
	var v535 int32
	_ = v535
	var v536 int32
	_ = v536
	var v542 int32
	_ = v542
	var v543 int32
	_ = v543
	var v544 int32
	_ = v544
	var v547 int32
	_ = v547
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v554 int32
	_ = v554
	var v557 int32
	_ = v557
	var v559 int32
	_ = v559
	var v562 int32
	_ = v562
	var v572 int32
	_ = v572
	var v621 int32
	_ = v621
	var v624 int32
	_ = v624
	var v628 int32
	_ = v628
	var v648 int32
	_ = v648
	var v651 int32
	_ = v651
	var v690 int32
	_ = v690
	var v692 int32
	_ = v692
	var v693 int32
	_ = v693
	var v716 int32
	_ = v716
	var v755 int32
	_ = v755
	var v760 int32
	_ = v760
	var v761 int32
	_ = v761
	var v762 int32
	_ = v762
	var v788 int32
	_ = v788
	var v824 int32
	_ = v824
	var v830 int32
	_ = v830
	var v833 int32
	_ = v833
	var v835 int32
	_ = v835
	var v838 int32
	_ = v838
	var v842 int32
	_ = v842
	var v843 int32
	_ = v843
	var v847 int32
	_ = v847
	var v849 int32
	_ = v849
	var v853 int32
	_ = v853
	var v854 int32
	_ = v854
	var v855 int32
	_ = v855
	var v860 int64
	_ = v860
	var v861 int32
	_ = v861
	var v868 int64
	_ = v868
	var v869 int32
	_ = v869
	var v870 int64
	_ = v870
	var v879 int32
	_ = v879
	var v881 int32
	_ = v881
	var v915 int32
	_ = v915
	var v916 int32
	_ = v916
	var v923 int32
	_ = v923
	var __phi923 int32
	_ = __phi923
	var v933 int32
	_ = v933
	var __phi933 int32
	_ = __phi933
	var v954 int32
	_ = v954
	var __phi954 int32
	_ = __phi954
	var v955 int32
	_ = v955
	var __phi955 int32
	_ = __phi955
	var v977 int32
	_ = v977
	var v987 int32
	_ = v987
	var v988 int32
	_ = v988
	var v1026 int32
	_ = v1026
	var v1029 int32
	_ = v1029
	var v1031 int32
	_ = v1031
	var v1036 int32
	_ = v1036
	var v1038 int32
	_ = v1038
	var v1041 int32
	_ = v1041
	var v1047 int32
	_ = v1047
	var v1049 int32
	_ = v1049
	var v1051 int32
	_ = v1051
	var v1053 int32
	_ = v1053
	var v1054 int32
	_ = v1054
	var v1055 int32
	_ = v1055
	var v1058 int32
	_ = v1058
	var v1062 int32
	_ = v1062
	var v1063 int32
	_ = v1063
	var v1067 int32
	_ = v1067
	var v1071 int32
	_ = v1071
	var v1072 int32
	_ = v1072
	var v1076 int32
	_ = v1076
	var v1078 int32
	_ = v1078
	var v1080 int32
	_ = v1080
	var v1081 int32
	_ = v1081
	var v1089 int32
	_ = v1089
	var v1090 int32
	_ = v1090
	var v1098 int32
	_ = v1098
	var v1099 int32
	_ = v1099
	var v1107 int32
	_ = v1107
	var v1108 int32
	_ = v1108
	var v1116 int32
	_ = v1116
	var v1117 int32
	_ = v1117
	var v1121 int32
	_ = v1121
	var v1123 int32
	_ = v1123
	var v1124 int32
	_ = v1124
	var v1130 int32
	_ = v1130
	var v1131 int32
	_ = v1131
	var v1140 int32
	_ = v1140
	var v1141 int32
	_ = v1141
	var v1148 int32
	_ = v1148
	var v1154 int32
	_ = v1154
	var v1155 int32
	_ = v1155
	var v1164 int32
	_ = v1164
	var v1165 int32
	_ = v1165
	var v1178 int32
	_ = v1178
	var v1179 int32
	_ = v1179
	var v1189 int32
	_ = v1189
	var v1192 int32
	_ = v1192
	var v1247 int32
	_ = v1247
	var v1248 int32
	_ = v1248
	var v1251 int32
	_ = v1251
	var v1255 int32
	_ = v1255
	var v1259 int32
	_ = v1259
	var v1263 int32
	_ = v1263
	var v1267 int32
	_ = v1267
	var v1271 int32
	_ = v1271
	var v1276 int32
	_ = v1276
	var v1278 int32
	_ = v1278
	var v1283 int32
	_ = v1283
	var v1295 int64
	_ = v1295
	var v1298 int64
	_ = v1298
	var v1299 int64
	_ = v1299
	var v1302 int64
	_ = v1302
	var v1305 int64
	_ = v1305
	var v1307 int64
	_ = v1307
	var v1308 int64
	_ = v1308
	var v1311 int32
	_ = v1311
	var v1314 int32
	_ = v1314
	var v1323 int64
	_ = v1323
	var v1326 int64
	_ = v1326
	var v1329 int64
	_ = v1329
	var v1331 int64
	_ = v1331
	var v1332 int64
	_ = v1332
	var v1343 int32
	_ = v1343
	var v1345 int32
	_ = v1345
	var v1350 int32
	_ = v1350
	var v1362 int64
	_ = v1362
	var v1365 int64
	_ = v1365
	var v1366 int64
	_ = v1366
	var v1369 int64
	_ = v1369
	var v1372 int64
	_ = v1372
	var v1374 int64
	_ = v1374
	var v1375 int64
	_ = v1375
	var v1378 int32
	_ = v1378
	var v1381 int32
	_ = v1381
	var v1390 int64
	_ = v1390
	var v1393 int64
	_ = v1393
	var v1396 int64
	_ = v1396
	var v1398 int64
	_ = v1398
	var v1399 int64
	_ = v1399
	var v1410 int32
	_ = v1410
	var v1412 int32
	_ = v1412
	var v1417 int32
	_ = v1417
	var v1429 int64
	_ = v1429
	var v1432 int64
	_ = v1432
	var v1433 int64
	_ = v1433
	var v1436 int64
	_ = v1436
	var v1439 int64
	_ = v1439
	var v1441 int64
	_ = v1441
	var v1442 int64
	_ = v1442
	var v1445 int32
	_ = v1445
	var v1448 int32
	_ = v1448
	var v1457 int64
	_ = v1457
	var v1460 int64
	_ = v1460
	var v1463 int64
	_ = v1463
	var v1465 int64
	_ = v1465
	var v1466 int64
	_ = v1466
	var v1477 int32
	_ = v1477
	var v1479 int32
	_ = v1479
	var v1484 int32
	_ = v1484
	var v1496 int64
	_ = v1496
	var v1499 int64
	_ = v1499
	var v1500 int64
	_ = v1500
	var v1503 int64
	_ = v1503
	var v1506 int64
	_ = v1506
	var v1508 int64
	_ = v1508
	var v1509 int64
	_ = v1509
	var v1512 int32
	_ = v1512
	var v1515 int32
	_ = v1515
	var v1524 int64
	_ = v1524
	var v1527 int64
	_ = v1527
	var v1530 int64
	_ = v1530
	var v1532 int64
	_ = v1532
	var v1533 int64
	_ = v1533
	var v1544 int32
	_ = v1544
	var v1546 int32
	_ = v1546
	var v1551 int32
	_ = v1551
	var v1563 int64
	_ = v1563
	var v1566 int64
	_ = v1566
	var v1567 int64
	_ = v1567
	var v1570 int64
	_ = v1570
	var v1573 int64
	_ = v1573
	var v1575 int64
	_ = v1575
	var v1576 int64
	_ = v1576
	var v1579 int32
	_ = v1579
	var v1582 int32
	_ = v1582
	var v1591 int64
	_ = v1591
	var v1594 int64
	_ = v1594
	var v1597 int64
	_ = v1597
	var v1599 int64
	_ = v1599
	var v1600 int64
	_ = v1600
	var v1611 int32
	_ = v1611
	var v1613 int32
	_ = v1613
	var v1618 int32
	_ = v1618
	var v1630 int64
	_ = v1630
	var v1633 int64
	_ = v1633
	var v1634 int64
	_ = v1634
	var v1637 int64
	_ = v1637
	var v1640 int64
	_ = v1640
	var v1642 int64
	_ = v1642
	var v1643 int64
	_ = v1643
	var v1646 int32
	_ = v1646
	var v1649 int32
	_ = v1649
	var v1658 int64
	_ = v1658
	var v1661 int64
	_ = v1661
	var v1664 int64
	_ = v1664
	var v1666 int64
	_ = v1666
	var v1667 int64
	_ = v1667
	var v1678 int32
	_ = v1678
	var v1680 int32
	_ = v1680
	var v1685 int32
	_ = v1685
	var v1697 int64
	_ = v1697
	var v1700 int64
	_ = v1700
	var v1701 int64
	_ = v1701
	var v1704 int64
	_ = v1704
	var v1707 int64
	_ = v1707
	var v1709 int64
	_ = v1709
	var v1710 int64
	_ = v1710
	var v1713 int32
	_ = v1713
	var v1716 int32
	_ = v1716
	var v1725 int64
	_ = v1725
	var v1728 int64
	_ = v1728
	var v1731 int64
	_ = v1731
	var v1733 int64
	_ = v1733
	var v1734 int64
	_ = v1734
	var v1745 int32
	_ = v1745
	var v1747 int32
	_ = v1747
	var v1752 int32
	_ = v1752
	var v1764 int64
	_ = v1764
	var v1767 int64
	_ = v1767
	var v1768 int64
	_ = v1768
	var v1771 int64
	_ = v1771
	var v1774 int64
	_ = v1774
	var v1776 int64
	_ = v1776
	var v1777 int64
	_ = v1777
	var v1780 int32
	_ = v1780
	var v1783 int32
	_ = v1783
	var v1792 int64
	_ = v1792
	var v1795 int64
	_ = v1795
	var v1798 int64
	_ = v1798
	var v1800 int64
	_ = v1800
	var v1801 int64
	_ = v1801
	var v1812 int32
	_ = v1812
	var v1814 int32
	_ = v1814
	var v1819 int32
	_ = v1819
	var v1831 int64
	_ = v1831
	var v1834 int64
	_ = v1834
	var v1835 int64
	_ = v1835
	var v1838 int64
	_ = v1838
	var v1841 int64
	_ = v1841
	var v1843 int64
	_ = v1843
	var v1844 int64
	_ = v1844
	var v1847 int32
	_ = v1847
	var v1850 int32
	_ = v1850
	var v1859 int64
	_ = v1859
	var v1862 int64
	_ = v1862
	var v1865 int64
	_ = v1865
	var v1867 int64
	_ = v1867
	var v1868 int64
	_ = v1868
	var v1879 int32
	_ = v1879
	var v1881 int32
	_ = v1881
	var v1886 int32
	_ = v1886
	var v1898 int64
	_ = v1898
	var v1901 int64
	_ = v1901
	var v1902 int64
	_ = v1902
	var v1905 int64
	_ = v1905
	var v1908 int64
	_ = v1908
	var v1910 int64
	_ = v1910
	var v1911 int64
	_ = v1911
	var v1914 int32
	_ = v1914
	var v1917 int32
	_ = v1917
	var v1926 int64
	_ = v1926
	var v1929 int64
	_ = v1929
	var v1932 int64
	_ = v1932
	var v1934 int64
	_ = v1934
	var v1935 int64
	_ = v1935
	var v1946 int32
	_ = v1946
	var v1948 int32
	_ = v1948
	var v1953 int32
	_ = v1953
	var v1965 int64
	_ = v1965
	var v1968 int64
	_ = v1968
	var v1969 int64
	_ = v1969
	var v1972 int64
	_ = v1972
	var v1975 int64
	_ = v1975
	var v1977 int64
	_ = v1977
	var v1978 int64
	_ = v1978
	var v1981 int32
	_ = v1981
	var v1984 int32
	_ = v1984
	var v1993 int64
	_ = v1993
	var v1996 int64
	_ = v1996
	var v1999 int64
	_ = v1999
	var v2001 int64
	_ = v2001
	var v2002 int64
	_ = v2002
	var v2013 int32
	_ = v2013
	var v2015 int32
	_ = v2015
	var v2020 int32
	_ = v2020
	var v2032 int64
	_ = v2032
	var v2035 int64
	_ = v2035
	var v2036 int64
	_ = v2036
	var v2039 int64
	_ = v2039
	var v2042 int64
	_ = v2042
	var v2044 int64
	_ = v2044
	var v2045 int64
	_ = v2045
	var v2048 int32
	_ = v2048
	var v2051 int32
	_ = v2051
	var v2060 int64
	_ = v2060
	var v2063 int64
	_ = v2063
	var v2066 int64
	_ = v2066
	var v2068 int64
	_ = v2068
	var v2069 int64
	_ = v2069
	var v2080 int32
	_ = v2080
	var v2082 int32
	_ = v2082
	var v2087 int32
	_ = v2087
	var v2099 int64
	_ = v2099
	var v2102 int64
	_ = v2102
	var v2103 int64
	_ = v2103
	var v2106 int64
	_ = v2106
	var v2109 int64
	_ = v2109
	var v2111 int64
	_ = v2111
	var v2112 int64
	_ = v2112
	var v2115 int32
	_ = v2115
	var v2118 int32
	_ = v2118
	var v2127 int64
	_ = v2127
	var v2130 int64
	_ = v2130
	var v2133 int64
	_ = v2133
	var v2135 int64
	_ = v2135
	var v2136 int64
	_ = v2136
	var v2141 int64
	_ = v2141
	var v2145 int32
	_ = v2145
	var v2147 int32
	_ = v2147
	var v2156 int64
	_ = v2156
	var v2157 int64
	_ = v2157
	var v2159 int64
	_ = v2159
	var v2168 int64
	_ = v2168
	var v2170 int64
	_ = v2170
	var v2172 int64
	_ = v2172
	var v2174 int64
	_ = v2174
	var v2178 int32
	_ = v2178
	var v2179 int32
	_ = v2179
	var v2182 int64
	_ = v2182
	var v2184 int64
	_ = v2184
	var v2189 int32
	_ = v2189
	var v2195 int64
	_ = v2195
	var v2197 int32
	_ = v2197
	var v2206 int64
	_ = v2206
	var v2207 int64
	_ = v2207
	var v2209 int32
	_ = v2209
	var v2210 int32
	_ = v2210
	var v2211 int32
	_ = v2211
	var v2216 int32
	_ = v2216
	var v2217 int32
	_ = v2217
	var v2218 int32
	_ = v2218
	var v2220 int32
	_ = v2220
	var v2229 int32
	_ = v2229
	var v2239 int32
	_ = v2239
	var v2241 int32
	_ = v2241
	var v2277 int32
	_ = v2277
	var v2278 int32
	_ = v2278
	var v2279 int32
	_ = v2279
	var v2281 int32
	_ = v2281
	var v2282 int32
	_ = v2282
	var v2285 int32
	_ = v2285
	var v2287 int32
	_ = v2287
	var v2289 int32
	_ = v2289
	var v2292 int32
	_ = v2292
	var v2294 int32
	_ = v2294
	var v2296 int32
	_ = v2296
	var v2299 int32
	_ = v2299
	var v2301 int32
	_ = v2301
	var v2306 int32
	_ = v2306
	var v2311 int32
	_ = v2311
	var v2332 int32
	_ = v2332
	var v2335 int32
	_ = v2335
	var v2366 int32
	_ = v2366
	var v2369 int32
	_ = v2369
	var v2372 int32
	_ = v2372
	var v2373 float32
	_ = v2373
	var v2403 int32
	_ = v2403
	var v2408 int32
	_ = v2408
	var v2410 int32
	_ = v2410
	var v2417 int32
	_ = v2417
	var v2420 int32
	_ = v2420
	var v2427 int32
	_ = v2427
	var v2428 int32
	_ = v2428
	var v2429 int32
	_ = v2429
	var v2430 int32
	_ = v2430
	var v2436 int32
	_ = v2436
	var v2437 int32
	_ = v2437
	var v2441 int32
	_ = v2441
	var v2442 int32
	_ = v2442
	var v2443 int32
	_ = v2443
	var v2447 int32
	_ = v2447
	var v2450 int32
	_ = v2450
	var v2451 int32
	_ = v2451
	var v2466 int32
	_ = v2466
	var v2476 int32
	_ = v2476
	var v2516 int32
	_ = v2516
	var v2527 int32
	_ = v2527
	var v2540 int32
	_ = v2540
	var v2580 int32
	_ = v2580
	var v2586 int32
	_ = v2586
	var v2604 int32
	_ = v2604
	var v2613 int32
	_ = v2613
	var v2670 int32
	_ = v2670
	var v2683 int32
	_ = v2683
	var v2718 int32
	_ = v2718
	var v2722 int32
	_ = v2722
	var v2729 int32
	_ = v2729
	var v2787 int32
	_ = v2787
	var v2804 int32
	_ = v2804
	var v2810 int32
	_ = v2810
	var v2811 int32
	_ = v2811
	var v2812 int32
	_ = v2812
	var v2845 int32
	_ = v2845
	var v2856 int32
	_ = v2856
	var v2859 int32
	_ = v2859
	var v2868 int32
	_ = v2868
	var v2871 int32
	_ = v2871
	var v2909 int32
	_ = v2909
	var v2910 int32
	_ = v2910
	var v2911 int32
	_ = v2911
	var v2942 int32
	_ = v2942
	var v2943 int32
	_ = v2943
	var v2944 int32
	_ = v2944
	var v2945 int32
	_ = v2945
	var v2946 int32
	_ = v2946
	var v2950 int32
	_ = v2950
	var v3010 int32
	_ = v3010
	var v3015 int32
	_ = v3015
	var v3019 int32
	_ = v3019
	var v3021 int32
	_ = v3021
	var v3024 int32
	_ = v3024
	var v3025 int64
	_ = v3025
	var v3036 int32
	_ = v3036
	var v3039 int32
	_ = v3039
	var v3042 int32
	_ = v3042
	var v3045 int32
	_ = v3045
	var v3046 int64
	_ = v3046
	var v3057 int32
	_ = v3057
	var v3060 int32
	_ = v3060
	var v3063 int32
	_ = v3063
	var v3066 int32
	_ = v3066
	var v3067 int64
	_ = v3067
	var v3078 int32
	_ = v3078
	var v3081 int32
	_ = v3081
	var v3084 int32
	_ = v3084
	var v3087 int32
	_ = v3087
	var v3088 int64
	_ = v3088
	var v3099 int32
	_ = v3099
	var v3102 int32
	_ = v3102
	var v3104 int32
	_ = v3104
	var v3109 int32
	_ = v3109
	var v3120 int32
	_ = v3120
	var v3123 int32
	_ = v3123
	var v3126 int32
	_ = v3126
	var v3128 int32
	_ = v3128
	var v3129 int32
	_ = v3129
	var v3141 int32
	_ = v3141
	var v3147 int32
	_ = v3147
	var v3152 int32
	_ = v3152
	var v3154 int32
	_ = v3154
	var v3157 int32
	_ = v3157
	var v3168 int32
	_ = v3168
	var v3171 int32
	_ = v3171
	var v3172 int32
	_ = v3172
	var v3174 int32
	_ = v3174
	var v3185 int32
	_ = v3185
	var v3186 int32
	_ = v3186
	var v3189 int32
	_ = v3189
	var v3198 int32
	_ = v3198
	var v3201 int32
	_ = v3201
	var v3205 int32
	_ = v3205
	var v3207 int32
	_ = v3207
	var v3208 int32
	_ = v3208
	var v3211 int32
	_ = v3211
	var v3213 int32
	_ = v3213
	var v3215 int32
	_ = v3215
	var v3220 int32
	_ = v3220
	var v3221 int32
	_ = v3221
	var v3224 int32
	_ = v3224
	var v3227 int32
	_ = v3227
	var v3230 int32
	_ = v3230
	var v3234 int32
	_ = v3234
	var v3237 int32
	_ = v3237
	var v3243 int32
	_ = v3243
	var v3244 int32
	_ = v3244
	var v3246 int32
	_ = v3246
	var v3254 int32
	_ = v3254
	var v3255 int32
	_ = v3255
	var v3256 int32
	_ = v3256
	var v3259 int32
	_ = v3259
	var v3272 int32
	_ = v3272
	var v3277 int32
	_ = v3277
	var v3278 int32
	_ = v3278
	var v3284 int32
	_ = v3284
	var v3294 int32
	_ = v3294
	var v3295 int32
	_ = v3295
	var v3296 int32
	_ = v3296
	var v3298 int32
	_ = v3298
	var v3299 int32
	_ = v3299
	var v3300 int32
	_ = v3300
	var v3302 int64
	_ = v3302
	var v3308 int32
	_ = v3308
	var v3309 int32
	_ = v3309
	var v3310 int32
	_ = v3310
	var v3317 int32
	_ = v3317
	var v3319 int32
	_ = v3319
	var v3323 int32
	_ = v3323
	var v3324 int32
	_ = v3324
	var v3329 int32
	_ = v3329
	var v3330 int32
	_ = v3330
	var v3331 int32
	_ = v3331
	var v3338 int32
	_ = v3338
	var v3339 int32
	_ = v3339
	var v3340 int32
	_ = v3340
	var v3341 int32
	_ = v3341
	var v3342 int64
	_ = v3342
	var v3344 int32
	_ = v3344
	var v3346 int32
	_ = v3346
	var v3354 int32
	_ = v3354
	var v3356 int32
	_ = v3356
	var v3362 int64
	_ = v3362
	var v3363 int32
	_ = v3363
	var v3370 int64
	_ = v3370
	var v3371 int32
	_ = v3371
	var v3372 int64
	_ = v3372
	var v3381 int32
	_ = v3381
	var v3383 int32
	_ = v3383
	var v3386 int32
	_ = v3386
	var v3388 int32
	_ = v3388
	var v3399 int32
	_ = v3399
	var v3400 int32
	_ = v3400
	var v3401 int32
	_ = v3401
	var v3402 int32
	_ = v3402
	var v3403 int32
	_ = v3403
	var v3407 int32
	_ = v3407
	var v3415 int32
	_ = v3415
	var v3419 int32
	_ = v3419
	var v3423 int32
	_ = v3423
	var v3425 int32
	_ = v3425
	var v3428 int32
	_ = v3428
	var v3429 int64
	_ = v3429
	var v3440 int32
	_ = v3440
	var v3443 int32
	_ = v3443
	var v3446 int32
	_ = v3446
	var v3449 int32
	_ = v3449
	var v3450 int64
	_ = v3450
	var v3461 int32
	_ = v3461
	var v3464 int32
	_ = v3464
	var v3467 int32
	_ = v3467
	var v3470 int32
	_ = v3470
	var v3471 int64
	_ = v3471
	var v3482 int32
	_ = v3482
	var v3485 int32
	_ = v3485
	var v3488 int32
	_ = v3488
	var v3491 int32
	_ = v3491
	var v3492 int64
	_ = v3492
	var v3503 int32
	_ = v3503
	var v3506 int32
	_ = v3506
	var v3508 int64
	_ = v3508
	var v3510 int32
	_ = v3510
	var v3512 int32
	_ = v3512
	var v3514 int32
	_ = v3514
	var v3518 int32
	_ = v3518
	var v3519 int32
	_ = v3519
	var v3523 int32
	_ = v3523
	var v3527 int32
	_ = v3527
	var v3531 int32
	_ = v3531
	var v3538 int32
	_ = v3538
	var v3544 int32
	_ = v3544
	var v3545 int32
	_ = v3545
	var v3547 int32
	_ = v3547
	var v3550 int32
	_ = v3550
	var v3556 int32
	_ = v3556
	var v3559 int32
	_ = v3559
	var v3563 int32
	_ = v3563
	var v3564 int32
	_ = v3564
	var v3568 int32
	_ = v3568
	var v3569 int32
	_ = v3569
	var v3572 int32
	_ = v3572
	var v3577 int32
	_ = v3577
	var v3578 int32
	_ = v3578
	var v3579 int32
	_ = v3579
	var v3582 int32
	_ = v3582
	var v3586 int32
	_ = v3586
	var v3587 int32
	_ = v3587
	var v3588 int32
	_ = v3588
	var v3591 int32
	_ = v3591
	var v3592 int32
	_ = v3592
	var v3594 int32
	_ = v3594
	var v3595 int32
	_ = v3595
	var v3597 int32
	_ = v3597
	var v3600 int32
	_ = v3600
	var v3601 int32
	_ = v3601
	var v3603 int32
	_ = v3603
	var v3611 int32
	_ = v3611
	var v3620 int64
	_ = v3620
	var v3621 int64
	_ = v3621
	var v3624 int32
	_ = v3624
	var v3625 int64
	_ = v3625
	var v3627 int32
	_ = v3627
	var v3628 int64
	_ = v3628
	var v3631 int32
	_ = v3631
	var v3632 int64
	_ = v3632
	var v3634 int32
	_ = v3634
	var v3635 int64
	_ = v3635
	var v3640 int32
	_ = v3640
	var v3646 int32
	_ = v3646
	var v3652 int32
	_ = v3652
	var v3701 int32
	_ = v3701
	var v3706 int32
	_ = v3706
	var v3707 int32
	_ = v3707
	var v3711 int64
	_ = v3711
	var v3719 int32
	_ = v3719
	var v3720 int32
	_ = v3720
	var v3725 int32
	_ = v3725
	var v3728 int32
	_ = v3728
	var v3731 int32
	_ = v3731
	var v3733 int32
	_ = v3733
	var v3735 int32
	_ = v3735
	var v3746 int32
	_ = v3746
	var v3747 int32
	_ = v3747
	var v3754 int32
	_ = v3754
	var v3757 int32
	_ = v3757
	var v3760 int32
	_ = v3760
	var v3762 int32
	_ = v3762
	var v3764 int32
	_ = v3764
	var v3775 int32
	_ = v3775
	var v3776 int32
	_ = v3776
	var v3783 int32
	_ = v3783
	var v3786 int32
	_ = v3786
	var v3789 int32
	_ = v3789
	var v3791 int32
	_ = v3791
	var v3793 int32
	_ = v3793
	var v3804 int32
	_ = v3804
	var v3805 int32
	_ = v3805
	var v3812 int32
	_ = v3812
	var v3815 int32
	_ = v3815
	var v3818 int32
	_ = v3818
	var v3820 int32
	_ = v3820
	var v3822 int32
	_ = v3822
	var v3833 int32
	_ = v3833
	var v3834 int32
	_ = v3834
	var v3840 int32
	_ = v3840
	var v3848 int32
	_ = v3848
	var v3849 int32
	_ = v3849
	var v3854 int32
	_ = v3854
	var v3857 int32
	_ = v3857
	var v3860 int32
	_ = v3860
	var v3862 int32
	_ = v3862
	var v3864 int32
	_ = v3864
	var v3875 int32
	_ = v3875
	var v3876 int32
	_ = v3876
	var v3883 int32
	_ = v3883
	var v3886 int32
	_ = v3886
	var v3889 int32
	_ = v3889
	var v3891 int32
	_ = v3891
	var v3893 int32
	_ = v3893
	var v3904 int32
	_ = v3904
	var v3905 int32
	_ = v3905
	var v3912 int32
	_ = v3912
	var v3915 int32
	_ = v3915
	var v3918 int32
	_ = v3918
	var v3920 int32
	_ = v3920
	var v3922 int32
	_ = v3922
	var v3933 int32
	_ = v3933
	var v3934 int32
	_ = v3934
	var v3941 int32
	_ = v3941
	var v3944 int32
	_ = v3944
	var v3947 int32
	_ = v3947
	var v3949 int32
	_ = v3949
	var v3951 int32
	_ = v3951
	var v3962 int32
	_ = v3962
	var v3963 int32
	_ = v3963
	var v3969 int32
	_ = v3969
	var v3974 int32
	_ = v3974
	var v3982 int32
	_ = v3982
	v54 = m.G0
	v56 = v54 - int32(1728)
	m.G0 = v56
	v58 = int64(1)
	v59 = int32(2224)
	goto L8
L1:
	;
	m.G0 = v56 + int32(1728)
	return v3982
L2:
	;
	v262 = v56 + int32(8)
	goto L62
L3:
	;
	v258 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v258 != 0 {
		goto L57
	} else {
		goto L58
	}
L4:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+12)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v79)+4)) = l1
	*(*int32)(unsafe.Add(mBase, uint32(v79))) = l0
	F_VP8LEncDspInit(m)
	mBase = m.M
	v94 = m.G1
	goto L16
L5:
	;
	if v79 != 0 {
		goto L4
	} else {
		goto L11
	}
L6:
	;
	goto L5
L7:
	;
	v77 = F_calloc(m, base.I32_wrap_i64(v58), v59)
	mBase = m.M
	v79 = v77
	goto L6
L8:
	;
	v66 = base.I64_div_u_s(int64(2147418112), v58)
	v67 = int32(0)
	v68 = base.I64_extend_i32_u(v59)
	if base.Ui64(int64(4294967295)) < base.Ui64(v68*v58) {
		v79 = v67
		goto L6
	} else {
		goto L9
	}
L9:
	;
	if base.Ui64(v66) < base.Ui64(v68) {
		v79 = v67
		goto L6
	} else {
		goto L10
	}
L10:
	;
	goto L7
L11:
	;
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v82 != 0 {
		goto L13
	} else {
		goto L14
	}
L12:
	;
	goto L15
L13:
	;
	goto L12
L14:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L13
L15:
	;
	goto L3
L16:
	;
	v98 = v56 + int32(180)
	v101 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v98))) = v101
	*(*int64)(unsafe.Add(mBase, uint32(v56+int32(196)))) = v101
	*(*int64)(unsafe.Add(mBase, uint32(v56+int32(188)))) = v101
	v114 = int32(1024)
	v116 = F_WebPSafeMalloc(m, int64(1), v114)
	mBase = m.M
	if v116 != 0 {
		goto L18
	} else {
		goto L19
	}
L17:
	;
	if v116 != 0 {
		goto L2
	} else {
		goto L20
	}
L18:
	;
	v120 = *(*int32)(unsafe.Add(mBase, uint32(v98)+8))
	F_WebPSafeFree(m, v120)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v98)+16)) = v116 + v114
	*(*int32)(unsafe.Add(mBase, uint32(v98)+12)) = v116
	*(*int32)(unsafe.Add(mBase, uint32(v98)+8)) = v116
	goto L17
L19:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v98)+20)) = int32(1)
	goto L17
L20:
	;
	v129 = v79 + int32(2216)
	v130 = *(*int32)(unsafe.Add(mBase, uint32(v129)))
	F_WebPSafeFree(m, v130)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v129))) = int64(0)
	goto L21
L21:
	;
	v135 = v79 + int32(2120)
	v138 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
	if v138 == int32(0) {
		goto L23
	} else {
		goto L24
	}
L22:
	;
	v164 = v79 + int32(2144)
	v167 = *(*int32)(unsafe.Add(mBase, uint32(v164)+12))
	if v167 == int32(0) {
		goto L31
	} else {
		goto L32
	}
L23:
	;
	v143 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v135)+20)) = v143
	v145 = *(*int32)(unsafe.Add(mBase, uint32(v135)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v135)+16)) = v145
	*(*int32)(unsafe.Add(mBase, uint32(v135)+12)) = v79 + int32(2128)
	*(*int32)(unsafe.Add(mBase, uint32(v135)+8)) = v143
	if v145 == v143 {
		goto L25
	} else {
		goto L26
	}
L24:
	;
	v141 = *(*int32)(unsafe.Add(mBase, uint32(v135)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v138))) = v141
	goto L23
L25:
	;
	goto L22
L26:
	;
	v156 = v145
	goto L27
L27:
	;
	v157 = *(*int32)(unsafe.Add(mBase, uint32(v156)))
	F_WebPSafeFree(m, v156)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v135)+16)) = v157
	if v157 != 0 {
		v156 = v157
		goto L27
	} else {
		goto L29
	}
L28:
	;
	goto L25
L29:
	;
	goto L28
L30:
	;
	v193 = v79 + int32(2168)
	v196 = *(*int32)(unsafe.Add(mBase, uint32(v193)+12))
	if v196 == int32(0) {
		goto L39
	} else {
		goto L40
	}
L31:
	;
	v172 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v164)+20)) = v172
	v174 = *(*int32)(unsafe.Add(mBase, uint32(v164)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v164)+16)) = v174
	*(*int32)(unsafe.Add(mBase, uint32(v164)+12)) = v79 + int32(2152)
	*(*int32)(unsafe.Add(mBase, uint32(v164)+8)) = v172
	if v174 == v172 {
		goto L33
	} else {
		goto L34
	}
L32:
	;
	v170 = *(*int32)(unsafe.Add(mBase, uint32(v164)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v167))) = v170
	goto L31
L33:
	;
	goto L30
L34:
	;
	v185 = v174
	goto L35
L35:
	;
	v186 = *(*int32)(unsafe.Add(mBase, uint32(v185)))
	F_WebPSafeFree(m, v185)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v164)+16)) = v186
	if v186 != 0 {
		v185 = v186
		goto L35
	} else {
		goto L37
	}
L36:
	;
	goto L33
L37:
	;
	goto L36
L38:
	;
	v222 = v79 + int32(2192)
	v225 = *(*int32)(unsafe.Add(mBase, uint32(v222)+12))
	if v225 == int32(0) {
		goto L47
	} else {
		goto L48
	}
L39:
	;
	v201 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v193)+20)) = v201
	v203 = *(*int32)(unsafe.Add(mBase, uint32(v193)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v193)+16)) = v203
	*(*int32)(unsafe.Add(mBase, uint32(v193)+12)) = v79 + int32(2176)
	*(*int32)(unsafe.Add(mBase, uint32(v193)+8)) = v201
	if v203 == v201 {
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v199 = *(*int32)(unsafe.Add(mBase, uint32(v193)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v196))) = v199
	goto L39
L41:
	;
	goto L38
L42:
	;
	v214 = v203
	goto L43
L43:
	;
	v215 = *(*int32)(unsafe.Add(mBase, uint32(v214)))
	F_WebPSafeFree(m, v214)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v193)+16)) = v215
	if v215 != 0 {
		v214 = v215
		goto L43
	} else {
		goto L45
	}
L44:
	;
	goto L41
L45:
	;
	goto L44
L46:
	;
	v250 = *(*int32)(unsafe.Add(mBase, uint32(v79)+24))
	F_free(m, v250)
	mBase = m.M
	goto L54
L47:
	;
	v230 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v222)+20)) = v230
	v232 = *(*int32)(unsafe.Add(mBase, uint32(v222)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v222)+16)) = v232
	*(*int32)(unsafe.Add(mBase, uint32(v222)+12)) = v79 + int32(2200)
	*(*int32)(unsafe.Add(mBase, uint32(v222)+8)) = v230
	if v232 == v230 {
		goto L49
	} else {
		goto L50
	}
L48:
	;
	v228 = *(*int32)(unsafe.Add(mBase, uint32(v222)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v225))) = v228
	goto L47
L49:
	;
	goto L46
L50:
	;
	v243 = v232
	goto L51
L51:
	;
	v244 = *(*int32)(unsafe.Add(mBase, uint32(v243)))
	F_WebPSafeFree(m, v243)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v222)+16)) = v244
	if v244 != 0 {
		v243 = v244
		goto L51
	} else {
		goto L53
	}
L52:
	;
	goto L49
L53:
	;
	goto L52
L54:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v79)+24)) = int64(0)
	F_free(m, v79)
	mBase = m.M
	goto L55
L55:
	;
	goto L3
L56:
	;
	v3982 = int32(0)
	goto L1
L57:
	;
	goto L56
L58:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L57
L59:
	;
	v3701 = v56 + int32(180)
	if v3701 == int32(0) {
		goto L589
	} else {
		goto L590
	}
L60:
	;
	goto L64
L61:
	;
	goto L60
L62:
	;
	if v262 == int32(0) {
		goto L61
	} else {
		goto L63
	}
L63:
	;
	v273 = int32(0)
	v275 = F_memset(m, v262, v273, int32(172))
	mBase = m.M
	v276 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v275)+72)) = v276 + int32(1)
	v281 = F_WebPEncodingSetError(m, v275, v273)
	mBase = m.M
	goto L61
L64:
	;
	v286 = *(*int32)(unsafe.Add(mBase, uint32(v79)+4))
	v287 = *(*int32)(unsafe.Add(mBase, uint32(v286)+12))
	v288 = *(*int32)(unsafe.Add(mBase, uint32(v286)+8))
	v289 = *(*int32)(unsafe.Add(mBase, uint32(v79)))
	v290 = *(*int32)(unsafe.Add(mBase, uint32(v289)+8))
	v292 = v79 + int32(1096)
	v293 = int32(0)
	v303 = m.G0
	v305 = v303 - int32(_a_F_VP8LEncodeStream_0)
	m.G0 = v305
	v308 = int32(_a_F_VP8LEncodeStream_1)
	v312 = F_memset(m, v305+v308, v293, int32(1024))
	mBase = m.M
	v315 = F_memset(m, v305, v293, v308)
	mBase = m.M
	v316 = *(*int32)(unsafe.Add(mBase, uint32(v286)+12))
	if v316 < int32(1) {
		v452 = v293
		goto L68
	} else {
		goto L69
	}
L65:
	;
	v535 = base.B2i32(v521 < int32(257))
	if v521 < int32(257) {
		goto L95
	} else {
		goto L96
	}
L66:
	;
	m.G0 = v315 + int32(_a_F_VP8LEncodeStream_0)
	goto L65
L67:
	;
	v521 = int32(257)
	goto L66
L68:
	;
	if v292 == int32(0) {
		v521 = v452
		goto L66
	} else {
		goto L87
	}
L69:
	;
	v319 = *(*int32)(unsafe.Add(mBase, uint32(v286)+8))
	if v319 < int32(1) {
		v452 = v293
		goto L68
	} else {
		goto L70
	}
L70:
	;
	v322 = *(*int32)(unsafe.Add(mBase, uint32(v286)+52))
	v323 = *(*int32)(unsafe.Add(mBase, uint32(v322)))
	v326 = int32(0)
	v330 = v323 ^ int32(-1)
	v331 = v326
	v335 = v322
	v336 = v326
	goto L71
L71:
	;
	v343 = v330
	v344 = v331
	v350 = int32(0)
	goto L73
L72:
	;
	v452 = v430
	goto L68
L73:
	;
	v356 = *(*int32)(unsafe.Add(mBase, uint32(v335+v350<<(uint(int32(2))%32))))
	if v356 == v343 {
		v429 = v343
		v430 = v344
		goto L75
	} else {
		goto L76
	}
L74:
	;
	v442 = *(*int32)(unsafe.Add(mBase, uint32(v286)+56))
	v447 = v336 + int32(1)
	if v447 != v316 {
		v330 = v429
		v331 = v430
		v335 = v335 + v442<<(uint(int32(2))%32)
		v336 = v447
		goto L71
	} else {
		goto L86
	}
L75:
	;
	v440 = v350 + int32(1)
	if v440 != v319 {
		v343 = v429
		v344 = v430
		v350 = v440
		goto L73
	} else {
		goto L85
	}
L76:
	;
	v363 = int32(base.Ui32(v356*int32(506832829)) >> (uint(int32(22)) % 32))
	v364 = v315 + int32(_a_F_VP8LEncodeStream_1) + v363
	v365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v364))))
	if v365 == int32(0) {
		v395 = v363
		v404 = v364
		goto L78
	} else {
		goto L79
	}
L77:
	;
	v429 = v356
	v430 = v418
	goto L75
L78:
	;
	v405 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v404))) = uint8(v405)
	*(*int32)(unsafe.Add(mBase, uint32(v315+v395<<(uint(int32(2))%32)))) = v356
	if int32(255) < v344 {
		goto L67
	} else {
		goto L84
	}
L79:
	;
	v370 = v363
	goto L80
L80:
	;
	v383 = *(*int32)(unsafe.Add(mBase, uint32(v315+v370<<(uint(int32(2))%32))))
	if v383 == v356 {
		v418 = v344
		goto L77
	} else {
		goto L82
	}
L81:
	;
	v395 = v390
	v404 = v391
	goto L78
L82:
	;
	v390 = (v370 + int32(1)) & int32(1023)
	v391 = v315 + int32(_a_F_VP8LEncodeStream_1) + v390
	v392 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v391))))
	if v392 != 0 {
		v370 = v390
		goto L80
	} else {
		goto L83
	}
L83:
	;
	goto L81
L84:
	;
	v418 = v344 + int32(1)
	goto L77
L85:
	;
	goto L74
L86:
	;
	goto L72
L87:
	;
	v463 = int32(0)
	v467 = v315
	v468 = v463
	v475 = v463
	goto L88
L88:
	;
	v479 = v315 + int32(_a_F_VP8LEncodeStream_1) + v475
	v480 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v479))))
	if v480 == int32(0) {
		v490 = v468
		goto L90
	} else {
		goto L91
	}
L89:
	;
	v513 = m.G2
	F_qsort(m, v292, v505, int32(4), v513+int32(170))
	mBase = m.M
	v521 = v505
	goto L66
L90:
	;
	v493 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v479+int32(1)))))
	if v493 == int32(0) {
		v505 = v490
		goto L92
	} else {
		goto L93
	}
L91:
	;
	v486 = *(*int32)(unsafe.Add(mBase, uint32(v467)))
	*(*int32)(unsafe.Add(mBase, uint32(v292+v468<<(uint(int32(2))%32)))) = v486
	v490 = v468 + int32(1)
	goto L90
L92:
	;
	v509 = v475 + int32(2)
	if v509 != int32(1024) {
		v467 = v467 + int32(8)
		v468 = v505
		v475 = v509
		goto L88
	} else {
		goto L94
	}
L93:
	;
	v501 = *(*int32)(unsafe.Add(mBase, uint32(v467+int32(4))))
	*(*int32)(unsafe.Add(mBase, uint32(v292+v490<<(uint(int32(2))%32)))) = v501
	v505 = v490 + int32(1)
	goto L92
L94:
	;
	goto L89
L95:
	;
	v536 = v521
	goto L97
L96:
	;
	v536 = int32(0)
	goto L97
L97:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+68)) = v536
	if v521 < int32(257) {
		goto L98
	} else {
		goto L99
	}
L98:
	;
	v542 = int32(9)
	goto L100
L99:
	;
	v542 = int32(7)
	goto L100
L100:
	;
	v543 = v542 - v290
	v544 = int32(9)
	if base.Ui32(v543) < base.Ui32(v544) {
		goto L101
	} else {
		goto L102
	}
L101:
	;
	v547 = v543
	goto L103
L102:
	;
	v547 = v544
	goto L103
L103:
	;
	if v543 < int32(2) {
		goto L104
	} else {
		goto L105
	}
L104:
	;
	v550 = int32(2)
	goto L106
L105:
	;
	v550 = v547
	goto L106
L106:
	;
	v551 = int32(1) << (uint(v550) % 32)
	v552 = *(*int32)(unsafe.Add(mBase, uint32(v286)+8))
	v553 = int32(-1)
	v554 = v552 + v553
	v557 = *(*int32)(unsafe.Add(mBase, uint32(v286)+12))
	v559 = v557 + v553
	v562 = int32(base.Ui32(v551+v554)>>(uint(v550)%32)) * int32(base.Ui32(v551+v559)>>(uint(v550)%32))
	if base.Ui32(int32(8)) < base.Ui32(v550) {
		goto L109
	} else {
		goto L110
	}
L107:
	;
	if base.Ui32(int32(3)) <= base.Ui32(v648) {
		goto L117
	} else {
		goto L118
	}
L108:
	;
	v572 = v550
	goto L112
L109:
	;
	v648 = v550
	v651 = v562
	goto L107
L110:
	;
	if int32(2601) <= v562 {
		goto L108
	} else {
		goto L111
	}
L111:
	;
	goto L109
L112:
	;
	v621 = int32(2) << (uint(v572) % 32)
	v624 = v572 + int32(1)
	v628 = int32(base.Ui32(v621+v554)>>(uint(v624)%32)) * int32(base.Ui32(v621+v559)>>(uint(v624)%32))
	if base.Ui32(int32(7)) < base.Ui32(v572) {
		v648 = v624
		v651 = v628
		goto L107
	} else {
		goto L114
	}
L113:
	;
	v648 = v624
	v651 = v628
	goto L107
L114:
	;
	if int32(2600) < v628 {
		v572 = v624
		goto L112
	} else {
		goto L115
	}
L115:
	;
	goto L113
L116:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+36)) = v788
	v824 = int32(4)
	if v290 == v824 {
		goto L128
	} else {
		goto L129
	}
L117:
	;
	if v651 == int32(1) {
		goto L119
	} else {
		goto L120
	}
L118:
	;
	v788 = v648
	goto L116
L119:
	;
	v690 = int32(1)
	v692 = v648 + int32(-1)
	v693 = v690 << (uint(v692) % 32)
	if int32(base.Ui32(v693+v554)>>(uint(v692)%32))*int32(base.Ui32(v693+v559)>>(uint(v692)%32)) == v690 {
		goto L121
	} else {
		goto L122
	}
L120:
	;
	v788 = v648
	goto L116
L121:
	;
	v716 = v648
	goto L123
L122:
	;
	v788 = v648
	goto L116
L123:
	;
	v755 = v716 + int32(-1)
	if int32(3) <= v755 {
		goto L125
	} else {
		goto L126
	}
L124:
	;
	v788 = v755
	goto L116
L125:
	;
	v760 = v716 + int32(-2)
	v761 = int32(1)
	v762 = v761 << (uint(v760) % 32)
	if int32(base.Ui32(v762+v554)>>(uint(v760)%32))*int32(base.Ui32(v762+v559)>>(uint(v760)%32)) == v761 {
		v716 = v755
		goto L123
	} else {
		goto L127
	}
L126:
	;
	v788 = int32(2)
	goto L116
L127:
	;
	goto L124
L128:
	;
	v830 = int32(5)
	goto L130
L129:
	;
	v830 = v824
	goto L130
L130:
	;
	if v290 < int32(4) {
		goto L131
	} else {
		goto L132
	}
L131:
	;
	v833 = int32(6)
	goto L133
L132:
	;
	v833 = v830
	goto L133
L133:
	;
	if base.Ui32(v833) < base.Ui32(v788) {
		goto L134
	} else {
		goto L135
	}
L134:
	;
	v835 = v833
	goto L136
L135:
	;
	v835 = v788
	goto L136
L136:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+44)) = v835
	*(*int32)(unsafe.Add(mBase, uint32(v79)+40)) = v835
	if v290 != 0 {
		goto L141
	} else {
		goto L142
	}
L137:
	;
	v3015 = v79 + int32(2120)
	v3019 = base.I32_div_s(v2946+int32(-1), int32(16))
	v3021 = v3019 + int32(1)
	v3024 = v79 + int32(2132)
	v3025 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3024))) = v3025
	*(*int64)(unsafe.Add(mBase, uint32(v3015)+4)) = v3025
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2140)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3024))) = v79 + int32(2128)
	v3036 = int32(256)
	if v3036 < v3021 {
		goto L461
	} else {
		goto L462
	}
L138:
	;
	v3010 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3010 != 0 {
		goto L458
	} else {
		goto L459
	}
L139:
	;
	v2942 = v79 + int32(2216)
	v2943 = *(*int32)(unsafe.Add(mBase, uint32(v79)+4))
	v2944 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+12))
	v2945 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+8))
	v2946 = v2944 * v2945
	v2950 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v2946), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v2942))) = v2950
	if v2950 != 0 {
		goto L454
	} else {
		goto L455
	}
L140:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1292)) = v2804
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1288)) = int32(3)
	v2845 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1304)) = v2845
	if v2810 == v2845 {
		v2909 = v2810
		v2910 = v2811
		v2911 = v2812
		goto L139
	} else {
		goto L450
	}
L141:
	;
	v853 = base.B2i32(v536 < int32(17))
	v854 = *(*int32)(unsafe.Add(mBase, uint32(v286)+56))
	v855 = *(*int32)(unsafe.Add(mBase, uint32(v286)+52))
	if int32(256) < v521 {
		goto L150
	} else {
		goto L151
	}
L142:
	;
	v838 = int32(0)
	v842 = base.B2i32(v521 < int32(257))
	if v521 < int32(257) {
		goto L143
	} else {
		goto L144
	}
L143:
	;
	v843 = v838
	goto L145
L144:
	;
	v843 = int32(3)
	goto L145
L145:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1284)) = v843
	if v521 < int32(257) {
		goto L146
	} else {
		goto L147
	}
L146:
	;
	v847 = int32(4)
	goto L148
L147:
	;
	v847 = int32(3)
	goto L148
L148:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1280)) = v847
	v849 = int32(1)
	v2804 = v838
	v2810 = v849
	v2811 = int32(0)
	v2812 = v849
	goto L140
L149:
	;
	v2366 = int32(1)
	if v536 < int32(17) {
		goto L414
	} else {
		goto L415
	}
L150:
	;
	v860 = int64(1)
	v861 = int32(_a_F_VP8LEncodeStream_2)
	goto L156
L151:
	;
	if v536 < int32(17) {
		v2332 = v824
		v2335 = int32(1)
		goto L149
	} else {
		goto L152
	}
L152:
	;
	goto L150
L153:
	;
	if v881 == int32(0) {
		goto L138
	} else {
		goto L159
	}
L154:
	;
	goto L153
L155:
	;
	v879 = F_calloc(m, base.I32_wrap_i64(v860), v861)
	mBase = m.M
	v881 = v879
	goto L154
L156:
	;
	v868 = base.I64_div_u_s(int64(2147418112), v860)
	v869 = int32(0)
	v870 = base.I64_extend_i32_u(v861)
	if base.Ui64(int64(4294967295)) < base.Ui64(v870*v860) {
		v881 = v869
		goto L154
	} else {
		goto L157
	}
L157:
	;
	if base.Ui64(v868) < base.Ui64(v870) {
		v881 = v869
		goto L154
	} else {
		goto L158
	}
L158:
	;
	goto L155
L159:
	;
	if v287 < int32(1) {
		goto L160
	} else {
		goto L161
	}
L160:
	;
	v1247 = *(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[0])))
	v1248 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[0]))) = v1247 + v1248
	v1251 = *(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[1])))
	*(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[1]))) = v1251 + v1248
	v1255 = *(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[2]))) = v1255 + v1248
	v1259 = *(*int32)(unsafe.Add(mBase, uint32(v881)+3072))
	*(*int32)(unsafe.Add(mBase, uint32(v881)+3072)) = v1259 + v1248
	v1263 = *(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v881)+uint32(_c_F_VP8LEncodeStream[3]))) = v1263 + v1248
	v1267 = *(*int32)(unsafe.Add(mBase, uint32(v881)+1024))
	*(*int32)(unsafe.Add(mBase, uint32(v881)+1024)) = v1267 + v1248
	v1271 = m.G61
	v1276 = m.G0
	v1278 = v1276 - int32(32)
	m.G0 = v1278
	F_VP8LBitsEntropyUnrefined(m, v881, int32(256), v1278+int32(8))
	mBase = m.M
	v1283 = *(*int32)(unsafe.Add(mBase, uint32(v1278)+20))
	if v1283 <= int32(4) {
		goto L177
	} else {
		goto L178
	}
L161:
	;
	if v288 < int32(1) {
		goto L160
	} else {
		goto L162
	}
L162:
	;
	v915 = *(*int32)(unsafe.Add(mBase, uint32(v855)))
	v916 = int32(0)
	__phi923 = v915
	__phi933 = v855
	__phi954 = v916
	__phi955 = v916
	v923 = __phi923
	v933 = __phi933
	v954 = __phi954
	v955 = __phi955
	goto L163
L163:
	;
	v977 = v923
	v987 = int32(0)
	v988 = v288
	goto L165
L164:
	;
	goto L160
L165:
	;
	v1026 = *(*int32)(unsafe.Add(mBase, uint32(v933+v987)))
	v1029 = int32(-16711936)
	v1031 = v1026 | int32(16711680) - v977&v1029
	v1036 = int32(16711935)
	v1038 = v1026 | int32(_a_F_VP8LEncodeStream_3) - v977&v1036
	v1041 = v1031&v1029 | v1038&v1036
	if v1041 == int32(0) {
		goto L167
	} else {
		goto L168
	}
L166:
	;
	v1192 = v954 + int32(1)
	if v1192 != v287 {
		__phi923 = v1026
		__phi933 = v933 + v854<<(uint(int32(2))%32)
		__phi954 = v1192
		__phi955 = v933
		v923 = __phi923
		v933 = __phi933
		v954 = __phi954
		v955 = __phi955
		goto L163
	} else {
		goto L173
	}
L167:
	;
	v1189 = v988 + int32(-1)
	if v1189 != 0 {
		v977 = v1026
		v987 = v987 + int32(4)
		v988 = v1189
		goto L165
	} else {
		goto L172
	}
L168:
	;
	if v955 == int32(0) {
		goto L169
	} else {
		goto L170
	}
L169:
	;
	v1049 = int32(22)
	v1051 = int32(1020)
	v1053 = v881 + int32(base.Ui32(v1026)>>(uint(v1049)%32))&v1051
	v1054 = *(*int32)(unsafe.Add(mBase, uint32(v1053)))
	v1055 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v1053))) = v1054 + v1055
	v1058 = int32(14)
	v1062 = v881 + int32(_a_F_VP8LEncodeStream_1) + int32(base.Ui32(v1026)>>(uint(v1058)%32))&v1051
	v1063 = *(*int32)(unsafe.Add(mBase, uint32(v1062)))
	*(*int32)(unsafe.Add(mBase, uint32(v1062))) = v1063 + v1055
	v1067 = int32(6)
	v1071 = v881 + int32(2048) + int32(base.Ui32(v1026)>>(uint(v1067)%32))&v1051
	v1072 = *(*int32)(unsafe.Add(mBase, uint32(v1071)))
	*(*int32)(unsafe.Add(mBase, uint32(v1071))) = v1072 + v1055
	v1076 = int32(255)
	v1078 = int32(2)
	v1080 = v881 + int32(_a_F_VP8LEncodeStream_4) + v1026&v1076<<(uint(v1078)%32)
	v1081 = *(*int32)(unsafe.Add(mBase, uint32(v1080)))
	*(*int32)(unsafe.Add(mBase, uint32(v1080))) = v1081 + v1055
	v1089 = v881 + int32(1024) + int32(base.Ui32(v1031)>>(uint(v1049)%32))&v1051
	v1090 = *(*int32)(unsafe.Add(mBase, uint32(v1089)))
	*(*int32)(unsafe.Add(mBase, uint32(v1089))) = v1090 + v1055
	v1098 = v881 + int32(_a_F_VP8LEncodeStream_0) + int32(base.Ui32(v1041)>>(uint(v1058)%32))&v1051
	v1099 = *(*int32)(unsafe.Add(mBase, uint32(v1098)))
	*(*int32)(unsafe.Add(mBase, uint32(v1098))) = v1099 + v1055
	v1107 = v881 + int32(3072) + int32(base.Ui32(v1031)>>(uint(v1067)%32))&v1051
	v1108 = *(*int32)(unsafe.Add(mBase, uint32(v1107)))
	*(*int32)(unsafe.Add(mBase, uint32(v1107))) = v1108 + v1055
	v1116 = v881 + int32(_a_F_VP8LEncodeStream_5) + v1038&v1076<<(uint(v1078)%32)
	v1117 = *(*int32)(unsafe.Add(mBase, uint32(v1116)))
	*(*int32)(unsafe.Add(mBase, uint32(v1116))) = v1117 + v1055
	v1121 = int32(16)
	v1123 = int32(8)
	v1124 = v1026 >> (uint(v1123) % 32)
	v1130 = v881 + int32(_a_F_VP8LEncodeStream_6) + (int32(base.Ui32(v1026)>>(uint(v1121)%32))-v1124)&v1076<<(uint(v1078)%32)
	v1131 = *(*int32)(unsafe.Add(mBase, uint32(v1130)))
	*(*int32)(unsafe.Add(mBase, uint32(v1130))) = v1131 + v1055
	v1140 = v881 + int32(_a_F_VP8LEncodeStream_7) + (v1026-v1124)&v1076<<(uint(v1078)%32)
	v1141 = *(*int32)(unsafe.Add(mBase, uint32(v1140)))
	*(*int32)(unsafe.Add(mBase, uint32(v1140))) = v1141 + v1055
	v1148 = v1041 >> (uint(v1123) % 32)
	v1154 = v881 + int32(_a_F_VP8LEncodeStream_8) + (int32(base.Ui32(v1041)>>(uint(v1121)%32))-v1148)&v1076<<(uint(v1078)%32)
	v1155 = *(*int32)(unsafe.Add(mBase, uint32(v1154)))
	*(*int32)(unsafe.Add(mBase, uint32(v1154))) = v1155 + v1055
	v1164 = v881 + int32(_a_F_VP8LEncodeStream_9) + (v1038-v1148)&v1076<<(uint(v1078)%32)
	v1165 = *(*int32)(unsafe.Add(mBase, uint32(v1164)))
	*(*int32)(unsafe.Add(mBase, uint32(v1164))) = v1165 + v1055
	v1178 = v881 + int32(_a_F_VP8LEncodeStream_10) + int32(base.Ui32((int32(base.Ui32(v1026)>>(uint(int32(19))%32))+v1026)*int32(969276327))>>(uint(v1049)%32))&v1051
	v1179 = *(*int32)(unsafe.Add(mBase, uint32(v1178)))
	*(*int32)(unsafe.Add(mBase, uint32(v1178))) = v1179 + v1055
	goto L167
L170:
	;
	v1047 = *(*int32)(unsafe.Add(mBase, uint32(v955+v987)))
	if v1026 == v1047 {
		goto L167
	} else {
		goto L171
	}
L171:
	;
	goto L169
L172:
	;
	goto L166
L173:
	;
	goto L164
L174:
	;
	v1343 = m.G0
	v1345 = v1343 - int32(32)
	m.G0 = v1345
	F_VP8LBitsEntropyUnrefined(m, v881+int32(1024), int32(256), v1345+int32(8))
	mBase = m.M
	v1350 = *(*int32)(unsafe.Add(mBase, uint32(v1345)+20))
	if v1350 <= int32(4) {
		goto L194
	} else {
		goto L195
	}
L175:
	;
	m.G0 = v1278 + int32(32)
	goto L174
L176:
	;
	v1308 = *(*int64)(unsafe.Add(mBase, uint32(v1278)+8))
	v1311 = *(*int32)(unsafe.Add(mBase, uint32(v1278)+16))
	v1314 = *(*int32)(unsafe.Add(mBase, uint32(v1278)+24))
	v1323 = v1307*base.I64_extend_i32_u(v1311<<(uint(int32(1))%32)-v1314)<<(uint(int64(23))%64) + v1308*(int64(1000)-v1307)
	if v1323 < int64(0) {
		goto L185
	} else {
		goto L186
	}
L177:
	;
	if v1283 < int32(2) {
		v1332 = int64(0)
		goto L175
	} else {
		goto L179
	}
L178:
	;
	v1307 = int64(627)
	goto L176
L179:
	;
	switch v1283 + int32(-2) {
	case 0:
		goto L181
	case 1:
		v1307 = int64(950)
		goto L176
	default:
		goto L180
	}
L180:
	;
	v1307 = int64(700)
	goto L176
L181:
	;
	v1295 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1278)+16)))
	v1298 = *(*int64)(unsafe.Add(mBase, uint32(v1278)+8))
	v1299 = v1295*int64(830472192) + v1298
	if v1299 < int64(0) {
		goto L182
	} else {
		goto L183
	}
L182:
	;
	v1302 = int64(-50)
	goto L184
L183:
	;
	v1302 = int64(50)
	goto L184
L184:
	;
	v1305 = base.I64_div_s(v1302+v1299, int64(100))
	v1332 = v1305
	goto L175
L185:
	;
	v1326 = int64(-500)
	goto L187
L186:
	;
	v1326 = int64(500)
	goto L187
L187:
	;
	v1329 = base.I64_div_s(v1326+v1323, int64(1000))
	if base.Ui64(v1329) < base.Ui64(v1308) {
		goto L188
	} else {
		goto L189
	}
L188:
	;
	v1331 = v1308
	goto L190
L189:
	;
	v1331 = v1329
	goto L190
L190:
	;
	v1332 = v1331
	goto L175
L191:
	;
	v1410 = m.G0
	v1412 = v1410 - int32(32)
	m.G0 = v1412
	F_VP8LBitsEntropyUnrefined(m, v881+int32(2048), int32(256), v1412+int32(8))
	mBase = m.M
	v1417 = *(*int32)(unsafe.Add(mBase, uint32(v1412)+20))
	if v1417 <= int32(4) {
		goto L211
	} else {
		goto L212
	}
L192:
	;
	m.G0 = v1345 + int32(32)
	goto L191
L193:
	;
	v1375 = *(*int64)(unsafe.Add(mBase, uint32(v1345)+8))
	v1378 = *(*int32)(unsafe.Add(mBase, uint32(v1345)+16))
	v1381 = *(*int32)(unsafe.Add(mBase, uint32(v1345)+24))
	v1390 = v1374*base.I64_extend_i32_u(v1378<<(uint(int32(1))%32)-v1381)<<(uint(int64(23))%64) + v1375*(int64(1000)-v1374)
	if v1390 < int64(0) {
		goto L202
	} else {
		goto L203
	}
L194:
	;
	if v1350 < int32(2) {
		v1399 = int64(0)
		goto L192
	} else {
		goto L196
	}
L195:
	;
	v1374 = int64(627)
	goto L193
L196:
	;
	switch v1350 + int32(-2) {
	case 0:
		goto L198
	case 1:
		v1374 = int64(950)
		goto L193
	default:
		goto L197
	}
L197:
	;
	v1374 = int64(700)
	goto L193
L198:
	;
	v1362 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1345)+16)))
	v1365 = *(*int64)(unsafe.Add(mBase, uint32(v1345)+8))
	v1366 = v1362*int64(830472192) + v1365
	if v1366 < int64(0) {
		goto L199
	} else {
		goto L200
	}
L199:
	;
	v1369 = int64(-50)
	goto L201
L200:
	;
	v1369 = int64(50)
	goto L201
L201:
	;
	v1372 = base.I64_div_s(v1369+v1366, int64(100))
	v1399 = v1372
	goto L192
L202:
	;
	v1393 = int64(-500)
	goto L204
L203:
	;
	v1393 = int64(500)
	goto L204
L204:
	;
	v1396 = base.I64_div_s(v1393+v1390, int64(1000))
	if base.Ui64(v1396) < base.Ui64(v1375) {
		goto L205
	} else {
		goto L206
	}
L205:
	;
	v1398 = v1375
	goto L207
L206:
	;
	v1398 = v1396
	goto L207
L207:
	;
	v1399 = v1398
	goto L192
L208:
	;
	v1477 = m.G0
	v1479 = v1477 - int32(32)
	m.G0 = v1479
	F_VP8LBitsEntropyUnrefined(m, v881+int32(3072), int32(256), v1479+int32(8))
	mBase = m.M
	v1484 = *(*int32)(unsafe.Add(mBase, uint32(v1479)+20))
	if v1484 <= int32(4) {
		goto L228
	} else {
		goto L229
	}
L209:
	;
	m.G0 = v1412 + int32(32)
	goto L208
L210:
	;
	v1442 = *(*int64)(unsafe.Add(mBase, uint32(v1412)+8))
	v1445 = *(*int32)(unsafe.Add(mBase, uint32(v1412)+16))
	v1448 = *(*int32)(unsafe.Add(mBase, uint32(v1412)+24))
	v1457 = v1441*base.I64_extend_i32_u(v1445<<(uint(int32(1))%32)-v1448)<<(uint(int64(23))%64) + v1442*(int64(1000)-v1441)
	if v1457 < int64(0) {
		goto L219
	} else {
		goto L220
	}
L211:
	;
	if v1417 < int32(2) {
		v1466 = int64(0)
		goto L209
	} else {
		goto L213
	}
L212:
	;
	v1441 = int64(627)
	goto L210
L213:
	;
	switch v1417 + int32(-2) {
	case 0:
		goto L215
	case 1:
		v1441 = int64(950)
		goto L210
	default:
		goto L214
	}
L214:
	;
	v1441 = int64(700)
	goto L210
L215:
	;
	v1429 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1412)+16)))
	v1432 = *(*int64)(unsafe.Add(mBase, uint32(v1412)+8))
	v1433 = v1429*int64(830472192) + v1432
	if v1433 < int64(0) {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	v1436 = int64(-50)
	goto L218
L217:
	;
	v1436 = int64(50)
	goto L218
L218:
	;
	v1439 = base.I64_div_s(v1436+v1433, int64(100))
	v1466 = v1439
	goto L209
L219:
	;
	v1460 = int64(-500)
	goto L221
L220:
	;
	v1460 = int64(500)
	goto L221
L221:
	;
	v1463 = base.I64_div_s(v1460+v1457, int64(1000))
	if base.Ui64(v1463) < base.Ui64(v1442) {
		goto L222
	} else {
		goto L223
	}
L222:
	;
	v1465 = v1442
	goto L224
L223:
	;
	v1465 = v1463
	goto L224
L224:
	;
	v1466 = v1465
	goto L209
L225:
	;
	v1544 = m.G0
	v1546 = v1544 - int32(32)
	m.G0 = v1546
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_1), int32(256), v1546+int32(8))
	mBase = m.M
	v1551 = *(*int32)(unsafe.Add(mBase, uint32(v1546)+20))
	if v1551 <= int32(4) {
		goto L245
	} else {
		goto L246
	}
L226:
	;
	m.G0 = v1479 + int32(32)
	goto L225
L227:
	;
	v1509 = *(*int64)(unsafe.Add(mBase, uint32(v1479)+8))
	v1512 = *(*int32)(unsafe.Add(mBase, uint32(v1479)+16))
	v1515 = *(*int32)(unsafe.Add(mBase, uint32(v1479)+24))
	v1524 = v1508*base.I64_extend_i32_u(v1512<<(uint(int32(1))%32)-v1515)<<(uint(int64(23))%64) + v1509*(int64(1000)-v1508)
	if v1524 < int64(0) {
		goto L236
	} else {
		goto L237
	}
L228:
	;
	if v1484 < int32(2) {
		v1533 = int64(0)
		goto L226
	} else {
		goto L230
	}
L229:
	;
	v1508 = int64(627)
	goto L227
L230:
	;
	switch v1484 + int32(-2) {
	case 0:
		goto L232
	case 1:
		v1508 = int64(950)
		goto L227
	default:
		goto L231
	}
L231:
	;
	v1508 = int64(700)
	goto L227
L232:
	;
	v1496 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1479)+16)))
	v1499 = *(*int64)(unsafe.Add(mBase, uint32(v1479)+8))
	v1500 = v1496*int64(830472192) + v1499
	if v1500 < int64(0) {
		goto L233
	} else {
		goto L234
	}
L233:
	;
	v1503 = int64(-50)
	goto L235
L234:
	;
	v1503 = int64(50)
	goto L235
L235:
	;
	v1506 = base.I64_div_s(v1503+v1500, int64(100))
	v1533 = v1506
	goto L226
L236:
	;
	v1527 = int64(-500)
	goto L238
L237:
	;
	v1527 = int64(500)
	goto L238
L238:
	;
	v1530 = base.I64_div_s(v1527+v1524, int64(1000))
	if base.Ui64(v1530) < base.Ui64(v1509) {
		goto L239
	} else {
		goto L240
	}
L239:
	;
	v1532 = v1509
	goto L241
L240:
	;
	v1532 = v1530
	goto L241
L241:
	;
	v1533 = v1532
	goto L226
L242:
	;
	v1611 = m.G0
	v1613 = v1611 - int32(32)
	m.G0 = v1613
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_0), int32(256), v1613+int32(8))
	mBase = m.M
	v1618 = *(*int32)(unsafe.Add(mBase, uint32(v1613)+20))
	if v1618 <= int32(4) {
		goto L262
	} else {
		goto L263
	}
L243:
	;
	m.G0 = v1546 + int32(32)
	goto L242
L244:
	;
	v1576 = *(*int64)(unsafe.Add(mBase, uint32(v1546)+8))
	v1579 = *(*int32)(unsafe.Add(mBase, uint32(v1546)+16))
	v1582 = *(*int32)(unsafe.Add(mBase, uint32(v1546)+24))
	v1591 = v1575*base.I64_extend_i32_u(v1579<<(uint(int32(1))%32)-v1582)<<(uint(int64(23))%64) + v1576*(int64(1000)-v1575)
	if v1591 < int64(0) {
		goto L253
	} else {
		goto L254
	}
L245:
	;
	if v1551 < int32(2) {
		v1600 = int64(0)
		goto L243
	} else {
		goto L247
	}
L246:
	;
	v1575 = int64(627)
	goto L244
L247:
	;
	switch v1551 + int32(-2) {
	case 0:
		goto L249
	case 1:
		v1575 = int64(950)
		goto L244
	default:
		goto L248
	}
L248:
	;
	v1575 = int64(700)
	goto L244
L249:
	;
	v1563 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1546)+16)))
	v1566 = *(*int64)(unsafe.Add(mBase, uint32(v1546)+8))
	v1567 = v1563*int64(830472192) + v1566
	if v1567 < int64(0) {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v1570 = int64(-50)
	goto L252
L251:
	;
	v1570 = int64(50)
	goto L252
L252:
	;
	v1573 = base.I64_div_s(v1570+v1567, int64(100))
	v1600 = v1573
	goto L243
L253:
	;
	v1594 = int64(-500)
	goto L255
L254:
	;
	v1594 = int64(500)
	goto L255
L255:
	;
	v1597 = base.I64_div_s(v1594+v1591, int64(1000))
	if base.Ui64(v1597) < base.Ui64(v1576) {
		goto L256
	} else {
		goto L257
	}
L256:
	;
	v1599 = v1576
	goto L258
L257:
	;
	v1599 = v1597
	goto L258
L258:
	;
	v1600 = v1599
	goto L243
L259:
	;
	v1678 = m.G0
	v1680 = v1678 - int32(32)
	m.G0 = v1680
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_4), int32(256), v1680+int32(8))
	mBase = m.M
	v1685 = *(*int32)(unsafe.Add(mBase, uint32(v1680)+20))
	if v1685 <= int32(4) {
		goto L279
	} else {
		goto L280
	}
L260:
	;
	m.G0 = v1613 + int32(32)
	goto L259
L261:
	;
	v1643 = *(*int64)(unsafe.Add(mBase, uint32(v1613)+8))
	v1646 = *(*int32)(unsafe.Add(mBase, uint32(v1613)+16))
	v1649 = *(*int32)(unsafe.Add(mBase, uint32(v1613)+24))
	v1658 = v1642*base.I64_extend_i32_u(v1646<<(uint(int32(1))%32)-v1649)<<(uint(int64(23))%64) + v1643*(int64(1000)-v1642)
	if v1658 < int64(0) {
		goto L270
	} else {
		goto L271
	}
L262:
	;
	if v1618 < int32(2) {
		v1667 = int64(0)
		goto L260
	} else {
		goto L264
	}
L263:
	;
	v1642 = int64(627)
	goto L261
L264:
	;
	switch v1618 + int32(-2) {
	case 0:
		goto L266
	case 1:
		v1642 = int64(950)
		goto L261
	default:
		goto L265
	}
L265:
	;
	v1642 = int64(700)
	goto L261
L266:
	;
	v1630 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1613)+16)))
	v1633 = *(*int64)(unsafe.Add(mBase, uint32(v1613)+8))
	v1634 = v1630*int64(830472192) + v1633
	if v1634 < int64(0) {
		goto L267
	} else {
		goto L268
	}
L267:
	;
	v1637 = int64(-50)
	goto L269
L268:
	;
	v1637 = int64(50)
	goto L269
L269:
	;
	v1640 = base.I64_div_s(v1637+v1634, int64(100))
	v1667 = v1640
	goto L260
L270:
	;
	v1661 = int64(-500)
	goto L272
L271:
	;
	v1661 = int64(500)
	goto L272
L272:
	;
	v1664 = base.I64_div_s(v1661+v1658, int64(1000))
	if base.Ui64(v1664) < base.Ui64(v1643) {
		goto L273
	} else {
		goto L274
	}
L273:
	;
	v1666 = v1643
	goto L275
L274:
	;
	v1666 = v1664
	goto L275
L275:
	;
	v1667 = v1666
	goto L260
L276:
	;
	v1745 = m.G0
	v1747 = v1745 - int32(32)
	m.G0 = v1747
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_5), int32(256), v1747+int32(8))
	mBase = m.M
	v1752 = *(*int32)(unsafe.Add(mBase, uint32(v1747)+20))
	if v1752 <= int32(4) {
		goto L296
	} else {
		goto L297
	}
L277:
	;
	m.G0 = v1680 + int32(32)
	goto L276
L278:
	;
	v1710 = *(*int64)(unsafe.Add(mBase, uint32(v1680)+8))
	v1713 = *(*int32)(unsafe.Add(mBase, uint32(v1680)+16))
	v1716 = *(*int32)(unsafe.Add(mBase, uint32(v1680)+24))
	v1725 = v1709*base.I64_extend_i32_u(v1713<<(uint(int32(1))%32)-v1716)<<(uint(int64(23))%64) + v1710*(int64(1000)-v1709)
	if v1725 < int64(0) {
		goto L287
	} else {
		goto L288
	}
L279:
	;
	if v1685 < int32(2) {
		v1734 = int64(0)
		goto L277
	} else {
		goto L281
	}
L280:
	;
	v1709 = int64(627)
	goto L278
L281:
	;
	switch v1685 + int32(-2) {
	case 0:
		goto L283
	case 1:
		v1709 = int64(950)
		goto L278
	default:
		goto L282
	}
L282:
	;
	v1709 = int64(700)
	goto L278
L283:
	;
	v1697 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1680)+16)))
	v1700 = *(*int64)(unsafe.Add(mBase, uint32(v1680)+8))
	v1701 = v1697*int64(830472192) + v1700
	if v1701 < int64(0) {
		goto L284
	} else {
		goto L285
	}
L284:
	;
	v1704 = int64(-50)
	goto L286
L285:
	;
	v1704 = int64(50)
	goto L286
L286:
	;
	v1707 = base.I64_div_s(v1704+v1701, int64(100))
	v1734 = v1707
	goto L277
L287:
	;
	v1728 = int64(-500)
	goto L289
L288:
	;
	v1728 = int64(500)
	goto L289
L289:
	;
	v1731 = base.I64_div_s(v1728+v1725, int64(1000))
	if base.Ui64(v1731) < base.Ui64(v1710) {
		goto L290
	} else {
		goto L291
	}
L290:
	;
	v1733 = v1710
	goto L292
L291:
	;
	v1733 = v1731
	goto L292
L292:
	;
	v1734 = v1733
	goto L277
L293:
	;
	v1812 = m.G0
	v1814 = v1812 - int32(32)
	m.G0 = v1814
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_6), int32(256), v1814+int32(8))
	mBase = m.M
	v1819 = *(*int32)(unsafe.Add(mBase, uint32(v1814)+20))
	if v1819 <= int32(4) {
		goto L313
	} else {
		goto L314
	}
L294:
	;
	m.G0 = v1747 + int32(32)
	goto L293
L295:
	;
	v1777 = *(*int64)(unsafe.Add(mBase, uint32(v1747)+8))
	v1780 = *(*int32)(unsafe.Add(mBase, uint32(v1747)+16))
	v1783 = *(*int32)(unsafe.Add(mBase, uint32(v1747)+24))
	v1792 = v1776*base.I64_extend_i32_u(v1780<<(uint(int32(1))%32)-v1783)<<(uint(int64(23))%64) + v1777*(int64(1000)-v1776)
	if v1792 < int64(0) {
		goto L304
	} else {
		goto L305
	}
L296:
	;
	if v1752 < int32(2) {
		v1801 = int64(0)
		goto L294
	} else {
		goto L298
	}
L297:
	;
	v1776 = int64(627)
	goto L295
L298:
	;
	switch v1752 + int32(-2) {
	case 0:
		goto L300
	case 1:
		v1776 = int64(950)
		goto L295
	default:
		goto L299
	}
L299:
	;
	v1776 = int64(700)
	goto L295
L300:
	;
	v1764 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1747)+16)))
	v1767 = *(*int64)(unsafe.Add(mBase, uint32(v1747)+8))
	v1768 = v1764*int64(830472192) + v1767
	if v1768 < int64(0) {
		goto L301
	} else {
		goto L302
	}
L301:
	;
	v1771 = int64(-50)
	goto L303
L302:
	;
	v1771 = int64(50)
	goto L303
L303:
	;
	v1774 = base.I64_div_s(v1771+v1768, int64(100))
	v1801 = v1774
	goto L294
L304:
	;
	v1795 = int64(-500)
	goto L306
L305:
	;
	v1795 = int64(500)
	goto L306
L306:
	;
	v1798 = base.I64_div_s(v1795+v1792, int64(1000))
	if base.Ui64(v1798) < base.Ui64(v1777) {
		goto L307
	} else {
		goto L308
	}
L307:
	;
	v1800 = v1777
	goto L309
L308:
	;
	v1800 = v1798
	goto L309
L309:
	;
	v1801 = v1800
	goto L294
L310:
	;
	v1879 = m.G0
	v1881 = v1879 - int32(32)
	m.G0 = v1881
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_8), int32(256), v1881+int32(8))
	mBase = m.M
	v1886 = *(*int32)(unsafe.Add(mBase, uint32(v1881)+20))
	if v1886 <= int32(4) {
		goto L330
	} else {
		goto L331
	}
L311:
	;
	m.G0 = v1814 + int32(32)
	goto L310
L312:
	;
	v1844 = *(*int64)(unsafe.Add(mBase, uint32(v1814)+8))
	v1847 = *(*int32)(unsafe.Add(mBase, uint32(v1814)+16))
	v1850 = *(*int32)(unsafe.Add(mBase, uint32(v1814)+24))
	v1859 = v1843*base.I64_extend_i32_u(v1847<<(uint(int32(1))%32)-v1850)<<(uint(int64(23))%64) + v1844*(int64(1000)-v1843)
	if v1859 < int64(0) {
		goto L321
	} else {
		goto L322
	}
L313:
	;
	if v1819 < int32(2) {
		v1868 = int64(0)
		goto L311
	} else {
		goto L315
	}
L314:
	;
	v1843 = int64(627)
	goto L312
L315:
	;
	switch v1819 + int32(-2) {
	case 0:
		goto L317
	case 1:
		v1843 = int64(950)
		goto L312
	default:
		goto L316
	}
L316:
	;
	v1843 = int64(700)
	goto L312
L317:
	;
	v1831 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1814)+16)))
	v1834 = *(*int64)(unsafe.Add(mBase, uint32(v1814)+8))
	v1835 = v1831*int64(830472192) + v1834
	if v1835 < int64(0) {
		goto L318
	} else {
		goto L319
	}
L318:
	;
	v1838 = int64(-50)
	goto L320
L319:
	;
	v1838 = int64(50)
	goto L320
L320:
	;
	v1841 = base.I64_div_s(v1838+v1835, int64(100))
	v1868 = v1841
	goto L311
L321:
	;
	v1862 = int64(-500)
	goto L323
L322:
	;
	v1862 = int64(500)
	goto L323
L323:
	;
	v1865 = base.I64_div_s(v1862+v1859, int64(1000))
	if base.Ui64(v1865) < base.Ui64(v1844) {
		goto L324
	} else {
		goto L325
	}
L324:
	;
	v1867 = v1844
	goto L326
L325:
	;
	v1867 = v1865
	goto L326
L326:
	;
	v1868 = v1867
	goto L311
L327:
	;
	v1946 = m.G0
	v1948 = v1946 - int32(32)
	m.G0 = v1948
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_7), int32(256), v1948+int32(8))
	mBase = m.M
	v1953 = *(*int32)(unsafe.Add(mBase, uint32(v1948)+20))
	if v1953 <= int32(4) {
		goto L347
	} else {
		goto L348
	}
L328:
	;
	m.G0 = v1881 + int32(32)
	goto L327
L329:
	;
	v1911 = *(*int64)(unsafe.Add(mBase, uint32(v1881)+8))
	v1914 = *(*int32)(unsafe.Add(mBase, uint32(v1881)+16))
	v1917 = *(*int32)(unsafe.Add(mBase, uint32(v1881)+24))
	v1926 = v1910*base.I64_extend_i32_u(v1914<<(uint(int32(1))%32)-v1917)<<(uint(int64(23))%64) + v1911*(int64(1000)-v1910)
	if v1926 < int64(0) {
		goto L338
	} else {
		goto L339
	}
L330:
	;
	if v1886 < int32(2) {
		v1935 = int64(0)
		goto L328
	} else {
		goto L332
	}
L331:
	;
	v1910 = int64(627)
	goto L329
L332:
	;
	switch v1886 + int32(-2) {
	case 0:
		goto L334
	case 1:
		v1910 = int64(950)
		goto L329
	default:
		goto L333
	}
L333:
	;
	v1910 = int64(700)
	goto L329
L334:
	;
	v1898 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1881)+16)))
	v1901 = *(*int64)(unsafe.Add(mBase, uint32(v1881)+8))
	v1902 = v1898*int64(830472192) + v1901
	if v1902 < int64(0) {
		goto L335
	} else {
		goto L336
	}
L335:
	;
	v1905 = int64(-50)
	goto L337
L336:
	;
	v1905 = int64(50)
	goto L337
L337:
	;
	v1908 = base.I64_div_s(v1905+v1902, int64(100))
	v1935 = v1908
	goto L328
L338:
	;
	v1929 = int64(-500)
	goto L340
L339:
	;
	v1929 = int64(500)
	goto L340
L340:
	;
	v1932 = base.I64_div_s(v1929+v1926, int64(1000))
	if base.Ui64(v1932) < base.Ui64(v1911) {
		goto L341
	} else {
		goto L342
	}
L341:
	;
	v1934 = v1911
	goto L343
L342:
	;
	v1934 = v1932
	goto L343
L343:
	;
	v1935 = v1934
	goto L328
L344:
	;
	v2013 = m.G0
	v2015 = v2013 - int32(32)
	m.G0 = v2015
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_9), int32(256), v2015+int32(8))
	mBase = m.M
	v2020 = *(*int32)(unsafe.Add(mBase, uint32(v2015)+20))
	if v2020 <= int32(4) {
		goto L364
	} else {
		goto L365
	}
L345:
	;
	m.G0 = v1948 + int32(32)
	goto L344
L346:
	;
	v1978 = *(*int64)(unsafe.Add(mBase, uint32(v1948)+8))
	v1981 = *(*int32)(unsafe.Add(mBase, uint32(v1948)+16))
	v1984 = *(*int32)(unsafe.Add(mBase, uint32(v1948)+24))
	v1993 = v1977*base.I64_extend_i32_u(v1981<<(uint(int32(1))%32)-v1984)<<(uint(int64(23))%64) + v1978*(int64(1000)-v1977)
	if v1993 < int64(0) {
		goto L355
	} else {
		goto L356
	}
L347:
	;
	if v1953 < int32(2) {
		v2002 = int64(0)
		goto L345
	} else {
		goto L349
	}
L348:
	;
	v1977 = int64(627)
	goto L346
L349:
	;
	switch v1953 + int32(-2) {
	case 0:
		goto L351
	case 1:
		v1977 = int64(950)
		goto L346
	default:
		goto L350
	}
L350:
	;
	v1977 = int64(700)
	goto L346
L351:
	;
	v1965 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1948)+16)))
	v1968 = *(*int64)(unsafe.Add(mBase, uint32(v1948)+8))
	v1969 = v1965*int64(830472192) + v1968
	if v1969 < int64(0) {
		goto L352
	} else {
		goto L353
	}
L352:
	;
	v1972 = int64(-50)
	goto L354
L353:
	;
	v1972 = int64(50)
	goto L354
L354:
	;
	v1975 = base.I64_div_s(v1972+v1969, int64(100))
	v2002 = v1975
	goto L345
L355:
	;
	v1996 = int64(-500)
	goto L357
L356:
	;
	v1996 = int64(500)
	goto L357
L357:
	;
	v1999 = base.I64_div_s(v1996+v1993, int64(1000))
	if base.Ui64(v1999) < base.Ui64(v1978) {
		goto L358
	} else {
		goto L359
	}
L358:
	;
	v2001 = v1978
	goto L360
L359:
	;
	v2001 = v1999
	goto L360
L360:
	;
	v2002 = v2001
	goto L345
L361:
	;
	v2080 = m.G0
	v2082 = v2080 - int32(32)
	m.G0 = v2082
	F_VP8LBitsEntropyUnrefined(m, v881+int32(_a_F_VP8LEncodeStream_10), int32(256), v2082+int32(8))
	mBase = m.M
	v2087 = *(*int32)(unsafe.Add(mBase, uint32(v2082)+20))
	if v2087 <= int32(4) {
		goto L381
	} else {
		goto L382
	}
L362:
	;
	m.G0 = v2015 + int32(32)
	goto L361
L363:
	;
	v2045 = *(*int64)(unsafe.Add(mBase, uint32(v2015)+8))
	v2048 = *(*int32)(unsafe.Add(mBase, uint32(v2015)+16))
	v2051 = *(*int32)(unsafe.Add(mBase, uint32(v2015)+24))
	v2060 = v2044*base.I64_extend_i32_u(v2048<<(uint(int32(1))%32)-v2051)<<(uint(int64(23))%64) + v2045*(int64(1000)-v2044)
	if v2060 < int64(0) {
		goto L372
	} else {
		goto L373
	}
L364:
	;
	if v2020 < int32(2) {
		v2069 = int64(0)
		goto L362
	} else {
		goto L366
	}
L365:
	;
	v2044 = int64(627)
	goto L363
L366:
	;
	switch v2020 + int32(-2) {
	case 0:
		goto L368
	case 1:
		v2044 = int64(950)
		goto L363
	default:
		goto L367
	}
L367:
	;
	v2044 = int64(700)
	goto L363
L368:
	;
	v2032 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2015)+16)))
	v2035 = *(*int64)(unsafe.Add(mBase, uint32(v2015)+8))
	v2036 = v2032*int64(830472192) + v2035
	if v2036 < int64(0) {
		goto L369
	} else {
		goto L370
	}
L369:
	;
	v2039 = int64(-50)
	goto L371
L370:
	;
	v2039 = int64(50)
	goto L371
L371:
	;
	v2042 = base.I64_div_s(v2039+v2036, int64(100))
	v2069 = v2042
	goto L362
L372:
	;
	v2063 = int64(-500)
	goto L374
L373:
	;
	v2063 = int64(500)
	goto L374
L374:
	;
	v2066 = base.I64_div_s(v2063+v2060, int64(1000))
	if base.Ui64(v2066) < base.Ui64(v2045) {
		goto L375
	} else {
		goto L376
	}
L375:
	;
	v2068 = v2045
	goto L377
L376:
	;
	v2068 = v2066
	goto L377
L377:
	;
	v2069 = v2068
	goto L362
L378:
	;
	v2141 = v1533 + v1399
	v2145 = int32(1) << (uint(v835) % 32)
	v2147 = int32(-1)
	v2156 = base.I64_extend_i32_u(int32(base.Ui32(v288+v2145+v2147)>>(uint(v835)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v287+v2145+v2147)>>(uint(v835)%32)))
	v2157 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1271)+96)))
	v2159 = v2069 + (v1935 + v2141) + v2156*v2157
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1704)) = v2159
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1712)) = v2136 + base.I64_extend_i32_s(v536)<<(uint(int64(26))%64)
	v2168 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v1271)+56)))
	v2170 = v1801 + (v1667 + v2141) + v2156*v2168
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1688)) = v2170
	v2172 = v1466 + v1332
	v2174 = v1734 + (v1600 + v2172)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1680)) = v2174
	v2178 = base.B2i32(base.Ui64(v2170) < base.Ui64(v2174))
	v2179 = int32(3)
	v2182 = *(*int64)(unsafe.Add(mBase, uint32(v56+int32(1680)|v2178<<(uint(v2179)%32))))
	v2184 = v2002 + (v1868 + v2172)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1696)) = v2184
	if base.Ui64(v2184) < base.Ui64(v2182) {
		goto L395
	} else {
		goto L396
	}
L379:
	;
	m.G0 = v2082 + int32(32)
	goto L378
L380:
	;
	v2112 = *(*int64)(unsafe.Add(mBase, uint32(v2082)+8))
	v2115 = *(*int32)(unsafe.Add(mBase, uint32(v2082)+16))
	v2118 = *(*int32)(unsafe.Add(mBase, uint32(v2082)+24))
	v2127 = v2111*base.I64_extend_i32_u(v2115<<(uint(int32(1))%32)-v2118)<<(uint(int64(23))%64) + v2112*(int64(1000)-v2111)
	if v2127 < int64(0) {
		goto L389
	} else {
		goto L390
	}
L381:
	;
	if v2087 < int32(2) {
		v2136 = int64(0)
		goto L379
	} else {
		goto L383
	}
L382:
	;
	v2111 = int64(627)
	goto L380
L383:
	;
	switch v2087 + int32(-2) {
	case 0:
		goto L385
	case 1:
		v2111 = int64(950)
		goto L380
	default:
		goto L384
	}
L384:
	;
	v2111 = int64(700)
	goto L380
L385:
	;
	v2099 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v2082)+16)))
	v2102 = *(*int64)(unsafe.Add(mBase, uint32(v2082)+8))
	v2103 = v2099*int64(830472192) + v2102
	if v2103 < int64(0) {
		goto L386
	} else {
		goto L387
	}
L386:
	;
	v2106 = int64(-50)
	goto L388
L387:
	;
	v2106 = int64(50)
	goto L388
L388:
	;
	v2109 = base.I64_div_s(v2106+v2103, int64(100))
	v2136 = v2109
	goto L379
L389:
	;
	v2130 = int64(-500)
	goto L391
L390:
	;
	v2130 = int64(500)
	goto L391
L391:
	;
	v2133 = base.I64_div_s(v2130+v2127, int64(1000))
	if base.Ui64(v2133) < base.Ui64(v2112) {
		goto L392
	} else {
		goto L393
	}
L392:
	;
	v2135 = v2112
	goto L394
L393:
	;
	v2135 = v2133
	goto L394
L394:
	;
	v2136 = v2135
	goto L379
L395:
	;
	v2189 = int32(2)
	goto L397
L396:
	;
	v2189 = v2178
	goto L397
L397:
	;
	v2195 = *(*int64)(unsafe.Add(mBase, uint32(v56+int32(1680)+v2189<<(uint(int32(3))%32))))
	if base.Ui64(v2159) < base.Ui64(v2195) {
		goto L398
	} else {
		goto L399
	}
L398:
	;
	v2197 = v2179
	goto L400
L399:
	;
	v2197 = v2189
	goto L400
L400:
	;
	if int32(256) < v521 {
		v2210 = v2197
		goto L401
	} else {
		goto L402
	}
L401:
	;
	v2211 = m.G1
	v2216 = v2211 + int32(_a_F_VP8LEncodeStream_11) + v2210<<(uint(int32(1))%32)
	v2217 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2216))))
	v2218 = int32(10)
	v2220 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2216)+1)))
	v2229 = v2217 << (uint(v2218) % 32)
	v2239 = v2220 << (uint(v2218) % 32)
	v2241 = int32(255)
	goto L407
L402:
	;
	v2206 = *(*int64)(unsafe.Add(mBase, uint32(v56+int32(1680)+v2197<<(uint(int32(3))%32))))
	v2207 = *(*int64)(unsafe.Add(mBase, uint32(v56)+1712))
	if base.Ui64(v2207) < base.Ui64(v2206) {
		goto L403
	} else {
		goto L404
	}
L403:
	;
	v2209 = int32(4)
	goto L405
L404:
	;
	v2209 = v2197
	goto L405
L405:
	;
	v2210 = v2209
	goto L401
L406:
	;
	F_free(m, v881)
	mBase = m.M
	goto L413
L407:
	;
	v2277 = int32(0)
	v2278 = v881 + v2239
	v2279 = int32(4)
	v2281 = *(*int32)(unsafe.Add(mBase, uint32(v2278+v2279)))
	v2282 = v881 + v2229
	v2285 = *(*int32)(unsafe.Add(mBase, uint32(v2282+v2279)))
	if v2281|v2285 != 0 {
		v2311 = v2277
		goto L406
	} else {
		goto L409
	}
L408:
	;
	v2311 = int32(1)
	goto L406
L409:
	;
	v2287 = int32(8)
	v2289 = *(*int32)(unsafe.Add(mBase, uint32(v2278+v2287)))
	v2292 = *(*int32)(unsafe.Add(mBase, uint32(v2282+v2287)))
	if v2289|v2292 != 0 {
		v2311 = v2277
		goto L406
	} else {
		goto L410
	}
L410:
	;
	v2294 = int32(12)
	v2296 = *(*int32)(unsafe.Add(mBase, uint32(v2278+v2294)))
	v2299 = *(*int32)(unsafe.Add(mBase, uint32(v2282+v2294)))
	if v2296|v2299 != 0 {
		v2311 = v2277
		goto L406
	} else {
		goto L411
	}
L411:
	;
	v2301 = int32(12)
	v2306 = v2241 + int32(-3)
	if v2306 != 0 {
		v2229 = v2229 + v2301
		v2239 = v2239 + v2301
		v2241 = v2306
		goto L407
	} else {
		goto L412
	}
L412:
	;
	goto L408
L413:
	;
	v2332 = v2210
	v2335 = v2311
	goto L149
L414:
	;
	v2369 = int32(2)
	goto L416
L415:
	;
	v2369 = v2366
	goto L416
L416:
	;
	if v536 < int32(1) {
		goto L417
	} else {
		goto L418
	}
L417:
	;
	v2372 = v2366
	goto L419
L418:
	;
	v2372 = v2369
	goto L419
L419:
	;
	v2373 = *(*float32)(unsafe.Add(mBase, uint32(v289)+4))
	if v290 != int32(6) {
		goto L423
	} else {
		goto L424
	}
L420:
	;
	if v2372 == int32(1) {
		v2804 = v2441
		v2810 = v2442
		v2811 = v2335
		v2812 = v2443
		goto L140
	} else {
		goto L436
	}
L421:
	;
	v2441 = v2436
	v2442 = v2437
	v2443 = int32(0)
	goto L420
L422:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2428))) = v2430
	*(*int32)(unsafe.Add(mBase, uint32(v2427))) = int32(5)
	v2436 = int32(1)
	v2437 = v2429
	goto L421
L423:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1280)) = v2332
	v2403 = int32(1)
	if v521 < int32(257) {
		goto L428
	} else {
		goto L429
	}
L424:
	;
	if base.F32_ne(v2373, float32(100)) != 0 {
		goto L423
	} else {
		goto L425
	}
L425:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1364)) = int64(12884901891)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1336)) = int64(12884901890)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1308)) = int64(12884901889)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1280)) = int64(12884901888)
	if v521 <= int32(256) {
		goto L426
	} else {
		goto L427
	}
L426:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1448)) = int64(4294967301)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1420)) = int64(8589934596)
	*(*int64)(unsafe.Add(mBase, uint32(v56)+1392)) = int64(4294967300)
	v2427 = v56 + int32(1476)
	v2428 = v56 + int32(1480)
	v2429 = int32(8)
	v2430 = int32(2)
	goto L422
L427:
	;
	v2436 = int32(1)
	v2437 = int32(4)
	goto L421
L428:
	;
	v2408 = v2403
	goto L430
L429:
	;
	v2408 = int32(3)
	goto L430
L430:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1284)) = v2408
	v2410 = int32(0)
	if base.F32_ge(v2373, float32(75)) != 0 {
		goto L431
	} else {
		goto L432
	}
L431:
	;
	if v290 != int32(5) {
		v2441 = v2410
		v2442 = int32(1)
		v2443 = v2403
		goto L420
	} else {
		goto L433
	}
L432:
	;
	v2441 = v2410
	v2442 = int32(1)
	v2443 = v2403
	goto L420
L433:
	;
	v2417 = int32(1)
	if v2332 == int32(4) {
		goto L434
	} else {
		goto L435
	}
L434:
	;
	v2427 = v56 + int32(1308)
	v2428 = v56 + int32(1312)
	v2429 = int32(2)
	v2430 = v2417
	goto L422
L435:
	;
	v2420 = int32(1)
	v2441 = v2420
	v2442 = v2420
	v2443 = v2417
	goto L420
L436:
	;
	v2447 = v2372 + int32(-1)
	v2450 = int32(3)
	v2451 = v2447 & v2450
	v2466 = int32(0)
	v2476 = v56 + int32(1280)
	goto L437
L437:
	;
	v2516 = v56 + int32(1280) + v2466*int32(28)
	*(*int32)(unsafe.Add(mBase, uint32(v2516)+12)) = v2441
	*(*int32)(unsafe.Add(mBase, uint32(v2516)+8)) = int32(3)
	if base.Ui32(v2372+int32(-2)) < base.Ui32(v2450) {
		v2613 = int32(1)
		goto L439
	} else {
		goto L440
	}
L439:
	;
	if v2451 == int32(0) {
		goto L444
	} else {
		goto L445
	}
L440:
	;
	v2527 = v2476
	v2540 = int32(0)
	goto L441
L441:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(44)))) = v2441
	v2580 = int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(40)))) = v2580
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(36)))) = v2441
	v2586 = v2527 + int32(32)
	*(*int32)(unsafe.Add(mBase, uint32(v2586))) = v2580
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(28)))) = v2441
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(24)))) = v2580
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(20)))) = v2441
	*(*int32)(unsafe.Add(mBase, uint32(v2527+int32(16)))) = v2580
	v2604 = v2540 + v2580
	if v2447&int32(-4) != v2604 {
		v2527 = v2586
		v2540 = v2604
		goto L441
	} else {
		goto L443
	}
L442:
	;
	v2613 = v2540 + int32(5)
	goto L439
L443:
	;
	goto L442
L444:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2516)+24)) = v2372
	v2787 = v2466 + int32(1)
	if v2787 != v2442 {
		v2466 = v2787
		v2476 = v2476 + int32(28)
		goto L437
	} else {
		goto L449
	}
L445:
	;
	v2670 = v2613 << (uint(int32(3)) % 32)
	v2683 = v2451
	goto L446
L446:
	;
	v2718 = v2476 + v2670
	*(*int32)(unsafe.Add(mBase, uint32(v2718+int32(12)))) = v2441
	v2722 = int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v2718+v2722))) = int32(4)
	v2729 = v2683 + int32(-1)
	if v2729 != 0 {
		v2670 = v2670 + v2722
		v2683 = v2729
		goto L446
	} else {
		goto L448
	}
L447:
	;
	goto L444
L448:
	;
	goto L447
L449:
	;
	v2909 = v2442
	v2910 = v2335
	v2911 = v2443
	goto L139
L450:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1332)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1320)) = v2804
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1316)) = int32(3)
	if v2810 == int32(2) {
		v2909 = v2810
		v2910 = v2811
		v2911 = v2812
		goto L139
	} else {
		goto L451
	}
L451:
	;
	v2856 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1388)) = v2856
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1376)) = v2804
	v2859 = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1372)) = v2859
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1360)) = v2856
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1348)) = v2804
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1344)) = v2859
	if v2810 == int32(4) {
		v2909 = v2810
		v2910 = v2811
		v2911 = v2812
		goto L139
	} else {
		goto L452
	}
L452:
	;
	v2868 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1500)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1488)) = v2804
	v2871 = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1484)) = v2871
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1472)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1460)) = v2804
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1456)) = v2871
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1444)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1432)) = v2804
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1428)) = v2871
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1416)) = v2868
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1404)) = v2804
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1400)) = v2871
	v2909 = v2810
	v2910 = v2811
	v2911 = v2812
	goto L139
L453:
	;
	if v2950 != 0 {
		goto L137
	} else {
		goto L456
	}
L454:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2942)+4)) = v2946
	goto L453
L455:
	;
	goto L453
L456:
	;
	goto L138
L457:
	;
	v3652 = int32(0)
	goto L59
L458:
	;
	goto L457
L459:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L458
L460:
	;
	v3042 = v79 + int32(2144)
	v3045 = v79 + int32(2156)
	v3046 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3045))) = v3046
	*(*int64)(unsafe.Add(mBase, uint32(v3042)+4)) = v3046
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2164)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3045))) = v79 + int32(2152)
	v3057 = int32(256)
	if v3057 < v3021 {
		goto L465
	} else {
		goto L466
	}
L461:
	;
	v3039 = v3021
	goto L463
L462:
	;
	v3039 = v3036
	goto L463
L463:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3015))) = v3039
	goto L460
L464:
	;
	v3063 = v79 + int32(2168)
	v3066 = v79 + int32(2180)
	v3067 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3066))) = v3067
	*(*int64)(unsafe.Add(mBase, uint32(v3063)+4)) = v3067
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2188)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3066))) = v79 + int32(2176)
	v3078 = int32(256)
	if v3078 < v3021 {
		goto L469
	} else {
		goto L470
	}
L465:
	;
	v3060 = v3021
	goto L467
L466:
	;
	v3060 = v3057
	goto L467
L467:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3042))) = v3060
	goto L464
L468:
	;
	v3084 = v79 + int32(2192)
	v3087 = v79 + int32(2204)
	v3088 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3087))) = v3088
	*(*int64)(unsafe.Add(mBase, uint32(v3084)+4)) = v3088
	*(*int32)(unsafe.Add(mBase, uint32(v79+int32(2212)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3087))) = v79 + int32(2200)
	v3099 = int32(256)
	if v3099 < v3021 {
		goto L473
	} else {
		goto L474
	}
L469:
	;
	v3081 = v3021
	goto L471
L470:
	;
	v3081 = v3078
	goto L471
L471:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3063))) = v3081
	goto L468
L472:
	;
	v3104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	if int32(1) <= v3104 {
		goto L477
	} else {
		goto L478
	}
L473:
	;
	v3102 = v3021
	goto L475
L474:
	;
	v3102 = v3099
	goto L475
L475:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3084))) = v3102
	goto L472
L476:
	;
	v3129 = v2909 - v3128
	if v3129 < int32(1) {
		goto L484
	} else {
		goto L485
	}
L477:
	;
	v3109 = int32(base.Ui32(v2909) >> (uint(int32(1)) % 32))
	if v2911 != 0 {
		goto L479
	} else {
		goto L480
	}
L478:
	;
	v3128 = int32(0)
	goto L476
L479:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+800)) = v3109
	v3128 = v3109
	goto L476
L480:
	;
	v3120 = int32(1)
	if base.Ui32(v3120) < base.Ui32(v3109) {
		goto L481
	} else {
		goto L482
	}
L481:
	;
	v3123 = v3109
	goto L483
L482:
	;
	v3123 = v3120
	goto L483
L483:
	;
	v3126 = F_memcpy(m, v56+int32(408), v56+int32(1280)+v3109*int32(-28)+v2909*int32(28), v3123*int32(28))
	mBase = m.M
	goto L479
L484:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1224)) = v2910
	*(*int32)(unsafe.Add(mBase, uint32(v56)+812)) = l0
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1220)) = v3129
	*(*int32)(unsafe.Add(mBase, uint32(v56)+820)) = l2
	*(*int32)(unsafe.Add(mBase, uint32(v56)+816)) = l1
	v3147 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1228)) = v3147
	*(*int32)(unsafe.Add(mBase, uint32(v56)+824)) = v79
	v3152 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[4])))
	m.T0[v3152].(func(*base.Module, int32))(m, v56+int32(1256))
	mBase = m.M
	v3154 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1272)) = v3154
	v3157 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1264)) = v3157 + int32(168)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1268)) = v56 + int32(812)
	if v3128 == v3154 {
		v3538 = v3154
		goto L486
	} else {
		goto L487
	}
L485:
	;
	v3141 = F_memcpy(m, v56+int32(828), v56+int32(1280), v3128*int32(-28)+v2909*int32(28))
	mBase = m.M
	goto L484
L486:
	;
	if v3128 == int32(0) {
		goto L567
	} else {
		goto L568
	}
L487:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+804)) = v2910
	*(*int32)(unsafe.Add(mBase, uint32(v56)+392)) = l0
	v3168 = int32(0)
	v3171 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	v3172 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	v3174 = v56 + int32(8)
	if l1 == v3168 {
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v3272 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+104)) = v3272
	v3277 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v3277 != 0 {
		goto L509
	} else {
		goto L510
	}
L489:
	;
	goto L488
L490:
	;
	if v3174 == int32(0) {
		goto L489
	} else {
		goto L491
	}
L491:
	;
	v3185 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	if v3185 != 0 {
		goto L492
	} else {
		goto L493
	}
L492:
	;
	v3186 = v3168
	goto L494
L493:
	;
	v3186 = int32(0)
	goto L494
L494:
	;
	if v3185 != 0 {
		goto L495
	} else {
		goto L496
	}
L495:
	;
	v3189 = v3168
	goto L497
L496:
	;
	v3189 = int32(0)
	goto L497
L497:
	;
	if v3186|v3189 < int32(0) {
		goto L489
	} else {
		goto L498
	}
L498:
	;
	if v3171 < int32(1) {
		goto L489
	} else {
		goto L499
	}
L499:
	;
	if v3172 < int32(1) {
		goto L489
	} else {
		goto L500
	}
L500:
	;
	v3198 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	if v3198 < v3186+v3171 {
		goto L489
	} else {
		goto L501
	}
L501:
	;
	v3201 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	if v3201 < v3189+v3172 {
		goto L489
	} else {
		goto L502
	}
L502:
	;
	if l1 == v3174 {
		v3208 = v3185
		goto L503
	} else {
		goto L504
	}
L503:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+12)) = v3172
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+8)) = v3171
	if v3208 != 0 {
		goto L506
	} else {
		goto L507
	}
L504:
	;
	v3205 = F_memcpy(m, v3174, l1, int32(172))
	mBase = m.M
	F_WebPPictureResetBuffers(m, v3205)
	mBase = m.M
	v3207 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v3208 = v3207
	goto L503
L505:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3174+v3256))) = v3254
	*(*int32)(unsafe.Add(mBase, uint32(v3174+v3255))) = v3259
	goto L489
L506:
	;
	v3243 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	v3244 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
	v3246 = int32(2)
	v3254 = v3244
	v3255 = int32(52)
	v3256 = int32(56)
	v3259 = v3243 + v3244*v3189<<(uint(v3246)%32) + v3186<<(uint(v3246)%32)
	goto L505
L507:
	;
	v3211 = *(*int32)(unsafe.Add(mBase, uint32(l1)+32))
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+32)) = v3211
	v3213 = *(*int32)(unsafe.Add(mBase, uint32(l1)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+28)) = v3213
	v3215 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+16)) = v3215 + v3213*v3189 + v3186
	v3220 = int32(1)
	v3221 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	v3224 = v3211 * int32(base.Ui32(v3189)>>(uint(v3220)%32))
	v3227 = int32(base.Ui32(v3186) >> (uint(v3220) % 32))
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+20)) = v3221 + v3224 + v3227
	v3230 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	*(*int32)(unsafe.Add(mBase, uint32(v3174)+24)) = v3230 + v3224 + v3227
	v3234 = *(*int32)(unsafe.Add(mBase, uint32(l1)+36))
	if v3234 == int32(0) {
		goto L489
	} else {
		goto L508
	}
L508:
	;
	v3237 = *(*int32)(unsafe.Add(mBase, uint32(l1)+40))
	v3254 = v3237
	v3255 = int32(36)
	v3256 = int32(40)
	v3259 = v3234 + v3237*v3189 + v3186
	goto L505
L509:
	;
	v3278 = v56 + int32(204)
	goto L511
L510:
	;
	v3278 = v3272
	goto L511
L511:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+808)) = v3278
	*(*int32)(unsafe.Add(mBase, uint32(v56)+396)) = v56 + int32(8)
	v3284 = v56 + int32(180)
	v3294 = *(*int32)(unsafe.Add(mBase, uint32(v3284)+12))
	v3295 = *(*int32)(unsafe.Add(mBase, uint32(v3284)+8))
	v3296 = v3294 - v3295
	v3298 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	v3299 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v3300 = v3298 - v3299
	v3302 = base.I64_extend_i32_u(v3296) + base.I64_extend_i32_u(v3300)
	if base.Ui64(v3302) < base.Ui64(int64(4294967296)) {
		goto L514
	} else {
		goto L515
	}
L512:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v56)+400)) = v56 + int32(180)
	v3362 = int64(1)
	v3363 = int32(2224)
	goto L537
L513:
	;
	if v3354 != 0 {
		goto L512
	} else {
		goto L527
	}
L514:
	;
	v3308 = base.I32_wrap_i64(v3302)
	v3309 = *(*int32)(unsafe.Add(mBase, uint32(v3284)+16))
	v3310 = v3309 - v3295
	if v3309 == v3295 {
		goto L517
	} else {
		goto L518
	}
L515:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+20)) = int32(1)
	v3354 = int32(0)
	goto L513
L516:
	;
	v3341 = F_memcpy(m, v3340, v3339, v3300)
	mBase = m.M
	v3342 = *(*int64)(unsafe.Add(mBase, uint32(l2)))
	*(*int64)(unsafe.Add(mBase, uint32(v3284))) = v3342
	v3344 = *(*int32)(unsafe.Add(mBase, uint32(l2)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+20)) = v3344
	v3346 = *(*int32)(unsafe.Add(mBase, uint32(v3284)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+12)) = v3346 + v3300
	v3354 = int32(1)
	goto L513
L517:
	;
	v3317 = int32(base.Ui32(v3310*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v3308) < base.Ui32(v3317) {
		goto L521
	} else {
		goto L522
	}
L518:
	;
	if base.Ui32(v3310) < base.Ui32(v3308) {
		goto L517
	} else {
		goto L519
	}
L519:
	;
	v3339 = v3299
	v3340 = v3295
	goto L516
L520:
	;
	if v3294 == v3295 {
		goto L525
	} else {
		goto L526
	}
L521:
	;
	v3319 = v3317
	goto L523
L522:
	;
	v3319 = v3308
	goto L523
L523:
	;
	v3323 = v3319&int32(-1024) + int32(1024)
	v3324 = F_WebPSafeMalloc(m, int64(1), v3323)
	mBase = m.M
	if v3324 != 0 {
		goto L520
	} else {
		goto L524
	}
L524:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+20)) = int32(1)
	v3354 = int32(0)
	goto L513
L525:
	;
	v3331 = *(*int32)(unsafe.Add(mBase, uint32(v3284)+8))
	F_WebPSafeFree(m, v3331)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+8)) = v3324
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+16)) = v3324 + v3323
	*(*int32)(unsafe.Add(mBase, uint32(v3284)+12)) = v3324 + v3296
	v3338 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v3339 = v3338
	v3340 = v3324
	goto L516
L526:
	;
	v3329 = *(*int32)(unsafe.Add(mBase, uint32(v3284)+8))
	v3330 = F_memcpy(m, v3324, v3329, v3296)
	mBase = m.M
	goto L525
L527:
	;
	v3356 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3356 != 0 {
		goto L529
	} else {
		goto L530
	}
L528:
	;
	v3652 = v3168
	goto L59
L529:
	;
	goto L528
L530:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L529
L531:
	;
	v3419 = v3383 + int32(2120)
	v3423 = base.I32_div_s(v3403+int32(-1), int32(16))
	v3425 = v3423 + int32(1)
	v3428 = v3383 + int32(2132)
	v3429 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3428))) = v3429
	*(*int64)(unsafe.Add(mBase, uint32(v3419)+4)) = v3429
	*(*int32)(unsafe.Add(mBase, uint32(v3383+int32(2140)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3428))) = v3383 + int32(2128)
	v3440 = int32(256)
	if v3440 < v3425 {
		goto L552
	} else {
		goto L553
	}
L532:
	;
	v3415 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3415 != 0 {
		goto L549
	} else {
		goto L550
	}
L533:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3383)+12)) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3383))) = l0
	*(*int32)(unsafe.Add(mBase, uint32(v3383)+4)) = v56 + int32(8)
	F_VP8LEncDspInit(m)
	mBase = m.M
	v3399 = v3383 + int32(2216)
	v3400 = *(*int32)(unsafe.Add(mBase, uint32(v3383)+4))
	v3401 = *(*int32)(unsafe.Add(mBase, uint32(v3400)+12))
	v3402 = *(*int32)(unsafe.Add(mBase, uint32(v3400)+8))
	v3403 = v3401 * v3402
	v3407 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v3403), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3399))) = v3407
	if v3407 != 0 {
		goto L545
	} else {
		goto L546
	}
L534:
	;
	if v3383 != 0 {
		goto L533
	} else {
		goto L540
	}
L535:
	;
	goto L534
L536:
	;
	v3381 = F_calloc(m, base.I32_wrap_i64(v3362), v3363)
	mBase = m.M
	v3383 = v3381
	goto L535
L537:
	;
	v3370 = base.I64_div_u_s(int64(2147418112), v3362)
	v3371 = int32(0)
	v3372 = base.I64_extend_i32_u(v3363)
	if base.Ui64(int64(4294967295)) < base.Ui64(v3372*v3362) {
		v3383 = v3371
		goto L535
	} else {
		goto L538
	}
L538:
	;
	if base.Ui64(v3370) < base.Ui64(v3372) {
		v3383 = v3371
		goto L535
	} else {
		goto L539
	}
L539:
	;
	goto L536
L540:
	;
	v3386 = v56 + int32(8)
	v3388 = *(*int32)(unsafe.Add(mBase, uint32(v3386)+92))
	if v3388 != 0 {
		goto L542
	} else {
		goto L543
	}
L541:
	;
	goto L532
L542:
	;
	goto L541
L543:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3386)+92)) = int32(1)
	goto L542
L544:
	;
	if v3407 != 0 {
		goto L531
	} else {
		goto L547
	}
L545:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3399)+4)) = v3403
	goto L544
L546:
	;
	goto L544
L547:
	;
	goto L532
L548:
	;
	v3652 = v3383
	goto L59
L549:
	;
	goto L548
L550:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L549
L551:
	;
	v3446 = v3383 + int32(2144)
	v3449 = v3383 + int32(2156)
	v3450 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3449))) = v3450
	*(*int64)(unsafe.Add(mBase, uint32(v3446)+4)) = v3450
	*(*int32)(unsafe.Add(mBase, uint32(v3383+int32(2164)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3449))) = v3383 + int32(2152)
	v3461 = int32(256)
	if v3461 < v3425 {
		goto L556
	} else {
		goto L557
	}
L552:
	;
	v3443 = v3425
	goto L554
L553:
	;
	v3443 = v3440
	goto L554
L554:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3419))) = v3443
	goto L551
L555:
	;
	v3467 = v3383 + int32(2168)
	v3470 = v3383 + int32(2180)
	v3471 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3470))) = v3471
	*(*int64)(unsafe.Add(mBase, uint32(v3467)+4)) = v3471
	*(*int32)(unsafe.Add(mBase, uint32(v3383+int32(2188)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3470))) = v3383 + int32(2176)
	v3482 = int32(256)
	if v3482 < v3425 {
		goto L560
	} else {
		goto L561
	}
L556:
	;
	v3464 = v3425
	goto L558
L557:
	;
	v3464 = v3461
	goto L558
L558:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3446))) = v3464
	goto L555
L559:
	;
	v3488 = v3383 + int32(2192)
	v3491 = v3383 + int32(2204)
	v3492 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3491))) = v3492
	*(*int64)(unsafe.Add(mBase, uint32(v3488)+4)) = v3492
	*(*int32)(unsafe.Add(mBase, uint32(v3383+int32(2212)))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3491))) = v3383 + int32(2200)
	v3503 = int32(256)
	if v3503 < v3425 {
		goto L564
	} else {
		goto L565
	}
L560:
	;
	v3485 = v3425
	goto L562
L561:
	;
	v3485 = v3482
	goto L562
L562:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3467))) = v3485
	goto L559
L563:
	;
	v3508 = *(*int64)(unsafe.Add(mBase, uint32(v79)+36))
	*(*int64)(unsafe.Add(mBase, uint32(v3383)+36)) = v3508
	v3510 = *(*int32)(unsafe.Add(mBase, uint32(v79)+44))
	*(*int32)(unsafe.Add(mBase, uint32(v3383)+44)) = v3510
	v3512 = *(*int32)(unsafe.Add(mBase, uint32(v79)+68))
	*(*int32)(unsafe.Add(mBase, uint32(v3383)+68)) = v3512
	v3514 = int32(72)
	v3518 = int32(1024)
	v3519 = F_memcpy(m, v3383+v3514, v79+v3514, v3518)
	mBase = m.M
	v3523 = F_memcpy(m, v3383+int32(1096), v292, v3518)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v56)+404)) = v3383
	v3527 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[4])))
	m.T0[v3527].(func(*base.Module, int32))(m, v56+int32(1232))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1248)) = int32(0)
	v3531 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1240)) = v3531 + int32(168)
	*(*int32)(unsafe.Add(mBase, uint32(v56)+1244)) = v56 + int32(392)
	v3538 = v3383
	goto L486
L564:
	;
	v3506 = v3425
	goto L566
L565:
	;
	v3506 = v3503
	goto L566
L566:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3488))) = v3506
	goto L563
L567:
	;
	v3563 = v56 + int32(1256)
	v3564 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[5])))
	m.T0[v3564].(func(*base.Module, int32))(m, v3563)
	mBase = m.M
	v3568 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[6])))
	v3569 = m.T0[v3568].(func(*base.Module, int32) int32)(m, v3563)
	mBase = m.M
	v3572 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[7])))
	m.T0[v3572].(func(*base.Module, int32))(m, v3563)
	mBase = m.M
	if v3128 == int32(0) {
		v3652 = v3538
		goto L59
	} else {
		goto L576
	}
L568:
	;
	v3544 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[8])))
	v3545 = m.T0[v3544].(func(*base.Module, int32) int32)(m, v56+int32(1232))
	mBase = m.M
	if v3545 != 0 {
		goto L569
	} else {
		goto L570
	}
L569:
	;
	v3550 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v3550 == int32(0) {
		goto L574
	} else {
		goto L575
	}
L570:
	;
	v3547 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3547 != 0 {
		goto L572
	} else {
		goto L573
	}
L571:
	;
	v3652 = v3538
	goto L59
L572:
	;
	goto L571
L573:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(1)
	goto L572
L574:
	;
	v3559 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[9])))
	m.T0[v3559].(func(*base.Module, int32))(m, v56+int32(1232))
	mBase = m.M
	goto L567
L575:
	;
	v3556 = F_memcpy(m, v56+int32(204), v3550, int32(188))
	mBase = m.M
	goto L574
L576:
	;
	v3577 = v56 + int32(1232)
	v3578 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[6])))
	v3579 = m.T0[v3578].(func(*base.Module, int32) int32)(m, v3577)
	mBase = m.M
	v3582 = *(*int32)(unsafe.Add(mBase, uint32(v94)+uint32(_c_F_VP8LEncodeStream[7])))
	m.T0[v3582].(func(*base.Module, int32))(m, v3577)
	mBase = m.M
	if v3569 == int32(0) {
		goto L578
	} else {
		goto L579
	}
L577:
	;
	v3591 = *(*int32)(unsafe.Add(mBase, uint32(v56)+192))
	v3592 = *(*int32)(unsafe.Add(mBase, uint32(v56)+188))
	v3594 = *(*int32)(unsafe.Add(mBase, uint32(v56)+184))
	v3595 = int32(7)
	v3597 = int32(3)
	v3600 = *(*int32)(unsafe.Add(mBase, uint32(l2)+12))
	v3601 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	v3603 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	if base.Ui32(v3600-v3601+(v3603+v3595)>>(uint(v3597)%32)) <= base.Ui32(v3591-v3592+(v3594+v3595)>>(uint(v3597)%32)) {
		v3652 = v3538
		goto L59
	} else {
		goto L585
	}
L578:
	;
	v3586 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3586 != 0 {
		v3652 = v3538
		goto L59
	} else {
		goto L581
	}
L579:
	;
	if v3579 != 0 {
		goto L577
	} else {
		goto L580
	}
L580:
	;
	goto L578
L581:
	;
	v3587 = *(*int32)(unsafe.Add(mBase, uint32(v56)+100))
	v3588 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v3588 != 0 {
		goto L583
	} else {
		goto L584
	}
L582:
	;
	v3652 = v3538
	goto L59
L583:
	;
	goto L582
L584:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = v3587
	goto L583
L585:
	;
	v3611 = v56 + int32(180)
	v3620 = *(*int64)(unsafe.Add(mBase, uint32(l2)))
	v3621 = *(*int64)(unsafe.Add(mBase, uint32(v3611)))
	*(*int64)(unsafe.Add(mBase, uint32(l2))) = v3621
	v3624 = l2 + int32(16)
	v3625 = *(*int64)(unsafe.Add(mBase, uint32(v3624)))
	v3627 = v56 + int32(196)
	v3628 = *(*int64)(unsafe.Add(mBase, uint32(v3627)))
	*(*int64)(unsafe.Add(mBase, uint32(v3624))) = v3628
	v3631 = l2 + int32(8)
	v3632 = *(*int64)(unsafe.Add(mBase, uint32(v3631)))
	v3634 = v56 + int32(188)
	v3635 = *(*int64)(unsafe.Add(mBase, uint32(v3634)))
	*(*int64)(unsafe.Add(mBase, uint32(v3631))) = v3635
	*(*int64)(unsafe.Add(mBase, uint32(v3611))) = v3620
	*(*int64)(unsafe.Add(mBase, uint32(v3627))) = v3625
	*(*int64)(unsafe.Add(mBase, uint32(v3634))) = v3632
	goto L586
L586:
	;
	v3640 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v3640 == int32(0) {
		v3652 = v3538
		goto L59
	} else {
		goto L587
	}
L587:
	;
	v3646 = F_memcpy(m, v3640, v56+int32(204), int32(188))
	mBase = m.M
	v3652 = v3538
	goto L59
L588:
	;
	v3719 = v79 + int32(2216)
	v3720 = *(*int32)(unsafe.Add(mBase, uint32(v3719)))
	F_WebPSafeFree(m, v3720)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3719))) = int64(0)
	goto L591
L589:
	;
	goto L588
L590:
	;
	v3706 = v56 + int32(188)
	v3707 = *(*int32)(unsafe.Add(mBase, uint32(v3706)))
	F_WebPSafeFree(m, v3707)
	mBase = m.M
	v3711 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v56+int32(196)))) = v3711
	*(*int64)(unsafe.Add(mBase, uint32(v3706))) = v3711
	*(*int64)(unsafe.Add(mBase, uint32(v3701))) = v3711
	goto L589
L591:
	;
	v3725 = v79 + int32(2120)
	v3728 = *(*int32)(unsafe.Add(mBase, uint32(v3725)+12))
	if v3728 == int32(0) {
		goto L593
	} else {
		goto L594
	}
L592:
	;
	v3754 = v79 + int32(2144)
	v3757 = *(*int32)(unsafe.Add(mBase, uint32(v3754)+12))
	if v3757 == int32(0) {
		goto L601
	} else {
		goto L602
	}
L593:
	;
	v3733 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3725)+20)) = v3733
	v3735 = *(*int32)(unsafe.Add(mBase, uint32(v3725)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3725)+16)) = v3735
	*(*int32)(unsafe.Add(mBase, uint32(v3725)+12)) = v79 + int32(2128)
	*(*int32)(unsafe.Add(mBase, uint32(v3725)+8)) = v3733
	if v3735 == v3733 {
		goto L595
	} else {
		goto L596
	}
L594:
	;
	v3731 = *(*int32)(unsafe.Add(mBase, uint32(v3725)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3728))) = v3731
	goto L593
L595:
	;
	goto L592
L596:
	;
	v3746 = v3735
	goto L597
L597:
	;
	v3747 = *(*int32)(unsafe.Add(mBase, uint32(v3746)))
	F_WebPSafeFree(m, v3746)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3725)+16)) = v3747
	if v3747 != 0 {
		v3746 = v3747
		goto L597
	} else {
		goto L599
	}
L598:
	;
	goto L595
L599:
	;
	goto L598
L600:
	;
	v3783 = v79 + int32(2168)
	v3786 = *(*int32)(unsafe.Add(mBase, uint32(v3783)+12))
	if v3786 == int32(0) {
		goto L609
	} else {
		goto L610
	}
L601:
	;
	v3762 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3754)+20)) = v3762
	v3764 = *(*int32)(unsafe.Add(mBase, uint32(v3754)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3754)+16)) = v3764
	*(*int32)(unsafe.Add(mBase, uint32(v3754)+12)) = v79 + int32(2152)
	*(*int32)(unsafe.Add(mBase, uint32(v3754)+8)) = v3762
	if v3764 == v3762 {
		goto L603
	} else {
		goto L604
	}
L602:
	;
	v3760 = *(*int32)(unsafe.Add(mBase, uint32(v3754)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3757))) = v3760
	goto L601
L603:
	;
	goto L600
L604:
	;
	v3775 = v3764
	goto L605
L605:
	;
	v3776 = *(*int32)(unsafe.Add(mBase, uint32(v3775)))
	F_WebPSafeFree(m, v3775)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3754)+16)) = v3776
	if v3776 != 0 {
		v3775 = v3776
		goto L605
	} else {
		goto L607
	}
L606:
	;
	goto L603
L607:
	;
	goto L606
L608:
	;
	v3812 = v79 + int32(2192)
	v3815 = *(*int32)(unsafe.Add(mBase, uint32(v3812)+12))
	if v3815 == int32(0) {
		goto L617
	} else {
		goto L618
	}
L609:
	;
	v3791 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3783)+20)) = v3791
	v3793 = *(*int32)(unsafe.Add(mBase, uint32(v3783)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3783)+16)) = v3793
	*(*int32)(unsafe.Add(mBase, uint32(v3783)+12)) = v79 + int32(2176)
	*(*int32)(unsafe.Add(mBase, uint32(v3783)+8)) = v3791
	if v3793 == v3791 {
		goto L611
	} else {
		goto L612
	}
L610:
	;
	v3789 = *(*int32)(unsafe.Add(mBase, uint32(v3783)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3786))) = v3789
	goto L609
L611:
	;
	goto L608
L612:
	;
	v3804 = v3793
	goto L613
L613:
	;
	v3805 = *(*int32)(unsafe.Add(mBase, uint32(v3804)))
	F_WebPSafeFree(m, v3804)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3783)+16)) = v3805
	if v3805 != 0 {
		v3804 = v3805
		goto L613
	} else {
		goto L615
	}
L614:
	;
	goto L611
L615:
	;
	goto L614
L616:
	;
	v3840 = *(*int32)(unsafe.Add(mBase, uint32(v79)+24))
	F_free(m, v3840)
	mBase = m.M
	goto L624
L617:
	;
	v3820 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3812)+20)) = v3820
	v3822 = *(*int32)(unsafe.Add(mBase, uint32(v3812)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3812)+16)) = v3822
	*(*int32)(unsafe.Add(mBase, uint32(v3812)+12)) = v79 + int32(2200)
	*(*int32)(unsafe.Add(mBase, uint32(v3812)+8)) = v3820
	if v3822 == v3820 {
		goto L619
	} else {
		goto L620
	}
L618:
	;
	v3818 = *(*int32)(unsafe.Add(mBase, uint32(v3812)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3815))) = v3818
	goto L617
L619:
	;
	goto L616
L620:
	;
	v3833 = v3822
	goto L621
L621:
	;
	v3834 = *(*int32)(unsafe.Add(mBase, uint32(v3833)))
	F_WebPSafeFree(m, v3833)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3812)+16)) = v3834
	if v3834 != 0 {
		v3833 = v3834
		goto L621
	} else {
		goto L623
	}
L622:
	;
	goto L619
L623:
	;
	goto L622
L624:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v79)+24)) = int64(0)
	F_free(m, v79)
	mBase = m.M
	goto L625
L625:
	;
	if v3652 == int32(0) {
		goto L626
	} else {
		goto L627
	}
L626:
	;
	v3974 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	v3982 = base.B2i32(v3974 == int32(0))
	goto L1
L627:
	;
	v3848 = v3652 + int32(2216)
	v3849 = *(*int32)(unsafe.Add(mBase, uint32(v3848)))
	F_WebPSafeFree(m, v3849)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3848))) = int64(0)
	goto L628
L628:
	;
	v3854 = v3652 + int32(2120)
	v3857 = *(*int32)(unsafe.Add(mBase, uint32(v3854)+12))
	if v3857 == int32(0) {
		goto L630
	} else {
		goto L631
	}
L629:
	;
	v3883 = v3652 + int32(2144)
	v3886 = *(*int32)(unsafe.Add(mBase, uint32(v3883)+12))
	if v3886 == int32(0) {
		goto L638
	} else {
		goto L639
	}
L630:
	;
	v3862 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3854)+20)) = v3862
	v3864 = *(*int32)(unsafe.Add(mBase, uint32(v3854)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3854)+16)) = v3864
	*(*int32)(unsafe.Add(mBase, uint32(v3854)+12)) = v3652 + int32(2128)
	*(*int32)(unsafe.Add(mBase, uint32(v3854)+8)) = v3862
	if v3864 == v3862 {
		goto L632
	} else {
		goto L633
	}
L631:
	;
	v3860 = *(*int32)(unsafe.Add(mBase, uint32(v3854)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3857))) = v3860
	goto L630
L632:
	;
	goto L629
L633:
	;
	v3875 = v3864
	goto L634
L634:
	;
	v3876 = *(*int32)(unsafe.Add(mBase, uint32(v3875)))
	F_WebPSafeFree(m, v3875)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3854)+16)) = v3876
	if v3876 != 0 {
		v3875 = v3876
		goto L634
	} else {
		goto L636
	}
L635:
	;
	goto L632
L636:
	;
	goto L635
L637:
	;
	v3912 = v3652 + int32(2168)
	v3915 = *(*int32)(unsafe.Add(mBase, uint32(v3912)+12))
	if v3915 == int32(0) {
		goto L646
	} else {
		goto L647
	}
L638:
	;
	v3891 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3883)+20)) = v3891
	v3893 = *(*int32)(unsafe.Add(mBase, uint32(v3883)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3883)+16)) = v3893
	*(*int32)(unsafe.Add(mBase, uint32(v3883)+12)) = v3652 + int32(2152)
	*(*int32)(unsafe.Add(mBase, uint32(v3883)+8)) = v3891
	if v3893 == v3891 {
		goto L640
	} else {
		goto L641
	}
L639:
	;
	v3889 = *(*int32)(unsafe.Add(mBase, uint32(v3883)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3886))) = v3889
	goto L638
L640:
	;
	goto L637
L641:
	;
	v3904 = v3893
	goto L642
L642:
	;
	v3905 = *(*int32)(unsafe.Add(mBase, uint32(v3904)))
	F_WebPSafeFree(m, v3904)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3883)+16)) = v3905
	if v3905 != 0 {
		v3904 = v3905
		goto L642
	} else {
		goto L644
	}
L643:
	;
	goto L640
L644:
	;
	goto L643
L645:
	;
	v3941 = v3652 + int32(2192)
	v3944 = *(*int32)(unsafe.Add(mBase, uint32(v3941)+12))
	if v3944 == int32(0) {
		goto L654
	} else {
		goto L655
	}
L646:
	;
	v3920 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3912)+20)) = v3920
	v3922 = *(*int32)(unsafe.Add(mBase, uint32(v3912)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3912)+16)) = v3922
	*(*int32)(unsafe.Add(mBase, uint32(v3912)+12)) = v3652 + int32(2176)
	*(*int32)(unsafe.Add(mBase, uint32(v3912)+8)) = v3920
	if v3922 == v3920 {
		goto L648
	} else {
		goto L649
	}
L647:
	;
	v3918 = *(*int32)(unsafe.Add(mBase, uint32(v3912)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3915))) = v3918
	goto L646
L648:
	;
	goto L645
L649:
	;
	v3933 = v3922
	goto L650
L650:
	;
	v3934 = *(*int32)(unsafe.Add(mBase, uint32(v3933)))
	F_WebPSafeFree(m, v3933)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3912)+16)) = v3934
	if v3934 != 0 {
		v3933 = v3934
		goto L650
	} else {
		goto L652
	}
L651:
	;
	goto L648
L652:
	;
	goto L651
L653:
	;
	v3969 = *(*int32)(unsafe.Add(mBase, uint32(v3652)+24))
	F_free(m, v3969)
	mBase = m.M
	goto L661
L654:
	;
	v3949 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3941)+20)) = v3949
	v3951 = *(*int32)(unsafe.Add(mBase, uint32(v3941)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v3941)+16)) = v3951
	*(*int32)(unsafe.Add(mBase, uint32(v3941)+12)) = v3652 + int32(2200)
	*(*int32)(unsafe.Add(mBase, uint32(v3941)+8)) = v3949
	if v3951 == v3949 {
		goto L656
	} else {
		goto L657
	}
L655:
	;
	v3947 = *(*int32)(unsafe.Add(mBase, uint32(v3941)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v3944))) = v3947
	goto L654
L656:
	;
	goto L653
L657:
	;
	v3962 = v3951
	goto L658
L658:
	;
	v3963 = *(*int32)(unsafe.Add(mBase, uint32(v3962)))
	F_WebPSafeFree(m, v3962)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v3941)+16)) = v3963
	if v3963 != 0 {
		v3962 = v3963
		goto L658
	} else {
		goto L660
	}
L659:
	;
	goto L656
L660:
	;
	goto L659
L661:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v3652)+24)) = int64(0)
	F_free(m, v3652)
	mBase = m.M
	goto L662
L662:
	;
	goto L626
}
func F_VP8LFreeHistogram(m *base.Module, l0 int32) {
	F_free(m, l0)
	return
}
func F_VP8LFreeHistogramSet(m *base.Module, l0 int32) {
	F_free(m, l0)
	return
}
func F_VP8LHashChainClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	F_free(m, v2)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(0)
	return
}
func F_VP8LHashChainFill(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) int32 {
	mBase := m.M
	_ = mBase
	var v33 int32
	_ = v33
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v56 int32
	_ = v56
	var v62 int32
	_ = v62
	var v63 int64
	_ = v63
	var v64 int32
	_ = v64
	var v71 int64
	_ = v71
	var v72 int32
	_ = v72
	var v73 int64
	_ = v73
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
	var v94 int32
	_ = v94
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v100 int32
	_ = v100
	var v226 int32
	_ = v226
	var v228 int32
	_ = v228
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v235 int32
	_ = v235
	var v255 int32
	_ = v255
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v269 int32
	_ = v269
	var v272 int32
	_ = v272
	var v274 int32
	_ = v274
	var v278 int32
	_ = v278
	var v282 int32
	_ = v282
	var v295 int32
	_ = v295
	var v303 int32
	_ = v303
	var v323 int32
	_ = v323
	var v328 int32
	_ = v328
	var v334 int32
	_ = v334
	var v337 int32
	_ = v337
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v342 int32
	_ = v342
	var v353 int32
	_ = v353
	var v378 int32
	_ = v378
	var v379 int32
	_ = v379
	var v383 int32
	_ = v383
	var v387 int32
	_ = v387
	var v388 int32
	_ = v388
	var v423 int32
	_ = v423
	var v424 int32
	_ = v424
	var v429 int64
	_ = v429
	var v432 int32
	_ = v432
	var v433 int32
	_ = v433
	var v445 int32
	_ = v445
	var v469 int32
	_ = v469
	var v477 int32
	_ = v477
	var v478 int32
	_ = v478
	var v483 int32
	_ = v483
	var v485 int32
	_ = v485
	var v488 int32
	_ = v488
	var v491 int32
	_ = v491
	var v500 int32
	_ = v500
	var v503 int32
	_ = v503
	var v504 int32
	_ = v504
	var v507 int32
	_ = v507
	var v509 int32
	_ = v509
	var v510 int32
	_ = v510
	var v513 int32
	_ = v513
	var v522 int32
	_ = v522
	var v531 int32
	_ = v531
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v565 int32
	_ = v565
	var v566 int32
	_ = v566
	var v567 int32
	_ = v567
	var v570 int32
	_ = v570
	var v572 int32
	_ = v572
	var v573 int32
	_ = v573
	var v575 int32
	_ = v575
	var v579 int32
	_ = v579
	var v583 int32
	_ = v583
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v592 int32
	_ = v592
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v597 int32
	_ = v597
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v601 int32
	_ = v601
	var v603 int32
	_ = v603
	var v605 int32
	_ = v605
	var v607 int32
	_ = v607
	var v609 int32
	_ = v609
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v613 int32
	_ = v613
	var v614 int32
	_ = v614
	var v617 int32
	_ = v617
	var v618 int32
	_ = v618
	var v621 int32
	_ = v621
	var v623 int32
	_ = v623
	var v626 int32
	_ = v626
	var v627 int32
	_ = v627
	var v628 int32
	_ = v628
	var v631 int32
	_ = v631
	var v634 int32
	_ = v634
	var v637 int32
	_ = v637
	var v641 int32
	_ = v641
	var v643 int32
	_ = v643
	var v646 int32
	_ = v646
	var v654 int32
	_ = v654
	var v663 int32
	_ = v663
	var v664 int32
	_ = v664
	var v675 int32
	_ = v675
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v680 int32
	_ = v680
	var v684 int32
	_ = v684
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v688 int32
	_ = v688
	var v690 int32
	_ = v690
	var v695 int32
	_ = v695
	var v696 int32
	_ = v696
	var v697 int32
	_ = v697
	var v698 int32
	_ = v698
	var v701 int32
	_ = v701
	var v704 int32
	_ = v704
	var v725 int32
	_ = v725
	var v735 int32
	_ = v735
	var v736 int32
	_ = v736
	var v738 int32
	_ = v738
	var v742 int32
	_ = v742
	var v758 int32
	_ = v758
	var v761 int32
	_ = v761
	var v769 int32
	_ = v769
	var v773 int32
	_ = v773
	var v780 int32
	_ = v780
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v799 int32
	_ = v799
	var v800 int32
	_ = v800
	var v803 int32
	_ = v803
	var v804 int32
	_ = v804
	var v809 int32
	_ = v809
	var v815 int32
	_ = v815
	var v818 int32
	_ = v818
	var v826 int32
	_ = v826
	var v830 int32
	_ = v830
	var v846 int32
	_ = v846
	var v849 int32
	_ = v849
	var v853 int32
	_ = v853
	var v857 int32
	_ = v857
	var v862 int32
	_ = v862
	var v868 int32
	_ = v868
	var v898 int32
	_ = v898
	var v900 int32
	_ = v900
	var v904 int32
	_ = v904
	var v907 int32
	_ = v907
	var v931 int32
	_ = v931
	var v939 int32
	_ = v939
	var v952 int32
	_ = v952
	var v953 int32
	_ = v953
	var v960 int32
	_ = v960
	var v962 int32
	_ = v962
	var v963 int32
	_ = v963
	var v973 int32
	_ = v973
	var v976 int32
	_ = v976
	var v984 int32
	_ = v984
	var v995 int32
	_ = v995
	var v1004 int32
	_ = v1004
	var v1006 int32
	_ = v1006
	var v1008 int32
	_ = v1008
	var v1009 int32
	_ = v1009
	var v1020 int32
	_ = v1020
	var v1021 int32
	_ = v1021
	var v1027 int32
	_ = v1027
	var v1033 int32
	_ = v1033
	var v1035 int32
	_ = v1035
	var v1055 int32
	_ = v1055
	var v1067 int32
	_ = v1067
	var v1069 int32
	_ = v1069
	v33 = l4 * l3
	if int32(75) < l1 {
		v47 = int32(1048456)
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v48 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if int32(2) < v33 {
		goto L8
	} else {
		goto L9
	}
L2:
	;
	if l1 < int32(51) {
		goto L3
	} else {
		goto L4
	}
L3:
	;
	if int32(25) < l1 {
		goto L5
	} else {
		goto L6
	}
L4:
	;
	v47 = l3 << (uint(int32(8)) % 32)
	goto L1
L5:
	;
	v45 = int32(6)
	goto L7
L6:
	;
	v45 = int32(4)
	goto L7
L7:
	;
	v47 = l3 << (uint(v45) % 32)
	goto L1
L8:
	;
	v62 = *(*int32)(unsafe.Add(mBase, uint32(l8)))
	v63 = int64(262144)
	v64 = int32(4)
	goto L14
L9:
	;
	v56 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v48+v33<<(uint(int32(2))%32)+int32(-4)))) = v56
	*(*int32)(unsafe.Add(mBase, uint32(v48))) = v56
	return int32(1)
L10:
	;
	v94 = int32(base.Ui32(l1*l1) >> (uint(int32(7)) % 32))
	v96 = v94 + int32(8)
	v97 = int32(1048456)
	if v47 < v97 {
		goto L21
	} else {
		goto L22
	}
L11:
	;
	if v85 != 0 {
		goto L10
	} else {
		goto L17
	}
L12:
	;
	goto L11
L13:
	;
	v83 = F_malloc(m, base.I32_wrap_i64(v63)*v64)
	mBase = m.M
	v85 = v83
	goto L12
L14:
	;
	v71 = base.I64_div_u_s(int64(2147418112), v63)
	v72 = int32(0)
	v73 = base.I64_extend_i32_u(v64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v73*v63) {
		v85 = v72
		goto L12
	} else {
		goto L15
	}
L15:
	;
	if base.Ui64(v71) < base.Ui64(v73) {
		v85 = v72
		goto L12
	} else {
		goto L16
	}
L16:
	;
	goto L13
L17:
	;
	v88 = *(*int32)(unsafe.Add(mBase, uint32(l6)+92))
	if v88 != 0 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	return int32(0)
L19:
	;
	goto L18
L20:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l6)+92)) = int32(1)
	goto L19
L21:
	;
	v100 = v47
	goto L23
L22:
	;
	v100 = v97
	goto L23
L23:
	;
	goto L26
L24:
	;
	v226 = v33 + int32(-2)
	v228 = base.I32_div_s(l7, int32(2))
	v230 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
	v231 = *(*int32)(unsafe.Add(mBase, uint32(l2)+4))
	v235 = int32(0)
	v255 = base.B2i32(v230 == v231)
	goto L38
L26:
	;
	base.MemoryFill(m, v85, int32(255), int32(1048576))
	goto L24
L38:
	;
	if v226 <= v235 {
		goto L42
	} else {
		goto L43
	}
L39:
	;
	F_free(m, v85)
	mBase = m.M
	goto L150
L40:
	;
	v1067 = base.I32_div_s(v1035*v228, v226)
	v1069 = F_WebPReportProgress(m, l6, v1067+v62, l8)
	mBase = m.M
	if v1069 != 0 {
		v235 = v1035
		v255 = v1055
		goto L38
	} else {
		goto L149
	}
L41:
	;
	v939 = v282 * int32(1540483478)
	if v931&int32(1) != 0 {
		goto L143
	} else {
		goto L144
	}
L42:
	;
	v483 = v235 << (uint(int32(2)) % 32)
	v485 = l2 + v483
	v488 = *(*int32)(unsafe.Add(mBase, uint32(v485+int32(4))))
	v491 = *(*int32)(unsafe.Add(mBase, uint32(v485)))
	v500 = *(*int32)(unsafe.Add(mBase, uint32(v85+int32(base.Ui32(v488*int32(-962287725)+v491*int32(1540483478))>>(uint(int32(12))%32))&int32(1048572))))
	*(*int32)(unsafe.Add(mBase, uint32(v48+v483))) = v500
	F_free(m, v85)
	mBase = m.M
	goto L71
L43:
	;
	v267 = int32(2)
	v268 = v235 << (uint(v267) % 32)
	v269 = l2 + v268
	v272 = *(*int32)(unsafe.Add(mBase, uint32(v269+int32(8))))
	v274 = v235 + int32(1)
	v278 = *(*int32)(unsafe.Add(mBase, uint32(l2+v274<<(uint(v267)%32))))
	if v255 == int32(0) {
		goto L44
	} else {
		goto L45
	}
L44:
	;
	v469 = *(*int32)(unsafe.Add(mBase, uint32(v269)))
	v477 = v85 + int32(base.Ui32(v278*int32(-962287725)+v469*int32(1540483478))>>(uint(int32(12))%32))&int32(1048572)
	v478 = *(*int32)(unsafe.Add(mBase, uint32(v477)))
	*(*int32)(unsafe.Add(mBase, uint32(v48+v268))) = v478
	*(*int32)(unsafe.Add(mBase, uint32(v477))) = v235
	v1035 = v274
	v1055 = base.B2i32(v278 == v272)
	goto L40
L45:
	;
	if v278 != v272 {
		goto L44
	} else {
		goto L46
	}
L46:
	;
	v282 = *(*int32)(unsafe.Add(mBase, uint32(v269)))
	if v33 <= v235+int32(3) {
		v907 = v235
		v931 = int32(1)
		goto L41
	} else {
		goto L47
	}
L47:
	;
	v295 = l2 + int32(12) + v268
	v303 = int32(3)
	goto L50
L48:
	;
	if base.Ui32(v334) < base.Ui32(int32(_a_F_VP8LHashChainFill_0)) {
		goto L54
	} else {
		goto L55
	}
L49:
	;
	v334 = v303 + int32(-2)
	goto L48
L50:
	;
	v323 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	if v323 != v282 {
		goto L49
	} else {
		goto L52
	}
L52:
	;
	v328 = v303 + int32(1)
	if v33-v235 != v328 {
		v295 = v295 + int32(4)
		v303 = v328
		goto L50
	} else {
		goto L53
	}
L53:
	;
	v334 = v226 - v235
	goto L48
L54:
	;
	if v334 != 0 {
		v907 = v235
		v931 = v334
		goto L41
	} else {
		goto L70
	}
L55:
	;
	v337 = v48 + v268
	v338 = int32(255)
	v340 = v334 + int32(-4095)
	v342 = v340 << (uint(int32(2)) % 32)
	if base.Ui32(v342) < base.Ui32(int32(33)) {
		goto L57
	} else {
		goto L58
	}
L56:
	;
	v907 = v340 + v235
	v931 = int32(4095)
	goto L41
L57:
	;
	if v342 == int32(0) {
		goto L59
	} else {
		goto L60
	}
L58:
	;
	base.MemoryFill(m, v337, v338, v342)
	goto L56
L59:
	;
	goto L56
L60:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v337))) = uint8(v338)
	v353 = v337 + v342
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-1)))) = uint8(v338)
	if base.Ui32(v342) < base.Ui32(int32(3)) {
		goto L59
	} else {
		goto L61
	}
L61:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v337)+2)) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v337)+1)) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-3)))) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-2)))) = uint8(v338)
	if base.Ui32(v342) < base.Ui32(int32(7)) {
		goto L59
	} else {
		goto L62
	}
L62:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v337)+3)) = uint8(v338)
	*(*uint8)(unsafe.Add(mBase, uint32(v353+int32(-4)))) = uint8(v338)
	if base.Ui32(v342) < base.Ui32(int32(9)) {
		goto L59
	} else {
		goto L63
	}
L63:
	;
	v378 = (int32(0) - v337) & int32(3)
	v379 = v337 + v378
	v383 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v379))) = v383
	v387 = (v342 - v378) & int32(60)
	v388 = v379 + v387
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-4)))) = v383
	if base.Ui32(v387) < base.Ui32(int32(9)) {
		goto L59
	} else {
		goto L64
	}
L64:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v379)+8)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+4)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-8)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-12)))) = v383
	if base.Ui32(v387) < base.Ui32(int32(25)) {
		goto L59
	} else {
		goto L65
	}
L65:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v379)+24)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+20)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+16)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v379)+12)) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-16)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-20)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-24)))) = v383
	*(*int32)(unsafe.Add(mBase, uint32(v388+int32(-28)))) = v383
	v423 = v379&int32(4) | int32(24)
	v424 = v387 - v423
	if base.Ui32(v424) < base.Ui32(int32(32)) {
		goto L59
	} else {
		goto L66
	}
L66:
	;
	v429 = base.I64_extend_i32_u(v383) * int64(4294967297)
	v432 = v424
	v433 = v379 + v423
	goto L67
L67:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v433)+24)) = v429
	*(*int64)(unsafe.Add(mBase, uint32(v433)+16)) = v429
	*(*int64)(unsafe.Add(mBase, uint32(v433)+8)) = v429
	*(*int64)(unsafe.Add(mBase, uint32(v433))) = v429
	v445 = v432 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v445) {
		v432 = v445
		v433 = v433 + int32(32)
		goto L67
	} else {
		goto L69
	}
L68:
	;
	goto L59
L69:
	;
	goto L68
L70:
	;
	v1035 = v235
	v1055 = int32(0)
	goto L40
L71:
	;
	v503 = v62 + v228
	v504 = F_WebPReportProgress(m, l6, v503, l8)
	mBase = m.M
	if v504 != 0 {
		goto L72
	} else {
		goto L73
	}
L72:
	;
	v507 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v509 = v33 + int32(-1)
	v510 = int32(2)
	v513 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v507+v509<<(uint(v510)%32)))) = v513
	*(*int32)(unsafe.Add(mBase, uint32(v507))) = v513
	v522 = v94 + int32(7)
	v531 = v226
	goto L75
L73:
	;
	return int32(0)
L74:
	;
	v904 = F_WebPReportProgress(m, l6, v62+l7, l8)
	mBase = m.M
	return v904
L75:
	;
	if v531 == int32(0) {
		goto L74
	} else {
		goto L77
	}
L76:
	;
	return int32(0)
L77:
	;
	v561 = int32(0)
	v563 = v531 - v100
	if base.Ui32(v531) < base.Ui32(v563) {
		goto L78
	} else {
		goto L79
	}
L78:
	;
	v565 = v561
	goto L80
L79:
	;
	v565 = v563
	goto L80
L80:
	;
	v566 = v509 - v531
	v567 = int32(4095)
	if v566 < v567 {
		goto L81
	} else {
		goto L82
	}
L81:
	;
	v570 = v566
	goto L83
L82:
	;
	v570 = v567
	goto L83
L83:
	;
	v572 = v531 << (uint(int32(2)) % 32)
	v573 = l2 + v572
	v575 = *(*int32)(unsafe.Add(mBase, uint32(v48+v572)))
	if l5 == int32(0) {
		goto L85
	} else {
		goto L86
	}
L84:
	;
	if v627 < v565 {
		v704 = v626
		v725 = v631
		goto L107
	} else {
		goto L108
	}
L85:
	;
	v579 = int32(0)
	if base.Ui32(v531) < base.Ui32(l3) {
		v598 = v579
		v599 = v579
		v601 = v96
		goto L87
	} else {
		goto L88
	}
L86:
	;
	v626 = v561
	v627 = v575
	v628 = v96
	v631 = int32(0)
	goto L84
L87:
	;
	v603 = v573 + int32(-4)
	v605 = v598 << (uint(int32(2)) % 32)
	v607 = *(*int32)(unsafe.Add(mBase, uint32(v603+v605)))
	v609 = *(*int32)(unsafe.Add(mBase, uint32(v573+v605)))
	if v607 != v609 {
		v614 = v579
		goto L96
	} else {
		goto L97
	}
L88:
	;
	v583 = int32(0)
	v585 = v573 + (v513-l3)<<(uint(v510)%32)
	v586 = *(*int32)(unsafe.Add(mBase, uint32(v585)))
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v573)))
	if v586 != v587 {
		v598 = v583
		v599 = v583
		v601 = v522
		goto L87
	} else {
		goto L89
	}
L89:
	;
	v589 = m.G63
	v590 = *(*int32)(unsafe.Add(mBase, uint32(v589)))
	v591 = m.T0[v590].(func(*base.Module, int32, int32, int32) int32)(m, v585, v573, v570)
	mBase = m.M
	v592 = int32(0)
	v594 = base.B2i32(v592 < v591)
	if v592 < v591 {
		goto L90
	} else {
		goto L91
	}
L90:
	;
	v595 = v591
	goto L92
L91:
	;
	v595 = v592
	goto L92
L92:
	;
	if v592 < v591 {
		goto L93
	} else {
		goto L94
	}
L93:
	;
	v597 = l3
	goto L95
L94:
	;
	v597 = int32(0)
	goto L95
L95:
	;
	v598 = v595
	v599 = v597
	v601 = v522
	goto L87
L96:
	;
	v617 = base.B2i32(v598 < v614)
	if v598 < v614 {
		goto L98
	} else {
		goto L99
	}
L97:
	;
	v611 = m.G63
	v612 = *(*int32)(unsafe.Add(mBase, uint32(v611)))
	v613 = m.T0[v612].(func(*base.Module, int32, int32, int32) int32)(m, v603, v573, v570)
	mBase = m.M
	v614 = v613
	goto L96
L98:
	;
	v618 = v614
	goto L100
L99:
	;
	v618 = v598
	goto L100
L100:
	;
	if v618 == int32(4095) {
		goto L101
	} else {
		goto L102
	}
L101:
	;
	v621 = v565 + int32(-1)
	goto L103
L102:
	;
	v621 = v575
	goto L103
L103:
	;
	if v598 < v614 {
		goto L104
	} else {
		goto L105
	}
L104:
	;
	v623 = int32(1)
	goto L106
L105:
	;
	v623 = v599
	goto L106
L106:
	;
	v626 = v618
	v627 = v621
	v628 = v601 + int32(-1)
	v631 = v623
	goto L84
L107:
	;
	v735 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v736 = v735 + v572
	v738 = v725 << (uint(int32(12)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v736))) = v704 | v738
	v742 = v531 + int32(-1)
	if v725 == int32(0) {
		v868 = v742
		goto L121
	} else {
		goto L122
	}
L108:
	;
	v634 = int32(256)
	if v566 < v634 {
		goto L109
	} else {
		goto L110
	}
L109:
	;
	v637 = v566
	goto L111
L110:
	;
	v637 = v634
	goto L111
L111:
	;
	v641 = *(*int32)(unsafe.Add(mBase, uint32(v573+v626<<(uint(int32(2))%32))))
	v643 = v626
	v646 = v627
	v654 = v628
	v663 = v641
	v664 = v631
	goto L112
L112:
	;
	v675 = v654 + int32(-1)
	if v675 == int32(0) {
		v704 = v643
		v725 = v664
		goto L107
	} else {
		goto L114
	}
L113:
	;
	v704 = v696
	v725 = v698
	goto L107
L114:
	;
	v678 = int32(2)
	v679 = v646 << (uint(v678) % 32)
	v680 = l2 + v679
	v684 = *(*int32)(unsafe.Add(mBase, uint32(v680+v643<<(uint(v678)%32))))
	if v684 != v663 {
		v696 = v643
		v697 = v663
		v698 = v664
		goto L115
	} else {
		goto L116
	}
L115:
	;
	v701 = *(*int32)(unsafe.Add(mBase, uint32(v48+v679)))
	if v565 <= v701 {
		v643 = v696
		v646 = v701
		v654 = v675
		v663 = v697
		v664 = v698
		goto L112
	} else {
		goto L120
	}
L116:
	;
	v686 = m.G63
	v687 = *(*int32)(unsafe.Add(mBase, uint32(v686)))
	v688 = m.T0[v687].(func(*base.Module, int32, int32, int32) int32)(m, v680, v573, v570)
	mBase = m.M
	if v688 <= v643 {
		v696 = v643
		v697 = v663
		v698 = v664
		goto L115
	} else {
		goto L117
	}
L117:
	;
	v690 = v531 - v646
	if v688 < v637 {
		goto L118
	} else {
		goto L119
	}
L118:
	;
	v695 = *(*int32)(unsafe.Add(mBase, uint32(v573+v688<<(uint(int32(2))%32))))
	v696 = v688
	v697 = v695
	v698 = v690
	goto L115
L119:
	;
	v704 = v688
	v725 = v690
	goto L107
L120:
	;
	goto L113
L121:
	;
	v898 = base.I32_div_u_s((v226-v868)*(l7-v228), v226)
	v900 = F_WebPReportProgress(m, l6, v898+v503, l8)
	mBase = m.M
	if v900 != 0 {
		v531 = v868
		goto L75
	} else {
		goto L141
	}
L122:
	;
	if v742 == int32(0) {
		v868 = v742
		goto L121
	} else {
		goto L123
	}
L123:
	;
	if base.Ui32(v742) < base.Ui32(v725) {
		v868 = v742
		goto L121
	} else {
		goto L124
	}
L124:
	;
	if v725 == int32(1) {
		goto L125
	} else {
		goto L126
	}
L125:
	;
	v815 = v704
	v818 = v742
	v826 = l2 + int32(-8) + v572
	v830 = v736 + int32(-4)
	goto L137
L126:
	;
	v758 = v704
	v761 = v742
	v769 = l2 + int32(-4) + v572
	v773 = v531
	v780 = v736 + int32(-4)
	goto L127
L127:
	;
	v790 = *(*int32)(unsafe.Add(mBase, uint32(v769+(int32(0)-v725<<(uint(int32(2))%32)))))
	v791 = *(*int32)(unsafe.Add(mBase, uint32(v769)))
	if v790 != v791 {
		v868 = v761
		goto L121
	} else {
		goto L129
	}
L129:
	;
	if v758 != int32(4095) {
		goto L130
	} else {
		goto L131
	}
L130:
	;
	v799 = base.B2i32(v758 < int32(4095))
	v800 = v758 + v799
	*(*int32)(unsafe.Add(mBase, uint32(v780))) = v800 | v738
	if v758 < int32(4095) {
		goto L133
	} else {
		goto L134
	}
L131:
	;
	if base.Ui32(v761+int32(4095)) < base.Ui32(v773) {
		v868 = v761
		goto L121
	} else {
		goto L132
	}
L132:
	;
	goto L130
L133:
	;
	v803 = v761
	goto L135
L134:
	;
	v803 = v773
	goto L135
L135:
	;
	v804 = int32(-4)
	v809 = v761 + int32(-1)
	if base.Ui32(v725) <= base.Ui32(v809) {
		v758 = v800
		v761 = v809
		v769 = v769 + v804
		v773 = v803
		v780 = v780 + v804
		goto L127
	} else {
		goto L136
	}
L136:
	;
	v868 = v809
	goto L121
L137:
	;
	v846 = *(*int32)(unsafe.Add(mBase, uint32(v826)))
	v849 = *(*int32)(unsafe.Add(mBase, uint32(v826+int32(4))))
	if v846 != v849 {
		v868 = v818
		goto L121
	} else {
		goto L139
	}
L138:
	;
	v868 = int32(0)
	goto L121
L139:
	;
	v853 = v815 + base.B2i32(v815 < int32(4095))
	*(*int32)(unsafe.Add(mBase, uint32(v830))) = v853 | int32(_a_F_VP8LHashChainFill_0)
	v857 = int32(-4)
	v862 = v818 + int32(-1)
	if v862 != 0 {
		v815 = v853
		v818 = v862
		v826 = v826 + v857
		v830 = v830 + v857
		goto L137
	} else {
		goto L140
	}
L140:
	;
	goto L138
L141:
	;
	goto L76
L142:
	;
	v963 = int32(0)
	if v931 == int32(1) {
		v1035 = v960
		v1055 = v963
		goto L40
	} else {
		goto L145
	}
L143:
	;
	v952 = v85 + int32(base.Ui32(v931*int32(-962287725)+v939)>>(uint(int32(12))%32))&int32(1048572)
	v953 = *(*int32)(unsafe.Add(mBase, uint32(v952)))
	*(*int32)(unsafe.Add(mBase, uint32(v48+v907<<(uint(int32(2))%32)))) = v953
	*(*int32)(unsafe.Add(mBase, uint32(v952))) = v907
	v960 = v907 + int32(1)
	v962 = v931 + int32(-1)
	goto L142
L144:
	;
	v960 = v907
	v962 = v931
	goto L142
L145:
	;
	v973 = v960
	v976 = v939 - v962*int32(962287725)
	v984 = v48 + v960<<(uint(int32(2))%32)
	v995 = v962
	goto L146
L146:
	;
	v1004 = int32(12)
	v1006 = int32(1048572)
	v1008 = v85 + int32(base.Ui32(v976)>>(uint(v1004)%32))&v1006
	v1009 = *(*int32)(unsafe.Add(mBase, uint32(v1008)))
	*(*int32)(unsafe.Add(mBase, uint32(v984))) = v1009
	*(*int32)(unsafe.Add(mBase, uint32(v1008))) = v973
	v1020 = v85 + int32(base.Ui32(v976+int32(962287725))>>(uint(v1004)%32))&v1006
	v1021 = *(*int32)(unsafe.Add(mBase, uint32(v1020)))
	*(*int32)(unsafe.Add(mBase, uint32(v984+int32(4)))) = v1021
	*(*int32)(unsafe.Add(mBase, uint32(v1020))) = v973 + int32(1)
	v1027 = v973 + int32(2)
	v1033 = v995 + int32(-2)
	if v1033 != 0 {
		v973 = v1027
		v976 = v976 + int32(1924575450)
		v984 = v984 + int32(8)
		v995 = v1033
		goto L146
	} else {
		goto L148
	}
L147:
	;
	v1035 = v1027
	v1055 = v963
	goto L40
L148:
	;
	goto L147
L149:
	;
	goto L39
L150:
	;
	return int32(0)
}
func F_VP8LHashChainInit(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int64
	_ = v4
	var v5 int32
	_ = v5
	var v12 int64
	_ = v12
	var v13 int32
	_ = v13
	var v14 int64
	_ = v14
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	v4 = base.I64_extend_i32_s(l1)
	v5 = int32(4)
	if v4 == int64(0) {
		v24 = F_malloc(m, base.I32_wrap_i64(v4)*v5)
		mBase = m.M
		v26 = v24
	} else {
		v12 = base.I64_div_u_s(int64(2147418112), v4)
		v13 = int32(0)
		v14 = base.I64_extend_i32_u(v5)
		if base.Ui64(int64(4294967295)) < base.Ui64(v14*v4) {
			v26 = v13
		} else {
			if base.Ui64(v12) < base.Ui64(v14) {
				v26 = v13
			} else {
				v24 = F_malloc(m, base.I32_wrap_i64(v4)*v5)
				mBase = m.M
				v26 = v24
			}
		}
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v26
	if v26 != 0 {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = l1
		return int32(1)
	} else {
		return int32(0)
	}
}
func F_VP8LHistogramCreate(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v12 int32
	_ = v12
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
	var v34 int32
	_ = v34
	var v56 int32
	_ = v56
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v110 int64
	_ = v110
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v126 int32
	_ = v126
	var v150 int32
	_ = v150
	var v156 int32
	_ = v156
	var v158 int64
	_ = v158
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v187 int32
	_ = v187
	var v188 int32
	_ = v188
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v201 int32
	_ = v201
	var v204 int32
	_ = v204
	var v207 int32
	_ = v207
	var v209 int32
	_ = v209
	var v211 int32
	_ = v211
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v222 int32
	_ = v222
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v228 int32
	_ = v228
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	v6 = m.G0
	v8 = v6 - int32(16)
	m.G0 = v8
	if int32(-1) < l2 {
		v13 = l2
	} else {
		v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3236))
		v13 = v12
	}
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v15 = int32(0)
	v18 = int32(_a_F_VP8LHistogramCreate_0)
	if v15 < v13 {
		v23 = int32(4)<<(uint(v13)%32) + v18
	} else {
		v23 = v18
	}
	if base.Ui32(v23) < base.Ui32(int32(33)) {
		if v23 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v15)
			v34 = l0 + v23
			*(*uint8)(unsafe.Add(mBase, uint32(v34+int32(-1)))) = uint8(v15)
			if base.Ui32(v23) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v15)
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v15)
				*(*uint8)(unsafe.Add(mBase, uint32(v34+int32(-3)))) = uint8(v15)
				*(*uint8)(unsafe.Add(mBase, uint32(v34+int32(-2)))) = uint8(v15)
				if base.Ui32(v23) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v15)
					*(*uint8)(unsafe.Add(mBase, uint32(v34+int32(-4)))) = uint8(v15)
					if base.Ui32(v23) < base.Ui32(int32(9)) {
					} else {
						v56 = int32(0)
						v59 = (v56 - l0) & int32(3)
						v60 = l0 + v59
						*(*int32)(unsafe.Add(mBase, uint32(v60))) = v56
						v68 = (v23 - v59) & int32(60)
						v69 = v60 + v68
						*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-4)))) = v56
						if base.Ui32(v68) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v60)+8)) = v56
							*(*int32)(unsafe.Add(mBase, uint32(v60)+4)) = v56
							*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-8)))) = v56
							*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-12)))) = v56
							if base.Ui32(v68) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v60)+24)) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v60)+20)) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v60)+16)) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v60)+12)) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-16)))) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-20)))) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-24)))) = v56
								*(*int32)(unsafe.Add(mBase, uint32(v69+int32(-28)))) = v56
								v104 = v60&int32(4) | int32(24)
								v105 = v68 - v104
								if base.Ui32(v105) < base.Ui32(int32(32)) {
								} else {
									v110 = base.I64_extend_i32_u(v56) * int64(4294967297)
									v113 = v105
									v114 = v60 + v104
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v114)+24)) = v110
										*(*int64)(unsafe.Add(mBase, uint32(v114)+16)) = v110
										*(*int64)(unsafe.Add(mBase, uint32(v114)+8)) = v110
										*(*int64)(unsafe.Add(mBase, uint32(v114))) = v110
										v126 = v113 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v126) {
											v113 = v126
											v114 = v114 + int32(32)
											continue
										} else {
											break
										}
										break
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		base.MemoryFill(m, l0, v15, v23)
	}
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v14
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3236)) = v13
	v150 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(3248)))) = uint16(v150)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3304)) = int32(16843009)
	v156 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3308)))) = uint8(v156)
	v158 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3256)) = v158
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3264)))) = v158
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3272)))) = v158
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3280)))) = v158
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3288)))) = v158
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3296)))) = v158
	v181 = v8 + int32(4)
	v183 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v181)+4)) = v183
	if v183 != 0 {
		v187 = *(*int32)(unsafe.Add(mBase, uint32(v183)+4))
		v188 = *(*int32)(unsafe.Add(mBase, uint32(v183)+8))
		v192 = v187 + v188<<(uint(int32(3))%32)
		v193 = v187
	} else {
		v185 = int32(0)
		v192 = v185
		v193 = v185
	}
	*(*int32)(unsafe.Add(mBase, uint32(v181)+8)) = v192
	*(*int32)(unsafe.Add(mBase, uint32(v181))) = v193
	v196 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
	if v196 == int32(0) {
	} else {
		v201 = v196
		for {
			v204 = int32(0)
			F_HistogramAddSinglePixOrCopy(m, l0, v201, v204, v204)
			mBase = m.M
			v207 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
			v209 = v207 + int32(8)
			*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v209
			v211 = *(*int32)(unsafe.Add(mBase, uint32(v8)+12))
			if v209 != v211 {
				v233 = v209
			} else {
				v214 = v8 + int32(4)
				v218 = *(*int32)(unsafe.Add(mBase, uint32(v214)+4))
				v219 = *(*int32)(unsafe.Add(mBase, uint32(v218)))
				if v219 != 0 {
					v222 = *(*int32)(unsafe.Add(mBase, uint32(v219)+4))
					v223 = *(*int32)(unsafe.Add(mBase, uint32(v219)+8))
					v227 = v222
					v228 = v222 + v223<<(uint(int32(3))%32)
				} else {
					v220 = int32(0)
					v227 = v220
					v228 = v220
				}
				*(*int32)(unsafe.Add(mBase, uint32(v214)+8)) = v228
				*(*int32)(unsafe.Add(mBase, uint32(v214))) = v227
				*(*int32)(unsafe.Add(mBase, uint32(v214)+4)) = v219
				v232 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				v233 = v232
			}
			if v233 != 0 {
				v201 = v233
				continue
			} else {
				break
			}
			break
		}
	}
	m.G0 = v8 + int32(16)
	return
}
func F_VP8LHistogramEstimateBits(m *base.Module, l0 int32) int64 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v10 int32
	_ = v10
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v20 int64
	_ = v20
	var v23 int32
	_ = v23
	var v26 int64
	_ = v26
	var v32 int64
	_ = v32
	var v38 int64
	_ = v38
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int64
	_ = v44
	var v45 int32
	_ = v45
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3236))
	v12 = int32(280)
	if int32(0) < v10 {
		v17 = int32(1)<<(uint(v10)%32) + v12
	} else {
		v17 = v12
	}
	v18 = int32(0)
	v20 = F_PopulationCost(m, v8, v17, v18, v18)
	mBase = m.M
	v23 = int32(256)
	v26 = F_PopulationCost(m, l0+int32(4), v23, v18, v18)
	mBase = m.M
	v32 = F_PopulationCost(m, l0+int32(1028), v23, v18, v18)
	mBase = m.M
	v38 = F_PopulationCost(m, l0+int32(2052), v23, v18, v18)
	mBase = m.M
	v40 = l0 + int32(3076)
	v41 = int32(40)
	v44 = F_PopulationCost(m, v40, v41, v18, v18)
	mBase = m.M
	v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v49 = m.G55
	v50 = *(*int32)(unsafe.Add(mBase, uint32(v49)))
	v51 = m.T0[v50].(func(*base.Module, int32, int32) int32)(m, v45+int32(1024), int32(24))
	mBase = m.M
	v53 = *(*int32)(unsafe.Add(mBase, uint32(v49)))
	v54 = m.T0[v53].(func(*base.Module, int32, int32) int32)(m, v40, v41)
	mBase = m.M
	return base.I64_extend_i32_u(v51+v54)<<(uint(int64(23))%64) + (v44 + (v38 + (v32 + (v26 + v20))))
}
func F_VP8LHistogramInit(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v28 int32
	_ = v28
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v104 int64
	_ = v104
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v120 int32
	_ = v120
	var v146 int64
	_ = v146
	var v150 int32
	_ = v150
	var v154 int32
	_ = v154
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3236)) = l1
	if l2 == int32(0) {
	} else {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		v9 = int32(0)
		v12 = int32(_a_F_VP8LHistogramInit_0)
		if v9 < l1 {
			v17 = int32(4)<<(uint(l1)%32) + v12
		} else {
			v17 = v12
		}
		if base.Ui32(v17) < base.Ui32(int32(33)) {
			if v17 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v9)
				v28 = l0 + v17
				*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(-1)))) = uint8(v9)
				if base.Ui32(v17) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v9)
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v9)
					*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(-3)))) = uint8(v9)
					*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(-2)))) = uint8(v9)
					if base.Ui32(v17) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v9)
						*(*uint8)(unsafe.Add(mBase, uint32(v28+int32(-4)))) = uint8(v9)
						if base.Ui32(v17) < base.Ui32(int32(9)) {
						} else {
							v50 = int32(0)
							v53 = (v50 - l0) & int32(3)
							v54 = l0 + v53
							*(*int32)(unsafe.Add(mBase, uint32(v54))) = v50
							v62 = (v17 - v53) & int32(60)
							v63 = v54 + v62
							*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-4)))) = v50
							if base.Ui32(v62) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v54)+8)) = v50
								*(*int32)(unsafe.Add(mBase, uint32(v54)+4)) = v50
								*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-8)))) = v50
								*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-12)))) = v50
								if base.Ui32(v62) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v54)+24)) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v54)+20)) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v54)+16)) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v54)+12)) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-16)))) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-20)))) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-24)))) = v50
									*(*int32)(unsafe.Add(mBase, uint32(v63+int32(-28)))) = v50
									v98 = v54&int32(4) | int32(24)
									v99 = v62 - v98
									if base.Ui32(v99) < base.Ui32(int32(32)) {
									} else {
										v104 = base.I64_extend_i32_u(v50) * int64(4294967297)
										v107 = v99
										v108 = v54 + v98
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v108)+24)) = v104
											*(*int64)(unsafe.Add(mBase, uint32(v108)+16)) = v104
											*(*int64)(unsafe.Add(mBase, uint32(v108)+8)) = v104
											*(*int64)(unsafe.Add(mBase, uint32(v108))) = v104
											v120 = v107 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v120) {
												v107 = v120
												v108 = v108 + int32(32)
												continue
											} else {
												break
											}
											break
										}
									}
								}
							}
						}
					}
				}
			}
		} else {
			base.MemoryFill(m, l0, v9, v17)
		}
		*(*int32)(unsafe.Add(mBase, uint32(l0))) = v8
		*(*int32)(unsafe.Add(mBase, uint32(l0)+3236)) = l1
	}
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3304)) = int32(16843009)
	v146 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3256)) = v146
	v150 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(l0+int32(3248)))) = uint16(v150)
	v154 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3308)))) = uint8(v154)
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3264)))) = v146
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3272)))) = v146
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3280)))) = v146
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3288)))) = v146
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(3296)))) = v146
	return
}
func F_VP8LHistogramSetClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
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
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v27 int32
	_ = v27
	var v38 int32
	_ = v38
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v72 int32
	_ = v72
	var v73 int32
	_ = v73
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v114 int64
	_ = v114
	var v117 int32
	_ = v117
	var v118 int32
	_ = v118
	var v130 int32
	_ = v130
	var v150 int32
	_ = v150
	var v155 int32
	_ = v155
	var v159 int32
	_ = v159
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v181 int32
	_ = v181
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v200 int32
	_ = v200
	var v202 int32
	_ = v202
	var v204 int32
	_ = v204
	var v210 int32
	_ = v210
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v221 int32
	_ = v221
	var v223 int32
	_ = v223
	var v227 int32
	_ = v227
	var v232 int32
	_ = v232
	var v234 int32
	_ = v234
	var v236 int32
	_ = v236
	var v241 int32
	_ = v241
	var v246 int32
	_ = v246
	var v253 int32
	_ = v253
	var v259 int32
	_ = v259
	var v263 int32
	_ = v263
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v271 int32
	_ = v271
	var v275 int32
	_ = v275
	var v280 int32
	_ = v280
	var v288 int32
	_ = v288
	var v297 int32
	_ = v297
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v312 int32
	_ = v312
	v2 = int32(0)
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(v15)+3236))
	v17 = int32(4) << (uint(v16) % 32)
	v18 = int32(_a_F_VP8LHistogramSetClear_0)
	v22 = base.B2i32(v2 < v16)
	if v2 < v16 {
		v23 = v17 + v18
	} else {
		v23 = v18
	}
	v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v27 = v23*v24 + int32(12)
	if base.Ui32(v27) < base.Ui32(int32(33)) {
		if v27 == int32(0) {
		} else {
			*(*uint8)(unsafe.Add(mBase, uint32(l0))) = uint8(v2)
			v38 = l0 + v27
			*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-1)))) = uint8(v2)
			if base.Ui32(v27) < base.Ui32(int32(3)) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+2)) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(l0)+1)) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-3)))) = uint8(v2)
				*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-2)))) = uint8(v2)
				if base.Ui32(v27) < base.Ui32(int32(7)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(l0)+3)) = uint8(v2)
					*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(-4)))) = uint8(v2)
					if base.Ui32(v27) < base.Ui32(int32(9)) {
					} else {
						v60 = int32(0)
						v63 = (v60 - l0) & int32(3)
						v64 = l0 + v63
						*(*int32)(unsafe.Add(mBase, uint32(v64))) = v60
						v72 = (v27 - v63) & int32(60)
						v73 = v64 + v72
						*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-4)))) = v60
						if base.Ui32(v72) < base.Ui32(int32(9)) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(v64)+8)) = v60
							*(*int32)(unsafe.Add(mBase, uint32(v64)+4)) = v60
							*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-8)))) = v60
							*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-12)))) = v60
							if base.Ui32(v72) < base.Ui32(int32(25)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v64)+24)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v64)+20)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v64)+16)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v64)+12)) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-16)))) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-20)))) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-24)))) = v60
								*(*int32)(unsafe.Add(mBase, uint32(v73+int32(-28)))) = v60
								v108 = v64&int32(4) | int32(24)
								v109 = v72 - v108
								if base.Ui32(v109) < base.Ui32(int32(32)) {
								} else {
									v114 = base.I64_extend_i32_u(v60) * int64(4294967297)
									v117 = v109
									v118 = v64 + v108
									for {
										*(*int64)(unsafe.Add(mBase, uint32(v118)+24)) = v114
										*(*int64)(unsafe.Add(mBase, uint32(v118)+16)) = v114
										*(*int64)(unsafe.Add(mBase, uint32(v118)+8)) = v114
										*(*int64)(unsafe.Add(mBase, uint32(v118))) = v114
										v130 = v117 + int32(-32)
										if base.Ui32(int32(31)) < base.Ui32(v130) {
											v117 = v130
											v118 = v118 + int32(32)
											continue
										} else {
											break
										}
										break
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		base.MemoryFill(m, l0, v2, v27)
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v24
	v150 = l0 + int32(12)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v150
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v24
	if v24 < int32(1) {
	} else {
		v155 = int32(1)
		v159 = v150 + v24<<(uint(int32(2))%32)
		if v24 == v155 {
			v215 = v2
			v221 = v150
			v223 = v159
		} else {
			v162 = int32(_a_F_VP8LHistogramSetClear_1)
			if v2 < v16 {
				v165 = v17 + v162
			} else {
				v165 = v162
			}
			v170 = int32(4)
			v171 = int32(0)
			v177 = v150
			v179 = v159
			for {
				v181 = int32(-4)
				v183 = int32(31)
				v185 = int32(-32)
				v186 = (v179 + v183) & v185
				*(*int32)(unsafe.Add(mBase, uint32(v177+v170+v181))) = v186
				v188 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v189 = v188 + v170
				v192 = *(*int32)(unsafe.Add(mBase, uint32(v189+v181)))
				v193 = int32(3312)
				*(*int32)(unsafe.Add(mBase, uint32(v192))) = v186 + v193
				v200 = (v186 + v165 + v183) & v185
				*(*int32)(unsafe.Add(mBase, uint32(v189))) = v200
				v202 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
				v204 = *(*int32)(unsafe.Add(mBase, uint32(v202+v170)))
				*(*int32)(unsafe.Add(mBase, uint32(v204))) = v200 + v193
				v210 = v200 + v165
				v212 = v171 + int32(2)
				if v24&int32(2147483646) != v212 {
					v170 = v170 + int32(8)
					v171 = v212
					v177 = v202
					v179 = v210
					continue
				} else {
					break
				}
				break
			}
			v215 = v212
			v221 = v202
			v223 = v210
		}
		if v24&v155 == int32(0) {
			v241 = v221
		} else {
			v227 = v215 << (uint(int32(2)) % 32)
			v232 = (v223 + int32(31)) & int32(-32)
			*(*int32)(unsafe.Add(mBase, uint32(v221+v227))) = v232
			v234 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			v236 = *(*int32)(unsafe.Add(mBase, uint32(v234+v227)))
			*(*int32)(unsafe.Add(mBase, uint32(v236))) = v232 + int32(3312)
			v241 = v234
		}
		if v24 < int32(1) {
		} else {
			v246 = v24 & int32(3)
			if base.Ui32(v24) < base.Ui32(int32(4)) {
				v288 = int32(0)
			} else {
				v253 = v241
				v259 = int32(0)
				for {
					v263 = *(*int32)(unsafe.Add(mBase, uint32(v253)))
					*(*int32)(unsafe.Add(mBase, uint32(v263)+3236)) = v16
					v265 = int32(4)
					v267 = *(*int32)(unsafe.Add(mBase, uint32(v253+v265)))
					*(*int32)(unsafe.Add(mBase, uint32(v267)+3236)) = v16
					v271 = *(*int32)(unsafe.Add(mBase, uint32(v253+int32(8))))
					*(*int32)(unsafe.Add(mBase, uint32(v271)+3236)) = v16
					v275 = *(*int32)(unsafe.Add(mBase, uint32(v253+int32(12))))
					*(*int32)(unsafe.Add(mBase, uint32(v275)+3236)) = v16
					v280 = v259 + v265
					if v24&int32(2147483644) != v280 {
						v253 = v253 + int32(16)
						v259 = v280
						continue
					} else {
						break
					}
					break
				}
				v288 = v280
			}
			if v246 == int32(0) {
			} else {
				v297 = v241 + v288<<(uint(int32(2))%32)
				v306 = v246
				for {
					v307 = *(*int32)(unsafe.Add(mBase, uint32(v297)))
					*(*int32)(unsafe.Add(mBase, uint32(v307)+3236)) = v16
					v312 = v306 + int32(-1)
					if v312 != 0 {
						v297 = v297 + int32(4)
						v306 = v312
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
func F_VP8LHistogramStoreRefs(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	v6 = m.G0
	v8 = v6 - int32(16)
	m.G0 = v8
	v11 = v8 + int32(4)
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v11)+4)) = v13
	if v13 != 0 {
		v17 = *(*int32)(unsafe.Add(mBase, uint32(v13)+4))
		v18 = *(*int32)(unsafe.Add(mBase, uint32(v13)+8))
		v22 = v17 + v18<<(uint(int32(3))%32)
		v23 = v17
	} else {
		v15 = int32(0)
		v22 = v15
		v23 = v15
	}
	*(*int32)(unsafe.Add(mBase, uint32(v11)+8)) = v22
	*(*int32)(unsafe.Add(mBase, uint32(v11))) = v23
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
	if v26 == int32(0) {
	} else {
		v29 = v26
		for {
			F_HistogramAddSinglePixOrCopy(m, l3, v29, l1, l2)
			mBase = m.M
			v35 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
			v37 = v35 + int32(8)
			*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v37
			v39 = *(*int32)(unsafe.Add(mBase, uint32(v8)+12))
			if v37 != v39 {
				v61 = v37
			} else {
				v42 = v8 + int32(4)
				v46 = *(*int32)(unsafe.Add(mBase, uint32(v42)+4))
				v47 = *(*int32)(unsafe.Add(mBase, uint32(v46)))
				if v47 != 0 {
					v50 = *(*int32)(unsafe.Add(mBase, uint32(v47)+4))
					v51 = *(*int32)(unsafe.Add(mBase, uint32(v47)+8))
					v55 = v50
					v56 = v50 + v51<<(uint(int32(3))%32)
				} else {
					v48 = int32(0)
					v55 = v48
					v56 = v48
				}
				*(*int32)(unsafe.Add(mBase, uint32(v42)+8)) = v56
				*(*int32)(unsafe.Add(mBase, uint32(v42))) = v55
				*(*int32)(unsafe.Add(mBase, uint32(v42)+4)) = v47
				v60 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				v61 = v60
			}
			if v61 != 0 {
				v29 = v61
				continue
			} else {
				break
			}
			break
		}
	}
	m.G0 = v8 + int32(16)
	return
}
func F_VP8LOptimizeSampling(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
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
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v49 int32
	_ = v49
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v96 int32
	_ = v96
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v125 int32
	_ = v125
	var v163 int32
	_ = v163
	var v180 int32
	_ = v180
	var v198 int32
	_ = v198
	var v211 int32
	_ = v211
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v254 int32
	_ = v254
	var v257 int32
	_ = v257
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v296 int32
	_ = v296
	var v298 int32
	_ = v298
	var v301 int32
	_ = v301
	var v304 int32
	_ = v304
	var v312 int32
	_ = v312
	var v349 int32
	_ = v349
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v371 int32
	_ = v371
	var v377 int32
	_ = v377
	var v378 int32
	_ = v378
	var v385 int32
	_ = v385
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v409 int32
	_ = v409
	var v422 int32
	_ = v422
	var v430 int32
	_ = v430
	var v432 int32
	_ = v432
	var v435 int32
	_ = v435
	var v437 int32
	_ = v437
	var v445 int32
	_ = v445
	var v455 int32
	_ = v455
	var v465 int32
	_ = v465
	var v470 int32
	_ = v470
	var v487 int32
	_ = v487
	var v503 int32
	_ = v503
	var v507 int32
	_ = v507
	var v516 int32
	_ = v516
	var v528 int32
	_ = v528
	var v535 int32
	_ = v535
	var v562 int32
	_ = v562
	*(*int32)(unsafe.Add(mBase, uint32(l5))) = l3
	if l4 <= l3 {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return
L2:
	;
	v27 = int32(1) << (uint(l3) % 32)
	v28 = int32(-1)
	v29 = l2 + v28
	v31 = int32(base.Ui32(v27+v29) >> (uint(l3) % 32))
	v33 = l1 + v28
	v35 = int32(base.Ui32(v27+v33) >> (uint(l3) % 32))
	v37 = v35 << (uint(int32(2)) % 32)
	v49 = l3
	goto L4
L3:
	;
	if v163 == l3 {
		goto L1
	} else {
		goto L21
	}
L4:
	;
	v61 = int32(1)
	v62 = v49 + v61
	v65 = v61 << (uint(v49-l3) % 32)
	if v31 <= v65 {
		goto L6
	} else {
		goto L7
	}
L5:
	;
	v163 = l4
	goto L3
L6:
	;
	if v62 != l4 {
		v49 = v62
		goto L4
	} else {
		goto L20
	}
L7:
	;
	v69 = int32(1) << (uint(v62-l3) % 32)
	v73 = v65
	v74 = l0
	goto L8
L8:
	;
	v96 = int32(0)
	if v37 == v96 {
		v121 = v96
		goto L11
	} else {
		goto L12
	}
L9:
	;
	goto L6
L10:
	;
	if v121 != 0 {
		v163 = v49
		goto L3
	} else {
		goto L18
	}
L11:
	;
	goto L10
L12:
	;
	v102 = v74
	v103 = v74 + v37*v65
	v104 = v37
	goto L14
L13:
	;
	v121 = v108 - v109
	goto L11
L14:
	;
	v108 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v102))))
	v109 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v103))))
	if v108 != v109 {
		goto L13
	} else {
		goto L16
	}
L16:
	;
	v111 = int32(1)
	v116 = v104 + int32(-1)
	if v116 != 0 {
		v102 = v102 + v111
		v103 = v103 + v111
		v104 = v116
		goto L14
	} else {
		goto L17
	}
L17:
	;
	v121 = v96
	goto L11
L18:
	;
	v125 = v73 + v69
	if v125 < v31 {
		v73 = v125
		v74 = v74 + v37*v69
		goto L8
	} else {
		goto L19
	}
L19:
	;
	goto L9
L20:
	;
	goto L5
L21:
	;
	if v31 < int32(1) {
		v349 = v163
		goto L22
	} else {
		goto L23
	}
L22:
	;
	if v349 == l3 {
		goto L1
	} else {
		goto L49
	}
L23:
	;
	if v35 < int32(1) {
		v349 = v163
		goto L22
	} else {
		goto L24
	}
L24:
	;
	if v163 < l3 {
		goto L25
	} else {
		goto L26
	}
L25:
	;
	v180 = v163
	goto L27
L26:
	;
	v180 = l3
	goto L27
L27:
	;
	if v163 <= l3 {
		goto L28
	} else {
		goto L29
	}
L28:
	;
	v349 = v180
	goto L22
L29:
	;
	v198 = v163
	goto L30
L30:
	;
	v211 = int32(1) << (uint(v198-l3) % 32)
	v236 = int32(0)
	v237 = l0 + int32(4)
	goto L32
L32:
	;
	v254 = v237
	v257 = int32(0)
	goto L34
L34:
	;
	v266 = v257 + v211
	if v266 < v35 {
		goto L36
	} else {
		goto L37
	}
L35:
	;
	v312 = v236 + int32(1)
	if v312 < v31 {
		v236 = v312
		v237 = v237 + v35<<(uint(int32(2))%32)
		goto L32
	} else {
		goto L48
	}
L36:
	;
	v268 = v266
	goto L38
L37:
	;
	v268 = v35
	goto L38
L38:
	;
	v273 = v257
	v274 = v254
	goto L40
L39:
	;
	if v296 < v268 {
		goto L45
	} else {
		goto L46
	}
L40:
	;
	v296 = v273 + int32(1)
	if v268 <= v296 {
		goto L39
	} else {
		goto L42
	}
L41:
	;
	v304 = v198 + int32(-1)
	if l3 < v304 {
		v198 = v304
		goto L30
	} else {
		goto L44
	}
L42:
	;
	v298 = *(*int32)(unsafe.Add(mBase, uint32(v274)))
	v301 = *(*int32)(unsafe.Add(mBase, uint32(l0+v236*v35<<(uint(int32(2))%32)+v257<<(uint(int32(2))%32))))
	if v298 == v301 {
		v273 = v296
		v274 = v274 + int32(4)
		goto L40
	} else {
		goto L43
	}
L43:
	;
	goto L41
L44:
	;
	goto L28
L45:
	;
	goto L35
L46:
	;
	if v266 < v35 {
		v254 = v254 + v211<<(uint(int32(2))%32)
		v257 = v266
		goto L34
	} else {
		goto L47
	}
L47:
	;
	goto L45
L48:
	;
	v349 = v198
	goto L22
L49:
	;
	v361 = int32(1)
	v362 = v361 << (uint(v349) % 32)
	v364 = int32(base.Ui32(v362+v29) >> (uint(v349) % 32))
	if v364 < v361 {
		goto L50
	} else {
		goto L51
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l5))) = v349
	goto L1
L51:
	;
	v368 = int32(base.Ui32(v362+v33) >> (uint(v349) % 32))
	if v368 < int32(1) {
		goto L50
	} else {
		goto L52
	}
L52:
	;
	v371 = v349 - l3
	v377 = v368 & int32(3)
	v378 = int32(0)
	v385 = v378
	v395 = v378
	v396 = l0
	goto L53
L53:
	;
	if base.Ui32(v368) < base.Ui32(int32(4)) {
		v487 = int32(0)
		goto L55
	} else {
		goto L56
	}
L54:
	;
	goto L50
L55:
	;
	if v377 == int32(0) {
		goto L60
	} else {
		goto L61
	}
L56:
	;
	v409 = v396
	v422 = int32(0)
	goto L57
L57:
	;
	v430 = v395 + v422
	v432 = int32(2)
	v435 = *(*int32)(unsafe.Add(mBase, uint32(l0+v430<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409))) = v435
	v437 = int32(4)
	v445 = *(*int32)(unsafe.Add(mBase, uint32(l0+(v430+int32(1))<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409+v437))) = v445
	v455 = *(*int32)(unsafe.Add(mBase, uint32(l0+(v430+v432)<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409+int32(8)))) = v455
	v465 = *(*int32)(unsafe.Add(mBase, uint32(l0+(v430+int32(3))<<(uint(v371)%32)<<(uint(v432)%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v409+int32(12)))) = v465
	v470 = v422 + v437
	if v368&int32(2147483644) != v470 {
		v409 = v409 + int32(16)
		v422 = v470
		goto L57
	} else {
		goto L59
	}
L58:
	;
	v487 = v470
	goto L55
L59:
	;
	goto L58
L60:
	;
	v562 = v385 + int32(1)
	if v562 != v364 {
		v385 = v562
		v395 = v395 + v35
		v396 = v396 + v368<<(uint(int32(2))%32)
		goto L53
	} else {
		goto L65
	}
L61:
	;
	v503 = v487 + v395
	v507 = v396 + v487<<(uint(int32(2))%32)
	v516 = v377
	goto L62
L62:
	;
	v528 = *(*int32)(unsafe.Add(mBase, uint32(l0+v503<<(uint(v371)%32)<<(uint(int32(2))%32))))
	*(*int32)(unsafe.Add(mBase, uint32(v507))) = v528
	v535 = v516 + int32(-1)
	if v535 != 0 {
		v503 = v503 + int32(1)
		v507 = v507 + int32(4)
		v516 = v535
		goto L62
	} else {
		goto L64
	}
L63:
	;
	goto L60
L64:
	;
	goto L63
L65:
	;
	goto L54
}
func F_VP8LPutBitsInternal(m *base.Module, l0 int32, l1 int32, l2 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
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
	var v35 int32
	_ = v35
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v55 int64
	_ = v55
	var v58 int32
	_ = v58
	var v60 int64
	_ = v60
	var v63 int32
	_ = v63
	var v66 int64
	_ = v66
	var v70 int32
	_ = v70
	var v72 int32
	_ = v72
	var v76 int32
	_ = v76
	var v83 int64
	_ = v83
	var v84 int32
	_ = v84
	var v85 int64
	_ = v85
	var v95 int32
	_ = v95
	var v97 int32
	_ = v97
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	var v126 int32
	_ = v126
	var v129 int32
	_ = v129
	var v131 int32
	_ = v131
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	if l2 < int32(1) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return
L2:
	;
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	if v16+l2 < int32(32) {
		goto L5
	} else {
		goto L6
	}
L3:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v138 + v136
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v135<<(uint(v138)%32) | v137
	goto L1
L4:
	;
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v39 = v31
	v40 = v32
	v41 = v35
	v42 = v34
	goto L8
L5:
	;
	if v16 < int32(16) {
		v135 = l1
		v136 = l2
		v137 = v15
		v138 = v16
		goto L3
	} else {
		goto L7
	}
L6:
	;
	v20 = int32(32)
	v21 = v20 - v16
	v29 = int32(base.Ui32(l1) >> (uint(v21) % 32))
	v30 = l2 - v21
	v31 = l1<<(uint(v16)%32) | v15
	v32 = v20
	goto L4
L7:
	;
	v29 = l1
	v30 = l2
	v31 = v15
	v32 = v16
	goto L4
L8:
	;
	if base.Ui32(v41+int32(2)) <= base.Ui32(v42) {
		v118 = v41
		v119 = v42
		goto L10
	} else {
		goto L11
	}
L9:
	;
	v135 = v29
	v136 = v30
	v137 = v131
	v138 = v129
	goto L3
L10:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v118))) = uint16(v39)
	v126 = v118 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v126
	v129 = v40 + int32(-16)
	v131 = int32(base.Ui32(v39) >> (uint(int32(16)) % 32))
	if int32(31) < v40 {
		v39 = v131
		v40 = v129
		v41 = v126
		v42 = v119
		goto L8
	} else {
		goto L32
	}
L11:
	;
	v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v52 = v42 - v51
	v55 = base.I64_extend_i32_s(v52) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v55) {
		v100 = v51
		goto L13
	} else {
		goto L14
	}
L12:
	;
	if v41 == v51 {
		goto L29
	} else {
		goto L30
	}
L13:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+12)) = v100
	return
L14:
	;
	v58 = v41 - v51
	v60 = v55 + base.I64_extend_i32_u(v58)
	if base.Ui64(int64(4294967295)) < base.Ui64(v60) {
		v100 = v51
		goto L13
	} else {
		goto L15
	}
L15:
	;
	v63 = base.I32_wrap_i64(v60)
	if v42 == v51 {
		goto L16
	} else {
		goto L17
	}
L16:
	;
	v66 = int64(1)
	v70 = int32(base.Ui32(v52*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v63) < base.Ui32(v70) {
		goto L19
	} else {
		goto L20
	}
L17:
	;
	if base.Ui32(v63) <= base.Ui32(v52) {
		v118 = v41
		v119 = v42
		goto L10
	} else {
		goto L18
	}
L18:
	;
	goto L16
L19:
	;
	v72 = v70
	goto L21
L20:
	;
	v72 = v63
	goto L21
L21:
	;
	v76 = v72&int32(-1024) + int32(1024)
	goto L25
L22:
	;
	if v97 != 0 {
		goto L12
	} else {
		goto L28
	}
L23:
	;
	goto L22
L24:
	;
	v95 = F_malloc(m, base.I32_wrap_i64(v66)*v76)
	mBase = m.M
	v97 = v95
	goto L23
L25:
	;
	v83 = base.I64_div_u_s(int64(2147418112), v66)
	v84 = int32(0)
	v85 = base.I64_extend_i32_u(v76)
	if base.Ui64(int64(4294967295)) < base.Ui64(v85*v66) {
		v97 = v84
		goto L23
	} else {
		goto L26
	}
L26:
	;
	if base.Ui64(v83) < base.Ui64(v85) {
		v97 = v84
		goto L23
	} else {
		goto L27
	}
L27:
	;
	goto L24
L28:
	;
	v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v100 = v99
	goto L13
L29:
	;
	v111 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	F_free(m, v111)
	mBase = m.M
	goto L31
L30:
	;
	v109 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v110 = F_memcpy(m, v97, v109, v58)
	mBase = m.M
	goto L29
L31:
	;
	v113 = v97 + v76
	*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v113
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v97
	v118 = v97 + v58
	v119 = v113
	goto L10
L32:
	;
	goto L9
}
func F_VP8LRefsCursorInit(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	v4 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v4
	if v4 != 0 {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(v4)+4))
		v9 = *(*int32)(unsafe.Add(mBase, uint32(v4)+8))
		v13 = v8 + v9<<(uint(int32(3))%32)
		v14 = v8
	} else {
		v6 = int32(0)
		v13 = v6
		v14 = v6
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v13
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v14
	return
}
func F_VP8LRefsCursorNextBlock(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v6 = *(*int32)(unsafe.Add(mBase, uint32(v5)))
	if v6 != 0 {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(v6)+4))
		v10 = *(*int32)(unsafe.Add(mBase, uint32(v6)+8))
		v14 = v9
		v15 = v9 + v10<<(uint(int32(3))%32)
	} else {
		v7 = int32(0)
		v14 = v7
		v15 = v7
	}
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = v15
	*(*int32)(unsafe.Add(mBase, uint32(l0))) = v14
	*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v6
	return
}
func F_VP8LResidualImage(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32) int32 {
	mBase := m.M
	_ = mBase
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
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
	var v115 int32
	_ = v115
	var v119 int32
	_ = v119
	var v134 int32
	_ = v134
	var v148 int32
	_ = v148
	var v217 int64
	_ = v217
	var v227 int32
	_ = v227
	var v234 int32
	_ = v234
	var v258 int32
	_ = v258
	var v340 int32
	_ = v340
	var v353 int32
	_ = v353
	var v428 int32
	_ = v428
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v526 int32
	_ = v526
	var v528 int32
	_ = v528
	var v529 int32
	_ = v529
	var v548 int32
	_ = v548
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v566 int32
	_ = v566
	var v631 int32
	_ = v631
	var v632 int32
	_ = v632
	var v637 int32
	_ = v637
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v649 int32
	_ = v649
	var v652 int32
	_ = v652
	var v654 int32
	_ = v654
	var v658 int32
	_ = v658
	var v667 int32
	_ = v667
	var v682 int32
	_ = v682
	var v756 int32
	_ = v756
	var v761 int32
	_ = v761
	var v787 int32
	_ = v787
	var v855 int64
	_ = v855
	var v856 int32
	_ = v856
	var v863 int64
	_ = v863
	var v864 int32
	_ = v864
	var v865 int64
	_ = v865
	var v875 int32
	_ = v875
	var v877 int32
	_ = v877
	var v885 int32
	_ = v885
	var v890 int32
	_ = v890
	var v894 int32
	_ = v894
	var v909 int32
	_ = v909
	var v922 int32
	_ = v922
	var v923 int32
	_ = v923
	var v924 int32
	_ = v924
	var v925 int32
	_ = v925
	var v992 int32
	_ = v992
	var v995 int32
	_ = v995
	var v997 int32
	_ = v997
	var v1002 int32
	_ = v1002
	var v1004 int32
	_ = v1004
	var v1013 int32
	_ = v1013
	var v1028 int32
	_ = v1028
	var v1100 int32
	_ = v1100
	var v1112 int32
	_ = v1112
	var v1127 int32
	_ = v1127
	var v1128 int32
	_ = v1128
	var v1130 int32
	_ = v1130
	var v1195 int32
	_ = v1195
	var v1196 int32
	_ = v1196
	var v1198 int32
	_ = v1198
	var v1199 int32
	_ = v1199
	var v1200 int32
	_ = v1200
	var v1202 int32
	_ = v1202
	var v1204 int32
	_ = v1204
	var v1208 int32
	_ = v1208
	var v1211 int32
	_ = v1211
	var v1213 int32
	_ = v1213
	var v1217 int32
	_ = v1217
	var v1220 int32
	_ = v1220
	var v1222 int32
	_ = v1222
	var v1226 int32
	_ = v1226
	var v1229 int32
	_ = v1229
	var v1236 int32
	_ = v1236
	var v1328 int32
	_ = v1328
	var v1330 int32
	_ = v1330
	var v1333 int64
	_ = v1333
	var v1334 int32
	_ = v1334
	var v1341 int64
	_ = v1341
	var v1342 int32
	_ = v1342
	var v1343 int64
	_ = v1343
	var v1352 int32
	_ = v1352
	var v1354 int32
	_ = v1354
	var v1356 int32
	_ = v1356
	var v1357 int32
	_ = v1357
	var v1361 int32
	_ = v1361
	var v1362 int32
	_ = v1362
	var v1364 int32
	_ = v1364
	var v1366 int32
	_ = v1366
	var v1368 int32
	_ = v1368
	var v1373 int32
	_ = v1373
	var v1377 int32
	_ = v1377
	var v1380 int32
	_ = v1380
	var v1382 int32
	_ = v1382
	var v1383 int32
	_ = v1383
	var v1385 int32
	_ = v1385
	var v1386 int32
	_ = v1386
	var v1389 int32
	_ = v1389
	var v1391 int32
	_ = v1391
	var v1393 int32
	_ = v1393
	var v1394 int32
	_ = v1394
	var v1397 int32
	_ = v1397
	var v1406 int32
	_ = v1406
	var v1411 int32
	_ = v1411
	var v1412 int32
	_ = v1412
	var v1415 int32
	_ = v1415
	var v1416 int32
	_ = v1416
	var v1473 int32
	_ = v1473
	var v1474 int32
	_ = v1474
	var v1475 int32
	_ = v1475
	var v1476 int32
	_ = v1476
	var v1477 int32
	_ = v1477
	var v1478 int32
	_ = v1478
	var v1513 int32
	_ = v1513
	var v1515 int32
	_ = v1515
	var v1516 int32
	_ = v1516
	var v1517 int32
	_ = v1517
	var v1518 int32
	_ = v1518
	var v1519 int32
	_ = v1519
	var v1520 int32
	_ = v1520
	var v1526 int32
	_ = v1526
	var v1527 int32
	_ = v1527
	var v1528 int32
	_ = v1528
	var v1529 int32
	_ = v1529
	var v1530 int32
	_ = v1530
	var v1534 int32
	_ = v1534
	var v1535 int32
	_ = v1535
	var v1537 int32
	_ = v1537
	var v1540 int32
	_ = v1540
	var v1541 int32
	_ = v1541
	var v1542 int32
	_ = v1542
	var v1554 int32
	_ = v1554
	var v1573 int32
	_ = v1573
	var v1614 int32
	_ = v1614
	var v1640 int32
	_ = v1640
	var v1643 int32
	_ = v1643
	var v1653 int32
	_ = v1653
	var __phi1653 int32
	_ = __phi1653
	var v1713 int32
	_ = v1713
	var __phi1713 int32
	_ = __phi1713
	var v1714 int32
	_ = v1714
	var __phi1714 int32
	_ = __phi1714
	var v1737 int32
	_ = v1737
	var v1739 int32
	_ = v1739
	var v1741 int32
	_ = v1741
	var v1743 int32
	_ = v1743
	var v1748 int32
	_ = v1748
	var v1768 int32
	_ = v1768
	var v1769 int32
	_ = v1769
	var v1772 int32
	_ = v1772
	var v1773 int32
	_ = v1773
	var v1774 int32
	_ = v1774
	var v1776 int32
	_ = v1776
	var v1779 int32
	_ = v1779
	var v1784 int32
	_ = v1784
	var v1788 int32
	_ = v1788
	var v1802 int32
	_ = v1802
	var v1803 int32
	_ = v1803
	var v1813 int32
	_ = v1813
	var v1817 int32
	_ = v1817
	var __phi1817 int32
	_ = __phi1817
	var v1820 int32
	_ = v1820
	var __phi1820 int32
	_ = __phi1820
	var v1821 int32
	_ = v1821
	var __phi1821 int32
	_ = __phi1821
	var v1823 int32
	_ = v1823
	var __phi1823 int32
	_ = __phi1823
	var v1824 int32
	_ = v1824
	var __phi1824 int32
	_ = __phi1824
	var v1833 int32
	_ = v1833
	var v1834 int32
	_ = v1834
	var v1836 int32
	_ = v1836
	var v1838 int32
	_ = v1838
	var v1841 int32
	_ = v1841
	var v1842 int32
	_ = v1842
	var v1843 int32
	_ = v1843
	var v1845 int32
	_ = v1845
	var v1848 int32
	_ = v1848
	var v1853 int32
	_ = v1853
	var v1857 int32
	_ = v1857
	var v1872 int32
	_ = v1872
	var v1886 int32
	_ = v1886
	var v1888 int32
	_ = v1888
	var v1889 int32
	_ = v1889
	var v1890 int32
	_ = v1890
	var v1891 int32
	_ = v1891
	var v1894 int32
	_ = v1894
	var v1895 int32
	_ = v1895
	var v1896 int32
	_ = v1896
	var v1898 int32
	_ = v1898
	var v1899 int32
	_ = v1899
	var v1901 int32
	_ = v1901
	var v1902 int32
	_ = v1902
	var v1907 int32
	_ = v1907
	var v1909 int32
	_ = v1909
	var v1911 int32
	_ = v1911
	var v1913 int32
	_ = v1913
	var v1914 int32
	_ = v1914
	var v1916 int32
	_ = v1916
	var v1917 int32
	_ = v1917
	var v1922 int32
	_ = v1922
	var v1923 int32
	_ = v1923
	var v1924 int32
	_ = v1924
	var v1926 int32
	_ = v1926
	var v1928 int32
	_ = v1928
	var v1931 int32
	_ = v1931
	var v1933 int32
	_ = v1933
	var v1935 int32
	_ = v1935
	var v1937 int32
	_ = v1937
	var v1939 int32
	_ = v1939
	var v1942 int32
	_ = v1942
	var v1943 int32
	_ = v1943
	var v1944 int32
	_ = v1944
	var v1946 int32
	_ = v1946
	var v1951 int32
	_ = v1951
	var v1953 int32
	_ = v1953
	var v1955 int32
	_ = v1955
	var v1957 int32
	_ = v1957
	var v1960 int32
	_ = v1960
	var v1962 int32
	_ = v1962
	var v1963 int32
	_ = v1963
	var v1964 int32
	_ = v1964
	var v1966 int32
	_ = v1966
	var v1969 int32
	_ = v1969
	var v1971 int32
	_ = v1971
	var v1973 int32
	_ = v1973
	var v1975 int32
	_ = v1975
	var v1977 int32
	_ = v1977
	var v1979 int32
	_ = v1979
	var v1982 int32
	_ = v1982
	var v1983 int32
	_ = v1983
	var v1984 int32
	_ = v1984
	var v1986 int32
	_ = v1986
	var v1991 int32
	_ = v1991
	var v1993 int32
	_ = v1993
	var v1995 int32
	_ = v1995
	var v1997 int32
	_ = v1997
	var v2000 int32
	_ = v2000
	var v2002 int32
	_ = v2002
	var v2003 int32
	_ = v2003
	var v2004 int32
	_ = v2004
	var v2006 int32
	_ = v2006
	var v2009 int32
	_ = v2009
	var v2011 int32
	_ = v2011
	var v2013 int32
	_ = v2013
	var v2015 int32
	_ = v2015
	var v2017 int32
	_ = v2017
	var v2020 int32
	_ = v2020
	var v2021 int32
	_ = v2021
	var v2022 int32
	_ = v2022
	var v2024 int32
	_ = v2024
	var v2029 int32
	_ = v2029
	var v2031 int32
	_ = v2031
	var v2033 int32
	_ = v2033
	var v2035 int32
	_ = v2035
	var v2038 int32
	_ = v2038
	var v2040 int32
	_ = v2040
	var v2041 int32
	_ = v2041
	var v2042 int32
	_ = v2042
	var v2044 int32
	_ = v2044
	var v2047 int32
	_ = v2047
	var v2049 int32
	_ = v2049
	var v2051 int32
	_ = v2051
	var v2053 int32
	_ = v2053
	var v2055 int32
	_ = v2055
	var v2057 int32
	_ = v2057
	var v2059 int32
	_ = v2059
	var v2064 int32
	_ = v2064
	var v2083 int32
	_ = v2083
	var v2108 int32
	_ = v2108
	var v2109 int32
	_ = v2109
	var v2178 int32
	_ = v2178
	var v2181 int32
	_ = v2181
	var v2183 int32
	_ = v2183
	var v2184 int32
	_ = v2184
	var v2185 int32
	_ = v2185
	var v2194 int32
	_ = v2194
	var v2195 int32
	_ = v2195
	var v2205 int32
	_ = v2205
	var v2206 int32
	_ = v2206
	var v2216 int32
	_ = v2216
	var v2217 int32
	_ = v2217
	var v2224 int32
	_ = v2224
	var v2253 int32
	_ = v2253
	var v2323 int32
	_ = v2323
	var v2348 int32
	_ = v2348
	var v2349 int32
	_ = v2349
	var v2417 int32
	_ = v2417
	var v2420 int32
	_ = v2420
	var v2422 int32
	_ = v2422
	var v2423 int32
	_ = v2423
	var v2424 int32
	_ = v2424
	var v2433 int32
	_ = v2433
	var v2434 int32
	_ = v2434
	var v2444 int32
	_ = v2444
	var v2445 int32
	_ = v2445
	var v2455 int32
	_ = v2455
	var v2456 int32
	_ = v2456
	var v2463 int32
	_ = v2463
	var v2465 int32
	_ = v2465
	var v2559 int32
	_ = v2559
	var v2562 int32
	_ = v2562
	var v2568 int32
	_ = v2568
	var v2599 int32
	_ = v2599
	var v2600 int32
	_ = v2600
	var v2602 int32
	_ = v2602
	var v2670 int32
	_ = v2670
	var v2679 int32
	_ = v2679
	var __phi2679 int32
	_ = __phi2679
	var v2692 int32
	_ = v2692
	var __phi2692 int32
	_ = __phi2692
	var v2694 int32
	_ = v2694
	var __phi2694 int32
	_ = __phi2694
	var v2695 int32
	_ = v2695
	var __phi2695 int32
	_ = __phi2695
	var v2696 int32
	_ = v2696
	var __phi2696 int32
	_ = __phi2696
	var v2764 int32
	_ = v2764
	var v2769 int32
	_ = v2769
	var v2775 int32
	_ = v2775
	var v2777 int32
	_ = v2777
	var v2781 int32
	_ = v2781
	var v2811 int32
	_ = v2811
	var v2812 int32
	_ = v2812
	var v2815 int32
	_ = v2815
	var v2883 int32
	_ = v2883
	var v2892 int32
	_ = v2892
	var __phi2892 int32
	_ = __phi2892
	var v2906 int32
	_ = v2906
	var __phi2906 int32
	_ = __phi2906
	var v2907 int32
	_ = v2907
	var __phi2907 int32
	_ = __phi2907
	var v2908 int32
	_ = v2908
	var __phi2908 int32
	_ = __phi2908
	var v2909 int32
	_ = v2909
	var __phi2909 int32
	_ = __phi2909
	var v2976 int32
	_ = v2976
	var v2977 int32
	_ = v2977
	var v2982 int32
	_ = v2982
	var v3000 int32
	_ = v3000
	var v3001 int32
	_ = v3001
	var v3004 int32
	_ = v3004
	var v3005 int32
	_ = v3005
	var v3006 int32
	_ = v3006
	var v3008 int32
	_ = v3008
	var v3011 int32
	_ = v3011
	var v3016 int32
	_ = v3016
	var v3020 int32
	_ = v3020
	var v3034 int32
	_ = v3034
	var v3035 int32
	_ = v3035
	var v3045 int32
	_ = v3045
	var v3049 int32
	_ = v3049
	var __phi3049 int32
	_ = __phi3049
	var v3052 int32
	_ = v3052
	var __phi3052 int32
	_ = __phi3052
	var v3053 int32
	_ = v3053
	var __phi3053 int32
	_ = __phi3053
	var v3055 int32
	_ = v3055
	var __phi3055 int32
	_ = __phi3055
	var v3056 int32
	_ = v3056
	var __phi3056 int32
	_ = __phi3056
	var v3065 int32
	_ = v3065
	var v3066 int32
	_ = v3066
	var v3068 int32
	_ = v3068
	var v3070 int32
	_ = v3070
	var v3073 int32
	_ = v3073
	var v3074 int32
	_ = v3074
	var v3075 int32
	_ = v3075
	var v3077 int32
	_ = v3077
	var v3080 int32
	_ = v3080
	var v3085 int32
	_ = v3085
	var v3089 int32
	_ = v3089
	var v3104 int32
	_ = v3104
	var v3118 int32
	_ = v3118
	var v3120 int32
	_ = v3120
	var v3121 int32
	_ = v3121
	var v3122 int32
	_ = v3122
	var v3123 int32
	_ = v3123
	var v3126 int32
	_ = v3126
	var v3127 int32
	_ = v3127
	var v3128 int32
	_ = v3128
	var v3130 int32
	_ = v3130
	var v3131 int32
	_ = v3131
	var v3133 int32
	_ = v3133
	var v3134 int32
	_ = v3134
	var v3139 int32
	_ = v3139
	var v3141 int32
	_ = v3141
	var v3143 int32
	_ = v3143
	var v3145 int32
	_ = v3145
	var v3146 int32
	_ = v3146
	var v3148 int32
	_ = v3148
	var v3149 int32
	_ = v3149
	var v3154 int32
	_ = v3154
	var v3155 int32
	_ = v3155
	var v3156 int32
	_ = v3156
	var v3158 int32
	_ = v3158
	var v3160 int32
	_ = v3160
	var v3163 int32
	_ = v3163
	var v3165 int32
	_ = v3165
	var v3167 int32
	_ = v3167
	var v3169 int32
	_ = v3169
	var v3171 int32
	_ = v3171
	var v3174 int32
	_ = v3174
	var v3175 int32
	_ = v3175
	var v3176 int32
	_ = v3176
	var v3178 int32
	_ = v3178
	var v3183 int32
	_ = v3183
	var v3185 int32
	_ = v3185
	var v3187 int32
	_ = v3187
	var v3189 int32
	_ = v3189
	var v3192 int32
	_ = v3192
	var v3194 int32
	_ = v3194
	var v3195 int32
	_ = v3195
	var v3196 int32
	_ = v3196
	var v3198 int32
	_ = v3198
	var v3201 int32
	_ = v3201
	var v3203 int32
	_ = v3203
	var v3205 int32
	_ = v3205
	var v3207 int32
	_ = v3207
	var v3209 int32
	_ = v3209
	var v3211 int32
	_ = v3211
	var v3214 int32
	_ = v3214
	var v3215 int32
	_ = v3215
	var v3216 int32
	_ = v3216
	var v3218 int32
	_ = v3218
	var v3223 int32
	_ = v3223
	var v3225 int32
	_ = v3225
	var v3227 int32
	_ = v3227
	var v3229 int32
	_ = v3229
	var v3232 int32
	_ = v3232
	var v3234 int32
	_ = v3234
	var v3235 int32
	_ = v3235
	var v3236 int32
	_ = v3236
	var v3238 int32
	_ = v3238
	var v3241 int32
	_ = v3241
	var v3243 int32
	_ = v3243
	var v3245 int32
	_ = v3245
	var v3247 int32
	_ = v3247
	var v3249 int32
	_ = v3249
	var v3252 int32
	_ = v3252
	var v3253 int32
	_ = v3253
	var v3254 int32
	_ = v3254
	var v3256 int32
	_ = v3256
	var v3261 int32
	_ = v3261
	var v3263 int32
	_ = v3263
	var v3265 int32
	_ = v3265
	var v3267 int32
	_ = v3267
	var v3270 int32
	_ = v3270
	var v3272 int32
	_ = v3272
	var v3273 int32
	_ = v3273
	var v3274 int32
	_ = v3274
	var v3276 int32
	_ = v3276
	var v3279 int32
	_ = v3279
	var v3281 int32
	_ = v3281
	var v3283 int32
	_ = v3283
	var v3285 int32
	_ = v3285
	var v3287 int32
	_ = v3287
	var v3289 int32
	_ = v3289
	var v3291 int32
	_ = v3291
	var v3296 int32
	_ = v3296
	var v3319 int32
	_ = v3319
	var v3321 int32
	_ = v3321
	var v3327 int32
	_ = v3327
	var v3328 int32
	_ = v3328
	var v3329 int32
	_ = v3329
	var v3330 int32
	_ = v3330
	var v3331 int32
	_ = v3331
	var v3332 int32
	_ = v3332
	var v3333 int32
	_ = v3333
	var v3334 int32
	_ = v3334
	var v3335 int32
	_ = v3335
	var v3336 int32
	_ = v3336
	var v3337 int32
	_ = v3337
	var v3338 int32
	_ = v3338
	var v3339 int32
	_ = v3339
	var v3340 int32
	_ = v3340
	var v3432 int32
	_ = v3432
	var v3505 int32
	_ = v3505
	var v3511 int32
	_ = v3511
	var v3512 int32
	_ = v3512
	var v3513 int32
	_ = v3513
	var v3514 int32
	_ = v3514
	var v3515 int32
	_ = v3515
	var v3529 int32
	_ = v3529
	var v3530 int32
	_ = v3530
	var v3533 int32
	_ = v3533
	var v3534 int32
	_ = v3534
	var v3538 int32
	_ = v3538
	var v3539 int32
	_ = v3539
	var v3541 int32
	_ = v3541
	var v3545 int32
	_ = v3545
	var v3553 int32
	_ = v3553
	var v3554 int32
	_ = v3554
	var v3557 int32
	_ = v3557
	var v3563 int32
	_ = v3563
	var v3569 int32
	_ = v3569
	var v3570 int32
	_ = v3570
	var v3578 int32
	_ = v3578
	var v3579 int32
	_ = v3579
	var v3581 int32
	_ = v3581
	var v3591 int32
	_ = v3591
	var v3603 int32
	_ = v3603
	var v3607 int32
	_ = v3607
	var v3608 int32
	_ = v3608
	var v3609 int32
	_ = v3609
	var v3610 int32
	_ = v3610
	var v3611 int32
	_ = v3611
	var v3652 int32
	_ = v3652
	var v3664 int64
	_ = v3664
	var v3676 int32
	_ = v3676
	var v3677 int64
	_ = v3677
	var v3703 int32
	_ = v3703
	var v3704 int32
	_ = v3704
	var v3764 int64
	_ = v3764
	var v3765 int64
	_ = v3765
	var v3775 int32
	_ = v3775
	var v3777 int32
	_ = v3777
	var v3780 int64
	_ = v3780
	var v3783 int64
	_ = v3783
	var v3786 int64
	_ = v3786
	var v3787 int64
	_ = v3787
	var v3794 int64
	_ = v3794
	var v3799 int64
	_ = v3799
	var v3801 int32
	_ = v3801
	var v3806 int32
	_ = v3806
	var v3807 int32
	_ = v3807
	var v3808 int64
	_ = v3808
	var v3809 int64
	_ = v3809
	var v3835 int32
	_ = v3835
	var v3836 int32
	_ = v3836
	var v3897 int64
	_ = v3897
	var v3898 int64
	_ = v3898
	var v3907 int32
	_ = v3907
	var v3909 int32
	_ = v3909
	var v3912 int64
	_ = v3912
	var v3915 int64
	_ = v3915
	var v3918 int64
	_ = v3918
	var v3919 int64
	_ = v3919
	var v3926 int64
	_ = v3926
	var v3931 int64
	_ = v3931
	var v3933 int32
	_ = v3933
	var v3938 int32
	_ = v3938
	var v3939 int32
	_ = v3939
	var v3940 int64
	_ = v3940
	var v3941 int64
	_ = v3941
	var v3967 int32
	_ = v3967
	var v3968 int32
	_ = v3968
	var v4029 int64
	_ = v4029
	var v4032 int64
	_ = v4032
	var v4039 int32
	_ = v4039
	var v4041 int32
	_ = v4041
	var v4044 int64
	_ = v4044
	var v4047 int64
	_ = v4047
	var v4050 int64
	_ = v4050
	var v4051 int64
	_ = v4051
	var v4058 int64
	_ = v4058
	var v4063 int64
	_ = v4063
	var v4065 int32
	_ = v4065
	var v4070 int32
	_ = v4070
	var v4071 int32
	_ = v4071
	var v4072 int64
	_ = v4072
	var v4073 int64
	_ = v4073
	var v4099 int32
	_ = v4099
	var v4100 int32
	_ = v4100
	var v4161 int64
	_ = v4161
	var v4166 int64
	_ = v4166
	var v4171 int32
	_ = v4171
	var v4173 int32
	_ = v4173
	var v4176 int64
	_ = v4176
	var v4179 int64
	_ = v4179
	var v4182 int64
	_ = v4182
	var v4183 int64
	_ = v4183
	var v4190 int64
	_ = v4190
	var v4195 int64
	_ = v4195
	var v4197 int32
	_ = v4197
	var v4204 int64
	_ = v4204
	var v4207 int64
	_ = v4207
	var v4213 int64
	_ = v4213
	var v4216 int64
	_ = v4216
	var v4223 int64
	_ = v4223
	var v4226 int64
	_ = v4226
	var v4233 int64
	_ = v4233
	var v4236 int64
	_ = v4236
	var v4238 int32
	_ = v4238
	var v4239 int32
	_ = v4239
	var v4240 int64
	_ = v4240
	var v4241 int64
	_ = v4241
	var v4245 int64
	_ = v4245
	var v4249 int64
	_ = v4249
	var v4250 int32
	_ = v4250
	var v4251 int64
	_ = v4251
	var v4252 int32
	_ = v4252
	var v4253 int32
	_ = v4253
	var v4254 int32
	_ = v4254
	var v4265 int32
	_ = v4265
	var v4269 int32
	_ = v4269
	var v4270 int32
	_ = v4270
	var v4273 int32
	_ = v4273
	var v4289 int32
	_ = v4289
	var v4290 int32
	_ = v4290
	var v4296 int32
	_ = v4296
	var v4303 int32
	_ = v4303
	var v4304 int32
	_ = v4304
	var v4307 int32
	_ = v4307
	var v4312 int32
	_ = v4312
	var v4317 int32
	_ = v4317
	var v4327 int32
	_ = v4327
	var v4339 int32
	_ = v4339
	var v4361 int32
	_ = v4361
	var v4364 int32
	_ = v4364
	var v4365 int32
	_ = v4365
	var v4373 int32
	_ = v4373
	var v4374 int32
	_ = v4374
	var v4409 int32
	_ = v4409
	var v4410 int32
	_ = v4410
	var v4415 int64
	_ = v4415
	var v4418 int32
	_ = v4418
	var v4419 int32
	_ = v4419
	var v4431 int32
	_ = v4431
	var v4452 int32
	_ = v4452
	var v4453 int32
	_ = v4453
	var v4456 int32
	_ = v4456
	var v4458 int32
	_ = v4458
	var v4469 int32
	_ = v4469
	var v4491 int32
	_ = v4491
	var v4494 int32
	_ = v4494
	var v4495 int32
	_ = v4495
	var v4503 int32
	_ = v4503
	var v4504 int32
	_ = v4504
	var v4539 int32
	_ = v4539
	var v4540 int32
	_ = v4540
	var v4545 int64
	_ = v4545
	var v4548 int32
	_ = v4548
	var v4549 int32
	_ = v4549
	var v4561 int32
	_ = v4561
	var v4579 int32
	_ = v4579
	var v4580 int32
	_ = v4580
	var v4582 int32
	_ = v4582
	var v4590 int32
	_ = v4590
	var v4593 int32
	_ = v4593
	var v4594 int32
	_ = v4594
	var v4600 int32
	_ = v4600
	var v4601 int32
	_ = v4601
	var v4602 int32
	_ = v4602
	var v4603 int32
	_ = v4603
	var v4605 int32
	_ = v4605
	var v4607 int32
	_ = v4607
	var v4609 int32
	_ = v4609
	var v4611 int32
	_ = v4611
	var v4706 int32
	_ = v4706
	var v4711 int32
	_ = v4711
	var v4718 int32
	_ = v4718
	var v4743 int32
	_ = v4743
	var v4744 int32
	_ = v4744
	var v4745 int32
	_ = v4745
	var v4746 int32
	_ = v4746
	var v4751 int32
	_ = v4751
	var v4804 int64
	_ = v4804
	var v4814 int32
	_ = v4814
	var v4815 int32
	_ = v4815
	var v4816 int64
	_ = v4816
	var v4817 int32
	_ = v4817
	var v4820 int32
	_ = v4820
	var v4821 int32
	_ = v4821
	var v4822 int64
	_ = v4822
	var v4827 int32
	_ = v4827
	var v4828 int64
	_ = v4828
	var v4833 int32
	_ = v4833
	var v4834 int64
	_ = v4834
	var v4839 int32
	_ = v4839
	var v4840 int64
	_ = v4840
	var v4841 int64
	_ = v4841
	var v4845 int32
	_ = v4845
	var v4846 int32
	_ = v4846
	var v4847 int64
	_ = v4847
	var v4855 int32
	_ = v4855
	var v4858 int32
	_ = v4858
	var v4861 int32
	_ = v4861
	var v4862 int32
	_ = v4862
	var v4865 int32
	_ = v4865
	var v4866 int32
	_ = v4866
	var v4867 int32
	_ = v4867
	var v4868 int32
	_ = v4868
	var v4870 int32
	_ = v4870
	var v4871 int32
	_ = v4871
	var v4872 int32
	_ = v4872
	var v4873 int32
	_ = v4873
	var v4874 int32
	_ = v4874
	var v4876 int32
	_ = v4876
	var v4877 int32
	_ = v4877
	var v4878 int32
	_ = v4878
	var v4879 int32
	_ = v4879
	var v4886 int32
	_ = v4886
	var v4890 int32
	_ = v4890
	var v4895 int32
	_ = v4895
	var v4952 int32
	_ = v4952
	var v4968 int32
	_ = v4968
	var v4982 int32
	_ = v4982
	var v5060 int32
	_ = v5060
	var v5138 int32
	_ = v5138
	var v5146 int32
	_ = v5146
	var v5149 int32
	_ = v5149
	var v5150 int32
	_ = v5150
	var v5153 int32
	_ = v5153
	var v5154 int32
	_ = v5154
	var v5156 int32
	_ = v5156
	var v5157 int32
	_ = v5157
	var v5158 int32
	_ = v5158
	var v5159 int32
	_ = v5159
	var v5160 int32
	_ = v5160
	var v5162 int32
	_ = v5162
	var v5163 int32
	_ = v5163
	var v5164 int32
	_ = v5164
	var v5165 int32
	_ = v5165
	var v5166 int32
	_ = v5166
	var v5168 int32
	_ = v5168
	var v5169 int32
	_ = v5169
	var v5170 int32
	_ = v5170
	var v5171 int32
	_ = v5171
	var v5172 int32
	_ = v5172
	var v5173 int32
	_ = v5173
	var v5247 int32
	_ = v5247
	var v5250 int32
	_ = v5250
	var v5253 int32
	_ = v5253
	var v5258 int32
	_ = v5258
	var v5259 int32
	_ = v5259
	var v5261 int32
	_ = v5261
	var v5262 int32
	_ = v5262
	var v5264 int32
	_ = v5264
	var v5269 int32
	_ = v5269
	var v5274 int32
	_ = v5274
	var v5283 int32
	_ = v5283
	var __phi5283 int32
	_ = __phi5283
	var v5285 int32
	_ = v5285
	var __phi5285 int32
	_ = __phi5285
	var v5298 int32
	_ = v5298
	var __phi5298 int32
	_ = __phi5298
	var v5300 int32
	_ = v5300
	var __phi5300 int32
	_ = __phi5300
	var v5368 int32
	_ = v5368
	var v5369 int32
	_ = v5369
	var v5374 int32
	_ = v5374
	var v5376 int32
	_ = v5376
	var v5377 int32
	_ = v5377
	var v5379 int32
	_ = v5379
	var v5385 int32
	_ = v5385
	var v5391 int32
	_ = v5391
	var v5392 int32
	_ = v5392
	var v5402 int32
	_ = v5402
	var __phi5402 int32
	_ = __phi5402
	var v5404 int32
	_ = v5404
	var __phi5404 int32
	_ = __phi5404
	var v5419 int32
	_ = v5419
	var __phi5419 int32
	_ = __phi5419
	var v5420 int32
	_ = v5420
	var __phi5420 int32
	_ = __phi5420
	var v5422 int32
	_ = v5422
	var __phi5422 int32
	_ = __phi5422
	var v5488 int32
	_ = v5488
	var v5490 int32
	_ = v5490
	var v5492 int32
	_ = v5492
	var v5497 int32
	_ = v5497
	var v5506 int32
	_ = v5506
	var v5521 int32
	_ = v5521
	var v5522 int32
	_ = v5522
	var v5525 int32
	_ = v5525
	var v5526 int32
	_ = v5526
	var v5527 int32
	_ = v5527
	var v5529 int32
	_ = v5529
	var v5532 int32
	_ = v5532
	var v5537 int32
	_ = v5537
	var v5541 int32
	_ = v5541
	var v5555 int32
	_ = v5555
	var v5556 int32
	_ = v5556
	var v5566 int32
	_ = v5566
	var v5570 int32
	_ = v5570
	var __phi5570 int32
	_ = __phi5570
	var v5573 int32
	_ = v5573
	var __phi5573 int32
	_ = __phi5573
	var v5574 int32
	_ = v5574
	var __phi5574 int32
	_ = __phi5574
	var v5576 int32
	_ = v5576
	var __phi5576 int32
	_ = __phi5576
	var v5577 int32
	_ = v5577
	var __phi5577 int32
	_ = __phi5577
	var v5586 int32
	_ = v5586
	var v5587 int32
	_ = v5587
	var v5589 int32
	_ = v5589
	var v5591 int32
	_ = v5591
	var v5594 int32
	_ = v5594
	var v5595 int32
	_ = v5595
	var v5596 int32
	_ = v5596
	var v5598 int32
	_ = v5598
	var v5601 int32
	_ = v5601
	var v5606 int32
	_ = v5606
	var v5610 int32
	_ = v5610
	var v5625 int32
	_ = v5625
	var v5639 int32
	_ = v5639
	var v5641 int32
	_ = v5641
	var v5642 int32
	_ = v5642
	var v5643 int32
	_ = v5643
	var v5644 int32
	_ = v5644
	var v5647 int32
	_ = v5647
	var v5648 int32
	_ = v5648
	var v5649 int32
	_ = v5649
	var v5651 int32
	_ = v5651
	var v5652 int32
	_ = v5652
	var v5654 int32
	_ = v5654
	var v5655 int32
	_ = v5655
	var v5660 int32
	_ = v5660
	var v5662 int32
	_ = v5662
	var v5664 int32
	_ = v5664
	var v5666 int32
	_ = v5666
	var v5667 int32
	_ = v5667
	var v5669 int32
	_ = v5669
	var v5670 int32
	_ = v5670
	var v5675 int32
	_ = v5675
	var v5676 int32
	_ = v5676
	var v5677 int32
	_ = v5677
	var v5679 int32
	_ = v5679
	var v5681 int32
	_ = v5681
	var v5684 int32
	_ = v5684
	var v5686 int32
	_ = v5686
	var v5688 int32
	_ = v5688
	var v5690 int32
	_ = v5690
	var v5692 int32
	_ = v5692
	var v5695 int32
	_ = v5695
	var v5696 int32
	_ = v5696
	var v5697 int32
	_ = v5697
	var v5699 int32
	_ = v5699
	var v5704 int32
	_ = v5704
	var v5706 int32
	_ = v5706
	var v5708 int32
	_ = v5708
	var v5710 int32
	_ = v5710
	var v5713 int32
	_ = v5713
	var v5715 int32
	_ = v5715
	var v5716 int32
	_ = v5716
	var v5717 int32
	_ = v5717
	var v5719 int32
	_ = v5719
	var v5722 int32
	_ = v5722
	var v5724 int32
	_ = v5724
	var v5726 int32
	_ = v5726
	var v5728 int32
	_ = v5728
	var v5730 int32
	_ = v5730
	var v5732 int32
	_ = v5732
	var v5735 int32
	_ = v5735
	var v5736 int32
	_ = v5736
	var v5737 int32
	_ = v5737
	var v5739 int32
	_ = v5739
	var v5744 int32
	_ = v5744
	var v5746 int32
	_ = v5746
	var v5748 int32
	_ = v5748
	var v5750 int32
	_ = v5750
	var v5753 int32
	_ = v5753
	var v5755 int32
	_ = v5755
	var v5756 int32
	_ = v5756
	var v5757 int32
	_ = v5757
	var v5759 int32
	_ = v5759
	var v5762 int32
	_ = v5762
	var v5764 int32
	_ = v5764
	var v5766 int32
	_ = v5766
	var v5768 int32
	_ = v5768
	var v5770 int32
	_ = v5770
	var v5773 int32
	_ = v5773
	var v5774 int32
	_ = v5774
	var v5775 int32
	_ = v5775
	var v5777 int32
	_ = v5777
	var v5782 int32
	_ = v5782
	var v5784 int32
	_ = v5784
	var v5786 int32
	_ = v5786
	var v5788 int32
	_ = v5788
	var v5791 int32
	_ = v5791
	var v5793 int32
	_ = v5793
	var v5794 int32
	_ = v5794
	var v5795 int32
	_ = v5795
	var v5797 int32
	_ = v5797
	var v5800 int32
	_ = v5800
	var v5802 int32
	_ = v5802
	var v5804 int32
	_ = v5804
	var v5806 int32
	_ = v5806
	var v5808 int32
	_ = v5808
	var v5810 int32
	_ = v5810
	var v5812 int32
	_ = v5812
	var v5817 int32
	_ = v5817
	var v5835 int32
	_ = v5835
	var v5836 int32
	_ = v5836
	var v5851 int32
	_ = v5851
	var v5938 int32
	_ = v5938
	var v5939 int32
	_ = v5939
	var v5940 int32
	_ = v5940
	var v5941 int32
	_ = v5941
	var v5950 int32
	_ = v5950
	var v5957 int32
	_ = v5957
	var __phi5957 int32
	_ = __phi5957
	var v5958 int32
	_ = v5958
	var __phi5958 int32
	_ = __phi5958
	var v5960 int32
	_ = v5960
	var __phi5960 int32
	_ = __phi5960
	var v5973 int32
	_ = v5973
	var __phi5973 int32
	_ = __phi5973
	var v5975 int32
	_ = v5975
	var __phi5975 int32
	_ = __phi5975
	var v5978 int32
	_ = v5978
	var __phi5978 int32
	_ = __phi5978
	var v6044 int32
	_ = v6044
	var v6047 int32
	_ = v6047
	var v6049 int32
	_ = v6049
	var v6053 int32
	_ = v6053
	var v6068 int32
	_ = v6068
	var v6069 int32
	_ = v6069
	var v6072 int32
	_ = v6072
	var v6073 int32
	_ = v6073
	var v6074 int32
	_ = v6074
	var v6076 int32
	_ = v6076
	var v6079 int32
	_ = v6079
	var v6084 int32
	_ = v6084
	var v6088 int32
	_ = v6088
	var v6102 int32
	_ = v6102
	var v6103 int32
	_ = v6103
	var v6113 int32
	_ = v6113
	var v6117 int32
	_ = v6117
	var __phi6117 int32
	_ = __phi6117
	var v6120 int32
	_ = v6120
	var __phi6120 int32
	_ = __phi6120
	var v6121 int32
	_ = v6121
	var __phi6121 int32
	_ = __phi6121
	var v6123 int32
	_ = v6123
	var __phi6123 int32
	_ = __phi6123
	var v6124 int32
	_ = v6124
	var __phi6124 int32
	_ = __phi6124
	var v6133 int32
	_ = v6133
	var v6134 int32
	_ = v6134
	var v6136 int32
	_ = v6136
	var v6138 int32
	_ = v6138
	var v6141 int32
	_ = v6141
	var v6142 int32
	_ = v6142
	var v6143 int32
	_ = v6143
	var v6145 int32
	_ = v6145
	var v6148 int32
	_ = v6148
	var v6153 int32
	_ = v6153
	var v6157 int32
	_ = v6157
	var v6172 int32
	_ = v6172
	var v6186 int32
	_ = v6186
	var v6188 int32
	_ = v6188
	var v6189 int32
	_ = v6189
	var v6190 int32
	_ = v6190
	var v6191 int32
	_ = v6191
	var v6194 int32
	_ = v6194
	var v6195 int32
	_ = v6195
	var v6196 int32
	_ = v6196
	var v6198 int32
	_ = v6198
	var v6199 int32
	_ = v6199
	var v6201 int32
	_ = v6201
	var v6202 int32
	_ = v6202
	var v6207 int32
	_ = v6207
	var v6209 int32
	_ = v6209
	var v6211 int32
	_ = v6211
	var v6213 int32
	_ = v6213
	var v6214 int32
	_ = v6214
	var v6216 int32
	_ = v6216
	var v6217 int32
	_ = v6217
	var v6222 int32
	_ = v6222
	var v6223 int32
	_ = v6223
	var v6224 int32
	_ = v6224
	var v6226 int32
	_ = v6226
	var v6228 int32
	_ = v6228
	var v6231 int32
	_ = v6231
	var v6233 int32
	_ = v6233
	var v6235 int32
	_ = v6235
	var v6237 int32
	_ = v6237
	var v6239 int32
	_ = v6239
	var v6242 int32
	_ = v6242
	var v6243 int32
	_ = v6243
	var v6244 int32
	_ = v6244
	var v6246 int32
	_ = v6246
	var v6251 int32
	_ = v6251
	var v6253 int32
	_ = v6253
	var v6255 int32
	_ = v6255
	var v6257 int32
	_ = v6257
	var v6260 int32
	_ = v6260
	var v6262 int32
	_ = v6262
	var v6263 int32
	_ = v6263
	var v6264 int32
	_ = v6264
	var v6266 int32
	_ = v6266
	var v6269 int32
	_ = v6269
	var v6271 int32
	_ = v6271
	var v6273 int32
	_ = v6273
	var v6275 int32
	_ = v6275
	var v6277 int32
	_ = v6277
	var v6279 int32
	_ = v6279
	var v6282 int32
	_ = v6282
	var v6283 int32
	_ = v6283
	var v6284 int32
	_ = v6284
	var v6286 int32
	_ = v6286
	var v6291 int32
	_ = v6291
	var v6293 int32
	_ = v6293
	var v6295 int32
	_ = v6295
	var v6297 int32
	_ = v6297
	var v6300 int32
	_ = v6300
	var v6302 int32
	_ = v6302
	var v6303 int32
	_ = v6303
	var v6304 int32
	_ = v6304
	var v6306 int32
	_ = v6306
	var v6309 int32
	_ = v6309
	var v6311 int32
	_ = v6311
	var v6313 int32
	_ = v6313
	var v6315 int32
	_ = v6315
	var v6317 int32
	_ = v6317
	var v6320 int32
	_ = v6320
	var v6321 int32
	_ = v6321
	var v6322 int32
	_ = v6322
	var v6324 int32
	_ = v6324
	var v6329 int32
	_ = v6329
	var v6331 int32
	_ = v6331
	var v6333 int32
	_ = v6333
	var v6335 int32
	_ = v6335
	var v6338 int32
	_ = v6338
	var v6340 int32
	_ = v6340
	var v6341 int32
	_ = v6341
	var v6342 int32
	_ = v6342
	var v6344 int32
	_ = v6344
	var v6347 int32
	_ = v6347
	var v6349 int32
	_ = v6349
	var v6351 int32
	_ = v6351
	var v6353 int32
	_ = v6353
	var v6355 int32
	_ = v6355
	var v6357 int32
	_ = v6357
	var v6359 int32
	_ = v6359
	var v6364 int32
	_ = v6364
	var v6384 int32
	_ = v6384
	var v6404 int32
	_ = v6404
	var v6417 int32
	_ = v6417
	var v6491 int32
	_ = v6491
	var v6493 int32
	_ = v6493
	var v6496 int32
	_ = v6496
	var v6501 int32
	_ = v6501
	var v6512 int32
	_ = v6512
	var v6598 int32
	_ = v6598
	var v6607 int32
	_ = v6607
	var v6700 int32
	_ = v6700
	var v6709 int32
	_ = v6709
	var v6716 int32
	_ = v6716
	v92 = m.G0
	v94 = v92 - int32(2144)
	m.G0 = v94
	v98 = base.I32_div_s(l8, int32(-20))
	v101 = int32(1) << (uint(v98+int32(5)) % 32)
	v102 = *(*int32)(unsafe.Add(mBase, uint32(l13)))
	if l4 == int32(0) {
		goto L3
	} else {
		goto L4
	}
L1:
	;
	m.G0 = v6716 + int32(2144)
	return v6709
L2:
	;
	if v5154 < int32(1) {
		goto L370
	} else {
		goto L371
	}
L3:
	;
	if l3 < l2 {
		v787 = int32(0)
		goto L16
	} else {
		goto L17
	}
L4:
	;
	v105 = int32(-1)
	v106 = l0 + v105
	v107 = int32(1)
	v108 = v107 << (uint(l3) % 32)
	v110 = int32(base.Ui32(v106+v108) >> (uint(l3) % 32))
	v115 = v110 * int32(base.Ui32(l1+v108+v105)>>(uint(l3)%32))
	if v115 < v107 {
		goto L5
	} else {
		goto L6
	}
L5:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l14))) = l3
	v5153 = l0
	v5154 = l1
	v5156 = l3
	v5157 = l4
	v5158 = l5
	v5159 = l6
	v5160 = l7
	v5162 = l9
	v5163 = l10
	v5164 = l11
	v5165 = l12
	v5166 = l13
	v5168 = v94
	v5169 = v101
	v5170 = v102
	v5171 = v106
	v5172 = v108
	v5173 = v110
	goto L2
L6:
	;
	v119 = v115 & int32(7)
	if base.Ui32(v115) < base.Ui32(int32(8)) {
		v258 = int32(0)
		goto L7
	} else {
		goto L8
	}
L7:
	;
	if v119 == int32(0) {
		goto L5
	} else {
		goto L12
	}
L8:
	;
	v134 = l7
	v148 = int32(0)
	goto L9
L9:
	;
	v217 = int64(-72045495131829504)
	*(*int64)(unsafe.Add(mBase, uint32(v134))) = v217
	*(*int64)(unsafe.Add(mBase, uint32(v134+int32(24)))) = v217
	*(*int64)(unsafe.Add(mBase, uint32(v134+int32(16)))) = v217
	v227 = int32(8)
	*(*int64)(unsafe.Add(mBase, uint32(v134+v227))) = v217
	v234 = v148 + v227
	if v115&int32(2147483640) != v234 {
		v134 = v134 + int32(32)
		v148 = v234
		goto L9
	} else {
		goto L11
	}
L10:
	;
	v258 = v234
	goto L7
L11:
	;
	goto L10
L12:
	;
	v340 = l7 + v258<<(uint(int32(2))%32)
	v353 = v119
	goto L13
L13:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v340))) = int32(-16774400)
	v428 = v353 + int32(-1)
	if v428 != 0 {
		v340 = v340 + int32(4)
		v353 = v428
		goto L13
	} else {
		goto L15
	}
L14:
	;
	goto L5
L15:
	;
	goto L14
L16:
	;
	v855 = base.I64_extend_i32_u(v787)
	v856 = int32(4)
	if v855 == int64(0) {
		goto L29
	} else {
		goto L30
	}
L17:
	;
	v523 = int32(-1)
	v524 = l1 + v523
	v526 = l0 + v523
	v528 = int32(1)
	v529 = l3 - l2 + v528
	if l3 != l2 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	if v529&v528 == int32(0) {
		v787 = v682
		goto L16
	} else {
		goto L24
	}
L19:
	;
	v548 = l2
	v561 = v94 + l2<<(uint(int32(2))%32)
	v563 = int32(0)
	v566 = v529 & int32(-2)
	goto L21
L20:
	;
	v667 = l2
	v682 = int32(0)
	goto L18
L21:
	;
	v631 = int32(1)
	v632 = v631 << (uint(v548) % 32)
	v637 = int32(base.Ui32(v526+v632)>>(uint(v548)%32)) * int32(base.Ui32(v524+v632)>>(uint(v548)%32))
	*(*int32)(unsafe.Add(mBase, uint32(v561))) = v637
	v643 = v548 + v631
	v644 = v631 << (uint(v643) % 32)
	v649 = int32(base.Ui32(v526+v644)>>(uint(v643)%32)) * int32(base.Ui32(v524+v644)>>(uint(v643)%32))
	*(*int32)(unsafe.Add(mBase, uint32(v561+int32(4)))) = v649
	v652 = v649 + (v637 + v563)
	v654 = v548 + int32(2)
	v658 = v566 + int32(-2)
	if v658 != 0 {
		v548 = v654
		v561 = v561 + int32(8)
		v563 = v652
		v566 = v658
		goto L21
	} else {
		goto L23
	}
L22:
	;
	v667 = v654
	v682 = v652
	goto L18
L23:
	;
	goto L22
L24:
	;
	v756 = int32(1) << (uint(v667) % 32)
	v761 = int32(base.Ui32(v526+v756)>>(uint(v667)%32)) * int32(base.Ui32(v524+v756)>>(uint(v667)%32))
	*(*int32)(unsafe.Add(mBase, uint32(v94+v667<<(uint(int32(2))%32)))) = v761
	v787 = v761 + v682
	goto L16
L25:
	;
	v5138 = int32(1) << (uint(v4952) % 32)
	v5146 = F_memcpy(m, v4868, v4886, int32(base.Ui32(v5138+v4895)>>(uint(v4952)%32))*int32(base.Ui32(v5138+v4879)>>(uint(v4952)%32))<<(uint(int32(2))%32))
	mBase = m.M
	F_free(m, v4890)
	mBase = m.M
	goto L369
L26:
	;
	v6709 = int32(0)
	v6716 = v5060
	goto L1
L27:
	;
	if v877 == int32(0) {
		v5060 = v94
		goto L26
	} else {
		goto L33
	}
L28:
	;
	goto L27
L29:
	;
	v875 = F_malloc(m, base.I32_wrap_i64(v855)*v856)
	mBase = m.M
	v877 = v875
	goto L28
L30:
	;
	v863 = base.I64_div_u_s(int64(2147418112), v855)
	v864 = int32(0)
	v865 = base.I64_extend_i32_u(v856)
	if base.Ui64(int64(4294967295)) < base.Ui64(v865*v855) {
		v877 = v864
		goto L28
	} else {
		goto L31
	}
L31:
	;
	if base.Ui64(v863) < base.Ui64(v865) {
		v877 = v864
		goto L28
	} else {
		goto L32
	}
L32:
	;
	goto L29
L33:
	;
	v885 = v94 + int32(48) + l2<<(uint(int32(2))%32)
	*(*int32)(unsafe.Add(mBase, uint32(v885))) = v877
	if l3 <= l2 {
		goto L34
	} else {
		goto L35
	}
L34:
	;
	v1328 = l3 - l2
	v1330 = v1328 + int32(1)
	v1333 = base.I64_extend_i32_s(v1330 * int32(_a_F_VP8LResidualImage_0))
	v1334 = int32(4)
	if v1333 == int64(0) {
		goto L47
	} else {
		goto L48
	}
L35:
	;
	v890 = (l3 - l2) & int32(3)
	if v890 == int32(0) {
		v1013 = v877
		v1028 = l2
		goto L36
	} else {
		goto L37
	}
L36:
	;
	if base.Ui32(int32(-4)) < base.Ui32(l2-l3) {
		goto L34
	} else {
		goto L41
	}
L37:
	;
	v894 = l2 << (uint(int32(2)) % 32)
	v909 = v877
	v922 = v894 + (v94 + int32(48)) + int32(4)
	v923 = v94 + v894
	v924 = l2
	v925 = v890
	goto L38
L38:
	;
	v992 = *(*int32)(unsafe.Add(mBase, uint32(v923)))
	v995 = v909 + v992<<(uint(int32(2))%32)
	*(*int32)(unsafe.Add(mBase, uint32(v922))) = v995
	v997 = int32(4)
	v1002 = v924 + int32(1)
	v1004 = v925 + int32(-1)
	if v1004 != 0 {
		v909 = v995
		v922 = v922 + v997
		v923 = v923 + v997
		v924 = v1002
		v925 = v1004
		goto L38
	} else {
		goto L40
	}
L39:
	;
	v1013 = v995
	v1028 = v1002
	goto L36
L40:
	;
	goto L39
L41:
	;
	v1100 = v1028 << (uint(int32(2)) % 32)
	v1112 = v1013
	v1127 = v94 + int32(48)
	v1128 = v94
	v1130 = l3 - v1028
	goto L42
L42:
	;
	v1195 = v1127 + v1100
	v1196 = int32(4)
	v1198 = v1128 + v1100
	v1199 = *(*int32)(unsafe.Add(mBase, uint32(v1198)))
	v1200 = int32(2)
	v1202 = v1112 + v1199<<(uint(v1200)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1195+v1196))) = v1202
	v1204 = int32(8)
	v1208 = *(*int32)(unsafe.Add(mBase, uint32(v1198+v1196)))
	v1211 = v1202 + v1208<<(uint(v1200)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1195+v1204))) = v1211
	v1213 = int32(12)
	v1217 = *(*int32)(unsafe.Add(mBase, uint32(v1198+v1204)))
	v1220 = v1211 + v1217<<(uint(v1200)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1195+v1213))) = v1220
	v1222 = int32(16)
	v1226 = *(*int32)(unsafe.Add(mBase, uint32(v1198+v1213)))
	v1229 = v1220 + v1226<<(uint(v1200)%32)
	*(*int32)(unsafe.Add(mBase, uint32(v1195+v1222))) = v1229
	v1236 = v1130 + int32(-4)
	if v1236 != 0 {
		v1112 = v1229
		v1127 = v1127 + v1222
		v1128 = v1128 + v1222
		v1130 = v1236
		goto L42
	} else {
		goto L44
	}
L43:
	;
	goto L34
L44:
	;
	goto L43
L45:
	;
	v1356 = *(*int32)(unsafe.Add(mBase, uint32(l13)))
	v1357 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l14))) = v1357
	if v1354 == v1357 {
		v4968 = v94
		v4982 = v877
		goto L51
	} else {
		goto L52
	}
L46:
	;
	goto L45
L47:
	;
	v1352 = F_calloc(m, base.I32_wrap_i64(v1333), v1334)
	mBase = m.M
	v1354 = v1352
	goto L46
L48:
	;
	v1341 = base.I64_div_u_s(int64(2147418112), v1333)
	v1342 = int32(0)
	v1343 = base.I64_extend_i32_u(v1334)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1343*v1333) {
		v1354 = v1342
		goto L46
	} else {
		goto L49
	}
L49:
	;
	if base.Ui64(v1341) < base.Ui64(v1343) {
		v1354 = v1342
		goto L46
	} else {
		goto L50
	}
L50:
	;
	goto L47
L51:
	;
	F_free(m, v4982)
	mBase = m.M
	goto L368
L52:
	;
	v1361 = int32(-1)
	v1362 = l0 + v1361
	v1364 = int32(1) << (uint(l2) % 32)
	v1366 = l1 + v1361
	v1368 = int32(base.Ui32(v1364+v1366) >> (uint(l2) % 32))
	if v1368 == int32(0) {
		goto L54
	} else {
		goto L55
	}
L53:
	;
	v4952 = *(*int32)(unsafe.Add(mBase, uint32(l14)))
	if v4952 != 0 {
		goto L25
	} else {
		goto L367
	}
L54:
	;
	v4706 = int32(_a_F_VP8LResidualImage_1)
	v4711 = int32(_a_F_VP8LResidualImage_2)
	v4718 = int32(0)
	v4743 = l3*v4706 - l2*v4706
	v4744 = l3*v4711 - l2*v4711 + v4711
	v4745 = v4718
	v4746 = v4718
	v4751 = v885
	v4804 = int64(9223372036854775807)
	goto L361
L55:
	;
	v1373 = v1354 + v1330*int32(_a_F_VP8LResidualImage_1)
	v1377 = int32(4)
	if l3 < v1377 {
		goto L56
	} else {
		goto L57
	}
L56:
	;
	v1380 = l3
	goto L58
L57:
	;
	v1380 = v1377
	goto L58
L58:
	;
	if l2 < v1380 {
		goto L59
	} else {
		goto L60
	}
L59:
	;
	v1382 = v1380
	goto L61
L60:
	;
	v1382 = l2
	goto L61
L61:
	;
	v1383 = v1382 - l2
	v1385 = v1383 + int32(1)
	v1386 = int32(2)
	if base.Ui32(v1386) < base.Ui32(v1385) {
		goto L62
	} else {
		goto L63
	}
L62:
	;
	v1389 = v1385
	goto L64
L63:
	;
	v1389 = v1386
	goto L64
L64:
	;
	v1391 = l0 << (uint(int32(2)) % 32)
	v1393 = int32(4)
	v1394 = l6 + v1391 + v1393
	v1397 = v1394 + v1391 + v1393
	v1406 = int32(-1)
	v1411 = int32(base.Ui32(v1364+v1362)>>(uint(l2)%32)) + v1406
	v1412 = int32(_a_F_VP8LResidualImage_1)
	v1415 = v1328*v1412 + v1412
	v1416 = int32(0)
	v1473 = v1416
	v1474 = v1416
	v1475 = v1416
	v1476 = v1416
	v1477 = v1416
	v1478 = v1416
	goto L65
L65:
	;
	v1513 = v1476 << (uint(l2) % 32)
	v1515 = base.B2i32(int32(0) < v1513)
	v1516 = v1513 - v1515
	v1517 = int32(2)
	v1518 = v1516 << (uint(v1517) % 32)
	v1519 = l5 + v1518
	v1520 = v1475 << (uint(l2) % 32)
	v1526 = v1519 + (v1520+int32(-1))*l0<<(uint(v1517)%32)
	v1527 = l0 - v1513
	v1528 = base.B2i32(v1364 < v1527)
	if v1364 < v1527 {
		goto L67
	} else {
		goto L68
	}
L66:
	;
	goto L54
L67:
	;
	v1529 = v1364
	goto L69
L68:
	;
	v1529 = v1527
	goto L69
L69:
	;
	v1530 = v1529 + v1515
	v1534 = v1530<<(uint(int32(2))%32) + int32(4)
	v1535 = l1 - v1520
	if v1364 < v1535 {
		goto L72
	} else {
		goto L73
	}
L70:
	;
	v3432 = int32(2147483646)
	v3505 = int32(0)
	v3511 = v1354 + int32(4092)
	v3512 = v1354 + int32(3068)
	v3513 = v1354 + int32(2044)
	v3514 = v1354
	v3515 = v1354 + int32(1020)
	goto L235
L71:
	;
	if v1520 < int32(1) {
		goto L70
	} else {
		goto L232
	}
L72:
	;
	v1537 = v1364
	goto L74
L73:
	;
	v1537 = v1535
	goto L74
L74:
	;
	if v1537 < int32(1) {
		goto L71
	} else {
		goto L75
	}
L75:
	;
	v1540 = v1397 + v1516
	v1541 = v1529 + v1513
	v1542 = v1530 + v1528
	if v1529 < int32(1) {
		goto L76
	} else {
		goto L77
	}
L76:
	;
	if int32(1) < v101 {
		goto L156
	} else {
		goto L157
	}
L77:
	;
	v1554 = l6
	v1573 = int32(0)
	v1614 = v1394
	goto L78
L78:
	;
	if v1520 < int32(1) {
		goto L80
	} else {
		goto L81
	}
L80:
	;
	v1643 = v1354 + v1573<<(uint(int32(12))%32)
	__phi1653 = v1554
	__phi1713 = v1614
	__phi1714 = int32(0)
	v1653 = __phi1653
	v1713 = __phi1713
	v1714 = __phi1714
	goto L82
L81:
	;
	v1640 = F_memcpy(m, v1614+v1518, v1526, v1534)
	mBase = m.M
	goto L80
L82:
	;
	v1737 = v1714 + v1520
	v1739 = int32(2)
	v1741 = v1519 + v1737*l0<<(uint(v1739)%32)
	v1743 = v1737 + int32(1)
	v1748 = F_memcpy(m, v1653+v1518, v1741, (v1530+base.B2i32(v1743 < l1))<<(uint(v1739)%32))
	mBase = m.M
	if v101 < v1739 {
		goto L84
	} else {
		goto L85
	}
L83:
	;
	v2562 = v1573 + int32(1)
	if v2562 != int32(14) {
		v1554 = v1713
		v1573 = v2562
		v1614 = v1653
		goto L78
	} else {
		goto L155
	}
L84:
	;
	v2083 = v94 + int32(96)
	F_GetResidual(m, l0, l1, v1713, v1653, v1397, v1573, v1513, v1541, v1737, v101, l9, l10, v2083)
	mBase = m.M
	v2108 = v2083
	v2109 = v1529
	goto L143
L85:
	;
	if v1737 < int32(1) {
		goto L84
	} else {
		goto L86
	}
L86:
	;
	if l1 <= v1743 {
		goto L84
	} else {
		goto L87
	}
L87:
	;
	if v1542 < int32(3) {
		goto L89
	} else {
		goto L90
	}
L88:
	;
	goto L84
L89:
	;
	goto L88
L90:
	;
	v1768 = *(*int32)(unsafe.Add(mBase, uint32(v1741)+4))
	v1769 = *(*int32)(unsafe.Add(mBase, uint32(v1741)))
	if l10 == int32(0) {
		v1802 = v1768
		v1803 = v1769
		goto L91
	} else {
		goto L92
	}
L91:
	;
	v1813 = l0 << (uint(int32(2)) % 32)
	__phi1817 = v1741 + int32(4)
	__phi1820 = v1802
	__phi1821 = v1803
	__phi1823 = v1540 + int32(1)
	__phi1824 = v1542 + int32(-2)
	v1817 = __phi1817
	v1820 = __phi1820
	v1821 = __phi1821
	v1823 = __phi1823
	v1824 = __phi1824
	goto L93
L92:
	;
	v1772 = int32(8)
	v1773 = int32(base.Ui32(v1768) >> (uint(v1772) % 32))
	v1774 = int32(255)
	v1776 = int32(16711935)
	v1779 = int32(16)
	v1784 = int32(-16711936)
	v1788 = int32(base.Ui32(v1769) >> (uint(v1772) % 32))
	v1802 = (v1773&v1774+v1768&v1776+v1773<<(uint(v1779)%32))&v1776 | v1768&v1784
	v1803 = (v1788&v1774+v1769&v1776+v1788<<(uint(v1779)%32))&v1776 | v1769&v1784
	goto L91
L93:
	;
	v1833 = v1817 + int32(4)
	v1834 = *(*int32)(unsafe.Add(mBase, uint32(v1833)))
	v1836 = *(*int32)(unsafe.Add(mBase, uint32(v1817+v1813)))
	v1838 = *(*int32)(unsafe.Add(mBase, uint32(v1817+(int32(0)-v1813))))
	if l10 == int32(0) {
		v1886 = v1836
		v1888 = v1838
		v1889 = v1834
		goto L95
	} else {
		goto L96
	}
L94:
	;
	goto L89
L95:
	;
	v1890 = int32(24)
	v1891 = int32(base.Ui32(v1820) >> (uint(v1890) % 32))
	v1894 = v1891 - int32(base.Ui32(v1888)>>(uint(v1890)%32))
	v1895 = int32(31)
	v1896 = v1894 >> (uint(v1895) % 32)
	v1898 = v1894 ^ v1896 - v1896
	v1899 = int32(16)
	v1901 = int32(255)
	v1902 = int32(base.Ui32(v1820)>>(uint(v1899)%32)) & v1901
	v1907 = v1902 - int32(base.Ui32(v1888)>>(uint(v1899)%32))&v1901
	v1909 = v1907 >> (uint(v1895) % 32)
	v1911 = v1907 ^ v1909 - v1909
	if base.Ui32(v1911) < base.Ui32(v1898) {
		goto L97
	} else {
		goto L98
	}
L96:
	;
	v1841 = int32(8)
	v1842 = int32(base.Ui32(v1834) >> (uint(v1841) % 32))
	v1843 = int32(255)
	v1845 = int32(16711935)
	v1848 = int32(16)
	v1853 = int32(-16711936)
	v1857 = int32(base.Ui32(v1836) >> (uint(v1841) % 32))
	v1872 = int32(base.Ui32(v1838) >> (uint(v1841) % 32))
	v1886 = (v1857&v1843+v1836&v1845+v1857<<(uint(v1848)%32))&v1845 | v1836&v1853
	v1888 = (v1872&v1843+v1838&v1845+v1872<<(uint(v1848)%32))&v1845 | v1838&v1853
	v1889 = (v1842&v1843+v1834&v1845+v1842<<(uint(v1848)%32))&v1845 | v1834&v1853
	goto L95
L97:
	;
	v1913 = v1898
	goto L99
L98:
	;
	v1913 = v1911
	goto L99
L99:
	;
	v1914 = int32(8)
	v1916 = int32(255)
	v1917 = int32(base.Ui32(v1820)>>(uint(v1914)%32)) & v1916
	v1922 = v1917 - int32(base.Ui32(v1888)>>(uint(v1914)%32))&v1916
	v1923 = int32(31)
	v1924 = v1922 >> (uint(v1923) % 32)
	v1926 = v1922 ^ v1924 - v1924
	v1928 = v1820 & v1916
	v1931 = v1928 - v1888&v1916
	v1933 = v1931 >> (uint(v1923) % 32)
	v1935 = v1931 ^ v1933 - v1933
	if base.Ui32(v1935) < base.Ui32(v1926) {
		goto L100
	} else {
		goto L101
	}
L100:
	;
	v1937 = v1926
	goto L102
L101:
	;
	v1937 = v1935
	goto L102
L102:
	;
	if base.Ui32(v1937) < base.Ui32(v1913) {
		goto L103
	} else {
		goto L104
	}
L103:
	;
	v1939 = v1913
	goto L105
L104:
	;
	v1939 = v1937
	goto L105
L105:
	;
	v1942 = v1891 - int32(base.Ui32(v1886)>>(uint(int32(24))%32))
	v1943 = int32(31)
	v1944 = v1942 >> (uint(v1943) % 32)
	v1946 = v1942 ^ v1944 - v1944
	v1951 = v1902 - int32(base.Ui32(v1886)>>(uint(int32(16))%32))&int32(255)
	v1953 = v1951 >> (uint(v1943) % 32)
	v1955 = v1951 ^ v1953 - v1953
	if base.Ui32(v1955) < base.Ui32(v1946) {
		goto L106
	} else {
		goto L107
	}
L106:
	;
	v1957 = v1946
	goto L108
L107:
	;
	v1957 = v1955
	goto L108
L108:
	;
	v1960 = int32(255)
	v1962 = v1917 - int32(base.Ui32(v1886)>>(uint(int32(8))%32))&v1960
	v1963 = int32(31)
	v1964 = v1962 >> (uint(v1963) % 32)
	v1966 = v1962 ^ v1964 - v1964
	v1969 = v1928 - v1886&v1960
	v1971 = v1969 >> (uint(v1963) % 32)
	v1973 = v1969 ^ v1971 - v1971
	if base.Ui32(v1973) < base.Ui32(v1966) {
		goto L109
	} else {
		goto L110
	}
L109:
	;
	v1975 = v1966
	goto L111
L110:
	;
	v1975 = v1973
	goto L111
L111:
	;
	if base.Ui32(v1975) < base.Ui32(v1957) {
		goto L112
	} else {
		goto L113
	}
L112:
	;
	v1977 = v1957
	goto L114
L113:
	;
	v1977 = v1975
	goto L114
L114:
	;
	if base.Ui32(v1977) < base.Ui32(v1939) {
		goto L115
	} else {
		goto L116
	}
L115:
	;
	v1979 = v1939
	goto L117
L116:
	;
	v1979 = v1977
	goto L117
L117:
	;
	v1982 = v1891 - int32(base.Ui32(v1821)>>(uint(int32(24))%32))
	v1983 = int32(31)
	v1984 = v1982 >> (uint(v1983) % 32)
	v1986 = v1982 ^ v1984 - v1984
	v1991 = v1902 - int32(base.Ui32(v1821)>>(uint(int32(16))%32))&int32(255)
	v1993 = v1991 >> (uint(v1983) % 32)
	v1995 = v1991 ^ v1993 - v1993
	if base.Ui32(v1995) < base.Ui32(v1986) {
		goto L118
	} else {
		goto L119
	}
L118:
	;
	v1997 = v1986
	goto L120
L119:
	;
	v1997 = v1995
	goto L120
L120:
	;
	v2000 = int32(255)
	v2002 = v1917 - int32(base.Ui32(v1821)>>(uint(int32(8))%32))&v2000
	v2003 = int32(31)
	v2004 = v2002 >> (uint(v2003) % 32)
	v2006 = v2002 ^ v2004 - v2004
	v2009 = v1928 - v1821&v2000
	v2011 = v2009 >> (uint(v2003) % 32)
	v2013 = v2009 ^ v2011 - v2011
	if base.Ui32(v2013) < base.Ui32(v2006) {
		goto L121
	} else {
		goto L122
	}
L121:
	;
	v2015 = v2006
	goto L123
L122:
	;
	v2015 = v2013
	goto L123
L123:
	;
	if base.Ui32(v2015) < base.Ui32(v1997) {
		goto L124
	} else {
		goto L125
	}
L124:
	;
	v2017 = v1997
	goto L126
L125:
	;
	v2017 = v2015
	goto L126
L126:
	;
	v2020 = v1891 - int32(base.Ui32(v1889)>>(uint(int32(24))%32))
	v2021 = int32(31)
	v2022 = v2020 >> (uint(v2021) % 32)
	v2024 = v2020 ^ v2022 - v2022
	v2029 = v1902 - int32(base.Ui32(v1889)>>(uint(int32(16))%32))&int32(255)
	v2031 = v2029 >> (uint(v2021) % 32)
	v2033 = v2029 ^ v2031 - v2031
	if base.Ui32(v2033) < base.Ui32(v2024) {
		goto L127
	} else {
		goto L128
	}
L127:
	;
	v2035 = v2024
	goto L129
L128:
	;
	v2035 = v2033
	goto L129
L129:
	;
	v2038 = int32(255)
	v2040 = v1917 - int32(base.Ui32(v1889)>>(uint(int32(8))%32))&v2038
	v2041 = int32(31)
	v2042 = v2040 >> (uint(v2041) % 32)
	v2044 = v2040 ^ v2042 - v2042
	v2047 = v1928 - v1889&v2038
	v2049 = v2047 >> (uint(v2041) % 32)
	v2051 = v2047 ^ v2049 - v2049
	if base.Ui32(v2051) < base.Ui32(v2044) {
		goto L130
	} else {
		goto L131
	}
L130:
	;
	v2053 = v2044
	goto L132
L131:
	;
	v2053 = v2051
	goto L132
L132:
	;
	if base.Ui32(v2053) < base.Ui32(v2035) {
		goto L133
	} else {
		goto L134
	}
L133:
	;
	v2055 = v2035
	goto L135
L134:
	;
	v2055 = v2053
	goto L135
L135:
	;
	if base.Ui32(v2055) < base.Ui32(v2017) {
		goto L136
	} else {
		goto L137
	}
L136:
	;
	v2057 = v2017
	goto L138
L137:
	;
	v2057 = v2055
	goto L138
L138:
	;
	if base.Ui32(v2057) < base.Ui32(v1979) {
		goto L139
	} else {
		goto L140
	}
L139:
	;
	v2059 = v1979
	goto L141
L140:
	;
	v2059 = v2057
	goto L141
L141:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1823))) = uint8(v2059)
	v2064 = v1824 + int32(-1)
	if v2064 != 0 {
		__phi1817 = v1833
		__phi1820 = v1889
		__phi1821 = v1820
		__phi1823 = v1823 + int32(1)
		__phi1824 = v2064
		v1817 = __phi1817
		v1820 = __phi1820
		v1821 = __phi1821
		v1823 = __phi1823
		v1824 = __phi1824
		goto L93
	} else {
		goto L142
	}
L142:
	;
	goto L94
L143:
	;
	v2178 = *(*int32)(unsafe.Add(mBase, uint32(v2108)))
	v2181 = int32(1020)
	v2183 = v1643 + int32(base.Ui32(v2178)>>(uint(int32(22))%32))&v2181
	v2184 = *(*int32)(unsafe.Add(mBase, uint32(v2183)))
	v2185 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2183))) = v2184 + v2185
	v2194 = v1643 + int32(base.Ui32(v2178)>>(uint(int32(14))%32))&v2181 + int32(1024)
	v2195 = *(*int32)(unsafe.Add(mBase, uint32(v2194)))
	*(*int32)(unsafe.Add(mBase, uint32(v2194))) = v2195 + v2185
	v2205 = v1643 + int32(base.Ui32(v2178)>>(uint(int32(6))%32))&v2181 + int32(2048)
	v2206 = *(*int32)(unsafe.Add(mBase, uint32(v2205)))
	*(*int32)(unsafe.Add(mBase, uint32(v2205))) = v2206 + v2185
	v2216 = v1643 + v2178&int32(255)<<(uint(int32(2))%32) + int32(3072)
	v2217 = *(*int32)(unsafe.Add(mBase, uint32(v2216)))
	*(*int32)(unsafe.Add(mBase, uint32(v2216))) = v2217 + v2185
	v2224 = v2109 + int32(-1)
	if v2224 != 0 {
		v2108 = v2108 + int32(4)
		v2109 = v2224
		goto L143
	} else {
		goto L145
	}
L144:
	;
	if v1380 <= l2 {
		goto L146
	} else {
		goto L147
	}
L145:
	;
	goto L144
L146:
	;
	v2559 = v1714 + int32(1)
	if v2559 != v1537 {
		__phi1653 = v1713
		__phi1713 = v1653
		__phi1714 = v2559
		v1653 = __phi1653
		v1713 = __phi1713
		v1714 = __phi1714
		goto L82
	} else {
		goto L154
	}
L147:
	;
	v2253 = int32(1)
	goto L148
L148:
	;
	v2323 = v1354 + (v2253*int32(14)+v1573)<<(uint(int32(12))%32)
	v2348 = v94 + int32(96)
	v2349 = v1529
	goto L150
L149:
	;
	goto L146
L150:
	;
	v2417 = *(*int32)(unsafe.Add(mBase, uint32(v2348)))
	v2420 = int32(1020)
	v2422 = v2323 + int32(base.Ui32(v2417)>>(uint(int32(22))%32))&v2420
	v2423 = *(*int32)(unsafe.Add(mBase, uint32(v2422)))
	v2424 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2422))) = v2423 + v2424
	v2433 = v2323 + int32(base.Ui32(v2417)>>(uint(int32(14))%32))&v2420 + int32(1024)
	v2434 = *(*int32)(unsafe.Add(mBase, uint32(v2433)))
	*(*int32)(unsafe.Add(mBase, uint32(v2433))) = v2434 + v2424
	v2444 = v2323 + int32(base.Ui32(v2417)>>(uint(int32(6))%32))&v2420 + int32(2048)
	v2445 = *(*int32)(unsafe.Add(mBase, uint32(v2444)))
	*(*int32)(unsafe.Add(mBase, uint32(v2444))) = v2445 + v2424
	v2455 = v2323 + v2417&int32(255)<<(uint(int32(2))%32) + int32(3072)
	v2456 = *(*int32)(unsafe.Add(mBase, uint32(v2455)))
	*(*int32)(unsafe.Add(mBase, uint32(v2455))) = v2456 + v2424
	v2463 = v2349 + int32(-1)
	if v2463 != 0 {
		v2348 = v2348 + int32(4)
		v2349 = v2463
		goto L150
	} else {
		goto L152
	}
L151:
	;
	v2465 = v2253 + int32(1)
	if v2465 != v1389 {
		v2253 = v2465
		goto L148
	} else {
		goto L153
	}
L152:
	;
	goto L151
L153:
	;
	goto L149
L154:
	;
	goto L83
L155:
	;
	goto L70
L156:
	;
	v2781 = int32(2)
	v2811 = v1394
	v2812 = l6
	v2815 = int32(0)
	goto L166
L157:
	;
	v2568 = int32(2)
	v2599 = v1394
	v2600 = l6
	v2602 = int32(0)
	goto L158
L158:
	;
	if v1520 < int32(1) {
		goto L160
	} else {
		goto L161
	}
L160:
	;
	__phi2679 = l5 + (v1391*v1520 + v1513<<(uint(v2568)%32) - v1515<<(uint(v2568)%32))
	__phi2692 = v1520
	__phi2694 = v2599
	__phi2695 = v2600
	__phi2696 = v1537
	v2679 = __phi2679
	v2692 = __phi2692
	v2694 = __phi2694
	v2695 = __phi2695
	v2696 = __phi2696
	goto L162
L161:
	;
	v2670 = F_memcpy(m, v2599+v1518, v1526, v1534)
	mBase = m.M
	goto L160
L162:
	;
	v2764 = v2692 + int32(1)
	v2769 = F_memcpy(m, v2695+v1518, v2679, (v1530+base.B2i32(v2764 < l1))<<(uint(int32(2))%32))
	mBase = m.M
	F_GetResidual(m, l0, l1, v2694, v2695, v1397, v2602, v1513, v1541, v2692, v101, l9, l10, v94+int32(96))
	mBase = m.M
	v2775 = v2696 + int32(-1)
	if v2775 != 0 {
		__phi2679 = v2679 + v1391
		__phi2692 = v2764
		__phi2694 = v2695
		__phi2695 = v2694
		__phi2696 = v2775
		v2679 = __phi2679
		v2692 = __phi2692
		v2694 = __phi2694
		v2695 = __phi2695
		v2696 = __phi2696
		goto L162
	} else {
		goto L164
	}
L163:
	;
	v2777 = v2602 + int32(1)
	if v2777 != int32(14) {
		v2599 = v2695
		v2600 = v2694
		v2602 = v2777
		goto L158
	} else {
		goto L165
	}
L164:
	;
	goto L163
L165:
	;
	goto L70
L166:
	;
	if v1520 < int32(1) {
		goto L168
	} else {
		goto L169
	}
L168:
	;
	__phi2892 = l5 + (v1391*v1520 + v1513<<(uint(v2781)%32) - v1515<<(uint(v2781)%32))
	__phi2906 = v2811
	__phi2907 = v2812
	__phi2908 = v1520
	__phi2909 = v1537
	v2892 = __phi2892
	v2906 = __phi2906
	v2907 = __phi2907
	v2908 = __phi2908
	v2909 = __phi2909
	goto L170
L169:
	;
	v2883 = F_memcpy(m, v2811+v1518, v1526, v1534)
	mBase = m.M
	goto L168
L170:
	;
	v2976 = int32(1)
	v2977 = v2908 + v2976
	v2982 = F_memcpy(m, v2907+v1518, v2892, (v1530+base.B2i32(v2977 < l1))<<(uint(int32(2))%32))
	mBase = m.M
	if v2908 < v2976 {
		goto L172
	} else {
		goto L173
	}
L171:
	;
	v3321 = v2815 + int32(1)
	if v3321 != int32(14) {
		v2811 = v2907
		v2812 = v2906
		v2815 = v3321
		goto L166
	} else {
		goto L231
	}
L172:
	;
	F_GetResidual(m, l0, l1, v2906, v2907, v1397, v2815, v1513, v1541, v2908, v101, l9, l10, v94+int32(96))
	mBase = m.M
	v3319 = v2909 + int32(-1)
	if v3319 != 0 {
		__phi2892 = v2892 + v1391
		__phi2906 = v2907
		__phi2907 = v2906
		__phi2908 = v2977
		__phi2909 = v3319
		v2892 = __phi2892
		v2906 = __phi2906
		v2907 = __phi2907
		v2908 = __phi2908
		v2909 = __phi2909
		goto L170
	} else {
		goto L230
	}
L173:
	;
	if l1 <= v2977 {
		goto L172
	} else {
		goto L174
	}
L174:
	;
	if v1542 < int32(3) {
		goto L176
	} else {
		goto L177
	}
L175:
	;
	goto L172
L176:
	;
	goto L175
L177:
	;
	v3000 = *(*int32)(unsafe.Add(mBase, uint32(v2892)+4))
	v3001 = *(*int32)(unsafe.Add(mBase, uint32(v2892)))
	if l10 == int32(0) {
		v3034 = v3000
		v3035 = v3001
		goto L178
	} else {
		goto L179
	}
L178:
	;
	v3045 = l0 << (uint(int32(2)) % 32)
	__phi3049 = v2892 + int32(4)
	__phi3052 = v3034
	__phi3053 = v3035
	__phi3055 = v1540 + int32(1)
	__phi3056 = v1542 + int32(-2)
	v3049 = __phi3049
	v3052 = __phi3052
	v3053 = __phi3053
	v3055 = __phi3055
	v3056 = __phi3056
	goto L180
L179:
	;
	v3004 = int32(8)
	v3005 = int32(base.Ui32(v3000) >> (uint(v3004) % 32))
	v3006 = int32(255)
	v3008 = int32(16711935)
	v3011 = int32(16)
	v3016 = int32(-16711936)
	v3020 = int32(base.Ui32(v3001) >> (uint(v3004) % 32))
	v3034 = (v3005&v3006+v3000&v3008+v3005<<(uint(v3011)%32))&v3008 | v3000&v3016
	v3035 = (v3020&v3006+v3001&v3008+v3020<<(uint(v3011)%32))&v3008 | v3001&v3016
	goto L178
L180:
	;
	v3065 = v3049 + int32(4)
	v3066 = *(*int32)(unsafe.Add(mBase, uint32(v3065)))
	v3068 = *(*int32)(unsafe.Add(mBase, uint32(v3049+v3045)))
	v3070 = *(*int32)(unsafe.Add(mBase, uint32(v3049+(int32(0)-v3045))))
	if l10 == int32(0) {
		v3118 = v3068
		v3120 = v3070
		v3121 = v3066
		goto L182
	} else {
		goto L183
	}
L181:
	;
	goto L176
L182:
	;
	v3122 = int32(24)
	v3123 = int32(base.Ui32(v3052) >> (uint(v3122) % 32))
	v3126 = v3123 - int32(base.Ui32(v3120)>>(uint(v3122)%32))
	v3127 = int32(31)
	v3128 = v3126 >> (uint(v3127) % 32)
	v3130 = v3126 ^ v3128 - v3128
	v3131 = int32(16)
	v3133 = int32(255)
	v3134 = int32(base.Ui32(v3052)>>(uint(v3131)%32)) & v3133
	v3139 = v3134 - int32(base.Ui32(v3120)>>(uint(v3131)%32))&v3133
	v3141 = v3139 >> (uint(v3127) % 32)
	v3143 = v3139 ^ v3141 - v3141
	if base.Ui32(v3143) < base.Ui32(v3130) {
		goto L184
	} else {
		goto L185
	}
L183:
	;
	v3073 = int32(8)
	v3074 = int32(base.Ui32(v3066) >> (uint(v3073) % 32))
	v3075 = int32(255)
	v3077 = int32(16711935)
	v3080 = int32(16)
	v3085 = int32(-16711936)
	v3089 = int32(base.Ui32(v3068) >> (uint(v3073) % 32))
	v3104 = int32(base.Ui32(v3070) >> (uint(v3073) % 32))
	v3118 = (v3089&v3075+v3068&v3077+v3089<<(uint(v3080)%32))&v3077 | v3068&v3085
	v3120 = (v3104&v3075+v3070&v3077+v3104<<(uint(v3080)%32))&v3077 | v3070&v3085
	v3121 = (v3074&v3075+v3066&v3077+v3074<<(uint(v3080)%32))&v3077 | v3066&v3085
	goto L182
L184:
	;
	v3145 = v3130
	goto L186
L185:
	;
	v3145 = v3143
	goto L186
L186:
	;
	v3146 = int32(8)
	v3148 = int32(255)
	v3149 = int32(base.Ui32(v3052)>>(uint(v3146)%32)) & v3148
	v3154 = v3149 - int32(base.Ui32(v3120)>>(uint(v3146)%32))&v3148
	v3155 = int32(31)
	v3156 = v3154 >> (uint(v3155) % 32)
	v3158 = v3154 ^ v3156 - v3156
	v3160 = v3052 & v3148
	v3163 = v3160 - v3120&v3148
	v3165 = v3163 >> (uint(v3155) % 32)
	v3167 = v3163 ^ v3165 - v3165
	if base.Ui32(v3167) < base.Ui32(v3158) {
		goto L187
	} else {
		goto L188
	}
L187:
	;
	v3169 = v3158
	goto L189
L188:
	;
	v3169 = v3167
	goto L189
L189:
	;
	if base.Ui32(v3169) < base.Ui32(v3145) {
		goto L190
	} else {
		goto L191
	}
L190:
	;
	v3171 = v3145
	goto L192
L191:
	;
	v3171 = v3169
	goto L192
L192:
	;
	v3174 = v3123 - int32(base.Ui32(v3118)>>(uint(int32(24))%32))
	v3175 = int32(31)
	v3176 = v3174 >> (uint(v3175) % 32)
	v3178 = v3174 ^ v3176 - v3176
	v3183 = v3134 - int32(base.Ui32(v3118)>>(uint(int32(16))%32))&int32(255)
	v3185 = v3183 >> (uint(v3175) % 32)
	v3187 = v3183 ^ v3185 - v3185
	if base.Ui32(v3187) < base.Ui32(v3178) {
		goto L193
	} else {
		goto L194
	}
L193:
	;
	v3189 = v3178
	goto L195
L194:
	;
	v3189 = v3187
	goto L195
L195:
	;
	v3192 = int32(255)
	v3194 = v3149 - int32(base.Ui32(v3118)>>(uint(int32(8))%32))&v3192
	v3195 = int32(31)
	v3196 = v3194 >> (uint(v3195) % 32)
	v3198 = v3194 ^ v3196 - v3196
	v3201 = v3160 - v3118&v3192
	v3203 = v3201 >> (uint(v3195) % 32)
	v3205 = v3201 ^ v3203 - v3203
	if base.Ui32(v3205) < base.Ui32(v3198) {
		goto L196
	} else {
		goto L197
	}
L196:
	;
	v3207 = v3198
	goto L198
L197:
	;
	v3207 = v3205
	goto L198
L198:
	;
	if base.Ui32(v3207) < base.Ui32(v3189) {
		goto L199
	} else {
		goto L200
	}
L199:
	;
	v3209 = v3189
	goto L201
L200:
	;
	v3209 = v3207
	goto L201
L201:
	;
	if base.Ui32(v3209) < base.Ui32(v3171) {
		goto L202
	} else {
		goto L203
	}
L202:
	;
	v3211 = v3171
	goto L204
L203:
	;
	v3211 = v3209
	goto L204
L204:
	;
	v3214 = v3123 - int32(base.Ui32(v3053)>>(uint(int32(24))%32))
	v3215 = int32(31)
	v3216 = v3214 >> (uint(v3215) % 32)
	v3218 = v3214 ^ v3216 - v3216
	v3223 = v3134 - int32(base.Ui32(v3053)>>(uint(int32(16))%32))&int32(255)
	v3225 = v3223 >> (uint(v3215) % 32)
	v3227 = v3223 ^ v3225 - v3225
	if base.Ui32(v3227) < base.Ui32(v3218) {
		goto L205
	} else {
		goto L206
	}
L205:
	;
	v3229 = v3218
	goto L207
L206:
	;
	v3229 = v3227
	goto L207
L207:
	;
	v3232 = int32(255)
	v3234 = v3149 - int32(base.Ui32(v3053)>>(uint(int32(8))%32))&v3232
	v3235 = int32(31)
	v3236 = v3234 >> (uint(v3235) % 32)
	v3238 = v3234 ^ v3236 - v3236
	v3241 = v3160 - v3053&v3232
	v3243 = v3241 >> (uint(v3235) % 32)
	v3245 = v3241 ^ v3243 - v3243
	if base.Ui32(v3245) < base.Ui32(v3238) {
		goto L208
	} else {
		goto L209
	}
L208:
	;
	v3247 = v3238
	goto L210
L209:
	;
	v3247 = v3245
	goto L210
L210:
	;
	if base.Ui32(v3247) < base.Ui32(v3229) {
		goto L211
	} else {
		goto L212
	}
L211:
	;
	v3249 = v3229
	goto L213
L212:
	;
	v3249 = v3247
	goto L213
L213:
	;
	v3252 = v3123 - int32(base.Ui32(v3121)>>(uint(int32(24))%32))
	v3253 = int32(31)
	v3254 = v3252 >> (uint(v3253) % 32)
	v3256 = v3252 ^ v3254 - v3254
	v3261 = v3134 - int32(base.Ui32(v3121)>>(uint(int32(16))%32))&int32(255)
	v3263 = v3261 >> (uint(v3253) % 32)
	v3265 = v3261 ^ v3263 - v3263
	if base.Ui32(v3265) < base.Ui32(v3256) {
		goto L214
	} else {
		goto L215
	}
L214:
	;
	v3267 = v3256
	goto L216
L215:
	;
	v3267 = v3265
	goto L216
L216:
	;
	v3270 = int32(255)
	v3272 = v3149 - int32(base.Ui32(v3121)>>(uint(int32(8))%32))&v3270
	v3273 = int32(31)
	v3274 = v3272 >> (uint(v3273) % 32)
	v3276 = v3272 ^ v3274 - v3274
	v3279 = v3160 - v3121&v3270
	v3281 = v3279 >> (uint(v3273) % 32)
	v3283 = v3279 ^ v3281 - v3281
	if base.Ui32(v3283) < base.Ui32(v3276) {
		goto L217
	} else {
		goto L218
	}
L217:
	;
	v3285 = v3276
	goto L219
L218:
	;
	v3285 = v3283
	goto L219
L219:
	;
	if base.Ui32(v3285) < base.Ui32(v3267) {
		goto L220
	} else {
		goto L221
	}
L220:
	;
	v3287 = v3267
	goto L222
L221:
	;
	v3287 = v3285
	goto L222
L222:
	;
	if base.Ui32(v3287) < base.Ui32(v3249) {
		goto L223
	} else {
		goto L224
	}
L223:
	;
	v3289 = v3249
	goto L225
L224:
	;
	v3289 = v3287
	goto L225
L225:
	;
	if base.Ui32(v3289) < base.Ui32(v3211) {
		goto L226
	} else {
		goto L227
	}
L226:
	;
	v3291 = v3211
	goto L228
L227:
	;
	v3291 = v3289
	goto L228
L228:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3055))) = uint8(v3291)
	v3296 = v3056 + int32(-1)
	if v3296 != 0 {
		__phi3049 = v3065
		__phi3052 = v3121
		__phi3053 = v3052
		__phi3055 = v3055 + int32(1)
		__phi3056 = v3296
		v3049 = __phi3049
		v3052 = __phi3052
		v3053 = __phi3053
		v3055 = __phi3055
		v3056 = __phi3056
		goto L180
	} else {
		goto L229
	}
L229:
	;
	goto L181
L230:
	;
	goto L171
L231:
	;
	goto L70
L232:
	;
	v3327 = F_memcpy(m, v1394+v1518, v1526, v1534)
	mBase = m.M
	v3328 = F_memcpy(m, v3327, v1526, v1534)
	mBase = m.M
	v3329 = F_memcpy(m, v3328, v1526, v1534)
	mBase = m.M
	v3330 = F_memcpy(m, v3329, v1526, v1534)
	mBase = m.M
	v3331 = F_memcpy(m, v3330, v1526, v1534)
	mBase = m.M
	v3332 = F_memcpy(m, v3331, v1526, v1534)
	mBase = m.M
	v3333 = F_memcpy(m, v3332, v1526, v1534)
	mBase = m.M
	v3334 = F_memcpy(m, v3333, v1526, v1534)
	mBase = m.M
	v3335 = F_memcpy(m, v3334, v1526, v1534)
	mBase = m.M
	v3336 = F_memcpy(m, v3335, v1526, v1534)
	mBase = m.M
	v3337 = F_memcpy(m, v3336, v1526, v1534)
	mBase = m.M
	v3338 = F_memcpy(m, v3337, v1526, v1534)
	mBase = m.M
	v3339 = F_memcpy(m, v3338, v1526, v1534)
	mBase = m.M
	v3340 = F_memcpy(m, v3339, v1526, v1534)
	mBase = m.M
	goto L70
L233:
	;
	v4605 = v4602<<(uint(v1328)%32) + v4601
	v4607 = v4600 + v4603<<(uint(v1328)%32)
	if v4607 != 0 {
		goto L356
	} else {
		goto L357
	}
L234:
	;
	v4456 = int32(0)
	v4458 = v3578 + int32(_a_F_VP8LResidualImage_1)
	if base.Ui32(v4458) < base.Ui32(int32(33)) {
		goto L336
	} else {
		goto L337
	}
L235:
	;
	v3529 = int32(1)
	v3530 = v3505 + l2
	v3533 = int32(base.Ui32(v3529<<(uint(v3530)%32)+v1362) >> (uint(v3530) % 32))
	v3534 = int32(base.Ui32(v1475) >> (uint(v3505) % 32))
	v3538 = *(*int32)(unsafe.Add(mBase, uint32(v885+v3505<<(uint(int32(2))%32))))
	v3539 = int32(255)
	v3541 = int32(base.Ui32(v1476) >> (uint(v3505) % 32))
	if v3541 < v3529 {
		v3554 = v3539
		goto L237
	} else {
		goto L238
	}
L236:
	;
	v4327 = int32(0)
	if base.Ui32(v1415) < base.Ui32(int32(33)) {
		goto L319
	} else {
		goto L320
	}
L237:
	;
	v3557 = v1373 + v3505<<(uint(int32(12))%32)
	if v3534 < int32(1) {
		v3570 = v3539
		goto L239
	} else {
		goto L240
	}
L238:
	;
	v3545 = int32(2)
	v3553 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3538+v3533*v3534<<(uint(v3545)%32)+v3541<<(uint(v3545)%32)+int32(-3)))))
	v3554 = v3553
	goto L237
L239:
	;
	v3578 = v3505 * int32(_a_F_VP8LResidualImage_1)
	v3579 = v1354 + v3578
	v3581 = int32(0)
	v3591 = v3514
	v3603 = v1354
	v3607 = v3581
	v3608 = v3511
	v3609 = v3512
	v3610 = v3513
	v3611 = v3515
	v3652 = v3581
	v3664 = int64(9223372036854775807)
	goto L241
L240:
	;
	v3563 = int32(2)
	v3569 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3538+v3533*(v3534+int32(-1))<<(uint(v3563)%32)+v3541<<(uint(v3563)%32))+1)))
	v3570 = v3569
	goto L239
L241:
	;
	v3676 = v3579 + v3607<<(uint(int32(12))%32)
	v3677 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3676))))
	v3703 = int32(4)
	v3704 = v3611
	v3764 = v3677 << (uint(int64(23)) % 64)
	v3765 = int64(788529152)
	goto L243
L242:
	;
	v4269 = m.G58
	v4270 = *(*int32)(unsafe.Add(mBase, uint32(v4269)))
	m.T0[v4270].(func(*base.Module, int32, int32, int32))(m, v4252, v3557, int32(1024))
	mBase = m.M
	v4273 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v3538+v3533*v3534<<(uint(v4273)%32)+v3541<<(uint(v4273)%32)))) = v4253<<(uint(int32(8))%32) | int32(-16777216)
	v4289 = v1373 + v1330<<(uint(int32(12))%32) + v3505*int32(56) + v4253<<(uint(v4273)%32)
	v4290 = *(*int32)(unsafe.Add(mBase, uint32(v4289)))
	*(*int32)(unsafe.Add(mBase, uint32(v4289))) = v4290 + int32(1)
	if v3505 == v1328 {
		goto L307
	} else {
		goto L308
	}
L243:
	;
	v3775 = *(*int32)(unsafe.Add(mBase, uint32(v3704)))
	v3777 = *(*int32)(unsafe.Add(mBase, uint32(v3591+v3703)))
	v3780 = v3765 * base.I64_extend_i32_u(v3775+v3777)
	if v3780 < int64(0) {
		goto L245
	} else {
		goto L246
	}
L244:
	;
	v3806 = m.G64
	v3807 = *(*int32)(unsafe.Add(mBase, uint32(v3806)))
	v3808 = m.T0[v3807].(func(*base.Module, int32, int32) int64)(m, v3676, v3557)
	mBase = m.M
	v3809 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3676)+1024)))
	v3835 = int32(1028)
	v3836 = v3610
	v3897 = int64(788529152)
	v3898 = v3809 << (uint(int64(23)) % 64)
	goto L252
L245:
	;
	v3783 = int64(-50)
	goto L247
L246:
	;
	v3783 = int64(50)
	goto L247
L247:
	;
	v3786 = base.I64_div_s(v3783+v3780, int64(100))
	v3787 = v3786 + v3764
	if v3765 < int64(0) {
		goto L248
	} else {
		goto L249
	}
L248:
	;
	v3794 = int64(-5)
	goto L250
L249:
	;
	v3794 = int64(5)
	goto L250
L250:
	;
	v3799 = base.I64_div_s(v3794+v3765*int64(6), int64(10))
	v3801 = v3703 + int32(4)
	if v3801 != int32(64) {
		v3703 = v3801
		v3704 = v3704 + int32(-4)
		v3764 = v3787
		v3765 = v3799
		goto L243
	} else {
		goto L251
	}
L251:
	;
	goto L244
L252:
	;
	v3907 = *(*int32)(unsafe.Add(mBase, uint32(v3836)))
	v3909 = *(*int32)(unsafe.Add(mBase, uint32(v3591+v3835)))
	v3912 = v3897 * base.I64_extend_i32_u(v3907+v3909)
	if v3912 < int64(0) {
		goto L254
	} else {
		goto L255
	}
L253:
	;
	v3938 = m.G64
	v3939 = *(*int32)(unsafe.Add(mBase, uint32(v3938)))
	v3940 = m.T0[v3939].(func(*base.Module, int32, int32) int64)(m, v3676+int32(1024), v3557+int32(1024))
	mBase = m.M
	v3941 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3676)+2048)))
	v3967 = int32(2052)
	v3968 = v3609
	v4029 = int64(788529152)
	v4032 = v3941 << (uint(int64(23)) % 64)
	goto L261
L254:
	;
	v3915 = int64(-50)
	goto L256
L255:
	;
	v3915 = int64(50)
	goto L256
L256:
	;
	v3918 = base.I64_div_s(v3915+v3912, int64(100))
	v3919 = v3918 + v3898
	if v3897 < int64(0) {
		goto L257
	} else {
		goto L258
	}
L257:
	;
	v3926 = int64(-5)
	goto L259
L258:
	;
	v3926 = int64(5)
	goto L259
L259:
	;
	v3931 = base.I64_div_s(v3926+v3897*int64(6), int64(10))
	v3933 = v3835 + int32(4)
	if v3933 != int32(1088) {
		v3835 = v3933
		v3836 = v3836 + int32(-4)
		v3897 = v3931
		v3898 = v3919
		goto L252
	} else {
		goto L260
	}
L260:
	;
	goto L253
L261:
	;
	v4039 = *(*int32)(unsafe.Add(mBase, uint32(v3968)))
	v4041 = *(*int32)(unsafe.Add(mBase, uint32(v3591+v3967)))
	v4044 = v4029 * base.I64_extend_i32_u(v4039+v4041)
	if v4044 < int64(0) {
		goto L263
	} else {
		goto L264
	}
L262:
	;
	v4070 = m.G64
	v4071 = *(*int32)(unsafe.Add(mBase, uint32(v4070)))
	v4072 = m.T0[v4071].(func(*base.Module, int32, int32) int64)(m, v3676+int32(2048), v3557+int32(2048))
	mBase = m.M
	v4073 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v3676)+3072)))
	v4099 = int32(3076)
	v4100 = v3608
	v4161 = int64(788529152)
	v4166 = v4073 << (uint(int64(23)) % 64)
	goto L270
L263:
	;
	v4047 = int64(-50)
	goto L265
L264:
	;
	v4047 = int64(50)
	goto L265
L265:
	;
	v4050 = base.I64_div_s(v4047+v4044, int64(100))
	v4051 = v4050 + v4032
	if v4029 < int64(0) {
		goto L266
	} else {
		goto L267
	}
L266:
	;
	v4058 = int64(-5)
	goto L268
L267:
	;
	v4058 = int64(5)
	goto L268
L268:
	;
	v4063 = base.I64_div_s(v4058+v4029*int64(6), int64(10))
	v4065 = v3967 + int32(4)
	if v4065 != int32(2112) {
		v3967 = v4065
		v3968 = v3968 + int32(-4)
		v4029 = v4063
		v4032 = v4051
		goto L261
	} else {
		goto L269
	}
L269:
	;
	goto L262
L270:
	;
	v4171 = *(*int32)(unsafe.Add(mBase, uint32(v4100)))
	v4173 = *(*int32)(unsafe.Add(mBase, uint32(v3591+v4099)))
	v4176 = v4161 * base.I64_extend_i32_u(v4171+v4173)
	if v4176 < int64(0) {
		goto L272
	} else {
		goto L273
	}
L271:
	;
	if v3787 < int64(0) {
		goto L279
	} else {
		goto L280
	}
L272:
	;
	v4179 = int64(-50)
	goto L274
L273:
	;
	v4179 = int64(50)
	goto L274
L274:
	;
	v4182 = base.I64_div_s(v4179+v4176, int64(100))
	v4183 = v4182 + v4166
	if v4161 < int64(0) {
		goto L275
	} else {
		goto L276
	}
L275:
	;
	v4190 = int64(-5)
	goto L277
L276:
	;
	v4190 = int64(5)
	goto L277
L277:
	;
	v4195 = base.I64_div_s(v4190+v4161*int64(6), int64(10))
	v4197 = v4099 + int32(4)
	if v4197 != int32(3136) {
		v4099 = v4197
		v4100 = v4100 + int32(-4)
		v4161 = v4195
		v4166 = v4183
		goto L270
	} else {
		goto L278
	}
L278:
	;
	goto L271
L279:
	;
	v4204 = int64(-5)
	goto L281
L280:
	;
	v4204 = int64(5)
	goto L281
L281:
	;
	v4207 = base.I64_div_s(v4204+v3787, int64(-10))
	if v3919 < int64(0) {
		goto L282
	} else {
		goto L283
	}
L282:
	;
	v4213 = int64(-5)
	goto L284
L283:
	;
	v4213 = int64(5)
	goto L284
L284:
	;
	v4216 = base.I64_div_s(v4213+v3919, int64(-10))
	if v4051 < int64(0) {
		goto L285
	} else {
		goto L286
	}
L285:
	;
	v4223 = int64(-5)
	goto L287
L286:
	;
	v4223 = int64(5)
	goto L287
L287:
	;
	v4226 = base.I64_div_s(v4223+v4051, int64(-10))
	if v4183 < int64(0) {
		goto L288
	} else {
		goto L289
	}
L288:
	;
	v4233 = int64(-5)
	goto L290
L289:
	;
	v4233 = int64(5)
	goto L290
L290:
	;
	v4236 = base.I64_div_s(v4233+v4183, int64(-10))
	v4238 = m.G64
	v4239 = *(*int32)(unsafe.Add(mBase, uint32(v4238)))
	v4240 = m.T0[v4239].(func(*base.Module, int32, int32) int64)(m, v3676+int32(3072), v3557+int32(3072))
	mBase = m.M
	v4241 = v3808 + v4207 + v4216 + v3940 + v4226 + v4072 + v4236 + v4240
	if v3607 == v3554 {
		goto L291
	} else {
		goto L292
	}
L291:
	;
	v4245 = v4241 + int64(-125829120)
	goto L293
L292:
	;
	v4245 = v4241
	goto L293
L293:
	;
	if v3607 == v3570 {
		goto L294
	} else {
		goto L295
	}
L294:
	;
	v4249 = v4245 + int64(-125829120)
	goto L296
L295:
	;
	v4249 = v4245
	goto L296
L296:
	;
	v4250 = base.B2i32(v4249 < v3664)
	if v4249 < v3664 {
		goto L297
	} else {
		goto L298
	}
L297:
	;
	v4251 = v4249
	goto L299
L298:
	;
	v4251 = v3664
	goto L299
L299:
	;
	if v4249 < v3664 {
		goto L300
	} else {
		goto L301
	}
L300:
	;
	v4252 = v3676
	goto L302
L301:
	;
	v4252 = v3603
	goto L302
L302:
	;
	if v4249 < v3664 {
		goto L303
	} else {
		goto L304
	}
L303:
	;
	v4253 = v3607
	goto L305
L304:
	;
	v4253 = v3652
	goto L305
L305:
	;
	v4254 = int32(_a_F_VP8LResidualImage_3)
	v4265 = v3607 + int32(1)
	if v4265 != int32(14) {
		v3591 = v3591 + v4254
		v3603 = v4252
		v3607 = v4265
		v3608 = v3608 + v4254
		v3609 = v3609 + v4254
		v3610 = v3610 + v4254
		v3611 = v3611 + v4254
		v3652 = v4253
		v3664 = v4251
		goto L241
	} else {
		goto L306
	}
L306:
	;
	goto L242
L307:
	;
	goto L236
L308:
	;
	v4296 = v3505 + int32(1)
	if base.Ui32(v4296) <= base.Ui32(v1383) {
		goto L309
	} else {
		goto L310
	}
L309:
	;
	if v1476 == v1411 {
		goto L312
	} else {
		goto L313
	}
L310:
	;
	if base.Ui32(v1328) < base.Ui32(v4296) {
		goto L309
	} else {
		goto L311
	}
L311:
	;
	v4303 = m.G58
	v4304 = *(*int32)(unsafe.Add(mBase, uint32(v4303)))
	m.T0[v4304].(func(*base.Module, int32, int32, int32))(m, v3579, v1354+v4296*int32(_a_F_VP8LResidualImage_1), int32(_a_F_VP8LResidualImage_4))
	mBase = m.M
	goto L309
L312:
	;
	if v1475 == v1368+v1406 {
		goto L315
	} else {
		goto L316
	}
L313:
	;
	v4307 = int32(-1)
	if v4307<<(uint(v4296)%32)|(v3432-v1473) != v4307 {
		goto L234
	} else {
		goto L314
	}
L314:
	;
	goto L312
L315:
	;
	v4317 = int32(_a_F_VP8LResidualImage_1)
	v3505 = v4296
	v3511 = v3511 + v4317
	v3512 = v3512 + v4317
	v3513 = v3513 + v4317
	v3514 = v3514 + v4317
	v3515 = v3515 + v4317
	goto L235
L316:
	;
	v4312 = int32(-1)
	if v4312<<(uint(v4296)%32)|(v3432-v1474) != v4312 {
		goto L234
	} else {
		goto L317
	}
L317:
	;
	goto L315
L318:
	;
	v4452 = base.B2i32(v1476 == v1411)
	if v1476 == v1411 {
		goto L332
	} else {
		goto L333
	}
L319:
	;
	if v1415 == int32(0) {
		goto L321
	} else {
		goto L322
	}
L320:
	;
	base.MemoryFill(m, v1354, v4327, v1415)
	goto L318
L321:
	;
	goto L318
L322:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1354))) = uint8(v4327)
	v4339 = v1354 + v1415
	*(*uint8)(unsafe.Add(mBase, uint32(v4339+int32(-1)))) = uint8(v4327)
	if base.Ui32(v1415) < base.Ui32(int32(3)) {
		goto L321
	} else {
		goto L323
	}
L323:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1354)+2)) = uint8(v4327)
	*(*uint8)(unsafe.Add(mBase, uint32(v1354)+1)) = uint8(v4327)
	*(*uint8)(unsafe.Add(mBase, uint32(v4339+int32(-3)))) = uint8(v4327)
	*(*uint8)(unsafe.Add(mBase, uint32(v4339+int32(-2)))) = uint8(v4327)
	if base.Ui32(v1415) < base.Ui32(int32(7)) {
		goto L321
	} else {
		goto L324
	}
L324:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1354)+3)) = uint8(v4327)
	*(*uint8)(unsafe.Add(mBase, uint32(v4339+int32(-4)))) = uint8(v4327)
	if base.Ui32(v1415) < base.Ui32(int32(9)) {
		goto L321
	} else {
		goto L325
	}
L325:
	;
	v4361 = int32(0)
	v4364 = (v4361 - v1354) & int32(3)
	v4365 = v1354 + v4364
	*(*int32)(unsafe.Add(mBase, uint32(v4365))) = v4361
	v4373 = (v1415 - v4364) & int32(60)
	v4374 = v4365 + v4373
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-4)))) = v4361
	if base.Ui32(v4373) < base.Ui32(int32(9)) {
		goto L321
	} else {
		goto L326
	}
L326:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4365)+8)) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4365)+4)) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-8)))) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-12)))) = v4361
	if base.Ui32(v4373) < base.Ui32(int32(25)) {
		goto L321
	} else {
		goto L327
	}
L327:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4365)+24)) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4365)+20)) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4365)+16)) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4365)+12)) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-16)))) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-20)))) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-24)))) = v4361
	*(*int32)(unsafe.Add(mBase, uint32(v4374+int32(-28)))) = v4361
	v4409 = v4365&int32(4) | int32(24)
	v4410 = v4373 - v4409
	if base.Ui32(v4410) < base.Ui32(int32(32)) {
		goto L321
	} else {
		goto L328
	}
L328:
	;
	v4415 = base.I64_extend_i32_u(v4361) * int64(4294967297)
	v4418 = v4410
	v4419 = v4365 + v4409
	goto L329
L329:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4419)+24)) = v4415
	*(*int64)(unsafe.Add(mBase, uint32(v4419)+16)) = v4415
	*(*int64)(unsafe.Add(mBase, uint32(v4419)+8)) = v4415
	*(*int64)(unsafe.Add(mBase, uint32(v4419))) = v4415
	v4431 = v4418 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v4431) {
		v4418 = v4431
		v4419 = v4419 + int32(32)
		goto L329
	} else {
		goto L331
	}
L330:
	;
	goto L321
L331:
	;
	goto L330
L332:
	;
	v4453 = int32(0)
	goto L334
L333:
	;
	v4453 = v1478 + int32(1)
	goto L334
L334:
	;
	v4600 = int32(0)
	v4601 = v4327
	v4602 = v1477 + v4452
	v4603 = v4453
	goto L233
L335:
	;
	v4579 = int32(base.Ui32(v1474) >> (uint(v3505) % 32))
	v4580 = int32(base.Ui32(v1473) >> (uint(v3505) % 32))
	v4582 = v4580 & int32(1)
	if v1476 != v1411 {
		goto L351
	} else {
		goto L352
	}
L336:
	;
	if v4458 == int32(0) {
		goto L338
	} else {
		goto L339
	}
L337:
	;
	base.MemoryFill(m, v1354, v4456, v4458)
	goto L335
L338:
	;
	goto L335
L339:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1354))) = uint8(v4456)
	v4469 = v1354 + v4458
	*(*uint8)(unsafe.Add(mBase, uint32(v4469+int32(-1)))) = uint8(v4456)
	if base.Ui32(v4458) < base.Ui32(int32(3)) {
		goto L338
	} else {
		goto L340
	}
L340:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1354)+2)) = uint8(v4456)
	*(*uint8)(unsafe.Add(mBase, uint32(v1354)+1)) = uint8(v4456)
	*(*uint8)(unsafe.Add(mBase, uint32(v4469+int32(-3)))) = uint8(v4456)
	*(*uint8)(unsafe.Add(mBase, uint32(v4469+int32(-2)))) = uint8(v4456)
	if base.Ui32(v4458) < base.Ui32(int32(7)) {
		goto L338
	} else {
		goto L341
	}
L341:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1354)+3)) = uint8(v4456)
	*(*uint8)(unsafe.Add(mBase, uint32(v4469+int32(-4)))) = uint8(v4456)
	if base.Ui32(v4458) < base.Ui32(int32(9)) {
		goto L338
	} else {
		goto L342
	}
L342:
	;
	v4491 = int32(0)
	v4494 = (v4491 - v1354) & int32(3)
	v4495 = v1354 + v4494
	*(*int32)(unsafe.Add(mBase, uint32(v4495))) = v4491
	v4503 = (v4458 - v4494) & int32(60)
	v4504 = v4495 + v4503
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-4)))) = v4491
	if base.Ui32(v4503) < base.Ui32(int32(9)) {
		goto L338
	} else {
		goto L343
	}
L343:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4495)+8)) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4495)+4)) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-8)))) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-12)))) = v4491
	if base.Ui32(v4503) < base.Ui32(int32(25)) {
		goto L338
	} else {
		goto L344
	}
L344:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4495)+24)) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4495)+20)) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4495)+16)) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4495)+12)) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-16)))) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-20)))) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-24)))) = v4491
	*(*int32)(unsafe.Add(mBase, uint32(v4504+int32(-28)))) = v4491
	v4539 = v4495&int32(4) | int32(24)
	v4540 = v4503 - v4539
	if base.Ui32(v4540) < base.Ui32(int32(32)) {
		goto L338
	} else {
		goto L345
	}
L345:
	;
	v4545 = base.I64_extend_i32_u(v4491) * int64(4294967297)
	v4548 = v4540
	v4549 = v4495 + v4539
	goto L346
L346:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4549)+24)) = v4545
	*(*int64)(unsafe.Add(mBase, uint32(v4549)+16)) = v4545
	*(*int64)(unsafe.Add(mBase, uint32(v4549)+8)) = v4545
	*(*int64)(unsafe.Add(mBase, uint32(v4549))) = v4545
	v4561 = v4548 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v4561) {
		v4548 = v4561
		v4549 = v4549 + int32(32)
		goto L346
	} else {
		goto L348
	}
L347:
	;
	goto L338
L348:
	;
	goto L347
L349:
	;
	v4600 = v4594 << (uint(v3505) % 32)
	v4601 = v4593 << (uint(v3505) % 32)
	v4602 = v1477
	v4603 = v1478
	goto L233
L350:
	;
	v4593 = v4579 + int32(1)
	v4594 = v4590
	goto L349
L351:
	;
	if v4582 != 0 {
		goto L354
	} else {
		goto L355
	}
L352:
	;
	if v4582 == int32(0) {
		v4590 = v4580
		goto L350
	} else {
		goto L353
	}
L353:
	;
	goto L351
L354:
	;
	v4590 = v4580 + int32(-1)
	goto L350
L355:
	;
	v4593 = v4579
	v4594 = v4580 | int32(1)
	goto L349
L356:
	;
	if base.Ui32(v4605) < base.Ui32(v1368) {
		v1473 = v4600
		v1474 = v4601
		v1475 = v4605
		v1476 = v4607
		v1477 = v4602
		v1478 = v4603
		goto L65
	} else {
		goto L360
	}
L357:
	;
	v4609 = base.I32_div_u_s(v4605*l12, v1368)
	v4611 = F_WebPReportProgress(m, l11, v4609+v1356, l13)
	mBase = m.M
	if v4611 != 0 {
		goto L356
	} else {
		goto L358
	}
L358:
	;
	F_free(m, v1354)
	mBase = m.M
	goto L359
L359:
	;
	v4861 = l0
	v4862 = l1
	v4865 = l4
	v4866 = l5
	v4867 = l6
	v4868 = l7
	v4870 = l9
	v4871 = l10
	v4872 = l11
	v4873 = l12
	v4874 = l13
	v4876 = v94
	v4877 = v101
	v4878 = v102
	v4879 = v1362
	v4886 = int32(0)
	v4890 = v877
	v4895 = v1366
	goto L53
L360:
	;
	goto L66
L361:
	;
	v4814 = m.G65
	v4815 = *(*int32)(unsafe.Add(mBase, uint32(v4814)))
	v4816 = m.T0[v4815].(func(*base.Module, int32, int32) int64)(m, v1354+v4744, int32(14))
	mBase = m.M
	v4817 = v1354 + v4743
	v4820 = int32(256)
	v4821 = *(*int32)(unsafe.Add(mBase, uint32(v4814)))
	v4822 = m.T0[v4821].(func(*base.Module, int32, int32) int64)(m, v4817+int32(_a_F_VP8LResidualImage_1), v4820)
	mBase = m.M
	v4827 = *(*int32)(unsafe.Add(mBase, uint32(v4814)))
	v4828 = m.T0[v4827].(func(*base.Module, int32, int32) int64)(m, v4817+int32(_a_F_VP8LResidualImage_5), v4820)
	mBase = m.M
	v4833 = *(*int32)(unsafe.Add(mBase, uint32(v4814)))
	v4834 = m.T0[v4833].(func(*base.Module, int32, int32) int64)(m, v4817+int32(_a_F_VP8LResidualImage_6), v4820)
	mBase = m.M
	v4839 = *(*int32)(unsafe.Add(mBase, uint32(v4814)))
	v4840 = m.T0[v4839].(func(*base.Module, int32, int32) int64)(m, v4817+int32(_a_F_VP8LResidualImage_7), v4820)
	mBase = m.M
	v4841 = v4816 + v4822 + v4828 + v4834 + v4840
	if v4804 <= v4841 {
		v4846 = v4746
		v4847 = v4804
		goto L363
	} else {
		goto L364
	}
L362:
	;
	F_free(m, v1354)
	mBase = m.M
	goto L366
L363:
	;
	v4855 = v4745 + int32(1)
	if base.Ui32(v4855) <= base.Ui32(v1328) {
		v4743 = v4743 + int32(_a_F_VP8LResidualImage_3)
		v4744 = v4744 + int32(56)
		v4745 = v4855
		v4746 = v4846
		v4751 = v4751 + int32(4)
		v4804 = v4847
		goto L361
	} else {
		goto L365
	}
L364:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l14))) = l2 + v4745
	v4845 = *(*int32)(unsafe.Add(mBase, uint32(v4751)))
	v4846 = v4845
	v4847 = v4841
	goto L363
L365:
	;
	goto L362
L366:
	;
	v4858 = *(*int32)(unsafe.Add(mBase, uint32(l14)))
	F_VP8LOptimizeSampling(m, v4846, l0, l1, v4858, int32(9), l14)
	mBase = m.M
	v4861 = l0
	v4862 = l1
	v4865 = l4
	v4866 = l5
	v4867 = l6
	v4868 = l7
	v4870 = l9
	v4871 = l10
	v4872 = l11
	v4873 = l12
	v4874 = l13
	v4876 = v94
	v4877 = v101
	v4878 = v102
	v4879 = v1362
	v4886 = v4846
	v4890 = v877
	v4895 = v1366
	goto L53
L367:
	;
	v4968 = v4876
	v4982 = v4890
	goto L51
L368:
	;
	v5060 = v4968
	goto L26
L369:
	;
	v5149 = *(*int32)(unsafe.Add(mBase, uint32(l14)))
	v5150 = int32(1) << (uint(v5149) % 32)
	v5153 = v4861
	v5154 = v4862
	v5156 = v5149
	v5157 = v4865
	v5158 = v4866
	v5159 = v4867
	v5160 = v4868
	v5162 = v4870
	v5163 = v4871
	v5164 = v4872
	v5165 = v4873
	v5166 = v4874
	v5168 = v4876
	v5169 = v4877
	v5170 = v4878
	v5171 = v4879
	v5172 = v5150
	v5173 = int32(base.Ui32(v5150+v4879) >> (uint(v5149) % 32))
	goto L2
L370:
	;
	v6700 = F_WebPReportProgress(m, v5164, v5170+v5165, v5166)
	mBase = m.M
	v6709 = v6700
	v6716 = v5168
	goto L1
L371:
	;
	v5247 = v5153 << (uint(int32(2)) % 32)
	v5250 = v5159 + v5247 + int32(4)
	if v5157 == int32(0) {
		goto L372
	} else {
		goto L373
	}
L372:
	;
	v5391 = v5250 + v5247 + int32(4)
	v5392 = v5391 + v5153
	if v5153 < int32(1) {
		goto L378
	} else {
		goto L379
	}
L373:
	;
	v5253 = int32(1)
	v5258 = F_memcpy(m, v5159, v5158, (base.B2i32(v5154 != v5253)+v5153)<<(uint(int32(2))%32))
	mBase = m.M
	v5259 = int32(0)
	v5261 = m.G66
	v5262 = *(*int32)(unsafe.Add(mBase, uint32(v5261)))
	m.T0[v5262].(func(*base.Module, int32, int32, int32, int32))(m, v5258, v5259, v5253, v5158)
	mBase = m.M
	v5264 = int32(4)
	v5269 = *(*int32)(unsafe.Add(mBase, uint32(v5261)+4))
	m.T0[v5269].(func(*base.Module, int32, int32, int32, int32))(m, v5258+v5264, v5259, v5171, v5158+v5264)
	mBase = m.M
	if v5154 == v5253 {
		goto L370
	} else {
		goto L374
	}
L374:
	;
	v5274 = v5153 << (uint(int32(2)) % 32)
	__phi5283 = v5159
	__phi5285 = v5158 + v5274
	__phi5298 = int32(1)
	__phi5300 = v5250
	v5283 = __phi5283
	v5285 = __phi5285
	v5298 = __phi5298
	v5300 = __phi5300
	goto L375
L375:
	;
	v5368 = int32(1)
	v5369 = v5298 + v5368
	v5374 = F_memcpy(m, v5300, v5285, (base.B2i32(v5369 < v5154)+v5153)<<(uint(int32(2))%32))
	mBase = m.M
	v5376 = m.G66
	v5377 = *(*int32)(unsafe.Add(mBase, uint32(v5376)+8))
	m.T0[v5377].(func(*base.Module, int32, int32, int32, int32))(m, v5374, v5283, v5368, v5285)
	mBase = m.M
	v5379 = int32(4)
	v5385 = *(*int32)(unsafe.Add(mBase, uint32(v5376)+44))
	m.T0[v5385].(func(*base.Module, int32, int32, int32, int32))(m, v5374+v5379, v5283+v5379, v5171, v5285+v5379)
	mBase = m.M
	if v5154 != v5369 {
		__phi5283 = v5374
		__phi5285 = v5285 + v5274
		__phi5298 = v5369
		__phi5300 = v5283
		v5283 = __phi5283
		v5285 = __phi5285
		v5298 = __phi5298
		v5300 = __phi5300
		goto L375
	} else {
		goto L377
	}
L377:
	;
	goto L370
L378:
	;
	if v5169 < int32(2) {
		goto L449
	} else {
		goto L450
	}
L379:
	;
	__phi5402 = v5159
	__phi5404 = v5392
	__phi5419 = v5250
	__phi5420 = int32(0)
	__phi5422 = v5391
	v5402 = __phi5402
	v5404 = __phi5404
	v5419 = __phi5419
	v5420 = __phi5420
	v5422 = __phi5422
	goto L380
L380:
	;
	v5488 = int32(2)
	v5490 = v5158 + v5420*v5153<<(uint(v5488)%32)
	v5492 = v5420 + int32(1)
	v5497 = F_memcpy(m, v5402, v5490, (base.B2i32(v5492 < v5154)+v5153)<<(uint(v5488)%32))
	mBase = m.M
	if v5488 <= v5169 {
		goto L383
	} else {
		goto L384
	}
L382:
	;
	v5851 = int32(0)
	goto L442
L383:
	;
	if v5154 <= v5420+int32(2) {
		goto L385
	} else {
		goto L386
	}
L384:
	;
	v5835 = v5422
	v5836 = v5404
	goto L382
L385:
	;
	v5835 = v5404
	v5836 = v5422
	goto L382
L386:
	;
	v5506 = v5158 + v5492*v5153<<(uint(int32(2))%32)
	if v5153 < int32(3) {
		goto L388
	} else {
		goto L389
	}
L387:
	;
	goto L385
L388:
	;
	goto L387
L389:
	;
	v5521 = *(*int32)(unsafe.Add(mBase, uint32(v5506)+4))
	v5522 = *(*int32)(unsafe.Add(mBase, uint32(v5506)))
	if v5163 == int32(0) {
		v5555 = v5521
		v5556 = v5522
		goto L390
	} else {
		goto L391
	}
L390:
	;
	v5566 = v5153 << (uint(int32(2)) % 32)
	__phi5570 = v5506 + int32(4)
	__phi5573 = v5555
	__phi5574 = v5556
	__phi5576 = v5422 + int32(1)
	__phi5577 = v5153 + int32(-2)
	v5570 = __phi5570
	v5573 = __phi5573
	v5574 = __phi5574
	v5576 = __phi5576
	v5577 = __phi5577
	goto L392
L391:
	;
	v5525 = int32(8)
	v5526 = int32(base.Ui32(v5521) >> (uint(v5525) % 32))
	v5527 = int32(255)
	v5529 = int32(16711935)
	v5532 = int32(16)
	v5537 = int32(-16711936)
	v5541 = int32(base.Ui32(v5522) >> (uint(v5525) % 32))
	v5555 = (v5526&v5527+v5521&v5529+v5526<<(uint(v5532)%32))&v5529 | v5521&v5537
	v5556 = (v5541&v5527+v5522&v5529+v5541<<(uint(v5532)%32))&v5529 | v5522&v5537
	goto L390
L392:
	;
	v5586 = v5570 + int32(4)
	v5587 = *(*int32)(unsafe.Add(mBase, uint32(v5586)))
	v5589 = *(*int32)(unsafe.Add(mBase, uint32(v5570+v5566)))
	v5591 = *(*int32)(unsafe.Add(mBase, uint32(v5570+(int32(0)-v5566))))
	if v5163 == int32(0) {
		v5639 = v5589
		v5641 = v5591
		v5642 = v5587
		goto L394
	} else {
		goto L395
	}
L393:
	;
	goto L388
L394:
	;
	v5643 = int32(24)
	v5644 = int32(base.Ui32(v5573) >> (uint(v5643) % 32))
	v5647 = v5644 - int32(base.Ui32(v5641)>>(uint(v5643)%32))
	v5648 = int32(31)
	v5649 = v5647 >> (uint(v5648) % 32)
	v5651 = v5647 ^ v5649 - v5649
	v5652 = int32(16)
	v5654 = int32(255)
	v5655 = int32(base.Ui32(v5573)>>(uint(v5652)%32)) & v5654
	v5660 = v5655 - int32(base.Ui32(v5641)>>(uint(v5652)%32))&v5654
	v5662 = v5660 >> (uint(v5648) % 32)
	v5664 = v5660 ^ v5662 - v5662
	if base.Ui32(v5664) < base.Ui32(v5651) {
		goto L396
	} else {
		goto L397
	}
L395:
	;
	v5594 = int32(8)
	v5595 = int32(base.Ui32(v5587) >> (uint(v5594) % 32))
	v5596 = int32(255)
	v5598 = int32(16711935)
	v5601 = int32(16)
	v5606 = int32(-16711936)
	v5610 = int32(base.Ui32(v5589) >> (uint(v5594) % 32))
	v5625 = int32(base.Ui32(v5591) >> (uint(v5594) % 32))
	v5639 = (v5610&v5596+v5589&v5598+v5610<<(uint(v5601)%32))&v5598 | v5589&v5606
	v5641 = (v5625&v5596+v5591&v5598+v5625<<(uint(v5601)%32))&v5598 | v5591&v5606
	v5642 = (v5595&v5596+v5587&v5598+v5595<<(uint(v5601)%32))&v5598 | v5587&v5606
	goto L394
L396:
	;
	v5666 = v5651
	goto L398
L397:
	;
	v5666 = v5664
	goto L398
L398:
	;
	v5667 = int32(8)
	v5669 = int32(255)
	v5670 = int32(base.Ui32(v5573)>>(uint(v5667)%32)) & v5669
	v5675 = v5670 - int32(base.Ui32(v5641)>>(uint(v5667)%32))&v5669
	v5676 = int32(31)
	v5677 = v5675 >> (uint(v5676) % 32)
	v5679 = v5675 ^ v5677 - v5677
	v5681 = v5573 & v5669
	v5684 = v5681 - v5641&v5669
	v5686 = v5684 >> (uint(v5676) % 32)
	v5688 = v5684 ^ v5686 - v5686
	if base.Ui32(v5688) < base.Ui32(v5679) {
		goto L399
	} else {
		goto L400
	}
L399:
	;
	v5690 = v5679
	goto L401
L400:
	;
	v5690 = v5688
	goto L401
L401:
	;
	if base.Ui32(v5690) < base.Ui32(v5666) {
		goto L402
	} else {
		goto L403
	}
L402:
	;
	v5692 = v5666
	goto L404
L403:
	;
	v5692 = v5690
	goto L404
L404:
	;
	v5695 = v5644 - int32(base.Ui32(v5639)>>(uint(int32(24))%32))
	v5696 = int32(31)
	v5697 = v5695 >> (uint(v5696) % 32)
	v5699 = v5695 ^ v5697 - v5697
	v5704 = v5655 - int32(base.Ui32(v5639)>>(uint(int32(16))%32))&int32(255)
	v5706 = v5704 >> (uint(v5696) % 32)
	v5708 = v5704 ^ v5706 - v5706
	if base.Ui32(v5708) < base.Ui32(v5699) {
		goto L405
	} else {
		goto L406
	}
L405:
	;
	v5710 = v5699
	goto L407
L406:
	;
	v5710 = v5708
	goto L407
L407:
	;
	v5713 = int32(255)
	v5715 = v5670 - int32(base.Ui32(v5639)>>(uint(int32(8))%32))&v5713
	v5716 = int32(31)
	v5717 = v5715 >> (uint(v5716) % 32)
	v5719 = v5715 ^ v5717 - v5717
	v5722 = v5681 - v5639&v5713
	v5724 = v5722 >> (uint(v5716) % 32)
	v5726 = v5722 ^ v5724 - v5724
	if base.Ui32(v5726) < base.Ui32(v5719) {
		goto L408
	} else {
		goto L409
	}
L408:
	;
	v5728 = v5719
	goto L410
L409:
	;
	v5728 = v5726
	goto L410
L410:
	;
	if base.Ui32(v5728) < base.Ui32(v5710) {
		goto L411
	} else {
		goto L412
	}
L411:
	;
	v5730 = v5710
	goto L413
L412:
	;
	v5730 = v5728
	goto L413
L413:
	;
	if base.Ui32(v5730) < base.Ui32(v5692) {
		goto L414
	} else {
		goto L415
	}
L414:
	;
	v5732 = v5692
	goto L416
L415:
	;
	v5732 = v5730
	goto L416
L416:
	;
	v5735 = v5644 - int32(base.Ui32(v5574)>>(uint(int32(24))%32))
	v5736 = int32(31)
	v5737 = v5735 >> (uint(v5736) % 32)
	v5739 = v5735 ^ v5737 - v5737
	v5744 = v5655 - int32(base.Ui32(v5574)>>(uint(int32(16))%32))&int32(255)
	v5746 = v5744 >> (uint(v5736) % 32)
	v5748 = v5744 ^ v5746 - v5746
	if base.Ui32(v5748) < base.Ui32(v5739) {
		goto L417
	} else {
		goto L418
	}
L417:
	;
	v5750 = v5739
	goto L419
L418:
	;
	v5750 = v5748
	goto L419
L419:
	;
	v5753 = int32(255)
	v5755 = v5670 - int32(base.Ui32(v5574)>>(uint(int32(8))%32))&v5753
	v5756 = int32(31)
	v5757 = v5755 >> (uint(v5756) % 32)
	v5759 = v5755 ^ v5757 - v5757
	v5762 = v5681 - v5574&v5753
	v5764 = v5762 >> (uint(v5756) % 32)
	v5766 = v5762 ^ v5764 - v5764
	if base.Ui32(v5766) < base.Ui32(v5759) {
		goto L420
	} else {
		goto L421
	}
L420:
	;
	v5768 = v5759
	goto L422
L421:
	;
	v5768 = v5766
	goto L422
L422:
	;
	if base.Ui32(v5768) < base.Ui32(v5750) {
		goto L423
	} else {
		goto L424
	}
L423:
	;
	v5770 = v5750
	goto L425
L424:
	;
	v5770 = v5768
	goto L425
L425:
	;
	v5773 = v5644 - int32(base.Ui32(v5642)>>(uint(int32(24))%32))
	v5774 = int32(31)
	v5775 = v5773 >> (uint(v5774) % 32)
	v5777 = v5773 ^ v5775 - v5775
	v5782 = v5655 - int32(base.Ui32(v5642)>>(uint(int32(16))%32))&int32(255)
	v5784 = v5782 >> (uint(v5774) % 32)
	v5786 = v5782 ^ v5784 - v5784
	if base.Ui32(v5786) < base.Ui32(v5777) {
		goto L426
	} else {
		goto L427
	}
L426:
	;
	v5788 = v5777
	goto L428
L427:
	;
	v5788 = v5786
	goto L428
L428:
	;
	v5791 = int32(255)
	v5793 = v5670 - int32(base.Ui32(v5642)>>(uint(int32(8))%32))&v5791
	v5794 = int32(31)
	v5795 = v5793 >> (uint(v5794) % 32)
	v5797 = v5793 ^ v5795 - v5795
	v5800 = v5681 - v5642&v5791
	v5802 = v5800 >> (uint(v5794) % 32)
	v5804 = v5800 ^ v5802 - v5802
	if base.Ui32(v5804) < base.Ui32(v5797) {
		goto L429
	} else {
		goto L430
	}
L429:
	;
	v5806 = v5797
	goto L431
L430:
	;
	v5806 = v5804
	goto L431
L431:
	;
	if base.Ui32(v5806) < base.Ui32(v5788) {
		goto L432
	} else {
		goto L433
	}
L432:
	;
	v5808 = v5788
	goto L434
L433:
	;
	v5808 = v5806
	goto L434
L434:
	;
	if base.Ui32(v5808) < base.Ui32(v5770) {
		goto L435
	} else {
		goto L436
	}
L435:
	;
	v5810 = v5770
	goto L437
L436:
	;
	v5810 = v5808
	goto L437
L437:
	;
	if base.Ui32(v5810) < base.Ui32(v5732) {
		goto L438
	} else {
		goto L439
	}
L438:
	;
	v5812 = v5732
	goto L440
L439:
	;
	v5812 = v5810
	goto L440
L440:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v5576))) = uint8(v5812)
	v5817 = v5577 + int32(-1)
	if v5817 != 0 {
		__phi5570 = v5586
		__phi5573 = v5642
		__phi5574 = v5573
		__phi5576 = v5576 + int32(1)
		__phi5577 = v5817
		v5570 = __phi5570
		v5573 = __phi5573
		v5574 = __phi5574
		v5576 = __phi5576
		v5577 = __phi5577
		goto L392
	} else {
		goto L441
	}
L441:
	;
	goto L393
L442:
	;
	v5938 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5160+int32(base.Ui32(v5420)>>(uint(v5156)%32))*v5173<<(uint(int32(2))%32)+v5851>>(uint(v5156)%32)<<(uint(int32(2))%32))+1)))
	v5939 = v5851 + v5172
	v5940 = base.B2i32(v5939 < v5153)
	if v5939 < v5153 {
		goto L444
	} else {
		goto L445
	}
L443:
	;
	if v5492 != v5154 {
		__phi5402 = v5419
		__phi5404 = v5836
		__phi5419 = v5497
		__phi5420 = v5492
		__phi5422 = v5835
		v5402 = __phi5402
		v5404 = __phi5404
		v5419 = __phi5419
		v5420 = __phi5420
		v5422 = __phi5422
		goto L380
	} else {
		goto L448
	}
L444:
	;
	v5941 = v5939
	goto L446
L445:
	;
	v5941 = v5153
	goto L446
L446:
	;
	F_GetResidual(m, v5153, v5154, v5419, v5497, v5835, v5938, v5851, v5941, v5420, v5169, v5162, v5163, v5490+v5851<<(uint(int32(2))%32))
	mBase = m.M
	if v5939 < v5153 {
		v5851 = v5941
		goto L442
	} else {
		goto L447
	}
L447:
	;
	goto L443
L448:
	;
	goto L370
L449:
	;
	v6384 = int32(1)
	if v5154 == v6384 {
		v6512 = int32(0)
		goto L511
	} else {
		goto L512
	}
L450:
	;
	v5950 = v5153 << (uint(int32(2)) % 32)
	__phi5957 = v5158
	__phi5958 = v5159
	__phi5960 = v5392
	__phi5973 = int32(0)
	__phi5975 = v5250
	__phi5978 = v5391
	v5957 = __phi5957
	v5958 = __phi5958
	v5960 = __phi5960
	v5973 = __phi5973
	v5975 = __phi5975
	v5978 = __phi5978
	goto L451
L451:
	;
	v6044 = v5973 + int32(1)
	v6047 = int32(2)
	v6049 = F_memcpy(m, v5958, v5957, (base.B2i32(v6044 < v5154)+v5153)<<(uint(v6047)%32))
	mBase = m.M
	if v5154 <= v5973+v6047 {
		goto L453
	} else {
		goto L454
	}
L453:
	;
	if v5154 != v6044 {
		__phi5957 = v5957 + v5950
		__phi5958 = v5975
		__phi5960 = v5978
		__phi5973 = v6044
		__phi5975 = v6049
		__phi5978 = v5960
		v5957 = __phi5957
		v5958 = __phi5958
		v5960 = __phi5960
		v5973 = __phi5973
		v5975 = __phi5975
		v5978 = __phi5978
		goto L451
	} else {
		goto L510
	}
L454:
	;
	v6053 = v5957 + v5950
	if v5153 < int32(3) {
		goto L456
	} else {
		goto L457
	}
L455:
	;
	goto L453
L456:
	;
	goto L455
L457:
	;
	v6068 = *(*int32)(unsafe.Add(mBase, uint32(v6053)+4))
	v6069 = *(*int32)(unsafe.Add(mBase, uint32(v6053)))
	if v5163 == int32(0) {
		v6102 = v6068
		v6103 = v6069
		goto L458
	} else {
		goto L459
	}
L458:
	;
	v6113 = v5153 << (uint(int32(2)) % 32)
	__phi6117 = v6053 + int32(4)
	__phi6120 = v6102
	__phi6121 = v6103
	__phi6123 = v5978 + int32(1)
	__phi6124 = v5153 + int32(-2)
	v6117 = __phi6117
	v6120 = __phi6120
	v6121 = __phi6121
	v6123 = __phi6123
	v6124 = __phi6124
	goto L460
L459:
	;
	v6072 = int32(8)
	v6073 = int32(base.Ui32(v6068) >> (uint(v6072) % 32))
	v6074 = int32(255)
	v6076 = int32(16711935)
	v6079 = int32(16)
	v6084 = int32(-16711936)
	v6088 = int32(base.Ui32(v6069) >> (uint(v6072) % 32))
	v6102 = (v6073&v6074+v6068&v6076+v6073<<(uint(v6079)%32))&v6076 | v6068&v6084
	v6103 = (v6088&v6074+v6069&v6076+v6088<<(uint(v6079)%32))&v6076 | v6069&v6084
	goto L458
L460:
	;
	v6133 = v6117 + int32(4)
	v6134 = *(*int32)(unsafe.Add(mBase, uint32(v6133)))
	v6136 = *(*int32)(unsafe.Add(mBase, uint32(v6117+v6113)))
	v6138 = *(*int32)(unsafe.Add(mBase, uint32(v6117+(int32(0)-v6113))))
	if v5163 == int32(0) {
		v6186 = v6136
		v6188 = v6138
		v6189 = v6134
		goto L462
	} else {
		goto L463
	}
L461:
	;
	goto L456
L462:
	;
	v6190 = int32(24)
	v6191 = int32(base.Ui32(v6120) >> (uint(v6190) % 32))
	v6194 = v6191 - int32(base.Ui32(v6188)>>(uint(v6190)%32))
	v6195 = int32(31)
	v6196 = v6194 >> (uint(v6195) % 32)
	v6198 = v6194 ^ v6196 - v6196
	v6199 = int32(16)
	v6201 = int32(255)
	v6202 = int32(base.Ui32(v6120)>>(uint(v6199)%32)) & v6201
	v6207 = v6202 - int32(base.Ui32(v6188)>>(uint(v6199)%32))&v6201
	v6209 = v6207 >> (uint(v6195) % 32)
	v6211 = v6207 ^ v6209 - v6209
	if base.Ui32(v6211) < base.Ui32(v6198) {
		goto L464
	} else {
		goto L465
	}
L463:
	;
	v6141 = int32(8)
	v6142 = int32(base.Ui32(v6134) >> (uint(v6141) % 32))
	v6143 = int32(255)
	v6145 = int32(16711935)
	v6148 = int32(16)
	v6153 = int32(-16711936)
	v6157 = int32(base.Ui32(v6136) >> (uint(v6141) % 32))
	v6172 = int32(base.Ui32(v6138) >> (uint(v6141) % 32))
	v6186 = (v6157&v6143+v6136&v6145+v6157<<(uint(v6148)%32))&v6145 | v6136&v6153
	v6188 = (v6172&v6143+v6138&v6145+v6172<<(uint(v6148)%32))&v6145 | v6138&v6153
	v6189 = (v6142&v6143+v6134&v6145+v6142<<(uint(v6148)%32))&v6145 | v6134&v6153
	goto L462
L464:
	;
	v6213 = v6198
	goto L466
L465:
	;
	v6213 = v6211
	goto L466
L466:
	;
	v6214 = int32(8)
	v6216 = int32(255)
	v6217 = int32(base.Ui32(v6120)>>(uint(v6214)%32)) & v6216
	v6222 = v6217 - int32(base.Ui32(v6188)>>(uint(v6214)%32))&v6216
	v6223 = int32(31)
	v6224 = v6222 >> (uint(v6223) % 32)
	v6226 = v6222 ^ v6224 - v6224
	v6228 = v6120 & v6216
	v6231 = v6228 - v6188&v6216
	v6233 = v6231 >> (uint(v6223) % 32)
	v6235 = v6231 ^ v6233 - v6233
	if base.Ui32(v6235) < base.Ui32(v6226) {
		goto L467
	} else {
		goto L468
	}
L467:
	;
	v6237 = v6226
	goto L469
L468:
	;
	v6237 = v6235
	goto L469
L469:
	;
	if base.Ui32(v6237) < base.Ui32(v6213) {
		goto L470
	} else {
		goto L471
	}
L470:
	;
	v6239 = v6213
	goto L472
L471:
	;
	v6239 = v6237
	goto L472
L472:
	;
	v6242 = v6191 - int32(base.Ui32(v6186)>>(uint(int32(24))%32))
	v6243 = int32(31)
	v6244 = v6242 >> (uint(v6243) % 32)
	v6246 = v6242 ^ v6244 - v6244
	v6251 = v6202 - int32(base.Ui32(v6186)>>(uint(int32(16))%32))&int32(255)
	v6253 = v6251 >> (uint(v6243) % 32)
	v6255 = v6251 ^ v6253 - v6253
	if base.Ui32(v6255) < base.Ui32(v6246) {
		goto L473
	} else {
		goto L474
	}
L473:
	;
	v6257 = v6246
	goto L475
L474:
	;
	v6257 = v6255
	goto L475
L475:
	;
	v6260 = int32(255)
	v6262 = v6217 - int32(base.Ui32(v6186)>>(uint(int32(8))%32))&v6260
	v6263 = int32(31)
	v6264 = v6262 >> (uint(v6263) % 32)
	v6266 = v6262 ^ v6264 - v6264
	v6269 = v6228 - v6186&v6260
	v6271 = v6269 >> (uint(v6263) % 32)
	v6273 = v6269 ^ v6271 - v6271
	if base.Ui32(v6273) < base.Ui32(v6266) {
		goto L476
	} else {
		goto L477
	}
L476:
	;
	v6275 = v6266
	goto L478
L477:
	;
	v6275 = v6273
	goto L478
L478:
	;
	if base.Ui32(v6275) < base.Ui32(v6257) {
		goto L479
	} else {
		goto L480
	}
L479:
	;
	v6277 = v6257
	goto L481
L480:
	;
	v6277 = v6275
	goto L481
L481:
	;
	if base.Ui32(v6277) < base.Ui32(v6239) {
		goto L482
	} else {
		goto L483
	}
L482:
	;
	v6279 = v6239
	goto L484
L483:
	;
	v6279 = v6277
	goto L484
L484:
	;
	v6282 = v6191 - int32(base.Ui32(v6121)>>(uint(int32(24))%32))
	v6283 = int32(31)
	v6284 = v6282 >> (uint(v6283) % 32)
	v6286 = v6282 ^ v6284 - v6284
	v6291 = v6202 - int32(base.Ui32(v6121)>>(uint(int32(16))%32))&int32(255)
	v6293 = v6291 >> (uint(v6283) % 32)
	v6295 = v6291 ^ v6293 - v6293
	if base.Ui32(v6295) < base.Ui32(v6286) {
		goto L485
	} else {
		goto L486
	}
L485:
	;
	v6297 = v6286
	goto L487
L486:
	;
	v6297 = v6295
	goto L487
L487:
	;
	v6300 = int32(255)
	v6302 = v6217 - int32(base.Ui32(v6121)>>(uint(int32(8))%32))&v6300
	v6303 = int32(31)
	v6304 = v6302 >> (uint(v6303) % 32)
	v6306 = v6302 ^ v6304 - v6304
	v6309 = v6228 - v6121&v6300
	v6311 = v6309 >> (uint(v6303) % 32)
	v6313 = v6309 ^ v6311 - v6311
	if base.Ui32(v6313) < base.Ui32(v6306) {
		goto L488
	} else {
		goto L489
	}
L488:
	;
	v6315 = v6306
	goto L490
L489:
	;
	v6315 = v6313
	goto L490
L490:
	;
	if base.Ui32(v6315) < base.Ui32(v6297) {
		goto L491
	} else {
		goto L492
	}
L491:
	;
	v6317 = v6297
	goto L493
L492:
	;
	v6317 = v6315
	goto L493
L493:
	;
	v6320 = v6191 - int32(base.Ui32(v6189)>>(uint(int32(24))%32))
	v6321 = int32(31)
	v6322 = v6320 >> (uint(v6321) % 32)
	v6324 = v6320 ^ v6322 - v6322
	v6329 = v6202 - int32(base.Ui32(v6189)>>(uint(int32(16))%32))&int32(255)
	v6331 = v6329 >> (uint(v6321) % 32)
	v6333 = v6329 ^ v6331 - v6331
	if base.Ui32(v6333) < base.Ui32(v6324) {
		goto L494
	} else {
		goto L495
	}
L494:
	;
	v6335 = v6324
	goto L496
L495:
	;
	v6335 = v6333
	goto L496
L496:
	;
	v6338 = int32(255)
	v6340 = v6217 - int32(base.Ui32(v6189)>>(uint(int32(8))%32))&v6338
	v6341 = int32(31)
	v6342 = v6340 >> (uint(v6341) % 32)
	v6344 = v6340 ^ v6342 - v6342
	v6347 = v6228 - v6189&v6338
	v6349 = v6347 >> (uint(v6341) % 32)
	v6351 = v6347 ^ v6349 - v6349
	if base.Ui32(v6351) < base.Ui32(v6344) {
		goto L497
	} else {
		goto L498
	}
L497:
	;
	v6353 = v6344
	goto L499
L498:
	;
	v6353 = v6351
	goto L499
L499:
	;
	if base.Ui32(v6353) < base.Ui32(v6335) {
		goto L500
	} else {
		goto L501
	}
L500:
	;
	v6355 = v6335
	goto L502
L501:
	;
	v6355 = v6353
	goto L502
L502:
	;
	if base.Ui32(v6355) < base.Ui32(v6317) {
		goto L503
	} else {
		goto L504
	}
L503:
	;
	v6357 = v6317
	goto L505
L504:
	;
	v6357 = v6355
	goto L505
L505:
	;
	if base.Ui32(v6357) < base.Ui32(v6279) {
		goto L506
	} else {
		goto L507
	}
L506:
	;
	v6359 = v6279
	goto L508
L507:
	;
	v6359 = v6357
	goto L508
L508:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v6123))) = uint8(v6359)
	v6364 = v6124 + int32(-1)
	if v6364 != 0 {
		__phi6117 = v6133
		__phi6120 = v6189
		__phi6121 = v6120
		__phi6123 = v6123 + int32(1)
		__phi6124 = v6364
		v6117 = __phi6117
		v6120 = __phi6120
		v6121 = __phi6121
		v6123 = __phi6123
		v6124 = __phi6124
		goto L460
	} else {
		goto L509
	}
L509:
	;
	goto L461
L510:
	;
	goto L370
L511:
	;
	if v5154&v6384 == int32(0) {
		goto L370
	} else {
		goto L516
	}
L512:
	;
	v6404 = int32(0)
	v6417 = v5158
	goto L513
L513:
	;
	v6491 = int32(2)
	v6493 = F_memcpy(m, v5159, v6417, (base.B2i32(v6404+int32(1) < v5154)+v5153)<<(uint(v6491)%32))
	mBase = m.M
	v6496 = v6404 + v6491
	v6501 = F_memcpy(m, v5250, v6417+v5153<<(uint(int32(2))%32), (base.B2i32(v6496 < v5154)+v5153)<<(uint(v6491)%32))
	mBase = m.M
	if v5154&int32(2147483646) != v6496 {
		v6404 = v6496
		v6417 = v6417 + v5153<<(uint(int32(3))%32)
		goto L513
	} else {
		goto L515
	}
L514:
	;
	v6512 = v6496
	goto L511
L515:
	;
	goto L514
L516:
	;
	v6598 = int32(2)
	v6607 = F_memcpy(m, v5159, v5158+v6512*v5153<<(uint(v6598)%32), (base.B2i32(v6512+int32(1) < v5154)+v5153)<<(uint(v6598)%32))
	mBase = m.M
	goto L370
}

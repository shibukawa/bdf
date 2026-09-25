//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_PackRGB_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v25 int32
	_ = v25
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v48 base.V128
	_ = v48
	var v55 base.V128
	_ = v55
	var v60 base.V128
	_ = v60
	var v63 int32
	_ = v63
	var v72 int32
	_ = v72
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v107 int32
	_ = v107
	var v111 int32
	_ = v111
	var v116 int32
	_ = v116
	var v125 int32
	_ = v125
	if l3 < int32(1) {
	} else {
		v17 = int32(0)
		if base.Ui32(int32(4)) <= base.Ui32(l3) {
			if l4 != int32(1) {
				v80 = v17
				v81 = int32(0)
				v98 = v80
				v101 = l5 + v81<<(uint(int32(2))%32)
				v102 = l3 - v81
				for {
					v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v98))))
					v111 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v98))))
					v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v98))))
					*(*int32)(unsafe.Add(mBase, uint32(v101))) = v107<<(uint(int32(16))%32) | v111<<(uint(int32(8))%32) | v116 | int32(-16777216)
					v125 = v102 + int32(-1)
					if v125 != 0 {
						v98 = v98 + l4
						v101 = v101 + int32(4)
						v102 = v125
						continue
					} else {
						break
					}
					break
				}
			} else {
				v25 = l3 & int32(2147483644)
				v35 = v25
				v36 = l2
				v37 = l1
				v38 = l0
				v39 = l5
				for {
					v41 = int32(0)
					v42 = base.Simd_g_v128_load32_zero(m, v38, v41)
					v45 = int32(16)
					v48 = base.Simd_g_v128_load32_zero(m, v37, v41)
					v55 = base.Simd_g_v128_load32_zero(m, v36, v41)
					v60 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v42, base.Simd_g_const(&F_PackRGB_C__k0)), v45), base.Simd_g_i32x4_shl(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v48)), int32(8))), base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i16x8_extend_low_i8x16_u(v55))), base.Simd_g_const(&F_PackRGB_C__k1))
					base.Simd_g_v128_store(m, v39, v41, v60)
					v63 = int32(4)
					v72 = v35 + int32(-4)
					if v72 != 0 {
						v35 = v72
						v36 = v36 + v63
						v37 = v37 + v63
						v38 = v38 + v63
						v39 = v39 + v45
						continue
					} else {
						break
					}
					break
				}
				if v25 == l3 {
				} else {
					v80 = v25 * l4
					v81 = v25
					v98 = v80
					v101 = l5 + v81<<(uint(int32(2))%32)
					v102 = l3 - v81
					for {
						v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v98))))
						v111 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v98))))
						v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v98))))
						*(*int32)(unsafe.Add(mBase, uint32(v101))) = v107<<(uint(int32(16))%32) | v111<<(uint(int32(8))%32) | v116 | int32(-16777216)
						v125 = v102 + int32(-1)
						if v125 != 0 {
							v98 = v98 + l4
							v101 = v101 + int32(4)
							v102 = v125
							continue
						} else {
							break
						}
						break
					}
				}
			}
		} else {
			v80 = v17
			v81 = int32(0)
			v98 = v80
			v101 = l5 + v81<<(uint(int32(2))%32)
			v102 = l3 - v81
			for {
				v107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v98))))
				v111 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1+v98))))
				v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2+v98))))
				*(*int32)(unsafe.Add(mBase, uint32(v101))) = v107<<(uint(int32(16))%32) | v111<<(uint(int32(8))%32) | v116 | int32(-16777216)
				v125 = v102 + int32(-1)
				if v125 != 0 {
					v98 = v98 + l4
					v101 = v101 + int32(4)
					v102 = v125
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

var F_PackRGB_C__k0 = [2]uint64{0x100000000, 0x300000002}
var F_PackRGB_C__k1 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}

func F_PaletteSort(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32 {
	mBase := m.M
	_ = mBase
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v43 int32
	_ = v43
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v66 int32
	_ = v66
	var v72 int32
	_ = v72
	var v73 base.V128
	_ = v73
	var v75 int32
	_ = v75
	var v85 int32
	_ = v85
	var v87 base.V128
	_ = v87
	var v89 base.V128
	_ = v89
	var v113 base.V128
	_ = v113
	var v114 base.V128
	_ = v114
	var v117 base.V128
	_ = v117
	var v120 base.V128
	_ = v120
	var v121 base.V128
	_ = v121
	var v122 base.V128
	_ = v122
	var v123 base.V128
	_ = v123
	var v125 base.V128
	_ = v125
	var v139 base.V128
	_ = v139
	var v151 int32
	_ = v151
	var v153 base.V128
	_ = v153
	var v162 base.V128
	_ = v162
	var v166 int32
	_ = v166
	var v169 base.V128
	_ = v169
	var v174 int32
	_ = v174
	var v178 int32
	_ = v178
	var v186 int32
	_ = v186
	var v187 int32
	_ = v187
	var v217 int32
	_ = v217
	var v225 int32
	_ = v225
	var v227 int32
	_ = v227
	var v238 int32
	_ = v238
	var v254 int32
	_ = v254
	var v259 int32
	_ = v259
	var v261 int32
	_ = v261
	var v264 int32
	_ = v264
	var v266 int32
	_ = v266
	var v268 int32
	_ = v268
	var v276 int32
	_ = v276
	var v279 int32
	_ = v279
	var v281 int32
	_ = v281
	var v288 int32
	_ = v288
	var v291 int32
	_ = v291
	var v293 int32
	_ = v293
	var v294 int32
	_ = v294
	var v298 int32
	_ = v298
	var v307 int32
	_ = v307
	var v343 int32
	_ = v343
	var v345 int32
	_ = v345
	var v348 int32
	_ = v348
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
	var v363 int32
	_ = v363
	var v386 int32
	_ = v386
	var v402 int32
	_ = v402
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v424 int32
	_ = v424
	var v437 int32
	_ = v437
	var v440 int32
	_ = v440
	var v444 int32
	_ = v444
	var v449 int32
	_ = v449
	var v452 int32
	_ = v452
	var v454 int32
	_ = v454
	var v459 int32
	_ = v459
	var v464 int32
	_ = v464
	var v469 int32
	_ = v469
	var v474 int32
	_ = v474
	var v479 int32
	_ = v479
	var v480 int32
	_ = v480
	var v482 int32
	_ = v482
	var v486 int32
	_ = v486
	var v488 int32
	_ = v488
	var v490 int32
	_ = v490
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v496 int32
	_ = v496
	var v497 int32
	_ = v497
	var v502 int32
	_ = v502
	var v504 int32
	_ = v504
	var v510 int64
	_ = v510
	var v511 int32
	_ = v511
	var v518 int64
	_ = v518
	var v519 int32
	_ = v519
	var v520 int64
	_ = v520
	var v529 int32
	_ = v529
	var v531 int32
	_ = v531
	var v534 int32
	_ = v534
	var v535 int32
	_ = v535
	var v536 int32
	_ = v536
	var v659 int32
	_ = v659
	var v662 int64
	_ = v662
	var v663 int32
	_ = v663
	var v670 int64
	_ = v670
	var v671 int32
	_ = v671
	var v672 int64
	_ = v672
	var v682 int32
	_ = v682
	var v684 int32
	_ = v684
	var v688 int32
	_ = v688
	var v690 int32
	_ = v690
	var v693 int32
	_ = v693
	var v697 int32
	_ = v697
	var v702 int32
	_ = v702
	var v709 int32
	_ = v709
	var v742 int32
	_ = v742
	var v754 int32
	_ = v754
	var v755 int32
	_ = v755
	var v782 int32
	_ = v782
	var v788 int32
	_ = v788
	var v789 int32
	_ = v789
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v793 int32
	_ = v793
	var v833 int32
	_ = v833
	var v835 int32
	_ = v835
	var v838 int32
	_ = v838
	var v846 int32
	_ = v846
	var v848 int32
	_ = v848
	var __phi848 int32
	_ = __phi848
	var v855 int32
	_ = v855
	var __phi855 int32
	_ = __phi855
	var v857 int32
	_ = v857
	var __phi857 int32
	_ = __phi857
	var v858 int32
	_ = v858
	var __phi858 int32
	_ = __phi858
	var v869 int32
	_ = v869
	var __phi869 int32
	_ = __phi869
	var v876 int32
	_ = v876
	var __phi876 int32
	_ = __phi876
	var v877 int32
	_ = v877
	var __phi877 int32
	_ = __phi877
	var v888 int32
	_ = v888
	var v894 int32
	_ = v894
	var v898 int32
	_ = v898
	var v924 int32
	_ = v924
	var v926 int32
	_ = v926
	var v929 int32
	_ = v929
	var v941 int32
	_ = v941
	var v942 int32
	_ = v942
	var v969 int32
	_ = v969
	var v975 int32
	_ = v975
	var v976 int32
	_ = v976
	var v977 int32
	_ = v977
	var v978 int32
	_ = v978
	var v980 int32
	_ = v980
	var v1018 int32
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1029 int32
	_ = v1029
	var v1054 int32
	_ = v1054
	var v1060 int32
	_ = v1060
	var v1063 int32
	_ = v1063
	var v1068 int32
	_ = v1068
	var v1069 int32
	_ = v1069
	var v1070 int32
	_ = v1070
	var v1079 int32
	_ = v1079
	var v1080 int32
	_ = v1080
	var v1087 int32
	_ = v1087
	var v1090 int32
	_ = v1090
	var v1095 int32
	_ = v1095
	var v1096 int32
	_ = v1096
	var v1097 int32
	_ = v1097
	var v1106 int32
	_ = v1106
	var v1107 int32
	_ = v1107
	var v1114 int32
	_ = v1114
	var v1115 int32
	_ = v1115
	var v1118 int32
	_ = v1118
	var v1124 int32
	_ = v1124
	var v1128 int32
	_ = v1128
	var v1154 int32
	_ = v1154
	var v1156 int32
	_ = v1156
	var v1159 int32
	_ = v1159
	var v1171 int32
	_ = v1171
	var v1172 int32
	_ = v1172
	var v1199 int32
	_ = v1199
	var v1205 int32
	_ = v1205
	var v1206 int32
	_ = v1206
	var v1207 int32
	_ = v1207
	var v1208 int32
	_ = v1208
	var v1210 int32
	_ = v1210
	var v1248 int32
	_ = v1248
	var v1249 int32
	_ = v1249
	var v1259 int32
	_ = v1259
	var v1284 int32
	_ = v1284
	var v1290 int32
	_ = v1290
	var v1293 int32
	_ = v1293
	var v1298 int32
	_ = v1298
	var v1299 int32
	_ = v1299
	var v1300 int32
	_ = v1300
	var v1309 int32
	_ = v1309
	var v1310 int32
	_ = v1310
	var v1317 int32
	_ = v1317
	var v1318 int32
	_ = v1318
	var v1320 int32
	_ = v1320
	var v1329 int32
	_ = v1329
	var v1330 int32
	_ = v1330
	var v1355 int32
	_ = v1355
	var v1360 int32
	_ = v1360
	var v1361 int32
	_ = v1361
	var v1402 int32
	_ = v1402
	var v1403 int32
	_ = v1403
	var v1414 int32
	_ = v1414
	var v1429 int32
	_ = v1429
	var v1430 int32
	_ = v1430
	var v1431 int32
	_ = v1431
	var v1443 int32
	_ = v1443
	var v1446 int32
	_ = v1446
	var v1456 int32
	_ = v1456
	var v1461 base.V128
	_ = v1461
	var v1482 base.V128
	_ = v1482
	var v1483 base.V128
	_ = v1483
	var v1487 int32
	_ = v1487
	var v1490 base.V128
	_ = v1490
	var v1495 int32
	_ = v1495
	var v1497 int32
	_ = v1497
	var v1507 int32
	_ = v1507
	var v1536 int32
	_ = v1536
	var v1545 int32
	_ = v1545
	var v1546 int32
	_ = v1546
	var v1571 int32
	_ = v1571
	var v1572 int32
	_ = v1572
	var v1576 int32
	_ = v1576
	var v1587 int32
	_ = v1587
	var v1612 int32
	_ = v1612
	var v1613 int32
	_ = v1613
	var v1614 int32
	_ = v1614
	var v1617 int32
	_ = v1617
	var v1619 int32
	_ = v1619
	var v1620 int32
	_ = v1620
	var v1623 int32
	_ = v1623
	var v1624 int32
	_ = v1624
	var v1626 int32
	_ = v1626
	var v1629 int32
	_ = v1629
	var v1635 int32
	_ = v1635
	var v1642 int32
	_ = v1642
	var v1651 int32
	_ = v1651
	var v1652 int32
	_ = v1652
	var v1663 int32
	_ = v1663
	var v1679 int32
	_ = v1679
	var v1682 int32
	_ = v1682
	var v1685 int32
	_ = v1685
	var v1686 int32
	_ = v1686
	var v1687 int32
	_ = v1687
	var v1688 int32
	_ = v1688
	var v1689 int32
	_ = v1689
	var v1690 int32
	_ = v1690
	var v1691 int32
	_ = v1691
	var v1692 int32
	_ = v1692
	var v1693 int32
	_ = v1693
	var v1694 int32
	_ = v1694
	var v1701 int32
	_ = v1701
	var v1702 int32
	_ = v1702
	var v1703 int32
	_ = v1703
	var v1704 int32
	_ = v1704
	var v1708 int32
	_ = v1708
	var v1714 int32
	_ = v1714
	var v1724 int32
	_ = v1724
	var v1735 int32
	_ = v1735
	var v1751 int32
	_ = v1751
	var v1757 int32
	_ = v1757
	var v1766 int32
	_ = v1766
	var v1767 int32
	_ = v1767
	var v1778 int32
	_ = v1778
	var v1779 int32
	_ = v1779
	var v1792 int32
	_ = v1792
	var v1793 int32
	_ = v1793
	var v1794 int32
	_ = v1794
	var v1795 int32
	_ = v1795
	var v1801 int32
	_ = v1801
	var v1823 int32
	_ = v1823
	var v1838 int32
	_ = v1838
	var v1841 int32
	_ = v1841
	var v1844 int32
	_ = v1844
	var v1849 int32
	_ = v1849
	var v1854 int32
	_ = v1854
	var v1858 int32
	_ = v1858
	var v1867 int32
	_ = v1867
	var v1868 int32
	_ = v1868
	var v1879 int32
	_ = v1879
	var v1897 int32
	_ = v1897
	var v1899 int32
	_ = v1899
	var v1901 int32
	_ = v1901
	var v1902 int32
	_ = v1902
	var v1905 int32
	_ = v1905
	var v1907 int32
	_ = v1907
	var v1911 int32
	_ = v1911
	var v1912 int32
	_ = v1912
	var v1916 int32
	_ = v1916
	var v1924 int32
	_ = v1924
	var v1925 int32
	_ = v1925
	var v1935 int32
	_ = v1935
	var v1938 int32
	_ = v1938
	var v1954 int32
	_ = v1954
	var v1955 int32
	_ = v1955
	var v1957 int32
	_ = v1957
	var v1961 int32
	_ = v1961
	var v1962 int32
	_ = v1962
	var v1963 int32
	_ = v1963
	var v1964 int32
	_ = v1964
	var v1966 int32
	_ = v1966
	var v1975 int32
	_ = v1975
	var v1984 int32
	_ = v1984
	var v1996 int32
	_ = v1996
	var v1997 int32
	_ = v1997
	var v2013 int32
	_ = v2013
	var v2017 int32
	_ = v2017
	var v2019 int32
	_ = v2019
	var v2020 int32
	_ = v2020
	var v2028 int32
	_ = v2028
	var v2030 int32
	_ = v2030
	var v2038 int32
	_ = v2038
	var v2054 int32
	_ = v2054
	var v2065 int32
	_ = v2065
	var v2071 int32
	_ = v2071
	var v2075 int64
	_ = v2075
	var v2082 int32
	_ = v2082
	var v2088 int32
	_ = v2088
	var v2098 int32
	_ = v2098
	var v2109 int32
	_ = v2109
	var v2124 int32
	_ = v2124
	var v2125 int32
	_ = v2125
	var v2127 int32
	_ = v2127
	var v2128 int32
	_ = v2128
	var v2129 int32
	_ = v2129
	var v2132 int32
	_ = v2132
	var v2133 int32
	_ = v2133
	var v2135 int32
	_ = v2135
	var v2137 int32
	_ = v2137
	var v2138 int32
	_ = v2138
	var v2139 int32
	_ = v2139
	var v2143 int32
	_ = v2143
	var v2144 int32
	_ = v2144
	var v2147 int32
	_ = v2147
	var v2148 int32
	_ = v2148
	var v2150 int32
	_ = v2150
	var v2154 int32
	_ = v2154
	var v2166 int32
	_ = v2166
	var v2177 int32
	_ = v2177
	var v2197 int32
	_ = v2197
	var v2198 int32
	_ = v2198
	var v2199 int32
	_ = v2199
	var v2203 int32
	_ = v2203
	var v2204 int32
	_ = v2204
	var v2206 int32
	_ = v2206
	var v2208 int32
	_ = v2208
	var v2211 int32
	_ = v2211
	var v2212 int32
	_ = v2212
	var v2227 int32
	_ = v2227
	var v2256 int32
	_ = v2256
	var v2264 int32
	_ = v2264
	var v2292 int32
	_ = v2292
	var v2300 int32
	_ = v2300
	var v2301 int32
	_ = v2301
	var v2302 int32
	_ = v2302
	var v2304 int32
	_ = v2304
	var v2305 int32
	_ = v2305
	var v2308 int32
	_ = v2308
	var v2316 int32
	_ = v2316
	var v2318 int32
	_ = v2318
	var v2322 int32
	_ = v2322
	var v2327 int32
	_ = v2327
	var v2357 int32
	_ = v2357
	var v2366 int32
	_ = v2366
	var v2372 int32
	_ = v2372
	var v2374 int32
	_ = v2374
	var v2378 int32
	_ = v2378
	var v2386 int32
	_ = v2386
	v36 = m.G0
	v38 = v36 - int32(3072)
	m.G0 = v38
	switch l0 {
	case 0:
		if base.Ui32(l3) < base.Ui32(int32(18)) {
			v57 = F_memcpy(m, l4, l2, l3<<(uint(int32(2))%32))
			mBase = m.M
			v2386 = int32(1)
		} else {
			v43 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
			if v43 != 0 {
				v57 = F_memcpy(m, l4, l2, l3<<(uint(int32(2))%32))
				mBase = m.M
				v2386 = int32(1)
			} else {
				v49 = l3<<(uint(int32(2))%32) + int32(-4)
				v50 = F_memcpy(m, l4, l2+int32(4), v49)
				mBase = m.M
				*(*int32)(unsafe.Add(mBase, uint32(v50+v49))) = int32(0)
				v2386 = int32(1)
			}
		}
	case 1:
		v61 = F_memcpy(m, l4, l2, l3<<(uint(int32(2))%32))
		mBase = m.M
		v62 = int32(1)
		if v62 <= l3 {
			v66 = int32(0)
			if base.Ui32(l3) < base.Ui32(int32(4)) {
				v178 = v66
				v186 = v66
				v187 = v66
				v217 = v178
				v225 = v186
				v227 = l2 + v187<<(uint(int32(2))%32)
				v238 = l3 - v187
				for {
					v254 = *(*int32)(unsafe.Add(mBase, uint32(v227)))
					v259 = v254 | int32(_a_F_PaletteSort_0) - v217&int32(16711935)
					v261 = v259 & int32(255)
					if base.Ui32(v261) < base.Ui32(int32(128)) {
						v264 = int32(64)
					} else {
						v264 = int32(-128)
					}
					if v261 != 0 {
						v266 = v264
					} else {
						v266 = int32(0)
					}
					v268 = int32(8)
					v276 = int32(base.Ui32(v254-v217&int32(_a_F_PaletteSort_0))>>(uint(v268)%32)) & int32(255)
					if base.Ui32(v276) < base.Ui32(int32(128)) {
						v279 = v268
					} else {
						v279 = int32(16)
					}
					if v276 != 0 {
						v281 = v279
					} else {
						v281 = int32(0)
					}
					v288 = int32(base.Ui32(v259)>>(uint(int32(16))%32)) & int32(255)
					if base.Ui32(v288) < base.Ui32(int32(128)) {
						v291 = int32(1)
					} else {
						v291 = int32(2)
					}
					if v288 != 0 {
						v293 = v291
					} else {
						v293 = int32(0)
					}
					v294 = v266 | v225 | v281 | v293
					v298 = v238 + int32(-1)
					if v298 != 0 {
						v217 = v254
						v225 = v294
						v227 = v227 + int32(4)
						v238 = v298
						continue
					} else {
						break
					}
					break
				}
				v307 = v294
			} else {
				v72 = l3 & int32(2147483644)
				v73 = base.Simd_g_const(&F_PaletteSort__k0)
				v75 = l2
				v85 = v72
				v87 = v73
				v89 = v73
				for {
					v113 = base.Simd_g_v128_load(m, v75, int32(0))
					v114 = base.Simd_g_const(&F_PaletteSort__k1)
					v117 = base.Simd_g_i8x16_shuffle2(v87, v113, base.Simd_g_const(&F_PaletteSort__k2), base.Simd_g_const(&F_PaletteSort__k3))
					v120 = base.Simd_g_i32x4_sub(base.Simd_g_v128_or(v113, v114), base.Simd_g_v128_and(v117, base.Simd_g_const(&F_PaletteSort__k4)))
					v121 = base.Simd_g_const(&F_PaletteSort__k5)
					v122 = base.Simd_g_v128_and(v120, v121)
					v123 = base.Simd_g_const(&F_PaletteSort__k6)
					v125 = base.Simd_g_const(&F_PaletteSort__k7)
					v139 = base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_sub(v113, base.Simd_g_v128_and(v117, v114)), int32(8)), v121)
					v151 = int32(16)
					v153 = base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(v120, v151), v121)
					v162 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_v128_bitselect(v73, base.Simd_g_v128_bitselect(base.Simd_g_const(&F_PaletteSort__k8), base.Simd_g_const(&F_PaletteSort__k9), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_u(v122, v123), v113, base.Simd_g_const(&F_PaletteSort__k7), base.Simd_g_const(&F_PaletteSort__k9))), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_eq(v122, v73), v113, base.Simd_g_const(&F_PaletteSort__k7), base.Simd_g_const(&F_PaletteSort__k9))), v89), base.Simd_g_v128_bitselect(v73, base.Simd_g_v128_bitselect(base.Simd_g_const(&F_PaletteSort__k10), base.Simd_g_const(&F_PaletteSort__k11), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_u(v139, v123), v113, base.Simd_g_const(&F_PaletteSort__k7), base.Simd_g_const(&F_PaletteSort__k9))), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_eq(v139, v73), v113, base.Simd_g_const(&F_PaletteSort__k7), base.Simd_g_const(&F_PaletteSort__k9)))), base.Simd_g_v128_bitselect(v73, base.Simd_g_v128_bitselect(base.Simd_g_const(&F_PaletteSort__k12), base.Simd_g_const(&F_PaletteSort__k13), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_lt_u(v153, v123), v113, base.Simd_g_const(&F_PaletteSort__k7), base.Simd_g_const(&F_PaletteSort__k9))), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_eq(v153, v73), v113, base.Simd_g_const(&F_PaletteSort__k7), base.Simd_g_const(&F_PaletteSort__k9))))
					v166 = v85 + int32(-4)
					if v166 != 0 {
						v75 = v75 + v151
						v85 = v166
						v87 = v113
						v89 = v162
						continue
					} else {
						break
					}
					break
				}
				v169 = base.Simd_g_v128_or(v162, base.Simd_g_i8x16_shuffle2(v162, v113, base.Simd_g_const(&F_PaletteSort__k14), base.Simd_g_const(&F_PaletteSort__k9)))
				v174 = base.Simd_g_i8x16_extract_lane_u_l0(base.Simd_g_v128_or(v169, base.Simd_g_i8x16_shuffle2(v169, v113, base.Simd_g_const(&F_PaletteSort__k15), base.Simd_g_const(&F_PaletteSort__k9))))
				if v72 == l3 {
					v307 = v174
				} else {
					v178 = base.Simd_g_i32x4_extract_lane_l3(v113)
					v186 = v174
					v187 = v72
					v217 = v178
					v225 = v186
					v227 = l2 + v187<<(uint(int32(2))%32)
					v238 = l3 - v187
					for {
						v254 = *(*int32)(unsafe.Add(mBase, uint32(v227)))
						v259 = v254 | int32(_a_F_PaletteSort_0) - v217&int32(16711935)
						v261 = v259 & int32(255)
						if base.Ui32(v261) < base.Ui32(int32(128)) {
							v264 = int32(64)
						} else {
							v264 = int32(-128)
						}
						if v261 != 0 {
							v266 = v264
						} else {
							v266 = int32(0)
						}
						v268 = int32(8)
						v276 = int32(base.Ui32(v254-v217&int32(_a_F_PaletteSort_0))>>(uint(v268)%32)) & int32(255)
						if base.Ui32(v276) < base.Ui32(int32(128)) {
							v279 = v268
						} else {
							v279 = int32(16)
						}
						if v276 != 0 {
							v281 = v279
						} else {
							v281 = int32(0)
						}
						v288 = int32(base.Ui32(v259)>>(uint(int32(16))%32)) & int32(255)
						if base.Ui32(v288) < base.Ui32(int32(128)) {
							v291 = int32(1)
						} else {
							v291 = int32(2)
						}
						if v288 != 0 {
							v293 = v291
						} else {
							v293 = int32(0)
						}
						v294 = v266 | v225 | v281 | v293
						v298 = v238 + int32(-1)
						if v298 != 0 {
							v217 = v254
							v225 = v294
							v227 = v227 + int32(4)
							v238 = v298
							continue
						} else {
							break
						}
						break
					}
					v307 = v294
				}
			}
		} else {
			v307 = int32(0)
		}
		if v307<<(uint(int32(1))%32)&v307&int32(254) == int32(0) {
			v2386 = v62
		} else {
			if l3 < int32(18) {
				v353 = int32(1)
				if l3 < v353 {
					v2386 = v353
				} else {
					v357 = l3
					v360 = int32(0)
					v362 = v360
					v363 = v61
					v386 = v360
					for {
						v402 = int32(-1)
						v411 = v386
						v412 = v363
						v424 = v386
						for {
							v437 = *(*int32)(unsafe.Add(mBase, uint32(v412)))
							v440 = v437 | int32(16711680) - v362&int32(-16711936)
							v444 = int32(base.Ui32(v440)>>(uint(int32(8))%32)) & int32(255)
							if base.Ui32(v444) < base.Ui32(int32(129)) {
								v449 = v444
							} else {
								v449 = int32(256) - v444
							}
							v452 = v437 | int32(_a_F_PaletteSort_0) - v362&int32(16711935)
							v454 = v452 & int32(255)
							if base.Ui32(v454) < base.Ui32(int32(129)) {
								v459 = v454
							} else {
								v459 = int32(256) - v454
							}
							v464 = int32(base.Ui32(v452)>>(uint(int32(16))%32)) & int32(255)
							if base.Ui32(v464) < base.Ui32(int32(129)) {
								v469 = v464
							} else {
								v469 = int32(256) - v464
							}
							v474 = int32(base.Ui32(v440) >> (uint(int32(24)) % 32))
							if base.Ui32(v440) < base.Ui32(int32(-2130706432)) {
								v479 = v474
							} else {
								v479 = int32(256) - v474
							}
							v480 = (v449+v459+v469)*int32(9) + v479
							if base.Ui32(v480) < base.Ui32(v402) {
								v482 = v411
							} else {
								v482 = v424
							}
							if base.Ui32(v402) < base.Ui32(v480) {
								v486 = v402
							} else {
								v486 = v480
							}
							v488 = v411 + int32(1)
							if v357 != v488 {
								v402 = v486
								v411 = v488
								v412 = v412 + int32(4)
								v424 = v482
								continue
							} else {
								break
							}
							break
						}
						v490 = int32(2)
						v492 = v61 + v482<<(uint(v490)%32)
						v493 = *(*int32)(unsafe.Add(mBase, uint32(v492)))
						v496 = v61 + v386<<(uint(v490)%32)
						v497 = *(*int32)(unsafe.Add(mBase, uint32(v496)))
						*(*int32)(unsafe.Add(mBase, uint32(v492))) = v497
						*(*int32)(unsafe.Add(mBase, uint32(v496))) = v493
						v502 = int32(1)
						v504 = v386 + v502
						if v504 != v357 {
							v362 = v493
							v363 = v363 + int32(4)
							v386 = v504
							continue
						} else {
							break
						}
						break
					}
					v2386 = v502
				}
			} else {
				v343 = *(*int32)(unsafe.Add(mBase, uint32(v61)))
				if v343 != 0 {
					v357 = l3
				} else {
					v345 = l3 + int32(-1)
					v348 = v61 + v345<<(uint(int32(2))%32)
					v349 = *(*int32)(unsafe.Add(mBase, uint32(v348)))
					*(*int32)(unsafe.Add(mBase, uint32(v348))) = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(v61))) = v349
					v357 = v345
				}
				v360 = int32(0)
				v362 = v360
				v363 = v61
				v386 = v360
				for {
					v402 = int32(-1)
					v411 = v386
					v412 = v363
					v424 = v386
					for {
						v437 = *(*int32)(unsafe.Add(mBase, uint32(v412)))
						v440 = v437 | int32(16711680) - v362&int32(-16711936)
						v444 = int32(base.Ui32(v440)>>(uint(int32(8))%32)) & int32(255)
						if base.Ui32(v444) < base.Ui32(int32(129)) {
							v449 = v444
						} else {
							v449 = int32(256) - v444
						}
						v452 = v437 | int32(_a_F_PaletteSort_0) - v362&int32(16711935)
						v454 = v452 & int32(255)
						if base.Ui32(v454) < base.Ui32(int32(129)) {
							v459 = v454
						} else {
							v459 = int32(256) - v454
						}
						v464 = int32(base.Ui32(v452)>>(uint(int32(16))%32)) & int32(255)
						if base.Ui32(v464) < base.Ui32(int32(129)) {
							v469 = v464
						} else {
							v469 = int32(256) - v464
						}
						v474 = int32(base.Ui32(v440) >> (uint(int32(24)) % 32))
						if base.Ui32(v440) < base.Ui32(int32(-2130706432)) {
							v479 = v474
						} else {
							v479 = int32(256) - v474
						}
						v480 = (v449+v459+v469)*int32(9) + v479
						if base.Ui32(v480) < base.Ui32(v402) {
							v482 = v411
						} else {
							v482 = v424
						}
						if base.Ui32(v402) < base.Ui32(v480) {
							v486 = v402
						} else {
							v486 = v480
						}
						v488 = v411 + int32(1)
						if v357 != v488 {
							v402 = v486
							v411 = v488
							v412 = v412 + int32(4)
							v424 = v482
							continue
						} else {
							break
						}
						break
					}
					v490 = int32(2)
					v492 = v61 + v482<<(uint(v490)%32)
					v493 = *(*int32)(unsafe.Add(mBase, uint32(v492)))
					v496 = v61 + v386<<(uint(v490)%32)
					v497 = *(*int32)(unsafe.Add(mBase, uint32(v496)))
					*(*int32)(unsafe.Add(mBase, uint32(v492))) = v497
					*(*int32)(unsafe.Add(mBase, uint32(v496))) = v493
					v502 = int32(1)
					v504 = v386 + v502
					if v504 != v357 {
						v362 = v493
						v363 = v363 + int32(4)
						v386 = v504
						continue
					} else {
						break
					}
					break
				}
				v2386 = v502
			}
		}
	case 2:
		if base.Ui32(int32(2)) <= base.Ui32(l3) {
			v510 = base.I64_extend_i32_u(l3 * l3)
			v511 = int32(4)
			if v510 == int64(0) {
				v529 = F_calloc(m, base.I32_wrap_i64(v510), v511)
				mBase = m.M
				v531 = v529
			} else {
				v518 = base.I64_div_u_s(int64(2147418112), v510)
				v519 = int32(0)
				v520 = base.I64_extend_i32_u(v511)
				if base.Ui64(int64(4294967295)) < base.Ui64(v520*v510) {
					v531 = v519
				} else {
					if base.Ui64(v518) < base.Ui64(v520) {
						v531 = v519
					} else {
						v529 = F_calloc(m, base.I32_wrap_i64(v510), v511)
						mBase = m.M
						v531 = v529
					}
				}
			}
			if v531 != 0 {
				v534 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
				v535 = *(*int32)(unsafe.Add(mBase, uint32(v534)))
				v536 = int32(0)
				base.MemoryFill(m, v38, v536, int32(1024))
				v659 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
				v662 = base.I64_extend_i32_s(v659 << (uint(int32(1)) % 32))
				v663 = int32(4)
				if v662 == int64(0) {
					v682 = F_malloc(m, base.I32_wrap_i64(v662)*v663)
					mBase = m.M
					v684 = v682
				} else {
					v670 = base.I64_div_u_s(int64(2147418112), v662)
					v671 = int32(0)
					v672 = base.I64_extend_i32_u(v663)
					if base.Ui64(int64(4294967295)) < base.Ui64(v672*v662) {
						v684 = v671
					} else {
						if base.Ui64(v670) < base.Ui64(v672) {
							v684 = v671
						} else {
							v682 = F_malloc(m, base.I32_wrap_i64(v662)*v663)
							mBase = m.M
							v684 = v682
						}
					}
				}
				if v684 == int32(0) {
					F_free(m, v531)
					mBase = m.M
					v2386 = int32(0)
				} else {
					v688 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
					v690 = v38 + int32(2048)
					v693 = F_memcpy(m, v690, l2, l3<<(uint(int32(2))%32))
					mBase = m.M
					v697 = m.G2
					F___qsort_r(m, v690, l3, int32(4), int32(344), v697+int32(328))
					mBase = m.M
					v702 = *(*int32)(unsafe.Add(mBase, uint32(v38)+2048))
					v709 = v536
					for {
						v742 = *(*int32)(unsafe.Add(mBase, uint32(l2+v709<<(uint(int32(2))%32))))
						if v702 == v742 {
							v793 = int32(0)
						} else {
							v754 = l3
							v755 = int32(0)
							for {
								v782 = (v754 + v755) >> (uint(int32(1)) % 32)
								v788 = *(*int32)(unsafe.Add(mBase, uint32(v38+int32(2048)+v782<<(uint(int32(2))%32))))
								v789 = base.B2i32(base.Ui32(v788) < base.Ui32(v742))
								if base.Ui32(v788) < base.Ui32(v742) {
									v790 = v754
								} else {
									v790 = v782
								}
								if base.Ui32(v788) < base.Ui32(v742) {
									v791 = v782
								} else {
									v791 = v755
								}
								if v788 != v742 {
									v754 = v790
									v755 = v791
									continue
								} else {
									break
								}
								break
							}
							v793 = v782
						}
						*(*int32)(unsafe.Add(mBase, uint32(v38+v793<<(uint(int32(2))%32)))) = v709
						v833 = v709 + int32(1)
						if v833 != l3 {
							v709 = v833
							continue
						} else {
							break
						}
						break
					}
					v835 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
					if v835 < int32(1) {
					} else {
						v838 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
						if v838 < int32(1) {
						} else {
							v846 = int32(0)
							__phi848 = v846
							__phi855 = v534
							__phi857 = v838
							__phi858 = v535 ^ int32(-1)
							__phi869 = v684
							__phi876 = v684 + v688<<(uint(int32(2))%32)
							__phi877 = v846
							v848 = __phi848
							v855 = __phi855
							v857 = __phi857
							v858 = __phi858
							v869 = __phi869
							v876 = __phi876
							v877 = __phi877
							for {
								if v857 < int32(1) {
									v1320 = v848
									v1329 = v857
									v1330 = v858
								} else {
									if v877 == int32(0) {
										v1118 = v848
										v1124 = int32(0)
										v1128 = v858
										for {
											v1154 = v1124 << (uint(int32(2)) % 32)
											v1156 = *(*int32)(unsafe.Add(mBase, uint32(v855+v1154)))
											if v1156 == v1128 {
												v1249 = v1118
												v1259 = v1128
											} else {
												v1159 = *(*int32)(unsafe.Add(mBase, uint32(v38)+2048))
												if v1159 == v1156 {
													v1210 = int32(0)
												} else {
													v1171 = l3
													v1172 = int32(0)
													for {
														v1199 = (v1171 + v1172) >> (uint(int32(1)) % 32)
														v1205 = *(*int32)(unsafe.Add(mBase, uint32(v38+int32(2048)+v1199<<(uint(int32(2))%32))))
														v1206 = base.B2i32(base.Ui32(v1205) < base.Ui32(v1156))
														if base.Ui32(v1205) < base.Ui32(v1156) {
															v1207 = v1171
														} else {
															v1207 = v1199
														}
														if base.Ui32(v1205) < base.Ui32(v1156) {
															v1208 = v1199
														} else {
															v1208 = v1172
														}
														if v1205 != v1156 {
															v1171 = v1207
															v1172 = v1208
															continue
														} else {
															break
														}
														break
													}
													v1210 = v1199
												}
												v1248 = *(*int32)(unsafe.Add(mBase, uint32(v38+v1210<<(uint(int32(2))%32))))
												v1249 = v1248
												v1259 = v1156
											}
											v1284 = v876 + v1154
											*(*int32)(unsafe.Add(mBase, uint32(v1284))) = v1249
											if v1124 == int32(0) {
											} else {
												v1290 = *(*int32)(unsafe.Add(mBase, uint32(v1284+int32(-4))))
												if v1249 == v1290 {
												} else {
													v1293 = int32(2)
													v1298 = v531 + v1249*l3<<(uint(v1293)%32) + v1290<<(uint(v1293)%32)
													v1299 = *(*int32)(unsafe.Add(mBase, uint32(v1298)))
													v1300 = int32(1)
													*(*int32)(unsafe.Add(mBase, uint32(v1298))) = v1299 + v1300
													v1309 = v531 + v1290*l3<<(uint(v1293)%32) + v1249<<(uint(v1293)%32)
													v1310 = *(*int32)(unsafe.Add(mBase, uint32(v1309)))
													*(*int32)(unsafe.Add(mBase, uint32(v1309))) = v1310 + v1300
												}
											}
											v1317 = v1124 + int32(1)
											v1318 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
											if v1317 < v1318 {
												v1118 = v1249
												v1124 = v1317
												v1128 = v1259
												continue
											} else {
												break
											}
											break
										}
										v1320 = v1249
										v1329 = v1318
										v1330 = v1259
									} else {
										v888 = v848
										v894 = int32(0)
										v898 = v858
										for {
											v924 = v894 << (uint(int32(2)) % 32)
											v926 = *(*int32)(unsafe.Add(mBase, uint32(v855+v924)))
											if v926 == v898 {
												v1019 = v888
												v1029 = v898
											} else {
												v929 = *(*int32)(unsafe.Add(mBase, uint32(v38)+2048))
												if v929 == v926 {
													v980 = int32(0)
												} else {
													v941 = l3
													v942 = int32(0)
													for {
														v969 = (v941 + v942) >> (uint(int32(1)) % 32)
														v975 = *(*int32)(unsafe.Add(mBase, uint32(v38+int32(2048)+v969<<(uint(int32(2))%32))))
														v976 = base.B2i32(base.Ui32(v975) < base.Ui32(v926))
														if base.Ui32(v975) < base.Ui32(v926) {
															v977 = v941
														} else {
															v977 = v969
														}
														if base.Ui32(v975) < base.Ui32(v926) {
															v978 = v969
														} else {
															v978 = v942
														}
														if v975 != v926 {
															v941 = v977
															v942 = v978
															continue
														} else {
															break
														}
														break
													}
													v980 = v969
												}
												v1018 = *(*int32)(unsafe.Add(mBase, uint32(v38+v980<<(uint(int32(2))%32))))
												v1019 = v1018
												v1029 = v926
											}
											v1054 = v876 + v924
											*(*int32)(unsafe.Add(mBase, uint32(v1054))) = v1019
											if v894 == int32(0) {
											} else {
												v1060 = *(*int32)(unsafe.Add(mBase, uint32(v1054+int32(-4))))
												if v1019 == v1060 {
												} else {
													v1063 = int32(2)
													v1068 = v531 + v1019*l3<<(uint(v1063)%32) + v1060<<(uint(v1063)%32)
													v1069 = *(*int32)(unsafe.Add(mBase, uint32(v1068)))
													v1070 = int32(1)
													*(*int32)(unsafe.Add(mBase, uint32(v1068))) = v1069 + v1070
													v1079 = v531 + v1060*l3<<(uint(v1063)%32) + v1019<<(uint(v1063)%32)
													v1080 = *(*int32)(unsafe.Add(mBase, uint32(v1079)))
													*(*int32)(unsafe.Add(mBase, uint32(v1079))) = v1080 + v1070
												}
											}
											v1087 = *(*int32)(unsafe.Add(mBase, uint32(v869+v924)))
											if v1019 == v1087 {
											} else {
												v1090 = int32(2)
												v1095 = v531 + v1019*l3<<(uint(v1090)%32) + v1087<<(uint(v1090)%32)
												v1096 = *(*int32)(unsafe.Add(mBase, uint32(v1095)))
												v1097 = int32(1)
												*(*int32)(unsafe.Add(mBase, uint32(v1095))) = v1096 + v1097
												v1106 = v531 + v1087*l3<<(uint(v1090)%32) + v1019<<(uint(v1090)%32)
												v1107 = *(*int32)(unsafe.Add(mBase, uint32(v1106)))
												*(*int32)(unsafe.Add(mBase, uint32(v1106))) = v1107 + v1097
											}
											v1114 = v894 + int32(1)
											v1115 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
											if v1114 < v1115 {
												v888 = v1019
												v894 = v1114
												v898 = v1029
												continue
											} else {
												break
											}
											break
										}
										v1320 = v1019
										v1329 = v1115
										v1330 = v1029
									}
								}
								v1355 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
								v1360 = v877 + int32(1)
								v1361 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
								if v1360 < v1361 {
									__phi848 = v1320
									__phi855 = v855 + v1355<<(uint(int32(2))%32)
									__phi857 = v1329
									__phi858 = v1330
									__phi869 = v876
									__phi876 = v869
									__phi877 = v1360
									v848 = __phi848
									v855 = __phi855
									v857 = __phi857
									v858 = __phi858
									v869 = __phi869
									v876 = __phi876
									v877 = __phi877
									continue
								} else {
									break
								}
								break
							}
						}
					}
					F_free(m, v684)
					mBase = m.M
					v1402 = l3 & int32(-4)
					v1403 = int32(0)
					v1414 = v1403
					v1429 = v531
					v1430 = v1403
					v1431 = v1403
					for {
						v1443 = int32(0)
						if base.Ui32(l3) < base.Ui32(int32(4)) {
							v1497 = v1443
							v1507 = v1443
							v1536 = v1429 + v1497<<(uint(int32(2))%32)
							v1545 = l3 - v1497
							v1546 = v1507
							for {
								v1571 = *(*int32)(unsafe.Add(mBase, uint32(v1536)))
								v1572 = v1571 + v1546
								v1576 = v1545 + int32(-1)
								if v1576 != 0 {
									v1536 = v1536 + int32(4)
									v1545 = v1576
									v1546 = v1572
									continue
								} else {
									break
								}
								break
							}
							v1587 = v1572
						} else {
							v1446 = v1429
							v1456 = v1402
							v1461 = base.Simd_g_const(&F_PaletteSort__k0)
							for {
								v1482 = base.Simd_g_v128_load(m, v1446, int32(0))
								v1483 = base.Simd_g_i32x4_add(v1482, v1461)
								v1487 = v1456 + int32(-4)
								if v1487 != 0 {
									v1446 = v1446 + int32(16)
									v1456 = v1487
									v1461 = v1483
									continue
								} else {
									break
								}
								break
							}
							v1490 = base.Simd_g_i32x4_add(v1483, base.Simd_g_i8x16_shuffle2(v1483, v1483, base.Simd_g_const(&F_PaletteSort__k16), base.Simd_g_const(&F_PaletteSort__k9)))
							v1495 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v1490, base.Simd_g_i8x16_shuffle2(v1490, v1490, base.Simd_g_const(&F_PaletteSort__k17), base.Simd_g_const(&F_PaletteSort__k9))))
							if v1402 == l3 {
								v1587 = v1495
							} else {
								v1497 = v1402
								v1507 = v1495
								v1536 = v1429 + v1497<<(uint(int32(2))%32)
								v1545 = l3 - v1497
								v1546 = v1507
								for {
									v1571 = *(*int32)(unsafe.Add(mBase, uint32(v1536)))
									v1572 = v1571 + v1546
									v1576 = v1545 + int32(-1)
									if v1576 != 0 {
										v1536 = v1536 + int32(4)
										v1545 = v1576
										v1546 = v1572
										continue
									} else {
										break
									}
									break
								}
								v1587 = v1572
							}
						}
						v1612 = base.B2i32(base.Ui32(v1431) < base.Ui32(v1587))
						if base.Ui32(v1431) < base.Ui32(v1587) {
							v1613 = v1587
						} else {
							v1613 = v1431
						}
						if base.Ui32(v1431) < base.Ui32(v1587) {
							v1614 = v1430
						} else {
							v1614 = v1414
						}
						v1617 = v1430 + int32(1)
						if v1617 != l3 {
							v1414 = v1614
							v1429 = v1429 + l3<<(uint(int32(2))%32)
							v1430 = v1617
							v1431 = v1613
							continue
						} else {
							break
						}
						break
					}
					v1619 = int32(3)
					v1620 = l3 & v1619
					*(*uint8)(unsafe.Add(mBase, uint32(v38)+2048)) = uint8(v1614)
					v1623 = v1614 & int32(255)
					v1624 = v1623 * l3
					v1626 = l3 + int32(-1)
					if base.Ui32(v1619) <= base.Ui32(v1626) {
						v1635 = int32(0)
						v1642 = v1635
						v1651 = v531 + v1624<<(uint(int32(2))%32)
						v1652 = int32(3)
						v1663 = v1635
						for {
							v1679 = *(*int32)(unsafe.Add(mBase, uint32(v1651+int32(12))))
							v1682 = *(*int32)(unsafe.Add(mBase, uint32(v1651+int32(8))))
							v1685 = *(*int32)(unsafe.Add(mBase, uint32(v1651+int32(4))))
							v1686 = *(*int32)(unsafe.Add(mBase, uint32(v1651)))
							v1687 = base.B2i32(base.Ui32(v1642) < base.Ui32(v1686))
							if base.Ui32(v1642) < base.Ui32(v1686) {
								v1688 = v1686
							} else {
								v1688 = v1642
							}
							v1689 = base.B2i32(base.Ui32(v1688) < base.Ui32(v1685))
							if base.Ui32(v1688) < base.Ui32(v1685) {
								v1690 = v1685
							} else {
								v1690 = v1688
							}
							v1691 = base.B2i32(base.Ui32(v1690) < base.Ui32(v1682))
							if base.Ui32(v1690) < base.Ui32(v1682) {
								v1692 = v1682
							} else {
								v1692 = v1690
							}
							v1693 = base.B2i32(base.Ui32(v1692) < base.Ui32(v1679))
							if base.Ui32(v1692) < base.Ui32(v1679) {
								v1694 = v1679
							} else {
								v1694 = v1692
							}
							if base.Ui32(v1642) < base.Ui32(v1686) {
								v1701 = v1652 + int32(-3)
							} else {
								v1701 = v1663
							}
							if base.Ui32(v1688) < base.Ui32(v1685) {
								v1702 = v1652 + int32(-2)
							} else {
								v1702 = v1701
							}
							if base.Ui32(v1690) < base.Ui32(v1682) {
								v1703 = v1652 + int32(-1)
							} else {
								v1703 = v1702
							}
							if base.Ui32(v1692) < base.Ui32(v1679) {
								v1704 = v1652
							} else {
								v1704 = v1703
							}
							v1708 = v1652 + int32(4)
							if v1635-l3&int32(-4)+v1708 != int32(3) {
								v1642 = v1694
								v1651 = v1651 + int32(16)
								v1652 = v1708
								v1663 = v1704
								continue
							} else {
								break
							}
							break
						}
						v1714 = v1694
						v1724 = v1652 + int32(1)
						v1735 = v1704
					} else {
						v1629 = int32(0)
						v1714 = v1629
						v1724 = v1629
						v1735 = v1629
					}
					if v1620 == int32(0) {
						v1823 = v1735
					} else {
						v1751 = int32(2)
						v1757 = v1714
						v1766 = v531 + v1724<<(uint(v1751)%32) + v1624<<(uint(v1751)%32)
						v1767 = v1724
						v1778 = v1735
						v1779 = v1620
						for {
							v1792 = *(*int32)(unsafe.Add(mBase, uint32(v1766)))
							v1793 = base.B2i32(base.Ui32(v1757) < base.Ui32(v1792))
							if base.Ui32(v1757) < base.Ui32(v1792) {
								v1794 = v1792
							} else {
								v1794 = v1757
							}
							if base.Ui32(v1757) < base.Ui32(v1792) {
								v1795 = v1767
							} else {
								v1795 = v1778
							}
							v1801 = v1779 + int32(-1)
							if v1801 != 0 {
								v1757 = v1794
								v1766 = v1766 + int32(4)
								v1767 = v1767 + int32(1)
								v1778 = v1795
								v1779 = v1801
								continue
							} else {
								break
							}
							break
						}
						v1823 = v1795
					}
					*(*uint8)(unsafe.Add(mBase, uint32(v38)+2049)) = uint8(v1823)
					v1838 = int32(0)
					v1841 = l3 + int32(-2)
					if v1841 == v1838 {
						v2227 = v1838
					} else {
						v1844 = int32(0)
						*(*int32)(unsafe.Add(mBase, uint32(v38)+4)) = v1844
						*(*uint8)(unsafe.Add(mBase, uint32(v38))) = uint8(v1844)
						v1849 = int32(2)
						v1854 = v1823 & int32(255)
						v1858 = v1844
						v1867 = v531
						v1868 = v38
						v1879 = v1844
						for {
							if v1623 == v1858 {
								v1911 = v1868
								v1912 = v1879
							} else {
								if v1854 == v1858 {
									v1911 = v1868
									v1912 = v1879
								} else {
									v1897 = v38 + v1879<<(uint(int32(3))%32)
									v1899 = *(*int32)(unsafe.Add(mBase, uint32(v1867+v1854<<(uint(v1849)%32))))
									v1901 = *(*int32)(unsafe.Add(mBase, uint32(v1867+v1623<<(uint(v1849)%32))))
									v1902 = v1899 + v1901
									*(*int32)(unsafe.Add(mBase, uint32(v1897)+4)) = v1902
									*(*uint8)(unsafe.Add(mBase, uint32(v1897))) = uint8(v1858)
									v1905 = *(*int32)(unsafe.Add(mBase, uint32(v1868)+4))
									if base.Ui32(v1905) < base.Ui32(v1902) {
										v1907 = v1897
									} else {
										v1907 = v1868
									}
									v1911 = v1907
									v1912 = v1879 + int32(1)
								}
							}
							v1916 = v1858 + int32(1)
							if l3 != v1916 {
								v1858 = v1916
								v1867 = v1867 + l3<<(uint(v1849)%32)
								v1868 = v1911
								v1879 = v1912
								continue
							} else {
								break
							}
							break
						}
						v1924 = int32(0)
						v1925 = int32(1)
						v1935 = v1924
						v1938 = v1911
						v1954 = l3 + int32(-3)
						v1955 = v1924
						v1957 = v1841
						v1961 = v1925
						v1962 = v1925
						for {
							v1963 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1938))))
							v1964 = base.I32_rem_u_s(v1935, l3)
							v1966 = v1962 + int32(1)
							if v1964 != v1966 {
								v1975 = v1964
								v1984 = v1935 + int32(1)
								v1996 = v1961
								v1997 = int32(0)
								for {
									v2013 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(2048)+v1975))))
									v2017 = *(*int32)(unsafe.Add(mBase, uint32(v531+v1963*l3<<(uint(int32(2))%32)+v2013<<(uint(int32(2))%32))))
									v2019 = v2017*v1996 + v1997
									v2020 = base.I32_rem_u_s(v1984, l3)
									if v2020 != v1966 {
										v1975 = v2020
										v1984 = v1984 + int32(1)
										v1996 = v1996 + int32(-2)
										v1997 = v2019
										continue
									} else {
										break
									}
									break
								}
								if int32(1) <= v2019 {
									if v1935 != 0 {
										v2028 = v1935
									} else {
										v2028 = l3
									}
									v2030 = v2028 + int32(-1)
									v2038 = v2030
									v2054 = v2030
									v2065 = v1962
								} else {
									v2038 = v1935
									v2054 = v1966
									v2065 = v1966
								}
							} else {
								v2038 = v1935
								v2054 = v1966
								v2065 = v1966
							}
							*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(2048)+v2054))) = uint8(v1963)
							v2071 = v1957 + int32(-1)
							v2075 = *(*int64)(unsafe.Add(mBase, uint32(v38+v2071<<(uint(int32(3))%32))))
							*(*int64)(unsafe.Add(mBase, uint32(v1938))) = v2075
							if v2071 == int32(0) {
								v2227 = v2038
								break
							} else {
								v2082 = v531 + v1963*l3<<(uint(int32(2))%32)
								if l3+int32(-4) != v1955 {
									v2088 = v38 | int32(8)
									v2098 = v38
									v2109 = int32(0)
									for {
										v2124 = v2088 + int32(-4)
										v2125 = *(*int32)(unsafe.Add(mBase, uint32(v2124)))
										v2127 = v2088 + int32(-8)
										v2128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2127))))
										v2129 = int32(2)
										v2132 = *(*int32)(unsafe.Add(mBase, uint32(v2082+v2128<<(uint(v2129)%32))))
										v2133 = v2125 + v2132
										*(*int32)(unsafe.Add(mBase, uint32(v2124))) = v2133
										v2135 = *(*int32)(unsafe.Add(mBase, uint32(v2098)+4))
										v2137 = v2088 + int32(4)
										v2138 = *(*int32)(unsafe.Add(mBase, uint32(v2137)))
										v2139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2088))))
										v2143 = *(*int32)(unsafe.Add(mBase, uint32(v2082+v2139<<(uint(v2129)%32))))
										v2144 = v2138 + v2143
										*(*int32)(unsafe.Add(mBase, uint32(v2137))) = v2144
										if base.Ui32(v2135) < base.Ui32(v2133) {
											v2147 = v2127
										} else {
											v2147 = v2098
										}
										v2148 = *(*int32)(unsafe.Add(mBase, uint32(v2147)+4))
										if base.Ui32(v2148) < base.Ui32(v2144) {
											v2150 = v2088
										} else {
											v2150 = v2147
										}
										v2154 = v2109 + int32(2)
										if v1954&int32(-2) != v2154 {
											v2088 = v2088 + int32(16)
											v2098 = v2150
											v2109 = v2154
											continue
										} else {
											break
										}
										break
									}
									v2166 = v2150
									v2177 = v2154
								} else {
									v2166 = v38
									v2177 = int32(0)
								}
								if v1954&int32(1) == int32(0) {
									v2211 = v2166
								} else {
									v2197 = v38 + v2177<<(uint(int32(3))%32)
									v2198 = *(*int32)(unsafe.Add(mBase, uint32(v2197)+4))
									v2199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2197))))
									v2203 = *(*int32)(unsafe.Add(mBase, uint32(v2082+v2199<<(uint(int32(2))%32))))
									v2204 = v2198 + v2203
									*(*int32)(unsafe.Add(mBase, uint32(v2197)+4)) = v2204
									v2206 = *(*int32)(unsafe.Add(mBase, uint32(v2166)+4))
									if base.Ui32(v2206) < base.Ui32(v2204) {
										v2208 = v2197
									} else {
										v2208 = v2166
									}
									v2211 = v2208
								}
								v2212 = int32(1)
								v1935 = v2038
								v1938 = v2211
								v1954 = v1954 + int32(-1)
								v1955 = v1955 + v2212
								v1957 = v2071
								v1961 = v1961 + v2212
								v1962 = v2065
								continue
							}
							break
						}
					}
					F_free(m, v531)
					mBase = m.M
					v2256 = int32(1)
					if v1626 == int32(0) {
						v2357 = v1838
					} else {
						v2264 = l4
						v2292 = int32(0)
						for {
							v2300 = v38 + int32(2048)
							v2301 = v2227 + v2292
							v2302 = base.I32_rem_u_s(v2301, l3)
							v2304 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2300+v2302))))
							v2305 = int32(2)
							v2308 = *(*int32)(unsafe.Add(mBase, uint32(l2+v2304<<(uint(v2305)%32))))
							*(*int32)(unsafe.Add(mBase, uint32(v2264))) = v2308
							v2316 = base.I32_rem_u_s(v2301+int32(1), l3)
							v2318 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2300+v2316))))
							v2322 = *(*int32)(unsafe.Add(mBase, uint32(l2+v2318<<(uint(v2305)%32))))
							*(*int32)(unsafe.Add(mBase, uint32(v2264+int32(4)))) = v2322
							v2327 = v2292 + v2305
							if l3&int32(-2) != v2327 {
								v2264 = v2264 + int32(8)
								v2292 = v2327
								continue
							} else {
								break
							}
							break
						}
						v2357 = v2327
					}
					if l3&v2256 == int32(0) {
						v2386 = v2256
					} else {
						v2366 = int32(2)
						v2372 = base.I32_rem_u_s(v2357+v2227, l3)
						v2374 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v38+int32(2048)+v2372))))
						v2378 = *(*int32)(unsafe.Add(mBase, uint32(l2+v2374<<(uint(v2366)%32))))
						*(*int32)(unsafe.Add(mBase, uint32(l4+v2357<<(uint(v2366)%32)))) = v2378
						v2386 = v2256
					}
				}
			} else {
				v2386 = int32(0)
			}
		} else {
			v2386 = int32(1)
		}
	default:
		v2386 = int32(0)
	}
	m.G0 = v38 + int32(3072)
	return v2386
}

var F_PaletteSort__k0 = [2]uint64{0x0, 0x0}
var F_PaletteSort__k1 = [2]uint64{0xff000000ff00, 0xff000000ff00}
var F_PaletteSort__k2 = [2]uint64{0x808080800f0e0d0c, 0x8080808080808080}
var F_PaletteSort__k3 = [2]uint64{0x302010080808080, 0xb0a090807060504}
var F_PaletteSort__k4 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_PaletteSort__k5 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_PaletteSort__k6 = [2]uint64{0x8000000080, 0x8000000080}
var F_PaletteSort__k7 = [2]uint64{0xc080400, 0x0}
var F_PaletteSort__k8 = [2]uint64{0x4040404040404040, 0x4040404040404040}
var F_PaletteSort__k9 = [2]uint64{0x8080808080808080, 0x8080808080808080}
var F_PaletteSort__k10 = [2]uint64{0x808080808080808, 0x808080808080808}
var F_PaletteSort__k11 = [2]uint64{0x1010101010101010, 0x1010101010101010}
var F_PaletteSort__k12 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_PaletteSort__k13 = [2]uint64{0x202020202020202, 0x202020202020202}
var F_PaletteSort__k14 = [2]uint64{0x302, 0x0}
var F_PaletteSort__k15 = [2]uint64{0x1, 0x0}
var F_PaletteSort__k16 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_PaletteSort__k17 = [2]uint64{0x302010007060504, 0x302010003020100}

func F_Predictor10_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v6 base.V128
	_ = v6
	var v8 base.V128
	_ = v8
	var v10 int32
	_ = v10
	var v11 base.V128
	_ = v11
	var v15 int32
	_ = v15
	var v16 base.V128
	_ = v16
	var v20 base.V128
	_ = v20
	var v24 base.V128
	_ = v24
	var v29 base.V128
	_ = v29
	var v34 base.V128
	_ = v34
	v3 = base.Simd_g_const(&F_Predictor10_SSE2__k0)
	v6 = base.Simd_g_v128_load32_zero(m, l1, int32(4))
	v8 = base.Simd_g_const(&F_Predictor10_SSE2__k1)
	v10 = int32(0)
	v11 = base.Simd_g_v128_load32_zero(m, l1, v10)
	v15 = int32(1)
	v16 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v6, v3, base.Simd_g_const(&F_Predictor10_SSE2__k2), base.Simd_g_const(&F_Predictor10_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v11, v3, base.Simd_g_const(&F_Predictor10_SSE2__k2), base.Simd_g_const(&F_Predictor10_SSE2__k3))), v15)
	v20 = base.Simd_g_v128_load32_zero(m, l1+int32(-4), v10)
	v24 = base.Simd_g_v128_load32_zero(m, l0, v10)
	v29 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v20, v3, base.Simd_g_const(&F_Predictor10_SSE2__k2), base.Simd_g_const(&F_Predictor10_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v24, v3, base.Simd_g_const(&F_Predictor10_SSE2__k2), base.Simd_g_const(&F_Predictor10_SSE2__k3))), v15)
	v34 = base.Simd_g_i16x8_add(base.Simd_g_v128_and(v16, v29), base.Simd_g_i16x8_shr_u(base.Simd_g_v128_xor(v16, v29), v15))
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_narrow_i16x8_u(v34, v34))
}

var F_Predictor10_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Predictor10_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Predictor10_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Predictor10_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_Predictor11_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 base.V128
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v12 base.V128
	_ = v12
	var v15 base.V128
	_ = v15
	var v20 base.V128
	_ = v20
	var v23 base.V128
	_ = v23
	var v29 base.V128
	_ = v29
	var v32 int32
	_ = v32
	var v43 int32
	_ = v43
	v3 = int32(0)
	v4 = base.Simd_g_const(&F_Predictor11_SSE2__k0)
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v12 = base.Simd_g_v128_load32_zero(m, l1+int32(-4), v3)
	v15 = base.Simd_g_i32x4_replace_lane_l0(v4, v8)
	v20 = base.Simd_g_const(&F_Predictor11_SSE2__k1)
	v23 = base.Simd_g_i32x4_replace_lane_l0(v4, v7)
	v29 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v12, v15), base.Simd_g_i8x16_sub_sat_u(v15, v12)), v4, base.Simd_g_const(&F_Predictor11_SSE2__k2), base.Simd_g_const(&F_Predictor11_SSE2__k3)), base.Simd_g_i8x16_shuffle2(base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v12, v23), base.Simd_g_i8x16_sub_sat_u(v23, v12)), v4, base.Simd_g_const(&F_Predictor11_SSE2__k2), base.Simd_g_const(&F_Predictor11_SSE2__k3)))
	v32 = int32(1)
	if base.Simd_g_i16x8_extract_lane_s_l0(v29)+base.Simd_g_i16x8_extract_lane_s_l1(v29)+base.Simd_g_i16x8_extract_lane_s_l2(v29)+base.Simd_g_i16x8_extract_lane_s_l3(v29) < v32 {
		v43 = v7
	} else {
		v43 = v8
	}
	return v43
}

var F_Predictor11_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Predictor11_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Predictor11_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Predictor11_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_Predictor12_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v4 int32
	_ = v4
	var v5 base.V128
	_ = v5
	var v7 base.V128
	_ = v7
	var v10 base.V128
	_ = v10
	var v17 base.V128
	_ = v17
	var v20 base.V128
	_ = v20
	v3 = base.Simd_g_const(&F_Predictor12_SSE2__k0)
	v4 = int32(0)
	v5 = base.Simd_g_v128_load32_zero(m, l1, v4)
	v7 = base.Simd_g_const(&F_Predictor12_SSE2__k1)
	v10 = base.Simd_g_v128_load32_zero(m, l0, v4)
	v17 = base.Simd_g_v128_load32_zero(m, l1+int32(-4), v4)
	v20 = base.Simd_g_i16x8_sub(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v5, v3, base.Simd_g_const(&F_Predictor12_SSE2__k2), base.Simd_g_const(&F_Predictor12_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v10, v3, base.Simd_g_const(&F_Predictor12_SSE2__k2), base.Simd_g_const(&F_Predictor12_SSE2__k3))), base.Simd_g_i8x16_shuffle2(v17, v3, base.Simd_g_const(&F_Predictor12_SSE2__k2), base.Simd_g_const(&F_Predictor12_SSE2__k3)))
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_narrow_i16x8_u(v20, v20))
}

var F_Predictor12_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Predictor12_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Predictor12_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Predictor12_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_Predictor13_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	var v8 base.V128
	_ = v8
	var v11 base.V128
	_ = v11
	var v15 int32
	_ = v15
	var v16 base.V128
	_ = v16
	var v20 base.V128
	_ = v20
	var v22 base.V128
	_ = v22
	var v28 base.V128
	_ = v28
	v3 = base.Simd_g_const(&F_Predictor13_SSE2__k0)
	v5 = int32(0)
	v6 = base.Simd_g_v128_load32_zero(m, l1, v5)
	v8 = base.Simd_g_const(&F_Predictor13_SSE2__k1)
	v11 = base.Simd_g_v128_load32_zero(m, l0, v5)
	v15 = int32(1)
	v16 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v6, v3, base.Simd_g_const(&F_Predictor13_SSE2__k2), base.Simd_g_const(&F_Predictor13_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v11, v3, base.Simd_g_const(&F_Predictor13_SSE2__k2), base.Simd_g_const(&F_Predictor13_SSE2__k3))), v15)
	v20 = base.Simd_g_v128_load32_zero(m, l1+int32(-4), v5)
	v22 = base.Simd_g_i8x16_shuffle2(v20, v3, base.Simd_g_const(&F_Predictor13_SSE2__k2), base.Simd_g_const(&F_Predictor13_SSE2__k3))
	v28 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(base.Simd_g_i16x8_sub(v16, v22), base.Simd_g_i16x8_lt_s(v16, v22)), v15), v16)
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_narrow_i16x8_u(v28, v28))
}

var F_Predictor13_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Predictor13_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Predictor13_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Predictor13_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_Predictor5_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v3 base.V128
	_ = v3
	var v5 base.V128
	_ = v5
	var v7 base.V128
	_ = v7
	var v9 int32
	_ = v9
	var v10 base.V128
	_ = v10
	var v14 int32
	_ = v14
	var v17 base.V128
	_ = v17
	var v22 base.V128
	_ = v22
	v3 = base.Simd_g_const(&F_Predictor5_SSE2__k0)
	v5 = base.Simd_g_v128_load32_zero(m, l1, int32(4))
	v7 = base.Simd_g_const(&F_Predictor5_SSE2__k1)
	v9 = int32(0)
	v10 = base.Simd_g_v128_load32_zero(m, l0, v9)
	v14 = int32(1)
	v17 = base.Simd_g_v128_load32_zero(m, l1, v9)
	v22 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v5, v3, base.Simd_g_const(&F_Predictor5_SSE2__k2), base.Simd_g_const(&F_Predictor5_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v10, v3, base.Simd_g_const(&F_Predictor5_SSE2__k2), base.Simd_g_const(&F_Predictor5_SSE2__k3))), v14), base.Simd_g_i8x16_shuffle2(v17, v3, base.Simd_g_const(&F_Predictor5_SSE2__k2), base.Simd_g_const(&F_Predictor5_SSE2__k3))), v14)
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_narrow_i16x8_u(v22, v22))
}

var F_Predictor5_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_Predictor5_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_Predictor5_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_Predictor5_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_Predictor6_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	var v10 base.V128
	_ = v10
	v5 = int32(0)
	v6 = base.Simd_g_v128_load32_zero(m, l0, v5)
	v10 = base.Simd_g_v128_load32_zero(m, l1+int32(-4), v5)
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v6, v10), base.Simd_g_v128_and(base.Simd_g_v128_xor(v10, v6), base.Simd_g_const(&F_Predictor6_SSE2__k0))))
}

var F_Predictor6_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_Predictor7_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	var v8 base.V128
	_ = v8
	v5 = int32(0)
	v6 = base.Simd_g_v128_load32_zero(m, l0, v5)
	v8 = base.Simd_g_v128_load32_zero(m, l1, v5)
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v6, v8), base.Simd_g_v128_and(base.Simd_g_v128_xor(v8, v6), base.Simd_g_const(&F_Predictor7_SSE2__k0))))
}

var F_Predictor7_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_Predictor8_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v7 int32
	_ = v7
	var v8 base.V128
	_ = v8
	var v10 base.V128
	_ = v10
	v7 = int32(0)
	v8 = base.Simd_g_v128_load32_zero(m, l1+int32(-4), v7)
	v10 = base.Simd_g_v128_load32_zero(m, l1, v7)
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v8, v10), base.Simd_g_v128_and(base.Simd_g_v128_xor(v10, v8), base.Simd_g_const(&F_Predictor8_SSE2__k0))))
}

var F_Predictor8_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_Predictor9_SSE2(m *base.Module, l0 int32, l1 int32) int32 {
	var v5 int32
	_ = v5
	var v6 base.V128
	_ = v6
	var v8 base.V128
	_ = v8
	v5 = int32(0)
	v6 = base.Simd_g_v128_load32_zero(m, l1, v5)
	v8 = base.Simd_g_v128_load32_zero(m, l1, int32(4))
	return base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v6, v8), base.Simd_g_v128_and(base.Simd_g_v128_xor(v8, v6), base.Simd_g_const(&F_Predictor9_SSE2__k0))))
}

var F_Predictor9_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorAdd0_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v30 int32
	_ = v30
	var v35 int32
	_ = v35
	var v41 int32
	_ = v41
	var v47 int32
	_ = v47
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v62 int32
	_ = v62
	var v67 int32
	_ = v67
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v41 = int32(0)
			v47 = v41 << (uint(int32(2)) % 32)
			v55 = l2 - v41
			v56 = l0 + v47
			v57 = l3 + v47
			for {
				v58 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
				*(*int32)(unsafe.Add(mBase, uint32(v57))) = v58 + int32(-16777216)
				v62 = int32(4)
				v67 = v55 + int32(-1)
				if v67 != 0 {
					v55 = v67
					v56 = v56 + v62
					v57 = v57 + v62
					continue
				} else {
					break
				}
				break
			}
		} else {
			v15 = l2 & int32(2147483644)
			v21 = v15
			v22 = l3
			v23 = l0
			for {
				v24 = int32(0)
				v25 = base.Simd_g_v128_load(m, v23, v24)
				v27 = base.Simd_g_i32x4_add(v25, base.Simd_g_const(&F_PredictorAdd0_C__k0))
				base.Simd_g_v128_store(m, v22, v24, v27)
				v30 = int32(16)
				v35 = v21 + int32(-4)
				if v35 != 0 {
					v21 = v35
					v22 = v22 + v30
					v23 = v23 + v30
					continue
				} else {
					break
				}
				break
			}
			if v15 == l2 {
			} else {
				v41 = v15
				v47 = v41 << (uint(int32(2)) % 32)
				v55 = l2 - v41
				v56 = l0 + v47
				v57 = l3 + v47
				for {
					v58 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
					*(*int32)(unsafe.Add(mBase, uint32(v57))) = v58 + int32(-16777216)
					v62 = int32(4)
					v67 = v55 + int32(-1)
					if v67 != 0 {
						v55 = v67
						v56 = v56 + v62
						v57 = v57 + v62
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

var F_PredictorAdd0_C__k0 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}

func F_PredictorAdd0_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
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
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v50 int32
	_ = v50
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v59 int32
	_ = v59
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
	var v68 int32
	_ = v68
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v102 int32
	_ = v102
	var __phi102 int32
	_ = __phi102
	var v103 int32
	_ = v103
	var __phi103 int32
	_ = __phi103
	var v105 int32
	_ = v105
	var __phi105 int32
	_ = __phi105
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v122 int32
	_ = v122
	var v130 int32
	_ = v130
	var v138 int32
	_ = v138
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	v12 = int32(4)
	if l2 < v12 {
		v130 = int32(0)
	} else {
		v17 = l2 + int32(-4)
		v21 = int32(base.Ui32(v17)>>(uint(int32(2))%32)) + int32(1)
		v23 = v21 & int32(3)
		if base.Ui32(int32(12)) <= base.Ui32(v17) {
			v35 = int32(4)
			v36 = int32(0)
			v39 = v21 & int32(2147483644)
			for {
				v42 = l3 + v36
				v43 = l0 + v36
				v44 = int32(0)
				v45 = base.Simd_g_v128_load(m, v43, v44)
				v46 = base.Simd_g_const(&F_PredictorAdd0_SSE2__k0)
				v47 = base.Simd_g_i8x16_add(v45, v46)
				base.Simd_g_v128_store(m, v42, v44, v47)
				v50 = int32(16)
				v55 = base.Simd_g_v128_load(m, v43+v50, v44)
				v56 = base.Simd_g_i8x16_add(v55, v46)
				base.Simd_g_v128_store(m, v42+v50, v44, v56)
				v59 = int32(32)
				v64 = base.Simd_g_v128_load(m, v43+v59, v44)
				v65 = base.Simd_g_i8x16_add(v64, v46)
				base.Simd_g_v128_store(m, v42+v59, v44, v65)
				v68 = int32(48)
				v73 = base.Simd_g_v128_load(m, v43+v68, v44)
				v74 = base.Simd_g_i8x16_add(v73, v46)
				base.Simd_g_v128_store(m, v42+v68, v44, v74)
				v80 = v35 + v50
				v82 = v39 + int32(-4)
				if v82 != 0 {
					v35 = v80
					v36 = v36 + int32(64)
					v39 = v82
					continue
				} else {
					break
				}
				break
			}
			v89 = v80
			v90 = v35 + int32(12)
		} else {
			v89 = v12
			v90 = int32(0)
		}
		if v23 == int32(0) {
			v130 = v90
		} else {
			__phi102 = v89
			__phi103 = v90
			__phi105 = v23
			v102 = __phi102
			v103 = __phi103
			v105 = __phi105
			for {
				v110 = v103 << (uint(int32(2)) % 32)
				v113 = int32(0)
				v114 = base.Simd_g_v128_load(m, l0+v110, v113)
				v116 = base.Simd_g_i8x16_add(v114, base.Simd_g_const(&F_PredictorAdd0_SSE2__k0))
				base.Simd_g_v128_store(m, l3+v110, v113, v116)
				v122 = v105 + int32(-1)
				if v122 != 0 {
					__phi102 = v102 + int32(4)
					__phi103 = v102
					__phi105 = v122
					v102 = __phi102
					v103 = __phi103
					v105 = __phi105
					continue
				} else {
					break
				}
				break
			}
			v130 = v102 + int32(0)
		}
	}
	if l2 == v130 {
	} else {
		v138 = v130 << (uint(int32(2)) % 32)
		v143 = m.G34
		v144 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
		m.T0[v144].(func(*base.Module, int32, int32, int32, int32))(m, l0+v138, int32(0), l2-v130, l3+v138)
		mBase = m.M
	}
	return
}

var F_PredictorAdd0_SSE2__k0 = [2]uint64{0xff000000ff000000, 0xff000000ff000000}

func F_PredictorAdd10_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v32 int32
	_ = v32
	var v33 base.V128
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v43 int32
	_ = v43
	var v45 int32
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 int32
	_ = v47
	var v50 base.V128
	_ = v50
	var v53 base.V128
	_ = v53
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v67 base.V128
	_ = v67
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v71 base.V128
	_ = v71
	var v73 base.V128
	_ = v73
	var v77 base.V128
	_ = v77
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v88 base.V128
	_ = v88
	var v90 base.V128
	_ = v90
	var v94 base.V128
	_ = v94
	var v100 base.V128
	_ = v100
	var v101 base.V128
	_ = v101
	var v105 base.V128
	_ = v105
	var v107 base.V128
	_ = v107
	var v111 base.V128
	_ = v111
	var v118 base.V128
	_ = v118
	var v120 base.V128
	_ = v120
	var v123 int32
	_ = v123
	var v132 int32
	_ = v132
	var v138 int32
	_ = v138
	var v152 int32
	_ = v152
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	if l2 < int32(4) {
		v138 = int32(0)
	} else {
		v20 = int32(-4)
		v24 = int32(0)
		v25 = base.Simd_g_v128_load32_zero(m, l3+v20, v24)
		v32 = l1 + v20
		v33 = v25
		v34 = v24
		v35 = l3
		v36 = l0
		for {
			v43 = int32(4)
			v45 = int32(0)
			v46 = base.Simd_g_v128_load_rng(m, v32+v43, v45, int32(-4), int32(24))
			v47 = int32(8)
			v50 = base.Simd_g_v128_load_nc(m, v32+v47, v45)
			v53 = base.Simd_g_const(&F_PredictorAdd10_SSE2__k0)
			v55 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v46, v50), base.Simd_g_v128_and(base.Simd_g_v128_xor(v50, v46), v53))
			v57 = base.Simd_g_v128_load_nc(m, v32, v45)
			v61 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v33, v57), base.Simd_g_v128_and(base.Simd_g_v128_xor(v57, v33), v53))
			v67 = base.Simd_g_v128_load(m, v36, v45)
			v68 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v55, v61), base.Simd_g_v128_and(base.Simd_g_v128_xor(v55, v61), v53)), v67)
			v69 = base.Simd_g_const(&F_PredictorAdd10_SSE2__k1)
			v70 = base.Simd_g_const(&F_PredictorAdd10_SSE2__k2)
			v71 = base.Simd_g_i8x16_shuffle2(v55, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v73 = base.Simd_g_i8x16_shuffle2(v57, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v77 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v68, v73), base.Simd_g_v128_and(base.Simd_g_v128_xor(v68, v73), v53))
			v83 = base.Simd_g_i8x16_shuffle2(v67, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v84 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v71, v77), base.Simd_g_v128_and(base.Simd_g_v128_xor(v77, v71), v53)), v83)
			v88 = base.Simd_g_i8x16_shuffle2(v71, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v90 = base.Simd_g_i8x16_shuffle2(v73, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v94 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v84, v90), base.Simd_g_v128_and(base.Simd_g_v128_xor(v84, v90), v53))
			v100 = base.Simd_g_i8x16_shuffle2(v83, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v101 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v88, v94), base.Simd_g_v128_and(base.Simd_g_v128_xor(v94, v88), v53)), v100)
			v105 = base.Simd_g_i8x16_shuffle2(v88, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v107 = base.Simd_g_i8x16_shuffle2(v90, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			v111 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v101, v107), base.Simd_g_v128_and(base.Simd_g_v128_xor(v101, v107), v53))
			v118 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(v105, v111), base.Simd_g_i8x16_shuffle2(v100, v69, base.Simd_g_const(&F_PredictorAdd10_SSE2__k3), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))), base.Simd_g_v128_and(base.Simd_g_v128_xor(v111, v105), v53))
			v120 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v68, v84, base.Simd_g_const(&F_PredictorAdd10_SSE2__k5), base.Simd_g_const(&F_PredictorAdd10_SSE2__k6)), v101, base.Simd_g_const(&F_PredictorAdd10_SSE2__k7), base.Simd_g_const(&F_PredictorAdd10_SSE2__k8)), v118, base.Simd_g_const(&F_PredictorAdd10_SSE2__k9), base.Simd_g_const(&F_PredictorAdd10_SSE2__k4))
			base.Simd_g_v128_store(m, v35, v45, v120)
			v123 = int32(16)
			v132 = v34 + v43
			if v34+v47 <= l2 {
				v32 = v32 + v123
				v33 = v118
				v34 = v132
				v35 = v35 + v123
				v36 = v36 + v123
				continue
			} else {
				break
			}
			break
		}
		v138 = v132
	}
	if l2 == v138 {
	} else {
		v152 = v138 << (uint(int32(2)) % 32)
		v157 = m.G34
		v158 = *(*int32)(unsafe.Add(mBase, uint32(v157)+40))
		m.T0[v158].(func(*base.Module, int32, int32, int32, int32))(m, l0+v152, l1+v152, l2-v138, l3+v152)
		mBase = m.M
	}
	return
}

var F_PredictorAdd10_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_PredictorAdd10_SSE2__k1 = [2]uint64{0x0, 0x0}
var F_PredictorAdd10_SSE2__k2 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_PredictorAdd10_SSE2__k3 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_PredictorAdd10_SSE2__k4 = [2]uint64{0x8080808080808080, 0x302010080808080}
var F_PredictorAdd10_SSE2__k5 = [2]uint64{0x8080808003020100, 0x302010003020100}
var F_PredictorAdd10_SSE2__k6 = [2]uint64{0x302010080808080, 0x8080808080808080}
var F_PredictorAdd10_SSE2__k7 = [2]uint64{0x706050403020100, 0x302010080808080}
var F_PredictorAdd10_SSE2__k8 = [2]uint64{0x8080808080808080, 0x8080808003020100}
var F_PredictorAdd10_SSE2__k9 = [2]uint64{0x706050403020100, 0x808080800b0a0908}

func F_PredictorAdd11_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v28 base.V128
	_ = v28
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v49 int32
	_ = v49
	var v50 base.V128
	_ = v50
	var v51 int32
	_ = v51
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v57 base.V128
	_ = v57
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v63 base.V128
	_ = v63
	var v64 int32
	_ = v64
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v69 int32
	_ = v69
	var v71 base.V128
	_ = v71
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v78 base.V128
	_ = v78
	var v81 base.V128
	_ = v81
	var v85 base.V128
	_ = v85
	var v88 base.V128
	_ = v88
	var v94 base.V128
	_ = v94
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v100 base.V128
	_ = v100
	var v104 base.V128
	_ = v104
	var v107 base.V128
	_ = v107
	var v113 base.V128
	_ = v113
	var v116 base.V128
	_ = v116
	var v117 base.V128
	_ = v117
	var v118 base.V128
	_ = v118
	var v119 base.V128
	_ = v119
	var v121 base.V128
	_ = v121
	var v123 base.V128
	_ = v123
	var v125 base.V128
	_ = v125
	var v127 base.V128
	_ = v127
	var v130 base.V128
	_ = v130
	var v134 base.V128
	_ = v134
	var v137 base.V128
	_ = v137
	var v144 base.V128
	_ = v144
	var v147 base.V128
	_ = v147
	var v151 base.V128
	_ = v151
	var v153 base.V128
	_ = v153
	var v155 base.V128
	_ = v155
	var v157 base.V128
	_ = v157
	var v159 base.V128
	_ = v159
	var v162 base.V128
	_ = v162
	var v166 base.V128
	_ = v166
	var v169 base.V128
	_ = v169
	var v176 base.V128
	_ = v176
	var v179 base.V128
	_ = v179
	var v185 base.V128
	_ = v185
	var v189 base.V128
	_ = v189
	var v191 base.V128
	_ = v191
	var v194 base.V128
	_ = v194
	var v198 base.V128
	_ = v198
	var v201 base.V128
	_ = v201
	var v211 base.V128
	_ = v211
	var v213 base.V128
	_ = v213
	var v225 int32
	_ = v225
	var v231 int32
	_ = v231
	var v248 int32
	_ = v248
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	if l2 < int32(4) {
		v231 = int32(0)
	} else {
		v23 = int32(-4)
		v27 = int32(0)
		v28 = base.Simd_g_v128_load32_zero(m, l3+v23, v27)
		v35 = l1 + v23
		v36 = v28
		v37 = v27
		v38 = l3
		v39 = l0
		for {
			v49 = int32(0)
			v50 = base.Simd_g_v128_load(m, v39, v49)
			v51 = int32(4)
			v54 = base.Simd_g_v128_load_rng(m, v35+v51, v49, int32(-4), int32(20))
			v56 = base.Simd_g_v128_load_nc(m, v35, v49)
			v57 = base.Simd_g_const(&F_PredictorAdd11_SSE2__k0)
			v58 = base.Simd_g_i8x16_shuffle2(v56, v54, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v60 = base.Simd_g_i8x16_shuffle2(v36, v54, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v63 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v58, v60), base.Simd_g_i8x16_sub_sat_u(v60, v58))
			v64 = int32(8)
			v66 = base.Simd_g_const(&F_PredictorAdd11_SSE2__k3)
			v68 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v63, v64), base.Simd_g_v128_and(v63, v66))
			v69 = int32(16)
			v71 = base.Simd_g_i16x8_add(v68, base.Simd_g_i32x4_shl(v68, v69))
			v72 = int32(32)
			v75 = int32(48)
			v78 = base.Simd_g_i8x16_shuffle2(v54, v54, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v81 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v58, v78), base.Simd_g_i8x16_sub_sat_u(v78, v58))
			v85 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v81, v64), base.Simd_g_v128_and(v81, v66))
			v88 = base.Simd_g_i16x8_add(v85, base.Simd_g_i32x4_shl(v85, v69))
			v94 = base.Simd_g_const(&F_PredictorAdd11_SSE2__k4)
			v95 = base.Simd_g_i8x16_shuffle2(v56, v54, base.Simd_g_const(&F_PredictorAdd11_SSE2__k5), base.Simd_g_const(&F_PredictorAdd11_SSE2__k6))
			v97 = base.Simd_g_i8x16_shuffle2(v54, v54, base.Simd_g_const(&F_PredictorAdd11_SSE2__k5), base.Simd_g_const(&F_PredictorAdd11_SSE2__k6))
			v100 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v95, v97), base.Simd_g_i8x16_sub_sat_u(v97, v95))
			v104 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v100, v64), base.Simd_g_v128_and(v100, v66))
			v107 = base.Simd_g_i16x8_add(v104, base.Simd_g_i32x4_shl(v104, v69))
			v113 = base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v88, base.Simd_g_i64x2_shl(v88, v72)), v75), base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v107, base.Simd_g_i64x2_shl(v107, v72)), v75))
			v116 = base.Simd_g_i8x16_add(v50, base.Simd_g_v128_bitselect(v36, v54, base.Simd_g_i32x4_gt_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v71, base.Simd_g_i64x2_shl(v71, v72)), v75), v113)))
			v117 = base.Simd_g_const(&F_PredictorAdd11_SSE2__k7)
			v118 = base.Simd_g_const(&F_PredictorAdd11_SSE2__k8)
			v119 = base.Simd_g_i8x16_shuffle2(v50, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v121 = base.Simd_g_i8x16_shuffle2(v54, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v123 = base.Simd_g_i8x16_shuffle2(v56, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v125 = base.Simd_g_i8x16_shuffle2(v123, v121, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v127 = base.Simd_g_i8x16_shuffle2(v116, v121, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v130 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v125, v127), base.Simd_g_i8x16_sub_sat_u(v127, v125))
			v134 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v130, v64), base.Simd_g_v128_and(v130, v66))
			v137 = base.Simd_g_i16x8_add(v134, base.Simd_g_i32x4_shl(v134, v69))
			v144 = base.Simd_g_i8x16_shuffle2(v113, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v147 = base.Simd_g_i8x16_add(v119, base.Simd_g_v128_bitselect(v116, v121, base.Simd_g_i32x4_gt_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v137, base.Simd_g_i64x2_shl(v137, v72)), v75), v144)))
			v151 = base.Simd_g_i8x16_shuffle2(v119, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v153 = base.Simd_g_i8x16_shuffle2(v121, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v155 = base.Simd_g_i8x16_shuffle2(v123, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v157 = base.Simd_g_i8x16_shuffle2(v155, v153, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v159 = base.Simd_g_i8x16_shuffle2(v147, v153, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v162 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v157, v159), base.Simd_g_i8x16_sub_sat_u(v159, v157))
			v166 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v162, v64), base.Simd_g_v128_and(v162, v66))
			v169 = base.Simd_g_i16x8_add(v166, base.Simd_g_i32x4_shl(v166, v69))
			v176 = base.Simd_g_i8x16_shuffle2(v144, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v179 = base.Simd_g_i8x16_add(v151, base.Simd_g_v128_bitselect(v147, v153, base.Simd_g_i32x4_gt_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v169, base.Simd_g_i64x2_shl(v169, v72)), v75), v176)))
			v185 = base.Simd_g_i8x16_shuffle2(v153, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			v189 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v155, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10)), v185, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v191 = base.Simd_g_i8x16_shuffle2(v179, v185, base.Simd_g_const(&F_PredictorAdd11_SSE2__k1), base.Simd_g_const(&F_PredictorAdd11_SSE2__k2))
			v194 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v189, v191), base.Simd_g_i8x16_sub_sat_u(v191, v189))
			v198 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v194, v64), base.Simd_g_v128_and(v194, v66))
			v201 = base.Simd_g_i16x8_add(v198, base.Simd_g_i32x4_shl(v198, v69))
			v211 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v151, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10)), base.Simd_g_v128_bitselect(v179, v185, base.Simd_g_i32x4_gt_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v201, base.Simd_g_i64x2_shl(v201, v72)), v75), base.Simd_g_i8x16_shuffle2(v176, v117, base.Simd_g_const(&F_PredictorAdd11_SSE2__k9), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10)))))
			v213 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v116, v147, base.Simd_g_const(&F_PredictorAdd11_SSE2__k11), base.Simd_g_const(&F_PredictorAdd11_SSE2__k12)), v179, base.Simd_g_const(&F_PredictorAdd11_SSE2__k13), base.Simd_g_const(&F_PredictorAdd11_SSE2__k14)), v211, base.Simd_g_const(&F_PredictorAdd11_SSE2__k15), base.Simd_g_const(&F_PredictorAdd11_SSE2__k10))
			base.Simd_g_v128_store(m, v38, v49, v213)
			v225 = v37 + v51
			if v37+v64 <= l2 {
				v35 = v35 + v69
				v36 = v211
				v37 = v225
				v38 = v38 + v69
				v39 = v39 + v69
				continue
			} else {
				break
			}
			break
		}
		v231 = v225
	}
	if l2 == v231 {
	} else {
		v248 = v231 << (uint(int32(2)) % 32)
		v253 = m.G34
		v254 = *(*int32)(unsafe.Add(mBase, uint32(v253)+44))
		m.T0[v254].(func(*base.Module, int32, int32, int32, int32))(m, l0+v248, l1+v248, l2-v231, l3+v248)
		mBase = m.M
	}
	return
}

var F_PredictorAdd11_SSE2__k0 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_PredictorAdd11_SSE2__k1 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_PredictorAdd11_SSE2__k2 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_PredictorAdd11_SSE2__k3 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_PredictorAdd11_SSE2__k4 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_PredictorAdd11_SSE2__k5 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_PredictorAdd11_SSE2__k6 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}
var F_PredictorAdd11_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_PredictorAdd11_SSE2__k8 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_PredictorAdd11_SSE2__k9 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_PredictorAdd11_SSE2__k10 = [2]uint64{0x8080808080808080, 0x302010080808080}
var F_PredictorAdd11_SSE2__k11 = [2]uint64{0x8080808003020100, 0x302010003020100}
var F_PredictorAdd11_SSE2__k12 = [2]uint64{0x302010080808080, 0x8080808080808080}
var F_PredictorAdd11_SSE2__k13 = [2]uint64{0x706050403020100, 0x302010080808080}
var F_PredictorAdd11_SSE2__k14 = [2]uint64{0x8080808080808080, 0x8080808003020100}
var F_PredictorAdd11_SSE2__k15 = [2]uint64{0x706050403020100, 0x808080800b0a0908}

func F_PredictorAdd12_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v26 base.V128
	_ = v26
	var v35 int32
	_ = v35
	var v37 base.V128
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v53 base.V128
	_ = v53
	var v56 base.V128
	_ = v56
	var v57 base.V128
	_ = v57
	var v60 base.V128
	_ = v60
	var v61 base.V128
	_ = v61
	var v64 base.V128
	_ = v64
	var v66 base.V128
	_ = v66
	var v68 base.V128
	_ = v68
	var v69 base.V128
	_ = v69
	var v70 base.V128
	_ = v70
	var v73 base.V128
	_ = v73
	var v77 base.V128
	_ = v77
	var v80 base.V128
	_ = v80
	var v83 base.V128
	_ = v83
	var v84 base.V128
	_ = v84
	var v91 base.V128
	_ = v91
	var v95 base.V128
	_ = v95
	var v97 base.V128
	_ = v97
	var v102 int32
	_ = v102
	var v111 int32
	_ = v111
	var v117 int32
	_ = v117
	var v131 int32
	_ = v131
	var v136 int32
	_ = v136
	var v137 int32
	_ = v137
	if l2 < int32(4) {
		v117 = int32(0)
	} else {
		v20 = int32(-4)
		v24 = int32(0)
		v25 = base.Simd_g_v128_load32_zero(m, l3+v20, v24)
		v26 = base.Simd_g_const(&F_PredictorAdd12_SSE2__k0)
		v35 = l1 + v20
		v37 = base.Simd_g_i8x16_shuffle2(v25, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2))
		v38 = v24
		v39 = l3
		v40 = l0
		for {
			v46 = int32(4)
			v48 = int32(0)
			v49 = base.Simd_g_v128_load_rng(m, v35+v46, v48, int32(-4), int32(20))
			v50 = base.Simd_g_const(&F_PredictorAdd12_SSE2__k3)
			v53 = base.Simd_g_v128_load_nc(m, v35, v48)
			v56 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v49, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2)), base.Simd_g_i8x16_shuffle2(v53, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2)))
			v57 = base.Simd_g_i16x8_add(v56, v37)
			v60 = base.Simd_g_v128_load(m, v40, v48)
			v61 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(v57, v57), v60)
			v64 = base.Simd_g_const(&F_PredictorAdd12_SSE2__k4)
			v66 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v61, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2)), base.Simd_g_i8x16_shuffle2(v56, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k5), base.Simd_g_const(&F_PredictorAdd12_SSE2__k6)))
			v68 = base.Simd_g_const(&F_PredictorAdd12_SSE2__k7)
			v69 = base.Simd_g_i8x16_shuffle2(v60, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k8), base.Simd_g_const(&F_PredictorAdd12_SSE2__k9))
			v70 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(v66, v66), v69)
			v73 = base.Simd_g_const(&F_PredictorAdd12_SSE2__k10)
			v77 = base.Simd_g_i16x8_sub(base.Simd_g_i8x16_shuffle2(v49, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k11), base.Simd_g_const(&F_PredictorAdd12_SSE2__k12)), base.Simd_g_i8x16_shuffle2(v53, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k11), base.Simd_g_const(&F_PredictorAdd12_SSE2__k12)))
			v80 = base.Simd_g_i16x8_add(v77, base.Simd_g_i8x16_shuffle2(v70, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2)))
			v83 = base.Simd_g_i8x16_shuffle2(v69, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k8), base.Simd_g_const(&F_PredictorAdd12_SSE2__k9))
			v84 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(v80, v80), v83)
			v91 = base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v84, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2)), base.Simd_g_i8x16_shuffle2(v77, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k5), base.Simd_g_const(&F_PredictorAdd12_SSE2__k6)))
			v95 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_narrow_i16x8_u(v91, v91), base.Simd_g_i8x16_shuffle2(v83, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k8), base.Simd_g_const(&F_PredictorAdd12_SSE2__k9)))
			v97 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(base.Simd_g_i8x16_shuffle2(v61, v70, base.Simd_g_const(&F_PredictorAdd12_SSE2__k13), base.Simd_g_const(&F_PredictorAdd12_SSE2__k14)), v84, base.Simd_g_const(&F_PredictorAdd12_SSE2__k15), base.Simd_g_const(&F_PredictorAdd12_SSE2__k16)), v95, base.Simd_g_const(&F_PredictorAdd12_SSE2__k17), base.Simd_g_const(&F_PredictorAdd12_SSE2__k9))
			base.Simd_g_v128_store(m, v39, v48, v97)
			v102 = int32(16)
			v111 = v38 + v46
			if v38+int32(8) <= l2 {
				v35 = v35 + v102
				v37 = base.Simd_g_i8x16_shuffle2(v95, v26, base.Simd_g_const(&F_PredictorAdd12_SSE2__k1), base.Simd_g_const(&F_PredictorAdd12_SSE2__k2))
				v38 = v111
				v39 = v39 + v102
				v40 = v40 + v102
				continue
			} else {
				break
			}
			break
		}
		v117 = v111
	}
	if l2 == v117 {
	} else {
		v131 = v117 << (uint(int32(2)) % 32)
		v136 = m.G34
		v137 = *(*int32)(unsafe.Add(mBase, uint32(v136)+48))
		m.T0[v137].(func(*base.Module, int32, int32, int32, int32))(m, l0+v131, l1+v131, l2-v117, l3+v131)
		mBase = m.M
	}
	return
}

var F_PredictorAdd12_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorAdd12_SSE2__k1 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_PredictorAdd12_SSE2__k2 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_PredictorAdd12_SSE2__k3 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_PredictorAdd12_SSE2__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x1716151413121110}
var F_PredictorAdd12_SSE2__k5 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_PredictorAdd12_SSE2__k6 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_PredictorAdd12_SSE2__k7 = [2]uint64{0xb0a090807060504, 0x131211100f0e0d0c}
var F_PredictorAdd12_SSE2__k8 = [2]uint64{0xb0a090807060504, 0x808080800f0e0d0c}
var F_PredictorAdd12_SSE2__k9 = [2]uint64{0x8080808080808080, 0x302010080808080}
var F_PredictorAdd12_SSE2__k10 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_PredictorAdd12_SSE2__k11 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_PredictorAdd12_SSE2__k12 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_PredictorAdd12_SSE2__k13 = [2]uint64{0x8080808003020100, 0x302010003020100}
var F_PredictorAdd12_SSE2__k14 = [2]uint64{0x302010080808080, 0x8080808080808080}
var F_PredictorAdd12_SSE2__k15 = [2]uint64{0x706050403020100, 0x302010080808080}
var F_PredictorAdd12_SSE2__k16 = [2]uint64{0x8080808080808080, 0x8080808003020100}
var F_PredictorAdd12_SSE2__k17 = [2]uint64{0x706050403020100, 0x808080800b0a0908}

func F_PredictorAdd13_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
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
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 base.V128
	_ = v27
	var v28 base.V128
	_ = v28
	var v29 base.V128
	_ = v29
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v42 base.V128
	_ = v42
	var v48 base.V128
	_ = v48
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
			v24 = int32(4)
			v25 = v17 + v24
			v26 = int32(0)
			v27 = base.Simd_g_v128_load32_zero(m, v25, v26)
			v28 = base.Simd_g_const(&F_PredictorAdd13_SSE2__k0)
			v29 = base.Simd_g_const(&F_PredictorAdd13_SSE2__k1)
			v37 = int32(1)
			v38 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v27, v28, base.Simd_g_const(&F_PredictorAdd13_SSE2__k2), base.Simd_g_const(&F_PredictorAdd13_SSE2__k3)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_replace_lane_l0(v28, v20), v28, base.Simd_g_const(&F_PredictorAdd13_SSE2__k2), base.Simd_g_const(&F_PredictorAdd13_SSE2__k3))), v37)
			v40 = base.Simd_g_v128_load32_zero(m, v17, v26)
			v42 = base.Simd_g_i8x16_shuffle2(v40, v28, base.Simd_g_const(&F_PredictorAdd13_SSE2__k2), base.Simd_g_const(&F_PredictorAdd13_SSE2__k3))
			v48 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(base.Simd_g_i16x8_sub(v38, v42), base.Simd_g_i16x8_lt_s(v38, v42)), v37), v38)
			v51 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_narrow_i16x8_u(v48, v48))
			v52 = int32(-16711936)
			v54 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
			v60 = int32(16711935)
			v67 = (v51&v52+v54&v52)&v52 | (v51&v60+v54&v60)&v60
			*(*int32)(unsafe.Add(mBase, uint32(v19))) = v67
			v74 = v18 + int32(-1)
			if v74 != 0 {
				v16 = v16 + v24
				v17 = v25
				v18 = v74
				v19 = v19 + v24
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

var F_PredictorAdd13_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorAdd13_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_PredictorAdd13_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_PredictorAdd13_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_PredictorAdd1_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v20 base.V128
	_ = v20
	var v22 int32
	_ = v22
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 base.V128
	_ = v39
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v45 base.V128
	_ = v45
	var v46 int32
	_ = v46
	var v47 base.V128
	_ = v47
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v52 base.V128
	_ = v52
	var v54 base.V128
	_ = v54
	var v57 int32
	_ = v57
	var v62 base.V128
	_ = v62
	var v65 base.V128
	_ = v65
	var v66 base.V128
	_ = v66
	var v71 base.V128
	_ = v71
	var v75 base.V128
	_ = v75
	var v76 int32
	_ = v76
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 base.V128
	_ = v92
	var v101 int32
	_ = v101
	var v103 base.V128
	_ = v103
	var v105 int32
	_ = v105
	var v106 base.V128
	_ = v106
	var v109 base.V128
	_ = v109
	var v113 base.V128
	_ = v113
	var v121 int32
	_ = v121
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	v13 = int32(4)
	if l2 < v13 {
		v121 = int32(0)
	} else {
		v17 = int32(-4)
		v20 = base.Simd_g_v128_load32_splat(m, l3+v17, int32(0))
		v22 = l2 + v17
		if base.Ui32(int32(4)) <= base.Ui32(v22) {
			v37 = int32(4)
			v38 = l0
			v39 = v20
			v41 = (int32(base.Ui32(v22)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			v42 = l3
			for {
				v45 = base.Simd_g_const(&F_PredictorAdd1_SSE2__k0)
				v46 = int32(0)
				v47 = base.Simd_g_v128_load(m, v38, v46)
				v48 = base.Simd_g_const(&F_PredictorAdd1_SSE2__k1)
				v50 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v45, v47, base.Simd_g_const(&F_PredictorAdd1_SSE2__k2), base.Simd_g_const(&F_PredictorAdd1_SSE2__k3)), v47)
				v52 = base.Simd_g_const(&F_PredictorAdd1_SSE2__k4)
				v54 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_add(v50, v39), base.Simd_g_i8x16_shuffle2(v45, v50, base.Simd_g_const(&F_PredictorAdd1_SSE2__k5), base.Simd_g_const(&F_PredictorAdd1_SSE2__k6)))
				base.Simd_g_v128_store(m, v42, v46, v54)
				v57 = int32(16)
				v62 = base.Simd_g_v128_load(m, v38+v57, v46)
				v65 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v45, v62, base.Simd_g_const(&F_PredictorAdd1_SSE2__k2), base.Simd_g_const(&F_PredictorAdd1_SSE2__k3)), v62)
				v66 = base.Simd_g_const(&F_PredictorAdd1_SSE2__k7)
				v71 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_add(v65, base.Simd_g_i8x16_swizzle_c(v54, base.Simd_g_const(&F_PredictorAdd1_SSE2__k7))), base.Simd_g_i8x16_shuffle2(v45, v65, base.Simd_g_const(&F_PredictorAdd1_SSE2__k5), base.Simd_g_const(&F_PredictorAdd1_SSE2__k6)))
				base.Simd_g_v128_store(m, v42+v57, v46, v71)
				v75 = base.Simd_g_i8x16_swizzle_c(v71, base.Simd_g_const(&F_PredictorAdd1_SSE2__k7))
				v76 = int32(32)
				v81 = v37 + int32(8)
				v83 = v41 + int32(-2)
				if v83 != 0 {
					v37 = v81
					v38 = v38 + v76
					v39 = v75
					v41 = v83
					v42 = v42 + v76
					continue
				} else {
					break
				}
				break
			}
			v90 = v81
			v91 = v37 + int32(4)
			v92 = v75
		} else {
			v90 = v13
			v91 = int32(0)
			v92 = v20
		}
		if v22&int32(4) != 0 {
			v121 = v91
		} else {
			v101 = v91 << (uint(int32(2)) % 32)
			v103 = base.Simd_g_const(&F_PredictorAdd1_SSE2__k0)
			v105 = int32(0)
			v106 = base.Simd_g_v128_load(m, l0+v101, v105)
			v109 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_shuffle2(v103, v106, base.Simd_g_const(&F_PredictorAdd1_SSE2__k2), base.Simd_g_const(&F_PredictorAdd1_SSE2__k3)), v106)
			v113 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_add(v109, v92), base.Simd_g_i8x16_shuffle2(v103, v109, base.Simd_g_const(&F_PredictorAdd1_SSE2__k5), base.Simd_g_const(&F_PredictorAdd1_SSE2__k6)))
			base.Simd_g_v128_store(m, l3+v101, v105, v113)
			v121 = v90
		}
	}
	if l2 == v121 {
	} else {
		v130 = v121 << (uint(int32(2)) % 32)
		v135 = m.G34
		v136 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
		m.T0[v136].(func(*base.Module, int32, int32, int32, int32))(m, l0+v130, l1+v130, l2-v121, l3+v130)
		mBase = m.M
	}
	return
}

var F_PredictorAdd1_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorAdd1_SSE2__k1 = [2]uint64{0x131211100f0e0d0c, 0x1b1a191817161514}
var F_PredictorAdd1_SSE2__k2 = [2]uint64{0x808080800f0e0d0c, 0x8080808080808080}
var F_PredictorAdd1_SSE2__k3 = [2]uint64{0x302010080808080, 0xb0a090807060504}
var F_PredictorAdd1_SSE2__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x1716151413121110}
var F_PredictorAdd1_SSE2__k5 = [2]uint64{0xf0e0d0c0b0a0908, 0x8080808080808080}
var F_PredictorAdd1_SSE2__k6 = [2]uint64{0x8080808080808080, 0x706050403020100}
var F_PredictorAdd1_SSE2__k7 = [2]uint64{0xf0e0d0c0f0e0d0c, 0xf0e0d0c0f0e0d0c}

func F_PredictorAdd2_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
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
	var v33 base.V128
	_ = v33
	var v34 base.V128
	_ = v34
	var v37 base.V128
	_ = v37
	var v41 base.V128
	_ = v41
	var v46 base.V128
	_ = v46
	var v49 int32
	_ = v49
	var v56 int32
	_ = v56
	var v62 int32
	_ = v62
	var v72 int32
	_ = v72
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v97 int32
	_ = v97
	var v106 int32
	_ = v106
	var v113 int32
	_ = v113
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v62 = int32(0)
			v72 = v62 << (uint(int32(2)) % 32)
			v81 = l2 - v62
			v82 = l3 + v72
			v83 = l1 + v72
			v84 = l0 + v72
			for {
				v88 = *(*int32)(unsafe.Add(mBase, uint32(v84)))
				v89 = int32(-16711936)
				v91 = *(*int32)(unsafe.Add(mBase, uint32(v83)))
				v97 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(v82))) = (v88&v89+v91&v89)&v89 | (v88&v97+v91&v97)&v97
				v106 = int32(4)
				v113 = v81 + int32(-1)
				if v113 != 0 {
					v81 = v113
					v82 = v82 + v106
					v83 = v83 + v106
					v84 = v84 + v106
					continue
				} else {
					break
				}
				break
			}
		} else {
			v19 = l2 & int32(2147483644)
			v25 = v19
			v26 = l0
			v27 = l1
			v28 = l3
			for {
				v32 = int32(0)
				v33 = base.Simd_g_v128_load(m, v26, v32)
				v34 = base.Simd_g_const(&F_PredictorAdd2_C__k0)
				v37 = base.Simd_g_v128_load(m, v27, v32)
				v41 = base.Simd_g_const(&F_PredictorAdd2_C__k1)
				v46 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v33, v34), base.Simd_g_v128_and(v37, v34)), v34), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v33, v41), base.Simd_g_v128_and(v37, v41)), v41))
				base.Simd_g_v128_store(m, v28, v32, v46)
				v49 = int32(16)
				v56 = v25 + int32(-4)
				if v56 != 0 {
					v25 = v56
					v26 = v26 + v49
					v27 = v27 + v49
					v28 = v28 + v49
					continue
				} else {
					break
				}
				break
			}
			if v19 == l2 {
			} else {
				v62 = v19
				v72 = v62 << (uint(int32(2)) % 32)
				v81 = l2 - v62
				v82 = l3 + v72
				v83 = l1 + v72
				v84 = l0 + v72
				for {
					v88 = *(*int32)(unsafe.Add(mBase, uint32(v84)))
					v89 = int32(-16711936)
					v91 = *(*int32)(unsafe.Add(mBase, uint32(v83)))
					v97 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v82))) = (v88&v89+v91&v89)&v89 | (v88&v97+v91&v97)&v97
					v106 = int32(4)
					v113 = v81 + int32(-1)
					if v113 != 0 {
						v81 = v113
						v82 = v82 + v106
						v83 = v83 + v106
						v84 = v84 + v106
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

var F_PredictorAdd2_C__k0 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_PredictorAdd2_C__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_PredictorAdd2_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v16 int32
	_ = v16
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v44 int32
	_ = v44
	var v49 base.V128
	_ = v49
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v57 int32
	_ = v57
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v82 int32
	_ = v82
	var v85 int32
	_ = v85
	var v86 base.V128
	_ = v86
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	v11 = int32(4)
	if l2 < v11 {
		v98 = int32(0)
	} else {
		v16 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v16) {
			v31 = int32(4)
			v32 = l0
			v34 = (int32(base.Ui32(v16)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			v35 = l1
			v36 = l3
			for {
				v37 = int32(0)
				v38 = base.Simd_g_v128_load(m, v35, v37)
				v40 = base.Simd_g_v128_load(m, v32, v37)
				v41 = base.Simd_g_i8x16_add(v38, v40)
				base.Simd_g_v128_store(m, v36, v37, v41)
				v44 = int32(16)
				v49 = base.Simd_g_v128_load(m, v35+v44, v37)
				v53 = base.Simd_g_v128_load(m, v32+v44, v37)
				v54 = base.Simd_g_i8x16_add(v49, v53)
				base.Simd_g_v128_store(m, v36+v44, v37, v54)
				v57 = int32(32)
				v64 = v31 + int32(8)
				v66 = v34 + int32(-2)
				if v66 != 0 {
					v31 = v64
					v32 = v32 + v57
					v34 = v66
					v35 = v35 + v57
					v36 = v36 + v57
					continue
				} else {
					break
				}
				break
			}
			v73 = v64
			v74 = v31 + int32(4)
		} else {
			v73 = v11
			v74 = int32(0)
		}
		if v16&int32(4) != 0 {
			v98 = v74
		} else {
			v82 = v74 << (uint(int32(2)) % 32)
			v85 = int32(0)
			v86 = base.Simd_g_v128_load(m, l1+v82, v85)
			v89 = base.Simd_g_v128_load(m, l0+v82, v85)
			v90 = base.Simd_g_i8x16_add(v86, v89)
			base.Simd_g_v128_store(m, l3+v82, v85, v90)
			v98 = v73
		}
	}
	if l2 == v98 {
	} else {
		v105 = v98 << (uint(int32(2)) % 32)
		v110 = m.G34
		v111 = *(*int32)(unsafe.Add(mBase, uint32(v110)+8))
		m.T0[v111].(func(*base.Module, int32, int32, int32, int32))(m, l0+v105, l1+v105, l2-v98, l3+v105)
		mBase = m.M
	}
	return
}
func F_PredictorAdd3_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v21 int32
	_ = v21
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v36 base.V128
	_ = v36
	var v39 base.V128
	_ = v39
	var v43 base.V128
	_ = v43
	var v48 base.V128
	_ = v48
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	var v64 int32
	_ = v64
	var v74 int32
	_ = v74
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v101 int32
	_ = v101
	var v110 int32
	_ = v110
	var v117 int32
	_ = v117
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v64 = int32(0)
			v74 = v64 << (uint(int32(2)) % 32)
			v85 = l3 + v74
			v86 = l2 - v64
			v87 = l0 + v74
			v88 = v74 + l1 + int32(4)
			for {
				v92 = *(*int32)(unsafe.Add(mBase, uint32(v87)))
				v93 = int32(-16711936)
				v95 = *(*int32)(unsafe.Add(mBase, uint32(v88)))
				v101 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(v85))) = (v92&v93+v95&v93)&v93 | (v92&v101+v95&v101)&v101
				v110 = int32(4)
				v117 = v86 + int32(-1)
				if v117 != 0 {
					v85 = v85 + v110
					v86 = v117
					v87 = v87 + v110
					v88 = v88 + v110
					continue
				} else {
					break
				}
				break
			}
		} else {
			v21 = l2 & int32(2147483644)
			v27 = l1 + int32(4)
			v28 = v21
			v29 = l0
			v30 = l3
			for {
				v34 = int32(0)
				v35 = base.Simd_g_v128_load(m, v29, v34)
				v36 = base.Simd_g_const(&F_PredictorAdd3_C__k0)
				v39 = base.Simd_g_v128_load(m, v27, v34)
				v43 = base.Simd_g_const(&F_PredictorAdd3_C__k1)
				v48 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v35, v36), base.Simd_g_v128_and(v39, v36)), v36), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v35, v43), base.Simd_g_v128_and(v39, v43)), v43))
				base.Simd_g_v128_store(m, v30, v34, v48)
				v51 = int32(16)
				v58 = v28 + int32(-4)
				if v58 != 0 {
					v27 = v27 + v51
					v28 = v58
					v29 = v29 + v51
					v30 = v30 + v51
					continue
				} else {
					break
				}
				break
			}
			if v21 == l2 {
			} else {
				v64 = v21
				v74 = v64 << (uint(int32(2)) % 32)
				v85 = l3 + v74
				v86 = l2 - v64
				v87 = l0 + v74
				v88 = v74 + l1 + int32(4)
				for {
					v92 = *(*int32)(unsafe.Add(mBase, uint32(v87)))
					v93 = int32(-16711936)
					v95 = *(*int32)(unsafe.Add(mBase, uint32(v88)))
					v101 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v85))) = (v92&v93+v95&v93)&v93 | (v92&v101+v95&v101)&v101
					v110 = int32(4)
					v117 = v86 + int32(-1)
					if v117 != 0 {
						v85 = v85 + v110
						v86 = v117
						v87 = v87 + v110
						v88 = v88 + v110
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

var F_PredictorAdd3_C__k0 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_PredictorAdd3_C__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_PredictorAdd3_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 int32
	_ = v46
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v52 int32
	_ = v52
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v87 int32
	_ = v87
	var v92 int32
	_ = v92
	var v93 base.V128
	_ = v93
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v105 int32
	_ = v105
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	v12 = int32(4)
	if l2 < v12 {
		v105 = int32(0)
	} else {
		v17 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v17) {
			v33 = int32(4)
			v34 = int32(0)
			v36 = (int32(base.Ui32(v17)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v40 = l3 + v34
				v41 = l1 + v34
				v44 = int32(0)
				v45 = base.Simd_g_v128_load(m, v41+int32(4), v44)
				v46 = l0 + v34
				v48 = base.Simd_g_v128_load(m, v46, v44)
				v49 = base.Simd_g_i8x16_add(v45, v48)
				base.Simd_g_v128_store(m, v40, v44, v49)
				v52 = int32(16)
				v57 = base.Simd_g_v128_load(m, v41+int32(20), v44)
				v61 = base.Simd_g_v128_load(m, v46+v52, v44)
				v62 = base.Simd_g_i8x16_add(v57, v61)
				base.Simd_g_v128_store(m, v40+v52, v44, v62)
				v68 = v33 + int32(8)
				v70 = v36 + int32(-2)
				if v70 != 0 {
					v33 = v68
					v34 = v34 + int32(32)
					v36 = v70
					continue
				} else {
					break
				}
				break
			}
			v77 = v68
			v78 = v33 + int32(4)
		} else {
			v77 = v12
			v78 = int32(0)
		}
		if v17&int32(4) != 0 {
			v105 = v78
		} else {
			v87 = v78 << (uint(int32(2)) % 32)
			v92 = int32(0)
			v93 = base.Simd_g_v128_load(m, l1+v87+int32(4), v92)
			v96 = base.Simd_g_v128_load(m, l0+v87, v92)
			v97 = base.Simd_g_i8x16_add(v93, v96)
			base.Simd_g_v128_store(m, l3+v87, v92, v97)
			v105 = v77
		}
	}
	if l2 == v105 {
	} else {
		v113 = v105 << (uint(int32(2)) % 32)
		v118 = m.G34
		v119 = *(*int32)(unsafe.Add(mBase, uint32(v118)+12))
		m.T0[v119].(func(*base.Module, int32, int32, int32, int32))(m, l0+v113, l1+v113, l2-v105, l3+v113)
		mBase = m.M
	}
	return
}
func F_PredictorAdd4_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v21 int32
	_ = v21
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v36 base.V128
	_ = v36
	var v39 base.V128
	_ = v39
	var v43 base.V128
	_ = v43
	var v48 base.V128
	_ = v48
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	var v64 int32
	_ = v64
	var v74 int32
	_ = v74
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v101 int32
	_ = v101
	var v110 int32
	_ = v110
	var v117 int32
	_ = v117
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v64 = int32(0)
			v74 = v64 << (uint(int32(2)) % 32)
			v85 = l3 + v74
			v86 = l2 - v64
			v87 = l0 + v74
			v88 = v74 + l1 + int32(-4)
			for {
				v92 = *(*int32)(unsafe.Add(mBase, uint32(v87)))
				v93 = int32(-16711936)
				v95 = *(*int32)(unsafe.Add(mBase, uint32(v88)))
				v101 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(v85))) = (v92&v93+v95&v93)&v93 | (v92&v101+v95&v101)&v101
				v110 = int32(4)
				v117 = v86 + int32(-1)
				if v117 != 0 {
					v85 = v85 + v110
					v86 = v117
					v87 = v87 + v110
					v88 = v88 + v110
					continue
				} else {
					break
				}
				break
			}
		} else {
			v21 = l2 & int32(2147483644)
			v27 = l1 + int32(-4)
			v28 = v21
			v29 = l0
			v30 = l3
			for {
				v34 = int32(0)
				v35 = base.Simd_g_v128_load(m, v29, v34)
				v36 = base.Simd_g_const(&F_PredictorAdd4_C__k0)
				v39 = base.Simd_g_v128_load(m, v27, v34)
				v43 = base.Simd_g_const(&F_PredictorAdd4_C__k1)
				v48 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v35, v36), base.Simd_g_v128_and(v39, v36)), v36), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v35, v43), base.Simd_g_v128_and(v39, v43)), v43))
				base.Simd_g_v128_store(m, v30, v34, v48)
				v51 = int32(16)
				v58 = v28 + int32(-4)
				if v58 != 0 {
					v27 = v27 + v51
					v28 = v58
					v29 = v29 + v51
					v30 = v30 + v51
					continue
				} else {
					break
				}
				break
			}
			if v21 == l2 {
			} else {
				v64 = v21
				v74 = v64 << (uint(int32(2)) % 32)
				v85 = l3 + v74
				v86 = l2 - v64
				v87 = l0 + v74
				v88 = v74 + l1 + int32(-4)
				for {
					v92 = *(*int32)(unsafe.Add(mBase, uint32(v87)))
					v93 = int32(-16711936)
					v95 = *(*int32)(unsafe.Add(mBase, uint32(v88)))
					v101 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v85))) = (v92&v93+v95&v93)&v93 | (v92&v101+v95&v101)&v101
					v110 = int32(4)
					v117 = v86 + int32(-1)
					if v117 != 0 {
						v85 = v85 + v110
						v86 = v117
						v87 = v87 + v110
						v88 = v88 + v110
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

var F_PredictorAdd4_C__k0 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_PredictorAdd4_C__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_PredictorAdd4_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 int32
	_ = v46
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v52 int32
	_ = v52
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v87 int32
	_ = v87
	var v92 int32
	_ = v92
	var v93 base.V128
	_ = v93
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v105 int32
	_ = v105
	var v113 int32
	_ = v113
	var v118 int32
	_ = v118
	var v119 int32
	_ = v119
	v12 = int32(4)
	if l2 < v12 {
		v105 = int32(0)
	} else {
		v17 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v17) {
			v33 = int32(4)
			v34 = int32(0)
			v36 = (int32(base.Ui32(v17)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v40 = l3 + v34
				v41 = l1 + v34
				v44 = int32(0)
				v45 = base.Simd_g_v128_load(m, v41+int32(-4), v44)
				v46 = l0 + v34
				v48 = base.Simd_g_v128_load(m, v46, v44)
				v49 = base.Simd_g_i8x16_add(v45, v48)
				base.Simd_g_v128_store(m, v40, v44, v49)
				v52 = int32(16)
				v57 = base.Simd_g_v128_load(m, v41+int32(12), v44)
				v61 = base.Simd_g_v128_load(m, v46+v52, v44)
				v62 = base.Simd_g_i8x16_add(v57, v61)
				base.Simd_g_v128_store(m, v40+v52, v44, v62)
				v68 = v33 + int32(8)
				v70 = v36 + int32(-2)
				if v70 != 0 {
					v33 = v68
					v34 = v34 + int32(32)
					v36 = v70
					continue
				} else {
					break
				}
				break
			}
			v77 = v68
			v78 = v33 + int32(4)
		} else {
			v77 = v12
			v78 = int32(0)
		}
		if v17&int32(4) != 0 {
			v105 = v78
		} else {
			v87 = v78 << (uint(int32(2)) % 32)
			v92 = int32(0)
			v93 = base.Simd_g_v128_load(m, l1+int32(-4)+v87, v92)
			v96 = base.Simd_g_v128_load(m, l0+v87, v92)
			v97 = base.Simd_g_i8x16_add(v93, v96)
			base.Simd_g_v128_store(m, l3+v87, v92, v97)
			v105 = v77
		}
	}
	if l2 == v105 {
	} else {
		v113 = v105 << (uint(int32(2)) % 32)
		v118 = m.G34
		v119 = *(*int32)(unsafe.Add(mBase, uint32(v118)+16))
		m.T0[v119].(func(*base.Module, int32, int32, int32, int32))(m, l0+v113, l1+v113, l2-v105, l3+v113)
		mBase = m.M
	}
	return
}
func F_PredictorAdd5_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
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
	var v23 base.V128
	_ = v23
	var v24 base.V128
	_ = v24
	var v25 base.V128
	_ = v25
	var v33 int32
	_ = v33
	var v36 base.V128
	_ = v36
	var v41 base.V128
	_ = v41
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v47 int32
	_ = v47
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
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
			v22 = int32(0)
			v23 = base.Simd_g_v128_load32_zero(m, v21, v22)
			v24 = base.Simd_g_const(&F_PredictorAdd5_SSE2__k0)
			v25 = base.Simd_g_const(&F_PredictorAdd5_SSE2__k1)
			v33 = int32(1)
			v36 = base.Simd_g_v128_load32_zero(m, v14, v22)
			v41 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v23, v24, base.Simd_g_const(&F_PredictorAdd5_SSE2__k2), base.Simd_g_const(&F_PredictorAdd5_SSE2__k3)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_replace_lane_l0(v24, v17), v24, base.Simd_g_const(&F_PredictorAdd5_SSE2__k2), base.Simd_g_const(&F_PredictorAdd5_SSE2__k3))), v33), base.Simd_g_i8x16_shuffle2(v36, v24, base.Simd_g_const(&F_PredictorAdd5_SSE2__k2), base.Simd_g_const(&F_PredictorAdd5_SSE2__k3))), v33)
			v44 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_narrow_i16x8_u(v41, v41))
			v45 = int32(-16711936)
			v47 = *(*int32)(unsafe.Add(mBase, uint32(v13)))
			v53 = int32(16711935)
			v60 = (v44&v45+v47&v45)&v45 | (v44&v53+v47&v53)&v53
			*(*int32)(unsafe.Add(mBase, uint32(v16))) = v60
			v67 = v15 + int32(-1)
			if v67 != 0 {
				v13 = v13 + v20
				v14 = v21
				v15 = v67
				v16 = v16 + v20
				v17 = v60
				continue
			} else {
				break
			}
			break
		}
	}
	return
}

var F_PredictorAdd5_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorAdd5_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_PredictorAdd5_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_PredictorAdd5_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}

func F_PredictorAdd6_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
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
	var v26 base.V128
	_ = v26
	var v28 base.V128
	_ = v28
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v38 int32
	_ = v38
	var v44 int32
	_ = v44
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v60 int32
	_ = v60
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
			v25 = int32(0)
			v26 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_PredictorAdd6_SSE2__k0), v20)
			v28 = base.Simd_g_v128_load32_zero(m, v17, v25)
			v35 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v26, v28), base.Simd_g_v128_and(base.Simd_g_v128_xor(v28, v26), base.Simd_g_const(&F_PredictorAdd6_SSE2__k1))))
			v36 = int32(-16711936)
			v38 = *(*int32)(unsafe.Add(mBase, uint32(v16)))
			v44 = int32(16711935)
			v51 = (v35&v36+v38&v36)&v36 | (v35&v44+v38&v44)&v44
			*(*int32)(unsafe.Add(mBase, uint32(v19))) = v51
			v53 = int32(4)
			v60 = v18 + int32(-1)
			if v60 != 0 {
				v16 = v16 + v53
				v17 = v17 + v53
				v18 = v60
				v19 = v19 + v53
				v20 = v51
				continue
			} else {
				break
			}
			break
		}
	}
	return
}

var F_PredictorAdd6_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorAdd6_SSE2__k1 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorAdd7_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
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
	var v23 int32
	_ = v23
	var v24 base.V128
	_ = v24
	var v26 base.V128
	_ = v26
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v42 int32
	_ = v42
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v58 int32
	_ = v58
	if l2 < int32(1) {
	} else {
		v13 = *(*int32)(unsafe.Add(mBase, uint32(l3+int32(-4))))
		v14 = l0
		v15 = l1
		v16 = l2
		v17 = l3
		v18 = v13
		for {
			v23 = int32(0)
			v24 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_PredictorAdd7_SSE2__k0), v18)
			v26 = base.Simd_g_v128_load32_zero(m, v15, v23)
			v33 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v24, v26), base.Simd_g_v128_and(base.Simd_g_v128_xor(v26, v24), base.Simd_g_const(&F_PredictorAdd7_SSE2__k1))))
			v34 = int32(-16711936)
			v36 = *(*int32)(unsafe.Add(mBase, uint32(v14)))
			v42 = int32(16711935)
			v49 = (v33&v34+v36&v34)&v34 | (v33&v42+v36&v42)&v42
			*(*int32)(unsafe.Add(mBase, uint32(v17))) = v49
			v51 = int32(4)
			v58 = v16 + int32(-1)
			if v58 != 0 {
				v14 = v14 + v51
				v15 = v15 + v51
				v16 = v58
				v17 = v17 + v51
				v18 = v49
				continue
			} else {
				break
			}
			break
		}
	}
	return
}

var F_PredictorAdd7_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorAdd7_SSE2__k1 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorAdd8_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v21 int32
	_ = v21
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v36 int32
	_ = v36
	var v37 base.V128
	_ = v37
	var v39 base.V128
	_ = v39
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v50 base.V128
	_ = v50
	var v54 base.V128
	_ = v54
	var v59 base.V128
	_ = v59
	var v62 int32
	_ = v62
	var v69 int32
	_ = v69
	var v75 int32
	_ = v75
	var v85 int32
	_ = v85
	var v96 int32
	_ = v96
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v106 int32
	_ = v106
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v116 int32
	_ = v116
	var v122 int32
	_ = v122
	var v136 int32
	_ = v136
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v75 = int32(0)
			v85 = v75 << (uint(int32(2)) % 32)
			v96 = l3 + v85
			v97 = l2 - v75
			v98 = l0 + v85
			v99 = v85 + l1 + int32(-4)
			for {
				v103 = int32(4)
				v104 = v99 + v103
				v105 = *(*int32)(unsafe.Add(mBase, uint32(v104)))
				v106 = *(*int32)(unsafe.Add(mBase, uint32(v99)))
				v113 = int32(base.Ui32(v105^v106)>>(uint(int32(1))%32))&int32(2139062143) + v105&v106
				v114 = int32(-16711936)
				v116 = *(*int32)(unsafe.Add(mBase, uint32(v98)))
				v122 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(v96))) = (v113&v114+v116&v114)&v114 | (v113&v122+v116&v122)&v122
				v136 = v97 + int32(-1)
				if v136 != 0 {
					v96 = v96 + v103
					v97 = v136
					v98 = v98 + v103
					v99 = v104
					continue
				} else {
					break
				}
				break
			}
		} else {
			v21 = l2 & int32(2147483644)
			v27 = l1 + int32(-4)
			v28 = v21
			v29 = l0
			v30 = l3
			for {
				v36 = int32(0)
				v37 = base.Simd_g_v128_load_rng(m, v27+int32(4), v36, int32(-4), int32(20))
				v39 = base.Simd_g_v128_load_nc(m, v27, v36)
				v46 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_v128_xor(v37, v39), int32(1)), base.Simd_g_const(&F_PredictorAdd8_C__k0)), base.Simd_g_v128_and(v37, v39))
				v47 = base.Simd_g_const(&F_PredictorAdd8_C__k1)
				v50 = base.Simd_g_v128_load(m, v29, v36)
				v54 = base.Simd_g_const(&F_PredictorAdd8_C__k2)
				v59 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v46, v47), base.Simd_g_v128_and(v50, v47)), v47), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v46, v54), base.Simd_g_v128_and(v50, v54)), v54))
				base.Simd_g_v128_store(m, v30, v36, v59)
				v62 = int32(16)
				v69 = v28 + int32(-4)
				if v69 != 0 {
					v27 = v27 + v62
					v28 = v69
					v29 = v29 + v62
					v30 = v30 + v62
					continue
				} else {
					break
				}
				break
			}
			if v21 == l2 {
			} else {
				v75 = v21
				v85 = v75 << (uint(int32(2)) % 32)
				v96 = l3 + v85
				v97 = l2 - v75
				v98 = l0 + v85
				v99 = v85 + l1 + int32(-4)
				for {
					v103 = int32(4)
					v104 = v99 + v103
					v105 = *(*int32)(unsafe.Add(mBase, uint32(v104)))
					v106 = *(*int32)(unsafe.Add(mBase, uint32(v99)))
					v113 = int32(base.Ui32(v105^v106)>>(uint(int32(1))%32))&int32(2139062143) + v105&v106
					v114 = int32(-16711936)
					v116 = *(*int32)(unsafe.Add(mBase, uint32(v98)))
					v122 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v96))) = (v113&v114+v116&v114)&v114 | (v113&v122+v116&v122)&v122
					v136 = v97 + int32(-1)
					if v136 != 0 {
						v96 = v96 + v103
						v97 = v136
						v98 = v98 + v103
						v99 = v104
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

var F_PredictorAdd8_C__k0 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_PredictorAdd8_C__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_PredictorAdd8_C__k2 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_PredictorAdd8_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v53 base.V128
	_ = v53
	var v55 int32
	_ = v55
	var v57 base.V128
	_ = v57
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v65 int32
	_ = v65
	var v70 base.V128
	_ = v70
	var v74 base.V128
	_ = v74
	var v79 base.V128
	_ = v79
	var v83 base.V128
	_ = v83
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 base.V128
	_ = v115
	var v119 base.V128
	_ = v119
	var v123 base.V128
	_ = v123
	var v128 base.V128
	_ = v128
	var v136 int32
	_ = v136
	var v147 int32
	_ = v147
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	v15 = int32(4)
	if l2 < v15 {
		v136 = int32(0)
	} else {
		v20 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v20) {
			v36 = int32(4)
			v37 = int32(0)
			v39 = (int32(base.Ui32(v20)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v46 = l3 + v37
				v47 = l1 + v37
				v48 = int32(0)
				v49 = base.Simd_g_v128_load(m, v47, v48)
				v53 = base.Simd_g_v128_load(m, v47+int32(-4), v48)
				v55 = l0 + v37
				v57 = base.Simd_g_v128_load(m, v55, v48)
				v60 = base.Simd_g_const(&F_PredictorAdd8_SSE2__k0)
				v62 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(v49, v53), v57), base.Simd_g_v128_and(base.Simd_g_v128_xor(v49, v53), v60))
				base.Simd_g_v128_store(m, v46, v48, v62)
				v65 = int32(16)
				v70 = base.Simd_g_v128_load_rng(m, v47+v65, v48, int32(-4), int32(20))
				v74 = base.Simd_g_v128_load_nc(m, v47+int32(12), v48)
				v79 = base.Simd_g_v128_load(m, v55+v65, v48)
				v83 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(v70, v74), v79), base.Simd_g_v128_and(base.Simd_g_v128_xor(v70, v74), v60))
				base.Simd_g_v128_store(m, v46+v65, v48, v83)
				v89 = v36 + int32(8)
				v91 = v39 + int32(-2)
				if v91 != 0 {
					v36 = v89
					v37 = v37 + int32(32)
					v39 = v91
					continue
				} else {
					break
				}
				break
			}
			v98 = v89
			v99 = v36 + int32(4)
		} else {
			v98 = v15
			v99 = int32(0)
		}
		if v20&int32(4) != 0 {
			v136 = v99
		} else {
			v111 = v99 << (uint(int32(2)) % 32)
			v113 = l1 + v111
			v114 = int32(0)
			v115 = base.Simd_g_v128_load(m, v113, v114)
			v119 = base.Simd_g_v128_load(m, v113+int32(-4), v114)
			v123 = base.Simd_g_v128_load(m, l0+v111, v114)
			v128 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(v115, v119), v123), base.Simd_g_v128_and(base.Simd_g_v128_xor(v115, v119), base.Simd_g_const(&F_PredictorAdd8_SSE2__k0)))
			base.Simd_g_v128_store(m, l3+v111, v114, v128)
			v136 = v98
		}
	}
	if l2 == v136 {
	} else {
		v147 = v136 << (uint(int32(2)) % 32)
		v152 = m.G34
		v153 = *(*int32)(unsafe.Add(mBase, uint32(v152)+32))
		m.T0[v153].(func(*base.Module, int32, int32, int32, int32))(m, l0+v147, l1+v147, l2-v136, l3+v147)
		mBase = m.M
	}
	return
}

var F_PredictorAdd8_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorAdd9_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v37 base.V128
	_ = v37
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v48 base.V128
	_ = v48
	var v52 base.V128
	_ = v52
	var v57 base.V128
	_ = v57
	var v60 int32
	_ = v60
	var v67 int32
	_ = v67
	var v73 int32
	_ = v73
	var v83 int32
	_ = v83
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v94 int32
	_ = v94
	var v95 int32
	_ = v95
	var v99 int32
	_ = v99
	var v100 int32
	_ = v100
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	var v109 int32
	_ = v109
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v118 int32
	_ = v118
	var v132 int32
	_ = v132
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v73 = int32(0)
			v83 = v73 << (uint(int32(2)) % 32)
			v92 = l2 - v73
			v93 = l1 + v83
			v94 = l3 + v83
			v95 = l0 + v83
			for {
				v99 = int32(4)
				v100 = v93 + v99
				v101 = *(*int32)(unsafe.Add(mBase, uint32(v100)))
				v102 = *(*int32)(unsafe.Add(mBase, uint32(v93)))
				v109 = int32(base.Ui32(v101^v102)>>(uint(int32(1))%32))&int32(2139062143) + v101&v102
				v110 = int32(-16711936)
				v112 = *(*int32)(unsafe.Add(mBase, uint32(v95)))
				v118 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(v94))) = (v109&v110+v112&v110)&v110 | (v109&v118+v112&v118)&v118
				v132 = v92 + int32(-1)
				if v132 != 0 {
					v92 = v132
					v93 = v100
					v94 = v94 + v99
					v95 = v95 + v99
					continue
				} else {
					break
				}
				break
			}
		} else {
			v19 = l2 & int32(2147483644)
			v25 = v19
			v26 = l0
			v27 = l1
			v28 = l3
			for {
				v34 = int32(0)
				v35 = base.Simd_g_v128_load_rng(m, v27+int32(4), v34, int32(-4), int32(20))
				v37 = base.Simd_g_v128_load_nc(m, v27, v34)
				v44 = base.Simd_g_i32x4_add(base.Simd_g_v128_and(base.Simd_g_i32x4_shr_u(base.Simd_g_v128_xor(v35, v37), int32(1)), base.Simd_g_const(&F_PredictorAdd9_C__k0)), base.Simd_g_v128_and(v35, v37))
				v45 = base.Simd_g_const(&F_PredictorAdd9_C__k1)
				v48 = base.Simd_g_v128_load(m, v26, v34)
				v52 = base.Simd_g_const(&F_PredictorAdd9_C__k2)
				v57 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v44, v45), base.Simd_g_v128_and(v48, v45)), v45), base.Simd_g_v128_and(base.Simd_g_i32x4_add(base.Simd_g_v128_and(v44, v52), base.Simd_g_v128_and(v48, v52)), v52))
				base.Simd_g_v128_store(m, v28, v34, v57)
				v60 = int32(16)
				v67 = v25 + int32(-4)
				if v67 != 0 {
					v25 = v67
					v26 = v26 + v60
					v27 = v27 + v60
					v28 = v28 + v60
					continue
				} else {
					break
				}
				break
			}
			if v19 == l2 {
			} else {
				v73 = v19
				v83 = v73 << (uint(int32(2)) % 32)
				v92 = l2 - v73
				v93 = l1 + v83
				v94 = l3 + v83
				v95 = l0 + v83
				for {
					v99 = int32(4)
					v100 = v93 + v99
					v101 = *(*int32)(unsafe.Add(mBase, uint32(v100)))
					v102 = *(*int32)(unsafe.Add(mBase, uint32(v93)))
					v109 = int32(base.Ui32(v101^v102)>>(uint(int32(1))%32))&int32(2139062143) + v101&v102
					v110 = int32(-16711936)
					v112 = *(*int32)(unsafe.Add(mBase, uint32(v95)))
					v118 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v94))) = (v109&v110+v112&v110)&v110 | (v109&v118+v112&v118)&v118
					v132 = v92 + int32(-1)
					if v132 != 0 {
						v92 = v132
						v93 = v100
						v94 = v94 + v99
						v95 = v95 + v99
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

var F_PredictorAdd9_C__k0 = [2]uint64{0x7f7f7f7f7f7f7f7f, 0x7f7f7f7f7f7f7f7f}
var F_PredictorAdd9_C__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_PredictorAdd9_C__k2 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}

func F_PredictorAdd9_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v30 base.V128
	_ = v30
	var v31 int32
	_ = v31
	var v34 base.V128
	_ = v34
	var v37 base.V128
	_ = v37
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	if l2 < int32(4) {
		v60 = int32(0)
	} else {
		v22 = int32(0)
		v23 = l3
		v24 = l0
		v25 = l1
		for {
			v29 = int32(0)
			v30 = base.Simd_g_v128_load_rng(m, v25, v29, int32(0), int32(20))
			v31 = int32(4)
			v34 = base.Simd_g_v128_load_nc(m, v25+v31, v29)
			v37 = base.Simd_g_v128_load(m, v24, v29)
			v42 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_add(base.Simd_g_i8x16_avgr_u(v30, v34), v37), base.Simd_g_v128_and(base.Simd_g_v128_xor(v30, v34), base.Simd_g_const(&F_PredictorAdd9_SSE2__k0)))
			base.Simd_g_v128_store(m, v23, v29, v42)
			v45 = int32(16)
			v54 = v22 + v31
			if v22+int32(8) <= l2 {
				v22 = v54
				v23 = v23 + v45
				v24 = v24 + v45
				v25 = v25 + v45
				continue
			} else {
				break
			}
			break
		}
		v60 = v54
	}
	if l2 == v60 {
	} else {
		v70 = v60 << (uint(int32(2)) % 32)
		v75 = m.G34
		v76 = *(*int32)(unsafe.Add(mBase, uint32(v75)+36))
		m.T0[v76].(func(*base.Module, int32, int32, int32, int32))(m, l0+v70, l1+v70, l2-v60, l3+v70)
		mBase = m.M
	}
	return
}

var F_PredictorAdd9_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorSub0_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 base.V128
	_ = v25
	var v27 base.V128
	_ = v27
	var v30 int32
	_ = v30
	var v35 int32
	_ = v35
	var v41 int32
	_ = v41
	var v47 int32
	_ = v47
	var v55 int32
	_ = v55
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v62 int32
	_ = v62
	var v67 int32
	_ = v67
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v41 = int32(0)
			v47 = v41 << (uint(int32(2)) % 32)
			v55 = l2 - v41
			v56 = l0 + v47
			v57 = l3 + v47
			for {
				v58 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
				*(*int32)(unsafe.Add(mBase, uint32(v57))) = v58 + int32(16777216)
				v62 = int32(4)
				v67 = v55 + int32(-1)
				if v67 != 0 {
					v55 = v67
					v56 = v56 + v62
					v57 = v57 + v62
					continue
				} else {
					break
				}
				break
			}
		} else {
			v15 = l2 & int32(2147483644)
			v21 = v15
			v22 = l3
			v23 = l0
			for {
				v24 = int32(0)
				v25 = base.Simd_g_v128_load(m, v23, v24)
				v27 = base.Simd_g_i32x4_add(v25, base.Simd_g_const(&F_PredictorSub0_C__k0))
				base.Simd_g_v128_store(m, v22, v24, v27)
				v30 = int32(16)
				v35 = v21 + int32(-4)
				if v35 != 0 {
					v21 = v35
					v22 = v22 + v30
					v23 = v23 + v30
					continue
				} else {
					break
				}
				break
			}
			if v15 == l2 {
			} else {
				v41 = v15
				v47 = v41 << (uint(int32(2)) % 32)
				v55 = l2 - v41
				v56 = l0 + v47
				v57 = l3 + v47
				for {
					v58 = *(*int32)(unsafe.Add(mBase, uint32(v56)))
					*(*int32)(unsafe.Add(mBase, uint32(v57))) = v58 + int32(16777216)
					v62 = int32(4)
					v67 = v55 + int32(-1)
					if v67 != 0 {
						v55 = v67
						v56 = v56 + v62
						v57 = v57 + v62
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

var F_PredictorSub0_C__k0 = [2]uint64{0x100000001000000, 0x100000001000000}

func F_PredictorSub0_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
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
	var v44 int32
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v47 base.V128
	_ = v47
	var v50 int32
	_ = v50
	var v55 base.V128
	_ = v55
	var v56 base.V128
	_ = v56
	var v59 int32
	_ = v59
	var v64 base.V128
	_ = v64
	var v65 base.V128
	_ = v65
	var v68 int32
	_ = v68
	var v73 base.V128
	_ = v73
	var v74 base.V128
	_ = v74
	var v80 int32
	_ = v80
	var v82 int32
	_ = v82
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v102 int32
	_ = v102
	var __phi102 int32
	_ = __phi102
	var v103 int32
	_ = v103
	var __phi103 int32
	_ = __phi103
	var v105 int32
	_ = v105
	var __phi105 int32
	_ = __phi105
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v114 base.V128
	_ = v114
	var v116 base.V128
	_ = v116
	var v122 int32
	_ = v122
	var v130 int32
	_ = v130
	var v138 int32
	_ = v138
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	v12 = int32(4)
	if l2 < v12 {
		v130 = int32(0)
	} else {
		v17 = l2 + int32(-4)
		v21 = int32(base.Ui32(v17)>>(uint(int32(2))%32)) + int32(1)
		v23 = v21 & int32(3)
		if base.Ui32(int32(12)) <= base.Ui32(v17) {
			v35 = int32(4)
			v36 = int32(0)
			v39 = v21 & int32(2147483644)
			for {
				v42 = l3 + v36
				v43 = l0 + v36
				v44 = int32(0)
				v45 = base.Simd_g_v128_load(m, v43, v44)
				v46 = base.Simd_g_const(&F_PredictorSub0_SSE2__k0)
				v47 = base.Simd_g_i8x16_add(v45, v46)
				base.Simd_g_v128_store(m, v42, v44, v47)
				v50 = int32(16)
				v55 = base.Simd_g_v128_load(m, v43+v50, v44)
				v56 = base.Simd_g_i8x16_add(v55, v46)
				base.Simd_g_v128_store(m, v42+v50, v44, v56)
				v59 = int32(32)
				v64 = base.Simd_g_v128_load(m, v43+v59, v44)
				v65 = base.Simd_g_i8x16_add(v64, v46)
				base.Simd_g_v128_store(m, v42+v59, v44, v65)
				v68 = int32(48)
				v73 = base.Simd_g_v128_load(m, v43+v68, v44)
				v74 = base.Simd_g_i8x16_add(v73, v46)
				base.Simd_g_v128_store(m, v42+v68, v44, v74)
				v80 = v35 + v50
				v82 = v39 + int32(-4)
				if v82 != 0 {
					v35 = v80
					v36 = v36 + int32(64)
					v39 = v82
					continue
				} else {
					break
				}
				break
			}
			v89 = v80
			v90 = v35 + int32(12)
		} else {
			v89 = v12
			v90 = int32(0)
		}
		if v23 == int32(0) {
			v130 = v90
		} else {
			__phi102 = v89
			__phi103 = v90
			__phi105 = v23
			v102 = __phi102
			v103 = __phi103
			v105 = __phi105
			for {
				v110 = v103 << (uint(int32(2)) % 32)
				v113 = int32(0)
				v114 = base.Simd_g_v128_load(m, l0+v110, v113)
				v116 = base.Simd_g_i8x16_add(v114, base.Simd_g_const(&F_PredictorSub0_SSE2__k0))
				base.Simd_g_v128_store(m, l3+v110, v113, v116)
				v122 = v105 + int32(-1)
				if v122 != 0 {
					__phi102 = v102 + int32(4)
					__phi103 = v102
					__phi105 = v122
					v102 = __phi102
					v103 = __phi103
					v105 = __phi105
					continue
				} else {
					break
				}
				break
			}
			v130 = v102 + int32(0)
		}
	}
	if l2 == v130 {
	} else {
		v138 = v130 << (uint(int32(2)) % 32)
		v143 = m.G111
		v144 = *(*int32)(unsafe.Add(mBase, uint32(v143)))
		m.T0[v144].(func(*base.Module, int32, int32, int32, int32))(m, l0+v138, int32(0), l2-v130, l3+v138)
		mBase = m.M
	}
	return
}

var F_PredictorSub0_SSE2__k0 = [2]uint64{0x100000001000000, 0x100000001000000}

func F_PredictorSub10_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v40 base.V128
	_ = v40
	var v44 base.V128
	_ = v44
	var v45 int32
	_ = v45
	var v48 base.V128
	_ = v48
	var v51 base.V128
	_ = v51
	var v53 base.V128
	_ = v53
	var v55 base.V128
	_ = v55
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v66 base.V128
	_ = v66
	var v69 int32
	_ = v69
	var v78 int32
	_ = v78
	var v84 int32
	_ = v84
	var v96 int32
	_ = v96
	var v101 int32
	_ = v101
	var v102 int32
	_ = v102
	if l2 < int32(4) {
		v84 = int32(0)
	} else {
		v18 = int32(-4)
		v28 = l1 + v18
		v29 = l0 + v18
		v30 = int32(0)
		v31 = l3
		for {
			v37 = int32(4)
			v39 = int32(0)
			v40 = base.Simd_g_v128_load_rng(m, v29+v37, v39, int32(-4), int32(20))
			v44 = base.Simd_g_v128_load_rng(m, v28+v37, v39, int32(-4), int32(24))
			v45 = int32(8)
			v48 = base.Simd_g_v128_load_nc(m, v28+v45, v39)
			v51 = base.Simd_g_const(&F_PredictorSub10_SSE2__k0)
			v53 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v44, v48), base.Simd_g_v128_and(base.Simd_g_v128_xor(v48, v44), v51))
			v55 = base.Simd_g_v128_load_nc(m, v29, v39)
			v57 = base.Simd_g_v128_load_nc(m, v28, v39)
			v61 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v55, v57), base.Simd_g_v128_and(base.Simd_g_v128_xor(v57, v55), v51))
			v66 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v40, base.Simd_g_i8x16_avgr_u(v53, v61)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v53, v61), v51))
			base.Simd_g_v128_store(m, v31, v39, v66)
			v69 = int32(16)
			v78 = v30 + v37
			if v30+v45 <= l2 {
				v28 = v28 + v69
				v29 = v29 + v69
				v30 = v78
				v31 = v31 + v69
				continue
			} else {
				break
			}
			break
		}
		v84 = v78
	}
	if l2 == v84 {
	} else {
		v96 = v84 << (uint(int32(2)) % 32)
		v101 = m.G111
		v102 = *(*int32)(unsafe.Add(mBase, uint32(v101)+40))
		m.T0[v102].(func(*base.Module, int32, int32, int32, int32))(m, l0+v96, l1+v96, l2-v84, l3+v96)
		mBase = m.M
	}
	return
}

var F_PredictorSub10_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorSub11_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v19 int32
	_ = v19
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v42 base.V128
	_ = v42
	var v44 base.V128
	_ = v44
	var v48 base.V128
	_ = v48
	var v50 base.V128
	_ = v50
	var v51 base.V128
	_ = v51
	var v52 base.V128
	_ = v52
	var v54 base.V128
	_ = v54
	var v57 base.V128
	_ = v57
	var v58 int32
	_ = v58
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v63 int32
	_ = v63
	var v65 base.V128
	_ = v65
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v71 base.V128
	_ = v71
	var v72 base.V128
	_ = v72
	var v74 base.V128
	_ = v74
	var v77 base.V128
	_ = v77
	var v81 base.V128
	_ = v81
	var v84 base.V128
	_ = v84
	var v92 base.V128
	_ = v92
	var v94 base.V128
	_ = v94
	var v97 base.V128
	_ = v97
	var v101 base.V128
	_ = v101
	var v104 base.V128
	_ = v104
	var v111 base.V128
	_ = v111
	var v113 base.V128
	_ = v113
	var v116 base.V128
	_ = v116
	var v120 base.V128
	_ = v120
	var v123 base.V128
	_ = v123
	var v132 base.V128
	_ = v132
	var v144 int32
	_ = v144
	var v150 int32
	_ = v150
	var v163 int32
	_ = v163
	var v168 int32
	_ = v168
	var v169 int32
	_ = v169
	if l2 < int32(4) {
		v150 = int32(0)
	} else {
		v19 = int32(-4)
		v29 = l0 + v19
		v30 = l1 + v19
		v31 = int32(0)
		v32 = l3
		for {
			v39 = int32(4)
			v41 = int32(0)
			v42 = base.Simd_g_v128_load_rng(m, v29+v39, v41, int32(-4), int32(20))
			v44 = base.Simd_g_v128_load_nc(m, v29, v41)
			v48 = base.Simd_g_v128_load_rng(m, v30+v39, v41, int32(-4), int32(20))
			v50 = base.Simd_g_v128_load_nc(m, v30, v41)
			v51 = base.Simd_g_const(&F_PredictorSub11_SSE2__k0)
			v52 = base.Simd_g_i8x16_shuffle2(v50, v44, base.Simd_g_const(&F_PredictorSub11_SSE2__k1), base.Simd_g_const(&F_PredictorSub11_SSE2__k2))
			v54 = base.Simd_g_i8x16_shuffle2(v44, v44, base.Simd_g_const(&F_PredictorSub11_SSE2__k1), base.Simd_g_const(&F_PredictorSub11_SSE2__k2))
			v57 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v52, v54), base.Simd_g_i8x16_sub_sat_u(v54, v52))
			v58 = int32(8)
			v60 = base.Simd_g_const(&F_PredictorSub11_SSE2__k3)
			v62 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v57, v58), base.Simd_g_v128_and(v57, v60))
			v63 = int32(16)
			v65 = base.Simd_g_i16x8_add(v62, base.Simd_g_i32x4_shl(v62, v63))
			v66 = int32(32)
			v69 = int32(48)
			v71 = base.Simd_g_const(&F_PredictorSub11_SSE2__k4)
			v72 = base.Simd_g_i8x16_shuffle2(v50, v44, base.Simd_g_const(&F_PredictorSub11_SSE2__k5), base.Simd_g_const(&F_PredictorSub11_SSE2__k6))
			v74 = base.Simd_g_i8x16_shuffle2(v44, v44, base.Simd_g_const(&F_PredictorSub11_SSE2__k5), base.Simd_g_const(&F_PredictorSub11_SSE2__k6))
			v77 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v72, v74), base.Simd_g_i8x16_sub_sat_u(v74, v72))
			v81 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v77, v58), base.Simd_g_v128_and(v77, v60))
			v84 = base.Simd_g_i16x8_add(v81, base.Simd_g_i32x4_shl(v81, v63))
			v92 = base.Simd_g_i8x16_shuffle2(v50, v48, base.Simd_g_const(&F_PredictorSub11_SSE2__k1), base.Simd_g_const(&F_PredictorSub11_SSE2__k2))
			v94 = base.Simd_g_i8x16_shuffle2(v48, v48, base.Simd_g_const(&F_PredictorSub11_SSE2__k1), base.Simd_g_const(&F_PredictorSub11_SSE2__k2))
			v97 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v92, v94), base.Simd_g_i8x16_sub_sat_u(v94, v92))
			v101 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v97, v58), base.Simd_g_v128_and(v97, v60))
			v104 = base.Simd_g_i16x8_add(v101, base.Simd_g_i32x4_shl(v101, v63))
			v111 = base.Simd_g_i8x16_shuffle2(v50, v48, base.Simd_g_const(&F_PredictorSub11_SSE2__k5), base.Simd_g_const(&F_PredictorSub11_SSE2__k6))
			v113 = base.Simd_g_i8x16_shuffle2(v48, v48, base.Simd_g_const(&F_PredictorSub11_SSE2__k5), base.Simd_g_const(&F_PredictorSub11_SSE2__k6))
			v116 = base.Simd_g_v128_or(base.Simd_g_i8x16_sub_sat_u(v111, v113), base.Simd_g_i8x16_sub_sat_u(v113, v111))
			v120 = base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_u(v116, v58), base.Simd_g_v128_and(v116, v60))
			v123 = base.Simd_g_i16x8_add(v120, base.Simd_g_i32x4_shl(v120, v63))
			v132 = base.Simd_g_i8x16_sub(v42, base.Simd_g_v128_bitselect(v44, v48, base.Simd_g_i32x4_gt_s(base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v65, base.Simd_g_i64x2_shl(v65, v66)), v69), base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v84, base.Simd_g_i64x2_shl(v84, v66)), v69)), base.Simd_g_i16x8_narrow_i32x4_s(base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v104, base.Simd_g_i64x2_shl(v104, v66)), v69), base.Simd_g_i64x2_shr_u(base.Simd_g_i16x8_add(v123, base.Simd_g_i64x2_shl(v123, v66)), v69)))))
			base.Simd_g_v128_store(m, v32, v41, v132)
			v144 = v31 + v39
			if v31+v58 <= l2 {
				v29 = v29 + v63
				v30 = v30 + v63
				v31 = v144
				v32 = v32 + v63
				continue
			} else {
				break
			}
			break
		}
		v150 = v144
	}
	if l2 == v150 {
	} else {
		v163 = v150 << (uint(int32(2)) % 32)
		v168 = m.G111
		v169 = *(*int32)(unsafe.Add(mBase, uint32(v168)+44))
		m.T0[v169].(func(*base.Module, int32, int32, int32, int32))(m, l0+v163, l1+v163, l2-v150, l3+v163)
		mBase = m.M
	}
	return
}

var F_PredictorSub11_SSE2__k0 = [2]uint64{0x1312111003020100, 0x1716151407060504}
var F_PredictorSub11_SSE2__k1 = [2]uint64{0x8080808003020100, 0x8080808007060504}
var F_PredictorSub11_SSE2__k2 = [2]uint64{0x302010080808080, 0x706050480808080}
var F_PredictorSub11_SSE2__k3 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_PredictorSub11_SSE2__k4 = [2]uint64{0x1b1a19180b0a0908, 0x1f1e1d1c0f0e0d0c}
var F_PredictorSub11_SSE2__k5 = [2]uint64{0x808080800b0a0908, 0x808080800f0e0d0c}
var F_PredictorSub11_SSE2__k6 = [2]uint64{0xb0a090880808080, 0xf0e0d0c80808080}

func F_PredictorSub12_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v40 base.V128
	_ = v40
	var v44 base.V128
	_ = v44
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v49 base.V128
	_ = v49
	var v54 base.V128
	_ = v54
	var v58 base.V128
	_ = v58
	var v67 base.V128
	_ = v67
	var v70 int32
	_ = v70
	var v79 int32
	_ = v79
	var v85 int32
	_ = v85
	var v97 int32
	_ = v97
	var v102 int32
	_ = v102
	var v103 int32
	_ = v103
	if l2 < int32(4) {
		v85 = int32(0)
	} else {
		v18 = int32(-4)
		v28 = l1 + v18
		v29 = l0 + v18
		v30 = int32(0)
		v31 = l3
		for {
			v37 = int32(4)
			v39 = int32(0)
			v40 = base.Simd_g_v128_load_rng(m, v29+v37, v39, int32(-4), int32(20))
			v44 = base.Simd_g_v128_load_rng(m, v28+v37, v39, int32(-4), int32(20))
			v45 = base.Simd_g_const(&F_PredictorSub12_SSE2__k0)
			v46 = base.Simd_g_const(&F_PredictorSub12_SSE2__k1)
			v49 = base.Simd_g_v128_load_nc(m, v29, v39)
			v54 = base.Simd_g_v128_load_nc(m, v28, v39)
			v58 = base.Simd_g_const(&F_PredictorSub12_SSE2__k2)
			v67 = base.Simd_g_i8x16_sub(v40, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_sub(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v44, v45, base.Simd_g_const(&F_PredictorSub12_SSE2__k3), base.Simd_g_const(&F_PredictorSub12_SSE2__k4)), base.Simd_g_i8x16_shuffle2(v49, v45, base.Simd_g_const(&F_PredictorSub12_SSE2__k3), base.Simd_g_const(&F_PredictorSub12_SSE2__k4))), base.Simd_g_i8x16_shuffle2(v54, v45, base.Simd_g_const(&F_PredictorSub12_SSE2__k3), base.Simd_g_const(&F_PredictorSub12_SSE2__k4))), base.Simd_g_i16x8_sub(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v44, v45, base.Simd_g_const(&F_PredictorSub12_SSE2__k5), base.Simd_g_const(&F_PredictorSub12_SSE2__k6)), base.Simd_g_i8x16_shuffle2(v49, v45, base.Simd_g_const(&F_PredictorSub12_SSE2__k5), base.Simd_g_const(&F_PredictorSub12_SSE2__k6))), base.Simd_g_i8x16_shuffle2(v54, v45, base.Simd_g_const(&F_PredictorSub12_SSE2__k5), base.Simd_g_const(&F_PredictorSub12_SSE2__k6)))))
			base.Simd_g_v128_store(m, v31, v39, v67)
			v70 = int32(16)
			v79 = v30 + v37
			if v30+int32(8) <= l2 {
				v28 = v28 + v70
				v29 = v29 + v70
				v30 = v79
				v31 = v31 + v70
				continue
			} else {
				break
			}
			break
		}
		v85 = v79
	}
	if l2 == v85 {
	} else {
		v97 = v85 << (uint(int32(2)) % 32)
		v102 = m.G111
		v103 = *(*int32)(unsafe.Add(mBase, uint32(v102)+48))
		m.T0[v103].(func(*base.Module, int32, int32, int32, int32))(m, l0+v97, l1+v97, l2-v85, l3+v97)
		mBase = m.M
	}
	return
}

var F_PredictorSub12_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorSub12_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_PredictorSub12_SSE2__k2 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_PredictorSub12_SSE2__k3 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_PredictorSub12_SSE2__k4 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_PredictorSub12_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_PredictorSub12_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}

func F_PredictorSub13_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v41 int32
	_ = v41
	var v43 int32
	_ = v43
	var v44 base.V128
	_ = v44
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 base.V128
	_ = v50
	var v53 base.V128
	_ = v53
	var v57 int32
	_ = v57
	var v58 base.V128
	_ = v58
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v69 base.V128
	_ = v69
	var v75 base.V128
	_ = v75
	var v77 base.V128
	_ = v77
	var v85 base.V128
	_ = v85
	var v88 int32
	_ = v88
	var v97 int32
	_ = v97
	var v103 int32
	_ = v103
	var v117 int32
	_ = v117
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	if l2 < int32(4) {
		v103 = int32(0)
	} else {
		v20 = int32(-4)
		v30 = l1 + v20
		v31 = l0 + v20
		v32 = int32(0)
		v33 = l3
		for {
			v41 = int32(4)
			v43 = int32(0)
			v44 = base.Simd_g_v128_load_rng(m, v31+v41, v43, int32(-4), int32(20))
			v48 = base.Simd_g_v128_load_rng(m, v30+v41, v43, int32(-4), int32(20))
			v49 = base.Simd_g_const(&F_PredictorSub13_SSE2__k0)
			v50 = base.Simd_g_const(&F_PredictorSub13_SSE2__k1)
			v53 = base.Simd_g_v128_load_nc(m, v31, v43)
			v57 = int32(1)
			v58 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v48, v49, base.Simd_g_const(&F_PredictorSub13_SSE2__k2), base.Simd_g_const(&F_PredictorSub13_SSE2__k3)), base.Simd_g_i8x16_shuffle2(v53, v49, base.Simd_g_const(&F_PredictorSub13_SSE2__k2), base.Simd_g_const(&F_PredictorSub13_SSE2__k3))), v57)
			v60 = base.Simd_g_v128_load_nc(m, v30, v43)
			v62 = base.Simd_g_i8x16_shuffle2(v60, v49, base.Simd_g_const(&F_PredictorSub13_SSE2__k2), base.Simd_g_const(&F_PredictorSub13_SSE2__k3))
			v69 = base.Simd_g_const(&F_PredictorSub13_SSE2__k4)
			v75 = base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(v48, v49, base.Simd_g_const(&F_PredictorSub13_SSE2__k5), base.Simd_g_const(&F_PredictorSub13_SSE2__k6)), base.Simd_g_i8x16_shuffle2(v53, v49, base.Simd_g_const(&F_PredictorSub13_SSE2__k5), base.Simd_g_const(&F_PredictorSub13_SSE2__k6))), v57)
			v77 = base.Simd_g_i8x16_shuffle2(v60, v49, base.Simd_g_const(&F_PredictorSub13_SSE2__k5), base.Simd_g_const(&F_PredictorSub13_SSE2__k6))
			v85 = base.Simd_g_i8x16_sub(v44, base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(base.Simd_g_i16x8_sub(v58, v62), base.Simd_g_i16x8_lt_s(v58, v62)), v57), v58), base.Simd_g_i16x8_add(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_sub(base.Simd_g_i16x8_sub(v75, v77), base.Simd_g_i16x8_lt_s(v75, v77)), v57), v75)))
			base.Simd_g_v128_store(m, v33, v43, v85)
			v88 = int32(16)
			v97 = v32 + v41
			if v32+int32(8) <= l2 {
				v30 = v30 + v88
				v31 = v31 + v88
				v32 = v97
				v33 = v33 + v88
				continue
			} else {
				break
			}
			break
		}
		v103 = v97
	}
	if l2 == v103 {
	} else {
		v117 = v103 << (uint(int32(2)) % 32)
		v122 = m.G111
		v123 = *(*int32)(unsafe.Add(mBase, uint32(v122)+52))
		m.T0[v123].(func(*base.Module, int32, int32, int32, int32))(m, l0+v117, l1+v117, l2-v103, l3+v117)
		mBase = m.M
	}
	return
}

var F_PredictorSub13_SSE2__k0 = [2]uint64{0x0, 0x0}
var F_PredictorSub13_SSE2__k1 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_PredictorSub13_SSE2__k2 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_PredictorSub13_SSE2__k3 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_PredictorSub13_SSE2__k4 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_PredictorSub13_SSE2__k5 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_PredictorSub13_SSE2__k6 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}

func F_PredictorSub1_C(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v20 int32
	_ = v20
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v34 int32
	_ = v34
	var v35 base.V128
	_ = v35
	var v39 base.V128
	_ = v39
	var v40 base.V128
	_ = v40
	var v46 base.V128
	_ = v46
	var v50 base.V128
	_ = v50
	var v53 int32
	_ = v53
	var v58 int32
	_ = v58
	var v64 int32
	_ = v64
	var v73 int32
	_ = v73
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v102 int32
	_ = v102
	var v112 int32
	_ = v112
	if l2 < int32(1) {
	} else {
		if base.Ui32(l2) < base.Ui32(int32(4)) {
			v64 = int32(0)
			v73 = v64 << (uint(int32(2)) % 32)
			v80 = l2 - v64
			v83 = l3 + v73
			v85 = v73 + l0 + int32(-4)
			for {
				v89 = *(*int32)(unsafe.Add(mBase, uint32(v85)))
				v90 = int32(4)
				v91 = v85 + v90
				v92 = *(*int32)(unsafe.Add(mBase, uint32(v91)))
				v95 = int32(-16711936)
				v102 = int32(16711935)
				*(*int32)(unsafe.Add(mBase, uint32(v83))) = (v92|int32(16711680)-v89&v95)&v95 | (v92|int32(_a_F_PredictorSub1_C_0)-v89&v102)&v102
				v112 = v80 + int32(-1)
				if v112 != 0 {
					v80 = v112
					v83 = v83 + v90
					v85 = v91
					continue
				} else {
					break
				}
				break
			}
		} else {
			v20 = l2 & int32(2147483644)
			v26 = l0 + int32(-4)
			v27 = v20
			v28 = l3
			for {
				v34 = int32(0)
				v35 = base.Simd_g_v128_load_rng(m, v26+int32(4), v34, int32(-4), int32(20))
				v39 = base.Simd_g_v128_load_nc(m, v26, v34)
				v40 = base.Simd_g_const(&F_PredictorSub1_C__k0)
				v46 = base.Simd_g_const(&F_PredictorSub1_C__k1)
				v50 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_sub(base.Simd_g_v128_or(v35, base.Simd_g_const(&F_PredictorSub1_C__k2)), base.Simd_g_v128_and(v39, v40)), v40), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(base.Simd_g_v128_or(v35, base.Simd_g_const(&F_PredictorSub1_C__k3)), base.Simd_g_v128_and(v39, v46)), v46))
				base.Simd_g_v128_store(m, v28, v34, v50)
				v53 = int32(16)
				v58 = v27 + int32(-4)
				if v58 != 0 {
					v26 = v26 + v53
					v27 = v58
					v28 = v28 + v53
					continue
				} else {
					break
				}
				break
			}
			if v20 == l2 {
			} else {
				v64 = v20
				v73 = v64 << (uint(int32(2)) % 32)
				v80 = l2 - v64
				v83 = l3 + v73
				v85 = v73 + l0 + int32(-4)
				for {
					v89 = *(*int32)(unsafe.Add(mBase, uint32(v85)))
					v90 = int32(4)
					v91 = v85 + v90
					v92 = *(*int32)(unsafe.Add(mBase, uint32(v91)))
					v95 = int32(-16711936)
					v102 = int32(16711935)
					*(*int32)(unsafe.Add(mBase, uint32(v83))) = (v92|int32(16711680)-v89&v95)&v95 | (v92|int32(_a_F_PredictorSub1_C_0)-v89&v102)&v102
					v112 = v80 + int32(-1)
					if v112 != 0 {
						v80 = v112
						v83 = v83 + v90
						v85 = v91
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

var F_PredictorSub1_C__k0 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_PredictorSub1_C__k1 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_PredictorSub1_C__k2 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_PredictorSub1_C__k3 = [2]uint64{0xff000000ff00, 0xff000000ff00}

func F_PredictorSub1_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v16 int32
	_ = v16
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v41 base.V128
	_ = v41
	var v45 base.V128
	_ = v45
	var v46 base.V128
	_ = v46
	var v49 int32
	_ = v49
	var v54 base.V128
	_ = v54
	var v58 base.V128
	_ = v58
	var v59 base.V128
	_ = v59
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v83 int32
	_ = v83
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 base.V128
	_ = v87
	var v91 base.V128
	_ = v91
	var v92 base.V128
	_ = v92
	var v100 int32
	_ = v100
	var v107 int32
	_ = v107
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v115 int32
	_ = v115
	v11 = int32(4)
	if l2 < v11 {
		v100 = int32(0)
	} else {
		v16 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v16) {
			v32 = int32(4)
			v33 = int32(0)
			v35 = (int32(base.Ui32(v16)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v38 = l3 + v33
				v39 = l0 + v33
				v40 = int32(0)
				v41 = base.Simd_g_v128_load(m, v39, v40)
				v45 = base.Simd_g_v128_load(m, v39+int32(-4), v40)
				v46 = base.Simd_g_i8x16_sub(v41, v45)
				base.Simd_g_v128_store(m, v38, v40, v46)
				v49 = int32(16)
				v54 = base.Simd_g_v128_load_rng(m, v39+v49, v40, int32(-4), int32(20))
				v58 = base.Simd_g_v128_load_nc(m, v39+int32(12), v40)
				v59 = base.Simd_g_i8x16_sub(v54, v58)
				base.Simd_g_v128_store(m, v38+v49, v40, v59)
				v65 = v32 + int32(8)
				v67 = v35 + int32(-2)
				if v67 != 0 {
					v32 = v65
					v33 = v33 + int32(32)
					v35 = v67
					continue
				} else {
					break
				}
				break
			}
			v74 = v65
			v75 = v32 + int32(4)
		} else {
			v74 = v11
			v75 = int32(0)
		}
		if v16&int32(4) != 0 {
			v100 = v75
		} else {
			v83 = v75 << (uint(int32(2)) % 32)
			v85 = l0 + v83
			v86 = int32(0)
			v87 = base.Simd_g_v128_load(m, v85, v86)
			v91 = base.Simd_g_v128_load(m, v85+int32(-4), v86)
			v92 = base.Simd_g_i8x16_sub(v87, v91)
			base.Simd_g_v128_store(m, l3+v83, v86, v92)
			v100 = v74
		}
	}
	if l2 == v100 {
	} else {
		v107 = v100 << (uint(int32(2)) % 32)
		if l1 != 0 {
			v111 = l1 + v107
		} else {
			v111 = int32(0)
		}
		v114 = m.G111
		v115 = *(*int32)(unsafe.Add(mBase, uint32(v114)+4))
		m.T0[v115].(func(*base.Module, int32, int32, int32, int32))(m, l0+v107, v111, l2-v100, l3+v107)
		mBase = m.M
	}
	return
}
func F_PredictorSub2_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v16 int32
	_ = v16
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v41 base.V128
	_ = v41
	var v44 int32
	_ = v44
	var v49 base.V128
	_ = v49
	var v53 base.V128
	_ = v53
	var v54 base.V128
	_ = v54
	var v57 int32
	_ = v57
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v82 int32
	_ = v82
	var v85 int32
	_ = v85
	var v86 base.V128
	_ = v86
	var v89 base.V128
	_ = v89
	var v90 base.V128
	_ = v90
	var v98 int32
	_ = v98
	var v105 int32
	_ = v105
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v113 int32
	_ = v113
	v11 = int32(4)
	if l2 < v11 {
		v98 = int32(0)
	} else {
		v16 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v16) {
			v31 = int32(4)
			v32 = l0
			v34 = (int32(base.Ui32(v16)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			v35 = l1
			v36 = l3
			for {
				v37 = int32(0)
				v38 = base.Simd_g_v128_load(m, v32, v37)
				v40 = base.Simd_g_v128_load(m, v35, v37)
				v41 = base.Simd_g_i8x16_sub(v38, v40)
				base.Simd_g_v128_store(m, v36, v37, v41)
				v44 = int32(16)
				v49 = base.Simd_g_v128_load(m, v32+v44, v37)
				v53 = base.Simd_g_v128_load(m, v35+v44, v37)
				v54 = base.Simd_g_i8x16_sub(v49, v53)
				base.Simd_g_v128_store(m, v36+v44, v37, v54)
				v57 = int32(32)
				v64 = v31 + int32(8)
				v66 = v34 + int32(-2)
				if v66 != 0 {
					v31 = v64
					v32 = v32 + v57
					v34 = v66
					v35 = v35 + v57
					v36 = v36 + v57
					continue
				} else {
					break
				}
				break
			}
			v73 = v64
			v74 = v31 + int32(4)
		} else {
			v73 = v11
			v74 = int32(0)
		}
		if v16&int32(4) != 0 {
			v98 = v74
		} else {
			v82 = v74 << (uint(int32(2)) % 32)
			v85 = int32(0)
			v86 = base.Simd_g_v128_load(m, l0+v82, v85)
			v89 = base.Simd_g_v128_load(m, l1+v82, v85)
			v90 = base.Simd_g_i8x16_sub(v86, v89)
			base.Simd_g_v128_store(m, l3+v82, v85, v90)
			v98 = v73
		}
	}
	if l2 == v98 {
	} else {
		v105 = v98 << (uint(int32(2)) % 32)
		if l1 != 0 {
			v109 = l1 + v105
		} else {
			v109 = int32(0)
		}
		v112 = m.G111
		v113 = *(*int32)(unsafe.Add(mBase, uint32(v112)+8))
		m.T0[v113].(func(*base.Module, int32, int32, int32, int32))(m, l0+v105, v109, l2-v98, l3+v105)
		mBase = m.M
	}
	return
}
func F_PredictorSub3_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v43 base.V128
	_ = v43
	var v44 int32
	_ = v44
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v52 int32
	_ = v52
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v91 base.V128
	_ = v91
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v105 int32
	_ = v105
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	v12 = int32(4)
	if l2 < v12 {
		v105 = int32(0)
	} else {
		v17 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v17) {
			v33 = int32(4)
			v34 = int32(0)
			v36 = (int32(base.Ui32(v17)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v40 = l3 + v34
				v41 = l0 + v34
				v42 = int32(0)
				v43 = base.Simd_g_v128_load(m, v41, v42)
				v44 = l1 + v34
				v48 = base.Simd_g_v128_load(m, v44+int32(4), v42)
				v49 = base.Simd_g_i8x16_sub(v43, v48)
				base.Simd_g_v128_store(m, v40, v42, v49)
				v52 = int32(16)
				v57 = base.Simd_g_v128_load(m, v41+v52, v42)
				v61 = base.Simd_g_v128_load(m, v44+int32(20), v42)
				v62 = base.Simd_g_i8x16_sub(v57, v61)
				base.Simd_g_v128_store(m, v40+v52, v42, v62)
				v68 = v33 + int32(8)
				v70 = v36 + int32(-2)
				if v70 != 0 {
					v33 = v68
					v34 = v34 + int32(32)
					v36 = v70
					continue
				} else {
					break
				}
				break
			}
			v77 = v68
			v78 = v33 + int32(4)
		} else {
			v77 = v12
			v78 = int32(0)
		}
		if v17&int32(4) != 0 {
			v105 = v78
		} else {
			v87 = v78 << (uint(int32(2)) % 32)
			v90 = int32(0)
			v91 = base.Simd_g_v128_load(m, l0+v87, v90)
			v96 = base.Simd_g_v128_load(m, l1+v87+int32(4), v90)
			v97 = base.Simd_g_i8x16_sub(v91, v96)
			base.Simd_g_v128_store(m, l3+v87, v90, v97)
			v105 = v77
		}
	}
	if l2 == v105 {
	} else {
		v113 = v105 << (uint(int32(2)) % 32)
		if l1 != 0 {
			v117 = l1 + v113
		} else {
			v117 = int32(0)
		}
		v120 = m.G111
		v121 = *(*int32)(unsafe.Add(mBase, uint32(v120)+12))
		m.T0[v121].(func(*base.Module, int32, int32, int32, int32))(m, l0+v113, v117, l2-v105, l3+v113)
		mBase = m.M
	}
	return
}
func F_PredictorSub4_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v17 int32
	_ = v17
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v36 int32
	_ = v36
	var v40 int32
	_ = v40
	var v41 int32
	_ = v41
	var v42 int32
	_ = v42
	var v43 base.V128
	_ = v43
	var v44 int32
	_ = v44
	var v48 base.V128
	_ = v48
	var v49 base.V128
	_ = v49
	var v52 int32
	_ = v52
	var v57 base.V128
	_ = v57
	var v61 base.V128
	_ = v61
	var v62 base.V128
	_ = v62
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v87 int32
	_ = v87
	var v90 int32
	_ = v90
	var v91 base.V128
	_ = v91
	var v96 base.V128
	_ = v96
	var v97 base.V128
	_ = v97
	var v105 int32
	_ = v105
	var v113 int32
	_ = v113
	var v117 int32
	_ = v117
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	v12 = int32(4)
	if l2 < v12 {
		v105 = int32(0)
	} else {
		v17 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v17) {
			v33 = int32(4)
			v34 = int32(0)
			v36 = (int32(base.Ui32(v17)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v40 = l3 + v34
				v41 = l0 + v34
				v42 = int32(0)
				v43 = base.Simd_g_v128_load(m, v41, v42)
				v44 = l1 + v34
				v48 = base.Simd_g_v128_load(m, v44+int32(-4), v42)
				v49 = base.Simd_g_i8x16_sub(v43, v48)
				base.Simd_g_v128_store(m, v40, v42, v49)
				v52 = int32(16)
				v57 = base.Simd_g_v128_load(m, v41+v52, v42)
				v61 = base.Simd_g_v128_load(m, v44+int32(12), v42)
				v62 = base.Simd_g_i8x16_sub(v57, v61)
				base.Simd_g_v128_store(m, v40+v52, v42, v62)
				v68 = v33 + int32(8)
				v70 = v36 + int32(-2)
				if v70 != 0 {
					v33 = v68
					v34 = v34 + int32(32)
					v36 = v70
					continue
				} else {
					break
				}
				break
			}
			v77 = v68
			v78 = v33 + int32(4)
		} else {
			v77 = v12
			v78 = int32(0)
		}
		if v17&int32(4) != 0 {
			v105 = v78
		} else {
			v87 = v78 << (uint(int32(2)) % 32)
			v90 = int32(0)
			v91 = base.Simd_g_v128_load(m, l0+v87, v90)
			v96 = base.Simd_g_v128_load(m, l1+int32(-4)+v87, v90)
			v97 = base.Simd_g_i8x16_sub(v91, v96)
			base.Simd_g_v128_store(m, l3+v87, v90, v97)
			v105 = v77
		}
	}
	if l2 == v105 {
	} else {
		v113 = v105 << (uint(int32(2)) % 32)
		if l1 != 0 {
			v117 = l1 + v113
		} else {
			v117 = int32(0)
		}
		v120 = m.G111
		v121 = *(*int32)(unsafe.Add(mBase, uint32(v120)+16))
		m.T0[v121].(func(*base.Module, int32, int32, int32, int32))(m, l0+v113, v117, l2-v105, l3+v113)
		mBase = m.M
	}
	return
}
func F_PredictorSub5_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v42 base.V128
	_ = v42
	var v45 base.V128
	_ = v45
	var v47 base.V128
	_ = v47
	var v49 base.V128
	_ = v49
	var v54 base.V128
	_ = v54
	var v57 int32
	_ = v57
	var v66 int32
	_ = v66
	var v72 int32
	_ = v72
	var v83 int32
	_ = v83
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	if l2 < int32(4) {
		v72 = int32(0)
	} else {
		v25 = l0 + int32(-4)
		v26 = int32(0)
		v27 = l3
		v28 = l1
		for {
			v33 = int32(4)
			v35 = int32(0)
			v36 = base.Simd_g_v128_load_rng(m, v25+v33, v35, int32(-4), int32(20))
			v38 = base.Simd_g_v128_load_nc(m, v25, v35)
			v42 = base.Simd_g_v128_load_rng(m, v28+v33, v35, int32(-4), int32(20))
			v45 = base.Simd_g_const(&F_PredictorSub5_SSE2__k0)
			v47 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v38, v42), base.Simd_g_v128_and(base.Simd_g_v128_xor(v42, v38), v45))
			v49 = base.Simd_g_v128_load_nc(m, v28, v35)
			v54 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v36, base.Simd_g_i8x16_avgr_u(v47, v49)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v47, v49), v45))
			base.Simd_g_v128_store(m, v27, v35, v54)
			v57 = int32(16)
			v66 = v26 + v33
			if v26+int32(8) <= l2 {
				v25 = v25 + v57
				v26 = v66
				v27 = v27 + v57
				v28 = v28 + v57
				continue
			} else {
				break
			}
			break
		}
		v72 = v66
	}
	if l2 == v72 {
	} else {
		v83 = v72 << (uint(int32(2)) % 32)
		v88 = m.G111
		v89 = *(*int32)(unsafe.Add(mBase, uint32(v88)+20))
		m.T0[v89].(func(*base.Module, int32, int32, int32, int32))(m, l0+v83, l1+v83, l2-v72, l3+v83)
		mBase = m.M
	}
	return
}

var F_PredictorSub5_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorSub6_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v16 int32
	_ = v16
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 base.V128
	_ = v36
	var v38 base.V128
	_ = v38
	var v40 base.V128
	_ = v40
	var v46 base.V128
	_ = v46
	var v49 int32
	_ = v49
	var v58 int32
	_ = v58
	var v64 int32
	_ = v64
	var v74 int32
	_ = v74
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	if l2 < int32(4) {
		v64 = int32(0)
	} else {
		v16 = int32(-4)
		v26 = l1 + v16
		v27 = l0 + v16
		v28 = int32(0)
		v29 = l3
		for {
			v33 = int32(4)
			v35 = int32(0)
			v36 = base.Simd_g_v128_load_rng(m, v27+v33, v35, int32(-4), int32(20))
			v38 = base.Simd_g_v128_load_nc(m, v27, v35)
			v40 = base.Simd_g_v128_load(m, v26, v35)
			v46 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v36, base.Simd_g_i8x16_avgr_u(v38, v40)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v40, v38), base.Simd_g_const(&F_PredictorSub6_SSE2__k0)))
			base.Simd_g_v128_store(m, v29, v35, v46)
			v49 = int32(16)
			v58 = v28 + v33
			if v28+int32(8) <= l2 {
				v26 = v26 + v49
				v27 = v27 + v49
				v28 = v58
				v29 = v29 + v49
				continue
			} else {
				break
			}
			break
		}
		v64 = v58
	}
	if l2 == v64 {
	} else {
		v74 = v64 << (uint(int32(2)) % 32)
		v79 = m.G111
		v80 = *(*int32)(unsafe.Add(mBase, uint32(v79)+24))
		m.T0[v80].(func(*base.Module, int32, int32, int32, int32))(m, l0+v74, l1+v74, l2-v64, l3+v74)
		mBase = m.M
	}
	return
}

var F_PredictorSub6_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorSub7_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v53 base.V128
	_ = v53
	var v54 int32
	_ = v54
	var v56 base.V128
	_ = v56
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v65 int32
	_ = v65
	var v70 base.V128
	_ = v70
	var v74 base.V128
	_ = v74
	var v78 base.V128
	_ = v78
	var v83 base.V128
	_ = v83
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v115 base.V128
	_ = v115
	var v119 base.V128
	_ = v119
	var v122 base.V128
	_ = v122
	var v128 base.V128
	_ = v128
	var v136 int32
	_ = v136
	var v147 int32
	_ = v147
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	v15 = int32(4)
	if l2 < v15 {
		v136 = int32(0)
	} else {
		v20 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v20) {
			v36 = int32(4)
			v37 = int32(0)
			v39 = (int32(base.Ui32(v20)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v46 = l3 + v37
				v47 = l0 + v37
				v48 = int32(0)
				v49 = base.Simd_g_v128_load(m, v47, v48)
				v53 = base.Simd_g_v128_load(m, v47+int32(-4), v48)
				v54 = l1 + v37
				v56 = base.Simd_g_v128_load(m, v54, v48)
				v60 = base.Simd_g_const(&F_PredictorSub7_SSE2__k0)
				v62 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v49, base.Simd_g_i8x16_avgr_u(v53, v56)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v56, v53), v60))
				base.Simd_g_v128_store(m, v46, v48, v62)
				v65 = int32(16)
				v70 = base.Simd_g_v128_load_rng(m, v47+v65, v48, int32(-4), int32(20))
				v74 = base.Simd_g_v128_load_nc(m, v47+int32(12), v48)
				v78 = base.Simd_g_v128_load(m, v54+v65, v48)
				v83 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v70, base.Simd_g_i8x16_avgr_u(v74, v78)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v78, v74), v60))
				base.Simd_g_v128_store(m, v46+v65, v48, v83)
				v89 = v36 + int32(8)
				v91 = v39 + int32(-2)
				if v91 != 0 {
					v36 = v89
					v37 = v37 + int32(32)
					v39 = v91
					continue
				} else {
					break
				}
				break
			}
			v98 = v89
			v99 = v36 + int32(4)
		} else {
			v98 = v15
			v99 = int32(0)
		}
		if v20&int32(4) != 0 {
			v136 = v99
		} else {
			v111 = v99 << (uint(int32(2)) % 32)
			v113 = l0 + v111
			v114 = int32(0)
			v115 = base.Simd_g_v128_load(m, v113, v114)
			v119 = base.Simd_g_v128_load(m, v113+int32(-4), v114)
			v122 = base.Simd_g_v128_load(m, l1+v111, v114)
			v128 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v115, base.Simd_g_i8x16_avgr_u(v119, v122)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v122, v119), base.Simd_g_const(&F_PredictorSub7_SSE2__k0)))
			base.Simd_g_v128_store(m, l3+v111, v114, v128)
			v136 = v98
		}
	}
	if l2 == v136 {
	} else {
		v147 = v136 << (uint(int32(2)) % 32)
		v152 = m.G111
		v153 = *(*int32)(unsafe.Add(mBase, uint32(v152)+28))
		m.T0[v153].(func(*base.Module, int32, int32, int32, int32))(m, l0+v147, l1+v147, l2-v136, l3+v147)
		mBase = m.M
	}
	return
}

var F_PredictorSub7_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorSub8_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v46 int32
	_ = v46
	var v47 int32
	_ = v47
	var v48 int32
	_ = v48
	var v49 base.V128
	_ = v49
	var v50 int32
	_ = v50
	var v54 base.V128
	_ = v54
	var v56 base.V128
	_ = v56
	var v60 base.V128
	_ = v60
	var v62 base.V128
	_ = v62
	var v65 int32
	_ = v65
	var v70 base.V128
	_ = v70
	var v74 base.V128
	_ = v74
	var v78 base.V128
	_ = v78
	var v83 base.V128
	_ = v83
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v111 int32
	_ = v111
	var v114 int32
	_ = v114
	var v115 base.V128
	_ = v115
	var v116 int32
	_ = v116
	var v120 base.V128
	_ = v120
	var v122 base.V128
	_ = v122
	var v128 base.V128
	_ = v128
	var v136 int32
	_ = v136
	var v147 int32
	_ = v147
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	v15 = int32(4)
	if l2 < v15 {
		v136 = int32(0)
	} else {
		v20 = l2 + int32(-4)
		if base.Ui32(int32(4)) <= base.Ui32(v20) {
			v36 = int32(4)
			v37 = int32(0)
			v39 = (int32(base.Ui32(v20)>>(uint(int32(2))%32)) + int32(1)) & int32(2147483646)
			for {
				v46 = l3 + v37
				v47 = l0 + v37
				v48 = int32(0)
				v49 = base.Simd_g_v128_load(m, v47, v48)
				v50 = l1 + v37
				v54 = base.Simd_g_v128_load(m, v50+int32(-4), v48)
				v56 = base.Simd_g_v128_load(m, v50, v48)
				v60 = base.Simd_g_const(&F_PredictorSub8_SSE2__k0)
				v62 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v49, base.Simd_g_i8x16_avgr_u(v54, v56)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v56, v54), v60))
				base.Simd_g_v128_store(m, v46, v48, v62)
				v65 = int32(16)
				v70 = base.Simd_g_v128_load(m, v47+v65, v48)
				v74 = base.Simd_g_v128_load_rng(m, v50+int32(12), v48, int32(0), int32(20))
				v78 = base.Simd_g_v128_load_nc(m, v50+v65, v48)
				v83 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v70, base.Simd_g_i8x16_avgr_u(v74, v78)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v78, v74), v60))
				base.Simd_g_v128_store(m, v46+v65, v48, v83)
				v89 = v36 + int32(8)
				v91 = v39 + int32(-2)
				if v91 != 0 {
					v36 = v89
					v37 = v37 + int32(32)
					v39 = v91
					continue
				} else {
					break
				}
				break
			}
			v98 = v89
			v99 = v36 + int32(4)
		} else {
			v98 = v15
			v99 = int32(0)
		}
		if v20&int32(4) != 0 {
			v136 = v99
		} else {
			v111 = v99 << (uint(int32(2)) % 32)
			v114 = int32(0)
			v115 = base.Simd_g_v128_load(m, l0+v111, v114)
			v116 = l1 + v111
			v120 = base.Simd_g_v128_load(m, v116+int32(-4), v114)
			v122 = base.Simd_g_v128_load(m, v116, v114)
			v128 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v115, base.Simd_g_i8x16_avgr_u(v120, v122)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v122, v120), base.Simd_g_const(&F_PredictorSub8_SSE2__k0)))
			base.Simd_g_v128_store(m, l3+v111, v114, v128)
			v136 = v98
		}
	}
	if l2 == v136 {
	} else {
		v147 = v136 << (uint(int32(2)) % 32)
		v152 = m.G111
		v153 = *(*int32)(unsafe.Add(mBase, uint32(v152)+32))
		m.T0[v153].(func(*base.Module, int32, int32, int32, int32))(m, l0+v147, l1+v147, l2-v136, l3+v147)
		mBase = m.M
	}
	return
}

var F_PredictorSub8_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

func F_PredictorSub9_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) {
	mBase := m.M
	_ = mBase
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	var v30 base.V128
	_ = v30
	var v32 base.V128
	_ = v32
	var v33 int32
	_ = v33
	var v36 base.V128
	_ = v36
	var v42 base.V128
	_ = v42
	var v45 int32
	_ = v45
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	if l2 < int32(4) {
		v60 = int32(0)
	} else {
		v22 = int32(0)
		v23 = l3
		v24 = l0
		v25 = l1
		for {
			v29 = int32(0)
			v30 = base.Simd_g_v128_load(m, v24, v29)
			v32 = base.Simd_g_v128_load_rng(m, v25, v29, int32(0), int32(20))
			v33 = int32(4)
			v36 = base.Simd_g_v128_load_nc(m, v25+v33, v29)
			v42 = base.Simd_g_i8x16_add(base.Simd_g_i8x16_sub(v30, base.Simd_g_i8x16_avgr_u(v32, v36)), base.Simd_g_v128_and(base.Simd_g_v128_xor(v36, v32), base.Simd_g_const(&F_PredictorSub9_SSE2__k0)))
			base.Simd_g_v128_store(m, v23, v29, v42)
			v45 = int32(16)
			v54 = v22 + v33
			if v22+int32(8) <= l2 {
				v22 = v54
				v23 = v23 + v45
				v24 = v24 + v45
				v25 = v25 + v45
				continue
			} else {
				break
			}
			break
		}
		v60 = v54
	}
	if l2 == v60 {
	} else {
		v70 = v60 << (uint(int32(2)) % 32)
		v75 = m.G111
		v76 = *(*int32)(unsafe.Add(mBase, uint32(v75)+36))
		m.T0[v76].(func(*base.Module, int32, int32, int32, int32))(m, l0+v70, l1+v70, l2-v60, l3+v70)
		mBase = m.M
	}
	return
}

var F_PredictorSub9_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}

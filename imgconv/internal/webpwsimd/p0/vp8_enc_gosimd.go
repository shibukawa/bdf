//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8EncAnalyze(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v27 int32
	_ = v27
	var v29 int32
	_ = v29
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
	var v41 int32
	_ = v41
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v54 int32
	_ = v54
	var v58 int32
	_ = v58
	var v61 int32
	_ = v61
	var v66 int32
	_ = v66
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v75 int32
	_ = v75
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v85 int32
	_ = v85
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v106 int64
	_ = v106
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v111 int32
	_ = v111
	var v115 int32
	_ = v115
	var v120 int32
	_ = v120
	var v121 int64
	_ = v121
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v168 int32
	_ = v168
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v178 int32
	_ = v178
	var v179 int32
	_ = v179
	var v181 int64
	_ = v181
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v190 int32
	_ = v190
	var v191 int32
	_ = v191
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v207 int32
	_ = v207
	var v211 int32
	_ = v211
	var v216 int32
	_ = v216
	var v217 int32
	_ = v217
	var v220 int32
	_ = v220
	var v224 int32
	_ = v224
	var v225 int64
	_ = v225
	var v231 int32
	_ = v231
	var v234 int32
	_ = v234
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v375 int32
	_ = v375
	var v376 int32
	_ = v376
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v384 int32
	_ = v384
	var v388 int32
	_ = v388
	var v390 int32
	_ = v390
	var v393 int32
	_ = v393
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v404 int32
	_ = v404
	var v412 int32
	_ = v412
	var v417 int32
	_ = v417
	var v436 int32
	_ = v436
	var v439 int32
	_ = v439
	var v442 int32
	_ = v442
	var v443 int32
	_ = v443
	var v447 int32
	_ = v447
	var v458 int32
	_ = v458
	var v467 int32
	_ = v467
	var v469 int32
	_ = v469
	var v491 int32
	_ = v491
	var v495 int32
	_ = v495
	var v499 int32
	_ = v499
	var v528 int32
	_ = v528
	var v534 int32
	_ = v534
	var v539 int32
	_ = v539
	var v557 int32
	_ = v557
	var v559 int32
	_ = v559
	var v563 int32
	_ = v563
	var v565 int32
	_ = v565
	var v567 int32
	_ = v567
	var v577 int32
	_ = v577
	var v578 int32
	_ = v578
	var v600 int32
	_ = v600
	var v604 int32
	_ = v604
	var v605 int32
	_ = v605
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v608 int32
	_ = v608
	var v611 int32
	_ = v611
	var v613 int32
	_ = v613
	var v621 int32
	_ = v621
	var v630 int32
	_ = v630
	var v631 int32
	_ = v631
	var v632 int32
	_ = v632
	var v638 int32
	_ = v638
	var v639 int32
	_ = v639
	var v646 int32
	_ = v646
	var v673 int32
	_ = v673
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v685 int32
	_ = v685
	var v686 int32
	_ = v686
	var v707 int32
	_ = v707
	var v708 int32
	_ = v708
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v717 int32
	_ = v717
	var v722 int32
	_ = v722
	var v723 int32
	_ = v723
	var v724 int32
	_ = v724
	var v753 int64
	_ = v753
	var v757 int32
	_ = v757
	var v758 int32
	_ = v758
	var v763 int32
	_ = v763
	var v769 int32
	_ = v769
	var v776 int32
	_ = v776
	var v777 int32
	_ = v777
	var v793 int32
	_ = v793
	var v794 int32
	_ = v794
	var v800 int32
	_ = v800
	var v801 int32
	_ = v801
	var v806 base.V128
	_ = v806
	var v815 int32
	_ = v815
	var v816 int32
	_ = v816
	var v829 base.V128
	_ = v829
	var v837 base.V128
	_ = v837
	var v838 int32
	_ = v838
	var v842 int32
	_ = v842
	var v844 int32
	_ = v844
	var v848 int32
	_ = v848
	var v851 int32
	_ = v851
	var v855 int32
	_ = v855
	var v858 int32
	_ = v858
	var v862 int32
	_ = v862
	var v865 base.V128
	_ = v865
	var v873 int32
	_ = v873
	var v886 int32
	_ = v886
	var v889 int32
	_ = v889
	var v913 int32
	_ = v913
	var v914 int32
	_ = v914
	var v915 int32
	_ = v915
	var v936 int32
	_ = v936
	var v943 int32
	_ = v943
	var v970 int32
	_ = v970
	var v971 int32
	_ = v971
	var v973 int32
	_ = v973
	var v998 int32
	_ = v998
	var v1002 int32
	_ = v1002
	var v1003 int32
	_ = v1003
	var v1014 int32
	_ = v1014
	var v1036 int32
	_ = v1036
	var v1039 int32
	_ = v1039
	var v1040 int32
	_ = v1040
	var v1048 int32
	_ = v1048
	var v1049 int32
	_ = v1049
	var v1084 int32
	_ = v1084
	var v1085 int32
	_ = v1085
	var v1090 int64
	_ = v1090
	var v1093 int32
	_ = v1093
	var v1094 int32
	_ = v1094
	var v1106 int32
	_ = v1106
	var v1125 int32
	_ = v1125
	var v1126 int32
	_ = v1126
	var v1137 int32
	_ = v1137
	var v1159 int32
	_ = v1159
	var v1162 int32
	_ = v1162
	var v1163 int32
	_ = v1163
	var v1171 int32
	_ = v1171
	var v1172 int32
	_ = v1172
	var v1207 int32
	_ = v1207
	var v1208 int32
	_ = v1208
	var v1213 int64
	_ = v1213
	var v1216 int32
	_ = v1216
	var v1217 int32
	_ = v1217
	var v1229 int32
	_ = v1229
	var v1251 int32
	_ = v1251
	var v1252 int32
	_ = v1252
	var v1275 int32
	_ = v1275
	var v1277 int32
	_ = v1277
	var v1281 int32
	_ = v1281
	var v1283 int32
	_ = v1283
	var v1285 int32
	_ = v1285
	var v1295 int32
	_ = v1295
	var v1296 int32
	_ = v1296
	var v1318 int32
	_ = v1318
	var v1322 int32
	_ = v1322
	var v1323 int32
	_ = v1323
	var v1324 int32
	_ = v1324
	var v1325 int32
	_ = v1325
	var v1326 int32
	_ = v1326
	var v1329 int32
	_ = v1329
	var v1331 int32
	_ = v1331
	var v1338 int32
	_ = v1338
	var v1348 int32
	_ = v1348
	var v1349 int32
	_ = v1349
	var v1350 int32
	_ = v1350
	var v1356 int32
	_ = v1356
	var v1357 int32
	_ = v1357
	var v1364 int32
	_ = v1364
	var v1415 int32
	_ = v1415
	var v1426 int32
	_ = v1426
	var v1427 int32
	_ = v1427
	var v1429 int32
	_ = v1429
	var v1434 int32
	_ = v1434
	var v1435 int32
	_ = v1435
	var v1437 int32
	_ = v1437
	var v1438 int32
	_ = v1438
	var v1450 int32
	_ = v1450
	var v1453 int32
	_ = v1453
	var v1454 int32
	_ = v1454
	var v1456 int32
	_ = v1456
	var v1458 int32
	_ = v1458
	var v1460 int32
	_ = v1460
	var v1462 int32
	_ = v1462
	var v1469 int32
	_ = v1469
	var v1470 int32
	_ = v1470
	var v1472 int32
	_ = v1472
	var v1474 int32
	_ = v1474
	var v1481 int32
	_ = v1481
	var v1485 int32
	_ = v1485
	var v1490 int32
	_ = v1490
	var v1492 int32
	_ = v1492
	var v1514 int32
	_ = v1514
	var v1519 int32
	_ = v1519
	var v1520 int32
	_ = v1520
	var v1521 int32
	_ = v1521
	var v1524 int32
	_ = v1524
	var v1529 int32
	_ = v1529
	var v1530 int32
	_ = v1530
	var v1552 int32
	_ = v1552
	var v1553 int32
	_ = v1553
	var v1556 int32
	_ = v1556
	var v1557 int32
	_ = v1557
	var v1558 int32
	_ = v1558
	var v1559 int32
	_ = v1559
	var v1562 int32
	_ = v1562
	var v1567 int32
	_ = v1567
	var v1570 int32
	_ = v1570
	var v1577 int32
	_ = v1577
	var v1582 int32
	_ = v1582
	var v1583 int32
	_ = v1583
	var v1584 int32
	_ = v1584
	var v1585 int32
	_ = v1585
	var v1589 int32
	_ = v1589
	var v1592 int32
	_ = v1592
	var v1611 int32
	_ = v1611
	var v1615 int32
	_ = v1615
	var v1616 int32
	_ = v1616
	var v1621 int64
	_ = v1621
	var v1622 int32
	_ = v1622
	var v1629 int64
	_ = v1629
	var v1630 int32
	_ = v1630
	var v1631 int64
	_ = v1631
	var v1641 int32
	_ = v1641
	var v1643 int32
	_ = v1643
	var v1652 int32
	_ = v1652
	var v1654 int32
	_ = v1654
	var v1659 int32
	_ = v1659
	var v1660 int32
	_ = v1660
	var v1671 int32
	_ = v1671
	var v1673 int32
	_ = v1673
	var v1676 int32
	_ = v1676
	var v1694 int32
	_ = v1694
	var v1696 int32
	_ = v1696
	var v1717 base.V128
	_ = v1717
	var v1718 int32
	_ = v1718
	var v1721 int32
	_ = v1721
	var v1722 int32
	_ = v1722
	var v1723 int32
	_ = v1723
	var v1724 int32
	_ = v1724
	var v1725 int32
	_ = v1725
	var v1727 int32
	_ = v1727
	var v1729 int32
	_ = v1729
	var v1730 int32
	_ = v1730
	var v1731 int32
	_ = v1731
	var v1736 int32
	_ = v1736
	var v1738 int32
	_ = v1738
	var v1743 int32
	_ = v1743
	var v1744 int32
	_ = v1744
	var v1750 int32
	_ = v1750
	var v1752 int32
	_ = v1752
	var v1757 int32
	_ = v1757
	var v1758 int32
	_ = v1758
	var v1764 int32
	_ = v1764
	var v1765 int32
	_ = v1765
	var v1770 int32
	_ = v1770
	var v1771 int32
	_ = v1771
	var v1779 int32
	_ = v1779
	var v1784 int32
	_ = v1784
	var v1785 int32
	_ = v1785
	var v1791 int32
	_ = v1791
	var v1792 int32
	_ = v1792
	var v1797 int32
	_ = v1797
	var v1798 int32
	_ = v1798
	var v1804 int32
	_ = v1804
	var v1809 int32
	_ = v1809
	var v1814 int32
	_ = v1814
	var v1815 int32
	_ = v1815
	var v1823 int32
	_ = v1823
	var v1828 int32
	_ = v1828
	var v1829 int32
	_ = v1829
	var v1833 int32
	_ = v1833
	var v1837 int32
	_ = v1837
	var v1841 int32
	_ = v1841
	var v1845 int32
	_ = v1845
	var v1850 int32
	_ = v1850
	var v1853 int32
	_ = v1853
	var v1854 int32
	_ = v1854
	var v1860 int32
	_ = v1860
	var v1865 int32
	_ = v1865
	var v1867 int32
	_ = v1867
	var v1872 int32
	_ = v1872
	var v1892 int32
	_ = v1892
	var v1893 int32
	_ = v1893
	var v1895 int32
	_ = v1895
	var v1914 int32
	_ = v1914
	var v1915 int32
	_ = v1915
	var v1936 int32
	_ = v1936
	var v1939 int32
	_ = v1939
	var v1940 int32
	_ = v1940
	var v1941 int32
	_ = v1941
	var v1943 int32
	_ = v1943
	var v1944 int32
	_ = v1944
	var v1945 int32
	_ = v1945
	var v1947 int32
	_ = v1947
	var v1949 int32
	_ = v1949
	var v1951 int32
	_ = v1951
	var v1952 int32
	_ = v1952
	var v1953 int32
	_ = v1953
	var v1958 int32
	_ = v1958
	var v1963 int32
	_ = v1963
	var v1968 int32
	_ = v1968
	var v1976 int32
	_ = v1976
	var v2000 int32
	_ = v2000
	var v2002 int32
	_ = v2002
	var v2005 int32
	_ = v2005
	var v2006 int32
	_ = v2006
	var v2010 int32
	_ = v2010
	var v2015 int32
	_ = v2015
	var v2022 int32
	_ = v2022
	var v2077 int32
	_ = v2077
	var v2078 int32
	_ = v2078
	var v2087 int32
	_ = v2087
	var v2088 base.V128
	_ = v2088
	var v2092 int32
	_ = v2092
	var v2093 int32
	_ = v2093
	var v2105 base.V128
	_ = v2105
	var v2108 base.V128
	_ = v2108
	var v2116 base.V128
	_ = v2116
	var v2117 base.V128
	_ = v2117
	var v2118 base.V128
	_ = v2118
	var v2122 int32
	_ = v2122
	var v2123 base.V128
	_ = v2123
	var v2125 base.V128
	_ = v2125
	var v2126 base.V128
	_ = v2126
	var v2129 int32
	_ = v2129
	var v2130 int32
	_ = v2130
	var v2133 base.V128
	_ = v2133
	var v2138 int32
	_ = v2138
	var v2142 int32
	_ = v2142
	var v2143 int32
	_ = v2143
	var v2145 int32
	_ = v2145
	var v2174 int32
	_ = v2174
	var v2175 int32
	_ = v2175
	var v2177 int32
	_ = v2177
	var v2185 int32
	_ = v2185
	var v2198 int32
	_ = v2198
	var v2200 int32
	_ = v2200
	var v2202 int32
	_ = v2202
	var v2206 int32
	_ = v2206
	var v2210 int32
	_ = v2210
	var v2212 int32
	_ = v2212
	var v2233 int32
	_ = v2233
	var v2239 int32
	_ = v2239
	var v2240 int32
	_ = v2240
	var v2245 base.V128
	_ = v2245
	var v2250 int32
	_ = v2250
	var v2254 int32
	_ = v2254
	var v2255 int32
	_ = v2255
	var v2264 int32
	_ = v2264
	var v2279 int32
	_ = v2279
	var v2280 base.V128
	_ = v2280
	var v2282 base.V128
	_ = v2282
	var v2283 base.V128
	_ = v2283
	var v2287 int32
	_ = v2287
	var v2288 int32
	_ = v2288
	var v2290 int32
	_ = v2290
	var v2293 int32
	_ = v2293
	var v2294 int32
	_ = v2294
	var v2297 int32
	_ = v2297
	var v2300 int32
	_ = v2300
	var v2301 int32
	_ = v2301
	var v2304 int32
	_ = v2304
	var v2307 int32
	_ = v2307
	var v2308 int32
	_ = v2308
	var v2314 base.V128
	_ = v2314
	var v2336 base.V128
	_ = v2336
	var v2339 int32
	_ = v2339
	var v2343 int32
	_ = v2343
	var v2348 int32
	_ = v2348
	var v2353 int32
	_ = v2353
	var v2358 base.V128
	_ = v2358
	var v2382 int32
	_ = v2382
	var v2390 int32
	_ = v2390
	var v2424 int32
	_ = v2424
	var v2425 int32
	_ = v2425
	var v2434 int32
	_ = v2434
	var v2447 int32
	_ = v2447
	var v2449 int32
	_ = v2449
	var v2451 int32
	_ = v2451
	var v2455 int32
	_ = v2455
	var v2456 int32
	_ = v2456
	var v2459 int32
	_ = v2459
	var v2466 int32
	_ = v2466
	var v2467 int32
	_ = v2467
	var v2470 int32
	_ = v2470
	var v2471 int32
	_ = v2471
	var v2474 int32
	_ = v2474
	var v2481 int32
	_ = v2481
	var v2489 int32
	_ = v2489
	v2 = int32(0)
	v27 = m.G0
	v29 = v27 - int32(_a_F_VP8EncAnalyze_0)
	m.G0 = v29
	v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v32 = *(*int32)(unsafe.Add(mBase, uint32(v31)+80))
	if v32 != 0 {
		goto L6
	} else {
		goto L7
	}
L1:
	;
	m.G0 = v29 + int32(_a_F_VP8EncAnalyze_0)
	return v2489
L2:
	;
	v1519 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v1520 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v1521 = v1519 * v1520
	if v1521 < int32(1) {
		v1589 = v1519
		v1592 = v1521
		v1611 = v1520
		goto L156
	} else {
		goto L157
	}
L3:
	;
	v793 = int32(1)
	v794 = v404 << (uint(v793) % 32)
	if base.Ui32(v404) <= base.Ui32(int32(3)) {
		v886 = int32(0)
		v889 = v793
		goto L90
	} else {
		goto L91
	}
L4:
	;
	v769 = int32(255)
	v776 = int32(256)
	v777 = int32(-1)
	goto L3
L5:
	;
	v673 = int32(1)
	v674 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v675 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	if v674*v675 < v673 {
		goto L84
	} else {
		goto L85
	}
L6:
	;
	v39 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v40 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v41 = m.G1
	goto L10
L7:
	;
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if int32(1) < v33 {
		goto L6
	} else {
		goto L8
	}
L8:
	;
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[0])))
	if int32(1) < v36 {
		goto L5
	} else {
		goto L9
	}
L9:
	;
	goto L6
L10:
	;
	v46 = m.G1
	goto L11
L11:
	;
	v49 = *(*int32)(unsafe.Add(mBase, uint32(v46)+uint32(_c_F_VP8EncAnalyze[1])))
	m.T0[v49].(func(*base.Module, int32))(m, v29+int32(8))
	mBase = m.M
	v52 = v29 + int32(1064)
	*(*int32)(unsafe.Add(mBase, uint32(v29)+24)) = v52
	v54 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v29)+16)) = v54 + int32(255)
	v58 = int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v29)+20)) = v29 + v58
	v61 = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(v52))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+24)) = l0
	v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+280)) = v66
	v68 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+296)) = v68
	v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+304)) = v70
	*(*int32)(unsafe.Add(mBase, uint32(v52)+32)) = l0 + int32(88)
	v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[4])))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+36)) = v75
	v79 = int32(-32)
	v80 = (v29 + int32(1511)) & v79
	*(*int32)(unsafe.Add(mBase, uint32(v52)+8)) = v80
	v85 = (v29 + int32(1424)) & v79
	*(*int32)(unsafe.Add(mBase, uint32(v52)+308)) = v85
	*(*int32)(unsafe.Add(mBase, uint32(v52)+20)) = v80 + int32(1536)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+16)) = v80 + int32(1024)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+12)) = v80 + int32(512)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+316)) = v85 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+312)) = v85 + int32(32)
	v102 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+40)) = v102
	v104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+28)) = v104
	v106 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v52)+320)) = v106
	v108 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v85)+47)) = uint8(v108)
	v110 = *(*int32)(unsafe.Add(mBase, uint32(v52)+312))
	v111 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v110+v111))) = uint8(v108)
	v115 = *(*int32)(unsafe.Add(mBase, uint32(v52)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v115+v111))) = uint8(v108)
	v120 = *(*int32)(unsafe.Add(mBase, uint32(v52)+308))
	v121 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v120))) = v121
	*(*int64)(unsafe.Add(mBase, uint32(v120+v58))) = v121
	v127 = *(*int32)(unsafe.Add(mBase, uint32(v52)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v127))) = v121
	v130 = *(*int32)(unsafe.Add(mBase, uint32(v52)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v130))) = v121
	*(*int32)(unsafe.Add(mBase, uint32(v52)+160)) = v61
	v135 = *(*int32)(unsafe.Add(mBase, uint32(v52)+304))
	if v135 == v61 {
		goto L13
	} else {
		goto L14
	}
L12:
	;
	v173 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+4)) = v173
	*(*int32)(unsafe.Add(mBase, uint32(v52))) = v173
	v178 = *(*int32)(unsafe.Add(mBase, uint32(v52)+24))
	v179 = *(*int32)(unsafe.Add(mBase, uint32(v178)+uint32(_c_F_VP8EncAnalyze[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+40)) = v179
	v181 = *(*int64)(unsafe.Add(mBase, uint32(v178)+uint32(_c_F_VP8EncAnalyze[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v52)+320)) = v181
	v183 = *(*int32)(unsafe.Add(mBase, uint32(v178)+uint32(_c_F_VP8EncAnalyze[4])))
	v184 = *(*int32)(unsafe.Add(mBase, uint32(v178)+48))
	v186 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+36)) = v183 + v173*v184<<(uint(v186)%32)
	v190 = *(*int32)(unsafe.Add(mBase, uint32(v178)+uint32(_c_F_VP8EncAnalyze[6])))
	v191 = *(*int32)(unsafe.Add(mBase, uint32(v178)+40))
	*(*int32)(unsafe.Add(mBase, uint32(v52)+28)) = v190 + v191*v173<<(uint(v186)%32)
	v197 = *(*int32)(unsafe.Add(mBase, uint32(v178)+52))
	v198 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+32)) = v178 + (v197+v198)&v173<<(uint(int32(5))%32) + int32(88)
	v207 = *(*int32)(unsafe.Add(mBase, uint32(v52)+316))
	v211 = int32(127)
	goto L19
L13:
	;
	v140 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v142 = v140 * v141
	*(*int32)(unsafe.Add(mBase, uint32(v52)+292)) = v142
	*(*int32)(unsafe.Add(mBase, uint32(v52)+288)) = v142
	v145 = *(*int32)(unsafe.Add(mBase, uint32(v52)+24))
	v146 = *(*int32)(unsafe.Add(mBase, uint32(v145)+uint32(_c_F_VP8EncAnalyze[7])))
	v148 = *(*int32)(unsafe.Add(mBase, uint32(v145)+40))
	v151 = F_memset(m, v146, int32(127), v148<<(uint(int32(5))%32))
	mBase = m.M
	v152 = *(*int32)(unsafe.Add(mBase, uint32(v145)+uint32(_c_F_VP8EncAnalyze[5])))
	v153 = int32(0)
	v154 = *(*int32)(unsafe.Add(mBase, uint32(v145)+40))
	v157 = F_memset(m, v152, v153, v154<<(uint(int32(2))%32))
	mBase = m.M
	v158 = *(*int32)(unsafe.Add(mBase, uint32(v145)+uint32(_c_F_VP8EncAnalyze[3])))
	if v158 == v153 {
		goto L15
	} else {
		goto L16
	}
L14:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v52)+300)) = int32(0)
	goto L13
L15:
	;
	v168 = int32(0)
	v170 = F_memset(m, v29+int32(1232), v168, int32(96))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v52)+284)) = v168
	goto L12
L16:
	;
	v162 = *(*int32)(unsafe.Add(mBase, uint32(v145)+40))
	v165 = F_memset(m, v158, int32(0), v162<<(uint(int32(2))%32))
	mBase = m.M
	goto L15
L17:
	;
	v244 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v245 = v40 * v244
	*(*int32)(unsafe.Add(mBase, uint32(v52)+288)) = v245
	*(*int32)(unsafe.Add(mBase, uint32(v52)+292)) = v245
	goto L23
L19:
	;
	goto L20
L20:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v207+v198))) = uint8(v211)
	v216 = *(*int32)(unsafe.Add(mBase, uint32(v52)+312))
	v217 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v216+v217))) = uint8(v211)
	v220 = *(*int32)(unsafe.Add(mBase, uint32(v52)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v220+v217))) = uint8(v211)
	v224 = *(*int32)(unsafe.Add(mBase, uint32(v52)+308))
	v225 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v224))) = v225
	*(*int64)(unsafe.Add(mBase, uint32(v224+int32(8)))) = v225
	v231 = *(*int32)(unsafe.Add(mBase, uint32(v52)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v231))) = v225
	v234 = *(*int32)(unsafe.Add(mBase, uint32(v52)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v234))) = v225
	v237 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v52)+160)) = v237
	v239 = *(*int32)(unsafe.Add(mBase, uint32(v52)+304))
	if v239 == v237 {
		goto L21
	} else {
		goto L22
	}
L21:
	;
	goto L17
L22:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v52)+300)) = int32(0)
	goto L21
L23:
	;
	v249 = v29 + int32(32)
	goto L26
L24:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[8]))) = int32(20)
	v375 = v29 + int32(8)
	v376 = *(*int32)(unsafe.Add(mBase, uint32(v41)+uint32(_c_F_VP8EncAnalyze[9])))
	m.T0[v376].(func(*base.Module, int32))(m, v375)
	mBase = m.M
	v380 = *(*int32)(unsafe.Add(mBase, uint32(v41)+uint32(_c_F_VP8EncAnalyze[10])))
	v381 = m.T0[v380].(func(*base.Module, int32) int32)(m, v375)
	mBase = m.M
	v384 = *(*int32)(unsafe.Add(mBase, uint32(v41)+uint32(_c_F_VP8EncAnalyze[11])))
	m.T0[v384].(func(*base.Module, int32))(m, v375)
	mBase = m.M
	if v381&int32(1) != 0 {
		goto L38
	} else {
		goto L39
	}
L26:
	;
	base.MemoryFill(m, v249, int32(0), int32(1032))
	goto L24
L38:
	;
	v393 = *(*int32)(unsafe.Add(mBase, uint32(v29)+1056))
	v394 = v39 * v40
	v395 = base.I32_div_s(v393, v394)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3388)) = v395
	v397 = *(*int32)(unsafe.Add(mBase, uint32(v29)+1060))
	v398 = base.I32_div_s(v397, v394)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3392)) = v398
	v400 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v401 = int32(4)
	if v400 < v401 {
		goto L43
	} else {
		goto L44
	}
L39:
	;
	v388 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v390 = *(*int32)(unsafe.Add(mBase, uint32(v388)+92))
	if v390 != 0 {
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v2489 = int32(0)
	goto L1
L41:
	;
	goto L40
L42:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v388)+92)) = int32(1)
	goto L41
L43:
	;
	v404 = v400
	goto L45
L44:
	;
	v404 = v401
	goto L45
L45:
	;
	v412 = v29 + int32(44)
	v417 = int32(3)
	goto L50
L46:
	;
	if base.Ui32(int32(254)) < base.Ui32(v458) {
		v499 = int32(255)
		goto L58
	} else {
		goto L59
	}
L47:
	;
	v458 = v417 + int32(-1)
	goto L46
L48:
	;
	v458 = v417 + int32(-2)
	goto L46
L49:
	;
	v458 = v417 + int32(-3)
	goto L46
L50:
	;
	v436 = *(*int32)(unsafe.Add(mBase, uint32(v412+int32(-12))))
	if v436 != 0 {
		goto L49
	} else {
		goto L52
	}
L51:
	;
	if int32(0) < v400 {
		goto L4
	} else {
		goto L57
	}
L52:
	;
	v439 = *(*int32)(unsafe.Add(mBase, uint32(v412+int32(-8))))
	if v439 != 0 {
		goto L48
	} else {
		goto L53
	}
L53:
	;
	v442 = *(*int32)(unsafe.Add(mBase, uint32(v412+int32(-4))))
	if v442 != 0 {
		goto L47
	} else {
		goto L54
	}
L54:
	;
	v443 = *(*int32)(unsafe.Add(mBase, uint32(v412)))
	if v443 != 0 {
		v458 = v417
		goto L46
	} else {
		goto L55
	}
L55:
	;
	v447 = v417 + int32(4)
	if v447 != int32(259) {
		v412 = v412 + int32(16)
		v417 = v447
		goto L50
	} else {
		goto L56
	}
L56:
	;
	goto L51
L57:
	;
	v1514 = v2
	goto L2
L58:
	;
	if v400 <= int32(0) {
		goto L64
	} else {
		goto L65
	}
L59:
	;
	v467 = int32(255)
	v469 = v29 + int32(1052)
	goto L60
L60:
	;
	v491 = *(*int32)(unsafe.Add(mBase, uint32(v469)))
	if v491 != 0 {
		v499 = v467
		goto L58
	} else {
		goto L62
	}
L61:
	;
	v499 = v458
	goto L58
L62:
	;
	v495 = v467 + int32(-1)
	if base.Ui32(v458) < base.Ui32(v495) {
		v467 = v495
		v469 = v469 + int32(-4)
		goto L60
	} else {
		goto L63
	}
L63:
	;
	goto L61
L64:
	;
	if v499 < v458 {
		v1514 = v2
		goto L2
	} else {
		goto L66
	}
L65:
	;
	v769 = v499
	v776 = v458
	v777 = v499 - v458
	goto L3
L66:
	;
	if v499 < v458 {
		goto L67
	} else {
		goto L68
	}
L67:
	;
	v528 = v458
	goto L69
L68:
	;
	v528 = v499
	goto L69
L69:
	;
	v534 = int32(0)
	v539 = v458
	goto L70
L70:
	;
	v557 = v539 << (uint(int32(2)) % 32)
	v559 = *(*int32)(unsafe.Add(mBase, uint32(v249+v557)))
	if v559 == int32(0) {
		v646 = v534
		goto L72
	} else {
		goto L73
	}
L72:
	;
	if base.B2i32(v539 == v528) == int32(0) {
		v534 = v646
		v539 = v539 + int32(1)
		goto L70
	} else {
		goto L83
	}
L73:
	;
	v563 = v534 + int32(1)
	if v563 < v404 {
		goto L74
	} else {
		goto L75
	}
L74:
	;
	v565 = v404
	goto L76
L75:
	;
	v565 = v563
	goto L76
L76:
	;
	v567 = v565 + int32(-1)
	v577 = v534
	v578 = v29 + int32(_a_F_VP8EncAnalyze_1) + v534<<(uint(int32(2))%32)
	goto L78
L77:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v29+int32(_a_F_VP8EncAnalyze_2)+v557))) = v621
	v630 = v621 << (uint(int32(2)) % 32)
	v631 = v29 + int32(_a_F_VP8EncAnalyze_3) + v630
	v632 = *(*int32)(unsafe.Add(mBase, uint32(v631)))
	*(*int32)(unsafe.Add(mBase, uint32(v631))) = v632 + v559*v539
	v638 = v29 + int32(_a_F_VP8EncAnalyze_4) + v630
	v639 = *(*int32)(unsafe.Add(mBase, uint32(v638)))
	*(*int32)(unsafe.Add(mBase, uint32(v638))) = v639 + v559
	v646 = v621
	goto L72
L78:
	;
	if v567 != v577 {
		goto L80
	} else {
		goto L81
	}
L79:
	;
	v621 = v577 + int32(0)
	goto L77
L80:
	;
	v600 = *(*int32)(unsafe.Add(mBase, uint32(v578)))
	v604 = v578 + int32(4)
	v605 = *(*int32)(unsafe.Add(mBase, uint32(v604)))
	v606 = v539 - v605
	v607 = int32(31)
	v608 = v606 >> (uint(v607) % 32)
	v611 = v539 - v600
	v613 = v611 >> (uint(v607) % 32)
	if base.Ui32(v606^v608-v608) < base.Ui32(v611^v613-v613) {
		v577 = v577 + int32(1)
		v578 = v604
		goto L78
	} else {
		goto L82
	}
L81:
	;
	v621 = v567
	goto L77
L82:
	;
	goto L79
L83:
	;
	v1514 = v2
	goto L2
L84:
	;
	v753 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3388)) = v753
	*(*int64)(unsafe.Add(mBase, uint32(l0)+1080)) = v753
	v757 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v758 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	v763 = F_WebPReportProgress(m, v757, v758+int32(20), l0+int32(368))
	mBase = m.M
	v2489 = v673
	goto L1
L85:
	;
	v685 = int32(1)
	v686 = int32(0)
	goto L86
L86:
	;
	v707 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v708 = v707 + v685
	v710 = v708 + int32(-1)
	v711 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v710))))
	v714 = int32(1)
	v715 = v711&int32(128) | v714
	*(*uint8)(unsafe.Add(mBase, uint32(v710))) = uint8(v715)
	v717 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v708))) = uint8(v717)
	v722 = v686 + v714
	v723 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v724 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	if v722 < v723*v724 {
		v685 = v685 + int32(4)
		v686 = v722
		goto L86
	} else {
		goto L88
	}
L87:
	;
	goto L84
L88:
	;
	goto L87
L89:
	;
	v970 = base.B2i32(v769 < v776)
	if v769 < v776 {
		goto L99
	} else {
		goto L100
	}
L90:
	;
	v913 = v404 - v886
	v914 = v889 * v777
	v915 = v29 + int32(_a_F_VP8EncAnalyze_1) + v886<<(uint(int32(2))%32)
	goto L96
L91:
	;
	v800 = v404 & int32(-4)
	v801 = int32(1)
	v806 = base.Simd_g_i32x4_splat(v794)
	v815 = v29 + int32(_a_F_VP8EncAnalyze_1)
	v816 = v800
	v829 = base.Simd_g_const(&F_VP8EncAnalyze__k0)
	goto L92
L92:
	;
	v837 = base.Simd_g_i32x4_mul(v829, base.Simd_g_i32x4_splat(v777))
	v838 = int32(0)
	v842 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l0(v837), base.Simd_g_i32x4_extract_lane_l0(v806))
	v844 = int32(1)
	v848 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l1(v837), base.Simd_g_i32x4_extract_lane_l1(v806))
	v851 = int32(2)
	v855 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l2(v837), base.Simd_g_i32x4_extract_lane_l2(v806))
	v858 = int32(3)
	v862 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l3(v837), base.Simd_g_i32x4_extract_lane_l3(v806))
	v865 = base.Simd_g_i32x4_add(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v842), v848), v855), v862), base.Simd_g_i32x4_splat(v776))
	base.Simd_g_v128_store(m, v815, v838, v865)
	v873 = v816 + int32(-4)
	if v873 != 0 {
		v815 = v815 + int32(16)
		v816 = v873
		v829 = base.Simd_g_i32x4_add(v829, base.Simd_g_const(&F_VP8EncAnalyze__k1))
		goto L92
	} else {
		goto L94
	}
L93:
	;
	if v404 == v800 {
		goto L89
	} else {
		goto L95
	}
L94:
	;
	goto L93
L95:
	;
	v886 = v800
	v889 = v800<<(uint(v801)%32) | v801
	goto L90
L96:
	;
	v936 = base.I32_div_s(v914, v794)
	*(*int32)(unsafe.Add(mBase, uint32(v915))) = v936 + v776
	v943 = v913 + int32(-1)
	if v943 != 0 {
		v913 = v943
		v914 = v914 + v777<<(uint(int32(1))%32)
		v915 = v915 + int32(4)
		goto L96
	} else {
		goto L98
	}
L97:
	;
	goto L89
L98:
	;
	goto L97
L99:
	;
	v971 = v776
	goto L101
L100:
	;
	v971 = v769
	goto L101
L101:
	;
	v973 = v404 << (uint(int32(2)) % 32)
	v998 = int32(0)
	goto L102
L102:
	;
	v1002 = v29 + int32(_a_F_VP8EncAnalyze_4)
	v1003 = int32(0)
	if base.Ui32(v973) < base.Ui32(int32(33)) {
		goto L105
	} else {
		goto L106
	}
L103:
	;
	v1490 = base.I32_div_s(v1469, int32(2))
	v1492 = base.I32_div_s(v1472+v1490, v1469)
	v1514 = v1492
	goto L2
L104:
	;
	v1125 = v29 + int32(_a_F_VP8EncAnalyze_3)
	v1126 = int32(0)
	if base.Ui32(v973) < base.Ui32(int32(33)) {
		goto L119
	} else {
		goto L120
	}
L105:
	;
	if v973 == int32(0) {
		goto L107
	} else {
		goto L108
	}
L106:
	;
	base.MemoryFill(m, v1002, v1003, v973)
	goto L104
L107:
	;
	goto L104
L108:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[12]))) = uint8(v1003)
	v1014 = v1002 + v973
	*(*uint8)(unsafe.Add(mBase, uint32(v1014+int32(-1)))) = uint8(v1003)
	if base.Ui32(v973) < base.Ui32(int32(3)) {
		goto L107
	} else {
		goto L109
	}
L109:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[13]))) = uint8(v1003)
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[14]))) = uint8(v1003)
	*(*uint8)(unsafe.Add(mBase, uint32(v1014+int32(-3)))) = uint8(v1003)
	*(*uint8)(unsafe.Add(mBase, uint32(v1014+int32(-2)))) = uint8(v1003)
	if base.Ui32(v973) < base.Ui32(int32(7)) {
		goto L107
	} else {
		goto L110
	}
L110:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[15]))) = uint8(v1003)
	*(*uint8)(unsafe.Add(mBase, uint32(v1014+int32(-4)))) = uint8(v1003)
	if base.Ui32(v973) < base.Ui32(int32(9)) {
		goto L107
	} else {
		goto L111
	}
L111:
	;
	v1036 = int32(0)
	v1039 = (v1036 - v1002) & int32(3)
	v1040 = v1002 + v1039
	*(*int32)(unsafe.Add(mBase, uint32(v1040))) = v1036
	v1048 = (v973 - v1039) & int32(60)
	v1049 = v1040 + v1048
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-4)))) = v1036
	if base.Ui32(v1048) < base.Ui32(int32(9)) {
		goto L107
	} else {
		goto L112
	}
L112:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1040)+8)) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1040)+4)) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-8)))) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-12)))) = v1036
	if base.Ui32(v1048) < base.Ui32(int32(25)) {
		goto L107
	} else {
		goto L113
	}
L113:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1040)+24)) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1040)+20)) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1040)+16)) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1040)+12)) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-16)))) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-20)))) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-24)))) = v1036
	*(*int32)(unsafe.Add(mBase, uint32(v1049+int32(-28)))) = v1036
	v1084 = v1040&int32(4) | int32(24)
	v1085 = v1048 - v1084
	if base.Ui32(v1085) < base.Ui32(int32(32)) {
		goto L107
	} else {
		goto L114
	}
L114:
	;
	v1090 = base.I64_extend_i32_u(v1036) * int64(4294967297)
	v1093 = v1085
	v1094 = v1040 + v1084
	goto L115
L115:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v1094)+24)) = v1090
	*(*int64)(unsafe.Add(mBase, uint32(v1094)+16)) = v1090
	*(*int64)(unsafe.Add(mBase, uint32(v1094)+8)) = v1090
	*(*int64)(unsafe.Add(mBase, uint32(v1094))) = v1090
	v1106 = v1093 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v1106) {
		v1093 = v1106
		v1094 = v1094 + int32(32)
		goto L115
	} else {
		goto L117
	}
L116:
	;
	goto L107
L117:
	;
	goto L116
L118:
	;
	if v769 < v776 {
		goto L132
	} else {
		goto L133
	}
L119:
	;
	if v973 == int32(0) {
		goto L121
	} else {
		goto L122
	}
L120:
	;
	base.MemoryFill(m, v1125, v1126, v973)
	goto L118
L121:
	;
	goto L118
L122:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[16]))) = uint8(v1126)
	v1137 = v1125 + v973
	*(*uint8)(unsafe.Add(mBase, uint32(v1137+int32(-1)))) = uint8(v1126)
	if base.Ui32(v973) < base.Ui32(int32(3)) {
		goto L121
	} else {
		goto L123
	}
L123:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[17]))) = uint8(v1126)
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[18]))) = uint8(v1126)
	*(*uint8)(unsafe.Add(mBase, uint32(v1137+int32(-3)))) = uint8(v1126)
	*(*uint8)(unsafe.Add(mBase, uint32(v1137+int32(-2)))) = uint8(v1126)
	if base.Ui32(v973) < base.Ui32(int32(7)) {
		goto L121
	} else {
		goto L124
	}
L124:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[19]))) = uint8(v1126)
	*(*uint8)(unsafe.Add(mBase, uint32(v1137+int32(-4)))) = uint8(v1126)
	if base.Ui32(v973) < base.Ui32(int32(9)) {
		goto L121
	} else {
		goto L125
	}
L125:
	;
	v1159 = int32(0)
	v1162 = (v1159 - v1125) & int32(3)
	v1163 = v1125 + v1162
	*(*int32)(unsafe.Add(mBase, uint32(v1163))) = v1159
	v1171 = (v973 - v1162) & int32(60)
	v1172 = v1163 + v1171
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-4)))) = v1159
	if base.Ui32(v1171) < base.Ui32(int32(9)) {
		goto L121
	} else {
		goto L126
	}
L126:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1163)+8)) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1163)+4)) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-8)))) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-12)))) = v1159
	if base.Ui32(v1171) < base.Ui32(int32(25)) {
		goto L121
	} else {
		goto L127
	}
L127:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1163)+24)) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1163)+20)) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1163)+16)) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1163)+12)) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-16)))) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-20)))) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-24)))) = v1159
	*(*int32)(unsafe.Add(mBase, uint32(v1172+int32(-28)))) = v1159
	v1207 = v1163&int32(4) | int32(24)
	v1208 = v1171 - v1207
	if base.Ui32(v1208) < base.Ui32(int32(32)) {
		goto L121
	} else {
		goto L128
	}
L128:
	;
	v1213 = base.I64_extend_i32_u(v1159) * int64(4294967297)
	v1216 = v1208
	v1217 = v1163 + v1207
	goto L129
L129:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v1217)+24)) = v1213
	*(*int64)(unsafe.Add(mBase, uint32(v1217)+16)) = v1213
	*(*int64)(unsafe.Add(mBase, uint32(v1217)+8)) = v1213
	*(*int64)(unsafe.Add(mBase, uint32(v1217))) = v1213
	v1229 = v1216 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v1229) {
		v1216 = v1229
		v1217 = v1217 + int32(32)
		goto L129
	} else {
		goto L131
	}
L130:
	;
	goto L121
L131:
	;
	goto L130
L132:
	;
	v1415 = int32(0)
	v1426 = v29 + int32(_a_F_VP8EncAnalyze_3)
	v1427 = v29 + int32(_a_F_VP8EncAnalyze_4)
	v1429 = v29 + int32(_a_F_VP8EncAnalyze_1)
	v1434 = v1415
	v1435 = v1415
	v1437 = v404
	v1438 = v1415
	goto L148
L133:
	;
	v1251 = v776
	v1252 = int32(0)
	goto L134
L134:
	;
	v1275 = v1251 << (uint(int32(2)) % 32)
	v1277 = *(*int32)(unsafe.Add(mBase, uint32(v249+v1275)))
	if v1277 == int32(0) {
		v1364 = v1252
		goto L136
	} else {
		goto L137
	}
L135:
	;
	goto L132
L136:
	;
	if v1251 != v971 {
		v1251 = v1251 + int32(1)
		v1252 = v1364
		goto L134
	} else {
		goto L147
	}
L137:
	;
	v1281 = v1252 + int32(1)
	if v1281 < v404 {
		goto L138
	} else {
		goto L139
	}
L138:
	;
	v1283 = v404
	goto L140
L139:
	;
	v1283 = v1281
	goto L140
L140:
	;
	v1285 = v1283 + int32(-1)
	v1295 = v1252
	v1296 = v29 + int32(_a_F_VP8EncAnalyze_1) + v1252<<(uint(int32(2))%32)
	goto L142
L141:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v29+int32(_a_F_VP8EncAnalyze_2)+v1275))) = v1338
	v1348 = v1338 << (uint(int32(2)) % 32)
	v1349 = v29 + int32(_a_F_VP8EncAnalyze_3) + v1348
	v1350 = *(*int32)(unsafe.Add(mBase, uint32(v1349)))
	*(*int32)(unsafe.Add(mBase, uint32(v1349))) = v1350 + v1277*v1251
	v1356 = v29 + int32(_a_F_VP8EncAnalyze_4) + v1348
	v1357 = *(*int32)(unsafe.Add(mBase, uint32(v1356)))
	*(*int32)(unsafe.Add(mBase, uint32(v1356))) = v1357 + v1277
	v1364 = v1338
	goto L136
L142:
	;
	if v1285 != v1295 {
		goto L144
	} else {
		goto L145
	}
L143:
	;
	v1338 = v1295 + int32(0)
	goto L141
L144:
	;
	v1318 = *(*int32)(unsafe.Add(mBase, uint32(v1296)))
	v1322 = v1296 + int32(4)
	v1323 = *(*int32)(unsafe.Add(mBase, uint32(v1322)))
	v1324 = v1251 - v1323
	v1325 = int32(31)
	v1326 = v1324 >> (uint(v1325) % 32)
	v1329 = v1251 - v1318
	v1331 = v1329 >> (uint(v1325) % 32)
	if base.Ui32(v1324^v1326-v1326) < base.Ui32(v1329^v1331-v1331) {
		v1295 = v1295 + int32(1)
		v1296 = v1322
		goto L142
	} else {
		goto L146
	}
L145:
	;
	v1338 = v1285
	goto L141
L146:
	;
	goto L143
L147:
	;
	goto L135
L148:
	;
	v1450 = *(*int32)(unsafe.Add(mBase, uint32(v1427)))
	if v1450 == int32(0) {
		v1469 = v1434
		v1470 = v1435
		v1472 = v1438
		goto L150
	} else {
		goto L151
	}
L149:
	;
	if v1470 < int32(5) {
		goto L153
	} else {
		goto L154
	}
L150:
	;
	v1474 = int32(4)
	v1481 = v1437 + int32(-1)
	if v1481 != 0 {
		v1426 = v1426 + v1474
		v1427 = v1427 + v1474
		v1429 = v1429 + v1474
		v1434 = v1469
		v1435 = v1470
		v1437 = v1481
		v1438 = v1472
		goto L148
	} else {
		goto L152
	}
L151:
	;
	v1453 = *(*int32)(unsafe.Add(mBase, uint32(v1429)))
	v1454 = *(*int32)(unsafe.Add(mBase, uint32(v1426)))
	v1456 = base.I32_div_s(v1450, int32(2))
	v1458 = base.I32_div_s(v1454+v1456, v1450)
	*(*int32)(unsafe.Add(mBase, uint32(v1429))) = v1458
	v1460 = v1453 - v1458
	v1462 = v1460 >> (uint(int32(31)) % 32)
	v1469 = v1450 + v1434
	v1470 = v1460 ^ v1462 - v1462 + v1435
	v1472 = v1458*v1450 + v1438
	goto L150
L152:
	;
	goto L149
L153:
	;
	goto L103
L154:
	;
	v1485 = v998 + int32(1)
	if v1485 != int32(6) {
		v998 = v1485
		goto L102
	} else {
		goto L155
	}
L155:
	;
	goto L153
L156:
	;
	if v400 < int32(2) {
		goto L161
	} else {
		goto L162
	}
L157:
	;
	v1524 = int32(0)
	v1529 = v1524
	v1530 = v1524
	goto L158
L158:
	;
	v1552 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1553 = v1552 + v1530
	v1556 = int32(1)
	v1557 = v1553 + v1556
	v1558 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1557))))
	v1559 = int32(2)
	v1562 = *(*int32)(unsafe.Add(mBase, uint32(v29+int32(_a_F_VP8EncAnalyze_2)+v1558<<(uint(v1559)%32))))
	v1567 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1553))))
	v1570 = v1562<<(uint(int32(5))%32)&int32(96) | v1567&int32(159)
	*(*uint8)(unsafe.Add(mBase, uint32(v1553))) = uint8(v1570)
	v1577 = *(*int32)(unsafe.Add(mBase, uint32(v29+int32(_a_F_VP8EncAnalyze_1)+v1562<<(uint(v1559)%32))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1557))) = uint8(v1577)
	v1582 = v1529 + v1556
	v1583 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v1584 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v1585 = v1583 * v1584
	if v1582 < v1585 {
		v1529 = v1582
		v1530 = v1530 + int32(4)
		goto L158
	} else {
		goto L160
	}
L159:
	;
	v1589 = v1583
	v1592 = v1585
	v1611 = v1584
	goto L156
L160:
	;
	goto L159
L161:
	;
	v2077 = *(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[20])))
	v2078 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if int32(2) <= v2078 {
		goto L202
	} else {
		goto L203
	}
L162:
	;
	v1615 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1616 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1615)+68)))
	if v1616&int32(1) == int32(0) {
		goto L161
	} else {
		goto L163
	}
L163:
	;
	v1621 = base.I64_extend_i32_s(v1592)
	v1622 = int32(1)
	if v1621 == int64(0) {
		goto L166
	} else {
		goto L167
	}
L164:
	;
	if v1643 == int32(0) {
		goto L161
	} else {
		goto L170
	}
L165:
	;
	goto L164
L166:
	;
	v1641 = F_malloc(m, base.I32_wrap_i64(v1621)*v1622)
	mBase = m.M
	v1643 = v1641
	goto L165
L167:
	;
	v1629 = base.I64_div_u_s(int64(2147418112), v1621)
	v1630 = int32(0)
	v1631 = base.I64_extend_i32_u(v1622)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1631*v1621) {
		v1643 = v1630
		goto L165
	} else {
		goto L168
	}
L168:
	;
	if base.Ui64(v1629) < base.Ui64(v1631) {
		v1643 = v1630
		goto L165
	} else {
		goto L169
	}
L169:
	;
	goto L166
L170:
	;
	if v1589 < int32(3) {
		goto L171
	} else {
		goto L172
	}
L171:
	;
	F_free(m, v1643)
	mBase = m.M
	goto L200
L172:
	;
	if v1611 < int32(3) {
		goto L171
	} else {
		goto L173
	}
L173:
	;
	v1652 = v1589 + int32(-1)
	v1654 = v1611 << (uint(int32(2)) % 32)
	v1659 = v1611 + v1643
	v1660 = int32(1)
	v1671 = int32(0)
	v1673 = v1660
	v1676 = v1659 + v1660
	goto L174
L174:
	;
	v1694 = int32(0)
	v1696 = v1671
	goto L176
L175:
	;
	v1867 = int32(1)
	v1872 = v1611 << (uint(int32(2)) % 32)
	v1892 = v1867
	v1893 = v1872 + int32(8)
	v1895 = v1659 + v1867
	goto L190
L176:
	;
	v1717 = base.Simd_g_const(&F_VP8EncAnalyze__k2)
	v1718 = int32(_a_F_VP8EncAnalyze_5)
	base.Simd_g_v128_store(m, v29, v1718, v1717)
	v1721 = v29 + v1718
	v1722 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1723 = v1722 + v1696
	v1724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1723))))
	v1725 = int32(3)
	v1727 = int32(12)
	v1729 = v1721 | int32(base.Ui32(v1724)>>(uint(v1725)%32))&v1727
	v1730 = *(*int32)(unsafe.Add(mBase, uint32(v1729)))
	v1731 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v1729))) = v1730 + v1731
	v1736 = int32(4)
	v1738 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1723+v1736))))
	v1743 = v1721 | int32(base.Ui32(v1738)>>(uint(v1725)%32))&v1727
	v1744 = *(*int32)(unsafe.Add(mBase, uint32(v1743)))
	*(*int32)(unsafe.Add(mBase, uint32(v1743))) = v1744 + v1731
	v1750 = int32(8)
	v1752 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1723+v1750))))
	v1757 = v1721 | int32(base.Ui32(v1752)>>(uint(v1725)%32))&v1727
	v1758 = *(*int32)(unsafe.Add(mBase, uint32(v1757)))
	*(*int32)(unsafe.Add(mBase, uint32(v1757))) = v1758 + v1731
	v1764 = v1723 + v1654
	v1765 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1764))))
	v1770 = v1721 | int32(base.Ui32(v1765)>>(uint(v1725)%32))&v1727
	v1771 = *(*int32)(unsafe.Add(mBase, uint32(v1770)))
	*(*int32)(unsafe.Add(mBase, uint32(v1770))) = v1771 + v1731
	v1779 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1764+v1750))))
	v1784 = v1721 | int32(base.Ui32(v1779)>>(uint(v1725)%32))&v1727
	v1785 = *(*int32)(unsafe.Add(mBase, uint32(v1784)))
	*(*int32)(unsafe.Add(mBase, uint32(v1784))) = v1785 + v1731
	v1791 = v1723 + v1611<<(uint(int32(3))%32)
	v1792 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1791))))
	v1797 = v1721 | int32(base.Ui32(v1792)>>(uint(v1725)%32))&v1727
	v1798 = *(*int32)(unsafe.Add(mBase, uint32(v1797)))
	*(*int32)(unsafe.Add(mBase, uint32(v1797))) = v1798 + v1731
	v1804 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1764+v1736))))
	v1809 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1791+v1736))))
	v1814 = v1721 | int32(base.Ui32(v1809)>>(uint(v1725)%32))&v1727
	v1815 = *(*int32)(unsafe.Add(mBase, uint32(v1814)))
	*(*int32)(unsafe.Add(mBase, uint32(v1814))) = v1815 + v1731
	v1823 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1791+v1750))))
	v1828 = v1721 | int32(base.Ui32(v1823)>>(uint(v1725)%32))&v1727
	v1829 = *(*int32)(unsafe.Add(mBase, uint32(v1828)))
	*(*int32)(unsafe.Add(mBase, uint32(v1828))) = v1829 + v1731
	v1833 = *(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[21])))
	if v1833 <= v1736 {
		goto L179
	} else {
		goto L180
	}
L177:
	;
	v1865 = v1673 + int32(1)
	if v1865 != v1652 {
		v1671 = v1671 + v1654
		v1673 = v1865
		v1676 = v1676 + v1611
		goto L174
	} else {
		goto L189
	}
L178:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1676+v1694))) = uint8(v1854)
	v1860 = v1694 + int32(1)
	if v1611+int32(-2) != v1860 {
		v1694 = v1860
		v1696 = v1696 + int32(4)
		goto L176
	} else {
		goto L188
	}
L179:
	;
	v1837 = *(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[22])))
	if v1837 <= int32(4) {
		goto L181
	} else {
		goto L182
	}
L180:
	;
	v1854 = int32(0)
	goto L178
L181:
	;
	v1841 = *(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[23])))
	if v1841 <= int32(4) {
		goto L183
	} else {
		goto L184
	}
L182:
	;
	v1854 = int32(1)
	goto L178
L183:
	;
	v1845 = int32(3)
	v1850 = *(*int32)(unsafe.Add(mBase, uint32(v29)+uint32(_c_F_VP8EncAnalyze[24])))
	if int32(4) < v1850 {
		goto L185
	} else {
		goto L186
	}
L184:
	;
	v1854 = int32(2)
	goto L178
L185:
	;
	v1853 = v1845
	goto L187
L186:
	;
	v1853 = int32(base.Ui32(v1804)>>(uint(int32(5))%32)) & v1845
	goto L187
L187:
	;
	v1854 = v1853
	goto L178
L188:
	;
	goto L177
L189:
	;
	goto L175
L190:
	;
	if v1611 == int32(3) {
		v1976 = int32(1)
		goto L192
	} else {
		goto L193
	}
L191:
	;
	goto L171
L192:
	;
	if v1611&v1867 == int32(0) {
		goto L197
	} else {
		goto L198
	}
L193:
	;
	v1914 = v1893
	v1915 = int32(0)
	goto L194
L194:
	;
	v1936 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1939 = v1936 + v1914 + int32(-4)
	v1940 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1939))))
	v1941 = int32(159)
	v1943 = v1895 + v1915
	v1944 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1943))))
	v1945 = int32(5)
	v1947 = int32(96)
	v1949 = v1940&v1941 | v1944<<(uint(v1945)%32)&v1947
	*(*uint8)(unsafe.Add(mBase, uint32(v1939))) = uint8(v1949)
	v1951 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1952 = v1951 + v1914
	v1953 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1952))))
	v1958 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1943+int32(1)))))
	v1963 = v1953&v1941 | v1958<<(uint(v1945)%32)&v1947
	*(*uint8)(unsafe.Add(mBase, uint32(v1952))) = uint8(v1963)
	v1968 = v1915 + int32(2)
	if v1611&int32(2147483646)+int32(-2) != v1968 {
		v1914 = v1914 + int32(8)
		v1915 = v1968
		goto L194
	} else {
		goto L196
	}
L195:
	;
	v1976 = v1915 + int32(3)
	goto L192
L196:
	;
	goto L195
L197:
	;
	v2022 = v1892 + int32(1)
	if v2022 != v1652 {
		v1892 = v2022
		v1893 = v1893 + v1872
		v1895 = v1895 + v1611
		goto L190
	} else {
		goto L199
	}
L198:
	;
	v2000 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v2002 = v1976 + v1892*v1611
	v2005 = v2000 + v2002<<(uint(int32(2))%32)
	v2006 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2005))))
	v2010 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1643+v2002))))
	v2015 = v2006&int32(159) | v2010<<(uint(int32(5))%32)&int32(96)
	*(*uint8)(unsafe.Add(mBase, uint32(v2005))) = uint8(v2015)
	goto L197
L199:
	;
	goto L191
L200:
	;
	goto L161
L201:
	;
	v2233 = int32(1)
	if v2078 < v2233 {
		v2489 = v2233
		goto L1
	} else {
		goto L220
	}
L202:
	;
	if base.Ui32(int32(4)) <= base.Ui32(v2078) {
		goto L205
	} else {
		goto L206
	}
L203:
	;
	v2210 = v2077
	v2212 = v2077
	goto L201
L204:
	;
	v2174 = v29 + int32(_a_F_VP8EncAnalyze_1) + v2142<<(uint(int32(2))%32)
	v2175 = v2143
	v2177 = v2145
	v2185 = v2078 - v2142
	goto L211
L205:
	;
	v2087 = v2078 & int32(2147483644)
	v2088 = base.Simd_g_i32x4_splat(v2077)
	v2092 = v2087
	v2093 = v29 + int32(_a_F_VP8EncAnalyze_1)
	v2105 = v2088
	v2108 = v2088
	goto L207
L206:
	;
	v2142 = int32(0)
	v2143 = v2077
	v2145 = v2077
	goto L204
L207:
	;
	v2116 = base.Simd_g_v128_load(m, v2093, int32(0))
	v2117 = base.Simd_g_i32x4_max_s(v2108, v2116)
	v2118 = base.Simd_g_i32x4_min_s(v2105, v2116)
	v2122 = v2092 + int32(-4)
	if v2122 != 0 {
		v2092 = v2122
		v2093 = v2093 + int32(16)
		v2105 = v2118
		v2108 = v2117
		goto L207
	} else {
		goto L209
	}
L208:
	;
	v2123 = base.Simd_g_const(&F_VP8EncAnalyze__k3)
	v2125 = base.Simd_g_i32x4_max_s(v2117, base.Simd_g_i8x16_swizzle_c(v2117, base.Simd_g_const(&F_VP8EncAnalyze__k3)))
	v2126 = base.Simd_g_const(&F_VP8EncAnalyze__k4)
	v2129 = int32(0)
	v2130 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_max_s(v2125, base.Simd_g_i8x16_swizzle_c(v2125, base.Simd_g_const(&F_VP8EncAnalyze__k4))))
	v2133 = base.Simd_g_i32x4_min_s(v2118, base.Simd_g_i8x16_swizzle_c(v2118, base.Simd_g_const(&F_VP8EncAnalyze__k3)))
	v2138 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_min_s(v2133, base.Simd_g_i8x16_swizzle_c(v2133, base.Simd_g_const(&F_VP8EncAnalyze__k4))))
	if v2078 == v2087 {
		v2210 = v2130
		v2212 = v2138
		goto L201
	} else {
		goto L210
	}
L209:
	;
	goto L208
L210:
	;
	v2142 = v2087
	v2143 = v2130
	v2145 = v2138
	goto L204
L211:
	;
	v2198 = *(*int32)(unsafe.Add(mBase, uint32(v2174)))
	if v2198 < v2175 {
		goto L213
	} else {
		goto L214
	}
L212:
	;
	v2210 = v2200
	v2212 = v2202
	goto L201
L213:
	;
	v2200 = v2175
	goto L215
L214:
	;
	v2200 = v2198
	goto L215
L215:
	;
	if v2177 < v2198 {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	v2202 = v2177
	goto L218
L217:
	;
	v2202 = v2198
	goto L218
L218:
	;
	v2206 = v2185 + int32(-1)
	if v2206 != 0 {
		v2174 = v2174 + int32(4)
		v2175 = v2200
		v2177 = v2202
		v2185 = v2206
		goto L211
	} else {
		goto L219
	}
L219:
	;
	goto L212
L220:
	;
	if v2210 == v2212 {
		goto L221
	} else {
		goto L222
	}
L221:
	;
	v2239 = v2212 + int32(1)
	goto L223
L222:
	;
	v2239 = v2210
	goto L223
L223:
	;
	v2240 = v2239 - v2212
	if base.Ui32(v2078) < base.Ui32(int32(4)) {
		v2390 = int32(0)
		goto L224
	} else {
		goto L225
	}
L224:
	;
	v2424 = v29 + int32(_a_F_VP8EncAnalyze_1) + v2390<<(uint(int32(2))%32)
	v2425 = v2390*int32(744) + l0 + int32(1084)
	v2434 = v2078 - v2390
	goto L230
L225:
	;
	v2245 = base.Simd_g_i32x4_splat(v2240)
	v2250 = v2078 & int32(2147483644)
	v2254 = v29 + int32(_a_F_VP8EncAnalyze_1)
	v2255 = l0
	v2264 = v2250
	goto L226
L226:
	;
	v2279 = int32(0)
	v2280 = base.Simd_g_v128_load(m, v2254, v2279)
	v2282 = base.Simd_g_const(&F_VP8EncAnalyze__k5)
	v2283 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_sub(v2280, base.Simd_g_i32x4_splat(v1514)), v2282)
	v2287 = base.Simd_g_i32x4_extract_lane_l0(v2245)
	v2288 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l0(v2283), v2287)
	v2290 = int32(1)
	v2293 = base.Simd_g_i32x4_extract_lane_l1(v2245)
	v2294 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l1(v2283), v2293)
	v2297 = int32(2)
	v2300 = base.Simd_g_i32x4_extract_lane_l2(v2245)
	v2301 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l2(v2283), v2300)
	v2304 = int32(3)
	v2307 = base.Simd_g_i32x4_extract_lane_l3(v2245)
	v2308 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l3(v2283), v2307)
	v2314 = base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v2288), v2294), v2301), v2308), base.Simd_g_const(&F_VP8EncAnalyze__k6)), base.Simd_g_const(&F_VP8EncAnalyze__k7))
	base.Simd_g_v128_store32_lane_l3(m, v2255+int32(3312), v2279, v2314)
	base.Simd_g_v128_store32_lane_l2(m, v2255+int32(2568), v2279, v2314)
	base.Simd_g_v128_store32_lane_l1(m, v2255+int32(1824), v2279, v2314)
	base.Simd_g_v128_store32_lane_l0(m, v2255+int32(1080), v2279, v2314)
	v2336 = base.Simd_g_i32x4_mul(base.Simd_g_i32x4_sub(v2280, base.Simd_g_i32x4_splat(v2212)), v2282)
	v2339 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l0(v2336), v2287)
	v2343 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l1(v2336), v2293)
	v2348 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l2(v2336), v2300)
	v2353 = base.I32_div_s(base.Simd_g_i32x4_extract_lane_l3(v2336), v2307)
	v2358 = base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v2339), v2343), v2348), v2353), v2282), base.Simd_g_const(&F_VP8EncAnalyze__k2))
	base.Simd_g_v128_store32_lane_l3(m, v2255+int32(3316), v2279, v2358)
	base.Simd_g_v128_store32_lane_l2(m, v2255+int32(2572), v2279, v2358)
	base.Simd_g_v128_store32_lane_l1(m, v2255+int32(1828), v2279, v2358)
	base.Simd_g_v128_store32_lane_l0(m, v2255+int32(1084), v2279, v2358)
	v2382 = v2264 + int32(-4)
	if v2382 != 0 {
		v2254 = v2254 + int32(16)
		v2255 = v2255 + int32(2976)
		v2264 = v2382
		goto L226
	} else {
		goto L228
	}
L227:
	;
	if v2078 == v2250 {
		v2489 = v2233
		goto L1
	} else {
		goto L229
	}
L228:
	;
	goto L227
L229:
	;
	v2390 = v2250
	goto L224
L230:
	;
	v2447 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	v2449 = int32(255)
	v2451 = base.I32_div_s((v2447-v2212)*v2449, v2240)
	if v2451 < v2449 {
		goto L232
	} else {
		goto L233
	}
L231:
	;
	v2489 = v2233
	goto L1
L232:
	;
	v2455 = v2451
	goto L234
L233:
	;
	v2455 = v2449
	goto L234
L234:
	;
	v2456 = int32(0)
	if v2456 < v2455 {
		goto L235
	} else {
		goto L236
	}
L235:
	;
	v2459 = v2455
	goto L237
L236:
	;
	v2459 = v2456
	goto L237
L237:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2425))) = v2459
	v2466 = base.I32_div_s((v2447-v1514)*int32(255), v2240)
	v2467 = int32(127)
	if v2466 < v2467 {
		goto L238
	} else {
		goto L239
	}
L238:
	;
	v2470 = v2466
	goto L240
L239:
	;
	v2470 = v2467
	goto L240
L240:
	;
	v2471 = int32(-127)
	if v2471 < v2470 {
		goto L241
	} else {
		goto L242
	}
L241:
	;
	v2474 = v2470
	goto L243
L242:
	;
	v2474 = v2471
	goto L243
L243:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2425+int32(-4)))) = v2474
	v2481 = v2434 + int32(-1)
	if v2481 != 0 {
		v2424 = v2424 + int32(4)
		v2425 = v2425 + int32(744)
		v2434 = v2481
		goto L230
	} else {
		goto L244
	}
L244:
	;
	goto L231
}

var F_VP8EncAnalyze__k0 = [2]uint64{0x300000001, 0x700000005}
var F_VP8EncAnalyze__k1 = [2]uint64{0x800000008, 0x800000008}
var F_VP8EncAnalyze__k2 = [2]uint64{0x0, 0x0}
var F_VP8EncAnalyze__k3 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_VP8EncAnalyze__k4 = [2]uint64{0x302010007060504, 0x302010003020100}
var F_VP8EncAnalyze__k5 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_VP8EncAnalyze__k6 = [2]uint64{0x7f0000007f, 0x7f0000007f}
var F_VP8EncAnalyze__k7 = [2]uint64{0xffffff81ffffff81, 0xffffff81ffffff81}

func F_VP8EncDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v28 int32
	_ = v28
	var v29 base.V128
	_ = v29
	var v30 base.V128
	_ = v30
	var v31 base.V128
	_ = v31
	var v32 base.V128
	_ = v32
	var v38 int32
	_ = v38
	var v44 base.V128
	_ = v44
	var v46 base.V128
	_ = v46
	var v56 base.V128
	_ = v56
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	var v67 base.V128
	_ = v67
	var v82 base.V128
	_ = v82
	var v85 base.V128
	_ = v85
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v170 int32
	_ = v170
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v224 int32
	_ = v224
	var v225 int32
	_ = v225
	var v229 int32
	_ = v229
	var v233 int32
	_ = v233
	var v237 int32
	_ = v237
	var v241 int32
	_ = v241
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v253 int32
	_ = v253
	var v257 int32
	_ = v257
	var v261 int32
	_ = v261
	var v265 int32
	_ = v265
	var v269 int32
	_ = v269
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v285 int32
	_ = v285
	var v289 int32
	_ = v289
	var v293 int32
	_ = v293
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v300 int32
	_ = v300
	var v304 int32
	_ = v304
	var v305 int32
	_ = v305
	var v309 int32
	_ = v309
	var v313 int32
	_ = v313
	var v317 int32
	_ = v317
	var v321 int32
	_ = v321
	var v325 int32
	_ = v325
	var v329 int32
	_ = v329
	var v332 int32
	_ = v332
	var v333 int32
	_ = v333
	v11 = m.G1
	v14 = *(*int32)(unsafe.Add(mBase, uint32(v11)+uint32(_c_F_VP8EncDspInit[0])))
	v15 = m.G13
	v16 = *(*int32)(unsafe.Add(mBase, uint32(v15)))
	if v14 == v16 {
	} else {
		v18 = m.G1
		F_VP8DspInit(m)
		mBase = m.M
		v22 = *(*int32)(unsafe.Add(mBase, uint32(v18)+uint32(_c_F_VP8EncDspInit[1])))
		if v22 != 0 {
		} else {
			v28 = int32(-736)
			v29 = base.Simd_g_const(&F_VP8EncDspInit__k0)
			v30 = base.Simd_g_const(&F_VP8EncDspInit__k1)
			v31 = base.Simd_g_const(&F_VP8EncDspInit__k2)
			v32 = base.Simd_g_const(&F_VP8EncDspInit__k3)
			for {
				v38 = m.G1
				v44 = base.Simd_g_const(&F_VP8EncDspInit__k4)
				v46 = base.Simd_g_const(&F_VP8EncDspInit__k5)
				v56 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(v32, v44), v46), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(v31, v44), v46)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(v30, v44), v46), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(v29, v44), v46)))
				v57 = int32(0)
				base.Simd_g_v128_store(m, v38+int32(_a_F_VP8EncDspInit_0)+v28+int32(736), v57, v56)
				if v28 == v57 {
					break
				} else {
					v61 = m.G1
					v67 = base.Simd_g_const(&F_VP8EncDspInit__k6)
					v82 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_add(v32, v67), v44), v46), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_add(v31, v67), v44), v46)), base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_add(v30, v67), v44), v46), base.Simd_g_i32x4_min_s(base.Simd_g_i32x4_max_s(base.Simd_g_i32x4_add(v29, v67), v44), v46)))
					base.Simd_g_v128_store(m, v61+int32(_a_F_VP8EncDspInit_0)+v28+int32(752), int32(0), v82)
					v85 = base.Simd_g_const(&F_VP8EncDspInit__k7)
					v28 = v28 + int32(32)
					v29 = base.Simd_g_i32x4_add(v29, v85)
					v30 = base.Simd_g_i32x4_add(v30, v85)
					v31 = base.Simd_g_i32x4_add(v31, v85)
					v32 = base.Simd_g_i32x4_add(v32, v85)
					continue
				}
				break
			}
			v92 = m.G1
			v95 = int32(_a_F_VP8EncDspInit_1)
			*(*uint16)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_VP8EncDspInit[2]))) = uint16(v95)
			*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_VP8EncDspInit[3]))) = int32(-1)
			*(*int64)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_VP8EncDspInit[4]))) = int64(-1)
			*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_VP8EncDspInit[1]))) = int32(1)
		}
		v115 = m.G2
		v116 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[5]))) = v115 + int32(232)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[6]))) = v115 + int32(233)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[7]))) = v115 + int32(234)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[8]))) = v115 + int32(235)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[9]))) = v115 + int32(236)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[10]))) = v115 + int32(237)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[11]))) = v115 + int32(238)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[12]))) = v115 + int32(239)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[13]))) = v115 + int32(240)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[14]))) = v115 + int32(241)
		v170 = v115 + int32(242)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[15]))) = v170
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[16]))) = v115 + int32(243)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[17]))) = v170
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[18]))) = v115 + int32(244)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[19]))) = v115 + int32(245)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[20]))) = v115 + int32(246)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[21]))) = v115 + int32(247)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[22]))) = v115 + int32(248)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[23]))) = v115 + int32(249)
		*(*int32)(unsafe.Add(mBase, uint32(v116)+uint32(_c_F_VP8EncDspInit[24]))) = v115 + int32(250)
		v215 = m.G13
		v216 = *(*int32)(unsafe.Add(mBase, uint32(v215)))
		if v216 == int32(0) {
		} else {
			v219 = int32(0)
			v220 = m.T0[v216].(func(*base.Module, int32) int32)(m, v219)
			mBase = m.M
			if v220 == v219 {
			} else {
				v224 = m.G2
				v225 = m.G58
				*(*int32)(unsafe.Add(mBase, uint32(v225))) = v224 + int32(208)
				v229 = m.G59
				*(*int32)(unsafe.Add(mBase, uint32(v229))) = v224 + int32(209)
				v233 = m.G60
				*(*int32)(unsafe.Add(mBase, uint32(v233))) = v224 + int32(210)
				v237 = m.G61
				*(*int32)(unsafe.Add(mBase, uint32(v237))) = v224 + int32(211)
				v241 = m.G62
				*(*int32)(unsafe.Add(mBase, uint32(v241))) = v224 + int32(212)
				v245 = m.G63
				*(*int32)(unsafe.Add(mBase, uint32(v245))) = v224 + int32(213)
				v249 = m.G64
				*(*int32)(unsafe.Add(mBase, uint32(v249))) = v224 + int32(214)
				v253 = m.G65
				*(*int32)(unsafe.Add(mBase, uint32(v253))) = v224 + int32(215)
				v257 = m.G66
				*(*int32)(unsafe.Add(mBase, uint32(v257))) = v224 + int32(216)
				v261 = m.G67
				*(*int32)(unsafe.Add(mBase, uint32(v261))) = v224 + int32(217)
				v265 = m.G68
				*(*int32)(unsafe.Add(mBase, uint32(v265))) = v224 + int32(218)
				v269 = m.G69
				*(*int32)(unsafe.Add(mBase, uint32(v269))) = v224 + int32(219)
				v273 = m.G70
				*(*int32)(unsafe.Add(mBase, uint32(v273))) = v224 + int32(220)
				v277 = m.G71
				*(*int32)(unsafe.Add(mBase, uint32(v277))) = v224 + int32(221)
				v281 = m.G72
				*(*int32)(unsafe.Add(mBase, uint32(v281))) = v224 + int32(222)
				v285 = m.G73
				*(*int32)(unsafe.Add(mBase, uint32(v285))) = v224 + int32(223)
				v289 = m.G74
				*(*int32)(unsafe.Add(mBase, uint32(v289))) = v224 + int32(224)
				v293 = m.G75
				*(*int32)(unsafe.Add(mBase, uint32(v293))) = v224 + int32(225)
				v298 = m.G13
				v299 = *(*int32)(unsafe.Add(mBase, uint32(v298)))
				v300 = m.T0[v299].(func(*base.Module, int32) int32)(m, int32(3))
				mBase = m.M
				if v300 == int32(0) {
				} else {
					v304 = m.G2
					v305 = m.G62
					*(*int32)(unsafe.Add(mBase, uint32(v305))) = v304 + int32(226)
					v309 = m.G59
					*(*int32)(unsafe.Add(mBase, uint32(v309))) = v304 + int32(227)
					v313 = m.G63
					*(*int32)(unsafe.Add(mBase, uint32(v313))) = v304 + int32(228)
					v317 = m.G64
					*(*int32)(unsafe.Add(mBase, uint32(v317))) = v304 + int32(229)
					v321 = m.G73
					*(*int32)(unsafe.Add(mBase, uint32(v321))) = v304 + int32(230)
					v325 = m.G74
					*(*int32)(unsafe.Add(mBase, uint32(v325))) = v304 + int32(231)
				}
			}
		}
		v329 = m.G1
		v332 = m.G13
		v333 = *(*int32)(unsafe.Add(mBase, uint32(v332)))
		*(*int32)(unsafe.Add(mBase, uint32(v329)+uint32(_c_F_VP8EncDspInit[0]))) = v333
	}
	return
}

var F_VP8EncDspInit__k0 = [2]uint64{0xffffff0effffff0d, 0xffffff10ffffff0f}
var F_VP8EncDspInit__k1 = [2]uint64{0xffffff0affffff09, 0xffffff0cffffff0b}
var F_VP8EncDspInit__k2 = [2]uint64{0xffffff06ffffff05, 0xffffff08ffffff07}
var F_VP8EncDspInit__k3 = [2]uint64{0xffffff02ffffff01, 0xffffff04ffffff03}
var F_VP8EncDspInit__k4 = [2]uint64{0x0, 0x0}
var F_VP8EncDspInit__k5 = [2]uint64{0xff000000ff, 0xff000000ff}
var F_VP8EncDspInit__k6 = [2]uint64{0x1000000010, 0x1000000010}
var F_VP8EncDspInit__k7 = [2]uint64{0x2000000020, 0x2000000020}

func F_VP8EncFreeBitWriters(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	var v14 base.V128
	_ = v14
	var v15 int32
	_ = v15
	var v21 int32
	_ = v21
	var v28 int32
	_ = v28
	var v29 int32
	_ = v29
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v38 base.V128
	_ = v38
	var v39 int32
	_ = v39
	var v48 int32
	_ = v48
	var v49 int32
	_ = v49
	v5 = l0 + int32(56)
	if v5 == int32(0) {
	} else {
		v11 = l0 + int32(72)
		v12 = *(*int32)(unsafe.Add(mBase, uint32(v11)))
		F_WebPSafeFree(m, v12)
		mBase = m.M
		v14 = base.Simd_g_const(&F_VP8EncFreeBitWriters__k0)
		v15 = int32(0)
		base.Simd_g_v128_store(m, v11, v15, v14)
		base.Simd_g_v128_store(m, v5, v15, v14)
	}
	v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v21 < int32(1) {
	} else {
		v28 = l0 + int32(88)
		v29 = int32(0)
		for {
			if v28 == int32(0) {
			} else {
				v35 = v28 + int32(16)
				v36 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
				F_WebPSafeFree(m, v36)
				mBase = m.M
				v38 = base.Simd_g_const(&F_VP8EncFreeBitWriters__k0)
				v39 = int32(0)
				base.Simd_g_v128_store(m, v35, v39, v38)
				base.Simd_g_v128_store(m, v28, v39, v38)
			}
			v48 = v29 + int32(1)
			v49 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
			if v48 < v49 {
				v28 = v28 + int32(32)
				v29 = v48
				continue
			} else {
				break
			}
			break
		}
	}
	return
}

var F_VP8EncFreeBitWriters__k0 = [2]uint64{0x0, 0x0}

func F_VP8EncTokenLoop(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v70 float32
	_ = v70
	var v71 int32
	_ = v71
	var v72 float32
	_ = v72
	var v73 float32
	_ = v73
	var v75 float32
	_ = v75
	var v77 float32
	_ = v77
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v84 int32
	_ = v84
	var v88 int32
	_ = v88
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v95 float32
	_ = v95
	var v100 float64
	_ = v100
	var v101 float64
	_ = v101
	var v108 int64
	_ = v108
	var v109 int32
	_ = v109
	var v112 int32
	_ = v112
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v175 int32
	_ = v175
	var v176 int32
	_ = v176
	var v178 int32
	_ = v178
	var v180 base.V128
	_ = v180
	var v183 base.V128
	_ = v183
	var v189 int32
	_ = v189
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v197 int32
	_ = v197
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v215 int32
	_ = v215
	var v223 int32
	_ = v223
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v236 int32
	_ = v236
	var v237 int32
	_ = v237
	var v242 int32
	_ = v242
	var v244 int32
	_ = v244
	var v254 int32
	_ = v254
	var v255 int32
	_ = v255
	var v258 int32
	_ = v258
	var v262 int32
	_ = v262
	var v316 int32
	_ = v316
	var v326 int32
	_ = v326
	var v329 float32
	_ = v329
	var v335 float32
	_ = v335
	var v341 int32
	_ = v341
	var v371 float64
	_ = v371
	var v372 int32
	_ = v372
	var v373 float32
	_ = v373
	var v384 int32
	_ = v384
	var v385 int32
	_ = v385
	var v392 int32
	_ = v392
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v398 int32
	_ = v398
	var v403 int32
	_ = v403
	var v405 int32
	_ = v405
	var v407 int32
	_ = v407
	var v412 int32
	_ = v412
	var v416 int32
	_ = v416
	var v417 int32
	_ = v417
	var v422 int32
	_ = v422
	var v439 int32
	_ = v439
	var v441 int32
	_ = v441
	var v443 int64
	_ = v443
	var v445 int32
	_ = v445
	var v447 int32
	_ = v447
	var v448 int32
	_ = v448
	var v452 int32
	_ = v452
	var v457 int32
	_ = v457
	var v458 int64
	_ = v458
	var v464 int32
	_ = v464
	var v467 int32
	_ = v467
	var v472 int32
	_ = v472
	var v477 int32
	_ = v477
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v482 int32
	_ = v482
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
	var v491 int32
	_ = v491
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v499 int32
	_ = v499
	var v502 int32
	_ = v502
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v513 int32
	_ = v513
	var v640 int32
	_ = v640
	var v645 int32
	_ = v645
	var v647 int32
	_ = v647
	var v652 int32
	_ = v652
	var v656 int32
	_ = v656
	var v658 int32
	_ = v658
	var v663 int32
	_ = v663
	var v669 int32
	_ = v669
	var v670 int32
	_ = v670
	var v673 int32
	_ = v673
	var v679 int64
	_ = v679
	var v691 int32
	_ = v691
	var v736 int64
	_ = v736
	var v737 int64
	_ = v737
	var v744 int32
	_ = v744
	var v750 int32
	_ = v750
	var v787 int32
	_ = v787
	var v792 int32
	_ = v792
	var v793 int32
	_ = v793
	var v794 int32
	_ = v794
	var v795 int32
	_ = v795
	var v796 int32
	_ = v796
	var v798 int32
	_ = v798
	var v799 int32
	_ = v799
	var v801 int32
	_ = v801
	var v828 int32
	_ = v828
	var v831 int32
	_ = v831
	var v832 int32
	_ = v832
	var v833 int32
	_ = v833
	var v834 int32
	_ = v834
	var v835 int32
	_ = v835
	var v836 int32
	_ = v836
	var v837 int32
	_ = v837
	var v857 int32
	_ = v857
	var v867 int32
	_ = v867
	var v868 int32
	_ = v868
	var v878 int32
	_ = v878
	var v880 int32
	_ = v880
	var v881 int32
	_ = v881
	var v883 int32
	_ = v883
	var v884 int32
	_ = v884
	var v886 int32
	_ = v886
	var v888 int32
	_ = v888
	var v890 int32
	_ = v890
	var v893 int32
	_ = v893
	var v895 int32
	_ = v895
	var v896 int32
	_ = v896
	var v899 int32
	_ = v899
	var v900 int32
	_ = v900
	var v903 int32
	_ = v903
	var v905 int32
	_ = v905
	var v910 int32
	_ = v910
	var v916 int32
	_ = v916
	var v919 int32
	_ = v919
	var v925 int32
	_ = v925
	var v930 int32
	_ = v930
	var v938 int32
	_ = v938
	var v965 int32
	_ = v965
	var v970 int32
	_ = v970
	var v982 int32
	_ = v982
	var v992 int32
	_ = v992
	var v998 int32
	_ = v998
	var v1002 int32
	_ = v1002
	var v1004 int32
	_ = v1004
	var v1005 int32
	_ = v1005
	var v1006 int32
	_ = v1006
	var v1008 int32
	_ = v1008
	var v1010 int32
	_ = v1010
	var v1011 int32
	_ = v1011
	var v1015 int32
	_ = v1015
	var v1017 int32
	_ = v1017
	var v1019 int32
	_ = v1019
	var v1022 int32
	_ = v1022
	var v1024 int32
	_ = v1024
	var v1025 int32
	_ = v1025
	var v1028 int32
	_ = v1028
	var v1029 int32
	_ = v1029
	var v1032 int32
	_ = v1032
	var v1034 int32
	_ = v1034
	var v1039 int32
	_ = v1039
	var v1045 int32
	_ = v1045
	var v1048 int32
	_ = v1048
	var v1054 int32
	_ = v1054
	var v1059 int32
	_ = v1059
	var v1067 int32
	_ = v1067
	var v1089 int32
	_ = v1089
	var v1094 int32
	_ = v1094
	var v1106 int32
	_ = v1106
	var v1116 int32
	_ = v1116
	var v1122 int32
	_ = v1122
	var v1126 int32
	_ = v1126
	var v1128 int32
	_ = v1128
	var v1129 int32
	_ = v1129
	var v1130 int32
	_ = v1130
	var v1132 int32
	_ = v1132
	var v1134 int32
	_ = v1134
	var v1135 int32
	_ = v1135
	var v1139 int32
	_ = v1139
	var v1141 int32
	_ = v1141
	var v1143 int32
	_ = v1143
	var v1146 int32
	_ = v1146
	var v1148 int32
	_ = v1148
	var v1149 int32
	_ = v1149
	var v1152 int32
	_ = v1152
	var v1153 int32
	_ = v1153
	var v1156 int32
	_ = v1156
	var v1158 int32
	_ = v1158
	var v1163 int32
	_ = v1163
	var v1169 int32
	_ = v1169
	var v1172 int32
	_ = v1172
	var v1178 int32
	_ = v1178
	var v1183 int32
	_ = v1183
	var v1191 int32
	_ = v1191
	var v1213 int32
	_ = v1213
	var v1218 int32
	_ = v1218
	var v1221 int32
	_ = v1221
	var v1223 int32
	_ = v1223
	var v1234 int32
	_ = v1234
	var v1237 int32
	_ = v1237
	var v1239 int32
	_ = v1239
	var v1250 int32
	_ = v1250
	var v1294 int32
	_ = v1294
	var v1309 int32
	_ = v1309
	var v1310 int32
	_ = v1310
	var v1311 int32
	_ = v1311
	var v1313 int32
	_ = v1313
	var v1316 int32
	_ = v1316
	var v1319 int32
	_ = v1319
	var v1322 int32
	_ = v1322
	var v1325 int32
	_ = v1325
	var v1328 int32
	_ = v1328
	var v1331 int32
	_ = v1331
	var v1334 int32
	_ = v1334
	var v1337 int32
	_ = v1337
	var v1340 int32
	_ = v1340
	var v1343 int32
	_ = v1343
	var v1346 int32
	_ = v1346
	var v1349 int32
	_ = v1349
	var v1352 int32
	_ = v1352
	var v1355 int32
	_ = v1355
	var v1360 int32
	_ = v1360
	var v1361 int32
	_ = v1361
	var v1362 int32
	_ = v1362
	var v1383 int32
	_ = v1383
	var v1402 int32
	_ = v1402
	var v1434 int32
	_ = v1434
	var v1435 int32
	_ = v1435
	var v1436 int32
	_ = v1436
	var v1437 int32
	_ = v1437
	var v1450 int32
	_ = v1450
	var v1451 int32
	_ = v1451
	var v1452 int32
	_ = v1452
	var v1455 int32
	_ = v1455
	var v1456 int32
	_ = v1456
	var v1460 int32
	_ = v1460
	var v1467 int32
	_ = v1467
	var v1500 int32
	_ = v1500
	var v1511 int32
	_ = v1511
	var v1516 int32
	_ = v1516
	var v1517 int32
	_ = v1517
	var v1519 int32
	_ = v1519
	var v1528 int32
	_ = v1528
	var v1555 int32
	_ = v1555
	var v1557 int32
	_ = v1557
	var v1558 int32
	_ = v1558
	var v1566 int32
	_ = v1566
	var v1567 int32
	_ = v1567
	var v1569 int32
	_ = v1569
	var v1578 int32
	_ = v1578
	var v1580 int32
	_ = v1580
	var v1582 int32
	_ = v1582
	var v1626 int32
	_ = v1626
	var v1631 int32
	_ = v1631
	var v1634 int32
	_ = v1634
	var v1637 int32
	_ = v1637
	var v1640 int32
	_ = v1640
	var v1641 int32
	_ = v1641
	var v1642 int32
	_ = v1642
	var v1646 int32
	_ = v1646
	var v1647 int32
	_ = v1647
	var v1648 int32
	_ = v1648
	var v1653 int32
	_ = v1653
	var v1654 int32
	_ = v1654
	var v1661 int32
	_ = v1661
	var v1697 int32
	_ = v1697
	var v1708 int32
	_ = v1708
	var v1713 int32
	_ = v1713
	var v1714 int32
	_ = v1714
	var v1716 int32
	_ = v1716
	var v1725 int32
	_ = v1725
	var v1752 int32
	_ = v1752
	var v1754 int32
	_ = v1754
	var v1755 int32
	_ = v1755
	var v1763 int32
	_ = v1763
	var v1764 int32
	_ = v1764
	var v1766 int32
	_ = v1766
	var v1775 int32
	_ = v1775
	var v1777 int32
	_ = v1777
	var v1779 int32
	_ = v1779
	var v1823 int32
	_ = v1823
	var v1828 int32
	_ = v1828
	var v1831 int32
	_ = v1831
	var v1834 int32
	_ = v1834
	var v1837 int32
	_ = v1837
	var v1838 int32
	_ = v1838
	var v1839 int32
	_ = v1839
	var v1843 int32
	_ = v1843
	var v1844 int32
	_ = v1844
	var v1845 int32
	_ = v1845
	var v1850 int32
	_ = v1850
	var v1851 int32
	_ = v1851
	var v1858 int32
	_ = v1858
	var v1894 int32
	_ = v1894
	var v1905 int32
	_ = v1905
	var v1910 int32
	_ = v1910
	var v1911 int32
	_ = v1911
	var v1913 int32
	_ = v1913
	var v1922 int32
	_ = v1922
	var v1949 int32
	_ = v1949
	var v1951 int32
	_ = v1951
	var v1952 int32
	_ = v1952
	var v1960 int32
	_ = v1960
	var v1961 int32
	_ = v1961
	var v1963 int32
	_ = v1963
	var v1972 int32
	_ = v1972
	var v1974 int32
	_ = v1974
	var v1976 int32
	_ = v1976
	var v2020 int32
	_ = v2020
	var v2025 int32
	_ = v2025
	var v2028 int32
	_ = v2028
	var v2031 int32
	_ = v2031
	var v2034 int32
	_ = v2034
	var v2041 int32
	_ = v2041
	var v2046 int32
	_ = v2046
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
	var v2061 int32
	_ = v2061
	var v2063 int32
	_ = v2063
	var v2065 int32
	_ = v2065
	var v2067 int32
	_ = v2067
	var v2069 int32
	_ = v2069
	var v2071 int32
	_ = v2071
	var v2073 int32
	_ = v2073
	var v2075 int32
	_ = v2075
	var v2077 int32
	_ = v2077
	var v2079 int32
	_ = v2079
	var v2082 int32
	_ = v2082
	var v2175 int32
	_ = v2175
	var v2182 int32
	_ = v2182
	var v2228 int32
	_ = v2228
	var v2230 int32
	_ = v2230
	var v2231 int32
	_ = v2231
	var v2232 int32
	_ = v2232
	var v2237 int32
	_ = v2237
	var v2240 int32
	_ = v2240
	var v2241 int32
	_ = v2241
	var v2244 int32
	_ = v2244
	var v2247 int32
	_ = v2247
	var v2257 int32
	_ = v2257
	var v2267 int32
	_ = v2267
	var v2327 int32
	_ = v2327
	var v2328 int32
	_ = v2328
	var v2333 int32
	_ = v2333
	var v2334 int32
	_ = v2334
	var v2360 int32
	_ = v2360
	var v2361 int32
	_ = v2361
	var v2362 int32
	_ = v2362
	var v2367 int32
	_ = v2367
	var v2420 int32
	_ = v2420
	var v2421 int32
	_ = v2421
	var v2422 int32
	_ = v2422
	var v2424 int32
	_ = v2424
	var v2425 int32
	_ = v2425
	var v2430 int32
	_ = v2430
	var v2433 int32
	_ = v2433
	var v2436 int32
	_ = v2436
	var v2441 int32
	_ = v2441
	var v2444 int32
	_ = v2444
	var v2447 int32
	_ = v2447
	var v2452 int32
	_ = v2452
	var v2455 int32
	_ = v2455
	var v2458 int32
	_ = v2458
	var v2463 int32
	_ = v2463
	var v2466 int32
	_ = v2466
	var v2467 int32
	_ = v2467
	var v2470 int32
	_ = v2470
	var v2475 int32
	_ = v2475
	var v2478 int32
	_ = v2478
	var v2481 int32
	_ = v2481
	var v2486 int32
	_ = v2486
	var v2489 int32
	_ = v2489
	var v2492 int32
	_ = v2492
	var v2497 int32
	_ = v2497
	var v2500 int32
	_ = v2500
	var v2503 int32
	_ = v2503
	var v2508 int32
	_ = v2508
	var v2511 int32
	_ = v2511
	var v2512 int32
	_ = v2512
	var v2515 int32
	_ = v2515
	var v2520 int32
	_ = v2520
	var v2523 int32
	_ = v2523
	var v2526 int32
	_ = v2526
	var v2531 int32
	_ = v2531
	var v2534 int32
	_ = v2534
	var v2537 int32
	_ = v2537
	var v2542 int32
	_ = v2542
	var v2545 int32
	_ = v2545
	var v2548 int32
	_ = v2548
	var v2553 int32
	_ = v2553
	var v2556 int32
	_ = v2556
	var v2557 int32
	_ = v2557
	var v2560 int32
	_ = v2560
	var v2565 int32
	_ = v2565
	var v2568 int32
	_ = v2568
	var v2571 int32
	_ = v2571
	var v2576 int32
	_ = v2576
	var v2579 int32
	_ = v2579
	var v2582 int32
	_ = v2582
	var v2587 int32
	_ = v2587
	var v2590 int32
	_ = v2590
	var v2593 int32
	_ = v2593
	var v2598 int32
	_ = v2598
	var v2625 int32
	_ = v2625
	var v2626 int32
	_ = v2626
	var v2628 int32
	_ = v2628
	var v2629 int32
	_ = v2629
	var v2634 int32
	_ = v2634
	var v2637 int32
	_ = v2637
	var v2640 int32
	_ = v2640
	var v2645 int32
	_ = v2645
	var v2648 int32
	_ = v2648
	var v2649 int32
	_ = v2649
	var v2652 int32
	_ = v2652
	var v2657 int32
	_ = v2657
	var v2660 int32
	_ = v2660
	var v2663 int32
	_ = v2663
	var v2668 int32
	_ = v2668
	var v2671 int32
	_ = v2671
	var v2672 int32
	_ = v2672
	var v2675 int32
	_ = v2675
	var v2680 int32
	_ = v2680
	var v2683 int32
	_ = v2683
	var v2686 int32
	_ = v2686
	var v2691 int32
	_ = v2691
	var v2694 int32
	_ = v2694
	var v2695 int32
	_ = v2695
	var v2698 int32
	_ = v2698
	var v2703 int32
	_ = v2703
	var v2706 int32
	_ = v2706
	var v2709 int32
	_ = v2709
	var v2714 int32
	_ = v2714
	var v2718 int32
	_ = v2718
	var v2719 int32
	_ = v2719
	var v2720 int32
	_ = v2720
	var v2723 int32
	_ = v2723
	var v2727 int32
	_ = v2727
	var v2731 int32
	_ = v2731
	var v2735 int32
	_ = v2735
	var v2739 int32
	_ = v2739
	var v2743 int32
	_ = v2743
	var v2747 int32
	_ = v2747
	var v2751 int32
	_ = v2751
	var v2755 int32
	_ = v2755
	var v2759 int32
	_ = v2759
	var v2763 int32
	_ = v2763
	var v2767 int32
	_ = v2767
	var v2771 int32
	_ = v2771
	var v2776 int32
	_ = v2776
	var v2779 int32
	_ = v2779
	var v2781 int32
	_ = v2781
	var v2785 int64
	_ = v2785
	var v2786 int64
	_ = v2786
	var v2788 int32
	_ = v2788
	var v2792 int32
	_ = v2792
	var v2793 int32
	_ = v2793
	var v2794 int32
	_ = v2794
	var v2795 int32
	_ = v2795
	var v2796 int32
	_ = v2796
	var v2800 int32
	_ = v2800
	var v2801 int32
	_ = v2801
	var v2803 int32
	_ = v2803
	var v2804 int32
	_ = v2804
	var v2806 int32
	_ = v2806
	var v2807 int32
	_ = v2807
	var v2809 int32
	_ = v2809
	var v2810 int32
	_ = v2810
	var v2812 int32
	_ = v2812
	var v2813 int32
	_ = v2813
	var v2815 int32
	_ = v2815
	var v2816 int32
	_ = v2816
	var v2818 int32
	_ = v2818
	var v2819 int32
	_ = v2819
	var v2821 int32
	_ = v2821
	var v2822 int32
	_ = v2822
	var v2824 int32
	_ = v2824
	var v2825 int32
	_ = v2825
	var v2827 int32
	_ = v2827
	var v2828 int32
	_ = v2828
	var v2830 int32
	_ = v2830
	var v2831 int32
	_ = v2831
	var v2833 int32
	_ = v2833
	var v2834 int32
	_ = v2834
	var v2836 int32
	_ = v2836
	var v2837 int32
	_ = v2837
	var v2839 int32
	_ = v2839
	var v2840 int32
	_ = v2840
	var v2842 int32
	_ = v2842
	var v2843 int32
	_ = v2843
	var v2845 int32
	_ = v2845
	var v2846 int32
	_ = v2846
	var v2848 int32
	_ = v2848
	var v2849 int32
	_ = v2849
	var v2851 int32
	_ = v2851
	var v2852 int32
	_ = v2852
	var v2854 int32
	_ = v2854
	var v2855 int32
	_ = v2855
	var v2857 int32
	_ = v2857
	var v2858 int32
	_ = v2858
	var v2860 int32
	_ = v2860
	var v2861 int32
	_ = v2861
	var v2863 int32
	_ = v2863
	var v2864 int32
	_ = v2864
	var v2866 int32
	_ = v2866
	var v2867 int32
	_ = v2867
	var v2869 int32
	_ = v2869
	var v2870 int32
	_ = v2870
	var v2872 int32
	_ = v2872
	var v2873 int32
	_ = v2873
	var v2875 int32
	_ = v2875
	var v2876 int32
	_ = v2876
	var v2878 int32
	_ = v2878
	var v2879 int32
	_ = v2879
	var v2881 int32
	_ = v2881
	var v2882 int32
	_ = v2882
	var v2884 int32
	_ = v2884
	var v2885 int32
	_ = v2885
	var v2887 int32
	_ = v2887
	var v2888 int32
	_ = v2888
	var v2890 int32
	_ = v2890
	var v2891 int32
	_ = v2891
	var v2893 int32
	_ = v2893
	var v2894 int32
	_ = v2894
	var v2896 int32
	_ = v2896
	var v2897 int32
	_ = v2897
	var v2899 int32
	_ = v2899
	var v2900 int32
	_ = v2900
	var v2902 int32
	_ = v2902
	var v2905 int32
	_ = v2905
	var v2906 int32
	_ = v2906
	var v2908 int32
	_ = v2908
	var v2911 int32
	_ = v2911
	var v2912 int32
	_ = v2912
	var v2914 int32
	_ = v2914
	var v2918 int32
	_ = v2918
	var v2920 base.V128
	_ = v2920
	var v2921 int32
	_ = v2921
	var v2923 int32
	_ = v2923
	var v2925 base.V128
	_ = v2925
	var v2929 int32
	_ = v2929
	var v2942 int32
	_ = v2942
	var v2945 int32
	_ = v2945
	var v2946 int32
	_ = v2946
	var v2947 int32
	_ = v2947
	var v2950 int32
	_ = v2950
	var v2953 int32
	_ = v2953
	var v2954 int32
	_ = v2954
	var v2957 int32
	_ = v2957
	var v2958 int32
	_ = v2958
	var v2960 int32
	_ = v2960
	var v2963 int32
	_ = v2963
	var v2965 int32
	_ = v2965
	var v2969 int32
	_ = v2969
	var v2973 int32
	_ = v2973
	var v2974 int32
	_ = v2974
	var v2975 int32
	_ = v2975
	var v2976 int32
	_ = v2976
	var v2977 int32
	_ = v2977
	var v2981 int32
	_ = v2981
	var v2982 int32
	_ = v2982
	var v2984 int32
	_ = v2984
	var v2985 int32
	_ = v2985
	var v2987 int32
	_ = v2987
	var v2988 int32
	_ = v2988
	var v2990 int32
	_ = v2990
	var v2991 int32
	_ = v2991
	var v2993 int32
	_ = v2993
	var v2994 int32
	_ = v2994
	var v2996 int32
	_ = v2996
	var v2997 int32
	_ = v2997
	var v2999 int32
	_ = v2999
	var v3000 int32
	_ = v3000
	var v3002 int32
	_ = v3002
	var v3003 int32
	_ = v3003
	var v3005 int32
	_ = v3005
	var v3006 int32
	_ = v3006
	var v3008 int32
	_ = v3008
	var v3009 int32
	_ = v3009
	var v3011 int32
	_ = v3011
	var v3012 int32
	_ = v3012
	var v3014 int32
	_ = v3014
	var v3015 int32
	_ = v3015
	var v3017 int32
	_ = v3017
	var v3018 int32
	_ = v3018
	var v3020 int32
	_ = v3020
	var v3021 int32
	_ = v3021
	var v3023 int32
	_ = v3023
	var v3024 int32
	_ = v3024
	var v3026 int32
	_ = v3026
	var v3027 int32
	_ = v3027
	var v3029 int32
	_ = v3029
	var v3030 int32
	_ = v3030
	var v3032 int32
	_ = v3032
	var v3033 int32
	_ = v3033
	var v3035 int32
	_ = v3035
	var v3036 int32
	_ = v3036
	var v3038 int32
	_ = v3038
	var v3039 int32
	_ = v3039
	var v3041 int32
	_ = v3041
	var v3042 int32
	_ = v3042
	var v3044 int32
	_ = v3044
	var v3045 int32
	_ = v3045
	var v3047 int32
	_ = v3047
	var v3048 int32
	_ = v3048
	var v3050 int32
	_ = v3050
	var v3051 int32
	_ = v3051
	var v3053 int32
	_ = v3053
	var v3054 int32
	_ = v3054
	var v3056 int32
	_ = v3056
	var v3057 int32
	_ = v3057
	var v3059 int32
	_ = v3059
	var v3060 int32
	_ = v3060
	var v3062 int32
	_ = v3062
	var v3063 int32
	_ = v3063
	var v3065 int32
	_ = v3065
	var v3066 int32
	_ = v3066
	var v3068 int32
	_ = v3068
	var v3069 int32
	_ = v3069
	var v3071 int32
	_ = v3071
	var v3072 int32
	_ = v3072
	var v3074 int32
	_ = v3074
	var v3075 int32
	_ = v3075
	var v3077 int32
	_ = v3077
	var v3078 int32
	_ = v3078
	var v3080 int32
	_ = v3080
	var v3081 int32
	_ = v3081
	var v3083 int32
	_ = v3083
	var v3086 int32
	_ = v3086
	var v3087 int32
	_ = v3087
	var v3089 int32
	_ = v3089
	var v3092 int32
	_ = v3092
	var v3093 int32
	_ = v3093
	var v3095 int32
	_ = v3095
	var v3099 int32
	_ = v3099
	var v3101 base.V128
	_ = v3101
	var v3102 int32
	_ = v3102
	var v3104 int32
	_ = v3104
	var v3106 base.V128
	_ = v3106
	var v3111 int64
	_ = v3111
	var v3112 int64
	_ = v3112
	var v3114 int32
	_ = v3114
	var v3119 int32
	_ = v3119
	var v3121 int32
	_ = v3121
	var v3123 int32
	_ = v3123
	var v3124 int32
	_ = v3124
	var v3128 int32
	_ = v3128
	var v3130 int64
	_ = v3130
	var v3132 int32
	_ = v3132
	var v3134 int32
	_ = v3134
	var v3136 int32
	_ = v3136
	var v3137 int32
	_ = v3137
	var v3139 int32
	_ = v3139
	var v3143 int32
	_ = v3143
	var v3149 int32
	_ = v3149
	var v3150 int32
	_ = v3150
	var v3159 int32
	_ = v3159
	var v3166 int32
	_ = v3166
	var v3168 int32
	_ = v3168
	var v3169 int32
	_ = v3169
	var v3172 int32
	_ = v3172
	var v3176 int32
	_ = v3176
	var v3177 int64
	_ = v3177
	var v3183 int32
	_ = v3183
	var v3186 int32
	_ = v3186
	var v3189 int32
	_ = v3189
	var v3191 int32
	_ = v3191
	var v3196 int32
	_ = v3196
	var v3197 int32
	_ = v3197
	var v3200 int32
	_ = v3200
	var v3204 int32
	_ = v3204
	var v3208 int32
	_ = v3208
	var v3209 int32
	_ = v3209
	var v3212 int32
	_ = v3212
	var v3219 int32
	_ = v3219
	var v3225 int64
	_ = v3225
	var v3226 int64
	_ = v3226
	var v3229 int32
	_ = v3229
	var v3266 int32
	_ = v3266
	var v3271 int32
	_ = v3271
	var v3272 int32
	_ = v3272
	var v3273 int32
	_ = v3273
	var v3274 int32
	_ = v3274
	var v3275 int32
	_ = v3275
	var v3277 int32
	_ = v3277
	var v3278 int32
	_ = v3278
	var v3279 int32
	_ = v3279
	var v3280 int32
	_ = v3280
	var v3307 int32
	_ = v3307
	var v3308 int32
	_ = v3308
	var v3310 int32
	_ = v3310
	var v3311 int32
	_ = v3311
	var v3312 int32
	_ = v3312
	var v3313 int32
	_ = v3313
	var v3314 int32
	_ = v3314
	var v3315 int32
	_ = v3315
	var v3316 int32
	_ = v3316
	var v3336 int32
	_ = v3336
	var v3337 int32
	_ = v3337
	var v3346 int32
	_ = v3346
	var v3347 int32
	_ = v3347
	var v3357 int32
	_ = v3357
	var v3359 int32
	_ = v3359
	var v3360 int32
	_ = v3360
	var v3362 int32
	_ = v3362
	var v3363 int32
	_ = v3363
	var v3365 int32
	_ = v3365
	var v3367 int32
	_ = v3367
	var v3369 int32
	_ = v3369
	var v3372 int32
	_ = v3372
	var v3374 int32
	_ = v3374
	var v3375 int32
	_ = v3375
	var v3378 int32
	_ = v3378
	var v3379 int32
	_ = v3379
	var v3382 int32
	_ = v3382
	var v3384 int32
	_ = v3384
	var v3389 int32
	_ = v3389
	var v3395 int32
	_ = v3395
	var v3396 int32
	_ = v3396
	var v3398 int32
	_ = v3398
	var v3404 int32
	_ = v3404
	var v3409 int32
	_ = v3409
	var v3417 int32
	_ = v3417
	var v3420 int32
	_ = v3420
	var v3422 int32
	_ = v3422
	var v3429 int32
	_ = v3429
	var v3430 int32
	_ = v3430
	var v3444 int32
	_ = v3444
	var v3445 int32
	_ = v3445
	var v3449 int32
	_ = v3449
	var v3461 int32
	_ = v3461
	var v3462 int32
	_ = v3462
	var v3471 int32
	_ = v3471
	var v3477 int32
	_ = v3477
	var v3481 int32
	_ = v3481
	var v3483 int32
	_ = v3483
	var v3484 int32
	_ = v3484
	var v3485 int32
	_ = v3485
	var v3487 int32
	_ = v3487
	var v3489 int32
	_ = v3489
	var v3490 int32
	_ = v3490
	var v3494 int32
	_ = v3494
	var v3496 int32
	_ = v3496
	var v3498 int32
	_ = v3498
	var v3501 int32
	_ = v3501
	var v3503 int32
	_ = v3503
	var v3504 int32
	_ = v3504
	var v3507 int32
	_ = v3507
	var v3508 int32
	_ = v3508
	var v3511 int32
	_ = v3511
	var v3513 int32
	_ = v3513
	var v3518 int32
	_ = v3518
	var v3524 int32
	_ = v3524
	var v3527 int32
	_ = v3527
	var v3533 int32
	_ = v3533
	var v3538 int32
	_ = v3538
	var v3546 int32
	_ = v3546
	var v3550 int32
	_ = v3550
	var v3551 int32
	_ = v3551
	var v3558 int32
	_ = v3558
	var v3559 int32
	_ = v3559
	var v3568 int32
	_ = v3568
	var v3569 int32
	_ = v3569
	var v3573 int32
	_ = v3573
	var v3585 int32
	_ = v3585
	var v3586 int32
	_ = v3586
	var v3595 int32
	_ = v3595
	var v3601 int32
	_ = v3601
	var v3605 int32
	_ = v3605
	var v3607 int32
	_ = v3607
	var v3608 int32
	_ = v3608
	var v3609 int32
	_ = v3609
	var v3611 int32
	_ = v3611
	var v3613 int32
	_ = v3613
	var v3614 int32
	_ = v3614
	var v3618 int32
	_ = v3618
	var v3620 int32
	_ = v3620
	var v3622 int32
	_ = v3622
	var v3625 int32
	_ = v3625
	var v3627 int32
	_ = v3627
	var v3628 int32
	_ = v3628
	var v3631 int32
	_ = v3631
	var v3632 int32
	_ = v3632
	var v3635 int32
	_ = v3635
	var v3637 int32
	_ = v3637
	var v3642 int32
	_ = v3642
	var v3648 int32
	_ = v3648
	var v3651 int32
	_ = v3651
	var v3657 int32
	_ = v3657
	var v3662 int32
	_ = v3662
	var v3670 int32
	_ = v3670
	var v3674 int32
	_ = v3674
	var v3675 int32
	_ = v3675
	var v3682 int32
	_ = v3682
	var v3683 int32
	_ = v3683
	var v3692 int32
	_ = v3692
	var v3693 int32
	_ = v3693
	var v3697 int32
	_ = v3697
	var v3700 int32
	_ = v3700
	var v3702 int32
	_ = v3702
	var v3713 int32
	_ = v3713
	var v3716 int32
	_ = v3716
	var v3718 int32
	_ = v3718
	var v3729 int32
	_ = v3729
	var v3743 int32
	_ = v3743
	var v3745 int32
	_ = v3745
	var v3753 int32
	_ = v3753
	var v3754 int32
	_ = v3754
	var v3762 int32
	_ = v3762
	var v3763 int32
	_ = v3763
	var v3764 int32
	_ = v3764
	var v3770 int32
	_ = v3770
	var v3774 int32
	_ = v3774
	var v3775 int32
	_ = v3775
	var v3777 int32
	_ = v3777
	var v3792 int32
	_ = v3792
	var v3799 int32
	_ = v3799
	var v3802 int32
	_ = v3802
	var v3806 int32
	_ = v3806
	var v3807 int32
	_ = v3807
	var v3809 int32
	_ = v3809
	var v3814 int32
	_ = v3814
	var v3824 int32
	_ = v3824
	var v3840 float64
	_ = v3840
	var v3846 float64
	_ = v3846
	var v3857 int64
	_ = v3857
	var v3872 int32
	_ = v3872
	var v3874 int64
	_ = v3874
	var v3883 int64
	_ = v3883
	var v3888 int64
	_ = v3888
	var v3889 int32
	_ = v3889
	var v3891 int32
	_ = v3891
	var v3893 int32
	_ = v3893
	var v3897 float64
	_ = v3897
	var v3899 float64
	_ = v3899
	var v3912 float64
	_ = v3912
	var v3915 float64
	_ = v3915
	var v3920 float64
	_ = v3920
	var v3921 float64
	_ = v3921
	var v3922 float64
	_ = v3922
	var v3923 float64
	_ = v3923
	var v3928 float64
	_ = v3928
	var v3929 float64
	_ = v3929
	var v3930 float64
	_ = v3930
	var v3955 float64
	_ = v3955
	var v3967 float64
	_ = v3967
	var v3989 float64
	_ = v3989
	var v3992 float64
	_ = v3992
	var v3993 int32
	_ = v3993
	var v4003 int32
	_ = v4003
	var v4004 int32
	_ = v4004
	var v4005 int32
	_ = v4005
	var v4014 int64
	_ = v4014
	var v4020 base.V128
	_ = v4020
	var v4029 float32
	_ = v4029
	var v4039 float32
	_ = v4039
	var v4042 float32
	_ = v4042
	var v4045 float32
	_ = v4045
	var v4048 float32
	_ = v4048
	var v4049 float32
	_ = v4049
	var v4051 float32
	_ = v4051
	var v4053 float32
	_ = v4053
	var v4054 float32
	_ = v4054
	var v4055 float32
	_ = v4055
	var v4056 float64
	_ = v4056
	var v4057 int32
	_ = v4057
	var v4058 float32
	_ = v4058
	var v4060 int32
	_ = v4060
	var v4061 float32
	_ = v4061
	var v4062 float32
	_ = v4062
	var v4063 float64
	_ = v4063
	var v4064 int32
	_ = v4064
	var v4065 float32
	_ = v4065
	var v4087 int32
	_ = v4087
	var v4129 int32
	_ = v4129
	var v4166 int32
	_ = v4166
	var v4171 int32
	_ = v4171
	var v4172 int32
	_ = v4172
	var v4173 int32
	_ = v4173
	var v4174 int32
	_ = v4174
	var v4175 int32
	_ = v4175
	var v4177 int32
	_ = v4177
	var v4178 int32
	_ = v4178
	var v4180 int32
	_ = v4180
	var v4207 int32
	_ = v4207
	var v4210 int32
	_ = v4210
	var v4211 int32
	_ = v4211
	var v4212 int32
	_ = v4212
	var v4213 int32
	_ = v4213
	var v4214 int32
	_ = v4214
	var v4215 int32
	_ = v4215
	var v4216 int32
	_ = v4216
	var v4236 int32
	_ = v4236
	var v4246 int32
	_ = v4246
	var v4247 int32
	_ = v4247
	var v4257 int32
	_ = v4257
	var v4259 int32
	_ = v4259
	var v4260 int32
	_ = v4260
	var v4262 int32
	_ = v4262
	var v4263 int32
	_ = v4263
	var v4265 int32
	_ = v4265
	var v4267 int32
	_ = v4267
	var v4269 int32
	_ = v4269
	var v4272 int32
	_ = v4272
	var v4274 int32
	_ = v4274
	var v4275 int32
	_ = v4275
	var v4278 int32
	_ = v4278
	var v4279 int32
	_ = v4279
	var v4282 int32
	_ = v4282
	var v4284 int32
	_ = v4284
	var v4289 int32
	_ = v4289
	var v4295 int32
	_ = v4295
	var v4298 int32
	_ = v4298
	var v4304 int32
	_ = v4304
	var v4309 int32
	_ = v4309
	var v4317 int32
	_ = v4317
	var v4344 int32
	_ = v4344
	var v4349 int32
	_ = v4349
	var v4361 int32
	_ = v4361
	var v4371 int32
	_ = v4371
	var v4377 int32
	_ = v4377
	var v4381 int32
	_ = v4381
	var v4383 int32
	_ = v4383
	var v4384 int32
	_ = v4384
	var v4385 int32
	_ = v4385
	var v4387 int32
	_ = v4387
	var v4389 int32
	_ = v4389
	var v4390 int32
	_ = v4390
	var v4394 int32
	_ = v4394
	var v4396 int32
	_ = v4396
	var v4398 int32
	_ = v4398
	var v4401 int32
	_ = v4401
	var v4403 int32
	_ = v4403
	var v4404 int32
	_ = v4404
	var v4407 int32
	_ = v4407
	var v4408 int32
	_ = v4408
	var v4411 int32
	_ = v4411
	var v4413 int32
	_ = v4413
	var v4418 int32
	_ = v4418
	var v4424 int32
	_ = v4424
	var v4427 int32
	_ = v4427
	var v4433 int32
	_ = v4433
	var v4438 int32
	_ = v4438
	var v4446 int32
	_ = v4446
	var v4468 int32
	_ = v4468
	var v4473 int32
	_ = v4473
	var v4485 int32
	_ = v4485
	var v4495 int32
	_ = v4495
	var v4501 int32
	_ = v4501
	var v4505 int32
	_ = v4505
	var v4507 int32
	_ = v4507
	var v4508 int32
	_ = v4508
	var v4509 int32
	_ = v4509
	var v4511 int32
	_ = v4511
	var v4513 int32
	_ = v4513
	var v4514 int32
	_ = v4514
	var v4518 int32
	_ = v4518
	var v4520 int32
	_ = v4520
	var v4522 int32
	_ = v4522
	var v4525 int32
	_ = v4525
	var v4527 int32
	_ = v4527
	var v4528 int32
	_ = v4528
	var v4531 int32
	_ = v4531
	var v4532 int32
	_ = v4532
	var v4535 int32
	_ = v4535
	var v4537 int32
	_ = v4537
	var v4542 int32
	_ = v4542
	var v4548 int32
	_ = v4548
	var v4551 int32
	_ = v4551
	var v4557 int32
	_ = v4557
	var v4562 int32
	_ = v4562
	var v4570 int32
	_ = v4570
	var v4592 int32
	_ = v4592
	var v4597 int32
	_ = v4597
	var v4600 int32
	_ = v4600
	var v4602 int32
	_ = v4602
	var v4613 int32
	_ = v4613
	var v4616 int32
	_ = v4616
	var v4618 int32
	_ = v4618
	var v4629 int32
	_ = v4629
	var v4634 int32
	_ = v4634
	var v4645 int32
	_ = v4645
	var v4654 int32
	_ = v4654
	var v4662 int32
	_ = v4662
	var v4663 int32
	_ = v4663
	var v4664 int32
	_ = v4664
	var v4665 int32
	_ = v4665
	var v4679 int32
	_ = v4679
	var v4680 int32
	_ = v4680
	var v4683 int32
	_ = v4683
	var v4695 int32
	_ = v4695
	var v4696 int32
	_ = v4696
	var v4697 int32
	_ = v4697
	var v4701 int32
	_ = v4701
	var v4796 int32
	_ = v4796
	var v4797 int32
	_ = v4797
	var v4801 int32
	_ = v4801
	var v4807 int32
	_ = v4807
	var v4866 int32
	_ = v4866
	var v4871 int32
	_ = v4871
	var v4874 int32
	_ = v4874
	var v4885 int32
	_ = v4885
	var v4886 int32
	_ = v4886
	var v4887 int32
	_ = v4887
	var v4888 int32
	_ = v4888
	var v4892 int32
	_ = v4892
	var v4893 int32
	_ = v4893
	var v4897 int32
	_ = v4897
	var v4898 int32
	_ = v4898
	var v4903 int32
	_ = v4903
	var v4908 int32
	_ = v4908
	var v4909 int32
	_ = v4909
	var v4912 int64
	_ = v4912
	var v4913 int64
	_ = v4913
	var v4915 int64
	_ = v4915
	var v4916 int64
	_ = v4916
	var v4918 int64
	_ = v4918
	var v4922 int64
	_ = v4922
	var v4924 int64
	_ = v4924
	var v4928 int64
	_ = v4928
	var v4930 int64
	_ = v4930
	var v4934 int64
	_ = v4934
	var v4936 int64
	_ = v4936
	var v4940 int64
	_ = v4940
	var v4942 int64
	_ = v4942
	var v4946 int64
	_ = v4946
	var v4948 int64
	_ = v4948
	var v4952 int64
	_ = v4952
	var v4954 int64
	_ = v4954
	var v4958 int64
	_ = v4958
	var v4960 int64
	_ = v4960
	var v4964 int64
	_ = v4964
	var v4966 int64
	_ = v4966
	var v4970 int64
	_ = v4970
	var v4972 int64
	_ = v4972
	var v4976 int64
	_ = v4976
	var v4978 int64
	_ = v4978
	var v4982 int64
	_ = v4982
	var v4992 int32
	_ = v4992
	var v4994 int32
	_ = v4994
	var v5001 int32
	_ = v5001
	var v5005 int32
	_ = v5005
	v62 = m.G0
	v64 = v62 - int32(_a_F_VP8EncTokenLoop_0)
	m.G0 = v64
	v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[0])))
	v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v68 = *(*int32)(unsafe.Add(mBase, uint32(v67)+60))
	v69 = *(*int32)(unsafe.Add(mBase, uint32(v67)+108))
	v70 = base.F32_convert_i32_s(v69)
	v71 = *(*int32)(unsafe.Add(mBase, uint32(v67)+112))
	v72 = base.F32_convert_i32_s(v71)
	v73 = *(*float32)(unsafe.Add(mBase, uint32(v67)+4))
	if base.F32_gt(v73, v72) != 0 {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v75 = v72
	goto L3
L2:
	;
	v75 = v73
	goto L3
L3:
	;
	if base.F32_lt(v73, v70) != 0 {
		goto L4
	} else {
		goto L5
	}
L4:
	;
	v77 = v70
	goto L6
L5:
	;
	v77 = v75
	goto L6
L6:
	;
	v78 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v80 = v78 * v79
	v81 = m.G1
	v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3384))
	v88 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v81+int32(_a_F_VP8EncTokenLoop_1)+v84>>(uint(int32(4))%32)))))
	v90 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v91 = base.I32_div_s(v80*v88, v90)
	v92 = *(*int32)(unsafe.Add(mBase, uint32(v67)+16))
	v95 = *(*float32)(unsafe.Add(mBase, uint32(v67)+20))
	if base.F32_gt(v95, float32(0)) != 0 {
		goto L7
	} else {
		goto L8
	}
L7:
	;
	v100 = base.F64_promote_f32(v95)
	goto L9
L8:
	;
	v100 = float64(40)
	goto L9
L9:
	;
	if v92 != 0 {
		goto L10
	} else {
		goto L11
	}
L10:
	;
	v101 = base.F64_convert_i64_u(base.I64_extend_i32_s(v92))
	goto L12
L11:
	;
	v101 = v100
	goto L12
L12:
	;
	v108 = base.I64_extend_i32_s(v79) * base.I64_extend_i32_s(v78) * int64(384)
	v109 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[1])))
	v112 = l0 + int32(88)
	v116 = v112
	v121 = int32(-1)
	goto L15
L13:
	;
	m.G0 = v64 + int32(_a_F_VP8EncTokenLoop_0)
	return v5005
L14:
	;
	if int32(1) <= v68 {
		goto L41
	} else {
		goto L42
	}
L15:
	;
	v175 = v121 + int32(1)
	v176 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v176 <= v175 {
		goto L14
	} else {
		goto L17
	}
L16:
	;
	F_VP8BitWriterWipeOut(m, l0+int32(56))
	mBase = m.M
	v223 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v223 < int32(1) {
		goto L30
	} else {
		goto L31
	}
L17:
	;
	v178 = int32(0)
	v180 = base.Simd_g_const(&F_VP8EncTokenLoop__k0)
	base.Simd_g_v128_store(m, v116, v178, v180)
	v183 = base.Simd_g_const(&F_VP8EncTokenLoop__k1)
	base.Simd_g_v128_store(m, v116, int32(16), v183)
	if v91 == v178 {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	if v215 != 0 {
		v116 = v116 + int32(32)
		v121 = v175
		goto L15
	} else {
		goto L28
	}
L19:
	;
	v215 = int32(1)
	goto L18
L20:
	;
	v189 = int32(1024)
	if base.Ui32(v189) < base.Ui32(v91) {
		goto L22
	} else {
		goto L23
	}
L21:
	;
	v197 = *(*int32)(unsafe.Add(mBase, uint32(v116)+20))
	if v197 == int32(0) {
		goto L26
	} else {
		goto L27
	}
L22:
	;
	v192 = v91
	goto L24
L23:
	;
	v192 = v189
	goto L24
L24:
	;
	v193 = F_WebPSafeMalloc(m, int64(1), v192)
	mBase = m.M
	if v193 != 0 {
		goto L21
	} else {
		goto L25
	}
L25:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v116)+28)) = int32(1)
	v215 = int32(0)
	goto L18
L26:
	;
	v204 = *(*int32)(unsafe.Add(mBase, uint32(v116)+16))
	F_WebPSafeFree(m, v204)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v116)+24)) = v192
	*(*int32)(unsafe.Add(mBase, uint32(v116)+16)) = v193
	goto L19
L27:
	;
	v202 = *(*int32)(unsafe.Add(mBase, uint32(v116+int32(16))))
	v203 = F_memcpy(m, v193, v202, v197)
	mBase = m.M
	goto L26
L28:
	;
	goto L16
L29:
	;
	v242 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v244 = *(*int32)(unsafe.Add(mBase, uint32(v242)+92))
	if v244 != 0 {
		goto L36
	} else {
		goto L37
	}
L30:
	;
	goto L29
L31:
	;
	v230 = l0 + int32(88)
	v231 = int32(0)
	goto L32
L32:
	;
	F_VP8BitWriterWipeOut(m, v230)
	mBase = m.M
	v236 = v231 + int32(1)
	v237 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v236 < v237 {
		v230 = v230 + int32(32)
		v231 = v236
		goto L32
	} else {
		goto L34
	}
L33:
	;
	goto L30
L34:
	;
	goto L33
L35:
	;
	goto L38
L36:
	;
	goto L35
L37:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v242)+92)) = int32(1)
	goto L36
L38:
	;
	v5005 = int32(0)
	goto L13
L39:
	;
	v4866 = v64 + int32(880)
	v4871 = *(*int32)(unsafe.Add(mBase, uint32(v4866)+24))
	if v4807 == int32(0) {
		goto L410
	} else {
		goto L411
	}
L40:
	;
	if v92 != 0 {
		goto L329
	} else {
		goto L330
	}
L41:
	;
	v254 = v80 >> (uint(int32(3)) % 32)
	v255 = int32(96)
	if v255 < v254 {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v4087 = int32(40)
	goto L40
L43:
	;
	v258 = v254
	goto L45
L44:
	;
	v258 = v255
	goto L45
L45:
	;
	v262 = l0 + int32(344)
	v316 = int32(40)
	v326 = v68
	v329 = v77
	v335 = float32(10)
	v341 = v316
	v371 = float64(0)
	v372 = int32(1)
	v373 = v77
	goto L46
L46:
	;
	v384 = v326 + int32(-1)
	v385 = int32(1)
	if base.F64_le(base.F64_promote_f32(base.F32_abs(v335)), float64(0.4)) != 0 {
		v395 = v385
		goto L48
	} else {
		goto L49
	}
L47:
	;
	v4087 = v647
	goto L40
L48:
	;
	v397 = v64 + int32(880)
	v398 = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(v397))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v397)+24)) = l0
	v403 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v397)+280)) = v403
	v405 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	*(*int32)(unsafe.Add(mBase, uint32(v397)+296)) = v405
	v407 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v397)+304)) = v407
	*(*int32)(unsafe.Add(mBase, uint32(v397)+32)) = l0 + int32(88)
	v412 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[4])))
	*(*int32)(unsafe.Add(mBase, uint32(v397)+36)) = v412
	v416 = int32(-32)
	v417 = (v64 + int32(1327)) & v416
	*(*int32)(unsafe.Add(mBase, uint32(v397)+8)) = v417
	v422 = (v64 + int32(1240)) & v416
	*(*int32)(unsafe.Add(mBase, uint32(v397)+308)) = v422
	*(*int32)(unsafe.Add(mBase, uint32(v397)+20)) = v417 + int32(1536)
	*(*int32)(unsafe.Add(mBase, uint32(v397)+16)) = v417 + int32(1024)
	*(*int32)(unsafe.Add(mBase, uint32(v397)+12)) = v417 + int32(512)
	*(*int32)(unsafe.Add(mBase, uint32(v397)+316)) = v422 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(v397)+312)) = v422 + int32(32)
	v439 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v397)+40)) = v439
	v441 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[6])))
	*(*int32)(unsafe.Add(mBase, uint32(v397)+28)) = v441
	v443 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v397)+320)) = v443
	v445 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v422)+47)) = uint8(v445)
	v447 = *(*int32)(unsafe.Add(mBase, uint32(v397)+312))
	v448 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v447+v448))) = uint8(v445)
	v452 = *(*int32)(unsafe.Add(mBase, uint32(v397)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v452+v448))) = uint8(v445)
	v457 = *(*int32)(unsafe.Add(mBase, uint32(v397)+308))
	v458 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v457))) = v458
	*(*int64)(unsafe.Add(mBase, uint32(v457+int32(8)))) = v458
	v464 = *(*int32)(unsafe.Add(mBase, uint32(v397)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v464))) = v458
	v467 = *(*int32)(unsafe.Add(mBase, uint32(v397)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v467))) = v458
	*(*int32)(unsafe.Add(mBase, uint32(v397)+160)) = v398
	v472 = *(*int32)(unsafe.Add(mBase, uint32(v397)+304))
	if v472 == v398 {
		goto L52
	} else {
		goto L53
	}
L49:
	;
	if v384 == int32(0) {
		v395 = v385
		goto L48
	} else {
		goto L50
	}
L50:
	;
	v392 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[8])))
	v395 = base.B2i32(v392 == int32(0))
	goto L48
L51:
	;
	F_SetLoopParams(m, l0, v329)
	mBase = m.M
	v513 = base.I32_div_s(v341, v326+int32(1))
	if v395 == int32(0) {
		goto L56
	} else {
		goto L57
	}
L52:
	;
	v477 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v478 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v479 = v477 * v478
	*(*int32)(unsafe.Add(mBase, uint32(v397)+292)) = v479
	*(*int32)(unsafe.Add(mBase, uint32(v397)+288)) = v479
	v482 = *(*int32)(unsafe.Add(mBase, uint32(v397)+24))
	v483 = *(*int32)(unsafe.Add(mBase, uint32(v482)+uint32(_c_F_VP8EncTokenLoop[7])))
	v485 = *(*int32)(unsafe.Add(mBase, uint32(v482)+40))
	v488 = F_memset(m, v483, int32(127), v485<<(uint(int32(5))%32))
	mBase = m.M
	v489 = *(*int32)(unsafe.Add(mBase, uint32(v482)+uint32(_c_F_VP8EncTokenLoop[5])))
	v490 = int32(0)
	v491 = *(*int32)(unsafe.Add(mBase, uint32(v482)+40))
	v494 = F_memset(m, v489, v490, v491<<(uint(int32(2))%32))
	mBase = m.M
	v495 = *(*int32)(unsafe.Add(mBase, uint32(v482)+uint32(_c_F_VP8EncTokenLoop[3])))
	if v495 == v490 {
		goto L54
	} else {
		goto L55
	}
L53:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v397)+300)) = int32(0)
	goto L52
L54:
	;
	v505 = int32(0)
	v507 = F_memset(m, v64+int32(1048), v505, int32(96))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v397)+284)) = v505
	goto L51
L55:
	;
	v499 = *(*int32)(unsafe.Add(mBase, uint32(v482)+40))
	v502 = F_memset(m, v495, int32(0), v499<<(uint(int32(2))%32))
	mBase = m.M
	goto L54
L56:
	;
	v647 = v341 - v513
	if v262 == int32(0) {
		goto L76
	} else {
		goto L77
	}
L57:
	;
	goto L60
L58:
	;
	v640 = *(*int32)(unsafe.Add(mBase, uint32(v64+int32(880))+280))
	if v640 == int32(0) {
		goto L73
	} else {
		goto L74
	}
L60:
	;
	base.MemoryFill(m, l0+int32(_a_F_VP8EncTokenLoop_2), int32(0), int32(_a_F_VP8EncTokenLoop_3))
	goto L58
L72:
	;
	goto L56
L73:
	;
	goto L72
L74:
	;
	v645 = F_memset(m, v640, int32(0), int32(2048))
	mBase = m.M
	F_VP8SSIMDspInit(m)
	mBase = m.M
	goto L73
L75:
	;
	v679 = int64(0)
	v691 = v258
	v736 = v679
	v737 = v679
	goto L86
L76:
	;
	goto L75
L77:
	;
	v652 = *(*int32)(unsafe.Add(mBase, uint32(v262)))
	if v652 == int32(0) {
		goto L78
	} else {
		goto L79
	}
L78:
	;
	v663 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v262)+20)) = v663
	*(*int64)(unsafe.Add(mBase, uint32(v262)+8)) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v262))) = v663
	v669 = *(*int32)(unsafe.Add(mBase, uint32(v262)+16))
	v670 = int32(_a_F_VP8EncTokenLoop_4)
	if v670 < v669 {
		goto L83
	} else {
		goto L84
	}
L79:
	;
	v656 = v652
	goto L80
L80:
	;
	v658 = *(*int32)(unsafe.Add(mBase, uint32(v656)))
	F_WebPSafeFree(m, v656)
	mBase = m.M
	if v658 != 0 {
		v656 = v658
		goto L80
	} else {
		goto L82
	}
L81:
	;
	goto L78
L82:
	;
	goto L81
L83:
	;
	v673 = v669
	goto L85
L84:
	;
	v673 = v670
	goto L85
L85:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v262)+16)) = v673
	*(*int32)(unsafe.Add(mBase, uint32(v262)+4)) = v262
	goto L76
L86:
	;
	v744 = int32(0)
	F_VP8IteratorImport(m, v64+int32(880), v744)
	mBase = m.M
	if v691 <= v744 {
		goto L89
	} else {
		goto L90
	}
L87:
	;
	v3225 = int64(*(*int32)(unsafe.Add(mBase, uint32(l0)+32)))
	v3226 = v3112 + v3225
	if v92 == int32(0) {
		goto L221
	} else {
		goto L222
	}
L88:
	;
	v2230 = v64 + int32(880)
	v2231 = F_VP8Decimate(m, v2230, v64, v109)
	mBase = m.M
	v2232 = *(*int32)(unsafe.Add(mBase, uint32(v64)+904))
	v2237 = *(*int32)(unsafe.Add(mBase, uint32(v2230)+40))
	v2240 = *(*int32)(unsafe.Add(mBase, uint32(v2237+int32(-4))))
	v2241 = *(*int32)(unsafe.Add(mBase, uint32(v2237)))
	v2244 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+124)) = int32(base.Ui32(v2241)>>(uint(int32(24))%32)) & v2244
	v2247 = int32(23)
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+120)) = int32(base.Ui32(v2241)>>(uint(v2247)%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+116)) = int32(base.Ui32(v2241)>>(uint(int32(22))%32)) & v2244
	v2257 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+112)) = int32(base.Ui32(v2241)>>(uint(v2257)%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+108)) = int32(base.Ui32(v2241)>>(uint(int32(18))%32)) & v2244
	v2267 = int32(15)
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+104)) = int32(base.Ui32(v2241)>>(uint(v2267)%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+100)) = int32(base.Ui32(v2241)>>(uint(int32(14))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+96)) = int32(base.Ui32(v2241)>>(uint(int32(13))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+92)) = int32(base.Ui32(v2241)>>(uint(int32(12))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+156)) = int32(base.Ui32(v2240)>>(uint(v2247)%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+152)) = int32(base.Ui32(v2240)>>(uint(int32(21))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+148)) = int32(base.Ui32(v2240)>>(uint(v2257)%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+144)) = int32(base.Ui32(v2240)>>(uint(int32(17))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+140)) = int32(base.Ui32(v2240)>>(uint(v2267)%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+136)) = int32(base.Ui32(v2240)>>(uint(int32(11))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+132)) = int32(base.Ui32(v2240)>>(uint(int32(7))%32)) & v2244
	*(*int32)(unsafe.Add(mBase, uint32(v2230)+128)) = int32(base.Ui32(v2240)>>(uint(int32(3))%32)) & v2244
	goto L176
L89:
	;
	v750 = int32(0)
	v787 = l0 + int32(_a_F_VP8EncTokenLoop_2)
	v792 = l0 + int32(_a_F_VP8EncTokenLoop_5)
	v793 = l0 + int32(3442)
	v794 = l0 + int32(_a_F_VP8EncTokenLoop_6)
	v795 = l0 + int32(3431)
	v796 = v750
	v798 = v787
	v799 = v750
	v801 = v750
	goto L92
L90:
	;
	v2228 = v691 + int32(-1)
	goto L88
L91:
	;
	v1294 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9])))
	if v1294 == int32(0) {
		goto L135
	} else {
		goto L136
	}
L92:
	;
	v828 = v799
	v831 = v792
	v832 = v793
	v833 = v794
	v834 = v795
	v835 = v796
	v836 = v798
	v837 = int32(0)
	goto L94
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = v1213
	goto L91
L94:
	;
	v857 = v828
	v867 = int32(0)
	v868 = v835
	goto L96
L95:
	;
	v1237 = int32(1056)
	v1239 = int32(264)
	v1250 = v801 + int32(1)
	if v1250 != int32(4) {
		v792 = v792 + v1237
		v793 = v793 + v1239
		v794 = v794 + v1237
		v795 = v795 + v1239
		v796 = v796 + v1239
		v798 = v798 + v1237
		v799 = v1213
		v801 = v1250
		goto L92
	} else {
		goto L133
	}
L96:
	;
	v878 = *(*int32)(unsafe.Add(mBase, uint32(v836+v867)))
	v880 = int32(base.Ui32(v878) >> (uint(int32(16)) % 32))
	v881 = m.G125
	v883 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v881+v868))))
	v884 = m.G126
	v886 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v884+v868))))
	v888 = v878 & int32(_a_F_VP8EncTokenLoop_7)
	if v888 != 0 {
		goto L99
	} else {
		goto L100
	}
L97:
	;
	v982 = v965
	v992 = v833
	v998 = int32(0)
	goto L108
L98:
	;
	v896 = m.G79
	v899 = v880 - v888
	v900 = int32(1)
	v903 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v896+v883<<(uint(v900)%32)))))
	v905 = int32(255)
	v910 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v896+(v883^v905)<<(uint(v900)%32)))))
	v916 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v896+v886<<(uint(v900)%32)))))
	v919 = v895 & v905
	v925 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v896+(v919^v905)<<(uint(v900)%32)))))
	v930 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v896+v919<<(uint(v900)%32)))))
	v938 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v896+(v886^v905)<<(uint(v900)%32)))))
	goto L101
L99:
	;
	v890 = int32(255)
	v893 = base.I32_div_u_s(v888*v890, v880)
	v895 = v890 - v893
	goto L98
L100:
	;
	v895 = int32(255)
	goto L98
L101:
	;
	goto L103
L103:
	;
	if v899*v903+v888*v910+v916 <= v888*v925+v899*v930+v938+int32(2048) {
		goto L105
	} else {
		goto L106
	}
L104:
	;
	v970 = v867 + int32(4)
	if v970 != int32(44) {
		v857 = v965
		v867 = v970
		v868 = v868 + int32(1)
		goto L96
	} else {
		goto L107
	}
L105:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v787+v868+int32(-1056)))) = uint8(v883)
	v965 = v857
	goto L104
L106:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v787+v868+int32(-1056)))) = uint8(v895)
	v965 = v857 | base.B2i32(v895 != v883)
	goto L104
L107:
	;
	goto L97
L108:
	;
	v1002 = *(*int32)(unsafe.Add(mBase, uint32(v992)))
	v1004 = int32(base.Ui32(v1002) >> (uint(int32(16)) % 32))
	v1005 = m.G125
	v1006 = v835 + v998
	v1008 = int32(11)
	v1010 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1005+v1006+v1008))))
	v1011 = m.G126
	v1015 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1011+v1006+v1008))))
	v1017 = v1002 & int32(_a_F_VP8EncTokenLoop_7)
	if v1017 != 0 {
		goto L111
	} else {
		goto L112
	}
L109:
	;
	v1106 = v1089
	v1116 = v831
	v1122 = int32(0)
	goto L120
L110:
	;
	v1025 = m.G79
	v1028 = v1004 - v1017
	v1029 = int32(1)
	v1032 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1025+v1010<<(uint(v1029)%32)))))
	v1034 = int32(255)
	v1039 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1025+(v1010^v1034)<<(uint(v1029)%32)))))
	v1045 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1025+v1015<<(uint(v1029)%32)))))
	v1048 = v1024 & v1034
	v1054 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1025+(v1048^v1034)<<(uint(v1029)%32)))))
	v1059 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1025+v1048<<(uint(v1029)%32)))))
	v1067 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1025+(v1015^v1034)<<(uint(v1029)%32)))))
	goto L113
L111:
	;
	v1019 = int32(255)
	v1022 = base.I32_div_u_s(v1017*v1019, v1004)
	v1024 = v1019 - v1022
	goto L110
L112:
	;
	v1024 = int32(255)
	goto L110
L113:
	;
	goto L115
L115:
	;
	if v1017*v1054+v1028*v1059+v1067+int32(2048) < v1028*v1032+v1017*v1039+v1045 {
		goto L117
	} else {
		goto L118
	}
L116:
	;
	v1094 = v998 + int32(1)
	if v1094 != int32(11) {
		v982 = v1089
		v992 = v992 + int32(4)
		v998 = v1094
		goto L108
	} else {
		goto L119
	}
L117:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v834+v998))) = uint8(v1024)
	v1089 = v982 | base.B2i32(v1024 != v1010)
	goto L116
L118:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v834+v998))) = uint8(v1010)
	v1089 = v982
	goto L116
L119:
	;
	goto L109
L120:
	;
	v1126 = *(*int32)(unsafe.Add(mBase, uint32(v1116)))
	v1128 = int32(base.Ui32(v1126) >> (uint(int32(16)) % 32))
	v1129 = m.G125
	v1130 = v835 + v1122
	v1132 = int32(22)
	v1134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1129+v1130+v1132))))
	v1135 = m.G126
	v1139 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1135+v1130+v1132))))
	v1141 = v1126 & int32(_a_F_VP8EncTokenLoop_7)
	if v1141 != 0 {
		goto L123
	} else {
		goto L124
	}
L121:
	;
	v1221 = int32(132)
	v1223 = int32(33)
	v1234 = v837 + int32(1)
	if v1234 != int32(8) {
		v828 = v1213
		v831 = v831 + v1221
		v832 = v832 + v1223
		v833 = v833 + v1221
		v834 = v834 + v1223
		v835 = v835 + v1223
		v836 = v836 + v1221
		v837 = v1234
		goto L94
	} else {
		goto L132
	}
L122:
	;
	v1149 = m.G79
	v1152 = v1128 - v1141
	v1153 = int32(1)
	v1156 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1149+v1134<<(uint(v1153)%32)))))
	v1158 = int32(255)
	v1163 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1149+(v1134^v1158)<<(uint(v1153)%32)))))
	v1169 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1149+v1139<<(uint(v1153)%32)))))
	v1172 = v1148 & v1158
	v1178 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1149+(v1172^v1158)<<(uint(v1153)%32)))))
	v1183 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1149+v1172<<(uint(v1153)%32)))))
	v1191 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1149+(v1139^v1158)<<(uint(v1153)%32)))))
	goto L125
L123:
	;
	v1143 = int32(255)
	v1146 = base.I32_div_u_s(v1141*v1143, v1128)
	v1148 = v1143 - v1146
	goto L122
L124:
	;
	v1148 = int32(255)
	goto L122
L125:
	;
	goto L127
L127:
	;
	if v1141*v1178+v1152*v1183+v1191+int32(2048) < v1152*v1156+v1141*v1163+v1169 {
		goto L129
	} else {
		goto L130
	}
L128:
	;
	v1218 = v1122 + int32(1)
	if v1218 != int32(11) {
		v1106 = v1213
		v1116 = v1116 + int32(4)
		v1122 = v1218
		goto L120
	} else {
		goto L131
	}
L129:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v832+v1122))) = uint8(v1148)
	v1213 = v1106 | base.B2i32(v1148 != v1134)
	goto L128
L130:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v832+v1122))) = uint8(v1134)
	v1213 = v1106
	goto L128
L131:
	;
	goto L121
L132:
	;
	goto L95
L133:
	;
	goto L93
L134:
	;
	v2228 = v258
	goto L88
L135:
	;
	goto L134
L136:
	;
	v1309 = m.G81
	v1310 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+15)))
	v1311 = int32(408)
	v1313 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+14)))
	v1316 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+13)))
	v1319 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+12)))
	v1322 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+11)))
	v1325 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+10)))
	v1328 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+9)))
	v1331 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+8)))
	v1334 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+7)))
	v1337 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+6)))
	v1340 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+5)))
	v1343 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+4)))
	v1346 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+3)))
	v1349 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+2)))
	v1352 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309)+1)))
	v1355 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1309))))
	v1360 = l0 + int32(3444)
	v1361 = l0 + int32(3433)
	v1362 = l0 + int32(3422)
	v1383 = int32(0)
	goto L137
L137:
	;
	v1402 = l0 + int32(_a_F_VP8EncTokenLoop_8) + v1383*int32(3264)
	v1434 = v1360
	v1435 = v1361
	v1436 = v1362
	v1437 = int32(0)
	goto L139
L138:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = int32(0)
	goto L135
L139:
	;
	v1450 = l0 + int32(3420) + v1383*int32(264) + v1437*int32(33)
	v1451 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1450)+1)))
	v1452 = int32(1)
	v1455 = v1402 + v1437*int32(408)
	v1456 = m.G79
	v1460 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1456+v1451<<(uint(v1452)%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1455))) = uint16(v1460)
	v1467 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1456+(v1451^int32(255))<<(uint(v1452)%32)))))
	v1500 = v1452
	goto L141
L140:
	;
	v2046 = l0 + int32(_a_F_VP8EncTokenLoop_9) + v1383*int32(192)
	v2047 = v1402 + v1310*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+180)) = v2047
	v2049 = v1402 + v1313*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+168)) = v2049
	v2051 = v1402 + v1316*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+156)) = v2051
	v2053 = v1402 + v1319*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+144)) = v2053
	v2055 = v1402 + v1322*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+132)) = v2055
	v2057 = v1402 + v1325*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+120)) = v2057
	v2059 = v1402 + v1328*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+108)) = v2059
	v2061 = v1402 + v1331*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+96)) = v2061
	v2063 = v1402 + v1334*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+84)) = v2063
	v2065 = v1402 + v1337*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+72)) = v2065
	v2067 = v1402 + v1340*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+60)) = v2067
	v2069 = v1402 + v1343*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+48)) = v2069
	v2071 = v1402 + v1346*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+36)) = v2071
	v2073 = v1402 + v1349*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+24)) = v2073
	v2075 = v1402 + v1352*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+12)) = v2075
	v2077 = v1402 + v1355*v1311
	*(*int32)(unsafe.Add(mBase, uint32(v2046))) = v2077
	v2079 = int32(272)
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+188)) = v2047 + v2079
	v2082 = int32(136)
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+184)) = v2047 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+176)) = v2049 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+172)) = v2049 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+164)) = v2051 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+160)) = v2051 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+152)) = v2053 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+148)) = v2053 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+140)) = v2055 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+136)) = v2055 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+128)) = v2057 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+124)) = v2057 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+116)) = v2059 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+112)) = v2059 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+104)) = v2061 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+100)) = v2061 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+92)) = v2063 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+88)) = v2063 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+80)) = v2065 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+76)) = v2065 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+68)) = v2067 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+64)) = v2067 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+56)) = v2069 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+52)) = v2069 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+44)) = v2071 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+40)) = v2071 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+32)) = v2073 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+28)) = v2073 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+20)) = v2075 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+16)) = v2075 + v2082
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+8)) = v2077 + v2079
	*(*int32)(unsafe.Add(mBase, uint32(v2046)+4)) = v2077 + v2082
	v2175 = int32(264)
	v2182 = v1383 + int32(1)
	if v2182 != int32(4) {
		v1360 = v1360 + v2175
		v1361 = v1361 + v2175
		v1362 = v1362 + v2175
		v1383 = v2182
		goto L137
	} else {
		goto L175
	}
L141:
	;
	v1511 = m.G1
	v1516 = v1500<<(uint(int32(2))%32) + (v1511 + int32(_a_F_VP8EncTokenLoop_10)) + int32(-4)
	v1517 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1516))))
	if v1517 != 0 {
		goto L144
	} else {
		goto L145
	}
L142:
	;
	v1640 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1450)+12)))
	v1641 = int32(1)
	v1642 = m.G79
	v1646 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1642+v1640<<(uint(v1641)%32)))))
	v1647 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1450)+11)))
	v1648 = int32(255)
	v1653 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1642+(v1647^v1648)<<(uint(v1641)%32)))))
	v1654 = v1646 + v1653
	*(*uint16)(unsafe.Add(mBase, uint32(v1455)+136)) = uint16(v1654)
	v1661 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1642+(v1640^v1648)<<(uint(v1641)%32)))))
	v1697 = v1641
	goto L152
L143:
	;
	v1631 = int32(1)
	v1634 = v1467 + v1626
	*(*uint16)(unsafe.Add(mBase, uint32(v1455+v1500<<(uint(v1631)%32)))) = uint16(v1634)
	v1637 = v1500 + v1631
	if v1637 != int32(68) {
		v1500 = v1637
		goto L141
	} else {
		goto L151
	}
L144:
	;
	v1519 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1516)+2)))
	v1528 = v1517
	v1555 = v1519
	v1557 = int32(0)
	v1558 = v1436
	goto L146
L145:
	;
	v1626 = int32(0)
	goto L143
L146:
	;
	if v1528&int32(1) == int32(0) {
		v1580 = v1557
		goto L148
	} else {
		goto L149
	}
L147:
	;
	v1626 = v1580
	goto L143
L148:
	;
	v1582 = int32(1)
	if base.Ui32(v1582) < base.Ui32(v1528) {
		v1528 = int32(base.Ui32(v1528) >> (uint(v1582) % 32))
		v1555 = int32(base.Ui32(v1555) >> (uint(v1582) % 32))
		v1557 = v1580
		v1558 = v1558 + v1582
		goto L146
	} else {
		goto L150
	}
L149:
	;
	v1566 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1558))))
	v1567 = m.G79
	v1569 = int32(1)
	v1578 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1567+(v1566^(int32(0)-v1555&v1569)&int32(255))<<(uint(v1569)%32)))))
	v1580 = v1557 + v1578
	goto L148
L150:
	;
	goto L147
L151:
	;
	goto L142
L152:
	;
	v1708 = m.G1
	v1713 = v1697<<(uint(int32(2))%32) + (v1708 + int32(_a_F_VP8EncTokenLoop_10)) + int32(-4)
	v1714 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1713))))
	if v1714 != 0 {
		goto L155
	} else {
		goto L156
	}
L153:
	;
	v1837 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1450)+23)))
	v1838 = int32(1)
	v1839 = m.G79
	v1843 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1839+v1837<<(uint(v1838)%32)))))
	v1844 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1450)+22)))
	v1845 = int32(255)
	v1850 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1839+(v1844^v1845)<<(uint(v1838)%32)))))
	v1851 = v1843 + v1850
	*(*uint16)(unsafe.Add(mBase, uint32(v1455)+272)) = uint16(v1851)
	v1858 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1839+(v1837^v1845)<<(uint(v1838)%32)))))
	v1894 = v1838
	goto L163
L154:
	;
	v1828 = int32(1)
	v1831 = v1653 + v1661 + v1823
	*(*uint16)(unsafe.Add(mBase, uint32(v1455+int32(136)+v1697<<(uint(v1828)%32)))) = uint16(v1831)
	v1834 = v1697 + v1828
	if v1834 != int32(68) {
		v1697 = v1834
		goto L152
	} else {
		goto L162
	}
L155:
	;
	v1716 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1713)+2)))
	v1725 = v1714
	v1752 = v1716
	v1754 = int32(0)
	v1755 = v1435
	goto L157
L156:
	;
	v1823 = int32(0)
	goto L154
L157:
	;
	if v1725&int32(1) == int32(0) {
		v1777 = v1754
		goto L159
	} else {
		goto L160
	}
L158:
	;
	v1823 = v1777
	goto L154
L159:
	;
	v1779 = int32(1)
	if base.Ui32(v1779) < base.Ui32(v1725) {
		v1725 = int32(base.Ui32(v1725) >> (uint(v1779) % 32))
		v1752 = int32(base.Ui32(v1752) >> (uint(v1779) % 32))
		v1754 = v1777
		v1755 = v1755 + v1779
		goto L157
	} else {
		goto L161
	}
L160:
	;
	v1763 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1755))))
	v1764 = m.G79
	v1766 = int32(1)
	v1775 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1764+(v1763^(int32(0)-v1752&v1766)&int32(255))<<(uint(v1766)%32)))))
	v1777 = v1754 + v1775
	goto L159
L161:
	;
	goto L158
L162:
	;
	goto L153
L163:
	;
	v1905 = m.G1
	v1910 = v1894<<(uint(int32(2))%32) + (v1905 + int32(_a_F_VP8EncTokenLoop_10)) + int32(-4)
	v1911 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1910))))
	if v1911 != 0 {
		goto L166
	} else {
		goto L167
	}
L164:
	;
	v2034 = int32(33)
	v2041 = v1437 + int32(1)
	if v2041 != int32(8) {
		v1434 = v1434 + v2034
		v1435 = v1435 + v2034
		v1436 = v1436 + v2034
		v1437 = v2041
		goto L139
	} else {
		goto L174
	}
L165:
	;
	v2025 = int32(1)
	v2028 = v1850 + v1858 + v2020
	*(*uint16)(unsafe.Add(mBase, uint32(v1455+int32(272)+v1894<<(uint(v2025)%32)))) = uint16(v2028)
	v2031 = v1894 + v2025
	if v2031 != int32(68) {
		v1894 = v2031
		goto L163
	} else {
		goto L173
	}
L166:
	;
	v1913 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1910)+2)))
	v1922 = v1911
	v1949 = v1913
	v1951 = int32(0)
	v1952 = v1434
	goto L168
L167:
	;
	v2020 = int32(0)
	goto L165
L168:
	;
	if v1922&int32(1) == int32(0) {
		v1974 = v1951
		goto L170
	} else {
		goto L171
	}
L169:
	;
	v2020 = v1974
	goto L165
L170:
	;
	v1976 = int32(1)
	if base.Ui32(v1976) < base.Ui32(v1922) {
		v1922 = int32(base.Ui32(v1922) >> (uint(v1976) % 32))
		v1949 = int32(base.Ui32(v1949) >> (uint(v1976) % 32))
		v1951 = v1974
		v1952 = v1952 + v1976
		goto L168
	} else {
		goto L172
	}
L171:
	;
	v1960 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1952))))
	v1961 = m.G79
	v1963 = int32(1)
	v1972 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1961+(v1960^(int32(0)-v1949&v1963)&int32(255))<<(uint(v1963)%32)))))
	v1974 = v1951 + v1972
	goto L170
L172:
	;
	goto L169
L173:
	;
	goto L164
L174:
	;
	goto L140
L175:
	;
	goto L138
L176:
	;
	v2327 = *(*int32)(unsafe.Add(mBase, uint32(v64)+908))
	v2328 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2327))))
	if v2328&int32(3) != int32(1) {
		goto L178
	} else {
		goto L179
	}
L177:
	;
	v2420 = m.G78
	v2421 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1008))
	v2422 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2424 = v64 + int32(_a_F_VP8EncTokenLoop_11)
	v2425 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2425].(func(*base.Module, int32, int32))(m, v64+int32(72), v2424)
	mBase = m.M
	v2430 = F_VP8RecordCoeffTokens(m, v2421+v2422, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2430
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2430
	v2433 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2436 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2436].(func(*base.Module, int32, int32))(m, v64+int32(104), v2424)
	mBase = m.M
	v2441 = F_VP8RecordCoeffTokens(m, v2433+v2430, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2441
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2441
	v2444 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2447 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2447].(func(*base.Module, int32, int32))(m, v64+int32(136), v2424)
	mBase = m.M
	v2452 = F_VP8RecordCoeffTokens(m, v2444+v2441, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2452
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2452
	v2455 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2458 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2458].(func(*base.Module, int32, int32))(m, v64+int32(168), v2424)
	mBase = m.M
	v2463 = F_VP8RecordCoeffTokens(m, v2455+v2452, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2463
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2463
	v2466 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2467 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1012))
	v2470 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2470].(func(*base.Module, int32, int32))(m, v64+int32(200), v2424)
	mBase = m.M
	v2475 = F_VP8RecordCoeffTokens(m, v2467+v2466, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2475
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2475
	v2478 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2481 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2481].(func(*base.Module, int32, int32))(m, v64+int32(232), v2424)
	mBase = m.M
	v2486 = F_VP8RecordCoeffTokens(m, v2478+v2475, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2486
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2486
	v2489 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2492 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2492].(func(*base.Module, int32, int32))(m, v64+int32(264), v2424)
	mBase = m.M
	v2497 = F_VP8RecordCoeffTokens(m, v2489+v2486, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2497
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2497
	v2500 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2503 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2503].(func(*base.Module, int32, int32))(m, v64+int32(296), v2424)
	mBase = m.M
	v2508 = F_VP8RecordCoeffTokens(m, v2500+v2497, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2508
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2508
	v2511 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2512 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1016))
	v2515 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2515].(func(*base.Module, int32, int32))(m, v64+int32(328), v2424)
	mBase = m.M
	v2520 = F_VP8RecordCoeffTokens(m, v2512+v2511, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2520
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2520
	v2523 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2526 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2526].(func(*base.Module, int32, int32))(m, v64+int32(360), v2424)
	mBase = m.M
	v2531 = F_VP8RecordCoeffTokens(m, v2523+v2520, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2531
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2531
	v2534 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2537 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2537].(func(*base.Module, int32, int32))(m, v64+int32(392), v2424)
	mBase = m.M
	v2542 = F_VP8RecordCoeffTokens(m, v2534+v2531, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2542
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2542
	v2545 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2548 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2548].(func(*base.Module, int32, int32))(m, v64+int32(424), v2424)
	mBase = m.M
	v2553 = F_VP8RecordCoeffTokens(m, v2545+v2542, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2553
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2553
	v2556 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2557 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1020))
	v2560 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2560].(func(*base.Module, int32, int32))(m, v64+int32(456), v2424)
	mBase = m.M
	v2565 = F_VP8RecordCoeffTokens(m, v2557+v2556, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2565
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2565
	v2568 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2571 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2571].(func(*base.Module, int32, int32))(m, v64+int32(488), v2424)
	mBase = m.M
	v2576 = F_VP8RecordCoeffTokens(m, v2568+v2565, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2576
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2576
	v2579 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2582 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2582].(func(*base.Module, int32, int32))(m, v64+int32(520), v2424)
	mBase = m.M
	v2587 = F_VP8RecordCoeffTokens(m, v2579+v2576, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2587
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2587
	v2590 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2593 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2593].(func(*base.Module, int32, int32))(m, v64+int32(552), v2424)
	mBase = m.M
	v2598 = F_VP8RecordCoeffTokens(m, v2590+v2587, v2424, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2598
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2598
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2232 + int32(_a_F_VP8EncTokenLoop_12)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2232 + int32(_a_F_VP8EncTokenLoop_13)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2232 + int32(3948)
	goto L183
L178:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2232 + int32(_a_F_VP8EncTokenLoop_14)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2232 + int32(_a_F_VP8EncTokenLoop_15)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2232 + int32(_a_F_VP8EncTokenLoop_16)
	goto L182
L179:
	;
	v2333 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1040))
	v2334 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1004))
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2232 + int32(_a_F_VP8EncTokenLoop_17)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2232 + int32(_a_F_VP8EncTokenLoop_18)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2232 + int32(3684)
	goto L180
L180:
	;
	v2360 = v64 + int32(_a_F_VP8EncTokenLoop_11)
	v2361 = m.G78
	v2362 = *(*int32)(unsafe.Add(mBase, uint32(v2361)))
	m.T0[v2362].(func(*base.Module, int32, int32))(m, v64+v316, v2360)
	mBase = m.M
	v2367 = F_VP8RecordCoeffTokens(m, v2333+v2334, v2360, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1004)) = v2367
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1040)) = v2367
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2232 + int32(_a_F_VP8EncTokenLoop_9)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2232 + int32(_a_F_VP8EncTokenLoop_2)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2232 + int32(3420)
	goto L181
L181:
	;
	goto L177
L182:
	;
	goto L177
L183:
	;
	v2625 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1024))
	v2626 = *(*int32)(unsafe.Add(mBase, uint32(v64)+988))
	v2628 = v64 + int32(_a_F_VP8EncTokenLoop_11)
	v2629 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2629].(func(*base.Module, int32, int32))(m, v64+int32(584), v2628)
	mBase = m.M
	v2634 = F_VP8RecordCoeffTokens(m, v2625+v2626, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+988)) = v2634
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1024)) = v2634
	v2637 = *(*int32)(unsafe.Add(mBase, uint32(v64)+992))
	v2640 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2640].(func(*base.Module, int32, int32))(m, v64+int32(616), v2628)
	mBase = m.M
	v2645 = F_VP8RecordCoeffTokens(m, v2637+v2634, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+992)) = v2645
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1024)) = v2645
	v2648 = *(*int32)(unsafe.Add(mBase, uint32(v64)+988))
	v2649 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1028))
	v2652 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2652].(func(*base.Module, int32, int32))(m, v64+int32(648), v2628)
	mBase = m.M
	v2657 = F_VP8RecordCoeffTokens(m, v2649+v2648, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+988)) = v2657
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1028)) = v2657
	v2660 = *(*int32)(unsafe.Add(mBase, uint32(v64)+992))
	v2663 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2663].(func(*base.Module, int32, int32))(m, v64+int32(680), v2628)
	mBase = m.M
	v2668 = F_VP8RecordCoeffTokens(m, v2660+v2657, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+992)) = v2668
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1028)) = v2668
	v2671 = *(*int32)(unsafe.Add(mBase, uint32(v64)+996))
	v2672 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1032))
	v2675 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2675].(func(*base.Module, int32, int32))(m, v64+int32(712), v2628)
	mBase = m.M
	v2680 = F_VP8RecordCoeffTokens(m, v2672+v2671, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+996)) = v2680
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1032)) = v2680
	v2683 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1000))
	v2686 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2686].(func(*base.Module, int32, int32))(m, v64+int32(744), v2628)
	mBase = m.M
	v2691 = F_VP8RecordCoeffTokens(m, v2683+v2680, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1000)) = v2691
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1032)) = v2691
	v2694 = *(*int32)(unsafe.Add(mBase, uint32(v64)+996))
	v2695 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1036))
	v2698 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2698].(func(*base.Module, int32, int32))(m, v64+int32(776), v2628)
	mBase = m.M
	v2703 = F_VP8RecordCoeffTokens(m, v2695+v2694, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+996)) = v2703
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1036)) = v2703
	v2706 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1000))
	v2709 = *(*int32)(unsafe.Add(mBase, uint32(v2420)))
	m.T0[v2709].(func(*base.Module, int32, int32))(m, v64+int32(808), v2628)
	mBase = m.M
	v2714 = F_VP8RecordCoeffTokens(m, v2706+v2703, v2628, v262)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1000)) = v2714
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1036)) = v2714
	v2718 = v64 + int32(880)
	v2719 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+40))
	v2720 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+96))
	v2723 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+92))
	v2727 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+100))
	v2731 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+104))
	v2735 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+108))
	v2739 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+112))
	v2743 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+116))
	v2747 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+120))
	v2751 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+124))
	v2755 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+128))
	v2759 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+132))
	v2763 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+136))
	v2767 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+144))
	v2771 = *(*int32)(unsafe.Add(mBase, uint32(v2718)+152))
	*(*int32)(unsafe.Add(mBase, uint32(v2719))) = v2720<<(uint(int32(13))%32) | v2723<<(uint(int32(12))%32) | v2727<<(uint(int32(14))%32) | v2731<<(uint(int32(15))%32) | v2735<<(uint(int32(18))%32) | v2739<<(uint(int32(19))%32) | v2743<<(uint(int32(22))%32) | v2747<<(uint(int32(23))%32) | v2751<<(uint(int32(24))%32) | v2755<<(uint(int32(3))%32) | v2759<<(uint(int32(7))%32) | v2763<<(uint(int32(11))%32) | v2767<<(uint(int32(17))%32) | v2771<<(uint(int32(21))%32)
	goto L184
L184:
	;
	v2776 = *(*int32)(unsafe.Add(mBase, uint32(l0)+364))
	if v2776 == int32(0) {
		goto L185
	} else {
		goto L186
	}
L185:
	;
	v2785 = *(*int64)(unsafe.Add(mBase, uint32(v64)))
	v2786 = *(*int64)(unsafe.Add(mBase, uint32(v64)+16))
	if v395 != 0 {
		goto L191
	} else {
		goto L192
	}
L186:
	;
	v2779 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v2781 = *(*int32)(unsafe.Add(mBase, uint32(v2779)+92))
	if v2781 != 0 {
		goto L188
	} else {
		goto L189
	}
L187:
	;
	v4807 = int32(0)
	goto L39
L188:
	;
	goto L187
L189:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2779)+92)) = int32(1)
	goto L188
L190:
	;
	v3111 = v2785 + v737
	v3112 = v2786 + v736
	v3114 = v64 + int32(880)
	v3119 = *(*int32)(unsafe.Add(mBase, uint32(v3114)))
	v3121 = v3119 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v3114))) = v3121
	v3123 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+24))
	v3124 = *(*int32)(unsafe.Add(mBase, uint32(v3123)+40))
	if v3121 != v3124 {
		goto L213
	} else {
		goto L214
	}
L191:
	;
	v2929 = v64 + int32(880)
	F_StoreSideInfo(m, v2929)
	mBase = m.M
	F_VP8StoreFilterStats(m, v2929)
	mBase = m.M
	F_VP8IteratorExport(m, v2929)
	mBase = m.M
	v2942 = int32(1)
	if v513 == int32(0) {
		v2965 = v2942
		goto L199
	} else {
		goto L200
	}
L192:
	;
	v2788 = v64 + int32(880)
	v2792 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+12))
	v2793 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+4))
	v2794 = *(*int32)(unsafe.Add(mBase, uint32(v2788)))
	v2795 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+24))
	v2796 = *(*int32)(unsafe.Add(mBase, uint32(v2795)+40))
	if v2796+int32(-1) <= v2794 {
		goto L194
	} else {
		goto L195
	}
L193:
	;
	goto L190
L194:
	;
	v2914 = *(*int32)(unsafe.Add(mBase, uint32(v2795)+44))
	if v2914+int32(-1) <= v2793 {
		goto L196
	} else {
		goto L197
	}
L195:
	;
	v2800 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2801 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2800))) = uint8(v2801)
	v2803 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2804 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+47)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2803)+1)) = uint8(v2804)
	v2806 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2807 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+79)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2806)+2)) = uint8(v2807)
	v2809 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2810 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+111)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2809)+3)) = uint8(v2810)
	v2812 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2813 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+143)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2812)+4)) = uint8(v2813)
	v2815 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2816 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+175)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2815)+5)) = uint8(v2816)
	v2818 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2819 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+207)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2818)+6)) = uint8(v2819)
	v2821 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2822 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+239)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2821)+7)) = uint8(v2822)
	v2824 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2825 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+271)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2824)+8)) = uint8(v2825)
	v2827 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2828 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+303)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2827)+9)) = uint8(v2828)
	v2830 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2831 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+335)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2830)+10)) = uint8(v2831)
	v2833 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2834 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+367)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2833)+11)) = uint8(v2834)
	v2836 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2837 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+399)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2836)+12)) = uint8(v2837)
	v2839 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2840 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+431)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2839)+13)) = uint8(v2840)
	v2842 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2843 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+463)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2842)+14)) = uint8(v2843)
	v2845 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2846 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+495)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2845)+15)) = uint8(v2846)
	v2848 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2849 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+23)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2848))) = uint8(v2849)
	v2851 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2852 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+31)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2851))) = uint8(v2852)
	v2854 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2855 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+55)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2854)+1)) = uint8(v2855)
	v2857 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2858 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+63)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2857)+1)) = uint8(v2858)
	v2860 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2861 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+87)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2860)+2)) = uint8(v2861)
	v2863 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2864 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+95)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2863)+2)) = uint8(v2864)
	v2866 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2867 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+119)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2866)+3)) = uint8(v2867)
	v2869 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2870 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+127)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2869)+3)) = uint8(v2870)
	v2872 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2873 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+151)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2872)+4)) = uint8(v2873)
	v2875 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2876 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+159)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2875)+4)) = uint8(v2876)
	v2878 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2879 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+183)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2878)+5)) = uint8(v2879)
	v2881 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2882 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+191)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2881)+5)) = uint8(v2882)
	v2884 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2885 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+215)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2884)+6)) = uint8(v2885)
	v2887 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2888 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+223)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2887)+6)) = uint8(v2888)
	v2890 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2891 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+247)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2890)+7)) = uint8(v2891)
	v2893 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2894 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2792)+255)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2893)+7)) = uint8(v2894)
	v2896 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+308))
	v2897 = int32(-1)
	v2899 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+320))
	v2900 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2899)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2896+v2897))) = uint8(v2900)
	v2902 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+312))
	v2905 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+324))
	v2906 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2905)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2902+v2897))) = uint8(v2906)
	v2908 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+316))
	v2911 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+324))
	v2912 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2911)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2908+v2897))) = uint8(v2912)
	goto L194
L196:
	;
	goto L193
L197:
	;
	v2918 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+320))
	v2920 = base.Simd_g_v128_load(m, v2792, int32(480))
	v2921 = int32(0)
	base.Simd_g_v128_store(m, v2918, v2921, v2920)
	v2923 = *(*int32)(unsafe.Add(mBase, uint32(v2788)+324))
	v2925 = base.Simd_g_v128_load(m, v2792, int32(240))
	base.Simd_g_v128_store(m, v2923, v2921, v2925)
	goto L196
L198:
	;
	v2969 = v64 + int32(880)
	v2973 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+12))
	v2974 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+4))
	v2975 = *(*int32)(unsafe.Add(mBase, uint32(v2969)))
	v2976 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+24))
	v2977 = *(*int32)(unsafe.Add(mBase, uint32(v2976)+40))
	if v2977+int32(-1) <= v2975 {
		goto L206
	} else {
		goto L207
	}
L199:
	;
	goto L198
L200:
	;
	v2945 = *(*int32)(unsafe.Add(mBase, uint32(v2929)+24))
	v2946 = *(*int32)(unsafe.Add(mBase, uint32(v2945)+4))
	v2947 = *(*int32)(unsafe.Add(mBase, uint32(v2946)+96))
	if v2947 == int32(0) {
		v2965 = v2942
		goto L199
	} else {
		goto L201
	}
L201:
	;
	v2950 = *(*int32)(unsafe.Add(mBase, uint32(v2929)+292))
	if int32(0) < v2950 {
		goto L203
	} else {
		goto L204
	}
L202:
	;
	v2963 = F_WebPReportProgress(m, v2946, v2960, v2945+int32(368))
	mBase = m.M
	v2965 = v2963
	goto L199
L203:
	;
	v2954 = *(*int32)(unsafe.Add(mBase, uint32(v2929)+288))
	v2957 = base.I32_div_s((v2950-v2954)*v513, v2950)
	v2958 = *(*int32)(unsafe.Add(mBase, uint32(v2929)+296))
	v2960 = v2957 + v2958
	goto L202
L204:
	;
	v2953 = *(*int32)(unsafe.Add(mBase, uint32(v2929)+296))
	v2960 = v2953
	goto L202
L205:
	;
	if v2965 != 0 {
		goto L190
	} else {
		goto L210
	}
L206:
	;
	v3095 = *(*int32)(unsafe.Add(mBase, uint32(v2976)+44))
	if v3095+int32(-1) <= v2974 {
		goto L208
	} else {
		goto L209
	}
L207:
	;
	v2981 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v2982 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2981))) = uint8(v2982)
	v2984 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v2985 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+47)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2984)+1)) = uint8(v2985)
	v2987 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v2988 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+79)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2987)+2)) = uint8(v2988)
	v2990 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v2991 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+111)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2990)+3)) = uint8(v2991)
	v2993 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v2994 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+143)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2993)+4)) = uint8(v2994)
	v2996 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v2997 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+175)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2996)+5)) = uint8(v2997)
	v2999 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3000 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+207)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2999)+6)) = uint8(v3000)
	v3002 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3003 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+239)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3002)+7)) = uint8(v3003)
	v3005 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3006 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+271)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3005)+8)) = uint8(v3006)
	v3008 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3009 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+303)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3008)+9)) = uint8(v3009)
	v3011 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3012 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+335)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3011)+10)) = uint8(v3012)
	v3014 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3015 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+367)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3014)+11)) = uint8(v3015)
	v3017 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3018 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+399)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3017)+12)) = uint8(v3018)
	v3020 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3021 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+431)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3020)+13)) = uint8(v3021)
	v3023 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3024 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+463)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3023)+14)) = uint8(v3024)
	v3026 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3027 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+495)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3026)+15)) = uint8(v3027)
	v3029 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3030 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+23)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3029))) = uint8(v3030)
	v3032 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3033 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+31)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3032))) = uint8(v3033)
	v3035 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3036 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+55)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3035)+1)) = uint8(v3036)
	v3038 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3039 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+63)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3038)+1)) = uint8(v3039)
	v3041 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3042 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+87)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3041)+2)) = uint8(v3042)
	v3044 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3045 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+95)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3044)+2)) = uint8(v3045)
	v3047 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3048 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+119)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3047)+3)) = uint8(v3048)
	v3050 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3051 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+127)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3050)+3)) = uint8(v3051)
	v3053 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3054 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+151)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3053)+4)) = uint8(v3054)
	v3056 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3057 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+159)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3056)+4)) = uint8(v3057)
	v3059 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3060 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+183)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3059)+5)) = uint8(v3060)
	v3062 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3063 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+191)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3062)+5)) = uint8(v3063)
	v3065 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3066 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+215)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3065)+6)) = uint8(v3066)
	v3068 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3069 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+223)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3068)+6)) = uint8(v3069)
	v3071 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3072 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+247)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3071)+7)) = uint8(v3072)
	v3074 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3075 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2973)+255)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3074)+7)) = uint8(v3075)
	v3077 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+308))
	v3078 = int32(-1)
	v3080 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+320))
	v3081 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3080)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3077+v3078))) = uint8(v3081)
	v3083 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+312))
	v3086 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+324))
	v3087 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3086)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3083+v3078))) = uint8(v3087)
	v3089 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+316))
	v3092 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+324))
	v3093 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3092)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3089+v3078))) = uint8(v3093)
	goto L206
L208:
	;
	goto L205
L209:
	;
	v3099 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+320))
	v3101 = base.Simd_g_v128_load(m, v2973, int32(480))
	v3102 = int32(0)
	base.Simd_g_v128_store(m, v3099, v3102, v3101)
	v3104 = *(*int32)(unsafe.Add(mBase, uint32(v2969)+324))
	v3106 = base.Simd_g_v128_load(m, v2973, int32(240))
	base.Simd_g_v128_store(m, v3104, v3102, v3106)
	goto L208
L210:
	;
	v4807 = int32(0)
	goto L39
L211:
	;
	if int32(1) < v3219 {
		v691 = v2228
		v736 = v3112
		v737 = v3111
		goto L86
	} else {
		goto L219
	}
L212:
	;
	v3219 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+288))
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+288)) = v3219 + int32(-1)
	goto L211
L213:
	;
	v3196 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+36))
	v3197 = int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+36)) = v3196 + v3197
	v3200 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+28)) = v3200 + v3197
	v3204 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+40))
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+40)) = v3204 + v3197
	v3208 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+320))
	v3209 = int32(16)
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+320)) = v3208 + v3209
	v3212 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+324))
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+324)) = v3212 + v3209
	goto L212
L214:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3114))) = int32(0)
	v3128 = *(*int32)(unsafe.Add(mBase, uint32(v3123)+uint32(_c_F_VP8EncTokenLoop[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+40)) = v3128
	v3130 = *(*int64)(unsafe.Add(mBase, uint32(v3123)+uint32(_c_F_VP8EncTokenLoop[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v3114)+320)) = v3130
	v3132 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+4))
	v3134 = v3132 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+4)) = v3134
	v3136 = *(*int32)(unsafe.Add(mBase, uint32(v3123)+uint32(_c_F_VP8EncTokenLoop[4])))
	v3137 = *(*int32)(unsafe.Add(mBase, uint32(v3123)+48))
	v3139 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+36)) = v3136 + v3134*v3137<<(uint(v3139)%32)
	v3143 = *(*int32)(unsafe.Add(mBase, uint32(v3123)+uint32(_c_F_VP8EncTokenLoop[6])))
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+28)) = v3143 + v3134*v3121<<(uint(v3139)%32)
	v3149 = *(*int32)(unsafe.Add(mBase, uint32(v3123)+52))
	v3150 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+32)) = v3123 + (v3149+v3150)&v3134<<(uint(int32(5))%32) + int32(88)
	v3159 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+316))
	if v3150 < v3132 {
		goto L215
	} else {
		goto L216
	}
L215:
	;
	v3166 = int32(-127)
	goto L217
L216:
	;
	v3166 = int32(127)
	goto L217
L217:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3159+v3150))) = uint8(v3166)
	v3168 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+312))
	v3169 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v3168+v3169))) = uint8(v3166)
	v3172 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v3172+v3169))) = uint8(v3166)
	v3176 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+308))
	v3177 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v3176))) = v3177
	*(*int64)(unsafe.Add(mBase, uint32(v3176+int32(8)))) = v3177
	v3183 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v3183))) = v3177
	v3186 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v3186))) = v3177
	v3189 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+160)) = v3189
	v3191 = *(*int32)(unsafe.Add(mBase, uint32(v3114)+304))
	if v3191 == v3189 {
		goto L212
	} else {
		goto L218
	}
L218:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3114)+300)) = int32(0)
	goto L212
L219:
	;
	goto L87
L220:
	;
	v3993 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[8])))
	if v3993 < int32(1) {
		goto L299
	} else {
		goto L300
	}
L221:
	;
	v3840 = float64(99)
	if v108 == int64(0) {
		v3992 = v3840
		goto L220
	} else {
		goto L283
	}
L222:
	;
	v3229 = int32(0)
	v3266 = l0 + int32(_a_F_VP8EncTokenLoop_2)
	v3271 = l0 + int32(_a_F_VP8EncTokenLoop_5)
	v3272 = l0 + int32(3442)
	v3273 = l0 + int32(_a_F_VP8EncTokenLoop_6)
	v3274 = l0 + int32(3431)
	v3275 = v3229
	v3277 = v3266
	v3278 = v3229
	v3279 = v3229
	v3280 = v3229
	goto L224
L223:
	;
	v3743 = *(*int32)(unsafe.Add(mBase, uint32(v262)))
	if v3743 != 0 {
		goto L268
	} else {
		goto L269
	}
L224:
	;
	v3307 = v3278
	v3308 = v3279
	v3310 = v3271
	v3311 = v3272
	v3312 = v3273
	v3313 = v3274
	v3314 = v3275
	v3315 = v3277
	v3316 = int32(0)
	goto L226
L225:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = v3692
	goto L223
L226:
	;
	v3336 = v3307
	v3337 = v3308
	v3346 = int32(0)
	v3347 = v3314
	goto L228
L227:
	;
	v3716 = int32(1056)
	v3718 = int32(264)
	v3729 = v3280 + int32(1)
	if v3729 != int32(4) {
		v3271 = v3271 + v3716
		v3272 = v3272 + v3718
		v3273 = v3273 + v3716
		v3274 = v3274 + v3718
		v3275 = v3275 + v3718
		v3277 = v3277 + v3716
		v3278 = v3692
		v3279 = v3693
		v3280 = v3729
		goto L224
	} else {
		goto L265
	}
L228:
	;
	v3357 = *(*int32)(unsafe.Add(mBase, uint32(v3315+v3346)))
	v3359 = int32(base.Ui32(v3357) >> (uint(int32(16)) % 32))
	v3360 = m.G125
	v3362 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3360+v3347))))
	v3363 = m.G126
	v3365 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3363+v3347))))
	v3367 = v3357 & int32(_a_F_VP8EncTokenLoop_7)
	if v3367 != 0 {
		goto L231
	} else {
		goto L232
	}
L229:
	;
	v3461 = v3444
	v3462 = v3445
	v3471 = v3312
	v3477 = int32(0)
	goto L240
L230:
	;
	v3375 = m.G79
	v3378 = v3359 - v3367
	v3379 = int32(1)
	v3382 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+v3362<<(uint(v3379)%32)))))
	v3384 = int32(255)
	v3389 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+(v3362^v3384)<<(uint(v3379)%32)))))
	v3395 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+v3365<<(uint(v3379)%32)))))
	v3396 = v3378*v3382 + v3367*v3389 + v3395
	v3398 = v3374 & v3384
	v3404 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+(v3398^v3384)<<(uint(v3379)%32)))))
	v3409 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+v3398<<(uint(v3379)%32)))))
	v3417 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+(v3365^v3384)<<(uint(v3379)%32)))))
	v3420 = v3367*v3404 + v3378*v3409 + v3417 + int32(2048)
	if v3420 < v3396 {
		goto L233
	} else {
		goto L234
	}
L231:
	;
	v3369 = int32(255)
	v3372 = base.I32_div_u_s(v3367*v3369, v3359)
	v3374 = v3369 - v3372
	goto L230
L232:
	;
	v3374 = int32(255)
	goto L230
L233:
	;
	v3422 = int32(-1)
	goto L235
L234:
	;
	v3422 = int32(0)
	goto L235
L235:
	;
	v3429 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3375+(v3365^v3422)&int32(255)<<(uint(int32(1))%32)))))
	v3430 = v3337 + v3429
	if v3396 <= v3420 {
		goto L237
	} else {
		goto L238
	}
L236:
	;
	v3449 = v3346 + int32(4)
	if v3449 != int32(44) {
		v3336 = v3444
		v3337 = v3445
		v3346 = v3449
		v3347 = v3347 + int32(1)
		goto L228
	} else {
		goto L239
	}
L237:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3266+v3347+int32(-1056)))) = uint8(v3362)
	v3444 = v3336
	v3445 = v3430
	goto L236
L238:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3266+v3347+int32(-1056)))) = uint8(v3374)
	v3444 = v3336 | base.B2i32(v3374 != v3362)
	v3445 = v3430 + int32(2048)
	goto L236
L239:
	;
	goto L229
L240:
	;
	v3481 = *(*int32)(unsafe.Add(mBase, uint32(v3471)))
	v3483 = int32(base.Ui32(v3481) >> (uint(int32(16)) % 32))
	v3484 = m.G125
	v3485 = v3314 + v3477
	v3487 = int32(11)
	v3489 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3484+v3485+v3487))))
	v3490 = m.G126
	v3494 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3490+v3485+v3487))))
	v3496 = v3481 & int32(_a_F_VP8EncTokenLoop_7)
	if v3496 != 0 {
		goto L243
	} else {
		goto L244
	}
L241:
	;
	v3585 = v3568
	v3586 = v3569
	v3595 = v3310
	v3601 = int32(0)
	goto L252
L242:
	;
	v3504 = m.G79
	v3507 = v3483 - v3496
	v3508 = int32(1)
	v3511 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+v3489<<(uint(v3508)%32)))))
	v3513 = int32(255)
	v3518 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+(v3489^v3513)<<(uint(v3508)%32)))))
	v3524 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+v3494<<(uint(v3508)%32)))))
	v3527 = v3503 & v3513
	v3533 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+(v3527^v3513)<<(uint(v3508)%32)))))
	v3538 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+v3527<<(uint(v3508)%32)))))
	v3546 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+(v3494^v3513)<<(uint(v3508)%32)))))
	v3550 = base.B2i32(v3496*v3533+v3507*v3538+v3546+int32(2048) < v3507*v3511+v3496*v3518+v3524)
	if v3496*v3533+v3507*v3538+v3546+int32(2048) < v3507*v3511+v3496*v3518+v3524 {
		goto L245
	} else {
		goto L246
	}
L243:
	;
	v3498 = int32(255)
	v3501 = base.I32_div_u_s(v3496*v3498, v3483)
	v3503 = v3498 - v3501
	goto L242
L244:
	;
	v3503 = int32(255)
	goto L242
L245:
	;
	v3551 = int32(-1)
	goto L247
L246:
	;
	v3551 = int32(0)
	goto L247
L247:
	;
	v3558 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3504+(v3494^v3551)&int32(255)<<(uint(int32(1))%32)))))
	v3559 = v3462 + v3558
	if v3496*v3533+v3507*v3538+v3546+int32(2048) < v3507*v3511+v3496*v3518+v3524 {
		goto L249
	} else {
		goto L250
	}
L248:
	;
	v3573 = v3477 + int32(1)
	if v3573 != int32(11) {
		v3461 = v3568
		v3462 = v3569
		v3471 = v3471 + int32(4)
		v3477 = v3573
		goto L240
	} else {
		goto L251
	}
L249:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3313+v3477))) = uint8(v3503)
	v3568 = v3461 | base.B2i32(v3503 != v3489)
	v3569 = v3559 + int32(2048)
	goto L248
L250:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3313+v3477))) = uint8(v3489)
	v3568 = v3461
	v3569 = v3559
	goto L248
L251:
	;
	goto L241
L252:
	;
	v3605 = *(*int32)(unsafe.Add(mBase, uint32(v3595)))
	v3607 = int32(base.Ui32(v3605) >> (uint(int32(16)) % 32))
	v3608 = m.G125
	v3609 = v3314 + v3601
	v3611 = int32(22)
	v3613 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3608+v3609+v3611))))
	v3614 = m.G126
	v3618 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3614+v3609+v3611))))
	v3620 = v3605 & int32(_a_F_VP8EncTokenLoop_7)
	if v3620 != 0 {
		goto L255
	} else {
		goto L256
	}
L253:
	;
	v3700 = int32(132)
	v3702 = int32(33)
	v3713 = v3316 + int32(1)
	if v3713 != int32(8) {
		v3307 = v3692
		v3308 = v3693
		v3310 = v3310 + v3700
		v3311 = v3311 + v3702
		v3312 = v3312 + v3700
		v3313 = v3313 + v3702
		v3314 = v3314 + v3702
		v3315 = v3315 + v3700
		v3316 = v3713
		goto L226
	} else {
		goto L264
	}
L254:
	;
	v3628 = m.G79
	v3631 = v3607 - v3620
	v3632 = int32(1)
	v3635 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+v3613<<(uint(v3632)%32)))))
	v3637 = int32(255)
	v3642 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+(v3613^v3637)<<(uint(v3632)%32)))))
	v3648 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+v3618<<(uint(v3632)%32)))))
	v3651 = v3627 & v3637
	v3657 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+(v3651^v3637)<<(uint(v3632)%32)))))
	v3662 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+v3651<<(uint(v3632)%32)))))
	v3670 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+(v3618^v3637)<<(uint(v3632)%32)))))
	v3674 = base.B2i32(v3620*v3657+v3631*v3662+v3670+int32(2048) < v3631*v3635+v3620*v3642+v3648)
	if v3620*v3657+v3631*v3662+v3670+int32(2048) < v3631*v3635+v3620*v3642+v3648 {
		goto L257
	} else {
		goto L258
	}
L255:
	;
	v3622 = int32(255)
	v3625 = base.I32_div_u_s(v3620*v3622, v3607)
	v3627 = v3622 - v3625
	goto L254
L256:
	;
	v3627 = int32(255)
	goto L254
L257:
	;
	v3675 = int32(-1)
	goto L259
L258:
	;
	v3675 = int32(0)
	goto L259
L259:
	;
	v3682 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3628+(v3618^v3675)&int32(255)<<(uint(int32(1))%32)))))
	v3683 = v3586 + v3682
	if v3620*v3657+v3631*v3662+v3670+int32(2048) < v3631*v3635+v3620*v3642+v3648 {
		goto L261
	} else {
		goto L262
	}
L260:
	;
	v3697 = v3601 + int32(1)
	if v3697 != int32(11) {
		v3585 = v3692
		v3586 = v3693
		v3595 = v3595 + int32(4)
		v3601 = v3697
		goto L252
	} else {
		goto L263
	}
L261:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3311+v3601))) = uint8(v3627)
	v3692 = v3585 | base.B2i32(v3627 != v3613)
	v3693 = v3683 + int32(2048)
	goto L260
L262:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3311+v3601))) = uint8(v3613)
	v3692 = v3585
	v3693 = v3683
	goto L260
L263:
	;
	goto L253
L264:
	;
	goto L227
L265:
	;
	goto L225
L266:
	;
	v3992 = base.F64_convert_i64_u(int64(base.Ui64(v3226+base.I64_extend_i32_s(v3693)+base.I64_extend_i32_u(v3824)+int64(1024))>>(uint(int64(11))%64)) + int64(30))
	goto L220
L267:
	;
	goto L266
L268:
	;
	v3745 = *(*int32)(unsafe.Add(mBase, uint32(v262)+16))
	v3753 = v3743
	v3754 = int32(0)
	goto L270
L269:
	;
	v3824 = int32(0)
	goto L267
L270:
	;
	v3762 = *(*int32)(unsafe.Add(mBase, uint32(v3753)))
	if v3762 != 0 {
		v3764 = int32(0)
		goto L272
	} else {
		goto L273
	}
L271:
	;
	v3824 = v3814
	goto L267
L272:
	;
	if v3745 <= v3764 {
		v3814 = v3754
		goto L274
	} else {
		goto L275
	}
L273:
	;
	v3763 = *(*int32)(unsafe.Add(mBase, uint32(v262)+12))
	v3764 = v3763
	goto L272
L274:
	;
	if v3762 != 0 {
		v3753 = v3762
		v3754 = v3814
		goto L270
	} else {
		goto L282
	}
L275:
	;
	v3770 = v3754
	v3774 = v3753 + (v3745<<(uint(int32(1))%32) + int32(2))
	v3775 = v3745
	goto L276
L276:
	;
	v3777 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3774))))
	if v3777&int32(_a_F_VP8EncTokenLoop_19) == int32(0) {
		goto L279
	} else {
		goto L280
	}
L277:
	;
	v3814 = v3807
	goto L274
L278:
	;
	v3802 = m.G79
	v3806 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3802+v3799<<(uint(int32(1))%32)))))
	v3807 = v3770 + v3806
	v3809 = v3775 + int32(-1)
	if v3764 < v3809 {
		v3770 = v3807
		v3774 = v3774 + int32(-2)
		v3775 = v3809
		goto L276
	} else {
		goto L281
	}
L279:
	;
	v3792 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3420)+v3777&int32(_a_F_VP8EncTokenLoop_20)))))
	v3799 = v3792 ^ int32(base.Ui32(base.I32_extend16_s(v3777))>>(uint(int32(15))%32))&int32(255)
	goto L278
L280:
	;
	v3799 = (v3777<<(uint(int32(16))%32)>>(uint(int32(31))%32) ^ v3777) & int32(255)
	goto L278
L281:
	;
	goto L277
L282:
	;
	goto L271
L283:
	;
	if v3111 == int64(0) {
		v3992 = v3840
		goto L220
	} else {
		goto L284
	}
L284:
	;
	v3846 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v108), float64(65025)), base.F64_convert_i64_u(v3111))
	v3857 = base.I64_reinterpret_f64(v3846)
	if int64(4503599627370495) < v3857 {
		goto L289
	} else {
		goto L290
	}
L285:
	;
	v3992 = base.F64_mul(v3989, float64(10))
	goto L220
L286:
	;
	v3989 = v3967
	goto L285
L287:
	;
	v3893 = v3891 + int32(614242)
	v3897 = base.F64_convert_i32_s(v3889 + int32(base.Ui32(v3893)>>(uint(int32(20))%32)))
	v3899 = base.F64_mul(v3897, float64(0.30102999566361177))
	v3912 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v3893&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v3888&int64(4294967295)), float64(-1))
	v3915 = base.F64_mul(v3912, base.F64_mul(v3912, float64(0.5)))
	v3920 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v3912, v3915)) & int64(-4294967296))
	v3921 = float64(0.4342944818781689)
	v3922 = base.F64_mul(v3920, v3921)
	v3923 = base.F64_add(v3899, v3922)
	v3928 = base.F64_div(v3912, base.F64_add(v3912, float64(2)))
	v3929 = base.F64_mul(v3928, v3928)
	v3930 = base.F64_mul(v3929, v3929)
	v3955 = base.F64_add(base.F64_mul(v3928, base.F64_add(v3915, base.F64_add(base.F64_mul(v3930, base.F64_add(base.F64_mul(v3930, base.F64_add(base.F64_mul(v3930, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v3929, base.F64_add(base.F64_mul(v3930, base.F64_add(base.F64_mul(v3930, base.F64_add(base.F64_mul(v3930, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v3912, v3920), v3915))
	v3967 = base.F64_add(v3923, base.F64_add(base.F64_add(v3922, base.F64_sub(v3899, v3923)), base.F64_add(base.F64_mul(v3955, v3921), base.F64_add(base.F64_mul(v3897, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v3955, v3920), float64(2.5082946711645275e-11))))))
	goto L286
L288:
	;
	v3883 = base.I64_reinterpret_f64(base.F64_mul(v3846, float64(1.8014398509481984e+16)))
	v3888 = v3883
	v3889 = int32(-1077)
	v3891 = base.I32_wrap_i64(int64(base.Ui64(v3883) >> (uint(int64(32)) % 64)))
	goto L287
L289:
	;
	if base.Ui64(int64(9218868437227405311)) < base.Ui64(v3857) {
		v3967 = v3846
		goto L286
	} else {
		goto L294
	}
L290:
	;
	if base.F64_ne(v3846, float64(0)) != 0 {
		goto L291
	} else {
		goto L292
	}
L291:
	;
	if int64(-1) < v3857 {
		goto L288
	} else {
		goto L293
	}
L292:
	;
	v3989 = base.F64_div(float64(-1), base.F64_mul(v3846, v3846))
	goto L285
L293:
	;
	v3989 = base.F64_div(base.F64_sub(v3846, v3846), float64(0))
	goto L285
L294:
	;
	v3872 = int32(-1023)
	v3874 = int64(base.Ui64(v3857) >> (uint(int64(32)) % 64))
	if v3874 == int64(1072693248) {
		goto L295
	} else {
		goto L296
	}
L295:
	;
	if base.I32_wrap_i64(v3857) != 0 {
		v3888 = v3857
		v3889 = v3872
		v3891 = int32(1072693248)
		goto L287
	} else {
		goto L297
	}
L296:
	;
	v3888 = v3857
	v3889 = v3872
	v3891 = base.I32_wrap_i64(v3874)
	goto L287
L297:
	;
	v3989 = float64(0)
	goto L285
L298:
	;
	if int32(0) < v4060 {
		v326 = v4060
		v329 = v4061
		v335 = v4062
		v341 = v647
		v371 = v4063
		v372 = v4064
		v373 = v4065
		goto L46
	} else {
		goto L328
	}
L299:
	;
	if v395 != 0 {
		v4087 = v647
		goto L40
	} else {
		goto L305
	}
L300:
	;
	if base.Ui64(v3226) < base.Ui64(int64(1069547521)) {
		goto L299
	} else {
		goto L301
	}
L301:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[8]))) = int32(base.Ui32(v3993) >> (uint(int32(1)) % 32))
	if v395 == int32(0) {
		v4060 = v326
		v4061 = v329
		v4062 = v335
		v4063 = v371
		v4064 = v372
		v4065 = v373
		goto L298
	} else {
		goto L302
	}
L302:
	;
	v4003 = *(*int32)(unsafe.Add(mBase, uint32(v64)+904))
	v4004 = *(*int32)(unsafe.Add(mBase, uint32(v4003)+4))
	v4005 = *(*int32)(unsafe.Add(mBase, uint32(v4004)+88))
	if v4005 == int32(0) {
		goto L303
	} else {
		goto L304
	}
L303:
	;
	v4014 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v4003)+uint32(_c_F_VP8EncTokenLoop[15]))) = v4014
	*(*int64)(unsafe.Add(mBase, uint32(v4003)+uint32(_c_F_VP8EncTokenLoop[16]))) = v4014
	v4020 = base.Simd_g_const(&F_VP8EncTokenLoop__k1)
	base.Simd_g_v128_store(m, v4003+int32(_a_F_VP8EncTokenLoop_21), int32(0), v4020)
	v4060 = v326
	v4061 = v329
	v4062 = v335
	v4063 = v371
	v4064 = v372
	v4065 = v373
	goto L298
L304:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4003)+uint32(_c_F_VP8EncTokenLoop[17]))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v4003)+uint32(_c_F_VP8EncTokenLoop[18]))) = int32(0)
	goto L303
L305:
	;
	if v66 == int32(0) {
		v4054 = v329
		v4055 = v335
		v4056 = v371
		v4057 = v372
		v4058 = v373
		goto L306
	} else {
		goto L307
	}
L306:
	;
	v4060 = v384
	v4061 = v4054
	v4062 = v4055
	v4063 = v4056
	v4064 = v4057
	v4065 = v4058
	goto L298
L307:
	;
	if v372 == int32(0) {
		goto L309
	} else {
		goto L310
	}
L308:
	;
	v4042 = float32(30)
	if base.F32_gt(v4039, v4042) != 0 {
		goto L316
	} else {
		goto L317
	}
L309:
	;
	if base.F64_ne(v3992, v371) != 0 {
		goto L314
	} else {
		goto L315
	}
L310:
	;
	if base.F64_gt(v3992, v101) != 0 {
		goto L311
	} else {
		goto L312
	}
L311:
	;
	v4029 = base.F32_neg(v335)
	goto L313
L312:
	;
	v4029 = v335
	goto L313
L313:
	;
	v4039 = v4029
	goto L308
L314:
	;
	v4039 = base.F32_demote_f64(base.F64_mul(base.F64_div(base.F64_sub(v101, v3992), base.F64_sub(v371, v3992)), base.F64_promote_f32(base.F32_sub(v373, v329))))
	goto L308
L315:
	;
	v4039 = float32(0)
	goto L308
L316:
	;
	v4045 = v4042
	goto L318
L317:
	;
	v4045 = v4039
	goto L318
L318:
	;
	if base.F32_lt(v4039, float32(-30)) != 0 {
		goto L319
	} else {
		goto L320
	}
L319:
	;
	v4048 = float32(-30)
	goto L321
L320:
	;
	v4048 = v4045
	goto L321
L321:
	;
	v4049 = base.F32_add(v329, v4048)
	if base.F32_gt(v4049, v72) != 0 {
		goto L322
	} else {
		goto L323
	}
L322:
	;
	v4051 = v72
	goto L324
L323:
	;
	v4051 = v4049
	goto L324
L324:
	;
	if base.F32_lt(v4049, v70) != 0 {
		goto L325
	} else {
		goto L326
	}
L325:
	;
	v4053 = v70
	goto L327
L326:
	;
	v4053 = v4051
	goto L327
L327:
	;
	v4054 = v4053
	v4055 = v4048
	v4056 = v3992
	v4057 = int32(0)
	v4058 = v329
	goto L306
L328:
	;
	goto L47
L329:
	;
	v4634 = l0 + int32(344)
	v4645 = *(*int32)(unsafe.Add(mBase, uint32(v4634)))
	if v4645 == int32(0) {
		goto L376
	} else {
		goto L377
	}
L330:
	;
	v4129 = int32(0)
	v4166 = l0 + int32(_a_F_VP8EncTokenLoop_2)
	v4171 = l0 + int32(_a_F_VP8EncTokenLoop_5)
	v4172 = l0 + int32(3442)
	v4173 = l0 + int32(_a_F_VP8EncTokenLoop_6)
	v4174 = l0 + int32(3431)
	v4175 = v4129
	v4177 = v4166
	v4178 = v4129
	v4180 = v4129
	goto L332
L331:
	;
	goto L329
L332:
	;
	v4207 = v4178
	v4210 = v4171
	v4211 = v4172
	v4212 = v4173
	v4213 = v4174
	v4214 = v4175
	v4215 = v4177
	v4216 = int32(0)
	goto L334
L333:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = v4592
	goto L331
L334:
	;
	v4236 = v4207
	v4246 = int32(0)
	v4247 = v4214
	goto L336
L335:
	;
	v4616 = int32(1056)
	v4618 = int32(264)
	v4629 = v4180 + int32(1)
	if v4629 != int32(4) {
		v4171 = v4171 + v4616
		v4172 = v4172 + v4618
		v4173 = v4173 + v4616
		v4174 = v4174 + v4618
		v4175 = v4175 + v4618
		v4177 = v4177 + v4616
		v4178 = v4592
		v4180 = v4629
		goto L332
	} else {
		goto L373
	}
L336:
	;
	v4257 = *(*int32)(unsafe.Add(mBase, uint32(v4215+v4246)))
	v4259 = int32(base.Ui32(v4257) >> (uint(int32(16)) % 32))
	v4260 = m.G125
	v4262 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4260+v4247))))
	v4263 = m.G126
	v4265 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4263+v4247))))
	v4267 = v4257 & int32(_a_F_VP8EncTokenLoop_7)
	if v4267 != 0 {
		goto L339
	} else {
		goto L340
	}
L337:
	;
	v4361 = v4344
	v4371 = v4212
	v4377 = int32(0)
	goto L348
L338:
	;
	v4275 = m.G79
	v4278 = v4259 - v4267
	v4279 = int32(1)
	v4282 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4275+v4262<<(uint(v4279)%32)))))
	v4284 = int32(255)
	v4289 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4275+(v4262^v4284)<<(uint(v4279)%32)))))
	v4295 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4275+v4265<<(uint(v4279)%32)))))
	v4298 = v4274 & v4284
	v4304 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4275+(v4298^v4284)<<(uint(v4279)%32)))))
	v4309 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4275+v4298<<(uint(v4279)%32)))))
	v4317 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4275+(v4265^v4284)<<(uint(v4279)%32)))))
	goto L341
L339:
	;
	v4269 = int32(255)
	v4272 = base.I32_div_u_s(v4267*v4269, v4259)
	v4274 = v4269 - v4272
	goto L338
L340:
	;
	v4274 = int32(255)
	goto L338
L341:
	;
	goto L343
L343:
	;
	if v4278*v4282+v4267*v4289+v4295 <= v4267*v4304+v4278*v4309+v4317+int32(2048) {
		goto L345
	} else {
		goto L346
	}
L344:
	;
	v4349 = v4246 + int32(4)
	if v4349 != int32(44) {
		v4236 = v4344
		v4246 = v4349
		v4247 = v4247 + int32(1)
		goto L336
	} else {
		goto L347
	}
L345:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4166+v4247+int32(-1056)))) = uint8(v4262)
	v4344 = v4236
	goto L344
L346:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4166+v4247+int32(-1056)))) = uint8(v4274)
	v4344 = v4236 | base.B2i32(v4274 != v4262)
	goto L344
L347:
	;
	goto L337
L348:
	;
	v4381 = *(*int32)(unsafe.Add(mBase, uint32(v4371)))
	v4383 = int32(base.Ui32(v4381) >> (uint(int32(16)) % 32))
	v4384 = m.G125
	v4385 = v4214 + v4377
	v4387 = int32(11)
	v4389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4384+v4385+v4387))))
	v4390 = m.G126
	v4394 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4390+v4385+v4387))))
	v4396 = v4381 & int32(_a_F_VP8EncTokenLoop_7)
	if v4396 != 0 {
		goto L351
	} else {
		goto L352
	}
L349:
	;
	v4485 = v4468
	v4495 = v4210
	v4501 = int32(0)
	goto L360
L350:
	;
	v4404 = m.G79
	v4407 = v4383 - v4396
	v4408 = int32(1)
	v4411 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4404+v4389<<(uint(v4408)%32)))))
	v4413 = int32(255)
	v4418 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4404+(v4389^v4413)<<(uint(v4408)%32)))))
	v4424 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4404+v4394<<(uint(v4408)%32)))))
	v4427 = v4403 & v4413
	v4433 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4404+(v4427^v4413)<<(uint(v4408)%32)))))
	v4438 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4404+v4427<<(uint(v4408)%32)))))
	v4446 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4404+(v4394^v4413)<<(uint(v4408)%32)))))
	goto L353
L351:
	;
	v4398 = int32(255)
	v4401 = base.I32_div_u_s(v4396*v4398, v4383)
	v4403 = v4398 - v4401
	goto L350
L352:
	;
	v4403 = int32(255)
	goto L350
L353:
	;
	goto L355
L355:
	;
	if v4396*v4433+v4407*v4438+v4446+int32(2048) < v4407*v4411+v4396*v4418+v4424 {
		goto L357
	} else {
		goto L358
	}
L356:
	;
	v4473 = v4377 + int32(1)
	if v4473 != int32(11) {
		v4361 = v4468
		v4371 = v4371 + int32(4)
		v4377 = v4473
		goto L348
	} else {
		goto L359
	}
L357:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4213+v4377))) = uint8(v4403)
	v4468 = v4361 | base.B2i32(v4403 != v4389)
	goto L356
L358:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4213+v4377))) = uint8(v4389)
	v4468 = v4361
	goto L356
L359:
	;
	goto L349
L360:
	;
	v4505 = *(*int32)(unsafe.Add(mBase, uint32(v4495)))
	v4507 = int32(base.Ui32(v4505) >> (uint(int32(16)) % 32))
	v4508 = m.G125
	v4509 = v4214 + v4501
	v4511 = int32(22)
	v4513 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4508+v4509+v4511))))
	v4514 = m.G126
	v4518 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4514+v4509+v4511))))
	v4520 = v4505 & int32(_a_F_VP8EncTokenLoop_7)
	if v4520 != 0 {
		goto L363
	} else {
		goto L364
	}
L361:
	;
	v4600 = int32(132)
	v4602 = int32(33)
	v4613 = v4216 + int32(1)
	if v4613 != int32(8) {
		v4207 = v4592
		v4210 = v4210 + v4600
		v4211 = v4211 + v4602
		v4212 = v4212 + v4600
		v4213 = v4213 + v4602
		v4214 = v4214 + v4602
		v4215 = v4215 + v4600
		v4216 = v4613
		goto L334
	} else {
		goto L372
	}
L362:
	;
	v4528 = m.G79
	v4531 = v4507 - v4520
	v4532 = int32(1)
	v4535 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4528+v4513<<(uint(v4532)%32)))))
	v4537 = int32(255)
	v4542 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4528+(v4513^v4537)<<(uint(v4532)%32)))))
	v4548 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4528+v4518<<(uint(v4532)%32)))))
	v4551 = v4527 & v4537
	v4557 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4528+(v4551^v4537)<<(uint(v4532)%32)))))
	v4562 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4528+v4551<<(uint(v4532)%32)))))
	v4570 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4528+(v4518^v4537)<<(uint(v4532)%32)))))
	goto L365
L363:
	;
	v4522 = int32(255)
	v4525 = base.I32_div_u_s(v4520*v4522, v4507)
	v4527 = v4522 - v4525
	goto L362
L364:
	;
	v4527 = int32(255)
	goto L362
L365:
	;
	goto L367
L367:
	;
	if v4520*v4557+v4531*v4562+v4570+int32(2048) < v4531*v4535+v4520*v4542+v4548 {
		goto L369
	} else {
		goto L370
	}
L368:
	;
	v4597 = v4501 + int32(1)
	if v4597 != int32(11) {
		v4485 = v4592
		v4495 = v4495 + int32(4)
		v4501 = v4597
		goto L360
	} else {
		goto L371
	}
L369:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4211+v4501))) = uint8(v4527)
	v4592 = v4485 | base.B2i32(v4527 != v4513)
	goto L368
L370:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4211+v4501))) = uint8(v4513)
	v4592 = v4485
	goto L368
L371:
	;
	goto L361
L372:
	;
	goto L335
L373:
	;
	goto L333
L374:
	;
	v4796 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v4797 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	v4801 = F_WebPReportProgress(m, v4796, v4797+v4087, l0+int32(368))
	mBase = m.M
	v4807 = base.B2i32(v4801 != int32(0))
	goto L39
L375:
	;
	goto L374
L376:
	;
	goto L407
L377:
	;
	goto L379
L379:
	;
	v4654 = v4645
	goto L380
L380:
	;
	v4662 = *(*int32)(unsafe.Add(mBase, uint32(v4654)))
	if v4662 != 0 {
		v4664 = int32(0)
		goto L382
	} else {
		goto L383
	}
L382:
	;
	v4665 = *(*int32)(unsafe.Add(mBase, uint32(v4634)+16))
	if v4665 <= v4664 {
		goto L384
	} else {
		goto L385
	}
L383:
	;
	v4663 = *(*int32)(unsafe.Add(mBase, uint32(v4634)+12))
	v4664 = v4663
	goto L382
L384:
	;
	F_WebPSafeFree(m, v4654)
	mBase = m.M
	if v4662 != 0 {
		v4654 = v4662
		goto L380
	} else {
		goto L392
	}
L385:
	;
	v4679 = v4665
	v4680 = v4654 + v4665<<(uint(int32(1))%32) + int32(2)
	goto L386
L386:
	;
	v4683 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4680))))
	if v4683&int32(_a_F_VP8EncTokenLoop_19) == int32(0) {
		goto L389
	} else {
		goto L390
	}
L387:
	;
	goto L384
L388:
	;
	v4697 = F_VP8PutBit(m, v112, int32(base.Ui32(v4683)>>(uint(int32(15))%32)), v4696)
	mBase = m.M
	v4701 = v4679 + int32(-1)
	if v4664 < v4701 {
		v4679 = v4701
		v4680 = v4680 + int32(-2)
		goto L386
	} else {
		goto L391
	}
L389:
	;
	v4695 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3420)+v4683&int32(_a_F_VP8EncTokenLoop_20)))))
	v4696 = v4695
	goto L388
L390:
	;
	v4696 = v4683 & int32(255)
	goto L388
L391:
	;
	goto L387
L392:
	;
	goto L376
L406:
	;
	goto L375
L407:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v4634))) = int32(0)
	goto L406
L409:
	;
	v5005 = v5001
	goto L13
L410:
	;
	F_VP8EncFreeBitWriters(m, v4871)
	mBase = m.M
	v4992 = *(*int32)(unsafe.Add(mBase, uint32(v4871)+4))
	v4994 = F_WebPEncodingSetError(m, v4992, int32(1))
	mBase = m.M
	v5001 = v4994
	goto L409
L411:
	;
	v4874 = *(*int32)(unsafe.Add(mBase, uint32(v4871)+52))
	if v4874 < int32(1) {
		v4903 = v4807
		goto L412
	} else {
		goto L413
	}
L412:
	;
	v4908 = *(*int32)(unsafe.Add(mBase, uint32(v4871)+4))
	v4909 = *(*int32)(unsafe.Add(mBase, uint32(v4908)+88))
	if v4909 == int32(0) {
		goto L421
	} else {
		goto L422
	}
L413:
	;
	v4885 = v4807 & int32(1)
	v4886 = v4871 + int32(88)
	v4887 = int32(0)
	goto L414
L414:
	;
	v4888 = F_VP8BitWriterFinish(m, v4886)
	mBase = m.M
	v4892 = *(*int32)(unsafe.Add(mBase, uint32(v4886+int32(28))))
	if v4892 != 0 {
		goto L416
	} else {
		goto L417
	}
L415:
	;
	if v4893 == int32(0) {
		goto L410
	} else {
		goto L420
	}
L416:
	;
	v4893 = int32(0)
	goto L418
L417:
	;
	v4893 = v4885
	goto L418
L418:
	;
	v4897 = v4887 + int32(1)
	v4898 = *(*int32)(unsafe.Add(mBase, uint32(v4871)+52))
	if v4897 < v4898 {
		v4885 = v4893
		v4886 = v4886 + int32(32)
		v4887 = v4897
		goto L414
	} else {
		goto L419
	}
L419:
	;
	goto L415
L420:
	;
	v4903 = v4885
	goto L412
L421:
	;
	F_VP8AdjustFilterStrength(m, v4866)
	mBase = m.M
	v5001 = v4903
	goto L409
L422:
	;
	v4912 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+168))
	v4913 = int64(7)
	v4915 = int64(3)
	v4916 = int64(base.Ui64(v4912+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[19]))) = uint32(v4916)
	v4918 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+192))
	v4922 = int64(base.Ui64(v4918+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[20]))) = uint32(v4922)
	v4924 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+216))
	v4928 = int64(base.Ui64(v4924+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[21]))) = uint32(v4928)
	v4930 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+240))
	v4934 = int64(base.Ui64(v4930+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[22]))) = uint32(v4934)
	v4936 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+176))
	v4940 = int64(base.Ui64(v4936+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[23]))) = uint32(v4940)
	v4942 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+200))
	v4946 = int64(base.Ui64(v4942+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[24]))) = uint32(v4946)
	v4948 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+224))
	v4952 = int64(base.Ui64(v4948+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[25]))) = uint32(v4952)
	v4954 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+248))
	v4958 = int64(base.Ui64(v4954+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[26]))) = uint32(v4958)
	v4960 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+184))
	v4964 = int64(base.Ui64(v4960+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[27]))) = uint32(v4964)
	v4966 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+208))
	v4970 = int64(base.Ui64(v4966+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[28]))) = uint32(v4970)
	v4972 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+232))
	v4976 = int64(base.Ui64(v4972+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[29]))) = uint32(v4976)
	v4978 = *(*int64)(unsafe.Add(mBase, uint32(v4866)+256))
	v4982 = int64(base.Ui64(v4978+v4913) >> (uint(v4915) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4871)+uint32(_c_F_VP8EncTokenLoop[30]))) = uint32(v4982)
	goto L421
}

var F_VP8EncTokenLoop__k0 = [2]uint64{0xfe, 0xfffffff800000000}
var F_VP8EncTokenLoop__k1 = [2]uint64{0x0, 0x0}

func F_VP8EncWrite(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v19 base.V128
	_ = v19
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v32 int32
	_ = v32
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v42 int32
	_ = v42
	var v45 base.V128
	_ = v45
	var v54 int32
	_ = v54
	var v57 int32
	_ = v57
	var v58 int32
	_ = v58
	var v62 int32
	_ = v62
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	var v69 int32
	_ = v69
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v83 int32
	_ = v83
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v112 int32
	_ = v112
	var v116 int32
	_ = v116
	var v125 int32
	_ = v125
	var v127 int32
	_ = v127
	var v140 int32
	_ = v140
	var v141 int32
	_ = v141
	var v144 int32
	_ = v144
	var v148 int32
	_ = v148
	var v152 int32
	_ = v152
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v172 int32
	_ = v172
	var v176 int32
	_ = v176
	var v180 int32
	_ = v180
	var v181 int32
	_ = v181
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v192 int32
	_ = v192
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v205 int32
	_ = v205
	var v208 int32
	_ = v208
	var v210 int32
	_ = v210
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v226 int32
	_ = v226
	var v230 int32
	_ = v230
	var v239 int32
	_ = v239
	var v244 int32
	_ = v244
	var v246 int32
	_ = v246
	var v249 int32
	_ = v249
	var v254 int32
	_ = v254
	var v255 int32
	_ = v255
	var v258 int32
	_ = v258
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v277 int32
	_ = v277
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v287 int32
	_ = v287
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v296 int32
	_ = v296
	var v300 int32
	_ = v300
	var v304 int32
	_ = v304
	var v310 int32
	_ = v310
	var v315 int32
	_ = v315
	var v317 int32
	_ = v317
	var v320 int32
	_ = v320
	var v322 int32
	_ = v322
	var v326 int32
	_ = v326
	var v330 int32
	_ = v330
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v338 int32
	_ = v338
	var v342 int32
	_ = v342
	var v352 int32
	_ = v352
	var v360 int32
	_ = v360
	var v364 int32
	_ = v364
	var v366 int32
	_ = v366
	var v370 int32
	_ = v370
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
	var v388 int32
	_ = v388
	var v392 int32
	_ = v392
	var v408 int32
	_ = v408
	var v412 int32
	_ = v412
	var v414 int32
	_ = v414
	var v418 int32
	_ = v418
	var v420 int32
	_ = v420
	var v424 int32
	_ = v424
	var v428 int32
	_ = v428
	var v429 int32
	_ = v429
	var v432 int32
	_ = v432
	var v436 int32
	_ = v436
	var v440 int32
	_ = v440
	var v446 int32
	_ = v446
	var v456 int32
	_ = v456
	var v461 int32
	_ = v461
	var v463 int32
	_ = v463
	var v466 int32
	_ = v466
	var v468 int32
	_ = v468
	var v472 int32
	_ = v472
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v480 int32
	_ = v480
	var v484 int32
	_ = v484
	var v488 int32
	_ = v488
	var v498 int32
	_ = v498
	var v506 int32
	_ = v506
	var v510 int32
	_ = v510
	var v512 int32
	_ = v512
	var v516 int32
	_ = v516
	var v518 int32
	_ = v518
	var v522 int32
	_ = v522
	var v526 int32
	_ = v526
	var v527 int32
	_ = v527
	var v530 int32
	_ = v530
	var v534 int32
	_ = v534
	var v538 int32
	_ = v538
	var v554 int32
	_ = v554
	var v558 int32
	_ = v558
	var v560 int32
	_ = v560
	var v564 int32
	_ = v564
	var v566 int32
	_ = v566
	var v570 int32
	_ = v570
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v578 int32
	_ = v578
	var v582 int32
	_ = v582
	var v586 int32
	_ = v586
	var v592 int32
	_ = v592
	var v602 int32
	_ = v602
	var v607 int32
	_ = v607
	var v609 int32
	_ = v609
	var v612 int32
	_ = v612
	var v614 int32
	_ = v614
	var v618 int32
	_ = v618
	var v622 int32
	_ = v622
	var v623 int32
	_ = v623
	var v626 int32
	_ = v626
	var v630 int32
	_ = v630
	var v634 int32
	_ = v634
	var v644 int32
	_ = v644
	var v652 int32
	_ = v652
	var v656 int32
	_ = v656
	var v658 int32
	_ = v658
	var v662 int32
	_ = v662
	var v664 int32
	_ = v664
	var v668 int32
	_ = v668
	var v672 int32
	_ = v672
	var v673 int32
	_ = v673
	var v676 int32
	_ = v676
	var v680 int32
	_ = v680
	var v684 int32
	_ = v684
	var v700 int32
	_ = v700
	var v704 int32
	_ = v704
	var v706 int32
	_ = v706
	var v710 int32
	_ = v710
	var v712 int32
	_ = v712
	var v716 int32
	_ = v716
	var v720 int32
	_ = v720
	var v721 int32
	_ = v721
	var v724 int32
	_ = v724
	var v728 int32
	_ = v728
	var v732 int32
	_ = v732
	var v738 int32
	_ = v738
	var v748 int32
	_ = v748
	var v753 int32
	_ = v753
	var v755 int32
	_ = v755
	var v758 int32
	_ = v758
	var v760 int32
	_ = v760
	var v764 int32
	_ = v764
	var v768 int32
	_ = v768
	var v769 int32
	_ = v769
	var v772 int32
	_ = v772
	var v776 int32
	_ = v776
	var v780 int32
	_ = v780
	var v790 int32
	_ = v790
	var v798 int32
	_ = v798
	var v802 int32
	_ = v802
	var v804 int32
	_ = v804
	var v808 int32
	_ = v808
	var v810 int32
	_ = v810
	var v814 int32
	_ = v814
	var v818 int32
	_ = v818
	var v819 int32
	_ = v819
	var v822 int32
	_ = v822
	var v826 int32
	_ = v826
	var v830 int32
	_ = v830
	var v846 int32
	_ = v846
	var v850 int32
	_ = v850
	var v852 int32
	_ = v852
	var v856 int32
	_ = v856
	var v858 int32
	_ = v858
	var v862 int32
	_ = v862
	var v866 int32
	_ = v866
	var v867 int32
	_ = v867
	var v870 int32
	_ = v870
	var v874 int32
	_ = v874
	var v878 int32
	_ = v878
	var v884 int32
	_ = v884
	var v894 int32
	_ = v894
	var v899 int32
	_ = v899
	var v901 int32
	_ = v901
	var v904 int32
	_ = v904
	var v906 int32
	_ = v906
	var v910 int32
	_ = v910
	var v914 int32
	_ = v914
	var v915 int32
	_ = v915
	var v918 int32
	_ = v918
	var v922 int32
	_ = v922
	var v926 int32
	_ = v926
	var v936 int32
	_ = v936
	var v944 int32
	_ = v944
	var v948 int32
	_ = v948
	var v950 int32
	_ = v950
	var v954 int32
	_ = v954
	var v956 int32
	_ = v956
	var v960 int32
	_ = v960
	var v964 int32
	_ = v964
	var v965 int32
	_ = v965
	var v968 int32
	_ = v968
	var v972 int32
	_ = v972
	var v976 int32
	_ = v976
	var v992 int32
	_ = v992
	var v996 int32
	_ = v996
	var v998 int32
	_ = v998
	var v1002 int32
	_ = v1002
	var v1004 int32
	_ = v1004
	var v1008 int32
	_ = v1008
	var v1012 int32
	_ = v1012
	var v1013 int32
	_ = v1013
	var v1016 int32
	_ = v1016
	var v1020 int32
	_ = v1020
	var v1024 int32
	_ = v1024
	var v1030 int32
	_ = v1030
	var v1040 int32
	_ = v1040
	var v1045 int32
	_ = v1045
	var v1047 int32
	_ = v1047
	var v1050 int32
	_ = v1050
	var v1052 int32
	_ = v1052
	var v1056 int32
	_ = v1056
	var v1060 int32
	_ = v1060
	var v1061 int32
	_ = v1061
	var v1064 int32
	_ = v1064
	var v1068 int32
	_ = v1068
	var v1072 int32
	_ = v1072
	var v1082 int32
	_ = v1082
	var v1090 int32
	_ = v1090
	var v1094 int32
	_ = v1094
	var v1096 int32
	_ = v1096
	var v1100 int32
	_ = v1100
	var v1102 int32
	_ = v1102
	var v1106 int32
	_ = v1106
	var v1110 int32
	_ = v1110
	var v1111 int32
	_ = v1111
	var v1114 int32
	_ = v1114
	var v1118 int32
	_ = v1118
	var v1122 int32
	_ = v1122
	var v1138 int32
	_ = v1138
	var v1142 int32
	_ = v1142
	var v1144 int32
	_ = v1144
	var v1148 int32
	_ = v1148
	var v1150 int32
	_ = v1150
	var v1154 int32
	_ = v1154
	var v1158 int32
	_ = v1158
	var v1159 int32
	_ = v1159
	var v1162 int32
	_ = v1162
	var v1166 int32
	_ = v1166
	var v1170 int32
	_ = v1170
	var v1176 int32
	_ = v1176
	var v1186 int32
	_ = v1186
	var v1191 int32
	_ = v1191
	var v1193 int32
	_ = v1193
	var v1196 int32
	_ = v1196
	var v1198 int32
	_ = v1198
	var v1202 int32
	_ = v1202
	var v1206 int32
	_ = v1206
	var v1207 int32
	_ = v1207
	var v1210 int32
	_ = v1210
	var v1214 int32
	_ = v1214
	var v1218 int32
	_ = v1218
	var v1228 int32
	_ = v1228
	var v1236 int32
	_ = v1236
	var v1240 int32
	_ = v1240
	var v1242 int32
	_ = v1242
	var v1246 int32
	_ = v1246
	var v1248 int32
	_ = v1248
	var v1252 int32
	_ = v1252
	var v1256 int32
	_ = v1256
	var v1257 int32
	_ = v1257
	var v1260 int32
	_ = v1260
	var v1264 int32
	_ = v1264
	var v1268 int32
	_ = v1268
	var v1284 int32
	_ = v1284
	var v1288 int32
	_ = v1288
	var v1290 int32
	_ = v1290
	var v1294 int32
	_ = v1294
	var v1296 int32
	_ = v1296
	var v1300 int32
	_ = v1300
	var v1304 int32
	_ = v1304
	var v1305 int32
	_ = v1305
	var v1308 int32
	_ = v1308
	var v1312 int32
	_ = v1312
	var v1316 int32
	_ = v1316
	var v1322 int32
	_ = v1322
	var v1332 int32
	_ = v1332
	var v1337 int32
	_ = v1337
	var v1339 int32
	_ = v1339
	var v1342 int32
	_ = v1342
	var v1344 int32
	_ = v1344
	var v1348 int32
	_ = v1348
	var v1352 int32
	_ = v1352
	var v1353 int32
	_ = v1353
	var v1356 int32
	_ = v1356
	var v1360 int32
	_ = v1360
	var v1364 int32
	_ = v1364
	var v1374 int32
	_ = v1374
	var v1382 int32
	_ = v1382
	var v1386 int32
	_ = v1386
	var v1388 int32
	_ = v1388
	var v1392 int32
	_ = v1392
	var v1394 int32
	_ = v1394
	var v1398 int32
	_ = v1398
	var v1402 int32
	_ = v1402
	var v1403 int32
	_ = v1403
	var v1406 int32
	_ = v1406
	var v1410 int32
	_ = v1410
	var v1414 int32
	_ = v1414
	var v1430 int32
	_ = v1430
	var v1434 int32
	_ = v1434
	var v1436 int32
	_ = v1436
	var v1440 int32
	_ = v1440
	var v1442 int32
	_ = v1442
	var v1446 int32
	_ = v1446
	var v1450 int32
	_ = v1450
	var v1451 int32
	_ = v1451
	var v1454 int32
	_ = v1454
	var v1458 int32
	_ = v1458
	var v1462 int32
	_ = v1462
	var v1468 int32
	_ = v1468
	var v1478 int32
	_ = v1478
	var v1481 int32
	_ = v1481
	var v1483 int32
	_ = v1483
	var v1486 int32
	_ = v1486
	var v1488 int32
	_ = v1488
	var v1491 int32
	_ = v1491
	var v1493 int32
	_ = v1493
	var v1497 int32
	_ = v1497
	var v1501 int32
	_ = v1501
	var v1502 int32
	_ = v1502
	var v1505 int32
	_ = v1505
	var v1509 int32
	_ = v1509
	var v1513 int32
	_ = v1513
	var v1521 int32
	_ = v1521
	var v1533 int32
	_ = v1533
	var v1536 int32
	_ = v1536
	var v1538 int32
	_ = v1538
	var v1542 int32
	_ = v1542
	var v1544 int32
	_ = v1544
	var v1548 int32
	_ = v1548
	var v1552 int32
	_ = v1552
	var v1553 int32
	_ = v1553
	var v1556 int32
	_ = v1556
	var v1560 int32
	_ = v1560
	var v1564 int32
	_ = v1564
	var v1570 int32
	_ = v1570
	var v1574 int32
	_ = v1574
	var v1576 int32
	_ = v1576
	var v1579 int32
	_ = v1579
	var v1581 int32
	_ = v1581
	var v1584 int32
	_ = v1584
	var v1586 int32
	_ = v1586
	var v1590 int32
	_ = v1590
	var v1594 int32
	_ = v1594
	var v1595 int32
	_ = v1595
	var v1598 int32
	_ = v1598
	var v1602 int32
	_ = v1602
	var v1606 int32
	_ = v1606
	var v1614 int32
	_ = v1614
	var v1626 int32
	_ = v1626
	var v1629 int32
	_ = v1629
	var v1631 int32
	_ = v1631
	var v1635 int32
	_ = v1635
	var v1637 int32
	_ = v1637
	var v1641 int32
	_ = v1641
	var v1645 int32
	_ = v1645
	var v1646 int32
	_ = v1646
	var v1649 int32
	_ = v1649
	var v1653 int32
	_ = v1653
	var v1657 int32
	_ = v1657
	var v1663 int32
	_ = v1663
	var v1667 int32
	_ = v1667
	var v1669 int32
	_ = v1669
	var v1672 int32
	_ = v1672
	var v1674 int32
	_ = v1674
	var v1677 int32
	_ = v1677
	var v1679 int32
	_ = v1679
	var v1683 int32
	_ = v1683
	var v1687 int32
	_ = v1687
	var v1688 int32
	_ = v1688
	var v1691 int32
	_ = v1691
	var v1695 int32
	_ = v1695
	var v1699 int32
	_ = v1699
	var v1707 int32
	_ = v1707
	var v1719 int32
	_ = v1719
	var v1722 int32
	_ = v1722
	var v1724 int32
	_ = v1724
	var v1728 int32
	_ = v1728
	var v1730 int32
	_ = v1730
	var v1734 int32
	_ = v1734
	var v1738 int32
	_ = v1738
	var v1739 int32
	_ = v1739
	var v1742 int32
	_ = v1742
	var v1746 int32
	_ = v1746
	var v1750 int32
	_ = v1750
	var v1756 int32
	_ = v1756
	var v1760 int32
	_ = v1760
	var v1761 int32
	_ = v1761
	var v1764 int32
	_ = v1764
	var v1766 int32
	_ = v1766
	var v1769 int32
	_ = v1769
	var v1771 int32
	_ = v1771
	var v1775 int32
	_ = v1775
	var v1779 int32
	_ = v1779
	var v1780 int32
	_ = v1780
	var v1783 int32
	_ = v1783
	var v1787 int32
	_ = v1787
	var v1791 int32
	_ = v1791
	var v1797 int32
	_ = v1797
	var v1809 int32
	_ = v1809
	var v1812 int32
	_ = v1812
	var v1814 int32
	_ = v1814
	var v1818 int32
	_ = v1818
	var v1820 int32
	_ = v1820
	var v1824 int32
	_ = v1824
	var v1828 int32
	_ = v1828
	var v1829 int32
	_ = v1829
	var v1832 int32
	_ = v1832
	var v1836 int32
	_ = v1836
	var v1840 int32
	_ = v1840
	var v1846 int32
	_ = v1846
	var v1851 int32
	_ = v1851
	var v1863 int32
	_ = v1863
	var v1866 int32
	_ = v1866
	var v1868 int32
	_ = v1868
	var v1872 int32
	_ = v1872
	var v1874 int32
	_ = v1874
	var v1878 int32
	_ = v1878
	var v1882 int32
	_ = v1882
	var v1883 int32
	_ = v1883
	var v1886 int32
	_ = v1886
	var v1890 int32
	_ = v1890
	var v1894 int32
	_ = v1894
	var v1900 int32
	_ = v1900
	var v1904 int32
	_ = v1904
	var v1905 int32
	_ = v1905
	var v1908 int32
	_ = v1908
	var v1910 int32
	_ = v1910
	var v1913 int32
	_ = v1913
	var v1915 int32
	_ = v1915
	var v1919 int32
	_ = v1919
	var v1923 int32
	_ = v1923
	var v1924 int32
	_ = v1924
	var v1927 int32
	_ = v1927
	var v1931 int32
	_ = v1931
	var v1935 int32
	_ = v1935
	var v1943 int32
	_ = v1943
	var v1944 int32
	_ = v1944
	var v1945 int32
	_ = v1945
	var v1948 int32
	_ = v1948
	var v1950 int32
	_ = v1950
	var v1953 int32
	_ = v1953
	var v1955 int32
	_ = v1955
	var v1959 int32
	_ = v1959
	var v1963 int32
	_ = v1963
	var v1964 int32
	_ = v1964
	var v1967 int32
	_ = v1967
	var v1971 int32
	_ = v1971
	var v1975 int32
	_ = v1975
	var v1995 int32
	_ = v1995
	var v1998 int32
	_ = v1998
	var v2000 int32
	_ = v2000
	var v2004 int32
	_ = v2004
	var v2006 int32
	_ = v2006
	var v2010 int32
	_ = v2010
	var v2014 int32
	_ = v2014
	var v2015 int32
	_ = v2015
	var v2018 int32
	_ = v2018
	var v2022 int32
	_ = v2022
	var v2026 int32
	_ = v2026
	var v2032 int32
	_ = v2032
	var v2036 int32
	_ = v2036
	var v2041 int32
	_ = v2041
	var v2043 int32
	_ = v2043
	var v2046 int32
	_ = v2046
	var v2048 int32
	_ = v2048
	var v2052 int32
	_ = v2052
	var v2056 int32
	_ = v2056
	var v2057 int32
	_ = v2057
	var v2060 int32
	_ = v2060
	var v2064 int32
	_ = v2064
	var v2068 int32
	_ = v2068
	var v2078 int32
	_ = v2078
	var v2086 int32
	_ = v2086
	var v2090 int32
	_ = v2090
	var v2092 int32
	_ = v2092
	var v2096 int32
	_ = v2096
	var v2098 int32
	_ = v2098
	var v2102 int32
	_ = v2102
	var v2106 int32
	_ = v2106
	var v2107 int32
	_ = v2107
	var v2110 int32
	_ = v2110
	var v2114 int32
	_ = v2114
	var v2118 int32
	_ = v2118
	var v2134 int32
	_ = v2134
	var v2138 int32
	_ = v2138
	var v2140 int32
	_ = v2140
	var v2144 int32
	_ = v2144
	var v2146 int32
	_ = v2146
	var v2150 int32
	_ = v2150
	var v2154 int32
	_ = v2154
	var v2155 int32
	_ = v2155
	var v2158 int32
	_ = v2158
	var v2162 int32
	_ = v2162
	var v2166 int32
	_ = v2166
	var v2172 int32
	_ = v2172
	var v2194 int32
	_ = v2194
	var v2197 int32
	_ = v2197
	var v2199 int32
	_ = v2199
	var v2203 int32
	_ = v2203
	var v2205 int32
	_ = v2205
	var v2209 int32
	_ = v2209
	var v2213 int32
	_ = v2213
	var v2214 int32
	_ = v2214
	var v2217 int32
	_ = v2217
	var v2221 int32
	_ = v2221
	var v2225 int32
	_ = v2225
	var v2231 int32
	_ = v2231
	var v2235 int32
	_ = v2235
	var v2241 int32
	_ = v2241
	var v2253 int32
	_ = v2253
	var v2256 int32
	_ = v2256
	var v2258 int32
	_ = v2258
	var v2262 int32
	_ = v2262
	var v2264 int32
	_ = v2264
	var v2268 int32
	_ = v2268
	var v2272 int32
	_ = v2272
	var v2273 int32
	_ = v2273
	var v2276 int32
	_ = v2276
	var v2280 int32
	_ = v2280
	var v2284 int32
	_ = v2284
	var v2290 int32
	_ = v2290
	var v2294 int32
	_ = v2294
	var v2306 int32
	_ = v2306
	var v2309 int32
	_ = v2309
	var v2311 int32
	_ = v2311
	var v2315 int32
	_ = v2315
	var v2317 int32
	_ = v2317
	var v2321 int32
	_ = v2321
	var v2325 int32
	_ = v2325
	var v2326 int32
	_ = v2326
	var v2329 int32
	_ = v2329
	var v2333 int32
	_ = v2333
	var v2337 int32
	_ = v2337
	var v2343 int32
	_ = v2343
	var v2347 int32
	_ = v2347
	var v2352 int32
	_ = v2352
	var v2354 int32
	_ = v2354
	var v2357 int32
	_ = v2357
	var v2359 int32
	_ = v2359
	var v2363 int32
	_ = v2363
	var v2367 int32
	_ = v2367
	var v2368 int32
	_ = v2368
	var v2371 int32
	_ = v2371
	var v2375 int32
	_ = v2375
	var v2379 int32
	_ = v2379
	var v2389 int32
	_ = v2389
	var v2397 int32
	_ = v2397
	var v2401 int32
	_ = v2401
	var v2403 int32
	_ = v2403
	var v2407 int32
	_ = v2407
	var v2409 int32
	_ = v2409
	var v2413 int32
	_ = v2413
	var v2417 int32
	_ = v2417
	var v2418 int32
	_ = v2418
	var v2421 int32
	_ = v2421
	var v2425 int32
	_ = v2425
	var v2429 int32
	_ = v2429
	var v2445 int32
	_ = v2445
	var v2449 int32
	_ = v2449
	var v2451 int32
	_ = v2451
	var v2455 int32
	_ = v2455
	var v2457 int32
	_ = v2457
	var v2461 int32
	_ = v2461
	var v2465 int32
	_ = v2465
	var v2466 int32
	_ = v2466
	var v2469 int32
	_ = v2469
	var v2473 int32
	_ = v2473
	var v2477 int32
	_ = v2477
	var v2483 int32
	_ = v2483
	var v2493 int32
	_ = v2493
	var v2498 int32
	_ = v2498
	var v2500 int32
	_ = v2500
	var v2503 int32
	_ = v2503
	var v2505 int32
	_ = v2505
	var v2509 int32
	_ = v2509
	var v2513 int32
	_ = v2513
	var v2514 int32
	_ = v2514
	var v2517 int32
	_ = v2517
	var v2521 int32
	_ = v2521
	var v2525 int32
	_ = v2525
	var v2535 int32
	_ = v2535
	var v2543 int32
	_ = v2543
	var v2547 int32
	_ = v2547
	var v2549 int32
	_ = v2549
	var v2553 int32
	_ = v2553
	var v2555 int32
	_ = v2555
	var v2559 int32
	_ = v2559
	var v2563 int32
	_ = v2563
	var v2564 int32
	_ = v2564
	var v2567 int32
	_ = v2567
	var v2571 int32
	_ = v2571
	var v2575 int32
	_ = v2575
	var v2591 int32
	_ = v2591
	var v2595 int32
	_ = v2595
	var v2597 int32
	_ = v2597
	var v2601 int32
	_ = v2601
	var v2603 int32
	_ = v2603
	var v2607 int32
	_ = v2607
	var v2611 int32
	_ = v2611
	var v2612 int32
	_ = v2612
	var v2615 int32
	_ = v2615
	var v2619 int32
	_ = v2619
	var v2623 int32
	_ = v2623
	var v2629 int32
	_ = v2629
	var v2639 int32
	_ = v2639
	var v2644 int32
	_ = v2644
	var v2646 int32
	_ = v2646
	var v2649 int32
	_ = v2649
	var v2651 int32
	_ = v2651
	var v2655 int32
	_ = v2655
	var v2659 int32
	_ = v2659
	var v2660 int32
	_ = v2660
	var v2663 int32
	_ = v2663
	var v2667 int32
	_ = v2667
	var v2671 int32
	_ = v2671
	var v2681 int32
	_ = v2681
	var v2689 int32
	_ = v2689
	var v2693 int32
	_ = v2693
	var v2695 int32
	_ = v2695
	var v2699 int32
	_ = v2699
	var v2701 int32
	_ = v2701
	var v2705 int32
	_ = v2705
	var v2709 int32
	_ = v2709
	var v2710 int32
	_ = v2710
	var v2713 int32
	_ = v2713
	var v2717 int32
	_ = v2717
	var v2721 int32
	_ = v2721
	var v2737 int32
	_ = v2737
	var v2741 int32
	_ = v2741
	var v2743 int32
	_ = v2743
	var v2747 int32
	_ = v2747
	var v2749 int32
	_ = v2749
	var v2753 int32
	_ = v2753
	var v2757 int32
	_ = v2757
	var v2758 int32
	_ = v2758
	var v2761 int32
	_ = v2761
	var v2765 int32
	_ = v2765
	var v2769 int32
	_ = v2769
	var v2775 int32
	_ = v2775
	var v2785 int32
	_ = v2785
	var v2790 int32
	_ = v2790
	var v2792 int32
	_ = v2792
	var v2795 int32
	_ = v2795
	var v2797 int32
	_ = v2797
	var v2801 int32
	_ = v2801
	var v2805 int32
	_ = v2805
	var v2806 int32
	_ = v2806
	var v2809 int32
	_ = v2809
	var v2813 int32
	_ = v2813
	var v2817 int32
	_ = v2817
	var v2827 int32
	_ = v2827
	var v2835 int32
	_ = v2835
	var v2839 int32
	_ = v2839
	var v2841 int32
	_ = v2841
	var v2845 int32
	_ = v2845
	var v2847 int32
	_ = v2847
	var v2851 int32
	_ = v2851
	var v2855 int32
	_ = v2855
	var v2856 int32
	_ = v2856
	var v2859 int32
	_ = v2859
	var v2863 int32
	_ = v2863
	var v2867 int32
	_ = v2867
	var v2883 int32
	_ = v2883
	var v2887 int32
	_ = v2887
	var v2889 int32
	_ = v2889
	var v2893 int32
	_ = v2893
	var v2895 int32
	_ = v2895
	var v2899 int32
	_ = v2899
	var v2903 int32
	_ = v2903
	var v2904 int32
	_ = v2904
	var v2907 int32
	_ = v2907
	var v2911 int32
	_ = v2911
	var v2915 int32
	_ = v2915
	var v2921 int32
	_ = v2921
	var v2931 int32
	_ = v2931
	var v2936 int32
	_ = v2936
	var v2938 int32
	_ = v2938
	var v2941 int32
	_ = v2941
	var v2943 int32
	_ = v2943
	var v2947 int32
	_ = v2947
	var v2951 int32
	_ = v2951
	var v2952 int32
	_ = v2952
	var v2955 int32
	_ = v2955
	var v2959 int32
	_ = v2959
	var v2963 int32
	_ = v2963
	var v2973 int32
	_ = v2973
	var v2981 int32
	_ = v2981
	var v2985 int32
	_ = v2985
	var v2987 int32
	_ = v2987
	var v2991 int32
	_ = v2991
	var v2993 int32
	_ = v2993
	var v2997 int32
	_ = v2997
	var v3001 int32
	_ = v3001
	var v3002 int32
	_ = v3002
	var v3005 int32
	_ = v3005
	var v3009 int32
	_ = v3009
	var v3013 int32
	_ = v3013
	var v3029 int32
	_ = v3029
	var v3033 int32
	_ = v3033
	var v3035 int32
	_ = v3035
	var v3039 int32
	_ = v3039
	var v3041 int32
	_ = v3041
	var v3045 int32
	_ = v3045
	var v3049 int32
	_ = v3049
	var v3050 int32
	_ = v3050
	var v3053 int32
	_ = v3053
	var v3057 int32
	_ = v3057
	var v3061 int32
	_ = v3061
	var v3067 int32
	_ = v3067
	var v3080 int32
	_ = v3080
	var v3082 int32
	_ = v3082
	var v3095 int32
	_ = v3095
	var v3096 int32
	_ = v3096
	var v3099 int32
	_ = v3099
	var v3103 int32
	_ = v3103
	var v3107 int32
	_ = v3107
	var v3115 int32
	_ = v3115
	var v3116 int32
	_ = v3116
	var v3117 int32
	_ = v3117
	var v3118 int32
	_ = v3118
	var v3125 int32
	_ = v3125
	var v3129 int32
	_ = v3129
	var v3132 int32
	_ = v3132
	var v3134 int32
	_ = v3134
	var v3138 int32
	_ = v3138
	var v3139 int32
	_ = v3139
	var v3142 int32
	_ = v3142
	var v3146 int32
	_ = v3146
	var v3150 int32
	_ = v3150
	var v3156 int32
	_ = v3156
	var v3164 int32
	_ = v3164
	var v3165 int32
	_ = v3165
	var v3168 int32
	_ = v3168
	var v3169 int32
	_ = v3169
	var v3174 int32
	_ = v3174
	var v3182 int32
	_ = v3182
	var v3185 int64
	_ = v3185
	var v3187 int64
	_ = v3187
	var v3189 int64
	_ = v3189
	var v3190 int64
	_ = v3190
	var v3196 int32
	_ = v3196
	var v3205 int64
	_ = v3205
	var v3207 int32
	_ = v3207
	var v3212 int32
	_ = v3212
	var v3216 int32
	_ = v3216
	var v3219 int32
	_ = v3219
	var v3241 int32
	_ = v3241
	var v3242 int32
	_ = v3242
	var v3247 int32
	_ = v3247
	var v3259 int32
	_ = v3259
	var v3269 int32
	_ = v3269
	var v3271 int32
	_ = v3271
	var v3278 base.V128
	_ = v3278
	var v3285 int32
	_ = v3285
	var v3286 base.V128
	_ = v3286
	var v3289 base.V128
	_ = v3289
	var v3292 base.V128
	_ = v3292
	var v3295 base.V128
	_ = v3295
	var v3296 base.V128
	_ = v3296
	var v3300 int32
	_ = v3300
	var v3303 base.V128
	_ = v3303
	var v3308 int32
	_ = v3308
	var v3315 int32
	_ = v3315
	var v3321 int32
	_ = v3321
	var v3339 int32
	_ = v3339
	var v3344 int32
	_ = v3344
	var v3346 int32
	_ = v3346
	var v3354 int32
	_ = v3354
	var v3355 int32
	_ = v3355
	var v3359 int32
	_ = v3359
	var v3371 int32
	_ = v3371
	var v3383 int32
	_ = v3383
	var v3384 int32
	_ = v3384
	var v3386 int32
	_ = v3386
	var v3387 int32
	_ = v3387
	var v3388 int32
	_ = v3388
	var v3391 int32
	_ = v3391
	var v3398 int32
	_ = v3398
	var v3403 int32
	_ = v3403
	var v3408 int32
	_ = v3408
	var v3409 int32
	_ = v3409
	var v3411 int32
	_ = v3411
	var v3420 int32
	_ = v3420
	var v3422 int64
	_ = v3422
	var v3428 int32
	_ = v3428
	var v3429 int32
	_ = v3429
	var v3434 int32
	_ = v3434
	var v3437 int32
	_ = v3437
	var v3438 int32
	_ = v3438
	var v3441 int32
	_ = v3441
	var v3442 base.V128
	_ = v3442
	var v3443 int32
	_ = v3443
	var v3447 int32
	_ = v3447
	var v3448 int32
	_ = v3448
	var v3449 int32
	_ = v3449
	var v3452 int32
	_ = v3452
	var v3454 int32
	_ = v3454
	var v3456 int32
	_ = v3456
	var v3458 int32
	_ = v3458
	var v3460 int32
	_ = v3460
	var v3463 int32
	_ = v3463
	var v3466 int32
	_ = v3466
	var v3471 int32
	_ = v3471
	var v3472 int32
	_ = v3472
	var v3477 int32
	_ = v3477
	var v3480 int32
	_ = v3480
	var v3483 int32
	_ = v3483
	var v3485 int32
	_ = v3485
	var v3489 int32
	_ = v3489
	var v3490 int32
	_ = v3490
	var v3493 int32
	_ = v3493
	var v3494 int32
	_ = v3494
	var v3495 int32
	_ = v3495
	var v3496 int32
	_ = v3496
	var v3499 int32
	_ = v3499
	var v3504 int32
	_ = v3504
	var v3509 int32
	_ = v3509
	var v3510 int32
	_ = v3510
	var v3520 int32
	_ = v3520
	var v3522 int32
	_ = v3522
	var v3523 int32
	_ = v3523
	var v3525 int32
	_ = v3525
	var v3527 int32
	_ = v3527
	var v3532 int32
	_ = v3532
	var v3533 int32
	_ = v3533
	var v3539 int32
	_ = v3539
	var v3540 int32
	_ = v3540
	var v3542 int32
	_ = v3542
	var v3544 int32
	_ = v3544
	var v3545 int32
	_ = v3545
	var v3552 int32
	_ = v3552
	var v3553 int32
	_ = v3553
	var v3554 int32
	_ = v3554
	var v3556 int32
	_ = v3556
	var v3558 int32
	_ = v3558
	var v3561 int32
	_ = v3561
	var v3566 int32
	_ = v3566
	var v3567 int32
	_ = v3567
	var v3570 int32
	_ = v3570
	var v3572 int32
	_ = v3572
	var v3580 int32
	_ = v3580
	var v3581 int32
	_ = v3581
	var v3583 int32
	_ = v3583
	var v3592 int32
	_ = v3592
	var v3597 int32
	_ = v3597
	var v3598 int32
	_ = v3598
	var v3602 int32
	_ = v3602
	var v3612 int32
	_ = v3612
	var v3616 int32
	_ = v3616
	var v3623 int32
	_ = v3623
	var v3628 int32
	_ = v3628
	var v3635 int32
	_ = v3635
	var v3640 int32
	_ = v3640
	var v3641 int32
	_ = v3641
	var v3646 int32
	_ = v3646
	var v3658 int32
	_ = v3658
	var v3681 int32
	_ = v3681
	var v3690 int32
	_ = v3690
	var v3696 int32
	_ = v3696
	var v3697 int32
	_ = v3697
	var v3699 base.V128
	_ = v3699
	var v3700 int32
	_ = v3700
	var v3706 int32
	_ = v3706
	var v3716 int32
	_ = v3716
	var v3721 int32
	_ = v3721
	var v3723 int32
	_ = v3723
	var v3731 int32
	_ = v3731
	var v3737 int32
	_ = v3737
	var v3738 int32
	_ = v3738
	var v3739 int32
	_ = v3739
	var v3742 int32
	_ = v3742
	var v3744 int32
	_ = v3744
	var v3750 int32
	_ = v3750
	var v3751 int32
	_ = v3751
	var v3753 base.V128
	_ = v3753
	var v3754 int32
	_ = v3754
	var v3761 int32
	_ = v3761
	var v3763 int32
	_ = v3763
	var v3766 int32
	_ = v3766
	var v3770 int32
	_ = v3770
	var v3771 int32
	_ = v3771
	var v3784 int32
	_ = v3784
	var v3796 int32
	_ = v3796
	var v3801 int32
	_ = v3801
	var v3802 int32
	_ = v3802
	var v3805 int32
	_ = v3805
	var v3813 int32
	_ = v3813
	var v3818 int32
	_ = v3818
	var v3841 int32
	_ = v3841
	v2 = int32(0)
	v19 = base.Simd_g_const(&F_VP8EncWrite__k0)
	v20 = m.G0
	v22 = v20 - int32(48)
	m.G0 = v22
	v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v27 = base.I32_div_s(int32(19), v26)
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(76))))
	v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v35 = l0 + int32(56)
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v42 = base.I32_div_s(v36*v37*int32(7), int32(8))
	v45 = base.Simd_g_const(&F_VP8EncWrite__k1)
	base.Simd_g_v128_store(m, v35, v2, v45)
	base.Simd_g_v128_store(m, v35, int32(16), v19)
	if v42 == v2 {
		goto L7
	} else {
		goto L8
	}
L1:
	;
	m.G0 = v22 + int32(48)
	return v3841
L2:
	;
	v3841 = int32(0)
	goto L1
L3:
	;
	v3241 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(76))))
	v3242 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v3247 = v3241 + v3242*int32(3) + int32(7)
	if v3242 < int32(1) {
		v3371 = v3247
		goto L573
	} else {
		goto L574
	}
L4:
	;
	if v3219 == int32(0) {
		goto L2
	} else {
		goto L572
	}
L5:
	;
	v89 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v91 = v89 >> (uint(int32(1)) % 32)
	goto L21
L6:
	;
	if v80 != 0 {
		goto L5
	} else {
		goto L16
	}
L7:
	;
	v80 = int32(1)
	goto L6
L8:
	;
	v54 = int32(1024)
	if base.Ui32(v54) < base.Ui32(v42) {
		goto L10
	} else {
		goto L11
	}
L9:
	;
	v62 = *(*int32)(unsafe.Add(mBase, uint32(v35)+20))
	if v62 == int32(0) {
		goto L14
	} else {
		goto L15
	}
L10:
	;
	v57 = v42
	goto L12
L11:
	;
	v57 = v54
	goto L12
L12:
	;
	v58 = F_WebPSafeMalloc(m, int64(1), v57)
	mBase = m.M
	if v58 != 0 {
		goto L9
	} else {
		goto L13
	}
L13:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35)+28)) = int32(1)
	v80 = int32(0)
	goto L6
L14:
	;
	v69 = *(*int32)(unsafe.Add(mBase, uint32(v35)+16))
	F_WebPSafeFree(m, v69)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v35)+24)) = v57
	*(*int32)(unsafe.Add(mBase, uint32(v35)+16)) = v58
	goto L7
L15:
	;
	v67 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(72))))
	v68 = F_memcpy(m, v58, v67, v62)
	mBase = m.M
	goto L14
L16:
	;
	v81 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v83 = *(*int32)(unsafe.Add(mBase, uint32(v81)+92))
	if v83 != 0 {
		goto L18
	} else {
		goto L19
	}
L17:
	;
	v3219 = int32(0)
	goto L4
L18:
	;
	goto L17
L19:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v81)+92)) = int32(1)
	goto L18
L20:
	;
	v125 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v127 = v125 >> (uint(int32(1)) % 32)
	goto L27
L21:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v91
	if int32(126) < v91 {
		goto L23
	} else {
		goto L24
	}
L23:
	;
	goto L20
L24:
	;
	v104 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v105 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v104 << (uint(v105) % 32)
	v108 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v108 + v105
	v112 = m.G1
	v116 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v112+int32(_a_F_VP8EncWrite_0)+v91))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v116
	if v108 < int32(0) {
		goto L23
	} else {
		goto L25
	}
L25:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L23
L26:
	;
	v159 = l0 + int32(3416)
	v160 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v161 = int32(1)
	v162 = base.B2i32(v161 < v160)
	v165 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v167 = v165 >> (uint(v161) % 32)
	if v162 == int32(0) {
		v176 = v167
		goto L34
	} else {
		goto L35
	}
L27:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v127
	if int32(126) < v127 {
		goto L29
	} else {
		goto L30
	}
L29:
	;
	goto L26
L30:
	;
	v140 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v141 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v140 << (uint(v141) % 32)
	v144 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v144 + v141
	v148 = m.G1
	v152 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v148+int32(_a_F_VP8EncWrite_0)+v127))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v152
	if v144 < int32(0) {
		goto L29
	} else {
		goto L31
	}
L31:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L29
L32:
	;
	v1760 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v1761 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1764 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1766 = v1764 >> (uint(int32(1)) % 32)
	if v1761 == int32(0) {
		v1775 = v1766
		goto L320
	} else {
		goto L321
	}
L33:
	;
	if v162 == int32(0) {
		goto L32
	} else {
		goto L39
	}
L34:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v176
	if int32(126) < v176 {
		goto L36
	} else {
		goto L37
	}
L35:
	;
	v170 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v172 = v167 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v170 + v172
	v176 = v165 - v172
	goto L34
L36:
	;
	goto L33
L37:
	;
	v180 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v181 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v180 << (uint(v181) % 32)
	v184 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v184 + v181
	v188 = m.G1
	v192 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v188+int32(_a_F_VP8EncWrite_0)+v176))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v192
	if v184 < int32(0) {
		goto L36
	} else {
		goto L38
	}
L38:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L36
L39:
	;
	v200 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v203 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v205 = v203 >> (uint(int32(1)) % 32)
	if v200 == int32(0) {
		v214 = v205
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v239 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	goto L49
L41:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v214
	if int32(126) < v214 {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v208 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v210 = v205 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v208 + v210
	v214 = v203 - v210
	goto L41
L43:
	;
	goto L40
L44:
	;
	v218 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v219 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v218 << (uint(v219) % 32)
	v222 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v222 + v219
	v226 = m.G1
	v230 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v226+int32(_a_F_VP8EncWrite_0)+v214))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v230
	if v222 < int32(0) {
		goto L43
	} else {
		goto L45
	}
L45:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L43
L46:
	;
	v1478 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	if v1478 == int32(0) {
		goto L32
	} else {
		goto L268
	}
L47:
	;
	goto L53
L48:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v249
	if int32(126) < v249 {
		goto L50
	} else {
		goto L51
	}
L49:
	;
	v244 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v246 = v239>>(uint(int32(1))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v244 + v246
	v249 = v239 - v246
	goto L48
L50:
	;
	goto L47
L51:
	;
	v254 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v255 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v254 << (uint(v255) % 32)
	v258 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v258 + v255
	v262 = m.G1
	v266 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v262+int32(_a_F_VP8EncWrite_0)+v249))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v266
	if v258 < int32(0) {
		goto L50
	} else {
		goto L52
	}
L52:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L50
L53:
	;
	v277 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	goto L56
L54:
	;
	v310 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
	v315 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v317 = v315 >> (uint(int32(1)) % 32)
	if v310 == int32(0) {
		v326 = v317
		goto L61
	} else {
		goto L62
	}
L55:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v287
	if int32(126) < v287 {
		goto L57
	} else {
		goto L58
	}
L56:
	;
	v282 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v284 = v277>>(uint(int32(1))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v282 + v284
	v287 = v277 - v284
	goto L55
L57:
	;
	goto L54
L58:
	;
	v292 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v293 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v292 << (uint(v293) % 32)
	v296 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v296 + v293
	v300 = m.G1
	v304 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v300+int32(_a_F_VP8EncWrite_0)+v287))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v304
	if v296 < int32(0) {
		goto L57
	} else {
		goto L59
	}
L59:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L57
L60:
	;
	v456 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1832))
	v461 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v463 = v461 >> (uint(int32(1)) % 32)
	if v456 == int32(0) {
		v472 = v463
		goto L87
	} else {
		goto L88
	}
L61:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v326
	if int32(126) < v326 {
		goto L63
	} else {
		goto L64
	}
L62:
	;
	v320 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v322 = v317 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v320 + v322
	v326 = v315 - v322
	goto L61
L63:
	;
	if v310 == int32(0) {
		goto L66
	} else {
		goto L67
	}
L64:
	;
	v330 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v331 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v330 << (uint(v331) % 32)
	v334 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v334 + v331
	v338 = m.G1
	v342 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v338+int32(_a_F_VP8EncWrite_0)+v326))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v342
	if v334 < int32(0) {
		goto L63
	} else {
		goto L65
	}
L65:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L63
L66:
	;
	goto L60
L67:
	;
	if int32(-1) < v310 {
		goto L68
	} else {
		goto L69
	}
L68:
	;
	v408 = int32(128)
	goto L78
L69:
	;
	v352 = int32(1)
	v360 = int32(128)
	goto L70
L70:
	;
	v364 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v366 = v364 >> (uint(int32(1)) % 32)
	if v360&(v352-v310<<(uint(v352)%32)) == int32(0) {
		v376 = v366
		goto L72
	} else {
		goto L73
	}
L72:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v376
	if int32(126) < v376 {
		goto L74
	} else {
		goto L75
	}
L73:
	;
	v370 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v372 = v366 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v370 + v372
	v376 = v364 - v372
	goto L72
L74:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v360) {
		v360 = int32(base.Ui32(v360) >> (uint(int32(1)) % 32))
		goto L70
	} else {
		goto L77
	}
L75:
	;
	v380 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v381 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v380 << (uint(v381) % 32)
	v384 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v384 + v381
	v388 = m.G1
	v392 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v388+int32(_a_F_VP8EncWrite_0)+v376))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v392
	if v384 < int32(0) {
		goto L74
	} else {
		goto L76
	}
L76:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L74
L77:
	;
	goto L66
L78:
	;
	v412 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v414 = v412 >> (uint(int32(1)) % 32)
	if v408&(v310<<(uint(int32(1))%32)) == int32(0) {
		v424 = v414
		goto L80
	} else {
		goto L81
	}
L79:
	;
	goto L66
L80:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v424
	if int32(126) < v424 {
		goto L82
	} else {
		goto L83
	}
L81:
	;
	v418 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v420 = v414 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v418 + v420
	v424 = v412 - v420
	goto L80
L82:
	;
	v446 = int32(1)
	if base.Ui32(v446) < base.Ui32(v408) {
		v408 = int32(base.Ui32(v408) >> (uint(v446) % 32))
		goto L78
	} else {
		goto L85
	}
L83:
	;
	v428 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v429 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v428 << (uint(v429) % 32)
	v432 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v432 + v429
	v436 = m.G1
	v440 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v436+int32(_a_F_VP8EncWrite_0)+v424))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v440
	if v432 < int32(0) {
		goto L82
	} else {
		goto L84
	}
L84:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L82
L85:
	;
	goto L79
L86:
	;
	v602 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2576))
	v607 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v609 = v607 >> (uint(int32(1)) % 32)
	if v602 == int32(0) {
		v618 = v609
		goto L113
	} else {
		goto L114
	}
L87:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v472
	if int32(126) < v472 {
		goto L89
	} else {
		goto L90
	}
L88:
	;
	v466 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v468 = v463 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v466 + v468
	v472 = v461 - v468
	goto L87
L89:
	;
	if v456 == int32(0) {
		goto L92
	} else {
		goto L93
	}
L90:
	;
	v476 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v477 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v476 << (uint(v477) % 32)
	v480 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v480 + v477
	v484 = m.G1
	v488 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v484+int32(_a_F_VP8EncWrite_0)+v472))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v488
	if v480 < int32(0) {
		goto L89
	} else {
		goto L91
	}
L91:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L89
L92:
	;
	goto L86
L93:
	;
	if int32(-1) < v456 {
		goto L94
	} else {
		goto L95
	}
L94:
	;
	v554 = int32(128)
	goto L104
L95:
	;
	v498 = int32(1)
	v506 = int32(128)
	goto L96
L96:
	;
	v510 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v512 = v510 >> (uint(int32(1)) % 32)
	if v506&(v498-v456<<(uint(v498)%32)) == int32(0) {
		v522 = v512
		goto L98
	} else {
		goto L99
	}
L98:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v522
	if int32(126) < v522 {
		goto L100
	} else {
		goto L101
	}
L99:
	;
	v516 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v518 = v512 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v516 + v518
	v522 = v510 - v518
	goto L98
L100:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v506) {
		v506 = int32(base.Ui32(v506) >> (uint(int32(1)) % 32))
		goto L96
	} else {
		goto L103
	}
L101:
	;
	v526 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v527 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v526 << (uint(v527) % 32)
	v530 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v530 + v527
	v534 = m.G1
	v538 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v534+int32(_a_F_VP8EncWrite_0)+v522))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v538
	if v530 < int32(0) {
		goto L100
	} else {
		goto L102
	}
L102:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L100
L103:
	;
	goto L92
L104:
	;
	v558 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v560 = v558 >> (uint(int32(1)) % 32)
	if v554&(v456<<(uint(int32(1))%32)) == int32(0) {
		v570 = v560
		goto L106
	} else {
		goto L107
	}
L105:
	;
	goto L92
L106:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v570
	if int32(126) < v570 {
		goto L108
	} else {
		goto L109
	}
L107:
	;
	v564 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v566 = v560 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v564 + v566
	v570 = v558 - v566
	goto L106
L108:
	;
	v592 = int32(1)
	if base.Ui32(v592) < base.Ui32(v554) {
		v554 = int32(base.Ui32(v554) >> (uint(v592) % 32))
		goto L104
	} else {
		goto L111
	}
L109:
	;
	v574 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v575 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v574 << (uint(v575) % 32)
	v578 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v578 + v575
	v582 = m.G1
	v586 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v582+int32(_a_F_VP8EncWrite_0)+v570))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v586
	if v578 < int32(0) {
		goto L108
	} else {
		goto L110
	}
L110:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L108
L111:
	;
	goto L105
L112:
	;
	v748 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3320))
	v753 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v755 = v753 >> (uint(int32(1)) % 32)
	if v748 == int32(0) {
		v764 = v755
		goto L139
	} else {
		goto L140
	}
L113:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v618
	if int32(126) < v618 {
		goto L115
	} else {
		goto L116
	}
L114:
	;
	v612 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v614 = v609 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v612 + v614
	v618 = v607 - v614
	goto L113
L115:
	;
	if v602 == int32(0) {
		goto L118
	} else {
		goto L119
	}
L116:
	;
	v622 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v623 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v622 << (uint(v623) % 32)
	v626 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v626 + v623
	v630 = m.G1
	v634 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v630+int32(_a_F_VP8EncWrite_0)+v618))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v634
	if v626 < int32(0) {
		goto L115
	} else {
		goto L117
	}
L117:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L115
L118:
	;
	goto L112
L119:
	;
	if int32(-1) < v602 {
		goto L120
	} else {
		goto L121
	}
L120:
	;
	v700 = int32(128)
	goto L130
L121:
	;
	v644 = int32(1)
	v652 = int32(128)
	goto L122
L122:
	;
	v656 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v658 = v656 >> (uint(int32(1)) % 32)
	if v652&(v644-v602<<(uint(v644)%32)) == int32(0) {
		v668 = v658
		goto L124
	} else {
		goto L125
	}
L124:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v668
	if int32(126) < v668 {
		goto L126
	} else {
		goto L127
	}
L125:
	;
	v662 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v664 = v658 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v662 + v664
	v668 = v656 - v664
	goto L124
L126:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v652) {
		v652 = int32(base.Ui32(v652) >> (uint(int32(1)) % 32))
		goto L122
	} else {
		goto L129
	}
L127:
	;
	v672 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v673 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v672 << (uint(v673) % 32)
	v676 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v676 + v673
	v680 = m.G1
	v684 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v680+int32(_a_F_VP8EncWrite_0)+v668))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v684
	if v676 < int32(0) {
		goto L126
	} else {
		goto L128
	}
L128:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L126
L129:
	;
	goto L118
L130:
	;
	v704 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v706 = v704 >> (uint(int32(1)) % 32)
	if v700&(v602<<(uint(int32(1))%32)) == int32(0) {
		v716 = v706
		goto L132
	} else {
		goto L133
	}
L131:
	;
	goto L118
L132:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v716
	if int32(126) < v716 {
		goto L134
	} else {
		goto L135
	}
L133:
	;
	v710 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v712 = v706 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v710 + v712
	v716 = v704 - v712
	goto L132
L134:
	;
	v738 = int32(1)
	if base.Ui32(v738) < base.Ui32(v700) {
		v700 = int32(base.Ui32(v700) >> (uint(v738) % 32))
		goto L130
	} else {
		goto L137
	}
L135:
	;
	v720 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v721 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v720 << (uint(v721) % 32)
	v724 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v724 + v721
	v728 = m.G1
	v732 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v728+int32(_a_F_VP8EncWrite_0)+v716))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v732
	if v724 < int32(0) {
		goto L134
	} else {
		goto L136
	}
L136:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L134
L137:
	;
	goto L131
L138:
	;
	v894 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1092))
	v899 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v901 = v899 >> (uint(int32(1)) % 32)
	if v894 == int32(0) {
		v910 = v901
		goto L165
	} else {
		goto L166
	}
L139:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v764
	if int32(126) < v764 {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v758 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v760 = v755 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v758 + v760
	v764 = v753 - v760
	goto L139
L141:
	;
	if v748 == int32(0) {
		goto L144
	} else {
		goto L145
	}
L142:
	;
	v768 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v769 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v768 << (uint(v769) % 32)
	v772 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v772 + v769
	v776 = m.G1
	v780 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v776+int32(_a_F_VP8EncWrite_0)+v764))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v780
	if v772 < int32(0) {
		goto L141
	} else {
		goto L143
	}
L143:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L141
L144:
	;
	goto L138
L145:
	;
	if int32(-1) < v748 {
		goto L146
	} else {
		goto L147
	}
L146:
	;
	v846 = int32(128)
	goto L156
L147:
	;
	v790 = int32(1)
	v798 = int32(128)
	goto L148
L148:
	;
	v802 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v804 = v802 >> (uint(int32(1)) % 32)
	if v798&(v790-v748<<(uint(v790)%32)) == int32(0) {
		v814 = v804
		goto L150
	} else {
		goto L151
	}
L150:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v814
	if int32(126) < v814 {
		goto L152
	} else {
		goto L153
	}
L151:
	;
	v808 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v810 = v804 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v808 + v810
	v814 = v802 - v810
	goto L150
L152:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v798) {
		v798 = int32(base.Ui32(v798) >> (uint(int32(1)) % 32))
		goto L148
	} else {
		goto L155
	}
L153:
	;
	v818 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v819 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v818 << (uint(v819) % 32)
	v822 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v822 + v819
	v826 = m.G1
	v830 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v826+int32(_a_F_VP8EncWrite_0)+v814))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v830
	if v822 < int32(0) {
		goto L152
	} else {
		goto L154
	}
L154:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L152
L155:
	;
	goto L144
L156:
	;
	v850 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v852 = v850 >> (uint(int32(1)) % 32)
	if v846&(v748<<(uint(int32(1))%32)) == int32(0) {
		v862 = v852
		goto L158
	} else {
		goto L159
	}
L157:
	;
	goto L144
L158:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v862
	if int32(126) < v862 {
		goto L160
	} else {
		goto L161
	}
L159:
	;
	v856 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v858 = v852 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v856 + v858
	v862 = v850 - v858
	goto L158
L160:
	;
	v884 = int32(1)
	if base.Ui32(v884) < base.Ui32(v846) {
		v846 = int32(base.Ui32(v846) >> (uint(v884) % 32))
		goto L156
	} else {
		goto L163
	}
L161:
	;
	v866 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v867 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v866 << (uint(v867) % 32)
	v870 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v870 + v867
	v874 = m.G1
	v878 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v874+int32(_a_F_VP8EncWrite_0)+v862))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v878
	if v870 < int32(0) {
		goto L160
	} else {
		goto L162
	}
L162:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L160
L163:
	;
	goto L157
L164:
	;
	v1040 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1836))
	v1045 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1047 = v1045 >> (uint(int32(1)) % 32)
	if v1040 == int32(0) {
		v1056 = v1047
		goto L191
	} else {
		goto L192
	}
L165:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v910
	if int32(126) < v910 {
		goto L167
	} else {
		goto L168
	}
L166:
	;
	v904 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v906 = v901 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v904 + v906
	v910 = v899 - v906
	goto L165
L167:
	;
	if v894 == int32(0) {
		goto L170
	} else {
		goto L171
	}
L168:
	;
	v914 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v915 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v914 << (uint(v915) % 32)
	v918 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v918 + v915
	v922 = m.G1
	v926 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v922+int32(_a_F_VP8EncWrite_0)+v910))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v926
	if v918 < int32(0) {
		goto L167
	} else {
		goto L169
	}
L169:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L167
L170:
	;
	goto L164
L171:
	;
	if int32(-1) < v894 {
		goto L172
	} else {
		goto L173
	}
L172:
	;
	v992 = int32(64)
	goto L182
L173:
	;
	v936 = int32(1)
	v944 = int32(64)
	goto L174
L174:
	;
	v948 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v950 = v948 >> (uint(int32(1)) % 32)
	if v944&(v936-v894<<(uint(v936)%32)) == int32(0) {
		v960 = v950
		goto L176
	} else {
		goto L177
	}
L176:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v960
	if int32(126) < v960 {
		goto L178
	} else {
		goto L179
	}
L177:
	;
	v954 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v956 = v950 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v954 + v956
	v960 = v948 - v956
	goto L176
L178:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v944) {
		v944 = int32(base.Ui32(v944) >> (uint(int32(1)) % 32))
		goto L174
	} else {
		goto L181
	}
L179:
	;
	v964 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v965 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v964 << (uint(v965) % 32)
	v968 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v968 + v965
	v972 = m.G1
	v976 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v972+int32(_a_F_VP8EncWrite_0)+v960))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v976
	if v968 < int32(0) {
		goto L178
	} else {
		goto L180
	}
L180:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L178
L181:
	;
	goto L170
L182:
	;
	v996 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v998 = v996 >> (uint(int32(1)) % 32)
	if v992&(v894<<(uint(int32(1))%32)) == int32(0) {
		v1008 = v998
		goto L184
	} else {
		goto L185
	}
L183:
	;
	goto L170
L184:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1008
	if int32(126) < v1008 {
		goto L186
	} else {
		goto L187
	}
L185:
	;
	v1002 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1004 = v998 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1002 + v1004
	v1008 = v996 - v1004
	goto L184
L186:
	;
	v1030 = int32(1)
	if base.Ui32(v1030) < base.Ui32(v992) {
		v992 = int32(base.Ui32(v992) >> (uint(v1030) % 32))
		goto L182
	} else {
		goto L189
	}
L187:
	;
	v1012 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1013 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1012 << (uint(v1013) % 32)
	v1016 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1016 + v1013
	v1020 = m.G1
	v1024 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1020+int32(_a_F_VP8EncWrite_0)+v1008))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1024
	if v1016 < int32(0) {
		goto L186
	} else {
		goto L188
	}
L188:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L186
L189:
	;
	goto L183
L190:
	;
	v1186 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2580))
	v1191 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1193 = v1191 >> (uint(int32(1)) % 32)
	if v1186 == int32(0) {
		v1202 = v1193
		goto L217
	} else {
		goto L218
	}
L191:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1056
	if int32(126) < v1056 {
		goto L193
	} else {
		goto L194
	}
L192:
	;
	v1050 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1052 = v1047 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1050 + v1052
	v1056 = v1045 - v1052
	goto L191
L193:
	;
	if v1040 == int32(0) {
		goto L196
	} else {
		goto L197
	}
L194:
	;
	v1060 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1061 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1060 << (uint(v1061) % 32)
	v1064 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1064 + v1061
	v1068 = m.G1
	v1072 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1068+int32(_a_F_VP8EncWrite_0)+v1056))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1072
	if v1064 < int32(0) {
		goto L193
	} else {
		goto L195
	}
L195:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L193
L196:
	;
	goto L190
L197:
	;
	if int32(-1) < v1040 {
		goto L198
	} else {
		goto L199
	}
L198:
	;
	v1138 = int32(64)
	goto L208
L199:
	;
	v1082 = int32(1)
	v1090 = int32(64)
	goto L200
L200:
	;
	v1094 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1096 = v1094 >> (uint(int32(1)) % 32)
	if v1090&(v1082-v1040<<(uint(v1082)%32)) == int32(0) {
		v1106 = v1096
		goto L202
	} else {
		goto L203
	}
L202:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1106
	if int32(126) < v1106 {
		goto L204
	} else {
		goto L205
	}
L203:
	;
	v1100 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1102 = v1096 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1100 + v1102
	v1106 = v1094 - v1102
	goto L202
L204:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v1090) {
		v1090 = int32(base.Ui32(v1090) >> (uint(int32(1)) % 32))
		goto L200
	} else {
		goto L207
	}
L205:
	;
	v1110 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1111 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1110 << (uint(v1111) % 32)
	v1114 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1114 + v1111
	v1118 = m.G1
	v1122 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1118+int32(_a_F_VP8EncWrite_0)+v1106))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1122
	if v1114 < int32(0) {
		goto L204
	} else {
		goto L206
	}
L206:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L204
L207:
	;
	goto L196
L208:
	;
	v1142 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1144 = v1142 >> (uint(int32(1)) % 32)
	if v1138&(v1040<<(uint(int32(1))%32)) == int32(0) {
		v1154 = v1144
		goto L210
	} else {
		goto L211
	}
L209:
	;
	goto L196
L210:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1154
	if int32(126) < v1154 {
		goto L212
	} else {
		goto L213
	}
L211:
	;
	v1148 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1150 = v1144 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1148 + v1150
	v1154 = v1142 - v1150
	goto L210
L212:
	;
	v1176 = int32(1)
	if base.Ui32(v1176) < base.Ui32(v1138) {
		v1138 = int32(base.Ui32(v1138) >> (uint(v1176) % 32))
		goto L208
	} else {
		goto L215
	}
L213:
	;
	v1158 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1159 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1158 << (uint(v1159) % 32)
	v1162 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1162 + v1159
	v1166 = m.G1
	v1170 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1166+int32(_a_F_VP8EncWrite_0)+v1154))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1170
	if v1162 < int32(0) {
		goto L212
	} else {
		goto L214
	}
L214:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L212
L215:
	;
	goto L209
L216:
	;
	v1332 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3324))
	v1337 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1339 = v1337 >> (uint(int32(1)) % 32)
	if v1332 == int32(0) {
		v1348 = v1339
		goto L243
	} else {
		goto L244
	}
L217:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1202
	if int32(126) < v1202 {
		goto L219
	} else {
		goto L220
	}
L218:
	;
	v1196 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1198 = v1193 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1196 + v1198
	v1202 = v1191 - v1198
	goto L217
L219:
	;
	if v1186 == int32(0) {
		goto L222
	} else {
		goto L223
	}
L220:
	;
	v1206 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1207 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1206 << (uint(v1207) % 32)
	v1210 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1210 + v1207
	v1214 = m.G1
	v1218 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1214+int32(_a_F_VP8EncWrite_0)+v1202))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1218
	if v1210 < int32(0) {
		goto L219
	} else {
		goto L221
	}
L221:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L219
L222:
	;
	goto L216
L223:
	;
	if int32(-1) < v1186 {
		goto L224
	} else {
		goto L225
	}
L224:
	;
	v1284 = int32(64)
	goto L234
L225:
	;
	v1228 = int32(1)
	v1236 = int32(64)
	goto L226
L226:
	;
	v1240 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1242 = v1240 >> (uint(int32(1)) % 32)
	if v1236&(v1228-v1186<<(uint(v1228)%32)) == int32(0) {
		v1252 = v1242
		goto L228
	} else {
		goto L229
	}
L228:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1252
	if int32(126) < v1252 {
		goto L230
	} else {
		goto L231
	}
L229:
	;
	v1246 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1248 = v1242 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1246 + v1248
	v1252 = v1240 - v1248
	goto L228
L230:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v1236) {
		v1236 = int32(base.Ui32(v1236) >> (uint(int32(1)) % 32))
		goto L226
	} else {
		goto L233
	}
L231:
	;
	v1256 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1257 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1256 << (uint(v1257) % 32)
	v1260 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1260 + v1257
	v1264 = m.G1
	v1268 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1264+int32(_a_F_VP8EncWrite_0)+v1252))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1268
	if v1260 < int32(0) {
		goto L230
	} else {
		goto L232
	}
L232:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L230
L233:
	;
	goto L222
L234:
	;
	v1288 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1290 = v1288 >> (uint(int32(1)) % 32)
	if v1284&(v1186<<(uint(int32(1))%32)) == int32(0) {
		v1300 = v1290
		goto L236
	} else {
		goto L237
	}
L235:
	;
	goto L222
L236:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1300
	if int32(126) < v1300 {
		goto L238
	} else {
		goto L239
	}
L237:
	;
	v1294 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1296 = v1290 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1294 + v1296
	v1300 = v1288 - v1296
	goto L236
L238:
	;
	v1322 = int32(1)
	if base.Ui32(v1322) < base.Ui32(v1284) {
		v1284 = int32(base.Ui32(v1284) >> (uint(v1322) % 32))
		goto L234
	} else {
		goto L241
	}
L239:
	;
	v1304 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1305 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1304 << (uint(v1305) % 32)
	v1308 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1308 + v1305
	v1312 = m.G1
	v1316 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1312+int32(_a_F_VP8EncWrite_0)+v1300))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1316
	if v1308 < int32(0) {
		goto L238
	} else {
		goto L240
	}
L240:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L238
L241:
	;
	goto L235
L242:
	;
	goto L46
L243:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1348
	if int32(126) < v1348 {
		goto L245
	} else {
		goto L246
	}
L244:
	;
	v1342 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1344 = v1339 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1342 + v1344
	v1348 = v1337 - v1344
	goto L243
L245:
	;
	if v1332 == int32(0) {
		goto L248
	} else {
		goto L249
	}
L246:
	;
	v1352 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1353 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1352 << (uint(v1353) % 32)
	v1356 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1356 + v1353
	v1360 = m.G1
	v1364 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1360+int32(_a_F_VP8EncWrite_0)+v1348))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1364
	if v1356 < int32(0) {
		goto L245
	} else {
		goto L247
	}
L247:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L245
L248:
	;
	goto L242
L249:
	;
	if int32(-1) < v1332 {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v1430 = int32(64)
	goto L260
L251:
	;
	v1374 = int32(1)
	v1382 = int32(64)
	goto L252
L252:
	;
	v1386 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1388 = v1386 >> (uint(int32(1)) % 32)
	if v1382&(v1374-v1332<<(uint(v1374)%32)) == int32(0) {
		v1398 = v1388
		goto L254
	} else {
		goto L255
	}
L254:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1398
	if int32(126) < v1398 {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	v1392 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1394 = v1388 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1392 + v1394
	v1398 = v1386 - v1394
	goto L254
L256:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v1382) {
		v1382 = int32(base.Ui32(v1382) >> (uint(int32(1)) % 32))
		goto L252
	} else {
		goto L259
	}
L257:
	;
	v1402 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1403 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1402 << (uint(v1403) % 32)
	v1406 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1406 + v1403
	v1410 = m.G1
	v1414 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1410+int32(_a_F_VP8EncWrite_0)+v1398))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1414
	if v1406 < int32(0) {
		goto L256
	} else {
		goto L258
	}
L258:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L256
L259:
	;
	goto L248
L260:
	;
	v1434 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1436 = v1434 >> (uint(int32(1)) % 32)
	if v1430&(v1332<<(uint(int32(1))%32)) == int32(0) {
		v1446 = v1436
		goto L262
	} else {
		goto L263
	}
L261:
	;
	goto L248
L262:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1446
	if int32(126) < v1446 {
		goto L264
	} else {
		goto L265
	}
L263:
	;
	v1440 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1442 = v1436 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1440 + v1442
	v1446 = v1434 - v1442
	goto L262
L264:
	;
	v1468 = int32(1)
	if base.Ui32(v1468) < base.Ui32(v1430) {
		v1430 = int32(base.Ui32(v1430) >> (uint(v1468) % 32))
		goto L260
	} else {
		goto L267
	}
L265:
	;
	v1450 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1451 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1450 << (uint(v1451) % 32)
	v1454 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1454 + v1451
	v1458 = m.G1
	v1462 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1458+int32(_a_F_VP8EncWrite_0)+v1446))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1462
	if v1454 < int32(0) {
		goto L264
	} else {
		goto L266
	}
L266:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L264
L267:
	;
	goto L261
L268:
	;
	v1481 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v159))))
	v1483 = base.B2i32(v1481 != int32(255))
	v1486 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1488 = v1486 >> (uint(int32(1)) % 32)
	if v1483 == int32(0) {
		v1497 = v1488
		goto L271
	} else {
		goto L272
	}
L269:
	;
	v1574 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)))
	v1576 = base.B2i32(v1574 != int32(255))
	v1579 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1581 = v1579 >> (uint(int32(1)) % 32)
	if v1576 == int32(0) {
		v1590 = v1581
		goto L288
	} else {
		goto L289
	}
L270:
	;
	if v1483 == int32(0) {
		goto L269
	} else {
		goto L276
	}
L271:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1497
	if int32(126) < v1497 {
		goto L273
	} else {
		goto L274
	}
L272:
	;
	v1491 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1493 = v1488 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1491 + v1493
	v1497 = v1486 - v1493
	goto L271
L273:
	;
	goto L270
L274:
	;
	v1501 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1502 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1501 << (uint(v1502) % 32)
	v1505 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1505 + v1502
	v1509 = m.G1
	v1513 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1509+int32(_a_F_VP8EncWrite_0)+v1497))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1513
	if v1505 < int32(0) {
		goto L273
	} else {
		goto L275
	}
L275:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L273
L276:
	;
	v1521 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v159))))
	v1533 = int32(128)
	goto L278
L277:
	;
	goto L269
L278:
	;
	v1536 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1538 = v1536 >> (uint(int32(1)) % 32)
	if v1533&v1521 == int32(0) {
		v1548 = v1538
		goto L280
	} else {
		goto L281
	}
L279:
	;
	goto L277
L280:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1548
	if int32(126) < v1548 {
		goto L282
	} else {
		goto L283
	}
L281:
	;
	v1542 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1544 = v1538 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1542 + v1544
	v1548 = v1536 - v1544
	goto L280
L282:
	;
	v1570 = int32(1)
	if base.Ui32(v1570) < base.Ui32(v1533) {
		v1533 = int32(base.Ui32(v1533) >> (uint(v1570) % 32))
		goto L278
	} else {
		goto L285
	}
L283:
	;
	v1552 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1553 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1552 << (uint(v1553) % 32)
	v1556 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1556 + v1553
	v1560 = m.G1
	v1564 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1560+int32(_a_F_VP8EncWrite_0)+v1548))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1564
	if v1556 < int32(0) {
		goto L282
	} else {
		goto L284
	}
L284:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L282
L285:
	;
	goto L279
L286:
	;
	v1667 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)))
	v1669 = base.B2i32(v1667 != int32(255))
	v1672 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1674 = v1672 >> (uint(int32(1)) % 32)
	if v1669 == int32(0) {
		v1683 = v1674
		goto L304
	} else {
		goto L305
	}
L287:
	;
	if v1576 == int32(0) {
		goto L286
	} else {
		goto L293
	}
L288:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1590
	if int32(126) < v1590 {
		goto L290
	} else {
		goto L291
	}
L289:
	;
	v1584 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1586 = v1581 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1584 + v1586
	v1590 = v1579 - v1586
	goto L288
L290:
	;
	goto L287
L291:
	;
	v1594 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1595 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1594 << (uint(v1595) % 32)
	v1598 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1598 + v1595
	v1602 = m.G1
	v1606 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1602+int32(_a_F_VP8EncWrite_0)+v1590))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1606
	if v1598 < int32(0) {
		goto L290
	} else {
		goto L292
	}
L292:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L290
L293:
	;
	v1614 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)))
	v1626 = int32(128)
	goto L295
L294:
	;
	goto L286
L295:
	;
	v1629 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1631 = v1629 >> (uint(int32(1)) % 32)
	if v1626&v1614 == int32(0) {
		v1641 = v1631
		goto L297
	} else {
		goto L298
	}
L296:
	;
	goto L294
L297:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1641
	if int32(126) < v1641 {
		goto L299
	} else {
		goto L300
	}
L298:
	;
	v1635 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1637 = v1631 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1635 + v1637
	v1641 = v1629 - v1637
	goto L297
L299:
	;
	v1663 = int32(1)
	if base.Ui32(v1663) < base.Ui32(v1626) {
		v1626 = int32(base.Ui32(v1626) >> (uint(v1663) % 32))
		goto L295
	} else {
		goto L302
	}
L300:
	;
	v1645 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1646 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1645 << (uint(v1646) % 32)
	v1649 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1649 + v1646
	v1653 = m.G1
	v1657 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1653+int32(_a_F_VP8EncWrite_0)+v1641))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1657
	if v1649 < int32(0) {
		goto L299
	} else {
		goto L301
	}
L301:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L299
L302:
	;
	goto L296
L303:
	;
	if v1669 == int32(0) {
		goto L32
	} else {
		goto L309
	}
L304:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1683
	if int32(126) < v1683 {
		goto L306
	} else {
		goto L307
	}
L305:
	;
	v1677 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1679 = v1674 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1677 + v1679
	v1683 = v1672 - v1679
	goto L304
L306:
	;
	goto L303
L307:
	;
	v1687 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1688 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1687 << (uint(v1688) % 32)
	v1691 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1691 + v1688
	v1695 = m.G1
	v1699 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1695+int32(_a_F_VP8EncWrite_0)+v1683))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1699
	if v1691 < int32(0) {
		goto L306
	} else {
		goto L308
	}
L308:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L306
L309:
	;
	v1707 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)))
	v1719 = int32(128)
	goto L311
L310:
	;
	goto L32
L311:
	;
	v1722 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1724 = v1722 >> (uint(int32(1)) % 32)
	if v1719&v1707 == int32(0) {
		v1734 = v1724
		goto L313
	} else {
		goto L314
	}
L312:
	;
	goto L310
L313:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1734
	if int32(126) < v1734 {
		goto L315
	} else {
		goto L316
	}
L314:
	;
	v1728 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1730 = v1724 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1728 + v1730
	v1734 = v1722 - v1730
	goto L313
L315:
	;
	v1756 = int32(1)
	if base.Ui32(v1756) < base.Ui32(v1719) {
		v1719 = int32(base.Ui32(v1719) >> (uint(v1756) % 32))
		goto L311
	} else {
		goto L318
	}
L316:
	;
	v1738 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1739 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1738 << (uint(v1739) % 32)
	v1742 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1742 + v1739
	v1746 = m.G1
	v1750 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1746+int32(_a_F_VP8EncWrite_0)+v1734))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1750
	if v1742 < int32(0) {
		goto L315
	} else {
		goto L317
	}
L317:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L315
L318:
	;
	goto L312
L319:
	;
	v1797 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1809 = int32(32)
	goto L326
L320:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1775
	if int32(126) < v1775 {
		goto L322
	} else {
		goto L323
	}
L321:
	;
	v1769 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1771 = v1766 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1769 + v1771
	v1775 = v1764 - v1771
	goto L320
L322:
	;
	goto L319
L323:
	;
	v1779 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1780 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1779 << (uint(v1780) % 32)
	v1783 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1783 + v1780
	v1787 = m.G1
	v1791 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1787+int32(_a_F_VP8EncWrite_0)+v1775))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1791
	if v1783 < int32(0) {
		goto L322
	} else {
		goto L324
	}
L324:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L322
L325:
	;
	v1851 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1863 = int32(4)
	goto L335
L326:
	;
	v1812 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1814 = v1812 >> (uint(int32(1)) % 32)
	if v1809&v1797 == int32(0) {
		v1824 = v1814
		goto L328
	} else {
		goto L329
	}
L327:
	;
	goto L325
L328:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1824
	if int32(126) < v1824 {
		goto L330
	} else {
		goto L331
	}
L329:
	;
	v1818 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1820 = v1814 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1818 + v1820
	v1824 = v1812 - v1820
	goto L328
L330:
	;
	v1846 = int32(1)
	if base.Ui32(v1846) < base.Ui32(v1809) {
		v1809 = int32(base.Ui32(v1809) >> (uint(v1846) % 32))
		goto L326
	} else {
		goto L333
	}
L331:
	;
	v1828 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1829 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1828 << (uint(v1829) % 32)
	v1832 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1832 + v1829
	v1836 = m.G1
	v1840 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1836+int32(_a_F_VP8EncWrite_0)+v1824))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1840
	if v1832 < int32(0) {
		goto L330
	} else {
		goto L332
	}
L332:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L330
L333:
	;
	goto L327
L334:
	;
	v1904 = int32(0)
	v1905 = base.B2i32(v1760 != v1904)
	v1908 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1910 = v1908 >> (uint(int32(1)) % 32)
	if v1905 == v1904 {
		v1919 = v1910
		goto L345
	} else {
		goto L346
	}
L335:
	;
	v1866 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1868 = v1866 >> (uint(int32(1)) % 32)
	if v1863&v1851 == int32(0) {
		v1878 = v1868
		goto L337
	} else {
		goto L338
	}
L336:
	;
	goto L334
L337:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1878
	if int32(126) < v1878 {
		goto L339
	} else {
		goto L340
	}
L338:
	;
	v1872 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1874 = v1868 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1872 + v1874
	v1878 = v1866 - v1874
	goto L337
L339:
	;
	v1900 = int32(1)
	if base.Ui32(v1900) < base.Ui32(v1863) {
		v1863 = int32(base.Ui32(v1863) >> (uint(v1900) % 32))
		goto L335
	} else {
		goto L342
	}
L340:
	;
	v1882 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1883 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1882 << (uint(v1883) % 32)
	v1886 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1886 + v1883
	v1890 = m.G1
	v1894 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1890+int32(_a_F_VP8EncWrite_0)+v1878))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1894
	if v1886 < int32(0) {
		goto L339
	} else {
		goto L341
	}
L341:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L339
L342:
	;
	goto L336
L343:
	;
	v2235 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	switch v2235 + int32(-4) {
	case 0:
		goto L403
	default:
		goto L404
	case 4:
		v2241 = int32(3)
		goto L402
	}
L344:
	;
	if v1905 == int32(0) {
		goto L343
	} else {
		goto L350
	}
L345:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1919
	if int32(126) < v1919 {
		goto L347
	} else {
		goto L348
	}
L346:
	;
	v1913 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1915 = v1910 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1913 + v1915
	v1919 = v1908 - v1915
	goto L345
L347:
	;
	goto L344
L348:
	;
	v1923 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1924 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1923 << (uint(v1924) % 32)
	v1927 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1927 + v1924
	v1931 = m.G1
	v1935 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1931+int32(_a_F_VP8EncWrite_0)+v1919))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1935
	if v1927 < int32(0) {
		goto L347
	} else {
		goto L349
	}
L349:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L347
L350:
	;
	v1943 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v1944 = int32(0)
	v1945 = base.B2i32(v1943 != v1944)
	v1948 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v1950 = v1948 >> (uint(int32(1)) % 32)
	if v1945 == v1944 {
		v1959 = v1950
		goto L352
	} else {
		goto L353
	}
L351:
	;
	if v1945 == int32(0) {
		goto L343
	} else {
		goto L357
	}
L352:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1959
	if int32(126) < v1959 {
		goto L354
	} else {
		goto L355
	}
L353:
	;
	v1953 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1955 = v1950 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1953 + v1955
	v1959 = v1948 - v1955
	goto L352
L354:
	;
	goto L351
L355:
	;
	v1963 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v1964 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v1963 << (uint(v1964) % 32)
	v1967 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v1967 + v1964
	v1971 = m.G1
	v1975 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1971+int32(_a_F_VP8EncWrite_0)+v1959))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v1975
	if v1967 < int32(0) {
		goto L354
	} else {
		goto L356
	}
L356:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L354
L357:
	;
	v1995 = int32(8)
	goto L359
L358:
	;
	v2036 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v2041 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2043 = v2041 >> (uint(int32(1)) % 32)
	if v2036 == int32(0) {
		v2052 = v2043
		goto L368
	} else {
		goto L369
	}
L359:
	;
	v1998 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2000 = v1998 >> (uint(int32(1)) % 32)
	if v1995&int32(0) == int32(0) {
		v2010 = v2000
		goto L361
	} else {
		goto L362
	}
L360:
	;
	goto L358
L361:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2010
	if int32(126) < v2010 {
		goto L363
	} else {
		goto L364
	}
L362:
	;
	v2004 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2006 = v2000 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2004 + v2006
	v2010 = v1998 - v2006
	goto L361
L363:
	;
	v2032 = int32(1)
	if base.Ui32(v2032) < base.Ui32(v1995) {
		v1995 = int32(base.Ui32(v1995) >> (uint(v2032) % 32))
		goto L359
	} else {
		goto L366
	}
L364:
	;
	v2014 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2015 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2014 << (uint(v2015) % 32)
	v2018 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2018 + v2015
	v2022 = m.G1
	v2026 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2022+int32(_a_F_VP8EncWrite_0)+v2010))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2026
	if v2018 < int32(0) {
		goto L363
	} else {
		goto L365
	}
L365:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L363
L366:
	;
	goto L360
L367:
	;
	v2194 = int32(4)
	goto L394
L368:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2052
	if int32(126) < v2052 {
		goto L370
	} else {
		goto L371
	}
L369:
	;
	v2046 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2048 = v2043 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2046 + v2048
	v2052 = v2041 - v2048
	goto L368
L370:
	;
	if v2036 == int32(0) {
		goto L373
	} else {
		goto L374
	}
L371:
	;
	v2056 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2057 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2056 << (uint(v2057) % 32)
	v2060 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2060 + v2057
	v2064 = m.G1
	v2068 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2064+int32(_a_F_VP8EncWrite_0)+v2052))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2068
	if v2060 < int32(0) {
		goto L370
	} else {
		goto L372
	}
L372:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L370
L373:
	;
	goto L367
L374:
	;
	if int32(-1) < v2036 {
		goto L375
	} else {
		goto L376
	}
L375:
	;
	v2134 = int32(64)
	goto L385
L376:
	;
	v2078 = int32(1)
	v2086 = int32(64)
	goto L377
L377:
	;
	v2090 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2092 = v2090 >> (uint(int32(1)) % 32)
	if v2086&(v2078-v2036<<(uint(v2078)%32)) == int32(0) {
		v2102 = v2092
		goto L379
	} else {
		goto L380
	}
L379:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2102
	if int32(126) < v2102 {
		goto L381
	} else {
		goto L382
	}
L380:
	;
	v2096 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2098 = v2092 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2096 + v2098
	v2102 = v2090 - v2098
	goto L379
L381:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2086) {
		v2086 = int32(base.Ui32(v2086) >> (uint(int32(1)) % 32))
		goto L377
	} else {
		goto L384
	}
L382:
	;
	v2106 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2107 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2106 << (uint(v2107) % 32)
	v2110 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2110 + v2107
	v2114 = m.G1
	v2118 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2114+int32(_a_F_VP8EncWrite_0)+v2102))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2118
	if v2110 < int32(0) {
		goto L381
	} else {
		goto L383
	}
L383:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L381
L384:
	;
	goto L373
L385:
	;
	v2138 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2140 = v2138 >> (uint(int32(1)) % 32)
	if v2134&(v2036<<(uint(int32(1))%32)) == int32(0) {
		v2150 = v2140
		goto L387
	} else {
		goto L388
	}
L386:
	;
	goto L373
L387:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2150
	if int32(126) < v2150 {
		goto L389
	} else {
		goto L390
	}
L388:
	;
	v2144 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2146 = v2140 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2144 + v2146
	v2150 = v2138 - v2146
	goto L387
L389:
	;
	v2172 = int32(1)
	if base.Ui32(v2172) < base.Ui32(v2134) {
		v2134 = int32(base.Ui32(v2134) >> (uint(v2172) % 32))
		goto L385
	} else {
		goto L392
	}
L390:
	;
	v2154 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2155 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2154 << (uint(v2155) % 32)
	v2158 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2158 + v2155
	v2162 = m.G1
	v2166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2162+int32(_a_F_VP8EncWrite_0)+v2150))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2166
	if v2158 < int32(0) {
		goto L389
	} else {
		goto L391
	}
L391:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L389
L392:
	;
	goto L386
L393:
	;
	goto L343
L394:
	;
	v2197 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2199 = v2197 >> (uint(int32(1)) % 32)
	if v2194&int32(0) == int32(0) {
		v2209 = v2199
		goto L396
	} else {
		goto L397
	}
L395:
	;
	goto L393
L396:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2209
	if int32(126) < v2209 {
		goto L398
	} else {
		goto L399
	}
L397:
	;
	v2203 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2205 = v2199 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2203 + v2205
	v2209 = v2197 - v2205
	goto L396
L398:
	;
	v2231 = int32(1)
	if base.Ui32(v2231) < base.Ui32(v2194) {
		v2194 = int32(base.Ui32(v2194) >> (uint(v2231) % 32))
		goto L394
	} else {
		goto L401
	}
L399:
	;
	v2213 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2214 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2213 << (uint(v2214) % 32)
	v2217 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2217 + v2214
	v2221 = m.G1
	v2225 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2221+int32(_a_F_VP8EncWrite_0)+v2209))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2225
	if v2217 < int32(0) {
		goto L398
	} else {
		goto L400
	}
L400:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L398
L401:
	;
	goto L395
L402:
	;
	v2253 = int32(2)
	goto L406
L403:
	;
	v2241 = int32(2)
	goto L402
L404:
	;
	v2241 = base.B2i32(v2235 == int32(2))
	goto L402
L405:
	;
	v2294 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3384))
	v2306 = int32(64)
	goto L415
L406:
	;
	v2256 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2258 = v2256 >> (uint(int32(1)) % 32)
	if v2253&v2241 == int32(0) {
		v2268 = v2258
		goto L408
	} else {
		goto L409
	}
L407:
	;
	goto L405
L408:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2268
	if int32(126) < v2268 {
		goto L410
	} else {
		goto L411
	}
L409:
	;
	v2262 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2264 = v2258 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2262 + v2264
	v2268 = v2256 - v2264
	goto L408
L410:
	;
	v2290 = int32(1)
	if base.Ui32(v2290) < base.Ui32(v2253) {
		v2253 = int32(base.Ui32(v2253) >> (uint(v2290) % 32))
		goto L406
	} else {
		goto L413
	}
L411:
	;
	v2272 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2273 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2272 << (uint(v2273) % 32)
	v2276 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2276 + v2273
	v2280 = m.G1
	v2284 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2280+int32(_a_F_VP8EncWrite_0)+v2268))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2284
	if v2276 < int32(0) {
		goto L410
	} else {
		goto L412
	}
L412:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L410
L413:
	;
	goto L407
L414:
	;
	v2347 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3396))
	v2352 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2354 = v2352 >> (uint(int32(1)) % 32)
	if v2347 == int32(0) {
		v2363 = v2354
		goto L424
	} else {
		goto L425
	}
L415:
	;
	v2309 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2311 = v2309 >> (uint(int32(1)) % 32)
	if v2306&v2294 == int32(0) {
		v2321 = v2311
		goto L417
	} else {
		goto L418
	}
L416:
	;
	goto L414
L417:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2321
	if int32(126) < v2321 {
		goto L419
	} else {
		goto L420
	}
L418:
	;
	v2315 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2317 = v2311 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2315 + v2317
	v2321 = v2309 - v2317
	goto L417
L419:
	;
	v2343 = int32(1)
	if base.Ui32(v2343) < base.Ui32(v2306) {
		v2306 = int32(base.Ui32(v2306) >> (uint(v2343) % 32))
		goto L415
	} else {
		goto L422
	}
L420:
	;
	v2325 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2326 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2325 << (uint(v2326) % 32)
	v2329 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2329 + v2326
	v2333 = m.G1
	v2337 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2333+int32(_a_F_VP8EncWrite_0)+v2321))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2337
	if v2329 < int32(0) {
		goto L419
	} else {
		goto L421
	}
L421:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L419
L422:
	;
	goto L416
L423:
	;
	v2493 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3400))
	v2498 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2500 = v2498 >> (uint(int32(1)) % 32)
	if v2493 == int32(0) {
		v2509 = v2500
		goto L450
	} else {
		goto L451
	}
L424:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2363
	if int32(126) < v2363 {
		goto L426
	} else {
		goto L427
	}
L425:
	;
	v2357 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2359 = v2354 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2357 + v2359
	v2363 = v2352 - v2359
	goto L424
L426:
	;
	if v2347 == int32(0) {
		goto L429
	} else {
		goto L430
	}
L427:
	;
	v2367 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2368 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2367 << (uint(v2368) % 32)
	v2371 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2371 + v2368
	v2375 = m.G1
	v2379 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2375+int32(_a_F_VP8EncWrite_0)+v2363))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2379
	if v2371 < int32(0) {
		goto L426
	} else {
		goto L428
	}
L428:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L426
L429:
	;
	goto L423
L430:
	;
	if int32(-1) < v2347 {
		goto L431
	} else {
		goto L432
	}
L431:
	;
	v2445 = int32(16)
	goto L441
L432:
	;
	v2389 = int32(1)
	v2397 = int32(16)
	goto L433
L433:
	;
	v2401 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2403 = v2401 >> (uint(int32(1)) % 32)
	if v2397&(v2389-v2347<<(uint(v2389)%32)) == int32(0) {
		v2413 = v2403
		goto L435
	} else {
		goto L436
	}
L435:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2413
	if int32(126) < v2413 {
		goto L437
	} else {
		goto L438
	}
L436:
	;
	v2407 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2409 = v2403 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2407 + v2409
	v2413 = v2401 - v2409
	goto L435
L437:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2397) {
		v2397 = int32(base.Ui32(v2397) >> (uint(int32(1)) % 32))
		goto L433
	} else {
		goto L440
	}
L438:
	;
	v2417 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2418 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2417 << (uint(v2418) % 32)
	v2421 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2421 + v2418
	v2425 = m.G1
	v2429 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2425+int32(_a_F_VP8EncWrite_0)+v2413))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2429
	if v2421 < int32(0) {
		goto L437
	} else {
		goto L439
	}
L439:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L437
L440:
	;
	goto L429
L441:
	;
	v2449 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2451 = v2449 >> (uint(int32(1)) % 32)
	if v2445&(v2347<<(uint(int32(1))%32)) == int32(0) {
		v2461 = v2451
		goto L443
	} else {
		goto L444
	}
L442:
	;
	goto L429
L443:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2461
	if int32(126) < v2461 {
		goto L445
	} else {
		goto L446
	}
L444:
	;
	v2455 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2457 = v2451 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2455 + v2457
	v2461 = v2449 - v2457
	goto L443
L445:
	;
	v2483 = int32(1)
	if base.Ui32(v2483) < base.Ui32(v2445) {
		v2445 = int32(base.Ui32(v2445) >> (uint(v2483) % 32))
		goto L441
	} else {
		goto L448
	}
L446:
	;
	v2465 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2466 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2465 << (uint(v2466) % 32)
	v2469 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2469 + v2466
	v2473 = m.G1
	v2477 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2473+int32(_a_F_VP8EncWrite_0)+v2461))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2477
	if v2469 < int32(0) {
		goto L445
	} else {
		goto L447
	}
L447:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L445
L448:
	;
	goto L442
L449:
	;
	v2639 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3404))
	v2644 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2646 = v2644 >> (uint(int32(1)) % 32)
	if v2639 == int32(0) {
		v2655 = v2646
		goto L476
	} else {
		goto L477
	}
L450:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2509
	if int32(126) < v2509 {
		goto L452
	} else {
		goto L453
	}
L451:
	;
	v2503 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2505 = v2500 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2503 + v2505
	v2509 = v2498 - v2505
	goto L450
L452:
	;
	if v2493 == int32(0) {
		goto L455
	} else {
		goto L456
	}
L453:
	;
	v2513 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2514 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2513 << (uint(v2514) % 32)
	v2517 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2517 + v2514
	v2521 = m.G1
	v2525 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2521+int32(_a_F_VP8EncWrite_0)+v2509))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2525
	if v2517 < int32(0) {
		goto L452
	} else {
		goto L454
	}
L454:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L452
L455:
	;
	goto L449
L456:
	;
	if int32(-1) < v2493 {
		goto L457
	} else {
		goto L458
	}
L457:
	;
	v2591 = int32(16)
	goto L467
L458:
	;
	v2535 = int32(1)
	v2543 = int32(16)
	goto L459
L459:
	;
	v2547 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2549 = v2547 >> (uint(int32(1)) % 32)
	if v2543&(v2535-v2493<<(uint(v2535)%32)) == int32(0) {
		v2559 = v2549
		goto L461
	} else {
		goto L462
	}
L461:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2559
	if int32(126) < v2559 {
		goto L463
	} else {
		goto L464
	}
L462:
	;
	v2553 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2555 = v2549 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2553 + v2555
	v2559 = v2547 - v2555
	goto L461
L463:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2543) {
		v2543 = int32(base.Ui32(v2543) >> (uint(int32(1)) % 32))
		goto L459
	} else {
		goto L466
	}
L464:
	;
	v2563 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2564 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2563 << (uint(v2564) % 32)
	v2567 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2567 + v2564
	v2571 = m.G1
	v2575 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2571+int32(_a_F_VP8EncWrite_0)+v2559))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2575
	if v2567 < int32(0) {
		goto L463
	} else {
		goto L465
	}
L465:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L463
L466:
	;
	goto L455
L467:
	;
	v2595 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2597 = v2595 >> (uint(int32(1)) % 32)
	if v2591&(v2493<<(uint(int32(1))%32)) == int32(0) {
		v2607 = v2597
		goto L469
	} else {
		goto L470
	}
L468:
	;
	goto L455
L469:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2607
	if int32(126) < v2607 {
		goto L471
	} else {
		goto L472
	}
L470:
	;
	v2601 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2603 = v2597 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2601 + v2603
	v2607 = v2595 - v2603
	goto L469
L471:
	;
	v2629 = int32(1)
	if base.Ui32(v2629) < base.Ui32(v2591) {
		v2591 = int32(base.Ui32(v2591) >> (uint(v2629) % 32))
		goto L467
	} else {
		goto L474
	}
L472:
	;
	v2611 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2612 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2611 << (uint(v2612) % 32)
	v2615 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2615 + v2612
	v2619 = m.G1
	v2623 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2619+int32(_a_F_VP8EncWrite_0)+v2607))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2623
	if v2615 < int32(0) {
		goto L471
	} else {
		goto L473
	}
L473:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L471
L474:
	;
	goto L468
L475:
	;
	v2785 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3408))
	v2790 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2792 = v2790 >> (uint(int32(1)) % 32)
	if v2785 == int32(0) {
		v2801 = v2792
		goto L502
	} else {
		goto L503
	}
L476:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2655
	if int32(126) < v2655 {
		goto L478
	} else {
		goto L479
	}
L477:
	;
	v2649 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2651 = v2646 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2649 + v2651
	v2655 = v2644 - v2651
	goto L476
L478:
	;
	if v2639 == int32(0) {
		goto L481
	} else {
		goto L482
	}
L479:
	;
	v2659 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2660 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2659 << (uint(v2660) % 32)
	v2663 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2663 + v2660
	v2667 = m.G1
	v2671 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2667+int32(_a_F_VP8EncWrite_0)+v2655))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2671
	if v2663 < int32(0) {
		goto L478
	} else {
		goto L480
	}
L480:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L478
L481:
	;
	goto L475
L482:
	;
	if int32(-1) < v2639 {
		goto L483
	} else {
		goto L484
	}
L483:
	;
	v2737 = int32(16)
	goto L493
L484:
	;
	v2681 = int32(1)
	v2689 = int32(16)
	goto L485
L485:
	;
	v2693 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2695 = v2693 >> (uint(int32(1)) % 32)
	if v2689&(v2681-v2639<<(uint(v2681)%32)) == int32(0) {
		v2705 = v2695
		goto L487
	} else {
		goto L488
	}
L487:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2705
	if int32(126) < v2705 {
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v2699 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2701 = v2695 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2699 + v2701
	v2705 = v2693 - v2701
	goto L487
L489:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2689) {
		v2689 = int32(base.Ui32(v2689) >> (uint(int32(1)) % 32))
		goto L485
	} else {
		goto L492
	}
L490:
	;
	v2709 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2710 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2709 << (uint(v2710) % 32)
	v2713 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2713 + v2710
	v2717 = m.G1
	v2721 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2717+int32(_a_F_VP8EncWrite_0)+v2705))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2721
	if v2713 < int32(0) {
		goto L489
	} else {
		goto L491
	}
L491:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L489
L492:
	;
	goto L481
L493:
	;
	v2741 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2743 = v2741 >> (uint(int32(1)) % 32)
	if v2737&(v2639<<(uint(int32(1))%32)) == int32(0) {
		v2753 = v2743
		goto L495
	} else {
		goto L496
	}
L494:
	;
	goto L481
L495:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2753
	if int32(126) < v2753 {
		goto L497
	} else {
		goto L498
	}
L496:
	;
	v2747 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2749 = v2743 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2747 + v2749
	v2753 = v2741 - v2749
	goto L495
L497:
	;
	v2775 = int32(1)
	if base.Ui32(v2775) < base.Ui32(v2737) {
		v2737 = int32(base.Ui32(v2737) >> (uint(v2775) % 32))
		goto L493
	} else {
		goto L500
	}
L498:
	;
	v2757 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2758 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2757 << (uint(v2758) % 32)
	v2761 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2761 + v2758
	v2765 = m.G1
	v2769 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2765+int32(_a_F_VP8EncWrite_0)+v2753))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2769
	if v2761 < int32(0) {
		goto L497
	} else {
		goto L499
	}
L499:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L497
L500:
	;
	goto L494
L501:
	;
	v2931 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3412))
	v2936 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2938 = v2936 >> (uint(int32(1)) % 32)
	if v2931 == int32(0) {
		v2947 = v2938
		goto L528
	} else {
		goto L529
	}
L502:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2801
	if int32(126) < v2801 {
		goto L504
	} else {
		goto L505
	}
L503:
	;
	v2795 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2797 = v2792 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2795 + v2797
	v2801 = v2790 - v2797
	goto L502
L504:
	;
	if v2785 == int32(0) {
		goto L507
	} else {
		goto L508
	}
L505:
	;
	v2805 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2806 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2805 << (uint(v2806) % 32)
	v2809 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2809 + v2806
	v2813 = m.G1
	v2817 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2813+int32(_a_F_VP8EncWrite_0)+v2801))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2817
	if v2809 < int32(0) {
		goto L504
	} else {
		goto L506
	}
L506:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L504
L507:
	;
	goto L501
L508:
	;
	if int32(-1) < v2785 {
		goto L509
	} else {
		goto L510
	}
L509:
	;
	v2883 = int32(16)
	goto L519
L510:
	;
	v2827 = int32(1)
	v2835 = int32(16)
	goto L511
L511:
	;
	v2839 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2841 = v2839 >> (uint(int32(1)) % 32)
	if v2835&(v2827-v2785<<(uint(v2827)%32)) == int32(0) {
		v2851 = v2841
		goto L513
	} else {
		goto L514
	}
L513:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2851
	if int32(126) < v2851 {
		goto L515
	} else {
		goto L516
	}
L514:
	;
	v2845 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2847 = v2841 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2845 + v2847
	v2851 = v2839 - v2847
	goto L513
L515:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2835) {
		v2835 = int32(base.Ui32(v2835) >> (uint(int32(1)) % 32))
		goto L511
	} else {
		goto L518
	}
L516:
	;
	v2855 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2856 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2855 << (uint(v2856) % 32)
	v2859 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2859 + v2856
	v2863 = m.G1
	v2867 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2863+int32(_a_F_VP8EncWrite_0)+v2851))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2867
	if v2859 < int32(0) {
		goto L515
	} else {
		goto L517
	}
L517:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L515
L518:
	;
	goto L507
L519:
	;
	v2887 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2889 = v2887 >> (uint(int32(1)) % 32)
	if v2883&(v2785<<(uint(int32(1))%32)) == int32(0) {
		v2899 = v2889
		goto L521
	} else {
		goto L522
	}
L520:
	;
	goto L507
L521:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2899
	if int32(126) < v2899 {
		goto L523
	} else {
		goto L524
	}
L522:
	;
	v2893 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2895 = v2889 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2893 + v2895
	v2899 = v2887 - v2895
	goto L521
L523:
	;
	v2921 = int32(1)
	if base.Ui32(v2921) < base.Ui32(v2883) {
		v2883 = int32(base.Ui32(v2883) >> (uint(v2921) % 32))
		goto L519
	} else {
		goto L526
	}
L524:
	;
	v2903 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2904 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2903 << (uint(v2904) % 32)
	v2907 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2907 + v2904
	v2911 = m.G1
	v2915 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2911+int32(_a_F_VP8EncWrite_0)+v2899))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2915
	if v2907 < int32(0) {
		goto L523
	} else {
		goto L525
	}
L525:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L523
L526:
	;
	goto L520
L527:
	;
	v3080 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v3082 = v3080 >> (uint(int32(1)) % 32)
	goto L554
L528:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2947
	if int32(126) < v2947 {
		goto L530
	} else {
		goto L531
	}
L529:
	;
	v2941 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2943 = v2938 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2941 + v2943
	v2947 = v2936 - v2943
	goto L528
L530:
	;
	if v2931 == int32(0) {
		goto L533
	} else {
		goto L534
	}
L531:
	;
	v2951 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2952 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2951 << (uint(v2952) % 32)
	v2955 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v2955 + v2952
	v2959 = m.G1
	v2963 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2959+int32(_a_F_VP8EncWrite_0)+v2947))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2963
	if v2955 < int32(0) {
		goto L530
	} else {
		goto L532
	}
L532:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L530
L533:
	;
	goto L527
L534:
	;
	if int32(-1) < v2931 {
		goto L535
	} else {
		goto L536
	}
L535:
	;
	v3029 = int32(16)
	goto L545
L536:
	;
	v2973 = int32(1)
	v2981 = int32(16)
	goto L537
L537:
	;
	v2985 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v2987 = v2985 >> (uint(int32(1)) % 32)
	if v2981&(v2973-v2931<<(uint(v2973)%32)) == int32(0) {
		v2997 = v2987
		goto L539
	} else {
		goto L540
	}
L539:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v2997
	if int32(126) < v2997 {
		goto L541
	} else {
		goto L542
	}
L540:
	;
	v2991 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v2993 = v2987 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v2991 + v2993
	v2997 = v2985 - v2993
	goto L539
L541:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2981) {
		v2981 = int32(base.Ui32(v2981) >> (uint(int32(1)) % 32))
		goto L537
	} else {
		goto L544
	}
L542:
	;
	v3001 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v3002 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v3001 << (uint(v3002) % 32)
	v3005 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v3005 + v3002
	v3009 = m.G1
	v3013 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3009+int32(_a_F_VP8EncWrite_0)+v2997))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3013
	if v3005 < int32(0) {
		goto L541
	} else {
		goto L543
	}
L543:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L541
L544:
	;
	goto L533
L545:
	;
	v3033 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v3035 = v3033 >> (uint(int32(1)) % 32)
	if v3029&(v2931<<(uint(int32(1))%32)) == int32(0) {
		v3045 = v3035
		goto L547
	} else {
		goto L548
	}
L546:
	;
	goto L533
L547:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3045
	if int32(126) < v3045 {
		goto L549
	} else {
		goto L550
	}
L548:
	;
	v3039 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v3041 = v3035 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v3039 + v3041
	v3045 = v3033 - v3041
	goto L547
L549:
	;
	v3067 = int32(1)
	if base.Ui32(v3067) < base.Ui32(v3029) {
		v3029 = int32(base.Ui32(v3029) >> (uint(v3067) % 32))
		goto L545
	} else {
		goto L552
	}
L550:
	;
	v3049 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v3050 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v3049 << (uint(v3050) % 32)
	v3053 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v3053 + v3050
	v3057 = m.G1
	v3061 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3057+int32(_a_F_VP8EncWrite_0)+v3045))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3061
	if v3053 < int32(0) {
		goto L549
	} else {
		goto L551
	}
L551:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L549
L552:
	;
	goto L546
L553:
	;
	F_VP8WriteProbas(m, v35, v159)
	mBase = m.M
	v3115 = l0 + int32(76)
	v3116 = *(*int32)(unsafe.Add(mBase, uint32(v3115)))
	v3117 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v3118 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	F_VP8CodeIntraModes(m, l0)
	mBase = m.M
	v3125 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	v3129 = int32(1) << (uint(int32(8)-v3125) % 32)
	goto L560
L554:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3082
	if int32(126) < v3082 {
		goto L556
	} else {
		goto L557
	}
L556:
	;
	goto L553
L557:
	;
	v3095 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v3096 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v3095 << (uint(v3096) % 32)
	v3099 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v3099 + v3096
	v3103 = m.G1
	v3107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3103+int32(_a_F_VP8EncWrite_0)+v3082))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3107
	if v3099 < int32(0) {
		goto L556
	} else {
		goto L558
	}
L558:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L556
L559:
	;
	v3164 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v3165 = *(*int32)(unsafe.Add(mBase, uint32(v3164)+88))
	if v3165 == int32(0) {
		goto L566
	} else {
		goto L567
	}
L560:
	;
	v3132 = *(*int32)(unsafe.Add(mBase, uint32(v35)))
	v3134 = v3132 >> (uint(int32(1)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3134
	if int32(126) < v3134 {
		goto L562
	} else {
		goto L563
	}
L561:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = int32(0)
	F_Flush(m, v35)
	mBase = m.M
	goto L559
L562:
	;
	v3156 = int32(1)
	if base.Ui32(v3156) < base.Ui32(v3129) {
		v3129 = int32(base.Ui32(v3129) >> (uint(v3156) % 32))
		goto L560
	} else {
		goto L565
	}
L563:
	;
	v3138 = *(*int32)(unsafe.Add(mBase, uint32(v35)+4))
	v3139 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v35)+4)) = v3138 << (uint(v3139) % 32)
	v3142 = *(*int32)(unsafe.Add(mBase, uint32(v35)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v35)+12)) = v3142 + v3139
	v3146 = m.G1
	v3150 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3146+int32(_a_F_VP8EncWrite_0)+v3134))))
	*(*int32)(unsafe.Add(mBase, uint32(v35))) = v3150
	if v3142 < int32(0) {
		goto L562
	} else {
		goto L564
	}
L564:
	;
	F_Flush(m, v35)
	mBase = m.M
	goto L562
L565:
	;
	goto L561
L566:
	;
	v3212 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	if v3212 == int32(0) {
		goto L3
	} else {
		goto L568
	}
L567:
	;
	v3168 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v3169 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	v3174 = int32(3)
	v3182 = int32(8)
	v3185 = base.I64_extend_i32_u((v3117+v3116)<<(uint(v3174)%32)) + base.I64_extend_i32_s(v3118+v3182)
	v3187 = int64(7)
	v3189 = int64(3)
	v3190 = int64(base.Ui64(base.I64_extend_i32_s(int32(-8)-v32)-base.I64_extend_i32_u((v31+v30)<<(uint(v3174)%32))+v3185+v3187) >> (uint(v3189) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v3165)+36)) = uint32(v3190)
	v3196 = *(*int32)(unsafe.Add(mBase, uint32(v3115)))
	v3205 = int64(base.Ui64(base.I64_extend_i32_s(v3169+v3182)-v3185+base.I64_extend_i32_u((v3168+v3196)<<(uint(v3174)%32))+v3187) >> (uint(v3189) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v3165)+40)) = uint32(v3205)
	v3207 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	*(*int32)(unsafe.Add(mBase, uint32(v3165)+140)) = v3207
	goto L566
L568:
	;
	v3216 = *(*int32)(unsafe.Add(mBase, uint32(v3164)+92))
	if v3216 != 0 {
		goto L570
	} else {
		goto L571
	}
L569:
	;
	v3219 = int32(0)
	goto L4
L570:
	;
	goto L569
L571:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3164)+92)) = int32(1)
	goto L570
L572:
	;
	goto L3
L573:
	;
	v3383 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(372))))
	if v3383 != 0 {
		goto L584
	} else {
		goto L585
	}
L574:
	;
	if base.Ui32(v3242) < base.Ui32(int32(4)) {
		v3315 = int32(0)
		v3321 = v3247
		goto L575
	} else {
		goto L576
	}
L575:
	;
	v3339 = v3242 - v3315
	v3344 = v3315<<(uint(int32(5))%32) + l0 + int32(108)
	v3346 = v3321
	goto L581
L576:
	;
	v3259 = v3242 & int32(2147483644)
	v3269 = l0 + int32(204)
	v3271 = v3259
	v3278 = base.Simd_g_i32x4_replace_lane_l0(base.Simd_g_const(&F_VP8EncWrite__k0), v3247)
	goto L577
L577:
	;
	v3285 = int32(0)
	v3286 = base.Simd_g_v128_load32_splat(m, v3269+int32(-96), v3285)
	v3289 = base.Simd_g_v128_load32_lane_l1(m, v3269+int32(-64), v3285, v3286)
	v3292 = base.Simd_g_v128_load32_lane_l2(m, v3269+int32(-32), v3285, v3289)
	v3295 = base.Simd_g_v128_load32_lane_l3(m, v3269, v3285, v3292)
	v3296 = base.Simd_g_i32x4_add(v3295, v3278)
	v3300 = v3271 + int32(-4)
	if v3300 != 0 {
		v3269 = v3269 + int32(128)
		v3271 = v3300
		v3278 = v3296
		goto L577
	} else {
		goto L579
	}
L578:
	;
	v3303 = base.Simd_g_i32x4_add(v3296, base.Simd_g_i8x16_swizzle_c(v3296, base.Simd_g_const(&F_VP8EncWrite__k2)))
	v3308 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_add(v3303, base.Simd_g_i8x16_swizzle_c(v3303, base.Simd_g_const(&F_VP8EncWrite__k3))))
	if v3242 == v3259 {
		v3371 = v3308
		goto L573
	} else {
		goto L580
	}
L579:
	;
	goto L578
L580:
	;
	v3315 = v3259
	v3321 = v3308
	goto L575
L581:
	;
	v3354 = *(*int32)(unsafe.Add(mBase, uint32(v3344)))
	v3355 = v3354 + v3346
	v3359 = v3339 + int32(-1)
	if v3359 != 0 {
		v3339 = v3359
		v3344 = v3344 + int32(32)
		v3346 = v3355
		goto L581
	} else {
		goto L583
	}
L582:
	;
	v3371 = v3355
	goto L573
L583:
	;
	goto L582
L584:
	;
	v3384 = int32(30)
	goto L586
L585:
	;
	v3384 = int32(12)
	goto L586
L586:
	;
	v3386 = v3371 & int32(1)
	v3387 = v3386 + v3371
	v3388 = v3384 + v3387
	if v3383 == int32(0) {
		v3398 = v3388
		goto L587
	} else {
		goto L588
	}
L587:
	;
	if v3398 != int32(-1) {
		goto L589
	} else {
		goto L590
	}
L588:
	;
	v3391 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	v3398 = v3388 + v3391 + v3391&int32(1) + int32(8)
	goto L587
L589:
	;
	v3408 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(72))))
	v3409 = m.G1
	v3411 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v3420 = *(*int32)(unsafe.Add(mBase, uint32(v3409)+uint32(_c_F_VP8EncWrite[0])))
	*(*int32)(unsafe.Add(mBase, uint32(v22+int32(24)))) = v3420
	v3422 = *(*int64)(unsafe.Add(mBase, uint32(v3409)+uint32(_c_F_VP8EncWrite[1])))
	*(*int64)(unsafe.Add(mBase, uint32(v22)+16)) = v3422
	*(*int32)(unsafe.Add(mBase, uint32(v22)+20)) = v3398
	v3428 = *(*int32)(unsafe.Add(mBase, uint32(v3411)+72))
	v3429 = m.T0[v3428].(func(*base.Module, int32, int32, int32) int32)(m, v22+int32(16), int32(12), v3411)
	mBase = m.M
	if v3429 == int32(0) {
		v3570 = int32(8)
		goto L596
	} else {
		goto L597
	}
L590:
	;
	v3403 = *(*int32)(unsafe.Add(mBase, uint32(v33)+92))
	if v3403 != 0 {
		goto L592
	} else {
		goto L593
	}
L591:
	;
	v3841 = int32(0)
	goto L1
L592:
	;
	goto L591
L593:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+92)) = int32(9)
	goto L592
L594:
	;
	v3690 = l0 + int32(368)
	if v35 == int32(0) {
		goto L632
	} else {
		goto L633
	}
L595:
	;
	v3580 = *(*int32)(unsafe.Add(mBase, uint32(v33)+72))
	v3581 = m.T0[v3580].(func(*base.Module, int32, int32, int32) int32)(m, v3408, v3241, v33)
	mBase = m.M
	if v3581 != 0 {
		goto L613
	} else {
		goto L614
	}
L596:
	;
	v3572 = *(*int32)(unsafe.Add(mBase, uint32(v3411)+92))
	if v3572 != 0 {
		goto L610
	} else {
		goto L611
	}
L597:
	;
	v3434 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(372))))
	if v3434 == int32(0) {
		goto L598
	} else {
		goto L599
	}
L598:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v22)+16)) = int64(540561494)
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+20)) = uint8(v3387)
	v3520 = int32(base.Ui32(v3387) >> (uint(int32(24)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+23)) = uint8(v3520)
	v3522 = int32(16)
	v3523 = int32(base.Ui32(v3387) >> (uint(v3522) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+22)) = uint8(v3523)
	v3525 = int32(8)
	v3527 = int32(base.Ui32(v3387) >> (uint(v3525) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+21)) = uint8(v3527)
	v3532 = *(*int32)(unsafe.Add(mBase, uint32(v3411)+72))
	v3533 = m.T0[v3532].(func(*base.Module, int32, int32, int32) int32)(m, v22+v3522, v3525, v3411)
	mBase = m.M
	if v3533 == int32(0) {
		v3570 = v3525
		goto L596
	} else {
		goto L606
	}
L599:
	;
	v3437 = m.G1
	v3438 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v3441 = int32(0)
	v3442 = base.Simd_g_v128_load(m, v3437+int32(_a_F_VP8EncWrite_1), v3441)
	v3443 = int32(16)
	base.Simd_g_v128_store(m, v22, v3443, v3442)
	*(*int64)(unsafe.Add(mBase, uint32(v22)+20)) = int64(68719476746)
	v3447 = *(*int32)(unsafe.Add(mBase, uint32(v3438)+8))
	v3448 = int32(-1)
	v3449 = v3447 + v3448
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+28)) = uint8(v3449)
	v3452 = int32(base.Ui32(v3449) >> (uint(v3443) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+30)) = uint8(v3452)
	v3454 = int32(8)
	v3456 = int32(base.Ui32(v3449) >> (uint(v3454) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+29)) = uint8(v3456)
	v3458 = *(*int32)(unsafe.Add(mBase, uint32(v3438)+12))
	v3460 = v3458 + v3448
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+31)) = uint8(v3460)
	v3463 = int32(base.Ui32(v3460) >> (uint(v3443) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+33)) = uint8(v3463)
	v3466 = int32(base.Ui32(v3460) >> (uint(v3454) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+32)) = uint8(v3466)
	v3471 = *(*int32)(unsafe.Add(mBase, uint32(v3438)+72))
	v3472 = m.T0[v3471].(func(*base.Module, int32, int32, int32) int32)(m, v22+v3443, int32(18), v3438)
	mBase = m.M
	if v3472 == v3441 {
		v3570 = v3454
		goto L596
	} else {
		goto L600
	}
L600:
	;
	v3477 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(372))))
	if v3477 == int32(0) {
		goto L598
	} else {
		goto L601
	}
L601:
	;
	v3480 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	*(*int64)(unsafe.Add(mBase, uint32(v22)+16)) = int64(1213221953)
	v3483 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	*(*int32)(unsafe.Add(mBase, uint32(v22)+20)) = v3483
	v3485 = int32(8)
	v3489 = *(*int32)(unsafe.Add(mBase, uint32(v3480)+72))
	v3490 = m.T0[v3489].(func(*base.Module, int32, int32, int32) int32)(m, v22+int32(16), v3485, v3480)
	mBase = m.M
	if v3490 == int32(0) {
		v3570 = v3485
		goto L596
	} else {
		goto L602
	}
L602:
	;
	v3493 = *(*int32)(unsafe.Add(mBase, uint32(l0)+376))
	v3494 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	v3495 = *(*int32)(unsafe.Add(mBase, uint32(v3480)+72))
	v3496 = m.T0[v3495].(func(*base.Module, int32, int32, int32) int32)(m, v3493, v3494, v3480)
	mBase = m.M
	if v3496 == int32(0) {
		v3570 = v3485
		goto L596
	} else {
		goto L603
	}
L603:
	;
	v3499 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+380)))
	if v3499&int32(1) == int32(0) {
		goto L598
	} else {
		goto L604
	}
L604:
	;
	v3504 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+15)) = uint8(v3504)
	v3509 = *(*int32)(unsafe.Add(mBase, uint32(v3480)+72))
	v3510 = m.T0[v3509].(func(*base.Module, int32, int32, int32) int32)(m, v22+int32(15), int32(1), v3480)
	mBase = m.M
	if v3510 == v3504 {
		v3570 = v3485
		goto L596
	} else {
		goto L605
	}
L605:
	;
	goto L598
L606:
	;
	if base.Ui32(int32(524287)) < base.Ui32(v3241) {
		v3570 = int32(6)
		goto L596
	} else {
		goto L607
	}
L607:
	;
	v3539 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	v3540 = int32(42)
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+21)) = uint8(v3540)
	v3542 = int32(413)
	*(*uint16)(unsafe.Add(mBase, uint32(v22)+19)) = uint16(v3542)
	v3544 = *(*int32)(unsafe.Add(mBase, uint32(v3411)+8))
	v3545 = *(*int32)(unsafe.Add(mBase, uint32(v3411)+12))
	*(*uint16)(unsafe.Add(mBase, uint32(v22)+24)) = uint16(v3545)
	*(*uint16)(unsafe.Add(mBase, uint32(v22)+22)) = uint16(v3544)
	v3552 = v3539<<(uint(int32(1))%32) | v3241<<(uint(int32(5))%32)
	v3553 = int32(16)
	v3554 = int32(base.Ui32(v3552) >> (uint(v3553) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+18)) = uint8(v3554)
	v3556 = int32(8)
	v3558 = int32(base.Ui32(v3552) >> (uint(v3556) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+17)) = uint8(v3558)
	v3561 = v3552 | v3553
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+16)) = uint8(v3561)
	v3566 = *(*int32)(unsafe.Add(mBase, uint32(v3411)+72))
	v3567 = m.T0[v3566].(func(*base.Module, int32, int32, int32) int32)(m, v22+v3553, int32(10), v3411)
	mBase = m.M
	if v3567 != 0 {
		goto L595
	} else {
		goto L608
	}
L608:
	;
	v3570 = v3556
	goto L596
L609:
	;
	goto L612
L610:
	;
	goto L609
L611:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3411)+92)) = v3570
	goto L610
L612:
	;
	v3681 = int32(0)
	goto L594
L613:
	;
	v3583 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if int32(2) <= v3583 {
		goto L616
	} else {
		goto L617
	}
L614:
	;
	v3681 = int32(0)
	goto L594
L615:
	;
	v3681 = base.B2i32(v3658 != int32(0))
	goto L594
L616:
	;
	v3592 = v3583 + int32(-1)
	v3597 = l0 + int32(108)
	v3598 = v3592
	v3602 = v22 + int32(16)
	goto L618
L617:
	;
	v3658 = int32(1)
	goto L615
L618:
	;
	v3612 = *(*int32)(unsafe.Add(mBase, uint32(v3597)))
	if base.Ui32(v3612) < base.Ui32(int32(16777216)) {
		goto L620
	} else {
		goto L621
	}
L619:
	;
	v3640 = *(*int32)(unsafe.Add(mBase, uint32(v33)+72))
	v3641 = m.T0[v3640].(func(*base.Module, int32, int32, int32) int32)(m, v22+int32(16), v3592*int32(3), v33)
	mBase = m.M
	if v3641 == int32(0) {
		goto L626
	} else {
		goto L627
	}
L620:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3602))) = uint8(v3612)
	v3623 = int32(base.Ui32(v3612) >> (uint(int32(16)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v3602+int32(2)))) = uint8(v3623)
	v3628 = int32(base.Ui32(v3612) >> (uint(int32(8)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v3602+int32(1)))) = uint8(v3628)
	v3635 = v3598 + int32(-1)
	if v3635 != 0 {
		v3597 = v3597 + int32(32)
		v3598 = v3635
		v3602 = v3602 + int32(3)
		goto L618
	} else {
		goto L625
	}
L621:
	;
	v3616 = *(*int32)(unsafe.Add(mBase, uint32(v33)+92))
	if v3616 != 0 {
		goto L623
	} else {
		goto L624
	}
L622:
	;
	v3658 = int32(0)
	goto L615
L623:
	;
	goto L622
L624:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+92)) = int32(7)
	goto L623
L625:
	;
	goto L619
L626:
	;
	v3646 = *(*int32)(unsafe.Add(mBase, uint32(v33)+92))
	if v3646 != 0 {
		goto L629
	} else {
		goto L630
	}
L627:
	;
	v3658 = int32(1)
	goto L615
L628:
	;
	v3658 = int32(0)
	goto L615
L629:
	;
	goto L628
L630:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+92)) = int32(8)
	goto L629
L631:
	;
	v3706 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v3706 < int32(1) {
		v3784 = v3681
		goto L634
	} else {
		goto L635
	}
L632:
	;
	goto L631
L633:
	;
	v3696 = l0 + int32(72)
	v3697 = *(*int32)(unsafe.Add(mBase, uint32(v3696)))
	F_WebPSafeFree(m, v3697)
	mBase = m.M
	v3699 = base.Simd_g_const(&F_VP8EncWrite__k0)
	v3700 = int32(0)
	base.Simd_g_v128_store(m, v3696, v3700, v3699)
	base.Simd_g_v128_store(m, v35, v3700, v3699)
	goto L632
L634:
	;
	if v3784 == int32(0) {
		v3805 = v3784
		goto L649
	} else {
		goto L650
	}
L635:
	;
	v3716 = int32(0)
	v3721 = l0 + int32(108)
	v3723 = v3681
	goto L636
L636:
	;
	v3731 = *(*int32)(unsafe.Add(mBase, uint32(v3721)))
	if v3731 == int32(0) {
		v3742 = v3723
		goto L638
	} else {
		goto L639
	}
L637:
	;
	v3784 = v3766
	goto L634
L638:
	;
	v3744 = v3721 + int32(-20)
	if v3744 == int32(0) {
		goto L643
	} else {
		goto L644
	}
L639:
	;
	if v3723 != 0 {
		goto L640
	} else {
		goto L641
	}
L640:
	;
	v3737 = *(*int32)(unsafe.Add(mBase, uint32(v3721+int32(-4))))
	v3738 = *(*int32)(unsafe.Add(mBase, uint32(v33)+72))
	v3739 = m.T0[v3738].(func(*base.Module, int32, int32, int32) int32)(m, v3737, v3731, v33)
	mBase = m.M
	v3742 = base.B2i32(v3739 != int32(0))
	goto L638
L641:
	;
	v3742 = int32(0)
	goto L638
L642:
	;
	if v3742 != 0 {
		goto L646
	} else {
		goto L647
	}
L643:
	;
	goto L642
L644:
	;
	v3750 = v3721 + int32(-4)
	v3751 = *(*int32)(unsafe.Add(mBase, uint32(v3750)))
	F_WebPSafeFree(m, v3751)
	mBase = m.M
	v3753 = base.Simd_g_const(&F_VP8EncWrite__k0)
	v3754 = int32(0)
	base.Simd_g_v128_store(m, v3750, v3754, v3753)
	base.Simd_g_v128_store(m, v3744, v3754, v3753)
	goto L643
L645:
	;
	v3770 = v3716 + int32(1)
	v3771 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v3770 < v3771 {
		v3716 = v3770
		v3721 = v3721 + int32(32)
		v3723 = v3766
		goto L636
	} else {
		goto L648
	}
L646:
	;
	v3761 = *(*int32)(unsafe.Add(mBase, uint32(v3690)))
	v3763 = F_WebPReportProgress(m, v33, v3761+v27, v3690)
	mBase = m.M
	v3766 = base.B2i32(v3763 != int32(0))
	goto L645
L647:
	;
	v3766 = int32(0)
	goto L645
L648:
	;
	goto L637
L649:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncWrite[2]))) = v3398 + int32(8)
	if v3805 == int32(0) {
		goto L652
	} else {
		goto L653
	}
L650:
	;
	if v3386 == int32(0) {
		v3805 = v3784
		goto L649
	} else {
		goto L651
	}
L651:
	;
	v3796 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v22)+16)) = uint8(v3796)
	v3801 = *(*int32)(unsafe.Add(mBase, uint32(v33)+72))
	v3802 = m.T0[v3801].(func(*base.Module, int32, int32, int32) int32)(m, v22+int32(16), int32(1), v33)
	mBase = m.M
	v3805 = base.B2i32(v3802 != v3796)
	goto L649
L652:
	;
	v3818 = *(*int32)(unsafe.Add(mBase, uint32(v33)+92))
	if v3818 != 0 {
		goto L656
	} else {
		goto L657
	}
L653:
	;
	v3813 = F_WebPReportProgress(m, v33, v24+int32(19), v3690)
	mBase = m.M
	if v3813 == int32(0) {
		goto L652
	} else {
		goto L654
	}
L654:
	;
	v3841 = int32(1)
	goto L1
L655:
	;
	goto L2
L656:
	;
	goto L655
L657:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v33)+92)) = int32(8)
	goto L656
}

var F_VP8EncWrite__k0 = [2]uint64{0x0, 0x0}
var F_VP8EncWrite__k1 = [2]uint64{0xfe, 0xfffffff800000000}
var F_VP8EncWrite__k2 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_VP8EncWrite__k3 = [2]uint64{0x302010007060504, 0x302010003020100}

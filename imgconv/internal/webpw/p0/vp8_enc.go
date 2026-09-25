//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_VP8EncAnalyze(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v52 int32
	_ = v52
	var v55 int32
	_ = v55
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v69 int32
	_ = v69
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v79 int32
	_ = v79
	var v96 int32
	_ = v96
	var v98 int32
	_ = v98
	var v100 int64
	_ = v100
	var v102 int32
	_ = v102
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v109 int32
	_ = v109
	var v114 int32
	_ = v114
	var v115 int64
	_ = v115
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v129 int32
	_ = v129
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v145 int32
	_ = v145
	var v146 int32
	_ = v146
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v156 int32
	_ = v156
	var v159 int32
	_ = v159
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v175 int64
	_ = v175
	var v177 int32
	_ = v177
	var v178 int32
	_ = v178
	var v180 int32
	_ = v180
	var v184 int32
	_ = v184
	var v185 int32
	_ = v185
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v210 int32
	_ = v210
	var v211 int32
	_ = v211
	var v214 int32
	_ = v214
	var v218 int32
	_ = v218
	var v219 int64
	_ = v219
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v231 int32
	_ = v231
	var v233 int32
	_ = v233
	var v238 int32
	_ = v238
	var v239 int32
	_ = v239
	var v243 int32
	_ = v243
	var v369 int32
	_ = v369
	var v370 int32
	_ = v370
	var v374 int32
	_ = v374
	var v375 int32
	_ = v375
	var v378 int32
	_ = v378
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
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v398 int32
	_ = v398
	var v406 int32
	_ = v406
	var v412 int32
	_ = v412
	var v424 int32
	_ = v424
	var v427 int32
	_ = v427
	var v430 int32
	_ = v430
	var v431 int32
	_ = v431
	var v435 int32
	_ = v435
	var v446 int32
	_ = v446
	var v457 int32
	_ = v457
	var v464 int32
	_ = v464
	var v473 int32
	_ = v473
	var v477 int32
	_ = v477
	var v490 int32
	_ = v490
	var v504 int32
	_ = v504
	var v510 int32
	_ = v510
	var v516 int32
	_ = v516
	var v527 int32
	_ = v527
	var v529 int32
	_ = v529
	var v533 int32
	_ = v533
	var v535 int32
	_ = v535
	var v537 int32
	_ = v537
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v564 int32
	_ = v564
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v570 int32
	_ = v570
	var v571 int32
	_ = v571
	var v572 int32
	_ = v572
	var v575 int32
	_ = v575
	var v577 int32
	_ = v577
	var v585 int32
	_ = v585
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v596 int32
	_ = v596
	var v602 int32
	_ = v602
	var v603 int32
	_ = v603
	var v610 int32
	_ = v610
	var v631 int32
	_ = v631
	var v632 int32
	_ = v632
	var v633 int32
	_ = v633
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v662 int32
	_ = v662
	var v663 int32
	_ = v663
	var v666 int32
	_ = v666
	var v667 int32
	_ = v667
	var v669 int32
	_ = v669
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v699 int64
	_ = v699
	var v703 int32
	_ = v703
	var v704 int32
	_ = v704
	var v709 int32
	_ = v709
	var v723 int32
	_ = v723
	var v724 int32
	_ = v724
	var v725 int32
	_ = v725
	var v733 int32
	_ = v733
	var v737 int32
	_ = v737
	var v746 int32
	_ = v746
	var v747 int32
	_ = v747
	var v751 int32
	_ = v751
	var v758 int32
	_ = v758
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v778 int32
	_ = v778
	var v782 int32
	_ = v782
	var v789 int32
	_ = v789
	var v794 int32
	_ = v794
	var v805 int32
	_ = v805
	var v819 int32
	_ = v819
	var v822 int32
	_ = v822
	var v823 int32
	_ = v823
	var v825 int32
	_ = v825
	var v846 int32
	_ = v846
	var v848 int32
	_ = v848
	var v849 int32
	_ = v849
	var v860 int32
	_ = v860
	var v882 int32
	_ = v882
	var v885 int32
	_ = v885
	var v886 int32
	_ = v886
	var v894 int32
	_ = v894
	var v895 int32
	_ = v895
	var v930 int32
	_ = v930
	var v931 int32
	_ = v931
	var v936 int64
	_ = v936
	var v939 int32
	_ = v939
	var v940 int32
	_ = v940
	var v952 int32
	_ = v952
	var v971 int32
	_ = v971
	var v972 int32
	_ = v972
	var v983 int32
	_ = v983
	var v1005 int32
	_ = v1005
	var v1008 int32
	_ = v1008
	var v1009 int32
	_ = v1009
	var v1017 int32
	_ = v1017
	var v1018 int32
	_ = v1018
	var v1053 int32
	_ = v1053
	var v1054 int32
	_ = v1054
	var v1059 int64
	_ = v1059
	var v1062 int32
	_ = v1062
	var v1063 int32
	_ = v1063
	var v1075 int32
	_ = v1075
	var v1097 int32
	_ = v1097
	var v1098 int32
	_ = v1098
	var v1115 int32
	_ = v1115
	var v1117 int32
	_ = v1117
	var v1121 int32
	_ = v1121
	var v1123 int32
	_ = v1123
	var v1125 int32
	_ = v1125
	var v1135 int32
	_ = v1135
	var v1136 int32
	_ = v1136
	var v1152 int32
	_ = v1152
	var v1156 int32
	_ = v1156
	var v1157 int32
	_ = v1157
	var v1158 int32
	_ = v1158
	var v1159 int32
	_ = v1159
	var v1160 int32
	_ = v1160
	var v1163 int32
	_ = v1163
	var v1165 int32
	_ = v1165
	var v1172 int32
	_ = v1172
	var v1182 int32
	_ = v1182
	var v1183 int32
	_ = v1183
	var v1184 int32
	_ = v1184
	var v1190 int32
	_ = v1190
	var v1191 int32
	_ = v1191
	var v1198 int32
	_ = v1198
	var v1237 int32
	_ = v1237
	var v1248 int32
	_ = v1248
	var v1249 int32
	_ = v1249
	var v1251 int32
	_ = v1251
	var v1253 int32
	_ = v1253
	var v1257 int32
	_ = v1257
	var v1259 int32
	_ = v1259
	var v1262 int32
	_ = v1262
	var v1266 int32
	_ = v1266
	var v1269 int32
	_ = v1269
	var v1270 int32
	_ = v1270
	var v1272 int32
	_ = v1272
	var v1274 int32
	_ = v1274
	var v1276 int32
	_ = v1276
	var v1278 int32
	_ = v1278
	var v1285 int32
	_ = v1285
	var v1287 int32
	_ = v1287
	var v1289 int32
	_ = v1289
	var v1290 int32
	_ = v1290
	var v1297 int32
	_ = v1297
	var v1301 int32
	_ = v1301
	var v1306 int32
	_ = v1306
	var v1308 int32
	_ = v1308
	var v1323 int32
	_ = v1323
	var v1329 int32
	_ = v1329
	var v1330 int32
	_ = v1330
	var v1331 int32
	_ = v1331
	var v1334 int32
	_ = v1334
	var v1339 int32
	_ = v1339
	var v1340 int32
	_ = v1340
	var v1356 int32
	_ = v1356
	var v1357 int32
	_ = v1357
	var v1360 int32
	_ = v1360
	var v1361 int32
	_ = v1361
	var v1362 int32
	_ = v1362
	var v1363 int32
	_ = v1363
	var v1366 int32
	_ = v1366
	var v1371 int32
	_ = v1371
	var v1374 int32
	_ = v1374
	var v1381 int32
	_ = v1381
	var v1386 int32
	_ = v1386
	var v1387 int32
	_ = v1387
	var v1388 int32
	_ = v1388
	var v1389 int32
	_ = v1389
	var v1393 int32
	_ = v1393
	var v1396 int32
	_ = v1396
	var v1400 int32
	_ = v1400
	var v1413 int32
	_ = v1413
	var v1414 int32
	_ = v1414
	var v1419 int64
	_ = v1419
	var v1420 int32
	_ = v1420
	var v1427 int64
	_ = v1427
	var v1428 int32
	_ = v1428
	var v1429 int64
	_ = v1429
	var v1439 int32
	_ = v1439
	var v1441 int32
	_ = v1441
	var v1450 int32
	_ = v1450
	var v1452 int32
	_ = v1452
	var v1457 int32
	_ = v1457
	var v1458 int32
	_ = v1458
	var v1468 int32
	_ = v1468
	var v1472 int32
	_ = v1472
	var v1474 int32
	_ = v1474
	var v1486 int32
	_ = v1486
	var v1488 int32
	_ = v1488
	var v1503 int64
	_ = v1503
	var v1508 int32
	_ = v1508
	var v1509 int32
	_ = v1509
	var v1510 int32
	_ = v1510
	var v1511 int32
	_ = v1511
	var v1512 int32
	_ = v1512
	var v1514 int32
	_ = v1514
	var v1516 int32
	_ = v1516
	var v1517 int32
	_ = v1517
	var v1518 int32
	_ = v1518
	var v1523 int32
	_ = v1523
	var v1525 int32
	_ = v1525
	var v1530 int32
	_ = v1530
	var v1531 int32
	_ = v1531
	var v1537 int32
	_ = v1537
	var v1539 int32
	_ = v1539
	var v1544 int32
	_ = v1544
	var v1545 int32
	_ = v1545
	var v1551 int32
	_ = v1551
	var v1552 int32
	_ = v1552
	var v1557 int32
	_ = v1557
	var v1558 int32
	_ = v1558
	var v1566 int32
	_ = v1566
	var v1571 int32
	_ = v1571
	var v1572 int32
	_ = v1572
	var v1578 int32
	_ = v1578
	var v1581 int32
	_ = v1581
	var v1582 int32
	_ = v1582
	var v1587 int32
	_ = v1587
	var v1588 int32
	_ = v1588
	var v1596 int32
	_ = v1596
	var v1601 int32
	_ = v1601
	var v1602 int32
	_ = v1602
	var v1610 int32
	_ = v1610
	var v1615 int32
	_ = v1615
	var v1616 int32
	_ = v1616
	var v1620 int32
	_ = v1620
	var v1624 int32
	_ = v1624
	var v1628 int32
	_ = v1628
	var v1632 int32
	_ = v1632
	var v1637 int32
	_ = v1637
	var v1640 int32
	_ = v1640
	var v1641 int32
	_ = v1641
	var v1647 int32
	_ = v1647
	var v1652 int32
	_ = v1652
	var v1654 int32
	_ = v1654
	var v1659 int32
	_ = v1659
	var v1676 int32
	_ = v1676
	var v1680 int32
	_ = v1680
	var v1682 int32
	_ = v1682
	var v1695 int32
	_ = v1695
	var v1696 int32
	_ = v1696
	var v1711 int32
	_ = v1711
	var v1714 int32
	_ = v1714
	var v1715 int32
	_ = v1715
	var v1716 int32
	_ = v1716
	var v1718 int32
	_ = v1718
	var v1719 int32
	_ = v1719
	var v1720 int32
	_ = v1720
	var v1722 int32
	_ = v1722
	var v1724 int32
	_ = v1724
	var v1726 int32
	_ = v1726
	var v1727 int32
	_ = v1727
	var v1728 int32
	_ = v1728
	var v1733 int32
	_ = v1733
	var v1738 int32
	_ = v1738
	var v1743 int32
	_ = v1743
	var v1751 int32
	_ = v1751
	var v1769 int32
	_ = v1769
	var v1771 int32
	_ = v1771
	var v1774 int32
	_ = v1774
	var v1775 int32
	_ = v1775
	var v1779 int32
	_ = v1779
	var v1784 int32
	_ = v1784
	var v1791 int32
	_ = v1791
	var v1834 int32
	_ = v1834
	var v1835 int32
	_ = v1835
	var v1839 int32
	_ = v1839
	var v1851 int32
	_ = v1851
	var v1852 int32
	_ = v1852
	var v1853 int32
	_ = v1853
	var v1857 int32
	_ = v1857
	var v1868 int32
	_ = v1868
	var v1870 int32
	_ = v1870
	var v1873 int32
	_ = v1873
	var v1875 int32
	_ = v1875
	var v1878 int32
	_ = v1878
	var v1880 int32
	_ = v1880
	var v1883 int32
	_ = v1883
	var v1885 int32
	_ = v1885
	var v1887 int32
	_ = v1887
	var v1889 int32
	_ = v1889
	var v1891 int32
	_ = v1891
	var v1893 int32
	_ = v1893
	var v1897 int32
	_ = v1897
	var v1903 int32
	_ = v1903
	var v1904 int32
	_ = v1904
	var v1908 int32
	_ = v1908
	var v1930 int32
	_ = v1930
	var v1931 int32
	_ = v1931
	var v1937 int32
	_ = v1937
	var v1942 int32
	_ = v1942
	var v1946 int32
	_ = v1946
	var v1948 int32
	_ = v1948
	var v1950 int32
	_ = v1950
	var v1954 int32
	_ = v1954
	var v1959 int32
	_ = v1959
	var v1960 int32
	_ = v1960
	var v1975 int32
	_ = v1975
	var v1981 int32
	_ = v1981
	var v1982 int32
	_ = v1982
	var v1989 int32
	_ = v1989
	var v1990 int32
	_ = v1990
	var v1991 int32
	_ = v1991
	var v2007 int32
	_ = v2007
	var v2009 int32
	_ = v2009
	var v2011 int32
	_ = v2011
	var v2015 int32
	_ = v2015
	var v2016 int32
	_ = v2016
	var v2019 int32
	_ = v2019
	var v2026 int32
	_ = v2026
	var v2027 int32
	_ = v2027
	var v2030 int32
	_ = v2030
	var v2031 int32
	_ = v2031
	var v2034 int32
	_ = v2034
	var v2041 int32
	_ = v2041
	var v2049 int32
	_ = v2049
	v2 = int32(0)
	v21 = m.G0
	v23 = v21 - int32(_a_F_VP8EncAnalyze_0)
	m.G0 = v23
	v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v26 = *(*int32)(unsafe.Add(mBase, uint32(v25)+80))
	if v26 != 0 {
		goto L6
	} else {
		goto L7
	}
L1:
	;
	m.G0 = v23 + int32(_a_F_VP8EncAnalyze_0)
	return v2049
L2:
	;
	v1329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v1330 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v1331 = v1329 * v1330
	if v1331 < int32(1) {
		v1393 = v1329
		v1396 = v1331
		v1400 = v1330
		goto L153
	} else {
		goto L154
	}
L3:
	;
	v733 = int32(1)
	v737 = v398 << (uint(v733) % 32)
	if v394 == v733 {
		v794 = int32(0)
		v805 = v733
		goto L89
	} else {
		goto L90
	}
L4:
	;
	v723 = int32(256)
	v724 = int32(255)
	v725 = int32(-1)
	goto L3
L5:
	;
	v631 = int32(1)
	v632 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v633 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	if v632*v633 < v631 {
		goto L84
	} else {
		goto L85
	}
L6:
	;
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v35 = m.G1
	goto L10
L7:
	;
	v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if int32(1) < v27 {
		goto L6
	} else {
		goto L8
	}
L8:
	;
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[0])))
	if int32(1) < v30 {
		goto L5
	} else {
		goto L9
	}
L9:
	;
	goto L6
L10:
	;
	v40 = m.G1
	goto L11
L11:
	;
	v43 = *(*int32)(unsafe.Add(mBase, uint32(v40)+uint32(_c_F_VP8EncAnalyze[1])))
	m.T0[v43].(func(*base.Module, int32))(m, v23+int32(8))
	mBase = m.M
	v46 = v23 + int32(1064)
	*(*int32)(unsafe.Add(mBase, uint32(v23)+24)) = v46
	v48 = m.G2
	*(*int32)(unsafe.Add(mBase, uint32(v23)+16)) = v48 + int32(130)
	v52 = int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v23)+20)) = v23 + v52
	v55 = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(v46))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+24)) = l0
	v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+280)) = v60
	v62 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+296)) = v62
	v64 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+304)) = v64
	*(*int32)(unsafe.Add(mBase, uint32(v46)+32)) = l0 + int32(88)
	v69 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[4])))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+36)) = v69
	v73 = int32(-32)
	v74 = (v23 + int32(1511)) & v73
	*(*int32)(unsafe.Add(mBase, uint32(v46)+8)) = v74
	v79 = (v23 + int32(1424)) & v73
	*(*int32)(unsafe.Add(mBase, uint32(v46)+308)) = v79
	*(*int32)(unsafe.Add(mBase, uint32(v46)+20)) = v74 + int32(1536)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+16)) = v74 + int32(1024)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+12)) = v74 + int32(512)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+316)) = v79 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+312)) = v79 + int32(32)
	v96 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+40)) = v96
	v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+28)) = v98
	v100 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v46)+320)) = v100
	v102 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v79)+47)) = uint8(v102)
	v104 = *(*int32)(unsafe.Add(mBase, uint32(v46)+312))
	v105 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v104+v105))) = uint8(v102)
	v109 = *(*int32)(unsafe.Add(mBase, uint32(v46)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v109+v105))) = uint8(v102)
	v114 = *(*int32)(unsafe.Add(mBase, uint32(v46)+308))
	v115 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v114))) = v115
	*(*int64)(unsafe.Add(mBase, uint32(v114+v52))) = v115
	v121 = *(*int32)(unsafe.Add(mBase, uint32(v46)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v121))) = v115
	v124 = *(*int32)(unsafe.Add(mBase, uint32(v46)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v124))) = v115
	*(*int32)(unsafe.Add(mBase, uint32(v46)+160)) = v55
	v129 = *(*int32)(unsafe.Add(mBase, uint32(v46)+304))
	if v129 == v55 {
		goto L13
	} else {
		goto L14
	}
L12:
	;
	v167 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+4)) = v167
	*(*int32)(unsafe.Add(mBase, uint32(v46))) = v167
	v172 = *(*int32)(unsafe.Add(mBase, uint32(v46)+24))
	v173 = *(*int32)(unsafe.Add(mBase, uint32(v172)+uint32(_c_F_VP8EncAnalyze[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+40)) = v173
	v175 = *(*int64)(unsafe.Add(mBase, uint32(v172)+uint32(_c_F_VP8EncAnalyze[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v46)+320)) = v175
	v177 = *(*int32)(unsafe.Add(mBase, uint32(v172)+uint32(_c_F_VP8EncAnalyze[4])))
	v178 = *(*int32)(unsafe.Add(mBase, uint32(v172)+48))
	v180 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+36)) = v177 + v167*v178<<(uint(v180)%32)
	v184 = *(*int32)(unsafe.Add(mBase, uint32(v172)+uint32(_c_F_VP8EncAnalyze[6])))
	v185 = *(*int32)(unsafe.Add(mBase, uint32(v172)+40))
	*(*int32)(unsafe.Add(mBase, uint32(v46)+28)) = v184 + v185*v167<<(uint(v180)%32)
	v191 = *(*int32)(unsafe.Add(mBase, uint32(v172)+52))
	v192 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+32)) = v172 + (v191+v192)&v167<<(uint(int32(5))%32) + int32(88)
	v201 = *(*int32)(unsafe.Add(mBase, uint32(v46)+316))
	v205 = int32(127)
	goto L19
L13:
	;
	v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v135 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v136 = v134 * v135
	*(*int32)(unsafe.Add(mBase, uint32(v46)+292)) = v136
	*(*int32)(unsafe.Add(mBase, uint32(v46)+288)) = v136
	v139 = *(*int32)(unsafe.Add(mBase, uint32(v46)+24))
	v140 = *(*int32)(unsafe.Add(mBase, uint32(v139)+uint32(_c_F_VP8EncAnalyze[7])))
	v142 = *(*int32)(unsafe.Add(mBase, uint32(v139)+40))
	v145 = F_memset(m, v140, int32(127), v142<<(uint(int32(5))%32))
	mBase = m.M
	v146 = *(*int32)(unsafe.Add(mBase, uint32(v139)+uint32(_c_F_VP8EncAnalyze[5])))
	v147 = int32(0)
	v148 = *(*int32)(unsafe.Add(mBase, uint32(v139)+40))
	v151 = F_memset(m, v146, v147, v148<<(uint(int32(2))%32))
	mBase = m.M
	v152 = *(*int32)(unsafe.Add(mBase, uint32(v139)+uint32(_c_F_VP8EncAnalyze[3])))
	if v152 == v147 {
		goto L15
	} else {
		goto L16
	}
L14:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v46)+300)) = int32(0)
	goto L13
L15:
	;
	v162 = int32(0)
	v164 = F_memset(m, v23+int32(1232), v162, int32(96))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v46)+284)) = v162
	goto L12
L16:
	;
	v156 = *(*int32)(unsafe.Add(mBase, uint32(v139)+40))
	v159 = F_memset(m, v152, int32(0), v156<<(uint(int32(2))%32))
	mBase = m.M
	goto L15
L17:
	;
	v238 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v239 = v34 * v238
	*(*int32)(unsafe.Add(mBase, uint32(v46)+288)) = v239
	*(*int32)(unsafe.Add(mBase, uint32(v46)+292)) = v239
	goto L23
L19:
	;
	goto L20
L20:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v201+v192))) = uint8(v205)
	v210 = *(*int32)(unsafe.Add(mBase, uint32(v46)+312))
	v211 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v210+v211))) = uint8(v205)
	v214 = *(*int32)(unsafe.Add(mBase, uint32(v46)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v214+v211))) = uint8(v205)
	v218 = *(*int32)(unsafe.Add(mBase, uint32(v46)+308))
	v219 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v218))) = v219
	*(*int64)(unsafe.Add(mBase, uint32(v218+int32(8)))) = v219
	v225 = *(*int32)(unsafe.Add(mBase, uint32(v46)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v225))) = v219
	v228 = *(*int32)(unsafe.Add(mBase, uint32(v46)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v228))) = v219
	v231 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v46)+160)) = v231
	v233 = *(*int32)(unsafe.Add(mBase, uint32(v46)+304))
	if v233 == v231 {
		goto L21
	} else {
		goto L22
	}
L21:
	;
	goto L17
L22:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v46)+300)) = int32(0)
	goto L21
L23:
	;
	v243 = v23 + int32(32)
	goto L26
L24:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[8]))) = int32(20)
	v369 = v23 + int32(8)
	v370 = *(*int32)(unsafe.Add(mBase, uint32(v35)+uint32(_c_F_VP8EncAnalyze[9])))
	m.T0[v370].(func(*base.Module, int32))(m, v369)
	mBase = m.M
	v374 = *(*int32)(unsafe.Add(mBase, uint32(v35)+uint32(_c_F_VP8EncAnalyze[10])))
	v375 = m.T0[v374].(func(*base.Module, int32) int32)(m, v369)
	mBase = m.M
	v378 = *(*int32)(unsafe.Add(mBase, uint32(v35)+uint32(_c_F_VP8EncAnalyze[11])))
	m.T0[v378].(func(*base.Module, int32))(m, v369)
	mBase = m.M
	if v375&int32(1) != 0 {
		goto L38
	} else {
		goto L39
	}
L26:
	;
	base.MemoryFill(m, v243, int32(0), int32(1032))
	goto L24
L38:
	;
	v387 = *(*int32)(unsafe.Add(mBase, uint32(v23)+1056))
	v388 = v33 * v34
	v389 = base.I32_div_s(v387, v388)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3388)) = v389
	v391 = *(*int32)(unsafe.Add(mBase, uint32(v23)+1060))
	v392 = base.I32_div_s(v391, v388)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+3392)) = v392
	v394 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v395 = int32(4)
	if v394 < v395 {
		goto L43
	} else {
		goto L44
	}
L39:
	;
	v382 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v384 = *(*int32)(unsafe.Add(mBase, uint32(v382)+92))
	if v384 != 0 {
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v2049 = int32(0)
	goto L1
L41:
	;
	goto L40
L42:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v382)+92)) = int32(1)
	goto L41
L43:
	;
	v398 = v394
	goto L45
L44:
	;
	v398 = v395
	goto L45
L45:
	;
	v406 = v23 + int32(44)
	v412 = int32(3)
	goto L50
L46:
	;
	if base.Ui32(int32(254)) < base.Ui32(v446) {
		v490 = int32(255)
		goto L58
	} else {
		goto L59
	}
L47:
	;
	v446 = v412 + int32(-1)
	goto L46
L48:
	;
	v446 = v412 + int32(-2)
	goto L46
L49:
	;
	v446 = v412 + int32(-3)
	goto L46
L50:
	;
	v424 = *(*int32)(unsafe.Add(mBase, uint32(v406+int32(-12))))
	if v424 != 0 {
		goto L49
	} else {
		goto L52
	}
L51:
	;
	if int32(0) < v394 {
		goto L4
	} else {
		goto L57
	}
L52:
	;
	v427 = *(*int32)(unsafe.Add(mBase, uint32(v406+int32(-8))))
	if v427 != 0 {
		goto L48
	} else {
		goto L53
	}
L53:
	;
	v430 = *(*int32)(unsafe.Add(mBase, uint32(v406+int32(-4))))
	if v430 != 0 {
		goto L47
	} else {
		goto L54
	}
L54:
	;
	v431 = *(*int32)(unsafe.Add(mBase, uint32(v406)))
	if v431 != 0 {
		v446 = v412
		goto L46
	} else {
		goto L55
	}
L55:
	;
	v435 = v412 + int32(4)
	if v435 != int32(259) {
		v406 = v406 + int32(16)
		v412 = v435
		goto L50
	} else {
		goto L56
	}
L56:
	;
	goto L51
L57:
	;
	v1323 = v2
	goto L2
L58:
	;
	if v394 <= int32(0) {
		goto L64
	} else {
		goto L65
	}
L59:
	;
	v457 = v23 + int32(1052)
	v464 = int32(255)
	goto L60
L60:
	;
	v473 = *(*int32)(unsafe.Add(mBase, uint32(v457)))
	if v473 != 0 {
		v490 = v464
		goto L58
	} else {
		goto L62
	}
L61:
	;
	v490 = v446
	goto L58
L62:
	;
	v477 = v464 + int32(-1)
	if base.Ui32(v446) < base.Ui32(v477) {
		v457 = v457 + int32(-4)
		v464 = v477
		goto L60
	} else {
		goto L63
	}
L63:
	;
	goto L61
L64:
	;
	if v490 < v446 {
		v1323 = v2
		goto L2
	} else {
		goto L66
	}
L65:
	;
	v723 = v446
	v724 = v490
	v725 = v490 - v446
	goto L3
L66:
	;
	if v490 < v446 {
		goto L67
	} else {
		goto L68
	}
L67:
	;
	v504 = v446
	goto L69
L68:
	;
	v504 = v490
	goto L69
L69:
	;
	v510 = int32(0)
	v516 = v446
	goto L70
L70:
	;
	v527 = v516 << (uint(int32(2)) % 32)
	v529 = *(*int32)(unsafe.Add(mBase, uint32(v243+v527)))
	if v529 == int32(0) {
		v610 = v510
		goto L72
	} else {
		goto L73
	}
L72:
	;
	if base.B2i32(v516 == v504) == int32(0) {
		v510 = v610
		v516 = v516 + int32(1)
		goto L70
	} else {
		goto L83
	}
L73:
	;
	v533 = v510 + int32(1)
	if v533 < v398 {
		goto L74
	} else {
		goto L75
	}
L74:
	;
	v535 = v398
	goto L76
L75:
	;
	v535 = v533
	goto L76
L76:
	;
	v537 = v535 + int32(-1)
	v547 = v510
	v548 = v23 + int32(_a_F_VP8EncAnalyze_1) + v510<<(uint(int32(2))%32)
	goto L78
L77:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v23+int32(_a_F_VP8EncAnalyze_2)+v527))) = v585
	v594 = v585 << (uint(int32(2)) % 32)
	v595 = v23 + int32(_a_F_VP8EncAnalyze_3) + v594
	v596 = *(*int32)(unsafe.Add(mBase, uint32(v595)))
	*(*int32)(unsafe.Add(mBase, uint32(v595))) = v596 + v529*v516
	v602 = v23 + int32(_a_F_VP8EncAnalyze_4) + v594
	v603 = *(*int32)(unsafe.Add(mBase, uint32(v602)))
	*(*int32)(unsafe.Add(mBase, uint32(v602))) = v603 + v529
	v610 = v585
	goto L72
L78:
	;
	if v537 != v547 {
		goto L80
	} else {
		goto L81
	}
L79:
	;
	v585 = v547 + int32(0)
	goto L77
L80:
	;
	v564 = *(*int32)(unsafe.Add(mBase, uint32(v548)))
	v568 = v548 + int32(4)
	v569 = *(*int32)(unsafe.Add(mBase, uint32(v568)))
	v570 = v516 - v569
	v571 = int32(31)
	v572 = v570 >> (uint(v571) % 32)
	v575 = v516 - v564
	v577 = v575 >> (uint(v571) % 32)
	if base.Ui32(v570^v572-v572) < base.Ui32(v575^v577-v577) {
		v547 = v547 + int32(1)
		v548 = v568
		goto L78
	} else {
		goto L82
	}
L81:
	;
	v585 = v537
	goto L77
L82:
	;
	goto L79
L83:
	;
	v1323 = v2
	goto L2
L84:
	;
	v699 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+3388)) = v699
	*(*int64)(unsafe.Add(mBase, uint32(l0)+1080)) = v699
	v703 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v704 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	v709 = F_WebPReportProgress(m, v703, v704+int32(20), l0+int32(368))
	mBase = m.M
	v2049 = v631
	goto L1
L85:
	;
	v643 = int32(1)
	v644 = int32(0)
	goto L86
L86:
	;
	v659 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v660 = v659 + v643
	v662 = v660 + int32(-1)
	v663 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v662))))
	v666 = int32(1)
	v667 = v663&int32(128) | v666
	*(*uint8)(unsafe.Add(mBase, uint32(v662))) = uint8(v667)
	v669 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v660))) = uint8(v669)
	v674 = v644 + v666
	v675 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v676 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	if v674 < v675*v676 {
		v643 = v643 + int32(4)
		v644 = v674
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
	if v398&v733 == int32(0) {
		goto L94
	} else {
		goto L95
	}
L90:
	;
	v746 = v398 & int32(-2)
	v747 = int32(1)
	v751 = int32(0)
	v758 = v751
	v759 = v751
	v760 = v23 + int32(_a_F_VP8EncAnalyze_1)
	goto L91
L91:
	;
	v778 = base.I32_div_s(v725*int32(3)+v759, v737)
	*(*int32)(unsafe.Add(mBase, uint32(v760+int32(4)))) = v778 + v723
	v782 = base.I32_div_s(v725+v759, v737)
	*(*int32)(unsafe.Add(mBase, uint32(v760))) = v782 + v723
	v789 = v758 + int32(2)
	if v746 != v789 {
		v758 = v789
		v759 = v759 + v725<<(uint(int32(2))%32)
		v760 = v760 + int32(8)
		goto L91
	} else {
		goto L93
	}
L92:
	;
	v794 = v789
	v805 = v746<<(uint(v747)%32) | v747
	goto L89
L93:
	;
	goto L92
L94:
	;
	v822 = base.B2i32(v724 < v723)
	if v724 < v723 {
		goto L96
	} else {
		goto L97
	}
L95:
	;
	v819 = base.I32_div_s(v805*v725, v737)
	*(*int32)(unsafe.Add(mBase, uint32(v23+int32(_a_F_VP8EncAnalyze_1)+v794<<(uint(int32(2))%32)))) = v819 + v723
	goto L94
L96:
	;
	v823 = v723
	goto L98
L97:
	;
	v823 = v724
	goto L98
L98:
	;
	v825 = v398 << (uint(int32(2)) % 32)
	v846 = int32(0)
	goto L99
L99:
	;
	v848 = v23 + int32(_a_F_VP8EncAnalyze_4)
	v849 = int32(0)
	if base.Ui32(v825) < base.Ui32(int32(33)) {
		goto L102
	} else {
		goto L103
	}
L100:
	;
	v1306 = base.I32_div_s(v1285, int32(2))
	v1308 = base.I32_div_s(v1289+v1306, v1285)
	v1323 = v1308
	goto L2
L101:
	;
	v971 = v23 + int32(_a_F_VP8EncAnalyze_3)
	v972 = int32(0)
	if base.Ui32(v825) < base.Ui32(int32(33)) {
		goto L116
	} else {
		goto L117
	}
L102:
	;
	if v825 == int32(0) {
		goto L104
	} else {
		goto L105
	}
L103:
	;
	base.MemoryFill(m, v848, v849, v825)
	goto L101
L104:
	;
	goto L101
L105:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[12]))) = uint8(v849)
	v860 = v848 + v825
	*(*uint8)(unsafe.Add(mBase, uint32(v860+int32(-1)))) = uint8(v849)
	if base.Ui32(v825) < base.Ui32(int32(3)) {
		goto L104
	} else {
		goto L106
	}
L106:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[13]))) = uint8(v849)
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[14]))) = uint8(v849)
	*(*uint8)(unsafe.Add(mBase, uint32(v860+int32(-3)))) = uint8(v849)
	*(*uint8)(unsafe.Add(mBase, uint32(v860+int32(-2)))) = uint8(v849)
	if base.Ui32(v825) < base.Ui32(int32(7)) {
		goto L104
	} else {
		goto L107
	}
L107:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[15]))) = uint8(v849)
	*(*uint8)(unsafe.Add(mBase, uint32(v860+int32(-4)))) = uint8(v849)
	if base.Ui32(v825) < base.Ui32(int32(9)) {
		goto L104
	} else {
		goto L108
	}
L108:
	;
	v882 = int32(0)
	v885 = (v882 - v848) & int32(3)
	v886 = v848 + v885
	*(*int32)(unsafe.Add(mBase, uint32(v886))) = v882
	v894 = (v825 - v885) & int32(60)
	v895 = v886 + v894
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-4)))) = v882
	if base.Ui32(v894) < base.Ui32(int32(9)) {
		goto L104
	} else {
		goto L109
	}
L109:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v886)+8)) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v886)+4)) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-8)))) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-12)))) = v882
	if base.Ui32(v894) < base.Ui32(int32(25)) {
		goto L104
	} else {
		goto L110
	}
L110:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v886)+24)) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v886)+20)) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v886)+16)) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v886)+12)) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-16)))) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-20)))) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-24)))) = v882
	*(*int32)(unsafe.Add(mBase, uint32(v895+int32(-28)))) = v882
	v930 = v886&int32(4) | int32(24)
	v931 = v894 - v930
	if base.Ui32(v931) < base.Ui32(int32(32)) {
		goto L104
	} else {
		goto L111
	}
L111:
	;
	v936 = base.I64_extend_i32_u(v882) * int64(4294967297)
	v939 = v931
	v940 = v886 + v930
	goto L112
L112:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v940)+24)) = v936
	*(*int64)(unsafe.Add(mBase, uint32(v940)+16)) = v936
	*(*int64)(unsafe.Add(mBase, uint32(v940)+8)) = v936
	*(*int64)(unsafe.Add(mBase, uint32(v940))) = v936
	v952 = v939 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v952) {
		v939 = v952
		v940 = v940 + int32(32)
		goto L112
	} else {
		goto L114
	}
L113:
	;
	goto L104
L114:
	;
	goto L113
L115:
	;
	if v724 < v723 {
		goto L129
	} else {
		goto L130
	}
L116:
	;
	if v825 == int32(0) {
		goto L118
	} else {
		goto L119
	}
L117:
	;
	base.MemoryFill(m, v971, v972, v825)
	goto L115
L118:
	;
	goto L115
L119:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[16]))) = uint8(v972)
	v983 = v971 + v825
	*(*uint8)(unsafe.Add(mBase, uint32(v983+int32(-1)))) = uint8(v972)
	if base.Ui32(v825) < base.Ui32(int32(3)) {
		goto L118
	} else {
		goto L120
	}
L120:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[17]))) = uint8(v972)
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[18]))) = uint8(v972)
	*(*uint8)(unsafe.Add(mBase, uint32(v983+int32(-3)))) = uint8(v972)
	*(*uint8)(unsafe.Add(mBase, uint32(v983+int32(-2)))) = uint8(v972)
	if base.Ui32(v825) < base.Ui32(int32(7)) {
		goto L118
	} else {
		goto L121
	}
L121:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[19]))) = uint8(v972)
	*(*uint8)(unsafe.Add(mBase, uint32(v983+int32(-4)))) = uint8(v972)
	if base.Ui32(v825) < base.Ui32(int32(9)) {
		goto L118
	} else {
		goto L122
	}
L122:
	;
	v1005 = int32(0)
	v1008 = (v1005 - v971) & int32(3)
	v1009 = v971 + v1008
	*(*int32)(unsafe.Add(mBase, uint32(v1009))) = v1005
	v1017 = (v825 - v1008) & int32(60)
	v1018 = v1009 + v1017
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-4)))) = v1005
	if base.Ui32(v1017) < base.Ui32(int32(9)) {
		goto L118
	} else {
		goto L123
	}
L123:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1009)+8)) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1009)+4)) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-8)))) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-12)))) = v1005
	if base.Ui32(v1017) < base.Ui32(int32(25)) {
		goto L118
	} else {
		goto L124
	}
L124:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1009)+24)) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1009)+20)) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1009)+16)) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1009)+12)) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-16)))) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-20)))) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-24)))) = v1005
	*(*int32)(unsafe.Add(mBase, uint32(v1018+int32(-28)))) = v1005
	v1053 = v1009&int32(4) | int32(24)
	v1054 = v1017 - v1053
	if base.Ui32(v1054) < base.Ui32(int32(32)) {
		goto L118
	} else {
		goto L125
	}
L125:
	;
	v1059 = base.I64_extend_i32_u(v1005) * int64(4294967297)
	v1062 = v1054
	v1063 = v1009 + v1053
	goto L126
L126:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v1063)+24)) = v1059
	*(*int64)(unsafe.Add(mBase, uint32(v1063)+16)) = v1059
	*(*int64)(unsafe.Add(mBase, uint32(v1063)+8)) = v1059
	*(*int64)(unsafe.Add(mBase, uint32(v1063))) = v1059
	v1075 = v1062 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v1075) {
		v1062 = v1075
		v1063 = v1063 + int32(32)
		goto L126
	} else {
		goto L128
	}
L127:
	;
	goto L118
L128:
	;
	goto L127
L129:
	;
	v1237 = int32(0)
	v1248 = v23 + int32(_a_F_VP8EncAnalyze_3)
	v1249 = v23 + int32(_a_F_VP8EncAnalyze_4)
	v1251 = v23 + int32(_a_F_VP8EncAnalyze_1)
	v1253 = v1237
	v1257 = v398
	v1259 = v1237
	v1262 = v1237
	goto L145
L130:
	;
	v1097 = v723
	v1098 = int32(0)
	goto L131
L131:
	;
	v1115 = v1097 << (uint(int32(2)) % 32)
	v1117 = *(*int32)(unsafe.Add(mBase, uint32(v243+v1115)))
	if v1117 == int32(0) {
		v1198 = v1098
		goto L133
	} else {
		goto L134
	}
L132:
	;
	goto L129
L133:
	;
	if v1097 != v823 {
		v1097 = v1097 + int32(1)
		v1098 = v1198
		goto L131
	} else {
		goto L144
	}
L134:
	;
	v1121 = v1098 + int32(1)
	if v1121 < v398 {
		goto L135
	} else {
		goto L136
	}
L135:
	;
	v1123 = v398
	goto L137
L136:
	;
	v1123 = v1121
	goto L137
L137:
	;
	v1125 = v1123 + int32(-1)
	v1135 = v1098
	v1136 = v23 + int32(_a_F_VP8EncAnalyze_1) + v1098<<(uint(int32(2))%32)
	goto L139
L138:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v23+int32(_a_F_VP8EncAnalyze_2)+v1115))) = v1172
	v1182 = v1172 << (uint(int32(2)) % 32)
	v1183 = v23 + int32(_a_F_VP8EncAnalyze_3) + v1182
	v1184 = *(*int32)(unsafe.Add(mBase, uint32(v1183)))
	*(*int32)(unsafe.Add(mBase, uint32(v1183))) = v1184 + v1117*v1097
	v1190 = v23 + int32(_a_F_VP8EncAnalyze_4) + v1182
	v1191 = *(*int32)(unsafe.Add(mBase, uint32(v1190)))
	*(*int32)(unsafe.Add(mBase, uint32(v1190))) = v1191 + v1117
	v1198 = v1172
	goto L133
L139:
	;
	if v1125 != v1135 {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v1172 = v1135 + int32(0)
	goto L138
L141:
	;
	v1152 = *(*int32)(unsafe.Add(mBase, uint32(v1136)))
	v1156 = v1136 + int32(4)
	v1157 = *(*int32)(unsafe.Add(mBase, uint32(v1156)))
	v1158 = v1097 - v1157
	v1159 = int32(31)
	v1160 = v1158 >> (uint(v1159) % 32)
	v1163 = v1097 - v1152
	v1165 = v1163 >> (uint(v1159) % 32)
	if base.Ui32(v1158^v1160-v1160) < base.Ui32(v1163^v1165-v1165) {
		v1135 = v1135 + int32(1)
		v1136 = v1156
		goto L139
	} else {
		goto L143
	}
L142:
	;
	v1172 = v1125
	goto L138
L143:
	;
	goto L140
L144:
	;
	goto L132
L145:
	;
	v1266 = *(*int32)(unsafe.Add(mBase, uint32(v1249)))
	if v1266 == int32(0) {
		v1285 = v1253
		v1287 = v1259
		v1289 = v1262
		goto L147
	} else {
		goto L148
	}
L146:
	;
	if v1287 < int32(5) {
		goto L150
	} else {
		goto L151
	}
L147:
	;
	v1290 = int32(4)
	v1297 = v1257 + int32(-1)
	if v1297 != 0 {
		v1248 = v1248 + v1290
		v1249 = v1249 + v1290
		v1251 = v1251 + v1290
		v1253 = v1285
		v1257 = v1297
		v1259 = v1287
		v1262 = v1289
		goto L145
	} else {
		goto L149
	}
L148:
	;
	v1269 = *(*int32)(unsafe.Add(mBase, uint32(v1251)))
	v1270 = *(*int32)(unsafe.Add(mBase, uint32(v1248)))
	v1272 = base.I32_div_s(v1266, int32(2))
	v1274 = base.I32_div_s(v1270+v1272, v1266)
	*(*int32)(unsafe.Add(mBase, uint32(v1251))) = v1274
	v1276 = v1269 - v1274
	v1278 = v1276 >> (uint(int32(31)) % 32)
	v1285 = v1266 + v1253
	v1287 = v1276 ^ v1278 - v1278 + v1259
	v1289 = v1274*v1266 + v1262
	goto L147
L149:
	;
	goto L146
L150:
	;
	goto L100
L151:
	;
	v1301 = v846 + int32(1)
	if v1301 != int32(6) {
		v846 = v1301
		goto L99
	} else {
		goto L152
	}
L152:
	;
	goto L150
L153:
	;
	if v394 < int32(2) {
		goto L158
	} else {
		goto L159
	}
L154:
	;
	v1334 = int32(0)
	v1339 = v1334
	v1340 = v1334
	goto L155
L155:
	;
	v1356 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1357 = v1356 + v1340
	v1360 = int32(1)
	v1361 = v1357 + v1360
	v1362 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1361))))
	v1363 = int32(2)
	v1366 = *(*int32)(unsafe.Add(mBase, uint32(v23+int32(_a_F_VP8EncAnalyze_2)+v1362<<(uint(v1363)%32))))
	v1371 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1357))))
	v1374 = v1366<<(uint(int32(5))%32)&int32(96) | v1371&int32(159)
	*(*uint8)(unsafe.Add(mBase, uint32(v1357))) = uint8(v1374)
	v1381 = *(*int32)(unsafe.Add(mBase, uint32(v23+int32(_a_F_VP8EncAnalyze_1)+v1366<<(uint(v1363)%32))))
	*(*uint8)(unsafe.Add(mBase, uint32(v1361))) = uint8(v1381)
	v1386 = v1339 + v1360
	v1387 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v1388 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v1389 = v1387 * v1388
	if v1386 < v1389 {
		v1339 = v1386
		v1340 = v1340 + int32(4)
		goto L155
	} else {
		goto L157
	}
L156:
	;
	v1393 = v1387
	v1396 = v1389
	v1400 = v1388
	goto L153
L157:
	;
	goto L156
L158:
	;
	v1834 = *(*int32)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[20])))
	v1835 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if int32(2) <= v1835 {
		goto L199
	} else {
		goto L200
	}
L159:
	;
	v1413 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v1414 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1413)+68)))
	if v1414&int32(1) == int32(0) {
		goto L158
	} else {
		goto L160
	}
L160:
	;
	v1419 = base.I64_extend_i32_s(v1396)
	v1420 = int32(1)
	if v1419 == int64(0) {
		goto L163
	} else {
		goto L164
	}
L161:
	;
	if v1441 == int32(0) {
		goto L158
	} else {
		goto L167
	}
L162:
	;
	goto L161
L163:
	;
	v1439 = F_malloc(m, base.I32_wrap_i64(v1419)*v1420)
	mBase = m.M
	v1441 = v1439
	goto L162
L164:
	;
	v1427 = base.I64_div_u_s(int64(2147418112), v1419)
	v1428 = int32(0)
	v1429 = base.I64_extend_i32_u(v1420)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1429*v1419) {
		v1441 = v1428
		goto L162
	} else {
		goto L165
	}
L165:
	;
	if base.Ui64(v1427) < base.Ui64(v1429) {
		v1441 = v1428
		goto L162
	} else {
		goto L166
	}
L166:
	;
	goto L163
L167:
	;
	if v1393 < int32(3) {
		goto L168
	} else {
		goto L169
	}
L168:
	;
	F_free(m, v1441)
	mBase = m.M
	goto L197
L169:
	;
	if v1400 < int32(3) {
		goto L168
	} else {
		goto L170
	}
L170:
	;
	v1450 = v1393 + int32(-1)
	v1452 = v1400 << (uint(int32(2)) % 32)
	v1457 = v1400 + v1441
	v1458 = int32(1)
	v1468 = int32(0)
	v1472 = v1458
	v1474 = v1457 + v1458
	goto L171
L171:
	;
	v1486 = int32(0)
	v1488 = v1468
	goto L173
L172:
	;
	v1654 = int32(1)
	v1659 = v1400 << (uint(int32(2)) % 32)
	v1676 = v1654
	v1680 = v1457 + v1654
	v1682 = v1659 + int32(8)
	goto L187
L173:
	;
	v1503 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[21]))) = v1503
	*(*int64)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[22]))) = v1503
	v1508 = v23 + int32(_a_F_VP8EncAnalyze_5)
	v1509 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1510 = v1509 + v1488
	v1511 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1510))))
	v1512 = int32(3)
	v1514 = int32(12)
	v1516 = v1508 | int32(base.Ui32(v1511)>>(uint(v1512)%32))&v1514
	v1517 = *(*int32)(unsafe.Add(mBase, uint32(v1516)))
	v1518 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v1516))) = v1517 + v1518
	v1523 = int32(4)
	v1525 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1510+v1523))))
	v1530 = v1508 | int32(base.Ui32(v1525)>>(uint(v1512)%32))&v1514
	v1531 = *(*int32)(unsafe.Add(mBase, uint32(v1530)))
	*(*int32)(unsafe.Add(mBase, uint32(v1530))) = v1531 + v1518
	v1537 = int32(8)
	v1539 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1510+v1537))))
	v1544 = v1508 | int32(base.Ui32(v1539)>>(uint(v1512)%32))&v1514
	v1545 = *(*int32)(unsafe.Add(mBase, uint32(v1544)))
	*(*int32)(unsafe.Add(mBase, uint32(v1544))) = v1545 + v1518
	v1551 = v1510 + v1452
	v1552 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1551))))
	v1557 = v1508 | int32(base.Ui32(v1552)>>(uint(v1512)%32))&v1514
	v1558 = *(*int32)(unsafe.Add(mBase, uint32(v1557)))
	*(*int32)(unsafe.Add(mBase, uint32(v1557))) = v1558 + v1518
	v1566 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1551+v1537))))
	v1571 = v1508 | int32(base.Ui32(v1566)>>(uint(v1512)%32))&v1514
	v1572 = *(*int32)(unsafe.Add(mBase, uint32(v1571)))
	*(*int32)(unsafe.Add(mBase, uint32(v1571))) = v1572 + v1518
	v1578 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1551+v1523))))
	v1581 = v1510 + v1400<<(uint(int32(3))%32)
	v1582 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1581))))
	v1587 = v1508 | int32(base.Ui32(v1582)>>(uint(v1512)%32))&v1514
	v1588 = *(*int32)(unsafe.Add(mBase, uint32(v1587)))
	*(*int32)(unsafe.Add(mBase, uint32(v1587))) = v1588 + v1518
	v1596 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1581+v1523))))
	v1601 = v1508 | int32(base.Ui32(v1596)>>(uint(v1512)%32))&v1514
	v1602 = *(*int32)(unsafe.Add(mBase, uint32(v1601)))
	*(*int32)(unsafe.Add(mBase, uint32(v1601))) = v1602 + v1518
	v1610 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1581+v1537))))
	v1615 = v1508 | int32(base.Ui32(v1610)>>(uint(v1512)%32))&v1514
	v1616 = *(*int32)(unsafe.Add(mBase, uint32(v1615)))
	*(*int32)(unsafe.Add(mBase, uint32(v1615))) = v1616 + v1518
	v1620 = *(*int32)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[22])))
	if v1620 <= v1523 {
		goto L176
	} else {
		goto L177
	}
L174:
	;
	v1652 = v1472 + int32(1)
	if v1652 != v1450 {
		v1468 = v1468 + v1452
		v1472 = v1652
		v1474 = v1474 + v1400
		goto L171
	} else {
		goto L186
	}
L175:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1474+v1486))) = uint8(v1641)
	v1647 = v1486 + int32(1)
	if v1400+int32(-2) != v1647 {
		v1486 = v1647
		v1488 = v1488 + int32(4)
		goto L173
	} else {
		goto L185
	}
L176:
	;
	v1624 = *(*int32)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[23])))
	if v1624 <= int32(4) {
		goto L178
	} else {
		goto L179
	}
L177:
	;
	v1641 = int32(0)
	goto L175
L178:
	;
	v1628 = *(*int32)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[21])))
	if v1628 <= int32(4) {
		goto L180
	} else {
		goto L181
	}
L179:
	;
	v1641 = int32(1)
	goto L175
L180:
	;
	v1632 = int32(3)
	v1637 = *(*int32)(unsafe.Add(mBase, uint32(v23)+uint32(_c_F_VP8EncAnalyze[24])))
	if int32(4) < v1637 {
		goto L182
	} else {
		goto L183
	}
L181:
	;
	v1641 = int32(2)
	goto L175
L182:
	;
	v1640 = v1632
	goto L184
L183:
	;
	v1640 = int32(base.Ui32(v1578)>>(uint(int32(5))%32)) & v1632
	goto L184
L184:
	;
	v1641 = v1640
	goto L175
L185:
	;
	goto L174
L186:
	;
	goto L172
L187:
	;
	if v1400 == int32(3) {
		v1751 = int32(1)
		goto L189
	} else {
		goto L190
	}
L188:
	;
	goto L168
L189:
	;
	if v1400&v1654 == int32(0) {
		goto L194
	} else {
		goto L195
	}
L190:
	;
	v1695 = v1682
	v1696 = int32(0)
	goto L191
L191:
	;
	v1711 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1714 = v1711 + v1695 + int32(-4)
	v1715 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1714))))
	v1716 = int32(159)
	v1718 = v1680 + v1696
	v1719 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1718))))
	v1720 = int32(5)
	v1722 = int32(96)
	v1724 = v1715&v1716 | v1719<<(uint(v1720)%32)&v1722
	*(*uint8)(unsafe.Add(mBase, uint32(v1714))) = uint8(v1724)
	v1726 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1727 = v1726 + v1695
	v1728 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1727))))
	v1733 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1718+int32(1)))))
	v1738 = v1728&v1716 | v1733<<(uint(v1720)%32)&v1722
	*(*uint8)(unsafe.Add(mBase, uint32(v1727))) = uint8(v1738)
	v1743 = v1696 + int32(2)
	if v1400&int32(2147483646)+int32(-2) != v1743 {
		v1695 = v1695 + int32(8)
		v1696 = v1743
		goto L191
	} else {
		goto L193
	}
L192:
	;
	v1751 = v1696 + int32(3)
	goto L189
L193:
	;
	goto L192
L194:
	;
	v1791 = v1676 + int32(1)
	if v1791 != v1450 {
		v1676 = v1791
		v1680 = v1680 + v1400
		v1682 = v1682 + v1659
		goto L187
	} else {
		goto L196
	}
L195:
	;
	v1769 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncAnalyze[6])))
	v1771 = v1751 + v1676*v1400
	v1774 = v1769 + v1771<<(uint(int32(2))%32)
	v1775 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1774))))
	v1779 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1441+v1771))))
	v1784 = v1775&int32(159) | v1779<<(uint(int32(5))%32)&int32(96)
	*(*uint8)(unsafe.Add(mBase, uint32(v1774))) = uint8(v1784)
	goto L194
L196:
	;
	goto L188
L197:
	;
	goto L158
L198:
	;
	v1975 = int32(1)
	if v1835 < v1975 {
		v2049 = v1975
		goto L1
	} else {
		goto L241
	}
L199:
	;
	v1839 = v1835 & int32(3)
	if base.Ui32(int32(4)) <= base.Ui32(v1835) {
		goto L202
	} else {
		goto L203
	}
L200:
	;
	v1959 = v1834
	v1960 = v1834
	goto L198
L201:
	;
	if v1839 == int32(0) {
		v1959 = v1903
		v1960 = v1904
		goto L198
	} else {
		goto L231
	}
L202:
	;
	v1851 = v23 + int32(_a_F_VP8EncAnalyze_1)
	v1852 = v1834
	v1853 = v1834
	v1857 = int32(0)
	goto L204
L203:
	;
	v1903 = v1834
	v1904 = v1834
	v1908 = int32(0)
	goto L201
L204:
	;
	v1868 = *(*int32)(unsafe.Add(mBase, uint32(v1851)))
	if v1868 < v1852 {
		goto L206
	} else {
		goto L207
	}
L205:
	;
	v1903 = v1885
	v1904 = v1893
	v1908 = v1897
	goto L201
L206:
	;
	v1870 = v1852
	goto L208
L207:
	;
	v1870 = v1868
	goto L208
L208:
	;
	v1873 = *(*int32)(unsafe.Add(mBase, uint32(v1851+int32(4))))
	if v1873 < v1870 {
		goto L209
	} else {
		goto L210
	}
L209:
	;
	v1875 = v1870
	goto L211
L210:
	;
	v1875 = v1873
	goto L211
L211:
	;
	v1878 = *(*int32)(unsafe.Add(mBase, uint32(v1851+int32(8))))
	if v1878 < v1875 {
		goto L212
	} else {
		goto L213
	}
L212:
	;
	v1880 = v1875
	goto L214
L213:
	;
	v1880 = v1878
	goto L214
L214:
	;
	v1883 = *(*int32)(unsafe.Add(mBase, uint32(v1851+int32(12))))
	if v1883 < v1880 {
		goto L215
	} else {
		goto L216
	}
L215:
	;
	v1885 = v1880
	goto L217
L216:
	;
	v1885 = v1883
	goto L217
L217:
	;
	if v1853 < v1868 {
		goto L218
	} else {
		goto L219
	}
L218:
	;
	v1887 = v1853
	goto L220
L219:
	;
	v1887 = v1868
	goto L220
L220:
	;
	if v1887 < v1873 {
		goto L221
	} else {
		goto L222
	}
L221:
	;
	v1889 = v1887
	goto L223
L222:
	;
	v1889 = v1873
	goto L223
L223:
	;
	if v1889 < v1878 {
		goto L224
	} else {
		goto L225
	}
L224:
	;
	v1891 = v1889
	goto L226
L225:
	;
	v1891 = v1878
	goto L226
L226:
	;
	if v1891 < v1883 {
		goto L227
	} else {
		goto L228
	}
L227:
	;
	v1893 = v1891
	goto L229
L228:
	;
	v1893 = v1883
	goto L229
L229:
	;
	v1897 = v1857 + int32(4)
	if v1835&int32(2147483644) != v1897 {
		v1851 = v1851 + int32(16)
		v1852 = v1885
		v1853 = v1893
		v1857 = v1897
		goto L204
	} else {
		goto L230
	}
L230:
	;
	goto L205
L231:
	;
	v1930 = v1903
	v1931 = v1904
	v1937 = v23 + int32(_a_F_VP8EncAnalyze_1) + v1908<<(uint(int32(2))%32)
	v1942 = v1839
	goto L232
L232:
	;
	v1946 = *(*int32)(unsafe.Add(mBase, uint32(v1937)))
	if v1946 < v1930 {
		goto L234
	} else {
		goto L235
	}
L233:
	;
	v1959 = v1948
	v1960 = v1950
	goto L198
L234:
	;
	v1948 = v1930
	goto L236
L235:
	;
	v1948 = v1946
	goto L236
L236:
	;
	if v1931 < v1946 {
		goto L237
	} else {
		goto L238
	}
L237:
	;
	v1950 = v1931
	goto L239
L238:
	;
	v1950 = v1946
	goto L239
L239:
	;
	v1954 = v1942 + int32(-1)
	if v1954 != 0 {
		v1930 = v1948
		v1931 = v1950
		v1937 = v1937 + int32(4)
		v1942 = v1954
		goto L232
	} else {
		goto L240
	}
L240:
	;
	goto L233
L241:
	;
	if v1959 == v1960 {
		goto L242
	} else {
		goto L243
	}
L242:
	;
	v1981 = v1960 + int32(1)
	goto L244
L243:
	;
	v1981 = v1959
	goto L244
L244:
	;
	v1982 = v1981 - v1960
	v1989 = v1835
	v1990 = v23 + int32(_a_F_VP8EncAnalyze_1)
	v1991 = l0 + int32(1084)
	goto L245
L245:
	;
	v2007 = *(*int32)(unsafe.Add(mBase, uint32(v1990)))
	v2009 = int32(255)
	v2011 = base.I32_div_s((v2007-v1960)*v2009, v1982)
	if v2011 < v2009 {
		goto L247
	} else {
		goto L248
	}
L246:
	;
	v2049 = v1975
	goto L1
L247:
	;
	v2015 = v2011
	goto L249
L248:
	;
	v2015 = v2009
	goto L249
L249:
	;
	v2016 = int32(0)
	if v2016 < v2015 {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v2019 = v2015
	goto L252
L251:
	;
	v2019 = v2016
	goto L252
L252:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1991))) = v2019
	v2026 = base.I32_div_s((v2007-v1323)*int32(255), v1982)
	v2027 = int32(127)
	if v2026 < v2027 {
		goto L253
	} else {
		goto L254
	}
L253:
	;
	v2030 = v2026
	goto L255
L254:
	;
	v2030 = v2027
	goto L255
L255:
	;
	v2031 = int32(-127)
	if v2031 < v2030 {
		goto L256
	} else {
		goto L257
	}
L256:
	;
	v2034 = v2030
	goto L258
L257:
	;
	v2034 = v2031
	goto L258
L258:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1991+int32(-4)))) = v2034
	v2041 = v1989 + int32(-1)
	if v2041 != 0 {
		v1989 = v2041
		v1990 = v1990 + int32(4)
		v1991 = v1991 + int32(744)
		goto L245
	} else {
		goto L259
	}
L259:
	;
	goto L246
}
func F_VP8EncDeleteAlpha(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v18 int32
	_ = v18
	var v20 int32
	_ = v20
	var v22 int32
	_ = v22
	v4 = int32(1)
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncDeleteAlpha[0])))
	if v5 < v4 {
		v20 = v4
	} else {
		v9 = l0 + int32(384)
		v10 = m.G1
		v13 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8EncDeleteAlpha[1])))
		v14 = m.T0[v13].(func(*base.Module, int32) int32)(m, v9)
		mBase = m.M
		v15 = m.G1
		v18 = *(*int32)(unsafe.Add(mBase, uint32(v15)+uint32(_c_F_VP8EncDeleteAlpha[2])))
		m.T0[v18].(func(*base.Module, int32))(m, v9)
		mBase = m.M
		v20 = v14
	}
	v22 = *(*int32)(unsafe.Add(mBase, uint32(l0)+376))
	F_free(m, v22)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(l0)+380)) = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+372)) = int64(0)
	return v20
}
func F_VP8EncDspCostInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v12 int32
	_ = v12
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_VP8EncDspCostInit[0])))
	v8 = m.G3
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8EncDspCostInit[1]))) = v11 + int32(128)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8EncDspCostInit[2]))) = v11 + int32(129)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_VP8EncDspCostInit[0]))) = v9
	}
	return
}
func F_VP8EncDspInit(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v7 int32
	_ = v7
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v11 int32
	_ = v11
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v43 int32
	_ = v43
	var v46 int32
	_ = v46
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v115 int32
	_ = v115
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	v4 = m.G1
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_VP8EncDspInit[0])))
	v8 = m.G3
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G1
		F_VP8DspInit(m)
		mBase = m.M
		v15 = *(*int32)(unsafe.Add(mBase, uint32(v11)+uint32(_c_F_VP8EncDspInit[1])))
		if v15 != 0 {
		} else {
			v17 = int32(-255)
			for {
				v20 = m.G1
				v23 = v20 + int32(_a_F_VP8EncDspInit_0) + v17
				v26 = int32(0)
				if v26 < v17 {
					v29 = v17
				} else {
					v29 = v26
				}
				v30 = int32(255)
				if v29 < v30 {
					v33 = v29
				} else {
					v33 = v30
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v23+int32(255)))) = uint8(v33)
				v37 = int32(-1)
				if v37 < v17 {
					v40 = v17
				} else {
					v40 = v37
				}
				v42 = v40 + int32(1)
				v43 = int32(255)
				if base.Ui32(v42) < base.Ui32(v43) {
					v46 = v42
				} else {
					v46 = v43
				}
				*(*uint8)(unsafe.Add(mBase, uint32(v23+int32(256)))) = uint8(v46)
				v49 = v17 + int32(2)
				if v49 != int32(511) {
					v17 = v49
					continue
				} else {
					break
				}
				break
			}
			v52 = m.G1
			*(*int32)(unsafe.Add(mBase, uint32(v52)+uint32(_c_F_VP8EncDspInit[1]))) = int32(1)
		}
		v60 = m.G2
		v61 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[2]))) = v60 + int32(109)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[3]))) = v60 + int32(110)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[4]))) = v60 + int32(111)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[5]))) = v60 + int32(112)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[6]))) = v60 + int32(113)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[7]))) = v60 + int32(114)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[8]))) = v60 + int32(115)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[9]))) = v60 + int32(116)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[10]))) = v60 + int32(117)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[11]))) = v60 + int32(118)
		v115 = v60 + int32(119)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[12]))) = v115
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[13]))) = v60 + int32(120)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[14]))) = v115
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[15]))) = v60 + int32(121)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[16]))) = v60 + int32(122)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[17]))) = v60 + int32(123)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[18]))) = v60 + int32(124)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[19]))) = v60 + int32(125)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[20]))) = v60 + int32(126)
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[21]))) = v60 + int32(127)
		v162 = m.G3
		v163 = *(*int32)(unsafe.Add(mBase, uint32(v162)))
		*(*int32)(unsafe.Add(mBase, uint32(v61)+uint32(_c_F_VP8EncDspInit[0]))) = v163
	}
	return
}
func F_VP8EncFinishAlpha(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v5 int32
	_ = v5
	var v10 int32
	_ = v10
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v23 int32
	_ = v23
	v2 = *(*int32)(unsafe.Add(mBase, uint32(l0)+372))
	if v2 == int32(0) {
		v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
		v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
		v23 = F_WebPReportProgress(m, v17, v18+int32(20), l0+int32(368))
		mBase = m.M
		return v23
	} else {
		v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncFinishAlpha[0])))
		if v5 < int32(1) {
			v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
			v23 = F_WebPReportProgress(m, v17, v18+int32(20), l0+int32(368))
			mBase = m.M
			return v23
		} else {
			v10 = m.G1
			v13 = *(*int32)(unsafe.Add(mBase, uint32(v10)+uint32(_c_F_VP8EncFinishAlpha[1])))
			v14 = m.T0[v13].(func(*base.Module, int32) int32)(m, l0+int32(384))
			mBase = m.M
			if v14 != 0 {
				v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
				v23 = F_WebPReportProgress(m, v17, v18+int32(20), l0+int32(368))
				mBase = m.M
				return v23
			} else {
				return int32(0)
			}
		}
	}
}
func F_VP8EncFreeBitWriters(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v15 int64
	_ = v15
	var v26 int32
	_ = v26
	var v33 int32
	_ = v33
	var v34 int32
	_ = v34
	var v39 int32
	_ = v39
	var v40 int32
	_ = v40
	var v44 int64
	_ = v44
	var v58 int32
	_ = v58
	var v59 int32
	_ = v59
	v5 = l0 + int32(56)
	if v5 == int32(0) {
	} else {
		v10 = l0 + int32(72)
		v11 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
		F_WebPSafeFree(m, v11)
		mBase = m.M
		v15 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(80)))) = v15
		*(*int64)(unsafe.Add(mBase, uint32(v10))) = v15
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(64)))) = v15
		*(*int64)(unsafe.Add(mBase, uint32(v5))) = v15
	}
	v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v26 < int32(1) {
	} else {
		v33 = l0 + int32(88)
		v34 = int32(0)
		for {
			if v33 == int32(0) {
			} else {
				v39 = v33 + int32(16)
				v40 = *(*int32)(unsafe.Add(mBase, uint32(v39)))
				F_WebPSafeFree(m, v40)
				mBase = m.M
				v44 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(v33+int32(24)))) = v44
				*(*int64)(unsafe.Add(mBase, uint32(v39))) = v44
				*(*int64)(unsafe.Add(mBase, uint32(v33+int32(8)))) = v44
				*(*int64)(unsafe.Add(mBase, uint32(v33))) = v44
			}
			v58 = v34 + int32(1)
			v59 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
			if v58 < v59 {
				v33 = v33 + int32(32)
				v34 = v58
				continue
			} else {
				break
			}
			break
		}
	}
	return
}
func F_VP8EncInitAlpha(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v10 int32
	_ = v10
	var v11 int32
	_ = v11
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v85 int32
	_ = v85
	var v90 int32
	_ = v90
	var v93 int32
	_ = v93
	var v97 int32
	_ = v97
	v6 = m.G1
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v6)+uint32(_c_F_VP8EncInitAlpha[0])))
	v10 = m.G3
	v11 = *(*int32)(unsafe.Add(mBase, uint32(v10)))
	if v9 == v11 {
	} else {
		v13 = m.G2
		v14 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[1]))) = v13 + int32(2)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[2]))) = v13 + int32(3)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[3]))) = v13 + int32(4)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[4]))) = v13 + int32(5)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[5]))) = v13 + int32(6)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[6]))) = v13 + int32(7)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[7]))) = v13 + int32(8)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[8]))) = v13 + int32(9)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[9]))) = v13 + int32(10)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[10]))) = v13 + int32(11)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[11]))) = v13 + int32(12)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[12]))) = v13 + int32(13)
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_VP8EncInitAlpha[0]))) = v11
	}
	v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v81 = F_WebPPictureHasTransparency(m, v80)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(l0)+376)) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(l0)+372)) = v81
	v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncInitAlpha[13])))
	if v85 < int32(1) {
	} else {
		v90 = m.G1
		v93 = *(*int32)(unsafe.Add(mBase, uint32(v90)+uint32(_c_F_VP8EncInitAlpha[14])))
		m.T0[v93].(func(*base.Module, int32))(m, l0+int32(384))
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0)+400)) = int32(0)
		v97 = m.G2
		*(*int32)(unsafe.Add(mBase, uint32(l0)+392)) = v97 + int32(169)
		*(*int32)(unsafe.Add(mBase, uint32(l0)+396)) = l0
	}
	return
}
func F_VP8EncStartAlpha(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v6 int32
	_ = v6
	var v10 int32
	_ = v10
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
	var v25 int32
	_ = v25
	var v29 int32
	_ = v29
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+372))
	if v3 != 0 {
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncStartAlpha[0])))
		if v6 < int32(1) {
			v29 = F_CompressAlphaJob(m, l0, l0)
			mBase = m.M
			return v29
		} else {
			v10 = l0 + int32(384)
			v11 = m.G1
			v14 = *(*int32)(unsafe.Add(mBase, uint32(v11)+uint32(_c_F_VP8EncStartAlpha[1])))
			v15 = m.T0[v14].(func(*base.Module, int32) int32)(m, v10)
			mBase = m.M
			if v15 != 0 {
				v22 = m.G1
				v25 = *(*int32)(unsafe.Add(mBase, uint32(v22)+uint32(_c_F_VP8EncStartAlpha[2])))
				m.T0[v25].(func(*base.Module, int32))(m, v10)
				mBase = m.M
				return int32(1)
			} else {
				v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				v18 = *(*int32)(unsafe.Add(mBase, uint32(v16)+92))
				if v18 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(v16)+92)) = int32(1)
				}
				return int32(0)
			}
		}
	} else {
		return int32(1)
	}
}
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
	var v180 int64
	_ = v180
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v201 int32
	_ = v201
	var v206 int32
	_ = v206
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v219 int32
	_ = v219
	var v227 int32
	_ = v227
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v240 int32
	_ = v240
	var v241 int32
	_ = v241
	var v246 int32
	_ = v246
	var v248 int32
	_ = v248
	var v258 int32
	_ = v258
	var v259 int32
	_ = v259
	var v262 int32
	_ = v262
	var v266 int32
	_ = v266
	var v320 int32
	_ = v320
	var v330 int32
	_ = v330
	var v333 float32
	_ = v333
	var v339 float32
	_ = v339
	var v345 int32
	_ = v345
	var v375 float64
	_ = v375
	var v376 int32
	_ = v376
	var v377 float32
	_ = v377
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v396 int32
	_ = v396
	var v399 int32
	_ = v399
	var v401 int32
	_ = v401
	var v402 int32
	_ = v402
	var v407 int32
	_ = v407
	var v409 int32
	_ = v409
	var v411 int32
	_ = v411
	var v416 int32
	_ = v416
	var v420 int32
	_ = v420
	var v421 int32
	_ = v421
	var v426 int32
	_ = v426
	var v443 int32
	_ = v443
	var v445 int32
	_ = v445
	var v447 int64
	_ = v447
	var v449 int32
	_ = v449
	var v451 int32
	_ = v451
	var v452 int32
	_ = v452
	var v456 int32
	_ = v456
	var v461 int32
	_ = v461
	var v462 int64
	_ = v462
	var v468 int32
	_ = v468
	var v471 int32
	_ = v471
	var v476 int32
	_ = v476
	var v481 int32
	_ = v481
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v486 int32
	_ = v486
	var v487 int32
	_ = v487
	var v489 int32
	_ = v489
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v498 int32
	_ = v498
	var v499 int32
	_ = v499
	var v503 int32
	_ = v503
	var v506 int32
	_ = v506
	var v509 int32
	_ = v509
	var v511 int32
	_ = v511
	var v517 int32
	_ = v517
	var v644 int32
	_ = v644
	var v649 int32
	_ = v649
	var v651 int32
	_ = v651
	var v656 int32
	_ = v656
	var v660 int32
	_ = v660
	var v662 int32
	_ = v662
	var v667 int32
	_ = v667
	var v673 int32
	_ = v673
	var v674 int32
	_ = v674
	var v677 int32
	_ = v677
	var v683 int64
	_ = v683
	var v695 int32
	_ = v695
	var v740 int64
	_ = v740
	var v741 int64
	_ = v741
	var v748 int32
	_ = v748
	var v754 int32
	_ = v754
	var v791 int32
	_ = v791
	var v796 int32
	_ = v796
	var v797 int32
	_ = v797
	var v798 int32
	_ = v798
	var v799 int32
	_ = v799
	var v800 int32
	_ = v800
	var v802 int32
	_ = v802
	var v803 int32
	_ = v803
	var v805 int32
	_ = v805
	var v832 int32
	_ = v832
	var v835 int32
	_ = v835
	var v836 int32
	_ = v836
	var v837 int32
	_ = v837
	var v838 int32
	_ = v838
	var v839 int32
	_ = v839
	var v840 int32
	_ = v840
	var v841 int32
	_ = v841
	var v861 int32
	_ = v861
	var v871 int32
	_ = v871
	var v872 int32
	_ = v872
	var v882 int32
	_ = v882
	var v884 int32
	_ = v884
	var v885 int32
	_ = v885
	var v887 int32
	_ = v887
	var v888 int32
	_ = v888
	var v890 int32
	_ = v890
	var v892 int32
	_ = v892
	var v894 int32
	_ = v894
	var v897 int32
	_ = v897
	var v899 int32
	_ = v899
	var v900 int32
	_ = v900
	var v903 int32
	_ = v903
	var v904 int32
	_ = v904
	var v907 int32
	_ = v907
	var v909 int32
	_ = v909
	var v914 int32
	_ = v914
	var v920 int32
	_ = v920
	var v923 int32
	_ = v923
	var v929 int32
	_ = v929
	var v934 int32
	_ = v934
	var v942 int32
	_ = v942
	var v969 int32
	_ = v969
	var v974 int32
	_ = v974
	var v986 int32
	_ = v986
	var v996 int32
	_ = v996
	var v1002 int32
	_ = v1002
	var v1006 int32
	_ = v1006
	var v1008 int32
	_ = v1008
	var v1009 int32
	_ = v1009
	var v1010 int32
	_ = v1010
	var v1012 int32
	_ = v1012
	var v1014 int32
	_ = v1014
	var v1015 int32
	_ = v1015
	var v1019 int32
	_ = v1019
	var v1021 int32
	_ = v1021
	var v1023 int32
	_ = v1023
	var v1026 int32
	_ = v1026
	var v1028 int32
	_ = v1028
	var v1029 int32
	_ = v1029
	var v1032 int32
	_ = v1032
	var v1033 int32
	_ = v1033
	var v1036 int32
	_ = v1036
	var v1038 int32
	_ = v1038
	var v1043 int32
	_ = v1043
	var v1049 int32
	_ = v1049
	var v1052 int32
	_ = v1052
	var v1058 int32
	_ = v1058
	var v1063 int32
	_ = v1063
	var v1071 int32
	_ = v1071
	var v1093 int32
	_ = v1093
	var v1098 int32
	_ = v1098
	var v1110 int32
	_ = v1110
	var v1120 int32
	_ = v1120
	var v1126 int32
	_ = v1126
	var v1130 int32
	_ = v1130
	var v1132 int32
	_ = v1132
	var v1133 int32
	_ = v1133
	var v1134 int32
	_ = v1134
	var v1136 int32
	_ = v1136
	var v1138 int32
	_ = v1138
	var v1139 int32
	_ = v1139
	var v1143 int32
	_ = v1143
	var v1145 int32
	_ = v1145
	var v1147 int32
	_ = v1147
	var v1150 int32
	_ = v1150
	var v1152 int32
	_ = v1152
	var v1153 int32
	_ = v1153
	var v1156 int32
	_ = v1156
	var v1157 int32
	_ = v1157
	var v1160 int32
	_ = v1160
	var v1162 int32
	_ = v1162
	var v1167 int32
	_ = v1167
	var v1173 int32
	_ = v1173
	var v1176 int32
	_ = v1176
	var v1182 int32
	_ = v1182
	var v1187 int32
	_ = v1187
	var v1195 int32
	_ = v1195
	var v1217 int32
	_ = v1217
	var v1222 int32
	_ = v1222
	var v1225 int32
	_ = v1225
	var v1227 int32
	_ = v1227
	var v1238 int32
	_ = v1238
	var v1241 int32
	_ = v1241
	var v1243 int32
	_ = v1243
	var v1254 int32
	_ = v1254
	var v1298 int32
	_ = v1298
	var v1313 int32
	_ = v1313
	var v1314 int32
	_ = v1314
	var v1315 int32
	_ = v1315
	var v1317 int32
	_ = v1317
	var v1320 int32
	_ = v1320
	var v1323 int32
	_ = v1323
	var v1326 int32
	_ = v1326
	var v1329 int32
	_ = v1329
	var v1332 int32
	_ = v1332
	var v1335 int32
	_ = v1335
	var v1338 int32
	_ = v1338
	var v1341 int32
	_ = v1341
	var v1344 int32
	_ = v1344
	var v1347 int32
	_ = v1347
	var v1350 int32
	_ = v1350
	var v1353 int32
	_ = v1353
	var v1356 int32
	_ = v1356
	var v1359 int32
	_ = v1359
	var v1364 int32
	_ = v1364
	var v1365 int32
	_ = v1365
	var v1366 int32
	_ = v1366
	var v1387 int32
	_ = v1387
	var v1406 int32
	_ = v1406
	var v1438 int32
	_ = v1438
	var v1439 int32
	_ = v1439
	var v1440 int32
	_ = v1440
	var v1441 int32
	_ = v1441
	var v1454 int32
	_ = v1454
	var v1455 int32
	_ = v1455
	var v1456 int32
	_ = v1456
	var v1459 int32
	_ = v1459
	var v1460 int32
	_ = v1460
	var v1464 int32
	_ = v1464
	var v1471 int32
	_ = v1471
	var v1504 int32
	_ = v1504
	var v1515 int32
	_ = v1515
	var v1520 int32
	_ = v1520
	var v1521 int32
	_ = v1521
	var v1523 int32
	_ = v1523
	var v1532 int32
	_ = v1532
	var v1559 int32
	_ = v1559
	var v1561 int32
	_ = v1561
	var v1562 int32
	_ = v1562
	var v1570 int32
	_ = v1570
	var v1571 int32
	_ = v1571
	var v1573 int32
	_ = v1573
	var v1582 int32
	_ = v1582
	var v1584 int32
	_ = v1584
	var v1586 int32
	_ = v1586
	var v1630 int32
	_ = v1630
	var v1635 int32
	_ = v1635
	var v1638 int32
	_ = v1638
	var v1641 int32
	_ = v1641
	var v1644 int32
	_ = v1644
	var v1645 int32
	_ = v1645
	var v1646 int32
	_ = v1646
	var v1650 int32
	_ = v1650
	var v1651 int32
	_ = v1651
	var v1652 int32
	_ = v1652
	var v1657 int32
	_ = v1657
	var v1658 int32
	_ = v1658
	var v1665 int32
	_ = v1665
	var v1701 int32
	_ = v1701
	var v1712 int32
	_ = v1712
	var v1717 int32
	_ = v1717
	var v1718 int32
	_ = v1718
	var v1720 int32
	_ = v1720
	var v1729 int32
	_ = v1729
	var v1756 int32
	_ = v1756
	var v1758 int32
	_ = v1758
	var v1759 int32
	_ = v1759
	var v1767 int32
	_ = v1767
	var v1768 int32
	_ = v1768
	var v1770 int32
	_ = v1770
	var v1779 int32
	_ = v1779
	var v1781 int32
	_ = v1781
	var v1783 int32
	_ = v1783
	var v1827 int32
	_ = v1827
	var v1832 int32
	_ = v1832
	var v1835 int32
	_ = v1835
	var v1838 int32
	_ = v1838
	var v1841 int32
	_ = v1841
	var v1842 int32
	_ = v1842
	var v1843 int32
	_ = v1843
	var v1847 int32
	_ = v1847
	var v1848 int32
	_ = v1848
	var v1849 int32
	_ = v1849
	var v1854 int32
	_ = v1854
	var v1855 int32
	_ = v1855
	var v1862 int32
	_ = v1862
	var v1898 int32
	_ = v1898
	var v1909 int32
	_ = v1909
	var v1914 int32
	_ = v1914
	var v1915 int32
	_ = v1915
	var v1917 int32
	_ = v1917
	var v1926 int32
	_ = v1926
	var v1953 int32
	_ = v1953
	var v1955 int32
	_ = v1955
	var v1956 int32
	_ = v1956
	var v1964 int32
	_ = v1964
	var v1965 int32
	_ = v1965
	var v1967 int32
	_ = v1967
	var v1976 int32
	_ = v1976
	var v1978 int32
	_ = v1978
	var v1980 int32
	_ = v1980
	var v2024 int32
	_ = v2024
	var v2029 int32
	_ = v2029
	var v2032 int32
	_ = v2032
	var v2035 int32
	_ = v2035
	var v2038 int32
	_ = v2038
	var v2045 int32
	_ = v2045
	var v2050 int32
	_ = v2050
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
	var v2081 int32
	_ = v2081
	var v2083 int32
	_ = v2083
	var v2086 int32
	_ = v2086
	var v2179 int32
	_ = v2179
	var v2186 int32
	_ = v2186
	var v2232 int32
	_ = v2232
	var v2234 int32
	_ = v2234
	var v2235 int32
	_ = v2235
	var v2236 int32
	_ = v2236
	var v2241 int32
	_ = v2241
	var v2244 int32
	_ = v2244
	var v2245 int32
	_ = v2245
	var v2248 int32
	_ = v2248
	var v2251 int32
	_ = v2251
	var v2261 int32
	_ = v2261
	var v2271 int32
	_ = v2271
	var v2331 int32
	_ = v2331
	var v2332 int32
	_ = v2332
	var v2337 int32
	_ = v2337
	var v2338 int32
	_ = v2338
	var v2364 int32
	_ = v2364
	var v2365 int32
	_ = v2365
	var v2366 int32
	_ = v2366
	var v2371 int32
	_ = v2371
	var v2424 int32
	_ = v2424
	var v2425 int32
	_ = v2425
	var v2426 int32
	_ = v2426
	var v2428 int32
	_ = v2428
	var v2429 int32
	_ = v2429
	var v2434 int32
	_ = v2434
	var v2437 int32
	_ = v2437
	var v2440 int32
	_ = v2440
	var v2445 int32
	_ = v2445
	var v2448 int32
	_ = v2448
	var v2451 int32
	_ = v2451
	var v2456 int32
	_ = v2456
	var v2459 int32
	_ = v2459
	var v2462 int32
	_ = v2462
	var v2467 int32
	_ = v2467
	var v2470 int32
	_ = v2470
	var v2471 int32
	_ = v2471
	var v2474 int32
	_ = v2474
	var v2479 int32
	_ = v2479
	var v2482 int32
	_ = v2482
	var v2485 int32
	_ = v2485
	var v2490 int32
	_ = v2490
	var v2493 int32
	_ = v2493
	var v2496 int32
	_ = v2496
	var v2501 int32
	_ = v2501
	var v2504 int32
	_ = v2504
	var v2507 int32
	_ = v2507
	var v2512 int32
	_ = v2512
	var v2515 int32
	_ = v2515
	var v2516 int32
	_ = v2516
	var v2519 int32
	_ = v2519
	var v2524 int32
	_ = v2524
	var v2527 int32
	_ = v2527
	var v2530 int32
	_ = v2530
	var v2535 int32
	_ = v2535
	var v2538 int32
	_ = v2538
	var v2541 int32
	_ = v2541
	var v2546 int32
	_ = v2546
	var v2549 int32
	_ = v2549
	var v2552 int32
	_ = v2552
	var v2557 int32
	_ = v2557
	var v2560 int32
	_ = v2560
	var v2561 int32
	_ = v2561
	var v2564 int32
	_ = v2564
	var v2569 int32
	_ = v2569
	var v2572 int32
	_ = v2572
	var v2575 int32
	_ = v2575
	var v2580 int32
	_ = v2580
	var v2583 int32
	_ = v2583
	var v2586 int32
	_ = v2586
	var v2591 int32
	_ = v2591
	var v2594 int32
	_ = v2594
	var v2597 int32
	_ = v2597
	var v2602 int32
	_ = v2602
	var v2629 int32
	_ = v2629
	var v2630 int32
	_ = v2630
	var v2632 int32
	_ = v2632
	var v2633 int32
	_ = v2633
	var v2638 int32
	_ = v2638
	var v2641 int32
	_ = v2641
	var v2644 int32
	_ = v2644
	var v2649 int32
	_ = v2649
	var v2652 int32
	_ = v2652
	var v2653 int32
	_ = v2653
	var v2656 int32
	_ = v2656
	var v2661 int32
	_ = v2661
	var v2664 int32
	_ = v2664
	var v2667 int32
	_ = v2667
	var v2672 int32
	_ = v2672
	var v2675 int32
	_ = v2675
	var v2676 int32
	_ = v2676
	var v2679 int32
	_ = v2679
	var v2684 int32
	_ = v2684
	var v2687 int32
	_ = v2687
	var v2690 int32
	_ = v2690
	var v2695 int32
	_ = v2695
	var v2698 int32
	_ = v2698
	var v2699 int32
	_ = v2699
	var v2702 int32
	_ = v2702
	var v2707 int32
	_ = v2707
	var v2710 int32
	_ = v2710
	var v2713 int32
	_ = v2713
	var v2718 int32
	_ = v2718
	var v2722 int32
	_ = v2722
	var v2723 int32
	_ = v2723
	var v2724 int32
	_ = v2724
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
	var v2775 int32
	_ = v2775
	var v2780 int32
	_ = v2780
	var v2783 int32
	_ = v2783
	var v2785 int32
	_ = v2785
	var v2789 int64
	_ = v2789
	var v2790 int64
	_ = v2790
	var v2792 int32
	_ = v2792
	var v2796 int32
	_ = v2796
	var v2797 int32
	_ = v2797
	var v2798 int32
	_ = v2798
	var v2799 int32
	_ = v2799
	var v2800 int32
	_ = v2800
	var v2804 int32
	_ = v2804
	var v2805 int32
	_ = v2805
	var v2807 int32
	_ = v2807
	var v2808 int32
	_ = v2808
	var v2810 int32
	_ = v2810
	var v2811 int32
	_ = v2811
	var v2813 int32
	_ = v2813
	var v2814 int32
	_ = v2814
	var v2816 int32
	_ = v2816
	var v2817 int32
	_ = v2817
	var v2819 int32
	_ = v2819
	var v2820 int32
	_ = v2820
	var v2822 int32
	_ = v2822
	var v2823 int32
	_ = v2823
	var v2825 int32
	_ = v2825
	var v2826 int32
	_ = v2826
	var v2828 int32
	_ = v2828
	var v2829 int32
	_ = v2829
	var v2831 int32
	_ = v2831
	var v2832 int32
	_ = v2832
	var v2834 int32
	_ = v2834
	var v2835 int32
	_ = v2835
	var v2837 int32
	_ = v2837
	var v2838 int32
	_ = v2838
	var v2840 int32
	_ = v2840
	var v2841 int32
	_ = v2841
	var v2843 int32
	_ = v2843
	var v2844 int32
	_ = v2844
	var v2846 int32
	_ = v2846
	var v2847 int32
	_ = v2847
	var v2849 int32
	_ = v2849
	var v2850 int32
	_ = v2850
	var v2852 int32
	_ = v2852
	var v2853 int32
	_ = v2853
	var v2855 int32
	_ = v2855
	var v2856 int32
	_ = v2856
	var v2858 int32
	_ = v2858
	var v2859 int32
	_ = v2859
	var v2861 int32
	_ = v2861
	var v2862 int32
	_ = v2862
	var v2864 int32
	_ = v2864
	var v2865 int32
	_ = v2865
	var v2867 int32
	_ = v2867
	var v2868 int32
	_ = v2868
	var v2870 int32
	_ = v2870
	var v2871 int32
	_ = v2871
	var v2873 int32
	_ = v2873
	var v2874 int32
	_ = v2874
	var v2876 int32
	_ = v2876
	var v2877 int32
	_ = v2877
	var v2879 int32
	_ = v2879
	var v2880 int32
	_ = v2880
	var v2882 int32
	_ = v2882
	var v2883 int32
	_ = v2883
	var v2885 int32
	_ = v2885
	var v2886 int32
	_ = v2886
	var v2888 int32
	_ = v2888
	var v2889 int32
	_ = v2889
	var v2891 int32
	_ = v2891
	var v2892 int32
	_ = v2892
	var v2894 int32
	_ = v2894
	var v2895 int32
	_ = v2895
	var v2897 int32
	_ = v2897
	var v2898 int32
	_ = v2898
	var v2900 int32
	_ = v2900
	var v2901 int32
	_ = v2901
	var v2903 int32
	_ = v2903
	var v2904 int32
	_ = v2904
	var v2906 int32
	_ = v2906
	var v2909 int32
	_ = v2909
	var v2910 int32
	_ = v2910
	var v2912 int32
	_ = v2912
	var v2915 int32
	_ = v2915
	var v2916 int32
	_ = v2916
	var v2918 int32
	_ = v2918
	var v2922 int32
	_ = v2922
	var v2923 int64
	_ = v2923
	var v2925 int32
	_ = v2925
	var v2929 int64
	_ = v2929
	var v2931 int32
	_ = v2931
	var v2932 int64
	_ = v2932
	var v2938 int64
	_ = v2938
	var v2943 int32
	_ = v2943
	var v2956 int32
	_ = v2956
	var v2959 int32
	_ = v2959
	var v2960 int32
	_ = v2960
	var v2961 int32
	_ = v2961
	var v2964 int32
	_ = v2964
	var v2967 int32
	_ = v2967
	var v2968 int32
	_ = v2968
	var v2971 int32
	_ = v2971
	var v2972 int32
	_ = v2972
	var v2974 int32
	_ = v2974
	var v2977 int32
	_ = v2977
	var v2979 int32
	_ = v2979
	var v2983 int32
	_ = v2983
	var v2987 int32
	_ = v2987
	var v2988 int32
	_ = v2988
	var v2989 int32
	_ = v2989
	var v2990 int32
	_ = v2990
	var v2991 int32
	_ = v2991
	var v2995 int32
	_ = v2995
	var v2996 int32
	_ = v2996
	var v2998 int32
	_ = v2998
	var v2999 int32
	_ = v2999
	var v3001 int32
	_ = v3001
	var v3002 int32
	_ = v3002
	var v3004 int32
	_ = v3004
	var v3005 int32
	_ = v3005
	var v3007 int32
	_ = v3007
	var v3008 int32
	_ = v3008
	var v3010 int32
	_ = v3010
	var v3011 int32
	_ = v3011
	var v3013 int32
	_ = v3013
	var v3014 int32
	_ = v3014
	var v3016 int32
	_ = v3016
	var v3017 int32
	_ = v3017
	var v3019 int32
	_ = v3019
	var v3020 int32
	_ = v3020
	var v3022 int32
	_ = v3022
	var v3023 int32
	_ = v3023
	var v3025 int32
	_ = v3025
	var v3026 int32
	_ = v3026
	var v3028 int32
	_ = v3028
	var v3029 int32
	_ = v3029
	var v3031 int32
	_ = v3031
	var v3032 int32
	_ = v3032
	var v3034 int32
	_ = v3034
	var v3035 int32
	_ = v3035
	var v3037 int32
	_ = v3037
	var v3038 int32
	_ = v3038
	var v3040 int32
	_ = v3040
	var v3041 int32
	_ = v3041
	var v3043 int32
	_ = v3043
	var v3044 int32
	_ = v3044
	var v3046 int32
	_ = v3046
	var v3047 int32
	_ = v3047
	var v3049 int32
	_ = v3049
	var v3050 int32
	_ = v3050
	var v3052 int32
	_ = v3052
	var v3053 int32
	_ = v3053
	var v3055 int32
	_ = v3055
	var v3056 int32
	_ = v3056
	var v3058 int32
	_ = v3058
	var v3059 int32
	_ = v3059
	var v3061 int32
	_ = v3061
	var v3062 int32
	_ = v3062
	var v3064 int32
	_ = v3064
	var v3065 int32
	_ = v3065
	var v3067 int32
	_ = v3067
	var v3068 int32
	_ = v3068
	var v3070 int32
	_ = v3070
	var v3071 int32
	_ = v3071
	var v3073 int32
	_ = v3073
	var v3074 int32
	_ = v3074
	var v3076 int32
	_ = v3076
	var v3077 int32
	_ = v3077
	var v3079 int32
	_ = v3079
	var v3080 int32
	_ = v3080
	var v3082 int32
	_ = v3082
	var v3083 int32
	_ = v3083
	var v3085 int32
	_ = v3085
	var v3086 int32
	_ = v3086
	var v3088 int32
	_ = v3088
	var v3089 int32
	_ = v3089
	var v3091 int32
	_ = v3091
	var v3092 int32
	_ = v3092
	var v3094 int32
	_ = v3094
	var v3095 int32
	_ = v3095
	var v3097 int32
	_ = v3097
	var v3100 int32
	_ = v3100
	var v3101 int32
	_ = v3101
	var v3103 int32
	_ = v3103
	var v3106 int32
	_ = v3106
	var v3107 int32
	_ = v3107
	var v3109 int32
	_ = v3109
	var v3113 int32
	_ = v3113
	var v3114 int64
	_ = v3114
	var v3116 int32
	_ = v3116
	var v3120 int64
	_ = v3120
	var v3122 int32
	_ = v3122
	var v3123 int64
	_ = v3123
	var v3129 int64
	_ = v3129
	var v3135 int64
	_ = v3135
	var v3136 int64
	_ = v3136
	var v3138 int32
	_ = v3138
	var v3143 int32
	_ = v3143
	var v3145 int32
	_ = v3145
	var v3147 int32
	_ = v3147
	var v3148 int32
	_ = v3148
	var v3152 int32
	_ = v3152
	var v3154 int64
	_ = v3154
	var v3156 int32
	_ = v3156
	var v3158 int32
	_ = v3158
	var v3160 int32
	_ = v3160
	var v3161 int32
	_ = v3161
	var v3163 int32
	_ = v3163
	var v3167 int32
	_ = v3167
	var v3173 int32
	_ = v3173
	var v3174 int32
	_ = v3174
	var v3183 int32
	_ = v3183
	var v3190 int32
	_ = v3190
	var v3192 int32
	_ = v3192
	var v3193 int32
	_ = v3193
	var v3196 int32
	_ = v3196
	var v3200 int32
	_ = v3200
	var v3201 int64
	_ = v3201
	var v3207 int32
	_ = v3207
	var v3210 int32
	_ = v3210
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
	var v3228 int32
	_ = v3228
	var v3232 int32
	_ = v3232
	var v3233 int32
	_ = v3233
	var v3236 int32
	_ = v3236
	var v3243 int32
	_ = v3243
	var v3249 int64
	_ = v3249
	var v3250 int64
	_ = v3250
	var v3253 int32
	_ = v3253
	var v3290 int32
	_ = v3290
	var v3295 int32
	_ = v3295
	var v3296 int32
	_ = v3296
	var v3297 int32
	_ = v3297
	var v3298 int32
	_ = v3298
	var v3299 int32
	_ = v3299
	var v3301 int32
	_ = v3301
	var v3302 int32
	_ = v3302
	var v3303 int32
	_ = v3303
	var v3304 int32
	_ = v3304
	var v3331 int32
	_ = v3331
	var v3332 int32
	_ = v3332
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
	var v3360 int32
	_ = v3360
	var v3361 int32
	_ = v3361
	var v3370 int32
	_ = v3370
	var v3371 int32
	_ = v3371
	var v3381 int32
	_ = v3381
	var v3383 int32
	_ = v3383
	var v3384 int32
	_ = v3384
	var v3386 int32
	_ = v3386
	var v3387 int32
	_ = v3387
	var v3389 int32
	_ = v3389
	var v3391 int32
	_ = v3391
	var v3393 int32
	_ = v3393
	var v3396 int32
	_ = v3396
	var v3398 int32
	_ = v3398
	var v3399 int32
	_ = v3399
	var v3402 int32
	_ = v3402
	var v3403 int32
	_ = v3403
	var v3406 int32
	_ = v3406
	var v3408 int32
	_ = v3408
	var v3413 int32
	_ = v3413
	var v3419 int32
	_ = v3419
	var v3420 int32
	_ = v3420
	var v3422 int32
	_ = v3422
	var v3428 int32
	_ = v3428
	var v3433 int32
	_ = v3433
	var v3441 int32
	_ = v3441
	var v3444 int32
	_ = v3444
	var v3446 int32
	_ = v3446
	var v3453 int32
	_ = v3453
	var v3454 int32
	_ = v3454
	var v3468 int32
	_ = v3468
	var v3469 int32
	_ = v3469
	var v3473 int32
	_ = v3473
	var v3485 int32
	_ = v3485
	var v3486 int32
	_ = v3486
	var v3495 int32
	_ = v3495
	var v3501 int32
	_ = v3501
	var v3505 int32
	_ = v3505
	var v3507 int32
	_ = v3507
	var v3508 int32
	_ = v3508
	var v3509 int32
	_ = v3509
	var v3511 int32
	_ = v3511
	var v3513 int32
	_ = v3513
	var v3514 int32
	_ = v3514
	var v3518 int32
	_ = v3518
	var v3520 int32
	_ = v3520
	var v3522 int32
	_ = v3522
	var v3525 int32
	_ = v3525
	var v3527 int32
	_ = v3527
	var v3528 int32
	_ = v3528
	var v3531 int32
	_ = v3531
	var v3532 int32
	_ = v3532
	var v3535 int32
	_ = v3535
	var v3537 int32
	_ = v3537
	var v3542 int32
	_ = v3542
	var v3548 int32
	_ = v3548
	var v3551 int32
	_ = v3551
	var v3557 int32
	_ = v3557
	var v3562 int32
	_ = v3562
	var v3570 int32
	_ = v3570
	var v3574 int32
	_ = v3574
	var v3575 int32
	_ = v3575
	var v3582 int32
	_ = v3582
	var v3583 int32
	_ = v3583
	var v3592 int32
	_ = v3592
	var v3593 int32
	_ = v3593
	var v3597 int32
	_ = v3597
	var v3609 int32
	_ = v3609
	var v3610 int32
	_ = v3610
	var v3619 int32
	_ = v3619
	var v3625 int32
	_ = v3625
	var v3629 int32
	_ = v3629
	var v3631 int32
	_ = v3631
	var v3632 int32
	_ = v3632
	var v3633 int32
	_ = v3633
	var v3635 int32
	_ = v3635
	var v3637 int32
	_ = v3637
	var v3638 int32
	_ = v3638
	var v3642 int32
	_ = v3642
	var v3644 int32
	_ = v3644
	var v3646 int32
	_ = v3646
	var v3649 int32
	_ = v3649
	var v3651 int32
	_ = v3651
	var v3652 int32
	_ = v3652
	var v3655 int32
	_ = v3655
	var v3656 int32
	_ = v3656
	var v3659 int32
	_ = v3659
	var v3661 int32
	_ = v3661
	var v3666 int32
	_ = v3666
	var v3672 int32
	_ = v3672
	var v3675 int32
	_ = v3675
	var v3681 int32
	_ = v3681
	var v3686 int32
	_ = v3686
	var v3694 int32
	_ = v3694
	var v3698 int32
	_ = v3698
	var v3699 int32
	_ = v3699
	var v3706 int32
	_ = v3706
	var v3707 int32
	_ = v3707
	var v3716 int32
	_ = v3716
	var v3717 int32
	_ = v3717
	var v3721 int32
	_ = v3721
	var v3724 int32
	_ = v3724
	var v3726 int32
	_ = v3726
	var v3737 int32
	_ = v3737
	var v3740 int32
	_ = v3740
	var v3742 int32
	_ = v3742
	var v3753 int32
	_ = v3753
	var v3767 int32
	_ = v3767
	var v3769 int32
	_ = v3769
	var v3777 int32
	_ = v3777
	var v3778 int32
	_ = v3778
	var v3786 int32
	_ = v3786
	var v3787 int32
	_ = v3787
	var v3788 int32
	_ = v3788
	var v3794 int32
	_ = v3794
	var v3798 int32
	_ = v3798
	var v3799 int32
	_ = v3799
	var v3801 int32
	_ = v3801
	var v3816 int32
	_ = v3816
	var v3823 int32
	_ = v3823
	var v3826 int32
	_ = v3826
	var v3830 int32
	_ = v3830
	var v3831 int32
	_ = v3831
	var v3833 int32
	_ = v3833
	var v3838 int32
	_ = v3838
	var v3848 int32
	_ = v3848
	var v3864 float64
	_ = v3864
	var v3870 float64
	_ = v3870
	var v3881 int64
	_ = v3881
	var v3896 int32
	_ = v3896
	var v3898 int64
	_ = v3898
	var v3907 int64
	_ = v3907
	var v3912 int64
	_ = v3912
	var v3913 int32
	_ = v3913
	var v3915 int32
	_ = v3915
	var v3917 int32
	_ = v3917
	var v3921 float64
	_ = v3921
	var v3923 float64
	_ = v3923
	var v3936 float64
	_ = v3936
	var v3939 float64
	_ = v3939
	var v3944 float64
	_ = v3944
	var v3945 float64
	_ = v3945
	var v3946 float64
	_ = v3946
	var v3947 float64
	_ = v3947
	var v3952 float64
	_ = v3952
	var v3953 float64
	_ = v3953
	var v3954 float64
	_ = v3954
	var v3979 float64
	_ = v3979
	var v3991 float64
	_ = v3991
	var v4013 float64
	_ = v4013
	var v4016 float64
	_ = v4016
	var v4017 int32
	_ = v4017
	var v4027 int32
	_ = v4027
	var v4028 int32
	_ = v4028
	var v4029 int32
	_ = v4029
	var v4038 int64
	_ = v4038
	var v4056 float32
	_ = v4056
	var v4066 float32
	_ = v4066
	var v4069 float32
	_ = v4069
	var v4072 float32
	_ = v4072
	var v4075 float32
	_ = v4075
	var v4076 float32
	_ = v4076
	var v4078 float32
	_ = v4078
	var v4080 float32
	_ = v4080
	var v4081 float32
	_ = v4081
	var v4082 float32
	_ = v4082
	var v4083 float64
	_ = v4083
	var v4084 int32
	_ = v4084
	var v4085 float32
	_ = v4085
	var v4087 int32
	_ = v4087
	var v4088 float32
	_ = v4088
	var v4089 float32
	_ = v4089
	var v4090 float64
	_ = v4090
	var v4091 int32
	_ = v4091
	var v4092 float32
	_ = v4092
	var v4114 int32
	_ = v4114
	var v4156 int32
	_ = v4156
	var v4193 int32
	_ = v4193
	var v4198 int32
	_ = v4198
	var v4199 int32
	_ = v4199
	var v4200 int32
	_ = v4200
	var v4201 int32
	_ = v4201
	var v4202 int32
	_ = v4202
	var v4204 int32
	_ = v4204
	var v4205 int32
	_ = v4205
	var v4207 int32
	_ = v4207
	var v4234 int32
	_ = v4234
	var v4237 int32
	_ = v4237
	var v4238 int32
	_ = v4238
	var v4239 int32
	_ = v4239
	var v4240 int32
	_ = v4240
	var v4241 int32
	_ = v4241
	var v4242 int32
	_ = v4242
	var v4243 int32
	_ = v4243
	var v4263 int32
	_ = v4263
	var v4273 int32
	_ = v4273
	var v4274 int32
	_ = v4274
	var v4284 int32
	_ = v4284
	var v4286 int32
	_ = v4286
	var v4287 int32
	_ = v4287
	var v4289 int32
	_ = v4289
	var v4290 int32
	_ = v4290
	var v4292 int32
	_ = v4292
	var v4294 int32
	_ = v4294
	var v4296 int32
	_ = v4296
	var v4299 int32
	_ = v4299
	var v4301 int32
	_ = v4301
	var v4302 int32
	_ = v4302
	var v4305 int32
	_ = v4305
	var v4306 int32
	_ = v4306
	var v4309 int32
	_ = v4309
	var v4311 int32
	_ = v4311
	var v4316 int32
	_ = v4316
	var v4322 int32
	_ = v4322
	var v4325 int32
	_ = v4325
	var v4331 int32
	_ = v4331
	var v4336 int32
	_ = v4336
	var v4344 int32
	_ = v4344
	var v4371 int32
	_ = v4371
	var v4376 int32
	_ = v4376
	var v4388 int32
	_ = v4388
	var v4398 int32
	_ = v4398
	var v4404 int32
	_ = v4404
	var v4408 int32
	_ = v4408
	var v4410 int32
	_ = v4410
	var v4411 int32
	_ = v4411
	var v4412 int32
	_ = v4412
	var v4414 int32
	_ = v4414
	var v4416 int32
	_ = v4416
	var v4417 int32
	_ = v4417
	var v4421 int32
	_ = v4421
	var v4423 int32
	_ = v4423
	var v4425 int32
	_ = v4425
	var v4428 int32
	_ = v4428
	var v4430 int32
	_ = v4430
	var v4431 int32
	_ = v4431
	var v4434 int32
	_ = v4434
	var v4435 int32
	_ = v4435
	var v4438 int32
	_ = v4438
	var v4440 int32
	_ = v4440
	var v4445 int32
	_ = v4445
	var v4451 int32
	_ = v4451
	var v4454 int32
	_ = v4454
	var v4460 int32
	_ = v4460
	var v4465 int32
	_ = v4465
	var v4473 int32
	_ = v4473
	var v4495 int32
	_ = v4495
	var v4500 int32
	_ = v4500
	var v4512 int32
	_ = v4512
	var v4522 int32
	_ = v4522
	var v4528 int32
	_ = v4528
	var v4532 int32
	_ = v4532
	var v4534 int32
	_ = v4534
	var v4535 int32
	_ = v4535
	var v4536 int32
	_ = v4536
	var v4538 int32
	_ = v4538
	var v4540 int32
	_ = v4540
	var v4541 int32
	_ = v4541
	var v4545 int32
	_ = v4545
	var v4547 int32
	_ = v4547
	var v4549 int32
	_ = v4549
	var v4552 int32
	_ = v4552
	var v4554 int32
	_ = v4554
	var v4555 int32
	_ = v4555
	var v4558 int32
	_ = v4558
	var v4559 int32
	_ = v4559
	var v4562 int32
	_ = v4562
	var v4564 int32
	_ = v4564
	var v4569 int32
	_ = v4569
	var v4575 int32
	_ = v4575
	var v4578 int32
	_ = v4578
	var v4584 int32
	_ = v4584
	var v4589 int32
	_ = v4589
	var v4597 int32
	_ = v4597
	var v4619 int32
	_ = v4619
	var v4624 int32
	_ = v4624
	var v4627 int32
	_ = v4627
	var v4629 int32
	_ = v4629
	var v4640 int32
	_ = v4640
	var v4643 int32
	_ = v4643
	var v4645 int32
	_ = v4645
	var v4656 int32
	_ = v4656
	var v4661 int32
	_ = v4661
	var v4672 int32
	_ = v4672
	var v4681 int32
	_ = v4681
	var v4689 int32
	_ = v4689
	var v4690 int32
	_ = v4690
	var v4691 int32
	_ = v4691
	var v4692 int32
	_ = v4692
	var v4706 int32
	_ = v4706
	var v4707 int32
	_ = v4707
	var v4710 int32
	_ = v4710
	var v4722 int32
	_ = v4722
	var v4723 int32
	_ = v4723
	var v4724 int32
	_ = v4724
	var v4728 int32
	_ = v4728
	var v4823 int32
	_ = v4823
	var v4824 int32
	_ = v4824
	var v4828 int32
	_ = v4828
	var v4834 int32
	_ = v4834
	var v4893 int32
	_ = v4893
	var v4898 int32
	_ = v4898
	var v4901 int32
	_ = v4901
	var v4912 int32
	_ = v4912
	var v4913 int32
	_ = v4913
	var v4914 int32
	_ = v4914
	var v4915 int32
	_ = v4915
	var v4919 int32
	_ = v4919
	var v4920 int32
	_ = v4920
	var v4924 int32
	_ = v4924
	var v4925 int32
	_ = v4925
	var v4930 int32
	_ = v4930
	var v4935 int32
	_ = v4935
	var v4936 int32
	_ = v4936
	var v4939 int64
	_ = v4939
	var v4940 int64
	_ = v4940
	var v4942 int64
	_ = v4942
	var v4943 int64
	_ = v4943
	var v4945 int64
	_ = v4945
	var v4949 int64
	_ = v4949
	var v4951 int64
	_ = v4951
	var v4955 int64
	_ = v4955
	var v4957 int64
	_ = v4957
	var v4961 int64
	_ = v4961
	var v4963 int64
	_ = v4963
	var v4967 int64
	_ = v4967
	var v4969 int64
	_ = v4969
	var v4973 int64
	_ = v4973
	var v4975 int64
	_ = v4975
	var v4979 int64
	_ = v4979
	var v4981 int64
	_ = v4981
	var v4985 int64
	_ = v4985
	var v4987 int64
	_ = v4987
	var v4991 int64
	_ = v4991
	var v4993 int64
	_ = v4993
	var v4997 int64
	_ = v4997
	var v4999 int64
	_ = v4999
	var v5003 int64
	_ = v5003
	var v5005 int64
	_ = v5005
	var v5009 int64
	_ = v5009
	var v5019 int32
	_ = v5019
	var v5021 int32
	_ = v5021
	var v5028 int32
	_ = v5028
	var v5032 int32
	_ = v5032
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
	return v5032
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
	v227 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v227 < int32(1) {
		goto L30
	} else {
		goto L31
	}
L17:
	;
	v180 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v116)+16)) = v180
	*(*int64)(unsafe.Add(mBase, uint32(v116)+8)) = int64(-34359738368)
	*(*int64)(unsafe.Add(mBase, uint32(v116))) = int64(254)
	*(*int64)(unsafe.Add(mBase, uint32(v116+int32(24)))) = v180
	if v91 == int32(0) {
		goto L19
	} else {
		goto L20
	}
L18:
	;
	if v219 != 0 {
		v116 = v116 + int32(32)
		v121 = v175
		goto L15
	} else {
		goto L28
	}
L19:
	;
	v219 = int32(1)
	goto L18
L20:
	;
	v193 = int32(1024)
	if base.Ui32(v193) < base.Ui32(v91) {
		goto L22
	} else {
		goto L23
	}
L21:
	;
	v201 = *(*int32)(unsafe.Add(mBase, uint32(v116)+20))
	if v201 == int32(0) {
		goto L26
	} else {
		goto L27
	}
L22:
	;
	v196 = v91
	goto L24
L23:
	;
	v196 = v193
	goto L24
L24:
	;
	v197 = F_WebPSafeMalloc(m, int64(1), v196)
	mBase = m.M
	if v197 != 0 {
		goto L21
	} else {
		goto L25
	}
L25:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v116)+28)) = int32(1)
	v219 = int32(0)
	goto L18
L26:
	;
	v208 = *(*int32)(unsafe.Add(mBase, uint32(v116)+16))
	F_WebPSafeFree(m, v208)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v116)+24)) = v196
	*(*int32)(unsafe.Add(mBase, uint32(v116)+16)) = v197
	goto L19
L27:
	;
	v206 = *(*int32)(unsafe.Add(mBase, uint32(v116+int32(16))))
	v207 = F_memcpy(m, v197, v206, v201)
	mBase = m.M
	goto L26
L28:
	;
	goto L16
L29:
	;
	v246 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v248 = *(*int32)(unsafe.Add(mBase, uint32(v246)+92))
	if v248 != 0 {
		goto L36
	} else {
		goto L37
	}
L30:
	;
	goto L29
L31:
	;
	v234 = l0 + int32(88)
	v235 = int32(0)
	goto L32
L32:
	;
	F_VP8BitWriterWipeOut(m, v234)
	mBase = m.M
	v240 = v235 + int32(1)
	v241 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v240 < v241 {
		v234 = v234 + int32(32)
		v235 = v240
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
	*(*int32)(unsafe.Add(mBase, uint32(v246)+92)) = int32(1)
	goto L36
L38:
	;
	v5032 = int32(0)
	goto L13
L39:
	;
	v4893 = v64 + int32(880)
	v4898 = *(*int32)(unsafe.Add(mBase, uint32(v4893)+24))
	if v4834 == int32(0) {
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
	v258 = v80 >> (uint(int32(3)) % 32)
	v259 = int32(96)
	if v259 < v258 {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v4114 = int32(40)
	goto L40
L43:
	;
	v262 = v258
	goto L45
L44:
	;
	v262 = v259
	goto L45
L45:
	;
	v266 = l0 + int32(344)
	v320 = int32(40)
	v330 = v68
	v333 = v77
	v339 = float32(10)
	v345 = v320
	v375 = float64(0)
	v376 = int32(1)
	v377 = v77
	goto L46
L46:
	;
	v388 = v330 + int32(-1)
	v389 = int32(1)
	if base.F64_le(base.F64_promote_f32(base.F32_abs(v339)), float64(0.4)) != 0 {
		v399 = v389
		goto L48
	} else {
		goto L49
	}
L47:
	;
	v4114 = v651
	goto L40
L48:
	;
	v401 = v64 + int32(880)
	v402 = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(v401))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v401)+24)) = l0
	v407 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[2])))
	*(*int32)(unsafe.Add(mBase, uint32(v401)+280)) = v407
	v409 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	*(*int32)(unsafe.Add(mBase, uint32(v401)+296)) = v409
	v411 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[3])))
	*(*int32)(unsafe.Add(mBase, uint32(v401)+304)) = v411
	*(*int32)(unsafe.Add(mBase, uint32(v401)+32)) = l0 + int32(88)
	v416 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[4])))
	*(*int32)(unsafe.Add(mBase, uint32(v401)+36)) = v416
	v420 = int32(-32)
	v421 = (v64 + int32(1327)) & v420
	*(*int32)(unsafe.Add(mBase, uint32(v401)+8)) = v421
	v426 = (v64 + int32(1240)) & v420
	*(*int32)(unsafe.Add(mBase, uint32(v401)+308)) = v426
	*(*int32)(unsafe.Add(mBase, uint32(v401)+20)) = v421 + int32(1536)
	*(*int32)(unsafe.Add(mBase, uint32(v401)+16)) = v421 + int32(1024)
	*(*int32)(unsafe.Add(mBase, uint32(v401)+12)) = v421 + int32(512)
	*(*int32)(unsafe.Add(mBase, uint32(v401)+316)) = v426 + int32(48)
	*(*int32)(unsafe.Add(mBase, uint32(v401)+312)) = v426 + int32(32)
	v443 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v401)+40)) = v443
	v445 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[6])))
	*(*int32)(unsafe.Add(mBase, uint32(v401)+28)) = v445
	v447 = *(*int64)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v401)+320)) = v447
	v449 = int32(127)
	*(*uint8)(unsafe.Add(mBase, uint32(v426)+47)) = uint8(v449)
	v451 = *(*int32)(unsafe.Add(mBase, uint32(v401)+312))
	v452 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v451+v452))) = uint8(v449)
	v456 = *(*int32)(unsafe.Add(mBase, uint32(v401)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v456+v452))) = uint8(v449)
	v461 = *(*int32)(unsafe.Add(mBase, uint32(v401)+308))
	v462 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v461))) = v462
	*(*int64)(unsafe.Add(mBase, uint32(v461+int32(8)))) = v462
	v468 = *(*int32)(unsafe.Add(mBase, uint32(v401)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v468))) = v462
	v471 = *(*int32)(unsafe.Add(mBase, uint32(v401)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v471))) = v462
	*(*int32)(unsafe.Add(mBase, uint32(v401)+160)) = v402
	v476 = *(*int32)(unsafe.Add(mBase, uint32(v401)+304))
	if v476 == v402 {
		goto L52
	} else {
		goto L53
	}
L49:
	;
	if v388 == int32(0) {
		v399 = v389
		goto L48
	} else {
		goto L50
	}
L50:
	;
	v396 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[8])))
	v399 = base.B2i32(v396 == int32(0))
	goto L48
L51:
	;
	F_SetLoopParams(m, l0, v333)
	mBase = m.M
	v517 = base.I32_div_s(v345, v330+int32(1))
	if v399 == int32(0) {
		goto L56
	} else {
		goto L57
	}
L52:
	;
	v481 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v482 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v483 = v481 * v482
	*(*int32)(unsafe.Add(mBase, uint32(v401)+292)) = v483
	*(*int32)(unsafe.Add(mBase, uint32(v401)+288)) = v483
	v486 = *(*int32)(unsafe.Add(mBase, uint32(v401)+24))
	v487 = *(*int32)(unsafe.Add(mBase, uint32(v486)+uint32(_c_F_VP8EncTokenLoop[7])))
	v489 = *(*int32)(unsafe.Add(mBase, uint32(v486)+40))
	v492 = F_memset(m, v487, int32(127), v489<<(uint(int32(5))%32))
	mBase = m.M
	v493 = *(*int32)(unsafe.Add(mBase, uint32(v486)+uint32(_c_F_VP8EncTokenLoop[5])))
	v494 = int32(0)
	v495 = *(*int32)(unsafe.Add(mBase, uint32(v486)+40))
	v498 = F_memset(m, v493, v494, v495<<(uint(int32(2))%32))
	mBase = m.M
	v499 = *(*int32)(unsafe.Add(mBase, uint32(v486)+uint32(_c_F_VP8EncTokenLoop[3])))
	if v499 == v494 {
		goto L54
	} else {
		goto L55
	}
L53:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v401)+300)) = int32(0)
	goto L52
L54:
	;
	v509 = int32(0)
	v511 = F_memset(m, v64+int32(1048), v509, int32(96))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v401)+284)) = v509
	goto L51
L55:
	;
	v503 = *(*int32)(unsafe.Add(mBase, uint32(v486)+40))
	v506 = F_memset(m, v499, int32(0), v503<<(uint(int32(2))%32))
	mBase = m.M
	goto L54
L56:
	;
	v651 = v345 - v517
	if v266 == int32(0) {
		goto L76
	} else {
		goto L77
	}
L57:
	;
	goto L60
L58:
	;
	v644 = *(*int32)(unsafe.Add(mBase, uint32(v64+int32(880))+280))
	if v644 == int32(0) {
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
	v649 = F_memset(m, v644, int32(0), int32(2048))
	mBase = m.M
	F_VP8SSIMDspInit(m)
	mBase = m.M
	goto L73
L75:
	;
	v683 = int64(0)
	v695 = v262
	v740 = v683
	v741 = v683
	goto L86
L76:
	;
	goto L75
L77:
	;
	v656 = *(*int32)(unsafe.Add(mBase, uint32(v266)))
	if v656 == int32(0) {
		goto L78
	} else {
		goto L79
	}
L78:
	;
	v667 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v266)+20)) = v667
	*(*int64)(unsafe.Add(mBase, uint32(v266)+8)) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v266))) = v667
	v673 = *(*int32)(unsafe.Add(mBase, uint32(v266)+16))
	v674 = int32(_a_F_VP8EncTokenLoop_4)
	if v674 < v673 {
		goto L83
	} else {
		goto L84
	}
L79:
	;
	v660 = v656
	goto L80
L80:
	;
	v662 = *(*int32)(unsafe.Add(mBase, uint32(v660)))
	F_WebPSafeFree(m, v660)
	mBase = m.M
	if v662 != 0 {
		v660 = v662
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
	v677 = v673
	goto L85
L84:
	;
	v677 = v674
	goto L85
L85:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v266)+16)) = v677
	*(*int32)(unsafe.Add(mBase, uint32(v266)+4)) = v266
	goto L76
L86:
	;
	v748 = int32(0)
	F_VP8IteratorImport(m, v64+int32(880), v748)
	mBase = m.M
	if v695 <= v748 {
		goto L89
	} else {
		goto L90
	}
L87:
	;
	v3249 = int64(*(*int32)(unsafe.Add(mBase, uint32(l0)+32)))
	v3250 = v3136 + v3249
	if v92 == int32(0) {
		goto L221
	} else {
		goto L222
	}
L88:
	;
	v2234 = v64 + int32(880)
	v2235 = F_VP8Decimate(m, v2234, v64, v109)
	mBase = m.M
	v2236 = *(*int32)(unsafe.Add(mBase, uint32(v64)+904))
	v2241 = *(*int32)(unsafe.Add(mBase, uint32(v2234)+40))
	v2244 = *(*int32)(unsafe.Add(mBase, uint32(v2241+int32(-4))))
	v2245 = *(*int32)(unsafe.Add(mBase, uint32(v2241)))
	v2248 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+124)) = int32(base.Ui32(v2245)>>(uint(int32(24))%32)) & v2248
	v2251 = int32(23)
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+120)) = int32(base.Ui32(v2245)>>(uint(v2251)%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+116)) = int32(base.Ui32(v2245)>>(uint(int32(22))%32)) & v2248
	v2261 = int32(19)
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+112)) = int32(base.Ui32(v2245)>>(uint(v2261)%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+108)) = int32(base.Ui32(v2245)>>(uint(int32(18))%32)) & v2248
	v2271 = int32(15)
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+104)) = int32(base.Ui32(v2245)>>(uint(v2271)%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+100)) = int32(base.Ui32(v2245)>>(uint(int32(14))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+96)) = int32(base.Ui32(v2245)>>(uint(int32(13))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+92)) = int32(base.Ui32(v2245)>>(uint(int32(12))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+156)) = int32(base.Ui32(v2244)>>(uint(v2251)%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+152)) = int32(base.Ui32(v2244)>>(uint(int32(21))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+148)) = int32(base.Ui32(v2244)>>(uint(v2261)%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+144)) = int32(base.Ui32(v2244)>>(uint(int32(17))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+140)) = int32(base.Ui32(v2244)>>(uint(v2271)%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+136)) = int32(base.Ui32(v2244)>>(uint(int32(11))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+132)) = int32(base.Ui32(v2244)>>(uint(int32(7))%32)) & v2248
	*(*int32)(unsafe.Add(mBase, uint32(v2234)+128)) = int32(base.Ui32(v2244)>>(uint(int32(3))%32)) & v2248
	goto L176
L89:
	;
	v754 = int32(0)
	v791 = l0 + int32(_a_F_VP8EncTokenLoop_2)
	v796 = l0 + int32(_a_F_VP8EncTokenLoop_5)
	v797 = l0 + int32(3442)
	v798 = l0 + int32(_a_F_VP8EncTokenLoop_6)
	v799 = l0 + int32(3431)
	v800 = v754
	v802 = v791
	v803 = v754
	v805 = v754
	goto L92
L90:
	;
	v2232 = v695 + int32(-1)
	goto L88
L91:
	;
	v1298 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9])))
	if v1298 == int32(0) {
		goto L135
	} else {
		goto L136
	}
L92:
	;
	v832 = v803
	v835 = v796
	v836 = v797
	v837 = v798
	v838 = v799
	v839 = v800
	v840 = v802
	v841 = int32(0)
	goto L94
L93:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = v1217
	goto L91
L94:
	;
	v861 = v832
	v871 = int32(0)
	v872 = v839
	goto L96
L95:
	;
	v1241 = int32(1056)
	v1243 = int32(264)
	v1254 = v805 + int32(1)
	if v1254 != int32(4) {
		v796 = v796 + v1241
		v797 = v797 + v1243
		v798 = v798 + v1241
		v799 = v799 + v1243
		v800 = v800 + v1243
		v802 = v802 + v1241
		v803 = v1217
		v805 = v1254
		goto L92
	} else {
		goto L133
	}
L96:
	;
	v882 = *(*int32)(unsafe.Add(mBase, uint32(v840+v871)))
	v884 = int32(base.Ui32(v882) >> (uint(int32(16)) % 32))
	v885 = m.G80
	v887 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v885+v872))))
	v888 = m.G81
	v890 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v888+v872))))
	v892 = v882 & int32(_a_F_VP8EncTokenLoop_7)
	if v892 != 0 {
		goto L99
	} else {
		goto L100
	}
L97:
	;
	v986 = v969
	v996 = v837
	v1002 = int32(0)
	goto L108
L98:
	;
	v900 = m.G24
	v903 = v884 - v892
	v904 = int32(1)
	v907 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v900+v887<<(uint(v904)%32)))))
	v909 = int32(255)
	v914 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v900+(v887^v909)<<(uint(v904)%32)))))
	v920 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v900+v890<<(uint(v904)%32)))))
	v923 = v899 & v909
	v929 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v900+(v923^v909)<<(uint(v904)%32)))))
	v934 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v900+v923<<(uint(v904)%32)))))
	v942 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v900+(v890^v909)<<(uint(v904)%32)))))
	goto L101
L99:
	;
	v894 = int32(255)
	v897 = base.I32_div_u_s(v892*v894, v884)
	v899 = v894 - v897
	goto L98
L100:
	;
	v899 = int32(255)
	goto L98
L101:
	;
	goto L103
L103:
	;
	if v903*v907+v892*v914+v920 <= v892*v929+v903*v934+v942+int32(2048) {
		goto L105
	} else {
		goto L106
	}
L104:
	;
	v974 = v871 + int32(4)
	if v974 != int32(44) {
		v861 = v969
		v871 = v974
		v872 = v872 + int32(1)
		goto L96
	} else {
		goto L107
	}
L105:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v791+v872+int32(-1056)))) = uint8(v887)
	v969 = v861
	goto L104
L106:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v791+v872+int32(-1056)))) = uint8(v899)
	v969 = v861 | base.B2i32(v899 != v887)
	goto L104
L107:
	;
	goto L97
L108:
	;
	v1006 = *(*int32)(unsafe.Add(mBase, uint32(v996)))
	v1008 = int32(base.Ui32(v1006) >> (uint(int32(16)) % 32))
	v1009 = m.G80
	v1010 = v839 + v1002
	v1012 = int32(11)
	v1014 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1009+v1010+v1012))))
	v1015 = m.G81
	v1019 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1015+v1010+v1012))))
	v1021 = v1006 & int32(_a_F_VP8EncTokenLoop_7)
	if v1021 != 0 {
		goto L111
	} else {
		goto L112
	}
L109:
	;
	v1110 = v1093
	v1120 = v835
	v1126 = int32(0)
	goto L120
L110:
	;
	v1029 = m.G24
	v1032 = v1008 - v1021
	v1033 = int32(1)
	v1036 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1029+v1014<<(uint(v1033)%32)))))
	v1038 = int32(255)
	v1043 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1029+(v1014^v1038)<<(uint(v1033)%32)))))
	v1049 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1029+v1019<<(uint(v1033)%32)))))
	v1052 = v1028 & v1038
	v1058 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1029+(v1052^v1038)<<(uint(v1033)%32)))))
	v1063 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1029+v1052<<(uint(v1033)%32)))))
	v1071 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1029+(v1019^v1038)<<(uint(v1033)%32)))))
	goto L113
L111:
	;
	v1023 = int32(255)
	v1026 = base.I32_div_u_s(v1021*v1023, v1008)
	v1028 = v1023 - v1026
	goto L110
L112:
	;
	v1028 = int32(255)
	goto L110
L113:
	;
	goto L115
L115:
	;
	if v1021*v1058+v1032*v1063+v1071+int32(2048) < v1032*v1036+v1021*v1043+v1049 {
		goto L117
	} else {
		goto L118
	}
L116:
	;
	v1098 = v1002 + int32(1)
	if v1098 != int32(11) {
		v986 = v1093
		v996 = v996 + int32(4)
		v1002 = v1098
		goto L108
	} else {
		goto L119
	}
L117:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v838+v1002))) = uint8(v1028)
	v1093 = v986 | base.B2i32(v1028 != v1014)
	goto L116
L118:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v838+v1002))) = uint8(v1014)
	v1093 = v986
	goto L116
L119:
	;
	goto L109
L120:
	;
	v1130 = *(*int32)(unsafe.Add(mBase, uint32(v1120)))
	v1132 = int32(base.Ui32(v1130) >> (uint(int32(16)) % 32))
	v1133 = m.G80
	v1134 = v839 + v1126
	v1136 = int32(22)
	v1138 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1133+v1134+v1136))))
	v1139 = m.G81
	v1143 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1139+v1134+v1136))))
	v1145 = v1130 & int32(_a_F_VP8EncTokenLoop_7)
	if v1145 != 0 {
		goto L123
	} else {
		goto L124
	}
L121:
	;
	v1225 = int32(132)
	v1227 = int32(33)
	v1238 = v841 + int32(1)
	if v1238 != int32(8) {
		v832 = v1217
		v835 = v835 + v1225
		v836 = v836 + v1227
		v837 = v837 + v1225
		v838 = v838 + v1227
		v839 = v839 + v1227
		v840 = v840 + v1225
		v841 = v1238
		goto L94
	} else {
		goto L132
	}
L122:
	;
	v1153 = m.G24
	v1156 = v1132 - v1145
	v1157 = int32(1)
	v1160 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1153+v1138<<(uint(v1157)%32)))))
	v1162 = int32(255)
	v1167 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1153+(v1138^v1162)<<(uint(v1157)%32)))))
	v1173 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1153+v1143<<(uint(v1157)%32)))))
	v1176 = v1152 & v1162
	v1182 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1153+(v1176^v1162)<<(uint(v1157)%32)))))
	v1187 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1153+v1176<<(uint(v1157)%32)))))
	v1195 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1153+(v1143^v1162)<<(uint(v1157)%32)))))
	goto L125
L123:
	;
	v1147 = int32(255)
	v1150 = base.I32_div_u_s(v1145*v1147, v1132)
	v1152 = v1147 - v1150
	goto L122
L124:
	;
	v1152 = int32(255)
	goto L122
L125:
	;
	goto L127
L127:
	;
	if v1145*v1182+v1156*v1187+v1195+int32(2048) < v1156*v1160+v1145*v1167+v1173 {
		goto L129
	} else {
		goto L130
	}
L128:
	;
	v1222 = v1126 + int32(1)
	if v1222 != int32(11) {
		v1110 = v1217
		v1120 = v1120 + int32(4)
		v1126 = v1222
		goto L120
	} else {
		goto L131
	}
L129:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v836+v1126))) = uint8(v1152)
	v1217 = v1110 | base.B2i32(v1152 != v1138)
	goto L128
L130:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v836+v1126))) = uint8(v1138)
	v1217 = v1110
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
	v2232 = v262
	goto L88
L135:
	;
	goto L134
L136:
	;
	v1313 = m.G23
	v1314 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+15)))
	v1315 = int32(408)
	v1317 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+14)))
	v1320 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+13)))
	v1323 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+12)))
	v1326 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+11)))
	v1329 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+10)))
	v1332 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+9)))
	v1335 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+8)))
	v1338 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+7)))
	v1341 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+6)))
	v1344 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+5)))
	v1347 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+4)))
	v1350 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+3)))
	v1353 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+2)))
	v1356 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313)+1)))
	v1359 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1313))))
	v1364 = l0 + int32(3444)
	v1365 = l0 + int32(3433)
	v1366 = l0 + int32(3422)
	v1387 = int32(0)
	goto L137
L137:
	;
	v1406 = l0 + int32(_a_F_VP8EncTokenLoop_8) + v1387*int32(3264)
	v1438 = v1364
	v1439 = v1365
	v1440 = v1366
	v1441 = int32(0)
	goto L139
L138:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = int32(0)
	goto L135
L139:
	;
	v1454 = l0 + int32(3420) + v1387*int32(264) + v1441*int32(33)
	v1455 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1454)+1)))
	v1456 = int32(1)
	v1459 = v1406 + v1441*int32(408)
	v1460 = m.G24
	v1464 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1460+v1455<<(uint(v1456)%32)))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1459))) = uint16(v1464)
	v1471 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1460+(v1455^int32(255))<<(uint(v1456)%32)))))
	v1504 = v1456
	goto L141
L140:
	;
	v2050 = l0 + int32(_a_F_VP8EncTokenLoop_9) + v1387*int32(192)
	v2051 = v1406 + v1314*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+180)) = v2051
	v2053 = v1406 + v1317*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+168)) = v2053
	v2055 = v1406 + v1320*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+156)) = v2055
	v2057 = v1406 + v1323*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+144)) = v2057
	v2059 = v1406 + v1326*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+132)) = v2059
	v2061 = v1406 + v1329*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+120)) = v2061
	v2063 = v1406 + v1332*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+108)) = v2063
	v2065 = v1406 + v1335*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+96)) = v2065
	v2067 = v1406 + v1338*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+84)) = v2067
	v2069 = v1406 + v1341*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+72)) = v2069
	v2071 = v1406 + v1344*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+60)) = v2071
	v2073 = v1406 + v1347*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+48)) = v2073
	v2075 = v1406 + v1350*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+36)) = v2075
	v2077 = v1406 + v1353*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+24)) = v2077
	v2079 = v1406 + v1356*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+12)) = v2079
	v2081 = v1406 + v1359*v1315
	*(*int32)(unsafe.Add(mBase, uint32(v2050))) = v2081
	v2083 = int32(272)
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+188)) = v2051 + v2083
	v2086 = int32(136)
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+184)) = v2051 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+176)) = v2053 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+172)) = v2053 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+164)) = v2055 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+160)) = v2055 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+152)) = v2057 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+148)) = v2057 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+140)) = v2059 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+136)) = v2059 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+128)) = v2061 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+124)) = v2061 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+116)) = v2063 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+112)) = v2063 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+104)) = v2065 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+100)) = v2065 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+92)) = v2067 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+88)) = v2067 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+80)) = v2069 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+76)) = v2069 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+68)) = v2071 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+64)) = v2071 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+56)) = v2073 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+52)) = v2073 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+44)) = v2075 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+40)) = v2075 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+32)) = v2077 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+28)) = v2077 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+20)) = v2079 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+16)) = v2079 + v2086
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+8)) = v2081 + v2083
	*(*int32)(unsafe.Add(mBase, uint32(v2050)+4)) = v2081 + v2086
	v2179 = int32(264)
	v2186 = v1387 + int32(1)
	if v2186 != int32(4) {
		v1364 = v1364 + v2179
		v1365 = v1365 + v2179
		v1366 = v1366 + v2179
		v1387 = v2186
		goto L137
	} else {
		goto L175
	}
L141:
	;
	v1515 = m.G1
	v1520 = v1504<<(uint(int32(2))%32) + (v1515 + int32(_a_F_VP8EncTokenLoop_10)) + int32(-4)
	v1521 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1520))))
	if v1521 != 0 {
		goto L144
	} else {
		goto L145
	}
L142:
	;
	v1644 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1454)+12)))
	v1645 = int32(1)
	v1646 = m.G24
	v1650 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1646+v1644<<(uint(v1645)%32)))))
	v1651 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1454)+11)))
	v1652 = int32(255)
	v1657 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1646+(v1651^v1652)<<(uint(v1645)%32)))))
	v1658 = v1650 + v1657
	*(*uint16)(unsafe.Add(mBase, uint32(v1459)+136)) = uint16(v1658)
	v1665 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1646+(v1644^v1652)<<(uint(v1645)%32)))))
	v1701 = v1645
	goto L152
L143:
	;
	v1635 = int32(1)
	v1638 = v1471 + v1630
	*(*uint16)(unsafe.Add(mBase, uint32(v1459+v1504<<(uint(v1635)%32)))) = uint16(v1638)
	v1641 = v1504 + v1635
	if v1641 != int32(68) {
		v1504 = v1641
		goto L141
	} else {
		goto L151
	}
L144:
	;
	v1523 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1520)+2)))
	v1532 = v1521
	v1559 = v1523
	v1561 = int32(0)
	v1562 = v1440
	goto L146
L145:
	;
	v1630 = int32(0)
	goto L143
L146:
	;
	if v1532&int32(1) == int32(0) {
		v1584 = v1561
		goto L148
	} else {
		goto L149
	}
L147:
	;
	v1630 = v1584
	goto L143
L148:
	;
	v1586 = int32(1)
	if base.Ui32(v1586) < base.Ui32(v1532) {
		v1532 = int32(base.Ui32(v1532) >> (uint(v1586) % 32))
		v1559 = int32(base.Ui32(v1559) >> (uint(v1586) % 32))
		v1561 = v1584
		v1562 = v1562 + v1586
		goto L146
	} else {
		goto L150
	}
L149:
	;
	v1570 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1562))))
	v1571 = m.G24
	v1573 = int32(1)
	v1582 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1571+(v1570^(int32(0)-v1559&v1573)&int32(255))<<(uint(v1573)%32)))))
	v1584 = v1561 + v1582
	goto L148
L150:
	;
	goto L147
L151:
	;
	goto L142
L152:
	;
	v1712 = m.G1
	v1717 = v1701<<(uint(int32(2))%32) + (v1712 + int32(_a_F_VP8EncTokenLoop_10)) + int32(-4)
	v1718 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1717))))
	if v1718 != 0 {
		goto L155
	} else {
		goto L156
	}
L153:
	;
	v1841 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1454)+23)))
	v1842 = int32(1)
	v1843 = m.G24
	v1847 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1843+v1841<<(uint(v1842)%32)))))
	v1848 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1454)+22)))
	v1849 = int32(255)
	v1854 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1843+(v1848^v1849)<<(uint(v1842)%32)))))
	v1855 = v1847 + v1854
	*(*uint16)(unsafe.Add(mBase, uint32(v1459)+272)) = uint16(v1855)
	v1862 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1843+(v1841^v1849)<<(uint(v1842)%32)))))
	v1898 = v1842
	goto L163
L154:
	;
	v1832 = int32(1)
	v1835 = v1657 + v1665 + v1827
	*(*uint16)(unsafe.Add(mBase, uint32(v1459+int32(136)+v1701<<(uint(v1832)%32)))) = uint16(v1835)
	v1838 = v1701 + v1832
	if v1838 != int32(68) {
		v1701 = v1838
		goto L152
	} else {
		goto L162
	}
L155:
	;
	v1720 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1717)+2)))
	v1729 = v1718
	v1756 = v1720
	v1758 = int32(0)
	v1759 = v1439
	goto L157
L156:
	;
	v1827 = int32(0)
	goto L154
L157:
	;
	if v1729&int32(1) == int32(0) {
		v1781 = v1758
		goto L159
	} else {
		goto L160
	}
L158:
	;
	v1827 = v1781
	goto L154
L159:
	;
	v1783 = int32(1)
	if base.Ui32(v1783) < base.Ui32(v1729) {
		v1729 = int32(base.Ui32(v1729) >> (uint(v1783) % 32))
		v1756 = int32(base.Ui32(v1756) >> (uint(v1783) % 32))
		v1758 = v1781
		v1759 = v1759 + v1783
		goto L157
	} else {
		goto L161
	}
L160:
	;
	v1767 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1759))))
	v1768 = m.G24
	v1770 = int32(1)
	v1779 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1768+(v1767^(int32(0)-v1756&v1770)&int32(255))<<(uint(v1770)%32)))))
	v1781 = v1758 + v1779
	goto L159
L161:
	;
	goto L158
L162:
	;
	goto L153
L163:
	;
	v1909 = m.G1
	v1914 = v1898<<(uint(int32(2))%32) + (v1909 + int32(_a_F_VP8EncTokenLoop_10)) + int32(-4)
	v1915 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1914))))
	if v1915 != 0 {
		goto L166
	} else {
		goto L167
	}
L164:
	;
	v2038 = int32(33)
	v2045 = v1441 + int32(1)
	if v2045 != int32(8) {
		v1438 = v1438 + v2038
		v1439 = v1439 + v2038
		v1440 = v1440 + v2038
		v1441 = v2045
		goto L139
	} else {
		goto L174
	}
L165:
	;
	v2029 = int32(1)
	v2032 = v1854 + v1862 + v2024
	*(*uint16)(unsafe.Add(mBase, uint32(v1459+int32(272)+v1898<<(uint(v2029)%32)))) = uint16(v2032)
	v2035 = v1898 + v2029
	if v2035 != int32(68) {
		v1898 = v2035
		goto L163
	} else {
		goto L173
	}
L166:
	;
	v1917 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1914)+2)))
	v1926 = v1915
	v1953 = v1917
	v1955 = int32(0)
	v1956 = v1438
	goto L168
L167:
	;
	v2024 = int32(0)
	goto L165
L168:
	;
	if v1926&int32(1) == int32(0) {
		v1978 = v1955
		goto L170
	} else {
		goto L171
	}
L169:
	;
	v2024 = v1978
	goto L165
L170:
	;
	v1980 = int32(1)
	if base.Ui32(v1980) < base.Ui32(v1926) {
		v1926 = int32(base.Ui32(v1926) >> (uint(v1980) % 32))
		v1953 = int32(base.Ui32(v1953) >> (uint(v1980) % 32))
		v1955 = v1978
		v1956 = v1956 + v1980
		goto L168
	} else {
		goto L172
	}
L171:
	;
	v1964 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1956))))
	v1965 = m.G24
	v1967 = int32(1)
	v1976 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1965+(v1964^(int32(0)-v1953&v1967)&int32(255))<<(uint(v1967)%32)))))
	v1978 = v1955 + v1976
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
	v2331 = *(*int32)(unsafe.Add(mBase, uint32(v64)+908))
	v2332 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2331))))
	if v2332&int32(3) != int32(1) {
		goto L178
	} else {
		goto L179
	}
L177:
	;
	v2424 = m.G25
	v2425 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1008))
	v2426 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2428 = v64 + int32(_a_F_VP8EncTokenLoop_11)
	v2429 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2429].(func(*base.Module, int32, int32))(m, v64+int32(72), v2428)
	mBase = m.M
	v2434 = F_VP8RecordCoeffTokens(m, v2425+v2426, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2434
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2434
	v2437 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2440 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2440].(func(*base.Module, int32, int32))(m, v64+int32(104), v2428)
	mBase = m.M
	v2445 = F_VP8RecordCoeffTokens(m, v2437+v2434, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2445
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2445
	v2448 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2451 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2451].(func(*base.Module, int32, int32))(m, v64+int32(136), v2428)
	mBase = m.M
	v2456 = F_VP8RecordCoeffTokens(m, v2448+v2445, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2456
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2456
	v2459 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2462 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2462].(func(*base.Module, int32, int32))(m, v64+int32(168), v2428)
	mBase = m.M
	v2467 = F_VP8RecordCoeffTokens(m, v2459+v2456, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2467
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1008)) = v2467
	v2470 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2471 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1012))
	v2474 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2474].(func(*base.Module, int32, int32))(m, v64+int32(200), v2428)
	mBase = m.M
	v2479 = F_VP8RecordCoeffTokens(m, v2471+v2470, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2479
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2479
	v2482 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2485 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2485].(func(*base.Module, int32, int32))(m, v64+int32(232), v2428)
	mBase = m.M
	v2490 = F_VP8RecordCoeffTokens(m, v2482+v2479, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2490
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2490
	v2493 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2496 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2496].(func(*base.Module, int32, int32))(m, v64+int32(264), v2428)
	mBase = m.M
	v2501 = F_VP8RecordCoeffTokens(m, v2493+v2490, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2501
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2501
	v2504 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2507 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2507].(func(*base.Module, int32, int32))(m, v64+int32(296), v2428)
	mBase = m.M
	v2512 = F_VP8RecordCoeffTokens(m, v2504+v2501, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2512
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1012)) = v2512
	v2515 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2516 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1016))
	v2519 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2519].(func(*base.Module, int32, int32))(m, v64+int32(328), v2428)
	mBase = m.M
	v2524 = F_VP8RecordCoeffTokens(m, v2516+v2515, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2524
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2524
	v2527 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2530 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2530].(func(*base.Module, int32, int32))(m, v64+int32(360), v2428)
	mBase = m.M
	v2535 = F_VP8RecordCoeffTokens(m, v2527+v2524, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2535
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2535
	v2538 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2541 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2541].(func(*base.Module, int32, int32))(m, v64+int32(392), v2428)
	mBase = m.M
	v2546 = F_VP8RecordCoeffTokens(m, v2538+v2535, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2546
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2546
	v2549 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2552 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2552].(func(*base.Module, int32, int32))(m, v64+int32(424), v2428)
	mBase = m.M
	v2557 = F_VP8RecordCoeffTokens(m, v2549+v2546, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2557
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1016)) = v2557
	v2560 = *(*int32)(unsafe.Add(mBase, uint32(v64)+972))
	v2561 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1020))
	v2564 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2564].(func(*base.Module, int32, int32))(m, v64+int32(456), v2428)
	mBase = m.M
	v2569 = F_VP8RecordCoeffTokens(m, v2561+v2560, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+972)) = v2569
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2569
	v2572 = *(*int32)(unsafe.Add(mBase, uint32(v64)+976))
	v2575 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2575].(func(*base.Module, int32, int32))(m, v64+int32(488), v2428)
	mBase = m.M
	v2580 = F_VP8RecordCoeffTokens(m, v2572+v2569, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+976)) = v2580
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2580
	v2583 = *(*int32)(unsafe.Add(mBase, uint32(v64)+980))
	v2586 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2586].(func(*base.Module, int32, int32))(m, v64+int32(520), v2428)
	mBase = m.M
	v2591 = F_VP8RecordCoeffTokens(m, v2583+v2580, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+980)) = v2591
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2591
	v2594 = *(*int32)(unsafe.Add(mBase, uint32(v64)+984))
	v2597 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2597].(func(*base.Module, int32, int32))(m, v64+int32(552), v2428)
	mBase = m.M
	v2602 = F_VP8RecordCoeffTokens(m, v2594+v2591, v2428, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+984)) = v2602
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1020)) = v2602
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2236 + int32(_a_F_VP8EncTokenLoop_12)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2236 + int32(_a_F_VP8EncTokenLoop_13)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2236 + int32(3948)
	goto L183
L178:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2236 + int32(_a_F_VP8EncTokenLoop_14)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2236 + int32(_a_F_VP8EncTokenLoop_15)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2236 + int32(_a_F_VP8EncTokenLoop_16)
	goto L182
L179:
	;
	v2337 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1040))
	v2338 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1004))
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2236 + int32(_a_F_VP8EncTokenLoop_17)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2236 + int32(_a_F_VP8EncTokenLoop_18)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2236 + int32(3684)
	goto L180
L180:
	;
	v2364 = v64 + int32(_a_F_VP8EncTokenLoop_11)
	v2365 = m.G25
	v2366 = *(*int32)(unsafe.Add(mBase, uint32(v2365)))
	m.T0[v2366].(func(*base.Module, int32, int32))(m, v64+v320, v2364)
	mBase = m.M
	v2371 = F_VP8RecordCoeffTokens(m, v2337+v2338, v2364, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1004)) = v2371
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1040)) = v2371
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[10]))) = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[11]))) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[12]))) = v2236 + int32(_a_F_VP8EncTokenLoop_9)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[13]))) = v2236 + int32(_a_F_VP8EncTokenLoop_2)
	*(*int32)(unsafe.Add(mBase, uint32(v64)+uint32(_c_F_VP8EncTokenLoop[14]))) = v2236 + int32(3420)
	goto L181
L181:
	;
	goto L177
L182:
	;
	goto L177
L183:
	;
	v2629 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1024))
	v2630 = *(*int32)(unsafe.Add(mBase, uint32(v64)+988))
	v2632 = v64 + int32(_a_F_VP8EncTokenLoop_11)
	v2633 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2633].(func(*base.Module, int32, int32))(m, v64+int32(584), v2632)
	mBase = m.M
	v2638 = F_VP8RecordCoeffTokens(m, v2629+v2630, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+988)) = v2638
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1024)) = v2638
	v2641 = *(*int32)(unsafe.Add(mBase, uint32(v64)+992))
	v2644 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2644].(func(*base.Module, int32, int32))(m, v64+int32(616), v2632)
	mBase = m.M
	v2649 = F_VP8RecordCoeffTokens(m, v2641+v2638, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+992)) = v2649
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1024)) = v2649
	v2652 = *(*int32)(unsafe.Add(mBase, uint32(v64)+988))
	v2653 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1028))
	v2656 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2656].(func(*base.Module, int32, int32))(m, v64+int32(648), v2632)
	mBase = m.M
	v2661 = F_VP8RecordCoeffTokens(m, v2653+v2652, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+988)) = v2661
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1028)) = v2661
	v2664 = *(*int32)(unsafe.Add(mBase, uint32(v64)+992))
	v2667 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2667].(func(*base.Module, int32, int32))(m, v64+int32(680), v2632)
	mBase = m.M
	v2672 = F_VP8RecordCoeffTokens(m, v2664+v2661, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+992)) = v2672
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1028)) = v2672
	v2675 = *(*int32)(unsafe.Add(mBase, uint32(v64)+996))
	v2676 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1032))
	v2679 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2679].(func(*base.Module, int32, int32))(m, v64+int32(712), v2632)
	mBase = m.M
	v2684 = F_VP8RecordCoeffTokens(m, v2676+v2675, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+996)) = v2684
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1032)) = v2684
	v2687 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1000))
	v2690 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2690].(func(*base.Module, int32, int32))(m, v64+int32(744), v2632)
	mBase = m.M
	v2695 = F_VP8RecordCoeffTokens(m, v2687+v2684, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1000)) = v2695
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1032)) = v2695
	v2698 = *(*int32)(unsafe.Add(mBase, uint32(v64)+996))
	v2699 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1036))
	v2702 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2702].(func(*base.Module, int32, int32))(m, v64+int32(776), v2632)
	mBase = m.M
	v2707 = F_VP8RecordCoeffTokens(m, v2699+v2698, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+996)) = v2707
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1036)) = v2707
	v2710 = *(*int32)(unsafe.Add(mBase, uint32(v64)+1000))
	v2713 = *(*int32)(unsafe.Add(mBase, uint32(v2424)))
	m.T0[v2713].(func(*base.Module, int32, int32))(m, v64+int32(808), v2632)
	mBase = m.M
	v2718 = F_VP8RecordCoeffTokens(m, v2710+v2707, v2632, v266)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1000)) = v2718
	*(*int32)(unsafe.Add(mBase, uint32(v64)+1036)) = v2718
	v2722 = v64 + int32(880)
	v2723 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+40))
	v2724 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+96))
	v2727 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+92))
	v2731 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+100))
	v2735 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+104))
	v2739 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+108))
	v2743 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+112))
	v2747 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+116))
	v2751 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+120))
	v2755 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+124))
	v2759 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+128))
	v2763 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+132))
	v2767 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+136))
	v2771 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+144))
	v2775 = *(*int32)(unsafe.Add(mBase, uint32(v2722)+152))
	*(*int32)(unsafe.Add(mBase, uint32(v2723))) = v2724<<(uint(int32(13))%32) | v2727<<(uint(int32(12))%32) | v2731<<(uint(int32(14))%32) | v2735<<(uint(int32(15))%32) | v2739<<(uint(int32(18))%32) | v2743<<(uint(int32(19))%32) | v2747<<(uint(int32(22))%32) | v2751<<(uint(int32(23))%32) | v2755<<(uint(int32(24))%32) | v2759<<(uint(int32(3))%32) | v2763<<(uint(int32(7))%32) | v2767<<(uint(int32(11))%32) | v2771<<(uint(int32(17))%32) | v2775<<(uint(int32(21))%32)
	goto L184
L184:
	;
	v2780 = *(*int32)(unsafe.Add(mBase, uint32(l0)+364))
	if v2780 == int32(0) {
		goto L185
	} else {
		goto L186
	}
L185:
	;
	v2789 = *(*int64)(unsafe.Add(mBase, uint32(v64)))
	v2790 = *(*int64)(unsafe.Add(mBase, uint32(v64)+16))
	if v399 != 0 {
		goto L191
	} else {
		goto L192
	}
L186:
	;
	v2783 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v2785 = *(*int32)(unsafe.Add(mBase, uint32(v2783)+92))
	if v2785 != 0 {
		goto L188
	} else {
		goto L189
	}
L187:
	;
	v4834 = int32(0)
	goto L39
L188:
	;
	goto L187
L189:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2783)+92)) = int32(1)
	goto L188
L190:
	;
	v3135 = v2789 + v741
	v3136 = v2790 + v740
	v3138 = v64 + int32(880)
	v3143 = *(*int32)(unsafe.Add(mBase, uint32(v3138)))
	v3145 = v3143 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v3138))) = v3145
	v3147 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+24))
	v3148 = *(*int32)(unsafe.Add(mBase, uint32(v3147)+40))
	if v3145 != v3148 {
		goto L213
	} else {
		goto L214
	}
L191:
	;
	v2943 = v64 + int32(880)
	F_StoreSideInfo(m, v2943)
	mBase = m.M
	F_VP8StoreFilterStats(m, v2943)
	mBase = m.M
	F_VP8IteratorExport(m, v2943)
	mBase = m.M
	v2956 = int32(1)
	if v517 == int32(0) {
		v2979 = v2956
		goto L199
	} else {
		goto L200
	}
L192:
	;
	v2792 = v64 + int32(880)
	v2796 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+12))
	v2797 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+4))
	v2798 = *(*int32)(unsafe.Add(mBase, uint32(v2792)))
	v2799 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+24))
	v2800 = *(*int32)(unsafe.Add(mBase, uint32(v2799)+40))
	if v2800+int32(-1) <= v2798 {
		goto L194
	} else {
		goto L195
	}
L193:
	;
	goto L190
L194:
	;
	v2918 = *(*int32)(unsafe.Add(mBase, uint32(v2799)+44))
	if v2918+int32(-1) <= v2797 {
		goto L196
	} else {
		goto L197
	}
L195:
	;
	v2804 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2805 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2804))) = uint8(v2805)
	v2807 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2808 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+47)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2807)+1)) = uint8(v2808)
	v2810 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2811 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+79)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2810)+2)) = uint8(v2811)
	v2813 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2814 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+111)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2813)+3)) = uint8(v2814)
	v2816 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2817 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+143)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2816)+4)) = uint8(v2817)
	v2819 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2820 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+175)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2819)+5)) = uint8(v2820)
	v2822 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2823 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+207)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2822)+6)) = uint8(v2823)
	v2825 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2826 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+239)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2825)+7)) = uint8(v2826)
	v2828 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2829 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+271)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2828)+8)) = uint8(v2829)
	v2831 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2832 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+303)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2831)+9)) = uint8(v2832)
	v2834 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2835 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+335)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2834)+10)) = uint8(v2835)
	v2837 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2838 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+367)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2837)+11)) = uint8(v2838)
	v2840 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2841 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+399)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2840)+12)) = uint8(v2841)
	v2843 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2844 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+431)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2843)+13)) = uint8(v2844)
	v2846 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2847 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+463)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2846)+14)) = uint8(v2847)
	v2849 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2850 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+495)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2849)+15)) = uint8(v2850)
	v2852 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2853 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+23)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2852))) = uint8(v2853)
	v2855 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2856 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+31)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2855))) = uint8(v2856)
	v2858 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2859 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+55)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2858)+1)) = uint8(v2859)
	v2861 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2862 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+63)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2861)+1)) = uint8(v2862)
	v2864 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2865 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+87)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2864)+2)) = uint8(v2865)
	v2867 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2868 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+95)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2867)+2)) = uint8(v2868)
	v2870 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2871 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+119)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2870)+3)) = uint8(v2871)
	v2873 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2874 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+127)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2873)+3)) = uint8(v2874)
	v2876 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2877 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+151)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2876)+4)) = uint8(v2877)
	v2879 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2880 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+159)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2879)+4)) = uint8(v2880)
	v2882 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2883 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+183)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2882)+5)) = uint8(v2883)
	v2885 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2886 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+191)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2885)+5)) = uint8(v2886)
	v2888 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2889 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+215)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2888)+6)) = uint8(v2889)
	v2891 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2892 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+223)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2891)+6)) = uint8(v2892)
	v2894 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2895 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+247)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2894)+7)) = uint8(v2895)
	v2897 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2898 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2796)+255)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2897)+7)) = uint8(v2898)
	v2900 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+308))
	v2901 = int32(-1)
	v2903 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+320))
	v2904 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2903)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2900+v2901))) = uint8(v2904)
	v2906 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+312))
	v2909 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+324))
	v2910 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2909)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2906+v2901))) = uint8(v2910)
	v2912 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+316))
	v2915 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+324))
	v2916 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2915)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2912+v2901))) = uint8(v2916)
	goto L194
L196:
	;
	goto L193
L197:
	;
	v2922 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+320))
	v2923 = *(*int64)(unsafe.Add(mBase, uint32(v2796)+480))
	*(*int64)(unsafe.Add(mBase, uint32(v2922))) = v2923
	v2925 = int32(8)
	v2929 = *(*int64)(unsafe.Add(mBase, uint32(v2796+int32(488))))
	*(*int64)(unsafe.Add(mBase, uint32(v2922+v2925))) = v2929
	v2931 = *(*int32)(unsafe.Add(mBase, uint32(v2792)+324))
	v2932 = *(*int64)(unsafe.Add(mBase, uint32(v2796)+240))
	*(*int64)(unsafe.Add(mBase, uint32(v2931))) = v2932
	v2938 = *(*int64)(unsafe.Add(mBase, uint32(v2796+int32(248))))
	*(*int64)(unsafe.Add(mBase, uint32(v2931+v2925))) = v2938
	goto L196
L198:
	;
	v2983 = v64 + int32(880)
	v2987 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+12))
	v2988 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+4))
	v2989 = *(*int32)(unsafe.Add(mBase, uint32(v2983)))
	v2990 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+24))
	v2991 = *(*int32)(unsafe.Add(mBase, uint32(v2990)+40))
	if v2991+int32(-1) <= v2989 {
		goto L206
	} else {
		goto L207
	}
L199:
	;
	goto L198
L200:
	;
	v2959 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+24))
	v2960 = *(*int32)(unsafe.Add(mBase, uint32(v2959)+4))
	v2961 = *(*int32)(unsafe.Add(mBase, uint32(v2960)+96))
	if v2961 == int32(0) {
		v2979 = v2956
		goto L199
	} else {
		goto L201
	}
L201:
	;
	v2964 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+292))
	if int32(0) < v2964 {
		goto L203
	} else {
		goto L204
	}
L202:
	;
	v2977 = F_WebPReportProgress(m, v2960, v2974, v2959+int32(368))
	mBase = m.M
	v2979 = v2977
	goto L199
L203:
	;
	v2968 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+288))
	v2971 = base.I32_div_s((v2964-v2968)*v517, v2964)
	v2972 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+296))
	v2974 = v2971 + v2972
	goto L202
L204:
	;
	v2967 = *(*int32)(unsafe.Add(mBase, uint32(v2943)+296))
	v2974 = v2967
	goto L202
L205:
	;
	if v2979 != 0 {
		goto L190
	} else {
		goto L210
	}
L206:
	;
	v3109 = *(*int32)(unsafe.Add(mBase, uint32(v2990)+44))
	if v3109+int32(-1) <= v2988 {
		goto L208
	} else {
		goto L209
	}
L207:
	;
	v2995 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v2996 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2995))) = uint8(v2996)
	v2998 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v2999 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+47)))
	*(*uint8)(unsafe.Add(mBase, uint32(v2998)+1)) = uint8(v2999)
	v3001 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3002 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+79)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3001)+2)) = uint8(v3002)
	v3004 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3005 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+111)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3004)+3)) = uint8(v3005)
	v3007 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3008 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+143)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3007)+4)) = uint8(v3008)
	v3010 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3011 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+175)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3010)+5)) = uint8(v3011)
	v3013 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3014 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+207)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3013)+6)) = uint8(v3014)
	v3016 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3017 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+239)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3016)+7)) = uint8(v3017)
	v3019 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3020 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+271)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3019)+8)) = uint8(v3020)
	v3022 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3023 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+303)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3022)+9)) = uint8(v3023)
	v3025 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3026 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+335)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3025)+10)) = uint8(v3026)
	v3028 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3029 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+367)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3028)+11)) = uint8(v3029)
	v3031 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3032 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+399)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3031)+12)) = uint8(v3032)
	v3034 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3035 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+431)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3034)+13)) = uint8(v3035)
	v3037 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3038 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+463)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3037)+14)) = uint8(v3038)
	v3040 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3041 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+495)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3040)+15)) = uint8(v3041)
	v3043 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3044 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+23)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3043))) = uint8(v3044)
	v3046 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3047 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+31)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3046))) = uint8(v3047)
	v3049 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3050 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+55)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3049)+1)) = uint8(v3050)
	v3052 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3053 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+63)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3052)+1)) = uint8(v3053)
	v3055 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3056 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+87)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3055)+2)) = uint8(v3056)
	v3058 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3059 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+95)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3058)+2)) = uint8(v3059)
	v3061 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3062 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+119)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3061)+3)) = uint8(v3062)
	v3064 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3065 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+127)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3064)+3)) = uint8(v3065)
	v3067 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3068 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+151)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3067)+4)) = uint8(v3068)
	v3070 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3071 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+159)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3070)+4)) = uint8(v3071)
	v3073 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3074 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+183)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3073)+5)) = uint8(v3074)
	v3076 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3077 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+191)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3076)+5)) = uint8(v3077)
	v3079 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3080 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+215)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3079)+6)) = uint8(v3080)
	v3082 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3083 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+223)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3082)+6)) = uint8(v3083)
	v3085 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3086 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+247)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3085)+7)) = uint8(v3086)
	v3088 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3089 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2987)+255)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3088)+7)) = uint8(v3089)
	v3091 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+308))
	v3092 = int32(-1)
	v3094 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+320))
	v3095 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3094)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3091+v3092))) = uint8(v3095)
	v3097 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+312))
	v3100 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+324))
	v3101 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3100)+7)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3097+v3092))) = uint8(v3101)
	v3103 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+316))
	v3106 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+324))
	v3107 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3106)+15)))
	*(*uint8)(unsafe.Add(mBase, uint32(v3103+v3092))) = uint8(v3107)
	goto L206
L208:
	;
	goto L205
L209:
	;
	v3113 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+320))
	v3114 = *(*int64)(unsafe.Add(mBase, uint32(v2987)+480))
	*(*int64)(unsafe.Add(mBase, uint32(v3113))) = v3114
	v3116 = int32(8)
	v3120 = *(*int64)(unsafe.Add(mBase, uint32(v2987+int32(488))))
	*(*int64)(unsafe.Add(mBase, uint32(v3113+v3116))) = v3120
	v3122 = *(*int32)(unsafe.Add(mBase, uint32(v2983)+324))
	v3123 = *(*int64)(unsafe.Add(mBase, uint32(v2987)+240))
	*(*int64)(unsafe.Add(mBase, uint32(v3122))) = v3123
	v3129 = *(*int64)(unsafe.Add(mBase, uint32(v2987+int32(248))))
	*(*int64)(unsafe.Add(mBase, uint32(v3122+v3116))) = v3129
	goto L208
L210:
	;
	v4834 = int32(0)
	goto L39
L211:
	;
	if int32(1) < v3243 {
		v695 = v2232
		v740 = v3136
		v741 = v3135
		goto L86
	} else {
		goto L219
	}
L212:
	;
	v3243 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+288))
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+288)) = v3243 + int32(-1)
	goto L211
L213:
	;
	v3220 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+36))
	v3221 = int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+36)) = v3220 + v3221
	v3224 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+28))
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+28)) = v3224 + v3221
	v3228 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+40))
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+40)) = v3228 + v3221
	v3232 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+320))
	v3233 = int32(16)
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+320)) = v3232 + v3233
	v3236 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+324))
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+324)) = v3236 + v3233
	goto L212
L214:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3138))) = int32(0)
	v3152 = *(*int32)(unsafe.Add(mBase, uint32(v3147)+uint32(_c_F_VP8EncTokenLoop[5])))
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+40)) = v3152
	v3154 = *(*int64)(unsafe.Add(mBase, uint32(v3147)+uint32(_c_F_VP8EncTokenLoop[7])))
	*(*int64)(unsafe.Add(mBase, uint32(v3138)+320)) = v3154
	v3156 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+4))
	v3158 = v3156 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+4)) = v3158
	v3160 = *(*int32)(unsafe.Add(mBase, uint32(v3147)+uint32(_c_F_VP8EncTokenLoop[4])))
	v3161 = *(*int32)(unsafe.Add(mBase, uint32(v3147)+48))
	v3163 = int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+36)) = v3160 + v3158*v3161<<(uint(v3163)%32)
	v3167 = *(*int32)(unsafe.Add(mBase, uint32(v3147)+uint32(_c_F_VP8EncTokenLoop[6])))
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+28)) = v3167 + v3158*v3145<<(uint(v3163)%32)
	v3173 = *(*int32)(unsafe.Add(mBase, uint32(v3147)+52))
	v3174 = int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+32)) = v3147 + (v3173+v3174)&v3158<<(uint(int32(5))%32) + int32(88)
	v3183 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+316))
	if v3174 < v3156 {
		goto L215
	} else {
		goto L216
	}
L215:
	;
	v3190 = int32(-127)
	goto L217
L216:
	;
	v3190 = int32(127)
	goto L217
L217:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3183+v3174))) = uint8(v3190)
	v3192 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+312))
	v3193 = int32(-1)
	*(*uint8)(unsafe.Add(mBase, uint32(v3192+v3193))) = uint8(v3190)
	v3196 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+308))
	*(*uint8)(unsafe.Add(mBase, uint32(v3196+v3193))) = uint8(v3190)
	v3200 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+308))
	v3201 = int64(-9114861777597660799)
	*(*int64)(unsafe.Add(mBase, uint32(v3200))) = v3201
	*(*int64)(unsafe.Add(mBase, uint32(v3200+int32(8)))) = v3201
	v3207 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+312))
	*(*int64)(unsafe.Add(mBase, uint32(v3207))) = v3201
	v3210 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+316))
	*(*int64)(unsafe.Add(mBase, uint32(v3210))) = v3201
	v3213 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+160)) = v3213
	v3215 = *(*int32)(unsafe.Add(mBase, uint32(v3138)+304))
	if v3215 == v3213 {
		goto L212
	} else {
		goto L218
	}
L218:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3138)+300)) = int32(0)
	goto L212
L219:
	;
	goto L87
L220:
	;
	v4017 = *(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[8])))
	if v4017 < int32(1) {
		goto L299
	} else {
		goto L300
	}
L221:
	;
	v3864 = float64(99)
	if v108 == int64(0) {
		v4016 = v3864
		goto L220
	} else {
		goto L283
	}
L222:
	;
	v3253 = int32(0)
	v3290 = l0 + int32(_a_F_VP8EncTokenLoop_2)
	v3295 = l0 + int32(_a_F_VP8EncTokenLoop_5)
	v3296 = l0 + int32(3442)
	v3297 = l0 + int32(_a_F_VP8EncTokenLoop_6)
	v3298 = l0 + int32(3431)
	v3299 = v3253
	v3301 = v3290
	v3302 = v3253
	v3303 = v3253
	v3304 = v3253
	goto L224
L223:
	;
	v3767 = *(*int32)(unsafe.Add(mBase, uint32(v266)))
	if v3767 != 0 {
		goto L268
	} else {
		goto L269
	}
L224:
	;
	v3331 = v3302
	v3332 = v3303
	v3334 = v3295
	v3335 = v3296
	v3336 = v3297
	v3337 = v3298
	v3338 = v3299
	v3339 = v3301
	v3340 = int32(0)
	goto L226
L225:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = v3716
	goto L223
L226:
	;
	v3360 = v3331
	v3361 = v3332
	v3370 = int32(0)
	v3371 = v3338
	goto L228
L227:
	;
	v3740 = int32(1056)
	v3742 = int32(264)
	v3753 = v3304 + int32(1)
	if v3753 != int32(4) {
		v3295 = v3295 + v3740
		v3296 = v3296 + v3742
		v3297 = v3297 + v3740
		v3298 = v3298 + v3742
		v3299 = v3299 + v3742
		v3301 = v3301 + v3740
		v3302 = v3716
		v3303 = v3717
		v3304 = v3753
		goto L224
	} else {
		goto L265
	}
L228:
	;
	v3381 = *(*int32)(unsafe.Add(mBase, uint32(v3339+v3370)))
	v3383 = int32(base.Ui32(v3381) >> (uint(int32(16)) % 32))
	v3384 = m.G80
	v3386 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3384+v3371))))
	v3387 = m.G81
	v3389 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3387+v3371))))
	v3391 = v3381 & int32(_a_F_VP8EncTokenLoop_7)
	if v3391 != 0 {
		goto L231
	} else {
		goto L232
	}
L229:
	;
	v3485 = v3468
	v3486 = v3469
	v3495 = v3336
	v3501 = int32(0)
	goto L240
L230:
	;
	v3399 = m.G24
	v3402 = v3383 - v3391
	v3403 = int32(1)
	v3406 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+v3386<<(uint(v3403)%32)))))
	v3408 = int32(255)
	v3413 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+(v3386^v3408)<<(uint(v3403)%32)))))
	v3419 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+v3389<<(uint(v3403)%32)))))
	v3420 = v3402*v3406 + v3391*v3413 + v3419
	v3422 = v3398 & v3408
	v3428 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+(v3422^v3408)<<(uint(v3403)%32)))))
	v3433 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+v3422<<(uint(v3403)%32)))))
	v3441 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+(v3389^v3408)<<(uint(v3403)%32)))))
	v3444 = v3391*v3428 + v3402*v3433 + v3441 + int32(2048)
	if v3444 < v3420 {
		goto L233
	} else {
		goto L234
	}
L231:
	;
	v3393 = int32(255)
	v3396 = base.I32_div_u_s(v3391*v3393, v3383)
	v3398 = v3393 - v3396
	goto L230
L232:
	;
	v3398 = int32(255)
	goto L230
L233:
	;
	v3446 = int32(-1)
	goto L235
L234:
	;
	v3446 = int32(0)
	goto L235
L235:
	;
	v3453 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3399+(v3389^v3446)&int32(255)<<(uint(int32(1))%32)))))
	v3454 = v3361 + v3453
	if v3420 <= v3444 {
		goto L237
	} else {
		goto L238
	}
L236:
	;
	v3473 = v3370 + int32(4)
	if v3473 != int32(44) {
		v3360 = v3468
		v3361 = v3469
		v3370 = v3473
		v3371 = v3371 + int32(1)
		goto L228
	} else {
		goto L239
	}
L237:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3290+v3371+int32(-1056)))) = uint8(v3386)
	v3468 = v3360
	v3469 = v3454
	goto L236
L238:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3290+v3371+int32(-1056)))) = uint8(v3398)
	v3468 = v3360 | base.B2i32(v3398 != v3386)
	v3469 = v3454 + int32(2048)
	goto L236
L239:
	;
	goto L229
L240:
	;
	v3505 = *(*int32)(unsafe.Add(mBase, uint32(v3495)))
	v3507 = int32(base.Ui32(v3505) >> (uint(int32(16)) % 32))
	v3508 = m.G80
	v3509 = v3338 + v3501
	v3511 = int32(11)
	v3513 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3508+v3509+v3511))))
	v3514 = m.G81
	v3518 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3514+v3509+v3511))))
	v3520 = v3505 & int32(_a_F_VP8EncTokenLoop_7)
	if v3520 != 0 {
		goto L243
	} else {
		goto L244
	}
L241:
	;
	v3609 = v3592
	v3610 = v3593
	v3619 = v3334
	v3625 = int32(0)
	goto L252
L242:
	;
	v3528 = m.G24
	v3531 = v3507 - v3520
	v3532 = int32(1)
	v3535 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+v3513<<(uint(v3532)%32)))))
	v3537 = int32(255)
	v3542 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+(v3513^v3537)<<(uint(v3532)%32)))))
	v3548 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+v3518<<(uint(v3532)%32)))))
	v3551 = v3527 & v3537
	v3557 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+(v3551^v3537)<<(uint(v3532)%32)))))
	v3562 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+v3551<<(uint(v3532)%32)))))
	v3570 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+(v3518^v3537)<<(uint(v3532)%32)))))
	v3574 = base.B2i32(v3520*v3557+v3531*v3562+v3570+int32(2048) < v3531*v3535+v3520*v3542+v3548)
	if v3520*v3557+v3531*v3562+v3570+int32(2048) < v3531*v3535+v3520*v3542+v3548 {
		goto L245
	} else {
		goto L246
	}
L243:
	;
	v3522 = int32(255)
	v3525 = base.I32_div_u_s(v3520*v3522, v3507)
	v3527 = v3522 - v3525
	goto L242
L244:
	;
	v3527 = int32(255)
	goto L242
L245:
	;
	v3575 = int32(-1)
	goto L247
L246:
	;
	v3575 = int32(0)
	goto L247
L247:
	;
	v3582 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3528+(v3518^v3575)&int32(255)<<(uint(int32(1))%32)))))
	v3583 = v3486 + v3582
	if v3520*v3557+v3531*v3562+v3570+int32(2048) < v3531*v3535+v3520*v3542+v3548 {
		goto L249
	} else {
		goto L250
	}
L248:
	;
	v3597 = v3501 + int32(1)
	if v3597 != int32(11) {
		v3485 = v3592
		v3486 = v3593
		v3495 = v3495 + int32(4)
		v3501 = v3597
		goto L240
	} else {
		goto L251
	}
L249:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3337+v3501))) = uint8(v3527)
	v3592 = v3485 | base.B2i32(v3527 != v3513)
	v3593 = v3583 + int32(2048)
	goto L248
L250:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3337+v3501))) = uint8(v3513)
	v3592 = v3485
	v3593 = v3583
	goto L248
L251:
	;
	goto L241
L252:
	;
	v3629 = *(*int32)(unsafe.Add(mBase, uint32(v3619)))
	v3631 = int32(base.Ui32(v3629) >> (uint(int32(16)) % 32))
	v3632 = m.G80
	v3633 = v3338 + v3625
	v3635 = int32(22)
	v3637 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3632+v3633+v3635))))
	v3638 = m.G81
	v3642 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3638+v3633+v3635))))
	v3644 = v3629 & int32(_a_F_VP8EncTokenLoop_7)
	if v3644 != 0 {
		goto L255
	} else {
		goto L256
	}
L253:
	;
	v3724 = int32(132)
	v3726 = int32(33)
	v3737 = v3340 + int32(1)
	if v3737 != int32(8) {
		v3331 = v3716
		v3332 = v3717
		v3334 = v3334 + v3724
		v3335 = v3335 + v3726
		v3336 = v3336 + v3724
		v3337 = v3337 + v3726
		v3338 = v3338 + v3726
		v3339 = v3339 + v3724
		v3340 = v3737
		goto L226
	} else {
		goto L264
	}
L254:
	;
	v3652 = m.G24
	v3655 = v3631 - v3644
	v3656 = int32(1)
	v3659 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+v3637<<(uint(v3656)%32)))))
	v3661 = int32(255)
	v3666 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+(v3637^v3661)<<(uint(v3656)%32)))))
	v3672 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+v3642<<(uint(v3656)%32)))))
	v3675 = v3651 & v3661
	v3681 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+(v3675^v3661)<<(uint(v3656)%32)))))
	v3686 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+v3675<<(uint(v3656)%32)))))
	v3694 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+(v3642^v3661)<<(uint(v3656)%32)))))
	v3698 = base.B2i32(v3644*v3681+v3655*v3686+v3694+int32(2048) < v3655*v3659+v3644*v3666+v3672)
	if v3644*v3681+v3655*v3686+v3694+int32(2048) < v3655*v3659+v3644*v3666+v3672 {
		goto L257
	} else {
		goto L258
	}
L255:
	;
	v3646 = int32(255)
	v3649 = base.I32_div_u_s(v3644*v3646, v3631)
	v3651 = v3646 - v3649
	goto L254
L256:
	;
	v3651 = int32(255)
	goto L254
L257:
	;
	v3699 = int32(-1)
	goto L259
L258:
	;
	v3699 = int32(0)
	goto L259
L259:
	;
	v3706 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3652+(v3642^v3699)&int32(255)<<(uint(int32(1))%32)))))
	v3707 = v3610 + v3706
	if v3644*v3681+v3655*v3686+v3694+int32(2048) < v3655*v3659+v3644*v3666+v3672 {
		goto L261
	} else {
		goto L262
	}
L260:
	;
	v3721 = v3625 + int32(1)
	if v3721 != int32(11) {
		v3609 = v3716
		v3610 = v3717
		v3619 = v3619 + int32(4)
		v3625 = v3721
		goto L252
	} else {
		goto L263
	}
L261:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3335+v3625))) = uint8(v3651)
	v3716 = v3609 | base.B2i32(v3651 != v3637)
	v3717 = v3707 + int32(2048)
	goto L260
L262:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3335+v3625))) = uint8(v3637)
	v3716 = v3609
	v3717 = v3707
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
	v4016 = base.F64_convert_i64_u(int64(base.Ui64(v3250+base.I64_extend_i32_s(v3717)+base.I64_extend_i32_u(v3848)+int64(1024))>>(uint(int64(11))%64)) + int64(30))
	goto L220
L267:
	;
	goto L266
L268:
	;
	v3769 = *(*int32)(unsafe.Add(mBase, uint32(v266)+16))
	v3777 = v3767
	v3778 = int32(0)
	goto L270
L269:
	;
	v3848 = int32(0)
	goto L267
L270:
	;
	v3786 = *(*int32)(unsafe.Add(mBase, uint32(v3777)))
	if v3786 != 0 {
		v3788 = int32(0)
		goto L272
	} else {
		goto L273
	}
L271:
	;
	v3848 = v3838
	goto L267
L272:
	;
	if v3769 <= v3788 {
		v3838 = v3778
		goto L274
	} else {
		goto L275
	}
L273:
	;
	v3787 = *(*int32)(unsafe.Add(mBase, uint32(v266)+12))
	v3788 = v3787
	goto L272
L274:
	;
	if v3786 != 0 {
		v3777 = v3786
		v3778 = v3838
		goto L270
	} else {
		goto L282
	}
L275:
	;
	v3794 = v3778
	v3798 = v3777 + (v3769<<(uint(int32(1))%32) + int32(2))
	v3799 = v3769
	goto L276
L276:
	;
	v3801 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3798))))
	if v3801&int32(_a_F_VP8EncTokenLoop_19) == int32(0) {
		goto L279
	} else {
		goto L280
	}
L277:
	;
	v3838 = v3831
	goto L274
L278:
	;
	v3826 = m.G24
	v3830 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3826+v3823<<(uint(int32(1))%32)))))
	v3831 = v3794 + v3830
	v3833 = v3799 + int32(-1)
	if v3788 < v3833 {
		v3794 = v3831
		v3798 = v3798 + int32(-2)
		v3799 = v3833
		goto L276
	} else {
		goto L281
	}
L279:
	;
	v3816 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3420)+v3801&int32(_a_F_VP8EncTokenLoop_20)))))
	v3823 = v3816 ^ int32(base.Ui32(base.I32_extend16_s(v3801))>>(uint(int32(15))%32))&int32(255)
	goto L278
L280:
	;
	v3823 = (v3801<<(uint(int32(16))%32)>>(uint(int32(31))%32) ^ v3801) & int32(255)
	goto L278
L281:
	;
	goto L277
L282:
	;
	goto L271
L283:
	;
	if v3135 == int64(0) {
		v4016 = v3864
		goto L220
	} else {
		goto L284
	}
L284:
	;
	v3870 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v108), float64(65025)), base.F64_convert_i64_u(v3135))
	v3881 = base.I64_reinterpret_f64(v3870)
	if int64(4503599627370495) < v3881 {
		goto L289
	} else {
		goto L290
	}
L285:
	;
	v4016 = base.F64_mul(v4013, float64(10))
	goto L220
L286:
	;
	v4013 = v3991
	goto L285
L287:
	;
	v3917 = v3915 + int32(614242)
	v3921 = base.F64_convert_i32_s(v3913 + int32(base.Ui32(v3917)>>(uint(int32(20))%32)))
	v3923 = base.F64_mul(v3921, float64(0.30102999566361177))
	v3936 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v3917&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v3912&int64(4294967295)), float64(-1))
	v3939 = base.F64_mul(v3936, base.F64_mul(v3936, float64(0.5)))
	v3944 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v3936, v3939)) & int64(-4294967296))
	v3945 = float64(0.4342944818781689)
	v3946 = base.F64_mul(v3944, v3945)
	v3947 = base.F64_add(v3923, v3946)
	v3952 = base.F64_div(v3936, base.F64_add(v3936, float64(2)))
	v3953 = base.F64_mul(v3952, v3952)
	v3954 = base.F64_mul(v3953, v3953)
	v3979 = base.F64_add(base.F64_mul(v3952, base.F64_add(v3939, base.F64_add(base.F64_mul(v3954, base.F64_add(base.F64_mul(v3954, base.F64_add(base.F64_mul(v3954, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v3953, base.F64_add(base.F64_mul(v3954, base.F64_add(base.F64_mul(v3954, base.F64_add(base.F64_mul(v3954, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v3936, v3944), v3939))
	v3991 = base.F64_add(v3947, base.F64_add(base.F64_add(v3946, base.F64_sub(v3923, v3947)), base.F64_add(base.F64_mul(v3979, v3945), base.F64_add(base.F64_mul(v3921, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v3979, v3944), float64(2.5082946711645275e-11))))))
	goto L286
L288:
	;
	v3907 = base.I64_reinterpret_f64(base.F64_mul(v3870, float64(1.8014398509481984e+16)))
	v3912 = v3907
	v3913 = int32(-1077)
	v3915 = base.I32_wrap_i64(int64(base.Ui64(v3907) >> (uint(int64(32)) % 64)))
	goto L287
L289:
	;
	if base.Ui64(int64(9218868437227405311)) < base.Ui64(v3881) {
		v3991 = v3870
		goto L286
	} else {
		goto L294
	}
L290:
	;
	if base.F64_ne(v3870, float64(0)) != 0 {
		goto L291
	} else {
		goto L292
	}
L291:
	;
	if int64(-1) < v3881 {
		goto L288
	} else {
		goto L293
	}
L292:
	;
	v4013 = base.F64_div(float64(-1), base.F64_mul(v3870, v3870))
	goto L285
L293:
	;
	v4013 = base.F64_div(base.F64_sub(v3870, v3870), float64(0))
	goto L285
L294:
	;
	v3896 = int32(-1023)
	v3898 = int64(base.Ui64(v3881) >> (uint(int64(32)) % 64))
	if v3898 == int64(1072693248) {
		goto L295
	} else {
		goto L296
	}
L295:
	;
	if base.I32_wrap_i64(v3881) != 0 {
		v3912 = v3881
		v3913 = v3896
		v3915 = int32(1072693248)
		goto L287
	} else {
		goto L297
	}
L296:
	;
	v3912 = v3881
	v3913 = v3896
	v3915 = base.I32_wrap_i64(v3898)
	goto L287
L297:
	;
	v4013 = float64(0)
	goto L285
L298:
	;
	if int32(0) < v4087 {
		v330 = v4087
		v333 = v4088
		v339 = v4089
		v345 = v651
		v375 = v4090
		v376 = v4091
		v377 = v4092
		goto L46
	} else {
		goto L328
	}
L299:
	;
	if v399 != 0 {
		v4114 = v651
		goto L40
	} else {
		goto L305
	}
L300:
	;
	if base.Ui64(v3250) < base.Ui64(int64(1069547521)) {
		goto L299
	} else {
		goto L301
	}
L301:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[8]))) = int32(base.Ui32(v4017) >> (uint(int32(1)) % 32))
	if v399 == int32(0) {
		v4087 = v330
		v4088 = v333
		v4089 = v339
		v4090 = v375
		v4091 = v376
		v4092 = v377
		goto L298
	} else {
		goto L302
	}
L302:
	;
	v4027 = *(*int32)(unsafe.Add(mBase, uint32(v64)+904))
	v4028 = *(*int32)(unsafe.Add(mBase, uint32(v4027)+4))
	v4029 = *(*int32)(unsafe.Add(mBase, uint32(v4028)+88))
	if v4029 == int32(0) {
		goto L303
	} else {
		goto L304
	}
L303:
	;
	v4038 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v4027)+uint32(_c_F_VP8EncTokenLoop[15]))) = v4038
	*(*int64)(unsafe.Add(mBase, uint32(v4027)+uint32(_c_F_VP8EncTokenLoop[16]))) = v4038
	*(*int64)(unsafe.Add(mBase, uint32(v4027)+uint32(_c_F_VP8EncTokenLoop[17]))) = v4038
	*(*int64)(unsafe.Add(mBase, uint32(v4027)+uint32(_c_F_VP8EncTokenLoop[18]))) = v4038
	v4087 = v330
	v4088 = v333
	v4089 = v339
	v4090 = v375
	v4091 = v376
	v4092 = v377
	goto L298
L304:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v4027)+uint32(_c_F_VP8EncTokenLoop[19]))) = int64(0)
	*(*int32)(unsafe.Add(mBase, uint32(v4027)+uint32(_c_F_VP8EncTokenLoop[20]))) = int32(0)
	goto L303
L305:
	;
	if v66 == int32(0) {
		v4081 = v333
		v4082 = v339
		v4083 = v375
		v4084 = v376
		v4085 = v377
		goto L306
	} else {
		goto L307
	}
L306:
	;
	v4087 = v388
	v4088 = v4081
	v4089 = v4082
	v4090 = v4083
	v4091 = v4084
	v4092 = v4085
	goto L298
L307:
	;
	if v376 == int32(0) {
		goto L309
	} else {
		goto L310
	}
L308:
	;
	v4069 = float32(30)
	if base.F32_gt(v4066, v4069) != 0 {
		goto L316
	} else {
		goto L317
	}
L309:
	;
	if base.F64_ne(v4016, v375) != 0 {
		goto L314
	} else {
		goto L315
	}
L310:
	;
	if base.F64_gt(v4016, v101) != 0 {
		goto L311
	} else {
		goto L312
	}
L311:
	;
	v4056 = base.F32_neg(v339)
	goto L313
L312:
	;
	v4056 = v339
	goto L313
L313:
	;
	v4066 = v4056
	goto L308
L314:
	;
	v4066 = base.F32_demote_f64(base.F64_mul(base.F64_div(base.F64_sub(v101, v4016), base.F64_sub(v375, v4016)), base.F64_promote_f32(base.F32_sub(v377, v333))))
	goto L308
L315:
	;
	v4066 = float32(0)
	goto L308
L316:
	;
	v4072 = v4069
	goto L318
L317:
	;
	v4072 = v4066
	goto L318
L318:
	;
	if base.F32_lt(v4066, float32(-30)) != 0 {
		goto L319
	} else {
		goto L320
	}
L319:
	;
	v4075 = float32(-30)
	goto L321
L320:
	;
	v4075 = v4072
	goto L321
L321:
	;
	v4076 = base.F32_add(v333, v4075)
	if base.F32_gt(v4076, v72) != 0 {
		goto L322
	} else {
		goto L323
	}
L322:
	;
	v4078 = v72
	goto L324
L323:
	;
	v4078 = v4076
	goto L324
L324:
	;
	if base.F32_lt(v4076, v70) != 0 {
		goto L325
	} else {
		goto L326
	}
L325:
	;
	v4080 = v70
	goto L327
L326:
	;
	v4080 = v4078
	goto L327
L327:
	;
	v4081 = v4080
	v4082 = v4075
	v4083 = v4016
	v4084 = int32(0)
	v4085 = v333
	goto L306
L328:
	;
	goto L47
L329:
	;
	v4661 = l0 + int32(344)
	v4672 = *(*int32)(unsafe.Add(mBase, uint32(v4661)))
	if v4672 == int32(0) {
		goto L376
	} else {
		goto L377
	}
L330:
	;
	v4156 = int32(0)
	v4193 = l0 + int32(_a_F_VP8EncTokenLoop_2)
	v4198 = l0 + int32(_a_F_VP8EncTokenLoop_5)
	v4199 = l0 + int32(3442)
	v4200 = l0 + int32(_a_F_VP8EncTokenLoop_6)
	v4201 = l0 + int32(3431)
	v4202 = v4156
	v4204 = v4193
	v4205 = v4156
	v4207 = v4156
	goto L332
L331:
	;
	goto L329
L332:
	;
	v4234 = v4205
	v4237 = v4198
	v4238 = v4199
	v4239 = v4200
	v4240 = v4201
	v4241 = v4202
	v4242 = v4204
	v4243 = int32(0)
	goto L334
L333:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncTokenLoop[9]))) = v4619
	goto L331
L334:
	;
	v4263 = v4234
	v4273 = int32(0)
	v4274 = v4241
	goto L336
L335:
	;
	v4643 = int32(1056)
	v4645 = int32(264)
	v4656 = v4207 + int32(1)
	if v4656 != int32(4) {
		v4198 = v4198 + v4643
		v4199 = v4199 + v4645
		v4200 = v4200 + v4643
		v4201 = v4201 + v4645
		v4202 = v4202 + v4645
		v4204 = v4204 + v4643
		v4205 = v4619
		v4207 = v4656
		goto L332
	} else {
		goto L373
	}
L336:
	;
	v4284 = *(*int32)(unsafe.Add(mBase, uint32(v4242+v4273)))
	v4286 = int32(base.Ui32(v4284) >> (uint(int32(16)) % 32))
	v4287 = m.G80
	v4289 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4287+v4274))))
	v4290 = m.G81
	v4292 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4290+v4274))))
	v4294 = v4284 & int32(_a_F_VP8EncTokenLoop_7)
	if v4294 != 0 {
		goto L339
	} else {
		goto L340
	}
L337:
	;
	v4388 = v4371
	v4398 = v4239
	v4404 = int32(0)
	goto L348
L338:
	;
	v4302 = m.G24
	v4305 = v4286 - v4294
	v4306 = int32(1)
	v4309 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4302+v4289<<(uint(v4306)%32)))))
	v4311 = int32(255)
	v4316 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4302+(v4289^v4311)<<(uint(v4306)%32)))))
	v4322 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4302+v4292<<(uint(v4306)%32)))))
	v4325 = v4301 & v4311
	v4331 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4302+(v4325^v4311)<<(uint(v4306)%32)))))
	v4336 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4302+v4325<<(uint(v4306)%32)))))
	v4344 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4302+(v4292^v4311)<<(uint(v4306)%32)))))
	goto L341
L339:
	;
	v4296 = int32(255)
	v4299 = base.I32_div_u_s(v4294*v4296, v4286)
	v4301 = v4296 - v4299
	goto L338
L340:
	;
	v4301 = int32(255)
	goto L338
L341:
	;
	goto L343
L343:
	;
	if v4305*v4309+v4294*v4316+v4322 <= v4294*v4331+v4305*v4336+v4344+int32(2048) {
		goto L345
	} else {
		goto L346
	}
L344:
	;
	v4376 = v4273 + int32(4)
	if v4376 != int32(44) {
		v4263 = v4371
		v4273 = v4376
		v4274 = v4274 + int32(1)
		goto L336
	} else {
		goto L347
	}
L345:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4193+v4274+int32(-1056)))) = uint8(v4289)
	v4371 = v4263
	goto L344
L346:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4193+v4274+int32(-1056)))) = uint8(v4301)
	v4371 = v4263 | base.B2i32(v4301 != v4289)
	goto L344
L347:
	;
	goto L337
L348:
	;
	v4408 = *(*int32)(unsafe.Add(mBase, uint32(v4398)))
	v4410 = int32(base.Ui32(v4408) >> (uint(int32(16)) % 32))
	v4411 = m.G80
	v4412 = v4241 + v4404
	v4414 = int32(11)
	v4416 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4411+v4412+v4414))))
	v4417 = m.G81
	v4421 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4417+v4412+v4414))))
	v4423 = v4408 & int32(_a_F_VP8EncTokenLoop_7)
	if v4423 != 0 {
		goto L351
	} else {
		goto L352
	}
L349:
	;
	v4512 = v4495
	v4522 = v4237
	v4528 = int32(0)
	goto L360
L350:
	;
	v4431 = m.G24
	v4434 = v4410 - v4423
	v4435 = int32(1)
	v4438 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4431+v4416<<(uint(v4435)%32)))))
	v4440 = int32(255)
	v4445 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4431+(v4416^v4440)<<(uint(v4435)%32)))))
	v4451 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4431+v4421<<(uint(v4435)%32)))))
	v4454 = v4430 & v4440
	v4460 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4431+(v4454^v4440)<<(uint(v4435)%32)))))
	v4465 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4431+v4454<<(uint(v4435)%32)))))
	v4473 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4431+(v4421^v4440)<<(uint(v4435)%32)))))
	goto L353
L351:
	;
	v4425 = int32(255)
	v4428 = base.I32_div_u_s(v4423*v4425, v4410)
	v4430 = v4425 - v4428
	goto L350
L352:
	;
	v4430 = int32(255)
	goto L350
L353:
	;
	goto L355
L355:
	;
	if v4423*v4460+v4434*v4465+v4473+int32(2048) < v4434*v4438+v4423*v4445+v4451 {
		goto L357
	} else {
		goto L358
	}
L356:
	;
	v4500 = v4404 + int32(1)
	if v4500 != int32(11) {
		v4388 = v4495
		v4398 = v4398 + int32(4)
		v4404 = v4500
		goto L348
	} else {
		goto L359
	}
L357:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4240+v4404))) = uint8(v4430)
	v4495 = v4388 | base.B2i32(v4430 != v4416)
	goto L356
L358:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4240+v4404))) = uint8(v4416)
	v4495 = v4388
	goto L356
L359:
	;
	goto L349
L360:
	;
	v4532 = *(*int32)(unsafe.Add(mBase, uint32(v4522)))
	v4534 = int32(base.Ui32(v4532) >> (uint(int32(16)) % 32))
	v4535 = m.G80
	v4536 = v4241 + v4528
	v4538 = int32(22)
	v4540 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4535+v4536+v4538))))
	v4541 = m.G81
	v4545 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4541+v4536+v4538))))
	v4547 = v4532 & int32(_a_F_VP8EncTokenLoop_7)
	if v4547 != 0 {
		goto L363
	} else {
		goto L364
	}
L361:
	;
	v4627 = int32(132)
	v4629 = int32(33)
	v4640 = v4243 + int32(1)
	if v4640 != int32(8) {
		v4234 = v4619
		v4237 = v4237 + v4627
		v4238 = v4238 + v4629
		v4239 = v4239 + v4627
		v4240 = v4240 + v4629
		v4241 = v4241 + v4629
		v4242 = v4242 + v4627
		v4243 = v4640
		goto L334
	} else {
		goto L372
	}
L362:
	;
	v4555 = m.G24
	v4558 = v4534 - v4547
	v4559 = int32(1)
	v4562 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4555+v4540<<(uint(v4559)%32)))))
	v4564 = int32(255)
	v4569 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4555+(v4540^v4564)<<(uint(v4559)%32)))))
	v4575 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4555+v4545<<(uint(v4559)%32)))))
	v4578 = v4554 & v4564
	v4584 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4555+(v4578^v4564)<<(uint(v4559)%32)))))
	v4589 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4555+v4578<<(uint(v4559)%32)))))
	v4597 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4555+(v4545^v4564)<<(uint(v4559)%32)))))
	goto L365
L363:
	;
	v4549 = int32(255)
	v4552 = base.I32_div_u_s(v4547*v4549, v4534)
	v4554 = v4549 - v4552
	goto L362
L364:
	;
	v4554 = int32(255)
	goto L362
L365:
	;
	goto L367
L367:
	;
	if v4547*v4584+v4558*v4589+v4597+int32(2048) < v4558*v4562+v4547*v4569+v4575 {
		goto L369
	} else {
		goto L370
	}
L368:
	;
	v4624 = v4528 + int32(1)
	if v4624 != int32(11) {
		v4512 = v4619
		v4522 = v4522 + int32(4)
		v4528 = v4624
		goto L360
	} else {
		goto L371
	}
L369:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4238+v4528))) = uint8(v4554)
	v4619 = v4512 | base.B2i32(v4554 != v4540)
	goto L368
L370:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4238+v4528))) = uint8(v4540)
	v4619 = v4512
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
	v4823 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v4824 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	v4828 = F_WebPReportProgress(m, v4823, v4824+v4114, l0+int32(368))
	mBase = m.M
	v4834 = base.B2i32(v4828 != int32(0))
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
	v4681 = v4672
	goto L380
L380:
	;
	v4689 = *(*int32)(unsafe.Add(mBase, uint32(v4681)))
	if v4689 != 0 {
		v4691 = int32(0)
		goto L382
	} else {
		goto L383
	}
L382:
	;
	v4692 = *(*int32)(unsafe.Add(mBase, uint32(v4661)+16))
	if v4692 <= v4691 {
		goto L384
	} else {
		goto L385
	}
L383:
	;
	v4690 = *(*int32)(unsafe.Add(mBase, uint32(v4661)+12))
	v4691 = v4690
	goto L382
L384:
	;
	F_WebPSafeFree(m, v4681)
	mBase = m.M
	if v4689 != 0 {
		v4681 = v4689
		goto L380
	} else {
		goto L392
	}
L385:
	;
	v4706 = v4692
	v4707 = v4681 + v4692<<(uint(int32(1))%32) + int32(2)
	goto L386
L386:
	;
	v4710 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4707))))
	if v4710&int32(_a_F_VP8EncTokenLoop_19) == int32(0) {
		goto L389
	} else {
		goto L390
	}
L387:
	;
	goto L384
L388:
	;
	v4724 = F_VP8PutBit(m, v112, int32(base.Ui32(v4710)>>(uint(int32(15))%32)), v4723)
	mBase = m.M
	v4728 = v4706 + int32(-1)
	if v4691 < v4728 {
		v4706 = v4728
		v4707 = v4707 + int32(-2)
		goto L386
	} else {
		goto L391
	}
L389:
	;
	v4722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+int32(3420)+v4710&int32(_a_F_VP8EncTokenLoop_20)))))
	v4723 = v4722
	goto L388
L390:
	;
	v4723 = v4710 & int32(255)
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
	*(*int32)(unsafe.Add(mBase, uint32(v4661))) = int32(0)
	goto L406
L409:
	;
	v5032 = v5028
	goto L13
L410:
	;
	F_VP8EncFreeBitWriters(m, v4898)
	mBase = m.M
	v5019 = *(*int32)(unsafe.Add(mBase, uint32(v4898)+4))
	v5021 = F_WebPEncodingSetError(m, v5019, int32(1))
	mBase = m.M
	v5028 = v5021
	goto L409
L411:
	;
	v4901 = *(*int32)(unsafe.Add(mBase, uint32(v4898)+52))
	if v4901 < int32(1) {
		v4930 = v4834
		goto L412
	} else {
		goto L413
	}
L412:
	;
	v4935 = *(*int32)(unsafe.Add(mBase, uint32(v4898)+4))
	v4936 = *(*int32)(unsafe.Add(mBase, uint32(v4935)+88))
	if v4936 == int32(0) {
		goto L421
	} else {
		goto L422
	}
L413:
	;
	v4912 = v4834 & int32(1)
	v4913 = v4898 + int32(88)
	v4914 = int32(0)
	goto L414
L414:
	;
	v4915 = F_VP8BitWriterFinish(m, v4913)
	mBase = m.M
	v4919 = *(*int32)(unsafe.Add(mBase, uint32(v4913+int32(28))))
	if v4919 != 0 {
		goto L416
	} else {
		goto L417
	}
L415:
	;
	if v4920 == int32(0) {
		goto L410
	} else {
		goto L420
	}
L416:
	;
	v4920 = int32(0)
	goto L418
L417:
	;
	v4920 = v4912
	goto L418
L418:
	;
	v4924 = v4914 + int32(1)
	v4925 = *(*int32)(unsafe.Add(mBase, uint32(v4898)+52))
	if v4924 < v4925 {
		v4912 = v4920
		v4913 = v4913 + int32(32)
		v4914 = v4924
		goto L414
	} else {
		goto L419
	}
L419:
	;
	goto L415
L420:
	;
	v4930 = v4912
	goto L412
L421:
	;
	F_VP8AdjustFilterStrength(m, v4893)
	mBase = m.M
	v5028 = v4930
	goto L409
L422:
	;
	v4939 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+168))
	v4940 = int64(7)
	v4942 = int64(3)
	v4943 = int64(base.Ui64(v4939+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[21]))) = uint32(v4943)
	v4945 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+192))
	v4949 = int64(base.Ui64(v4945+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[22]))) = uint32(v4949)
	v4951 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+216))
	v4955 = int64(base.Ui64(v4951+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[23]))) = uint32(v4955)
	v4957 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+240))
	v4961 = int64(base.Ui64(v4957+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[24]))) = uint32(v4961)
	v4963 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+176))
	v4967 = int64(base.Ui64(v4963+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[25]))) = uint32(v4967)
	v4969 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+200))
	v4973 = int64(base.Ui64(v4969+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[26]))) = uint32(v4973)
	v4975 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+224))
	v4979 = int64(base.Ui64(v4975+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[27]))) = uint32(v4979)
	v4981 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+248))
	v4985 = int64(base.Ui64(v4981+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[28]))) = uint32(v4985)
	v4987 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+184))
	v4991 = int64(base.Ui64(v4987+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[29]))) = uint32(v4991)
	v4993 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+208))
	v4997 = int64(base.Ui64(v4993+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[30]))) = uint32(v4997)
	v4999 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+232))
	v5003 = int64(base.Ui64(v4999+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[31]))) = uint32(v5003)
	v5005 = *(*int64)(unsafe.Add(mBase, uint32(v4893)+256))
	v5009 = int64(base.Ui64(v5005+v4940) >> (uint(v4942) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v4898)+uint32(_c_F_VP8EncTokenLoop[32]))) = uint32(v5009)
	goto L421
}
func F_VP8EncWrite(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v18 int64
	_ = v18
	var v19 int32
	_ = v19
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v25 int32
	_ = v25
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
	var v35 int32
	_ = v35
	var v36 int32
	_ = v36
	var v41 int32
	_ = v41
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v65 int32
	_ = v65
	var v70 int32
	_ = v70
	var v71 int32
	_ = v71
	var v72 int32
	_ = v72
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v86 int32
	_ = v86
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v111 int32
	_ = v111
	var v115 int32
	_ = v115
	var v119 int32
	_ = v119
	var v128 int32
	_ = v128
	var v130 int32
	_ = v130
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v151 int32
	_ = v151
	var v155 int32
	_ = v155
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v168 int32
	_ = v168
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v175 int32
	_ = v175
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v187 int32
	_ = v187
	var v191 int32
	_ = v191
	var v195 int32
	_ = v195
	var v203 int32
	_ = v203
	var v206 int32
	_ = v206
	var v208 int32
	_ = v208
	var v211 int32
	_ = v211
	var v213 int32
	_ = v213
	var v217 int32
	_ = v217
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v225 int32
	_ = v225
	var v229 int32
	_ = v229
	var v233 int32
	_ = v233
	var v242 int32
	_ = v242
	var v247 int32
	_ = v247
	var v249 int32
	_ = v249
	var v252 int32
	_ = v252
	var v257 int32
	_ = v257
	var v258 int32
	_ = v258
	var v261 int32
	_ = v261
	var v265 int32
	_ = v265
	var v269 int32
	_ = v269
	var v280 int32
	_ = v280
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v290 int32
	_ = v290
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v307 int32
	_ = v307
	var v313 int32
	_ = v313
	var v318 int32
	_ = v318
	var v320 int32
	_ = v320
	var v323 int32
	_ = v323
	var v325 int32
	_ = v325
	var v329 int32
	_ = v329
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v337 int32
	_ = v337
	var v341 int32
	_ = v341
	var v345 int32
	_ = v345
	var v355 int32
	_ = v355
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v369 int32
	_ = v369
	var v373 int32
	_ = v373
	var v375 int32
	_ = v375
	var v379 int32
	_ = v379
	var v383 int32
	_ = v383
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v391 int32
	_ = v391
	var v395 int32
	_ = v395
	var v411 int32
	_ = v411
	var v415 int32
	_ = v415
	var v417 int32
	_ = v417
	var v421 int32
	_ = v421
	var v423 int32
	_ = v423
	var v427 int32
	_ = v427
	var v431 int32
	_ = v431
	var v432 int32
	_ = v432
	var v435 int32
	_ = v435
	var v439 int32
	_ = v439
	var v443 int32
	_ = v443
	var v449 int32
	_ = v449
	var v459 int32
	_ = v459
	var v464 int32
	_ = v464
	var v466 int32
	_ = v466
	var v469 int32
	_ = v469
	var v471 int32
	_ = v471
	var v475 int32
	_ = v475
	var v479 int32
	_ = v479
	var v480 int32
	_ = v480
	var v483 int32
	_ = v483
	var v487 int32
	_ = v487
	var v491 int32
	_ = v491
	var v501 int32
	_ = v501
	var v509 int32
	_ = v509
	var v513 int32
	_ = v513
	var v515 int32
	_ = v515
	var v519 int32
	_ = v519
	var v521 int32
	_ = v521
	var v525 int32
	_ = v525
	var v529 int32
	_ = v529
	var v530 int32
	_ = v530
	var v533 int32
	_ = v533
	var v537 int32
	_ = v537
	var v541 int32
	_ = v541
	var v557 int32
	_ = v557
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v567 int32
	_ = v567
	var v569 int32
	_ = v569
	var v573 int32
	_ = v573
	var v577 int32
	_ = v577
	var v578 int32
	_ = v578
	var v581 int32
	_ = v581
	var v585 int32
	_ = v585
	var v589 int32
	_ = v589
	var v595 int32
	_ = v595
	var v605 int32
	_ = v605
	var v610 int32
	_ = v610
	var v612 int32
	_ = v612
	var v615 int32
	_ = v615
	var v617 int32
	_ = v617
	var v621 int32
	_ = v621
	var v625 int32
	_ = v625
	var v626 int32
	_ = v626
	var v629 int32
	_ = v629
	var v633 int32
	_ = v633
	var v637 int32
	_ = v637
	var v647 int32
	_ = v647
	var v655 int32
	_ = v655
	var v659 int32
	_ = v659
	var v661 int32
	_ = v661
	var v665 int32
	_ = v665
	var v667 int32
	_ = v667
	var v671 int32
	_ = v671
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v679 int32
	_ = v679
	var v683 int32
	_ = v683
	var v687 int32
	_ = v687
	var v703 int32
	_ = v703
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v713 int32
	_ = v713
	var v715 int32
	_ = v715
	var v719 int32
	_ = v719
	var v723 int32
	_ = v723
	var v724 int32
	_ = v724
	var v727 int32
	_ = v727
	var v731 int32
	_ = v731
	var v735 int32
	_ = v735
	var v741 int32
	_ = v741
	var v751 int32
	_ = v751
	var v756 int32
	_ = v756
	var v758 int32
	_ = v758
	var v761 int32
	_ = v761
	var v763 int32
	_ = v763
	var v767 int32
	_ = v767
	var v771 int32
	_ = v771
	var v772 int32
	_ = v772
	var v775 int32
	_ = v775
	var v779 int32
	_ = v779
	var v783 int32
	_ = v783
	var v793 int32
	_ = v793
	var v801 int32
	_ = v801
	var v805 int32
	_ = v805
	var v807 int32
	_ = v807
	var v811 int32
	_ = v811
	var v813 int32
	_ = v813
	var v817 int32
	_ = v817
	var v821 int32
	_ = v821
	var v822 int32
	_ = v822
	var v825 int32
	_ = v825
	var v829 int32
	_ = v829
	var v833 int32
	_ = v833
	var v849 int32
	_ = v849
	var v853 int32
	_ = v853
	var v855 int32
	_ = v855
	var v859 int32
	_ = v859
	var v861 int32
	_ = v861
	var v865 int32
	_ = v865
	var v869 int32
	_ = v869
	var v870 int32
	_ = v870
	var v873 int32
	_ = v873
	var v877 int32
	_ = v877
	var v881 int32
	_ = v881
	var v887 int32
	_ = v887
	var v897 int32
	_ = v897
	var v902 int32
	_ = v902
	var v904 int32
	_ = v904
	var v907 int32
	_ = v907
	var v909 int32
	_ = v909
	var v913 int32
	_ = v913
	var v917 int32
	_ = v917
	var v918 int32
	_ = v918
	var v921 int32
	_ = v921
	var v925 int32
	_ = v925
	var v929 int32
	_ = v929
	var v939 int32
	_ = v939
	var v947 int32
	_ = v947
	var v951 int32
	_ = v951
	var v953 int32
	_ = v953
	var v957 int32
	_ = v957
	var v959 int32
	_ = v959
	var v963 int32
	_ = v963
	var v967 int32
	_ = v967
	var v968 int32
	_ = v968
	var v971 int32
	_ = v971
	var v975 int32
	_ = v975
	var v979 int32
	_ = v979
	var v995 int32
	_ = v995
	var v999 int32
	_ = v999
	var v1001 int32
	_ = v1001
	var v1005 int32
	_ = v1005
	var v1007 int32
	_ = v1007
	var v1011 int32
	_ = v1011
	var v1015 int32
	_ = v1015
	var v1016 int32
	_ = v1016
	var v1019 int32
	_ = v1019
	var v1023 int32
	_ = v1023
	var v1027 int32
	_ = v1027
	var v1033 int32
	_ = v1033
	var v1043 int32
	_ = v1043
	var v1048 int32
	_ = v1048
	var v1050 int32
	_ = v1050
	var v1053 int32
	_ = v1053
	var v1055 int32
	_ = v1055
	var v1059 int32
	_ = v1059
	var v1063 int32
	_ = v1063
	var v1064 int32
	_ = v1064
	var v1067 int32
	_ = v1067
	var v1071 int32
	_ = v1071
	var v1075 int32
	_ = v1075
	var v1085 int32
	_ = v1085
	var v1093 int32
	_ = v1093
	var v1097 int32
	_ = v1097
	var v1099 int32
	_ = v1099
	var v1103 int32
	_ = v1103
	var v1105 int32
	_ = v1105
	var v1109 int32
	_ = v1109
	var v1113 int32
	_ = v1113
	var v1114 int32
	_ = v1114
	var v1117 int32
	_ = v1117
	var v1121 int32
	_ = v1121
	var v1125 int32
	_ = v1125
	var v1141 int32
	_ = v1141
	var v1145 int32
	_ = v1145
	var v1147 int32
	_ = v1147
	var v1151 int32
	_ = v1151
	var v1153 int32
	_ = v1153
	var v1157 int32
	_ = v1157
	var v1161 int32
	_ = v1161
	var v1162 int32
	_ = v1162
	var v1165 int32
	_ = v1165
	var v1169 int32
	_ = v1169
	var v1173 int32
	_ = v1173
	var v1179 int32
	_ = v1179
	var v1189 int32
	_ = v1189
	var v1194 int32
	_ = v1194
	var v1196 int32
	_ = v1196
	var v1199 int32
	_ = v1199
	var v1201 int32
	_ = v1201
	var v1205 int32
	_ = v1205
	var v1209 int32
	_ = v1209
	var v1210 int32
	_ = v1210
	var v1213 int32
	_ = v1213
	var v1217 int32
	_ = v1217
	var v1221 int32
	_ = v1221
	var v1231 int32
	_ = v1231
	var v1239 int32
	_ = v1239
	var v1243 int32
	_ = v1243
	var v1245 int32
	_ = v1245
	var v1249 int32
	_ = v1249
	var v1251 int32
	_ = v1251
	var v1255 int32
	_ = v1255
	var v1259 int32
	_ = v1259
	var v1260 int32
	_ = v1260
	var v1263 int32
	_ = v1263
	var v1267 int32
	_ = v1267
	var v1271 int32
	_ = v1271
	var v1287 int32
	_ = v1287
	var v1291 int32
	_ = v1291
	var v1293 int32
	_ = v1293
	var v1297 int32
	_ = v1297
	var v1299 int32
	_ = v1299
	var v1303 int32
	_ = v1303
	var v1307 int32
	_ = v1307
	var v1308 int32
	_ = v1308
	var v1311 int32
	_ = v1311
	var v1315 int32
	_ = v1315
	var v1319 int32
	_ = v1319
	var v1325 int32
	_ = v1325
	var v1335 int32
	_ = v1335
	var v1340 int32
	_ = v1340
	var v1342 int32
	_ = v1342
	var v1345 int32
	_ = v1345
	var v1347 int32
	_ = v1347
	var v1351 int32
	_ = v1351
	var v1355 int32
	_ = v1355
	var v1356 int32
	_ = v1356
	var v1359 int32
	_ = v1359
	var v1363 int32
	_ = v1363
	var v1367 int32
	_ = v1367
	var v1377 int32
	_ = v1377
	var v1385 int32
	_ = v1385
	var v1389 int32
	_ = v1389
	var v1391 int32
	_ = v1391
	var v1395 int32
	_ = v1395
	var v1397 int32
	_ = v1397
	var v1401 int32
	_ = v1401
	var v1405 int32
	_ = v1405
	var v1406 int32
	_ = v1406
	var v1409 int32
	_ = v1409
	var v1413 int32
	_ = v1413
	var v1417 int32
	_ = v1417
	var v1433 int32
	_ = v1433
	var v1437 int32
	_ = v1437
	var v1439 int32
	_ = v1439
	var v1443 int32
	_ = v1443
	var v1445 int32
	_ = v1445
	var v1449 int32
	_ = v1449
	var v1453 int32
	_ = v1453
	var v1454 int32
	_ = v1454
	var v1457 int32
	_ = v1457
	var v1461 int32
	_ = v1461
	var v1465 int32
	_ = v1465
	var v1471 int32
	_ = v1471
	var v1481 int32
	_ = v1481
	var v1484 int32
	_ = v1484
	var v1486 int32
	_ = v1486
	var v1489 int32
	_ = v1489
	var v1491 int32
	_ = v1491
	var v1494 int32
	_ = v1494
	var v1496 int32
	_ = v1496
	var v1500 int32
	_ = v1500
	var v1504 int32
	_ = v1504
	var v1505 int32
	_ = v1505
	var v1508 int32
	_ = v1508
	var v1512 int32
	_ = v1512
	var v1516 int32
	_ = v1516
	var v1524 int32
	_ = v1524
	var v1536 int32
	_ = v1536
	var v1539 int32
	_ = v1539
	var v1541 int32
	_ = v1541
	var v1545 int32
	_ = v1545
	var v1547 int32
	_ = v1547
	var v1551 int32
	_ = v1551
	var v1555 int32
	_ = v1555
	var v1556 int32
	_ = v1556
	var v1559 int32
	_ = v1559
	var v1563 int32
	_ = v1563
	var v1567 int32
	_ = v1567
	var v1573 int32
	_ = v1573
	var v1577 int32
	_ = v1577
	var v1579 int32
	_ = v1579
	var v1582 int32
	_ = v1582
	var v1584 int32
	_ = v1584
	var v1587 int32
	_ = v1587
	var v1589 int32
	_ = v1589
	var v1593 int32
	_ = v1593
	var v1597 int32
	_ = v1597
	var v1598 int32
	_ = v1598
	var v1601 int32
	_ = v1601
	var v1605 int32
	_ = v1605
	var v1609 int32
	_ = v1609
	var v1617 int32
	_ = v1617
	var v1629 int32
	_ = v1629
	var v1632 int32
	_ = v1632
	var v1634 int32
	_ = v1634
	var v1638 int32
	_ = v1638
	var v1640 int32
	_ = v1640
	var v1644 int32
	_ = v1644
	var v1648 int32
	_ = v1648
	var v1649 int32
	_ = v1649
	var v1652 int32
	_ = v1652
	var v1656 int32
	_ = v1656
	var v1660 int32
	_ = v1660
	var v1666 int32
	_ = v1666
	var v1670 int32
	_ = v1670
	var v1672 int32
	_ = v1672
	var v1675 int32
	_ = v1675
	var v1677 int32
	_ = v1677
	var v1680 int32
	_ = v1680
	var v1682 int32
	_ = v1682
	var v1686 int32
	_ = v1686
	var v1690 int32
	_ = v1690
	var v1691 int32
	_ = v1691
	var v1694 int32
	_ = v1694
	var v1698 int32
	_ = v1698
	var v1702 int32
	_ = v1702
	var v1710 int32
	_ = v1710
	var v1722 int32
	_ = v1722
	var v1725 int32
	_ = v1725
	var v1727 int32
	_ = v1727
	var v1731 int32
	_ = v1731
	var v1733 int32
	_ = v1733
	var v1737 int32
	_ = v1737
	var v1741 int32
	_ = v1741
	var v1742 int32
	_ = v1742
	var v1745 int32
	_ = v1745
	var v1749 int32
	_ = v1749
	var v1753 int32
	_ = v1753
	var v1759 int32
	_ = v1759
	var v1763 int32
	_ = v1763
	var v1764 int32
	_ = v1764
	var v1767 int32
	_ = v1767
	var v1769 int32
	_ = v1769
	var v1772 int32
	_ = v1772
	var v1774 int32
	_ = v1774
	var v1778 int32
	_ = v1778
	var v1782 int32
	_ = v1782
	var v1783 int32
	_ = v1783
	var v1786 int32
	_ = v1786
	var v1790 int32
	_ = v1790
	var v1794 int32
	_ = v1794
	var v1800 int32
	_ = v1800
	var v1812 int32
	_ = v1812
	var v1815 int32
	_ = v1815
	var v1817 int32
	_ = v1817
	var v1821 int32
	_ = v1821
	var v1823 int32
	_ = v1823
	var v1827 int32
	_ = v1827
	var v1831 int32
	_ = v1831
	var v1832 int32
	_ = v1832
	var v1835 int32
	_ = v1835
	var v1839 int32
	_ = v1839
	var v1843 int32
	_ = v1843
	var v1849 int32
	_ = v1849
	var v1854 int32
	_ = v1854
	var v1866 int32
	_ = v1866
	var v1869 int32
	_ = v1869
	var v1871 int32
	_ = v1871
	var v1875 int32
	_ = v1875
	var v1877 int32
	_ = v1877
	var v1881 int32
	_ = v1881
	var v1885 int32
	_ = v1885
	var v1886 int32
	_ = v1886
	var v1889 int32
	_ = v1889
	var v1893 int32
	_ = v1893
	var v1897 int32
	_ = v1897
	var v1903 int32
	_ = v1903
	var v1907 int32
	_ = v1907
	var v1908 int32
	_ = v1908
	var v1911 int32
	_ = v1911
	var v1913 int32
	_ = v1913
	var v1916 int32
	_ = v1916
	var v1918 int32
	_ = v1918
	var v1922 int32
	_ = v1922
	var v1926 int32
	_ = v1926
	var v1927 int32
	_ = v1927
	var v1930 int32
	_ = v1930
	var v1934 int32
	_ = v1934
	var v1938 int32
	_ = v1938
	var v1946 int32
	_ = v1946
	var v1947 int32
	_ = v1947
	var v1948 int32
	_ = v1948
	var v1951 int32
	_ = v1951
	var v1953 int32
	_ = v1953
	var v1956 int32
	_ = v1956
	var v1958 int32
	_ = v1958
	var v1962 int32
	_ = v1962
	var v1966 int32
	_ = v1966
	var v1967 int32
	_ = v1967
	var v1970 int32
	_ = v1970
	var v1974 int32
	_ = v1974
	var v1978 int32
	_ = v1978
	var v1998 int32
	_ = v1998
	var v2001 int32
	_ = v2001
	var v2003 int32
	_ = v2003
	var v2007 int32
	_ = v2007
	var v2009 int32
	_ = v2009
	var v2013 int32
	_ = v2013
	var v2017 int32
	_ = v2017
	var v2018 int32
	_ = v2018
	var v2021 int32
	_ = v2021
	var v2025 int32
	_ = v2025
	var v2029 int32
	_ = v2029
	var v2035 int32
	_ = v2035
	var v2039 int32
	_ = v2039
	var v2044 int32
	_ = v2044
	var v2046 int32
	_ = v2046
	var v2049 int32
	_ = v2049
	var v2051 int32
	_ = v2051
	var v2055 int32
	_ = v2055
	var v2059 int32
	_ = v2059
	var v2060 int32
	_ = v2060
	var v2063 int32
	_ = v2063
	var v2067 int32
	_ = v2067
	var v2071 int32
	_ = v2071
	var v2081 int32
	_ = v2081
	var v2089 int32
	_ = v2089
	var v2093 int32
	_ = v2093
	var v2095 int32
	_ = v2095
	var v2099 int32
	_ = v2099
	var v2101 int32
	_ = v2101
	var v2105 int32
	_ = v2105
	var v2109 int32
	_ = v2109
	var v2110 int32
	_ = v2110
	var v2113 int32
	_ = v2113
	var v2117 int32
	_ = v2117
	var v2121 int32
	_ = v2121
	var v2137 int32
	_ = v2137
	var v2141 int32
	_ = v2141
	var v2143 int32
	_ = v2143
	var v2147 int32
	_ = v2147
	var v2149 int32
	_ = v2149
	var v2153 int32
	_ = v2153
	var v2157 int32
	_ = v2157
	var v2158 int32
	_ = v2158
	var v2161 int32
	_ = v2161
	var v2165 int32
	_ = v2165
	var v2169 int32
	_ = v2169
	var v2175 int32
	_ = v2175
	var v2197 int32
	_ = v2197
	var v2200 int32
	_ = v2200
	var v2202 int32
	_ = v2202
	var v2206 int32
	_ = v2206
	var v2208 int32
	_ = v2208
	var v2212 int32
	_ = v2212
	var v2216 int32
	_ = v2216
	var v2217 int32
	_ = v2217
	var v2220 int32
	_ = v2220
	var v2224 int32
	_ = v2224
	var v2228 int32
	_ = v2228
	var v2234 int32
	_ = v2234
	var v2238 int32
	_ = v2238
	var v2244 int32
	_ = v2244
	var v2256 int32
	_ = v2256
	var v2259 int32
	_ = v2259
	var v2261 int32
	_ = v2261
	var v2265 int32
	_ = v2265
	var v2267 int32
	_ = v2267
	var v2271 int32
	_ = v2271
	var v2275 int32
	_ = v2275
	var v2276 int32
	_ = v2276
	var v2279 int32
	_ = v2279
	var v2283 int32
	_ = v2283
	var v2287 int32
	_ = v2287
	var v2293 int32
	_ = v2293
	var v2297 int32
	_ = v2297
	var v2309 int32
	_ = v2309
	var v2312 int32
	_ = v2312
	var v2314 int32
	_ = v2314
	var v2318 int32
	_ = v2318
	var v2320 int32
	_ = v2320
	var v2324 int32
	_ = v2324
	var v2328 int32
	_ = v2328
	var v2329 int32
	_ = v2329
	var v2332 int32
	_ = v2332
	var v2336 int32
	_ = v2336
	var v2340 int32
	_ = v2340
	var v2346 int32
	_ = v2346
	var v2350 int32
	_ = v2350
	var v2355 int32
	_ = v2355
	var v2357 int32
	_ = v2357
	var v2360 int32
	_ = v2360
	var v2362 int32
	_ = v2362
	var v2366 int32
	_ = v2366
	var v2370 int32
	_ = v2370
	var v2371 int32
	_ = v2371
	var v2374 int32
	_ = v2374
	var v2378 int32
	_ = v2378
	var v2382 int32
	_ = v2382
	var v2392 int32
	_ = v2392
	var v2400 int32
	_ = v2400
	var v2404 int32
	_ = v2404
	var v2406 int32
	_ = v2406
	var v2410 int32
	_ = v2410
	var v2412 int32
	_ = v2412
	var v2416 int32
	_ = v2416
	var v2420 int32
	_ = v2420
	var v2421 int32
	_ = v2421
	var v2424 int32
	_ = v2424
	var v2428 int32
	_ = v2428
	var v2432 int32
	_ = v2432
	var v2448 int32
	_ = v2448
	var v2452 int32
	_ = v2452
	var v2454 int32
	_ = v2454
	var v2458 int32
	_ = v2458
	var v2460 int32
	_ = v2460
	var v2464 int32
	_ = v2464
	var v2468 int32
	_ = v2468
	var v2469 int32
	_ = v2469
	var v2472 int32
	_ = v2472
	var v2476 int32
	_ = v2476
	var v2480 int32
	_ = v2480
	var v2486 int32
	_ = v2486
	var v2496 int32
	_ = v2496
	var v2501 int32
	_ = v2501
	var v2503 int32
	_ = v2503
	var v2506 int32
	_ = v2506
	var v2508 int32
	_ = v2508
	var v2512 int32
	_ = v2512
	var v2516 int32
	_ = v2516
	var v2517 int32
	_ = v2517
	var v2520 int32
	_ = v2520
	var v2524 int32
	_ = v2524
	var v2528 int32
	_ = v2528
	var v2538 int32
	_ = v2538
	var v2546 int32
	_ = v2546
	var v2550 int32
	_ = v2550
	var v2552 int32
	_ = v2552
	var v2556 int32
	_ = v2556
	var v2558 int32
	_ = v2558
	var v2562 int32
	_ = v2562
	var v2566 int32
	_ = v2566
	var v2567 int32
	_ = v2567
	var v2570 int32
	_ = v2570
	var v2574 int32
	_ = v2574
	var v2578 int32
	_ = v2578
	var v2594 int32
	_ = v2594
	var v2598 int32
	_ = v2598
	var v2600 int32
	_ = v2600
	var v2604 int32
	_ = v2604
	var v2606 int32
	_ = v2606
	var v2610 int32
	_ = v2610
	var v2614 int32
	_ = v2614
	var v2615 int32
	_ = v2615
	var v2618 int32
	_ = v2618
	var v2622 int32
	_ = v2622
	var v2626 int32
	_ = v2626
	var v2632 int32
	_ = v2632
	var v2642 int32
	_ = v2642
	var v2647 int32
	_ = v2647
	var v2649 int32
	_ = v2649
	var v2652 int32
	_ = v2652
	var v2654 int32
	_ = v2654
	var v2658 int32
	_ = v2658
	var v2662 int32
	_ = v2662
	var v2663 int32
	_ = v2663
	var v2666 int32
	_ = v2666
	var v2670 int32
	_ = v2670
	var v2674 int32
	_ = v2674
	var v2684 int32
	_ = v2684
	var v2692 int32
	_ = v2692
	var v2696 int32
	_ = v2696
	var v2698 int32
	_ = v2698
	var v2702 int32
	_ = v2702
	var v2704 int32
	_ = v2704
	var v2708 int32
	_ = v2708
	var v2712 int32
	_ = v2712
	var v2713 int32
	_ = v2713
	var v2716 int32
	_ = v2716
	var v2720 int32
	_ = v2720
	var v2724 int32
	_ = v2724
	var v2740 int32
	_ = v2740
	var v2744 int32
	_ = v2744
	var v2746 int32
	_ = v2746
	var v2750 int32
	_ = v2750
	var v2752 int32
	_ = v2752
	var v2756 int32
	_ = v2756
	var v2760 int32
	_ = v2760
	var v2761 int32
	_ = v2761
	var v2764 int32
	_ = v2764
	var v2768 int32
	_ = v2768
	var v2772 int32
	_ = v2772
	var v2778 int32
	_ = v2778
	var v2788 int32
	_ = v2788
	var v2793 int32
	_ = v2793
	var v2795 int32
	_ = v2795
	var v2798 int32
	_ = v2798
	var v2800 int32
	_ = v2800
	var v2804 int32
	_ = v2804
	var v2808 int32
	_ = v2808
	var v2809 int32
	_ = v2809
	var v2812 int32
	_ = v2812
	var v2816 int32
	_ = v2816
	var v2820 int32
	_ = v2820
	var v2830 int32
	_ = v2830
	var v2838 int32
	_ = v2838
	var v2842 int32
	_ = v2842
	var v2844 int32
	_ = v2844
	var v2848 int32
	_ = v2848
	var v2850 int32
	_ = v2850
	var v2854 int32
	_ = v2854
	var v2858 int32
	_ = v2858
	var v2859 int32
	_ = v2859
	var v2862 int32
	_ = v2862
	var v2866 int32
	_ = v2866
	var v2870 int32
	_ = v2870
	var v2886 int32
	_ = v2886
	var v2890 int32
	_ = v2890
	var v2892 int32
	_ = v2892
	var v2896 int32
	_ = v2896
	var v2898 int32
	_ = v2898
	var v2902 int32
	_ = v2902
	var v2906 int32
	_ = v2906
	var v2907 int32
	_ = v2907
	var v2910 int32
	_ = v2910
	var v2914 int32
	_ = v2914
	var v2918 int32
	_ = v2918
	var v2924 int32
	_ = v2924
	var v2934 int32
	_ = v2934
	var v2939 int32
	_ = v2939
	var v2941 int32
	_ = v2941
	var v2944 int32
	_ = v2944
	var v2946 int32
	_ = v2946
	var v2950 int32
	_ = v2950
	var v2954 int32
	_ = v2954
	var v2955 int32
	_ = v2955
	var v2958 int32
	_ = v2958
	var v2962 int32
	_ = v2962
	var v2966 int32
	_ = v2966
	var v2976 int32
	_ = v2976
	var v2984 int32
	_ = v2984
	var v2988 int32
	_ = v2988
	var v2990 int32
	_ = v2990
	var v2994 int32
	_ = v2994
	var v2996 int32
	_ = v2996
	var v3000 int32
	_ = v3000
	var v3004 int32
	_ = v3004
	var v3005 int32
	_ = v3005
	var v3008 int32
	_ = v3008
	var v3012 int32
	_ = v3012
	var v3016 int32
	_ = v3016
	var v3032 int32
	_ = v3032
	var v3036 int32
	_ = v3036
	var v3038 int32
	_ = v3038
	var v3042 int32
	_ = v3042
	var v3044 int32
	_ = v3044
	var v3048 int32
	_ = v3048
	var v3052 int32
	_ = v3052
	var v3053 int32
	_ = v3053
	var v3056 int32
	_ = v3056
	var v3060 int32
	_ = v3060
	var v3064 int32
	_ = v3064
	var v3070 int32
	_ = v3070
	var v3083 int32
	_ = v3083
	var v3085 int32
	_ = v3085
	var v3098 int32
	_ = v3098
	var v3099 int32
	_ = v3099
	var v3102 int32
	_ = v3102
	var v3106 int32
	_ = v3106
	var v3110 int32
	_ = v3110
	var v3118 int32
	_ = v3118
	var v3119 int32
	_ = v3119
	var v3120 int32
	_ = v3120
	var v3121 int32
	_ = v3121
	var v3128 int32
	_ = v3128
	var v3132 int32
	_ = v3132
	var v3135 int32
	_ = v3135
	var v3137 int32
	_ = v3137
	var v3141 int32
	_ = v3141
	var v3142 int32
	_ = v3142
	var v3145 int32
	_ = v3145
	var v3149 int32
	_ = v3149
	var v3153 int32
	_ = v3153
	var v3159 int32
	_ = v3159
	var v3167 int32
	_ = v3167
	var v3168 int32
	_ = v3168
	var v3171 int32
	_ = v3171
	var v3172 int32
	_ = v3172
	var v3177 int32
	_ = v3177
	var v3185 int32
	_ = v3185
	var v3188 int64
	_ = v3188
	var v3190 int64
	_ = v3190
	var v3192 int64
	_ = v3192
	var v3193 int64
	_ = v3193
	var v3199 int32
	_ = v3199
	var v3208 int64
	_ = v3208
	var v3210 int32
	_ = v3210
	var v3215 int32
	_ = v3215
	var v3219 int32
	_ = v3219
	var v3222 int32
	_ = v3222
	var v3244 int32
	_ = v3244
	var v3245 int32
	_ = v3245
	var v3250 int32
	_ = v3250
	var v3254 int32
	_ = v3254
	var v3268 int32
	_ = v3268
	var v3272 int32
	_ = v3272
	var v3274 int32
	_ = v3274
	var v3281 int32
	_ = v3281
	var v3284 int32
	_ = v3284
	var v3287 int32
	_ = v3287
	var v3290 int32
	_ = v3290
	var v3294 int32
	_ = v3294
	var v3298 int32
	_ = v3298
	var v3305 int32
	_ = v3305
	var v3311 int32
	_ = v3311
	var v3329 int32
	_ = v3329
	var v3334 int32
	_ = v3334
	var v3336 int32
	_ = v3336
	var v3343 int32
	_ = v3343
	var v3344 int32
	_ = v3344
	var v3348 int32
	_ = v3348
	var v3360 int32
	_ = v3360
	var v3371 int32
	_ = v3371
	var v3372 int32
	_ = v3372
	var v3374 int32
	_ = v3374
	var v3375 int32
	_ = v3375
	var v3376 int32
	_ = v3376
	var v3379 int32
	_ = v3379
	var v3386 int32
	_ = v3386
	var v3391 int32
	_ = v3391
	var v3396 int32
	_ = v3396
	var v3397 int32
	_ = v3397
	var v3399 int32
	_ = v3399
	var v3408 int32
	_ = v3408
	var v3410 int64
	_ = v3410
	var v3416 int32
	_ = v3416
	var v3417 int32
	_ = v3417
	var v3422 int32
	_ = v3422
	var v3425 int32
	_ = v3425
	var v3426 int32
	_ = v3426
	var v3427 int32
	_ = v3427
	var v3432 int64
	_ = v3432
	var v3434 int64
	_ = v3434
	var v3436 int32
	_ = v3436
	var v3438 int32
	_ = v3438
	var v3439 int32
	_ = v3439
	var v3440 int32
	_ = v3440
	var v3443 int32
	_ = v3443
	var v3446 int32
	_ = v3446
	var v3450 int32
	_ = v3450
	var v3452 int32
	_ = v3452
	var v3455 int32
	_ = v3455
	var v3458 int32
	_ = v3458
	var v3463 int32
	_ = v3463
	var v3464 int32
	_ = v3464
	var v3469 int32
	_ = v3469
	var v3472 int32
	_ = v3472
	var v3475 int32
	_ = v3475
	var v3477 int32
	_ = v3477
	var v3481 int32
	_ = v3481
	var v3482 int32
	_ = v3482
	var v3485 int32
	_ = v3485
	var v3486 int32
	_ = v3486
	var v3487 int32
	_ = v3487
	var v3488 int32
	_ = v3488
	var v3491 int32
	_ = v3491
	var v3496 int32
	_ = v3496
	var v3501 int32
	_ = v3501
	var v3502 int32
	_ = v3502
	var v3512 int32
	_ = v3512
	var v3514 int32
	_ = v3514
	var v3515 int32
	_ = v3515
	var v3517 int32
	_ = v3517
	var v3519 int32
	_ = v3519
	var v3524 int32
	_ = v3524
	var v3525 int32
	_ = v3525
	var v3531 int32
	_ = v3531
	var v3532 int32
	_ = v3532
	var v3534 int32
	_ = v3534
	var v3536 int32
	_ = v3536
	var v3537 int32
	_ = v3537
	var v3544 int32
	_ = v3544
	var v3545 int32
	_ = v3545
	var v3546 int32
	_ = v3546
	var v3548 int32
	_ = v3548
	var v3550 int32
	_ = v3550
	var v3553 int32
	_ = v3553
	var v3558 int32
	_ = v3558
	var v3559 int32
	_ = v3559
	var v3562 int32
	_ = v3562
	var v3564 int32
	_ = v3564
	var v3572 int32
	_ = v3572
	var v3573 int32
	_ = v3573
	var v3575 int32
	_ = v3575
	var v3584 int32
	_ = v3584
	var v3589 int32
	_ = v3589
	var v3590 int32
	_ = v3590
	var v3594 int32
	_ = v3594
	var v3603 int32
	_ = v3603
	var v3607 int32
	_ = v3607
	var v3614 int32
	_ = v3614
	var v3619 int32
	_ = v3619
	var v3626 int32
	_ = v3626
	var v3631 int32
	_ = v3631
	var v3632 int32
	_ = v3632
	var v3637 int32
	_ = v3637
	var v3649 int32
	_ = v3649
	var v3671 int32
	_ = v3671
	var v3679 int32
	_ = v3679
	var v3684 int32
	_ = v3684
	var v3685 int32
	_ = v3685
	var v3689 int64
	_ = v3689
	var v3700 int32
	_ = v3700
	var v3710 int32
	_ = v3710
	var v3715 int32
	_ = v3715
	var v3717 int32
	_ = v3717
	var v3724 int32
	_ = v3724
	var v3730 int32
	_ = v3730
	var v3731 int32
	_ = v3731
	var v3732 int32
	_ = v3732
	var v3735 int32
	_ = v3735
	var v3737 int32
	_ = v3737
	var v3742 int32
	_ = v3742
	var v3743 int32
	_ = v3743
	var v3747 int64
	_ = v3747
	var v3759 int32
	_ = v3759
	var v3761 int32
	_ = v3761
	var v3764 int32
	_ = v3764
	var v3768 int32
	_ = v3768
	var v3769 int32
	_ = v3769
	var v3782 int32
	_ = v3782
	var v3793 int32
	_ = v3793
	var v3798 int32
	_ = v3798
	var v3799 int32
	_ = v3799
	var v3802 int32
	_ = v3802
	var v3810 int32
	_ = v3810
	var v3815 int32
	_ = v3815
	var v3837 int32
	_ = v3837
	v18 = int64(0)
	v19 = m.G0
	v21 = v19 - int32(48)
	m.G0 = v21
	v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)+368))
	v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v26 = base.I32_div_s(int32(19), v25)
	v29 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(76))))
	v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v34 = l0 + int32(56)
	v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	v41 = base.I32_div_s(v35*v36*int32(7), int32(8))
	*(*int64)(unsafe.Add(mBase, uint32(v34)+16)) = v18
	*(*int64)(unsafe.Add(mBase, uint32(v34)+8)) = int64(-34359738368)
	*(*int64)(unsafe.Add(mBase, uint32(v34))) = int64(254)
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(80)))) = v18
	if v41 == int32(0) {
		goto L7
	} else {
		goto L8
	}
L1:
	;
	m.G0 = v21 + int32(48)
	return v3837
L2:
	;
	v3837 = int32(0)
	goto L1
L3:
	;
	v3244 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(76))))
	v3245 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	v3250 = v3244 + v3245*int32(3) + int32(7)
	if v3245 < int32(1) {
		v3360 = v3250
		goto L573
	} else {
		goto L574
	}
L4:
	;
	if v3222 == int32(0) {
		goto L2
	} else {
		goto L572
	}
L5:
	;
	v92 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v94 = v92 >> (uint(int32(1)) % 32)
	goto L21
L6:
	;
	if v83 != 0 {
		goto L5
	} else {
		goto L16
	}
L7:
	;
	v83 = int32(1)
	goto L6
L8:
	;
	v57 = int32(1024)
	if base.Ui32(v57) < base.Ui32(v41) {
		goto L10
	} else {
		goto L11
	}
L9:
	;
	v65 = *(*int32)(unsafe.Add(mBase, uint32(v34)+20))
	if v65 == int32(0) {
		goto L14
	} else {
		goto L15
	}
L10:
	;
	v60 = v41
	goto L12
L11:
	;
	v60 = v57
	goto L12
L12:
	;
	v61 = F_WebPSafeMalloc(m, int64(1), v60)
	mBase = m.M
	if v61 != 0 {
		goto L9
	} else {
		goto L13
	}
L13:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34)+28)) = int32(1)
	v83 = int32(0)
	goto L6
L14:
	;
	v72 = *(*int32)(unsafe.Add(mBase, uint32(v34)+16))
	F_WebPSafeFree(m, v72)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v34)+24)) = v60
	*(*int32)(unsafe.Add(mBase, uint32(v34)+16)) = v61
	goto L7
L15:
	;
	v70 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(72))))
	v71 = F_memcpy(m, v61, v70, v65)
	mBase = m.M
	goto L14
L16:
	;
	v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v86 = *(*int32)(unsafe.Add(mBase, uint32(v84)+92))
	if v86 != 0 {
		goto L18
	} else {
		goto L19
	}
L17:
	;
	v3222 = int32(0)
	goto L4
L18:
	;
	goto L17
L19:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v84)+92)) = int32(1)
	goto L18
L20:
	;
	v128 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v130 = v128 >> (uint(int32(1)) % 32)
	goto L27
L21:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v94
	if int32(126) < v94 {
		goto L23
	} else {
		goto L24
	}
L23:
	;
	goto L20
L24:
	;
	v107 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v108 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v107 << (uint(v108) % 32)
	v111 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v111 + v108
	v115 = m.G1
	v119 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v115+int32(_a_F_VP8EncWrite_0)+v94))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v119
	if v111 < int32(0) {
		goto L23
	} else {
		goto L25
	}
L25:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L23
L26:
	;
	v162 = l0 + int32(3416)
	v163 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	v164 = int32(1)
	v165 = base.B2i32(v164 < v163)
	v168 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v170 = v168 >> (uint(v164) % 32)
	if v165 == int32(0) {
		v179 = v170
		goto L34
	} else {
		goto L35
	}
L27:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v130
	if int32(126) < v130 {
		goto L29
	} else {
		goto L30
	}
L29:
	;
	goto L26
L30:
	;
	v143 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v144 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v143 << (uint(v144) % 32)
	v147 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v147 + v144
	v151 = m.G1
	v155 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v151+int32(_a_F_VP8EncWrite_0)+v130))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v155
	if v147 < int32(0) {
		goto L29
	} else {
		goto L31
	}
L31:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L29
L32:
	;
	v1763 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v1764 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v1767 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1769 = v1767 >> (uint(int32(1)) % 32)
	if v1764 == int32(0) {
		v1778 = v1769
		goto L320
	} else {
		goto L321
	}
L33:
	;
	if v165 == int32(0) {
		goto L32
	} else {
		goto L39
	}
L34:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v179
	if int32(126) < v179 {
		goto L36
	} else {
		goto L37
	}
L35:
	;
	v173 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v175 = v170 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v173 + v175
	v179 = v168 - v175
	goto L34
L36:
	;
	goto L33
L37:
	;
	v183 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v184 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v183 << (uint(v184) % 32)
	v187 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v187 + v184
	v191 = m.G1
	v195 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v191+int32(_a_F_VP8EncWrite_0)+v179))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v195
	if v187 < int32(0) {
		goto L36
	} else {
		goto L38
	}
L38:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L36
L39:
	;
	v203 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	v206 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v208 = v206 >> (uint(int32(1)) % 32)
	if v203 == int32(0) {
		v217 = v208
		goto L41
	} else {
		goto L42
	}
L40:
	;
	v242 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	goto L49
L41:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v217
	if int32(126) < v217 {
		goto L43
	} else {
		goto L44
	}
L42:
	;
	v211 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v213 = v208 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v211 + v213
	v217 = v206 - v213
	goto L41
L43:
	;
	goto L40
L44:
	;
	v221 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v222 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v221 << (uint(v222) % 32)
	v225 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v225 + v222
	v229 = m.G1
	v233 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v229+int32(_a_F_VP8EncWrite_0)+v217))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v233
	if v225 < int32(0) {
		goto L43
	} else {
		goto L45
	}
L45:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L43
L46:
	;
	v1481 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	if v1481 == int32(0) {
		goto L32
	} else {
		goto L268
	}
L47:
	;
	goto L53
L48:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v252
	if int32(126) < v252 {
		goto L50
	} else {
		goto L51
	}
L49:
	;
	v247 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v249 = v242>>(uint(int32(1))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v247 + v249
	v252 = v242 - v249
	goto L48
L50:
	;
	goto L47
L51:
	;
	v257 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v258 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v257 << (uint(v258) % 32)
	v261 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v261 + v258
	v265 = m.G1
	v269 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v265+int32(_a_F_VP8EncWrite_0)+v252))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v269
	if v261 < int32(0) {
		goto L50
	} else {
		goto L52
	}
L52:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L50
L53:
	;
	v280 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	goto L56
L54:
	;
	v313 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1088))
	v318 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v320 = v318 >> (uint(int32(1)) % 32)
	if v313 == int32(0) {
		v329 = v320
		goto L61
	} else {
		goto L62
	}
L55:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v290
	if int32(126) < v290 {
		goto L57
	} else {
		goto L58
	}
L56:
	;
	v285 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v287 = v280>>(uint(int32(1))%32) + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v285 + v287
	v290 = v280 - v287
	goto L55
L57:
	;
	goto L54
L58:
	;
	v295 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v296 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v295 << (uint(v296) % 32)
	v299 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v299 + v296
	v303 = m.G1
	v307 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v303+int32(_a_F_VP8EncWrite_0)+v290))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v307
	if v299 < int32(0) {
		goto L57
	} else {
		goto L59
	}
L59:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L57
L60:
	;
	v459 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1832))
	v464 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v466 = v464 >> (uint(int32(1)) % 32)
	if v459 == int32(0) {
		v475 = v466
		goto L87
	} else {
		goto L88
	}
L61:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v329
	if int32(126) < v329 {
		goto L63
	} else {
		goto L64
	}
L62:
	;
	v323 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v325 = v320 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v323 + v325
	v329 = v318 - v325
	goto L61
L63:
	;
	if v313 == int32(0) {
		goto L66
	} else {
		goto L67
	}
L64:
	;
	v333 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v334 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v333 << (uint(v334) % 32)
	v337 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v337 + v334
	v341 = m.G1
	v345 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v341+int32(_a_F_VP8EncWrite_0)+v329))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v345
	if v337 < int32(0) {
		goto L63
	} else {
		goto L65
	}
L65:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L63
L66:
	;
	goto L60
L67:
	;
	if int32(-1) < v313 {
		goto L68
	} else {
		goto L69
	}
L68:
	;
	v411 = int32(128)
	goto L78
L69:
	;
	v355 = int32(1)
	v363 = int32(128)
	goto L70
L70:
	;
	v367 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v369 = v367 >> (uint(int32(1)) % 32)
	if v363&(v355-v313<<(uint(v355)%32)) == int32(0) {
		v379 = v369
		goto L72
	} else {
		goto L73
	}
L72:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v379
	if int32(126) < v379 {
		goto L74
	} else {
		goto L75
	}
L73:
	;
	v373 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v375 = v369 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v373 + v375
	v379 = v367 - v375
	goto L72
L74:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v363) {
		v363 = int32(base.Ui32(v363) >> (uint(int32(1)) % 32))
		goto L70
	} else {
		goto L77
	}
L75:
	;
	v383 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v384 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v383 << (uint(v384) % 32)
	v387 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v387 + v384
	v391 = m.G1
	v395 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v391+int32(_a_F_VP8EncWrite_0)+v379))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v395
	if v387 < int32(0) {
		goto L74
	} else {
		goto L76
	}
L76:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L74
L77:
	;
	goto L66
L78:
	;
	v415 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v417 = v415 >> (uint(int32(1)) % 32)
	if v411&(v313<<(uint(int32(1))%32)) == int32(0) {
		v427 = v417
		goto L80
	} else {
		goto L81
	}
L79:
	;
	goto L66
L80:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v427
	if int32(126) < v427 {
		goto L82
	} else {
		goto L83
	}
L81:
	;
	v421 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v423 = v417 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v421 + v423
	v427 = v415 - v423
	goto L80
L82:
	;
	v449 = int32(1)
	if base.Ui32(v449) < base.Ui32(v411) {
		v411 = int32(base.Ui32(v411) >> (uint(v449) % 32))
		goto L78
	} else {
		goto L85
	}
L83:
	;
	v431 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v432 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v431 << (uint(v432) % 32)
	v435 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v435 + v432
	v439 = m.G1
	v443 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v439+int32(_a_F_VP8EncWrite_0)+v427))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v443
	if v435 < int32(0) {
		goto L82
	} else {
		goto L84
	}
L84:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L82
L85:
	;
	goto L79
L86:
	;
	v605 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2576))
	v610 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v612 = v610 >> (uint(int32(1)) % 32)
	if v605 == int32(0) {
		v621 = v612
		goto L113
	} else {
		goto L114
	}
L87:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v475
	if int32(126) < v475 {
		goto L89
	} else {
		goto L90
	}
L88:
	;
	v469 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v471 = v466 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v469 + v471
	v475 = v464 - v471
	goto L87
L89:
	;
	if v459 == int32(0) {
		goto L92
	} else {
		goto L93
	}
L90:
	;
	v479 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v480 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v479 << (uint(v480) % 32)
	v483 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v483 + v480
	v487 = m.G1
	v491 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v487+int32(_a_F_VP8EncWrite_0)+v475))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v491
	if v483 < int32(0) {
		goto L89
	} else {
		goto L91
	}
L91:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L89
L92:
	;
	goto L86
L93:
	;
	if int32(-1) < v459 {
		goto L94
	} else {
		goto L95
	}
L94:
	;
	v557 = int32(128)
	goto L104
L95:
	;
	v501 = int32(1)
	v509 = int32(128)
	goto L96
L96:
	;
	v513 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v515 = v513 >> (uint(int32(1)) % 32)
	if v509&(v501-v459<<(uint(v501)%32)) == int32(0) {
		v525 = v515
		goto L98
	} else {
		goto L99
	}
L98:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v525
	if int32(126) < v525 {
		goto L100
	} else {
		goto L101
	}
L99:
	;
	v519 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v521 = v515 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v519 + v521
	v525 = v513 - v521
	goto L98
L100:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v509) {
		v509 = int32(base.Ui32(v509) >> (uint(int32(1)) % 32))
		goto L96
	} else {
		goto L103
	}
L101:
	;
	v529 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v530 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v529 << (uint(v530) % 32)
	v533 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v533 + v530
	v537 = m.G1
	v541 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v537+int32(_a_F_VP8EncWrite_0)+v525))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v541
	if v533 < int32(0) {
		goto L100
	} else {
		goto L102
	}
L102:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L100
L103:
	;
	goto L92
L104:
	;
	v561 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v563 = v561 >> (uint(int32(1)) % 32)
	if v557&(v459<<(uint(int32(1))%32)) == int32(0) {
		v573 = v563
		goto L106
	} else {
		goto L107
	}
L105:
	;
	goto L92
L106:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v573
	if int32(126) < v573 {
		goto L108
	} else {
		goto L109
	}
L107:
	;
	v567 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v569 = v563 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v567 + v569
	v573 = v561 - v569
	goto L106
L108:
	;
	v595 = int32(1)
	if base.Ui32(v595) < base.Ui32(v557) {
		v557 = int32(base.Ui32(v557) >> (uint(v595) % 32))
		goto L104
	} else {
		goto L111
	}
L109:
	;
	v577 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v578 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v577 << (uint(v578) % 32)
	v581 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v581 + v578
	v585 = m.G1
	v589 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v585+int32(_a_F_VP8EncWrite_0)+v573))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v589
	if v581 < int32(0) {
		goto L108
	} else {
		goto L110
	}
L110:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L108
L111:
	;
	goto L105
L112:
	;
	v751 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3320))
	v756 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v758 = v756 >> (uint(int32(1)) % 32)
	if v751 == int32(0) {
		v767 = v758
		goto L139
	} else {
		goto L140
	}
L113:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v621
	if int32(126) < v621 {
		goto L115
	} else {
		goto L116
	}
L114:
	;
	v615 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v617 = v612 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v615 + v617
	v621 = v610 - v617
	goto L113
L115:
	;
	if v605 == int32(0) {
		goto L118
	} else {
		goto L119
	}
L116:
	;
	v625 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v626 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v625 << (uint(v626) % 32)
	v629 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v629 + v626
	v633 = m.G1
	v637 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v633+int32(_a_F_VP8EncWrite_0)+v621))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v637
	if v629 < int32(0) {
		goto L115
	} else {
		goto L117
	}
L117:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L115
L118:
	;
	goto L112
L119:
	;
	if int32(-1) < v605 {
		goto L120
	} else {
		goto L121
	}
L120:
	;
	v703 = int32(128)
	goto L130
L121:
	;
	v647 = int32(1)
	v655 = int32(128)
	goto L122
L122:
	;
	v659 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v661 = v659 >> (uint(int32(1)) % 32)
	if v655&(v647-v605<<(uint(v647)%32)) == int32(0) {
		v671 = v661
		goto L124
	} else {
		goto L125
	}
L124:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v671
	if int32(126) < v671 {
		goto L126
	} else {
		goto L127
	}
L125:
	;
	v665 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v667 = v661 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v665 + v667
	v671 = v659 - v667
	goto L124
L126:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v655) {
		v655 = int32(base.Ui32(v655) >> (uint(int32(1)) % 32))
		goto L122
	} else {
		goto L129
	}
L127:
	;
	v675 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v676 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v675 << (uint(v676) % 32)
	v679 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v679 + v676
	v683 = m.G1
	v687 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v683+int32(_a_F_VP8EncWrite_0)+v671))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v687
	if v679 < int32(0) {
		goto L126
	} else {
		goto L128
	}
L128:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L126
L129:
	;
	goto L118
L130:
	;
	v707 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v709 = v707 >> (uint(int32(1)) % 32)
	if v703&(v605<<(uint(int32(1))%32)) == int32(0) {
		v719 = v709
		goto L132
	} else {
		goto L133
	}
L131:
	;
	goto L118
L132:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v719
	if int32(126) < v719 {
		goto L134
	} else {
		goto L135
	}
L133:
	;
	v713 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v715 = v709 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v713 + v715
	v719 = v707 - v715
	goto L132
L134:
	;
	v741 = int32(1)
	if base.Ui32(v741) < base.Ui32(v703) {
		v703 = int32(base.Ui32(v703) >> (uint(v741) % 32))
		goto L130
	} else {
		goto L137
	}
L135:
	;
	v723 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v724 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v723 << (uint(v724) % 32)
	v727 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v727 + v724
	v731 = m.G1
	v735 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v731+int32(_a_F_VP8EncWrite_0)+v719))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v735
	if v727 < int32(0) {
		goto L134
	} else {
		goto L136
	}
L136:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L134
L137:
	;
	goto L131
L138:
	;
	v897 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1092))
	v902 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v904 = v902 >> (uint(int32(1)) % 32)
	if v897 == int32(0) {
		v913 = v904
		goto L165
	} else {
		goto L166
	}
L139:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v767
	if int32(126) < v767 {
		goto L141
	} else {
		goto L142
	}
L140:
	;
	v761 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v763 = v758 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v761 + v763
	v767 = v756 - v763
	goto L139
L141:
	;
	if v751 == int32(0) {
		goto L144
	} else {
		goto L145
	}
L142:
	;
	v771 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v772 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v771 << (uint(v772) % 32)
	v775 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v775 + v772
	v779 = m.G1
	v783 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v779+int32(_a_F_VP8EncWrite_0)+v767))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v783
	if v775 < int32(0) {
		goto L141
	} else {
		goto L143
	}
L143:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L141
L144:
	;
	goto L138
L145:
	;
	if int32(-1) < v751 {
		goto L146
	} else {
		goto L147
	}
L146:
	;
	v849 = int32(128)
	goto L156
L147:
	;
	v793 = int32(1)
	v801 = int32(128)
	goto L148
L148:
	;
	v805 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v807 = v805 >> (uint(int32(1)) % 32)
	if v801&(v793-v751<<(uint(v793)%32)) == int32(0) {
		v817 = v807
		goto L150
	} else {
		goto L151
	}
L150:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v817
	if int32(126) < v817 {
		goto L152
	} else {
		goto L153
	}
L151:
	;
	v811 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v813 = v807 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v811 + v813
	v817 = v805 - v813
	goto L150
L152:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v801) {
		v801 = int32(base.Ui32(v801) >> (uint(int32(1)) % 32))
		goto L148
	} else {
		goto L155
	}
L153:
	;
	v821 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v822 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v821 << (uint(v822) % 32)
	v825 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v825 + v822
	v829 = m.G1
	v833 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v829+int32(_a_F_VP8EncWrite_0)+v817))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v833
	if v825 < int32(0) {
		goto L152
	} else {
		goto L154
	}
L154:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L152
L155:
	;
	goto L144
L156:
	;
	v853 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v855 = v853 >> (uint(int32(1)) % 32)
	if v849&(v751<<(uint(int32(1))%32)) == int32(0) {
		v865 = v855
		goto L158
	} else {
		goto L159
	}
L157:
	;
	goto L144
L158:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v865
	if int32(126) < v865 {
		goto L160
	} else {
		goto L161
	}
L159:
	;
	v859 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v861 = v855 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v859 + v861
	v865 = v853 - v861
	goto L158
L160:
	;
	v887 = int32(1)
	if base.Ui32(v887) < base.Ui32(v849) {
		v849 = int32(base.Ui32(v849) >> (uint(v887) % 32))
		goto L156
	} else {
		goto L163
	}
L161:
	;
	v869 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v870 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v869 << (uint(v870) % 32)
	v873 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v873 + v870
	v877 = m.G1
	v881 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v877+int32(_a_F_VP8EncWrite_0)+v865))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v881
	if v873 < int32(0) {
		goto L160
	} else {
		goto L162
	}
L162:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L160
L163:
	;
	goto L157
L164:
	;
	v1043 = *(*int32)(unsafe.Add(mBase, uint32(l0)+1836))
	v1048 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1050 = v1048 >> (uint(int32(1)) % 32)
	if v1043 == int32(0) {
		v1059 = v1050
		goto L191
	} else {
		goto L192
	}
L165:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v913
	if int32(126) < v913 {
		goto L167
	} else {
		goto L168
	}
L166:
	;
	v907 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v909 = v904 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v907 + v909
	v913 = v902 - v909
	goto L165
L167:
	;
	if v897 == int32(0) {
		goto L170
	} else {
		goto L171
	}
L168:
	;
	v917 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v918 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v917 << (uint(v918) % 32)
	v921 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v921 + v918
	v925 = m.G1
	v929 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v925+int32(_a_F_VP8EncWrite_0)+v913))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v929
	if v921 < int32(0) {
		goto L167
	} else {
		goto L169
	}
L169:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L167
L170:
	;
	goto L164
L171:
	;
	if int32(-1) < v897 {
		goto L172
	} else {
		goto L173
	}
L172:
	;
	v995 = int32(64)
	goto L182
L173:
	;
	v939 = int32(1)
	v947 = int32(64)
	goto L174
L174:
	;
	v951 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v953 = v951 >> (uint(int32(1)) % 32)
	if v947&(v939-v897<<(uint(v939)%32)) == int32(0) {
		v963 = v953
		goto L176
	} else {
		goto L177
	}
L176:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v963
	if int32(126) < v963 {
		goto L178
	} else {
		goto L179
	}
L177:
	;
	v957 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v959 = v953 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v957 + v959
	v963 = v951 - v959
	goto L176
L178:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v947) {
		v947 = int32(base.Ui32(v947) >> (uint(int32(1)) % 32))
		goto L174
	} else {
		goto L181
	}
L179:
	;
	v967 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v968 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v967 << (uint(v968) % 32)
	v971 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v971 + v968
	v975 = m.G1
	v979 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v975+int32(_a_F_VP8EncWrite_0)+v963))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v979
	if v971 < int32(0) {
		goto L178
	} else {
		goto L180
	}
L180:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L178
L181:
	;
	goto L170
L182:
	;
	v999 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1001 = v999 >> (uint(int32(1)) % 32)
	if v995&(v897<<(uint(int32(1))%32)) == int32(0) {
		v1011 = v1001
		goto L184
	} else {
		goto L185
	}
L183:
	;
	goto L170
L184:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1011
	if int32(126) < v1011 {
		goto L186
	} else {
		goto L187
	}
L185:
	;
	v1005 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1007 = v1001 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1005 + v1007
	v1011 = v999 - v1007
	goto L184
L186:
	;
	v1033 = int32(1)
	if base.Ui32(v1033) < base.Ui32(v995) {
		v995 = int32(base.Ui32(v995) >> (uint(v1033) % 32))
		goto L182
	} else {
		goto L189
	}
L187:
	;
	v1015 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1016 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1015 << (uint(v1016) % 32)
	v1019 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1019 + v1016
	v1023 = m.G1
	v1027 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1023+int32(_a_F_VP8EncWrite_0)+v1011))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1027
	if v1019 < int32(0) {
		goto L186
	} else {
		goto L188
	}
L188:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L186
L189:
	;
	goto L183
L190:
	;
	v1189 = *(*int32)(unsafe.Add(mBase, uint32(l0)+2580))
	v1194 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1196 = v1194 >> (uint(int32(1)) % 32)
	if v1189 == int32(0) {
		v1205 = v1196
		goto L217
	} else {
		goto L218
	}
L191:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1059
	if int32(126) < v1059 {
		goto L193
	} else {
		goto L194
	}
L192:
	;
	v1053 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1055 = v1050 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1053 + v1055
	v1059 = v1048 - v1055
	goto L191
L193:
	;
	if v1043 == int32(0) {
		goto L196
	} else {
		goto L197
	}
L194:
	;
	v1063 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1064 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1063 << (uint(v1064) % 32)
	v1067 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1067 + v1064
	v1071 = m.G1
	v1075 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1071+int32(_a_F_VP8EncWrite_0)+v1059))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1075
	if v1067 < int32(0) {
		goto L193
	} else {
		goto L195
	}
L195:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L193
L196:
	;
	goto L190
L197:
	;
	if int32(-1) < v1043 {
		goto L198
	} else {
		goto L199
	}
L198:
	;
	v1141 = int32(64)
	goto L208
L199:
	;
	v1085 = int32(1)
	v1093 = int32(64)
	goto L200
L200:
	;
	v1097 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1099 = v1097 >> (uint(int32(1)) % 32)
	if v1093&(v1085-v1043<<(uint(v1085)%32)) == int32(0) {
		v1109 = v1099
		goto L202
	} else {
		goto L203
	}
L202:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1109
	if int32(126) < v1109 {
		goto L204
	} else {
		goto L205
	}
L203:
	;
	v1103 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1105 = v1099 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1103 + v1105
	v1109 = v1097 - v1105
	goto L202
L204:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v1093) {
		v1093 = int32(base.Ui32(v1093) >> (uint(int32(1)) % 32))
		goto L200
	} else {
		goto L207
	}
L205:
	;
	v1113 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1114 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1113 << (uint(v1114) % 32)
	v1117 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1117 + v1114
	v1121 = m.G1
	v1125 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1121+int32(_a_F_VP8EncWrite_0)+v1109))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1125
	if v1117 < int32(0) {
		goto L204
	} else {
		goto L206
	}
L206:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L204
L207:
	;
	goto L196
L208:
	;
	v1145 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1147 = v1145 >> (uint(int32(1)) % 32)
	if v1141&(v1043<<(uint(int32(1))%32)) == int32(0) {
		v1157 = v1147
		goto L210
	} else {
		goto L211
	}
L209:
	;
	goto L196
L210:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1157
	if int32(126) < v1157 {
		goto L212
	} else {
		goto L213
	}
L211:
	;
	v1151 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1153 = v1147 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1151 + v1153
	v1157 = v1145 - v1153
	goto L210
L212:
	;
	v1179 = int32(1)
	if base.Ui32(v1179) < base.Ui32(v1141) {
		v1141 = int32(base.Ui32(v1141) >> (uint(v1179) % 32))
		goto L208
	} else {
		goto L215
	}
L213:
	;
	v1161 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1162 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1161 << (uint(v1162) % 32)
	v1165 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1165 + v1162
	v1169 = m.G1
	v1173 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1169+int32(_a_F_VP8EncWrite_0)+v1157))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1173
	if v1165 < int32(0) {
		goto L212
	} else {
		goto L214
	}
L214:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L212
L215:
	;
	goto L209
L216:
	;
	v1335 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3324))
	v1340 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1342 = v1340 >> (uint(int32(1)) % 32)
	if v1335 == int32(0) {
		v1351 = v1342
		goto L243
	} else {
		goto L244
	}
L217:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1205
	if int32(126) < v1205 {
		goto L219
	} else {
		goto L220
	}
L218:
	;
	v1199 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1201 = v1196 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1199 + v1201
	v1205 = v1194 - v1201
	goto L217
L219:
	;
	if v1189 == int32(0) {
		goto L222
	} else {
		goto L223
	}
L220:
	;
	v1209 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1210 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1209 << (uint(v1210) % 32)
	v1213 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1213 + v1210
	v1217 = m.G1
	v1221 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1217+int32(_a_F_VP8EncWrite_0)+v1205))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1221
	if v1213 < int32(0) {
		goto L219
	} else {
		goto L221
	}
L221:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L219
L222:
	;
	goto L216
L223:
	;
	if int32(-1) < v1189 {
		goto L224
	} else {
		goto L225
	}
L224:
	;
	v1287 = int32(64)
	goto L234
L225:
	;
	v1231 = int32(1)
	v1239 = int32(64)
	goto L226
L226:
	;
	v1243 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1245 = v1243 >> (uint(int32(1)) % 32)
	if v1239&(v1231-v1189<<(uint(v1231)%32)) == int32(0) {
		v1255 = v1245
		goto L228
	} else {
		goto L229
	}
L228:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1255
	if int32(126) < v1255 {
		goto L230
	} else {
		goto L231
	}
L229:
	;
	v1249 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1251 = v1245 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1249 + v1251
	v1255 = v1243 - v1251
	goto L228
L230:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v1239) {
		v1239 = int32(base.Ui32(v1239) >> (uint(int32(1)) % 32))
		goto L226
	} else {
		goto L233
	}
L231:
	;
	v1259 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1260 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1259 << (uint(v1260) % 32)
	v1263 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1263 + v1260
	v1267 = m.G1
	v1271 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1267+int32(_a_F_VP8EncWrite_0)+v1255))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1271
	if v1263 < int32(0) {
		goto L230
	} else {
		goto L232
	}
L232:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L230
L233:
	;
	goto L222
L234:
	;
	v1291 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1293 = v1291 >> (uint(int32(1)) % 32)
	if v1287&(v1189<<(uint(int32(1))%32)) == int32(0) {
		v1303 = v1293
		goto L236
	} else {
		goto L237
	}
L235:
	;
	goto L222
L236:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1303
	if int32(126) < v1303 {
		goto L238
	} else {
		goto L239
	}
L237:
	;
	v1297 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1299 = v1293 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1297 + v1299
	v1303 = v1291 - v1299
	goto L236
L238:
	;
	v1325 = int32(1)
	if base.Ui32(v1325) < base.Ui32(v1287) {
		v1287 = int32(base.Ui32(v1287) >> (uint(v1325) % 32))
		goto L234
	} else {
		goto L241
	}
L239:
	;
	v1307 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1308 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1307 << (uint(v1308) % 32)
	v1311 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1311 + v1308
	v1315 = m.G1
	v1319 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1315+int32(_a_F_VP8EncWrite_0)+v1303))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1319
	if v1311 < int32(0) {
		goto L238
	} else {
		goto L240
	}
L240:
	;
	F_Flush(m, v34)
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
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1351
	if int32(126) < v1351 {
		goto L245
	} else {
		goto L246
	}
L244:
	;
	v1345 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1347 = v1342 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1345 + v1347
	v1351 = v1340 - v1347
	goto L243
L245:
	;
	if v1335 == int32(0) {
		goto L248
	} else {
		goto L249
	}
L246:
	;
	v1355 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1356 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1355 << (uint(v1356) % 32)
	v1359 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1359 + v1356
	v1363 = m.G1
	v1367 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1363+int32(_a_F_VP8EncWrite_0)+v1351))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1367
	if v1359 < int32(0) {
		goto L245
	} else {
		goto L247
	}
L247:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L245
L248:
	;
	goto L242
L249:
	;
	if int32(-1) < v1335 {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v1433 = int32(64)
	goto L260
L251:
	;
	v1377 = int32(1)
	v1385 = int32(64)
	goto L252
L252:
	;
	v1389 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1391 = v1389 >> (uint(int32(1)) % 32)
	if v1385&(v1377-v1335<<(uint(v1377)%32)) == int32(0) {
		v1401 = v1391
		goto L254
	} else {
		goto L255
	}
L254:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1401
	if int32(126) < v1401 {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	v1395 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1397 = v1391 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1395 + v1397
	v1401 = v1389 - v1397
	goto L254
L256:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v1385) {
		v1385 = int32(base.Ui32(v1385) >> (uint(int32(1)) % 32))
		goto L252
	} else {
		goto L259
	}
L257:
	;
	v1405 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1406 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1405 << (uint(v1406) % 32)
	v1409 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1409 + v1406
	v1413 = m.G1
	v1417 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1413+int32(_a_F_VP8EncWrite_0)+v1401))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1417
	if v1409 < int32(0) {
		goto L256
	} else {
		goto L258
	}
L258:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L256
L259:
	;
	goto L248
L260:
	;
	v1437 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1439 = v1437 >> (uint(int32(1)) % 32)
	if v1433&(v1335<<(uint(int32(1))%32)) == int32(0) {
		v1449 = v1439
		goto L262
	} else {
		goto L263
	}
L261:
	;
	goto L248
L262:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1449
	if int32(126) < v1449 {
		goto L264
	} else {
		goto L265
	}
L263:
	;
	v1443 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1445 = v1439 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1443 + v1445
	v1449 = v1437 - v1445
	goto L262
L264:
	;
	v1471 = int32(1)
	if base.Ui32(v1471) < base.Ui32(v1433) {
		v1433 = int32(base.Ui32(v1433) >> (uint(v1471) % 32))
		goto L260
	} else {
		goto L267
	}
L265:
	;
	v1453 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1454 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1453 << (uint(v1454) % 32)
	v1457 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1457 + v1454
	v1461 = m.G1
	v1465 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1461+int32(_a_F_VP8EncWrite_0)+v1449))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1465
	if v1457 < int32(0) {
		goto L264
	} else {
		goto L266
	}
L266:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L264
L267:
	;
	goto L261
L268:
	;
	v1484 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v162))))
	v1486 = base.B2i32(v1484 != int32(255))
	v1489 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1491 = v1489 >> (uint(int32(1)) % 32)
	if v1486 == int32(0) {
		v1500 = v1491
		goto L271
	} else {
		goto L272
	}
L269:
	;
	v1577 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)))
	v1579 = base.B2i32(v1577 != int32(255))
	v1582 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1584 = v1582 >> (uint(int32(1)) % 32)
	if v1579 == int32(0) {
		v1593 = v1584
		goto L288
	} else {
		goto L289
	}
L270:
	;
	if v1486 == int32(0) {
		goto L269
	} else {
		goto L276
	}
L271:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1500
	if int32(126) < v1500 {
		goto L273
	} else {
		goto L274
	}
L272:
	;
	v1494 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1496 = v1491 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1494 + v1496
	v1500 = v1489 - v1496
	goto L271
L273:
	;
	goto L270
L274:
	;
	v1504 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1505 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1504 << (uint(v1505) % 32)
	v1508 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1508 + v1505
	v1512 = m.G1
	v1516 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1512+int32(_a_F_VP8EncWrite_0)+v1500))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1516
	if v1508 < int32(0) {
		goto L273
	} else {
		goto L275
	}
L275:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L273
L276:
	;
	v1524 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v162))))
	v1536 = int32(128)
	goto L278
L277:
	;
	goto L269
L278:
	;
	v1539 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1541 = v1539 >> (uint(int32(1)) % 32)
	if v1536&v1524 == int32(0) {
		v1551 = v1541
		goto L280
	} else {
		goto L281
	}
L279:
	;
	goto L277
L280:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1551
	if int32(126) < v1551 {
		goto L282
	} else {
		goto L283
	}
L281:
	;
	v1545 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1547 = v1541 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1545 + v1547
	v1551 = v1539 - v1547
	goto L280
L282:
	;
	v1573 = int32(1)
	if base.Ui32(v1573) < base.Ui32(v1536) {
		v1536 = int32(base.Ui32(v1536) >> (uint(v1573) % 32))
		goto L278
	} else {
		goto L285
	}
L283:
	;
	v1555 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1556 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1555 << (uint(v1556) % 32)
	v1559 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1559 + v1556
	v1563 = m.G1
	v1567 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1563+int32(_a_F_VP8EncWrite_0)+v1551))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1567
	if v1559 < int32(0) {
		goto L282
	} else {
		goto L284
	}
L284:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L282
L285:
	;
	goto L279
L286:
	;
	v1670 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)))
	v1672 = base.B2i32(v1670 != int32(255))
	v1675 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1677 = v1675 >> (uint(int32(1)) % 32)
	if v1672 == int32(0) {
		v1686 = v1677
		goto L304
	} else {
		goto L305
	}
L287:
	;
	if v1579 == int32(0) {
		goto L286
	} else {
		goto L293
	}
L288:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1593
	if int32(126) < v1593 {
		goto L290
	} else {
		goto L291
	}
L289:
	;
	v1587 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1589 = v1584 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1587 + v1589
	v1593 = v1582 - v1589
	goto L288
L290:
	;
	goto L287
L291:
	;
	v1597 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1598 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1597 << (uint(v1598) % 32)
	v1601 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1601 + v1598
	v1605 = m.G1
	v1609 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1605+int32(_a_F_VP8EncWrite_0)+v1593))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1609
	if v1601 < int32(0) {
		goto L290
	} else {
		goto L292
	}
L292:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L290
L293:
	;
	v1617 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3417)))
	v1629 = int32(128)
	goto L295
L294:
	;
	goto L286
L295:
	;
	v1632 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1634 = v1632 >> (uint(int32(1)) % 32)
	if v1629&v1617 == int32(0) {
		v1644 = v1634
		goto L297
	} else {
		goto L298
	}
L296:
	;
	goto L294
L297:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1644
	if int32(126) < v1644 {
		goto L299
	} else {
		goto L300
	}
L298:
	;
	v1638 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1640 = v1634 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1638 + v1640
	v1644 = v1632 - v1640
	goto L297
L299:
	;
	v1666 = int32(1)
	if base.Ui32(v1666) < base.Ui32(v1629) {
		v1629 = int32(base.Ui32(v1629) >> (uint(v1666) % 32))
		goto L295
	} else {
		goto L302
	}
L300:
	;
	v1648 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1649 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1648 << (uint(v1649) % 32)
	v1652 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1652 + v1649
	v1656 = m.G1
	v1660 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1656+int32(_a_F_VP8EncWrite_0)+v1644))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1660
	if v1652 < int32(0) {
		goto L299
	} else {
		goto L301
	}
L301:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L299
L302:
	;
	goto L296
L303:
	;
	if v1672 == int32(0) {
		goto L32
	} else {
		goto L309
	}
L304:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1686
	if int32(126) < v1686 {
		goto L306
	} else {
		goto L307
	}
L305:
	;
	v1680 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1682 = v1677 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1680 + v1682
	v1686 = v1675 - v1682
	goto L304
L306:
	;
	goto L303
L307:
	;
	v1690 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1691 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1690 << (uint(v1691) % 32)
	v1694 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1694 + v1691
	v1698 = m.G1
	v1702 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1698+int32(_a_F_VP8EncWrite_0)+v1686))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1702
	if v1694 < int32(0) {
		goto L306
	} else {
		goto L308
	}
L308:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L306
L309:
	;
	v1710 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+3418)))
	v1722 = int32(128)
	goto L311
L310:
	;
	goto L32
L311:
	;
	v1725 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1727 = v1725 >> (uint(int32(1)) % 32)
	if v1722&v1710 == int32(0) {
		v1737 = v1727
		goto L313
	} else {
		goto L314
	}
L312:
	;
	goto L310
L313:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1737
	if int32(126) < v1737 {
		goto L315
	} else {
		goto L316
	}
L314:
	;
	v1731 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1733 = v1727 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1731 + v1733
	v1737 = v1725 - v1733
	goto L313
L315:
	;
	v1759 = int32(1)
	if base.Ui32(v1759) < base.Ui32(v1722) {
		v1722 = int32(base.Ui32(v1722) >> (uint(v1759) % 32))
		goto L311
	} else {
		goto L318
	}
L316:
	;
	v1741 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1742 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1741 << (uint(v1742) % 32)
	v1745 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1745 + v1742
	v1749 = m.G1
	v1753 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1749+int32(_a_F_VP8EncWrite_0)+v1737))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1753
	if v1745 < int32(0) {
		goto L315
	} else {
		goto L317
	}
L317:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L315
L318:
	;
	goto L312
L319:
	;
	v1800 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v1812 = int32(32)
	goto L326
L320:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1778
	if int32(126) < v1778 {
		goto L322
	} else {
		goto L323
	}
L321:
	;
	v1772 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1774 = v1769 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1772 + v1774
	v1778 = v1767 - v1774
	goto L320
L322:
	;
	goto L319
L323:
	;
	v1782 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1783 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1782 << (uint(v1783) % 32)
	v1786 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1786 + v1783
	v1790 = m.G1
	v1794 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1790+int32(_a_F_VP8EncWrite_0)+v1778))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1794
	if v1786 < int32(0) {
		goto L322
	} else {
		goto L324
	}
L324:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L322
L325:
	;
	v1854 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	v1866 = int32(4)
	goto L335
L326:
	;
	v1815 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1817 = v1815 >> (uint(int32(1)) % 32)
	if v1812&v1800 == int32(0) {
		v1827 = v1817
		goto L328
	} else {
		goto L329
	}
L327:
	;
	goto L325
L328:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1827
	if int32(126) < v1827 {
		goto L330
	} else {
		goto L331
	}
L329:
	;
	v1821 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1823 = v1817 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1821 + v1823
	v1827 = v1815 - v1823
	goto L328
L330:
	;
	v1849 = int32(1)
	if base.Ui32(v1849) < base.Ui32(v1812) {
		v1812 = int32(base.Ui32(v1812) >> (uint(v1849) % 32))
		goto L326
	} else {
		goto L333
	}
L331:
	;
	v1831 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1832 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1831 << (uint(v1832) % 32)
	v1835 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1835 + v1832
	v1839 = m.G1
	v1843 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1839+int32(_a_F_VP8EncWrite_0)+v1827))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1843
	if v1835 < int32(0) {
		goto L330
	} else {
		goto L332
	}
L332:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L330
L333:
	;
	goto L327
L334:
	;
	v1907 = int32(0)
	v1908 = base.B2i32(v1763 != v1907)
	v1911 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1913 = v1911 >> (uint(int32(1)) % 32)
	if v1908 == v1907 {
		v1922 = v1913
		goto L345
	} else {
		goto L346
	}
L335:
	;
	v1869 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1871 = v1869 >> (uint(int32(1)) % 32)
	if v1866&v1854 == int32(0) {
		v1881 = v1871
		goto L337
	} else {
		goto L338
	}
L336:
	;
	goto L334
L337:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1881
	if int32(126) < v1881 {
		goto L339
	} else {
		goto L340
	}
L338:
	;
	v1875 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1877 = v1871 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1875 + v1877
	v1881 = v1869 - v1877
	goto L337
L339:
	;
	v1903 = int32(1)
	if base.Ui32(v1903) < base.Ui32(v1866) {
		v1866 = int32(base.Ui32(v1866) >> (uint(v1903) % 32))
		goto L335
	} else {
		goto L342
	}
L340:
	;
	v1885 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1886 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1885 << (uint(v1886) % 32)
	v1889 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1889 + v1886
	v1893 = m.G1
	v1897 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1893+int32(_a_F_VP8EncWrite_0)+v1881))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1897
	if v1889 < int32(0) {
		goto L339
	} else {
		goto L341
	}
L341:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L339
L342:
	;
	goto L336
L343:
	;
	v2238 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	switch v2238 + int32(-4) {
	case 0:
		goto L403
	default:
		goto L404
	case 4:
		v2244 = int32(3)
		goto L402
	}
L344:
	;
	if v1908 == int32(0) {
		goto L343
	} else {
		goto L350
	}
L345:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1922
	if int32(126) < v1922 {
		goto L347
	} else {
		goto L348
	}
L346:
	;
	v1916 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1918 = v1913 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1916 + v1918
	v1922 = v1911 - v1918
	goto L345
L347:
	;
	goto L344
L348:
	;
	v1926 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1927 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1926 << (uint(v1927) % 32)
	v1930 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1930 + v1927
	v1934 = m.G1
	v1938 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1934+int32(_a_F_VP8EncWrite_0)+v1922))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1938
	if v1930 < int32(0) {
		goto L347
	} else {
		goto L349
	}
L349:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L347
L350:
	;
	v1946 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v1947 = int32(0)
	v1948 = base.B2i32(v1946 != v1947)
	v1951 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v1953 = v1951 >> (uint(int32(1)) % 32)
	if v1948 == v1947 {
		v1962 = v1953
		goto L352
	} else {
		goto L353
	}
L351:
	;
	if v1948 == int32(0) {
		goto L343
	} else {
		goto L357
	}
L352:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1962
	if int32(126) < v1962 {
		goto L354
	} else {
		goto L355
	}
L353:
	;
	v1956 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1958 = v1953 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1956 + v1958
	v1962 = v1951 - v1958
	goto L352
L354:
	;
	goto L351
L355:
	;
	v1966 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v1967 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v1966 << (uint(v1967) % 32)
	v1970 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v1970 + v1967
	v1974 = m.G1
	v1978 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1974+int32(_a_F_VP8EncWrite_0)+v1962))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v1978
	if v1970 < int32(0) {
		goto L354
	} else {
		goto L356
	}
L356:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L354
L357:
	;
	v1998 = int32(8)
	goto L359
L358:
	;
	v2039 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
	v2044 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2046 = v2044 >> (uint(int32(1)) % 32)
	if v2039 == int32(0) {
		v2055 = v2046
		goto L368
	} else {
		goto L369
	}
L359:
	;
	v2001 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2003 = v2001 >> (uint(int32(1)) % 32)
	if v1998&int32(0) == int32(0) {
		v2013 = v2003
		goto L361
	} else {
		goto L362
	}
L360:
	;
	goto L358
L361:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2013
	if int32(126) < v2013 {
		goto L363
	} else {
		goto L364
	}
L362:
	;
	v2007 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2009 = v2003 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2007 + v2009
	v2013 = v2001 - v2009
	goto L361
L363:
	;
	v2035 = int32(1)
	if base.Ui32(v2035) < base.Ui32(v1998) {
		v1998 = int32(base.Ui32(v1998) >> (uint(v2035) % 32))
		goto L359
	} else {
		goto L366
	}
L364:
	;
	v2017 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2018 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2017 << (uint(v2018) % 32)
	v2021 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2021 + v2018
	v2025 = m.G1
	v2029 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2025+int32(_a_F_VP8EncWrite_0)+v2013))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2029
	if v2021 < int32(0) {
		goto L363
	} else {
		goto L365
	}
L365:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L363
L366:
	;
	goto L360
L367:
	;
	v2197 = int32(4)
	goto L394
L368:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2055
	if int32(126) < v2055 {
		goto L370
	} else {
		goto L371
	}
L369:
	;
	v2049 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2051 = v2046 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2049 + v2051
	v2055 = v2044 - v2051
	goto L368
L370:
	;
	if v2039 == int32(0) {
		goto L373
	} else {
		goto L374
	}
L371:
	;
	v2059 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2060 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2059 << (uint(v2060) % 32)
	v2063 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2063 + v2060
	v2067 = m.G1
	v2071 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2067+int32(_a_F_VP8EncWrite_0)+v2055))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2071
	if v2063 < int32(0) {
		goto L370
	} else {
		goto L372
	}
L372:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L370
L373:
	;
	goto L367
L374:
	;
	if int32(-1) < v2039 {
		goto L375
	} else {
		goto L376
	}
L375:
	;
	v2137 = int32(64)
	goto L385
L376:
	;
	v2081 = int32(1)
	v2089 = int32(64)
	goto L377
L377:
	;
	v2093 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2095 = v2093 >> (uint(int32(1)) % 32)
	if v2089&(v2081-v2039<<(uint(v2081)%32)) == int32(0) {
		v2105 = v2095
		goto L379
	} else {
		goto L380
	}
L379:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2105
	if int32(126) < v2105 {
		goto L381
	} else {
		goto L382
	}
L380:
	;
	v2099 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2101 = v2095 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2099 + v2101
	v2105 = v2093 - v2101
	goto L379
L381:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2089) {
		v2089 = int32(base.Ui32(v2089) >> (uint(int32(1)) % 32))
		goto L377
	} else {
		goto L384
	}
L382:
	;
	v2109 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2110 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2109 << (uint(v2110) % 32)
	v2113 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2113 + v2110
	v2117 = m.G1
	v2121 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2117+int32(_a_F_VP8EncWrite_0)+v2105))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2121
	if v2113 < int32(0) {
		goto L381
	} else {
		goto L383
	}
L383:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L381
L384:
	;
	goto L373
L385:
	;
	v2141 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2143 = v2141 >> (uint(int32(1)) % 32)
	if v2137&(v2039<<(uint(int32(1))%32)) == int32(0) {
		v2153 = v2143
		goto L387
	} else {
		goto L388
	}
L386:
	;
	goto L373
L387:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2153
	if int32(126) < v2153 {
		goto L389
	} else {
		goto L390
	}
L388:
	;
	v2147 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2149 = v2143 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2147 + v2149
	v2153 = v2141 - v2149
	goto L387
L389:
	;
	v2175 = int32(1)
	if base.Ui32(v2175) < base.Ui32(v2137) {
		v2137 = int32(base.Ui32(v2137) >> (uint(v2175) % 32))
		goto L385
	} else {
		goto L392
	}
L390:
	;
	v2157 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2158 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2157 << (uint(v2158) % 32)
	v2161 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2161 + v2158
	v2165 = m.G1
	v2169 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2165+int32(_a_F_VP8EncWrite_0)+v2153))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2169
	if v2161 < int32(0) {
		goto L389
	} else {
		goto L391
	}
L391:
	;
	F_Flush(m, v34)
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
	v2200 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2202 = v2200 >> (uint(int32(1)) % 32)
	if v2197&int32(0) == int32(0) {
		v2212 = v2202
		goto L396
	} else {
		goto L397
	}
L395:
	;
	goto L393
L396:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2212
	if int32(126) < v2212 {
		goto L398
	} else {
		goto L399
	}
L397:
	;
	v2206 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2208 = v2202 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2206 + v2208
	v2212 = v2200 - v2208
	goto L396
L398:
	;
	v2234 = int32(1)
	if base.Ui32(v2234) < base.Ui32(v2197) {
		v2197 = int32(base.Ui32(v2197) >> (uint(v2234) % 32))
		goto L394
	} else {
		goto L401
	}
L399:
	;
	v2216 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2217 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2216 << (uint(v2217) % 32)
	v2220 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2220 + v2217
	v2224 = m.G1
	v2228 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2224+int32(_a_F_VP8EncWrite_0)+v2212))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2228
	if v2220 < int32(0) {
		goto L398
	} else {
		goto L400
	}
L400:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L398
L401:
	;
	goto L395
L402:
	;
	v2256 = int32(2)
	goto L406
L403:
	;
	v2244 = int32(2)
	goto L402
L404:
	;
	v2244 = base.B2i32(v2238 == int32(2))
	goto L402
L405:
	;
	v2297 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3384))
	v2309 = int32(64)
	goto L415
L406:
	;
	v2259 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2261 = v2259 >> (uint(int32(1)) % 32)
	if v2256&v2244 == int32(0) {
		v2271 = v2261
		goto L408
	} else {
		goto L409
	}
L407:
	;
	goto L405
L408:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2271
	if int32(126) < v2271 {
		goto L410
	} else {
		goto L411
	}
L409:
	;
	v2265 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2267 = v2261 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2265 + v2267
	v2271 = v2259 - v2267
	goto L408
L410:
	;
	v2293 = int32(1)
	if base.Ui32(v2293) < base.Ui32(v2256) {
		v2256 = int32(base.Ui32(v2256) >> (uint(v2293) % 32))
		goto L406
	} else {
		goto L413
	}
L411:
	;
	v2275 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2276 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2275 << (uint(v2276) % 32)
	v2279 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2279 + v2276
	v2283 = m.G1
	v2287 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2283+int32(_a_F_VP8EncWrite_0)+v2271))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2287
	if v2279 < int32(0) {
		goto L410
	} else {
		goto L412
	}
L412:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L410
L413:
	;
	goto L407
L414:
	;
	v2350 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3396))
	v2355 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2357 = v2355 >> (uint(int32(1)) % 32)
	if v2350 == int32(0) {
		v2366 = v2357
		goto L424
	} else {
		goto L425
	}
L415:
	;
	v2312 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2314 = v2312 >> (uint(int32(1)) % 32)
	if v2309&v2297 == int32(0) {
		v2324 = v2314
		goto L417
	} else {
		goto L418
	}
L416:
	;
	goto L414
L417:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2324
	if int32(126) < v2324 {
		goto L419
	} else {
		goto L420
	}
L418:
	;
	v2318 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2320 = v2314 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2318 + v2320
	v2324 = v2312 - v2320
	goto L417
L419:
	;
	v2346 = int32(1)
	if base.Ui32(v2346) < base.Ui32(v2309) {
		v2309 = int32(base.Ui32(v2309) >> (uint(v2346) % 32))
		goto L415
	} else {
		goto L422
	}
L420:
	;
	v2328 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2329 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2328 << (uint(v2329) % 32)
	v2332 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2332 + v2329
	v2336 = m.G1
	v2340 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2336+int32(_a_F_VP8EncWrite_0)+v2324))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2340
	if v2332 < int32(0) {
		goto L419
	} else {
		goto L421
	}
L421:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L419
L422:
	;
	goto L416
L423:
	;
	v2496 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3400))
	v2501 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2503 = v2501 >> (uint(int32(1)) % 32)
	if v2496 == int32(0) {
		v2512 = v2503
		goto L450
	} else {
		goto L451
	}
L424:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2366
	if int32(126) < v2366 {
		goto L426
	} else {
		goto L427
	}
L425:
	;
	v2360 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2362 = v2357 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2360 + v2362
	v2366 = v2355 - v2362
	goto L424
L426:
	;
	if v2350 == int32(0) {
		goto L429
	} else {
		goto L430
	}
L427:
	;
	v2370 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2371 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2370 << (uint(v2371) % 32)
	v2374 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2374 + v2371
	v2378 = m.G1
	v2382 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2378+int32(_a_F_VP8EncWrite_0)+v2366))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2382
	if v2374 < int32(0) {
		goto L426
	} else {
		goto L428
	}
L428:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L426
L429:
	;
	goto L423
L430:
	;
	if int32(-1) < v2350 {
		goto L431
	} else {
		goto L432
	}
L431:
	;
	v2448 = int32(16)
	goto L441
L432:
	;
	v2392 = int32(1)
	v2400 = int32(16)
	goto L433
L433:
	;
	v2404 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2406 = v2404 >> (uint(int32(1)) % 32)
	if v2400&(v2392-v2350<<(uint(v2392)%32)) == int32(0) {
		v2416 = v2406
		goto L435
	} else {
		goto L436
	}
L435:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2416
	if int32(126) < v2416 {
		goto L437
	} else {
		goto L438
	}
L436:
	;
	v2410 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2412 = v2406 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2410 + v2412
	v2416 = v2404 - v2412
	goto L435
L437:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2400) {
		v2400 = int32(base.Ui32(v2400) >> (uint(int32(1)) % 32))
		goto L433
	} else {
		goto L440
	}
L438:
	;
	v2420 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2421 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2420 << (uint(v2421) % 32)
	v2424 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2424 + v2421
	v2428 = m.G1
	v2432 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2428+int32(_a_F_VP8EncWrite_0)+v2416))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2432
	if v2424 < int32(0) {
		goto L437
	} else {
		goto L439
	}
L439:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L437
L440:
	;
	goto L429
L441:
	;
	v2452 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2454 = v2452 >> (uint(int32(1)) % 32)
	if v2448&(v2350<<(uint(int32(1))%32)) == int32(0) {
		v2464 = v2454
		goto L443
	} else {
		goto L444
	}
L442:
	;
	goto L429
L443:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2464
	if int32(126) < v2464 {
		goto L445
	} else {
		goto L446
	}
L444:
	;
	v2458 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2460 = v2454 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2458 + v2460
	v2464 = v2452 - v2460
	goto L443
L445:
	;
	v2486 = int32(1)
	if base.Ui32(v2486) < base.Ui32(v2448) {
		v2448 = int32(base.Ui32(v2448) >> (uint(v2486) % 32))
		goto L441
	} else {
		goto L448
	}
L446:
	;
	v2468 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2469 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2468 << (uint(v2469) % 32)
	v2472 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2472 + v2469
	v2476 = m.G1
	v2480 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2476+int32(_a_F_VP8EncWrite_0)+v2464))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2480
	if v2472 < int32(0) {
		goto L445
	} else {
		goto L447
	}
L447:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L445
L448:
	;
	goto L442
L449:
	;
	v2642 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3404))
	v2647 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2649 = v2647 >> (uint(int32(1)) % 32)
	if v2642 == int32(0) {
		v2658 = v2649
		goto L476
	} else {
		goto L477
	}
L450:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2512
	if int32(126) < v2512 {
		goto L452
	} else {
		goto L453
	}
L451:
	;
	v2506 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2508 = v2503 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2506 + v2508
	v2512 = v2501 - v2508
	goto L450
L452:
	;
	if v2496 == int32(0) {
		goto L455
	} else {
		goto L456
	}
L453:
	;
	v2516 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2517 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2516 << (uint(v2517) % 32)
	v2520 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2520 + v2517
	v2524 = m.G1
	v2528 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2524+int32(_a_F_VP8EncWrite_0)+v2512))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2528
	if v2520 < int32(0) {
		goto L452
	} else {
		goto L454
	}
L454:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L452
L455:
	;
	goto L449
L456:
	;
	if int32(-1) < v2496 {
		goto L457
	} else {
		goto L458
	}
L457:
	;
	v2594 = int32(16)
	goto L467
L458:
	;
	v2538 = int32(1)
	v2546 = int32(16)
	goto L459
L459:
	;
	v2550 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2552 = v2550 >> (uint(int32(1)) % 32)
	if v2546&(v2538-v2496<<(uint(v2538)%32)) == int32(0) {
		v2562 = v2552
		goto L461
	} else {
		goto L462
	}
L461:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2562
	if int32(126) < v2562 {
		goto L463
	} else {
		goto L464
	}
L462:
	;
	v2556 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2558 = v2552 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2556 + v2558
	v2562 = v2550 - v2558
	goto L461
L463:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2546) {
		v2546 = int32(base.Ui32(v2546) >> (uint(int32(1)) % 32))
		goto L459
	} else {
		goto L466
	}
L464:
	;
	v2566 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2567 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2566 << (uint(v2567) % 32)
	v2570 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2570 + v2567
	v2574 = m.G1
	v2578 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2574+int32(_a_F_VP8EncWrite_0)+v2562))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2578
	if v2570 < int32(0) {
		goto L463
	} else {
		goto L465
	}
L465:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L463
L466:
	;
	goto L455
L467:
	;
	v2598 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2600 = v2598 >> (uint(int32(1)) % 32)
	if v2594&(v2496<<(uint(int32(1))%32)) == int32(0) {
		v2610 = v2600
		goto L469
	} else {
		goto L470
	}
L468:
	;
	goto L455
L469:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2610
	if int32(126) < v2610 {
		goto L471
	} else {
		goto L472
	}
L470:
	;
	v2604 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2606 = v2600 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2604 + v2606
	v2610 = v2598 - v2606
	goto L469
L471:
	;
	v2632 = int32(1)
	if base.Ui32(v2632) < base.Ui32(v2594) {
		v2594 = int32(base.Ui32(v2594) >> (uint(v2632) % 32))
		goto L467
	} else {
		goto L474
	}
L472:
	;
	v2614 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2615 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2614 << (uint(v2615) % 32)
	v2618 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2618 + v2615
	v2622 = m.G1
	v2626 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2622+int32(_a_F_VP8EncWrite_0)+v2610))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2626
	if v2618 < int32(0) {
		goto L471
	} else {
		goto L473
	}
L473:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L471
L474:
	;
	goto L468
L475:
	;
	v2788 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3408))
	v2793 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2795 = v2793 >> (uint(int32(1)) % 32)
	if v2788 == int32(0) {
		v2804 = v2795
		goto L502
	} else {
		goto L503
	}
L476:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2658
	if int32(126) < v2658 {
		goto L478
	} else {
		goto L479
	}
L477:
	;
	v2652 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2654 = v2649 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2652 + v2654
	v2658 = v2647 - v2654
	goto L476
L478:
	;
	if v2642 == int32(0) {
		goto L481
	} else {
		goto L482
	}
L479:
	;
	v2662 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2663 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2662 << (uint(v2663) % 32)
	v2666 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2666 + v2663
	v2670 = m.G1
	v2674 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2670+int32(_a_F_VP8EncWrite_0)+v2658))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2674
	if v2666 < int32(0) {
		goto L478
	} else {
		goto L480
	}
L480:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L478
L481:
	;
	goto L475
L482:
	;
	if int32(-1) < v2642 {
		goto L483
	} else {
		goto L484
	}
L483:
	;
	v2740 = int32(16)
	goto L493
L484:
	;
	v2684 = int32(1)
	v2692 = int32(16)
	goto L485
L485:
	;
	v2696 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2698 = v2696 >> (uint(int32(1)) % 32)
	if v2692&(v2684-v2642<<(uint(v2684)%32)) == int32(0) {
		v2708 = v2698
		goto L487
	} else {
		goto L488
	}
L487:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2708
	if int32(126) < v2708 {
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v2702 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2704 = v2698 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2702 + v2704
	v2708 = v2696 - v2704
	goto L487
L489:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2692) {
		v2692 = int32(base.Ui32(v2692) >> (uint(int32(1)) % 32))
		goto L485
	} else {
		goto L492
	}
L490:
	;
	v2712 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2713 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2712 << (uint(v2713) % 32)
	v2716 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2716 + v2713
	v2720 = m.G1
	v2724 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2720+int32(_a_F_VP8EncWrite_0)+v2708))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2724
	if v2716 < int32(0) {
		goto L489
	} else {
		goto L491
	}
L491:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L489
L492:
	;
	goto L481
L493:
	;
	v2744 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2746 = v2744 >> (uint(int32(1)) % 32)
	if v2740&(v2642<<(uint(int32(1))%32)) == int32(0) {
		v2756 = v2746
		goto L495
	} else {
		goto L496
	}
L494:
	;
	goto L481
L495:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2756
	if int32(126) < v2756 {
		goto L497
	} else {
		goto L498
	}
L496:
	;
	v2750 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2752 = v2746 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2750 + v2752
	v2756 = v2744 - v2752
	goto L495
L497:
	;
	v2778 = int32(1)
	if base.Ui32(v2778) < base.Ui32(v2740) {
		v2740 = int32(base.Ui32(v2740) >> (uint(v2778) % 32))
		goto L493
	} else {
		goto L500
	}
L498:
	;
	v2760 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2761 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2760 << (uint(v2761) % 32)
	v2764 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2764 + v2761
	v2768 = m.G1
	v2772 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2768+int32(_a_F_VP8EncWrite_0)+v2756))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2772
	if v2764 < int32(0) {
		goto L497
	} else {
		goto L499
	}
L499:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L497
L500:
	;
	goto L494
L501:
	;
	v2934 = *(*int32)(unsafe.Add(mBase, uint32(l0)+3412))
	v2939 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2941 = v2939 >> (uint(int32(1)) % 32)
	if v2934 == int32(0) {
		v2950 = v2941
		goto L528
	} else {
		goto L529
	}
L502:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2804
	if int32(126) < v2804 {
		goto L504
	} else {
		goto L505
	}
L503:
	;
	v2798 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2800 = v2795 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2798 + v2800
	v2804 = v2793 - v2800
	goto L502
L504:
	;
	if v2788 == int32(0) {
		goto L507
	} else {
		goto L508
	}
L505:
	;
	v2808 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2809 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2808 << (uint(v2809) % 32)
	v2812 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2812 + v2809
	v2816 = m.G1
	v2820 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2816+int32(_a_F_VP8EncWrite_0)+v2804))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2820
	if v2812 < int32(0) {
		goto L504
	} else {
		goto L506
	}
L506:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L504
L507:
	;
	goto L501
L508:
	;
	if int32(-1) < v2788 {
		goto L509
	} else {
		goto L510
	}
L509:
	;
	v2886 = int32(16)
	goto L519
L510:
	;
	v2830 = int32(1)
	v2838 = int32(16)
	goto L511
L511:
	;
	v2842 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2844 = v2842 >> (uint(int32(1)) % 32)
	if v2838&(v2830-v2788<<(uint(v2830)%32)) == int32(0) {
		v2854 = v2844
		goto L513
	} else {
		goto L514
	}
L513:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2854
	if int32(126) < v2854 {
		goto L515
	} else {
		goto L516
	}
L514:
	;
	v2848 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2850 = v2844 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2848 + v2850
	v2854 = v2842 - v2850
	goto L513
L515:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2838) {
		v2838 = int32(base.Ui32(v2838) >> (uint(int32(1)) % 32))
		goto L511
	} else {
		goto L518
	}
L516:
	;
	v2858 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2859 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2858 << (uint(v2859) % 32)
	v2862 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2862 + v2859
	v2866 = m.G1
	v2870 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2866+int32(_a_F_VP8EncWrite_0)+v2854))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2870
	if v2862 < int32(0) {
		goto L515
	} else {
		goto L517
	}
L517:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L515
L518:
	;
	goto L507
L519:
	;
	v2890 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2892 = v2890 >> (uint(int32(1)) % 32)
	if v2886&(v2788<<(uint(int32(1))%32)) == int32(0) {
		v2902 = v2892
		goto L521
	} else {
		goto L522
	}
L520:
	;
	goto L507
L521:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2902
	if int32(126) < v2902 {
		goto L523
	} else {
		goto L524
	}
L522:
	;
	v2896 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2898 = v2892 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2896 + v2898
	v2902 = v2890 - v2898
	goto L521
L523:
	;
	v2924 = int32(1)
	if base.Ui32(v2924) < base.Ui32(v2886) {
		v2886 = int32(base.Ui32(v2886) >> (uint(v2924) % 32))
		goto L519
	} else {
		goto L526
	}
L524:
	;
	v2906 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2907 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2906 << (uint(v2907) % 32)
	v2910 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2910 + v2907
	v2914 = m.G1
	v2918 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2914+int32(_a_F_VP8EncWrite_0)+v2902))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2918
	if v2910 < int32(0) {
		goto L523
	} else {
		goto L525
	}
L525:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L523
L526:
	;
	goto L520
L527:
	;
	v3083 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v3085 = v3083 >> (uint(int32(1)) % 32)
	goto L554
L528:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2950
	if int32(126) < v2950 {
		goto L530
	} else {
		goto L531
	}
L529:
	;
	v2944 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2946 = v2941 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2944 + v2946
	v2950 = v2939 - v2946
	goto L528
L530:
	;
	if v2934 == int32(0) {
		goto L533
	} else {
		goto L534
	}
L531:
	;
	v2954 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2955 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2954 << (uint(v2955) % 32)
	v2958 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v2958 + v2955
	v2962 = m.G1
	v2966 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2962+int32(_a_F_VP8EncWrite_0)+v2950))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v2966
	if v2958 < int32(0) {
		goto L530
	} else {
		goto L532
	}
L532:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L530
L533:
	;
	goto L527
L534:
	;
	if int32(-1) < v2934 {
		goto L535
	} else {
		goto L536
	}
L535:
	;
	v3032 = int32(16)
	goto L545
L536:
	;
	v2976 = int32(1)
	v2984 = int32(16)
	goto L537
L537:
	;
	v2988 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v2990 = v2988 >> (uint(int32(1)) % 32)
	if v2984&(v2976-v2934<<(uint(v2976)%32)) == int32(0) {
		v3000 = v2990
		goto L539
	} else {
		goto L540
	}
L539:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3000
	if int32(126) < v3000 {
		goto L541
	} else {
		goto L542
	}
L540:
	;
	v2994 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v2996 = v2990 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v2994 + v2996
	v3000 = v2988 - v2996
	goto L539
L541:
	;
	if base.Ui32(int32(2)) <= base.Ui32(v2984) {
		v2984 = int32(base.Ui32(v2984) >> (uint(int32(1)) % 32))
		goto L537
	} else {
		goto L544
	}
L542:
	;
	v3004 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v3005 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v3004 << (uint(v3005) % 32)
	v3008 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v3008 + v3005
	v3012 = m.G1
	v3016 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3012+int32(_a_F_VP8EncWrite_0)+v3000))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3016
	if v3008 < int32(0) {
		goto L541
	} else {
		goto L543
	}
L543:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L541
L544:
	;
	goto L533
L545:
	;
	v3036 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v3038 = v3036 >> (uint(int32(1)) % 32)
	if v3032&(v2934<<(uint(int32(1))%32)) == int32(0) {
		v3048 = v3038
		goto L547
	} else {
		goto L548
	}
L546:
	;
	goto L533
L547:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3048
	if int32(126) < v3048 {
		goto L549
	} else {
		goto L550
	}
L548:
	;
	v3042 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v3044 = v3038 + int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v3042 + v3044
	v3048 = v3036 - v3044
	goto L547
L549:
	;
	v3070 = int32(1)
	if base.Ui32(v3070) < base.Ui32(v3032) {
		v3032 = int32(base.Ui32(v3032) >> (uint(v3070) % 32))
		goto L545
	} else {
		goto L552
	}
L550:
	;
	v3052 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v3053 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v3052 << (uint(v3053) % 32)
	v3056 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v3056 + v3053
	v3060 = m.G1
	v3064 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3060+int32(_a_F_VP8EncWrite_0)+v3048))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3064
	if v3056 < int32(0) {
		goto L549
	} else {
		goto L551
	}
L551:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L549
L552:
	;
	goto L546
L553:
	;
	F_VP8WriteProbas(m, v34, v162)
	mBase = m.M
	v3118 = l0 + int32(76)
	v3119 = *(*int32)(unsafe.Add(mBase, uint32(v3118)))
	v3120 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v3121 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	F_VP8CodeIntraModes(m, l0)
	mBase = m.M
	v3128 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	v3132 = int32(1) << (uint(int32(8)-v3128) % 32)
	goto L560
L554:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3085
	if int32(126) < v3085 {
		goto L556
	} else {
		goto L557
	}
L556:
	;
	goto L553
L557:
	;
	v3098 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v3099 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v3098 << (uint(v3099) % 32)
	v3102 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v3102 + v3099
	v3106 = m.G1
	v3110 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3106+int32(_a_F_VP8EncWrite_0)+v3085))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3110
	if v3102 < int32(0) {
		goto L556
	} else {
		goto L558
	}
L558:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L556
L559:
	;
	v3167 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v3168 = *(*int32)(unsafe.Add(mBase, uint32(v3167)+88))
	if v3168 == int32(0) {
		goto L566
	} else {
		goto L567
	}
L560:
	;
	v3135 = *(*int32)(unsafe.Add(mBase, uint32(v34)))
	v3137 = v3135 >> (uint(int32(1)) % 32)
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3137
	if int32(126) < v3137 {
		goto L562
	} else {
		goto L563
	}
L561:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = int32(0)
	F_Flush(m, v34)
	mBase = m.M
	goto L559
L562:
	;
	v3159 = int32(1)
	if base.Ui32(v3159) < base.Ui32(v3132) {
		v3132 = int32(base.Ui32(v3132) >> (uint(v3159) % 32))
		goto L560
	} else {
		goto L565
	}
L563:
	;
	v3141 = *(*int32)(unsafe.Add(mBase, uint32(v34)+4))
	v3142 = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v34)+4)) = v3141 << (uint(v3142) % 32)
	v3145 = *(*int32)(unsafe.Add(mBase, uint32(v34)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v34)+12)) = v3145 + v3142
	v3149 = m.G1
	v3153 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3149+int32(_a_F_VP8EncWrite_0)+v3137))))
	*(*int32)(unsafe.Add(mBase, uint32(v34))) = v3153
	if v3145 < int32(0) {
		goto L562
	} else {
		goto L564
	}
L564:
	;
	F_Flush(m, v34)
	mBase = m.M
	goto L562
L565:
	;
	goto L561
L566:
	;
	v3215 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	if v3215 == int32(0) {
		goto L3
	} else {
		goto L568
	}
L567:
	;
	v3171 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	v3172 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	v3177 = int32(3)
	v3185 = int32(8)
	v3188 = base.I64_extend_i32_u((v3120+v3119)<<(uint(v3177)%32)) + base.I64_extend_i32_s(v3121+v3185)
	v3190 = int64(7)
	v3192 = int64(3)
	v3193 = int64(base.Ui64(base.I64_extend_i32_s(int32(-8)-v31)-base.I64_extend_i32_u((v30+v29)<<(uint(v3177)%32))+v3188+v3190) >> (uint(v3192) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v3168)+36)) = uint32(v3193)
	v3199 = *(*int32)(unsafe.Add(mBase, uint32(v3118)))
	v3208 = int64(base.Ui64(base.I64_extend_i32_s(v3172+v3185)-v3188+base.I64_extend_i32_u((v3171+v3199)<<(uint(v3177)%32))+v3190) >> (uint(v3192) % 64))
	*(*uint32)(unsafe.Add(mBase, uint32(v3168)+40)) = uint32(v3208)
	v3210 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	*(*int32)(unsafe.Add(mBase, uint32(v3168)+140)) = v3210
	goto L566
L568:
	;
	v3219 = *(*int32)(unsafe.Add(mBase, uint32(v3167)+92))
	if v3219 != 0 {
		goto L570
	} else {
		goto L571
	}
L569:
	;
	v3222 = int32(0)
	goto L4
L570:
	;
	goto L569
L571:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3167)+92)) = int32(1)
	goto L570
L572:
	;
	goto L3
L573:
	;
	v3371 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(372))))
	if v3371 != 0 {
		goto L585
	} else {
		goto L586
	}
L574:
	;
	v3254 = v3245 & int32(3)
	if base.Ui32(int32(4)) <= base.Ui32(v3245) {
		goto L576
	} else {
		goto L577
	}
L575:
	;
	if v3254 == int32(0) {
		v3360 = v3311
		goto L573
	} else {
		goto L581
	}
L576:
	;
	v3268 = int32(0)
	v3272 = l0 + int32(204)
	v3274 = v3250
	goto L578
L577:
	;
	v3305 = int32(0)
	v3311 = v3250
	goto L575
L578:
	;
	v3281 = *(*int32)(unsafe.Add(mBase, uint32(v3272)))
	v3284 = *(*int32)(unsafe.Add(mBase, uint32(v3272+int32(-32))))
	v3287 = *(*int32)(unsafe.Add(mBase, uint32(v3272+int32(-64))))
	v3290 = *(*int32)(unsafe.Add(mBase, uint32(v3272+int32(-96))))
	v3294 = v3281 + (v3284 + (v3287 + (v3290 + v3274)))
	v3298 = v3268 + int32(4)
	if v3245&int32(2147483644) != v3298 {
		v3268 = v3298
		v3272 = v3272 + int32(128)
		v3274 = v3294
		goto L578
	} else {
		goto L580
	}
L579:
	;
	v3305 = v3298
	v3311 = v3294
	goto L575
L580:
	;
	goto L579
L581:
	;
	v3329 = v3254
	v3334 = v3305<<(uint(int32(5))%32) + l0 + int32(108)
	v3336 = v3311
	goto L582
L582:
	;
	v3343 = *(*int32)(unsafe.Add(mBase, uint32(v3334)))
	v3344 = v3343 + v3336
	v3348 = v3329 + int32(-1)
	if v3348 != 0 {
		v3329 = v3348
		v3334 = v3334 + int32(32)
		v3336 = v3344
		goto L582
	} else {
		goto L584
	}
L583:
	;
	v3360 = v3344
	goto L573
L584:
	;
	goto L583
L585:
	;
	v3372 = int32(30)
	goto L587
L586:
	;
	v3372 = int32(12)
	goto L587
L587:
	;
	v3374 = v3360 & int32(1)
	v3375 = v3374 + v3360
	v3376 = v3372 + v3375
	if v3371 == int32(0) {
		v3386 = v3376
		goto L588
	} else {
		goto L589
	}
L588:
	;
	if v3386 != int32(-1) {
		goto L590
	} else {
		goto L591
	}
L589:
	;
	v3379 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	v3386 = v3376 + v3379 + v3379&int32(1) + int32(8)
	goto L588
L590:
	;
	v3396 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(72))))
	v3397 = m.G1
	v3399 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v3408 = *(*int32)(unsafe.Add(mBase, uint32(v3397)+uint32(_c_F_VP8EncWrite[0])))
	*(*int32)(unsafe.Add(mBase, uint32(v21+int32(24)))) = v3408
	v3410 = *(*int64)(unsafe.Add(mBase, uint32(v3397)+uint32(_c_F_VP8EncWrite[1])))
	*(*int64)(unsafe.Add(mBase, uint32(v21)+16)) = v3410
	*(*int32)(unsafe.Add(mBase, uint32(v21)+20)) = v3386
	v3416 = *(*int32)(unsafe.Add(mBase, uint32(v3399)+72))
	v3417 = m.T0[v3416].(func(*base.Module, int32, int32, int32) int32)(m, v21+int32(16), int32(12), v3399)
	mBase = m.M
	if v3417 == int32(0) {
		v3562 = int32(8)
		goto L597
	} else {
		goto L598
	}
L591:
	;
	v3391 = *(*int32)(unsafe.Add(mBase, uint32(v32)+92))
	if v3391 != 0 {
		goto L593
	} else {
		goto L594
	}
L592:
	;
	v3837 = int32(0)
	goto L1
L593:
	;
	goto L592
L594:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v32)+92)) = int32(9)
	goto L593
L595:
	;
	v3679 = l0 + int32(368)
	if v34 == int32(0) {
		goto L633
	} else {
		goto L634
	}
L596:
	;
	v3572 = *(*int32)(unsafe.Add(mBase, uint32(v32)+72))
	v3573 = m.T0[v3572].(func(*base.Module, int32, int32, int32) int32)(m, v3396, v3244, v32)
	mBase = m.M
	if v3573 != 0 {
		goto L614
	} else {
		goto L615
	}
L597:
	;
	v3564 = *(*int32)(unsafe.Add(mBase, uint32(v3399)+92))
	if v3564 != 0 {
		goto L611
	} else {
		goto L612
	}
L598:
	;
	v3422 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(372))))
	if v3422 == int32(0) {
		goto L599
	} else {
		goto L600
	}
L599:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v21)+16)) = int64(540561494)
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+20)) = uint8(v3375)
	v3512 = int32(base.Ui32(v3375) >> (uint(int32(24)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+23)) = uint8(v3512)
	v3514 = int32(16)
	v3515 = int32(base.Ui32(v3375) >> (uint(v3514) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+22)) = uint8(v3515)
	v3517 = int32(8)
	v3519 = int32(base.Ui32(v3375) >> (uint(v3517) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+21)) = uint8(v3519)
	v3524 = *(*int32)(unsafe.Add(mBase, uint32(v3399)+72))
	v3525 = m.T0[v3524].(func(*base.Module, int32, int32, int32) int32)(m, v21+v3514, v3517, v3399)
	mBase = m.M
	if v3525 == int32(0) {
		v3562 = v3517
		goto L597
	} else {
		goto L607
	}
L600:
	;
	v3425 = m.G1
	v3426 = int32(8)
	v3427 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v3432 = *(*int64)(unsafe.Add(mBase, uint32(v3425)+uint32(_c_F_VP8EncWrite[2])))
	*(*int64)(unsafe.Add(mBase, uint32(v21)+24)) = v3432
	v3434 = *(*int64)(unsafe.Add(mBase, uint32(v3425)+uint32(_c_F_VP8EncWrite[3])))
	*(*int64)(unsafe.Add(mBase, uint32(v21)+16)) = v3434
	v3436 = int32(16)
	*(*int32)(unsafe.Add(mBase, uint32(v21)+24)) = v3436
	v3438 = *(*int32)(unsafe.Add(mBase, uint32(v3427)+8))
	v3439 = int32(-1)
	v3440 = v3438 + v3439
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+28)) = uint8(v3440)
	v3443 = int32(base.Ui32(v3440) >> (uint(v3436) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+30)) = uint8(v3443)
	v3446 = int32(base.Ui32(v3440) >> (uint(v3426) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+29)) = uint8(v3446)
	*(*int32)(unsafe.Add(mBase, uint32(v21)+20)) = int32(10)
	v3450 = *(*int32)(unsafe.Add(mBase, uint32(v3427)+12))
	v3452 = v3450 + v3439
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+31)) = uint8(v3452)
	v3455 = int32(base.Ui32(v3452) >> (uint(v3436) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+33)) = uint8(v3455)
	v3458 = int32(base.Ui32(v3452) >> (uint(v3426) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+32)) = uint8(v3458)
	v3463 = *(*int32)(unsafe.Add(mBase, uint32(v3427)+72))
	v3464 = m.T0[v3463].(func(*base.Module, int32, int32, int32) int32)(m, v21+v3436, int32(18), v3427)
	mBase = m.M
	if v3464 == int32(0) {
		v3562 = v3426
		goto L597
	} else {
		goto L601
	}
L601:
	;
	v3469 = *(*int32)(unsafe.Add(mBase, uint32(l0+int32(372))))
	if v3469 == int32(0) {
		goto L599
	} else {
		goto L602
	}
L602:
	;
	v3472 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	*(*int64)(unsafe.Add(mBase, uint32(v21)+16)) = int64(1213221953)
	v3475 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	*(*int32)(unsafe.Add(mBase, uint32(v21)+20)) = v3475
	v3477 = int32(8)
	v3481 = *(*int32)(unsafe.Add(mBase, uint32(v3472)+72))
	v3482 = m.T0[v3481].(func(*base.Module, int32, int32, int32) int32)(m, v21+int32(16), v3477, v3472)
	mBase = m.M
	if v3482 == int32(0) {
		v3562 = v3477
		goto L597
	} else {
		goto L603
	}
L603:
	;
	v3485 = *(*int32)(unsafe.Add(mBase, uint32(l0)+376))
	v3486 = *(*int32)(unsafe.Add(mBase, uint32(l0)+380))
	v3487 = *(*int32)(unsafe.Add(mBase, uint32(v3472)+72))
	v3488 = m.T0[v3487].(func(*base.Module, int32, int32, int32) int32)(m, v3485, v3486, v3472)
	mBase = m.M
	if v3488 == int32(0) {
		v3562 = v3477
		goto L597
	} else {
		goto L604
	}
L604:
	;
	v3491 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+380)))
	if v3491&int32(1) == int32(0) {
		goto L599
	} else {
		goto L605
	}
L605:
	;
	v3496 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+15)) = uint8(v3496)
	v3501 = *(*int32)(unsafe.Add(mBase, uint32(v3472)+72))
	v3502 = m.T0[v3501].(func(*base.Module, int32, int32, int32) int32)(m, v21+int32(15), int32(1), v3472)
	mBase = m.M
	if v3502 == v3496 {
		v3562 = v3477
		goto L597
	} else {
		goto L606
	}
L606:
	;
	goto L599
L607:
	;
	if base.Ui32(int32(524287)) < base.Ui32(v3244) {
		v3562 = int32(6)
		goto L597
	} else {
		goto L608
	}
L608:
	;
	v3531 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	v3532 = int32(42)
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+21)) = uint8(v3532)
	v3534 = int32(413)
	*(*uint16)(unsafe.Add(mBase, uint32(v21)+19)) = uint16(v3534)
	v3536 = *(*int32)(unsafe.Add(mBase, uint32(v3399)+8))
	v3537 = *(*int32)(unsafe.Add(mBase, uint32(v3399)+12))
	*(*uint16)(unsafe.Add(mBase, uint32(v21)+24)) = uint16(v3537)
	*(*uint16)(unsafe.Add(mBase, uint32(v21)+22)) = uint16(v3536)
	v3544 = v3531<<(uint(int32(1))%32) | v3244<<(uint(int32(5))%32)
	v3545 = int32(16)
	v3546 = int32(base.Ui32(v3544) >> (uint(v3545) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+18)) = uint8(v3546)
	v3548 = int32(8)
	v3550 = int32(base.Ui32(v3544) >> (uint(v3548) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+17)) = uint8(v3550)
	v3553 = v3544 | v3545
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+16)) = uint8(v3553)
	v3558 = *(*int32)(unsafe.Add(mBase, uint32(v3399)+72))
	v3559 = m.T0[v3558].(func(*base.Module, int32, int32, int32) int32)(m, v21+v3545, int32(10), v3399)
	mBase = m.M
	if v3559 != 0 {
		goto L596
	} else {
		goto L609
	}
L609:
	;
	v3562 = v3548
	goto L597
L610:
	;
	goto L613
L611:
	;
	goto L610
L612:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3399)+92)) = v3562
	goto L611
L613:
	;
	v3671 = int32(0)
	goto L595
L614:
	;
	v3575 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if int32(2) <= v3575 {
		goto L617
	} else {
		goto L618
	}
L615:
	;
	v3671 = int32(0)
	goto L595
L616:
	;
	v3671 = base.B2i32(v3649 != int32(0))
	goto L595
L617:
	;
	v3584 = v3575 + int32(-1)
	v3589 = l0 + int32(108)
	v3590 = v3584
	v3594 = v21 + int32(16)
	goto L619
L618:
	;
	v3649 = int32(1)
	goto L616
L619:
	;
	v3603 = *(*int32)(unsafe.Add(mBase, uint32(v3589)))
	if base.Ui32(v3603) < base.Ui32(int32(16777216)) {
		goto L621
	} else {
		goto L622
	}
L620:
	;
	v3631 = *(*int32)(unsafe.Add(mBase, uint32(v32)+72))
	v3632 = m.T0[v3631].(func(*base.Module, int32, int32, int32) int32)(m, v21+int32(16), v3584*int32(3), v32)
	mBase = m.M
	if v3632 == int32(0) {
		goto L627
	} else {
		goto L628
	}
L621:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3594))) = uint8(v3603)
	v3614 = int32(base.Ui32(v3603) >> (uint(int32(16)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v3594+int32(2)))) = uint8(v3614)
	v3619 = int32(base.Ui32(v3603) >> (uint(int32(8)) % 32))
	*(*uint8)(unsafe.Add(mBase, uint32(v3594+int32(1)))) = uint8(v3619)
	v3626 = v3590 + int32(-1)
	if v3626 != 0 {
		v3589 = v3589 + int32(32)
		v3590 = v3626
		v3594 = v3594 + int32(3)
		goto L619
	} else {
		goto L626
	}
L622:
	;
	v3607 = *(*int32)(unsafe.Add(mBase, uint32(v32)+92))
	if v3607 != 0 {
		goto L624
	} else {
		goto L625
	}
L623:
	;
	v3649 = int32(0)
	goto L616
L624:
	;
	goto L623
L625:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v32)+92)) = int32(7)
	goto L624
L626:
	;
	goto L620
L627:
	;
	v3637 = *(*int32)(unsafe.Add(mBase, uint32(v32)+92))
	if v3637 != 0 {
		goto L630
	} else {
		goto L631
	}
L628:
	;
	v3649 = int32(1)
	goto L616
L629:
	;
	v3649 = int32(0)
	goto L616
L630:
	;
	goto L629
L631:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v32)+92)) = int32(8)
	goto L630
L632:
	;
	v3700 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v3700 < int32(1) {
		v3782 = v3671
		goto L635
	} else {
		goto L636
	}
L633:
	;
	goto L632
L634:
	;
	v3684 = l0 + int32(72)
	v3685 = *(*int32)(unsafe.Add(mBase, uint32(v3684)))
	F_WebPSafeFree(m, v3685)
	mBase = m.M
	v3689 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(80)))) = v3689
	*(*int64)(unsafe.Add(mBase, uint32(v3684))) = v3689
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(64)))) = v3689
	*(*int64)(unsafe.Add(mBase, uint32(v34))) = v3689
	goto L633
L635:
	;
	if v3782 == int32(0) {
		v3802 = v3782
		goto L650
	} else {
		goto L651
	}
L636:
	;
	v3710 = int32(0)
	v3715 = l0 + int32(108)
	v3717 = v3671
	goto L637
L637:
	;
	v3724 = *(*int32)(unsafe.Add(mBase, uint32(v3715)))
	if v3724 == int32(0) {
		v3735 = v3717
		goto L639
	} else {
		goto L640
	}
L638:
	;
	v3782 = v3764
	goto L635
L639:
	;
	v3737 = v3715 + int32(-20)
	if v3737 == int32(0) {
		goto L644
	} else {
		goto L645
	}
L640:
	;
	if v3717 != 0 {
		goto L641
	} else {
		goto L642
	}
L641:
	;
	v3730 = *(*int32)(unsafe.Add(mBase, uint32(v3715+int32(-4))))
	v3731 = *(*int32)(unsafe.Add(mBase, uint32(v32)+72))
	v3732 = m.T0[v3731].(func(*base.Module, int32, int32, int32) int32)(m, v3730, v3724, v32)
	mBase = m.M
	v3735 = base.B2i32(v3732 != int32(0))
	goto L639
L642:
	;
	v3735 = int32(0)
	goto L639
L643:
	;
	if v3735 != 0 {
		goto L647
	} else {
		goto L648
	}
L644:
	;
	goto L643
L645:
	;
	v3742 = v3715 + int32(-4)
	v3743 = *(*int32)(unsafe.Add(mBase, uint32(v3742)))
	F_WebPSafeFree(m, v3743)
	mBase = m.M
	v3747 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v3715+int32(4)))) = v3747
	*(*int64)(unsafe.Add(mBase, uint32(v3742))) = v3747
	*(*int64)(unsafe.Add(mBase, uint32(v3715+int32(-12)))) = v3747
	*(*int64)(unsafe.Add(mBase, uint32(v3737))) = v3747
	goto L644
L646:
	;
	v3768 = v3710 + int32(1)
	v3769 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v3768 < v3769 {
		v3710 = v3768
		v3715 = v3715 + int32(32)
		v3717 = v3764
		goto L637
	} else {
		goto L649
	}
L647:
	;
	v3759 = *(*int32)(unsafe.Add(mBase, uint32(v3679)))
	v3761 = F_WebPReportProgress(m, v32, v3759+v26, v3679)
	mBase = m.M
	v3764 = base.B2i32(v3761 != int32(0))
	goto L646
L648:
	;
	v3764 = int32(0)
	goto L646
L649:
	;
	goto L638
L650:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l0)+uint32(_c_F_VP8EncWrite[4]))) = v3386 + int32(8)
	if v3802 == int32(0) {
		goto L653
	} else {
		goto L654
	}
L651:
	;
	if v3374 == int32(0) {
		v3802 = v3782
		goto L650
	} else {
		goto L652
	}
L652:
	;
	v3793 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v21)+16)) = uint8(v3793)
	v3798 = *(*int32)(unsafe.Add(mBase, uint32(v32)+72))
	v3799 = m.T0[v3798].(func(*base.Module, int32, int32, int32) int32)(m, v21+int32(16), int32(1), v32)
	mBase = m.M
	v3802 = base.B2i32(v3799 != v3793)
	goto L650
L653:
	;
	v3815 = *(*int32)(unsafe.Add(mBase, uint32(v32)+92))
	if v3815 != 0 {
		goto L657
	} else {
		goto L658
	}
L654:
	;
	v3810 = F_WebPReportProgress(m, v32, v23+int32(19), v3679)
	mBase = m.M
	if v3810 == int32(0) {
		goto L653
	} else {
		goto L655
	}
L655:
	;
	v3837 = int32(1)
	goto L1
L656:
	;
	goto L2
L657:
	;
	goto L656
L658:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v32)+92)) = int32(8)
	goto L657
}

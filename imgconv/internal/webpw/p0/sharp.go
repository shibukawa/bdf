//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"math"
	"unsafe"
)

func F_SharpYuvConvert(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32, l15 int32) int32 {
	mBase := m.M
	_ = mBase
	var v18 int32
	_ = v18
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v27 int32
	_ = v27
	v18 = m.G0
	v19 = int32(16)
	v20 = v18 - v19
	m.G0 = v20
	*(*int32)(unsafe.Add(mBase, uint32(v20)+12)) = int32(13)
	*(*int32)(unsafe.Add(mBase, uint32(v20)+8)) = l15
	v27 = F_SharpYuvConvertWithOptions(m, l0, l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11, l12, l13, l14, v20+int32(8))
	mBase = m.M
	m.G0 = v20 + v19
	return v27
}
func F_SharpYuvConvertWithOptions(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32, l15 int32) int32 {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v134 int32
	_ = v134
	var v135 int32
	_ = v135
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v153 int32
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v179 int32
	_ = v179
	var v180 int32
	_ = v180
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v189 int32
	_ = v189
	var v191 int32
	_ = v191
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v209 int32
	_ = v209
	var v210 int32
	_ = v210
	var v213 int32
	_ = v213
	var v214 int32
	_ = v214
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v221 int32
	_ = v221
	var v222 int32
	_ = v222
	var v225 int32
	_ = v225
	var v226 int32
	_ = v226
	var v229 int32
	_ = v229
	var v230 int32
	_ = v230
	var v231 int32
	_ = v231
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
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
	var v241 int32
	_ = v241
	var v242 int32
	_ = v242
	var v243 int32
	_ = v243
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
	var v249 int32
	_ = v249
	var v250 int32
	_ = v250
	var v251 int32
	_ = v251
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	var v255 int32
	_ = v255
	var v256 int32
	_ = v256
	var v259 int32
	_ = v259
	var v261 int32
	_ = v261
	var v263 int32
	_ = v263
	var v266 int64
	_ = v266
	var v270 int32
	_ = v270
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v274 int64
	_ = v274
	var v278 int64
	_ = v278
	var v281 int32
	_ = v281
	var v282 int32
	_ = v282
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v285 int32
	_ = v285
	var v287 int32
	_ = v287
	var v288 int32
	_ = v288
	var v290 int32
	_ = v290
	var v291 int32
	_ = v291
	var v294 int64
	_ = v294
	var v298 int32
	_ = v298
	var v299 int32
	_ = v299
	var v303 int32
	_ = v303
	var v304 int64
	_ = v304
	var v307 int64
	_ = v307
	var v310 int32
	_ = v310
	var v311 int32
	_ = v311
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
	var v315 int32
	_ = v315
	var v318 int64
	_ = v318
	var v322 int32
	_ = v322
	var v323 int32
	_ = v323
	var v328 float64
	_ = v328
	var v336 int64
	_ = v336
	var v338 int64
	_ = v338
	var v339 int32
	_ = v339
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v361 int32
	_ = v361
	var v364 int32
	_ = v364
	var v370 int32
	_ = v370
	var v374 int32
	_ = v374
	var v377 int32
	_ = v377
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v393 int32
	_ = v393
	var v394 int32
	_ = v394
	var v401 int32
	_ = v401
	var v404 int32
	_ = v404
	var v405 int32
	_ = v405
	var v410 int32
	_ = v410
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v413 int32
	_ = v413
	var v468 int32
	_ = v468
	var v478 int32
	_ = v478
	var v479 int32
	_ = v479
	var v483 int32
	_ = v483
	var v484 int32
	_ = v484
	var v485 int32
	_ = v485
	var v486 int32
	_ = v486
	var v515 int32
	_ = v515
	var v518 int32
	_ = v518
	var v521 int32
	_ = v521
	var v523 int32
	_ = v523
	var v525 int32
	_ = v525
	var v528 int32
	_ = v528
	var v531 int32
	_ = v531
	var v539 int32
	_ = v539
	var v540 int32
	_ = v540
	var v547 int32
	_ = v547
	var v551 int32
	_ = v551
	var v552 int32
	_ = v552
	var v553 int32
	_ = v553
	var v557 int32
	_ = v557
	var v559 int32
	_ = v559
	var v563 int32
	_ = v563
	var v565 int32
	_ = v565
	var v571 int32
	_ = v571
	var v572 int32
	_ = v572
	var v575 int32
	_ = v575
	var v581 int32
	_ = v581
	var v584 int32
	_ = v584
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v602 int32
	_ = v602
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v611 int32
	_ = v611
	var v612 int32
	_ = v612
	var v616 int32
	_ = v616
	var v617 int32
	_ = v617
	var v623 int32
	_ = v623
	var v624 int32
	_ = v624
	var v625 int32
	_ = v625
	var v626 int32
	_ = v626
	var v636 int32
	_ = v636
	var v637 int32
	_ = v637
	var v644 int32
	_ = v644
	var v648 int32
	_ = v648
	var v649 int32
	_ = v649
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v665 int32
	_ = v665
	var v684 int32
	_ = v684
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v689 int32
	_ = v689
	var v693 int32
	_ = v693
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v703 int32
	_ = v703
	var v708 int32
	_ = v708
	var v709 int32
	_ = v709
	var v710 int32
	_ = v710
	var v719 int32
	_ = v719
	var v722 int32
	_ = v722
	var v725 int32
	_ = v725
	var v727 int32
	_ = v727
	var v729 int32
	_ = v729
	var v732 int32
	_ = v732
	var v735 int32
	_ = v735
	var v743 int32
	_ = v743
	var v744 int32
	_ = v744
	var v751 int32
	_ = v751
	var v755 int32
	_ = v755
	var v756 int32
	_ = v756
	var v757 int32
	_ = v757
	var v761 int32
	_ = v761
	var v763 int32
	_ = v763
	var v767 int32
	_ = v767
	var v769 int32
	_ = v769
	var v775 int32
	_ = v775
	var v776 int32
	_ = v776
	var v779 int32
	_ = v779
	var v785 int32
	_ = v785
	var v788 int32
	_ = v788
	var v798 int32
	_ = v798
	var v799 int32
	_ = v799
	var v806 int32
	_ = v806
	var v810 int32
	_ = v810
	var v811 int32
	_ = v811
	var v815 int32
	_ = v815
	var v816 int32
	_ = v816
	var v820 int32
	_ = v820
	var v821 int32
	_ = v821
	var v827 int32
	_ = v827
	var v828 int32
	_ = v828
	var v829 int32
	_ = v829
	var v830 int32
	_ = v830
	var v840 int32
	_ = v840
	var v841 int32
	_ = v841
	var v848 int32
	_ = v848
	var v852 int32
	_ = v852
	var v853 int32
	_ = v853
	var v857 int32
	_ = v857
	var v858 int32
	_ = v858
	var v862 int32
	_ = v862
	var v863 int32
	_ = v863
	var v869 int32
	_ = v869
	var v888 int32
	_ = v888
	var v890 int32
	_ = v890
	var v891 int32
	_ = v891
	var v893 int32
	_ = v893
	var v897 int32
	_ = v897
	var v900 int32
	_ = v900
	var v904 int32
	_ = v904
	var v907 int32
	_ = v907
	var v911 int32
	_ = v911
	var v928 int32
	_ = v928
	var v931 int32
	_ = v931
	var v1011 int64
	_ = v1011
	var v1015 int64
	_ = v1015
	var v1020 int64
	_ = v1020
	var v1027 int64
	_ = v1027
	var v1032 int32
	_ = v1032
	var v1048 int32
	_ = v1048
	var v1051 int32
	_ = v1051
	var v1062 int32
	_ = v1062
	var v1072 int32
	_ = v1072
	var v1109 int32
	_ = v1109
	var v1129 int64
	_ = v1129
	var v1132 int64
	_ = v1132
	var v1136 int64
	_ = v1136
	var v1143 int64
	_ = v1143
	var v1145 int32
	_ = v1145
	var v1154 int32
	_ = v1154
	var v1171 int32
	_ = v1171
	var v1174 int32
	_ = v1174
	var v1254 int32
	_ = v1254
	var v1255 int32
	_ = v1255
	var v1260 int32
	_ = v1260
	var v1261 int32
	_ = v1261
	var v1267 int32
	_ = v1267
	var v1268 int32
	_ = v1268
	var v1278 int32
	_ = v1278
	var v1283 int32
	_ = v1283
	var v1299 int32
	_ = v1299
	var v1302 int32
	_ = v1302
	var v1313 int32
	_ = v1313
	var v1323 int32
	_ = v1323
	var v1360 int32
	_ = v1360
	var v1380 int32
	_ = v1380
	var v1381 int32
	_ = v1381
	var v1385 int32
	_ = v1385
	var v1386 int32
	_ = v1386
	var v1391 int32
	_ = v1391
	var v1392 int32
	_ = v1392
	var v1402 int32
	_ = v1402
	var v1404 int32
	_ = v1404
	var v1413 int32
	_ = v1413
	var v1432 int32
	_ = v1432
	var v1433 int32
	_ = v1433
	var v1434 int32
	_ = v1434
	var v1437 int32
	_ = v1437
	var v1441 int32
	_ = v1441
	var v1444 int32
	_ = v1444
	var v1451 int32
	_ = v1451
	var v1452 int32
	_ = v1452
	var v1455 int32
	_ = v1455
	var v1468 int32
	_ = v1468
	var v1469 int32
	_ = v1469
	var v1471 int32
	_ = v1471
	var v1472 int32
	_ = v1472
	var v1473 int32
	_ = v1473
	var v1476 int32
	_ = v1476
	var v1477 int32
	_ = v1477
	var v1478 int32
	_ = v1478
	var v1479 int32
	_ = v1479
	var v1481 int32
	_ = v1481
	var v1483 int32
	_ = v1483
	var v1489 int32
	_ = v1489
	var v1490 int32
	_ = v1490
	var v1493 int32
	_ = v1493
	var v1494 int32
	_ = v1494
	var v1495 int32
	_ = v1495
	var v1498 int32
	_ = v1498
	var v1499 int32
	_ = v1499
	var v1500 int32
	_ = v1500
	var v1501 int32
	_ = v1501
	var v1503 int32
	_ = v1503
	var v1505 int32
	_ = v1505
	var v1511 int32
	_ = v1511
	var v1512 int32
	_ = v1512
	var v1515 int32
	_ = v1515
	var v1516 int32
	_ = v1516
	var v1517 int32
	_ = v1517
	var v1520 int32
	_ = v1520
	var v1528 int32
	_ = v1528
	var v1529 int32
	_ = v1529
	var v1530 int32
	_ = v1530
	var v1532 int32
	_ = v1532
	var v1534 int32
	_ = v1534
	var v1540 int32
	_ = v1540
	var v1549 int32
	_ = v1549
	var v1550 int32
	_ = v1550
	var v1553 int32
	_ = v1553
	var v1556 int32
	_ = v1556
	var v1563 int32
	_ = v1563
	var v1564 int32
	_ = v1564
	var v1573 int32
	_ = v1573
	var v1577 int32
	_ = v1577
	var v1637 int32
	_ = v1637
	var v1638 int32
	_ = v1638
	var v1639 int32
	_ = v1639
	var v1675 int32
	_ = v1675
	var v1676 int32
	_ = v1676
	var v1677 int32
	_ = v1677
	var v1678 int32
	_ = v1678
	var v1679 int32
	_ = v1679
	var v1680 int32
	_ = v1680
	var v1681 int32
	_ = v1681
	var v1682 int32
	_ = v1682
	var v1684 int32
	_ = v1684
	var v1688 int32
	_ = v1688
	var v1690 int32
	_ = v1690
	var v1692 int32
	_ = v1692
	var v1694 int32
	_ = v1694
	var v1698 int32
	_ = v1698
	var v1709 int32
	_ = v1709
	var v1711 int32
	_ = v1711
	var v1713 int32
	_ = v1713
	var v1715 int32
	_ = v1715
	var v1732 int32
	_ = v1732
	var v1748 int32
	_ = v1748
	var v1768 int64
	_ = v1768
	var v1834 int32
	_ = v1834
	var __phi1834 int32
	_ = __phi1834
	var v1835 int32
	_ = v1835
	var __phi1835 int32
	_ = __phi1835
	var v1866 int64
	_ = v1866
	var __phi1866 int64
	_ = __phi1866
	var v1870 int32
	_ = v1870
	var __phi1870 int32
	_ = __phi1870
	var v1892 int32
	_ = v1892
	var __phi1892 int32
	_ = __phi1892
	var v1893 int32
	_ = v1893
	var __phi1893 int32
	_ = __phi1893
	var v1899 int32
	_ = v1899
	var __phi1899 int32
	_ = __phi1899
	var v1904 int32
	_ = v1904
	var __phi1904 int32
	_ = __phi1904
	var v1927 int32
	_ = v1927
	var v1929 int32
	_ = v1929
	var v1932 int32
	_ = v1932
	var v1933 int32
	_ = v1933
	var v1934 int32
	_ = v1934
	var v1938 int32
	_ = v1938
	var v1939 int32
	_ = v1939
	var v1942 int32
	_ = v1942
	var v1944 int32
	_ = v1944
	var v1946 int32
	_ = v1946
	var v1949 int32
	_ = v1949
	var v1952 int32
	_ = v1952
	var v1953 int32
	_ = v1953
	var v1957 int32
	_ = v1957
	var v1958 int32
	_ = v1958
	var v1959 int32
	_ = v1959
	var v1962 int32
	_ = v1962
	var v1964 int32
	_ = v1964
	var v1966 int32
	_ = v1966
	var v1967 int32
	_ = v1967
	var v1968 int32
	_ = v1968
	var v1969 int32
	_ = v1969
	var v1972 int32
	_ = v1972
	var v1973 int32
	_ = v1973
	var v1975 int32
	_ = v1975
	var v1977 int32
	_ = v1977
	var v1981 int32
	_ = v1981
	var v1983 int32
	_ = v1983
	var v1987 int32
	_ = v1987
	var v1988 int32
	_ = v1988
	var v1989 int32
	_ = v1989
	var v1992 int32
	_ = v1992
	var v1994 int32
	_ = v1994
	var v1996 int32
	_ = v1996
	var v1998 int32
	_ = v1998
	var v2002 int32
	_ = v2002
	var v2003 int32
	_ = v2003
	var v2004 int32
	_ = v2004
	var v2007 int32
	_ = v2007
	var v2009 int32
	_ = v2009
	var v2011 int32
	_ = v2011
	var v2013 int32
	_ = v2013
	var v2014 int32
	_ = v2014
	var v2015 int32
	_ = v2015
	var v2018 int32
	_ = v2018
	var v2019 int32
	_ = v2019
	var v2020 int32
	_ = v2020
	var v2021 int32
	_ = v2021
	var v2025 int32
	_ = v2025
	var v2026 int32
	_ = v2026
	var v2029 int32
	_ = v2029
	var v2031 int32
	_ = v2031
	var v2033 int32
	_ = v2033
	var v2034 int32
	_ = v2034
	var v2035 int32
	_ = v2035
	var v2039 int32
	_ = v2039
	var v2040 int32
	_ = v2040
	var v2043 int32
	_ = v2043
	var v2045 int32
	_ = v2045
	var v2047 int32
	_ = v2047
	var v2049 int32
	_ = v2049
	var v2051 int32
	_ = v2051
	var v2053 int32
	_ = v2053
	var v2056 int32
	_ = v2056
	var v2057 int32
	_ = v2057
	var v2059 int32
	_ = v2059
	var v2063 int32
	_ = v2063
	var v2064 int32
	_ = v2064
	var v2067 int32
	_ = v2067
	var v2069 int32
	_ = v2069
	var v2071 int32
	_ = v2071
	var v2073 int32
	_ = v2073
	var v2077 int32
	_ = v2077
	var v2078 int32
	_ = v2078
	var v2081 int32
	_ = v2081
	var v2083 int32
	_ = v2083
	var v2085 int32
	_ = v2085
	var v2086 int32
	_ = v2086
	var v2087 int32
	_ = v2087
	var v2090 int32
	_ = v2090
	var v2091 int32
	_ = v2091
	var v2092 int32
	_ = v2092
	var v2093 int32
	_ = v2093
	var v2097 int32
	_ = v2097
	var v2098 int32
	_ = v2098
	var v2101 int32
	_ = v2101
	var v2103 int32
	_ = v2103
	var v2105 int32
	_ = v2105
	var v2106 int32
	_ = v2106
	var v2107 int32
	_ = v2107
	var v2111 int32
	_ = v2111
	var v2112 int32
	_ = v2112
	var v2115 int32
	_ = v2115
	var v2117 int32
	_ = v2117
	var v2119 int32
	_ = v2119
	var v2121 int32
	_ = v2121
	var v2123 int32
	_ = v2123
	var v2125 int32
	_ = v2125
	var v2128 int32
	_ = v2128
	var v2129 int32
	_ = v2129
	var v2131 int32
	_ = v2131
	var v2135 int32
	_ = v2135
	var v2136 int32
	_ = v2136
	var v2139 int32
	_ = v2139
	var v2141 int32
	_ = v2141
	var v2143 int32
	_ = v2143
	var v2145 int32
	_ = v2145
	var v2149 int32
	_ = v2149
	var v2150 int32
	_ = v2150
	var v2153 int32
	_ = v2153
	var v2155 int32
	_ = v2155
	var v2172 int32
	_ = v2172
	var v2175 int32
	_ = v2175
	var v2255 int32
	_ = v2255
	var v2256 int32
	_ = v2256
	var v2261 int32
	_ = v2261
	var v2262 int32
	_ = v2262
	var v2268 int32
	_ = v2268
	var v2269 int32
	_ = v2269
	var v2279 int32
	_ = v2279
	var v2284 int32
	_ = v2284
	var v2300 int32
	_ = v2300
	var v2303 int32
	_ = v2303
	var v2314 int32
	_ = v2314
	var v2324 int32
	_ = v2324
	var v2381 int32
	_ = v2381
	var v2382 int32
	_ = v2382
	var v2386 int32
	_ = v2386
	var v2387 int32
	_ = v2387
	var v2393 int32
	_ = v2393
	var v2394 int32
	_ = v2394
	var v2404 int32
	_ = v2404
	var v2406 int32
	_ = v2406
	var v2413 int32
	_ = v2413
	var v2432 int32
	_ = v2432
	var v2433 int32
	_ = v2433
	var v2434 int32
	_ = v2434
	var v2437 int32
	_ = v2437
	var v2441 int32
	_ = v2441
	var v2444 int32
	_ = v2444
	var v2451 int32
	_ = v2451
	var v2452 int32
	_ = v2452
	var v2455 int32
	_ = v2455
	var v2468 int32
	_ = v2468
	var v2469 int32
	_ = v2469
	var v2471 int32
	_ = v2471
	var v2472 int32
	_ = v2472
	var v2473 int32
	_ = v2473
	var v2476 int32
	_ = v2476
	var v2477 int32
	_ = v2477
	var v2478 int32
	_ = v2478
	var v2479 int32
	_ = v2479
	var v2481 int32
	_ = v2481
	var v2483 int32
	_ = v2483
	var v2489 int32
	_ = v2489
	var v2490 int32
	_ = v2490
	var v2493 int32
	_ = v2493
	var v2494 int32
	_ = v2494
	var v2495 int32
	_ = v2495
	var v2498 int32
	_ = v2498
	var v2499 int32
	_ = v2499
	var v2500 int32
	_ = v2500
	var v2501 int32
	_ = v2501
	var v2503 int32
	_ = v2503
	var v2505 int32
	_ = v2505
	var v2511 int32
	_ = v2511
	var v2512 int32
	_ = v2512
	var v2515 int32
	_ = v2515
	var v2516 int32
	_ = v2516
	var v2517 int32
	_ = v2517
	var v2520 int32
	_ = v2520
	var v2528 int32
	_ = v2528
	var v2529 int32
	_ = v2529
	var v2530 int32
	_ = v2530
	var v2532 int32
	_ = v2532
	var v2534 int32
	_ = v2534
	var v2540 int32
	_ = v2540
	var v2549 int32
	_ = v2549
	var v2550 int32
	_ = v2550
	var v2553 int32
	_ = v2553
	var v2556 int32
	_ = v2556
	var v2563 int32
	_ = v2563
	var v2564 int32
	_ = v2564
	var v2565 int32
	_ = v2565
	var v2566 int32
	_ = v2566
	var v2567 int64
	_ = v2567
	var v2568 int32
	_ = v2568
	var v2570 int64
	_ = v2570
	var v2571 int32
	_ = v2571
	var v2572 int32
	_ = v2572
	var v2575 int32
	_ = v2575
	var v2580 int32
	_ = v2580
	var v2592 int32
	_ = v2592
	var v2595 int32
	_ = v2595
	var v2596 int32
	_ = v2596
	var v2599 int32
	_ = v2599
	var v2601 int32
	_ = v2601
	var v2605 int32
	_ = v2605
	var v2619 int32
	_ = v2619
	var v2628 int32
	_ = v2628
	var v2682 int32
	_ = v2682
	var v2683 int32
	_ = v2683
	var v2726 int32
	_ = v2726
	var v2728 int32
	_ = v2728
	var v2810 int32
	_ = v2810
	var v2811 int32
	_ = v2811
	var v2815 int32
	_ = v2815
	var v2816 int32
	_ = v2816
	var v2821 int32
	_ = v2821
	var v2829 int32
	_ = v2829
	var v2834 int32
	_ = v2834
	var v2844 int32
	_ = v2844
	var v2849 int32
	_ = v2849
	var v2857 int32
	_ = v2857
	var v2859 int32
	_ = v2859
	var v2863 int32
	_ = v2863
	var v2877 int32
	_ = v2877
	var v2889 int32
	_ = v2889
	var v2940 int32
	_ = v2940
	var v2943 int32
	_ = v2943
	var v2967 int32
	_ = v2967
	var v2984 int32
	_ = v2984
	var v2985 int32
	_ = v2985
	var v3070 int32
	_ = v3070
	var v3072 int32
	_ = v3072
	var v3076 int32
	_ = v3076
	var v3077 int32
	_ = v3077
	var v3082 int32
	_ = v3082
	var v3090 int32
	_ = v3090
	var v3094 int32
	_ = v3094
	var v3097 int32
	_ = v3097
	var v3100 int32
	_ = v3100
	var v3105 int32
	_ = v3105
	var v3113 int32
	_ = v3113
	var v3115 int32
	_ = v3115
	var v3119 int32
	_ = v3119
	var v3217 int32
	_ = v3217
	var v3220 int32
	_ = v3220
	var v3221 int32
	_ = v3221
	var v3224 int32
	_ = v3224
	var v3226 int32
	_ = v3226
	var v3229 int32
	_ = v3229
	var v3236 int32
	_ = v3236
	var v3244 int32
	_ = v3244
	var v3246 int32
	_ = v3246
	var v3309 int32
	_ = v3309
	var v3348 int32
	_ = v3348
	var v3349 int32
	_ = v3349
	var v3430 int32
	_ = v3430
	var v3434 int32
	_ = v3434
	var v3438 int32
	_ = v3438
	var v3442 int32
	_ = v3442
	var v3452 int32
	_ = v3452
	var v3462 int32
	_ = v3462
	var v3472 int32
	_ = v3472
	var v3477 int32
	_ = v3477
	var v3482 int32
	_ = v3482
	var v3484 int32
	_ = v3484
	var v3486 int32
	_ = v3486
	var v3499 int32
	_ = v3499
	var v3507 int32
	_ = v3507
	var v3509 int32
	_ = v3509
	var v3559 int32
	_ = v3559
	var v3568 int32
	_ = v3568
	var v3575 int32
	_ = v3575
	var v3612 int32
	_ = v3612
	var v3695 int32
	_ = v3695
	var v3699 int32
	_ = v3699
	var v3703 int32
	_ = v3703
	var v3706 int32
	_ = v3706
	var v3709 int32
	_ = v3709
	var v3712 int32
	_ = v3712
	var v3722 int32
	_ = v3722
	var v3725 int32
	_ = v3725
	var v3728 int32
	_ = v3728
	var v3731 int32
	_ = v3731
	var v3738 int32
	_ = v3738
	var v3740 int32
	_ = v3740
	var v3758 int32
	_ = v3758
	var v3861 int32
	_ = v3861
	v17 = int32(0)
	if l10 == v17 {
		v3861 = v17
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v3861
L2:
	;
	if l8 == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L3
	}
L3:
	;
	if l6 == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L4
	}
L4:
	;
	if l2 == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L5
	}
L5:
	;
	if l1 == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L6
	}
L6:
	;
	if l0 == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L7
	}
L7:
	;
	if base.Ui32(l13+int32(-2147483647)) < base.Ui32(int32(-2147483646)) {
		v3861 = v17
		goto L1
	} else {
		goto L8
	}
L8:
	;
	if base.Ui32(l14+int32(-2147483647)) < base.Ui32(int32(-2147483646)) {
		v3861 = v17
		goto L1
	} else {
		goto L9
	}
L9:
	;
	if base.Ui32(int32(16)) < base.Ui32(l5) {
		v3861 = v17
		goto L1
	} else {
		goto L10
	}
L10:
	;
	if int32(1)<<(uint(l5)%32)&int32(_a_F_SharpYuvConvertWithOptions_0) == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L11
	}
L11:
	;
	if base.Ui32(int32(12)) < base.Ui32(l12) {
		v3861 = v17
		goto L1
	} else {
		goto L12
	}
L12:
	;
	if int32(1)<<(uint(l12)%32)&int32(_a_F_SharpYuvConvertWithOptions_1) == int32(0) {
		v3861 = v17
		goto L1
	} else {
		goto L13
	}
L13:
	;
	v134 = *(*int32)(unsafe.Add(mBase, uint32(l15)+4))
	v135 = *(*int32)(unsafe.Add(mBase, uint32(l15)))
	if base.Ui32(l5) < base.Ui32(int32(9)) {
		goto L14
	} else {
		goto L15
	}
L14:
	;
	if base.Ui32(l12) < base.Ui32(int32(9)) {
		goto L17
	} else {
		goto L18
	}
L15:
	;
	if (l4|l3)&int32(1) != 0 {
		v3861 = v17
		goto L1
	} else {
		goto L16
	}
L16:
	;
	goto L14
L17:
	;
	v149 = m.G1
	v152 = *(*int32)(unsafe.Add(mBase, uint32(v149)+uint32(_c_F_SharpYuvConvertWithOptions[0])))
	v153 = m.G83
	v154 = *(*int32)(unsafe.Add(mBase, uint32(v153)))
	if v152 == v154 {
		goto L20
	} else {
		goto L21
	}
L18:
	;
	if (l9|l7|l11)&int32(1) != 0 {
		v3861 = v17
		goto L1
	} else {
		goto L19
	}
L19:
	;
	goto L17
L20:
	;
	v184 = int32(-1)<<(uint(l12)%32) ^ int32(-1)
	if l5 == l12 {
		goto L24
	} else {
		goto L25
	}
L21:
	;
	v156 = m.G1
	v159 = m.G2
	v160 = m.G1
	*(*int32)(unsafe.Add(mBase, uint32(v160)+uint32(_c_F_SharpYuvConvertWithOptions[1]))) = v159 + int32(178)
	*(*int32)(unsafe.Add(mBase, uint32(v160)+uint32(_c_F_SharpYuvConvertWithOptions[2]))) = v159 + int32(179)
	*(*int32)(unsafe.Add(mBase, uint32(v160)+uint32(_c_F_SharpYuvConvertWithOptions[3]))) = v159 + int32(180)
	goto L22
L22:
	;
	F_SharpYuvInitGammaTables(m)
	mBase = m.M
	v179 = m.G83
	v180 = *(*int32)(unsafe.Add(mBase, uint32(v179)))
	*(*int32)(unsafe.Add(mBase, uint32(v156)+uint32(_c_F_SharpYuvConvertWithOptions[0]))) = v180
	goto L20
L23:
	;
	v250 = *(*int32)(unsafe.Add(mBase, uint32(v135)+44))
	v251 = *(*int32)(unsafe.Add(mBase, uint32(v135)+28))
	v252 = *(*int32)(unsafe.Add(mBase, uint32(v135)+12))
	v253 = int32(1)
	v254 = l14 + v253
	v255 = int32(-2)
	v256 = v254 & v255
	v259 = l13 + v253
	v261 = v259 & v255
	v263 = v261 * int32(3)
	v266 = base.I64_extend_i32_s(v263) << (uint(int64(2)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v266) {
		v271 = int32(0)
		goto L26
	} else {
		goto L27
	}
L24:
	;
	v230 = *(*int32)(unsafe.Add(mBase, uint32(v135)+40))
	v231 = *(*int32)(unsafe.Add(mBase, uint32(v135)+36))
	v232 = *(*int32)(unsafe.Add(mBase, uint32(v135)+32))
	v233 = *(*int32)(unsafe.Add(mBase, uint32(v135)+24))
	v234 = *(*int32)(unsafe.Add(mBase, uint32(v135)+20))
	v235 = *(*int32)(unsafe.Add(mBase, uint32(v135)+16))
	v236 = *(*int32)(unsafe.Add(mBase, uint32(v135)+8))
	v237 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
	v238 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
	v241 = v230
	v242 = v233
	v243 = v236
	v244 = v231
	v245 = v234
	v246 = v237
	v247 = v232
	v248 = v235
	v249 = v238
	goto L23
L25:
	;
	v186 = *(*int32)(unsafe.Add(mBase, uint32(v135)+40))
	v189 = int32(-1)
	v191 = int32(1) << (uint(l5+v189) % 32)
	v196 = v189<<(uint(l5)%32) ^ v189
	v197 = base.I32_div_s(v186*v184+v191, v196)
	v198 = *(*int32)(unsafe.Add(mBase, uint32(v135)+24))
	v201 = base.I32_div_s(v198*v184+v191, v196)
	v202 = *(*int32)(unsafe.Add(mBase, uint32(v135)+8))
	v205 = base.I32_div_s(v202*v184+v191, v196)
	v206 = *(*int32)(unsafe.Add(mBase, uint32(v135)+36))
	v209 = base.I32_div_s(v206*v184+v191, v196)
	v210 = *(*int32)(unsafe.Add(mBase, uint32(v135)+20))
	v213 = base.I32_div_s(v210*v184+v191, v196)
	v214 = *(*int32)(unsafe.Add(mBase, uint32(v135)+4))
	v217 = base.I32_div_s(v214*v184+v191, v196)
	v218 = *(*int32)(unsafe.Add(mBase, uint32(v135)+32))
	v221 = base.I32_div_s(v218*v184+v191, v196)
	v222 = *(*int32)(unsafe.Add(mBase, uint32(v135)+16))
	v225 = base.I32_div_s(v222*v184+v191, v196)
	v226 = *(*int32)(unsafe.Add(mBase, uint32(v135)))
	v229 = base.I32_div_s(v226*v184+v191, v196)
	v241 = v197
	v242 = v201
	v243 = v205
	v244 = v209
	v245 = v213
	v246 = v217
	v247 = v221
	v248 = v225
	v249 = v229
	goto L23
L26:
	;
	v272 = int32(0)
	v274 = base.I64_extend_i32_s(v261)
	v278 = v274 * base.I64_extend_i32_s(v256) << (uint(int64(1)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v278) {
		v284 = v272
		v285 = v272
		goto L29
	} else {
		goto L30
	}
L27:
	;
	v270 = F_dlmalloc(m, base.I32_wrap_i64(v266))
	mBase = m.M
	goto L28
L28:
	;
	v271 = v270
	goto L26
L29:
	;
	v287 = int32(1)
	v288 = v254 >> (uint(v287) % 32)
	v290 = v259 >> (uint(v287) % 32)
	v291 = int32(0)
	v294 = v274 << (uint(int64(2)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v294) {
		v299 = v291
		goto L33
	} else {
		goto L34
	}
L30:
	;
	v281 = base.I32_wrap_i64(v278)
	v282 = F_dlmalloc(m, v281)
	mBase = m.M
	goto L31
L31:
	;
	v283 = F_dlmalloc(m, v281)
	mBase = m.M
	goto L32
L32:
	;
	v284 = v282
	v285 = v283
	goto L29
L33:
	;
	v303 = v290 * int32(3)
	v304 = base.I64_extend_i32_s(v303)
	v307 = base.I64_extend_i32_s(v288) * v304 << (uint(int64(1)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v307) {
		v314 = v291
		v315 = int32(0)
		goto L36
	} else {
		goto L37
	}
L34:
	;
	v298 = F_dlmalloc(m, base.I32_wrap_i64(v294))
	mBase = m.M
	goto L35
L35:
	;
	v299 = v298
	goto L33
L36:
	;
	v318 = v304 << (uint(int64(1)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v318) {
		v323 = int32(0)
		goto L40
	} else {
		goto L41
	}
L37:
	;
	v310 = base.I32_wrap_i64(v307)
	v311 = F_dlmalloc(m, v310)
	mBase = m.M
	goto L38
L38:
	;
	v312 = F_dlmalloc(m, v310)
	mBase = m.M
	goto L39
L39:
	;
	v314 = v311
	v315 = v312
	goto L36
L40:
	;
	v328 = base.F64_mul(base.F64_mul(base.F64_convert_i32_s(v261), float64(3)), base.F64_convert_i32_s(v256))
	if base.F64_lt(v328, float64(1.8446744073709552e+19))&base.F64_ge(v328, float64(0)) == int32(0) {
		goto L44
	} else {
		goto L45
	}
L41:
	;
	v322 = F_dlmalloc(m, base.I32_wrap_i64(v318))
	mBase = m.M
	goto L42
L42:
	;
	v323 = v322
	goto L40
L43:
	;
	v339 = int32(0)
	if v271 == v339 {
		v3758 = v339
		goto L46
	} else {
		goto L47
	}
L44:
	;
	v338 = int64(0)
	goto L43
L45:
	;
	v336 = base.I64_trunc_f64_u(v328)
	v338 = v336
	goto L43
L46:
	;
	F_dlfree(m, v284)
	mBase = m.M
	goto L336
L47:
	;
	if v299 == int32(0) {
		v3758 = v339
		goto L46
	} else {
		goto L48
	}
L48:
	;
	if v315 == int32(0) {
		v3758 = v339
		goto L46
	} else {
		goto L49
	}
L49:
	;
	if v285 == int32(0) {
		v3758 = v339
		goto L46
	} else {
		goto L50
	}
L50:
	;
	if v284 == int32(0) {
		v3758 = v339
		goto L46
	} else {
		goto L51
	}
L51:
	;
	if v314 == int32(0) {
		v3758 = v339
		goto L46
	} else {
		goto L52
	}
L52:
	;
	if v323 == int32(0) {
		v3758 = v339
		goto L46
	} else {
		goto L53
	}
L53:
	;
	if l5 < int32(13) {
		goto L54
	} else {
		goto L55
	}
L54:
	;
	v359 = int32(2)
	goto L56
L55:
	;
	v359 = int32(14) - l5
	goto L56
L56:
	;
	v360 = int32(0)
	v361 = base.B2i32(v359 < v360)
	v364 = v360 - v359
	v370 = v359 + l5
	if v360 < l14 {
		goto L58
	} else {
		goto L59
	}
L57:
	;
	if v359 < v360 {
		goto L152
	} else {
		goto L153
	}
L58:
	;
	v382 = int32(1)
	v384 = v271 + v263<<(uint(v382)%32)
	v386 = v261 << (uint(int32(2)) % 32)
	v387 = v384 + v386
	if v382 < v261 {
		goto L63
	} else {
		goto L64
	}
L59:
	;
	v374 = int32(1)
	if v374 < v261 {
		goto L60
	} else {
		goto L61
	}
L60:
	;
	v377 = v261
	goto L62
L61:
	;
	v377 = v374
	goto L62
L62:
	;
	v1637 = v377
	v1638 = v261 << (uint(int32(1)) % 32)
	v1639 = v261 << (uint(int32(2)) % 32)
	goto L57
L63:
	;
	v391 = v261
	goto L65
L64:
	;
	v391 = v382
	goto L65
L65:
	;
	v392 = int32(1)
	v393 = l4 << (uint(v392) % 32)
	v394 = int32(6)
	v401 = v261 << (uint(v392) % 32)
	v404 = v271 + v401
	v405 = v271 + v386
	v410 = v271 + v259<<(uint(int32(3))%32)&int32(-16)
	v411 = l0
	v412 = l1
	v413 = l2
	v468 = v360
	v478 = v285 + v401
	v479 = v284 + v401
	v483 = v285
	v484 = v284
	v485 = v315
	v486 = v314
	goto L66
L66:
	;
	v515 = base.I32_div_s(l3, int32(2))
	if int32(8) < l5 {
		goto L69
	} else {
		goto L70
	}
L67:
	;
	v1637 = v391
	v1638 = v401
	v1639 = v386
	goto L57
L68:
	;
	if v468 == l14+int32(-1) {
		goto L98
	} else {
		goto L99
	}
L69:
	;
	v518 = v515
	goto L71
L70:
	;
	v518 = l3
	goto L71
L71:
	;
	v521 = l13 + int32(1)
	v523 = v521 & int32(-2)
	v525 = v523 << (uint(int32(2)) % 32)
	if l5 != int32(8) {
		goto L73
	} else {
		goto L74
	}
L72:
	;
	if l13&int32(1) == int32(0) {
		goto L95
	} else {
		goto L96
	}
L73:
	;
	v572 = int32(1)
	if v572 < l13 {
		goto L81
	} else {
		goto L82
	}
L74:
	;
	v528 = int32(1)
	if v528 < l13 {
		goto L75
	} else {
		goto L76
	}
L75:
	;
	v531 = l13
	goto L77
L76:
	;
	v531 = v528
	goto L77
L77:
	;
	v539 = int32(0)
	v540 = v271
	v547 = v531
	goto L78
L78:
	;
	v551 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v411+v539))))
	v552 = int32(2)
	v553 = v551 << (uint(v552) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v540))) = uint16(v553)
	v557 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v412+v539))))
	v559 = v557 << (uint(v552) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v540+v521<<(uint(int32(1))%32)&int32(-4)))) = uint16(v559)
	v563 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v413+v539))))
	v565 = v563 << (uint(v552) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v540+v525))) = uint16(v565)
	v571 = v547 + int32(-1)
	if v571 != 0 {
		v539 = v539 + v518
		v540 = v540 + v552
		v547 = v571
		goto L78
	} else {
		goto L80
	}
L80:
	;
	goto L72
L81:
	;
	v575 = l13
	goto L83
L82:
	;
	v575 = v572
	goto L83
L83:
	;
	if l5 < int32(13) {
		goto L85
	} else {
		goto L86
	}
L84:
	;
	v624 = int32(0)
	v625 = v624 - v581
	v626 = int32(1)
	v636 = v624
	v637 = v271
	v644 = v575
	goto L92
L85:
	;
	v581 = int32(2)
	goto L87
L86:
	;
	v581 = int32(14) - l5
	goto L87
L87:
	;
	if v581 < int32(0) {
		goto L84
	} else {
		goto L88
	}
L88:
	;
	v584 = int32(1)
	v594 = int32(0)
	v595 = v271
	v602 = v575
	goto L89
L89:
	;
	v606 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v411+v594))))
	v607 = v606 << (uint(v581) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v595))) = uint16(v607)
	v611 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v412+v594))))
	v612 = v611 << (uint(v581) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v595+v521<<(uint(v584)%32)&int32(-4)))) = uint16(v612)
	v616 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v413+v594))))
	v617 = v616 << (uint(v581) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v595+v525))) = uint16(v617)
	v623 = v602 + int32(-1)
	if v623 != 0 {
		v594 = v594 + v518<<(uint(v584)%32)
		v595 = v595 + int32(2)
		v602 = v623
		goto L89
	} else {
		goto L91
	}
L91:
	;
	goto L72
L92:
	;
	v648 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v411+v636))))
	v649 = int32(base.Ui32(v648) >> (uint(v625) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v637))) = uint16(v649)
	v653 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v412+v636))))
	v654 = int32(base.Ui32(v653) >> (uint(v625) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v637+v521<<(uint(v626)%32)&int32(-4)))) = uint16(v654)
	v658 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v413+v636))))
	v659 = int32(base.Ui32(v658) >> (uint(v625) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v637+v525))) = uint16(v659)
	v665 = v644 + int32(-1)
	if v665 != 0 {
		v636 = v636 + v518<<(uint(v626)%32)
		v637 = v637 + int32(2)
		v644 = v665
		goto L92
	} else {
		goto L94
	}
L93:
	;
	goto L72
L94:
	;
	goto L93
L95:
	;
	goto L68
L96:
	;
	v684 = int32(1)
	v686 = v271 + l13<<(uint(v684)%32)
	v687 = int32(-2)
	v689 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v686+v687))))
	*(*uint16)(unsafe.Add(mBase, uint32(v686))) = uint16(v689)
	v693 = v686 + v523<<(uint(v684)%32)
	v696 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v693+v687))))
	*(*uint16)(unsafe.Add(mBase, uint32(v693))) = uint16(v696)
	v700 = v686 + v521<<(uint(int32(2))%32)
	v703 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v700+v687))))
	*(*uint16)(unsafe.Add(mBase, uint32(v700))) = uint16(v703)
	goto L95
L97:
	;
	v928 = int32(0)
	v931 = v391
	goto L129
L98:
	;
	v911 = F_memcpy(m, v384, v271, v261*v394)
	mBase = m.M
	goto L97
L99:
	;
	v708 = v411 + l4
	v709 = v412 + l4
	v710 = v413 + l4
	v719 = base.I32_div_s(l3, int32(2))
	if int32(8) < l5 {
		goto L101
	} else {
		goto L102
	}
L100:
	;
	goto L97
L101:
	;
	v722 = v719
	goto L103
L102:
	;
	v722 = l3
	goto L103
L103:
	;
	v725 = l13 + int32(1)
	v727 = v725 & int32(-2)
	v729 = v727 << (uint(int32(2)) % 32)
	if l5 != int32(8) {
		goto L105
	} else {
		goto L106
	}
L104:
	;
	if l13&int32(1) == int32(0) {
		goto L127
	} else {
		goto L128
	}
L105:
	;
	v776 = int32(1)
	if v776 < l13 {
		goto L113
	} else {
		goto L114
	}
L106:
	;
	v732 = int32(1)
	if v732 < l13 {
		goto L107
	} else {
		goto L108
	}
L107:
	;
	v735 = l13
	goto L109
L108:
	;
	v735 = v732
	goto L109
L109:
	;
	v743 = int32(0)
	v744 = v384
	v751 = v735
	goto L110
L110:
	;
	v755 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v708+v743))))
	v756 = int32(2)
	v757 = v755 << (uint(v756) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v744))) = uint16(v757)
	v761 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v709+v743))))
	v763 = v761 << (uint(v756) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v744+v725<<(uint(int32(1))%32)&int32(-4)))) = uint16(v763)
	v767 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v710+v743))))
	v769 = v767 << (uint(v756) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v744+v729))) = uint16(v769)
	v775 = v751 + int32(-1)
	if v775 != 0 {
		v743 = v743 + v722
		v744 = v744 + v756
		v751 = v775
		goto L110
	} else {
		goto L112
	}
L112:
	;
	goto L104
L113:
	;
	v779 = l13
	goto L115
L114:
	;
	v779 = v776
	goto L115
L115:
	;
	if l5 < int32(13) {
		goto L117
	} else {
		goto L118
	}
L116:
	;
	v828 = int32(0)
	v829 = v828 - v785
	v830 = int32(1)
	v840 = v828
	v841 = v384
	v848 = v779
	goto L124
L117:
	;
	v785 = int32(2)
	goto L119
L118:
	;
	v785 = int32(14) - l5
	goto L119
L119:
	;
	if v785 < int32(0) {
		goto L116
	} else {
		goto L120
	}
L120:
	;
	v788 = int32(1)
	v798 = int32(0)
	v799 = v384
	v806 = v779
	goto L121
L121:
	;
	v810 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v708+v798))))
	v811 = v810 << (uint(v785) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v799))) = uint16(v811)
	v815 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v709+v798))))
	v816 = v815 << (uint(v785) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v799+v725<<(uint(v788)%32)&int32(-4)))) = uint16(v816)
	v820 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v710+v798))))
	v821 = v820 << (uint(v785) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v799+v729))) = uint16(v821)
	v827 = v806 + int32(-1)
	if v827 != 0 {
		v798 = v798 + v722<<(uint(v788)%32)
		v799 = v799 + int32(2)
		v806 = v827
		goto L121
	} else {
		goto L123
	}
L123:
	;
	goto L104
L124:
	;
	v852 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v708+v840))))
	v853 = int32(base.Ui32(v852) >> (uint(v829) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v841))) = uint16(v853)
	v857 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v709+v840))))
	v858 = int32(base.Ui32(v857) >> (uint(v829) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v841+v725<<(uint(v830)%32)&int32(-4)))) = uint16(v858)
	v862 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v710+v840))))
	v863 = int32(base.Ui32(v862) >> (uint(v829) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v841+v729))) = uint16(v863)
	v869 = v848 + int32(-1)
	if v869 != 0 {
		v840 = v840 + v722<<(uint(v830)%32)
		v841 = v841 + int32(2)
		v848 = v869
		goto L124
	} else {
		goto L126
	}
L125:
	;
	goto L104
L126:
	;
	goto L125
L127:
	;
	goto L100
L128:
	;
	v888 = int32(1)
	v890 = v384 + l13<<(uint(v888)%32)
	v891 = int32(-2)
	v893 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v890+v891))))
	*(*uint16)(unsafe.Add(mBase, uint32(v890))) = uint16(v893)
	v897 = v890 + v727<<(uint(v888)%32)
	v900 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v897+v891))))
	*(*uint16)(unsafe.Add(mBase, uint32(v897))) = uint16(v900)
	v904 = v890 + v725<<(uint(int32(2))%32)
	v907 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v904+v891))))
	*(*uint16)(unsafe.Add(mBase, uint32(v904))) = uint16(v907)
	goto L127
L129:
	;
	v1011 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v271+v928))))
	v1015 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v404+v928))))
	v1020 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v405+v928))))
	v1027 = int64(base.Ui64(v1011*int64(13933)+v1015*int64(46871)+v1020*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64))
	*(*uint16)(unsafe.Add(mBase, uint32(v484+v928))) = uint16(v1027)
	v1032 = v931 + int32(-1)
	if v1032 != 0 {
		v928 = v928 + int32(2)
		v931 = v1032
		goto L129
	} else {
		goto L131
	}
L130:
	;
	v1048 = v384
	v1051 = v410
	v1062 = v387
	v1072 = v479
	v1109 = v391
	goto L132
L131:
	;
	goto L130
L132:
	;
	v1129 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v1048))))
	v1132 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v1051))))
	v1136 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v1062))))
	v1143 = int64(base.Ui64(v1129*int64(13933)+v1132*int64(46871)+v1136*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64))
	*(*uint16)(unsafe.Add(mBase, uint32(v1072))) = uint16(v1143)
	v1145 = int32(2)
	v1154 = v1109 + int32(-1)
	if v1154 != 0 {
		v1048 = v1048 + v1145
		v1051 = v1051 + v1145
		v1062 = v1062 + v1145
		v1072 = v1072 + v1145
		v1109 = v1154
		goto L132
	} else {
		goto L134
	}
L133:
	;
	v1171 = int32(0)
	v1174 = v391
	goto L135
L134:
	;
	goto L133
L135:
	;
	v1254 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+v1171))))
	v1255 = F_SharpYuvGammaToLinear(m, v1254, v370, v134)
	mBase = m.M
	v1260 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v404+v1171))))
	v1261 = F_SharpYuvGammaToLinear(m, v1260, v370, v134)
	mBase = m.M
	v1267 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v405+v1171))))
	v1268 = F_SharpYuvGammaToLinear(m, v1267, v370, v134)
	mBase = m.M
	v1278 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v1255)*int64(13933)+base.I64_extend_i32_u(v1261)*int64(46871)+base.I64_extend_i32_u(v1268)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v370, v134)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v483+v1171))) = uint16(v1278)
	v1283 = v1174 + int32(-1)
	if v1283 != 0 {
		v1171 = v1171 + int32(2)
		v1174 = v1283
		goto L135
	} else {
		goto L137
	}
L136:
	;
	v1299 = v384
	v1302 = v410
	v1313 = v387
	v1323 = v478
	v1360 = v391
	goto L138
L137:
	;
	goto L136
L138:
	;
	v1380 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1299))))
	v1381 = F_SharpYuvGammaToLinear(m, v1380, v370, v134)
	mBase = m.M
	v1385 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1302))))
	v1386 = F_SharpYuvGammaToLinear(m, v1385, v370, v134)
	mBase = m.M
	v1391 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1313))))
	v1392 = F_SharpYuvGammaToLinear(m, v1391, v370, v134)
	mBase = m.M
	v1402 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v1381)*int64(13933)+base.I64_extend_i32_u(v1386)*int64(46871)+base.I64_extend_i32_u(v1392)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v370, v134)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v1323))) = uint16(v1402)
	v1404 = int32(2)
	v1413 = v1360 + int32(-1)
	if v1413 != 0 {
		v1299 = v1299 + v1404
		v1302 = v1302 + v1404
		v1313 = v1313 + v1404
		v1323 = v1323 + v1404
		v1360 = v1413
		goto L138
	} else {
		goto L140
	}
L139:
	;
	if l5 < int32(13) {
		goto L142
	} else {
		goto L143
	}
L140:
	;
	goto L139
L141:
	;
	v1564 = F_memcpy(m, v486, v485, v290*v394)
	mBase = m.M
	v1573 = v303 << (uint(int32(1)) % 32)
	v1577 = v468 + int32(2)
	if v1577 < l14 {
		v411 = v411 + v393
		v412 = v412 + v393
		v413 = v413 + v393
		v468 = v1577
		v478 = v478 + v386
		v479 = v479 + v386
		v483 = v483 + v386
		v484 = v484 + v386
		v485 = v485 + v1573
		v486 = v1564 + v1573
		goto L66
	} else {
		goto L151
	}
L142:
	;
	v1432 = int32(2)
	goto L144
L143:
	;
	v1432 = int32(14) - l5
	goto L144
L144:
	;
	v1433 = v1432 + l5
	v1434 = int32(1)
	if v1434 < v290 {
		goto L145
	} else {
		goto L146
	}
L145:
	;
	v1437 = v290
	goto L147
L146:
	;
	v1437 = v1434
	goto L147
L147:
	;
	v1441 = v290 << (uint(int32(2)) % 32)
	v1444 = v290 << (uint(int32(3)) % 32)
	v1451 = v485
	v1452 = int32(0)
	v1455 = v1437
	goto L148
L148:
	;
	v1468 = v384 + v1452
	v1469 = int32(2)
	v1471 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1468+v1469))))
	v1472 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1468))))
	v1473 = v271 + v1452
	v1476 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1473+v1469))))
	v1477 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1473))))
	v1478 = F_SharpYuvGammaToLinear(m, v1477, v1433, v134)
	mBase = m.M
	v1479 = F_SharpYuvGammaToLinear(m, v1476, v1433, v134)
	mBase = m.M
	v1481 = F_SharpYuvGammaToLinear(m, v1472, v1433, v134)
	mBase = m.M
	v1483 = F_SharpYuvGammaToLinear(m, v1471, v1433, v134)
	mBase = m.M
	v1489 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v1478+v1479+v1481+v1483+v1469)>>(uint(v1469)%32)), v1433, v134)
	mBase = m.M
	v1490 = v384 + v1441 + v1452
	v1493 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1490+v1469))))
	v1494 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1490))))
	v1495 = v271 + v1441 + v1452
	v1498 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1495+v1469))))
	v1499 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1495))))
	v1500 = F_SharpYuvGammaToLinear(m, v1499, v1433, v134)
	mBase = m.M
	v1501 = F_SharpYuvGammaToLinear(m, v1498, v1433, v134)
	mBase = m.M
	v1503 = F_SharpYuvGammaToLinear(m, v1494, v1433, v134)
	mBase = m.M
	v1505 = F_SharpYuvGammaToLinear(m, v1493, v1433, v134)
	mBase = m.M
	v1511 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v1500+v1501+v1503+v1505+v1469)>>(uint(v1469)%32)), v1433, v134)
	mBase = m.M
	v1512 = v384 + v1444 + v1452
	v1515 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1512+v1469))))
	v1516 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1512))))
	v1517 = v271 + v1444 + v1452
	v1520 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1517+v1469))))
	v1528 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1517))))
	v1529 = F_SharpYuvGammaToLinear(m, v1528, v1433, v134)
	mBase = m.M
	v1530 = F_SharpYuvGammaToLinear(m, v1520, v1433, v134)
	mBase = m.M
	v1532 = F_SharpYuvGammaToLinear(m, v1516, v1433, v134)
	mBase = m.M
	v1534 = F_SharpYuvGammaToLinear(m, v1515, v1433, v134)
	mBase = m.M
	v1540 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v1529+v1530+v1532+v1534+v1469)>>(uint(v1469)%32)), v1433, v134)
	mBase = m.M
	v1549 = base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v1489)*int64(13933)+base.I64_extend_i32_u(v1511)*int64(46871)+base.I64_extend_i32_u(v1540)*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64)))
	v1550 = v1489 - v1549
	*(*uint16)(unsafe.Add(mBase, uint32(v1451))) = uint16(v1550)
	v1553 = v1511 - v1549
	*(*uint16)(unsafe.Add(mBase, uint32(v1451+v290<<(uint(int32(1))%32)))) = uint16(v1553)
	v1556 = v1540 - v1549
	*(*uint16)(unsafe.Add(mBase, uint32(v1451+v1441))) = uint16(v1556)
	v1563 = v1455 + int32(-1)
	if v1563 != 0 {
		v1451 = v1451 + v1469
		v1452 = v1452 + int32(4)
		v1455 = v1563
		goto L148
	} else {
		goto L150
	}
L149:
	;
	goto L141
L150:
	;
	goto L149
L151:
	;
	goto L67
L152:
	;
	v1675 = v250 >> (uint(v364) % 32)
	goto L154
L153:
	;
	v1675 = v250 << (uint(v359) % 32)
	goto L154
L154:
	;
	if v359 < v360 {
		goto L155
	} else {
		goto L156
	}
L155:
	;
	v1676 = v251 >> (uint(v364) % 32)
	goto L157
L156:
	;
	v1676 = v251 << (uint(v359) % 32)
	goto L157
L157:
	;
	if v359 < v360 {
		goto L158
	} else {
		goto L159
	}
L158:
	;
	v1677 = v252 >> (uint(v364) % 32)
	goto L160
L159:
	;
	v1677 = v252 << (uint(v359) % 32)
	goto L160
L160:
	;
	v1678 = int32(1)
	v1679 = v261 << (uint(v1678) % 32)
	v1680 = v271 + v1679
	v1681 = int32(-1)
	v1682 = v261 + v1681
	v1684 = v1682 << (uint(v1678) % 32)
	v1688 = v271 + v263<<(uint(v1678)%32)
	v1690 = v1680 + v1679
	v1692 = int32(2)
	v1694 = v1688 + v1679
	v1698 = v1694 + v1679
	v1709 = v1681 << (uint(v370) % 32)
	v1711 = v1709 ^ v1681
	v1713 = v1682 >> (uint(v1678) % 32)
	v1715 = int32(base.Ui32(v259) >> (uint(v1678) % 32))
	v1732 = (v290 + v1681) << (uint(v1678) % 32)
	v1748 = int32(0)
	v1768 = int64(-1)
	goto L161
L161:
	;
	__phi1834 = v315
	__phi1835 = v314
	__phi1866 = int64(0)
	__phi1870 = v314
	__phi1892 = v314
	__phi1893 = int32(0)
	__phi1899 = v285
	__phi1904 = v284
	v1834 = __phi1834
	v1835 = __phi1835
	v1866 = __phi1866
	v1870 = __phi1870
	v1892 = __phi1892
	v1893 = __phi1893
	v1899 = __phi1899
	v1904 = __phi1904
	goto L163
L162:
	;
	v2592 = int32(1)
	if v2592 < l14 {
		goto L262
	} else {
		goto L263
	}
L163:
	;
	v1927 = int32(0)
	v1929 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1892))))
	v1932 = int32(2)
	v1933 = v1929*int32(3) + v1932
	v1934 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1870))))
	v1938 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1904))))
	v1939 = (v1933+v1934)>>(uint(v1932)%32) + v1938
	if v1939 < v1927 {
		goto L165
	} else {
		goto L166
	}
L164:
	;
	if v1748 == int32(0) {
		v1748 = int32(1)
		v1768 = v2570
		goto L161
	} else {
		goto L257
	}
L165:
	;
	v1942 = v1927
	goto L167
L166:
	;
	v1942 = v1711
	goto L167
L167:
	;
	if v1939&v1709 != 0 {
		goto L168
	} else {
		goto L169
	}
L168:
	;
	v1944 = v1942
	goto L170
L169:
	;
	v1944 = v1939
	goto L170
L170:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v271))) = uint16(v1944)
	v1946 = int32(0)
	if v1893 < v256+int32(-2) {
		goto L171
	} else {
		goto L172
	}
L171:
	;
	v1949 = v303
	goto L173
L172:
	;
	v1949 = v1946
	goto L173
L173:
	;
	v1952 = v1892 + v1949<<(uint(int32(1))%32)
	v1953 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1952))))
	v1957 = v1904 + v1679
	v1958 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1957))))
	v1959 = (v1933+v1953)>>(uint(int32(2))%32) + v1958
	if v1959 < int32(0) {
		goto L174
	} else {
		goto L175
	}
L174:
	;
	v1962 = v1946
	goto L176
L175:
	;
	v1962 = v1711
	goto L176
L176:
	;
	if v1959&v1709 != 0 {
		goto L177
	} else {
		goto L178
	}
L177:
	;
	v1964 = v1962
	goto L179
L178:
	;
	v1964 = v1959
	goto L179
L179:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1688))) = uint16(v1964)
	v1966 = int32(2)
	v1967 = v1904 + v1966
	v1968 = m.G84
	v1969 = *(*int32)(unsafe.Add(mBase, uint32(v1968)))
	m.T0[v1969].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v1892, v1870, v1713, v1967, v271+v1692, v370)
	mBase = m.M
	v1972 = v1957 + v1966
	v1973 = *(*int32)(unsafe.Add(mBase, uint32(v1968)))
	m.T0[v1973].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v1892, v1952, v1713, v1972, v1688+v1692, v370)
	mBase = m.M
	v1975 = int32(0)
	v1977 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1892+v1732))))
	v1981 = v1977*int32(3) + v1966
	v1983 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1870+v1732))))
	v1987 = v1904 + v1684
	v1988 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1987))))
	v1989 = (v1981+v1983)>>(uint(v1966)%32) + v1988
	if v1989 < v1975 {
		goto L180
	} else {
		goto L181
	}
L180:
	;
	v1992 = v1975
	goto L182
L181:
	;
	v1992 = v1711
	goto L182
L182:
	;
	if v1989&v1709 != 0 {
		goto L183
	} else {
		goto L184
	}
L183:
	;
	v1994 = v1992
	goto L185
L184:
	;
	v1994 = v1989
	goto L185
L185:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v271+v1684))) = uint16(v1994)
	v1996 = int32(0)
	v1998 = int32(*(*int16)(unsafe.Add(mBase, uint32(v1952+v1732))))
	v2002 = v1987 + v1679
	v2003 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2002))))
	v2004 = (v1981+v1998)>>(uint(int32(2))%32) + v2003
	if v2004 < v1996 {
		goto L186
	} else {
		goto L187
	}
L186:
	;
	v2007 = v1996
	goto L188
L187:
	;
	v2007 = v1711
	goto L188
L188:
	;
	if v2004&v1709 != 0 {
		goto L189
	} else {
		goto L190
	}
L189:
	;
	v2009 = v2007
	goto L191
L190:
	;
	v2009 = v2004
	goto L191
L191:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1688+v1684))) = uint16(v2009)
	v2011 = int32(0)
	v2013 = v290 << (uint(int32(1)) % 32)
	v2014 = v1892 + v2013
	v2015 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2014))))
	v2018 = int32(2)
	v2019 = v2015*int32(3) + v2018
	v2020 = v1870 + v2013
	v2021 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2020))))
	v2025 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1904))))
	v2026 = (v2019+v2021)>>(uint(v2018)%32) + v2025
	if v2026 < v2011 {
		goto L192
	} else {
		goto L193
	}
L192:
	;
	v2029 = v2011
	goto L194
L193:
	;
	v2029 = v1711
	goto L194
L194:
	;
	if v2026&v1709 != 0 {
		goto L195
	} else {
		goto L196
	}
L195:
	;
	v2031 = v2029
	goto L197
L196:
	;
	v2031 = v2026
	goto L197
L197:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1680))) = uint16(v2031)
	v2033 = int32(0)
	v2034 = v1952 + v2013
	v2035 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2034))))
	v2039 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1957))))
	v2040 = (v2019+v2035)>>(uint(int32(2))%32) + v2039
	if v2040 < v2033 {
		goto L198
	} else {
		goto L199
	}
L198:
	;
	v2043 = v2033
	goto L200
L199:
	;
	v2043 = v1711
	goto L200
L200:
	;
	if v2040&v1709 != 0 {
		goto L201
	} else {
		goto L202
	}
L201:
	;
	v2045 = v2043
	goto L203
L202:
	;
	v2045 = v2040
	goto L203
L203:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1694))) = uint16(v2045)
	v2047 = *(*int32)(unsafe.Add(mBase, uint32(v1968)))
	m.T0[v2047].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v2014, v2020, v1713, v1967, v1680+v1692, v370)
	mBase = m.M
	v2049 = *(*int32)(unsafe.Add(mBase, uint32(v1968)))
	m.T0[v2049].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v2014, v2034, v1713, v1972, v1694+v1692, v370)
	mBase = m.M
	v2051 = int32(0)
	v2053 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2014+v1732))))
	v2056 = int32(2)
	v2057 = v2053*int32(3) + v2056
	v2059 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2020+v1732))))
	v2063 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1987))))
	v2064 = (v2057+v2059)>>(uint(v2056)%32) + v2063
	if v2064 < v2051 {
		goto L204
	} else {
		goto L205
	}
L204:
	;
	v2067 = v2051
	goto L206
L205:
	;
	v2067 = v1711
	goto L206
L206:
	;
	if v2064&v1709 != 0 {
		goto L207
	} else {
		goto L208
	}
L207:
	;
	v2069 = v2067
	goto L209
L208:
	;
	v2069 = v2064
	goto L209
L209:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1680+v1684))) = uint16(v2069)
	v2071 = int32(0)
	v2073 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2034+v1732))))
	v2077 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2002))))
	v2078 = (v2057+v2073)>>(uint(int32(2))%32) + v2077
	if v2078 < v2071 {
		goto L210
	} else {
		goto L211
	}
L210:
	;
	v2081 = v2071
	goto L212
L211:
	;
	v2081 = v1711
	goto L212
L212:
	;
	if v2078&v1709 != 0 {
		goto L213
	} else {
		goto L214
	}
L213:
	;
	v2083 = v2081
	goto L215
L214:
	;
	v2083 = v2078
	goto L215
L215:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1694+v1684))) = uint16(v2083)
	v2085 = int32(0)
	v2086 = v2014 + v2013
	v2087 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2086))))
	v2090 = int32(2)
	v2091 = v2087*int32(3) + v2090
	v2092 = v2020 + v2013
	v2093 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2092))))
	v2097 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1904))))
	v2098 = (v2091+v2093)>>(uint(v2090)%32) + v2097
	if v2098 < v2085 {
		goto L216
	} else {
		goto L217
	}
L216:
	;
	v2101 = v2085
	goto L218
L217:
	;
	v2101 = v1711
	goto L218
L218:
	;
	if v2098&v1709 != 0 {
		goto L219
	} else {
		goto L220
	}
L219:
	;
	v2103 = v2101
	goto L221
L220:
	;
	v2103 = v2098
	goto L221
L221:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1690))) = uint16(v2103)
	v2105 = int32(0)
	v2106 = v2034 + v2013
	v2107 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2106))))
	v2111 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1957))))
	v2112 = (v2091+v2107)>>(uint(int32(2))%32) + v2111
	if v2112 < v2105 {
		goto L222
	} else {
		goto L223
	}
L222:
	;
	v2115 = v2105
	goto L224
L223:
	;
	v2115 = v1711
	goto L224
L224:
	;
	if v2112&v1709 != 0 {
		goto L225
	} else {
		goto L226
	}
L225:
	;
	v2117 = v2115
	goto L227
L226:
	;
	v2117 = v2112
	goto L227
L227:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1698))) = uint16(v2117)
	v2119 = *(*int32)(unsafe.Add(mBase, uint32(v1968)))
	m.T0[v2119].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v2086, v2092, v1713, v1967, v1690+v1692, v370)
	mBase = m.M
	v2121 = *(*int32)(unsafe.Add(mBase, uint32(v1968)))
	m.T0[v2121].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v2086, v2106, v1713, v1972, v1698+v1692, v370)
	mBase = m.M
	v2123 = int32(0)
	v2125 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2086+v1732))))
	v2128 = int32(2)
	v2129 = v2125*int32(3) + v2128
	v2131 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2092+v1732))))
	v2135 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1987))))
	v2136 = (v2129+v2131)>>(uint(v2128)%32) + v2135
	if v2136 < v2123 {
		goto L228
	} else {
		goto L229
	}
L228:
	;
	v2139 = v2123
	goto L230
L229:
	;
	v2139 = v1711
	goto L230
L230:
	;
	if v2136&v1709 != 0 {
		goto L231
	} else {
		goto L232
	}
L231:
	;
	v2141 = v2139
	goto L233
L232:
	;
	v2141 = v2136
	goto L233
L233:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1690+v1684))) = uint16(v2141)
	v2143 = int32(0)
	v2145 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2106+v1732))))
	v2149 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2002))))
	v2150 = (v2129+v2145)>>(uint(int32(2))%32) + v2149
	if v2150 < v2143 {
		goto L234
	} else {
		goto L235
	}
L234:
	;
	v2153 = v2143
	goto L236
L235:
	;
	v2153 = v1711
	goto L236
L236:
	;
	if v2150&v1709 != 0 {
		goto L237
	} else {
		goto L238
	}
L237:
	;
	v2155 = v2153
	goto L239
L238:
	;
	v2155 = v2150
	goto L239
L239:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1698+v1684))) = uint16(v2155)
	v2172 = v1927
	v2175 = v1637
	goto L240
L240:
	;
	v2255 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+v2172))))
	v2256 = F_SharpYuvGammaToLinear(m, v2255, v370, v134)
	mBase = m.M
	v2261 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+v1715<<(uint(v1692)%32)+v2172))))
	v2262 = F_SharpYuvGammaToLinear(m, v2261, v370, v134)
	mBase = m.M
	v2268 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+v1639+v2172))))
	v2269 = F_SharpYuvGammaToLinear(m, v2268, v370, v134)
	mBase = m.M
	v2279 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v2256)*int64(13933)+base.I64_extend_i32_u(v2262)*int64(46871)+base.I64_extend_i32_u(v2269)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v370, v134)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v299+v2172))) = uint16(v2279)
	v2284 = v2175 + int32(-1)
	if v2284 != 0 {
		v2172 = v2172 + int32(2)
		v2175 = v2284
		goto L240
	} else {
		goto L242
	}
L241:
	;
	v2300 = v1688
	v2303 = v271 + v1715<<(uint(int32(4))%32)
	v2314 = v299 + v1679
	v2324 = v1637
	goto L243
L242:
	;
	goto L241
L243:
	;
	v2381 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2300))))
	v2382 = F_SharpYuvGammaToLinear(m, v2381, v370, v134)
	mBase = m.M
	v2386 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2303))))
	v2387 = F_SharpYuvGammaToLinear(m, v2386, v370, v134)
	mBase = m.M
	v2393 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2300+v1639))))
	v2394 = F_SharpYuvGammaToLinear(m, v2393, v370, v134)
	mBase = m.M
	v2404 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v2382)*int64(13933)+base.I64_extend_i32_u(v2387)*int64(46871)+base.I64_extend_i32_u(v2394)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v370, v134)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v2314))) = uint16(v2404)
	v2406 = int32(2)
	v2413 = v2324 + int32(-1)
	if v2413 != 0 {
		v2300 = v2300 + v2406
		v2303 = v2303 + v2406
		v2314 = v2314 + v2406
		v2324 = v2413
		goto L243
	} else {
		goto L245
	}
L244:
	;
	if l5 < int32(13) {
		goto L247
	} else {
		goto L248
	}
L245:
	;
	goto L244
L246:
	;
	v2564 = m.G85
	v2565 = *(*int32)(unsafe.Add(mBase, uint32(v2564)))
	v2566 = m.G86
	v2567 = m.T0[v2565].(func(*base.Module, int32, int32, int32, int32, int32) int64)(m, v1899, v299, v1904, v1638, v370)
	mBase = m.M
	v2568 = *(*int32)(unsafe.Add(mBase, uint32(v2566)))
	m.T0[v2568].(func(*base.Module, int32, int32, int32, int32))(m, v1834, v323, v1835, v303)
	mBase = m.M
	v2570 = v2567 + v1866
	v2571 = int32(1)
	v2572 = v303 << (uint(v2571) % 32)
	v2575 = v1638 << (uint(v2571) % 32)
	v2580 = v1893 + int32(2)
	if v2580 < v256 {
		__phi1834 = v1834 + v2572
		__phi1835 = v1835 + v2572
		__phi1866 = v2570
		__phi1870 = v1892
		__phi1892 = v1952
		__phi1893 = v2580
		__phi1899 = v1899 + v2575
		__phi1904 = v1904 + v2575
		v1834 = __phi1834
		v1835 = __phi1835
		v1866 = __phi1866
		v1870 = __phi1870
		v1892 = __phi1892
		v1893 = __phi1893
		v1899 = __phi1899
		v1904 = __phi1904
		goto L163
	} else {
		goto L256
	}
L247:
	;
	v2432 = int32(2)
	goto L249
L248:
	;
	v2432 = int32(14) - l5
	goto L249
L249:
	;
	v2433 = v2432 + l5
	v2434 = int32(1)
	if v2434 < v290 {
		goto L250
	} else {
		goto L251
	}
L250:
	;
	v2437 = v290
	goto L252
L251:
	;
	v2437 = v2434
	goto L252
L252:
	;
	v2441 = v290 << (uint(int32(2)) % 32)
	v2444 = v290 << (uint(int32(3)) % 32)
	v2451 = v323
	v2452 = int32(0)
	v2455 = v2437
	goto L253
L253:
	;
	v2468 = v1688 + v2452
	v2469 = int32(2)
	v2471 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2468+v2469))))
	v2472 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2468))))
	v2473 = v271 + v2452
	v2476 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2473+v2469))))
	v2477 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2473))))
	v2478 = F_SharpYuvGammaToLinear(m, v2477, v2433, v134)
	mBase = m.M
	v2479 = F_SharpYuvGammaToLinear(m, v2476, v2433, v134)
	mBase = m.M
	v2481 = F_SharpYuvGammaToLinear(m, v2472, v2433, v134)
	mBase = m.M
	v2483 = F_SharpYuvGammaToLinear(m, v2471, v2433, v134)
	mBase = m.M
	v2489 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v2478+v2479+v2481+v2483+v2469)>>(uint(v2469)%32)), v2433, v134)
	mBase = m.M
	v2490 = v1688 + v2441 + v2452
	v2493 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2490+v2469))))
	v2494 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2490))))
	v2495 = v271 + v2441 + v2452
	v2498 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2495+v2469))))
	v2499 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2495))))
	v2500 = F_SharpYuvGammaToLinear(m, v2499, v2433, v134)
	mBase = m.M
	v2501 = F_SharpYuvGammaToLinear(m, v2498, v2433, v134)
	mBase = m.M
	v2503 = F_SharpYuvGammaToLinear(m, v2494, v2433, v134)
	mBase = m.M
	v2505 = F_SharpYuvGammaToLinear(m, v2493, v2433, v134)
	mBase = m.M
	v2511 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v2500+v2501+v2503+v2505+v2469)>>(uint(v2469)%32)), v2433, v134)
	mBase = m.M
	v2512 = v1688 + v2444 + v2452
	v2515 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2512+v2469))))
	v2516 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2512))))
	v2517 = v271 + v2444 + v2452
	v2520 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2517+v2469))))
	v2528 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2517))))
	v2529 = F_SharpYuvGammaToLinear(m, v2528, v2433, v134)
	mBase = m.M
	v2530 = F_SharpYuvGammaToLinear(m, v2520, v2433, v134)
	mBase = m.M
	v2532 = F_SharpYuvGammaToLinear(m, v2516, v2433, v134)
	mBase = m.M
	v2534 = F_SharpYuvGammaToLinear(m, v2515, v2433, v134)
	mBase = m.M
	v2540 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v2529+v2530+v2532+v2534+v2469)>>(uint(v2469)%32)), v2433, v134)
	mBase = m.M
	v2549 = base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v2489)*int64(13933)+base.I64_extend_i32_u(v2511)*int64(46871)+base.I64_extend_i32_u(v2540)*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64)))
	v2550 = v2489 - v2549
	*(*uint16)(unsafe.Add(mBase, uint32(v2451))) = uint16(v2550)
	v2553 = v2511 - v2549
	*(*uint16)(unsafe.Add(mBase, uint32(v2451+v290<<(uint(int32(1))%32)))) = uint16(v2553)
	v2556 = v2540 - v2549
	*(*uint16)(unsafe.Add(mBase, uint32(v2451+v2441))) = uint16(v2556)
	v2563 = v2455 + int32(-1)
	if v2563 != 0 {
		v2451 = v2451 + v2469
		v2452 = v2452 + int32(4)
		v2455 = v2563
		goto L253
	} else {
		goto L255
	}
L254:
	;
	goto L246
L255:
	;
	goto L254
L256:
	;
	goto L164
L257:
	;
	if base.Ui64(v2570) < base.Ui64(v338) {
		goto L258
	} else {
		goto L259
	}
L258:
	;
	goto L162
L259:
	;
	if base.Ui64(v1768) < base.Ui64(v2570) {
		goto L258
	} else {
		goto L260
	}
L260:
	;
	if base.Ui32(v1748) < base.Ui32(int32(3)) {
		v1748 = v1748 + int32(1)
		v1768 = v2570
		goto L161
	} else {
		goto L261
	}
L261:
	;
	goto L258
L262:
	;
	v2595 = l14
	goto L264
L263:
	;
	v2595 = v2592
	goto L264
L264:
	;
	v2596 = int32(1)
	if v2596 < l13 {
		goto L265
	} else {
		goto L266
	}
L265:
	;
	v2599 = l13
	goto L267
L266:
	;
	v2599 = v2596
	goto L267
L267:
	;
	v2601 = v359 + int32(16)
	v2605 = int32(1) << (uint(v359+int32(15)) % 32)
	if int32(8) < l12 {
		goto L269
	} else {
		goto L270
	}
L268:
	;
	v3217 = int32(1)
	if v3217 < v288 {
		goto L298
	} else {
		goto L299
	}
L269:
	;
	v2877 = l6
	v2889 = v314
	v2940 = v284
	v2943 = int32(0)
	goto L283
L270:
	;
	v2619 = l6
	v2628 = v314
	v2682 = v284
	v2683 = int32(0)
	goto L271
L271:
	;
	v2726 = int32(0)
	v2728 = v2682
	goto L273
L273:
	;
	v2810 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2628+v2726&int32(-2)))))
	v2811 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2728))))
	v2815 = int32(1)
	v2816 = int32(base.Ui32(v2726) >> (uint(v2815) % 32))
	v2821 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2628+(v2816+v290)<<(uint(v2815)%32)))))
	v2829 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2628+(v2816+v261)<<(uint(v2815)%32)))))
	v2834 = (v249*(v2810+v2811) + v2605 + v246*(v2821+v2811) + v243*(v2829+v2811) + v1677) >> (uint(v2601) % 32)
	if base.Ui32(v2834&int32(_a_F_SharpYuvConvertWithOptions_2)) < base.Ui32(int32(256)) {
		goto L275
	} else {
		goto L276
	}
L274:
	;
	if v2683&int32(1) != 0 {
		goto L279
	} else {
		goto L280
	}
L275:
	;
	v2844 = v2834
	goto L277
L276:
	;
	v2844 = int32(base.Ui32(base.I32_extend16_s(v2834)^int32(-1)) >> (uint(int32(15)) % 32))
	goto L277
L277:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v2619+v2726))) = uint8(v2844)
	v2849 = v2726 + int32(1)
	if v2599 != v2849 {
		v2726 = v2849
		v2728 = v2728 + int32(2)
		goto L273
	} else {
		goto L278
	}
L278:
	;
	goto L274
L279:
	;
	v2857 = int32(3)
	goto L281
L280:
	;
	v2857 = int32(0)
	goto L281
L281:
	;
	v2859 = int32(1)
	v2863 = v2683 + v2859
	if v2863 != v2595 {
		v2619 = v2619 + l7
		v2628 = v2628 + v2857*v290<<(uint(v2859)%32)
		v2682 = v2682 + v259<<(uint(int32(1))%32)&int32(-4)
		v2683 = v2863
		goto L271
	} else {
		goto L282
	}
L282:
	;
	goto L268
L283:
	;
	v2967 = int32(0)
	v2984 = v2967
	v2985 = v2967
	goto L285
L284:
	;
	goto L268
L285:
	;
	v3070 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2889+v2984&int32(-2)))))
	v3072 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2940+v2985))))
	v3076 = int32(1)
	v3077 = int32(base.Ui32(v2984) >> (uint(v3076) % 32))
	v3082 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2889+(v3077+v290)<<(uint(v3076)%32)))))
	v3090 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2889+(v3077+v261)<<(uint(v3076)%32)))))
	v3094 = (v1677 + v2605 + (v3070+v3072)*v249 + (v3082+v3072)*v246 + (v3090+v3072)*v243) >> (uint(v2601) % 32)
	if v184 < base.I32_extend16_s(v3094) {
		goto L287
	} else {
		goto L288
	}
L286:
	;
	if v2943&int32(1) != 0 {
		goto L294
	} else {
		goto L295
	}
L287:
	;
	v3097 = v184
	goto L289
L288:
	;
	v3097 = v3094
	goto L289
L289:
	;
	if v3094&int32(_a_F_SharpYuvConvertWithOptions_3) != 0 {
		goto L290
	} else {
		goto L291
	}
L290:
	;
	v3100 = int32(0)
	goto L292
L291:
	;
	v3100 = v3097
	goto L292
L292:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2877+v2985))) = uint16(v3100)
	v3105 = v2984 + int32(1)
	if v2599 != v3105 {
		v2984 = v3105
		v2985 = v2985 + int32(2)
		goto L285
	} else {
		goto L293
	}
L293:
	;
	goto L286
L294:
	;
	v3113 = int32(3)
	goto L296
L295:
	;
	v3113 = int32(0)
	goto L296
L296:
	;
	v3115 = int32(1)
	v3119 = v2943 + v3115
	if v3119 != v2595 {
		v2877 = v2877 + l7
		v2889 = v2889 + v3113*v290<<(uint(v3115)%32)
		v2940 = v2940 + v259<<(uint(int32(1))%32)&int32(-4)
		v2943 = v3119
		goto L283
	} else {
		goto L297
	}
L297:
	;
	goto L284
L298:
	;
	v3220 = v288
	goto L300
L299:
	;
	v3220 = v3217
	goto L300
L300:
	;
	v3221 = int32(1)
	if v3221 < v290 {
		goto L301
	} else {
		goto L302
	}
L301:
	;
	v3224 = v290
	goto L303
L302:
	;
	v3224 = v3221
	goto L303
L303:
	;
	v3226 = v290 * int32(6)
	if int32(8) < l12 {
		goto L304
	} else {
		goto L305
	}
L304:
	;
	v3486 = int32(1)
	v3499 = int32(0)
	v3507 = l8
	v3509 = l10
	v3559 = v314 + v290<<(uint(v3486)%32)
	v3568 = v314
	v3575 = v314 + v259<<(uint(v3486)%32)&int32(-4)
	goto L318
L305:
	;
	v3229 = int32(1)
	v3236 = int32(0)
	v3244 = l8
	v3246 = l10
	v3309 = v314
	goto L306
L306:
	;
	v3348 = int32(0)
	v3349 = v3309
	goto L308
L308:
	;
	v3430 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3349))))
	v3434 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3349+v290<<(uint(v3229)%32)))))
	v3438 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3349+v259<<(uint(v3229)%32)&int32(-4)))))
	v3442 = (v248*v3430 + v2605 + v245*v3434 + v242*v3438 + v1676) >> (uint(v2601) % 32)
	if base.Ui32(v3442&int32(_a_F_SharpYuvConvertWithOptions_2)) < base.Ui32(int32(256)) {
		goto L310
	} else {
		goto L311
	}
L309:
	;
	v3482 = int32(1)
	v3484 = v3236 + v3482
	if v3484 != v3220 {
		v3236 = v3484
		v3244 = v3244 + l9
		v3246 = v3246 + l11
		v3309 = v3309 + v3226
		goto L306
	} else {
		goto L317
	}
L310:
	;
	v3452 = v3442
	goto L312
L311:
	;
	v3452 = int32(base.Ui32(base.I32_extend16_s(v3442)^int32(-1)) >> (uint(int32(15)) % 32))
	goto L312
L312:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3244+v3348))) = uint8(v3452)
	v3462 = (v247*v3430 + v2605 + v244*v3434 + v241*v3438 + v1675) >> (uint(v2601) % 32)
	if base.Ui32(v3462&int32(_a_F_SharpYuvConvertWithOptions_2)) < base.Ui32(int32(256)) {
		goto L313
	} else {
		goto L314
	}
L313:
	;
	v3472 = v3462
	goto L315
L314:
	;
	v3472 = int32(base.Ui32(base.I32_extend16_s(v3462)^int32(-1)) >> (uint(int32(15)) % 32))
	goto L315
L315:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3246+v3348))) = uint8(v3472)
	v3477 = v3348 + int32(1)
	if v3224 != v3477 {
		v3348 = v3477
		v3349 = v3349 + int32(2)
		goto L308
	} else {
		goto L316
	}
L316:
	;
	goto L309
L317:
	;
	v3758 = v3482
	goto L46
L318:
	;
	v3612 = int32(0)
	goto L320
L319:
	;
	v3758 = v3738
	goto L46
L320:
	;
	v3695 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3568+v3612))))
	v3699 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3559+v3612))))
	v3703 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3575+v3612))))
	v3706 = (v1676 + v2605 + v248*v3695 + v245*v3699 + v242*v3703) >> (uint(v2601) % 32)
	if v184 < base.I32_extend16_s(v3706) {
		goto L322
	} else {
		goto L323
	}
L321:
	;
	v3738 = int32(1)
	v3740 = v3499 + v3738
	if v3740 != v3220 {
		v3499 = v3740
		v3507 = v3507 + l9
		v3509 = v3509 + l11
		v3559 = v3559 + v3226
		v3568 = v3568 + v3226
		v3575 = v3575 + v3226
		goto L318
	} else {
		goto L335
	}
L322:
	;
	v3709 = v184
	goto L324
L323:
	;
	v3709 = v3706
	goto L324
L324:
	;
	if v3706&int32(_a_F_SharpYuvConvertWithOptions_3) != 0 {
		goto L325
	} else {
		goto L326
	}
L325:
	;
	v3712 = int32(0)
	goto L327
L326:
	;
	v3712 = v3709
	goto L327
L327:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3507+v3612))) = uint16(v3712)
	v3722 = (v1675 + v2605 + v247*v3695 + v244*v3699 + v241*v3703) >> (uint(v2601) % 32)
	if v184 < base.I32_extend16_s(v3722) {
		goto L328
	} else {
		goto L329
	}
L328:
	;
	v3725 = v184
	goto L330
L329:
	;
	v3725 = v3722
	goto L330
L330:
	;
	if v3722&int32(_a_F_SharpYuvConvertWithOptions_3) != 0 {
		goto L331
	} else {
		goto L332
	}
L331:
	;
	v3728 = int32(0)
	goto L333
L332:
	;
	v3728 = v3725
	goto L333
L333:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3509+v3612))) = uint16(v3728)
	v3731 = v3612 + int32(2)
	if v3224<<(uint(v3486)%32) != v3731 {
		v3612 = v3731
		goto L320
	} else {
		goto L334
	}
L334:
	;
	goto L321
L335:
	;
	goto L319
L336:
	;
	F_dlfree(m, v314)
	mBase = m.M
	goto L337
L337:
	;
	F_dlfree(m, v285)
	mBase = m.M
	goto L338
L338:
	;
	F_dlfree(m, v315)
	mBase = m.M
	goto L339
L339:
	;
	F_dlfree(m, v299)
	mBase = m.M
	goto L340
L340:
	;
	F_dlfree(m, v323)
	mBase = m.M
	goto L341
L341:
	;
	F_dlfree(m, v271)
	mBase = m.M
	goto L342
L342:
	;
	v3861 = v3758
	goto L1
}
func F_SharpYuvGammaToLinear(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v12 int32
	_ = v12
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v35 int32
	_ = v35
	var v45 int32
	_ = v45
	var v51 int32
	_ = v51
	var v56 float32
	_ = v56
	var v66 float32
	_ = v66
	var v77 float64
	_ = v77
	var v79 float32
	_ = v79
	var v82 float32
	_ = v82
	var v85 float64
	_ = v85
	var v87 float32
	_ = v87
	var v90 float32
	_ = v90
	var v93 float64
	_ = v93
	var v101 float32
	_ = v101
	var v112 float64
	_ = v112
	var v120 float32
	_ = v120
	var v123 float32
	_ = v123
	var v125 float32
	_ = v125
	var v128 float64
	_ = v128
	var v136 float32
	_ = v136
	var v139 float32
	_ = v139
	var v145 float64
	_ = v145
	var v159 float64
	_ = v159
	var v167 float32
	_ = v167
	var v178 float64
	_ = v178
	var v180 float32
	_ = v180
	var v187 float64
	_ = v187
	var v188 float32
	_ = v188
	var v190 float32
	_ = v190
	var v191 float32
	_ = v191
	var v194 float32
	_ = v194
	var v198 float32
	_ = v198
	var v199 float32
	_ = v199
	var v202 float32
	_ = v202
	var v206 float64
	_ = v206
	var v208 float32
	_ = v208
	var v211 float32
	_ = v211
	var v214 float64
	_ = v214
	var v227 float64
	_ = v227
	var v232 float32
	_ = v232
	var v243 int32
	_ = v243
	var v257 float32
	_ = v257
	var v263 float32
	_ = v263
	var v265 int32
	_ = v265
	var v266 float64
	_ = v266
	var v268 float64
	_ = v268
	var v270 float64
	_ = v270
	var v272 float64
	_ = v272
	var v273 float64
	_ = v273
	var v275 float64
	_ = v275
	var v278 float64
	_ = v278
	var v283 float64
	_ = v283
	var v288 int64
	_ = v288
	var v298 int64
	_ = v298
	var v303 float32
	_ = v303
	var v313 float32
	_ = v313
	var v320 float64
	_ = v320
	var v323 float32
	_ = v323
	var v325 float32
	_ = v325
	var v334 float32
	_ = v334
	var v342 int32
	_ = v342
	var v345 int32
	_ = v345
	if l2 != int32(13) {
		v51 = int32(-1)
		v56 = base.F32_div(base.F32_convert_i32_u(l0), base.F32_convert_i32_u(v51<<(uint(l1)%32)^v51))
		switch l2 + v51 {
		case 0, 5, 13, 14:
			if base.F32_lt(v56, float32(0.08124286)) == int32(0) {
				v66 = float32(1)
				if base.F32_lt(v56, v66) == int32(0) {
					v323 = v66
				} else {
					v77 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.09929682)), float32(1.0992968))), float64(2.222222328186035))
					mBase = m.M
					v323 = base.F32_demote_f64(v77)
				}
			} else {
				v323 = base.F32_div(v56, float32(4.5))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		default:
			v323 = float32(0)
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 3:
			v79 = float32(1)
			if base.F32_gt(v56, v79) != 0 {
				v82 = v79
			} else {
				v82 = v56
			}
			v85 = F_pow(m, base.F64_promote_f32(v82), float64(2.200000047683716))
			mBase = m.M
			v323 = base.F32_demote_f64(v85)
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 4:
			v87 = float32(1)
			if base.F32_gt(v56, v87) != 0 {
				v90 = v87
			} else {
				v90 = v56
			}
			v93 = F_pow(m, base.F64_promote_f32(v90), float64(2.799999952316284))
			mBase = m.M
			v323 = base.F32_demote_f64(v93)
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 6:
			if base.F32_lt(v56, float32(0.09128634)) == int32(0) {
				v101 = float32(1)
				if base.F32_lt(v56, v101) == int32(0) {
					v323 = v101
				} else {
					v112 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.1115722)), float32(1.1115721))), float64(2.222222328186035))
					mBase = m.M
					v323 = base.F32_demote_f64(v112)
				}
			} else {
				v323 = base.F32_mul(v56, float32(0.25))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 7:
			v345 = l0
			return v345
		case 8:
			if base.F32_le(v56, float32(0)) == int32(0) {
				v120 = float32(1)
				if base.F32_lt(v56, v120) != 0 {
					v123 = v56
				} else {
					v123 = v120
				}
				v125 = base.F32_add(v123, float32(-1))
				v128 = F_pow(m, float64(10), base.F64_promote_f32(base.F32_add(v125, v125)))
				mBase = m.M
				v323 = base.F32_demote_f64(v128)
			} else {
				v323 = float32(0.005)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 9:
			if base.F32_le(v56, float32(0)) == int32(0) {
				v136 = float32(1)
				if base.F32_lt(v56, v136) != 0 {
					v139 = v56
				} else {
					v139 = v136
				}
				v145 = F_pow(m, float64(10), base.F64_promote_f32(base.F32_mul(base.F32_add(v139, float32(-1)), float32(2.5))))
				mBase = m.M
				v323 = base.F32_demote_f64(v145)
			} else {
				v323 = float32(0.0015811388)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 10:
			if base.F32_lt(v56, float32(0.08124286)) == int32(0) {
				v159 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.09929682)), float32(1.0992968))), float64(2.222222328186035))
				mBase = m.M
				v323 = base.F32_demote_f64(v159)
			} else {
				v323 = base.F32_div(v56, float32(4.5))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 11:
			if base.F32_lt(v56, float32(0.08124286)) == int32(0) {
				v167 = float32(1)
				if base.F32_lt(v56, v167) == int32(0) {
					v323 = v167
				} else {
					v178 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v56, float32(0.09929682)), float32(1.0992968))), float64(2.222222328186035))
					mBase = m.M
					v323 = base.F32_demote_f64(v178)
				}
			} else {
				v323 = base.F32_div(v56, float32(4.5))
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 15:
			v180 = float32(0)
			if base.F32_gt(v56, v180) == int32(0) {
				v323 = v180
			} else {
				v187 = F_pow(m, base.F64_promote_f32(v56), float64(0.012683313339948654))
				mBase = m.M
				v188 = base.F32_demote_f64(v187)
				v190 = base.F32_add(v188, float32(-0.8359375))
				v191 = float32(0)
				if base.F32_gt(v190, v191) != 0 {
					v194 = v190
				} else {
					v194 = v191
				}
				v198 = base.F32_add(base.F32_mul(v188, float32(-18.6875)), float32(18.851562))
				v199 = float32(1.1754944e-38)
				if base.F32_gt(v198, v199) != 0 {
					v202 = v198
				} else {
					v202 = v199
				}
				v206 = F_pow(m, base.F64_promote_f32(base.F32_div(v194, v202)), float64(6.27258825302124))
				mBase = m.M
				v323 = base.F32_demote_f64(v206)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 16:
			v208 = float32(0)
			if base.F32_gt(v56, v208) != 0 {
				v211 = v56
			} else {
				v211 = v208
			}
			v214 = F_pow(m, base.F64_promote_f32(v211), float64(2.5999999046325684))
			mBase = m.M
			v323 = base.F32_div(base.F32_demote_f64(v214), float32(0.9165553))
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		case 17:
			if base.F32_le(v56, float32(0.5)) == int32(0) {
				v232 = base.F32_div(base.F32_add(v56, float32(-0.5599107)), float32(0.17883277))
				v243 = int32(base.Ui32(base.I32_reinterpret_f32(v232))>>(uint(int32(20))%32)) & int32(2047)
				if base.Ui32(v243) < base.Ui32(int32(1067)) {
					v265 = int32(0)
					v266 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[0]))
					v268 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[1]))
					v270 = base.F64_mul(v268, base.F64_promote_f32(v232))
					v272 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[2]))
					v273 = base.F64_add(v270, v272)
					v275 = base.F64_sub(v270, base.F64_sub(v273, v272))
					v278 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[3]))
					v283 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[4]))
					v288 = base.I64_reinterpret_f64(v273)
					v298 = *(*int64)(unsafe.Add(mBase, uint32(base.I32_wrap_i64(v288)&int32(31)<<(uint(int32(3))%32))+uint32(_c_F_SharpYuvGammaToLinear[5])))
					v303 = base.F32_demote_f64(base.F64_mul(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v266, v275), v278), base.F64_mul(v275, v275)), base.F64_add(base.F64_mul(v283, v275), float64(1))), base.F64_reinterpret_i64(v288<<(uint(int64(47))%64)+v298)))
					v313 = v303
				} else {
					if base.F32_eq(v232, math.Float32frombits(uint32(0xff800000))) != 0 {
						v303 = float32(0)
						v313 = v303
					} else {
						if base.Ui32(v243) < base.Ui32(int32(2040)) {
							if base.F32_gt(v232, float32(88.72283)) == int32(0) {
								if base.F32_lt(v232, float32(-103.97208)) == int32(0) {
									v265 = int32(0)
									v266 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[0]))
									v268 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[1]))
									v270 = base.F64_mul(v268, base.F64_promote_f32(v232))
									v272 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[2]))
									v273 = base.F64_add(v270, v272)
									v275 = base.F64_sub(v270, base.F64_sub(v273, v272))
									v278 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[3]))
									v283 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvGammaToLinear[4]))
									v288 = base.I64_reinterpret_f64(v273)
									v298 = *(*int64)(unsafe.Add(mBase, uint32(base.I32_wrap_i64(v288)&int32(31)<<(uint(int32(3))%32))+uint32(_c_F_SharpYuvGammaToLinear[5])))
									v303 = base.F32_demote_f64(base.F64_mul(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v266, v275), v278), base.F64_mul(v275, v275)), base.F64_add(base.F64_mul(v283, v275), float64(1))), base.F64_reinterpret_i64(v288<<(uint(int64(47))%64)+v298)))
									v313 = v303
								} else {
									v263 = F___math_uflowf(m, int32(0))
									mBase = m.M
									v313 = v263
								}
							} else {
								v257 = F___math_oflowf(m, int32(0))
								mBase = m.M
								v313 = v257
							}
						} else {
							v313 = base.F32_add(v232, v232)
						}
					}
				}
				v320 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(v313, float32(0.28466892)), float32(12))), float64(1.2000000476837158))
				mBase = m.M
				v323 = base.F32_demote_f64(v320)
			} else {
				v227 = F_pow(m, base.F64_promote_f32(base.F32_mul(base.F32_mul(v56, v56), float32(0.33333334))), float64(1.2000000476837158))
				mBase = m.M
				v323 = base.F32_demote_f64(v227)
			}
			v325 = base.F32_mul(v323, float32(65535))
			if base.F32_lt(v325, float32(0)) != 0 {
				v334 = base.F32_ceil(base.F32_add(v325, float32(-0.5)))
			} else {
				v334 = base.F32_floor(base.F32_add(v325, float32(0.5)))
			}
			if base.F32_lt(v334, float32(4.2949673e+09))&base.F32_ge(v334, float32(0)) == int32(0) {
				v345 = int32(0)
				return v345
			} else {
				v342 = base.I32_trunc_f32_u(v334)
				return v342
			}
		}
	} else {
		if int32(9) < l1 {
			v23 = m.G1
			v27 = l1 + int32(-10)
			v28 = int32(base.Ui32(l0) >> (uint(v27) % 32))
			v31 = v23 + int32(_a_F_SharpYuvGammaToLinear_0) + v28<<(uint(int32(2))%32)
			v34 = *(*int32)(unsafe.Add(mBase, uint32(v31+int32(4))))
			v35 = *(*int32)(unsafe.Add(mBase, uint32(v31)))
			if v27 != 0 {
				v45 = int32(1) << (uint(l1+int32(-11)) % 32)
			} else {
				v45 = int32(0)
			}
			return int32(base.Ui32((v34-v35)*(l0-v28<<(uint(v27)%32))+v45)>>(uint(v27)%32)) + v35
		} else {
			v12 = m.G1
			v21 = *(*int32)(unsafe.Add(mBase, uint32(v12+int32(_a_F_SharpYuvGammaToLinear_0)+l0<<(uint(int32(10)-l1)%32)<<(uint(int32(2))%32))))
			return v21
		}
	}
}
func F_SharpYuvGetConversionMatrix(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	if base.Ui32(int32(5)) < base.Ui32(l0) {
		v13 = int32(0)
	} else {
		v6 = m.G1
		v12 = *(*int32)(unsafe.Add(mBase, uint32(v6+int32(_a_F_SharpYuvGetConversionMatrix_0)+l0<<(uint(int32(2))%32))))
		v13 = v12
	}
	return v13
}
func F_SharpYuvInit(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v4 int32
	_ = v4
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v8 int32
	_ = v8
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	v2 = m.G83
	if l0 != v2 {
		v6 = m.G83
		*(*int32)(unsafe.Add(mBase, uint32(v6))) = l0
		v8 = l0
	} else {
		v4 = m.G83
		v5 = *(*int32)(unsafe.Add(mBase, uint32(v4)))
		v8 = v5
	}
	v9 = m.G1
	v12 = *(*int32)(unsafe.Add(mBase, uint32(v9)+uint32(_c_F_SharpYuvInit[0])))
	if v12 == v8 {
	} else {
		v14 = m.G1
		v17 = m.G2
		v18 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v18)+uint32(_c_F_SharpYuvInit[1]))) = v17 + int32(178)
		*(*int32)(unsafe.Add(mBase, uint32(v18)+uint32(_c_F_SharpYuvInit[2]))) = v17 + int32(179)
		*(*int32)(unsafe.Add(mBase, uint32(v18)+uint32(_c_F_SharpYuvInit[3]))) = v17 + int32(180)
		F_SharpYuvInitGammaTables(m)
		mBase = m.M
		v37 = m.G83
		v38 = *(*int32)(unsafe.Add(mBase, uint32(v37)))
		*(*int32)(unsafe.Add(mBase, uint32(v14)+uint32(_c_F_SharpYuvInit[0]))) = v38
	}
	return
}
func F_SharpYuvInitDsp(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v4 int32
	_ = v4
	v3 = m.G2
	v4 = m.G1
	*(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_SharpYuvInitDsp[0]))) = v3 + int32(178)
	*(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_SharpYuvInitDsp[1]))) = v3 + int32(179)
	*(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_SharpYuvInitDsp[2]))) = v3 + int32(180)
	return
}
func F_SharpYuvInitGammaTables(m *base.Module) {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v9 int32
	_ = v9
	var v12 float64
	_ = v12
	var v13 int32
	_ = v13
	var v18 float64
	_ = v18
	var v30 float64
	_ = v30
	var v31 float64
	_ = v31
	var v32 int32
	_ = v32
	var v41 float64
	_ = v41
	var v49 int32
	_ = v49
	var v51 int32
	_ = v51
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v64 float64
	_ = v64
	var v65 int32
	_ = v65
	var v70 float64
	_ = v70
	var v78 float64
	_ = v78
	var v83 float64
	_ = v83
	var v84 int32
	_ = v84
	var v93 float64
	_ = v93
	var v101 int32
	_ = v101
	var v103 int32
	_ = v103
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v116 int32
	_ = v116
	v6 = m.G1
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v6)+uint32(_c_F_SharpYuvInitGammaTables[0])))
	if v9 != 0 {
	} else {
		v12 = float64(0)
		v13 = int32(-4100)
		for {
			v18 = base.F64_mul(v12, float64(0.0009765625))
			if base.F64_le(v18, float64(0.08124285829863151)) == int32(0) {
				v30 = F_pow(m, base.F64_mul(base.F64_add(v18, float64(0.09929682680944)), float64(0.909672415686275)), float64(2.2222222222222223))
				mBase = m.M
				v31 = v30
			} else {
				v31 = base.F64_div(v18, float64(4.5))
			}
			v32 = m.G1
			v41 = base.F64_add(base.F64_mul(v31, float64(65536)), float64(0.5))
			if base.F64_lt(v41, float64(4.294967296e+09))&base.F64_ge(v41, float64(0)) == int32(0) {
				v51 = int32(0)
			} else {
				v49 = base.I32_trunc_f64_u(v41)
				v51 = v49
			}
			*(*int32)(unsafe.Add(mBase, uint32(v32+int32(_a_F_SharpYuvInitGammaTables_0)+v13)+uint32(_c_F_SharpYuvInitGammaTables[1]))) = v51
			v56 = v13 + int32(4)
			if v56 != 0 {
				v12 = base.F64_add(v12, float64(1))
				v13 = v56
				continue
			} else {
				break
			}
			break
		}
		v57 = m.G1
		v60 = *(*int32)(unsafe.Add(mBase, uint32(v57)+uint32(_c_F_SharpYuvInitGammaTables[2])))
		*(*int32)(unsafe.Add(mBase, uint32(v57)+uint32(_c_F_SharpYuvInitGammaTables[3]))) = v60
		v64 = float64(0)
		v65 = int32(-2052)
		for {
			v70 = base.F64_mul(v64, float64(0.001953125))
			if base.F64_le(v70, float64(0.018053968510807)) == int32(0) {
				v78 = F_pow(m, v70, float64(0.44999999999999996))
				mBase = m.M
				v83 = base.F64_add(base.F64_mul(v78, float64(1.09929682680944)), float64(-0.09929682680944))
			} else {
				v83 = base.F64_mul(v70, float64(4.5))
			}
			v84 = m.G1
			v93 = base.F64_add(base.F64_mul(v83, float64(65536)), float64(0.5))
			if base.F64_lt(v93, float64(4.294967296e+09))&base.F64_ge(v93, float64(0)) == int32(0) {
				v103 = int32(0)
			} else {
				v101 = base.I32_trunc_f64_u(v93)
				v103 = v101
			}
			*(*int32)(unsafe.Add(mBase, uint32(v84+int32(_a_F_SharpYuvInitGammaTables_1)+v65+int32(2052)))) = v103
			v108 = v65 + int32(4)
			if v108 != 0 {
				v64 = base.F64_add(v64, float64(1))
				v65 = v108
				continue
			} else {
				break
			}
			break
		}
		v109 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v109)+uint32(_c_F_SharpYuvInitGammaTables[0]))) = int32(1)
		v116 = *(*int32)(unsafe.Add(mBase, uint32(v109)+uint32(_c_F_SharpYuvInitGammaTables[4])))
		*(*int32)(unsafe.Add(mBase, uint32(v109)+uint32(_c_F_SharpYuvInitGammaTables[5]))) = v116
	}
	return
}
func F_SharpYuvLinearToGamma(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v11 int32
	_ = v11
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v30 int32
	_ = v30
	var v31 int32
	_ = v31
	var v34 int32
	_ = v34
	var v46 float32
	_ = v46
	var v56 float32
	_ = v56
	var v63 float64
	_ = v63
	var v69 float32
	_ = v69
	var v72 float32
	_ = v72
	var v75 float64
	_ = v75
	var v77 float32
	_ = v77
	var v80 float32
	_ = v80
	var v83 float64
	_ = v83
	var v91 float32
	_ = v91
	var v98 float64
	_ = v98
	var v107 float32
	_ = v107
	var v110 float32
	_ = v110
	var v111 float64
	_ = v111
	var v122 int64
	_ = v122
	var v137 int32
	_ = v137
	var v139 int64
	_ = v139
	var v148 int64
	_ = v148
	var v153 int64
	_ = v153
	var v154 int32
	_ = v154
	var v156 int32
	_ = v156
	var v158 int32
	_ = v158
	var v162 float64
	_ = v162
	var v164 float64
	_ = v164
	var v177 float64
	_ = v177
	var v180 float64
	_ = v180
	var v185 float64
	_ = v185
	var v186 float64
	_ = v186
	var v187 float64
	_ = v187
	var v188 float64
	_ = v188
	var v193 float64
	_ = v193
	var v194 float64
	_ = v194
	var v195 float64
	_ = v195
	var v220 float64
	_ = v220
	var v232 float64
	_ = v232
	var v254 float64
	_ = v254
	var v263 float32
	_ = v263
	var v266 float32
	_ = v266
	var v267 float64
	_ = v267
	var v278 int64
	_ = v278
	var v293 int32
	_ = v293
	var v295 int64
	_ = v295
	var v304 int64
	_ = v304
	var v309 int64
	_ = v309
	var v310 int32
	_ = v310
	var v312 int32
	_ = v312
	var v314 int32
	_ = v314
	var v318 float64
	_ = v318
	var v320 float64
	_ = v320
	var v333 float64
	_ = v333
	var v336 float64
	_ = v336
	var v341 float64
	_ = v341
	var v342 float64
	_ = v342
	var v343 float64
	_ = v343
	var v344 float64
	_ = v344
	var v349 float64
	_ = v349
	var v350 float64
	_ = v350
	var v351 float64
	_ = v351
	var v376 float64
	_ = v376
	var v388 float64
	_ = v388
	var v410 float64
	_ = v410
	var v424 float64
	_ = v424
	var v436 float32
	_ = v436
	var v443 float64
	_ = v443
	var v449 float32
	_ = v449
	var v456 float64
	_ = v456
	var v457 float32
	_ = v457
	var v469 float64
	_ = v469
	var v471 float32
	_ = v471
	var v474 float32
	_ = v474
	var v479 float64
	_ = v479
	var v481 float32
	_ = v481
	var v484 float64
	_ = v484
	var v485 float32
	_ = v485
	var v498 float32
	_ = v498
	var v503 int32
	_ = v503
	var v509 int32
	_ = v509
	var v511 float32
	_ = v511
	var v518 float32
	_ = v518
	var v524 int32
	_ = v524
	var v526 int32
	_ = v526
	var v527 float64
	_ = v527
	var v529 int32
	_ = v529
	var v538 int32
	_ = v538
	var v541 float64
	_ = v541
	var v544 float64
	_ = v544
	var v545 float64
	_ = v545
	var v548 float64
	_ = v548
	var v551 float64
	_ = v551
	var v559 float64
	_ = v559
	var v563 float64
	_ = v563
	var v568 float32
	_ = v568
	var v578 float32
	_ = v578
	var v584 float32
	_ = v584
	var v585 int32
	_ = v585
	var v590 float32
	_ = v590
	var v599 float32
	_ = v599
	var v607 int32
	_ = v607
	var v609 int32
	_ = v609
	if l2 != int32(13) {
		v46 = base.F32_div(base.F32_convert_i32_u(l0), float32(65535))
		switch l2 + int32(-1) {
		case 0, 5, 13, 14:
			if base.F32_lt(v46, float32(0.01805397)) == int32(0) {
				v56 = float32(1)
				if base.F32_lt(v46, v56) == int32(0) {
					v584 = v56
				} else {
					v63 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
					mBase = m.M
					v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v63), float32(1.0992968)), float32(-0.09929682))
				}
			} else {
				v584 = base.F32_mul(v46, float32(4.5))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		default:
			v584 = float32(0)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 3:
			v69 = float32(1)
			if base.F32_gt(v46, v69) != 0 {
				v72 = v69
			} else {
				v72 = v46
			}
			v75 = F_pow(m, base.F64_promote_f32(v72), float64(0.45454543828964233))
			mBase = m.M
			v584 = base.F32_demote_f64(v75)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 4:
			v77 = float32(1)
			if base.F32_gt(v46, v77) != 0 {
				v80 = v77
			} else {
				v80 = v46
			}
			v83 = F_pow(m, base.F64_promote_f32(v80), float64(0.3571428656578064))
			mBase = m.M
			v584 = base.F32_demote_f64(v83)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 6:
			if base.F32_lt(v46, float32(0.022821585)) == int32(0) {
				v91 = float32(1)
				if base.F32_lt(v46, v91) == int32(0) {
					v584 = v91
				} else {
					v98 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
					mBase = m.M
					v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v98), float32(1.1115721)), float32(-0.1115722))
				}
			} else {
				v584 = base.F32_mul(v46, float32(4))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 7:
			v609 = l0
		case 8:
			if base.F32_lt(v46, float32(0.01)) != 0 {
				v584 = float32(0)
			} else {
				v107 = float32(1)
				if base.F32_lt(v46, v107) != 0 {
					v110 = v46
				} else {
					v110 = v107
				}
				v111 = base.F64_promote_f32(v110)
				v122 = base.I64_reinterpret_f64(v111)
				if int64(4503599627370495) < v122 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v122) {
						v232 = v111
						v254 = v232
					} else {
						v137 = int32(-1023)
						v139 = int64(base.Ui64(v122) >> (uint(int64(32)) % 64))
						if v139 == int64(1072693248) {
							if base.I32_wrap_i64(v122) != 0 {
								v153 = v122
								v154 = v137
								v156 = int32(1072693248)
								v158 = v156 + int32(614242)
								v162 = base.F64_convert_i32_s(v154 + int32(base.Ui32(v158)>>(uint(int32(20))%32)))
								v164 = base.F64_mul(v162, float64(0.30102999566361177))
								v177 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v158&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v153&int64(4294967295)), float64(-1))
								v180 = base.F64_mul(v177, base.F64_mul(v177, float64(0.5)))
								v185 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v177, v180)) & int64(-4294967296))
								v186 = float64(0.4342944818781689)
								v187 = base.F64_mul(v185, v186)
								v188 = base.F64_add(v164, v187)
								v193 = base.F64_div(v177, base.F64_add(v177, float64(2)))
								v194 = base.F64_mul(v193, v193)
								v195 = base.F64_mul(v194, v194)
								v220 = base.F64_add(base.F64_mul(v193, base.F64_add(v180, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v194, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v177, v185), v180))
								v232 = base.F64_add(v188, base.F64_add(base.F64_add(v187, base.F64_sub(v164, v188)), base.F64_add(base.F64_mul(v220, v186), base.F64_add(base.F64_mul(v162, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v220, v185), float64(2.5082946711645275e-11))))))
								v254 = v232
							} else {
								v254 = float64(0)
							}
						} else {
							v153 = v122
							v154 = v137
							v156 = base.I32_wrap_i64(v139)
							v158 = v156 + int32(614242)
							v162 = base.F64_convert_i32_s(v154 + int32(base.Ui32(v158)>>(uint(int32(20))%32)))
							v164 = base.F64_mul(v162, float64(0.30102999566361177))
							v177 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v158&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v153&int64(4294967295)), float64(-1))
							v180 = base.F64_mul(v177, base.F64_mul(v177, float64(0.5)))
							v185 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v177, v180)) & int64(-4294967296))
							v186 = float64(0.4342944818781689)
							v187 = base.F64_mul(v185, v186)
							v188 = base.F64_add(v164, v187)
							v193 = base.F64_div(v177, base.F64_add(v177, float64(2)))
							v194 = base.F64_mul(v193, v193)
							v195 = base.F64_mul(v194, v194)
							v220 = base.F64_add(base.F64_mul(v193, base.F64_add(v180, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v194, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v177, v185), v180))
							v232 = base.F64_add(v188, base.F64_add(base.F64_add(v187, base.F64_sub(v164, v188)), base.F64_add(base.F64_mul(v220, v186), base.F64_add(base.F64_mul(v162, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v220, v185), float64(2.5082946711645275e-11))))))
							v254 = v232
						}
					}
				} else {
					if base.F64_ne(v111, float64(0)) != 0 {
						if int64(-1) < v122 {
							v148 = base.I64_reinterpret_f64(base.F64_mul(v111, float64(1.8014398509481984e+16)))
							v153 = v148
							v154 = int32(-1077)
							v156 = base.I32_wrap_i64(int64(base.Ui64(v148) >> (uint(int64(32)) % 64)))
							v158 = v156 + int32(614242)
							v162 = base.F64_convert_i32_s(v154 + int32(base.Ui32(v158)>>(uint(int32(20))%32)))
							v164 = base.F64_mul(v162, float64(0.30102999566361177))
							v177 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v158&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v153&int64(4294967295)), float64(-1))
							v180 = base.F64_mul(v177, base.F64_mul(v177, float64(0.5)))
							v185 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v177, v180)) & int64(-4294967296))
							v186 = float64(0.4342944818781689)
							v187 = base.F64_mul(v185, v186)
							v188 = base.F64_add(v164, v187)
							v193 = base.F64_div(v177, base.F64_add(v177, float64(2)))
							v194 = base.F64_mul(v193, v193)
							v195 = base.F64_mul(v194, v194)
							v220 = base.F64_add(base.F64_mul(v193, base.F64_add(v180, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v194, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, base.F64_add(base.F64_mul(v195, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v177, v185), v180))
							v232 = base.F64_add(v188, base.F64_add(base.F64_add(v187, base.F64_sub(v164, v188)), base.F64_add(base.F64_mul(v220, v186), base.F64_add(base.F64_mul(v162, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v220, v185), float64(2.5082946711645275e-11))))))
							v254 = v232
						} else {
							v254 = base.F64_div(base.F64_sub(v111, v111), float64(0))
						}
					} else {
						v254 = base.F64_div(float64(-1), base.F64_mul(v111, v111))
					}
				}
				v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v254), float32(0.5)), float32(1))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 9:
			if base.F32_lt(v46, float32(0.0031622776)) != 0 {
				v584 = float32(0)
			} else {
				v263 = float32(1)
				if base.F32_lt(v46, v263) != 0 {
					v266 = v46
				} else {
					v266 = v263
				}
				v267 = base.F64_promote_f32(v266)
				v278 = base.I64_reinterpret_f64(v267)
				if int64(4503599627370495) < v278 {
					if base.Ui64(int64(9218868437227405311)) < base.Ui64(v278) {
						v388 = v267
						v410 = v388
					} else {
						v293 = int32(-1023)
						v295 = int64(base.Ui64(v278) >> (uint(int64(32)) % 64))
						if v295 == int64(1072693248) {
							if base.I32_wrap_i64(v278) != 0 {
								v309 = v278
								v310 = v293
								v312 = int32(1072693248)
								v314 = v312 + int32(614242)
								v318 = base.F64_convert_i32_s(v310 + int32(base.Ui32(v314)>>(uint(int32(20))%32)))
								v320 = base.F64_mul(v318, float64(0.30102999566361177))
								v333 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v314&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v309&int64(4294967295)), float64(-1))
								v336 = base.F64_mul(v333, base.F64_mul(v333, float64(0.5)))
								v341 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v333, v336)) & int64(-4294967296))
								v342 = float64(0.4342944818781689)
								v343 = base.F64_mul(v341, v342)
								v344 = base.F64_add(v320, v343)
								v349 = base.F64_div(v333, base.F64_add(v333, float64(2)))
								v350 = base.F64_mul(v349, v349)
								v351 = base.F64_mul(v350, v350)
								v376 = base.F64_add(base.F64_mul(v349, base.F64_add(v336, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v350, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v333, v341), v336))
								v388 = base.F64_add(v344, base.F64_add(base.F64_add(v343, base.F64_sub(v320, v344)), base.F64_add(base.F64_mul(v376, v342), base.F64_add(base.F64_mul(v318, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v376, v341), float64(2.5082946711645275e-11))))))
								v410 = v388
							} else {
								v410 = float64(0)
							}
						} else {
							v309 = v278
							v310 = v293
							v312 = base.I32_wrap_i64(v295)
							v314 = v312 + int32(614242)
							v318 = base.F64_convert_i32_s(v310 + int32(base.Ui32(v314)>>(uint(int32(20))%32)))
							v320 = base.F64_mul(v318, float64(0.30102999566361177))
							v333 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v314&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v309&int64(4294967295)), float64(-1))
							v336 = base.F64_mul(v333, base.F64_mul(v333, float64(0.5)))
							v341 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v333, v336)) & int64(-4294967296))
							v342 = float64(0.4342944818781689)
							v343 = base.F64_mul(v341, v342)
							v344 = base.F64_add(v320, v343)
							v349 = base.F64_div(v333, base.F64_add(v333, float64(2)))
							v350 = base.F64_mul(v349, v349)
							v351 = base.F64_mul(v350, v350)
							v376 = base.F64_add(base.F64_mul(v349, base.F64_add(v336, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v350, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v333, v341), v336))
							v388 = base.F64_add(v344, base.F64_add(base.F64_add(v343, base.F64_sub(v320, v344)), base.F64_add(base.F64_mul(v376, v342), base.F64_add(base.F64_mul(v318, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v376, v341), float64(2.5082946711645275e-11))))))
							v410 = v388
						}
					}
				} else {
					if base.F64_ne(v267, float64(0)) != 0 {
						if int64(-1) < v278 {
							v304 = base.I64_reinterpret_f64(base.F64_mul(v267, float64(1.8014398509481984e+16)))
							v309 = v304
							v310 = int32(-1077)
							v312 = base.I32_wrap_i64(int64(base.Ui64(v304) >> (uint(int64(32)) % 64)))
							v314 = v312 + int32(614242)
							v318 = base.F64_convert_i32_s(v310 + int32(base.Ui32(v314)>>(uint(int32(20))%32)))
							v320 = base.F64_mul(v318, float64(0.30102999566361177))
							v333 = base.F64_add(base.F64_reinterpret_i64(base.I64_extend_i32_u(v314&int32(1048575)+int32(1072079006))<<(uint(int64(32))%64)|v309&int64(4294967295)), float64(-1))
							v336 = base.F64_mul(v333, base.F64_mul(v333, float64(0.5)))
							v341 = base.F64_reinterpret_i64(base.I64_reinterpret_f64(base.F64_sub(v333, v336)) & int64(-4294967296))
							v342 = float64(0.4342944818781689)
							v343 = base.F64_mul(v341, v342)
							v344 = base.F64_add(v320, v343)
							v349 = base.F64_div(v333, base.F64_add(v333, float64(2)))
							v350 = base.F64_mul(v349, v349)
							v351 = base.F64_mul(v350, v350)
							v376 = base.F64_add(base.F64_mul(v349, base.F64_add(v336, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.15313837699209373)), float64(0.22222198432149784))), float64(0.3999999999940942))), base.F64_mul(v350, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, base.F64_add(base.F64_mul(v351, float64(0.14798198605116586)), float64(0.1818357216161805))), float64(0.2857142874366239))), float64(0.6666666666666735)))))), base.F64_sub(base.F64_sub(v333, v341), v336))
							v388 = base.F64_add(v344, base.F64_add(base.F64_add(v343, base.F64_sub(v320, v344)), base.F64_add(base.F64_mul(v376, v342), base.F64_add(base.F64_mul(v318, float64(3.694239077158931e-13)), base.F64_mul(base.F64_add(v376, v341), float64(2.5082946711645275e-11))))))
							v410 = v388
						} else {
							v410 = base.F64_div(base.F64_sub(v267, v267), float64(0))
						}
					} else {
						v410 = base.F64_div(float64(-1), base.F64_mul(v267, v267))
					}
				}
				v584 = base.F32_add(base.F32_div(base.F32_demote_f64(v410), float32(2.5)), float32(1))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 10:
			if base.F32_lt(v46, float32(0.01805397)) == int32(0) {
				v424 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
				mBase = m.M
				v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v424), float32(1.0992968)), float32(-0.09929682))
			} else {
				v584 = base.F32_mul(v46, float32(4.5))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 11:
			if base.F32_lt(v46, float32(0.01805397)) == int32(0) {
				v436 = float32(1)
				if base.F32_lt(v46, v436) == int32(0) {
					v584 = v436
				} else {
					v443 = F_pow(m, base.F64_promote_f32(v46), float64(0.44999998807907104))
					mBase = m.M
					v584 = base.F32_add(base.F32_mul(base.F32_demote_f64(v443), float32(1.0992968)), float32(-0.09929682))
				}
			} else {
				v584 = base.F32_mul(v46, float32(4.5))
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 15:
			v449 = float32(0)
			if base.F32_gt(v46, v449) == int32(0) {
				v584 = v449
			} else {
				v456 = F_pow(m, base.F64_promote_f32(v46), float64(0.159423828125))
				mBase = m.M
				v457 = base.F32_demote_f64(v456)
				v469 = F_pow(m, base.F64_promote_f32(base.F32_div(base.F32_add(base.F32_mul(v457, float32(18.851562)), float32(0.8359375)), base.F32_add(base.F32_mul(v457, float32(18.6875)), float32(1)))), float64(78.84375))
				mBase = m.M
				v584 = base.F32_demote_f64(v469)
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 16:
			v471 = float32(0)
			if base.F32_gt(v46, v471) != 0 {
				v474 = v46
			} else {
				v474 = v471
			}
			v479 = F_pow(m, base.F64_promote_f32(base.F32_mul(v474, float32(0.9165553))), float64(0.38461539149284363))
			mBase = m.M
			v584 = base.F32_demote_f64(v479)
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		case 17:
			v481 = float32(0)
			v484 = F_pow(m, base.F64_promote_f32(v46), float64(0.8333333134651184))
			mBase = m.M
			v485 = base.F32_demote_f64(v484)
			if base.F32_lt(v485, v481) != 0 {
				v584 = v481
			} else {
				if base.F32_le(v485, float32(0.083333336)) == int32(0) {
					v498 = base.F32_add(base.F32_mul(v485, float32(12)), float32(-0.28466892))
					v503 = base.I32_reinterpret_f32(v498)
					if base.Ui32(int32(-2130706433)) < base.Ui32(v503+int32(-2139095040)) {
						v524 = v503
						v526 = int32(0)
						v527 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[0]))
						v529 = v524 + int32(-1060306944)
						v538 = int32(base.Ui32(v529)>>(uint(int32(15))%32)) & int32(240)
						v541 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[1])))
						v544 = base.F64_add(base.F64_mul(base.F64_promote_f32(base.F32_reinterpret_i32(v524-v529&int32(-8388608))), v541), float64(-1))
						v545 = base.F64_mul(v544, v544)
						v548 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[2]))
						v551 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[3]))
						v559 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[4]))
						v563 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[5])))
						v568 = base.F32_demote_f64(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v527, v545), base.F64_add(base.F64_mul(v548, v544), v551)), v545), base.F64_add(base.F64_add(base.F64_mul(base.F64_convert_i32_s(v529>>(uint(int32(23))%32)), v559), v563), v544)))
						v578 = v568
					} else {
						v509 = v503 << (uint(int32(1)) % 32)
						if v509 != 0 {
							if v503 == int32(2139095040) {
								v568 = v498
								v578 = v568
							} else {
								if v503 < int32(0) {
									v518 = F___math_invalidf(m, v498)
									mBase = m.M
									v578 = v518
								} else {
									if base.Ui32(v509) < base.Ui32(int32(-16777216)) {
										v524 = base.I32_reinterpret_f32(base.F32_mul(v498, float32(8.388608e+06))) + int32(-192937984)
										v526 = int32(0)
										v527 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[0]))
										v529 = v524 + int32(-1060306944)
										v538 = int32(base.Ui32(v529)>>(uint(int32(15))%32)) & int32(240)
										v541 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[1])))
										v544 = base.F64_add(base.F64_mul(base.F64_promote_f32(base.F32_reinterpret_i32(v524-v529&int32(-8388608))), v541), float64(-1))
										v545 = base.F64_mul(v544, v544)
										v548 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[2]))
										v551 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[3]))
										v559 = *(*float64)(unsafe.Add(mBase, _c_F_SharpYuvLinearToGamma[4]))
										v563 = *(*float64)(unsafe.Add(mBase, uint32(v538)+uint32(_c_F_SharpYuvLinearToGamma[5])))
										v568 = base.F32_demote_f64(base.F64_add(base.F64_mul(base.F64_add(base.F64_mul(v527, v545), base.F64_add(base.F64_mul(v548, v544), v551)), v545), base.F64_add(base.F64_add(base.F64_mul(base.F64_convert_i32_s(v529>>(uint(int32(23))%32)), v559), v563), v544)))
										v578 = v568
									} else {
										v518 = F___math_invalidf(m, v498)
										mBase = m.M
										v578 = v518
									}
								}
							}
						} else {
							v511 = F___math_divzerof(m, int32(1))
							mBase = m.M
							v578 = v511
						}
					}
					v584 = base.F32_add(base.F32_mul(v578, float32(0.17883277)), float32(0.5599107))
				} else {
					v584 = base.F32_sqrt(base.F32_mul(v485, float32(3)))
				}
			}
			v585 = int32(-1)
			v590 = base.F32_mul(v584, base.F32_convert_i32_u(v585<<(uint(l1)%32)^v585))
			if base.F32_lt(v590, float32(0)) != 0 {
				v599 = base.F32_ceil(base.F32_add(v590, float32(-0.5)))
			} else {
				v599 = base.F32_floor(base.F32_add(v590, float32(0.5)))
			}
			if base.F32_lt(v599, float32(4.2949673e+09))&base.F32_ge(v599, float32(0)) == int32(0) {
				v609 = int32(0)
			} else {
				v607 = base.I32_trunc_f32_u(v599)
				v609 = v607
			}
		}
	} else {
		v11 = m.G1
		v18 = v11 + int32(_a_F_SharpYuvLinearToGamma_0) + l0>>(uint(int32(7))%32)<<(uint(int32(2))%32)
		v21 = *(*int32)(unsafe.Add(mBase, uint32(v18+int32(4))))
		v22 = int32(16)
		v23 = v22 - l1
		v26 = l1 + int32(-16)
		v29 = base.B2i32(l1 < v22)
		if l1 < v22 {
			v30 = v21 >> (uint(v23) % 32)
		} else {
			v30 = v21 << (uint(v26) % 32)
		}
		v31 = *(*int32)(unsafe.Add(mBase, uint32(v18)))
		if l1 < v22 {
			v34 = v31 >> (uint(v23) % 32)
		} else {
			v34 = v31 << (uint(v26) % 32)
		}
		v609 = int32(base.Ui32((v30-v34)*(l0&int32(127))+int32(64))>>(uint(int32(7))%32)) + v34
	}
	return v609 & int32(_a_F_SharpYuvLinearToGamma_1)
}

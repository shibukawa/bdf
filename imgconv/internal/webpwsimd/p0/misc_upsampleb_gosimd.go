//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_UpsampleBgrLinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v210 int32
	_ = v210
	var v212 int32
	_ = v212
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v226 int32
	_ = v226
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v261 int32
	_ = v261
	var v263 int32
	_ = v263
	var v270 int32
	_ = v270
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v291 int32
	_ = v291
	var v294 int32
	_ = v294
	var v298 int32
	_ = v298
	var v305 int32
	_ = v305
	var v307 int32
	_ = v307
	var v314 int32
	_ = v314
	var v317 int32
	_ = v317
	var v325 int32
	_ = v325
	var v327 int32
	_ = v327
	var v330 int32
	_ = v330
	var v332 int32
	_ = v332
	var v339 int32
	_ = v339
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v354 int32
	_ = v354
	var v355 int32
	_ = v355
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v378 base.V128
	_ = v378
	var v379 int32
	_ = v379
	var v381 base.V128
	_ = v381
	var v382 int32
	_ = v382
	var v383 base.V128
	_ = v383
	var v384 base.V128
	_ = v384
	var v386 base.V128
	_ = v386
	var v387 base.V128
	_ = v387
	var v389 base.V128
	_ = v389
	var v390 base.V128
	_ = v390
	var v392 base.V128
	_ = v392
	var v394 base.V128
	_ = v394
	var v396 base.V128
	_ = v396
	var v402 base.V128
	_ = v402
	var v403 base.V128
	_ = v403
	var v409 base.V128
	_ = v409
	var v410 base.V128
	_ = v410
	var v411 base.V128
	_ = v411
	var v412 base.V128
	_ = v412
	var v413 int32
	_ = v413
	var v415 base.V128
	_ = v415
	var v416 base.V128
	_ = v416
	var v417 int32
	_ = v417
	var v419 base.V128
	_ = v419
	var v420 base.V128
	_ = v420
	var v422 base.V128
	_ = v422
	var v423 int32
	_ = v423
	var v426 base.V128
	_ = v426
	var v429 int32
	_ = v429
	var v431 base.V128
	_ = v431
	var v432 int32
	_ = v432
	var v434 base.V128
	_ = v434
	var v436 base.V128
	_ = v436
	var v437 base.V128
	_ = v437
	var v439 base.V128
	_ = v439
	var v440 base.V128
	_ = v440
	var v442 base.V128
	_ = v442
	var v443 base.V128
	_ = v443
	var v445 base.V128
	_ = v445
	var v448 base.V128
	_ = v448
	var v454 base.V128
	_ = v454
	var v455 base.V128
	_ = v455
	var v461 base.V128
	_ = v461
	var v462 base.V128
	_ = v462
	var v464 base.V128
	_ = v464
	var v468 base.V128
	_ = v468
	var v471 base.V128
	_ = v471
	var v472 base.V128
	_ = v472
	var v474 base.V128
	_ = v474
	var v475 int32
	_ = v475
	var v478 base.V128
	_ = v478
	var v479 int32
	_ = v479
	var v481 int32
	_ = v481
	var v482 base.V128
	_ = v482
	var v515 base.V128
	_ = v515
	var v517 base.V128
	_ = v517
	var v519 base.V128
	_ = v519
	var v523 base.V128
	_ = v523
	var v524 base.V128
	_ = v524
	var v526 base.V128
	_ = v526
	var v528 base.V128
	_ = v528
	var v529 base.V128
	_ = v529
	var v530 base.V128
	_ = v530
	var v532 base.V128
	_ = v532
	var v537 base.V128
	_ = v537
	var v539 int32
	_ = v539
	var v541 int32
	_ = v541
	var v542 base.V128
	_ = v542
	var v544 base.V128
	_ = v544
	var v550 base.V128
	_ = v550
	var v552 base.V128
	_ = v552
	var v554 base.V128
	_ = v554
	var v555 base.V128
	_ = v555
	var v557 base.V128
	_ = v557
	var v565 base.V128
	_ = v565
	var v566 base.V128
	_ = v566
	var v569 base.V128
	_ = v569
	var v571 base.V128
	_ = v571
	var v577 base.V128
	_ = v577
	var v579 base.V128
	_ = v579
	var v581 base.V128
	_ = v581
	var v582 base.V128
	_ = v582
	var v584 base.V128
	_ = v584
	var v592 int32
	_ = v592
	var v593 base.V128
	_ = v593
	var v595 base.V128
	_ = v595
	var v601 base.V128
	_ = v601
	var v603 base.V128
	_ = v603
	var v605 base.V128
	_ = v605
	var v606 base.V128
	_ = v606
	var v608 base.V128
	_ = v608
	var v616 base.V128
	_ = v616
	var v618 base.V128
	_ = v618
	var v622 base.V128
	_ = v622
	var v624 base.V128
	_ = v624
	var v625 base.V128
	_ = v625
	var v626 base.V128
	_ = v626
	var v628 base.V128
	_ = v628
	var v633 base.V128
	_ = v633
	var v638 base.V128
	_ = v638
	var v640 base.V128
	_ = v640
	var v641 base.V128
	_ = v641
	var v643 base.V128
	_ = v643
	var v651 base.V128
	_ = v651
	var v655 base.V128
	_ = v655
	var v657 base.V128
	_ = v657
	var v658 base.V128
	_ = v658
	var v660 base.V128
	_ = v660
	var v669 base.V128
	_ = v669
	var v671 base.V128
	_ = v671
	var v672 base.V128
	_ = v672
	var v674 base.V128
	_ = v674
	var v682 base.V128
	_ = v682
	var v685 base.V128
	_ = v685
	var v688 base.V128
	_ = v688
	var v690 base.V128
	_ = v690
	var v695 base.V128
	_ = v695
	var v702 base.V128
	_ = v702
	var v719 base.V128
	_ = v719
	var v748 base.V128
	_ = v748
	var v751 base.V128
	_ = v751
	var v758 base.V128
	_ = v758
	var v761 base.V128
	_ = v761
	var v763 base.V128
	_ = v763
	var v768 base.V128
	_ = v768
	var v772 base.V128
	_ = v772
	var v774 base.V128
	_ = v774
	var v779 base.V128
	_ = v779
	var v782 base.V128
	_ = v782
	var v785 base.V128
	_ = v785
	var v790 base.V128
	_ = v790
	var v797 base.V128
	_ = v797
	var v800 base.V128
	_ = v800
	var v807 base.V128
	_ = v807
	var v810 base.V128
	_ = v810
	var v813 base.V128
	_ = v813
	var v818 base.V128
	_ = v818
	var v823 base.V128
	_ = v823
	var v828 base.V128
	_ = v828
	var v831 base.V128
	_ = v831
	var v834 base.V128
	_ = v834
	var v839 base.V128
	_ = v839
	var v844 base.V128
	_ = v844
	var v847 base.V128
	_ = v847
	var v852 base.V128
	_ = v852
	var v857 base.V128
	_ = v857
	var v862 base.V128
	_ = v862
	var v867 int32
	_ = v867
	var v868 base.V128
	_ = v868
	var v900 int32
	_ = v900
	var v901 base.V128
	_ = v901
	var v902 base.V128
	_ = v902
	var v903 base.V128
	_ = v903
	var v905 base.V128
	_ = v905
	var v909 base.V128
	_ = v909
	var v910 base.V128
	_ = v910
	var v912 base.V128
	_ = v912
	var v914 base.V128
	_ = v914
	var v915 base.V128
	_ = v915
	var v916 base.V128
	_ = v916
	var v918 base.V128
	_ = v918
	var v923 base.V128
	_ = v923
	var v925 int32
	_ = v925
	var v927 int32
	_ = v927
	var v928 base.V128
	_ = v928
	var v930 base.V128
	_ = v930
	var v936 base.V128
	_ = v936
	var v938 base.V128
	_ = v938
	var v940 base.V128
	_ = v940
	var v941 base.V128
	_ = v941
	var v943 base.V128
	_ = v943
	var v951 base.V128
	_ = v951
	var v952 base.V128
	_ = v952
	var v954 int32
	_ = v954
	var v955 base.V128
	_ = v955
	var v957 base.V128
	_ = v957
	var v963 base.V128
	_ = v963
	var v965 base.V128
	_ = v965
	var v967 base.V128
	_ = v967
	var v968 base.V128
	_ = v968
	var v970 base.V128
	_ = v970
	var v978 int32
	_ = v978
	var v979 base.V128
	_ = v979
	var v981 base.V128
	_ = v981
	var v987 base.V128
	_ = v987
	var v989 base.V128
	_ = v989
	var v991 base.V128
	_ = v991
	var v992 base.V128
	_ = v992
	var v994 base.V128
	_ = v994
	var v1002 base.V128
	_ = v1002
	var v1004 base.V128
	_ = v1004
	var v1008 base.V128
	_ = v1008
	var v1010 base.V128
	_ = v1010
	var v1011 base.V128
	_ = v1011
	var v1012 base.V128
	_ = v1012
	var v1014 base.V128
	_ = v1014
	var v1019 base.V128
	_ = v1019
	var v1024 base.V128
	_ = v1024
	var v1026 base.V128
	_ = v1026
	var v1027 base.V128
	_ = v1027
	var v1029 base.V128
	_ = v1029
	var v1037 base.V128
	_ = v1037
	var v1041 base.V128
	_ = v1041
	var v1043 base.V128
	_ = v1043
	var v1044 base.V128
	_ = v1044
	var v1046 base.V128
	_ = v1046
	var v1055 base.V128
	_ = v1055
	var v1057 base.V128
	_ = v1057
	var v1058 base.V128
	_ = v1058
	var v1060 base.V128
	_ = v1060
	var v1068 base.V128
	_ = v1068
	var v1071 base.V128
	_ = v1071
	var v1074 base.V128
	_ = v1074
	var v1076 base.V128
	_ = v1076
	var v1081 base.V128
	_ = v1081
	var v1088 base.V128
	_ = v1088
	var v1105 base.V128
	_ = v1105
	var v1134 base.V128
	_ = v1134
	var v1137 base.V128
	_ = v1137
	var v1144 base.V128
	_ = v1144
	var v1147 base.V128
	_ = v1147
	var v1149 base.V128
	_ = v1149
	var v1154 base.V128
	_ = v1154
	var v1158 base.V128
	_ = v1158
	var v1160 base.V128
	_ = v1160
	var v1165 base.V128
	_ = v1165
	var v1168 base.V128
	_ = v1168
	var v1171 base.V128
	_ = v1171
	var v1176 base.V128
	_ = v1176
	var v1183 base.V128
	_ = v1183
	var v1186 base.V128
	_ = v1186
	var v1193 base.V128
	_ = v1193
	var v1196 base.V128
	_ = v1196
	var v1199 base.V128
	_ = v1199
	var v1204 base.V128
	_ = v1204
	var v1209 base.V128
	_ = v1209
	var v1214 base.V128
	_ = v1214
	var v1217 base.V128
	_ = v1217
	var v1220 base.V128
	_ = v1220
	var v1225 base.V128
	_ = v1225
	var v1230 base.V128
	_ = v1230
	var v1233 base.V128
	_ = v1233
	var v1238 base.V128
	_ = v1238
	var v1243 base.V128
	_ = v1243
	var v1248 base.V128
	_ = v1248
	var v1251 int32
	_ = v1251
	var v1256 int32
	_ = v1256
	var v1274 int32
	_ = v1274
	var v1275 int32
	_ = v1275
	var v1298 int32
	_ = v1298
	var v1301 int32
	_ = v1301
	var v1307 int32
	_ = v1307
	var v1308 int32
	_ = v1308
	var v1310 int32
	_ = v1310
	var v1312 int32
	_ = v1312
	var v1313 int32
	_ = v1313
	var v1317 int32
	_ = v1317
	var v1318 int32
	_ = v1318
	var v1319 int32
	_ = v1319
	var v1321 int32
	_ = v1321
	var v1332 int32
	_ = v1332
	var v1357 int32
	_ = v1357
	var v1358 int32
	_ = v1358
	var v1362 int32
	_ = v1362
	var v1366 int32
	_ = v1366
	var v1367 int32
	_ = v1367
	var v1402 int32
	_ = v1402
	var v1403 int32
	_ = v1403
	var v1408 int64
	_ = v1408
	var v1411 int32
	_ = v1411
	var v1412 int32
	_ = v1412
	var v1424 int32
	_ = v1424
	var v1442 int32
	_ = v1442
	var v1443 int32
	_ = v1443
	var v1444 int32
	_ = v1444
	var v1455 int32
	_ = v1455
	var v1480 int32
	_ = v1480
	var v1481 int32
	_ = v1481
	var v1485 int32
	_ = v1485
	var v1489 int32
	_ = v1489
	var v1490 int32
	_ = v1490
	var v1525 int32
	_ = v1525
	var v1526 int32
	_ = v1526
	var v1531 int64
	_ = v1531
	var v1534 int32
	_ = v1534
	var v1535 int32
	_ = v1535
	var v1547 int32
	_ = v1547
	var v1565 int32
	_ = v1565
	var v1566 base.V128
	_ = v1566
	var v1567 int32
	_ = v1567
	var v1568 base.V128
	_ = v1568
	var v1570 base.V128
	_ = v1570
	var v1571 base.V128
	_ = v1571
	var v1572 int32
	_ = v1572
	var v1573 base.V128
	_ = v1573
	var v1574 base.V128
	_ = v1574
	var v1576 base.V128
	_ = v1576
	var v1577 base.V128
	_ = v1577
	var v1579 base.V128
	_ = v1579
	var v1581 base.V128
	_ = v1581
	var v1583 base.V128
	_ = v1583
	var v1589 base.V128
	_ = v1589
	var v1590 base.V128
	_ = v1590
	var v1596 base.V128
	_ = v1596
	var v1597 base.V128
	_ = v1597
	var v1598 base.V128
	_ = v1598
	var v1599 base.V128
	_ = v1599
	var v1602 base.V128
	_ = v1602
	var v1603 base.V128
	_ = v1603
	var v1606 base.V128
	_ = v1606
	var v1607 base.V128
	_ = v1607
	var v1609 base.V128
	_ = v1609
	var v1613 base.V128
	_ = v1613
	var v1619 int32
	_ = v1619
	var v1621 int32
	_ = v1621
	var v1622 int32
	_ = v1622
	var v1633 int32
	_ = v1633
	var v1658 int32
	_ = v1658
	var v1659 int32
	_ = v1659
	var v1663 int32
	_ = v1663
	var v1667 int32
	_ = v1667
	var v1668 int32
	_ = v1668
	var v1703 int32
	_ = v1703
	var v1704 int32
	_ = v1704
	var v1709 int64
	_ = v1709
	var v1712 int32
	_ = v1712
	var v1713 int32
	_ = v1713
	var v1725 int32
	_ = v1725
	var v1743 int32
	_ = v1743
	var v1754 int32
	_ = v1754
	var v1779 int32
	_ = v1779
	var v1780 int32
	_ = v1780
	var v1784 int32
	_ = v1784
	var v1788 int32
	_ = v1788
	var v1789 int32
	_ = v1789
	var v1824 int32
	_ = v1824
	var v1825 int32
	_ = v1825
	var v1830 int64
	_ = v1830
	var v1833 int32
	_ = v1833
	var v1834 int32
	_ = v1834
	var v1846 int32
	_ = v1846
	var v1865 base.V128
	_ = v1865
	var v1866 int32
	_ = v1866
	var v1867 base.V128
	_ = v1867
	var v1869 base.V128
	_ = v1869
	var v1870 base.V128
	_ = v1870
	var v1872 base.V128
	_ = v1872
	var v1873 base.V128
	_ = v1873
	var v1875 base.V128
	_ = v1875
	var v1876 base.V128
	_ = v1876
	var v1878 base.V128
	_ = v1878
	var v1881 base.V128
	_ = v1881
	var v1887 base.V128
	_ = v1887
	var v1888 base.V128
	_ = v1888
	var v1894 base.V128
	_ = v1894
	var v1895 base.V128
	_ = v1895
	var v1896 base.V128
	_ = v1896
	var v1897 base.V128
	_ = v1897
	var v1900 base.V128
	_ = v1900
	var v1901 base.V128
	_ = v1901
	var v1904 base.V128
	_ = v1904
	var v1905 base.V128
	_ = v1905
	var v1907 base.V128
	_ = v1907
	var v1911 base.V128
	_ = v1911
	var v1917 int32
	_ = v1917
	var v1918 int32
	_ = v1918
	var v1920 int32
	_ = v1920
	var v1921 base.V128
	_ = v1921
	var v1953 int32
	_ = v1953
	var v1954 base.V128
	_ = v1954
	var v1955 base.V128
	_ = v1955
	var v1956 base.V128
	_ = v1956
	var v1958 base.V128
	_ = v1958
	var v1962 base.V128
	_ = v1962
	var v1963 base.V128
	_ = v1963
	var v1965 base.V128
	_ = v1965
	var v1967 base.V128
	_ = v1967
	var v1968 base.V128
	_ = v1968
	var v1969 base.V128
	_ = v1969
	var v1971 base.V128
	_ = v1971
	var v1976 base.V128
	_ = v1976
	var v1978 int32
	_ = v1978
	var v1980 int32
	_ = v1980
	var v1981 base.V128
	_ = v1981
	var v1983 base.V128
	_ = v1983
	var v1989 base.V128
	_ = v1989
	var v1991 base.V128
	_ = v1991
	var v1993 base.V128
	_ = v1993
	var v1994 base.V128
	_ = v1994
	var v1996 base.V128
	_ = v1996
	var v2004 base.V128
	_ = v2004
	var v2005 base.V128
	_ = v2005
	var v2007 int32
	_ = v2007
	var v2008 base.V128
	_ = v2008
	var v2010 base.V128
	_ = v2010
	var v2016 base.V128
	_ = v2016
	var v2018 base.V128
	_ = v2018
	var v2020 base.V128
	_ = v2020
	var v2021 base.V128
	_ = v2021
	var v2023 base.V128
	_ = v2023
	var v2031 int32
	_ = v2031
	var v2032 base.V128
	_ = v2032
	var v2034 base.V128
	_ = v2034
	var v2040 base.V128
	_ = v2040
	var v2042 base.V128
	_ = v2042
	var v2044 base.V128
	_ = v2044
	var v2045 base.V128
	_ = v2045
	var v2047 base.V128
	_ = v2047
	var v2055 base.V128
	_ = v2055
	var v2057 base.V128
	_ = v2057
	var v2061 base.V128
	_ = v2061
	var v2063 base.V128
	_ = v2063
	var v2064 base.V128
	_ = v2064
	var v2065 base.V128
	_ = v2065
	var v2067 base.V128
	_ = v2067
	var v2072 base.V128
	_ = v2072
	var v2077 base.V128
	_ = v2077
	var v2079 base.V128
	_ = v2079
	var v2080 base.V128
	_ = v2080
	var v2082 base.V128
	_ = v2082
	var v2090 base.V128
	_ = v2090
	var v2094 base.V128
	_ = v2094
	var v2096 base.V128
	_ = v2096
	var v2097 base.V128
	_ = v2097
	var v2099 base.V128
	_ = v2099
	var v2108 base.V128
	_ = v2108
	var v2110 base.V128
	_ = v2110
	var v2111 base.V128
	_ = v2111
	var v2113 base.V128
	_ = v2113
	var v2121 base.V128
	_ = v2121
	var v2124 base.V128
	_ = v2124
	var v2127 base.V128
	_ = v2127
	var v2129 base.V128
	_ = v2129
	var v2134 base.V128
	_ = v2134
	var v2141 base.V128
	_ = v2141
	var v2158 base.V128
	_ = v2158
	var v2187 base.V128
	_ = v2187
	var v2190 base.V128
	_ = v2190
	var v2197 base.V128
	_ = v2197
	var v2200 base.V128
	_ = v2200
	var v2202 base.V128
	_ = v2202
	var v2207 base.V128
	_ = v2207
	var v2211 base.V128
	_ = v2211
	var v2213 base.V128
	_ = v2213
	var v2218 base.V128
	_ = v2218
	var v2221 base.V128
	_ = v2221
	var v2224 base.V128
	_ = v2224
	var v2229 base.V128
	_ = v2229
	var v2236 base.V128
	_ = v2236
	var v2239 base.V128
	_ = v2239
	var v2246 base.V128
	_ = v2246
	var v2249 base.V128
	_ = v2249
	var v2252 base.V128
	_ = v2252
	var v2257 base.V128
	_ = v2257
	var v2262 base.V128
	_ = v2262
	var v2267 base.V128
	_ = v2267
	var v2270 base.V128
	_ = v2270
	var v2273 base.V128
	_ = v2273
	var v2278 base.V128
	_ = v2278
	var v2283 base.V128
	_ = v2283
	var v2286 base.V128
	_ = v2286
	var v2291 base.V128
	_ = v2291
	var v2296 base.V128
	_ = v2296
	var v2301 base.V128
	_ = v2301
	var v2304 int32
	_ = v2304
	var v2309 int32
	_ = v2309
	var v2313 int32
	_ = v2313
	var v2314 base.V128
	_ = v2314
	var v2346 int32
	_ = v2346
	var v2347 base.V128
	_ = v2347
	var v2348 base.V128
	_ = v2348
	var v2349 base.V128
	_ = v2349
	var v2351 base.V128
	_ = v2351
	var v2355 base.V128
	_ = v2355
	var v2356 base.V128
	_ = v2356
	var v2358 base.V128
	_ = v2358
	var v2360 base.V128
	_ = v2360
	var v2361 base.V128
	_ = v2361
	var v2362 base.V128
	_ = v2362
	var v2364 base.V128
	_ = v2364
	var v2369 base.V128
	_ = v2369
	var v2371 int32
	_ = v2371
	var v2373 int32
	_ = v2373
	var v2374 base.V128
	_ = v2374
	var v2376 base.V128
	_ = v2376
	var v2382 base.V128
	_ = v2382
	var v2384 base.V128
	_ = v2384
	var v2386 base.V128
	_ = v2386
	var v2387 base.V128
	_ = v2387
	var v2389 base.V128
	_ = v2389
	var v2397 base.V128
	_ = v2397
	var v2398 base.V128
	_ = v2398
	var v2400 int32
	_ = v2400
	var v2401 base.V128
	_ = v2401
	var v2403 base.V128
	_ = v2403
	var v2409 base.V128
	_ = v2409
	var v2411 base.V128
	_ = v2411
	var v2413 base.V128
	_ = v2413
	var v2414 base.V128
	_ = v2414
	var v2416 base.V128
	_ = v2416
	var v2424 int32
	_ = v2424
	var v2425 base.V128
	_ = v2425
	var v2427 base.V128
	_ = v2427
	var v2433 base.V128
	_ = v2433
	var v2435 base.V128
	_ = v2435
	var v2437 base.V128
	_ = v2437
	var v2438 base.V128
	_ = v2438
	var v2440 base.V128
	_ = v2440
	var v2448 base.V128
	_ = v2448
	var v2450 base.V128
	_ = v2450
	var v2454 base.V128
	_ = v2454
	var v2456 base.V128
	_ = v2456
	var v2457 base.V128
	_ = v2457
	var v2458 base.V128
	_ = v2458
	var v2460 base.V128
	_ = v2460
	var v2465 base.V128
	_ = v2465
	var v2470 base.V128
	_ = v2470
	var v2472 base.V128
	_ = v2472
	var v2473 base.V128
	_ = v2473
	var v2475 base.V128
	_ = v2475
	var v2483 base.V128
	_ = v2483
	var v2487 base.V128
	_ = v2487
	var v2489 base.V128
	_ = v2489
	var v2490 base.V128
	_ = v2490
	var v2492 base.V128
	_ = v2492
	var v2501 base.V128
	_ = v2501
	var v2503 base.V128
	_ = v2503
	var v2504 base.V128
	_ = v2504
	var v2506 base.V128
	_ = v2506
	var v2514 base.V128
	_ = v2514
	var v2517 base.V128
	_ = v2517
	var v2520 base.V128
	_ = v2520
	var v2522 base.V128
	_ = v2522
	var v2527 base.V128
	_ = v2527
	var v2534 base.V128
	_ = v2534
	var v2551 base.V128
	_ = v2551
	var v2580 base.V128
	_ = v2580
	var v2583 base.V128
	_ = v2583
	var v2590 base.V128
	_ = v2590
	var v2593 base.V128
	_ = v2593
	var v2595 base.V128
	_ = v2595
	var v2600 base.V128
	_ = v2600
	var v2604 base.V128
	_ = v2604
	var v2606 base.V128
	_ = v2606
	var v2611 base.V128
	_ = v2611
	var v2614 base.V128
	_ = v2614
	var v2617 base.V128
	_ = v2617
	var v2622 base.V128
	_ = v2622
	var v2629 base.V128
	_ = v2629
	var v2632 base.V128
	_ = v2632
	var v2639 base.V128
	_ = v2639
	var v2642 base.V128
	_ = v2642
	var v2645 base.V128
	_ = v2645
	var v2650 base.V128
	_ = v2650
	var v2655 base.V128
	_ = v2655
	var v2660 base.V128
	_ = v2660
	var v2663 base.V128
	_ = v2663
	var v2666 base.V128
	_ = v2666
	var v2671 base.V128
	_ = v2671
	var v2676 base.V128
	_ = v2676
	var v2679 base.V128
	_ = v2679
	var v2684 base.V128
	_ = v2684
	var v2689 base.V128
	_ = v2689
	var v2694 base.V128
	_ = v2694
	var v2698 int32
	_ = v2698
	var v2700 int32
	_ = v2700
	var v2702 int32
	_ = v2702
	var v2703 base.V128
	_ = v2703
	var v2735 int32
	_ = v2735
	var v2736 base.V128
	_ = v2736
	var v2737 base.V128
	_ = v2737
	var v2738 base.V128
	_ = v2738
	var v2740 base.V128
	_ = v2740
	var v2744 base.V128
	_ = v2744
	var v2745 base.V128
	_ = v2745
	var v2747 base.V128
	_ = v2747
	var v2749 base.V128
	_ = v2749
	var v2750 base.V128
	_ = v2750
	var v2751 base.V128
	_ = v2751
	var v2753 base.V128
	_ = v2753
	var v2758 base.V128
	_ = v2758
	var v2760 int32
	_ = v2760
	var v2762 int32
	_ = v2762
	var v2763 base.V128
	_ = v2763
	var v2765 base.V128
	_ = v2765
	var v2771 base.V128
	_ = v2771
	var v2773 base.V128
	_ = v2773
	var v2775 base.V128
	_ = v2775
	var v2776 base.V128
	_ = v2776
	var v2778 base.V128
	_ = v2778
	var v2786 base.V128
	_ = v2786
	var v2787 base.V128
	_ = v2787
	var v2789 int32
	_ = v2789
	var v2790 base.V128
	_ = v2790
	var v2792 base.V128
	_ = v2792
	var v2798 base.V128
	_ = v2798
	var v2800 base.V128
	_ = v2800
	var v2802 base.V128
	_ = v2802
	var v2803 base.V128
	_ = v2803
	var v2805 base.V128
	_ = v2805
	var v2813 int32
	_ = v2813
	var v2814 base.V128
	_ = v2814
	var v2816 base.V128
	_ = v2816
	var v2822 base.V128
	_ = v2822
	var v2824 base.V128
	_ = v2824
	var v2826 base.V128
	_ = v2826
	var v2827 base.V128
	_ = v2827
	var v2829 base.V128
	_ = v2829
	var v2837 base.V128
	_ = v2837
	var v2839 base.V128
	_ = v2839
	var v2843 base.V128
	_ = v2843
	var v2845 base.V128
	_ = v2845
	var v2846 base.V128
	_ = v2846
	var v2847 base.V128
	_ = v2847
	var v2849 base.V128
	_ = v2849
	var v2854 base.V128
	_ = v2854
	var v2859 base.V128
	_ = v2859
	var v2861 base.V128
	_ = v2861
	var v2862 base.V128
	_ = v2862
	var v2864 base.V128
	_ = v2864
	var v2872 base.V128
	_ = v2872
	var v2876 base.V128
	_ = v2876
	var v2878 base.V128
	_ = v2878
	var v2879 base.V128
	_ = v2879
	var v2881 base.V128
	_ = v2881
	var v2890 base.V128
	_ = v2890
	var v2892 base.V128
	_ = v2892
	var v2893 base.V128
	_ = v2893
	var v2895 base.V128
	_ = v2895
	var v2903 base.V128
	_ = v2903
	var v2906 base.V128
	_ = v2906
	var v2909 base.V128
	_ = v2909
	var v2911 base.V128
	_ = v2911
	var v2916 base.V128
	_ = v2916
	var v2923 base.V128
	_ = v2923
	var v2940 base.V128
	_ = v2940
	var v2969 base.V128
	_ = v2969
	var v2972 base.V128
	_ = v2972
	var v2979 base.V128
	_ = v2979
	var v2982 base.V128
	_ = v2982
	var v2984 base.V128
	_ = v2984
	var v2989 base.V128
	_ = v2989
	var v2993 base.V128
	_ = v2993
	var v2995 base.V128
	_ = v2995
	var v3000 base.V128
	_ = v3000
	var v3003 base.V128
	_ = v3003
	var v3006 base.V128
	_ = v3006
	var v3011 base.V128
	_ = v3011
	var v3018 base.V128
	_ = v3018
	var v3021 base.V128
	_ = v3021
	var v3028 base.V128
	_ = v3028
	var v3031 base.V128
	_ = v3031
	var v3034 base.V128
	_ = v3034
	var v3039 base.V128
	_ = v3039
	var v3044 base.V128
	_ = v3044
	var v3049 base.V128
	_ = v3049
	var v3052 base.V128
	_ = v3052
	var v3055 base.V128
	_ = v3055
	var v3060 base.V128
	_ = v3060
	var v3065 base.V128
	_ = v3065
	var v3068 base.V128
	_ = v3068
	var v3073 base.V128
	_ = v3073
	var v3078 base.V128
	_ = v3078
	var v3083 base.V128
	_ = v3083
	var v3086 int32
	_ = v3086
	var v3087 int32
	_ = v3087
	var v3090 int32
	_ = v3090
	var v3091 int32
	_ = v3091
	var v3093 int32
	_ = v3093
	v10 = int32(0)
	v33 = m.G0
	v35 = v33 - int32(528)
	m.G0 = v35
	base.MemoryFill(m, v35+int32(64), v10, int32(463))
	v162 = int32(1)
	v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v169 = int32(base.Ui32(v163+v164)>>(uint(v162)%32)) + v162
	v172 = int32(base.Ui32(v169+v164) >> (uint(v162) % 32))
	v175 = int32(8)
	v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v181 = int32(base.Ui32(v177*int32(_a_F_UpsampleBgrLinePair_SSE2_0)) >> (uint(v175) % 32))
	v182 = int32(base.Ui32(v172*int32(_a_F_UpsampleBgrLinePair_SSE2_1))>>(uint(v175)%32)) + v181
	v184 = v182 + int32(-14234)
	if base.Ui32(v182) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_2)) {
		v191 = int32(0)
	} else {
		v191 = int32(255)
	}
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_3)) {
		v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
	} else {
		v194 = v191
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v194)
	v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v199 = int32(1)
	v202 = int32(base.Ui32(v196+v197)>>(uint(v199)%32)) + v199
	v205 = int32(base.Ui32(v202+v197) >> (uint(v199) % 32))
	v210 = int32(base.Ui32(v205*int32(_a_F_UpsampleBgrLinePair_SSE2_4))>>(uint(int32(8))%32)) + v181
	v212 = v210 + int32(-17685)
	if base.Ui32(v210) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_5)) {
		v219 = int32(0)
	} else {
		v219 = int32(255)
	}
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_3)) {
		v222 = int32(base.Ui32(v212) >> (uint(int32(6)) % 32))
	} else {
		v222 = v219
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v222)
	v226 = int32(8)
	v233 = v181 - (int32(base.Ui32(v205*int32(_a_F_UpsampleBgrLinePair_SSE2_6))>>(uint(v226)%32)) + int32(base.Ui32(v172*int32(_a_F_UpsampleBgrLinePair_SSE2_7))>>(uint(v226)%32)))
	v235 = v233 + int32(_a_F_UpsampleBgrLinePair_SSE2_8)
	if v233 < int32(-8708) {
		v242 = int32(0)
	} else {
		v242 = int32(255)
	}
	if base.Ui32(v235) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_3)) {
		v245 = int32(base.Ui32(v235) >> (uint(int32(6)) % 32))
	} else {
		v245 = v242
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v245)
	if l1 == int32(0) {
	} else {
		v249 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v252 = int32(8)
		v253 = int32(base.Ui32(v249*int32(_a_F_UpsampleBgrLinePair_SSE2_0)) >> (uint(v252) % 32))
		v256 = int32(base.Ui32(v169+v163) >> (uint(int32(1)) % 32))
		v261 = v253 + int32(base.Ui32(v256*int32(_a_F_UpsampleBgrLinePair_SSE2_1))>>(uint(v252)%32))
		v263 = v261 + int32(-14234)
		if base.Ui32(v261) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_2)) {
			v270 = int32(0)
		} else {
			v270 = int32(255)
		}
		if base.Ui32(v263) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_3)) {
			v273 = int32(base.Ui32(v263) >> (uint(int32(6)) % 32))
		} else {
			v273 = v270
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v273)
		v277 = int32(base.Ui32(v202+v196) >> (uint(int32(1)) % 32))
		v282 = v253 + int32(base.Ui32(v277*int32(_a_F_UpsampleBgrLinePair_SSE2_4))>>(uint(int32(8))%32))
		v284 = v282 + int32(-17685)
		if base.Ui32(v282) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_5)) {
			v291 = int32(0)
		} else {
			v291 = int32(255)
		}
		if base.Ui32(v284) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_3)) {
			v294 = int32(base.Ui32(v284) >> (uint(int32(6)) % 32))
		} else {
			v294 = v291
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v294)
		v298 = int32(8)
		v305 = v253 - (int32(base.Ui32(v256*int32(_a_F_UpsampleBgrLinePair_SSE2_7))>>(uint(v298)%32)) + int32(base.Ui32(v277*int32(_a_F_UpsampleBgrLinePair_SSE2_6))>>(uint(v298)%32)))
		v307 = v305 + int32(_a_F_UpsampleBgrLinePair_SSE2_8)
		if v305 < int32(-8708) {
			v314 = int32(0)
		} else {
			v314 = int32(255)
		}
		if base.Ui32(v307) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE2_3)) {
			v317 = int32(base.Ui32(v307) >> (uint(int32(6)) % 32))
		} else {
			v317 = v314
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v317)
	}
	v325 = v35 + int32(96)
	v327 = v35 + int32(64)
	if l8 < int32(34) {
		v1274 = v10
		v1275 = v162
	} else {
		v330 = int32(1)
		v332 = int32(3)
		v339 = v35 + int32(160)
		v341 = v35 + int32(128)
		v342 = int32(0)
		v354 = v342
		v355 = l7 + v332
		v359 = v342
		v360 = l6 + v332
		for {
			v376 = l4 + v354
			v377 = int32(0)
			v378 = base.Simd_g_v128_load_rng(m, v376, v377, int32(0), int32(17))
			v379 = l2 + v354
			v381 = base.Simd_g_v128_load_rng(m, v379, v377, int32(0), int32(17))
			v382 = int32(1)
			v383 = base.Simd_g_v128_load_nc(m, v376, v382)
			v384 = base.Simd_g_i8x16_avgr_u(v381, v383)
			v386 = base.Simd_g_v128_load_nc(m, v379, v382)
			v387 = base.Simd_g_i8x16_avgr_u(v378, v386)
			v389 = base.Simd_g_v128_xor(v383, v381)
			v390 = base.Simd_g_v128_xor(v378, v386)
			v392 = base.Simd_g_v128_xor(v384, v387)
			v394 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k0)
			v396 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v384, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v389, v390), v392), v394))
			v402 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v384), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v384), base.Simd_g_v128_and(v392, v389)), v394))
			v403 = base.Simd_g_i8x16_avgr_u(v378, v402)
			v409 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v387), base.Simd_g_v128_and(v392, v390)), v394))
			v410 = base.Simd_g_i8x16_avgr_u(v383, v409)
			v411 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k1)
			v412 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
			v413 = int32(80)
			base.Simd_g_v128_store(m, v327, v413, v412)
			v415 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
			v416 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v417 = int32(64)
			base.Simd_g_v128_store(m, v327, v417, v416)
			v419 = base.Simd_g_i8x16_avgr_u(v381, v409)
			v420 = base.Simd_g_i8x16_avgr_u(v386, v402)
			v422 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
			v423 = int32(16)
			base.Simd_g_v128_store(m, v327, v423, v422)
			v426 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v327, v377, v426)
			v429 = l5 + v354
			v431 = base.Simd_g_v128_load_rng(m, v429, v377, int32(0), int32(17))
			v432 = l3 + v354
			v434 = base.Simd_g_v128_load_rng(m, v432, v377, int32(0), int32(17))
			v436 = base.Simd_g_v128_load_nc(m, v429, v382)
			v437 = base.Simd_g_i8x16_avgr_u(v434, v436)
			v439 = base.Simd_g_v128_load_nc(m, v432, v382)
			v440 = base.Simd_g_i8x16_avgr_u(v431, v439)
			v442 = base.Simd_g_v128_xor(v436, v434)
			v443 = base.Simd_g_v128_xor(v431, v439)
			v445 = base.Simd_g_v128_xor(v437, v440)
			v448 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v437, v440), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v442, v443), v445), v394))
			v454 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v448, v437), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v448, v437), base.Simd_g_v128_and(v445, v442)), v394))
			v455 = base.Simd_g_i8x16_avgr_u(v431, v454)
			v461 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v448, v440), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v448, v440), base.Simd_g_v128_and(v445, v443)), v394))
			v462 = base.Simd_g_i8x16_avgr_u(v436, v461)
			v464 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v327, int32(112), v464)
			v468 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v327, int32(96), v468)
			v471 = base.Simd_g_i8x16_avgr_u(v434, v461)
			v472 = base.Simd_g_i8x16_avgr_u(v439, v454)
			v474 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
			v475 = int32(48)
			base.Simd_g_v128_store(m, v327, v475, v474)
			v478 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v479 = int32(32)
			base.Simd_g_v128_store(m, v327, v479, v478)
			v481 = l0 + v330 + v359
			v482 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k7)
			v515 = base.Simd_g_v128_load64_zero(m, v481, v377)
			v517 = base.Simd_g_i8x16_shuffle2(v482, v515, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v519 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k8)
			v523 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k9)
			v524 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v517), v519), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v517), v519), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v526 = base.Simd_g_v128_load64_zero(m, v325, v377)
			v528 = base.Simd_g_i8x16_shuffle2(v482, v526, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v529 = base.Simd_g_i32x4_extend_low_i16x8_u(v528)
			v530 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k12)
			v532 = base.Simd_g_i32x4_extend_high_i16x8_u(v528)
			v537 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k13)
			v539 = int32(6)
			v541 = int32(8)
			v542 = base.Simd_g_v128_load64_zero(m, v481, v541)
			v544 = base.Simd_g_i8x16_shuffle2(v482, v542, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v550 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v544), v519), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v544), v519), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v552 = base.Simd_g_v128_load64_zero(m, v325, v541)
			v554 = base.Simd_g_i8x16_shuffle2(v482, v552, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v555 = base.Simd_g_i32x4_extend_low_i16x8_u(v554)
			v557 = base.Simd_g_i32x4_extend_high_i16x8_u(v554)
			v565 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v524, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v529, v530), base.Simd_g_i32x4_mul(v532, v530), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v537), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v550, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v555, v530), base.Simd_g_i32x4_mul(v557, v530), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v537), v539))
			v566 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k14)
			v569 = base.Simd_g_v128_load64_zero(m, v481, v423)
			v571 = base.Simd_g_i8x16_shuffle2(v482, v569, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v577 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v571), v519), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v571), v519), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v579 = base.Simd_g_v128_load64_zero(m, v325, v423)
			v581 = base.Simd_g_i8x16_shuffle2(v482, v579, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v582 = base.Simd_g_i32x4_extend_low_i16x8_u(v581)
			v584 = base.Simd_g_i32x4_extend_high_i16x8_u(v581)
			v592 = int32(24)
			v593 = base.Simd_g_v128_load64_zero(m, v481, v592)
			v595 = base.Simd_g_i8x16_shuffle2(v482, v593, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v601 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v595), v519), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v595), v519), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v603 = base.Simd_g_v128_load64_zero(m, v325, v592)
			v605 = base.Simd_g_i8x16_shuffle2(v482, v603, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v606 = base.Simd_g_i32x4_extend_low_i16x8_u(v605)
			v608 = base.Simd_g_i32x4_extend_high_i16x8_u(v605)
			v616 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v577, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v582, v530), base.Simd_g_i32x4_mul(v584, v530), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v537), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v601, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v606, v530), base.Simd_g_i32x4_mul(v608, v530), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v537), v539))
			v618 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v565, v566), base.Simd_g_v128_and(v616, v566))
			v622 = base.Simd_g_v128_load64_zero(m, v327, v377)
			v624 = base.Simd_g_i8x16_shuffle2(v482, v622, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v625 = base.Simd_g_i32x4_extend_low_i16x8_u(v624)
			v626 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k15)
			v628 = base.Simd_g_i32x4_extend_high_i16x8_u(v624)
			v633 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k16)
			v638 = base.Simd_g_v128_load64_zero(m, v327, v541)
			v640 = base.Simd_g_i8x16_shuffle2(v482, v638, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v641 = base.Simd_g_i32x4_extend_low_i16x8_u(v640)
			v643 = base.Simd_g_i32x4_extend_high_i16x8_u(v640)
			v651 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v625, v626), base.Simd_g_i32x4_mul(v628, v626), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v524), v633), v539), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v641, v626), base.Simd_g_i32x4_mul(v643, v626), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v550), v633), v539))
			v655 = base.Simd_g_v128_load64_zero(m, v327, v423)
			v657 = base.Simd_g_i8x16_shuffle2(v482, v655, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v658 = base.Simd_g_i32x4_extend_low_i16x8_u(v657)
			v660 = base.Simd_g_i32x4_extend_high_i16x8_u(v657)
			v669 = base.Simd_g_v128_load64_zero(m, v327, v592)
			v671 = base.Simd_g_i8x16_shuffle2(v482, v669, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v672 = base.Simd_g_i32x4_extend_low_i16x8_u(v671)
			v674 = base.Simd_g_i32x4_extend_high_i16x8_u(v671)
			v682 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v658, v626), base.Simd_g_i32x4_mul(v660, v626), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v577), v633), v539), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v672, v626), base.Simd_g_i32x4_mul(v674, v626), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v601), v633), v539))
			v685 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v651, v541), base.Simd_g_i16x8_shr_u(v682, v541))
			v688 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v618, v541), base.Simd_g_i16x8_shr_u(v685, v541))
			v690 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k17)
			v695 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k18)
			v702 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k19)
			v719 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v524, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v625, v690), base.Simd_g_i32x4_mul(v628, v690), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v529, v695), base.Simd_g_i32x4_mul(v532, v695), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v702), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v550, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v641, v690), base.Simd_g_i32x4_mul(v643, v690), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v555, v695), base.Simd_g_i32x4_mul(v557, v695), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v702), v539))
			v748 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v577, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v658, v690), base.Simd_g_i32x4_mul(v660, v690), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v582, v695), base.Simd_g_i32x4_mul(v584, v695), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v702), v539), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v601, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v672, v690), base.Simd_g_i32x4_mul(v674, v690), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v606, v695), base.Simd_g_i32x4_mul(v608, v695), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v702), v539))
			v751 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v719, v541), base.Simd_g_i16x8_shr_u(v748, v541))
			v758 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v565, v541), base.Simd_g_i16x8_shr_u(v616, v541))
			v761 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v751, v541), base.Simd_g_i16x8_shr_u(v758, v541))
			v763 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v688, v566), base.Simd_g_v128_and(v761, v566))
			v768 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v651, v566), base.Simd_g_v128_and(v682, v566))
			v772 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v719, v566), base.Simd_g_v128_and(v748, v566))
			v774 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v768, v566), base.Simd_g_v128_and(v772, v566))
			v779 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v618, v566), base.Simd_g_v128_and(v685, v566))
			v782 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v774, v541), base.Simd_g_i16x8_shr_u(v779, v541))
			v785 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v763, v541), base.Simd_g_i16x8_shr_u(v782, v541))
			v790 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v751, v566), base.Simd_g_v128_and(v758, v566))
			v797 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v768, v541), base.Simd_g_i16x8_shr_u(v772, v541))
			v800 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v790, v541), base.Simd_g_i16x8_shr_u(v797, v541))
			v807 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v688, v541), base.Simd_g_i16x8_shr_u(v761, v541))
			v810 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v800, v541), base.Simd_g_i16x8_shr_u(v807, v541))
			v813 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v785, v541), base.Simd_g_i16x8_shr_u(v810, v541))
			base.Simd_g_v128_store(m, v360, v413, v813)
			v818 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v800, v566), base.Simd_g_v128_and(v807, v566))
			v823 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v774, v566), base.Simd_g_v128_and(v779, v566))
			v828 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v790, v566), base.Simd_g_v128_and(v797, v566))
			v831 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v823, v541), base.Simd_g_i16x8_shr_u(v828, v541))
			v834 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v818, v541), base.Simd_g_i16x8_shr_u(v831, v541))
			base.Simd_g_v128_store(m, v360, v417, v834)
			v839 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v823, v566), base.Simd_g_v128_and(v828, v566))
			v844 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v763, v566), base.Simd_g_v128_and(v782, v566))
			v847 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v839, v541), base.Simd_g_i16x8_shr_u(v844, v541))
			base.Simd_g_v128_store(m, v360, v475, v847)
			v852 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v785, v566), base.Simd_g_v128_and(v810, v566))
			base.Simd_g_v128_store(m, v360, v479, v852)
			v857 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v818, v566), base.Simd_g_v128_and(v831, v566))
			base.Simd_g_v128_store(m, v360, v423, v857)
			v862 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v839, v566), base.Simd_g_v128_and(v844, v566))
			base.Simd_g_v128_store(m, v360, v377, v862)
			if l1 == int32(0) {
			} else {
				v867 = l1 + v330 + v359
				v868 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k7)
				v900 = int32(0)
				v901 = base.Simd_g_v128_load64_zero(m, v867, v900)
				v902 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
				v903 = base.Simd_g_i8x16_shuffle2(v868, v901, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v905 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k8)
				v909 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k9)
				v910 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v903), v905), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v903), v905), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
				v912 = base.Simd_g_v128_load64_zero(m, v339, v900)
				v914 = base.Simd_g_i8x16_shuffle2(v868, v912, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v915 = base.Simd_g_i32x4_extend_low_i16x8_u(v914)
				v916 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k12)
				v918 = base.Simd_g_i32x4_extend_high_i16x8_u(v914)
				v923 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k13)
				v925 = int32(6)
				v927 = int32(8)
				v928 = base.Simd_g_v128_load64_zero(m, v867, v927)
				v930 = base.Simd_g_i8x16_shuffle2(v868, v928, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v936 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v930), v905), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v930), v905), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
				v938 = base.Simd_g_v128_load64_zero(m, v339, v927)
				v940 = base.Simd_g_i8x16_shuffle2(v868, v938, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v941 = base.Simd_g_i32x4_extend_low_i16x8_u(v940)
				v943 = base.Simd_g_i32x4_extend_high_i16x8_u(v940)
				v951 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v910, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v915, v916), base.Simd_g_i32x4_mul(v918, v916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v923), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v936, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v941, v916), base.Simd_g_i32x4_mul(v943, v916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v923), v925))
				v952 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k14)
				v954 = int32(16)
				v955 = base.Simd_g_v128_load64_zero(m, v867, v954)
				v957 = base.Simd_g_i8x16_shuffle2(v868, v955, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v963 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v957), v905), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v957), v905), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
				v965 = base.Simd_g_v128_load64_zero(m, v339, v954)
				v967 = base.Simd_g_i8x16_shuffle2(v868, v965, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v968 = base.Simd_g_i32x4_extend_low_i16x8_u(v967)
				v970 = base.Simd_g_i32x4_extend_high_i16x8_u(v967)
				v978 = int32(24)
				v979 = base.Simd_g_v128_load64_zero(m, v867, v978)
				v981 = base.Simd_g_i8x16_shuffle2(v868, v979, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v987 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v981), v905), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v981), v905), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
				v989 = base.Simd_g_v128_load64_zero(m, v339, v978)
				v991 = base.Simd_g_i8x16_shuffle2(v868, v989, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v992 = base.Simd_g_i32x4_extend_low_i16x8_u(v991)
				v994 = base.Simd_g_i32x4_extend_high_i16x8_u(v991)
				v1002 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v963, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v968, v916), base.Simd_g_i32x4_mul(v970, v916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v923), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v987, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v992, v916), base.Simd_g_i32x4_mul(v994, v916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v923), v925))
				v1004 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v951, v952), base.Simd_g_v128_and(v1002, v952))
				v1008 = base.Simd_g_v128_load64_zero(m, v341, v900)
				v1010 = base.Simd_g_i8x16_shuffle2(v868, v1008, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v1011 = base.Simd_g_i32x4_extend_low_i16x8_u(v1010)
				v1012 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k15)
				v1014 = base.Simd_g_i32x4_extend_high_i16x8_u(v1010)
				v1019 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k16)
				v1024 = base.Simd_g_v128_load64_zero(m, v341, v927)
				v1026 = base.Simd_g_i8x16_shuffle2(v868, v1024, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v1027 = base.Simd_g_i32x4_extend_low_i16x8_u(v1026)
				v1029 = base.Simd_g_i32x4_extend_high_i16x8_u(v1026)
				v1037 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1011, v1012), base.Simd_g_i32x4_mul(v1014, v1012), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v910), v1019), v925), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1027, v1012), base.Simd_g_i32x4_mul(v1029, v1012), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v936), v1019), v925))
				v1041 = base.Simd_g_v128_load64_zero(m, v341, v954)
				v1043 = base.Simd_g_i8x16_shuffle2(v868, v1041, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v1044 = base.Simd_g_i32x4_extend_low_i16x8_u(v1043)
				v1046 = base.Simd_g_i32x4_extend_high_i16x8_u(v1043)
				v1055 = base.Simd_g_v128_load64_zero(m, v341, v978)
				v1057 = base.Simd_g_i8x16_shuffle2(v868, v1055, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
				v1058 = base.Simd_g_i32x4_extend_low_i16x8_u(v1057)
				v1060 = base.Simd_g_i32x4_extend_high_i16x8_u(v1057)
				v1068 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1044, v1012), base.Simd_g_i32x4_mul(v1046, v1012), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v963), v1019), v925), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1058, v1012), base.Simd_g_i32x4_mul(v1060, v1012), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v987), v1019), v925))
				v1071 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1037, v927), base.Simd_g_i16x8_shr_u(v1068, v927))
				v1074 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1004, v927), base.Simd_g_i16x8_shr_u(v1071, v927))
				v1076 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k17)
				v1081 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k18)
				v1088 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k19)
				v1105 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v910, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1011, v1076), base.Simd_g_i32x4_mul(v1014, v1076), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v915, v1081), base.Simd_g_i32x4_mul(v918, v1081), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v1088), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v936, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1027, v1076), base.Simd_g_i32x4_mul(v1029, v1076), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v941, v1081), base.Simd_g_i32x4_mul(v943, v1081), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v1088), v925))
				v1134 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v963, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1044, v1076), base.Simd_g_i32x4_mul(v1046, v1076), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v968, v1081), base.Simd_g_i32x4_mul(v970, v1081), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v1088), v925), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v987, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1058, v1076), base.Simd_g_i32x4_mul(v1060, v1076), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v992, v1081), base.Simd_g_i32x4_mul(v994, v1081), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v1088), v925))
				v1137 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1105, v927), base.Simd_g_i16x8_shr_u(v1134, v927))
				v1144 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v951, v927), base.Simd_g_i16x8_shr_u(v1002, v927))
				v1147 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1137, v927), base.Simd_g_i16x8_shr_u(v1144, v927))
				v1149 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1074, v952), base.Simd_g_v128_and(v1147, v952))
				v1154 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1037, v952), base.Simd_g_v128_and(v1068, v952))
				v1158 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1105, v952), base.Simd_g_v128_and(v1134, v952))
				v1160 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1154, v952), base.Simd_g_v128_and(v1158, v952))
				v1165 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1004, v952), base.Simd_g_v128_and(v1071, v952))
				v1168 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1160, v927), base.Simd_g_i16x8_shr_u(v1165, v927))
				v1171 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1149, v927), base.Simd_g_i16x8_shr_u(v1168, v927))
				v1176 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1137, v952), base.Simd_g_v128_and(v1144, v952))
				v1183 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1154, v927), base.Simd_g_i16x8_shr_u(v1158, v927))
				v1186 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1176, v927), base.Simd_g_i16x8_shr_u(v1183, v927))
				v1193 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1074, v927), base.Simd_g_i16x8_shr_u(v1147, v927))
				v1196 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1186, v927), base.Simd_g_i16x8_shr_u(v1193, v927))
				v1199 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1171, v927), base.Simd_g_i16x8_shr_u(v1196, v927))
				base.Simd_g_v128_store(m, v355, int32(80), v1199)
				v1204 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1186, v952), base.Simd_g_v128_and(v1193, v952))
				v1209 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1160, v952), base.Simd_g_v128_and(v1165, v952))
				v1214 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1176, v952), base.Simd_g_v128_and(v1183, v952))
				v1217 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1209, v927), base.Simd_g_i16x8_shr_u(v1214, v927))
				v1220 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1204, v927), base.Simd_g_i16x8_shr_u(v1217, v927))
				base.Simd_g_v128_store(m, v355, int32(64), v1220)
				v1225 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1209, v952), base.Simd_g_v128_and(v1214, v952))
				v1230 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1149, v952), base.Simd_g_v128_and(v1168, v952))
				v1233 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v1225, v927), base.Simd_g_i16x8_shr_u(v1230, v927))
				base.Simd_g_v128_store(m, v355, int32(48), v1233)
				v1238 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1171, v952), base.Simd_g_v128_and(v1196, v952))
				base.Simd_g_v128_store(m, v355, int32(32), v1238)
				v1243 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1204, v952), base.Simd_g_v128_and(v1217, v952))
				base.Simd_g_v128_store(m, v355, v954, v1243)
				v1248 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v1225, v952), base.Simd_g_v128_and(v1230, v952))
				base.Simd_g_v128_store(m, v355, v900, v1248)
			}
			v1251 = int32(96)
			v1256 = v354 + int32(16)
			if v359+int32(66) <= l8 {
				v354 = v1256
				v355 = v355 + v1251
				v359 = v359 + int32(32)
				v360 = v360 + v1251
				continue
			} else {
				break
			}
			break
		}
		v1274 = v1256
		v1275 = v359 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v1298 = int32(32)
		v1301 = int32(1)
		v1307 = int32(base.Ui32(l8+v1301)>>(uint(v1301)%32)) - int32(base.Ui32(v1275)>>(uint(v1301)%32))
		v1308 = F_memcpy(m, v35+v1298, l2+v1274, v1307)
		mBase = m.M
		v1310 = F_memcpy(m, v35, l4+v1274, v1307)
		mBase = m.M
		v1312 = v1310 + v1298
		v1313 = v1312 + v1307
		v1317 = v1307 + int32(-1)
		v1318 = v1312 + v1317
		v1319 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1318))))
		v1321 = int32(17) - v1307
		if base.Ui32(v1321) < base.Ui32(int32(33)) {
			if v1321 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1313))) = uint8(v1319)
				v1332 = v1313 + v1321
				*(*uint8)(unsafe.Add(mBase, uint32(v1332+int32(-1)))) = uint8(v1319)
				if base.Ui32(v1321) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1313)+2)) = uint8(v1319)
					*(*uint8)(unsafe.Add(mBase, uint32(v1313)+1)) = uint8(v1319)
					*(*uint8)(unsafe.Add(mBase, uint32(v1332+int32(-3)))) = uint8(v1319)
					*(*uint8)(unsafe.Add(mBase, uint32(v1332+int32(-2)))) = uint8(v1319)
					if base.Ui32(v1321) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1313)+3)) = uint8(v1319)
						*(*uint8)(unsafe.Add(mBase, uint32(v1332+int32(-4)))) = uint8(v1319)
						if base.Ui32(v1321) < base.Ui32(int32(9)) {
						} else {
							v1357 = (int32(0) - v1313) & int32(3)
							v1358 = v1313 + v1357
							v1362 = v1319 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1358))) = v1362
							v1366 = (v1321 - v1357) & int32(60)
							v1367 = v1358 + v1366
							*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-4)))) = v1362
							if base.Ui32(v1366) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1358)+8)) = v1362
								*(*int32)(unsafe.Add(mBase, uint32(v1358)+4)) = v1362
								*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-8)))) = v1362
								*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-12)))) = v1362
								if base.Ui32(v1366) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1358)+24)) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1358)+20)) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1358)+16)) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1358)+12)) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-16)))) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-20)))) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-24)))) = v1362
									*(*int32)(unsafe.Add(mBase, uint32(v1367+int32(-28)))) = v1362
									v1402 = v1358&int32(4) | int32(24)
									v1403 = v1366 - v1402
									if base.Ui32(v1403) < base.Ui32(int32(32)) {
									} else {
										v1408 = base.I64_extend_i32_u(v1362) * int64(4294967297)
										v1411 = v1403
										v1412 = v1358 + v1402
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1412)+24)) = v1408
											*(*int64)(unsafe.Add(mBase, uint32(v1412)+16)) = v1408
											*(*int64)(unsafe.Add(mBase, uint32(v1412)+8)) = v1408
											*(*int64)(unsafe.Add(mBase, uint32(v1412))) = v1408
											v1424 = v1411 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1424) {
												v1411 = v1424
												v1412 = v1412 + int32(32)
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
			base.MemoryFill(m, v1313, v1319, v1321)
		}
		v1442 = v1310 + v1307
		v1443 = v1310 + v1317
		v1444 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1443))))
		if base.Ui32(v1321) < base.Ui32(int32(33)) {
			if v1321 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1442))) = uint8(v1444)
				v1455 = v1442 + v1321
				*(*uint8)(unsafe.Add(mBase, uint32(v1455+int32(-1)))) = uint8(v1444)
				if base.Ui32(v1321) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1442)+2)) = uint8(v1444)
					*(*uint8)(unsafe.Add(mBase, uint32(v1442)+1)) = uint8(v1444)
					*(*uint8)(unsafe.Add(mBase, uint32(v1455+int32(-3)))) = uint8(v1444)
					*(*uint8)(unsafe.Add(mBase, uint32(v1455+int32(-2)))) = uint8(v1444)
					if base.Ui32(v1321) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1442)+3)) = uint8(v1444)
						*(*uint8)(unsafe.Add(mBase, uint32(v1455+int32(-4)))) = uint8(v1444)
						if base.Ui32(v1321) < base.Ui32(int32(9)) {
						} else {
							v1480 = (int32(0) - v1442) & int32(3)
							v1481 = v1442 + v1480
							v1485 = v1444 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1481))) = v1485
							v1489 = (v1321 - v1480) & int32(60)
							v1490 = v1481 + v1489
							*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-4)))) = v1485
							if base.Ui32(v1489) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1481)+8)) = v1485
								*(*int32)(unsafe.Add(mBase, uint32(v1481)+4)) = v1485
								*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-8)))) = v1485
								*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-12)))) = v1485
								if base.Ui32(v1489) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1481)+24)) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1481)+20)) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1481)+16)) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1481)+12)) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-16)))) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-20)))) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-24)))) = v1485
									*(*int32)(unsafe.Add(mBase, uint32(v1490+int32(-28)))) = v1485
									v1525 = v1481&int32(4) | int32(24)
									v1526 = v1489 - v1525
									if base.Ui32(v1526) < base.Ui32(int32(32)) {
									} else {
										v1531 = base.I64_extend_i32_u(v1485) * int64(4294967297)
										v1534 = v1526
										v1535 = v1481 + v1525
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1535)+24)) = v1531
											*(*int64)(unsafe.Add(mBase, uint32(v1535)+16)) = v1531
											*(*int64)(unsafe.Add(mBase, uint32(v1535)+8)) = v1531
											*(*int64)(unsafe.Add(mBase, uint32(v1535))) = v1531
											v1547 = v1534 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1547) {
												v1534 = v1547
												v1535 = v1535 + int32(32)
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
			base.MemoryFill(m, v1442, v1444, v1321)
		}
		v1565 = int32(0)
		v1566 = base.Simd_g_v128_load_rng(m, v1310, v1565, int32(0), int32(49))
		v1567 = int32(32)
		v1568 = base.Simd_g_v128_load_nc(m, v1310, v1567)
		v1570 = base.Simd_g_v128_load_nc(m, v1310, int32(1))
		v1571 = base.Simd_g_i8x16_avgr_u(v1568, v1570)
		v1572 = int32(33)
		v1573 = base.Simd_g_v128_load_nc(m, v1310, v1572)
		v1574 = base.Simd_g_i8x16_avgr_u(v1566, v1573)
		v1576 = base.Simd_g_v128_xor(v1570, v1568)
		v1577 = base.Simd_g_v128_xor(v1566, v1573)
		v1579 = base.Simd_g_v128_xor(v1571, v1574)
		v1581 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k0)
		v1583 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1571, v1574), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1576, v1577), v1579), v1581))
		v1589 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1583, v1571), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1583, v1571), base.Simd_g_v128_and(v1579, v1576)), v1581))
		v1590 = base.Simd_g_i8x16_avgr_u(v1566, v1589)
		v1596 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1583, v1574), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1583, v1574), base.Simd_g_v128_and(v1579, v1577)), v1581))
		v1597 = base.Simd_g_i8x16_avgr_u(v1570, v1596)
		v1598 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k1)
		v1599 = base.Simd_g_i8x16_shuffle2(v1590, v1597, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(80), v1599)
		v1602 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
		v1603 = base.Simd_g_i8x16_shuffle2(v1590, v1597, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, int32(64), v1603)
		v1606 = base.Simd_g_i8x16_avgr_u(v1568, v1596)
		v1607 = base.Simd_g_i8x16_avgr_u(v1573, v1589)
		v1609 = base.Simd_g_i8x16_shuffle2(v1606, v1607, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(16), v1609)
		v1613 = base.Simd_g_i8x16_shuffle2(v1606, v1607, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, v1565, v1613)
		v1619 = F_memcpy(m, v1310+v1567, l3+v1274, v1307)
		mBase = m.M
		v1621 = F_memcpy(m, v1310, l5+v1274, v1307)
		mBase = m.M
		v1622 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1318))))
		if base.Ui32(v1321) < base.Ui32(v1572) {
			if v1321 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1313))) = uint8(v1622)
				v1633 = v1313 + v1321
				*(*uint8)(unsafe.Add(mBase, uint32(v1633+int32(-1)))) = uint8(v1622)
				if base.Ui32(v1321) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1313)+2)) = uint8(v1622)
					*(*uint8)(unsafe.Add(mBase, uint32(v1313)+1)) = uint8(v1622)
					*(*uint8)(unsafe.Add(mBase, uint32(v1633+int32(-3)))) = uint8(v1622)
					*(*uint8)(unsafe.Add(mBase, uint32(v1633+int32(-2)))) = uint8(v1622)
					if base.Ui32(v1321) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1313)+3)) = uint8(v1622)
						*(*uint8)(unsafe.Add(mBase, uint32(v1633+int32(-4)))) = uint8(v1622)
						if base.Ui32(v1321) < base.Ui32(int32(9)) {
						} else {
							v1658 = (int32(0) - v1313) & int32(3)
							v1659 = v1313 + v1658
							v1663 = v1622 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1659))) = v1663
							v1667 = (v1321 - v1658) & int32(60)
							v1668 = v1659 + v1667
							*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-4)))) = v1663
							if base.Ui32(v1667) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1659)+8)) = v1663
								*(*int32)(unsafe.Add(mBase, uint32(v1659)+4)) = v1663
								*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-8)))) = v1663
								*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-12)))) = v1663
								if base.Ui32(v1667) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1659)+24)) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1659)+20)) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1659)+16)) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1659)+12)) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-16)))) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-20)))) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-24)))) = v1663
									*(*int32)(unsafe.Add(mBase, uint32(v1668+int32(-28)))) = v1663
									v1703 = v1659&int32(4) | int32(24)
									v1704 = v1667 - v1703
									if base.Ui32(v1704) < base.Ui32(int32(32)) {
									} else {
										v1709 = base.I64_extend_i32_u(v1663) * int64(4294967297)
										v1712 = v1704
										v1713 = v1659 + v1703
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1713)+24)) = v1709
											*(*int64)(unsafe.Add(mBase, uint32(v1713)+16)) = v1709
											*(*int64)(unsafe.Add(mBase, uint32(v1713)+8)) = v1709
											*(*int64)(unsafe.Add(mBase, uint32(v1713))) = v1709
											v1725 = v1712 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1725) {
												v1712 = v1725
												v1713 = v1713 + int32(32)
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
			base.MemoryFill(m, v1313, v1622, v1321)
		}
		v1743 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1443))))
		if base.Ui32(v1321) < base.Ui32(int32(33)) {
			if v1321 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1442))) = uint8(v1743)
				v1754 = v1442 + v1321
				*(*uint8)(unsafe.Add(mBase, uint32(v1754+int32(-1)))) = uint8(v1743)
				if base.Ui32(v1321) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1442)+2)) = uint8(v1743)
					*(*uint8)(unsafe.Add(mBase, uint32(v1442)+1)) = uint8(v1743)
					*(*uint8)(unsafe.Add(mBase, uint32(v1754+int32(-3)))) = uint8(v1743)
					*(*uint8)(unsafe.Add(mBase, uint32(v1754+int32(-2)))) = uint8(v1743)
					if base.Ui32(v1321) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1442)+3)) = uint8(v1743)
						*(*uint8)(unsafe.Add(mBase, uint32(v1754+int32(-4)))) = uint8(v1743)
						if base.Ui32(v1321) < base.Ui32(int32(9)) {
						} else {
							v1779 = (int32(0) - v1442) & int32(3)
							v1780 = v1442 + v1779
							v1784 = v1743 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1780))) = v1784
							v1788 = (v1321 - v1779) & int32(60)
							v1789 = v1780 + v1788
							*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-4)))) = v1784
							if base.Ui32(v1788) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1780)+8)) = v1784
								*(*int32)(unsafe.Add(mBase, uint32(v1780)+4)) = v1784
								*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-8)))) = v1784
								*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-12)))) = v1784
								if base.Ui32(v1788) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1780)+24)) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1780)+20)) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1780)+16)) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1780)+12)) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-16)))) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-20)))) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-24)))) = v1784
									*(*int32)(unsafe.Add(mBase, uint32(v1789+int32(-28)))) = v1784
									v1824 = v1780&int32(4) | int32(24)
									v1825 = v1788 - v1824
									if base.Ui32(v1825) < base.Ui32(int32(32)) {
									} else {
										v1830 = base.I64_extend_i32_u(v1784) * int64(4294967297)
										v1833 = v1825
										v1834 = v1780 + v1824
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1834)+24)) = v1830
											*(*int64)(unsafe.Add(mBase, uint32(v1834)+16)) = v1830
											*(*int64)(unsafe.Add(mBase, uint32(v1834)+8)) = v1830
											*(*int64)(unsafe.Add(mBase, uint32(v1834))) = v1830
											v1846 = v1833 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1846) {
												v1833 = v1846
												v1834 = v1834 + int32(32)
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
			base.MemoryFill(m, v1442, v1743, v1321)
		}
		v1865 = base.Simd_g_v128_load_rng(m, v1621, int32(0), int32(0), int32(49))
		v1866 = int32(32)
		v1867 = base.Simd_g_v128_load_nc(m, v1621, v1866)
		v1869 = base.Simd_g_v128_load_nc(m, v1621, int32(1))
		v1870 = base.Simd_g_i8x16_avgr_u(v1867, v1869)
		v1872 = base.Simd_g_v128_load_nc(m, v1621, int32(33))
		v1873 = base.Simd_g_i8x16_avgr_u(v1865, v1872)
		v1875 = base.Simd_g_v128_xor(v1869, v1867)
		v1876 = base.Simd_g_v128_xor(v1865, v1872)
		v1878 = base.Simd_g_v128_xor(v1870, v1873)
		v1881 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1870, v1873), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1875, v1876), v1878), v1581))
		v1887 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1881, v1870), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1881, v1870), base.Simd_g_v128_and(v1878, v1875)), v1581))
		v1888 = base.Simd_g_i8x16_avgr_u(v1865, v1887)
		v1894 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1881, v1873), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1881, v1873), base.Simd_g_v128_and(v1878, v1876)), v1581))
		v1895 = base.Simd_g_i8x16_avgr_u(v1869, v1894)
		v1896 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k1)
		v1897 = base.Simd_g_i8x16_shuffle2(v1888, v1895, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(112), v1897)
		v1900 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
		v1901 = base.Simd_g_i8x16_shuffle2(v1888, v1895, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, int32(96), v1901)
		v1904 = base.Simd_g_i8x16_avgr_u(v1867, v1894)
		v1905 = base.Simd_g_i8x16_avgr_u(v1872, v1887)
		v1907 = base.Simd_g_i8x16_shuffle2(v1904, v1905, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v327, int32(48), v1907)
		v1911 = base.Simd_g_i8x16_shuffle2(v1904, v1905, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v327, v1866, v1911)
		v1917 = l8 - v1275
		v1918 = F_memcpy(m, v35+int32(448), l0+v1275, v1917)
		mBase = m.M
		v1920 = v35 + int32(192)
		if l1 != 0 {
			v2313 = F_memcpy(m, v35+int32(480), l1+v1275, v1917)
			mBase = m.M
			v2314 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k7)
			v2346 = int32(0)
			v2347 = base.Simd_g_v128_load64_zero(m, v1918, v2346)
			v2348 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
			v2349 = base.Simd_g_i8x16_shuffle2(v2314, v2347, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2351 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k8)
			v2355 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k9)
			v2356 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2349), v2351), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2349), v2351), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2358 = base.Simd_g_v128_load64_zero(m, v325, v2346)
			v2360 = base.Simd_g_i8x16_shuffle2(v2314, v2358, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2361 = base.Simd_g_i32x4_extend_low_i16x8_u(v2360)
			v2362 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k12)
			v2364 = base.Simd_g_i32x4_extend_high_i16x8_u(v2360)
			v2369 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k13)
			v2371 = int32(6)
			v2373 = int32(8)
			v2374 = base.Simd_g_v128_load64_zero(m, v1918, v2373)
			v2376 = base.Simd_g_i8x16_shuffle2(v2314, v2374, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2382 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2376), v2351), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2376), v2351), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2384 = base.Simd_g_v128_load64_zero(m, v325, v2373)
			v2386 = base.Simd_g_i8x16_shuffle2(v2314, v2384, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2387 = base.Simd_g_i32x4_extend_low_i16x8_u(v2386)
			v2389 = base.Simd_g_i32x4_extend_high_i16x8_u(v2386)
			v2397 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2356, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2361, v2362), base.Simd_g_i32x4_mul(v2364, v2362), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2369), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2382, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2387, v2362), base.Simd_g_i32x4_mul(v2389, v2362), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2369), v2371))
			v2398 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k14)
			v2400 = int32(16)
			v2401 = base.Simd_g_v128_load64_zero(m, v1918, v2400)
			v2403 = base.Simd_g_i8x16_shuffle2(v2314, v2401, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2409 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2403), v2351), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2403), v2351), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2411 = base.Simd_g_v128_load64_zero(m, v325, v2400)
			v2413 = base.Simd_g_i8x16_shuffle2(v2314, v2411, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2414 = base.Simd_g_i32x4_extend_low_i16x8_u(v2413)
			v2416 = base.Simd_g_i32x4_extend_high_i16x8_u(v2413)
			v2424 = int32(24)
			v2425 = base.Simd_g_v128_load64_zero(m, v1918, v2424)
			v2427 = base.Simd_g_i8x16_shuffle2(v2314, v2425, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2433 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2427), v2351), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2427), v2351), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2435 = base.Simd_g_v128_load64_zero(m, v325, v2424)
			v2437 = base.Simd_g_i8x16_shuffle2(v2314, v2435, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2438 = base.Simd_g_i32x4_extend_low_i16x8_u(v2437)
			v2440 = base.Simd_g_i32x4_extend_high_i16x8_u(v2437)
			v2448 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2409, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2414, v2362), base.Simd_g_i32x4_mul(v2416, v2362), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2369), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2433, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2438, v2362), base.Simd_g_i32x4_mul(v2440, v2362), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2369), v2371))
			v2450 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2397, v2398), base.Simd_g_v128_and(v2448, v2398))
			v2454 = base.Simd_g_v128_load64_zero(m, v327, v2346)
			v2456 = base.Simd_g_i8x16_shuffle2(v2314, v2454, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2457 = base.Simd_g_i32x4_extend_low_i16x8_u(v2456)
			v2458 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k15)
			v2460 = base.Simd_g_i32x4_extend_high_i16x8_u(v2456)
			v2465 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k16)
			v2470 = base.Simd_g_v128_load64_zero(m, v327, v2373)
			v2472 = base.Simd_g_i8x16_shuffle2(v2314, v2470, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2473 = base.Simd_g_i32x4_extend_low_i16x8_u(v2472)
			v2475 = base.Simd_g_i32x4_extend_high_i16x8_u(v2472)
			v2483 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2457, v2458), base.Simd_g_i32x4_mul(v2460, v2458), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2356), v2465), v2371), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2473, v2458), base.Simd_g_i32x4_mul(v2475, v2458), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2382), v2465), v2371))
			v2487 = base.Simd_g_v128_load64_zero(m, v327, v2400)
			v2489 = base.Simd_g_i8x16_shuffle2(v2314, v2487, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2490 = base.Simd_g_i32x4_extend_low_i16x8_u(v2489)
			v2492 = base.Simd_g_i32x4_extend_high_i16x8_u(v2489)
			v2501 = base.Simd_g_v128_load64_zero(m, v327, v2424)
			v2503 = base.Simd_g_i8x16_shuffle2(v2314, v2501, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2504 = base.Simd_g_i32x4_extend_low_i16x8_u(v2503)
			v2506 = base.Simd_g_i32x4_extend_high_i16x8_u(v2503)
			v2514 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2490, v2458), base.Simd_g_i32x4_mul(v2492, v2458), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2409), v2465), v2371), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2504, v2458), base.Simd_g_i32x4_mul(v2506, v2458), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2433), v2465), v2371))
			v2517 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2483, v2373), base.Simd_g_i16x8_shr_u(v2514, v2373))
			v2520 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2450, v2373), base.Simd_g_i16x8_shr_u(v2517, v2373))
			v2522 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k17)
			v2527 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k18)
			v2534 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k19)
			v2551 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2356, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2457, v2522), base.Simd_g_i32x4_mul(v2460, v2522), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2361, v2527), base.Simd_g_i32x4_mul(v2364, v2527), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2534), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2382, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2473, v2522), base.Simd_g_i32x4_mul(v2475, v2522), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2387, v2527), base.Simd_g_i32x4_mul(v2389, v2527), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2534), v2371))
			v2580 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2409, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2490, v2522), base.Simd_g_i32x4_mul(v2492, v2522), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2414, v2527), base.Simd_g_i32x4_mul(v2416, v2527), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2534), v2371), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2433, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2504, v2522), base.Simd_g_i32x4_mul(v2506, v2522), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2438, v2527), base.Simd_g_i32x4_mul(v2440, v2527), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2534), v2371))
			v2583 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2551, v2373), base.Simd_g_i16x8_shr_u(v2580, v2373))
			v2590 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2397, v2373), base.Simd_g_i16x8_shr_u(v2448, v2373))
			v2593 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2583, v2373), base.Simd_g_i16x8_shr_u(v2590, v2373))
			v2595 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2520, v2398), base.Simd_g_v128_and(v2593, v2398))
			v2600 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2483, v2398), base.Simd_g_v128_and(v2514, v2398))
			v2604 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2551, v2398), base.Simd_g_v128_and(v2580, v2398))
			v2606 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2600, v2398), base.Simd_g_v128_and(v2604, v2398))
			v2611 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2450, v2398), base.Simd_g_v128_and(v2517, v2398))
			v2614 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2606, v2373), base.Simd_g_i16x8_shr_u(v2611, v2373))
			v2617 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2595, v2373), base.Simd_g_i16x8_shr_u(v2614, v2373))
			v2622 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2583, v2398), base.Simd_g_v128_and(v2590, v2398))
			v2629 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2600, v2373), base.Simd_g_i16x8_shr_u(v2604, v2373))
			v2632 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2622, v2373), base.Simd_g_i16x8_shr_u(v2629, v2373))
			v2639 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2520, v2373), base.Simd_g_i16x8_shr_u(v2593, v2373))
			v2642 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2632, v2373), base.Simd_g_i16x8_shr_u(v2639, v2373))
			v2645 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2617, v2373), base.Simd_g_i16x8_shr_u(v2642, v2373))
			base.Simd_g_v128_store(m, v1920, int32(80), v2645)
			v2650 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2632, v2398), base.Simd_g_v128_and(v2639, v2398))
			v2655 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2606, v2398), base.Simd_g_v128_and(v2611, v2398))
			v2660 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2622, v2398), base.Simd_g_v128_and(v2629, v2398))
			v2663 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2655, v2373), base.Simd_g_i16x8_shr_u(v2660, v2373))
			v2666 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2650, v2373), base.Simd_g_i16x8_shr_u(v2663, v2373))
			base.Simd_g_v128_store(m, v1920, int32(64), v2666)
			v2671 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2655, v2398), base.Simd_g_v128_and(v2660, v2398))
			v2676 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2595, v2398), base.Simd_g_v128_and(v2614, v2398))
			v2679 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2671, v2373), base.Simd_g_i16x8_shr_u(v2676, v2373))
			base.Simd_g_v128_store(m, v1920, int32(48), v2679)
			v2684 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2617, v2398), base.Simd_g_v128_and(v2642, v2398))
			base.Simd_g_v128_store(m, v1920, int32(32), v2684)
			v2689 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2650, v2398), base.Simd_g_v128_and(v2663, v2398))
			base.Simd_g_v128_store(m, v1920, v2400, v2689)
			v2694 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2671, v2398), base.Simd_g_v128_and(v2676, v2398))
			base.Simd_g_v128_store(m, v1920, v2346, v2694)
			v2698 = v35 + int32(128)
			v2700 = v35 + int32(160)
			v2702 = v35 + int32(320)
			v2703 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k7)
			v2735 = int32(0)
			v2736 = base.Simd_g_v128_load64_zero(m, v2313, v2735)
			v2737 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
			v2738 = base.Simd_g_i8x16_shuffle2(v2703, v2736, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2740 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k8)
			v2744 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k9)
			v2745 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2738), v2740), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2738), v2740), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2747 = base.Simd_g_v128_load64_zero(m, v2700, v2735)
			v2749 = base.Simd_g_i8x16_shuffle2(v2703, v2747, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2750 = base.Simd_g_i32x4_extend_low_i16x8_u(v2749)
			v2751 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k12)
			v2753 = base.Simd_g_i32x4_extend_high_i16x8_u(v2749)
			v2758 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k13)
			v2760 = int32(6)
			v2762 = int32(8)
			v2763 = base.Simd_g_v128_load64_zero(m, v2313, v2762)
			v2765 = base.Simd_g_i8x16_shuffle2(v2703, v2763, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2771 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2765), v2740), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2765), v2740), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2773 = base.Simd_g_v128_load64_zero(m, v2700, v2762)
			v2775 = base.Simd_g_i8x16_shuffle2(v2703, v2773, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2776 = base.Simd_g_i32x4_extend_low_i16x8_u(v2775)
			v2778 = base.Simd_g_i32x4_extend_high_i16x8_u(v2775)
			v2786 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2745, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2750, v2751), base.Simd_g_i32x4_mul(v2753, v2751), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2758), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2771, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2776, v2751), base.Simd_g_i32x4_mul(v2778, v2751), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2758), v2760))
			v2787 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k14)
			v2789 = int32(16)
			v2790 = base.Simd_g_v128_load64_zero(m, v2313, v2789)
			v2792 = base.Simd_g_i8x16_shuffle2(v2703, v2790, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2798 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2792), v2740), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2792), v2740), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2800 = base.Simd_g_v128_load64_zero(m, v2700, v2789)
			v2802 = base.Simd_g_i8x16_shuffle2(v2703, v2800, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2803 = base.Simd_g_i32x4_extend_low_i16x8_u(v2802)
			v2805 = base.Simd_g_i32x4_extend_high_i16x8_u(v2802)
			v2813 = int32(24)
			v2814 = base.Simd_g_v128_load64_zero(m, v2313, v2813)
			v2816 = base.Simd_g_i8x16_shuffle2(v2703, v2814, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2822 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2816), v2740), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2816), v2740), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2824 = base.Simd_g_v128_load64_zero(m, v2700, v2813)
			v2826 = base.Simd_g_i8x16_shuffle2(v2703, v2824, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2827 = base.Simd_g_i32x4_extend_low_i16x8_u(v2826)
			v2829 = base.Simd_g_i32x4_extend_high_i16x8_u(v2826)
			v2837 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2798, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2803, v2751), base.Simd_g_i32x4_mul(v2805, v2751), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2758), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2822, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2827, v2751), base.Simd_g_i32x4_mul(v2829, v2751), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v2758), v2760))
			v2839 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2786, v2787), base.Simd_g_v128_and(v2837, v2787))
			v2843 = base.Simd_g_v128_load64_zero(m, v2698, v2735)
			v2845 = base.Simd_g_i8x16_shuffle2(v2703, v2843, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2846 = base.Simd_g_i32x4_extend_low_i16x8_u(v2845)
			v2847 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k15)
			v2849 = base.Simd_g_i32x4_extend_high_i16x8_u(v2845)
			v2854 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k16)
			v2859 = base.Simd_g_v128_load64_zero(m, v2698, v2762)
			v2861 = base.Simd_g_i8x16_shuffle2(v2703, v2859, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2862 = base.Simd_g_i32x4_extend_low_i16x8_u(v2861)
			v2864 = base.Simd_g_i32x4_extend_high_i16x8_u(v2861)
			v2872 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2846, v2847), base.Simd_g_i32x4_mul(v2849, v2847), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2745), v2854), v2760), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2862, v2847), base.Simd_g_i32x4_mul(v2864, v2847), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2771), v2854), v2760))
			v2876 = base.Simd_g_v128_load64_zero(m, v2698, v2789)
			v2878 = base.Simd_g_i8x16_shuffle2(v2703, v2876, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2879 = base.Simd_g_i32x4_extend_low_i16x8_u(v2878)
			v2881 = base.Simd_g_i32x4_extend_high_i16x8_u(v2878)
			v2890 = base.Simd_g_v128_load64_zero(m, v2698, v2813)
			v2892 = base.Simd_g_i8x16_shuffle2(v2703, v2890, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2893 = base.Simd_g_i32x4_extend_low_i16x8_u(v2892)
			v2895 = base.Simd_g_i32x4_extend_high_i16x8_u(v2892)
			v2903 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2879, v2847), base.Simd_g_i32x4_mul(v2881, v2847), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2798), v2854), v2760), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2893, v2847), base.Simd_g_i32x4_mul(v2895, v2847), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2822), v2854), v2760))
			v2906 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2872, v2762), base.Simd_g_i16x8_shr_u(v2903, v2762))
			v2909 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2839, v2762), base.Simd_g_i16x8_shr_u(v2906, v2762))
			v2911 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k17)
			v2916 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k18)
			v2923 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k19)
			v2940 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2745, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2846, v2911), base.Simd_g_i32x4_mul(v2849, v2911), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2750, v2916), base.Simd_g_i32x4_mul(v2753, v2916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2923), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2771, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2862, v2911), base.Simd_g_i32x4_mul(v2864, v2911), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2776, v2916), base.Simd_g_i32x4_mul(v2778, v2916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2923), v2760))
			v2969 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2798, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2879, v2911), base.Simd_g_i32x4_mul(v2881, v2911), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2803, v2916), base.Simd_g_i32x4_mul(v2805, v2916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2923), v2760), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2822, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2893, v2911), base.Simd_g_i32x4_mul(v2895, v2911), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2827, v2916), base.Simd_g_i32x4_mul(v2829, v2916), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2923), v2760))
			v2972 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2940, v2762), base.Simd_g_i16x8_shr_u(v2969, v2762))
			v2979 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2786, v2762), base.Simd_g_i16x8_shr_u(v2837, v2762))
			v2982 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2972, v2762), base.Simd_g_i16x8_shr_u(v2979, v2762))
			v2984 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2909, v2787), base.Simd_g_v128_and(v2982, v2787))
			v2989 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2872, v2787), base.Simd_g_v128_and(v2903, v2787))
			v2993 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2940, v2787), base.Simd_g_v128_and(v2969, v2787))
			v2995 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2989, v2787), base.Simd_g_v128_and(v2993, v2787))
			v3000 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2839, v2787), base.Simd_g_v128_and(v2906, v2787))
			v3003 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2995, v2762), base.Simd_g_i16x8_shr_u(v3000, v2762))
			v3006 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2984, v2762), base.Simd_g_i16x8_shr_u(v3003, v2762))
			v3011 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2972, v2787), base.Simd_g_v128_and(v2979, v2787))
			v3018 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2989, v2762), base.Simd_g_i16x8_shr_u(v2993, v2762))
			v3021 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v3011, v2762), base.Simd_g_i16x8_shr_u(v3018, v2762))
			v3028 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2909, v2762), base.Simd_g_i16x8_shr_u(v2982, v2762))
			v3031 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v3021, v2762), base.Simd_g_i16x8_shr_u(v3028, v2762))
			v3034 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v3006, v2762), base.Simd_g_i16x8_shr_u(v3031, v2762))
			base.Simd_g_v128_store(m, v2702, int32(80), v3034)
			v3039 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v3021, v2787), base.Simd_g_v128_and(v3028, v2787))
			v3044 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2995, v2787), base.Simd_g_v128_and(v3000, v2787))
			v3049 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v3011, v2787), base.Simd_g_v128_and(v3018, v2787))
			v3052 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v3044, v2762), base.Simd_g_i16x8_shr_u(v3049, v2762))
			v3055 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v3039, v2762), base.Simd_g_i16x8_shr_u(v3052, v2762))
			base.Simd_g_v128_store(m, v2702, int32(64), v3055)
			v3060 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v3044, v2787), base.Simd_g_v128_and(v3049, v2787))
			v3065 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2984, v2787), base.Simd_g_v128_and(v3003, v2787))
			v3068 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v3060, v2762), base.Simd_g_i16x8_shr_u(v3065, v2762))
			base.Simd_g_v128_store(m, v2702, int32(48), v3068)
			v3073 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v3006, v2787), base.Simd_g_v128_and(v3031, v2787))
			base.Simd_g_v128_store(m, v2702, int32(32), v3073)
			v3078 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v3039, v2787), base.Simd_g_v128_and(v3052, v2787))
			base.Simd_g_v128_store(m, v2702, v2789, v3078)
			v3083 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v3060, v2787), base.Simd_g_v128_and(v3065, v2787))
			base.Simd_g_v128_store(m, v2702, v2735, v3083)
			v3086 = int32(3)
			v3087 = v1275 * v3086
			v3090 = v1917 * v3086
			v3091 = F_memcpy(m, l6+v3087, v1920, v3090)
			mBase = m.M
			v3093 = F_memcpy(m, l7+v3087, v2702, v3090)
			mBase = m.M
		} else {
			v1921 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k7)
			v1953 = int32(0)
			v1954 = base.Simd_g_v128_load64_zero(m, v1918, v1953)
			v1955 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k4)
			v1956 = base.Simd_g_i8x16_shuffle2(v1921, v1954, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v1958 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k8)
			v1962 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k9)
			v1963 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1956), v1958), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1956), v1958), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v1965 = base.Simd_g_v128_load64_zero(m, v325, v1953)
			v1967 = base.Simd_g_i8x16_shuffle2(v1921, v1965, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v1968 = base.Simd_g_i32x4_extend_low_i16x8_u(v1967)
			v1969 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k12)
			v1971 = base.Simd_g_i32x4_extend_high_i16x8_u(v1967)
			v1976 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k13)
			v1978 = int32(6)
			v1980 = int32(8)
			v1981 = base.Simd_g_v128_load64_zero(m, v1918, v1980)
			v1983 = base.Simd_g_i8x16_shuffle2(v1921, v1981, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v1989 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1983), v1958), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1983), v1958), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v1991 = base.Simd_g_v128_load64_zero(m, v325, v1980)
			v1993 = base.Simd_g_i8x16_shuffle2(v1921, v1991, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v1994 = base.Simd_g_i32x4_extend_low_i16x8_u(v1993)
			v1996 = base.Simd_g_i32x4_extend_high_i16x8_u(v1993)
			v2004 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1963, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1968, v1969), base.Simd_g_i32x4_mul(v1971, v1969), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v1976), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1989, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1994, v1969), base.Simd_g_i32x4_mul(v1996, v1969), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v1976), v1978))
			v2005 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k14)
			v2007 = int32(16)
			v2008 = base.Simd_g_v128_load64_zero(m, v1918, v2007)
			v2010 = base.Simd_g_i8x16_shuffle2(v1921, v2008, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2016 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2010), v1958), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2010), v1958), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2018 = base.Simd_g_v128_load64_zero(m, v325, v2007)
			v2020 = base.Simd_g_i8x16_shuffle2(v1921, v2018, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2021 = base.Simd_g_i32x4_extend_low_i16x8_u(v2020)
			v2023 = base.Simd_g_i32x4_extend_high_i16x8_u(v2020)
			v2031 = int32(24)
			v2032 = base.Simd_g_v128_load64_zero(m, v1918, v2031)
			v2034 = base.Simd_g_i8x16_shuffle2(v1921, v2032, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2040 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2034), v1958), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2034), v1958), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))
			v2042 = base.Simd_g_v128_load64_zero(m, v325, v2031)
			v2044 = base.Simd_g_i8x16_shuffle2(v1921, v2042, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2045 = base.Simd_g_i32x4_extend_low_i16x8_u(v2044)
			v2047 = base.Simd_g_i32x4_extend_high_i16x8_u(v2044)
			v2055 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2016, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2021, v1969), base.Simd_g_i32x4_mul(v2023, v1969), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v1976), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2040, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2045, v1969), base.Simd_g_i32x4_mul(v2047, v1969), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11))), v1976), v1978))
			v2057 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2004, v2005), base.Simd_g_v128_and(v2055, v2005))
			v2061 = base.Simd_g_v128_load64_zero(m, v327, v1953)
			v2063 = base.Simd_g_i8x16_shuffle2(v1921, v2061, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2064 = base.Simd_g_i32x4_extend_low_i16x8_u(v2063)
			v2065 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k15)
			v2067 = base.Simd_g_i32x4_extend_high_i16x8_u(v2063)
			v2072 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k16)
			v2077 = base.Simd_g_v128_load64_zero(m, v327, v1980)
			v2079 = base.Simd_g_i8x16_shuffle2(v1921, v2077, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2080 = base.Simd_g_i32x4_extend_low_i16x8_u(v2079)
			v2082 = base.Simd_g_i32x4_extend_high_i16x8_u(v2079)
			v2090 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2064, v2065), base.Simd_g_i32x4_mul(v2067, v2065), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v1963), v2072), v1978), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2080, v2065), base.Simd_g_i32x4_mul(v2082, v2065), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v1989), v2072), v1978))
			v2094 = base.Simd_g_v128_load64_zero(m, v327, v2007)
			v2096 = base.Simd_g_i8x16_shuffle2(v1921, v2094, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2097 = base.Simd_g_i32x4_extend_low_i16x8_u(v2096)
			v2099 = base.Simd_g_i32x4_extend_high_i16x8_u(v2096)
			v2108 = base.Simd_g_v128_load64_zero(m, v327, v2031)
			v2110 = base.Simd_g_i8x16_shuffle2(v1921, v2108, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k6))
			v2111 = base.Simd_g_i32x4_extend_low_i16x8_u(v2110)
			v2113 = base.Simd_g_i32x4_extend_high_i16x8_u(v2110)
			v2121 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2097, v2065), base.Simd_g_i32x4_mul(v2099, v2065), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2016), v2072), v1978), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2111, v2065), base.Simd_g_i32x4_mul(v2113, v2065), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), v2040), v2072), v1978))
			v2124 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2090, v1980), base.Simd_g_i16x8_shr_u(v2121, v1980))
			v2127 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2057, v1980), base.Simd_g_i16x8_shr_u(v2124, v1980))
			v2129 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k17)
			v2134 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k18)
			v2141 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k19)
			v2158 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1963, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2064, v2129), base.Simd_g_i32x4_mul(v2067, v2129), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1968, v2134), base.Simd_g_i32x4_mul(v1971, v2134), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2141), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1989, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2080, v2129), base.Simd_g_i32x4_mul(v2082, v2129), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1994, v2134), base.Simd_g_i32x4_mul(v1996, v2134), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2141), v1978))
			v2187 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2016, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2097, v2129), base.Simd_g_i32x4_mul(v2099, v2129), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2021, v2134), base.Simd_g_i32x4_mul(v2023, v2134), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2141), v1978), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2040, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2111, v2129), base.Simd_g_i32x4_mul(v2113, v2129), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2045, v2134), base.Simd_g_i32x4_mul(v2047, v2134), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE2__k11)))), v2141), v1978))
			v2190 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2158, v1980), base.Simd_g_i16x8_shr_u(v2187, v1980))
			v2197 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2004, v1980), base.Simd_g_i16x8_shr_u(v2055, v1980))
			v2200 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2190, v1980), base.Simd_g_i16x8_shr_u(v2197, v1980))
			v2202 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2127, v2005), base.Simd_g_v128_and(v2200, v2005))
			v2207 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2090, v2005), base.Simd_g_v128_and(v2121, v2005))
			v2211 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2158, v2005), base.Simd_g_v128_and(v2187, v2005))
			v2213 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2207, v2005), base.Simd_g_v128_and(v2211, v2005))
			v2218 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2057, v2005), base.Simd_g_v128_and(v2124, v2005))
			v2221 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2213, v1980), base.Simd_g_i16x8_shr_u(v2218, v1980))
			v2224 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2202, v1980), base.Simd_g_i16x8_shr_u(v2221, v1980))
			v2229 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2190, v2005), base.Simd_g_v128_and(v2197, v2005))
			v2236 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2207, v1980), base.Simd_g_i16x8_shr_u(v2211, v1980))
			v2239 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2229, v1980), base.Simd_g_i16x8_shr_u(v2236, v1980))
			v2246 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2127, v1980), base.Simd_g_i16x8_shr_u(v2200, v1980))
			v2249 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2239, v1980), base.Simd_g_i16x8_shr_u(v2246, v1980))
			v2252 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2224, v1980), base.Simd_g_i16x8_shr_u(v2249, v1980))
			base.Simd_g_v128_store(m, v1920, int32(80), v2252)
			v2257 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2239, v2005), base.Simd_g_v128_and(v2246, v2005))
			v2262 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2213, v2005), base.Simd_g_v128_and(v2218, v2005))
			v2267 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2229, v2005), base.Simd_g_v128_and(v2236, v2005))
			v2270 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2262, v1980), base.Simd_g_i16x8_shr_u(v2267, v1980))
			v2273 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2257, v1980), base.Simd_g_i16x8_shr_u(v2270, v1980))
			base.Simd_g_v128_store(m, v1920, int32(64), v2273)
			v2278 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2262, v2005), base.Simd_g_v128_and(v2267, v2005))
			v2283 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2202, v2005), base.Simd_g_v128_and(v2221, v2005))
			v2286 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(v2278, v1980), base.Simd_g_i16x8_shr_u(v2283, v1980))
			base.Simd_g_v128_store(m, v1920, int32(48), v2286)
			v2291 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2224, v2005), base.Simd_g_v128_and(v2249, v2005))
			base.Simd_g_v128_store(m, v1920, int32(32), v2291)
			v2296 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2257, v2005), base.Simd_g_v128_and(v2270, v2005))
			base.Simd_g_v128_store(m, v1920, v2007, v2296)
			v2301 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_v128_and(v2278, v2005), base.Simd_g_v128_and(v2283, v2005))
			base.Simd_g_v128_store(m, v1920, v1953, v2301)
			v2304 = int32(3)
			v2309 = F_memcpy(m, l6+v1275*v2304, v1920, v1917*v2304)
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleBgrLinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleBgrLinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleBgrLinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleBgrLinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleBgrLinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleBgrLinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleBgrLinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleBgrLinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleBgrLinePair_SSE2__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleBgrLinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleBgrLinePair_SSE2__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleBgrLinePair_SSE2__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleBgrLinePair_SSE2__k12 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleBgrLinePair_SSE2__k13 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleBgrLinePair_SSE2__k14 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_UpsampleBgrLinePair_SSE2__k15 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleBgrLinePair_SSE2__k16 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleBgrLinePair_SSE2__k17 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleBgrLinePair_SSE2__k18 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleBgrLinePair_SSE2__k19 = [2]uint64{0x2204220422042204, 0x2204220422042204}

func F_UpsampleBgrLinePair_SSE41(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v175 int32
	_ = v175
	var v177 int32
	_ = v177
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v191 int32
	_ = v191
	var v194 int32
	_ = v194
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v199 int32
	_ = v199
	var v202 int32
	_ = v202
	var v205 int32
	_ = v205
	var v210 int32
	_ = v210
	var v212 int32
	_ = v212
	var v219 int32
	_ = v219
	var v222 int32
	_ = v222
	var v226 int32
	_ = v226
	var v233 int32
	_ = v233
	var v235 int32
	_ = v235
	var v242 int32
	_ = v242
	var v245 int32
	_ = v245
	var v249 int32
	_ = v249
	var v252 int32
	_ = v252
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v261 int32
	_ = v261
	var v263 int32
	_ = v263
	var v270 int32
	_ = v270
	var v273 int32
	_ = v273
	var v277 int32
	_ = v277
	var v282 int32
	_ = v282
	var v284 int32
	_ = v284
	var v291 int32
	_ = v291
	var v294 int32
	_ = v294
	var v298 int32
	_ = v298
	var v305 int32
	_ = v305
	var v307 int32
	_ = v307
	var v314 int32
	_ = v314
	var v317 int32
	_ = v317
	var v325 int32
	_ = v325
	var v327 int32
	_ = v327
	var v330 int32
	_ = v330
	var v332 int32
	_ = v332
	var v339 int32
	_ = v339
	var v341 int32
	_ = v341
	var v342 int32
	_ = v342
	var v354 int32
	_ = v354
	var v355 int32
	_ = v355
	var v359 int32
	_ = v359
	var v360 int32
	_ = v360
	var v376 int32
	_ = v376
	var v377 int32
	_ = v377
	var v378 base.V128
	_ = v378
	var v379 int32
	_ = v379
	var v381 base.V128
	_ = v381
	var v382 int32
	_ = v382
	var v383 base.V128
	_ = v383
	var v384 base.V128
	_ = v384
	var v386 base.V128
	_ = v386
	var v387 base.V128
	_ = v387
	var v389 base.V128
	_ = v389
	var v390 base.V128
	_ = v390
	var v392 base.V128
	_ = v392
	var v394 base.V128
	_ = v394
	var v396 base.V128
	_ = v396
	var v402 base.V128
	_ = v402
	var v403 base.V128
	_ = v403
	var v409 base.V128
	_ = v409
	var v410 base.V128
	_ = v410
	var v411 base.V128
	_ = v411
	var v412 base.V128
	_ = v412
	var v413 int32
	_ = v413
	var v415 base.V128
	_ = v415
	var v416 base.V128
	_ = v416
	var v417 int32
	_ = v417
	var v419 base.V128
	_ = v419
	var v420 base.V128
	_ = v420
	var v422 base.V128
	_ = v422
	var v423 int32
	_ = v423
	var v426 base.V128
	_ = v426
	var v429 int32
	_ = v429
	var v431 base.V128
	_ = v431
	var v432 int32
	_ = v432
	var v434 base.V128
	_ = v434
	var v436 base.V128
	_ = v436
	var v437 base.V128
	_ = v437
	var v439 base.V128
	_ = v439
	var v440 base.V128
	_ = v440
	var v442 base.V128
	_ = v442
	var v443 base.V128
	_ = v443
	var v445 base.V128
	_ = v445
	var v448 base.V128
	_ = v448
	var v454 base.V128
	_ = v454
	var v455 base.V128
	_ = v455
	var v461 base.V128
	_ = v461
	var v462 base.V128
	_ = v462
	var v464 base.V128
	_ = v464
	var v468 base.V128
	_ = v468
	var v471 base.V128
	_ = v471
	var v472 base.V128
	_ = v472
	var v474 base.V128
	_ = v474
	var v475 int32
	_ = v475
	var v478 base.V128
	_ = v478
	var v479 int32
	_ = v479
	var v481 int32
	_ = v481
	var v482 base.V128
	_ = v482
	var v511 base.V128
	_ = v511
	var v513 base.V128
	_ = v513
	var v515 base.V128
	_ = v515
	var v519 base.V128
	_ = v519
	var v520 base.V128
	_ = v520
	var v522 base.V128
	_ = v522
	var v524 base.V128
	_ = v524
	var v525 base.V128
	_ = v525
	var v526 base.V128
	_ = v526
	var v528 base.V128
	_ = v528
	var v533 base.V128
	_ = v533
	var v535 base.V128
	_ = v535
	var v536 base.V128
	_ = v536
	var v537 base.V128
	_ = v537
	var v539 base.V128
	_ = v539
	var v545 base.V128
	_ = v545
	var v547 int32
	_ = v547
	var v549 int32
	_ = v549
	var v550 base.V128
	_ = v550
	var v552 base.V128
	_ = v552
	var v558 base.V128
	_ = v558
	var v560 base.V128
	_ = v560
	var v562 base.V128
	_ = v562
	var v563 base.V128
	_ = v563
	var v565 base.V128
	_ = v565
	var v570 base.V128
	_ = v570
	var v572 base.V128
	_ = v572
	var v573 base.V128
	_ = v573
	var v575 base.V128
	_ = v575
	var v584 base.V128
	_ = v584
	var v585 base.V128
	_ = v585
	var v587 base.V128
	_ = v587
	var v593 base.V128
	_ = v593
	var v605 base.V128
	_ = v605
	var v606 base.V128
	_ = v606
	var v609 base.V128
	_ = v609
	var v615 base.V128
	_ = v615
	var v627 base.V128
	_ = v627
	var v628 base.V128
	_ = v628
	var v630 base.V128
	_ = v630
	var v633 base.V128
	_ = v633
	var v635 base.V128
	_ = v635
	var v638 base.V128
	_ = v638
	var v640 base.V128
	_ = v640
	var v643 base.V128
	_ = v643
	var v645 base.V128
	_ = v645
	var v648 base.V128
	_ = v648
	var v650 base.V128
	_ = v650
	var v654 base.V128
	_ = v654
	var v656 base.V128
	_ = v656
	var v662 base.V128
	_ = v662
	var v664 base.V128
	_ = v664
	var v666 base.V128
	_ = v666
	var v667 base.V128
	_ = v667
	var v669 base.V128
	_ = v669
	var v674 base.V128
	_ = v674
	var v676 base.V128
	_ = v676
	var v677 base.V128
	_ = v677
	var v679 base.V128
	_ = v679
	var v688 int32
	_ = v688
	var v689 base.V128
	_ = v689
	var v691 base.V128
	_ = v691
	var v697 base.V128
	_ = v697
	var v699 base.V128
	_ = v699
	var v701 base.V128
	_ = v701
	var v702 base.V128
	_ = v702
	var v704 base.V128
	_ = v704
	var v709 base.V128
	_ = v709
	var v711 base.V128
	_ = v711
	var v712 base.V128
	_ = v712
	var v714 base.V128
	_ = v714
	var v723 base.V128
	_ = v723
	var v741 base.V128
	_ = v741
	var v760 base.V128
	_ = v760
	var v762 base.V128
	_ = v762
	var v769 base.V128
	_ = v769
	var v776 base.V128
	_ = v776
	var v781 int32
	_ = v781
	var v782 base.V128
	_ = v782
	var v810 int32
	_ = v810
	var v811 base.V128
	_ = v811
	var v812 base.V128
	_ = v812
	var v813 base.V128
	_ = v813
	var v815 base.V128
	_ = v815
	var v819 base.V128
	_ = v819
	var v820 base.V128
	_ = v820
	var v822 base.V128
	_ = v822
	var v824 base.V128
	_ = v824
	var v825 base.V128
	_ = v825
	var v826 base.V128
	_ = v826
	var v828 base.V128
	_ = v828
	var v833 base.V128
	_ = v833
	var v835 base.V128
	_ = v835
	var v836 base.V128
	_ = v836
	var v837 base.V128
	_ = v837
	var v839 base.V128
	_ = v839
	var v845 base.V128
	_ = v845
	var v847 int32
	_ = v847
	var v849 int32
	_ = v849
	var v850 base.V128
	_ = v850
	var v852 base.V128
	_ = v852
	var v858 base.V128
	_ = v858
	var v860 base.V128
	_ = v860
	var v862 base.V128
	_ = v862
	var v863 base.V128
	_ = v863
	var v865 base.V128
	_ = v865
	var v870 base.V128
	_ = v870
	var v872 base.V128
	_ = v872
	var v873 base.V128
	_ = v873
	var v875 base.V128
	_ = v875
	var v884 base.V128
	_ = v884
	var v885 base.V128
	_ = v885
	var v887 base.V128
	_ = v887
	var v893 base.V128
	_ = v893
	var v905 base.V128
	_ = v905
	var v906 base.V128
	_ = v906
	var v909 base.V128
	_ = v909
	var v915 base.V128
	_ = v915
	var v927 base.V128
	_ = v927
	var v928 base.V128
	_ = v928
	var v930 base.V128
	_ = v930
	var v933 base.V128
	_ = v933
	var v935 base.V128
	_ = v935
	var v938 base.V128
	_ = v938
	var v940 base.V128
	_ = v940
	var v943 base.V128
	_ = v943
	var v945 base.V128
	_ = v945
	var v948 base.V128
	_ = v948
	var v950 base.V128
	_ = v950
	var v953 int32
	_ = v953
	var v954 base.V128
	_ = v954
	var v956 base.V128
	_ = v956
	var v962 base.V128
	_ = v962
	var v964 base.V128
	_ = v964
	var v966 base.V128
	_ = v966
	var v967 base.V128
	_ = v967
	var v969 base.V128
	_ = v969
	var v974 base.V128
	_ = v974
	var v976 base.V128
	_ = v976
	var v977 base.V128
	_ = v977
	var v979 base.V128
	_ = v979
	var v988 int32
	_ = v988
	var v989 base.V128
	_ = v989
	var v991 base.V128
	_ = v991
	var v997 base.V128
	_ = v997
	var v999 base.V128
	_ = v999
	var v1001 base.V128
	_ = v1001
	var v1002 base.V128
	_ = v1002
	var v1004 base.V128
	_ = v1004
	var v1009 base.V128
	_ = v1009
	var v1011 base.V128
	_ = v1011
	var v1012 base.V128
	_ = v1012
	var v1014 base.V128
	_ = v1014
	var v1023 base.V128
	_ = v1023
	var v1041 base.V128
	_ = v1041
	var v1060 base.V128
	_ = v1060
	var v1062 base.V128
	_ = v1062
	var v1069 base.V128
	_ = v1069
	var v1076 base.V128
	_ = v1076
	var v1079 int32
	_ = v1079
	var v1084 int32
	_ = v1084
	var v1102 int32
	_ = v1102
	var v1103 int32
	_ = v1103
	var v1126 int32
	_ = v1126
	var v1129 int32
	_ = v1129
	var v1135 int32
	_ = v1135
	var v1136 int32
	_ = v1136
	var v1138 int32
	_ = v1138
	var v1140 int32
	_ = v1140
	var v1141 int32
	_ = v1141
	var v1145 int32
	_ = v1145
	var v1146 int32
	_ = v1146
	var v1147 int32
	_ = v1147
	var v1149 int32
	_ = v1149
	var v1160 int32
	_ = v1160
	var v1185 int32
	_ = v1185
	var v1186 int32
	_ = v1186
	var v1190 int32
	_ = v1190
	var v1194 int32
	_ = v1194
	var v1195 int32
	_ = v1195
	var v1230 int32
	_ = v1230
	var v1231 int32
	_ = v1231
	var v1236 int64
	_ = v1236
	var v1239 int32
	_ = v1239
	var v1240 int32
	_ = v1240
	var v1252 int32
	_ = v1252
	var v1270 int32
	_ = v1270
	var v1271 int32
	_ = v1271
	var v1272 int32
	_ = v1272
	var v1283 int32
	_ = v1283
	var v1308 int32
	_ = v1308
	var v1309 int32
	_ = v1309
	var v1313 int32
	_ = v1313
	var v1317 int32
	_ = v1317
	var v1318 int32
	_ = v1318
	var v1353 int32
	_ = v1353
	var v1354 int32
	_ = v1354
	var v1359 int64
	_ = v1359
	var v1362 int32
	_ = v1362
	var v1363 int32
	_ = v1363
	var v1375 int32
	_ = v1375
	var v1393 int32
	_ = v1393
	var v1394 base.V128
	_ = v1394
	var v1395 int32
	_ = v1395
	var v1396 base.V128
	_ = v1396
	var v1398 base.V128
	_ = v1398
	var v1399 base.V128
	_ = v1399
	var v1400 int32
	_ = v1400
	var v1401 base.V128
	_ = v1401
	var v1402 base.V128
	_ = v1402
	var v1404 base.V128
	_ = v1404
	var v1405 base.V128
	_ = v1405
	var v1407 base.V128
	_ = v1407
	var v1409 base.V128
	_ = v1409
	var v1411 base.V128
	_ = v1411
	var v1417 base.V128
	_ = v1417
	var v1418 base.V128
	_ = v1418
	var v1424 base.V128
	_ = v1424
	var v1425 base.V128
	_ = v1425
	var v1426 base.V128
	_ = v1426
	var v1427 base.V128
	_ = v1427
	var v1430 base.V128
	_ = v1430
	var v1431 base.V128
	_ = v1431
	var v1434 base.V128
	_ = v1434
	var v1435 base.V128
	_ = v1435
	var v1437 base.V128
	_ = v1437
	var v1441 base.V128
	_ = v1441
	var v1447 int32
	_ = v1447
	var v1449 int32
	_ = v1449
	var v1450 int32
	_ = v1450
	var v1461 int32
	_ = v1461
	var v1486 int32
	_ = v1486
	var v1487 int32
	_ = v1487
	var v1491 int32
	_ = v1491
	var v1495 int32
	_ = v1495
	var v1496 int32
	_ = v1496
	var v1531 int32
	_ = v1531
	var v1532 int32
	_ = v1532
	var v1537 int64
	_ = v1537
	var v1540 int32
	_ = v1540
	var v1541 int32
	_ = v1541
	var v1553 int32
	_ = v1553
	var v1571 int32
	_ = v1571
	var v1582 int32
	_ = v1582
	var v1607 int32
	_ = v1607
	var v1608 int32
	_ = v1608
	var v1612 int32
	_ = v1612
	var v1616 int32
	_ = v1616
	var v1617 int32
	_ = v1617
	var v1652 int32
	_ = v1652
	var v1653 int32
	_ = v1653
	var v1658 int64
	_ = v1658
	var v1661 int32
	_ = v1661
	var v1662 int32
	_ = v1662
	var v1674 int32
	_ = v1674
	var v1693 base.V128
	_ = v1693
	var v1694 int32
	_ = v1694
	var v1695 base.V128
	_ = v1695
	var v1697 base.V128
	_ = v1697
	var v1698 base.V128
	_ = v1698
	var v1700 base.V128
	_ = v1700
	var v1701 base.V128
	_ = v1701
	var v1703 base.V128
	_ = v1703
	var v1704 base.V128
	_ = v1704
	var v1706 base.V128
	_ = v1706
	var v1709 base.V128
	_ = v1709
	var v1715 base.V128
	_ = v1715
	var v1716 base.V128
	_ = v1716
	var v1722 base.V128
	_ = v1722
	var v1723 base.V128
	_ = v1723
	var v1724 base.V128
	_ = v1724
	var v1725 base.V128
	_ = v1725
	var v1728 base.V128
	_ = v1728
	var v1729 base.V128
	_ = v1729
	var v1732 base.V128
	_ = v1732
	var v1733 base.V128
	_ = v1733
	var v1735 base.V128
	_ = v1735
	var v1739 base.V128
	_ = v1739
	var v1745 int32
	_ = v1745
	var v1746 int32
	_ = v1746
	var v1748 int32
	_ = v1748
	var v1749 base.V128
	_ = v1749
	var v1777 int32
	_ = v1777
	var v1778 base.V128
	_ = v1778
	var v1779 base.V128
	_ = v1779
	var v1780 base.V128
	_ = v1780
	var v1782 base.V128
	_ = v1782
	var v1786 base.V128
	_ = v1786
	var v1787 base.V128
	_ = v1787
	var v1789 base.V128
	_ = v1789
	var v1791 base.V128
	_ = v1791
	var v1792 base.V128
	_ = v1792
	var v1793 base.V128
	_ = v1793
	var v1795 base.V128
	_ = v1795
	var v1800 base.V128
	_ = v1800
	var v1802 base.V128
	_ = v1802
	var v1803 base.V128
	_ = v1803
	var v1804 base.V128
	_ = v1804
	var v1806 base.V128
	_ = v1806
	var v1812 base.V128
	_ = v1812
	var v1814 int32
	_ = v1814
	var v1816 int32
	_ = v1816
	var v1817 base.V128
	_ = v1817
	var v1819 base.V128
	_ = v1819
	var v1825 base.V128
	_ = v1825
	var v1827 base.V128
	_ = v1827
	var v1829 base.V128
	_ = v1829
	var v1830 base.V128
	_ = v1830
	var v1832 base.V128
	_ = v1832
	var v1837 base.V128
	_ = v1837
	var v1839 base.V128
	_ = v1839
	var v1840 base.V128
	_ = v1840
	var v1842 base.V128
	_ = v1842
	var v1851 base.V128
	_ = v1851
	var v1852 base.V128
	_ = v1852
	var v1854 base.V128
	_ = v1854
	var v1860 base.V128
	_ = v1860
	var v1872 base.V128
	_ = v1872
	var v1873 base.V128
	_ = v1873
	var v1876 base.V128
	_ = v1876
	var v1882 base.V128
	_ = v1882
	var v1894 base.V128
	_ = v1894
	var v1895 base.V128
	_ = v1895
	var v1897 base.V128
	_ = v1897
	var v1900 base.V128
	_ = v1900
	var v1902 base.V128
	_ = v1902
	var v1905 base.V128
	_ = v1905
	var v1907 base.V128
	_ = v1907
	var v1910 base.V128
	_ = v1910
	var v1912 base.V128
	_ = v1912
	var v1915 base.V128
	_ = v1915
	var v1917 base.V128
	_ = v1917
	var v1920 int32
	_ = v1920
	var v1921 base.V128
	_ = v1921
	var v1923 base.V128
	_ = v1923
	var v1929 base.V128
	_ = v1929
	var v1931 base.V128
	_ = v1931
	var v1933 base.V128
	_ = v1933
	var v1934 base.V128
	_ = v1934
	var v1936 base.V128
	_ = v1936
	var v1941 base.V128
	_ = v1941
	var v1943 base.V128
	_ = v1943
	var v1944 base.V128
	_ = v1944
	var v1946 base.V128
	_ = v1946
	var v1955 int32
	_ = v1955
	var v1956 base.V128
	_ = v1956
	var v1958 base.V128
	_ = v1958
	var v1964 base.V128
	_ = v1964
	var v1966 base.V128
	_ = v1966
	var v1968 base.V128
	_ = v1968
	var v1969 base.V128
	_ = v1969
	var v1971 base.V128
	_ = v1971
	var v1976 base.V128
	_ = v1976
	var v1978 base.V128
	_ = v1978
	var v1979 base.V128
	_ = v1979
	var v1981 base.V128
	_ = v1981
	var v1990 base.V128
	_ = v1990
	var v2008 base.V128
	_ = v2008
	var v2027 base.V128
	_ = v2027
	var v2029 base.V128
	_ = v2029
	var v2036 base.V128
	_ = v2036
	var v2043 base.V128
	_ = v2043
	var v2046 int32
	_ = v2046
	var v2051 int32
	_ = v2051
	var v2055 int32
	_ = v2055
	var v2056 base.V128
	_ = v2056
	var v2084 int32
	_ = v2084
	var v2085 base.V128
	_ = v2085
	var v2086 base.V128
	_ = v2086
	var v2087 base.V128
	_ = v2087
	var v2089 base.V128
	_ = v2089
	var v2093 base.V128
	_ = v2093
	var v2094 base.V128
	_ = v2094
	var v2096 base.V128
	_ = v2096
	var v2098 base.V128
	_ = v2098
	var v2099 base.V128
	_ = v2099
	var v2100 base.V128
	_ = v2100
	var v2102 base.V128
	_ = v2102
	var v2107 base.V128
	_ = v2107
	var v2109 base.V128
	_ = v2109
	var v2110 base.V128
	_ = v2110
	var v2111 base.V128
	_ = v2111
	var v2113 base.V128
	_ = v2113
	var v2119 base.V128
	_ = v2119
	var v2121 int32
	_ = v2121
	var v2123 int32
	_ = v2123
	var v2124 base.V128
	_ = v2124
	var v2126 base.V128
	_ = v2126
	var v2132 base.V128
	_ = v2132
	var v2134 base.V128
	_ = v2134
	var v2136 base.V128
	_ = v2136
	var v2137 base.V128
	_ = v2137
	var v2139 base.V128
	_ = v2139
	var v2144 base.V128
	_ = v2144
	var v2146 base.V128
	_ = v2146
	var v2147 base.V128
	_ = v2147
	var v2149 base.V128
	_ = v2149
	var v2158 base.V128
	_ = v2158
	var v2159 base.V128
	_ = v2159
	var v2161 base.V128
	_ = v2161
	var v2167 base.V128
	_ = v2167
	var v2179 base.V128
	_ = v2179
	var v2180 base.V128
	_ = v2180
	var v2183 base.V128
	_ = v2183
	var v2189 base.V128
	_ = v2189
	var v2201 base.V128
	_ = v2201
	var v2202 base.V128
	_ = v2202
	var v2204 base.V128
	_ = v2204
	var v2207 base.V128
	_ = v2207
	var v2209 base.V128
	_ = v2209
	var v2212 base.V128
	_ = v2212
	var v2214 base.V128
	_ = v2214
	var v2217 base.V128
	_ = v2217
	var v2219 base.V128
	_ = v2219
	var v2222 base.V128
	_ = v2222
	var v2224 base.V128
	_ = v2224
	var v2227 int32
	_ = v2227
	var v2228 base.V128
	_ = v2228
	var v2230 base.V128
	_ = v2230
	var v2236 base.V128
	_ = v2236
	var v2238 base.V128
	_ = v2238
	var v2240 base.V128
	_ = v2240
	var v2241 base.V128
	_ = v2241
	var v2243 base.V128
	_ = v2243
	var v2248 base.V128
	_ = v2248
	var v2250 base.V128
	_ = v2250
	var v2251 base.V128
	_ = v2251
	var v2253 base.V128
	_ = v2253
	var v2262 int32
	_ = v2262
	var v2263 base.V128
	_ = v2263
	var v2265 base.V128
	_ = v2265
	var v2271 base.V128
	_ = v2271
	var v2273 base.V128
	_ = v2273
	var v2275 base.V128
	_ = v2275
	var v2276 base.V128
	_ = v2276
	var v2278 base.V128
	_ = v2278
	var v2283 base.V128
	_ = v2283
	var v2285 base.V128
	_ = v2285
	var v2286 base.V128
	_ = v2286
	var v2288 base.V128
	_ = v2288
	var v2297 base.V128
	_ = v2297
	var v2315 base.V128
	_ = v2315
	var v2334 base.V128
	_ = v2334
	var v2336 base.V128
	_ = v2336
	var v2343 base.V128
	_ = v2343
	var v2350 base.V128
	_ = v2350
	var v2354 int32
	_ = v2354
	var v2356 int32
	_ = v2356
	var v2358 int32
	_ = v2358
	var v2359 base.V128
	_ = v2359
	var v2387 int32
	_ = v2387
	var v2388 base.V128
	_ = v2388
	var v2389 base.V128
	_ = v2389
	var v2390 base.V128
	_ = v2390
	var v2392 base.V128
	_ = v2392
	var v2396 base.V128
	_ = v2396
	var v2397 base.V128
	_ = v2397
	var v2399 base.V128
	_ = v2399
	var v2401 base.V128
	_ = v2401
	var v2402 base.V128
	_ = v2402
	var v2403 base.V128
	_ = v2403
	var v2405 base.V128
	_ = v2405
	var v2410 base.V128
	_ = v2410
	var v2412 base.V128
	_ = v2412
	var v2413 base.V128
	_ = v2413
	var v2414 base.V128
	_ = v2414
	var v2416 base.V128
	_ = v2416
	var v2422 base.V128
	_ = v2422
	var v2424 int32
	_ = v2424
	var v2426 int32
	_ = v2426
	var v2427 base.V128
	_ = v2427
	var v2429 base.V128
	_ = v2429
	var v2435 base.V128
	_ = v2435
	var v2437 base.V128
	_ = v2437
	var v2439 base.V128
	_ = v2439
	var v2440 base.V128
	_ = v2440
	var v2442 base.V128
	_ = v2442
	var v2447 base.V128
	_ = v2447
	var v2449 base.V128
	_ = v2449
	var v2450 base.V128
	_ = v2450
	var v2452 base.V128
	_ = v2452
	var v2461 base.V128
	_ = v2461
	var v2462 base.V128
	_ = v2462
	var v2464 base.V128
	_ = v2464
	var v2470 base.V128
	_ = v2470
	var v2482 base.V128
	_ = v2482
	var v2483 base.V128
	_ = v2483
	var v2486 base.V128
	_ = v2486
	var v2492 base.V128
	_ = v2492
	var v2504 base.V128
	_ = v2504
	var v2505 base.V128
	_ = v2505
	var v2507 base.V128
	_ = v2507
	var v2510 base.V128
	_ = v2510
	var v2512 base.V128
	_ = v2512
	var v2515 base.V128
	_ = v2515
	var v2517 base.V128
	_ = v2517
	var v2520 base.V128
	_ = v2520
	var v2522 base.V128
	_ = v2522
	var v2525 base.V128
	_ = v2525
	var v2527 base.V128
	_ = v2527
	var v2530 int32
	_ = v2530
	var v2531 base.V128
	_ = v2531
	var v2533 base.V128
	_ = v2533
	var v2539 base.V128
	_ = v2539
	var v2541 base.V128
	_ = v2541
	var v2543 base.V128
	_ = v2543
	var v2544 base.V128
	_ = v2544
	var v2546 base.V128
	_ = v2546
	var v2551 base.V128
	_ = v2551
	var v2553 base.V128
	_ = v2553
	var v2554 base.V128
	_ = v2554
	var v2556 base.V128
	_ = v2556
	var v2565 int32
	_ = v2565
	var v2566 base.V128
	_ = v2566
	var v2568 base.V128
	_ = v2568
	var v2574 base.V128
	_ = v2574
	var v2576 base.V128
	_ = v2576
	var v2578 base.V128
	_ = v2578
	var v2579 base.V128
	_ = v2579
	var v2581 base.V128
	_ = v2581
	var v2586 base.V128
	_ = v2586
	var v2588 base.V128
	_ = v2588
	var v2589 base.V128
	_ = v2589
	var v2591 base.V128
	_ = v2591
	var v2600 base.V128
	_ = v2600
	var v2618 base.V128
	_ = v2618
	var v2637 base.V128
	_ = v2637
	var v2639 base.V128
	_ = v2639
	var v2646 base.V128
	_ = v2646
	var v2653 base.V128
	_ = v2653
	var v2656 int32
	_ = v2656
	var v2657 int32
	_ = v2657
	var v2660 int32
	_ = v2660
	var v2661 int32
	_ = v2661
	var v2663 int32
	_ = v2663
	v10 = int32(0)
	v33 = m.G0
	v35 = v33 - int32(528)
	m.G0 = v35
	base.MemoryFill(m, v35+int32(64), v10, int32(463))
	v162 = int32(1)
	v163 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v164 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v169 = int32(base.Ui32(v163+v164)>>(uint(v162)%32)) + v162
	v172 = int32(base.Ui32(v169+v164) >> (uint(v162) % 32))
	v175 = int32(8)
	v177 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v181 = int32(base.Ui32(v177*int32(_a_F_UpsampleBgrLinePair_SSE41_0)) >> (uint(v175) % 32))
	v182 = int32(base.Ui32(v172*int32(_a_F_UpsampleBgrLinePair_SSE41_1))>>(uint(v175)%32)) + v181
	v184 = v182 + int32(-14234)
	if base.Ui32(v182) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_2)) {
		v191 = int32(0)
	} else {
		v191 = int32(255)
	}
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_3)) {
		v194 = int32(base.Ui32(v184) >> (uint(int32(6)) % 32))
	} else {
		v194 = v191
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v194)
	v196 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v197 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v199 = int32(1)
	v202 = int32(base.Ui32(v196+v197)>>(uint(v199)%32)) + v199
	v205 = int32(base.Ui32(v202+v197) >> (uint(v199) % 32))
	v210 = int32(base.Ui32(v205*int32(_a_F_UpsampleBgrLinePair_SSE41_4))>>(uint(int32(8))%32)) + v181
	v212 = v210 + int32(-17685)
	if base.Ui32(v210) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_5)) {
		v219 = int32(0)
	} else {
		v219 = int32(255)
	}
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_3)) {
		v222 = int32(base.Ui32(v212) >> (uint(int32(6)) % 32))
	} else {
		v222 = v219
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v222)
	v226 = int32(8)
	v233 = v181 - (int32(base.Ui32(v205*int32(_a_F_UpsampleBgrLinePair_SSE41_6))>>(uint(v226)%32)) + int32(base.Ui32(v172*int32(_a_F_UpsampleBgrLinePair_SSE41_7))>>(uint(v226)%32)))
	v235 = v233 + int32(_a_F_UpsampleBgrLinePair_SSE41_8)
	if v233 < int32(-8708) {
		v242 = int32(0)
	} else {
		v242 = int32(255)
	}
	if base.Ui32(v235) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_3)) {
		v245 = int32(base.Ui32(v235) >> (uint(int32(6)) % 32))
	} else {
		v245 = v242
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v245)
	if l1 == int32(0) {
	} else {
		v249 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v252 = int32(8)
		v253 = int32(base.Ui32(v249*int32(_a_F_UpsampleBgrLinePair_SSE41_0)) >> (uint(v252) % 32))
		v256 = int32(base.Ui32(v169+v163) >> (uint(int32(1)) % 32))
		v261 = v253 + int32(base.Ui32(v256*int32(_a_F_UpsampleBgrLinePair_SSE41_1))>>(uint(v252)%32))
		v263 = v261 + int32(-14234)
		if base.Ui32(v261) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_2)) {
			v270 = int32(0)
		} else {
			v270 = int32(255)
		}
		if base.Ui32(v263) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_3)) {
			v273 = int32(base.Ui32(v263) >> (uint(int32(6)) % 32))
		} else {
			v273 = v270
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v273)
		v277 = int32(base.Ui32(v202+v196) >> (uint(int32(1)) % 32))
		v282 = v253 + int32(base.Ui32(v277*int32(_a_F_UpsampleBgrLinePair_SSE41_4))>>(uint(int32(8))%32))
		v284 = v282 + int32(-17685)
		if base.Ui32(v282) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_5)) {
			v291 = int32(0)
		} else {
			v291 = int32(255)
		}
		if base.Ui32(v284) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_3)) {
			v294 = int32(base.Ui32(v284) >> (uint(int32(6)) % 32))
		} else {
			v294 = v291
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v294)
		v298 = int32(8)
		v305 = v253 - (int32(base.Ui32(v256*int32(_a_F_UpsampleBgrLinePair_SSE41_7))>>(uint(v298)%32)) + int32(base.Ui32(v277*int32(_a_F_UpsampleBgrLinePair_SSE41_6))>>(uint(v298)%32)))
		v307 = v305 + int32(_a_F_UpsampleBgrLinePair_SSE41_8)
		if v305 < int32(-8708) {
			v314 = int32(0)
		} else {
			v314 = int32(255)
		}
		if base.Ui32(v307) < base.Ui32(int32(_a_F_UpsampleBgrLinePair_SSE41_3)) {
			v317 = int32(base.Ui32(v307) >> (uint(int32(6)) % 32))
		} else {
			v317 = v314
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v317)
	}
	v325 = v35 + int32(96)
	v327 = v35 + int32(64)
	if l8 < int32(34) {
		v1102 = v10
		v1103 = v162
	} else {
		v330 = int32(1)
		v332 = int32(3)
		v339 = v35 + int32(160)
		v341 = v35 + int32(128)
		v342 = int32(0)
		v354 = v342
		v355 = l7 + v332
		v359 = v342
		v360 = l6 + v332
		for {
			v376 = l4 + v354
			v377 = int32(0)
			v378 = base.Simd_g_v128_load_rng(m, v376, v377, int32(0), int32(17))
			v379 = l2 + v354
			v381 = base.Simd_g_v128_load_rng(m, v379, v377, int32(0), int32(17))
			v382 = int32(1)
			v383 = base.Simd_g_v128_load_nc(m, v376, v382)
			v384 = base.Simd_g_i8x16_avgr_u(v381, v383)
			v386 = base.Simd_g_v128_load_nc(m, v379, v382)
			v387 = base.Simd_g_i8x16_avgr_u(v378, v386)
			v389 = base.Simd_g_v128_xor(v383, v381)
			v390 = base.Simd_g_v128_xor(v378, v386)
			v392 = base.Simd_g_v128_xor(v384, v387)
			v394 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k0)
			v396 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v384, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v389, v390), v392), v394))
			v402 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v384), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v384), base.Simd_g_v128_and(v392, v389)), v394))
			v403 = base.Simd_g_i8x16_avgr_u(v378, v402)
			v409 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v396, v387), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v396, v387), base.Simd_g_v128_and(v392, v390)), v394))
			v410 = base.Simd_g_i8x16_avgr_u(v383, v409)
			v411 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k1)
			v412 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
			v413 = int32(80)
			base.Simd_g_v128_store(m, v327, v413, v412)
			v415 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
			v416 = base.Simd_g_i8x16_shuffle2(v403, v410, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v417 = int32(64)
			base.Simd_g_v128_store(m, v327, v417, v416)
			v419 = base.Simd_g_i8x16_avgr_u(v381, v409)
			v420 = base.Simd_g_i8x16_avgr_u(v386, v402)
			v422 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
			v423 = int32(16)
			base.Simd_g_v128_store(m, v327, v423, v422)
			v426 = base.Simd_g_i8x16_shuffle2(v419, v420, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			base.Simd_g_v128_store(m, v327, v377, v426)
			v429 = l5 + v354
			v431 = base.Simd_g_v128_load_rng(m, v429, v377, int32(0), int32(17))
			v432 = l3 + v354
			v434 = base.Simd_g_v128_load_rng(m, v432, v377, int32(0), int32(17))
			v436 = base.Simd_g_v128_load_nc(m, v429, v382)
			v437 = base.Simd_g_i8x16_avgr_u(v434, v436)
			v439 = base.Simd_g_v128_load_nc(m, v432, v382)
			v440 = base.Simd_g_i8x16_avgr_u(v431, v439)
			v442 = base.Simd_g_v128_xor(v436, v434)
			v443 = base.Simd_g_v128_xor(v431, v439)
			v445 = base.Simd_g_v128_xor(v437, v440)
			v448 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v437, v440), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v442, v443), v445), v394))
			v454 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v448, v437), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v448, v437), base.Simd_g_v128_and(v445, v442)), v394))
			v455 = base.Simd_g_i8x16_avgr_u(v431, v454)
			v461 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v448, v440), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v448, v440), base.Simd_g_v128_and(v445, v443)), v394))
			v462 = base.Simd_g_i8x16_avgr_u(v436, v461)
			v464 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
			base.Simd_g_v128_store(m, v327, int32(112), v464)
			v468 = base.Simd_g_i8x16_shuffle2(v455, v462, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			base.Simd_g_v128_store(m, v327, int32(96), v468)
			v471 = base.Simd_g_i8x16_avgr_u(v434, v461)
			v472 = base.Simd_g_i8x16_avgr_u(v439, v454)
			v474 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
			v475 = int32(48)
			base.Simd_g_v128_store(m, v327, v475, v474)
			v478 = base.Simd_g_i8x16_shuffle2(v471, v472, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v479 = int32(32)
			base.Simd_g_v128_store(m, v327, v479, v478)
			v481 = l0 + v330 + v359
			v482 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k7)
			v511 = base.Simd_g_v128_load64_zero(m, v481, v423)
			v513 = base.Simd_g_i8x16_shuffle2(v482, v511, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v515 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k8)
			v519 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k9)
			v520 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v513), v515), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v513), v515), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v522 = base.Simd_g_v128_load64_zero(m, v327, v423)
			v524 = base.Simd_g_i8x16_shuffle2(v482, v522, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v525 = base.Simd_g_i32x4_extend_low_i16x8_u(v524)
			v526 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k12)
			v528 = base.Simd_g_i32x4_extend_high_i16x8_u(v524)
			v533 = base.Simd_g_v128_load64_zero(m, v325, v423)
			v535 = base.Simd_g_i8x16_shuffle2(v482, v533, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v536 = base.Simd_g_i32x4_extend_low_i16x8_u(v535)
			v537 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k13)
			v539 = base.Simd_g_i32x4_extend_high_i16x8_u(v535)
			v545 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k14)
			v547 = int32(6)
			v549 = int32(24)
			v550 = base.Simd_g_v128_load64_zero(m, v481, v549)
			v552 = base.Simd_g_i8x16_shuffle2(v482, v550, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v558 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v552), v515), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v552), v515), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v560 = base.Simd_g_v128_load64_zero(m, v327, v549)
			v562 = base.Simd_g_i8x16_shuffle2(v482, v560, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v563 = base.Simd_g_i32x4_extend_low_i16x8_u(v562)
			v565 = base.Simd_g_i32x4_extend_high_i16x8_u(v562)
			v570 = base.Simd_g_v128_load64_zero(m, v325, v549)
			v572 = base.Simd_g_i8x16_shuffle2(v482, v570, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v573 = base.Simd_g_i32x4_extend_low_i16x8_u(v572)
			v575 = base.Simd_g_i32x4_extend_high_i16x8_u(v572)
			v584 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v520, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v525, v526), base.Simd_g_i32x4_mul(v528, v526), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v536, v537), base.Simd_g_i32x4_mul(v539, v537), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v545), v547), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v558, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v563, v526), base.Simd_g_i32x4_mul(v565, v526), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v573, v537), base.Simd_g_i32x4_mul(v575, v537), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v545), v547))
			v585 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k15)
			v587 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k16)
			v593 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k17)
			v605 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v525, v587), base.Simd_g_i32x4_mul(v528, v587), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v520), v593), v547), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v563, v587), base.Simd_g_i32x4_mul(v565, v587), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v558), v593), v547))
			v606 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k18)
			v609 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k19)
			v615 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k20)
			v627 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v520, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v536, v609), base.Simd_g_i32x4_mul(v539, v609), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v615), v547), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v558, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v573, v609), base.Simd_g_i32x4_mul(v575, v609), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v615), v547))
			v628 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k21)
			v630 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v584, v585), base.Simd_g_i8x16_swizzle(v605, v606)), base.Simd_g_i8x16_swizzle(v627, v628))
			base.Simd_g_v128_store(m, v360, v413, v630)
			v633 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k22)
			v635 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k23)
			v638 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k24)
			v640 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v584, v633), base.Simd_g_i8x16_swizzle(v605, v635)), base.Simd_g_i8x16_swizzle(v627, v638))
			base.Simd_g_v128_store(m, v360, v417, v640)
			v643 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k25)
			v645 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k26)
			v648 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k27)
			v650 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v584, v643), base.Simd_g_i8x16_swizzle(v605, v645)), base.Simd_g_i8x16_swizzle(v627, v648))
			base.Simd_g_v128_store(m, v360, v475, v650)
			v654 = base.Simd_g_v128_load64_zero(m, v481, v377)
			v656 = base.Simd_g_i8x16_shuffle2(v482, v654, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v662 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v656), v515), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v656), v515), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v664 = base.Simd_g_v128_load64_zero(m, v327, v377)
			v666 = base.Simd_g_i8x16_shuffle2(v482, v664, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v667 = base.Simd_g_i32x4_extend_low_i16x8_u(v666)
			v669 = base.Simd_g_i32x4_extend_high_i16x8_u(v666)
			v674 = base.Simd_g_v128_load64_zero(m, v325, v377)
			v676 = base.Simd_g_i8x16_shuffle2(v482, v674, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v677 = base.Simd_g_i32x4_extend_low_i16x8_u(v676)
			v679 = base.Simd_g_i32x4_extend_high_i16x8_u(v676)
			v688 = int32(8)
			v689 = base.Simd_g_v128_load64_zero(m, v481, v688)
			v691 = base.Simd_g_i8x16_shuffle2(v482, v689, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v697 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v691), v515), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v691), v515), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v699 = base.Simd_g_v128_load64_zero(m, v327, v688)
			v701 = base.Simd_g_i8x16_shuffle2(v482, v699, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v702 = base.Simd_g_i32x4_extend_low_i16x8_u(v701)
			v704 = base.Simd_g_i32x4_extend_high_i16x8_u(v701)
			v709 = base.Simd_g_v128_load64_zero(m, v325, v688)
			v711 = base.Simd_g_i8x16_shuffle2(v482, v709, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v712 = base.Simd_g_i32x4_extend_low_i16x8_u(v711)
			v714 = base.Simd_g_i32x4_extend_high_i16x8_u(v711)
			v723 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v662, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v667, v526), base.Simd_g_i32x4_mul(v669, v526), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v677, v537), base.Simd_g_i32x4_mul(v679, v537), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v545), v547), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v697, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v702, v526), base.Simd_g_i32x4_mul(v704, v526), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v712, v537), base.Simd_g_i32x4_mul(v714, v537), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v545), v547))
			v741 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v667, v587), base.Simd_g_i32x4_mul(v669, v587), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v662), v593), v547), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v702, v587), base.Simd_g_i32x4_mul(v704, v587), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v697), v593), v547))
			v760 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v662, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v677, v609), base.Simd_g_i32x4_mul(v679, v609), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v615), v547), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v697, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v712, v609), base.Simd_g_i32x4_mul(v714, v609), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v615), v547))
			v762 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v723, v585), base.Simd_g_i8x16_swizzle(v741, v606)), base.Simd_g_i8x16_swizzle(v760, v628))
			base.Simd_g_v128_store(m, v360, v479, v762)
			v769 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v723, v633), base.Simd_g_i8x16_swizzle(v741, v635)), base.Simd_g_i8x16_swizzle(v760, v638))
			base.Simd_g_v128_store(m, v360, v423, v769)
			v776 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v723, v643), base.Simd_g_i8x16_swizzle(v741, v645)), base.Simd_g_i8x16_swizzle(v760, v648))
			base.Simd_g_v128_store(m, v360, v377, v776)
			if l1 == int32(0) {
			} else {
				v781 = l1 + v330 + v359
				v782 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k7)
				v810 = int32(16)
				v811 = base.Simd_g_v128_load64_zero(m, v781, v810)
				v812 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
				v813 = base.Simd_g_i8x16_shuffle2(v782, v811, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v815 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k8)
				v819 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k9)
				v820 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v813), v815), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v813), v815), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
				v822 = base.Simd_g_v128_load64_zero(m, v341, v810)
				v824 = base.Simd_g_i8x16_shuffle2(v782, v822, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v825 = base.Simd_g_i32x4_extend_low_i16x8_u(v824)
				v826 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k12)
				v828 = base.Simd_g_i32x4_extend_high_i16x8_u(v824)
				v833 = base.Simd_g_v128_load64_zero(m, v339, v810)
				v835 = base.Simd_g_i8x16_shuffle2(v782, v833, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v836 = base.Simd_g_i32x4_extend_low_i16x8_u(v835)
				v837 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k13)
				v839 = base.Simd_g_i32x4_extend_high_i16x8_u(v835)
				v845 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k14)
				v847 = int32(6)
				v849 = int32(24)
				v850 = base.Simd_g_v128_load64_zero(m, v781, v849)
				v852 = base.Simd_g_i8x16_shuffle2(v782, v850, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v858 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v852), v815), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v852), v815), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
				v860 = base.Simd_g_v128_load64_zero(m, v341, v849)
				v862 = base.Simd_g_i8x16_shuffle2(v782, v860, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v863 = base.Simd_g_i32x4_extend_low_i16x8_u(v862)
				v865 = base.Simd_g_i32x4_extend_high_i16x8_u(v862)
				v870 = base.Simd_g_v128_load64_zero(m, v339, v849)
				v872 = base.Simd_g_i8x16_shuffle2(v782, v870, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v873 = base.Simd_g_i32x4_extend_low_i16x8_u(v872)
				v875 = base.Simd_g_i32x4_extend_high_i16x8_u(v872)
				v884 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v820, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v825, v826), base.Simd_g_i32x4_mul(v828, v826), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v836, v837), base.Simd_g_i32x4_mul(v839, v837), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v845), v847), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v858, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v863, v826), base.Simd_g_i32x4_mul(v865, v826), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v873, v837), base.Simd_g_i32x4_mul(v875, v837), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v845), v847))
				v885 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k15)
				v887 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k16)
				v893 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k17)
				v905 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v825, v887), base.Simd_g_i32x4_mul(v828, v887), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v820), v893), v847), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v863, v887), base.Simd_g_i32x4_mul(v865, v887), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v858), v893), v847))
				v906 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k18)
				v909 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k19)
				v915 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k20)
				v927 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v820, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v836, v909), base.Simd_g_i32x4_mul(v839, v909), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v915), v847), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v858, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v873, v909), base.Simd_g_i32x4_mul(v875, v909), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v915), v847))
				v928 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k21)
				v930 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v884, v885), base.Simd_g_i8x16_swizzle(v905, v906)), base.Simd_g_i8x16_swizzle(v927, v928))
				base.Simd_g_v128_store(m, v355, int32(80), v930)
				v933 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k22)
				v935 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k23)
				v938 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k24)
				v940 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v884, v933), base.Simd_g_i8x16_swizzle(v905, v935)), base.Simd_g_i8x16_swizzle(v927, v938))
				base.Simd_g_v128_store(m, v355, int32(64), v940)
				v943 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k25)
				v945 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k26)
				v948 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k27)
				v950 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v884, v943), base.Simd_g_i8x16_swizzle(v905, v945)), base.Simd_g_i8x16_swizzle(v927, v948))
				base.Simd_g_v128_store(m, v355, int32(48), v950)
				v953 = int32(0)
				v954 = base.Simd_g_v128_load64_zero(m, v781, v953)
				v956 = base.Simd_g_i8x16_shuffle2(v782, v954, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v962 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v956), v815), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v956), v815), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
				v964 = base.Simd_g_v128_load64_zero(m, v341, v953)
				v966 = base.Simd_g_i8x16_shuffle2(v782, v964, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v967 = base.Simd_g_i32x4_extend_low_i16x8_u(v966)
				v969 = base.Simd_g_i32x4_extend_high_i16x8_u(v966)
				v974 = base.Simd_g_v128_load64_zero(m, v339, v953)
				v976 = base.Simd_g_i8x16_shuffle2(v782, v974, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v977 = base.Simd_g_i32x4_extend_low_i16x8_u(v976)
				v979 = base.Simd_g_i32x4_extend_high_i16x8_u(v976)
				v988 = int32(8)
				v989 = base.Simd_g_v128_load64_zero(m, v781, v988)
				v991 = base.Simd_g_i8x16_shuffle2(v782, v989, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v997 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v991), v815), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v991), v815), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
				v999 = base.Simd_g_v128_load64_zero(m, v341, v988)
				v1001 = base.Simd_g_i8x16_shuffle2(v782, v999, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v1002 = base.Simd_g_i32x4_extend_low_i16x8_u(v1001)
				v1004 = base.Simd_g_i32x4_extend_high_i16x8_u(v1001)
				v1009 = base.Simd_g_v128_load64_zero(m, v339, v988)
				v1011 = base.Simd_g_i8x16_shuffle2(v782, v1009, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
				v1012 = base.Simd_g_i32x4_extend_low_i16x8_u(v1011)
				v1014 = base.Simd_g_i32x4_extend_high_i16x8_u(v1011)
				v1023 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v962, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v967, v826), base.Simd_g_i32x4_mul(v969, v826), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v977, v837), base.Simd_g_i32x4_mul(v979, v837), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v845), v847), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v997, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1002, v826), base.Simd_g_i32x4_mul(v1004, v826), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1012, v837), base.Simd_g_i32x4_mul(v1014, v837), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v845), v847))
				v1041 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v967, v887), base.Simd_g_i32x4_mul(v969, v887), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v962), v893), v847), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1002, v887), base.Simd_g_i32x4_mul(v1004, v887), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v997), v893), v847))
				v1060 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v962, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v977, v909), base.Simd_g_i32x4_mul(v979, v909), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v915), v847), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v997, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1012, v909), base.Simd_g_i32x4_mul(v1014, v909), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v915), v847))
				v1062 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1023, v885), base.Simd_g_i8x16_swizzle(v1041, v906)), base.Simd_g_i8x16_swizzle(v1060, v928))
				base.Simd_g_v128_store(m, v355, int32(32), v1062)
				v1069 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1023, v933), base.Simd_g_i8x16_swizzle(v1041, v935)), base.Simd_g_i8x16_swizzle(v1060, v938))
				base.Simd_g_v128_store(m, v355, v810, v1069)
				v1076 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1023, v943), base.Simd_g_i8x16_swizzle(v1041, v945)), base.Simd_g_i8x16_swizzle(v1060, v948))
				base.Simd_g_v128_store(m, v355, v953, v1076)
			}
			v1079 = int32(96)
			v1084 = v354 + int32(16)
			if v359+int32(66) <= l8 {
				v354 = v1084
				v355 = v355 + v1079
				v359 = v359 + int32(32)
				v360 = v360 + v1079
				continue
			} else {
				break
			}
			break
		}
		v1102 = v1084
		v1103 = v359 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v1126 = int32(32)
		v1129 = int32(1)
		v1135 = int32(base.Ui32(l8+v1129)>>(uint(v1129)%32)) - int32(base.Ui32(v1103)>>(uint(v1129)%32))
		v1136 = F_memcpy(m, v35+v1126, l2+v1102, v1135)
		mBase = m.M
		v1138 = F_memcpy(m, v35, l4+v1102, v1135)
		mBase = m.M
		v1140 = v1138 + v1126
		v1141 = v1140 + v1135
		v1145 = v1135 + int32(-1)
		v1146 = v1140 + v1145
		v1147 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1146))))
		v1149 = int32(17) - v1135
		if base.Ui32(v1149) < base.Ui32(int32(33)) {
			if v1149 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1141))) = uint8(v1147)
				v1160 = v1141 + v1149
				*(*uint8)(unsafe.Add(mBase, uint32(v1160+int32(-1)))) = uint8(v1147)
				if base.Ui32(v1149) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1141)+2)) = uint8(v1147)
					*(*uint8)(unsafe.Add(mBase, uint32(v1141)+1)) = uint8(v1147)
					*(*uint8)(unsafe.Add(mBase, uint32(v1160+int32(-3)))) = uint8(v1147)
					*(*uint8)(unsafe.Add(mBase, uint32(v1160+int32(-2)))) = uint8(v1147)
					if base.Ui32(v1149) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1141)+3)) = uint8(v1147)
						*(*uint8)(unsafe.Add(mBase, uint32(v1160+int32(-4)))) = uint8(v1147)
						if base.Ui32(v1149) < base.Ui32(int32(9)) {
						} else {
							v1185 = (int32(0) - v1141) & int32(3)
							v1186 = v1141 + v1185
							v1190 = v1147 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1186))) = v1190
							v1194 = (v1149 - v1185) & int32(60)
							v1195 = v1186 + v1194
							*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-4)))) = v1190
							if base.Ui32(v1194) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1186)+8)) = v1190
								*(*int32)(unsafe.Add(mBase, uint32(v1186)+4)) = v1190
								*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-8)))) = v1190
								*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-12)))) = v1190
								if base.Ui32(v1194) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1186)+24)) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1186)+20)) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1186)+16)) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1186)+12)) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-16)))) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-20)))) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-24)))) = v1190
									*(*int32)(unsafe.Add(mBase, uint32(v1195+int32(-28)))) = v1190
									v1230 = v1186&int32(4) | int32(24)
									v1231 = v1194 - v1230
									if base.Ui32(v1231) < base.Ui32(int32(32)) {
									} else {
										v1236 = base.I64_extend_i32_u(v1190) * int64(4294967297)
										v1239 = v1231
										v1240 = v1186 + v1230
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1240)+24)) = v1236
											*(*int64)(unsafe.Add(mBase, uint32(v1240)+16)) = v1236
											*(*int64)(unsafe.Add(mBase, uint32(v1240)+8)) = v1236
											*(*int64)(unsafe.Add(mBase, uint32(v1240))) = v1236
											v1252 = v1239 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1252) {
												v1239 = v1252
												v1240 = v1240 + int32(32)
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
			base.MemoryFill(m, v1141, v1147, v1149)
		}
		v1270 = v1138 + v1135
		v1271 = v1138 + v1145
		v1272 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1271))))
		if base.Ui32(v1149) < base.Ui32(int32(33)) {
			if v1149 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1270))) = uint8(v1272)
				v1283 = v1270 + v1149
				*(*uint8)(unsafe.Add(mBase, uint32(v1283+int32(-1)))) = uint8(v1272)
				if base.Ui32(v1149) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1270)+2)) = uint8(v1272)
					*(*uint8)(unsafe.Add(mBase, uint32(v1270)+1)) = uint8(v1272)
					*(*uint8)(unsafe.Add(mBase, uint32(v1283+int32(-3)))) = uint8(v1272)
					*(*uint8)(unsafe.Add(mBase, uint32(v1283+int32(-2)))) = uint8(v1272)
					if base.Ui32(v1149) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1270)+3)) = uint8(v1272)
						*(*uint8)(unsafe.Add(mBase, uint32(v1283+int32(-4)))) = uint8(v1272)
						if base.Ui32(v1149) < base.Ui32(int32(9)) {
						} else {
							v1308 = (int32(0) - v1270) & int32(3)
							v1309 = v1270 + v1308
							v1313 = v1272 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1309))) = v1313
							v1317 = (v1149 - v1308) & int32(60)
							v1318 = v1309 + v1317
							*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-4)))) = v1313
							if base.Ui32(v1317) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1309)+8)) = v1313
								*(*int32)(unsafe.Add(mBase, uint32(v1309)+4)) = v1313
								*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-8)))) = v1313
								*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-12)))) = v1313
								if base.Ui32(v1317) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1309)+24)) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1309)+20)) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1309)+16)) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1309)+12)) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-16)))) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-20)))) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-24)))) = v1313
									*(*int32)(unsafe.Add(mBase, uint32(v1318+int32(-28)))) = v1313
									v1353 = v1309&int32(4) | int32(24)
									v1354 = v1317 - v1353
									if base.Ui32(v1354) < base.Ui32(int32(32)) {
									} else {
										v1359 = base.I64_extend_i32_u(v1313) * int64(4294967297)
										v1362 = v1354
										v1363 = v1309 + v1353
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1363)+24)) = v1359
											*(*int64)(unsafe.Add(mBase, uint32(v1363)+16)) = v1359
											*(*int64)(unsafe.Add(mBase, uint32(v1363)+8)) = v1359
											*(*int64)(unsafe.Add(mBase, uint32(v1363))) = v1359
											v1375 = v1362 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1375) {
												v1362 = v1375
												v1363 = v1363 + int32(32)
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
			base.MemoryFill(m, v1270, v1272, v1149)
		}
		v1393 = int32(0)
		v1394 = base.Simd_g_v128_load_rng(m, v1138, v1393, int32(0), int32(49))
		v1395 = int32(32)
		v1396 = base.Simd_g_v128_load_nc(m, v1138, v1395)
		v1398 = base.Simd_g_v128_load_nc(m, v1138, int32(1))
		v1399 = base.Simd_g_i8x16_avgr_u(v1396, v1398)
		v1400 = int32(33)
		v1401 = base.Simd_g_v128_load_nc(m, v1138, v1400)
		v1402 = base.Simd_g_i8x16_avgr_u(v1394, v1401)
		v1404 = base.Simd_g_v128_xor(v1398, v1396)
		v1405 = base.Simd_g_v128_xor(v1394, v1401)
		v1407 = base.Simd_g_v128_xor(v1399, v1402)
		v1409 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k0)
		v1411 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1399, v1402), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1404, v1405), v1407), v1409))
		v1417 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1411, v1399), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1411, v1399), base.Simd_g_v128_and(v1407, v1404)), v1409))
		v1418 = base.Simd_g_i8x16_avgr_u(v1394, v1417)
		v1424 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1411, v1402), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1411, v1402), base.Simd_g_v128_and(v1407, v1405)), v1409))
		v1425 = base.Simd_g_i8x16_avgr_u(v1398, v1424)
		v1426 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k1)
		v1427 = base.Simd_g_i8x16_shuffle2(v1418, v1425, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(80), v1427)
		v1430 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
		v1431 = base.Simd_g_i8x16_shuffle2(v1418, v1425, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, int32(64), v1431)
		v1434 = base.Simd_g_i8x16_avgr_u(v1396, v1424)
		v1435 = base.Simd_g_i8x16_avgr_u(v1401, v1417)
		v1437 = base.Simd_g_i8x16_shuffle2(v1434, v1435, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(16), v1437)
		v1441 = base.Simd_g_i8x16_shuffle2(v1434, v1435, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, v1393, v1441)
		v1447 = F_memcpy(m, v1138+v1395, l3+v1102, v1135)
		mBase = m.M
		v1449 = F_memcpy(m, v1138, l5+v1102, v1135)
		mBase = m.M
		v1450 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1146))))
		if base.Ui32(v1149) < base.Ui32(v1400) {
			if v1149 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1141))) = uint8(v1450)
				v1461 = v1141 + v1149
				*(*uint8)(unsafe.Add(mBase, uint32(v1461+int32(-1)))) = uint8(v1450)
				if base.Ui32(v1149) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1141)+2)) = uint8(v1450)
					*(*uint8)(unsafe.Add(mBase, uint32(v1141)+1)) = uint8(v1450)
					*(*uint8)(unsafe.Add(mBase, uint32(v1461+int32(-3)))) = uint8(v1450)
					*(*uint8)(unsafe.Add(mBase, uint32(v1461+int32(-2)))) = uint8(v1450)
					if base.Ui32(v1149) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1141)+3)) = uint8(v1450)
						*(*uint8)(unsafe.Add(mBase, uint32(v1461+int32(-4)))) = uint8(v1450)
						if base.Ui32(v1149) < base.Ui32(int32(9)) {
						} else {
							v1486 = (int32(0) - v1141) & int32(3)
							v1487 = v1141 + v1486
							v1491 = v1450 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1487))) = v1491
							v1495 = (v1149 - v1486) & int32(60)
							v1496 = v1487 + v1495
							*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-4)))) = v1491
							if base.Ui32(v1495) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1487)+8)) = v1491
								*(*int32)(unsafe.Add(mBase, uint32(v1487)+4)) = v1491
								*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-8)))) = v1491
								*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-12)))) = v1491
								if base.Ui32(v1495) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1487)+24)) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1487)+20)) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1487)+16)) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1487)+12)) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-16)))) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-20)))) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-24)))) = v1491
									*(*int32)(unsafe.Add(mBase, uint32(v1496+int32(-28)))) = v1491
									v1531 = v1487&int32(4) | int32(24)
									v1532 = v1495 - v1531
									if base.Ui32(v1532) < base.Ui32(int32(32)) {
									} else {
										v1537 = base.I64_extend_i32_u(v1491) * int64(4294967297)
										v1540 = v1532
										v1541 = v1487 + v1531
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1541)+24)) = v1537
											*(*int64)(unsafe.Add(mBase, uint32(v1541)+16)) = v1537
											*(*int64)(unsafe.Add(mBase, uint32(v1541)+8)) = v1537
											*(*int64)(unsafe.Add(mBase, uint32(v1541))) = v1537
											v1553 = v1540 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1553) {
												v1540 = v1553
												v1541 = v1541 + int32(32)
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
			base.MemoryFill(m, v1141, v1450, v1149)
		}
		v1571 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1271))))
		if base.Ui32(v1149) < base.Ui32(int32(33)) {
			if v1149 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v1270))) = uint8(v1571)
				v1582 = v1270 + v1149
				*(*uint8)(unsafe.Add(mBase, uint32(v1582+int32(-1)))) = uint8(v1571)
				if base.Ui32(v1149) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v1270)+2)) = uint8(v1571)
					*(*uint8)(unsafe.Add(mBase, uint32(v1270)+1)) = uint8(v1571)
					*(*uint8)(unsafe.Add(mBase, uint32(v1582+int32(-3)))) = uint8(v1571)
					*(*uint8)(unsafe.Add(mBase, uint32(v1582+int32(-2)))) = uint8(v1571)
					if base.Ui32(v1149) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v1270)+3)) = uint8(v1571)
						*(*uint8)(unsafe.Add(mBase, uint32(v1582+int32(-4)))) = uint8(v1571)
						if base.Ui32(v1149) < base.Ui32(int32(9)) {
						} else {
							v1607 = (int32(0) - v1270) & int32(3)
							v1608 = v1270 + v1607
							v1612 = v1571 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1608))) = v1612
							v1616 = (v1149 - v1607) & int32(60)
							v1617 = v1608 + v1616
							*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-4)))) = v1612
							if base.Ui32(v1616) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1608)+8)) = v1612
								*(*int32)(unsafe.Add(mBase, uint32(v1608)+4)) = v1612
								*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-8)))) = v1612
								*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-12)))) = v1612
								if base.Ui32(v1616) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1608)+24)) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1608)+20)) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1608)+16)) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1608)+12)) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-16)))) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-20)))) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-24)))) = v1612
									*(*int32)(unsafe.Add(mBase, uint32(v1617+int32(-28)))) = v1612
									v1652 = v1608&int32(4) | int32(24)
									v1653 = v1616 - v1652
									if base.Ui32(v1653) < base.Ui32(int32(32)) {
									} else {
										v1658 = base.I64_extend_i32_u(v1612) * int64(4294967297)
										v1661 = v1653
										v1662 = v1608 + v1652
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1662)+24)) = v1658
											*(*int64)(unsafe.Add(mBase, uint32(v1662)+16)) = v1658
											*(*int64)(unsafe.Add(mBase, uint32(v1662)+8)) = v1658
											*(*int64)(unsafe.Add(mBase, uint32(v1662))) = v1658
											v1674 = v1661 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1674) {
												v1661 = v1674
												v1662 = v1662 + int32(32)
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
			base.MemoryFill(m, v1270, v1571, v1149)
		}
		v1693 = base.Simd_g_v128_load_rng(m, v1449, int32(0), int32(0), int32(49))
		v1694 = int32(32)
		v1695 = base.Simd_g_v128_load_nc(m, v1449, v1694)
		v1697 = base.Simd_g_v128_load_nc(m, v1449, int32(1))
		v1698 = base.Simd_g_i8x16_avgr_u(v1695, v1697)
		v1700 = base.Simd_g_v128_load_nc(m, v1449, int32(33))
		v1701 = base.Simd_g_i8x16_avgr_u(v1693, v1700)
		v1703 = base.Simd_g_v128_xor(v1697, v1695)
		v1704 = base.Simd_g_v128_xor(v1693, v1700)
		v1706 = base.Simd_g_v128_xor(v1698, v1701)
		v1709 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1698, v1701), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1703, v1704), v1706), v1409))
		v1715 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1709, v1698), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1709, v1698), base.Simd_g_v128_and(v1706, v1703)), v1409))
		v1716 = base.Simd_g_i8x16_avgr_u(v1693, v1715)
		v1722 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1709, v1701), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1709, v1701), base.Simd_g_v128_and(v1706, v1704)), v1409))
		v1723 = base.Simd_g_i8x16_avgr_u(v1697, v1722)
		v1724 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k1)
		v1725 = base.Simd_g_i8x16_shuffle2(v1716, v1723, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(112), v1725)
		v1728 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
		v1729 = base.Simd_g_i8x16_shuffle2(v1716, v1723, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, int32(96), v1729)
		v1732 = base.Simd_g_i8x16_avgr_u(v1695, v1722)
		v1733 = base.Simd_g_i8x16_avgr_u(v1700, v1715)
		v1735 = base.Simd_g_i8x16_shuffle2(v1732, v1733, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k2), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k3))
		base.Simd_g_v128_store(m, v327, int32(48), v1735)
		v1739 = base.Simd_g_i8x16_shuffle2(v1732, v1733, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
		base.Simd_g_v128_store(m, v327, v1694, v1739)
		v1745 = l8 - v1103
		v1746 = F_memcpy(m, v35+int32(448), l0+v1103, v1745)
		mBase = m.M
		v1748 = v35 + int32(192)
		if l1 != 0 {
			v2055 = F_memcpy(m, v35+int32(480), l1+v1103, v1745)
			mBase = m.M
			v2056 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k7)
			v2084 = int32(16)
			v2085 = base.Simd_g_v128_load64_zero(m, v1746, v2084)
			v2086 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
			v2087 = base.Simd_g_i8x16_shuffle2(v2056, v2085, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2089 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k8)
			v2093 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k9)
			v2094 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2087), v2089), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2087), v2089), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2096 = base.Simd_g_v128_load64_zero(m, v327, v2084)
			v2098 = base.Simd_g_i8x16_shuffle2(v2056, v2096, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2099 = base.Simd_g_i32x4_extend_low_i16x8_u(v2098)
			v2100 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k12)
			v2102 = base.Simd_g_i32x4_extend_high_i16x8_u(v2098)
			v2107 = base.Simd_g_v128_load64_zero(m, v325, v2084)
			v2109 = base.Simd_g_i8x16_shuffle2(v2056, v2107, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2110 = base.Simd_g_i32x4_extend_low_i16x8_u(v2109)
			v2111 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k13)
			v2113 = base.Simd_g_i32x4_extend_high_i16x8_u(v2109)
			v2119 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k14)
			v2121 = int32(6)
			v2123 = int32(24)
			v2124 = base.Simd_g_v128_load64_zero(m, v1746, v2123)
			v2126 = base.Simd_g_i8x16_shuffle2(v2056, v2124, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2132 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2126), v2089), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2126), v2089), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2134 = base.Simd_g_v128_load64_zero(m, v327, v2123)
			v2136 = base.Simd_g_i8x16_shuffle2(v2056, v2134, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2137 = base.Simd_g_i32x4_extend_low_i16x8_u(v2136)
			v2139 = base.Simd_g_i32x4_extend_high_i16x8_u(v2136)
			v2144 = base.Simd_g_v128_load64_zero(m, v325, v2123)
			v2146 = base.Simd_g_i8x16_shuffle2(v2056, v2144, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2147 = base.Simd_g_i32x4_extend_low_i16x8_u(v2146)
			v2149 = base.Simd_g_i32x4_extend_high_i16x8_u(v2146)
			v2158 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2094, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2099, v2100), base.Simd_g_i32x4_mul(v2102, v2100), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2110, v2111), base.Simd_g_i32x4_mul(v2113, v2111), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2119), v2121), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2132, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2137, v2100), base.Simd_g_i32x4_mul(v2139, v2100), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2147, v2111), base.Simd_g_i32x4_mul(v2149, v2111), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2119), v2121))
			v2159 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k15)
			v2161 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k16)
			v2167 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k17)
			v2179 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2099, v2161), base.Simd_g_i32x4_mul(v2102, v2161), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2094), v2167), v2121), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2137, v2161), base.Simd_g_i32x4_mul(v2139, v2161), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2132), v2167), v2121))
			v2180 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k18)
			v2183 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k19)
			v2189 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k20)
			v2201 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2094, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2110, v2183), base.Simd_g_i32x4_mul(v2113, v2183), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2189), v2121), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2132, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2147, v2183), base.Simd_g_i32x4_mul(v2149, v2183), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2189), v2121))
			v2202 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k21)
			v2204 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2158, v2159), base.Simd_g_i8x16_swizzle(v2179, v2180)), base.Simd_g_i8x16_swizzle(v2201, v2202))
			base.Simd_g_v128_store(m, v1748, int32(80), v2204)
			v2207 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k22)
			v2209 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k23)
			v2212 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k24)
			v2214 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2158, v2207), base.Simd_g_i8x16_swizzle(v2179, v2209)), base.Simd_g_i8x16_swizzle(v2201, v2212))
			base.Simd_g_v128_store(m, v1748, int32(64), v2214)
			v2217 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k25)
			v2219 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k26)
			v2222 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k27)
			v2224 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2158, v2217), base.Simd_g_i8x16_swizzle(v2179, v2219)), base.Simd_g_i8x16_swizzle(v2201, v2222))
			base.Simd_g_v128_store(m, v1748, int32(48), v2224)
			v2227 = int32(0)
			v2228 = base.Simd_g_v128_load64_zero(m, v1746, v2227)
			v2230 = base.Simd_g_i8x16_shuffle2(v2056, v2228, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2236 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2230), v2089), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2230), v2089), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2238 = base.Simd_g_v128_load64_zero(m, v327, v2227)
			v2240 = base.Simd_g_i8x16_shuffle2(v2056, v2238, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2241 = base.Simd_g_i32x4_extend_low_i16x8_u(v2240)
			v2243 = base.Simd_g_i32x4_extend_high_i16x8_u(v2240)
			v2248 = base.Simd_g_v128_load64_zero(m, v325, v2227)
			v2250 = base.Simd_g_i8x16_shuffle2(v2056, v2248, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2251 = base.Simd_g_i32x4_extend_low_i16x8_u(v2250)
			v2253 = base.Simd_g_i32x4_extend_high_i16x8_u(v2250)
			v2262 = int32(8)
			v2263 = base.Simd_g_v128_load64_zero(m, v1746, v2262)
			v2265 = base.Simd_g_i8x16_shuffle2(v2056, v2263, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2271 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2265), v2089), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2265), v2089), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2273 = base.Simd_g_v128_load64_zero(m, v327, v2262)
			v2275 = base.Simd_g_i8x16_shuffle2(v2056, v2273, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2276 = base.Simd_g_i32x4_extend_low_i16x8_u(v2275)
			v2278 = base.Simd_g_i32x4_extend_high_i16x8_u(v2275)
			v2283 = base.Simd_g_v128_load64_zero(m, v325, v2262)
			v2285 = base.Simd_g_i8x16_shuffle2(v2056, v2283, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2286 = base.Simd_g_i32x4_extend_low_i16x8_u(v2285)
			v2288 = base.Simd_g_i32x4_extend_high_i16x8_u(v2285)
			v2297 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2236, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2241, v2100), base.Simd_g_i32x4_mul(v2243, v2100), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2251, v2111), base.Simd_g_i32x4_mul(v2253, v2111), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2119), v2121), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2271, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2276, v2100), base.Simd_g_i32x4_mul(v2278, v2100), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2286, v2111), base.Simd_g_i32x4_mul(v2288, v2111), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2119), v2121))
			v2315 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2241, v2161), base.Simd_g_i32x4_mul(v2243, v2161), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2236), v2167), v2121), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2276, v2161), base.Simd_g_i32x4_mul(v2278, v2161), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2271), v2167), v2121))
			v2334 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2236, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2251, v2183), base.Simd_g_i32x4_mul(v2253, v2183), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2189), v2121), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2271, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2286, v2183), base.Simd_g_i32x4_mul(v2288, v2183), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2189), v2121))
			v2336 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2297, v2159), base.Simd_g_i8x16_swizzle(v2315, v2180)), base.Simd_g_i8x16_swizzle(v2334, v2202))
			base.Simd_g_v128_store(m, v1748, int32(32), v2336)
			v2343 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2297, v2207), base.Simd_g_i8x16_swizzle(v2315, v2209)), base.Simd_g_i8x16_swizzle(v2334, v2212))
			base.Simd_g_v128_store(m, v1748, v2084, v2343)
			v2350 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2297, v2217), base.Simd_g_i8x16_swizzle(v2315, v2219)), base.Simd_g_i8x16_swizzle(v2334, v2222))
			base.Simd_g_v128_store(m, v1748, v2227, v2350)
			v2354 = v35 + int32(128)
			v2356 = v35 + int32(160)
			v2358 = v35 + int32(320)
			v2359 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k7)
			v2387 = int32(16)
			v2388 = base.Simd_g_v128_load64_zero(m, v2055, v2387)
			v2389 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
			v2390 = base.Simd_g_i8x16_shuffle2(v2359, v2388, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2392 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k8)
			v2396 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k9)
			v2397 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2390), v2392), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2390), v2392), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2399 = base.Simd_g_v128_load64_zero(m, v2354, v2387)
			v2401 = base.Simd_g_i8x16_shuffle2(v2359, v2399, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2402 = base.Simd_g_i32x4_extend_low_i16x8_u(v2401)
			v2403 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k12)
			v2405 = base.Simd_g_i32x4_extend_high_i16x8_u(v2401)
			v2410 = base.Simd_g_v128_load64_zero(m, v2356, v2387)
			v2412 = base.Simd_g_i8x16_shuffle2(v2359, v2410, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2413 = base.Simd_g_i32x4_extend_low_i16x8_u(v2412)
			v2414 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k13)
			v2416 = base.Simd_g_i32x4_extend_high_i16x8_u(v2412)
			v2422 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k14)
			v2424 = int32(6)
			v2426 = int32(24)
			v2427 = base.Simd_g_v128_load64_zero(m, v2055, v2426)
			v2429 = base.Simd_g_i8x16_shuffle2(v2359, v2427, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2435 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2429), v2392), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2429), v2392), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2437 = base.Simd_g_v128_load64_zero(m, v2354, v2426)
			v2439 = base.Simd_g_i8x16_shuffle2(v2359, v2437, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2440 = base.Simd_g_i32x4_extend_low_i16x8_u(v2439)
			v2442 = base.Simd_g_i32x4_extend_high_i16x8_u(v2439)
			v2447 = base.Simd_g_v128_load64_zero(m, v2356, v2426)
			v2449 = base.Simd_g_i8x16_shuffle2(v2359, v2447, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2450 = base.Simd_g_i32x4_extend_low_i16x8_u(v2449)
			v2452 = base.Simd_g_i32x4_extend_high_i16x8_u(v2449)
			v2461 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2397, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2402, v2403), base.Simd_g_i32x4_mul(v2405, v2403), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2413, v2414), base.Simd_g_i32x4_mul(v2416, v2414), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2422), v2424), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2435, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2440, v2403), base.Simd_g_i32x4_mul(v2442, v2403), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2450, v2414), base.Simd_g_i32x4_mul(v2452, v2414), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2422), v2424))
			v2462 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k15)
			v2464 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k16)
			v2470 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k17)
			v2482 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2402, v2464), base.Simd_g_i32x4_mul(v2405, v2464), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2397), v2470), v2424), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2440, v2464), base.Simd_g_i32x4_mul(v2442, v2464), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2435), v2470), v2424))
			v2483 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k18)
			v2486 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k19)
			v2492 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k20)
			v2504 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2397, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2413, v2486), base.Simd_g_i32x4_mul(v2416, v2486), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2492), v2424), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2435, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2450, v2486), base.Simd_g_i32x4_mul(v2452, v2486), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2492), v2424))
			v2505 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k21)
			v2507 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2461, v2462), base.Simd_g_i8x16_swizzle(v2482, v2483)), base.Simd_g_i8x16_swizzle(v2504, v2505))
			base.Simd_g_v128_store(m, v2358, int32(80), v2507)
			v2510 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k22)
			v2512 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k23)
			v2515 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k24)
			v2517 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2461, v2510), base.Simd_g_i8x16_swizzle(v2482, v2512)), base.Simd_g_i8x16_swizzle(v2504, v2515))
			base.Simd_g_v128_store(m, v2358, int32(64), v2517)
			v2520 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k25)
			v2522 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k26)
			v2525 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k27)
			v2527 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2461, v2520), base.Simd_g_i8x16_swizzle(v2482, v2522)), base.Simd_g_i8x16_swizzle(v2504, v2525))
			base.Simd_g_v128_store(m, v2358, int32(48), v2527)
			v2530 = int32(0)
			v2531 = base.Simd_g_v128_load64_zero(m, v2055, v2530)
			v2533 = base.Simd_g_i8x16_shuffle2(v2359, v2531, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2539 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2533), v2392), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2533), v2392), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2541 = base.Simd_g_v128_load64_zero(m, v2354, v2530)
			v2543 = base.Simd_g_i8x16_shuffle2(v2359, v2541, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2544 = base.Simd_g_i32x4_extend_low_i16x8_u(v2543)
			v2546 = base.Simd_g_i32x4_extend_high_i16x8_u(v2543)
			v2551 = base.Simd_g_v128_load64_zero(m, v2356, v2530)
			v2553 = base.Simd_g_i8x16_shuffle2(v2359, v2551, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2554 = base.Simd_g_i32x4_extend_low_i16x8_u(v2553)
			v2556 = base.Simd_g_i32x4_extend_high_i16x8_u(v2553)
			v2565 = int32(8)
			v2566 = base.Simd_g_v128_load64_zero(m, v2055, v2565)
			v2568 = base.Simd_g_i8x16_shuffle2(v2359, v2566, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2574 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2568), v2392), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v2568), v2392), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v2576 = base.Simd_g_v128_load64_zero(m, v2354, v2565)
			v2578 = base.Simd_g_i8x16_shuffle2(v2359, v2576, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2579 = base.Simd_g_i32x4_extend_low_i16x8_u(v2578)
			v2581 = base.Simd_g_i32x4_extend_high_i16x8_u(v2578)
			v2586 = base.Simd_g_v128_load64_zero(m, v2356, v2565)
			v2588 = base.Simd_g_i8x16_shuffle2(v2359, v2586, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v2589 = base.Simd_g_i32x4_extend_low_i16x8_u(v2588)
			v2591 = base.Simd_g_i32x4_extend_high_i16x8_u(v2588)
			v2600 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2539, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2544, v2403), base.Simd_g_i32x4_mul(v2546, v2403), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2554, v2414), base.Simd_g_i32x4_mul(v2556, v2414), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2422), v2424), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v2574, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2579, v2403), base.Simd_g_i32x4_mul(v2581, v2403), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2589, v2414), base.Simd_g_i32x4_mul(v2591, v2414), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v2422), v2424))
			v2618 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2544, v2464), base.Simd_g_i32x4_mul(v2546, v2464), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2539), v2470), v2424), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2579, v2464), base.Simd_g_i32x4_mul(v2581, v2464), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v2574), v2470), v2424))
			v2637 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2539, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2554, v2486), base.Simd_g_i32x4_mul(v2556, v2486), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2492), v2424), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v2574, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v2589, v2486), base.Simd_g_i32x4_mul(v2591, v2486), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v2492), v2424))
			v2639 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2600, v2462), base.Simd_g_i8x16_swizzle(v2618, v2483)), base.Simd_g_i8x16_swizzle(v2637, v2505))
			base.Simd_g_v128_store(m, v2358, int32(32), v2639)
			v2646 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2600, v2510), base.Simd_g_i8x16_swizzle(v2618, v2512)), base.Simd_g_i8x16_swizzle(v2637, v2515))
			base.Simd_g_v128_store(m, v2358, v2387, v2646)
			v2653 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v2600, v2520), base.Simd_g_i8x16_swizzle(v2618, v2522)), base.Simd_g_i8x16_swizzle(v2637, v2525))
			base.Simd_g_v128_store(m, v2358, v2530, v2653)
			v2656 = int32(3)
			v2657 = v1103 * v2656
			v2660 = v1745 * v2656
			v2661 = F_memcpy(m, l6+v2657, v1748, v2660)
			mBase = m.M
			v2663 = F_memcpy(m, l7+v2657, v2358, v2660)
			mBase = m.M
		} else {
			v1749 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k7)
			v1777 = int32(16)
			v1778 = base.Simd_g_v128_load64_zero(m, v1746, v1777)
			v1779 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k4)
			v1780 = base.Simd_g_i8x16_shuffle2(v1749, v1778, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1782 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k8)
			v1786 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k9)
			v1787 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1780), v1782), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1780), v1782), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v1789 = base.Simd_g_v128_load64_zero(m, v327, v1777)
			v1791 = base.Simd_g_i8x16_shuffle2(v1749, v1789, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1792 = base.Simd_g_i32x4_extend_low_i16x8_u(v1791)
			v1793 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k12)
			v1795 = base.Simd_g_i32x4_extend_high_i16x8_u(v1791)
			v1800 = base.Simd_g_v128_load64_zero(m, v325, v1777)
			v1802 = base.Simd_g_i8x16_shuffle2(v1749, v1800, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1803 = base.Simd_g_i32x4_extend_low_i16x8_u(v1802)
			v1804 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k13)
			v1806 = base.Simd_g_i32x4_extend_high_i16x8_u(v1802)
			v1812 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k14)
			v1814 = int32(6)
			v1816 = int32(24)
			v1817 = base.Simd_g_v128_load64_zero(m, v1746, v1816)
			v1819 = base.Simd_g_i8x16_shuffle2(v1749, v1817, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1825 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1819), v1782), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1819), v1782), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v1827 = base.Simd_g_v128_load64_zero(m, v327, v1816)
			v1829 = base.Simd_g_i8x16_shuffle2(v1749, v1827, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1830 = base.Simd_g_i32x4_extend_low_i16x8_u(v1829)
			v1832 = base.Simd_g_i32x4_extend_high_i16x8_u(v1829)
			v1837 = base.Simd_g_v128_load64_zero(m, v325, v1816)
			v1839 = base.Simd_g_i8x16_shuffle2(v1749, v1837, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1840 = base.Simd_g_i32x4_extend_low_i16x8_u(v1839)
			v1842 = base.Simd_g_i32x4_extend_high_i16x8_u(v1839)
			v1851 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1787, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1792, v1793), base.Simd_g_i32x4_mul(v1795, v1793), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1803, v1804), base.Simd_g_i32x4_mul(v1806, v1804), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v1812), v1814), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1825, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1830, v1793), base.Simd_g_i32x4_mul(v1832, v1793), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1840, v1804), base.Simd_g_i32x4_mul(v1842, v1804), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v1812), v1814))
			v1852 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k15)
			v1854 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k16)
			v1860 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k17)
			v1872 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1792, v1854), base.Simd_g_i32x4_mul(v1795, v1854), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v1787), v1860), v1814), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1830, v1854), base.Simd_g_i32x4_mul(v1832, v1854), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v1825), v1860), v1814))
			v1873 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k18)
			v1876 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k19)
			v1882 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k20)
			v1894 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1787, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1803, v1876), base.Simd_g_i32x4_mul(v1806, v1876), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v1882), v1814), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1825, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1840, v1876), base.Simd_g_i32x4_mul(v1842, v1876), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v1882), v1814))
			v1895 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k21)
			v1897 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1851, v1852), base.Simd_g_i8x16_swizzle(v1872, v1873)), base.Simd_g_i8x16_swizzle(v1894, v1895))
			base.Simd_g_v128_store(m, v1748, int32(80), v1897)
			v1900 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k22)
			v1902 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k23)
			v1905 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k24)
			v1907 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1851, v1900), base.Simd_g_i8x16_swizzle(v1872, v1902)), base.Simd_g_i8x16_swizzle(v1894, v1905))
			base.Simd_g_v128_store(m, v1748, int32(64), v1907)
			v1910 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k25)
			v1912 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k26)
			v1915 = base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k27)
			v1917 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1851, v1910), base.Simd_g_i8x16_swizzle(v1872, v1912)), base.Simd_g_i8x16_swizzle(v1894, v1915))
			base.Simd_g_v128_store(m, v1748, int32(48), v1917)
			v1920 = int32(0)
			v1921 = base.Simd_g_v128_load64_zero(m, v1746, v1920)
			v1923 = base.Simd_g_i8x16_shuffle2(v1749, v1921, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1929 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1923), v1782), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1923), v1782), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v1931 = base.Simd_g_v128_load64_zero(m, v327, v1920)
			v1933 = base.Simd_g_i8x16_shuffle2(v1749, v1931, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1934 = base.Simd_g_i32x4_extend_low_i16x8_u(v1933)
			v1936 = base.Simd_g_i32x4_extend_high_i16x8_u(v1933)
			v1941 = base.Simd_g_v128_load64_zero(m, v325, v1920)
			v1943 = base.Simd_g_i8x16_shuffle2(v1749, v1941, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1944 = base.Simd_g_i32x4_extend_low_i16x8_u(v1943)
			v1946 = base.Simd_g_i32x4_extend_high_i16x8_u(v1943)
			v1955 = int32(8)
			v1956 = base.Simd_g_v128_load64_zero(m, v1746, v1955)
			v1958 = base.Simd_g_i8x16_shuffle2(v1749, v1956, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1964 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1958), v1782), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1958), v1782), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))
			v1966 = base.Simd_g_v128_load64_zero(m, v327, v1955)
			v1968 = base.Simd_g_i8x16_shuffle2(v1749, v1966, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1969 = base.Simd_g_i32x4_extend_low_i16x8_u(v1968)
			v1971 = base.Simd_g_i32x4_extend_high_i16x8_u(v1968)
			v1976 = base.Simd_g_v128_load64_zero(m, v325, v1955)
			v1978 = base.Simd_g_i8x16_shuffle2(v1749, v1976, base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k5), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k6))
			v1979 = base.Simd_g_i32x4_extend_low_i16x8_u(v1978)
			v1981 = base.Simd_g_i32x4_extend_high_i16x8_u(v1978)
			v1990 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1929, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1934, v1793), base.Simd_g_i32x4_mul(v1936, v1793), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1944, v1804), base.Simd_g_i32x4_mul(v1946, v1804), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v1812), v1814), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1964, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1969, v1793), base.Simd_g_i32x4_mul(v1971, v1793), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1979, v1804), base.Simd_g_i32x4_mul(v1981, v1804), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)))), v1812), v1814))
			v2008 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1934, v1854), base.Simd_g_i32x4_mul(v1936, v1854), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v1929), v1860), v1814), base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1969, v1854), base.Simd_g_i32x4_mul(v1971, v1854), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11)), v1964), v1860), v1814))
			v2027 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1929, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1944, v1876), base.Simd_g_i32x4_mul(v1946, v1876), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v1882), v1814), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1964, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1979, v1876), base.Simd_g_i32x4_mul(v1981, v1876), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k10), base.Simd_g_const(&F_UpsampleBgrLinePair_SSE41__k11))), v1882), v1814))
			v2029 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1990, v1852), base.Simd_g_i8x16_swizzle(v2008, v1873)), base.Simd_g_i8x16_swizzle(v2027, v1895))
			base.Simd_g_v128_store(m, v1748, int32(32), v2029)
			v2036 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1990, v1900), base.Simd_g_i8x16_swizzle(v2008, v1902)), base.Simd_g_i8x16_swizzle(v2027, v1905))
			base.Simd_g_v128_store(m, v1748, v1777, v2036)
			v2043 = base.Simd_g_v128_or(base.Simd_g_v128_or(base.Simd_g_i8x16_swizzle(v1990, v1910), base.Simd_g_i8x16_swizzle(v2008, v1912)), base.Simd_g_i8x16_swizzle(v2027, v1915))
			base.Simd_g_v128_store(m, v1748, v1920, v2043)
			v2046 = int32(3)
			v2051 = F_memcpy(m, l6+v1103*v2046, v1748, v1745*v2046)
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleBgrLinePair_SSE41__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleBgrLinePair_SSE41__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleBgrLinePair_SSE41__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleBgrLinePair_SSE41__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleBgrLinePair_SSE41__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleBgrLinePair_SSE41__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleBgrLinePair_SSE41__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleBgrLinePair_SSE41__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleBgrLinePair_SSE41__k8 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleBgrLinePair_SSE41__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleBgrLinePair_SSE41__k10 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleBgrLinePair_SSE41__k11 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleBgrLinePair_SSE41__k12 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleBgrLinePair_SSE41__k13 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleBgrLinePair_SSE41__k14 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleBgrLinePair_SSE41__k15 = [2]uint64{0x8f8f0c8f8f0b8f8f, 0x8f0f8f8f0e8f8f0d}
var F_UpsampleBgrLinePair_SSE41__k16 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleBgrLinePair_SSE41__k17 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleBgrLinePair_SSE41__k18 = [2]uint64{0xd8f8f0c8f8f0b8f, 0x8f8f0f8f8f0e8f8f}
var F_UpsampleBgrLinePair_SSE41__k19 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleBgrLinePair_SSE41__k20 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleBgrLinePair_SSE41__k21 = [2]uint64{0x8f0c8f8f0b8f8f0a, 0xf8f8f0e8f8f0d8f}
var F_UpsampleBgrLinePair_SSE41__k22 = [2]uint64{0x8f078f8f068f8f05, 0xa8f8f098f8f088f}
var F_UpsampleBgrLinePair_SSE41__k23 = [2]uint64{0x8f8f078f8f068f8f, 0x8f0a8f8f098f8f08}
var F_UpsampleBgrLinePair_SSE41__k24 = [2]uint64{0x78f8f068f8f058f, 0x8f8f098f8f088f8f}
var F_UpsampleBgrLinePair_SSE41__k25 = [2]uint64{0x28f8f018f8f008f, 0x8f8f048f8f038f8f}
var F_UpsampleBgrLinePair_SSE41__k26 = [2]uint64{0x8f028f8f018f8f00, 0x58f8f048f8f038f}
var F_UpsampleBgrLinePair_SSE41__k27 = [2]uint64{0x8f8f018f8f008f8f, 0x8f048f8f038f8f02}

func F_UpsampleBgraLinePair_SSE2(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) {
	mBase := m.M
	_ = mBase
	var v10 int32
	_ = v10
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v162 int32
	_ = v162
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v166 int32
	_ = v166
	var v171 int32
	_ = v171
	var v174 int32
	_ = v174
	var v177 int32
	_ = v177
	var v179 int32
	_ = v179
	var v183 int32
	_ = v183
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v193 int32
	_ = v193
	var v196 int32
	_ = v196
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v201 int32
	_ = v201
	var v204 int32
	_ = v204
	var v207 int32
	_ = v207
	var v212 int32
	_ = v212
	var v214 int32
	_ = v214
	var v221 int32
	_ = v221
	var v224 int32
	_ = v224
	var v228 int32
	_ = v228
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v244 int32
	_ = v244
	var v247 int32
	_ = v247
	var v251 int32
	_ = v251
	var v253 int32
	_ = v253
	var v256 int32
	_ = v256
	var v257 int32
	_ = v257
	var v260 int32
	_ = v260
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v274 int32
	_ = v274
	var v277 int32
	_ = v277
	var v281 int32
	_ = v281
	var v286 int32
	_ = v286
	var v288 int32
	_ = v288
	var v295 int32
	_ = v295
	var v298 int32
	_ = v298
	var v302 int32
	_ = v302
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v318 int32
	_ = v318
	var v321 int32
	_ = v321
	var v329 int32
	_ = v329
	var v331 int32
	_ = v331
	var v334 int32
	_ = v334
	var v338 int32
	_ = v338
	var v346 int32
	_ = v346
	var v358 int32
	_ = v358
	var v359 int32
	_ = v359
	var v363 int32
	_ = v363
	var v364 int32
	_ = v364
	var v380 int32
	_ = v380
	var v381 int32
	_ = v381
	var v382 base.V128
	_ = v382
	var v383 int32
	_ = v383
	var v385 base.V128
	_ = v385
	var v386 int32
	_ = v386
	var v387 base.V128
	_ = v387
	var v388 base.V128
	_ = v388
	var v390 base.V128
	_ = v390
	var v391 base.V128
	_ = v391
	var v393 base.V128
	_ = v393
	var v394 base.V128
	_ = v394
	var v396 base.V128
	_ = v396
	var v398 base.V128
	_ = v398
	var v400 base.V128
	_ = v400
	var v406 base.V128
	_ = v406
	var v407 base.V128
	_ = v407
	var v413 base.V128
	_ = v413
	var v414 base.V128
	_ = v414
	var v415 base.V128
	_ = v415
	var v416 base.V128
	_ = v416
	var v419 base.V128
	_ = v419
	var v420 base.V128
	_ = v420
	var v423 base.V128
	_ = v423
	var v424 base.V128
	_ = v424
	var v426 base.V128
	_ = v426
	var v430 base.V128
	_ = v430
	var v433 int32
	_ = v433
	var v435 base.V128
	_ = v435
	var v436 int32
	_ = v436
	var v438 base.V128
	_ = v438
	var v440 base.V128
	_ = v440
	var v441 base.V128
	_ = v441
	var v443 base.V128
	_ = v443
	var v444 base.V128
	_ = v444
	var v446 base.V128
	_ = v446
	var v447 base.V128
	_ = v447
	var v449 base.V128
	_ = v449
	var v452 base.V128
	_ = v452
	var v458 base.V128
	_ = v458
	var v459 base.V128
	_ = v459
	var v465 base.V128
	_ = v465
	var v466 base.V128
	_ = v466
	var v468 base.V128
	_ = v468
	var v472 base.V128
	_ = v472
	var v475 base.V128
	_ = v475
	var v476 base.V128
	_ = v476
	var v478 base.V128
	_ = v478
	var v482 base.V128
	_ = v482
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v509 base.V128
	_ = v509
	var v511 int32
	_ = v511
	var v512 base.V128
	_ = v512
	var v513 base.V128
	_ = v513
	var v514 base.V128
	_ = v514
	var v515 base.V128
	_ = v515
	var v516 base.V128
	_ = v516
	var v518 base.V128
	_ = v518
	var v520 base.V128
	_ = v520
	var v524 base.V128
	_ = v524
	var v526 base.V128
	_ = v526
	var v528 base.V128
	_ = v528
	var v533 base.V128
	_ = v533
	var v537 int32
	_ = v537
	var v541 base.V128
	_ = v541
	var v543 base.V128
	_ = v543
	var v544 base.V128
	_ = v544
	var v545 base.V128
	_ = v545
	var v547 base.V128
	_ = v547
	var v556 base.V128
	_ = v556
	var v557 base.V128
	_ = v557
	var v562 base.V128
	_ = v562
	var v574 base.V128
	_ = v574
	var v576 base.V128
	_ = v576
	var v578 base.V128
	_ = v578
	var v580 base.V128
	_ = v580
	var v584 base.V128
	_ = v584
	var v609 int32
	_ = v609
	var v610 int32
	_ = v610
	var v619 base.V128
	_ = v619
	var v621 int32
	_ = v621
	var v622 base.V128
	_ = v622
	var v623 base.V128
	_ = v623
	var v624 base.V128
	_ = v624
	var v625 base.V128
	_ = v625
	var v626 base.V128
	_ = v626
	var v628 base.V128
	_ = v628
	var v630 base.V128
	_ = v630
	var v634 base.V128
	_ = v634
	var v636 base.V128
	_ = v636
	var v638 base.V128
	_ = v638
	var v643 base.V128
	_ = v643
	var v647 int32
	_ = v647
	var v651 base.V128
	_ = v651
	var v653 base.V128
	_ = v653
	var v654 base.V128
	_ = v654
	var v655 base.V128
	_ = v655
	var v657 base.V128
	_ = v657
	var v666 base.V128
	_ = v666
	var v667 base.V128
	_ = v667
	var v672 base.V128
	_ = v672
	var v684 base.V128
	_ = v684
	var v686 base.V128
	_ = v686
	var v688 base.V128
	_ = v688
	var v690 base.V128
	_ = v690
	var v694 base.V128
	_ = v694
	var v703 int32
	_ = v703
	var v708 int32
	_ = v708
	var v726 int32
	_ = v726
	var v727 int32
	_ = v727
	var v750 int32
	_ = v750
	var v753 int32
	_ = v753
	var v759 int32
	_ = v759
	var v760 int32
	_ = v760
	var v762 int32
	_ = v762
	var v764 int32
	_ = v764
	var v765 int32
	_ = v765
	var v769 int32
	_ = v769
	var v770 int32
	_ = v770
	var v771 int32
	_ = v771
	var v773 int32
	_ = v773
	var v784 int32
	_ = v784
	var v809 int32
	_ = v809
	var v810 int32
	_ = v810
	var v814 int32
	_ = v814
	var v818 int32
	_ = v818
	var v819 int32
	_ = v819
	var v854 int32
	_ = v854
	var v855 int32
	_ = v855
	var v860 int64
	_ = v860
	var v863 int32
	_ = v863
	var v864 int32
	_ = v864
	var v876 int32
	_ = v876
	var v894 int32
	_ = v894
	var v895 int32
	_ = v895
	var v896 int32
	_ = v896
	var v907 int32
	_ = v907
	var v932 int32
	_ = v932
	var v933 int32
	_ = v933
	var v937 int32
	_ = v937
	var v941 int32
	_ = v941
	var v942 int32
	_ = v942
	var v977 int32
	_ = v977
	var v978 int32
	_ = v978
	var v983 int64
	_ = v983
	var v986 int32
	_ = v986
	var v987 int32
	_ = v987
	var v999 int32
	_ = v999
	var v1017 int32
	_ = v1017
	var v1018 base.V128
	_ = v1018
	var v1019 int32
	_ = v1019
	var v1020 base.V128
	_ = v1020
	var v1022 base.V128
	_ = v1022
	var v1023 base.V128
	_ = v1023
	var v1024 int32
	_ = v1024
	var v1025 base.V128
	_ = v1025
	var v1026 base.V128
	_ = v1026
	var v1028 base.V128
	_ = v1028
	var v1029 base.V128
	_ = v1029
	var v1031 base.V128
	_ = v1031
	var v1033 base.V128
	_ = v1033
	var v1035 base.V128
	_ = v1035
	var v1041 base.V128
	_ = v1041
	var v1042 base.V128
	_ = v1042
	var v1048 base.V128
	_ = v1048
	var v1049 base.V128
	_ = v1049
	var v1050 base.V128
	_ = v1050
	var v1051 base.V128
	_ = v1051
	var v1054 base.V128
	_ = v1054
	var v1055 base.V128
	_ = v1055
	var v1058 base.V128
	_ = v1058
	var v1059 base.V128
	_ = v1059
	var v1061 base.V128
	_ = v1061
	var v1065 base.V128
	_ = v1065
	var v1071 int32
	_ = v1071
	var v1073 int32
	_ = v1073
	var v1074 int32
	_ = v1074
	var v1085 int32
	_ = v1085
	var v1110 int32
	_ = v1110
	var v1111 int32
	_ = v1111
	var v1115 int32
	_ = v1115
	var v1119 int32
	_ = v1119
	var v1120 int32
	_ = v1120
	var v1155 int32
	_ = v1155
	var v1156 int32
	_ = v1156
	var v1161 int64
	_ = v1161
	var v1164 int32
	_ = v1164
	var v1165 int32
	_ = v1165
	var v1177 int32
	_ = v1177
	var v1195 int32
	_ = v1195
	var v1206 int32
	_ = v1206
	var v1231 int32
	_ = v1231
	var v1232 int32
	_ = v1232
	var v1236 int32
	_ = v1236
	var v1240 int32
	_ = v1240
	var v1241 int32
	_ = v1241
	var v1276 int32
	_ = v1276
	var v1277 int32
	_ = v1277
	var v1282 int64
	_ = v1282
	var v1285 int32
	_ = v1285
	var v1286 int32
	_ = v1286
	var v1298 int32
	_ = v1298
	var v1317 base.V128
	_ = v1317
	var v1318 int32
	_ = v1318
	var v1319 base.V128
	_ = v1319
	var v1321 base.V128
	_ = v1321
	var v1322 base.V128
	_ = v1322
	var v1324 base.V128
	_ = v1324
	var v1325 base.V128
	_ = v1325
	var v1327 base.V128
	_ = v1327
	var v1328 base.V128
	_ = v1328
	var v1330 base.V128
	_ = v1330
	var v1333 base.V128
	_ = v1333
	var v1339 base.V128
	_ = v1339
	var v1340 base.V128
	_ = v1340
	var v1346 base.V128
	_ = v1346
	var v1347 base.V128
	_ = v1347
	var v1348 base.V128
	_ = v1348
	var v1349 base.V128
	_ = v1349
	var v1352 base.V128
	_ = v1352
	var v1353 base.V128
	_ = v1353
	var v1356 base.V128
	_ = v1356
	var v1357 base.V128
	_ = v1357
	var v1359 base.V128
	_ = v1359
	var v1363 base.V128
	_ = v1363
	var v1369 int32
	_ = v1369
	var v1370 int32
	_ = v1370
	var v1372 int32
	_ = v1372
	var v1386 int32
	_ = v1386
	var v1387 int32
	_ = v1387
	var v1396 base.V128
	_ = v1396
	var v1398 int32
	_ = v1398
	var v1399 base.V128
	_ = v1399
	var v1400 base.V128
	_ = v1400
	var v1401 base.V128
	_ = v1401
	var v1402 base.V128
	_ = v1402
	var v1403 base.V128
	_ = v1403
	var v1405 base.V128
	_ = v1405
	var v1407 base.V128
	_ = v1407
	var v1411 base.V128
	_ = v1411
	var v1413 base.V128
	_ = v1413
	var v1415 base.V128
	_ = v1415
	var v1420 base.V128
	_ = v1420
	var v1424 int32
	_ = v1424
	var v1428 base.V128
	_ = v1428
	var v1430 base.V128
	_ = v1430
	var v1431 base.V128
	_ = v1431
	var v1432 base.V128
	_ = v1432
	var v1434 base.V128
	_ = v1434
	var v1443 base.V128
	_ = v1443
	var v1444 base.V128
	_ = v1444
	var v1449 base.V128
	_ = v1449
	var v1461 base.V128
	_ = v1461
	var v1463 base.V128
	_ = v1463
	var v1465 base.V128
	_ = v1465
	var v1467 base.V128
	_ = v1467
	var v1471 base.V128
	_ = v1471
	var v1480 int32
	_ = v1480
	var v1485 int32
	_ = v1485
	var v1489 int32
	_ = v1489
	var v1503 int32
	_ = v1503
	var v1504 int32
	_ = v1504
	var v1513 base.V128
	_ = v1513
	var v1515 int32
	_ = v1515
	var v1516 base.V128
	_ = v1516
	var v1517 base.V128
	_ = v1517
	var v1518 base.V128
	_ = v1518
	var v1519 base.V128
	_ = v1519
	var v1520 base.V128
	_ = v1520
	var v1522 base.V128
	_ = v1522
	var v1524 base.V128
	_ = v1524
	var v1528 base.V128
	_ = v1528
	var v1530 base.V128
	_ = v1530
	var v1532 base.V128
	_ = v1532
	var v1537 base.V128
	_ = v1537
	var v1541 int32
	_ = v1541
	var v1545 base.V128
	_ = v1545
	var v1547 base.V128
	_ = v1547
	var v1548 base.V128
	_ = v1548
	var v1549 base.V128
	_ = v1549
	var v1551 base.V128
	_ = v1551
	var v1560 base.V128
	_ = v1560
	var v1561 base.V128
	_ = v1561
	var v1566 base.V128
	_ = v1566
	var v1578 base.V128
	_ = v1578
	var v1580 base.V128
	_ = v1580
	var v1582 base.V128
	_ = v1582
	var v1584 base.V128
	_ = v1584
	var v1588 base.V128
	_ = v1588
	var v1602 int32
	_ = v1602
	var v1616 int32
	_ = v1616
	var v1617 int32
	_ = v1617
	var v1626 base.V128
	_ = v1626
	var v1628 int32
	_ = v1628
	var v1629 base.V128
	_ = v1629
	var v1630 base.V128
	_ = v1630
	var v1631 base.V128
	_ = v1631
	var v1632 base.V128
	_ = v1632
	var v1633 base.V128
	_ = v1633
	var v1635 base.V128
	_ = v1635
	var v1637 base.V128
	_ = v1637
	var v1641 base.V128
	_ = v1641
	var v1643 base.V128
	_ = v1643
	var v1645 base.V128
	_ = v1645
	var v1650 base.V128
	_ = v1650
	var v1654 int32
	_ = v1654
	var v1658 base.V128
	_ = v1658
	var v1660 base.V128
	_ = v1660
	var v1661 base.V128
	_ = v1661
	var v1662 base.V128
	_ = v1662
	var v1664 base.V128
	_ = v1664
	var v1673 base.V128
	_ = v1673
	var v1674 base.V128
	_ = v1674
	var v1679 base.V128
	_ = v1679
	var v1691 base.V128
	_ = v1691
	var v1693 base.V128
	_ = v1693
	var v1695 base.V128
	_ = v1695
	var v1697 base.V128
	_ = v1697
	var v1701 base.V128
	_ = v1701
	var v1710 int32
	_ = v1710
	var v1711 int32
	_ = v1711
	var v1714 int32
	_ = v1714
	var v1715 int32
	_ = v1715
	var v1717 int32
	_ = v1717
	v10 = int32(0)
	v33 = m.G0
	v35 = v33 - int32(528)
	m.G0 = v35
	base.MemoryFill(m, v35+int32(64), v10, int32(463))
	v162 = int32(255)
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+3)) = uint8(v162)
	v164 = int32(1)
	v165 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l5))))
	v166 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l3))))
	v171 = int32(base.Ui32(v165+v166)>>(uint(v164)%32)) + v164
	v174 = int32(base.Ui32(v171+v166) >> (uint(v164) % 32))
	v177 = int32(8)
	v179 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0))))
	v183 = int32(base.Ui32(v179*int32(_a_F_UpsampleBgraLinePair_SSE2_0)) >> (uint(v177) % 32))
	v184 = int32(base.Ui32(v174*int32(_a_F_UpsampleBgraLinePair_SSE2_1))>>(uint(v177)%32)) + v183
	v186 = v184 + int32(-14234)
	if base.Ui32(v184) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_2)) {
		v193 = int32(0)
	} else {
		v193 = v162
	}
	if base.Ui32(v186) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_3)) {
		v196 = int32(base.Ui32(v186) >> (uint(int32(6)) % 32))
	} else {
		v196 = v193
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+2)) = uint8(v196)
	v198 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l4))))
	v199 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l2))))
	v201 = int32(1)
	v204 = int32(base.Ui32(v198+v199)>>(uint(v201)%32)) + v201
	v207 = int32(base.Ui32(v204+v199) >> (uint(v201) % 32))
	v212 = int32(base.Ui32(v207*int32(_a_F_UpsampleBgraLinePair_SSE2_4))>>(uint(int32(8))%32)) + v183
	v214 = v212 + int32(-17685)
	if base.Ui32(v212) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_5)) {
		v221 = int32(0)
	} else {
		v221 = int32(255)
	}
	if base.Ui32(v214) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_3)) {
		v224 = int32(base.Ui32(v214) >> (uint(int32(6)) % 32))
	} else {
		v224 = v221
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6))) = uint8(v224)
	v228 = int32(8)
	v235 = v183 - (int32(base.Ui32(v207*int32(_a_F_UpsampleBgraLinePair_SSE2_6))>>(uint(v228)%32)) + int32(base.Ui32(v174*int32(_a_F_UpsampleBgraLinePair_SSE2_7))>>(uint(v228)%32)))
	v237 = v235 + int32(_a_F_UpsampleBgraLinePair_SSE2_8)
	if v235 < int32(-8708) {
		v244 = int32(0)
	} else {
		v244 = int32(255)
	}
	if base.Ui32(v237) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_3)) {
		v247 = int32(base.Ui32(v237) >> (uint(int32(6)) % 32))
	} else {
		v247 = v244
	}
	*(*uint8)(unsafe.Add(mBase, uint32(l6)+1)) = uint8(v247)
	if l1 == int32(0) {
	} else {
		v251 = int32(255)
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+3)) = uint8(v251)
		v253 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l1))))
		v256 = int32(8)
		v257 = int32(base.Ui32(v253*int32(_a_F_UpsampleBgraLinePair_SSE2_0)) >> (uint(v256) % 32))
		v260 = int32(base.Ui32(v171+v165) >> (uint(int32(1)) % 32))
		v265 = v257 + int32(base.Ui32(v260*int32(_a_F_UpsampleBgraLinePair_SSE2_1))>>(uint(v256)%32))
		v267 = v265 + int32(-14234)
		if base.Ui32(v265) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_2)) {
			v274 = int32(0)
		} else {
			v274 = v251
		}
		if base.Ui32(v267) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_3)) {
			v277 = int32(base.Ui32(v267) >> (uint(int32(6)) % 32))
		} else {
			v277 = v274
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v277)
		v281 = int32(base.Ui32(v204+v198) >> (uint(int32(1)) % 32))
		v286 = v257 + int32(base.Ui32(v281*int32(_a_F_UpsampleBgraLinePair_SSE2_4))>>(uint(int32(8))%32))
		v288 = v286 + int32(-17685)
		if base.Ui32(v286) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_5)) {
			v295 = int32(0)
		} else {
			v295 = int32(255)
		}
		if base.Ui32(v288) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_3)) {
			v298 = int32(base.Ui32(v288) >> (uint(int32(6)) % 32))
		} else {
			v298 = v295
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v298)
		v302 = int32(8)
		v309 = v257 - (int32(base.Ui32(v260*int32(_a_F_UpsampleBgraLinePair_SSE2_7))>>(uint(v302)%32)) + int32(base.Ui32(v281*int32(_a_F_UpsampleBgraLinePair_SSE2_6))>>(uint(v302)%32)))
		v311 = v309 + int32(_a_F_UpsampleBgraLinePair_SSE2_8)
		if v309 < int32(-8708) {
			v318 = int32(0)
		} else {
			v318 = int32(255)
		}
		if base.Ui32(v311) < base.Ui32(int32(_a_F_UpsampleBgraLinePair_SSE2_3)) {
			v321 = int32(base.Ui32(v311) >> (uint(int32(6)) % 32))
		} else {
			v321 = v318
		}
		*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v321)
	}
	v329 = v35 + int32(96)
	v331 = v35 + int32(64)
	if l8 < int32(34) {
		v726 = v10
		v727 = v164
	} else {
		v334 = int32(4)
		v338 = int32(1)
		v346 = int32(0)
		v358 = v346
		v359 = l7 + v334
		v363 = v346
		v364 = l6 + v334
		for {
			v380 = l4 + v358
			v381 = int32(0)
			v382 = base.Simd_g_v128_load_rng(m, v380, v381, int32(0), int32(17))
			v383 = l2 + v358
			v385 = base.Simd_g_v128_load_rng(m, v383, v381, int32(0), int32(17))
			v386 = int32(1)
			v387 = base.Simd_g_v128_load_nc(m, v380, v386)
			v388 = base.Simd_g_i8x16_avgr_u(v385, v387)
			v390 = base.Simd_g_v128_load_nc(m, v383, v386)
			v391 = base.Simd_g_i8x16_avgr_u(v382, v390)
			v393 = base.Simd_g_v128_xor(v387, v385)
			v394 = base.Simd_g_v128_xor(v382, v390)
			v396 = base.Simd_g_v128_xor(v388, v391)
			v398 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k0)
			v400 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v388, v391), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v393, v394), v396), v398))
			v406 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v400, v388), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v400, v388), base.Simd_g_v128_and(v396, v393)), v398))
			v407 = base.Simd_g_i8x16_avgr_u(v382, v406)
			v413 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v400, v391), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v400, v391), base.Simd_g_v128_and(v396, v394)), v398))
			v414 = base.Simd_g_i8x16_avgr_u(v387, v413)
			v415 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k1)
			v416 = base.Simd_g_i8x16_shuffle2(v407, v414, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(80), v416)
			v419 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
			v420 = base.Simd_g_i8x16_shuffle2(v407, v414, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(64), v420)
			v423 = base.Simd_g_i8x16_avgr_u(v385, v413)
			v424 = base.Simd_g_i8x16_avgr_u(v390, v406)
			v426 = base.Simd_g_i8x16_shuffle2(v423, v424, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(16), v426)
			v430 = base.Simd_g_i8x16_shuffle2(v423, v424, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, v381, v430)
			v433 = l5 + v358
			v435 = base.Simd_g_v128_load_rng(m, v433, v381, int32(0), int32(17))
			v436 = l3 + v358
			v438 = base.Simd_g_v128_load_rng(m, v436, v381, int32(0), int32(17))
			v440 = base.Simd_g_v128_load_nc(m, v433, v386)
			v441 = base.Simd_g_i8x16_avgr_u(v438, v440)
			v443 = base.Simd_g_v128_load_nc(m, v436, v386)
			v444 = base.Simd_g_i8x16_avgr_u(v435, v443)
			v446 = base.Simd_g_v128_xor(v440, v438)
			v447 = base.Simd_g_v128_xor(v435, v443)
			v449 = base.Simd_g_v128_xor(v441, v444)
			v452 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v441, v444), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v446, v447), v449), v398))
			v458 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v452, v441), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v452, v441), base.Simd_g_v128_and(v449, v446)), v398))
			v459 = base.Simd_g_i8x16_avgr_u(v435, v458)
			v465 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v452, v444), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v452, v444), base.Simd_g_v128_and(v449, v447)), v398))
			v466 = base.Simd_g_i8x16_avgr_u(v440, v465)
			v468 = base.Simd_g_i8x16_shuffle2(v459, v466, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(112), v468)
			v472 = base.Simd_g_i8x16_shuffle2(v459, v466, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(96), v472)
			v475 = base.Simd_g_i8x16_avgr_u(v438, v465)
			v476 = base.Simd_g_i8x16_avgr_u(v443, v458)
			v478 = base.Simd_g_i8x16_shuffle2(v475, v476, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
			base.Simd_g_v128_store(m, v331, int32(48), v478)
			v482 = base.Simd_g_i8x16_shuffle2(v475, v476, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
			base.Simd_g_v128_store(m, v331, int32(32), v482)
			v499 = v364
			v500 = v381
			for {
				v509 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k7)
				v511 = int32(0)
				v512 = base.Simd_g_v128_load64_zero(m, v331+v500, v511)
				v513 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
				v514 = base.Simd_g_i8x16_shuffle2(v509, v512, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v515 = base.Simd_g_i32x4_extend_low_i16x8_u(v514)
				v516 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k8)
				v518 = base.Simd_g_i32x4_extend_high_i16x8_u(v514)
				v520 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k9)
				v524 = base.Simd_g_v128_load64_zero(m, l0+v338+v363+v500, v511)
				v526 = base.Simd_g_i8x16_shuffle2(v509, v524, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v528 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k10)
				v533 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v526), v528), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v526), v528), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))
				v537 = int32(6)
				v541 = base.Simd_g_v128_load64_zero(m, v329+v500, v511)
				v543 = base.Simd_g_i8x16_shuffle2(v509, v541, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v544 = base.Simd_g_i32x4_extend_low_i16x8_u(v543)
				v545 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k13)
				v547 = base.Simd_g_i32x4_extend_high_i16x8_u(v543)
				v556 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v515, v516), base.Simd_g_i32x4_mul(v518, v516), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), v533), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k14)), v537), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v533, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v544, v545), base.Simd_g_i32x4_mul(v547, v545), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k15)), v537))
				v557 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k16)
				v562 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k17)
				v574 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v533, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v515, v557), base.Simd_g_i32x4_mul(v518, v557), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v544, v562), base.Simd_g_i32x4_mul(v547, v562), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k18)), v537), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k19))
				v576 = base.Simd_g_i8x16_shuffle2(v556, v574, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v578 = base.Simd_g_i8x16_shuffle2(v556, v574, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
				v580 = base.Simd_g_i8x16_shuffle2(v576, v578, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v499, int32(16), v580)
				v584 = base.Simd_g_i8x16_shuffle2(v576, v578, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v499, v511, v584)
				if base.Ui32(v500) < base.Ui32(int32(24)) {
					v499 = v499 + int32(32)
					v500 = v500 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			if l1 == int32(0) {
			} else {
				v609 = v359
				v610 = int32(0)
				for {
					v619 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k7)
					v621 = int32(0)
					v622 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v610, v621)
					v623 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
					v624 = base.Simd_g_i8x16_shuffle2(v619, v622, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
					v625 = base.Simd_g_i32x4_extend_low_i16x8_u(v624)
					v626 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k8)
					v628 = base.Simd_g_i32x4_extend_high_i16x8_u(v624)
					v630 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k9)
					v634 = base.Simd_g_v128_load64_zero(m, l1+v338+v363+v610, v621)
					v636 = base.Simd_g_i8x16_shuffle2(v619, v634, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
					v638 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k10)
					v643 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v636), v638), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v636), v638), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))
					v647 = int32(6)
					v651 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v610, v621)
					v653 = base.Simd_g_i8x16_shuffle2(v619, v651, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
					v654 = base.Simd_g_i32x4_extend_low_i16x8_u(v653)
					v655 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k13)
					v657 = base.Simd_g_i32x4_extend_high_i16x8_u(v653)
					v666 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v625, v626), base.Simd_g_i32x4_mul(v628, v626), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), v643), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k14)), v647), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v643, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v654, v655), base.Simd_g_i32x4_mul(v657, v655), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k15)), v647))
					v667 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k16)
					v672 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k17)
					v684 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v643, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v625, v667), base.Simd_g_i32x4_mul(v628, v667), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v654, v672), base.Simd_g_i32x4_mul(v657, v672), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k18)), v647), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k19))
					v686 = base.Simd_g_i8x16_shuffle2(v666, v684, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
					v688 = base.Simd_g_i8x16_shuffle2(v666, v684, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
					v690 = base.Simd_g_i8x16_shuffle2(v686, v688, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k21))
					base.Simd_g_v128_store(m, v609, int32(16), v690)
					v694 = base.Simd_g_i8x16_shuffle2(v686, v688, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k23))
					base.Simd_g_v128_store(m, v609, v621, v694)
					if base.Ui32(v610) < base.Ui32(int32(24)) {
						v609 = v609 + int32(32)
						v610 = v610 + int32(8)
						continue
					} else {
						break
					}
					break
				}
			}
			v703 = int32(128)
			v708 = v358 + int32(16)
			if v363+int32(66) <= l8 {
				v358 = v708
				v359 = v359 + v703
				v363 = v363 + int32(32)
				v364 = v364 + v703
				continue
			} else {
				break
			}
			break
		}
		v726 = v708
		v727 = v363 + int32(33)
	}
	if l8 < int32(2) {
	} else {
		v750 = int32(32)
		v753 = int32(1)
		v759 = int32(base.Ui32(l8+v753)>>(uint(v753)%32)) - int32(base.Ui32(v727)>>(uint(v753)%32))
		v760 = F_memcpy(m, v35+v750, l2+v726, v759)
		mBase = m.M
		v762 = F_memcpy(m, v35, l4+v726, v759)
		mBase = m.M
		v764 = v762 + v750
		v765 = v764 + v759
		v769 = v759 + int32(-1)
		v770 = v764 + v769
		v771 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v770))))
		v773 = int32(17) - v759
		if base.Ui32(v773) < base.Ui32(int32(33)) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v765))) = uint8(v771)
				v784 = v765 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-1)))) = uint8(v771)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+2)) = uint8(v771)
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+1)) = uint8(v771)
					*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-3)))) = uint8(v771)
					*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-2)))) = uint8(v771)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v765)+3)) = uint8(v771)
						*(*uint8)(unsafe.Add(mBase, uint32(v784+int32(-4)))) = uint8(v771)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v809 = (int32(0) - v765) & int32(3)
							v810 = v765 + v809
							v814 = v771 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v810))) = v814
							v818 = (v773 - v809) & int32(60)
							v819 = v810 + v818
							*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-4)))) = v814
							if base.Ui32(v818) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v810)+8)) = v814
								*(*int32)(unsafe.Add(mBase, uint32(v810)+4)) = v814
								*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-8)))) = v814
								*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-12)))) = v814
								if base.Ui32(v818) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v810)+24)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v810)+20)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v810)+16)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v810)+12)) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-16)))) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-20)))) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-24)))) = v814
									*(*int32)(unsafe.Add(mBase, uint32(v819+int32(-28)))) = v814
									v854 = v810&int32(4) | int32(24)
									v855 = v818 - v854
									if base.Ui32(v855) < base.Ui32(int32(32)) {
									} else {
										v860 = base.I64_extend_i32_u(v814) * int64(4294967297)
										v863 = v855
										v864 = v810 + v854
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v864)+24)) = v860
											*(*int64)(unsafe.Add(mBase, uint32(v864)+16)) = v860
											*(*int64)(unsafe.Add(mBase, uint32(v864)+8)) = v860
											*(*int64)(unsafe.Add(mBase, uint32(v864))) = v860
											v876 = v863 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v876) {
												v863 = v876
												v864 = v864 + int32(32)
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
			base.MemoryFill(m, v765, v771, v773)
		}
		v894 = v762 + v759
		v895 = v762 + v769
		v896 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v895))))
		if base.Ui32(v773) < base.Ui32(int32(33)) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v894))) = uint8(v896)
				v907 = v894 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-1)))) = uint8(v896)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+2)) = uint8(v896)
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+1)) = uint8(v896)
					*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-3)))) = uint8(v896)
					*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-2)))) = uint8(v896)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v894)+3)) = uint8(v896)
						*(*uint8)(unsafe.Add(mBase, uint32(v907+int32(-4)))) = uint8(v896)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v932 = (int32(0) - v894) & int32(3)
							v933 = v894 + v932
							v937 = v896 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v933))) = v937
							v941 = (v773 - v932) & int32(60)
							v942 = v933 + v941
							*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-4)))) = v937
							if base.Ui32(v941) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v933)+8)) = v937
								*(*int32)(unsafe.Add(mBase, uint32(v933)+4)) = v937
								*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-8)))) = v937
								*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-12)))) = v937
								if base.Ui32(v941) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v933)+24)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v933)+20)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v933)+16)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v933)+12)) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-16)))) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-20)))) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-24)))) = v937
									*(*int32)(unsafe.Add(mBase, uint32(v942+int32(-28)))) = v937
									v977 = v933&int32(4) | int32(24)
									v978 = v941 - v977
									if base.Ui32(v978) < base.Ui32(int32(32)) {
									} else {
										v983 = base.I64_extend_i32_u(v937) * int64(4294967297)
										v986 = v978
										v987 = v933 + v977
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v987)+24)) = v983
											*(*int64)(unsafe.Add(mBase, uint32(v987)+16)) = v983
											*(*int64)(unsafe.Add(mBase, uint32(v987)+8)) = v983
											*(*int64)(unsafe.Add(mBase, uint32(v987))) = v983
											v999 = v986 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v999) {
												v986 = v999
												v987 = v987 + int32(32)
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
			base.MemoryFill(m, v894, v896, v773)
		}
		v1017 = int32(0)
		v1018 = base.Simd_g_v128_load_rng(m, v762, v1017, int32(0), int32(49))
		v1019 = int32(32)
		v1020 = base.Simd_g_v128_load_nc(m, v762, v1019)
		v1022 = base.Simd_g_v128_load_nc(m, v762, int32(1))
		v1023 = base.Simd_g_i8x16_avgr_u(v1020, v1022)
		v1024 = int32(33)
		v1025 = base.Simd_g_v128_load_nc(m, v762, v1024)
		v1026 = base.Simd_g_i8x16_avgr_u(v1018, v1025)
		v1028 = base.Simd_g_v128_xor(v1022, v1020)
		v1029 = base.Simd_g_v128_xor(v1018, v1025)
		v1031 = base.Simd_g_v128_xor(v1023, v1026)
		v1033 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k0)
		v1035 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1023, v1026), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1028, v1029), v1031), v1033))
		v1041 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1035, v1023), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1035, v1023), base.Simd_g_v128_and(v1031, v1028)), v1033))
		v1042 = base.Simd_g_i8x16_avgr_u(v1018, v1041)
		v1048 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1035, v1026), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1035, v1026), base.Simd_g_v128_and(v1031, v1029)), v1033))
		v1049 = base.Simd_g_i8x16_avgr_u(v1022, v1048)
		v1050 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k1)
		v1051 = base.Simd_g_i8x16_shuffle2(v1042, v1049, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(80), v1051)
		v1054 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
		v1055 = base.Simd_g_i8x16_shuffle2(v1042, v1049, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, int32(64), v1055)
		v1058 = base.Simd_g_i8x16_avgr_u(v1020, v1048)
		v1059 = base.Simd_g_i8x16_avgr_u(v1025, v1041)
		v1061 = base.Simd_g_i8x16_shuffle2(v1058, v1059, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(16), v1061)
		v1065 = base.Simd_g_i8x16_shuffle2(v1058, v1059, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, v1017, v1065)
		v1071 = F_memcpy(m, v762+v1019, l3+v726, v759)
		mBase = m.M
		v1073 = F_memcpy(m, v762, l5+v726, v759)
		mBase = m.M
		v1074 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v770))))
		if base.Ui32(v773) < base.Ui32(v1024) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v765))) = uint8(v1074)
				v1085 = v765 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-1)))) = uint8(v1074)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+2)) = uint8(v1074)
					*(*uint8)(unsafe.Add(mBase, uint32(v765)+1)) = uint8(v1074)
					*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-3)))) = uint8(v1074)
					*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-2)))) = uint8(v1074)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v765)+3)) = uint8(v1074)
						*(*uint8)(unsafe.Add(mBase, uint32(v1085+int32(-4)))) = uint8(v1074)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v1110 = (int32(0) - v765) & int32(3)
							v1111 = v765 + v1110
							v1115 = v1074 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1111))) = v1115
							v1119 = (v773 - v1110) & int32(60)
							v1120 = v1111 + v1119
							*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-4)))) = v1115
							if base.Ui32(v1119) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1111)+8)) = v1115
								*(*int32)(unsafe.Add(mBase, uint32(v1111)+4)) = v1115
								*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-8)))) = v1115
								*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-12)))) = v1115
								if base.Ui32(v1119) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+24)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+20)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+16)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1111)+12)) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-16)))) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-20)))) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-24)))) = v1115
									*(*int32)(unsafe.Add(mBase, uint32(v1120+int32(-28)))) = v1115
									v1155 = v1111&int32(4) | int32(24)
									v1156 = v1119 - v1155
									if base.Ui32(v1156) < base.Ui32(int32(32)) {
									} else {
										v1161 = base.I64_extend_i32_u(v1115) * int64(4294967297)
										v1164 = v1156
										v1165 = v1111 + v1155
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1165)+24)) = v1161
											*(*int64)(unsafe.Add(mBase, uint32(v1165)+16)) = v1161
											*(*int64)(unsafe.Add(mBase, uint32(v1165)+8)) = v1161
											*(*int64)(unsafe.Add(mBase, uint32(v1165))) = v1161
											v1177 = v1164 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1177) {
												v1164 = v1177
												v1165 = v1165 + int32(32)
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
			base.MemoryFill(m, v765, v1074, v773)
		}
		v1195 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v895))))
		if base.Ui32(v773) < base.Ui32(int32(33)) {
			if v773 == int32(0) {
			} else {
				*(*uint8)(unsafe.Add(mBase, uint32(v894))) = uint8(v1195)
				v1206 = v894 + v773
				*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-1)))) = uint8(v1195)
				if base.Ui32(v773) < base.Ui32(int32(3)) {
				} else {
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+2)) = uint8(v1195)
					*(*uint8)(unsafe.Add(mBase, uint32(v894)+1)) = uint8(v1195)
					*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-3)))) = uint8(v1195)
					*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-2)))) = uint8(v1195)
					if base.Ui32(v773) < base.Ui32(int32(7)) {
					} else {
						*(*uint8)(unsafe.Add(mBase, uint32(v894)+3)) = uint8(v1195)
						*(*uint8)(unsafe.Add(mBase, uint32(v1206+int32(-4)))) = uint8(v1195)
						if base.Ui32(v773) < base.Ui32(int32(9)) {
						} else {
							v1231 = (int32(0) - v894) & int32(3)
							v1232 = v894 + v1231
							v1236 = v1195 & int32(255) * int32(16843009)
							*(*int32)(unsafe.Add(mBase, uint32(v1232))) = v1236
							v1240 = (v773 - v1231) & int32(60)
							v1241 = v1232 + v1240
							*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-4)))) = v1236
							if base.Ui32(v1240) < base.Ui32(int32(9)) {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(v1232)+8)) = v1236
								*(*int32)(unsafe.Add(mBase, uint32(v1232)+4)) = v1236
								*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-8)))) = v1236
								*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-12)))) = v1236
								if base.Ui32(v1240) < base.Ui32(int32(25)) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+24)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+20)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+16)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1232)+12)) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-16)))) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-20)))) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-24)))) = v1236
									*(*int32)(unsafe.Add(mBase, uint32(v1241+int32(-28)))) = v1236
									v1276 = v1232&int32(4) | int32(24)
									v1277 = v1240 - v1276
									if base.Ui32(v1277) < base.Ui32(int32(32)) {
									} else {
										v1282 = base.I64_extend_i32_u(v1236) * int64(4294967297)
										v1285 = v1277
										v1286 = v1232 + v1276
										for {
											*(*int64)(unsafe.Add(mBase, uint32(v1286)+24)) = v1282
											*(*int64)(unsafe.Add(mBase, uint32(v1286)+16)) = v1282
											*(*int64)(unsafe.Add(mBase, uint32(v1286)+8)) = v1282
											*(*int64)(unsafe.Add(mBase, uint32(v1286))) = v1282
											v1298 = v1285 + int32(-32)
											if base.Ui32(int32(31)) < base.Ui32(v1298) {
												v1285 = v1298
												v1286 = v1286 + int32(32)
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
			base.MemoryFill(m, v894, v1195, v773)
		}
		v1317 = base.Simd_g_v128_load_rng(m, v1073, int32(0), int32(0), int32(49))
		v1318 = int32(32)
		v1319 = base.Simd_g_v128_load_nc(m, v1073, v1318)
		v1321 = base.Simd_g_v128_load_nc(m, v1073, int32(1))
		v1322 = base.Simd_g_i8x16_avgr_u(v1319, v1321)
		v1324 = base.Simd_g_v128_load_nc(m, v1073, int32(33))
		v1325 = base.Simd_g_i8x16_avgr_u(v1317, v1324)
		v1327 = base.Simd_g_v128_xor(v1321, v1319)
		v1328 = base.Simd_g_v128_xor(v1317, v1324)
		v1330 = base.Simd_g_v128_xor(v1322, v1325)
		v1333 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1322, v1325), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_or(v1327, v1328), v1330), v1033))
		v1339 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1333, v1322), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1333, v1322), base.Simd_g_v128_and(v1330, v1327)), v1033))
		v1340 = base.Simd_g_i8x16_avgr_u(v1317, v1339)
		v1346 = base.Simd_g_i8x16_sub(base.Simd_g_i8x16_avgr_u(v1333, v1325), base.Simd_g_v128_and(base.Simd_g_v128_or(base.Simd_g_v128_xor(v1333, v1325), base.Simd_g_v128_and(v1330, v1328)), v1033))
		v1347 = base.Simd_g_i8x16_avgr_u(v1321, v1346)
		v1348 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k1)
		v1349 = base.Simd_g_i8x16_shuffle2(v1340, v1347, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(112), v1349)
		v1352 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
		v1353 = base.Simd_g_i8x16_shuffle2(v1340, v1347, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, int32(96), v1353)
		v1356 = base.Simd_g_i8x16_avgr_u(v1319, v1346)
		v1357 = base.Simd_g_i8x16_avgr_u(v1324, v1339)
		v1359 = base.Simd_g_i8x16_shuffle2(v1356, v1357, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
		base.Simd_g_v128_store(m, v331, int32(48), v1359)
		v1363 = base.Simd_g_i8x16_shuffle2(v1356, v1357, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
		base.Simd_g_v128_store(m, v331, v1318, v1363)
		v1369 = l8 - v727
		v1370 = F_memcpy(m, v35+int32(448), l0+v727, v1369)
		mBase = m.M
		v1372 = v35 + int32(192)
		if l1 != 0 {
			v1489 = F_memcpy(m, v35+int32(480), l1+v727, v1369)
			mBase = m.M
			v1503 = v1372
			v1504 = int32(0)
			for {
				v1513 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k7)
				v1515 = int32(0)
				v1516 = base.Simd_g_v128_load64_zero(m, v331+v1504, v1515)
				v1517 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
				v1518 = base.Simd_g_i8x16_shuffle2(v1513, v1516, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1519 = base.Simd_g_i32x4_extend_low_i16x8_u(v1518)
				v1520 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k8)
				v1522 = base.Simd_g_i32x4_extend_high_i16x8_u(v1518)
				v1524 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k9)
				v1528 = base.Simd_g_v128_load64_zero(m, v1370+v1504, v1515)
				v1530 = base.Simd_g_i8x16_shuffle2(v1513, v1528, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1532 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k10)
				v1537 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1530), v1532), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1530), v1532), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))
				v1541 = int32(6)
				v1545 = base.Simd_g_v128_load64_zero(m, v329+v1504, v1515)
				v1547 = base.Simd_g_i8x16_shuffle2(v1513, v1545, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1548 = base.Simd_g_i32x4_extend_low_i16x8_u(v1547)
				v1549 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k13)
				v1551 = base.Simd_g_i32x4_extend_high_i16x8_u(v1547)
				v1560 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1519, v1520), base.Simd_g_i32x4_mul(v1522, v1520), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), v1537), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k14)), v1541), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1537, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1548, v1549), base.Simd_g_i32x4_mul(v1551, v1549), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k15)), v1541))
				v1561 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k16)
				v1566 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k17)
				v1578 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1537, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1519, v1561), base.Simd_g_i32x4_mul(v1522, v1561), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1548, v1566), base.Simd_g_i32x4_mul(v1551, v1566), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k18)), v1541), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k19))
				v1580 = base.Simd_g_i8x16_shuffle2(v1560, v1578, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1582 = base.Simd_g_i8x16_shuffle2(v1560, v1578, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
				v1584 = base.Simd_g_i8x16_shuffle2(v1580, v1582, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1503, int32(16), v1584)
				v1588 = base.Simd_g_i8x16_shuffle2(v1580, v1582, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v1503, v1515, v1588)
				if base.Ui32(v1504) < base.Ui32(int32(24)) {
					v1503 = v1503 + int32(32)
					v1504 = v1504 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1602 = v35 + int32(320)
			v1616 = v1602
			v1617 = int32(0)
			for {
				v1626 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k7)
				v1628 = int32(0)
				v1629 = base.Simd_g_v128_load64_zero(m, v35+int32(128)+v1617, v1628)
				v1630 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
				v1631 = base.Simd_g_i8x16_shuffle2(v1626, v1629, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1632 = base.Simd_g_i32x4_extend_low_i16x8_u(v1631)
				v1633 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k8)
				v1635 = base.Simd_g_i32x4_extend_high_i16x8_u(v1631)
				v1637 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k9)
				v1641 = base.Simd_g_v128_load64_zero(m, v1489+v1617, v1628)
				v1643 = base.Simd_g_i8x16_shuffle2(v1626, v1641, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1645 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k10)
				v1650 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1643), v1645), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1643), v1645), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))
				v1654 = int32(6)
				v1658 = base.Simd_g_v128_load64_zero(m, v35+int32(160)+v1617, v1628)
				v1660 = base.Simd_g_i8x16_shuffle2(v1626, v1658, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1661 = base.Simd_g_i32x4_extend_low_i16x8_u(v1660)
				v1662 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k13)
				v1664 = base.Simd_g_i32x4_extend_high_i16x8_u(v1660)
				v1673 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1632, v1633), base.Simd_g_i32x4_mul(v1635, v1633), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), v1650), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k14)), v1654), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1650, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1661, v1662), base.Simd_g_i32x4_mul(v1664, v1662), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k15)), v1654))
				v1674 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k16)
				v1679 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k17)
				v1691 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1650, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1632, v1674), base.Simd_g_i32x4_mul(v1635, v1674), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1661, v1679), base.Simd_g_i32x4_mul(v1664, v1679), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k18)), v1654), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k19))
				v1693 = base.Simd_g_i8x16_shuffle2(v1673, v1691, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1695 = base.Simd_g_i8x16_shuffle2(v1673, v1691, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
				v1697 = base.Simd_g_i8x16_shuffle2(v1693, v1695, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1616, int32(16), v1697)
				v1701 = base.Simd_g_i8x16_shuffle2(v1693, v1695, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v1616, v1628, v1701)
				if base.Ui32(v1617) < base.Ui32(int32(24)) {
					v1616 = v1616 + int32(32)
					v1617 = v1617 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1710 = int32(2)
			v1711 = v727 << (uint(v1710) % 32)
			v1714 = v1369 << (uint(v1710) % 32)
			v1715 = F_memcpy(m, l6+v1711, v1372, v1714)
			mBase = m.M
			v1717 = F_memcpy(m, l7+v1711, v1602, v1714)
			mBase = m.M
		} else {
			v1386 = v1372
			v1387 = int32(0)
			for {
				v1396 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k7)
				v1398 = int32(0)
				v1399 = base.Simd_g_v128_load64_zero(m, v331+v1387, v1398)
				v1400 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k4)
				v1401 = base.Simd_g_i8x16_shuffle2(v1396, v1399, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1402 = base.Simd_g_i32x4_extend_low_i16x8_u(v1401)
				v1403 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k8)
				v1405 = base.Simd_g_i32x4_extend_high_i16x8_u(v1401)
				v1407 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k9)
				v1411 = base.Simd_g_v128_load64_zero(m, v1370+v1387, v1398)
				v1413 = base.Simd_g_i8x16_shuffle2(v1396, v1411, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1415 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k10)
				v1420 = base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v1413), v1415), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_high_i16x8_u(v1413), v1415), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))
				v1424 = int32(6)
				v1428 = base.Simd_g_v128_load64_zero(m, v329+v1387, v1398)
				v1430 = base.Simd_g_i8x16_shuffle2(v1396, v1428, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1431 = base.Simd_g_i32x4_extend_low_i16x8_u(v1430)
				v1432 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k13)
				v1434 = base.Simd_g_i32x4_extend_high_i16x8_u(v1430)
				v1443 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_u(base.Simd_g_i16x8_sub_sat_u(base.Simd_g_i16x8_add_sat_u(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1402, v1403), base.Simd_g_i32x4_mul(v1405, v1403), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), v1420), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k14)), v1424), base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_add(v1420, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1431, v1432), base.Simd_g_i32x4_mul(v1434, v1432), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k15)), v1424))
				v1444 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k16)
				v1449 = base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k17)
				v1461 = base.Simd_g_i8x16_narrow_i16x8_u(base.Simd_g_i16x8_shr_s(base.Simd_g_i16x8_add(base.Simd_g_i16x8_sub(v1420, base.Simd_g_i16x8_add(base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1402, v1444), base.Simd_g_i32x4_mul(v1405, v1444), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)), base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_mul(v1431, v1449), base.Simd_g_i32x4_mul(v1434, v1449), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k11), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k12)))), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k18)), v1424), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k19))
				v1463 = base.Simd_g_i8x16_shuffle2(v1443, v1461, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k5), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k6))
				v1465 = base.Simd_g_i8x16_shuffle2(v1443, v1461, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k2), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k3))
				v1467 = base.Simd_g_i8x16_shuffle2(v1463, v1465, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k20), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k21))
				base.Simd_g_v128_store(m, v1386, int32(16), v1467)
				v1471 = base.Simd_g_i8x16_shuffle2(v1463, v1465, base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k22), base.Simd_g_const(&F_UpsampleBgraLinePair_SSE2__k23))
				base.Simd_g_v128_store(m, v1386, v1398, v1471)
				if base.Ui32(v1387) < base.Ui32(int32(24)) {
					v1386 = v1386 + int32(32)
					v1387 = v1387 + int32(8)
					continue
				} else {
					break
				}
				break
			}
			v1480 = int32(2)
			v1485 = F_memcpy(m, l6+v727<<(uint(v1480)%32), v1372, v1369<<(uint(v1480)%32))
			mBase = m.M
		}
	}
	m.G0 = v35 + int32(528)
	return
}

var F_UpsampleBgraLinePair_SSE2__k0 = [2]uint64{0x101010101010101, 0x101010101010101}
var F_UpsampleBgraLinePair_SSE2__k1 = [2]uint64{0x1b0b1a0a19091808, 0x1f0f1e0e1d0d1c0c}
var F_UpsampleBgraLinePair_SSE2__k2 = [2]uint64{0x800b800a80098008, 0x800f800e800d800c}
var F_UpsampleBgraLinePair_SSE2__k3 = [2]uint64{0xb800a8009800880, 0xf800e800d800c80}
var F_UpsampleBgraLinePair_SSE2__k4 = [2]uint64{0x1303120211011000, 0x1707160615051404}
var F_UpsampleBgraLinePair_SSE2__k5 = [2]uint64{0x8003800280018000, 0x8007800680058004}
var F_UpsampleBgraLinePair_SSE2__k6 = [2]uint64{0x380028001800080, 0x780068005800480}
var F_UpsampleBgraLinePair_SSE2__k7 = [2]uint64{0x0, 0x0}
var F_UpsampleBgraLinePair_SSE2__k8 = [2]uint64{0x811a0000811a, 0x811a0000811a}
var F_UpsampleBgraLinePair_SSE2__k9 = [2]uint64{0xf0e0b0a07060302, 0x1f1e1b1a17161312}
var F_UpsampleBgraLinePair_SSE2__k10 = [2]uint64{0x4a8500004a85, 0x4a8500004a85}
var F_UpsampleBgraLinePair_SSE2__k11 = [2]uint64{0xf0e0b0a07060302, 0x8080808080808080}
var F_UpsampleBgraLinePair_SSE2__k12 = [2]uint64{0x8080808080808080, 0xf0e0b0a07060302}
var F_UpsampleBgraLinePair_SSE2__k13 = [2]uint64{0x662500006625, 0x662500006625}
var F_UpsampleBgraLinePair_SSE2__k14 = [2]uint64{0x4515451545154515, 0x4515451545154515}
var F_UpsampleBgraLinePair_SSE2__k15 = [2]uint64{0xc866c866c866c866, 0xc866c866c866c866}
var F_UpsampleBgraLinePair_SSE2__k16 = [2]uint64{0x191300001913, 0x191300001913}
var F_UpsampleBgraLinePair_SSE2__k17 = [2]uint64{0x340800003408, 0x340800003408}
var F_UpsampleBgraLinePair_SSE2__k18 = [2]uint64{0x2204220422042204, 0x2204220422042204}
var F_UpsampleBgraLinePair_SSE2__k19 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_UpsampleBgraLinePair_SSE2__k20 = [2]uint64{0x80800b0a80800908, 0x80800f0e80800d0c}
var F_UpsampleBgraLinePair_SSE2__k21 = [2]uint64{0xb0a808009088080, 0xf0e80800d0c8080}
var F_UpsampleBgraLinePair_SSE2__k22 = [2]uint64{0x8080030280800100, 0x8080070680800504}
var F_UpsampleBgraLinePair_SSE2__k23 = [2]uint64{0x302808001008080, 0x706808005048080}

//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_WebPCleanupTransparentArea(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v66 int32
	_ = v66
	var v67 int32
	_ = v67
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v140 int32
	_ = v140
	var v144 int32
	_ = v144
	var v145 int32
	_ = v145
	var v195 int32
	_ = v195
	var v200 int32
	_ = v200
	var v201 int32
	_ = v201
	var v205 int32
	_ = v205
	var v209 int32
	_ = v209
	var v213 int32
	_ = v213
	var v217 int32
	_ = v217
	var v221 int32
	_ = v221
	var v225 int32
	_ = v225
	var v229 int32
	_ = v229
	var v234 int32
	_ = v234
	var v235 int32
	_ = v235
	var v236 int32
	_ = v236
	var v240 int32
	_ = v240
	var v244 int32
	_ = v244
	var v248 int32
	_ = v248
	var v252 int32
	_ = v252
	var v256 int32
	_ = v256
	var v260 int32
	_ = v260
	var v264 int32
	_ = v264
	var v268 int32
	_ = v268
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
	var v297 int32
	_ = v297
	var v301 int32
	_ = v301
	var v302 int32
	_ = v302
	var v306 int32
	_ = v306
	var v310 int32
	_ = v310
	var v314 int32
	_ = v314
	var v318 int32
	_ = v318
	var v322 int32
	_ = v322
	var v326 int32
	_ = v326
	var v330 int32
	_ = v330
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v339 int32
	_ = v339
	var v343 int32
	_ = v343
	var v347 int32
	_ = v347
	var v351 int32
	_ = v351
	var v355 int32
	_ = v355
	var v359 int32
	_ = v359
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v380 int32
	_ = v380
	var v384 int32
	_ = v384
	var v388 int32
	_ = v388
	var v392 int32
	_ = v392
	var v396 int32
	_ = v396
	var v400 int32
	_ = v400
	var v401 int32
	_ = v401
	var v405 int32
	_ = v405
	var v409 int32
	_ = v409
	var v413 int32
	_ = v413
	var v417 int32
	_ = v417
	var v421 int32
	_ = v421
	var v425 int32
	_ = v425
	var v429 int32
	_ = v429
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v438 int32
	_ = v438
	var v442 int32
	_ = v442
	var v446 int32
	_ = v446
	var v450 int32
	_ = v450
	var v454 int32
	_ = v454
	var v458 int32
	_ = v458
	var v462 int32
	_ = v462
	var v466 int32
	_ = v466
	var v532 int32
	_ = v532
	var v533 int32
	_ = v533
	var v542 int32
	_ = v542
	var v545 int32
	_ = v545
	var v547 int32
	_ = v547
	var v550 int32
	_ = v550
	var v553 int32
	_ = v553
	var v556 int32
	_ = v556
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v565 int32
	_ = v565
	var v568 int32
	_ = v568
	var v569 int32
	_ = v569
	var v572 int32
	_ = v572
	var v573 int32
	_ = v573
	var v574 int32
	_ = v574
	var v575 int32
	_ = v575
	var v576 int32
	_ = v576
	var v578 int32
	_ = v578
	var v579 int32
	_ = v579
	var v585 int32
	_ = v585
	var v591 int32
	_ = v591
	var v615 int32
	_ = v615
	var v633 int32
	_ = v633
	var v634 int32
	_ = v634
	var v637 int32
	_ = v637
	var v638 int32
	_ = v638
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v648 int32
	_ = v648
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v655 int32
	_ = v655
	var v656 int32
	_ = v656
	var v657 int32
	_ = v657
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v661 int32
	_ = v661
	var v685 int32
	_ = v685
	var __phi685 int32
	_ = __phi685
	var v690 int32
	_ = v690
	var __phi690 int32
	_ = __phi690
	var v696 int32
	_ = v696
	var __phi696 int32
	_ = __phi696
	var v722 int32
	_ = v722
	var __phi722 int32
	_ = __phi722
	var v724 int32
	_ = v724
	var __phi724 int32
	_ = __phi724
	var v725 int32
	_ = v725
	var __phi725 int32
	_ = __phi725
	var v742 int32
	_ = v742
	var v743 int32
	_ = v743
	var v744 int32
	_ = v744
	var v746 int32
	_ = v746
	var v752 int32
	_ = v752
	var v753 int32
	_ = v753
	var v755 int32
	_ = v755
	var v756 int32
	_ = v756
	var v804 int32
	_ = v804
	var v809 int32
	_ = v809
	var v811 int32
	_ = v811
	var v812 int32
	_ = v812
	var v813 int32
	_ = v813
	var v818 int32
	_ = v818
	var v820 int32
	_ = v820
	var v821 int32
	_ = v821
	var v822 int32
	_ = v822
	var v827 int32
	_ = v827
	var v829 int32
	_ = v829
	var v830 int32
	_ = v830
	var v831 int32
	_ = v831
	var v836 int32
	_ = v836
	var v838 int32
	_ = v838
	var v839 int32
	_ = v839
	var v840 int32
	_ = v840
	var v845 int32
	_ = v845
	var v847 int32
	_ = v847
	var v848 int32
	_ = v848
	var v849 int32
	_ = v849
	var v854 int32
	_ = v854
	var v856 int32
	_ = v856
	var v857 int32
	_ = v857
	var v858 int32
	_ = v858
	var v863 int32
	_ = v863
	var v865 int32
	_ = v865
	var v866 int32
	_ = v866
	var v867 int32
	_ = v867
	var v872 int32
	_ = v872
	var v874 int32
	_ = v874
	var v875 int32
	_ = v875
	var v879 int32
	_ = v879
	var v884 int32
	_ = v884
	var v890 int32
	_ = v890
	var v894 int32
	_ = v894
	var v896 int32
	_ = v896
	var v944 int32
	_ = v944
	var v946 int32
	_ = v946
	var v948 int32
	_ = v948
	var v950 int32
	_ = v950
	var v952 int32
	_ = v952
	var v954 int32
	_ = v954
	var v956 int32
	_ = v956
	var v958 int32
	_ = v958
	var v963 int32
	_ = v963
	var v971 int32
	_ = v971
	var v973 int32
	_ = v973
	var v975 int32
	_ = v975
	var v976 int32
	_ = v976
	var v977 int32
	_ = v977
	var v978 int32
	_ = v978
	var v979 int32
	_ = v979
	var v980 int32
	_ = v980
	var v985 int64
	_ = v985
	var v987 int32
	_ = v987
	var v989 int32
	_ = v989
	var v991 int32
	_ = v991
	var v993 int32
	_ = v993
	var v995 int32
	_ = v995
	var v997 int32
	_ = v997
	var v1001 int32
	_ = v1001
	var v1002 int32
	_ = v1002
	var v1004 int32
	_ = v1004
	var v1005 int32
	_ = v1005
	var v1007 int32
	_ = v1007
	var v1009 int32
	_ = v1009
	var v1013 int32
	_ = v1013
	var v1017 int32
	_ = v1017
	var v1019 int32
	_ = v1019
	var v1021 int32
	_ = v1021
	var v1039 int32
	_ = v1039
	var v1065 int32
	_ = v1065
	var v1067 int32
	_ = v1067
	var v1068 int32
	_ = v1068
	var v1085 int32
	_ = v1085
	var v1089 int32
	_ = v1089
	var v1126 int32
	_ = v1126
	var v1128 int32
	_ = v1128
	var v1129 int32
	_ = v1129
	var v1146 int32
	_ = v1146
	var v1149 int32
	_ = v1149
	var v1150 int32
	_ = v1150
	var v1153 int32
	_ = v1153
	var v1154 int32
	_ = v1154
	var v1158 int32
	_ = v1158
	var v1167 int32
	_ = v1167
	var v1168 int32
	_ = v1168
	var v1171 int32
	_ = v1171
	var v1219 int32
	_ = v1219
	var v1220 int32
	_ = v1220
	var v1226 int32
	_ = v1226
	var v1228 int32
	_ = v1228
	var v1229 int32
	_ = v1229
	var v1232 int32
	_ = v1232
	var v1235 int32
	_ = v1235
	var v1240 int32
	_ = v1240
	var v1242 int32
	_ = v1242
	var v1243 int32
	_ = v1243
	var v1245 int32
	_ = v1245
	var v1253 int32
	_ = v1253
	var v1254 int32
	_ = v1254
	var v1257 int32
	_ = v1257
	var v1308 int32
	_ = v1308
	var v1314 int32
	_ = v1314
	var v1316 int32
	_ = v1316
	var v1317 int32
	_ = v1317
	var v1319 int32
	_ = v1319
	var v1329 int32
	_ = v1329
	var v1330 int32
	_ = v1330
	var v1333 int32
	_ = v1333
	var v1381 int32
	_ = v1381
	var v1382 int32
	_ = v1382
	var v1388 int32
	_ = v1388
	var v1390 int32
	_ = v1390
	var v1391 int32
	_ = v1391
	var v1394 int32
	_ = v1394
	var v1397 int32
	_ = v1397
	var v1402 int32
	_ = v1402
	var v1404 int32
	_ = v1404
	var v1405 int32
	_ = v1405
	var v1407 int32
	_ = v1407
	var v1415 int32
	_ = v1415
	var v1416 int32
	_ = v1416
	var v1419 int32
	_ = v1419
	var v1467 int32
	_ = v1467
	var v1468 int32
	_ = v1468
	var v1472 int32
	_ = v1472
	var v1478 int32
	_ = v1478
	var v1480 int32
	_ = v1480
	var v1481 int32
	_ = v1481
	var v1483 int32
	_ = v1483
	var v1493 int32
	_ = v1493
	var v1494 int32
	_ = v1494
	var v1497 int32
	_ = v1497
	var v1545 int32
	_ = v1545
	var v1546 int32
	_ = v1546
	var v1552 int32
	_ = v1552
	var v1554 int32
	_ = v1554
	var v1555 int32
	_ = v1555
	var v1558 int32
	_ = v1558
	var v1561 int32
	_ = v1561
	var v1566 int32
	_ = v1566
	var v1568 int32
	_ = v1568
	var v1569 int32
	_ = v1569
	var v1571 int32
	_ = v1571
	var v1579 int32
	_ = v1579
	var v1580 int32
	_ = v1580
	var v1583 int32
	_ = v1583
	var v1631 int32
	_ = v1631
	var v1632 int32
	_ = v1632
	var v1636 int32
	_ = v1636
	var v1642 int32
	_ = v1642
	var v1644 int32
	_ = v1644
	var v1645 int32
	_ = v1645
	var v1647 int32
	_ = v1647
	var v1657 int32
	_ = v1657
	var v1658 int32
	_ = v1658
	var v1661 int32
	_ = v1661
	var v1709 int32
	_ = v1709
	var v1710 int32
	_ = v1710
	var v1716 int32
	_ = v1716
	var v1718 int32
	_ = v1718
	var v1719 int32
	_ = v1719
	var v1722 int32
	_ = v1722
	var v1725 int32
	_ = v1725
	var v1730 int32
	_ = v1730
	var v1732 int32
	_ = v1732
	var v1733 int32
	_ = v1733
	var v1735 int32
	_ = v1735
	var v1743 int32
	_ = v1743
	var v1744 int32
	_ = v1744
	var v1747 int32
	_ = v1747
	var v1795 int32
	_ = v1795
	var v1796 int32
	_ = v1796
	var v1800 int32
	_ = v1800
	var v1806 int32
	_ = v1806
	var v1808 int32
	_ = v1808
	var v1809 int32
	_ = v1809
	var v1811 int32
	_ = v1811
	var v1821 int32
	_ = v1821
	var v1822 int32
	_ = v1822
	var v1825 int32
	_ = v1825
	var v1873 int32
	_ = v1873
	var v1874 int32
	_ = v1874
	var v1880 int32
	_ = v1880
	var v1882 int32
	_ = v1882
	var v1883 int32
	_ = v1883
	var v1886 int32
	_ = v1886
	var v1889 int32
	_ = v1889
	var v1894 int32
	_ = v1894
	var v1896 int32
	_ = v1896
	var v1897 int32
	_ = v1897
	var v1899 int32
	_ = v1899
	var v1907 int32
	_ = v1907
	var v1908 int32
	_ = v1908
	var v1911 int32
	_ = v1911
	var v1959 int32
	_ = v1959
	var v1960 int32
	_ = v1960
	var v1964 int32
	_ = v1964
	var v1970 int32
	_ = v1970
	var v1972 int32
	_ = v1972
	var v1973 int32
	_ = v1973
	var v1975 int32
	_ = v1975
	var v1985 int32
	_ = v1985
	var v1986 int32
	_ = v1986
	var v1989 int32
	_ = v1989
	var v2037 int32
	_ = v2037
	var v2038 int32
	_ = v2038
	var v2044 int32
	_ = v2044
	var v2046 int32
	_ = v2046
	var v2047 int32
	_ = v2047
	var v2050 int32
	_ = v2050
	var v2053 int32
	_ = v2053
	var v2058 int32
	_ = v2058
	var v2060 int32
	_ = v2060
	var v2061 int32
	_ = v2061
	var v2063 int32
	_ = v2063
	var v2071 int32
	_ = v2071
	var v2072 int32
	_ = v2072
	var v2075 int32
	_ = v2075
	var v2123 int32
	_ = v2123
	var v2124 int32
	_ = v2124
	var v2128 int32
	_ = v2128
	var v2134 int32
	_ = v2134
	var v2136 int32
	_ = v2136
	var v2137 int32
	_ = v2137
	var v2139 int32
	_ = v2139
	var v2149 int32
	_ = v2149
	var v2150 int32
	_ = v2150
	var v2153 int32
	_ = v2153
	var v2201 int32
	_ = v2201
	var v2202 int32
	_ = v2202
	var v2208 int32
	_ = v2208
	var v2210 int32
	_ = v2210
	var v2211 int32
	_ = v2211
	var v2214 int32
	_ = v2214
	var v2217 int32
	_ = v2217
	var v2222 int32
	_ = v2222
	var v2224 int32
	_ = v2224
	var v2225 int32
	_ = v2225
	var v2227 int32
	_ = v2227
	var v2235 int32
	_ = v2235
	var v2236 int32
	_ = v2236
	var v2239 int32
	_ = v2239
	var v2287 int32
	_ = v2287
	var v2288 int32
	_ = v2288
	var v2292 int32
	_ = v2292
	var v2298 int32
	_ = v2298
	var v2300 int32
	_ = v2300
	var v2301 int32
	_ = v2301
	var v2303 int32
	_ = v2303
	var v2313 int32
	_ = v2313
	var v2314 int32
	_ = v2314
	var v2317 int32
	_ = v2317
	var v2365 int32
	_ = v2365
	var v2366 int32
	_ = v2366
	var v2372 int32
	_ = v2372
	var v2374 int32
	_ = v2374
	var v2375 int32
	_ = v2375
	var v2378 int32
	_ = v2378
	var v2381 int32
	_ = v2381
	var v2386 int32
	_ = v2386
	var v2388 int32
	_ = v2388
	var v2389 int32
	_ = v2389
	var v2391 int32
	_ = v2391
	var v2399 int32
	_ = v2399
	var v2400 int32
	_ = v2400
	var v2403 int32
	_ = v2403
	var v2451 int32
	_ = v2451
	var v2452 int32
	_ = v2452
	var v2456 int32
	_ = v2456
	var v2462 int32
	_ = v2462
	var v2464 int32
	_ = v2464
	var v2465 int32
	_ = v2465
	var v2471 int32
	_ = v2471
	var v2485 int32
	_ = v2485
	var v2533 int32
	_ = v2533
	var v2534 int32
	_ = v2534
	var v2539 int32
	_ = v2539
	var v2545 int32
	_ = v2545
	var v2557 int32
	_ = v2557
	var v2608 int32
	_ = v2608
	var v2612 int32
	_ = v2612
	var v2626 int32
	_ = v2626
	var v2674 int32
	_ = v2674
	var v2675 int32
	_ = v2675
	var v2680 int32
	_ = v2680
	var v2686 int32
	_ = v2686
	var v2698 int32
	_ = v2698
	var v2749 int32
	_ = v2749
	var v2753 int32
	_ = v2753
	var v2767 int32
	_ = v2767
	var v2815 int32
	_ = v2815
	var v2816 int32
	_ = v2816
	var v2821 int32
	_ = v2821
	var v2827 int32
	_ = v2827
	var v2839 int32
	_ = v2839
	var v2890 int32
	_ = v2890
	var v2894 int32
	_ = v2894
	var v2908 int32
	_ = v2908
	var v2956 int32
	_ = v2956
	var v2957 int32
	_ = v2957
	var v2962 int32
	_ = v2962
	var v2968 int32
	_ = v2968
	var v2980 int32
	_ = v2980
	var v3031 int32
	_ = v3031
	var v3035 int32
	_ = v3035
	var v3049 int32
	_ = v3049
	var v3097 int32
	_ = v3097
	var v3098 int32
	_ = v3098
	var v3103 int32
	_ = v3103
	var v3109 int32
	_ = v3109
	var v3121 int32
	_ = v3121
	var v3172 int32
	_ = v3172
	var v3176 int32
	_ = v3176
	var v3190 int32
	_ = v3190
	var v3238 int32
	_ = v3238
	var v3239 int32
	_ = v3239
	var v3244 int32
	_ = v3244
	var v3250 int32
	_ = v3250
	var v3262 int32
	_ = v3262
	var v3313 int32
	_ = v3313
	var v3317 int32
	_ = v3317
	var v3331 int32
	_ = v3331
	var v3379 int32
	_ = v3379
	var v3380 int32
	_ = v3380
	var v3385 int32
	_ = v3385
	var v3391 int32
	_ = v3391
	var v3403 int32
	_ = v3403
	var v3454 int32
	_ = v3454
	var v3458 int32
	_ = v3458
	var v3472 int32
	_ = v3472
	var v3520 int32
	_ = v3520
	var v3521 int32
	_ = v3521
	var v3526 int32
	_ = v3526
	var v3532 int32
	_ = v3532
	var v3544 int32
	_ = v3544
	var v3595 int32
	_ = v3595
	var v3672 int32
	_ = v3672
	var v3673 int32
	_ = v3673
	var v3675 int32
	_ = v3675
	var v3691 int32
	_ = v3691
	var v3692 int32
	_ = v3692
	var v3697 int32
	_ = v3697
	var v3736 int32
	_ = v3736
	var v3753 int32
	_ = v3753
	var v3754 int32
	_ = v3754
	var v3804 int32
	_ = v3804
	var v3808 int32
	_ = v3808
	var v3810 int32
	_ = v3810
	var v3813 int32
	_ = v3813
	var v3815 int32
	_ = v3815
	var v3816 int32
	_ = v3816
	var v3864 int32
	_ = v3864
	var v3865 int32
	_ = v3865
	var v3866 int32
	_ = v3866
	var v3871 int32
	_ = v3871
	var v3873 int32
	_ = v3873
	var v3874 int32
	_ = v3874
	var v3877 int32
	_ = v3877
	var v3880 int32
	_ = v3880
	var v3884 int32
	_ = v3884
	var v3886 int32
	_ = v3886
	var v3887 int32
	_ = v3887
	var v3890 int32
	_ = v3890
	var v3897 int32
	_ = v3897
	var v3899 int32
	_ = v3899
	var v3900 int32
	_ = v3900
	var v3903 int32
	_ = v3903
	var v3910 int32
	_ = v3910
	var v3912 int32
	_ = v3912
	var v3913 int32
	_ = v3913
	var v3916 int32
	_ = v3916
	var v3923 int32
	_ = v3923
	var v3925 int32
	_ = v3925
	var v3926 int32
	_ = v3926
	var v3929 int32
	_ = v3929
	var v3936 int32
	_ = v3936
	var v3938 int32
	_ = v3938
	var v3939 int32
	_ = v3939
	var v3942 int32
	_ = v3942
	var v3949 int32
	_ = v3949
	var v3951 int32
	_ = v3951
	var v3952 int32
	_ = v3952
	var v3955 int32
	_ = v3955
	var v3962 int32
	_ = v3962
	var v3964 int32
	_ = v3964
	var v3965 int32
	_ = v3965
	var v3969 int32
	_ = v3969
	var v3973 int32
	_ = v3973
	var v3974 int32
	_ = v3974
	var v3978 int32
	_ = v3978
	var v3983 int32
	_ = v3983
	var v4032 int32
	_ = v4032
	var v4033 int32
	_ = v4033
	var v4034 int32
	_ = v4034
	var v4038 int32
	_ = v4038
	var v4044 int32
	_ = v4044
	var v4050 int32
	_ = v4050
	var v4056 int32
	_ = v4056
	var v4062 int32
	_ = v4062
	var v4068 int32
	_ = v4068
	var v4074 int32
	_ = v4074
	var v4081 int32
	_ = v4081
	var v4140 int32
	_ = v4140
	var v4143 int32
	_ = v4143
	var v4163 int32
	_ = v4163
	var v4208 int32
	_ = v4208
	var v4218 int32
	_ = v4218
	var v4219 int32
	_ = v4219
	var v4220 int32
	_ = v4220
	var v4221 int32
	_ = v4221
	var v4228 int32
	_ = v4228
	var v4230 int32
	_ = v4230
	var v4232 int32
	_ = v4232
	var v4234 int32
	_ = v4234
	var v4238 int32
	_ = v4238
	var v4290 int32
	_ = v4290
	var v4291 int32
	_ = v4291
	var v4294 int32
	_ = v4294
	var v4342 int32
	_ = v4342
	var v4343 int32
	_ = v4343
	var v4349 int32
	_ = v4349
	var v4351 int32
	_ = v4351
	var v4352 int32
	_ = v4352
	var v4355 int32
	_ = v4355
	var v4358 int32
	_ = v4358
	var v4363 int32
	_ = v4363
	var v4365 int32
	_ = v4365
	var v4366 int32
	_ = v4366
	var v4368 int32
	_ = v4368
	var v4376 int32
	_ = v4376
	var v4377 int32
	_ = v4377
	var v4380 int32
	_ = v4380
	var v4431 int32
	_ = v4431
	var v4437 int32
	_ = v4437
	var v4439 int32
	_ = v4439
	var v4440 int32
	_ = v4440
	var v4444 int32
	_ = v4444
	var v4452 int32
	_ = v4452
	var v4459 int32
	_ = v4459
	var v4465 int32
	_ = v4465
	var v4469 int32
	_ = v4469
	var v4515 int32
	_ = v4515
	var v4526 int32
	_ = v4526
	var v4577 int32
	_ = v4577
	var v4578 int32
	_ = v4578
	var v4583 int32
	_ = v4583
	var v4589 int32
	_ = v4589
	var v4598 int32
	_ = v4598
	var v4652 int32
	_ = v4652
	var v4658 int32
	_ = v4658
	if l0 == int32(0) {
	} else {
		v61 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
		v62 = int32(8)
		v63 = base.I32_div_s(v61, v62)
		v64 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		v66 = base.I32_div_s(v64, v62)
		v67 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v67 == int32(0) {
			v547 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
			if v547 == int32(0) {
			} else {
				v550 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				if v550 == int32(0) {
				} else {
					v553 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
					if v553 == int32(0) {
					} else {
						v556 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
						if v556 == int32(0) {
						} else {
							v559 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
							v560 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
							if v61 < int32(8) {
								v3691 = v547
								v3692 = v550
								v3697 = int32(0)
							} else {
								v565 = v64 & int32(-2)
								v568 = int32(1)
								v569 = v64 & v568
								v572 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
								v573 = int32(2)
								v574 = v572 << (uint(v573) % 32)
								v575 = int32(3)
								v576 = v560 << (uint(v575) % 32)
								v578 = v559 << (uint(v575) % 32)
								v579 = int32(7)
								v585 = int32(6)
								v591 = int32(5)
								v615 = int32(0)
								v633 = v547
								v634 = v550
								v637 = v553
								v638 = v556
								v641 = v550 + v560
								v642 = v547 + v559
								v646 = v550 + v560*v579
								v647 = v547 + v559*v579
								v648 = v550 + v560*v585
								v649 = v547 + v559*v585
								v650 = v550 + v560*v591
								v651 = v547 + v559*v591
								v652 = v550 + v560<<(uint(v573)%32)
								v653 = v547 + v559<<(uint(v573)%32)
								v654 = v550 + v560*v575
								v655 = v547 + v559*v575
								v656 = v550 + v560<<(uint(v568)%32)
								v657 = v547 + v559<<(uint(v568)%32)
								v658 = v615
								v659 = int32(8)
								v660 = v615
								v661 = v615
								for {
									if int32(8) <= v64 {
										__phi685 = int32(0)
										__phi690 = int32(8)
										__phi696 = int32(1)
										__phi722 = v658
										__phi724 = v660
										__phi725 = v661
										v685 = __phi685
										v690 = __phi690
										v696 = __phi696
										v722 = __phi722
										v724 = __phi724
										v725 = __phi725
										for {
											v742 = int32(0)
											v743 = v634 + v685
											v744 = v633 + v685
											v746 = v743
											v752 = v744
											v753 = v742
											v755 = int32(8)
											v756 = v742
											for {
												v804 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752))))
												if v804 == int32(0) {
													v811 = v753
													v812 = v756
												} else {
													v809 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746))))
													v811 = v753 + int32(1)
													v812 = v756 + v809
												}
												v813 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+1)))
												if v813 == int32(0) {
													v820 = v811
													v821 = v812
												} else {
													v818 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+1)))
													v820 = v811 + int32(1)
													v821 = v812 + v818
												}
												v822 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+2)))
												if v822 == int32(0) {
													v829 = v820
													v830 = v821
												} else {
													v827 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+2)))
													v829 = v820 + int32(1)
													v830 = v821 + v827
												}
												v831 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+3)))
												if v831 == int32(0) {
													v838 = v829
													v839 = v830
												} else {
													v836 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+3)))
													v838 = v829 + int32(1)
													v839 = v830 + v836
												}
												v840 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+4)))
												if v840 == int32(0) {
													v847 = v838
													v848 = v839
												} else {
													v845 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+4)))
													v847 = v838 + int32(1)
													v848 = v839 + v845
												}
												v849 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+5)))
												if v849 == int32(0) {
													v856 = v847
													v857 = v848
												} else {
													v854 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+5)))
													v856 = v847 + int32(1)
													v857 = v848 + v854
												}
												v858 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+6)))
												if v858 == int32(0) {
													v865 = v856
													v866 = v857
												} else {
													v863 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+6)))
													v865 = v856 + int32(1)
													v866 = v857 + v863
												}
												v867 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v752)+7)))
												if v867 == int32(0) {
													v874 = v865
													v875 = v866
												} else {
													v872 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v746)+7)))
													v874 = v865 + int32(1)
													v875 = v866 + v872
												}
												v879 = v755 + int32(-1)
												if v879 != 0 {
													v746 = v746 + v560
													v752 = v752 + v559
													v753 = v874
													v755 = v879
													v756 = v875
													continue
												} else {
													break
												}
												break
											}
											if base.Ui32(int32(62)) < base.Ui32(v874+int32(-1)) {
												if v874 == int32(0) {
													if v696 != 0 {
														v971 = int32(base.Ui32(v685) >> (uint(int32(1)) % 32))
														v973 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v638+v971))))
														v975 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v637+v971))))
														v976 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v743))))
														v977 = v971
														v978 = v973
														v979 = v975
														v980 = v976
													} else {
														v977 = int32(base.Ui32(v685) >> (uint(int32(1)) % 32))
														v978 = v722
														v979 = v724
														v980 = v725
													}
													v985 = base.I64_extend_i32_u(v980) & int64(255) * int64(72340172838076673)
													*(*int64)(unsafe.Add(mBase, uint32(v743))) = v985
													v987 = v743 + v560
													*(*int64)(unsafe.Add(mBase, uint32(v987))) = v985
													v989 = v987 + v560
													*(*int64)(unsafe.Add(mBase, uint32(v989))) = v985
													v991 = v989 + v560
													*(*int64)(unsafe.Add(mBase, uint32(v991))) = v985
													v993 = v991 + v560
													*(*int64)(unsafe.Add(mBase, uint32(v993))) = v985
													v995 = v993 + v560
													*(*int64)(unsafe.Add(mBase, uint32(v995))) = v985
													v997 = v995 + v560
													*(*int64)(unsafe.Add(mBase, uint32(v997))) = v985
													*(*int64)(unsafe.Add(mBase, uint32(v997+v560))) = v985
													v1001 = v637 + v977
													v1002 = int32(255)
													v1004 = int32(16843009)
													v1005 = v979 & v1002 * v1004
													*(*int32)(unsafe.Add(mBase, uint32(v1001))) = v1005
													v1007 = v1001 + v572
													*(*int32)(unsafe.Add(mBase, uint32(v1007))) = v1005
													v1009 = v1007 + v572
													*(*int32)(unsafe.Add(mBase, uint32(v1009))) = v1005
													*(*int32)(unsafe.Add(mBase, uint32(v1009+v572))) = v1005
													v1013 = v638 + v977
													v1017 = v978 & v1002 * v1004
													*(*int32)(unsafe.Add(mBase, uint32(v1013))) = v1017
													v1019 = v1013 + v572
													*(*int32)(unsafe.Add(mBase, uint32(v1019))) = v1017
													v1021 = v1019 + v572
													*(*int32)(unsafe.Add(mBase, uint32(v1021))) = v1017
													*(*int32)(unsafe.Add(mBase, uint32(v1021+v572))) = v1017
													v1039 = int32(0)
													v1065 = v978
													v1067 = v979
													v1068 = v980
												} else {
													v1039 = int32(1)
													v1065 = v722
													v1067 = v724
													v1068 = v725
												}
											} else {
												v884 = base.I32_div_s(v875, v874)
												v890 = v743
												v894 = v744
												v896 = int32(8)
												for {
													v944 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894))))
													if v944 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890))) = uint8(v884)
													}
													v946 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+1)))
													if v946 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+1)) = uint8(v884)
													}
													v948 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+2)))
													if v948 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+2)) = uint8(v884)
													}
													v950 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+3)))
													if v950 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+3)) = uint8(v884)
													}
													v952 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+4)))
													if v952 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+4)) = uint8(v884)
													}
													v954 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+5)))
													if v954 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+5)) = uint8(v884)
													}
													v956 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+6)))
													if v956 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+6)) = uint8(v884)
													}
													v958 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v894)+7)))
													if v958 != 0 {
													} else {
														*(*uint8)(unsafe.Add(mBase, uint32(v890)+7)) = uint8(v884)
													}
													v963 = v896 + int32(-1)
													if v963 != 0 {
														v890 = v890 + v560
														v894 = v894 + v559
														v896 = v963
														continue
													} else {
														break
													}
													break
												}
												v1039 = int32(1)
												v1065 = v722
												v1067 = v724
												v1068 = v725
											}
											v1085 = v690 + int32(8)
											if v1085 <= v64 {
												__phi685 = v690
												__phi690 = v1085
												__phi696 = v1039
												__phi722 = v1065
												__phi724 = v1067
												__phi725 = v1068
												v685 = __phi685
												v690 = __phi690
												v696 = __phi696
												v722 = __phi722
												v724 = __phi724
												v725 = __phi725
												continue
											} else {
												break
											}
											break
										}
										v1089 = v690
										v1126 = v1065
										v1128 = v1067
										v1129 = v1068
									} else {
										v1089 = int32(0)
										v1126 = v658
										v1128 = v660
										v1129 = v661
									}
									if v64 <= v1089 {
									} else {
										v1146 = v64 - v1089
										if v1146 < int32(1) {
										} else {
											v1149 = v634 + v1089
											v1150 = v633 + v1089
											v1153 = v64 + (v1089 ^ int32(-1))
											if v1153 != 0 {
												v1158 = int32(0)
												v1167 = v1158
												v1168 = v1158
												v1171 = v1158
												for {
													v1219 = v1150 + v1167
													v1220 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1219))))
													if v1220 == int32(0) {
														v1228 = v1168
														v1229 = v1171
													} else {
														v1226 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1149+v1167))))
														v1228 = v1168 + int32(1)
														v1229 = v1171 + v1226
													}
													v1232 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1219+int32(1)))))
													if v1232 == int32(0) {
														v1242 = v1228
														v1243 = v1229
													} else {
														v1235 = int32(1)
														v1240 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1149+v1167+v1235))))
														v1242 = v1228 + v1235
														v1243 = v1229 + v1240
													}
													v1245 = v1167 + int32(2)
													if v565-v1089 != v1245 {
														v1167 = v1245
														v1168 = v1242
														v1171 = v1243
														continue
													} else {
														break
													}
													break
												}
												v1253 = v1245
												v1254 = v1242
												v1257 = v1243
											} else {
												v1154 = int32(0)
												v1253 = v1154
												v1254 = v1154
												v1257 = v1154
											}
											if v569 == int32(0) {
												v1316 = v1254
												v1317 = v1257
											} else {
												v1308 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1150+v1253))))
												if v1308 == int32(0) {
													v1316 = v1254
													v1317 = v1257
												} else {
													v1314 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1149+v1253))))
													v1316 = v1254 + int32(1)
													v1317 = v1257 + v1314
												}
											}
											if v1153 != 0 {
												v1319 = v641 + v1089
												v1329 = int32(0)
												v1330 = v1316
												v1333 = v1317
												for {
													v1381 = v642 + v1089 + v1329
													v1382 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1381))))
													if v1382 == int32(0) {
														v1390 = v1330
														v1391 = v1333
													} else {
														v1388 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1319+v1329))))
														v1390 = v1330 + int32(1)
														v1391 = v1333 + v1388
													}
													v1394 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1381+int32(1)))))
													if v1394 == int32(0) {
														v1404 = v1390
														v1405 = v1391
													} else {
														v1397 = int32(1)
														v1402 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1319+v1329+v1397))))
														v1404 = v1390 + v1397
														v1405 = v1391 + v1402
													}
													v1407 = v1329 + int32(2)
													if v565-v1089 != v1407 {
														v1329 = v1407
														v1330 = v1404
														v1333 = v1405
														continue
													} else {
														break
													}
													break
												}
												v1415 = v1407
												v1416 = v1404
												v1419 = v1405
											} else {
												v1415 = int32(0)
												v1416 = v1316
												v1419 = v1317
											}
											v1467 = v1149 + v560
											v1468 = v1150 + v559
											if v569 == int32(0) {
												v1480 = v1416
												v1481 = v1419
											} else {
												v1472 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1468+v1415))))
												if v1472 == int32(0) {
													v1480 = v1416
													v1481 = v1419
												} else {
													v1478 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1467+v1415))))
													v1480 = v1416 + int32(1)
													v1481 = v1419 + v1478
												}
											}
											if v1153 != 0 {
												v1483 = v656 + v1089
												v1493 = int32(0)
												v1494 = v1480
												v1497 = v1481
												for {
													v1545 = v657 + v1089 + v1493
													v1546 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1545))))
													if v1546 == int32(0) {
														v1554 = v1494
														v1555 = v1497
													} else {
														v1552 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1483+v1493))))
														v1554 = v1494 + int32(1)
														v1555 = v1497 + v1552
													}
													v1558 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1545+int32(1)))))
													if v1558 == int32(0) {
														v1568 = v1554
														v1569 = v1555
													} else {
														v1561 = int32(1)
														v1566 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1483+v1493+v1561))))
														v1568 = v1554 + v1561
														v1569 = v1555 + v1566
													}
													v1571 = v1493 + int32(2)
													if v565-v1089 != v1571 {
														v1493 = v1571
														v1494 = v1568
														v1497 = v1569
														continue
													} else {
														break
													}
													break
												}
												v1579 = v1571
												v1580 = v1568
												v1583 = v1569
											} else {
												v1579 = int32(0)
												v1580 = v1480
												v1583 = v1481
											}
											v1631 = v1467 + v560
											v1632 = v1468 + v559
											if v569 == int32(0) {
												v1644 = v1580
												v1645 = v1583
											} else {
												v1636 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1632+v1579))))
												if v1636 == int32(0) {
													v1644 = v1580
													v1645 = v1583
												} else {
													v1642 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1631+v1579))))
													v1644 = v1580 + int32(1)
													v1645 = v1583 + v1642
												}
											}
											if v1153 != 0 {
												v1647 = v654 + v1089
												v1657 = int32(0)
												v1658 = v1644
												v1661 = v1645
												for {
													v1709 = v655 + v1089 + v1657
													v1710 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1709))))
													if v1710 == int32(0) {
														v1718 = v1658
														v1719 = v1661
													} else {
														v1716 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1647+v1657))))
														v1718 = v1658 + int32(1)
														v1719 = v1661 + v1716
													}
													v1722 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1709+int32(1)))))
													if v1722 == int32(0) {
														v1732 = v1718
														v1733 = v1719
													} else {
														v1725 = int32(1)
														v1730 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1647+v1657+v1725))))
														v1732 = v1718 + v1725
														v1733 = v1719 + v1730
													}
													v1735 = v1657 + int32(2)
													if v565-v1089 != v1735 {
														v1657 = v1735
														v1658 = v1732
														v1661 = v1733
														continue
													} else {
														break
													}
													break
												}
												v1743 = v1735
												v1744 = v1732
												v1747 = v1733
											} else {
												v1743 = int32(0)
												v1744 = v1644
												v1747 = v1645
											}
											v1795 = v1631 + v560
											v1796 = v1632 + v559
											if v569 == int32(0) {
												v1808 = v1744
												v1809 = v1747
											} else {
												v1800 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1796+v1743))))
												if v1800 == int32(0) {
													v1808 = v1744
													v1809 = v1747
												} else {
													v1806 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1795+v1743))))
													v1808 = v1744 + int32(1)
													v1809 = v1747 + v1806
												}
											}
											if v1153 != 0 {
												v1811 = v652 + v1089
												v1821 = int32(0)
												v1822 = v1808
												v1825 = v1809
												for {
													v1873 = v653 + v1089 + v1821
													v1874 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1873))))
													if v1874 == int32(0) {
														v1882 = v1822
														v1883 = v1825
													} else {
														v1880 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1811+v1821))))
														v1882 = v1822 + int32(1)
														v1883 = v1825 + v1880
													}
													v1886 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1873+int32(1)))))
													if v1886 == int32(0) {
														v1896 = v1882
														v1897 = v1883
													} else {
														v1889 = int32(1)
														v1894 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1811+v1821+v1889))))
														v1896 = v1882 + v1889
														v1897 = v1883 + v1894
													}
													v1899 = v1821 + int32(2)
													if v565-v1089 != v1899 {
														v1821 = v1899
														v1822 = v1896
														v1825 = v1897
														continue
													} else {
														break
													}
													break
												}
												v1907 = v1899
												v1908 = v1896
												v1911 = v1897
											} else {
												v1907 = int32(0)
												v1908 = v1808
												v1911 = v1809
											}
											v1959 = v1795 + v560
											v1960 = v1796 + v559
											if v569 == int32(0) {
												v1972 = v1908
												v1973 = v1911
											} else {
												v1964 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1960+v1907))))
												if v1964 == int32(0) {
													v1972 = v1908
													v1973 = v1911
												} else {
													v1970 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1959+v1907))))
													v1972 = v1908 + int32(1)
													v1973 = v1911 + v1970
												}
											}
											if v1153 != 0 {
												v1975 = v650 + v1089
												v1985 = int32(0)
												v1986 = v1972
												v1989 = v1973
												for {
													v2037 = v651 + v1089 + v1985
													v2038 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2037))))
													if v2038 == int32(0) {
														v2046 = v1986
														v2047 = v1989
													} else {
														v2044 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1975+v1985))))
														v2046 = v1986 + int32(1)
														v2047 = v1989 + v2044
													}
													v2050 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2037+int32(1)))))
													if v2050 == int32(0) {
														v2060 = v2046
														v2061 = v2047
													} else {
														v2053 = int32(1)
														v2058 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1975+v1985+v2053))))
														v2060 = v2046 + v2053
														v2061 = v2047 + v2058
													}
													v2063 = v1985 + int32(2)
													if v565-v1089 != v2063 {
														v1985 = v2063
														v1986 = v2060
														v1989 = v2061
														continue
													} else {
														break
													}
													break
												}
												v2071 = v2063
												v2072 = v2060
												v2075 = v2061
											} else {
												v2071 = int32(0)
												v2072 = v1972
												v2075 = v1973
											}
											v2123 = v1959 + v560
											v2124 = v1960 + v559
											if v569 == int32(0) {
												v2136 = v2072
												v2137 = v2075
											} else {
												v2128 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2124+v2071))))
												if v2128 == int32(0) {
													v2136 = v2072
													v2137 = v2075
												} else {
													v2134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2123+v2071))))
													v2136 = v2072 + int32(1)
													v2137 = v2075 + v2134
												}
											}
											if v1153 != 0 {
												v2139 = v648 + v1089
												v2149 = int32(0)
												v2150 = v2136
												v2153 = v2137
												for {
													v2201 = v649 + v1089 + v2149
													v2202 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2201))))
													if v2202 == int32(0) {
														v2210 = v2150
														v2211 = v2153
													} else {
														v2208 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2139+v2149))))
														v2210 = v2150 + int32(1)
														v2211 = v2153 + v2208
													}
													v2214 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2201+int32(1)))))
													if v2214 == int32(0) {
														v2224 = v2210
														v2225 = v2211
													} else {
														v2217 = int32(1)
														v2222 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2139+v2149+v2217))))
														v2224 = v2210 + v2217
														v2225 = v2211 + v2222
													}
													v2227 = v2149 + int32(2)
													if v565-v1089 != v2227 {
														v2149 = v2227
														v2150 = v2224
														v2153 = v2225
														continue
													} else {
														break
													}
													break
												}
												v2235 = v2227
												v2236 = v2224
												v2239 = v2225
											} else {
												v2235 = int32(0)
												v2236 = v2136
												v2239 = v2137
											}
											v2287 = v2123 + v560
											v2288 = v2124 + v559
											if v569 == int32(0) {
												v2300 = v2236
												v2301 = v2239
											} else {
												v2292 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2288+v2235))))
												if v2292 == int32(0) {
													v2300 = v2236
													v2301 = v2239
												} else {
													v2298 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2287+v2235))))
													v2300 = v2236 + int32(1)
													v2301 = v2239 + v2298
												}
											}
											if v1153 != 0 {
												v2303 = v646 + v1089
												v2313 = int32(0)
												v2314 = v2300
												v2317 = v2301
												for {
													v2365 = v647 + v1089 + v2313
													v2366 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2365))))
													if v2366 == int32(0) {
														v2374 = v2314
														v2375 = v2317
													} else {
														v2372 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2303+v2313))))
														v2374 = v2314 + int32(1)
														v2375 = v2317 + v2372
													}
													v2378 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2365+int32(1)))))
													if v2378 == int32(0) {
														v2388 = v2374
														v2389 = v2375
													} else {
														v2381 = int32(1)
														v2386 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2303+v2313+v2381))))
														v2388 = v2374 + v2381
														v2389 = v2375 + v2386
													}
													v2391 = v2313 + int32(2)
													if v565-v1089 != v2391 {
														v2313 = v2391
														v2314 = v2388
														v2317 = v2389
														continue
													} else {
														break
													}
													break
												}
												v2399 = v2391
												v2400 = v2388
												v2403 = v2389
											} else {
												v2399 = int32(0)
												v2400 = v2300
												v2403 = v2301
											}
											v2451 = v2287 + v560
											v2452 = v2288 + v559
											if v569 == int32(0) {
												v2464 = v2400
												v2465 = v2403
											} else {
												v2456 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2452+v2399))))
												if v2456 == int32(0) {
													v2464 = v2400
													v2465 = v2403
												} else {
													v2462 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2451+v2399))))
													v2464 = v2400 + int32(1)
													v2465 = v2403 + v2462
												}
											}
											if v2464 < int32(1) {
											} else {
												if v1146<<(uint(int32(3))%32) <= v2464 {
												} else {
													v2471 = base.I32_div_s(v2465, v2464)
													if v1153 != 0 {
														v2485 = int32(0)
														for {
															v2533 = v1150 + v2485
															v2534 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2533))))
															if v2534 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v1149+v2485))) = uint8(v2471)
															}
															v2539 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2533+int32(1)))))
															if v2539 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v1149+v2485+int32(1)))) = uint8(v2471)
															}
															v2545 = v2485 + int32(2)
															if v565-v1089 != v2545 {
																v2485 = v2545
																continue
															} else {
																break
															}
															break
														}
														v2557 = v2545
													} else {
														v2557 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v2608 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1150+v2557))))
														if v2608 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v1149+v2557))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v2612 = v641 + v1089
														v2626 = int32(0)
														for {
															v2674 = v642 + v1089 + v2626
															v2675 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2674))))
															if v2675 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v2612+v2626))) = uint8(v2471)
															}
															v2680 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2674+int32(1)))))
															if v2680 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v2612+v2626+int32(1)))) = uint8(v2471)
															}
															v2686 = v2626 + int32(2)
															if v565-v1089 != v2686 {
																v2626 = v2686
																continue
															} else {
																break
															}
															break
														}
														v2698 = v2686
													} else {
														v2698 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v2749 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1468+v2698))))
														if v2749 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v1467+v2698))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v2753 = v656 + v1089
														v2767 = int32(0)
														for {
															v2815 = v657 + v1089 + v2767
															v2816 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2815))))
															if v2816 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v2753+v2767))) = uint8(v2471)
															}
															v2821 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2815+int32(1)))))
															if v2821 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v2753+v2767+int32(1)))) = uint8(v2471)
															}
															v2827 = v2767 + int32(2)
															if v565-v1089 != v2827 {
																v2767 = v2827
																continue
															} else {
																break
															}
															break
														}
														v2839 = v2827
													} else {
														v2839 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v2890 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1632+v2839))))
														if v2890 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v1631+v2839))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v2894 = v654 + v1089
														v2908 = int32(0)
														for {
															v2956 = v655 + v1089 + v2908
															v2957 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2956))))
															if v2957 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v2894+v2908))) = uint8(v2471)
															}
															v2962 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2956+int32(1)))))
															if v2962 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v2894+v2908+int32(1)))) = uint8(v2471)
															}
															v2968 = v2908 + int32(2)
															if v565-v1089 != v2968 {
																v2908 = v2968
																continue
															} else {
																break
															}
															break
														}
														v2980 = v2968
													} else {
														v2980 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v3031 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1796+v2980))))
														if v3031 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v1795+v2980))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v3035 = v652 + v1089
														v3049 = int32(0)
														for {
															v3097 = v653 + v1089 + v3049
															v3098 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3097))))
															if v3098 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3035+v3049))) = uint8(v2471)
															}
															v3103 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3097+int32(1)))))
															if v3103 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3035+v3049+int32(1)))) = uint8(v2471)
															}
															v3109 = v3049 + int32(2)
															if v565-v1089 != v3109 {
																v3049 = v3109
																continue
															} else {
																break
															}
															break
														}
														v3121 = v3109
													} else {
														v3121 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v3172 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1960+v3121))))
														if v3172 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v1959+v3121))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v3176 = v650 + v1089
														v3190 = int32(0)
														for {
															v3238 = v651 + v1089 + v3190
															v3239 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3238))))
															if v3239 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3176+v3190))) = uint8(v2471)
															}
															v3244 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3238+int32(1)))))
															if v3244 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3176+v3190+int32(1)))) = uint8(v2471)
															}
															v3250 = v3190 + int32(2)
															if v565-v1089 != v3250 {
																v3190 = v3250
																continue
															} else {
																break
															}
															break
														}
														v3262 = v3250
													} else {
														v3262 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v3313 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2124+v3262))))
														if v3313 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v2123+v3262))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v3317 = v648 + v1089
														v3331 = int32(0)
														for {
															v3379 = v649 + v1089 + v3331
															v3380 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3379))))
															if v3380 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3317+v3331))) = uint8(v2471)
															}
															v3385 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3379+int32(1)))))
															if v3385 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3317+v3331+int32(1)))) = uint8(v2471)
															}
															v3391 = v3331 + int32(2)
															if v565-v1089 != v3391 {
																v3331 = v3391
																continue
															} else {
																break
															}
															break
														}
														v3403 = v3391
													} else {
														v3403 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v3454 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2288+v3403))))
														if v3454 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v2287+v3403))) = uint8(v2471)
														}
													}
													if v1153 != 0 {
														v3458 = v646 + v1089
														v3472 = int32(0)
														for {
															v3520 = v647 + v1089 + v3472
															v3521 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3520))))
															if v3521 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3458+v3472))) = uint8(v2471)
															}
															v3526 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3520+int32(1)))))
															if v3526 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v3458+v3472+int32(1)))) = uint8(v2471)
															}
															v3532 = v3472 + int32(2)
															if v565-v1089 != v3532 {
																v3472 = v3532
																continue
															} else {
																break
															}
															break
														}
														v3544 = v3532
													} else {
														v3544 = int32(0)
													}
													if v569 == int32(0) {
													} else {
														v3595 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2452+v3544))))
														if v3595 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v2451+v3544))) = uint8(v2471)
														}
													}
												}
											}
										}
									}
									v3672 = v634 + v576
									v3673 = v633 + v578
									v3675 = v659 + int32(8)
									if v3675 <= v61 {
										v633 = v3673
										v634 = v3672
										v637 = v637 + v574
										v638 = v638 + v574
										v641 = v641 + v576
										v642 = v642 + v578
										v646 = v646 + v576
										v647 = v647 + v578
										v648 = v648 + v576
										v649 = v649 + v578
										v650 = v650 + v576
										v651 = v651 + v578
										v652 = v652 + v576
										v653 = v653 + v578
										v654 = v654 + v576
										v655 = v655 + v578
										v656 = v656 + v576
										v657 = v657 + v578
										v658 = v1126
										v659 = v3675
										v660 = v1128
										v661 = v1129
										continue
									} else {
										break
									}
									break
								}
								v3691 = v3673
								v3692 = v3672
								v3697 = v61 & int32(2147483640)
							}
							if v61 <= v3697 {
							} else {
								v3736 = v61 - v3697
								if v64 < int32(8) {
									v4163 = int32(0)
								} else {
									if v3736 < int32(1) {
										v4163 = v64 & int32(2147483640)
									} else {
										v3753 = int32(8)
										v3754 = int32(0)
										for {
											v3804 = int32(0)
											v3808 = v3736
											v3810 = v3691
											v3813 = v3804
											v3815 = v3692
											v3816 = v3804
											for {
												v3864 = v3815 + v3754
												v3865 = v3810 + v3754
												v3866 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865))))
												if v3866 == int32(0) {
													v3873 = v3813
													v3874 = v3816
												} else {
													v3871 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864))))
													v3873 = v3813 + int32(1)
													v3874 = v3816 + v3871
												}
												v3877 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(1)))))
												if v3877 == int32(0) {
													v3886 = v3873
													v3887 = v3874
												} else {
													v3880 = int32(1)
													v3884 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+v3880))))
													v3886 = v3873 + v3880
													v3887 = v3874 + v3884
												}
												v3890 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(2)))))
												if v3890 == int32(0) {
													v3899 = v3886
													v3900 = v3887
												} else {
													v3897 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+int32(2)))))
													v3899 = v3886 + int32(1)
													v3900 = v3887 + v3897
												}
												v3903 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(3)))))
												if v3903 == int32(0) {
													v3912 = v3899
													v3913 = v3900
												} else {
													v3910 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+int32(3)))))
													v3912 = v3899 + int32(1)
													v3913 = v3900 + v3910
												}
												v3916 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(4)))))
												if v3916 == int32(0) {
													v3925 = v3912
													v3926 = v3913
												} else {
													v3923 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+int32(4)))))
													v3925 = v3912 + int32(1)
													v3926 = v3913 + v3923
												}
												v3929 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(5)))))
												if v3929 == int32(0) {
													v3938 = v3925
													v3939 = v3926
												} else {
													v3936 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+int32(5)))))
													v3938 = v3925 + int32(1)
													v3939 = v3926 + v3936
												}
												v3942 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(6)))))
												if v3942 == int32(0) {
													v3951 = v3938
													v3952 = v3939
												} else {
													v3949 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+int32(6)))))
													v3951 = v3938 + int32(1)
													v3952 = v3939 + v3949
												}
												v3955 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3865+int32(7)))))
												if v3955 == int32(0) {
													v3964 = v3951
													v3965 = v3952
												} else {
													v3962 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v3864+int32(7)))))
													v3964 = v3951 + int32(1)
													v3965 = v3952 + v3962
												}
												v3969 = v3808 + int32(-1)
												if v3969 != 0 {
													v3808 = v3969
													v3810 = v3810 + v559
													v3813 = v3964
													v3815 = v3815 + v560
													v3816 = v3965
													continue
												} else {
													break
												}
												break
											}
											if v3964 < int32(1) {
											} else {
												if v3736<<(uint(int32(3))%32) <= v3964 {
												} else {
													v3973 = base.I32_div_s(v3965, v3964)
													v3974 = v3691
													v3978 = v3692
													v3983 = v3736
													for {
														v4032 = v3978 + v3754
														v4033 = v3974 + v3754
														v4034 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033))))
														if v4034 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032))) = uint8(v3973)
														}
														v4038 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(1)))))
														if v4038 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(1)))) = uint8(v3973)
														}
														v4044 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(2)))))
														if v4044 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(2)))) = uint8(v3973)
														}
														v4050 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(3)))))
														if v4050 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(3)))) = uint8(v3973)
														}
														v4056 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(4)))))
														if v4056 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(4)))) = uint8(v3973)
														}
														v4062 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(5)))))
														if v4062 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(5)))) = uint8(v3973)
														}
														v4068 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(6)))))
														if v4068 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(6)))) = uint8(v3973)
														}
														v4074 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4033+int32(7)))))
														if v4074 != 0 {
														} else {
															*(*uint8)(unsafe.Add(mBase, uint32(v4032+int32(7)))) = uint8(v3973)
														}
														v4081 = v3983 + int32(-1)
														if v4081 != 0 {
															v3974 = v3974 + v559
															v3978 = v3978 + v560
															v3983 = v4081
															continue
														} else {
															break
														}
														break
													}
												}
											}
											v4140 = int32(8)
											v4143 = v3753 + v4140
											if v4143 <= v64 {
												v3753 = v4143
												v3754 = v3754 + v4140
												continue
											} else {
												break
											}
											break
										}
										v4163 = v3753
									}
								}
								if v64 <= v4163 {
								} else {
									if v3736 < int32(1) {
									} else {
										v4208 = v64 - v4163
										if v4208 < int32(1) {
										} else {
											v4218 = v64 + (v4163 ^ int32(-1))
											v4219 = int32(0)
											v4220 = v3692 + v4163
											v4221 = v3691 + v4163
											v4228 = v4221
											v4230 = v4219
											v4232 = v4220
											v4234 = v4219
											v4238 = v4219
											for {
												if v4218 != 0 {
													v4290 = v4230
													v4291 = int32(0)
													v4294 = v4234
													for {
														v4342 = v4228 + v4291
														v4343 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4342))))
														if v4343 == int32(0) {
															v4351 = v4290
															v4352 = v4294
														} else {
															v4349 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4232+v4291))))
															v4351 = v4290 + v4349
															v4352 = v4294 + int32(1)
														}
														v4355 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4342+int32(1)))))
														if v4355 == int32(0) {
															v4365 = v4351
															v4366 = v4352
														} else {
															v4358 = int32(1)
															v4363 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4232+v4291+v4358))))
															v4365 = v4351 + v4363
															v4366 = v4352 + v4358
														}
														v4368 = v4291 + int32(2)
														if v64&int32(-2)-v4163 != v4368 {
															v4290 = v4365
															v4291 = v4368
															v4294 = v4366
															continue
														} else {
															break
														}
														break
													}
													v4376 = v4365
													v4377 = v4368
													v4380 = v4366
												} else {
													v4376 = v4230
													v4377 = int32(0)
													v4380 = v4234
												}
												if v64&int32(1) == int32(0) {
													v4439 = v4376
													v4440 = v4380
												} else {
													v4431 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4228+v4377))))
													if v4431 == int32(0) {
														v4439 = v4376
														v4440 = v4380
													} else {
														v4437 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4232+v4377))))
														v4439 = v4376 + v4437
														v4440 = v4380 + int32(1)
													}
												}
												v4444 = v4238 + int32(1)
												if v4444 != v3736 {
													v4228 = v4228 + v559
													v4230 = v4439
													v4232 = v4232 + v560
													v4234 = v4440
													v4238 = v4444
													continue
												} else {
													break
												}
												break
											}
											if v4440 < int32(1) {
											} else {
												if v4208*v3736 <= v4440 {
												} else {
													v4452 = base.I32_div_s(v4439, v4440)
													v4459 = v4221
													v4465 = int32(0)
													v4469 = v4220
													for {
														v4515 = int32(0)
														if v4218 == v4515 {
															v4598 = v4515
														} else {
															v4526 = int32(0)
															for {
																v4577 = v4459 + v4526
																v4578 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4577))))
																if v4578 != 0 {
																} else {
																	*(*uint8)(unsafe.Add(mBase, uint32(v4469+v4526))) = uint8(v4452)
																}
																v4583 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4577+int32(1)))))
																if v4583 != 0 {
																} else {
																	*(*uint8)(unsafe.Add(mBase, uint32(v4469+v4526+int32(1)))) = uint8(v4452)
																}
																v4589 = v4526 + int32(2)
																if v64&int32(-2)-v4163 != v4589 {
																	v4526 = v4589
																	continue
																} else {
																	break
																}
																break
															}
															v4598 = v4589
														}
														if v64&int32(1) == int32(0) {
														} else {
															v4652 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4459+v4598))))
															if v4652 != 0 {
															} else {
																*(*uint8)(unsafe.Add(mBase, uint32(v4469+v4598))) = uint8(v4452)
															}
														}
														v4658 = v4465 + int32(1)
														if v4658 != v3736 {
															v4459 = v4459 + v559
															v4465 = v4658
															v4469 = v4469 + v560
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
					}
				}
			}
		} else {
			if v61 < int32(8) {
			} else {
				if v64 < int32(8) {
				} else {
					v74 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
					v75 = int32(0)
					v80 = v75
					v83 = v75
					for {
						v140 = v80
						v144 = int32(0)
						v145 = int32(1)
						for {
							v195 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
							v200 = v74 + (v144+v83*v195)<<(uint(int32(5))%32)
							v201 = *(*int32)(unsafe.Add(mBase, uint32(v200)))
							if base.Ui32(v201) <= base.Ui32(int32(16777215)) {
								v205 = *(*int32)(unsafe.Add(mBase, uint32(v200)+4))
								if base.Ui32(v205) <= base.Ui32(int32(16777215)) {
									v209 = *(*int32)(unsafe.Add(mBase, uint32(v200)+8))
									if base.Ui32(v209) <= base.Ui32(int32(16777215)) {
										v213 = *(*int32)(unsafe.Add(mBase, uint32(v200)+12))
										if base.Ui32(v213) <= base.Ui32(int32(16777215)) {
											v217 = *(*int32)(unsafe.Add(mBase, uint32(v200)+16))
											if base.Ui32(v217) <= base.Ui32(int32(16777215)) {
												v221 = *(*int32)(unsafe.Add(mBase, uint32(v200)+20))
												if base.Ui32(v221) <= base.Ui32(int32(16777215)) {
													v225 = *(*int32)(unsafe.Add(mBase, uint32(v200)+24))
													if base.Ui32(v225) <= base.Ui32(int32(16777215)) {
														v229 = *(*int32)(unsafe.Add(mBase, uint32(v200)+28))
														if base.Ui32(v229) <= base.Ui32(int32(16777215)) {
															v234 = v195 << (uint(int32(2)) % 32)
															v235 = v200 + v234
															v236 = *(*int32)(unsafe.Add(mBase, uint32(v235)))
															if base.Ui32(v236) <= base.Ui32(int32(16777215)) {
																v240 = *(*int32)(unsafe.Add(mBase, uint32(v235)+4))
																if base.Ui32(v240) <= base.Ui32(int32(16777215)) {
																	v244 = *(*int32)(unsafe.Add(mBase, uint32(v235)+8))
																	if base.Ui32(v244) <= base.Ui32(int32(16777215)) {
																		v248 = *(*int32)(unsafe.Add(mBase, uint32(v235)+12))
																		if base.Ui32(v248) <= base.Ui32(int32(16777215)) {
																			v252 = *(*int32)(unsafe.Add(mBase, uint32(v235)+16))
																			if base.Ui32(v252) <= base.Ui32(int32(16777215)) {
																				v256 = *(*int32)(unsafe.Add(mBase, uint32(v235)+20))
																				if base.Ui32(v256) <= base.Ui32(int32(16777215)) {
																					v260 = *(*int32)(unsafe.Add(mBase, uint32(v235)+24))
																					if base.Ui32(v260) <= base.Ui32(int32(16777215)) {
																						v264 = *(*int32)(unsafe.Add(mBase, uint32(v235)+28))
																						if base.Ui32(v264) <= base.Ui32(int32(16777215)) {
																							v268 = v235 + v234
																							v269 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
																							if base.Ui32(v269) <= base.Ui32(int32(16777215)) {
																								v273 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
																								if base.Ui32(v273) <= base.Ui32(int32(16777215)) {
																									v277 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
																									if base.Ui32(v277) <= base.Ui32(int32(16777215)) {
																										v281 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
																										if base.Ui32(v281) <= base.Ui32(int32(16777215)) {
																											v285 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
																											if base.Ui32(v285) <= base.Ui32(int32(16777215)) {
																												v289 = *(*int32)(unsafe.Add(mBase, uint32(v268)+20))
																												if base.Ui32(v289) <= base.Ui32(int32(16777215)) {
																													v293 = *(*int32)(unsafe.Add(mBase, uint32(v268)+24))
																													if base.Ui32(v293) <= base.Ui32(int32(16777215)) {
																														v297 = *(*int32)(unsafe.Add(mBase, uint32(v268)+28))
																														if base.Ui32(v297) <= base.Ui32(int32(16777215)) {
																															v301 = v268 + v234
																															v302 = *(*int32)(unsafe.Add(mBase, uint32(v301)))
																															if base.Ui32(v302) <= base.Ui32(int32(16777215)) {
																																v306 = *(*int32)(unsafe.Add(mBase, uint32(v301)+4))
																																if base.Ui32(v306) <= base.Ui32(int32(16777215)) {
																																	v310 = *(*int32)(unsafe.Add(mBase, uint32(v301)+8))
																																	if base.Ui32(v310) <= base.Ui32(int32(16777215)) {
																																		v314 = *(*int32)(unsafe.Add(mBase, uint32(v301)+12))
																																		if base.Ui32(v314) <= base.Ui32(int32(16777215)) {
																																			v318 = *(*int32)(unsafe.Add(mBase, uint32(v301)+16))
																																			if base.Ui32(v318) <= base.Ui32(int32(16777215)) {
																																				v322 = *(*int32)(unsafe.Add(mBase, uint32(v301)+20))
																																				if base.Ui32(v322) <= base.Ui32(int32(16777215)) {
																																					v326 = *(*int32)(unsafe.Add(mBase, uint32(v301)+24))
																																					if base.Ui32(v326) <= base.Ui32(int32(16777215)) {
																																						v330 = *(*int32)(unsafe.Add(mBase, uint32(v301)+28))
																																						if base.Ui32(v330) <= base.Ui32(int32(16777215)) {
																																							v334 = v301 + v234
																																							v335 = *(*int32)(unsafe.Add(mBase, uint32(v334)))
																																							if base.Ui32(v335) <= base.Ui32(int32(16777215)) {
																																								v339 = *(*int32)(unsafe.Add(mBase, uint32(v334)+4))
																																								if base.Ui32(v339) <= base.Ui32(int32(16777215)) {
																																									v343 = *(*int32)(unsafe.Add(mBase, uint32(v334)+8))
																																									if base.Ui32(v343) <= base.Ui32(int32(16777215)) {
																																										v347 = *(*int32)(unsafe.Add(mBase, uint32(v334)+12))
																																										if base.Ui32(v347) <= base.Ui32(int32(16777215)) {
																																											v351 = *(*int32)(unsafe.Add(mBase, uint32(v334)+16))
																																											if base.Ui32(v351) <= base.Ui32(int32(16777215)) {
																																												v355 = *(*int32)(unsafe.Add(mBase, uint32(v334)+20))
																																												if base.Ui32(v355) <= base.Ui32(int32(16777215)) {
																																													v359 = *(*int32)(unsafe.Add(mBase, uint32(v334)+24))
																																													if base.Ui32(v359) <= base.Ui32(int32(16777215)) {
																																														v363 = *(*int32)(unsafe.Add(mBase, uint32(v334)+28))
																																														if base.Ui32(v363) <= base.Ui32(int32(16777215)) {
																																															v367 = v334 + v234
																																															v368 = *(*int32)(unsafe.Add(mBase, uint32(v367)))
																																															if base.Ui32(v368) <= base.Ui32(int32(16777215)) {
																																																v372 = *(*int32)(unsafe.Add(mBase, uint32(v367)+4))
																																																if base.Ui32(v372) <= base.Ui32(int32(16777215)) {
																																																	v376 = *(*int32)(unsafe.Add(mBase, uint32(v367)+8))
																																																	if base.Ui32(v376) <= base.Ui32(int32(16777215)) {
																																																		v380 = *(*int32)(unsafe.Add(mBase, uint32(v367)+12))
																																																		if base.Ui32(v380) <= base.Ui32(int32(16777215)) {
																																																			v384 = *(*int32)(unsafe.Add(mBase, uint32(v367)+16))
																																																			if base.Ui32(v384) <= base.Ui32(int32(16777215)) {
																																																				v388 = *(*int32)(unsafe.Add(mBase, uint32(v367)+20))
																																																				if base.Ui32(v388) <= base.Ui32(int32(16777215)) {
																																																					v392 = *(*int32)(unsafe.Add(mBase, uint32(v367)+24))
																																																					if base.Ui32(v392) <= base.Ui32(int32(16777215)) {
																																																						v396 = *(*int32)(unsafe.Add(mBase, uint32(v367)+28))
																																																						if base.Ui32(v396) <= base.Ui32(int32(16777215)) {
																																																							v400 = v367 + v234
																																																							v401 = *(*int32)(unsafe.Add(mBase, uint32(v400)))
																																																							if base.Ui32(v401) <= base.Ui32(int32(16777215)) {
																																																								v405 = *(*int32)(unsafe.Add(mBase, uint32(v400)+4))
																																																								if base.Ui32(v405) <= base.Ui32(int32(16777215)) {
																																																									v409 = *(*int32)(unsafe.Add(mBase, uint32(v400)+8))
																																																									if base.Ui32(v409) <= base.Ui32(int32(16777215)) {
																																																										v413 = *(*int32)(unsafe.Add(mBase, uint32(v400)+12))
																																																										if base.Ui32(v413) <= base.Ui32(int32(16777215)) {
																																																											v417 = *(*int32)(unsafe.Add(mBase, uint32(v400)+16))
																																																											if base.Ui32(v417) <= base.Ui32(int32(16777215)) {
																																																												v421 = *(*int32)(unsafe.Add(mBase, uint32(v400)+20))
																																																												if base.Ui32(v421) <= base.Ui32(int32(16777215)) {
																																																													v425 = *(*int32)(unsafe.Add(mBase, uint32(v400)+24))
																																																													if base.Ui32(v425) <= base.Ui32(int32(16777215)) {
																																																														v429 = *(*int32)(unsafe.Add(mBase, uint32(v400)+28))
																																																														if base.Ui32(v429) <= base.Ui32(int32(16777215)) {
																																																															v433 = v400 + v234
																																																															v434 = *(*int32)(unsafe.Add(mBase, uint32(v433)))
																																																															if base.Ui32(v434) <= base.Ui32(int32(16777215)) {
																																																																v438 = *(*int32)(unsafe.Add(mBase, uint32(v433)+4))
																																																																if base.Ui32(v438) <= base.Ui32(int32(16777215)) {
																																																																	v442 = *(*int32)(unsafe.Add(mBase, uint32(v433)+8))
																																																																	if base.Ui32(v442) <= base.Ui32(int32(16777215)) {
																																																																		v446 = *(*int32)(unsafe.Add(mBase, uint32(v433)+12))
																																																																		if base.Ui32(v446) <= base.Ui32(int32(16777215)) {
																																																																			v450 = *(*int32)(unsafe.Add(mBase, uint32(v433)+16))
																																																																			if base.Ui32(v450) <= base.Ui32(int32(16777215)) {
																																																																				v454 = *(*int32)(unsafe.Add(mBase, uint32(v433)+20))
																																																																				if base.Ui32(v454) <= base.Ui32(int32(16777215)) {
																																																																					v458 = *(*int32)(unsafe.Add(mBase, uint32(v433)+24))
																																																																					if base.Ui32(v458) <= base.Ui32(int32(16777215)) {
																																																																						v462 = *(*int32)(unsafe.Add(mBase, uint32(v433)+28))
																																																																						if base.Ui32(v462) <= base.Ui32(int32(16777215)) {
																																																																							if v145 != 0 {
																																																																								v466 = v201
																																																																							} else {
																																																																								v466 = v140
																																																																							}
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v200))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v235))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v268))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v301))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v334))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v367))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v400))) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+28)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+24)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+20)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+16)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+12)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+8)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433)+4)) = v466
																																																																							*(*int32)(unsafe.Add(mBase, uint32(v433))) = v466
																																																																							v532 = v466
																																																																							v533 = int32(0)
																																																																						} else {
																																																																							v532 = v140
																																																																							v533 = int32(1)
																																																																						}
																																																																					} else {
																																																																						v532 = v140
																																																																						v533 = int32(1)
																																																																					}
																																																																				} else {
																																																																					v532 = v140
																																																																					v533 = int32(1)
																																																																				}
																																																																			} else {
																																																																				v532 = v140
																																																																				v533 = int32(1)
																																																																			}
																																																																		} else {
																																																																			v532 = v140
																																																																			v533 = int32(1)
																																																																		}
																																																																	} else {
																																																																		v532 = v140
																																																																		v533 = int32(1)
																																																																	}
																																																																} else {
																																																																	v532 = v140
																																																																	v533 = int32(1)
																																																																}
																																																															} else {
																																																																v532 = v140
																																																																v533 = int32(1)
																																																															}
																																																														} else {
																																																															v532 = v140
																																																															v533 = int32(1)
																																																														}
																																																													} else {
																																																														v532 = v140
																																																														v533 = int32(1)
																																																													}
																																																												} else {
																																																													v532 = v140
																																																													v533 = int32(1)
																																																												}
																																																											} else {
																																																												v532 = v140
																																																												v533 = int32(1)
																																																											}
																																																										} else {
																																																											v532 = v140
																																																											v533 = int32(1)
																																																										}
																																																									} else {
																																																										v532 = v140
																																																										v533 = int32(1)
																																																									}
																																																								} else {
																																																									v532 = v140
																																																									v533 = int32(1)
																																																								}
																																																							} else {
																																																								v532 = v140
																																																								v533 = int32(1)
																																																							}
																																																						} else {
																																																							v532 = v140
																																																							v533 = int32(1)
																																																						}
																																																					} else {
																																																						v532 = v140
																																																						v533 = int32(1)
																																																					}
																																																				} else {
																																																					v532 = v140
																																																					v533 = int32(1)
																																																				}
																																																			} else {
																																																				v532 = v140
																																																				v533 = int32(1)
																																																			}
																																																		} else {
																																																			v532 = v140
																																																			v533 = int32(1)
																																																		}
																																																	} else {
																																																		v532 = v140
																																																		v533 = int32(1)
																																																	}
																																																} else {
																																																	v532 = v140
																																																	v533 = int32(1)
																																																}
																																															} else {
																																																v532 = v140
																																																v533 = int32(1)
																																															}
																																														} else {
																																															v532 = v140
																																															v533 = int32(1)
																																														}
																																													} else {
																																														v532 = v140
																																														v533 = int32(1)
																																													}
																																												} else {
																																													v532 = v140
																																													v533 = int32(1)
																																												}
																																											} else {
																																												v532 = v140
																																												v533 = int32(1)
																																											}
																																										} else {
																																											v532 = v140
																																											v533 = int32(1)
																																										}
																																									} else {
																																										v532 = v140
																																										v533 = int32(1)
																																									}
																																								} else {
																																									v532 = v140
																																									v533 = int32(1)
																																								}
																																							} else {
																																								v532 = v140
																																								v533 = int32(1)
																																							}
																																						} else {
																																							v532 = v140
																																							v533 = int32(1)
																																						}
																																					} else {
																																						v532 = v140
																																						v533 = int32(1)
																																					}
																																				} else {
																																					v532 = v140
																																					v533 = int32(1)
																																				}
																																			} else {
																																				v532 = v140
																																				v533 = int32(1)
																																			}
																																		} else {
																																			v532 = v140
																																			v533 = int32(1)
																																		}
																																	} else {
																																		v532 = v140
																																		v533 = int32(1)
																																	}
																																} else {
																																	v532 = v140
																																	v533 = int32(1)
																																}
																															} else {
																																v532 = v140
																																v533 = int32(1)
																															}
																														} else {
																															v532 = v140
																															v533 = int32(1)
																														}
																													} else {
																														v532 = v140
																														v533 = int32(1)
																													}
																												} else {
																													v532 = v140
																													v533 = int32(1)
																												}
																											} else {
																												v532 = v140
																												v533 = int32(1)
																											}
																										} else {
																											v532 = v140
																											v533 = int32(1)
																										}
																									} else {
																										v532 = v140
																										v533 = int32(1)
																									}
																								} else {
																									v532 = v140
																									v533 = int32(1)
																								}
																							} else {
																								v532 = v140
																								v533 = int32(1)
																							}
																						} else {
																							v532 = v140
																							v533 = int32(1)
																						}
																					} else {
																						v532 = v140
																						v533 = int32(1)
																					}
																				} else {
																					v532 = v140
																					v533 = int32(1)
																				}
																			} else {
																				v532 = v140
																				v533 = int32(1)
																			}
																		} else {
																			v532 = v140
																			v533 = int32(1)
																		}
																	} else {
																		v532 = v140
																		v533 = int32(1)
																	}
																} else {
																	v532 = v140
																	v533 = int32(1)
																}
															} else {
																v532 = v140
																v533 = int32(1)
															}
														} else {
															v532 = v140
															v533 = int32(1)
														}
													} else {
														v532 = v140
														v533 = int32(1)
													}
												} else {
													v532 = v140
													v533 = int32(1)
												}
											} else {
												v532 = v140
												v533 = int32(1)
											}
										} else {
											v532 = v140
											v533 = int32(1)
										}
									} else {
										v532 = v140
										v533 = int32(1)
									}
								} else {
									v532 = v140
									v533 = int32(1)
								}
							} else {
								v532 = v140
								v533 = int32(1)
							}
							v542 = v144 + int32(1)
							if v66 != v542 {
								v140 = v532
								v144 = v542
								v145 = v533
								continue
							} else {
								break
							}
							break
						}
						v545 = v83 + int32(1)
						if v545 != v63 {
							v80 = v532
							v83 = v545
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
	return
}
func F_WebPConfigInitInternal(m *base.Module, l0 int32, l1 int32, l2 float32, l3 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v13 int64
	_ = v13
	var v20 int64
	_ = v20
	var v34 int32
	_ = v34
	var v66 int32
	_ = v66
	var v73 float32
	_ = v73
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v82 float32
	_ = v82
	var v85 int32
	_ = v85
	var v88 int32
	_ = v88
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
	var v113 int32
	_ = v113
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v121 int32
	_ = v121
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v130 int32
	_ = v130
	var v133 int32
	_ = v133
	var v134 int32
	_ = v134
	var v137 int32
	_ = v137
	var v138 int32
	_ = v138
	var v141 int32
	_ = v141
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v159 int32
	_ = v159
	var v162 int32
	_ = v162
	var v165 int32
	_ = v165
	var v168 int32
	_ = v168
	var v172 int32
	_ = v172
	v5 = int32(0)
	if l0 == v5 {
		v172 = v5
	} else {
		if l3&int32(-256) != int32(512) {
			v172 = v5
		} else {
			v13 = int64(0)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v13
			*(*float32)(unsafe.Add(mBase, uint32(l0)+4)) = l2
			*(*int64)(unsafe.Add(mBase, uint32(l0)+32)) = int64(60)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+8)) = int64(4)
			v20 = int64(429496729600)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+108)) = v20
			*(*int64)(unsafe.Add(mBase, uint32(l0)+24)) = int64(214748364804)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+64)) = v13
			*(*int64)(unsafe.Add(mBase, uint32(l0)+72)) = v13
			*(*int64)(unsafe.Add(mBase, uint32(l0)+40)) = int64(1)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+56)) = int64(4294967396)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+48)) = int64(4294967297)
			v34 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(l0))) = v34
			*(*int32)(unsafe.Add(mBase, uint32(l0)+96)) = v34
			*(*int64)(unsafe.Add(mBase, uint32(l0)+88)) = v20
			*(*int64)(unsafe.Add(mBase, uint32(l0)+80)) = v13
			*(*int32)(unsafe.Add(mBase, uint32(l0)+104)) = v34
			switch l1 + int32(-1) {
			case 0:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = int32(4)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(150323855440)
			case 1:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = int32(3)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(128849018960)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+68)) = int32(2)
			case 2:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = int32(6)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(42949672985)
			case 3:
				*(*int64)(unsafe.Add(mBase, uint32(l0)+28)) = int64(0)
			case 4:
				*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = int32(0)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+24)) = int64(2)
			default:
			}
			v66 = int32(0)
			if l0 == v66 {
				v168 = v66
			} else {
				v73 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
				if base.F32_lt(v73, float32(0)) != 0 {
					v168 = v66
				} else {
					if base.F32_gt(v73, float32(100)) != 0 {
						v168 = v66
					} else {
						v78 = int32(0)
						v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
						if v79 < v78 {
							v168 = v78
						} else {
							v82 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
							if base.F32_lt(v82, float32(0)) != 0 {
								v168 = v78
							} else {
								v85 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								if base.Ui32(int32(6)) < base.Ui32(v85) {
									v168 = v78
								} else {
									v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
									if base.Ui32(v88+int32(-5)) < base.Ui32(int32(-4)) {
										v168 = v78
									} else {
										v93 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
										if base.Ui32(int32(100)) < base.Ui32(v93) {
											v168 = v78
										} else {
											v96 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
											if base.Ui32(int32(100)) < base.Ui32(v96) {
												v168 = v78
											} else {
												v99 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
												if base.Ui32(int32(7)) < base.Ui32(v99) {
													v168 = v78
												} else {
													v102 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
													if base.Ui32(int32(1)) < base.Ui32(v102) {
														v168 = v78
													} else {
														v105 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
														if base.Ui32(int32(1)) < base.Ui32(v105) {
															v168 = v78
														} else {
															v108 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
															if base.Ui32(v108+int32(-11)) < base.Ui32(int32(-10)) {
																v168 = v78
															} else {
																v113 = int32(0)
																v114 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
																if v114 < v113 {
																	v168 = v113
																} else {
																	v117 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
																	if int32(100) < v117 {
																		v168 = v113
																	} else {
																		if v117 < v114 {
																			v168 = v113
																		} else {
																			v121 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
																			if base.Ui32(int32(1)) < base.Ui32(v121) {
																				v168 = v113
																			} else {
																				v124 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
																				if base.Ui32(int32(7)) < base.Ui32(v124) {
																					v168 = v113
																				} else {
																					v127 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
																					if base.Ui32(int32(3)) < base.Ui32(v127) {
																						v168 = v113
																					} else {
																						v130 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
																						if base.Ui32(int32(100)) < base.Ui32(v130) {
																							v168 = v113
																						} else {
																							v133 = int32(0)
																							v134 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
																							if v134 < v133 {
																								v168 = v133
																							} else {
																								v137 = int32(0)
																								v138 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
																								if v138 < v137 {
																									v168 = v137
																								} else {
																									v141 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
																									if base.Ui32(int32(100)) < base.Ui32(v141) {
																										v168 = v137
																									} else {
																										v144 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
																										if base.Ui32(int32(1)) < base.Ui32(v144) {
																											v168 = v137
																										} else {
																											v147 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
																											if base.Ui32(int32(100)) < base.Ui32(v147) {
																												v168 = v137
																											} else {
																												v150 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
																												if base.Ui32(int32(3)) < base.Ui32(v150) {
																													v168 = v137
																												} else {
																													v153 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
																													if base.Ui32(int32(1)) < base.Ui32(v153) {
																														v168 = v137
																													} else {
																														v156 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
																														if base.Ui32(int32(1)) < base.Ui32(v156) {
																															v168 = v137
																														} else {
																															v159 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
																															if base.Ui32(int32(1)) < base.Ui32(v159) {
																																v168 = v137
																															} else {
																																v162 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
																																if base.Ui32(int32(1)) < base.Ui32(v162) {
																																	v168 = v137
																																} else {
																																	v165 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
																																	v168 = base.B2i32(base.Ui32(v165) < base.Ui32(int32(2)))
																																}
																															}
																														}
																													}
																												}
																											}
																										}
																									}
																								}
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
			v172 = v168
		}
	}
	return v172
}
func F_WebPCopyPlane(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) {
	var v12 int32
	_ = v12
	var v15 int32
	_ = v15
	var v17 int32
	_ = v17
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v30 int32
	_ = v30
	var v35 int32
	_ = v35
	var v40 int32
	_ = v40
	var v42 int32
	_ = v42
	var v46 int32
	_ = v46
	var v48 int32
	_ = v48
	var v50 int32
	_ = v50
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v61 int32
	_ = v61
	if l5 < int32(1) {
	} else {
		v12 = l5 & int32(3)
		if v12 != 0 {
			v15 = l0
			v17 = l2
			v21 = v12
			for {
				v23 = F_memcpy(m, v17, v15, l4)
				v24 = v23 + l3
				v25 = v15 + l1
				v27 = v21 + int32(-1)
				if v27 != 0 {
					v15 = v25
					v17 = v24
					v21 = v27
					continue
				} else {
					break
				}
				break
			}
			v28 = v25
			v30 = v24
			v35 = l5 & int32(-4)
		} else {
			v28 = l0
			v30 = l2
			v35 = l5
		}
		if base.Ui32(l5) < base.Ui32(int32(4)) {
		} else {
			v40 = v28
			v42 = v30
			v46 = v35 + int32(-1)
			for {
				v48 = F_memcpy(m, v42, v40, l4)
				v50 = v40 + l1
				v51 = F_memcpy(m, v48+l3, v50, l4)
				v53 = v50 + l1
				v54 = F_memcpy(m, v51+l3, v53, l4)
				v56 = v53 + l1
				v57 = F_memcpy(m, v54+l3, v56, l4)
				v61 = v46 + int32(-4)
				if base.Ui32(v61) < base.Ui32(int32(-2)) {
					v40 = v56 + l1
					v42 = v57 + l3
					v46 = v61
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
func F_WebPEncode(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v14 int32
	_ = v14
	var v21 float32
	_ = v21
	var v26 int32
	_ = v26
	var v27 int32
	_ = v27
	var v30 float32
	_ = v30
	var v33 int32
	_ = v33
	var v36 int32
	_ = v36
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v47 int32
	_ = v47
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v65 int32
	_ = v65
	var v69 int32
	_ = v69
	var v72 int32
	_ = v72
	var v75 int32
	_ = v75
	var v78 int32
	_ = v78
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	var v107 int32
	_ = v107
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v116 int32
	_ = v116
	var v121 int32
	_ = v121
	var v130 int32
	_ = v130
	var v131 int32
	_ = v131
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v140 int32
	_ = v140
	var v142 int32
	_ = v142
	var v144 int32
	_ = v144
	var v147 int32
	_ = v147
	var v150 int32
	_ = v150
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v161 int32
	_ = v161
	var v286 int32
	_ = v286
	var v287 int32
	_ = v287
	var v288 int32
	_ = v288
	var v291 int32
	_ = v291
	var v294 int32
	_ = v294
	var v295 int32
	_ = v295
	var v296 int32
	_ = v296
	var v304 int32
	_ = v304
	var v306 int32
	_ = v306
	var v309 int32
	_ = v309
	var v311 int32
	_ = v311
	var v316 int32
	_ = v316
	var v321 int32
	_ = v321
	var v323 int32
	_ = v323
	var v331 float32
	_ = v331
	var v333 float32
	_ = v333
	var v334 float32
	_ = v334
	var v340 float32
	_ = v340
	var v345 int32
	_ = v345
	var v347 int32
	_ = v347
	var v354 int32
	_ = v354
	var v356 int32
	_ = v356
	var v363 int32
	_ = v363
	var v367 int32
	_ = v367
	var v369 int32
	_ = v369
	var v374 int32
	_ = v374
	var v376 int32
	_ = v376
	var v379 int32
	_ = v379
	var v382 int32
	_ = v382
	var v383 int32
	_ = v383
	var v386 int32
	_ = v386
	var v389 int32
	_ = v389
	var v390 int32
	_ = v390
	var v393 int32
	_ = v393
	var v396 int32
	_ = v396
	var v404 int32
	_ = v404
	var v411 int32
	_ = v411
	var v412 int32
	_ = v412
	var v417 int32
	_ = v417
	var v418 int32
	_ = v418
	var v424 int32
	_ = v424
	var v432 int32
	_ = v432
	var v440 int32
	_ = v440
	var v441 int32
	_ = v441
	var v442 int32
	_ = v442
	var v446 int32
	_ = v446
	var v449 int32
	_ = v449
	var v450 int32
	_ = v450
	v3 = int32(0)
	if l1 == v3 {
		v450 = v3
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v450
L2:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(0)
	if l0 != 0 {
		goto L3
	} else {
		goto L4
	}
L3:
	;
	v14 = int32(0)
	if l0 == v14 {
		v116 = v14
		goto L7
	} else {
		goto L8
	}
L4:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(3)
	return int32(0)
L5:
	;
	if l1 != 0 {
		goto L42
	} else {
		goto L43
	}
L6:
	;
	if v116 != 0 {
		goto L5
	} else {
		goto L38
	}
L7:
	;
	goto L6
L8:
	;
	v21 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
	if base.F32_lt(v21, float32(0)) != 0 {
		v116 = v14
		goto L7
	} else {
		goto L9
	}
L9:
	;
	if base.F32_gt(v21, float32(100)) != 0 {
		v116 = v14
		goto L7
	} else {
		goto L10
	}
L10:
	;
	v26 = int32(0)
	v27 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
	if v27 < v26 {
		v116 = v26
		goto L7
	} else {
		goto L11
	}
L11:
	;
	v30 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
	if base.F32_lt(v30, float32(0)) != 0 {
		v116 = v26
		goto L7
	} else {
		goto L12
	}
L12:
	;
	v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if base.Ui32(int32(6)) < base.Ui32(v33) {
		v116 = v26
		goto L7
	} else {
		goto L13
	}
L13:
	;
	v36 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
	if base.Ui32(v36+int32(-5)) < base.Ui32(int32(-4)) {
		v116 = v26
		goto L7
	} else {
		goto L14
	}
L14:
	;
	v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
	if base.Ui32(int32(100)) < base.Ui32(v41) {
		v116 = v26
		goto L7
	} else {
		goto L15
	}
L15:
	;
	v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
	if base.Ui32(int32(100)) < base.Ui32(v44) {
		v116 = v26
		goto L7
	} else {
		goto L16
	}
L16:
	;
	v47 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	if base.Ui32(int32(7)) < base.Ui32(v47) {
		v116 = v26
		goto L7
	} else {
		goto L17
	}
L17:
	;
	v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	if base.Ui32(int32(1)) < base.Ui32(v50) {
		v116 = v26
		goto L7
	} else {
		goto L18
	}
L18:
	;
	v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
	if base.Ui32(int32(1)) < base.Ui32(v53) {
		v116 = v26
		goto L7
	} else {
		goto L19
	}
L19:
	;
	v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
	if base.Ui32(v56+int32(-11)) < base.Ui32(int32(-10)) {
		v116 = v26
		goto L7
	} else {
		goto L20
	}
L20:
	;
	v61 = int32(0)
	v62 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
	if v62 < v61 {
		v116 = v61
		goto L7
	} else {
		goto L21
	}
L21:
	;
	v65 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
	if int32(100) < v65 {
		v116 = v61
		goto L7
	} else {
		goto L22
	}
L22:
	;
	if v65 < v62 {
		v116 = v61
		goto L7
	} else {
		goto L23
	}
L23:
	;
	v69 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
	if base.Ui32(int32(1)) < base.Ui32(v69) {
		v116 = v61
		goto L7
	} else {
		goto L24
	}
L24:
	;
	v72 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	if base.Ui32(int32(7)) < base.Ui32(v72) {
		v116 = v61
		goto L7
	} else {
		goto L25
	}
L25:
	;
	v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
	if base.Ui32(int32(3)) < base.Ui32(v75) {
		v116 = v61
		goto L7
	} else {
		goto L26
	}
L26:
	;
	v78 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
	if base.Ui32(int32(100)) < base.Ui32(v78) {
		v116 = v61
		goto L7
	} else {
		goto L27
	}
L27:
	;
	v81 = int32(0)
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
	if v82 < v81 {
		v116 = v81
		goto L7
	} else {
		goto L28
	}
L28:
	;
	v85 = int32(0)
	v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v86 < v85 {
		v116 = v85
		goto L7
	} else {
		goto L29
	}
L29:
	;
	v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	if base.Ui32(int32(100)) < base.Ui32(v89) {
		v116 = v85
		goto L7
	} else {
		goto L30
	}
L30:
	;
	v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if base.Ui32(int32(1)) < base.Ui32(v92) {
		v116 = v85
		goto L7
	} else {
		goto L31
	}
L31:
	;
	v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
	if base.Ui32(int32(100)) < base.Ui32(v95) {
		v116 = v85
		goto L7
	} else {
		goto L32
	}
L32:
	;
	v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	if base.Ui32(int32(3)) < base.Ui32(v98) {
		v116 = v85
		goto L7
	} else {
		goto L33
	}
L33:
	;
	v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
	if base.Ui32(int32(1)) < base.Ui32(v101) {
		v116 = v85
		goto L7
	} else {
		goto L34
	}
L34:
	;
	v104 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
	if base.Ui32(int32(1)) < base.Ui32(v104) {
		v116 = v85
		goto L7
	} else {
		goto L35
	}
L35:
	;
	v107 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
	if base.Ui32(int32(1)) < base.Ui32(v107) {
		v116 = v85
		goto L7
	} else {
		goto L36
	}
L36:
	;
	v110 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	if base.Ui32(int32(1)) < base.Ui32(v110) {
		v116 = v85
		goto L7
	} else {
		goto L37
	}
L37:
	;
	v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
	v116 = base.B2i32(base.Ui32(v113) < base.Ui32(int32(2)))
	goto L7
L38:
	;
	v121 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v121 != 0 {
		v450 = int32(0)
		goto L1
	} else {
		goto L39
	}
L39:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(4)
	return int32(0)
L40:
	;
	v150 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	if int32(_a_F_WebPEncode_0) < v150 {
		goto L50
	} else {
		goto L51
	}
L41:
	;
	if v147 != 0 {
		goto L40
	} else {
		goto L48
	}
L42:
	;
	v130 = int32(5)
	v131 = *(*int32)(unsafe.Add(mBase, uint32(l1)+8))
	if v131 < int32(1) {
		v140 = v130
		goto L45
	} else {
		goto L46
	}
L43:
	;
	v147 = int32(0)
	goto L41
L44:
	;
	v147 = v144
	goto L41
L45:
	;
	v142 = F_WebPEncodingSetError(m, l1, v140)
	mBase = m.M
	v144 = v142
	goto L44
L46:
	;
	v134 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	if v134 < int32(1) {
		v140 = v130
		goto L45
	} else {
		goto L47
	}
L47:
	;
	v139 = *(*int32)(unsafe.Add(mBase, uint32(l1)+4))
	switch v139 {
	case 0, 4:
		v144 = int32(1)
		goto L44
	default:
		v140 = int32(4)
		goto L45
	}
L48:
	;
	return int32(0)
L49:
	;
	v161 = *(*int32)(unsafe.Add(mBase, uint32(l1)+88))
	if v161 == int32(0) {
		goto L54
	} else {
		goto L55
	}
L50:
	;
	v156 = *(*int32)(unsafe.Add(mBase, uint32(l1)+92))
	if v156 != 0 {
		v450 = int32(0)
		goto L1
	} else {
		goto L53
	}
L51:
	;
	v153 = *(*int32)(unsafe.Add(mBase, uint32(l1)+12))
	if v153 < int32(_a_F_WebPEncode_1) {
		goto L49
	} else {
		goto L52
	}
L52:
	;
	goto L50
L53:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+92)) = int32(5)
	return int32(0)
L54:
	;
	v286 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v286 != 0 {
		goto L70
	} else {
		goto L71
	}
L55:
	;
	goto L58
L56:
	;
	goto L54
L58:
	;
	base.MemoryFill(m, v161, int32(0), int32(188))
	goto L56
L70:
	;
	v440 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v440 != 0 {
		goto L123
	} else {
		goto L124
	}
L71:
	;
	v287 = *(*int32)(unsafe.Add(mBase, uint32(l1)))
	if v287 != 0 {
		goto L73
	} else {
		goto L74
	}
L72:
	;
	v374 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	if v374 != 0 {
		goto L98
	} else {
		goto L99
	}
L73:
	;
	v295 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
	if v295 != 0 {
		goto L79
	} else {
		goto L80
	}
L74:
	;
	v288 = *(*int32)(unsafe.Add(mBase, uint32(l1)+16))
	if v288 == int32(0) {
		goto L73
	} else {
		goto L75
	}
L75:
	;
	v291 = *(*int32)(unsafe.Add(mBase, uint32(l1)+20))
	if v291 == int32(0) {
		goto L73
	} else {
		goto L76
	}
L76:
	;
	v294 = *(*int32)(unsafe.Add(mBase, uint32(l1)+24))
	if v294 != 0 {
		goto L72
	} else {
		goto L77
	}
L77:
	;
	goto L73
L78:
	;
	if v296&int32(2) == int32(0) {
		v340 = float32(0)
		goto L88
	} else {
		goto L89
	}
L79:
	;
	if l1 != 0 {
		goto L83
	} else {
		goto L84
	}
L80:
	;
	v296 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
	if v296&int32(4) == int32(0) {
		goto L78
	} else {
		goto L81
	}
L81:
	;
	goto L79
L82:
	;
	if v323 != 0 {
		goto L72
	} else {
		goto L87
	}
L83:
	;
	v304 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v304 != 0 {
		goto L85
	} else {
		goto L86
	}
L84:
	;
	v323 = int32(0)
	goto L82
L85:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = int32(0)
	v309 = int32(2)
	v311 = int32(1)
	v316 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
	v321 = F_ImportYUVAFromRGBA(m, v304+v309, v304+v311, v304, v304+int32(3), int32(4), v316<<(uint(v309)%32), float32(0), v311, l1)
	mBase = m.M
	v323 = v321
	goto L82
L86:
	;
	v306 = F_WebPEncodingSetError(m, l1, int32(3))
	mBase = m.M
	v323 = v306
	goto L82
L87:
	;
	return int32(0)
L88:
	;
	if l1 != 0 {
		goto L91
	} else {
		goto L92
	}
L89:
	;
	v331 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
	v333 = base.F32_div(v331, float32(100))
	v334 = base.F32_mul(v333, v333)
	v340 = base.F32_add(base.F32_mul(base.F32_mul(v334, float32(-0.5)), v334), float32(1))
	goto L88
L90:
	;
	if v369 == int32(0) {
		v450 = int32(0)
		goto L1
	} else {
		goto L97
	}
L91:
	;
	v345 = *(*int32)(unsafe.Add(mBase, uint32(l1)+52))
	if v345 != 0 {
		goto L93
	} else {
		goto L94
	}
L92:
	;
	v369 = int32(0)
	goto L90
L93:
	;
	goto L95
L94:
	;
	v347 = F_WebPEncodingSetError(m, l1, int32(3))
	mBase = m.M
	v369 = v347
	goto L90
L95:
	;
	v354 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l1)+4)) = v354
	v356 = int32(2)
	v363 = *(*int32)(unsafe.Add(mBase, uint32(l1)+56))
	v367 = F_ImportYUVAFromRGBA(m, v345+v356, v345+int32(1), v345, v345+int32(3), int32(4), v363<<(uint(v356)%32), v340, v354, l1)
	mBase = m.M
	v369 = v367
	goto L90
L97:
	;
	goto L72
L98:
	;
	v376 = F_InitVP8Encoder(m, l0, l1)
	mBase = m.M
	if v376 != 0 {
		goto L100
	} else {
		goto L101
	}
L99:
	;
	F_WebPCleanupTransparentArea(m, l1)
	mBase = m.M
	goto L98
L100:
	;
	v379 = F_VP8EncAnalyze(m, v376)
	mBase = m.M
	if v379 == int32(0) {
		goto L105
	} else {
		goto L106
	}
L101:
	;
	return int32(0)
L102:
	;
	v432 = F_VP8EncDeleteAlpha(m, v376)
	mBase = m.M
	F_VP8TBufferClear(m, v376+int32(344))
	mBase = m.M
	F_WebPSafeFree(m, v376)
	mBase = m.M
	goto L122
L103:
	;
	F_VP8BitWriterWipeOut(m, v376+int32(56))
	mBase = m.M
	v404 = *(*int32)(unsafe.Add(mBase, uint32(v376)+52))
	if v404 < int32(1) {
		goto L116
	} else {
		goto L117
	}
L104:
	;
	v396 = F_VP8EncWrite(m, v376)
	mBase = m.M
	F_StoreStats(m, v376)
	mBase = m.M
	if v396 != 0 {
		goto L102
	} else {
		goto L114
	}
L105:
	;
	F_StoreStats(m, v376)
	mBase = m.M
	goto L103
L106:
	;
	v382 = F_VP8EncStartAlpha(m, v376)
	mBase = m.M
	v383 = *(*int32)(unsafe.Add(mBase, uint32(v376)+uint32(_c_F_WebPEncode[0])))
	if v383 != 0 {
		goto L108
	} else {
		goto L109
	}
L107:
	;
	if v390 == int32(0) {
		goto L105
	} else {
		goto L112
	}
L108:
	;
	if v382 == int32(0) {
		goto L105
	} else {
		goto L111
	}
L109:
	;
	if v382 == int32(0) {
		goto L105
	} else {
		goto L110
	}
L110:
	;
	v386 = F_VP8EncLoop(m, v376)
	mBase = m.M
	v390 = v386
	goto L107
L111:
	;
	v389 = F_VP8EncTokenLoop(m, v376)
	mBase = m.M
	v390 = v389
	goto L107
L112:
	;
	v393 = F_VP8EncFinishAlpha(m, v376)
	mBase = m.M
	if v393 != 0 {
		goto L104
	} else {
		goto L113
	}
L113:
	;
	goto L105
L114:
	;
	goto L103
L115:
	;
	v424 = F_VP8EncDeleteAlpha(m, v376)
	mBase = m.M
	F_VP8TBufferClear(m, v376+int32(344))
	mBase = m.M
	F_WebPSafeFree(m, v376)
	mBase = m.M
	goto L121
L116:
	;
	goto L115
L117:
	;
	v411 = v376 + int32(88)
	v412 = int32(0)
	goto L118
L118:
	;
	F_VP8BitWriterWipeOut(m, v411)
	mBase = m.M
	v417 = v412 + int32(1)
	v418 = *(*int32)(unsafe.Add(mBase, uint32(v376)+52))
	if v417 < v418 {
		v411 = v411 + int32(32)
		v412 = v417
		goto L118
	} else {
		goto L120
	}
L119:
	;
	goto L116
L120:
	;
	goto L119
L121:
	;
	return int32(0)
L122:
	;
	return v432 & int32(1)
L123:
	;
	v446 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
	if v446 != 0 {
		goto L126
	} else {
		goto L127
	}
L124:
	;
	v441 = int32(0)
	v442 = F_WebPPictureYUVAToARGB(m, l1)
	mBase = m.M
	if v442 == v441 {
		v450 = v441
		goto L1
	} else {
		goto L125
	}
L125:
	;
	goto L123
L126:
	;
	v449 = F_VP8LEncodeImage(m, l0, l1)
	mBase = m.M
	v450 = v449
	goto L1
L127:
	;
	F_WebPReplaceTransparentPixels(m, l1, int32(0))
	mBase = m.M
	goto L126
}
func F_WebPEncodingSetError(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	v3 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
	if v3 != 0 {
	} else {
		*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = l1
	}
	return int32(0)
}
func F_WebPEstimateBestFilter(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32) int32 {
	mBase := m.M
	_ = mBase
	var v62 int32
	_ = v62
	var v63 int32
	_ = v63
	var v64 int32
	_ = v64
	var v188 int32
	_ = v188
	var v211 int32
	_ = v211
	var v234 int32
	_ = v234
	var v244 int32
	_ = v244
	var v245 int32
	_ = v245
	var v247 int32
	_ = v247
	var v264 int32
	_ = v264
	var v265 int32
	_ = v265
	var v270 int32
	_ = v270
	var v315 int32
	_ = v315
	var v321 int32
	_ = v321
	var v324 int32
	_ = v324
	var v378 int32
	_ = v378
	var v379 int32
	_ = v379
	var v381 int32
	_ = v381
	var v382 int32
	_ = v382
	var v383 int32
	_ = v383
	var v384 int32
	_ = v384
	var v387 int32
	_ = v387
	var v394 int32
	_ = v394
	var v395 int32
	_ = v395
	var v397 int32
	_ = v397
	var v402 int32
	_ = v402
	var v407 int32
	_ = v407
	var v410 int32
	_ = v410
	var v411 int32
	_ = v411
	var v413 int32
	_ = v413
	var v424 int32
	_ = v424
	var v425 int32
	_ = v425
	var v426 int32
	_ = v426
	var v429 int32
	_ = v429
	var v430 int32
	_ = v430
	var v433 int32
	_ = v433
	var v434 int32
	_ = v434
	var v436 int32
	_ = v436
	var v439 int32
	_ = v439
	var v461 int32
	_ = v461
	var v464 int32
	_ = v464
	var v465 int32
	_ = v465
	var v468 int32
	_ = v468
	var v470 int32
	_ = v470
	var v471 int32
	_ = v471
	var v474 int32
	_ = v474
	var v476 int32
	_ = v476
	var v477 int32
	_ = v477
	var v480 int32
	_ = v480
	var v482 int32
	_ = v482
	var v483 int32
	_ = v483
	var v486 int32
	_ = v486
	var v488 int32
	_ = v488
	var v489 int32
	_ = v489
	var v492 int32
	_ = v492
	var v494 int32
	_ = v494
	var v495 int32
	_ = v495
	var v498 int32
	_ = v498
	var v500 int32
	_ = v500
	var v501 int32
	_ = v501
	var v504 int32
	_ = v504
	var v506 int32
	_ = v506
	var v507 int32
	_ = v507
	var v510 int32
	_ = v510
	var v512 int32
	_ = v512
	var v513 int32
	_ = v513
	var v516 int32
	_ = v516
	var v518 int32
	_ = v518
	var v519 int32
	_ = v519
	var v522 int32
	_ = v522
	var v524 int32
	_ = v524
	var v525 int32
	_ = v525
	var v528 int32
	_ = v528
	var v530 int32
	_ = v530
	var v531 int32
	_ = v531
	var v534 int32
	_ = v534
	var v536 int32
	_ = v536
	var v537 int32
	_ = v537
	var v540 int32
	_ = v540
	var v542 int32
	_ = v542
	var v543 int32
	_ = v543
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v550 int32
	_ = v550
	var v553 int32
	_ = v553
	var v556 int32
	_ = v556
	var v558 int32
	_ = v558
	var v561 int32
	_ = v561
	var v563 int32
	_ = v563
	var v568 int32
	_ = v568
	var v573 int32
	_ = v573
	var v582 int32
	_ = v582
	var v584 int32
	_ = v584
	var v585 int32
	_ = v585
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v588 int32
	_ = v588
	var v589 int32
	_ = v589
	var v590 int32
	_ = v590
	var v591 int32
	_ = v591
	var v592 int32
	_ = v592
	var v593 int32
	_ = v593
	var v594 int32
	_ = v594
	var v595 int32
	_ = v595
	var v596 int32
	_ = v596
	var v597 int32
	_ = v597
	var v598 int32
	_ = v598
	var v599 int32
	_ = v599
	var v600 int32
	_ = v600
	var v601 int32
	_ = v601
	var v602 int32
	_ = v602
	var v603 int32
	_ = v603
	var v639 int32
	_ = v639
	var v640 int32
	_ = v640
	var v641 int32
	_ = v641
	var v642 int32
	_ = v642
	var v643 int32
	_ = v643
	var v644 int32
	_ = v644
	var v645 int32
	_ = v645
	var v646 int32
	_ = v646
	var v647 int32
	_ = v647
	var v648 int32
	_ = v648
	var v649 int32
	_ = v649
	var v650 int32
	_ = v650
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v653 int32
	_ = v653
	var v654 int32
	_ = v654
	var v655 int32
	_ = v655
	var v656 int32
	_ = v656
	var v657 int32
	_ = v657
	var v658 int32
	_ = v658
	var v659 int32
	_ = v659
	var v660 int32
	_ = v660
	var v661 int32
	_ = v661
	var v662 int32
	_ = v662
	var v663 int32
	_ = v663
	var v664 int32
	_ = v664
	var v665 int32
	_ = v665
	var v666 int32
	_ = v666
	var v667 int32
	_ = v667
	var v668 int32
	_ = v668
	var v669 int32
	_ = v669
	var v670 int32
	_ = v670
	var v671 int32
	_ = v671
	var v672 int32
	_ = v672
	var v673 int32
	_ = v673
	var v674 int32
	_ = v674
	var v675 int32
	_ = v675
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v684 int32
	_ = v684
	var v687 int32
	_ = v687
	var v689 int32
	_ = v689
	var v692 int32
	_ = v692
	var v694 int32
	_ = v694
	var v697 int32
	_ = v697
	var v699 int32
	_ = v699
	var v702 int32
	_ = v702
	var v704 int32
	_ = v704
	var v707 int32
	_ = v707
	var v709 int32
	_ = v709
	var v712 int32
	_ = v712
	var v714 int32
	_ = v714
	var v717 int32
	_ = v717
	var v718 int32
	_ = v718
	var v726 int32
	_ = v726
	var v740 int32
	_ = v740
	var v754 int32
	_ = v754
	var v755 int32
	_ = v755
	var v757 int32
	_ = v757
	var v760 int32
	_ = v760
	var v762 int32
	_ = v762
	var v765 int32
	_ = v765
	var v767 int32
	_ = v767
	var v770 int32
	_ = v770
	var v772 int32
	_ = v772
	var v775 int32
	_ = v775
	var v777 int32
	_ = v777
	var v780 int32
	_ = v780
	var v782 int32
	_ = v782
	var v785 int32
	_ = v785
	var v787 int32
	_ = v787
	var v790 int32
	_ = v790
	var v791 int32
	_ = v791
	var v799 int32
	_ = v799
	var v801 int32
	_ = v801
	var v804 int32
	_ = v804
	var v806 int32
	_ = v806
	var v809 int32
	_ = v809
	var v810 int32
	_ = v810
	var v818 int32
	_ = v818
	var v819 int32
	_ = v819
	var v838 int32
	_ = v838
	var v839 int32
	_ = v839
	var v840 int32
	_ = v840
	var v841 int32
	_ = v841
	var v843 int32
	_ = v843
	var v846 int32
	_ = v846
	var v848 int32
	_ = v848
	var v851 int32
	_ = v851
	var v853 int32
	_ = v853
	var v856 int32
	_ = v856
	var v858 int32
	_ = v858
	var v861 int32
	_ = v861
	var v863 int32
	_ = v863
	var v866 int32
	_ = v866
	var v868 int32
	_ = v868
	var v871 int32
	_ = v871
	var v873 int32
	_ = v873
	var v876 int32
	_ = v876
	var v877 int32
	_ = v877
	var v885 int32
	_ = v885
	var v887 int32
	_ = v887
	var v890 int32
	_ = v890
	var v892 int32
	_ = v892
	var v895 int32
	_ = v895
	var v896 int32
	_ = v896
	var v904 int32
	_ = v904
	var v905 int32
	_ = v905
	var v925 int32
	_ = v925
	var v927 int32
	_ = v927
	v62 = m.G0
	v63 = int32(256)
	v64 = v62 - v63
	m.G0 = v64
	base.MemoryFill(m, v64, int32(0), v63)
	v188 = int32(0)
	if l2 < int32(4) {
		v582 = v188
		v584 = v188
		v585 = v188
		v586 = v188
		v587 = v188
		v588 = v188
		v589 = v188
		v590 = v188
		v591 = v188
		v592 = v188
		v593 = v188
		v594 = v188
		v595 = v188
		v596 = v188
		v597 = v188
		v598 = v188
		v599 = v188
		v600 = v188
		v601 = v188
		v602 = v188
		v603 = v188
	} else {
		v211 = int32(0)
		if l1 < int32(4) {
			v582 = v211
			v584 = v211
			v585 = v211
			v586 = v211
			v587 = v211
			v588 = v211
			v589 = v211
			v590 = v211
			v591 = v211
			v592 = v211
			v593 = v211
			v594 = v211
			v595 = v211
			v596 = v211
			v597 = v211
			v598 = v211
			v599 = v211
			v600 = v211
			v601 = v211
			v602 = v211
			v603 = v211
		} else {
			v234 = int32(-1)
			v244 = int32(1)
			v245 = l3 << (uint(v244) % 32)
			v247 = v245 | v244
			v264 = l0 + v247
			v265 = l0 + (v247 - l1)
			v270 = int32(2)
			for {
				v315 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0+v270*l3))))
				v321 = int32(0)
				v324 = v315
				for {
					v378 = v264 + v321
					v379 = int32(1)
					v381 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v378+v379))))
					v382 = v381 - v324
					v383 = int32(31)
					v384 = v382 >> (uint(v383) % 32)
					v387 = int32(2)
					*(*int32)(unsafe.Add(mBase, uint32(v64+int32(base.Ui32(v382^v384-v384)>>(uint(v387)%32))&int32(1073741820)))) = v379
					v394 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v378))))
					v395 = v381 - v394
					v397 = v395 >> (uint(v383) % 32)
					v402 = int32(252)
					*(*int32)(unsafe.Add(mBase, uint32(v64+int32(64)+int32(base.Ui32(v395^v397-v397)>>(uint(v387)%32))&v402))) = v379
					v407 = v265 + v321
					v410 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v407+v379))))
					v411 = v381 - v410
					v413 = v411 >> (uint(v383) % 32)
					*(*int32)(unsafe.Add(mBase, uint32(v64+int32(128)+int32(base.Ui32(v411^v413-v413)>>(uint(v387)%32))&v402))) = v379
					v424 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v407))))
					v425 = v410 + v394 - v424
					v426 = int32(0)
					if v426 < v425 {
						v429 = v425
					} else {
						v429 = v426
					}
					v430 = int32(255)
					if v429 < v430 {
						v433 = v429
					} else {
						v433 = v430
					}
					v434 = v381 - v433
					v436 = v434 >> (uint(int32(31)) % 32)
					v439 = int32(2)
					*(*int32)(unsafe.Add(mBase, uint32(v64+int32(192)+int32(base.Ui32(v434^v436-v436)>>(uint(v439)%32))&int32(252)))) = int32(1)
					if v321+int32(4) < l1+v234 {
						v321 = v321 + v439
						v324 = int32(base.Ui32(v381+v324*int32(3)+v439) >> (uint(v439) % 32))
						continue
					} else {
						break
					}
					break
				}
				v461 = v270 + int32(2)
				if v461 < l2+v234 {
					v264 = v264 + v245
					v265 = v265 + v245
					v270 = v461
					continue
				} else {
					break
				}
				break
			}
			v464 = int32(0)
			v465 = *(*int32)(unsafe.Add(mBase, uint32(v64)+88))
			if v464 < v465 {
				v468 = int32(6)
			} else {
				v468 = v464
			}
			v470 = int32(0)
			v471 = *(*int32)(unsafe.Add(mBase, uint32(v64)+84))
			if v470 < v471 {
				v474 = int32(5)
			} else {
				v474 = v470
			}
			v476 = int32(0)
			v477 = *(*int32)(unsafe.Add(mBase, uint32(v64)+76))
			if v476 < v477 {
				v480 = int32(3)
			} else {
				v480 = v476
			}
			v482 = int32(0)
			v483 = *(*int32)(unsafe.Add(mBase, uint32(v64)+60))
			if v482 < v483 {
				v486 = int32(15)
			} else {
				v486 = v482
			}
			v488 = int32(0)
			v489 = *(*int32)(unsafe.Add(mBase, uint32(v64)+56))
			if v488 < v489 {
				v492 = int32(14)
			} else {
				v492 = v488
			}
			v494 = int32(0)
			v495 = *(*int32)(unsafe.Add(mBase, uint32(v64)+52))
			if v494 < v495 {
				v498 = int32(13)
			} else {
				v498 = v494
			}
			v500 = int32(0)
			v501 = *(*int32)(unsafe.Add(mBase, uint32(v64)+48))
			if v500 < v501 {
				v504 = int32(12)
			} else {
				v504 = v500
			}
			v506 = int32(0)
			v507 = *(*int32)(unsafe.Add(mBase, uint32(v64)+44))
			if v506 < v507 {
				v510 = int32(11)
			} else {
				v510 = v506
			}
			v512 = int32(0)
			v513 = *(*int32)(unsafe.Add(mBase, uint32(v64)+40))
			if v512 < v513 {
				v516 = int32(10)
			} else {
				v516 = v512
			}
			v518 = int32(0)
			v519 = *(*int32)(unsafe.Add(mBase, uint32(v64)+36))
			if v518 < v519 {
				v522 = int32(9)
			} else {
				v522 = v518
			}
			v524 = int32(0)
			v525 = *(*int32)(unsafe.Add(mBase, uint32(v64)+28))
			if v524 < v525 {
				v528 = int32(7)
			} else {
				v528 = v524
			}
			v530 = int32(0)
			v531 = *(*int32)(unsafe.Add(mBase, uint32(v64)+24))
			if v530 < v531 {
				v534 = int32(6)
			} else {
				v534 = v530
			}
			v536 = int32(0)
			v537 = *(*int32)(unsafe.Add(mBase, uint32(v64)+20))
			if v536 < v537 {
				v540 = int32(5)
			} else {
				v540 = v536
			}
			v542 = int32(0)
			v543 = *(*int32)(unsafe.Add(mBase, uint32(v64)+12))
			if v542 < v543 {
				v546 = int32(3)
			} else {
				v546 = v542
			}
			v547 = *(*int32)(unsafe.Add(mBase, uint32(v64)+68))
			v548 = int32(0)
			v550 = *(*int32)(unsafe.Add(mBase, uint32(v64)+4))
			v553 = *(*int32)(unsafe.Add(mBase, uint32(v64)+80))
			v556 = int32(2)
			v558 = *(*int32)(unsafe.Add(mBase, uint32(v64)+72))
			v561 = int32(1)
			v563 = *(*int32)(unsafe.Add(mBase, uint32(v64)+32))
			v568 = *(*int32)(unsafe.Add(mBase, uint32(v64)+16))
			v573 = *(*int32)(unsafe.Add(mBase, uint32(v64)+8))
			v582 = v468
			v584 = v474
			v585 = base.B2i32(v548 < v553) << (uint(v556) % 32)
			v586 = v480
			v587 = base.B2i32(v548 < v558) << (uint(v561) % 32)
			v588 = base.B2i32(v548 < v547)
			v589 = v486
			v590 = v492
			v591 = v498
			v592 = v504
			v593 = v510
			v594 = v516
			v595 = v522
			v596 = base.B2i32(v548 < v563) << (uint(int32(3)) % 32)
			v597 = v528
			v598 = v534
			v599 = v540
			v600 = base.B2i32(v548 < v568) << (uint(v556) % 32)
			v601 = v546
			v602 = base.B2i32(v548 < v573) << (uint(v561) % 32)
			v603 = base.B2i32(v548 < v550)
		}
	}
	v639 = *(*int32)(unsafe.Add(mBase, uint32(v64)+92))
	v640 = *(*int32)(unsafe.Add(mBase, uint32(v64)+96))
	v641 = *(*int32)(unsafe.Add(mBase, uint32(v64)+100))
	v642 = *(*int32)(unsafe.Add(mBase, uint32(v64)+104))
	v643 = *(*int32)(unsafe.Add(mBase, uint32(v64)+108))
	v644 = *(*int32)(unsafe.Add(mBase, uint32(v64)+112))
	v645 = *(*int32)(unsafe.Add(mBase, uint32(v64)+116))
	v646 = *(*int32)(unsafe.Add(mBase, uint32(v64)+120))
	v647 = *(*int32)(unsafe.Add(mBase, uint32(v64)+124))
	v648 = *(*int32)(unsafe.Add(mBase, uint32(v64)+132))
	v649 = *(*int32)(unsafe.Add(mBase, uint32(v64)+136))
	v650 = *(*int32)(unsafe.Add(mBase, uint32(v64)+140))
	v651 = *(*int32)(unsafe.Add(mBase, uint32(v64)+144))
	v652 = *(*int32)(unsafe.Add(mBase, uint32(v64)+148))
	v653 = *(*int32)(unsafe.Add(mBase, uint32(v64)+152))
	v654 = *(*int32)(unsafe.Add(mBase, uint32(v64)+156))
	v655 = *(*int32)(unsafe.Add(mBase, uint32(v64)+160))
	v656 = *(*int32)(unsafe.Add(mBase, uint32(v64)+164))
	v657 = *(*int32)(unsafe.Add(mBase, uint32(v64)+168))
	v658 = *(*int32)(unsafe.Add(mBase, uint32(v64)+172))
	v659 = *(*int32)(unsafe.Add(mBase, uint32(v64)+176))
	v660 = *(*int32)(unsafe.Add(mBase, uint32(v64)+180))
	v661 = *(*int32)(unsafe.Add(mBase, uint32(v64)+184))
	v662 = *(*int32)(unsafe.Add(mBase, uint32(v64)+188))
	v663 = *(*int32)(unsafe.Add(mBase, uint32(v64)+196))
	v664 = *(*int32)(unsafe.Add(mBase, uint32(v64)+200))
	v665 = *(*int32)(unsafe.Add(mBase, uint32(v64)+204))
	v666 = *(*int32)(unsafe.Add(mBase, uint32(v64)+208))
	v667 = *(*int32)(unsafe.Add(mBase, uint32(v64)+212))
	v668 = *(*int32)(unsafe.Add(mBase, uint32(v64)+216))
	v669 = *(*int32)(unsafe.Add(mBase, uint32(v64)+220))
	v670 = *(*int32)(unsafe.Add(mBase, uint32(v64)+224))
	v671 = *(*int32)(unsafe.Add(mBase, uint32(v64)+228))
	v672 = *(*int32)(unsafe.Add(mBase, uint32(v64)+232))
	v673 = *(*int32)(unsafe.Add(mBase, uint32(v64)+236))
	v674 = *(*int32)(unsafe.Add(mBase, uint32(v64)+240))
	v675 = *(*int32)(unsafe.Add(mBase, uint32(v64)+244))
	v676 = *(*int32)(unsafe.Add(mBase, uint32(v64)+248))
	v677 = *(*int32)(unsafe.Add(mBase, uint32(v64)+252))
	m.G0 = v64 + int32(256)
	v684 = int32(0)
	if v684 < v647 {
		v687 = int32(15)
	} else {
		v687 = v684
	}
	v689 = int32(0)
	if v689 < v646 {
		v692 = int32(14)
	} else {
		v692 = v689
	}
	v694 = int32(0)
	if v694 < v645 {
		v697 = int32(13)
	} else {
		v697 = v694
	}
	v699 = int32(0)
	if v699 < v644 {
		v702 = int32(12)
	} else {
		v702 = v699
	}
	v704 = int32(0)
	if v704 < v643 {
		v707 = int32(11)
	} else {
		v707 = v704
	}
	v709 = int32(0)
	if v709 < v642 {
		v712 = int32(10)
	} else {
		v712 = v709
	}
	v714 = int32(0)
	if v714 < v641 {
		v717 = int32(9)
	} else {
		v717 = v714
	}
	v718 = int32(0)
	if v718 < v639 {
		v726 = int32(7)
	} else {
		v726 = v718
	}
	v740 = v687 + (v692 + (v697 + (v702 + (v707 + (v712 + (v717 + (base.B2i32(v718 < v640)<<(uint(int32(3))%32) + (v726 + (v582 + (v584 + (v585 + (v586 + (v587 | v588)))))))))))))
	v754 = v589 + (v590 + (v591 + (v592 + (v593 + (v594 + (v595 + (v596 + (v597 + (v598 + (v599 + (v600 + (v601 + (v602 | v603)))))))))))))
	v755 = base.B2i32(base.Ui32(v740) < base.Ui32(v754))
	v757 = int32(0)
	if v757 < v662 {
		v760 = int32(15)
	} else {
		v760 = v757
	}
	v762 = int32(0)
	if v762 < v661 {
		v765 = int32(14)
	} else {
		v765 = v762
	}
	v767 = int32(0)
	if v767 < v660 {
		v770 = int32(13)
	} else {
		v770 = v767
	}
	v772 = int32(0)
	if v772 < v659 {
		v775 = int32(12)
	} else {
		v775 = v772
	}
	v777 = int32(0)
	if v777 < v658 {
		v780 = int32(11)
	} else {
		v780 = v777
	}
	v782 = int32(0)
	if v782 < v657 {
		v785 = int32(10)
	} else {
		v785 = v782
	}
	v787 = int32(0)
	if v787 < v656 {
		v790 = int32(9)
	} else {
		v790 = v787
	}
	v791 = int32(0)
	if v791 < v654 {
		v799 = int32(7)
	} else {
		v799 = v791
	}
	v801 = int32(0)
	if v801 < v653 {
		v804 = int32(6)
	} else {
		v804 = v801
	}
	v806 = int32(0)
	if v806 < v652 {
		v809 = int32(5)
	} else {
		v809 = v806
	}
	v810 = int32(0)
	if v810 < v650 {
		v818 = int32(3)
	} else {
		v818 = v810
	}
	v819 = int32(0)
	v838 = v760 + (v765 + (v770 + (v775 + (v780 + (v785 + (v790 + (base.B2i32(v791 < v655)<<(uint(int32(3))%32) + (v799 + (v804 + (v809 + (base.B2i32(v810 < v651)<<(uint(int32(2))%32) + (v818 + (base.B2i32(v819 < v649)<<(uint(int32(1))%32) | base.B2i32(v819 < v648))))))))))))))
	if base.Ui32(v740) < base.Ui32(v754) {
		v839 = v740
	} else {
		v839 = v754
	}
	v840 = base.B2i32(base.Ui32(v838) < base.Ui32(v839))
	if base.Ui32(v838) < base.Ui32(v839) {
		v841 = int32(2)
	} else {
		v841 = v755
	}
	v843 = int32(0)
	if v843 < v677 {
		v846 = int32(15)
	} else {
		v846 = v843
	}
	v848 = int32(0)
	if v848 < v676 {
		v851 = int32(14)
	} else {
		v851 = v848
	}
	v853 = int32(0)
	if v853 < v675 {
		v856 = int32(13)
	} else {
		v856 = v853
	}
	v858 = int32(0)
	if v858 < v674 {
		v861 = int32(12)
	} else {
		v861 = v858
	}
	v863 = int32(0)
	if v863 < v673 {
		v866 = int32(11)
	} else {
		v866 = v863
	}
	v868 = int32(0)
	if v868 < v672 {
		v871 = int32(10)
	} else {
		v871 = v868
	}
	v873 = int32(0)
	if v873 < v671 {
		v876 = int32(9)
	} else {
		v876 = v873
	}
	v877 = int32(0)
	if v877 < v669 {
		v885 = int32(7)
	} else {
		v885 = v877
	}
	v887 = int32(0)
	if v887 < v668 {
		v890 = int32(6)
	} else {
		v890 = v887
	}
	v892 = int32(0)
	if v892 < v667 {
		v895 = int32(5)
	} else {
		v895 = v892
	}
	v896 = int32(0)
	if v896 < v665 {
		v904 = int32(3)
	} else {
		v904 = v896
	}
	v905 = int32(0)
	if base.Ui32(v838) < base.Ui32(v839) {
		v925 = v838
	} else {
		v925 = v839
	}
	if base.Ui32(v846+(v851+(v856+(v861+(v866+(v871+(v876+(base.B2i32(v877 < v670)<<(uint(int32(3))%32)+(v885+(v890+(v895+(base.B2i32(v896 < v666)<<(uint(int32(2))%32)+(v904+(base.B2i32(v905 < v664)<<(uint(int32(1))%32)|base.B2i32(v905 < v663))))))))))))))) < base.Ui32(v925) {
		v927 = int32(3)
	} else {
		v927 = v841
	}
	return v927
}
func F_WebPGetLinePairConverter(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
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
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v28 int32
	_ = v28
	var v31 int32
	_ = v31
	var v52 int32
	_ = v52
	var v57 int32
	_ = v57
	var v59 int32
	_ = v59
	v8 = m.G1
	v11 = *(*int32)(unsafe.Add(mBase, uint32(v8)+uint32(_c_F_WebPGetLinePairConverter[0])))
	v12 = m.G3
	v13 = *(*int32)(unsafe.Add(mBase, uint32(v12)))
	if v11 == v13 {
	} else {
		v15 = m.G2
		v16 = m.G1
		v20 = v15 + int32(14)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[1]))) = v20
		v23 = v15 + int32(15)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[2]))) = v23
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[3]))) = v20
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[4]))) = v23
		v28 = v15 + int32(16)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[5]))) = v28
		v31 = v15 + int32(17)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[6]))) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[7]))) = v15 + int32(18)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[8]))) = v28
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[9]))) = v31
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[10]))) = v15 + int32(19)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[11]))) = v15 + int32(20)
		*(*int32)(unsafe.Add(mBase, uint32(v16)+uint32(_c_F_WebPGetLinePairConverter[0]))) = v13
	}
	v52 = m.G1
	if l0 != 0 {
		v57 = int32(12)
	} else {
		v57 = int32(16)
	}
	v59 = *(*int32)(unsafe.Add(mBase, uint32(v52+int32(_a_F_WebPGetLinePairConverter_0)+v57)))
	return v59
}
func F_WebPGetWorkerInterface(m *base.Module) int32 {
	var v1 int32
	_ = v1
	v1 = m.G1
	return v1 + int32(_a_F_WebPGetWorkerInterface_0)
}
func F_WebPInitAlphaProcessing(m *base.Module) {
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
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_WebPInitAlphaProcessing[0])))
	v8 = m.G3
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[1]))) = v11 + int32(2)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[2]))) = v11 + int32(3)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[3]))) = v11 + int32(4)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[4]))) = v11 + int32(5)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[5]))) = v11 + int32(6)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[6]))) = v11 + int32(7)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[7]))) = v11 + int32(8)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[8]))) = v11 + int32(9)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[9]))) = v11 + int32(10)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[10]))) = v11 + int32(11)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[11]))) = v11 + int32(12)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[12]))) = v11 + int32(13)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitAlphaProcessing[0]))) = v9
	}
	return
}
func F_WebPInitConvertARGBToYUV(m *base.Module) {
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
	v7 = *(*int32)(unsafe.Add(mBase, uint32(v4)+uint32(_c_F_WebPInitConvertARGBToYUV[0])))
	v8 = m.G3
	v9 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
	if v7 == v9 {
	} else {
		v11 = m.G2
		v12 = m.G1
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[1]))) = v11 + int32(58)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[2]))) = v11 + int32(59)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[3]))) = v11 + int32(60)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[4]))) = v11 + int32(61)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[5]))) = v11 + int32(62)
		*(*int32)(unsafe.Add(mBase, uint32(v12)+uint32(_c_F_WebPInitConvertARGBToYUV[0]))) = v9
	}
	return
}
func F_WebPMemoryWrite(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v11 int64
	_ = v11
	var v13 int64
	_ = v13
	var v14 int64
	_ = v14
	var v17 int64
	_ = v17
	var v19 int64
	_ = v19
	var v20 int64
	_ = v20
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
	var v49 int32
	_ = v49
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v54 int32
	_ = v54
	var v64 int32
	_ = v64
	var v65 int32
	_ = v65
	var v67 int32
	_ = v67
	var v68 int32
	_ = v68
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l2)+76))
	if v8 == int32(0) {
		return int32(1)
	} else {
		v11 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v8)+4)))
		v13 = v11 + base.I64_extend_i32_u(l1)
		v14 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v8)+8)))
		if base.Ui64(v13) <= base.Ui64(v14) {
			if l1 == int32(0) {
			} else {
				v64 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
				v65 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				v67 = F_memcpy(m, v64+v65, l0, l1)
				mBase = m.M
				v68 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v68 + l1
			}
			return int32(1)
		} else {
			v17 = v14 << (uint(int64(1)) % 64)
			if base.Ui64(v13) < base.Ui64(v17) {
				v19 = v17
			} else {
				v19 = v13
			}
			v20 = int64(8192)
			if base.Ui64(v20) < base.Ui64(v19) {
				v23 = v19
			} else {
				v23 = v20
			}
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
			if v45 != 0 {
				v49 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
				if v49 == int32(0) {
				} else {
					v52 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
					v53 = F_memcpy(m, v45, v52, v49)
					mBase = m.M
				}
				v54 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
				F_free(m, v54)
				mBase = m.M
				*(*uint32)(unsafe.Add(mBase, uint32(v8)+8)) = uint32(v23)
				*(*int32)(unsafe.Add(mBase, uint32(v8))) = v45
				if l1 == int32(0) {
				} else {
					v64 = *(*int32)(unsafe.Add(mBase, uint32(v8)))
					v65 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
					v67 = F_memcpy(m, v64+v65, l0, l1)
					mBase = m.M
					v68 = *(*int32)(unsafe.Add(mBase, uint32(v8)+4))
					*(*int32)(unsafe.Add(mBase, uint32(v8)+4)) = v68 + l1
				}
				return int32(1)
			} else {
				return int32(0)
			}
		}
	}
}
func F_WebPMemoryWriterClear(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	if l0 == int32(0) {
	} else {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		F_free(m, v4)
		mBase = m.M
		*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = int32(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(0)
	}
	return
}
func F_WebPMemoryWriterInit(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	*(*int32)(unsafe.Add(mBase, uint32(l0)+8)) = int32(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0))) = int64(0)
	return
}
func F_WebPPictureARGBToYUVADithered(m *base.Module, l0 int32, l1 int32, l2 float32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v9 int32
	_ = v9
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v31 int32
	_ = v31
	var v35 int32
	_ = v35
	if l0 != 0 {
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
		if v7 != 0 {
			if l1&int32(3) == int32(0) {
				v22 = int32(0)
				*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = v22
				v24 = int32(2)
				v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
				v35 = F_ImportYUVAFromRGBA(m, v7+v24, v7+int32(1), v7, v7+int32(3), int32(4), v31<<(uint(v24)%32), l2, v22, l0)
				mBase = m.M
				return v35
			} else {
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v18 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(4)
				}
				return int32(0)
			}
		} else {
			v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v9 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(3)
			}
			return int32(0)
		}
	} else {
		return int32(0)
	}
}
func F_WebPPictureAlloc(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int64
	_ = v8
	var v26 int32
	_ = v26
	var v37 int32
	_ = v37
	var v38 int32
	_ = v38
	var v39 int32
	_ = v39
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v55 int64
	_ = v55
	var v57 int32
	_ = v57
	var v71 int64
	_ = v71
	var v72 int64
	_ = v72
	var v75 int64
	_ = v75
	var v76 int32
	_ = v76
	var v80 int64
	_ = v80
	var v81 int64
	_ = v81
	var v82 int64
	_ = v82
	var v91 int32
	_ = v91
	var v96 int32
	_ = v96
	var v98 int64
	_ = v98
	var v99 int64
	_ = v99
	var v101 int64
	_ = v101
	var v104 int64
	_ = v104
	var v109 int32
	_ = v109
	var v111 int32
	_ = v111
	var v118 int32
	_ = v118
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v137 int32
	_ = v137
	var v142 int32
	_ = v142
	var v143 int32
	_ = v143
	var v144 int32
	_ = v144
	var v150 int32
	_ = v150
	var v151 int32
	_ = v151
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v167 int32
	_ = v167
	var v169 int32
	_ = v169
	var v180 int32
	_ = v180
	if l0 != 0 {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
		F_free(m, v4)
		mBase = m.M
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
		F_free(m, v6)
		mBase = m.M
		v8 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v8
		*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = int32(0)
		v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v26 != 0 {
			v142 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v143 = int32(5)
			v144 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			if v144 < int32(1) {
				v151 = v143
				v152 = F_WebPEncodingSetError(m, l0, v151)
				mBase = m.M
				if v152 != 0 {
					v155 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
					F_WebPSafeFree(m, v155)
					mBase = m.M
					*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
					v167 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v142)*base.I64_extend_i32_s(v144)+int64(31), int32(4))
					mBase = m.M
					if v167 != 0 {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v144
						*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v167
						*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v167 + int32(31)) & int32(-32)
						v180 = int32(1)
					} else {
						v169 = F_WebPEncodingSetError(m, l0, int32(1))
						mBase = m.M
						v180 = v169
					}
				} else {
					v180 = int32(0)
				}
			} else {
				if v142 < int32(1) {
					v151 = v143
					v152 = F_WebPEncodingSetError(m, l0, v151)
					mBase = m.M
					if v152 != 0 {
						v155 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
						F_WebPSafeFree(m, v155)
						mBase = m.M
						*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
						v167 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v142)*base.I64_extend_i32_s(v144)+int64(31), int32(4))
						mBase = m.M
						if v167 != 0 {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v144
							*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v167
							*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v167 + int32(31)) & int32(-32)
							v180 = int32(1)
						} else {
							v169 = F_WebPEncodingSetError(m, l0, int32(1))
							mBase = m.M
							v180 = v169
						}
					} else {
						v180 = int32(0)
					}
				} else {
					v150 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					switch v150 {
					case 0, 4:
						v155 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
						F_WebPSafeFree(m, v155)
						mBase = m.M
						*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
						v167 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v142)*base.I64_extend_i32_s(v144)+int64(31), int32(4))
						mBase = m.M
						if v167 != 0 {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v144
							*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v167
							*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v167 + int32(31)) & int32(-32)
							v180 = int32(1)
						} else {
							v169 = F_WebPEncodingSetError(m, l0, int32(1))
							mBase = m.M
							v180 = v169
						}
					default:
						v151 = int32(4)
						v152 = F_WebPEncodingSetError(m, l0, v151)
						mBase = m.M
						if v152 != 0 {
							v155 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
							F_WebPSafeFree(m, v155)
							mBase = m.M
							*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
							*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
							v167 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v142)*base.I64_extend_i32_s(v144)+int64(31), int32(4))
							mBase = m.M
							if v167 != 0 {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v144
								*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v167
								*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v167 + int32(31)) & int32(-32)
								v180 = int32(1)
							} else {
								v169 = F_WebPEncodingSetError(m, l0, int32(1))
								mBase = m.M
								v180 = v169
							}
						} else {
							v180 = int32(0)
						}
					}
				}
			}
			return v180
		} else {
			v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			v39 = int32(1)
			v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v44 = base.B2i32(v38 < v39) | base.B2i32(v41 < v39)
			if v44 == int32(0) {
				switch v37 {
				case 0, 4:
					v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
					F_WebPSafeFree(m, v53)
					mBase = m.M
					v55 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v55
					v57 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v57
					*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v55
					*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v55
					*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v57
					v71 = base.I64_extend_i32_s(v38)
					v72 = int64(1)
					v75 = int64(base.Ui64(v71+v72) >> (uint(v72) % 64))
					v76 = base.I32_wrap_i64(v75)
					if v44|base.B2i32(v76 < int32(1)) != 0 {
						v91 = F_WebPEncodingSetError(m, l0, int32(5))
						mBase = m.M
						v137 = v91
					} else {
						v80 = base.I64_extend_i32_s(v41)
						v81 = int64(1)
						v82 = v80 + v81
						if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v82)>>(uint(v81)%64))) {
							v96 = v37 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v38
							v98 = base.I64_extend_i32_s(v96) * v80
							v99 = v80 * v71
							v101 = int64(1)
							v104 = v82 >> (uint(v101) % 64) * base.I64_extend32_s(v75)
							v109 = F_WebPSafeMalloc(m, v98+v99+v104<<(uint(v101)%64), int32(1))
							mBase = m.M
							if v109 != 0 {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v96
								*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v76
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v38
								*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v109
								*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v109
								v118 = v109 + base.I32_wrap_i64(v99)
								*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v118
								v120 = base.I32_wrap_i64(v104)
								v121 = v118 + v120
								*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v121
								if v98 == int64(0) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v121 + v120
								}
								v137 = int32(1)
							} else {
								v111 = F_WebPEncodingSetError(m, l0, int32(1))
								mBase = m.M
								v137 = v111
							}
						} else {
							v91 = F_WebPEncodingSetError(m, l0, int32(5))
							mBase = m.M
							v137 = v91
						}
					}
				default:
					v49 = int32(4)
					v50 = F_WebPEncodingSetError(m, l0, v49)
					mBase = m.M
					if v50 != 0 {
						v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
						F_WebPSafeFree(m, v53)
						mBase = m.M
						v55 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v55
						v57 = int32(0)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v57
						*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v55
						*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v55
						*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v57
						v71 = base.I64_extend_i32_s(v38)
						v72 = int64(1)
						v75 = int64(base.Ui64(v71+v72) >> (uint(v72) % 64))
						v76 = base.I32_wrap_i64(v75)
						if v44|base.B2i32(v76 < int32(1)) != 0 {
							v91 = F_WebPEncodingSetError(m, l0, int32(5))
							mBase = m.M
							v137 = v91
						} else {
							v80 = base.I64_extend_i32_s(v41)
							v81 = int64(1)
							v82 = v80 + v81
							if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v82)>>(uint(v81)%64))) {
								v96 = v37 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v38
								v98 = base.I64_extend_i32_s(v96) * v80
								v99 = v80 * v71
								v101 = int64(1)
								v104 = v82 >> (uint(v101) % 64) * base.I64_extend32_s(v75)
								v109 = F_WebPSafeMalloc(m, v98+v99+v104<<(uint(v101)%64), int32(1))
								mBase = m.M
								if v109 != 0 {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v96
									*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v76
									*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v38
									*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v109
									*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v109
									v118 = v109 + base.I32_wrap_i64(v99)
									*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v118
									v120 = base.I32_wrap_i64(v104)
									v121 = v118 + v120
									*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v121
									if v98 == int64(0) {
									} else {
										*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v121 + v120
									}
									v137 = int32(1)
								} else {
									v111 = F_WebPEncodingSetError(m, l0, int32(1))
									mBase = m.M
									v137 = v111
								}
							} else {
								v91 = F_WebPEncodingSetError(m, l0, int32(5))
								mBase = m.M
								v137 = v91
							}
						}
					} else {
						v137 = int32(0)
					}
				}
			} else {
				v49 = int32(5)
				v50 = F_WebPEncodingSetError(m, l0, v49)
				mBase = m.M
				if v50 != 0 {
					v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
					F_WebPSafeFree(m, v53)
					mBase = m.M
					v55 = int64(0)
					*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v55
					v57 = int32(0)
					*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v57
					*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v55
					*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v55
					*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v57
					v71 = base.I64_extend_i32_s(v38)
					v72 = int64(1)
					v75 = int64(base.Ui64(v71+v72) >> (uint(v72) % 64))
					v76 = base.I32_wrap_i64(v75)
					if v44|base.B2i32(v76 < int32(1)) != 0 {
						v91 = F_WebPEncodingSetError(m, l0, int32(5))
						mBase = m.M
						v137 = v91
					} else {
						v80 = base.I64_extend_i32_s(v41)
						v81 = int64(1)
						v82 = v80 + v81
						if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v82)>>(uint(v81)%64))) {
							v96 = v37 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v38
							v98 = base.I64_extend_i32_s(v96) * v80
							v99 = v80 * v71
							v101 = int64(1)
							v104 = v82 >> (uint(v101) % 64) * base.I64_extend32_s(v75)
							v109 = F_WebPSafeMalloc(m, v98+v99+v104<<(uint(v101)%64), int32(1))
							mBase = m.M
							if v109 != 0 {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v96
								*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v76
								*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v38
								*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v109
								*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v109
								v118 = v109 + base.I32_wrap_i64(v99)
								*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v118
								v120 = base.I32_wrap_i64(v104)
								v121 = v118 + v120
								*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v121
								if v98 == int64(0) {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v121 + v120
								}
								v137 = int32(1)
							} else {
								v111 = F_WebPEncodingSetError(m, l0, int32(1))
								mBase = m.M
								v137 = v111
							}
						} else {
							v91 = F_WebPEncodingSetError(m, l0, int32(5))
							mBase = m.M
							v137 = v91
						}
					}
				} else {
					v137 = int32(0)
				}
			}
			return v137
		}
	} else {
		return int32(1)
	}
}
func F_WebPPictureAllocARGB(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v13 int32
	_ = v13
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v21 int32
	_ = v21
	var v23 int64
	_ = v23
	var v31 int64
	_ = v31
	var v32 int32
	_ = v32
	var v39 int64
	_ = v39
	var v40 int32
	_ = v40
	var v41 int64
	_ = v41
	var v51 int32
	_ = v51
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v6 = int32(5)
	v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	if v7 < int32(1) {
		v14 = v6
		v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		if v15 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v14
		}
		return int32(0)
	} else {
		if v5 < int32(1) {
			v14 = v6
			v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v15 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v14
			}
			return int32(0)
		} else {
			v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
			switch v13 {
			case 0, 4:
				v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
				F_free(m, v21)
				mBase = m.M
				v23 = int64(0)
				*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v23
				*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
				v31 = base.I64_extend_i32_s(v5)*base.I64_extend_i32_s(v7) + int64(31)
				v32 = int32(4)
				if v31 == v23 {
					v51 = F_malloc(m, base.I32_wrap_i64(v31)*v32)
					mBase = m.M
					v53 = v51
				} else {
					v39 = base.I64_div_u_s(int64(2147418112), v31)
					v40 = int32(0)
					v41 = base.I64_extend_i32_u(v32)
					if base.Ui64(int64(4294967295)) < base.Ui64(v41*v31) {
						v53 = v40
					} else {
						if base.Ui64(v39) < base.Ui64(v41) {
							v53 = v40
						} else {
							v51 = F_malloc(m, base.I32_wrap_i64(v31)*v32)
							mBase = m.M
							v53 = v51
						}
					}
				}
				if v53 != 0 {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v7
					*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v53
					*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v53 + int32(31)) & int32(-32)
					return int32(1)
				} else {
					v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v56 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(1)
					}
					return int32(0)
				}
			default:
				v14 = int32(4)
				v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v15 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v14
				}
				return int32(0)
			}
		}
	}
}
func F_WebPPictureAllocYUVA(m *base.Module, l0 int32) int32 {
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
	var v19 int32
	_ = v19
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v31 int32
	_ = v31
	var v33 int64
	_ = v33
	var v35 int32
	_ = v35
	var v49 int64
	_ = v49
	var v50 int64
	_ = v50
	var v53 int64
	_ = v53
	var v54 int32
	_ = v54
	var v58 int64
	_ = v58
	var v59 int64
	_ = v59
	var v60 int64
	_ = v60
	var v69 int32
	_ = v69
	var v77 int32
	_ = v77
	var v79 int64
	_ = v79
	var v80 int64
	_ = v80
	var v82 int64
	_ = v82
	var v85 int64
	_ = v85
	var v88 int64
	_ = v88
	var v89 int32
	_ = v89
	var v96 int64
	_ = v96
	var v97 int32
	_ = v97
	var v98 int64
	_ = v98
	var v108 int32
	_ = v108
	var v110 int32
	_ = v110
	var v113 int32
	_ = v113
	var v123 int32
	_ = v123
	var v125 int32
	_ = v125
	var v126 int32
	_ = v126
	v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v14 = int32(1)
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v19 = base.B2i32(v13 < v14) | base.B2i32(v16 < v14)
	if v19 == int32(0) {
		switch v12 {
		case 0, 4:
			v31 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
			F_free(m, v31)
			mBase = m.M
			v33 = int64(0)
			*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v33
			v35 = int32(0)
			*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v35
			*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v33
			*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v33
			*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = v35
			v49 = base.I64_extend_i32_s(v13)
			v50 = int64(1)
			v53 = int64(base.Ui64(v49+v50) >> (uint(v50) % 64))
			v54 = base.I32_wrap_i64(v53)
			if v19|base.B2i32(v54 < int32(1)) != 0 {
				v69 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v69 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(5)
				}
				return int32(0)
			} else {
				v58 = base.I64_extend_i32_s(v16)
				v59 = int64(1)
				v60 = v58 + v59
				if int32(0) < base.I32_wrap_i64(int64(base.Ui64(v60)>>(uint(v59)%64))) {
					v77 = v12 << (uint(int32(29)) % 32) >> (uint(int32(31)) % 32) & v13
					v79 = base.I64_extend_i32_s(v77) * v58
					v80 = v58 * v49
					v82 = int64(1)
					v85 = v60 >> (uint(v82) % 64) * base.I64_extend32_s(v53)
					v88 = v79 + v80 + v85<<(uint(v82)%64)
					v89 = int32(1)
					if v88 == int64(0) {
						v108 = F_malloc(m, base.I32_wrap_i64(v88)*v89)
						mBase = m.M
						v110 = v108
					} else {
						v96 = base.I64_div_u_s(int64(2147418112), v88)
						v97 = int32(0)
						v98 = base.I64_extend_i32_u(v89)
						if base.Ui64(int64(4294967295)) < base.Ui64(v98*v88) {
							v110 = v97
						} else {
							if base.Ui64(v96) < base.Ui64(v98) {
								v110 = v97
							} else {
								v108 = F_malloc(m, base.I32_wrap_i64(v88)*v89)
								mBase = m.M
								v110 = v108
							}
						}
					}
					if v110 != 0 {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+40)) = v77
						*(*int32)(unsafe.Add(mBase, uint32(l0)+32)) = v54
						*(*int32)(unsafe.Add(mBase, uint32(l0)+28)) = v13
						*(*int32)(unsafe.Add(mBase, uint32(l0)+156)) = v110
						*(*int32)(unsafe.Add(mBase, uint32(l0)+16)) = v110
						v123 = v110 + base.I32_wrap_i64(v80)
						*(*int32)(unsafe.Add(mBase, uint32(l0)+20)) = v123
						v125 = base.I32_wrap_i64(v85)
						v126 = v123 + v125
						*(*int32)(unsafe.Add(mBase, uint32(l0)+24)) = v126
						if v79 == int64(0) {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+36)) = v126 + v125
						}
						return int32(1)
					} else {
						v113 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
						if v113 != 0 {
						} else {
							*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(1)
						}
						return int32(0)
					}
				} else {
					v69 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v69 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(5)
					}
					return int32(0)
				}
			}
		default:
			v24 = int32(4)
			v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v25 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v24
			}
			return int32(0)
		}
	} else {
		v24 = int32(5)
		v25 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
		if v25 != 0 {
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v24
		}
		return int32(0)
	}
}
func F_WebPPictureFree(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v6 int32
	_ = v6
	var v8 int64
	_ = v8
	if l0 == int32(0) {
	} else {
		v4 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
		F_free(m, v4)
		mBase = m.M
		v6 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
		F_free(m, v6)
		mBase = m.M
		v8 = int64(0)
		*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v8
		*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v8
		*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = int32(0)
	}
	return
}
func F_WebPPictureHasTransparency(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v27 int32
	_ = v27
	var v28 int32
	_ = v28
	var v102 int32
	_ = v102
	var v105 int32
	_ = v105
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v109 int32
	_ = v109
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v122 int32
	_ = v122
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v128 int32
	_ = v128
	var v131 int32
	_ = v131
	var v132 int32
	_ = v132
	var v133 int32
	_ = v133
	var v135 int32
	_ = v135
	var v136 int32
	_ = v136
	var v206 int32
	_ = v206
	var v209 int32
	_ = v209
	var v211 int32
	_ = v211
	var v212 int32
	_ = v212
	var v213 int32
	_ = v213
	var v220 int32
	_ = v220
	if l0 == int32(0) {
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return int32(0)
L2:
	;
	v8 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	if v8 == int32(0) {
		goto L3
	} else {
		goto L4
	}
L3:
	;
	v119 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
	if v119 == int32(0) {
		goto L1
	} else {
		goto L15
	}
L4:
	;
	v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
	if v11 == int32(0) {
		goto L1
	} else {
		goto L5
	}
L5:
	;
	v14 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
	v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v20 = m.G1
	v23 = *(*int32)(unsafe.Add(mBase, uint32(v20)+uint32(_c_F_WebPPictureHasTransparency[0])))
	v24 = m.G3
	v25 = *(*int32)(unsafe.Add(mBase, uint32(v24)))
	if v23 == v25 {
		goto L7
	} else {
		goto L8
	}
L6:
	;
	if v15 < int32(1) {
		goto L1
	} else {
		goto L9
	}
L7:
	;
	goto L6
L8:
	;
	v27 = m.G2
	v28 = m.G1
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[1]))) = v27 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[2]))) = v27 + int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[3]))) = v27 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[4]))) = v27 + int32(5)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[5]))) = v27 + int32(6)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[6]))) = v27 + int32(7)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[7]))) = v27 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[8]))) = v27 + int32(9)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[9]))) = v27 + int32(10)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[10]))) = v27 + int32(11)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[11]))) = v27 + int32(12)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[12]))) = v27 + int32(13)
	*(*int32)(unsafe.Add(mBase, uint32(v28)+uint32(_c_F_WebPPictureHasTransparency[0]))) = v25
	goto L7
L9:
	;
	v102 = v15 + int32(1)
	v105 = v11 + int32(3)
	goto L10
L10:
	;
	v107 = m.G4
	v108 = *(*int32)(unsafe.Add(mBase, uint32(v107)))
	v109 = m.T0[v108].(func(*base.Module, int32, int32) int32)(m, v105, v16)
	mBase = m.M
	if v109 == int32(0) {
		goto L12
	} else {
		goto L13
	}
L12:
	;
	v116 = v102 + int32(-1)
	if int32(1) < v116 {
		v102 = v116
		v105 = v105 + v14<<(uint(int32(2))%32)
		goto L10
	} else {
		goto L14
	}
L13:
	;
	return int32(1)
L14:
	;
	goto L1
L15:
	;
	v122 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
	v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v124 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v128 = m.G1
	v131 = *(*int32)(unsafe.Add(mBase, uint32(v128)+uint32(_c_F_WebPPictureHasTransparency[0])))
	v132 = m.G3
	v133 = *(*int32)(unsafe.Add(mBase, uint32(v132)))
	if v131 == v133 {
		goto L17
	} else {
		goto L18
	}
L16:
	;
	if v124 < int32(1) {
		goto L1
	} else {
		goto L19
	}
L17:
	;
	goto L16
L18:
	;
	v135 = m.G2
	v136 = m.G1
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[1]))) = v135 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[2]))) = v135 + int32(3)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[3]))) = v135 + int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[4]))) = v135 + int32(5)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[5]))) = v135 + int32(6)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[6]))) = v135 + int32(7)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[7]))) = v135 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[8]))) = v135 + int32(9)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[9]))) = v135 + int32(10)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[10]))) = v135 + int32(11)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[11]))) = v135 + int32(12)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[12]))) = v135 + int32(13)
	*(*int32)(unsafe.Add(mBase, uint32(v136)+uint32(_c_F_WebPPictureHasTransparency[0]))) = v133
	goto L17
L19:
	;
	v206 = v124 + int32(1)
	v209 = v119
	goto L20
L20:
	;
	v211 = m.G5
	v212 = *(*int32)(unsafe.Add(mBase, uint32(v211)))
	v213 = m.T0[v212].(func(*base.Module, int32, int32) int32)(m, v209, v123)
	mBase = m.M
	if v213 == int32(0) {
		goto L22
	} else {
		goto L23
	}
L21:
	;
	goto L1
L22:
	;
	v220 = v206 + int32(-1)
	if int32(1) < v220 {
		v206 = v220
		v209 = v209 + v122
		goto L20
	} else {
		goto L24
	}
L23:
	;
	return int32(1)
L24:
	;
	goto L21
}
func F_WebPPictureImportRGBA(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v4 int32
	_ = v4
	var v16 int32
	_ = v16
	var v19 int32
	_ = v19
	var v23 int32
	_ = v23
	var v33 int32
	_ = v33
	var v35 int32
	_ = v35
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v41 int64
	_ = v41
	var v59 int32
	_ = v59
	var v60 int32
	_ = v60
	var v61 int32
	_ = v61
	var v62 int32
	_ = v62
	var v84 int32
	_ = v84
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v96 int32
	_ = v96
	var v142 int32
	_ = v142
	var v146 int32
	_ = v146
	var v149 int32
	_ = v149
	var v152 int32
	_ = v152
	var v155 int32
	_ = v155
	var v158 int32
	_ = v158
	var v161 int32
	_ = v161
	var v164 int32
	_ = v164
	var v167 int32
	_ = v167
	var v170 int32
	_ = v170
	var v173 int32
	_ = v173
	var v176 int32
	_ = v176
	var v179 int32
	_ = v179
	var v182 int32
	_ = v182
	var v271 int32
	_ = v271
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v347 int32
	_ = v347
	var v350 int32
	_ = v350
	var v356 int32
	_ = v356
	var v359 int32
	_ = v359
	var v361 int32
	_ = v361
	var v364 int32
	_ = v364
	var v365 int32
	_ = v365
	var v367 int32
	_ = v367
	var v368 int32
	_ = v368
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v372 int32
	_ = v372
	var v374 int32
	_ = v374
	var v377 int32
	_ = v377
	var v378 int32
	_ = v378
	var v380 int32
	_ = v380
	var v382 int32
	_ = v382
	var v385 int32
	_ = v385
	var v392 int32
	_ = v392
	var v393 int32
	_ = v393
	var v398 int32
	_ = v398
	v4 = int32(0)
	if l0 == v4 {
		v398 = v4
		return v398
	} else {
		if l1 == int32(0) {
			v398 = v4
			return v398
		} else {
			v16 = l2 >> (uint(int32(31)) % 32)
			v19 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
			if l2^v16-v16 < v19<<(uint(int32(2))%32) {
				v398 = v4
				return v398
			} else {
				v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
				if v23 != 0 {
					v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
					if l0 != 0 {
						v37 = *(*int32)(unsafe.Add(mBase, uint32(l0)+156))
						F_WebPSafeFree(m, v37)
						mBase = m.M
						v39 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
						F_WebPSafeFree(m, v39)
						mBase = m.M
						v41 = int64(0)
						*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v41
						*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v41
						*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v41
						*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v41
						*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v41
						*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = int32(0)
						v59 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
						if v59 != 0 {
							v61 = F_WebPPictureAllocARGB(m, l0)
							mBase = m.M
							v62 = v61
						} else {
							v60 = F_WebPPictureAllocYUVA(m, l0)
							mBase = m.M
							v62 = v60
						}
					} else {
						v62 = int32(1)
					}
					if v62 == int32(0) {
						v398 = v4
					} else {
						v84 = m.G1
						v87 = *(*int32)(unsafe.Add(mBase, uint32(v84)+uint32(_c_F_WebPPictureImportRGBA[0])))
						v88 = m.G3
						v89 = *(*int32)(unsafe.Add(mBase, uint32(v88)))
						if v87 == v89 {
						} else {
							v91 = m.G2
							v92 = m.G1
							v96 = v91 + int32(21)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[1]))) = v96
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[2]))) = v96
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[3]))) = v91 + int32(22)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[4]))) = v91 + int32(23)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[5]))) = v91 + int32(24)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[6]))) = v91 + int32(25)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[7]))) = v91 + int32(26)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[8]))) = v91 + int32(27)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[9]))) = v91 + int32(28)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[10]))) = v91 + int32(29)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[11]))) = v91 + int32(30)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[12]))) = v91 + int32(31)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[13]))) = v91 + int32(32)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[14]))) = v91 + int32(33)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[15]))) = v91 + int32(34)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[16]))) = v96
							v142 = v91 + int32(35)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[17]))) = v142
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[18]))) = v142
							v146 = v91 + int32(36)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[19]))) = v146
							v149 = v91 + int32(37)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[20]))) = v149
							v152 = v91 + int32(38)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[21]))) = v152
							v155 = v91 + int32(39)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[22]))) = v155
							v158 = v91 + int32(40)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[23]))) = v158
							v161 = v91 + int32(41)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[24]))) = v161
							v164 = v91 + int32(42)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[25]))) = v164
							v167 = v91 + int32(43)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[26]))) = v167
							v170 = v91 + int32(44)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[27]))) = v170
							v173 = v91 + int32(45)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[28]))) = v173
							v176 = v91 + int32(46)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[29]))) = v176
							v179 = v91 + int32(47)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[30]))) = v179
							v182 = v91 + int32(48)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[31]))) = v182
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[32]))) = v142
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[33]))) = v142
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[34]))) = v142
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[35]))) = v146
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[36]))) = v149
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[37]))) = v152
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[38]))) = v155
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[39]))) = v158
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[40]))) = v161
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[41]))) = v164
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[42]))) = v167
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[43]))) = v170
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[44]))) = v173
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[45]))) = v176
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[46]))) = v179
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[47]))) = v182
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[48]))) = v142
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[49]))) = v91 + int32(49)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[50]))) = v91 + int32(50)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[51]))) = v91 + int32(51)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[52]))) = v91 + int32(52)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[53]))) = v91 + int32(53)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[54]))) = v91 + int32(54)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[55]))) = v91 + int32(55)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[56]))) = v91 + int32(56)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[57]))) = v91 + int32(57)
							*(*int32)(unsafe.Add(mBase, uint32(v92)+uint32(_c_F_WebPPictureImportRGBA[0]))) = v89
						}
						v271 = m.G1
						v274 = *(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_WebPPictureImportRGBA[58])))
						v275 = m.G3
						v276 = *(*int32)(unsafe.Add(mBase, uint32(v275)))
						if v274 == v276 {
						} else {
							v278 = m.G2
							v279 = m.G1
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[59]))) = v278 + int32(2)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[60]))) = v278 + int32(3)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[61]))) = v278 + int32(4)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[62]))) = v278 + int32(5)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[63]))) = v278 + int32(6)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[64]))) = v278 + int32(7)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[65]))) = v278 + int32(8)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[66]))) = v278 + int32(9)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[67]))) = v278 + int32(10)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[68]))) = v278 + int32(11)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[69]))) = v278 + int32(12)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[70]))) = v278 + int32(13)
							*(*int32)(unsafe.Add(mBase, uint32(v279)+uint32(_c_F_WebPPictureImportRGBA[58]))) = v276
						}
						if v35 < int32(1) {
							v398 = int32(1)
						} else {
							v347 = int32(1)
							v350 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
							if v35 == v347 {
								v382 = l1
								v385 = v350
							} else {
								v356 = l1
								v359 = v350
								v361 = v35 & int32(-2)
								for {
									v364 = m.G10
									v365 = *(*int32)(unsafe.Add(mBase, uint32(v364)))
									m.T0[v365].(func(*base.Module, int32, int32, int32))(m, v356, v19, v359)
									mBase = m.M
									v367 = v356 + l2
									v368 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v369 = int32(2)
									v371 = v359 + v368<<(uint(v369)%32)
									v372 = *(*int32)(unsafe.Add(mBase, uint32(v364)))
									m.T0[v372].(func(*base.Module, int32, int32, int32))(m, v367, v19, v371)
									mBase = m.M
									v374 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v377 = v371 + v374<<(uint(v369)%32)
									v378 = v367 + l2
									v380 = v361 + int32(-2)
									if v380 != 0 {
										v356 = v378
										v359 = v377
										v361 = v380
										continue
									} else {
										break
									}
									break
								}
								v382 = v378
								v385 = v377
							}
							if v35&v347 == int32(0) {
								v398 = v347
							} else {
								v392 = m.G10
								v393 = *(*int32)(unsafe.Add(mBase, uint32(v392)))
								m.T0[v393].(func(*base.Module, int32, int32, int32))(m, v382, v19, v385)
								mBase = m.M
								v398 = v347
							}
						}
					}
					return v398
				} else {
					v33 = F_ImportYUVAFromRGBA(m, l1, l1+int32(1), l1+int32(2), l1+int32(3), int32(4), l2, float32(0), int32(0), l0)
					mBase = m.M
					return v33
				}
			}
		}
	}
}
func F_WebPPictureInitInternal(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v134 int32
	_ = v134
	var v139 int32
	_ = v139
	var v143 int32
	_ = v143
	if l1&int32(-256) != int32(512) {
		v143 = int32(0)
	} else {
		v9 = int32(1)
		if l0 == int32(0) {
			v143 = v9
		} else {
			base.MemoryFill(m, l0, int32(0), int32(172))
			v134 = m.G2
			*(*int32)(unsafe.Add(mBase, uint32(l0)+72)) = v134 + int32(1)
			v139 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v139 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(0)
			}
			v143 = v9
		}
	}
	return v143
}
func F_WebPPictureResetBuffers(m *base.Module, l0 int32) {
	mBase := m.M
	_ = mBase
	var v2 int64
	_ = v2
	v2 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = v2
	*(*int64)(unsafe.Add(mBase, uint32(l0)+156)) = v2
	*(*int64)(unsafe.Add(mBase, uint32(l0)+16)) = v2
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(24)))) = v2
	*(*int64)(unsafe.Add(mBase, uint32(l0+int32(32)))) = v2
	*(*int32)(unsafe.Add(mBase, uint32(l0+int32(40)))) = int32(0)
	return
}
func F_WebPPictureSharpARGBToYUVA(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v7 int32
	_ = v7
	var v13 int32
	_ = v13
	var v15 int32
	_ = v15
	var v20 int32
	_ = v20
	var v25 int32
	_ = v25
	if l0 != 0 {
		v5 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
		if v5 != 0 {
			*(*int32)(unsafe.Add(mBase, uint32(l0)+4)) = int32(0)
			v13 = int32(2)
			v15 = int32(1)
			v20 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
			v25 = F_ImportYUVAFromRGBA(m, v5+v13, v5+v15, v5, v5+int32(3), int32(4), v20<<(uint(v13)%32), float32(0), v15, l0)
			mBase = m.M
			return v25
		} else {
			v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v7 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(3)
			}
			return int32(0)
		}
	} else {
		return int32(0)
	}
}
func F_WebPPictureView(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32 {
	mBase := m.M
	_ = mBase
	var v7 int32
	_ = v7
	var v14 int32
	_ = v14
	var v17 int32
	_ = v17
	var v18 int32
	_ = v18
	var v21 int32
	_ = v21
	var v30 int32
	_ = v30
	var v33 int32
	_ = v33
	var v37 int32
	_ = v37
	var v38 int64
	_ = v38
	var v56 int32
	_ = v56
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v64 int32
	_ = v64
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v92 int32
	_ = v92
	var v93 int32
	_ = v93
	var v95 int32
	_ = v95
	var v103 int32
	_ = v103
	var v104 int32
	_ = v104
	var v105 int32
	_ = v105
	var v108 int32
	_ = v108
	var v119 int32
	_ = v119
	v7 = int32(0)
	if l0 == v7 {
		v119 = v7
	} else {
		if l5 == int32(0) {
			v119 = v7
		} else {
			v14 = int32(0)
			v17 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
			if v17 != 0 {
				v18 = l1
			} else {
				v18 = l1 & int32(-2)
			}
			if v17 != 0 {
				v21 = l2
			} else {
				v21 = l2 & int32(-2)
			}
			if v18|v21 < int32(0) {
				v119 = v14
			} else {
				if l3 < int32(1) {
					v119 = v14
				} else {
					if l4 < int32(1) {
						v119 = v14
					} else {
						v30 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if v30 < v18+l3 {
							v119 = v14
						} else {
							v33 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							if v33 < v21+l4 {
								v119 = v14
							} else {
								if l0 == l5 {
									v57 = v17
								} else {
									v37 = F_memcpy(m, l5, l0, int32(172))
									mBase = m.M
									v38 = int64(0)
									*(*int64)(unsafe.Add(mBase, uint32(v37)+52)) = v38
									*(*int64)(unsafe.Add(mBase, uint32(v37)+156)) = v38
									*(*int64)(unsafe.Add(mBase, uint32(v37)+16)) = v38
									*(*int64)(unsafe.Add(mBase, uint32(v37+int32(24)))) = v38
									*(*int64)(unsafe.Add(mBase, uint32(v37+int32(32)))) = v38
									*(*int32)(unsafe.Add(mBase, uint32(v37+int32(40)))) = int32(0)
									v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
									v57 = v56
								}
								*(*int32)(unsafe.Add(mBase, uint32(l5)+12)) = l4
								*(*int32)(unsafe.Add(mBase, uint32(l5)+8)) = l3
								if v57 != 0 {
									v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
									v93 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v95 = int32(2)
									v103 = v93
									v104 = int32(52)
									v105 = int32(56)
									v108 = v92 + v93*v21<<(uint(v95)%32) + v18<<(uint(v95)%32)
									*(*int32)(unsafe.Add(mBase, uint32(l5+v105))) = v103
									*(*int32)(unsafe.Add(mBase, uint32(l5+v104))) = v108
									v119 = int32(1)
								} else {
									v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+32)) = v60
									v62 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+28)) = v62
									v64 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+16)) = v64 + v62*v21 + v18
									v69 = int32(1)
									v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
									v73 = v60 * int32(base.Ui32(v21)>>(uint(v69)%32))
									v76 = int32(base.Ui32(v18) >> (uint(v69) % 32))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+20)) = v70 + v73 + v76
									v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
									*(*int32)(unsafe.Add(mBase, uint32(l5)+24)) = v79 + v73 + v76
									v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
									if v83 == int32(0) {
										v119 = v69
									} else {
										v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
										v103 = v86
										v104 = int32(36)
										v105 = int32(40)
										v108 = v83 + v86*v21 + v18
										*(*int32)(unsafe.Add(mBase, uint32(l5+v105))) = v103
										*(*int32)(unsafe.Add(mBase, uint32(l5+v104))) = v108
										v119 = int32(1)
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return v119
}
func F_WebPPictureYUVAToARGB(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v19 int32
	_ = v19
	var v20 int32
	_ = v20
	var v23 int32
	_ = v23
	var v26 int32
	_ = v26
	var v29 int32
	_ = v29
	var v34 int32
	_ = v34
	var v43 int32
	_ = v43
	var v44 int32
	_ = v44
	var v45 int32
	_ = v45
	var v51 int32
	_ = v51
	var v52 int32
	_ = v52
	var v53 int32
	_ = v53
	var v56 int32
	_ = v56
	var v68 int32
	_ = v68
	var v70 int32
	_ = v70
	var v81 int32
	_ = v81
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 int32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v92 int32
	_ = v92
	var v94 int32
	_ = v94
	var v102 int32
	_ = v102
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
	var v114 int32
	_ = v114
	var v117 int32
	_ = v117
	var v122 int32
	_ = v122
	var v125 int32
	_ = v125
	var v146 int32
	_ = v146
	var v153 int32
	_ = v153
	var v156 int32
	_ = v156
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v166 int32
	_ = v166
	var v167 int32
	_ = v167
	var v168 int32
	_ = v168
	var v170 int32
	_ = v170
	var v171 int32
	_ = v171
	var v175 int32
	_ = v175
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v183 int32
	_ = v183
	var v185 int32
	_ = v185
	var v187 int32
	_ = v187
	var v190 int32
	_ = v190
	var v192 int32
	_ = v192
	var v196 int32
	_ = v196
	var v197 int32
	_ = v197
	var v198 int32
	_ = v198
	var v201 int32
	_ = v201
	var v214 int32
	_ = v214
	var v217 int32
	_ = v217
	var v218 int32
	_ = v218
	var v225 int32
	_ = v225
	var v228 int32
	_ = v228
	var v229 int32
	_ = v229
	var v234 int32
	_ = v234
	var v241 int32
	_ = v241
	var v248 int32
	_ = v248
	var v251 int32
	_ = v251
	var v253 int32
	_ = v253
	var v254 int32
	_ = v254
	var v262 int32
	_ = v262
	var v263 int32
	_ = v263
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v285 int32
	_ = v285
	var v290 int32
	_ = v290
	var v294 int32
	_ = v294
	var v310 int32
	_ = v310
	var v319 int32
	_ = v319
	var v323 int32
	_ = v323
	var v325 int32
	_ = v325
	var v327 int32
	_ = v327
	var v329 int32
	_ = v329
	var v333 int32
	_ = v333
	v2 = int32(0)
	if l0 == v2 {
		v333 = v2
	} else {
		v19 = int32(3)
		v20 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
		if v20 == int32(0) {
			v327 = v19
			v329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v329 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v327
			}
			v333 = int32(0)
		} else {
			v23 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
			if v23 == int32(0) {
				v327 = v19
				v329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v329 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v327
				}
				v333 = int32(0)
			} else {
				v26 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
				if v26 == int32(0) {
					v327 = v19
					v329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v329 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v327
					}
					v333 = int32(0)
				} else {
					v29 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
					if v29&int32(4) == int32(0) {
						if v29&int32(3) != 0 {
							v327 = int32(4)
							v329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
							if v329 != 0 {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v327
							}
							v333 = int32(0)
						} else {
							v43 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
							v44 = int32(5)
							v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
							if v45 < int32(1) {
								v52 = v44
								v53 = F_WebPEncodingSetError(m, l0, v52)
								mBase = m.M
								if v53 != 0 {
									v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
									F_WebPSafeFree(m, v56)
									mBase = m.M
									*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
									*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
									v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
									mBase = m.M
									if v68 != 0 {
										*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
										*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
										v81 = int32(1)
									} else {
										v70 = F_WebPEncodingSetError(m, l0, int32(1))
										mBase = m.M
										v81 = v70
									}
								} else {
									v81 = int32(0)
								}
							} else {
								if v43 < int32(1) {
									v52 = v44
									v53 = F_WebPEncodingSetError(m, l0, v52)
									mBase = m.M
									if v53 != 0 {
										v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
										F_WebPSafeFree(m, v56)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
										v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
										mBase = m.M
										if v68 != 0 {
											*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
											*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
											v81 = int32(1)
										} else {
											v70 = F_WebPEncodingSetError(m, l0, int32(1))
											mBase = m.M
											v81 = v70
										}
									} else {
										v81 = int32(0)
									}
								} else {
									v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
									switch v51 {
									case 0, 4:
										v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
										F_WebPSafeFree(m, v56)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
										v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
										mBase = m.M
										if v68 != 0 {
											*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
											*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
											v81 = int32(1)
										} else {
											v70 = F_WebPEncodingSetError(m, l0, int32(1))
											mBase = m.M
											v81 = v70
										}
									default:
										v52 = int32(4)
										v53 = F_WebPEncodingSetError(m, l0, v52)
										mBase = m.M
										if v53 != 0 {
											v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
											F_WebPSafeFree(m, v56)
											mBase = m.M
											*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
											v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
											mBase = m.M
											if v68 != 0 {
												*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
												*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
												v81 = int32(1)
											} else {
												v70 = F_WebPEncodingSetError(m, l0, int32(1))
												mBase = m.M
												v81 = v70
											}
										} else {
											v81 = int32(0)
										}
									}
								}
							}
							if v81 == int32(0) {
								v333 = v2
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(1)
								v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
								v87 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
								v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
								v89 = int32(0)
								v90 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
								v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
								v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
								v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								v102 = m.G1
								v105 = *(*int32)(unsafe.Add(mBase, uint32(v102)+uint32(_c_F_WebPPictureYUVAToARGB[0])))
								v106 = m.G3
								v107 = *(*int32)(unsafe.Add(mBase, uint32(v106)))
								if v105 == v107 {
								} else {
									v109 = m.G2
									v110 = m.G1
									v114 = v109 + int32(14)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[1]))) = v114
									v117 = v109 + int32(15)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[2]))) = v117
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[3]))) = v114
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[4]))) = v117
									v122 = v109 + int32(16)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[5]))) = v122
									v125 = v109 + int32(17)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[6]))) = v125
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[7]))) = v109 + int32(18)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[8]))) = v122
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[9]))) = v125
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[10]))) = v109 + int32(19)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[11]))) = v109 + int32(20)
									*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[0]))) = v107
								}
								v146 = m.G1
								v153 = *(*int32)(unsafe.Add(mBase, uint32(v146)+uint32(_c_F_WebPPictureYUVAToARGB[3])))
								m.T0[v153].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v88, v89, v90, v91, v90, v91, v92, v89, v94)
								mBase = m.M
								v156 = v87 << (uint(int32(2)) % 32)
								v157 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
								v158 = v88 + v157
								if int32(3) <= v86 {
									v166 = v158
									v167 = v90
									v168 = v91
									v170 = int32(2)
									v171 = v92
									v175 = v157
									for {
										v181 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
										v182 = v167 + v181
										v183 = v168 + v181
										v185 = v171 + v87<<(uint(int32(3))%32)
										m.T0[v153].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v166, v166+v175, v167, v168, v182, v183, v171+v156, v185, v94)
										mBase = m.M
										v187 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
										v190 = v166 + v187<<(uint(int32(1))%32)
										v192 = v170 + int32(2)
										if v192 < v86 {
											v166 = v190
											v167 = v182
											v168 = v183
											v170 = v192
											v171 = v185
											v175 = v187
											continue
										} else {
											break
										}
										break
									}
									v196 = v190
									v197 = v182
									v198 = v183
									v201 = v185 + v156
								} else {
									v196 = v158
									v197 = v90
									v198 = v91
									v201 = v92 + v156
								}
								if v86 < int32(2) {
								} else {
									if v86&int32(1) != 0 {
									} else {
										v214 = int32(0)
										m.T0[v153].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v196, v214, v197, v198, v197, v198, v201, v214, v94)
										mBase = m.M
									}
								}
								v217 = int32(1)
								v218 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+4)))
								if v218&int32(4) == int32(0) {
									v333 = v217
								} else {
									if v86 < int32(1) {
										v333 = v217
									} else {
										v225 = int32(1)
										if v94 < v225 {
											v333 = v225
										} else {
											v228 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
											v229 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
											v234 = int32(0)
											v241 = v234
											v248 = v234
											for {
												v251 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
												v253 = v228 + v251*v241
												v254 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
												if v94 == int32(1) {
													v294 = int32(0)
												} else {
													v262 = v229 + v254*v248
													v263 = int32(0)
													for {
														v278 = v253 + v263
														v279 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v278))))
														*(*uint8)(unsafe.Add(mBase, uint32(v262+int32(3)))) = uint8(v279)
														v285 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v278+int32(1)))))
														*(*uint8)(unsafe.Add(mBase, uint32(v262+int32(7)))) = uint8(v285)
														v290 = v263 + int32(2)
														if v94&int32(2147483646) != v290 {
															v262 = v262 + int32(8)
															v263 = v290
															continue
														} else {
															break
														}
														break
													}
													v294 = v290
												}
												if v94&int32(1) == int32(0) {
												} else {
													v310 = int32(2)
													v319 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v253+v294))))
													*(*uint8)(unsafe.Add(mBase, uint32(v229+v254*v241<<(uint(v310)%32)+v294<<(uint(v310)%32)+int32(3)))) = uint8(v319)
												}
												v323 = int32(1)
												v325 = v241 + v323
												if v325 != v86 {
													v241 = v325
													v248 = v248 + int32(4)
													continue
												} else {
													break
												}
												break
											}
											v333 = v323
										}
									}
								}
							}
						}
					} else {
						v34 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
						if v34 == int32(0) {
							v327 = v19
							v329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
							if v329 != 0 {
							} else {
								*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v327
							}
							v333 = int32(0)
						} else {
							if v29&int32(3) != 0 {
								v327 = int32(4)
								v329 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
								if v329 != 0 {
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v327
								}
								v333 = int32(0)
							} else {
								v43 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
								v44 = int32(5)
								v45 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
								if v45 < int32(1) {
									v52 = v44
									v53 = F_WebPEncodingSetError(m, l0, v52)
									mBase = m.M
									if v53 != 0 {
										v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
										F_WebPSafeFree(m, v56)
										mBase = m.M
										*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
										*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
										v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
										mBase = m.M
										if v68 != 0 {
											*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
											*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
											v81 = int32(1)
										} else {
											v70 = F_WebPEncodingSetError(m, l0, int32(1))
											mBase = m.M
											v81 = v70
										}
									} else {
										v81 = int32(0)
									}
								} else {
									if v43 < int32(1) {
										v52 = v44
										v53 = F_WebPEncodingSetError(m, l0, v52)
										mBase = m.M
										if v53 != 0 {
											v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
											F_WebPSafeFree(m, v56)
											mBase = m.M
											*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
											v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
											mBase = m.M
											if v68 != 0 {
												*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
												*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
												v81 = int32(1)
											} else {
												v70 = F_WebPEncodingSetError(m, l0, int32(1))
												mBase = m.M
												v81 = v70
											}
										} else {
											v81 = int32(0)
										}
									} else {
										v51 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
										switch v51 {
										case 0, 4:
											v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
											F_WebPSafeFree(m, v56)
											mBase = m.M
											*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
											*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
											v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
											mBase = m.M
											if v68 != 0 {
												*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
												*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
												v81 = int32(1)
											} else {
												v70 = F_WebPEncodingSetError(m, l0, int32(1))
												mBase = m.M
												v81 = v70
											}
										default:
											v52 = int32(4)
											v53 = F_WebPEncodingSetError(m, l0, v52)
											mBase = m.M
											if v53 != 0 {
												v56 = *(*int32)(unsafe.Add(mBase, uint32(l0)+160))
												F_WebPSafeFree(m, v56)
												mBase = m.M
												*(*int64)(unsafe.Add(mBase, uint32(l0)+52)) = int64(0)
												*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = int32(0)
												v68 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v43)*base.I64_extend_i32_s(v45)+int64(31), int32(4))
												mBase = m.M
												if v68 != 0 {
													*(*int32)(unsafe.Add(mBase, uint32(l0)+56)) = v45
													*(*int32)(unsafe.Add(mBase, uint32(l0)+160)) = v68
													*(*int32)(unsafe.Add(mBase, uint32(l0)+52)) = (v68 + int32(31)) & int32(-32)
													v81 = int32(1)
												} else {
													v70 = F_WebPEncodingSetError(m, l0, int32(1))
													mBase = m.M
													v81 = v70
												}
											} else {
												v81 = int32(0)
											}
										}
									}
								}
								if v81 == int32(0) {
									v333 = v2
								} else {
									*(*int32)(unsafe.Add(mBase, uint32(l0))) = int32(1)
									v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
									v87 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
									v88 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
									v89 = int32(0)
									v90 = *(*int32)(unsafe.Add(mBase, uint32(l0)+20))
									v91 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
									v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
									v94 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
									v102 = m.G1
									v105 = *(*int32)(unsafe.Add(mBase, uint32(v102)+uint32(_c_F_WebPPictureYUVAToARGB[0])))
									v106 = m.G3
									v107 = *(*int32)(unsafe.Add(mBase, uint32(v106)))
									if v105 == v107 {
									} else {
										v109 = m.G2
										v110 = m.G1
										v114 = v109 + int32(14)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[1]))) = v114
										v117 = v109 + int32(15)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[2]))) = v117
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[3]))) = v114
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[4]))) = v117
										v122 = v109 + int32(16)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[5]))) = v122
										v125 = v109 + int32(17)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[6]))) = v125
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[7]))) = v109 + int32(18)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[8]))) = v122
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[9]))) = v125
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[10]))) = v109 + int32(19)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[11]))) = v109 + int32(20)
										*(*int32)(unsafe.Add(mBase, uint32(v110)+uint32(_c_F_WebPPictureYUVAToARGB[0]))) = v107
									}
									v146 = m.G1
									v153 = *(*int32)(unsafe.Add(mBase, uint32(v146)+uint32(_c_F_WebPPictureYUVAToARGB[3])))
									m.T0[v153].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v88, v89, v90, v91, v90, v91, v92, v89, v94)
									mBase = m.M
									v156 = v87 << (uint(int32(2)) % 32)
									v157 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
									v158 = v88 + v157
									if int32(3) <= v86 {
										v166 = v158
										v167 = v90
										v168 = v91
										v170 = int32(2)
										v171 = v92
										v175 = v157
										for {
											v181 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
											v182 = v167 + v181
											v183 = v168 + v181
											v185 = v171 + v87<<(uint(int32(3))%32)
											m.T0[v153].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v166, v166+v175, v167, v168, v182, v183, v171+v156, v185, v94)
											mBase = m.M
											v187 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
											v190 = v166 + v187<<(uint(int32(1))%32)
											v192 = v170 + int32(2)
											if v192 < v86 {
												v166 = v190
												v167 = v182
												v168 = v183
												v170 = v192
												v171 = v185
												v175 = v187
												continue
											} else {
												break
											}
											break
										}
										v196 = v190
										v197 = v182
										v198 = v183
										v201 = v185 + v156
									} else {
										v196 = v158
										v197 = v90
										v198 = v91
										v201 = v92 + v156
									}
									if v86 < int32(2) {
									} else {
										if v86&int32(1) != 0 {
										} else {
											v214 = int32(0)
											m.T0[v153].(func(*base.Module, int32, int32, int32, int32, int32, int32, int32, int32, int32))(m, v196, v214, v197, v198, v197, v198, v201, v214, v94)
											mBase = m.M
										}
									}
									v217 = int32(1)
									v218 = int32(*(*uint8)(unsafe.Add(mBase, uint32(l0)+4)))
									if v218&int32(4) == int32(0) {
										v333 = v217
									} else {
										if v86 < int32(1) {
											v333 = v217
										} else {
											v225 = int32(1)
											if v94 < v225 {
												v333 = v225
											} else {
												v228 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
												v229 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
												v234 = int32(0)
												v241 = v234
												v248 = v234
												for {
													v251 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
													v253 = v228 + v251*v241
													v254 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
													if v94 == int32(1) {
														v294 = int32(0)
													} else {
														v262 = v229 + v254*v248
														v263 = int32(0)
														for {
															v278 = v253 + v263
															v279 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v278))))
															*(*uint8)(unsafe.Add(mBase, uint32(v262+int32(3)))) = uint8(v279)
															v285 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v278+int32(1)))))
															*(*uint8)(unsafe.Add(mBase, uint32(v262+int32(7)))) = uint8(v285)
															v290 = v263 + int32(2)
															if v94&int32(2147483646) != v290 {
																v262 = v262 + int32(8)
																v263 = v290
																continue
															} else {
																break
															}
															break
														}
														v294 = v290
													}
													if v94&int32(1) == int32(0) {
													} else {
														v310 = int32(2)
														v319 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v253+v294))))
														*(*uint8)(unsafe.Add(mBase, uint32(v229+v254*v241<<(uint(v310)%32)+v294<<(uint(v310)%32)+int32(3)))) = uint8(v319)
													}
													v323 = int32(1)
													v325 = v241 + v323
													if v325 != v86 {
														v241 = v325
														v248 = v248 + int32(4)
														continue
													} else {
														break
													}
													break
												}
												v333 = v323
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return v333
}
func F_WebPReplaceTransparentPixels(m *base.Module, l0 int32, l1 int32) {
	mBase := m.M
	_ = mBase
	var v9 int32
	_ = v9
	var v12 int32
	_ = v12
	var v13 int32
	_ = v13
	var v17 int32
	_ = v17
	var v20 int32
	_ = v20
	var v21 int32
	_ = v21
	var v22 int32
	_ = v22
	var v24 int32
	_ = v24
	var v25 int32
	_ = v25
	var v94 int32
	_ = v94
	var v97 int32
	_ = v97
	var v98 int32
	_ = v98
	var v99 int32
	_ = v99
	var v103 int32
	_ = v103
	var v107 int32
	_ = v107
	var v108 int32
	_ = v108
	var v115 int32
	_ = v115
	var v116 int32
	_ = v116
	var v119 int32
	_ = v119
	var v120 int32
	_ = v120
	var v121 int32
	_ = v121
	var v123 int32
	_ = v123
	var v124 int32
	_ = v124
	var v126 int32
	_ = v126
	var v127 int32
	_ = v127
	var v128 int32
	_ = v128
	var v130 int32
	_ = v130
	var v135 int32
	_ = v135
	if l0 == int32(0) {
	} else {
		v9 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
		if v9 == int32(0) {
		} else {
			v12 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			v13 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
			v17 = m.G1
			v20 = *(*int32)(unsafe.Add(mBase, uint32(v17)+uint32(_c_F_WebPReplaceTransparentPixels[0])))
			v21 = m.G3
			v22 = *(*int32)(unsafe.Add(mBase, uint32(v21)))
			if v20 == v22 {
			} else {
				v24 = m.G2
				v25 = m.G1
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[1]))) = v24 + int32(2)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[2]))) = v24 + int32(3)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[3]))) = v24 + int32(4)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[4]))) = v24 + int32(5)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[5]))) = v24 + int32(6)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[6]))) = v24 + int32(7)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[7]))) = v24 + int32(8)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[8]))) = v24 + int32(9)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[9]))) = v24 + int32(10)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[10]))) = v24 + int32(11)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[11]))) = v24 + int32(12)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[12]))) = v24 + int32(13)
				*(*int32)(unsafe.Add(mBase, uint32(v25)+uint32(_c_F_WebPReplaceTransparentPixels[0]))) = v22
			}
			if v12 < int32(1) {
			} else {
				v94 = l1 & int32(16777215)
				if v12&int32(1) != 0 {
					v97 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
					v98 = m.G11
					v99 = *(*int32)(unsafe.Add(mBase, uint32(v98)))
					m.T0[v99].(func(*base.Module, int32, int32, int32))(m, v13, v97, v94)
					mBase = m.M
					v103 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
					v107 = v13 + v103<<(uint(int32(2))%32)
					v108 = v12 + int32(-1)
				} else {
					v107 = v13
					v108 = v12
				}
				if v12 == int32(1) {
				} else {
					v115 = v108 + int32(-1)
					v116 = v107
					for {
						v119 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v120 = m.G11
						v121 = *(*int32)(unsafe.Add(mBase, uint32(v120)))
						m.T0[v121].(func(*base.Module, int32, int32, int32))(m, v116, v119, v94)
						mBase = m.M
						v123 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v124 = int32(2)
						v126 = v116 + v123<<(uint(v124)%32)
						v127 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						v128 = *(*int32)(unsafe.Add(mBase, uint32(v120)))
						m.T0[v128].(func(*base.Module, int32, int32, int32))(m, v126, v127, v94)
						mBase = m.M
						v130 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
						v135 = v115 + int32(-2)
						if base.Ui32(v135) < base.Ui32(v115) {
							v115 = v135
							v116 = v126 + v130<<(uint(v124)%32)
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
	return
}
func F_WebPReportProgress(m *base.Module, l0 int32, l1 int32, l2 int32) int32 {
	mBase := m.M
	_ = mBase
	var v5 int32
	_ = v5
	var v8 int32
	_ = v8
	var v11 int32
	_ = v11
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v20 int32
	_ = v20
	v5 = int32(1)
	if l2 == int32(0) {
		v20 = v5
	} else {
		v8 = *(*int32)(unsafe.Add(mBase, uint32(l2)))
		if v8 == l1 {
			v20 = v5
		} else {
			*(*int32)(unsafe.Add(mBase, uint32(l2))) = l1
			v11 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
			if v11 == int32(0) {
				v20 = v5
			} else {
				v14 = m.T0[v11].(func(*base.Module, int32, int32) int32)(m, l1, l0)
				mBase = m.M
				if v14 != 0 {
					v20 = v5
				} else {
					v15 = int32(0)
					v16 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v16 != 0 {
						v20 = v15
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = int32(10)
						v20 = v15
					}
				}
			}
		}
	}
	return v20
}
func F_WebPSafeCalloc(m *base.Module, l0 int64, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v11 int64
	_ = v11
	var v19 int32
	_ = v19
	var v25 int64
	_ = v25
	var v26 int32
	_ = v26
	var v36 int32
	_ = v36
	var v37 int32
	_ = v37
	var v39 int32
	_ = v39
	var v44 int32
	_ = v44
	var v50 int32
	_ = v50
	var v52 int32
	_ = v52
	if l0 == int64(0) {
		v19 = base.I32_wrap_i64(l0)
		if v19 != 0 {
			v25 = base.I64_extend_i32_u(v19) * base.I64_extend_i32_u(l1)
			v26 = base.I32_wrap_i64(v25)
			if base.Ui32(l1|v19) < base.Ui32(int32(_a_F_WebPSafeCalloc_0)) {
				v37 = v26
			} else {
				if base.I32_wrap_i64(int64(base.Ui64(v25)>>(uint(int64(32))%64))) != int32(0) {
					v36 = int32(-1)
				} else {
					v36 = v26
				}
				v37 = v36
			}
		} else {
			v37 = int32(0)
		}
		v39 = F_dlmalloc(m, v37)
		mBase = m.M
		if v39 == int32(0) {
		} else {
			v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v39+int32(-4)))))
			if v44&int32(3) == int32(0) {
			} else {
				v50 = F_memset(m, v39, int32(0), v37)
				mBase = m.M
			}
		}
		v52 = v39
	} else {
		v9 = base.I64_div_u_s(int64(2147418112), l0)
		v10 = int32(0)
		v11 = base.I64_extend_i32_u(l1)
		if base.Ui64(int64(4294967295)) < base.Ui64(v11*l0) {
			v52 = v10
		} else {
			if base.Ui64(v9) < base.Ui64(v11) {
				v52 = v10
			} else {
				v19 = base.I32_wrap_i64(l0)
				if v19 != 0 {
					v25 = base.I64_extend_i32_u(v19) * base.I64_extend_i32_u(l1)
					v26 = base.I32_wrap_i64(v25)
					if base.Ui32(l1|v19) < base.Ui32(int32(_a_F_WebPSafeCalloc_0)) {
						v37 = v26
					} else {
						if base.I32_wrap_i64(int64(base.Ui64(v25)>>(uint(int64(32))%64))) != int32(0) {
							v36 = int32(-1)
						} else {
							v36 = v26
						}
						v37 = v36
					}
				} else {
					v37 = int32(0)
				}
				v39 = F_dlmalloc(m, v37)
				mBase = m.M
				if v39 == int32(0) {
				} else {
					v44 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v39+int32(-4)))))
					if v44&int32(3) == int32(0) {
					} else {
						v50 = F_memset(m, v39, int32(0), v37)
						mBase = m.M
					}
				}
				v52 = v39
			}
		}
	}
	return v52
}
func F_WebPSafeFree(m *base.Module, l0 int32) {
	F_dlfree(m, l0)
	return
}
func F_WebPSafeMalloc(m *base.Module, l0 int64, l1 int32) int32 {
	var v9 int64
	_ = v9
	var v10 int32
	_ = v10
	var v11 int64
	_ = v11
	var v21 int32
	_ = v21
	var v23 int32
	_ = v23
	if l0 == int64(0) {
		v21 = F_dlmalloc(m, base.I32_wrap_i64(l0)*l1)
		v23 = v21
	} else {
		v9 = base.I64_div_u_s(int64(2147418112), l0)
		v10 = int32(0)
		v11 = base.I64_extend_i32_u(l1)
		if base.Ui64(int64(4294967295)) < base.Ui64(v11*l0) {
			v23 = v10
		} else {
			if base.Ui64(v9) < base.Ui64(v11) {
				v23 = v10
			} else {
				v21 = F_dlmalloc(m, base.I32_wrap_i64(l0)*l1)
				v23 = v21
			}
		}
	}
	return v23
}
func F_WebPValidateConfig(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v2 int32
	_ = v2
	var v9 float32
	_ = v9
	var v14 int32
	_ = v14
	var v15 int32
	_ = v15
	var v18 float32
	_ = v18
	var v21 int32
	_ = v21
	var v24 int32
	_ = v24
	var v29 int32
	_ = v29
	var v32 int32
	_ = v32
	var v35 int32
	_ = v35
	var v38 int32
	_ = v38
	var v41 int32
	_ = v41
	var v44 int32
	_ = v44
	var v49 int32
	_ = v49
	var v50 int32
	_ = v50
	var v53 int32
	_ = v53
	var v57 int32
	_ = v57
	var v60 int32
	_ = v60
	var v63 int32
	_ = v63
	var v66 int32
	_ = v66
	var v69 int32
	_ = v69
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v77 int32
	_ = v77
	var v80 int32
	_ = v80
	var v83 int32
	_ = v83
	var v86 int32
	_ = v86
	var v89 int32
	_ = v89
	var v92 int32
	_ = v92
	var v95 int32
	_ = v95
	var v98 int32
	_ = v98
	var v101 int32
	_ = v101
	var v104 int32
	_ = v104
	v2 = int32(0)
	if l0 == v2 {
		v104 = v2
	} else {
		v9 = *(*float32)(unsafe.Add(mBase, uint32(l0)+4))
		if base.F32_lt(v9, float32(0)) != 0 {
			v104 = v2
		} else {
			if base.F32_gt(v9, float32(100)) != 0 {
				v104 = v2
			} else {
				v14 = int32(0)
				v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+16))
				if v15 < v14 {
					v104 = v14
				} else {
					v18 = *(*float32)(unsafe.Add(mBase, uint32(l0)+20))
					if base.F32_lt(v18, float32(0)) != 0 {
						v104 = v14
					} else {
						v21 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
						if base.Ui32(int32(6)) < base.Ui32(v21) {
							v104 = v14
						} else {
							v24 = *(*int32)(unsafe.Add(mBase, uint32(l0)+24))
							if base.Ui32(v24+int32(-5)) < base.Ui32(int32(-4)) {
								v104 = v14
							} else {
								v29 = *(*int32)(unsafe.Add(mBase, uint32(l0)+28))
								if base.Ui32(int32(100)) < base.Ui32(v29) {
									v104 = v14
								} else {
									v32 = *(*int32)(unsafe.Add(mBase, uint32(l0)+32))
									if base.Ui32(int32(100)) < base.Ui32(v32) {
										v104 = v14
									} else {
										v35 = *(*int32)(unsafe.Add(mBase, uint32(l0)+36))
										if base.Ui32(int32(7)) < base.Ui32(v35) {
											v104 = v14
										} else {
											v38 = *(*int32)(unsafe.Add(mBase, uint32(l0)+40))
											if base.Ui32(int32(1)) < base.Ui32(v38) {
												v104 = v14
											} else {
												v41 = *(*int32)(unsafe.Add(mBase, uint32(l0)+44))
												if base.Ui32(int32(1)) < base.Ui32(v41) {
													v104 = v14
												} else {
													v44 = *(*int32)(unsafe.Add(mBase, uint32(l0)+60))
													if base.Ui32(v44+int32(-11)) < base.Ui32(int32(-10)) {
														v104 = v14
													} else {
														v49 = int32(0)
														v50 = *(*int32)(unsafe.Add(mBase, uint32(l0)+108))
														if v50 < v49 {
															v104 = v49
														} else {
															v53 = *(*int32)(unsafe.Add(mBase, uint32(l0)+112))
															if int32(100) < v53 {
																v104 = v49
															} else {
																if v53 < v50 {
																	v104 = v49
																} else {
																	v57 = *(*int32)(unsafe.Add(mBase, uint32(l0)+64))
																	if base.Ui32(int32(1)) < base.Ui32(v57) {
																		v104 = v49
																	} else {
																		v60 = *(*int32)(unsafe.Add(mBase, uint32(l0)+68))
																		if base.Ui32(int32(7)) < base.Ui32(v60) {
																			v104 = v49
																		} else {
																			v63 = *(*int32)(unsafe.Add(mBase, uint32(l0)+72))
																			if base.Ui32(int32(3)) < base.Ui32(v63) {
																				v104 = v49
																			} else {
																				v66 = *(*int32)(unsafe.Add(mBase, uint32(l0)+76))
																				if base.Ui32(int32(100)) < base.Ui32(v66) {
																					v104 = v49
																				} else {
																					v69 = int32(0)
																					v70 = *(*int32)(unsafe.Add(mBase, uint32(l0)+48))
																					if v70 < v69 {
																						v104 = v69
																					} else {
																						v73 = int32(0)
																						v74 = *(*int32)(unsafe.Add(mBase, uint32(l0)+52))
																						if v74 < v73 {
																							v104 = v73
																						} else {
																							v77 = *(*int32)(unsafe.Add(mBase, uint32(l0)+56))
																							if base.Ui32(int32(100)) < base.Ui32(v77) {
																								v104 = v73
																							} else {
																								v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
																								if base.Ui32(int32(1)) < base.Ui32(v80) {
																									v104 = v73
																								} else {
																									v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
																									if base.Ui32(int32(100)) < base.Ui32(v83) {
																										v104 = v73
																									} else {
																										v86 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
																										if base.Ui32(int32(3)) < base.Ui32(v86) {
																											v104 = v73
																										} else {
																											v89 = *(*int32)(unsafe.Add(mBase, uint32(l0)+80))
																											if base.Ui32(int32(1)) < base.Ui32(v89) {
																												v104 = v73
																											} else {
																												v92 = *(*int32)(unsafe.Add(mBase, uint32(l0)+84))
																												if base.Ui32(int32(1)) < base.Ui32(v92) {
																													v104 = v73
																												} else {
																													v95 = *(*int32)(unsafe.Add(mBase, uint32(l0)+88))
																													if base.Ui32(int32(1)) < base.Ui32(v95) {
																														v104 = v73
																													} else {
																														v98 = *(*int32)(unsafe.Add(mBase, uint32(l0)+96))
																														if base.Ui32(int32(1)) < base.Ui32(v98) {
																															v104 = v73
																														} else {
																															v101 = *(*int32)(unsafe.Add(mBase, uint32(l0)+104))
																															v104 = base.B2i32(base.Ui32(v101) < base.Ui32(int32(2)))
																														}
																													}
																												}
																											}
																										}
																									}
																								}
																							}
																						}
																					}
																				}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return v104
}
func F_WebPValidatePicture(m *base.Module, l0 int32) int32 {
	mBase := m.M
	_ = mBase
	var v6 int32
	_ = v6
	var v7 int32
	_ = v7
	var v10 int32
	_ = v10
	var v15 int32
	_ = v15
	var v16 int32
	_ = v16
	var v18 int32
	_ = v18
	var v22 int32
	_ = v22
	if l0 != 0 {
		v6 = int32(5)
		v7 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
		if v7 < int32(1) {
			v16 = v6
			v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
			if v18 != 0 {
			} else {
				*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v16
			}
			v22 = int32(0)
		} else {
			v10 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
			if v10 < int32(1) {
				v16 = v6
				v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
				if v18 != 0 {
				} else {
					*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v16
				}
				v22 = int32(0)
			} else {
				v15 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
				switch v15 {
				case 0, 4:
					v22 = int32(1)
				default:
					v16 = int32(4)
					v18 = *(*int32)(unsafe.Add(mBase, uint32(l0)+92))
					if v18 != 0 {
					} else {
						*(*int32)(unsafe.Add(mBase, uint32(l0)+92)) = v16
					}
					v22 = int32(0)
				}
			}
		}
		return v22
	} else {
		return int32(0)
	}
}

//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_VP8LGetHistoImageSymbols(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32) int32 {
	mBase := m.M
	_ = mBase
	var v52 int32
	_ = v52
	var v54 int32
	_ = v54
	var v60 int32
	_ = v60
	var v62 int32
	_ = v62
	var v70 int32
	_ = v70
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v83 int64
	_ = v83
	var v84 int32
	_ = v84
	var v91 int64
	_ = v91
	var v92 int32
	_ = v92
	var v93 int64
	_ = v93
	var v103 int32
	_ = v103
	var v105 int32
	_ = v105
	var v112 int32
	_ = v112
	var v120 int32
	_ = v120
	var v122 int32
	_ = v122
	var v124 int32
	_ = v124
	var v127 int32
	_ = v127
	var v133 int32
	_ = v133
	var v147 int32
	_ = v147
	var v148 int32
	_ = v148
	var v151 int32
	_ = v151
	var v184 int32
	_ = v184
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v191 int32
	_ = v191
	var v192 int32
	_ = v192
	var v195 int32
	_ = v195
	var v196 int32
	_ = v196
	var v203 int32
	_ = v203
	var v205 int32
	_ = v205
	var v207 int32
	_ = v207
	var v213 int32
	_ = v213
	var v215 int32
	_ = v215
	var v232 int32
	_ = v232
	var v233 int32
	_ = v233
	var v236 int32
	_ = v236
	var v271 int32
	_ = v271
	var v276 int32
	_ = v276
	var v278 int32
	_ = v278
	var v280 int32
	_ = v280
	var v308 int32
	_ = v308
	var v312 int32
	_ = v312
	var v344 int32
	_ = v344
	var v348 int32
	_ = v348
	var v349 int64
	_ = v349
	var v351 int32
	_ = v351
	var v353 base.V128
	_ = v353
	var v354 int32
	_ = v354
	var v356 int32
	_ = v356
	var v358 int32
	_ = v358
	var v360 int32
	_ = v360
	var v362 int32
	_ = v362
	var v364 int32
	_ = v364
	var v366 int32
	_ = v366
	var v368 int32
	_ = v368
	var v373 int32
	_ = v373
	var v375 int32
	_ = v375
	var v402 int32
	_ = v402
	var v423 int32
	_ = v423
	var v457 int32
	_ = v457
	var v461 int32
	_ = v461
	var v467 base.V128
	_ = v467
	var v472 int32
	_ = v472
	var v476 int32
	_ = v476
	var v480 int32
	_ = v480
	var v488 int32
	_ = v488
	var v508 int32
	_ = v508
	var v541 int32
	_ = v541
	var v543 int32
	_ = v543
	var v545 int32
	_ = v545
	var v546 int32
	_ = v546
	var v550 int32
	_ = v550
	var v551 int32
	_ = v551
	var v554 int32
	_ = v554
	var v557 int32
	_ = v557
	var v558 int32
	_ = v558
	var v559 int32
	_ = v559
	var v560 int32
	_ = v560
	var v561 int32
	_ = v561
	var v565 int32
	_ = v565
	var v566 int32
	_ = v566
	var v567 int32
	_ = v567
	var v570 int32
	_ = v570
	var v581 int32
	_ = v581
	var v603 int32
	_ = v603
	var v606 int32
	_ = v606
	var v607 int32
	_ = v607
	var v615 int32
	_ = v615
	var v616 int32
	_ = v616
	var v651 int32
	_ = v651
	var v652 int32
	_ = v652
	var v657 int64
	_ = v657
	var v660 int32
	_ = v660
	var v661 int32
	_ = v661
	var v673 int32
	_ = v673
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v703 int32
	_ = v703
	var v706 int32
	_ = v706
	var v712 int32
	_ = v712
	var v717 int32
	_ = v717
	var v727 int32
	_ = v727
	var v729 int32
	_ = v729
	var v763 int32
	_ = v763
	var v765 int32
	_ = v765
	var v767 int32
	_ = v767
	var v768 int32
	_ = v768
	var v770 int32
	_ = v770
	var v771 int32
	_ = v771
	var v774 int32
	_ = v774
	var v775 int32
	_ = v775
	var v782 int32
	_ = v782
	var v784 int32
	_ = v784
	var v786 int32
	_ = v786
	var v792 int32
	_ = v792
	var v794 int32
	_ = v794
	var v802 int32
	_ = v802
	var v812 int32
	_ = v812
	var v814 int32
	_ = v814
	var v850 int32
	_ = v850
	var v855 int32
	_ = v855
	var v857 int32
	_ = v857
	var v859 int32
	_ = v859
	var v865 int32
	_ = v865
	var v870 int32
	_ = v870
	var v872 int32
	_ = v872
	var v877 int32
	_ = v877
	var v924 int32
	_ = v924
	var v927 int32
	_ = v927
	var v930 int32
	_ = v930
	var v931 int32
	_ = v931
	var v939 int32
	_ = v939
	var v961 int32
	_ = v961
	var v997 int32
	_ = v997
	var v1002 int32
	_ = v1002
	var v1047 int32
	_ = v1047
	var v1052 int32
	_ = v1052
	var v1104 int32
	_ = v1104
	var v1113 int32
	_ = v1113
	var v1116 int32
	_ = v1116
	var v1121 int32
	_ = v1121
	var v1130 int32
	_ = v1130
	var v1168 int32
	_ = v1168
	var v1175 int32
	_ = v1175
	var v1176 int32
	_ = v1176
	var v1181 int32
	_ = v1181
	var v1182 int32
	_ = v1182
	var v1185 int32
	_ = v1185
	var v1199 int32
	_ = v1199
	var v1236 int32
	_ = v1236
	var v1237 int32
	_ = v1237
	var v1240 int32
	_ = v1240
	var v1254 int32
	_ = v1254
	var v1290 int32
	_ = v1290
	var v1292 int32
	_ = v1292
	var v1294 int32
	_ = v1294
	var v1299 int32
	_ = v1299
	var v1300 int32
	_ = v1300
	var v1301 int32
	_ = v1301
	var v1303 int32
	_ = v1303
	var v1304 int32
	_ = v1304
	var v1308 int32
	_ = v1308
	var v1309 int32
	_ = v1309
	var v1313 int32
	_ = v1313
	var v1314 int32
	_ = v1314
	var v1366 int32
	_ = v1366
	var v1369 int32
	_ = v1369
	var v1374 int32
	_ = v1374
	var v1375 int32
	_ = v1375
	var v1376 int32
	_ = v1376
	var v1378 int32
	_ = v1378
	var v1383 int32
	_ = v1383
	var v1429 int32
	_ = v1429
	var v1430 int32
	_ = v1430
	var v1432 int32
	_ = v1432
	var v1434 int32
	_ = v1434
	var v1439 int32
	_ = v1439
	var v1444 int64
	_ = v1444
	var v1448 int32
	_ = v1448
	var v1452 int32
	_ = v1452
	var v1453 int64
	_ = v1453
	var v1461 int32
	_ = v1461
	var v1462 int64
	_ = v1462
	var v1470 int32
	_ = v1470
	var v1471 int64
	_ = v1471
	var v1479 int32
	_ = v1479
	var v1480 int64
	_ = v1480
	var v1482 int64
	_ = v1482
	var v1484 int64
	_ = v1484
	var v1486 int64
	_ = v1486
	var v1488 int64
	_ = v1488
	var v1491 int32
	_ = v1491
	var v1492 int32
	_ = v1492
	var v1493 int32
	_ = v1493
	var v1494 int32
	_ = v1494
	var v1495 int32
	_ = v1495
	var v1498 int32
	_ = v1498
	var v1502 int32
	_ = v1502
	var v1503 int32
	_ = v1503
	var v1505 int32
	_ = v1505
	var v1506 int32
	_ = v1506
	var v1507 int32
	_ = v1507
	var v1511 int32
	_ = v1511
	var v1512 int32
	_ = v1512
	var v1513 int32
	_ = v1513
	var v1515 int32
	_ = v1515
	var v1516 int32
	_ = v1516
	var v1519 int32
	_ = v1519
	var v1520 int32
	_ = v1520
	var v1521 int32
	_ = v1521
	var v1532 int32
	_ = v1532
	var v1533 int32
	_ = v1533
	var v1535 int32
	_ = v1535
	var v1548 int32
	_ = v1548
	var v1558 int32
	_ = v1558
	var v1559 int32
	_ = v1559
	var v1564 int64
	_ = v1564
	var v1565 int64
	_ = v1565
	var v1572 int32
	_ = v1572
	var v1573 base.V128
	_ = v1573
	var v1574 base.V128
	_ = v1574
	var v1576 int32
	_ = v1576
	var v1580 int32
	_ = v1580
	var v1599 base.V128
	_ = v1599
	var v1608 base.V128
	_ = v1608
	var v1609 base.V128
	_ = v1609
	var v1610 base.V128
	_ = v1610
	var v1611 base.V128
	_ = v1611
	var v1612 base.V128
	_ = v1612
	var v1626 int32
	_ = v1626
	var v1627 int64
	_ = v1627
	var v1631 int32
	_ = v1631
	var v1632 int64
	_ = v1632
	var v1634 base.V128
	_ = v1634
	var v1640 int64
	_ = v1640
	var v1647 int64
	_ = v1647
	var v1650 base.V128
	_ = v1650
	var v1656 int64
	_ = v1656
	var v1663 int64
	_ = v1663
	var v1664 int32
	_ = v1664
	var v1666 base.V128
	_ = v1666
	var v1667 int64
	_ = v1667
	var v1669 int64
	_ = v1669
	var v1671 base.V128
	_ = v1671
	var v1677 int64
	_ = v1677
	var v1684 int64
	_ = v1684
	var v1687 base.V128
	_ = v1687
	var v1693 int64
	_ = v1693
	var v1700 int64
	_ = v1700
	var v1701 int32
	_ = v1701
	var v1703 base.V128
	_ = v1703
	var v1704 int64
	_ = v1704
	var v1706 int64
	_ = v1706
	var v1708 base.V128
	_ = v1708
	var v1714 int64
	_ = v1714
	var v1721 int64
	_ = v1721
	var v1724 base.V128
	_ = v1724
	var v1730 int64
	_ = v1730
	var v1737 int64
	_ = v1737
	var v1740 base.V128
	_ = v1740
	var v1744 int32
	_ = v1744
	var v1746 base.V128
	_ = v1746
	var v1747 base.V128
	_ = v1747
	var v1750 int32
	_ = v1750
	var v1755 int64
	_ = v1755
	var v1756 int32
	_ = v1756
	var v1760 int64
	_ = v1760
	var v1762 base.V128
	_ = v1762
	var v1770 int64
	_ = v1770
	var v1771 int32
	_ = v1771
	var v1775 int64
	_ = v1775
	var v1777 base.V128
	_ = v1777
	var v1785 int64
	_ = v1785
	var v1786 int32
	_ = v1786
	var v1790 int64
	_ = v1790
	var v1792 base.V128
	_ = v1792
	var v1800 int64
	_ = v1800
	var v1801 int32
	_ = v1801
	var v1805 int64
	_ = v1805
	var v1807 base.V128
	_ = v1807
	var v1815 int64
	_ = v1815
	var v1816 int32
	_ = v1816
	var v1820 int64
	_ = v1820
	var v1822 base.V128
	_ = v1822
	var v1830 int64
	_ = v1830
	var v1831 int32
	_ = v1831
	var v1835 int64
	_ = v1835
	var v1838 int32
	_ = v1838
	var v1859 int32
	_ = v1859
	var v1868 int64
	_ = v1868
	var v1869 int64
	_ = v1869
	var v1870 int64
	_ = v1870
	var v1871 int64
	_ = v1871
	var v1872 int64
	_ = v1872
	var v1873 int64
	_ = v1873
	var v1897 int32
	_ = v1897
	var v1902 int32
	_ = v1902
	var v1923 int64
	_ = v1923
	var v1924 int64
	_ = v1924
	var v1925 int64
	_ = v1925
	var v1926 int64
	_ = v1926
	var v1927 int64
	_ = v1927
	var v1928 int64
	_ = v1928
	var v1947 int32
	_ = v1947
	var v1948 int64
	_ = v1948
	var v1950 int64
	_ = v1950
	var v1952 int64
	_ = v1952
	var v1953 int64
	_ = v1953
	var v1955 int64
	_ = v1955
	var v1957 int64
	_ = v1957
	var v1958 int64
	_ = v1958
	var v1960 int64
	_ = v1960
	var v1962 int64
	_ = v1962
	var v1966 int32
	_ = v1966
	var v1994 int64
	_ = v1994
	var v1995 int64
	_ = v1995
	var v1996 int64
	_ = v1996
	var v1997 int64
	_ = v1997
	var v1998 int64
	_ = v1998
	var v1999 int64
	_ = v1999
	var v2019 float64
	_ = v2019
	var v2040 int32
	_ = v2040
	var v2043 int32
	_ = v2043
	var v2076 int32
	_ = v2076
	var v2077 int32
	_ = v2077
	var v2079 int64
	_ = v2079
	var v2084 float64
	_ = v2084
	var v2090 int32
	_ = v2090
	var v2092 int32
	_ = v2092
	var v2095 int64
	_ = v2095
	var v2100 float64
	_ = v2100
	var v2106 int32
	_ = v2106
	var v2108 int32
	_ = v2108
	var v2111 int32
	_ = v2111
	var v2118 int64
	_ = v2118
	var v2123 float64
	_ = v2123
	var v2129 int32
	_ = v2129
	var v2131 int32
	_ = v2131
	var v2133 int32
	_ = v2133
	var v2138 int32
	_ = v2138
	var v2140 int32
	_ = v2140
	var v2149 int32
	_ = v2149
	var v2163 int32
	_ = v2163
	var v2199 int32
	_ = v2199
	var v2200 int64
	_ = v2200
	var v2205 float64
	_ = v2205
	var v2211 int32
	_ = v2211
	var v2213 int32
	_ = v2213
	var v2217 int32
	_ = v2217
	var v2218 int64
	_ = v2218
	var v2223 float64
	_ = v2223
	var v2229 int32
	_ = v2229
	var v2231 int32
	_ = v2231
	var v2236 int32
	_ = v2236
	var v2253 int32
	_ = v2253
	var v2294 int32
	_ = v2294
	var v2295 int64
	_ = v2295
	var v2300 float64
	_ = v2300
	var v2306 int32
	_ = v2306
	var v2308 int32
	_ = v2308
	var v2314 int32
	_ = v2314
	var v2316 int32
	_ = v2316
	var v2330 int32
	_ = v2330
	var v2368 int32
	_ = v2368
	var v2371 int32
	_ = v2371
	var v2374 int32
	_ = v2374
	var v2375 int32
	_ = v2375
	var v2376 int32
	_ = v2376
	var v2387 int32
	_ = v2387
	var v2408 int32
	_ = v2408
	var v2445 int32
	_ = v2445
	var v2459 int32
	_ = v2459
	var v2495 int32
	_ = v2495
	var v2496 int32
	_ = v2496
	var v2501 int32
	_ = v2501
	var v2557 int32
	_ = v2557
	var v2607 int32
	_ = v2607
	var v2608 int64
	_ = v2608
	var v2615 int32
	_ = v2615
	var v2623 int32
	_ = v2623
	var v2642 int32
	_ = v2642
	var v2673 int32
	_ = v2673
	var v2674 int32
	_ = v2674
	var v2676 int32
	_ = v2676
	var v2677 int32
	_ = v2677
	var v2680 int32
	_ = v2680
	var v2681 int32
	_ = v2681
	var v2687 int32
	_ = v2687
	var v2690 int32
	_ = v2690
	var v2692 int32
	_ = v2692
	var v2694 int32
	_ = v2694
	var v2696 int32
	_ = v2696
	var v2703 int32
	_ = v2703
	var v2705 int32
	_ = v2705
	var v2708 int32
	_ = v2708
	var v2713 int32
	_ = v2713
	var v2720 int32
	_ = v2720
	var v2763 int32
	_ = v2763
	var v2764 int32
	_ = v2764
	var v2765 int32
	_ = v2765
	var v2766 int32
	_ = v2766
	var v2767 int32
	_ = v2767
	var v2770 int32
	_ = v2770
	var v2771 int32
	_ = v2771
	var v2776 int32
	_ = v2776
	var v2777 int32
	_ = v2777
	var v2781 int64
	_ = v2781
	var v2782 int64
	_ = v2782
	var v2785 int64
	_ = v2785
	var v2786 int64
	_ = v2786
	var v2788 int64
	_ = v2788
	var v2789 int64
	_ = v2789
	var v2790 int64
	_ = v2790
	var v2791 int64
	_ = v2791
	var v2795 int64
	_ = v2795
	var v2798 int64
	_ = v2798
	var v2803 int32
	_ = v2803
	var v2807 int64
	_ = v2807
	var v2810 base.V128
	_ = v2810
	var v2814 base.V128
	_ = v2814
	var v2817 int64
	_ = v2817
	var v2819 int32
	_ = v2819
	var v2822 int32
	_ = v2822
	var v2825 int32
	_ = v2825
	var v2828 int32
	_ = v2828
	var v2829 int32
	_ = v2829
	var v2831 int32
	_ = v2831
	var v2832 int32
	_ = v2832
	var v2836 int32
	_ = v2836
	var v2839 int32
	_ = v2839
	var v2842 int32
	_ = v2842
	var v2843 int32
	_ = v2843
	var v2844 int32
	_ = v2844
	var v2847 int32
	_ = v2847
	var v2850 int32
	_ = v2850
	var v2853 int32
	_ = v2853
	var v2855 int32
	_ = v2855
	var v2859 int32
	_ = v2859
	var v2873 int32
	_ = v2873
	var v2876 int32
	_ = v2876
	var v2880 int32
	_ = v2880
	var v2887 int32
	_ = v2887
	var v2889 int32
	_ = v2889
	var v2892 int32
	_ = v2892
	var v2900 int32
	_ = v2900
	var v2922 int32
	_ = v2922
	var v2958 int32
	_ = v2958
	var v2973 int32
	_ = v2973
	var v3009 int32
	_ = v3009
	var v3010 int32
	_ = v3010
	var v3012 int32
	_ = v3012
	var v3014 int32
	_ = v3014
	var v3019 int32
	_ = v3019
	var v3024 int64
	_ = v3024
	var v3026 int32
	_ = v3026
	var v3028 int32
	_ = v3028
	var v3033 int64
	_ = v3033
	var v3042 int64
	_ = v3042
	var v3051 int64
	_ = v3051
	var v3060 int64
	_ = v3060
	var v3062 int64
	_ = v3062
	var v3064 int64
	_ = v3064
	var v3066 int64
	_ = v3066
	var v3068 int64
	_ = v3068
	var v3074 int32
	_ = v3074
	var v3075 int32
	_ = v3075
	var v3097 int32
	_ = v3097
	var v3148 int32
	_ = v3148
	var v3180 int32
	_ = v3180
	var v3181 int32
	_ = v3181
	var v3185 int32
	_ = v3185
	var v3188 int64
	_ = v3188
	var v3194 int64
	_ = v3194
	var v3195 int32
	_ = v3195
	var v3197 int64
	_ = v3197
	var v3198 int32
	_ = v3198
	var v3205 int64
	_ = v3205
	var v3206 int32
	_ = v3206
	var v3207 int64
	_ = v3207
	var v3217 int32
	_ = v3217
	var v3219 int32
	_ = v3219
	var v3223 int32
	_ = v3223
	var v3228 int32
	_ = v3228
	var v3232 int32
	_ = v3232
	var v3236 int32
	_ = v3236
	var v3237 int32
	_ = v3237
	var v3241 int32
	_ = v3241
	var v3263 int32
	_ = v3263
	var v3288 int32
	_ = v3288
	var v3289 int32
	_ = v3289
	var v3292 int32
	_ = v3292
	var v3295 int32
	_ = v3295
	var v3298 int64
	_ = v3298
	var v3299 int64
	_ = v3299
	var v3300 int32
	_ = v3300
	var v3301 int32
	_ = v3301
	var v3308 int32
	_ = v3308
	var v3313 int32
	_ = v3313
	var v3314 int32
	_ = v3314
	var v3330 int32
	_ = v3330
	var v3334 int64
	_ = v3334
	var v3363 int64
	_ = v3363
	var v3364 int32
	_ = v3364
	var v3367 int32
	_ = v3367
	var v3369 int32
	_ = v3369
	var v3370 int32
	_ = v3370
	var v3372 int32
	_ = v3372
	var v3374 int32
	_ = v3374
	var v3376 int32
	_ = v3376
	var v3381 int32
	_ = v3381
	var v3382 int64
	_ = v3382
	var v3384 int32
	_ = v3384
	var v3389 int32
	_ = v3389
	var v3390 int64
	_ = v3390
	var v3391 int64
	_ = v3391
	var v3392 int64
	_ = v3392
	var v3393 int64
	_ = v3393
	var v3397 int64
	_ = v3397
	var v3400 int64
	_ = v3400
	var v3401 int32
	_ = v3401
	var v3404 int32
	_ = v3404
	var v3406 int32
	_ = v3406
	var v3413 int32
	_ = v3413
	var v3414 base.V128
	_ = v3414
	var v3424 base.V128
	_ = v3424
	var v3429 int64
	_ = v3429
	var v3433 int64
	_ = v3433
	var v3437 int64
	_ = v3437
	var v3439 int64
	_ = v3439
	var v3447 int32
	_ = v3447
	var v3450 int32
	_ = v3450
	var v3453 int64
	_ = v3453
	var v3454 int64
	_ = v3454
	var v3456 int32
	_ = v3456
	var v3457 base.V128
	_ = v3457
	var v3459 int32
	_ = v3459
	var v3461 base.V128
	_ = v3461
	var v3464 int32
	_ = v3464
	var v3465 int32
	_ = v3465
	var v3467 base.V128
	_ = v3467
	var v3469 int32
	_ = v3469
	var v3471 base.V128
	_ = v3471
	var v3474 int32
	_ = v3474
	var v3475 int32
	_ = v3475
	var v3477 base.V128
	_ = v3477
	var v3479 int32
	_ = v3479
	var v3481 base.V128
	_ = v3481
	var v3484 int32
	_ = v3484
	var v3485 int32
	_ = v3485
	var v3487 base.V128
	_ = v3487
	var v3489 int32
	_ = v3489
	var v3491 base.V128
	_ = v3491
	var v3495 int32
	_ = v3495
	var v3499 int32
	_ = v3499
	var v3503 int32
	_ = v3503
	var v3508 int64
	_ = v3508
	var v3510 int64
	_ = v3510
	var v3516 int64
	_ = v3516
	var v3518 int64
	_ = v3518
	var v3520 int64
	_ = v3520
	var v3526 int64
	_ = v3526
	var v3532 int64
	_ = v3532
	var v3538 int64
	_ = v3538
	var v3553 int64
	_ = v3553
	var v3559 int32
	_ = v3559
	var v3567 int64
	_ = v3567
	var v3574 int32
	_ = v3574
	var v3580 int32
	_ = v3580
	var v3602 int32
	_ = v3602
	var v3632 int32
	_ = v3632
	var v3654 int32
	_ = v3654
	var v3683 int32
	_ = v3683
	var v3684 int32
	_ = v3684
	var v3685 int32
	_ = v3685
	var v3687 int32
	_ = v3687
	var v3688 int32
	_ = v3688
	var v3691 int32
	_ = v3691
	var v3692 int32
	_ = v3692
	var v3694 int32
	_ = v3694
	var v3695 int64
	_ = v3695
	var v3697 int64
	_ = v3697
	var v3699 int64
	_ = v3699
	var v3701 int64
	_ = v3701
	var v3703 int64
	_ = v3703
	var v3705 int32
	_ = v3705
	var v3709 int32
	_ = v3709
	var v3716 int32
	_ = v3716
	var v3718 int64
	_ = v3718
	var v3723 int32
	_ = v3723
	var v3742 int32
	_ = v3742
	var v3776 int32
	_ = v3776
	var v3777 int32
	_ = v3777
	var v3780 int32
	_ = v3780
	var v3781 int32
	_ = v3781
	var v3784 int32
	_ = v3784
	var v3791 int32
	_ = v3791
	var v3792 int32
	_ = v3792
	var v3793 base.V128
	_ = v3793
	var v3796 int32
	_ = v3796
	var v3801 base.V128
	_ = v3801
	var v3804 int32
	_ = v3804
	var v3809 base.V128
	_ = v3809
	var v3812 int32
	_ = v3812
	var v3817 base.V128
	_ = v3817
	var v3827 int32
	_ = v3827
	var v3830 int32
	_ = v3830
	var v3834 int32
	_ = v3834
	var v3835 int32
	_ = v3835
	var v3836 int32
	_ = v3836
	var v3839 int32
	_ = v3839
	var v3843 int32
	_ = v3843
	var v3844 int64
	_ = v3844
	var v3845 int64
	_ = v3845
	var v3846 int64
	_ = v3846
	var v3848 int32
	_ = v3848
	var v3851 int32
	_ = v3851
	var v3854 int64
	_ = v3854
	var v3857 int32
	_ = v3857
	var v3858 int32
	_ = v3858
	var v3862 int32
	_ = v3862
	var v3865 int32
	_ = v3865
	var v3866 int32
	_ = v3866
	var v3869 int32
	_ = v3869
	var v3873 int64
	_ = v3873
	var v3874 int64
	_ = v3874
	var v3876 int32
	_ = v3876
	var v3877 base.V128
	_ = v3877
	var v3879 base.V128
	_ = v3879
	var v3882 int32
	_ = v3882
	var v3883 int32
	_ = v3883
	var v3885 base.V128
	_ = v3885
	var v3887 int32
	_ = v3887
	var v3889 base.V128
	_ = v3889
	var v3892 int32
	_ = v3892
	var v3893 int32
	_ = v3893
	var v3895 base.V128
	_ = v3895
	var v3897 int32
	_ = v3897
	var v3899 base.V128
	_ = v3899
	var v3902 int32
	_ = v3902
	var v3903 int32
	_ = v3903
	var v3905 base.V128
	_ = v3905
	var v3907 int32
	_ = v3907
	var v3909 base.V128
	_ = v3909
	var v3913 int32
	_ = v3913
	var v3917 int32
	_ = v3917
	var v3921 int32
	_ = v3921
	var v3926 int64
	_ = v3926
	var v3928 int32
	_ = v3928
	var v3932 int64
	_ = v3932
	var v3934 int64
	_ = v3934
	var v3936 int32
	_ = v3936
	var v3940 int64
	_ = v3940
	var v3942 int64
	_ = v3942
	var v3944 int64
	_ = v3944
	var v3946 int32
	_ = v3946
	var v3950 int64
	_ = v3950
	var v3952 int32
	_ = v3952
	var v3956 int64
	_ = v3956
	var v3972 int32
	_ = v3972
	var v3973 int32
	_ = v3973
	var v3974 base.V128
	_ = v3974
	var v3977 int32
	_ = v3977
	var v3982 base.V128
	_ = v3982
	var v3985 int32
	_ = v3985
	var v3990 base.V128
	_ = v3990
	var v3996 base.V128
	_ = v3996
	var v4001 int32
	_ = v4001
	var v4006 int32
	_ = v4006
	var v4015 int32
	_ = v4015
	var v4067 int32
	_ = v4067
	var v4089 int32
	_ = v4089
	var v4115 int32
	_ = v4115
	var v4119 int32
	_ = v4119
	var v4172 int32
	_ = v4172
	var v4174 int32
	_ = v4174
	var v4176 int32
	_ = v4176
	var v4197 int32
	_ = v4197
	var v4198 int32
	_ = v4198
	var v4228 int32
	_ = v4228
	var v4231 int64
	_ = v4231
	var v4232 int32
	_ = v4232
	var v4239 int64
	_ = v4239
	var v4240 int32
	_ = v4240
	var v4241 int64
	_ = v4241
	var v4251 int32
	_ = v4251
	var v4253 int32
	_ = v4253
	var v4259 int32
	_ = v4259
	var v4262 int32
	_ = v4262
	var v4311 int32
	_ = v4311
	var v4313 int32
	_ = v4313
	var v4316 int32
	_ = v4316
	var v4337 int32
	_ = v4337
	var v4368 int32
	_ = v4368
	var v4373 int32
	_ = v4373
	var v4374 int64
	_ = v4374
	var v4376 int32
	_ = v4376
	var v4381 int32
	_ = v4381
	var v4382 int64
	_ = v4382
	var v4383 int64
	_ = v4383
	var v4386 int32
	_ = v4386
	var v4390 int32
	_ = v4390
	var v4394 int32
	_ = v4394
	var v4398 int32
	_ = v4398
	var v4406 int32
	_ = v4406
	var v4407 int32
	_ = v4407
	var v4421 int32
	_ = v4421
	var v4436 int64
	_ = v4436
	var v4461 int32
	_ = v4461
	var v4466 int32
	_ = v4466
	var v4468 int32
	_ = v4468
	var v4469 int32
	_ = v4469
	var v4479 int32
	_ = v4479
	var v4487 int64
	_ = v4487
	var v4489 int64
	_ = v4489
	var v4490 int32
	_ = v4490
	var v4492 int32
	_ = v4492
	var v4494 int32
	_ = v4494
	var v4499 int32
	_ = v4499
	var v4500 int32
	_ = v4500
	var v4501 int32
	_ = v4501
	var v4503 int32
	_ = v4503
	var v4504 int32
	_ = v4504
	var v4505 int32
	_ = v4505
	var v4508 int32
	_ = v4508
	var v4509 int32
	_ = v4509
	var v4511 int32
	_ = v4511
	var v4523 int64
	_ = v4523
	var v4526 int64
	_ = v4526
	var v4527 int64
	_ = v4527
	var v4530 int64
	_ = v4530
	var v4533 int64
	_ = v4533
	var v4535 int64
	_ = v4535
	var v4536 int64
	_ = v4536
	var v4539 int32
	_ = v4539
	var v4542 int32
	_ = v4542
	var v4551 int64
	_ = v4551
	var v4554 int64
	_ = v4554
	var v4557 int64
	_ = v4557
	var v4559 int64
	_ = v4559
	var v4560 int64
	_ = v4560
	var v4562 int32
	_ = v4562
	var v4565 int32
	_ = v4565
	var v4569 int32
	_ = v4569
	var v4573 int32
	_ = v4573
	var v4577 int32
	_ = v4577
	var v4581 int32
	_ = v4581
	var v4594 int64
	_ = v4594
	var v4602 int64
	_ = v4602
	var v4609 int32
	_ = v4609
	var v4612 int32
	_ = v4612
	var v4614 int32
	_ = v4614
	var v4621 int32
	_ = v4621
	var v4622 base.V128
	_ = v4622
	var v4632 base.V128
	_ = v4632
	var v4641 int64
	_ = v4641
	var v4648 int64
	_ = v4648
	var v4660 int64
	_ = v4660
	var v4663 int32
	_ = v4663
	var v4666 int32
	_ = v4666
	var v4669 int64
	_ = v4669
	var v4670 int64
	_ = v4670
	var v4672 int32
	_ = v4672
	var v4673 base.V128
	_ = v4673
	var v4675 int32
	_ = v4675
	var v4677 base.V128
	_ = v4677
	var v4680 int32
	_ = v4680
	var v4681 int32
	_ = v4681
	var v4683 base.V128
	_ = v4683
	var v4685 int32
	_ = v4685
	var v4687 base.V128
	_ = v4687
	var v4690 int32
	_ = v4690
	var v4691 int32
	_ = v4691
	var v4693 base.V128
	_ = v4693
	var v4695 int32
	_ = v4695
	var v4697 base.V128
	_ = v4697
	var v4700 int32
	_ = v4700
	var v4701 int32
	_ = v4701
	var v4703 base.V128
	_ = v4703
	var v4705 int32
	_ = v4705
	var v4707 base.V128
	_ = v4707
	var v4711 int32
	_ = v4711
	var v4715 int32
	_ = v4715
	var v4719 int32
	_ = v4719
	var v4724 int64
	_ = v4724
	var v4726 int64
	_ = v4726
	var v4732 int64
	_ = v4732
	var v4734 int64
	_ = v4734
	var v4736 int64
	_ = v4736
	var v4742 int64
	_ = v4742
	var v4748 int64
	_ = v4748
	var v4754 int64
	_ = v4754
	var v4757 int32
	_ = v4757
	var v4808 int32
	_ = v4808
	var v4811 int32
	_ = v4811
	var v4867 int32
	_ = v4867
	var v4917 int32
	_ = v4917
	var v4918 int32
	_ = v4918
	var v4919 int32
	_ = v4919
	var v4921 int32
	_ = v4921
	var v4922 int32
	_ = v4922
	var v4925 int32
	_ = v4925
	var v4926 int32
	_ = v4926
	var v4928 int32
	_ = v4928
	var v4929 int64
	_ = v4929
	var v4931 int64
	_ = v4931
	var v4933 int64
	_ = v4933
	var v4935 int64
	_ = v4935
	var v4937 int64
	_ = v4937
	var v4939 int32
	_ = v4939
	var v4943 int32
	_ = v4943
	var v4950 int32
	_ = v4950
	var v4952 int64
	_ = v4952
	var v4956 int32
	_ = v4956
	var v4961 int32
	_ = v4961
	var v5008 int32
	_ = v5008
	var v5009 int32
	_ = v5009
	var v5011 int32
	_ = v5011
	var v5018 int32
	_ = v5018
	var v5019 int32
	_ = v5019
	var v5020 base.V128
	_ = v5020
	var v5023 int32
	_ = v5023
	var v5028 base.V128
	_ = v5028
	var v5031 int32
	_ = v5031
	var v5036 base.V128
	_ = v5036
	var v5039 int32
	_ = v5039
	var v5044 base.V128
	_ = v5044
	var v5049 int32
	_ = v5049
	var v5052 int32
	_ = v5052
	var v5055 int32
	_ = v5055
	var v5059 int64
	_ = v5059
	var v5060 int64
	_ = v5060
	var v5062 int32
	_ = v5062
	var v5063 base.V128
	_ = v5063
	var v5065 base.V128
	_ = v5065
	var v5068 int32
	_ = v5068
	var v5069 int32
	_ = v5069
	var v5071 base.V128
	_ = v5071
	var v5073 int32
	_ = v5073
	var v5075 base.V128
	_ = v5075
	var v5078 int32
	_ = v5078
	var v5079 int32
	_ = v5079
	var v5081 base.V128
	_ = v5081
	var v5083 int32
	_ = v5083
	var v5085 base.V128
	_ = v5085
	var v5088 int32
	_ = v5088
	var v5089 int32
	_ = v5089
	var v5091 base.V128
	_ = v5091
	var v5093 int32
	_ = v5093
	var v5095 base.V128
	_ = v5095
	var v5099 int32
	_ = v5099
	var v5103 int32
	_ = v5103
	var v5107 int32
	_ = v5107
	var v5112 int64
	_ = v5112
	var v5114 int32
	_ = v5114
	var v5118 int64
	_ = v5118
	var v5120 int64
	_ = v5120
	var v5122 int32
	_ = v5122
	var v5126 int64
	_ = v5126
	var v5128 int64
	_ = v5128
	var v5130 int64
	_ = v5130
	var v5132 int32
	_ = v5132
	var v5136 int64
	_ = v5136
	var v5138 int32
	_ = v5138
	var v5142 int64
	_ = v5142
	var v5156 int32
	_ = v5156
	var v5158 int32
	_ = v5158
	var v5169 int32
	_ = v5169
	var v5173 int32
	_ = v5173
	var v5181 int32
	_ = v5181
	var v5195 int32
	_ = v5195
	var v5226 int32
	_ = v5226
	var v5228 int32
	_ = v5228
	var v5233 int32
	_ = v5233
	var v5234 int64
	_ = v5234
	var v5236 int32
	_ = v5236
	var v5241 int32
	_ = v5241
	var v5242 int64
	_ = v5242
	var v5243 int64
	_ = v5243
	var v5246 int32
	_ = v5246
	var v5250 int32
	_ = v5250
	var v5254 int32
	_ = v5254
	var v5258 int32
	_ = v5258
	var v5266 int32
	_ = v5266
	var v5267 int32
	_ = v5267
	var v5281 int32
	_ = v5281
	var v5296 int64
	_ = v5296
	var v5321 int32
	_ = v5321
	var v5326 int32
	_ = v5326
	var v5328 int32
	_ = v5328
	var v5329 int32
	_ = v5329
	var v5339 int32
	_ = v5339
	var v5347 int64
	_ = v5347
	var v5349 int64
	_ = v5349
	var v5350 int32
	_ = v5350
	var v5352 int32
	_ = v5352
	var v5354 int32
	_ = v5354
	var v5359 int32
	_ = v5359
	var v5360 int32
	_ = v5360
	var v5361 int32
	_ = v5361
	var v5363 int32
	_ = v5363
	var v5364 int32
	_ = v5364
	var v5365 int32
	_ = v5365
	var v5368 int32
	_ = v5368
	var v5369 int32
	_ = v5369
	var v5371 int32
	_ = v5371
	var v5383 int64
	_ = v5383
	var v5386 int64
	_ = v5386
	var v5387 int64
	_ = v5387
	var v5390 int64
	_ = v5390
	var v5393 int64
	_ = v5393
	var v5395 int64
	_ = v5395
	var v5396 int64
	_ = v5396
	var v5399 int32
	_ = v5399
	var v5402 int32
	_ = v5402
	var v5411 int64
	_ = v5411
	var v5414 int64
	_ = v5414
	var v5417 int64
	_ = v5417
	var v5419 int64
	_ = v5419
	var v5420 int64
	_ = v5420
	var v5422 int32
	_ = v5422
	var v5425 int32
	_ = v5425
	var v5429 int32
	_ = v5429
	var v5433 int32
	_ = v5433
	var v5437 int32
	_ = v5437
	var v5441 int32
	_ = v5441
	var v5454 int64
	_ = v5454
	var v5462 int64
	_ = v5462
	var v5469 int32
	_ = v5469
	var v5472 int32
	_ = v5472
	var v5474 int32
	_ = v5474
	var v5481 int32
	_ = v5481
	var v5482 base.V128
	_ = v5482
	var v5492 base.V128
	_ = v5492
	var v5501 int64
	_ = v5501
	var v5508 int64
	_ = v5508
	var v5520 int64
	_ = v5520
	var v5523 int32
	_ = v5523
	var v5526 int32
	_ = v5526
	var v5529 int64
	_ = v5529
	var v5530 int64
	_ = v5530
	var v5532 int32
	_ = v5532
	var v5533 base.V128
	_ = v5533
	var v5535 int32
	_ = v5535
	var v5537 base.V128
	_ = v5537
	var v5540 int32
	_ = v5540
	var v5541 int32
	_ = v5541
	var v5543 base.V128
	_ = v5543
	var v5545 int32
	_ = v5545
	var v5547 base.V128
	_ = v5547
	var v5550 int32
	_ = v5550
	var v5551 int32
	_ = v5551
	var v5553 base.V128
	_ = v5553
	var v5555 int32
	_ = v5555
	var v5557 base.V128
	_ = v5557
	var v5560 int32
	_ = v5560
	var v5561 int32
	_ = v5561
	var v5563 base.V128
	_ = v5563
	var v5565 int32
	_ = v5565
	var v5567 base.V128
	_ = v5567
	var v5571 int32
	_ = v5571
	var v5575 int32
	_ = v5575
	var v5579 int32
	_ = v5579
	var v5584 int64
	_ = v5584
	var v5586 int64
	_ = v5586
	var v5592 int64
	_ = v5592
	var v5594 int64
	_ = v5594
	var v5596 int64
	_ = v5596
	var v5602 int64
	_ = v5602
	var v5608 int64
	_ = v5608
	var v5614 int64
	_ = v5614
	var v5617 int32
	_ = v5617
	var v5667 int32
	_ = v5667
	var v5668 int32
	_ = v5668
	var v5669 int32
	_ = v5669
	var v5676 int32
	_ = v5676
	var v5720 int32
	_ = v5720
	var v5723 int32
	_ = v5723
	var v5827 int32
	_ = v5827
	var v5848 int32
	_ = v5848
	var v5879 int32
	_ = v5879
	var v5880 int32
	_ = v5880
	var v5881 int32
	_ = v5881
	var v5886 int32
	_ = v5886
	var v5888 int32
	_ = v5888
	var v5899 int32
	_ = v5899
	var v5921 int32
	_ = v5921
	var v5924 int32
	_ = v5924
	var v5925 int32
	_ = v5925
	var v5933 int32
	_ = v5933
	var v5934 int32
	_ = v5934
	var v5969 int32
	_ = v5969
	var v5970 int32
	_ = v5970
	var v5975 int64
	_ = v5975
	var v5978 int32
	_ = v5978
	var v5979 int32
	_ = v5979
	var v5991 int32
	_ = v5991
	var v6037 int32
	_ = v6037
	var v6064 int32
	_ = v6064
	var v6065 int32
	_ = v6065
	var v6066 int32
	_ = v6066
	var v6070 int32
	_ = v6070
	var v6086 int32
	_ = v6086
	var v6101 int64
	_ = v6101
	var v6117 int32
	_ = v6117
	var v6123 int64
	_ = v6123
	var v6127 int32
	_ = v6127
	var v6128 int64
	_ = v6128
	var v6129 int64
	_ = v6129
	var v6133 int64
	_ = v6133
	var v6136 int64
	_ = v6136
	var v6139 int32
	_ = v6139
	var v6140 int32
	_ = v6140
	var v6144 int32
	_ = v6144
	var v6148 int32
	_ = v6148
	var v6152 int32
	_ = v6152
	var v6160 int32
	_ = v6160
	var v6161 int32
	_ = v6161
	var v6166 int32
	_ = v6166
	var v6186 int64
	_ = v6186
	var v6215 int32
	_ = v6215
	var v6220 int32
	_ = v6220
	var v6222 int32
	_ = v6222
	var v6223 int32
	_ = v6223
	var v6233 int32
	_ = v6233
	var v6239 int64
	_ = v6239
	var v6241 int64
	_ = v6241
	var v6242 int32
	_ = v6242
	var v6245 int32
	_ = v6245
	var v6247 int32
	_ = v6247
	var v6252 int32
	_ = v6252
	var v6253 int32
	_ = v6253
	var v6254 int32
	_ = v6254
	var v6255 int32
	_ = v6255
	var v6256 int32
	_ = v6256
	var v6257 int32
	_ = v6257
	var v6260 int32
	_ = v6260
	var v6261 int32
	_ = v6261
	var v6263 int32
	_ = v6263
	var v6275 int64
	_ = v6275
	var v6278 int64
	_ = v6278
	var v6279 int64
	_ = v6279
	var v6282 int64
	_ = v6282
	var v6285 int64
	_ = v6285
	var v6287 int64
	_ = v6287
	var v6288 int64
	_ = v6288
	var v6291 int32
	_ = v6291
	var v6294 int32
	_ = v6294
	var v6303 int64
	_ = v6303
	var v6306 int64
	_ = v6306
	var v6309 int64
	_ = v6309
	var v6311 int64
	_ = v6311
	var v6313 int64
	_ = v6313
	var v6314 int32
	_ = v6314
	var v6317 int32
	_ = v6317
	var v6321 int32
	_ = v6321
	var v6325 int32
	_ = v6325
	var v6329 int32
	_ = v6329
	var v6333 int32
	_ = v6333
	var v6347 int64
	_ = v6347
	var v6348 int64
	_ = v6348
	var v6355 int32
	_ = v6355
	var v6358 int64
	_ = v6358
	var v6389 int64
	_ = v6389
	var v6405 int32
	_ = v6405
	var v6412 int32
	_ = v6412
	var v6416 int32
	_ = v6416
	var v6419 int32
	_ = v6419
	var v6473 int32
	_ = v6473
	var v6475 int32
	_ = v6475
	var v6477 int32
	_ = v6477
	var v6527 int32
	_ = v6527
	var v6529 int32
	_ = v6529
	var v6550 int32
	_ = v6550
	var v6579 int32
	_ = v6579
	var v6582 int32
	_ = v6582
	var v6583 int32
	_ = v6583
	var v6584 int32
	_ = v6584
	var v6585 int32
	_ = v6585
	var v6589 int32
	_ = v6589
	var v6590 int32
	_ = v6590
	var v6593 int32
	_ = v6593
	var v6604 int32
	_ = v6604
	var v6626 int32
	_ = v6626
	var v6629 int32
	_ = v6629
	var v6630 int32
	_ = v6630
	var v6638 int32
	_ = v6638
	var v6639 int32
	_ = v6639
	var v6674 int32
	_ = v6674
	var v6675 int32
	_ = v6675
	var v6680 int64
	_ = v6680
	var v6683 int32
	_ = v6683
	var v6684 int32
	_ = v6684
	var v6696 int32
	_ = v6696
	var v6716 int32
	_ = v6716
	var v6720 int32
	_ = v6720
	var v6724 int32
	_ = v6724
	var v6727 int32
	_ = v6727
	var v6730 int32
	_ = v6730
	var v6736 int32
	_ = v6736
	var v6740 int32
	_ = v6740
	var v6750 int32
	_ = v6750
	var v6753 int32
	_ = v6753
	var v6787 int32
	_ = v6787
	var v6789 int32
	_ = v6789
	var v6791 int32
	_ = v6791
	var v6792 int32
	_ = v6792
	var v6794 int32
	_ = v6794
	var v6795 int32
	_ = v6795
	var v6798 int32
	_ = v6798
	var v6799 int32
	_ = v6799
	var v6806 int32
	_ = v6806
	var v6808 int32
	_ = v6808
	var v6810 int32
	_ = v6810
	var v6816 int32
	_ = v6816
	var v6818 int32
	_ = v6818
	var v6825 int32
	_ = v6825
	var v6835 int32
	_ = v6835
	var v6838 int32
	_ = v6838
	var v6874 int32
	_ = v6874
	var v6879 int32
	_ = v6879
	var v6881 int32
	_ = v6881
	var v6883 int32
	_ = v6883
	var v6888 int32
	_ = v6888
	var v6894 int32
	_ = v6894
	var v6896 int32
	_ = v6896
	var v6910 int32
	_ = v6910
	var v6948 int32
	_ = v6948
	var v6951 int32
	_ = v6951
	var v6954 int32
	_ = v6954
	var v6955 int32
	_ = v6955
	var v6963 int32
	_ = v6963
	var v6988 int32
	_ = v6988
	var v7021 int32
	_ = v7021
	var v7035 int32
	_ = v7035
	var v7071 int32
	_ = v7071
	var v7076 int32
	_ = v7076
	var v7140 int32
	_ = v7140
	var v7175 int32
	_ = v7175
	var v7177 int32
	_ = v7177
	var v7182 int32
	_ = v7182
	var v7185 int32
	_ = v7185
	var v7189 int32
	_ = v7189
	var v7192 int32
	_ = v7192
	var v7197 int32
	_ = v7197
	var v7249 int32
	_ = v7249
	var v7251 int32
	_ = v7251
	var v7255 int32
	_ = v7255
	var v7261 int32
	_ = v7261
	var v7316 int32
	_ = v7316
	v52 = m.G0
	v54 = v52 - int32(352)
	m.G0 = v54
	if l5 == int32(0) {
		v70 = int32(1)
		goto L1
	} else {
		goto L2
	}
L1:
	;
	v73 = int32(4) << (uint(l6) % 32)
	v74 = int32(_a_F_VP8LGetHistoImageSymbols_0)
	v78 = base.B2i32(int32(0) < l6)
	if int32(0) < l6 {
		goto L12
	} else {
		goto L13
	}
L2:
	;
	v60 = int32(1) << (uint(l5) % 32)
	v62 = int32(-1)
	v70 = int32(base.Ui32(l0+v60+v62)>>(uint(l5)%32)) * int32(base.Ui32(l1+v60+v62)>>(uint(l5)%32))
	goto L1
L3:
	;
	F_free(m, v105)
	mBase = m.M
	goto L677
L4:
	;
	F_free(m, int32(0))
	mBase = m.M
	goto L673
L5:
	;
	F_free(m, int32(0))
	mBase = m.M
	goto L669
L6:
	;
	v5879 = *(*int32)(unsafe.Add(mBase, uint32(l7)+4))
	v5880 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v5881 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	if int32(1) < v5848 {
		goto L550
	} else {
		goto L551
	}
L7:
	;
	v3180 = base.I32_div_s(v3148, int32(2))
	v3181 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v3185 = l3 * l3 * l3
	if v3185 < int32(0) {
		goto L314
	} else {
		goto L315
	}
L8:
	;
	v2557 = int32(0)
	goto L260
L9:
	;
	v2019 = base.F64_convert_i64_u(v1995 - v1994)
	if l4 != 0 {
		goto L216
	} else {
		goto L217
	}
L10:
	;
	v1897 = v1559 + v1859<<(uint(int32(2))%32)
	v1902 = v1535 - v1859
	v1923 = v1868
	v1924 = v1869
	v1925 = v1870
	v1926 = v1871
	v1927 = v1872
	v1928 = v1873
	goto L195
L11:
	;
	v1838 = *(*int32)(unsafe.Add(mBase, uint32(l10)+92))
	if v1838 != 0 {
		goto L193
	} else {
		goto L194
	}
L12:
	;
	v79 = v73 + v74
	goto L14
L13:
	;
	v79 = v74
	goto L14
L14:
	;
	v83 = base.I64_extend_i32_u(v70*v79 + int32(12))
	v84 = int32(1)
	if v83 == int64(0) {
		goto L17
	} else {
		goto L18
	}
L15:
	;
	if v105 == int32(0) {
		goto L11
	} else {
		goto L21
	}
L16:
	;
	goto L15
L17:
	;
	v103 = F_malloc(m, base.I32_wrap_i64(v83)*v84)
	mBase = m.M
	v105 = v103
	goto L16
L18:
	;
	v91 = base.I64_div_u_s(int64(2147418112), v83)
	v92 = int32(0)
	v93 = base.I64_extend_i32_u(v84)
	if base.Ui64(int64(4294967295)) < base.Ui64(v93*v83) {
		v105 = v92
		goto L16
	} else {
		goto L19
	}
L19:
	;
	if base.Ui64(v91) < base.Ui64(v93) {
		v105 = v92
		goto L16
	} else {
		goto L20
	}
L20:
	;
	goto L17
L21:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v105)+4)) = v70
	*(*int32)(unsafe.Add(mBase, uint32(v105))) = v70
	v112 = v105 + int32(12)
	*(*int32)(unsafe.Add(mBase, uint32(v105)+8)) = v112
	if v70 == int32(0) {
		v508 = v112
		goto L22
	} else {
		goto L23
	}
L22:
	;
	v541 = *(*int32)(unsafe.Add(mBase, uint32(l2)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v54)+4)) = v541
	if v541 != 0 {
		goto L45
	} else {
		goto L46
	}
L23:
	;
	v120 = v112 + v70<<(uint(int32(2))%32)
	v122 = v70 + int32(-1)
	if v122 != 0 {
		goto L25
	} else {
		goto L26
	}
L24:
	;
	if v70&int32(1) == int32(0) {
		goto L33
	} else {
		goto L34
	}
L25:
	;
	v124 = int32(_a_F_VP8LGetHistoImageSymbols_1)
	if int32(0) < l6 {
		goto L27
	} else {
		goto L28
	}
L26:
	;
	v232 = v120
	v233 = v112
	v236 = int32(0)
	goto L24
L27:
	;
	v127 = v73 + v124
	goto L29
L28:
	;
	v127 = v124
	goto L29
L29:
	;
	v133 = int32(4)
	v147 = v120
	v148 = v112
	v151 = int32(0)
	goto L30
L30:
	;
	v184 = int32(-4)
	v186 = int32(31)
	v188 = int32(-32)
	v189 = (v147 + v186) & v188
	*(*int32)(unsafe.Add(mBase, uint32(v148+v133+v184))) = v189
	v191 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v192 = v191 + v133
	v195 = *(*int32)(unsafe.Add(mBase, uint32(v192+v184)))
	v196 = int32(3312)
	*(*int32)(unsafe.Add(mBase, uint32(v195))) = v189 + v196
	v203 = (v189 + v127 + v186) & v188
	*(*int32)(unsafe.Add(mBase, uint32(v192))) = v203
	v205 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v207 = *(*int32)(unsafe.Add(mBase, uint32(v205+v133)))
	*(*int32)(unsafe.Add(mBase, uint32(v207))) = v203 + v196
	v213 = v203 + v127
	v215 = v151 + int32(2)
	if v70&int32(-2) != v215 {
		v133 = v133 + int32(8)
		v147 = v213
		v148 = v205
		v151 = v215
		goto L30
	} else {
		goto L32
	}
L31:
	;
	v232 = v213
	v233 = v205
	v236 = v215
	goto L24
L32:
	;
	goto L31
L33:
	;
	if v122 != 0 {
		goto L36
	} else {
		goto L37
	}
L34:
	;
	v271 = v236 << (uint(int32(2)) % 32)
	v276 = (v232 + int32(31)) & int32(-32)
	*(*int32)(unsafe.Add(mBase, uint32(v233+v271))) = v276
	v278 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v280 = *(*int32)(unsafe.Add(mBase, uint32(v278+v271)))
	*(*int32)(unsafe.Add(mBase, uint32(v280))) = v276 + int32(3312)
	goto L33
L35:
	;
	if v70&int32(1) == int32(0) {
		goto L41
	} else {
		goto L42
	}
L36:
	;
	v308 = int32(4)
	v312 = int32(0)
	goto L38
L37:
	;
	v423 = int32(0)
	goto L35
L38:
	;
	v344 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v348 = *(*int32)(unsafe.Add(mBase, uint32(v344+v308+int32(-4))))
	v349 = int64(-1)
	*(*int64)(unsafe.Add(mBase, uint32(v348)+3240)) = v349
	v351 = int32(16843009)
	*(*int32)(unsafe.Add(mBase, uint32(v348)+3304)) = v351
	v353 = base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k0)
	v354 = int32(3256)
	base.Simd_g_v128_store(m, v348, v354, v353)
	v356 = int32(3248)
	v358 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v348+v356))) = uint16(v358)
	v360 = int32(3308)
	v362 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v348+v360))) = uint8(v362)
	v364 = int32(3272)
	v366 = int32(0)
	base.Simd_g_v128_store(m, v348+v364, v366, v353)
	v368 = int32(3288)
	base.Simd_g_v128_store(m, v348+v368, v366, v353)
	*(*int32)(unsafe.Add(mBase, uint32(v348)+3236)) = l6
	v373 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v375 = *(*int32)(unsafe.Add(mBase, uint32(v373+v308)))
	*(*int64)(unsafe.Add(mBase, uint32(v375)+3240)) = v349
	*(*int32)(unsafe.Add(mBase, uint32(v375)+3304)) = v351
	base.Simd_g_v128_store(m, v375, v354, v353)
	*(*int32)(unsafe.Add(mBase, uint32(v375)+3236)) = l6
	*(*uint16)(unsafe.Add(mBase, uint32(v375+v356))) = uint16(v358)
	*(*uint8)(unsafe.Add(mBase, uint32(v375+v360))) = uint8(v362)
	base.Simd_g_v128_store(m, v375+v364, v366, v353)
	base.Simd_g_v128_store(m, v375+v368, v366, v353)
	v402 = v312 + int32(2)
	if v70&int32(-2) != v402 {
		v308 = v308 + int32(8)
		v312 = v402
		goto L38
	} else {
		goto L40
	}
L39:
	;
	v423 = v402
	goto L35
L40:
	;
	goto L39
L41:
	;
	v488 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v508 = v488
	goto L22
L42:
	;
	v457 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v461 = *(*int32)(unsafe.Add(mBase, uint32(v457+v423<<(uint(int32(2))%32))))
	*(*int64)(unsafe.Add(mBase, uint32(v461)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v461)+3236)) = l6
	*(*int32)(unsafe.Add(mBase, uint32(v461)+3304)) = int32(16843009)
	v467 = base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k0)
	base.Simd_g_v128_store(m, v461, int32(3256), v467)
	v472 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v461+int32(3248)))) = uint16(v472)
	v476 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v461+int32(3308)))) = uint8(v476)
	v480 = int32(0)
	base.Simd_g_v128_store(m, v461+int32(3272), v480, v467)
	base.Simd_g_v128_store(m, v461+int32(3288), v480, v467)
	goto L41
L43:
	;
	v554 = int32(0)
	v557 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v558 = *(*int32)(unsafe.Add(mBase, uint32(v557)))
	v559 = *(*int32)(unsafe.Add(mBase, uint32(v558)+3236))
	v560 = int32(4) << (uint(v559) % 32)
	v561 = int32(_a_F_VP8LGetHistoImageSymbols_0)
	v565 = base.B2i32(v554 < v559)
	if v554 < v559 {
		goto L47
	} else {
		goto L48
	}
L44:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+8)) = v550
	*(*int32)(unsafe.Add(mBase, uint32(v54))) = v551
	goto L43
L45:
	;
	v545 = *(*int32)(unsafe.Add(mBase, uint32(v541)+4))
	v546 = *(*int32)(unsafe.Add(mBase, uint32(v541)+8))
	v550 = v545 + v546<<(uint(int32(3))%32)
	v551 = v545
	goto L44
L46:
	;
	v543 = int32(0)
	v550 = v543
	v551 = v543
	goto L44
L47:
	;
	v566 = v560 + v561
	goto L49
L48:
	;
	v566 = v561
	goto L49
L49:
	;
	v567 = *(*int32)(unsafe.Add(mBase, uint32(v105)+4))
	v570 = v566*v567 + int32(12)
	if base.Ui32(v570) < base.Ui32(int32(33)) {
		goto L51
	} else {
		goto L52
	}
L50:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v105)+4)) = v567
	*(*int32)(unsafe.Add(mBase, uint32(v105)+8)) = v112
	*(*int32)(unsafe.Add(mBase, uint32(v105))) = v567
	if v567 < int32(1) {
		goto L64
	} else {
		goto L65
	}
L51:
	;
	if v570 == int32(0) {
		goto L53
	} else {
		goto L54
	}
L52:
	;
	base.MemoryFill(m, v105, v554, v570)
	goto L50
L53:
	;
	goto L50
L54:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v105))) = uint8(v554)
	v581 = v105 + v570
	*(*uint8)(unsafe.Add(mBase, uint32(v581+int32(-1)))) = uint8(v554)
	if base.Ui32(v570) < base.Ui32(int32(3)) {
		goto L53
	} else {
		goto L55
	}
L55:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v105)+2)) = uint8(v554)
	*(*uint8)(unsafe.Add(mBase, uint32(v105)+1)) = uint8(v554)
	*(*uint8)(unsafe.Add(mBase, uint32(v581+int32(-3)))) = uint8(v554)
	*(*uint8)(unsafe.Add(mBase, uint32(v581+int32(-2)))) = uint8(v554)
	if base.Ui32(v570) < base.Ui32(int32(7)) {
		goto L53
	} else {
		goto L56
	}
L56:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v105)+3)) = uint8(v554)
	*(*uint8)(unsafe.Add(mBase, uint32(v581+int32(-4)))) = uint8(v554)
	if base.Ui32(v570) < base.Ui32(int32(9)) {
		goto L53
	} else {
		goto L57
	}
L57:
	;
	v603 = int32(0)
	v606 = (v603 - v105) & int32(3)
	v607 = v105 + v606
	*(*int32)(unsafe.Add(mBase, uint32(v607))) = v603
	v615 = (v570 - v606) & int32(60)
	v616 = v607 + v615
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-4)))) = v603
	if base.Ui32(v615) < base.Ui32(int32(9)) {
		goto L53
	} else {
		goto L58
	}
L58:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v607)+8)) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v607)+4)) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-8)))) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-12)))) = v603
	if base.Ui32(v615) < base.Ui32(int32(25)) {
		goto L53
	} else {
		goto L59
	}
L59:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v607)+24)) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v607)+20)) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v607)+16)) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v607)+12)) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-16)))) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-20)))) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-24)))) = v603
	*(*int32)(unsafe.Add(mBase, uint32(v616+int32(-28)))) = v603
	v651 = v607&int32(4) | int32(24)
	v652 = v615 - v651
	if base.Ui32(v652) < base.Ui32(int32(32)) {
		goto L53
	} else {
		goto L60
	}
L60:
	;
	v657 = base.I64_extend_i32_u(v603) * int64(4294967297)
	v660 = v652
	v661 = v607 + v651
	goto L61
L61:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v661)+24)) = v657
	*(*int64)(unsafe.Add(mBase, uint32(v661)+16)) = v657
	*(*int64)(unsafe.Add(mBase, uint32(v661)+8)) = v657
	*(*int64)(unsafe.Add(mBase, uint32(v661))) = v657
	v673 = v660 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v673) {
		v660 = v673
		v661 = v661 + int32(32)
		goto L61
	} else {
		goto L63
	}
L62:
	;
	goto L53
L63:
	;
	goto L62
L64:
	;
	v1104 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	if v1104 == int32(0) {
		goto L85
	} else {
		goto L86
	}
L65:
	;
	v696 = int32(1)
	v700 = v112 + v567<<(uint(int32(2))%32)
	if v567 == v696 {
		v802 = v700
		v812 = v554
		v814 = v112
		goto L66
	} else {
		goto L67
	}
L66:
	;
	if v567&v696 == int32(0) {
		v865 = v814
		goto L74
	} else {
		goto L75
	}
L67:
	;
	v703 = int32(_a_F_VP8LGetHistoImageSymbols_1)
	if v554 < v559 {
		goto L68
	} else {
		goto L69
	}
L68:
	;
	v706 = v560 + v703
	goto L70
L69:
	;
	v706 = v703
	goto L70
L70:
	;
	v712 = int32(4)
	v717 = v700
	v727 = int32(0)
	v729 = v112
	goto L71
L71:
	;
	v763 = int32(-4)
	v765 = int32(31)
	v767 = int32(-32)
	v768 = (v717 + v765) & v767
	*(*int32)(unsafe.Add(mBase, uint32(v729+v712+v763))) = v768
	v770 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v771 = v770 + v712
	v774 = *(*int32)(unsafe.Add(mBase, uint32(v771+v763)))
	v775 = int32(3312)
	*(*int32)(unsafe.Add(mBase, uint32(v774))) = v768 + v775
	v782 = (v768 + v706 + v765) & v767
	*(*int32)(unsafe.Add(mBase, uint32(v771))) = v782
	v784 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v786 = *(*int32)(unsafe.Add(mBase, uint32(v784+v712)))
	*(*int32)(unsafe.Add(mBase, uint32(v786))) = v782 + v775
	v792 = v782 + v706
	v794 = v727 + int32(2)
	if v567&int32(2147483646) != v794 {
		v712 = v712 + int32(8)
		v717 = v792
		v727 = v794
		v729 = v784
		goto L71
	} else {
		goto L73
	}
L72:
	;
	v802 = v792
	v812 = v794
	v814 = v784
	goto L66
L73:
	;
	goto L72
L74:
	;
	if base.Ui32(v567) < base.Ui32(int32(4)) {
		v961 = int32(0)
		goto L76
	} else {
		goto L77
	}
L75:
	;
	v850 = v812 << (uint(int32(2)) % 32)
	v855 = (v802 + int32(31)) & int32(-32)
	*(*int32)(unsafe.Add(mBase, uint32(v814+v850))) = v855
	v857 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v859 = *(*int32)(unsafe.Add(mBase, uint32(v857+v850)))
	*(*int32)(unsafe.Add(mBase, uint32(v859))) = v855 + int32(3312)
	v865 = v857
	goto L74
L76:
	;
	v997 = v865 + v961<<(uint(int32(2))%32)
	v1002 = v567 - v961
	goto L82
L77:
	;
	v870 = v567 & int32(2147483644)
	v872 = v865
	v877 = v870
	goto L78
L78:
	;
	v924 = *(*int32)(unsafe.Add(mBase, uint32(v872+int32(12))))
	v927 = *(*int32)(unsafe.Add(mBase, uint32(v872+int32(8))))
	v930 = *(*int32)(unsafe.Add(mBase, uint32(v872+int32(4))))
	v931 = *(*int32)(unsafe.Add(mBase, uint32(v872)))
	*(*int32)(unsafe.Add(mBase, uint32(v931)+3236)) = v559
	*(*int32)(unsafe.Add(mBase, uint32(v930)+3236)) = v559
	*(*int32)(unsafe.Add(mBase, uint32(v927)+3236)) = v559
	*(*int32)(unsafe.Add(mBase, uint32(v924)+3236)) = v559
	v939 = v877 + int32(-4)
	if v939 != 0 {
		v872 = v872 + int32(16)
		v877 = v939
		goto L78
	} else {
		goto L80
	}
L79:
	;
	if v567 == v870 {
		goto L64
	} else {
		goto L81
	}
L80:
	;
	goto L79
L81:
	;
	v961 = v870
	goto L76
L82:
	;
	v1047 = *(*int32)(unsafe.Add(mBase, uint32(v997)))
	*(*int32)(unsafe.Add(mBase, uint32(v1047)+3236)) = v559
	v1052 = v1002 + int32(-1)
	if v1052 != 0 {
		v997 = v997 + int32(4)
		v1002 = v1052
		goto L82
	} else {
		goto L84
	}
L83:
	;
	goto L64
L84:
	;
	goto L83
L85:
	;
	v1366 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v1366
	v1369 = *(*int32)(unsafe.Add(mBase, uint32(v105)+4))
	if v1369 < int32(1) {
		v3148 = v1366
		goto L7
	} else {
		goto L101
	}
L86:
	;
	v1113 = int32(0)
	v1116 = v1113
	v1121 = v1104
	v1130 = v1113
	goto L87
L87:
	;
	v1168 = int32(2)
	v1175 = *(*int32)(unsafe.Add(mBase, uint32(v508+v1130>>(uint(l5)%32)*int32(base.Ui32(l0+int32(1)<<(uint(l5)%32)+int32(-1))>>(uint(l5)%32))<<(uint(v1168)%32)+v1116>>(uint(l5)%32)<<(uint(v1168)%32))))
	v1176 = int32(0)
	F_HistogramAddSinglePixOrCopy(m, v1175, v1121, v1176, v1176)
	mBase = m.M
	v1181 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1121+v1168))))
	v1182 = v1116 + v1181
	if v1182 < l0 {
		v1240 = v1182
		v1254 = v1130
		goto L89
	} else {
		goto L90
	}
L88:
	;
	goto L85
L89:
	;
	v1290 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	v1292 = v1290 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v54))) = v1292
	v1294 = *(*int32)(unsafe.Add(mBase, uint32(v54)+8))
	if v1292 != v1294 {
		v1314 = v1292
		goto L94
	} else {
		goto L95
	}
L90:
	;
	v1185 = v1182
	v1199 = v1130
	goto L91
L91:
	;
	v1236 = v1199 + int32(1)
	v1237 = v1185 - l0
	if l0 <= v1237 {
		v1185 = v1237
		v1199 = v1236
		goto L91
	} else {
		goto L93
	}
L92:
	;
	v1240 = v1237
	v1254 = v1236
	goto L89
L93:
	;
	goto L92
L94:
	;
	if v1314 != 0 {
		v1116 = v1240
		v1121 = v1314
		v1130 = v1254
		goto L87
	} else {
		goto L100
	}
L95:
	;
	v1299 = *(*int32)(unsafe.Add(mBase, uint32(v54)+4))
	v1300 = *(*int32)(unsafe.Add(mBase, uint32(v1299)))
	if v1300 != 0 {
		goto L98
	} else {
		goto L99
	}
L96:
	;
	v1313 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	v1314 = v1313
	goto L94
L97:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+8)) = v1309
	*(*int32)(unsafe.Add(mBase, uint32(v54))) = v1308
	*(*int32)(unsafe.Add(mBase, uint32(v54)+4)) = v1300
	goto L96
L98:
	;
	v1303 = *(*int32)(unsafe.Add(mBase, uint32(v1300)+4))
	v1304 = *(*int32)(unsafe.Add(mBase, uint32(v1300)+8))
	v1308 = v1303
	v1309 = v1303 + v1304<<(uint(int32(3))%32)
	goto L97
L99:
	;
	v1301 = int32(0)
	v1308 = v1301
	v1309 = v1301
	goto L97
L100:
	;
	goto L88
L101:
	;
	if l4 != 0 {
		goto L102
	} else {
		goto L103
	}
L102:
	;
	v1374 = int32(4)
	goto L104
L103:
	;
	v1374 = int32(64)
	goto L104
L104:
	;
	v1375 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v1376 = *(*int32)(unsafe.Add(mBase, uint32(v105)+8))
	v1378 = v1376
	v1383 = int32(0)
	goto L105
L105:
	;
	v1429 = *(*int32)(unsafe.Add(mBase, uint32(v1378)))
	v1430 = *(*int32)(unsafe.Add(mBase, uint32(v1429)))
	v1432 = *(*int32)(unsafe.Add(mBase, uint32(v1429)+3236))
	v1434 = int32(280)
	if int32(0) < v1432 {
		goto L107
	} else {
		goto L108
	}
L106:
	;
	v1535 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if int32(99) < l3 {
		v3148 = v1535
		goto L7
	} else {
		goto L124
	}
L107:
	;
	v1439 = int32(1)<<(uint(v1432)%32) + v1434
	goto L109
L108:
	;
	v1439 = v1434
	goto L109
L109:
	;
	v1444 = F_PopulationCost(m, v1430, v1439, v1429+int32(3240), v1429+int32(3304))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v1429)+3264)) = v1444
	v1448 = int32(256)
	v1452 = v1429 + int32(3305)
	v1453 = F_PopulationCost(m, v1429+int32(4), v1448, v1429+int32(3242), v1452)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v1429)+3272)) = v1453
	v1461 = v1429 + int32(3306)
	v1462 = F_PopulationCost(m, v1429+int32(1028), v1448, v1429+int32(3244), v1461)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v1429)+3280)) = v1462
	v1470 = v1429 + int32(3307)
	v1471 = F_PopulationCost(m, v1429+int32(2052), v1448, v1429+int32(3246), v1470)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v1429)+3288)) = v1471
	v1479 = v1429 + int32(3308)
	v1480 = F_PopulationCost(m, v1429+int32(3076), int32(40), v1429+int32(3248), v1479)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v1429)+3296)) = v1480
	v1482 = *(*int64)(unsafe.Add(mBase, uint32(v1429)+3264))
	v1484 = *(*int64)(unsafe.Add(mBase, uint32(v1429)+3272))
	v1486 = *(*int64)(unsafe.Add(mBase, uint32(v1429)+3280))
	v1488 = *(*int64)(unsafe.Add(mBase, uint32(v1429)+3288))
	*(*int64)(unsafe.Add(mBase, uint32(v1429)+3256)) = v1480 + v1482 + v1484 + v1486 + v1488
	v1491 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1429)+3304)))
	if v1491 != 0 {
		goto L111
	} else {
		goto L112
	}
L110:
	;
	v1532 = v1383 + int32(1)
	v1533 = *(*int32)(unsafe.Add(mBase, uint32(v105)+4))
	if v1532 < v1533 {
		v1378 = v1378 + int32(4)
		v1383 = v1532
		goto L105
	} else {
		goto L123
	}
L111:
	;
	v1498 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	v1502 = *(*int32)(unsafe.Add(mBase, uint32(v1375+v1498<<(uint(int32(2))%32))))
	v1503 = *(*int32)(unsafe.Add(mBase, uint32(v1502)))
	v1505 = *(*int32)(unsafe.Add(mBase, uint32(v1502)+3236))
	v1506 = int32(4) << (uint(v1505) % 32)
	v1507 = int32(_a_F_VP8LGetHistoImageSymbols_1)
	v1511 = base.B2i32(int32(0) < v1505)
	if int32(0) < v1505 {
		goto L117
	} else {
		goto L118
	}
L112:
	;
	v1492 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1452))))
	if v1492 != 0 {
		goto L111
	} else {
		goto L113
	}
L113:
	;
	v1493 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1461))))
	if v1493 != 0 {
		goto L111
	} else {
		goto L114
	}
L114:
	;
	v1494 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1470))))
	if v1494 != 0 {
		goto L111
	} else {
		goto L115
	}
L115:
	;
	v1495 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1479))))
	if v1495 != 0 {
		goto L111
	} else {
		goto L116
	}
L116:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1378))) = int32(0)
	goto L110
L117:
	;
	v1512 = v1506 + v1507
	goto L119
L118:
	;
	v1512 = v1507
	goto L119
L119:
	;
	v1513 = F_memcpy(m, v1502, v1429, v1512)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v1513))) = v1503
	v1515 = *(*int32)(unsafe.Add(mBase, uint32(v1429)))
	v1516 = int32(1120)
	if int32(0) < v1505 {
		goto L120
	} else {
		goto L121
	}
L120:
	;
	v1519 = v1506 + v1516
	goto L122
L121:
	;
	v1519 = v1516
	goto L122
L122:
	;
	v1520 = F_memcpy(m, v1503, v1515, v1519)
	mBase = m.M
	v1521 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v1521 + int32(1)
	goto L110
L123:
	;
	goto L106
L124:
	;
	if v1535 <= v1374<<(uint(int32(1))%32) {
		v3148 = v1535
		goto L7
	} else {
		goto L125
	}
L125:
	;
	if int32(89) < l3 {
		v1558 = int32(16)
		goto L126
	} else {
		goto L127
	}
L126:
	;
	v1559 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	if v1535 < int32(1) {
		goto L8
	} else {
		goto L131
	}
L127:
	;
	if int32(256) < v70 {
		goto L128
	} else {
		goto L129
	}
L128:
	;
	v1548 = int32(8)
	goto L130
L129:
	;
	v1548 = int32(16)
	goto L130
L130:
	;
	v1558 = int32(base.Ui32(int32(base.Ui32(int32(base.Ui32(v1548)>>(uint(base.B2i32(int32(512) < v70))%32)))>>(uint(base.B2i32(int32(1024) < v70))%32))) >> (uint(base.B2i32(l3 < int32(51))) % 32))
	goto L126
L131:
	;
	if v1535 != int32(1) {
		goto L132
	} else {
		goto L133
	}
L132:
	;
	v1572 = v1535 & int32(2147483646)
	v1573 = base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k0)
	v1574 = base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k1)
	v1576 = v1559
	v1580 = v1572
	v1599 = v1573
	v1608 = v1573
	v1609 = v1574
	v1610 = v1574
	v1611 = v1573
	v1612 = v1574
	goto L134
L133:
	;
	v1564 = int64(-1)
	v1565 = int64(0)
	v1859 = int32(0)
	v1868 = v1564
	v1869 = v1565
	v1870 = v1565
	v1871 = v1564
	v1872 = v1565
	v1873 = v1564
	goto L10
L134:
	;
	v1626 = *(*int32)(unsafe.Add(mBase, uint32(v1576)))
	v1627 = *(*int64)(unsafe.Add(mBase, uint32(v1626)+3280))
	v1631 = *(*int32)(unsafe.Add(mBase, uint32(v1576+int32(4))))
	v1632 = *(*int64)(unsafe.Add(mBase, uint32(v1631)+3280))
	v1634 = base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1627), v1632)
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1612)) < base.Ui64(v1627) {
		goto L136
	} else {
		goto L137
	}
L135:
	;
	v1746 = base.Simd_g_i8x16_swizzle_c(v1650, base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k2))
	v1747 = base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k0)
	v1750 = int32(0)
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1650)) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1746)) {
		goto L173
	} else {
		goto L174
	}
L136:
	;
	v1640 = int64(-1)
	goto L138
L137:
	;
	v1640 = int64(0)
	goto L138
L138:
	;
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l1(v1612)) < base.Ui64(v1632) {
		goto L139
	} else {
		goto L140
	}
L139:
	;
	v1647 = int64(-1)
	goto L141
L140:
	;
	v1647 = int64(0)
	goto L141
L141:
	;
	v1650 = base.Simd_g_v128_bitselect(v1612, v1634, base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1640), v1647))
	if base.Ui64(v1627) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1611)) {
		goto L142
	} else {
		goto L143
	}
L142:
	;
	v1656 = int64(-1)
	goto L144
L143:
	;
	v1656 = int64(0)
	goto L144
L144:
	;
	if base.Ui64(v1632) < base.Ui64(base.Simd_g_i64x2_extract_lane_l1(v1611)) {
		goto L145
	} else {
		goto L146
	}
L145:
	;
	v1663 = int64(-1)
	goto L147
L146:
	;
	v1663 = int64(0)
	goto L147
L147:
	;
	v1664 = int32(1)
	v1666 = base.Simd_g_v128_bitselect(v1611, v1634, base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1656), v1663))
	v1667 = *(*int64)(unsafe.Add(mBase, uint32(v1626)+3272))
	v1669 = *(*int64)(unsafe.Add(mBase, uint32(v1631)+3272))
	v1671 = base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1667), v1669)
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1610)) < base.Ui64(v1667) {
		goto L148
	} else {
		goto L149
	}
L148:
	;
	v1677 = int64(-1)
	goto L150
L149:
	;
	v1677 = int64(0)
	goto L150
L150:
	;
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l1(v1610)) < base.Ui64(v1669) {
		goto L151
	} else {
		goto L152
	}
L151:
	;
	v1684 = int64(-1)
	goto L153
L152:
	;
	v1684 = int64(0)
	goto L153
L153:
	;
	v1687 = base.Simd_g_v128_bitselect(v1610, v1671, base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1677), v1684))
	if base.Ui64(v1667) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1608)) {
		goto L154
	} else {
		goto L155
	}
L154:
	;
	v1693 = int64(-1)
	goto L156
L155:
	;
	v1693 = int64(0)
	goto L156
L156:
	;
	if base.Ui64(v1669) < base.Ui64(base.Simd_g_i64x2_extract_lane_l1(v1608)) {
		goto L157
	} else {
		goto L158
	}
L157:
	;
	v1700 = int64(-1)
	goto L159
L158:
	;
	v1700 = int64(0)
	goto L159
L159:
	;
	v1701 = int32(1)
	v1703 = base.Simd_g_v128_bitselect(v1608, v1671, base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1693), v1700))
	v1704 = *(*int64)(unsafe.Add(mBase, uint32(v1626)+3264))
	v1706 = *(*int64)(unsafe.Add(mBase, uint32(v1631)+3264))
	v1708 = base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1704), v1706)
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1609)) < base.Ui64(v1704) {
		goto L160
	} else {
		goto L161
	}
L160:
	;
	v1714 = int64(-1)
	goto L162
L161:
	;
	v1714 = int64(0)
	goto L162
L162:
	;
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l1(v1609)) < base.Ui64(v1706) {
		goto L163
	} else {
		goto L164
	}
L163:
	;
	v1721 = int64(-1)
	goto L165
L164:
	;
	v1721 = int64(0)
	goto L165
L165:
	;
	v1724 = base.Simd_g_v128_bitselect(v1609, v1708, base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1714), v1721))
	if base.Ui64(v1704) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1599)) {
		goto L166
	} else {
		goto L167
	}
L166:
	;
	v1730 = int64(-1)
	goto L168
L167:
	;
	v1730 = int64(0)
	goto L168
L168:
	;
	if base.Ui64(v1706) < base.Ui64(base.Simd_g_i64x2_extract_lane_l1(v1599)) {
		goto L169
	} else {
		goto L170
	}
L169:
	;
	v1737 = int64(-1)
	goto L171
L170:
	;
	v1737 = int64(0)
	goto L171
L171:
	;
	v1740 = base.Simd_g_v128_bitselect(v1599, v1708, base.Simd_g_i64x2_replace_lane_l1(base.Simd_g_i64x2_splat(v1730), v1737))
	v1744 = v1580 + int32(-2)
	if v1744 != 0 {
		v1576 = v1576 + int32(8)
		v1580 = v1744
		v1599 = v1740
		v1608 = v1703
		v1609 = v1724
		v1610 = v1687
		v1611 = v1666
		v1612 = v1650
		goto L134
	} else {
		goto L172
	}
L172:
	;
	goto L135
L173:
	;
	v1755 = int64(-1)
	goto L175
L174:
	;
	v1755 = int64(0)
	goto L175
L175:
	;
	v1756 = int32(0)
	v1760 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_v128_bitselect(v1650, v1746, base.Simd_g_i64x2_replace_lane_l0(v1747, v1755)))
	v1762 = base.Simd_g_i8x16_swizzle_c(v1666, base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k2))
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1762)) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1666)) {
		goto L176
	} else {
		goto L177
	}
L176:
	;
	v1770 = int64(-1)
	goto L178
L177:
	;
	v1770 = int64(0)
	goto L178
L178:
	;
	v1771 = int32(0)
	v1775 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_v128_bitselect(v1666, v1762, base.Simd_g_i64x2_replace_lane_l0(v1747, v1770)))
	v1777 = base.Simd_g_i8x16_swizzle_c(v1687, base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k2))
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1687)) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1777)) {
		goto L179
	} else {
		goto L180
	}
L179:
	;
	v1785 = int64(-1)
	goto L181
L180:
	;
	v1785 = int64(0)
	goto L181
L181:
	;
	v1786 = int32(0)
	v1790 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_v128_bitselect(v1687, v1777, base.Simd_g_i64x2_replace_lane_l0(v1747, v1785)))
	v1792 = base.Simd_g_i8x16_swizzle_c(v1703, base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k2))
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1792)) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1703)) {
		goto L182
	} else {
		goto L183
	}
L182:
	;
	v1800 = int64(-1)
	goto L184
L183:
	;
	v1800 = int64(0)
	goto L184
L184:
	;
	v1801 = int32(0)
	v1805 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_v128_bitselect(v1703, v1792, base.Simd_g_i64x2_replace_lane_l0(v1747, v1800)))
	v1807 = base.Simd_g_i8x16_swizzle_c(v1724, base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k2))
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1724)) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1807)) {
		goto L185
	} else {
		goto L186
	}
L185:
	;
	v1815 = int64(-1)
	goto L187
L186:
	;
	v1815 = int64(0)
	goto L187
L187:
	;
	v1816 = int32(0)
	v1820 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_v128_bitselect(v1724, v1807, base.Simd_g_i64x2_replace_lane_l0(v1747, v1815)))
	v1822 = base.Simd_g_i8x16_swizzle_c(v1740, base.Simd_g_const(&F_VP8LGetHistoImageSymbols__k2))
	if base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1822)) < base.Ui64(base.Simd_g_i64x2_extract_lane_l0(v1740)) {
		goto L188
	} else {
		goto L189
	}
L188:
	;
	v1830 = int64(-1)
	goto L190
L189:
	;
	v1830 = int64(0)
	goto L190
L190:
	;
	v1831 = int32(0)
	v1835 = base.Simd_g_i64x2_extract_lane_l0(base.Simd_g_v128_bitselect(v1740, v1822, base.Simd_g_i64x2_replace_lane_l0(v1747, v1830)))
	if v1535 != v1572 {
		v1859 = v1572
		v1868 = v1820
		v1869 = v1835
		v1870 = v1805
		v1871 = v1790
		v1872 = v1775
		v1873 = v1760
		goto L10
	} else {
		goto L191
	}
L191:
	;
	v1994 = v1820
	v1995 = v1835
	v1996 = v1805
	v1997 = v1790
	v1998 = v1775
	v1999 = v1760
	goto L9
L192:
	;
	goto L3
L193:
	;
	goto L192
L194:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l10)+92)) = int32(1)
	goto L193
L195:
	;
	v1947 = *(*int32)(unsafe.Add(mBase, uint32(v1897)))
	v1948 = *(*int64)(unsafe.Add(mBase, uint32(v1947)+3280))
	if base.Ui64(v1928) < base.Ui64(v1948) {
		goto L197
	} else {
		goto L198
	}
L196:
	;
	v1994 = v1960
	v1995 = v1962
	v1996 = v1957
	v1997 = v1955
	v1998 = v1952
	v1999 = v1950
	goto L9
L197:
	;
	v1950 = v1928
	goto L199
L198:
	;
	v1950 = v1948
	goto L199
L199:
	;
	if base.Ui64(v1948) < base.Ui64(v1927) {
		goto L200
	} else {
		goto L201
	}
L200:
	;
	v1952 = v1927
	goto L202
L201:
	;
	v1952 = v1948
	goto L202
L202:
	;
	v1953 = *(*int64)(unsafe.Add(mBase, uint32(v1947)+3272))
	if base.Ui64(v1926) < base.Ui64(v1953) {
		goto L203
	} else {
		goto L204
	}
L203:
	;
	v1955 = v1926
	goto L205
L204:
	;
	v1955 = v1953
	goto L205
L205:
	;
	if base.Ui64(v1953) < base.Ui64(v1925) {
		goto L206
	} else {
		goto L207
	}
L206:
	;
	v1957 = v1925
	goto L208
L207:
	;
	v1957 = v1953
	goto L208
L208:
	;
	v1958 = *(*int64)(unsafe.Add(mBase, uint32(v1947)+3264))
	if base.Ui64(v1923) < base.Ui64(v1958) {
		goto L209
	} else {
		goto L210
	}
L209:
	;
	v1960 = v1923
	goto L211
L210:
	;
	v1960 = v1958
	goto L211
L211:
	;
	if base.Ui64(v1958) < base.Ui64(v1924) {
		goto L212
	} else {
		goto L213
	}
L212:
	;
	v1962 = v1924
	goto L214
L213:
	;
	v1962 = v1958
	goto L214
L214:
	;
	v1966 = v1902 + int32(-1)
	if v1966 != 0 {
		v1897 = v1897 + int32(4)
		v1902 = v1966
		v1923 = v1960
		v1924 = v1962
		v1925 = v1957
		v1926 = v1955
		v1927 = v1952
		v1928 = v1950
		goto L195
	} else {
		goto L215
	}
L215:
	;
	goto L196
L216:
	;
	if v1995 == v1994 {
		goto L234
	} else {
		goto L235
	}
L217:
	;
	v2040 = v1559
	v2043 = v1535
	goto L218
L218:
	;
	v2076 = *(*int32)(unsafe.Add(mBase, uint32(v2040)))
	v2077 = int32(0)
	if v1995 == v1994 {
		v2092 = v2077
		goto L220
	} else {
		goto L221
	}
L220:
	;
	if v1996 == v1997 {
		v2111 = v2077
		goto L224
	} else {
		goto L225
	}
L221:
	;
	v2079 = *(*int64)(unsafe.Add(mBase, uint32(v2076)+3264))
	v2084 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v2079-v1994), float64(3.999999)), v2019)
	if base.F64_lt(base.F64_abs(v2084), float64(2.147483648e+09)) == int32(0) {
		goto L222
	} else {
		goto L223
	}
L222:
	;
	v2092 = int32(-2147483648)
	goto L220
L223:
	;
	v2090 = base.I32_trunc_f64_s(v2084)
	v2092 = v2090
	goto L220
L224:
	;
	if v1998 == v1999 {
		v2131 = int32(0)
		goto L229
	} else {
		goto L230
	}
L225:
	;
	v2095 = *(*int64)(unsafe.Add(mBase, uint32(v2076)+3272))
	v2100 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v2095-v1997), float64(3.999999)), base.F64_convert_i64_u(v1996-v1997))
	if base.F64_lt(base.F64_abs(v2100), float64(2.147483648e+09)) == int32(0) {
		goto L227
	} else {
		goto L228
	}
L226:
	;
	v2111 = v2108 << (uint(int32(2)) % 32)
	goto L224
L227:
	;
	v2108 = int32(-2147483648)
	goto L226
L228:
	;
	v2106 = base.I32_trunc_f64_s(v2100)
	v2108 = v2106
	goto L226
L229:
	;
	v2133 = v2111 + v2092<<(uint(int32(4))%32) + v2131
	*(*uint16)(unsafe.Add(mBase, uint32(v2076)+3310)) = uint16(v2133)
	v2138 = v2043 + int32(-1)
	if v2138 != 0 {
		v2040 = v2040 + int32(4)
		v2043 = v2138
		goto L218
	} else {
		goto L233
	}
L230:
	;
	v2118 = *(*int64)(unsafe.Add(mBase, uint32(v2076)+3280))
	v2123 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v2118-v1999), float64(3.999999)), base.F64_convert_i64_u(v1998-v1999))
	if base.F64_lt(base.F64_abs(v2123), float64(2.147483648e+09)) == int32(0) {
		goto L231
	} else {
		goto L232
	}
L231:
	;
	v2131 = int32(-2147483648)
	goto L229
L232:
	;
	v2129 = base.I32_trunc_f64_s(v2123)
	v2131 = v2129
	goto L229
L233:
	;
	goto L8
L234:
	;
	if base.Ui32(v1535) < base.Ui32(int32(4)) {
		v2408 = int32(0)
		goto L251
	} else {
		goto L252
	}
L235:
	;
	v2140 = int32(1)
	if v1535 == v2140 {
		v2253 = int32(0)
		goto L236
	} else {
		goto L237
	}
L236:
	;
	if v1535&v2140 == int32(0) {
		goto L8
	} else {
		goto L247
	}
L237:
	;
	v2149 = v1559
	v2163 = int32(0)
	goto L238
L238:
	;
	v2199 = *(*int32)(unsafe.Add(mBase, uint32(v2149)))
	v2200 = *(*int64)(unsafe.Add(mBase, uint32(v2199)+3264))
	v2205 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v2200-v1994), float64(3.999999)), v2019)
	if base.F64_lt(base.F64_abs(v2205), float64(2.147483648e+09)) == int32(0) {
		goto L241
	} else {
		goto L242
	}
L239:
	;
	v2253 = v2236
	goto L236
L240:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2199)+3310)) = uint16(v2213)
	v2217 = *(*int32)(unsafe.Add(mBase, uint32(v2149+int32(4))))
	v2218 = *(*int64)(unsafe.Add(mBase, uint32(v2217)+3264))
	v2223 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v2218-v1994), float64(3.999999)), v2019)
	if base.F64_lt(base.F64_abs(v2223), float64(2.147483648e+09)) == int32(0) {
		goto L244
	} else {
		goto L245
	}
L241:
	;
	v2213 = int32(-2147483648)
	goto L240
L242:
	;
	v2211 = base.I32_trunc_f64_s(v2205)
	v2213 = v2211
	goto L240
L243:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2217)+3310)) = uint16(v2231)
	v2236 = v2163 + int32(2)
	if v1535&int32(2147483646) != v2236 {
		v2149 = v2149 + int32(8)
		v2163 = v2236
		goto L238
	} else {
		goto L246
	}
L244:
	;
	v2231 = int32(-2147483648)
	goto L243
L245:
	;
	v2229 = base.I32_trunc_f64_s(v2223)
	v2231 = v2229
	goto L243
L246:
	;
	goto L239
L247:
	;
	v2294 = *(*int32)(unsafe.Add(mBase, uint32(v1559+v2253<<(uint(int32(2))%32))))
	v2295 = *(*int64)(unsafe.Add(mBase, uint32(v2294)+3264))
	v2300 = base.F64_div(base.F64_mul(base.F64_convert_i64_u(v2295-v1994), float64(3.999999)), v2019)
	if base.F64_lt(base.F64_abs(v2300), float64(2.147483648e+09)) == int32(0) {
		goto L249
	} else {
		goto L250
	}
L248:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2294)+3310)) = uint16(v2308)
	goto L8
L249:
	;
	v2308 = int32(-2147483648)
	goto L248
L250:
	;
	v2306 = base.I32_trunc_f64_s(v2300)
	v2308 = v2306
	goto L248
L251:
	;
	v2445 = v1559 + v2408<<(uint(int32(2))%32)
	v2459 = v1535 - v2408
	goto L257
L252:
	;
	v2314 = v1535 & int32(2147483644)
	v2316 = v1559
	v2330 = v2314
	goto L253
L253:
	;
	v2368 = *(*int32)(unsafe.Add(mBase, uint32(v2316+int32(12))))
	v2371 = *(*int32)(unsafe.Add(mBase, uint32(v2316+int32(8))))
	v2374 = *(*int32)(unsafe.Add(mBase, uint32(v2316+int32(4))))
	v2375 = *(*int32)(unsafe.Add(mBase, uint32(v2316)))
	v2376 = int32(0)
	*(*uint16)(unsafe.Add(mBase, uint32(v2375)+3310)) = uint16(v2376)
	*(*uint16)(unsafe.Add(mBase, uint32(v2374)+3310)) = uint16(v2376)
	*(*uint16)(unsafe.Add(mBase, uint32(v2371)+3310)) = uint16(v2376)
	*(*uint16)(unsafe.Add(mBase, uint32(v2368)+3310)) = uint16(v2376)
	v2387 = v2330 + int32(-4)
	if v2387 != 0 {
		v2316 = v2316 + int32(16)
		v2330 = v2387
		goto L253
	} else {
		goto L255
	}
L254:
	;
	if v1535 == v2314 {
		goto L8
	} else {
		goto L256
	}
L255:
	;
	goto L254
L256:
	;
	v2408 = v2314
	goto L251
L257:
	;
	v2495 = *(*int32)(unsafe.Add(mBase, uint32(v2445)))
	v2496 = int32(0)
	*(*uint16)(unsafe.Add(mBase, uint32(v2495)+3310)) = uint16(v2496)
	v2501 = v2459 + int32(-1)
	if v2501 != 0 {
		v2445 = v2445 + int32(4)
		v2459 = v2501
		goto L257
	} else {
		goto L259
	}
L258:
	;
	goto L8
L259:
	;
	goto L258
L260:
	;
	v2607 = v54 + v2557
	v2608 = int64(281470681808895)
	*(*int64)(unsafe.Add(mBase, uint32(v2607))) = v2608
	*(*int64)(unsafe.Add(mBase, uint32(v2607+int32(8)))) = v2608
	v2615 = v2557 + int32(16)
	if v1374<<(uint(int32(2))%32) != v2615 {
		v2557 = v2615
		goto L260
	} else {
		goto L262
	}
L261:
	;
	if v1535 < int32(1) {
		v3097 = v1535
		goto L263
	} else {
		goto L264
	}
L262:
	;
	goto L261
L263:
	;
	if l4 != 0 {
		v5848 = v3097
		goto L6
	} else {
		goto L312
	}
L264:
	;
	if l4 == int32(0) {
		goto L266
	} else {
		goto L267
	}
L265:
	;
	if l4 == int32(0) {
		v3097 = v2922
		goto L263
	} else {
		goto L304
	}
L266:
	;
	v2713 = int32(0)
	v2720 = l8
	goto L274
L267:
	;
	v2623 = int32(0)
	v2642 = v1535
	goto L268
L268:
	;
	v2673 = int32(2)
	v2674 = v2623 << (uint(v2673) % 32)
	v2676 = *(*int32)(unsafe.Add(mBase, uint32(v1559+v2674)))
	v2677 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2676)+3310)))
	v2680 = v54 + v2677<<(uint(v2673)%32)
	v2681 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2680))))
	if v2681 != int32(-1) {
		goto L271
	} else {
		goto L272
	}
L270:
	;
	if v2705 < v2708 {
		v2623 = v2705
		v2642 = v2708
		goto L268
	} else {
		goto L273
	}
L271:
	;
	v2687 = int32(2)
	v2690 = *(*int32)(unsafe.Add(mBase, uint32(v1559+v2681<<(uint(v2687)%32))))
	F_HistogramAdd(m, v2676, v2690, v2690)
	mBase = m.M
	v2692 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	v2694 = v2692 + int32(-1)
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v2694
	v2696 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v2703 = *(*int32)(unsafe.Add(mBase, uint32(v2696+v2692<<(uint(v2687)%32)+int32(-4))))
	*(*int32)(unsafe.Add(mBase, uint32(v2696+v2674))) = v2703
	v2705 = v2623
	v2708 = v2694
	goto L270
L272:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2680))) = uint16(v2623)
	v2705 = v2623 + int32(1)
	v2708 = v2642
	goto L270
L273:
	;
	v2922 = v2708
	goto L265
L274:
	;
	v2763 = int32(2)
	v2764 = v2713 << (uint(v2763) % 32)
	v2765 = v1559 + v2764
	v2766 = *(*int32)(unsafe.Add(mBase, uint32(v2765)))
	v2767 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2766)+3310)))
	v2770 = v54 + v2767<<(uint(v2763)%32)
	v2771 = int32(*(*int16)(unsafe.Add(mBase, uint32(v2770))))
	if v2771 == int32(-1) {
		goto L279
	} else {
		goto L280
	}
L275:
	;
	v2922 = v2900
	goto L265
L276:
	;
	v2900 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v2889 < v2900 {
		v2713 = v2889
		v2720 = v2892
		goto L274
	} else {
		goto L303
	}
L277:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v2776))) = v2720
	v2876 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v2876 + int32(-1)
	v2880 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v2887 = *(*int32)(unsafe.Add(mBase, uint32(v2880+v2876<<(uint(int32(2))%32)+int32(-4))))
	*(*int32)(unsafe.Add(mBase, uint32(v2880+v2764))) = v2887
	v2889 = v2713
	v2892 = v2873
	goto L276
L278:
	;
	v2889 = v2713 + int32(1)
	v2892 = v2720
	goto L276
L279:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2770))) = uint16(v2713)
	goto L278
L280:
	;
	v2776 = v1559 + v2771<<(uint(int32(2))%32)
	v2777 = *(*int32)(unsafe.Add(mBase, uint32(v2776)))
	v2781 = *(*int64)(unsafe.Add(mBase, uint32(v2766)+3256))
	v2782 = v2781 * base.I64_extend_i32_u(v1558)
	if v2782 < int64(0) {
		goto L281
	} else {
		goto L282
	}
L281:
	;
	v2785 = int64(-50)
	goto L283
L282:
	;
	v2785 = int64(50)
	goto L283
L283:
	;
	v2786 = v2785 + v2782
	v2788 = base.I64_div_s(v2786, int64(-100))
	v2789 = *(*int64)(unsafe.Add(mBase, uint32(v2777)+3256))
	v2790 = v2789 + v2781
	v2791 = v2788 + v2790
	if v2788^int64(9223372036854775807) < v2790 {
		goto L284
	} else {
		goto L285
	}
L284:
	;
	v2795 = int64(9223372036854775807)
	goto L286
L285:
	;
	v2795 = v2791
	goto L286
L286:
	;
	if v2786 < int64(100) {
		goto L287
	} else {
		goto L288
	}
L287:
	;
	v2798 = v2795
	goto L289
L288:
	;
	v2798 = v2791
	goto L289
L289:
	;
	v2803 = F_GetCombinedHistogramEntropy(m, v2777, v2766, v2798, v54+int32(328), v54+int32(256))
	mBase = m.M
	if v2803 == int32(0) {
		goto L278
	} else {
		goto L290
	}
L290:
	;
	F_HistogramAdd(m, v2777, v2766, v2720)
	mBase = m.M
	v2807 = *(*int64)(unsafe.Add(mBase, uint32(v54)+328))
	*(*int64)(unsafe.Add(mBase, uint32(v2720)+3256)) = v2807
	v2810 = base.Simd_g_v128_load(m, v54, int32(256))
	base.Simd_g_v128_store(m, v2720, int32(3264), v2810)
	v2814 = base.Simd_g_v128_load(m, v54, int32(272))
	base.Simd_g_v128_store(m, v2720, int32(3280), v2814)
	v2817 = *(*int64)(unsafe.Add(mBase, uint32(v54)+288))
	*(*int64)(unsafe.Add(mBase, uint32(v2720)+3296)) = v2817
	v2819 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2720)+3242)))
	if v2819 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		goto L291
	} else {
		goto L292
	}
L291:
	;
	v2829 = int32(1)
	v2831 = *(*int32)(unsafe.Add(mBase, uint32(v2765)))
	v2832 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2831)+3242)))
	if v2832 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		v2842 = v2829
		goto L295
	} else {
		goto L296
	}
L292:
	;
	v2822 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2720)+3244)))
	if v2822 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		goto L291
	} else {
		goto L293
	}
L293:
	;
	v2825 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2720)+3246)))
	if v2825 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		goto L291
	} else {
		goto L294
	}
L294:
	;
	v2828 = *(*int32)(unsafe.Add(mBase, uint32(v2776)))
	v2873 = v2828
	goto L277
L295:
	;
	v2843 = *(*int32)(unsafe.Add(mBase, uint32(v2776)))
	v2844 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2843)+3242)))
	if v2844 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		v2853 = v2829
		goto L298
	} else {
		goto L299
	}
L296:
	;
	v2836 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2831)+3244)))
	if v2836 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		v2842 = int32(1)
		goto L295
	} else {
		goto L297
	}
L297:
	;
	v2839 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2831)+3246)))
	v2842 = base.B2i32(v2839 == int32(_a_F_VP8LGetHistoImageSymbols_2))
	goto L295
L298:
	;
	if v2842&v2853 != 0 {
		v2873 = v2843
		goto L277
	} else {
		goto L301
	}
L299:
	;
	v2847 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2843)+3244)))
	if v2847 == int32(_a_F_VP8LGetHistoImageSymbols_2) {
		v2853 = v2829
		goto L298
	} else {
		goto L300
	}
L300:
	;
	v2850 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2843)+3246)))
	v2853 = base.B2i32(v2850 == int32(_a_F_VP8LGetHistoImageSymbols_2))
	goto L298
L301:
	;
	v2855 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2770)+2)))
	if base.Ui32(int32(31)) < base.Ui32(v2855) {
		v2873 = v2843
		goto L277
	} else {
		goto L302
	}
L302:
	;
	v2859 = v2855 + int32(1)
	*(*uint16)(unsafe.Add(mBase, uint32(v2770)+2)) = uint16(v2859)
	goto L278
L303:
	;
	goto L275
L304:
	;
	if v2922 < int32(1) {
		v3097 = v2922
		goto L263
	} else {
		goto L305
	}
L305:
	;
	v2958 = v1559
	v2973 = int32(0)
	goto L306
L306:
	;
	v3009 = *(*int32)(unsafe.Add(mBase, uint32(v2958)))
	v3010 = *(*int32)(unsafe.Add(mBase, uint32(v3009)))
	v3012 = *(*int32)(unsafe.Add(mBase, uint32(v3009)+3236))
	v3014 = int32(280)
	if int32(0) < v3012 {
		goto L308
	} else {
		goto L309
	}
L307:
	;
	v3097 = v3075
	goto L263
L308:
	;
	v3019 = int32(1)<<(uint(v3012)%32) + v3014
	goto L310
L309:
	;
	v3019 = v3014
	goto L310
L310:
	;
	v3024 = F_PopulationCost(m, v3010, v3019, v3009+int32(3240), v3009+int32(3304))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3009)+3264)) = v3024
	v3026 = int32(4)
	v3028 = int32(256)
	v3033 = F_PopulationCost(m, v3009+v3026, v3028, v3009+int32(3242), v3009+int32(3305))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3009)+3272)) = v3033
	v3042 = F_PopulationCost(m, v3009+int32(1028), v3028, v3009+int32(3244), v3009+int32(3306))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3009)+3280)) = v3042
	v3051 = F_PopulationCost(m, v3009+int32(2052), v3028, v3009+int32(3246), v3009+int32(3307))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3009)+3288)) = v3051
	v3060 = F_PopulationCost(m, v3009+int32(3076), int32(40), v3009+int32(3248), v3009+int32(3308))
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v3009)+3296)) = v3060
	v3062 = *(*int64)(unsafe.Add(mBase, uint32(v3009)+3264))
	v3064 = *(*int64)(unsafe.Add(mBase, uint32(v3009)+3272))
	v3066 = *(*int64)(unsafe.Add(mBase, uint32(v3009)+3280))
	v3068 = *(*int64)(unsafe.Add(mBase, uint32(v3009)+3288))
	*(*int64)(unsafe.Add(mBase, uint32(v3009)+3256)) = v3060 + v3062 + v3064 + v3066 + v3068
	v3074 = v2973 + int32(1)
	v3075 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v3074 < v3075 {
		v2958 = v2958 + v3026
		v2973 = v3074
		goto L306
	} else {
		goto L311
	}
L311:
	;
	goto L307
L312:
	;
	v3148 = v3097
	goto L7
L313:
	;
	v4228 = v4197 * v4197
	v4231 = base.I64_extend_i32_u(v4228 + int32(1))
	v4232 = int32(64)
	if v4231 == int64(0) {
		goto L398
	} else {
		goto L399
	}
L314:
	;
	v3188 = int64(-500000)
	goto L316
L315:
	;
	v3188 = int64(500000)
	goto L316
L316:
	;
	v3194 = base.I64_div_s(v3188+base.I64_extend_i32_s(v3185*int32(99)), int64(1000000))
	v3195 = base.I32_wrap_i64(v3194)
	if v3148 <= v3195 {
		v4197 = v3148
		v4198 = v3181
		goto L313
	} else {
		goto L317
	}
L317:
	;
	v3197 = int64(10)
	v3198 = int32(64)
	goto L321
L318:
	;
	if v3219 == int32(0) {
		goto L5
	} else {
		goto L324
	}
L319:
	;
	goto L318
L320:
	;
	v3217 = F_malloc(m, base.I32_wrap_i64(v3197)*v3198)
	mBase = m.M
	v3219 = v3217
	goto L319
L321:
	;
	v3205 = base.I64_div_u_s(int64(2147418112), v3197)
	v3206 = int32(0)
	v3207 = base.I64_extend_i32_u(v3198)
	if base.Ui64(int64(4294967295)) < base.Ui64(v3207*v3197) {
		v3219 = v3206
		goto L319
	} else {
		goto L322
	}
L322:
	;
	if base.Ui64(v3205) < base.Ui64(v3207) {
		v3219 = v3206
		goto L319
	} else {
		goto L323
	}
L323:
	;
	goto L320
L324:
	;
	v3223 = int32(1)
	if v3148 < v3223 {
		goto L325
	} else {
		goto L326
	}
L325:
	;
	v4172 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	F_free(m, v3219)
	mBase = m.M
	goto L394
L326:
	;
	v3228 = v3219 + int32(-64)
	v3232 = v54 + int32(280)
	v3236 = v54 + int32(272)
	v3237 = int32(0)
	v3241 = v3237
	v3263 = int32(1)
	v3288 = v3237
	v3289 = v3237
	goto L327
L327:
	;
	v3292 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v3292 <= v3195 {
		goto L325
	} else {
		goto L329
	}
L328:
	;
	goto L325
L329:
	;
	v3295 = v3289 + int32(1)
	if v3180 <= v3295 {
		goto L325
	} else {
		goto L330
	}
L330:
	;
	if v3241 != 0 {
		goto L332
	} else {
		goto L333
	}
L331:
	;
	v3300 = int32(2)
	v3301 = base.I32_div_s(v3292, v3300)
	if v3292 < v3300 {
		v3580 = v3241
		v3602 = v3263
		goto L336
	} else {
		goto L337
	}
L332:
	;
	v3298 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+8))
	v3299 = v3298
	goto L331
L333:
	;
	v3299 = int64(0)
	goto L331
L334:
	;
	v4119 = v3288 + int32(1)
	if v4119 != v3148 {
		v3241 = v4067
		v3263 = v4089
		v3288 = v4119
		v3289 = v4115
		goto L327
	} else {
		goto L393
	}
L335:
	;
	v3683 = *(*int32)(unsafe.Add(mBase, uint32(v3219)+4))
	v3684 = int32(2)
	v3685 = v3683 << (uint(v3684) % 32)
	v3687 = *(*int32)(unsafe.Add(mBase, uint32(v3181+v3685)))
	v3688 = *(*int32)(unsafe.Add(mBase, uint32(v3219)))
	v3691 = v3181 + v3688<<(uint(v3684)%32)
	v3692 = *(*int32)(unsafe.Add(mBase, uint32(v3691)))
	F_HistogramAdd(m, v3687, v3692, v3692)
	mBase = m.M
	v3694 = *(*int32)(unsafe.Add(mBase, uint32(v3691)))
	v3695 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+16))
	*(*int64)(unsafe.Add(mBase, uint32(v3694)+3256)) = v3695
	v3697 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+24))
	*(*int64)(unsafe.Add(mBase, uint32(v3694)+3264)) = v3697
	v3699 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+32))
	*(*int64)(unsafe.Add(mBase, uint32(v3694)+3272)) = v3699
	v3701 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+40))
	*(*int64)(unsafe.Add(mBase, uint32(v3694)+3280)) = v3701
	v3703 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+48))
	*(*int64)(unsafe.Add(mBase, uint32(v3694)+3288)) = v3703
	v3705 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v3705 + int32(-1)
	v3709 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v3716 = *(*int32)(unsafe.Add(mBase, uint32(v3709+v3705<<(uint(v3684)%32)+int32(-4))))
	*(*int32)(unsafe.Add(mBase, uint32(v3709+v3685))) = v3716
	v3718 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+56))
	*(*int64)(unsafe.Add(mBase, uint32(v3694)+3296)) = v3718
	if v3632 < int32(1) {
		v4015 = v3632
		goto L365
	} else {
		goto L366
	}
L336:
	;
	if v3580 != 0 {
		v3632 = v3580
		v3654 = v3602
		goto L335
	} else {
		goto L364
	}
L337:
	;
	v3308 = v3241
	v3313 = v3292
	v3314 = int32(1)
	v3330 = v3263
	v3334 = v3299
	goto L338
L338:
	;
	v3363 = base.I64_rem_u_s(base.I64_extend_i32_u(v3330)*int64(48271), int64(2147483647))
	v3364 = base.I32_wrap_i64(v3363)
	if v3308 == int32(9) {
		v3559 = v3308
		v3567 = v3334
		goto L340
	} else {
		goto L341
	}
L339:
	;
	v3580 = v3559
	v3602 = v3364
	goto L336
L340:
	;
	v3574 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v3574 < int32(2) {
		v3580 = v3559
		v3602 = v3364
		goto L336
	} else {
		goto L362
	}
L341:
	;
	v3367 = base.I32_rem_u_s(v3364, (v3292+int32(-1))*v3292)
	v3369 = v3313 + int32(-1)
	v3370 = base.I32_div_u_s(v3367, v3369)
	v3372 = v3367 - v3370*v3369
	v3374 = v3372 + base.B2i32(base.Ui32(v3370) <= base.Ui32(v3372))
	if base.Ui32(v3370) < base.Ui32(v3374) {
		goto L342
	} else {
		goto L343
	}
L342:
	;
	v3376 = v3370
	goto L344
L343:
	;
	v3376 = v3374
	goto L344
L344:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+256)) = v3376
	v3381 = *(*int32)(unsafe.Add(mBase, uint32(v3181+v3376<<(uint(int32(2))%32))))
	v3382 = *(*int64)(unsafe.Add(mBase, uint32(v3381)+3256))
	if base.Ui32(v3374) < base.Ui32(v3370) {
		goto L345
	} else {
		goto L346
	}
L345:
	;
	v3384 = v3370
	goto L347
L346:
	;
	v3384 = v3374
	goto L347
L347:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+260)) = v3384
	v3389 = *(*int32)(unsafe.Add(mBase, uint32(v3181+v3384<<(uint(int32(2))%32))))
	v3390 = int64(9223372036854775807)
	v3391 = *(*int64)(unsafe.Add(mBase, uint32(v3389)+3256))
	v3392 = v3382 + v3391
	v3393 = v3392 + v3334
	if v3334^v3390 < v3392 {
		goto L348
	} else {
		goto L349
	}
L348:
	;
	v3397 = v3390
	goto L350
L349:
	;
	v3397 = v3393
	goto L350
L350:
	;
	if int64(-1) < v3334 {
		goto L351
	} else {
		goto L352
	}
L351:
	;
	v3400 = v3397
	goto L353
L352:
	;
	v3400 = v3393
	goto L353
L353:
	;
	v3401 = F_GetCombinedHistogramEntropy(m, v3381, v3389, v3400, v3236, v3232)
	mBase = m.M
	if v3401 == int32(0) {
		v3559 = v3308
		v3567 = v3334
		goto L340
	} else {
		goto L354
	}
L354:
	;
	v3404 = int32(6)
	v3406 = v3219 + v3308<<(uint(v3404)%32)
	v3413 = int32(0)
	v3414 = base.Simd_g_v128_load(m, v54+int32(304), v3413)
	base.Simd_g_v128_store(m, v3406+int32(48), v3413, v3414)
	v3424 = base.Simd_g_v128_load(m, v54+int32(288), v3413)
	base.Simd_g_v128_store(m, v3406+int32(32), v3413, v3424)
	v3429 = *(*int64)(unsafe.Add(mBase, uint32(v3232)))
	*(*int64)(unsafe.Add(mBase, uint32(v3406+int32(24)))) = v3429
	v3433 = *(*int64)(unsafe.Add(mBase, uint32(v3236)))
	*(*int64)(unsafe.Add(mBase, uint32(v3406+int32(16)))) = v3433
	v3437 = v3433 - v3392
	*(*int64)(unsafe.Add(mBase, uint32(v3406+int32(8)))) = v3437
	v3439 = *(*int64)(unsafe.Add(mBase, uint32(v54)+256))
	*(*int64)(unsafe.Add(mBase, uint32(v3406))) = v3439
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(264)))) = v3437
	v3447 = v3308 + int32(1)
	v3450 = v3219 + v3447<<(uint(v3404)%32)
	v3453 = *(*int64)(unsafe.Add(mBase, uint32(v3450+int32(-56))))
	v3454 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+8))
	if v3454 <= v3453 {
		goto L355
	} else {
		goto L356
	}
L355:
	;
	if int64(-1) < v3437 {
		goto L357
	} else {
		goto L358
	}
L356:
	;
	v3456 = int32(0)
	v3457 = base.Simd_g_v128_load(m, v3219, v3456)
	v3459 = v3450 + int32(-64)
	v3461 = base.Simd_g_v128_load(m, v3459, v3456)
	base.Simd_g_v128_store(m, v3219, v3456, v3461)
	v3464 = int32(48)
	v3465 = v3219 + v3464
	v3467 = base.Simd_g_v128_load(m, v3465, v3456)
	v3469 = v3450 + int32(-16)
	v3471 = base.Simd_g_v128_load(m, v3469, v3456)
	base.Simd_g_v128_store(m, v3465, v3456, v3471)
	v3474 = int32(32)
	v3475 = v3219 + v3474
	v3477 = base.Simd_g_v128_load(m, v3475, v3456)
	v3479 = v3450 + int32(-32)
	v3481 = base.Simd_g_v128_load(m, v3479, v3456)
	base.Simd_g_v128_store(m, v3475, v3456, v3481)
	v3484 = int32(16)
	v3485 = v3219 + v3484
	v3487 = base.Simd_g_v128_load(m, v3485, v3456)
	v3489 = v3450 + int32(-48)
	v3491 = base.Simd_g_v128_load(m, v3489, v3456)
	base.Simd_g_v128_store(m, v3485, v3456, v3491)
	v3495 = v54 + v3484
	base.Simd_g_v128_store(m, v3495, v3456, v3487)
	v3499 = v54 + v3474
	base.Simd_g_v128_store(m, v3499, v3456, v3477)
	v3503 = v54 + v3464
	base.Simd_g_v128_store(m, v3503, v3456, v3467)
	base.Simd_g_v128_store(m, v54, v3456, v3457)
	v3508 = *(*int64)(unsafe.Add(mBase, uint32(v54)))
	*(*int64)(unsafe.Add(mBase, uint32(v3459))) = v3508
	v3510 = *(*int64)(unsafe.Add(mBase, uint32(v3495)))
	*(*int64)(unsafe.Add(mBase, uint32(v3489))) = v3510
	v3516 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(24))))
	*(*int64)(unsafe.Add(mBase, uint32(v3450+int32(-40)))) = v3516
	v3518 = *(*int64)(unsafe.Add(mBase, uint32(v3499)))
	*(*int64)(unsafe.Add(mBase, uint32(v3479))) = v3518
	v3520 = *(*int64)(unsafe.Add(mBase, uint32(v3503)))
	*(*int64)(unsafe.Add(mBase, uint32(v3469))) = v3520
	v3526 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(56))))
	*(*int64)(unsafe.Add(mBase, uint32(v3450+int32(-8)))) = v3526
	v3532 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(40))))
	*(*int64)(unsafe.Add(mBase, uint32(v3450+int32(-24)))) = v3532
	v3538 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(8))))
	*(*int64)(unsafe.Add(mBase, uint32(v3450+int32(-56)))) = v3538
	goto L355
L357:
	;
	v3553 = v3334
	goto L359
L358:
	;
	v3553 = v3437
	goto L359
L359:
	;
	if v3447 != int32(9) {
		v3559 = v3447
		v3567 = v3553
		goto L340
	} else {
		goto L360
	}
L360:
	;
	if int64(0) <= v3437 {
		v3559 = v3447
		v3567 = v3553
		goto L340
	} else {
		goto L361
	}
L361:
	;
	v3632 = int32(9)
	v3654 = v3364
	goto L335
L362:
	;
	if v3314 < v3301 {
		v3308 = v3559
		v3313 = v3574
		v3314 = v3314 + int32(1)
		v3330 = v3364
		v3334 = v3567
		goto L338
	} else {
		goto L363
	}
L363:
	;
	goto L339
L364:
	;
	v4067 = int32(0)
	v4089 = v3602
	v4115 = v3295
	goto L334
L365:
	;
	v4067 = v4015
	v4089 = v3654
	v4115 = int32(0)
	goto L334
L366:
	;
	v3723 = v3632
	v3742 = int32(0)
	goto L367
L367:
	;
	v3776 = v3219 + v3742<<(uint(int32(6))%32)
	v3777 = *(*int32)(unsafe.Add(mBase, uint32(v3776)+4))
	v3780 = base.B2i32(v3777 == v3688) | base.B2i32(v3777 == v3683)
	v3781 = *(*int32)(unsafe.Add(mBase, uint32(v3776)))
	v3784 = base.B2i32(v3781 == v3688) | base.B2i32(v3781 == v3683)
	if v3784 != int32(1) {
		goto L370
	} else {
		goto L371
	}
L368:
	;
	v4015 = v4001
	goto L365
L369:
	;
	if v4006 < v4001 {
		v3723 = v4001
		v3742 = v4006
		goto L367
	} else {
		goto L392
	}
L370:
	;
	if v3784|v3780 != int32(1) {
		v3858 = v3781
		goto L374
	} else {
		goto L375
	}
L371:
	;
	if v3780 == int32(0) {
		goto L370
	} else {
		goto L372
	}
L372:
	;
	v3791 = v3228 + v3723<<(uint(int32(6))%32)
	v3792 = int32(0)
	v3793 = base.Simd_g_v128_load(m, v3791, v3792)
	base.Simd_g_v128_store(m, v3776, v3792, v3793)
	v3796 = int32(48)
	v3801 = base.Simd_g_v128_load(m, v3791+v3796, v3792)
	base.Simd_g_v128_store(m, v3776+v3796, v3792, v3801)
	v3804 = int32(32)
	v3809 = base.Simd_g_v128_load(m, v3791+v3804, v3792)
	base.Simd_g_v128_store(m, v3776+v3804, v3792, v3809)
	v3812 = int32(16)
	v3817 = base.Simd_g_v128_load(m, v3791+v3812, v3792)
	base.Simd_g_v128_store(m, v3776+v3812, v3792, v3817)
	v4001 = v3723 + int32(-1)
	v4006 = v3742
	goto L369
L373:
	;
	v3972 = v3228 + v3723<<(uint(int32(6))%32)
	v3973 = int32(0)
	v3974 = base.Simd_g_v128_load(m, v3972, v3973)
	base.Simd_g_v128_store(m, v3776, v3973, v3974)
	v3977 = int32(48)
	v3982 = base.Simd_g_v128_load(m, v3972+v3977, v3973)
	base.Simd_g_v128_store(m, v3776+v3977, v3973, v3982)
	v3985 = int32(32)
	v3990 = base.Simd_g_v128_load(m, v3972+v3985, v3973)
	base.Simd_g_v128_store(m, v3776+v3985, v3973, v3990)
	v3996 = base.Simd_g_v128_load(m, v3972+int32(16), v3973)
	base.Simd_g_v128_store(m, v3848, v3973, v3996)
	v4001 = v3723 + int32(-1)
	v4006 = v3742
	goto L369
L374:
	;
	v3862 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v3858 != v3862 {
		v3865 = v3858
		goto L384
	} else {
		goto L385
	}
L375:
	;
	if v3781 != v3683 {
		v3827 = v3781
		goto L376
	} else {
		goto L377
	}
L376:
	;
	if v3777 != v3683 {
		v3830 = v3777
		goto L378
	} else {
		goto L379
	}
L377:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3776))) = v3688
	v3827 = v3688
	goto L376
L378:
	;
	if v3830 < v3827 {
		goto L381
	} else {
		goto L382
	}
L379:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3776)+4)) = v3688
	v3830 = v3688
	goto L378
L380:
	;
	v3836 = int32(2)
	v3839 = *(*int32)(unsafe.Add(mBase, uint32(v3181+v3834<<(uint(v3836)%32))))
	v3843 = *(*int32)(unsafe.Add(mBase, uint32(v3181+v3835<<(uint(v3836)%32))))
	v3844 = *(*int64)(unsafe.Add(mBase, uint32(v3843)+3256))
	v3845 = *(*int64)(unsafe.Add(mBase, uint32(v3839)+3256))
	v3846 = v3844 + v3845
	v3848 = v3776 + int32(16)
	v3851 = F_GetCombinedHistogramEntropy(m, v3839, v3843, v3846, v3848, v3776+int32(24))
	mBase = m.M
	if v3851 == int32(0) {
		goto L373
	} else {
		goto L383
	}
L381:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3776)+4)) = v3827
	*(*int32)(unsafe.Add(mBase, uint32(v3776))) = v3830
	v3834 = v3830
	v3835 = v3827
	goto L380
L382:
	;
	v3834 = v3827
	v3835 = v3830
	goto L380
L383:
	;
	v3854 = *(*int64)(unsafe.Add(mBase, uint32(v3776)+16))
	*(*int64)(unsafe.Add(mBase, uint32(v3776)+8)) = v3854 - v3846
	v3857 = *(*int32)(unsafe.Add(mBase, uint32(v3776)))
	v3858 = v3857
	goto L374
L384:
	;
	v3866 = *(*int32)(unsafe.Add(mBase, uint32(v3776)+4))
	if v3866 != v3862 {
		v3869 = v3866
		goto L386
	} else {
		goto L387
	}
L385:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3776))) = v3683
	v3865 = v3683
	goto L384
L386:
	;
	if v3865 <= v3869 {
		goto L388
	} else {
		goto L389
	}
L387:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3776)+4)) = v3683
	v3869 = v3683
	goto L386
L388:
	;
	v3873 = *(*int64)(unsafe.Add(mBase, uint32(v3776)+8))
	v3874 = *(*int64)(unsafe.Add(mBase, uint32(v3219)+8))
	if v3874 <= v3873 {
		goto L390
	} else {
		goto L391
	}
L389:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v3776)+4)) = v3865
	*(*int32)(unsafe.Add(mBase, uint32(v3776))) = v3869
	goto L388
L390:
	;
	v4001 = v3723
	v4006 = v3742 + int32(1)
	goto L369
L391:
	;
	v3876 = int32(0)
	v3877 = base.Simd_g_v128_load(m, v3219, v3876)
	v3879 = base.Simd_g_v128_load(m, v3776, v3876)
	base.Simd_g_v128_store(m, v3219, v3876, v3879)
	v3882 = int32(48)
	v3883 = v3219 + v3882
	v3885 = base.Simd_g_v128_load(m, v3883, v3876)
	v3887 = v3776 + v3882
	v3889 = base.Simd_g_v128_load(m, v3887, v3876)
	base.Simd_g_v128_store(m, v3883, v3876, v3889)
	v3892 = int32(32)
	v3893 = v3219 + v3892
	v3895 = base.Simd_g_v128_load(m, v3893, v3876)
	v3897 = v3776 + v3892
	v3899 = base.Simd_g_v128_load(m, v3897, v3876)
	base.Simd_g_v128_store(m, v3893, v3876, v3899)
	v3902 = int32(16)
	v3903 = v3219 + v3902
	v3905 = base.Simd_g_v128_load(m, v3903, v3876)
	v3907 = v3776 + v3902
	v3909 = base.Simd_g_v128_load(m, v3907, v3876)
	base.Simd_g_v128_store(m, v3903, v3876, v3909)
	v3913 = v54 + v3902
	base.Simd_g_v128_store(m, v3913, v3876, v3905)
	v3917 = v54 + v3892
	base.Simd_g_v128_store(m, v3917, v3876, v3895)
	v3921 = v54 + v3882
	base.Simd_g_v128_store(m, v3921, v3876, v3885)
	base.Simd_g_v128_store(m, v54, v3876, v3877)
	v3926 = *(*int64)(unsafe.Add(mBase, uint32(v54)))
	*(*int64)(unsafe.Add(mBase, uint32(v3776))) = v3926
	v3928 = int32(8)
	v3932 = *(*int64)(unsafe.Add(mBase, uint32(v54+v3928)))
	*(*int64)(unsafe.Add(mBase, uint32(v3776+v3928))) = v3932
	v3934 = *(*int64)(unsafe.Add(mBase, uint32(v3913)))
	*(*int64)(unsafe.Add(mBase, uint32(v3907))) = v3934
	v3936 = int32(24)
	v3940 = *(*int64)(unsafe.Add(mBase, uint32(v54+v3936)))
	*(*int64)(unsafe.Add(mBase, uint32(v3776+v3936))) = v3940
	v3942 = *(*int64)(unsafe.Add(mBase, uint32(v3917)))
	*(*int64)(unsafe.Add(mBase, uint32(v3897))) = v3942
	v3944 = *(*int64)(unsafe.Add(mBase, uint32(v3921)))
	*(*int64)(unsafe.Add(mBase, uint32(v3887))) = v3944
	v3946 = int32(56)
	v3950 = *(*int64)(unsafe.Add(mBase, uint32(v54+v3946)))
	*(*int64)(unsafe.Add(mBase, uint32(v3776+v3946))) = v3950
	v3952 = int32(40)
	v3956 = *(*int64)(unsafe.Add(mBase, uint32(v54+v3952)))
	*(*int64)(unsafe.Add(mBase, uint32(v3776+v3952))) = v3956
	goto L390
L392:
	;
	goto L368
L393:
	;
	goto L328
L394:
	;
	v4174 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v3195+v3223 < v4172 {
		v5848 = v4174
		goto L6
	} else {
		goto L395
	}
L395:
	;
	v4176 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v4197 = v4174
	v4198 = v4176
	goto L313
L396:
	;
	if v4253 == int32(0) {
		goto L4
	} else {
		goto L402
	}
L397:
	;
	goto L396
L398:
	;
	v4251 = F_malloc(m, base.I32_wrap_i64(v4231)*v4232)
	mBase = m.M
	v4253 = v4251
	goto L397
L399:
	;
	v4239 = base.I64_div_u_s(int64(2147418112), v4231)
	v4240 = int32(0)
	v4241 = base.I64_extend_i32_u(v4232)
	if base.Ui64(int64(4294967295)) < base.Ui64(v4241*v4231) {
		v4253 = v4240
		goto L397
	} else {
		goto L400
	}
L400:
	;
	if base.Ui64(v4239) < base.Ui64(v4241) {
		v4253 = v4240
		goto L397
	} else {
		goto L401
	}
L401:
	;
	goto L398
L402:
	;
	if v4197 < int32(1) {
		goto L403
	} else {
		goto L404
	}
L403:
	;
	F_free(m, v4253)
	mBase = m.M
	goto L547
L404:
	;
	v4259 = int32(0)
	v4262 = v4259
	v4311 = v4259
	goto L405
L405:
	;
	v4313 = v4311 + int32(1)
	if v4197 <= v4313 {
		v4811 = v4262
		goto L407
	} else {
		goto L408
	}
L406:
	;
	if v4811 < int32(1) {
		goto L403
	} else {
		goto L465
	}
L407:
	;
	if v4313 != v4197 {
		v4262 = v4811
		v4311 = v4313
		goto L405
	} else {
		goto L464
	}
L408:
	;
	v4316 = v4262
	v4337 = v4313
	goto L409
L409:
	;
	if v4316 == v4228 {
		v4757 = v4228
		goto L411
	} else {
		goto L412
	}
L410:
	;
	v4811 = v4757
	goto L407
L411:
	;
	v4808 = v4337 + int32(1)
	if v4808 != v4197 {
		v4316 = v4757
		v4337 = v4808
		goto L409
	} else {
		goto L463
	}
L412:
	;
	if v4311 < v4337 {
		goto L413
	} else {
		goto L414
	}
L413:
	;
	v4368 = v4311
	goto L415
L414:
	;
	v4368 = v4337
	goto L415
L415:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+256)) = v4368
	v4373 = *(*int32)(unsafe.Add(mBase, uint32(v4198+v4368<<(uint(int32(2))%32))))
	v4374 = *(*int64)(unsafe.Add(mBase, uint32(v4373)+3256))
	if v4337 < v4311 {
		goto L416
	} else {
		goto L417
	}
L416:
	;
	v4376 = v4311
	goto L418
L417:
	;
	v4376 = v4337
	goto L418
L418:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+260)) = v4376
	v4381 = *(*int32)(unsafe.Add(mBase, uint32(v4198+v4376<<(uint(int32(2))%32))))
	v4382 = *(*int64)(unsafe.Add(mBase, uint32(v4381)+3256))
	v4383 = v4374 + v4382
	if int64(1) <= v4383 {
		goto L419
	} else {
		goto L420
	}
L419:
	;
	v4386 = int32(4)
	v4390 = int32(1028)
	v4394 = int32(2052)
	v4398 = int32(3076)
	v4406 = int32(3240)
	v4407 = int32(0)
	v4421 = int32(3264)
	v4436 = int64(0)
	goto L422
L420:
	;
	v4757 = v4316
	goto L411
L421:
	;
	v4612 = int32(6)
	v4614 = v4253 + v4316<<(uint(v4612)%32)
	v4621 = int32(0)
	v4622 = base.Simd_g_v128_load(m, v54+int32(304), v4621)
	base.Simd_g_v128_store(m, v4614+int32(48), v4621, v4622)
	v4632 = base.Simd_g_v128_load(m, v54+int32(288), v4621)
	base.Simd_g_v128_store(m, v4614+int32(32), v4621, v4632)
	v4641 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(280))))
	*(*int64)(unsafe.Add(mBase, uint32(v4614+int32(24)))) = v4641
	*(*int64)(unsafe.Add(mBase, uint32(v4614+int32(16)))) = v4602
	v4648 = v4602 - v4383
	*(*int64)(unsafe.Add(mBase, uint32(v4614+int32(8)))) = v4648
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(264)))) = v4648
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(272)))) = v4602
	v4660 = *(*int64)(unsafe.Add(mBase, uint32(v54)+256))
	*(*int64)(unsafe.Add(mBase, uint32(v4614))) = v4660
	v4663 = v4316 + int32(1)
	v4666 = v4253 + v4663<<(uint(v4612)%32)
	v4669 = *(*int64)(unsafe.Add(mBase, uint32(v4666+int32(-56))))
	v4670 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+8))
	if v4670 <= v4669 {
		v4757 = v4663
		goto L411
	} else {
		goto L462
	}
L422:
	;
	v4461 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4373+v4406))))
	if v4461 != int32(_a_F_VP8LGetHistoImageSymbols_2) {
		goto L425
	} else {
		goto L426
	}
L423:
	;
	v4757 = v4316
	goto L411
L424:
	;
	v4469 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4373+v4407+int32(3304)))))
	if v4468 == int32(0) {
		goto L429
	} else {
		goto L430
	}
L425:
	;
	v4466 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4381+v4406))))
	v4468 = base.B2i32(v4461 != v4466)
	goto L424
L426:
	;
	v4468 = int32(1)
	goto L424
L427:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(256)+v4421+int32(-3240)))) = v4594
	v4602 = v4594 + v4436
	if base.Ui64(v4383) <= base.Ui64(v4602) {
		goto L459
	} else {
		goto L460
	}
L428:
	;
	v4490 = int32(256)
	switch v4407 {
	default:
		goto L439
	case 1:
		v4503 = v4373 + v4386
		v4504 = v4490
		v4505 = v4381 + v4386
		goto L435
	case 2:
		goto L438
	case 3:
		goto L437
	case 4:
		goto L436
	}
L429:
	;
	if v4469&int32(255) == int32(0) {
		goto L433
	} else {
		goto L434
	}
L430:
	;
	if v4469&int32(255) == int32(0) {
		goto L429
	} else {
		goto L431
	}
L431:
	;
	v4479 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v4381+v4407+int32(3304)))))
	if v4479&int32(255) != 0 {
		goto L428
	} else {
		goto L432
	}
L432:
	;
	goto L429
L433:
	;
	v4489 = *(*int64)(unsafe.Add(mBase, uint32(v4381+v4421)))
	v4594 = v4489
	goto L427
L434:
	;
	v4487 = *(*int64)(unsafe.Add(mBase, uint32(v4373+v4421)))
	v4594 = v4487
	goto L427
L435:
	;
	v4508 = m.G115
	v4509 = *(*int32)(unsafe.Add(mBase, uint32(v4508)))
	m.T0[v4509].(func(*base.Module, int32, int32, int32, int32, int32))(m, v4503, v4505, v4504, v54+int32(328), v54)
	mBase = m.M
	v4511 = *(*int32)(unsafe.Add(mBase, uint32(v54)+340))
	if v4511 <= int32(4) {
		goto L445
	} else {
		goto L446
	}
L436:
	;
	v4503 = v4373 + v4398
	v4504 = int32(40)
	v4505 = v4381 + v4398
	goto L435
L437:
	;
	v4503 = v4373 + v4394
	v4504 = v4490
	v4505 = v4381 + v4394
	goto L435
L438:
	;
	v4503 = v4373 + v4390
	v4504 = v4490
	v4505 = v4381 + v4390
	goto L435
L439:
	;
	v4492 = *(*int32)(unsafe.Add(mBase, uint32(v4381)+3236))
	v4494 = int32(280)
	if int32(0) < v4492 {
		goto L440
	} else {
		goto L441
	}
L440:
	;
	v4499 = int32(1)<<(uint(v4492)%32) + v4494
	goto L442
L441:
	;
	v4499 = v4494
	goto L442
L442:
	;
	v4500 = *(*int32)(unsafe.Add(mBase, uint32(v4381)))
	v4501 = *(*int32)(unsafe.Add(mBase, uint32(v4373)))
	v4503 = v4501
	v4504 = v4499
	v4505 = v4500
	goto L435
L443:
	;
	v4562 = *(*int32)(unsafe.Add(mBase, uint32(v54)+12))
	v4565 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	v4569 = *(*int32)(unsafe.Add(mBase, uint32(v54)+4))
	v4573 = *(*int32)(unsafe.Add(mBase, uint32(v54)+20))
	v4577 = *(*int32)(unsafe.Add(mBase, uint32(v54)+8))
	v4581 = *(*int32)(unsafe.Add(mBase, uint32(v54)+16))
	v4594 = v4560 + base.I64_extend_i32_u(v4562*int32(240)+v4565*int32(1600)+v4569*int32(2640)+v4573*int32(720)+v4577*int32(1840)+v4581*int32(3360))<<(uint(int64(13))%64) + int64(401814323)
	goto L427
L444:
	;
	v4536 = *(*int64)(unsafe.Add(mBase, uint32(v54)+328))
	v4539 = *(*int32)(unsafe.Add(mBase, uint32(v54)+336))
	v4542 = *(*int32)(unsafe.Add(mBase, uint32(v54)+344))
	v4551 = v4535*base.I64_extend_i32_u(v4539<<(uint(int32(1))%32)-v4542)<<(uint(int64(23))%64) + v4536*(int64(1000)-v4535)
	if v4551 < int64(0) {
		goto L453
	} else {
		goto L454
	}
L445:
	;
	if v4511 < int32(2) {
		v4560 = int64(0)
		goto L443
	} else {
		goto L447
	}
L446:
	;
	v4535 = int64(627)
	goto L444
L447:
	;
	switch v4511 + int32(-2) {
	case 0:
		goto L449
	case 1:
		v4535 = int64(950)
		goto L444
	default:
		goto L448
	}
L448:
	;
	v4535 = int64(700)
	goto L444
L449:
	;
	v4523 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v54)+336)))
	v4526 = *(*int64)(unsafe.Add(mBase, uint32(v54)+328))
	v4527 = v4523*int64(830472192) + v4526
	if v4527 < int64(0) {
		goto L450
	} else {
		goto L451
	}
L450:
	;
	v4530 = int64(-50)
	goto L452
L451:
	;
	v4530 = int64(50)
	goto L452
L452:
	;
	v4533 = base.I64_div_s(v4530+v4527, int64(100))
	v4560 = v4533
	goto L443
L453:
	;
	v4554 = int64(-500)
	goto L455
L454:
	;
	v4554 = int64(500)
	goto L455
L455:
	;
	v4557 = base.I64_div_s(v4554+v4551, int64(1000))
	if base.Ui64(v4557) < base.Ui64(v4536) {
		goto L456
	} else {
		goto L457
	}
L456:
	;
	v4559 = v4536
	goto L458
L457:
	;
	v4559 = v4557
	goto L458
L458:
	;
	v4560 = v4559
	goto L443
L459:
	;
	goto L423
L460:
	;
	v4609 = v4407 + int32(1)
	if v4609 == int32(5) {
		goto L421
	} else {
		goto L461
	}
L461:
	;
	v4406 = v4406 + int32(2)
	v4407 = v4609
	v4421 = v4421 + int32(8)
	v4436 = v4602
	goto L422
L462:
	;
	v4672 = int32(0)
	v4673 = base.Simd_g_v128_load(m, v4253, v4672)
	v4675 = v4666 + int32(-64)
	v4677 = base.Simd_g_v128_load(m, v4675, v4672)
	base.Simd_g_v128_store(m, v4253, v4672, v4677)
	v4680 = int32(48)
	v4681 = v4253 + v4680
	v4683 = base.Simd_g_v128_load(m, v4681, v4672)
	v4685 = v4666 + int32(-16)
	v4687 = base.Simd_g_v128_load(m, v4685, v4672)
	base.Simd_g_v128_store(m, v4681, v4672, v4687)
	v4690 = int32(32)
	v4691 = v4253 + v4690
	v4693 = base.Simd_g_v128_load(m, v4691, v4672)
	v4695 = v4666 + int32(-32)
	v4697 = base.Simd_g_v128_load(m, v4695, v4672)
	base.Simd_g_v128_store(m, v4691, v4672, v4697)
	v4700 = int32(16)
	v4701 = v4253 + v4700
	v4703 = base.Simd_g_v128_load(m, v4701, v4672)
	v4705 = v4666 + int32(-48)
	v4707 = base.Simd_g_v128_load(m, v4705, v4672)
	base.Simd_g_v128_store(m, v4701, v4672, v4707)
	v4711 = v54 + v4700
	base.Simd_g_v128_store(m, v4711, v4672, v4703)
	v4715 = v54 + v4690
	base.Simd_g_v128_store(m, v4715, v4672, v4693)
	v4719 = v54 + v4680
	base.Simd_g_v128_store(m, v4719, v4672, v4683)
	base.Simd_g_v128_store(m, v54, v4672, v4673)
	v4724 = *(*int64)(unsafe.Add(mBase, uint32(v54)))
	*(*int64)(unsafe.Add(mBase, uint32(v4675))) = v4724
	v4726 = *(*int64)(unsafe.Add(mBase, uint32(v4711)))
	*(*int64)(unsafe.Add(mBase, uint32(v4705))) = v4726
	v4732 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(24))))
	*(*int64)(unsafe.Add(mBase, uint32(v4666+int32(-40)))) = v4732
	v4734 = *(*int64)(unsafe.Add(mBase, uint32(v4715)))
	*(*int64)(unsafe.Add(mBase, uint32(v4695))) = v4734
	v4736 = *(*int64)(unsafe.Add(mBase, uint32(v4719)))
	*(*int64)(unsafe.Add(mBase, uint32(v4685))) = v4736
	v4742 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(56))))
	*(*int64)(unsafe.Add(mBase, uint32(v4666+int32(-8)))) = v4742
	v4748 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(40))))
	*(*int64)(unsafe.Add(mBase, uint32(v4666+int32(-24)))) = v4748
	v4754 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(8))))
	*(*int64)(unsafe.Add(mBase, uint32(v4666+int32(-56)))) = v4754
	v4757 = v4663
	goto L411
L463:
	;
	goto L410
L464:
	;
	goto L406
L465:
	;
	v4867 = v4811
	goto L466
L466:
	;
	v4917 = *(*int32)(unsafe.Add(mBase, uint32(v4253)+4))
	v4918 = int32(2)
	v4919 = v4917 << (uint(v4918) % 32)
	v4921 = *(*int32)(unsafe.Add(mBase, uint32(v4198+v4919)))
	v4922 = *(*int32)(unsafe.Add(mBase, uint32(v4253)))
	v4925 = v4198 + v4922<<(uint(v4918)%32)
	v4926 = *(*int32)(unsafe.Add(mBase, uint32(v4925)))
	F_HistogramAdd(m, v4921, v4926, v4926)
	mBase = m.M
	v4928 = *(*int32)(unsafe.Add(mBase, uint32(v4925)))
	v4929 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+16))
	*(*int64)(unsafe.Add(mBase, uint32(v4928)+3256)) = v4929
	v4931 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+24))
	*(*int64)(unsafe.Add(mBase, uint32(v4928)+3264)) = v4931
	v4933 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+32))
	*(*int64)(unsafe.Add(mBase, uint32(v4928)+3272)) = v4933
	v4935 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+40))
	*(*int64)(unsafe.Add(mBase, uint32(v4928)+3280)) = v4935
	v4937 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+48))
	*(*int64)(unsafe.Add(mBase, uint32(v4928)+3288)) = v4937
	v4939 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v4939 + int32(-1)
	v4943 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v4950 = *(*int32)(unsafe.Add(mBase, uint32(v4943+v4939<<(uint(v4918)%32)+int32(-4))))
	*(*int32)(unsafe.Add(mBase, uint32(v4943+v4919))) = v4950
	v4952 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+56))
	*(*int64)(unsafe.Add(mBase, uint32(v4928)+3296)) = v4952
	v4956 = v4867
	v4961 = int32(0)
	goto L468
L467:
	;
	goto L403
L468:
	;
	v5008 = v4253 + v4961<<(uint(int32(6))%32)
	v5009 = *(*int32)(unsafe.Add(mBase, uint32(v5008)))
	if v5009 == v4922 {
		goto L472
	} else {
		goto L473
	}
L469:
	;
	v5169 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v5169 < int32(1) {
		v5723 = v5156
		goto L486
	} else {
		goto L487
	}
L470:
	;
	if v5158 < v5156 {
		v4956 = v5156
		v4961 = v5158
		goto L468
	} else {
		goto L485
	}
L471:
	;
	v5049 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	if v5009 != v5049 {
		v5052 = v5009
		goto L477
	} else {
		goto L478
	}
L472:
	;
	v5018 = v4253 + int32(-64) + v4956<<(uint(int32(6))%32)
	v5019 = int32(0)
	v5020 = base.Simd_g_v128_load(m, v5018, v5019)
	base.Simd_g_v128_store(m, v5008, v5019, v5020)
	v5023 = int32(48)
	v5028 = base.Simd_g_v128_load(m, v5018+v5023, v5019)
	base.Simd_g_v128_store(m, v5008+v5023, v5019, v5028)
	v5031 = int32(32)
	v5036 = base.Simd_g_v128_load(m, v5018+v5031, v5019)
	base.Simd_g_v128_store(m, v5008+v5031, v5019, v5036)
	v5039 = int32(16)
	v5044 = base.Simd_g_v128_load(m, v5018+v5039, v5019)
	base.Simd_g_v128_store(m, v5008+v5039, v5019, v5044)
	v5156 = v4956 + int32(-1)
	v5158 = v4961
	goto L470
L473:
	;
	v5011 = *(*int32)(unsafe.Add(mBase, uint32(v5008)+4))
	if v5011 == v4917 {
		goto L472
	} else {
		goto L474
	}
L474:
	;
	if v5009 == v4917 {
		goto L472
	} else {
		goto L475
	}
L475:
	;
	if v5011 != v4922 {
		goto L471
	} else {
		goto L476
	}
L476:
	;
	goto L472
L477:
	;
	if v5011 != v5049 {
		v5055 = v5011
		goto L479
	} else {
		goto L480
	}
L478:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5008))) = v4917
	v5052 = v4917
	goto L477
L479:
	;
	if v5052 <= v5055 {
		goto L481
	} else {
		goto L482
	}
L480:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5008)+4)) = v4917
	v5055 = v4917
	goto L479
L481:
	;
	v5059 = *(*int64)(unsafe.Add(mBase, uint32(v5008)+8))
	v5060 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+8))
	if v5060 <= v5059 {
		goto L483
	} else {
		goto L484
	}
L482:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5008)+4)) = v5052
	*(*int32)(unsafe.Add(mBase, uint32(v5008))) = v5055
	goto L481
L483:
	;
	v5156 = v4956
	v5158 = v4961 + int32(1)
	goto L470
L484:
	;
	v5062 = int32(0)
	v5063 = base.Simd_g_v128_load(m, v4253, v5062)
	v5065 = base.Simd_g_v128_load(m, v5008, v5062)
	base.Simd_g_v128_store(m, v4253, v5062, v5065)
	v5068 = int32(48)
	v5069 = v4253 + v5068
	v5071 = base.Simd_g_v128_load(m, v5069, v5062)
	v5073 = v5008 + v5068
	v5075 = base.Simd_g_v128_load(m, v5073, v5062)
	base.Simd_g_v128_store(m, v5069, v5062, v5075)
	v5078 = int32(32)
	v5079 = v4253 + v5078
	v5081 = base.Simd_g_v128_load(m, v5079, v5062)
	v5083 = v5008 + v5078
	v5085 = base.Simd_g_v128_load(m, v5083, v5062)
	base.Simd_g_v128_store(m, v5079, v5062, v5085)
	v5088 = int32(16)
	v5089 = v4253 + v5088
	v5091 = base.Simd_g_v128_load(m, v5089, v5062)
	v5093 = v5008 + v5088
	v5095 = base.Simd_g_v128_load(m, v5093, v5062)
	base.Simd_g_v128_store(m, v5089, v5062, v5095)
	v5099 = v54 + v5088
	base.Simd_g_v128_store(m, v5099, v5062, v5091)
	v5103 = v54 + v5078
	base.Simd_g_v128_store(m, v5103, v5062, v5081)
	v5107 = v54 + v5068
	base.Simd_g_v128_store(m, v5107, v5062, v5071)
	base.Simd_g_v128_store(m, v54, v5062, v5063)
	v5112 = *(*int64)(unsafe.Add(mBase, uint32(v54)))
	*(*int64)(unsafe.Add(mBase, uint32(v5008))) = v5112
	v5114 = int32(8)
	v5118 = *(*int64)(unsafe.Add(mBase, uint32(v54+v5114)))
	*(*int64)(unsafe.Add(mBase, uint32(v5008+v5114))) = v5118
	v5120 = *(*int64)(unsafe.Add(mBase, uint32(v5099)))
	*(*int64)(unsafe.Add(mBase, uint32(v5093))) = v5120
	v5122 = int32(24)
	v5126 = *(*int64)(unsafe.Add(mBase, uint32(v54+v5122)))
	*(*int64)(unsafe.Add(mBase, uint32(v5008+v5122))) = v5126
	v5128 = *(*int64)(unsafe.Add(mBase, uint32(v5103)))
	*(*int64)(unsafe.Add(mBase, uint32(v5083))) = v5128
	v5130 = *(*int64)(unsafe.Add(mBase, uint32(v5107)))
	*(*int64)(unsafe.Add(mBase, uint32(v5073))) = v5130
	v5132 = int32(56)
	v5136 = *(*int64)(unsafe.Add(mBase, uint32(v54+v5132)))
	*(*int64)(unsafe.Add(mBase, uint32(v5008+v5132))) = v5136
	v5138 = int32(40)
	v5142 = *(*int64)(unsafe.Add(mBase, uint32(v54+v5138)))
	*(*int64)(unsafe.Add(mBase, uint32(v5008+v5138))) = v5142
	goto L483
L485:
	;
	goto L469
L486:
	;
	if int32(1) <= v5723 {
		v4867 = v5723
		goto L466
	} else {
		goto L546
	}
L487:
	;
	v5173 = v5169
	v5181 = v5156
	v5195 = int32(0)
	goto L488
L488:
	;
	if v5195 != v4922 {
		goto L491
	} else {
		goto L492
	}
L489:
	;
	v5723 = v5669
	goto L486
L490:
	;
	v5720 = v5195 + int32(1)
	if v5720 < v5668 {
		v5173 = v5668
		v5181 = v5676
		v5195 = v5720
		goto L488
	} else {
		goto L545
	}
L491:
	;
	if v5181 == v4228 {
		v5617 = v4228
		goto L493
	} else {
		goto L494
	}
L492:
	;
	v5668 = v5173
	v5669 = v5181
	v5676 = v5181
	goto L490
L493:
	;
	v5667 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	v5668 = v5667
	v5669 = v5617
	v5676 = v5617
	goto L490
L494:
	;
	v5226 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	if v4922 < v5195 {
		goto L495
	} else {
		goto L496
	}
L495:
	;
	v5228 = v4922
	goto L497
L496:
	;
	v5228 = v5195
	goto L497
L497:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+256)) = v5228
	v5233 = *(*int32)(unsafe.Add(mBase, uint32(v5226+v5228<<(uint(int32(2))%32))))
	v5234 = *(*int64)(unsafe.Add(mBase, uint32(v5233)+3256))
	if v5195 < v4922 {
		goto L498
	} else {
		goto L499
	}
L498:
	;
	v5236 = v4922
	goto L500
L499:
	;
	v5236 = v5195
	goto L500
L500:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v54)+260)) = v5236
	v5241 = *(*int32)(unsafe.Add(mBase, uint32(v5226+v5236<<(uint(int32(2))%32))))
	v5242 = *(*int64)(unsafe.Add(mBase, uint32(v5241)+3256))
	v5243 = v5234 + v5242
	if int64(1) <= v5243 {
		goto L501
	} else {
		goto L502
	}
L501:
	;
	v5246 = int32(4)
	v5250 = int32(1028)
	v5254 = int32(2052)
	v5258 = int32(3076)
	v5266 = int32(3240)
	v5267 = int32(0)
	v5281 = int32(3264)
	v5296 = int64(0)
	goto L504
L502:
	;
	v5617 = v5181
	goto L493
L503:
	;
	v5472 = int32(6)
	v5474 = v4253 + v5181<<(uint(v5472)%32)
	v5481 = int32(0)
	v5482 = base.Simd_g_v128_load(m, v54+int32(304), v5481)
	base.Simd_g_v128_store(m, v5474+int32(48), v5481, v5482)
	v5492 = base.Simd_g_v128_load(m, v54+int32(288), v5481)
	base.Simd_g_v128_store(m, v5474+int32(32), v5481, v5492)
	v5501 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(280))))
	*(*int64)(unsafe.Add(mBase, uint32(v5474+int32(24)))) = v5501
	*(*int64)(unsafe.Add(mBase, uint32(v5474+int32(16)))) = v5462
	v5508 = v5462 - v5243
	*(*int64)(unsafe.Add(mBase, uint32(v5474+int32(8)))) = v5508
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(264)))) = v5508
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(272)))) = v5462
	v5520 = *(*int64)(unsafe.Add(mBase, uint32(v54)+256))
	*(*int64)(unsafe.Add(mBase, uint32(v5474))) = v5520
	v5523 = v5181 + int32(1)
	v5526 = v4253 + v5523<<(uint(v5472)%32)
	v5529 = *(*int64)(unsafe.Add(mBase, uint32(v5526+int32(-56))))
	v5530 = *(*int64)(unsafe.Add(mBase, uint32(v4253)+8))
	if v5530 <= v5529 {
		v5617 = v5523
		goto L493
	} else {
		goto L544
	}
L504:
	;
	v5321 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v5233+v5266))))
	if v5321 != int32(_a_F_VP8LGetHistoImageSymbols_2) {
		goto L507
	} else {
		goto L508
	}
L505:
	;
	v5617 = v5181
	goto L493
L506:
	;
	v5329 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5233+v5267+int32(3304)))))
	if v5328 == int32(0) {
		goto L511
	} else {
		goto L512
	}
L507:
	;
	v5326 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v5241+v5266))))
	v5328 = base.B2i32(v5321 != v5326)
	goto L506
L508:
	;
	v5328 = int32(1)
	goto L506
L509:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v54+int32(256)+v5281+int32(-3240)))) = v5454
	v5462 = v5454 + v5296
	if base.Ui64(v5243) <= base.Ui64(v5462) {
		goto L541
	} else {
		goto L542
	}
L510:
	;
	v5350 = int32(256)
	switch v5267 {
	default:
		goto L521
	case 1:
		v5363 = v5233 + v5246
		v5364 = v5350
		v5365 = v5241 + v5246
		goto L517
	case 2:
		goto L520
	case 3:
		goto L519
	case 4:
		goto L518
	}
L511:
	;
	if v5329&int32(255) == int32(0) {
		goto L515
	} else {
		goto L516
	}
L512:
	;
	if v5329&int32(255) == int32(0) {
		goto L511
	} else {
		goto L513
	}
L513:
	;
	v5339 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v5241+v5267+int32(3304)))))
	if v5339&int32(255) != 0 {
		goto L510
	} else {
		goto L514
	}
L514:
	;
	goto L511
L515:
	;
	v5349 = *(*int64)(unsafe.Add(mBase, uint32(v5241+v5281)))
	v5454 = v5349
	goto L509
L516:
	;
	v5347 = *(*int64)(unsafe.Add(mBase, uint32(v5233+v5281)))
	v5454 = v5347
	goto L509
L517:
	;
	v5368 = m.G115
	v5369 = *(*int32)(unsafe.Add(mBase, uint32(v5368)))
	m.T0[v5369].(func(*base.Module, int32, int32, int32, int32, int32))(m, v5363, v5365, v5364, v54+int32(328), v54)
	mBase = m.M
	v5371 = *(*int32)(unsafe.Add(mBase, uint32(v54)+340))
	if v5371 <= int32(4) {
		goto L527
	} else {
		goto L528
	}
L518:
	;
	v5363 = v5233 + v5258
	v5364 = int32(40)
	v5365 = v5241 + v5258
	goto L517
L519:
	;
	v5363 = v5233 + v5254
	v5364 = v5350
	v5365 = v5241 + v5254
	goto L517
L520:
	;
	v5363 = v5233 + v5250
	v5364 = v5350
	v5365 = v5241 + v5250
	goto L517
L521:
	;
	v5352 = *(*int32)(unsafe.Add(mBase, uint32(v5241)+3236))
	v5354 = int32(280)
	if int32(0) < v5352 {
		goto L522
	} else {
		goto L523
	}
L522:
	;
	v5359 = int32(1)<<(uint(v5352)%32) + v5354
	goto L524
L523:
	;
	v5359 = v5354
	goto L524
L524:
	;
	v5360 = *(*int32)(unsafe.Add(mBase, uint32(v5241)))
	v5361 = *(*int32)(unsafe.Add(mBase, uint32(v5233)))
	v5363 = v5361
	v5364 = v5359
	v5365 = v5360
	goto L517
L525:
	;
	v5422 = *(*int32)(unsafe.Add(mBase, uint32(v54)+12))
	v5425 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	v5429 = *(*int32)(unsafe.Add(mBase, uint32(v54)+4))
	v5433 = *(*int32)(unsafe.Add(mBase, uint32(v54)+20))
	v5437 = *(*int32)(unsafe.Add(mBase, uint32(v54)+8))
	v5441 = *(*int32)(unsafe.Add(mBase, uint32(v54)+16))
	v5454 = v5420 + base.I64_extend_i32_u(v5422*int32(240)+v5425*int32(1600)+v5429*int32(2640)+v5433*int32(720)+v5437*int32(1840)+v5441*int32(3360))<<(uint(int64(13))%64) + int64(401814323)
	goto L509
L526:
	;
	v5396 = *(*int64)(unsafe.Add(mBase, uint32(v54)+328))
	v5399 = *(*int32)(unsafe.Add(mBase, uint32(v54)+336))
	v5402 = *(*int32)(unsafe.Add(mBase, uint32(v54)+344))
	v5411 = v5395*base.I64_extend_i32_u(v5399<<(uint(int32(1))%32)-v5402)<<(uint(int64(23))%64) + v5396*(int64(1000)-v5395)
	if v5411 < int64(0) {
		goto L535
	} else {
		goto L536
	}
L527:
	;
	if v5371 < int32(2) {
		v5420 = int64(0)
		goto L525
	} else {
		goto L529
	}
L528:
	;
	v5395 = int64(627)
	goto L526
L529:
	;
	switch v5371 + int32(-2) {
	case 0:
		goto L531
	case 1:
		v5395 = int64(950)
		goto L526
	default:
		goto L530
	}
L530:
	;
	v5395 = int64(700)
	goto L526
L531:
	;
	v5383 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v54)+336)))
	v5386 = *(*int64)(unsafe.Add(mBase, uint32(v54)+328))
	v5387 = v5383*int64(830472192) + v5386
	if v5387 < int64(0) {
		goto L532
	} else {
		goto L533
	}
L532:
	;
	v5390 = int64(-50)
	goto L534
L533:
	;
	v5390 = int64(50)
	goto L534
L534:
	;
	v5393 = base.I64_div_s(v5390+v5387, int64(100))
	v5420 = v5393
	goto L525
L535:
	;
	v5414 = int64(-500)
	goto L537
L536:
	;
	v5414 = int64(500)
	goto L537
L537:
	;
	v5417 = base.I64_div_s(v5414+v5411, int64(1000))
	if base.Ui64(v5417) < base.Ui64(v5396) {
		goto L538
	} else {
		goto L539
	}
L538:
	;
	v5419 = v5396
	goto L540
L539:
	;
	v5419 = v5417
	goto L540
L540:
	;
	v5420 = v5419
	goto L525
L541:
	;
	goto L505
L542:
	;
	v5469 = v5267 + int32(1)
	if v5469 == int32(5) {
		goto L503
	} else {
		goto L543
	}
L543:
	;
	v5266 = v5266 + int32(2)
	v5267 = v5469
	v5281 = v5281 + int32(8)
	v5296 = v5462
	goto L504
L544:
	;
	v5532 = int32(0)
	v5533 = base.Simd_g_v128_load(m, v4253, v5532)
	v5535 = v5526 + int32(-64)
	v5537 = base.Simd_g_v128_load(m, v5535, v5532)
	base.Simd_g_v128_store(m, v4253, v5532, v5537)
	v5540 = int32(48)
	v5541 = v4253 + v5540
	v5543 = base.Simd_g_v128_load(m, v5541, v5532)
	v5545 = v5526 + int32(-16)
	v5547 = base.Simd_g_v128_load(m, v5545, v5532)
	base.Simd_g_v128_store(m, v5541, v5532, v5547)
	v5550 = int32(32)
	v5551 = v4253 + v5550
	v5553 = base.Simd_g_v128_load(m, v5551, v5532)
	v5555 = v5526 + int32(-32)
	v5557 = base.Simd_g_v128_load(m, v5555, v5532)
	base.Simd_g_v128_store(m, v5551, v5532, v5557)
	v5560 = int32(16)
	v5561 = v4253 + v5560
	v5563 = base.Simd_g_v128_load(m, v5561, v5532)
	v5565 = v5526 + int32(-48)
	v5567 = base.Simd_g_v128_load(m, v5565, v5532)
	base.Simd_g_v128_store(m, v5561, v5532, v5567)
	v5571 = v54 + v5560
	base.Simd_g_v128_store(m, v5571, v5532, v5563)
	v5575 = v54 + v5550
	base.Simd_g_v128_store(m, v5575, v5532, v5553)
	v5579 = v54 + v5540
	base.Simd_g_v128_store(m, v5579, v5532, v5543)
	base.Simd_g_v128_store(m, v54, v5532, v5533)
	v5584 = *(*int64)(unsafe.Add(mBase, uint32(v54)))
	*(*int64)(unsafe.Add(mBase, uint32(v5535))) = v5584
	v5586 = *(*int64)(unsafe.Add(mBase, uint32(v5571)))
	*(*int64)(unsafe.Add(mBase, uint32(v5565))) = v5586
	v5592 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(24))))
	*(*int64)(unsafe.Add(mBase, uint32(v5526+int32(-40)))) = v5592
	v5594 = *(*int64)(unsafe.Add(mBase, uint32(v5575)))
	*(*int64)(unsafe.Add(mBase, uint32(v5555))) = v5594
	v5596 = *(*int64)(unsafe.Add(mBase, uint32(v5579)))
	*(*int64)(unsafe.Add(mBase, uint32(v5545))) = v5596
	v5602 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(56))))
	*(*int64)(unsafe.Add(mBase, uint32(v5526+int32(-8)))) = v5602
	v5608 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(40))))
	*(*int64)(unsafe.Add(mBase, uint32(v5526+int32(-24)))) = v5608
	v5614 = *(*int64)(unsafe.Add(mBase, uint32(v54+int32(8))))
	*(*int64)(unsafe.Add(mBase, uint32(v5526+int32(-56)))) = v5614
	v5617 = v5523
	goto L493
L545:
	;
	goto L489
L546:
	;
	goto L467
L547:
	;
	v5827 = *(*int32)(unsafe.Add(mBase, uint32(l7)))
	v5848 = v5827
	goto L6
L548:
	;
	v6579 = int32(0)
	v6582 = *(*int32)(unsafe.Add(mBase, uint32(v6529)))
	v6583 = *(*int32)(unsafe.Add(mBase, uint32(v6582)+3236))
	v6584 = int32(4) << (uint(v6583) % 32)
	v6585 = int32(_a_F_VP8LGetHistoImageSymbols_0)
	v6589 = base.B2i32(v6579 < v6583)
	if v6579 < v6583 {
		goto L624
	} else {
		goto L625
	}
L549:
	;
	v6527 = *(*int32)(unsafe.Add(mBase, uint32(l7)+4))
	v6529 = v6477
	v6550 = v6527
	goto L548
L550:
	;
	if v5879 < int32(1) {
		v6529 = v5880
		v6550 = v5879
		goto L548
	} else {
		goto L567
	}
L551:
	;
	if v5879 < int32(1) {
		v6529 = v5880
		v6550 = v5879
		goto L548
	} else {
		goto L552
	}
L552:
	;
	v5886 = int32(0)
	v5888 = v5879 << (uint(int32(2)) % 32)
	if base.Ui32(v5888) < base.Ui32(int32(33)) {
		goto L554
	} else {
		goto L555
	}
L553:
	;
	v6477 = v5880
	goto L549
L554:
	;
	if v5888 == int32(0) {
		goto L556
	} else {
		goto L557
	}
L555:
	;
	base.MemoryFill(m, l9, v5886, v5888)
	goto L553
L556:
	;
	goto L553
L557:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l9))) = uint8(v5886)
	v5899 = l9 + v5888
	*(*uint8)(unsafe.Add(mBase, uint32(v5899+int32(-1)))) = uint8(v5886)
	if base.Ui32(v5888) < base.Ui32(int32(3)) {
		goto L556
	} else {
		goto L558
	}
L558:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l9)+2)) = uint8(v5886)
	*(*uint8)(unsafe.Add(mBase, uint32(l9)+1)) = uint8(v5886)
	*(*uint8)(unsafe.Add(mBase, uint32(v5899+int32(-3)))) = uint8(v5886)
	*(*uint8)(unsafe.Add(mBase, uint32(v5899+int32(-2)))) = uint8(v5886)
	if base.Ui32(v5888) < base.Ui32(int32(7)) {
		goto L556
	} else {
		goto L559
	}
L559:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l9)+3)) = uint8(v5886)
	*(*uint8)(unsafe.Add(mBase, uint32(v5899+int32(-4)))) = uint8(v5886)
	if base.Ui32(v5888) < base.Ui32(int32(9)) {
		goto L556
	} else {
		goto L560
	}
L560:
	;
	v5921 = int32(0)
	v5924 = (v5921 - l9) & int32(3)
	v5925 = l9 + v5924
	*(*int32)(unsafe.Add(mBase, uint32(v5925))) = v5921
	v5933 = (v5888 - v5924) & int32(60)
	v5934 = v5925 + v5933
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-4)))) = v5921
	if base.Ui32(v5933) < base.Ui32(int32(9)) {
		goto L556
	} else {
		goto L561
	}
L561:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5925)+8)) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5925)+4)) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-8)))) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-12)))) = v5921
	if base.Ui32(v5933) < base.Ui32(int32(25)) {
		goto L556
	} else {
		goto L562
	}
L562:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v5925)+24)) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5925)+20)) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5925)+16)) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5925)+12)) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-16)))) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-20)))) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-24)))) = v5921
	*(*int32)(unsafe.Add(mBase, uint32(v5934+int32(-28)))) = v5921
	v5969 = v5925&int32(4) | int32(24)
	v5970 = v5933 - v5969
	if base.Ui32(v5970) < base.Ui32(int32(32)) {
		goto L556
	} else {
		goto L563
	}
L563:
	;
	v5975 = base.I64_extend_i32_u(v5921) * int64(4294967297)
	v5978 = v5970
	v5979 = v5925 + v5969
	goto L564
L564:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v5979)+24)) = v5975
	*(*int64)(unsafe.Add(mBase, uint32(v5979)+16)) = v5975
	*(*int64)(unsafe.Add(mBase, uint32(v5979)+8)) = v5975
	*(*int64)(unsafe.Add(mBase, uint32(v5979))) = v5975
	v5991 = v5978 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v5991) {
		v5978 = v5991
		v5979 = v5979 + int32(32)
		goto L564
	} else {
		goto L566
	}
L565:
	;
	goto L556
L566:
	;
	goto L565
L567:
	;
	v6037 = int32(0)
	goto L568
L568:
	;
	v6064 = v6037 << (uint(int32(2)) % 32)
	v6065 = v5881 + v6064
	v6066 = *(*int32)(unsafe.Add(mBase, uint32(v6065)))
	if v6066 == int32(0) {
		goto L571
	} else {
		goto L572
	}
L569:
	;
	v6475 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v6477 = v6475
	goto L549
L570:
	;
	v6473 = v6037 + int32(1)
	if v6473 != v5879 {
		v6037 = v6473
		goto L568
	} else {
		goto L623
	}
L571:
	;
	v6416 = l9 + v6064
	v6419 = *(*int32)(unsafe.Add(mBase, uint32(v6416+int32(-4))))
	*(*int32)(unsafe.Add(mBase, uint32(v6416))) = v6419
	goto L570
L572:
	;
	v6070 = int32(0)
	v6086 = v6070
	v6101 = int64(9223372036854775807)
	v6117 = v6070
	goto L573
L573:
	;
	v6123 = int64(9223372036854775807)
	v6127 = *(*int32)(unsafe.Add(mBase, uint32(v5880+v6086<<(uint(int32(2))%32))))
	v6128 = *(*int64)(unsafe.Add(mBase, uint32(v6127)+3256))
	v6129 = v6128 + v6101
	if v6101^v6123 < v6128 {
		goto L576
	} else {
		goto L577
	}
L574:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l9+v6064))) = v6405
	goto L570
L575:
	;
	v6412 = v6086 + int32(1)
	if v6412 != v5848 {
		v6086 = v6412
		v6101 = v6389
		v6117 = v6405
		goto L573
	} else {
		goto L622
	}
L576:
	;
	v6133 = v6123
	goto L578
L577:
	;
	v6133 = v6129
	goto L578
L578:
	;
	if int64(-1) < v6101 {
		goto L579
	} else {
		goto L580
	}
L579:
	;
	v6136 = v6133
	goto L581
L580:
	;
	v6136 = v6129
	goto L581
L581:
	;
	if v6136 < int64(1) {
		v6389 = v6101
		v6405 = v6117
		goto L575
	} else {
		goto L582
	}
L582:
	;
	v6139 = *(*int32)(unsafe.Add(mBase, uint32(v6065)))
	v6140 = int32(4)
	v6144 = int32(1028)
	v6148 = int32(2052)
	v6152 = int32(3076)
	v6160 = int32(3240)
	v6161 = int32(0)
	v6166 = int32(3264)
	v6186 = int64(0)
	goto L583
L583:
	;
	v6215 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v6127+v6160))))
	if v6215 != int32(_a_F_VP8LGetHistoImageSymbols_2) {
		goto L586
	} else {
		goto L587
	}
L584:
	;
	v6358 = *(*int64)(unsafe.Add(mBase, uint32(v6127)+3256))
	v6389 = v6348 - v6358
	v6405 = v6086
	goto L575
L585:
	;
	v6223 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v6127+v6161+int32(3304)))))
	if v6222 == int32(0) {
		goto L590
	} else {
		goto L591
	}
L586:
	;
	v6220 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v6139+v6160))))
	v6222 = base.B2i32(v6215 != v6220)
	goto L585
L587:
	;
	v6222 = int32(1)
	goto L585
L588:
	;
	v6348 = v6347 + v6186
	if base.Ui64(v6136) <= base.Ui64(v6348) {
		v6389 = v6101
		v6405 = v6117
		goto L575
	} else {
		goto L620
	}
L589:
	;
	v6242 = int32(256)
	switch v6161 {
	default:
		goto L597
	case 1:
		v6255 = v6127 + v6140
		v6256 = v6242
		v6257 = v6139 + v6140
		goto L596
	case 2:
		goto L598
	case 3:
		goto L599
	case 4:
		goto L600
	}
L590:
	;
	if v6223&int32(255) != 0 {
		goto L594
	} else {
		goto L595
	}
L591:
	;
	if v6223&int32(255) == int32(0) {
		goto L590
	} else {
		goto L592
	}
L592:
	;
	v6233 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v6139+v6161+int32(3304)))))
	if v6233&int32(255) != 0 {
		goto L589
	} else {
		goto L593
	}
L593:
	;
	goto L590
L594:
	;
	v6241 = *(*int64)(unsafe.Add(mBase, uint32(v6127+v6166)))
	v6347 = v6241
	goto L588
L595:
	;
	v6239 = *(*int64)(unsafe.Add(mBase, uint32(v6139+v6166)))
	v6347 = v6239
	goto L588
L596:
	;
	v6260 = m.G115
	v6261 = *(*int32)(unsafe.Add(mBase, uint32(v6260)))
	m.T0[v6261].(func(*base.Module, int32, int32, int32, int32, int32))(m, v6255, v6257, v6256, v54+int32(256), v54)
	mBase = m.M
	v6263 = *(*int32)(unsafe.Add(mBase, uint32(v54)+268))
	if v6263 <= int32(4) {
		goto L606
	} else {
		goto L607
	}
L597:
	;
	v6245 = *(*int32)(unsafe.Add(mBase, uint32(v6139)+3236))
	v6247 = int32(280)
	if int32(0) < v6245 {
		goto L601
	} else {
		goto L602
	}
L598:
	;
	v6255 = v6127 + v6144
	v6256 = v6242
	v6257 = v6139 + v6144
	goto L596
L599:
	;
	v6255 = v6127 + v6148
	v6256 = v6242
	v6257 = v6139 + v6148
	goto L596
L600:
	;
	v6255 = v6127 + v6152
	v6256 = int32(40)
	v6257 = v6139 + v6152
	goto L596
L601:
	;
	v6252 = int32(1)<<(uint(v6245)%32) + v6247
	goto L603
L602:
	;
	v6252 = v6247
	goto L603
L603:
	;
	v6253 = *(*int32)(unsafe.Add(mBase, uint32(v6139)))
	v6254 = *(*int32)(unsafe.Add(mBase, uint32(v6127)))
	v6255 = v6254
	v6256 = v6252
	v6257 = v6253
	goto L596
L604:
	;
	v6314 = *(*int32)(unsafe.Add(mBase, uint32(v54)+12))
	v6317 = *(*int32)(unsafe.Add(mBase, uint32(v54)))
	v6321 = *(*int32)(unsafe.Add(mBase, uint32(v54)+4))
	v6325 = *(*int32)(unsafe.Add(mBase, uint32(v54)+20))
	v6329 = *(*int32)(unsafe.Add(mBase, uint32(v54)+8))
	v6333 = *(*int32)(unsafe.Add(mBase, uint32(v54)+16))
	v6347 = v6313 + base.I64_extend_i32_u(v6314*int32(240)+v6317*int32(1600)+v6321*int32(2640)+v6325*int32(720)+v6329*int32(1840)+v6333*int32(3360))<<(uint(int64(13))%64) + int64(401814323)
	goto L588
L605:
	;
	v6288 = *(*int64)(unsafe.Add(mBase, uint32(v54)+256))
	v6291 = *(*int32)(unsafe.Add(mBase, uint32(v54)+264))
	v6294 = *(*int32)(unsafe.Add(mBase, uint32(v54)+272))
	v6303 = v6287*base.I64_extend_i32_u(v6291<<(uint(int32(1))%32)-v6294)<<(uint(int64(23))%64) + v6288*(int64(1000)-v6287)
	if v6303 < int64(0) {
		goto L614
	} else {
		goto L615
	}
L606:
	;
	if v6263 < int32(2) {
		v6313 = int64(0)
		goto L604
	} else {
		goto L608
	}
L607:
	;
	v6287 = int64(627)
	goto L605
L608:
	;
	switch v6263 + int32(-2) {
	case 0:
		goto L610
	case 1:
		v6287 = int64(950)
		goto L605
	default:
		goto L609
	}
L609:
	;
	v6287 = int64(700)
	goto L605
L610:
	;
	v6275 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v54)+264)))
	v6278 = *(*int64)(unsafe.Add(mBase, uint32(v54)+256))
	v6279 = v6275*int64(830472192) + v6278
	if v6279 < int64(0) {
		goto L611
	} else {
		goto L612
	}
L611:
	;
	v6282 = int64(-50)
	goto L613
L612:
	;
	v6282 = int64(50)
	goto L613
L613:
	;
	v6285 = base.I64_div_s(v6282+v6279, int64(100))
	v6313 = v6285
	goto L604
L614:
	;
	v6306 = int64(-500)
	goto L616
L615:
	;
	v6306 = int64(500)
	goto L616
L616:
	;
	v6309 = base.I64_div_s(v6306+v6303, int64(1000))
	if base.Ui64(v6309) < base.Ui64(v6288) {
		goto L617
	} else {
		goto L618
	}
L617:
	;
	v6311 = v6288
	goto L619
L618:
	;
	v6311 = v6309
	goto L619
L619:
	;
	v6313 = v6311
	goto L604
L620:
	;
	v6355 = v6161 + int32(1)
	if v6355 != int32(5) {
		v6160 = v6160 + int32(2)
		v6161 = v6355
		v6166 = v6166 + int32(8)
		v6186 = v6348
		goto L583
	} else {
		goto L621
	}
L621:
	;
	goto L584
L622:
	;
	goto L574
L623:
	;
	goto L569
L624:
	;
	v6590 = v6584 + v6585
	goto L626
L625:
	;
	v6590 = v6585
	goto L626
L626:
	;
	v6593 = v6590*v6550 + int32(12)
	if base.Ui32(v6593) < base.Ui32(int32(33)) {
		goto L628
	} else {
		goto L629
	}
L627:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l7)+4)) = v6550
	v6716 = l7 + int32(12)
	*(*int32)(unsafe.Add(mBase, uint32(l7)+8)) = v6716
	if v6550 < int32(1) {
		goto L641
	} else {
		goto L642
	}
L628:
	;
	if v6593 == int32(0) {
		goto L630
	} else {
		goto L631
	}
L629:
	;
	base.MemoryFill(m, l7, v6579, v6593)
	goto L627
L630:
	;
	goto L627
L631:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l7))) = uint8(v6579)
	v6604 = l7 + v6593
	*(*uint8)(unsafe.Add(mBase, uint32(v6604+int32(-1)))) = uint8(v6579)
	if base.Ui32(v6593) < base.Ui32(int32(3)) {
		goto L630
	} else {
		goto L632
	}
L632:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l7)+2)) = uint8(v6579)
	*(*uint8)(unsafe.Add(mBase, uint32(l7)+1)) = uint8(v6579)
	*(*uint8)(unsafe.Add(mBase, uint32(v6604+int32(-3)))) = uint8(v6579)
	*(*uint8)(unsafe.Add(mBase, uint32(v6604+int32(-2)))) = uint8(v6579)
	if base.Ui32(v6593) < base.Ui32(int32(7)) {
		goto L630
	} else {
		goto L633
	}
L633:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(l7)+3)) = uint8(v6579)
	*(*uint8)(unsafe.Add(mBase, uint32(v6604+int32(-4)))) = uint8(v6579)
	if base.Ui32(v6593) < base.Ui32(int32(9)) {
		goto L630
	} else {
		goto L634
	}
L634:
	;
	v6626 = int32(0)
	v6629 = (v6626 - l7) & int32(3)
	v6630 = l7 + v6629
	*(*int32)(unsafe.Add(mBase, uint32(v6630))) = v6626
	v6638 = (v6593 - v6629) & int32(60)
	v6639 = v6630 + v6638
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-4)))) = v6626
	if base.Ui32(v6638) < base.Ui32(int32(9)) {
		goto L630
	} else {
		goto L635
	}
L635:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v6630)+8)) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6630)+4)) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-8)))) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-12)))) = v6626
	if base.Ui32(v6638) < base.Ui32(int32(25)) {
		goto L630
	} else {
		goto L636
	}
L636:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v6630)+24)) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6630)+20)) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6630)+16)) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6630)+12)) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-16)))) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-20)))) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-24)))) = v6626
	*(*int32)(unsafe.Add(mBase, uint32(v6639+int32(-28)))) = v6626
	v6674 = v6630&int32(4) | int32(24)
	v6675 = v6638 - v6674
	if base.Ui32(v6675) < base.Ui32(int32(32)) {
		goto L630
	} else {
		goto L637
	}
L637:
	;
	v6680 = base.I64_extend_i32_u(v6626) * int64(4294967297)
	v6683 = v6675
	v6684 = v6630 + v6674
	goto L638
L638:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v6684)+24)) = v6680
	*(*int64)(unsafe.Add(mBase, uint32(v6684)+16)) = v6680
	*(*int64)(unsafe.Add(mBase, uint32(v6684)+8)) = v6680
	*(*int64)(unsafe.Add(mBase, uint32(v6684))) = v6680
	v6696 = v6683 + int32(-32)
	if base.Ui32(int32(31)) < base.Ui32(v6696) {
		v6683 = v6696
		v6684 = v6684 + int32(32)
		goto L638
	} else {
		goto L640
	}
L639:
	;
	goto L630
L640:
	;
	goto L639
L641:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l7))) = v5848
	if v5879 < int32(1) {
		goto L662
	} else {
		goto L663
	}
L642:
	;
	v6720 = int32(1)
	v6724 = v6716 + v6550<<(uint(int32(2))%32)
	if v6550 == v6720 {
		v6825 = v6716
		v6835 = v6724
		v6838 = v6579
		goto L643
	} else {
		goto L644
	}
L643:
	;
	if v6550&v6720 == int32(0) {
		v6888 = v6825
		goto L651
	} else {
		goto L652
	}
L644:
	;
	v6727 = int32(_a_F_VP8LGetHistoImageSymbols_1)
	if v6579 < v6583 {
		goto L645
	} else {
		goto L646
	}
L645:
	;
	v6730 = v6584 + v6727
	goto L647
L646:
	;
	v6730 = v6727
	goto L647
L647:
	;
	v6736 = int32(4)
	v6740 = v6716
	v6750 = v6724
	v6753 = int32(0)
	goto L648
L648:
	;
	v6787 = int32(-4)
	v6789 = int32(31)
	v6791 = int32(-32)
	v6792 = (v6750 + v6789) & v6791
	*(*int32)(unsafe.Add(mBase, uint32(v6740+v6736+v6787))) = v6792
	v6794 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v6795 = v6794 + v6736
	v6798 = *(*int32)(unsafe.Add(mBase, uint32(v6795+v6787)))
	v6799 = int32(3312)
	*(*int32)(unsafe.Add(mBase, uint32(v6798))) = v6792 + v6799
	v6806 = (v6792 + v6730 + v6789) & v6791
	*(*int32)(unsafe.Add(mBase, uint32(v6795))) = v6806
	v6808 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v6810 = *(*int32)(unsafe.Add(mBase, uint32(v6808+v6736)))
	*(*int32)(unsafe.Add(mBase, uint32(v6810))) = v6806 + v6799
	v6816 = v6806 + v6730
	v6818 = v6753 + int32(2)
	if v6550&int32(2147483646) != v6818 {
		v6736 = v6736 + int32(8)
		v6740 = v6808
		v6750 = v6816
		v6753 = v6818
		goto L648
	} else {
		goto L650
	}
L649:
	;
	v6825 = v6808
	v6835 = v6816
	v6838 = v6818
	goto L643
L650:
	;
	goto L649
L651:
	;
	if base.Ui32(v6550) < base.Ui32(int32(4)) {
		v6988 = int32(0)
		goto L653
	} else {
		goto L654
	}
L652:
	;
	v6874 = v6838 << (uint(int32(2)) % 32)
	v6879 = (v6835 + int32(31)) & int32(-32)
	*(*int32)(unsafe.Add(mBase, uint32(v6825+v6874))) = v6879
	v6881 = *(*int32)(unsafe.Add(mBase, uint32(l7)+8))
	v6883 = *(*int32)(unsafe.Add(mBase, uint32(v6881+v6874)))
	*(*int32)(unsafe.Add(mBase, uint32(v6883))) = v6879 + int32(3312)
	v6888 = v6881
	goto L651
L653:
	;
	v7021 = v6888 + v6988<<(uint(int32(2))%32)
	v7035 = v6550 - v6988
	goto L659
L654:
	;
	v6894 = v6550 & int32(2147483644)
	v6896 = v6888
	v6910 = v6894
	goto L655
L655:
	;
	v6948 = *(*int32)(unsafe.Add(mBase, uint32(v6896+int32(12))))
	v6951 = *(*int32)(unsafe.Add(mBase, uint32(v6896+int32(8))))
	v6954 = *(*int32)(unsafe.Add(mBase, uint32(v6896+int32(4))))
	v6955 = *(*int32)(unsafe.Add(mBase, uint32(v6896)))
	*(*int32)(unsafe.Add(mBase, uint32(v6955)+3236)) = v6583
	*(*int32)(unsafe.Add(mBase, uint32(v6954)+3236)) = v6583
	*(*int32)(unsafe.Add(mBase, uint32(v6951)+3236)) = v6583
	*(*int32)(unsafe.Add(mBase, uint32(v6948)+3236)) = v6583
	v6963 = v6910 + int32(-4)
	if v6963 != 0 {
		v6896 = v6896 + int32(16)
		v6910 = v6963
		goto L655
	} else {
		goto L657
	}
L656:
	;
	if v6550 == v6894 {
		goto L641
	} else {
		goto L658
	}
L657:
	;
	goto L656
L658:
	;
	v6988 = v6894
	goto L653
L659:
	;
	v7071 = *(*int32)(unsafe.Add(mBase, uint32(v7021)))
	*(*int32)(unsafe.Add(mBase, uint32(v7071)+3236)) = v6583
	v7076 = v7035 + int32(-1)
	if v7076 != 0 {
		v7021 = v7021 + int32(4)
		v7035 = v7076
		goto L659
	} else {
		goto L661
	}
L660:
	;
	goto L641
L661:
	;
	goto L660
L662:
	;
	v7249 = *(*int32)(unsafe.Add(mBase, uint32(l12)))
	v7251 = F_WebPReportProgress(m, l10, v7249+l11, l12)
	mBase = m.M
	goto L3
L663:
	;
	v7140 = l9
	v7175 = v5881
	v7177 = v5879
	goto L664
L664:
	;
	v7182 = *(*int32)(unsafe.Add(mBase, uint32(v7175)))
	if v7182 == int32(0) {
		goto L666
	} else {
		goto L667
	}
L665:
	;
	goto L662
L666:
	;
	v7192 = int32(4)
	v7197 = v7177 + int32(-1)
	if v7197 != 0 {
		v7140 = v7140 + v7192
		v7175 = v7175 + v7192
		v7177 = v7197
		goto L664
	} else {
		goto L668
	}
L667:
	;
	v7185 = *(*int32)(unsafe.Add(mBase, uint32(v7140)))
	v7189 = *(*int32)(unsafe.Add(mBase, uint32(v5880+v7185<<(uint(int32(2))%32))))
	F_HistogramAdd(m, v7182, v7189, v7189)
	mBase = m.M
	goto L666
L668:
	;
	goto L665
L669:
	;
	v7255 = *(*int32)(unsafe.Add(mBase, uint32(l10)+92))
	if v7255 != 0 {
		goto L671
	} else {
		goto L672
	}
L670:
	;
	goto L3
L671:
	;
	goto L670
L672:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l10)+92)) = int32(1)
	goto L671
L673:
	;
	v7261 = *(*int32)(unsafe.Add(mBase, uint32(l10)+92))
	if v7261 != 0 {
		goto L675
	} else {
		goto L676
	}
L674:
	;
	goto L3
L675:
	;
	goto L674
L676:
	;
	*(*int32)(unsafe.Add(mBase, uint32(l10)+92)) = int32(1)
	goto L675
L677:
	;
	v7316 = *(*int32)(unsafe.Add(mBase, uint32(l10)+92))
	m.G0 = v54 + int32(352)
	return base.B2i32(v7316 == int32(0))
}

var F_VP8LGetHistoImageSymbols__k0 = [2]uint64{0x0, 0x0}
var F_VP8LGetHistoImageSymbols__k1 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}
var F_VP8LGetHistoImageSymbols__k2 = [2]uint64{0xf0e0d0c0b0a0908, 0x706050403020100}

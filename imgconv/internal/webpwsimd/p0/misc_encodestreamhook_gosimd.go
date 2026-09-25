//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_EncodeStreamHook(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v3 int32
	_ = v3
	var v72 int32
	_ = v72
	var v74 int32
	_ = v74
	var v76 int32
	_ = v76
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 int32
	_ = v83
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v87 int32
	_ = v87
	var v88 float32
	_ = v88
	var v89 int32
	_ = v89
	var v90 int32
	_ = v90
	var v91 int32
	_ = v91
	var v100 int64
	_ = v100
	var v103 base.V128
	_ = v103
	var v111 int32
	_ = v111
	var v113 int32
	_ = v113
	var v115 int32
	_ = v115
	var v118 base.V128
	_ = v118
	var v128 int32
	_ = v128
	var v130 int32
	_ = v130
	var v134 int32
	_ = v134
	var v141 int32
	_ = v141
	var v147 int32
	_ = v147
	var v157 int32
	_ = v157
	var v158 int32
	_ = v158
	var v159 int32
	_ = v159
	var v161 int32
	_ = v161
	var v162 int32
	_ = v162
	var v163 int32
	_ = v163
	var v165 int64
	_ = v165
	var v171 int32
	_ = v171
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v180 int32
	_ = v180
	var v182 int32
	_ = v182
	var v186 int32
	_ = v186
	var v187 int32
	_ = v187
	var v192 int32
	_ = v192
	var v193 int32
	_ = v193
	var v194 int32
	_ = v194
	var v201 int32
	_ = v201
	var v202 int32
	_ = v202
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 int64
	_ = v205
	var v207 int32
	_ = v207
	var v209 int32
	_ = v209
	var v217 int32
	_ = v217
	var v219 int32
	_ = v219
	var v232 int32
	_ = v232
	var v235 int32
	_ = v235
	var v237 int32
	_ = v237
	var v247 int32
	_ = v247
	var v269 int32
	_ = v269
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v274 int32
	_ = v274
	var v275 int32
	_ = v275
	var v276 int32
	_ = v276
	var v277 int32
	_ = v277
	var v278 int32
	_ = v278
	var v279 int32
	_ = v279
	var v280 int32
	_ = v280
	var v281 int32
	_ = v281
	var v283 int32
	_ = v283
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
	var v299 int32
	_ = v299
	var v300 int32
	_ = v300
	var v301 int32
	_ = v301
	var v302 int32
	_ = v302
	var v303 int32
	_ = v303
	var v304 int32
	_ = v304
	var v305 int32
	_ = v305
	var v306 int32
	_ = v306
	var v307 int32
	_ = v307
	var v342 int32
	_ = v342
	var v343 int32
	_ = v343
	var v344 int32
	_ = v344
	var v347 int32
	_ = v347
	var v349 int32
	_ = v349
	var v357 int32
	_ = v357
	var v366 int32
	_ = v366
	var v369 int32
	_ = v369
	var v371 int32
	_ = v371
	var v373 int32
	_ = v373
	var v377 int32
	_ = v377
	var v380 int32
	_ = v380
	var v382 int32
	_ = v382
	var v384 int32
	_ = v384
	var v395 int32
	_ = v395
	var v396 int32
	_ = v396
	var v404 int32
	_ = v404
	var v407 int32
	_ = v407
	var v409 int32
	_ = v409
	var v411 int32
	_ = v411
	var v422 int32
	_ = v422
	var v423 int32
	_ = v423
	var v430 int32
	_ = v430
	var v433 int32
	_ = v433
	var v438 int32
	_ = v438
	var v449 int64
	_ = v449
	var v450 int32
	_ = v450
	var v454 int32
	_ = v454
	var v457 int32
	_ = v457
	var v461 int32
	_ = v461
	var v463 int32
	_ = v463
	var v464 int64
	_ = v464
	var v466 int32
	_ = v466
	var v469 int32
	_ = v469
	var v470 int64
	_ = v470
	var v472 int32
	_ = v472
	var v474 int32
	_ = v474
	var v483 int32
	_ = v483
	var v484 int64
	_ = v484
	var v485 int64
	_ = v485
	var v489 int64
	_ = v489
	var v490 int32
	_ = v490
	var v493 int64
	_ = v493
	var v499 int32
	_ = v499
	var v500 int32
	_ = v500
	var v502 int32
	_ = v502
	var v507 int32
	_ = v507
	var v511 int32
	_ = v511
	var v514 int32
	_ = v514
	var v516 int32
	_ = v516
	var v517 int32
	_ = v517
	var v529 int32
	_ = v529
	var v532 int32
	_ = v532
	var v533 int32
	_ = v533
	var v536 int32
	_ = v536
	var v537 int32
	_ = v537
	var v546 int32
	_ = v546
	var v547 int32
	_ = v547
	var v548 int32
	_ = v548
	var v549 int32
	_ = v549
	var v554 int32
	_ = v554
	var v556 int32
	_ = v556
	var v565 int32
	_ = v565
	var v568 int32
	_ = v568
	var v580 int32
	_ = v580
	var v582 int32
	_ = v582
	var v587 int32
	_ = v587
	var v588 int32
	_ = v588
	var v590 int32
	_ = v590
	var v593 int32
	_ = v593
	var v595 int32
	_ = v595
	var v596 int32
	_ = v596
	var v603 int32
	_ = v603
	var v606 int32
	_ = v606
	var v611 int32
	_ = v611
	var v621 int32
	_ = v621
	var v624 int32
	_ = v624
	var v625 int32
	_ = v625
	var v630 int32
	_ = v630
	var v631 int32
	_ = v631
	var v636 int32
	_ = v636
	var v639 int32
	_ = v639
	var v651 int32
	_ = v651
	var v674 int32
	_ = v674
	var v676 int32
	_ = v676
	var v679 int32
	_ = v679
	var v681 int32
	_ = v681
	var v685 int32
	_ = v685
	var v686 int32
	_ = v686
	var v687 int32
	_ = v687
	var v688 int32
	_ = v688
	var v689 int32
	_ = v689
	var v691 int32
	_ = v691
	var v694 int32
	_ = v694
	var v696 int32
	_ = v696
	var v700 int32
	_ = v700
	var v701 int32
	_ = v701
	var v714 int32
	_ = v714
	var v715 int32
	_ = v715
	var v719 int32
	_ = v719
	var v720 int32
	_ = v720
	var v728 int32
	_ = v728
	var v729 int32
	_ = v729
	var v730 int32
	_ = v730
	var v731 int32
	_ = v731
	var v733 int32
	_ = v733
	var v734 int32
	_ = v734
	var v738 int32
	_ = v738
	var v739 int32
	_ = v739
	var v740 int32
	_ = v740
	var v741 int32
	_ = v741
	var v750 int32
	_ = v750
	var v751 int32
	_ = v751
	var v754 int64
	_ = v754
	var v757 int32
	_ = v757
	var v759 int64
	_ = v759
	var v762 int32
	_ = v762
	var v769 int32
	_ = v769
	var v771 int32
	_ = v771
	var v775 int32
	_ = v775
	var v776 int32
	_ = v776
	var v777 int32
	_ = v777
	var v778 int32
	_ = v778
	var v787 int32
	_ = v787
	var v788 int32
	_ = v788
	var v789 int32
	_ = v789
	var v791 int32
	_ = v791
	var v796 int32
	_ = v796
	var v797 int32
	_ = v797
	var v804 int32
	_ = v804
	var v807 int32
	_ = v807
	var v809 int32
	_ = v809
	var v813 int32
	_ = v813
	var v814 int32
	_ = v814
	var v815 int32
	_ = v815
	var v816 int32
	_ = v816
	var v853 int32
	_ = v853
	var v854 int32
	_ = v854
	var v866 int32
	_ = v866
	var v867 int32
	_ = v867
	var v871 int32
	_ = v871
	var v872 int32
	_ = v872
	var v880 int32
	_ = v880
	var v881 int32
	_ = v881
	var v882 int32
	_ = v882
	var v883 int32
	_ = v883
	var v885 int32
	_ = v885
	var v886 int32
	_ = v886
	var v890 int32
	_ = v890
	var v891 int32
	_ = v891
	var v892 int32
	_ = v892
	var v893 int32
	_ = v893
	var v902 int32
	_ = v902
	var v903 int32
	_ = v903
	var v906 int64
	_ = v906
	var v909 int32
	_ = v909
	var v911 int64
	_ = v911
	var v914 int32
	_ = v914
	var v921 int32
	_ = v921
	var v923 int32
	_ = v923
	var v927 int32
	_ = v927
	var v928 int32
	_ = v928
	var v929 int32
	_ = v929
	var v930 int32
	_ = v930
	var v939 int32
	_ = v939
	var v940 int32
	_ = v940
	var v941 int32
	_ = v941
	var v943 int32
	_ = v943
	var v948 int32
	_ = v948
	var v949 int32
	_ = v949
	var v956 int32
	_ = v956
	var v959 int32
	_ = v959
	var v961 int32
	_ = v961
	var v965 int32
	_ = v965
	var v966 int32
	_ = v966
	var v967 int32
	_ = v967
	var v968 int32
	_ = v968
	var v1005 int32
	_ = v1005
	var v1008 int32
	_ = v1008
	var v1010 int32
	_ = v1010
	var v1011 int32
	_ = v1011
	var v1023 int32
	_ = v1023
	var v1024 int32
	_ = v1024
	var v1028 int32
	_ = v1028
	var v1029 int32
	_ = v1029
	var v1037 int32
	_ = v1037
	var v1038 int32
	_ = v1038
	var v1039 int32
	_ = v1039
	var v1040 int32
	_ = v1040
	var v1042 int32
	_ = v1042
	var v1043 int32
	_ = v1043
	var v1047 int32
	_ = v1047
	var v1048 int32
	_ = v1048
	var v1049 int32
	_ = v1049
	var v1050 int32
	_ = v1050
	var v1059 int32
	_ = v1059
	var v1060 int32
	_ = v1060
	var v1063 int64
	_ = v1063
	var v1066 int32
	_ = v1066
	var v1068 int64
	_ = v1068
	var v1071 int32
	_ = v1071
	var v1078 int32
	_ = v1078
	var v1080 int32
	_ = v1080
	var v1084 int32
	_ = v1084
	var v1085 int32
	_ = v1085
	var v1086 int32
	_ = v1086
	var v1087 int32
	_ = v1087
	var v1096 int32
	_ = v1096
	var v1097 int32
	_ = v1097
	var v1098 int32
	_ = v1098
	var v1100 int32
	_ = v1100
	var v1105 int32
	_ = v1105
	var v1106 int32
	_ = v1106
	var v1113 int32
	_ = v1113
	var v1116 int32
	_ = v1116
	var v1118 int32
	_ = v1118
	var v1122 int32
	_ = v1122
	var v1123 int32
	_ = v1123
	var v1124 int32
	_ = v1124
	var v1125 int32
	_ = v1125
	var v1167 int32
	_ = v1167
	var v1179 int32
	_ = v1179
	var v1184 int32
	_ = v1184
	var v1195 int32
	_ = v1195
	var v1196 int32
	_ = v1196
	var v1253 int32
	_ = v1253
	var v1254 base.V128
	_ = v1254
	var v1258 base.V128
	_ = v1258
	var v1259 base.V128
	_ = v1259
	var v1265 base.V128
	_ = v1265
	var v1269 base.V128
	_ = v1269
	var v1272 int32
	_ = v1272
	var v1277 int32
	_ = v1277
	var v1321 int32
	_ = v1321
	var v1353 int32
	_ = v1353
	var v1362 int32
	_ = v1362
	var v1373 int32
	_ = v1373
	var v1374 int32
	_ = v1374
	var v1431 int32
	_ = v1431
	var v1434 int32
	_ = v1434
	var v1435 int32
	_ = v1435
	var v1442 int32
	_ = v1442
	var v1449 int32
	_ = v1449
	var v1454 int32
	_ = v1454
	var v1528 int32
	_ = v1528
	var v1534 int32
	_ = v1534
	var v1537 int32
	_ = v1537
	var v1540 int32
	_ = v1540
	var v1541 int32
	_ = v1541
	var v1543 int32
	_ = v1543
	var v1545 int32
	_ = v1545
	var v1548 int32
	_ = v1548
	var v1553 int32
	_ = v1553
	var v1558 int32
	_ = v1558
	var v1559 int32
	_ = v1559
	var v1566 int64
	_ = v1566
	var v1567 int32
	_ = v1567
	var v1571 int32
	_ = v1571
	var v1574 int32
	_ = v1574
	var v1578 int32
	_ = v1578
	var v1580 int32
	_ = v1580
	var v1581 int64
	_ = v1581
	var v1583 int32
	_ = v1583
	var v1586 int32
	_ = v1586
	var v1587 int64
	_ = v1587
	var v1589 int32
	_ = v1589
	var v1591 int32
	_ = v1591
	var v1600 int32
	_ = v1600
	var v1601 int64
	_ = v1601
	var v1602 int64
	_ = v1602
	var v1606 int64
	_ = v1606
	var v1607 int32
	_ = v1607
	var v1610 int64
	_ = v1610
	var v1616 int32
	_ = v1616
	var v1617 int32
	_ = v1617
	var v1619 int32
	_ = v1619
	var v1624 int32
	_ = v1624
	var v1628 int32
	_ = v1628
	var v1631 int32
	_ = v1631
	var v1633 int32
	_ = v1633
	var v1634 int32
	_ = v1634
	var v1646 int32
	_ = v1646
	var v1649 int32
	_ = v1649
	var v1650 int32
	_ = v1650
	var v1651 int32
	_ = v1651
	var v1652 int32
	_ = v1652
	var v1653 int64
	_ = v1653
	var v1654 int32
	_ = v1654
	var v1661 int64
	_ = v1661
	var v1662 int32
	_ = v1662
	var v1663 int64
	_ = v1663
	var v1673 int32
	_ = v1673
	var v1675 int32
	_ = v1675
	var v1685 int32
	_ = v1685
	var v1686 int32
	_ = v1686
	var v1728 int32
	_ = v1728
	var v1730 int32
	_ = v1730
	var v1739 int32
	_ = v1739
	var v1740 int32
	_ = v1740
	var v1742 int32
	_ = v1742
	var v1765 int32
	_ = v1765
	var v1776 int32
	_ = v1776
	var v1801 int32
	_ = v1801
	var v1802 int32
	_ = v1802
	var v1803 int32
	_ = v1803
	var v1832 int32
	_ = v1832
	var v1834 int32
	_ = v1834
	var v1837 int32
	_ = v1837
	var v1842 int32
	_ = v1842
	var v1844 int32
	_ = v1844
	var v1845 int32
	_ = v1845
	var v1846 int32
	_ = v1846
	var v1847 int32
	_ = v1847
	var v1854 int32
	_ = v1854
	var v1895 int32
	_ = v1895
	var v1897 int32
	_ = v1897
	var v1926 int32
	_ = v1926
	var v1927 int32
	_ = v1927
	var v1934 int32
	_ = v1934
	var v2065 int32
	_ = v2065
	var v2076 int32
	_ = v2076
	var v2134 int32
	_ = v2134
	var v2137 int32
	_ = v2137
	var v2138 int32
	_ = v2138
	var v2270 int32
	_ = v2270
	var v2281 int32
	_ = v2281
	var v2339 int32
	_ = v2339
	var v2348 int32
	_ = v2348
	var v2349 int32
	_ = v2349
	var v2481 int32
	_ = v2481
	var v2492 int32
	_ = v2492
	var v2550 int32
	_ = v2550
	var v2559 int32
	_ = v2559
	var v2560 int32
	_ = v2560
	var v2567 int32
	_ = v2567
	var v2573 int32
	_ = v2573
	var v2587 int32
	_ = v2587
	var v2634 int32
	_ = v2634
	var v2654 int32
	_ = v2654
	var v2655 int32
	_ = v2655
	var v2658 int32
	_ = v2658
	var v2660 int32
	_ = v2660
	var v2662 int32
	_ = v2662
	var v2714 int32
	_ = v2714
	var v2736 int32
	_ = v2736
	var v2737 int32
	_ = v2737
	var v2744 int32
	_ = v2744
	var v2750 int32
	_ = v2750
	var v2764 int32
	_ = v2764
	var v2811 int32
	_ = v2811
	var v2831 int32
	_ = v2831
	var v2832 int32
	_ = v2832
	var v2835 int32
	_ = v2835
	var v2837 int32
	_ = v2837
	var v2839 int32
	_ = v2839
	var v2891 int32
	_ = v2891
	var v2913 int32
	_ = v2913
	var v2914 int32
	_ = v2914
	var v2923 int32
	_ = v2923
	var v2924 int32
	_ = v2924
	var v2942 int32
	_ = v2942
	var v2967 int32
	_ = v2967
	var v2977 int32
	_ = v2977
	var v2978 int32
	_ = v2978
	var v2979 int32
	_ = v2979
	var v3005 int32
	_ = v3005
	var v3016 int32
	_ = v3016
	var v3017 int32
	_ = v3017
	var v3042 int32
	_ = v3042
	var v3072 int32
	_ = v3072
	var v3085 int32
	_ = v3085
	var v3086 int32
	_ = v3086
	var v3087 int32
	_ = v3087
	var v3088 int32
	_ = v3088
	var v3092 int32
	_ = v3092
	var v3105 int32
	_ = v3105
	var v3106 int32
	_ = v3106
	var v3107 int32
	_ = v3107
	var v3114 int32
	_ = v3114
	var v3131 int32
	_ = v3131
	var v3132 int32
	_ = v3132
	var v3157 int32
	_ = v3157
	var v3192 int32
	_ = v3192
	var v3205 int32
	_ = v3205
	var v3206 int32
	_ = v3206
	var v3207 int32
	_ = v3207
	var v3211 int32
	_ = v3211
	var v3212 int32
	_ = v3212
	var v3213 int32
	_ = v3213
	var v3214 int32
	_ = v3214
	var v3221 int32
	_ = v3221
	var v3229 int32
	_ = v3229
	var v3230 int32
	_ = v3230
	var v3248 int32
	_ = v3248
	var v3273 int32
	_ = v3273
	var v3283 int32
	_ = v3283
	var v3284 int32
	_ = v3284
	var v3285 int32
	_ = v3285
	var v3311 int32
	_ = v3311
	var v3322 int32
	_ = v3322
	var v3323 int32
	_ = v3323
	var v3348 int32
	_ = v3348
	var v3378 int32
	_ = v3378
	var v3391 int32
	_ = v3391
	var v3392 int32
	_ = v3392
	var v3393 int32
	_ = v3393
	var v3394 int32
	_ = v3394
	var v3398 int32
	_ = v3398
	var v3411 int32
	_ = v3411
	var v3412 int32
	_ = v3412
	var v3413 int32
	_ = v3413
	var v3420 int32
	_ = v3420
	var v3437 int32
	_ = v3437
	var v3438 int32
	_ = v3438
	var v3463 int32
	_ = v3463
	var v3498 int32
	_ = v3498
	var v3511 int32
	_ = v3511
	var v3512 int32
	_ = v3512
	var v3513 int32
	_ = v3513
	var v3517 int32
	_ = v3517
	var v3518 int32
	_ = v3518
	var v3519 int32
	_ = v3519
	var v3520 int32
	_ = v3520
	var v3527 int32
	_ = v3527
	var v3542 int32
	_ = v3542
	var v3544 int32
	_ = v3544
	var v3556 int32
	_ = v3556
	var v3563 int32
	_ = v3563
	var v3567 int32
	_ = v3567
	var v3577 int32
	_ = v3577
	var v3578 int32
	_ = v3578
	var v3583 int32
	_ = v3583
	var v3587 int32
	_ = v3587
	var v3588 int32
	_ = v3588
	var v3589 int32
	_ = v3589
	var v3590 int32
	_ = v3590
	var v3594 int32
	_ = v3594
	var v3608 int32
	_ = v3608
	var v3625 int32
	_ = v3625
	var v3639 int32
	_ = v3639
	var v3686 int32
	_ = v3686
	var v3706 int32
	_ = v3706
	var v3707 int32
	_ = v3707
	var v3710 int32
	_ = v3710
	var v3712 int32
	_ = v3712
	var v3714 int32
	_ = v3714
	var v3766 int32
	_ = v3766
	var v3788 int32
	_ = v3788
	var v3789 int32
	_ = v3789
	var v3794 int32
	_ = v3794
	var v3795 int32
	_ = v3795
	var v3838 int32
	_ = v3838
	var v3839 int32
	_ = v3839
	var v3848 int32
	_ = v3848
	var v3849 int32
	_ = v3849
	var v3851 int32
	_ = v3851
	var v3873 int32
	_ = v3873
	var v3884 int32
	_ = v3884
	var v3910 int32
	_ = v3910
	var v3911 int32
	_ = v3911
	var v3940 int32
	_ = v3940
	var v3945 int32
	_ = v3945
	var v3951 int32
	_ = v3951
	var v3956 int32
	_ = v3956
	var v3958 int32
	_ = v3958
	var v3963 int32
	_ = v3963
	var v3967 int32
	_ = v3967
	var v3968 int32
	_ = v3968
	var v3969 int32
	_ = v3969
	var v3970 int32
	_ = v3970
	var v3975 int32
	_ = v3975
	var v3982 int32
	_ = v3982
	var v3983 int32
	_ = v3983
	var v3984 int32
	_ = v3984
	var v3990 int32
	_ = v3990
	var v3992 int32
	_ = v3992
	var v3993 int32
	_ = v3993
	var v4000 int32
	_ = v4000
	var v4006 int32
	_ = v4006
	var v4012 int32
	_ = v4012
	var v4026 int32
	_ = v4026
	var v4073 int32
	_ = v4073
	var v4093 int32
	_ = v4093
	var v4094 int32
	_ = v4094
	var v4097 int32
	_ = v4097
	var v4099 int32
	_ = v4099
	var v4101 int32
	_ = v4101
	var v4153 int32
	_ = v4153
	var v4175 int32
	_ = v4175
	var v4176 int32
	_ = v4176
	var v4185 int32
	_ = v4185
	var v4186 int32
	_ = v4186
	var v4204 int32
	_ = v4204
	var v4229 int32
	_ = v4229
	var v4239 int32
	_ = v4239
	var v4240 int32
	_ = v4240
	var v4241 int32
	_ = v4241
	var v4267 int32
	_ = v4267
	var v4278 int32
	_ = v4278
	var v4279 int32
	_ = v4279
	var v4304 int32
	_ = v4304
	var v4334 int32
	_ = v4334
	var v4343 int32
	_ = v4343
	var v4344 int32
	_ = v4344
	var v4345 int32
	_ = v4345
	var v4346 int32
	_ = v4346
	var v4350 int32
	_ = v4350
	var v4359 int32
	_ = v4359
	var v4360 int32
	_ = v4360
	var v4361 int32
	_ = v4361
	var v4368 int32
	_ = v4368
	var v4385 int32
	_ = v4385
	var v4386 int32
	_ = v4386
	var v4411 int32
	_ = v4411
	var v4446 int32
	_ = v4446
	var v4455 int32
	_ = v4455
	var v4456 int32
	_ = v4456
	var v4457 int32
	_ = v4457
	var v4461 int32
	_ = v4461
	var v4462 int32
	_ = v4462
	var v4463 int32
	_ = v4463
	var v4464 int32
	_ = v4464
	var v4471 int32
	_ = v4471
	var v4549 int32
	_ = v4549
	var v4550 int32
	_ = v4550
	var v4551 int32
	_ = v4551
	var v4558 int64
	_ = v4558
	var v4559 int32
	_ = v4559
	var v4563 int32
	_ = v4563
	var v4566 int32
	_ = v4566
	var v4570 int32
	_ = v4570
	var v4572 int32
	_ = v4572
	var v4573 int64
	_ = v4573
	var v4575 int32
	_ = v4575
	var v4578 int32
	_ = v4578
	var v4579 int64
	_ = v4579
	var v4581 int32
	_ = v4581
	var v4583 int32
	_ = v4583
	var v4592 int32
	_ = v4592
	var v4593 int64
	_ = v4593
	var v4594 int64
	_ = v4594
	var v4598 int64
	_ = v4598
	var v4599 int32
	_ = v4599
	var v4602 int64
	_ = v4602
	var v4608 int32
	_ = v4608
	var v4609 int32
	_ = v4609
	var v4611 int32
	_ = v4611
	var v4616 int32
	_ = v4616
	var v4620 int32
	_ = v4620
	var v4623 int32
	_ = v4623
	var v4625 int32
	_ = v4625
	var v4626 int32
	_ = v4626
	var v4638 int32
	_ = v4638
	var v4641 int32
	_ = v4641
	var v4646 int32
	_ = v4646
	var v4649 int32
	_ = v4649
	var v4650 int32
	_ = v4650
	var v4651 int32
	_ = v4651
	var v4662 int32
	_ = v4662
	var v4699 int32
	_ = v4699
	var v4700 int32
	_ = v4700
	var v4729 int32
	_ = v4729
	var v4731 int32
	_ = v4731
	var v4732 int32
	_ = v4732
	var v4734 int32
	_ = v4734
	var v4735 int32
	_ = v4735
	var v4736 int32
	_ = v4736
	var v4739 int32
	_ = v4739
	var v4740 int32
	_ = v4740
	var v4742 int32
	_ = v4742
	var v4747 int32
	_ = v4747
	var v4784 int32
	_ = v4784
	var v4816 int32
	_ = v4816
	var v4891 int32
	_ = v4891
	var v4969 int32
	_ = v4969
	var v5021 int32
	_ = v5021
	var v5049 int32
	_ = v5049
	var v5052 int32
	_ = v5052
	var v5053 int32
	_ = v5053
	var v5066 int32
	_ = v5066
	var v5067 int32
	_ = v5067
	var v5071 int32
	_ = v5071
	var v5072 int32
	_ = v5072
	var v5080 int32
	_ = v5080
	var v5081 int32
	_ = v5081
	var v5082 int32
	_ = v5082
	var v5083 int32
	_ = v5083
	var v5085 int32
	_ = v5085
	var v5086 int32
	_ = v5086
	var v5090 int32
	_ = v5090
	var v5091 int32
	_ = v5091
	var v5092 int32
	_ = v5092
	var v5093 int32
	_ = v5093
	var v5102 int32
	_ = v5102
	var v5103 int32
	_ = v5103
	var v5106 int64
	_ = v5106
	var v5109 int32
	_ = v5109
	var v5111 int64
	_ = v5111
	var v5114 int32
	_ = v5114
	var v5121 int32
	_ = v5121
	var v5123 int32
	_ = v5123
	var v5127 int32
	_ = v5127
	var v5128 int32
	_ = v5128
	var v5129 int32
	_ = v5129
	var v5130 int32
	_ = v5130
	var v5139 int32
	_ = v5139
	var v5140 int32
	_ = v5140
	var v5141 int32
	_ = v5141
	var v5143 int32
	_ = v5143
	var v5148 int32
	_ = v5148
	var v5149 int32
	_ = v5149
	var v5156 int32
	_ = v5156
	var v5159 int32
	_ = v5159
	var v5161 int32
	_ = v5161
	var v5165 int32
	_ = v5165
	var v5166 int32
	_ = v5166
	var v5167 int32
	_ = v5167
	var v5168 int32
	_ = v5168
	var v5205 int32
	_ = v5205
	var v5218 int32
	_ = v5218
	var v5219 int32
	_ = v5219
	var v5223 int32
	_ = v5223
	var v5224 int32
	_ = v5224
	var v5232 int32
	_ = v5232
	var v5233 int32
	_ = v5233
	var v5234 int32
	_ = v5234
	var v5235 int32
	_ = v5235
	var v5237 int32
	_ = v5237
	var v5238 int32
	_ = v5238
	var v5242 int32
	_ = v5242
	var v5243 int32
	_ = v5243
	var v5244 int32
	_ = v5244
	var v5245 int32
	_ = v5245
	var v5254 int32
	_ = v5254
	var v5255 int32
	_ = v5255
	var v5258 int64
	_ = v5258
	var v5261 int32
	_ = v5261
	var v5263 int64
	_ = v5263
	var v5266 int32
	_ = v5266
	var v5273 int32
	_ = v5273
	var v5275 int32
	_ = v5275
	var v5279 int32
	_ = v5279
	var v5280 int32
	_ = v5280
	var v5281 int32
	_ = v5281
	var v5282 int32
	_ = v5282
	var v5291 int32
	_ = v5291
	var v5292 int32
	_ = v5292
	var v5293 int32
	_ = v5293
	var v5295 int32
	_ = v5295
	var v5300 int32
	_ = v5300
	var v5301 int32
	_ = v5301
	var v5308 int32
	_ = v5308
	var v5311 int32
	_ = v5311
	var v5313 int32
	_ = v5313
	var v5317 int32
	_ = v5317
	var v5318 int32
	_ = v5318
	var v5319 int32
	_ = v5319
	var v5320 int32
	_ = v5320
	var v5357 int32
	_ = v5357
	var v5359 int32
	_ = v5359
	var v5360 int32
	_ = v5360
	var v5363 int32
	_ = v5363
	var v5367 int32
	_ = v5367
	var v5369 int32
	_ = v5369
	var v5370 int32
	_ = v5370
	var v5371 int32
	_ = v5371
	var v5372 int32
	_ = v5372
	var v5373 int32
	_ = v5373
	var v5376 int32
	_ = v5376
	var v5377 int32
	_ = v5377
	var v5380 int32
	_ = v5380
	var v5383 int32
	_ = v5383
	var v5384 int32
	_ = v5384
	var v5385 int32
	_ = v5385
	var v5387 int32
	_ = v5387
	var v5392 int32
	_ = v5392
	var v5412 int32
	_ = v5412
	var v5469 int32
	_ = v5469
	var v5472 int32
	_ = v5472
	var v5476 int32
	_ = v5476
	var v5485 int32
	_ = v5485
	var v5522 int32
	_ = v5522
	var v5556 int32
	_ = v5556
	var v5558 int32
	_ = v5558
	var v5559 int32
	_ = v5559
	var v5571 int32
	_ = v5571
	var v5639 int32
	_ = v5639
	var v5644 int32
	_ = v5644
	var v5645 int32
	_ = v5645
	var v5646 int32
	_ = v5646
	var v5696 int32
	_ = v5696
	var v5725 int32
	_ = v5725
	var v5727 int32
	_ = v5727
	var v5728 int32
	_ = v5728
	var v5736 int32
	_ = v5736
	var v5737 int32
	_ = v5737
	var v5738 int32
	_ = v5738
	var v5741 int32
	_ = v5741
	var v5744 int32
	_ = v5744
	var v5745 int32
	_ = v5745
	var v5750 int32
	_ = v5750
	var v5770 int32
	_ = v5770
	var v5827 int32
	_ = v5827
	var v5830 int32
	_ = v5830
	var v5834 int32
	_ = v5834
	var v5843 int32
	_ = v5843
	var v5880 int32
	_ = v5880
	var v5914 int32
	_ = v5914
	var v5916 int32
	_ = v5916
	var v5917 int32
	_ = v5917
	var v5929 int32
	_ = v5929
	var v5997 int32
	_ = v5997
	var v6002 int32
	_ = v6002
	var v6003 int32
	_ = v6003
	var v6004 int32
	_ = v6004
	var v6053 int32
	_ = v6053
	var v6083 int32
	_ = v6083
	var v6084 int32
	_ = v6084
	var v6085 int32
	_ = v6085
	var v6086 int32
	_ = v6086
	var v6087 int32
	_ = v6087
	var v6089 int32
	_ = v6089
	var v6094 int32
	_ = v6094
	var v6097 int32
	_ = v6097
	var v6110 int32
	_ = v6110
	var v6111 int32
	_ = v6111
	var v6115 int32
	_ = v6115
	var v6116 int32
	_ = v6116
	var v6124 int32
	_ = v6124
	var v6125 int32
	_ = v6125
	var v6126 int32
	_ = v6126
	var v6127 int32
	_ = v6127
	var v6129 int32
	_ = v6129
	var v6130 int32
	_ = v6130
	var v6134 int32
	_ = v6134
	var v6135 int32
	_ = v6135
	var v6136 int32
	_ = v6136
	var v6137 int32
	_ = v6137
	var v6146 int32
	_ = v6146
	var v6147 int32
	_ = v6147
	var v6150 int64
	_ = v6150
	var v6153 int32
	_ = v6153
	var v6155 int64
	_ = v6155
	var v6158 int32
	_ = v6158
	var v6165 int32
	_ = v6165
	var v6167 int32
	_ = v6167
	var v6171 int32
	_ = v6171
	var v6172 int32
	_ = v6172
	var v6173 int32
	_ = v6173
	var v6174 int32
	_ = v6174
	var v6183 int32
	_ = v6183
	var v6184 int32
	_ = v6184
	var v6185 int32
	_ = v6185
	var v6187 int32
	_ = v6187
	var v6192 int32
	_ = v6192
	var v6193 int32
	_ = v6193
	var v6200 int32
	_ = v6200
	var v6203 int32
	_ = v6203
	var v6205 int32
	_ = v6205
	var v6209 int32
	_ = v6209
	var v6210 int32
	_ = v6210
	var v6211 int32
	_ = v6211
	var v6212 int32
	_ = v6212
	var v6249 int32
	_ = v6249
	var v6250 int32
	_ = v6250
	var v6262 int32
	_ = v6262
	var v6263 int32
	_ = v6263
	var v6267 int32
	_ = v6267
	var v6268 int32
	_ = v6268
	var v6276 int32
	_ = v6276
	var v6277 int32
	_ = v6277
	var v6278 int32
	_ = v6278
	var v6279 int32
	_ = v6279
	var v6281 int32
	_ = v6281
	var v6282 int32
	_ = v6282
	var v6286 int32
	_ = v6286
	var v6287 int32
	_ = v6287
	var v6288 int32
	_ = v6288
	var v6289 int32
	_ = v6289
	var v6298 int32
	_ = v6298
	var v6299 int32
	_ = v6299
	var v6302 int64
	_ = v6302
	var v6305 int32
	_ = v6305
	var v6307 int64
	_ = v6307
	var v6310 int32
	_ = v6310
	var v6317 int32
	_ = v6317
	var v6319 int32
	_ = v6319
	var v6323 int32
	_ = v6323
	var v6324 int32
	_ = v6324
	var v6325 int32
	_ = v6325
	var v6326 int32
	_ = v6326
	var v6335 int32
	_ = v6335
	var v6336 int32
	_ = v6336
	var v6337 int32
	_ = v6337
	var v6339 int32
	_ = v6339
	var v6344 int32
	_ = v6344
	var v6345 int32
	_ = v6345
	var v6352 int32
	_ = v6352
	var v6355 int32
	_ = v6355
	var v6357 int32
	_ = v6357
	var v6361 int32
	_ = v6361
	var v6362 int32
	_ = v6362
	var v6363 int32
	_ = v6363
	var v6364 int32
	_ = v6364
	var v6401 int32
	_ = v6401
	var v6403 int32
	_ = v6403
	var v6404 int32
	_ = v6404
	var v6416 int32
	_ = v6416
	var v6417 int32
	_ = v6417
	var v6421 int32
	_ = v6421
	var v6422 int32
	_ = v6422
	var v6430 int32
	_ = v6430
	var v6431 int32
	_ = v6431
	var v6432 int32
	_ = v6432
	var v6433 int32
	_ = v6433
	var v6435 int32
	_ = v6435
	var v6436 int32
	_ = v6436
	var v6440 int32
	_ = v6440
	var v6441 int32
	_ = v6441
	var v6442 int32
	_ = v6442
	var v6443 int32
	_ = v6443
	var v6452 int32
	_ = v6452
	var v6453 int32
	_ = v6453
	var v6456 int64
	_ = v6456
	var v6459 int32
	_ = v6459
	var v6461 int64
	_ = v6461
	var v6464 int32
	_ = v6464
	var v6471 int32
	_ = v6471
	var v6473 int32
	_ = v6473
	var v6477 int32
	_ = v6477
	var v6478 int32
	_ = v6478
	var v6479 int32
	_ = v6479
	var v6480 int32
	_ = v6480
	var v6489 int32
	_ = v6489
	var v6490 int32
	_ = v6490
	var v6491 int32
	_ = v6491
	var v6493 int32
	_ = v6493
	var v6498 int32
	_ = v6498
	var v6499 int32
	_ = v6499
	var v6506 int32
	_ = v6506
	var v6509 int32
	_ = v6509
	var v6511 int32
	_ = v6511
	var v6515 int32
	_ = v6515
	var v6516 int32
	_ = v6516
	var v6517 int32
	_ = v6517
	var v6518 int32
	_ = v6518
	var v6555 int32
	_ = v6555
	var v6557 int32
	_ = v6557
	var v6558 int32
	_ = v6558
	var v6563 int32
	_ = v6563
	var v6567 int32
	_ = v6567
	var v6614 int32
	_ = v6614
	var v6642 int32
	_ = v6642
	var v6646 int32
	_ = v6646
	var v6647 int32
	_ = v6647
	var v6648 int32
	_ = v6648
	var v6649 int32
	_ = v6649
	var v6650 int32
	_ = v6650
	var v6651 int32
	_ = v6651
	var v6653 int32
	_ = v6653
	var v6656 int32
	_ = v6656
	var v6659 int32
	_ = v6659
	var v6672 int32
	_ = v6672
	var v6673 int32
	_ = v6673
	var v6677 int32
	_ = v6677
	var v6678 int32
	_ = v6678
	var v6686 int32
	_ = v6686
	var v6687 int32
	_ = v6687
	var v6688 int32
	_ = v6688
	var v6689 int32
	_ = v6689
	var v6691 int32
	_ = v6691
	var v6692 int32
	_ = v6692
	var v6696 int32
	_ = v6696
	var v6697 int32
	_ = v6697
	var v6698 int32
	_ = v6698
	var v6699 int32
	_ = v6699
	var v6708 int32
	_ = v6708
	var v6709 int32
	_ = v6709
	var v6712 int64
	_ = v6712
	var v6715 int32
	_ = v6715
	var v6717 int64
	_ = v6717
	var v6720 int32
	_ = v6720
	var v6727 int32
	_ = v6727
	var v6729 int32
	_ = v6729
	var v6733 int32
	_ = v6733
	var v6734 int32
	_ = v6734
	var v6735 int32
	_ = v6735
	var v6736 int32
	_ = v6736
	var v6745 int32
	_ = v6745
	var v6746 int32
	_ = v6746
	var v6747 int32
	_ = v6747
	var v6749 int32
	_ = v6749
	var v6754 int32
	_ = v6754
	var v6755 int32
	_ = v6755
	var v6762 int32
	_ = v6762
	var v6765 int32
	_ = v6765
	var v6767 int32
	_ = v6767
	var v6771 int32
	_ = v6771
	var v6772 int32
	_ = v6772
	var v6773 int32
	_ = v6773
	var v6774 int32
	_ = v6774
	var v6811 int32
	_ = v6811
	var v6812 int32
	_ = v6812
	var v6824 int32
	_ = v6824
	var v6825 int32
	_ = v6825
	var v6829 int32
	_ = v6829
	var v6830 int32
	_ = v6830
	var v6838 int32
	_ = v6838
	var v6839 int32
	_ = v6839
	var v6840 int32
	_ = v6840
	var v6841 int32
	_ = v6841
	var v6843 int32
	_ = v6843
	var v6844 int32
	_ = v6844
	var v6848 int32
	_ = v6848
	var v6849 int32
	_ = v6849
	var v6850 int32
	_ = v6850
	var v6851 int32
	_ = v6851
	var v6860 int32
	_ = v6860
	var v6861 int32
	_ = v6861
	var v6864 int64
	_ = v6864
	var v6867 int32
	_ = v6867
	var v6869 int64
	_ = v6869
	var v6872 int32
	_ = v6872
	var v6879 int32
	_ = v6879
	var v6881 int32
	_ = v6881
	var v6885 int32
	_ = v6885
	var v6886 int32
	_ = v6886
	var v6887 int32
	_ = v6887
	var v6888 int32
	_ = v6888
	var v6897 int32
	_ = v6897
	var v6898 int32
	_ = v6898
	var v6899 int32
	_ = v6899
	var v6901 int32
	_ = v6901
	var v6906 int32
	_ = v6906
	var v6907 int32
	_ = v6907
	var v6914 int32
	_ = v6914
	var v6917 int32
	_ = v6917
	var v6919 int32
	_ = v6919
	var v6923 int32
	_ = v6923
	var v6924 int32
	_ = v6924
	var v6925 int32
	_ = v6925
	var v6926 int32
	_ = v6926
	var v6963 int32
	_ = v6963
	var v6965 int32
	_ = v6965
	var v6966 int32
	_ = v6966
	var v6978 int32
	_ = v6978
	var v6979 int32
	_ = v6979
	var v6983 int32
	_ = v6983
	var v6984 int32
	_ = v6984
	var v6992 int32
	_ = v6992
	var v6993 int32
	_ = v6993
	var v6994 int32
	_ = v6994
	var v6995 int32
	_ = v6995
	var v6997 int32
	_ = v6997
	var v6998 int32
	_ = v6998
	var v7002 int32
	_ = v7002
	var v7003 int32
	_ = v7003
	var v7004 int32
	_ = v7004
	var v7005 int32
	_ = v7005
	var v7014 int32
	_ = v7014
	var v7015 int32
	_ = v7015
	var v7018 int64
	_ = v7018
	var v7021 int32
	_ = v7021
	var v7023 int64
	_ = v7023
	var v7026 int32
	_ = v7026
	var v7033 int32
	_ = v7033
	var v7035 int32
	_ = v7035
	var v7039 int32
	_ = v7039
	var v7040 int32
	_ = v7040
	var v7041 int32
	_ = v7041
	var v7042 int32
	_ = v7042
	var v7051 int32
	_ = v7051
	var v7052 int32
	_ = v7052
	var v7053 int32
	_ = v7053
	var v7055 int32
	_ = v7055
	var v7060 int32
	_ = v7060
	var v7061 int32
	_ = v7061
	var v7068 int32
	_ = v7068
	var v7071 int32
	_ = v7071
	var v7073 int32
	_ = v7073
	var v7077 int32
	_ = v7077
	var v7078 int32
	_ = v7078
	var v7079 int32
	_ = v7079
	var v7080 int32
	_ = v7080
	var v7117 int32
	_ = v7117
	var v7119 int32
	_ = v7119
	var v7120 int32
	_ = v7120
	var v7127 int32
	_ = v7127
	var v7131 int32
	_ = v7131
	var v7140 int32
	_ = v7140
	var v7141 int32
	_ = v7141
	var v7142 int32
	_ = v7142
	var v7154 int32
	_ = v7154
	var v7155 int32
	_ = v7155
	var v7159 int32
	_ = v7159
	var v7160 int32
	_ = v7160
	var v7168 int32
	_ = v7168
	var v7169 int32
	_ = v7169
	var v7170 int32
	_ = v7170
	var v7171 int32
	_ = v7171
	var v7173 int32
	_ = v7173
	var v7174 int32
	_ = v7174
	var v7178 int32
	_ = v7178
	var v7179 int32
	_ = v7179
	var v7180 int32
	_ = v7180
	var v7181 int32
	_ = v7181
	var v7190 int32
	_ = v7190
	var v7191 int32
	_ = v7191
	var v7194 int64
	_ = v7194
	var v7197 int32
	_ = v7197
	var v7199 int64
	_ = v7199
	var v7202 int32
	_ = v7202
	var v7209 int32
	_ = v7209
	var v7211 int32
	_ = v7211
	var v7215 int32
	_ = v7215
	var v7216 int32
	_ = v7216
	var v7217 int32
	_ = v7217
	var v7218 int32
	_ = v7218
	var v7227 int32
	_ = v7227
	var v7228 int32
	_ = v7228
	var v7229 int32
	_ = v7229
	var v7231 int32
	_ = v7231
	var v7236 int32
	_ = v7236
	var v7237 int32
	_ = v7237
	var v7244 int32
	_ = v7244
	var v7247 int32
	_ = v7247
	var v7249 int32
	_ = v7249
	var v7253 int32
	_ = v7253
	var v7254 int32
	_ = v7254
	var v7255 int32
	_ = v7255
	var v7256 int32
	_ = v7256
	var v7293 int32
	_ = v7293
	var v7294 int32
	_ = v7294
	var v7295 int32
	_ = v7295
	var v7296 int32
	_ = v7296
	var v7297 int64
	_ = v7297
	var v7298 int32
	_ = v7298
	var v7305 int64
	_ = v7305
	var v7306 int32
	_ = v7306
	var v7307 int64
	_ = v7307
	var v7317 int32
	_ = v7317
	var v7319 int32
	_ = v7319
	var v7322 int32
	_ = v7322
	var v7324 int32
	_ = v7324
	var v7329 int32
	_ = v7329
	var v7330 int64
	_ = v7330
	var v7331 int32
	_ = v7331
	var v7338 int64
	_ = v7338
	var v7339 int32
	_ = v7339
	var v7340 int64
	_ = v7340
	var v7350 int32
	_ = v7350
	var v7352 int32
	_ = v7352
	var v7360 int64
	_ = v7360
	var v7362 int32
	_ = v7362
	var v7363 base.V128
	_ = v7363
	var v7366 int64
	_ = v7366
	var v7369 int32
	_ = v7369
	var v7372 base.V128
	_ = v7372
	var v7382 int32
	_ = v7382
	var v7384 int32
	_ = v7384
	var v7388 int32
	_ = v7388
	var v7397 int32
	_ = v7397
	var v7405 int32
	_ = v7405
	var v7409 int32
	_ = v7409
	var v7416 int32
	_ = v7416
	var v7419 int32
	_ = v7419
	var v7421 int32
	_ = v7421
	var v7424 int32
	_ = v7424
	var v7427 int32
	_ = v7427
	var v7428 int32
	_ = v7428
	var v7431 int32
	_ = v7431
	var v7435 int32
	_ = v7435
	var v7445 int32
	_ = v7445
	var v7446 int32
	_ = v7446
	var v7447 int32
	_ = v7447
	var v7449 int32
	_ = v7449
	var v7450 int32
	_ = v7450
	var v7451 int32
	_ = v7451
	var v7453 int64
	_ = v7453
	var v7459 int32
	_ = v7459
	var v7460 int32
	_ = v7460
	var v7461 int32
	_ = v7461
	var v7468 int32
	_ = v7468
	var v7470 int32
	_ = v7470
	var v7474 int32
	_ = v7474
	var v7475 int32
	_ = v7475
	var v7480 int32
	_ = v7480
	var v7481 int32
	_ = v7481
	var v7482 int32
	_ = v7482
	var v7489 int32
	_ = v7489
	var v7490 int32
	_ = v7490
	var v7491 int32
	_ = v7491
	var v7492 int32
	_ = v7492
	var v7493 int64
	_ = v7493
	var v7495 int32
	_ = v7495
	var v7497 int32
	_ = v7497
	var v7505 int32
	_ = v7505
	var v7508 int32
	_ = v7508
	var v7509 int32
	_ = v7509
	var v7514 int32
	_ = v7514
	var v7527 int32
	_ = v7527
	var v7560 int32
	_ = v7560
	var v7561 int32
	_ = v7561
	var v7587 int32
	_ = v7587
	var v7588 int32
	_ = v7588
	var v7596 int32
	_ = v7596
	var v7597 int32
	_ = v7597
	var v7598 int32
	_ = v7598
	var v7601 int32
	_ = v7601
	var v7603 int32
	_ = v7603
	var v7606 int32
	_ = v7606
	var v7628 int32
	_ = v7628
	var v7649 int32
	_ = v7649
	var v7650 int32
	_ = v7650
	var v7676 int32
	_ = v7676
	var v7679 int32
	_ = v7679
	var v7680 int32
	_ = v7680
	var v7683 int32
	_ = v7683
	var v7692 int32
	_ = v7692
	var v7693 int64
	_ = v7693
	var v7695 int32
	_ = v7695
	var v7697 int32
	_ = v7697
	var v7698 int32
	_ = v7698
	var v7699 int32
	_ = v7699
	var v7703 int32
	_ = v7703
	var v7707 int32
	_ = v7707
	var v7708 int32
	_ = v7708
	var v7713 int32
	_ = v7713
	var v7718 int32
	_ = v7718
	var v7721 int32
	_ = v7721
	var v7729 base.V128
	_ = v7729
	var v7737 int32
	_ = v7737
	var v7741 int32
	_ = v7741
	var v7745 int32
	_ = v7745
	var v7758 int32
	_ = v7758
	var v7759 int32
	_ = v7759
	var v7761 int32
	_ = v7761
	var v7764 int32
	_ = v7764
	var v7765 int32
	_ = v7765
	var v7767 int32
	_ = v7767
	var v7770 int64
	_ = v7770
	var v7771 int32
	_ = v7771
	var v7778 int64
	_ = v7778
	var v7779 int32
	_ = v7779
	var v7780 int64
	_ = v7780
	var v7789 int32
	_ = v7789
	var v7791 int32
	_ = v7791
	var v7798 int32
	_ = v7798
	var v7801 int32
	_ = v7801
	var v7803 int32
	_ = v7803
	var v7810 int32
	_ = v7810
	var v7823 int32
	_ = v7823
	var v7824 int32
	_ = v7824
	var v7828 int32
	_ = v7828
	var v7829 int32
	_ = v7829
	var v7837 int32
	_ = v7837
	var v7838 int32
	_ = v7838
	var v7839 int32
	_ = v7839
	var v7840 int32
	_ = v7840
	var v7842 int32
	_ = v7842
	var v7843 int32
	_ = v7843
	var v7847 int32
	_ = v7847
	var v7848 int32
	_ = v7848
	var v7849 int32
	_ = v7849
	var v7850 int32
	_ = v7850
	var v7859 int32
	_ = v7859
	var v7860 int32
	_ = v7860
	var v7863 int64
	_ = v7863
	var v7866 int32
	_ = v7866
	var v7868 int64
	_ = v7868
	var v7871 int32
	_ = v7871
	var v7878 int32
	_ = v7878
	var v7880 int32
	_ = v7880
	var v7884 int32
	_ = v7884
	var v7885 int32
	_ = v7885
	var v7886 int32
	_ = v7886
	var v7887 int32
	_ = v7887
	var v7896 int32
	_ = v7896
	var v7897 int32
	_ = v7897
	var v7898 int32
	_ = v7898
	var v7900 int32
	_ = v7900
	var v7905 int32
	_ = v7905
	var v7906 int32
	_ = v7906
	var v7913 int32
	_ = v7913
	var v7916 int32
	_ = v7916
	var v7918 int32
	_ = v7918
	var v7922 int32
	_ = v7922
	var v7923 int32
	_ = v7923
	var v7924 int32
	_ = v7924
	var v7925 int32
	_ = v7925
	var v7962 int32
	_ = v7962
	var v7974 int32
	_ = v7974
	var v7975 int32
	_ = v7975
	var v7979 int32
	_ = v7979
	var v7980 int32
	_ = v7980
	var v7988 int32
	_ = v7988
	var v7989 int32
	_ = v7989
	var v7990 int32
	_ = v7990
	var v7991 int32
	_ = v7991
	var v7993 int32
	_ = v7993
	var v7994 int32
	_ = v7994
	var v7998 int32
	_ = v7998
	var v7999 int32
	_ = v7999
	var v8000 int32
	_ = v8000
	var v8001 int32
	_ = v8001
	var v8010 int32
	_ = v8010
	var v8011 int32
	_ = v8011
	var v8014 int64
	_ = v8014
	var v8017 int32
	_ = v8017
	var v8019 int64
	_ = v8019
	var v8022 int32
	_ = v8022
	var v8029 int32
	_ = v8029
	var v8031 int32
	_ = v8031
	var v8035 int32
	_ = v8035
	var v8036 int32
	_ = v8036
	var v8037 int32
	_ = v8037
	var v8038 int32
	_ = v8038
	var v8047 int32
	_ = v8047
	var v8048 int32
	_ = v8048
	var v8049 int32
	_ = v8049
	var v8051 int32
	_ = v8051
	var v8056 int32
	_ = v8056
	var v8057 int32
	_ = v8057
	var v8064 int32
	_ = v8064
	var v8067 int32
	_ = v8067
	var v8069 int32
	_ = v8069
	var v8073 int32
	_ = v8073
	var v8074 int32
	_ = v8074
	var v8075 int32
	_ = v8075
	var v8076 int32
	_ = v8076
	var v8113 int32
	_ = v8113
	var v8114 int32
	_ = v8114
	var v8126 int32
	_ = v8126
	var v8127 int32
	_ = v8127
	var v8131 int32
	_ = v8131
	var v8132 int32
	_ = v8132
	var v8140 int32
	_ = v8140
	var v8141 int32
	_ = v8141
	var v8142 int32
	_ = v8142
	var v8143 int32
	_ = v8143
	var v8145 int32
	_ = v8145
	var v8146 int32
	_ = v8146
	var v8150 int32
	_ = v8150
	var v8151 int32
	_ = v8151
	var v8152 int32
	_ = v8152
	var v8153 int32
	_ = v8153
	var v8162 int32
	_ = v8162
	var v8163 int32
	_ = v8163
	var v8166 int64
	_ = v8166
	var v8169 int32
	_ = v8169
	var v8171 int64
	_ = v8171
	var v8174 int32
	_ = v8174
	var v8181 int32
	_ = v8181
	var v8183 int32
	_ = v8183
	var v8187 int32
	_ = v8187
	var v8188 int32
	_ = v8188
	var v8189 int32
	_ = v8189
	var v8190 int32
	_ = v8190
	var v8199 int32
	_ = v8199
	var v8200 int32
	_ = v8200
	var v8201 int32
	_ = v8201
	var v8203 int32
	_ = v8203
	var v8208 int32
	_ = v8208
	var v8209 int32
	_ = v8209
	var v8216 int32
	_ = v8216
	var v8219 int32
	_ = v8219
	var v8221 int32
	_ = v8221
	var v8225 int32
	_ = v8225
	var v8226 int32
	_ = v8226
	var v8227 int32
	_ = v8227
	var v8228 int32
	_ = v8228
	var v8265 int32
	_ = v8265
	var v8266 int32
	_ = v8266
	var v8268 int32
	_ = v8268
	var v8274 int32
	_ = v8274
	var v8285 int32
	_ = v8285
	var v8312 int32
	_ = v8312
	var v8341 int32
	_ = v8341
	var v8342 int32
	_ = v8342
	var v8346 int32
	_ = v8346
	var v8347 int32
	_ = v8347
	var v8354 int32
	_ = v8354
	var v8358 int32
	_ = v8358
	var v8362 int32
	_ = v8362
	var v8379 int32
	_ = v8379
	var v8406 int32
	_ = v8406
	var v8439 int32
	_ = v8439
	var v8440 int32
	_ = v8440
	var v8447 int32
	_ = v8447
	var v8449 int32
	_ = v8449
	var v8450 int32
	_ = v8450
	var v8451 int32
	_ = v8451
	var v8464 int32
	_ = v8464
	var v8465 int32
	_ = v8465
	var v8469 int32
	_ = v8469
	var v8470 int32
	_ = v8470
	var v8478 int32
	_ = v8478
	var v8479 int32
	_ = v8479
	var v8480 int32
	_ = v8480
	var v8481 int32
	_ = v8481
	var v8483 int32
	_ = v8483
	var v8484 int32
	_ = v8484
	var v8488 int32
	_ = v8488
	var v8489 int32
	_ = v8489
	var v8490 int32
	_ = v8490
	var v8491 int32
	_ = v8491
	var v8500 int32
	_ = v8500
	var v8501 int32
	_ = v8501
	var v8504 int64
	_ = v8504
	var v8507 int32
	_ = v8507
	var v8509 int64
	_ = v8509
	var v8512 int32
	_ = v8512
	var v8519 int32
	_ = v8519
	var v8521 int32
	_ = v8521
	var v8525 int32
	_ = v8525
	var v8526 int32
	_ = v8526
	var v8527 int32
	_ = v8527
	var v8528 int32
	_ = v8528
	var v8537 int32
	_ = v8537
	var v8538 int32
	_ = v8538
	var v8539 int32
	_ = v8539
	var v8541 int32
	_ = v8541
	var v8546 int32
	_ = v8546
	var v8547 int32
	_ = v8547
	var v8554 int32
	_ = v8554
	var v8557 int32
	_ = v8557
	var v8559 int32
	_ = v8559
	var v8563 int32
	_ = v8563
	var v8564 int32
	_ = v8564
	var v8565 int32
	_ = v8565
	var v8566 int32
	_ = v8566
	var v8609 int32
	_ = v8609
	var v8611 int32
	_ = v8611
	var v8612 int32
	_ = v8612
	var v8624 int32
	_ = v8624
	var v8625 int32
	_ = v8625
	var v8629 int32
	_ = v8629
	var v8630 int32
	_ = v8630
	var v8638 int32
	_ = v8638
	var v8639 int32
	_ = v8639
	var v8640 int32
	_ = v8640
	var v8641 int32
	_ = v8641
	var v8643 int32
	_ = v8643
	var v8644 int32
	_ = v8644
	var v8648 int32
	_ = v8648
	var v8649 int32
	_ = v8649
	var v8650 int32
	_ = v8650
	var v8651 int32
	_ = v8651
	var v8660 int32
	_ = v8660
	var v8661 int32
	_ = v8661
	var v8664 int64
	_ = v8664
	var v8667 int32
	_ = v8667
	var v8669 int64
	_ = v8669
	var v8672 int32
	_ = v8672
	var v8679 int32
	_ = v8679
	var v8681 int32
	_ = v8681
	var v8685 int32
	_ = v8685
	var v8686 int32
	_ = v8686
	var v8687 int32
	_ = v8687
	var v8688 int32
	_ = v8688
	var v8697 int32
	_ = v8697
	var v8698 int32
	_ = v8698
	var v8699 int32
	_ = v8699
	var v8701 int32
	_ = v8701
	var v8706 int32
	_ = v8706
	var v8707 int32
	_ = v8707
	var v8714 int32
	_ = v8714
	var v8717 int32
	_ = v8717
	var v8719 int32
	_ = v8719
	var v8723 int32
	_ = v8723
	var v8724 int32
	_ = v8724
	var v8725 int32
	_ = v8725
	var v8726 int32
	_ = v8726
	var v8766 int32
	_ = v8766
	var v8767 int32
	_ = v8767
	var v8773 int32
	_ = v8773
	var v8776 int32
	_ = v8776
	var v8783 int32
	_ = v8783
	var v8785 int32
	_ = v8785
	var v8786 int32
	_ = v8786
	var v8789 int32
	_ = v8789
	var v8790 int32
	_ = v8790
	var v8796 int32
	_ = v8796
	var v8801 int32
	_ = v8801
	var v8812 int32
	_ = v8812
	var v8841 base.V128
	_ = v8841
	var v8874 int32
	_ = v8874
	var v8875 base.V128
	_ = v8875
	var v8878 base.V128
	_ = v8878
	var v8881 base.V128
	_ = v8881
	var v8884 base.V128
	_ = v8884
	var v8885 base.V128
	_ = v8885
	var v8889 int32
	_ = v8889
	var v8892 base.V128
	_ = v8892
	var v8897 int32
	_ = v8897
	var v8903 int32
	_ = v8903
	var v8940 int32
	_ = v8940
	var v8978 int32
	_ = v8978
	var v8989 int32
	_ = v8989
	var v8990 int32
	_ = v8990
	var v9045 int32
	_ = v9045
	var v9047 int32
	_ = v9047
	var v9051 int32
	_ = v9051
	var v9056 int32
	_ = v9056
	var v9123 int64
	_ = v9123
	var v9124 int32
	_ = v9124
	var v9131 int64
	_ = v9131
	var v9132 int32
	_ = v9132
	var v9133 int64
	_ = v9133
	var v9143 int32
	_ = v9143
	var v9145 int32
	_ = v9145
	var v9190 int32
	_ = v9190
	var v9223 int32
	_ = v9223
	var v9225 int32
	_ = v9225
	var v9228 int32
	_ = v9228
	var v9245 int32
	_ = v9245
	var v9246 int32
	_ = v9246
	var v9271 int32
	_ = v9271
	var v9301 int32
	_ = v9301
	var v9307 int32
	_ = v9307
	var v9312 int32
	_ = v9312
	var v9313 int32
	_ = v9313
	var v9330 int32
	_ = v9330
	var v9331 int32
	_ = v9331
	var v9386 int32
	_ = v9386
	var v9388 int32
	_ = v9388
	var v9390 int32
	_ = v9390
	var v9397 int32
	_ = v9397
	var v9398 int32
	_ = v9398
	var v9472 int32
	_ = v9472
	var v9494 int32
	_ = v9494
	var v9570 int32
	_ = v9570
	var v9572 int32
	_ = v9572
	var v9584 int32
	_ = v9584
	var v9585 int32
	_ = v9585
	var v9589 int32
	_ = v9589
	var v9590 int32
	_ = v9590
	var v9598 int32
	_ = v9598
	var v9599 int32
	_ = v9599
	var v9600 int32
	_ = v9600
	var v9601 int32
	_ = v9601
	var v9603 int32
	_ = v9603
	var v9604 int32
	_ = v9604
	var v9608 int32
	_ = v9608
	var v9609 int32
	_ = v9609
	var v9610 int32
	_ = v9610
	var v9611 int32
	_ = v9611
	var v9620 int32
	_ = v9620
	var v9621 int32
	_ = v9621
	var v9624 int64
	_ = v9624
	var v9627 int32
	_ = v9627
	var v9629 int64
	_ = v9629
	var v9632 int32
	_ = v9632
	var v9639 int32
	_ = v9639
	var v9641 int32
	_ = v9641
	var v9645 int32
	_ = v9645
	var v9646 int32
	_ = v9646
	var v9647 int32
	_ = v9647
	var v9648 int32
	_ = v9648
	var v9657 int32
	_ = v9657
	var v9658 int32
	_ = v9658
	var v9659 int32
	_ = v9659
	var v9661 int32
	_ = v9661
	var v9666 int32
	_ = v9666
	var v9667 int32
	_ = v9667
	var v9674 int32
	_ = v9674
	var v9677 int32
	_ = v9677
	var v9679 int32
	_ = v9679
	var v9683 int32
	_ = v9683
	var v9684 int32
	_ = v9684
	var v9685 int32
	_ = v9685
	var v9686 int32
	_ = v9686
	var v9743 int32
	_ = v9743
	var v9762 int32
	_ = v9762
	var v9820 int32
	_ = v9820
	var v9875 int32
	_ = v9875
	var v9890 int32
	_ = v9890
	var v9894 int32
	_ = v9894
	var v9895 int32
	_ = v9895
	var v9896 int32
	_ = v9896
	var v9897 int32
	_ = v9897
	var v9898 int32
	_ = v9898
	var v9900 int32
	_ = v9900
	var v9905 int32
	_ = v9905
	var v9906 int32
	_ = v9906
	var v9907 int32
	_ = v9907
	var v9908 int32
	_ = v9908
	var v9915 int32
	_ = v9915
	var v9918 int32
	_ = v9918
	var v9920 int32
	_ = v9920
	var v9923 int32
	_ = v9923
	var v9927 int64
	_ = v9927
	var v9928 int64
	_ = v9928
	var v9931 int32
	_ = v9931
	var v9933 base.V128
	_ = v9933
	var v9935 int32
	_ = v9935
	var v9937 base.V128
	_ = v9937
	var v9951 int32
	_ = v9951
	var v9952 int32
	_ = v9952
	var v9953 int32
	_ = v9953
	var v9955 int32
	_ = v9955
	var v9958 int32
	_ = v9958
	var v9959 int32
	_ = v9959
	var v10001 int32
	_ = v10001
	var v10002 int32
	_ = v10002
	var v10028 int32
	_ = v10028
	var v10036 int32
	_ = v10036
	var v10037 int32
	_ = v10037
	var v10076 int32
	_ = v10076
	var v10077 int32
	_ = v10077
	var v10111 int32
	_ = v10111
	var v10114 int32
	_ = v10114
	var v10118 int64
	_ = v10118
	var v10119 int64
	_ = v10119
	var v10122 int32
	_ = v10122
	var v10124 base.V128
	_ = v10124
	var v10126 int32
	_ = v10126
	var v10128 base.V128
	_ = v10128
	var v10137 int32
	_ = v10137
	var v10139 int32
	_ = v10139
	var v10179 int32
	_ = v10179
	var v10180 int32
	_ = v10180
	var v10251 int32
	_ = v10251
	var v10252 int32
	_ = v10252
	var v10285 int32
	_ = v10285
	var v10287 int32
	_ = v10287
	var v10289 int32
	_ = v10289
	var v10290 int32
	_ = v10290
	var v10292 int32
	_ = v10292
	var v10293 int32
	_ = v10293
	var v10294 int32
	_ = v10294
	var v10295 int32
	_ = v10295
	var v10296 int32
	_ = v10296
	var v10297 int32
	_ = v10297
	var v10298 int32
	_ = v10298
	var v10299 int32
	_ = v10299
	var v10301 int32
	_ = v10301
	var v10304 int32
	_ = v10304
	var v10305 int32
	_ = v10305
	var v10306 int32
	_ = v10306
	var v10307 int32
	_ = v10307
	var v10308 int32
	_ = v10308
	var v10309 int32
	_ = v10309
	var v10310 int32
	_ = v10310
	var v10311 int32
	_ = v10311
	var v10312 int32
	_ = v10312
	var v10313 int32
	_ = v10313
	var v10314 int32
	_ = v10314
	var v10315 int32
	_ = v10315
	var v10316 int32
	_ = v10316
	var v10317 int32
	_ = v10317
	var v10318 int32
	_ = v10318
	var v10319 int32
	_ = v10319
	var v10320 int32
	_ = v10320
	var v10321 int32
	_ = v10321
	var v10322 int32
	_ = v10322
	var v10323 int32
	_ = v10323
	var v10324 int32
	_ = v10324
	var v10325 int32
	_ = v10325
	var v10327 int32
	_ = v10327
	var v10334 int32
	_ = v10334
	var v10336 int32
	_ = v10336
	var v10337 int32
	_ = v10337
	var v10359 int32
	_ = v10359
	var v10361 int32
	_ = v10361
	var v10362 int32
	_ = v10362
	var v10364 int32
	_ = v10364
	var v10365 int32
	_ = v10365
	var v10366 int32
	_ = v10366
	var v10367 int32
	_ = v10367
	var v10368 int32
	_ = v10368
	var v10369 int32
	_ = v10369
	var v10370 int32
	_ = v10370
	var v10371 int32
	_ = v10371
	var v10373 int32
	_ = v10373
	var v10374 int32
	_ = v10374
	var v10376 int32
	_ = v10376
	var v10377 int32
	_ = v10377
	var v10378 int32
	_ = v10378
	var v10379 int32
	_ = v10379
	var v10380 int32
	_ = v10380
	var v10381 int32
	_ = v10381
	var v10382 int32
	_ = v10382
	var v10383 int32
	_ = v10383
	var v10384 int32
	_ = v10384
	var v10385 int32
	_ = v10385
	var v10386 int32
	_ = v10386
	var v10387 int32
	_ = v10387
	var v10388 int32
	_ = v10388
	var v10389 int32
	_ = v10389
	var v10390 int32
	_ = v10390
	var v10391 int32
	_ = v10391
	var v10392 int32
	_ = v10392
	var v10393 int32
	_ = v10393
	var v10394 int32
	_ = v10394
	var v10395 int32
	_ = v10395
	var v10396 int32
	_ = v10396
	var v10397 int32
	_ = v10397
	var v10399 int32
	_ = v10399
	var v10406 int32
	_ = v10406
	var v10408 int32
	_ = v10408
	var v10409 int32
	_ = v10409
	var v10431 int32
	_ = v10431
	var v10433 int32
	_ = v10433
	var v10434 int32
	_ = v10434
	var v10435 int32
	_ = v10435
	var v10436 int32
	_ = v10436
	var v10437 int32
	_ = v10437
	var v10438 int32
	_ = v10438
	var v10439 int32
	_ = v10439
	var v10440 int32
	_ = v10440
	var v10441 int32
	_ = v10441
	var v10442 int32
	_ = v10442
	var v10443 int32
	_ = v10443
	var v10445 int32
	_ = v10445
	var v10446 int32
	_ = v10446
	var v10448 int32
	_ = v10448
	var v10449 int32
	_ = v10449
	var v10450 int32
	_ = v10450
	var v10451 int32
	_ = v10451
	var v10452 int32
	_ = v10452
	var v10453 int32
	_ = v10453
	var v10454 int32
	_ = v10454
	var v10455 int32
	_ = v10455
	var v10456 int32
	_ = v10456
	var v10457 int32
	_ = v10457
	var v10458 int32
	_ = v10458
	var v10459 int32
	_ = v10459
	var v10460 int32
	_ = v10460
	var v10461 int32
	_ = v10461
	var v10462 int32
	_ = v10462
	var v10463 int32
	_ = v10463
	var v10464 int32
	_ = v10464
	var v10465 int32
	_ = v10465
	var v10466 int32
	_ = v10466
	var v10467 int32
	_ = v10467
	var v10468 int32
	_ = v10468
	var v10469 int32
	_ = v10469
	var v10471 int32
	_ = v10471
	var v10478 int32
	_ = v10478
	var v10480 int32
	_ = v10480
	var v10481 int32
	_ = v10481
	var v10507 int32
	_ = v10507
	var v10508 int32
	_ = v10508
	var v10514 int32
	_ = v10514
	var v10519 int32
	_ = v10519
	var v10524 int32
	_ = v10524
	var v10530 base.V128
	_ = v10530
	var v10533 int32
	_ = v10533
	var v10534 int32
	_ = v10534
	var v10539 int32
	_ = v10539
	var v10540 int32
	_ = v10540
	var v10542 int32
	_ = v10542
	var v10544 int32
	_ = v10544
	var v10545 int32
	_ = v10545
	var v10548 int32
	_ = v10548
	var v10552 int64
	_ = v10552
	var v10553 int64
	_ = v10553
	var v10556 int32
	_ = v10556
	var v10558 base.V128
	_ = v10558
	var v10560 int32
	_ = v10560
	var v10562 base.V128
	_ = v10562
	var v10570 int32
	_ = v10570
	var v10571 int32
	_ = v10571
	var v10572 int32
	_ = v10572
	var v10574 int32
	_ = v10574
	var v10578 int32
	_ = v10578
	var v10580 int32
	_ = v10580
	var v10581 int32
	_ = v10581
	var v10585 int32
	_ = v10585
	var v10587 int32
	_ = v10587
	var v10588 int32
	_ = v10588
	var v10594 int32
	_ = v10594
	var v10596 int32
	_ = v10596
	var v10600 int32
	_ = v10600
	var v10602 int32
	_ = v10602
	var v10606 int32
	_ = v10606
	var v10610 int32
	_ = v10610
	var v10614 int32
	_ = v10614
	var v10615 int64
	_ = v10615
	var v10617 int32
	_ = v10617
	var v10619 int32
	_ = v10619
	var v10620 int32
	_ = v10620
	var v10621 int32
	_ = v10621
	var v10626 int32
	_ = v10626
	var v10628 int32
	_ = v10628
	var v10630 int32
	_ = v10630
	var v10631 int32
	_ = v10631
	var v10699 int32
	_ = v10699
	var v10700 int32
	_ = v10700
	var v10703 int32
	_ = v10703
	var v10707 int64
	_ = v10707
	var v10708 int64
	_ = v10708
	var v10711 int32
	_ = v10711
	var v10713 base.V128
	_ = v10713
	var v10715 int32
	_ = v10715
	var v10717 base.V128
	_ = v10717
	var v10723 int32
	_ = v10723
	var v10725 int32
	_ = v10725
	var v10795 int32
	_ = v10795
	var v10800 int32
	_ = v10800
	var v10806 base.V128
	_ = v10806
	var v10809 int32
	_ = v10809
	var v10810 int32
	_ = v10810
	v3 = int32(0)
	v72 = m.G0
	v74 = v72 - int32(_a_F_EncodeStreamHook_0)
	m.G0 = v74
	v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v79 = *(*int32)(unsafe.Add(mBase, uint32(v76+int32(8))))
	v80 = *(*int32)(unsafe.Add(mBase, uint32(l0)+416))
	v81 = *(*int32)(unsafe.Add(mBase, uint32(l0)+412))
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)+408))
	v83 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v84 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v85 = *(*int32)(unsafe.Add(mBase, uint32(v84)+12))
	v86 = *(*int32)(unsafe.Add(mBase, uint32(v84)+8))
	v87 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v88 = *(*float32)(unsafe.Add(mBase, uint32(v87)+4))
	v89 = *(*int32)(unsafe.Add(mBase, uint32(v87)+8))
	v90 = *(*int32)(unsafe.Add(mBase, uint32(v76)+4))
	v91 = *(*int32)(unsafe.Add(mBase, uint32(v76)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v74)+60)) = int32(2)
	v100 = *(*int64)(unsafe.Add(mBase, uint32(v76+int32(16))))
	*(*int64)(unsafe.Add(mBase, uint32(v74+int32(48)))) = v100
	v103 = base.Simd_g_v128_load(m, v76, v3)
	base.Simd_g_v128_store(m, v74, int32(32), v103)
	if base.F32_lt(base.F32_abs(v88), float32(2.1474836e+09)) == v3 {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	v115 = v74 + int32(8)
	v118 = base.Simd_g_const(&F_EncodeStreamHook__k0)
	base.Simd_g_v128_store(m, v115, int32(0), v118)
	*(*int64)(unsafe.Add(mBase, uint32(v74+int32(24)))) = int64(0)
	v128 = int32(1024)
	v130 = F_WebPSafeMalloc(m, int64(1), v128)
	mBase = m.M
	if v130 != 0 {
		goto L10
	} else {
		goto L11
	}
L2:
	;
	v113 = int32(-2147483648)
	goto L1
L3:
	;
	v111 = base.I32_trunc_f32_s(v88)
	v113 = v111
	goto L1
L4:
	;
	v10795 = v10725 + int32(8)
	if v10795 == int32(0) {
		goto L1163
	} else {
		goto L1164
	}
L5:
	;
	v10699 = int32(8)
	v10700 = v10630 + v10699
	v10703 = int32(0)
	v10707 = *(*int64)(unsafe.Add(mBase, uint32(v10700)))
	v10708 = *(*int64)(unsafe.Add(mBase, uint32(v10631)))
	*(*int64)(unsafe.Add(mBase, uint32(v10700))) = v10708
	v10711 = v10630 + int32(16)
	v10713 = base.Simd_g_v128_load(m, v10711, v10703)
	v10715 = v10631 + v10699
	v10717 = base.Simd_g_v128_load(m, v10715, v10703)
	base.Simd_g_v128_store(m, v10711, v10703, v10717)
	*(*int64)(unsafe.Add(mBase, uint32(v10631))) = v10707
	base.Simd_g_v128_store(m, v10715, v10703, v10713)
	goto L1161
L6:
	;
	v232 = int32(0)
	v235 = base.I32_div_u_s(int32(97), v82)
	v237 = int32(base.Ui32(v235) >> (uint(int32(2)) % 32))
	v247 = int32(-1)
	v269 = l0
	v271 = v74
	v272 = v76
	v274 = v80
	v275 = v81
	v276 = v82
	v277 = v83
	v278 = v84
	v279 = v85
	v280 = v86
	v281 = v87
	v283 = v89
	v286 = v113
	v287 = (v90+int32(7))>>(uint(int32(3))%32) + (v91 - v79)
	v288 = l0 + int32(16)
	v289 = base.B2i32(v89 == v232)
	v290 = v235
	v291 = v237
	v292 = v235 - v237
	v293 = v83 + int32(68)
	v294 = v83 + int32(52)
	v295 = v83 + int32(72)
	v296 = v83 + int32(2168)
	v297 = v85 + v247
	v298 = v83 + int32(2216)
	v299 = v83 + int32(1096)
	v300 = v83 + int32(2144)
	v301 = v83 + int32(2120)
	v302 = v74 + int32(2096)
	v303 = v74 + int32(2112) | int32(4)
	v304 = v247
	v305 = v232
	v306 = v232
	v307 = v232
	goto L33
L7:
	;
	if v82 != int32(1) {
		v10628 = l0
		v10630 = v74
		v10631 = v76
		goto L5
	} else {
		goto L32
	}
L8:
	;
	v219 = *(*int32)(unsafe.Add(mBase, uint32(v84)+92))
	if v219 != 0 {
		goto L30
	} else {
		goto L31
	}
L9:
	;
	if v141 == int32(0) {
		goto L8
	} else {
		goto L12
	}
L10:
	;
	v134 = *(*int32)(unsafe.Add(mBase, uint32(v115)+8))
	F_WebPSafeFree(m, v134)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v115)+16)) = v130 + v128
	*(*int32)(unsafe.Add(mBase, uint32(v115)+12)) = v130
	*(*int32)(unsafe.Add(mBase, uint32(v115)+8)) = v130
	v141 = int32(1)
	goto L9
L11:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v115)+20)) = int32(1)
	v141 = int32(0)
	goto L9
L12:
	;
	if v82 < int32(2) {
		goto L7
	} else {
		goto L13
	}
L13:
	;
	v147 = v74 + int32(8)
	v157 = *(*int32)(unsafe.Add(mBase, uint32(v147)+12))
	v158 = *(*int32)(unsafe.Add(mBase, uint32(v147)+8))
	v159 = v157 - v158
	v161 = *(*int32)(unsafe.Add(mBase, uint32(v76)+12))
	v162 = *(*int32)(unsafe.Add(mBase, uint32(v76)+8))
	v163 = v161 - v162
	v165 = base.I64_extend_i32_u(v159) + base.I64_extend_i32_u(v163)
	if base.Ui64(v165) < base.Ui64(int64(4294967296)) {
		goto L15
	} else {
		goto L16
	}
L14:
	;
	if v217 != 0 {
		goto L6
	} else {
		goto L28
	}
L15:
	;
	v171 = base.I32_wrap_i64(v165)
	v172 = *(*int32)(unsafe.Add(mBase, uint32(v147)+16))
	v173 = v172 - v158
	if v172 == v158 {
		goto L18
	} else {
		goto L19
	}
L16:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v147)+20)) = int32(1)
	v217 = int32(0)
	goto L14
L17:
	;
	v204 = F_memcpy(m, v203, v202, v163)
	mBase = m.M
	v205 = *(*int64)(unsafe.Add(mBase, uint32(v76)))
	*(*int64)(unsafe.Add(mBase, uint32(v147))) = v205
	v207 = *(*int32)(unsafe.Add(mBase, uint32(v76)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v147)+20)) = v207
	v209 = *(*int32)(unsafe.Add(mBase, uint32(v147)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v147)+12)) = v209 + v163
	v217 = int32(1)
	goto L14
L18:
	;
	v180 = int32(base.Ui32(v173*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v171) < base.Ui32(v180) {
		goto L22
	} else {
		goto L23
	}
L19:
	;
	if base.Ui32(v173) < base.Ui32(v171) {
		goto L18
	} else {
		goto L20
	}
L20:
	;
	v202 = v162
	v203 = v158
	goto L17
L21:
	;
	if v157 == v158 {
		goto L26
	} else {
		goto L27
	}
L22:
	;
	v182 = v180
	goto L24
L23:
	;
	v182 = v171
	goto L24
L24:
	;
	v186 = v182&int32(-1024) + int32(1024)
	v187 = F_WebPSafeMalloc(m, int64(1), v186)
	mBase = m.M
	if v187 != 0 {
		goto L21
	} else {
		goto L25
	}
L25:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v147)+20)) = int32(1)
	v217 = int32(0)
	goto L14
L26:
	;
	v194 = *(*int32)(unsafe.Add(mBase, uint32(v147)+8))
	F_WebPSafeFree(m, v194)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v147)+8)) = v187
	*(*int32)(unsafe.Add(mBase, uint32(v147)+16)) = v187 + v186
	*(*int32)(unsafe.Add(mBase, uint32(v147)+12)) = v187 + v159
	v201 = *(*int32)(unsafe.Add(mBase, uint32(v76)+8))
	v202 = v201
	v203 = v187
	goto L17
L27:
	;
	v192 = *(*int32)(unsafe.Add(mBase, uint32(v147)+8))
	v193 = F_memcpy(m, v187, v192, v159)
	mBase = m.M
	goto L26
L28:
	;
	goto L8
L29:
	;
	v10723 = l0
	v10725 = v74
	goto L4
L30:
	;
	goto L29
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v84)+92)) = int32(1)
	goto L30
L32:
	;
	goto L6
L33:
	;
	v342 = v288 + v305*int32(28)
	v343 = *(*int32)(unsafe.Add(mBase, uint32(v342)))
	v344 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v271)+4)) = v344
	v347 = v343 & int32(-2)
	v349 = base.B2i32(v347 == int32(4))
	*(*int32)(unsafe.Add(mBase, uint32(v277)+64)) = v349
	v357 = base.B2i32(v343 == int32(5)) | base.B2i32(v343&int32(-3) == int32(1))
	*(*int32)(unsafe.Add(mBase, uint32(v277)+60)) = v357
	*(*int32)(unsafe.Add(mBase, uint32(v277)+56)) = base.B2i32(v347 == int32(2))
	*(*int32)(unsafe.Add(mBase, uint32(v277)+48)) = v344
	if v275 != 0 {
		goto L35
	} else {
		goto L36
	}
L34:
	;
	v10628 = v10431
	v10630 = v10433
	v10631 = v10434
	goto L5
L35:
	;
	v366 = v344
	goto L37
L36:
	;
	v366 = v357
	goto L37
L37:
	;
	if v349&int32(1) != 0 {
		goto L38
	} else {
		goto L39
	}
L38:
	;
	v369 = v344
	goto L40
L39:
	;
	v369 = v366
	goto L40
L40:
	;
	if v283 != 0 {
		goto L41
	} else {
		goto L42
	}
L41:
	;
	v371 = v369
	goto L43
L42:
	;
	v371 = int32(0)
	goto L43
L43:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+52)) = v371
	v373 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v271))) = v373
	v377 = *(*int32)(unsafe.Add(mBase, uint32(v301)+12))
	if v377 == v373 {
		goto L45
	} else {
		goto L46
	}
L44:
	;
	v404 = *(*int32)(unsafe.Add(mBase, uint32(v300)+12))
	if v404 == int32(0) {
		goto L53
	} else {
		goto L54
	}
L45:
	;
	v382 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v301)+20)) = v382
	v384 = *(*int32)(unsafe.Add(mBase, uint32(v301)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v301)+16)) = v384
	*(*int32)(unsafe.Add(mBase, uint32(v301)+12)) = v301 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v301)+8)) = v382
	if v384 == v382 {
		goto L47
	} else {
		goto L48
	}
L46:
	;
	v380 = *(*int32)(unsafe.Add(mBase, uint32(v301)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v377))) = v380
	goto L45
L47:
	;
	goto L44
L48:
	;
	v395 = v384
	goto L49
L49:
	;
	v396 = *(*int32)(unsafe.Add(mBase, uint32(v395)))
	F_WebPSafeFree(m, v395)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v301)+16)) = v396
	if v396 != 0 {
		v395 = v396
		goto L49
	} else {
		goto L51
	}
L50:
	;
	goto L47
L51:
	;
	goto L50
L52:
	;
	v430 = *(*int32)(unsafe.Add(mBase, uint32(v281)+92))
	if int32(99) < v430 {
		v679 = int32(0)
		goto L66
	} else {
		goto L67
	}
L53:
	;
	v409 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v300)+20)) = v409
	v411 = *(*int32)(unsafe.Add(mBase, uint32(v300)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v300)+16)) = v411
	*(*int32)(unsafe.Add(mBase, uint32(v300)+12)) = v300 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v300)+8)) = v409
	if v411 == v409 {
		goto L55
	} else {
		goto L56
	}
L54:
	;
	v407 = *(*int32)(unsafe.Add(mBase, uint32(v300)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v404))) = v407
	goto L53
L55:
	;
	goto L52
L56:
	;
	v422 = v411
	goto L57
L57:
	;
	v423 = *(*int32)(unsafe.Add(mBase, uint32(v422)))
	F_WebPSafeFree(m, v422)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v300)+16)) = v423
	if v423 != 0 {
		v422 = v423
		goto L57
	} else {
		goto L59
	}
L58:
	;
	goto L55
L59:
	;
	goto L58
L60:
	;
	v5049 = *(*int32)(unsafe.Add(mBase, uint32(v277)+56))
	if v5049 == int32(0) {
		goto L477
	} else {
		goto L478
	}
L61:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(3)
	v4969 = *(*int32)(unsafe.Add(mBase, uint32(v277)+68))
	if int32(1023) < v4969 {
		v5021 = v292
		goto L60
	} else {
		goto L476
	}
L62:
	;
	v4891 = *(*int32)(unsafe.Add(mBase, uint32(v1540)+92))
	if v4891 != 0 {
		goto L473
	} else {
		goto L474
	}
L63:
	;
	v4549 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v4550 = *(*int32)(unsafe.Add(mBase, uint32(v4549)+8))
	v4551 = *(*int32)(unsafe.Add(mBase, uint32(v4549)+12))
	v4558 = base.I64_extend_i32_s(v4551) * base.I64_extend_i32_s(v4550)
	v4559 = *(*int32)(unsafe.Add(mBase, uint32(v277)+60))
	if v4559 == int32(0) {
		goto L453
	} else {
		goto L454
	}
L64:
	;
	if v430 < int32(100) {
		v5021 = v290
		goto L60
	} else {
		goto L449
	}
L65:
	;
	v685 = *(*int32)(unsafe.Add(mBase, uint32(v342)+4))
	v686 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v687 = *(*int32)(unsafe.Add(mBase, uint32(v277)+68))
	v688 = F_PaletteSort(m, v685, v686, v299, v687, v295)
	mBase = m.M
	if v688 != 0 {
		goto L110
	} else {
		goto L111
	}
L66:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = v679
	v681 = *(*int32)(unsafe.Add(mBase, uint32(v277)+64))
	if v681 == int32(0) {
		goto L64
	} else {
		goto L109
	}
L67:
	;
	v433 = *(*int32)(unsafe.Add(mBase, uint32(v277)+64))
	if v433 == int32(0) {
		goto L68
	} else {
		goto L69
	}
L68:
	;
	v438 = *(*int32)(unsafe.Add(mBase, uint32(v277)+60))
	if v438 == int32(0) {
		goto L70
	} else {
		goto L71
	}
L69:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(0)
	goto L65
L70:
	;
	v449 = base.I64_extend_i32_s(v279) * base.I64_extend_i32_s(v280)
	v450 = *(*int32)(unsafe.Add(mBase, uint32(v277)+60))
	if v450 == int32(0) {
		goto L75
	} else {
		goto L76
	}
L71:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(0)
	goto L63
L72:
	;
	if v529 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L84
	}
L73:
	;
	v489 = v449 + v484 + v485 + int64(16)
	v490 = *(*int32)(unsafe.Add(mBase, uint32(v277)+24))
	if v490 == int32(0) {
		goto L79
	} else {
		goto L80
	}
L74:
	;
	v472 = int32(3)
	v474 = int32(2)
	v483 = v469
	v484 = v470
	v485 = base.I64_extend_i32_u(int32(base.Ui32(v279+v472)>>(uint(v474)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v280+v472)>>(uint(v474)%32)))
	goto L73
L75:
	;
	v463 = int32(0)
	v464 = int64(0)
	v466 = *(*int32)(unsafe.Add(mBase, uint32(v277)+52))
	if v466 == v463 {
		v483 = v463
		v484 = v464
		v485 = v464
		goto L73
	} else {
		goto L77
	}
L76:
	;
	v454 = v280 << (uint(int32(1)) % 32)
	v457 = int32(2)
	v461 = v454 + int32(base.Ui32(v454+int32(3))>>(uint(v457)%32)) + v457
	v469 = v461
	v470 = base.I64_extend_i32_u(v461)
	goto L74
L77:
	;
	v469 = v463
	v470 = v464
	goto L74
L78:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+32)) = v280
	*(*int32)(unsafe.Add(mBase, uint32(v277)+8)) = v507
	v511 = int32(2)
	v514 = int32(31)
	v516 = int32(-32)
	v517 = (v507 + base.I32_wrap_i64(v449)<<(uint(v511)%32) + v514) & v516
	*(*int32)(unsafe.Add(mBase, uint32(v277)+16)) = v517
	*(*int32)(unsafe.Add(mBase, uint32(v277)+20)) = (v517 + v483<<(uint(v511)%32) + v514) & v516
	v529 = int32(1)
	goto L72
L79:
	;
	F_WebPSafeFree(m, v490)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v277)+24)) = int64(0)
	v499 = F_WebPSafeMalloc(m, v489, int32(4))
	mBase = m.M
	if v499 != 0 {
		goto L82
	} else {
		goto L83
	}
L80:
	;
	v493 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v277)+28)))
	if base.Ui64(v489) <= base.Ui64(v493) {
		v507 = v490
		goto L78
	} else {
		goto L81
	}
L81:
	;
	goto L79
L82:
	;
	*(*uint32)(unsafe.Add(mBase, uint32(v277)+28)) = uint32(v489)
	*(*int32)(unsafe.Add(mBase, uint32(v277)+24)) = v499
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(0)
	v507 = v499
	goto L78
L83:
	;
	v500 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v502 = F_WebPEncodingSetError(m, v500, int32(1))
	mBase = m.M
	v529 = v502
	goto L72
L84:
	;
	v532 = int32(2)
	v533 = *(*int32)(unsafe.Add(mBase, uint32(v277)+12))
	if v533 == v532 {
		v679 = v532
		goto L66
	} else {
		goto L85
	}
L85:
	;
	v536 = *(*int32)(unsafe.Add(mBase, uint32(v281)+92))
	v537 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v546 = base.I32_div_s(v536, int32(-20))
	v547 = *(*int32)(unsafe.Add(mBase, uint32(v278)+56))
	v548 = *(*int32)(unsafe.Add(mBase, uint32(v278)+12))
	v549 = *(*int32)(unsafe.Add(mBase, uint32(v278)+8))
	v554 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v549*int32(3)), int32(4))
	mBase = m.M
	if v554 != 0 {
		goto L87
	} else {
		goto L88
	}
L86:
	;
	if v674 != 0 {
		v679 = v532
		goto L66
	} else {
		goto L105
	}
L87:
	;
	v556 = int32(64)
	if base.B2i32(v549 < v556)&base.B2i32(v548 < v556) != 0 {
		goto L91
	} else {
		goto L92
	}
L88:
	;
	v674 = int32(0)
	goto L86
L89:
	;
	F_WebPSafeFree(m, v554)
	mBase = m.M
	v674 = int32(1)
	goto L86
L90:
	;
	v631 = *(*int32)(unsafe.Add(mBase, uint32(v278)+52))
	F_NearLossless(m, v549, v548, v631, v547, v546+int32(5), v554, v537)
	mBase = m.M
	v636 = v546 + int32(4)
	if v636 == int32(0) {
		goto L89
	} else {
		goto L101
	}
L91:
	;
	if v548 < int32(1) {
		goto L89
	} else {
		goto L94
	}
L92:
	;
	if int32(2) < v548 {
		goto L90
	} else {
		goto L93
	}
L93:
	;
	goto L91
L94:
	;
	v565 = int32(1)
	v568 = v549 << (uint(int32(2)) % 32)
	if v548 == v565 {
		v611 = int32(0)
		goto L95
	} else {
		goto L96
	}
L95:
	;
	if v548&v565 == int32(0) {
		goto L89
	} else {
		goto L100
	}
L96:
	;
	v580 = int32(0)
	v582 = v537
	goto L97
L97:
	;
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v278)+52))
	v588 = *(*int32)(unsafe.Add(mBase, uint32(v278)+56))
	v590 = int32(2)
	v593 = F_memcpy(m, v582, v587+v588*v580<<(uint(v590)%32), v568)
	mBase = m.M
	v595 = *(*int32)(unsafe.Add(mBase, uint32(v278)+52))
	v596 = *(*int32)(unsafe.Add(mBase, uint32(v278)+56))
	v603 = F_memcpy(m, v593+v568, v595+v596*(v580+int32(1))<<(uint(v590)%32), v568)
	mBase = m.M
	v606 = v580 + v590
	if v548&int32(2147483646) != v606 {
		v580 = v606
		v582 = v593 + v549<<(uint(int32(3))%32)
		goto L97
	} else {
		goto L99
	}
L98:
	;
	v611 = v606
	goto L95
L99:
	;
	goto L98
L100:
	;
	v621 = int32(2)
	v624 = *(*int32)(unsafe.Add(mBase, uint32(v278)+52))
	v625 = *(*int32)(unsafe.Add(mBase, uint32(v278)+56))
	v630 = F_memcpy(m, v537+v611*v549<<(uint(v621)%32), v624+v625*v611<<(uint(v621)%32), v568)
	mBase = m.M
	goto L89
L101:
	;
	v639 = v636
	goto L102
L102:
	;
	F_NearLossless(m, v549, v548, v537, v549, v639, v554, v537)
	mBase = m.M
	v651 = v639 + int32(-1)
	if v651 != 0 {
		v639 = v651
		goto L102
	} else {
		goto L104
	}
L103:
	;
	goto L89
L104:
	;
	goto L103
L105:
	;
	v676 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v676 != 0 {
		goto L107
	} else {
		goto L108
	}
L106:
	;
	v10723 = v269
	v10725 = v271
	goto L4
L107:
	;
	goto L106
L108:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L107
L109:
	;
	goto L65
L110:
	;
	v694 = *(*int32)(unsafe.Add(mBase, uint32(v277)+68))
	v696 = v694 + int32(-1)
	v700 = *(*int32)(unsafe.Add(mBase, uint32(v295+v696<<(uint(int32(2))%32))))
	v701 = int32(1)
	goto L117
L111:
	;
	v689 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v691 = *(*int32)(unsafe.Add(mBase, uint32(v689)+92))
	if v691 != 0 {
		goto L113
	} else {
		goto L114
	}
L112:
	;
	v10723 = v269
	v10725 = v271
	goto L4
L113:
	;
	goto L112
L114:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v689)+92)) = int32(1)
	goto L113
L115:
	;
	v853 = int32(3)
	v854 = int32(2)
	goto L143
L116:
	;
	goto L115
L117:
	;
	v714 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v715 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v715+v701 < int32(32) {
		goto L120
	} else {
		goto L121
	}
L118:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v816 + v814
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v813<<(uint(v816)%32) | v815
	goto L116
L119:
	;
	v733 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v734 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v738 = v730
	v739 = v731
	v740 = v734
	v741 = v733
	goto L123
L120:
	;
	if v715 < int32(16) {
		v813 = v701
		v814 = v701
		v815 = v714
		v816 = v715
		goto L118
	} else {
		goto L122
	}
L121:
	;
	v719 = int32(32)
	v720 = v719 - v715
	v728 = int32(base.Ui32(v701) >> (uint(v720) % 32))
	v729 = v701 - v720
	v730 = v701<<(uint(v715)%32) | v714
	v731 = v719
	goto L119
L122:
	;
	v728 = v701
	v729 = v701
	v730 = v714
	v731 = v715
	goto L119
L123:
	;
	if base.Ui32(v740+int32(2)) <= base.Ui32(v741) {
		v796 = v740
		v797 = v741
		goto L125
	} else {
		goto L126
	}
L124:
	;
	v813 = v728
	v814 = v729
	v815 = v809
	v816 = v807
	goto L118
L125:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v796))) = uint16(v738)
	v804 = v796 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v804
	v807 = v739 + int32(-16)
	v809 = int32(base.Ui32(v738) >> (uint(int32(16)) % 32))
	if int32(31) < v739 {
		v738 = v809
		v739 = v807
		v740 = v804
		v741 = v797
		goto L123
	} else {
		goto L140
	}
L126:
	;
	v750 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v751 = v741 - v750
	v754 = base.I64_extend_i32_s(v751) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v754) {
		v778 = v750
		goto L128
	} else {
		goto L129
	}
L127:
	;
	if v740 == v750 {
		goto L138
	} else {
		goto L139
	}
L128:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v778
	goto L115
L129:
	;
	v757 = v740 - v750
	v759 = v754 + base.I64_extend_i32_u(v757)
	if base.Ui64(int64(4294967295)) < base.Ui64(v759) {
		v778 = v750
		goto L128
	} else {
		goto L130
	}
L130:
	;
	v762 = base.I32_wrap_i64(v759)
	if v741 == v750 {
		goto L131
	} else {
		goto L132
	}
L131:
	;
	v769 = int32(base.Ui32(v751*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v762) < base.Ui32(v769) {
		goto L134
	} else {
		goto L135
	}
L132:
	;
	if base.Ui32(v762) <= base.Ui32(v751) {
		v796 = v740
		v797 = v741
		goto L125
	} else {
		goto L133
	}
L133:
	;
	goto L131
L134:
	;
	v771 = v769
	goto L136
L135:
	;
	v771 = v762
	goto L136
L136:
	;
	v775 = v771&int32(-1024) + int32(1024)
	v776 = F_WebPSafeMalloc(m, int64(1), v775)
	mBase = m.M
	if v776 != 0 {
		goto L127
	} else {
		goto L137
	}
L137:
	;
	v777 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v778 = v777
	goto L128
L138:
	;
	v789 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v789)
	mBase = m.M
	v791 = v776 + v775
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v791
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v776
	v796 = v776 + v757
	v797 = v791
	goto L125
L139:
	;
	v787 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v788 = F_memcpy(m, v776, v787, v757)
	mBase = m.M
	goto L138
L140:
	;
	goto L124
L141:
	;
	if v700 != 0 {
		goto L167
	} else {
		goto L168
	}
L142:
	;
	goto L141
L143:
	;
	v866 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v867 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v867+v854 < int32(32) {
		goto L146
	} else {
		goto L147
	}
L144:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v968 + v966
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v965<<(uint(v968)%32) | v967
	goto L142
L145:
	;
	v885 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v886 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v890 = v882
	v891 = v883
	v892 = v886
	v893 = v885
	goto L149
L146:
	;
	if v867 < int32(16) {
		v965 = v853
		v966 = v854
		v967 = v866
		v968 = v867
		goto L144
	} else {
		goto L148
	}
L147:
	;
	v871 = int32(32)
	v872 = v871 - v867
	v880 = int32(base.Ui32(v853) >> (uint(v872) % 32))
	v881 = v854 - v872
	v882 = v853<<(uint(v867)%32) | v866
	v883 = v871
	goto L145
L148:
	;
	v880 = v853
	v881 = v854
	v882 = v866
	v883 = v867
	goto L145
L149:
	;
	if base.Ui32(v892+int32(2)) <= base.Ui32(v893) {
		v948 = v892
		v949 = v893
		goto L151
	} else {
		goto L152
	}
L150:
	;
	v965 = v880
	v966 = v881
	v967 = v961
	v968 = v959
	goto L144
L151:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v948))) = uint16(v890)
	v956 = v948 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v956
	v959 = v891 + int32(-16)
	v961 = int32(base.Ui32(v890) >> (uint(int32(16)) % 32))
	if int32(31) < v891 {
		v890 = v961
		v891 = v959
		v892 = v956
		v893 = v949
		goto L149
	} else {
		goto L166
	}
L152:
	;
	v902 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v903 = v893 - v902
	v906 = base.I64_extend_i32_s(v903) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v906) {
		v930 = v902
		goto L154
	} else {
		goto L155
	}
L153:
	;
	if v892 == v902 {
		goto L164
	} else {
		goto L165
	}
L154:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v930
	goto L141
L155:
	;
	v909 = v892 - v902
	v911 = v906 + base.I64_extend_i32_u(v909)
	if base.Ui64(int64(4294967295)) < base.Ui64(v911) {
		v930 = v902
		goto L154
	} else {
		goto L156
	}
L156:
	;
	v914 = base.I32_wrap_i64(v911)
	if v893 == v902 {
		goto L157
	} else {
		goto L158
	}
L157:
	;
	v921 = int32(base.Ui32(v903*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v914) < base.Ui32(v921) {
		goto L160
	} else {
		goto L161
	}
L158:
	;
	if base.Ui32(v914) <= base.Ui32(v903) {
		v948 = v892
		v949 = v893
		goto L151
	} else {
		goto L159
	}
L159:
	;
	goto L157
L160:
	;
	v923 = v921
	goto L162
L161:
	;
	v923 = v914
	goto L162
L162:
	;
	v927 = v923&int32(-1024) + int32(1024)
	v928 = F_WebPSafeMalloc(m, int64(1), v927)
	mBase = m.M
	if v928 != 0 {
		goto L153
	} else {
		goto L163
	}
L163:
	;
	v929 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v930 = v929
	goto L154
L164:
	;
	v941 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v941)
	mBase = m.M
	v943 = v928 + v927
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v943
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v928
	v948 = v928 + v909
	v949 = v943
	goto L151
L165:
	;
	v939 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v940 = F_memcpy(m, v928, v939, v909)
	mBase = m.M
	goto L164
L166:
	;
	goto L150
L167:
	;
	v1005 = v694
	goto L169
L168:
	;
	v1005 = v696
	goto L169
L169:
	;
	if int32(17) < v694 {
		goto L170
	} else {
		goto L171
	}
L170:
	;
	v1008 = v1005
	goto L172
L171:
	;
	v1008 = v694
	goto L172
L172:
	;
	v1010 = v1008 + int32(-1)
	v1011 = int32(8)
	goto L175
L173:
	;
	if v1010 < int32(1) {
		goto L199
	} else {
		goto L200
	}
L174:
	;
	goto L173
L175:
	;
	v1023 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v1024 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v1024+v1011 < int32(32) {
		goto L178
	} else {
		goto L179
	}
L176:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v1125 + v1123
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v1122<<(uint(v1125)%32) | v1124
	goto L174
L177:
	;
	v1042 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v1043 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v1047 = v1039
	v1048 = v1040
	v1049 = v1043
	v1050 = v1042
	goto L181
L178:
	;
	if v1024 < int32(16) {
		v1122 = v1010
		v1123 = v1011
		v1124 = v1023
		v1125 = v1024
		goto L176
	} else {
		goto L180
	}
L179:
	;
	v1028 = int32(32)
	v1029 = v1028 - v1024
	v1037 = int32(base.Ui32(v1010) >> (uint(v1029) % 32))
	v1038 = v1011 - v1029
	v1039 = v1010<<(uint(v1024)%32) | v1023
	v1040 = v1028
	goto L177
L180:
	;
	v1037 = v1010
	v1038 = v1011
	v1039 = v1023
	v1040 = v1024
	goto L177
L181:
	;
	if base.Ui32(v1049+int32(2)) <= base.Ui32(v1050) {
		v1105 = v1049
		v1106 = v1050
		goto L183
	} else {
		goto L184
	}
L182:
	;
	v1122 = v1037
	v1123 = v1038
	v1124 = v1118
	v1125 = v1116
	goto L176
L183:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1105))) = uint16(v1047)
	v1113 = v1105 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v1113
	v1116 = v1048 + int32(-16)
	v1118 = int32(base.Ui32(v1047) >> (uint(int32(16)) % 32))
	if int32(31) < v1048 {
		v1047 = v1118
		v1048 = v1116
		v1049 = v1113
		v1050 = v1106
		goto L181
	} else {
		goto L198
	}
L184:
	;
	v1059 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v1060 = v1050 - v1059
	v1063 = base.I64_extend_i32_s(v1060) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1063) {
		v1087 = v1059
		goto L186
	} else {
		goto L187
	}
L185:
	;
	if v1049 == v1059 {
		goto L196
	} else {
		goto L197
	}
L186:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v1087
	goto L173
L187:
	;
	v1066 = v1049 - v1059
	v1068 = v1063 + base.I64_extend_i32_u(v1066)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1068) {
		v1087 = v1059
		goto L186
	} else {
		goto L188
	}
L188:
	;
	v1071 = base.I32_wrap_i64(v1068)
	if v1050 == v1059 {
		goto L189
	} else {
		goto L190
	}
L189:
	;
	v1078 = int32(base.Ui32(v1060*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1071) < base.Ui32(v1078) {
		goto L192
	} else {
		goto L193
	}
L190:
	;
	if base.Ui32(v1071) <= base.Ui32(v1060) {
		v1105 = v1049
		v1106 = v1050
		goto L183
	} else {
		goto L191
	}
L191:
	;
	goto L189
L192:
	;
	v1080 = v1078
	goto L194
L193:
	;
	v1080 = v1071
	goto L194
L194:
	;
	v1084 = v1080&int32(-1024) + int32(1024)
	v1085 = F_WebPSafeMalloc(m, int64(1), v1084)
	mBase = m.M
	if v1085 != 0 {
		goto L185
	} else {
		goto L195
	}
L195:
	;
	v1086 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v1087 = v1086
	goto L186
L196:
	;
	v1098 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v1098)
	mBase = m.M
	v1100 = v1085 + v1084
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v1100
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v1085
	v1105 = v1085 + v1066
	v1106 = v1100
	goto L183
L197:
	;
	v1096 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v1097 = F_memcpy(m, v1085, v1096, v1066)
	mBase = m.M
	goto L196
L198:
	;
	goto L182
L199:
	;
	v1528 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	*(*int32)(unsafe.Add(mBase, uint32(v271)+2112)) = v1528
	v1534 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v1537 = F_EncodeImageNoHuffman(m, v272, v271+int32(2112), v298, v301, v1008, int32(1), int32(20), v289, v1534, v291, v271+int32(60))
	mBase = m.M
	if v1537 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L214
	}
L200:
	;
	if base.Ui32(int32(4)) <= base.Ui32(v1010) {
		goto L202
	} else {
		goto L203
	}
L201:
	;
	v1353 = v1321 << (uint(int32(2)) % 32)
	v1362 = v293 + v1353
	v1373 = v271 + int32(2112) + v1353
	v1374 = v1321 + int32(1)
	goto L211
L202:
	;
	v1167 = v1008 << (uint(int32(2)) % 32)
	if base.Ui32(v295+v1167) <= base.Ui32(v303) {
		goto L204
	} else {
		goto L205
	}
L203:
	;
	v1321 = v1010
	goto L201
L204:
	;
	v1179 = v1010 & int32(2147483644)
	v1184 = v294 + v1167
	v1195 = v302 + v1167
	v1196 = v1179
	goto L207
L205:
	;
	if base.Ui32(v271+int32(2112)+v1167) <= base.Ui32(v295) {
		goto L204
	} else {
		goto L206
	}
L206:
	;
	v1321 = v1010
	goto L201
L207:
	;
	v1253 = int32(0)
	v1254 = base.Simd_g_v128_load_rng(m, v1184+int32(4), v1253, int32(-4), int32(20))
	v1258 = base.Simd_g_v128_load_nc(m, v1184, v1253)
	v1259 = base.Simd_g_const(&F_EncodeStreamHook__k1)
	v1265 = base.Simd_g_const(&F_EncodeStreamHook__k2)
	v1269 = base.Simd_g_v128_or(base.Simd_g_v128_and(base.Simd_g_i32x4_sub(base.Simd_g_v128_or(v1254, base.Simd_g_const(&F_EncodeStreamHook__k3)), base.Simd_g_v128_and(v1258, v1259)), v1259), base.Simd_g_v128_and(base.Simd_g_i32x4_sub(base.Simd_g_v128_or(v1254, base.Simd_g_const(&F_EncodeStreamHook__k4)), base.Simd_g_v128_and(v1258, v1265)), v1265))
	base.Simd_g_v128_store(m, v1195, v1253, v1269)
	v1272 = int32(-16)
	v1277 = v1196 + int32(-4)
	if v1277 != 0 {
		v1184 = v1184 + v1272
		v1195 = v1195 + v1272
		v1196 = v1277
		goto L207
	} else {
		goto L209
	}
L208:
	;
	if v1010 == v1179 {
		goto L199
	} else {
		goto L210
	}
L209:
	;
	goto L208
L210:
	;
	v1321 = v1010 & int32(3)
	goto L201
L211:
	;
	v1431 = *(*int32)(unsafe.Add(mBase, uint32(v1362+int32(4))))
	v1434 = *(*int32)(unsafe.Add(mBase, uint32(v1362)))
	v1435 = int32(-16711936)
	v1442 = int32(16711935)
	*(*int32)(unsafe.Add(mBase, uint32(v1373))) = (v1431|int32(16711680)-v1434&v1435)&v1435 | (v1431|int32(_a_F_EncodeStreamHook_1)-v1434&v1442)&v1442
	v1449 = int32(-4)
	v1454 = v1374 + int32(-1)
	if base.Ui32(int32(1)) < base.Ui32(v1454) {
		v1362 = v1362 + v1449
		v1373 = v1373 + v1449
		v1374 = v1454
		goto L211
	} else {
		goto L213
	}
L212:
	;
	goto L199
L213:
	;
	goto L212
L214:
	;
	v1540 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v1541 = *(*int32)(unsafe.Add(mBase, uint32(v1540)+8))
	v1543 = int32(3)
	v1545 = *(*int32)(unsafe.Add(mBase, uint32(v277)+68))
	if v1545 < v1543 {
		goto L215
	} else {
		goto L216
	}
L215:
	;
	v1548 = v1543
	goto L217
L216:
	;
	v1548 = int32(2)
	goto L217
L217:
	;
	if v1545 < int32(5) {
		goto L218
	} else {
		goto L219
	}
L218:
	;
	v1553 = v1548
	goto L220
L219:
	;
	v1553 = base.B2i32(v1545 < int32(17))
	goto L220
L220:
	;
	v1558 = int32(base.Ui32(v1541+int32(1)<<(uint(v1553)%32)+int32(-1)) >> (uint(v1553) % 32))
	v1559 = *(*int32)(unsafe.Add(mBase, uint32(v1540)+12))
	v1566 = base.I64_extend_i32_s(v1559) * base.I64_extend_i32_s(v1558)
	v1567 = *(*int32)(unsafe.Add(mBase, uint32(v277)+60))
	if v1567 == int32(0) {
		goto L224
	} else {
		goto L225
	}
L221:
	;
	if v1646 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L233
	}
L222:
	;
	v1606 = v1566 + v1601 + v1602 + int64(16)
	v1607 = *(*int32)(unsafe.Add(mBase, uint32(v277)+24))
	if v1607 == int32(0) {
		goto L228
	} else {
		goto L229
	}
L223:
	;
	v1589 = int32(3)
	v1591 = int32(2)
	v1600 = v1586
	v1601 = v1587
	v1602 = base.I64_extend_i32_u(int32(base.Ui32(v1559+v1589)>>(uint(v1591)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v1558+v1589)>>(uint(v1591)%32)))
	goto L222
L224:
	;
	v1580 = int32(0)
	v1581 = int64(0)
	v1583 = *(*int32)(unsafe.Add(mBase, uint32(v277)+52))
	if v1583 == v1580 {
		v1600 = v1580
		v1601 = v1581
		v1602 = v1581
		goto L222
	} else {
		goto L226
	}
L225:
	;
	v1571 = v1558 << (uint(int32(1)) % 32)
	v1574 = int32(2)
	v1578 = v1571 + int32(base.Ui32(v1571+int32(3))>>(uint(v1574)%32)) + v1574
	v1586 = v1578
	v1587 = base.I64_extend_i32_u(v1578)
	goto L223
L226:
	;
	v1586 = v1580
	v1587 = v1581
	goto L223
L227:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+32)) = v1558
	*(*int32)(unsafe.Add(mBase, uint32(v277)+8)) = v1624
	v1628 = int32(2)
	v1631 = int32(31)
	v1633 = int32(-32)
	v1634 = (v1624 + base.I32_wrap_i64(v1566)<<(uint(v1628)%32) + v1631) & v1633
	*(*int32)(unsafe.Add(mBase, uint32(v277)+16)) = v1634
	*(*int32)(unsafe.Add(mBase, uint32(v277)+20)) = (v1634 + v1600<<(uint(v1628)%32) + v1631) & v1633
	v1646 = int32(1)
	goto L221
L228:
	;
	F_WebPSafeFree(m, v1607)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v277)+24)) = int64(0)
	v1616 = F_WebPSafeMalloc(m, v1606, int32(4))
	mBase = m.M
	if v1616 != 0 {
		goto L231
	} else {
		goto L232
	}
L229:
	;
	v1610 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v277)+28)))
	if base.Ui64(v1606) <= base.Ui64(v1610) {
		v1624 = v1607
		goto L227
	} else {
		goto L230
	}
L230:
	;
	goto L228
L231:
	;
	*(*uint32)(unsafe.Add(mBase, uint32(v277)+28)) = uint32(v1606)
	*(*int32)(unsafe.Add(mBase, uint32(v277)+24)) = v1616
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(0)
	v1624 = v1616
	goto L227
L232:
	;
	v1617 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v1619 = F_WebPEncodingSetError(m, v1617, int32(1))
	mBase = m.M
	v1646 = v1619
	goto L221
L233:
	;
	v1649 = *(*int32)(unsafe.Add(mBase, uint32(v277)+32))
	v1650 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v1651 = *(*int32)(unsafe.Add(mBase, uint32(v1540)+56))
	v1652 = *(*int32)(unsafe.Add(mBase, uint32(v1540)+52))
	v1653 = base.I64_extend_i32_s(v1541)
	v1654 = int32(1)
	if v1653 == int64(0) {
		goto L236
	} else {
		goto L237
	}
L234:
	;
	if v1675 == int32(0) {
		goto L62
	} else {
		goto L240
	}
L235:
	;
	goto L234
L236:
	;
	v1673 = F_malloc(m, base.I32_wrap_i64(v1653)*v1654)
	mBase = m.M
	v1675 = v1673
	goto L235
L237:
	;
	v1661 = base.I64_div_u_s(int64(2147418112), v1653)
	v1662 = int32(0)
	v1663 = base.I64_extend_i32_u(v1654)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1663*v1653) {
		v1675 = v1662
		goto L235
	} else {
		goto L238
	}
L238:
	;
	if base.Ui64(v1661) < base.Ui64(v1663) {
		v1675 = v1662
		goto L235
	} else {
		goto L239
	}
L239:
	;
	goto L236
L240:
	;
	if int32(3) < v1545 {
		goto L242
	} else {
		goto L243
	}
L241:
	;
	F_free(m, v1675)
	mBase = m.M
	goto L448
L242:
	;
	goto L265
L243:
	;
	if v1559 < int32(1) {
		goto L241
	} else {
		goto L244
	}
L244:
	;
	v1685 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v1686 = int32(0)
	v1728 = v1686
	v1730 = v1685
	v1739 = v1650
	v1740 = v1652
	v1742 = v1686
	goto L245
L245:
	;
	if v1541 < int32(1) {
		v1895 = v1728
		v1897 = v1730
		goto L247
	} else {
		goto L248
	}
L247:
	;
	v1926 = m.G101
	v1927 = *(*int32)(unsafe.Add(mBase, uint32(v1926)))
	m.T0[v1927].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v1739)
	mBase = m.M
	v1934 = v1742 + int32(1)
	if v1934 != v1559 {
		v1728 = v1895
		v1730 = v1897
		v1739 = v1739 + v1649<<(uint(int32(2))%32)
		v1740 = v1740 + v1651<<(uint(int32(2))%32)
		v1742 = v1934
		goto L245
	} else {
		goto L262
	}
L248:
	;
	v1765 = v1740
	v1776 = v1675
	v1801 = v1728
	v1802 = v1541
	v1803 = v1730
	goto L249
L249:
	;
	v1832 = *(*int32)(unsafe.Add(mBase, uint32(v1765)))
	if v1832 == v1803 {
		v1846 = v1801
		v1847 = v1803
		goto L251
	} else {
		goto L252
	}
L250:
	;
	v1895 = v1846
	v1897 = v1847
	goto L247
L251:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1776))) = uint8(v1846)
	v1854 = v1802 + int32(-1)
	if v1854 != 0 {
		v1765 = v1765 + int32(4)
		v1776 = v1776 + int32(1)
		v1801 = v1846
		v1802 = v1854
		v1803 = v1847
		goto L249
	} else {
		goto L261
	}
L252:
	;
	v1834 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	if v1834 != v1832 {
		goto L254
	} else {
		goto L255
	}
L253:
	;
	v1846 = v1845
	v1847 = v1832
	goto L251
L254:
	;
	v1837 = *(*int32)(unsafe.Add(mBase, uint32(v277)+76))
	if v1837 != v1832 {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	v1845 = int32(0)
	goto L253
L256:
	;
	v1842 = *(*int32)(unsafe.Add(mBase, uint32(v277)+80))
	if v1842 == v1832 {
		goto L258
	} else {
		goto L259
	}
L257:
	;
	v1845 = int32(1)
	goto L253
L258:
	;
	v1844 = int32(2)
	goto L260
L259:
	;
	v1844 = int32(3)
	goto L260
L260:
	;
	v1845 = v1844
	goto L253
L261:
	;
	goto L250
L262:
	;
	goto L241
L263:
	;
	v2065 = v295
	v2076 = int32(0)
	goto L277
L265:
	;
	base.MemoryFill(m, v271+int32(2112), int32(255), int32(_a_F_EncodeStreamHook_2))
	goto L263
L277:
	;
	v2134 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v2065)+1)))
	v2137 = v271 + int32(2112) + v2134<<(uint(int32(1))%32)
	v2138 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2137))))
	if v2138 == int32(_a_F_EncodeStreamHook_3) {
		goto L279
	} else {
		goto L280
	}
L278:
	;
	if v1559 < int32(1) {
		goto L241
	} else {
		goto L421
	}
L279:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2137))) = uint16(v2076)
	v4006 = v2076 + int32(1)
	if v1545 != v4006 {
		v2065 = v2065 + int32(4)
		v2076 = v4006
		goto L277
	} else {
		goto L420
	}
L280:
	;
	goto L283
L281:
	;
	v2270 = v295
	v2281 = int32(0)
	goto L297
L283:
	;
	base.MemoryFill(m, v271+int32(2112), int32(255), int32(_a_F_EncodeStreamHook_2))
	goto L281
L295:
	;
	v3542 = F_memcpy(m, v271+int32(64), v295, v1545<<(uint(int32(2))%32))
	mBase = m.M
	v3544 = m.G2
	F_qsort(m, v3542, v1545, int32(4), v3544+int32(328))
	mBase = m.M
	if v1545 == int32(0) {
		goto L374
	} else {
		goto L375
	}
L296:
	;
	v3229 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v3230 = int32(0)
	v3248 = v3229
	v3273 = v3230
	v3283 = v1650
	v3284 = v1652
	v3285 = v3230
	goto L356
L297:
	;
	v2339 = *(*int32)(unsafe.Add(mBase, uint32(v2270)))
	v2348 = v271 + int32(2112) + int32(base.Ui32(v2339&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)
	v2349 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2348))))
	if v2349 == int32(_a_F_EncodeStreamHook_3) {
		goto L299
	} else {
		goto L300
	}
L298:
	;
	if v1559 < int32(1) {
		goto L241
	} else {
		goto L329
	}
L299:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2348))) = uint16(v2281)
	v2744 = v2281 + int32(1)
	if v1545 != v2744 {
		v2270 = v2270 + int32(4)
		v2281 = v2744
		goto L297
	} else {
		goto L328
	}
L300:
	;
	goto L303
L301:
	;
	v2481 = v295
	v2492 = int32(0)
	goto L315
L303:
	;
	base.MemoryFill(m, v271+int32(2112), int32(255), int32(_a_F_EncodeStreamHook_2))
	goto L301
L315:
	;
	v2550 = *(*int32)(unsafe.Add(mBase, uint32(v2481)))
	v2559 = v271 + int32(2112) + int32(base.Ui32(v2550&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)
	v2560 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2559))))
	if v2560 != int32(_a_F_EncodeStreamHook_3) {
		goto L295
	} else {
		goto L317
	}
L316:
	;
	if v1559 < int32(1) {
		goto L241
	} else {
		goto L319
	}
L317:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2559))) = uint16(v2492)
	v2567 = v2492 + int32(1)
	if v1545 != v2567 {
		v2481 = v2481 + int32(4)
		v2492 = v2567
		goto L315
	} else {
		goto L318
	}
L318:
	;
	goto L316
L319:
	;
	if int32(0) < v1541 {
		goto L296
	} else {
		goto L320
	}
L320:
	;
	v2573 = int32(1)
	if v1559 == v2573 {
		v2714 = v1650
		goto L321
	} else {
		goto L322
	}
L321:
	;
	if v1559&v2573 == int32(0) {
		goto L241
	} else {
		goto L326
	}
L322:
	;
	v2587 = v1559 & int32(2147483646)
	v2634 = v1650
	goto L323
L323:
	;
	v2654 = m.G101
	v2655 = *(*int32)(unsafe.Add(mBase, uint32(v2654)))
	m.T0[v2655].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2634)
	mBase = m.M
	v2658 = *(*int32)(unsafe.Add(mBase, uint32(v2654)))
	m.T0[v2658].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2634+v1649<<(uint(int32(2))%32))
	mBase = m.M
	v2660 = v2634 + v1649<<(uint(int32(3))%32)
	v2662 = v2587 + int32(-2)
	if v2662 != 0 {
		v2587 = v2662
		v2634 = v2660
		goto L323
	} else {
		goto L325
	}
L324:
	;
	v2714 = v2660
	goto L321
L325:
	;
	goto L324
L326:
	;
	v2736 = m.G101
	v2737 = *(*int32)(unsafe.Add(mBase, uint32(v2736)))
	m.T0[v2737].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2714)
	mBase = m.M
	F_free(m, v1675)
	mBase = m.M
	goto L327
L327:
	;
	goto L61
L328:
	;
	goto L298
L329:
	;
	if int32(0) < v1541 {
		goto L330
	} else {
		goto L331
	}
L330:
	;
	v2923 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v2924 = int32(0)
	v2942 = v2923
	v2967 = v2924
	v2977 = v1650
	v2978 = v1652
	v2979 = v2924
	goto L339
L331:
	;
	v2750 = int32(1)
	if v1559 == v2750 {
		v2891 = v1650
		goto L332
	} else {
		goto L333
	}
L332:
	;
	if v1559&v2750 == int32(0) {
		goto L241
	} else {
		goto L337
	}
L333:
	;
	v2764 = v1559 & int32(2147483646)
	v2811 = v1650
	goto L334
L334:
	;
	v2831 = m.G101
	v2832 = *(*int32)(unsafe.Add(mBase, uint32(v2831)))
	m.T0[v2832].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2811)
	mBase = m.M
	v2835 = *(*int32)(unsafe.Add(mBase, uint32(v2831)))
	m.T0[v2835].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2811+v1649<<(uint(int32(2))%32))
	mBase = m.M
	v2837 = v2811 + v1649<<(uint(int32(3))%32)
	v2839 = v2764 + int32(-2)
	if v2839 != 0 {
		v2764 = v2839
		v2811 = v2837
		goto L334
	} else {
		goto L336
	}
L335:
	;
	v2891 = v2837
	goto L332
L336:
	;
	goto L335
L337:
	;
	v2913 = m.G101
	v2914 = *(*int32)(unsafe.Add(mBase, uint32(v2913)))
	m.T0[v2914].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2891)
	mBase = m.M
	F_free(m, v1675)
	mBase = m.M
	goto L338
L338:
	;
	goto L61
L339:
	;
	if v1541 != int32(1) {
		goto L342
	} else {
		goto L343
	}
L341:
	;
	if v1541&int32(1) == int32(0) {
		v3211 = v3132
		v3212 = v3157
		goto L351
	} else {
		goto L352
	}
L342:
	;
	v3005 = v2978
	v3016 = int32(0)
	v3017 = v2942
	v3042 = v2967
	goto L344
L343:
	;
	v3131 = int32(0)
	v3132 = v2942
	v3157 = v2967
	goto L341
L344:
	;
	v3072 = *(*int32)(unsafe.Add(mBase, uint32(v3005)))
	if v3072 == v3017 {
		v3086 = v3017
		v3087 = v3042
		goto L346
	} else {
		goto L347
	}
L345:
	;
	v3131 = v3114
	v3132 = v3106
	v3157 = v3107
	goto L341
L346:
	;
	v3088 = v1675 + v3016
	*(*uint8)(unsafe.Add(mBase, uint32(v3088))) = uint8(v3087)
	v3092 = *(*int32)(unsafe.Add(mBase, uint32(v3005+int32(4))))
	if v3092 == v3086 {
		v3106 = v3086
		v3107 = v3087
		goto L348
	} else {
		goto L349
	}
L347:
	;
	v3085 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v3072&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)))))
	v3086 = v3072
	v3087 = v3085
	goto L346
L348:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3088+int32(1)))) = uint8(v3107)
	v3114 = v3016 + int32(2)
	if v1541&int32(2147483646) != v3114 {
		v3005 = v3005 + int32(8)
		v3016 = v3114
		v3017 = v3106
		v3042 = v3107
		goto L344
	} else {
		goto L350
	}
L349:
	;
	v3105 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v3092&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)))))
	v3106 = v3092
	v3107 = v3105
	goto L348
L350:
	;
	goto L345
L351:
	;
	v3213 = m.G101
	v3214 = *(*int32)(unsafe.Add(mBase, uint32(v3213)))
	m.T0[v3214].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v2977)
	mBase = m.M
	v3221 = v2979 + int32(1)
	if v3221 != v1559 {
		v2942 = v3211
		v2967 = v3212
		v2977 = v2977 + v1649<<(uint(int32(2))%32)
		v2978 = v2978 + v1651<<(uint(int32(2))%32)
		v2979 = v3221
		goto L339
	} else {
		goto L355
	}
L352:
	;
	v3192 = *(*int32)(unsafe.Add(mBase, uint32(v2978+v3131<<(uint(int32(2))%32))))
	if v3192 == v3132 {
		v3206 = v3132
		v3207 = v3157
		goto L353
	} else {
		goto L354
	}
L353:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1675+v3131))) = uint8(v3207)
	v3211 = v3206
	v3212 = v3207
	goto L351
L354:
	;
	v3205 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v3192&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)))))
	v3206 = v3192
	v3207 = v3205
	goto L353
L355:
	;
	goto L241
L356:
	;
	if v1541 != int32(1) {
		goto L359
	} else {
		goto L360
	}
L358:
	;
	if v1541&int32(1) == int32(0) {
		v3517 = v3438
		v3518 = v3463
		goto L368
	} else {
		goto L369
	}
L359:
	;
	v3311 = v3284
	v3322 = int32(0)
	v3323 = v3248
	v3348 = v3273
	goto L361
L360:
	;
	v3437 = int32(0)
	v3438 = v3248
	v3463 = v3273
	goto L358
L361:
	;
	v3378 = *(*int32)(unsafe.Add(mBase, uint32(v3311)))
	if v3378 == v3323 {
		v3392 = v3323
		v3393 = v3348
		goto L363
	} else {
		goto L364
	}
L362:
	;
	v3437 = v3420
	v3438 = v3412
	v3463 = v3413
	goto L358
L363:
	;
	v3394 = v1675 + v3322
	*(*uint8)(unsafe.Add(mBase, uint32(v3394))) = uint8(v3393)
	v3398 = *(*int32)(unsafe.Add(mBase, uint32(v3311+int32(4))))
	if v3398 == v3392 {
		v3412 = v3392
		v3413 = v3393
		goto L365
	} else {
		goto L366
	}
L364:
	;
	v3391 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v3378&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)))))
	v3392 = v3378
	v3393 = v3391
	goto L363
L365:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3394+int32(1)))) = uint8(v3413)
	v3420 = v3322 + int32(2)
	if v1541&int32(2147483646) != v3420 {
		v3311 = v3311 + int32(8)
		v3322 = v3420
		v3323 = v3412
		v3348 = v3413
		goto L361
	} else {
		goto L367
	}
L366:
	;
	v3411 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v3398&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)))))
	v3412 = v3398
	v3413 = v3411
	goto L365
L367:
	;
	goto L362
L368:
	;
	v3519 = m.G101
	v3520 = *(*int32)(unsafe.Add(mBase, uint32(v3519)))
	m.T0[v3520].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v3283)
	mBase = m.M
	v3527 = v3285 + int32(1)
	if v3527 != v1559 {
		v3248 = v3517
		v3273 = v3518
		v3283 = v3283 + v1649<<(uint(int32(2))%32)
		v3284 = v3284 + v1651<<(uint(int32(2))%32)
		v3285 = v3527
		goto L356
	} else {
		goto L372
	}
L369:
	;
	v3498 = *(*int32)(unsafe.Add(mBase, uint32(v3284+v3437<<(uint(int32(2))%32))))
	if v3498 == v3438 {
		v3512 = v3438
		v3513 = v3463
		goto L370
	} else {
		goto L371
	}
L370:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1675+v3437))) = uint8(v3513)
	v3517 = v3512
	v3518 = v3513
	goto L368
L371:
	;
	v3511 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v3498&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)))))
	v3512 = v3498
	v3513 = v3511
	goto L370
L372:
	;
	goto L241
L373:
	;
	if v1559 < int32(1) {
		goto L241
	} else {
		goto L390
	}
L374:
	;
	goto L373
L375:
	;
	v3556 = int32(0)
	goto L376
L376:
	;
	v3563 = *(*int32)(unsafe.Add(mBase, uint32(v3542)))
	v3567 = *(*int32)(unsafe.Add(mBase, uint32(v295+v3556<<(uint(int32(2))%32))))
	if v3563 == v3567 {
		v3594 = int32(0)
		goto L378
	} else {
		goto L379
	}
L377:
	;
	goto L374
L378:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v271+int32(1088)+v3594<<(uint(int32(2))%32)))) = v3556
	v3608 = v3556 + int32(1)
	if v3608 != v1545 {
		v3556 = v3608
		goto L376
	} else {
		goto L389
	}
L379:
	;
	v3577 = int32(0)
	v3578 = v1545
	goto L380
L380:
	;
	v3583 = (v3578 + v3577) >> (uint(int32(1)) % 32)
	v3587 = *(*int32)(unsafe.Add(mBase, uint32(v3542+v3583<<(uint(int32(2))%32))))
	v3588 = base.B2i32(base.Ui32(v3587) < base.Ui32(v3567))
	if base.Ui32(v3587) < base.Ui32(v3567) {
		goto L382
	} else {
		goto L383
	}
L381:
	;
	v3594 = v3583
	goto L378
L382:
	;
	v3589 = v3578
	goto L384
L383:
	;
	v3589 = v3583
	goto L384
L384:
	;
	if base.Ui32(v3587) < base.Ui32(v3567) {
		goto L385
	} else {
		goto L386
	}
L385:
	;
	v3590 = v3583
	goto L387
L386:
	;
	v3590 = v3577
	goto L387
L387:
	;
	if v3587 != v3567 {
		v3577 = v3590
		v3578 = v3589
		goto L380
	} else {
		goto L388
	}
L388:
	;
	goto L381
L389:
	;
	goto L377
L390:
	;
	if int32(0) < v1541 {
		goto L391
	} else {
		goto L392
	}
L391:
	;
	v3794 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v3795 = int32(0)
	v3838 = v3794
	v3839 = v3795
	v3848 = v1650
	v3849 = v1652
	v3851 = v3795
	goto L400
L392:
	;
	v3625 = int32(1)
	if v1559 == v3625 {
		v3766 = v1650
		goto L393
	} else {
		goto L394
	}
L393:
	;
	if v1559&v3625 == int32(0) {
		goto L241
	} else {
		goto L398
	}
L394:
	;
	v3639 = v1559 & int32(2147483646)
	v3686 = v1650
	goto L395
L395:
	;
	v3706 = m.G101
	v3707 = *(*int32)(unsafe.Add(mBase, uint32(v3706)))
	m.T0[v3707].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v3686)
	mBase = m.M
	v3710 = *(*int32)(unsafe.Add(mBase, uint32(v3706)))
	m.T0[v3710].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v3686+v1649<<(uint(int32(2))%32))
	mBase = m.M
	v3712 = v3686 + v1649<<(uint(int32(3))%32)
	v3714 = v3639 + int32(-2)
	if v3714 != 0 {
		v3639 = v3714
		v3686 = v3712
		goto L395
	} else {
		goto L397
	}
L396:
	;
	v3766 = v3712
	goto L393
L397:
	;
	goto L396
L398:
	;
	v3788 = m.G101
	v3789 = *(*int32)(unsafe.Add(mBase, uint32(v3788)))
	m.T0[v3789].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v3766)
	mBase = m.M
	F_free(m, v1675)
	mBase = m.M
	goto L399
L399:
	;
	goto L61
L400:
	;
	v3873 = v3849
	v3884 = int32(0)
	v3910 = v3838
	v3911 = v3839
	goto L402
L402:
	;
	v3940 = *(*int32)(unsafe.Add(mBase, uint32(v3873)))
	if v3940 == v3910 {
		v3983 = v3910
		v3984 = v3911
		goto L404
	} else {
		goto L405
	}
L403:
	;
	v3992 = m.G101
	v3993 = *(*int32)(unsafe.Add(mBase, uint32(v3992)))
	m.T0[v3993].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v3848)
	mBase = m.M
	v4000 = v3851 + int32(1)
	if v4000 != v1559 {
		v3838 = v3983
		v3839 = v3984
		v3848 = v3848 + v1649<<(uint(int32(2))%32)
		v3849 = v3849 + v1651<<(uint(int32(2))%32)
		v3851 = v4000
		goto L400
	} else {
		goto L419
	}
L404:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1675+v3884))) = uint8(v3984)
	v3990 = v3884 + int32(1)
	if v1541 != v3990 {
		v3873 = v3873 + int32(4)
		v3884 = v3990
		v3910 = v3983
		v3911 = v3984
		goto L402
	} else {
		goto L418
	}
L405:
	;
	v3945 = v271 + int32(64)
	v3951 = *(*int32)(unsafe.Add(mBase, uint32(v3945)))
	if v3951 == v3940 {
		v3975 = int32(0)
		goto L407
	} else {
		goto L408
	}
L406:
	;
	v3982 = *(*int32)(unsafe.Add(mBase, uint32(v271+int32(1088)+v3975<<(uint(int32(2))%32))))
	v3983 = v3940
	v3984 = v3982
	goto L404
L407:
	;
	goto L406
L408:
	;
	v3956 = v1545
	v3958 = int32(0)
	goto L409
L409:
	;
	v3963 = (v3956 + v3958) >> (uint(int32(1)) % 32)
	v3967 = *(*int32)(unsafe.Add(mBase, uint32(v3945+v3963<<(uint(int32(2))%32))))
	v3968 = base.B2i32(base.Ui32(v3967) < base.Ui32(v3940))
	if base.Ui32(v3967) < base.Ui32(v3940) {
		goto L411
	} else {
		goto L412
	}
L410:
	;
	v3975 = v3963
	goto L407
L411:
	;
	v3969 = v3956
	goto L413
L412:
	;
	v3969 = v3963
	goto L413
L413:
	;
	if base.Ui32(v3967) < base.Ui32(v3940) {
		goto L414
	} else {
		goto L415
	}
L414:
	;
	v3970 = v3963
	goto L416
L415:
	;
	v3970 = v3958
	goto L416
L416:
	;
	if v3967 != v3940 {
		v3956 = v3969
		v3958 = v3970
		goto L409
	} else {
		goto L417
	}
L417:
	;
	goto L410
L418:
	;
	goto L403
L419:
	;
	goto L241
L420:
	;
	goto L278
L421:
	;
	if int32(0) < v1541 {
		goto L422
	} else {
		goto L423
	}
L422:
	;
	v4185 = *(*int32)(unsafe.Add(mBase, uint32(v295)))
	v4186 = int32(0)
	v4204 = v4186
	v4229 = v4185
	v4239 = v1650
	v4240 = v1652
	v4241 = v4186
	goto L431
L423:
	;
	v4012 = int32(1)
	if v1559 == v4012 {
		v4153 = v1650
		goto L424
	} else {
		goto L425
	}
L424:
	;
	if v1559&v4012 == int32(0) {
		goto L241
	} else {
		goto L429
	}
L425:
	;
	v4026 = v1559 & int32(2147483646)
	v4073 = v1650
	goto L426
L426:
	;
	v4093 = m.G101
	v4094 = *(*int32)(unsafe.Add(mBase, uint32(v4093)))
	m.T0[v4094].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v4073)
	mBase = m.M
	v4097 = *(*int32)(unsafe.Add(mBase, uint32(v4093)))
	m.T0[v4097].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v4073+v1649<<(uint(int32(2))%32))
	mBase = m.M
	v4099 = v4073 + v1649<<(uint(int32(3))%32)
	v4101 = v4026 + int32(-2)
	if v4101 != 0 {
		v4026 = v4101
		v4073 = v4099
		goto L426
	} else {
		goto L428
	}
L427:
	;
	v4153 = v4099
	goto L424
L428:
	;
	goto L427
L429:
	;
	v4175 = m.G101
	v4176 = *(*int32)(unsafe.Add(mBase, uint32(v4175)))
	m.T0[v4176].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v4153)
	mBase = m.M
	F_free(m, v1675)
	mBase = m.M
	goto L430
L430:
	;
	goto L61
L431:
	;
	if v1541 != int32(1) {
		goto L434
	} else {
		goto L435
	}
L432:
	;
	goto L241
L433:
	;
	if v1541&int32(1) == int32(0) {
		v4461 = v4386
		v4462 = v4411
		goto L443
	} else {
		goto L444
	}
L434:
	;
	v4267 = v4240
	v4278 = int32(0)
	v4279 = v4204
	v4304 = v4229
	goto L436
L435:
	;
	v4385 = int32(0)
	v4386 = v4204
	v4411 = v4229
	goto L433
L436:
	;
	v4334 = *(*int32)(unsafe.Add(mBase, uint32(v4267)))
	if v4334 == v4304 {
		v4344 = v4279
		v4345 = v4304
		goto L438
	} else {
		goto L439
	}
L437:
	;
	v4385 = v4368
	v4386 = v4360
	v4411 = v4361
	goto L433
L438:
	;
	v4346 = v1675 + v4278
	*(*uint8)(unsafe.Add(mBase, uint32(v4346))) = uint8(v4344)
	v4350 = *(*int32)(unsafe.Add(mBase, uint32(v4267+int32(4))))
	if v4350 == v4345 {
		v4360 = v4344
		v4361 = v4345
		goto L440
	} else {
		goto L441
	}
L439:
	;
	v4343 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v4334)>>(uint(int32(7))%32))&int32(510)))))
	v4344 = v4343
	v4345 = v4334
	goto L438
L440:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4346+int32(1)))) = uint8(v4360)
	v4368 = v4278 + int32(2)
	if v1541&int32(2147483646) != v4368 {
		v4267 = v4267 + int32(8)
		v4278 = v4368
		v4279 = v4360
		v4304 = v4361
		goto L436
	} else {
		goto L442
	}
L441:
	;
	v4359 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v4350)>>(uint(int32(7))%32))&int32(510)))))
	v4360 = v4359
	v4361 = v4350
	goto L440
L442:
	;
	goto L437
L443:
	;
	v4463 = m.G101
	v4464 = *(*int32)(unsafe.Add(mBase, uint32(v4463)))
	m.T0[v4464].(func(*base.Module, int32, int32, int32, int32))(m, v1675, v1541, v1553, v4239)
	mBase = m.M
	v4471 = v4241 + int32(1)
	if v4471 != v1559 {
		v4204 = v4461
		v4229 = v4462
		v4239 = v4239 + v1649<<(uint(int32(2))%32)
		v4240 = v4240 + v1651<<(uint(int32(2))%32)
		v4241 = v4471
		goto L431
	} else {
		goto L447
	}
L444:
	;
	v4446 = *(*int32)(unsafe.Add(mBase, uint32(v4240+v4385<<(uint(int32(2))%32))))
	if v4446 == v4411 {
		v4456 = v4386
		v4457 = v4411
		goto L445
	} else {
		goto L446
	}
L445:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1675+v4385))) = uint8(v4456)
	v4461 = v4456
	v4462 = v4457
	goto L443
L446:
	;
	v4455 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v271+int32(2112)+int32(base.Ui32(v4446)>>(uint(int32(7))%32))&int32(510)))))
	v4456 = v4455
	v4457 = v4446
	goto L445
L447:
	;
	goto L432
L448:
	;
	goto L61
L449:
	;
	goto L63
L450:
	;
	if v4638 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L462
	}
L451:
	;
	v4598 = v4558 + v4593 + v4594 + int64(16)
	v4599 = *(*int32)(unsafe.Add(mBase, uint32(v277)+24))
	if v4599 == int32(0) {
		goto L457
	} else {
		goto L458
	}
L452:
	;
	v4581 = int32(3)
	v4583 = int32(2)
	v4592 = v4578
	v4593 = v4579
	v4594 = base.I64_extend_i32_u(int32(base.Ui32(v4551+v4581)>>(uint(v4583)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v4550+v4581)>>(uint(v4583)%32)))
	goto L451
L453:
	;
	v4572 = int32(0)
	v4573 = int64(0)
	v4575 = *(*int32)(unsafe.Add(mBase, uint32(v277)+52))
	if v4575 == v4572 {
		v4592 = v4572
		v4593 = v4573
		v4594 = v4573
		goto L451
	} else {
		goto L455
	}
L454:
	;
	v4563 = v4550 << (uint(int32(1)) % 32)
	v4566 = int32(2)
	v4570 = v4563 + int32(base.Ui32(v4563+int32(3))>>(uint(v4566)%32)) + v4566
	v4578 = v4570
	v4579 = base.I64_extend_i32_u(v4570)
	goto L452
L455:
	;
	v4578 = v4572
	v4579 = v4573
	goto L452
L456:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+32)) = v4550
	*(*int32)(unsafe.Add(mBase, uint32(v277)+8)) = v4616
	v4620 = int32(2)
	v4623 = int32(31)
	v4625 = int32(-32)
	v4626 = (v4616 + base.I32_wrap_i64(v4558)<<(uint(v4620)%32) + v4623) & v4625
	*(*int32)(unsafe.Add(mBase, uint32(v277)+16)) = v4626
	*(*int32)(unsafe.Add(mBase, uint32(v277)+20)) = (v4626 + v4592<<(uint(v4620)%32) + v4623) & v4625
	v4638 = int32(1)
	goto L450
L457:
	;
	F_WebPSafeFree(m, v4599)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v277)+24)) = int64(0)
	v4608 = F_WebPSafeMalloc(m, v4598, int32(4))
	mBase = m.M
	if v4608 != 0 {
		goto L460
	} else {
		goto L461
	}
L458:
	;
	v4602 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v277)+28)))
	if base.Ui64(v4598) <= base.Ui64(v4602) {
		v4616 = v4599
		goto L456
	} else {
		goto L459
	}
L459:
	;
	goto L457
L460:
	;
	*(*uint32)(unsafe.Add(mBase, uint32(v277)+28)) = uint32(v4598)
	*(*int32)(unsafe.Add(mBase, uint32(v277)+24)) = v4608
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(0)
	v4616 = v4608
	goto L456
L461:
	;
	v4609 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v4611 = F_WebPEncodingSetError(m, v4609, int32(1))
	mBase = m.M
	v4638 = v4611
	goto L450
L462:
	;
	v4641 = *(*int32)(unsafe.Add(mBase, uint32(v277)+12))
	if v4641 == int32(1) {
		v5021 = v290
		goto L60
	} else {
		goto L463
	}
L463:
	;
	if v4551 < int32(1) {
		goto L464
	} else {
		goto L465
	}
L464:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+12)) = int32(1)
	v5021 = v290
	goto L60
L465:
	;
	v4646 = int32(1)
	v4649 = v4550 << (uint(int32(2)) % 32)
	v4650 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v4651 = *(*int32)(unsafe.Add(mBase, uint32(v4549)+52))
	if v4551 == v4646 {
		v4747 = v4651
		v4784 = v4650
		goto L466
	} else {
		goto L467
	}
L466:
	;
	if v4551&v4646 == int32(0) {
		goto L464
	} else {
		goto L471
	}
L467:
	;
	v4662 = v4651
	v4699 = v4650
	v4700 = v4551 & int32(2147483646)
	goto L468
L468:
	;
	v4729 = F_memcpy(m, v4699, v4662, v4649)
	mBase = m.M
	v4731 = *(*int32)(unsafe.Add(mBase, uint32(v4549)+56))
	v4732 = int32(2)
	v4734 = v4662 + v4731<<(uint(v4732)%32)
	v4735 = F_memcpy(m, v4729+v4649, v4734, v4649)
	mBase = m.M
	v4736 = *(*int32)(unsafe.Add(mBase, uint32(v4549)+56))
	v4739 = v4734 + v4736<<(uint(v4732)%32)
	v4740 = v4729 + v4550<<(uint(int32(3))%32)
	v4742 = v4700 + int32(-2)
	if v4742 != 0 {
		v4662 = v4739
		v4699 = v4740
		v4700 = v4742
		goto L468
	} else {
		goto L470
	}
L469:
	;
	v4747 = v4739
	v4784 = v4740
	goto L466
L470:
	;
	goto L469
L471:
	;
	v4816 = F_memcpy(m, v4784, v4747, v4649)
	mBase = m.M
	goto L464
L472:
	;
	v10723 = v269
	v10725 = v271
	goto L4
L473:
	;
	goto L472
L474:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1540)+92)) = int32(1)
	goto L473
L476:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v277)+48)) = base.I32_clz(v4969) ^ int32(31) + int32(1)
	v5021 = v292
	goto L60
L477:
	;
	v5363 = *(*int32)(unsafe.Add(mBase, uint32(v277)+60))
	if v5363 == int32(0) {
		v6614 = v5021
		goto L531
	} else {
		goto L532
	}
L478:
	;
	v5052 = *(*int32)(unsafe.Add(mBase, uint32(v277)+32))
	v5053 = int32(1)
	goto L481
L479:
	;
	v5205 = int32(2)
	goto L507
L480:
	;
	goto L479
L481:
	;
	v5066 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v5067 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v5067+v5053 < int32(32) {
		goto L484
	} else {
		goto L485
	}
L482:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v5168 + v5166
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v5165<<(uint(v5168)%32) | v5167
	goto L480
L483:
	;
	v5085 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v5086 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v5090 = v5082
	v5091 = v5083
	v5092 = v5086
	v5093 = v5085
	goto L487
L484:
	;
	if v5067 < int32(16) {
		v5165 = v5053
		v5166 = v5053
		v5167 = v5066
		v5168 = v5067
		goto L482
	} else {
		goto L486
	}
L485:
	;
	v5071 = int32(32)
	v5072 = v5071 - v5067
	v5080 = int32(base.Ui32(v5053) >> (uint(v5072) % 32))
	v5081 = v5053 - v5072
	v5082 = v5053<<(uint(v5067)%32) | v5066
	v5083 = v5071
	goto L483
L486:
	;
	v5080 = v5053
	v5081 = v5053
	v5082 = v5066
	v5083 = v5067
	goto L483
L487:
	;
	if base.Ui32(v5092+int32(2)) <= base.Ui32(v5093) {
		v5148 = v5092
		v5149 = v5093
		goto L489
	} else {
		goto L490
	}
L488:
	;
	v5165 = v5080
	v5166 = v5081
	v5167 = v5161
	v5168 = v5159
	goto L482
L489:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v5148))) = uint16(v5090)
	v5156 = v5148 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v5156
	v5159 = v5091 + int32(-16)
	v5161 = int32(base.Ui32(v5090) >> (uint(int32(16)) % 32))
	if int32(31) < v5091 {
		v5090 = v5161
		v5091 = v5159
		v5092 = v5156
		v5093 = v5149
		goto L487
	} else {
		goto L504
	}
L490:
	;
	v5102 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v5103 = v5093 - v5102
	v5106 = base.I64_extend_i32_s(v5103) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5106) {
		v5130 = v5102
		goto L492
	} else {
		goto L493
	}
L491:
	;
	if v5092 == v5102 {
		goto L502
	} else {
		goto L503
	}
L492:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v5130
	goto L479
L493:
	;
	v5109 = v5092 - v5102
	v5111 = v5106 + base.I64_extend_i32_u(v5109)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5111) {
		v5130 = v5102
		goto L492
	} else {
		goto L494
	}
L494:
	;
	v5114 = base.I32_wrap_i64(v5111)
	if v5093 == v5102 {
		goto L495
	} else {
		goto L496
	}
L495:
	;
	v5121 = int32(base.Ui32(v5103*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v5114) < base.Ui32(v5121) {
		goto L498
	} else {
		goto L499
	}
L496:
	;
	if base.Ui32(v5114) <= base.Ui32(v5103) {
		v5148 = v5092
		v5149 = v5093
		goto L489
	} else {
		goto L497
	}
L497:
	;
	goto L495
L498:
	;
	v5123 = v5121
	goto L500
L499:
	;
	v5123 = v5114
	goto L500
L500:
	;
	v5127 = v5123&int32(-1024) + int32(1024)
	v5128 = F_WebPSafeMalloc(m, int64(1), v5127)
	mBase = m.M
	if v5128 != 0 {
		goto L491
	} else {
		goto L501
	}
L501:
	;
	v5129 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v5130 = v5129
	goto L492
L502:
	;
	v5141 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v5141)
	mBase = m.M
	v5143 = v5128 + v5127
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v5143
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v5128
	v5148 = v5128 + v5109
	v5149 = v5143
	goto L489
L503:
	;
	v5139 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v5140 = F_memcpy(m, v5128, v5139, v5109)
	mBase = m.M
	goto L502
L504:
	;
	goto L488
L505:
	;
	v5357 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v5359 = m.G94
	v5360 = *(*int32)(unsafe.Add(mBase, uint32(v5359)))
	m.T0[v5360].(func(*base.Module, int32, int32))(m, v5357, v5052*v279)
	mBase = m.M
	goto L477
L506:
	;
	goto L505
L507:
	;
	v5218 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v5219 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v5219+v5205 < int32(32) {
		goto L510
	} else {
		goto L511
	}
L508:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v5320 + v5318
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v5317<<(uint(v5320)%32) | v5319
	goto L506
L509:
	;
	v5237 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v5238 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v5242 = v5234
	v5243 = v5235
	v5244 = v5238
	v5245 = v5237
	goto L513
L510:
	;
	if v5219 < int32(16) {
		v5317 = v5205
		v5318 = v5205
		v5319 = v5218
		v5320 = v5219
		goto L508
	} else {
		goto L512
	}
L511:
	;
	v5223 = int32(32)
	v5224 = v5223 - v5219
	v5232 = int32(base.Ui32(v5205) >> (uint(v5224) % 32))
	v5233 = v5205 - v5224
	v5234 = v5205<<(uint(v5219)%32) | v5218
	v5235 = v5223
	goto L509
L512:
	;
	v5232 = v5205
	v5233 = v5205
	v5234 = v5218
	v5235 = v5219
	goto L509
L513:
	;
	if base.Ui32(v5244+int32(2)) <= base.Ui32(v5245) {
		v5300 = v5244
		v5301 = v5245
		goto L515
	} else {
		goto L516
	}
L514:
	;
	v5317 = v5232
	v5318 = v5233
	v5319 = v5313
	v5320 = v5311
	goto L508
L515:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v5300))) = uint16(v5242)
	v5308 = v5300 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v5308
	v5311 = v5243 + int32(-16)
	v5313 = int32(base.Ui32(v5242) >> (uint(int32(16)) % 32))
	if int32(31) < v5243 {
		v5242 = v5313
		v5243 = v5311
		v5244 = v5308
		v5245 = v5301
		goto L513
	} else {
		goto L530
	}
L516:
	;
	v5254 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v5255 = v5245 - v5254
	v5258 = base.I64_extend_i32_s(v5255) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5258) {
		v5282 = v5254
		goto L518
	} else {
		goto L519
	}
L517:
	;
	if v5244 == v5254 {
		goto L528
	} else {
		goto L529
	}
L518:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v5282
	goto L505
L519:
	;
	v5261 = v5244 - v5254
	v5263 = v5258 + base.I64_extend_i32_u(v5261)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5263) {
		v5282 = v5254
		goto L518
	} else {
		goto L520
	}
L520:
	;
	v5266 = base.I32_wrap_i64(v5263)
	if v5245 == v5254 {
		goto L521
	} else {
		goto L522
	}
L521:
	;
	v5273 = int32(base.Ui32(v5255*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v5266) < base.Ui32(v5273) {
		goto L524
	} else {
		goto L525
	}
L522:
	;
	if base.Ui32(v5266) <= base.Ui32(v5255) {
		v5300 = v5244
		v5301 = v5245
		goto L515
	} else {
		goto L523
	}
L523:
	;
	goto L521
L524:
	;
	v5275 = v5273
	goto L526
L525:
	;
	v5275 = v5266
	goto L526
L526:
	;
	v5279 = v5275&int32(-1024) + int32(1024)
	v5280 = F_WebPSafeMalloc(m, int64(1), v5279)
	mBase = m.M
	if v5280 != 0 {
		goto L517
	} else {
		goto L527
	}
L527:
	;
	v5281 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v5282 = v5281
	goto L518
L528:
	;
	v5293 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v5293)
	mBase = m.M
	v5295 = v5280 + v5279
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v5295
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v5280
	v5300 = v5280 + v5261
	v5301 = v5295
	goto L515
L529:
	;
	v5291 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v5292 = F_memcpy(m, v5280, v5291, v5261)
	mBase = m.M
	goto L528
L530:
	;
	goto L514
L531:
	;
	v6642 = *(*int32)(unsafe.Add(mBase, uint32(v277)+52))
	if v6642 == int32(0) {
		v7140 = v6614
		goto L666
	} else {
		goto L667
	}
L532:
	;
	v5367 = base.I32_div_s(v5021, int32(3))
	v5369 = *(*int32)(unsafe.Add(mBase, uint32(v277)+64))
	if v5369 != 0 {
		v5372 = int32(100)
		goto L533
	} else {
		goto L534
	}
L533:
	;
	v5373 = *(*int32)(unsafe.Add(mBase, uint32(v277)+56))
	v5376 = *(*int32)(unsafe.Add(mBase, uint32(v277)+40))
	v5377 = int32(9)
	if base.Ui32(v5376) < base.Ui32(v5377) {
		goto L535
	} else {
		goto L536
	}
L534:
	;
	v5370 = *(*int32)(unsafe.Add(mBase, uint32(v277)))
	v5371 = *(*int32)(unsafe.Add(mBase, uint32(v5370)+92))
	v5372 = v5371
	goto L533
L535:
	;
	v5380 = v5376
	goto L537
L536:
	;
	v5380 = v5377
	goto L537
L537:
	;
	if v5376 < int32(2) {
		goto L538
	} else {
		goto L539
	}
L538:
	;
	v5383 = int32(2)
	goto L540
L539:
	;
	v5383 = v5380
	goto L540
L540:
	;
	v5384 = int32(1) << (uint(v5383) % 32)
	v5385 = *(*int32)(unsafe.Add(mBase, uint32(v277)+32))
	v5387 = v5385 + int32(-1)
	v5392 = int32(base.Ui32(v5384+v5387)>>(uint(v5383)%32)) * int32(base.Ui32(v5384+v297)>>(uint(v5383)%32))
	if base.Ui32(int32(8)) < base.Ui32(v5383) {
		v5485 = v5383
		v5522 = v5392
		goto L541
	} else {
		goto L542
	}
L541:
	;
	if base.Ui32(int32(3)) <= base.Ui32(v5485) {
		goto L549
	} else {
		goto L550
	}
L542:
	;
	if v5392 < int32(_a_F_EncodeStreamHook_4) {
		v5485 = v5383
		v5522 = v5392
		goto L541
	} else {
		goto L543
	}
L543:
	;
	v5412 = v5383
	goto L544
L544:
	;
	v5469 = int32(2) << (uint(v5412) % 32)
	v5472 = v5412 + int32(1)
	v5476 = int32(base.Ui32(v5469+v5387)>>(uint(v5472)%32)) * int32(base.Ui32(v5469+v297)>>(uint(v5472)%32))
	if base.Ui32(int32(7)) < base.Ui32(v5412) {
		v5485 = v5472
		v5522 = v5476
		goto L541
	} else {
		goto L546
	}
L545:
	;
	v5485 = v5472
	v5522 = v5476
	goto L541
L546:
	;
	if int32(_a_F_EncodeStreamHook_5) < v5476 {
		v5412 = v5472
		goto L544
	} else {
		goto L547
	}
L547:
	;
	goto L545
L548:
	;
	v5725 = int32(1)
	v5727 = *(*int32)(unsafe.Add(mBase, uint32(v277)))
	v5728 = *(*int32)(unsafe.Add(mBase, uint32(v5727)+8))
	if int32(4) < v5728 {
		goto L559
	} else {
		goto L560
	}
L549:
	;
	if v5522 == int32(1) {
		goto L551
	} else {
		goto L552
	}
L550:
	;
	v5696 = v5485
	goto L548
L551:
	;
	v5556 = int32(1)
	v5558 = v5485 + int32(-1)
	v5559 = v5556 << (uint(v5558) % 32)
	if int32(base.Ui32(v5559+v5387)>>(uint(v5558)%32))*int32(base.Ui32(v5559+v297)>>(uint(v5558)%32)) != v5556 {
		v5696 = v5485
		goto L548
	} else {
		goto L553
	}
L552:
	;
	v5696 = v5485
	goto L548
L553:
	;
	v5571 = v5485
	goto L554
L554:
	;
	v5639 = v5571 + int32(-1)
	if int32(3) <= v5639 {
		goto L556
	} else {
		goto L557
	}
L555:
	;
	v5696 = v5639
	goto L548
L556:
	;
	v5644 = v5571 + int32(-2)
	v5645 = int32(1)
	v5646 = v5645 << (uint(v5644) % 32)
	if int32(base.Ui32(v5646+v5387)>>(uint(v5644)%32))*int32(base.Ui32(v5646+v297)>>(uint(v5644)%32)) == v5645 {
		v5571 = v5639
		goto L554
	} else {
		goto L558
	}
L557:
	;
	v5696 = int32(2)
	goto L548
L558:
	;
	goto L555
L559:
	;
	v5736 = v5728<<(uint(v5725)%32) + int32(-8)
	goto L561
L560:
	;
	v5736 = int32(0)
	goto L561
L561:
	;
	v5737 = v5696 - v5736
	v5738 = int32(9)
	if base.Ui32(v5737) < base.Ui32(v5738) {
		goto L562
	} else {
		goto L563
	}
L562:
	;
	v5741 = v5737
	goto L564
L563:
	;
	v5741 = v5738
	goto L564
L564:
	;
	if v5737 < int32(2) {
		goto L565
	} else {
		goto L566
	}
L565:
	;
	v5744 = int32(2)
	goto L567
L566:
	;
	v5744 = v5741
	goto L567
L567:
	;
	v5745 = v5725 << (uint(v5744) % 32)
	v5750 = int32(base.Ui32(v5745+v5387)>>(uint(v5744)%32)) * int32(base.Ui32(v5745+v297)>>(uint(v5744)%32))
	if base.Ui32(int32(8)) < base.Ui32(v5744) {
		v5843 = v5744
		v5880 = v5750
		goto L568
	} else {
		goto L569
	}
L568:
	;
	if base.Ui32(int32(3)) <= base.Ui32(v5843) {
		goto L576
	} else {
		goto L577
	}
L569:
	;
	if v5750 < int32(_a_F_EncodeStreamHook_4) {
		v5843 = v5744
		v5880 = v5750
		goto L568
	} else {
		goto L570
	}
L570:
	;
	v5770 = v5744
	goto L571
L571:
	;
	v5827 = int32(2) << (uint(v5770) % 32)
	v5830 = v5770 + int32(1)
	v5834 = int32(base.Ui32(v5827+v5387)>>(uint(v5830)%32)) * int32(base.Ui32(v5827+v297)>>(uint(v5830)%32))
	if base.Ui32(int32(7)) < base.Ui32(v5770) {
		v5843 = v5830
		v5880 = v5834
		goto L568
	} else {
		goto L573
	}
L572:
	;
	v5843 = v5830
	v5880 = v5834
	goto L568
L573:
	;
	if int32(_a_F_EncodeStreamHook_5) < v5834 {
		v5770 = v5830
		goto L571
	} else {
		goto L574
	}
L574:
	;
	goto L572
L575:
	;
	v6083 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v6084 = *(*int32)(unsafe.Add(mBase, uint32(v277)+16))
	v6085 = *(*int32)(unsafe.Add(mBase, uint32(v277)+20))
	v6086 = *(*int32)(unsafe.Add(mBase, uint32(v5727)+96))
	v6087 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v6089 = base.I32_div_s(v5021, int32(6))
	v6094 = F_VP8LResidualImage(m, v5385, v279, v6053, v5696, v289, v6083, v6084, v6085, v5372, v6086, v5373, v6087, v6089, v271+int32(60), v271+int32(4))
	mBase = m.M
	if v6094 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L586
	}
L576:
	;
	if v5880 == int32(1) {
		goto L578
	} else {
		goto L579
	}
L577:
	;
	v6053 = v5843
	goto L575
L578:
	;
	v5914 = int32(1)
	v5916 = v5843 + int32(-1)
	v5917 = v5914 << (uint(v5916) % 32)
	if int32(base.Ui32(v5917+v5387)>>(uint(v5916)%32))*int32(base.Ui32(v5917+v297)>>(uint(v5916)%32)) != v5914 {
		v6053 = v5843
		goto L575
	} else {
		goto L580
	}
L579:
	;
	v6053 = v5843
	goto L575
L580:
	;
	v5929 = v5843
	goto L581
L581:
	;
	v5997 = v5929 + int32(-1)
	if int32(3) <= v5997 {
		goto L583
	} else {
		goto L584
	}
L582:
	;
	v6053 = v5997
	goto L575
L583:
	;
	v6002 = v5929 + int32(-2)
	v6003 = int32(1)
	v6004 = v6003 << (uint(v6002) % 32)
	if int32(base.Ui32(v6004+v5387)>>(uint(v6002)%32))*int32(base.Ui32(v6004+v297)>>(uint(v6002)%32)) == v6003 {
		v5929 = v5997
		goto L581
	} else {
		goto L585
	}
L584:
	;
	v6053 = int32(2)
	goto L575
L585:
	;
	goto L582
L586:
	;
	v6097 = int32(1)
	goto L589
L587:
	;
	v6249 = int32(0)
	v6250 = int32(2)
	goto L615
L588:
	;
	goto L587
L589:
	;
	v6110 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v6111 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v6111+v6097 < int32(32) {
		goto L592
	} else {
		goto L593
	}
L590:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v6212 + v6210
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v6209<<(uint(v6212)%32) | v6211
	goto L588
L591:
	;
	v6129 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v6130 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v6134 = v6126
	v6135 = v6127
	v6136 = v6130
	v6137 = v6129
	goto L595
L592:
	;
	if v6111 < int32(16) {
		v6209 = v6097
		v6210 = v6097
		v6211 = v6110
		v6212 = v6111
		goto L590
	} else {
		goto L594
	}
L593:
	;
	v6115 = int32(32)
	v6116 = v6115 - v6111
	v6124 = int32(base.Ui32(v6097) >> (uint(v6116) % 32))
	v6125 = v6097 - v6116
	v6126 = v6097<<(uint(v6111)%32) | v6110
	v6127 = v6115
	goto L591
L594:
	;
	v6124 = v6097
	v6125 = v6097
	v6126 = v6110
	v6127 = v6111
	goto L591
L595:
	;
	if base.Ui32(v6136+int32(2)) <= base.Ui32(v6137) {
		v6192 = v6136
		v6193 = v6137
		goto L597
	} else {
		goto L598
	}
L596:
	;
	v6209 = v6124
	v6210 = v6125
	v6211 = v6205
	v6212 = v6203
	goto L590
L597:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6192))) = uint16(v6134)
	v6200 = v6192 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6200
	v6203 = v6135 + int32(-16)
	v6205 = int32(base.Ui32(v6134) >> (uint(int32(16)) % 32))
	if int32(31) < v6135 {
		v6134 = v6205
		v6135 = v6203
		v6136 = v6200
		v6137 = v6193
		goto L595
	} else {
		goto L612
	}
L598:
	;
	v6146 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6147 = v6137 - v6146
	v6150 = base.I64_extend_i32_s(v6147) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6150) {
		v6174 = v6146
		goto L600
	} else {
		goto L601
	}
L599:
	;
	if v6136 == v6146 {
		goto L610
	} else {
		goto L611
	}
L600:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6174
	goto L587
L601:
	;
	v6153 = v6136 - v6146
	v6155 = v6150 + base.I64_extend_i32_u(v6153)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6155) {
		v6174 = v6146
		goto L600
	} else {
		goto L602
	}
L602:
	;
	v6158 = base.I32_wrap_i64(v6155)
	if v6137 == v6146 {
		goto L603
	} else {
		goto L604
	}
L603:
	;
	v6165 = int32(base.Ui32(v6147*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6158) < base.Ui32(v6165) {
		goto L606
	} else {
		goto L607
	}
L604:
	;
	if base.Ui32(v6158) <= base.Ui32(v6147) {
		v6192 = v6136
		v6193 = v6137
		goto L597
	} else {
		goto L605
	}
L605:
	;
	goto L603
L606:
	;
	v6167 = v6165
	goto L608
L607:
	;
	v6167 = v6158
	goto L608
L608:
	;
	v6171 = v6167&int32(-1024) + int32(1024)
	v6172 = F_WebPSafeMalloc(m, int64(1), v6171)
	mBase = m.M
	if v6172 != 0 {
		goto L599
	} else {
		goto L609
	}
L609:
	;
	v6173 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6174 = v6173
	goto L600
L610:
	;
	v6185 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v6185)
	mBase = m.M
	v6187 = v6172 + v6171
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v6187
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v6172
	v6192 = v6172 + v6153
	v6193 = v6187
	goto L597
L611:
	;
	v6183 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6184 = F_memcpy(m, v6172, v6183, v6153)
	mBase = m.M
	goto L610
L612:
	;
	goto L596
L613:
	;
	v6401 = *(*int32)(unsafe.Add(mBase, uint32(v271)+4))
	v6403 = v6401 + int32(-2)
	v6404 = int32(3)
	goto L641
L614:
	;
	goto L613
L615:
	;
	v6262 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v6263 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v6263+v6250 < int32(32) {
		goto L618
	} else {
		goto L619
	}
L616:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v6364 + v6362
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v6361<<(uint(v6364)%32) | v6363
	goto L614
L617:
	;
	v6281 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v6282 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v6286 = v6278
	v6287 = v6279
	v6288 = v6282
	v6289 = v6281
	goto L621
L618:
	;
	if v6263 < int32(16) {
		v6361 = v6249
		v6362 = v6250
		v6363 = v6262
		v6364 = v6263
		goto L616
	} else {
		goto L620
	}
L619:
	;
	v6267 = int32(32)
	v6268 = v6267 - v6263
	v6276 = int32(base.Ui32(v6249) >> (uint(v6268) % 32))
	v6277 = v6250 - v6268
	v6278 = v6249<<(uint(v6263)%32) | v6262
	v6279 = v6267
	goto L617
L620:
	;
	v6276 = v6249
	v6277 = v6250
	v6278 = v6262
	v6279 = v6263
	goto L617
L621:
	;
	if base.Ui32(v6288+int32(2)) <= base.Ui32(v6289) {
		v6344 = v6288
		v6345 = v6289
		goto L623
	} else {
		goto L624
	}
L622:
	;
	v6361 = v6276
	v6362 = v6277
	v6363 = v6357
	v6364 = v6355
	goto L616
L623:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6344))) = uint16(v6286)
	v6352 = v6344 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6352
	v6355 = v6287 + int32(-16)
	v6357 = int32(base.Ui32(v6286) >> (uint(int32(16)) % 32))
	if int32(31) < v6287 {
		v6286 = v6357
		v6287 = v6355
		v6288 = v6352
		v6289 = v6345
		goto L621
	} else {
		goto L638
	}
L624:
	;
	v6298 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6299 = v6289 - v6298
	v6302 = base.I64_extend_i32_s(v6299) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6302) {
		v6326 = v6298
		goto L626
	} else {
		goto L627
	}
L625:
	;
	if v6288 == v6298 {
		goto L636
	} else {
		goto L637
	}
L626:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6326
	goto L613
L627:
	;
	v6305 = v6288 - v6298
	v6307 = v6302 + base.I64_extend_i32_u(v6305)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6307) {
		v6326 = v6298
		goto L626
	} else {
		goto L628
	}
L628:
	;
	v6310 = base.I32_wrap_i64(v6307)
	if v6289 == v6298 {
		goto L629
	} else {
		goto L630
	}
L629:
	;
	v6317 = int32(base.Ui32(v6299*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6310) < base.Ui32(v6317) {
		goto L632
	} else {
		goto L633
	}
L630:
	;
	if base.Ui32(v6310) <= base.Ui32(v6299) {
		v6344 = v6288
		v6345 = v6289
		goto L623
	} else {
		goto L631
	}
L631:
	;
	goto L629
L632:
	;
	v6319 = v6317
	goto L634
L633:
	;
	v6319 = v6310
	goto L634
L634:
	;
	v6323 = v6319&int32(-1024) + int32(1024)
	v6324 = F_WebPSafeMalloc(m, int64(1), v6323)
	mBase = m.M
	if v6324 != 0 {
		goto L625
	} else {
		goto L635
	}
L635:
	;
	v6325 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6326 = v6325
	goto L626
L636:
	;
	v6337 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v6337)
	mBase = m.M
	v6339 = v6324 + v6323
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v6339
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v6324
	v6344 = v6324 + v6305
	v6345 = v6339
	goto L623
L637:
	;
	v6335 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6336 = F_memcpy(m, v6324, v6335, v6305)
	mBase = m.M
	goto L636
L638:
	;
	goto L622
L639:
	;
	v6555 = *(*int32)(unsafe.Add(mBase, uint32(v277)+20))
	v6557 = *(*int32)(unsafe.Add(mBase, uint32(v271)+4))
	v6558 = int32(1) << (uint(v6557) % 32)
	v6563 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v6567 = F_EncodeImageNoHuffman(m, v272, v6555, v298, v301, int32(base.Ui32(v6558+v5387)>>(uint(v6557)%32)), int32(base.Ui32(v6558+v297)>>(uint(v6557)%32)), v286, v289, v6563, v5367-v6089, v271+int32(60))
	mBase = m.M
	if v6567 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L665
	}
L640:
	;
	goto L639
L641:
	;
	v6416 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v6417 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v6417+v6404 < int32(32) {
		goto L644
	} else {
		goto L645
	}
L642:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v6518 + v6516
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v6515<<(uint(v6518)%32) | v6517
	goto L640
L643:
	;
	v6435 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v6436 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v6440 = v6432
	v6441 = v6433
	v6442 = v6436
	v6443 = v6435
	goto L647
L644:
	;
	if v6417 < int32(16) {
		v6515 = v6403
		v6516 = v6404
		v6517 = v6416
		v6518 = v6417
		goto L642
	} else {
		goto L646
	}
L645:
	;
	v6421 = int32(32)
	v6422 = v6421 - v6417
	v6430 = int32(base.Ui32(v6403) >> (uint(v6422) % 32))
	v6431 = v6404 - v6422
	v6432 = v6403<<(uint(v6417)%32) | v6416
	v6433 = v6421
	goto L643
L646:
	;
	v6430 = v6403
	v6431 = v6404
	v6432 = v6416
	v6433 = v6417
	goto L643
L647:
	;
	if base.Ui32(v6442+int32(2)) <= base.Ui32(v6443) {
		v6498 = v6442
		v6499 = v6443
		goto L649
	} else {
		goto L650
	}
L648:
	;
	v6515 = v6430
	v6516 = v6431
	v6517 = v6511
	v6518 = v6509
	goto L642
L649:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6498))) = uint16(v6440)
	v6506 = v6498 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6506
	v6509 = v6441 + int32(-16)
	v6511 = int32(base.Ui32(v6440) >> (uint(int32(16)) % 32))
	if int32(31) < v6441 {
		v6440 = v6511
		v6441 = v6509
		v6442 = v6506
		v6443 = v6499
		goto L647
	} else {
		goto L664
	}
L650:
	;
	v6452 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6453 = v6443 - v6452
	v6456 = base.I64_extend_i32_s(v6453) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6456) {
		v6480 = v6452
		goto L652
	} else {
		goto L653
	}
L651:
	;
	if v6442 == v6452 {
		goto L662
	} else {
		goto L663
	}
L652:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6480
	goto L639
L653:
	;
	v6459 = v6442 - v6452
	v6461 = v6456 + base.I64_extend_i32_u(v6459)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6461) {
		v6480 = v6452
		goto L652
	} else {
		goto L654
	}
L654:
	;
	v6464 = base.I32_wrap_i64(v6461)
	if v6443 == v6452 {
		goto L655
	} else {
		goto L656
	}
L655:
	;
	v6471 = int32(base.Ui32(v6453*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6464) < base.Ui32(v6471) {
		goto L658
	} else {
		goto L659
	}
L656:
	;
	if base.Ui32(v6464) <= base.Ui32(v6453) {
		v6498 = v6442
		v6499 = v6443
		goto L649
	} else {
		goto L657
	}
L657:
	;
	goto L655
L658:
	;
	v6473 = v6471
	goto L660
L659:
	;
	v6473 = v6464
	goto L660
L660:
	;
	v6477 = v6473&int32(-1024) + int32(1024)
	v6478 = F_WebPSafeMalloc(m, int64(1), v6477)
	mBase = m.M
	if v6478 != 0 {
		goto L651
	} else {
		goto L661
	}
L661:
	;
	v6479 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6480 = v6479
	goto L652
L662:
	;
	v6491 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v6491)
	mBase = m.M
	v6493 = v6478 + v6477
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v6493
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v6478
	v6498 = v6478 + v6459
	v6499 = v6493
	goto L649
L663:
	;
	v6489 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6490 = F_memcpy(m, v6478, v6489, v6459)
	mBase = m.M
	goto L662
L664:
	;
	goto L648
L665:
	;
	v6614 = v5021 - v5367
	goto L531
L666:
	;
	v7141 = int32(0)
	v7142 = int32(1)
	goto L750
L667:
	;
	v6646 = base.I32_div_s(v6614, int32(2))
	v6647 = *(*int32)(unsafe.Add(mBase, uint32(v277)+32))
	v6648 = *(*int32)(unsafe.Add(mBase, uint32(v277)+44))
	v6649 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v6650 = *(*int32)(unsafe.Add(mBase, uint32(v277)+20))
	v6651 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v6653 = base.I32_div_s(v6614, int32(4))
	v6656 = F_VP8LColorSpaceTransform(m, v6647, v279, v6648, v286, v6649, v6650, v6651, v6653, v271+int32(60), v271)
	mBase = m.M
	if v6656 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L668
	}
L668:
	;
	v6659 = int32(1)
	goto L671
L669:
	;
	v6811 = int32(1)
	v6812 = int32(2)
	goto L697
L670:
	;
	goto L669
L671:
	;
	v6672 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v6673 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v6673+v6659 < int32(32) {
		goto L674
	} else {
		goto L675
	}
L672:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v6774 + v6772
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v6771<<(uint(v6774)%32) | v6773
	goto L670
L673:
	;
	v6691 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v6692 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v6696 = v6688
	v6697 = v6689
	v6698 = v6692
	v6699 = v6691
	goto L677
L674:
	;
	if v6673 < int32(16) {
		v6771 = v6659
		v6772 = v6659
		v6773 = v6672
		v6774 = v6673
		goto L672
	} else {
		goto L676
	}
L675:
	;
	v6677 = int32(32)
	v6678 = v6677 - v6673
	v6686 = int32(base.Ui32(v6659) >> (uint(v6678) % 32))
	v6687 = v6659 - v6678
	v6688 = v6659<<(uint(v6673)%32) | v6672
	v6689 = v6677
	goto L673
L676:
	;
	v6686 = v6659
	v6687 = v6659
	v6688 = v6672
	v6689 = v6673
	goto L673
L677:
	;
	if base.Ui32(v6698+int32(2)) <= base.Ui32(v6699) {
		v6754 = v6698
		v6755 = v6699
		goto L679
	} else {
		goto L680
	}
L678:
	;
	v6771 = v6686
	v6772 = v6687
	v6773 = v6767
	v6774 = v6765
	goto L672
L679:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6754))) = uint16(v6696)
	v6762 = v6754 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6762
	v6765 = v6697 + int32(-16)
	v6767 = int32(base.Ui32(v6696) >> (uint(int32(16)) % 32))
	if int32(31) < v6697 {
		v6696 = v6767
		v6697 = v6765
		v6698 = v6762
		v6699 = v6755
		goto L677
	} else {
		goto L694
	}
L680:
	;
	v6708 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6709 = v6699 - v6708
	v6712 = base.I64_extend_i32_s(v6709) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6712) {
		v6736 = v6708
		goto L682
	} else {
		goto L683
	}
L681:
	;
	if v6698 == v6708 {
		goto L692
	} else {
		goto L693
	}
L682:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6736
	goto L669
L683:
	;
	v6715 = v6698 - v6708
	v6717 = v6712 + base.I64_extend_i32_u(v6715)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6717) {
		v6736 = v6708
		goto L682
	} else {
		goto L684
	}
L684:
	;
	v6720 = base.I32_wrap_i64(v6717)
	if v6699 == v6708 {
		goto L685
	} else {
		goto L686
	}
L685:
	;
	v6727 = int32(base.Ui32(v6709*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6720) < base.Ui32(v6727) {
		goto L688
	} else {
		goto L689
	}
L686:
	;
	if base.Ui32(v6720) <= base.Ui32(v6709) {
		v6754 = v6698
		v6755 = v6699
		goto L679
	} else {
		goto L687
	}
L687:
	;
	goto L685
L688:
	;
	v6729 = v6727
	goto L690
L689:
	;
	v6729 = v6720
	goto L690
L690:
	;
	v6733 = v6729&int32(-1024) + int32(1024)
	v6734 = F_WebPSafeMalloc(m, int64(1), v6733)
	mBase = m.M
	if v6734 != 0 {
		goto L681
	} else {
		goto L691
	}
L691:
	;
	v6735 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6736 = v6735
	goto L682
L692:
	;
	v6747 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v6747)
	mBase = m.M
	v6749 = v6734 + v6733
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v6749
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v6734
	v6754 = v6734 + v6715
	v6755 = v6749
	goto L679
L693:
	;
	v6745 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6746 = F_memcpy(m, v6734, v6745, v6715)
	mBase = m.M
	goto L692
L694:
	;
	goto L678
L695:
	;
	v6963 = *(*int32)(unsafe.Add(mBase, uint32(v271)))
	v6965 = v6963 + int32(-2)
	v6966 = int32(3)
	goto L723
L696:
	;
	goto L695
L697:
	;
	v6824 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v6825 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v6825+v6812 < int32(32) {
		goto L700
	} else {
		goto L701
	}
L698:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v6926 + v6924
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v6923<<(uint(v6926)%32) | v6925
	goto L696
L699:
	;
	v6843 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v6844 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v6848 = v6840
	v6849 = v6841
	v6850 = v6844
	v6851 = v6843
	goto L703
L700:
	;
	if v6825 < int32(16) {
		v6923 = v6811
		v6924 = v6812
		v6925 = v6824
		v6926 = v6825
		goto L698
	} else {
		goto L702
	}
L701:
	;
	v6829 = int32(32)
	v6830 = v6829 - v6825
	v6838 = int32(base.Ui32(v6811) >> (uint(v6830) % 32))
	v6839 = v6812 - v6830
	v6840 = v6811<<(uint(v6825)%32) | v6824
	v6841 = v6829
	goto L699
L702:
	;
	v6838 = v6811
	v6839 = v6812
	v6840 = v6824
	v6841 = v6825
	goto L699
L703:
	;
	if base.Ui32(v6850+int32(2)) <= base.Ui32(v6851) {
		v6906 = v6850
		v6907 = v6851
		goto L705
	} else {
		goto L706
	}
L704:
	;
	v6923 = v6838
	v6924 = v6839
	v6925 = v6919
	v6926 = v6917
	goto L698
L705:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6906))) = uint16(v6848)
	v6914 = v6906 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6914
	v6917 = v6849 + int32(-16)
	v6919 = int32(base.Ui32(v6848) >> (uint(int32(16)) % 32))
	if int32(31) < v6849 {
		v6848 = v6919
		v6849 = v6917
		v6850 = v6914
		v6851 = v6907
		goto L703
	} else {
		goto L720
	}
L706:
	;
	v6860 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6861 = v6851 - v6860
	v6864 = base.I64_extend_i32_s(v6861) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6864) {
		v6888 = v6860
		goto L708
	} else {
		goto L709
	}
L707:
	;
	if v6850 == v6860 {
		goto L718
	} else {
		goto L719
	}
L708:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v6888
	goto L695
L709:
	;
	v6867 = v6850 - v6860
	v6869 = v6864 + base.I64_extend_i32_u(v6867)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6869) {
		v6888 = v6860
		goto L708
	} else {
		goto L710
	}
L710:
	;
	v6872 = base.I32_wrap_i64(v6869)
	if v6851 == v6860 {
		goto L711
	} else {
		goto L712
	}
L711:
	;
	v6879 = int32(base.Ui32(v6861*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6872) < base.Ui32(v6879) {
		goto L714
	} else {
		goto L715
	}
L712:
	;
	if base.Ui32(v6872) <= base.Ui32(v6861) {
		v6906 = v6850
		v6907 = v6851
		goto L705
	} else {
		goto L713
	}
L713:
	;
	goto L711
L714:
	;
	v6881 = v6879
	goto L716
L715:
	;
	v6881 = v6872
	goto L716
L716:
	;
	v6885 = v6881&int32(-1024) + int32(1024)
	v6886 = F_WebPSafeMalloc(m, int64(1), v6885)
	mBase = m.M
	if v6886 != 0 {
		goto L707
	} else {
		goto L717
	}
L717:
	;
	v6887 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6888 = v6887
	goto L708
L718:
	;
	v6899 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v6899)
	mBase = m.M
	v6901 = v6886 + v6885
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v6901
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v6886
	v6906 = v6886 + v6867
	v6907 = v6901
	goto L705
L719:
	;
	v6897 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v6898 = F_memcpy(m, v6886, v6897, v6867)
	mBase = m.M
	goto L718
L720:
	;
	goto L704
L721:
	;
	v7117 = *(*int32)(unsafe.Add(mBase, uint32(v277)+20))
	v7119 = *(*int32)(unsafe.Add(mBase, uint32(v271)))
	v7120 = int32(1) << (uint(v7119) % 32)
	v7127 = *(*int32)(unsafe.Add(mBase, uint32(v277)+4))
	v7131 = F_EncodeImageNoHuffman(m, v272, v7117, v298, v301, int32(base.Ui32(v6647+v7120+int32(-1))>>(uint(v7119)%32)), int32(base.Ui32(v297+v7120)>>(uint(v7119)%32)), v286, v289, v7127, v6646-v6653, v271+int32(60))
	mBase = m.M
	if v7131 == int32(0) {
		v10723 = v269
		v10725 = v271
		goto L4
	} else {
		goto L747
	}
L722:
	;
	goto L721
L723:
	;
	v6978 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v6979 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v6979+v6966 < int32(32) {
		goto L726
	} else {
		goto L727
	}
L724:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v7080 + v7078
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v7077<<(uint(v7080)%32) | v7079
	goto L722
L725:
	;
	v6997 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v6998 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v7002 = v6994
	v7003 = v6995
	v7004 = v6998
	v7005 = v6997
	goto L729
L726:
	;
	if v6979 < int32(16) {
		v7077 = v6965
		v7078 = v6966
		v7079 = v6978
		v7080 = v6979
		goto L724
	} else {
		goto L728
	}
L727:
	;
	v6983 = int32(32)
	v6984 = v6983 - v6979
	v6992 = int32(base.Ui32(v6965) >> (uint(v6984) % 32))
	v6993 = v6966 - v6984
	v6994 = v6965<<(uint(v6979)%32) | v6978
	v6995 = v6983
	goto L725
L728:
	;
	v6992 = v6965
	v6993 = v6966
	v6994 = v6978
	v6995 = v6979
	goto L725
L729:
	;
	if base.Ui32(v7004+int32(2)) <= base.Ui32(v7005) {
		v7060 = v7004
		v7061 = v7005
		goto L731
	} else {
		goto L732
	}
L730:
	;
	v7077 = v6992
	v7078 = v6993
	v7079 = v7073
	v7080 = v7071
	goto L724
L731:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7060))) = uint16(v7002)
	v7068 = v7060 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7068
	v7071 = v7003 + int32(-16)
	v7073 = int32(base.Ui32(v7002) >> (uint(int32(16)) % 32))
	if int32(31) < v7003 {
		v7002 = v7073
		v7003 = v7071
		v7004 = v7068
		v7005 = v7061
		goto L729
	} else {
		goto L746
	}
L732:
	;
	v7014 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7015 = v7005 - v7014
	v7018 = base.I64_extend_i32_s(v7015) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7018) {
		v7042 = v7014
		goto L734
	} else {
		goto L735
	}
L733:
	;
	if v7004 == v7014 {
		goto L744
	} else {
		goto L745
	}
L734:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7042
	goto L721
L735:
	;
	v7021 = v7004 - v7014
	v7023 = v7018 + base.I64_extend_i32_u(v7021)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7023) {
		v7042 = v7014
		goto L734
	} else {
		goto L736
	}
L736:
	;
	v7026 = base.I32_wrap_i64(v7023)
	if v7005 == v7014 {
		goto L737
	} else {
		goto L738
	}
L737:
	;
	v7033 = int32(base.Ui32(v7015*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7026) < base.Ui32(v7033) {
		goto L740
	} else {
		goto L741
	}
L738:
	;
	if base.Ui32(v7026) <= base.Ui32(v7015) {
		v7060 = v7004
		v7061 = v7005
		goto L731
	} else {
		goto L739
	}
L739:
	;
	goto L737
L740:
	;
	v7035 = v7033
	goto L742
L741:
	;
	v7035 = v7026
	goto L742
L742:
	;
	v7039 = v7035&int32(-1024) + int32(1024)
	v7040 = F_WebPSafeMalloc(m, int64(1), v7039)
	mBase = m.M
	if v7040 != 0 {
		goto L733
	} else {
		goto L743
	}
L743:
	;
	v7041 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7042 = v7041
	goto L734
L744:
	;
	v7053 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v7053)
	mBase = m.M
	v7055 = v7040 + v7039
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v7055
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v7040
	v7060 = v7040 + v7021
	v7061 = v7055
	goto L731
L745:
	;
	v7051 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7052 = F_memcpy(m, v7040, v7051, v7021)
	mBase = m.M
	goto L744
L746:
	;
	goto L730
L747:
	;
	v7140 = v6614 - v6646
	goto L666
L748:
	;
	v7293 = *(*int32)(unsafe.Add(mBase, uint32(v277)+32))
	v7294 = *(*int32)(unsafe.Add(mBase, uint32(v277)+36))
	v7295 = *(*int32)(unsafe.Add(mBase, uint32(v277)+8))
	v7296 = *(*int32)(unsafe.Add(mBase, uint32(v271)+60))
	v7297 = int64(57)
	v7298 = int32(16)
	goto L777
L749:
	;
	goto L748
L750:
	;
	v7154 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v7155 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v7155+v7142 < int32(32) {
		goto L753
	} else {
		goto L754
	}
L751:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v7256 + v7254
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v7253<<(uint(v7256)%32) | v7255
	goto L749
L752:
	;
	v7173 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v7174 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v7178 = v7170
	v7179 = v7171
	v7180 = v7174
	v7181 = v7173
	goto L756
L753:
	;
	if v7155 < int32(16) {
		v7253 = v7141
		v7254 = v7142
		v7255 = v7154
		v7256 = v7155
		goto L751
	} else {
		goto L755
	}
L754:
	;
	v7159 = int32(32)
	v7160 = v7159 - v7155
	v7168 = int32(base.Ui32(v7141) >> (uint(v7160) % 32))
	v7169 = v7142 - v7160
	v7170 = v7141<<(uint(v7155)%32) | v7154
	v7171 = v7159
	goto L752
L755:
	;
	v7168 = v7141
	v7169 = v7142
	v7170 = v7154
	v7171 = v7155
	goto L752
L756:
	;
	if base.Ui32(v7180+int32(2)) <= base.Ui32(v7181) {
		v7236 = v7180
		v7237 = v7181
		goto L758
	} else {
		goto L759
	}
L757:
	;
	v7253 = v7168
	v7254 = v7169
	v7255 = v7249
	v7256 = v7247
	goto L751
L758:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7236))) = uint16(v7178)
	v7244 = v7236 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7244
	v7247 = v7179 + int32(-16)
	v7249 = int32(base.Ui32(v7178) >> (uint(int32(16)) % 32))
	if int32(31) < v7179 {
		v7178 = v7249
		v7179 = v7247
		v7180 = v7244
		v7181 = v7237
		goto L756
	} else {
		goto L773
	}
L759:
	;
	v7190 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7191 = v7181 - v7190
	v7194 = base.I64_extend_i32_s(v7191) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7194) {
		v7218 = v7190
		goto L761
	} else {
		goto L762
	}
L760:
	;
	if v7180 == v7190 {
		goto L771
	} else {
		goto L772
	}
L761:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7218
	goto L748
L762:
	;
	v7197 = v7180 - v7190
	v7199 = v7194 + base.I64_extend_i32_u(v7197)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7199) {
		v7218 = v7190
		goto L761
	} else {
		goto L763
	}
L763:
	;
	v7202 = base.I32_wrap_i64(v7199)
	if v7181 == v7190 {
		goto L764
	} else {
		goto L765
	}
L764:
	;
	v7209 = int32(base.Ui32(v7191*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7202) < base.Ui32(v7209) {
		goto L767
	} else {
		goto L768
	}
L765:
	;
	if base.Ui32(v7202) <= base.Ui32(v7191) {
		v7236 = v7180
		v7237 = v7181
		goto L758
	} else {
		goto L766
	}
L766:
	;
	goto L764
L767:
	;
	v7211 = v7209
	goto L769
L768:
	;
	v7211 = v7202
	goto L769
L769:
	;
	v7215 = v7211&int32(-1024) + int32(1024)
	v7216 = F_WebPSafeMalloc(m, int64(1), v7215)
	mBase = m.M
	if v7216 != 0 {
		goto L760
	} else {
		goto L770
	}
L770:
	;
	v7217 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7218 = v7217
	goto L761
L771:
	;
	v7229 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v7229)
	mBase = m.M
	v7231 = v7216 + v7215
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v7231
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v7216
	v7236 = v7216 + v7197
	v7237 = v7231
	goto L758
L772:
	;
	v7227 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7228 = F_memcpy(m, v7216, v7227, v7197)
	mBase = m.M
	goto L771
L773:
	;
	goto L757
L774:
	;
	v7322 = int32(1) << (uint(v7294) % 32)
	v7324 = v7293 + int32(-1)
	v7329 = int32(base.Ui32(v7322+v7324)>>(uint(v7294)%32)) * int32(base.Ui32(v7322+v297)>>(uint(v7294)%32))
	v7330 = base.I64_extend_i32_u(v7329)
	v7331 = int32(4)
	if v7330 == int64(0) {
		goto L782
	} else {
		goto L783
	}
L775:
	;
	goto L774
L776:
	;
	v7317 = F_malloc(m, base.I32_wrap_i64(v7297)*v7298)
	mBase = m.M
	v7319 = v7317
	goto L775
L777:
	;
	v7305 = base.I64_div_u_s(int64(2147418112), v7297)
	v7306 = int32(0)
	v7307 = base.I64_extend_i32_u(v7298)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7307*v7297) {
		v7319 = v7306
		goto L775
	} else {
		goto L778
	}
L778:
	;
	if base.Ui64(v7305) < base.Ui64(v7307) {
		v7319 = v7306
		goto L775
	} else {
		goto L779
	}
L779:
	;
	goto L776
L780:
	;
	v7360 = *(*int64)(unsafe.Add(mBase, uint32(v272+int32(16))))
	*(*int64)(unsafe.Add(mBase, uint32(v271+int32(2128)))) = v7360
	v7362 = int32(0)
	v7363 = base.Simd_g_v128_load(m, v272, v7362)
	base.Simd_g_v128_store(m, v271, int32(2112), v7363)
	v7366 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v271)+64)) = v7366
	v7369 = v271 + int32(1088)
	v7372 = base.Simd_g_const(&F_EncodeStreamHook__k0)
	base.Simd_g_v128_store(m, v7369, v7362, v7372)
	*(*int64)(unsafe.Add(mBase, uint32(v271+int32(1104)))) = v7366
	v7382 = int32(1024)
	v7384 = F_WebPSafeMalloc(m, int64(1), v7382)
	mBase = m.M
	if v7384 != 0 {
		goto L792
	} else {
		goto L793
	}
L781:
	;
	goto L780
L782:
	;
	v7350 = F_malloc(m, base.I32_wrap_i64(v7330)*v7331)
	mBase = m.M
	v7352 = v7350
	goto L781
L783:
	;
	v7338 = base.I64_div_u_s(int64(2147418112), v7330)
	v7339 = int32(0)
	v7340 = base.I64_extend_i32_u(v7331)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7340*v7330) {
		v7352 = v7339
		goto L781
	} else {
		goto L784
	}
L784:
	;
	if base.Ui64(v7338) < base.Ui64(v7340) {
		v7352 = v7339
		goto L781
	} else {
		goto L785
	}
L785:
	;
	goto L782
L786:
	;
	F_free(m, v10471)
	mBase = m.M
	goto L1132
L787:
	;
	v10431 = v10359
	v10433 = v10361
	v10434 = v10362
	v10435 = int32(0)
	v10436 = v10364
	v10437 = v10365
	v10438 = v10366
	v10439 = v10367
	v10440 = v10368
	v10441 = v10369
	v10442 = v10370
	v10443 = v10371
	v10445 = v10373
	v10446 = v10374
	v10448 = v10376
	v10449 = v10377
	v10450 = v10378
	v10451 = v10379
	v10452 = v10380
	v10453 = v10381
	v10454 = v10382
	v10455 = v10383
	v10456 = v10384
	v10457 = v10385
	v10458 = v10386
	v10459 = v10387
	v10460 = v10388
	v10461 = v10389
	v10462 = v10390
	v10463 = v10391
	v10464 = v10392
	v10465 = v10393
	v10466 = v10394
	v10467 = v10395
	v10468 = v10396
	v10469 = v10397
	v10471 = v10399
	v10478 = v10406
	v10480 = v10408
	v10481 = v10409
	goto L786
L788:
	;
	v10359 = v10287
	v10361 = v10289
	v10362 = v10290
	v10364 = v10292
	v10365 = v10293
	v10366 = v10294
	v10367 = v10295
	v10368 = v10296
	v10369 = v10297
	v10370 = v10298
	v10371 = v10299
	v10373 = v10301
	v10374 = int32(0)
	v10376 = v10304
	v10377 = v10305
	v10378 = v10306
	v10379 = v10307
	v10380 = v10308
	v10381 = v10309
	v10382 = v10310
	v10383 = v10311
	v10384 = v10312
	v10385 = v10313
	v10386 = v10314
	v10387 = v10315
	v10388 = v10316
	v10389 = v10317
	v10390 = v10318
	v10391 = v10319
	v10392 = v10320
	v10393 = v10321
	v10394 = v10322
	v10395 = v10323
	v10396 = v10324
	v10397 = v10325
	v10399 = v10327
	v10406 = v10334
	v10408 = v10336
	v10409 = v10337
	goto L787
L789:
	;
	v10285 = int32(0)
	v10287 = v269
	v10289 = v271
	v10290 = v272
	v10292 = v274
	v10293 = v275
	v10294 = v276
	v10295 = v277
	v10296 = v278
	v10297 = v279
	v10298 = v280
	v10299 = v281
	v10301 = v283
	v10304 = v286
	v10305 = v287
	v10306 = v288
	v10307 = v289
	v10308 = v290
	v10309 = v291
	v10310 = v292
	v10311 = v293
	v10312 = v294
	v10313 = v295
	v10314 = v296
	v10315 = v297
	v10316 = v298
	v10317 = v299
	v10318 = v300
	v10319 = v301
	v10320 = v302
	v10321 = v303
	v10322 = v304
	v10323 = v305
	v10324 = v10251
	v10325 = v10252
	v10327 = v10285
	v10334 = v7319
	v10336 = v10285
	v10337 = v7352
	goto L788
L790:
	;
	if v7319 == int32(0) {
		goto L799
	} else {
		goto L800
	}
L791:
	;
	if v7384 != 0 {
		goto L790
	} else {
		goto L794
	}
L792:
	;
	v7388 = *(*int32)(unsafe.Add(mBase, uint32(v7369)+8))
	F_WebPSafeFree(m, v7388)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v7369)+16)) = v7384 + v7382
	*(*int32)(unsafe.Add(mBase, uint32(v7369)+12)) = v7384
	*(*int32)(unsafe.Add(mBase, uint32(v7369)+8)) = v7384
	goto L791
L793:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7369)+20)) = int32(1)
	goto L791
L794:
	;
	v7397 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v7397 != 0 {
		goto L796
	} else {
		goto L797
	}
L795:
	;
	v10251 = v306
	v10252 = v307
	goto L789
L796:
	;
	goto L795
L797:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L796
L798:
	;
	v7419 = int32(0)
	v7421 = base.I32_div_s(v7140, int32(5))
	v7424 = F_VP8LHashChainFill(m, v298, v286, v7295, v7293, v279, v289, v278, v7421, v271+int32(60))
	mBase = m.M
	if v7424 == v7419 {
		v10179 = v306
		v10180 = v307
		goto L809
	} else {
		goto L810
	}
L799:
	;
	v7416 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v7416 != 0 {
		goto L807
	} else {
		goto L808
	}
L800:
	;
	if v7352 == int32(0) {
		goto L799
	} else {
		goto L801
	}
L801:
	;
	v7405 = v271 + int32(64)
	v7409 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v7329), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v7405))) = v7409
	if v7409 != 0 {
		goto L803
	} else {
		goto L804
	}
L802:
	;
	if v7409 != 0 {
		goto L798
	} else {
		goto L805
	}
L803:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7405)+4)) = v7329
	goto L802
L804:
	;
	goto L802
L805:
	;
	goto L799
L806:
	;
	v10251 = v306
	v10252 = v307
	goto L789
L807:
	;
	goto L806
L808:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L807
L809:
	;
	v10287 = v269
	v10289 = v271
	v10290 = v272
	v10292 = v274
	v10293 = v275
	v10294 = v276
	v10295 = v277
	v10296 = v278
	v10297 = v279
	v10298 = v280
	v10299 = v281
	v10301 = v283
	v10304 = v286
	v10305 = v287
	v10306 = v288
	v10307 = v289
	v10308 = v290
	v10309 = v291
	v10310 = v292
	v10311 = v293
	v10312 = v294
	v10313 = v295
	v10314 = v296
	v10315 = v297
	v10316 = v298
	v10317 = v299
	v10318 = v300
	v10319 = v301
	v10320 = v302
	v10321 = v303
	v10322 = v304
	v10323 = v305
	v10324 = v10179
	v10325 = v10180
	v10327 = int32(0)
	v10334 = v7319
	v10336 = v7419
	v10337 = v7352
	goto L788
L810:
	;
	v7427 = *(*int32)(unsafe.Add(mBase, uint32(v277)+48))
	v7428 = *(*int32)(unsafe.Add(mBase, uint32(v342)+24))
	if int32(1) < v7428 {
		goto L813
	} else {
		goto L814
	}
L811:
	;
	v10139 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v10139 != 0 {
		goto L1130
	} else {
		goto L1131
	}
L812:
	;
	if v7509 < int32(1) {
		v10076 = v306
		v10077 = v307
		goto L831
	} else {
		goto L832
	}
L813:
	;
	v7435 = v271 + int32(1088)
	v7445 = *(*int32)(unsafe.Add(mBase, uint32(v7435)+12))
	v7446 = *(*int32)(unsafe.Add(mBase, uint32(v7435)+8))
	v7447 = v7445 - v7446
	v7449 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v7450 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7451 = v7449 - v7450
	v7453 = base.I64_extend_i32_u(v7447) + base.I64_extend_i32_u(v7451)
	if base.Ui64(v7453) < base.Ui64(int64(4294967296)) {
		goto L817
	} else {
		goto L818
	}
L814:
	;
	v7431 = *(*int32)(unsafe.Add(mBase, uint32(v342)+12))
	if v7431 == int32(0) {
		v7509 = v7428
		goto L812
	} else {
		goto L815
	}
L815:
	;
	goto L813
L816:
	;
	if v7505 == int32(0) {
		goto L811
	} else {
		goto L830
	}
L817:
	;
	v7459 = base.I32_wrap_i64(v7453)
	v7460 = *(*int32)(unsafe.Add(mBase, uint32(v7435)+16))
	v7461 = v7460 - v7446
	if v7460 == v7446 {
		goto L820
	} else {
		goto L821
	}
L818:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+20)) = int32(1)
	v7505 = int32(0)
	goto L816
L819:
	;
	v7492 = F_memcpy(m, v7491, v7490, v7451)
	mBase = m.M
	v7493 = *(*int64)(unsafe.Add(mBase, uint32(v272)))
	*(*int64)(unsafe.Add(mBase, uint32(v7435))) = v7493
	v7495 = *(*int32)(unsafe.Add(mBase, uint32(v272)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+20)) = v7495
	v7497 = *(*int32)(unsafe.Add(mBase, uint32(v7435)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+12)) = v7497 + v7451
	v7505 = int32(1)
	goto L816
L820:
	;
	v7468 = int32(base.Ui32(v7461*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7459) < base.Ui32(v7468) {
		goto L824
	} else {
		goto L825
	}
L821:
	;
	if base.Ui32(v7461) < base.Ui32(v7459) {
		goto L820
	} else {
		goto L822
	}
L822:
	;
	v7490 = v7450
	v7491 = v7446
	goto L819
L823:
	;
	if v7445 == v7446 {
		goto L828
	} else {
		goto L829
	}
L824:
	;
	v7470 = v7468
	goto L826
L825:
	;
	v7470 = v7459
	goto L826
L826:
	;
	v7474 = v7470&int32(-1024) + int32(1024)
	v7475 = F_WebPSafeMalloc(m, int64(1), v7474)
	mBase = m.M
	if v7475 != 0 {
		goto L823
	} else {
		goto L827
	}
L827:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+20)) = int32(1)
	v7505 = int32(0)
	goto L816
L828:
	;
	v7482 = *(*int32)(unsafe.Add(mBase, uint32(v7435)+8))
	F_WebPSafeFree(m, v7482)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+8)) = v7475
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+16)) = v7475 + v7474
	*(*int32)(unsafe.Add(mBase, uint32(v7435)+12)) = v7475 + v7447
	v7489 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7490 = v7489
	v7491 = v7475
	goto L819
L829:
	;
	v7480 = *(*int32)(unsafe.Add(mBase, uint32(v7435)+8))
	v7481 = F_memcpy(m, v7475, v7480, v7447)
	mBase = m.M
	goto L828
L830:
	;
	v7508 = *(*int32)(unsafe.Add(mBase, uint32(v342)+24))
	v7509 = v7508
	goto L812
L831:
	;
	v10111 = v271 + int32(1088)
	v10114 = int32(0)
	v10118 = *(*int64)(unsafe.Add(mBase, uint32(v272)))
	v10119 = *(*int64)(unsafe.Add(mBase, uint32(v10111)))
	*(*int64)(unsafe.Add(mBase, uint32(v272))) = v10119
	v10122 = v272 + int32(8)
	v10124 = base.Simd_g_v128_load(m, v10122, v10114)
	v10126 = v271 + int32(1096)
	v10128 = base.Simd_g_v128_load(m, v10126, v10114)
	base.Simd_g_v128_store(m, v10122, v10114, v10128)
	*(*int64)(unsafe.Add(mBase, uint32(v10111))) = v10118
	base.Simd_g_v128_store(m, v10126, v10114, v10124)
	goto L1128
L832:
	;
	if v7427 != 0 {
		goto L833
	} else {
		goto L834
	}
L833:
	;
	v7514 = v7427
	goto L835
L834:
	;
	v7514 = int32(10)
	goto L835
L835:
	;
	v7527 = v7509
	v7560 = v306
	v7561 = v307
	v7587 = int32(-1)
	v7588 = int32(0)
	goto L836
L836:
	;
	v7596 = v342 + int32(8) + v7588<<(uint(int32(3))%32)
	v7597 = *(*int32)(unsafe.Add(mBase, uint32(v7596)))
	v7598 = *(*int32)(unsafe.Add(mBase, uint32(v7596)+4))
	v7601 = base.I32_div_s(v7140-v7421, v7527)
	v7603 = base.I32_div_s(v7601, int32(4))
	v7606 = F_VP8LGetBackwardReferences(m, v7293, v279, v7295, v286, v289, v7597, v7514, v7598, v298, v301, v271+int32(_a_F_EncodeStreamHook_6), v278, v7603, v271+int32(60))
	mBase = m.M
	if v7606 == int32(0) {
		v10251 = v7560
		v10252 = v7561
		goto L789
	} else {
		goto L838
	}
L837:
	;
	v10076 = v10001
	v10077 = v10002
	goto L831
L838:
	;
	v7628 = int32(0)
	v7649 = v7560
	v7650 = v7561
	v7676 = v7587
	v7679 = v7601 - v7603
	v7680 = int32(1)
	goto L840
L839:
	;
	v10036 = v7588 + int32(1)
	v10037 = *(*int32)(unsafe.Add(mBase, uint32(v342)+24))
	if v10036 < v10037 {
		v7527 = v10037
		v7560 = v10001
		v7561 = v10002
		v7587 = v10028
		v7588 = v10036
		goto L836
	} else {
		goto L1127
	}
L840:
	;
	v7683 = *(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_EncodeStreamHook[0])))
	*(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_EncodeStreamHook[1]))) = v7294
	if v7628&int32(1) == int32(0) {
		goto L842
	} else {
		goto L843
	}
L841:
	;
	v10001 = v9951
	v10002 = v9952
	v10028 = v9953
	goto L839
L842:
	;
	v7692 = v271 + int32(2112)
	v7693 = *(*int64)(unsafe.Add(mBase, uint32(v7692)))
	*(*int64)(unsafe.Add(mBase, uint32(v272))) = v7693
	v7695 = *(*int32)(unsafe.Add(mBase, uint32(v7692)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = v7695
	v7697 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7698 = *(*int32)(unsafe.Add(mBase, uint32(v7692)+12))
	v7699 = *(*int32)(unsafe.Add(mBase, uint32(v7692)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7697 + (v7698 - v7699)
	goto L845
L843:
	;
	if v7683 == int32(0) {
		v10001 = v7649
		v10002 = v7650
		v10028 = v7676
		goto L839
	} else {
		goto L844
	}
L844:
	;
	goto L842
L845:
	;
	v7703 = int32(0)
	if v7680&int32(1) != 0 {
		goto L846
	} else {
		goto L847
	}
L846:
	;
	v7707 = v7683
	goto L848
L847:
	;
	v7707 = v7703
	goto L848
L848:
	;
	v7708 = F_VP8LAllocateHistogramSet(m, v7329, v7707)
	mBase = m.M
	v7713 = int32(_a_F_EncodeStreamHook_7)
	if int32(0) < v7707 {
		goto L851
	} else {
		goto L852
	}
L849:
	;
	if v7708 == int32(0) {
		goto L856
	} else {
		goto L857
	}
L850:
	;
	goto L849
L851:
	;
	v7718 = int32(4)<<(uint(v7707)%32) + v7713
	goto L853
L852:
	;
	v7718 = v7713
	goto L853
L853:
	;
	v7721 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v7718), int32(1))
	mBase = m.M
	if v7721 == int32(0) {
		goto L850
	} else {
		goto L854
	}
L854:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v7721)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v7721)+3236)) = v7707
	*(*int32)(unsafe.Add(mBase, uint32(v7721)+3304)) = int32(16843009)
	v7729 = base.Simd_g_const(&F_EncodeStreamHook__k0)
	base.Simd_g_v128_store(m, v7721, int32(3256), v7729)
	*(*int32)(unsafe.Add(mBase, uint32(v7721))) = v7721 + int32(3312)
	v7737 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v7721+int32(3248)))) = uint16(v7737)
	v7741 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v7721+int32(3308)))) = uint8(v7741)
	v7745 = int32(0)
	base.Simd_g_v128_store(m, v7721+int32(3272), v7745, v7729)
	base.Simd_g_v128_store(m, v7721+int32(3288), v7745, v7729)
	goto L850
L855:
	;
	v7801 = F_GetHuffBitLengthsAndCodes(m, v7708, v7791)
	mBase = m.M
	if v7801 != 0 {
		goto L871
	} else {
		goto L872
	}
L856:
	;
	v7798 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v7798 != 0 {
		goto L869
	} else {
		goto L870
	}
L857:
	;
	if v7721 == int32(0) {
		goto L856
	} else {
		goto L858
	}
L858:
	;
	v7758 = v301 + v7628*int32(24)
	v7759 = *(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_EncodeStreamHook[1])))
	v7761 = base.I32_div_s(v7679, int32(3))
	v7764 = F_VP8LGetHistoImageSymbols(m, v7293, v279, v7758, v286, v289, v7759, v7707, v7708, v7721, v7352, v278, v7761, v271+int32(60))
	mBase = m.M
	if v7764 != 0 {
		goto L859
	} else {
		goto L860
	}
L859:
	;
	v7767 = *(*int32)(unsafe.Add(mBase, uint32(v7708)))
	v7770 = base.I64_extend_i32_u(v7767 * int32(5))
	v7771 = int32(12)
	if v7770 == int64(0) {
		goto L863
	} else {
		goto L864
	}
L860:
	;
	v7765 = int32(0)
	v10431 = v269
	v10433 = v271
	v10434 = v272
	v10435 = v7708
	v10436 = v274
	v10437 = v275
	v10438 = v276
	v10439 = v277
	v10440 = v278
	v10441 = v279
	v10442 = v280
	v10443 = v281
	v10445 = v283
	v10446 = v7721
	v10448 = v286
	v10449 = v287
	v10450 = v288
	v10451 = v289
	v10452 = v290
	v10453 = v291
	v10454 = v292
	v10455 = v293
	v10456 = v294
	v10457 = v295
	v10458 = v296
	v10459 = v297
	v10460 = v298
	v10461 = v299
	v10462 = v300
	v10463 = v301
	v10464 = v302
	v10465 = v303
	v10466 = v304
	v10467 = v305
	v10468 = v7649
	v10469 = v7650
	v10471 = v7765
	v10478 = v7319
	v10480 = v7765
	v10481 = v7352
	goto L786
L861:
	;
	if v7791 != 0 {
		goto L855
	} else {
		goto L867
	}
L862:
	;
	goto L861
L863:
	;
	v7789 = F_calloc(m, base.I32_wrap_i64(v7770), v7771)
	mBase = m.M
	v7791 = v7789
	goto L862
L864:
	;
	v7778 = base.I64_div_u_s(int64(2147418112), v7770)
	v7779 = int32(0)
	v7780 = base.I64_extend_i32_u(v7771)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7780*v7770) {
		v7791 = v7779
		goto L862
	} else {
		goto L865
	}
L865:
	;
	if base.Ui64(v7778) < base.Ui64(v7780) {
		v7791 = v7779
		goto L862
	} else {
		goto L866
	}
L866:
	;
	goto L863
L867:
	;
	goto L856
L868:
	;
	v10431 = v269
	v10433 = v271
	v10434 = v272
	v10435 = v7708
	v10436 = v274
	v10437 = v275
	v10438 = v276
	v10439 = v277
	v10440 = v278
	v10441 = v279
	v10442 = v280
	v10443 = v281
	v10445 = v283
	v10446 = v7721
	v10448 = v286
	v10449 = v287
	v10450 = v288
	v10451 = v289
	v10452 = v290
	v10453 = v291
	v10454 = v292
	v10455 = v293
	v10456 = v294
	v10457 = v295
	v10458 = v296
	v10459 = v297
	v10460 = v298
	v10461 = v299
	v10462 = v300
	v10463 = v301
	v10464 = v302
	v10465 = v303
	v10466 = v304
	v10467 = v305
	v10468 = v7649
	v10469 = v7650
	v10471 = v7703
	v10478 = v7319
	v10480 = int32(0)
	v10481 = v7352
	goto L786
L869:
	;
	goto L868
L870:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L869
L871:
	;
	F_WebPSafeFree(m, v7708)
	mBase = m.M
	goto L876
L872:
	;
	v7803 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v7803 != 0 {
		goto L874
	} else {
		goto L875
	}
L873:
	;
	v10431 = v269
	v10433 = v271
	v10434 = v272
	v10435 = v7708
	v10436 = v274
	v10437 = v275
	v10438 = v276
	v10439 = v277
	v10440 = v278
	v10441 = v279
	v10442 = v280
	v10443 = v281
	v10445 = v283
	v10446 = v7721
	v10448 = v286
	v10449 = v287
	v10450 = v288
	v10451 = v289
	v10452 = v290
	v10453 = v291
	v10454 = v292
	v10455 = v293
	v10456 = v294
	v10457 = v295
	v10458 = v296
	v10459 = v297
	v10460 = v298
	v10461 = v299
	v10462 = v300
	v10463 = v301
	v10464 = v302
	v10465 = v303
	v10466 = v304
	v10467 = v305
	v10468 = v7649
	v10469 = v7650
	v10471 = v7703
	v10478 = v7319
	v10480 = v7791
	v10481 = v7352
	goto L786
L874:
	;
	goto L873
L875:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L874
L876:
	;
	F_WebPSafeFree(m, v7721)
	mBase = m.M
	goto L877
L877:
	;
	if v7707 < int32(1) {
		goto L879
	} else {
		goto L880
	}
L878:
	;
	v8265 = v7679 - v7761
	v8266 = int32(0)
	switch v7329 {
	case 0:
		goto L961
	case 1:
		v8379 = v8266
		v8406 = v8266
		goto L962
	default:
		goto L963
	}
L879:
	;
	v8113 = int32(0)
	v8114 = int32(1)
	goto L935
L880:
	;
	v7810 = int32(1)
	goto L883
L881:
	;
	v7962 = int32(4)
	goto L909
L882:
	;
	goto L881
L883:
	;
	v7823 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v7824 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v7824+v7810 < int32(32) {
		goto L886
	} else {
		goto L887
	}
L884:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v7925 + v7923
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v7922<<(uint(v7925)%32) | v7924
	goto L882
L885:
	;
	v7842 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v7843 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v7847 = v7839
	v7848 = v7840
	v7849 = v7843
	v7850 = v7842
	goto L889
L886:
	;
	if v7824 < int32(16) {
		v7922 = v7810
		v7923 = v7810
		v7924 = v7823
		v7925 = v7824
		goto L884
	} else {
		goto L888
	}
L887:
	;
	v7828 = int32(32)
	v7829 = v7828 - v7824
	v7837 = int32(base.Ui32(v7810) >> (uint(v7829) % 32))
	v7838 = v7810 - v7829
	v7839 = v7810<<(uint(v7824)%32) | v7823
	v7840 = v7828
	goto L885
L888:
	;
	v7837 = v7810
	v7838 = v7810
	v7839 = v7823
	v7840 = v7824
	goto L885
L889:
	;
	if base.Ui32(v7849+int32(2)) <= base.Ui32(v7850) {
		v7905 = v7849
		v7906 = v7850
		goto L891
	} else {
		goto L892
	}
L890:
	;
	v7922 = v7837
	v7923 = v7838
	v7924 = v7918
	v7925 = v7916
	goto L884
L891:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7905))) = uint16(v7847)
	v7913 = v7905 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7913
	v7916 = v7848 + int32(-16)
	v7918 = int32(base.Ui32(v7847) >> (uint(int32(16)) % 32))
	if int32(31) < v7848 {
		v7847 = v7918
		v7848 = v7916
		v7849 = v7913
		v7850 = v7906
		goto L889
	} else {
		goto L906
	}
L892:
	;
	v7859 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7860 = v7850 - v7859
	v7863 = base.I64_extend_i32_s(v7860) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7863) {
		v7887 = v7859
		goto L894
	} else {
		goto L895
	}
L893:
	;
	if v7849 == v7859 {
		goto L904
	} else {
		goto L905
	}
L894:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v7887
	goto L881
L895:
	;
	v7866 = v7849 - v7859
	v7868 = v7863 + base.I64_extend_i32_u(v7866)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7868) {
		v7887 = v7859
		goto L894
	} else {
		goto L896
	}
L896:
	;
	v7871 = base.I32_wrap_i64(v7868)
	if v7850 == v7859 {
		goto L897
	} else {
		goto L898
	}
L897:
	;
	v7878 = int32(base.Ui32(v7860*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7871) < base.Ui32(v7878) {
		goto L900
	} else {
		goto L901
	}
L898:
	;
	if base.Ui32(v7871) <= base.Ui32(v7860) {
		v7905 = v7849
		v7906 = v7850
		goto L891
	} else {
		goto L899
	}
L899:
	;
	goto L897
L900:
	;
	v7880 = v7878
	goto L902
L901:
	;
	v7880 = v7871
	goto L902
L902:
	;
	v7884 = v7880&int32(-1024) + int32(1024)
	v7885 = F_WebPSafeMalloc(m, int64(1), v7884)
	mBase = m.M
	if v7885 != 0 {
		goto L893
	} else {
		goto L903
	}
L903:
	;
	v7886 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7887 = v7886
	goto L894
L904:
	;
	v7898 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v7898)
	mBase = m.M
	v7900 = v7885 + v7884
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v7900
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v7885
	v7905 = v7885 + v7866
	v7906 = v7900
	goto L891
L905:
	;
	v7896 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v7897 = F_memcpy(m, v7885, v7896, v7866)
	mBase = m.M
	goto L904
L906:
	;
	goto L890
L907:
	;
	goto L878
L908:
	;
	goto L907
L909:
	;
	v7974 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v7975 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v7975+v7962 < int32(32) {
		goto L912
	} else {
		goto L913
	}
L910:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v8076 + v8074
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v8073<<(uint(v8076)%32) | v8075
	goto L908
L911:
	;
	v7993 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v7994 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v7998 = v7990
	v7999 = v7991
	v8000 = v7994
	v8001 = v7993
	goto L915
L912:
	;
	if v7975 < int32(16) {
		v8073 = v7707
		v8074 = v7962
		v8075 = v7974
		v8076 = v7975
		goto L910
	} else {
		goto L914
	}
L913:
	;
	v7979 = int32(32)
	v7980 = v7979 - v7975
	v7988 = int32(base.Ui32(v7707) >> (uint(v7980) % 32))
	v7989 = v7962 - v7980
	v7990 = v7707<<(uint(v7975)%32) | v7974
	v7991 = v7979
	goto L911
L914:
	;
	v7988 = v7707
	v7989 = v7962
	v7990 = v7974
	v7991 = v7975
	goto L911
L915:
	;
	if base.Ui32(v8000+int32(2)) <= base.Ui32(v8001) {
		v8056 = v8000
		v8057 = v8001
		goto L917
	} else {
		goto L918
	}
L916:
	;
	v8073 = v7988
	v8074 = v7989
	v8075 = v8069
	v8076 = v8067
	goto L910
L917:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v8056))) = uint16(v7998)
	v8064 = v8056 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8064
	v8067 = v7999 + int32(-16)
	v8069 = int32(base.Ui32(v7998) >> (uint(int32(16)) % 32))
	if int32(31) < v7999 {
		v7998 = v8069
		v7999 = v8067
		v8000 = v8064
		v8001 = v8057
		goto L915
	} else {
		goto L932
	}
L918:
	;
	v8010 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8011 = v8001 - v8010
	v8014 = base.I64_extend_i32_s(v8011) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8014) {
		v8038 = v8010
		goto L920
	} else {
		goto L921
	}
L919:
	;
	if v8000 == v8010 {
		goto L930
	} else {
		goto L931
	}
L920:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8038
	goto L907
L921:
	;
	v8017 = v8000 - v8010
	v8019 = v8014 + base.I64_extend_i32_u(v8017)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8019) {
		v8038 = v8010
		goto L920
	} else {
		goto L922
	}
L922:
	;
	v8022 = base.I32_wrap_i64(v8019)
	if v8001 == v8010 {
		goto L923
	} else {
		goto L924
	}
L923:
	;
	v8029 = int32(base.Ui32(v8011*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v8022) < base.Ui32(v8029) {
		goto L926
	} else {
		goto L927
	}
L924:
	;
	if base.Ui32(v8022) <= base.Ui32(v8011) {
		v8056 = v8000
		v8057 = v8001
		goto L917
	} else {
		goto L925
	}
L925:
	;
	goto L923
L926:
	;
	v8031 = v8029
	goto L928
L927:
	;
	v8031 = v8022
	goto L928
L928:
	;
	v8035 = v8031&int32(-1024) + int32(1024)
	v8036 = F_WebPSafeMalloc(m, int64(1), v8035)
	mBase = m.M
	if v8036 != 0 {
		goto L919
	} else {
		goto L929
	}
L929:
	;
	v8037 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8038 = v8037
	goto L920
L930:
	;
	v8049 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v8049)
	mBase = m.M
	v8051 = v8036 + v8035
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v8051
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v8036
	v8056 = v8036 + v8017
	v8057 = v8051
	goto L917
L931:
	;
	v8047 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8048 = F_memcpy(m, v8036, v8047, v8017)
	mBase = m.M
	goto L930
L932:
	;
	goto L916
L933:
	;
	goto L878
L934:
	;
	goto L933
L935:
	;
	v8126 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v8127 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v8127+v8114 < int32(32) {
		goto L938
	} else {
		goto L939
	}
L936:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v8228 + v8226
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v8225<<(uint(v8228)%32) | v8227
	goto L934
L937:
	;
	v8145 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v8146 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v8150 = v8142
	v8151 = v8143
	v8152 = v8146
	v8153 = v8145
	goto L941
L938:
	;
	if v8127 < int32(16) {
		v8225 = v8113
		v8226 = v8114
		v8227 = v8126
		v8228 = v8127
		goto L936
	} else {
		goto L940
	}
L939:
	;
	v8131 = int32(32)
	v8132 = v8131 - v8127
	v8140 = int32(base.Ui32(v8113) >> (uint(v8132) % 32))
	v8141 = v8114 - v8132
	v8142 = v8113<<(uint(v8127)%32) | v8126
	v8143 = v8131
	goto L937
L940:
	;
	v8140 = v8113
	v8141 = v8114
	v8142 = v8126
	v8143 = v8127
	goto L937
L941:
	;
	if base.Ui32(v8152+int32(2)) <= base.Ui32(v8153) {
		v8208 = v8152
		v8209 = v8153
		goto L943
	} else {
		goto L944
	}
L942:
	;
	v8225 = v8140
	v8226 = v8141
	v8227 = v8221
	v8228 = v8219
	goto L936
L943:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v8208))) = uint16(v8150)
	v8216 = v8208 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8216
	v8219 = v8151 + int32(-16)
	v8221 = int32(base.Ui32(v8150) >> (uint(int32(16)) % 32))
	if int32(31) < v8151 {
		v8150 = v8221
		v8151 = v8219
		v8152 = v8216
		v8153 = v8209
		goto L941
	} else {
		goto L958
	}
L944:
	;
	v8162 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8163 = v8153 - v8162
	v8166 = base.I64_extend_i32_s(v8163) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8166) {
		v8190 = v8162
		goto L946
	} else {
		goto L947
	}
L945:
	;
	if v8152 == v8162 {
		goto L956
	} else {
		goto L957
	}
L946:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8190
	goto L933
L947:
	;
	v8169 = v8152 - v8162
	v8171 = v8166 + base.I64_extend_i32_u(v8169)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8171) {
		v8190 = v8162
		goto L946
	} else {
		goto L948
	}
L948:
	;
	v8174 = base.I32_wrap_i64(v8171)
	if v8153 == v8162 {
		goto L949
	} else {
		goto L950
	}
L949:
	;
	v8181 = int32(base.Ui32(v8163*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v8174) < base.Ui32(v8181) {
		goto L952
	} else {
		goto L953
	}
L950:
	;
	if base.Ui32(v8174) <= base.Ui32(v8163) {
		v8208 = v8152
		v8209 = v8153
		goto L943
	} else {
		goto L951
	}
L951:
	;
	goto L949
L952:
	;
	v8183 = v8181
	goto L954
L953:
	;
	v8183 = v8174
	goto L954
L954:
	;
	v8187 = v8183&int32(-1024) + int32(1024)
	v8188 = F_WebPSafeMalloc(m, int64(1), v8187)
	mBase = m.M
	if v8188 != 0 {
		goto L945
	} else {
		goto L955
	}
L955:
	;
	v8189 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8190 = v8189
	goto L946
L956:
	;
	v8201 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v8201)
	mBase = m.M
	v8203 = v8188 + v8187
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v8203
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v8188
	v8208 = v8188 + v8169
	v8209 = v8203
	goto L943
L957:
	;
	v8199 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8200 = F_memcpy(m, v8188, v8199, v8169)
	mBase = m.M
	goto L956
L958:
	;
	goto L942
L959:
	;
	v9894 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v9895 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v9896 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	v9897 = *(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_EncodeStreamHook[1])))
	v9898 = F_StoreImageToBitMask(m, v272, v7293, v9897, v7758, v7352, v7791, v278)
	mBase = m.M
	if v9898 != 0 {
		goto L1118
	} else {
		goto L1119
	}
L960:
	;
	v9820 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if v9820 != 0 {
		goto L1116
	} else {
		goto L1117
	}
L961:
	;
	v9570 = int32(0)
	v9572 = int32(1)
	goto L1084
L962:
	;
	if v7329&int32(1) == int32(0) {
		v8449 = v8379
		goto L973
	} else {
		goto L974
	}
L963:
	;
	v8268 = int32(0)
	v8274 = v7352
	v8285 = v8268
	v8312 = v8268
	goto L964
L964:
	;
	v8341 = *(*int32)(unsafe.Add(mBase, uint32(v8274)))
	v8342 = int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v8274))) = v8341 << (uint(v8342) % 32)
	v8346 = v8274 + int32(4)
	v8347 = *(*int32)(unsafe.Add(mBase, uint32(v8346)))
	*(*int32)(unsafe.Add(mBase, uint32(v8346))) = v8347 << (uint(v8342) % 32)
	if base.Ui32(v8341) < base.Ui32(v8285) {
		goto L966
	} else {
		goto L967
	}
L965:
	;
	v8379 = v8358
	v8406 = v8362
	goto L962
L966:
	;
	v8354 = v8285
	goto L968
L967:
	;
	v8354 = v8341 + int32(1)
	goto L968
L968:
	;
	if base.Ui32(v8347) < base.Ui32(v8354) {
		goto L969
	} else {
		goto L970
	}
L969:
	;
	v8358 = v8354
	goto L971
L970:
	;
	v8358 = v8347 + int32(1)
	goto L971
L971:
	;
	v8362 = v8312 + int32(2)
	if v7329&int32(-2) != v8362 {
		v8274 = v8274 + int32(8)
		v8285 = v8358
		v8312 = v8362
		goto L964
	} else {
		goto L972
	}
L972:
	;
	goto L965
L973:
	;
	v8450 = int32(1)
	v8451 = base.B2i32(base.Ui32(v8450) < base.Ui32(v8449))
	goto L980
L974:
	;
	v8439 = v7352 + v8406<<(uint(int32(2))%32)
	v8440 = *(*int32)(unsafe.Add(mBase, uint32(v8439)))
	*(*int32)(unsafe.Add(mBase, uint32(v8439))) = v8440 << (uint(int32(8)) % 32)
	if base.Ui32(v8440) < base.Ui32(v8379) {
		goto L975
	} else {
		goto L976
	}
L975:
	;
	v8447 = v8379
	goto L977
L976:
	;
	v8447 = v8440 + int32(1)
	goto L977
L977:
	;
	v8449 = v8447
	goto L973
L978:
	;
	if base.Ui32(v8449) < base.Ui32(int32(2)) {
		goto L1007
	} else {
		goto L1008
	}
L979:
	;
	goto L978
L980:
	;
	v8464 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v8465 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v8465+v8450 < int32(32) {
		goto L983
	} else {
		goto L984
	}
L981:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v8566 + v8564
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v8563<<(uint(v8566)%32) | v8565
	goto L979
L982:
	;
	v8483 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v8484 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v8488 = v8480
	v8489 = v8481
	v8490 = v8484
	v8491 = v8483
	goto L986
L983:
	;
	if v8465 < int32(16) {
		v8563 = v8451
		v8564 = v8450
		v8565 = v8464
		v8566 = v8465
		goto L981
	} else {
		goto L985
	}
L984:
	;
	v8469 = int32(32)
	v8470 = v8469 - v8465
	v8478 = int32(base.Ui32(v8451) >> (uint(v8470) % 32))
	v8479 = v8450 - v8470
	v8480 = v8451<<(uint(v8465)%32) | v8464
	v8481 = v8469
	goto L982
L985:
	;
	v8478 = v8451
	v8479 = v8450
	v8480 = v8464
	v8481 = v8465
	goto L982
L986:
	;
	if base.Ui32(v8490+int32(2)) <= base.Ui32(v8491) {
		v8546 = v8490
		v8547 = v8491
		goto L988
	} else {
		goto L989
	}
L987:
	;
	v8563 = v8478
	v8564 = v8479
	v8565 = v8559
	v8566 = v8557
	goto L981
L988:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v8546))) = uint16(v8488)
	v8554 = v8546 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8554
	v8557 = v8489 + int32(-16)
	v8559 = int32(base.Ui32(v8488) >> (uint(int32(16)) % 32))
	if int32(31) < v8489 {
		v8488 = v8559
		v8489 = v8557
		v8490 = v8554
		v8491 = v8547
		goto L986
	} else {
		goto L1003
	}
L989:
	;
	v8500 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8501 = v8491 - v8500
	v8504 = base.I64_extend_i32_s(v8501) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8504) {
		v8528 = v8500
		goto L991
	} else {
		goto L992
	}
L990:
	;
	if v8490 == v8500 {
		goto L1001
	} else {
		goto L1002
	}
L991:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8528
	goto L978
L992:
	;
	v8507 = v8490 - v8500
	v8509 = v8504 + base.I64_extend_i32_u(v8507)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8509) {
		v8528 = v8500
		goto L991
	} else {
		goto L993
	}
L993:
	;
	v8512 = base.I32_wrap_i64(v8509)
	if v8491 == v8500 {
		goto L994
	} else {
		goto L995
	}
L994:
	;
	v8519 = int32(base.Ui32(v8501*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v8512) < base.Ui32(v8519) {
		goto L997
	} else {
		goto L998
	}
L995:
	;
	if base.Ui32(v8512) <= base.Ui32(v8501) {
		v8546 = v8490
		v8547 = v8491
		goto L988
	} else {
		goto L996
	}
L996:
	;
	goto L994
L997:
	;
	v8521 = v8519
	goto L999
L998:
	;
	v8521 = v8512
	goto L999
L999:
	;
	v8525 = v8521&int32(-1024) + int32(1024)
	v8526 = F_WebPSafeMalloc(m, int64(1), v8525)
	mBase = m.M
	if v8526 != 0 {
		goto L990
	} else {
		goto L1000
	}
L1000:
	;
	v8527 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8528 = v8527
	goto L991
L1001:
	;
	v8539 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v8539)
	mBase = m.M
	v8541 = v8526 + v8525
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v8541
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v8526
	v8546 = v8526 + v8507
	v8547 = v8541
	goto L988
L1002:
	;
	v8537 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8538 = F_memcpy(m, v8526, v8537, v8507)
	mBase = m.M
	goto L1001
L1003:
	;
	goto L987
L1004:
	;
	v9762 = int32(0)
	goto L960
L1005:
	;
	goto L1077
L1006:
	;
	v8785 = v8449 * int32(5)
	v8786 = int32(1)
	if base.Ui32(v8786) < base.Ui32(v8785) {
		goto L1038
	} else {
		goto L1039
	}
L1007:
	;
	if v8449 == int32(0) {
		goto L1005
	} else {
		goto L1037
	}
L1008:
	;
	F_VP8LOptimizeSampling(m, v7352, v7293, v279, v7294, int32(9), v271+int32(_a_F_EncodeStreamHook_8))
	mBase = m.M
	v8609 = *(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_EncodeStreamHook[1])))
	v8611 = v8609 + int32(-2)
	v8612 = int32(3)
	goto L1011
L1009:
	;
	v8766 = *(*int32)(unsafe.Add(mBase, uint32(v271)+uint32(_c_F_EncodeStreamHook[1])))
	v8767 = int32(1) << (uint(v8766) % 32)
	v8773 = base.I32_div_s(v8265, int32(2))
	v8776 = F_EncodeImageNoHuffman(m, v272, v7352, v271+int32(64), v296, int32(base.Ui32(v8767+v7324)>>(uint(v8766)%32)), int32(base.Ui32(v8767+v297)>>(uint(v8766)%32)), v286, v289, v278, v8773, v271+int32(60))
	mBase = m.M
	if v8776 != 0 {
		goto L1035
	} else {
		goto L1036
	}
L1010:
	;
	goto L1009
L1011:
	;
	v8624 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v8625 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v8625+v8612 < int32(32) {
		goto L1014
	} else {
		goto L1015
	}
L1012:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v8726 + v8724
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v8723<<(uint(v8726)%32) | v8725
	goto L1010
L1013:
	;
	v8643 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v8644 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v8648 = v8640
	v8649 = v8641
	v8650 = v8644
	v8651 = v8643
	goto L1017
L1014:
	;
	if v8625 < int32(16) {
		v8723 = v8611
		v8724 = v8612
		v8725 = v8624
		v8726 = v8625
		goto L1012
	} else {
		goto L1016
	}
L1015:
	;
	v8629 = int32(32)
	v8630 = v8629 - v8625
	v8638 = int32(base.Ui32(v8611) >> (uint(v8630) % 32))
	v8639 = v8612 - v8630
	v8640 = v8611<<(uint(v8625)%32) | v8624
	v8641 = v8629
	goto L1013
L1016:
	;
	v8638 = v8611
	v8639 = v8612
	v8640 = v8624
	v8641 = v8625
	goto L1013
L1017:
	;
	if base.Ui32(v8650+int32(2)) <= base.Ui32(v8651) {
		v8706 = v8650
		v8707 = v8651
		goto L1019
	} else {
		goto L1020
	}
L1018:
	;
	v8723 = v8638
	v8724 = v8639
	v8725 = v8719
	v8726 = v8717
	goto L1012
L1019:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v8706))) = uint16(v8648)
	v8714 = v8706 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8714
	v8717 = v8649 + int32(-16)
	v8719 = int32(base.Ui32(v8648) >> (uint(int32(16)) % 32))
	if int32(31) < v8649 {
		v8648 = v8719
		v8649 = v8717
		v8650 = v8714
		v8651 = v8707
		goto L1017
	} else {
		goto L1034
	}
L1020:
	;
	v8660 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8661 = v8651 - v8660
	v8664 = base.I64_extend_i32_s(v8661) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8664) {
		v8688 = v8660
		goto L1022
	} else {
		goto L1023
	}
L1021:
	;
	if v8650 == v8660 {
		goto L1032
	} else {
		goto L1033
	}
L1022:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v8688
	goto L1009
L1023:
	;
	v8667 = v8650 - v8660
	v8669 = v8664 + base.I64_extend_i32_u(v8667)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8669) {
		v8688 = v8660
		goto L1022
	} else {
		goto L1024
	}
L1024:
	;
	v8672 = base.I32_wrap_i64(v8669)
	if v8651 == v8660 {
		goto L1025
	} else {
		goto L1026
	}
L1025:
	;
	v8679 = int32(base.Ui32(v8661*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v8672) < base.Ui32(v8679) {
		goto L1028
	} else {
		goto L1029
	}
L1026:
	;
	if base.Ui32(v8672) <= base.Ui32(v8661) {
		v8706 = v8650
		v8707 = v8651
		goto L1019
	} else {
		goto L1027
	}
L1027:
	;
	goto L1025
L1028:
	;
	v8681 = v8679
	goto L1030
L1029:
	;
	v8681 = v8672
	goto L1030
L1030:
	;
	v8685 = v8681&int32(-1024) + int32(1024)
	v8686 = F_WebPSafeMalloc(m, int64(1), v8685)
	mBase = m.M
	if v8686 != 0 {
		goto L1021
	} else {
		goto L1031
	}
L1031:
	;
	v8687 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8688 = v8687
	goto L1022
L1032:
	;
	v8699 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v8699)
	mBase = m.M
	v8701 = v8686 + v8685
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v8701
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v8686
	v8706 = v8686 + v8667
	v8707 = v8701
	goto L1019
L1033:
	;
	v8697 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v8698 = F_memcpy(m, v8686, v8697, v8667)
	mBase = m.M
	goto L1032
L1034:
	;
	goto L1018
L1035:
	;
	v8783 = v8265 - v8773
	goto L1006
L1036:
	;
	v10287 = v269
	v10289 = v271
	v10290 = v272
	v10292 = v274
	v10293 = v275
	v10294 = v276
	v10295 = v277
	v10296 = v278
	v10297 = v279
	v10298 = v280
	v10299 = v281
	v10301 = v283
	v10304 = v286
	v10305 = v287
	v10306 = v288
	v10307 = v289
	v10308 = v290
	v10309 = v291
	v10310 = v292
	v10311 = v293
	v10312 = v294
	v10313 = v295
	v10314 = v296
	v10315 = v297
	v10316 = v298
	v10317 = v299
	v10318 = v300
	v10319 = v301
	v10320 = v302
	v10321 = v303
	v10322 = v304
	v10323 = v305
	v10324 = v7649
	v10325 = v7650
	v10327 = int32(0)
	v10334 = v7319
	v10336 = v7791
	v10337 = v7352
	goto L788
L1037:
	;
	v8783 = v8265
	goto L1006
L1038:
	;
	v8789 = v8785
	goto L1040
L1039:
	;
	v8789 = v8786
	goto L1040
L1040:
	;
	v8790 = int32(0)
	if base.Ui32(v8785) < base.Ui32(int32(4)) {
		v8903 = v8790
		v8940 = v8790
		goto L1042
	} else {
		goto L1043
	}
L1041:
	;
	v9123 = base.I64_extend_i32_u(v9056)
	v9124 = int32(2)
	if v9123 == int64(0) {
		goto L1056
	} else {
		goto L1057
	}
L1042:
	;
	v8978 = v8903
	v8989 = v7791 + v8940*int32(12)
	v8990 = v8789 - v8940
	goto L1048
L1043:
	;
	v8796 = v8789 & int32(-4)
	v8801 = v7791
	v8812 = v8796
	v8841 = base.Simd_g_const(&F_EncodeStreamHook__k0)
	goto L1044
L1044:
	;
	v8874 = int32(0)
	v8875 = base.Simd_g_v128_load32_splat(m, v8801, v8874)
	v8878 = base.Simd_g_v128_load32_lane_l1(m, v8801+int32(12), v8874, v8875)
	v8881 = base.Simd_g_v128_load32_lane_l2(m, v8801+int32(24), v8874, v8878)
	v8884 = base.Simd_g_v128_load32_lane_l3(m, v8801+int32(36), v8874, v8881)
	v8885 = base.Simd_g_i32x4_max_s(v8841, v8884)
	v8889 = v8812 + int32(-4)
	if v8889 != 0 {
		v8801 = v8801 + int32(48)
		v8812 = v8889
		v8841 = v8885
		goto L1044
	} else {
		goto L1046
	}
L1045:
	;
	v8892 = base.Simd_g_i32x4_max_s(v8885, base.Simd_g_i8x16_swizzle_c(v8885, base.Simd_g_const(&F_EncodeStreamHook__k5)))
	v8897 = base.Simd_g_i32x4_extract_lane_l0(base.Simd_g_i32x4_max_s(v8892, base.Simd_g_i8x16_swizzle_c(v8892, base.Simd_g_const(&F_EncodeStreamHook__k6))))
	if v8785 == v8796 {
		v9056 = v8897
		goto L1041
	} else {
		goto L1047
	}
L1046:
	;
	goto L1045
L1047:
	;
	v8903 = v8897
	v8940 = v8796
	goto L1042
L1048:
	;
	v9045 = *(*int32)(unsafe.Add(mBase, uint32(v8989)))
	if v9045 < v8978 {
		goto L1050
	} else {
		goto L1051
	}
L1049:
	;
	v9056 = v9047
	goto L1041
L1050:
	;
	v9047 = v8978
	goto L1052
L1051:
	;
	v9047 = v9045
	goto L1052
L1052:
	;
	v9051 = v8990 + int32(-1)
	if v9051 != 0 {
		v8978 = v9047
		v8989 = v8989 + int32(12)
		v8990 = v9051
		goto L1048
	} else {
		goto L1053
	}
L1053:
	;
	goto L1049
L1054:
	;
	if v9145 == int32(0) {
		goto L1004
	} else {
		goto L1060
	}
L1055:
	;
	goto L1054
L1056:
	;
	v9143 = F_malloc(m, base.I32_wrap_i64(v9123)*v9124)
	mBase = m.M
	v9145 = v9143
	goto L1055
L1057:
	;
	v9131 = base.I64_div_u_s(int64(2147418112), v9123)
	v9132 = int32(0)
	v9133 = base.I64_extend_i32_u(v9124)
	if base.Ui64(int64(4294967295)) < base.Ui64(v9133*v9123) {
		v9145 = v9132
		goto L1055
	} else {
		goto L1058
	}
L1058:
	;
	if base.Ui64(v9131) < base.Ui64(v9133) {
		v9145 = v9132
		goto L1055
	} else {
		goto L1059
	}
L1059:
	;
	goto L1056
L1060:
	;
	v9190 = int32(0)
	goto L1061
L1061:
	;
	v9223 = v7791 + v9190*int32(12)
	F_StoreHuffmanCode(m, v272, v7319, v9145, v9223)
	mBase = m.M
	v9225 = *(*int32)(unsafe.Add(mBase, uint32(v9223)))
	if v9225 < int32(1) {
		goto L1063
	} else {
		goto L1064
	}
L1063:
	;
	v9472 = v9190 + int32(1)
	if v9472 != v8789 {
		v9190 = v9472
		goto L1061
	} else {
		goto L1074
	}
L1064:
	;
	v9228 = *(*int32)(unsafe.Add(mBase, uint32(v9223)+4))
	v9245 = v9228
	v9246 = v9225
	v9271 = int32(0)
	goto L1065
L1065:
	;
	v9301 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v9245))))
	if v9301 == int32(0) {
		v9307 = v9271
		goto L1067
	} else {
		goto L1068
	}
L1066:
	;
	v9313 = int32(0)
	v9330 = v9313
	v9331 = v9313
	goto L1071
L1067:
	;
	v9312 = v9246 + int32(-1)
	if v9312 != 0 {
		v9245 = v9245 + int32(1)
		v9246 = v9312
		v9271 = v9307
		goto L1065
	} else {
		goto L1070
	}
L1068:
	;
	if int32(0) < v9271 {
		goto L1063
	} else {
		goto L1069
	}
L1069:
	;
	v9307 = int32(1)
	goto L1067
L1070:
	;
	goto L1066
L1071:
	;
	v9386 = *(*int32)(unsafe.Add(mBase, uint32(v9223)+4))
	v9388 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v9386+v9331))) = uint8(v9388)
	v9390 = *(*int32)(unsafe.Add(mBase, uint32(v9223)+8))
	*(*uint16)(unsafe.Add(mBase, uint32(v9390+v9330))) = uint16(v9388)
	v9397 = v9331 + int32(1)
	v9398 = *(*int32)(unsafe.Add(mBase, uint32(v9223)))
	if v9397 < v9398 {
		v9330 = v9330 + int32(2)
		v9331 = v9397
		goto L1071
	} else {
		goto L1073
	}
L1072:
	;
	goto L1063
L1073:
	;
	goto L1072
L1074:
	;
	v9875 = v9145
	v9890 = v8783
	goto L959
L1075:
	;
	if v9494 != 0 {
		v9875 = v9494
		v9890 = v8265
		goto L959
	} else {
		goto L1081
	}
L1076:
	;
	goto L1075
L1077:
	;
	v9494 = F_malloc(m, base.I32_wrap_i64(int64(0))*int32(2))
	mBase = m.M
	goto L1076
L1081:
	;
	goto L1004
L1082:
	;
	goto L1110
L1083:
	;
	goto L1082
L1084:
	;
	v9584 = *(*int32)(unsafe.Add(mBase, uint32(v272)))
	v9585 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	if v9585+v9572 < int32(32) {
		goto L1087
	} else {
		goto L1088
	}
L1085:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+4)) = v9686 + v9684
	*(*int32)(unsafe.Add(mBase, uint32(v272))) = v9683<<(uint(v9686)%32) | v9685
	goto L1083
L1086:
	;
	v9603 = *(*int32)(unsafe.Add(mBase, uint32(v272)+16))
	v9604 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v9608 = v9600
	v9609 = v9601
	v9610 = v9604
	v9611 = v9603
	goto L1090
L1087:
	;
	if v9585 < int32(16) {
		v9683 = v9570
		v9684 = v9572
		v9685 = v9584
		v9686 = v9585
		goto L1085
	} else {
		goto L1089
	}
L1088:
	;
	v9589 = int32(32)
	v9590 = v9589 - v9585
	v9598 = int32(base.Ui32(v9570) >> (uint(v9590) % 32))
	v9599 = v9572 - v9590
	v9600 = v9570<<(uint(v9585)%32) | v9584
	v9601 = v9589
	goto L1086
L1089:
	;
	v9598 = v9570
	v9599 = v9572
	v9600 = v9584
	v9601 = v9585
	goto L1086
L1090:
	;
	if base.Ui32(v9610+int32(2)) <= base.Ui32(v9611) {
		v9666 = v9610
		v9667 = v9611
		goto L1092
	} else {
		goto L1093
	}
L1091:
	;
	v9683 = v9598
	v9684 = v9599
	v9685 = v9679
	v9686 = v9677
	goto L1085
L1092:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v9666))) = uint16(v9608)
	v9674 = v9666 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v9674
	v9677 = v9609 + int32(-16)
	v9679 = int32(base.Ui32(v9608) >> (uint(int32(16)) % 32))
	if int32(31) < v9609 {
		v9608 = v9679
		v9609 = v9677
		v9610 = v9674
		v9611 = v9667
		goto L1090
	} else {
		goto L1107
	}
L1093:
	;
	v9620 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v9621 = v9611 - v9620
	v9624 = base.I64_extend_i32_s(v9621) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v9624) {
		v9648 = v9620
		goto L1095
	} else {
		goto L1096
	}
L1094:
	;
	if v9610 == v9620 {
		goto L1105
	} else {
		goto L1106
	}
L1095:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v272)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v272)+12)) = v9648
	goto L1082
L1096:
	;
	v9627 = v9610 - v9620
	v9629 = v9624 + base.I64_extend_i32_u(v9627)
	if base.Ui64(int64(4294967295)) < base.Ui64(v9629) {
		v9648 = v9620
		goto L1095
	} else {
		goto L1097
	}
L1097:
	;
	v9632 = base.I32_wrap_i64(v9629)
	if v9611 == v9620 {
		goto L1098
	} else {
		goto L1099
	}
L1098:
	;
	v9639 = int32(base.Ui32(v9621*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v9632) < base.Ui32(v9639) {
		goto L1101
	} else {
		goto L1102
	}
L1099:
	;
	if base.Ui32(v9632) <= base.Ui32(v9621) {
		v9666 = v9610
		v9667 = v9611
		goto L1092
	} else {
		goto L1100
	}
L1100:
	;
	goto L1098
L1101:
	;
	v9641 = v9639
	goto L1103
L1102:
	;
	v9641 = v9632
	goto L1103
L1103:
	;
	v9645 = v9641&int32(-1024) + int32(1024)
	v9646 = F_WebPSafeMalloc(m, int64(1), v9645)
	mBase = m.M
	if v9646 != 0 {
		goto L1094
	} else {
		goto L1104
	}
L1104:
	;
	v9647 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v9648 = v9647
	goto L1095
L1105:
	;
	v9659 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	F_WebPSafeFree(m, v9659)
	mBase = m.M
	v9661 = v9646 + v9645
	*(*int32)(unsafe.Add(mBase, uint32(v272)+16)) = v9661
	*(*int32)(unsafe.Add(mBase, uint32(v272)+8)) = v9646
	v9666 = v9646 + v9627
	v9667 = v9661
	goto L1092
L1106:
	;
	v9657 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v9658 = F_memcpy(m, v9646, v9657, v9627)
	mBase = m.M
	goto L1105
L1107:
	;
	goto L1091
L1108:
	;
	if v9743 != 0 {
		v9875 = v9743
		v9890 = v8265
		goto L959
	} else {
		goto L1114
	}
L1109:
	;
	goto L1108
L1110:
	;
	v9743 = F_malloc(m, base.I32_wrap_i64(int64(0))*int32(2))
	mBase = m.M
	goto L1109
L1114:
	;
	v9762 = v9570
	goto L960
L1115:
	;
	v10431 = v269
	v10433 = v271
	v10434 = v272
	v10435 = int32(0)
	v10436 = v274
	v10437 = v275
	v10438 = v276
	v10439 = v277
	v10440 = v278
	v10441 = v279
	v10442 = v280
	v10443 = v281
	v10445 = v283
	v10446 = v9762
	v10448 = v286
	v10449 = v287
	v10450 = v288
	v10451 = v289
	v10452 = v290
	v10453 = v291
	v10454 = v292
	v10455 = v293
	v10456 = v294
	v10457 = v295
	v10458 = v296
	v10459 = v297
	v10460 = v298
	v10461 = v299
	v10462 = v300
	v10463 = v301
	v10464 = v302
	v10465 = v303
	v10466 = v304
	v10467 = v305
	v10468 = v7649
	v10469 = v7650
	v10471 = v7703
	v10478 = v7319
	v10480 = v7791
	v10481 = v7352
	goto L786
L1116:
	;
	goto L1115
L1117:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L1116
L1118:
	;
	v9900 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	v9905 = *(*int32)(unsafe.Add(mBase, uint32(v272)+12))
	v9906 = *(*int32)(unsafe.Add(mBase, uint32(v272)+8))
	v9907 = v9905 - v9906
	v9908 = (v9900+int32(7))>>(uint(int32(3))%32) + v9907
	if base.Ui32(v7676) <= base.Ui32(v9908) {
		v9951 = v7649
		v9952 = v7650
		v9953 = v7676
		goto L1120
	} else {
		goto L1121
	}
L1119:
	;
	v10359 = v269
	v10361 = v271
	v10362 = v272
	v10364 = v274
	v10365 = v275
	v10366 = v276
	v10367 = v277
	v10368 = v278
	v10369 = v279
	v10370 = v280
	v10371 = v281
	v10373 = v283
	v10374 = int32(0)
	v10376 = v286
	v10377 = v287
	v10378 = v288
	v10379 = v289
	v10380 = v290
	v10381 = v291
	v10382 = v292
	v10383 = v293
	v10384 = v294
	v10385 = v295
	v10386 = v296
	v10387 = v297
	v10388 = v298
	v10389 = v299
	v10390 = v300
	v10391 = v301
	v10392 = v302
	v10393 = v303
	v10394 = v304
	v10395 = v305
	v10396 = v7649
	v10397 = v7650
	v10399 = v9875
	v10406 = v7319
	v10408 = v7791
	v10409 = v7352
	goto L787
L1120:
	;
	F_free(m, v9875)
	mBase = m.M
	goto L1123
L1121:
	;
	v9915 = (v9896+int32(7))>>(uint(int32(3))%32) + (v9895 - v9894)
	*(*int32)(unsafe.Add(mBase, uint32(v277)+48)) = v7707
	v9918 = *(*int32)(unsafe.Add(mBase, uint32(v272)+4))
	v9920 = v271 + int32(1088)
	v9923 = int32(0)
	v9927 = *(*int64)(unsafe.Add(mBase, uint32(v272)))
	v9928 = *(*int64)(unsafe.Add(mBase, uint32(v9920)))
	*(*int64)(unsafe.Add(mBase, uint32(v272))) = v9928
	v9931 = v272 + int32(8)
	v9933 = base.Simd_g_v128_load(m, v9931, v9923)
	v9935 = v271 + int32(1096)
	v9937 = base.Simd_g_v128_load(m, v9935, v9923)
	base.Simd_g_v128_store(m, v9931, v9923, v9937)
	*(*int64)(unsafe.Add(mBase, uint32(v9920))) = v9927
	base.Simd_g_v128_store(m, v9935, v9923, v9933)
	goto L1122
L1122:
	;
	v9951 = v9907 - v9915 + (v9918+int32(7))>>(uint(int32(3))%32)
	v9952 = v9915 - v287
	v9953 = v9908
	goto L1120
L1123:
	;
	v9955 = *(*int32)(unsafe.Add(mBase, uint32(v7791)+8))
	F_free(m, v9955)
	mBase = m.M
	goto L1124
L1124:
	;
	F_free(m, v7791)
	mBase = m.M
	goto L1125
L1125:
	;
	v9958 = *(*int32)(unsafe.Add(mBase, uint32(v7596)+4))
	v9959 = int32(0)
	if v7680&base.B2i32(v9958 != v9959) != 0 {
		v7628 = int32(1)
		v7649 = v9951
		v7650 = v9952
		v7676 = v9953
		v7679 = v9890
		v7680 = v9959
		goto L840
	} else {
		goto L1126
	}
L1126:
	;
	goto L841
L1127:
	;
	goto L837
L1128:
	;
	v10137 = F_WebPReportProgress(m, v278, v7296+v7140, v271+int32(60))
	mBase = m.M
	v10179 = v10076
	v10180 = v10077
	goto L809
L1129:
	;
	v10179 = v306
	v10180 = v307
	goto L809
L1130:
	;
	goto L1129
L1131:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v278)+92)) = int32(1)
	goto L1130
L1132:
	;
	F_free(m, v10478)
	mBase = m.M
	goto L1133
L1133:
	;
	F_WebPSafeFree(m, v10435)
	mBase = m.M
	goto L1134
L1134:
	;
	F_WebPSafeFree(m, v10446)
	mBase = m.M
	goto L1135
L1135:
	;
	v10507 = v10433 + int32(64)
	v10508 = *(*int32)(unsafe.Add(mBase, uint32(v10507)))
	F_WebPSafeFree(m, v10508)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v10507))) = int64(0)
	goto L1136
L1136:
	;
	if v10480 == int32(0) {
		goto L1137
	} else {
		goto L1138
	}
L1137:
	;
	F_free(m, v10481)
	mBase = m.M
	goto L1141
L1138:
	;
	v10514 = *(*int32)(unsafe.Add(mBase, uint32(v10480)+8))
	F_free(m, v10514)
	mBase = m.M
	goto L1139
L1139:
	;
	F_free(m, v10480)
	mBase = m.M
	goto L1140
L1140:
	;
	goto L1137
L1141:
	;
	v10519 = v10433 + int32(1088)
	if v10519 == int32(0) {
		goto L1143
	} else {
		goto L1144
	}
L1142:
	;
	v10533 = *(*int32)(unsafe.Add(mBase, uint32(v10440)+92))
	if v10533 != 0 {
		v10723 = v10431
		v10725 = v10433
		goto L4
	} else {
		goto L1145
	}
L1143:
	;
	goto L1142
L1144:
	;
	v10524 = *(*int32)(unsafe.Add(mBase, uint32(v10433+int32(1096))))
	F_WebPSafeFree(m, v10524)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v10433+int32(1104)))) = int64(0)
	v10530 = base.Simd_g_const(&F_EncodeStreamHook__k0)
	base.Simd_g_v128_store(m, v10519, int32(0), v10530)
	goto L1143
L1145:
	;
	v10534 = *(*int32)(unsafe.Add(mBase, uint32(v10434)+4))
	v10539 = *(*int32)(unsafe.Add(mBase, uint32(v10434)+12))
	v10540 = *(*int32)(unsafe.Add(mBase, uint32(v10434)+8))
	v10542 = (v10534+int32(7))>>(uint(int32(3))%32) + (v10539 - v10540)
	if base.Ui32(v10466) <= base.Ui32(v10542) {
		v10610 = v10466
		goto L1146
	} else {
		goto L1147
	}
L1146:
	;
	if v10438 < int32(2) {
		goto L1157
	} else {
		goto L1158
	}
L1147:
	;
	v10544 = int32(8)
	v10545 = v10433 + v10544
	v10548 = int32(0)
	v10552 = *(*int64)(unsafe.Add(mBase, uint32(v10434)))
	v10553 = *(*int64)(unsafe.Add(mBase, uint32(v10545)))
	*(*int64)(unsafe.Add(mBase, uint32(v10434))) = v10553
	v10556 = v10434 + v10544
	v10558 = base.Simd_g_v128_load(m, v10556, v10548)
	v10560 = v10433 + int32(16)
	v10562 = base.Simd_g_v128_load(m, v10560, v10548)
	base.Simd_g_v128_store(m, v10556, v10548, v10562)
	*(*int64)(unsafe.Add(mBase, uint32(v10545))) = v10552
	base.Simd_g_v128_store(m, v10560, v10548, v10558)
	goto L1148
L1148:
	;
	if v10436 == int32(0) {
		goto L1149
	} else {
		goto L1150
	}
L1149:
	;
	v10610 = v10542
	goto L1146
L1150:
	;
	v10570 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+60))
	v10571 = int32(0)
	v10572 = base.B2i32(v10570 != v10571)
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+148)) = v10572
	v10574 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+52))
	if v10574 == v10571 {
		v10580 = v10572
		goto L1151
	} else {
		goto L1152
	}
L1151:
	;
	v10581 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+56))
	if v10581 == int32(0) {
		v10587 = v10580
		goto L1153
	} else {
		goto L1154
	}
L1152:
	;
	v10578 = v10572 | int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+148)) = v10578
	v10580 = v10578
	goto L1151
L1153:
	;
	v10588 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+64))
	if v10588 == int32(0) {
		goto L1155
	} else {
		goto L1156
	}
L1154:
	;
	v10585 = v10580 | int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+148)) = v10585
	v10587 = v10585
	goto L1153
L1155:
	;
	v10594 = *(*int32)(unsafe.Add(mBase, uint32(v10433)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+156)) = v10594
	v10596 = *(*int32)(unsafe.Add(mBase, uint32(v10433)))
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+180)) = v10596
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+176)) = v10468
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+172)) = v10469
	v10600 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+36))
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+152)) = v10600
	v10602 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+48))
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+160)) = v10602
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+168)) = v10542 - v10449
	v10606 = *(*int32)(unsafe.Add(mBase, uint32(v10439)+68))
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+164)) = v10606
	goto L1149
L1156:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v10436)+148)) = v10587 | int32(8)
	goto L1155
L1157:
	;
	v10626 = v10467 + int32(1)
	if v10626 != v10438 {
		v269 = v10431
		v271 = v10433
		v272 = v10434
		v274 = v10436
		v275 = v10437
		v276 = v10438
		v277 = v10439
		v278 = v10440
		v279 = v10441
		v280 = v10442
		v281 = v10443
		v283 = v10445
		v286 = v10448
		v287 = v10449
		v288 = v10450
		v289 = v10451
		v290 = v10452
		v291 = v10453
		v292 = v10454
		v293 = v10455
		v294 = v10456
		v295 = v10457
		v296 = v10458
		v297 = v10459
		v298 = v10460
		v299 = v10461
		v300 = v10462
		v301 = v10463
		v302 = v10464
		v303 = v10465
		v304 = v10610
		v305 = v10626
		v306 = v10468
		v307 = v10469
		goto L33
	} else {
		goto L1160
	}
L1158:
	;
	v10614 = v10433 + int32(32)
	v10615 = *(*int64)(unsafe.Add(mBase, uint32(v10614)))
	*(*int64)(unsafe.Add(mBase, uint32(v10434))) = v10615
	v10617 = *(*int32)(unsafe.Add(mBase, uint32(v10614)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v10434)+20)) = v10617
	v10619 = *(*int32)(unsafe.Add(mBase, uint32(v10434)+8))
	v10620 = *(*int32)(unsafe.Add(mBase, uint32(v10614)+12))
	v10621 = *(*int32)(unsafe.Add(mBase, uint32(v10614)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v10434)+12)) = v10619 + (v10620 - v10621)
	goto L1159
L1159:
	;
	goto L1157
L1160:
	;
	goto L34
L1161:
	;
	v10723 = v10628
	v10725 = v10630
	goto L4
L1162:
	;
	v10809 = *(*int32)(unsafe.Add(mBase, uint32(v10723)+4))
	v10810 = *(*int32)(unsafe.Add(mBase, uint32(v10809)+92))
	m.G0 = v10725 + int32(_a_F_EncodeStreamHook_0)
	return base.B2i32(v10810 == int32(0))
L1163:
	;
	goto L1162
L1164:
	;
	v10800 = *(*int32)(unsafe.Add(mBase, uint32(v10725+int32(16))))
	F_WebPSafeFree(m, v10800)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v10725+int32(24)))) = int64(0)
	v10806 = base.Simd_g_const(&F_EncodeStreamHook__k0)
	base.Simd_g_v128_store(m, v10795, int32(0), v10806)
	goto L1163
}

var F_EncodeStreamHook__k0 = [2]uint64{0x0, 0x0}
var F_EncodeStreamHook__k1 = [2]uint64{0xff00ff00ff00ff00, 0xff00ff00ff00ff00}
var F_EncodeStreamHook__k2 = [2]uint64{0xff00ff00ff00ff, 0xff00ff00ff00ff}
var F_EncodeStreamHook__k3 = [2]uint64{0xff000000ff0000, 0xff000000ff0000}
var F_EncodeStreamHook__k4 = [2]uint64{0xff000000ff00, 0xff000000ff00}
var F_EncodeStreamHook__k5 = [2]uint64{0xf0e0d0c0b0a0908, 0x302010003020100}
var F_EncodeStreamHook__k6 = [2]uint64{0x302010007060504, 0x302010003020100}

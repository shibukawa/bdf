//go:build !bdf_noconv

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpw/base"
	"unsafe"
)

func F_EncodeStreamHook(m *base.Module, l0 int32, l1 int32) int32 {
	mBase := m.M
	_ = mBase
	var v67 int32
	_ = v67
	var v69 int32
	_ = v69
	var v71 int32
	_ = v71
	var v73 int32
	_ = v73
	var v74 int32
	_ = v74
	var v75 int32
	_ = v75
	var v76 int32
	_ = v76
	var v77 int32
	_ = v77
	var v78 int32
	_ = v78
	var v79 int32
	_ = v79
	var v80 int32
	_ = v80
	var v81 int32
	_ = v81
	var v82 int32
	_ = v82
	var v83 float32
	_ = v83
	var v84 int32
	_ = v84
	var v85 int32
	_ = v85
	var v86 int32
	_ = v86
	var v95 int64
	_ = v95
	var v101 int64
	_ = v101
	var v103 int64
	_ = v103
	var v110 int32
	_ = v110
	var v112 int32
	_ = v112
	var v114 int32
	_ = v114
	var v117 int64
	_ = v117
	var v130 int32
	_ = v130
	var v132 int32
	_ = v132
	var v136 int32
	_ = v136
	var v143 int32
	_ = v143
	var v149 int32
	_ = v149
	var v159 int32
	_ = v159
	var v160 int32
	_ = v160
	var v161 int32
	_ = v161
	var v163 int32
	_ = v163
	var v164 int32
	_ = v164
	var v165 int32
	_ = v165
	var v167 int64
	_ = v167
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v175 int32
	_ = v175
	var v182 int32
	_ = v182
	var v184 int32
	_ = v184
	var v188 int32
	_ = v188
	var v189 int32
	_ = v189
	var v194 int32
	_ = v194
	var v195 int32
	_ = v195
	var v196 int32
	_ = v196
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v205 int32
	_ = v205
	var v206 int32
	_ = v206
	var v207 int64
	_ = v207
	var v209 int32
	_ = v209
	var v211 int32
	_ = v211
	var v219 int32
	_ = v219
	var v221 int32
	_ = v221
	var v234 int32
	_ = v234
	var v237 int32
	_ = v237
	var v239 int32
	_ = v239
	var v245 int32
	_ = v245
	var v265 int32
	_ = v265
	var v267 int32
	_ = v267
	var v268 int32
	_ = v268
	var v271 int32
	_ = v271
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
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
	var v280 int32
	_ = v280
	var v283 int32
	_ = v283
	var v284 int32
	_ = v284
	var v285 int32
	_ = v285
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
	var v333 int32
	_ = v333
	var v334 int32
	_ = v334
	var v335 int32
	_ = v335
	var v338 int32
	_ = v338
	var v340 int32
	_ = v340
	var v348 int32
	_ = v348
	var v357 int32
	_ = v357
	var v360 int32
	_ = v360
	var v362 int32
	_ = v362
	var v364 int32
	_ = v364
	var v368 int32
	_ = v368
	var v371 int32
	_ = v371
	var v373 int32
	_ = v373
	var v375 int32
	_ = v375
	var v386 int32
	_ = v386
	var v387 int32
	_ = v387
	var v395 int32
	_ = v395
	var v398 int32
	_ = v398
	var v400 int32
	_ = v400
	var v402 int32
	_ = v402
	var v413 int32
	_ = v413
	var v414 int32
	_ = v414
	var v421 int32
	_ = v421
	var v424 int32
	_ = v424
	var v429 int32
	_ = v429
	var v440 int64
	_ = v440
	var v441 int32
	_ = v441
	var v445 int32
	_ = v445
	var v448 int32
	_ = v448
	var v452 int32
	_ = v452
	var v454 int32
	_ = v454
	var v455 int64
	_ = v455
	var v457 int32
	_ = v457
	var v460 int32
	_ = v460
	var v461 int64
	_ = v461
	var v463 int32
	_ = v463
	var v465 int32
	_ = v465
	var v474 int32
	_ = v474
	var v475 int64
	_ = v475
	var v476 int64
	_ = v476
	var v480 int64
	_ = v480
	var v481 int32
	_ = v481
	var v484 int64
	_ = v484
	var v490 int32
	_ = v490
	var v491 int32
	_ = v491
	var v493 int32
	_ = v493
	var v498 int32
	_ = v498
	var v502 int32
	_ = v502
	var v505 int32
	_ = v505
	var v507 int32
	_ = v507
	var v508 int32
	_ = v508
	var v520 int32
	_ = v520
	var v523 int32
	_ = v523
	var v524 int32
	_ = v524
	var v527 int32
	_ = v527
	var v528 int32
	_ = v528
	var v537 int32
	_ = v537
	var v538 int32
	_ = v538
	var v539 int32
	_ = v539
	var v540 int32
	_ = v540
	var v545 int32
	_ = v545
	var v547 int32
	_ = v547
	var v556 int32
	_ = v556
	var v559 int32
	_ = v559
	var v571 int32
	_ = v571
	var v573 int32
	_ = v573
	var v578 int32
	_ = v578
	var v579 int32
	_ = v579
	var v581 int32
	_ = v581
	var v584 int32
	_ = v584
	var v586 int32
	_ = v586
	var v587 int32
	_ = v587
	var v594 int32
	_ = v594
	var v597 int32
	_ = v597
	var v602 int32
	_ = v602
	var v612 int32
	_ = v612
	var v615 int32
	_ = v615
	var v616 int32
	_ = v616
	var v621 int32
	_ = v621
	var v622 int32
	_ = v622
	var v627 int32
	_ = v627
	var v630 int32
	_ = v630
	var v642 int32
	_ = v642
	var v665 int32
	_ = v665
	var v667 int32
	_ = v667
	var v670 int32
	_ = v670
	var v672 int32
	_ = v672
	var v676 int32
	_ = v676
	var v677 int32
	_ = v677
	var v678 int32
	_ = v678
	var v679 int32
	_ = v679
	var v680 int32
	_ = v680
	var v682 int32
	_ = v682
	var v685 int32
	_ = v685
	var v687 int32
	_ = v687
	var v691 int32
	_ = v691
	var v692 int32
	_ = v692
	var v705 int32
	_ = v705
	var v706 int32
	_ = v706
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v719 int32
	_ = v719
	var v720 int32
	_ = v720
	var v721 int32
	_ = v721
	var v722 int32
	_ = v722
	var v724 int32
	_ = v724
	var v725 int32
	_ = v725
	var v729 int32
	_ = v729
	var v730 int32
	_ = v730
	var v731 int32
	_ = v731
	var v732 int32
	_ = v732
	var v741 int32
	_ = v741
	var v742 int32
	_ = v742
	var v745 int64
	_ = v745
	var v748 int32
	_ = v748
	var v750 int64
	_ = v750
	var v753 int32
	_ = v753
	var v760 int32
	_ = v760
	var v762 int32
	_ = v762
	var v766 int32
	_ = v766
	var v767 int32
	_ = v767
	var v768 int32
	_ = v768
	var v769 int32
	_ = v769
	var v778 int32
	_ = v778
	var v779 int32
	_ = v779
	var v780 int32
	_ = v780
	var v782 int32
	_ = v782
	var v787 int32
	_ = v787
	var v788 int32
	_ = v788
	var v795 int32
	_ = v795
	var v798 int32
	_ = v798
	var v800 int32
	_ = v800
	var v804 int32
	_ = v804
	var v805 int32
	_ = v805
	var v806 int32
	_ = v806
	var v807 int32
	_ = v807
	var v844 int32
	_ = v844
	var v845 int32
	_ = v845
	var v857 int32
	_ = v857
	var v858 int32
	_ = v858
	var v862 int32
	_ = v862
	var v863 int32
	_ = v863
	var v871 int32
	_ = v871
	var v872 int32
	_ = v872
	var v873 int32
	_ = v873
	var v874 int32
	_ = v874
	var v876 int32
	_ = v876
	var v877 int32
	_ = v877
	var v881 int32
	_ = v881
	var v882 int32
	_ = v882
	var v883 int32
	_ = v883
	var v884 int32
	_ = v884
	var v893 int32
	_ = v893
	var v894 int32
	_ = v894
	var v897 int64
	_ = v897
	var v900 int32
	_ = v900
	var v902 int64
	_ = v902
	var v905 int32
	_ = v905
	var v912 int32
	_ = v912
	var v914 int32
	_ = v914
	var v918 int32
	_ = v918
	var v919 int32
	_ = v919
	var v920 int32
	_ = v920
	var v921 int32
	_ = v921
	var v930 int32
	_ = v930
	var v931 int32
	_ = v931
	var v932 int32
	_ = v932
	var v934 int32
	_ = v934
	var v939 int32
	_ = v939
	var v940 int32
	_ = v940
	var v947 int32
	_ = v947
	var v950 int32
	_ = v950
	var v952 int32
	_ = v952
	var v956 int32
	_ = v956
	var v957 int32
	_ = v957
	var v958 int32
	_ = v958
	var v959 int32
	_ = v959
	var v996 int32
	_ = v996
	var v999 int32
	_ = v999
	var v1001 int32
	_ = v1001
	var v1002 int32
	_ = v1002
	var v1014 int32
	_ = v1014
	var v1015 int32
	_ = v1015
	var v1019 int32
	_ = v1019
	var v1020 int32
	_ = v1020
	var v1028 int32
	_ = v1028
	var v1029 int32
	_ = v1029
	var v1030 int32
	_ = v1030
	var v1031 int32
	_ = v1031
	var v1033 int32
	_ = v1033
	var v1034 int32
	_ = v1034
	var v1038 int32
	_ = v1038
	var v1039 int32
	_ = v1039
	var v1040 int32
	_ = v1040
	var v1041 int32
	_ = v1041
	var v1050 int32
	_ = v1050
	var v1051 int32
	_ = v1051
	var v1054 int64
	_ = v1054
	var v1057 int32
	_ = v1057
	var v1059 int64
	_ = v1059
	var v1062 int32
	_ = v1062
	var v1069 int32
	_ = v1069
	var v1071 int32
	_ = v1071
	var v1075 int32
	_ = v1075
	var v1076 int32
	_ = v1076
	var v1077 int32
	_ = v1077
	var v1078 int32
	_ = v1078
	var v1087 int32
	_ = v1087
	var v1088 int32
	_ = v1088
	var v1089 int32
	_ = v1089
	var v1091 int32
	_ = v1091
	var v1096 int32
	_ = v1096
	var v1097 int32
	_ = v1097
	var v1104 int32
	_ = v1104
	var v1107 int32
	_ = v1107
	var v1109 int32
	_ = v1109
	var v1113 int32
	_ = v1113
	var v1114 int32
	_ = v1114
	var v1115 int32
	_ = v1115
	var v1116 int32
	_ = v1116
	var v1156 int32
	_ = v1156
	var v1163 int32
	_ = v1163
	var v1164 int32
	_ = v1164
	var v1175 int32
	_ = v1175
	var v1227 int32
	_ = v1227
	var v1230 int32
	_ = v1230
	var v1231 int32
	_ = v1231
	var v1238 int32
	_ = v1238
	var v1245 int32
	_ = v1245
	var v1250 int32
	_ = v1250
	var v1321 int32
	_ = v1321
	var v1327 int32
	_ = v1327
	var v1330 int32
	_ = v1330
	var v1333 int32
	_ = v1333
	var v1334 int32
	_ = v1334
	var v1336 int32
	_ = v1336
	var v1338 int32
	_ = v1338
	var v1341 int32
	_ = v1341
	var v1346 int32
	_ = v1346
	var v1351 int32
	_ = v1351
	var v1352 int32
	_ = v1352
	var v1359 int64
	_ = v1359
	var v1360 int32
	_ = v1360
	var v1364 int32
	_ = v1364
	var v1367 int32
	_ = v1367
	var v1371 int32
	_ = v1371
	var v1373 int32
	_ = v1373
	var v1374 int64
	_ = v1374
	var v1376 int32
	_ = v1376
	var v1379 int32
	_ = v1379
	var v1380 int64
	_ = v1380
	var v1382 int32
	_ = v1382
	var v1384 int32
	_ = v1384
	var v1393 int32
	_ = v1393
	var v1394 int64
	_ = v1394
	var v1395 int64
	_ = v1395
	var v1399 int64
	_ = v1399
	var v1400 int32
	_ = v1400
	var v1403 int64
	_ = v1403
	var v1409 int32
	_ = v1409
	var v1410 int32
	_ = v1410
	var v1412 int32
	_ = v1412
	var v1417 int32
	_ = v1417
	var v1421 int32
	_ = v1421
	var v1424 int32
	_ = v1424
	var v1426 int32
	_ = v1426
	var v1427 int32
	_ = v1427
	var v1439 int32
	_ = v1439
	var v1442 int32
	_ = v1442
	var v1443 int32
	_ = v1443
	var v1444 int32
	_ = v1444
	var v1445 int32
	_ = v1445
	var v1446 int64
	_ = v1446
	var v1447 int32
	_ = v1447
	var v1454 int64
	_ = v1454
	var v1455 int32
	_ = v1455
	var v1456 int64
	_ = v1456
	var v1466 int32
	_ = v1466
	var v1468 int32
	_ = v1468
	var v1478 int32
	_ = v1478
	var v1479 int32
	_ = v1479
	var v1520 int32
	_ = v1520
	var v1521 int32
	_ = v1521
	var v1526 int32
	_ = v1526
	var v1527 int32
	_ = v1527
	var v1530 int32
	_ = v1530
	var v1553 int32
	_ = v1553
	var v1554 int32
	_ = v1554
	var v1566 int32
	_ = v1566
	var v1588 int32
	_ = v1588
	var v1589 int32
	_ = v1589
	var v1615 int32
	_ = v1615
	var v1617 int32
	_ = v1617
	var v1620 int32
	_ = v1620
	var v1625 int32
	_ = v1625
	var v1627 int32
	_ = v1627
	var v1628 int32
	_ = v1628
	var v1629 int32
	_ = v1629
	var v1630 int32
	_ = v1630
	var v1637 int32
	_ = v1637
	var v1677 int32
	_ = v1677
	var v1678 int32
	_ = v1678
	var v1704 int32
	_ = v1704
	var v1705 int32
	_ = v1705
	var v1712 int32
	_ = v1712
	var v1843 int32
	_ = v1843
	var v1844 int32
	_ = v1844
	var v1907 int32
	_ = v1907
	var v1910 int32
	_ = v1910
	var v1911 int32
	_ = v1911
	var v2043 int32
	_ = v2043
	var v2044 int32
	_ = v2044
	var v2107 int32
	_ = v2107
	var v2116 int32
	_ = v2116
	var v2117 int32
	_ = v2117
	var v2249 int32
	_ = v2249
	var v2250 int32
	_ = v2250
	var v2313 int32
	_ = v2313
	var v2322 int32
	_ = v2322
	var v2323 int32
	_ = v2323
	var v2330 int32
	_ = v2330
	var v2336 int32
	_ = v2336
	var v2350 int32
	_ = v2350
	var v2391 int32
	_ = v2391
	var v2412 int32
	_ = v2412
	var v2413 int32
	_ = v2413
	var v2416 int32
	_ = v2416
	var v2418 int32
	_ = v2418
	var v2420 int32
	_ = v2420
	var v2466 int32
	_ = v2466
	var v2489 int32
	_ = v2489
	var v2490 int32
	_ = v2490
	var v2497 int32
	_ = v2497
	var v2503 int32
	_ = v2503
	var v2517 int32
	_ = v2517
	var v2558 int32
	_ = v2558
	var v2579 int32
	_ = v2579
	var v2580 int32
	_ = v2580
	var v2583 int32
	_ = v2583
	var v2585 int32
	_ = v2585
	var v2587 int32
	_ = v2587
	var v2633 int32
	_ = v2633
	var v2656 int32
	_ = v2656
	var v2657 int32
	_ = v2657
	var v2666 int32
	_ = v2666
	var v2667 int32
	_ = v2667
	var v2685 int32
	_ = v2685
	var v2686 int32
	_ = v2686
	var v2714 int32
	_ = v2714
	var v2715 int32
	_ = v2715
	var v2717 int32
	_ = v2717
	var v2743 int32
	_ = v2743
	var v2744 int32
	_ = v2744
	var v2755 int32
	_ = v2755
	var v2756 int32
	_ = v2756
	var v2805 int32
	_ = v2805
	var v2818 int32
	_ = v2818
	var v2819 int32
	_ = v2819
	var v2820 int32
	_ = v2820
	var v2821 int32
	_ = v2821
	var v2825 int32
	_ = v2825
	var v2838 int32
	_ = v2838
	var v2839 int32
	_ = v2839
	var v2840 int32
	_ = v2840
	var v2847 int32
	_ = v2847
	var v2854 int32
	_ = v2854
	var v2865 int32
	_ = v2865
	var v2866 int32
	_ = v2866
	var v2920 int32
	_ = v2920
	var v2933 int32
	_ = v2933
	var v2934 int32
	_ = v2934
	var v2935 int32
	_ = v2935
	var v2939 int32
	_ = v2939
	var v2940 int32
	_ = v2940
	var v2941 int32
	_ = v2941
	var v2942 int32
	_ = v2942
	var v2949 int32
	_ = v2949
	var v2957 int32
	_ = v2957
	var v2958 int32
	_ = v2958
	var v2976 int32
	_ = v2976
	var v2977 int32
	_ = v2977
	var v3005 int32
	_ = v3005
	var v3006 int32
	_ = v3006
	var v3008 int32
	_ = v3008
	var v3034 int32
	_ = v3034
	var v3035 int32
	_ = v3035
	var v3046 int32
	_ = v3046
	var v3047 int32
	_ = v3047
	var v3096 int32
	_ = v3096
	var v3109 int32
	_ = v3109
	var v3110 int32
	_ = v3110
	var v3111 int32
	_ = v3111
	var v3112 int32
	_ = v3112
	var v3116 int32
	_ = v3116
	var v3129 int32
	_ = v3129
	var v3130 int32
	_ = v3130
	var v3131 int32
	_ = v3131
	var v3138 int32
	_ = v3138
	var v3145 int32
	_ = v3145
	var v3156 int32
	_ = v3156
	var v3157 int32
	_ = v3157
	var v3211 int32
	_ = v3211
	var v3224 int32
	_ = v3224
	var v3225 int32
	_ = v3225
	var v3226 int32
	_ = v3226
	var v3230 int32
	_ = v3230
	var v3231 int32
	_ = v3231
	var v3232 int32
	_ = v3232
	var v3233 int32
	_ = v3233
	var v3240 int32
	_ = v3240
	var v3255 int32
	_ = v3255
	var v3257 int32
	_ = v3257
	var v3269 int32
	_ = v3269
	var v3276 int32
	_ = v3276
	var v3280 int32
	_ = v3280
	var v3290 int32
	_ = v3290
	var v3291 int32
	_ = v3291
	var v3296 int32
	_ = v3296
	var v3300 int32
	_ = v3300
	var v3301 int32
	_ = v3301
	var v3302 int32
	_ = v3302
	var v3303 int32
	_ = v3303
	var v3307 int32
	_ = v3307
	var v3321 int32
	_ = v3321
	var v3338 int32
	_ = v3338
	var v3352 int32
	_ = v3352
	var v3393 int32
	_ = v3393
	var v3414 int32
	_ = v3414
	var v3415 int32
	_ = v3415
	var v3418 int32
	_ = v3418
	var v3420 int32
	_ = v3420
	var v3422 int32
	_ = v3422
	var v3468 int32
	_ = v3468
	var v3491 int32
	_ = v3491
	var v3492 int32
	_ = v3492
	var v3497 int32
	_ = v3497
	var v3498 int32
	_ = v3498
	var v3517 int32
	_ = v3517
	var v3540 int32
	_ = v3540
	var v3545 int32
	_ = v3545
	var v3546 int32
	_ = v3546
	var v3549 int32
	_ = v3549
	var v3571 int32
	_ = v3571
	var v3572 int32
	_ = v3572
	var v3584 int32
	_ = v3584
	var v3607 int32
	_ = v3607
	var v3633 int32
	_ = v3633
	var v3638 int32
	_ = v3638
	var v3644 int32
	_ = v3644
	var v3649 int32
	_ = v3649
	var v3651 int32
	_ = v3651
	var v3656 int32
	_ = v3656
	var v3660 int32
	_ = v3660
	var v3661 int32
	_ = v3661
	var v3662 int32
	_ = v3662
	var v3663 int32
	_ = v3663
	var v3668 int32
	_ = v3668
	var v3675 int32
	_ = v3675
	var v3676 int32
	_ = v3676
	var v3677 int32
	_ = v3677
	var v3683 int32
	_ = v3683
	var v3685 int32
	_ = v3685
	var v3686 int32
	_ = v3686
	var v3693 int32
	_ = v3693
	var v3699 int32
	_ = v3699
	var v3705 int32
	_ = v3705
	var v3719 int32
	_ = v3719
	var v3760 int32
	_ = v3760
	var v3781 int32
	_ = v3781
	var v3782 int32
	_ = v3782
	var v3785 int32
	_ = v3785
	var v3787 int32
	_ = v3787
	var v3789 int32
	_ = v3789
	var v3835 int32
	_ = v3835
	var v3858 int32
	_ = v3858
	var v3859 int32
	_ = v3859
	var v3868 int32
	_ = v3868
	var v3869 int32
	_ = v3869
	var v3887 int32
	_ = v3887
	var v3888 int32
	_ = v3888
	var v3916 int32
	_ = v3916
	var v3917 int32
	_ = v3917
	var v3919 int32
	_ = v3919
	var v3945 int32
	_ = v3945
	var v3946 int32
	_ = v3946
	var v3957 int32
	_ = v3957
	var v3958 int32
	_ = v3958
	var v4007 int32
	_ = v4007
	var v4016 int32
	_ = v4016
	var v4017 int32
	_ = v4017
	var v4018 int32
	_ = v4018
	var v4019 int32
	_ = v4019
	var v4023 int32
	_ = v4023
	var v4032 int32
	_ = v4032
	var v4033 int32
	_ = v4033
	var v4034 int32
	_ = v4034
	var v4041 int32
	_ = v4041
	var v4048 int32
	_ = v4048
	var v4059 int32
	_ = v4059
	var v4060 int32
	_ = v4060
	var v4114 int32
	_ = v4114
	var v4123 int32
	_ = v4123
	var v4124 int32
	_ = v4124
	var v4125 int32
	_ = v4125
	var v4129 int32
	_ = v4129
	var v4130 int32
	_ = v4130
	var v4131 int32
	_ = v4131
	var v4132 int32
	_ = v4132
	var v4139 int32
	_ = v4139
	var v4212 int32
	_ = v4212
	var v4213 int32
	_ = v4213
	var v4214 int32
	_ = v4214
	var v4221 int64
	_ = v4221
	var v4222 int32
	_ = v4222
	var v4226 int32
	_ = v4226
	var v4229 int32
	_ = v4229
	var v4233 int32
	_ = v4233
	var v4235 int32
	_ = v4235
	var v4236 int64
	_ = v4236
	var v4238 int32
	_ = v4238
	var v4241 int32
	_ = v4241
	var v4242 int64
	_ = v4242
	var v4244 int32
	_ = v4244
	var v4246 int32
	_ = v4246
	var v4255 int32
	_ = v4255
	var v4256 int64
	_ = v4256
	var v4257 int64
	_ = v4257
	var v4261 int64
	_ = v4261
	var v4262 int32
	_ = v4262
	var v4265 int64
	_ = v4265
	var v4271 int32
	_ = v4271
	var v4272 int32
	_ = v4272
	var v4274 int32
	_ = v4274
	var v4279 int32
	_ = v4279
	var v4283 int32
	_ = v4283
	var v4286 int32
	_ = v4286
	var v4288 int32
	_ = v4288
	var v4289 int32
	_ = v4289
	var v4301 int32
	_ = v4301
	var v4304 int32
	_ = v4304
	var v4309 int32
	_ = v4309
	var v4312 int32
	_ = v4312
	var v4313 int32
	_ = v4313
	var v4314 int32
	_ = v4314
	var v4325 int32
	_ = v4325
	var v4338 int32
	_ = v4338
	var v4361 int32
	_ = v4361
	var v4387 int32
	_ = v4387
	var v4389 int32
	_ = v4389
	var v4390 int32
	_ = v4390
	var v4392 int32
	_ = v4392
	var v4393 int32
	_ = v4393
	var v4394 int32
	_ = v4394
	var v4397 int32
	_ = v4397
	var v4398 int32
	_ = v4398
	var v4400 int32
	_ = v4400
	var v4405 int32
	_ = v4405
	var v4418 int32
	_ = v4418
	var v4469 int32
	_ = v4469
	var v4539 int32
	_ = v4539
	var v4612 int32
	_ = v4612
	var v4668 int32
	_ = v4668
	var v4687 int32
	_ = v4687
	var v4690 int32
	_ = v4690
	var v4691 int32
	_ = v4691
	var v4704 int32
	_ = v4704
	var v4705 int32
	_ = v4705
	var v4709 int32
	_ = v4709
	var v4710 int32
	_ = v4710
	var v4718 int32
	_ = v4718
	var v4719 int32
	_ = v4719
	var v4720 int32
	_ = v4720
	var v4721 int32
	_ = v4721
	var v4723 int32
	_ = v4723
	var v4724 int32
	_ = v4724
	var v4728 int32
	_ = v4728
	var v4729 int32
	_ = v4729
	var v4730 int32
	_ = v4730
	var v4731 int32
	_ = v4731
	var v4740 int32
	_ = v4740
	var v4741 int32
	_ = v4741
	var v4744 int64
	_ = v4744
	var v4747 int32
	_ = v4747
	var v4749 int64
	_ = v4749
	var v4752 int32
	_ = v4752
	var v4759 int32
	_ = v4759
	var v4761 int32
	_ = v4761
	var v4765 int32
	_ = v4765
	var v4766 int32
	_ = v4766
	var v4767 int32
	_ = v4767
	var v4768 int32
	_ = v4768
	var v4777 int32
	_ = v4777
	var v4778 int32
	_ = v4778
	var v4779 int32
	_ = v4779
	var v4781 int32
	_ = v4781
	var v4786 int32
	_ = v4786
	var v4787 int32
	_ = v4787
	var v4794 int32
	_ = v4794
	var v4797 int32
	_ = v4797
	var v4799 int32
	_ = v4799
	var v4803 int32
	_ = v4803
	var v4804 int32
	_ = v4804
	var v4805 int32
	_ = v4805
	var v4806 int32
	_ = v4806
	var v4843 int32
	_ = v4843
	var v4856 int32
	_ = v4856
	var v4857 int32
	_ = v4857
	var v4861 int32
	_ = v4861
	var v4862 int32
	_ = v4862
	var v4870 int32
	_ = v4870
	var v4871 int32
	_ = v4871
	var v4872 int32
	_ = v4872
	var v4873 int32
	_ = v4873
	var v4875 int32
	_ = v4875
	var v4876 int32
	_ = v4876
	var v4880 int32
	_ = v4880
	var v4881 int32
	_ = v4881
	var v4882 int32
	_ = v4882
	var v4883 int32
	_ = v4883
	var v4892 int32
	_ = v4892
	var v4893 int32
	_ = v4893
	var v4896 int64
	_ = v4896
	var v4899 int32
	_ = v4899
	var v4901 int64
	_ = v4901
	var v4904 int32
	_ = v4904
	var v4911 int32
	_ = v4911
	var v4913 int32
	_ = v4913
	var v4917 int32
	_ = v4917
	var v4918 int32
	_ = v4918
	var v4919 int32
	_ = v4919
	var v4920 int32
	_ = v4920
	var v4929 int32
	_ = v4929
	var v4930 int32
	_ = v4930
	var v4931 int32
	_ = v4931
	var v4933 int32
	_ = v4933
	var v4938 int32
	_ = v4938
	var v4939 int32
	_ = v4939
	var v4946 int32
	_ = v4946
	var v4949 int32
	_ = v4949
	var v4951 int32
	_ = v4951
	var v4955 int32
	_ = v4955
	var v4956 int32
	_ = v4956
	var v4957 int32
	_ = v4957
	var v4958 int32
	_ = v4958
	var v4995 int32
	_ = v4995
	var v4997 int32
	_ = v4997
	var v4998 int32
	_ = v4998
	var v5001 int32
	_ = v5001
	var v5005 int32
	_ = v5005
	var v5007 int32
	_ = v5007
	var v5008 int32
	_ = v5008
	var v5009 int32
	_ = v5009
	var v5010 int32
	_ = v5010
	var v5011 int32
	_ = v5011
	var v5014 int32
	_ = v5014
	var v5015 int32
	_ = v5015
	var v5018 int32
	_ = v5018
	var v5021 int32
	_ = v5021
	var v5022 int32
	_ = v5022
	var v5023 int32
	_ = v5023
	var v5025 int32
	_ = v5025
	var v5030 int32
	_ = v5030
	var v5040 int32
	_ = v5040
	var v5102 int32
	_ = v5102
	var v5105 int32
	_ = v5105
	var v5109 int32
	_ = v5109
	var v5118 int32
	_ = v5118
	var v5131 int32
	_ = v5131
	var v5184 int32
	_ = v5184
	var v5186 int32
	_ = v5186
	var v5187 int32
	_ = v5187
	var v5199 int32
	_ = v5199
	var v5262 int32
	_ = v5262
	var v5267 int32
	_ = v5267
	var v5268 int32
	_ = v5268
	var v5269 int32
	_ = v5269
	var v5317 int32
	_ = v5317
	var v5343 int32
	_ = v5343
	var v5345 int32
	_ = v5345
	var v5346 int32
	_ = v5346
	var v5354 int32
	_ = v5354
	var v5355 int32
	_ = v5355
	var v5356 int32
	_ = v5356
	var v5359 int32
	_ = v5359
	var v5362 int32
	_ = v5362
	var v5363 int32
	_ = v5363
	var v5368 int32
	_ = v5368
	var v5378 int32
	_ = v5378
	var v5440 int32
	_ = v5440
	var v5443 int32
	_ = v5443
	var v5447 int32
	_ = v5447
	var v5456 int32
	_ = v5456
	var v5469 int32
	_ = v5469
	var v5522 int32
	_ = v5522
	var v5524 int32
	_ = v5524
	var v5525 int32
	_ = v5525
	var v5537 int32
	_ = v5537
	var v5600 int32
	_ = v5600
	var v5605 int32
	_ = v5605
	var v5606 int32
	_ = v5606
	var v5607 int32
	_ = v5607
	var v5632 int32
	_ = v5632
	var v5681 int32
	_ = v5681
	var v5682 int32
	_ = v5682
	var v5683 int32
	_ = v5683
	var v5684 int32
	_ = v5684
	var v5685 int32
	_ = v5685
	var v5687 int32
	_ = v5687
	var v5692 int32
	_ = v5692
	var v5695 int32
	_ = v5695
	var v5708 int32
	_ = v5708
	var v5709 int32
	_ = v5709
	var v5713 int32
	_ = v5713
	var v5714 int32
	_ = v5714
	var v5722 int32
	_ = v5722
	var v5723 int32
	_ = v5723
	var v5724 int32
	_ = v5724
	var v5725 int32
	_ = v5725
	var v5727 int32
	_ = v5727
	var v5728 int32
	_ = v5728
	var v5732 int32
	_ = v5732
	var v5733 int32
	_ = v5733
	var v5734 int32
	_ = v5734
	var v5735 int32
	_ = v5735
	var v5744 int32
	_ = v5744
	var v5745 int32
	_ = v5745
	var v5748 int64
	_ = v5748
	var v5751 int32
	_ = v5751
	var v5753 int64
	_ = v5753
	var v5756 int32
	_ = v5756
	var v5763 int32
	_ = v5763
	var v5765 int32
	_ = v5765
	var v5769 int32
	_ = v5769
	var v5770 int32
	_ = v5770
	var v5771 int32
	_ = v5771
	var v5772 int32
	_ = v5772
	var v5781 int32
	_ = v5781
	var v5782 int32
	_ = v5782
	var v5783 int32
	_ = v5783
	var v5785 int32
	_ = v5785
	var v5790 int32
	_ = v5790
	var v5791 int32
	_ = v5791
	var v5798 int32
	_ = v5798
	var v5801 int32
	_ = v5801
	var v5803 int32
	_ = v5803
	var v5807 int32
	_ = v5807
	var v5808 int32
	_ = v5808
	var v5809 int32
	_ = v5809
	var v5810 int32
	_ = v5810
	var v5847 int32
	_ = v5847
	var v5848 int32
	_ = v5848
	var v5860 int32
	_ = v5860
	var v5861 int32
	_ = v5861
	var v5865 int32
	_ = v5865
	var v5866 int32
	_ = v5866
	var v5874 int32
	_ = v5874
	var v5875 int32
	_ = v5875
	var v5876 int32
	_ = v5876
	var v5877 int32
	_ = v5877
	var v5879 int32
	_ = v5879
	var v5880 int32
	_ = v5880
	var v5884 int32
	_ = v5884
	var v5885 int32
	_ = v5885
	var v5886 int32
	_ = v5886
	var v5887 int32
	_ = v5887
	var v5896 int32
	_ = v5896
	var v5897 int32
	_ = v5897
	var v5900 int64
	_ = v5900
	var v5903 int32
	_ = v5903
	var v5905 int64
	_ = v5905
	var v5908 int32
	_ = v5908
	var v5915 int32
	_ = v5915
	var v5917 int32
	_ = v5917
	var v5921 int32
	_ = v5921
	var v5922 int32
	_ = v5922
	var v5923 int32
	_ = v5923
	var v5924 int32
	_ = v5924
	var v5933 int32
	_ = v5933
	var v5934 int32
	_ = v5934
	var v5935 int32
	_ = v5935
	var v5937 int32
	_ = v5937
	var v5942 int32
	_ = v5942
	var v5943 int32
	_ = v5943
	var v5950 int32
	_ = v5950
	var v5953 int32
	_ = v5953
	var v5955 int32
	_ = v5955
	var v5959 int32
	_ = v5959
	var v5960 int32
	_ = v5960
	var v5961 int32
	_ = v5961
	var v5962 int32
	_ = v5962
	var v5999 int32
	_ = v5999
	var v6001 int32
	_ = v6001
	var v6002 int32
	_ = v6002
	var v6014 int32
	_ = v6014
	var v6015 int32
	_ = v6015
	var v6019 int32
	_ = v6019
	var v6020 int32
	_ = v6020
	var v6028 int32
	_ = v6028
	var v6029 int32
	_ = v6029
	var v6030 int32
	_ = v6030
	var v6031 int32
	_ = v6031
	var v6033 int32
	_ = v6033
	var v6034 int32
	_ = v6034
	var v6038 int32
	_ = v6038
	var v6039 int32
	_ = v6039
	var v6040 int32
	_ = v6040
	var v6041 int32
	_ = v6041
	var v6050 int32
	_ = v6050
	var v6051 int32
	_ = v6051
	var v6054 int64
	_ = v6054
	var v6057 int32
	_ = v6057
	var v6059 int64
	_ = v6059
	var v6062 int32
	_ = v6062
	var v6069 int32
	_ = v6069
	var v6071 int32
	_ = v6071
	var v6075 int32
	_ = v6075
	var v6076 int32
	_ = v6076
	var v6077 int32
	_ = v6077
	var v6078 int32
	_ = v6078
	var v6087 int32
	_ = v6087
	var v6088 int32
	_ = v6088
	var v6089 int32
	_ = v6089
	var v6091 int32
	_ = v6091
	var v6096 int32
	_ = v6096
	var v6097 int32
	_ = v6097
	var v6104 int32
	_ = v6104
	var v6107 int32
	_ = v6107
	var v6109 int32
	_ = v6109
	var v6113 int32
	_ = v6113
	var v6114 int32
	_ = v6114
	var v6115 int32
	_ = v6115
	var v6116 int32
	_ = v6116
	var v6153 int32
	_ = v6153
	var v6155 int32
	_ = v6155
	var v6156 int32
	_ = v6156
	var v6161 int32
	_ = v6161
	var v6165 int32
	_ = v6165
	var v6216 int32
	_ = v6216
	var v6235 int32
	_ = v6235
	var v6239 int32
	_ = v6239
	var v6240 int32
	_ = v6240
	var v6241 int32
	_ = v6241
	var v6242 int32
	_ = v6242
	var v6243 int32
	_ = v6243
	var v6244 int32
	_ = v6244
	var v6246 int32
	_ = v6246
	var v6249 int32
	_ = v6249
	var v6252 int32
	_ = v6252
	var v6265 int32
	_ = v6265
	var v6266 int32
	_ = v6266
	var v6270 int32
	_ = v6270
	var v6271 int32
	_ = v6271
	var v6279 int32
	_ = v6279
	var v6280 int32
	_ = v6280
	var v6281 int32
	_ = v6281
	var v6282 int32
	_ = v6282
	var v6284 int32
	_ = v6284
	var v6285 int32
	_ = v6285
	var v6289 int32
	_ = v6289
	var v6290 int32
	_ = v6290
	var v6291 int32
	_ = v6291
	var v6292 int32
	_ = v6292
	var v6301 int32
	_ = v6301
	var v6302 int32
	_ = v6302
	var v6305 int64
	_ = v6305
	var v6308 int32
	_ = v6308
	var v6310 int64
	_ = v6310
	var v6313 int32
	_ = v6313
	var v6320 int32
	_ = v6320
	var v6322 int32
	_ = v6322
	var v6326 int32
	_ = v6326
	var v6327 int32
	_ = v6327
	var v6328 int32
	_ = v6328
	var v6329 int32
	_ = v6329
	var v6338 int32
	_ = v6338
	var v6339 int32
	_ = v6339
	var v6340 int32
	_ = v6340
	var v6342 int32
	_ = v6342
	var v6347 int32
	_ = v6347
	var v6348 int32
	_ = v6348
	var v6355 int32
	_ = v6355
	var v6358 int32
	_ = v6358
	var v6360 int32
	_ = v6360
	var v6364 int32
	_ = v6364
	var v6365 int32
	_ = v6365
	var v6366 int32
	_ = v6366
	var v6367 int32
	_ = v6367
	var v6404 int32
	_ = v6404
	var v6405 int32
	_ = v6405
	var v6417 int32
	_ = v6417
	var v6418 int32
	_ = v6418
	var v6422 int32
	_ = v6422
	var v6423 int32
	_ = v6423
	var v6431 int32
	_ = v6431
	var v6432 int32
	_ = v6432
	var v6433 int32
	_ = v6433
	var v6434 int32
	_ = v6434
	var v6436 int32
	_ = v6436
	var v6437 int32
	_ = v6437
	var v6441 int32
	_ = v6441
	var v6442 int32
	_ = v6442
	var v6443 int32
	_ = v6443
	var v6444 int32
	_ = v6444
	var v6453 int32
	_ = v6453
	var v6454 int32
	_ = v6454
	var v6457 int64
	_ = v6457
	var v6460 int32
	_ = v6460
	var v6462 int64
	_ = v6462
	var v6465 int32
	_ = v6465
	var v6472 int32
	_ = v6472
	var v6474 int32
	_ = v6474
	var v6478 int32
	_ = v6478
	var v6479 int32
	_ = v6479
	var v6480 int32
	_ = v6480
	var v6481 int32
	_ = v6481
	var v6490 int32
	_ = v6490
	var v6491 int32
	_ = v6491
	var v6492 int32
	_ = v6492
	var v6494 int32
	_ = v6494
	var v6499 int32
	_ = v6499
	var v6500 int32
	_ = v6500
	var v6507 int32
	_ = v6507
	var v6510 int32
	_ = v6510
	var v6512 int32
	_ = v6512
	var v6516 int32
	_ = v6516
	var v6517 int32
	_ = v6517
	var v6518 int32
	_ = v6518
	var v6519 int32
	_ = v6519
	var v6556 int32
	_ = v6556
	var v6558 int32
	_ = v6558
	var v6559 int32
	_ = v6559
	var v6571 int32
	_ = v6571
	var v6572 int32
	_ = v6572
	var v6576 int32
	_ = v6576
	var v6577 int32
	_ = v6577
	var v6585 int32
	_ = v6585
	var v6586 int32
	_ = v6586
	var v6587 int32
	_ = v6587
	var v6588 int32
	_ = v6588
	var v6590 int32
	_ = v6590
	var v6591 int32
	_ = v6591
	var v6595 int32
	_ = v6595
	var v6596 int32
	_ = v6596
	var v6597 int32
	_ = v6597
	var v6598 int32
	_ = v6598
	var v6607 int32
	_ = v6607
	var v6608 int32
	_ = v6608
	var v6611 int64
	_ = v6611
	var v6614 int32
	_ = v6614
	var v6616 int64
	_ = v6616
	var v6619 int32
	_ = v6619
	var v6626 int32
	_ = v6626
	var v6628 int32
	_ = v6628
	var v6632 int32
	_ = v6632
	var v6633 int32
	_ = v6633
	var v6634 int32
	_ = v6634
	var v6635 int32
	_ = v6635
	var v6644 int32
	_ = v6644
	var v6645 int32
	_ = v6645
	var v6646 int32
	_ = v6646
	var v6648 int32
	_ = v6648
	var v6653 int32
	_ = v6653
	var v6654 int32
	_ = v6654
	var v6661 int32
	_ = v6661
	var v6664 int32
	_ = v6664
	var v6666 int32
	_ = v6666
	var v6670 int32
	_ = v6670
	var v6671 int32
	_ = v6671
	var v6672 int32
	_ = v6672
	var v6673 int32
	_ = v6673
	var v6710 int32
	_ = v6710
	var v6712 int32
	_ = v6712
	var v6713 int32
	_ = v6713
	var v6720 int32
	_ = v6720
	var v6724 int32
	_ = v6724
	var v6733 int32
	_ = v6733
	var v6734 int32
	_ = v6734
	var v6735 int32
	_ = v6735
	var v6747 int32
	_ = v6747
	var v6748 int32
	_ = v6748
	var v6752 int32
	_ = v6752
	var v6753 int32
	_ = v6753
	var v6761 int32
	_ = v6761
	var v6762 int32
	_ = v6762
	var v6763 int32
	_ = v6763
	var v6764 int32
	_ = v6764
	var v6766 int32
	_ = v6766
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
	var v6783 int32
	_ = v6783
	var v6784 int32
	_ = v6784
	var v6787 int64
	_ = v6787
	var v6790 int32
	_ = v6790
	var v6792 int64
	_ = v6792
	var v6795 int32
	_ = v6795
	var v6802 int32
	_ = v6802
	var v6804 int32
	_ = v6804
	var v6808 int32
	_ = v6808
	var v6809 int32
	_ = v6809
	var v6810 int32
	_ = v6810
	var v6811 int32
	_ = v6811
	var v6820 int32
	_ = v6820
	var v6821 int32
	_ = v6821
	var v6822 int32
	_ = v6822
	var v6824 int32
	_ = v6824
	var v6829 int32
	_ = v6829
	var v6830 int32
	_ = v6830
	var v6837 int32
	_ = v6837
	var v6840 int32
	_ = v6840
	var v6842 int32
	_ = v6842
	var v6846 int32
	_ = v6846
	var v6847 int32
	_ = v6847
	var v6848 int32
	_ = v6848
	var v6849 int32
	_ = v6849
	var v6886 int32
	_ = v6886
	var v6887 int32
	_ = v6887
	var v6888 int32
	_ = v6888
	var v6889 int32
	_ = v6889
	var v6890 int64
	_ = v6890
	var v6891 int32
	_ = v6891
	var v6898 int64
	_ = v6898
	var v6899 int32
	_ = v6899
	var v6900 int64
	_ = v6900
	var v6910 int32
	_ = v6910
	var v6912 int32
	_ = v6912
	var v6915 int32
	_ = v6915
	var v6917 int32
	_ = v6917
	var v6922 int32
	_ = v6922
	var v6923 int64
	_ = v6923
	var v6924 int32
	_ = v6924
	var v6931 int64
	_ = v6931
	var v6932 int32
	_ = v6932
	var v6933 int64
	_ = v6933
	var v6943 int32
	_ = v6943
	var v6945 int32
	_ = v6945
	var v6953 int64
	_ = v6953
	var v6961 int64
	_ = v6961
	var v6963 int64
	_ = v6963
	var v6965 int64
	_ = v6965
	var v6968 int32
	_ = v6968
	var v6984 int32
	_ = v6984
	var v6986 int32
	_ = v6986
	var v6990 int32
	_ = v6990
	var v6999 int32
	_ = v6999
	var v7007 int32
	_ = v7007
	var v7011 int32
	_ = v7011
	var v7018 int32
	_ = v7018
	var v7021 int32
	_ = v7021
	var v7023 int32
	_ = v7023
	var v7026 int32
	_ = v7026
	var v7029 int32
	_ = v7029
	var v7030 int32
	_ = v7030
	var v7033 int32
	_ = v7033
	var v7037 int32
	_ = v7037
	var v7047 int32
	_ = v7047
	var v7048 int32
	_ = v7048
	var v7049 int32
	_ = v7049
	var v7051 int32
	_ = v7051
	var v7052 int32
	_ = v7052
	var v7053 int32
	_ = v7053
	var v7055 int64
	_ = v7055
	var v7061 int32
	_ = v7061
	var v7062 int32
	_ = v7062
	var v7063 int32
	_ = v7063
	var v7070 int32
	_ = v7070
	var v7072 int32
	_ = v7072
	var v7076 int32
	_ = v7076
	var v7077 int32
	_ = v7077
	var v7082 int32
	_ = v7082
	var v7083 int32
	_ = v7083
	var v7084 int32
	_ = v7084
	var v7091 int32
	_ = v7091
	var v7092 int32
	_ = v7092
	var v7093 int32
	_ = v7093
	var v7094 int32
	_ = v7094
	var v7095 int64
	_ = v7095
	var v7097 int32
	_ = v7097
	var v7099 int32
	_ = v7099
	var v7107 int32
	_ = v7107
	var v7110 int32
	_ = v7110
	var v7111 int32
	_ = v7111
	var v7116 int32
	_ = v7116
	var v7129 int32
	_ = v7129
	var v7161 int32
	_ = v7161
	var v7162 int32
	_ = v7162
	var v7184 int32
	_ = v7184
	var v7185 int32
	_ = v7185
	var v7193 int32
	_ = v7193
	var v7194 int32
	_ = v7194
	var v7195 int32
	_ = v7195
	var v7198 int32
	_ = v7198
	var v7200 int32
	_ = v7200
	var v7203 int32
	_ = v7203
	var v7225 int32
	_ = v7225
	var v7245 int32
	_ = v7245
	var v7246 int32
	_ = v7246
	var v7268 int32
	_ = v7268
	var v7271 int32
	_ = v7271
	var v7272 int32
	_ = v7272
	var v7275 int32
	_ = v7275
	var v7284 int32
	_ = v7284
	var v7285 int64
	_ = v7285
	var v7287 int32
	_ = v7287
	var v7289 int32
	_ = v7289
	var v7290 int32
	_ = v7290
	var v7291 int32
	_ = v7291
	var v7295 int32
	_ = v7295
	var v7299 int32
	_ = v7299
	var v7300 int32
	_ = v7300
	var v7304 int32
	_ = v7304
	var v7309 int32
	_ = v7309
	var v7312 int32
	_ = v7312
	var v7320 int64
	_ = v7320
	var v7327 int32
	_ = v7327
	var v7331 int32
	_ = v7331
	var v7359 int32
	_ = v7359
	var v7360 int32
	_ = v7360
	var v7362 int32
	_ = v7362
	var v7365 int32
	_ = v7365
	var v7366 int32
	_ = v7366
	var v7368 int32
	_ = v7368
	var v7371 int64
	_ = v7371
	var v7372 int32
	_ = v7372
	var v7379 int64
	_ = v7379
	var v7380 int32
	_ = v7380
	var v7381 int64
	_ = v7381
	var v7390 int32
	_ = v7390
	var v7392 int32
	_ = v7392
	var v7399 int32
	_ = v7399
	var v7402 int32
	_ = v7402
	var v7404 int32
	_ = v7404
	var v7411 int32
	_ = v7411
	var v7424 int32
	_ = v7424
	var v7425 int32
	_ = v7425
	var v7429 int32
	_ = v7429
	var v7430 int32
	_ = v7430
	var v7438 int32
	_ = v7438
	var v7439 int32
	_ = v7439
	var v7440 int32
	_ = v7440
	var v7441 int32
	_ = v7441
	var v7443 int32
	_ = v7443
	var v7444 int32
	_ = v7444
	var v7448 int32
	_ = v7448
	var v7449 int32
	_ = v7449
	var v7450 int32
	_ = v7450
	var v7451 int32
	_ = v7451
	var v7460 int32
	_ = v7460
	var v7461 int32
	_ = v7461
	var v7464 int64
	_ = v7464
	var v7467 int32
	_ = v7467
	var v7469 int64
	_ = v7469
	var v7472 int32
	_ = v7472
	var v7479 int32
	_ = v7479
	var v7481 int32
	_ = v7481
	var v7485 int32
	_ = v7485
	var v7486 int32
	_ = v7486
	var v7487 int32
	_ = v7487
	var v7488 int32
	_ = v7488
	var v7497 int32
	_ = v7497
	var v7498 int32
	_ = v7498
	var v7499 int32
	_ = v7499
	var v7501 int32
	_ = v7501
	var v7506 int32
	_ = v7506
	var v7507 int32
	_ = v7507
	var v7514 int32
	_ = v7514
	var v7517 int32
	_ = v7517
	var v7519 int32
	_ = v7519
	var v7523 int32
	_ = v7523
	var v7524 int32
	_ = v7524
	var v7525 int32
	_ = v7525
	var v7526 int32
	_ = v7526
	var v7563 int32
	_ = v7563
	var v7575 int32
	_ = v7575
	var v7576 int32
	_ = v7576
	var v7580 int32
	_ = v7580
	var v7581 int32
	_ = v7581
	var v7589 int32
	_ = v7589
	var v7590 int32
	_ = v7590
	var v7591 int32
	_ = v7591
	var v7592 int32
	_ = v7592
	var v7594 int32
	_ = v7594
	var v7595 int32
	_ = v7595
	var v7599 int32
	_ = v7599
	var v7600 int32
	_ = v7600
	var v7601 int32
	_ = v7601
	var v7602 int32
	_ = v7602
	var v7611 int32
	_ = v7611
	var v7612 int32
	_ = v7612
	var v7615 int64
	_ = v7615
	var v7618 int32
	_ = v7618
	var v7620 int64
	_ = v7620
	var v7623 int32
	_ = v7623
	var v7630 int32
	_ = v7630
	var v7632 int32
	_ = v7632
	var v7636 int32
	_ = v7636
	var v7637 int32
	_ = v7637
	var v7638 int32
	_ = v7638
	var v7639 int32
	_ = v7639
	var v7648 int32
	_ = v7648
	var v7649 int32
	_ = v7649
	var v7650 int32
	_ = v7650
	var v7652 int32
	_ = v7652
	var v7657 int32
	_ = v7657
	var v7658 int32
	_ = v7658
	var v7665 int32
	_ = v7665
	var v7668 int32
	_ = v7668
	var v7670 int32
	_ = v7670
	var v7674 int32
	_ = v7674
	var v7675 int32
	_ = v7675
	var v7676 int32
	_ = v7676
	var v7677 int32
	_ = v7677
	var v7714 int32
	_ = v7714
	var v7715 int32
	_ = v7715
	var v7727 int32
	_ = v7727
	var v7728 int32
	_ = v7728
	var v7732 int32
	_ = v7732
	var v7733 int32
	_ = v7733
	var v7741 int32
	_ = v7741
	var v7742 int32
	_ = v7742
	var v7743 int32
	_ = v7743
	var v7744 int32
	_ = v7744
	var v7746 int32
	_ = v7746
	var v7747 int32
	_ = v7747
	var v7751 int32
	_ = v7751
	var v7752 int32
	_ = v7752
	var v7753 int32
	_ = v7753
	var v7754 int32
	_ = v7754
	var v7763 int32
	_ = v7763
	var v7764 int32
	_ = v7764
	var v7767 int64
	_ = v7767
	var v7770 int32
	_ = v7770
	var v7772 int64
	_ = v7772
	var v7775 int32
	_ = v7775
	var v7782 int32
	_ = v7782
	var v7784 int32
	_ = v7784
	var v7788 int32
	_ = v7788
	var v7789 int32
	_ = v7789
	var v7790 int32
	_ = v7790
	var v7791 int32
	_ = v7791
	var v7800 int32
	_ = v7800
	var v7801 int32
	_ = v7801
	var v7802 int32
	_ = v7802
	var v7804 int32
	_ = v7804
	var v7809 int32
	_ = v7809
	var v7810 int32
	_ = v7810
	var v7817 int32
	_ = v7817
	var v7820 int32
	_ = v7820
	var v7822 int32
	_ = v7822
	var v7826 int32
	_ = v7826
	var v7827 int32
	_ = v7827
	var v7828 int32
	_ = v7828
	var v7829 int32
	_ = v7829
	var v7866 int32
	_ = v7866
	var v7867 int32
	_ = v7867
	var v7869 int32
	_ = v7869
	var v7875 int32
	_ = v7875
	var v7876 int32
	_ = v7876
	var v7911 int32
	_ = v7911
	var v7937 int32
	_ = v7937
	var v7938 int32
	_ = v7938
	var v7942 int32
	_ = v7942
	var v7943 int32
	_ = v7943
	var v7950 int32
	_ = v7950
	var v7954 int32
	_ = v7954
	var v7958 int32
	_ = v7958
	var v7965 int32
	_ = v7965
	var v8000 int32
	_ = v8000
	var v8030 int32
	_ = v8030
	var v8031 int32
	_ = v8031
	var v8038 int32
	_ = v8038
	var v8040 int32
	_ = v8040
	var v8041 int32
	_ = v8041
	var v8042 int32
	_ = v8042
	var v8055 int32
	_ = v8055
	var v8056 int32
	_ = v8056
	var v8060 int32
	_ = v8060
	var v8061 int32
	_ = v8061
	var v8069 int32
	_ = v8069
	var v8070 int32
	_ = v8070
	var v8071 int32
	_ = v8071
	var v8072 int32
	_ = v8072
	var v8074 int32
	_ = v8074
	var v8075 int32
	_ = v8075
	var v8079 int32
	_ = v8079
	var v8080 int32
	_ = v8080
	var v8081 int32
	_ = v8081
	var v8082 int32
	_ = v8082
	var v8091 int32
	_ = v8091
	var v8092 int32
	_ = v8092
	var v8095 int64
	_ = v8095
	var v8098 int32
	_ = v8098
	var v8100 int64
	_ = v8100
	var v8103 int32
	_ = v8103
	var v8110 int32
	_ = v8110
	var v8112 int32
	_ = v8112
	var v8116 int32
	_ = v8116
	var v8117 int32
	_ = v8117
	var v8118 int32
	_ = v8118
	var v8119 int32
	_ = v8119
	var v8128 int32
	_ = v8128
	var v8129 int32
	_ = v8129
	var v8130 int32
	_ = v8130
	var v8132 int32
	_ = v8132
	var v8137 int32
	_ = v8137
	var v8138 int32
	_ = v8138
	var v8145 int32
	_ = v8145
	var v8148 int32
	_ = v8148
	var v8150 int32
	_ = v8150
	var v8154 int32
	_ = v8154
	var v8155 int32
	_ = v8155
	var v8156 int32
	_ = v8156
	var v8157 int32
	_ = v8157
	var v8200 int32
	_ = v8200
	var v8202 int32
	_ = v8202
	var v8203 int32
	_ = v8203
	var v8215 int32
	_ = v8215
	var v8216 int32
	_ = v8216
	var v8220 int32
	_ = v8220
	var v8221 int32
	_ = v8221
	var v8229 int32
	_ = v8229
	var v8230 int32
	_ = v8230
	var v8231 int32
	_ = v8231
	var v8232 int32
	_ = v8232
	var v8234 int32
	_ = v8234
	var v8235 int32
	_ = v8235
	var v8239 int32
	_ = v8239
	var v8240 int32
	_ = v8240
	var v8241 int32
	_ = v8241
	var v8242 int32
	_ = v8242
	var v8251 int32
	_ = v8251
	var v8252 int32
	_ = v8252
	var v8255 int64
	_ = v8255
	var v8258 int32
	_ = v8258
	var v8260 int64
	_ = v8260
	var v8263 int32
	_ = v8263
	var v8270 int32
	_ = v8270
	var v8272 int32
	_ = v8272
	var v8276 int32
	_ = v8276
	var v8277 int32
	_ = v8277
	var v8278 int32
	_ = v8278
	var v8279 int32
	_ = v8279
	var v8288 int32
	_ = v8288
	var v8289 int32
	_ = v8289
	var v8290 int32
	_ = v8290
	var v8292 int32
	_ = v8292
	var v8297 int32
	_ = v8297
	var v8298 int32
	_ = v8298
	var v8305 int32
	_ = v8305
	var v8308 int32
	_ = v8308
	var v8310 int32
	_ = v8310
	var v8314 int32
	_ = v8314
	var v8315 int32
	_ = v8315
	var v8316 int32
	_ = v8316
	var v8317 int32
	_ = v8317
	var v8357 int32
	_ = v8357
	var v8358 int32
	_ = v8358
	var v8364 int32
	_ = v8364
	var v8367 int32
	_ = v8367
	var v8374 int32
	_ = v8374
	var v8376 int32
	_ = v8376
	var v8377 int32
	_ = v8377
	var v8380 int32
	_ = v8380
	var v8382 int32
	_ = v8382
	var v8383 int32
	_ = v8383
	var v8389 int32
	_ = v8389
	var v8395 int32
	_ = v8395
	var v8396 int32
	_ = v8396
	var v8408 int32
	_ = v8408
	var v8457 int32
	_ = v8457
	var v8459 int32
	_ = v8459
	var v8462 int32
	_ = v8462
	var v8464 int32
	_ = v8464
	var v8467 int32
	_ = v8467
	var v8469 int32
	_ = v8469
	var v8472 int32
	_ = v8472
	var v8474 int32
	_ = v8474
	var v8478 int32
	_ = v8478
	var v8484 int32
	_ = v8484
	var v8497 int32
	_ = v8497
	var v8555 int32
	_ = v8555
	var v8556 int32
	_ = v8556
	var v8567 int32
	_ = v8567
	var v8617 int32
	_ = v8617
	var v8619 int32
	_ = v8619
	var v8623 int32
	_ = v8623
	var v8628 int32
	_ = v8628
	var v8690 int64
	_ = v8690
	var v8691 int32
	_ = v8691
	var v8698 int64
	_ = v8698
	var v8699 int32
	_ = v8699
	var v8700 int64
	_ = v8700
	var v8710 int32
	_ = v8710
	var v8712 int32
	_ = v8712
	var v8756 int32
	_ = v8756
	var v8785 int32
	_ = v8785
	var v8787 int32
	_ = v8787
	var v8790 int32
	_ = v8790
	var v8797 int32
	_ = v8797
	var v8808 int32
	_ = v8808
	var v8809 int32
	_ = v8809
	var v8858 int32
	_ = v8858
	var v8864 int32
	_ = v8864
	var v8869 int32
	_ = v8869
	var v8870 int32
	_ = v8870
	var v8877 int32
	_ = v8877
	var v8888 int32
	_ = v8888
	var v8938 int32
	_ = v8938
	var v8940 int32
	_ = v8940
	var v8942 int32
	_ = v8942
	var v8949 int32
	_ = v8949
	var v8950 int32
	_ = v8950
	var v9019 int32
	_ = v9019
	var v9041 int32
	_ = v9041
	var v9112 int32
	_ = v9112
	var v9114 int32
	_ = v9114
	var v9126 int32
	_ = v9126
	var v9127 int32
	_ = v9127
	var v9131 int32
	_ = v9131
	var v9132 int32
	_ = v9132
	var v9140 int32
	_ = v9140
	var v9141 int32
	_ = v9141
	var v9142 int32
	_ = v9142
	var v9143 int32
	_ = v9143
	var v9145 int32
	_ = v9145
	var v9146 int32
	_ = v9146
	var v9150 int32
	_ = v9150
	var v9151 int32
	_ = v9151
	var v9152 int32
	_ = v9152
	var v9153 int32
	_ = v9153
	var v9162 int32
	_ = v9162
	var v9163 int32
	_ = v9163
	var v9166 int64
	_ = v9166
	var v9169 int32
	_ = v9169
	var v9171 int64
	_ = v9171
	var v9174 int32
	_ = v9174
	var v9181 int32
	_ = v9181
	var v9183 int32
	_ = v9183
	var v9187 int32
	_ = v9187
	var v9188 int32
	_ = v9188
	var v9189 int32
	_ = v9189
	var v9190 int32
	_ = v9190
	var v9199 int32
	_ = v9199
	var v9200 int32
	_ = v9200
	var v9201 int32
	_ = v9201
	var v9203 int32
	_ = v9203
	var v9208 int32
	_ = v9208
	var v9209 int32
	_ = v9209
	var v9216 int32
	_ = v9216
	var v9219 int32
	_ = v9219
	var v9221 int32
	_ = v9221
	var v9225 int32
	_ = v9225
	var v9226 int32
	_ = v9226
	var v9227 int32
	_ = v9227
	var v9228 int32
	_ = v9228
	var v9285 int32
	_ = v9285
	var v9294 int32
	_ = v9294
	var v9357 int32
	_ = v9357
	var v9406 int32
	_ = v9406
	var v9422 int32
	_ = v9422
	var v9426 int32
	_ = v9426
	var v9427 int32
	_ = v9427
	var v9428 int32
	_ = v9428
	var v9429 int32
	_ = v9429
	var v9430 int32
	_ = v9430
	var v9432 int32
	_ = v9432
	var v9437 int32
	_ = v9437
	var v9438 int32
	_ = v9438
	var v9439 int32
	_ = v9439
	var v9440 int32
	_ = v9440
	var v9447 int32
	_ = v9447
	var v9450 int32
	_ = v9450
	var v9452 int32
	_ = v9452
	var v9461 int64
	_ = v9461
	var v9462 int64
	_ = v9462
	var v9465 int32
	_ = v9465
	var v9466 int64
	_ = v9466
	var v9468 int32
	_ = v9468
	var v9469 int64
	_ = v9469
	var v9472 int32
	_ = v9472
	var v9473 int64
	_ = v9473
	var v9475 int32
	_ = v9475
	var v9476 int64
	_ = v9476
	var v9489 int32
	_ = v9489
	var v9490 int32
	_ = v9490
	var v9491 int32
	_ = v9491
	var v9493 int32
	_ = v9493
	var v9496 int32
	_ = v9496
	var v9497 int32
	_ = v9497
	var v9538 int32
	_ = v9538
	var v9539 int32
	_ = v9539
	var v9561 int32
	_ = v9561
	var v9569 int32
	_ = v9569
	var v9570 int32
	_ = v9570
	var v9608 int32
	_ = v9608
	var v9609 int32
	_ = v9609
	var v9639 int32
	_ = v9639
	var v9648 int64
	_ = v9648
	var v9649 int64
	_ = v9649
	var v9652 int32
	_ = v9652
	var v9653 int64
	_ = v9653
	var v9655 int32
	_ = v9655
	var v9656 int64
	_ = v9656
	var v9659 int32
	_ = v9659
	var v9660 int64
	_ = v9660
	var v9662 int32
	_ = v9662
	var v9663 int64
	_ = v9663
	var v9671 int32
	_ = v9671
	var v9673 int32
	_ = v9673
	var v9712 int32
	_ = v9712
	var v9713 int32
	_ = v9713
	var v9779 int32
	_ = v9779
	var v9780 int32
	_ = v9780
	var v9809 int32
	_ = v9809
	var v9811 int32
	_ = v9811
	var v9813 int32
	_ = v9813
	var v9814 int32
	_ = v9814
	var v9817 int32
	_ = v9817
	var v9818 int32
	_ = v9818
	var v9819 int32
	_ = v9819
	var v9820 int32
	_ = v9820
	var v9821 int32
	_ = v9821
	var v9822 int32
	_ = v9822
	var v9823 int32
	_ = v9823
	var v9824 int32
	_ = v9824
	var v9826 int32
	_ = v9826
	var v9829 int32
	_ = v9829
	var v9830 int32
	_ = v9830
	var v9831 int32
	_ = v9831
	var v9832 int32
	_ = v9832
	var v9833 int32
	_ = v9833
	var v9834 int32
	_ = v9834
	var v9835 int32
	_ = v9835
	var v9836 int32
	_ = v9836
	var v9837 int32
	_ = v9837
	var v9838 int32
	_ = v9838
	var v9839 int32
	_ = v9839
	var v9840 int32
	_ = v9840
	var v9841 int32
	_ = v9841
	var v9842 int32
	_ = v9842
	var v9843 int32
	_ = v9843
	var v9844 int32
	_ = v9844
	var v9845 int32
	_ = v9845
	var v9846 int32
	_ = v9846
	var v9847 int32
	_ = v9847
	var v9848 int32
	_ = v9848
	var v9850 int32
	_ = v9850
	var v9852 int32
	_ = v9852
	var v9854 int32
	_ = v9854
	var v9855 int32
	_ = v9855
	var v9878 int32
	_ = v9878
	var v9880 int32
	_ = v9880
	var v9881 int32
	_ = v9881
	var v9883 int32
	_ = v9883
	var v9884 int32
	_ = v9884
	var v9885 int32
	_ = v9885
	var v9886 int32
	_ = v9886
	var v9887 int32
	_ = v9887
	var v9888 int32
	_ = v9888
	var v9889 int32
	_ = v9889
	var v9890 int32
	_ = v9890
	var v9891 int32
	_ = v9891
	var v9893 int32
	_ = v9893
	var v9896 int32
	_ = v9896
	var v9897 int32
	_ = v9897
	var v9898 int32
	_ = v9898
	var v9899 int32
	_ = v9899
	var v9900 int32
	_ = v9900
	var v9901 int32
	_ = v9901
	var v9902 int32
	_ = v9902
	var v9903 int32
	_ = v9903
	var v9904 int32
	_ = v9904
	var v9905 int32
	_ = v9905
	var v9906 int32
	_ = v9906
	var v9907 int32
	_ = v9907
	var v9908 int32
	_ = v9908
	var v9909 int32
	_ = v9909
	var v9910 int32
	_ = v9910
	var v9911 int32
	_ = v9911
	var v9912 int32
	_ = v9912
	var v9913 int32
	_ = v9913
	var v9914 int32
	_ = v9914
	var v9915 int32
	_ = v9915
	var v9917 int32
	_ = v9917
	var v9919 int32
	_ = v9919
	var v9921 int32
	_ = v9921
	var v9922 int32
	_ = v9922
	var v9945 int32
	_ = v9945
	var v9947 int32
	_ = v9947
	var v9948 int32
	_ = v9948
	var v9949 int32
	_ = v9949
	var v9950 int32
	_ = v9950
	var v9951 int32
	_ = v9951
	var v9952 int32
	_ = v9952
	var v9953 int32
	_ = v9953
	var v9954 int32
	_ = v9954
	var v9955 int32
	_ = v9955
	var v9956 int32
	_ = v9956
	var v9957 int32
	_ = v9957
	var v9958 int32
	_ = v9958
	var v9960 int32
	_ = v9960
	var v9963 int32
	_ = v9963
	var v9964 int32
	_ = v9964
	var v9965 int32
	_ = v9965
	var v9966 int32
	_ = v9966
	var v9967 int32
	_ = v9967
	var v9968 int32
	_ = v9968
	var v9969 int32
	_ = v9969
	var v9970 int32
	_ = v9970
	var v9971 int32
	_ = v9971
	var v9972 int32
	_ = v9972
	var v9973 int32
	_ = v9973
	var v9974 int32
	_ = v9974
	var v9975 int32
	_ = v9975
	var v9976 int32
	_ = v9976
	var v9977 int32
	_ = v9977
	var v9978 int32
	_ = v9978
	var v9979 int32
	_ = v9979
	var v9980 int32
	_ = v9980
	var v9981 int32
	_ = v9981
	var v9982 int32
	_ = v9982
	var v9984 int32
	_ = v9984
	var v9986 int32
	_ = v9986
	var v9988 int32
	_ = v9988
	var v9989 int32
	_ = v9989
	var v10016 int32
	_ = v10016
	var v10017 int32
	_ = v10017
	var v10023 int32
	_ = v10023
	var v10028 int32
	_ = v10028
	var v10033 int32
	_ = v10033
	var v10034 int32
	_ = v10034
	var v10038 int64
	_ = v10038
	var v10045 int32
	_ = v10045
	var v10046 int32
	_ = v10046
	var v10051 int32
	_ = v10051
	var v10052 int32
	_ = v10052
	var v10054 int32
	_ = v10054
	var v10056 int32
	_ = v10056
	var v10057 int32
	_ = v10057
	var v10066 int64
	_ = v10066
	var v10067 int64
	_ = v10067
	var v10069 int32
	_ = v10069
	var v10070 int32
	_ = v10070
	var v10071 int64
	_ = v10071
	var v10073 int32
	_ = v10073
	var v10074 int64
	_ = v10074
	var v10077 int32
	_ = v10077
	var v10078 int64
	_ = v10078
	var v10080 int32
	_ = v10080
	var v10081 int64
	_ = v10081
	var v10088 int32
	_ = v10088
	var v10089 int32
	_ = v10089
	var v10090 int32
	_ = v10090
	var v10092 int32
	_ = v10092
	var v10096 int32
	_ = v10096
	var v10098 int32
	_ = v10098
	var v10099 int32
	_ = v10099
	var v10103 int32
	_ = v10103
	var v10105 int32
	_ = v10105
	var v10106 int32
	_ = v10106
	var v10112 int32
	_ = v10112
	var v10114 int32
	_ = v10114
	var v10118 int32
	_ = v10118
	var v10120 int32
	_ = v10120
	var v10124 int32
	_ = v10124
	var v10128 int32
	_ = v10128
	var v10132 int32
	_ = v10132
	var v10133 int64
	_ = v10133
	var v10135 int32
	_ = v10135
	var v10137 int32
	_ = v10137
	var v10138 int32
	_ = v10138
	var v10139 int32
	_ = v10139
	var v10144 int32
	_ = v10144
	var v10146 int32
	_ = v10146
	var v10148 int32
	_ = v10148
	var v10149 int32
	_ = v10149
	var v10212 int32
	_ = v10212
	var v10213 int32
	_ = v10213
	var v10222 int64
	_ = v10222
	var v10223 int64
	_ = v10223
	var v10225 int32
	_ = v10225
	var v10226 int32
	_ = v10226
	var v10227 int64
	_ = v10227
	var v10229 int32
	_ = v10229
	var v10230 int64
	_ = v10230
	var v10233 int32
	_ = v10233
	var v10234 int64
	_ = v10234
	var v10236 int32
	_ = v10236
	var v10237 int64
	_ = v10237
	var v10242 int32
	_ = v10242
	var v10244 int32
	_ = v10244
	var v10309 int32
	_ = v10309
	var v10314 int32
	_ = v10314
	var v10315 int32
	_ = v10315
	var v10319 int64
	_ = v10319
	var v10326 int32
	_ = v10326
	var v10327 int32
	_ = v10327
	v67 = m.G0
	v69 = v67 - int32(_a_F_EncodeStreamHook_0)
	m.G0 = v69
	v71 = *(*int32)(unsafe.Add(mBase, uint32(l0)+8))
	v73 = v71 + int32(8)
	v74 = *(*int32)(unsafe.Add(mBase, uint32(v73)))
	v75 = *(*int32)(unsafe.Add(mBase, uint32(l0)+416))
	v76 = *(*int32)(unsafe.Add(mBase, uint32(l0)+412))
	v77 = *(*int32)(unsafe.Add(mBase, uint32(l0)+408))
	v78 = *(*int32)(unsafe.Add(mBase, uint32(l0)+12))
	v79 = *(*int32)(unsafe.Add(mBase, uint32(l0)+4))
	v80 = *(*int32)(unsafe.Add(mBase, uint32(v79)+12))
	v81 = *(*int32)(unsafe.Add(mBase, uint32(v79)+8))
	v82 = *(*int32)(unsafe.Add(mBase, uint32(l0)))
	v83 = *(*float32)(unsafe.Add(mBase, uint32(v82)+4))
	v84 = *(*int32)(unsafe.Add(mBase, uint32(v82)+8))
	v85 = *(*int32)(unsafe.Add(mBase, uint32(v71)+4))
	v86 = *(*int32)(unsafe.Add(mBase, uint32(v71)+12))
	*(*int32)(unsafe.Add(mBase, uint32(v69)+60)) = int32(2)
	v95 = *(*int64)(unsafe.Add(mBase, uint32(v71+int32(16))))
	*(*int64)(unsafe.Add(mBase, uint32(v69+int32(48)))) = v95
	v101 = *(*int64)(unsafe.Add(mBase, uint32(v73)))
	*(*int64)(unsafe.Add(mBase, uint32(v69+int32(40)))) = v101
	v103 = *(*int64)(unsafe.Add(mBase, uint32(v71)))
	*(*int64)(unsafe.Add(mBase, uint32(v69)+32)) = v103
	if base.F32_lt(base.F32_abs(v83), float32(2.1474836e+09)) == int32(0) {
		goto L2
	} else {
		goto L3
	}
L1:
	;
	v114 = v69 + int32(8)
	v117 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v114))) = v117
	*(*int64)(unsafe.Add(mBase, uint32(v69+int32(24)))) = v117
	*(*int64)(unsafe.Add(mBase, uint32(v69+int32(16)))) = v117
	v130 = int32(1024)
	v132 = F_WebPSafeMalloc(m, int64(1), v130)
	mBase = m.M
	if v132 != 0 {
		goto L10
	} else {
		goto L11
	}
L2:
	;
	v112 = int32(-2147483648)
	goto L1
L3:
	;
	v110 = base.I32_trunc_f32_s(v83)
	v112 = v110
	goto L1
L4:
	;
	v10309 = v10244 + int32(8)
	if v10309 == int32(0) {
		goto L1165
	} else {
		goto L1166
	}
L5:
	;
	v10212 = int32(8)
	v10213 = v10148 + v10212
	v10222 = *(*int64)(unsafe.Add(mBase, uint32(v10213)))
	v10223 = *(*int64)(unsafe.Add(mBase, uint32(v10149)))
	*(*int64)(unsafe.Add(mBase, uint32(v10213))) = v10223
	v10225 = int32(16)
	v10226 = v10148 + int32(24)
	v10227 = *(*int64)(unsafe.Add(mBase, uint32(v10226)))
	v10229 = v10149 + v10225
	v10230 = *(*int64)(unsafe.Add(mBase, uint32(v10229)))
	*(*int64)(unsafe.Add(mBase, uint32(v10226))) = v10230
	v10233 = v10148 + v10225
	v10234 = *(*int64)(unsafe.Add(mBase, uint32(v10233)))
	v10236 = v10149 + v10212
	v10237 = *(*int64)(unsafe.Add(mBase, uint32(v10236)))
	*(*int64)(unsafe.Add(mBase, uint32(v10233))) = v10237
	*(*int64)(unsafe.Add(mBase, uint32(v10149))) = v10222
	*(*int64)(unsafe.Add(mBase, uint32(v10229))) = v10227
	*(*int64)(unsafe.Add(mBase, uint32(v10236))) = v10234
	goto L1163
L6:
	;
	v234 = int32(0)
	v237 = base.I32_div_u_s(int32(97), v77)
	v239 = int32(base.Ui32(v237) >> (uint(int32(2)) % 32))
	v245 = int32(-1)
	v265 = l0
	v267 = v69
	v268 = v71
	v271 = v75
	v272 = v76
	v273 = v77
	v274 = v78
	v275 = v79
	v276 = v80
	v277 = v81
	v278 = v82
	v280 = v84
	v283 = v112
	v284 = (v85+int32(7))>>(uint(int32(3))%32) + (v86 - v74)
	v285 = l0 + int32(16)
	v286 = base.B2i32(v84 == v234)
	v287 = v237
	v288 = v239
	v289 = v237 - v239
	v290 = v78 + int32(64)
	v291 = v78 + int32(2168)
	v292 = v80 + v245
	v293 = v78 + int32(2216)
	v294 = v78 + int32(72)
	v295 = v78 + int32(1096)
	v296 = v78 + int32(2144)
	v297 = v78 + int32(2120)
	v298 = v69 + int32(2108)
	v299 = v245
	v300 = v234
	v301 = v234
	v302 = v234
	goto L33
L7:
	;
	if v77 != int32(1) {
		v10146 = l0
		v10148 = v69
		v10149 = v71
		goto L5
	} else {
		goto L32
	}
L8:
	;
	v221 = *(*int32)(unsafe.Add(mBase, uint32(v79)+92))
	if v221 != 0 {
		goto L30
	} else {
		goto L31
	}
L9:
	;
	if v143 == int32(0) {
		goto L8
	} else {
		goto L12
	}
L10:
	;
	v136 = *(*int32)(unsafe.Add(mBase, uint32(v114)+8))
	F_WebPSafeFree(m, v136)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v114)+16)) = v132 + v130
	*(*int32)(unsafe.Add(mBase, uint32(v114)+12)) = v132
	*(*int32)(unsafe.Add(mBase, uint32(v114)+8)) = v132
	v143 = int32(1)
	goto L9
L11:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v114)+20)) = int32(1)
	v143 = int32(0)
	goto L9
L12:
	;
	if v77 < int32(2) {
		goto L7
	} else {
		goto L13
	}
L13:
	;
	v149 = v69 + int32(8)
	v159 = *(*int32)(unsafe.Add(mBase, uint32(v149)+12))
	v160 = *(*int32)(unsafe.Add(mBase, uint32(v149)+8))
	v161 = v159 - v160
	v163 = *(*int32)(unsafe.Add(mBase, uint32(v71)+12))
	v164 = *(*int32)(unsafe.Add(mBase, uint32(v71)+8))
	v165 = v163 - v164
	v167 = base.I64_extend_i32_u(v161) + base.I64_extend_i32_u(v165)
	if base.Ui64(v167) < base.Ui64(int64(4294967296)) {
		goto L15
	} else {
		goto L16
	}
L14:
	;
	if v219 != 0 {
		goto L6
	} else {
		goto L28
	}
L15:
	;
	v173 = base.I32_wrap_i64(v167)
	v174 = *(*int32)(unsafe.Add(mBase, uint32(v149)+16))
	v175 = v174 - v160
	if v174 == v160 {
		goto L18
	} else {
		goto L19
	}
L16:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v149)+20)) = int32(1)
	v219 = int32(0)
	goto L14
L17:
	;
	v206 = F_memcpy(m, v205, v204, v165)
	mBase = m.M
	v207 = *(*int64)(unsafe.Add(mBase, uint32(v71)))
	*(*int64)(unsafe.Add(mBase, uint32(v149))) = v207
	v209 = *(*int32)(unsafe.Add(mBase, uint32(v71)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v149)+20)) = v209
	v211 = *(*int32)(unsafe.Add(mBase, uint32(v149)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v149)+12)) = v211 + v165
	v219 = int32(1)
	goto L14
L18:
	;
	v182 = int32(base.Ui32(v175*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v173) < base.Ui32(v182) {
		goto L22
	} else {
		goto L23
	}
L19:
	;
	if base.Ui32(v175) < base.Ui32(v173) {
		goto L18
	} else {
		goto L20
	}
L20:
	;
	v204 = v164
	v205 = v160
	goto L17
L21:
	;
	if v159 == v160 {
		goto L26
	} else {
		goto L27
	}
L22:
	;
	v184 = v182
	goto L24
L23:
	;
	v184 = v173
	goto L24
L24:
	;
	v188 = v184&int32(-1024) + int32(1024)
	v189 = F_WebPSafeMalloc(m, int64(1), v188)
	mBase = m.M
	if v189 != 0 {
		goto L21
	} else {
		goto L25
	}
L25:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v149)+20)) = int32(1)
	v219 = int32(0)
	goto L14
L26:
	;
	v196 = *(*int32)(unsafe.Add(mBase, uint32(v149)+8))
	F_WebPSafeFree(m, v196)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v149)+8)) = v189
	*(*int32)(unsafe.Add(mBase, uint32(v149)+16)) = v189 + v188
	*(*int32)(unsafe.Add(mBase, uint32(v149)+12)) = v189 + v161
	v203 = *(*int32)(unsafe.Add(mBase, uint32(v71)+8))
	v204 = v203
	v205 = v189
	goto L17
L27:
	;
	v194 = *(*int32)(unsafe.Add(mBase, uint32(v149)+8))
	v195 = F_memcpy(m, v189, v194, v161)
	mBase = m.M
	goto L26
L28:
	;
	goto L8
L29:
	;
	v10242 = l0
	v10244 = v69
	goto L4
L30:
	;
	goto L29
L31:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v79)+92)) = int32(1)
	goto L30
L32:
	;
	goto L6
L33:
	;
	v333 = v285 + v300*int32(28)
	v334 = *(*int32)(unsafe.Add(mBase, uint32(v333)))
	v335 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v267)+4)) = v335
	v338 = v334 & int32(-2)
	v340 = base.B2i32(v338 == int32(4))
	*(*int32)(unsafe.Add(mBase, uint32(v274)+64)) = v340
	v348 = base.B2i32(v334 == int32(5)) | base.B2i32(v334&int32(-3) == int32(1))
	*(*int32)(unsafe.Add(mBase, uint32(v274)+60)) = v348
	*(*int32)(unsafe.Add(mBase, uint32(v274)+56)) = base.B2i32(v338 == int32(2))
	*(*int32)(unsafe.Add(mBase, uint32(v274)+48)) = v335
	if v272 != 0 {
		goto L35
	} else {
		goto L36
	}
L34:
	;
	v10146 = v9945
	v10148 = v9947
	v10149 = v9948
	goto L5
L35:
	;
	v357 = v335
	goto L37
L36:
	;
	v357 = v348
	goto L37
L37:
	;
	if v340&int32(1) != 0 {
		goto L38
	} else {
		goto L39
	}
L38:
	;
	v360 = v335
	goto L40
L39:
	;
	v360 = v357
	goto L40
L40:
	;
	if v280 != 0 {
		goto L41
	} else {
		goto L42
	}
L41:
	;
	v362 = v360
	goto L43
L42:
	;
	v362 = int32(0)
	goto L43
L43:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+52)) = v362
	v364 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v267))) = v364
	v368 = *(*int32)(unsafe.Add(mBase, uint32(v297)+12))
	if v368 == v364 {
		goto L45
	} else {
		goto L46
	}
L44:
	;
	v395 = *(*int32)(unsafe.Add(mBase, uint32(v296)+12))
	if v395 == int32(0) {
		goto L53
	} else {
		goto L54
	}
L45:
	;
	v373 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v297)+20)) = v373
	v375 = *(*int32)(unsafe.Add(mBase, uint32(v297)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v297)+16)) = v375
	*(*int32)(unsafe.Add(mBase, uint32(v297)+12)) = v297 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v297)+8)) = v373
	if v375 == v373 {
		goto L47
	} else {
		goto L48
	}
L46:
	;
	v371 = *(*int32)(unsafe.Add(mBase, uint32(v297)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v368))) = v371
	goto L45
L47:
	;
	goto L44
L48:
	;
	v386 = v375
	goto L49
L49:
	;
	v387 = *(*int32)(unsafe.Add(mBase, uint32(v386)))
	F_WebPSafeFree(m, v386)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v297)+16)) = v387
	if v387 != 0 {
		v386 = v387
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
	v421 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	if int32(99) < v421 {
		v670 = int32(0)
		goto L66
	} else {
		goto L67
	}
L53:
	;
	v400 = int32(0)
	*(*int32)(unsafe.Add(mBase, uint32(v296)+20)) = v400
	v402 = *(*int32)(unsafe.Add(mBase, uint32(v296)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v296)+16)) = v402
	*(*int32)(unsafe.Add(mBase, uint32(v296)+12)) = v296 + int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v296)+8)) = v400
	if v402 == v400 {
		goto L55
	} else {
		goto L56
	}
L54:
	;
	v398 = *(*int32)(unsafe.Add(mBase, uint32(v296)+16))
	*(*int32)(unsafe.Add(mBase, uint32(v395))) = v398
	goto L53
L55:
	;
	goto L52
L56:
	;
	v413 = v402
	goto L57
L57:
	;
	v414 = *(*int32)(unsafe.Add(mBase, uint32(v413)))
	F_WebPSafeFree(m, v413)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v296)+16)) = v414
	if v414 != 0 {
		v413 = v414
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
	v4687 = *(*int32)(unsafe.Add(mBase, uint32(v274)+56))
	if v4687 == int32(0) {
		goto L467
	} else {
		goto L468
	}
L61:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(3)
	v4612 = *(*int32)(unsafe.Add(mBase, uint32(v274)+68))
	if int32(1023) < v4612 {
		v4668 = v289
		goto L60
	} else {
		goto L466
	}
L62:
	;
	v4539 = *(*int32)(unsafe.Add(mBase, uint32(v1333)+92))
	if v4539 != 0 {
		goto L463
	} else {
		goto L464
	}
L63:
	;
	v4212 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v4213 = *(*int32)(unsafe.Add(mBase, uint32(v4212)+8))
	v4214 = *(*int32)(unsafe.Add(mBase, uint32(v4212)+12))
	v4221 = base.I64_extend_i32_s(v4214) * base.I64_extend_i32_s(v4213)
	v4222 = *(*int32)(unsafe.Add(mBase, uint32(v274)+60))
	if v4222 == int32(0) {
		goto L443
	} else {
		goto L444
	}
L64:
	;
	if v421 < int32(100) {
		v4668 = v287
		goto L60
	} else {
		goto L439
	}
L65:
	;
	v676 = *(*int32)(unsafe.Add(mBase, uint32(v333)+4))
	v677 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v678 = *(*int32)(unsafe.Add(mBase, uint32(v274)+68))
	v679 = F_PaletteSort(m, v676, v677, v295, v678, v294)
	mBase = m.M
	if v679 != 0 {
		goto L110
	} else {
		goto L111
	}
L66:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = v670
	v672 = *(*int32)(unsafe.Add(mBase, uint32(v274)+64))
	if v672 == int32(0) {
		goto L64
	} else {
		goto L109
	}
L67:
	;
	v424 = *(*int32)(unsafe.Add(mBase, uint32(v274)+64))
	if v424 == int32(0) {
		goto L68
	} else {
		goto L69
	}
L68:
	;
	v429 = *(*int32)(unsafe.Add(mBase, uint32(v274)+60))
	if v429 == int32(0) {
		goto L70
	} else {
		goto L71
	}
L69:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(0)
	goto L65
L70:
	;
	v440 = base.I64_extend_i32_s(v276) * base.I64_extend_i32_s(v277)
	v441 = *(*int32)(unsafe.Add(mBase, uint32(v274)+60))
	if v441 == int32(0) {
		goto L75
	} else {
		goto L76
	}
L71:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(0)
	goto L63
L72:
	;
	if v520 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L84
	}
L73:
	;
	v480 = v440 + v475 + v476 + int64(16)
	v481 = *(*int32)(unsafe.Add(mBase, uint32(v274)+24))
	if v481 == int32(0) {
		goto L79
	} else {
		goto L80
	}
L74:
	;
	v463 = int32(3)
	v465 = int32(2)
	v474 = v460
	v475 = v461
	v476 = base.I64_extend_i32_u(int32(base.Ui32(v276+v463)>>(uint(v465)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v277+v463)>>(uint(v465)%32)))
	goto L73
L75:
	;
	v454 = int32(0)
	v455 = int64(0)
	v457 = *(*int32)(unsafe.Add(mBase, uint32(v274)+52))
	if v457 == v454 {
		v474 = v454
		v475 = v455
		v476 = v455
		goto L73
	} else {
		goto L77
	}
L76:
	;
	v445 = v277 << (uint(int32(1)) % 32)
	v448 = int32(2)
	v452 = v445 + int32(base.Ui32(v445+int32(3))>>(uint(v448)%32)) + v448
	v460 = v452
	v461 = base.I64_extend_i32_u(v452)
	goto L74
L77:
	;
	v460 = v454
	v461 = v455
	goto L74
L78:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+32)) = v277
	*(*int32)(unsafe.Add(mBase, uint32(v274)+8)) = v498
	v502 = int32(2)
	v505 = int32(31)
	v507 = int32(-32)
	v508 = (v498 + base.I32_wrap_i64(v440)<<(uint(v502)%32) + v505) & v507
	*(*int32)(unsafe.Add(mBase, uint32(v274)+16)) = v508
	*(*int32)(unsafe.Add(mBase, uint32(v274)+20)) = (v508 + v474<<(uint(v502)%32) + v505) & v507
	v520 = int32(1)
	goto L72
L79:
	;
	F_WebPSafeFree(m, v481)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v274)+24)) = int64(0)
	v490 = F_WebPSafeMalloc(m, v480, int32(4))
	mBase = m.M
	if v490 != 0 {
		goto L82
	} else {
		goto L83
	}
L80:
	;
	v484 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v274)+28)))
	if base.Ui64(v480) <= base.Ui64(v484) {
		v498 = v481
		goto L78
	} else {
		goto L81
	}
L81:
	;
	goto L79
L82:
	;
	*(*uint32)(unsafe.Add(mBase, uint32(v274)+28)) = uint32(v480)
	*(*int32)(unsafe.Add(mBase, uint32(v274)+24)) = v490
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(0)
	v498 = v490
	goto L78
L83:
	;
	v491 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v493 = F_WebPEncodingSetError(m, v491, int32(1))
	mBase = m.M
	v520 = v493
	goto L72
L84:
	;
	v523 = int32(2)
	v524 = *(*int32)(unsafe.Add(mBase, uint32(v274)+12))
	if v524 == v523 {
		v670 = v523
		goto L66
	} else {
		goto L85
	}
L85:
	;
	v527 = *(*int32)(unsafe.Add(mBase, uint32(v278)+92))
	v528 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v537 = base.I32_div_s(v527, int32(-20))
	v538 = *(*int32)(unsafe.Add(mBase, uint32(v275)+56))
	v539 = *(*int32)(unsafe.Add(mBase, uint32(v275)+12))
	v540 = *(*int32)(unsafe.Add(mBase, uint32(v275)+8))
	v545 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v540*int32(3)), int32(4))
	mBase = m.M
	if v545 != 0 {
		goto L87
	} else {
		goto L88
	}
L86:
	;
	if v665 != 0 {
		v670 = v523
		goto L66
	} else {
		goto L105
	}
L87:
	;
	v547 = int32(64)
	if base.B2i32(v540 < v547)&base.B2i32(v539 < v547) != 0 {
		goto L91
	} else {
		goto L92
	}
L88:
	;
	v665 = int32(0)
	goto L86
L89:
	;
	F_WebPSafeFree(m, v545)
	mBase = m.M
	v665 = int32(1)
	goto L86
L90:
	;
	v622 = *(*int32)(unsafe.Add(mBase, uint32(v275)+52))
	F_NearLossless(m, v540, v539, v622, v538, v537+int32(5), v545, v528)
	mBase = m.M
	v627 = v537 + int32(4)
	if v627 == int32(0) {
		goto L89
	} else {
		goto L101
	}
L91:
	;
	if v539 < int32(1) {
		goto L89
	} else {
		goto L94
	}
L92:
	;
	if int32(2) < v539 {
		goto L90
	} else {
		goto L93
	}
L93:
	;
	goto L91
L94:
	;
	v556 = int32(1)
	v559 = v540 << (uint(int32(2)) % 32)
	if v539 == v556 {
		v602 = int32(0)
		goto L95
	} else {
		goto L96
	}
L95:
	;
	if v539&v556 == int32(0) {
		goto L89
	} else {
		goto L100
	}
L96:
	;
	v571 = int32(0)
	v573 = v528
	goto L97
L97:
	;
	v578 = *(*int32)(unsafe.Add(mBase, uint32(v275)+52))
	v579 = *(*int32)(unsafe.Add(mBase, uint32(v275)+56))
	v581 = int32(2)
	v584 = F_memcpy(m, v573, v578+v579*v571<<(uint(v581)%32), v559)
	mBase = m.M
	v586 = *(*int32)(unsafe.Add(mBase, uint32(v275)+52))
	v587 = *(*int32)(unsafe.Add(mBase, uint32(v275)+56))
	v594 = F_memcpy(m, v584+v559, v586+v587*(v571+int32(1))<<(uint(v581)%32), v559)
	mBase = m.M
	v597 = v571 + v581
	if v539&int32(2147483646) != v597 {
		v571 = v597
		v573 = v584 + v540<<(uint(int32(3))%32)
		goto L97
	} else {
		goto L99
	}
L98:
	;
	v602 = v597
	goto L95
L99:
	;
	goto L98
L100:
	;
	v612 = int32(2)
	v615 = *(*int32)(unsafe.Add(mBase, uint32(v275)+52))
	v616 = *(*int32)(unsafe.Add(mBase, uint32(v275)+56))
	v621 = F_memcpy(m, v528+v602*v540<<(uint(v612)%32), v615+v616*v602<<(uint(v612)%32), v559)
	mBase = m.M
	goto L89
L101:
	;
	v630 = v627
	goto L102
L102:
	;
	F_NearLossless(m, v540, v539, v528, v540, v630, v545, v528)
	mBase = m.M
	v642 = v630 + int32(-1)
	if v642 != 0 {
		v630 = v642
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
	v667 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v667 != 0 {
		goto L107
	} else {
		goto L108
	}
L106:
	;
	v10242 = v265
	v10244 = v267
	goto L4
L107:
	;
	goto L106
L108:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L107
L109:
	;
	goto L65
L110:
	;
	v685 = *(*int32)(unsafe.Add(mBase, uint32(v274)+68))
	v687 = v685 + int32(-1)
	v691 = *(*int32)(unsafe.Add(mBase, uint32(v294+v687<<(uint(int32(2))%32))))
	v692 = int32(1)
	goto L117
L111:
	;
	v680 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v682 = *(*int32)(unsafe.Add(mBase, uint32(v680)+92))
	if v682 != 0 {
		goto L113
	} else {
		goto L114
	}
L112:
	;
	v10242 = v265
	v10244 = v267
	goto L4
L113:
	;
	goto L112
L114:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v680)+92)) = int32(1)
	goto L113
L115:
	;
	v844 = int32(3)
	v845 = int32(2)
	goto L143
L116:
	;
	goto L115
L117:
	;
	v705 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v706 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v706+v692 < int32(32) {
		goto L120
	} else {
		goto L121
	}
L118:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v807 + v805
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v804<<(uint(v807)%32) | v806
	goto L116
L119:
	;
	v724 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v725 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v729 = v721
	v730 = v722
	v731 = v725
	v732 = v724
	goto L123
L120:
	;
	if v706 < int32(16) {
		v804 = v692
		v805 = v692
		v806 = v705
		v807 = v706
		goto L118
	} else {
		goto L122
	}
L121:
	;
	v710 = int32(32)
	v711 = v710 - v706
	v719 = int32(base.Ui32(v692) >> (uint(v711) % 32))
	v720 = v692 - v711
	v721 = v692<<(uint(v706)%32) | v705
	v722 = v710
	goto L119
L122:
	;
	v719 = v692
	v720 = v692
	v721 = v705
	v722 = v706
	goto L119
L123:
	;
	if base.Ui32(v731+int32(2)) <= base.Ui32(v732) {
		v787 = v731
		v788 = v732
		goto L125
	} else {
		goto L126
	}
L124:
	;
	v804 = v719
	v805 = v720
	v806 = v800
	v807 = v798
	goto L118
L125:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v787))) = uint16(v729)
	v795 = v787 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v795
	v798 = v730 + int32(-16)
	v800 = int32(base.Ui32(v729) >> (uint(int32(16)) % 32))
	if int32(31) < v730 {
		v729 = v800
		v730 = v798
		v731 = v795
		v732 = v788
		goto L123
	} else {
		goto L140
	}
L126:
	;
	v741 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v742 = v732 - v741
	v745 = base.I64_extend_i32_s(v742) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v745) {
		v769 = v741
		goto L128
	} else {
		goto L129
	}
L127:
	;
	if v731 == v741 {
		goto L138
	} else {
		goto L139
	}
L128:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v769
	goto L115
L129:
	;
	v748 = v731 - v741
	v750 = v745 + base.I64_extend_i32_u(v748)
	if base.Ui64(int64(4294967295)) < base.Ui64(v750) {
		v769 = v741
		goto L128
	} else {
		goto L130
	}
L130:
	;
	v753 = base.I32_wrap_i64(v750)
	if v732 == v741 {
		goto L131
	} else {
		goto L132
	}
L131:
	;
	v760 = int32(base.Ui32(v742*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v753) < base.Ui32(v760) {
		goto L134
	} else {
		goto L135
	}
L132:
	;
	if base.Ui32(v753) <= base.Ui32(v742) {
		v787 = v731
		v788 = v732
		goto L125
	} else {
		goto L133
	}
L133:
	;
	goto L131
L134:
	;
	v762 = v760
	goto L136
L135:
	;
	v762 = v753
	goto L136
L136:
	;
	v766 = v762&int32(-1024) + int32(1024)
	v767 = F_WebPSafeMalloc(m, int64(1), v766)
	mBase = m.M
	if v767 != 0 {
		goto L127
	} else {
		goto L137
	}
L137:
	;
	v768 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v769 = v768
	goto L128
L138:
	;
	v780 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v780)
	mBase = m.M
	v782 = v767 + v766
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v782
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v767
	v787 = v767 + v748
	v788 = v782
	goto L125
L139:
	;
	v778 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v779 = F_memcpy(m, v767, v778, v748)
	mBase = m.M
	goto L138
L140:
	;
	goto L124
L141:
	;
	if v691 != 0 {
		goto L167
	} else {
		goto L168
	}
L142:
	;
	goto L141
L143:
	;
	v857 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v858 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v858+v845 < int32(32) {
		goto L146
	} else {
		goto L147
	}
L144:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v959 + v957
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v956<<(uint(v959)%32) | v958
	goto L142
L145:
	;
	v876 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v877 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v881 = v873
	v882 = v874
	v883 = v877
	v884 = v876
	goto L149
L146:
	;
	if v858 < int32(16) {
		v956 = v844
		v957 = v845
		v958 = v857
		v959 = v858
		goto L144
	} else {
		goto L148
	}
L147:
	;
	v862 = int32(32)
	v863 = v862 - v858
	v871 = int32(base.Ui32(v844) >> (uint(v863) % 32))
	v872 = v845 - v863
	v873 = v844<<(uint(v858)%32) | v857
	v874 = v862
	goto L145
L148:
	;
	v871 = v844
	v872 = v845
	v873 = v857
	v874 = v858
	goto L145
L149:
	;
	if base.Ui32(v883+int32(2)) <= base.Ui32(v884) {
		v939 = v883
		v940 = v884
		goto L151
	} else {
		goto L152
	}
L150:
	;
	v956 = v871
	v957 = v872
	v958 = v952
	v959 = v950
	goto L144
L151:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v939))) = uint16(v881)
	v947 = v939 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v947
	v950 = v882 + int32(-16)
	v952 = int32(base.Ui32(v881) >> (uint(int32(16)) % 32))
	if int32(31) < v882 {
		v881 = v952
		v882 = v950
		v883 = v947
		v884 = v940
		goto L149
	} else {
		goto L166
	}
L152:
	;
	v893 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v894 = v884 - v893
	v897 = base.I64_extend_i32_s(v894) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v897) {
		v921 = v893
		goto L154
	} else {
		goto L155
	}
L153:
	;
	if v883 == v893 {
		goto L164
	} else {
		goto L165
	}
L154:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v921
	goto L141
L155:
	;
	v900 = v883 - v893
	v902 = v897 + base.I64_extend_i32_u(v900)
	if base.Ui64(int64(4294967295)) < base.Ui64(v902) {
		v921 = v893
		goto L154
	} else {
		goto L156
	}
L156:
	;
	v905 = base.I32_wrap_i64(v902)
	if v884 == v893 {
		goto L157
	} else {
		goto L158
	}
L157:
	;
	v912 = int32(base.Ui32(v894*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v905) < base.Ui32(v912) {
		goto L160
	} else {
		goto L161
	}
L158:
	;
	if base.Ui32(v905) <= base.Ui32(v894) {
		v939 = v883
		v940 = v884
		goto L151
	} else {
		goto L159
	}
L159:
	;
	goto L157
L160:
	;
	v914 = v912
	goto L162
L161:
	;
	v914 = v905
	goto L162
L162:
	;
	v918 = v914&int32(-1024) + int32(1024)
	v919 = F_WebPSafeMalloc(m, int64(1), v918)
	mBase = m.M
	if v919 != 0 {
		goto L153
	} else {
		goto L163
	}
L163:
	;
	v920 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v921 = v920
	goto L154
L164:
	;
	v932 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v932)
	mBase = m.M
	v934 = v919 + v918
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v934
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v919
	v939 = v919 + v900
	v940 = v934
	goto L151
L165:
	;
	v930 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v931 = F_memcpy(m, v919, v930, v900)
	mBase = m.M
	goto L164
L166:
	;
	goto L150
L167:
	;
	v996 = v685
	goto L169
L168:
	;
	v996 = v687
	goto L169
L169:
	;
	if int32(17) < v685 {
		goto L170
	} else {
		goto L171
	}
L170:
	;
	v999 = v996
	goto L172
L171:
	;
	v999 = v685
	goto L172
L172:
	;
	v1001 = v999 + int32(-1)
	v1002 = int32(8)
	goto L175
L173:
	;
	if v1001 < int32(1) {
		goto L199
	} else {
		goto L200
	}
L174:
	;
	goto L173
L175:
	;
	v1014 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v1015 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v1015+v1002 < int32(32) {
		goto L178
	} else {
		goto L179
	}
L176:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v1116 + v1114
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v1113<<(uint(v1116)%32) | v1115
	goto L174
L177:
	;
	v1033 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v1034 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v1038 = v1030
	v1039 = v1031
	v1040 = v1034
	v1041 = v1033
	goto L181
L178:
	;
	if v1015 < int32(16) {
		v1113 = v1001
		v1114 = v1002
		v1115 = v1014
		v1116 = v1015
		goto L176
	} else {
		goto L180
	}
L179:
	;
	v1019 = int32(32)
	v1020 = v1019 - v1015
	v1028 = int32(base.Ui32(v1001) >> (uint(v1020) % 32))
	v1029 = v1002 - v1020
	v1030 = v1001<<(uint(v1015)%32) | v1014
	v1031 = v1019
	goto L177
L180:
	;
	v1028 = v1001
	v1029 = v1002
	v1030 = v1014
	v1031 = v1015
	goto L177
L181:
	;
	if base.Ui32(v1040+int32(2)) <= base.Ui32(v1041) {
		v1096 = v1040
		v1097 = v1041
		goto L183
	} else {
		goto L184
	}
L182:
	;
	v1113 = v1028
	v1114 = v1029
	v1115 = v1109
	v1116 = v1107
	goto L176
L183:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1096))) = uint16(v1038)
	v1104 = v1096 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v1104
	v1107 = v1039 + int32(-16)
	v1109 = int32(base.Ui32(v1038) >> (uint(int32(16)) % 32))
	if int32(31) < v1039 {
		v1038 = v1109
		v1039 = v1107
		v1040 = v1104
		v1041 = v1097
		goto L181
	} else {
		goto L198
	}
L184:
	;
	v1050 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v1051 = v1041 - v1050
	v1054 = base.I64_extend_i32_s(v1051) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1054) {
		v1078 = v1050
		goto L186
	} else {
		goto L187
	}
L185:
	;
	if v1040 == v1050 {
		goto L196
	} else {
		goto L197
	}
L186:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v1078
	goto L173
L187:
	;
	v1057 = v1040 - v1050
	v1059 = v1054 + base.I64_extend_i32_u(v1057)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1059) {
		v1078 = v1050
		goto L186
	} else {
		goto L188
	}
L188:
	;
	v1062 = base.I32_wrap_i64(v1059)
	if v1041 == v1050 {
		goto L189
	} else {
		goto L190
	}
L189:
	;
	v1069 = int32(base.Ui32(v1051*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v1062) < base.Ui32(v1069) {
		goto L192
	} else {
		goto L193
	}
L190:
	;
	if base.Ui32(v1062) <= base.Ui32(v1051) {
		v1096 = v1040
		v1097 = v1041
		goto L183
	} else {
		goto L191
	}
L191:
	;
	goto L189
L192:
	;
	v1071 = v1069
	goto L194
L193:
	;
	v1071 = v1062
	goto L194
L194:
	;
	v1075 = v1071&int32(-1024) + int32(1024)
	v1076 = F_WebPSafeMalloc(m, int64(1), v1075)
	mBase = m.M
	if v1076 != 0 {
		goto L185
	} else {
		goto L195
	}
L195:
	;
	v1077 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v1078 = v1077
	goto L186
L196:
	;
	v1089 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v1089)
	mBase = m.M
	v1091 = v1076 + v1075
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v1091
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v1076
	v1096 = v1076 + v1057
	v1097 = v1091
	goto L183
L197:
	;
	v1087 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v1088 = F_memcpy(m, v1076, v1087, v1057)
	mBase = m.M
	goto L196
L198:
	;
	goto L182
L199:
	;
	v1321 = *(*int32)(unsafe.Add(mBase, uint32(v274+int32(72))))
	*(*int32)(unsafe.Add(mBase, uint32(v267)+2112)) = v1321
	v1327 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v1330 = F_EncodeImageNoHuffman(m, v268, v267+int32(2112), v293, v297, v999, int32(1), int32(20), v286, v1327, v288, v267+int32(60))
	mBase = m.M
	if v1330 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L204
	}
L200:
	;
	v1156 = v999 << (uint(int32(2)) % 32)
	v1163 = v290 + v1156
	v1164 = v298 + v1156
	v1175 = v999
	goto L201
L201:
	;
	v1227 = *(*int32)(unsafe.Add(mBase, uint32(v1163+int32(4))))
	v1230 = *(*int32)(unsafe.Add(mBase, uint32(v1163)))
	v1231 = int32(-16711936)
	v1238 = int32(16711935)
	*(*int32)(unsafe.Add(mBase, uint32(v1164))) = (v1227|int32(16711680)-v1230&v1231)&v1231 | (v1227|int32(_a_F_EncodeStreamHook_1)-v1230&v1238)&v1238
	v1245 = int32(-4)
	v1250 = v1175 + int32(-1)
	if base.Ui32(int32(1)) < base.Ui32(v1250) {
		v1163 = v1163 + v1245
		v1164 = v1164 + v1245
		v1175 = v1250
		goto L201
	} else {
		goto L203
	}
L202:
	;
	goto L199
L203:
	;
	goto L202
L204:
	;
	v1333 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v1334 = *(*int32)(unsafe.Add(mBase, uint32(v1333)+8))
	v1336 = int32(3)
	v1338 = *(*int32)(unsafe.Add(mBase, uint32(v274)+68))
	if v1338 < v1336 {
		goto L205
	} else {
		goto L206
	}
L205:
	;
	v1341 = v1336
	goto L207
L206:
	;
	v1341 = int32(2)
	goto L207
L207:
	;
	if v1338 < int32(5) {
		goto L208
	} else {
		goto L209
	}
L208:
	;
	v1346 = v1341
	goto L210
L209:
	;
	v1346 = base.B2i32(v1338 < int32(17))
	goto L210
L210:
	;
	v1351 = int32(base.Ui32(v1334+int32(1)<<(uint(v1346)%32)+int32(-1)) >> (uint(v1346) % 32))
	v1352 = *(*int32)(unsafe.Add(mBase, uint32(v1333)+12))
	v1359 = base.I64_extend_i32_s(v1352) * base.I64_extend_i32_s(v1351)
	v1360 = *(*int32)(unsafe.Add(mBase, uint32(v274)+60))
	if v1360 == int32(0) {
		goto L214
	} else {
		goto L215
	}
L211:
	;
	if v1439 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L223
	}
L212:
	;
	v1399 = v1359 + v1394 + v1395 + int64(16)
	v1400 = *(*int32)(unsafe.Add(mBase, uint32(v274)+24))
	if v1400 == int32(0) {
		goto L218
	} else {
		goto L219
	}
L213:
	;
	v1382 = int32(3)
	v1384 = int32(2)
	v1393 = v1379
	v1394 = v1380
	v1395 = base.I64_extend_i32_u(int32(base.Ui32(v1352+v1382)>>(uint(v1384)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v1351+v1382)>>(uint(v1384)%32)))
	goto L212
L214:
	;
	v1373 = int32(0)
	v1374 = int64(0)
	v1376 = *(*int32)(unsafe.Add(mBase, uint32(v274)+52))
	if v1376 == v1373 {
		v1393 = v1373
		v1394 = v1374
		v1395 = v1374
		goto L212
	} else {
		goto L216
	}
L215:
	;
	v1364 = v1351 << (uint(int32(1)) % 32)
	v1367 = int32(2)
	v1371 = v1364 + int32(base.Ui32(v1364+int32(3))>>(uint(v1367)%32)) + v1367
	v1379 = v1371
	v1380 = base.I64_extend_i32_u(v1371)
	goto L213
L216:
	;
	v1379 = v1373
	v1380 = v1374
	goto L213
L217:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+32)) = v1351
	*(*int32)(unsafe.Add(mBase, uint32(v274)+8)) = v1417
	v1421 = int32(2)
	v1424 = int32(31)
	v1426 = int32(-32)
	v1427 = (v1417 + base.I32_wrap_i64(v1359)<<(uint(v1421)%32) + v1424) & v1426
	*(*int32)(unsafe.Add(mBase, uint32(v274)+16)) = v1427
	*(*int32)(unsafe.Add(mBase, uint32(v274)+20)) = (v1427 + v1393<<(uint(v1421)%32) + v1424) & v1426
	v1439 = int32(1)
	goto L211
L218:
	;
	F_WebPSafeFree(m, v1400)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v274)+24)) = int64(0)
	v1409 = F_WebPSafeMalloc(m, v1399, int32(4))
	mBase = m.M
	if v1409 != 0 {
		goto L221
	} else {
		goto L222
	}
L219:
	;
	v1403 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v274)+28)))
	if base.Ui64(v1399) <= base.Ui64(v1403) {
		v1417 = v1400
		goto L217
	} else {
		goto L220
	}
L220:
	;
	goto L218
L221:
	;
	*(*uint32)(unsafe.Add(mBase, uint32(v274)+28)) = uint32(v1399)
	*(*int32)(unsafe.Add(mBase, uint32(v274)+24)) = v1409
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(0)
	v1417 = v1409
	goto L217
L222:
	;
	v1410 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v1412 = F_WebPEncodingSetError(m, v1410, int32(1))
	mBase = m.M
	v1439 = v1412
	goto L211
L223:
	;
	v1442 = *(*int32)(unsafe.Add(mBase, uint32(v274)+32))
	v1443 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v1444 = *(*int32)(unsafe.Add(mBase, uint32(v1333)+56))
	v1445 = *(*int32)(unsafe.Add(mBase, uint32(v1333)+52))
	v1446 = base.I64_extend_i32_s(v1334)
	v1447 = int32(1)
	if v1446 == int64(0) {
		goto L226
	} else {
		goto L227
	}
L224:
	;
	if v1468 == int32(0) {
		goto L62
	} else {
		goto L230
	}
L225:
	;
	goto L224
L226:
	;
	v1466 = F_malloc(m, base.I32_wrap_i64(v1446)*v1447)
	mBase = m.M
	v1468 = v1466
	goto L225
L227:
	;
	v1454 = base.I64_div_u_s(int64(2147418112), v1446)
	v1455 = int32(0)
	v1456 = base.I64_extend_i32_u(v1447)
	if base.Ui64(int64(4294967295)) < base.Ui64(v1456*v1446) {
		v1468 = v1455
		goto L225
	} else {
		goto L228
	}
L228:
	;
	if base.Ui64(v1454) < base.Ui64(v1456) {
		v1468 = v1455
		goto L225
	} else {
		goto L229
	}
L229:
	;
	goto L226
L230:
	;
	if int32(3) < v1338 {
		goto L232
	} else {
		goto L233
	}
L231:
	;
	F_free(m, v1468)
	mBase = m.M
	goto L438
L232:
	;
	goto L255
L233:
	;
	if v1352 < int32(1) {
		goto L231
	} else {
		goto L234
	}
L234:
	;
	v1478 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v1479 = int32(0)
	v1520 = v1479
	v1521 = v1478
	v1526 = v1443
	v1527 = v1445
	v1530 = v1479
	goto L235
L235:
	;
	if v1334 < int32(1) {
		v1677 = v1520
		v1678 = v1521
		goto L237
	} else {
		goto L238
	}
L237:
	;
	v1704 = m.G71
	v1705 = *(*int32)(unsafe.Add(mBase, uint32(v1704)))
	m.T0[v1705].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v1526)
	mBase = m.M
	v1712 = v1530 + int32(1)
	if v1712 != v1352 {
		v1520 = v1677
		v1521 = v1678
		v1526 = v1526 + v1442<<(uint(int32(2))%32)
		v1527 = v1527 + v1444<<(uint(int32(2))%32)
		v1530 = v1712
		goto L235
	} else {
		goto L252
	}
L238:
	;
	v1553 = v1527
	v1554 = v1468
	v1566 = v1334
	v1588 = v1520
	v1589 = v1521
	goto L239
L239:
	;
	v1615 = *(*int32)(unsafe.Add(mBase, uint32(v1553)))
	if v1615 == v1589 {
		v1629 = v1588
		v1630 = v1589
		goto L241
	} else {
		goto L242
	}
L240:
	;
	v1677 = v1629
	v1678 = v1630
	goto L237
L241:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1554))) = uint8(v1629)
	v1637 = v1566 + int32(-1)
	if v1637 != 0 {
		v1553 = v1553 + int32(4)
		v1554 = v1554 + int32(1)
		v1566 = v1637
		v1588 = v1629
		v1589 = v1630
		goto L239
	} else {
		goto L251
	}
L242:
	;
	v1617 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	if v1617 != v1615 {
		goto L244
	} else {
		goto L245
	}
L243:
	;
	v1629 = v1628
	v1630 = v1615
	goto L241
L244:
	;
	v1620 = *(*int32)(unsafe.Add(mBase, uint32(v274)+76))
	if v1620 != v1615 {
		goto L246
	} else {
		goto L247
	}
L245:
	;
	v1628 = int32(0)
	goto L243
L246:
	;
	v1625 = *(*int32)(unsafe.Add(mBase, uint32(v274)+80))
	if v1625 == v1615 {
		goto L248
	} else {
		goto L249
	}
L247:
	;
	v1628 = int32(1)
	goto L243
L248:
	;
	v1627 = int32(2)
	goto L250
L249:
	;
	v1627 = int32(3)
	goto L250
L250:
	;
	v1628 = v1627
	goto L243
L251:
	;
	goto L240
L252:
	;
	goto L231
L253:
	;
	v1843 = v294
	v1844 = int32(0)
	goto L267
L255:
	;
	base.MemoryFill(m, v267+int32(2112), int32(255), int32(_a_F_EncodeStreamHook_2))
	goto L253
L267:
	;
	v1907 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1843)+1)))
	v1910 = v267 + int32(2112) + v1907<<(uint(int32(1))%32)
	v1911 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1910))))
	if v1911 == int32(_a_F_EncodeStreamHook_3) {
		goto L269
	} else {
		goto L270
	}
L268:
	;
	if v1352 < int32(1) {
		goto L231
	} else {
		goto L411
	}
L269:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v1910))) = uint16(v1844)
	v3699 = v1844 + int32(1)
	if v1338 != v3699 {
		v1843 = v1843 + int32(4)
		v1844 = v3699
		goto L267
	} else {
		goto L410
	}
L270:
	;
	goto L273
L271:
	;
	v2043 = v294
	v2044 = int32(0)
	goto L287
L273:
	;
	base.MemoryFill(m, v267+int32(2112), int32(255), int32(_a_F_EncodeStreamHook_2))
	goto L271
L285:
	;
	v3255 = F_memcpy(m, v267+int32(64), v294, v1338<<(uint(int32(2))%32))
	mBase = m.M
	v3257 = m.G2
	F_qsort(m, v3255, v1338, int32(4), v3257+int32(170))
	mBase = m.M
	if v1338 == int32(0) {
		goto L364
	} else {
		goto L365
	}
L286:
	;
	v2957 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v2958 = int32(0)
	v2976 = v2957
	v2977 = v2958
	v3005 = v1443
	v3006 = v1445
	v3008 = v2958
	goto L346
L287:
	;
	v2107 = *(*int32)(unsafe.Add(mBase, uint32(v2043)))
	v2116 = v267 + int32(2112) + int32(base.Ui32(v2107&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)
	v2117 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2116))))
	if v2117 == int32(_a_F_EncodeStreamHook_3) {
		goto L289
	} else {
		goto L290
	}
L288:
	;
	if v1352 < int32(1) {
		goto L231
	} else {
		goto L319
	}
L289:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2116))) = uint16(v2044)
	v2497 = v2044 + int32(1)
	if v1338 != v2497 {
		v2043 = v2043 + int32(4)
		v2044 = v2497
		goto L287
	} else {
		goto L318
	}
L290:
	;
	goto L293
L291:
	;
	v2249 = v294
	v2250 = int32(0)
	goto L305
L293:
	;
	base.MemoryFill(m, v267+int32(2112), int32(255), int32(_a_F_EncodeStreamHook_2))
	goto L291
L305:
	;
	v2313 = *(*int32)(unsafe.Add(mBase, uint32(v2249)))
	v2322 = v267 + int32(2112) + int32(base.Ui32(v2313&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)
	v2323 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v2322))))
	if v2323 != int32(_a_F_EncodeStreamHook_3) {
		goto L285
	} else {
		goto L307
	}
L306:
	;
	if v1352 < int32(1) {
		goto L231
	} else {
		goto L309
	}
L307:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v2322))) = uint16(v2250)
	v2330 = v2250 + int32(1)
	if v1338 != v2330 {
		v2249 = v2249 + int32(4)
		v2250 = v2330
		goto L305
	} else {
		goto L308
	}
L308:
	;
	goto L306
L309:
	;
	if int32(0) < v1334 {
		goto L286
	} else {
		goto L310
	}
L310:
	;
	v2336 = int32(1)
	if v1352 == v2336 {
		v2466 = v1443
		goto L311
	} else {
		goto L312
	}
L311:
	;
	if v1352&v2336 == int32(0) {
		goto L231
	} else {
		goto L316
	}
L312:
	;
	v2350 = v1352 & int32(2147483646)
	v2391 = v1443
	goto L313
L313:
	;
	v2412 = m.G71
	v2413 = *(*int32)(unsafe.Add(mBase, uint32(v2412)))
	m.T0[v2413].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2391)
	mBase = m.M
	v2416 = *(*int32)(unsafe.Add(mBase, uint32(v2412)))
	m.T0[v2416].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2391+v1442<<(uint(int32(2))%32))
	mBase = m.M
	v2418 = v2391 + v1442<<(uint(int32(3))%32)
	v2420 = v2350 + int32(-2)
	if v2420 != 0 {
		v2350 = v2420
		v2391 = v2418
		goto L313
	} else {
		goto L315
	}
L314:
	;
	v2466 = v2418
	goto L311
L315:
	;
	goto L314
L316:
	;
	v2489 = m.G71
	v2490 = *(*int32)(unsafe.Add(mBase, uint32(v2489)))
	m.T0[v2490].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2466)
	mBase = m.M
	F_free(m, v1468)
	mBase = m.M
	goto L317
L317:
	;
	goto L61
L318:
	;
	goto L288
L319:
	;
	if int32(0) < v1334 {
		goto L320
	} else {
		goto L321
	}
L320:
	;
	v2666 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v2667 = int32(0)
	v2685 = v2666
	v2686 = v2667
	v2714 = v1443
	v2715 = v1445
	v2717 = v2667
	goto L329
L321:
	;
	v2503 = int32(1)
	if v1352 == v2503 {
		v2633 = v1443
		goto L322
	} else {
		goto L323
	}
L322:
	;
	if v1352&v2503 == int32(0) {
		goto L231
	} else {
		goto L327
	}
L323:
	;
	v2517 = v1352 & int32(2147483646)
	v2558 = v1443
	goto L324
L324:
	;
	v2579 = m.G71
	v2580 = *(*int32)(unsafe.Add(mBase, uint32(v2579)))
	m.T0[v2580].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2558)
	mBase = m.M
	v2583 = *(*int32)(unsafe.Add(mBase, uint32(v2579)))
	m.T0[v2583].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2558+v1442<<(uint(int32(2))%32))
	mBase = m.M
	v2585 = v2558 + v1442<<(uint(int32(3))%32)
	v2587 = v2517 + int32(-2)
	if v2587 != 0 {
		v2517 = v2587
		v2558 = v2585
		goto L324
	} else {
		goto L326
	}
L325:
	;
	v2633 = v2585
	goto L322
L326:
	;
	goto L325
L327:
	;
	v2656 = m.G71
	v2657 = *(*int32)(unsafe.Add(mBase, uint32(v2656)))
	m.T0[v2657].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2633)
	mBase = m.M
	F_free(m, v1468)
	mBase = m.M
	goto L328
L328:
	;
	goto L61
L329:
	;
	if v1334 != int32(1) {
		goto L332
	} else {
		goto L333
	}
L331:
	;
	if v1334&int32(1) == int32(0) {
		v2939 = v2865
		v2940 = v2866
		goto L341
	} else {
		goto L342
	}
L332:
	;
	v2743 = v2715
	v2744 = int32(0)
	v2755 = v2685
	v2756 = v2686
	goto L334
L333:
	;
	v2854 = int32(0)
	v2865 = v2685
	v2866 = v2686
	goto L331
L334:
	;
	v2805 = *(*int32)(unsafe.Add(mBase, uint32(v2743)))
	if v2805 == v2755 {
		v2819 = v2755
		v2820 = v2756
		goto L336
	} else {
		goto L337
	}
L335:
	;
	v2854 = v2847
	v2865 = v2839
	v2866 = v2840
	goto L331
L336:
	;
	v2821 = v1468 + v2744
	*(*uint8)(unsafe.Add(mBase, uint32(v2821))) = uint8(v2820)
	v2825 = *(*int32)(unsafe.Add(mBase, uint32(v2743+int32(4))))
	if v2825 == v2819 {
		v2839 = v2819
		v2840 = v2820
		goto L338
	} else {
		goto L339
	}
L337:
	;
	v2818 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v2805&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)))))
	v2819 = v2805
	v2820 = v2818
	goto L336
L338:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v2821+int32(1)))) = uint8(v2840)
	v2847 = v2744 + int32(2)
	if v1334&int32(2147483646) != v2847 {
		v2743 = v2743 + int32(8)
		v2744 = v2847
		v2755 = v2839
		v2756 = v2840
		goto L334
	} else {
		goto L340
	}
L339:
	;
	v2838 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v2825&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)))))
	v2839 = v2825
	v2840 = v2838
	goto L338
L340:
	;
	goto L335
L341:
	;
	v2941 = m.G71
	v2942 = *(*int32)(unsafe.Add(mBase, uint32(v2941)))
	m.T0[v2942].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v2714)
	mBase = m.M
	v2949 = v2717 + int32(1)
	if v2949 != v1352 {
		v2685 = v2939
		v2686 = v2940
		v2714 = v2714 + v1442<<(uint(int32(2))%32)
		v2715 = v2715 + v1444<<(uint(int32(2))%32)
		v2717 = v2949
		goto L329
	} else {
		goto L345
	}
L342:
	;
	v2920 = *(*int32)(unsafe.Add(mBase, uint32(v2715+v2854<<(uint(int32(2))%32))))
	if v2920 == v2865 {
		v2934 = v2865
		v2935 = v2866
		goto L343
	} else {
		goto L344
	}
L343:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1468+v2854))) = uint8(v2935)
	v2939 = v2934
	v2940 = v2935
	goto L341
L344:
	;
	v2933 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v2920&int32(16777215)*int32(-72723225))>>(uint(int32(20))%32))&int32(4094)))))
	v2934 = v2920
	v2935 = v2933
	goto L343
L345:
	;
	goto L231
L346:
	;
	if v1334 != int32(1) {
		goto L349
	} else {
		goto L350
	}
L348:
	;
	if v1334&int32(1) == int32(0) {
		v3230 = v3156
		v3231 = v3157
		goto L358
	} else {
		goto L359
	}
L349:
	;
	v3034 = v3006
	v3035 = int32(0)
	v3046 = v2976
	v3047 = v2977
	goto L351
L350:
	;
	v3145 = int32(0)
	v3156 = v2976
	v3157 = v2977
	goto L348
L351:
	;
	v3096 = *(*int32)(unsafe.Add(mBase, uint32(v3034)))
	if v3096 == v3046 {
		v3110 = v3046
		v3111 = v3047
		goto L353
	} else {
		goto L354
	}
L352:
	;
	v3145 = v3138
	v3156 = v3130
	v3157 = v3131
	goto L348
L353:
	;
	v3112 = v1468 + v3035
	*(*uint8)(unsafe.Add(mBase, uint32(v3112))) = uint8(v3111)
	v3116 = *(*int32)(unsafe.Add(mBase, uint32(v3034+int32(4))))
	if v3116 == v3110 {
		v3130 = v3110
		v3131 = v3111
		goto L355
	} else {
		goto L356
	}
L354:
	;
	v3109 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v3096&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)))))
	v3110 = v3096
	v3111 = v3109
	goto L353
L355:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v3112+int32(1)))) = uint8(v3131)
	v3138 = v3035 + int32(2)
	if v1334&int32(2147483646) != v3138 {
		v3034 = v3034 + int32(8)
		v3035 = v3138
		v3046 = v3130
		v3047 = v3131
		goto L351
	} else {
		goto L357
	}
L356:
	;
	v3129 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v3116&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)))))
	v3130 = v3116
	v3131 = v3129
	goto L355
L357:
	;
	goto L352
L358:
	;
	v3232 = m.G71
	v3233 = *(*int32)(unsafe.Add(mBase, uint32(v3232)))
	m.T0[v3233].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3005)
	mBase = m.M
	v3240 = v3008 + int32(1)
	if v3240 != v1352 {
		v2976 = v3230
		v2977 = v3231
		v3005 = v3005 + v1442<<(uint(int32(2))%32)
		v3006 = v3006 + v1444<<(uint(int32(2))%32)
		v3008 = v3240
		goto L346
	} else {
		goto L362
	}
L359:
	;
	v3211 = *(*int32)(unsafe.Add(mBase, uint32(v3006+v3145<<(uint(int32(2))%32))))
	if v3211 == v3156 {
		v3225 = v3156
		v3226 = v3157
		goto L360
	} else {
		goto L361
	}
L360:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1468+v3145))) = uint8(v3226)
	v3230 = v3225
	v3231 = v3226
	goto L358
L361:
	;
	v3224 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v3211&int32(16777215)*int32(2147483647))>>(uint(int32(20))%32))&int32(4094)))))
	v3225 = v3211
	v3226 = v3224
	goto L360
L362:
	;
	goto L231
L363:
	;
	if v1352 < int32(1) {
		goto L231
	} else {
		goto L380
	}
L364:
	;
	goto L363
L365:
	;
	v3269 = int32(0)
	goto L366
L366:
	;
	v3276 = *(*int32)(unsafe.Add(mBase, uint32(v3255)))
	v3280 = *(*int32)(unsafe.Add(mBase, uint32(v294+v3269<<(uint(int32(2))%32))))
	if v3276 == v3280 {
		v3307 = int32(0)
		goto L368
	} else {
		goto L369
	}
L367:
	;
	goto L364
L368:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v267+int32(1088)+v3307<<(uint(int32(2))%32)))) = v3269
	v3321 = v3269 + int32(1)
	if v3321 != v1338 {
		v3269 = v3321
		goto L366
	} else {
		goto L379
	}
L369:
	;
	v3290 = int32(0)
	v3291 = v1338
	goto L370
L370:
	;
	v3296 = (v3291 + v3290) >> (uint(int32(1)) % 32)
	v3300 = *(*int32)(unsafe.Add(mBase, uint32(v3255+v3296<<(uint(int32(2))%32))))
	v3301 = base.B2i32(base.Ui32(v3300) < base.Ui32(v3280))
	if base.Ui32(v3300) < base.Ui32(v3280) {
		goto L372
	} else {
		goto L373
	}
L371:
	;
	v3307 = v3296
	goto L368
L372:
	;
	v3302 = v3291
	goto L374
L373:
	;
	v3302 = v3296
	goto L374
L374:
	;
	if base.Ui32(v3300) < base.Ui32(v3280) {
		goto L375
	} else {
		goto L376
	}
L375:
	;
	v3303 = v3296
	goto L377
L376:
	;
	v3303 = v3290
	goto L377
L377:
	;
	if v3300 != v3280 {
		v3290 = v3303
		v3291 = v3302
		goto L370
	} else {
		goto L378
	}
L378:
	;
	goto L371
L379:
	;
	goto L367
L380:
	;
	if int32(0) < v1334 {
		goto L381
	} else {
		goto L382
	}
L381:
	;
	v3497 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v3498 = int32(0)
	v3517 = v3497
	v3540 = v3498
	v3545 = v1443
	v3546 = v1445
	v3549 = v3498
	goto L390
L382:
	;
	v3338 = int32(1)
	if v1352 == v3338 {
		v3468 = v1443
		goto L383
	} else {
		goto L384
	}
L383:
	;
	if v1352&v3338 == int32(0) {
		goto L231
	} else {
		goto L388
	}
L384:
	;
	v3352 = v1352 & int32(2147483646)
	v3393 = v1443
	goto L385
L385:
	;
	v3414 = m.G71
	v3415 = *(*int32)(unsafe.Add(mBase, uint32(v3414)))
	m.T0[v3415].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3393)
	mBase = m.M
	v3418 = *(*int32)(unsafe.Add(mBase, uint32(v3414)))
	m.T0[v3418].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3393+v1442<<(uint(int32(2))%32))
	mBase = m.M
	v3420 = v3393 + v1442<<(uint(int32(3))%32)
	v3422 = v3352 + int32(-2)
	if v3422 != 0 {
		v3352 = v3422
		v3393 = v3420
		goto L385
	} else {
		goto L387
	}
L386:
	;
	v3468 = v3420
	goto L383
L387:
	;
	goto L386
L388:
	;
	v3491 = m.G71
	v3492 = *(*int32)(unsafe.Add(mBase, uint32(v3491)))
	m.T0[v3492].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3468)
	mBase = m.M
	F_free(m, v1468)
	mBase = m.M
	goto L389
L389:
	;
	goto L61
L390:
	;
	v3571 = v3546
	v3572 = int32(0)
	v3584 = v3517
	v3607 = v3540
	goto L392
L392:
	;
	v3633 = *(*int32)(unsafe.Add(mBase, uint32(v3571)))
	if v3633 == v3584 {
		v3676 = v3584
		v3677 = v3607
		goto L394
	} else {
		goto L395
	}
L393:
	;
	v3685 = m.G71
	v3686 = *(*int32)(unsafe.Add(mBase, uint32(v3685)))
	m.T0[v3686].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3545)
	mBase = m.M
	v3693 = v3549 + int32(1)
	if v3693 != v1352 {
		v3517 = v3676
		v3540 = v3677
		v3545 = v3545 + v1442<<(uint(int32(2))%32)
		v3546 = v3546 + v1444<<(uint(int32(2))%32)
		v3549 = v3693
		goto L390
	} else {
		goto L409
	}
L394:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1468+v3572))) = uint8(v3677)
	v3683 = v3572 + int32(1)
	if v1334 != v3683 {
		v3571 = v3571 + int32(4)
		v3572 = v3683
		v3584 = v3676
		v3607 = v3677
		goto L392
	} else {
		goto L408
	}
L395:
	;
	v3638 = v267 + int32(64)
	v3644 = *(*int32)(unsafe.Add(mBase, uint32(v3638)))
	if v3644 == v3633 {
		v3668 = int32(0)
		goto L397
	} else {
		goto L398
	}
L396:
	;
	v3675 = *(*int32)(unsafe.Add(mBase, uint32(v267+int32(1088)+v3668<<(uint(int32(2))%32))))
	v3676 = v3633
	v3677 = v3675
	goto L394
L397:
	;
	goto L396
L398:
	;
	v3649 = v1338
	v3651 = int32(0)
	goto L399
L399:
	;
	v3656 = (v3649 + v3651) >> (uint(int32(1)) % 32)
	v3660 = *(*int32)(unsafe.Add(mBase, uint32(v3638+v3656<<(uint(int32(2))%32))))
	v3661 = base.B2i32(base.Ui32(v3660) < base.Ui32(v3633))
	if base.Ui32(v3660) < base.Ui32(v3633) {
		goto L401
	} else {
		goto L402
	}
L400:
	;
	v3668 = v3656
	goto L397
L401:
	;
	v3662 = v3649
	goto L403
L402:
	;
	v3662 = v3656
	goto L403
L403:
	;
	if base.Ui32(v3660) < base.Ui32(v3633) {
		goto L404
	} else {
		goto L405
	}
L404:
	;
	v3663 = v3656
	goto L406
L405:
	;
	v3663 = v3651
	goto L406
L406:
	;
	if v3660 != v3633 {
		v3649 = v3662
		v3651 = v3663
		goto L399
	} else {
		goto L407
	}
L407:
	;
	goto L400
L408:
	;
	goto L393
L409:
	;
	goto L231
L410:
	;
	goto L268
L411:
	;
	if int32(0) < v1334 {
		goto L412
	} else {
		goto L413
	}
L412:
	;
	v3868 = *(*int32)(unsafe.Add(mBase, uint32(v294)))
	v3869 = int32(0)
	v3887 = v3869
	v3888 = v3868
	v3916 = v1443
	v3917 = v1445
	v3919 = v3869
	goto L421
L413:
	;
	v3705 = int32(1)
	if v1352 == v3705 {
		v3835 = v1443
		goto L414
	} else {
		goto L415
	}
L414:
	;
	if v1352&v3705 == int32(0) {
		goto L231
	} else {
		goto L419
	}
L415:
	;
	v3719 = v1352 & int32(2147483646)
	v3760 = v1443
	goto L416
L416:
	;
	v3781 = m.G71
	v3782 = *(*int32)(unsafe.Add(mBase, uint32(v3781)))
	m.T0[v3782].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3760)
	mBase = m.M
	v3785 = *(*int32)(unsafe.Add(mBase, uint32(v3781)))
	m.T0[v3785].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3760+v1442<<(uint(int32(2))%32))
	mBase = m.M
	v3787 = v3760 + v1442<<(uint(int32(3))%32)
	v3789 = v3719 + int32(-2)
	if v3789 != 0 {
		v3719 = v3789
		v3760 = v3787
		goto L416
	} else {
		goto L418
	}
L417:
	;
	v3835 = v3787
	goto L414
L418:
	;
	goto L417
L419:
	;
	v3858 = m.G71
	v3859 = *(*int32)(unsafe.Add(mBase, uint32(v3858)))
	m.T0[v3859].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3835)
	mBase = m.M
	F_free(m, v1468)
	mBase = m.M
	goto L420
L420:
	;
	goto L61
L421:
	;
	if v1334 != int32(1) {
		goto L424
	} else {
		goto L425
	}
L422:
	;
	goto L231
L423:
	;
	if v1334&int32(1) == int32(0) {
		v4129 = v4059
		v4130 = v4060
		goto L433
	} else {
		goto L434
	}
L424:
	;
	v3945 = v3917
	v3946 = int32(0)
	v3957 = v3887
	v3958 = v3888
	goto L426
L425:
	;
	v4048 = int32(0)
	v4059 = v3887
	v4060 = v3888
	goto L423
L426:
	;
	v4007 = *(*int32)(unsafe.Add(mBase, uint32(v3945)))
	if v4007 == v3958 {
		v4017 = v3957
		v4018 = v3958
		goto L428
	} else {
		goto L429
	}
L427:
	;
	v4048 = v4041
	v4059 = v4033
	v4060 = v4034
	goto L423
L428:
	;
	v4019 = v1468 + v3946
	*(*uint8)(unsafe.Add(mBase, uint32(v4019))) = uint8(v4017)
	v4023 = *(*int32)(unsafe.Add(mBase, uint32(v3945+int32(4))))
	if v4023 == v4018 {
		v4033 = v4017
		v4034 = v4018
		goto L430
	} else {
		goto L431
	}
L429:
	;
	v4016 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v4007)>>(uint(int32(7))%32))&int32(510)))))
	v4017 = v4016
	v4018 = v4007
	goto L428
L430:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4019+int32(1)))) = uint8(v4033)
	v4041 = v3946 + int32(2)
	if v1334&int32(2147483646) != v4041 {
		v3945 = v3945 + int32(8)
		v3946 = v4041
		v3957 = v4033
		v3958 = v4034
		goto L426
	} else {
		goto L432
	}
L431:
	;
	v4032 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v4023)>>(uint(int32(7))%32))&int32(510)))))
	v4033 = v4032
	v4034 = v4023
	goto L430
L432:
	;
	goto L427
L433:
	;
	v4131 = m.G71
	v4132 = *(*int32)(unsafe.Add(mBase, uint32(v4131)))
	m.T0[v4132].(func(*base.Module, int32, int32, int32, int32))(m, v1468, v1334, v1346, v3916)
	mBase = m.M
	v4139 = v3919 + int32(1)
	if v4139 != v1352 {
		v3887 = v4129
		v3888 = v4130
		v3916 = v3916 + v1442<<(uint(int32(2))%32)
		v3917 = v3917 + v1444<<(uint(int32(2))%32)
		v3919 = v4139
		goto L421
	} else {
		goto L437
	}
L434:
	;
	v4114 = *(*int32)(unsafe.Add(mBase, uint32(v3917+v4048<<(uint(int32(2))%32))))
	if v4114 == v4060 {
		v4124 = v4059
		v4125 = v4060
		goto L435
	} else {
		goto L436
	}
L435:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v1468+v4048))) = uint8(v4124)
	v4129 = v4124
	v4130 = v4125
	goto L433
L436:
	;
	v4123 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v267+int32(2112)+int32(base.Ui32(v4114)>>(uint(int32(7))%32))&int32(510)))))
	v4124 = v4123
	v4125 = v4114
	goto L435
L437:
	;
	goto L422
L438:
	;
	goto L61
L439:
	;
	goto L63
L440:
	;
	if v4301 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L452
	}
L441:
	;
	v4261 = v4221 + v4256 + v4257 + int64(16)
	v4262 = *(*int32)(unsafe.Add(mBase, uint32(v274)+24))
	if v4262 == int32(0) {
		goto L447
	} else {
		goto L448
	}
L442:
	;
	v4244 = int32(3)
	v4246 = int32(2)
	v4255 = v4241
	v4256 = v4242
	v4257 = base.I64_extend_i32_u(int32(base.Ui32(v4214+v4244)>>(uint(v4246)%32))) * base.I64_extend_i32_u(int32(base.Ui32(v4213+v4244)>>(uint(v4246)%32)))
	goto L441
L443:
	;
	v4235 = int32(0)
	v4236 = int64(0)
	v4238 = *(*int32)(unsafe.Add(mBase, uint32(v274)+52))
	if v4238 == v4235 {
		v4255 = v4235
		v4256 = v4236
		v4257 = v4236
		goto L441
	} else {
		goto L445
	}
L444:
	;
	v4226 = v4213 << (uint(int32(1)) % 32)
	v4229 = int32(2)
	v4233 = v4226 + int32(base.Ui32(v4226+int32(3))>>(uint(v4229)%32)) + v4229
	v4241 = v4233
	v4242 = base.I64_extend_i32_u(v4233)
	goto L442
L445:
	;
	v4241 = v4235
	v4242 = v4236
	goto L442
L446:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+32)) = v4213
	*(*int32)(unsafe.Add(mBase, uint32(v274)+8)) = v4279
	v4283 = int32(2)
	v4286 = int32(31)
	v4288 = int32(-32)
	v4289 = (v4279 + base.I32_wrap_i64(v4221)<<(uint(v4283)%32) + v4286) & v4288
	*(*int32)(unsafe.Add(mBase, uint32(v274)+16)) = v4289
	*(*int32)(unsafe.Add(mBase, uint32(v274)+20)) = (v4289 + v4255<<(uint(v4283)%32) + v4286) & v4288
	v4301 = int32(1)
	goto L440
L447:
	;
	F_WebPSafeFree(m, v4262)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v274)+24)) = int64(0)
	v4271 = F_WebPSafeMalloc(m, v4261, int32(4))
	mBase = m.M
	if v4271 != 0 {
		goto L450
	} else {
		goto L451
	}
L448:
	;
	v4265 = int64(*(*uint32)(unsafe.Add(mBase, uint32(v274)+28)))
	if base.Ui64(v4261) <= base.Ui64(v4265) {
		v4279 = v4262
		goto L446
	} else {
		goto L449
	}
L449:
	;
	goto L447
L450:
	;
	*(*uint32)(unsafe.Add(mBase, uint32(v274)+28)) = uint32(v4261)
	*(*int32)(unsafe.Add(mBase, uint32(v274)+24)) = v4271
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(0)
	v4279 = v4271
	goto L446
L451:
	;
	v4272 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v4274 = F_WebPEncodingSetError(m, v4272, int32(1))
	mBase = m.M
	v4301 = v4274
	goto L440
L452:
	;
	v4304 = *(*int32)(unsafe.Add(mBase, uint32(v274)+12))
	if v4304 == int32(1) {
		v4668 = v287
		goto L60
	} else {
		goto L453
	}
L453:
	;
	if v4214 < int32(1) {
		goto L454
	} else {
		goto L455
	}
L454:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+12)) = int32(1)
	v4668 = v287
	goto L60
L455:
	;
	v4309 = int32(1)
	v4312 = v4213 << (uint(int32(2)) % 32)
	v4313 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v4314 = *(*int32)(unsafe.Add(mBase, uint32(v4212)+52))
	if v4214 == v4309 {
		v4405 = v4314
		v4418 = v4313
		goto L456
	} else {
		goto L457
	}
L456:
	;
	if v4214&v4309 == int32(0) {
		goto L454
	} else {
		goto L461
	}
L457:
	;
	v4325 = v4314
	v4338 = v4313
	v4361 = v4214 & int32(2147483646)
	goto L458
L458:
	;
	v4387 = F_memcpy(m, v4338, v4325, v4312)
	mBase = m.M
	v4389 = *(*int32)(unsafe.Add(mBase, uint32(v4212)+56))
	v4390 = int32(2)
	v4392 = v4325 + v4389<<(uint(v4390)%32)
	v4393 = F_memcpy(m, v4387+v4312, v4392, v4312)
	mBase = m.M
	v4394 = *(*int32)(unsafe.Add(mBase, uint32(v4212)+56))
	v4397 = v4392 + v4394<<(uint(v4390)%32)
	v4398 = v4387 + v4213<<(uint(int32(3))%32)
	v4400 = v4361 + int32(-2)
	if v4400 != 0 {
		v4325 = v4397
		v4338 = v4398
		v4361 = v4400
		goto L458
	} else {
		goto L460
	}
L459:
	;
	v4405 = v4397
	v4418 = v4398
	goto L456
L460:
	;
	goto L459
L461:
	;
	v4469 = F_memcpy(m, v4418, v4405, v4312)
	mBase = m.M
	goto L454
L462:
	;
	v10242 = v265
	v10244 = v267
	goto L4
L463:
	;
	goto L462
L464:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v1333)+92)) = int32(1)
	goto L463
L466:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v274)+48)) = base.I32_clz(v4612) ^ int32(31) + int32(1)
	v4668 = v289
	goto L60
L467:
	;
	v5001 = *(*int32)(unsafe.Add(mBase, uint32(v274)+60))
	if v5001 == int32(0) {
		v6216 = v4668
		goto L521
	} else {
		goto L522
	}
L468:
	;
	v4690 = *(*int32)(unsafe.Add(mBase, uint32(v274)+32))
	v4691 = int32(1)
	goto L471
L469:
	;
	v4843 = int32(2)
	goto L497
L470:
	;
	goto L469
L471:
	;
	v4704 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v4705 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v4705+v4691 < int32(32) {
		goto L474
	} else {
		goto L475
	}
L472:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v4806 + v4804
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v4803<<(uint(v4806)%32) | v4805
	goto L470
L473:
	;
	v4723 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v4724 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v4728 = v4720
	v4729 = v4721
	v4730 = v4724
	v4731 = v4723
	goto L477
L474:
	;
	if v4705 < int32(16) {
		v4803 = v4691
		v4804 = v4691
		v4805 = v4704
		v4806 = v4705
		goto L472
	} else {
		goto L476
	}
L475:
	;
	v4709 = int32(32)
	v4710 = v4709 - v4705
	v4718 = int32(base.Ui32(v4691) >> (uint(v4710) % 32))
	v4719 = v4691 - v4710
	v4720 = v4691<<(uint(v4705)%32) | v4704
	v4721 = v4709
	goto L473
L476:
	;
	v4718 = v4691
	v4719 = v4691
	v4720 = v4704
	v4721 = v4705
	goto L473
L477:
	;
	if base.Ui32(v4730+int32(2)) <= base.Ui32(v4731) {
		v4786 = v4730
		v4787 = v4731
		goto L479
	} else {
		goto L480
	}
L478:
	;
	v4803 = v4718
	v4804 = v4719
	v4805 = v4799
	v4806 = v4797
	goto L472
L479:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v4786))) = uint16(v4728)
	v4794 = v4786 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v4794
	v4797 = v4729 + int32(-16)
	v4799 = int32(base.Ui32(v4728) >> (uint(int32(16)) % 32))
	if int32(31) < v4729 {
		v4728 = v4799
		v4729 = v4797
		v4730 = v4794
		v4731 = v4787
		goto L477
	} else {
		goto L494
	}
L480:
	;
	v4740 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v4741 = v4731 - v4740
	v4744 = base.I64_extend_i32_s(v4741) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v4744) {
		v4768 = v4740
		goto L482
	} else {
		goto L483
	}
L481:
	;
	if v4730 == v4740 {
		goto L492
	} else {
		goto L493
	}
L482:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v4768
	goto L469
L483:
	;
	v4747 = v4730 - v4740
	v4749 = v4744 + base.I64_extend_i32_u(v4747)
	if base.Ui64(int64(4294967295)) < base.Ui64(v4749) {
		v4768 = v4740
		goto L482
	} else {
		goto L484
	}
L484:
	;
	v4752 = base.I32_wrap_i64(v4749)
	if v4731 == v4740 {
		goto L485
	} else {
		goto L486
	}
L485:
	;
	v4759 = int32(base.Ui32(v4741*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v4752) < base.Ui32(v4759) {
		goto L488
	} else {
		goto L489
	}
L486:
	;
	if base.Ui32(v4752) <= base.Ui32(v4741) {
		v4786 = v4730
		v4787 = v4731
		goto L479
	} else {
		goto L487
	}
L487:
	;
	goto L485
L488:
	;
	v4761 = v4759
	goto L490
L489:
	;
	v4761 = v4752
	goto L490
L490:
	;
	v4765 = v4761&int32(-1024) + int32(1024)
	v4766 = F_WebPSafeMalloc(m, int64(1), v4765)
	mBase = m.M
	if v4766 != 0 {
		goto L481
	} else {
		goto L491
	}
L491:
	;
	v4767 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v4768 = v4767
	goto L482
L492:
	;
	v4779 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v4779)
	mBase = m.M
	v4781 = v4766 + v4765
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v4781
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v4766
	v4786 = v4766 + v4747
	v4787 = v4781
	goto L479
L493:
	;
	v4777 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v4778 = F_memcpy(m, v4766, v4777, v4747)
	mBase = m.M
	goto L492
L494:
	;
	goto L478
L495:
	;
	v4995 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v4997 = m.G72
	v4998 = *(*int32)(unsafe.Add(mBase, uint32(v4997)))
	m.T0[v4998].(func(*base.Module, int32, int32))(m, v4995, v4690*v276)
	mBase = m.M
	goto L467
L496:
	;
	goto L495
L497:
	;
	v4856 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v4857 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v4857+v4843 < int32(32) {
		goto L500
	} else {
		goto L501
	}
L498:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v4958 + v4956
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v4955<<(uint(v4958)%32) | v4957
	goto L496
L499:
	;
	v4875 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v4876 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v4880 = v4872
	v4881 = v4873
	v4882 = v4876
	v4883 = v4875
	goto L503
L500:
	;
	if v4857 < int32(16) {
		v4955 = v4843
		v4956 = v4843
		v4957 = v4856
		v4958 = v4857
		goto L498
	} else {
		goto L502
	}
L501:
	;
	v4861 = int32(32)
	v4862 = v4861 - v4857
	v4870 = int32(base.Ui32(v4843) >> (uint(v4862) % 32))
	v4871 = v4843 - v4862
	v4872 = v4843<<(uint(v4857)%32) | v4856
	v4873 = v4861
	goto L499
L502:
	;
	v4870 = v4843
	v4871 = v4843
	v4872 = v4856
	v4873 = v4857
	goto L499
L503:
	;
	if base.Ui32(v4882+int32(2)) <= base.Ui32(v4883) {
		v4938 = v4882
		v4939 = v4883
		goto L505
	} else {
		goto L506
	}
L504:
	;
	v4955 = v4870
	v4956 = v4871
	v4957 = v4951
	v4958 = v4949
	goto L498
L505:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v4938))) = uint16(v4880)
	v4946 = v4938 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v4946
	v4949 = v4881 + int32(-16)
	v4951 = int32(base.Ui32(v4880) >> (uint(int32(16)) % 32))
	if int32(31) < v4881 {
		v4880 = v4951
		v4881 = v4949
		v4882 = v4946
		v4883 = v4939
		goto L503
	} else {
		goto L520
	}
L506:
	;
	v4892 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v4893 = v4883 - v4892
	v4896 = base.I64_extend_i32_s(v4893) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v4896) {
		v4920 = v4892
		goto L508
	} else {
		goto L509
	}
L507:
	;
	if v4882 == v4892 {
		goto L518
	} else {
		goto L519
	}
L508:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v4920
	goto L495
L509:
	;
	v4899 = v4882 - v4892
	v4901 = v4896 + base.I64_extend_i32_u(v4899)
	if base.Ui64(int64(4294967295)) < base.Ui64(v4901) {
		v4920 = v4892
		goto L508
	} else {
		goto L510
	}
L510:
	;
	v4904 = base.I32_wrap_i64(v4901)
	if v4883 == v4892 {
		goto L511
	} else {
		goto L512
	}
L511:
	;
	v4911 = int32(base.Ui32(v4893*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v4904) < base.Ui32(v4911) {
		goto L514
	} else {
		goto L515
	}
L512:
	;
	if base.Ui32(v4904) <= base.Ui32(v4893) {
		v4938 = v4882
		v4939 = v4883
		goto L505
	} else {
		goto L513
	}
L513:
	;
	goto L511
L514:
	;
	v4913 = v4911
	goto L516
L515:
	;
	v4913 = v4904
	goto L516
L516:
	;
	v4917 = v4913&int32(-1024) + int32(1024)
	v4918 = F_WebPSafeMalloc(m, int64(1), v4917)
	mBase = m.M
	if v4918 != 0 {
		goto L507
	} else {
		goto L517
	}
L517:
	;
	v4919 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v4920 = v4919
	goto L508
L518:
	;
	v4931 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v4931)
	mBase = m.M
	v4933 = v4918 + v4917
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v4933
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v4918
	v4938 = v4918 + v4899
	v4939 = v4933
	goto L505
L519:
	;
	v4929 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v4930 = F_memcpy(m, v4918, v4929, v4899)
	mBase = m.M
	goto L518
L520:
	;
	goto L504
L521:
	;
	v6235 = *(*int32)(unsafe.Add(mBase, uint32(v274)+52))
	if v6235 == int32(0) {
		v6733 = v6216
		goto L656
	} else {
		goto L657
	}
L522:
	;
	v5005 = base.I32_div_s(v4668, int32(3))
	v5007 = *(*int32)(unsafe.Add(mBase, uint32(v274)+64))
	if v5007 != 0 {
		v5010 = int32(100)
		goto L523
	} else {
		goto L524
	}
L523:
	;
	v5011 = *(*int32)(unsafe.Add(mBase, uint32(v274)+56))
	v5014 = *(*int32)(unsafe.Add(mBase, uint32(v274)+40))
	v5015 = int32(9)
	if base.Ui32(v5014) < base.Ui32(v5015) {
		goto L525
	} else {
		goto L526
	}
L524:
	;
	v5008 = *(*int32)(unsafe.Add(mBase, uint32(v274)))
	v5009 = *(*int32)(unsafe.Add(mBase, uint32(v5008)+92))
	v5010 = v5009
	goto L523
L525:
	;
	v5018 = v5014
	goto L527
L526:
	;
	v5018 = v5015
	goto L527
L527:
	;
	if v5014 < int32(2) {
		goto L528
	} else {
		goto L529
	}
L528:
	;
	v5021 = int32(2)
	goto L530
L529:
	;
	v5021 = v5018
	goto L530
L530:
	;
	v5022 = int32(1) << (uint(v5021) % 32)
	v5023 = *(*int32)(unsafe.Add(mBase, uint32(v274)+32))
	v5025 = v5023 + int32(-1)
	v5030 = int32(base.Ui32(v5022+v5025)>>(uint(v5021)%32)) * int32(base.Ui32(v5022+v292)>>(uint(v5021)%32))
	if base.Ui32(int32(8)) < base.Ui32(v5021) {
		v5118 = v5021
		v5131 = v5030
		goto L531
	} else {
		goto L532
	}
L531:
	;
	if base.Ui32(int32(3)) <= base.Ui32(v5118) {
		goto L539
	} else {
		goto L540
	}
L532:
	;
	if v5030 < int32(_a_F_EncodeStreamHook_4) {
		v5118 = v5021
		v5131 = v5030
		goto L531
	} else {
		goto L533
	}
L533:
	;
	v5040 = v5021
	goto L534
L534:
	;
	v5102 = int32(2) << (uint(v5040) % 32)
	v5105 = v5040 + int32(1)
	v5109 = int32(base.Ui32(v5102+v5025)>>(uint(v5105)%32)) * int32(base.Ui32(v5102+v292)>>(uint(v5105)%32))
	if base.Ui32(int32(7)) < base.Ui32(v5040) {
		v5118 = v5105
		v5131 = v5109
		goto L531
	} else {
		goto L536
	}
L535:
	;
	v5118 = v5105
	v5131 = v5109
	goto L531
L536:
	;
	if int32(_a_F_EncodeStreamHook_5) < v5109 {
		v5040 = v5105
		goto L534
	} else {
		goto L537
	}
L537:
	;
	goto L535
L538:
	;
	v5343 = int32(1)
	v5345 = *(*int32)(unsafe.Add(mBase, uint32(v274)))
	v5346 = *(*int32)(unsafe.Add(mBase, uint32(v5345)+8))
	if int32(4) < v5346 {
		goto L549
	} else {
		goto L550
	}
L539:
	;
	if v5131 == int32(1) {
		goto L541
	} else {
		goto L542
	}
L540:
	;
	v5317 = v5118
	goto L538
L541:
	;
	v5184 = int32(1)
	v5186 = v5118 + int32(-1)
	v5187 = v5184 << (uint(v5186) % 32)
	if int32(base.Ui32(v5187+v5025)>>(uint(v5186)%32))*int32(base.Ui32(v5187+v292)>>(uint(v5186)%32)) != v5184 {
		v5317 = v5118
		goto L538
	} else {
		goto L543
	}
L542:
	;
	v5317 = v5118
	goto L538
L543:
	;
	v5199 = v5118
	goto L544
L544:
	;
	v5262 = v5199 + int32(-1)
	if int32(3) <= v5262 {
		goto L546
	} else {
		goto L547
	}
L545:
	;
	v5317 = v5262
	goto L538
L546:
	;
	v5267 = v5199 + int32(-2)
	v5268 = int32(1)
	v5269 = v5268 << (uint(v5267) % 32)
	if int32(base.Ui32(v5269+v5025)>>(uint(v5267)%32))*int32(base.Ui32(v5269+v292)>>(uint(v5267)%32)) == v5268 {
		v5199 = v5262
		goto L544
	} else {
		goto L548
	}
L547:
	;
	v5317 = int32(2)
	goto L538
L548:
	;
	goto L545
L549:
	;
	v5354 = v5346<<(uint(v5343)%32) + int32(-8)
	goto L551
L550:
	;
	v5354 = int32(0)
	goto L551
L551:
	;
	v5355 = v5317 - v5354
	v5356 = int32(9)
	if base.Ui32(v5355) < base.Ui32(v5356) {
		goto L552
	} else {
		goto L553
	}
L552:
	;
	v5359 = v5355
	goto L554
L553:
	;
	v5359 = v5356
	goto L554
L554:
	;
	if v5355 < int32(2) {
		goto L555
	} else {
		goto L556
	}
L555:
	;
	v5362 = int32(2)
	goto L557
L556:
	;
	v5362 = v5359
	goto L557
L557:
	;
	v5363 = v5343 << (uint(v5362) % 32)
	v5368 = int32(base.Ui32(v5363+v5025)>>(uint(v5362)%32)) * int32(base.Ui32(v5363+v292)>>(uint(v5362)%32))
	if base.Ui32(int32(8)) < base.Ui32(v5362) {
		v5456 = v5362
		v5469 = v5368
		goto L558
	} else {
		goto L559
	}
L558:
	;
	if base.Ui32(int32(3)) <= base.Ui32(v5456) {
		goto L566
	} else {
		goto L567
	}
L559:
	;
	if v5368 < int32(_a_F_EncodeStreamHook_4) {
		v5456 = v5362
		v5469 = v5368
		goto L558
	} else {
		goto L560
	}
L560:
	;
	v5378 = v5362
	goto L561
L561:
	;
	v5440 = int32(2) << (uint(v5378) % 32)
	v5443 = v5378 + int32(1)
	v5447 = int32(base.Ui32(v5440+v5025)>>(uint(v5443)%32)) * int32(base.Ui32(v5440+v292)>>(uint(v5443)%32))
	if base.Ui32(int32(7)) < base.Ui32(v5378) {
		v5456 = v5443
		v5469 = v5447
		goto L558
	} else {
		goto L563
	}
L562:
	;
	v5456 = v5443
	v5469 = v5447
	goto L558
L563:
	;
	if int32(_a_F_EncodeStreamHook_5) < v5447 {
		v5378 = v5443
		goto L561
	} else {
		goto L564
	}
L564:
	;
	goto L562
L565:
	;
	v5681 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v5682 = *(*int32)(unsafe.Add(mBase, uint32(v274)+16))
	v5683 = *(*int32)(unsafe.Add(mBase, uint32(v274)+20))
	v5684 = *(*int32)(unsafe.Add(mBase, uint32(v5345)+96))
	v5685 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v5687 = base.I32_div_s(v4668, int32(6))
	v5692 = F_VP8LResidualImage(m, v5023, v276, v5632, v5317, v286, v5681, v5682, v5683, v5010, v5684, v5011, v5685, v5687, v267+int32(60), v267+int32(4))
	mBase = m.M
	if v5692 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L576
	}
L566:
	;
	if v5469 == int32(1) {
		goto L568
	} else {
		goto L569
	}
L567:
	;
	v5632 = v5456
	goto L565
L568:
	;
	v5522 = int32(1)
	v5524 = v5456 + int32(-1)
	v5525 = v5522 << (uint(v5524) % 32)
	if int32(base.Ui32(v5525+v5025)>>(uint(v5524)%32))*int32(base.Ui32(v5525+v292)>>(uint(v5524)%32)) != v5522 {
		v5632 = v5456
		goto L565
	} else {
		goto L570
	}
L569:
	;
	v5632 = v5456
	goto L565
L570:
	;
	v5537 = v5456
	goto L571
L571:
	;
	v5600 = v5537 + int32(-1)
	if int32(3) <= v5600 {
		goto L573
	} else {
		goto L574
	}
L572:
	;
	v5632 = v5600
	goto L565
L573:
	;
	v5605 = v5537 + int32(-2)
	v5606 = int32(1)
	v5607 = v5606 << (uint(v5605) % 32)
	if int32(base.Ui32(v5607+v5025)>>(uint(v5605)%32))*int32(base.Ui32(v5607+v292)>>(uint(v5605)%32)) == v5606 {
		v5537 = v5600
		goto L571
	} else {
		goto L575
	}
L574:
	;
	v5632 = int32(2)
	goto L565
L575:
	;
	goto L572
L576:
	;
	v5695 = int32(1)
	goto L579
L577:
	;
	v5847 = int32(0)
	v5848 = int32(2)
	goto L605
L578:
	;
	goto L577
L579:
	;
	v5708 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v5709 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v5709+v5695 < int32(32) {
		goto L582
	} else {
		goto L583
	}
L580:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v5810 + v5808
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v5807<<(uint(v5810)%32) | v5809
	goto L578
L581:
	;
	v5727 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v5728 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v5732 = v5724
	v5733 = v5725
	v5734 = v5728
	v5735 = v5727
	goto L585
L582:
	;
	if v5709 < int32(16) {
		v5807 = v5695
		v5808 = v5695
		v5809 = v5708
		v5810 = v5709
		goto L580
	} else {
		goto L584
	}
L583:
	;
	v5713 = int32(32)
	v5714 = v5713 - v5709
	v5722 = int32(base.Ui32(v5695) >> (uint(v5714) % 32))
	v5723 = v5695 - v5714
	v5724 = v5695<<(uint(v5709)%32) | v5708
	v5725 = v5713
	goto L581
L584:
	;
	v5722 = v5695
	v5723 = v5695
	v5724 = v5708
	v5725 = v5709
	goto L581
L585:
	;
	if base.Ui32(v5734+int32(2)) <= base.Ui32(v5735) {
		v5790 = v5734
		v5791 = v5735
		goto L587
	} else {
		goto L588
	}
L586:
	;
	v5807 = v5722
	v5808 = v5723
	v5809 = v5803
	v5810 = v5801
	goto L580
L587:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v5790))) = uint16(v5732)
	v5798 = v5790 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v5798
	v5801 = v5733 + int32(-16)
	v5803 = int32(base.Ui32(v5732) >> (uint(int32(16)) % 32))
	if int32(31) < v5733 {
		v5732 = v5803
		v5733 = v5801
		v5734 = v5798
		v5735 = v5791
		goto L585
	} else {
		goto L602
	}
L588:
	;
	v5744 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v5745 = v5735 - v5744
	v5748 = base.I64_extend_i32_s(v5745) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5748) {
		v5772 = v5744
		goto L590
	} else {
		goto L591
	}
L589:
	;
	if v5734 == v5744 {
		goto L600
	} else {
		goto L601
	}
L590:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v5772
	goto L577
L591:
	;
	v5751 = v5734 - v5744
	v5753 = v5748 + base.I64_extend_i32_u(v5751)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5753) {
		v5772 = v5744
		goto L590
	} else {
		goto L592
	}
L592:
	;
	v5756 = base.I32_wrap_i64(v5753)
	if v5735 == v5744 {
		goto L593
	} else {
		goto L594
	}
L593:
	;
	v5763 = int32(base.Ui32(v5745*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v5756) < base.Ui32(v5763) {
		goto L596
	} else {
		goto L597
	}
L594:
	;
	if base.Ui32(v5756) <= base.Ui32(v5745) {
		v5790 = v5734
		v5791 = v5735
		goto L587
	} else {
		goto L595
	}
L595:
	;
	goto L593
L596:
	;
	v5765 = v5763
	goto L598
L597:
	;
	v5765 = v5756
	goto L598
L598:
	;
	v5769 = v5765&int32(-1024) + int32(1024)
	v5770 = F_WebPSafeMalloc(m, int64(1), v5769)
	mBase = m.M
	if v5770 != 0 {
		goto L589
	} else {
		goto L599
	}
L599:
	;
	v5771 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v5772 = v5771
	goto L590
L600:
	;
	v5783 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v5783)
	mBase = m.M
	v5785 = v5770 + v5769
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v5785
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v5770
	v5790 = v5770 + v5751
	v5791 = v5785
	goto L587
L601:
	;
	v5781 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v5782 = F_memcpy(m, v5770, v5781, v5751)
	mBase = m.M
	goto L600
L602:
	;
	goto L586
L603:
	;
	v5999 = *(*int32)(unsafe.Add(mBase, uint32(v267)+4))
	v6001 = v5999 + int32(-2)
	v6002 = int32(3)
	goto L631
L604:
	;
	goto L603
L605:
	;
	v5860 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v5861 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v5861+v5848 < int32(32) {
		goto L608
	} else {
		goto L609
	}
L606:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v5962 + v5960
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v5959<<(uint(v5962)%32) | v5961
	goto L604
L607:
	;
	v5879 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v5880 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v5884 = v5876
	v5885 = v5877
	v5886 = v5880
	v5887 = v5879
	goto L611
L608:
	;
	if v5861 < int32(16) {
		v5959 = v5847
		v5960 = v5848
		v5961 = v5860
		v5962 = v5861
		goto L606
	} else {
		goto L610
	}
L609:
	;
	v5865 = int32(32)
	v5866 = v5865 - v5861
	v5874 = int32(base.Ui32(v5847) >> (uint(v5866) % 32))
	v5875 = v5848 - v5866
	v5876 = v5847<<(uint(v5861)%32) | v5860
	v5877 = v5865
	goto L607
L610:
	;
	v5874 = v5847
	v5875 = v5848
	v5876 = v5860
	v5877 = v5861
	goto L607
L611:
	;
	if base.Ui32(v5886+int32(2)) <= base.Ui32(v5887) {
		v5942 = v5886
		v5943 = v5887
		goto L613
	} else {
		goto L614
	}
L612:
	;
	v5959 = v5874
	v5960 = v5875
	v5961 = v5955
	v5962 = v5953
	goto L606
L613:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v5942))) = uint16(v5884)
	v5950 = v5942 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v5950
	v5953 = v5885 + int32(-16)
	v5955 = int32(base.Ui32(v5884) >> (uint(int32(16)) % 32))
	if int32(31) < v5885 {
		v5884 = v5955
		v5885 = v5953
		v5886 = v5950
		v5887 = v5943
		goto L611
	} else {
		goto L628
	}
L614:
	;
	v5896 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v5897 = v5887 - v5896
	v5900 = base.I64_extend_i32_s(v5897) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5900) {
		v5924 = v5896
		goto L616
	} else {
		goto L617
	}
L615:
	;
	if v5886 == v5896 {
		goto L626
	} else {
		goto L627
	}
L616:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v5924
	goto L603
L617:
	;
	v5903 = v5886 - v5896
	v5905 = v5900 + base.I64_extend_i32_u(v5903)
	if base.Ui64(int64(4294967295)) < base.Ui64(v5905) {
		v5924 = v5896
		goto L616
	} else {
		goto L618
	}
L618:
	;
	v5908 = base.I32_wrap_i64(v5905)
	if v5887 == v5896 {
		goto L619
	} else {
		goto L620
	}
L619:
	;
	v5915 = int32(base.Ui32(v5897*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v5908) < base.Ui32(v5915) {
		goto L622
	} else {
		goto L623
	}
L620:
	;
	if base.Ui32(v5908) <= base.Ui32(v5897) {
		v5942 = v5886
		v5943 = v5887
		goto L613
	} else {
		goto L621
	}
L621:
	;
	goto L619
L622:
	;
	v5917 = v5915
	goto L624
L623:
	;
	v5917 = v5908
	goto L624
L624:
	;
	v5921 = v5917&int32(-1024) + int32(1024)
	v5922 = F_WebPSafeMalloc(m, int64(1), v5921)
	mBase = m.M
	if v5922 != 0 {
		goto L615
	} else {
		goto L625
	}
L625:
	;
	v5923 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v5924 = v5923
	goto L616
L626:
	;
	v5935 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v5935)
	mBase = m.M
	v5937 = v5922 + v5921
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v5937
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v5922
	v5942 = v5922 + v5903
	v5943 = v5937
	goto L613
L627:
	;
	v5933 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v5934 = F_memcpy(m, v5922, v5933, v5903)
	mBase = m.M
	goto L626
L628:
	;
	goto L612
L629:
	;
	v6153 = *(*int32)(unsafe.Add(mBase, uint32(v274)+20))
	v6155 = *(*int32)(unsafe.Add(mBase, uint32(v267)+4))
	v6156 = int32(1) << (uint(v6155) % 32)
	v6161 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v6165 = F_EncodeImageNoHuffman(m, v268, v6153, v293, v297, int32(base.Ui32(v6156+v5025)>>(uint(v6155)%32)), int32(base.Ui32(v6156+v292)>>(uint(v6155)%32)), v283, v286, v6161, v5005-v5687, v267+int32(60))
	mBase = m.M
	if v6165 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L655
	}
L630:
	;
	goto L629
L631:
	;
	v6014 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v6015 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v6015+v6002 < int32(32) {
		goto L634
	} else {
		goto L635
	}
L632:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v6116 + v6114
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v6113<<(uint(v6116)%32) | v6115
	goto L630
L633:
	;
	v6033 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v6034 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v6038 = v6030
	v6039 = v6031
	v6040 = v6034
	v6041 = v6033
	goto L637
L634:
	;
	if v6015 < int32(16) {
		v6113 = v6001
		v6114 = v6002
		v6115 = v6014
		v6116 = v6015
		goto L632
	} else {
		goto L636
	}
L635:
	;
	v6019 = int32(32)
	v6020 = v6019 - v6015
	v6028 = int32(base.Ui32(v6001) >> (uint(v6020) % 32))
	v6029 = v6002 - v6020
	v6030 = v6001<<(uint(v6015)%32) | v6014
	v6031 = v6019
	goto L633
L636:
	;
	v6028 = v6001
	v6029 = v6002
	v6030 = v6014
	v6031 = v6015
	goto L633
L637:
	;
	if base.Ui32(v6040+int32(2)) <= base.Ui32(v6041) {
		v6096 = v6040
		v6097 = v6041
		goto L639
	} else {
		goto L640
	}
L638:
	;
	v6113 = v6028
	v6114 = v6029
	v6115 = v6109
	v6116 = v6107
	goto L632
L639:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6096))) = uint16(v6038)
	v6104 = v6096 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6104
	v6107 = v6039 + int32(-16)
	v6109 = int32(base.Ui32(v6038) >> (uint(int32(16)) % 32))
	if int32(31) < v6039 {
		v6038 = v6109
		v6039 = v6107
		v6040 = v6104
		v6041 = v6097
		goto L637
	} else {
		goto L654
	}
L640:
	;
	v6050 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6051 = v6041 - v6050
	v6054 = base.I64_extend_i32_s(v6051) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6054) {
		v6078 = v6050
		goto L642
	} else {
		goto L643
	}
L641:
	;
	if v6040 == v6050 {
		goto L652
	} else {
		goto L653
	}
L642:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6078
	goto L629
L643:
	;
	v6057 = v6040 - v6050
	v6059 = v6054 + base.I64_extend_i32_u(v6057)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6059) {
		v6078 = v6050
		goto L642
	} else {
		goto L644
	}
L644:
	;
	v6062 = base.I32_wrap_i64(v6059)
	if v6041 == v6050 {
		goto L645
	} else {
		goto L646
	}
L645:
	;
	v6069 = int32(base.Ui32(v6051*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6062) < base.Ui32(v6069) {
		goto L648
	} else {
		goto L649
	}
L646:
	;
	if base.Ui32(v6062) <= base.Ui32(v6051) {
		v6096 = v6040
		v6097 = v6041
		goto L639
	} else {
		goto L647
	}
L647:
	;
	goto L645
L648:
	;
	v6071 = v6069
	goto L650
L649:
	;
	v6071 = v6062
	goto L650
L650:
	;
	v6075 = v6071&int32(-1024) + int32(1024)
	v6076 = F_WebPSafeMalloc(m, int64(1), v6075)
	mBase = m.M
	if v6076 != 0 {
		goto L641
	} else {
		goto L651
	}
L651:
	;
	v6077 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6078 = v6077
	goto L642
L652:
	;
	v6089 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v6089)
	mBase = m.M
	v6091 = v6076 + v6075
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v6091
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v6076
	v6096 = v6076 + v6057
	v6097 = v6091
	goto L639
L653:
	;
	v6087 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6088 = F_memcpy(m, v6076, v6087, v6057)
	mBase = m.M
	goto L652
L654:
	;
	goto L638
L655:
	;
	v6216 = v4668 - v5005
	goto L521
L656:
	;
	v6734 = int32(0)
	v6735 = int32(1)
	goto L740
L657:
	;
	v6239 = base.I32_div_s(v6216, int32(2))
	v6240 = *(*int32)(unsafe.Add(mBase, uint32(v274)+32))
	v6241 = *(*int32)(unsafe.Add(mBase, uint32(v274)+44))
	v6242 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v6243 = *(*int32)(unsafe.Add(mBase, uint32(v274)+20))
	v6244 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v6246 = base.I32_div_s(v6216, int32(4))
	v6249 = F_VP8LColorSpaceTransform(m, v6240, v276, v6241, v283, v6242, v6243, v6244, v6246, v267+int32(60), v267)
	mBase = m.M
	if v6249 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L658
	}
L658:
	;
	v6252 = int32(1)
	goto L661
L659:
	;
	v6404 = int32(1)
	v6405 = int32(2)
	goto L687
L660:
	;
	goto L659
L661:
	;
	v6265 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v6266 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v6266+v6252 < int32(32) {
		goto L664
	} else {
		goto L665
	}
L662:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v6367 + v6365
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v6364<<(uint(v6367)%32) | v6366
	goto L660
L663:
	;
	v6284 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v6285 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v6289 = v6281
	v6290 = v6282
	v6291 = v6285
	v6292 = v6284
	goto L667
L664:
	;
	if v6266 < int32(16) {
		v6364 = v6252
		v6365 = v6252
		v6366 = v6265
		v6367 = v6266
		goto L662
	} else {
		goto L666
	}
L665:
	;
	v6270 = int32(32)
	v6271 = v6270 - v6266
	v6279 = int32(base.Ui32(v6252) >> (uint(v6271) % 32))
	v6280 = v6252 - v6271
	v6281 = v6252<<(uint(v6266)%32) | v6265
	v6282 = v6270
	goto L663
L666:
	;
	v6279 = v6252
	v6280 = v6252
	v6281 = v6265
	v6282 = v6266
	goto L663
L667:
	;
	if base.Ui32(v6291+int32(2)) <= base.Ui32(v6292) {
		v6347 = v6291
		v6348 = v6292
		goto L669
	} else {
		goto L670
	}
L668:
	;
	v6364 = v6279
	v6365 = v6280
	v6366 = v6360
	v6367 = v6358
	goto L662
L669:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6347))) = uint16(v6289)
	v6355 = v6347 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6355
	v6358 = v6290 + int32(-16)
	v6360 = int32(base.Ui32(v6289) >> (uint(int32(16)) % 32))
	if int32(31) < v6290 {
		v6289 = v6360
		v6290 = v6358
		v6291 = v6355
		v6292 = v6348
		goto L667
	} else {
		goto L684
	}
L670:
	;
	v6301 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6302 = v6292 - v6301
	v6305 = base.I64_extend_i32_s(v6302) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6305) {
		v6329 = v6301
		goto L672
	} else {
		goto L673
	}
L671:
	;
	if v6291 == v6301 {
		goto L682
	} else {
		goto L683
	}
L672:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6329
	goto L659
L673:
	;
	v6308 = v6291 - v6301
	v6310 = v6305 + base.I64_extend_i32_u(v6308)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6310) {
		v6329 = v6301
		goto L672
	} else {
		goto L674
	}
L674:
	;
	v6313 = base.I32_wrap_i64(v6310)
	if v6292 == v6301 {
		goto L675
	} else {
		goto L676
	}
L675:
	;
	v6320 = int32(base.Ui32(v6302*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6313) < base.Ui32(v6320) {
		goto L678
	} else {
		goto L679
	}
L676:
	;
	if base.Ui32(v6313) <= base.Ui32(v6302) {
		v6347 = v6291
		v6348 = v6292
		goto L669
	} else {
		goto L677
	}
L677:
	;
	goto L675
L678:
	;
	v6322 = v6320
	goto L680
L679:
	;
	v6322 = v6313
	goto L680
L680:
	;
	v6326 = v6322&int32(-1024) + int32(1024)
	v6327 = F_WebPSafeMalloc(m, int64(1), v6326)
	mBase = m.M
	if v6327 != 0 {
		goto L671
	} else {
		goto L681
	}
L681:
	;
	v6328 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6329 = v6328
	goto L672
L682:
	;
	v6340 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v6340)
	mBase = m.M
	v6342 = v6327 + v6326
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v6342
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v6327
	v6347 = v6327 + v6308
	v6348 = v6342
	goto L669
L683:
	;
	v6338 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6339 = F_memcpy(m, v6327, v6338, v6308)
	mBase = m.M
	goto L682
L684:
	;
	goto L668
L685:
	;
	v6556 = *(*int32)(unsafe.Add(mBase, uint32(v267)))
	v6558 = v6556 + int32(-2)
	v6559 = int32(3)
	goto L713
L686:
	;
	goto L685
L687:
	;
	v6417 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v6418 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v6418+v6405 < int32(32) {
		goto L690
	} else {
		goto L691
	}
L688:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v6519 + v6517
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v6516<<(uint(v6519)%32) | v6518
	goto L686
L689:
	;
	v6436 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v6437 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v6441 = v6433
	v6442 = v6434
	v6443 = v6437
	v6444 = v6436
	goto L693
L690:
	;
	if v6418 < int32(16) {
		v6516 = v6404
		v6517 = v6405
		v6518 = v6417
		v6519 = v6418
		goto L688
	} else {
		goto L692
	}
L691:
	;
	v6422 = int32(32)
	v6423 = v6422 - v6418
	v6431 = int32(base.Ui32(v6404) >> (uint(v6423) % 32))
	v6432 = v6405 - v6423
	v6433 = v6404<<(uint(v6418)%32) | v6417
	v6434 = v6422
	goto L689
L692:
	;
	v6431 = v6404
	v6432 = v6405
	v6433 = v6417
	v6434 = v6418
	goto L689
L693:
	;
	if base.Ui32(v6443+int32(2)) <= base.Ui32(v6444) {
		v6499 = v6443
		v6500 = v6444
		goto L695
	} else {
		goto L696
	}
L694:
	;
	v6516 = v6431
	v6517 = v6432
	v6518 = v6512
	v6519 = v6510
	goto L688
L695:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6499))) = uint16(v6441)
	v6507 = v6499 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6507
	v6510 = v6442 + int32(-16)
	v6512 = int32(base.Ui32(v6441) >> (uint(int32(16)) % 32))
	if int32(31) < v6442 {
		v6441 = v6512
		v6442 = v6510
		v6443 = v6507
		v6444 = v6500
		goto L693
	} else {
		goto L710
	}
L696:
	;
	v6453 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6454 = v6444 - v6453
	v6457 = base.I64_extend_i32_s(v6454) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6457) {
		v6481 = v6453
		goto L698
	} else {
		goto L699
	}
L697:
	;
	if v6443 == v6453 {
		goto L708
	} else {
		goto L709
	}
L698:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6481
	goto L685
L699:
	;
	v6460 = v6443 - v6453
	v6462 = v6457 + base.I64_extend_i32_u(v6460)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6462) {
		v6481 = v6453
		goto L698
	} else {
		goto L700
	}
L700:
	;
	v6465 = base.I32_wrap_i64(v6462)
	if v6444 == v6453 {
		goto L701
	} else {
		goto L702
	}
L701:
	;
	v6472 = int32(base.Ui32(v6454*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6465) < base.Ui32(v6472) {
		goto L704
	} else {
		goto L705
	}
L702:
	;
	if base.Ui32(v6465) <= base.Ui32(v6454) {
		v6499 = v6443
		v6500 = v6444
		goto L695
	} else {
		goto L703
	}
L703:
	;
	goto L701
L704:
	;
	v6474 = v6472
	goto L706
L705:
	;
	v6474 = v6465
	goto L706
L706:
	;
	v6478 = v6474&int32(-1024) + int32(1024)
	v6479 = F_WebPSafeMalloc(m, int64(1), v6478)
	mBase = m.M
	if v6479 != 0 {
		goto L697
	} else {
		goto L707
	}
L707:
	;
	v6480 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6481 = v6480
	goto L698
L708:
	;
	v6492 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v6492)
	mBase = m.M
	v6494 = v6479 + v6478
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v6494
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v6479
	v6499 = v6479 + v6460
	v6500 = v6494
	goto L695
L709:
	;
	v6490 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6491 = F_memcpy(m, v6479, v6490, v6460)
	mBase = m.M
	goto L708
L710:
	;
	goto L694
L711:
	;
	v6710 = *(*int32)(unsafe.Add(mBase, uint32(v274)+20))
	v6712 = *(*int32)(unsafe.Add(mBase, uint32(v267)))
	v6713 = int32(1) << (uint(v6712) % 32)
	v6720 = *(*int32)(unsafe.Add(mBase, uint32(v274)+4))
	v6724 = F_EncodeImageNoHuffman(m, v268, v6710, v293, v297, int32(base.Ui32(v6240+v6713+int32(-1))>>(uint(v6712)%32)), int32(base.Ui32(v292+v6713)>>(uint(v6712)%32)), v283, v286, v6720, v6239-v6246, v267+int32(60))
	mBase = m.M
	if v6724 == int32(0) {
		v10242 = v265
		v10244 = v267
		goto L4
	} else {
		goto L737
	}
L712:
	;
	goto L711
L713:
	;
	v6571 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v6572 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v6572+v6559 < int32(32) {
		goto L716
	} else {
		goto L717
	}
L714:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v6673 + v6671
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v6670<<(uint(v6673)%32) | v6672
	goto L712
L715:
	;
	v6590 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v6591 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v6595 = v6587
	v6596 = v6588
	v6597 = v6591
	v6598 = v6590
	goto L719
L716:
	;
	if v6572 < int32(16) {
		v6670 = v6558
		v6671 = v6559
		v6672 = v6571
		v6673 = v6572
		goto L714
	} else {
		goto L718
	}
L717:
	;
	v6576 = int32(32)
	v6577 = v6576 - v6572
	v6585 = int32(base.Ui32(v6558) >> (uint(v6577) % 32))
	v6586 = v6559 - v6577
	v6587 = v6558<<(uint(v6572)%32) | v6571
	v6588 = v6576
	goto L715
L718:
	;
	v6585 = v6558
	v6586 = v6559
	v6587 = v6571
	v6588 = v6572
	goto L715
L719:
	;
	if base.Ui32(v6597+int32(2)) <= base.Ui32(v6598) {
		v6653 = v6597
		v6654 = v6598
		goto L721
	} else {
		goto L722
	}
L720:
	;
	v6670 = v6585
	v6671 = v6586
	v6672 = v6666
	v6673 = v6664
	goto L714
L721:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6653))) = uint16(v6595)
	v6661 = v6653 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6661
	v6664 = v6596 + int32(-16)
	v6666 = int32(base.Ui32(v6595) >> (uint(int32(16)) % 32))
	if int32(31) < v6596 {
		v6595 = v6666
		v6596 = v6664
		v6597 = v6661
		v6598 = v6654
		goto L719
	} else {
		goto L736
	}
L722:
	;
	v6607 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6608 = v6598 - v6607
	v6611 = base.I64_extend_i32_s(v6608) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6611) {
		v6635 = v6607
		goto L724
	} else {
		goto L725
	}
L723:
	;
	if v6597 == v6607 {
		goto L734
	} else {
		goto L735
	}
L724:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6635
	goto L711
L725:
	;
	v6614 = v6597 - v6607
	v6616 = v6611 + base.I64_extend_i32_u(v6614)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6616) {
		v6635 = v6607
		goto L724
	} else {
		goto L726
	}
L726:
	;
	v6619 = base.I32_wrap_i64(v6616)
	if v6598 == v6607 {
		goto L727
	} else {
		goto L728
	}
L727:
	;
	v6626 = int32(base.Ui32(v6608*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6619) < base.Ui32(v6626) {
		goto L730
	} else {
		goto L731
	}
L728:
	;
	if base.Ui32(v6619) <= base.Ui32(v6608) {
		v6653 = v6597
		v6654 = v6598
		goto L721
	} else {
		goto L729
	}
L729:
	;
	goto L727
L730:
	;
	v6628 = v6626
	goto L732
L731:
	;
	v6628 = v6619
	goto L732
L732:
	;
	v6632 = v6628&int32(-1024) + int32(1024)
	v6633 = F_WebPSafeMalloc(m, int64(1), v6632)
	mBase = m.M
	if v6633 != 0 {
		goto L723
	} else {
		goto L733
	}
L733:
	;
	v6634 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6635 = v6634
	goto L724
L734:
	;
	v6646 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v6646)
	mBase = m.M
	v6648 = v6633 + v6632
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v6648
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v6633
	v6653 = v6633 + v6614
	v6654 = v6648
	goto L721
L735:
	;
	v6644 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6645 = F_memcpy(m, v6633, v6644, v6614)
	mBase = m.M
	goto L734
L736:
	;
	goto L720
L737:
	;
	v6733 = v6216 - v6239
	goto L656
L738:
	;
	v6886 = *(*int32)(unsafe.Add(mBase, uint32(v274)+32))
	v6887 = *(*int32)(unsafe.Add(mBase, uint32(v274)+36))
	v6888 = *(*int32)(unsafe.Add(mBase, uint32(v274)+8))
	v6889 = *(*int32)(unsafe.Add(mBase, uint32(v267)+60))
	v6890 = int64(57)
	v6891 = int32(16)
	goto L767
L739:
	;
	goto L738
L740:
	;
	v6747 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v6748 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v6748+v6735 < int32(32) {
		goto L743
	} else {
		goto L744
	}
L741:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v6849 + v6847
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v6846<<(uint(v6849)%32) | v6848
	goto L739
L742:
	;
	v6766 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v6767 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v6771 = v6763
	v6772 = v6764
	v6773 = v6767
	v6774 = v6766
	goto L746
L743:
	;
	if v6748 < int32(16) {
		v6846 = v6734
		v6847 = v6735
		v6848 = v6747
		v6849 = v6748
		goto L741
	} else {
		goto L745
	}
L744:
	;
	v6752 = int32(32)
	v6753 = v6752 - v6748
	v6761 = int32(base.Ui32(v6734) >> (uint(v6753) % 32))
	v6762 = v6735 - v6753
	v6763 = v6734<<(uint(v6748)%32) | v6747
	v6764 = v6752
	goto L742
L745:
	;
	v6761 = v6734
	v6762 = v6735
	v6763 = v6747
	v6764 = v6748
	goto L742
L746:
	;
	if base.Ui32(v6773+int32(2)) <= base.Ui32(v6774) {
		v6829 = v6773
		v6830 = v6774
		goto L748
	} else {
		goto L749
	}
L747:
	;
	v6846 = v6761
	v6847 = v6762
	v6848 = v6842
	v6849 = v6840
	goto L741
L748:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v6829))) = uint16(v6771)
	v6837 = v6829 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6837
	v6840 = v6772 + int32(-16)
	v6842 = int32(base.Ui32(v6771) >> (uint(int32(16)) % 32))
	if int32(31) < v6772 {
		v6771 = v6842
		v6772 = v6840
		v6773 = v6837
		v6774 = v6830
		goto L746
	} else {
		goto L763
	}
L749:
	;
	v6783 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6784 = v6774 - v6783
	v6787 = base.I64_extend_i32_s(v6784) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6787) {
		v6811 = v6783
		goto L751
	} else {
		goto L752
	}
L750:
	;
	if v6773 == v6783 {
		goto L761
	} else {
		goto L762
	}
L751:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v6811
	goto L738
L752:
	;
	v6790 = v6773 - v6783
	v6792 = v6787 + base.I64_extend_i32_u(v6790)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6792) {
		v6811 = v6783
		goto L751
	} else {
		goto L753
	}
L753:
	;
	v6795 = base.I32_wrap_i64(v6792)
	if v6774 == v6783 {
		goto L754
	} else {
		goto L755
	}
L754:
	;
	v6802 = int32(base.Ui32(v6784*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v6795) < base.Ui32(v6802) {
		goto L757
	} else {
		goto L758
	}
L755:
	;
	if base.Ui32(v6795) <= base.Ui32(v6784) {
		v6829 = v6773
		v6830 = v6774
		goto L748
	} else {
		goto L756
	}
L756:
	;
	goto L754
L757:
	;
	v6804 = v6802
	goto L759
L758:
	;
	v6804 = v6795
	goto L759
L759:
	;
	v6808 = v6804&int32(-1024) + int32(1024)
	v6809 = F_WebPSafeMalloc(m, int64(1), v6808)
	mBase = m.M
	if v6809 != 0 {
		goto L750
	} else {
		goto L760
	}
L760:
	;
	v6810 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6811 = v6810
	goto L751
L761:
	;
	v6822 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v6822)
	mBase = m.M
	v6824 = v6809 + v6808
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v6824
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v6809
	v6829 = v6809 + v6790
	v6830 = v6824
	goto L748
L762:
	;
	v6820 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v6821 = F_memcpy(m, v6809, v6820, v6790)
	mBase = m.M
	goto L761
L763:
	;
	goto L747
L764:
	;
	v6915 = int32(1) << (uint(v6887) % 32)
	v6917 = v6886 + int32(-1)
	v6922 = int32(base.Ui32(v6915+v6917)>>(uint(v6887)%32)) * int32(base.Ui32(v6915+v292)>>(uint(v6887)%32))
	v6923 = base.I64_extend_i32_u(v6922)
	v6924 = int32(4)
	if v6923 == int64(0) {
		goto L772
	} else {
		goto L773
	}
L765:
	;
	goto L764
L766:
	;
	v6910 = F_malloc(m, base.I32_wrap_i64(v6890)*v6891)
	mBase = m.M
	v6912 = v6910
	goto L765
L767:
	;
	v6898 = base.I64_div_u_s(int64(2147418112), v6890)
	v6899 = int32(0)
	v6900 = base.I64_extend_i32_u(v6891)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6900*v6890) {
		v6912 = v6899
		goto L765
	} else {
		goto L768
	}
L768:
	;
	if base.Ui64(v6898) < base.Ui64(v6900) {
		v6912 = v6899
		goto L765
	} else {
		goto L769
	}
L769:
	;
	goto L766
L770:
	;
	v6953 = *(*int64)(unsafe.Add(mBase, uint32(v268+int32(16))))
	*(*int64)(unsafe.Add(mBase, uint32(v267+int32(2128)))) = v6953
	v6961 = *(*int64)(unsafe.Add(mBase, uint32(v268+int32(8))))
	*(*int64)(unsafe.Add(mBase, uint32(v267+int32(2120)))) = v6961
	v6963 = *(*int64)(unsafe.Add(mBase, uint32(v268)))
	*(*int64)(unsafe.Add(mBase, uint32(v267)+2112)) = v6963
	v6965 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v267)+64)) = v6965
	v6968 = v267 + int32(1088)
	*(*int64)(unsafe.Add(mBase, uint32(v6968))) = v6965
	*(*int64)(unsafe.Add(mBase, uint32(v267+int32(1104)))) = v6965
	*(*int64)(unsafe.Add(mBase, uint32(v267+int32(1096)))) = v6965
	v6984 = int32(1024)
	v6986 = F_WebPSafeMalloc(m, int64(1), v6984)
	mBase = m.M
	if v6986 != 0 {
		goto L782
	} else {
		goto L783
	}
L771:
	;
	goto L770
L772:
	;
	v6943 = F_malloc(m, base.I32_wrap_i64(v6923)*v6924)
	mBase = m.M
	v6945 = v6943
	goto L771
L773:
	;
	v6931 = base.I64_div_u_s(int64(2147418112), v6923)
	v6932 = int32(0)
	v6933 = base.I64_extend_i32_u(v6924)
	if base.Ui64(int64(4294967295)) < base.Ui64(v6933*v6923) {
		v6945 = v6932
		goto L771
	} else {
		goto L774
	}
L774:
	;
	if base.Ui64(v6931) < base.Ui64(v6933) {
		v6945 = v6932
		goto L771
	} else {
		goto L775
	}
L775:
	;
	goto L772
L776:
	;
	F_free(m, v9984)
	mBase = m.M
	goto L1134
L777:
	;
	v9945 = v9878
	v9947 = v9880
	v9948 = v9881
	v9949 = int32(0)
	v9950 = v9883
	v9951 = v9884
	v9952 = v9885
	v9953 = v9886
	v9954 = v9887
	v9955 = v9888
	v9956 = v9889
	v9957 = v9890
	v9958 = v9891
	v9960 = v9893
	v9963 = v9896
	v9964 = v9897
	v9965 = v9898
	v9966 = v9899
	v9967 = v9900
	v9968 = v9901
	v9969 = v9902
	v9970 = v9903
	v9971 = v9904
	v9972 = v9905
	v9973 = v9906
	v9974 = v9907
	v9975 = v9908
	v9976 = v9909
	v9977 = v9910
	v9978 = v9911
	v9979 = v9912
	v9980 = v9913
	v9981 = v9914
	v9982 = v9915
	v9984 = v9917
	v9986 = v9919
	v9988 = v9921
	v9989 = v9922
	goto L776
L778:
	;
	v9878 = v9811
	v9880 = v9813
	v9881 = v9814
	v9883 = int32(0)
	v9884 = v9817
	v9885 = v9818
	v9886 = v9819
	v9887 = v9820
	v9888 = v9821
	v9889 = v9822
	v9890 = v9823
	v9891 = v9824
	v9893 = v9826
	v9896 = v9829
	v9897 = v9830
	v9898 = v9831
	v9899 = v9832
	v9900 = v9833
	v9901 = v9834
	v9902 = v9835
	v9903 = v9836
	v9904 = v9837
	v9905 = v9838
	v9906 = v9839
	v9907 = v9840
	v9908 = v9841
	v9909 = v9842
	v9910 = v9843
	v9911 = v9844
	v9912 = v9845
	v9913 = v9846
	v9914 = v9847
	v9915 = v9848
	v9917 = v9850
	v9919 = v9852
	v9921 = v9854
	v9922 = v9855
	goto L777
L779:
	;
	v9809 = int32(0)
	v9811 = v265
	v9813 = v267
	v9814 = v268
	v9817 = v271
	v9818 = v272
	v9819 = v273
	v9820 = v274
	v9821 = v275
	v9822 = v276
	v9823 = v277
	v9824 = v278
	v9826 = v280
	v9829 = v283
	v9830 = v284
	v9831 = v285
	v9832 = v286
	v9833 = v287
	v9834 = v288
	v9835 = v289
	v9836 = v290
	v9837 = v291
	v9838 = v292
	v9839 = v293
	v9840 = v294
	v9841 = v295
	v9842 = v296
	v9843 = v297
	v9844 = v298
	v9845 = v299
	v9846 = v300
	v9847 = v9779
	v9848 = v9780
	v9850 = v9809
	v9852 = v6912
	v9854 = v9809
	v9855 = v6945
	goto L778
L780:
	;
	if v6912 == int32(0) {
		goto L789
	} else {
		goto L790
	}
L781:
	;
	if v6986 != 0 {
		goto L780
	} else {
		goto L784
	}
L782:
	;
	v6990 = *(*int32)(unsafe.Add(mBase, uint32(v6968)+8))
	F_WebPSafeFree(m, v6990)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v6968)+16)) = v6986 + v6984
	*(*int32)(unsafe.Add(mBase, uint32(v6968)+12)) = v6986
	*(*int32)(unsafe.Add(mBase, uint32(v6968)+8)) = v6986
	goto L781
L783:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v6968)+20)) = int32(1)
	goto L781
L784:
	;
	v6999 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v6999 != 0 {
		goto L786
	} else {
		goto L787
	}
L785:
	;
	v9779 = v301
	v9780 = v302
	goto L779
L786:
	;
	goto L785
L787:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L786
L788:
	;
	v7021 = int32(0)
	v7023 = base.I32_div_s(v6733, int32(5))
	v7026 = F_VP8LHashChainFill(m, v293, v283, v6888, v6886, v276, v286, v275, v7023, v267+int32(60))
	mBase = m.M
	if v7026 == v7021 {
		v9712 = v301
		v9713 = v302
		goto L799
	} else {
		goto L800
	}
L789:
	;
	v7018 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v7018 != 0 {
		goto L797
	} else {
		goto L798
	}
L790:
	;
	if v6945 == int32(0) {
		goto L789
	} else {
		goto L791
	}
L791:
	;
	v7007 = v267 + int32(64)
	v7011 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v6922), int32(4))
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v7007))) = v7011
	if v7011 != 0 {
		goto L793
	} else {
		goto L794
	}
L792:
	;
	if v7011 != 0 {
		goto L788
	} else {
		goto L795
	}
L793:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7007)+4)) = v6922
	goto L792
L794:
	;
	goto L792
L795:
	;
	goto L789
L796:
	;
	v9779 = v301
	v9780 = v302
	goto L779
L797:
	;
	goto L796
L798:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L797
L799:
	;
	v9811 = v265
	v9813 = v267
	v9814 = v268
	v9817 = v271
	v9818 = v272
	v9819 = v273
	v9820 = v274
	v9821 = v275
	v9822 = v276
	v9823 = v277
	v9824 = v278
	v9826 = v280
	v9829 = v283
	v9830 = v284
	v9831 = v285
	v9832 = v286
	v9833 = v287
	v9834 = v288
	v9835 = v289
	v9836 = v290
	v9837 = v291
	v9838 = v292
	v9839 = v293
	v9840 = v294
	v9841 = v295
	v9842 = v296
	v9843 = v297
	v9844 = v298
	v9845 = v299
	v9846 = v300
	v9847 = v9712
	v9848 = v9713
	v9850 = int32(0)
	v9852 = v6912
	v9854 = v7021
	v9855 = v6945
	goto L778
L800:
	;
	v7029 = *(*int32)(unsafe.Add(mBase, uint32(v274)+48))
	v7030 = *(*int32)(unsafe.Add(mBase, uint32(v333)+24))
	if int32(1) < v7030 {
		goto L803
	} else {
		goto L804
	}
L801:
	;
	v9673 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v9673 != 0 {
		goto L1132
	} else {
		goto L1133
	}
L802:
	;
	if v7111 < int32(1) {
		v9608 = v301
		v9609 = v302
		goto L821
	} else {
		goto L822
	}
L803:
	;
	v7037 = v267 + int32(1088)
	v7047 = *(*int32)(unsafe.Add(mBase, uint32(v7037)+12))
	v7048 = *(*int32)(unsafe.Add(mBase, uint32(v7037)+8))
	v7049 = v7047 - v7048
	v7051 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v7052 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7053 = v7051 - v7052
	v7055 = base.I64_extend_i32_u(v7049) + base.I64_extend_i32_u(v7053)
	if base.Ui64(v7055) < base.Ui64(int64(4294967296)) {
		goto L807
	} else {
		goto L808
	}
L804:
	;
	v7033 = *(*int32)(unsafe.Add(mBase, uint32(v333)+12))
	if v7033 == int32(0) {
		v7111 = v7030
		goto L802
	} else {
		goto L805
	}
L805:
	;
	goto L803
L806:
	;
	if v7107 == int32(0) {
		goto L801
	} else {
		goto L820
	}
L807:
	;
	v7061 = base.I32_wrap_i64(v7055)
	v7062 = *(*int32)(unsafe.Add(mBase, uint32(v7037)+16))
	v7063 = v7062 - v7048
	if v7062 == v7048 {
		goto L810
	} else {
		goto L811
	}
L808:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+20)) = int32(1)
	v7107 = int32(0)
	goto L806
L809:
	;
	v7094 = F_memcpy(m, v7093, v7092, v7053)
	mBase = m.M
	v7095 = *(*int64)(unsafe.Add(mBase, uint32(v268)))
	*(*int64)(unsafe.Add(mBase, uint32(v7037))) = v7095
	v7097 = *(*int32)(unsafe.Add(mBase, uint32(v268)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+20)) = v7097
	v7099 = *(*int32)(unsafe.Add(mBase, uint32(v7037)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+12)) = v7099 + v7053
	v7107 = int32(1)
	goto L806
L810:
	;
	v7070 = int32(base.Ui32(v7063*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7061) < base.Ui32(v7070) {
		goto L814
	} else {
		goto L815
	}
L811:
	;
	if base.Ui32(v7063) < base.Ui32(v7061) {
		goto L810
	} else {
		goto L812
	}
L812:
	;
	v7092 = v7052
	v7093 = v7048
	goto L809
L813:
	;
	if v7047 == v7048 {
		goto L818
	} else {
		goto L819
	}
L814:
	;
	v7072 = v7070
	goto L816
L815:
	;
	v7072 = v7061
	goto L816
L816:
	;
	v7076 = v7072&int32(-1024) + int32(1024)
	v7077 = F_WebPSafeMalloc(m, int64(1), v7076)
	mBase = m.M
	if v7077 != 0 {
		goto L813
	} else {
		goto L817
	}
L817:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+20)) = int32(1)
	v7107 = int32(0)
	goto L806
L818:
	;
	v7084 = *(*int32)(unsafe.Add(mBase, uint32(v7037)+8))
	F_WebPSafeFree(m, v7084)
	mBase = m.M
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+8)) = v7077
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+16)) = v7077 + v7076
	*(*int32)(unsafe.Add(mBase, uint32(v7037)+12)) = v7077 + v7049
	v7091 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7092 = v7091
	v7093 = v7077
	goto L809
L819:
	;
	v7082 = *(*int32)(unsafe.Add(mBase, uint32(v7037)+8))
	v7083 = F_memcpy(m, v7077, v7082, v7049)
	mBase = m.M
	goto L818
L820:
	;
	v7110 = *(*int32)(unsafe.Add(mBase, uint32(v333)+24))
	v7111 = v7110
	goto L802
L821:
	;
	v9639 = v267 + int32(1088)
	v9648 = *(*int64)(unsafe.Add(mBase, uint32(v268)))
	v9649 = *(*int64)(unsafe.Add(mBase, uint32(v9639)))
	*(*int64)(unsafe.Add(mBase, uint32(v268))) = v9649
	v9652 = v268 + int32(16)
	v9653 = *(*int64)(unsafe.Add(mBase, uint32(v9652)))
	v9655 = v267 + int32(1104)
	v9656 = *(*int64)(unsafe.Add(mBase, uint32(v9655)))
	*(*int64)(unsafe.Add(mBase, uint32(v9652))) = v9656
	v9659 = v268 + int32(8)
	v9660 = *(*int64)(unsafe.Add(mBase, uint32(v9659)))
	v9662 = v267 + int32(1096)
	v9663 = *(*int64)(unsafe.Add(mBase, uint32(v9662)))
	*(*int64)(unsafe.Add(mBase, uint32(v9659))) = v9663
	*(*int64)(unsafe.Add(mBase, uint32(v9639))) = v9648
	*(*int64)(unsafe.Add(mBase, uint32(v9655))) = v9653
	*(*int64)(unsafe.Add(mBase, uint32(v9662))) = v9660
	goto L1130
L822:
	;
	if v7029 != 0 {
		goto L823
	} else {
		goto L824
	}
L823:
	;
	v7116 = v7029
	goto L825
L824:
	;
	v7116 = int32(10)
	goto L825
L825:
	;
	v7129 = v7111
	v7161 = v301
	v7162 = v302
	v7184 = int32(-1)
	v7185 = int32(0)
	goto L826
L826:
	;
	v7193 = v333 + int32(8) + v7185<<(uint(int32(3))%32)
	v7194 = *(*int32)(unsafe.Add(mBase, uint32(v7193)))
	v7195 = *(*int32)(unsafe.Add(mBase, uint32(v7193)+4))
	v7198 = base.I32_div_s(v6733-v7023, v7129)
	v7200 = base.I32_div_s(v7198, int32(4))
	v7203 = F_VP8LGetBackwardReferences(m, v6886, v276, v6888, v283, v286, v7194, v7116, v7195, v293, v297, v267+int32(_a_F_EncodeStreamHook_6), v275, v7200, v267+int32(60))
	mBase = m.M
	if v7203 == int32(0) {
		v9779 = v7161
		v9780 = v7162
		goto L779
	} else {
		goto L828
	}
L827:
	;
	v9608 = v9538
	v9609 = v9539
	goto L821
L828:
	;
	v7225 = int32(0)
	v7245 = v7161
	v7246 = v7162
	v7268 = v7184
	v7271 = v7198 - v7200
	v7272 = int32(1)
	goto L830
L829:
	;
	v9569 = v7185 + int32(1)
	v9570 = *(*int32)(unsafe.Add(mBase, uint32(v333)+24))
	if v9569 < v9570 {
		v7129 = v9570
		v7161 = v9538
		v7162 = v9539
		v7184 = v9561
		v7185 = v9569
		goto L826
	} else {
		goto L1129
	}
L830:
	;
	v7275 = *(*int32)(unsafe.Add(mBase, uint32(v267)+uint32(_c_F_EncodeStreamHook[0])))
	*(*int32)(unsafe.Add(mBase, uint32(v267)+uint32(_c_F_EncodeStreamHook[1]))) = v6887
	if v7225&int32(1) == int32(0) {
		goto L832
	} else {
		goto L833
	}
L831:
	;
	v9538 = v9489
	v9539 = v9490
	v9561 = v9491
	goto L829
L832:
	;
	v7284 = v267 + int32(2112)
	v7285 = *(*int64)(unsafe.Add(mBase, uint32(v7284)))
	*(*int64)(unsafe.Add(mBase, uint32(v268))) = v7285
	v7287 = *(*int32)(unsafe.Add(mBase, uint32(v7284)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = v7287
	v7289 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7290 = *(*int32)(unsafe.Add(mBase, uint32(v7284)+12))
	v7291 = *(*int32)(unsafe.Add(mBase, uint32(v7284)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7289 + (v7290 - v7291)
	goto L835
L833:
	;
	if v7275 == int32(0) {
		v9538 = v7245
		v9539 = v7246
		v9561 = v7268
		goto L829
	} else {
		goto L834
	}
L834:
	;
	goto L832
L835:
	;
	v7295 = int32(0)
	if v7272&int32(1) != 0 {
		goto L836
	} else {
		goto L837
	}
L836:
	;
	v7299 = v7275
	goto L838
L837:
	;
	v7299 = v7295
	goto L838
L838:
	;
	v7300 = F_VP8LAllocateHistogramSet(m, v6922, v7299)
	mBase = m.M
	v7304 = int32(_a_F_EncodeStreamHook_7)
	if int32(0) < v7299 {
		goto L841
	} else {
		goto L842
	}
L839:
	;
	if v7300 == int32(0) {
		goto L846
	} else {
		goto L847
	}
L840:
	;
	goto L839
L841:
	;
	v7309 = int32(4)<<(uint(v7299)%32) + v7304
	goto L843
L842:
	;
	v7309 = v7304
	goto L843
L843:
	;
	v7312 = F_WebPSafeMalloc(m, base.I64_extend_i32_s(v7309), int32(1))
	mBase = m.M
	if v7312 == int32(0) {
		goto L840
	} else {
		goto L844
	}
L844:
	;
	*(*int64)(unsafe.Add(mBase, uint32(v7312)+3240)) = int64(-1)
	*(*int32)(unsafe.Add(mBase, uint32(v7312)+3236)) = v7299
	*(*int32)(unsafe.Add(mBase, uint32(v7312)+3304)) = int32(16843009)
	v7320 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v7312)+3256)) = v7320
	*(*int32)(unsafe.Add(mBase, uint32(v7312))) = v7312 + int32(3312)
	v7327 = int32(-1)
	*(*uint16)(unsafe.Add(mBase, uint32(v7312+int32(3248)))) = uint16(v7327)
	v7331 = int32(1)
	*(*uint8)(unsafe.Add(mBase, uint32(v7312+int32(3308)))) = uint8(v7331)
	*(*int64)(unsafe.Add(mBase, uint32(v7312+int32(3264)))) = v7320
	*(*int64)(unsafe.Add(mBase, uint32(v7312+int32(3272)))) = v7320
	*(*int64)(unsafe.Add(mBase, uint32(v7312+int32(3280)))) = v7320
	*(*int64)(unsafe.Add(mBase, uint32(v7312+int32(3288)))) = v7320
	*(*int64)(unsafe.Add(mBase, uint32(v7312+int32(3296)))) = v7320
	goto L840
L845:
	;
	v7402 = F_GetHuffBitLengthsAndCodes(m, v7300, v7392)
	mBase = m.M
	if v7402 != 0 {
		goto L861
	} else {
		goto L862
	}
L846:
	;
	v7399 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v7399 != 0 {
		goto L859
	} else {
		goto L860
	}
L847:
	;
	if v7312 == int32(0) {
		goto L846
	} else {
		goto L848
	}
L848:
	;
	v7359 = v297 + v7225*int32(24)
	v7360 = *(*int32)(unsafe.Add(mBase, uint32(v267)+uint32(_c_F_EncodeStreamHook[1])))
	v7362 = base.I32_div_s(v7271, int32(3))
	v7365 = F_VP8LGetHistoImageSymbols(m, v6886, v276, v7359, v283, v286, v7360, v7299, v7300, v7312, v6945, v275, v7362, v267+int32(60))
	mBase = m.M
	if v7365 != 0 {
		goto L849
	} else {
		goto L850
	}
L849:
	;
	v7368 = *(*int32)(unsafe.Add(mBase, uint32(v7300)))
	v7371 = base.I64_extend_i32_u(v7368 * int32(5))
	v7372 = int32(12)
	if v7371 == int64(0) {
		goto L853
	} else {
		goto L854
	}
L850:
	;
	v7366 = int32(0)
	v9945 = v265
	v9947 = v267
	v9948 = v268
	v9949 = v7300
	v9950 = v7312
	v9951 = v271
	v9952 = v272
	v9953 = v273
	v9954 = v274
	v9955 = v275
	v9956 = v276
	v9957 = v277
	v9958 = v278
	v9960 = v280
	v9963 = v283
	v9964 = v284
	v9965 = v285
	v9966 = v286
	v9967 = v287
	v9968 = v288
	v9969 = v289
	v9970 = v290
	v9971 = v291
	v9972 = v292
	v9973 = v293
	v9974 = v294
	v9975 = v295
	v9976 = v296
	v9977 = v297
	v9978 = v298
	v9979 = v299
	v9980 = v300
	v9981 = v7245
	v9982 = v7246
	v9984 = v7366
	v9986 = v6912
	v9988 = v7366
	v9989 = v6945
	goto L776
L851:
	;
	if v7392 != 0 {
		goto L845
	} else {
		goto L857
	}
L852:
	;
	goto L851
L853:
	;
	v7390 = F_calloc(m, base.I32_wrap_i64(v7371), v7372)
	mBase = m.M
	v7392 = v7390
	goto L852
L854:
	;
	v7379 = base.I64_div_u_s(int64(2147418112), v7371)
	v7380 = int32(0)
	v7381 = base.I64_extend_i32_u(v7372)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7381*v7371) {
		v7392 = v7380
		goto L852
	} else {
		goto L855
	}
L855:
	;
	if base.Ui64(v7379) < base.Ui64(v7381) {
		v7392 = v7380
		goto L852
	} else {
		goto L856
	}
L856:
	;
	goto L853
L857:
	;
	goto L846
L858:
	;
	v9945 = v265
	v9947 = v267
	v9948 = v268
	v9949 = v7300
	v9950 = v7312
	v9951 = v271
	v9952 = v272
	v9953 = v273
	v9954 = v274
	v9955 = v275
	v9956 = v276
	v9957 = v277
	v9958 = v278
	v9960 = v280
	v9963 = v283
	v9964 = v284
	v9965 = v285
	v9966 = v286
	v9967 = v287
	v9968 = v288
	v9969 = v289
	v9970 = v290
	v9971 = v291
	v9972 = v292
	v9973 = v293
	v9974 = v294
	v9975 = v295
	v9976 = v296
	v9977 = v297
	v9978 = v298
	v9979 = v299
	v9980 = v300
	v9981 = v7245
	v9982 = v7246
	v9984 = v7295
	v9986 = v6912
	v9988 = int32(0)
	v9989 = v6945
	goto L776
L859:
	;
	goto L858
L860:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L859
L861:
	;
	F_WebPSafeFree(m, v7300)
	mBase = m.M
	goto L866
L862:
	;
	v7404 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v7404 != 0 {
		goto L864
	} else {
		goto L865
	}
L863:
	;
	v9945 = v265
	v9947 = v267
	v9948 = v268
	v9949 = v7300
	v9950 = v7312
	v9951 = v271
	v9952 = v272
	v9953 = v273
	v9954 = v274
	v9955 = v275
	v9956 = v276
	v9957 = v277
	v9958 = v278
	v9960 = v280
	v9963 = v283
	v9964 = v284
	v9965 = v285
	v9966 = v286
	v9967 = v287
	v9968 = v288
	v9969 = v289
	v9970 = v290
	v9971 = v291
	v9972 = v292
	v9973 = v293
	v9974 = v294
	v9975 = v295
	v9976 = v296
	v9977 = v297
	v9978 = v298
	v9979 = v299
	v9980 = v300
	v9981 = v7245
	v9982 = v7246
	v9984 = v7295
	v9986 = v6912
	v9988 = v7392
	v9989 = v6945
	goto L776
L864:
	;
	goto L863
L865:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L864
L866:
	;
	F_WebPSafeFree(m, v7312)
	mBase = m.M
	goto L867
L867:
	;
	if v7299 < int32(1) {
		goto L869
	} else {
		goto L870
	}
L868:
	;
	v7866 = v7271 - v7362
	v7867 = int32(0)
	switch v6922 {
	case 0:
		goto L951
	case 1:
		v7965 = v7867
		v8000 = v7867
		goto L952
	default:
		goto L953
	}
L869:
	;
	v7714 = int32(0)
	v7715 = int32(1)
	goto L925
L870:
	;
	v7411 = int32(1)
	goto L873
L871:
	;
	v7563 = int32(4)
	goto L899
L872:
	;
	goto L871
L873:
	;
	v7424 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v7425 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v7425+v7411 < int32(32) {
		goto L876
	} else {
		goto L877
	}
L874:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v7526 + v7524
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v7523<<(uint(v7526)%32) | v7525
	goto L872
L875:
	;
	v7443 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v7444 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v7448 = v7440
	v7449 = v7441
	v7450 = v7444
	v7451 = v7443
	goto L879
L876:
	;
	if v7425 < int32(16) {
		v7523 = v7411
		v7524 = v7411
		v7525 = v7424
		v7526 = v7425
		goto L874
	} else {
		goto L878
	}
L877:
	;
	v7429 = int32(32)
	v7430 = v7429 - v7425
	v7438 = int32(base.Ui32(v7411) >> (uint(v7430) % 32))
	v7439 = v7411 - v7430
	v7440 = v7411<<(uint(v7425)%32) | v7424
	v7441 = v7429
	goto L875
L878:
	;
	v7438 = v7411
	v7439 = v7411
	v7440 = v7424
	v7441 = v7425
	goto L875
L879:
	;
	if base.Ui32(v7450+int32(2)) <= base.Ui32(v7451) {
		v7506 = v7450
		v7507 = v7451
		goto L881
	} else {
		goto L882
	}
L880:
	;
	v7523 = v7438
	v7524 = v7439
	v7525 = v7519
	v7526 = v7517
	goto L874
L881:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7506))) = uint16(v7448)
	v7514 = v7506 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7514
	v7517 = v7449 + int32(-16)
	v7519 = int32(base.Ui32(v7448) >> (uint(int32(16)) % 32))
	if int32(31) < v7449 {
		v7448 = v7519
		v7449 = v7517
		v7450 = v7514
		v7451 = v7507
		goto L879
	} else {
		goto L896
	}
L882:
	;
	v7460 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7461 = v7451 - v7460
	v7464 = base.I64_extend_i32_s(v7461) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7464) {
		v7488 = v7460
		goto L884
	} else {
		goto L885
	}
L883:
	;
	if v7450 == v7460 {
		goto L894
	} else {
		goto L895
	}
L884:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7488
	goto L871
L885:
	;
	v7467 = v7450 - v7460
	v7469 = v7464 + base.I64_extend_i32_u(v7467)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7469) {
		v7488 = v7460
		goto L884
	} else {
		goto L886
	}
L886:
	;
	v7472 = base.I32_wrap_i64(v7469)
	if v7451 == v7460 {
		goto L887
	} else {
		goto L888
	}
L887:
	;
	v7479 = int32(base.Ui32(v7461*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7472) < base.Ui32(v7479) {
		goto L890
	} else {
		goto L891
	}
L888:
	;
	if base.Ui32(v7472) <= base.Ui32(v7461) {
		v7506 = v7450
		v7507 = v7451
		goto L881
	} else {
		goto L889
	}
L889:
	;
	goto L887
L890:
	;
	v7481 = v7479
	goto L892
L891:
	;
	v7481 = v7472
	goto L892
L892:
	;
	v7485 = v7481&int32(-1024) + int32(1024)
	v7486 = F_WebPSafeMalloc(m, int64(1), v7485)
	mBase = m.M
	if v7486 != 0 {
		goto L883
	} else {
		goto L893
	}
L893:
	;
	v7487 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7488 = v7487
	goto L884
L894:
	;
	v7499 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v7499)
	mBase = m.M
	v7501 = v7486 + v7485
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v7501
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v7486
	v7506 = v7486 + v7467
	v7507 = v7501
	goto L881
L895:
	;
	v7497 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7498 = F_memcpy(m, v7486, v7497, v7467)
	mBase = m.M
	goto L894
L896:
	;
	goto L880
L897:
	;
	goto L868
L898:
	;
	goto L897
L899:
	;
	v7575 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v7576 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v7576+v7563 < int32(32) {
		goto L902
	} else {
		goto L903
	}
L900:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v7677 + v7675
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v7674<<(uint(v7677)%32) | v7676
	goto L898
L901:
	;
	v7594 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v7595 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v7599 = v7591
	v7600 = v7592
	v7601 = v7595
	v7602 = v7594
	goto L905
L902:
	;
	if v7576 < int32(16) {
		v7674 = v7299
		v7675 = v7563
		v7676 = v7575
		v7677 = v7576
		goto L900
	} else {
		goto L904
	}
L903:
	;
	v7580 = int32(32)
	v7581 = v7580 - v7576
	v7589 = int32(base.Ui32(v7299) >> (uint(v7581) % 32))
	v7590 = v7563 - v7581
	v7591 = v7299<<(uint(v7576)%32) | v7575
	v7592 = v7580
	goto L901
L904:
	;
	v7589 = v7299
	v7590 = v7563
	v7591 = v7575
	v7592 = v7576
	goto L901
L905:
	;
	if base.Ui32(v7601+int32(2)) <= base.Ui32(v7602) {
		v7657 = v7601
		v7658 = v7602
		goto L907
	} else {
		goto L908
	}
L906:
	;
	v7674 = v7589
	v7675 = v7590
	v7676 = v7670
	v7677 = v7668
	goto L900
L907:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7657))) = uint16(v7599)
	v7665 = v7657 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7665
	v7668 = v7600 + int32(-16)
	v7670 = int32(base.Ui32(v7599) >> (uint(int32(16)) % 32))
	if int32(31) < v7600 {
		v7599 = v7670
		v7600 = v7668
		v7601 = v7665
		v7602 = v7658
		goto L905
	} else {
		goto L922
	}
L908:
	;
	v7611 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7612 = v7602 - v7611
	v7615 = base.I64_extend_i32_s(v7612) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7615) {
		v7639 = v7611
		goto L910
	} else {
		goto L911
	}
L909:
	;
	if v7601 == v7611 {
		goto L920
	} else {
		goto L921
	}
L910:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7639
	goto L897
L911:
	;
	v7618 = v7601 - v7611
	v7620 = v7615 + base.I64_extend_i32_u(v7618)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7620) {
		v7639 = v7611
		goto L910
	} else {
		goto L912
	}
L912:
	;
	v7623 = base.I32_wrap_i64(v7620)
	if v7602 == v7611 {
		goto L913
	} else {
		goto L914
	}
L913:
	;
	v7630 = int32(base.Ui32(v7612*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7623) < base.Ui32(v7630) {
		goto L916
	} else {
		goto L917
	}
L914:
	;
	if base.Ui32(v7623) <= base.Ui32(v7612) {
		v7657 = v7601
		v7658 = v7602
		goto L907
	} else {
		goto L915
	}
L915:
	;
	goto L913
L916:
	;
	v7632 = v7630
	goto L918
L917:
	;
	v7632 = v7623
	goto L918
L918:
	;
	v7636 = v7632&int32(-1024) + int32(1024)
	v7637 = F_WebPSafeMalloc(m, int64(1), v7636)
	mBase = m.M
	if v7637 != 0 {
		goto L909
	} else {
		goto L919
	}
L919:
	;
	v7638 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7639 = v7638
	goto L910
L920:
	;
	v7650 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v7650)
	mBase = m.M
	v7652 = v7637 + v7636
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v7652
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v7637
	v7657 = v7637 + v7618
	v7658 = v7652
	goto L907
L921:
	;
	v7648 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7649 = F_memcpy(m, v7637, v7648, v7618)
	mBase = m.M
	goto L920
L922:
	;
	goto L906
L923:
	;
	goto L868
L924:
	;
	goto L923
L925:
	;
	v7727 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v7728 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v7728+v7715 < int32(32) {
		goto L928
	} else {
		goto L929
	}
L926:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v7829 + v7827
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v7826<<(uint(v7829)%32) | v7828
	goto L924
L927:
	;
	v7746 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v7747 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v7751 = v7743
	v7752 = v7744
	v7753 = v7747
	v7754 = v7746
	goto L931
L928:
	;
	if v7728 < int32(16) {
		v7826 = v7714
		v7827 = v7715
		v7828 = v7727
		v7829 = v7728
		goto L926
	} else {
		goto L930
	}
L929:
	;
	v7732 = int32(32)
	v7733 = v7732 - v7728
	v7741 = int32(base.Ui32(v7714) >> (uint(v7733) % 32))
	v7742 = v7715 - v7733
	v7743 = v7714<<(uint(v7728)%32) | v7727
	v7744 = v7732
	goto L927
L930:
	;
	v7741 = v7714
	v7742 = v7715
	v7743 = v7727
	v7744 = v7728
	goto L927
L931:
	;
	if base.Ui32(v7753+int32(2)) <= base.Ui32(v7754) {
		v7809 = v7753
		v7810 = v7754
		goto L933
	} else {
		goto L934
	}
L932:
	;
	v7826 = v7741
	v7827 = v7742
	v7828 = v7822
	v7829 = v7820
	goto L926
L933:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7809))) = uint16(v7751)
	v7817 = v7809 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7817
	v7820 = v7752 + int32(-16)
	v7822 = int32(base.Ui32(v7751) >> (uint(int32(16)) % 32))
	if int32(31) < v7752 {
		v7751 = v7822
		v7752 = v7820
		v7753 = v7817
		v7754 = v7810
		goto L931
	} else {
		goto L948
	}
L934:
	;
	v7763 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7764 = v7754 - v7763
	v7767 = base.I64_extend_i32_s(v7764) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7767) {
		v7791 = v7763
		goto L936
	} else {
		goto L937
	}
L935:
	;
	if v7753 == v7763 {
		goto L946
	} else {
		goto L947
	}
L936:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v7791
	goto L923
L937:
	;
	v7770 = v7753 - v7763
	v7772 = v7767 + base.I64_extend_i32_u(v7770)
	if base.Ui64(int64(4294967295)) < base.Ui64(v7772) {
		v7791 = v7763
		goto L936
	} else {
		goto L938
	}
L938:
	;
	v7775 = base.I32_wrap_i64(v7772)
	if v7754 == v7763 {
		goto L939
	} else {
		goto L940
	}
L939:
	;
	v7782 = int32(base.Ui32(v7764*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v7775) < base.Ui32(v7782) {
		goto L942
	} else {
		goto L943
	}
L940:
	;
	if base.Ui32(v7775) <= base.Ui32(v7764) {
		v7809 = v7753
		v7810 = v7754
		goto L933
	} else {
		goto L941
	}
L941:
	;
	goto L939
L942:
	;
	v7784 = v7782
	goto L944
L943:
	;
	v7784 = v7775
	goto L944
L944:
	;
	v7788 = v7784&int32(-1024) + int32(1024)
	v7789 = F_WebPSafeMalloc(m, int64(1), v7788)
	mBase = m.M
	if v7789 != 0 {
		goto L935
	} else {
		goto L945
	}
L945:
	;
	v7790 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7791 = v7790
	goto L936
L946:
	;
	v7802 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v7802)
	mBase = m.M
	v7804 = v7789 + v7788
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v7804
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v7789
	v7809 = v7789 + v7770
	v7810 = v7804
	goto L933
L947:
	;
	v7800 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v7801 = F_memcpy(m, v7789, v7800, v7770)
	mBase = m.M
	goto L946
L948:
	;
	goto L932
L949:
	;
	v9426 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v9427 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v9428 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	v9429 = *(*int32)(unsafe.Add(mBase, uint32(v267)+uint32(_c_F_EncodeStreamHook[1])))
	v9430 = F_StoreImageToBitMask(m, v268, v6886, v9429, v7359, v6945, v7392, v275)
	mBase = m.M
	if v9430 != 0 {
		goto L1120
	} else {
		goto L1121
	}
L950:
	;
	v9357 = *(*int32)(unsafe.Add(mBase, uint32(v275)+92))
	if v9357 != 0 {
		goto L1118
	} else {
		goto L1119
	}
L951:
	;
	v9112 = int32(0)
	v9114 = int32(1)
	goto L1086
L952:
	;
	if v6922&int32(1) == int32(0) {
		v8040 = v7965
		goto L963
	} else {
		goto L964
	}
L953:
	;
	v7869 = int32(0)
	v7875 = v6945
	v7876 = v7869
	v7911 = v7869
	goto L954
L954:
	;
	v7937 = *(*int32)(unsafe.Add(mBase, uint32(v7875)))
	v7938 = int32(8)
	*(*int32)(unsafe.Add(mBase, uint32(v7875))) = v7937 << (uint(v7938) % 32)
	v7942 = v7875 + int32(4)
	v7943 = *(*int32)(unsafe.Add(mBase, uint32(v7942)))
	*(*int32)(unsafe.Add(mBase, uint32(v7942))) = v7943 << (uint(v7938) % 32)
	if base.Ui32(v7937) < base.Ui32(v7876) {
		goto L956
	} else {
		goto L957
	}
L955:
	;
	v7965 = v7954
	v8000 = v7958
	goto L952
L956:
	;
	v7950 = v7876
	goto L958
L957:
	;
	v7950 = v7937 + int32(1)
	goto L958
L958:
	;
	if base.Ui32(v7943) < base.Ui32(v7950) {
		goto L959
	} else {
		goto L960
	}
L959:
	;
	v7954 = v7950
	goto L961
L960:
	;
	v7954 = v7943 + int32(1)
	goto L961
L961:
	;
	v7958 = v7911 + int32(2)
	if v6922&int32(-2) != v7958 {
		v7875 = v7875 + int32(8)
		v7876 = v7954
		v7911 = v7958
		goto L954
	} else {
		goto L962
	}
L962:
	;
	goto L955
L963:
	;
	v8041 = int32(1)
	v8042 = base.B2i32(base.Ui32(v8041) < base.Ui32(v8040))
	goto L970
L964:
	;
	v8030 = v6945 + v8000<<(uint(int32(2))%32)
	v8031 = *(*int32)(unsafe.Add(mBase, uint32(v8030)))
	*(*int32)(unsafe.Add(mBase, uint32(v8030))) = v8031 << (uint(int32(8)) % 32)
	if base.Ui32(v8031) < base.Ui32(v7965) {
		goto L965
	} else {
		goto L966
	}
L965:
	;
	v8038 = v7965
	goto L967
L966:
	;
	v8038 = v8031 + int32(1)
	goto L967
L967:
	;
	v8040 = v8038
	goto L963
L968:
	;
	if base.Ui32(v8040) < base.Ui32(int32(2)) {
		goto L997
	} else {
		goto L998
	}
L969:
	;
	goto L968
L970:
	;
	v8055 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v8056 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v8056+v8041 < int32(32) {
		goto L973
	} else {
		goto L974
	}
L971:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v8157 + v8155
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v8154<<(uint(v8157)%32) | v8156
	goto L969
L972:
	;
	v8074 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v8075 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v8079 = v8071
	v8080 = v8072
	v8081 = v8075
	v8082 = v8074
	goto L976
L973:
	;
	if v8056 < int32(16) {
		v8154 = v8042
		v8155 = v8041
		v8156 = v8055
		v8157 = v8056
		goto L971
	} else {
		goto L975
	}
L974:
	;
	v8060 = int32(32)
	v8061 = v8060 - v8056
	v8069 = int32(base.Ui32(v8042) >> (uint(v8061) % 32))
	v8070 = v8041 - v8061
	v8071 = v8042<<(uint(v8056)%32) | v8055
	v8072 = v8060
	goto L972
L975:
	;
	v8069 = v8042
	v8070 = v8041
	v8071 = v8055
	v8072 = v8056
	goto L972
L976:
	;
	if base.Ui32(v8081+int32(2)) <= base.Ui32(v8082) {
		v8137 = v8081
		v8138 = v8082
		goto L978
	} else {
		goto L979
	}
L977:
	;
	v8154 = v8069
	v8155 = v8070
	v8156 = v8150
	v8157 = v8148
	goto L971
L978:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v8137))) = uint16(v8079)
	v8145 = v8137 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v8145
	v8148 = v8080 + int32(-16)
	v8150 = int32(base.Ui32(v8079) >> (uint(int32(16)) % 32))
	if int32(31) < v8080 {
		v8079 = v8150
		v8080 = v8148
		v8081 = v8145
		v8082 = v8138
		goto L976
	} else {
		goto L993
	}
L979:
	;
	v8091 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v8092 = v8082 - v8091
	v8095 = base.I64_extend_i32_s(v8092) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8095) {
		v8119 = v8091
		goto L981
	} else {
		goto L982
	}
L980:
	;
	if v8081 == v8091 {
		goto L991
	} else {
		goto L992
	}
L981:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v8119
	goto L968
L982:
	;
	v8098 = v8081 - v8091
	v8100 = v8095 + base.I64_extend_i32_u(v8098)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8100) {
		v8119 = v8091
		goto L981
	} else {
		goto L983
	}
L983:
	;
	v8103 = base.I32_wrap_i64(v8100)
	if v8082 == v8091 {
		goto L984
	} else {
		goto L985
	}
L984:
	;
	v8110 = int32(base.Ui32(v8092*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v8103) < base.Ui32(v8110) {
		goto L987
	} else {
		goto L988
	}
L985:
	;
	if base.Ui32(v8103) <= base.Ui32(v8092) {
		v8137 = v8081
		v8138 = v8082
		goto L978
	} else {
		goto L986
	}
L986:
	;
	goto L984
L987:
	;
	v8112 = v8110
	goto L989
L988:
	;
	v8112 = v8103
	goto L989
L989:
	;
	v8116 = v8112&int32(-1024) + int32(1024)
	v8117 = F_WebPSafeMalloc(m, int64(1), v8116)
	mBase = m.M
	if v8117 != 0 {
		goto L980
	} else {
		goto L990
	}
L990:
	;
	v8118 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v8119 = v8118
	goto L981
L991:
	;
	v8130 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v8130)
	mBase = m.M
	v8132 = v8117 + v8116
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v8132
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v8117
	v8137 = v8117 + v8098
	v8138 = v8132
	goto L978
L992:
	;
	v8128 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v8129 = F_memcpy(m, v8117, v8128, v8098)
	mBase = m.M
	goto L991
L993:
	;
	goto L977
L994:
	;
	v9294 = int32(0)
	goto L950
L995:
	;
	goto L1079
L996:
	;
	v8376 = v8040 * int32(5)
	v8377 = int32(1)
	if base.Ui32(v8377) < base.Ui32(v8376) {
		goto L1028
	} else {
		goto L1029
	}
L997:
	;
	if v8040 == int32(0) {
		goto L995
	} else {
		goto L1027
	}
L998:
	;
	F_VP8LOptimizeSampling(m, v6945, v6886, v276, v6887, int32(9), v267+int32(_a_F_EncodeStreamHook_8))
	mBase = m.M
	v8200 = *(*int32)(unsafe.Add(mBase, uint32(v267)+uint32(_c_F_EncodeStreamHook[1])))
	v8202 = v8200 + int32(-2)
	v8203 = int32(3)
	goto L1001
L999:
	;
	v8357 = *(*int32)(unsafe.Add(mBase, uint32(v267)+uint32(_c_F_EncodeStreamHook[1])))
	v8358 = int32(1) << (uint(v8357) % 32)
	v8364 = base.I32_div_s(v7866, int32(2))
	v8367 = F_EncodeImageNoHuffman(m, v268, v6945, v267+int32(64), v291, int32(base.Ui32(v8358+v6917)>>(uint(v8357)%32)), int32(base.Ui32(v8358+v292)>>(uint(v8357)%32)), v283, v286, v275, v8364, v267+int32(60))
	mBase = m.M
	if v8367 != 0 {
		goto L1025
	} else {
		goto L1026
	}
L1000:
	;
	goto L999
L1001:
	;
	v8215 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v8216 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v8216+v8203 < int32(32) {
		goto L1004
	} else {
		goto L1005
	}
L1002:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v8317 + v8315
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v8314<<(uint(v8317)%32) | v8316
	goto L1000
L1003:
	;
	v8234 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v8235 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v8239 = v8231
	v8240 = v8232
	v8241 = v8235
	v8242 = v8234
	goto L1007
L1004:
	;
	if v8216 < int32(16) {
		v8314 = v8202
		v8315 = v8203
		v8316 = v8215
		v8317 = v8216
		goto L1002
	} else {
		goto L1006
	}
L1005:
	;
	v8220 = int32(32)
	v8221 = v8220 - v8216
	v8229 = int32(base.Ui32(v8202) >> (uint(v8221) % 32))
	v8230 = v8203 - v8221
	v8231 = v8202<<(uint(v8216)%32) | v8215
	v8232 = v8220
	goto L1003
L1006:
	;
	v8229 = v8202
	v8230 = v8203
	v8231 = v8215
	v8232 = v8216
	goto L1003
L1007:
	;
	if base.Ui32(v8241+int32(2)) <= base.Ui32(v8242) {
		v8297 = v8241
		v8298 = v8242
		goto L1009
	} else {
		goto L1010
	}
L1008:
	;
	v8314 = v8229
	v8315 = v8230
	v8316 = v8310
	v8317 = v8308
	goto L1002
L1009:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v8297))) = uint16(v8239)
	v8305 = v8297 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v8305
	v8308 = v8240 + int32(-16)
	v8310 = int32(base.Ui32(v8239) >> (uint(int32(16)) % 32))
	if int32(31) < v8240 {
		v8239 = v8310
		v8240 = v8308
		v8241 = v8305
		v8242 = v8298
		goto L1007
	} else {
		goto L1024
	}
L1010:
	;
	v8251 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v8252 = v8242 - v8251
	v8255 = base.I64_extend_i32_s(v8252) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8255) {
		v8279 = v8251
		goto L1012
	} else {
		goto L1013
	}
L1011:
	;
	if v8241 == v8251 {
		goto L1022
	} else {
		goto L1023
	}
L1012:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v8279
	goto L999
L1013:
	;
	v8258 = v8241 - v8251
	v8260 = v8255 + base.I64_extend_i32_u(v8258)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8260) {
		v8279 = v8251
		goto L1012
	} else {
		goto L1014
	}
L1014:
	;
	v8263 = base.I32_wrap_i64(v8260)
	if v8242 == v8251 {
		goto L1015
	} else {
		goto L1016
	}
L1015:
	;
	v8270 = int32(base.Ui32(v8252*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v8263) < base.Ui32(v8270) {
		goto L1018
	} else {
		goto L1019
	}
L1016:
	;
	if base.Ui32(v8263) <= base.Ui32(v8252) {
		v8297 = v8241
		v8298 = v8242
		goto L1009
	} else {
		goto L1017
	}
L1017:
	;
	goto L1015
L1018:
	;
	v8272 = v8270
	goto L1020
L1019:
	;
	v8272 = v8263
	goto L1020
L1020:
	;
	v8276 = v8272&int32(-1024) + int32(1024)
	v8277 = F_WebPSafeMalloc(m, int64(1), v8276)
	mBase = m.M
	if v8277 != 0 {
		goto L1011
	} else {
		goto L1021
	}
L1021:
	;
	v8278 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v8279 = v8278
	goto L1012
L1022:
	;
	v8290 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v8290)
	mBase = m.M
	v8292 = v8277 + v8276
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v8292
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v8277
	v8297 = v8277 + v8258
	v8298 = v8292
	goto L1009
L1023:
	;
	v8288 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v8289 = F_memcpy(m, v8277, v8288, v8258)
	mBase = m.M
	goto L1022
L1024:
	;
	goto L1008
L1025:
	;
	v8374 = v7866 - v8364
	goto L996
L1026:
	;
	v9811 = v265
	v9813 = v267
	v9814 = v268
	v9817 = v271
	v9818 = v272
	v9819 = v273
	v9820 = v274
	v9821 = v275
	v9822 = v276
	v9823 = v277
	v9824 = v278
	v9826 = v280
	v9829 = v283
	v9830 = v284
	v9831 = v285
	v9832 = v286
	v9833 = v287
	v9834 = v288
	v9835 = v289
	v9836 = v290
	v9837 = v291
	v9838 = v292
	v9839 = v293
	v9840 = v294
	v9841 = v295
	v9842 = v296
	v9843 = v297
	v9844 = v298
	v9845 = v299
	v9846 = v300
	v9847 = v7245
	v9848 = v7246
	v9850 = int32(0)
	v9852 = v6912
	v9854 = v7392
	v9855 = v6945
	goto L778
L1027:
	;
	v8374 = v7866
	goto L996
L1028:
	;
	v8380 = v8376
	goto L1030
L1029:
	;
	v8380 = v8377
	goto L1030
L1030:
	;
	v8382 = v8380 & int32(3)
	v8383 = int32(0)
	if base.Ui32(v8376) < base.Ui32(int32(4)) {
		v8484 = v8383
		v8497 = v8383
		goto L1031
	} else {
		goto L1032
	}
L1031:
	;
	if v8382 == int32(0) {
		v8628 = v8484
		goto L1048
	} else {
		goto L1049
	}
L1032:
	;
	v8389 = int32(0)
	v8395 = v8389
	v8396 = v7392
	v8408 = v8389
	goto L1033
L1033:
	;
	v8457 = *(*int32)(unsafe.Add(mBase, uint32(v8396)))
	if v8457 < v8395 {
		goto L1035
	} else {
		goto L1036
	}
L1034:
	;
	v8484 = v8474
	v8497 = v8478
	goto L1031
L1035:
	;
	v8459 = v8395
	goto L1037
L1036:
	;
	v8459 = v8457
	goto L1037
L1037:
	;
	v8462 = *(*int32)(unsafe.Add(mBase, uint32(v8396+int32(12))))
	if v8462 < v8459 {
		goto L1038
	} else {
		goto L1039
	}
L1038:
	;
	v8464 = v8459
	goto L1040
L1039:
	;
	v8464 = v8462
	goto L1040
L1040:
	;
	v8467 = *(*int32)(unsafe.Add(mBase, uint32(v8396+int32(24))))
	if v8467 < v8464 {
		goto L1041
	} else {
		goto L1042
	}
L1041:
	;
	v8469 = v8464
	goto L1043
L1042:
	;
	v8469 = v8467
	goto L1043
L1043:
	;
	v8472 = *(*int32)(unsafe.Add(mBase, uint32(v8396+int32(36))))
	if v8472 < v8469 {
		goto L1044
	} else {
		goto L1045
	}
L1044:
	;
	v8474 = v8469
	goto L1046
L1045:
	;
	v8474 = v8472
	goto L1046
L1046:
	;
	v8478 = v8408 + int32(4)
	if v8380&int32(-4) != v8478 {
		v8395 = v8474
		v8396 = v8396 + int32(48)
		v8408 = v8478
		goto L1033
	} else {
		goto L1047
	}
L1047:
	;
	goto L1034
L1048:
	;
	v8690 = base.I64_extend_i32_u(v8628)
	v8691 = int32(2)
	if v8690 == int64(0) {
		goto L1058
	} else {
		goto L1059
	}
L1049:
	;
	v8555 = v8484
	v8556 = v7392 + v8497*int32(12)
	v8567 = v8382
	goto L1050
L1050:
	;
	v8617 = *(*int32)(unsafe.Add(mBase, uint32(v8556)))
	if v8617 < v8555 {
		goto L1052
	} else {
		goto L1053
	}
L1051:
	;
	v8628 = v8619
	goto L1048
L1052:
	;
	v8619 = v8555
	goto L1054
L1053:
	;
	v8619 = v8617
	goto L1054
L1054:
	;
	v8623 = v8567 + int32(-1)
	if v8623 != 0 {
		v8555 = v8619
		v8556 = v8556 + int32(12)
		v8567 = v8623
		goto L1050
	} else {
		goto L1055
	}
L1055:
	;
	goto L1051
L1056:
	;
	if v8712 == int32(0) {
		goto L994
	} else {
		goto L1062
	}
L1057:
	;
	goto L1056
L1058:
	;
	v8710 = F_malloc(m, base.I32_wrap_i64(v8690)*v8691)
	mBase = m.M
	v8712 = v8710
	goto L1057
L1059:
	;
	v8698 = base.I64_div_u_s(int64(2147418112), v8690)
	v8699 = int32(0)
	v8700 = base.I64_extend_i32_u(v8691)
	if base.Ui64(int64(4294967295)) < base.Ui64(v8700*v8690) {
		v8712 = v8699
		goto L1057
	} else {
		goto L1060
	}
L1060:
	;
	if base.Ui64(v8698) < base.Ui64(v8700) {
		v8712 = v8699
		goto L1057
	} else {
		goto L1061
	}
L1061:
	;
	goto L1058
L1062:
	;
	v8756 = int32(0)
	goto L1063
L1063:
	;
	v8785 = v7392 + v8756*int32(12)
	F_StoreHuffmanCode(m, v268, v6912, v8712, v8785)
	mBase = m.M
	v8787 = *(*int32)(unsafe.Add(mBase, uint32(v8785)))
	if v8787 < int32(1) {
		goto L1065
	} else {
		goto L1066
	}
L1065:
	;
	v9019 = v8756 + int32(1)
	if v9019 != v8380 {
		v8756 = v9019
		goto L1063
	} else {
		goto L1076
	}
L1066:
	;
	v8790 = *(*int32)(unsafe.Add(mBase, uint32(v8785)+4))
	v8797 = v8790
	v8808 = v8787
	v8809 = int32(0)
	goto L1067
L1067:
	;
	v8858 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v8797))))
	if v8858 == int32(0) {
		v8864 = v8809
		goto L1069
	} else {
		goto L1070
	}
L1068:
	;
	v8870 = int32(0)
	v8877 = v8870
	v8888 = v8870
	goto L1073
L1069:
	;
	v8869 = v8808 + int32(-1)
	if v8869 != 0 {
		v8797 = v8797 + int32(1)
		v8808 = v8869
		v8809 = v8864
		goto L1067
	} else {
		goto L1072
	}
L1070:
	;
	if int32(0) < v8809 {
		goto L1065
	} else {
		goto L1071
	}
L1071:
	;
	v8864 = int32(1)
	goto L1069
L1072:
	;
	goto L1068
L1073:
	;
	v8938 = *(*int32)(unsafe.Add(mBase, uint32(v8785)+4))
	v8940 = int32(0)
	*(*uint8)(unsafe.Add(mBase, uint32(v8938+v8888))) = uint8(v8940)
	v8942 = *(*int32)(unsafe.Add(mBase, uint32(v8785)+8))
	*(*uint16)(unsafe.Add(mBase, uint32(v8942+v8877))) = uint16(v8940)
	v8949 = v8888 + int32(1)
	v8950 = *(*int32)(unsafe.Add(mBase, uint32(v8785)))
	if v8949 < v8950 {
		v8877 = v8877 + int32(2)
		v8888 = v8949
		goto L1073
	} else {
		goto L1075
	}
L1074:
	;
	goto L1065
L1075:
	;
	goto L1074
L1076:
	;
	v9406 = v8712
	v9422 = v8374
	goto L949
L1077:
	;
	if v9041 != 0 {
		v9406 = v9041
		v9422 = v7866
		goto L949
	} else {
		goto L1083
	}
L1078:
	;
	goto L1077
L1079:
	;
	v9041 = F_malloc(m, base.I32_wrap_i64(int64(0))*int32(2))
	mBase = m.M
	goto L1078
L1083:
	;
	goto L994
L1084:
	;
	goto L1112
L1085:
	;
	goto L1084
L1086:
	;
	v9126 = *(*int32)(unsafe.Add(mBase, uint32(v268)))
	v9127 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	if v9127+v9114 < int32(32) {
		goto L1089
	} else {
		goto L1090
	}
L1087:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+4)) = v9228 + v9226
	*(*int32)(unsafe.Add(mBase, uint32(v268))) = v9225<<(uint(v9228)%32) | v9227
	goto L1085
L1088:
	;
	v9145 = *(*int32)(unsafe.Add(mBase, uint32(v268)+16))
	v9146 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v9150 = v9142
	v9151 = v9143
	v9152 = v9146
	v9153 = v9145
	goto L1092
L1089:
	;
	if v9127 < int32(16) {
		v9225 = v9112
		v9226 = v9114
		v9227 = v9126
		v9228 = v9127
		goto L1087
	} else {
		goto L1091
	}
L1090:
	;
	v9131 = int32(32)
	v9132 = v9131 - v9127
	v9140 = int32(base.Ui32(v9112) >> (uint(v9132) % 32))
	v9141 = v9114 - v9132
	v9142 = v9112<<(uint(v9127)%32) | v9126
	v9143 = v9131
	goto L1088
L1091:
	;
	v9140 = v9112
	v9141 = v9114
	v9142 = v9126
	v9143 = v9127
	goto L1088
L1092:
	;
	if base.Ui32(v9152+int32(2)) <= base.Ui32(v9153) {
		v9208 = v9152
		v9209 = v9153
		goto L1094
	} else {
		goto L1095
	}
L1093:
	;
	v9225 = v9140
	v9226 = v9141
	v9227 = v9221
	v9228 = v9219
	goto L1087
L1094:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v9208))) = uint16(v9150)
	v9216 = v9208 + int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v9216
	v9219 = v9151 + int32(-16)
	v9221 = int32(base.Ui32(v9150) >> (uint(int32(16)) % 32))
	if int32(31) < v9151 {
		v9150 = v9221
		v9151 = v9219
		v9152 = v9216
		v9153 = v9209
		goto L1092
	} else {
		goto L1109
	}
L1095:
	;
	v9162 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v9163 = v9153 - v9162
	v9166 = base.I64_extend_i32_s(v9163) + int64(32768)
	if base.Ui64(int64(4294967295)) < base.Ui64(v9166) {
		v9190 = v9162
		goto L1097
	} else {
		goto L1098
	}
L1096:
	;
	if v9152 == v9162 {
		goto L1107
	} else {
		goto L1108
	}
L1097:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v268)+20)) = int32(1)
	*(*int32)(unsafe.Add(mBase, uint32(v268)+12)) = v9190
	goto L1084
L1098:
	;
	v9169 = v9152 - v9162
	v9171 = v9166 + base.I64_extend_i32_u(v9169)
	if base.Ui64(int64(4294967295)) < base.Ui64(v9171) {
		v9190 = v9162
		goto L1097
	} else {
		goto L1099
	}
L1099:
	;
	v9174 = base.I32_wrap_i64(v9171)
	if v9153 == v9162 {
		goto L1100
	} else {
		goto L1101
	}
L1100:
	;
	v9181 = int32(base.Ui32(v9163*int32(3)) >> (uint(int32(1)) % 32))
	if base.Ui32(v9174) < base.Ui32(v9181) {
		goto L1103
	} else {
		goto L1104
	}
L1101:
	;
	if base.Ui32(v9174) <= base.Ui32(v9163) {
		v9208 = v9152
		v9209 = v9153
		goto L1094
	} else {
		goto L1102
	}
L1102:
	;
	goto L1100
L1103:
	;
	v9183 = v9181
	goto L1105
L1104:
	;
	v9183 = v9174
	goto L1105
L1105:
	;
	v9187 = v9183&int32(-1024) + int32(1024)
	v9188 = F_WebPSafeMalloc(m, int64(1), v9187)
	mBase = m.M
	if v9188 != 0 {
		goto L1096
	} else {
		goto L1106
	}
L1106:
	;
	v9189 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v9190 = v9189
	goto L1097
L1107:
	;
	v9201 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	F_WebPSafeFree(m, v9201)
	mBase = m.M
	v9203 = v9188 + v9187
	*(*int32)(unsafe.Add(mBase, uint32(v268)+16)) = v9203
	*(*int32)(unsafe.Add(mBase, uint32(v268)+8)) = v9188
	v9208 = v9188 + v9169
	v9209 = v9203
	goto L1094
L1108:
	;
	v9199 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v9200 = F_memcpy(m, v9188, v9199, v9169)
	mBase = m.M
	goto L1107
L1109:
	;
	goto L1093
L1110:
	;
	if v9285 != 0 {
		v9406 = v9285
		v9422 = v7866
		goto L949
	} else {
		goto L1116
	}
L1111:
	;
	goto L1110
L1112:
	;
	v9285 = F_malloc(m, base.I32_wrap_i64(int64(0))*int32(2))
	mBase = m.M
	goto L1111
L1116:
	;
	v9294 = v9112
	goto L950
L1117:
	;
	v9945 = v265
	v9947 = v267
	v9948 = v268
	v9949 = int32(0)
	v9950 = v9294
	v9951 = v271
	v9952 = v272
	v9953 = v273
	v9954 = v274
	v9955 = v275
	v9956 = v276
	v9957 = v277
	v9958 = v278
	v9960 = v280
	v9963 = v283
	v9964 = v284
	v9965 = v285
	v9966 = v286
	v9967 = v287
	v9968 = v288
	v9969 = v289
	v9970 = v290
	v9971 = v291
	v9972 = v292
	v9973 = v293
	v9974 = v294
	v9975 = v295
	v9976 = v296
	v9977 = v297
	v9978 = v298
	v9979 = v299
	v9980 = v300
	v9981 = v7245
	v9982 = v7246
	v9984 = v7295
	v9986 = v6912
	v9988 = v7392
	v9989 = v6945
	goto L776
L1118:
	;
	goto L1117
L1119:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L1118
L1120:
	;
	v9432 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	v9437 = *(*int32)(unsafe.Add(mBase, uint32(v268)+12))
	v9438 = *(*int32)(unsafe.Add(mBase, uint32(v268)+8))
	v9439 = v9437 - v9438
	v9440 = (v9432+int32(7))>>(uint(int32(3))%32) + v9439
	if base.Ui32(v7268) <= base.Ui32(v9440) {
		v9489 = v7245
		v9490 = v7246
		v9491 = v7268
		goto L1122
	} else {
		goto L1123
	}
L1121:
	;
	v9878 = v265
	v9880 = v267
	v9881 = v268
	v9883 = int32(0)
	v9884 = v271
	v9885 = v272
	v9886 = v273
	v9887 = v274
	v9888 = v275
	v9889 = v276
	v9890 = v277
	v9891 = v278
	v9893 = v280
	v9896 = v283
	v9897 = v284
	v9898 = v285
	v9899 = v286
	v9900 = v287
	v9901 = v288
	v9902 = v289
	v9903 = v290
	v9904 = v291
	v9905 = v292
	v9906 = v293
	v9907 = v294
	v9908 = v295
	v9909 = v296
	v9910 = v297
	v9911 = v298
	v9912 = v299
	v9913 = v300
	v9914 = v7245
	v9915 = v7246
	v9917 = v9406
	v9919 = v6912
	v9921 = v7392
	v9922 = v6945
	goto L777
L1122:
	;
	F_free(m, v9406)
	mBase = m.M
	goto L1125
L1123:
	;
	v9447 = (v9428+int32(7))>>(uint(int32(3))%32) + (v9427 - v9426)
	*(*int32)(unsafe.Add(mBase, uint32(v274)+48)) = v7299
	v9450 = *(*int32)(unsafe.Add(mBase, uint32(v268)+4))
	v9452 = v267 + int32(1088)
	v9461 = *(*int64)(unsafe.Add(mBase, uint32(v268)))
	v9462 = *(*int64)(unsafe.Add(mBase, uint32(v9452)))
	*(*int64)(unsafe.Add(mBase, uint32(v268))) = v9462
	v9465 = v268 + int32(16)
	v9466 = *(*int64)(unsafe.Add(mBase, uint32(v9465)))
	v9468 = v267 + int32(1104)
	v9469 = *(*int64)(unsafe.Add(mBase, uint32(v9468)))
	*(*int64)(unsafe.Add(mBase, uint32(v9465))) = v9469
	v9472 = v268 + int32(8)
	v9473 = *(*int64)(unsafe.Add(mBase, uint32(v9472)))
	v9475 = v267 + int32(1096)
	v9476 = *(*int64)(unsafe.Add(mBase, uint32(v9475)))
	*(*int64)(unsafe.Add(mBase, uint32(v9472))) = v9476
	*(*int64)(unsafe.Add(mBase, uint32(v9452))) = v9461
	*(*int64)(unsafe.Add(mBase, uint32(v9468))) = v9466
	*(*int64)(unsafe.Add(mBase, uint32(v9475))) = v9473
	goto L1124
L1124:
	;
	v9489 = v9439 - v9447 + (v9450+int32(7))>>(uint(int32(3))%32)
	v9490 = v9447 - v284
	v9491 = v9440
	goto L1122
L1125:
	;
	v9493 = *(*int32)(unsafe.Add(mBase, uint32(v7392)+8))
	F_free(m, v9493)
	mBase = m.M
	goto L1126
L1126:
	;
	F_free(m, v7392)
	mBase = m.M
	goto L1127
L1127:
	;
	v9496 = *(*int32)(unsafe.Add(mBase, uint32(v7193)+4))
	v9497 = int32(0)
	if v7272&base.B2i32(v9496 != v9497) != 0 {
		v7225 = int32(1)
		v7245 = v9489
		v7246 = v9490
		v7268 = v9491
		v7271 = v9422
		v7272 = v9497
		goto L830
	} else {
		goto L1128
	}
L1128:
	;
	goto L831
L1129:
	;
	goto L827
L1130:
	;
	v9671 = F_WebPReportProgress(m, v275, v6889+v6733, v267+int32(60))
	mBase = m.M
	v9712 = v9608
	v9713 = v9609
	goto L799
L1131:
	;
	v9712 = v301
	v9713 = v302
	goto L799
L1132:
	;
	goto L1131
L1133:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v275)+92)) = int32(1)
	goto L1132
L1134:
	;
	F_free(m, v9986)
	mBase = m.M
	goto L1135
L1135:
	;
	F_WebPSafeFree(m, v9949)
	mBase = m.M
	goto L1136
L1136:
	;
	F_WebPSafeFree(m, v9950)
	mBase = m.M
	goto L1137
L1137:
	;
	v10016 = v9947 + int32(64)
	v10017 = *(*int32)(unsafe.Add(mBase, uint32(v10016)))
	F_WebPSafeFree(m, v10017)
	mBase = m.M
	*(*int64)(unsafe.Add(mBase, uint32(v10016))) = int64(0)
	goto L1138
L1138:
	;
	if v9988 == int32(0) {
		goto L1139
	} else {
		goto L1140
	}
L1139:
	;
	F_free(m, v9989)
	mBase = m.M
	goto L1143
L1140:
	;
	v10023 = *(*int32)(unsafe.Add(mBase, uint32(v9988)+8))
	F_free(m, v10023)
	mBase = m.M
	goto L1141
L1141:
	;
	F_free(m, v9988)
	mBase = m.M
	goto L1142
L1142:
	;
	goto L1139
L1143:
	;
	v10028 = v9947 + int32(1088)
	if v10028 == int32(0) {
		goto L1145
	} else {
		goto L1146
	}
L1144:
	;
	v10045 = *(*int32)(unsafe.Add(mBase, uint32(v9955)+92))
	if v10045 != 0 {
		v10242 = v9945
		v10244 = v9947
		goto L4
	} else {
		goto L1147
	}
L1145:
	;
	goto L1144
L1146:
	;
	v10033 = v9947 + int32(1096)
	v10034 = *(*int32)(unsafe.Add(mBase, uint32(v10033)))
	F_WebPSafeFree(m, v10034)
	mBase = m.M
	v10038 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v9947+int32(1104)))) = v10038
	*(*int64)(unsafe.Add(mBase, uint32(v10033))) = v10038
	*(*int64)(unsafe.Add(mBase, uint32(v10028))) = v10038
	goto L1145
L1147:
	;
	v10046 = *(*int32)(unsafe.Add(mBase, uint32(v9948)+4))
	v10051 = *(*int32)(unsafe.Add(mBase, uint32(v9948)+12))
	v10052 = *(*int32)(unsafe.Add(mBase, uint32(v9948)+8))
	v10054 = (v10046+int32(7))>>(uint(int32(3))%32) + (v10051 - v10052)
	if base.Ui32(v9979) <= base.Ui32(v10054) {
		v10128 = v9979
		goto L1148
	} else {
		goto L1149
	}
L1148:
	;
	if v9953 < int32(2) {
		goto L1159
	} else {
		goto L1160
	}
L1149:
	;
	v10056 = int32(8)
	v10057 = v9947 + v10056
	v10066 = *(*int64)(unsafe.Add(mBase, uint32(v9948)))
	v10067 = *(*int64)(unsafe.Add(mBase, uint32(v10057)))
	*(*int64)(unsafe.Add(mBase, uint32(v9948))) = v10067
	v10069 = int32(16)
	v10070 = v9948 + v10069
	v10071 = *(*int64)(unsafe.Add(mBase, uint32(v10070)))
	v10073 = v9947 + int32(24)
	v10074 = *(*int64)(unsafe.Add(mBase, uint32(v10073)))
	*(*int64)(unsafe.Add(mBase, uint32(v10070))) = v10074
	v10077 = v9948 + v10056
	v10078 = *(*int64)(unsafe.Add(mBase, uint32(v10077)))
	v10080 = v9947 + v10069
	v10081 = *(*int64)(unsafe.Add(mBase, uint32(v10080)))
	*(*int64)(unsafe.Add(mBase, uint32(v10077))) = v10081
	*(*int64)(unsafe.Add(mBase, uint32(v10057))) = v10066
	*(*int64)(unsafe.Add(mBase, uint32(v10073))) = v10071
	*(*int64)(unsafe.Add(mBase, uint32(v10080))) = v10078
	goto L1150
L1150:
	;
	if v9951 == int32(0) {
		goto L1151
	} else {
		goto L1152
	}
L1151:
	;
	v10128 = v10054
	goto L1148
L1152:
	;
	v10088 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+60))
	v10089 = int32(0)
	v10090 = base.B2i32(v10088 != v10089)
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+148)) = v10090
	v10092 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+52))
	if v10092 == v10089 {
		v10098 = v10090
		goto L1153
	} else {
		goto L1154
	}
L1153:
	;
	v10099 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+56))
	if v10099 == int32(0) {
		v10105 = v10098
		goto L1155
	} else {
		goto L1156
	}
L1154:
	;
	v10096 = v10090 | int32(2)
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+148)) = v10096
	v10098 = v10096
	goto L1153
L1155:
	;
	v10106 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+64))
	if v10106 == int32(0) {
		goto L1157
	} else {
		goto L1158
	}
L1156:
	;
	v10103 = v10098 | int32(4)
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+148)) = v10103
	v10105 = v10103
	goto L1155
L1157:
	;
	v10112 = *(*int32)(unsafe.Add(mBase, uint32(v9947)+4))
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+156)) = v10112
	v10114 = *(*int32)(unsafe.Add(mBase, uint32(v9947)))
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+180)) = v10114
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+176)) = v9981
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+172)) = v9982
	v10118 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+36))
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+152)) = v10118
	v10120 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+48))
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+160)) = v10120
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+168)) = v10054 - v9964
	v10124 = *(*int32)(unsafe.Add(mBase, uint32(v9954)+68))
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+164)) = v10124
	goto L1151
L1158:
	;
	*(*int32)(unsafe.Add(mBase, uint32(v9951)+148)) = v10105 | int32(8)
	goto L1157
L1159:
	;
	v10144 = v9980 + int32(1)
	if v10144 != v9953 {
		v265 = v9945
		v267 = v9947
		v268 = v9948
		v271 = v9951
		v272 = v9952
		v273 = v9953
		v274 = v9954
		v275 = v9955
		v276 = v9956
		v277 = v9957
		v278 = v9958
		v280 = v9960
		v283 = v9963
		v284 = v9964
		v285 = v9965
		v286 = v9966
		v287 = v9967
		v288 = v9968
		v289 = v9969
		v290 = v9970
		v291 = v9971
		v292 = v9972
		v293 = v9973
		v294 = v9974
		v295 = v9975
		v296 = v9976
		v297 = v9977
		v298 = v9978
		v299 = v10128
		v300 = v10144
		v301 = v9981
		v302 = v9982
		goto L33
	} else {
		goto L1162
	}
L1160:
	;
	v10132 = v9947 + int32(32)
	v10133 = *(*int64)(unsafe.Add(mBase, uint32(v10132)))
	*(*int64)(unsafe.Add(mBase, uint32(v9948))) = v10133
	v10135 = *(*int32)(unsafe.Add(mBase, uint32(v10132)+20))
	*(*int32)(unsafe.Add(mBase, uint32(v9948)+20)) = v10135
	v10137 = *(*int32)(unsafe.Add(mBase, uint32(v9948)+8))
	v10138 = *(*int32)(unsafe.Add(mBase, uint32(v10132)+12))
	v10139 = *(*int32)(unsafe.Add(mBase, uint32(v10132)+8))
	*(*int32)(unsafe.Add(mBase, uint32(v9948)+12)) = v10137 + (v10138 - v10139)
	goto L1161
L1161:
	;
	goto L1159
L1162:
	;
	goto L34
L1163:
	;
	v10242 = v10146
	v10244 = v10148
	goto L4
L1164:
	;
	v10326 = *(*int32)(unsafe.Add(mBase, uint32(v10242)+4))
	v10327 = *(*int32)(unsafe.Add(mBase, uint32(v10326)+92))
	m.G0 = v10244 + int32(_a_F_EncodeStreamHook_0)
	return base.B2i32(v10327 == int32(0))
L1165:
	;
	goto L1164
L1166:
	;
	v10314 = v10244 + int32(16)
	v10315 = *(*int32)(unsafe.Add(mBase, uint32(v10314)))
	F_WebPSafeFree(m, v10315)
	mBase = m.M
	v10319 = int64(0)
	*(*int64)(unsafe.Add(mBase, uint32(v10244+int32(24)))) = v10319
	*(*int64)(unsafe.Add(mBase, uint32(v10314))) = v10319
	*(*int64)(unsafe.Add(mBase, uint32(v10309))) = v10319
	goto L1165
}

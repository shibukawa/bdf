//go:build !bdf_noconv && goexperiment.simd && go1.27 && !go1.28 && (amd64 || arm64)

package p0

import (
	base "github.com/shibukawa/bdf/imgconv/internal/webpwsimd/base"
	"unsafe"
)

func F_SharpYuvConvertWithOptions(m *base.Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32, l15 int32) int32 {
	mBase := m.M
	_ = mBase
	var v17 int32
	_ = v17
	var v154 int32
	_ = v154
	var v155 int32
	_ = v155
	var v169 int32
	_ = v169
	var v172 int32
	_ = v172
	var v173 int32
	_ = v173
	var v174 int32
	_ = v174
	var v176 int32
	_ = v176
	var v181 int32
	_ = v181
	var v182 int32
	_ = v182
	var v186 int32
	_ = v186
	var v188 int32
	_ = v188
	var v191 int32
	_ = v191
	var v193 int32
	_ = v193
	var v198 int32
	_ = v198
	var v199 int32
	_ = v199
	var v200 int32
	_ = v200
	var v203 int32
	_ = v203
	var v204 int32
	_ = v204
	var v207 int32
	_ = v207
	var v208 int32
	_ = v208
	var v211 int32
	_ = v211
	var v212 int32
	_ = v212
	var v215 int32
	_ = v215
	var v216 int32
	_ = v216
	var v219 int32
	_ = v219
	var v220 int32
	_ = v220
	var v223 int32
	_ = v223
	var v224 int32
	_ = v224
	var v227 int32
	_ = v227
	var v228 int32
	_ = v228
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
	var v239 int32
	_ = v239
	var v240 int32
	_ = v240
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
	var v257 int32
	_ = v257
	var v258 int32
	_ = v258
	var v261 int32
	_ = v261
	var v263 int32
	_ = v263
	var v265 int32
	_ = v265
	var v268 int64
	_ = v268
	var v272 int32
	_ = v272
	var v273 int32
	_ = v273
	var v274 int32
	_ = v274
	var v276 int64
	_ = v276
	var v280 int64
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
	var v289 int32
	_ = v289
	var v290 int32
	_ = v290
	var v292 int32
	_ = v292
	var v293 int32
	_ = v293
	var v296 int64
	_ = v296
	var v300 int32
	_ = v300
	var v301 int32
	_ = v301
	var v305 int32
	_ = v305
	var v306 int64
	_ = v306
	var v309 int64
	_ = v309
	var v312 int32
	_ = v312
	var v313 int32
	_ = v313
	var v314 int32
	_ = v314
	var v316 int32
	_ = v316
	var v317 int32
	_ = v317
	var v320 int64
	_ = v320
	var v324 int32
	_ = v324
	var v325 int32
	_ = v325
	var v330 float64
	_ = v330
	var v338 int64
	_ = v338
	var v340 int64
	_ = v340
	var v341 int32
	_ = v341
	var v361 int32
	_ = v361
	var v362 int32
	_ = v362
	var v363 int32
	_ = v363
	var v366 int32
	_ = v366
	var v372 int32
	_ = v372
	var v376 int32
	_ = v376
	var v379 int32
	_ = v379
	var v384 int32
	_ = v384
	var v386 int32
	_ = v386
	var v388 int32
	_ = v388
	var v389 int32
	_ = v389
	var v391 int32
	_ = v391
	var v392 int32
	_ = v392
	var v399 int32
	_ = v399
	var v402 int32
	_ = v402
	var v403 int32
	_ = v403
	var v408 int32
	_ = v408
	var v412 int32
	_ = v412
	var v414 int32
	_ = v414
	var v416 int32
	_ = v416
	var v417 int32
	_ = v417
	var v418 int32
	_ = v418
	var v419 int32
	_ = v419
	var v474 int32
	_ = v474
	var v484 int32
	_ = v484
	var v485 int32
	_ = v485
	var v491 int32
	_ = v491
	var v492 int32
	_ = v492
	var v493 int32
	_ = v493
	var v494 int32
	_ = v494
	var v548 int32
	_ = v548
	var v551 int32
	_ = v551
	var v552 int32
	_ = v552
	var v554 int32
	_ = v554
	var v556 int32
	_ = v556
	var v558 int32
	_ = v558
	var v561 int32
	_ = v561
	var v564 int32
	_ = v564
	var v569 int32
	_ = v569
	var v573 int32
	_ = v573
	var v575 int32
	_ = v575
	var v576 int32
	_ = v576
	var v578 int32
	_ = v578
	var v581 int32
	_ = v581
	var v582 int32
	_ = v582
	var v586 int32
	_ = v586
	var v590 int32
	_ = v590
	var v594 int32
	_ = v594
	var v624 int32
	_ = v624
	var v628 int32
	_ = v628
	var v629 int32
	_ = v629
	var v638 int32
	_ = v638
	var v639 int32
	_ = v639
	var v640 int32
	_ = v640
	var v646 int32
	_ = v646
	var v647 base.V128
	_ = v647
	var v648 int32
	_ = v648
	var v649 base.V128
	_ = v649
	var v654 base.V128
	_ = v654
	var v656 base.V128
	_ = v656
	var v661 base.V128
	_ = v661
	var v663 base.V128
	_ = v663
	var v666 int32
	_ = v666
	var v675 int32
	_ = v675
	var v685 int32
	_ = v685
	var v700 int32
	_ = v700
	var v710 int32
	_ = v710
	var v711 int32
	_ = v711
	var v722 int32
	_ = v722
	var v729 int32
	_ = v729
	var v730 int32
	_ = v730
	var v731 int32
	_ = v731
	var v735 int32
	_ = v735
	var v737 int32
	_ = v737
	var v741 int32
	_ = v741
	var v743 int32
	_ = v743
	var v749 int32
	_ = v749
	var v752 int32
	_ = v752
	var v754 int32
	_ = v754
	var v758 int32
	_ = v758
	var v764 int32
	_ = v764
	var v767 int32
	_ = v767
	var v770 int32
	_ = v770
	var v773 int32
	_ = v773
	var v774 int32
	_ = v774
	var v777 int32
	_ = v777
	var v782 int32
	_ = v782
	var v785 int32
	_ = v785
	var v787 base.V128
	_ = v787
	var v789 base.V128
	_ = v789
	var v818 int32
	_ = v818
	var v822 int32
	_ = v822
	var v823 int32
	_ = v823
	var v832 int32
	_ = v832
	var v833 int32
	_ = v833
	var v835 int32
	_ = v835
	var v840 int32
	_ = v840
	var v841 base.V128
	_ = v841
	var v842 base.V128
	_ = v842
	var v846 base.V128
	_ = v846
	var v850 base.V128
	_ = v850
	var v855 base.V128
	_ = v855
	var v864 base.V128
	_ = v864
	var v869 base.V128
	_ = v869
	var v878 base.V128
	_ = v878
	var v881 int32
	_ = v881
	var v890 int32
	_ = v890
	var v892 int32
	_ = v892
	var v893 int32
	_ = v893
	var v897 int32
	_ = v897
	var v900 int32
	_ = v900
	var v901 int32
	_ = v901
	var v904 int32
	_ = v904
	var v909 int32
	_ = v909
	var v912 int32
	_ = v912
	var v914 base.V128
	_ = v914
	var v916 base.V128
	_ = v916
	var v945 int32
	_ = v945
	var v949 int32
	_ = v949
	var v950 int32
	_ = v950
	var v959 int32
	_ = v959
	var v960 int32
	_ = v960
	var v962 int32
	_ = v962
	var v967 int32
	_ = v967
	var v968 base.V128
	_ = v968
	var v971 base.V128
	_ = v971
	var v975 base.V128
	_ = v975
	var v980 base.V128
	_ = v980
	var v987 base.V128
	_ = v987
	var v992 base.V128
	_ = v992
	var v999 base.V128
	_ = v999
	var v1002 int32
	_ = v1002
	var v1011 int32
	_ = v1011
	var v1025 int32
	_ = v1025
	var v1034 int32
	_ = v1034
	var v1050 int32
	_ = v1050
	var v1051 int32
	_ = v1051
	var v1061 int32
	_ = v1061
	var v1069 int32
	_ = v1069
	var v1070 int32
	_ = v1070
	var v1074 int32
	_ = v1074
	var v1075 int32
	_ = v1075
	var v1079 int32
	_ = v1079
	var v1080 int32
	_ = v1080
	var v1086 int32
	_ = v1086
	var v1099 int32
	_ = v1099
	var v1108 int32
	_ = v1108
	var v1124 int32
	_ = v1124
	var v1125 int32
	_ = v1125
	var v1135 int32
	_ = v1135
	var v1143 int32
	_ = v1143
	var v1144 int32
	_ = v1144
	var v1148 int32
	_ = v1148
	var v1149 int32
	_ = v1149
	var v1153 int32
	_ = v1153
	var v1154 int32
	_ = v1154
	var v1160 int32
	_ = v1160
	var v1186 int32
	_ = v1186
	var v1188 int32
	_ = v1188
	var v1189 int32
	_ = v1189
	var v1191 int32
	_ = v1191
	var v1195 int32
	_ = v1195
	var v1198 int32
	_ = v1198
	var v1202 int32
	_ = v1202
	var v1205 int32
	_ = v1205
	var v1210 int32
	_ = v1210
	var v1211 int32
	_ = v1211
	var v1212 int32
	_ = v1212
	var v1228 int32
	_ = v1228
	var v1231 int32
	_ = v1231
	var v1232 int32
	_ = v1232
	var v1234 int32
	_ = v1234
	var v1236 int32
	_ = v1236
	var v1238 int32
	_ = v1238
	var v1241 int32
	_ = v1241
	var v1244 int32
	_ = v1244
	var v1249 int32
	_ = v1249
	var v1253 int32
	_ = v1253
	var v1255 int32
	_ = v1255
	var v1256 int32
	_ = v1256
	var v1258 int32
	_ = v1258
	var v1261 int32
	_ = v1261
	var v1262 int32
	_ = v1262
	var v1266 int32
	_ = v1266
	var v1270 int32
	_ = v1270
	var v1274 int32
	_ = v1274
	var v1304 int32
	_ = v1304
	var v1308 int32
	_ = v1308
	var v1309 int32
	_ = v1309
	var v1318 int32
	_ = v1318
	var v1319 int32
	_ = v1319
	var v1320 int32
	_ = v1320
	var v1326 int32
	_ = v1326
	var v1327 base.V128
	_ = v1327
	var v1328 int32
	_ = v1328
	var v1329 base.V128
	_ = v1329
	var v1334 base.V128
	_ = v1334
	var v1336 base.V128
	_ = v1336
	var v1341 base.V128
	_ = v1341
	var v1343 base.V128
	_ = v1343
	var v1346 int32
	_ = v1346
	var v1355 int32
	_ = v1355
	var v1365 int32
	_ = v1365
	var v1380 int32
	_ = v1380
	var v1390 int32
	_ = v1390
	var v1391 int32
	_ = v1391
	var v1402 int32
	_ = v1402
	var v1409 int32
	_ = v1409
	var v1410 int32
	_ = v1410
	var v1411 int32
	_ = v1411
	var v1415 int32
	_ = v1415
	var v1417 int32
	_ = v1417
	var v1421 int32
	_ = v1421
	var v1423 int32
	_ = v1423
	var v1429 int32
	_ = v1429
	var v1432 int32
	_ = v1432
	var v1434 int32
	_ = v1434
	var v1438 int32
	_ = v1438
	var v1444 int32
	_ = v1444
	var v1447 int32
	_ = v1447
	var v1450 int32
	_ = v1450
	var v1453 int32
	_ = v1453
	var v1454 int32
	_ = v1454
	var v1457 int32
	_ = v1457
	var v1462 int32
	_ = v1462
	var v1465 int32
	_ = v1465
	var v1467 base.V128
	_ = v1467
	var v1469 base.V128
	_ = v1469
	var v1498 int32
	_ = v1498
	var v1502 int32
	_ = v1502
	var v1503 int32
	_ = v1503
	var v1512 int32
	_ = v1512
	var v1513 int32
	_ = v1513
	var v1515 int32
	_ = v1515
	var v1520 int32
	_ = v1520
	var v1521 base.V128
	_ = v1521
	var v1522 base.V128
	_ = v1522
	var v1526 base.V128
	_ = v1526
	var v1530 base.V128
	_ = v1530
	var v1535 base.V128
	_ = v1535
	var v1544 base.V128
	_ = v1544
	var v1549 base.V128
	_ = v1549
	var v1558 base.V128
	_ = v1558
	var v1561 int32
	_ = v1561
	var v1570 int32
	_ = v1570
	var v1572 int32
	_ = v1572
	var v1573 int32
	_ = v1573
	var v1577 int32
	_ = v1577
	var v1580 int32
	_ = v1580
	var v1581 int32
	_ = v1581
	var v1584 int32
	_ = v1584
	var v1589 int32
	_ = v1589
	var v1592 int32
	_ = v1592
	var v1594 base.V128
	_ = v1594
	var v1596 base.V128
	_ = v1596
	var v1625 int32
	_ = v1625
	var v1629 int32
	_ = v1629
	var v1630 int32
	_ = v1630
	var v1639 int32
	_ = v1639
	var v1640 int32
	_ = v1640
	var v1642 int32
	_ = v1642
	var v1647 int32
	_ = v1647
	var v1648 base.V128
	_ = v1648
	var v1651 base.V128
	_ = v1651
	var v1655 base.V128
	_ = v1655
	var v1660 base.V128
	_ = v1660
	var v1667 base.V128
	_ = v1667
	var v1672 base.V128
	_ = v1672
	var v1679 base.V128
	_ = v1679
	var v1682 int32
	_ = v1682
	var v1691 int32
	_ = v1691
	var v1705 int32
	_ = v1705
	var v1714 int32
	_ = v1714
	var v1730 int32
	_ = v1730
	var v1731 int32
	_ = v1731
	var v1741 int32
	_ = v1741
	var v1749 int32
	_ = v1749
	var v1750 int32
	_ = v1750
	var v1754 int32
	_ = v1754
	var v1755 int32
	_ = v1755
	var v1759 int32
	_ = v1759
	var v1760 int32
	_ = v1760
	var v1766 int32
	_ = v1766
	var v1779 int32
	_ = v1779
	var v1788 int32
	_ = v1788
	var v1804 int32
	_ = v1804
	var v1805 int32
	_ = v1805
	var v1815 int32
	_ = v1815
	var v1823 int32
	_ = v1823
	var v1824 int32
	_ = v1824
	var v1828 int32
	_ = v1828
	var v1829 int32
	_ = v1829
	var v1833 int32
	_ = v1833
	var v1834 int32
	_ = v1834
	var v1840 int32
	_ = v1840
	var v1866 int32
	_ = v1866
	var v1868 int32
	_ = v1868
	var v1869 int32
	_ = v1869
	var v1871 int32
	_ = v1871
	var v1875 int32
	_ = v1875
	var v1878 int32
	_ = v1878
	var v1882 int32
	_ = v1882
	var v1885 int32
	_ = v1885
	var v1889 int32
	_ = v1889
	var v1907 int32
	_ = v1907
	var v1910 int32
	_ = v1910
	var v2010 int32
	_ = v2010
	var v2011 base.V128
	_ = v2011
	var v2013 base.V128
	_ = v2013
	var v2017 base.V128
	_ = v2017
	var v2019 base.V128
	_ = v2019
	var v2024 base.V128
	_ = v2024
	var v2026 base.V128
	_ = v2026
	var v2029 base.V128
	_ = v2029
	var v2031 int32
	_ = v2031
	var v2033 base.V128
	_ = v2033
	var v2050 base.V128
	_ = v2050
	var v2056 int32
	_ = v2056
	var v2076 int32
	_ = v2076
	var v2192 int32
	_ = v2192
	var v2195 int32
	_ = v2195
	var v2295 int64
	_ = v2295
	var v2299 int64
	_ = v2299
	var v2304 int64
	_ = v2304
	var v2311 int64
	_ = v2311
	var v2316 int32
	_ = v2316
	var v2449 int32
	_ = v2449
	var v2452 int32
	_ = v2452
	var v2463 int32
	_ = v2463
	var v2473 int32
	_ = v2473
	var v2519 int32
	_ = v2519
	var v2550 int32
	_ = v2550
	var v2551 base.V128
	_ = v2551
	var v2553 base.V128
	_ = v2553
	var v2556 base.V128
	_ = v2556
	var v2558 base.V128
	_ = v2558
	var v2562 base.V128
	_ = v2562
	var v2564 base.V128
	_ = v2564
	var v2567 base.V128
	_ = v2567
	var v2569 int32
	_ = v2569
	var v2571 base.V128
	_ = v2571
	var v2588 base.V128
	_ = v2588
	var v2600 int32
	_ = v2600
	var v2620 int32
	_ = v2620
	var v2736 int32
	_ = v2736
	var v2739 int32
	_ = v2739
	var v2839 int64
	_ = v2839
	var v2843 int64
	_ = v2843
	var v2848 int64
	_ = v2848
	var v2855 int64
	_ = v2855
	var v2860 int32
	_ = v2860
	var v2993 int32
	_ = v2993
	var v2996 int32
	_ = v2996
	var v3096 int32
	_ = v3096
	var v3097 int32
	_ = v3097
	var v3102 int32
	_ = v3102
	var v3103 int32
	_ = v3103
	var v3109 int32
	_ = v3109
	var v3110 int32
	_ = v3110
	var v3120 int32
	_ = v3120
	var v3125 int32
	_ = v3125
	var v3141 int32
	_ = v3141
	var v3144 int32
	_ = v3144
	var v3155 int32
	_ = v3155
	var v3165 int32
	_ = v3165
	var v3211 int32
	_ = v3211
	var v3242 int32
	_ = v3242
	var v3243 int32
	_ = v3243
	var v3247 int32
	_ = v3247
	var v3248 int32
	_ = v3248
	var v3253 int32
	_ = v3253
	var v3254 int32
	_ = v3254
	var v3264 int32
	_ = v3264
	var v3266 int32
	_ = v3266
	var v3275 int32
	_ = v3275
	var v3294 int32
	_ = v3294
	var v3295 int32
	_ = v3295
	var v3296 int32
	_ = v3296
	var v3299 int32
	_ = v3299
	var v3303 int32
	_ = v3303
	var v3306 int32
	_ = v3306
	var v3313 int32
	_ = v3313
	var v3314 int32
	_ = v3314
	var v3317 int32
	_ = v3317
	var v3330 int32
	_ = v3330
	var v3331 int32
	_ = v3331
	var v3333 int32
	_ = v3333
	var v3334 int32
	_ = v3334
	var v3335 int32
	_ = v3335
	var v3338 int32
	_ = v3338
	var v3339 int32
	_ = v3339
	var v3340 int32
	_ = v3340
	var v3341 int32
	_ = v3341
	var v3343 int32
	_ = v3343
	var v3345 int32
	_ = v3345
	var v3351 int32
	_ = v3351
	var v3352 int32
	_ = v3352
	var v3355 int32
	_ = v3355
	var v3356 int32
	_ = v3356
	var v3357 int32
	_ = v3357
	var v3360 int32
	_ = v3360
	var v3361 int32
	_ = v3361
	var v3362 int32
	_ = v3362
	var v3363 int32
	_ = v3363
	var v3365 int32
	_ = v3365
	var v3367 int32
	_ = v3367
	var v3373 int32
	_ = v3373
	var v3374 int32
	_ = v3374
	var v3377 int32
	_ = v3377
	var v3378 int32
	_ = v3378
	var v3379 int32
	_ = v3379
	var v3382 int32
	_ = v3382
	var v3390 int32
	_ = v3390
	var v3391 int32
	_ = v3391
	var v3392 int32
	_ = v3392
	var v3394 int32
	_ = v3394
	var v3396 int32
	_ = v3396
	var v3402 int32
	_ = v3402
	var v3411 int32
	_ = v3411
	var v3412 int32
	_ = v3412
	var v3415 int32
	_ = v3415
	var v3418 int32
	_ = v3418
	var v3425 int32
	_ = v3425
	var v3426 int32
	_ = v3426
	var v3435 int32
	_ = v3435
	var v3439 int32
	_ = v3439
	var v3499 int32
	_ = v3499
	var v3500 int32
	_ = v3500
	var v3501 int32
	_ = v3501
	var v3557 int32
	_ = v3557
	var v3558 int32
	_ = v3558
	var v3559 int32
	_ = v3559
	var v3560 int32
	_ = v3560
	var v3561 int32
	_ = v3561
	var v3562 int32
	_ = v3562
	var v3563 int32
	_ = v3563
	var v3564 int32
	_ = v3564
	var v3566 int32
	_ = v3566
	var v3570 int32
	_ = v3570
	var v3572 int32
	_ = v3572
	var v3574 int32
	_ = v3574
	var v3576 int32
	_ = v3576
	var v3580 int32
	_ = v3580
	var v3591 int32
	_ = v3591
	var v3593 int32
	_ = v3593
	var v3595 int32
	_ = v3595
	var v3597 int32
	_ = v3597
	var v3614 int32
	_ = v3614
	var v3630 int32
	_ = v3630
	var v3650 int64
	_ = v3650
	var v3768 int64
	_ = v3768
	var __phi3768 int64
	_ = __phi3768
	var v3772 int32
	_ = v3772
	var __phi3772 int32
	_ = __phi3772
	var v3790 int32
	_ = v3790
	var __phi3790 int32
	_ = __phi3790
	var v3796 int32
	_ = v3796
	var __phi3796 int32
	_ = __phi3796
	var v3800 int32
	_ = v3800
	var __phi3800 int32
	_ = __phi3800
	var v3804 int32
	_ = v3804
	var __phi3804 int32
	_ = __phi3804
	var v3805 int32
	_ = v3805
	var __phi3805 int32
	_ = __phi3805
	var v3807 int32
	_ = v3807
	var __phi3807 int32
	_ = __phi3807
	var v3849 int32
	_ = v3849
	var v3851 int32
	_ = v3851
	var v3854 int32
	_ = v3854
	var v3855 int32
	_ = v3855
	var v3856 int32
	_ = v3856
	var v3860 int32
	_ = v3860
	var v3861 int32
	_ = v3861
	var v3864 int32
	_ = v3864
	var v3866 int32
	_ = v3866
	var v3868 int32
	_ = v3868
	var v3871 int32
	_ = v3871
	var v3874 int32
	_ = v3874
	var v3875 int32
	_ = v3875
	var v3879 int32
	_ = v3879
	var v3880 int32
	_ = v3880
	var v3881 int32
	_ = v3881
	var v3884 int32
	_ = v3884
	var v3886 int32
	_ = v3886
	var v3888 int32
	_ = v3888
	var v3889 int32
	_ = v3889
	var v3890 int32
	_ = v3890
	var v3891 int32
	_ = v3891
	var v3894 int32
	_ = v3894
	var v3895 int32
	_ = v3895
	var v3897 int32
	_ = v3897
	var v3899 int32
	_ = v3899
	var v3903 int32
	_ = v3903
	var v3905 int32
	_ = v3905
	var v3909 int32
	_ = v3909
	var v3910 int32
	_ = v3910
	var v3911 int32
	_ = v3911
	var v3914 int32
	_ = v3914
	var v3916 int32
	_ = v3916
	var v3918 int32
	_ = v3918
	var v3920 int32
	_ = v3920
	var v3924 int32
	_ = v3924
	var v3925 int32
	_ = v3925
	var v3926 int32
	_ = v3926
	var v3929 int32
	_ = v3929
	var v3931 int32
	_ = v3931
	var v3933 int32
	_ = v3933
	var v3935 int32
	_ = v3935
	var v3936 int32
	_ = v3936
	var v3937 int32
	_ = v3937
	var v3940 int32
	_ = v3940
	var v3941 int32
	_ = v3941
	var v3942 int32
	_ = v3942
	var v3943 int32
	_ = v3943
	var v3947 int32
	_ = v3947
	var v3948 int32
	_ = v3948
	var v3951 int32
	_ = v3951
	var v3953 int32
	_ = v3953
	var v3955 int32
	_ = v3955
	var v3956 int32
	_ = v3956
	var v3957 int32
	_ = v3957
	var v3961 int32
	_ = v3961
	var v3962 int32
	_ = v3962
	var v3965 int32
	_ = v3965
	var v3967 int32
	_ = v3967
	var v3969 int32
	_ = v3969
	var v3971 int32
	_ = v3971
	var v3973 int32
	_ = v3973
	var v3975 int32
	_ = v3975
	var v3978 int32
	_ = v3978
	var v3979 int32
	_ = v3979
	var v3981 int32
	_ = v3981
	var v3985 int32
	_ = v3985
	var v3986 int32
	_ = v3986
	var v3989 int32
	_ = v3989
	var v3991 int32
	_ = v3991
	var v3993 int32
	_ = v3993
	var v3995 int32
	_ = v3995
	var v3999 int32
	_ = v3999
	var v4000 int32
	_ = v4000
	var v4003 int32
	_ = v4003
	var v4005 int32
	_ = v4005
	var v4007 int32
	_ = v4007
	var v4008 int32
	_ = v4008
	var v4009 int32
	_ = v4009
	var v4012 int32
	_ = v4012
	var v4013 int32
	_ = v4013
	var v4014 int32
	_ = v4014
	var v4015 int32
	_ = v4015
	var v4019 int32
	_ = v4019
	var v4020 int32
	_ = v4020
	var v4023 int32
	_ = v4023
	var v4025 int32
	_ = v4025
	var v4027 int32
	_ = v4027
	var v4028 int32
	_ = v4028
	var v4029 int32
	_ = v4029
	var v4033 int32
	_ = v4033
	var v4034 int32
	_ = v4034
	var v4037 int32
	_ = v4037
	var v4039 int32
	_ = v4039
	var v4041 int32
	_ = v4041
	var v4043 int32
	_ = v4043
	var v4045 int32
	_ = v4045
	var v4047 int32
	_ = v4047
	var v4050 int32
	_ = v4050
	var v4051 int32
	_ = v4051
	var v4053 int32
	_ = v4053
	var v4057 int32
	_ = v4057
	var v4058 int32
	_ = v4058
	var v4061 int32
	_ = v4061
	var v4063 int32
	_ = v4063
	var v4065 int32
	_ = v4065
	var v4067 int32
	_ = v4067
	var v4071 int32
	_ = v4071
	var v4072 int32
	_ = v4072
	var v4075 int32
	_ = v4075
	var v4077 int32
	_ = v4077
	var v4094 int32
	_ = v4094
	var v4097 int32
	_ = v4097
	var v4197 int32
	_ = v4197
	var v4198 int32
	_ = v4198
	var v4203 int32
	_ = v4203
	var v4204 int32
	_ = v4204
	var v4210 int32
	_ = v4210
	var v4211 int32
	_ = v4211
	var v4221 int32
	_ = v4221
	var v4226 int32
	_ = v4226
	var v4242 int32
	_ = v4242
	var v4245 int32
	_ = v4245
	var v4256 int32
	_ = v4256
	var v4266 int32
	_ = v4266
	var v4343 int32
	_ = v4343
	var v4344 int32
	_ = v4344
	var v4348 int32
	_ = v4348
	var v4349 int32
	_ = v4349
	var v4355 int32
	_ = v4355
	var v4356 int32
	_ = v4356
	var v4366 int32
	_ = v4366
	var v4368 int32
	_ = v4368
	var v4375 int32
	_ = v4375
	var v4394 int32
	_ = v4394
	var v4395 int32
	_ = v4395
	var v4396 int32
	_ = v4396
	var v4399 int32
	_ = v4399
	var v4403 int32
	_ = v4403
	var v4406 int32
	_ = v4406
	var v4413 int32
	_ = v4413
	var v4414 int32
	_ = v4414
	var v4417 int32
	_ = v4417
	var v4430 int32
	_ = v4430
	var v4431 int32
	_ = v4431
	var v4433 int32
	_ = v4433
	var v4434 int32
	_ = v4434
	var v4435 int32
	_ = v4435
	var v4438 int32
	_ = v4438
	var v4439 int32
	_ = v4439
	var v4440 int32
	_ = v4440
	var v4441 int32
	_ = v4441
	var v4443 int32
	_ = v4443
	var v4445 int32
	_ = v4445
	var v4451 int32
	_ = v4451
	var v4452 int32
	_ = v4452
	var v4455 int32
	_ = v4455
	var v4456 int32
	_ = v4456
	var v4457 int32
	_ = v4457
	var v4460 int32
	_ = v4460
	var v4461 int32
	_ = v4461
	var v4462 int32
	_ = v4462
	var v4463 int32
	_ = v4463
	var v4465 int32
	_ = v4465
	var v4467 int32
	_ = v4467
	var v4473 int32
	_ = v4473
	var v4474 int32
	_ = v4474
	var v4477 int32
	_ = v4477
	var v4478 int32
	_ = v4478
	var v4479 int32
	_ = v4479
	var v4482 int32
	_ = v4482
	var v4490 int32
	_ = v4490
	var v4491 int32
	_ = v4491
	var v4492 int32
	_ = v4492
	var v4494 int32
	_ = v4494
	var v4496 int32
	_ = v4496
	var v4502 int32
	_ = v4502
	var v4511 int32
	_ = v4511
	var v4512 int32
	_ = v4512
	var v4515 int32
	_ = v4515
	var v4518 int32
	_ = v4518
	var v4525 int32
	_ = v4525
	var v4526 int32
	_ = v4526
	var v4527 int32
	_ = v4527
	var v4528 int32
	_ = v4528
	var v4529 int64
	_ = v4529
	var v4530 int32
	_ = v4530
	var v4532 int64
	_ = v4532
	var v4533 int32
	_ = v4533
	var v4534 int32
	_ = v4534
	var v4537 int32
	_ = v4537
	var v4542 int32
	_ = v4542
	var v4554 int32
	_ = v4554
	var v4557 int32
	_ = v4557
	var v4558 int32
	_ = v4558
	var v4561 int32
	_ = v4561
	var v4563 int32
	_ = v4563
	var v4567 int32
	_ = v4567
	var v4570 int32
	_ = v4570
	var v4573 int32
	_ = v4573
	var v4577 int32
	_ = v4577
	var v4592 int32
	_ = v4592
	var v4601 int32
	_ = v4601
	var v4661 int32
	_ = v4661
	var v4671 int32
	_ = v4671
	var v4722 int32
	_ = v4722
	var v4724 int32
	_ = v4724
	var v4824 int32
	_ = v4824
	var v4825 base.V128
	_ = v4825
	var v4828 base.V128
	_ = v4828
	var v4829 base.V128
	_ = v4829
	var v4835 base.V128
	_ = v4835
	var v4842 base.V128
	_ = v4842
	var v4848 base.V128
	_ = v4848
	var v4852 base.V128
	_ = v4852
	var v4855 base.V128
	_ = v4855
	var v4868 int32
	_ = v4868
	var v4887 int32
	_ = v4887
	var v5006 int32
	_ = v5006
	var v5008 int32
	_ = v5008
	var v5110 int32
	_ = v5110
	var v5111 int32
	_ = v5111
	var v5115 int32
	_ = v5115
	var v5116 int32
	_ = v5116
	var v5121 int32
	_ = v5121
	var v5129 int32
	_ = v5129
	var v5134 int32
	_ = v5134
	var v5144 int32
	_ = v5144
	var v5149 int32
	_ = v5149
	var v5273 int32
	_ = v5273
	var v5275 int32
	_ = v5275
	var v5279 int32
	_ = v5279
	var v5282 int32
	_ = v5282
	var v5287 int32
	_ = v5287
	var v5288 base.V128
	_ = v5288
	var v5289 base.V128
	_ = v5289
	var v5290 base.V128
	_ = v5290
	var v5291 base.V128
	_ = v5291
	var v5292 base.V128
	_ = v5292
	var v5293 base.V128
	_ = v5293
	var v5294 base.V128
	_ = v5294
	var v5304 int32
	_ = v5304
	var v5314 int32
	_ = v5314
	var v5368 int32
	_ = v5368
	var v5373 int32
	_ = v5373
	var v5433 int32
	_ = v5433
	var v5436 int32
	_ = v5436
	var v5499 base.V128
	_ = v5499
	var v5500 base.V128
	_ = v5500
	var v5535 int32
	_ = v5535
	var v5536 base.V128
	_ = v5536
	var v5537 int32
	_ = v5537
	var v5542 int32
	_ = v5542
	var v5552 int32
	_ = v5552
	var v5558 base.V128
	_ = v5558
	var v5561 base.V128
	_ = v5561
	var v5564 base.V128
	_ = v5564
	var v5567 base.V128
	_ = v5567
	var v5571 base.V128
	_ = v5571
	var v5572 base.V128
	_ = v5572
	var v5576 base.V128
	_ = v5576
	var v5598 base.V128
	_ = v5598
	var v5601 base.V128
	_ = v5601
	var v5604 base.V128
	_ = v5604
	var v5607 base.V128
	_ = v5607
	var v5612 base.V128
	_ = v5612
	var v5634 base.V128
	_ = v5634
	var v5637 base.V128
	_ = v5637
	var v5640 base.V128
	_ = v5640
	var v5643 base.V128
	_ = v5643
	var v5648 base.V128
	_ = v5648
	var v5649 int32
	_ = v5649
	var v5655 base.V128
	_ = v5655
	var v5658 base.V128
	_ = v5658
	var v5680 base.V128
	_ = v5680
	var v5683 base.V128
	_ = v5683
	var v5686 base.V128
	_ = v5686
	var v5689 base.V128
	_ = v5689
	var v5693 base.V128
	_ = v5693
	var v5697 base.V128
	_ = v5697
	var v5719 base.V128
	_ = v5719
	var v5722 base.V128
	_ = v5722
	var v5725 base.V128
	_ = v5725
	var v5728 base.V128
	_ = v5728
	var v5733 base.V128
	_ = v5733
	var v5755 base.V128
	_ = v5755
	var v5758 base.V128
	_ = v5758
	var v5761 base.V128
	_ = v5761
	var v5764 base.V128
	_ = v5764
	var v5769 base.V128
	_ = v5769
	var v5778 base.V128
	_ = v5778
	var v5779 base.V128
	_ = v5779
	var v5787 base.V128
	_ = v5787
	var v5790 base.V128
	_ = v5790
	var v5796 int32
	_ = v5796
	var v5813 int32
	_ = v5813
	var v5931 int32
	_ = v5931
	var v5934 int32
	_ = v5934
	var v6037 int32
	_ = v6037
	var v6039 int32
	_ = v6039
	var v6043 int32
	_ = v6043
	var v6044 int32
	_ = v6044
	var v6049 int32
	_ = v6049
	var v6057 int32
	_ = v6057
	var v6061 int32
	_ = v6061
	var v6064 int32
	_ = v6064
	var v6067 int32
	_ = v6067
	var v6072 int32
	_ = v6072
	var v6196 int32
	_ = v6196
	var v6198 int32
	_ = v6198
	var v6202 int32
	_ = v6202
	var v6320 int32
	_ = v6320
	var v6323 int32
	_ = v6323
	var v6324 int32
	_ = v6324
	var v6327 int32
	_ = v6327
	var v6329 int32
	_ = v6329
	var v6332 int32
	_ = v6332
	var v6334 int32
	_ = v6334
	var v6335 int32
	_ = v6335
	var v6336 int32
	_ = v6336
	var v6341 int32
	_ = v6341
	var v6343 base.V128
	_ = v6343
	var v6344 base.V128
	_ = v6344
	var v6345 base.V128
	_ = v6345
	var v6346 base.V128
	_ = v6346
	var v6347 base.V128
	_ = v6347
	var v6348 base.V128
	_ = v6348
	var v6349 base.V128
	_ = v6349
	var v6350 base.V128
	_ = v6350
	var v6351 base.V128
	_ = v6351
	var v6356 int32
	_ = v6356
	var v6373 int32
	_ = v6373
	var v6380 int32
	_ = v6380
	var v6382 int32
	_ = v6382
	var v6442 int32
	_ = v6442
	var v6447 int32
	_ = v6447
	var v6457 int32
	_ = v6457
	var v6505 int32
	_ = v6505
	var v6506 int32
	_ = v6506
	var v6607 int32
	_ = v6607
	var v6608 base.V128
	_ = v6608
	var v6609 base.V128
	_ = v6609
	var v6614 base.V128
	_ = v6614
	var v6615 base.V128
	_ = v6615
	var v6620 base.V128
	_ = v6620
	var v6621 base.V128
	_ = v6621
	var v6625 base.V128
	_ = v6625
	var v6626 base.V128
	_ = v6626
	var v6628 base.V128
	_ = v6628
	var v6633 base.V128
	_ = v6633
	var v6638 base.V128
	_ = v6638
	var v6642 base.V128
	_ = v6642
	var v6643 base.V128
	_ = v6643
	var v6645 base.V128
	_ = v6645
	var v6648 base.V128
	_ = v6648
	var v6649 base.V128
	_ = v6649
	var v6651 base.V128
	_ = v6651
	var v6653 base.V128
	_ = v6653
	var v6669 base.V128
	_ = v6669
	var v6677 base.V128
	_ = v6677
	var v6682 base.V128
	_ = v6682
	var v6696 int32
	_ = v6696
	var v6714 int32
	_ = v6714
	var v6832 int32
	_ = v6832
	var v6833 int32
	_ = v6833
	var v6935 int32
	_ = v6935
	var v6939 int32
	_ = v6939
	var v6943 int32
	_ = v6943
	var v6947 int32
	_ = v6947
	var v6957 int32
	_ = v6957
	var v6967 int32
	_ = v6967
	var v6977 int32
	_ = v6977
	var v6982 int32
	_ = v6982
	var v7105 int32
	_ = v7105
	var v7107 int32
	_ = v7107
	var v7110 int32
	_ = v7110
	var v7111 int32
	_ = v7111
	var v7112 int32
	_ = v7112
	var v7115 int32
	_ = v7115
	var v7121 int32
	_ = v7121
	var v7123 int32
	_ = v7123
	var v7124 base.V128
	_ = v7124
	var v7125 int32
	_ = v7125
	var v7126 base.V128
	_ = v7126
	var v7127 base.V128
	_ = v7127
	var v7128 base.V128
	_ = v7128
	var v7129 base.V128
	_ = v7129
	var v7130 base.V128
	_ = v7130
	var v7131 base.V128
	_ = v7131
	var v7132 base.V128
	_ = v7132
	var v7133 base.V128
	_ = v7133
	var v7137 int32
	_ = v7137
	var v7148 int32
	_ = v7148
	var v7163 int32
	_ = v7163
	var v7165 int32
	_ = v7165
	var v7217 int32
	_ = v7217
	var v7224 int32
	_ = v7224
	var v7225 int32
	_ = v7225
	var v7240 int32
	_ = v7240
	var v7291 int32
	_ = v7291
	var v7393 int32
	_ = v7393
	var v7394 base.V128
	_ = v7394
	var v7395 base.V128
	_ = v7395
	var v7400 base.V128
	_ = v7400
	var v7401 base.V128
	_ = v7401
	var v7406 base.V128
	_ = v7406
	var v7407 base.V128
	_ = v7407
	var v7410 base.V128
	_ = v7410
	var v7411 int32
	_ = v7411
	var v7417 base.V128
	_ = v7417
	var v7419 base.V128
	_ = v7419
	var v7421 base.V128
	_ = v7421
	var v7426 base.V128
	_ = v7426
	var v7431 base.V128
	_ = v7431
	var v7434 base.V128
	_ = v7434
	var v7443 base.V128
	_ = v7443
	var v7444 base.V128
	_ = v7444
	var v7450 base.V128
	_ = v7450
	var v7452 base.V128
	_ = v7452
	var v7462 base.V128
	_ = v7462
	var v7476 base.V128
	_ = v7476
	var v7491 base.V128
	_ = v7491
	var v7495 int32
	_ = v7495
	var v7513 int32
	_ = v7513
	var v7633 int32
	_ = v7633
	var v7635 int32
	_ = v7635
	var v7736 int32
	_ = v7736
	var v7740 int32
	_ = v7740
	var v7744 int32
	_ = v7744
	var v7747 int32
	_ = v7747
	var v7750 int32
	_ = v7750
	var v7753 int32
	_ = v7753
	var v7763 int32
	_ = v7763
	var v7766 int32
	_ = v7766
	var v7769 int32
	_ = v7769
	var v7774 int32
	_ = v7774
	var v7896 int32
	_ = v7896
	var v7898 int32
	_ = v7898
	var v7916 int32
	_ = v7916
	var v8039 int32
	_ = v8039
	v17 = int32(0)
	if l10 == v17 {
		v8039 = v17
		goto L1
	} else {
		goto L2
	}
L1:
	;
	return v8039
L2:
	;
	if l8 == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L3
	}
L3:
	;
	if l6 == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L4
	}
L4:
	;
	if l2 == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L5
	}
L5:
	;
	if l1 == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L6
	}
L6:
	;
	if l0 == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L7
	}
L7:
	;
	if base.Ui32(l13+int32(-2147483647)) < base.Ui32(int32(-2147483646)) {
		v8039 = v17
		goto L1
	} else {
		goto L8
	}
L8:
	;
	if base.Ui32(l14+int32(-2147483647)) < base.Ui32(int32(-2147483646)) {
		v8039 = v17
		goto L1
	} else {
		goto L9
	}
L9:
	;
	if base.Ui32(int32(16)) < base.Ui32(l5) {
		v8039 = v17
		goto L1
	} else {
		goto L10
	}
L10:
	;
	if int32(1)<<(uint(l5)%32)&int32(_a_F_SharpYuvConvertWithOptions_0) == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L11
	}
L11:
	;
	if base.Ui32(int32(12)) < base.Ui32(l12) {
		v8039 = v17
		goto L1
	} else {
		goto L12
	}
L12:
	;
	if int32(1)<<(uint(l12)%32)&int32(_a_F_SharpYuvConvertWithOptions_1) == int32(0) {
		v8039 = v17
		goto L1
	} else {
		goto L13
	}
L13:
	;
	v154 = *(*int32)(unsafe.Add(mBase, uint32(l15)+4))
	v155 = *(*int32)(unsafe.Add(mBase, uint32(l15)))
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
		v8039 = v17
		goto L1
	} else {
		goto L16
	}
L16:
	;
	goto L14
L17:
	;
	v169 = m.G1
	v172 = *(*int32)(unsafe.Add(mBase, uint32(v169)+uint32(_c_F_SharpYuvConvertWithOptions[0])))
	v173 = m.G130
	v174 = *(*int32)(unsafe.Add(mBase, uint32(v173)))
	if v172 == v174 {
		goto L20
	} else {
		goto L21
	}
L18:
	;
	if (l9|l7|l11)&int32(1) != 0 {
		v8039 = v17
		goto L1
	} else {
		goto L19
	}
L19:
	;
	goto L17
L20:
	;
	v186 = int32(-1)<<(uint(l12)%32) ^ int32(-1)
	if l5 == l12 {
		goto L23
	} else {
		goto L24
	}
L21:
	;
	v176 = m.G1
	F_SharpYuvInitDsp(m)
	mBase = m.M
	F_SharpYuvInitGammaTables(m)
	mBase = m.M
	v181 = m.G130
	v182 = *(*int32)(unsafe.Add(mBase, uint32(v181)))
	*(*int32)(unsafe.Add(mBase, uint32(v176)+uint32(_c_F_SharpYuvConvertWithOptions[0]))) = v182
	goto L20
L22:
	;
	v252 = *(*int32)(unsafe.Add(mBase, uint32(v155)+44))
	v253 = *(*int32)(unsafe.Add(mBase, uint32(v155)+28))
	v254 = *(*int32)(unsafe.Add(mBase, uint32(v155)+12))
	v255 = int32(1)
	v256 = l14 + v255
	v257 = int32(-2)
	v258 = v256 & v257
	v261 = l13 + v255
	v263 = v261 & v257
	v265 = v263 * int32(3)
	v268 = base.I64_extend_i32_s(v265) << (uint(int64(2)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v268) {
		v273 = int32(0)
		goto L25
	} else {
		goto L26
	}
L23:
	;
	v232 = *(*int32)(unsafe.Add(mBase, uint32(v155)+40))
	v233 = *(*int32)(unsafe.Add(mBase, uint32(v155)+36))
	v234 = *(*int32)(unsafe.Add(mBase, uint32(v155)+32))
	v235 = *(*int32)(unsafe.Add(mBase, uint32(v155)+24))
	v236 = *(*int32)(unsafe.Add(mBase, uint32(v155)+20))
	v237 = *(*int32)(unsafe.Add(mBase, uint32(v155)+16))
	v238 = *(*int32)(unsafe.Add(mBase, uint32(v155)+8))
	v239 = *(*int32)(unsafe.Add(mBase, uint32(v155)+4))
	v240 = *(*int32)(unsafe.Add(mBase, uint32(v155)))
	v243 = v232
	v244 = v235
	v245 = v238
	v246 = v233
	v247 = v236
	v248 = v239
	v249 = v234
	v250 = v237
	v251 = v240
	goto L22
L24:
	;
	v188 = *(*int32)(unsafe.Add(mBase, uint32(v155)+40))
	v191 = int32(-1)
	v193 = int32(1) << (uint(l5+v191) % 32)
	v198 = v191<<(uint(l5)%32) ^ v191
	v199 = base.I32_div_s(v188*v186+v193, v198)
	v200 = *(*int32)(unsafe.Add(mBase, uint32(v155)+24))
	v203 = base.I32_div_s(v200*v186+v193, v198)
	v204 = *(*int32)(unsafe.Add(mBase, uint32(v155)+8))
	v207 = base.I32_div_s(v204*v186+v193, v198)
	v208 = *(*int32)(unsafe.Add(mBase, uint32(v155)+36))
	v211 = base.I32_div_s(v208*v186+v193, v198)
	v212 = *(*int32)(unsafe.Add(mBase, uint32(v155)+20))
	v215 = base.I32_div_s(v212*v186+v193, v198)
	v216 = *(*int32)(unsafe.Add(mBase, uint32(v155)+4))
	v219 = base.I32_div_s(v216*v186+v193, v198)
	v220 = *(*int32)(unsafe.Add(mBase, uint32(v155)+32))
	v223 = base.I32_div_s(v220*v186+v193, v198)
	v224 = *(*int32)(unsafe.Add(mBase, uint32(v155)+16))
	v227 = base.I32_div_s(v224*v186+v193, v198)
	v228 = *(*int32)(unsafe.Add(mBase, uint32(v155)))
	v231 = base.I32_div_s(v228*v186+v193, v198)
	v243 = v199
	v244 = v203
	v245 = v207
	v246 = v211
	v247 = v215
	v248 = v219
	v249 = v223
	v250 = v227
	v251 = v231
	goto L22
L25:
	;
	v274 = int32(0)
	v276 = base.I64_extend_i32_s(v263)
	v280 = v276 * base.I64_extend_i32_s(v258) << (uint(int64(1)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v280) {
		v286 = v274
		v287 = v274
		goto L28
	} else {
		goto L29
	}
L26:
	;
	v272 = F_dlmalloc(m, base.I32_wrap_i64(v268))
	mBase = m.M
	goto L27
L27:
	;
	v273 = v272
	goto L25
L28:
	;
	v289 = int32(1)
	v290 = v256 >> (uint(v289) % 32)
	v292 = v261 >> (uint(v289) % 32)
	v293 = int32(0)
	v296 = v276 << (uint(int64(2)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v296) {
		v301 = v293
		goto L32
	} else {
		goto L33
	}
L29:
	;
	v283 = base.I32_wrap_i64(v280)
	v284 = F_dlmalloc(m, v283)
	mBase = m.M
	goto L30
L30:
	;
	v285 = F_dlmalloc(m, v283)
	mBase = m.M
	goto L31
L31:
	;
	v286 = v284
	v287 = v285
	goto L28
L32:
	;
	v305 = v292 * int32(3)
	v306 = base.I64_extend_i32_s(v305)
	v309 = base.I64_extend_i32_s(v290) * v306 << (uint(int64(1)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v309) {
		v316 = v293
		v317 = int32(0)
		goto L35
	} else {
		goto L36
	}
L33:
	;
	v300 = F_dlmalloc(m, base.I32_wrap_i64(v296))
	mBase = m.M
	goto L34
L34:
	;
	v301 = v300
	goto L32
L35:
	;
	v320 = v306 << (uint(int64(1)) % 64)
	if base.Ui64(int64(4294967295)) < base.Ui64(v320) {
		v325 = int32(0)
		goto L39
	} else {
		goto L40
	}
L36:
	;
	v312 = base.I32_wrap_i64(v309)
	v313 = F_dlmalloc(m, v312)
	mBase = m.M
	goto L37
L37:
	;
	v314 = F_dlmalloc(m, v312)
	mBase = m.M
	goto L38
L38:
	;
	v316 = v313
	v317 = v314
	goto L35
L39:
	;
	v330 = base.F64_mul(base.F64_mul(base.F64_convert_i32_s(v263), float64(3)), base.F64_convert_i32_s(v258))
	if base.F64_lt(v330, float64(1.8446744073709552e+19))&base.F64_ge(v330, float64(0)) == int32(0) {
		goto L43
	} else {
		goto L44
	}
L40:
	;
	v324 = F_dlmalloc(m, base.I32_wrap_i64(v320))
	mBase = m.M
	goto L41
L41:
	;
	v325 = v324
	goto L39
L42:
	;
	v341 = int32(0)
	if v273 == v341 {
		v7916 = v341
		goto L45
	} else {
		goto L46
	}
L43:
	;
	v340 = int64(0)
	goto L42
L44:
	;
	v338 = base.I64_trunc_f64_u(v330)
	v340 = v338
	goto L42
L45:
	;
	F_dlfree(m, v286)
	mBase = m.M
	goto L455
L46:
	;
	if v301 == int32(0) {
		v7916 = v341
		goto L45
	} else {
		goto L47
	}
L47:
	;
	if v317 == int32(0) {
		v7916 = v341
		goto L45
	} else {
		goto L48
	}
L48:
	;
	if v287 == int32(0) {
		v7916 = v341
		goto L45
	} else {
		goto L49
	}
L49:
	;
	if v286 == int32(0) {
		v7916 = v341
		goto L45
	} else {
		goto L50
	}
L50:
	;
	if v316 == int32(0) {
		v7916 = v341
		goto L45
	} else {
		goto L51
	}
L51:
	;
	if v325 == int32(0) {
		v7916 = v341
		goto L45
	} else {
		goto L52
	}
L52:
	;
	if l5 < int32(13) {
		goto L53
	} else {
		goto L54
	}
L53:
	;
	v361 = int32(2)
	goto L55
L54:
	;
	v361 = int32(14) - l5
	goto L55
L55:
	;
	v362 = int32(0)
	v363 = base.B2i32(v361 < v362)
	v366 = v362 - v361
	v372 = v361 + l5
	if v362 < l14 {
		goto L57
	} else {
		goto L58
	}
L56:
	;
	if v361 < v362 {
		goto L243
	} else {
		goto L244
	}
L57:
	;
	v384 = int32(1)
	v386 = v273 + v265<<(uint(v384)%32)
	v388 = v263 << (uint(int32(2)) % 32)
	v389 = v386 + v388
	v391 = l4 << (uint(v384) % 32)
	v392 = int32(6)
	v399 = v263 << (uint(v384) % 32)
	v402 = v273 + v388
	v403 = v273 + v399
	v408 = v273 + v261<<(uint(int32(3))%32)&int32(-16)
	if v384 < v263 {
		goto L62
	} else {
		goto L63
	}
L58:
	;
	v376 = int32(1)
	if v376 < v263 {
		goto L59
	} else {
		goto L60
	}
L59:
	;
	v379 = v263
	goto L61
L60:
	;
	v379 = v376
	goto L61
L61:
	;
	v3499 = v379
	v3500 = v263 << (uint(int32(1)) % 32)
	v3501 = v263 << (uint(int32(2)) % 32)
	goto L56
L62:
	;
	v412 = v263
	goto L64
L63:
	;
	v412 = v384
	goto L64
L64:
	;
	v414 = v412 & int32(2147483640)
	v416 = base.B2i32(v263 < int32(8))
	v417 = l0
	v418 = l1
	v419 = l2
	v474 = v362
	v484 = v287 + v399
	v485 = v286 + v399
	v491 = v287
	v492 = v286
	v493 = v317
	v494 = v316
	goto L65
L65:
	;
	v548 = base.I32_div_s(l3, int32(2))
	if int32(8) < l5 {
		goto L68
	} else {
		goto L69
	}
L66:
	;
	v3499 = v412
	v3500 = v399
	v3501 = v388
	goto L56
L67:
	;
	if v474 == l14+int32(-1) {
		goto L136
	} else {
		goto L137
	}
L68:
	;
	v551 = v548
	goto L70
L69:
	;
	v551 = l3
	goto L70
L70:
	;
	v552 = int32(0)
	v554 = l13 + int32(1)
	v556 = v554 & int32(-2)
	v558 = v556 << (uint(int32(2)) % 32)
	if l5 != int32(8) {
		goto L72
	} else {
		goto L73
	}
L71:
	;
	if l13&int32(1) == int32(0) {
		goto L133
	} else {
		goto L134
	}
L72:
	;
	v752 = int32(1)
	v754 = base.B2i32(int32(31) < l13) & base.B2i32(v551 == v752)
	if v752 < l13 {
		goto L99
	} else {
		goto L100
	}
L73:
	;
	v561 = int32(1)
	if v561 < l13 {
		goto L74
	} else {
		goto L75
	}
L74:
	;
	v564 = l13
	goto L76
L75:
	;
	v564 = v561
	goto L76
L76:
	;
	if l13 < int32(48) {
		v685 = v552
		goto L77
	} else {
		goto L78
	}
L77:
	;
	v700 = int32(1)
	v710 = v273 + v685<<(uint(v700)%32)
	v711 = v551 * v685
	v722 = v564 - v685
	goto L96
L78:
	;
	if v551 != int32(1) {
		v685 = v552
		goto L77
	} else {
		goto L79
	}
L79:
	;
	v569 = int32(1)
	v573 = v273 + v554<<(uint(v569)%32)&int32(-4)
	v575 = v564 << (uint(v569) % 32)
	v576 = v573 + v575
	v578 = v273 + v575
	if base.B2i32(base.Ui32(v273) < base.Ui32(v576))&base.B2i32(base.Ui32(v573) < base.Ui32(v578)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L80
	}
L80:
	;
	v581 = v273 + v558
	v582 = v581 + v575
	if base.B2i32(base.Ui32(v273) < base.Ui32(v582))&base.B2i32(base.Ui32(v581) < base.Ui32(v578)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L81
	}
L81:
	;
	v586 = v417 + v564
	if base.B2i32(base.Ui32(v273) < base.Ui32(v586))&base.B2i32(base.Ui32(v417) < base.Ui32(v578)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L82
	}
L82:
	;
	v590 = v418 + v564
	if base.B2i32(base.Ui32(v273) < base.Ui32(v590))&base.B2i32(base.Ui32(v418) < base.Ui32(v578)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L83
	}
L83:
	;
	v594 = v419 + v564
	if base.B2i32(base.Ui32(v273) < base.Ui32(v594))&base.B2i32(base.Ui32(v419) < base.Ui32(v578)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L84
	}
L84:
	;
	if base.B2i32(base.Ui32(v573) < base.Ui32(v582))&base.B2i32(base.Ui32(v581) < base.Ui32(v576)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L85
	}
L85:
	;
	if base.B2i32(base.Ui32(v573) < base.Ui32(v586))&base.B2i32(base.Ui32(v417) < base.Ui32(v576)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L86
	}
L86:
	;
	if base.B2i32(base.Ui32(v573) < base.Ui32(v590))&base.B2i32(base.Ui32(v418) < base.Ui32(v576)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L87
	}
L87:
	;
	if base.B2i32(base.Ui32(v573) < base.Ui32(v594))&base.B2i32(base.Ui32(v419) < base.Ui32(v576)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L88
	}
L88:
	;
	if base.B2i32(base.Ui32(v581) < base.Ui32(v586))&base.B2i32(base.Ui32(v417) < base.Ui32(v582)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L89
	}
L89:
	;
	if base.B2i32(base.Ui32(v581) < base.Ui32(v590))&base.B2i32(base.Ui32(v418) < base.Ui32(v582)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L90
	}
L90:
	;
	if base.B2i32(base.Ui32(v581) < base.Ui32(v594))&base.B2i32(base.Ui32(v419) < base.Ui32(v582)) != 0 {
		v685 = v552
		goto L77
	} else {
		goto L91
	}
L91:
	;
	v624 = v564 & int32(2147483640)
	v628 = v419
	v629 = v273
	v638 = v624
	v639 = v417
	v640 = v418
	goto L92
L92:
	;
	v646 = int32(0)
	v647 = base.Simd_g_v128_load8x8_u(m, v639, v646)
	v648 = int32(2)
	v649 = base.Simd_g_i16x8_shl(v647, v648)
	base.Simd_g_v128_store(m, v629, v646, v649)
	v654 = base.Simd_g_v128_load8x8_u(m, v640, v646)
	v656 = base.Simd_g_i16x8_shl(v654, v648)
	base.Simd_g_v128_store(m, v629+v554<<(uint(int32(1))%32)&int32(-4), v646, v656)
	v661 = base.Simd_g_v128_load8x8_u(m, v628, v646)
	v663 = base.Simd_g_i16x8_shl(v661, v648)
	base.Simd_g_v128_store(m, v629+v558, v646, v663)
	v666 = int32(8)
	v675 = v638 + int32(-8)
	if v675 != 0 {
		v628 = v628 + v666
		v629 = v629 + int32(16)
		v638 = v675
		v639 = v639 + v666
		v640 = v640 + v666
		goto L92
	} else {
		goto L94
	}
L93:
	;
	if v564 == v624 {
		goto L71
	} else {
		goto L95
	}
L94:
	;
	goto L93
L95:
	;
	v685 = v624
	goto L77
L96:
	;
	v729 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v417+v711))))
	v730 = int32(2)
	v731 = v729 << (uint(v730) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v710))) = uint16(v731)
	v735 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v418+v711))))
	v737 = v735 << (uint(v730) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v710+v554<<(uint(v700)%32)&int32(-4)))) = uint16(v737)
	v741 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v419+v711))))
	v743 = v741 << (uint(v730) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v710+v558))) = uint16(v743)
	v749 = v722 + int32(-1)
	if v749 != 0 {
		v710 = v710 + v730
		v711 = v711 + v551
		v722 = v749
		goto L96
	} else {
		goto L98
	}
L98:
	;
	goto L71
L99:
	;
	v758 = l13
	goto L101
L100:
	;
	v758 = v752
	goto L101
L101:
	;
	if l5 < int32(13) {
		goto L104
	} else {
		goto L105
	}
L102:
	;
	v1108 = int32(1)
	v1124 = v273 + v1099<<(uint(v1108)%32)
	v1125 = v551 * v1099 << (uint(v1108) % 32)
	v1135 = v758 - v1099
	goto L130
L103:
	;
	v892 = int32(0)
	v893 = v892 - v764
	if v754 == v892 {
		v1025 = v892
		goto L117
	} else {
		goto L118
	}
L104:
	;
	v764 = int32(2)
	goto L106
L105:
	;
	v764 = int32(14) - l5
	goto L106
L106:
	;
	if v764 < int32(0) {
		goto L103
	} else {
		goto L107
	}
L107:
	;
	v767 = int32(0)
	if v754 == v767 {
		v1099 = v767
		goto L102
	} else {
		goto L108
	}
L108:
	;
	v770 = v558 + v273
	v773 = int32(1)
	v774 = v554 << (uint(v773) % 32)
	v777 = v774&int32(-4) + v273
	v782 = int32(2)
	v785 = int32(3)
	v787 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k0)
	v789 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	if base.Simd_g_v128_any_true(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v770-v419), v777-v417), v777-v418), v419-v777), v787), v789), base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v273-v417), v418-v273), v419-v273), v774), v787), v789))) != 0 {
		v1099 = v767
		goto L102
	} else {
		goto L109
	}
L109:
	;
	if base.Ui32(v770-v417) < base.Ui32(int32(16)) {
		v1099 = v767
		goto L102
	} else {
		goto L110
	}
L110:
	;
	if base.Ui32(v770-v418) < base.Ui32(int32(16)) {
		v1099 = v767
		goto L102
	} else {
		goto L111
	}
L111:
	;
	if base.Ui32(v558) < base.Ui32(int32(16)) {
		v1099 = v767
		goto L102
	} else {
		goto L112
	}
L112:
	;
	v818 = v758 & int32(2147483640)
	v822 = v419
	v823 = v273
	v832 = v417
	v833 = v418
	v835 = v818
	goto L113
L113:
	;
	v840 = int32(0)
	v841 = base.Simd_g_v128_load(m, v832, v840)
	v842 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)
	v846 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)
	v850 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v841, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)), v764), v789), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v841, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)), v764), v789))
	base.Simd_g_v128_store(m, v823, v840, v850)
	v855 = base.Simd_g_v128_load(m, v833, v840)
	v864 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v855, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)), v764), v789), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v855, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)), v764), v789))
	base.Simd_g_v128_store(m, v823+v554<<(uint(int32(1))%32)&int32(-4), v840, v864)
	v869 = base.Simd_g_v128_load(m, v822, v840)
	v878 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v869, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)), v764), v789), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v869, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)), v764), v789))
	base.Simd_g_v128_store(m, v823+v558, v840, v878)
	v881 = int32(16)
	v890 = v835 + int32(-8)
	if v890 != 0 {
		v822 = v822 + v881
		v823 = v823 + v881
		v832 = v832 + v881
		v833 = v833 + v881
		v835 = v890
		goto L113
	} else {
		goto L115
	}
L114:
	;
	if v758 != v818 {
		v1099 = v818
		goto L102
	} else {
		goto L116
	}
L115:
	;
	goto L114
L116:
	;
	goto L71
L117:
	;
	v1034 = int32(1)
	v1050 = v273 + v1025<<(uint(v1034)%32)
	v1051 = v551 * v1025 << (uint(v1034) % 32)
	v1061 = v758 - v1025
	goto L127
L118:
	;
	v897 = v558 + v273
	v900 = int32(1)
	v901 = v554 << (uint(v900) % 32)
	v904 = v901&int32(-4) + v273
	v909 = int32(2)
	v912 = int32(3)
	v914 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k0)
	v916 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	if base.Simd_g_v128_any_true(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v897-v419), v904-v417), v904-v418), v419-v904), v914), v916), base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v273-v417), v418-v273), v419-v273), v901), v914), v916))) != 0 {
		v1025 = v892
		goto L117
	} else {
		goto L119
	}
L119:
	;
	if base.Ui32(v897-v417) < base.Ui32(int32(16)) {
		v1025 = v892
		goto L117
	} else {
		goto L120
	}
L120:
	;
	if base.Ui32(v897-v418) < base.Ui32(int32(16)) {
		v1025 = v892
		goto L117
	} else {
		goto L121
	}
L121:
	;
	if base.Ui32(v558) < base.Ui32(int32(16)) {
		v1025 = v892
		goto L117
	} else {
		goto L122
	}
L122:
	;
	v945 = v758 & int32(2147483640)
	v949 = v419
	v950 = v273
	v959 = v417
	v960 = v418
	v962 = v945
	goto L123
L123:
	;
	v967 = int32(0)
	v968 = base.Simd_g_v128_load(m, v959, v967)
	v971 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)
	v975 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(v968), v893), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v968, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v893))
	base.Simd_g_v128_store(m, v950, v967, v975)
	v980 = base.Simd_g_v128_load(m, v960, v967)
	v987 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(v980), v893), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v980, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v893))
	base.Simd_g_v128_store(m, v950+v554<<(uint(int32(1))%32)&int32(-4), v967, v987)
	v992 = base.Simd_g_v128_load(m, v949, v967)
	v999 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(v992), v893), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v992, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v893))
	base.Simd_g_v128_store(m, v950+v558, v967, v999)
	v1002 = int32(16)
	v1011 = v962 + int32(-8)
	if v1011 != 0 {
		v949 = v949 + v1002
		v950 = v950 + v1002
		v959 = v959 + v1002
		v960 = v960 + v1002
		v962 = v1011
		goto L123
	} else {
		goto L125
	}
L124:
	;
	if v758 == v945 {
		goto L71
	} else {
		goto L126
	}
L125:
	;
	goto L124
L126:
	;
	v1025 = v945
	goto L117
L127:
	;
	v1069 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v417+v1051))))
	v1070 = int32(base.Ui32(v1069) >> (uint(v893) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1050))) = uint16(v1070)
	v1074 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v418+v1051))))
	v1075 = int32(base.Ui32(v1074) >> (uint(v893) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1050+v554<<(uint(v1034)%32)&int32(-4)))) = uint16(v1075)
	v1079 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v419+v1051))))
	v1080 = int32(base.Ui32(v1079) >> (uint(v893) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1050+v558))) = uint16(v1080)
	v1086 = v1061 + int32(-1)
	if v1086 != 0 {
		v1050 = v1050 + int32(2)
		v1051 = v1051 + v551<<(uint(v1034)%32)
		v1061 = v1086
		goto L127
	} else {
		goto L129
	}
L129:
	;
	goto L71
L130:
	;
	v1143 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v417+v1125))))
	v1144 = v1143 << (uint(v764) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1124))) = uint16(v1144)
	v1148 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v418+v1125))))
	v1149 = v1148 << (uint(v764) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1124+v554<<(uint(v1108)%32)&int32(-4)))) = uint16(v1149)
	v1153 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v419+v1125))))
	v1154 = v1153 << (uint(v764) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1124+v558))) = uint16(v1154)
	v1160 = v1135 + int32(-1)
	if v1160 != 0 {
		v1124 = v1124 + int32(2)
		v1125 = v1125 + v551<<(uint(v1108)%32)
		v1135 = v1160
		goto L130
	} else {
		goto L132
	}
L131:
	;
	goto L71
L132:
	;
	goto L131
L133:
	;
	goto L67
L134:
	;
	v1186 = int32(1)
	v1188 = v273 + l13<<(uint(v1186)%32)
	v1189 = int32(-2)
	v1191 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1188+v1189))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1188))) = uint16(v1191)
	v1195 = v1188 + v556<<(uint(v1186)%32)
	v1198 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1195+v1189))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1195))) = uint16(v1198)
	v1202 = v1188 + v554<<(uint(int32(2))%32)
	v1205 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1202+v1189))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1202))) = uint16(v1205)
	goto L133
L135:
	;
	if v263 < int32(8) {
		v2076 = int32(0)
		goto L207
	} else {
		goto L208
	}
L136:
	;
	v1889 = F_memcpy(m, v386, v273, v263*v392)
	mBase = m.M
	goto L135
L137:
	;
	v1210 = v417 + l4
	v1211 = v418 + l4
	v1212 = v419 + l4
	v1228 = base.I32_div_s(l3, int32(2))
	if int32(8) < l5 {
		goto L139
	} else {
		goto L140
	}
L138:
	;
	goto L135
L139:
	;
	v1231 = v1228
	goto L141
L140:
	;
	v1231 = l3
	goto L141
L141:
	;
	v1232 = int32(0)
	v1234 = l13 + int32(1)
	v1236 = v1234 & int32(-2)
	v1238 = v1236 << (uint(int32(2)) % 32)
	if l5 != int32(8) {
		goto L143
	} else {
		goto L144
	}
L142:
	;
	if l13&int32(1) == int32(0) {
		goto L204
	} else {
		goto L205
	}
L143:
	;
	v1432 = int32(1)
	v1434 = base.B2i32(int32(31) < l13) & base.B2i32(v1231 == v1432)
	if v1432 < l13 {
		goto L170
	} else {
		goto L171
	}
L144:
	;
	v1241 = int32(1)
	if v1241 < l13 {
		goto L145
	} else {
		goto L146
	}
L145:
	;
	v1244 = l13
	goto L147
L146:
	;
	v1244 = v1241
	goto L147
L147:
	;
	if l13 < int32(48) {
		v1365 = v1232
		goto L148
	} else {
		goto L149
	}
L148:
	;
	v1380 = int32(1)
	v1390 = v386 + v1365<<(uint(v1380)%32)
	v1391 = v1231 * v1365
	v1402 = v1244 - v1365
	goto L167
L149:
	;
	if v1231 != int32(1) {
		v1365 = v1232
		goto L148
	} else {
		goto L150
	}
L150:
	;
	v1249 = int32(1)
	v1253 = v386 + v1234<<(uint(v1249)%32)&int32(-4)
	v1255 = v1244 << (uint(v1249) % 32)
	v1256 = v1253 + v1255
	v1258 = v386 + v1255
	if base.B2i32(base.Ui32(v386) < base.Ui32(v1256))&base.B2i32(base.Ui32(v1253) < base.Ui32(v1258)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L151
	}
L151:
	;
	v1261 = v386 + v1238
	v1262 = v1261 + v1255
	if base.B2i32(base.Ui32(v386) < base.Ui32(v1262))&base.B2i32(base.Ui32(v1261) < base.Ui32(v1258)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L152
	}
L152:
	;
	v1266 = v1210 + v1244
	if base.B2i32(base.Ui32(v386) < base.Ui32(v1266))&base.B2i32(base.Ui32(v1210) < base.Ui32(v1258)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L153
	}
L153:
	;
	v1270 = v1211 + v1244
	if base.B2i32(base.Ui32(v386) < base.Ui32(v1270))&base.B2i32(base.Ui32(v1211) < base.Ui32(v1258)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L154
	}
L154:
	;
	v1274 = v1212 + v1244
	if base.B2i32(base.Ui32(v386) < base.Ui32(v1274))&base.B2i32(base.Ui32(v1212) < base.Ui32(v1258)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L155
	}
L155:
	;
	if base.B2i32(base.Ui32(v1253) < base.Ui32(v1262))&base.B2i32(base.Ui32(v1261) < base.Ui32(v1256)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L156
	}
L156:
	;
	if base.B2i32(base.Ui32(v1253) < base.Ui32(v1266))&base.B2i32(base.Ui32(v1210) < base.Ui32(v1256)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L157
	}
L157:
	;
	if base.B2i32(base.Ui32(v1253) < base.Ui32(v1270))&base.B2i32(base.Ui32(v1211) < base.Ui32(v1256)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L158
	}
L158:
	;
	if base.B2i32(base.Ui32(v1253) < base.Ui32(v1274))&base.B2i32(base.Ui32(v1212) < base.Ui32(v1256)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L159
	}
L159:
	;
	if base.B2i32(base.Ui32(v1261) < base.Ui32(v1266))&base.B2i32(base.Ui32(v1210) < base.Ui32(v1262)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L160
	}
L160:
	;
	if base.B2i32(base.Ui32(v1261) < base.Ui32(v1270))&base.B2i32(base.Ui32(v1211) < base.Ui32(v1262)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L161
	}
L161:
	;
	if base.B2i32(base.Ui32(v1261) < base.Ui32(v1274))&base.B2i32(base.Ui32(v1212) < base.Ui32(v1262)) != 0 {
		v1365 = v1232
		goto L148
	} else {
		goto L162
	}
L162:
	;
	v1304 = v1244 & int32(2147483640)
	v1308 = v1212
	v1309 = v386
	v1318 = v1304
	v1319 = v1210
	v1320 = v1211
	goto L163
L163:
	;
	v1326 = int32(0)
	v1327 = base.Simd_g_v128_load8x8_u(m, v1319, v1326)
	v1328 = int32(2)
	v1329 = base.Simd_g_i16x8_shl(v1327, v1328)
	base.Simd_g_v128_store(m, v1309, v1326, v1329)
	v1334 = base.Simd_g_v128_load8x8_u(m, v1320, v1326)
	v1336 = base.Simd_g_i16x8_shl(v1334, v1328)
	base.Simd_g_v128_store(m, v1309+v1234<<(uint(int32(1))%32)&int32(-4), v1326, v1336)
	v1341 = base.Simd_g_v128_load8x8_u(m, v1308, v1326)
	v1343 = base.Simd_g_i16x8_shl(v1341, v1328)
	base.Simd_g_v128_store(m, v1309+v1238, v1326, v1343)
	v1346 = int32(8)
	v1355 = v1318 + int32(-8)
	if v1355 != 0 {
		v1308 = v1308 + v1346
		v1309 = v1309 + int32(16)
		v1318 = v1355
		v1319 = v1319 + v1346
		v1320 = v1320 + v1346
		goto L163
	} else {
		goto L165
	}
L164:
	;
	if v1244 == v1304 {
		goto L142
	} else {
		goto L166
	}
L165:
	;
	goto L164
L166:
	;
	v1365 = v1304
	goto L148
L167:
	;
	v1409 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1210+v1391))))
	v1410 = int32(2)
	v1411 = v1409 << (uint(v1410) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1390))) = uint16(v1411)
	v1415 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1211+v1391))))
	v1417 = v1415 << (uint(v1410) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1390+v1234<<(uint(v1380)%32)&int32(-4)))) = uint16(v1417)
	v1421 = int32(*(*uint8)(unsafe.Add(mBase, uint32(v1212+v1391))))
	v1423 = v1421 << (uint(v1410) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1390+v1238))) = uint16(v1423)
	v1429 = v1402 + int32(-1)
	if v1429 != 0 {
		v1390 = v1390 + v1410
		v1391 = v1391 + v1231
		v1402 = v1429
		goto L167
	} else {
		goto L169
	}
L169:
	;
	goto L142
L170:
	;
	v1438 = l13
	goto L172
L171:
	;
	v1438 = v1432
	goto L172
L172:
	;
	if l5 < int32(13) {
		goto L175
	} else {
		goto L176
	}
L173:
	;
	v1788 = int32(1)
	v1804 = v386 + v1779<<(uint(v1788)%32)
	v1805 = v1231 * v1779 << (uint(v1788) % 32)
	v1815 = v1438 - v1779
	goto L201
L174:
	;
	v1572 = int32(0)
	v1573 = v1572 - v1444
	if v1434 == v1572 {
		v1705 = v1572
		goto L188
	} else {
		goto L189
	}
L175:
	;
	v1444 = int32(2)
	goto L177
L176:
	;
	v1444 = int32(14) - l5
	goto L177
L177:
	;
	if v1444 < int32(0) {
		goto L174
	} else {
		goto L178
	}
L178:
	;
	v1447 = int32(0)
	if v1434 == v1447 {
		v1779 = v1447
		goto L173
	} else {
		goto L179
	}
L179:
	;
	v1450 = v1238 + v386
	v1453 = int32(1)
	v1454 = v1234 << (uint(v1453) % 32)
	v1457 = v1454&int32(-4) + v386
	v1462 = int32(2)
	v1465 = int32(3)
	v1467 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k0)
	v1469 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	if base.Simd_g_v128_any_true(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v1450-v1212), v1457-v1210), v1457-v1211), v1212-v1457), v1467), v1469), base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v386-v1210), v1211-v386), v1212-v386), v1454), v1467), v1469))) != 0 {
		v1779 = v1447
		goto L173
	} else {
		goto L180
	}
L180:
	;
	if base.Ui32(v1450-v1210) < base.Ui32(int32(16)) {
		v1779 = v1447
		goto L173
	} else {
		goto L181
	}
L181:
	;
	if base.Ui32(v1450-v1211) < base.Ui32(int32(16)) {
		v1779 = v1447
		goto L173
	} else {
		goto L182
	}
L182:
	;
	if base.Ui32(v1238) < base.Ui32(int32(16)) {
		v1779 = v1447
		goto L173
	} else {
		goto L183
	}
L183:
	;
	v1498 = v1438 & int32(2147483640)
	v1502 = v1212
	v1503 = v386
	v1512 = v1210
	v1513 = v1211
	v1515 = v1498
	goto L184
L184:
	;
	v1520 = int32(0)
	v1521 = base.Simd_g_v128_load(m, v1512, v1520)
	v1522 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)
	v1526 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)
	v1530 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v1521, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)), v1444), v1469), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v1521, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)), v1444), v1469))
	base.Simd_g_v128_store(m, v1503, v1520, v1530)
	v1535 = base.Simd_g_v128_load(m, v1513, v1520)
	v1544 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v1535, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)), v1444), v1469), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v1535, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)), v1444), v1469))
	base.Simd_g_v128_store(m, v1503+v1234<<(uint(int32(1))%32)&int32(-4), v1520, v1544)
	v1549 = base.Simd_g_v128_load(m, v1502, v1520)
	v1558 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v1549, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k2)), v1444), v1469), base.Simd_g_v128_and(base.Simd_g_i32x4_shl(base.Simd_g_i8x16_swizzle_c(v1549, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k3)), v1444), v1469))
	base.Simd_g_v128_store(m, v1503+v1238, v1520, v1558)
	v1561 = int32(16)
	v1570 = v1515 + int32(-8)
	if v1570 != 0 {
		v1502 = v1502 + v1561
		v1503 = v1503 + v1561
		v1512 = v1512 + v1561
		v1513 = v1513 + v1561
		v1515 = v1570
		goto L184
	} else {
		goto L186
	}
L185:
	;
	if v1438 != v1498 {
		v1779 = v1498
		goto L173
	} else {
		goto L187
	}
L186:
	;
	goto L185
L187:
	;
	goto L142
L188:
	;
	v1714 = int32(1)
	v1730 = v386 + v1705<<(uint(v1714)%32)
	v1731 = v1231 * v1705 << (uint(v1714) % 32)
	v1741 = v1438 - v1705
	goto L198
L189:
	;
	v1577 = v1238 + v386
	v1580 = int32(1)
	v1581 = v1234 << (uint(v1580) % 32)
	v1584 = v1581&int32(-4) + v386
	v1589 = int32(2)
	v1592 = int32(3)
	v1594 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k0)
	v1596 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	if base.Simd_g_v128_any_true(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v1577-v1212), v1584-v1210), v1584-v1211), v1212-v1584), v1594), v1596), base.Simd_g_v128_and(base.Simd_g_i32x4_lt_u(base.Simd_g_i32x4_replace_lane_l3(base.Simd_g_i32x4_replace_lane_l2(base.Simd_g_i32x4_replace_lane_l1(base.Simd_g_i32x4_splat(v386-v1210), v1211-v386), v1212-v386), v1581), v1594), v1596))) != 0 {
		v1705 = v1572
		goto L188
	} else {
		goto L190
	}
L190:
	;
	if base.Ui32(v1577-v1210) < base.Ui32(int32(16)) {
		v1705 = v1572
		goto L188
	} else {
		goto L191
	}
L191:
	;
	if base.Ui32(v1577-v1211) < base.Ui32(int32(16)) {
		v1705 = v1572
		goto L188
	} else {
		goto L192
	}
L192:
	;
	if base.Ui32(v1238) < base.Ui32(int32(16)) {
		v1705 = v1572
		goto L188
	} else {
		goto L193
	}
L193:
	;
	v1625 = v1438 & int32(2147483640)
	v1629 = v1212
	v1630 = v386
	v1639 = v1210
	v1640 = v1211
	v1642 = v1625
	goto L194
L194:
	;
	v1647 = int32(0)
	v1648 = base.Simd_g_v128_load(m, v1639, v1647)
	v1651 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)
	v1655 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(v1648), v1573), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v1648, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v1573))
	base.Simd_g_v128_store(m, v1630, v1647, v1655)
	v1660 = base.Simd_g_v128_load(m, v1640, v1647)
	v1667 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(v1660), v1573), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v1660, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v1573))
	base.Simd_g_v128_store(m, v1630+v1234<<(uint(int32(1))%32)&int32(-4), v1647, v1667)
	v1672 = base.Simd_g_v128_load(m, v1629, v1647)
	v1679 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(v1672), v1573), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v1672, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v1573))
	base.Simd_g_v128_store(m, v1630+v1238, v1647, v1679)
	v1682 = int32(16)
	v1691 = v1642 + int32(-8)
	if v1691 != 0 {
		v1629 = v1629 + v1682
		v1630 = v1630 + v1682
		v1639 = v1639 + v1682
		v1640 = v1640 + v1682
		v1642 = v1691
		goto L194
	} else {
		goto L196
	}
L195:
	;
	if v1438 == v1625 {
		goto L142
	} else {
		goto L197
	}
L196:
	;
	goto L195
L197:
	;
	v1705 = v1625
	goto L188
L198:
	;
	v1749 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1210+v1731))))
	v1750 = int32(base.Ui32(v1749) >> (uint(v1573) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1730))) = uint16(v1750)
	v1754 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1211+v1731))))
	v1755 = int32(base.Ui32(v1754) >> (uint(v1573) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1730+v1234<<(uint(v1714)%32)&int32(-4)))) = uint16(v1755)
	v1759 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1212+v1731))))
	v1760 = int32(base.Ui32(v1759) >> (uint(v1573) % 32))
	*(*uint16)(unsafe.Add(mBase, uint32(v1730+v1238))) = uint16(v1760)
	v1766 = v1741 + int32(-1)
	if v1766 != 0 {
		v1730 = v1730 + int32(2)
		v1731 = v1731 + v1231<<(uint(v1714)%32)
		v1741 = v1766
		goto L198
	} else {
		goto L200
	}
L200:
	;
	goto L142
L201:
	;
	v1823 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1210+v1805))))
	v1824 = v1823 << (uint(v1444) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1804))) = uint16(v1824)
	v1828 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1211+v1805))))
	v1829 = v1828 << (uint(v1444) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1804+v1234<<(uint(v1788)%32)&int32(-4)))) = uint16(v1829)
	v1833 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1212+v1805))))
	v1834 = v1833 << (uint(v1444) % 32)
	*(*uint16)(unsafe.Add(mBase, uint32(v1804+v1238))) = uint16(v1834)
	v1840 = v1815 + int32(-1)
	if v1840 != 0 {
		v1804 = v1804 + int32(2)
		v1805 = v1805 + v1231<<(uint(v1788)%32)
		v1815 = v1840
		goto L201
	} else {
		goto L203
	}
L202:
	;
	goto L142
L203:
	;
	goto L202
L204:
	;
	goto L138
L205:
	;
	v1866 = int32(1)
	v1868 = v386 + l13<<(uint(v1866)%32)
	v1869 = int32(-2)
	v1871 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1868+v1869))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1868))) = uint16(v1871)
	v1875 = v1868 + v1236<<(uint(v1866)%32)
	v1878 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1875+v1869))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1875))) = uint16(v1878)
	v1882 = v1868 + v1234<<(uint(int32(2))%32)
	v1885 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v1882+v1869))))
	*(*uint16)(unsafe.Add(mBase, uint32(v1882))) = uint16(v1885)
	goto L204
L206:
	;
	if v263 < int32(8) {
		v2620 = int32(0)
		goto L217
	} else {
		goto L218
	}
L207:
	;
	v2192 = v2076 << (uint(int32(1)) % 32)
	v2195 = v412 - v2076
	goto L213
L208:
	;
	v1907 = int32(0)
	v1910 = v414
	goto L209
L209:
	;
	v2010 = int32(0)
	v2011 = base.Simd_g_v128_load(m, v273+v1907, v2010)
	v2013 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k5)
	v2017 = base.Simd_g_v128_load(m, v403+v1907, v2010)
	v2019 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k6)
	v2024 = base.Simd_g_v128_load(m, v402+v1907, v2010)
	v2026 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k7)
	v2029 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k8)
	v2031 = int32(16)
	v2033 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)
	v2050 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2011), v2013), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2017), v2019)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2024), v2026)), v2029), v2031), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v2011, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v2013), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v2017, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v2019)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v2024, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v2026)), v2029), v2031))
	base.Simd_g_v128_store(m, v492+v1907, v2010, v2050)
	v2056 = v1910 + int32(-8)
	if v2056 != 0 {
		v1907 = v1907 + v2031
		v1910 = v2056
		goto L209
	} else {
		goto L211
	}
L210:
	;
	if v412 == v414 {
		goto L206
	} else {
		goto L212
	}
L211:
	;
	goto L210
L212:
	;
	v2076 = v414
	goto L207
L213:
	;
	v2295 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v273+v2192))))
	v2299 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v403+v2192))))
	v2304 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v402+v2192))))
	v2311 = int64(base.Ui64(v2295*int64(13933)+v2299*int64(46871)+v2304*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64))
	*(*uint16)(unsafe.Add(mBase, uint32(v492+v2192))) = uint16(v2311)
	v2316 = v2195 + int32(-1)
	if v2316 != 0 {
		v2192 = v2192 + int32(2)
		v2195 = v2316
		goto L213
	} else {
		goto L215
	}
L214:
	;
	goto L206
L215:
	;
	goto L214
L216:
	;
	v2993 = int32(0)
	v2996 = v412
	goto L226
L217:
	;
	v2736 = v2620 << (uint(int32(1)) % 32)
	v2739 = v412 - v2620
	goto L223
L218:
	;
	v2449 = v485
	v2452 = v389
	v2463 = v408
	v2473 = v386
	v2519 = v414
	goto L219
L219:
	;
	v2550 = int32(0)
	v2551 = base.Simd_g_v128_load(m, v2473, v2550)
	v2553 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k5)
	v2556 = base.Simd_g_v128_load(m, v2463, v2550)
	v2558 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k6)
	v2562 = base.Simd_g_v128_load(m, v2452, v2550)
	v2564 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k7)
	v2567 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k8)
	v2569 = int32(16)
	v2571 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)
	v2588 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2551), v2553), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2556), v2558)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(v2562), v2564)), v2567), v2569), base.Simd_g_i32x4_shr_u(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v2551, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v2553), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v2556, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v2558)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v2562, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4))), v2564)), v2567), v2569))
	base.Simd_g_v128_store(m, v2449, v2550, v2588)
	v2600 = v2519 + int32(-8)
	if v2600 != 0 {
		v2449 = v2449 + v2569
		v2452 = v2452 + v2569
		v2463 = v2463 + v2569
		v2473 = v2473 + v2569
		v2519 = v2600
		goto L219
	} else {
		goto L221
	}
L220:
	;
	if v412 == v414 {
		goto L216
	} else {
		goto L222
	}
L221:
	;
	goto L220
L222:
	;
	v2620 = v414
	goto L217
L223:
	;
	v2839 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v386+v2736))))
	v2843 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v408+v2736))))
	v2848 = int64(*(*uint16)(unsafe.Add(mBase, uint32(v389+v2736))))
	v2855 = int64(base.Ui64(v2839*int64(13933)+v2843*int64(46871)+v2848*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64))
	*(*uint16)(unsafe.Add(mBase, uint32(v485+v2736))) = uint16(v2855)
	v2860 = v2739 + int32(-1)
	if v2860 != 0 {
		v2736 = v2736 + int32(2)
		v2739 = v2860
		goto L223
	} else {
		goto L225
	}
L224:
	;
	goto L216
L225:
	;
	goto L224
L226:
	;
	v3096 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v273+v2993))))
	v3097 = F_SharpYuvGammaToLinear(m, v3096, v372, v154)
	mBase = m.M
	v3102 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v403+v2993))))
	v3103 = F_SharpYuvGammaToLinear(m, v3102, v372, v154)
	mBase = m.M
	v3109 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v402+v2993))))
	v3110 = F_SharpYuvGammaToLinear(m, v3109, v372, v154)
	mBase = m.M
	v3120 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v3097)*int64(13933)+base.I64_extend_i32_u(v3103)*int64(46871)+base.I64_extend_i32_u(v3110)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v372, v154)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v491+v2993))) = uint16(v3120)
	v3125 = v2996 + int32(-1)
	if v3125 != 0 {
		v2993 = v2993 + int32(2)
		v2996 = v3125
		goto L226
	} else {
		goto L228
	}
L227:
	;
	v3141 = v386
	v3144 = v408
	v3155 = v389
	v3165 = v484
	v3211 = v412
	goto L229
L228:
	;
	goto L227
L229:
	;
	v3242 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3141))))
	v3243 = F_SharpYuvGammaToLinear(m, v3242, v372, v154)
	mBase = m.M
	v3247 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3144))))
	v3248 = F_SharpYuvGammaToLinear(m, v3247, v372, v154)
	mBase = m.M
	v3253 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3155))))
	v3254 = F_SharpYuvGammaToLinear(m, v3253, v372, v154)
	mBase = m.M
	v3264 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v3243)*int64(13933)+base.I64_extend_i32_u(v3248)*int64(46871)+base.I64_extend_i32_u(v3254)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v372, v154)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v3165))) = uint16(v3264)
	v3266 = int32(2)
	v3275 = v3211 + int32(-1)
	if v3275 != 0 {
		v3141 = v3141 + v3266
		v3144 = v3144 + v3266
		v3155 = v3155 + v3266
		v3165 = v3165 + v3266
		v3211 = v3275
		goto L229
	} else {
		goto L231
	}
L230:
	;
	if l5 < int32(13) {
		goto L233
	} else {
		goto L234
	}
L231:
	;
	goto L230
L232:
	;
	v3426 = F_memcpy(m, v494, v493, v292*v392)
	mBase = m.M
	v3435 = v305 << (uint(int32(1)) % 32)
	v3439 = v474 + int32(2)
	if v3439 < l14 {
		v417 = v417 + v391
		v418 = v418 + v391
		v419 = v419 + v391
		v474 = v3439
		v484 = v484 + v388
		v485 = v485 + v388
		v491 = v491 + v388
		v492 = v492 + v388
		v493 = v493 + v3435
		v494 = v3426 + v3435
		goto L65
	} else {
		goto L242
	}
L233:
	;
	v3294 = int32(2)
	goto L235
L234:
	;
	v3294 = int32(14) - l5
	goto L235
L235:
	;
	v3295 = v3294 + l5
	v3296 = int32(1)
	if v3296 < v292 {
		goto L236
	} else {
		goto L237
	}
L236:
	;
	v3299 = v292
	goto L238
L237:
	;
	v3299 = v3296
	goto L238
L238:
	;
	v3303 = v292 << (uint(int32(2)) % 32)
	v3306 = v292 << (uint(int32(3)) % 32)
	v3313 = v493
	v3314 = int32(0)
	v3317 = v3299
	goto L239
L239:
	;
	v3330 = v386 + v3314
	v3331 = int32(2)
	v3333 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3330+v3331))))
	v3334 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3330))))
	v3335 = v273 + v3314
	v3338 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3335+v3331))))
	v3339 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3335))))
	v3340 = F_SharpYuvGammaToLinear(m, v3339, v3295, v154)
	mBase = m.M
	v3341 = F_SharpYuvGammaToLinear(m, v3338, v3295, v154)
	mBase = m.M
	v3343 = F_SharpYuvGammaToLinear(m, v3334, v3295, v154)
	mBase = m.M
	v3345 = F_SharpYuvGammaToLinear(m, v3333, v3295, v154)
	mBase = m.M
	v3351 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v3340+v3341+v3343+v3345+v3331)>>(uint(v3331)%32)), v3295, v154)
	mBase = m.M
	v3352 = v386 + v3303 + v3314
	v3355 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3352+v3331))))
	v3356 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3352))))
	v3357 = v273 + v3303 + v3314
	v3360 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3357+v3331))))
	v3361 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3357))))
	v3362 = F_SharpYuvGammaToLinear(m, v3361, v3295, v154)
	mBase = m.M
	v3363 = F_SharpYuvGammaToLinear(m, v3360, v3295, v154)
	mBase = m.M
	v3365 = F_SharpYuvGammaToLinear(m, v3356, v3295, v154)
	mBase = m.M
	v3367 = F_SharpYuvGammaToLinear(m, v3355, v3295, v154)
	mBase = m.M
	v3373 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v3362+v3363+v3365+v3367+v3331)>>(uint(v3331)%32)), v3295, v154)
	mBase = m.M
	v3374 = v386 + v3306 + v3314
	v3377 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3374+v3331))))
	v3378 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3374))))
	v3379 = v273 + v3306 + v3314
	v3382 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3379+v3331))))
	v3390 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3379))))
	v3391 = F_SharpYuvGammaToLinear(m, v3390, v3295, v154)
	mBase = m.M
	v3392 = F_SharpYuvGammaToLinear(m, v3382, v3295, v154)
	mBase = m.M
	v3394 = F_SharpYuvGammaToLinear(m, v3378, v3295, v154)
	mBase = m.M
	v3396 = F_SharpYuvGammaToLinear(m, v3377, v3295, v154)
	mBase = m.M
	v3402 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v3391+v3392+v3394+v3396+v3331)>>(uint(v3331)%32)), v3295, v154)
	mBase = m.M
	v3411 = base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v3351)*int64(13933)+base.I64_extend_i32_u(v3373)*int64(46871)+base.I64_extend_i32_u(v3402)*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64)))
	v3412 = v3351 - v3411
	*(*uint16)(unsafe.Add(mBase, uint32(v3313))) = uint16(v3412)
	v3415 = v3373 - v3411
	*(*uint16)(unsafe.Add(mBase, uint32(v3313+v292<<(uint(int32(1))%32)))) = uint16(v3415)
	v3418 = v3402 - v3411
	*(*uint16)(unsafe.Add(mBase, uint32(v3313+v3303))) = uint16(v3418)
	v3425 = v3317 + int32(-1)
	if v3425 != 0 {
		v3313 = v3313 + v3331
		v3314 = v3314 + int32(4)
		v3317 = v3425
		goto L239
	} else {
		goto L241
	}
L240:
	;
	goto L232
L241:
	;
	goto L240
L242:
	;
	goto L66
L243:
	;
	v3557 = v252 >> (uint(v366) % 32)
	goto L245
L244:
	;
	v3557 = v252 << (uint(v361) % 32)
	goto L245
L245:
	;
	if v361 < v362 {
		goto L246
	} else {
		goto L247
	}
L246:
	;
	v3558 = v253 >> (uint(v366) % 32)
	goto L248
L247:
	;
	v3558 = v253 << (uint(v361) % 32)
	goto L248
L248:
	;
	if v361 < v362 {
		goto L249
	} else {
		goto L250
	}
L249:
	;
	v3559 = v254 >> (uint(v366) % 32)
	goto L251
L250:
	;
	v3559 = v254 << (uint(v361) % 32)
	goto L251
L251:
	;
	v3560 = int32(1)
	v3561 = v263 << (uint(v3560) % 32)
	v3562 = v273 + v3561
	v3563 = int32(-1)
	v3564 = v263 + v3563
	v3566 = v3564 << (uint(v3560) % 32)
	v3570 = v273 + v265<<(uint(v3560)%32)
	v3572 = v3562 + v3561
	v3574 = int32(2)
	v3576 = v3570 + v3561
	v3580 = v3576 + v3561
	v3591 = v3563 << (uint(v372) % 32)
	v3593 = v3591 ^ v3563
	v3595 = v3564 >> (uint(v3560) % 32)
	v3597 = int32(base.Ui32(v261) >> (uint(v3560) % 32))
	v3614 = (v292 + v3563) << (uint(v3560) % 32)
	v3630 = int32(0)
	v3650 = int64(-1)
	goto L252
L252:
	;
	__phi3768 = int64(0)
	__phi3772 = v316
	__phi3790 = v287
	__phi3796 = v316
	__phi3800 = v317
	__phi3804 = v316
	__phi3805 = int32(0)
	__phi3807 = v286
	v3768 = __phi3768
	v3772 = __phi3772
	v3790 = __phi3790
	v3796 = __phi3796
	v3800 = __phi3800
	v3804 = __phi3804
	v3805 = __phi3805
	v3807 = __phi3807
	goto L254
L253:
	;
	v4554 = int32(1)
	if v4554 < l14 {
		goto L353
	} else {
		goto L354
	}
L254:
	;
	v3849 = int32(0)
	v3851 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3804))))
	v3854 = int32(2)
	v3855 = v3851*int32(3) + v3854
	v3856 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3772))))
	v3860 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3807))))
	v3861 = (v3855+v3856)>>(uint(v3854)%32) + v3860
	if v3861 < v3849 {
		goto L256
	} else {
		goto L257
	}
L255:
	;
	if v3630 == int32(0) {
		v3630 = int32(1)
		v3650 = v4532
		goto L252
	} else {
		goto L348
	}
L256:
	;
	v3864 = v3849
	goto L258
L257:
	;
	v3864 = v3593
	goto L258
L258:
	;
	if v3861&v3591 != 0 {
		goto L259
	} else {
		goto L260
	}
L259:
	;
	v3866 = v3864
	goto L261
L260:
	;
	v3866 = v3861
	goto L261
L261:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v273))) = uint16(v3866)
	v3868 = int32(0)
	if v3805 < v258+int32(-2) {
		goto L262
	} else {
		goto L263
	}
L262:
	;
	v3871 = v305
	goto L264
L263:
	;
	v3871 = v3868
	goto L264
L264:
	;
	v3874 = v3804 + v3871<<(uint(int32(1))%32)
	v3875 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3874))))
	v3879 = v3807 + v3561
	v3880 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3879))))
	v3881 = (v3855+v3875)>>(uint(int32(2))%32) + v3880
	if v3881 < int32(0) {
		goto L265
	} else {
		goto L266
	}
L265:
	;
	v3884 = v3868
	goto L267
L266:
	;
	v3884 = v3593
	goto L267
L267:
	;
	if v3881&v3591 != 0 {
		goto L268
	} else {
		goto L269
	}
L268:
	;
	v3886 = v3884
	goto L270
L269:
	;
	v3886 = v3881
	goto L270
L270:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3570))) = uint16(v3886)
	v3888 = int32(2)
	v3889 = v3807 + v3888
	v3890 = m.G129
	v3891 = *(*int32)(unsafe.Add(mBase, uint32(v3890)))
	m.T0[v3891].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v3804, v3772, v3595, v3889, v273+v3574, v372)
	mBase = m.M
	v3894 = v3879 + v3888
	v3895 = *(*int32)(unsafe.Add(mBase, uint32(v3890)))
	m.T0[v3895].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v3804, v3874, v3595, v3894, v3570+v3574, v372)
	mBase = m.M
	v3897 = int32(0)
	v3899 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3804+v3614))))
	v3903 = v3899*int32(3) + v3888
	v3905 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3772+v3614))))
	v3909 = v3807 + v3566
	v3910 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3909))))
	v3911 = (v3903+v3905)>>(uint(v3888)%32) + v3910
	if v3911 < v3897 {
		goto L271
	} else {
		goto L272
	}
L271:
	;
	v3914 = v3897
	goto L273
L272:
	;
	v3914 = v3593
	goto L273
L273:
	;
	if v3911&v3591 != 0 {
		goto L274
	} else {
		goto L275
	}
L274:
	;
	v3916 = v3914
	goto L276
L275:
	;
	v3916 = v3911
	goto L276
L276:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v273+v3566))) = uint16(v3916)
	v3918 = int32(0)
	v3920 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3874+v3614))))
	v3924 = v3909 + v3561
	v3925 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3924))))
	v3926 = (v3903+v3920)>>(uint(int32(2))%32) + v3925
	if v3926 < v3918 {
		goto L277
	} else {
		goto L278
	}
L277:
	;
	v3929 = v3918
	goto L279
L278:
	;
	v3929 = v3593
	goto L279
L279:
	;
	if v3926&v3591 != 0 {
		goto L280
	} else {
		goto L281
	}
L280:
	;
	v3931 = v3929
	goto L282
L281:
	;
	v3931 = v3926
	goto L282
L282:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3570+v3566))) = uint16(v3931)
	v3933 = int32(0)
	v3935 = v292 << (uint(int32(1)) % 32)
	v3936 = v3804 + v3935
	v3937 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3936))))
	v3940 = int32(2)
	v3941 = v3937*int32(3) + v3940
	v3942 = v3772 + v3935
	v3943 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3942))))
	v3947 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3807))))
	v3948 = (v3941+v3943)>>(uint(v3940)%32) + v3947
	if v3948 < v3933 {
		goto L283
	} else {
		goto L284
	}
L283:
	;
	v3951 = v3933
	goto L285
L284:
	;
	v3951 = v3593
	goto L285
L285:
	;
	if v3948&v3591 != 0 {
		goto L286
	} else {
		goto L287
	}
L286:
	;
	v3953 = v3951
	goto L288
L287:
	;
	v3953 = v3948
	goto L288
L288:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3562))) = uint16(v3953)
	v3955 = int32(0)
	v3956 = v3874 + v3935
	v3957 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3956))))
	v3961 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3879))))
	v3962 = (v3941+v3957)>>(uint(int32(2))%32) + v3961
	if v3962 < v3955 {
		goto L289
	} else {
		goto L290
	}
L289:
	;
	v3965 = v3955
	goto L291
L290:
	;
	v3965 = v3593
	goto L291
L291:
	;
	if v3962&v3591 != 0 {
		goto L292
	} else {
		goto L293
	}
L292:
	;
	v3967 = v3965
	goto L294
L293:
	;
	v3967 = v3962
	goto L294
L294:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3576))) = uint16(v3967)
	v3969 = *(*int32)(unsafe.Add(mBase, uint32(v3890)))
	m.T0[v3969].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v3936, v3942, v3595, v3889, v3562+v3574, v372)
	mBase = m.M
	v3971 = *(*int32)(unsafe.Add(mBase, uint32(v3890)))
	m.T0[v3971].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v3936, v3956, v3595, v3894, v3576+v3574, v372)
	mBase = m.M
	v3973 = int32(0)
	v3975 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3936+v3614))))
	v3978 = int32(2)
	v3979 = v3975*int32(3) + v3978
	v3981 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3942+v3614))))
	v3985 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3909))))
	v3986 = (v3979+v3981)>>(uint(v3978)%32) + v3985
	if v3986 < v3973 {
		goto L295
	} else {
		goto L296
	}
L295:
	;
	v3989 = v3973
	goto L297
L296:
	;
	v3989 = v3593
	goto L297
L297:
	;
	if v3986&v3591 != 0 {
		goto L298
	} else {
		goto L299
	}
L298:
	;
	v3991 = v3989
	goto L300
L299:
	;
	v3991 = v3986
	goto L300
L300:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3562+v3566))) = uint16(v3991)
	v3993 = int32(0)
	v3995 = int32(*(*int16)(unsafe.Add(mBase, uint32(v3956+v3614))))
	v3999 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3924))))
	v4000 = (v3979+v3995)>>(uint(int32(2))%32) + v3999
	if v4000 < v3993 {
		goto L301
	} else {
		goto L302
	}
L301:
	;
	v4003 = v3993
	goto L303
L302:
	;
	v4003 = v3593
	goto L303
L303:
	;
	if v4000&v3591 != 0 {
		goto L304
	} else {
		goto L305
	}
L304:
	;
	v4005 = v4003
	goto L306
L305:
	;
	v4005 = v4000
	goto L306
L306:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3576+v3566))) = uint16(v4005)
	v4007 = int32(0)
	v4008 = v3936 + v3935
	v4009 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4008))))
	v4012 = int32(2)
	v4013 = v4009*int32(3) + v4012
	v4014 = v3942 + v3935
	v4015 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4014))))
	v4019 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3807))))
	v4020 = (v4013+v4015)>>(uint(v4012)%32) + v4019
	if v4020 < v4007 {
		goto L307
	} else {
		goto L308
	}
L307:
	;
	v4023 = v4007
	goto L309
L308:
	;
	v4023 = v3593
	goto L309
L309:
	;
	if v4020&v3591 != 0 {
		goto L310
	} else {
		goto L311
	}
L310:
	;
	v4025 = v4023
	goto L312
L311:
	;
	v4025 = v4020
	goto L312
L312:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3572))) = uint16(v4025)
	v4027 = int32(0)
	v4028 = v3956 + v3935
	v4029 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4028))))
	v4033 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3879))))
	v4034 = (v4013+v4029)>>(uint(int32(2))%32) + v4033
	if v4034 < v4027 {
		goto L313
	} else {
		goto L314
	}
L313:
	;
	v4037 = v4027
	goto L315
L314:
	;
	v4037 = v3593
	goto L315
L315:
	;
	if v4034&v3591 != 0 {
		goto L316
	} else {
		goto L317
	}
L316:
	;
	v4039 = v4037
	goto L318
L317:
	;
	v4039 = v4034
	goto L318
L318:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3580))) = uint16(v4039)
	v4041 = *(*int32)(unsafe.Add(mBase, uint32(v3890)))
	m.T0[v4041].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v4008, v4014, v3595, v3889, v3572+v3574, v372)
	mBase = m.M
	v4043 = *(*int32)(unsafe.Add(mBase, uint32(v3890)))
	m.T0[v4043].(func(*base.Module, int32, int32, int32, int32, int32, int32))(m, v4008, v4028, v3595, v3894, v3580+v3574, v372)
	mBase = m.M
	v4045 = int32(0)
	v4047 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4008+v3614))))
	v4050 = int32(2)
	v4051 = v4047*int32(3) + v4050
	v4053 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4014+v3614))))
	v4057 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3909))))
	v4058 = (v4051+v4053)>>(uint(v4050)%32) + v4057
	if v4058 < v4045 {
		goto L319
	} else {
		goto L320
	}
L319:
	;
	v4061 = v4045
	goto L321
L320:
	;
	v4061 = v3593
	goto L321
L321:
	;
	if v4058&v3591 != 0 {
		goto L322
	} else {
		goto L323
	}
L322:
	;
	v4063 = v4061
	goto L324
L323:
	;
	v4063 = v4058
	goto L324
L324:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3572+v3566))) = uint16(v4063)
	v4065 = int32(0)
	v4067 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4028+v3614))))
	v4071 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v3924))))
	v4072 = (v4051+v4067)>>(uint(int32(2))%32) + v4071
	if v4072 < v4065 {
		goto L325
	} else {
		goto L326
	}
L325:
	;
	v4075 = v4065
	goto L327
L326:
	;
	v4075 = v3593
	goto L327
L327:
	;
	if v4072&v3591 != 0 {
		goto L328
	} else {
		goto L329
	}
L328:
	;
	v4077 = v4075
	goto L330
L329:
	;
	v4077 = v4072
	goto L330
L330:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v3580+v3566))) = uint16(v4077)
	v4094 = v3849
	v4097 = v3499
	goto L331
L331:
	;
	v4197 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v273+v4094))))
	v4198 = F_SharpYuvGammaToLinear(m, v4197, v372, v154)
	mBase = m.M
	v4203 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v273+v3597<<(uint(v3574)%32)+v4094))))
	v4204 = F_SharpYuvGammaToLinear(m, v4203, v372, v154)
	mBase = m.M
	v4210 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v273+v3501+v4094))))
	v4211 = F_SharpYuvGammaToLinear(m, v4210, v372, v154)
	mBase = m.M
	v4221 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v4198)*int64(13933)+base.I64_extend_i32_u(v4204)*int64(46871)+base.I64_extend_i32_u(v4211)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v372, v154)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v301+v4094))) = uint16(v4221)
	v4226 = v4097 + int32(-1)
	if v4226 != 0 {
		v4094 = v4094 + int32(2)
		v4097 = v4226
		goto L331
	} else {
		goto L333
	}
L332:
	;
	v4242 = v3570
	v4245 = v273 + v3597<<(uint(int32(4))%32)
	v4256 = v301 + v3561
	v4266 = v3499
	goto L334
L333:
	;
	goto L332
L334:
	;
	v4343 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4242))))
	v4344 = F_SharpYuvGammaToLinear(m, v4343, v372, v154)
	mBase = m.M
	v4348 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4245))))
	v4349 = F_SharpYuvGammaToLinear(m, v4348, v372, v154)
	mBase = m.M
	v4355 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4242+v3501))))
	v4356 = F_SharpYuvGammaToLinear(m, v4355, v372, v154)
	mBase = m.M
	v4366 = F_SharpYuvLinearToGamma(m, base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v4344)*int64(13933)+base.I64_extend_i32_u(v4349)*int64(46871)+base.I64_extend_i32_u(v4356)*int64(4732)+int64(32768))>>(uint(int64(16))%64))), v372, v154)
	mBase = m.M
	*(*uint16)(unsafe.Add(mBase, uint32(v4256))) = uint16(v4366)
	v4368 = int32(2)
	v4375 = v4266 + int32(-1)
	if v4375 != 0 {
		v4242 = v4242 + v4368
		v4245 = v4245 + v4368
		v4256 = v4256 + v4368
		v4266 = v4375
		goto L334
	} else {
		goto L336
	}
L335:
	;
	if l5 < int32(13) {
		goto L338
	} else {
		goto L339
	}
L336:
	;
	goto L335
L337:
	;
	v4526 = m.G128
	v4527 = *(*int32)(unsafe.Add(mBase, uint32(v4526)))
	v4528 = m.G127
	v4529 = m.T0[v4527].(func(*base.Module, int32, int32, int32, int32, int32) int64)(m, v3790, v301, v3807, v3500, v372)
	mBase = m.M
	v4530 = *(*int32)(unsafe.Add(mBase, uint32(v4528)))
	m.T0[v4530].(func(*base.Module, int32, int32, int32, int32))(m, v3800, v325, v3796, v305)
	mBase = m.M
	v4532 = v4529 + v3768
	v4533 = int32(1)
	v4534 = v305 << (uint(v4533) % 32)
	v4537 = v3500 << (uint(v4533) % 32)
	v4542 = v3805 + int32(2)
	if v4542 < v258 {
		__phi3768 = v4532
		__phi3772 = v3804
		__phi3790 = v3790 + v4537
		__phi3796 = v3796 + v4534
		__phi3800 = v3800 + v4534
		__phi3804 = v3874
		__phi3805 = v4542
		__phi3807 = v3807 + v4537
		v3768 = __phi3768
		v3772 = __phi3772
		v3790 = __phi3790
		v3796 = __phi3796
		v3800 = __phi3800
		v3804 = __phi3804
		v3805 = __phi3805
		v3807 = __phi3807
		goto L254
	} else {
		goto L347
	}
L338:
	;
	v4394 = int32(2)
	goto L340
L339:
	;
	v4394 = int32(14) - l5
	goto L340
L340:
	;
	v4395 = v4394 + l5
	v4396 = int32(1)
	if v4396 < v292 {
		goto L341
	} else {
		goto L342
	}
L341:
	;
	v4399 = v292
	goto L343
L342:
	;
	v4399 = v4396
	goto L343
L343:
	;
	v4403 = v292 << (uint(int32(2)) % 32)
	v4406 = v292 << (uint(int32(3)) % 32)
	v4413 = v325
	v4414 = int32(0)
	v4417 = v4399
	goto L344
L344:
	;
	v4430 = v3570 + v4414
	v4431 = int32(2)
	v4433 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4430+v4431))))
	v4434 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4430))))
	v4435 = v273 + v4414
	v4438 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4435+v4431))))
	v4439 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4435))))
	v4440 = F_SharpYuvGammaToLinear(m, v4439, v4395, v154)
	mBase = m.M
	v4441 = F_SharpYuvGammaToLinear(m, v4438, v4395, v154)
	mBase = m.M
	v4443 = F_SharpYuvGammaToLinear(m, v4434, v4395, v154)
	mBase = m.M
	v4445 = F_SharpYuvGammaToLinear(m, v4433, v4395, v154)
	mBase = m.M
	v4451 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v4440+v4441+v4443+v4445+v4431)>>(uint(v4431)%32)), v4395, v154)
	mBase = m.M
	v4452 = v3570 + v4403 + v4414
	v4455 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4452+v4431))))
	v4456 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4452))))
	v4457 = v273 + v4403 + v4414
	v4460 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4457+v4431))))
	v4461 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4457))))
	v4462 = F_SharpYuvGammaToLinear(m, v4461, v4395, v154)
	mBase = m.M
	v4463 = F_SharpYuvGammaToLinear(m, v4460, v4395, v154)
	mBase = m.M
	v4465 = F_SharpYuvGammaToLinear(m, v4456, v4395, v154)
	mBase = m.M
	v4467 = F_SharpYuvGammaToLinear(m, v4455, v4395, v154)
	mBase = m.M
	v4473 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v4462+v4463+v4465+v4467+v4431)>>(uint(v4431)%32)), v4395, v154)
	mBase = m.M
	v4474 = v3570 + v4406 + v4414
	v4477 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4474+v4431))))
	v4478 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4474))))
	v4479 = v273 + v4406 + v4414
	v4482 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4479+v4431))))
	v4490 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v4479))))
	v4491 = F_SharpYuvGammaToLinear(m, v4490, v4395, v154)
	mBase = m.M
	v4492 = F_SharpYuvGammaToLinear(m, v4482, v4395, v154)
	mBase = m.M
	v4494 = F_SharpYuvGammaToLinear(m, v4478, v4395, v154)
	mBase = m.M
	v4496 = F_SharpYuvGammaToLinear(m, v4477, v4395, v154)
	mBase = m.M
	v4502 = F_SharpYuvLinearToGamma(m, int32(base.Ui32(v4491+v4492+v4494+v4496+v4431)>>(uint(v4431)%32)), v4395, v154)
	mBase = m.M
	v4511 = base.I32_wrap_i64(int64(base.Ui64(base.I64_extend_i32_u(v4451)*int64(13933)+base.I64_extend_i32_u(v4473)*int64(46871)+base.I64_extend_i32_u(v4502)*int64(4732)+int64(32768)) >> (uint(int64(16)) % 64)))
	v4512 = v4451 - v4511
	*(*uint16)(unsafe.Add(mBase, uint32(v4413))) = uint16(v4512)
	v4515 = v4473 - v4511
	*(*uint16)(unsafe.Add(mBase, uint32(v4413+v292<<(uint(int32(1))%32)))) = uint16(v4515)
	v4518 = v4502 - v4511
	*(*uint16)(unsafe.Add(mBase, uint32(v4413+v4403))) = uint16(v4518)
	v4525 = v4417 + int32(-1)
	if v4525 != 0 {
		v4413 = v4413 + v4431
		v4414 = v4414 + int32(4)
		v4417 = v4525
		goto L344
	} else {
		goto L346
	}
L345:
	;
	goto L337
L346:
	;
	goto L345
L347:
	;
	goto L255
L348:
	;
	if base.Ui64(v4532) < base.Ui64(v340) {
		goto L349
	} else {
		goto L350
	}
L349:
	;
	goto L253
L350:
	;
	if base.Ui64(v3650) < base.Ui64(v4532) {
		goto L349
	} else {
		goto L351
	}
L351:
	;
	if base.Ui32(v3630) < base.Ui32(int32(3)) {
		v3630 = v3630 + int32(1)
		v3650 = v4532
		goto L252
	} else {
		goto L352
	}
L352:
	;
	goto L349
L353:
	;
	v4557 = l14
	goto L355
L354:
	;
	v4557 = v4554
	goto L355
L355:
	;
	v4558 = int32(1)
	if v4558 < l13 {
		goto L356
	} else {
		goto L357
	}
L356:
	;
	v4561 = l13
	goto L358
L357:
	;
	v4561 = v4558
	goto L358
L358:
	;
	v4563 = v361 + int32(16)
	v4567 = int32(1) << (uint(v361+int32(15)) % 32)
	if int32(8) < l12 {
		goto L360
	} else {
		goto L361
	}
L359:
	;
	v6320 = int32(1)
	if v6320 < v290 {
		goto L403
	} else {
		goto L404
	}
L360:
	;
	v5282 = v4561 & int32(2147483640)
	v5287 = v3559 + v4567
	v5288 = base.Simd_g_i32x4_splat(v5287)
	v5289 = base.Simd_g_i32x4_splat(v186)
	v5290 = base.Simd_g_i32x4_splat(v245)
	v5291 = base.Simd_g_i32x4_splat(v248)
	v5292 = base.Simd_g_i32x4_splat(v251)
	v5293 = base.Simd_g_i32x4_splat(v263)
	v5294 = base.Simd_g_i32x4_splat(v292)
	v5304 = l6
	v5314 = v316
	v5368 = v286
	v5373 = int32(0)
	goto L381
L361:
	;
	v4570 = int32(1)
	v4573 = v4561 & int32(2147483646)
	v4577 = v261 << (uint(v4570) % 32) & int32(-4)
	v4592 = l6
	v4601 = v316
	v4661 = int32(0)
	v4671 = v286
	goto L362
L362:
	;
	if l13 < int32(2) {
		v4887 = int32(0)
		goto L365
	} else {
		goto L366
	}
L364:
	;
	if v4661&int32(1) != 0 {
		goto L377
	} else {
		goto L378
	}
L365:
	;
	v5006 = v4887
	v5008 = v4671 + v4887<<(uint(int32(1))%32)
	goto L371
L366:
	;
	v4722 = int32(0)
	v4724 = v4671
	goto L367
L367:
	;
	v4824 = int32(0)
	v4825 = base.Simd_g_v128_load16_splat(m, v4601+v4722, v4824)
	v4828 = base.Simd_g_v128_load32_zero(m, v4724, v4824)
	v4829 = base.Simd_g_i32x4_extend_low_i16x8_u(v4828)
	v4835 = base.Simd_g_v128_load16_splat(m, v4601+v292<<(uint(v4570)%32)+v4722, v4824)
	v4842 = base.Simd_g_v128_load16_splat(m, v4601+v4577+v4722, v4824)
	v4848 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(base.Simd_g_i32x4_splat(v251), base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v4825), v4829)), base.Simd_g_i32x4_splat(v4567)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_splat(v248), base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v4835), v4829))), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_splat(v245), base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v4842), v4829))), base.Simd_g_i32x4_splat(v3559)), v4563)
	v4852 = base.Simd_g_i8x16_swizzle_c(v4848, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k9))
	v4855 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k10)
	base.Simd_g_v128_store16_lane_l0(m, v4592+v4722, v4824, base.Simd_g_v128_bitselect(base.Simd_g_i8x16_swizzle_c(v4848, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k11)), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_gt_s(v4852, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k12)), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k10)), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_lt_u(v4852, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k13)), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k10))))
	v4868 = v4722 + int32(2)
	if v4573 != v4868 {
		v4722 = v4868
		v4724 = v4724 + int32(4)
		goto L367
	} else {
		goto L369
	}
L368:
	;
	if v4561 == v4573 {
		goto L364
	} else {
		goto L370
	}
L369:
	;
	goto L368
L370:
	;
	v4887 = v4573
	goto L365
L371:
	;
	v5110 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4601+v5006&int32(-2)))))
	v5111 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v5008))))
	v5115 = int32(1)
	v5116 = int32(base.Ui32(v5006) >> (uint(v5115) % 32))
	v5121 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4601+(v5116+v292)<<(uint(v5115)%32)))))
	v5129 = int32(*(*int16)(unsafe.Add(mBase, uint32(v4601+(v5116+v263)<<(uint(v5115)%32)))))
	v5134 = (v251*(v5110+v5111) + v4567 + v248*(v5121+v5111) + v245*(v5129+v5111) + v3559) >> (uint(v4563) % 32)
	if base.Ui32(v5134&int32(_a_F_SharpYuvConvertWithOptions_2)) < base.Ui32(int32(256)) {
		goto L373
	} else {
		goto L374
	}
L372:
	;
	goto L364
L373:
	;
	v5144 = v5134
	goto L375
L374:
	;
	v5144 = int32(base.Ui32(base.I32_extend16_s(v5134)^int32(-1)) >> (uint(int32(15)) % 32))
	goto L375
L375:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v4592+v5006))) = uint8(v5144)
	v5149 = v5006 + int32(1)
	if v4561 != v5149 {
		v5006 = v5149
		v5008 = v5008 + int32(2)
		goto L371
	} else {
		goto L376
	}
L376:
	;
	goto L372
L377:
	;
	v5273 = int32(3)
	goto L379
L378:
	;
	v5273 = int32(0)
	goto L379
L379:
	;
	v5275 = int32(1)
	v5279 = v4661 + v5275
	if v5279 != v4557 {
		v4592 = v4592 + l7
		v4601 = v4601 + v5273*v292<<(uint(v5275)%32)
		v4661 = v5279
		v4671 = v4671 + v4577
		goto L362
	} else {
		goto L380
	}
L380:
	;
	goto L359
L381:
	;
	if l13 < int32(8) {
		v5813 = int32(0)
		goto L384
	} else {
		goto L385
	}
L382:
	;
	goto L359
L383:
	;
	if v5373&int32(1) != 0 {
		goto L399
	} else {
		goto L400
	}
L384:
	;
	v5931 = v5813
	v5934 = v5813 << (uint(int32(1)) % 32)
	goto L390
L385:
	;
	v5433 = int32(0)
	v5436 = v5282
	v5499 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k14)
	v5500 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k15)
	goto L386
L386:
	;
	v5535 = int32(1)
	v5536 = base.Simd_g_i32x4_shr_u(v5500, v5535)
	v5537 = int32(3)
	v5542 = int32(2)
	v5552 = int32(0)
	v5558 = base.Simd_g_v128_load16_splat(m, v5314+base.Simd_g_i32x4_extract_lane_l0(v5536)<<(uint(v5535)%32), v5552)
	v5561 = base.Simd_g_v128_load16_lane_l1(m, v5314+base.Simd_g_i32x4_extract_lane_l1(v5536)<<(uint(v5535)%32), v5552, v5558)
	v5564 = base.Simd_g_v128_load16_lane_l2(m, v5314+base.Simd_g_i32x4_extract_lane_l2(v5536)<<(uint(v5535)%32), v5552, v5561)
	v5567 = base.Simd_g_v128_load16_lane_l3(m, v5314+base.Simd_g_i32x4_extract_lane_l3(v5536)<<(uint(v5535)%32), v5552, v5564)
	v5571 = base.Simd_g_v128_load(m, v5368+v5433, v5552)
	v5572 = base.Simd_g_i32x4_extend_low_i16x8_u(v5571)
	v5576 = base.Simd_g_i32x4_add(v5536, v5294)
	v5598 = base.Simd_g_v128_load16_splat(m, v5314+base.Simd_g_i32x4_extract_lane_l0(v5576)<<(uint(v5535)%32), v5552)
	v5601 = base.Simd_g_v128_load16_lane_l1(m, v5314+base.Simd_g_i32x4_extract_lane_l1(v5576)<<(uint(v5535)%32), v5552, v5598)
	v5604 = base.Simd_g_v128_load16_lane_l2(m, v5314+base.Simd_g_i32x4_extract_lane_l2(v5576)<<(uint(v5535)%32), v5552, v5601)
	v5607 = base.Simd_g_v128_load16_lane_l3(m, v5314+base.Simd_g_i32x4_extract_lane_l3(v5576)<<(uint(v5535)%32), v5552, v5604)
	v5612 = base.Simd_g_i32x4_add(v5536, v5293)
	v5634 = base.Simd_g_v128_load16_splat(m, v5314+base.Simd_g_i32x4_extract_lane_l0(v5612)<<(uint(v5535)%32), v5552)
	v5637 = base.Simd_g_v128_load16_lane_l1(m, v5314+base.Simd_g_i32x4_extract_lane_l1(v5612)<<(uint(v5535)%32), v5552, v5634)
	v5640 = base.Simd_g_v128_load16_lane_l2(m, v5314+base.Simd_g_i32x4_extract_lane_l2(v5612)<<(uint(v5535)%32), v5552, v5637)
	v5643 = base.Simd_g_v128_load16_lane_l3(m, v5314+base.Simd_g_i32x4_extract_lane_l3(v5612)<<(uint(v5535)%32), v5552, v5640)
	v5648 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v5288, base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v5567), v5572), v5292)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v5607), v5572), v5291)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v5643), v5572), v5290)), v4563)
	v5649 = int32(16)
	v5655 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	v5658 = base.Simd_g_i32x4_shr_u(v5499, v5535)
	v5680 = base.Simd_g_v128_load16_splat(m, v5314+base.Simd_g_i32x4_extract_lane_l0(v5658)<<(uint(v5535)%32), v5552)
	v5683 = base.Simd_g_v128_load16_lane_l1(m, v5314+base.Simd_g_i32x4_extract_lane_l1(v5658)<<(uint(v5535)%32), v5552, v5680)
	v5686 = base.Simd_g_v128_load16_lane_l2(m, v5314+base.Simd_g_i32x4_extract_lane_l2(v5658)<<(uint(v5535)%32), v5552, v5683)
	v5689 = base.Simd_g_v128_load16_lane_l3(m, v5314+base.Simd_g_i32x4_extract_lane_l3(v5658)<<(uint(v5535)%32), v5552, v5686)
	v5693 = base.Simd_g_i32x4_extend_low_i16x8_u(base.Simd_g_i8x16_swizzle_c(v5571, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v5697 = base.Simd_g_i32x4_add(v5658, v5294)
	v5719 = base.Simd_g_v128_load16_splat(m, v5314+base.Simd_g_i32x4_extract_lane_l0(v5697)<<(uint(v5535)%32), v5552)
	v5722 = base.Simd_g_v128_load16_lane_l1(m, v5314+base.Simd_g_i32x4_extract_lane_l1(v5697)<<(uint(v5535)%32), v5552, v5719)
	v5725 = base.Simd_g_v128_load16_lane_l2(m, v5314+base.Simd_g_i32x4_extract_lane_l2(v5697)<<(uint(v5535)%32), v5552, v5722)
	v5728 = base.Simd_g_v128_load16_lane_l3(m, v5314+base.Simd_g_i32x4_extract_lane_l3(v5697)<<(uint(v5535)%32), v5552, v5725)
	v5733 = base.Simd_g_i32x4_add(v5658, v5293)
	v5755 = base.Simd_g_v128_load16_splat(m, v5314+base.Simd_g_i32x4_extract_lane_l0(v5733)<<(uint(v5535)%32), v5552)
	v5758 = base.Simd_g_v128_load16_lane_l1(m, v5314+base.Simd_g_i32x4_extract_lane_l1(v5733)<<(uint(v5535)%32), v5552, v5755)
	v5761 = base.Simd_g_v128_load16_lane_l2(m, v5314+base.Simd_g_i32x4_extract_lane_l2(v5733)<<(uint(v5535)%32), v5552, v5758)
	v5764 = base.Simd_g_v128_load16_lane_l3(m, v5314+base.Simd_g_i32x4_extract_lane_l3(v5733)<<(uint(v5535)%32), v5552, v5761)
	v5769 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v5288, base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v5689), v5693), v5292)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v5728), v5693), v5291)), base.Simd_g_i32x4_mul(base.Simd_g_i32x4_add(base.Simd_g_i32x4_extend_low_i16x8_s(v5764), v5693), v5290)), v4563)
	v5778 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k16)
	v5779 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k8)
	v5787 = base.Simd_g_v128_bitselect(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_v128_bitselect(v5289, v5648, base.Simd_g_i32x4_gt_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v5648, v5649), v5649), v5289)), v5655), base.Simd_g_v128_and(base.Simd_g_v128_bitselect(v5289, v5769, base.Simd_g_i32x4_gt_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v5769, v5649), v5649), v5289)), v5655)), v5778, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_eq(base.Simd_g_v128_and(v5648, v5779), v5778), base.Simd_g_i32x4_eq(base.Simd_g_v128_and(v5769, v5779), v5778), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k17), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k18)))
	base.Simd_g_v128_store(m, v5304+v5433, v5552, v5787)
	v5790 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k19)
	v5796 = v5436 + int32(-8)
	if v5796 != 0 {
		v5433 = v5433 + v5649
		v5436 = v5796
		v5499 = base.Simd_g_i32x4_add(v5499, v5790)
		v5500 = base.Simd_g_i32x4_add(v5500, v5790)
		goto L386
	} else {
		goto L388
	}
L387:
	;
	if v4561 == v5282 {
		goto L383
	} else {
		goto L389
	}
L388:
	;
	goto L387
L389:
	;
	v5813 = v5282
	goto L384
L390:
	;
	v6037 = int32(*(*int16)(unsafe.Add(mBase, uint32(v5314+v5931&int32(-2)))))
	v6039 = int32(*(*uint16)(unsafe.Add(mBase, uint32(v5368+v5934))))
	v6043 = int32(1)
	v6044 = int32(base.Ui32(v5931) >> (uint(v6043) % 32))
	v6049 = int32(*(*int16)(unsafe.Add(mBase, uint32(v5314+(v6044+v292)<<(uint(v6043)%32)))))
	v6057 = int32(*(*int16)(unsafe.Add(mBase, uint32(v5314+(v6044+v263)<<(uint(v6043)%32)))))
	v6061 = (v5287 + (v6037+v6039)*v251 + (v6049+v6039)*v248 + (v6057+v6039)*v245) >> (uint(v4563) % 32)
	if v186 < base.I32_extend16_s(v6061) {
		goto L392
	} else {
		goto L393
	}
L391:
	;
	goto L383
L392:
	;
	v6064 = v186
	goto L394
L393:
	;
	v6064 = v6061
	goto L394
L394:
	;
	if v6061&int32(_a_F_SharpYuvConvertWithOptions_3) != 0 {
		goto L395
	} else {
		goto L396
	}
L395:
	;
	v6067 = int32(0)
	goto L397
L396:
	;
	v6067 = v6064
	goto L397
L397:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v5304+v5934))) = uint16(v6067)
	v6072 = v5931 + int32(1)
	if v4561 != v6072 {
		v5931 = v6072
		v5934 = v5934 + int32(2)
		goto L390
	} else {
		goto L398
	}
L398:
	;
	goto L391
L399:
	;
	v6196 = int32(3)
	goto L401
L400:
	;
	v6196 = int32(0)
	goto L401
L401:
	;
	v6198 = int32(1)
	v6202 = v5373 + v6198
	if v6202 != v4557 {
		v5304 = v5304 + l7
		v5314 = v5314 + v6196*v292<<(uint(v6198)%32)
		v5368 = v5368 + v261<<(uint(int32(1))%32)&int32(-4)
		v5373 = v6202
		goto L381
	} else {
		goto L402
	}
L402:
	;
	goto L382
L403:
	;
	v6323 = v290
	goto L405
L404:
	;
	v6323 = v6320
	goto L405
L405:
	;
	v6324 = int32(1)
	if v6324 < v292 {
		goto L406
	} else {
		goto L407
	}
L406:
	;
	v6327 = v292
	goto L408
L407:
	;
	v6327 = v6324
	goto L408
L408:
	;
	v6329 = v292 * int32(6)
	if int32(8) < l12 {
		goto L409
	} else {
		goto L410
	}
L409:
	;
	v7110 = v6327 & int32(2147483640)
	v7111 = int32(1)
	v7112 = v292 << (uint(v7111) % 32)
	v7115 = v6327 << (uint(v7111) % 32)
	v7121 = v261 << (uint(v7111) % 32) & int32(-4)
	v7123 = v3557 + v4567
	v7124 = base.Simd_g_i32x4_splat(v7123)
	v7125 = v3558 + v4567
	v7126 = base.Simd_g_i32x4_splat(v7125)
	v7127 = base.Simd_g_i32x4_splat(v186)
	v7128 = base.Simd_g_i32x4_splat(v243)
	v7129 = base.Simd_g_i32x4_splat(v246)
	v7130 = base.Simd_g_i32x4_splat(v249)
	v7131 = base.Simd_g_i32x4_splat(v244)
	v7132 = base.Simd_g_i32x4_splat(v247)
	v7133 = base.Simd_g_i32x4_splat(v250)
	v7137 = v6323 + int32(-1)
	v7148 = int32(0)
	v7163 = l8
	v7165 = l10
	v7217 = v7148
	v7224 = v316
	v7225 = v316 + v7112
	v7240 = v316 + v7121
	goto L430
L410:
	;
	v6332 = int32(0)
	v6334 = v6327 & int32(2147483640)
	v6335 = int32(1)
	v6336 = v292 << (uint(v6335) % 32)
	v6341 = v261 << (uint(v6335) % 32) & int32(-4)
	v6343 = base.Simd_g_i32x4_splat(v3557)
	v6344 = base.Simd_g_i32x4_splat(v243)
	v6345 = base.Simd_g_i32x4_splat(v246)
	v6346 = base.Simd_g_i32x4_splat(v249)
	v6347 = base.Simd_g_i32x4_splat(v3558)
	v6348 = base.Simd_g_i32x4_splat(v4567)
	v6349 = base.Simd_g_i32x4_splat(v244)
	v6350 = base.Simd_g_i32x4_splat(v247)
	v6351 = base.Simd_g_i32x4_splat(v250)
	v6356 = v6323 + int32(-1)
	v6373 = v6332
	v6380 = l8
	v6382 = l10
	v6442 = v316
	v6447 = v316 + v6341
	v6457 = v316 + v6336
	goto L411
L411:
	;
	if (base.B2i32(v292 < int32(8))|(base.B2i32(base.Ui32(l8) < base.Ui32(l10+v6327+v6356*l11))&base.B2i32(base.Ui32(l10) < base.Ui32(l8+v6327+v6356*l9))|base.B2i32(l9|l11 < v6332)))&v6335 != 0 {
		v6714 = int32(0)
		goto L414
	} else {
		goto L415
	}
L413:
	;
	v7105 = int32(1)
	v7107 = v6373 + v7105
	if v7107 != v6323 {
		v6373 = v7107
		v6380 = v6380 + l9
		v6382 = v6382 + l11
		v6442 = v6442 + v6329
		v6447 = v6447 + v6329
		v6457 = v6457 + v6329
		goto L411
	} else {
		goto L429
	}
L414:
	;
	v6832 = v6714
	v6833 = v6714 << (uint(int32(1)) % 32)
	goto L420
L415:
	;
	v6505 = int32(0)
	v6506 = v6442
	goto L416
L416:
	;
	v6607 = int32(0)
	v6608 = base.Simd_g_v128_load(m, v6506, v6607)
	v6609 = base.Simd_g_i32x4_extend_low_i16x8_s(v6608)
	v6614 = base.Simd_g_v128_load(m, v6506+v6336, v6607)
	v6615 = base.Simd_g_i32x4_extend_low_i16x8_s(v6614)
	v6620 = base.Simd_g_v128_load(m, v6506+v6341, v6607)
	v6621 = base.Simd_g_i32x4_extend_low_i16x8_s(v6620)
	v6625 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v6351, v6609), v6348), base.Simd_g_i32x4_mul(v6350, v6615)), base.Simd_g_i32x4_mul(v6349, v6621)), v6347), v4563)
	v6626 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)
	v6628 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_swizzle_c(v6608, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v6633 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_swizzle_c(v6614, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v6638 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_swizzle_c(v6620, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v6642 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v6351, v6628), v6348), base.Simd_g_i32x4_mul(v6350, v6633)), base.Simd_g_i32x4_mul(v6349, v6638)), v6347), v4563)
	v6643 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k20)
	v6645 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	v6648 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(v6625, v6645), base.Simd_g_v128_and(v6642, v6645))
	v6649 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k12)
	v6651 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k21)
	v6653 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k13)
	base.Simd_g_v128_store64_lane_l0(m, v6380+v6505, v6607, base.Simd_g_v128_bitselect(base.Simd_g_i8x16_shuffle2(v6625, v6642, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k22), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k23)), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_gt_s(v6648, v6649), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k21)), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_lt_u(v6648, v6653), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k21))))
	v6669 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v6346, v6609), v6348), base.Simd_g_i32x4_mul(v6345, v6615)), base.Simd_g_i32x4_mul(v6344, v6621)), v6343), v4563)
	v6677 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_mul(v6346, v6628), v6348), base.Simd_g_i32x4_mul(v6345, v6633)), base.Simd_g_i32x4_mul(v6344, v6638)), v6343), v4563)
	v6682 = base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(v6669, v6645), base.Simd_g_v128_and(v6677, v6645))
	base.Simd_g_v128_store64_lane_l0(m, v6382+v6505, v6607, base.Simd_g_v128_bitselect(base.Simd_g_i8x16_shuffle2(v6669, v6677, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k22), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k23)), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_gt_s(v6682, v6649), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k21)), base.Simd_g_i8x16_swizzle_c(base.Simd_g_i16x8_lt_u(v6682, v6653), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k21))))
	v6696 = v6505 + int32(8)
	if v6334 != v6696 {
		v6505 = v6696
		v6506 = v6506 + int32(16)
		goto L416
	} else {
		goto L418
	}
L417:
	;
	if v6327 == v6334 {
		goto L413
	} else {
		goto L419
	}
L418:
	;
	goto L417
L419:
	;
	v6714 = v6334
	goto L414
L420:
	;
	v6935 = int32(*(*int16)(unsafe.Add(mBase, uint32(v6442+v6833))))
	v6939 = int32(*(*int16)(unsafe.Add(mBase, uint32(v6457+v6833))))
	v6943 = int32(*(*int16)(unsafe.Add(mBase, uint32(v6447+v6833))))
	v6947 = (v250*v6935 + v4567 + v247*v6939 + v244*v6943 + v3558) >> (uint(v4563) % 32)
	if base.Ui32(v6947&int32(_a_F_SharpYuvConvertWithOptions_2)) < base.Ui32(int32(256)) {
		goto L422
	} else {
		goto L423
	}
L421:
	;
	goto L413
L422:
	;
	v6957 = v6947
	goto L424
L423:
	;
	v6957 = int32(base.Ui32(base.I32_extend16_s(v6947)^int32(-1)) >> (uint(int32(15)) % 32))
	goto L424
L424:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v6380+v6832))) = uint8(v6957)
	v6967 = (v249*v6935 + v4567 + v246*v6939 + v243*v6943 + v3557) >> (uint(v4563) % 32)
	if base.Ui32(v6967&int32(_a_F_SharpYuvConvertWithOptions_2)) < base.Ui32(int32(256)) {
		goto L425
	} else {
		goto L426
	}
L425:
	;
	v6977 = v6967
	goto L427
L426:
	;
	v6977 = int32(base.Ui32(base.I32_extend16_s(v6967)^int32(-1)) >> (uint(int32(15)) % 32))
	goto L427
L427:
	;
	*(*uint8)(unsafe.Add(mBase, uint32(v6382+v6832))) = uint8(v6977)
	v6982 = v6832 + int32(1)
	if v6327 != v6982 {
		v6832 = v6982
		v6833 = v6833 + int32(2)
		goto L420
	} else {
		goto L428
	}
L428:
	;
	goto L421
L429:
	;
	v7916 = v7105
	goto L45
L430:
	;
	if (base.B2i32(v292 < int32(8))|(base.B2i32(base.Ui32(l8) < base.Ui32(l10+v7137*l11+v7115))&base.B2i32(base.Ui32(l10) < base.Ui32(l8+v7137*l9+v7115))|base.B2i32(l9|l11 < v7148)))&v7111 != 0 {
		v7513 = int32(0)
		goto L433
	} else {
		goto L434
	}
L431:
	;
	v7916 = v7896
	goto L45
L432:
	;
	v7896 = int32(1)
	v7898 = v7217 + v7896
	if v7898 != v6323 {
		v7163 = v7163 + l9
		v7165 = v7165 + l11
		v7217 = v7898
		v7224 = v7224 + v6329
		v7225 = v7225 + v6329
		v7240 = v7240 + v6329
		goto L430
	} else {
		goto L454
	}
L433:
	;
	v7633 = v7513 << (uint(int32(1)) % 32)
	v7635 = v6327 - v7513
	goto L439
L434:
	;
	v7291 = int32(0)
	goto L435
L435:
	;
	v7393 = int32(0)
	v7394 = base.Simd_g_v128_load(m, v7224+v7291, v7393)
	v7395 = base.Simd_g_i32x4_extend_low_i16x8_s(v7394)
	v7400 = base.Simd_g_v128_load(m, v7224+v7112+v7291, v7393)
	v7401 = base.Simd_g_i32x4_extend_low_i16x8_s(v7400)
	v7406 = base.Simd_g_v128_load(m, v7224+v7121+v7291, v7393)
	v7407 = base.Simd_g_i32x4_extend_low_i16x8_s(v7406)
	v7410 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v7126, base.Simd_g_i32x4_mul(v7133, v7395)), base.Simd_g_i32x4_mul(v7132, v7401)), base.Simd_g_i32x4_mul(v7131, v7407)), v4563)
	v7411 = int32(16)
	v7417 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k1)
	v7419 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)
	v7421 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_swizzle_c(v7394, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v7426 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_swizzle_c(v7400, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v7431 = base.Simd_g_i32x4_extend_low_i16x8_s(base.Simd_g_i8x16_swizzle_c(v7406, base.Simd_g_const(&F_SharpYuvConvertWithOptions__k4)))
	v7434 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v7126, base.Simd_g_i32x4_mul(v7133, v7421)), base.Simd_g_i32x4_mul(v7132, v7426)), base.Simd_g_i32x4_mul(v7131, v7431)), v4563)
	v7443 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k16)
	v7444 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k8)
	v7450 = base.Simd_g_const(&F_SharpYuvConvertWithOptions__k24)
	v7452 = base.Simd_g_v128_bitselect(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_v128_bitselect(v7127, v7410, base.Simd_g_i32x4_gt_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v7410, v7411), v7411), v7127)), v7417), base.Simd_g_v128_and(base.Simd_g_v128_bitselect(v7127, v7434, base.Simd_g_i32x4_gt_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v7434, v7411), v7411), v7127)), v7417)), v7443, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_eq(base.Simd_g_v128_and(v7410, v7444), v7443), base.Simd_g_i32x4_eq(base.Simd_g_v128_and(v7434, v7444), v7443), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k17), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k18)))
	base.Simd_g_v128_store(m, v7163+v7291, v7393, v7452)
	v7462 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v7124, base.Simd_g_i32x4_mul(v7130, v7395)), base.Simd_g_i32x4_mul(v7129, v7401)), base.Simd_g_i32x4_mul(v7128, v7407)), v4563)
	v7476 = base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(base.Simd_g_i32x4_add(v7124, base.Simd_g_i32x4_mul(v7130, v7421)), base.Simd_g_i32x4_mul(v7129, v7426)), base.Simd_g_i32x4_mul(v7128, v7431)), v4563)
	v7491 = base.Simd_g_v128_bitselect(base.Simd_g_i16x8_narrow_i32x4_u(base.Simd_g_v128_and(base.Simd_g_v128_bitselect(v7127, v7462, base.Simd_g_i32x4_gt_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v7462, v7411), v7411), v7127)), v7417), base.Simd_g_v128_and(base.Simd_g_v128_bitselect(v7127, v7476, base.Simd_g_i32x4_gt_s(base.Simd_g_i32x4_shr_s(base.Simd_g_i32x4_shl(v7476, v7411), v7411), v7127)), v7417)), v7443, base.Simd_g_i8x16_shuffle2(base.Simd_g_i32x4_eq(base.Simd_g_v128_and(v7462, v7444), v7443), base.Simd_g_i32x4_eq(base.Simd_g_v128_and(v7476, v7444), v7443), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k17), base.Simd_g_const(&F_SharpYuvConvertWithOptions__k18)))
	base.Simd_g_v128_store(m, v7165+v7291, v7393, v7491)
	v7495 = v7291 + v7411
	if v7115&int32(-16) != v7495 {
		v7291 = v7495
		goto L435
	} else {
		goto L437
	}
L436:
	;
	if v6327 == v7110 {
		goto L432
	} else {
		goto L438
	}
L437:
	;
	goto L436
L438:
	;
	v7513 = v7110
	goto L433
L439:
	;
	v7736 = int32(*(*int16)(unsafe.Add(mBase, uint32(v7224+v7633))))
	v7740 = int32(*(*int16)(unsafe.Add(mBase, uint32(v7225+v7633))))
	v7744 = int32(*(*int16)(unsafe.Add(mBase, uint32(v7240+v7633))))
	v7747 = (v7125 + v250*v7736 + v247*v7740 + v244*v7744) >> (uint(v4563) % 32)
	if v186 < base.I32_extend16_s(v7747) {
		goto L441
	} else {
		goto L442
	}
L440:
	;
	goto L432
L441:
	;
	v7750 = v186
	goto L443
L442:
	;
	v7750 = v7747
	goto L443
L443:
	;
	if v7747&int32(_a_F_SharpYuvConvertWithOptions_3) != 0 {
		goto L444
	} else {
		goto L445
	}
L444:
	;
	v7753 = int32(0)
	goto L446
L445:
	;
	v7753 = v7750
	goto L446
L446:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7163+v7633))) = uint16(v7753)
	v7763 = (v7123 + v249*v7736 + v246*v7740 + v243*v7744) >> (uint(v4563) % 32)
	if v186 < base.I32_extend16_s(v7763) {
		goto L447
	} else {
		goto L448
	}
L447:
	;
	v7766 = v186
	goto L449
L448:
	;
	v7766 = v7763
	goto L449
L449:
	;
	if v7763&int32(_a_F_SharpYuvConvertWithOptions_3) != 0 {
		goto L450
	} else {
		goto L451
	}
L450:
	;
	v7769 = int32(0)
	goto L452
L451:
	;
	v7769 = v7766
	goto L452
L452:
	;
	*(*uint16)(unsafe.Add(mBase, uint32(v7165+v7633))) = uint16(v7769)
	v7774 = v7635 + int32(-1)
	if v7774 != 0 {
		v7633 = v7633 + int32(2)
		v7635 = v7774
		goto L439
	} else {
		goto L453
	}
L453:
	;
	goto L440
L454:
	;
	goto L431
L455:
	;
	F_dlfree(m, v316)
	mBase = m.M
	goto L456
L456:
	;
	F_dlfree(m, v287)
	mBase = m.M
	goto L457
L457:
	;
	F_dlfree(m, v317)
	mBase = m.M
	goto L458
L458:
	;
	F_dlfree(m, v301)
	mBase = m.M
	goto L459
L459:
	;
	F_dlfree(m, v325)
	mBase = m.M
	goto L460
L460:
	;
	F_dlfree(m, v273)
	mBase = m.M
	goto L461
L461:
	;
	v8039 = v7916
	goto L1
}

var F_SharpYuvConvertWithOptions__k0 = [2]uint64{0x1000000010, 0x1000000010}
var F_SharpYuvConvertWithOptions__k1 = [2]uint64{0xffff0000ffff, 0xffff0000ffff}
var F_SharpYuvConvertWithOptions__k2 = [2]uint64{0x100030201000100, 0x100070601000504}
var F_SharpYuvConvertWithOptions__k3 = [2]uint64{0x1000b0a01000908, 0x1000f0e01000d0c}
var F_SharpYuvConvertWithOptions__k4 = [2]uint64{0xf0e0d0c0b0a0908, 0x100010001000100}
var F_SharpYuvConvertWithOptions__k5 = [2]uint64{0x366d0000366d, 0x366d0000366d}
var F_SharpYuvConvertWithOptions__k6 = [2]uint64{0xb7170000b717, 0xb7170000b717}
var F_SharpYuvConvertWithOptions__k7 = [2]uint64{0x127c0000127c, 0x127c0000127c}
var F_SharpYuvConvertWithOptions__k8 = [2]uint64{0x800000008000, 0x800000008000}
var F_SharpYuvConvertWithOptions__k9 = [2]uint64{0x100010005040100, 0x100010001000100}
var F_SharpYuvConvertWithOptions__k10 = [2]uint64{0x200, 0x0}
var F_SharpYuvConvertWithOptions__k11 = [2]uint64{0x400, 0x0}
var F_SharpYuvConvertWithOptions__k12 = [2]uint64{0xffffffffffffffff, 0xffffffffffffffff}
var F_SharpYuvConvertWithOptions__k13 = [2]uint64{0x100010001000100, 0x100010001000100}
var F_SharpYuvConvertWithOptions__k14 = [2]uint64{0x500000004, 0x700000006}
var F_SharpYuvConvertWithOptions__k15 = [2]uint64{0x100000000, 0x300000002}
var F_SharpYuvConvertWithOptions__k16 = [2]uint64{0x0, 0x0}
var F_SharpYuvConvertWithOptions__k17 = [2]uint64{0xd0c090805040100, 0x8080808080808080}
var F_SharpYuvConvertWithOptions__k18 = [2]uint64{0x8080808080808080, 0xd0c090805040100}
var F_SharpYuvConvertWithOptions__k19 = [2]uint64{0x800000008, 0x800000008}
var F_SharpYuvConvertWithOptions__k20 = [2]uint64{0x1c1814100c080400, 0x0}
var F_SharpYuvConvertWithOptions__k21 = [2]uint64{0xe0c0a0806040200, 0x0}
var F_SharpYuvConvertWithOptions__k22 = [2]uint64{0x808080800c080400, 0x0}
var F_SharpYuvConvertWithOptions__k23 = [2]uint64{0xc08040080808080, 0x8080808080808080}
var F_SharpYuvConvertWithOptions__k24 = [2]uint64{0xd0c090805040100, 0x1d1c191815141110}
